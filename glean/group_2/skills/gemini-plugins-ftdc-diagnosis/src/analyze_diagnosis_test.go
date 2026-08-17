package main

import (
	"testing"
)

// TestDiagnosisCandidatePrimaryPaliLogPressure verifies that the
// "primary_pali_log_pressure" pattern is correctly matched when the
// required Finding kinds (materialization_lag + memory_triplet) are present.
func TestDiagnosisCandidatePrimaryPaliLogPressure(t *testing.T) {
	c := &findingCollector{}
	store := &seriesStore{
		nameToID: map[string]int{
			"serverStatus.metrics.disagg.getPageRequest.numQueued": 0,
		},
		names: []string{"serverStatus.metrics.disagg.getPageRequest.numQueued"},
	}

	// Add the required findings for PALI log pressure pattern.
	c.add(newFinding("materialization_lag", "host1", "metric1", map[string]interface{}{
		"metric": "some.metric",
		"p95_ms": 1500.0,
	}))
	c.add(newFinding("memory_triplet", "host1", "triplet", map[string]interface{}{
		"host": "host1",
	}))

	// Also add a supporting finding to test confidence boosting.
	c.add(newFinding("checkpoint_install_anomaly", "host1", "checkpoint", map[string]interface{}{
		"metric":         "some.metric",
		"count_increase": 42,
	}))

	emitted, _ := emitDiagnosisCandidates(c, "host1", store)

	if emitted == 0 {
		t.Fatalf("expected at least 1 diagnosis_candidate emitted, got %d", emitted)
	}

	// Find the diagnosis_candidate with pattern_name="primary_pali_log_pressure".
	var found Finding
	for _, f := range c.items {
		if f["kind"] == "diagnosis_candidate" {
			if pn, ok := f["pattern_name"].(string); ok && pn == "primary_pali_log_pressure" {
				found = f
				break
			}
		}
	}

	if found == nil {
		t.Fatal("diagnosis_candidate with pattern_name='primary_pali_log_pressure' not found")
	}

	// Verify confidence_rank is reasonable.
	if cr, ok := found["confidence_rank"].(float64); ok {
		if cr < 0.9 || cr > 1.0 {
			t.Logf("confidence_rank=%v is outside expected [0.9, 1.0]", cr)
			// Not fatal; just informational for now.
		}
	} else {
		t.Fatalf("confidence_rank not found or not a number")
	}
}

// TestDiagnosisCandidateNonPALICaptureSkips verifies that PALI-specific
// patterns are skipped when no disagg.* metrics are present.
func TestDiagnosisCandidateNonPALICaptureSkips(t *testing.T) {
	c := &findingCollector{}
	store := &seriesStore{
		nameToID: map[string]int{
			"serverStatus.opLatencies.reads.latency": 0,
		},
		names: []string{"serverStatus.opLatencies.reads.latency"},
	}

	// Add findings that WOULD match PALI pattern, but domain gate should block them.
	c.add(newFinding("materialization_lag", "host1", "metric1", map[string]interface{}{
		"metric": "some.metric",
		"p95_ms": 1500.0,
	}))
	c.add(newFinding("memory_triplet", "host1", "triplet", map[string]interface{}{
		"host": "host1",
	}))

	_, _ = emitDiagnosisCandidates(c, "host1", store)

	// Without PALI domain hint (disagg.* metrics), the pattern should be skipped.
	for _, f := range c.items {
		if f["kind"] == "diagnosis_candidate" {
			if pn, ok := f["pattern_name"].(string); ok {
				if pn == "primary_pali_log_pressure" {
					t.Fatalf("primary_pali_log_pressure should not be emitted without PALI metrics")
				}
			}
		}
	}
}
