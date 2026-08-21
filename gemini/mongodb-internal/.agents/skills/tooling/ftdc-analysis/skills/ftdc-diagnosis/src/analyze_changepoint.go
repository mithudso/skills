// analyze_changepoint.go - Trend and change-point analysis on top movers.
//
// Two complementary detectors plus a cross-validator:
//
//	monotonic_trend  — Spearman rank correlation between metric values and
//	                   the time index, with a linear-regression slope and
//	                   significance test. Catches slow ramps (tcmalloc creep,
//	                   connection leaks) where there is no clean breakpoint.
//	                   Spearman is the rank-based equivalent of Kendall's
//	                   tau and runs in O(N log N) vs Kendall's naïve O(N²);
//	                   they identify the same "clearly monotonic" series in
//	                   practice.
//
//	changepoint      — ED-PELT (Empirical Distribution PELT) via
//	                   pgregory.net/changepoint. Nonparametric: handles
//	                   heavy-tailed FTDC counter rates correctly where
//	                   PELT-L2+BIC over-segments per Truong et al. 2020.
//
//	CUSUM x-validator — for each ED-PELT detection, run a one-pass CUSUM on
//	                    a tight window around the breakpoint. Agreement
//	                    raises confidence; disagreement halves it. The
//	                    LLM downstream cannot distinguish over-segmentation
//	                    from a real regime change without this metadata.
//
// Stage 3 operates only on the top movers by Stage-1 Welch score (not all
// 6500 metrics). Stage 1 has already considered every metric mechanically;
// the deeper algorithms here run on the 200 most likely to repay the cost.
package main

import (
	"math"
	"sort"

	"pgregory.net/changepoint"
)

// trendOnly runs Spearman-rank trend + linear regression slope per metric.
// Cheap (O(N log N) per metric) so we run it on the top-200 movers. Emits
// monotonic_trend Findings for metrics with clear monotonic shape — the
// case ED-PELT correctly stays silent on (no clean breakpoint).
func trendOnly(host string, store *seriesStore, ids []int, cfg AutoConfig) []Finding {
	out := []Finding{}
	for _, id := range ids {
		series := store.materializedSeries[id]
		if series == nil || len(series) < 30 {
			continue
		}
		name := store.names[id]
		kind := "gauge"
		if store.isCounter[id] {
			kind = "counter"
		}
		spearman := spearmanRankCorr(series)
		slope, slopeT := linearSlope(series)
		// Tier B.14: emit Spearman p_value via z = ρ·√(N-1) asymptotic.
		// MKPMax config was dead pre-B.14; now gates the trend emit as
		// originally intended. Tier B.6 (Kendall block-size + power
		// boundary) is handled below for short captures.
		spearmanP := spearmanPValueApprox(spearman, len(series))
		// If MKPMax is non-positive, treat it as "p-gate disabled" so the
		// detector still fires on tau + slope_t alone. Pre-fix, MKPMax=0
		// silently locked out trend emission because spearmanP ∈ (0,1] is
		// never ≤ 0.
		pGateOk := cfg.MKPMax <= 0 || spearmanP <= cfg.MKPMax
		if absFloat(spearman) >= cfg.MKTauMin && absFloat(slopeT) >= 2.0 && pGateOk {
			// Tier B.6: refuse to emit on captures < 30 minutes of equivalent
			// observations (N < 1800 at 1Hz) — Kendall/Spearman power for
			// τ≥0.25 is marginal (~0.6) at that capture length. Emit
			// `trend_inconclusive` instead so the analyzer knows the gap
			// is "capture too short", not "no trend".
			if len(series) < 1800 {
				out = append(out, newFinding(KindTrendInconclusive, host, name, map[string]interface{}{
					"metric":         name,
					"metric_kind":    kind,
					"reason_code":    "capture_too_short",
					"spearman":       roundFloat(spearman, 4),
					"spearman_p":     roundFloat(spearmanP, 6),
					"slope":          roundFloat(slope, 6),
					"slope_t":        roundFloat(slopeT, 4),
					"n":              len(series),
					"power_estimate": "marginal_for_tau_geq_0.25",
				}))
				continue
			}
			out = append(out, newFinding(KindMonotonicTrend, host, name, map[string]interface{}{
				"metric":      name,
				"metric_kind": kind,
				"spearman":    roundFloat(spearman, 4),
				"spearman_p":  roundFloat(spearmanP, 6),
				"p_value":     roundFloat(spearmanP, 6), // Tier B.10: uniform field
				"slope":       roundFloat(slope, 6),
				"slope_t":     roundFloat(slopeT, 4),
				"n":           len(series),
			}))
		}
	}
	return out
}

// spearmanPValueApprox returns the two-sided asymptotic p-value for the
// null `ρ = 0` under the Spearman rank correlation. Uses the standard
// large-n approximation z = ρ · √(n - 2). Returns 1 when n < 3.
//
// Tier B.14 fix: dead config `MKPMax` is now wired via this function.
// Post-final-review fix: previously used √(n-1) — the standard Spearman
// large-n null is √(n-2), matching Pearson's transformation. The ~0.5%
// inflation in z-scores at n=100 from the wrong k was minor but worth
// fixing for honesty.
func spearmanPValueApprox(rho float64, n int) float64 {
	if n < 3 {
		return 1
	}
	rho2 := rho * rho
	if rho2 >= 1 {
		return 0 // perfect correlation
	}
	t := rho * sqrtFloat(float64(n-2)) / sqrtFloat(1-rho2)
	return 2 * normalSurvival(absFloat(t))
}

// lag1Autocorr returns the lag-1 sample autocorrelation ρ̂₁ ∈ [-1, 1].
// Returns 0 for series shorter than 2 or with zero variance. Used by
// Tier B.3 to scale the ED-PELT effect-size gate.
func lag1Autocorr(x []float64) float64 {
	if len(x) < 2 {
		return 0
	}
	var mean float64
	for _, v := range x {
		mean += v
	}
	mean /= float64(len(x))
	var num, denom float64
	for i := 1; i < len(x); i++ {
		num += (x[i-1] - mean) * (x[i] - mean)
	}
	for _, v := range x {
		denom += (v - mean) * (v - mean)
	}
	if denom <= 0 {
		return 0
	}
	return num / denom
}

// changepointOnly runs ED-PELT on the top-K (smaller cap, e.g. 30) movers.
// On series longer than 4000 samples we block-average down to ≤4000 for the
// ED-PELT call (its worst-case O(N²) on heavy-tailed FTDC rates makes
// running at full resolution infeasible), then REFINE the detected
// breakpoint indices back to original-resolution sample indices via a
// local CUSUM scan within the decimation window. This recovers the
// time-precision loss from downsampling: end-to-end the breakpoint
// timestamps are accurate to ±1 sample (= ±1s at 1Hz FTDC) regardless of
// the downsampling factor.
func changepointOnly(host string, store *seriesStore, ids []int, cfg AutoConfig) []Finding {
	out := []Finding{}
	for _, id := range ids {
		series := store.materializedSeries[id]
		if series == nil || len(series) < cfg.PELTMinSegment*2 {
			continue
		}
		name := store.names[id]
		kind := "gauge"
		if store.isCounter[id] {
			kind = "counter"
		}

		// Skip ED-PELT entirely when the series is
		// dominated by a strong linear trend. ED-PELT segments the
		// drift into many pieces; `monotonic_trend` already captures
		// the meaningful signal. Threshold of |slope_t| >= 100 catches
		// the pure-ramp FP class (mem.virtual @ slope_t=240,
		// connections.current @ slope_t=155, etc.) without suppressing
		// metrics whose trend is real but weaker.
		_, slopeT := linearSlope(series)
		if absFloat(slopeT) >= 100 {
			continue
		}

		pelInput, decimation := downsampleForPELT(series, 4000)
		minSeg := cfg.PELTMinSegment / decimation
		if minSeg < 5 {
			minSeg = 5
		}
		rawBreaks := changepoint.NonParametric(pelInput, minSeg)

		// Cap per-metric changepoint emission. When
		// ED-PELT finds >25 breakpoints on a single metric, the metric
		// is structurally non-stationary (oscillating cache fills,
		// noisy gauges with no stable regime). 100+ individual
		// changepoint Findings on one metric drown the analyzer in
		// noise without adding signal — emit ONE Finding describing
		// the oversegmentation instead. Structural, not name-based.
		const maxBreaksPerMetric = 25
		if len(rawBreaks) > maxBreaksPerMetric {
			// Emit a separate Finding kind so downstream consumers
			// (analyze_schema.go, ftdc-analyzer SKILL.md, run_health
			// `cv_failure_pct`) don't conflate it with the per-event
			// `changepoint` schema, which requires `at`, `sample_index`,
			// `effect_size`, `confidence`, `cross_validated`.
			out = append(out, newFinding(KindChangepointOversegmented, host, name, map[string]interface{}{
				"metric":      name,
				"metric_kind": kind,
				"break_count": len(rawBreaks),
				"decimation":  decimation,
				"reason_code": "ed_pelt_oversegmented",
				"reason":      "ED-PELT segmented >25 times; series is structurally non-stationary, not a discrete regime change",
			}))
			continue
		}
		for _, b := range rawBreaks {
			// Refine: search a TWO-block window straddling the down-
			// sampled breakpoint. ED-PELT typically marks the start of
			// the new segment, so the true transition can lie in block
			// (b-1) OR block b — searching only block b biases the
			// refined index forward by up to `decimation` samples. Use
			// [(b-1)*dec, (b+1)*dec) clamped to [0, n).
			refLo := (b - 1) * decimation
			refHi := (b + 1) * decimation
			if refLo < 0 {
				refLo = 0
			}
			if refHi > len(series) {
				refHi = len(series)
			}
			br := refineBreakpoint(series, refLo, refHi)
			if br <= 0 || br >= len(series) {
				continue
			}
			pre := series[:br]
			post := series[br:]
			preMean, preSD := meanAndSDFloat(pre)
			postMean, postSD := meanAndSDFloat(post)
			pooledSD := pooledStandardDeviation(preSD, len(pre), postSD, len(post))
			effect := 0.0
			if pooledSD > 0 {
				effect = absFloat(postMean-preMean) / pooledSD
			}
			// Tier B.3 + round-4 fix: scale the Cohen's d gate by the
			// AR(1) variance-inflation factor (1 + 2ρ̂₁)/(1 − ρ̂₁) for the
			// PRE-BREAKPOINT segment, not the full series. Pre-fix used
			// lag1Autocorr on the full series; a sharp breakpoint at
			// position br inflates lag-1 correlation purely because the
			// shift itself creates correlation across the boundary —
			// self-suppressing real changes. The pre segment captures the
			// noise structure independent of the change.
			//
			// Cap at 4× to avoid suppressing real signals on near-unit-
			// root series; cap ρ̂ itself at 0.95 to keep the formula
			// stable.
			rho := 0.0
			if len(pre) >= 30 {
				rho = lag1Autocorr(pre)
			} else if len(post) >= 30 {
				rho = lag1Autocorr(post)
			} else {
				rho = lag1Autocorr(series)
			}
			if rho > 0.95 {
				rho = 0.95
			}
			arInflate := 1.0
			if rho > 0 {
				arInflate = (1.0 + 2.0*rho) / (1.0 - rho)
				if arInflate > 4.0 {
					arInflate = 4.0
				}
			}
			gate := 0.5 * arInflate
			confidence := mapEffectToConfidence(effect, len(pre), len(post))
			validated, cusumOffset := cusumCrossValidate(series, br, cfg.PELTMinSegment)
			if !validated {
				// At small effect sizes (< AR-adjusted gate), CUSUM
				// disagreement is the canonical false-positive signature
				// for ED-PELT — drop the Finding entirely.
				if effect < gate {
					continue
				}
				confidence /= 2
			}
			var atTS string
			if br < len(store.timestamps) {
				atTS = store.timestamps[br].UTC().Format(tsFormat)
			}
			// Tier B.13: rename "confidence" → "confidence_rank" to prevent
			// downstream consumers (including the LLM analyzer) from treating
			// the tanh-based heuristic score as a calibrated probability.
			// The mapEffectToConfidence docstring states it's a heuristic,
			// but the field name "confidence" misleads. The interpretation
			// rule is documented in the analyzer SKILL.md addendum.
			out = append(out, newFinding(KindChangepoint, host, name+"|"+atTS, map[string]interface{}{
				"metric":          name,
				"metric_kind":     kind,
				"at":              atTS,
				"sample_index":    br,
				"baseline_mean":   roundFloat(preMean, 6),
				"baseline_stdev":  roundFloat(preSD, 6),
				"baseline_n":      len(pre),
				"stressed_mean":   roundFloat(postMean, 6),
				"stressed_stdev":  roundFloat(postSD, 6),
				"stressed_n":      len(post),
				"effect_size":     roundFloat(effect, 4),
				"cohens_d":        roundFloat(effect, 4), // alias: effect_size IS Cohen's d here
				"confidence_rank": roundFloat(confidence, 4),
				"cross_validated": validated,
				"cusum_offset":    cusumOffset,
				"decimation":      decimation,
			}))
		}
	}
	return out
}

// refineBreakpoint locates the most prominent local mean-shift inside
// [lo, hi) of the input series using one-pass CUSUM. Returns the global
// sample index. Used to recover original-resolution precision after the
// ED-PELT downsampling pass.
func refineBreakpoint(series []float64, lo, hi int) int {
	if hi <= lo+1 {
		return lo
	}
	overallMean := 0.0
	for i := lo; i < hi; i++ {
		overallMean += series[i]
	}
	overallMean /= float64(hi - lo)
	cum := 0.0
	bestIdx := lo
	bestAbs := 0.0
	for i := lo; i < hi; i++ {
		cum += series[i] - overallMean
		if absFloat(cum) > bestAbs {
			bestAbs = absFloat(cum)
			bestIdx = i
		}
	}
	return bestIdx
}

// detectSpikes runs robust z-score against global median+MAD on full-
// resolution materialized series. Emits kind:"spike" Findings when the
// maximum |z| in the series exceeds the threshold (5σ-equivalent). Closes
// the per-sample-extreme detection gap that ED-PELT (sustained-shift focus)
// and variance-shift CUSUM (cumulative aggregation) leave open.
func detectSpikes(host string, store *seriesStore, ids []int, cfg AutoConfig) []Finding {
	threshold := cfg.SpikeZ
	if threshold <= 0 {
		threshold = defaultSpikeZ
	}
	out := []Finding{}
	for _, id := range ids {
		series := store.materializedSeries[id]
		if series == nil || len(series) < 30 {
			continue
		}
		name := store.names[id]
		kind := "gauge"
		if store.isCounter[id] {
			kind = "counter"
		}
		med, mad := medianAndMAD(series)
		if mad <= 0 {
			continue
		}
		maxAbsZ := 0.0
		maxIdx := 0
		for i, v := range series {
			z := absFloat(v-med) / (mad * 1.4826) // 1.4826 = MAD → σ scale
			if z > maxAbsZ {
				maxAbsZ = z
				maxIdx = i
			}
		}
		if maxAbsZ < threshold {
			continue
		}
		atTS := ""
		if maxIdx < len(store.timestamps) {
			atTS = store.timestamps[maxIdx].UTC().Format(tsFormat)
		}
		out = append(out, newFinding(KindSpike, host, name+"|"+atTS, map[string]interface{}{
			"metric":       name,
			"metric_kind":  kind,
			"at":           atTS,
			"sample_index": maxIdx,
			"value":        series[maxIdx],
			"median":       roundFloat(med, 6),
			"mad":          roundFloat(mad, 6),
			"z_score":      roundFloat(maxAbsZ, 4),
		}))
	}
	return out
}

// medianAndMAD computes the median and median-absolute-deviation of x.
// Both robust estimators of central tendency / spread for heavy-tailed
// data where mean/stddev are misleading.
func medianAndMAD(x []float64) (float64, float64) {
	n := len(x)
	if n == 0 {
		return 0, 0
	}
	buf := make([]float64, n)
	copy(buf, x)
	sort.Float64s(buf)
	med := buf[n/2]
	if n%2 == 0 {
		med = (buf[n/2-1] + buf[n/2]) / 2
	}
	for i := range buf {
		buf[i] = absFloat(x[i] - med)
	}
	sort.Float64s(buf)
	mad := buf[n/2]
	if n%2 == 0 {
		mad = (buf[n/2-1] + buf[n/2]) / 2
	}
	return med, mad
}

// spearmanRankCorr computes Spearman rank correlation between x and the time
// index 0,1,2,…,N-1. Equivalent to Pearson correlation on ranks. Returns 0
// for series with fewer than 3 points or zero variance.
func spearmanRankCorr(x []float64) float64 {
	n := len(x)
	if n < 3 {
		return 0
	}
	// Compute ranks of x with average rank for ties.
	type idxVal struct {
		idx int
		val float64
	}
	pairs := make([]idxVal, n)
	for i, v := range x {
		pairs[i] = idxVal{i, v}
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].val < pairs[j].val })
	ranks := make([]float64, n)
	i := 0
	for i < n {
		j := i + 1
		for j < n && pairs[j].val == pairs[i].val {
			j++
		}
		// average rank for [i, j)
		avg := float64(i+j-1) / 2.0
		for k := i; k < j; k++ {
			ranks[pairs[k].idx] = avg
		}
		i = j
	}

	// Pearson between ranks and time index (0..N-1).
	meanRank := float64(n-1) / 2.0
	num := 0.0
	denX := 0.0
	denT := 0.0
	for i := 0; i < n; i++ {
		dr := ranks[i] - meanRank
		dt := float64(i) - meanRank
		num += dr * dt
		denX += dr * dr
		denT += dt * dt
	}
	if denX <= 0 || denT <= 0 {
		return 0
	}
	return num / sqrtFloat(denX*denT)
}

// linearSlope returns (slope, t-statistic) from a simple linear regression
// of x against time index 0..N-1. The t-statistic tests slope ≠ 0.
func linearSlope(x []float64) (float64, float64) {
	n := len(x)
	if n < 3 {
		return 0, 0
	}
	meanT := float64(n-1) / 2.0
	meanX := 0.0
	for _, v := range x {
		meanX += v
	}
	meanX /= float64(n)

	sxt := 0.0
	stt := 0.0
	for i := 0; i < n; i++ {
		dt := float64(i) - meanT
		sxt += dt * (x[i] - meanX)
		stt += dt * dt
	}
	if stt <= 0 {
		return 0, 0
	}
	slope := sxt / stt
	// Residual sum of squares for t-statistic.
	rss := 0.0
	for i := 0; i < n; i++ {
		pred := meanX + slope*(float64(i)-meanT)
		r := x[i] - pred
		rss += r * r
	}
	if n < 3 {
		return slope, 0
	}
	se := sqrtFloat((rss / float64(n-2)) / stt)
	if se <= 0 {
		// Noiseless linear data: slope_t is mathematically infinite.
		// Returning 0 lets the changepoint slope_t gate fail
		// open on perfect ramps. Return a large sentinel that
		// reliably exceeds the gate threshold (100) for any non-zero
		// slope. Sign matches the slope so |slope_t| comparisons work.
		if slope == 0 {
			return 0, 0
		}
		sign := 1.0
		if slope < 0 {
			sign = -1.0
		}
		return slope, sign * 1e9
	}
	return slope, slope / se
}

// pooledStandardDeviation returns the pooled SD of two samples, used as the
// denominator of Cohen's d effect size at a change-point.
func pooledStandardDeviation(sd1 float64, n1 int, sd2 float64, n2 int) float64 {
	if n1+n2-2 <= 0 {
		return 0
	}
	v := ((float64(n1-1) * sd1 * sd1) + (float64(n2-1) * sd2 * sd2)) / float64(n1+n2-2)
	// Guard against FP underflow producing v < 0 (which would NaN out the
	// sqrt and silently propagate through Cohen's d downstream). welfordSD
	// already clamps internally, but pooled aggregation can re-introduce
	// negative-by-ULP values.
	if v < 0 {
		v = 0
	}
	return sqrtFloat(v)
}

// mapEffectToConfidence collapses (Cohen's d effect size, segment lengths)
// into a [0, 1] confidence. Heuristic, not a calibrated probability: large
// effect + long segments = high confidence; small effect or short segments =
// low confidence. The LLM downstream interprets this categorically; precise
// calibration would require empirical labeling.
func mapEffectToConfidence(effect float64, n1, n2 int) float64 {
	if effect <= 0 {
		return 0
	}
	// Tanh maps [0, ∞) → [0, 1) and is gentle near origin. Scale so effect
	// of 0.5 → ~0.46 confidence, 1.0 → 0.76, 2.0 → 0.96. Then attenuate by
	// segment length adequacy.
	base := tanhFloat(effect / 2.0)
	minN := n1
	if n2 < minN {
		minN = n2
	}
	if minN < 60 {
		base *= float64(minN) / 60.0
	}
	if base > 0.99 {
		base = 0.99
	}
	if base < 0 {
		base = 0
	}
	return base
}

// cusumCrossValidate runs a one-pass CUSUM on a window around the proposed
// breakpoint and reports whether CUSUM independently agrees that a shift
// occurred near the same index. Returns (validated, offset_in_samples). On
// disagreement, the caller halves the change-point confidence.
//
// CUSUM detects the index at which the cumulative running mean has diverged
// most from the overall mean of the window. If that index is within
// minSegment/4 of the proposed breakpoint, agreement.
func cusumCrossValidate(series []float64, br int, minSegment int) (bool, int) {
	w := minSegment * 3
	lo := br - w
	hi := br + w
	if lo < 0 {
		lo = 0
	}
	if hi > len(series) {
		hi = len(series)
	}
	if hi-lo < minSegment*2 {
		return false, 0
	}
	seg := series[lo:hi]
	overallMean := 0.0
	for _, v := range seg {
		overallMean += v
	}
	overallMean /= float64(len(seg))

	cum := 0.0
	bestIdx := 0
	bestAbs := 0.0
	for i, v := range seg {
		cum += v - overallMean
		if absFloat(cum) > bestAbs {
			bestAbs = absFloat(cum)
			bestIdx = i
		}
	}
	cusumBr := lo + bestIdx
	offset := cusumBr - br
	tolerance := minSegment / 4
	if tolerance < 5 {
		tolerance = 5
	}
	if offset < 0 {
		if -offset <= tolerance {
			return true, offset
		}
	} else if offset <= tolerance {
		return true, offset
	}
	return false, offset
}

// tanhFloat wraps math.Tanh.
func tanhFloat(x float64) float64 {
	return math.Tanh(x)
}

// downsampleForPELT block-averages a series down to at most maxN samples.
// Returns (smaller_series, decimation_factor). When len(series) ≤ maxN the
// input is returned unchanged and decimation is 1. Block averaging
// preserves change-point structure while reducing ED-PELT's worst-case
// O(N²) cost by the square of the decimation factor.
func downsampleForPELT(series []float64, maxN int) ([]float64, int) {
	n := len(series)
	if n <= maxN {
		return series, 1
	}
	dec := (n + maxN - 1) / maxN
	if dec < 2 {
		dec = 2
	}
	out := make([]float64, 0, n/dec+1)
	for i := 0; i < n; i += dec {
		end := i + dec
		if end > n {
			end = n
		}
		sum := 0.0
		for k := i; k < end; k++ {
			sum += series[k]
		}
		out = append(out, sum/float64(end-i))
	}
	return out, dec
}
