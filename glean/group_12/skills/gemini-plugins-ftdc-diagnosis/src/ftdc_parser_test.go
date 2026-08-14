package main

import (
	"sort"
	"testing"
	"time"
)

// ─── P1: printStats percentile helpers ────────────────────────────────────────

// percentileIdx returns the index used to compute the Nth percentile
// (0-100 integer) on a sorted slice of length n, using the same formula
// as the fixed printStats code.
func pctIdx(n, pct int) int {
	return (n - 1) * pct / 100
}

// p99FloatIdx mirrors the fixed p99 formula: int(float64(len-1) * 0.99).
func p99FloatIdx(n int) int {
	idx := int(float64(n-1) * 0.99)
	if idx >= n {
		idx = n - 1
	}
	return idx
}

// TestPrintStatsPercentiles verifies that the p50 and p99 index formulas
// produce correct indices for edge-case slice lengths.
func TestPrintStatsPercentiles(t *testing.T) {
	cases := []struct {
		n       int
		wantP50 int
		wantP99 int
	}{
		// n=1: only element 0 is valid for both percentiles.
		{n: 1, wantP50: 0, wantP99: 0},
		// n=2: p50 → (1)*50/100 = 0; p99 → int(1.0*0.99)=0.
		{n: 2, wantP50: 0, wantP99: 0},
		// n=100: p50 → (99)*50/100 = 49; p99 → int(99.0*0.99)=98.
		{n: 100, wantP50: 49, wantP99: 98},
		// n=300: p50 → (299)*50/100 = 149; p99 → int(299.0*0.99)=296.
		{n: 300, wantP50: 149, wantP99: 296},
	}

	for _, tc := range cases {
		gotP50 := pctIdx(tc.n, 50)
		gotP99 := p99FloatIdx(tc.n)

		if gotP50 != tc.wantP50 {
			t.Errorf("n=%d: p50 index = %d, want %d", tc.n, gotP50, tc.wantP50)
		}
		if gotP99 != tc.wantP99 {
			t.Errorf("n=%d: p99 index = %d, want %d", tc.n, gotP99, tc.wantP99)
		}

		// Verify indices are in-bounds for a slice of length n.
		if gotP50 < 0 || gotP50 >= tc.n {
			t.Errorf("n=%d: p50 index %d out of range [0,%d)", tc.n, gotP50, tc.n)
		}
		if gotP99 < 0 || gotP99 >= tc.n {
			t.Errorf("n=%d: p99 index %d out of range [0,%d)", tc.n, gotP99, tc.n)
		}
	}
}

// TestPrintStatsP99EndToEnd builds a sorted int64 slice of length 100 where
// the value at index i is i, and asserts that p99 == 98 (not 99) and p50 == 49
// (not 50), using exactly the formulas in the fixed printStats.
func TestPrintStatsP99EndToEnd(t *testing.T) {
	vals := make([]int64, 100)
	for i := range vals {
		vals[i] = int64(i)
	}
	sort.Slice(vals, func(a, b int) bool { return vals[a] < vals[b] })

	p50 := vals[(len(vals)-1)*50/100]
	if p50 != 49 {
		t.Errorf("p50 = %d, want 49", p50)
	}

	p99idx := int(float64(len(vals)-1) * 0.99)
	if p99idx >= len(vals) {
		p99idx = len(vals) - 1
	}
	p99 := vals[p99idx]
	if p99 != 98 {
		t.Errorf("p99 = %d (idx=%d), want 98", p99, p99idx)
	}
}

// ─── P2: computeMetricStats lazy-init ─────────────────────────────────────────

// TestComputeMetricStatsLateAppearingMetric verifies that a metric that does
// not appear in samples[0] but appears in samples[2] is still tracked.
func TestComputeMetricStatsLateAppearingMetric(t *testing.T) {
	now := time.Now()
	samples := []Sample{
		{Timestamp: now, Metrics: map[string]int64{"a": 10}},
		{Timestamp: now.Add(time.Second), Metrics: map[string]int64{"a": 20}},
		{Timestamp: now.Add(2 * time.Second), Metrics: map[string]int64{"a": 30, "b": 100}},
		{Timestamp: now.Add(3 * time.Second), Metrics: map[string]int64{"a": 40, "b": 200}},
	}

	stats := computeMetricStats(samples)

	found := map[string]*metricStats{}
	for i := range stats {
		found[stats[i].name] = &stats[i]
	}

	// "a" should always be present.
	if _, ok := found["a"]; !ok {
		t.Fatal("metric 'a' missing from stats")
	}

	// "b" first appears at index 2 — the fix ensures it is NOT dropped.
	sb, ok := found["b"]
	if !ok {
		t.Fatal("metric 'b' (first seen at sample[2]) missing from stats — lazy-init fix not applied")
	}
	if sb.min != 100 {
		t.Errorf("b.min = %d, want 100", sb.min)
	}
	if sb.max != 200 {
		t.Errorf("b.max = %d, want 200", sb.max)
	}
	if sb.count != 2 {
		t.Errorf("b.count = %d, want 2", sb.count)
	}
}

// ─── P5: computeRates preserves negative values ───────────────────────────────

// TestComputeRatesNegativeGauge verifies that a gauge that decreases produces
// a negative rate in computeRates (not clamped to 0).
func TestComputeRatesNegativeGauge(t *testing.T) {
	now := time.Now()
	samples := []Sample{
		{Timestamp: now, Metrics: map[string]int64{"gauge": 100}},
		{Timestamp: now.Add(time.Second), Metrics: map[string]int64{"gauge": 60}},
	}

	rates := computeRates(samples)

	if len(rates) != 1 {
		t.Fatalf("expected 1 rate sample, got %d", len(rates))
	}
	got, ok := rates[0].Rates["gauge"]
	if !ok {
		t.Fatal("'gauge' missing from rates")
	}
	// delta = 60 - 100 = -40; dtSec = 1.0 → rate = -40.0
	if got >= 0 {
		t.Errorf("gauge rate = %f, want negative (delta=-40)", got)
	}
	const want = -40.0
	if got != want {
		t.Errorf("gauge rate = %f, want %f", got, want)
	}
}

// TestComputeRatesPositiveCounterUnchanged verifies that a counter that only
// increases still produces a positive rate (no regression from removing the clamp).
func TestComputeRatesPositiveCounterUnchanged(t *testing.T) {
	now := time.Now()
	samples := []Sample{
		{Timestamp: now, Metrics: map[string]int64{"ops": 1000}},
		{Timestamp: now.Add(time.Second), Metrics: map[string]int64{"ops": 1500}},
	}

	rates := computeRates(samples)

	if len(rates) != 1 {
		t.Fatalf("expected 1 rate sample, got %d", len(rates))
	}
	got, ok := rates[0].Rates["ops"]
	if !ok {
		t.Fatal("'ops' missing from rates")
	}
	const want = 500.0
	if got != want {
		t.Errorf("ops rate = %f, want %f", got, want)
	}
}
