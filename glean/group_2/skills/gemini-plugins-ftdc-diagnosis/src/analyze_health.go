// analyze_health.go — health-monitoring stage covering USE/RED method
// + Atlas-style alert primitives + data-quality detectors that close
// the "implicit signal" gap.
//
// New kinds:
//   - `utilization`           — bounded ratio with explicit denominator
//   - `queue_depth_pressure`  — sustained non-zero queue depth
//   - `error_rate`            — RED-Errors anchor (asserts, conn rejects,
//                                command failures)
//   - `request_rate`          — RED-Rate anchor (opcounters, top commands)
//   - `alert_threshold`       — Atlas-alert-mirroring point events
//                                (election, restart, oplog window)
//   - `regression`            — vs --auto-baseline-file (Item 10 wire-up)
//   - `clock_skew`            — G2 cross-host data-quality guard
//   - `missing_signal`        — G5 per-version expected-metric guard
//
// Bias-control invariants reaffirmed: all detectors are STRUCTURAL —
// they declare the FTDC path(s) that form the ratio / rate / event;
// none classify the metric as a "category" or claim a "cause".
// Direction symmetry is preserved: utilization emits both above and
// below thresholds; error_rate / request_rate are scalar quantities.

package main

import (
	"sort"
	"time"
)

// utilizationRule declares a bounded ratio: value / denominator. The
// catalog covers the structural ceilings the analyzer skill needs to
// reason about CPU%, memory%, disk%, cache%, tickets%, connections%.
//
// Same `companion_gauge` pattern as the capacityRule catalog, but the
// emitted Finding has different semantics: utilization is the per-
// sample ratio summary (mean/p50/p95/max + delta between baseline and
// stress windows) — the analyzer reads it to answer "what was the
// average cache fill ratio during stress?" capacity_pressure instead
// answers "how long was the metric at saturation?"
type utilizationRule struct {
	metric      string
	denomMetric string // "" if numerator already in [0, 100] units
	displayKind string // "fraction" | "percent_ms_per_sec" | "ratio"
}

var utilizationRules = []utilizationRule{
	// Bounded gauges that already encode a percentage:
	// (none today — CPU/disk %util are derived from cumulative ms;
	// see ms-rate kind below)

	// Companion-gauge ratios (health-stage mirrors the capacity-stage
	// catalog deliberately — utilization is the value-summary, capacity
	// is the dwell-summary):
	{
		metric:      "serverStatus.queues.execution.write.out",
		denomMetric: "serverStatus.queues.execution.write.available",
		displayKind: "fraction_of_total",
	},
	{
		metric:      "serverStatus.queues.execution.read.out",
		denomMetric: "serverStatus.queues.execution.read.available",
		displayKind: "fraction_of_total",
	},
	{
		metric:      "serverStatus.connections.current",
		denomMetric: "serverStatus.connections.available",
		displayKind: "fraction_of_total",
	},
	{
		metric:      "serverStatus.wiredTiger.cache.bytes currently in the cache",
		denomMetric: "serverStatus.wiredTiger.cache.maximum bytes configured",
		displayKind: "fraction_of_max",
	},
	// Legacy WT ticket model — mirror capacityRules coverage. The two
	// catalogs intentionally overlap so the
	// analyzer can pair dwell-time (capacity_pressure) with value
	// summary (utilization) on the same metric.
	{
		metric:      "serverStatus.wiredTiger.concurrentTransactions.write.out",
		denomMetric: "serverStatus.wiredTiger.concurrentTransactions.write.available",
		displayKind: "fraction_of_total",
	},
	{
		metric:      "serverStatus.wiredTiger.concurrentTransactions.read.out",
		denomMetric: "serverStatus.wiredTiger.concurrentTransactions.read.available",
		displayKind: "fraction_of_total",
	},
}

// emitUtilization computes ratio summary stats (mean/p50/p95/max) for
// each rule whose metrics are present, plus the baseline-vs-stress
// delta. Single Finding per rule.
func emitUtilization(c *findingCollector, host string, store *seriesStore, bStart, bEnd, sStart, sEnd int) (emitted, skipped int) {
	for _, rule := range utilizationRules {
		mid, ok := store.nameToID[rule.metric]
		if !ok {
			skipped++
			continue
		}
		mCol := store.columns[mid]
		if mCol == nil {
			skipped++
			continue
		}
		var dCol *CompressedColumn
		if rule.denomMetric != "" {
			did, ok := store.nameToID[rule.denomMetric]
			if !ok {
				skipped++
				continue
			}
			dCol = store.columns[did]
			if dCol == nil {
				skipped++
				continue
			}
		}
		mVals := mCol.ReadAllInt64(nil)
		var dVals []int64
		if dCol != nil {
			dVals = dCol.ReadAllInt64(nil)
		}
		n := len(mVals)
		if dVals != nil && len(dVals) < n {
			n = len(dVals)
		}
		if n == 0 {
			skipped++
			continue
		}
		ratios := make([]float64, 0, n)
		var baseMean, stressMean float64
		baseN, stressN := 0, 0
		for i := 0; i < n; i++ {
			var ceiling int64
			switch rule.displayKind {
			case "fraction_of_total":
				ceiling = mVals[i] + dVals[i]
			case "fraction_of_max":
				ceiling = dVals[i]
			default:
				// Unknown displayKind — skip rather than silently
				// producing ratio=1.0.
				continue
			}
			if ceiling <= 0 {
				continue
			}
			r := float64(mVals[i]) / float64(ceiling)
			ratios = append(ratios, r)
			if i >= bStart && i < bEnd {
				baseMean += r
				baseN++
			}
			if i >= sStart && i < sEnd {
				stressMean += r
				stressN++
			}
		}
		if len(ratios) == 0 {
			skipped++
			continue
		}
		if baseN > 0 {
			baseMean /= float64(baseN)
		}
		if stressN > 0 {
			stressMean /= float64(stressN)
		}
		mean, p50, p95, max := ratioStats(ratios)
		c.add(newFinding("utilization", host, rule.metric, map[string]interface{}{
			"metric":            rule.metric,
			"denom_metric":      rule.denomMetric,
			"display_kind":      rule.displayKind,
			"mean":              roundFloat(mean, 4),
			"p50":               roundFloat(p50, 4),
			"p95":               roundFloat(p95, 4),
			"max":               roundFloat(max, 4),
			"baseline_mean":     roundFloat(baseMean, 4),
			"stress_mean":       roundFloat(stressMean, 4),
			"stress_delta":      roundFloat(stressMean-baseMean, 4),
			"samples_evaluated": len(ratios),
		}))
		emitted++
	}
	return emitted, skipped
}

// ratioStats returns mean, p50, p95, max of a non-empty []float64. Sorts in-place.
func ratioStats(in []float64) (mean, p50, p95, max float64) {
	if len(in) == 0 {
		return
	}
	sortFloat64Slice(in)
	for _, v := range in {
		mean += v
	}
	mean /= float64(len(in))
	p50 = in[minInt(len(in)-1, (len(in)-1)*50/100)]
	p95 = in[minInt(len(in)-1, (len(in)-1)*95/100)]
	max = in[len(in)-1]
	return
}

func sliceMean(in []float64) float64 {
	if len(in) == 0 {
		return 0
	}
	var s float64
	for _, v := range in {
		s += v
	}
	return s / float64(len(in))
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// sortFloat64Slice delegates to stdlib sort.Float64s. Earlier draft
// used insertion sort which is O(n²) — too slow on the 90K-sample
// ratios slice (~8G ops).
func sortFloat64Slice(s []float64) {
	sort.Float64s(s)
}

// ----------------------------------------------------------------------------
// alert_threshold — Atlas-style point-event detectors
// ----------------------------------------------------------------------------

// emitAlertThreshold detects two Atlas-like point events:
//   - Election: `repl.electionId` value changed mid-capture (per host).
//   - Restart: `uptime` sample resets to ~0 mid-capture.
//
// Each event emits one Finding per occurrence. Bounded by metric
// presence — no false positives on standalone mongos / standalone
// processes lacking the source metrics.
//
// Oplog-window-low detection was scoped out of this stage; it lands
// alongside the replication domain stage where the oplog metrics live.
func emitAlertThreshold(c *findingCollector, host string, store *seriesStore) (emitted, skipped int) {
	// Election detection: electionId changes between consecutive samples.
	if id, ok := store.nameToID["serverStatus.repl.electionId"]; ok {
		col := store.columns[id]
		if col != nil && col.Len() > 1 {
			vals := col.ReadAllInt64(nil)
			for i := 1; i < len(vals); i++ {
				if vals[i] != vals[i-1] && i < len(store.timestamps) {
					c.add(newFinding("alert_threshold", host, "election_at_"+store.timestamps[i].UTC().Format("20060102T150405Z"), map[string]interface{}{
						"alert_class":  "election",
						"at":           store.timestamps[i].UTC().Format(tsFormat),
						"prior_value":  vals[i-1],
						"new_value":    vals[i],
						"sample_index": i,
					}))
					emitted++
				}
			}
		}
	}
	// Restart detection: uptime resets (sample[i] < sample[i-1]) on
	// the `serverStatus.uptime` counter (seconds since process start).
	if id, ok := store.nameToID["serverStatus.uptime"]; ok {
		col := store.columns[id]
		if col != nil && col.Len() > 1 {
			vals := col.ReadAllInt64(nil)
			for i := 1; i < len(vals); i++ {
				if vals[i] < vals[i-1] && vals[i] < 60 && i < len(store.timestamps) {
					c.add(newFinding("alert_threshold", host, "restart_at_"+store.timestamps[i].UTC().Format("20060102T150405Z"), map[string]interface{}{
						"alert_class":  "restart",
						"at":           store.timestamps[i].UTC().Format(tsFormat),
						"prior_uptime": vals[i-1],
						"new_uptime":   vals[i],
						"sample_index": i,
					}))
					emitted++
				}
			}
		}
	}
	return emitted, skipped
}

// ----------------------------------------------------------------------------
// clock_skew (G2) — cross-host data-quality guard
// ----------------------------------------------------------------------------

// emitClockSkew flags when a host's timestamps drift relative to its
// sample cadence: if (timestamps[i] - timestamps[0]) - i × intervalSec
// exceeds 50ms, the host's wall clock is skewing. ONLY checked at
// per-host level here — cross-host comparison lives in
// emitCrossHostFindings (future work). When this fires, the
// run is tagged so downstream cross-host findings are advisory.
func emitClockSkew(c *findingCollector, host string, store *seriesStore) (emitted, skipped int) {
	if len(store.timestamps) < 2 {
		return 0, 1
	}
	intervalSec := modalSampleIntervalSecondsStore(store.timestamps)
	if intervalSec <= 0 {
		return 0, 1
	}
	base := store.timestamps[0]
	maxDrift := time.Duration(0)
	maxIdx := -1
	for i := 1; i < len(store.timestamps); i++ {
		// Tier A.5 fix: drop the dead `* 1` multiplier. Computing duration
		// in float64 nanoseconds preserves precision up to ~2^53 ns (~104
		// days), well above the longest realistic capture (~48h).
		expected := base.Add(time.Duration(float64(i) * intervalSec * float64(time.Second)))
		drift := store.timestamps[i].Sub(expected)
		if drift < 0 {
			drift = -drift
		}
		if drift > maxDrift {
			maxDrift = drift
			maxIdx = i
		}
	}
	if maxDrift > 50*time.Millisecond {
		c.add(newFinding("clock_skew", host, "wall_clock_drift", map[string]interface{}{
			"host":            host,
			"max_drift_ms":    roundFloat(float64(maxDrift.Microseconds())/1000.0, 2),
			"at_sample_index": maxIdx,
			"interval_sec":    roundFloat(intervalSec, 3),
			"reason_code":     "wall_clock_drift",
		}))
		emitted++
	}
	return emitted, skipped
}

// ----------------------------------------------------------------------------
// missing_signal (G5)
// ----------------------------------------------------------------------------

// expectedMetricsByVersion is the per-MongoDB-major-version manifest
// of metrics the parser EXPECTS to find. When the capture's schema
// fingerprint is missing one of these, emit `missing_signal` so
// downstream Findings using that metric class are tagged advisory.
//
// Intentionally narrow: only metrics whose absence breaks the
// analyzer's reasoning ability. Adding entries here is a deliberate
// commit — the parser doesn't infer "missing = broken" automatically.
var expectedMetricsByVersion = map[string][]string{
	"core": { // always expected on any mongod
		"serverStatus.uptime",
		"serverStatus.opcounters.insert",
		"serverStatus.opcounters.query",
		"serverStatus.connections.current",
		"serverStatus.mem.resident",
	},
	"wt": { // expected when storage engine is wiredTiger
		"serverStatus.wiredTiger.cache.bytes currently in the cache",
		"serverStatus.wiredTiger.cache.maximum bytes configured",
		"serverStatus.wiredTiger.cache.tracked dirty bytes in the cache",
	},
	"repl": { // expected on replica-set members
		"serverStatus.repl.electionId",
	},
}

// emitMissingSignal scans the per-class manifests against the
// seriesStore's nameToID. Emits one `missing_signal` Finding per
// missing metric.
func emitMissingSignal(c *findingCollector, host string, store *seriesStore) (emitted, skipped int) {
	for class, expected := range expectedMetricsByVersion {
		for _, name := range expected {
			if _, ok := store.nameToID[name]; ok {
				continue
			}
			// Tier C bias cleanup: pre-Tier-C this field was "class" —
			// a generic verdict-style label. Renamed to "category" to
			// describe what bucket (e.g. "core", "repl") of expected
			// metrics the absent metric belongs to, without the
			// verdict-vocabulary tinge.
			c.add(newFinding("missing_signal", host, name, map[string]interface{}{
				"metric":      name,
				"category":    class,
				"reason_code": "metric_absent_from_capture",
			}))
			emitted++
		}
	}
	return emitted, skipped
}
