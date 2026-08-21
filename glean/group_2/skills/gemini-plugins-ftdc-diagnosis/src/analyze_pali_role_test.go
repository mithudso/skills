package main

import (
	"testing"
	"time"
)

// buildPALIRoleStore constructs a seriesStore with the given named metrics
// and per-metric value slices. All slices must be the same length.
func buildPALIRoleStore(names []string, valsByMetric [][]int64) *seriesStore {
	n := len(names)
	s := &seriesStore{nameToID: map[string]int{}}
	s.names = make([]string, n)
	s.isCounter = make([]bool, n)
	s.isDegenerate = make([]bool, n)
	s.observedMin = make([]int64, n)
	s.observedMax = make([]int64, n)
	s.fullMean = make([]float64, n)
	s.fullSD = make([]float64, n)
	s.materializedSeries = make([][]float64, n)
	s.columns = make([]*CompressedColumn, n)
	var nSamples int
	for i, name := range names {
		s.names[i] = name
		s.nameToID[name] = i
		vals := valsByMetric[i]
		if nSamples == 0 {
			nSamples = len(vals)
		}
		col := NewCompressedColumn(len(vals))
		for _, v := range vals {
			col.Add(v, true)
		}
		col.Finish()
		s.columns[i] = col
	}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < nSamples; i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	return s
}

// TestPALIRole_ZeroLagStandby exercises D5: a standby that has fully caught
// up (PhylogMaterLagNonZeroN == 0) must still be classified as standby, not
// unknown — the presence of the lag metric alone is sufficient.
func TestPALIRole_ZeroLagStandby(t *testing.T) {
	const (
		lagMetric    = "serverStatus.metrics.disagg.phylog.committedToMaterializedLagMillis"
		// A disagg anchor metric is needed for isPALICapture to gate in.
		anchorMetric = "serverStatus.metrics.disagg.logServer.appendLogTotal"
	)
	// Build a store where the lag metric is always 0 (fully caught up)
	// and appendLogSuccess is absent (so primaryLikely is false).
	n := 20
	lagVals := make([]int64, n)    // all zeros: lag is always 0
	anchorVals := make([]int64, n) // constant, just for the isPALI gate
	for i := 0; i < n; i++ {
		anchorVals[i] = 100
	}
	store := buildPALIRoleStore(
		[]string{anchorMetric, lagMetric},
		[][]int64{anchorVals, lagVals},
	)
	role, ev := paliRole(store)
	if role != PALIRoleStandby {
		t.Errorf("paliRole: got %q, want %q; evidence=%+v", role, PALIRoleStandby, ev)
	}
	if !ev.HasPhylogMaterLag {
		t.Errorf("evidence.HasPhylogMaterLag should be true")
	}
	if ev.PhylogMaterLagNonZeroN != 0 {
		t.Errorf("evidence.PhylogMaterLagNonZeroN should be 0 (zero-lag standby), got %d", ev.PhylogMaterLagNonZeroN)
	}
}

// TestPALIRole_ActiveLagStandby: standby with non-zero lag samples is
// still PALIRoleStandby (regression guard for the pre-fix behaviour).
func TestPALIRole_ActiveLagStandby(t *testing.T) {
	const (
		lagMetric    = "serverStatus.metrics.disagg.phylog.committedToMaterializedLagMillis"
		anchorMetric = "serverStatus.metrics.disagg.logServer.appendLogTotal"
	)
	n := 20
	lagVals := make([]int64, n)
	anchorVals := make([]int64, n)
	for i := 0; i < n; i++ {
		anchorVals[i] = 100
		lagVals[i] = int64(50 + i*10) // non-zero lag throughout
	}
	store := buildPALIRoleStore(
		[]string{anchorMetric, lagMetric},
		[][]int64{anchorVals, lagVals},
	)
	role, ev := paliRole(store)
	if role != PALIRoleStandby {
		t.Errorf("paliRole: got %q, want %q; evidence=%+v", role, PALIRoleStandby, ev)
	}
	if ev.PhylogMaterLagNonZeroN == 0 {
		t.Errorf("evidence.PhylogMaterLagNonZeroN should be > 0 for active-lag standby")
	}
}
