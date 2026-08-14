// analyze_seal.go — PALI seal-event detection (P0.NEW-PALI-1).
//
// "Seal" is the PALI distributed-log transition where the log generation
// changes (a new primary term begins writing, the prior log is sealed).
// Without a structural anchor for seal windows the parser cannot compose
// any of the PALI-specific patterns that follow (role_asymmetry within a
// seal window; standby_seal_topology_skew; standby_read_log_wedge).
//
// Detection rule (D1 in the P0 plan):
//   A `seal_event` Finding fires when ALL of these hold within a 60-second
//   window centered on a `term_change` Finding:
//     (1) `term_change` present.
//     (2) `local_spike` on at least one of:
//             disagg.logServer.appendLogTotal
//             disagg.logServer.appendLogSuccess
//             disagg.logServer.discontinuousCheckpointInstallCount
//     (3) EITHER
//           (3a) `gap` Finding on a `disagg.phylog.*` rate stream, OR
//           (3b) `discontinuousCheckpointInstallCount` delta >= 1 over
//                 the seal window.
//     (4) NOT preceded within 30 s by a `counter_inversion` on a
//          `disagg.*` metric (rules out plain process restart, which can
//          mimic conditions 1-3 spuriously).
//
// Confidence tiers:
//    4/4 conditions:      confidence_rank = 0.85, kind = "seal_event"
//    3/4 (cond4 missing): confidence_rank = 0.60, kind = "seal_event"
//    3/4 (cond3 missing): confidence_rank = 0.55, kind = "seal_event"
//    3/4 (cond2 missing): confidence_rank = 0.45, kind = "seal_event"
//   ≤2/4:                 confidence_rank = 0.35, kind = "seal_event_candidate"
//
// Bias-control: the Finding declares the structural evidence (the
// supporting_finding_ids it composed) and lets the analyzer narrate.
// No verdicts, no severity, no role.
package main

import (
	"strings"
	"time"
)

// sealEventWindow is half the time window the composer considers around
// each candidate term_change anchor.
const sealEventHalfWindow = 30 * time.Second

// sealEventLogServerMetrics is the closed set of disagg.logServer metric
// suffixes whose local_spike Findings the composer treats as condition #2.
// Order is the order they're cited in supporting_finding_ids if multiple
// fire.
var sealEventLogServerMetrics = []string{
	"serverStatus.metrics.disagg.logServer.appendLogTotal",
	"serverStatus.metrics.disagg.logServer.appendLogSuccess",
	"serverStatus.metrics.disagg.logServer.discontinuousCheckpointInstallCount",
}

// disaggMetricPrefix is the FTDC path prefix for any PALI metric. Used by
// condition #4 (counter_inversion on disagg.*).
const disaggMetricPrefix = "serverStatus.metrics.disagg."

// phylogMetricPrefix narrows to the phylog subtree (the standby's lag /
// rate observations). Gap Findings on metrics under this prefix satisfy
// condition #3a.
const phylogMetricPrefix = "serverStatus.metrics.disagg.phylog."

// discontinuousCheckpointInstallCount is the metric whose delta satisfies
// condition #3b. The PALI domain stage emits a `checkpoint_install_anomaly`
// Finding when this counter increases by ≥1 — the seal composer can read
// either the column directly or the Finding's presence as evidence.
const discontinuousCheckpointMetric = "serverStatus.metrics.disagg.logServer.discontinuousCheckpointInstallCount"

// emitSealEvents runs after the replication, election_storm, PALI, and
// counter-inversion stages have all emitted to the collector. It scans
// the existing Findings, composes seal_event(s), and emits them.
//
// Side effect: the collector grows by 0..N seal_event[/seal_event_candidate]
// Findings. Pure function of the collector + store; no globals.
//
// Gated on isPALICapture — non-PALI captures have nothing meaningful here.
func emitSealEvents(c *findingCollector, host string, store *seriesStore) (emitted int, skipped int) {
	if !isPALICapture(store) {
		return 0, 1
	}

	// Index Findings we care about by parsed timestamp. We use the
	// already-emitted Findings as evidence rather than re-running any
	// detection — keeps the composer pure and respects the bias-control
	// invariant that detection stages emit shapes; composers compose.
	type tsEntry struct {
		f  Finding
		ts time.Time
	}
	var termChanges, logServerSpikes, phylogGaps, disaggInversions []tsEntry

	for _, f := range c.items {
		k, _ := f["kind"].(string)
		hf, _ := f["host"].(string)
		if hf != "" && hf != host {
			continue
		}
		switch k {
		case KindTermChange:
			if ts, ok := parseFindingAt(f); ok {
				termChanges = append(termChanges, tsEntry{f: f, ts: ts})
			}
		case "local_spike":
			metric, _ := f["metric"].(string)
			if metric == "" {
				metric, _ = f["name"].(string)
			}
			if isSealEventLogServerMetric(metric) {
				if ts, ok := parseFindingAt(f); ok {
					logServerSpikes = append(logServerSpikes, tsEntry{f: f, ts: ts})
				}
			}
		case KindGap:
			metric, _ := f["metric"].(string)
			if strings.HasPrefix(metric, phylogMetricPrefix) {
				if ts, ok := parseFindingAt(f); ok {
					phylogGaps = append(phylogGaps, tsEntry{f: f, ts: ts})
				}
			}
		case "counter_inversion":
			metric, _ := f["metric"].(string)
			if strings.HasPrefix(metric, disaggMetricPrefix) {
				if ts, ok := parseFindingAt(f); ok {
					disaggInversions = append(disaggInversions, tsEntry{f: f, ts: ts})
				}
			}
		}
	}

	if len(termChanges) == 0 {
		return 0, 1
	}

	// Condition #3b helper: did the discontinuousCheckpointInstallCount
	// delta within the seal window exceed 0? Read directly from the
	// store column so we don't have to depend on the PALI stage's
	// checkpoint_install_anomaly Finding being present (analyzer-side
	// gates may suppress it).
	discontiID, hasDiscontiCol := store.nameToID[discontinuousCheckpointMetric]
	var discontiCol []int64
	if hasDiscontiCol {
		discontiCol = store.readColumnInt64(discontiID, nil)
	}

	for _, tc := range termChanges {
		windowFrom := tc.ts.Add(-sealEventHalfWindow)
		windowTo := tc.ts.Add(sealEventHalfWindow)

		// Condition #2: local_spike on at least one logServer metric.
		spikeIDs := []string{}
		spikeMetrics := []string{}
		earliestEvidence := tc.ts
		latestEvidence := tc.ts
		for _, sp := range logServerSpikes {
			if sp.ts.Before(windowFrom) || sp.ts.After(windowTo) {
				continue
			}
			if fid, ok := sp.f["id"].(string); ok {
				spikeIDs = append(spikeIDs, fid)
			}
			if m, ok := sp.f["metric"].(string); ok && m != "" {
				spikeMetrics = append(spikeMetrics, m)
			} else if m, ok := sp.f["name"].(string); ok {
				spikeMetrics = append(spikeMetrics, m)
			}
			if sp.ts.Before(earliestEvidence) {
				earliestEvidence = sp.ts
			}
			if sp.ts.After(latestEvidence) {
				latestEvidence = sp.ts
			}
		}
		cond2 := len(spikeIDs) > 0

		// Condition #3: phylog gap OR discontinuousCheckpointInstallCount
		// delta within the window.
		gapIDs := []string{}
		for _, g := range phylogGaps {
			if g.ts.Before(windowFrom) || g.ts.After(windowTo) {
				continue
			}
			if fid, ok := g.f["id"].(string); ok {
				gapIDs = append(gapIDs, fid)
			}
			if g.ts.Before(earliestEvidence) {
				earliestEvidence = g.ts
			}
			if g.ts.After(latestEvidence) {
				latestEvidence = g.ts
			}
		}
		discontiDelta := int64(0)
		if hasDiscontiCol && len(discontiCol) > 1 {
			loIdx, hiIdx := indicesForWindow(store, windowFrom, windowTo)
			if hiIdx > loIdx && hiIdx <= len(discontiCol) {
				// Baseline is the last sample BEFORE the window so that any
				// increment occurring at loIdx itself is captured. When
				// loIdx==0 there is no prior sample; use discontiCol[0].
				baseline := discontiCol[loIdx]
				if loIdx > 0 {
					baseline = discontiCol[loIdx-1]
				}
				discontiDelta = discontiCol[hiIdx-1] - baseline
			}
		}
		cond3 := len(gapIDs) > 0 || discontiDelta >= 1

		// Condition #4 (negative): no disagg counter_inversion in the
		// 30s preceding the term_change anchor. If present, downgrade.
		var cond4Block []string
		for _, inv := range disaggInversions {
			if inv.ts.Before(tc.ts.Add(-sealEventHalfWindow)) || inv.ts.After(tc.ts) {
				continue
			}
			if fid, ok := inv.f["id"].(string); ok {
				cond4Block = append(cond4Block, fid)
			}
		}
		cond4 := len(cond4Block) == 0

		// Score conditions. #1 (term_change) is always present at this
		// point in the loop body. Count the truthy among 2, 3, 4.
		conditionsMet := 1 // condition 1 always met (we have a term_change)
		if cond2 {
			conditionsMet++
		}
		if cond3 {
			conditionsMet++
		}
		if cond4 {
			conditionsMet++
		}

		var (
			kind            string
			confidence      float64
			supportingNotes []string
		)
		switch {
		case conditionsMet == 4:
			kind = KindSealEvent
			confidence = 0.85
		case conditionsMet == 3:
			kind = KindSealEvent
			// Which condition is missing determines confidence.
			// cond4=false: counter_inversion present (restart-masking concern) → 0.60
			// cond2=false: no log-server spike (weaker structural evidence)   → 0.45
			// cond3=false: no phylog gap / checkpoint delta                    → 0.55
			if !cond4 {
				confidence = 0.60
			} else if !cond2 {
				confidence = 0.45
			} else {
				confidence = 0.55 // cond3 missing
			}
		default:
			kind = KindSealEventCandidate
			confidence = 0.35
		}

		// Compose supporting_finding_ids in deterministic order: the
		// term_change first, then logServer spikes, then phylog gaps.
		supportingIDs := []string{}
		if tcID, ok := tc.f["id"].(string); ok {
			supportingIDs = append(supportingIDs, tcID)
		}
		supportingIDs = append(supportingIDs, spikeIDs...)
		supportingIDs = append(supportingIDs, gapIDs...)

		if discontiDelta >= 1 {
			supportingNotes = append(supportingNotes,
				"discontinuousCheckpointInstallCount delta within window")
		}
		if !cond4 {
			supportingNotes = append(supportingNotes,
				"counter_inversion on disagg.* in the 30s preceding the term_change (kind downgraded; see contradicting_finding_ids)")
		}

		atStr := tc.ts.UTC().Format(tsFormat)
		fields := map[string]interface{}{
			"at":                    atStr,
			"seal_window":           map[string]string{"from": earliestEvidence.UTC().Format(tsFormat), "to": latestEvidence.UTC().Format(tsFormat)},
			"supporting_finding_ids": supportingIDs,
			"confidence_rank":       confidence,
			"conditions_met":        conditionsMet,
			"evidence_summary": map[string]interface{}{
				"term_change":                          true,
				"local_spike_metrics":                  spikeMetrics,
				"phylog_gap_count":                     len(gapIDs),
				"discontinuous_checkpoint_install_delta": discontiDelta,
				"counter_inversion_blockers":           len(cond4Block),
			},
		}
		if len(supportingNotes) > 0 {
			fields["notes"] = supportingNotes
		}
		if len(cond4Block) > 0 {
			fields["contradicting_finding_ids"] = cond4Block
		}
		c.add(newFinding(kind, host, atStr, fields))
		emitted++
	}
	return emitted, 0
}

// parseFindingAt returns the time.Time encoded in Finding["at"], or zero
// + false if absent / unparseable.
func parseFindingAt(f Finding) (time.Time, bool) {
	at, ok := f["at"].(string)
	if !ok || at == "" {
		return time.Time{}, false
	}
	// Two formats are used in the codebase:
	//   tsFormat (millisecond precision, gap/local_spike)
	//   RFC3339Nano (run_metadata window_*)
	// Try the explicit millisecond format first; fall back to RFC3339Nano.
	if t, err := time.Parse(tsFormat, at); err == nil {
		return t.UTC(), true
	}
	if t, err := time.Parse(time.RFC3339Nano, at); err == nil {
		return t.UTC(), true
	}
	return time.Time{}, false
}

// isSealEventLogServerMetric reports whether the given metric path is one
// of the three logServer counters tracked by condition #2.
func isSealEventLogServerMetric(name string) bool {
	for _, m := range sealEventLogServerMetrics {
		if name == m {
			return true
		}
	}
	return false
}

// indicesForWindow maps a [from, to] wall-clock window to the half-open
// sample-index bounds [lo, hi). Returns (0, 0) when the store is empty.
func indicesForWindow(store *seriesStore, from, to time.Time) (lo, hi int) {
	return timeRangeToSampleIdx(store.timestamps, from, to)
}
