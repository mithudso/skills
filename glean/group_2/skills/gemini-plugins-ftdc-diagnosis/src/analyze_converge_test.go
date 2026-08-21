package main

import (
	"os"
	"path/filepath"
	"testing"
)

// makeSiblingDir creates a directory that passes looksLikeCaptureDir.
func makeSiblingDir(t *testing.T, parent, name string) string {
	t.Helper()
	dir := filepath.Join(parent, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	// Drop a non-golden file so looksLikeCaptureDir returns true.
	if err := os.WriteFile(filepath.Join(dir, "metrics.bson"), []byte("placeholder"), 0o644); err != nil {
		t.Fatalf("write placeholder: %v", err)
	}
	return dir
}

// TestEmitConvergePlan_SelectsHighestConfidenceCandidate: when multiple
// diagnosis_candidate Findings exist in [0.3, 0.9), emitConvergePlan
// must attach the converge_plan to the HIGHEST-confidence one (closest
// to 0.9), not the lowest-confidence one.
//
// Before the fix the comparison was `cr < targetCR`, which selected the
// minimum confidence. The fix changes it to `cr > targetCR`.
func TestEmitConvergePlan_SelectsHighestConfidenceCandidate(t *testing.T) {
	// discoverSiblingCaptures(inputPath) calls filepath.Dir(inputPath) to
	// find the parent, then walks that parent for peer subdirectories.
	// Structure:
	//   captureParent/            ← parent directory walked for siblings
	//     metrics.bson            ← the input file (so Dir = captureParent)
	//     sibling1/               ← peer capture dir
	//     sibling2/               ← peer capture dir
	captureParent := t.TempDir()
	// Write the "input" file directly in captureParent.
	inputPath := filepath.Join(captureParent, "metrics.bson")
	if err := os.WriteFile(inputPath, []byte("placeholder"), 0o644); err != nil {
		t.Fatalf("write input file: %v", err)
	}
	_ = makeSiblingDir(t, captureParent, "sibling1")
	_ = makeSiblingDir(t, captureParent, "sibling2")

	// Inject three diagnosis_candidate Findings with different confidence
	// ranks in the valid [0.3, 0.9) window.
	c := &findingCollector{}
	c.add(Finding{
		"kind":            "diagnosis_candidate",
		"host":            "h",
		"id":              "cand-low",
		"pattern_name":    "pattern_low",
		"confidence_rank": float64(0.35),
	})
	c.add(Finding{
		"kind":            "diagnosis_candidate",
		"host":            "h",
		"id":              "cand-mid",
		"pattern_name":    "pattern_mid",
		"confidence_rank": float64(0.60),
	})
	c.add(Finding{
		"kind":            "diagnosis_candidate",
		"host":            "h",
		"id":              "cand-high",
		"pattern_name":    "pattern_high",
		"confidence_rank": float64(0.82),
	})

	store := &seriesStore{nameToID: map[string]int{}}
	cfg := AutoConfig{InputPath: inputPath}

	emitted, _ := emitConvergePlan(c, "h", store, cfg)
	if emitted == 0 {
		t.Fatal("expected at least one Finding from emitConvergePlan, got 0")
	}

	// Locate the converge_plan Finding.
	var plan Finding
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "converge_plan" {
			plan = f
			break
		}
	}
	if plan == nil {
		// Dump what was emitted for diagnosis.
		for _, f := range c.items {
			t.Logf("emitted: %v", f)
		}
		t.Fatal("no converge_plan Finding emitted")
	}

	targetID, _ := plan["target_candidate_id"].(string)
	if targetID != "cand-high" {
		t.Errorf("target_candidate_id = %q, want cand-high (highest-confidence candidate below 0.9)", targetID)
	}
	targetCR, _ := plan["current_confidence_rank"].(float64)
	if targetCR < 0.81 || targetCR > 0.83 {
		t.Errorf("current_confidence_rank = %v, want ~0.82", targetCR)
	}
}

// TestEmitConvergePlan_IgnoresBelowFloor: candidates with confidence < 0.3
// must not be selected.
func TestEmitConvergePlan_IgnoresBelowFloor(t *testing.T) {
	// Below-floor test only checks that no converge_plan is emitted;
	// we don't need siblings for that path — the function exits early
	// at "no_low_confidence_candidate".
	captureParent := t.TempDir()
	inputPath := filepath.Join(captureParent, "metrics.bson")
	if err := os.WriteFile(inputPath, []byte("placeholder"), 0o644); err != nil {
		t.Fatalf("write input file: %v", err)
	}

	c := &findingCollector{}
	// Only a below-floor candidate — should produce convergence_failed.
	c.add(Finding{
		"kind":            "diagnosis_candidate",
		"host":            "h",
		"id":              "cand-tiny",
		"pattern_name":    "pattern_tiny",
		"confidence_rank": float64(0.1),
	})

	store := &seriesStore{nameToID: map[string]int{}}
	cfg := AutoConfig{InputPath: inputPath}

	emitConvergePlan(c, "h", store, cfg)

	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "converge_plan" {
			t.Errorf("expected no converge_plan for below-floor candidates, but got one: %v", f)
		}
	}
}
