// Tier G: capability extensions surviving the round-3 gap audit.
//
// New detectors (G.1):
//   - apply_saturation: secondary replication apply-side bottleneck
//   - election_storm: ≥3 term_change Findings clustered within 5 min
//   - checkpoint_stall: dual-path — generation counter stalled (primary) +
//     completed-checkpoint exceeded syncdelay (secondary)
//   - eviction_divergence: WT cache read-into vs evicted divergence
//
// New observability fields (G.3):
//   - capture_quality_score: composite per-run 0-100 health score
//
// Each detector branches on a Tier-5-class structural prefix (replication,
// disagg / wiredTiger). Tier C.2 lint allowlists those branches via the
// per-file allowlist; the patterns are NOT generic-stage hardcoded paths.

package main

import (
	"sort"
	"time"
)

// emitApplySaturation: secondary replication apply-side bottleneck.
// Fires when the apply-batch-latency rate is sustained-high.
//
// Metrics:
//   - serverStatus.metrics.repl.apply.batches.num (counter, batches applied)
//   - serverStatus.metrics.repl.apply.batches.totalMillis (counter, ms in apply)
//
// Heuristic: if `totalMillis_rate / num_rate` (avg ms per batch) is > 100 ms
// sustained over the capture, the secondary's applyOps is the bottleneck.
// Only emits on replica-set captures (gated by the existence of these
// counters).
func emitApplySaturation(c *findingCollector, host string, store *seriesStore) (emitted, skipped int) {
	numID, ok1 := store.nameToID["serverStatus.metrics.repl.apply.batches.num"]
	totID, ok2 := store.nameToID["serverStatus.metrics.repl.apply.batches.totalMillis"]
	if !ok1 || !ok2 {
		return 0, 1
	}
	numCol := store.columns[numID]
	totCol := store.columns[totID]
	if numCol == nil || totCol == nil || numCol.Len() < 60 {
		return 0, 1
	}
	numVals := numCol.ReadAllInt64(nil)
	totVals := totCol.ReadAllInt64(nil)
	if len(numVals) != len(totVals) || len(numVals) < 2 {
		return 0, 1
	}
	// Pre-fix used first-to-last delta which is vulnerable to a single
	// slow batch inflating the average across an otherwise idle capture.
	// Walk in 60-sample (≈ 1 minute at 1Hz) windows; emit only when the
	// MAJORITY of windows exceed the threshold (sustained-high signal).
	const (
		applyLatencyThresholdMs = 100.0
		windowSamples           = 60
		majorityFraction        = 0.5
	)
	var windowsTotal, windowsAboveThreshold int
	for lo := 0; lo+windowSamples <= len(numVals); lo += windowSamples {
		hi := lo + windowSamples
		dn := numVals[hi-1] - numVals[lo]
		dt := totVals[hi-1] - totVals[lo]
		// Count all windows in the denominator so idle windows don't
		// inflate fracAbove. An idle window (dn <= 0) contributes 0 to
		// windowsAboveThreshold but 1 to windowsTotal.
		windowsTotal++
		if dn <= 0 {
			continue
		}
		ms := float64(dt) / float64(dn)
		if ms >= applyLatencyThresholdMs {
			windowsAboveThreshold++
		}
	}
	globalDeltaNum := numVals[len(numVals)-1] - numVals[0]
	globalDeltaTot := totVals[len(totVals)-1] - totVals[0]
	if windowsTotal == 0 {
		return 0, 1
	}
	fracAbove := float64(windowsAboveThreshold) / float64(windowsTotal)
	if fracAbove < majorityFraction {
		return 0, 1
	}
	avgMsPerBatch := 0.0
	if globalDeltaNum > 0 {
		avgMsPerBatch = float64(globalDeltaTot) / float64(globalDeltaNum)
	}
	c.add(newFinding(KindApplySaturation, host, "global", map[string]interface{}{
		"metric":                  "serverStatus.metrics.repl.apply.batches.totalMillis",
		"denom_metric":            "serverStatus.metrics.repl.apply.batches.num",
		"avg_ms_per_batch":        roundFloat(avgMsPerBatch, 2),
		"threshold_ms_per_batch":  applyLatencyThresholdMs,
		"window_samples":          windowSamples,
		"windows_above_threshold": windowsAboveThreshold,
		"windows_total":           windowsTotal,
		"fraction_above":          roundFloat(fracAbove, 3),
		"batches_applied":         globalDeltaNum,
		"total_apply_ms":          globalDeltaTot,
	}))
	return 1, 0
}

// emitElectionStorm: overlay on existing term_change Findings. Clusters
// ≥3 term changes within 5 minutes.
func emitElectionStorm(c *findingCollector, host string) (emitted, skipped int) {
	var termTimes []time.Time
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		if k != KindTermChange {
			continue
		}
		atStr, _ := f["at"].(string)
		t, err := time.Parse(tsFormat, atStr)
		if err != nil {
			continue
		}
		termTimes = append(termTimes, t)
	}
	if len(termTimes) < 3 {
		return 0, 1
	}
	// Sort term-change timestamps before the sliding-window scan.
	// term_change Findings arrive in collector order (kind-grouped, then
	// stableSortKey order), which is NOT guaranteed to be time-sorted.
	// Pre-fix, an out-of-time-order collector could miss real storms.
	// SliceStable for byte-equality on tied timestamps (rare, but possible
	// when two term_changes record the same millisecond).
	sort.SliceStable(termTimes, func(i, j int) bool { return termTimes[i].Before(termTimes[j]) })
	// Emit ONE Finding per NON-OVERLAPPING 5-min window with ≥3 term
	// changes. Pre-fix returned after the first window, silently dropping
	// later storms. Non-overlapping: after emitting at index i, jump
	// past i+2 so we don't double-count.
	for i := 0; i+2 < len(termTimes); {
		if termTimes[i+2].Sub(termTimes[i]) <= 5*time.Minute {
			// Greedy: extend the window to capture all term_changes
			// within 5 minutes of termTimes[i].
			j := i + 2
			for j+1 < len(termTimes) && termTimes[j+1].Sub(termTimes[i]) <= 5*time.Minute {
				j++
			}
			c.add(newFinding(KindElectionStorm, host, termTimes[i].UTC().Format(tsFormat), map[string]interface{}{
				"window_start":      termTimes[i].UTC().Format(tsFormat),
				"window_end":        termTimes[j].UTC().Format(tsFormat),
				"term_change_count": j - i + 1,
				"reason_code":       ReasonCodeMultipleTermChanges,
			}))
			emitted++
			i = j + 1
		} else {
			i++
		}
	}
	if emitted == 0 {
		return 0, 1
	}
	return emitted, 0
}

// emitCheckpointStall: WT checkpoint anomaly detector — two independent paths.
//
// Path A (primary): generation-counter stall — fires when the checkpoint
// generation counter stops advancing for > 2× syncdelay. This catches the
// common PALI/disagg failure mode where checkpoints stop *starting* (the
// previous detector only checked completed-checkpoint duration, which is fine
// when checkpoints run but is 0 when they never start). Uses:
//   - serverStatus.wiredTiger.checkpoint.generation (preferred)
//   - serverStatus.wiredTiger.checkpoint.total succeed number of checkpoints (fallback)
//
// Path B (secondary): duration threshold — fires when the most-recent
// completed checkpoint exceeds 80% of syncdelay. Uses:
//   - serverStatus.wiredTiger.checkpoint.most recent time (msecs)
//
// NOTE: the previous implementation used the non-existent path
//   serverStatus.wiredTiger.transaction.transaction checkpoint most recent time (msecs)
// which caused this stage to always return skipped=1. The correct paths are
// under serverStatus.wiredTiger.checkpoint.* (verified against live FTDC data).
func emitCheckpointStall(c *findingCollector, host string, store *seriesStore) (emitted, skipped int) {
	syncDelaySec := int64(60)
	if sid, ok := store.nameToID["getCmdLineOpts.parsed.storage.syncPeriodSecs"]; ok {
		if sCol := store.columns[sid]; sCol != nil && sCol.Len() > 0 {
			sVals := sCol.ReadAllInt64(nil)
			if len(sVals) > 0 && sVals[0] > 0 {
				syncDelaySec = sVals[0]
			}
		}
	}

	// Path A: generation counter stall.
	genMetric := ""
	var genID int
	var genOK bool
	for _, candidate := range []string{
		"serverStatus.wiredTiger.checkpoint.generation",
		"serverStatus.wiredTiger.checkpoint.total succeed number of checkpoints",
	} {
		if id, ok := store.nameToID[candidate]; ok {
			genMetric, genID, genOK = candidate, id, true
			break
		}
	}
	if genOK {
		genCol := store.columns[genID]
		if genCol != nil && genCol.Len() >= 2 {
			genVals := genCol.ReadAllInt64(nil)
			// Find the last sample where the counter advanced.
			lastAdvanceIdx := 0
			for i := 1; i < len(genVals); i++ {
				if genVals[i] > genVals[i-1] {
					lastAdvanceIdx = i
				}
			}
			// Stall duration = time from last advance to end of capture.
			if lastAdvanceIdx < len(store.timestamps)-1 && lastAdvanceIdx < len(genVals) {
				stalledDur := store.timestamps[len(store.timestamps)-1].Sub(store.timestamps[lastAdvanceIdx])
				stallThreshold := time.Duration(syncDelaySec*2) * time.Second
				if stalledDur >= stallThreshold {
					atTS := store.timestamps[lastAdvanceIdx].UTC().Format(tsFormat)
					c.add(newFinding(KindCheckpointStall, host, "global", map[string]interface{}{
						"metric":            genMetric,
						"at":                atTS,
						"sample_index":      lastAdvanceIdx,
						"last_generation":   genVals[lastAdvanceIdx],
						"stalled_duration_s": int64(stalledDur.Seconds()),
						"syncdelay_s":       syncDelaySec,
						"reason_code":       "generation_counter_stalled",
					}))
					emitted++
				}
			}
		}
	}

	// Path B: completed-checkpoint duration threshold.
	const durMetric = "serverStatus.wiredTiger.checkpoint.most recent time (msecs)"
	if durID, ok := store.nameToID[durMetric]; ok {
		durCol := store.columns[durID]
		if durCol != nil && durCol.Len() >= 1 {
			durVals := durCol.ReadAllInt64(nil)
			thresholdMs := int64(float64(syncDelaySec) * 1000 * 0.8)
			var maxDur int64
			var maxIdx int
			for i, v := range durVals {
				if v > maxDur {
					maxDur = v
					maxIdx = i
				}
			}
			if maxDur >= thresholdMs {
				atTS := ""
				if maxIdx < len(store.timestamps) {
					atTS = store.timestamps[maxIdx].UTC().Format(tsFormat)
				}
				c.add(newFinding(KindCheckpointStall, host, "global", map[string]interface{}{
					"metric":          durMetric,
					"at":              atTS,
					"sample_index":    maxIdx,
					"max_duration_ms": maxDur,
					"threshold_ms":    thresholdMs,
					"syncdelay_s":     syncDelaySec,
					"reason_code":     ReasonCodeCheckpointExceededSyncdelay,
				}))
				emitted++
			}
		}
	}

	if emitted == 0 {
		if !genOK {
			return 0, 1
		}
		return 0, 1
	}
	return emitted, 0
}

// heartbeatCounter describes a monotonic counter that represents ongoing work
// in a subsystem. If it stops advancing for longer than StallSec, the
// subsystem has frozen.
type heartbeatCounter struct {
	metric   string
	subsystem string // human label used in Finding fields
	stallSec int64  // seconds of flat-line that constitutes a stall
	// corroborator: optional metric that must still be advancing to confirm
	// the local process is alive (rules out the entire node being down).
	corroborator string
}

// emitHeartbeatCounterStall: generalized counter-stall detector.
//
// Scans for flat regions (stall segments) in subsystem heartbeat counters —
// metrics where advancing means "this subsystem is alive" and stalling means
// it has frozen. Emits one counter_stall Finding per longest stall segment per
// metric (above the per-metric threshold).
//
// Key difference from the simpler "find last advancement" approach used in
// emitCheckpointStall Path A: this scanner finds ANY flat region in the
// capture, including mid-capture stalls that recovered (e.g. via process
// restart). The previous "tail-only" approach missed the PALI producer freeze
// at 2026-05-20T21:11:50Z because PALI recovered after the 10:45 restart.
//
// Stall segment algorithm:
//  1. Walk samples left-to-right.
//  2. When vals[i] > vals[i-1] the counter advanced — close any open stall.
//  3. When vals[i] <= vals[i-1] (flat or reset) — open/continue a stall.
//  4. At end of series, close any open stall.
//  5. Emit the LONGEST stall whose duration exceeds stallSec AND whose
//     corroborator metric was advancing during the stall window.
//
// Corroboration rule: the corroborator must have advanced between stallStart
// and the minimum of (stallEnd, stallStart + 5 min). This confirms the local
// process was alive during the stall — ruling out a whole-node outage.
func emitHeartbeatCounterStall(c *findingCollector, host string, store *seriesStore) (emitted, skipped int) {
	type heartbeatCounter struct {
		metric       string
		subsystem    string
		stallSec     int64
		corroborator string
	}
	heartbeats := []heartbeatCounter{
		{
			metric:    "serverStatus.metrics.disagg.phylog.generatedMillis",
			subsystem: "pali_log_producer",
			stallSec:  10,
			// opcounters.query is independent of PALI; oplogLSN comes from
			// logServerManagerMetrics which freezes WITH PALI, so it is not
			// a valid corroborator for PALI-specific stalls.
			corroborator: "serverStatus.opcounters.query",
		},
		{
			metric:       "serverStatus.disagg.paliRateLimiter.successfulAdmissions",
			subsystem:    "pali_rate_limiter",
			stallSec:     10,
			corroborator: "serverStatus.opcounters.query",
		},
		{
			metric:       "serverStatus.metrics.repl.apply.batches.num",
			subsystem:    "repl_apply",
			stallSec:     120,
			corroborator: "serverStatus.connections.current",
		},
	}

	anyPresent := false
	for _, hb := range heartbeats {
		id, ok := store.nameToID[hb.metric]
		if !ok {
			continue
		}
		anyPresent = true
		col := store.columns[id]
		if col == nil || col.Len() < 2 {
			continue
		}
		vals := col.ReadAllInt64(nil)
		ts := store.timestamps

		// Find all stall segments (maximal flat or decreasing runs).
		type seg struct{ startIdx, endIdx int }
		var rawSegs []seg
		stallStart := -1
		for i := 1; i < len(vals); i++ {
			if vals[i] > vals[i-1] {
				if stallStart >= 0 {
					rawSegs = append(rawSegs, seg{stallStart, i - 1})
					stallStart = -1
				}
			} else {
				if stallStart < 0 {
					stallStart = i - 1
				}
			}
		}
		if stallStart >= 0 {
			rawSegs = append(rawSegs, seg{stallStart, len(vals) - 1})
		}

		// Phase 1: keep only segments that individually exceed stallSec.
		// Normal counter update cadence produces many short flat periods (e.g.,
		// generatedMillis updates every ~5s, leaving 4-5s flat gaps). These are
		// never stalls — they're measurement artifacts. Filter them out first.
		var longSegs []seg
		for _, s := range rawSegs {
			if s.startIdx < len(ts) && s.endIdx < len(ts) {
				dur := ts[s.endIdx].Sub(ts[s.startIdx])
				if int64(dur.Seconds()) >= hb.stallSec {
					longSegs = append(longSegs, s)
				}
			}
		}

		// Phase 2: merge adjacent qualifying segments whose inter-segment gap
		// is shorter than stallSec. A 1-2 sample "recovery" (e.g., one FTDC
		// chunk reference document updating while the subsystem is still frozen)
		// is not a meaningful recovery and should not split one long stall into
		// two shorter ones.
		segs := longSegs[:0:len(longSegs)]
		for _, s := range longSegs {
			if len(segs) > 0 {
				last := &segs[len(segs)-1]
				if last.endIdx < len(ts) && s.startIdx < len(ts) {
					gap := ts[s.startIdx].Sub(ts[last.endIdx])
					if int64(gap.Seconds()) <= hb.stallSec {
						last.endIdx = s.endIdx
						continue
					}
				}
			}
			segs = append(segs, s)
		}

		// Find the longest merged segment above the stall threshold with corroboration.
		var cVals []int64
		if hb.corroborator != "" {
			if cid, cok := store.nameToID[hb.corroborator]; cok {
				if cCol := store.columns[cid]; cCol != nil && cCol.Len() > 0 {
					cVals = cCol.ReadAllInt64(nil)
				}
			}
		}

		bestDur := time.Duration(0)
		bestStartIdx := -1
		for _, s := range segs {
			if s.startIdx >= len(ts) || s.endIdx >= len(ts) {
				continue
			}
			dur := ts[s.endIdx].Sub(ts[s.startIdx])
			if int64(dur.Seconds()) < hb.stallSec {
				continue
			}
			// Corroboration: corroborator must advance within the stall window.
			if len(cVals) > 0 {
				// Check up to 5 minutes into the stall for corroborator advance.
				windowEnd := s.endIdx
				fiveMin := ts[s.startIdx].Add(5 * time.Minute)
				for wi := s.startIdx; wi <= s.endIdx && wi < len(ts); wi++ {
					if ts[wi].After(fiveMin) {
						windowEnd = wi
						break
					}
				}
				baseVal := int64(0)
				if s.startIdx < len(cVals) {
					baseVal = cVals[s.startIdx]
				}
				corroborated := false
				for ci := s.startIdx + 1; ci <= windowEnd && ci < len(cVals); ci++ {
					if cVals[ci] > baseVal {
						corroborated = true
						break
					}
				}
				if !corroborated {
					continue
				}
			}
			if dur > bestDur {
				bestDur = dur
				bestStartIdx = s.startIdx
			}
		}

		if bestStartIdx < 0 {
			continue
		}
		atTS := ts[bestStartIdx].UTC().Format(tsFormat)
		fields := map[string]interface{}{
			"metric":             hb.metric,
			"subsystem":         hb.subsystem,
			"at":                atTS,
			"sample_index":      bestStartIdx,
			"last_value":        vals[bestStartIdx],
			"stalled_duration_s": int64(bestDur.Seconds()),
			"stall_threshold_s": hb.stallSec,
		}
		// For PALI log-producer stalls, check whether the replication waiter
		// queue grew during the stall. This is the downstream effect of PALI
		// being frozen: operations queue up waiting for replication to complete.
		// The waiter metric doesn't get a changepoint Finding (it's a linear
		// ramp from zero, not a step change), so we attach it here as evidence.
		if hb.subsystem == "pali_log_producer" {
			const waiterMetric = "serverStatus.metrics.disagg.waiters.replication"
			if wid, wok := store.nameToID[waiterMetric]; wok {
				if wCol := store.columns[wid]; wCol != nil && wCol.Len() > bestStartIdx+1 {
					wVals := wCol.ReadAllInt64(nil)
					if bestStartIdx < len(wVals) {
						endIdx := len(wVals) - 1
						peakWaiters := int64(0)
						for _, v := range wVals[bestStartIdx:] {
							if v > peakWaiters {
								peakWaiters = v
							}
						}
						baseWaiters := wVals[bestStartIdx]
						if peakWaiters > baseWaiters {
							fields["waiter_buildup_metric"] = waiterMetric
							fields["waiter_base"] = baseWaiters
							fields["waiter_peak"] = peakWaiters
							fields["waiter_end"] = wVals[endIdx]
						}
					}
				}
			}
		}
		c.add(newFinding(KindCounterStall, host, hb.metric, fields))
		emitted++
	}

	if !anyPresent {
		return 0, 1
	}
	return emitted, 0
}

// emitEvictionDivergence: WT cache fill rate vs eviction rate divergence.
// Fires when bytes_read_into / bytes_written_from_for_eviction is sustained
// > 2.0, indicating eviction can't keep up with cache fill.
func emitEvictionDivergence(c *findingCollector, host string, store *seriesStore) (emitted, skipped int) {
	rID, okR := store.nameToID["serverStatus.wiredTiger.cache.bytes read into cache"]
	wID, okW := store.nameToID["serverStatus.wiredTiger.cache.bytes written from cache"]
	if !okR || !okW {
		return 0, 1
	}
	rCol := store.columns[rID]
	wCol := store.columns[wID]
	if rCol == nil || wCol == nil || rCol.Len() < 2 || wCol.Len() < 2 {
		return 0, 1
	}
	rVals := rCol.ReadAllInt64(nil)
	wVals := wCol.ReadAllInt64(nil)
	deltaR := rVals[len(rVals)-1] - rVals[0]
	deltaW := wVals[len(wVals)-1] - wVals[0]
	if deltaW <= 0 {
		return 0, 1
	}
	ratio := float64(deltaR) / float64(deltaW)
	const evictionRatioThreshold = 2.0
	if ratio < evictionRatioThreshold {
		return 0, 1
	}
	c.add(newFinding(KindEvictionDivergence, host, "global", map[string]interface{}{
		"metric":          "serverStatus.wiredTiger.cache.bytes read into cache",
		"denom_metric":    "serverStatus.wiredTiger.cache.bytes written from cache",
		"ratio":           roundFloat(ratio, 4),
		"threshold":       evictionRatioThreshold,
		"bytes_read_into": deltaR,
		"bytes_evicted":   deltaW,
		"reason_code":     ReasonCodeReadIntoExceedsEvicted,
	}))
	return 1, 0
}

// emitCaptureQualityScore (G.3): composite per-run 0-100 score. Single
// number summarizing capture health; analyzer can fast-fail on bad inputs.
// 100 = pristine; lower = quality concerns present.
//
// Inputs (each contributes -X points if present):
//   - cv_failure_pct (run_health metric): -1 per percentage point
//   - degenerate variance metrics: -1 per 10 such Findings, capped at -20
//   - suspicious_smoothness Findings: -2 per Finding, capped at -10
//   - missing_signal Findings: -3 per Finding, capped at -20
//   - gap Findings: -2 per Finding, capped at -15
//   - decode_error Findings: -10 per Finding, capped at -30
func emitCaptureQualityScore(c *findingCollector, host string) {
	score := 100.0
	count := map[string]int{}
	degenerateCount := 0
	var cpTotal, cpFailed int
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		count[k]++
		if k == KindSuspiciousSmoothness {
			if rc, _ := f["reason_code"].(string); rc == ReasonCodeConstantZero || rc == ReasonCodeVarianceTooLow {
				degenerateCount++
			}
		}
		if k == KindChangepoint {
			// Only count Findings that carry cross_validated; react-mode
			// changepoints omit this field and must not skew the ratio.
			if validated, ok := f["cross_validated"].(bool); ok {
				cpTotal++
				if !validated {
					cpFailed++
				}
			}
		}
	}
	var cvFailurePct float64
	if cpTotal > 0 {
		cvFailurePct = 100.0 * float64(cpFailed) / float64(cpTotal)
	}
	deductionSmoothness := minInt(count[KindSuspiciousSmoothness]*2, 10)
	deductionMissing := minInt(count[KindMissingSignal]*3, 20)
	deductionGap := minInt(count[KindGap]*2, 15)
	deductionDecode := minInt(count[KindChunkDecodeError]*10, 30)
	deductionCVFailure := minInt(int(cvFailurePct), 20)
	// Extra penalty for degenerate_variance on top of suspicious_smoothness:
	// a metric that never varies is structurally useless, so the double-count
	// is intentional policy. -1 per 10 Findings, capped at 20.
	deductionDegenerate := minInt(degenerateCount/10, 20)

	score -= float64(deductionSmoothness + deductionMissing + deductionGap + deductionDecode + deductionCVFailure + deductionDegenerate)
	if score < 0 {
		score = 0
	}
	c.add(newFinding(KindCaptureQualityScore, host, "global", map[string]interface{}{
		"score":                    roundFloat(score, 1),
		"deduction_smoothness":     deductionSmoothness,
		"deduction_missing":        deductionMissing,
		"deduction_gap":            deductionGap,
		"deduction_decode":         deductionDecode,
		"deduction_cv_failure":     deductionCVFailure,
		"deduction_degenerate":     deductionDegenerate,
		"suspicious_smoothness":    count[KindSuspiciousSmoothness],
		"missing_signal":           count[KindMissingSignal],
		"gap":                      count[KindGap],
		"chunk_decode_error":       count[KindChunkDecodeError],
		"cv_failure_pct":           roundFloat(cvFailurePct, 2),
		"degenerate_variance_count": degenerateCount,
		"interpretation":           "0-100 composite; lower = quality concerns. Analyzer may fast-fail at low scores.",
	}))
}
