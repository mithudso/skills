// analyze_compare.go — P2.3 capture-vs-capture comparison.
//
// --auto-baseline-file compares against rolling EWMA from one or more
// past captures. The simpler, more common question is "this run vs
// THAT run" — two specific captures. This stage builds two seriesStores,
// drops the column bulk on both (keeping only the per-metric Welford
// aggregates fullMean/fullSD), and emits comparative Findings ranked
// by effect size.
//
// CLI: --auto-compare OTHER_PATH
//   ftdc_parser ./capA --auto --auto-compare ./capB
//
// Output:
//   kind:"capture_diff"          — one per metric satisfying threshold,
//                                  sorted by |cohens_d| descending
//   kind:"capture_diff_summary"  — one per run, with top-K metric paths
//                                  and a single line summary
//
// Bias-control: Findings carry direction (capA_higher / capB_higher) but
// NEVER any "good/bad/regressed/improved" framing — the parser doesn't
// know which capture is which (it just received two paths).
package main

import (
	"math"
	"sort"
	"time"
)

// compareMaxFindings caps the per-comparison emission so very-noisy
// runs don't flood the output. Top-K by |cohens_d| descending.
const compareMaxFindings = 200

// compareMinCohensD is the threshold for emitting a capture_diff. 0.5
// is the "medium effect" floor from BIAS_CONTROL §B.5.
const compareMinCohensD = 0.5

// runCaptureCompare is the entry point invoked from main when --auto-compare
// is set. It builds a second seriesStore for the comparison capture and
// emits diff Findings against the current store. The caller has already
// invoked the normal --auto pipeline on the primary store, so we only
// add the comparative Findings here.
func runCaptureCompare(c *findingCollector, primaryHost string, primary *seriesStore, secondaryPath string, jobs int, hasFrom, hasTo bool, fromTime, toTime time.Time) (err error) {
	secondary, _, err := buildSeriesStoreStreaming(secondaryPath, nil, jobs, hasFrom, hasTo, fromTime, toTime)
	if err != nil {
		return err
	}
	// Drop the bulk on the secondary store immediately — we only need
	// the aggregates for the comparison.
	secondary.columns = nil
	secondary.timestamps = nil
	secondary.materializedSeries = nil

	emitCaptureDiff(c, primaryHost, primary, secondary, secondaryPath)
	return nil
}

// emitCaptureDiff produces capture_diff Findings + a single capture_diff_summary
// for the (primary, secondary) pair. Comparison is over the full-capture
// Welford aggregates fullMean / fullSD, which survive the column null.
func emitCaptureDiff(c *findingCollector, host string, primary, secondary *seriesStore, secondaryPath string) {
	type diff struct {
		metric  string
		d       float64
		meanA   float64
		meanB   float64
		sdA     float64
		sdB     float64
	}

	// Intersection of metric names; sort for determinism.
	names := make([]string, 0, len(primary.nameToID))
	for n := range primary.nameToID {
		if _, ok := secondary.nameToID[n]; ok {
			names = append(names, n)
		}
	}
	sort.Strings(names)

	var diffs []diff
	for _, name := range names {
		idA := primary.nameToID[name]
		idB := secondary.nameToID[name]
		if primary.isDegenerate[idA] && secondary.isDegenerate[idB] {
			continue
		}
		meanA := primary.fullMean[idA]
		meanB := secondary.fullMean[idB]
		sdA := primary.fullSD[idA]
		sdB := secondary.fullSD[idB]
		if sdA == 0 && sdB == 0 {
			continue
		}
		nA := fullNOrDefault(primary, idA)
		nB := fullNOrDefault(secondary, idB)
		pooledSD := pooledStandardDeviation(sdA, nA, sdB, nB)
		if pooledSD <= 0 {
			continue
		}
		d := (meanA - meanB) / pooledSD
		if math.IsNaN(d) || math.IsInf(d, 0) {
			continue
		}
		if math.Abs(d) < compareMinCohensD {
			continue
		}
		diffs = append(diffs, diff{metric: name, d: d, meanA: meanA, meanB: meanB, sdA: sdA, sdB: sdB})
	}
	sort.Slice(diffs, func(i, j int) bool {
		di := math.Abs(diffs[i].d)
		dj := math.Abs(diffs[j].d)
		if di != dj {
			return di > dj
		}
		return diffs[i].metric < diffs[j].metric
	})
	if len(diffs) > compareMaxFindings {
		diffs = diffs[:compareMaxFindings]
	}

	for _, x := range diffs {
		direction := "primary_higher"
		if x.d < 0 {
			direction = "secondary_higher"
		}
		c.add(newFinding("capture_diff", host, x.metric, map[string]interface{}{
			"metric":           x.metric,
			"cohens_d":         roundFloat(x.d, 4),
			"direction":        direction,
			"primary_mean":     roundFloat(x.meanA, 6),
			"secondary_mean":   roundFloat(x.meanB, 6),
			"primary_sd":       roundFloat(x.sdA, 6),
			"secondary_sd":     roundFloat(x.sdB, 6),
			"secondary_path":   secondaryPath,
		}))
	}

	// Summary.
	topMetrics := make([]string, 0, 10)
	for _, x := range diffs {
		if len(topMetrics) >= 10 {
			break
		}
		topMetrics = append(topMetrics, x.metric)
	}
	c.add(newFinding("capture_diff_summary", host, "global", map[string]interface{}{
		"compared_against_path": secondaryPath,
		"metrics_compared":      len(names),
		"diffs_emitted":         len(diffs),
		"top_diverging_metrics": topMetrics,
		"min_cohens_d":          compareMinCohensD,
		"interpretation":        "Two specific captures compared on full-capture per-metric Welford aggregates. Directions are mechanical (primary_higher / secondary_higher); the parser does not know which capture is 'good' or 'bad'.",
	}))
}
