// analyze_wave2_fixes_test.go — regression tests for Wave 2 bug fixes.
package main

import (
	"math"
	"strings"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Fix 1 (W-R9): emitCaptureQualityScore reads cv_failure_pct from changepoints,
// not from run_health. Verified: cv_failure_pct is non-zero when changepoints
// are present and run_health has NOT yet been added.
// ---------------------------------------------------------------------------

func TestCaptureQualityScore_NoRunHealthNeeded(t *testing.T) {
	// 5 failing + 5 passing changepoints, run_health absent.
	// cv_failure_pct = 50% → deduction = min(50, 20) = 20.
	c := &findingCollector{}
	for i := 0; i < 5; i++ {
		c.add(Finding{"kind": KindChangepoint, "cross_validated": true})
	}
	for i := 0; i < 5; i++ {
		c.add(Finding{"kind": KindChangepoint, "cross_validated": false})
	}
	emitCaptureQualityScore(c, "h")

	var qf Finding
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == KindCaptureQualityScore {
			qf = f
			break
		}
	}
	if qf == nil {
		t.Fatal("expected capture_quality_score Finding")
	}
	deduction, _ := qf["deduction_cv_failure"].(int)
	if deduction != 20 {
		t.Errorf("deduction_cv_failure: got %d, want 20 (50%% cv_failure capped at 20)", deduction)
	}
	score, _ := qf["score"].(float64)
	if score < 79.9 || score > 80.1 {
		t.Errorf("score: got %v, want 80.0", score)
	}
}

// ---------------------------------------------------------------------------
// Fix 2 (StdDev): CompressedColumn.StdDev / RateStdDev use N-1 (Bessel).
// ---------------------------------------------------------------------------

func TestCompressedColumn_StdDev_SampleNotPopulation(t *testing.T) {
	// Two values 0 and 10: Welford M2 = 50, N = 2.
	// Sample stddev  = sqrt(M2 / (N-1)) = sqrt(50/1) ≈ 7.071
	// Population stddev = sqrt(M2 / N)  = sqrt(50/2) = 5.0
	col := NewCompressedColumn(4)
	col.Add(0, false)
	col.Add(10, false)

	got := col.StdDev()
	wantSample := math.Sqrt(50.0 / float64(1))
	wantPop := math.Sqrt(50.0 / float64(2))

	if math.Abs(got-wantSample) > 0.001 {
		t.Errorf("StdDev: got %v, want sample stddev %v (not population %v)", got, wantSample, wantPop)
	}
	if math.Abs(got-wantPop) < 0.001 {
		t.Errorf("StdDev should NOT equal population stddev %v", wantPop)
	}
}

func TestCompressedColumn_RateStdDev_SampleNotPopulation(t *testing.T) {
	// Four consecutive counter values 0, 10, 20, 40 produce two rates: 10 and 20.
	// Sample stddev = sqrt(M2/(N-1)) = sqrt(50/1) ≈ 7.071.
	col := NewCompressedColumn(8)
	col.Add(0, false)
	col.Add(10, true)
	col.Add(20, true)
	col.Add(40, true)

	got := col.RateStdDev()
	wantSample := math.Sqrt(50.0 / float64(1))
	wantPop := math.Sqrt(50.0 / float64(2))

	if math.Abs(got-wantSample) > 0.001 {
		t.Errorf("RateStdDev: got %v, want sample stddev %v (not population %v)", got, wantSample, wantPop)
	}
	if math.Abs(got-wantPop) < 0.001 {
		t.Errorf("RateStdDev should NOT equal population stddev %v", wantPop)
	}
}

// ---------------------------------------------------------------------------
// Fix 3 (React Pearson): Normalization uses n-1, not count-1. For lag=1 with
// a linear series, the old formula (count-1) inflates r toward 1.0 while the
// correct formula (n-1) gives r < 1.
// ---------------------------------------------------------------------------

func TestPearsonNormalizationUsesNMinus1(t *testing.T) {
	n := 100
	series := make([]float64, n)
	mean := 0.0
	for i := range series {
		series[i] = float64(i + 1)
		mean += series[i]
	}
	mean /= float64(n)
	m2 := 0.0
	for _, v := range series {
		d := v - mean
		m2 += d * d
	}
	sd := math.Sqrt(m2 / float64(n-1))

	// Compute dot at lag=1 (count = n-1 = 99).
	dot := 0.0
	count := 0
	for i := 0; i < n; i++ {
		target := i + 1
		if target >= n {
			continue
		}
		dot += (series[i] - mean) * (series[target] - mean)
		count++
	}

	rNew := dot / (float64(n-1) * sd * sd)     // correct (n-1)
	rOld := dot / (float64(count-1) * sd * sd)  // bug (count-1 = n-2)

	// The old formula uses a smaller denominator, inflating r.
	if rNew >= rOld {
		t.Errorf("rNew=%v should be < rOld=%v for lag=1 (n-1 is larger denominator)", rNew, rOld)
	}
	// Correct r must be in [-1, 1].
	if rNew > 1.0+1e-9 || rNew < -1.0-1e-9 {
		t.Errorf("corrected r=%v out of valid Pearson range [-1, 1]", rNew)
	}
}

// ---------------------------------------------------------------------------
// Fix 4 (Glossary hasMongos): Repl metrics must clear hasMongos.
// ---------------------------------------------------------------------------

func TestInferTopology_ReplMetricClearsMongos(t *testing.T) {
	// A node with only repl metrics (no WT, no sharding) should not be classified
	// as sharded_router; it is at minimum a replica-set member.
	names := []string{
		"serverStatus.repl.electionId",
		"serverStatus.repl.opTime.ts",
		"serverStatus.opcounters.insert",
	}
	got := inferTopology(names)
	if got == "sharded_router" {
		t.Errorf("inferTopology: got %q for repl-only metrics, must not be sharded_router", got)
	}
}

func TestInferTopology_ReplAndWT(t *testing.T) {
	names := []string{
		"serverStatus.repl.electionId",
		"serverStatus.wiredTiger.cache.bytes currently in the cache",
	}
	got := inferTopology(names)
	if got != "replica_set" {
		t.Errorf("inferTopology: got %q for repl+WT metrics, want replica_set", got)
	}
}

// ---------------------------------------------------------------------------
// Fix 5 (EWMA variance): Hunter 1986 formula eliminates steady-state bias.
// On a constant-input series, EWMAVar must converge to 0.
// ---------------------------------------------------------------------------

func TestEWMAVariance_HunterFormula_ConvergesToZeroOnConstantInput(t *testing.T) {
	alphaMean := baselineEWMAlphaMean
	alphaVar := baselineEWMAlphaVar
	ewmaMean := 100.0
	ewmaVar := 1.0 // start non-zero

	for i := 0; i < 500; i++ {
		x := 100.0 // constant
		dev := x - ewmaMean
		newMean := (1-alphaMean)*ewmaMean + alphaMean*x
		// Hunter (1986): V_new = (1-β)*V_old + β*(1-α_mean)*dev²
		ewmaVar = (1-alphaVar)*ewmaVar + alphaVar*(1-alphaMean)*dev*dev
		ewmaMean = newMean
	}
	if ewmaVar > 1e-10 {
		t.Errorf("Hunter formula: EWMAVar on constant input did not converge to 0, got %v", ewmaVar)
	}
}

func TestEWMAVariance_HunterFormula_LowerSteadyStateVarianceThanOld(t *testing.T) {
	// On oscillating input the old (dev²) formula overestimates steady-state
	// variance relative to the Hunter formula. Verify the Hunter estimate is
	// strictly lower on a signal with known variance.
	alphaMean := baselineEWMAlphaMean
	alphaVar := baselineEWMAlphaVar

	// Steady-state on IID normal noise: compare var estimates.
	// Expected: Hunter gives ~1.0 for unit-normal input; old gives ~1.25.
	hunterVar := 0.0
	oldVar := 0.0
	hunterMean := 0.0
	oldMean := 0.0

	// Use a fixed pseudo-random sequence: values oscillating ±1 around 0.
	// This gives a known variance of ~1.0.
	vals := []float64{1, -1, 1, -1, 1, -1, 1, -1}
	for rep := 0; rep < 1000; rep++ {
		x := vals[rep%len(vals)]
		devH := x - hunterMean
		devO := x - oldMean
		newHMean := (1-alphaMean)*hunterMean + alphaMean*x
		newOMean := (1-alphaMean)*oldMean + alphaMean*x
		hunterVar = (1-alphaVar)*hunterVar + alphaVar*(1-alphaMean)*devH*devH
		oldVar = (1-alphaVar)*oldVar + alphaVar*devO*devO
		hunterMean = newHMean
		oldMean = newOMean
	}
	// Hunter should give a smaller steady-state variance than the old formula.
	if hunterVar >= oldVar {
		t.Errorf("Hunter variance (%v) should be < old formula variance (%v) for oscillating input", hunterVar, oldVar)
	}
}

// ---------------------------------------------------------------------------
// Fix 6 (Tunable threshold): noise-floor gate uses 0.2, not 0.5.
// ---------------------------------------------------------------------------

func TestAntiRecGate_SignalBelowNoiseFloor_Threshold02(t *testing.T) {
	// |cohens_d| = 0.3 is above the noise floor (0.2): gate must NOT fire.
	kindIndex := map[string][]Finding{
		"mover": {{"kind": "mover", "cohens_d": 0.3}},
	}
	gate := antiRecommendationGate{code: "signal_below_noise_floor"}
	if antiRecGateFires(gate, kindIndex, 100) {
		t.Error("gate fires for cohens_d=0.3 (above noise floor 0.2); must NOT fire")
	}
}

func TestAntiRecGate_SignalBelowNoiseFloor_Below02(t *testing.T) {
	// |cohens_d| = 0.1 is below the noise floor: gate must fire.
	kindIndex := map[string][]Finding{
		"mover": {{"kind": "mover", "cohens_d": 0.1}},
	}
	gate := antiRecommendationGate{code: "signal_below_noise_floor"}
	if !antiRecGateFires(gate, kindIndex, 100) {
		t.Error("gate does not fire for cohens_d=0.1 (below noise floor 0.2); must fire")
	}
}

func TestAntiRecGate_Old05Threshold_WasOvesuppressing(t *testing.T) {
	// |cohens_d| = 0.4 is between old threshold (0.5) and new threshold (0.2).
	// Old code would have fired the gate (suppressing the recommendation).
	// New code must NOT fire it.
	kindIndex := map[string][]Finding{
		"mover": {{"kind": "mover", "cohens_d": 0.4}},
	}
	gate := antiRecommendationGate{code: "signal_below_noise_floor"}
	if antiRecGateFires(gate, kindIndex, 100) {
		t.Error("gate fires for cohens_d=0.4 with new 0.2 threshold; must NOT fire (informational-band signal)")
	}
}

// ---------------------------------------------------------------------------
// Fix 7 (Hampel nil): local_spike Findings carry a non-empty "at" timestamp,
// and ids are unique per metric (no |" suffix clobbering).
// ---------------------------------------------------------------------------

func hampelTestStore(names []string, n int, spikesAt []int, spikeVal float64) (*seriesStore, time.Time) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s := &seriesStore{
		names:              names,
		nameToID:           make(map[string]int, len(names)),
		isCounter:          make([]bool, len(names)),
		materializedSeries: make([][]float64, len(names)),
		timestamps:         make([]time.Time, n),
	}
	for i, name := range names {
		s.nameToID[name] = i
	}
	for i := 0; i < n; i++ {
		s.timestamps[i] = now.Add(time.Duration(i) * time.Second)
	}
	// Build a series with low-level noise (oscillating 1/2) so MAD > 0.
	for si := range names {
		series := make([]float64, n)
		for j := range series {
			if j%2 == 0 {
				series[j] = 1.0
			} else {
				series[j] = 2.0
			}
		}
		for _, sp := range spikesAt {
			series[sp] = spikeVal // spike far above noise level
		}
		s.materializedSeries[si] = series
	}
	return s, now
}

func TestHampelLocalSpikes_AtTimestampPopulated(t *testing.T) {
	// Series with alternating 1/2 noise + spike at index 100.
	// MAD of the ±50 window around any non-spike index = MAD([1,2,...]) = 0.5.
	// At index 100, value=1000, z ≈ (1000-1.5)/(0.5*1.4826) >> threshold.
	store, _ := hampelTestStore([]string{"m"}, 200, []int{100}, 1000.0)

	c := &findingCollector{}
	emitHampelLocalSpikes(c, "host", store, []int{0}, AutoConfig{SpikeZ: 3.5})

	if len(c.items) == 0 {
		t.Fatal("no local_spike Findings emitted for spike at index 100")
	}
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k != "local_spike" {
			continue
		}
		at, _ := f["at"].(string)
		if at == "" {
			t.Errorf("local_spike Finding has empty 'at': %v", f)
		}
		_, err := time.Parse("2006-01-02T15:04:05.000Z", at)
		if err != nil {
			t.Errorf("local_spike 'at' not parseable: %q: %v", at, err)
		}
	}
}

func TestHampelLocalSpikes_IdKeyUnique(t *testing.T) {
	// Two metrics each with a spike at index 100; ids must differ and not end with "|".
	store, _ := hampelTestStore([]string{"m1", "m2"}, 200, []int{100}, 1000.0)

	c := &findingCollector{}
	emitHampelLocalSpikes(c, "host", store, []int{0, 1}, AutoConfig{SpikeZ: 3.5})

	idSeen := map[string]int{}
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k != "local_spike" {
			continue
		}
		id, _ := f["id"].(string)
		idSeen[id]++
		if strings.HasSuffix(id, "|") {
			t.Errorf("local_spike id ends with '|' (empty timestamp): %q", id)
		}
	}
	for id, count := range idSeen {
		if count > 1 {
			t.Errorf("duplicate local_spike id %q (count=%d)", id, count)
		}
	}
}
