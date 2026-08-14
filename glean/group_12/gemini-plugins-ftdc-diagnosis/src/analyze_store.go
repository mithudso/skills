// analyze_store.go - Shared metric-index store used by every --auto stage.
//
// The performance problem this solves: Go map operations on Sample.Metrics
// are ~50ns each. Without coordination, every stage independently iterates
// samples and looks up metric values by name, doing billions of map ops on
// real fixtures (19800 samples × 6500 metrics × multiple passes per stage
// = ~3B map ops, ~150s).
//
// seriesStore is built ONCE at the top of runAuto. It does a single pass
// over all samples to:
//   - Discover all metric names and assign each a stable integer ID.
//   - Classify each metric (counter/gauge/degenerate) from its time series.
//   - Compute full-series mean and stddev via streaming Welford.
//   - Record observed min/max for the saturation-deviation component.
//
// Subsequent stages address metrics by ID (array indexing, ~0.5ns each)
// instead of by name (map indexing, ~50ns each). The total per-stage cost
// drops from "N × M × P map ops" to "N × M map ops + N × M × P array ops".
//
// Optional materialization: for the top-K metrics by Stage-1 score, we
// also build flat []float64 series so Stage 3 (Mann-Kendall, ED-PELT) and
// variance-shift CUSUM can iterate without map lookups at all.
package main

import (
	"sort"
	"sync"
	"time"
)

// seriesStore is the shared per-(--auto invocation) cache of per-metric
// metadata + indices. Built once, read by every stage.
//
// Two construction paths populate this struct:
//   1. buildSeriesStore(samples []Sample) — legacy row-oriented path,
//      kept for callers that still hold a []Sample slice.
//   2. buildSeriesStoreStreaming(path, matcher, ...) — columnar path used
//      by --auto. Builds the per-metric CompressedColumn series during
//      chunk decode, never materializing a row-oriented []Sample slice.
type seriesStore struct {
	// Indexed by metric ID. names[i] is the metric name; nameToID[names[i]] = i.
	nameToID map[string]int
	names    []string

	// Per-metric classification.
	isCounter    []bool // true ⇒ monotonically non-decreasing across full series
	isDegenerate []bool // min == max across the full series (constant counter or gauge)

	// Per-metric observed extents.
	observedMin []int64
	observedMax []int64

	// Per-metric Welford accumulators over the FULL series (counter →
	// first-difference, gauge → raw value). Used by Stage 2 correlation
	// and any stage needing centred values.
	fullMean []float64
	fullSD   []float64
	fullN    []int // Welford N used for fullSD; needed for pooled-SD in cross-capture compare

	// Optional materialized float64 series for the top-K metrics. Index
	// is the metric ID; nil entries mean "not materialized". Stage 3 and
	// variance_shift use these to avoid per-metric sample scans.
	materializedSeries [][]float64

	// Streaming-columnar path only: per-metric compressed column +
	// per-sample timestamps. Populated by buildSeriesStoreStreaming;
	// nil/empty when the legacy []Sample path is used.
	columns    []*CompressedColumn
	timestamps []time.Time
}

// buildSeriesStore runs the one-time classification pass. Two phases:
//
//   1. Sequential name discovery + per-metric min/max + chunk-boundary
//      previous values. Sequential because IDs are assigned incrementally
//      as new names appear; doing it in parallel requires merging chunk-
//      local name maps and patching IDs everywhere. The sequential pass is
//      cheap relative to phase 2 (only does ID assignment, min/max
//      updates, and one prev-value snapshot per chunk boundary).
//
//   2. Parallel Welford accumulation: each worker scans its sample chunk
//      and updates per-worker Welford accumulators for BOTH gauge (raw)
//      and counter (rate) interpretations. Per-chunk rate computation
//      uses the chunk-boundary prev values from phase 1. After phase 2,
//      partial accumulators are merged via Chan's parallel-Welford formula.
//
// The "did this metric ever decrease?" classification is also tracked
// during phase 1.
func buildSeriesStore(samples []Sample) *seriesStore {
	s := &seriesStore{nameToID: map[string]int{}}
	if len(samples) == 0 {
		return s
	}

	workers := defaultJobs()
	if workers < 1 {
		workers = 1
	}
	if len(samples) < workers*200 {
		workers = 1
	}
	chunkSize := (len(samples) + workers - 1) / workers

	// Phase 1: sequential ID + min/max + decrease detection.
	// Also snapshot the "previous value" per metric at every chunk boundary
	// so phase 2 can compute counter rates within its chunk starting from
	// the right initial value.
	type metricMeta struct {
		// Per-chunk boundary "prev value" snapshots. boundaryPrev[w] is the
		// metric's value at the last sample BEFORE chunk w begins, or
		// (0, false) if the metric hadn't appeared yet at that point.
		boundaryPrev    []int64
		boundaryHasPrev []bool
	}
	var (
		mins        []int64
		maxs        []int64
		sawDecrease []bool
		curPrev     []int64
		curHasPrev  []bool
		meta        []*metricMeta
	)
	addMetric := func(v int64, w int) int {
		s.names = append(s.names, "")
		mins = append(mins, v)
		maxs = append(maxs, v)
		sawDecrease = append(sawDecrease, false)
		curPrev = append(curPrev, 0)
		curHasPrev = append(curHasPrev, false)
		mm := &metricMeta{
			boundaryPrev:    make([]int64, workers),
			boundaryHasPrev: make([]bool, workers),
		}
		meta = append(meta, mm)
		return len(s.names) - 1
	}
	// Walk samples chunk-by-chunk so we can snapshot boundary values
	// without two passes.
	for w := 0; w < workers; w++ {
		lo := w * chunkSize
		hi := lo + chunkSize
		if hi > len(samples) {
			hi = len(samples)
		}
		if lo >= len(samples) {
			break
		}
		// Snapshot current prev state into the boundary for chunk w.
		for id := range meta {
			meta[id].boundaryPrev[w] = curPrev[id]
			meta[id].boundaryHasPrev[w] = curHasPrev[id]
		}
		for i := lo; i < hi; i++ {
			for name, v := range samples[i].Metrics {
				id, ok := s.nameToID[name]
				if !ok {
					id = addMetric(v, w)
					s.nameToID[name] = id
					s.names[id] = name
				} else {
					if v < mins[id] {
						mins[id] = v
					}
					if v > maxs[id] {
						maxs[id] = v
					}
				}
				if curHasPrev[id] && v < curPrev[id] {
					sawDecrease[id] = true
				}
				curPrev[id] = v
				curHasPrev[id] = true
			}
		}
	}
	M := len(s.names)
	s.observedMin = mins
	s.observedMax = maxs
	s.isDegenerate = make([]bool, M)
	s.isCounter = make([]bool, M)
	for id := 0; id < M; id++ {
		s.isDegenerate[id] = mins[id] == maxs[id]
		s.isCounter[id] = !sawDecrease[id] && !s.isDegenerate[id]
	}

	// Phase 2: parallel Welford on rates (counters) and raw values
	// (gauges). Each worker has its own accumulators for BOTH
	// interpretations; we pick the right one per metric after merge.
	type chunkAcc struct {
		gN, cN []int
		gMean  []float64
		cMean  []float64
		gM2    []float64
		cM2    []float64
	}
	parts := make([]chunkAcc, workers)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		lo := w * chunkSize
		hi := lo + chunkSize
		if hi > len(samples) {
			hi = len(samples)
		}
		if lo >= len(samples) {
			break
		}
		parts[w] = chunkAcc{
			gN:    make([]int, M),
			cN:    make([]int, M),
			gMean: make([]float64, M),
			cMean: make([]float64, M),
			gM2:   make([]float64, M),
			cM2:   make([]float64, M),
		}
		wg.Add(1)
		go func(w, lo, hi int) {
			defer wg.Done()
			ac := parts[w]
			prev := make([]int64, M)
			hasPrev := make([]bool, M)
			// Seed from boundary snapshots so counter rates connect across chunks.
			for id := 0; id < M; id++ {
				prev[id] = meta[id].boundaryPrev[w]
				hasPrev[id] = meta[id].boundaryHasPrev[w]
			}
			// Track which metrics appeared in this sample so we can reset
			// hasPrev for any metric that went MISSING (skips Welford
			// rate computation across a gap — otherwise the cross-gap
			// delta gets recorded as a single 1-sample rate, inflating
			// the rate distribution).
			seenThisSample := make([]bool, M)
			for i := lo; i < hi; i++ {
				for j := range seenThisSample {
					seenThisSample[j] = false
				}
				for name, v := range samples[i].Metrics {
					id := s.nameToID[name]
					seenThisSample[id] = true
					// Gauge Welford.
					welfordUpdate(&ac.gN[id], &ac.gMean[id], &ac.gM2[id], float64(v))
					// Counter Welford on rate.
					if hasPrev[id] {
						welfordUpdate(&ac.cN[id], &ac.cMean[id], &ac.cM2[id], float64(v-prev[id]))
					}
					prev[id] = v
					hasPrev[id] = true
				}
				// Reset hasPrev for metrics absent from this sample so the
				// next observation starts a fresh rate sequence.
				for id := 0; id < M; id++ {
					if !seenThisSample[id] && hasPrev[id] {
						hasPrev[id] = false
					}
				}
			}
		}(w, lo, hi)
	}
	wg.Wait()

	// Merge partial Welford via Chan's parallel formula.
	gN := make([]int, M)
	gMean := make([]float64, M)
	gM2 := make([]float64, M)
	cN := make([]int, M)
	cMean := make([]float64, M)
	cM2 := make([]float64, M)
	combine := func(nA *int, meanA *float64, m2A *float64, nB int, meanB float64, m2B float64) {
		if nB == 0 {
			return
		}
		if *nA == 0 {
			*nA = nB
			*meanA = meanB
			*m2A = m2B
			return
		}
		combinedN := *nA + nB
		delta := meanB - *meanA
		newMean := *meanA + delta*float64(nB)/float64(combinedN)
		newM2 := *m2A + m2B + delta*delta*float64(*nA)*float64(nB)/float64(combinedN)
		*nA = combinedN
		*meanA = newMean
		*m2A = newM2
	}
	for _, p := range parts {
		if p.gN == nil {
			continue
		}
		for id := 0; id < M; id++ {
			combine(&gN[id], &gMean[id], &gM2[id], p.gN[id], p.gMean[id], p.gM2[id])
			combine(&cN[id], &cMean[id], &cM2[id], p.cN[id], p.cMean[id], p.cM2[id])
		}
	}

	s.fullMean = make([]float64, M)
	s.fullSD = make([]float64, M)
	s.fullN = make([]int, M)
	for id := 0; id < M; id++ {
		if s.isCounter[id] {
			s.fullMean[id] = cMean[id]
			s.fullSD[id] = welfordSD(cN[id], cM2[id])
			s.fullN[id] = cN[id]
		} else {
			s.fullMean[id] = gMean[id]
			s.fullSD[id] = welfordSD(gN[id], gM2[id])
			s.fullN[id] = gN[id]
		}
	}
	s.materializedSeries = make([][]float64, M)
	return s
}

// materialize fills s.materializedSeries for the given metric IDs.
// Parallelized across metrics: each goroutine owns a slice of IDs and walks
// every sample to extract those IDs' values via map[string]int64 lookup.
// This is faster than the obvious sample-outer-loop because each goroutine
// only does a small set of lookups per sample (len(idsForWorker) lookups)
// instead of iterating the entire ~6500-entry Metrics map.
func (s *seriesStore) materialize(samples []Sample, ids []int) {
	if len(ids) == 0 {
		return
	}
	// Pre-allocate target series for IDs we haven't materialized yet.
	need := make([]int, 0, len(ids))
	for _, id := range ids {
		if id >= 0 && id < len(s.materializedSeries) && s.materializedSeries[id] == nil {
			s.materializedSeries[id] = make([]float64, len(samples))
			need = append(need, id)
		}
	}
	if len(need) == 0 {
		return
	}
	// Build name lookups once for the IDs we need (avoids name→ID work
	// inside the hot loop).
	names := make([]string, len(need))
	for i, id := range need {
		names[i] = s.names[id]
	}

	workers := defaultJobs()
	if workers > len(need) {
		workers = len(need)
	}
	if workers < 1 {
		workers = 1
	}
	per := (len(need) + workers - 1) / workers
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		lo := w * per
		hi := lo + per
		if hi > len(need) {
			hi = len(need)
		}
		if lo >= len(need) {
			break
		}
		wg.Add(1)
		go func(lo, hi int) {
			defer wg.Done()
			workerNames := names[lo:hi]
			workerIDs := need[lo:hi]
			prev := make([]int64, len(workerIDs))
			hasPrev := make([]bool, len(workerIDs))
			for i := range samples {
				for k, name := range workerNames {
					v, ok := samples[i].Metrics[name]
					if !ok {
						// Counter rate is undefined across a missing
						// sample — would otherwise produce an inflated
						// "rate" spanning the gap. Reset hasPrev so the
						// NEXT observation starts a fresh rate sequence.
						if s.isCounter[workerIDs[k]] {
							hasPrev[k] = false
						}
						continue
					}
					id := workerIDs[k]
					ser := s.materializedSeries[id]
					if s.isCounter[id] {
						if hasPrev[k] {
							ser[i] = float64(v - prev[k])
						}
						prev[k] = v
						hasPrev[k] = true
					} else {
						ser[i] = float64(v)
					}
				}
			}
		}(lo, hi)
	}
	wg.Wait()
}

// sortedIDsByAbsScore returns metric IDs sorted by absolute score
// descending. Used to pick top-K movers across stages.
func sortedIDsByAbsScore(scores []float64) []int {
	return sortedIDsByAbsScoreWithNames(scores, nil)
}

// sortedIDsByAbsScoreWithNames is the name-aware version. When `names`
// is non-nil and has the same length as `scores`, ties on |score| are
// broken by lexicographic name order — which is invariant across runs.
// The plain `sortedIDsByAbsScore` falls back to ID order; sufficient for
// most call sites but NOT for the trendIDs top-N cut where two highly-
// correlated metrics may share a score and ID order is non-deterministic
// (the store assigns IDs in first-encounter order, which depends on
// map iteration randomness over samples[i].Metrics).
func sortedIDsByAbsScoreWithNames(scores []float64, names []string) []int {
	ids := make([]int, len(scores))
	for i := range ids {
		ids[i] = i
	}
	hasNames := names != nil && len(names) == len(scores)
	sort.SliceStable(ids, func(i, j int) bool {
		a, b := absFloat(scores[ids[i]]), absFloat(scores[ids[j]])
		if a != b {
			return a > b
		}
		if hasNames {
			return names[ids[i]] < names[ids[j]]
		}
		return ids[i] < ids[j]
	})
	return ids
}
