package main

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

// buildCapacityStore helper: builds a minimal seriesStore with two
// metrics (rule.metric and rule.companion) and a per-sample timestamp
// series. Values come from the supplied arrays. Used to exercise
// emitCapacityPressure / emitNearCeiling deterministically without
// touching real FTDC data.
func buildCapacityStore(t *testing.T, metric, companion string, mVals, cVals []int64) *seriesStore {
	t.Helper()
	if len(mVals) != len(cVals) {
		t.Fatalf("buildCapacityStore: mVals and cVals must be same length, got %d/%d", len(mVals), len(cVals))
	}
	s := &seriesStore{nameToID: map[string]int{}}
	s.names = []string{metric, companion}
	s.nameToID[metric] = 0
	s.nameToID[companion] = 1
	s.isCounter = []bool{false, false}
	s.isDegenerate = []bool{false, false}
	s.observedMin = []int64{0, 0}
	s.observedMax = []int64{0, 0}
	s.fullMean = []float64{0, 0}
	s.fullSD = []float64{0, 0}
	s.materializedSeries = make([][]float64, 2)
	mCol := NewCompressedColumn(len(mVals))
	for _, v := range mVals {
		mCol.Add(v, true)
	}
	mCol.Finish()
	cCol := NewCompressedColumn(len(cVals))
	for _, v := range cVals {
		cCol.Add(v, true)
	}
	cCol.Finish()
	s.columns = []*CompressedColumn{mCol, cCol}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < len(mVals); i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	return s
}

// findKind scans a collector for the first Finding of a given kind.
func findKind(c *findingCollector, kind string) Finding {
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == kind {
			return f
		}
	}
	return nil
}

// countKind counts how many Findings of a given kind are in the collector.
func countKind(c *findingCollector, kind string) int {
	n := 0
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == kind {
			n++
		}
	}
	return n
}

// pickRule returns the catalog rule for a given metric path. Lets tests
// pin behavior to the actual rules that ship.
func pickRule(metric string) capacityRule {
	for _, r := range capacityRules {
		if r.metric == metric {
			return r
		}
	}
	panic("rule not found: " + metric)
}

// Tests ----------------------------------------------------------------------

func TestCapacityPressure_LongestRunAtCeiling(t *testing.T) {
	rule := pickRule("serverStatus.queues.execution.write.out")
	const n = 600
	mVals := make([]int64, n)
	cVals := make([]int64, n)
	for i := 0; i < n; i++ {
		// Total = 128 tickets. Sit at 30 (out=30, available=98) for
		// most of the capture, then saturate for 30 consecutive
		// seconds in the middle.
		if i >= 200 && i < 230 {
			mVals[i] = 128
			cVals[i] = 0
		} else {
			mVals[i] = 30
			cVals[i] = 98
		}
	}
	store := buildCapacityStore(t, rule.metric, rule.companion, mVals, cVals)
	c := &findingCollector{}
	emitCapacityPressure(c, "h", store, defaultAutoConfig())
	if got := countKind(c, "capacity_pressure"); got != 1 {
		t.Fatalf("expected 1 capacity_pressure Finding, got %d", got)
	}
	f := findKind(c, "capacity_pressure")
	if v, _ := f["longest_run_seconds"].(float64); v < 29 || v > 31 {
		t.Errorf("longest_run_seconds: got %v, want ~30", v)
	}
}

func TestCapacityPressure_PctGate(t *testing.T) {
	// 1000 samples, 60 scattered at-ceiling samples (no single long
	// run > 5s). pct_at = 60/1000 = 0.06, above the 0.05 gate for
	// ticket queues.
	rule := pickRule("serverStatus.queues.execution.write.out")
	const n = 1000
	mVals := make([]int64, n)
	cVals := make([]int64, n)
	for i := 0; i < n; i++ {
		if i%17 == 0 && i < 17*60 {
			mVals[i] = 128
			cVals[i] = 0
		} else {
			mVals[i] = 20
			cVals[i] = 108
		}
	}
	store := buildCapacityStore(t, rule.metric, rule.companion, mVals, cVals)
	c := &findingCollector{}
	emitCapacityPressure(c, "h", store, defaultAutoConfig())
	if got := countKind(c, "capacity_pressure"); got != 1 {
		t.Fatalf("expected 1 capacity_pressure Finding via pct gate, got %d", got)
	}
}

func TestCapacityPressure_NoCompanionDoesNotFire(t *testing.T) {
	// Construct a store with ONLY the metric, no companion. Catalog
	// rule requires companion presence — should be silently skipped.
	rule := pickRule("serverStatus.queues.execution.write.out")
	s := &seriesStore{nameToID: map[string]int{rule.metric: 0}}
	s.names = []string{rule.metric}
	s.isCounter = []bool{false}
	s.isDegenerate = []bool{false}
	s.observedMin = []int64{0}
	s.observedMax = []int64{0}
	s.fullMean = []float64{0}
	s.fullSD = []float64{0}
	col := NewCompressedColumn(10)
	for i := 0; i < 10; i++ {
		col.Add(128, true)
	}
	col.Finish()
	s.columns = []*CompressedColumn{col}
	s.materializedSeries = make([][]float64, 1)
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < 10; i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	c := &findingCollector{}
	emitCapacityPressure(c, "h", s, defaultAutoConfig())
	if countKind(c, "capacity_pressure") != 0 {
		t.Errorf("expected no Finding when companion is absent")
	}
}

// TestCapacityPressure_NilCompanionColumnNoPanic exercises the D2 fix:
// the companion column ID is registered in nameToID but store.columns[cid]
// is nil. emitCapacityPressure must skip the rule without panicking.
func TestCapacityPressure_NilCompanionColumnNoPanic(t *testing.T) {
	rule := pickRule("serverStatus.queues.execution.write.out")
	// Register both metric and companion in nameToID, but only allocate
	// the metric column. Leave the companion slot nil.
	s := &seriesStore{nameToID: map[string]int{rule.metric: 0, rule.companion: 1}}
	s.names = []string{rule.metric, rule.companion}
	s.isCounter = []bool{false, false}
	s.isDegenerate = []bool{false, false}
	s.observedMin = []int64{0, 0}
	s.observedMax = []int64{0, 0}
	s.fullMean = []float64{0, 0}
	s.fullSD = []float64{0, 0}
	col := NewCompressedColumn(10)
	for i := 0; i < 10; i++ {
		col.Add(128, true) // metric at full capacity
	}
	col.Finish()
	// Two slots: metric column real, companion column nil.
	s.columns = []*CompressedColumn{col, nil}
	s.materializedSeries = make([][]float64, 2)
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < 10; i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	c := &findingCollector{}
	// Must not panic.
	emitCapacityPressure(c, "h", s, defaultAutoConfig())
	if countKind(c, "capacity_pressure") != 0 {
		t.Errorf("expected no capacity_pressure Finding when companion column is nil")
	}
}

func TestCapacityPressure_CeilingZeroSkipped(t *testing.T) {
	// All samples at m=0,c=0 → ceiling=0 → ratio undefined → skipped.
	rule := pickRule("serverStatus.queues.execution.write.out")
	store := buildCapacityStore(t, rule.metric, rule.companion,
		[]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	c := &findingCollector{}
	emitCapacityPressure(c, "h", store, defaultAutoConfig())
	if countKind(c, "capacity_pressure") != 0 {
		t.Errorf("expected no Finding when ceiling is 0")
	}
}

func TestNearCeiling_BoundaryAt85Pct(t *testing.T) {
	rule := pickRule("serverStatus.queues.execution.write.out")
	cases := []struct {
		name   string
		m, c   int64
		expect bool
	}{
		{"below threshold (0.84)", 84, 16, false},
		{"at threshold (0.85)", 85, 15, true},
		{"above threshold (0.86)", 86, 14, true},
	}
	for _, tc := range cases {
		store := buildCapacityStore(t, rule.metric, rule.companion,
			[]int64{tc.m}, []int64{tc.c})
		col := &findingCollector{}
		emitNearCeiling(col, "h", store, defaultAutoConfig())
		got := countKind(col, "near_ceiling")
		if tc.expect && got != 1 {
			t.Errorf("%s: expected 1 near_ceiling Finding, got %d", tc.name, got)
		}
		if !tc.expect && got != 0 {
			t.Errorf("%s: expected 0 near_ceiling Findings, got %d", tc.name, got)
		}
	}
}

func TestConsistencyWarning_Agreeing(t *testing.T) {
	c := &findingCollector{}
	c.add(Finding{"kind": KindMover, "metric": "m1", "direction": "up", "host": "h"})
	c.add(Finding{"kind": KindMonotonicTrend, "metric": "m1", "slope": 5.0, "host": "h"})
	emitConsistencyWarning(c, "h")
	if countKind(c, "consistency_warning") != 0 {
		t.Errorf("expected no warning when mover up + slope > 0 agree")
	}
}

func TestConsistencyWarning_Disagreeing(t *testing.T) {
	c := &findingCollector{}
	c.add(Finding{"kind": KindMover, "metric": "m1", "direction": "up", "host": "h"})
	c.add(Finding{"kind": KindMonotonicTrend, "metric": "m1", "slope": -3.0, "host": "h"})
	emitConsistencyWarning(c, "h")
	if got := countKind(c, "consistency_warning"); got != 1 {
		t.Errorf("expected 1 warning on direction/slope disagree, got %d", got)
	}
}

func TestConsistencyWarning_MissingTrend(t *testing.T) {
	c := &findingCollector{}
	c.add(Finding{"kind": KindMover, "metric": "m1", "direction": "up", "host": "h"})
	emitConsistencyWarning(c, "h")
	if got := countKind(c, "consistency_warning"); got != 0 {
		t.Errorf("expected no warning when monotonic_trend is absent, got %d", got)
	}
}

func TestSyntheticEvent_MinimumThresholds(t *testing.T) {
	// Build a store with 100 1Hz timestamps.
	store := &seriesStore{nameToID: map[string]int{}}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < 100; i++ {
		store.timestamps = append(store.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	// 2 metrics, 1 family — should NOT fire.
	c := &findingCollector{}
	atStr := base.UTC().Format("2006-01-02T15:04:05.000Z")
	c.add(Finding{"kind": KindSpike, "metric": "serverStatus.opcounters.insert", "at": atStr, "value": 10.0, "median": 1.0, "host": "h"})
	c.add(Finding{"kind": KindSpike, "metric": "serverStatus.opcounters.query", "at": atStr, "value": 10.0, "median": 1.0, "host": "h"})
	emitSyntheticEvent(c, "h", store)
	if got := countKind(c, "synthetic_event"); got != 0 {
		t.Errorf("expected no synthetic_event with only 1 family, got %d", got)
	}
	// Add 1 metric from a different family — should now fire.
	c.add(Finding{"kind": KindSpike, "metric": "serverStatus.wiredTiger.cache.X", "at": atStr, "value": 50.0, "median": 5.0, "host": "h"})
	emitSyntheticEvent(c, "h", store)
	if got := countKind(c, "synthetic_event"); got != 1 {
		t.Errorf("expected 1 synthetic_event with 3 metrics + 2 families, got %d", got)
	}
}

func TestSyntheticEvent_WindowRadius(t *testing.T) {
	// With modal interval = 5s, radius = max(4*5, 5) = 20s. Two
	// Findings 18s apart should cluster; 22s apart should not.
	store := &seriesStore{nameToID: map[string]int{}}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < 100; i++ {
		store.timestamps = append(store.timestamps, base.Add(time.Duration(i)*5*time.Second))
	}
	mk := func(metric string, secOffset int) Finding {
		ts := base.Add(time.Duration(secOffset) * time.Second).UTC().Format("2006-01-02T15:04:05.000Z")
		return Finding{
			"kind": KindSpike, "metric": metric, "at": ts,
			"value": 10.0, "median": 1.0, "host": "h",
		}
	}
	// 3 metrics from 3 different families, spaced 8s apart — within
	// the 20s radius, single cluster.
	c := &findingCollector{}
	c.add(mk("serverStatus.opcounters.insert", 0))
	c.add(mk("serverStatus.wiredTiger.cache.X", 8))
	c.add(mk("config.transactions.stats.Y", 16))
	emitSyntheticEvent(c, "h", store)
	if got := countKind(c, "synthetic_event"); got != 1 {
		t.Errorf("expected 1 synthetic_event for 3 Findings within 20s radius, got %d", got)
	}
}

func TestSyntheticEvent_DirectionCounts(t *testing.T) {
	// Verify that mixed direction (up + down) emits both counts.
	store := &seriesStore{nameToID: map[string]int{}}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < 100; i++ {
		store.timestamps = append(store.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	atStr := base.UTC().Format("2006-01-02T15:04:05.000Z")
	c := &findingCollector{}
	// 2 up-spike from family A, 1 down-spike from family B.
	c.add(Finding{"kind": KindSpike, "metric": "serverStatus.opcounters.insert", "at": atStr, "value": 50.0, "median": 5.0, "host": "h"})
	c.add(Finding{"kind": KindSpike, "metric": "serverStatus.opcounters.query", "at": atStr, "value": 50.0, "median": 5.0, "host": "h"})
	c.add(Finding{"kind": KindSpike, "metric": "serverStatus.wiredTiger.cache.X", "at": atStr, "value": 1.0, "median": 10.0, "host": "h"})
	emitSyntheticEvent(c, "h", store)
	f := findKind(c, "synthetic_event")
	if f == nil {
		t.Fatalf("expected a synthetic_event Finding")
	}
	if up, _ := f["up_count"].(int); up != 2 {
		t.Errorf("up_count: got %d, want 2", up)
	}
	if down, _ := f["down_count"].(int); down != 1 {
		t.Errorf("down_count: got %d, want 1", down)
	}
}

func TestMetricFamily_Cases(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", ""},
		{"foo", "foo"},
		{"foo.bar", "foo.bar"},
		{"foo.bar.baz", "foo.bar"},
		{"serverStatus.wiredTiger.cache.bytes", "serverStatus.wiredTiger"},
		{"local.oplog.rs.stats.X", "local.oplog"},
	}
	for _, c := range cases {
		got := metricFamily(c.in)
		if got != c.want {
			t.Errorf("metricFamily(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestInferFindingDirection_AllPaths(t *testing.T) {
	cases := []struct {
		name string
		f    Finding
		want string
	}{
		{"direction field wins", Finding{"direction": "down", "value": 5.0, "median": 1.0}, "down"},
		{"baseline_mean / stressed_mean up", Finding{"baseline_mean": 1.0, "stressed_mean": 5.0}, "up"},
		{"baseline_mean / stressed_mean down", Finding{"baseline_mean": 5.0, "stressed_mean": 1.0}, "down"},
		{"value > median = up", Finding{"value": 10.0, "median": 1.0}, "up"},
		{"value < median = down", Finding{"value": 1.0, "median": 10.0}, "down"},
		{"value == median = ''", Finding{"value": 5.0, "median": 5.0}, ""},
		{"none set", Finding{"metric": "foo"}, ""},
	}
	for _, c := range cases {
		got := inferFindingDirection(c.f)
		if got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}

func TestNumericField_AllTypes(t *testing.T) {
	cases := []struct {
		name string
		v    interface{}
		want float64
		ok   bool
	}{
		{"float64", float64(3.5), 3.5, true},
		{"float32", float32(2.5), 2.5, true},
		{"int64", int64(7), 7, true},
		{"uint64", uint64(8), 8, true},
		{"int32", int32(9), 9, true},
		{"uint32", uint32(10), 10, true},
		{"int", int(11), 11, true},
		{"string rejected", "5", 0, false},
		{"bool rejected", true, 0, false},
	}
	for _, c := range cases {
		f := Finding{"k": c.v}
		got, ok := numericField(f, "k")
		if ok != c.ok || got != c.want {
			t.Errorf("%s: got (%v, %v), want (%v, %v)", c.name, got, ok, c.want, c.ok)
		}
	}
	// Missing key
	if _, ok := numericField(Finding{}, "missing"); ok {
		t.Errorf("missing key should return ok=false")
	}
}

func TestSyntheticEvent_Determinism(t *testing.T) {
	// Run the same input through emitSyntheticEvent twice; assert the
	// member_finding_ids array is identical across both runs. The
	// tie-break sort on (timestamp, kind, metric, id) guarantees this
	// even when Findings come in in arbitrary order — we shuffle
	// inputs by ordering same-timestamp Findings in REVERSE one run
	// vs. natural order the other, and verify both runs produce the
	// same cluster.
	store := &seriesStore{nameToID: map[string]int{}}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < 100; i++ {
		store.timestamps = append(store.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	atStr := base.UTC().Format("2006-01-02T15:04:05.000Z")
	mk := func(metric, id string, v float64) Finding {
		return Finding{
			"kind": KindSpike, "metric": metric, "at": atStr,
			"value": v, "median": 1.0, "host": "h", "id": id,
		}
	}
	run := func(reverse bool) []string {
		c := &findingCollector{}
		seeds := []Finding{
			mk("serverStatus.opcounters.insert", "f-aaaa", 10),
			mk("serverStatus.opcounters.query", "f-bbbb", 12),
			mk("serverStatus.wiredTiger.cache.X", "f-cccc", 8),
			mk("serverStatus.wiredTiger.cache.Y", "f-dddd", 15),
		}
		if reverse {
			for i, j := 0, len(seeds)-1; i < j; i, j = i+1, j-1 {
				seeds[i], seeds[j] = seeds[j], seeds[i]
			}
		}
		for _, s := range seeds {
			c.add(s)
		}
		emitSyntheticEvent(c, "h", store)
		f := findKind(c, "synthetic_event")
		if f == nil {
			t.Fatal("expected synthetic_event")
		}
		ids, _ := f["member_finding_ids"].([]string)
		return ids
	}
	a := run(false)
	b := run(true)
	if !reflect.DeepEqual(a, b) {
		t.Errorf("synthetic_event member_finding_ids non-deterministic: %v vs %v", a, b)
	}
}

func TestTopByMagnitude_KBounds(t *testing.T) {
	refs := []sortableRef{
		{id: "a", mag: 1.0},
		{id: "b", mag: 5.0},
		{id: "c", mag: 3.0},
		{id: "d", mag: 8.0},
		{id: "e", mag: 2.0},
	}
	cases := []struct {
		name string
		k    int
		want []string
	}{
		{"k=0 empty", 0, []string{}},
		{"k=1 top", 1, []string{"d"}},
		{"k=3 top", 3, []string{"d", "b", "c"}},
		{"k>len returns all", 99, []string{"d", "b", "c", "e", "a"}},
	}
	for _, c := range cases {
		got := topByMagnitude(refs, c.k)
		ids := []string{}
		for _, r := range got {
			ids = append(ids, r.id)
		}
		// Build comparable strings.
		gj := strings.Join(ids, ",")
		wj := strings.Join(c.want, ",")
		if gj != wj {
			t.Errorf("%s: got %q, want %q", c.name, gj, wj)
		}
	}
}
