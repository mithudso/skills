// Tier B.7: counter wrap vs reset detection. Walks the materialized
// series for each gauge-classified metric (post-build), looks for
// monotone-decrease events, and emits either kind:"counter_wrap" (when
// the previous value is near uint32 max, suggesting wraparound) or
// kind:"counter_inversion" (when the magnitude is small enough that
// it's likely a real reset, e.g. process restart or shard movement).
//
// This does NOT change the existing isCounter classification — it adds
// diagnostic Findings alongside. The store's classification is left
// alone to preserve goldens; analyzers can use the new Findings to
// recover counter semantics when a wrap was the only reason a metric
// was reclassified to gauge.
//
// Scope (per Tier B.7 docstring): "intended for legacy uint32 surfaces
// / driver-aggregated counters; modern int64 server-side counters
// never trip this." The 0.9 × 2³² threshold is the heuristic boundary;
// values above are likely wraps, values below likely resets.

package main

const counterWrapThreshold = 0.9 * float64(uint32(0xFFFFFFFF))

// emitCounterEvents walks the seriesStore looking for inversions and
// emits counter_wrap / counter_inversion Findings.
//
// CRITICAL: must read the RAW counter column (`store.columns[id].ReadAllInt64`),
// NOT `store.materializedSeries[id]`. For counters, materializedSeries holds
// the rate (first-difference) series, where `series[i] - series[i-1]` is the
// rate-of-rate (second derivative) and the threshold `0.9*2^32` applied to
// rates is unreachable in practice. The original Tier B.7 implementation made
// this error; this fix reads the raw counter values.
func emitCounterEvents(c *findingCollector, host string, store *seriesStore) (emitted, skipped int) {
	if store == nil {
		return 0, 0
	}
	var vals []int64
	for id, name := range store.names {
		// isCounter=false means a decrease was observed during classification
		// (counter wrapped or reset). Skip true monotone counters (isCounter=true).
		if id >= len(store.isCounter) || store.isCounter[id] {
			continue
		}
		col := store.columns[id]
		if col == nil || col.Len() < 2 {
			continue
		}
		vals = col.ReadAllInt64(vals[:0])
		if len(vals) < 2 {
			continue
		}
		var nonNeg, neg int
		var firstDecIdx = -1
		var firstDecPrev, firstDecCur int64
		for i := 1; i < len(vals); i++ {
			d := vals[i] - vals[i-1]
			if d >= 0 {
				nonNeg++
			} else {
				neg++
				if firstDecIdx < 0 {
					firstDecIdx = i
					firstDecPrev = vals[i-1]
					firstDecCur = vals[i]
				}
			}
		}
		if neg == 0 || firstDecIdx < 0 {
			continue
		}
		// Require ≥90% non-decreasing — otherwise it's a noisy gauge
		// the store happened to classify as counter, not a wrap/reset
		// candidate.
		if float64(nonNeg)/float64(nonNeg+neg) < 0.9 {
			continue
		}
		// Classify the first observed decrease.
		atTS := ""
		if firstDecIdx < len(store.timestamps) {
			atTS = store.timestamps[firstDecIdx].UTC().Format(tsFormat)
		}
		if float64(firstDecPrev) >= counterWrapThreshold {
			c.add(newFinding("counter_wrap", host, name+"|"+atTS, map[string]interface{}{
				"metric":       name,
				"at":           atTS,
				"sample_index": firstDecIdx,
				"prev_value":   firstDecPrev,
				"new_value":    firstDecCur,
				"reason_code":  "uint32_wraparound_suspected",
				"explanation":  "previous value near 2^32 ceiling; treat as wrap and preserve counter accumulator",
			}))
			emitted++
		} else {
			c.add(newFinding("counter_inversion", host, name+"|"+atTS, map[string]interface{}{
				"metric":       name,
				"at":           atTS,
				"sample_index": firstDecIdx,
				"prev_value":   firstDecPrev,
				"new_value":    firstDecCur,
				"reason_code":  "below_wrap_threshold",
				"explanation":  "decrease magnitude smaller than uint32 wrap threshold; likely process restart or shard movement",
			}))
			emitted++
		}
	}
	return emitted, skipped
}
