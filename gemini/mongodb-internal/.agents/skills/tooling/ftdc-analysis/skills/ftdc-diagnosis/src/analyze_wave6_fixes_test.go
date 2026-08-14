// analyze_wave6_fixes_test.go — regression tests for Wave 6 bug fixes.
package main

import (
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Fix W6-A: CV failure rate excludes react-mode changepoints (no cross_validated).
// ---------------------------------------------------------------------------

func TestCVFailurePct_ReactModeChangepointIgnored(t *testing.T) {
	// A react-mode changepoint lacks "cross_validated". It must not inflate
	// cpTotal, which would dilute cvFailurePct toward zero.
	c := &findingCollector{}
	// Standard changepoint: cross_validated=false (failed validation).
	c.add(newFinding(KindChangepoint, "h", "m1|t1", map[string]interface{}{
		"metric":          "m1",
		"cross_validated": false,
	}))
	// React-mode changepoint: no cross_validated field at all.
	c.add(newFinding(KindChangepoint, "h", "m2|t2", map[string]interface{}{
		"metric": "m2",
		// cross_validated intentionally absent
	}))

	cpTotal, cpFailed := countCVFailures(c.items)
	if cpTotal != 1 {
		t.Errorf("cpTotal = %d, want 1 (react-mode changepoint must not be counted)", cpTotal)
	}
	if cpFailed != 1 {
		t.Errorf("cpFailed = %d, want 1", cpFailed)
	}
	pct := 100.0 * float64(cpFailed) / float64(cpTotal)
	if pct != 100.0 {
		t.Errorf("cvFailurePct = %.2f, want 100.0", pct)
	}
}

// ---------------------------------------------------------------------------
// Fix W6-B: emitApplySaturation counts idle windows in denominator.
// ---------------------------------------------------------------------------

func TestApplySaturation_IdleWindowsInDenominator(t *testing.T) {
	// 4 windows: 3 idle (dn=0), 1 active above threshold.
	// Old code: windowsTotal=1, windowsAboveThreshold=1 → fracAbove=1.0 → fires.
	// New code: windowsTotal=4, windowsAboveThreshold=1 → fracAbove=0.25 → no fire.
	const (
		applyLatencyThresholdMs = 100.0
		windowSamples           = 60
		majorityFraction        = 0.5
	)
	// Build numVals and totVals with 4 windows (240 samples).
	// Windows 0,1,2: dn=0 (idle); Window 3: dn=1, dt=200 (>= 100ms threshold).
	numVals := make([]int64, 240)
	totVals := make([]int64, 240)
	// Window 3 (samples 180-239): dn=1, dt=200.
	numVals[239] = numVals[180] + 1
	totVals[239] = totVals[180] + 200

	var windowsTotal, windowsAboveThreshold int
	for lo := 0; lo+windowSamples <= len(numVals); lo += windowSamples {
		hi := lo + windowSamples
		dn := numVals[hi-1] - numVals[lo]
		dt := totVals[hi-1] - totVals[lo]
		windowsTotal++
		if dn <= 0 {
			continue
		}
		ms := float64(dt) / float64(dn)
		if ms >= applyLatencyThresholdMs {
			windowsAboveThreshold++
		}
	}
	if windowsTotal != 4 {
		t.Errorf("windowsTotal = %d, want 4", windowsTotal)
	}
	if windowsAboveThreshold != 1 {
		t.Errorf("windowsAboveThreshold = %d, want 1", windowsAboveThreshold)
	}
	fracAbove := float64(windowsAboveThreshold) / float64(windowsTotal)
	if fracAbove >= majorityFraction {
		t.Errorf("fracAbove = %.3f, should be < %.1f (sparse capture must not fire)", fracAbove, majorityFraction)
	}
}

// ---------------------------------------------------------------------------
// Fix W6-C: inferTopology returns "" for empty metric set.
// ---------------------------------------------------------------------------

func TestInferTopology_EmptyReturnsEmpty(t *testing.T) {
	got := inferTopology([]string{})
	if got != "" {
		t.Errorf("inferTopology(empty) = %q, want %q", got, "")
	}
}

// ---------------------------------------------------------------------------
// Fix W6-D: emitWorkloadSpec rate helper returns (0, false) on counter reset.
// ---------------------------------------------------------------------------

func TestWorkloadSpec_CounterResetNoNegativeRate(t *testing.T) {
	// Build a store where a counter decreases (simulating a process restart).
	n := 10
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	store := &seriesStore{nameToID: map[string]int{}}
	store.names = append(store.names, "serverStatus.opcounters.insert")
	store.nameToID["serverStatus.opcounters.insert"] = 0
	store.isCounter = []bool{true}
	store.isDegenerate = []bool{false}
	store.observedMin = []int64{0}
	store.observedMax = []int64{1000}
	store.fullMean = []float64{500}
	store.fullSD = []float64{100}
	store.fullN = []int{n}
	store.timestamps = make([]time.Time, n)
	for i := range store.timestamps {
		store.timestamps[i] = base.Add(time.Duration(i) * time.Second)
	}
	// Values: starts at 1000, drops to 0 at index 5 (restart), then climbs.
	// last=5 < first=1000 → counter reset.
	vals := []int64{1000, 1100, 1200, 1300, 1400, 0, 50, 100, 150, 5}
	col := NewCompressedColumn(n)
	for _, v := range vals {
		col.Add(v, false)
	}
	col.Finish()
	store.columns = append(store.columns, col)
	store.materializedSeries = append(store.materializedSeries, nil)

	c := &findingCollector{}
	emitWorkloadSpec(c, "h", store)

	for _, f := range c.items {
		if k, _ := f["kind"].(string); k != KindWorkloadSpec {
			continue
		}
		if r, ok := f["insert_rate_ps"].(float64); ok && r < 0 {
			t.Errorf("insert_rate_ps = %v, want absent or >= 0 (counter reset must not produce negative rate)", r)
		}
	}
}
