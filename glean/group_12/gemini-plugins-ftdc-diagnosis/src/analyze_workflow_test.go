package main

import (
	"testing"
	"time"
)

// TestEmitMemoryTriplet_AllThreeMetrics: a store with all 3 canonical
// memory metrics emits one Finding with slopes for each.
func TestEmitMemoryTriplet_AllThreeMetrics(t *testing.T) {
	s := &seriesStore{nameToID: map[string]int{}}
	metrics := []string{
		"serverStatus.tcmalloc.generic.current_allocated_bytes",
		"serverStatus.mem.resident",
		"serverStatus.tcmalloc.generic.heap_size",
	}
	for i, n := range metrics {
		s.names = append(s.names, n)
		s.nameToID[n] = i
		s.isCounter = append(s.isCounter, false)
		s.isDegenerate = append(s.isDegenerate, false)
		s.observedMin = append(s.observedMin, 0)
		s.observedMax = append(s.observedMax, 0)
		s.fullMean = append(s.fullMean, 0)
		s.fullSD = append(s.fullSD, 0)
	}
	s.materializedSeries = make([][]float64, 3)
	// 100 samples over 100 seconds.
	// current_allocated_bytes: 1e9 + i * 1e6 → slope = 1e6 bytes/sec
	// mem.resident_mb: 1000 + i*0 (steady) → slope = 0
	// heap_size: 2e9 + i * 2e6 → slope = 2e6 bytes/sec
	for i := 0; i < 3; i++ {
		col := NewCompressedColumn(100)
		for j := 0; j < 100; j++ {
			var v int64
			switch i {
			case 0:
				v = int64(1e9) + int64(j)*int64(1e6)
			case 1:
				v = 1000
			case 2:
				v = int64(2e9) + int64(j)*int64(2e6)
			}
			col.Add(v, true)
		}
		col.Finish()
		s.columns = append(s.columns, col)
	}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for j := 0; j < 100; j++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(j)*time.Second))
	}
	c := &findingCollector{}
	emitted, _ := emitMemoryTriplet(c, "h", s)
	if emitted != 1 {
		t.Fatalf("expected 1 memory_triplet Finding, got %d", emitted)
	}
	var f Finding
	for _, x := range c.items {
		if k, _ := x["kind"].(string); k == "memory_triplet" {
			f = x
		}
	}
	if f == nil {
		t.Fatal("expected memory_triplet Finding")
	}
	if s, _ := f["current_allocated_bytes_slope_per_sec"].(float64); s < 0.99e6 || s > 1.01e6 {
		t.Errorf("current_allocated_bytes_slope_per_sec: got %v, want ~1e6", s)
	}
	if s, _ := f["mem_resident_mb_slope_per_sec"].(float64); s != 0 {
		t.Errorf("mem_resident_mb_slope_per_sec: got %v, want 0", s)
	}
}

// TestEmitMemoryTriplet_OnlyOneMetric: with only one of the three
// metrics present, the Finding is NOT emitted (need ≥2 for joint shape).
func TestEmitMemoryTriplet_OnlyOneMetric(t *testing.T) {
	s := makeStoreWithSingle("serverStatus.mem.resident", []int64{1000, 1001, 1002})
	c := &findingCollector{}
	emitted, skipped := emitMemoryTriplet(c, "h", s)
	if emitted != 0 {
		t.Errorf("expected 0 with only 1 metric, got %d", emitted)
	}
	if skipped != 1 {
		t.Errorf("expected skipped=1, got %d", skipped)
	}
}

// TestEmitRecurrenceFingerprint_Deterministic: two runs with identical
// inputs produce identical fingerprints.
func TestEmitRecurrenceFingerprint_Deterministic(t *testing.T) {
	mkCollector := func() *findingCollector {
		c := &findingCollector{}
		c.add(Finding{
			"kind":                "synthetic_event",
			"id":                  "f-1",
			"host":                "h",
			"metric_family_count": 3,
			"members": []interface{}{
				map[string]interface{}{"kind": KindSpike},
				map[string]interface{}{"kind": KindChangepoint},
			},
		})
		c.add(Finding{
			"kind":           "workload_spec",
			"id":             "f-2",
			"host":           "h",
			"query_rate_ps":  1500.0,
			"insert_rate_ps": 200.0,
		})
		c.add(Finding{
			"kind":     "missing_signal",
			"id":       "f-3",
			"host":     "h",
			"category": "repl", // Tier C: was "class"
		})
		return c
	}
	c1 := mkCollector()
	emitRecurrenceFingerprint(c1, "h")
	c2 := mkCollector()
	emitRecurrenceFingerprint(c2, "h")
	var fp1, fp2 string
	for _, f := range c1.items {
		if k, _ := f["kind"].(string); k == "recurrence_fingerprint" {
			fp1, _ = f["fingerprint"].(string)
		}
	}
	for _, f := range c2.items {
		if k, _ := f["kind"].(string); k == "recurrence_fingerprint" {
			fp2, _ = f["fingerprint"].(string)
		}
	}
	if fp1 == "" || fp1 != fp2 {
		t.Errorf("fingerprints not deterministic: %q vs %q", fp1, fp2)
	}
}

// TestEmitRecurrenceFingerprint_DiffersOnDifferentInput: a different
// workload produces a different fingerprint.
func TestEmitRecurrenceFingerprint_DiffersOnDifferentInput(t *testing.T) {
	c1 := &findingCollector{}
	c1.add(Finding{"kind": "workload_spec", "host": "h", "query_rate_ps": 1500.0})
	emitRecurrenceFingerprint(c1, "h")

	c2 := &findingCollector{}
	c2.add(Finding{"kind": "workload_spec", "host": "h", "insert_rate_ps": 1500.0})
	emitRecurrenceFingerprint(c2, "h")

	var fp1, fp2 string
	for _, f := range c1.items {
		if k, _ := f["kind"].(string); k == "recurrence_fingerprint" {
			fp1, _ = f["fingerprint"].(string)
		}
	}
	for _, f := range c2.items {
		if k, _ := f["kind"].(string); k == "recurrence_fingerprint" {
			fp2, _ = f["fingerprint"].(string)
		}
	}
	if fp1 == fp2 {
		t.Errorf("expected different fingerprints for different dominant ops, both = %q", fp1)
	}
}

// TestEmitRecurrenceFingerprint_NoShapeBits: collector with no relevant
// Findings emits no fingerprint.
func TestEmitRecurrenceFingerprint_NoShapeBits(t *testing.T) {
	c := &findingCollector{}
	c.add(Finding{"kind": KindMover, "host": "h"})
	emitted := emitRecurrenceFingerprint(c, "h")
	if emitted != 0 {
		t.Errorf("expected 0 fingerprint with no shape bits, got %d", emitted)
	}
}
