package main

import (
	"testing"
	"time"
)

// TestEmitWorkloadSpec_RateMath: emits a Finding with per-second rates
// derived from cumulative-counter first/last delta over capture
// duration.
func TestEmitWorkloadSpec_RateMath(t *testing.T) {
	s := &seriesStore{nameToID: map[string]int{}}
	// Insert counter grows by 1000 over 100 seconds → rate = 10/s.
	s.names = []string{"serverStatus.opcounters.insert"}
	s.nameToID["serverStatus.opcounters.insert"] = 0
	s.isCounter = []bool{true}
	s.isDegenerate = []bool{false}
	s.observedMin = []int64{0}
	s.observedMax = []int64{0}
	s.fullMean = []float64{0}
	s.fullSD = []float64{0}
	s.materializedSeries = make([][]float64, 1)
	col := NewCompressedColumn(101)
	for i := 0; i <= 100; i++ {
		col.Add(int64(i*10), true)
	}
	col.Finish()
	s.columns = []*CompressedColumn{col}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i <= 100; i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	c := &findingCollector{}
	emitted, _ := emitWorkloadSpec(c, "h", s)
	if emitted != 1 {
		t.Fatalf("expected 1 workload_spec Finding, got %d", emitted)
	}
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "workload_spec" {
			if r, _ := f["insert_rate_ps"].(float64); r < 9.9 || r > 10.1 {
				t.Errorf("insert_rate_ps: got %v, want ~10", r)
			}
		}
	}
}

// TestEmitWorkloadSpec_NoCounters: store with no opcounter metrics →
// no Finding (no useful fields to emit).
func TestEmitWorkloadSpec_NoCounters(t *testing.T) {
	s := &seriesStore{nameToID: map[string]int{}}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < 10; i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	c := &findingCollector{}
	emitted, _ := emitWorkloadSpec(c, "h", s)
	if emitted != 0 {
		t.Errorf("expected 0 emits when no counters available, got %d", emitted)
	}
}

// TestEmitPrecursorCandidates_PrecedenceRanking: for a changepoint at
// T, precursors are Findings on OTHER metrics that fired EARLIER
// within 300s, ranked by magnitude desc.
func TestEmitPrecursorCandidates_PrecedenceRanking(t *testing.T) {
	c := &findingCollector{}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	// Three precursor spikes on different metrics, then a target
	// changepoint at T+200s.
	at := func(off int) string {
		return base.Add(time.Duration(off) * time.Second).UTC().Format("2006-01-02T15:04:05.000Z")
	}
	c.add(Finding{
		"kind": KindSpike, "id": "f-a", "metric": "m.a", "at": at(50),
		"z_score": 12.0, "host": "h",
	})
	c.add(Finding{
		"kind": KindSpike, "id": "f-b", "metric": "m.b", "at": at(100),
		"z_score": 5.0, "host": "h",
	})
	c.add(Finding{
		"kind": KindVarianceShift, "id": "f-c", "metric": "m.c", "at": at(150),
		"cusum_peak": 3.0, "host": "h",
	})
	c.add(Finding{
		"kind": KindChangepoint, "id": "f-target", "metric": "m.target", "at": at(200),
		"effect_size": 3.0, "host": "h",
	})
	emitted, _ := emitPrecursorCandidates(c, "h")
	if emitted != 1 {
		t.Fatalf("expected 1 precursor_candidates Finding, got %d", emitted)
	}
	var f Finding
	for _, x := range c.items {
		if k, _ := x["kind"].(string); k == "precursor_candidates" {
			f = x
		}
	}
	if f == nil {
		t.Fatal("expected a precursor_candidates Finding")
	}
	pres, _ := f["precursors"].([]interface{})
	if len(pres) != 3 {
		t.Errorf("expected 3 precursors, got %d", len(pres))
	}
	// Top precursor should be the spike with highest magnitude (m.a, z=12).
	first, _ := pres[0].(map[string]interface{})
	if m, _ := first["metric"].(string); m != "m.a" {
		t.Errorf("top precursor metric: got %q, want m.a", m)
	}
}

// TestEmitPrecursorCandidates_ExcludesSameMetric: a precursor on the
// same metric as the target is filtered out (avoid self-precedence).
func TestEmitPrecursorCandidates_ExcludesSameMetric(t *testing.T) {
	c := &findingCollector{}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	at := func(off int) string {
		return base.Add(time.Duration(off) * time.Second).UTC().Format("2006-01-02T15:04:05.000Z")
	}
	c.add(Finding{"kind": KindSpike, "id": "f-same", "metric": "m.x", "at": at(50), "z_score": 10.0, "host": "h"})
	c.add(Finding{"kind": KindSpike, "id": "f-other", "metric": "m.y", "at": at(60), "z_score": 5.0, "host": "h"})
	c.add(Finding{"kind": KindChangepoint, "id": "f-target", "metric": "m.x", "at": at(200), "effect_size": 3.0, "host": "h"})
	emitPrecursorCandidates(c, "h")
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "precursor_candidates" {
			pres, _ := f["precursors"].([]interface{})
			for _, p := range pres {
				pm, _ := p.(map[string]interface{})
				if m, _ := pm["metric"].(string); m == "m.x" {
					t.Errorf("same-metric precursor leaked into output: m.x")
				}
			}
		}
	}
}

// TestEmitPrecursorCandidates_OutsideLookback: a precursor older than
// 300s is excluded.
func TestEmitPrecursorCandidates_OutsideLookback(t *testing.T) {
	c := &findingCollector{}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	at := func(off int) string {
		return base.Add(time.Duration(off) * time.Second).UTC().Format("2006-01-02T15:04:05.000Z")
	}
	// Precursor at t=0; target at t=500. Lookback=300s → precursor
	// excluded.
	c.add(Finding{"kind": KindSpike, "id": "f-far", "metric": "m.a", "at": at(0), "z_score": 10.0, "host": "h"})
	c.add(Finding{"kind": KindChangepoint, "id": "f-target", "metric": "m.target", "at": at(500), "effect_size": 3.0, "host": "h"})
	emitted, _ := emitPrecursorCandidates(c, "h")
	if emitted != 0 {
		t.Errorf("expected 0 precursor_candidates Finding (only precursor outside window), got %d", emitted)
	}
}

// TestEmitPrecursorCandidates_OnlyForChangepointTargets: a spike
// shouldn't be an anchor for precursor ranking — only changepoints.
func TestEmitPrecursorCandidates_OnlyForChangepointTargets(t *testing.T) {
	c := &findingCollector{}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	at := func(off int) string {
		return base.Add(time.Duration(off) * time.Second).UTC().Format("2006-01-02T15:04:05.000Z")
	}
	c.add(Finding{"kind": KindSpike, "id": "f-a", "metric": "m.a", "at": at(50), "z_score": 10.0, "host": "h"})
	c.add(Finding{"kind": KindSpike, "id": "f-target", "metric": "m.target", "at": at(200), "z_score": 12.0, "host": "h"})
	emitted, _ := emitPrecursorCandidates(c, "h")
	if emitted != 0 {
		t.Errorf("expected no precursor_candidates Finding when target is a spike (not changepoint), got %d", emitted)
	}
}
