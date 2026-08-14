// analyze_columnar.go - Column-store-native variants of the analyze
// stages. Each *Store / *Col function consumes a *seriesStore built by
// buildSeriesStoreStreaming (columns + timestamps already populated) and
// avoids any reference to a row-oriented []Sample slice.
//
// The legacy samples-iterating functions in analyze.go and the other
// analyze_*.go files are still callable for react modes that hold a
// pre-filtered samples slice (they go through buildSeriesStoreFromSamples
// → runAutoWithHostStore so they share this column path internally).

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

// buildSeriesStoreFromSamples is the back-compat shim used by react modes
// that already hold a filtered []Sample slice. It runs the streaming
// builder over an in-memory slice and returns the populated seriesStore.
func buildSeriesStoreFromSamples(samples []Sample) *seriesStore {
	b := newStreamingSeriesStore()
	for i := range samples {
		b.appendSample(&samples[i])
	}
	b.finalize()
	return b.store
}

// emitRunMetadataStore is the column-store-native run_metadata emitter.
func emitRunMetadataStore(c *findingCollector, host string, store *seriesStore, metadata []MetadataDoc, cfg AutoConfig) {
	N := store.nSamples()
	var startTS, endTS string
	if N > 0 && len(store.timestamps) > 0 {
		startTS = store.timestamps[0].UTC().Format(time.RFC3339Nano)
		endTS = store.timestamps[len(store.timestamps)-1].UTC().Format(time.RFC3339Nano)
	}
	fields := map[string]interface{}{
		"sample_count":      N,
		"metadata_count":    len(metadata),
		"window_start":      startTS,
		"window_end":        endTS,
		"baseline_fraction": defaultBaselineFraction,
		"stress_fraction":   defaultStressFraction,
		"lag_grid_seconds":  []int{0, 1, 5, 10, 30, 60, 120, 300},
		"library_versions": map[string]string{
			"changepoint":    "pgregory.net/changepoint@v1.0.0",
			"parser_version": ParserVersion,
		},
		"histogram_recognizers": histogramRecognizerNames(),
		"mode":                  cfg.Mode,
		"thresholds": map[string]interface{}{
			"top_k":                           cfg.TopK,
			"drift_ratio":                     cfg.DriftRatio,
			"gap_threshold_sec":               cfg.GapThresholdSec,
			"correlation_min_r":               cfg.CorrelationMinR,
			"correlation_per_follower_cap":    30,
			"pelt_min_segment":                cfg.PELTMinSegment,
			"pelt_downsample_cap":             4000,
			"mk_tau_min":                      cfg.MKTauMin,
			"mk_p_max":                        cfg.MKPMax,
			"slope_t_threshold":               2.0,
			"spike_z_threshold":               cfg.SpikeZ,
			"variance_shift_sigma_multiplier": 3.0,
			"saturation_weight":               4.0,
			// Mirror P0.5 audit-pass lifts from emitRunMetadata so both
			// metadata code paths echo the same algorithmic choices.
			"capacity_pressure_default_run_seconds":  10,
			"capacity_pressure_default_pct_at":       0.05,
			"capacity_pressure_event_radius_min_sec": 5,
			"capacity_pressure_max_synthetic_events": 50,
			"changepoint_max_breaks_per_metric":      25,
			"changepoint_refinement_min_tolerance":   5,
			"apply_saturation_latency_threshold_ms":  100.0,
			"apply_saturation_window_samples":        60,
			"apply_saturation_majority_fraction":     0.5,
			"eviction_divergence_ratio_threshold":    2.0,
			"baseline_ttl_days":                      30,
			"baseline_regression_sigma":              3.0,
			"periodicity_alpha":                      0.001,
			"glossary_top_per_kind":                  3,
			"selfcheck_min_window_size":              30,
			"fdr_alpha":                              0.05,
		},
		"topk_tiers": map[string]int{
			"score_mover_each_direction": cfg.TopK,
			"trend":                      200,
			"changepoint":                30,
			"variance_shift":             30,
			"spike":                      200,
		},
		"core_metrics_always_analyzed": coreMetricsAlwaysAnalyzed,
	}
	if N >= 4 && !cfg.HasBaseline && !cfg.HasStress {
		bCount := int(float64(N) * defaultBaselineFraction)
		if bCount < 1 {
			bCount = 1
		}
		sStart := N - int(float64(N)*defaultStressFraction)
		if sStart >= N {
			sStart = N - 1
		}
		if sStart < bCount {
			sStart = bCount
		}
		fields["effective_baseline"] = map[string]interface{}{
			"sample_idx_lo": 0,
			"sample_idx_hi": bCount,
			"from":          store.timestamps[0].UTC().Format(time.RFC3339Nano),
			"to":            store.timestamps[bCount-1].UTC().Format(time.RFC3339Nano),
		}
		fields["effective_stress"] = map[string]interface{}{
			"sample_idx_lo": sStart,
			"sample_idx_hi": N,
			"from":          store.timestamps[sStart].UTC().Format(time.RFC3339Nano),
			"to":            store.timestamps[N-1].UTC().Format(time.RFC3339Nano),
		}
	}
	if cfg.HasBaseline {
		fields["explicit_baseline"] = map[string]string{
			"from": cfg.BaselineFrom.UTC().Format(time.RFC3339Nano),
			"to":   cfg.BaselineTo.UTC().Format(time.RFC3339Nano),
		}
	}
	if cfg.HasStress {
		fields["explicit_stress"] = map[string]string{
			"from": cfg.StressFrom.UTC().Format(time.RFC3339Nano),
			"to":   cfg.StressTo.UTC().Format(time.RFC3339Nano),
		}
	}
	if len(cfg.FollowerMetrics) > 0 {
		fields["user_follower_metrics"] = cfg.FollowerMetrics
	}
	// P0.3: input provenance (mirror of emitRunMetadata).
	if cfg.InputFingerprint != nil {
		fields["input_sha256"] = cfg.InputFingerprint.TopSHA256
		if cfg.InputFingerprint.IsDirectory && len(cfg.InputFingerprint.PerFileSHA256) > 0 {
			fields["input_files"] = cfg.InputFingerprint.PerFileSHA256
		}
	}
	if cfg.InputPath != "" {
		fields["input_path_basename"] = filepathBase(cfg.InputPath)
	}
	if len(cfg.CmdLine) > 0 {
		fields["cmd_line"] = cfg.CmdLine
	}
	c.add(newFinding(KindRunMetadata, host, "global", fields))
}

// emitGapFindingsStore detects FTDC sample gaps from store.timestamps.
func emitGapFindingsStore(c *findingCollector, host string, store *seriesStore, cfg AutoConfig) {
	if len(store.timestamps) < 2 {
		return
	}
	threshold := time.Duration(cfg.GapThresholdSec * float64(time.Second))
	if threshold <= 0 {
		threshold = 5 * time.Second
	}
	ts := store.timestamps
	for i := 1; i < len(ts); i++ {
		gap := ts[i].Sub(ts[i-1])
		if gap > threshold {
			at := ts[i-1].UTC().Format(tsFormat)
			c.add(newFinding(KindGap, host, at, map[string]interface{}{
				"at":               at,
				"end_at":           ts[i].UTC().Format(tsFormat),
				"duration_seconds": roundFloat(gap.Seconds(), 3),
				"prior_sample_idx": i - 1,
				"next_sample_idx":  i,
			}))
		}
	}
}

// emitSchemaFingerprintStore builds the (name, kind) hash from store fields.
// No sample iteration needed — the store already has every metric's min/max
// and counter classification.
func emitSchemaFingerprintStore(c *findingCollector, host string, store *seriesStore, metadata []MetadataDoc) {
	if len(store.names) == 0 {
		return
	}
	sortedNames := append([]string(nil), store.names...)
	sort.Strings(sortedNames)

	h := sha256.New()
	for _, name := range sortedNames {
		id := store.nameToID[name]
		kind := "gauge"
		if store.isDegenerate[id] {
			kind = "constant"
		} else if store.isCounter[id] {
			kind = "counter"
		}
		fmt.Fprintf(h, "%s|%s\n", name, kind)
	}
	sum := h.Sum(nil)
	fp := hex.EncodeToString(sum[:8])

	c.add(newFinding(KindSchemaFingerprint, host, "global", map[string]interface{}{
		"fingerprint":        fp,
		"metric_count":       len(sortedNames),
		"metadata_doc_count": len(metadata),
	}))
}

// selectWindowsStore picks baseline and stress sample index ranges. Honours
// --baseline / --stress overrides when supplied (fixing the analyze.go:608
// no-op stub).
func selectWindowsStore(store *seriesStore, cfg AutoConfig) (bStart, bEnd, sStart, sEnd int, windowFinding Finding) {
	n := store.nSamples()
	if n < 4 {
		return 0, 0, 0, 0, nil
	}
	// Honor explicit windows when provided.
	if cfg.HasBaseline {
		bStart, bEnd = timeRangeToSampleIdx(store.timestamps, cfg.BaselineFrom, cfg.BaselineTo)
	} else {
		bCount := int(float64(n) * defaultBaselineFraction)
		if bCount < 1 {
			bCount = 1
		}
		bStart, bEnd = 0, bCount
	}
	if cfg.HasStress {
		sStart, sEnd = timeRangeToSampleIdx(store.timestamps, cfg.StressFrom, cfg.StressTo)
	} else {
		sStartI := n - int(float64(n)*defaultStressFraction)
		if sStartI >= n {
			sStartI = n - 1
		}
		if sStartI < bEnd {
			sStartI = bEnd
		}
		sStart, sEnd = sStartI, n
	}
	// Final guards.
	if bEnd > n {
		bEnd = n
	}
	if sEnd > n {
		sEnd = n
	}
	if bStart < 0 {
		bStart = 0
	}
	if sStart < 0 {
		sStart = 0
	}
	return bStart, bEnd, sStart, sEnd, nil
}

// timeRangeToSampleIdx maps an [from, to] timestamp range into the closest
// sample index range [start, end) via binary search on the sorted
// timestamps slice. Used by --baseline / --stress overrides.
func timeRangeToSampleIdx(timestamps []time.Time, from, to time.Time) (int, int) {
	if len(timestamps) == 0 {
		return 0, 0
	}
	start := sort.Search(len(timestamps), func(i int) bool {
		return !timestamps[i].Before(from)
	})
	end := sort.Search(len(timestamps), func(i int) bool {
		return timestamps[i].After(to)
	})
	if start > end {
		start = end
	}
	return start, end
}

// welfordWindowCol runs Welford streaming variance over a metric's column
// in [start, end). For counters, the series is first-differences (rate);
// for gauges, raw values plus a saturation accumulator. Returns
// (n, mean, m2, satSum). Caller is responsible for parallelism across
// metrics.
//
// scratch is a reusable []int64 buffer; pass nil to let it allocate.
func welfordWindowCol(store *seriesStore, id, start, end int, scratch []int64) (
	int, float64, float64, float64,
) {
	col := store.columns[id]
	if col == nil || end <= start {
		return 0, 0, 0, 0
	}
	// For counter rates we need value at (start-1) to compute the first
	// rate at index start.
	readStart := start
	if store.isCounter[id] && start > 0 {
		readStart = start - 1
	}
	vals := col.ReadRange(readStart, end, scratch)
	var n int
	var mean, m2, satSum float64
	if store.isCounter[id] {
		// First-difference Welford over vals[1:] paired with vals[:len-1].
		for i := 1; i < len(vals); i++ {
			x := float64(vals[i] - vals[i-1])
			n++
			delta := x - mean
			mean += delta / float64(n)
			m2 += delta * (x - mean)
		}
	} else {
		omax := store.observedMax[id]
		for i := 0; i < len(vals); i++ {
			x := float64(vals[i])
			n++
			delta := x - mean
			mean += delta / float64(n)
			m2 += delta * (x - mean)
			if omax > 0 {
				satSum += x / float64(omax)
			}
		}
	}
	return n, mean, m2, satSum
}

// parallelWelfordWindowCol fans Welford across metrics in parallel and
// returns per-metric (n, mean, m2, satSum) arrays.
func parallelWelfordWindowCol(store *seriesStore, start, end int) (
	[]int, []float64, []float64, []float64,
) {
	M := len(store.names)
	n := make([]int, M)
	mean := make([]float64, M)
	m2 := make([]float64, M)
	satSum := make([]float64, M)
	if M == 0 || end <= start {
		return n, mean, m2, satSum
	}
	// Tier A.4: pin Welford worker count for deterministic reduction tree.
	workers := welfordWorkers
	if M < workers*4 {
		workers = 1
	}
	per := (M + workers - 1) / workers
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		lo := w * per
		hi := lo + per
		if hi > M {
			hi = M
		}
		if lo >= M {
			break
		}
		wg.Add(1)
		go func(lo, hi int) {
			defer wg.Done()
			scratch := make([]int64, 0, end-start+1)
			for id := lo; id < hi; id++ {
				if store.isDegenerate[id] {
					continue
				}
				ni, mi, m2i, sati := welfordWindowCol(store, id, start, end, scratch)
				n[id] = ni
				mean[id] = mi
				m2[id] = m2i
				satSum[id] = sati
				scratch = scratch[:0]
			}
		}(lo, hi)
	}
	wg.Wait()
	return n, mean, m2, satSum
}

// scoreAndRankMetricsWithStoreCol is the column-native Stage 1.
func scoreAndRankMetricsWithStoreCol(host string, store *seriesStore, bStart, bEnd, sStart, sEnd int, cfg AutoConfig) (Finding, []Finding, []float64) {
	M := len(store.names)
	if M == 0 || bEnd <= bStart || sEnd <= sStart {
		return newFinding(KindMetricScoreTable, host, "global", map[string]interface{}{
			"entries":      []interface{}{},
			"empty_reason": "no metrics or empty window",
		}), nil, nil
	}

	baseN, baseMean, baseM2, baseSatSum := parallelWelfordWindowCol(store, bStart, bEnd)
	stressN, stressMean, stressM2, stressSatSum := parallelWelfordWindowCol(store, sStart, sEnd)

	scores := make([]float64, M)
	for id := 0; id < M; id++ {
		if store.isDegenerate[id] {
			continue
		}
		bSD := welfordSD(baseN[id], baseM2[id])
		sSD := welfordSD(stressN[id], stressM2[id])
		t := welchT(baseMean[id], bSD, baseN[id], stressMean[id], sSD, stressN[id])
		var satDelta float64
		if !store.isCounter[id] && baseN[id] > 0 && stressN[id] > 0 {
			satBase := baseSatSum[id] / float64(baseN[id])
			satStress := stressSatSum[id] / float64(stressN[id])
			satDelta = satBase - satStress
		}
		combined := t
		if !store.isCounter[id] {
			combined = t + 4.0*satDelta
		}
		scores[id] = combined
	}

	type scored struct {
		id    int
		name  string
		kind  string
		score float64
	}
	scoredAll := make([]scored, 0, M)
	nDegenerate := 0
	for id := 0; id < M; id++ {
		kind := "gauge"
		if store.isDegenerate[id] {
			kind = "constant"
			nDegenerate++
		} else if store.isCounter[id] {
			kind = "counter"
		}
		scoredAll = append(scoredAll, scored{id, store.names[id], kind, scores[id]})
	}
	sort.Slice(scoredAll, func(i, j int) bool { return scoredAll[i].name < scoredAll[j].name })
	entries := make([]interface{}, 0, len(scoredAll))
	for _, s := range scoredAll {
		entries = append(entries, map[string]interface{}{
			"metric": s.name,
			"kind":   s.kind,
			"score":  roundFloat(s.score, 4),
		})
	}
	// P0.NEW-7 setup: compute FDR-BH q-values across all non-degenerate
	// metrics, then attach the summary (count surviving + p-threshold) to
	// the score-table Finding. The per-id q-value lookup `qByID` is reused
	// by emitMover below to attach q_value to each mover.
	scoreTable := newFinding(KindMetricScoreTable, host, "global", map[string]interface{}{
		"entries":      entries,
		"n_metrics":    len(scoredAll),
		"n_degenerate": nDegenerate,
	})

	movableScored := make([]scored, 0, len(scoredAll))
	for _, s := range scoredAll {
		if s.kind == "constant" {
			continue
		}
		movableScored = append(movableScored, s)
	}
	scoredAll = movableScored
	// P0.NEW-7: Compute Welch p-values for every non-degenerate metric, then
	// apply Benjamini-Hochberg FDR adjustment. The q-value lookup is keyed
	// by metric id (s.id) so emitMover can attach it. With ~6500 tests per
	// capture, BH is the right correction: Bonferroni is way too
	// conservative (single-test α/m ≈ 7.7e-6) and uncorrected p_value
	// gives the solo user a flood of false positives. Tier B.11 already
	// does Bonferroni on periodicity; this completes the multiple-testing
	// posture for movers.
	const fdrAlpha = 0.05
	pByID := make([]float64, M)
	for i := range pByID {
		pByID[i] = math.NaN()
	}
	for _, s := range scoredAll {
		id := s.id
		bSD := welfordSD(baseN[id], baseM2[id])
		sSD := welfordSD(stressN[id], stressM2[id])
		df := welchDF(bSD, baseN[id], sSD, stressN[id])
		if df <= 0 {
			continue
		}
		t := welchT(baseMean[id], bSD, baseN[id], stressMean[id], sSD, stressN[id])
		pv := clampFloat01(2 * studentTSurvivalApprox(t, df))
		pByID[id] = pv
	}
	qByID := benjaminiHochberg(pByID)
	bhSig := countSignificantBH(qByID, fdrAlpha)
	fdrThreshold := bhSignificanceThreshold(pByID, fdrAlpha)
	// Attach FDR summary to scoreTable (the per-run aggregate Finding) so
	// downstream consumers see one place to read significance accounting.
	// Note: scoreTable was constructed earlier in this function; we
	// augment its fields here rather than rebuilding. Finding is just a
	// map[string]interface{}; mutation is safe before emission.
	scoreTable["fdr_alpha"] = fdrAlpha
	// -1.0 is the bhSignificanceThreshold sentinel meaning no test passed BH correction.
	scoreTable["fdr_adjusted_threshold_p"] = roundFloat(fdrThreshold, 6)
	scoreTable["bh_significant_count"] = bhSig // int, not float64 — assert accordingly
	scoreTable["bh_total_tested"] = countNonNaN(pByID)
	// Tier A.4 fix: stable sort with explicit tiebreaker.
	sort.SliceStable(scoredAll, func(i, j int) bool {
		if scoredAll[i].score != scoredAll[j].score {
			return scoredAll[i].score > scoredAll[j].score
		}
		return scoredAll[i].name < scoredAll[j].name
	})
	movers := make([]Finding, 0, 2*cfg.TopK)
	emitMover := func(s scored, direction string) {
		id := s.id
		bSD := welfordSD(baseN[id], baseM2[id])
		sSD := welfordSD(stressN[id], stressM2[id])
		t := welchT(baseMean[id], bSD, baseN[id], stressMean[id], sSD, stressN[id])
		// Tier B.9 / B.10 / B.5: emit Welch-Satterthwaite df, calibrated
		// two-sided p_value, and Cohen's d alongside the existing welch_t.
		df := welchDF(bSD, baseN[id], sSD, stressN[id])
		var pValue float64
		if df > 0 {
			pValue = clampFloat01(2 * studentTSurvivalApprox(t, df))
		}
		var cohenD float64
		if baseN[id] >= 2 && stressN[id] >= 2 {
			pooledSD := pooledStandardDeviation(bSD, baseN[id], sSD, stressN[id])
			if pooledSD > 0 {
				cohenD = (stressMean[id] - baseMean[id]) / pooledSD
			}
		}
		var satDelta float64
		if !store.isCounter[id] && baseN[id] > 0 && stressN[id] > 0 {
			satDelta = (baseSatSum[id] / float64(baseN[id])) - (stressSatSum[id] / float64(stressN[id]))
		}
		fields := map[string]interface{}{
			"metric":      s.name,
			"metric_kind": s.kind,
			"direction":   direction,
			"score":       roundFloat(s.score, 4),
			"welch_t":     roundFloat(t, 4),
			"welch_df":    roundFloat(df, 2),
			"p_value":     roundFloat(pValue, 6),
			"cohens_d":    roundFloat(cohenD, 4),
			"baseline": map[string]interface{}{
				"mean":  roundFloat(baseMean[id], 6),
				"stdev": roundFloat(bSD, 6),
				"n":     baseN[id],
			},
			"stress": map[string]interface{}{
				"mean":  roundFloat(stressMean[id], 6),
				"stdev": roundFloat(sSD, 6),
				"n":     stressN[id],
			},
		}
		// P0.NEW-7: Benjamini-Hochberg q-value, computed once across all
		// ~6500 metrics. Trust q_value (not p_value) for significance —
		// the parser ran m tests, p_value alone is uncorrected.
		if qv := qByID[id]; !math.IsNaN(qv) {
			fields["q_value"] = roundFloat(qv, 6)
		}
		if s.kind == "gauge" {
			fields["saturation_delta"] = roundFloat(satDelta, 6)
		}
		movers = append(movers, newFinding(KindMover, host, s.name+"|"+direction, fields))
	}
	k := cfg.TopK
	if k <= 0 {
		k = defaultAutoTopK
	}
	for i := 0; i < k && i < len(scoredAll); i++ {
		emitMover(scoredAll[i], "up")
	}
	for i := 0; i < k && i < len(scoredAll); i++ {
		emitMover(scoredAll[len(scoredAll)-1-i], "down")
	}
	return scoreTable, movers, scores
}

// modalSampleIntervalSecondsStore picks the most common 1-step delta from
// store.timestamps to determine the effective sample rate.
func modalSampleIntervalSecondsStore(timestamps []time.Time) float64 {
	if len(timestamps) < 2 {
		return 1.0
	}
	counts := map[int64]int{}
	for i := 1; i < len(timestamps); i++ {
		dt := timestamps[i].Sub(timestamps[i-1]).Milliseconds()
		if dt <= 0 {
			continue
		}
		counts[dt]++
	}
	// Tier A.5 fix: iterate in sorted key order so ties break deterministically
	// (smaller dt wins). Pre-fix, Go's randomized map iteration made the modal
	// interval — and downstream lagSamples / clock-skew / periodogram outputs
	// — non-reproducible across runs on tied counts.
	dts := make([]int64, 0, len(counts))
	for dt := range counts {
		dts = append(dts, dt)
	}
	sort.Slice(dts, func(i, j int) bool { return dts[i] < dts[j] })
	bestDt := int64(1000)
	bestN := 0
	for _, dt := range dts {
		n := counts[dt]
		if n > bestN {
			bestN = n
			bestDt = dt
		}
	}
	return float64(bestDt) / 1000.0
}

// correlateAgainstFollowersWithStoreCol is the column-native Stage 2.
// Walks each metric's column directly, accumulates lagged dot products
// against follower columns (already in store.materializedSeries).
func correlateAgainstFollowersWithStoreCol(host string, store *seriesStore, cfg AutoConfig) []Finding {
	N := store.nSamples()
	if N < 8 {
		return nil
	}
	intervalSec := modalSampleIntervalSecondsStore(store.timestamps)
	if intervalSec <= 0 {
		intervalSec = 1
	}
	// Build follower vector — defaults + user extension, deduplicated.
	allFollowers := append([]string(nil), defaultFollowerMetrics...)
	allFollowers = append(allFollowers, cfg.FollowerMetrics...)
	sort.Strings(allFollowers)
	allFollowers = dedupStrings(allFollowers)

	followerIDs := []int{}
	followerNames := []string{}
	for _, fname := range allFollowers {
		id, ok := store.nameToID[fname]
		if !ok {
			continue
		}
		followerIDs = append(followerIDs, id)
		followerNames = append(followerNames, fname)
	}
	if len(followerIDs) == 0 {
		return nil
	}
	store.materializeFromColumns(followerIDs)
	followerSeries := make([][]float64, len(followerIDs))
	for i, id := range followerIDs {
		followerSeries[i] = store.materializedSeries[id]
	}
	// followerSD is an early-exit guard (fSD<=0 skips degenerate followers
	// before the lag loop); it is NOT used in the Pearson normalization —
	// the denom<=0 check in the normalization handles that case too.
	followerSD := make([]float64, len(followerIDs))
	for i, id := range followerIDs {
		followerSD[i] = store.fullSD[id]
	}

	// Lag grid in samples. Use rounded division (half-up) so a fractional
	// intervalSec doesn't truncate small lags to 0. Mirrors the fix in
	// analyze_correlation.go; pre-fix, intervalSec=2.0 collapsed the {1s}
	// entry into the {0s} entry.
	lagSamples := make([]int, 0, len(lagGridSeconds))
	seenLag := map[int]bool{}
	for _, l := range lagGridSeconds {
		var ls int
		if l == 0 {
			ls = 0
		} else {
			ls = int(float64(l)/intervalSec + 0.5)
			if ls == 0 {
				ls = 1
			}
		}
		if ls < 0 || seenLag[ls] {
			continue
		}
		lagSamples = append(lagSamples, ls)
		seenLag[ls] = true
	}
	if len(lagSamples) == 0 {
		return nil
	}

	M := len(store.names)
	F := len(followerSeries)
	L := len(lagSamples)

	// Metric-outer parallelism. Each worker takes a slice of metric IDs;
	// for each metric it decodes the column once, then accumulates
	// per-(follower, lag) dot products against the centred follower series.
	//
	// Tier A.4: pin worker count for deterministic FP summation order.
	workers := welfordWorkers
	if M < workers*4 {
		workers = 1
	}
	per := (M + workers - 1) / workers

	// Per-(metric, follower, lag) raw-sum accumulators for overlap-window
	// computational Pearson r (W7-R1 fix: full-series SDs deflate r at
	// non-zero lags; overlap-window formula eliminates the deflation).
	sz := M * F * L
	sumXBuf := make([]float64, sz)
	sumYBuf := make([]float64, sz)
	sumXYBuf := make([]float64, sz)
	sumX2Buf := make([]float64, sz)
	sumY2Buf := make([]float64, sz)
	nCounts := make([]int32, sz)

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		lo := w * per
		hi := lo + per
		if hi > M {
			hi = M
		}
		if lo >= M {
			break
		}
		wg.Add(1)
		go func(lo, hi int) {
			defer wg.Done()
			vals := make([]float64, 0, N)
			ints := make([]int64, 0, N)
			for idx := lo; idx < hi; idx++ {
				if store.isDegenerate[idx] {
					continue
				}
				col := store.columns[idx]
				if col == nil {
					continue
				}
				// Decode the metric column. For counters we need the
				// first-difference series to match what fullMean/fullSD
				// describe (which were computed on rates).
				if store.isCounter[idx] {
					ints = col.ReadAllInt64(ints[:0])
					if cap(vals) < len(ints) {
						vals = make([]float64, len(ints))
					} else {
						vals = vals[:len(ints)]
					}
					var prev int64
					hasPrev := false
					for i, v := range ints {
						if hasPrev {
							vals[i] = float64(v - prev)
						} else {
							vals[i] = 0
						}
						prev = v
						hasPrev = true
					}
				} else {
					vals = col.ReadAll(vals[:0])
				}
				for i := 0; i < len(vals); i++ {
					// For counters, sample i=0 has no defined rate, skip.
					if store.isCounter[idx] && i == 0 {
						continue
					}
					mv := vals[i]
					mvSq := mv * mv
					for fIdx := 0; fIdx < F; fIdx++ {
						fSeries := followerSeries[fIdx]
						base := (idx*F + fIdx) * L
						for lagIdx, lag := range lagSamples {
							target := i + lag
							if target < 0 || target >= N {
								continue
							}
							fv := fSeries[target]
							sumXBuf[base+lagIdx] += mv
							sumYBuf[base+lagIdx] += fv
							sumXYBuf[base+lagIdx] += mv * fv
							sumX2Buf[base+lagIdx] += mvSq
							sumY2Buf[base+lagIdx] += fv * fv
							nCounts[base+lagIdx]++
						}
					}
				}
			}
		}(lo, hi)
	}
	wg.Wait()

	type candidate struct {
		metric, follower string
		lagSec           int
		r                float64
		nUsed            int
	}
	candidates := []candidate{}
	for idx := 0; idx < M; idx++ {
		if store.isDegenerate[idx] {
			continue
		}
		mSD := store.fullSD[idx]
		if mSD <= 0 {
			continue
		}
		for fIdx := 0; fIdx < F; fIdx++ {
			if followerIDs[fIdx] == idx {
				continue // skip self-correlation
			}
			fSD := followerSD[fIdx]
			if fSD <= 0 {
				continue
			}
			for lagIdx, lag := range lagSamples {
				base := (idx*F + fIdx) * L
				nu := int(nCounts[base+lagIdx])
				if nu < 4 {
					continue
				}
				// Overlap-window computational Pearson r (W7-R1).
				fc := float64(nu)
				mx := sumXBuf[base+lagIdx] / fc
				my := sumYBuf[base+lagIdx] / fc
				ssX := sumX2Buf[base+lagIdx] - fc*mx*mx
				ssY := sumY2Buf[base+lagIdx] - fc*my*my
				dotC := sumXYBuf[base+lagIdx] - fc*mx*my
				denom := math.Sqrt(ssX * ssY)
				if denom <= 0 {
					continue
				}
				r := dotC / denom
				if r > 1 {
					r = 1
				}
				if r < -1 {
					r = -1
				}
				lagSec := int(math.Round(float64(lag) * intervalSec))
				if absFloat(r) < cfg.CorrelationMinR {
					continue
				}
				candidates = append(candidates, candidate{
					metric:   store.names[idx],
					follower: followerNames[fIdx],
					lagSec:   lagSec,
					r:        r,
					nUsed:    nu,
				})
			}
		}
	}
	// Sort by |r| desc, ties broken by (follower, metric, lag) for
	// determinism. Without tie-breaking, the per-follower cap below picks
	// arbitrarily between equally-correlated metrics, making goldens racy.
	// Tier A.4: SliceStable for byte-equality across --jobs settings.
	sort.SliceStable(candidates, func(i, j int) bool {
		ai, aj := absFloat(candidates[i].r), absFloat(candidates[j].r)
		if ai != aj {
			return ai > aj
		}
		if candidates[i].follower != candidates[j].follower {
			return candidates[i].follower < candidates[j].follower
		}
		if candidates[i].metric != candidates[j].metric {
			return candidates[i].metric < candidates[j].metric
		}
		return candidates[i].lagSec < candidates[j].lagSec
	})
	const perFollowerCap = 30
	perFollower := map[string]int{}
	out := []Finding{}
	for _, c := range candidates {
		if perFollower[c.follower] >= perFollowerCap {
			continue
		}
		perFollower[c.follower]++
		// Align field names with analyze_correlation.go's legacy emit AND
		// the schema (analyze_schema.go correlation oneOf requires
		// leader/follower/lag_seconds/pearson). Pre-fix the columnar
		// path emitted metric/follower/lag_sec/r — a schema violation
		// that --auto-self-validate would catch but slipped through.
		out = append(out, newFinding(KindCorrelation, host, c.metric+"|"+c.follower, map[string]interface{}{
			"leader":      c.metric,
			"follower":    c.follower,
			"lag_seconds": c.lagSec,
			"pearson":     roundFloat(c.r, 4),
			"n_samples":   c.nUsed,
		}))
	}
	return out
}

// inferHistogramPercentilesStore is the column-native Stage 4.
func inferHistogramPercentilesStore(host string, store *seriesStore, bStart, bEnd, sStart, sEnd int, cfg AutoConfig) []Finding {
	_ = cfg
	if bEnd-bStart < 2 || sEnd-sStart < 2 {
		return nil
	}
	keys := make(map[string]bool, len(store.names))
	for _, n := range store.names {
		keys[n] = true
	}
	families := discoverHistogramFamilies(keys)
	if len(families) == 0 {
		return nil
	}
	type famResult struct {
		baseline Finding
		stress   Finding
	}
	results := make([]famResult, len(families))
	var wg sync.WaitGroup
	for idx, fam := range families {
		idx, fam := idx, fam
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[idx].baseline = inferFamilyWindowStore(host, fam, "baseline", store, bStart, bEnd)
			results[idx].stress = inferFamilyWindowStore(host, fam, "stress", store, sStart, sEnd)
		}()
	}
	wg.Wait()
	out := make([]Finding, 0, 2*len(families))
	for _, r := range results {
		if r.baseline != nil {
			out = append(out, r.baseline)
		}
		if r.stress != nil {
			out = append(out, r.stress)
		}
	}
	return out
}

// inferFamilyWindowStore is the column-native variant of inferFamilyWindow.
// Reads each bucket's first/last value in [start, end) via column.ReadRange.
func inferFamilyWindowStore(host string, fam histogramFamily, windowName string, store *seriesStore, start, end int) Finding {
	if end-start < 2 || end > store.nSamples() {
		return nil
	}
	readFirst := func(name string) (int64, bool) {
		id, ok := store.nameToID[name]
		if !ok {
			return 0, false
		}
		v := store.columns[id].ReadRange(start, start+1, nil)
		if len(v) == 0 {
			return 0, false
		}
		return v[0], true
	}
	readLast := func(name string) (int64, bool) {
		id, ok := store.nameToID[name]
		if !ok {
			return 0, false
		}
		v := store.columns[id].ReadRange(end-1, end, nil)
		if len(v) == 0 {
			return 0, false
		}
		return v[0], true
	}
	startTS := store.timestamps[start].UTC().Format(tsFormat)
	endTS := store.timestamps[end-1].UTC().Format(tsFormat)

	switch fam.pattern {
	case "indexed":
		counts := make([]int64, len(fam.countKeys))
		bucketBounds := make([]float64, len(fam.countKeys))
		hasBound := make([]bool, len(fam.countKeys))
		invalid := false
		for i, ck := range fam.countKeys {
			f, fok := readFirst(ck)
			l, lok := readLast(ck)
			if !fok || !lok {
				continue
			}
			if l < f {
				invalid = true
				break
			}
			counts[i] = l - f
			if fam.microsKeys[i] != "" {
				if mv, ok := readLast(fam.microsKeys[i]); ok {
					bucketBounds[i] = float64(mv)
					hasBound[i] = true
				}
			}
		}
		if invalid {
			return newFinding(KindHistogramInvalid, host, fam.prefix+"|"+windowName, map[string]interface{}{
				"metric":  fam.prefix,
				"window":  windowName,
				"pattern": "indexed",
				"reason":  "bucket counter decreased (counter reset or shard movement)",
			})
		}
		var total int64
		for _, c := range counts {
			total += c
		}
		if total < 4 {
			return nil
		}
		fields := map[string]interface{}{
			"metric":       fam.prefix,
			"window":       windowName,
			"pattern":      "indexed",
			"ops":          total,
			"buckets":      len(counts),
			"window_start": startTS,
			"window_end":   endTS,
		}
		anyBound := false
		for _, h := range hasBound {
			if h {
				anyBound = true
				break
			}
		}
		if anyBound {
			bounds := make([]float64, 0, len(counts)+1)
			bounds = append(bounds, 0)
			for i := range counts {
				if hasBound[i] {
					bounds = append(bounds, bucketBounds[i])
				} else if len(bounds) > 0 {
					bounds = append(bounds, bounds[len(bounds)-1]*2)
				}
			}
			fields["p50_us"] = roundFloat(interpolatePercentile(counts, bounds, total, 0.50), 1)
			fields["p95_us"] = roundFloat(interpolatePercentile(counts, bounds, total, 0.95), 1)
			fields["p99_us"] = roundFloat(interpolatePercentile(counts, bounds, total, 0.99), 1)
			fields["bucket_upper_bounds"] = bucketBounds
		} else {
			dist := make([]int64, len(counts))
			copy(dist, counts)
			fields["bucket_counts"] = dist
		}
		return newFinding(KindHistogramPercentile, host, fam.prefix+"|"+windowName, fields)

	case "bracket":
		counts := make([]int64, len(fam.buckets))
		invalid := false
		var total int64
		for i, b := range fam.buckets {
			f, fok := readFirst(b.key)
			l, lok := readLast(b.key)
			if !fok || !lok {
				continue
			}
			if l < f {
				invalid = true
				break
			}
			counts[i] = l - f
			total += counts[i]
		}
		if invalid {
			return newFinding(KindHistogramInvalid, host, fam.prefix+"|"+windowName, map[string]interface{}{
				"metric":  fam.prefix,
				"window":  windowName,
				"pattern": "bracket",
				"reason":  "bucket counter decreased (counter reset or shard movement)",
			})
		}
		if total < 4 {
			return nil
		}
		bounds := make([]float64, 0, len(fam.buckets)+1)
		bounds = append(bounds, fam.buckets[0].lo)
		for _, b := range fam.buckets {
			bounds = append(bounds, b.hi)
		}
		p50 := interpolatePercentile(counts, bounds, total, 0.50)
		p95 := interpolatePercentile(counts, bounds, total, 0.95)
		p99 := interpolatePercentile(counts, bounds, total, 0.99)
		return newFinding(KindHistogramPercentile, host, fam.prefix+"|"+windowName, map[string]interface{}{
			"metric":       fam.prefix,
			"window":       windowName,
			"pattern":      "bracket",
			"ops":          total,
			"buckets":      len(fam.buckets),
			"p50_us":       roundFloat(p50, 4),
			"p95_us":       roundFloat(p95, 4),
			"p99_us":       roundFloat(p99, 4),
			"window_start": startTS,
			"window_end":   endTS,
		})
	}
	return nil
}

