package main

import (
	"testing"
	"time"
)

// makeStoreForCoverage builds a seriesStore with the named metrics and
// their observed min/max. Pass same value for min and max to simulate a
// "silent" (all-zero-range) metric; different values for an active one.
func makeStoreForCoverage(metrics []string, minVals, maxVals []int64) *seriesStore {
	s := &seriesStore{nameToID: map[string]int{}}
	for i, name := range metrics {
		s.names = append(s.names, name)
		s.nameToID[name] = i
	}
	s.isCounter = make([]bool, len(metrics))
	s.isDegenerate = make([]bool, len(metrics))
	s.observedMin = minVals
	s.observedMax = maxVals
	s.fullMean = make([]float64, len(metrics))
	s.fullSD = make([]float64, len(metrics))
	s.materializedSeries = make([][]float64, len(metrics))
	// One timestamp so the store is non-empty.
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	s.timestamps = []time.Time{base}
	for range metrics {
		col := NewCompressedColumn(1)
		col.Add(0, true)
		col.Finish()
		s.columns = append(s.columns, col)
	}
	return s
}

// TestEmitCoverageMap_GapForAbsentArea: an expected-active area that has
// ZERO metrics in the store (not merely all-static) must still emit a
// coverage_gap Finding with reason "entirely_absent".
//
// Before the fix, only areas present in silentSet (metrics present but
// all-zero-range) were flagged; completely absent areas were silently
// ignored, producing a false negative.
func TestEmitCoverageMap_GapForAbsentArea(t *testing.T) {
	// Build a replica_set capture: has repl_state and core_opcounters
	// metrics present (and active), but core_oplatencies and core_queues
	// are completely absent from the store.
	metrics := []string{
		"serverStatus.repl.electionId",          // repl_state (active)
		"serverStatus.opcounters.insert",         // core_opcounters (active)
		"serverStatus.opcounters.query",          // core_opcounters (active)
	}
	minVals := []int64{1, 0, 0}
	maxVals := []int64{2, 100, 50} // all have nonzero range → active

	store := makeStoreForCoverage(metrics, minVals, maxVals)
	c := &findingCollector{}
	emitted, _ := emitCoverageMap(c, "host1", store)
	if emitted == 0 {
		t.Fatal("expected at least one Finding from emitCoverageMap, got 0")
	}

	// Collect coverage_gap kinds.
	var gapAreas []string
	gapReasons := map[string]string{}
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k != "coverage_gap" {
			continue
		}
		area, _ := f["area"].(string)
		reason, _ := f["reason"].(string)
		gapAreas = append(gapAreas, area)
		gapReasons[area] = reason
	}

	// For replica_set topology expectedActive =
	// {"core_opcounters","core_oplatencies","core_queues","repl_state"}.
	// core_oplatencies and core_queues are entirely absent.
	for _, expected := range []string{"core_oplatencies", "core_queues"} {
		found := false
		for _, a := range gapAreas {
			if a == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected coverage_gap for absent area %q, but none was emitted (gap areas: %v)", expected, gapAreas)
			continue
		}
		if gapReasons[expected] != "entirely_absent" {
			t.Errorf("area %q: reason = %q, want entirely_absent", expected, gapReasons[expected])
		}
	}
}

// TestEmitCoverageMap_GapForAllStaticArea: an area present in the store
// but entirely static (min==max for all metrics) must still emit
// coverage_gap with reason "all_static" (pre-existing behavior; regression guard).
func TestEmitCoverageMap_GapForAllStaticArea(t *testing.T) {
	// replica_set with core_oplatencies present but all zeros.
	metrics := []string{
		"serverStatus.repl.electionId",       // repl_state active
		"serverStatus.opcounters.insert",      // core_opcounters active
		"serverStatus.opcounters.query",       // core_opcounters active
		"serverStatus.opLatencies.reads.ops",  // core_oplatencies — static
		"serverStatus.queues.execution.read.out", // core_queues — static
	}
	minVals := []int64{1, 0, 0, 0, 0}
	maxVals := []int64{2, 100, 50, 0, 0} // last two: min==max → silent

	store := makeStoreForCoverage(metrics, minVals, maxVals)
	c := &findingCollector{}
	emitCoverageMap(c, "host1", store)

	gapReasons := map[string]string{}
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k != "coverage_gap" {
			continue
		}
		area, _ := f["area"].(string)
		reason, _ := f["reason"].(string)
		gapReasons[area] = reason
	}

	for _, area := range []string{"core_oplatencies", "core_queues"} {
		r, ok := gapReasons[area]
		if !ok {
			t.Errorf("expected coverage_gap for static area %q, but none emitted", area)
			continue
		}
		if r != "all_static" {
			t.Errorf("area %q: reason = %q, want all_static", area, r)
		}
	}
}
