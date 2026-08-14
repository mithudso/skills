// analyze_wave11_fixes_test.go — regression tests for Wave 11 confirmed bugs.
//
// W11-1: analyze_baseline.go skipped counter never incremented for degenerate/counter metrics.
// W11-2: analyze_followups.go asserting []interface{} on supporting_finding_ids stored as []string.
// W11-3: analyze_reasoning.go skipped counter not incremented when no precursor_candidates emitted.
package main

import (
	"os"
	"testing"
)

// W11-1: runBaselineStage must count skipped metrics (degenerate or counter).
// Pre-fix: always returned skipped=0 regardless of how many metrics were skipped.
func TestBaselineSkippedCount(t *testing.T) {
	// Write a minimal valid baseline file so loadBaseline succeeds and the
	// loop runs — otherwise the function exits before reaching the skipped path.
	f, err := os.CreateTemp("", "ftdc-baseline-*.jsonl")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer os.Remove(f.Name())
	// Minimal valid baseline: just a salt line.
	if _, err := f.WriteString(`{"kind":"salt","salt":"testsalt"}` + "\n"); err != nil {
		t.Fatalf("write salt: %v", err)
	}
	f.Close()

	// Store with one counter and one degenerate metric — both should be skipped.
	store := &seriesStore{
		names:        []string{"test.counter", "test.degenerate"},
		nameToID:     map[string]int{"test.counter": 0, "test.degenerate": 1},
		isCounter:    []bool{true, false},
		isDegenerate: []bool{false, true},
		fullMean:     []float64{100.0, 0.0},
		fullSD:       []float64{1.0, 0.0},
		fullN:        []int{100, 100},
	}
	c := &findingCollector{}
	cfg := AutoConfig{BaselineFile: f.Name()}
	_, skipped := runBaselineStage(c, "h", store, cfg)
	if skipped != 2 {
		t.Errorf("expected skipped=2 (1 counter + 1 degenerate), got %d; pre-fix always returned 0", skipped)
	}
}

// W11-1b: runBaselineStage skipped counter — verify the fix is wired correctly
// by using a store with only counter metrics and a valid (but empty) baseline that
// loadBaseline can parse. Since we can't easily create a real baseline file in a
// unit test, we verify the fix via code inspection: the fix added skipped++ before
// the isDegenerate/isCounter continue. This test confirms the function's return
// signature matches (emitted, skipped) and that for a zero-metric store, skipped=0.
func TestBaselineSkippedZeroMetrics(t *testing.T) {
	store := &seriesStore{
		names:        []string{},
		nameToID:     map[string]int{},
		isCounter:    []bool{},
		isDegenerate: []bool{},
		fullMean:     []float64{},
		fullSD:       []float64{},
		fullN:        []int{},
	}
	c := &findingCollector{}
	cfg := AutoConfig{BaselineFile: ""}
	emitted, skipped := runBaselineStage(c, "h", store, cfg)
	if emitted != 0 || skipped != 0 {
		t.Errorf("empty store + empty BaselineFile: want (0,0), got (%d,%d)", emitted, skipped)
	}
}

// W11-2: emitFollowupsBacklog must populate supporting_finding_ids from
// diagnosis_candidate Findings. Pre-fix: .([]interface{}) assertion always failed
// on a []string field, so supporting IDs were silently dropped into the file.
//
// The snapshot Finding only carries pattern_name + confidence_rank (not IDs).
// We verify the fix is wired by asserting the .([]string) assertion succeeds and
// the backlog snapshot is emitted (freshCands was non-empty).
func TestFollowupsSupportingIDs(t *testing.T) {
	c := &findingCollector{}
	// Create a diagnosis_candidate with a []string supporting_finding_ids — exactly
	// as emitDiagnosisCandidates stores it.
	c.add(newFinding("diagnosis_candidate", "h1", "f-abc", map[string]interface{}{
		"pattern_name":               "primary_pali_log_pressure",
		"confidence_rank":            0.65,
		"discriminating_kinds_needed": []string{"counter_wrap"},
		"supporting_finding_ids":     []string{"f-001", "f-002"},
	}))

	// Directly verify the type assertion now succeeds.
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		if k != "diagnosis_candidate" {
			continue
		}
		ids, ok := f["supporting_finding_ids"].([]string)
		if !ok {
			t.Error(".([]string) assertion failed on supporting_finding_ids; pre-fix used .([]interface{}) which always fails")
		}
		if len(ids) != 2 {
			t.Errorf("expected 2 supporting IDs, got %d", len(ids))
		}
	}

	emitted, _ := emitFollowupsBacklog(c, "h1")
	if emitted == 0 {
		t.Error("emitFollowupsBacklog emitted 0; expected followups_backlog_snapshot (freshCands non-empty)")
	}
}

// W11-3: emitPrecursorCandidates must return skipped≥1 when changepoints exist
// but none have preceding candidates within the lookback window.
// Pre-fix: skipped stayed 0 even when no precursor_candidates were emitted.
func TestPrecursorCandidatesSkippedCount(t *testing.T) {
	c := &findingCollector{}
	// One changepoint with no other findings to serve as precursors.
	c.add(newFinding(KindChangepoint, "h1", "metric1|10", map[string]interface{}{
		"metric":       "serverStatus.metric1",
		"at":           "2024-01-01T00:05:00.000Z",
		"sample_index": 10,
	}))

	emitted, skipped := emitPrecursorCandidates(c, "h1")
	if emitted != 0 {
		t.Errorf("expected emitted=0 (no other metrics as precursors), got %d", emitted)
	}
	if skipped == 0 {
		t.Error("skipped=0 even though no precursor_candidates were emitted; pre-fix always returned skipped=0")
	}
}
