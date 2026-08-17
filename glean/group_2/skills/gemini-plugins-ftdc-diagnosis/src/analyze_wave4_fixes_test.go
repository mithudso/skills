// analyze_wave4_fixes_test.go — regression tests for Wave 4 bug fixes.
package main

import (
	"math"
	"testing"
	"time"
)

// addColToStore appends one metric (gauge by default; set isCounter=true for
// counters) to store. The caller must set store.timestamps separately.
func addColToStore(store *seriesStore, name string, vals []int64, isCounter bool) {
	id := len(store.names)
	store.names = append(store.names, name)
	store.nameToID[name] = id
	store.isCounter = append(store.isCounter, isCounter)
	store.isDegenerate = append(store.isDegenerate, false)
	store.observedMin = append(store.observedMin, 0)
	store.observedMax = append(store.observedMax, 0)
	store.fullMean = append(store.fullMean, 0)
	store.fullSD = append(store.fullSD, 0)
	store.fullN = append(store.fullN, len(vals))
	col := NewCompressedColumn(len(vals))
	for _, v := range vals {
		col.Add(v, false)
	}
	col.Finish()
	store.columns = append(store.columns, col)
}

// ---------------------------------------------------------------------------
// Fix W4-A: emitConsistencyWarning emits in deterministic order.
// ---------------------------------------------------------------------------

func TestConsistencyWarning_DeterministicOrder(t *testing.T) {
	// Two metrics both have mover/trend mismatches. The order of the emitted
	// consistency_warning Findings must be alphabetical (deterministic).
	c := &findingCollector{}
	c.add(Finding{"kind": KindMover, "metric": "z.metric", "direction": "up"})
	c.add(Finding{"kind": KindMover, "metric": "a.metric", "direction": "down"})
	c.add(Finding{"kind": KindMonotonicTrend, "metric": "z.metric", "slope": float64(-1.0)})
	c.add(Finding{"kind": KindMonotonicTrend, "metric": "a.metric", "slope": float64(1.0)})

	emitConsistencyWarning(c, "h")

	var warnings []string
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "consistency_warning" {
			warnings = append(warnings, f["metric"].(string))
		}
	}
	if len(warnings) != 2 {
		t.Fatalf("expected 2 consistency_warning Findings, got %d", len(warnings))
	}
	if warnings[0] != "a.metric" || warnings[1] != "z.metric" {
		t.Errorf("consistency_warning order wrong: got %v, want [a.metric, z.metric]", warnings)
	}
}

// ---------------------------------------------------------------------------
// Fix W4-D/W4-E: emitUtilization uses correct baseline/stress window means.
// ---------------------------------------------------------------------------

func TestUtilization_BaselineStressMeansUseCorrectSamples(t *testing.T) {
	// Build a store with the wiredTiger cache utilization rule metrics.
	// Samples 0..99 = low ratio (baseline); samples 100..199 = high ratio (stress).
	// bStart=0, bEnd=50, sStart=150, sEnd=200.
	// Expected baseline_mean ≈ 0.1, stress_mean ≈ 0.9.
	n := 200
	store := &seriesStore{nameToID: map[string]int{}}

	// cache.bytes currently in cache: low in baseline (10), high in stress (90).
	// cache.maximum bytes configured: always 100.
	metricVals := make([]int64, n)
	maxVals := make([]int64, n)
	for i := 0; i < n; i++ {
		maxVals[i] = 100
		if i < 100 {
			metricVals[i] = 10 // baseline: 10% fill
		} else {
			metricVals[i] = 90 // stress: 90% fill
		}
	}
	addColToStore(store, "serverStatus.wiredTiger.cache.bytes currently in the cache", metricVals, false)
	addColToStore(store, "serverStatus.wiredTiger.cache.maximum bytes configured", maxVals, false)

	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	store.timestamps = make([]time.Time, n)
	for i := range store.timestamps {
		store.timestamps[i] = base.Add(time.Duration(i) * time.Second)
	}

	c := &findingCollector{}
	// bStart=0, bEnd=50 (baseline in low-ratio zone)
	// sStart=150, sEnd=200 (stress in high-ratio zone)
	emitUtilization(c, "h", store, 0, 50, 150, 200)

	var baseMean, stressMean float64
	found := false
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "utilization" {
			if m, _ := f["metric"].(string); m == "serverStatus.wiredTiger.cache.bytes currently in the cache" {
				baseMean, _ = f["baseline_mean"].(float64)
				stressMean, _ = f["stress_mean"].(float64)
				found = true
			}
		}
	}
	if !found {
		t.Fatal("no utilization Finding for wiredTiger.cache metric")
	}
	// Baseline window samples 0..49: all 10/100 = 0.1.
	if math.Abs(baseMean-0.1) > 0.01 {
		t.Errorf("baseline_mean = %v, want ~0.1 (before fix: sort scrambled the window)", baseMean)
	}
	// Stress window samples 150..199: all 90/100 = 0.9.
	if math.Abs(stressMean-0.9) > 0.01 {
		t.Errorf("stress_mean = %v, want ~0.9 (before fix: sort scrambled the window)", stressMean)
	}
}

// ---------------------------------------------------------------------------
// Fix W4-G: periodic detector skips emission when len(powers) < 2.
// ---------------------------------------------------------------------------

func TestPeriodicSignal_SingleCandidateSkippedNotNaN(t *testing.T) {
	// At 1 Hz with N=60 samples the period filter keeps only the 30 s
	// candidate (30*2=60 ≤ 60); the 45 s candidate requires N≥90 and is
	// dropped.  That gives len(powers)==1, which before the fix produced
	// NaN SNR via 0/0.  After the fix the metric must be skipped
	// (skipped==1, emitted==0) and no NaN Finding may be emitted.
	//
	// We build a column-backed store so nSamples() returns 60 and
	// materializeFromColumns can populate materializedSeries.
	n := 60
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")

	vals := make([]int64, n)
	for i := range vals {
		vals[i] = 100
		if i%30 == 0 {
			vals[i] = 200
		}
	}
	store := &seriesStore{nameToID: map[string]int{}}
	addColToStore(store, "test.periodic30", vals, false)

	store.timestamps = make([]time.Time, n)
	for i := range store.timestamps {
		store.timestamps[i] = base.Add(time.Duration(i) * time.Second)
	}
	store.materializedSeries = make([][]float64, 1)
	store.materializeFromColumns([]int{0})

	cfg := defaultAutoConfig()
	c := &findingCollector{}
	_, skipped := emitPeriodic(c, "h", store, []int{0}, cfg)

	// The single-candidate path must skip rather than emit a NaN SNR.
	if skipped == 0 {
		t.Error("expected skipped > 0 for single-candidate series (len(powers)==1), got 0")
	}
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k != "periodic" {
			continue
		}
		snr, ok := f["snr"].(float64)
		if !ok {
			continue
		}
		if math.IsNaN(snr) {
			t.Errorf("periodic Finding emitted with NaN snr: %v", f)
		}
		if math.IsInf(snr, 0) {
			t.Errorf("periodic Finding emitted with Inf snr: %v", f)
		}
	}
}
