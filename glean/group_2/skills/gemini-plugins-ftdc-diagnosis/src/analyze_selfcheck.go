// analyze_selfcheck.go - Self-observability layer.
//
// Two responsibilities:
//
//   1. Variance-shift CUSUM: for each top mover by Stage-1 score, compute the
//      rolling-median deviation series |x - rolling_median(x)| and run CUSUM
//      on it. Emits kind:"variance_shift" Findings. This is V1's substitute
//      for the E-divisive change-point detector (deferred to V2 per agent
//      research: E-divisive fabricates spurious breakpoints on monotonic
//      drift, which is the dominant FTDC shape). CUSUM-on-deviation catches
//      the variance-shift case E-divisive was meant to cover at zero
//      portability cost.
//
//   2. End-of-run kind:"run_health" Finding: per-Finding cross-validation
//      outcomes plus stage timings and resource peak. The LLM downstream
//      cannot distinguish "real signal" from "algorithm freaking out" without
//      this metadata, so we surface the analyzer's own self-doubt.
//
// Other per-Finding cross-validators (mover subsample recompute, correlation
// split-half) are deferred to a future enhancement; the changepoint CUSUM
// cross-validator runs inline in Stage 3 and covers the highest-risk class.
package main

import (
	"sort"
)

// detectVarianceShiftsWithStore is the store-aware variant. Iterates the
// pre-materialized top-K series in the seriesStore rather than re-scanning
// samples per metric.
func detectVarianceShiftsWithStore(host string, store *seriesStore, topIDs []int, cfg AutoConfig) []Finding {
	if len(store.timestamps) < cfg.PELTMinSegment*3 {
		return nil
	}
	windowSize := cfg.PELTMinSegment
	if windowSize < 30 {
		windowSize = 30
	}
	out := []Finding{}
	for _, id := range topIDs {
		series := store.materializedSeries[id]
		if series == nil || len(series) < windowSize*3 {
			continue
		}
		dev := rollingMedianDeviation(series, windowSize)
		if len(dev) < windowSize*2 {
			continue
		}
		overallMean := 0.0
		for _, v := range dev {
			overallMean += v
		}
		overallMean /= float64(len(dev))
		cum := 0.0
		bestIdx := 0
		bestAbs := 0.0
		for i, v := range dev {
			cum += v - overallMean
			if absFloat(cum) > bestAbs {
				bestAbs = absFloat(cum)
				bestIdx = i
			}
		}
		_, devSD := meanAndSDFloat(dev)
		if devSD <= 0 || bestAbs < 3.0*devSD*sqrtFloat(float64(len(dev))) {
			continue
		}
		sampleIdx := bestIdx
		if sampleIdx >= len(store.timestamps) {
			sampleIdx = len(store.timestamps) - 1
		}
		name := store.names[id]
		kind := "gauge"
		if store.isCounter[id] {
			kind = "counter"
		}
		atTS := store.timestamps[sampleIdx].UTC().Format(tsFormat)
		out = append(out, newFinding(KindVarianceShift, host, name+"|"+atTS, map[string]interface{}{
			"metric":       name,
			"metric_kind":  kind,
			"at":           atTS,
			"sample_index": sampleIdx,
			"cusum_peak":   roundFloat(bestAbs, 4),
			"deviation_sd": roundFloat(devSD, 6),
			"window_size":  windowSize,
		}))
	}
	return out
}

// rollingMedianDeviation returns |x[i] - median(x[i-w/2 : i+w/2])| for each
// i, where w is windowSize. Edges use the available window. Uses an O(N·w·log w)
// implementation; at w=60 the inner sort operates entirely in L1 cache, making
// the heap alternative's asymptotic win irrelevant in practice.
func rollingMedianDeviation(x []float64, windowSize int) []float64 {
	n := len(x)
	out := make([]float64, n)
	half := windowSize / 2
	if half < 1 {
		half = 1
	}
	buf := make([]float64, 0, windowSize)
	for i := 0; i < n; i++ {
		lo := i - half
		hi := i + half + 1
		if lo < 0 {
			lo = 0
		}
		if hi > n {
			hi = n
		}
		buf = buf[:0]
		buf = append(buf, x[lo:hi]...)
		sort.Float64s(buf)
		var med float64
		if len(buf)%2 == 1 {
			med = buf[len(buf)/2]
		} else {
			med = (buf[len(buf)/2-1] + buf[len(buf)/2]) / 2
		}
		out[i] = absFloat(x[i] - med)
	}
	return out
}

// countCVFailures counts changepoint Findings that carry a cross_validated
// field, returning (total, failed). React-mode changepoints that omit
// cross_validated are excluded from both counts so they don't skew the ratio.
func countCVFailures(findings []Finding) (total, failed int) {
	for _, f := range findings {
		if k, _ := f["kind"].(string); k == KindChangepoint {
			if validated, ok := f["cross_validated"].(bool); ok {
				total++
				if !validated {
					failed++
				}
			}
		}
	}
	return
}

// buildRunHealth produces the end-of-run kind:"run_health" Finding. Counts
// Findings by kind plus an aggregate cv_failure_pct (% of changepoint
// Findings whose CUSUM cross-validator disagreed with ED-PELT). Stage
// timings would require timing hooks woven through the pipeline; v1 omits
// them and the field will be empty.
func buildRunHealth(host string, allFindings []Finding, sampleCount, metadataCount int) Finding {
	countsByKind := map[string]int{}
	var cpTotal, cpFailed int
	for _, f := range allFindings {
		k, ok := f["kind"].(string)
		if !ok {
			continue
		}
		countsByKind[k]++
		if k == KindChangepoint {
			if validated, ok := f["cross_validated"].(bool); ok {
				cpTotal++
				if !validated {
					cpFailed++
				}
			}
		}
	}
	var cvFailPct float64
	if cpTotal > 0 {
		cvFailPct = 100.0 * float64(cpFailed) / float64(cpTotal)
	}

	fields := map[string]interface{}{
		"findings_emitted_by_kind": countsByKind,
		"sample_count":             sampleCount,
		"metadata_count":           metadataCount,
		"changepoint_count":        cpTotal,
		"changepoint_cv_failures":  cpFailed,
		"cv_failure_pct":           roundFloat(cvFailPct, 2),
	}
	return newFinding(KindRunHealth, host, "global", fields)
}
