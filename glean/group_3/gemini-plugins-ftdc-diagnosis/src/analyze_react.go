// analyze_react.go - ReAct hypothesis-loop CLI primitives.
//
// The analyzer skill drives targeted re-analysis by re-invoking the parser
// with one of three flags. The three primitives let the LLM compose
// follow-up queries when the initial --auto pass leaves multiple hypotheses
// at low confidence:
//
//   --auto-around-event TIMESTAMP --auto-window DURATION
//       Re-run the full --auto pipeline with samples narrowed to
//       [TIMESTAMP - window/2, TIMESTAMP + window/2]. Forces change-point
//       and correlation detection into the local segment.
//
//   --auto-segment-by METRIC
//       Detect change-points on METRIC (full resolution, no downsampling
//       cap since we run on one metric only), then emit per-segment Welch
//       stats for ALL metrics around each detected break. Answers "what
//       else shifted when this metric did?"
//
//   --auto-correlate LEADER FOLLOWER [--auto-lag-resolution MS]
//       Fine-grained sub-second xcorr between exactly two metrics. Sample-
//       rate-limited but useful when the LLM has a candidate causal pair.
package main

import (
	"io"
	"sort"
	"time"

	"pgregory.net/changepoint"
)

// runAutoAroundEvent narrows the sample window around a timestamp and runs
// the full --auto pipeline. Implemented by setting --baseline / --stress
// to halves of the window. The narrowing also makes ED-PELT's
// computational cost trivial.
func runAutoAroundEvent(w io.Writer, path string, samples []Sample, metadata []MetadataDoc, cfg AutoConfig) error {
	if !cfg.HasAroundEvent || cfg.AroundWindow <= 0 {
		return runAuto(w, path, samples, metadata, cfg)
	}
	half := cfg.AroundWindow / 2
	from := cfg.AroundEvent.Add(-half)
	to := cfg.AroundEvent.Add(half)
	// Tier A.5 fix: defensively allocate a fresh slice rather than reuse
	// the caller's backing array (samples[:0] aliasing). The in-place
	// pattern is dormant today (no caller re-reads the original samples)
	// but is a footgun for future call sites.
	filtered := make([]Sample, 0, len(samples))
	for _, s := range samples {
		if s.Timestamp.Before(from) || s.Timestamp.After(to) {
			continue
		}
		filtered = append(filtered, s)
	}
	if len(filtered) < 4 {
		c := &findingCollector{}
		c.add(newFinding(KindSkipped, hostIDFromPath(path), "around_event_empty", map[string]interface{}{
			"reason": "around-event window contains fewer than 4 samples",
			"event":  cfg.AroundEvent.UTC().Format(time.RFC3339Nano),
			"window": cfg.AroundWindow.String(),
		}))
		return c.emitAll(w)
	}
	return runAuto(w, path, filtered, metadata, cfg)
}

// runAutoSegmentBy detects change-points on the named metric and emits
// per-segment Welch stats for ALL metrics on each side of each break.
// Catches the "what else moved when X moved?" question without forcing
// the LLM to guess which metrics to inspect.
func runAutoSegmentBy(w io.Writer, path string, samples []Sample, metadata []MetadataDoc, cfg AutoConfig) error {
	if cfg.SegmentByMetric == "" {
		return runAuto(w, path, samples, metadata, cfg)
	}
	host := hostIDFromPath(path)
	c := &findingCollector{}
	emitRunMetadata(c, host, samples, metadata, cfg)

	store := buildSeriesStore(samples)
	id, ok := store.nameToID[cfg.SegmentByMetric]
	if !ok {
		c.add(newFinding(KindSkipped, host, "segment_by_unknown", map[string]interface{}{
			"reason": "segment-by metric not found in capture",
			"metric": cfg.SegmentByMetric,
		}))
		return c.emitAll(w)
	}
	store.materialize(samples, []int{id})
	series := store.materializedSeries[id]
	if series == nil || len(series) < cfg.PELTMinSegment*2 {
		c.add(newFinding(KindSkipped, host, "segment_by_short", map[string]interface{}{
			"reason": "segment-by metric series too short",
			"metric": cfg.SegmentByMetric,
		}))
		return c.emitAll(w)
	}
	// No downsampling here — running on ONE metric we can afford full
	// resolution ED-PELT.
	breaks := changepoint.NonParametric(series, cfg.PELTMinSegment)
	if len(breaks) == 0 {
		c.add(newFinding(KindSkipped, host, "segment_by_no_breaks", map[string]interface{}{
			"reason": "no change-points detected in segment-by metric",
			"metric": cfg.SegmentByMetric,
		}))
		return c.emitAll(w)
	}
	// For each break, emit a Finding describing what shifted across that
	// break (Welch stats for the top-K most-moving metrics).
	for _, br := range breaks {
		if br < cfg.PELTMinSegment/2 || br >= len(samples)-cfg.PELTMinSegment/2 {
			continue
		}
		preSamples := samples[:br]
		postSamples := samples[br:]
		_, movers, _ := scoreAndRankMetricsWithStore(host, store, preSamples, postSamples, cfg)
		// Emit the segment-anchor Finding.
		atTS := samples[br].Timestamp.UTC().Format(tsFormat)
		c.add(newFinding(KindChangepoint, host, cfg.SegmentByMetric+"|"+atTS+"|segment", map[string]interface{}{
			"metric":       cfg.SegmentByMetric,
			"at":           atTS,
			"sample_index": br,
			"segment_at":   atTS,
			"context":      "segment_by_anchor",
		}))
		// Re-emit the top movers; relabel their kind so the LLM
		// understands these are per-segment movers, not global ones.
		for _, m := range movers {
			// Tag for clarity.
			m["context"] = "segment_by:" + cfg.SegmentByMetric
			m["segment_at"] = atTS
			c.add(m)
		}
	}
	return c.emitAll(w)
}

// runAutoCorrelate runs fine-grained sub-second cross-correlation between
// exactly two metrics. Used when the LLM has identified a candidate causal
// pair from a prior --auto pass and wants to nail down the lag at higher
// precision than the default lag grid offers.
func runAutoCorrelate(w io.Writer, path string, samples []Sample, metadata []MetadataDoc, cfg AutoConfig) error {
	if cfg.CorrelateLeader == "" || cfg.CorrelateFollower == "" {
		return runAuto(w, path, samples, metadata, cfg)
	}
	host := hostIDFromPath(path)
	c := &findingCollector{}
	emitRunMetadata(c, host, samples, metadata, cfg)

	// Reject self-correlation: Pearson(X, X) ≡ 1.0 at every lag; emitting
	// it just adds noise to the output. Pre-fix the parser silently
	// emitted the trivial r=1 result.
	if cfg.CorrelateLeader == cfg.CorrelateFollower {
		c.add(newFinding(KindSkipped, host, "correlate_self", map[string]interface{}{
			"reason": "correlate leader == follower (self-correlation is trivially 1.0)",
			"metric": cfg.CorrelateLeader,
		}))
		return c.emitAll(w)
	}

	store := buildSeriesStore(samples)
	leaderID, lok := store.nameToID[cfg.CorrelateLeader]
	followerID, fok := store.nameToID[cfg.CorrelateFollower]
	if !lok || !fok {
		c.add(newFinding(KindSkipped, host, "correlate_unknown", map[string]interface{}{
			"reason":   "correlate leader or follower metric not in capture",
			"leader":   cfg.CorrelateLeader,
			"follower": cfg.CorrelateFollower,
		}))
		return c.emitAll(w)
	}
	store.materialize(samples, []int{leaderID, followerID})
	leader := store.materializedSeries[leaderID]
	follower := store.materializedSeries[followerID]
	if leader == nil || follower == nil {
		c.add(newFinding(KindSkipped, host, "correlate_no_series", map[string]interface{}{
			"reason": "could not materialize one or both series",
		}))
		return c.emitAll(w)
	}
	if store.fullSD[leaderID] <= 0 || store.fullSD[followerID] <= 0 {
		c.add(newFinding(KindSkipped, host, "correlate_degenerate", map[string]interface{}{
			"reason": "one or both series have zero variance",
		}))
		return c.emitAll(w)
	}
	// Fine-grained lag grid: from 0 to ±60 samples at 1-sample resolution
	// (sub-second when FTDC is faster than 1Hz). The current FTDC nominal
	// 1s sample interval makes "sub-second" lag impossible without
	// higher-rate FTDC capture; the flag's intent is "give me every
	// integer lag" rather than the default-pruned 8-lag grid.
	lags := []int{}
	for l := -60; l <= 60; l++ {
		lags = append(lags, l)
	}
	type lagPoint struct {
		lag int
		r   float64
		n   int
	}
	points := make([]lagPoint, 0, len(lags))
	n := len(leader)
	for _, lag := range lags {
		// count = n - |lag|, exact for a contiguous series.
		absLag := lag
		if absLag < 0 {
			absLag = -absLag
		}
		count := n - absLag
		if count < 4 {
			continue
		}
		// Pass 1: overlap-window means.
		var sumX, sumY float64
		for i := 0; i < n; i++ {
			target := i + lag
			if target < 0 || target >= n {
				continue
			}
			sumX += leader[i]
			sumY += follower[target]
		}
		meanX := sumX / float64(count)
		meanY := sumY / float64(count)
		// Pass 2: SS and centred dot product → exact Pearson r on overlap window.
		// Using full-series means/SDs would deflate r by count/(n-1) at non-zero
		// lags (e.g. r=0.5 for a perfect lag-60 pair with n=120).
		var ssX, ssY, dot float64
		for i := 0; i < n; i++ {
			target := i + lag
			if target < 0 || target >= n {
				continue
			}
			dx := leader[i] - meanX
			dy := follower[target] - meanY
			ssX += dx * dx
			ssY += dy * dy
			dot += dx * dy
		}
		denom := sqrtFloat(ssX * ssY)
		if denom <= 0 {
			continue
		}
		r := dot / denom
		if r > 1 {
			r = 1
		}
		if r < -1 {
			r = -1
		}
		points = append(points, lagPoint{lag, r, count})
	}
	if len(points) == 0 {
		c.add(newFinding(KindSkipped, host, cfg.CorrelateLeader+"|"+cfg.CorrelateFollower, map[string]interface{}{
			"reason":   "no lag had >= 4 overlapping samples",
			"leader":   cfg.CorrelateLeader,
			"follower": cfg.CorrelateFollower,
		}))
		return c.emitAll(w)
	}
	// Sort by lag ascending for the output Finding.
	sort.Slice(points, func(i, j int) bool { return points[i].lag < points[j].lag })
	bestLag := 0
	bestR := 0.0
	bestN := 0
	for _, p := range points {
		if absFloat(p.r) > absFloat(bestR) {
			bestR = p.r
			bestLag = p.lag
			bestN = p.n
		}
	}
	// Pack the full lag scan into a single Finding so the LLM can see
	// the entire profile, not just the peak.
	profile := make([]map[string]interface{}, 0, len(points))
	for _, p := range points {
		profile = append(profile, map[string]interface{}{
			"lag":     p.lag,
			"pearson": roundFloat(p.r, 4),
			"n":       p.n,
		})
	}
	c.add(newFinding(KindCorrelation, host, cfg.CorrelateLeader+"|"+cfg.CorrelateFollower, map[string]interface{}{
		"leader":       cfg.CorrelateLeader,
		"follower":     cfg.CorrelateFollower,
		"best_lag":     bestLag,
		"best_pearson": roundFloat(bestR, 4),
		"n_samples":    bestN,
		"profile":      profile,
		"resolution":   "fine-grained (sample-by-sample, lags -60..+60)",
	}))
	return c.emitAll(w)
}
