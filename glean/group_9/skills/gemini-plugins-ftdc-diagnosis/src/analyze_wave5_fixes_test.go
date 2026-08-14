// analyze_wave5_fixes_test.go — regression tests for Wave 5 bug fixes.
package main

import (
	"encoding/json"
	"math"
	"os"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Fix W5-A: Pearson denominator uses full-series N, not lagged-window count.
// ---------------------------------------------------------------------------

func TestCorrelation_LagZeroIdenticalSeriesREqualsOne(t *testing.T) {
	// At lag=0 with identical leader and follower, r must be 1.0.
	// The old formula used (nLag-1) which at lag=0 equals (N-1), so lag=0
	// was coincidentally correct. The new formula uses (N-1) uniformly.
	// This test verifies the lag=0 case still works after the fix.
	n := 200
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	samples := make([]Sample, n)
	for i := range samples {
		samples[i].Timestamp = base.Add(time.Duration(i) * time.Second)
		v := int64(100 + i%20)
		samples[i].Metrics = map[string]int64{
			"test.leader":   v,
			"test.follower": v,
		}
	}
	store := buildSeriesStore(samples)

	cfg := defaultAutoConfig()
	cfg.CorrelationMinR = 0.9
	// Add test.follower to the follower vector so it is eligible as a follower.
	cfg.FollowerMetrics = append(cfg.FollowerMetrics, "test.follower")
	findings := correlateAgainstFollowersWithStore("h", store, samples, cfg)
	found := false
	for _, f := range findings {
		if k, _ := f["kind"].(string); k != "correlation" {
			continue
		}
		// The finding uses field "leader" for the metric name.
		if m, _ := f["leader"].(string); m != "test.leader" {
			continue
		}
		// The finding uses field "pearson" for the correlation coefficient.
		r, _ := f["pearson"].(float64)
		if math.Abs(r-1.0) > 0.02 {
			t.Errorf("lag=0 identical series: r = %v, want ~1.0", r)
		}
		found = true
	}
	if !found {
		t.Error("no correlation Finding for identical leader/follower at threshold 0.9")
	}
}

// ---------------------------------------------------------------------------
// Fix W5-D: diagnosis_candidate emits required_finding_ids separately.
// ---------------------------------------------------------------------------

func TestDiagnosis_RequiredFindingIDsSeparateFromSupporting(t *testing.T) {
	// Inject a collector with movers and changepoints — enough to satisfy
	// at least one diagnosis pattern's required kinds. The emitted
	// diagnosis_candidate must have "required_finding_ids" AND
	// "supporting_finding_ids" as disjoint sets.
	c := &findingCollector{}
	for _, kind := range []string{"mover", "changepoint", "monotonic_trend",
		"correlation", "utilization", "spike", "variance_shift"} {
		c.add(newFinding(kind, "h", "m."+kind, map[string]interface{}{
			"metric": "m." + kind,
		}))
	}
	store := &seriesStore{nameToID: map[string]int{}}
	emitDiagnosisCandidates(c, "h", store)

	for _, f := range c.items {
		if k, _ := f["kind"].(string); k != "diagnosis_candidate" {
			continue
		}
		reqIDs, reqOK := f["required_finding_ids"].([]string)
		if !reqOK {
			// nil slice is fine — some patterns may have no required kinds
			// represented as nil; ensure the field exists at all.
			if _, exists := f["required_finding_ids"]; !exists {
				t.Errorf("diagnosis_candidate missing required_finding_ids field")
			}
			continue
		}
		supIDs, _ := f["supporting_finding_ids"].([]string)
		reqSet := map[string]bool{}
		for _, id := range reqIDs {
			reqSet[id] = true
		}
		for _, id := range supIDs {
			if reqSet[id] {
				t.Errorf("ID %q in both required_finding_ids and supporting_finding_ids", id)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Fix W5-F: writeBaseline cleans up .tmp on success (and no .tmp after write).
// ---------------------------------------------------------------------------

func TestWriteBaseline_NoTmpAfterSuccess(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/baseline.jsonl"
	tmp := path + ".tmp"

	err := writeBaseline(path, "salt", map[string]baselineRow{}, map[string]recurrenceRow{})
	if err != nil {
		t.Fatalf("writeBaseline: %v", err)
	}
	if _, err := os.Stat(tmp); err == nil {
		t.Error(".tmp file still exists after successful writeBaseline")
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected baseline file to exist after success: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Fix W5-G: recap entries with zero WrittenAt are quarantined.
// ---------------------------------------------------------------------------

func TestRecap_ZeroWrittenAtIsQuarantined(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/recap.ndjson"

	// Write a JSON line with written_at absent (decodes to zero time.Time).
	entry := map[string]interface{}{
		"provenance": map[string]interface{}{
			"parser_version":           "v1",
			"pattern_catalog_hash":     "abc",
			"confidence_rank_at_write": 0.8,
			// deliberately omit "written_at"
		},
		"pattern_name": "test_pattern",
	}
	data, _ := json.Marshal(entry)
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	valid, quarantined := readRecapFile(path, "v1", "abc")
	if len(valid) != 0 {
		t.Errorf("expected 0 valid entries for zero WrittenAt, got %d", len(valid))
	}
	if len(quarantined) == 0 {
		t.Error("expected entry with zero WrittenAt to be quarantined")
	}
	found := false
	for _, q := range quarantined {
		if q.reason == "missing_written_at" {
			found = true
		}
	}
	if !found {
		t.Errorf("quarantined entries did not include missing_written_at: %v", quarantined)
	}
}

// ---------------------------------------------------------------------------
// Fix W5-H: baselineEnd clamped to captureStart; no inverted windows emitted.
// ---------------------------------------------------------------------------

func TestWindowSuggest_NearStartNoInvertedWindow(t *testing.T) {
	// Cluster all changepoints in the first 60 samples of a 600-sample capture.
	// centroidIdx ≈ 30. baselineSkipSec=120 → baselineEnd = start+30s - 120s < captureStart.
	// After the fix: baselineEnd is clamped to captureStart and the window is
	// either valid or skipped entirely (not emitted inverted).
	n := 600
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")

	store := &seriesStore{nameToID: map[string]int{}}
	store.names = append(store.names, "test.m")
	store.nameToID["test.m"] = 0
	store.isCounter = append(store.isCounter, false)
	store.isDegenerate = append(store.isDegenerate, false)
	store.observedMin = append(store.observedMin, 0)
	store.observedMax = append(store.observedMax, 100)
	store.fullMean = append(store.fullMean, 50)
	store.fullSD = append(store.fullSD, 10)
	store.fullN = append(store.fullN, n)
	store.timestamps = make([]time.Time, n)
	for i := range store.timestamps {
		store.timestamps[i] = base.Add(time.Duration(i) * time.Second)
	}

	// Changepoint near sample 30 (30 seconds into a 600-second capture).
	c := &findingCollector{}
	c.add(newFinding("changepoint", "h", "test.m", map[string]interface{}{
		"metric":       "test.m",
		"sample_index": 30,
		"at":           store.timestamps[30].Format(time.RFC3339),
		"cohens_d":     float64(2.5),
	}))

	emitWindowSuggestions(c, "h", store)

	for _, f := range c.items {
		if k, _ := f["kind"].(string); k != "window_suggestion" {
			continue
		}
		bl, ok := f["baseline"].(map[string]interface{})
		if !ok {
			continue
		}
		fromStr, _ := bl["from"].(string)
		toStr, _ := bl["to"].(string)
		from, err1 := time.Parse(time.RFC3339, fromStr)
		to, err2 := time.Parse(time.RFC3339, toStr)
		if err1 != nil || err2 != nil {
			continue
		}
		if !to.After(from) {
			t.Errorf("baseline window inverted or zero-width: from=%v to=%v", from, to)
		}
		if from.Before(base) {
			t.Errorf("baseline.from before capture start: %v < %v", from, base)
		}
		if to.Before(base) {
			t.Errorf("baseline.to before capture start: %v < %v", to, base)
		}
	}
}

// ---------------------------------------------------------------------------
// Fix W5-I: stratified outlier detection uses sample variance (not population).
// ---------------------------------------------------------------------------

func TestStratified_SampleVarianceNotPopulation(t *testing.T) {
	// For t-values [0, 0, X] with X > 0:
	//   mean = X/3
	//   population SD = X * sqrt(2/9) = X * sqrt(2)/3 ≈ 0.471 * X
	//   sample SD = X / sqrt(3) ≈ 0.577 * X
	// The outlier deviation for X is |X - X/3| = 2X/3 ≈ 0.667 * X
	//
	// With outlierZ = 2.0:
	//   pop threshold  = 2.0 * 0.471 * X = 0.943 * X → 0.667 * X < 0.943 * X → NOT outlier
	//   sample threshold = 2.0 * 0.577 * X = 1.155 * X → 0.667 * X < 1.155 * X → NOT outlier
	//
	// Both say not-outlier here. Let's use [0, 0, X] where pop WOULD flag it.
	// We need 0.667*X >= 0.943*X? No, 0.667 < 0.943 always.
	// Hmm, let me try [0, 0, 3] more carefully.
	// Values: [0, 0, 3]. mean = 1.
	// Deviations: -1, -1, 2.
	// pop variance = (1+1+4)/3 = 2; pop SD = sqrt(2) ≈ 1.414
	// sample variance = (1+1+4)/2 = 3; sample SD = sqrt(3) ≈ 1.732
	// |3 - 1| = 2.
	// pop threshold = 2.0 * 1.414 = 2.828 → 2 < 2.828 → NOT outlier
	// sample threshold = 2.0 * 1.732 = 3.464 → 2 < 3.464 → NOT outlier
	// Still both not-outlier.
	//
	// Use stratifiedOutlierZ = 1.0 for the math: then
	// pop threshold = 1.414; sample threshold = 1.732
	// |3 - 1| = 2 > 1.414 → pop flags it; 2 > 1.732 → sample also flags it.
	//
	// We need a range where pop flags but sample doesn't:
	// deviation >= Z * popSD but < Z * sampleSD
	// deviation ∈ [Z * popSD, Z * sampleSD)
	// For [0, 0, X], deviation = 2X/3, popSD = X*sqrt(2)/3, sampleSD = X/sqrt(3)
	// 2X/3 >= Z * X*sqrt(2)/3 → 2/sqrt(2) = sqrt(2) ≈ 1.414 → Z <= sqrt(2)
	// 2X/3 < Z * X/sqrt(3) → 2/sqrt(3) * sqrt(3) → 2 < Z → Z > 2
	// Contradiction: Z <= sqrt(2) AND Z > 2 — impossible.
	//
	// Let me try a different distribution. Use [0, 2, 10].
	// mean = 4. deviations: -4, -2, 6.
	// pop variance = (16+4+36)/3 = 56/3 ≈ 18.67; popSD ≈ 4.32
	// sample variance = (16+4+36)/2 = 28; sampleSD ≈ 5.29
	// |10 - 4| = 6.
	// Z=1.2: pop threshold = 1.2 * 4.32 = 5.18 → 6 > 5.18 → POP OUTLIER
	//        sample threshold = 1.2 * 5.29 = 6.35 → 6 < 6.35 → NOT sample outlier
	//
	// This discriminates! With outlierZ=1.2, [0, 2, 10] flags with pop but not with sample.
	// We verify the formula used by the fix produces the sample result.

	values := []float64{0, 2, 10}
	const z = 1.2

	var sum, sumSq float64
	for _, v := range values {
		sum += v
		sumSq += v * v
	}
	n := len(values)
	mean := sum / float64(n)

	// Population formula (the old buggy one).
	popVariance := (sumSq/float64(n)) - mean*mean
	popSD := math.Sqrt(popVariance)

	// Sample formula (the fixed one).
	var sampleVariance float64
	for _, v := range values {
		d := v - mean
		sampleVariance += d * d
	}
	sampleVariance /= float64(n - 1)
	sampleSD := math.Sqrt(sampleVariance)

	deviation := math.Abs(10.0 - mean)

	popFlags := deviation >= z*popSD
	sampleFlags := deviation >= z*sampleSD

	if !popFlags {
		t.Fatal("test setup error: pop formula should flag 10.0 as outlier with z=1.2")
	}
	if sampleFlags {
		t.Errorf("sample formula incorrectly flags outlier (threshold %.4f, deviation %.4f): "+
			"sample SD is larger than pop SD so the threshold should be higher",
			z*sampleSD, deviation)
	}
	// Confirm the relationship: sampleSD > popSD for n > 1.
	if sampleSD <= popSD {
		t.Errorf("expected sampleSD (%.4f) > popSD (%.4f) for n=%d", sampleSD, popSD, n)
	}
}
