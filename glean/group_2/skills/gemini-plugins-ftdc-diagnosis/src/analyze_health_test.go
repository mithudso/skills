package main

import (
	"testing"
	"time"
)

// utilStoreFromVals builds a 2-metric seriesStore (metric + denom) with
// the supplied 1Hz timestamps and int64 value series. Used by health-stage
// utilization tests.
func utilStoreFromVals(t *testing.T, metric, denom string, mVals, dVals []int64) *seriesStore {
	t.Helper()
	if len(mVals) != len(dVals) {
		t.Fatalf("len mVals=%d != len dVals=%d", len(mVals), len(dVals))
	}
	s := &seriesStore{nameToID: map[string]int{}}
	s.names = []string{metric, denom}
	s.nameToID[metric] = 0
	s.nameToID[denom] = 1
	s.isCounter = []bool{false, false}
	s.isDegenerate = []bool{false, false}
	s.observedMin = []int64{0, 0}
	s.observedMax = []int64{0, 0}
	s.fullMean = []float64{0, 0}
	s.fullSD = []float64{0, 0}
	s.materializedSeries = make([][]float64, 2)
	mCol := NewCompressedColumn(len(mVals))
	dCol := NewCompressedColumn(len(dVals))
	for i := range mVals {
		mCol.Add(mVals[i], true)
		dCol.Add(dVals[i], true)
	}
	mCol.Finish()
	dCol.Finish()
	s.columns = []*CompressedColumn{mCol, dCol}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < len(mVals); i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	return s
}

// TestEmitUtilization_RatioMath: closed-form expected values on a
// known fraction_of_total series.
func TestEmitUtilization_RatioMath(t *testing.T) {
	// Pick a metric from the catalog so the rule fires.
	metric := "serverStatus.queues.execution.write.out"
	denom := "serverStatus.queues.execution.write.available"
	mVals := []int64{1, 2, 3, 4, 5}
	dVals := []int64{9, 8, 7, 6, 5}
	// Ratios = {0.1, 0.2, 0.3, 0.4, 0.5}
	store := utilStoreFromVals(t, metric, denom, mVals, dVals)
	c := &findingCollector{}
	emitUtilization(c, "h", store, 0, 5, 0, 5)
	var f Finding
	for _, x := range c.items {
		if k, _ := x["kind"].(string); k == "utilization" {
			f = x
		}
	}
	if f == nil {
		t.Fatal("expected utilization Finding")
	}
	if mean, _ := f["mean"].(float64); mean < 0.29 || mean > 0.31 {
		t.Errorf("mean: got %v, want ~0.3", mean)
	}
	if max, _ := f["max"].(float64); max < 0.49 || max > 0.51 {
		t.Errorf("max: got %v, want ~0.5", max)
	}
}

// TestEmitUtilization_StressDelta: baseline window at 0.2, stress
// window at 0.8 → stress_delta ≈ 0.6.
func TestEmitUtilization_StressDelta(t *testing.T) {
	metric := "serverStatus.queues.execution.write.out"
	denom := "serverStatus.queues.execution.write.available"
	mVals := make([]int64, 100)
	dVals := make([]int64, 100)
	for i := 0; i < 50; i++ {
		mVals[i] = 2
		dVals[i] = 8 // ratio 0.2
	}
	for i := 50; i < 100; i++ {
		mVals[i] = 8
		dVals[i] = 2 // ratio 0.8
	}
	store := utilStoreFromVals(t, metric, denom, mVals, dVals)
	c := &findingCollector{}
	emitUtilization(c, "h", store, 0, 50, 50, 100)
	var f Finding
	for _, x := range c.items {
		if k, _ := x["kind"].(string); k == "utilization" {
			f = x
		}
	}
	if f == nil {
		t.Fatal("expected utilization Finding")
	}
	if d, _ := f["stress_delta"].(float64); d < 0.55 || d > 0.65 {
		t.Errorf("stress_delta: got %v, want ~0.6", d)
	}
}

// TestEmitUtilization_FractionOfMax: cache uses denom as absolute
// ceiling, not (m+denom).
func TestEmitUtilization_FractionOfMax(t *testing.T) {
	metric := "serverStatus.wiredTiger.cache.bytes currently in the cache"
	denom := "serverStatus.wiredTiger.cache.maximum bytes configured"
	// m = 50, denom = 100 → ratio = 50/100 = 0.5 (not 50/150 = 0.33)
	mVals := []int64{50, 50, 50, 50, 50}
	dVals := []int64{100, 100, 100, 100, 100}
	store := utilStoreFromVals(t, metric, denom, mVals, dVals)
	c := &findingCollector{}
	emitUtilization(c, "h", store, 0, 5, 0, 5)
	var f Finding
	for _, x := range c.items {
		if k, _ := x["kind"].(string); k == "utilization" {
			f = x
		}
	}
	if f == nil {
		t.Fatal("expected utilization Finding")
	}
	if mean, _ := f["mean"].(float64); mean < 0.49 || mean > 0.51 {
		t.Errorf("fraction_of_max mean: got %v, want ~0.5 (not 0.33)", mean)
	}
}

// TestEmitUtilization_ZeroCeilingSkipped: when ceiling collapses to 0
// the row is skipped without divide-by-zero.
func TestEmitUtilization_ZeroCeilingSkipped(t *testing.T) {
	metric := "serverStatus.queues.execution.write.out"
	denom := "serverStatus.queues.execution.write.available"
	mVals := []int64{0, 0, 0}
	dVals := []int64{0, 0, 0}
	store := utilStoreFromVals(t, metric, denom, mVals, dVals)
	c := &findingCollector{}
	// emitUtilization should observe ceiling<=0 every iteration and
	// emit nothing (samples_evaluated == 0 → skipped).
	emitUtilization(c, "h", store, 0, 3, 0, 3)
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "utilization" {
			t.Errorf("expected no utilization Finding when ceiling is 0 everywhere")
		}
	}
}

// TestEmitAlertThreshold_Election: electionId changes mid-capture.
func TestEmitAlertThreshold_Election(t *testing.T) {
	s := &seriesStore{nameToID: map[string]int{}}
	s.names = []string{"serverStatus.repl.electionId"}
	s.nameToID["serverStatus.repl.electionId"] = 0
	s.isCounter = []bool{false}
	s.isDegenerate = []bool{false}
	s.observedMin = []int64{0}
	s.observedMax = []int64{0}
	s.fullMean = []float64{0}
	s.fullSD = []float64{0}
	s.materializedSeries = make([][]float64, 1)
	col := NewCompressedColumn(10)
	// Election at sample 5: 100 → 200.
	for i := 0; i < 10; i++ {
		v := int64(100)
		if i >= 5 {
			v = 200
		}
		col.Add(v, true)
	}
	col.Finish()
	s.columns = []*CompressedColumn{col}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < 10; i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	c := &findingCollector{}
	emitted, _ := emitAlertThreshold(c, "h", s)
	if emitted != 1 {
		t.Errorf("expected 1 alert_threshold Finding for election, got %d", emitted)
	}
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "alert_threshold" {
			if ac, _ := f["alert_class"].(string); ac != "election" {
				t.Errorf("alert_class: got %q, want election", ac)
			}
			if v, _ := f["new_value"].(int64); v != 200 {
				t.Errorf("new_value: got %v, want 200", v)
			}
		}
	}
}

// TestEmitAlertThreshold_Restart: uptime drops to < 60 mid-capture.
func TestEmitAlertThreshold_Restart(t *testing.T) {
	s := &seriesStore{nameToID: map[string]int{}}
	s.names = []string{"serverStatus.uptime"}
	s.nameToID["serverStatus.uptime"] = 0
	s.isCounter = []bool{false}
	s.isDegenerate = []bool{false}
	s.observedMin = []int64{0}
	s.observedMax = []int64{0}
	s.fullMean = []float64{0}
	s.fullSD = []float64{0}
	s.materializedSeries = make([][]float64, 1)
	col := NewCompressedColumn(20)
	// Uptime grows linearly then resets to 1.
	for i := 0; i < 10; i++ {
		col.Add(int64(1000+i), true)
	}
	for i := 0; i < 10; i++ {
		col.Add(int64(1+i), true)
	}
	col.Finish()
	s.columns = []*CompressedColumn{col}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < 20; i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	c := &findingCollector{}
	emitted, _ := emitAlertThreshold(c, "h", s)
	if emitted < 1 {
		t.Errorf("expected ≥1 alert_threshold Finding for restart, got %d", emitted)
	}
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "alert_threshold" {
			if ac, _ := f["alert_class"].(string); ac != "restart" {
				t.Errorf("alert_class: got %q, want restart", ac)
			}
		}
	}
}

// TestEmitAlertThreshold_NoMetric: store without electionId/uptime →
// 0 emitted, no panic.
func TestEmitAlertThreshold_NoMetric(t *testing.T) {
	s := &seriesStore{nameToID: map[string]int{}}
	c := &findingCollector{}
	emitted, _ := emitAlertThreshold(c, "h", s)
	if emitted != 0 {
		t.Errorf("expected 0 emits with no source metrics, got %d", emitted)
	}
}

// TestEmitClockSkew_FiresOnDrift: an isolated 100ms jump (clock
// adjustment / NTP step) drifts the timestamps away from the modal-
// interval reference. Monotonic per-sample drift doesn't trigger
// because it gets absorbed into the modal interval — that's by
// design; this detector targets irregular wall-clock events.
func TestEmitClockSkew_FiresOnDrift(t *testing.T) {
	s := &seriesStore{}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < 200; i++ {
		offset := time.Duration(i) * time.Second
		if i == 100 {
			// One-time 200ms jump at sample 100.
			offset += 200 * time.Millisecond
		}
		s.timestamps = append(s.timestamps, base.Add(offset))
	}
	c := &findingCollector{}
	emitted, _ := emitClockSkew(c, "h", s)
	if emitted < 1 {
		t.Errorf("expected ≥1 clock_skew Finding on 200ms jump, got %d", emitted)
	}
}

// TestEmitClockSkew_DoesNotFireOnSmallDrift: drift ≤ 50ms is tolerated.
func TestEmitClockSkew_DoesNotFireOnSmallDrift(t *testing.T) {
	s := &seriesStore{}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	// Drift 0.1ms per sample → 20ms drift at sample 200.
	for i := 0; i < 200; i++ {
		s.timestamps = append(s.timestamps, base.Add(
			time.Duration(i)*time.Second+time.Duration(i)*100*time.Microsecond,
		))
	}
	c := &findingCollector{}
	emitted, _ := emitClockSkew(c, "h", s)
	if emitted != 0 {
		t.Errorf("expected 0 clock_skew Findings on 20ms drift, got %d", emitted)
	}
}

// TestEmitMissingSignal_ClassDetection: a store missing only repl
// fires the repl class; missing nothing → no Findings.
func TestEmitMissingSignal_ClassDetection(t *testing.T) {
	// Build a store with all `core` + `wt` metrics, but NOT
	// `repl.electionId`.
	s := &seriesStore{nameToID: map[string]int{}}
	present := []string{
		"serverStatus.uptime",
		"serverStatus.opcounters.insert",
		"serverStatus.opcounters.query",
		"serverStatus.connections.current",
		"serverStatus.mem.resident",
		"serverStatus.wiredTiger.cache.bytes currently in the cache",
		"serverStatus.wiredTiger.cache.maximum bytes configured",
		"serverStatus.wiredTiger.cache.tracked dirty bytes in the cache",
	}
	for i, n := range present {
		s.names = append(s.names, n)
		s.nameToID[n] = i
	}
	c := &findingCollector{}
	emitted, _ := emitMissingSignal(c, "h", s)
	if emitted != 1 {
		t.Fatalf("expected exactly 1 missing_signal (repl.electionId), got %d", emitted)
	}
	var f Finding
	for _, x := range c.items {
		if k, _ := x["kind"].(string); k == KindMissingSignal {
			f = x
		}
	}
	if f == nil {
		t.Fatal("expected a missing_signal Finding")
	}
	if class, _ := f["category"].(string); class != "repl" {
		t.Errorf("category: got %q, want repl", class)
	}
}

// TestSortFloat64Slice: regression for the perf fix at the end of
// analyze_health.go. Asserts the slice is sorted in-place.
func TestSortFloat64Slice(t *testing.T) {
	in := []float64{3.5, 1.0, 2.7, 0.5, 4.2}
	want := []float64{0.5, 1.0, 2.7, 3.5, 4.2}
	sortFloat64Slice(in)
	for i := range want {
		if in[i] != want[i] {
			t.Errorf("idx %d: got %v, want %v", i, in[i], want[i])
		}
	}
}

// TestRatioStats_Empty: empty input must not panic or return NaN.
func TestRatioStats_Empty(t *testing.T) {
	mean, p50, p95, max := ratioStats(nil)
	if mean != 0 || p50 != 0 || p95 != 0 || max != 0 {
		t.Errorf("empty ratioStats: got %v %v %v %v, want all zero", mean, p50, p95, max)
	}
}

// TestRatioStats_P50_EvenOdd: verify p50 returns the correct median for
// n=1, n=2, n=3, and n=4 (the even-length case was the off-by-one site).
func TestRatioStats_P50_EvenOdd(t *testing.T) {
	cases := []struct {
		in      []float64
		wantP50 float64
		desc    string
	}{
		{[]float64{5.0}, 5.0, "n=1 single element"},
		// n=2: sorted=[1,3]; (n-1)*50/100 = 1*50/100 = 0 → index 0 = 1.0
		// (lower of the two, not the max — that's the fix)
		{[]float64{3.0, 1.0}, 1.0, "n=2 lower element, not max"},
		// n=3: sorted=[1,2,3]; (3-1)*50/100 = 1 → index 1 = 2.0
		{[]float64{3.0, 1.0, 2.0}, 2.0, "n=3 true median"},
		// n=4: sorted=[1,2,3,4]; (4-1)*50/100 = 1 → index 1 = 2.0
		{[]float64{4.0, 2.0, 1.0, 3.0}, 2.0, "n=4 lower-middle element"},
	}
	for _, tc := range cases {
		_, p50, _, _ := ratioStats(tc.in)
		if p50 != tc.wantP50 {
			t.Errorf("%s: p50 got %v, want %v", tc.desc, p50, tc.wantP50)
		}
	}
}

// TestRatioStats_P95: verify p95 uses (n-1)*95/100 (not n*95/100).
func TestRatioStats_P95(t *testing.T) {
	cases := []struct {
		in      []float64
		wantP95 float64
		desc    string
	}{
		// n=1: (1-1)*95/100 = 0 → index 0 = 7.0
		{[]float64{7.0}, 7.0, "n=1 single element"},
		// n=2: sorted=[1,3]; (2-1)*95/100 = 0 → index 0 = 1.0
		// (old formula: 2*95/100=1 → index 1 = 3.0, wrong)
		{[]float64{3.0, 1.0}, 1.0, "n=2 lower element"},
		// n=20: sorted 1..20; (20-1)*95/100 = 18 → index 18 = 19.0
		{func() []float64 {
			s := make([]float64, 20)
			for i := range s {
				s[i] = float64(i + 1)
			}
			return s
		}(), 19.0, "n=20 index 18"},
	}
	for _, tc := range cases {
		_, _, p95, _ := ratioStats(tc.in)
		if p95 != tc.wantP95 {
			t.Errorf("%s: p95 got %v, want %v", tc.desc, p95, tc.wantP95)
		}
	}
}

// TestCaptureQualityScore_CVFailurePctDeduction verifies that
// emitCaptureQualityScore deducts 1 point per percentage point of
// cv_failure_pct (computed from changepoint cross_validated flags), capped at 20.
func TestCaptureQualityScore_CVFailurePctDeduction(t *testing.T) {
	// 3 failed + 17 passed = 20 total changepoints → cv_failure_pct = 15%.
	c := &findingCollector{}
	for i := 0; i < 17; i++ {
		c.add(Finding{"kind": KindChangepoint, "cross_validated": true})
	}
	for i := 0; i < 3; i++ {
		c.add(Finding{"kind": KindChangepoint, "cross_validated": false})
	}
	emitCaptureQualityScore(c, "h")

	var qf Finding
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == KindCaptureQualityScore {
			qf = f
			break
		}
	}
	if qf == nil {
		t.Fatal("expected a capture_quality_score Finding")
	}
	// With only a cv_failure_pct of 15 and no other deductions, score = 100 - 15 = 85.
	score, _ := qf["score"].(float64)
	if score < 84.9 || score > 85.1 {
		t.Errorf("score: got %v, want 85.0 (100 - 15 cv_failure_pct deduction)", score)
	}
	deduction, _ := qf["deduction_cv_failure"].(int)
	if deduction != 15 {
		t.Errorf("deduction_cv_failure: got %v, want 15", deduction)
	}
}

// TestCaptureQualityScore_CVFailurePctCap verifies the cv_failure_pct
// deduction is capped at 20 even when cv_failure_pct is very high.
func TestCaptureQualityScore_CVFailurePctCap(t *testing.T) {
	// 75 failed + 25 passed = 100 total changepoints → cv_failure_pct = 75%.
	// Deduction should be capped at 20.
	c := &findingCollector{}
	for i := 0; i < 25; i++ {
		c.add(Finding{"kind": KindChangepoint, "cross_validated": true})
	}
	for i := 0; i < 75; i++ {
		c.add(Finding{"kind": KindChangepoint, "cross_validated": false})
	}
	emitCaptureQualityScore(c, "h")

	var qf Finding
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == KindCaptureQualityScore {
			qf = f
			break
		}
	}
	if qf == nil {
		t.Fatal("expected a capture_quality_score Finding")
	}
	deduction, _ := qf["deduction_cv_failure"].(int)
	if deduction != 20 {
		t.Errorf("deduction_cv_failure: got %v, want 20 (capped)", deduction)
	}
}

// buildMMcStore builds a seriesStore for the M/M/c counterfactual test.
// It creates the 4 required metric series: arrival counter, latency counter,
// ops counter, available gauge, and out gauge.
//
// Design values:
//   - arrivals: 10 ops/sample → meanArrivals = 10
//   - latency: 1000 µs per op → srvMeanMicros = 1000 → µ = 1000 ops/sec
//   - capacity: available=90, out=10 → total capacity = 100
//   - ρ = λ / (c * µ) = 10 / (100 * 1000) ≈ 0.0001 (very healthy)
func buildMMcStore(t *testing.T) *seriesStore {
	t.Helper()
	const n = 120 // ≥60 samples required
	metrics := []string{
		"serverStatus.opcounters.update",
		"serverStatus.opLatencies.writes.latency",
		"serverStatus.opLatencies.writes.ops",
		"serverStatus.wiredTiger.concurrentTransactions.write.available",
		"serverStatus.wiredTiger.concurrentTransactions.write.out",
	}
	s := &seriesStore{nameToID: make(map[string]int, len(metrics))}
	for i, name := range metrics {
		s.names = append(s.names, name)
		s.nameToID[name] = i
	}
	s.isCounter = make([]bool, len(metrics))
	s.isDegenerate = make([]bool, len(metrics))
	s.observedMin = make([]int64, len(metrics))
	s.observedMax = make([]int64, len(metrics))
	s.fullMean = make([]float64, len(metrics))
	s.fullSD = make([]float64, len(metrics))
	s.materializedSeries = make([][]float64, len(metrics))

	// Arrival counter: 10 ops/sample (cumulative).
	arrCol := NewCompressedColumn(n)
	for i := 0; i < n; i++ {
		arrCol.Add(int64(i*10), true)
	}
	arrCol.Finish()

	// Latency counter: 1000 µs/op × 10 ops/sample = 10000 µs cumulative per sample.
	latCol := NewCompressedColumn(n)
	for i := 0; i < n; i++ {
		latCol.Add(int64(i*10000), true)
	}
	latCol.Finish()

	// Ops counter for latency: same as arrival (10 ops/sample).
	opsCol := NewCompressedColumn(n)
	for i := 0; i < n; i++ {
		opsCol.Add(int64(i*10), true)
	}
	opsCol.Finish()

	// Available tickets gauge: 90 per sample (stable).
	availCol := NewCompressedColumn(n)
	for i := 0; i < n; i++ {
		availCol.Add(int64(90), true)
	}
	availCol.Finish()

	// Out (in-use) tickets gauge: 10 per sample (stable).
	outCol := NewCompressedColumn(n)
	for i := 0; i < n; i++ {
		outCol.Add(int64(10), true)
	}
	outCol.Finish()

	s.columns = []*CompressedColumn{arrCol, latCol, opsCol, availCol, outCol}
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	for i := 0; i < n; i++ {
		s.timestamps = append(s.timestamps, base.Add(time.Duration(i)*time.Second))
	}
	return s
}

// TestMMcCounterfactual_UtilizationSane verifies that the M/M/c model
// produces a reasonable (well-below 0.99 clamp) utilization for a healthy
// workload with a known service rate.
func TestMMcCounterfactual_UtilizationSane(t *testing.T) {
	store := buildMMcStore(t)
	// Build a collector that already has a tunable_observation for the
	// write-tickets tunable (emitTunableCounterfactuals reads these).
	c := &findingCollector{}
	c.add(Finding{
		"schema_version": "1",
		"kind":           "tunable_observation",
		"id":             "f-abcd1234",
		"host":           "h",
		"tunable":        "wiredTigerConcurrentWriteTransactions",
		"direction":      "increase",
	})

	cfg := AutoConfig{RecommendTunables: true}
	emitted, _ := emitTunableCounterfactuals(c, "h", store, cfg)
	if emitted == 0 {
		t.Fatal("expected at least one tunable_counterfactual or refused Finding")
	}

	var cf Finding
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k == "tunable_counterfactual" {
			cf = f
			break
		}
	}
	if cf == nil {
		// Check for refused — if it refused, print the reason so the test is informative.
		for _, f := range c.items {
			if k, _ := f["kind"].(string); k == "tunable_counterfactual_refused" {
				t.Fatalf("model refused: reason_code=%v reason=%v", f["reason_code"], f["reason"])
			}
		}
		t.Fatal("expected a tunable_counterfactual Finding")
	}

	rho, _ := cf["current_utilization"].(float64)
	// With 10 arrivals/sample, 1000µs service time, and 100 capacity,
	// ρ should be very small (≪ 0.99 clamp).
	if rho >= 0.99 {
		t.Errorf("utilization hit the 0.99 clamp; expected a sane value, got %v", rho)
	}
	if rho <= 0 {
		t.Errorf("utilization is zero or negative: %v", rho)
	}
	// Capacity should be 100 (available=90 + out=10).
	cap_, _ := cf["current_capacity_observed"].(float64)
	if cap_ < 99 || cap_ > 101 {
		t.Errorf("current_capacity_observed: got %v, want 100 (available+out)", cap_)
	}
}
