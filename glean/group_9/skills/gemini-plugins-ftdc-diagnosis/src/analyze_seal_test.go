package main

import (
	"fmt"
	"testing"
	"time"
)

// buildSealStore returns a minimal PALI seriesStore sufficient for
// isPALICapture to return true and for the timestamp index to be usable
// by emitSealEvents. The store covers [base, base+N seconds).
func buildSealStore(base time.Time, n int) *seriesStore {
	const anchor = "serverStatus.metrics.disagg.logServer.appendLogSuccess"
	s := &seriesStore{nameToID: map[string]int{anchor: 0}}
	s.names = []string{anchor}
	s.isCounter = []bool{false}
	s.isDegenerate = []bool{false}
	s.observedMin = []int64{0}
	s.observedMax = []int64{0}
	s.fullMean = []float64{0}
	s.fullSD = []float64{0}
	s.materializedSeries = make([][]float64, 1)
	col := NewCompressedColumn(n)
	for i := 0; i < n; i++ {
		col.Add(int64(100+i), true)
	}
	col.Finish()
	s.columns = []*CompressedColumn{col}
	for i := 0; i < n; i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	return s
}

// atStr formats a time.Time in the standard Finding "at" format.
func atStr(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}

// makeSealTermChange returns a synthetic term_change Finding at time t.
func makeSealTermChange(t time.Time, host string) Finding {
	return Finding{
		"kind": "term_change",
		"host": host,
		"at":   atStr(t),
		"id":   fmt.Sprintf("tc-%d", t.Unix()),
	}
}

// makeSealLogSpike returns a synthetic local_spike Finding on a logServer
// metric at time t (satisfies cond2).
func makeSealLogSpike(t time.Time, host string) Finding {
	return Finding{
		"kind":   "local_spike",
		"host":   host,
		"at":     atStr(t),
		"metric": "serverStatus.metrics.disagg.logServer.appendLogSuccess",
		"id":     fmt.Sprintf("spike-%d", t.Unix()),
	}
}

// makeSealPhylogGap returns a synthetic gap Finding on a phylog metric at
// time t (satisfies cond3a).
func makeSealPhylogGap(t time.Time, host string) Finding {
	return Finding{
		"kind":   "gap",
		"host":   host,
		"at":     atStr(t),
		"metric": "serverStatus.metrics.disagg.phylog.committedToMaterializedLagMillis",
		"id":     fmt.Sprintf("gap-%d", t.Unix()),
	}
}

// makeSealDisaggInversion returns a synthetic counter_inversion on a disagg
// metric at time t (blocks cond4 when within 30s before the term_change).
func makeSealDisaggInversion(t time.Time, host string) Finding {
	return Finding{
		"kind":   "counter_inversion",
		"host":   host,
		"at":     atStr(t),
		"metric": "serverStatus.metrics.disagg.logServer.appendLogTotal",
		"id":     fmt.Sprintf("inv-%d", t.Unix()),
	}
}

// TestSealConfidence_AllFourConditions: 4/4 conditions → confidence 0.85.
func TestSealConfidence_AllFourConditions(t *testing.T) {
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:01:00Z")
	store := buildSealStore(base.Add(-60*time.Second), 180)
	host := "h1"

	tcTime := base
	c := &findingCollector{}
	c.add(makeSealTermChange(tcTime, host))
	c.add(makeSealLogSpike(tcTime, host))         // cond2
	c.add(makeSealPhylogGap(tcTime, host))         // cond3
	// cond4: no counter_inversion → satisfied

	emitSealEvents(c, host, store)

	var got *Finding
	for i := range c.items {
		if k, _ := c.items[i]["kind"].(string); k == "seal_event" {
			f := c.items[i]
			got = &f
			break
		}
	}
	if got == nil {
		t.Fatal("expected a seal_event Finding, got none")
	}
	if conf, _ := (*got)["confidence_rank"].(float64); conf != 0.85 {
		t.Errorf("4/4 conditions: confidence_rank = %v, want 0.85", conf)
	}
	if n, _ := (*got)["conditions_met"].(int); n != 4 {
		t.Errorf("conditions_met = %d, want 4", n)
	}
}

// TestSealConfidence_MissingCond4: cond4 absent (counter_inversion present)
// → confidence 0.60, kind still "seal_event".
func TestSealConfidence_MissingCond4(t *testing.T) {
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:01:00Z")
	store := buildSealStore(base.Add(-60*time.Second), 180)
	host := "h1"

	tcTime := base
	c := &findingCollector{}
	c.add(makeSealTermChange(tcTime, host))
	c.add(makeSealLogSpike(tcTime, host))  // cond2
	c.add(makeSealPhylogGap(tcTime, host)) // cond3
	// Block cond4: inject a disagg counter_inversion 10s before the term_change.
	c.add(makeSealDisaggInversion(tcTime.Add(-10*time.Second), host))

	emitSealEvents(c, host, store)

	var got *Finding
	for i := range c.items {
		if k, _ := c.items[i]["kind"].(string); k == "seal_event" {
			f := c.items[i]
			got = &f
			break
		}
	}
	if got == nil {
		t.Fatal("expected a seal_event Finding, got none")
	}
	if conf, _ := (*got)["confidence_rank"].(float64); conf != 0.60 {
		t.Errorf("cond4 missing: confidence_rank = %v, want 0.60", conf)
	}
	if n, _ := (*got)["conditions_met"].(int); n != 3 {
		t.Errorf("conditions_met = %d, want 3", n)
	}
}

// TestSealConfidence_MissingCond2: cond2 absent (no logServer spike)
// → confidence 0.45, kind still "seal_event".
func TestSealConfidence_MissingCond2(t *testing.T) {
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:01:00Z")
	store := buildSealStore(base.Add(-60*time.Second), 180)
	host := "h1"

	tcTime := base
	c := &findingCollector{}
	c.add(makeSealTermChange(tcTime, host))
	// cond2: NO log spike
	c.add(makeSealPhylogGap(tcTime, host)) // cond3
	// cond4: no inversion → satisfied

	emitSealEvents(c, host, store)

	var got *Finding
	for i := range c.items {
		if k, _ := c.items[i]["kind"].(string); k == "seal_event" {
			f := c.items[i]
			got = &f
			break
		}
	}
	if got == nil {
		t.Fatal("expected a seal_event Finding, got none")
	}
	if conf, _ := (*got)["confidence_rank"].(float64); conf != 0.45 {
		t.Errorf("cond2 missing: confidence_rank = %v, want 0.45", conf)
	}
	if n, _ := (*got)["conditions_met"].(int); n != 3 {
		t.Errorf("conditions_met = %d, want 3", n)
	}
}

// TestSealConfidence_MissingCond3: cond3 absent (no phylog gap, no
// discontinuous checkpoint delta) → confidence 0.55.
func TestSealConfidence_MissingCond3(t *testing.T) {
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:01:00Z")
	store := buildSealStore(base.Add(-60*time.Second), 180)
	host := "h1"

	tcTime := base
	c := &findingCollector{}
	c.add(makeSealTermChange(tcTime, host))
	c.add(makeSealLogSpike(tcTime, host)) // cond2
	// cond3: NO phylog gap and no discontinuousCheckpointInstallCount delta (store has no such column)
	// cond4: no inversion → satisfied

	emitSealEvents(c, host, store)

	var got *Finding
	for i := range c.items {
		if k, _ := c.items[i]["kind"].(string); k == "seal_event" {
			f := c.items[i]
			got = &f
			break
		}
	}
	if got == nil {
		t.Fatal("expected a seal_event Finding, got none")
	}
	if conf, _ := (*got)["confidence_rank"].(float64); conf != 0.55 {
		t.Errorf("cond3 missing: confidence_rank = %v, want 0.55", conf)
	}
	if n, _ := (*got)["conditions_met"].(int); n != 3 {
		t.Errorf("conditions_met = %d, want 3", n)
	}
}
