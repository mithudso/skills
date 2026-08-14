package main

import (
	"math"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

// TestAbsFloat_NaN guards the NaN-safety guard: math.Abs propagates
// NaN, while the prior hand-rolled `if x < 0 { return -x }` returned NaN
// unchanged but silently made downstream comparisons (`absFloat(x) >= 100`)
// return false. The fix is correctness-critical — a future "optimization"
// reverting to the hand-rolled form would silently regress.
func TestAbsFloat_NaN(t *testing.T) {
	if !math.IsNaN(absFloat(math.NaN())) {
		t.Errorf("absFloat(NaN) must propagate NaN; got %v", absFloat(math.NaN()))
	}
	// Comparison check (the actual hazard): NaN compared against any
	// threshold returns false. Confirm via the function's typical call site.
	if absFloat(math.NaN()) >= 100 {
		t.Errorf("absFloat(NaN) >= 100 must be false (NaN comparison)")
	}
	// Sanity: positive infinity is still finite-comparable.
	if !math.IsInf(absFloat(math.Inf(1)), 1) {
		t.Errorf("absFloat(+Inf) must be +Inf")
	}
	if !math.IsInf(absFloat(math.Inf(-1)), 1) {
		t.Errorf("absFloat(-Inf) must be +Inf")
	}
}

// TestInferTopology pins the four-state inference for the metric_context
// Finding. Each case validates the classifier against a minimal name set.
func TestInferTopology(t *testing.T) {
	cases := []struct {
		name  string
		names []string
		want  string
	}{
		{"mongos", []string{
			"serverStatus.connections.current",
			"serverStatus.network.bytesIn",
		}, "sharded_router"},
		{"shard_primary", []string{
			"serverStatus.wiredTiger.cache.bytes currently in the cache",
			"serverStatus.repl.electionId",
			"serverStatus.shardingStatistics.countDonorMoveChunkStarted",
		}, "sharded_member"},
		{"replica_set", []string{
			"serverStatus.wiredTiger.cache.bytes currently in the cache",
			"serverStatus.repl.electionId",
		}, "replica_set"},
		{"standalone", []string{
			"serverStatus.wiredTiger.cache.bytes currently in the cache",
		}, "standalone"},
		{"empty", []string{}, ""},
	}
	for _, c := range cases {
		got := inferTopology(c.names)
		if got != c.want {
			t.Errorf("inferTopology(%s): got %q, want %q", c.name, got, c.want)
		}
	}
}

// TestFindingMagnitude verifies the rank-order key extracts the right
// field when present and falls back cleanly otherwise.
func TestFindingMagnitude(t *testing.T) {
	cases := []struct {
		name string
		f    Finding
		want float64
	}{
		{"effect_size wins", Finding{"effect_size": 3.5, "score": 100.0}, 3.5},
		{"score next", Finding{"score": 7.2}, 7.2},
		{"r next", Finding{"r": -0.8}, 0.8},
		{"welch_t", Finding{"welch_t": 12.0}, 12.0},
		{"slope_t", Finding{"slope_t": 5.5}, 5.5},
		{"no key", Finding{"metric": "foo"}, 0},
		{"non-float64 ignored", Finding{"effect_size": 3}, 0}, // int(3) not float64
	}
	for _, c := range cases {
		got := findingMagnitude(c.f)
		if got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

// TestApplyOutputGates_ScoreTableFilter: by default the
// metric_score_table Finding is removed; with EmitScoreTable=true it stays.
func TestApplyOutputGates_ScoreTableFilter(t *testing.T) {
	mkCollector := func() *findingCollector {
		return &findingCollector{items: []Finding{
			{"kind": KindRunMetadata, "host": "h"},
			{"kind": KindMetricScoreTable, "host": "h", "entries": []interface{}{}},
			{"kind": KindMover, "host": "h", "score": 1.0},
		}}
	}

	// Default: score table filtered, budget gate inactive.
	c1 := mkCollector()
	applyOutputGates(c1, AutoConfig{BudgetTokens: 0, EmitScoreTable: false})
	for _, f := range c1.items {
		if k, _ := f["kind"].(string); k == KindMetricScoreTable {
			t.Errorf("score table not filtered out with EmitScoreTable=false")
		}
	}
	if len(c1.items) != 2 {
		t.Errorf("expected 2 items after filter, got %d", len(c1.items))
	}

	// With EmitScoreTable: score table retained.
	c2 := mkCollector()
	applyOutputGates(c2, AutoConfig{BudgetTokens: 0, EmitScoreTable: true})
	hasScoreTable := false
	for _, f := range c2.items {
		if k, _ := f["kind"].(string); k == KindMetricScoreTable {
			hasScoreTable = true
		}
	}
	if !hasScoreTable {
		t.Errorf("score table should be retained with EmitScoreTable=true")
	}
}

// TestApplyOutputGates_BudgetGate: when output exceeds the budget,
// non-always-keep kinds are capped at top-3 and a `truncated` Finding
// summarizes what was elided.
func TestApplyOutputGates_BudgetGate(t *testing.T) {
	// Build a collector with 10 mover Findings + the always-keep kinds.
	c := &findingCollector{}
	c.add(Finding{"kind": KindRunMetadata, "host": "h"})
	c.add(Finding{"kind": KindRunHealth, "host": "h"})
	c.add(Finding{"kind": KindSchemaFingerprint, "host": "h"})
	for i := 0; i < 10; i++ {
		c.add(Finding{
			"kind":   KindMover,
			"host":   "h",
			"metric": "m" + intToString(i),
			"score":  float64(i),
		})
	}
	// Very small budget — guaranteed to trigger.
	applyOutputGates(c, AutoConfig{BudgetTokens: 100, EmitScoreTable: false})

	// Check: kept includes all 3 always-keep + top-3 movers + 1 truncated.
	moverCount := 0
	truncCount := 0
	var truncFinding Finding
	for _, f := range c.items {
		switch k, _ := f["kind"].(string); k {
		case KindMover:
			moverCount++
		case "truncated":
			truncCount++
			truncFinding = f
		}
	}
	if moverCount != 3 {
		t.Errorf("expected 3 mover Findings after budget gate, got %d", moverCount)
	}
	if truncCount != 1 {
		t.Errorf("expected 1 truncated Finding, got %d", truncCount)
	}
	if truncFinding != nil {
		elided, _ := truncFinding["kinds_elided"].(map[string]int)
		if elided["mover"] != 7 {
			t.Errorf("truncated.kinds_elided[mover]: got %d, want 7", elided["mover"])
		}
	}
}

// TestApplyOutputGates_BudgetGate_TopByMagnitude: when culling, the
// HIGHEST-magnitude entries should be kept (not just the first 3).
func TestApplyOutputGates_BudgetGate_TopByMagnitude(t *testing.T) {
	c := &findingCollector{}
	c.add(Finding{"kind": KindRunMetadata, "host": "h"})
	for i := 0; i < 10; i++ {
		c.add(Finding{
			"kind":   KindMover,
			"host":   "h",
			"metric": "m" + intToString(i),
			"score":  float64(i), // index 9 has highest score
		})
	}
	applyOutputGates(c, AutoConfig{BudgetTokens: 50, EmitScoreTable: false})

	// Confirm the 3 retained movers are the top-3 by score (7, 8, 9).
	scores := []float64{}
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == KindMover {
			s, _ := f["score"].(float64)
			scores = append(scores, s)
		}
	}
	if len(scores) != 3 {
		t.Fatalf("expected 3 retained movers, got %d", len(scores))
	}
	// Should be 9, 8, 7 in some order.
	gotSet := map[float64]bool{}
	for _, s := range scores {
		gotSet[s] = true
	}
	for _, want := range []float64{7, 8, 9} {
		if !gotSet[want] {
			t.Errorf("expected mover with score=%v retained; got %v", want, scores)
		}
	}
}

// TestApplyOutputGates_EmptyKept: degenerate case — no Findings whose
// kind is in the always-keep set AND every other kind has ≤3 entries.
// `kept` may still be non-empty (≤3 of each), but `elided` is empty so
// no truncated Finding emits. The kept[0] guard must not panic on
// genuinely empty paths.
func TestApplyOutputGates_EmptyKept(t *testing.T) {
	c := &findingCollector{}
	c.add(Finding{"kind": KindMover, "host": "h", "score": 1.0})
	// Tiny budget but ≤3 of each kind — no elision needed.
	applyOutputGates(c, AutoConfig{BudgetTokens: 1, EmitScoreTable: false})
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "truncated" {
			t.Errorf("unexpected truncated Finding when no elision was needed")
		}
	}
}

// TestEmitMetricGlossary covers (a) empty input fast path, (b)
// no-overlap-with-catalog early return, (c) dedup + sorted entries.
func TestEmitMetricGlossary(t *testing.T) {
	// (a) Empty
	c1 := &findingCollector{}
	emitMetricGlossary(c1, "h")
	if len(c1.items) != 0 {
		t.Errorf("empty collector: expected no glossary Finding, got %d", len(c1.items))
	}

	// (b) No overlap with catalog
	c2 := &findingCollector{items: []Finding{
		{"kind": KindMover, "metric": "unknown.metric.path", "host": "h"},
	}}
	emitMetricGlossary(c2, "h")
	hasGlossary := false
	for _, f := range c2.items {
		if k, _ := f["kind"].(string); k == "metric_glossary" {
			hasGlossary = true
		}
	}
	if hasGlossary {
		t.Errorf("expected no glossary Finding when no cited metric is in catalog")
	}

	// (c) Cited metrics with catalog overlap
	c3 := &findingCollector{items: []Finding{
		{"kind": KindMover, "metric": "serverStatus.connections.current", "host": "h"},
		{"kind": KindCorrelation, "leader": "serverStatus.connections.current", "follower": "serverStatus.opcounters.query", "host": "h"},
	}}
	emitMetricGlossary(c3, "h")
	var glossary Finding
	for _, f := range c3.items {
		if k, _ := f["kind"].(string); k == "metric_glossary" {
			glossary = f
		}
	}
	if glossary == nil {
		t.Fatalf("expected glossary Finding to be emitted")
	}
	entries, _ := glossary["entries"].([]interface{})
	if len(entries) != 2 {
		t.Errorf("expected 2 dedup'd entries, got %d", len(entries))
	}
	// Confirm sorted by metric name.
	prev := ""
	for _, e := range entries {
		m, _ := e.(map[string]interface{})
		mn, _ := m["metric"].(string)
		if prev != "" && mn < prev {
			t.Errorf("entries not sorted: %q < %q", mn, prev)
		}
		prev = mn
	}
}

// TestWalkMetadataMap_NestedBSON exercises the bson.D / bson.A traversal
// the metric_context emitter relies on. Without the bson.A case, nested
// D inside arrays would be invisible.
func TestWalkMetadataMap_NestedBSON(t *testing.T) {
	doc := bson.D{
		{Key: "version", Value: "8.0.5"},
		{Key: "system", Value: bson.D{
			{Key: "numCores", Value: int32(16)},
			{Key: "memSizeMB", Value: int32(32768)},
		}},
		{Key: "versionArray", Value: bson.A{int32(8), int32(0), int32(5)}},
		{Key: "modules", Value: bson.A{
			bson.D{{Key: "name", Value: "enterprise"}},
		}},
	}
	got := map[string]interface{}{}
	walkMetadataMap(doc, func(k string, v interface{}) {
		got[k] = v
	})
	if got["version"] != "8.0.5" {
		t.Errorf("missing version: got %v", got["version"])
	}
	if v, _ := toInt(got["system.numCores"]); v != 16 {
		t.Errorf("missing system.numCores: got %v", got["system.numCores"])
	}
	if v, _ := toInt(got["system.memSizeMB"]); v != 32768 {
		t.Errorf("missing system.memSizeMB: got %v", got["system.memSizeMB"])
	}
	// versionArray entries reachable via index suffix.
	if v, _ := toInt(got["versionArray.0"]); v != 8 {
		t.Errorf("missing versionArray.0: got %v (keys present: %v)", got["versionArray.0"], got)
	}
	// Nested D inside A reachable.
	if got["modules.0.name"] != "enterprise" {
		t.Errorf("missing modules.0.name: got %v", got["modules.0.name"])
	}
}

// TestEmitMetricContext_EmptyMetadata: with no metadata docs, the
// Finding still emits with at-minimum host + metric_count (required by
// schema).
func TestEmitMetricContext_EmptyMetadata(t *testing.T) {
	store := &seriesStore{
		nameToID: map[string]int{
			"serverStatus.wiredTiger.cache.bytes currently in the cache": 0,
			"serverStatus.repl.electionId":                               1,
		},
		names: []string{
			"serverStatus.wiredTiger.cache.bytes currently in the cache",
			"serverStatus.repl.electionId",
		},
	}
	c := &findingCollector{}
	emitMetricContext(c, "h", store, nil)
	if len(c.items) != 1 {
		t.Fatalf("expected 1 metric_context Finding, got %d", len(c.items))
	}
	f := c.items[0]
	if k, _ := f["kind"].(string); k != "metric_context" {
		t.Errorf("expected kind:metric_context, got %v", f["kind"])
	}
	if f["host"] != "h" {
		t.Errorf("expected host=h, got %v", f["host"])
	}
	if v, _ := f["metric_count"].(int); v != 2 {
		t.Errorf("expected metric_count=2, got %v", f["metric_count"])
	}
	// Topology should be inferred from the metric set.
	if topo, _ := f["topology"].(string); topo != "replica_set" {
		t.Errorf("expected topology=replica_set, got %v", topo)
	}
}
