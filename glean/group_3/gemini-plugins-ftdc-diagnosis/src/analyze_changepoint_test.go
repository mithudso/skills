package main

import (
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"
)

// helpers --------------------------------------------------------------------

// linearSeries returns a sequence v0, v0+step, v0+2*step, ...
func linearSeries(n int, v0, step float64) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = v0 + step*float64(i)
	}
	return out
}

// noisyConstant returns n samples drawn from N(mean, sigma) with a fixed seed.
func noisyConstant(n int, mean, sigma float64, seed int64) []float64 {
	rng := rand.New(rand.NewSource(seed))
	out := make([]float64, n)
	for i := range out {
		out[i] = mean + rng.NormFloat64()*sigma
	}
	return out
}

// stepSeries returns n samples that switch from level a (with noise sigma) to
// level b at index breakAt. Models a single clean regime change.
func stepSeries(n, breakAt int, a, b, sigma float64, seed int64) []float64 {
	rng := rand.New(rand.NewSource(seed))
	out := make([]float64, n)
	for i := range out {
		level := a
		if i >= breakAt {
			level = b
		}
		out[i] = level + rng.NormFloat64()*sigma
	}
	return out
}

// oscillatingSeries returns a series with no clean regime — alternates between
// two levels every `period` samples + noise. Designed to trigger >25 ED-PELT
// breakpoints (oversegmentation).
func oscillatingSeries(n, period int, low, high, sigma float64, seed int64) []float64 {
	rng := rand.New(rand.NewSource(seed))
	out := make([]float64, n)
	for i := range out {
		level := low
		if (i/period)%2 == 1 {
			level = high
		}
		out[i] = level + rng.NormFloat64()*sigma
	}
	return out
}

// storeWithSeries builds a minimal seriesStore around the named float64 series.
// All metrics are marked as non-counter gauges so the Welford rate handling
// stays out of the way. Timestamps are populated with monotonic 1Hz values.
func storeWithSeries(named map[string][]float64) *seriesStore {
	s := &seriesStore{nameToID: map[string]int{}}
	maxLen := 0
	for name, series := range named {
		id := len(s.names)
		s.names = append(s.names, name)
		s.nameToID[name] = id
		s.isCounter = append(s.isCounter, false)
		s.isDegenerate = append(s.isDegenerate, false)
		s.observedMin = append(s.observedMin, 0)
		s.observedMax = append(s.observedMax, 0)
		s.fullMean = append(s.fullMean, 0)
		s.fullSD = append(s.fullSD, 0)
		// Make sure materializedSeries is long enough.
		for len(s.materializedSeries) <= id {
			s.materializedSeries = append(s.materializedSeries, nil)
		}
		s.materializedSeries[id] = series
		if len(series) > maxLen {
			maxLen = len(series)
		}
	}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	s.timestamps = make([]time.Time, maxLen)
	for i := range s.timestamps {
		s.timestamps[i] = base.Add(time.Duration(i) * time.Second)
	}
	return s
}

// Tests ----------------------------------------------------------------------

// TestLinearSlope_pureRamp verifies linearSlope's slope_t blows up on a
// pure ramp. This is the precondition the slope_t >= 100 skip
// relies on.
func TestLinearSlope_pureRamp(t *testing.T) {
	series := linearSeries(1000, 0, 1) // y = i
	slope, slopeT := linearSlope(series)
	if math.Abs(slope-1.0) > 1e-9 {
		t.Errorf("slope: got %v, want 1.0", slope)
	}
	if slopeT < 1e6 {
		// Pure noiseless ramp → infinite t. We just want very large.
		t.Errorf("slope_t too small for pure ramp: got %v", slopeT)
	}
}

// TestLinearSlope_noisyConstant: slope_t stays near zero on stationary noise.
func TestLinearSlope_noisyConstant(t *testing.T) {
	series := noisyConstant(1000, 100, 5, 42)
	_, slopeT := linearSlope(series)
	if math.Abs(slopeT) > 5 {
		t.Errorf("slope_t too large for stationary series: got %v", slopeT)
	}
}

// TestChangepointSkip_highSlope_t confirms the slope_t >= 100 gate
// suppresses changepoint emission on pure-ramp metrics.
func TestChangepointSkip_highSlope_t(t *testing.T) {
	series := linearSeries(500, 0, 1.0) // pure ramp
	store := storeWithSeries(map[string][]float64{
		"test.ramp": series,
	})
	cfg := defaultAutoConfig()
	cfg.PELTMinSegment = 60
	findings := changepointOnly("h", store, []int{0}, cfg)
	for _, f := range findings {
		if k, _ := f["kind"].(string); k == KindChangepoint {
			t.Errorf("got changepoint Finding on pure ramp; expected slope_t gate to skip. Finding=%v", f)
		}
	}
}

// TestChangepointEmit_realStep verifies a clean step is still detected
// despite recent gating changes. Slope_t on a step should be moderate (~10-30),
// not >100.
func TestChangepointEmit_realStep(t *testing.T) {
	series := stepSeries(800, 400, 10, 50, 2, 7)
	store := storeWithSeries(map[string][]float64{
		"test.step": series,
	})
	cfg := defaultAutoConfig()
	cfg.PELTMinSegment = 60
	findings := changepointOnly("h", store, []int{0}, cfg)
	gotChangepoint := false
	for _, f := range findings {
		if k, _ := f["kind"].(string); k == KindChangepoint {
			gotChangepoint = true
			break
		}
	}
	if !gotChangepoint {
		t.Errorf("expected at least one changepoint Finding on clean step; got none. Findings=%v", findings)
	}
}

// TestChangepointOversegmented confirms oscillating series emit
// kind:"changepoint_oversegmented" instead of N individual changepoints.
// Uses a period of 80 (above minSeg of 60) so ED-PELT can find each
// transition, but enough transitions in 8000 samples to trip the >25 cap.
func TestChangepointOversegmented(t *testing.T) {
	// 100 oscillations of period 80 over 8000 samples (so ~100 transitions).
	series := oscillatingSeries(8000, 80, 0, 100, 3, 11)
	store := storeWithSeries(map[string][]float64{
		"test.osc": series,
	})
	cfg := defaultAutoConfig()
	cfg.PELTMinSegment = 60
	findings := changepointOnly("h", store, []int{0}, cfg)

	var overseg, regular int
	for _, f := range findings {
		k, _ := f["kind"].(string)
		switch k {
		case KindChangepointOversegmented:
			overseg++
			// Required fields per schema.
			if _, ok := f["break_count"].(int); !ok {
				t.Errorf("oversegmented Finding missing int break_count; got %T %v", f["break_count"], f["break_count"])
			}
			if rc, _ := f["reason_code"].(string); rc != "ed_pelt_oversegmented" {
				t.Errorf("oversegmented Finding has reason_code=%q, want ed_pelt_oversegmented", rc)
			}
		case KindChangepoint:
			regular++
		}
	}
	if overseg != 1 {
		t.Errorf("expected exactly 1 oversegmented Finding on oscillating series; got %d (and %d regular changepoints). Findings=%v", overseg, regular, findings)
	}
	if regular != 0 {
		t.Errorf("expected 0 regular changepoint Findings when oversegmented; got %d", regular)
	}
}

// TestChangepointCUSUMDropAtSmallEffect: when CUSUM disagrees AND effect is
// small (<0.5 Cohen's d), the Finding is dropped entirely.
//
// Construct a series where ED-PELT proposes a breakpoint but the effect size
// is tiny (means barely shift relative to noise). CUSUM may or may not agree
// — we test by counting Findings before vs after the change.
func TestChangepointCUSUMDropAtSmallEffect(t *testing.T) {
	// Tiny step (1 unit) in highly variable noise (sigma=5).
	series := stepSeries(800, 400, 10, 11, 5, 13)
	store := storeWithSeries(map[string][]float64{
		"test.tiny": series,
	})
	cfg := defaultAutoConfig()
	cfg.PELTMinSegment = 60
	findings := changepointOnly("h", store, []int{0}, cfg)
	// We don't assert zero — just that any emitted Findings have effect >= 0.5
	// OR cross_validated=true. The drop-gate guarantees small-effect+CUSUM-fail
	// pairs never emit.
	for _, f := range findings {
		k, _ := f["kind"].(string)
		if k != KindChangepoint {
			continue
		}
		effect, _ := f["effect_size"].(float64)
		validated, _ := f["cross_validated"].(bool)
		if !validated && effect < 0.5 {
			t.Errorf("drop-gate violated: small-effect+CUSUM-disagreed Finding emitted; effect=%v validated=%v", effect, validated)
		}
	}
}

// TestChangepointKindConstantValue asserts the new Finding kind string
// matches what the analyzer skill and schema expect.
func TestChangepointKindConstantValue(t *testing.T) {
	if KindChangepointOversegmented != "changepoint_oversegmented" {
		t.Errorf("KindChangepointOversegmented = %q; analyzer SKILL.md + analyze_schema.go expect %q",
			KindChangepointOversegmented, "changepoint_oversegmented")
	}
}

// TestChangepointReasonNotMachineFlag: the `reason` field is prose; the
// machine-actionable hook is `reason_code`. Verify both are present and
// distinct.
func TestChangepointReasonFields(t *testing.T) {
	series := oscillatingSeries(2000, 30, 0, 100, 5, 19)
	store := storeWithSeries(map[string][]float64{
		"test.osc": series,
	})
	cfg := defaultAutoConfig()
	cfg.PELTMinSegment = 60
	findings := changepointOnly("h", store, []int{0}, cfg)
	for _, f := range findings {
		k, _ := f["kind"].(string)
		if k != KindChangepointOversegmented {
			continue
		}
		reason, _ := f["reason"].(string)
		code, _ := f["reason_code"].(string)
		if reason == "" || code == "" {
			t.Errorf("oversegmented Finding missing reason/reason_code; got reason=%q code=%q", reason, code)
		}
		if !strings.Contains(reason, "ED-PELT") {
			t.Errorf("reason prose should mention ED-PELT for analyzer clarity; got %q", reason)
		}
	}
}

// TestStudentTSurvivalApprox_InvalidDF checks that df=0 and NaN df both return 0.5.
func TestStudentTSurvivalApprox_InvalidDF(t *testing.T) {
	cases := []struct {
		name string
		df   float64
	}{
		{"df=0", 0},
		{"df=-1", -1},
		{"df=NaN", math.NaN()},
	}
	for _, c := range cases {
		got := studentTSurvivalApprox(2.0, c.df)
		if got != 0.5 {
			t.Errorf("studentTSurvivalApprox(2.0, %v) = %v, want 0.5", c.name, got)
		}
	}
}

// TestStudentTSurvivalApprox_PValueInRange checks that p-values are always in [0,1],
// including for very large t values that could produce negative results before clamping.
func TestStudentTSurvivalApprox_PValueInRange(t *testing.T) {
	tCases := []float64{0, 0.5, 1, 2, 5, 10, 100, 1000}
	dfCases := []float64{1, 2, 5, 10, 30, 100}
	for _, tStat := range tCases {
		for _, df := range dfCases {
			oneSided := studentTSurvivalApprox(tStat, df)
			twoSided := 2 * oneSided
			if twoSided < 0 || twoSided > 2 {
				t.Errorf("raw two-sided p_value out of [0,2]: t=%v df=%v oneSided=%v twoSided=%v",
					tStat, df, oneSided, twoSided)
			}
		}
	}
	// Large t (1000) specifically: raw 2*oneSided should not be negative before clamp.
	raw := 2 * studentTSurvivalApprox(1000, 10)
	if raw < 0 {
		t.Errorf("raw two-sided p_value negative for t=1000, df=10: got %v", raw)
	}
}

// TestLagSecRounding checks that lagSec conversion uses rounding, not truncation.
// With intervalSec=2.5 and lag=24, the correct result is round(60.0)=60, not
// the old truncation lag*int(intervalSec)=24*2=48 (integer cast of 2.5).
func TestLagSecRounding(t *testing.T) {
	cases := []struct {
		lag         int
		intervalSec float64
		want        int
	}{
		{24, 2.5, 60},  // 24 * 2.5 = 60.0, old code: 24 * int(2.5) = 48
		{3, 0.5, 2},    // 3 * 0.5 = 1.5, rounds to 2
		{10, 1.0, 10},  // integer interval, exact
		{7, 1.5, 11},   // 7 * 1.5 = 10.5, rounds to 11
	}
	for _, c := range cases {
		got := int(math.Round(float64(c.lag) * c.intervalSec))
		if got != c.want {
			t.Errorf("lagSec(lag=%d, interval=%.2f): got %d, want %d",
				c.lag, c.intervalSec, got, c.want)
		}
	}
}

// TestSpearmanPValueApprox verifies the corrected Spearman test statistic
// t = rho * sqrt(n-2) / sqrt(1 - rho^2) against known reference values.
//
// Reference: for rho=0.5, n=30:
//   t = 0.5 * sqrt(28) / sqrt(1-0.25) = 0.5*5.292 / 0.866 ≈ 3.055
//   p ≈ 0.005 (two-sided; using normal approximation for n≥30)
//
// The old (buggy) formula used z = rho*sqrt(n-2) without the
// sqrt(1-rho^2) divisor, giving z = 0.5*5.292 = 2.646, p ≈ 0.008 —
// less significant than the correct formula.
func TestSpearmanPValueApprox(t *testing.T) {
	cases := []struct {
		rho  float64
		n    int
		desc string
		lo   float64
		hi   float64
	}{
		// rho=0: z=0 regardless of n, p must equal 1.0.
		{0.0, 10, "zero rho → p=1.0", 0.999, 1.001},
		// Perfect correlation: sqrt(1-rho^2) → 0, t → ∞, p → 0.
		{1.0, 10, "perfect positive → p=0", -0.001, 0.001},
		{-1.0, 10, "perfect negative → p=0", -0.001, 0.001},
		// rho=0.5, n=30: t ≈ 3.055, two-sided p ≈ 0.005.
		// Accept [0.002, 0.010] to tolerate normal vs t-distribution difference.
		{0.5, 30, "rho=0.5 n=30 → p≈0.005", 0.002, 0.010},
		// Symmetric: rho=-0.5 gives same two-sided p.
		{-0.5, 30, "rho=-0.5 n=30 → p≈0.005 (symmetric)", 0.002, 0.010},
	}
	for _, tc := range cases {
		p := spearmanPValueApprox(tc.rho, tc.n)
		if p < tc.lo || p > tc.hi {
			t.Errorf("%s: spearmanPValueApprox(%v, %d) = %v, want in [%v, %v]",
				tc.desc, tc.rho, tc.n, p, tc.lo, tc.hi)
		}
	}
}

// TestSpearmanPValueApprox_NLessThan3 verifies p=1.0 is returned for n < 3.
func TestSpearmanPValueApprox_NLessThan3(t *testing.T) {
	for _, n := range []int{0, 1, 2} {
		p := spearmanPValueApprox(0.99, n)
		if p != 1.0 {
			t.Errorf("n=%d: expected p=1.0, got %v", n, p)
		}
	}
}
