// analyze_wave3_fixes_test.go — regression tests for Wave 3 bug fixes.
package main

import (
	"math"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Fix W3-C: emitPrecursorCandidates must filter by host in multi-host mode.
// ---------------------------------------------------------------------------

func TestPrecursorCandidates_HostFilter(t *testing.T) {
	// Two changepoints: one on hostA (the target), one on hostB (different host,
	// 30 seconds earlier). Without the host filter, hostB's changepoint would
	// appear as a precursor for hostA's changepoint.
	base, _ := time.Parse("2006-01-02T15:04:05.000Z", "2026-01-01T00:01:00.000Z")
	earlier := base.Add(-30 * time.Second)

	c := &findingCollector{}
	// changepoint on hostA at T+60s (the target)
	c.add(Finding{
		"kind":       KindChangepoint,
		"host":       "hostA",
		"metric":     "m1",
		"at":         base.UTC().Format("2006-01-02T15:04:05.000Z"),
		"effect_size": 2.0,
		"id":          "f-host-a",
	})
	// changepoint on hostB at T+30s (earlier — would be a precursor without host filter)
	c.add(Finding{
		"kind":       KindChangepoint,
		"host":       "hostB",
		"metric":     "m2",
		"at":         earlier.UTC().Format("2006-01-02T15:04:05.000Z"),
		"effect_size": 2.0,
		"id":          "f-host-b",
	})

	emitPrecursorCandidates(c, "hostA")

	for _, f := range c.items {
		if k, _ := f["kind"].(string); k != "precursor_candidates" {
			continue
		}
		precursors, _ := f["precursors"].([]interface{})
		for _, p := range precursors {
			pm, _ := p.(map[string]interface{})
			metric, _ := pm["metric"].(string)
			if metric == "m2" {
				t.Errorf("cross-host precursor leaked into hostA precursor_candidates: %v", pm)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Fix W3-F: recurrence fingerprint is deterministic on tied _rate_ps values.
// ---------------------------------------------------------------------------

func TestRecurrenceFingerprintDeterministicOnTiedRates(t *testing.T) {
	// Construct a collector with a workload_spec finding whose two _rate_ps
	// fields tie at the same value. The fix uses lexicographic tie-breaking,
	// so "insert_rate_ps" < "query_rate_ps" — insert must always win.
	c := &findingCollector{}
	c.add(Finding{
		"kind":           "workload_spec",
		"host":           "h",
		"insert_rate_ps": float64(1500.0),
		"query_rate_ps":  float64(1500.0),
	})

	// Run multiple times: map iteration is randomized, so without tie-breaking
	// the dominant_op would differ. With the fix it must be stable.
	var firstFP string
	for i := 0; i < 200; i++ {
		// Re-create the collector each iteration so the map is fresh.
		ci := &findingCollector{}
		ci.add(Finding{
			"kind":           "workload_spec",
			"host":           "h",
			"insert_rate_ps": float64(1500.0),
			"query_rate_ps":  float64(1500.0),
		})
		emitRecurrenceFingerprint(ci, "h")
		for _, f := range ci.items {
			if k, _ := f["kind"].(string); k == "recurrence_fingerprint" {
				fp, _ := f["fingerprint"].(string)
				if firstFP == "" {
					firstFP = fp
				} else if fp != firstFP {
					t.Errorf("non-deterministic recurrence_fingerprint on tied rates: %q vs %q (iteration %d)", firstFP, fp, i)
					return
				}
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Fix W3-G: bhSignificanceThreshold returns -1.0 (not 0.0) when no test passes.
// ---------------------------------------------------------------------------

func TestBHSignificanceThreshold_NoPassReturnsNegativeOne(t *testing.T) {
	// All p-values above the BH threshold for alpha=0.05: no test passes.
	// Old code returned 0.0 (ambiguous: could also mean "threshold = 0").
	// New code returns -1.0.
	got := bhSignificanceThreshold([]float64{0.5, 0.9, 0.99}, 0.05)
	if got != -1.0 {
		t.Errorf("bhSignificanceThreshold: got %v for no-passing tests, want -1.0", got)
	}
}

func TestBHSignificanceThreshold_ZeroPValuePassesAndIsDistinguishable(t *testing.T) {
	// A p-value of exactly 0.0 should pass BH (0.0 <= bhCut always) and be
	// returned as 0.0 — distinguishable from the "no test passed" -1.0 sentinel.
	got := bhSignificanceThreshold([]float64{0.0, 0.5, 0.9}, 0.05)
	if got != 0.0 {
		t.Errorf("bhSignificanceThreshold: got %v for p=0 input, want 0.0 (zero p-value should pass)", got)
	}
}

// ---------------------------------------------------------------------------
// Fix W3-H: seriesStore.fullN populated; pooled SD uses actual sample count.
// ---------------------------------------------------------------------------

func TestSeriesStore_FullNPopulated(t *testing.T) {
	// Build a 100-sample store for one gauge and one counter. fullN should
	// reflect the number of Welford samples seen (gauge: 100, counter: 99).
	n := 100
	samples := make([]Sample, n)
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := range samples {
		samples[i].Timestamp = base.Add(time.Duration(i) * time.Second)
		samples[i].Metrics = map[string]int64{
			"gauge.metric":   int64(i + 1),
			"counter.metric": int64(i * 10), // monotonically increasing → counter
		}
	}
	store := buildSeriesStore(samples)

	gid, ok := store.nameToID["gauge.metric"]
	if !ok {
		t.Fatal("gauge.metric not in store")
	}
	cid, ok := store.nameToID["counter.metric"]
	if !ok {
		t.Fatal("counter.metric not in store")
	}

	if store.fullN[gid] < n-1 {
		t.Errorf("gauge fullN = %d, want >= %d", store.fullN[gid], n-1)
	}
	// Counter Welford is on rates (differences between consecutive samples).
	// With 100 samples there are 99 rate deltas.
	if store.fullN[cid] < n-2 {
		t.Errorf("counter fullN = %d, want >= %d", store.fullN[cid], n-2)
	}
}

func TestPooledSD_LargeNDiffersFromN2(t *testing.T) {
	// pooledStandardDeviation(sd, 2, sd, 2) gives sqrt((sd^2+sd^2)/2) = sd.
	// pooledStandardDeviation(sd1, 1000, sd2, 10) should weight sd1 much more.
	// Concretely: when n1=1000, n2=10, sd1=1.0, sd2=100.0:
	//   n=2 formula: sqrt((1^2+100^2)/2) = sqrt(5000.5) ≈ 70.7
	//   correct:     sqrt((999*1^2 + 9*100^2)/(1000+10-2)) = sqrt((999+90000)/1008) = sqrt(90.7) ≈ 9.52
	sdA, sdB := 1.0, 100.0
	nA, nB := 1000, 10

	withN2 := pooledStandardDeviation(sdA, 2, sdB, 2)
	withActualN := pooledStandardDeviation(sdA, nA, sdB, nB)

	// The two should differ significantly when N is unequal.
	if math.Abs(withN2-withActualN)/withN2 < 0.5 {
		t.Errorf("pooledSD with n=2 (%v) and actual n (%v) are too close; expected large difference for 100:1 n ratio", withN2, withActualN)
	}
	// Actual N result should be closer to sdA (the high-N group dominates).
	if withActualN > withN2 {
		t.Errorf("actualN pooledSD (%v) should be < n=2 pooledSD (%v) since high-N low-SD group dominates", withActualN, withN2)
	}
}
