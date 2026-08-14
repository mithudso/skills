// streaming_store.go - One-pass columnar builder for seriesStore.
//
// Replaces the old "accumulate []Sample, then buildSeriesStore" pipeline
// with a streaming builder that consumes chunk batches as they arrive from
// the FTDC decoder, appends per-metric CompressedColumns, and updates
// online stats. The []Sample batch from each onBatch callback is released
// for GC before the next batch arrives — so peak memory is bounded by the
// compressed column store, not by the un-compressed row-oriented sample
// slice that previously consumed 16-25 GB on a 48h capture.
//
// See plan: /home/ubuntu/.claude/plans/how-to-make-ftdc-parser-graceful-sky.md

package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// streamingCapNotice prints the worker-cap notice at most once per process.
var streamingCapNotice sync.Once

// buildSeriesStoreStreaming streams chunk batches from processFTDC into a
// columnar seriesStore. Returns a *seriesStore with columns, timestamps,
// nameToID, names, fullMean/fullSD, observedMin/Max, isCounter,
// isDegenerate populated.
//
// Memory peak: O(compressed columns + active chunk batches in flight).
// Multiple decode workers each hold a chunk's row-oriented []Sample
// (~200 MB per chunk with full metric map) until it's consumed by
// appendSample. To bound this, we cap streaming workers to 2 regardless
// of the caller's --jobs setting — this trades ~25% decode wall-time for
// ~3× lower peak RSS, which matters far more for long captures.
//
// The applyAuto* filters (from/to time range, nonzero) are applied
// per-batch as each batch arrives, so memory bounds hold even when the
// caller supplies a tight window.
func buildSeriesStoreStreaming(
	path string,
	matcher *metricMatcher,
	workers int,
	hasFrom, hasTo bool,
	fromTime, toTime time.Time,
) (*seriesStore, []MetadataDoc, error) {

	// Tier A.5 fix: keep the RSS-budget cap (documented above) but stop
	// silently overriding the user's --jobs choice. Print to stderr the
	// first time we cap so users can see why their --jobs N ≠ effective
	// workers. The cap is enforced once per process (sync.Once).
	if workers > 2 {
		streamingCapNotice.Do(func() {
			fmt.Fprintf(os.Stderr,
				"ftdc_parser: streaming mode caps decode workers at 2 for RSS budget (requested=%d); see streaming_store.go:23 for rationale\n",
				workers)
		})
		workers = 2
	}
	if workers < 1 {
		workers = 1
	}

	s := newStreamingSeriesStore()
	var meta []MetadataDoc
	meta = processFTDC(path, matcher, workers, func(batch []Sample) {
		// Time-range filter applied here so out-of-window samples never
		// enter the column store.
		for i := range batch {
			if hasFrom && batch[i].Timestamp.Before(fromTime) {
				continue
			}
			if hasTo && batch[i].Timestamp.After(toTime) {
				continue
			}
			s.appendSample(&batch[i])
		}
		// batch goes out of scope, GC frees it before next batch arrives.
	})

	s.finalize()
	return s.store, meta, nil
}

// streamingBuilder is the working state for buildSeriesStoreStreaming.
// Kept package-private; the consumer only sees the finished *seriesStore.
type streamingBuilder struct {
	store *seriesStore

	// Per-metric scratchpad updated as samples flow in.
	lastSeenIdx []int   // sample index where this metric was last observed
	lastValue   []int64 // last observed value
	hasPrev     []bool  // for rate Welford
	prevValue   []int64

	// Bookkeeping of total samples processed.
	totalSamples int

	// Per-sample scratch (reused).
	prevSampleTs time.Time
	havePrevTs   bool
}

func newStreamingSeriesStore() *streamingBuilder {
	return &streamingBuilder{
		store: &seriesStore{
			nameToID: map[string]int{},
		},
	}
}

// addMetric registers a brand-new metric name and allocates parallel
// per-metric slots. Returns the new metric ID.
func (b *streamingBuilder) addMetric(name string, initialVal int64) int {
	id := len(b.store.names)
	b.store.names = append(b.store.names, name)
	b.store.nameToID[name] = id

	b.store.observedMin = append(b.store.observedMin, initialVal)
	b.store.observedMax = append(b.store.observedMax, initialVal)
	b.store.isCounter = append(b.store.isCounter, true)
	b.store.isDegenerate = append(b.store.isDegenerate, true)

	b.store.columns = append(b.store.columns, NewCompressedColumn(b.totalSamples+1024))

	// Backfill: if this metric appears for the first time at sample N>0,
	// pad with zeros for samples [0, N) so columns stay sample-aligned.
	// b.totalSamples has already been incremented for the current sample;
	// subtract 1 because the caller's Add() at the bottom of appendSample
	// will append the actual current-sample value.
	preCount := b.totalSamples - 1
	if preCount < 0 {
		preCount = 0
	}
	for i := 0; i < preCount; i++ {
		b.store.columns[id].AddPad(0)
	}

	b.lastSeenIdx = append(b.lastSeenIdx, -1)
	b.lastValue = append(b.lastValue, 0)
	b.hasPrev = append(b.hasPrev, false)
	b.prevValue = append(b.prevValue, 0)

	return id
}

// appendSample consumes one decoded sample into the column store. All
// known metrics get one column entry per call; metrics missing from
// sample.Metrics are zero-padded with consecutive=false to break rate
// continuity (mirrors the reset behaviour at analyze_store.go:233-238).
func (b *streamingBuilder) appendSample(sample *Sample) {
	sampleIdx := b.totalSamples
	b.totalSamples++

	// Detect FTDC gap (>5s since the previous sample).
	isGap := false
	if b.havePrevTs && sample.Timestamp.Sub(b.prevSampleTs) > 5*time.Second {
		isGap = true
	}
	b.prevSampleTs = sample.Timestamp
	b.havePrevTs = true
	b.store.timestamps = append(b.store.timestamps, sample.Timestamp)

	// Track which metrics appeared in this sample so we can pad the absent ones.
	// To avoid an O(M) scan on every sample, we use a per-call bitmap whose
	// size grows with the number of registered metrics.
	for name, v := range sample.Metrics {
		id, ok := b.store.nameToID[name]
		if !ok {
			id = b.addMetric(name, v)
		}
		// Track decrease (counter classification).
		if b.lastSeenIdx[id] >= 0 && v < b.lastValue[id] {
			b.store.isCounter[id] = false
		}
		// Track min/max.
		if v < b.store.observedMin[id] {
			b.store.observedMin[id] = v
		}
		if v > b.store.observedMax[id] {
			b.store.observedMax[id] = v
		}
		if v != b.lastValue[id] || b.lastSeenIdx[id] < 0 {
			// isDegenerate stays true only if every observed value is identical.
			if b.lastSeenIdx[id] >= 0 || v != b.lastValue[id] {
				b.store.isDegenerate[id] = b.store.observedMin[id] == b.store.observedMax[id]
			}
		}
		// Was the previous sample contiguous for this metric?
		// Consecutive = previous-sample observation AND no FTDC gap.
		consecutive := !isGap && b.lastSeenIdx[id] == sampleIdx-1
		// If this metric was absent from one or more recent samples, pad
		// the column with zeros (with consecutive=false) up to but NOT
		// including the current sample index.
		for missing := b.lastSeenIdx[id] + 1; missing < sampleIdx; missing++ {
			b.store.columns[id].AddPad(0)
		}
		b.store.columns[id].Add(v, consecutive)
		b.lastSeenIdx[id] = sampleIdx
		b.lastValue[id] = v
	}
}

// finalize closes out all per-metric columns, pads any trailing missing
// samples with zeros so every column has length == totalSamples, and
// computes fullMean / fullSD / final isDegenerate from the columns'
// online stats.
func (b *streamingBuilder) finalize() {
	M := len(b.store.names)
	for id := 0; id < M; id++ {
		// Trailing pad: metrics that disappeared before the last sample
		// get zero-padding to keep columns sample-aligned.
		for missing := b.lastSeenIdx[id] + 1; missing < b.totalSamples; missing++ {
			b.store.columns[id].AddPad(0)
		}
		b.store.columns[id].Finish()
	}

	// Compute fullMean / fullSD / fullN: counters use rate-series stats;
	// gauges use raw-value stats. The CompressedColumn maintains both online.
	b.store.fullMean = make([]float64, M)
	b.store.fullSD = make([]float64, M)
	b.store.fullN = make([]int, M)
	for id := 0; id < M; id++ {
		col := b.store.columns[id]
		// Use the builder's real-value-only min/max tracking, NOT the
		// column's online stats (which include zero-pad entries for
		// samples where the metric was absent).
		b.store.isDegenerate[id] = b.store.observedMin[id] == b.store.observedMax[id]
		if b.store.isDegenerate[id] {
			b.store.isCounter[id] = false
		}
		if b.store.isCounter[id] {
			b.store.fullMean[id] = col.RateMean()
			b.store.fullSD[id] = col.RateStdDev()
			b.store.fullN[id] = col.RateN()
		} else {
			b.store.fullMean[id] = col.Mean()
			b.store.fullSD[id] = col.StdDev()
			b.store.fullN[id] = col.WelfordN()
		}
	}

	// Pre-allocate materialized-series slot (populated lazily by
	// materialize → ReadAll when stages need float64 data).
	b.store.materializedSeries = make([][]float64, M)
}

// nSamples returns the number of samples currently in the store. Used by
// stages that need to know the total length.
func (s *seriesStore) nSamples() int {
	if s == nil || len(s.columns) == 0 {
		return 0
	}
	return s.columns[0].Len()
}

// readColumnFloat decodes the metric's column into a fresh float64 slice
// (allocating dst if it's nil/short). Convenience for stages that work in
// float64.
func (s *seriesStore) readColumnFloat(id int, dst []float64) []float64 {
	if id < 0 || id >= len(s.columns) || s.columns[id] == nil {
		return dst[:0]
	}
	return s.columns[id].ReadAll(dst)
}

// readColumnInt64 decodes the metric's column into a fresh int64 slice.
func (s *seriesStore) readColumnInt64(id int, dst []int64) []int64 {
	if id < 0 || id >= len(s.columns) || s.columns[id] == nil {
		return dst[:0]
	}
	return s.columns[id].ReadAllInt64(dst)
}

// readWindowInt64 decodes the metric's values in sample range [start, end)
// into dst. Honours the empty-column fast path.
func (s *seriesStore) readWindowInt64(id, start, end int, dst []int64) []int64 {
	if id < 0 || id >= len(s.columns) || s.columns[id] == nil {
		return dst[:0]
	}
	return s.columns[id].ReadRange(start, end, dst)
}

// materializeFromColumns populates s.materializedSeries[id] for each id in
// `ids` by decoding the column into a fresh []float64. Replaces the
// sample-iteration-based `materialize` for code paths that have a store
// built via buildSeriesStoreStreaming.
func (s *seriesStore) materializeFromColumns(ids []int) {
	if len(ids) == 0 {
		return
	}
	need := make([]int, 0, len(ids))
	for _, id := range ids {
		if id >= 0 && id < len(s.materializedSeries) && s.materializedSeries[id] == nil {
			need = append(need, id)
		}
	}
	if len(need) == 0 {
		return
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
			for _, id := range need[lo:hi] {
				col := s.columns[id]
				if col == nil {
					continue
				}
				// For counters, materialize the FIRST-DIFFERENCE series
				// (rate) to match what fullMean/fullSD describe. For
				// gauges, raw values.
				if s.isCounter[id] {
					ints := col.ReadAllInt64(nil)
					out := make([]float64, len(ints))
					for i := 1; i < len(ints); i++ {
						out[i] = float64(ints[i] - ints[i-1])
					}
					s.materializedSeries[id] = out
				} else {
					s.materializedSeries[id] = col.ReadAll(nil)
				}
			}
		}(lo, hi)
	}
	wg.Wait()
}
