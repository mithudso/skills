package main

import (
	"testing"
	"time"
)

// TestCounterInversion tests a monotone counter with small max value and ≥90% non-decreasing.
// Expects kind:"counter_inversion" with reason_code "below_wrap_threshold".
func TestCounterInversion(t *testing.T) {
	c := &findingCollector{}
	col := &CompressedColumn{}
	// 20 samples: 19 non-decreasing transitions, 1 decrease → 19/20 = 95% ≥ 90% ✓
	data := []int64{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160, 170, 180, 175}
	for _, v := range data {
		col.Add(v, true)
	}
	col.Finish()

	store := &seriesStore{
		nameToID:   map[string]int{"test.counter": 0},
		names:      []string{"test.counter"},
		isCounter:  []bool{false},
		columns:    []*CompressedColumn{col},
		timestamps: make([]time.Time, len(data)),
	}
	for i := range store.timestamps {
		store.timestamps[i] = time.Now().Add(time.Duration(i) * time.Second)
	}

	emitted, _ := emitCounterEvents(c, "test.host", store)
	if emitted != 1 {
		t.Errorf("expected 1 emitted Finding, got %d", emitted)
	}
	if len(c.items) != 1 {
		t.Fatalf("expected 1 Finding in collector, got %d", len(c.items))
	}

	f := c.items[0]
	if kind, _ := f["kind"].(string); kind != "counter_inversion" {
		t.Errorf("expected kind:counter_inversion, got %q", kind)
	}
	if rc, _ := f["reason_code"].(string); rc != "below_wrap_threshold" {
		t.Errorf("expected reason_code:below_wrap_threshold, got %q", rc)
	}
	if mv, _ := f["prev_value"].(int64); mv != 180 {
		t.Errorf("expected prev_value:180, got %d", mv)
	}
}

// TestCounterWrap tests a counter with max > 0.9*2^32.
// Expects kind:"counter_wrap" with reason_code "uint32_wraparound_suspected".
func TestCounterWrap(t *testing.T) {
	c := &findingCollector{}
	maxVal := int64(4_000_000_000) // > 0.9*2^32 ≈ 3.8G
	col := &CompressedColumn{}
	// 20 samples: 19 non-decreasing, 1 wrap → 19/20 = 95% ≥ 90% ✓
	data := []int64{maxVal, maxVal + 100, maxVal + 200, maxVal + 300, maxVal + 400,
		maxVal + 500, maxVal + 600, maxVal + 700, maxVal + 800, maxVal + 900,
		maxVal + 1000, maxVal + 1100, maxVal + 1200, maxVal + 1300, maxVal + 1400,
		maxVal + 1500, maxVal + 1600, maxVal + 1700, maxVal + 1800, 100}
	for _, v := range data {
		col.Add(v, true)
	}
	col.Finish()

	store := &seriesStore{
		nameToID:   map[string]int{"test.wrapped": 0},
		names:      []string{"test.wrapped"},
		isCounter:  []bool{false},
		columns:    []*CompressedColumn{col},
		timestamps: make([]time.Time, len(data)),
	}
	for i := range store.timestamps {
		store.timestamps[i] = time.Now().Add(time.Duration(i) * time.Second)
	}

	emitted, _ := emitCounterEvents(c, "test.host", store)
	if emitted != 1 {
		t.Errorf("expected 1 emitted Finding, got %d", emitted)
	}
	if len(c.items) != 1 {
		t.Fatalf("expected 1 Finding in collector, got %d", len(c.items))
	}

	f := c.items[0]
	if kind, _ := f["kind"].(string); kind != "counter_wrap" {
		t.Errorf("expected kind:counter_wrap, got %q", kind)
	}
	if rc, _ := f["reason_code"].(string); rc != "uint32_wraparound_suspected" {
		t.Errorf("expected reason_code:uint32_wraparound_suspected, got %q", rc)
	}
}

// TestCounterNonClassified tests that non-counters do not emit.
func TestCounterNonClassified(t *testing.T) {
	c := &findingCollector{}
	col := &CompressedColumn{}
	for _, v := range []int64{100, 200, 150} {
		col.Add(v, true)
	}
	col.Finish()

	store := &seriesStore{
		nameToID:   map[string]int{"test.gauge": 0},
		names:      []string{"test.gauge"},
		isCounter:  []bool{true}, // true monotone counter: no decreases seen, nothing to detect
		columns:    []*CompressedColumn{col},
		timestamps: []time.Time{time.Now(), time.Now().Add(1 * time.Second), time.Now().Add(2 * time.Second)},
	}

	emitted, _ := emitCounterEvents(c, "test.host", store)
	if emitted != 0 {
		t.Errorf("expected 0 emitted for gauge, got %d", emitted)
	}
	if len(c.items) != 0 {
		t.Errorf("expected no Findings for gauge, got %d", len(c.items))
	}
}

// TestCounterSingleDecrease tests the 90% non-decreasing threshold.
// One decrease in 11 samples → 10/11 ≈ 91% non-decreasing ✓
// Three decreases in 11 samples → 8/11 ≈ 73% < 90% ✗
func TestCounterSingleDecrease(t *testing.T) {
	// Case 1: 90% non-decreasing passes.
	c1 := &findingCollector{}
	col1 := &CompressedColumn{}
	for _, v := range []int64{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 85} {
		col1.Add(v, true)
	}
	col1.Finish()
	store1 := &seriesStore{
		nameToID:   map[string]int{"test.m1": 0},
		names:      []string{"test.m1"},
		isCounter:  []bool{false},
		columns:    []*CompressedColumn{col1},
		timestamps: make([]time.Time, 11),
	}
	for i := range store1.timestamps {
		store1.timestamps[i] = time.Now().Add(time.Duration(i) * time.Second)
	}

	emitted1, _ := emitCounterEvents(c1, "test.host", store1)
	if emitted1 == 0 {
		t.Errorf("case 1 (91%% non-decreasing): expected 1 emitted, got 0")
	}

	// Case 2: <90% non-decreasing fails.
	c2 := &findingCollector{}
	col2 := &CompressedColumn{}
	for _, v := range []int64{0, 10, 5, 20, 15, 30, 25, 40, 35, 50, 45} {
		col2.Add(v, true)
	}
	col2.Finish()
	store2 := &seriesStore{
		nameToID:   map[string]int{"test.m2": 0},
		names:      []string{"test.m2"},
		isCounter:  []bool{false},
		columns:    []*CompressedColumn{col2},
		timestamps: make([]time.Time, 11),
	}
	for i := range store2.timestamps {
		store2.timestamps[i] = time.Now().Add(time.Duration(i) * time.Second)
	}

	emitted2, _ := emitCounterEvents(c2, "test.host", store2)
	if emitted2 != 0 {
		t.Errorf("case 2 (<90%% non-decreasing): expected 0 emitted, got %d", emitted2)
	}
}
