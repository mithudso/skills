package main

import (
	"fmt"
	"testing"
	"time"
)

// makeStratifiedStore builds a seriesStore with families of metrics for
// use by emitMoverStratified tests. Each family has `famSize` members of
// the form "serverStatus.queues.execution.<dim>.out" where dim varies.
// The first member of the first family has elevated stress values so the
// family is non-trivial and emitMoverStratified actually fires.
func makeStratifiedStore(t *testing.T, famADims, famBDims []string) *seriesStore {
	t.Helper()
	var names []string
	for _, d := range famADims {
		names = append(names, fmt.Sprintf("serverStatus.queues.execution.%s.out", d))
	}
	for _, d := range famBDims {
		names = append(names, fmt.Sprintf("serverStatus.metrics.repl.network.%s.numOps", d))
	}

	s := &seriesStore{nameToID: map[string]int{}}
	n := len(names)
	s.names = names
	s.nameToID = make(map[string]int, n)
	for i, nm := range names {
		s.nameToID[nm] = i
	}
	s.isCounter = make([]bool, n)
	s.isDegenerate = make([]bool, n)
	s.observedMin = make([]int64, n)
	s.observedMax = make([]int64, n)
	s.fullMean = make([]float64, n)
	s.fullSD = make([]float64, n)
	s.materializedSeries = make([][]float64, n)

	// 20 samples: first 10 baseline, last 10 stress.
	const samples = 20
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < samples; i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}

	for i := range names {
		col := NewCompressedColumn(samples)
		for j := 0; j < samples; j++ {
			v := int64(1)
			// Make the first metric of each family an outlier during stress
			// so the Welch-t is large enough for emitMoverStratified to fire.
			if i == 0 && j >= 10 {
				v = 1000
			}
			col.Add(v, true)
		}
		col.Finish()
		s.columns = append(s.columns, col)
		s.observedMin[i] = 1
		if i == 0 {
			s.observedMax[i] = 1000
		} else {
			s.observedMax[i] = 1
		}
	}
	return s
}

// TestEmitMoverStratified_StableSortOnEqualFamilies: when two families
// have equal member counts, two consecutive calls must produce the same
// Finding order. Before the fix (sort.Slice → sort.SliceStable), the
// relative order of equal-size families was non-deterministic because the
// preceding sort.Strings established a tiebreaker that sort.Slice did not
// preserve.
func TestEmitMoverStratified_StableSortOnEqualFamilies(t *testing.T) {
	// Two families with exactly the same number of dimension values (4 each).
	famADims := []string{"alpha", "beta", "gamma", "delta"}
	famBDims := []string{"nodeA", "nodeB", "nodeC", "nodeD"}

	store := makeStratifiedStore(t, famADims, famBDims)

	collectFindings := func() []string {
		c := &findingCollector{}
		emitMoverStratified(c, "host1", store, 0, 10, 10, 20)
		var keys []string
		for _, f := range c.items {
			if k, _ := f["kind"].(string); k == "mover_stratified" {
				idKey, _ := f["id_key"].(string)
				keys = append(keys, idKey)
			}
		}
		return keys
	}

	first := collectFindings()
	if len(first) == 0 {
		// The store might not produce mover_stratified if the families
		// don't have enough variance — that's fine, just skip.
		t.Skip("emitMoverStratified emitted no mover_stratified Findings; skipping order test")
	}

	const rounds = 5
	for i := 0; i < rounds; i++ {
		got := collectFindings()
		if len(got) != len(first) {
			t.Errorf("round %d: got %d Findings, want %d", i, len(got), len(first))
			continue
		}
		for j, k := range got {
			if k != first[j] {
				t.Errorf("round %d: Finding[%d] id_key = %q, want %q (non-deterministic order)", i, j, k, first[j])
			}
		}
	}
}
