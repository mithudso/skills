// analyze_wave10_fixes_test.go — regression tests for Wave 10 confirmed bugs.
//
// W10-1: analyze_counter.go filter inverted (isCounter vs !isCounter).
// W10-2: analyze_hypothesis.go discriminating_kinds_needed asserted as []interface{} not []string.
// W10-2b: analyze_uncertainty.go same []interface{} bug.
// W10-3: analyze_joint_changepoint.go direction always "up" (cohens_d is absFloat, never negative).
// W10-4: analyze_tunable.go skipped counter hardcoded to 0.
// W10-6: pattern_dsl.go loadUserPatterns accepts empty required_kinds.
package main

import (
	"strconv"
	"testing"
	"time"
)

// W10-1: emitCounterEvents must process gauge-classified metrics (isCounter=false),
// not counter-classified ones. A metric that wrapped or reset has isCounter=false
// because the decrease was observed during store build.
func TestCounterFilterInversionFix(t *testing.T) {
	// A metric with isCounter=false (has a decrease) must emit a finding.
	col := &CompressedColumn{}
	for _, v := range []int64{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160, 170, 180, 175} {
		col.Add(v, true)
	}
	col.Finish()
	ts := make([]time.Time, 20)
	for i := range ts {
		ts[i] = time.Now().Add(time.Duration(i) * time.Second)
	}
	store := &seriesStore{
		nameToID:   map[string]int{"test.resetting": 0},
		names:      []string{"test.resetting"},
		isCounter:  []bool{false}, // classified as gauge due to the observed decrease
		columns:    []*CompressedColumn{col},
		timestamps: ts,
	}
	c := &findingCollector{}
	emitted, _ := emitCounterEvents(c, "h", store)
	if emitted == 0 {
		t.Error("emitCounterEvents emitted 0 findings for gauge-classified metric with decrease; pre-fix filter would also skip it")
	}

	// A metric with isCounter=true (no decreases) must NOT emit.
	colMonotone := &CompressedColumn{}
	for _, v := range []int64{0, 1, 2, 3, 4, 5} {
		colMonotone.Add(v, true)
	}
	colMonotone.Finish()
	ts2 := make([]time.Time, 6)
	for i := range ts2 {
		ts2[i] = time.Now().Add(time.Duration(i) * time.Second)
	}
	store2 := &seriesStore{
		nameToID:   map[string]int{"test.monotone": 0},
		names:      []string{"test.monotone"},
		isCounter:  []bool{true}, // never decreased: true counter, nothing to detect
		columns:    []*CompressedColumn{colMonotone},
		timestamps: ts2,
	}
	c2 := &findingCollector{}
	emitted2, _ := emitCounterEvents(c2, "h", store2)
	if emitted2 != 0 {
		t.Errorf("emitCounterEvents emitted %d findings for truly monotone counter; want 0", emitted2)
	}
}

// W10-2: emitHypothesisTests must emit hypothesis_test Findings when
// discriminating_kinds_needed is non-empty and confidence < 0.9.
// Pre-fix: []interface{} assertion failed silently → always 0 emitted.
func TestHypothesisTestsEmitted(t *testing.T) {
	c := &findingCollector{}
	// Add a diagnosis_candidate with confidence < 0.9 and a discriminating kind.
	c.add(newFinding("diagnosis_candidate", "h1", "f-abc", map[string]interface{}{
		"pattern_name":               "primary_pali_log_pressure",
		"confidence_rank":            0.72,
		"discriminating_kinds_needed": []string{"counter_wrap", "regression"},
	}))

	emitted, _ := emitHypothesisTests(c, "h1")
	if emitted == 0 {
		t.Error("emitHypothesisTests emitted 0 hypothesis_test Findings; pre-fix []interface{} assertion always failed")
	}
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		if k != "hypothesis_test" {
			continue
		}
		if _, ok := f["recipe"].(string); !ok {
			t.Error("hypothesis_test Finding missing recipe field")
		}
	}
}

// W10-2b: emitDiagnosisUncertainty must emit diagnosis_uncertainty Findings.
// Pre-fix: []interface{} assertion failed → always 0 emitted.
// Also verifies skipped count: a high-confidence candidate (≥0.9) must be counted as skipped.
func TestDiagnosisUncertaintyEmitted(t *testing.T) {
	c := &findingCollector{}
	c.add(newFinding("diagnosis_candidate", "h1", "f-def", map[string]interface{}{
		"pattern_name":               "primary_pali_log_pressure",
		"confidence_rank":            0.65,
		"discriminating_kinds_needed": []string{"counter_wrap"},
	}))
	// High-confidence candidate: must be skipped, not emitted.
	c.add(newFinding("diagnosis_candidate", "h1", "f-high", map[string]interface{}{
		"pattern_name":               "primary_pali_log_pressure",
		"confidence_rank":            0.95,
		"discriminating_kinds_needed": []string{"counter_wrap"},
	}))

	emitted, skipped := emitDiagnosisUncertainty(c, "h1")
	if emitted == 0 {
		t.Error("emitDiagnosisUncertainty emitted 0 diagnosis_uncertainty Findings; pre-fix []interface{} assertion always failed")
	}
	if skipped == 0 {
		t.Error("emitDiagnosisUncertainty skipped=0; high-confidence candidate must be counted as skipped")
	}
}

// W10-3: joint_changepoint direction must use stressed_mean vs baseline_mean.
// Pre-fix: cohens_d is always ≥ 0 (absFloat), so dir was always "up", n_down always 0.
func TestJointChangepointDirection(t *testing.T) {
	c := &findingCollector{}
	// Add 5 changepoint Findings all at sample_index=10, all with stressed < baseline
	// (metric went DOWN). cohens_d is positive (absFloat), which was the pre-fix trap.
	for i := 0; i < 5; i++ {
		metricName := "serverStatus.metric" + strconv.Itoa(i)
		c.add(newFinding(KindChangepoint, "h1", metricName+"|10", map[string]interface{}{
			"metric":         metricName,
			"sample_index":   10,
			"cohens_d":       1.5,    // always positive — the pre-fix direction check used this
			"baseline_mean":  200.0,  // before: high
			"stressed_mean":  100.0,  // after: low → direction should be "down"
		}))
	}

	emitted, _ := emitJointChangepoint(c, "h1")
	if emitted == 0 {
		t.Fatal("emitJointChangepoint emitted 0 findings; expected joint_changepoint for cluster of 5")
	}

	for _, f := range c.items {
		k, _ := f["kind"].(string)
		if k != "joint_changepoint" {
			continue
		}
		nDown, _ := toIntField(f["n_down"])
		nUp, _ := toIntField(f["n_up"])
		if nDown == 0 {
			t.Errorf("n_down = 0, want 5 (stressed < baseline for all 5; pre-fix code used cohens_d sign which is always ≥ 0); n_up = %d", nUp)
		}
		if nUp != 0 {
			t.Errorf("n_up = %d, want 0 (all metrics went down)", nUp)
		}
	}
}

// W10-4: emitTunableObservations must count skipped entries correctly.
// Pre-fix: always returned (emitted, 0) regardless of how many were skipped.
func TestTunableSkippedCount(t *testing.T) {
	// A store with no PALI, no shard, no replicaset signals → topology = "standalone".
	// tunableCatalog has no entries for "standalone" topology, so everything is skipped.
	store := &seriesStore{
		nameToID: map[string]int{},
		names:    []string{},
	}
	cfg := AutoConfig{RecommendTunables: true}
	c := &findingCollector{}
	c.add(newFinding("run_metadata", "h1", "run", map[string]interface{}{
		"sample_count": 1000,
	}))

	_, skipped := emitTunableObservations(c, "h1", store, cfg)
	if skipped == 0 && len(tunableCatalog) > 0 {
		t.Errorf("skipped = 0 but tunableCatalog has %d entries, at least some should be skipped for standalone topology; pre-fix always returned 0", len(tunableCatalog))
	}
}
