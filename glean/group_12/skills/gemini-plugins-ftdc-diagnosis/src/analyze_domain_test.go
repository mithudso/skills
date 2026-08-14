package main

import (
	"testing"
	"time"
)

// makeStoreWithSingle builds a 1-metric seriesStore with the named
// metric and supplied int64 column. Used by domain-stage tests.
func makeStoreWithSingle(name string, vals []int64) *seriesStore {
	s := &seriesStore{nameToID: map[string]int{}}
	s.names = []string{name}
	s.nameToID[name] = 0
	s.isCounter = []bool{false}
	s.isDegenerate = []bool{false}
	s.observedMin = []int64{0}
	s.observedMax = []int64{0}
	s.fullMean = []float64{0}
	s.fullSD = []float64{0}
	s.materializedSeries = make([][]float64, 1)
	col := NewCompressedColumn(len(vals))
	for _, v := range vals {
		col.Add(v, true)
	}
	col.Finish()
	s.columns = []*CompressedColumn{col}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < len(vals); i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	return s
}

// TestSchemaFingerprint_Gates: each domain detector returns true ONLY
// when its anchor metric is present.
func TestSchemaFingerprint_Gates(t *testing.T) {
	empty := &seriesStore{nameToID: map[string]int{}}
	if isReplicaSetCapture(empty) {
		t.Error("isReplicaSetCapture should be false on empty store")
	}
	if isPALICapture(empty) {
		t.Error("isPALICapture should be false on empty store")
	}
	if isShardedCapture(empty) {
		t.Error("isShardedCapture should be false on empty store")
	}
	if hasHostOSMetrics(empty) {
		t.Error("hasHostOSMetrics should be false on empty store")
	}

	repl := makeStoreWithSingle("serverStatus.repl.electionId", []int64{1})
	if !isReplicaSetCapture(repl) {
		t.Error("isReplicaSetCapture should be true with electionId present")
	}

	pali := makeStoreWithSingle("serverStatus.metrics.disagg.foo", []int64{1})
	if !isPALICapture(pali) {
		t.Error("isPALICapture should be true with metrics.disagg.* present")
	}

	sharded := makeStoreWithSingle("serverStatus.shardingStatistics.bar", []int64{1})
	if !isShardedCapture(sharded) {
		t.Error("isShardedCapture should be true with shardingStatistics.* present")
	}

	host := makeStoreWithSingle("systemMetrics.cpu.user_ms", []int64{1})
	if !hasHostOSMetrics(host) {
		t.Error("hasHostOSMetrics should be true with systemMetrics.* present")
	}
}

// TestEmitReplicationFindings_TermChange: term increments → term_change
// Findings (one per increment).
func TestEmitReplicationFindings_TermChange(t *testing.T) {
	// Build a 2-metric store with electionId (gate) + term.
	s := &seriesStore{nameToID: map[string]int{}}
	s.names = []string{"serverStatus.repl.electionId", "replSetGetStatus.term"}
	s.nameToID["serverStatus.repl.electionId"] = 0
	s.nameToID["replSetGetStatus.term"] = 1
	s.isCounter = []bool{false, false}
	s.isDegenerate = []bool{false, false}
	s.observedMin = []int64{0, 0}
	s.observedMax = []int64{0, 0}
	s.fullMean = []float64{0, 0}
	s.fullSD = []float64{0, 0}
	s.materializedSeries = make([][]float64, 2)
	// electionId constant (just gates the stage); term increments 5→6 at sample 3.
	idCol := NewCompressedColumn(10)
	tCol := NewCompressedColumn(10)
	for i := 0; i < 10; i++ {
		idCol.Add(100, true)
		v := int64(5)
		if i >= 3 {
			v = 6
		}
		tCol.Add(v, true)
	}
	idCol.Finish()
	tCol.Finish()
	s.columns = []*CompressedColumn{idCol, tCol}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < 10; i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	c := &findingCollector{}
	emitted, _ := emitReplicationFindings(c, "h", s)
	if emitted < 1 {
		t.Errorf("expected ≥1 term_change Finding, got %d emitted total", emitted)
	}
	gotTermChange := false
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "term_change" {
			gotTermChange = true
			if v, _ := f["new_term"].(int64); v != 6 {
				t.Errorf("new_term: got %v, want 6", v)
			}
		}
	}
	if !gotTermChange {
		t.Errorf("term_change kind not emitted")
	}
}

// TestEmitReplicationFindings_GatedOff: non-replica-set capture
// (no electionId metric) → 0 emits, skipped=1.
func TestEmitReplicationFindings_GatedOff(t *testing.T) {
	s := &seriesStore{nameToID: map[string]int{}}
	c := &findingCollector{}
	emitted, skipped := emitReplicationFindings(c, "h", s)
	if emitted != 0 {
		t.Errorf("expected 0 emits when gate is off, got %d", emitted)
	}
	if skipped != 1 {
		t.Errorf("expected skipped=1 (whole-stage), got %d", skipped)
	}
}

// TestEmitReplicationFindings_OplogWindowLow: timeDiffHours dropping
// below 24 fires oplog_window_low.
func TestEmitReplicationFindings_OplogWindowLow(t *testing.T) {
	s := &seriesStore{nameToID: map[string]int{}}
	s.names = []string{"serverStatus.repl.electionId", "serverStatus.oplog.timeDiffHours"}
	s.nameToID["serverStatus.repl.electionId"] = 0
	s.nameToID["serverStatus.oplog.timeDiffHours"] = 1
	s.isCounter = []bool{false, false}
	s.isDegenerate = []bool{false, false}
	s.observedMin = []int64{0, 0}
	s.observedMax = []int64{0, 0}
	s.fullMean = []float64{0, 0}
	s.fullSD = []float64{0, 0}
	s.materializedSeries = make([][]float64, 2)
	idCol := NewCompressedColumn(10)
	tCol := NewCompressedColumn(10)
	for i := 0; i < 10; i++ {
		idCol.Add(100, true)
		tCol.Add(int64(48-i*2), true) // 48, 46, 44, ... 30
	}
	idCol.Finish()
	tCol.Finish()
	s.columns = []*CompressedColumn{idCol, tCol}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < 10; i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	c := &findingCollector{}
	emitReplicationFindings(c, "h", s)
	gotOplog := false
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "oplog_window_low" {
			gotOplog = true
			// At sample 9, timeDiffHours = 48 - 18 = 30. Not < 24,
			// so we shouldn't fire. Reverse: emit at min, but min is 30.
			// Actually min in our series is 30. Should NOT fire.
			_ = f
		}
	}
	if gotOplog {
		t.Errorf("oplog_window_low fired when min=30 (above 24-hour threshold)")
	}

	// Now build a series that goes below 24 — should fire.
	c2 := &findingCollector{}
	s2 := &seriesStore{nameToID: map[string]int{
		"serverStatus.repl.electionId":     0,
		"serverStatus.oplog.timeDiffHours": 1,
	}}
	s2.names = []string{"serverStatus.repl.electionId", "serverStatus.oplog.timeDiffHours"}
	s2.isCounter = []bool{false, false}
	s2.isDegenerate = []bool{false, false}
	s2.observedMin = []int64{0, 0}
	s2.observedMax = []int64{0, 0}
	s2.fullMean = []float64{0, 0}
	s2.fullSD = []float64{0, 0}
	s2.materializedSeries = make([][]float64, 2)
	idCol2 := NewCompressedColumn(10)
	tCol2 := NewCompressedColumn(10)
	for i := 0; i < 10; i++ {
		idCol2.Add(100, true)
		tCol2.Add(int64(30-i*4), true) // 30, 26, 22, ...
	}
	idCol2.Finish()
	tCol2.Finish()
	s2.columns = []*CompressedColumn{idCol2, tCol2}
	for i := 0; i < 10; i++ {
		s2.timestamps = append(s2.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	emitReplicationFindings(c2, "h", s2)
	gotOplog2 := false
	for _, f := range c2.items {
		if k, _ := f["kind"].(string); k == "oplog_window_low" {
			gotOplog2 = true
			if m, _ := f["min_hours"].(int64); m > 24 {
				t.Errorf("min_hours: got %v, expected <= 24", m)
			}
		}
	}
	if !gotOplog2 {
		t.Errorf("expected oplog_window_low to fire when min < 24")
	}
}

// TestEmitPALIFindings_CheckpointAnomaly: discontinuousCheckpointInstallCount
// increasing → fires.
func TestEmitPALIFindings_CheckpointAnomaly(t *testing.T) {
	s := makeStoreWithSingle("serverStatus.metrics.disagg.logServer.discontinuousCheckpointInstallCount",
		[]int64{0, 0, 0, 1, 1, 1, 2, 2, 2, 3})
	c := &findingCollector{}
	emitted, _ := emitPALIFindings(c, "h", s)
	if emitted < 1 {
		t.Errorf("expected ≥1 PALI Finding, got %d", emitted)
	}
	gotAnomaly := false
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "checkpoint_install_anomaly" {
			gotAnomaly = true
			if delta, _ := f["count_increase"].(int64); delta != 3 {
				t.Errorf("count_increase: got %v, want 3", delta)
			}
		}
	}
	if !gotAnomaly {
		t.Errorf("checkpoint_install_anomaly not emitted")
	}
}

// TestEmitPALIFindings_IORoutingNegativeDelta: a counter reset (lDelta < 0)
// must NOT produce an io_routing Finding (D3 fix).
func TestEmitPALIFindings_IORoutingNegativeDelta(t *testing.T) {
	const (
		localMetric    = "serverStatus.metrics.disagg.pageServer.local.success"
		nonLocalMetric = "serverStatus.metrics.disagg.pageServer.nonLocal.success"
		// A disagg anchor metric is needed for isPALICapture to gate in.
		anchorMetric = "serverStatus.metrics.disagg.logServer.appendLogSuccess"
	)
	s := &seriesStore{nameToID: map[string]int{}}
	s.names = []string{anchorMetric, localMetric, nonLocalMetric}
	s.nameToID[anchorMetric] = 0
	s.nameToID[localMetric] = 1
	s.nameToID[nonLocalMetric] = 2
	s.isCounter = []bool{false, false, false}
	s.isDegenerate = []bool{false, false, false}
	s.observedMin = []int64{0, 0, 0}
	s.observedMax = []int64{0, 0, 0}
	s.fullMean = []float64{0, 0, 0}
	s.fullSD = []float64{0, 0, 0}
	s.materializedSeries = make([][]float64, 3)

	anchorCol := NewCompressedColumn(10)
	localCol := NewCompressedColumn(10)
	nonLocalCol := NewCompressedColumn(10)
	for i := 0; i < 10; i++ {
		anchorCol.Add(int64(100+i), true)
		// local counter resets: starts high, ends low → lDelta < 0.
		localCol.Add(int64(5000-i*100), true)
		// nonLocal counter is large positive.
		nonLocalCol.Add(int64(i*200), true)
	}
	anchorCol.Finish()
	localCol.Finish()
	nonLocalCol.Finish()
	s.columns = []*CompressedColumn{anchorCol, localCol, nonLocalCol}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < 10; i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}

	c := &findingCollector{}
	emitPALIFindings(c, "h", s)
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "io_routing" {
			t.Errorf("io_routing Finding must not fire when lDelta < 0 (counter reset); got: %v", f)
		}
	}
}

// TestEmitPALIFindings_GatedOff: no metrics.disagg.* anchor → 0 emits.
func TestEmitPALIFindings_GatedOff(t *testing.T) {
	s := &seriesStore{nameToID: map[string]int{}}
	c := &findingCollector{}
	emitted, skipped := emitPALIFindings(c, "h", s)
	if emitted != 0 || skipped != 1 {
		t.Errorf("expected (0, 1), got (%d, %d)", emitted, skipped)
	}
}

// TestEmitShardedFindings_ChunkMigration: countDonorMoveChunkStarted
// rising → chunk_migration_event fires.
func TestEmitShardedFindings_ChunkMigration(t *testing.T) {
	s := makeStoreWithSingle("serverStatus.shardingStatistics.countDonorMoveChunkStarted",
		[]int64{10, 10, 10, 12, 12, 15, 18, 18, 18, 20})
	c := &findingCollector{}
	emitted, _ := emitShardedFindings(c, "h", s)
	if emitted < 1 {
		t.Errorf("expected ≥1 sharded Finding, got %d", emitted)
	}
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "chunk_migration_event" {
			if d, _ := f["migrations_started"].(int64); d != 10 {
				t.Errorf("migrations_started: got %v, want 10", d)
			}
		}
	}
}

// TestEmitHostOSFindings_CPUSteal: steal > 5% of total → host_pathology
// fires with subkind cpu_steal.
func TestEmitHostOSFindings_CPUSteal(t *testing.T) {
	s := &seriesStore{nameToID: map[string]int{}}
	s.names = []string{
		"systemMetrics.cpu.user_ms",
		"systemMetrics.cpu.system_ms",
		"systemMetrics.cpu.steal_ms",
	}
	for i, n := range s.names {
		s.nameToID[n] = i
	}
	s.isCounter = []bool{false, false, false}
	s.isDegenerate = []bool{false, false, false}
	s.observedMin = []int64{0, 0, 0}
	s.observedMax = []int64{0, 0, 0}
	s.fullMean = []float64{0, 0, 0}
	s.fullSD = []float64{0, 0, 0}
	s.materializedSeries = make([][]float64, 3)
	uCol := NewCompressedColumn(5)
	sysCol := NewCompressedColumn(5)
	stCol := NewCompressedColumn(5)
	for i := 0; i < 5; i++ {
		// Each sample adds 100 user + 20 system + 30 steal (total 150 deltas, steal=30/150=20%)
		uCol.Add(int64(100*i), true)
		sysCol.Add(int64(20*i), true)
		stCol.Add(int64(30*i), true)
	}
	uCol.Finish()
	sysCol.Finish()
	stCol.Finish()
	s.columns = []*CompressedColumn{uCol, sysCol, stCol}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < 5; i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	c := &findingCollector{}
	emitted, _ := emitHostOSFindings(c, "h", s)
	if emitted < 1 {
		t.Errorf("expected ≥1 host_pathology Finding, got %d", emitted)
	}
	gotSteal := false
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "host_pathology" {
			gotSteal = true
			if sk, _ := f["subkind"].(string); sk != "cpu_steal" {
				t.Errorf("subkind: got %q, want cpu_steal", sk)
			}
		}
	}
	if !gotSteal {
		t.Errorf("host_pathology not emitted on 20%% steal")
	}
}
