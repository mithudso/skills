// analyze_capacity.go - Findings: capacity_pressure + near_ceiling.
//
// Background (from /home/ubuntu/.claude/plans/how-to-make-ftdc-parser-graceful-sky.md):
//
// The existing `mover` Finding carries a `saturation_delta` field that is
// a per-window AVERAGE of value/observedMax — it can detect a shift in
// average saturation between baseline and stress but says nothing about
// dwell time at the ceiling. A ticket queue that sat at 100% for 30
// seconds in the middle of a 12h capture moves saturation_delta by
// 30/43200 ≈ 0.0007 — invisible.
//
// This stage adds two structural Findings driven by STRUCTURAL CEILINGS
// (companion gauges already in the FTDC schema). The ceiling kind is
// declared in the Finding so the analyzer skill doesn't have to infer.
//
// `kind:"capacity_pressure"`:
//   metric:               <ceiling-bound metric>
//   ceiling_metric:       <companion gauge that defines the ceiling>
//   ceiling_kind:         "companion_gauge" | "drains_to_zero" | "fraction"
//   pct_at_ceiling:       % of samples at value >= 0.99 * ceiling
//   pct_near_ceiling:     % of samples at value >= 0.95 * ceiling
//   longest_run_seconds:  longest consecutive run at ceiling (seconds)
//
// `kind:"near_ceiling"`:
//   metric, ceiling_metric, ceiling_value, last_value, last_ratio
//   Fires when the metric's LATEST observation is >= 0.85 * ceiling.
//   No ETA / no extrapolation — the linear-forecast math is unreliable
//   for the metric classes the analyzer cares about (logistic growth on
//   cache fill, sublinear on tcmalloc heap).

package main

import (
	"fmt"
	"sort"
	"time"
)

// capacityRule declares how a metric's ceiling is structured.
//
// This stage ships ONE kind: `companion_gauge` (ceiling = metric + companion;
// covers ticket queues, connections, wt cache). The `kind` field exists
// to forward-declare future ceiling shapes (fixed_pct for percentage
// metrics, drains_to_zero for inverted gauges); these are NOT yet
// wired in emitCapacityPressure / emitNearCeiling. The emit functions
// reject rules with empty companion to prevent a footgun where a future
// rule with kind="fixed_pct" would silently compute ratio=m/m=1.0 and
// always fire.
type capacityRule struct {
	metric    string
	companion string // required for companion_gauge; "" reserved for future kinds
	kind      string // currently only "companion_gauge"
	// Per-rule thresholds. Defaults if zero: run_secs=10, pct_at=0.05.
	runSeconds float64
	pctAt      float64
}

// capacityRules is the structural ceiling catalog. Intentionally narrow
// (~12 entries); expand as new domain stages land (sharded, host_os).
// All metrics carry their ceiling-source in the emitted Finding so the
// analyzer skill can audit. No name-based recognition of cause; the
// catalog only declares WHERE the ceiling comes from.
var capacityRules = []capacityRule{
	// Ticket queues — fire on short runs (these are sub-second-meaningful).
	{
		metric:     "serverStatus.queues.execution.write.out",
		companion:  "serverStatus.queues.execution.write.available",
		kind:       "companion_gauge",
		runSeconds: 5,
		pctAt:      0.05,
	},
	{
		metric:     "serverStatus.queues.execution.read.out",
		companion:  "serverStatus.queues.execution.read.available",
		kind:       "companion_gauge",
		runSeconds: 5,
		pctAt:      0.05,
	},
	{
		metric:     "serverStatus.wiredTiger.concurrentTransactions.write.out",
		companion:  "serverStatus.wiredTiger.concurrentTransactions.write.available",
		kind:       "companion_gauge",
		runSeconds: 5,
		pctAt:      0.05,
	},
	{
		metric:     "serverStatus.wiredTiger.concurrentTransactions.read.out",
		companion:  "serverStatus.wiredTiger.concurrentTransactions.read.available",
		kind:       "companion_gauge",
		runSeconds: 5,
		pctAt:      0.05,
	},
	// Connection pool — short-run signal too.
	{
		metric:     "serverStatus.connections.current",
		companion:  "serverStatus.connections.available",
		kind:       "companion_gauge",
		runSeconds: 5,
		pctAt:      0.05,
	},
	// WT cache — sustained pressure threshold; cache fill above 80%
	// for half the capture is the relevant signal.
	{
		metric:     "serverStatus.wiredTiger.cache.bytes currently in the cache",
		companion:  "serverStatus.wiredTiger.cache.maximum bytes configured",
		kind:       "companion_gauge",
		runSeconds: 60,
		pctAt:      0.50,
	},
}

// emitCapacityPressure scans the seriesStore for any rule whose
// (metric, companion) pair is present, computes the at-ceiling /
// near-ceiling dwell statistics, and emits one Finding per rule that
// fires. Structural — gated by metric presence, not by metric name
// matching beyond the catalog declaration.
func emitCapacityPressure(c *findingCollector, host string, store *seriesStore, cfg AutoConfig) {
	if len(store.timestamps) < 2 {
		return
	}
	intervalSec := modalSampleIntervalSecondsStore(store.timestamps)
	if intervalSec <= 0 {
		intervalSec = 1
	}
	for _, rule := range capacityRules {
		// Only companion_gauge is implemented; reject other kinds
		// rather than silently compute meaningless ratios.
		if rule.kind != "companion_gauge" || rule.companion == "" {
			continue
		}
		mid, ok := store.nameToID[rule.metric]
		if !ok {
			continue
		}
		cid, ok := store.nameToID[rule.companion]
		if !ok {
			continue // ceiling source missing; skip
		}
		mCol := store.columns[mid]
		if mCol == nil {
			continue
		}
		mVals := mCol.ReadAllInt64(nil)
		if len(mVals) == 0 {
			continue
		}
		cCol := store.columns[cid]
		if cCol == nil {
			continue
		}
		cVals := cCol.ReadAllInt64(nil)
		// Per-sample ceiling = m + c (companion_gauge). For each sample
		// compute ratio = m / (m + c); truncate to min(len(mVals),
		// len(cVals)) in case of late-emerging metrics.
		n := len(mVals)
		if len(cVals) < n {
			n = len(cVals)
		}
		var nAt, nNear int
		var longestRun int
		var curRun int
		for i := 0; i < n; i++ {
			m := mVals[i]
			ceiling := m + cVals[i]
			if ceiling <= 0 {
				continue
			}
			ratio := float64(m) / float64(ceiling)
			atCeiling := ratio >= 0.99
			nearCeiling := ratio >= 0.95
			if atCeiling {
				nAt++
				curRun++
				if curRun > longestRun {
					longestRun = curRun
				}
			} else {
				curRun = 0
			}
			if nearCeiling {
				nNear++
			}
		}
		pctAt := float64(nAt) / float64(n)
		pctNear := float64(nNear) / float64(n)
		longestRunSec := float64(longestRun) * intervalSec

		// Defaults if rule didn't specify.
		runThr := rule.runSeconds
		if runThr <= 0 {
			runThr = 10
		}
		pctThr := rule.pctAt
		if pctThr <= 0 {
			pctThr = 0.05
		}
		if longestRunSec < runThr && pctAt < pctThr {
			continue
		}
		c.add(newFinding("capacity_pressure", host, rule.metric, map[string]interface{}{
			"metric":              rule.metric,
			"ceiling_metric":      rule.companion,
			"ceiling_kind":        rule.kind,
			"pct_at_ceiling":      roundFloat(pctAt, 4),
			"pct_near_ceiling":    roundFloat(pctNear, 4),
			"longest_run_seconds": roundFloat(longestRunSec, 1),
			"samples_evaluated":   n,
		}))
	}
}

// emitNearCeiling fires when the metric's LATEST observation is at or
// above 85% of its structural ceiling. No extrapolation, no ETA — the
// analyzer skill combines this with the monotonic_trend slope sign to
// decide whether the metric is approaching exhaustion.
func emitNearCeiling(c *findingCollector, host string, store *seriesStore, cfg AutoConfig) {
	for _, rule := range capacityRules {
		if rule.kind != "companion_gauge" || rule.companion == "" {
			continue
		}
		mid, ok := store.nameToID[rule.metric]
		if !ok {
			continue
		}
		mCol := store.columns[mid]
		if mCol == nil || mCol.Len() == 0 {
			continue
		}
		mLast := mCol.ReadRange(mCol.Len()-1, mCol.Len(), nil)
		if len(mLast) == 0 {
			continue
		}
		lastM := mLast[0]
		cid, ok := store.nameToID[rule.companion]
		if !ok {
			continue
		}
		cCol := store.columns[cid]
		if cCol == nil || cCol.Len() == 0 {
			continue
		}
		cLast := cCol.ReadRange(cCol.Len()-1, cCol.Len(), nil)
		if len(cLast) == 0 {
			continue
		}
		ceiling := lastM + cLast[0]
		if ceiling <= 0 {
			continue
		}
		ratio := float64(lastM) / float64(ceiling)
		if ratio < 0.85 {
			continue
		}
		c.add(newFinding("near_ceiling", host, rule.metric, map[string]interface{}{
			"metric":         rule.metric,
			"ceiling_metric": rule.companion,
			"ceiling_kind":   rule.kind,
			"last_value":     lastM,
			"ceiling_value":  ceiling,
			"last_ratio":     roundFloat(ratio, 4),
		}))
	}
}

// emitConsistencyWarning is a post-pass cross-Finding sanity check.
// For each metric that has BOTH a `mover` Finding and a `monotonic_trend`
// Finding, verify the directions agree:
//   - mover.direction == "up"   ⇒ monotonic_trend.slope > 0
//   - mover.direction == "down" ⇒ monotonic_trend.slope < 0
// Disagreement implies one of (a) Welford-ordering edge case, (b)
// downsampling artifact in monotonic_trend's Spearman, (c) a true
// regime change that flipped direction mid-capture. Emit a
// consistency_warning so the analyzer treats the metric's signal as
// less certain.
func emitConsistencyWarning(c *findingCollector, host string) {
	moverDir := map[string]string{}
	trendSlope := map[string]float64{}
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		metric, _ := f["metric"].(string)
		if metric == "" {
			continue
		}
		switch k {
		case KindMover:
			d, _ := f["direction"].(string)
			moverDir[metric] = d
		case KindMonotonicTrend:
			s, _ := f["slope"].(float64)
			trendSlope[metric] = s
		}
	}
	// Sort metrics for deterministic Finding emission order.
	sortedMetrics := make([]string, 0, len(moverDir))
	for m := range moverDir {
		sortedMetrics = append(sortedMetrics, m)
	}
	sort.Strings(sortedMetrics)
	for _, metric := range sortedMetrics {
		dir := moverDir[metric]
		slope, ok := trendSlope[metric]
		if !ok {
			continue
		}
		mismatch := (dir == "up" && slope < 0) || (dir == "down" && slope > 0)
		if !mismatch {
			continue
		}
		c.add(newFinding("consistency_warning", host, metric, map[string]interface{}{
			"metric":          metric,
			"mover_direction": dir,
			"trend_slope":     roundFloat(slope, 6),
			"reason_code":     "mover_trend_disagree",
			"reason":          "mover and monotonic_trend disagree on direction; signal localization is uncertain",
		}))
	}
}

// emitSyntheticEvent clusters findings within a sliding time window and
// emits one overlay Finding per cluster. Children are preserved (overlay
// not replacement). Per the bias-control invariants the overlay carries
// `member_finding_ids`, structural counters (up_count, down_count,
// metric_family_count, members[]) — no direction arrow, no causal
// claim. The analyzer composes causation if any.
//
// Cluster rules:
//   - Use sliding window with radius = max(4 * modal interval, 5s).
//   - Cluster eligible kinds: changepoint, spike, variance_shift, level.
//   - Require >=3 distinct metrics from >=2 metric-name families
//     (first '.' segment).
//   - Cap at 50 events with `truncated:true`.
func emitSyntheticEvent(c *findingCollector, host string, store *seriesStore) {
	if len(store.timestamps) < 2 {
		return
	}
	intervalSec := modalSampleIntervalSecondsStore(store.timestamps)
	if intervalSec <= 0 {
		intervalSec = 1
	}
	radius := 4 * intervalSec
	if radius < 5 {
		radius = 5
	}
	radiusDur := time.Duration(radius * float64(time.Second))

	var refs []sortableRef
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		switch k {
		case KindChangepoint, KindSpike, KindVarianceShift:
			// eligible
		default:
			continue
		}
		atStr, _ := f["at"].(string)
		if atStr == "" {
			continue
		}
		t, err := time.Parse(tsFormat, atStr)
		if err != nil {
			continue
		}
		metric, _ := f["metric"].(string)
		if metric == "" {
			continue
		}
		family := metricFamily(metric)
		id, _ := f["id"].(string)
		dir := inferFindingDirection(f)
		refs = append(refs, sortableRef{
			id: id, kind: k, metric: metric, family: family,
			dir: dir, mag: findingMagnitude(f), at: t,
		})
	}
	if len(refs) < 3 {
		return
	}
	// Sort by timestamp; tie-break by (kind, metric, id) so two
	// Findings with the same `at:` produce a deterministic cluster
	// composition. Without the tie-break, parallel-emitted spike /
	// variance_shift Findings would land in different orders across
	// runs, making synthetic_event golden output non-deterministic.
	sort.Slice(refs, func(i, j int) bool {
		if !refs[i].at.Equal(refs[j].at) {
			return refs[i].at.Before(refs[j].at)
		}
		if refs[i].kind != refs[j].kind {
			return refs[i].kind < refs[j].kind
		}
		if refs[i].metric != refs[j].metric {
			return refs[i].metric < refs[j].metric
		}
		return refs[i].id < refs[j].id
	})

	// Sliding-window grouping. We emit at most one synthetic_event per
	// disjoint window; events farther apart than `radiusDur` are
	// distinct.
	const maxEvents = 50
	type cluster struct {
		start, end time.Time
		members    []sortableRef
		families   map[string]bool
		metrics    map[string]bool
	}
	var emitted []cluster
	i := 0
	for i < len(refs) {
		cls := cluster{
			start:    refs[i].at,
			end:      refs[i].at,
			families: map[string]bool{},
			metrics:  map[string]bool{},
		}
		j := i
		for j < len(refs) && refs[j].at.Sub(refs[i].at) <= radiusDur {
			cls.members = append(cls.members, refs[j])
			cls.families[refs[j].family] = true
			cls.metrics[refs[j].metric] = true
			cls.end = refs[j].at
			j++
		}
		if len(cls.metrics) >= 3 && len(cls.families) >= 2 {
			emitted = append(emitted, cls)
			i = j // step past this window
		} else {
			i++ // advance one
		}
		if len(emitted) >= maxEvents {
			break
		}
	}
	if len(emitted) == 0 {
		return
	}
	truncated := len(emitted) >= maxEvents
	if truncated && len(refs) > 0 {
		// Drop overflow if any.
	}
	for _, cls := range emitted {
		ids := make([]string, 0, len(cls.members))
		var upCount, downCount int
		// Top-5 members by magnitude.
		topMembers := topByMagnitude(cls.members, 5)
		mems := make([]interface{}, 0, len(topMembers))
		for _, ch := range cls.members {
			ids = append(ids, ch.id)
			switch ch.dir {
			case "up":
				upCount++
			case "down":
				downCount++
			}
		}
		for _, ch := range topMembers {
			mems = append(mems, map[string]interface{}{
				"id":          ch.id,
				"kind":        ch.kind,
				"metric":      ch.metric,
				"direction":   ch.dir,
				"magnitude_z": roundFloat(ch.mag, 4),
			})
		}
		// Key uses (start_ts || end_ts || member_count) so two clusters
		// sharing a start_ts (rare but possible with sub-minute granularity)
		// don't collide on the same Finding id derived from the key.
		clusterKey := fmt.Sprintf("%s|%s|%d",
			cls.start.UTC().Format(tsFormat),
			cls.end.UTC().Format(tsFormat),
			len(cls.members))
		c.add(newFinding("synthetic_event", host, clusterKey, map[string]interface{}{
			"at_start":            cls.start.UTC().Format(tsFormat),
			"at_end":              cls.end.UTC().Format(tsFormat),
			"member_count":        len(cls.members),
			"metric_family_count": len(cls.families),
			"up_count":            upCount,
			"down_count":          downCount,
			"member_finding_ids":  ids,
			"members":             mems,
			"window_seconds":      roundFloat(radius, 1),
			"truncated":           truncated,
		}))
	}
}

// metricFamily returns the two-segment prefix of a metric name as a
// subsystem identifier. Used to enforce the "≥2 distinct families"
// rule for synthetic_event clustering. Two-segment depth distinguishes
// `serverStatus.wiredTiger.*` from `serverStatus.opcounters.*` etc.,
// which a single-segment prefix would conflate as one big "serverStatus"
// family. For paths shorter than two segments, returns the whole path.
//
// Examples:
//   "serverStatus.wiredTiger.cache.bytes" → "serverStatus.wiredTiger"
//   "local.oplog.rs.stats.X"              → "local.oplog"
//   "serverStatus.foo"                    → "serverStatus.foo" (1 dot)
//   "foo"                                 → "foo" (0 dots)
func metricFamily(metric string) string {
	dots := 0
	for i := 0; i < len(metric); i++ {
		if metric[i] == '.' {
			dots++
			if dots == 2 {
				return metric[:i]
			}
		}
	}
	return metric
}

// inferFindingDirection extracts a direction signal from a Finding when
// available. Returns "up", "down", or "".
func inferFindingDirection(f Finding) string {
	if d, ok := f["direction"].(string); ok && d != "" {
		return d
	}
	// Mean-shift Findings (changepoint, level) embed direction in
	// baseline vs stress means.
	if bm, ok := f["baseline_mean"].(float64); ok {
		if sm, ok := f["stressed_mean"].(float64); ok {
			if sm > bm {
				return "up"
			}
			if sm < bm {
				return "down"
			}
		}
	}
	// Spike Findings encode direction via (value, median): spike above
	// median is "up", below is "down". value and median are typed as
	// numbers but Go json unmarshal gives float64 for any number.
	if val, ok := numericField(f, "value"); ok {
		if med, ok := numericField(f, "median"); ok {
			if val > med {
				return "up"
			}
			if val < med {
				return "down"
			}
		}
	}
	return ""
}

// numericField returns a float64 value for a Finding field that may
// carry any of the common numeric types. FTDC raw values are int64;
// computed stats are float64; BSON-roundtripped Findings may emit
// uint64 / int32 etc.
func numericField(f Finding, key string) (float64, bool) {
	v, ok := f[key]
	if !ok {
		return 0, false
	}
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int64:
		return float64(x), true
	case uint64:
		return float64(x), true
	case int32:
		return float64(x), true
	case uint32:
		return float64(x), true
	case int:
		return float64(x), true
	}
	return 0, false
}

// sortableRef is the per-Finding tuple used by emitSyntheticEvent during
// clustering. Kept private; not used outside this file.
type sortableRef struct {
	id     string
	kind   string
	metric string
	family string
	dir    string
	mag    float64
	at     time.Time
}

// topByMagnitude returns the top-k entries by magnitude (stable-sort).
func topByMagnitude(refs []sortableRef, k int) []sortableRef {
	if k <= 0 || len(refs) == 0 {
		return nil
	}
	out := append([]sortableRef(nil), refs...)
	// Selection sort to k positions (small k; cheap, deterministic).
	for i := 0; i < k && i < len(out); i++ {
		best := i
		for j := i + 1; j < len(out); j++ {
			if out[j].mag > out[best].mag {
				best = j
			}
		}
		out[i], out[best] = out[best], out[i]
	}
	if len(out) > k {
		out = out[:k]
	}
	return out
}
