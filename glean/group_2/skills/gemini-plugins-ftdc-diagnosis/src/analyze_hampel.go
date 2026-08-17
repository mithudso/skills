// Tier B.8: Hampel filter for local spike detection. The existing
// `detectSpikes` uses global median+MAD; on autocorrelated series with
// a real regime shift, the global MAD inflates and absorbs single-sample
// extrema. Hampel uses a rolling window so that local outliers stand out
// against their immediate neighborhood. Emits `kind:"local_spike"` with
// `algorithm:"hampel"` to distinguish from the global-MAD `kind:"spike"`.
//
// Per CHANGELOG Tier-3 deferred item resolution: Hampel adds value on
// autocorrelated series where the existing global-MAD detector loses
// sensitivity. Cost: O(N · log W) for window W ≈ 50-100.

package main

import (
	"sort"
)

const (
	hampelWindow = 50 // samples on each side of the index under test
)

// emitHampelLocalSpikes walks the materialized series for each metric
// and emits a Finding when |x[i] - rolling_median(i)| / (1.4826 * rolling_mad(i))
// exceeds cfg.SpikeZ (falls back to defaultSpikeZ=5.0). Only the strongest
// hit per metric is emitted (analogous to detectSpikes' behavior) to bound
// output volume.
//
// Edge samples (first/last hampelWindow indices) are excluded from
// detection because they don't have a full ±W neighborhood. Series
// shorter than 2W+10 samples (~110 at W=50) are silently skipped.
func emitHampelLocalSpikes(c *findingCollector, host string, store *seriesStore, ids []int, cfg AutoConfig) (emitted, skipped int) {
	threshold := cfg.SpikeZ
	if threshold <= 0 {
		threshold = defaultSpikeZ
	}
	if store == nil {
		return 0, 0
	}
	for _, id := range ids {
		series := store.materializedSeries[id]
		if series == nil || len(series) < 2*hampelWindow+10 {
			skipped++
			continue
		}
		name := store.names[id]
		kind := "gauge"
		if store.isCounter[id] {
			kind = "counter"
		}
		maxZ := 0.0
		maxIdx := -1
		// Rolling window via sliding sort. For N=90K and W=100, this is
		// ~4M comparisons — acceptable. A heap-based two-heap median would
		// be faster but the constant factor is small at this window size.
		// Both buf and diffs are pre-allocated outside the loop to avoid
		// ~17M heap allocations per metric on a full 48h capture.
		buf := make([]float64, 2*hampelWindow+1)
		diffs := make([]float64, 2*hampelWindow)
		for i := hampelWindow; i < len(series)-hampelWindow; i++ {
			// Build window excluding the center sample so the spike
			// candidate doesn't pull the median toward itself.
			j := 0
			for k := i - hampelWindow; k <= i+hampelWindow; k++ {
				if k == i {
					continue
				}
				buf[j] = series[k]
				j++
			}
			window := buf[:j]
			sort.Float64s(window)
			med := medianSorted(window)
			// MAD over the same window.
			for k, v := range window {
				diffs[k] = absFloat(v - med)
			}
			sort.Float64s(diffs[:j])
			mad := medianSorted(diffs[:j])
			if mad <= 0 {
				continue
			}
			z := absFloat(series[i]-med) / (mad * 1.4826)
			if z > maxZ {
				maxZ = z
				maxIdx = i
			}
		}
		if maxZ < threshold || maxIdx < 0 {
			continue
		}
		atTS := ""
		if maxIdx < len(store.timestamps) {
			atTS = store.timestamps[maxIdx].UTC().Format(tsFormat)
		}
		c.add(newFinding("local_spike", host, name+"|"+atTS, map[string]interface{}{
			"metric":       name,
			"metric_kind":  kind,
			"algorithm":    "hampel",
			"at":           atTS,
			"sample_index": maxIdx,
			"value":        series[maxIdx],
			"z_score":      roundFloat(maxZ, 4),
			"window":       2*hampelWindow + 1,
		}))
		emitted++
	}
	return emitted, skipped
}

// medianSorted returns the median of an already-sorted slice. Behavior on
// empty slice: returns 0 (caller checks emptiness).
func medianSorted(s []float64) float64 {
	if len(s) == 0 {
		return 0
	}
	n := len(s)
	if n%2 == 1 {
		return s[n/2]
	}
	return (s[n/2-1] + s[n/2]) / 2.0
}
