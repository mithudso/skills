// analyze_window_suggest.go — P1.2 --auto-suggest-windows.
//
// Default `--auto` uses baseline = first 25% of samples, stress = last 25%.
// A localized regression in a long capture (e.g. a 30s event in 2 hours)
// blurs out of view under that default. This stage scans the changepoint
// Findings already emitted by Stage 3, clusters breakpoints that happen
// near each other in sample index, ranks the clusters by aggregate
// effect, and emits a small set of suggested (baseline, stress) windows
// the user can re-run --auto with.
//
// Output: kind:"window_suggestion" Findings (top-K clusters). Each
// carries an embedded rerun_command field with a ready-to-run
// `--baseline FROM..TO --stress FROM..TO --auto-around-event T --auto-window W`
// invocation.
//
// The stage is currently always-on under --auto. A flag-gated variant
// could be added later but the suggestions are advisory — they never
// gate other Findings — and adding them by default keeps the user from
// having to know they exist.
package main

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// windowSuggestionClusterSpan is the maximum sample-index distance over
// which two changepoint breakpoints are considered to belong to the same
// cluster. Per the P1.2 plan: ±3 samples.
const windowSuggestionClusterSpan = 3

// windowSuggestionTopK is the number of (baseline, stress) suggestions
// emitted per host. Three is enough for a solo user to pick the most
// promising window without burying them in alternatives.
const windowSuggestionTopK = 3

// windowSuggestionMinCluster is the minimum cluster size (number of
// changepoints aligned within span) to emit a suggestion. Single-metric
// changepoints aren't interesting — Stage 3 already emits those.
const windowSuggestionMinCluster = 3

// emitWindowSuggestions reads changepoint Findings from the collector and
// emits the top-K window suggestions. Pure function of the collector +
// store; no globals. Skipped silently when there's no useful signal.
func emitWindowSuggestions(c *findingCollector, host string, store *seriesStore) (emitted int, skipped int) {
	type cpEntry struct {
		sampleIdx       int
		ts              time.Time
		metric          string
		cohensD         float64
		confidenceRank  float64
		fid             string
	}

	var entries []cpEntry
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		if k != KindChangepoint {
			continue
		}
		hf, _ := f["host"].(string)
		if hf != "" && hf != host {
			continue
		}
		idx, _ := toIntField(f["sample_index"])
		if idx < 0 {
			continue
		}
		var ts time.Time
		if t, ok := parseFindingAt(f); ok {
			ts = t
		}
		metric, _ := f["metric"].(string)
		cd, _ := f["cohens_d"].(float64)
		cr, _ := f["confidence_rank"].(float64)
		fid, _ := f["id"].(string)
		entries = append(entries, cpEntry{
			sampleIdx:      idx,
			ts:             ts,
			metric:         metric,
			cohensD:        cd,
			confidenceRank: cr,
			fid:            fid,
		})
	}
	if len(entries) == 0 {
		return 0, 1
	}

	// Sort by sample_index ascending so we can sweep clusters in one pass.
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].sampleIdx != entries[j].sampleIdx {
			return entries[i].sampleIdx < entries[j].sampleIdx
		}
		return entries[i].metric < entries[j].metric
	})

	type cluster struct {
		minIdx, maxIdx int
		score          float64 // sum of |cohens_d| * confidence_rank
		members        []cpEntry
	}
	var clusters []cluster
	var cur cluster
	flush := func() {
		if len(cur.members) > 0 {
			clusters = append(clusters, cur)
		}
		cur = cluster{}
	}
	for _, e := range entries {
		w := math.Abs(e.cohensD) * math.Max(0.1, e.confidenceRank)
		if len(cur.members) == 0 {
			cur.minIdx, cur.maxIdx = e.sampleIdx, e.sampleIdx
			cur.members = []cpEntry{e}
			cur.score = w
			continue
		}
		// Expand cluster only when the new index is within span of the
		// CURRENT maxIdx — keeps clusters compact instead of letting a
		// long chain of pairwise-near breakpoints fuse.
		if e.sampleIdx-cur.maxIdx <= windowSuggestionClusterSpan {
			cur.maxIdx = e.sampleIdx
			cur.members = append(cur.members, e)
			cur.score += w
		} else {
			flush()
			cur.minIdx, cur.maxIdx = e.sampleIdx, e.sampleIdx
			cur.members = []cpEntry{e}
			cur.score = w
		}
	}
	flush()

	// Filter by minimum cluster size + rank by score (desc, then earliest first).
	var ranked []cluster
	for _, cl := range clusters {
		if len(cl.members) >= windowSuggestionMinCluster {
			ranked = append(ranked, cl)
		}
	}
	if len(ranked) == 0 {
		return 0, 1
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].score != ranked[j].score {
			return ranked[i].score > ranked[j].score
		}
		return ranked[i].minIdx < ranked[j].minIdx
	})
	if len(ranked) > windowSuggestionTopK {
		ranked = ranked[:windowSuggestionTopK]
	}

	// Translate each cluster into baseline/stress windows.
	// Baseline: 30 minutes ending 2 minutes before cluster centroid (or
	// earliest sample in capture if 30 min isn't available before).
	// Stress: ±2 minutes around cluster centroid.
	// Window-around-event: ±2 min so --auto-around-event downstream
	// re-runs see the same neighborhood.
	intervalSec := modalSampleIntervalSecondsStore(store.timestamps)
	if intervalSec <= 0 {
		intervalSec = 1
	}
	const baselineDurSec = 30 * 60
	const stressHalfSec = 2 * 60
	const eventHalfSec = 2 * 60
	const baselineSkipSec = 2 * 60
	N := store.nSamples()
	if N == 0 {
		return 0, 1
	}

	for _, cl := range ranked {
		centroidIdx := (cl.minIdx + cl.maxIdx) / 2
		if centroidIdx >= N {
			centroidIdx = N - 1
		}
		centroidTS := store.timestamps[centroidIdx]
		baselineEnd := centroidTS.Add(-time.Duration(baselineSkipSec) * time.Second)
		baselineStart := baselineEnd.Add(-time.Duration(baselineDurSec) * time.Second)
		// Clamp to capture bounds.
		captureStart := store.timestamps[0]
		captureEnd := store.timestamps[N-1]
		if baselineStart.Before(captureStart) {
			baselineStart = captureStart
		}
		if baselineEnd.Before(captureStart) {
			baselineEnd = captureStart
		}
		if baselineEnd.After(captureEnd) {
			baselineEnd = captureEnd
		}
		// Skip this cluster if the clamping produced a degenerate window.
		if !baselineEnd.After(baselineStart) {
			continue
		}
		stressStart := centroidTS.Add(-time.Duration(stressHalfSec) * time.Second)
		stressEnd := centroidTS.Add(time.Duration(stressHalfSec) * time.Second)
		if stressStart.Before(captureStart) {
			stressStart = captureStart
		}
		if stressEnd.After(captureEnd) {
			stressEnd = captureEnd
		}
		eventWindow := time.Duration(eventHalfSec*2) * time.Second

		// Collect top metrics in the cluster for the rationale.
		clMetrics := []string{}
		clFindingIDs := []string{}
		for _, m := range cl.members {
			if m.metric != "" {
				clMetrics = append(clMetrics, m.metric)
			}
			if m.fid != "" {
				clFindingIDs = append(clFindingIDs, m.fid)
			}
		}
		// Deduplicate metrics while preserving order.
		seen := map[string]bool{}
		uniqMetrics := make([]string, 0, len(clMetrics))
		for _, m := range clMetrics {
			if seen[m] {
				continue
			}
			seen[m] = true
			uniqMetrics = append(uniqMetrics, m)
		}
		// Cap rationale list — we're advisory, not exhaustive.
		if len(uniqMetrics) > 8 {
			uniqMetrics = uniqMetrics[:8]
		}

		atStr := centroidTS.UTC().Format(tsFormat)
		baselineStartStr := baselineStart.UTC().Format(time.RFC3339)
		baselineEndStr := baselineEnd.UTC().Format(time.RFC3339)
		stressStartStr := stressStart.UTC().Format(time.RFC3339)
		stressEndStr := stressEnd.UTC().Format(time.RFC3339)
		fields := map[string]interface{}{
			"at":                     atStr,
			"centroid_sample_index":  centroidIdx,
			"cluster_size":           len(cl.members),
			"cluster_score":          roundFloat(cl.score, 4),
			"baseline":               map[string]string{"from": baselineStartStr, "to": baselineEndStr},
			"stress":                 map[string]string{"from": stressStartStr, "to": stressEndStr},
			"suggested_around_event": atStr,
			"suggested_window":       eventWindow.String(),
			"supporting_finding_ids": clFindingIDs,
			"top_metrics":            uniqMetrics,
			"rationale": fmt.Sprintf(
				"%d changepoints clustered within ±%d samples around index %d (score %.2f); narrowed windows isolate the event.",
				len(cl.members), windowSuggestionClusterSpan, centroidIdx, cl.score,
			),
			"rerun_command": fmt.Sprintf(
				"ftdc_parser <path> --auto --baseline %s..%s --stress %s..%s --auto-around-event %s --auto-window %s",
				baselineStartStr, baselineEndStr, stressStartStr, stressEndStr, atStr, eventWindow.String(),
			),
		}
		c.add(newFinding("window_suggestion", host, atStr, fields))
		emitted++
	}
	return emitted, 0
}

// toIntField extracts an int from a Finding-field value that may be int,
// int64, float64, or absent. Returns (-1, false) when extraction fails.
func toIntField(v interface{}) (int, bool) {
	switch x := v.(type) {
	case int:
		return x, true
	case int64:
		return int(x), true
	case float64:
		return int(x), true
	}
	return -1, false
}
