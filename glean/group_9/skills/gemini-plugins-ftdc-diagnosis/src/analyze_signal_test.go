package main

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

// TestDFTPower_Sinusoid: a pure sin(2π·t/P) signal should have ≥10×
// power at period P compared to a non-matching candidate period.
func TestDFTPower_Sinusoid(t *testing.T) {
	cases := []struct{ period int }{{8}, {16}, {32}, {60}}
	for _, c := range cases {
		const N = 600
		x := make([]float64, N)
		for i := range x {
			x[i] = math.Sin(2 * math.Pi * float64(i) / float64(c.period))
		}
		pMatch := dftPower(x, c.period)
		pOther := dftPower(x, c.period*3+5) // non-harmonic mismatch
		if pMatch < 10*pOther {
			t.Errorf("period=%d: P(matching)=%v, P(other)=%v — expected ≥10× ratio",
				c.period, pMatch, pOther)
		}
	}
}

// TestDetrendFloat_PureRamp: y = a + b·i ⇒ residual is uniformly ~0.
func TestDetrendFloat_PureRamp(t *testing.T) {
	const N = 500
	x := make([]float64, N)
	for i := range x {
		x[i] = 3.0 + 2.5*float64(i)
	}
	r := detrendFloat(x)
	if len(r) != N {
		t.Fatalf("length: got %d, want %d", len(r), N)
	}
	for i, v := range r {
		if math.Abs(v) > 1e-9 {
			t.Errorf("detrendFloat residual[%d] = %v, want ~0", i, v)
		}
	}
}

// TestDetrendFloat_NoisyRamp: residual stddev ≈ noise stddev within 5%;
// residual mean ≈ 0.
func TestDetrendFloat_NoisyRamp(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	const N = 5000
	const sigma = 2.0
	x := make([]float64, N)
	for i := range x {
		x[i] = 10.0 + 0.5*float64(i) + rng.NormFloat64()*sigma
	}
	r := detrendFloat(x)
	var mean, ssq float64
	for _, v := range r {
		mean += v
	}
	mean /= float64(N)
	for _, v := range r {
		d := v - mean
		ssq += d * d
	}
	sd := math.Sqrt(ssq / float64(N))
	if math.Abs(mean) > 0.1 {
		t.Errorf("residual mean = %v, want |mean| < 0.1", mean)
	}
	if sd < 0.9*sigma || sd > 1.1*sigma {
		t.Errorf("residual sd = %v, want %v (±5%%)", sd, sigma)
	}
}

// TestEmitPeriodic_FiresOnSynthetic: inject a clean sinusoid into a
// metric column, assert one periodic Finding emitted with the right
// period. Uses N=3000 (50 cycles of period=60) to ensure SNR clearly
// exceeds the Bonferroni-Schuster threshold; shorter N can produce
// SNR right at the boundary and flake.
func TestEmitPeriodic_FiresOnSynthetic(t *testing.T) {
	const N = 3000
	const period = 60
	series := make([]float64, N)
	for i := range series {
		series[i] = 100 + 20*math.Sin(2*math.Pi*float64(i)/float64(period))
	}
	store := &seriesStore{nameToID: map[string]int{}}
	store.names = []string{"test.metric"}
	store.nameToID["test.metric"] = 0
	store.isCounter = []bool{false}
	store.isDegenerate = []bool{false}
	store.observedMin = []int64{80}
	store.observedMax = []int64{120}
	store.fullMean = []float64{100}
	store.fullSD = []float64{14}
	store.materializedSeries = [][]float64{series}
	store.columns = []*CompressedColumn{NewCompressedColumn(N)}
	for _, v := range series {
		store.columns[0].Add(int64(v), true)
	}
	store.columns[0].Finish()
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < N; i++ {
		store.timestamps = append(store.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	c := &findingCollector{}
	emitted, _ := emitPeriodic(c, "h", store, []int{0}, defaultAutoConfig())
	if emitted < 1 {
		t.Fatalf("expected ≥1 periodic Finding on clean sinusoid, got %d", emitted)
	}
	var f Finding
	for _, x := range c.items {
		if k, _ := x["kind"].(string); k == "periodic" {
			f = x
			break
		}
	}
	if f == nil {
		t.Fatal("no periodic Finding in collector")
	}
	if p, _ := f["period_seconds"].(int); p != period {
		t.Errorf("period_seconds: got %d, want %d", p, period)
	}
}

// TestEmitPeriodic_DoesNotFireOnNoise: pure Gaussian noise should emit
// zero periodic Findings under the α=0.001 / Bonferroni threshold.
// (Slightly larger N matches the synthetic sinusoid test for parity.)
func TestEmitPeriodic_DoesNotFireOnNoise(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	const N = 3000
	series := make([]float64, N)
	for i := range series {
		series[i] = rng.NormFloat64()
	}
	store := &seriesStore{nameToID: map[string]int{"test.metric": 0}}
	store.names = []string{"test.metric"}
	store.isCounter = []bool{false}
	store.isDegenerate = []bool{false}
	store.observedMin = []int64{-3}
	store.observedMax = []int64{3}
	store.fullMean = []float64{0}
	store.fullSD = []float64{1}
	store.materializedSeries = [][]float64{series}
	store.columns = []*CompressedColumn{NewCompressedColumn(N)}
	for _, v := range series {
		store.columns[0].Add(int64(v*1000), true) // scale for int64 storage
	}
	store.columns[0].Finish()
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < N; i++ {
		store.timestamps = append(store.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	c := &findingCollector{}
	emitted, _ := emitPeriodic(c, "h", store, []int{0}, defaultAutoConfig())
	if emitted > 0 {
		t.Errorf("expected 0 periodic Findings on pure noise, got %d", emitted)
	}
}

// TestEmitSuspiciousSmoothness_VarianceTooLow: a near-constant series
// (mean=100, sd=1e-12) fires variance_too_low.
func TestEmitSuspiciousSmoothness_VarianceTooLow(t *testing.T) {
	const N = 200
	series := make([]float64, N)
	for i := range series {
		series[i] = 100 + 1e-13*float64(i%2)
	}
	store := buildStoreWithMaterialized("test.flat", series, false)
	c := &findingCollector{}
	emitSuspiciousSmoothness(c, "h", store, []int{0})
	found := false
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "suspicious_smoothness" {
			if d, _ := f["detector"].(string); d == "variance_too_low" {
				found = true
			}
		}
	}
	if !found {
		t.Errorf("expected variance_too_low detector to fire on near-constant series")
	}
}

// TestEmitSuspiciousSmoothness_PerfectlySteadyRate: a counter with
// exactly-constant rate fires perfectly_steady_rate.
//
// materializedSeries for a counter is already first-differenced (rate values):
// series[0]=0 (no prior sample), series[1..N-1] = R (constant rate).
// The detector must compare rate values directly, not second differences.
func TestEmitSuspiciousSmoothness_PerfectlySteadyRate(t *testing.T) {
	const N = 100
	const R = 7.0
	series := make([]float64, N)
	// series[0] = 0 (Go zero init, no prior sample); rates start at index 1.
	for i := 1; i < N; i++ {
		series[i] = R
	}
	store := buildStoreWithMaterialized("test.counter", series, true)
	c := &findingCollector{}
	emitSuspiciousSmoothness(c, "h", store, []int{0})
	found := false
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "suspicious_smoothness" {
			if d, _ := f["detector"].(string); d == "perfectly_steady_rate" {
				found = true
			}
		}
	}
	if !found {
		t.Errorf("expected perfectly_steady_rate detector to fire on uniform-rate counter (rate format: [0, R, R, ...])")
	}
}

// TestEmitSuspiciousSmoothness_RoundNumberDensity: >90% of values are
// multiples of 10.
func TestEmitSuspiciousSmoothness_RoundNumberDensity(t *testing.T) {
	const N = 100
	series := make([]float64, N)
	for i := range series {
		series[i] = float64(((i * 31) % 50) * 10) // all multiples of 10
	}
	store := buildStoreWithMaterialized("test.round", series, false)
	c := &findingCollector{}
	emitSuspiciousSmoothness(c, "h", store, []int{0})
	found := false
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "suspicious_smoothness" {
			if d, _ := f["detector"].(string); d == "round_number_density" {
				found = true
			}
		}
	}
	if !found {
		t.Errorf("expected round_number_density detector to fire on multiples-of-10 series")
	}
}

// TestEmitEvidenceConvergence_DistinctKinds: synthetic_event with
// members across 3 distinct kinds produces distinct_kinds_count = 3.
func TestEmitEvidenceConvergence_DistinctKinds(t *testing.T) {
	c := &findingCollector{}
	// Three child Findings with distinct kinds.
	c.add(Finding{"kind": KindSpike, "id": "f-aaaa", "host": "h"})
	c.add(Finding{"kind": KindChangepoint, "id": "f-bbbb", "host": "h"})
	c.add(Finding{"kind": KindVarianceShift, "id": "f-cccc", "host": "h"})
	// And the synthetic_event referencing them.
	c.add(Finding{
		"kind":               "synthetic_event",
		"id":                 "f-dddd",
		"host":               "h",
		"member_finding_ids": []string{"f-aaaa", "f-bbbb", "f-cccc"},
		"at_start":           "2026-01-01T00:00:00.000Z",
	})
	emitted := emitEvidenceConvergence(c, "h")
	if emitted != 1 {
		t.Fatalf("expected 1 evidence_convergence Finding, got %d", emitted)
	}
	var ec Finding
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "evidence_convergence" {
			ec = f
		}
	}
	if ec == nil {
		t.Fatal("no evidence_convergence Finding emitted")
	}
	if d, _ := ec["distinct_kinds_count"].(int); d != 3 {
		t.Errorf("distinct_kinds_count: got %d, want 3", d)
	}
}

// TestAutoStats_Records confirms timeStage records emitted + duration.
func TestAutoStats_Records(t *testing.T) {
	c := &findingCollector{}
	s := &AutoStats{}
	s.timeStage(c, "test", 0, func() {
		c.add(Finding{"kind": "foo"})
		c.add(Finding{"kind": "bar"})
	})
	if len(s.Stages) != 1 {
		t.Fatalf("expected 1 stage recorded, got %d", len(s.Stages))
	}
	st := s.Stages[0]
	if st.Name != "test" {
		t.Errorf("name: got %q, want test", st.Name)
	}
	if st.EmittedN != 2 {
		t.Errorf("emitted: got %d, want 2", st.EmittedN)
	}
	if st.DurationMs < 0 {
		t.Errorf("duration negative: %v", st.DurationMs)
	}
}

// TestAutoStats_AsJSON returns the per-stage struct as a slice of maps.
func TestAutoStats_AsJSON(t *testing.T) {
	s := &AutoStats{Stages: []AutoStageStat{
		{Name: "a", EmittedN: 5, SkippedN: 1, DurationMs: 2.5},
		{Name: "b", EmittedN: 0, SkippedN: 0, DurationMs: 0.1},
	}}
	out := s.asJSON()
	if len(out) != 2 {
		t.Fatalf("len: got %d, want 2", len(out))
	}
	a := out[0].(map[string]interface{})
	if a["name"] != "a" || a["emitted"] != 5 {
		t.Errorf("stage 0 unexpected: %v", a)
	}
}

// buildStoreWithMaterialized constructs a 1-metric seriesStore with a
// pre-populated materialized series + column (int64-quantized).
func buildStoreWithMaterialized(name string, series []float64, isCounter bool) *seriesStore {
	store := &seriesStore{nameToID: map[string]int{name: 0}}
	store.names = []string{name}
	store.isCounter = []bool{isCounter}
	store.isDegenerate = []bool{false}
	store.observedMin = []int64{0}
	store.observedMax = []int64{0}
	store.fullMean = []float64{0}
	store.fullSD = []float64{1}
	store.materializedSeries = [][]float64{series}
	col := NewCompressedColumn(len(series))
	for _, v := range series {
		col.Add(int64(math.Round(v)), true)
	}
	col.Finish()
	store.columns = []*CompressedColumn{col}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < len(series); i++ {
		store.timestamps = append(store.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	return store
}
