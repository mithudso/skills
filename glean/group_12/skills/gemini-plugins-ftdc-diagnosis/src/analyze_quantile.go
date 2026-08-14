// analyze_quantile.go — P3.1 Theil-Sen robust slope on tail-latency series.
//
// monotonic_trend uses Spearman + OLS slope. OLS is biased upward by
// occasional p99 spikes; latency series are EXACTLY where you want a
// robust slope estimator. Theil-Sen takes the median of all pairwise
// slopes — robust to up to ~29.3% outliers, no Gaussian assumption.
//
// This stage picks out the tail-latency-shaped metrics (opLatencies.*
// and queues.execution.*) and emits `kind:"quantile_trend"` Findings
// reporting:
//   - slope (units per sample-interval)
//   - intercept
//   - n (sample count used)
//   - sign confidence (% of pairwise slopes with matching sign — a
//     non-parametric significance proxy)
//
// Cost is O(n^2) per metric for the pairwise slopes; we cap n at
// quantileMaxN by uniform decimation so the worst case stays bounded.
package main

import (
	"math"
	"sort"
	"strings"
)

// quantileMaxN caps the per-metric series length passed into Theil-Sen.
// 800 samples × n*(n-1)/2 ≈ 320K slopes per metric — bounded and fast.
const quantileMaxN = 800

// quantileMinN sets the minimum series length below which we don't
// attempt a trend estimate. Theil-Sen needs at least 10 points to be
// reasonably stable; the slope-sign confidence is undefined at n<4.
const quantileMinN = 30

// quantileMetricPrefixes is the closed set of metric path prefixes the
// stage scans. Limiting to opLatencies + queues.execution keeps the
// quadratic cost bounded while covering the high-value tail-latency
// signals.
var quantileMetricPrefixes = []string{
	"serverStatus.opLatencies.",
	"serverStatus.queues.execution.",
	// PALI tail-latency families — high-value for Gregory's domain.
	"serverStatus.metrics.disagg.phylog.committedToMaterializedLagMillis",
	"serverStatus.metrics.disagg.phylog.generatedToCommittedLagMillis",
}

// quantileMetricLatencyContains: an additional substring filter applied
// on top of the prefix gate. We want LATENCY-shaped metrics, not raw
// counters (which would happily trend up just by virtue of being
// counters). Empty string disables the substring requirement.
var quantileMetricLatencyContains = []string{
	"latency",
	"out",
	"LagMillis",
}

// emitQuantileTrend scans the store for tail-latency-shaped metrics
// and emits a quantile_trend Finding per metric where a robust slope
// could be computed. Pure function of the store; no globals.
func emitQuantileTrend(c *findingCollector, host string, store *seriesStore) (emitted int, skipped int) {
	N := store.nSamples()
	if N < quantileMinN {
		return 0, 1
	}

	// Snapshot metric names in a deterministic order so the emission
	// order is stable regardless of map iteration.
	names := make([]string, 0, len(store.nameToID))
	for n := range store.nameToID {
		if !quantileMetricMatches(n) {
			continue
		}
		names = append(names, n)
	}
	sort.Strings(names)

	for _, name := range names {
		id, ok := store.nameToID[name]
		if !ok {
			continue
		}
		if store.isDegenerate[id] {
			continue
		}
		col := store.readColumnInt64(id, nil)
		if len(col) < quantileMinN {
			continue
		}

		// Decimate to quantileMaxN if needed — uniform sampling
		// preserves trend shape, just lowers temporal resolution.
		series := col
		if len(series) > quantileMaxN {
			stride := len(series) / quantileMaxN
			if stride < 1 {
				stride = 1
			}
			dec := make([]int64, 0, quantileMaxN)
			for i := 0; i < len(series); i += stride {
				dec = append(dec, series[i])
				if len(dec) >= quantileMaxN {
					break
				}
			}
			series = dec
		}

		slope, intercept, signFraction := theilSenSlope(series)
		if math.IsNaN(slope) || math.IsInf(slope, 0) {
			continue
		}
		// Direction follows slope sign. signFraction is the fraction of
		// pairwise slopes whose sign matches the median's — a Kendall-tau-ish
		// non-parametric confidence proxy. ≥0.66 is "robustly directional";
		// <0.55 is essentially noise even when the median is nonzero.
		direction := "up"
		if slope < 0 {
			direction = "down"
		}
		idKey := name + "|quantile_trend"
		c.add(newFinding("quantile_trend", host, idKey, map[string]interface{}{
			"metric":             name,
			"slope_per_sample":   roundFloat(slope, 6),
			"intercept":          roundFloat(intercept, 6),
			"n":                  len(series),
			"direction":          direction,
			"sign_fraction":      roundFloat(signFraction, 4),
			"estimator":          "theil_sen",
			"outlier_robust":     true,
		}))
		emitted++
	}
	if emitted == 0 {
		return 0, 1
	}
	return emitted, 0
}

// quantileMetricMatches returns true if the given metric path is in
// scope for the quantile-trend stage.
func quantileMetricMatches(name string) bool {
	prefixOK := false
	for _, p := range quantileMetricPrefixes {
		if strings.HasPrefix(name, p) {
			prefixOK = true
			break
		}
	}
	if !prefixOK {
		return false
	}
	if len(quantileMetricLatencyContains) == 0 {
		return true
	}
	for _, sub := range quantileMetricLatencyContains {
		if sub == "" || strings.Contains(name, sub) {
			return true
		}
	}
	return false
}

// theilSenSlope returns (median_slope, intercept, sign_fraction) for the
// series y_i indexed by i. The "intercept" is the median of (y_i -
// median_slope * i), the Theil-Sen-recommended companion estimator.
// sign_fraction is the fraction of pairwise slopes whose sign matches
// the returned median slope's sign (between 0.5 and 1.0; higher is
// more robustly directional).
func theilSenSlope(series []int64) (slope, intercept, signFraction float64) {
	n := len(series)
	if n < 2 {
		return math.NaN(), math.NaN(), math.NaN()
	}
	pairCount := n * (n - 1) / 2
	slopes := make([]float64, 0, pairCount)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dy := float64(series[j] - series[i])
			dx := float64(j - i)
			slopes = append(slopes, dy/dx)
		}
	}
	sort.Float64s(slopes)
	slope = medianSorted(slopes)
	if slope == 0 {
		signFraction = 0.5
	} else {
		matching := 0
		for _, s := range slopes {
			if (slope > 0 && s > 0) || (slope < 0 && s < 0) {
				matching++
			}
		}
		signFraction = float64(matching) / float64(len(slopes))
	}
	// Intercept: median(y_i - slope * i).
	resid := make([]float64, n)
	for i, y := range series {
		resid[i] = float64(y) - slope*float64(i)
	}
	sort.Float64s(resid)
	intercept = medianSorted(resid)
	return slope, intercept, signFraction
}
