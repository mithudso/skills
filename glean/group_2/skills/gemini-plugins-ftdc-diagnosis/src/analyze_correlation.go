// analyze_correlation.go - Lagged Pearson cross-correlation against a small
// follower vector.
//
// Design note (revised from plan): time-domain Pearson at the 8-element lag
// grid {0, 1, 5, 10, 30, 60, 120, 300}s beats FFT-based xcorr. FFT amortizes
// when you want every lag; we want 8 specific lags, so a single dot-product
// per (leader, follower, lag) triple at O(N) per triple is cheaper than FFT
// (O(N log N) per metric) plus O(N log N) per pair for inverse FFT to
// extract the 8 lag points. Time-domain is also simpler and avoids depending
// on a stable in-process FFT implementation.
//
// Pearson at lag L is computed over the overlap window only:
//
//	r(L) = Σ_valid (x - μ_x)(y - μ_y) / sqrt(SS_x · SS_y)
//
// where μ_x, μ_y, SS_x, SS_y are computed from the count = n - |L| overlap
// pairs, not from the full-series statistics. Using full-series means/SDs
// deflates r by count/(n-1) at non-zero lags — a perfectly correlated pair
// at lag=60 with n=120 would return r=0.5 instead of r=1.0. Workers
// accumulate raw sums (sumX, sumY, sumXY, sumX2, sumY2) per
// (metric, follower, lag) triple; Pearson is computed post-merge in O(1)
// arithmetic using the computational formula for variance. The accumulation
// pass walks samples outside, metrics inside, so each sample's Metrics map
// is iterated exactly once regardless of how many leaders we track.
package main

import (
	"math"
	"sort"
	"sync"
	"time"
)

// defaultFollowerMetrics are the canonical "outcome" metrics we correlate
// every leader against. These names are checked at runtime; if a default
// doesn't exist in the captured metric set (e.g. version skew renamed it),
// it's silently skipped. The bias-control invariant requires that the LLM
// downstream never see a single hand-picked follower — emit correlation
// results against the full vector and let the LLM choose.
var defaultFollowerMetrics = []string{
	"serverStatus.opLatencies.writes.latency",
	"serverStatus.opLatencies.reads.latency",
	"serverStatus.queues.execution.write.out",
	"serverStatus.queues.execution.read.out",
	"serverStatus.connections.current",
}

// lag grid in seconds. Forensics chose these values: tighter than the prior
// {10,30,60,120,300} default because real FTDC is 1Hz and sub-minute cascades
// exist; dropped {600,1800} because correlation at those lags is buried in
// the noise floor at typical capture sizes.
var lagGridSeconds = []int{0, 1, 5, 10, 30, 60, 120, 300}

// correlateAgainstFollowersWithStore is the store-aware Stage 2 variant.
// Uses store.fullMean/fullSD directly so we don't re-Welford the full series.
// Hot loop is pure array work after a single sample-outer / metrics-inner pass.
func correlateAgainstFollowersWithStore(host string, store *seriesStore, samples []Sample, cfg AutoConfig) []Finding {
	if len(samples) < 8 {
		return nil
	}
	intervalSec := modalSampleIntervalSeconds(samples)
	if intervalSec <= 0 {
		intervalSec = 1
	}
	// Build follower vector — defaults + user extension, deduplicated.
	allFollowers := append([]string(nil), defaultFollowerMetrics...)
	allFollowers = append(allFollowers, cfg.FollowerMetrics...)
	sort.Strings(allFollowers)
	allFollowers = dedupStrings(allFollowers)

	// Materialize follower series (small set, only 5-10 usually).
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
	store.materialize(samples, followerIDs)
	followerSeries := make([][]float64, len(followerIDs))
	for i, id := range followerIDs {
		followerSeries[i] = store.materializedSeries[id]
	}
	followerSD := make([]float64, len(followerIDs))
	for i, id := range followerIDs {
		followerSD[i] = store.fullSD[id]
	}

	// Lag grid in samples. Use rounded division (half-up) so fractional
	// intervalSec doesn't truncate small lags to 0 — pre-fix, intervalSec=2.0
	// collapsed the {1s} entry into the {0s} entry, losing information.
	lagSamples := make([]int, 0, len(lagGridSeconds))
	seenLag := map[int]bool{}
	for _, l := range lagGridSeconds {
		var ls int
		if l == 0 {
			ls = 0
		} else {
			ls = int(float64(l)/intervalSec + 0.5)
			// Preserve at least one sample for non-zero requested lags;
			// the rank order matters more than exact match.
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
	n := len(samples)

	// Sample-chunk parallelism. Each worker gets a contiguous range of
	// samples and owns its own accumulator buffer (5 float64 arrays +
	// 1 int32 array). Memory per worker: 5*M*F*L float64 + M*F*L int32
	// ≈ 11 MB on a typical 6500-metric, 5-follower, 8-lag run;
	// 8 workers ≈ 91 MB — the 5-array layout is required for the
	// overlap-window computational Pearson formula (W7-R1).
	//
	// Tier A.4: pin worker count for deterministic FP accumulation.
	workers := welfordWorkers
	if n < workers*100 {
		// Tiny dataset; serial is fine.
		workers = 1
	}
	type partial struct {
		sumX  []float64 // Σ leaderVal over valid overlap pairs
		sumY  []float64 // Σ followerVal over valid overlap pairs
		sumXY []float64 // Σ leaderVal * followerVal
		sumX2 []float64 // Σ leaderVal²
		sumY2 []float64 // Σ followerVal²
		cnt   []int32
	}
	parts := make([]partial, workers)
	chunkSize := (n + workers - 1) / workers

	// Critical correctness: for counter metrics, `store.fullMean[idx]` and
	// `store.fullSD[idx]` are the RATE-series statistics (first
	// differences), not the raw-value statistics. The Pearson r at lag L
	// must therefore compute the leader value AS A RATE for counters,
	// matching what the mean/SD were computed against. Without this, the
	// numerator is `(raw_counter - rate_mean) × (follower_rate - follower_mean)`
	// which is meaningless dimensional nonsense.
	//
	// We compute per-chunk-boundary snapshots of each counter's "previous
	// value" so workers can compute rates within their chunk while
	// preserving cross-chunk rate continuity (same pattern as
	// parallelWelfordWindow in analyze.go).
	boundaryPrev := make([][]int64, workers)
	boundaryHas := make([][]bool, workers)
	for w := 0; w < workers; w++ {
		boundaryPrev[w] = make([]int64, M)
		boundaryHas[w] = make([]bool, M)
	}
	runningPrev := make([]int64, M)
	runningHas := make([]bool, M)
	seenThisSample := make([]bool, M)
	for w := 0; w < workers; w++ {
		copy(boundaryPrev[w], runningPrev)
		copy(boundaryHas[w], runningHas)
		lo := w * chunkSize
		hi := lo + chunkSize
		if hi > n {
			hi = n
		}
		for i := lo; i < hi; i++ {
			for j := range seenThisSample {
				seenThisSample[j] = false
			}
			for name, v := range samples[i].Metrics {
				idx, ok := store.nameToID[name]
				if !ok {
					continue
				}
				seenThisSample[idx] = true
				runningPrev[idx] = v
				runningHas[idx] = true
			}
			for idx := 0; idx < M; idx++ {
				if !seenThisSample[idx] && runningHas[idx] {
					runningHas[idx] = false
				}
			}
		}
	}

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		lo := w * chunkSize
		hi := lo + chunkSize
		if hi > n {
			hi = n
		}
		if lo >= n {
			break
		}
		parts[w] = partial{
			sumX:  make([]float64, M*F*L),
			sumY:  make([]float64, M*F*L),
			sumXY: make([]float64, M*F*L),
			sumX2: make([]float64, M*F*L),
			sumY2: make([]float64, M*F*L),
			cnt:   make([]int32, M*F*L),
		}
		wg.Add(1)
		go func(w, lo, hi int) {
			defer wg.Done()
			psx := parts[w].sumX
			psy := parts[w].sumY
			psxy := parts[w].sumXY
			psx2 := parts[w].sumX2
			psy2 := parts[w].sumY2
			pc := parts[w].cnt
			// Per-worker rate computation state for counter leaders.
			prev := make([]int64, M)
			hasPrev := make([]bool, M)
			copy(prev, boundaryPrev[w])
			copy(hasPrev, boundaryHas[w])
			seen := make([]bool, M)
			for i := lo; i < hi; i++ {
				for j := range seen {
					seen[j] = false
				}
				for name, v := range samples[i].Metrics {
					idx, ok := store.nameToID[name]
					if !ok {
						continue
					}
					seen[idx] = true
					var leaderVal float64
					if store.isCounter[idx] {
						// Rate = first difference. Skip the first
						// observation of each contiguous run; otherwise
						// the rate is undefined.
						if !hasPrev[idx] {
							prev[idx] = v
							hasPrev[idx] = true
							continue
						}
						leaderVal = float64(v - prev[idx])
						prev[idx] = v
					} else {
						leaderVal = float64(v)
					}
					lvSq := leaderVal * leaderVal
					for fIdx := 0; fIdx < F; fIdx++ {
						fSeries := followerSeries[fIdx]
						base := (idx*F + fIdx) * L
						for lagIdx, lag := range lagSamples {
							target := i + lag
							if target < 0 || target >= n {
								continue
							}
							fy := fSeries[target]
							psx[base+lagIdx] += leaderVal
							psy[base+lagIdx] += fy
							psxy[base+lagIdx] += leaderVal * fy
							psx2[base+lagIdx] += lvSq
							psy2[base+lagIdx] += fy * fy
							pc[base+lagIdx]++
						}
					}
				}
				// Reset hasPrev for counters absent from this sample —
				// next observation starts a fresh rate sequence.
				for idx := 0; idx < M; idx++ {
					if !seen[idx] && hasPrev[idx] && store.isCounter[idx] {
						hasPrev[idx] = false
					}
				}
			}
		}(w, lo, hi)
	}
	wg.Wait()

	// Merge partial accumulators.
	sumX  := make([]float64, M*F*L)
	sumY  := make([]float64, M*F*L)
	sumXY := make([]float64, M*F*L)
	sumX2 := make([]float64, M*F*L)
	sumY2 := make([]float64, M*F*L)
	nCounts := make([]int32, M*F*L)
	for _, p := range parts {
		if p.cnt == nil {
			continue
		}
		for k := range nCounts {
			sumX[k]  += p.sumX[k]
			sumY[k]  += p.sumY[k]
			sumXY[k] += p.sumXY[k]
			sumX2[k] += p.sumX2[k]
			sumY2[k] += p.sumY2[k]
			nCounts[k] += p.cnt[k]
		}
	}

	type candidate struct {
		metric, follower string
		lagSec           int
		r                float64
		nUsed            int
	}
	candidates := []candidate{}
	for idx := 0; idx < M; idx++ {
		mSD := store.fullSD[idx]
		if mSD <= 0 {
			continue
		}
		for fIdx := 0; fIdx < F; fIdx++ {
			fSD := followerSD[fIdx]
			if fSD <= 0 {
				continue
			}
			base := (idx*F + fIdx) * L
			bestLag := 0
			bestR := 0.0
			bestN := int32(0)
			for lagIdx := 0; lagIdx < L; lagIdx++ {
				k := base + lagIdx
				cnt := nCounts[k]
				if cnt < 4 {
					continue
				}
				// Compute exact Pearson r over the overlap window using
				// raw sums. Using full-series means/SDs deflates r by
				// count/(n-1) at non-zero lags (a perfect lag-60 pair
				// with n=120 would return r=0.5 instead of r=1.0).
				fc := float64(cnt)
				mx := sumX[k] / fc
				my := sumY[k] / fc
				ssX := sumX2[k] - fc*mx*mx
				ssY := sumY2[k] - fc*my*my
				dotC := sumXY[k] - fc*mx*my
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
				if absFloat(r) > absFloat(bestR) {
					bestR = r
					bestLag = lagIdx
					bestN = cnt
				}
			}
			if absFloat(bestR) >= cfg.CorrelationMinR {
				candidates = append(candidates, candidate{
					metric:   store.names[idx],
					follower: followerNames[fIdx],
					lagSec:   int(float64(lagSamples[bestLag]) * intervalSec),
					r:        bestR,
					nUsed:    int(bestN),
				})
			}
		}
	}
	byFollower := map[string][]candidate{}
	for _, c := range candidates {
		byFollower[c.follower] = append(byFollower[c.follower], c)
	}
	out := []Finding{}
	for _, fName := range followerNames {
		group := byFollower[fName]
		sort.SliceStable(group, func(i, j int) bool {
			ai, aj := absFloat(group[i].r), absFloat(group[j].r)
			if ai != aj {
				return ai > aj
			}
			// Tie-break by metric name so the top-30 cap admits the
			// same metrics across runs even when |r| is identical
			// (common for redundant counter pairs).
			if group[i].metric != group[j].metric {
				return group[i].metric < group[j].metric
			}
			return group[i].lagSec < group[j].lagSec
		})
		if len(group) > 30 {
			group = group[:30]
		}
		for _, c := range group {
			out = append(out, newFinding(KindCorrelation, host, c.metric+"|"+c.follower, map[string]interface{}{
				"leader":      c.metric,
				"follower":    c.follower,
				"lag_seconds": c.lagSec,
				"pearson":     roundFloat(c.r, 4),
				"n_samples":   c.nUsed,
			}))
		}
	}
	return out
}

// modalSampleIntervalSeconds returns the most common interval between
// consecutive samples (rounded to whole seconds). Used to convert
// lagGridSeconds to lag-in-samples.
func modalSampleIntervalSeconds(samples []Sample) float64 {
	if len(samples) < 2 {
		return 1
	}
	counts := map[int]int{}
	for i := 1; i < len(samples); i++ {
		d := samples[i].Timestamp.Sub(samples[i-1].Timestamp)
		if d < 0 {
			continue
		}
		// Round to nearest second.
		secs := int((d + 500*time.Millisecond) / time.Second)
		if secs < 1 {
			secs = 1
		}
		counts[secs]++
	}
	// Tier A.5 fix: iterate in sorted order so ties break deterministically.
	secsList := make([]int, 0, len(counts))
	for s := range counts {
		secsList = append(secsList, s)
	}
	sort.Ints(secsList)
	bestSecs := 1
	bestCount := 0
	for _, s := range secsList {
		c := counts[s]
		if c > bestCount {
			bestCount = c
			bestSecs = s
		}
	}
	return float64(bestSecs)
}

// meanAndSDFloat returns sample mean and sample standard deviation of x.
// Returns (0, 0) when len(x) < 2.
func meanAndSDFloat(x []float64) (float64, float64) {
	n := 0
	mean := 0.0
	m2 := 0.0
	for _, v := range x {
		n++
		delta := v - mean
		mean += delta / float64(n)
		delta2 := v - mean
		m2 += delta * delta2
	}
	if n < 2 {
		return mean, 0
	}
	return mean, sqrtFloat(m2 / float64(n-1))
}

// dedupStrings returns s with adjacent duplicates removed. Input must be
// sorted.
func dedupStrings(s []string) []string {
	if len(s) == 0 {
		return s
	}
	out := s[:1]
	for i := 1; i < len(s); i++ {
		if s[i] != s[i-1] {
			out = append(out, s[i])
		}
	}
	return out
}

// absFloat returns the absolute value of x. Delegates to math.Abs so NaN
// propagates cleanly (`math.Abs(NaN) == NaN`); a hand-rolled `x < 0`
// check would return NaN unchanged but make downstream comparisons (e.g.
// `absFloat(x) >= 100`) silently false for any NaN input.
func absFloat(x float64) float64 {
	return math.Abs(x)
}
