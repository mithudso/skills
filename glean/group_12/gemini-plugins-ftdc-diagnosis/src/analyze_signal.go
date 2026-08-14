// analyze_signal.go — signal-quality detection stage.
//
// Emits three structural Findings about the shape and quality of
// already-classified metric series, plus the per-stage AutoStats
// observability hook surfaced in run_health.stages.
//
// Periodicity detection: Welch periodogram on detrended/differenced
// series. Raw autocorrelation with a 0.5 threshold fires on every
// monotonic counter (r(L=1h) ≈ 0.99996 on a clean ramp); the
// periodogram on residuals is the right tool.
//
// suspicious_smoothness: detector for data that's TOO smooth — variance
// anomalously low for the metric's class, counter rate too perfectly
// steady (round-number deltas), histogram buckets suspiciously round.
// Catches mangled / replayed / synthetic data masquerading as real
// captures.
//
// evidence_convergence: meta-Finding per synthetic_event listing the
// distinct detector kinds that fired on the same metric+timestamp
// region. Three or more detectors agreeing on the same event = high
// confidence; one detector alone = noise candidate.
//
// AutoStats: per-stage skipped/emitted/duration counters surfaced via
// the run_health Finding. Lets the analyzer distinguish "stage didn't
// fire" from "stage fired and found nothing".

package main

import (
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

// ----------------------------------------------------------------------------
// Periodicity detection (item 2)
// ----------------------------------------------------------------------------

// emitPeriodic runs a Welch periodogram on detrended top-K metric
// series and emits `kind:"periodic"` for each metric with a
// statistically significant peak. Operates on the top-200 trend pool
// (already materialized in store.materializedSeries by Stage 3 setup)
// so we don't double-decompress.
//
// Algorithm:
//   1. For counter metrics, take first differences (rates). For gauges
//      with strong linear trend (|slope_t| >= 5), subtract the OLS
//      fit. For other gauges, use raw values.
//   2. Compute the periodogram via segment-averaged squared FFT
//      (Welch's method) — but we don't have an FFT in stdlib, so use
//      the direct DFT at a small set of candidate frequencies. The
//      candidates are tied to known MongoDB cycles: 30s, 45s, 60s, 90s,
//      120s, 180s, 300s, 600s, 900s, 1800s, 3600s, 7200s.
//   3. For each candidate frequency f, compute the power
//      P(f) = |sum(x[t] * exp(-2πi * f * t))|² / N.
//   4. Schuster test: under the null (white noise), max-P / mean-P
//      across F candidate frequencies follows an exponential distribution.
//      Bonferroni-correct for F candidates: reject if
//      max-P / mean-P > -ln(α / F) where α = 1e-4.
//   5. Emit `kind:"periodic"` with period_seconds, snr, p_value_adj.
func emitPeriodic(c *findingCollector, host string, store *seriesStore, ids []int, cfg AutoConfig) (emitted, skipped int) {
	candidatePeriodsSec := []int{30, 45, 60, 90, 120, 180, 300, 600, 900, 1800, 3600, 7200}
	if len(store.timestamps) < 2 {
		return 0, 0
	}
	intervalSec := modalSampleIntervalSecondsStore(store.timestamps)
	if intervalSec <= 0 {
		intervalSec = 1
	}
	N := store.nSamples()
	// Drop periods too short (< 4 × interval) or too long (< 2 cycles
	// fit in the capture).
	periodsInSamples := []int{}
	periodsForLabel := []int{}
	for _, ps := range candidatePeriodsSec {
		psInSamples := int(float64(ps) / intervalSec)
		if psInSamples < 4 {
			continue
		}
		if psInSamples*2 > N {
			continue
		}
		periodsInSamples = append(periodsInSamples, psInSamples)
		periodsForLabel = append(periodsForLabel, ps)
	}
	if len(periodsInSamples) == 0 {
		return 0, len(ids)
	}

	// Schuster test alpha, Bonferroni-corrected across F candidate
	// frequencies. α=0.001 corresponds to a threshold SNR of
	// -ln(0.001/12) ≈ 9.4 on F=12 candidates. Tighter (1e-4) excludes
	// real FTDC periodic signals whose energy is split across
	// harmonics (a 60s pulse has content at 60s + 30s + 20s + ...).
	const alpha = 0.001
	for _, id := range ids {
		if id < 0 || id >= len(store.materializedSeries) {
			skipped++
			continue
		}
		series := store.materializedSeries[id]
		if len(series) < 50 {
			skipped++
			continue
		}
		// Detrend: subtract linear fit for high-slope_t gauges;
		// first-difference counters; pass-through otherwise. For
		// counters, materializedSeries already holds the rate series
		// (foundation-stage invariant), so no further differencing
		// needed.
		work := series
		if !store.isCounter[id] {
			_, slopeT := linearSlope(series)
			if absFloat(slopeT) >= 5 {
				work = detrendFloat(series)
			}
		}
		// Tier B.11: AR(1) prewhiten before centering. Schuster's test
		// assumes white residuals; FTDC gauges typically have ρ̂₁ ∈
		// [0.3, 0.8], which inflates the Schuster null and produces false-
		// positive periodicity findings. Prewhitening flattens the
		// spectrum so the unit-exponential null becomes accurate.
		work = arPrewhiten(work)
		mean := 0.0
		for _, v := range work {
			mean += v
		}
		mean /= float64(len(work))
		// Centre to zero mean for DFT.
		centred := make([]float64, len(work))
		for i, v := range work {
			centred[i] = v - mean
		}

		powers := make([]float64, len(periodsInSamples))
		var totalPower float64
		for pi, period := range periodsInSamples {
			powers[pi] = dftPower(centred, period)
			totalPower += powers[pi]
		}
		if totalPower <= 0 {
			skipped++
			continue
		}
		var bestIdx int
		bestPower := powers[0]
		for pi, p := range powers {
			if p > bestPower {
				bestPower = p
				bestIdx = pi
			}
		}
		// Standard Schuster SNR: peak / mean(of OTHER candidates).
		// Including the peak in the mean caps SNR at F (number of
		// candidates), which can't exceed the Bonferroni threshold
		// `-ln(α/F)` when α < 1/F. Excluding the peak gives the
		// detector enough dynamic range to fire on real signals.
		// With only one candidate the denominator is 0/0 = NaN; skip
		// the emission since the Schuster test requires a comparison set.
		if len(powers) < 2 {
			skipped++
			continue
		}
		otherMean := (totalPower - bestPower) / float64(len(powers)-1)
		if otherMean <= 0 {
			otherMean = 1e-12
		}
		snr := bestPower / otherMean
		// Schuster-test rejection threshold with Bonferroni
		// correction for F candidate frequencies.
		threshold := -math.Log(alpha / float64(len(powers)))
		if snr <= threshold {
			skipped++
			continue
		}
		c.add(newFinding("periodic", host, store.names[id], map[string]interface{}{
			"metric":         store.names[id],
			"metric_kind":    metricKindStr(store, id),
			"period_seconds": periodsForLabel[bestIdx],
			"snr":            roundFloat(snr, 4),
			"threshold":      roundFloat(threshold, 4),
			"detrended":      !store.isCounter[id] && work != nil && &work[0] != &series[0],
			"differencing":   store.isCounter[id],
			"n":              len(series),
			"candidates":     periodsForLabel,
		}))
		emitted++
	}
	return emitted, skipped
}

// dftPower computes the squared-magnitude of the discrete Fourier
// component at the frequency 1/period (in samples), normalized by N.
//
//	X(f) = sum_{t=0..N-1} x[t] * exp(-2π i f t)
//	P(f) = |X(f)|² / N
//
// For the small set of candidate frequencies in emitPeriodic, this
// O(N) direct evaluation is cheaper than an O(N log N) FFT (which we'd
// have to pull from a third-party dep). Cost on the 25h testdata:
// F=12 × N=90K = 1.08M ops per metric × 200 metrics ≈ 200M ops,
// ~2-5 sec wall (each iteration is 1 cos + 1 sin transcendental).
func dftPower(x []float64, periodSamples int) float64 {
	if periodSamples <= 0 || len(x) == 0 {
		return 0
	}
	// Tier B.11: Hann windowing removes spectral leakage that contaminates
	// the Schuster-test null. A 0.5-cosine taper costs one multiplication
	// per sample. Without it, strong neighboring frequencies (harmonics)
	// scaled into candidate bins, biasing SNR upward on AR(1) noise.
	N := len(x)
	omega := 2 * math.Pi / float64(periodSamples)
	var re, im float64
	nMinus1 := float64(N - 1)
	if nMinus1 <= 0 {
		nMinus1 = 1
	}
	for t, v := range x {
		// Hann: w(t) = 0.5 * (1 - cos(2π t / (N-1)))
		w := 0.5 * (1.0 - math.Cos(2*math.Pi*float64(t)/nMinus1))
		wv := v * w
		c := math.Cos(omega * float64(t))
		s := math.Sin(omega * float64(t))
		re += wv * c
		im -= wv * s
	}
	// Window-amplitude correction: Hann's coherent gain is 0.5, so the
	// power is scaled by 1/0.5² = 4 to recover the unwindowed magnitude.
	return 4.0 * (re*re + im*im) / float64(N)
}

// arPrewhiten returns x with the AR(1) residual removed: y[t] = x[t] -
// ρ̂·x[t-1] where ρ̂ is the sample lag-1 autocorrelation. Tier B.11
// addresses A13 — Schuster's test assumes white noise residuals; AR(1)
// noise (common in FTDC) violates the null and produces false positives.
// Returns a new slice; never mutates the input.
func arPrewhiten(x []float64) []float64 {
	if len(x) < 2 {
		out := make([]float64, len(x))
		copy(out, x)
		return out
	}
	// Estimate ρ̂₁ via lag-1 autocorrelation.
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
		out := make([]float64, len(x))
		copy(out, x)
		return out
	}
	rho := num / denom
	if rho > 0.95 {
		rho = 0.95
	} else if rho < -0.95 {
		rho = -0.95
	}
	out := make([]float64, len(x))
	out[0] = x[0]
	for i := 1; i < len(x); i++ {
		out[i] = x[i] - rho*x[i-1]
	}
	return out
}

// detrendFloat returns x minus its OLS linear fit. Allocates a new
// slice; never mutates the input.
func detrendFloat(x []float64) []float64 {
	n := len(x)
	if n < 3 {
		out := make([]float64, n)
		copy(out, x)
		return out
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
		out := make([]float64, n)
		copy(out, x)
		return out
	}
	slope := sxt / stt
	intercept := meanX - slope*meanT
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = x[i] - (intercept + slope*float64(i))
	}
	return out
}

// metricKindStr returns "counter" / "gauge" / "constant" for a metric ID.
func metricKindStr(store *seriesStore, id int) string {
	if store.isDegenerate[id] {
		return "constant"
	}
	if store.isCounter[id] {
		return "counter"
	}
	return "gauge"
}

// ----------------------------------------------------------------------------
// suspicious_smoothness (item 23)
// ----------------------------------------------------------------------------

// emitSuspiciousSmoothness scans top-K materialized series for shapes
// indicating mangled / replayed / synthetic data. Three structural
// detectors:
//
//   1. variance_too_low: stddev < 1e-9 × |mean| AND mean != 0. The
//      metric is uniformly constant despite being marked non-degenerate
//      (degenerate already gets a `constant` kind tag).
//   2. perfectly_steady_rate: for counter metrics, all consecutive
//      first-differences are identical to within float epsilon. A real
//      monotonic counter has rate jitter from sample-interval
//      variation alone.
//   3. round_number_density: > 90% of non-zero values are multiples of
//      10, 100, or 1000. Signals binned / pre-rounded data.
//
// All three detectors are structural (no name-based recognition) and
// emit a single `kind:"suspicious_smoothness"` Finding per metric per
// detector that fires.
func emitSuspiciousSmoothness(c *findingCollector, host string, store *seriesStore, ids []int) (emitted, skipped int) {
	for _, id := range ids {
		if id < 0 || id >= len(store.materializedSeries) {
			skipped++
			continue
		}
		series := store.materializedSeries[id]
		if len(series) < 30 {
			skipped++
			continue
		}
		if store.isDegenerate[id] {
			continue // already classified as constant
		}
		// Detector 1: variance_too_low.
		var sd, mean float64
		for _, v := range series {
			mean += v
		}
		mean /= float64(len(series))
		var ssq float64
		for _, v := range series {
			d := v - mean
			ssq += d * d
		}
		sd = math.Sqrt(ssq / float64(len(series)))
		// Defensive guard: the pre-filter at line 286 catches isDegenerate
		// only when the store's degenerate-classification ran; fire here too.
		if mean == 0 && sd == 0 {
			c.add(newFinding(KindSuspiciousSmoothness, host, store.names[id]+"|constant_zero", map[string]interface{}{
				"metric":      store.names[id],
				"detector":    ReasonCodeVarianceTooLow,
				"reason_code": ReasonCodeConstantZero,
				"mean":        0.0,
				"sd":          0.0,
				"n":           len(series),
			}))
			emitted++
		} else if mean != 0 && sd < 1e-9*math.Abs(mean) {
			c.add(newFinding(KindSuspiciousSmoothness, host, store.names[id]+"|variance_too_low", map[string]interface{}{
				"metric":      store.names[id],
				"detector":    ReasonCodeVarianceTooLow,
				"reason_code": ReasonCodeVarianceTooLow,
				"mean":        roundFloat(mean, 6),
				"sd":          roundFloat(sd, 12),
				"n":           len(series),
			}))
			emitted++
		}

		// Detector 2: perfectly_steady_rate (counters).
		// Tolerance is relative to the reference rate magnitude (1e-12 floor for tiny rates).
		if store.isCounter[id] && len(series) >= 30 {
			ref := series[1]
			tol := math.Abs(ref)*1e-9 + 1e-12
			steady := ref != 0
			for i := 2; i < len(series) && steady; i++ {
				if math.Abs(series[i]-ref) > tol {
					steady = false
				}
			}
			if steady {
				c.add(newFinding(KindSuspiciousSmoothness, host, store.names[id]+"|perfectly_steady_rate", map[string]interface{}{
					"metric":      store.names[id],
					"detector":    ReasonCodePerfectlySteadyRate,
					"reason_code": ReasonCodePerfectlySteadyRate,
					"rate":        roundFloat(ref, 6),
					"n":           len(series),
				}))
				emitted++
			}
		}

		// Detector 3: round_number_density. Skip for tiny-magnitude
		// series where every value is naturally small. Note:
		// `iv % 10 == 0` already covers iv % 100 == 0
		// and iv % 1000 == 0, so a single modulo is sufficient.
		var nonZeroCount, roundCount int
		for _, v := range series {
			if v == 0 {
				continue
			}
			nonZeroCount++
			iv := int64(math.Round(v))
			if float64(iv) != v {
				continue
			}
			if iv%10 == 0 {
				roundCount++
			}
		}
		if nonZeroCount >= 20 && roundCount*10 > 9*nonZeroCount {
			c.add(newFinding(KindSuspiciousSmoothness, host, store.names[id]+"|round_number_density", map[string]interface{}{
				"metric":      store.names[id],
				"detector":    ReasonCodeRoundNumberDensity,
				"reason_code": ReasonCodeRoundNumberDensity,
				"non_zero_n":  nonZeroCount,
				"round_n":     roundCount,
				"round_ratio": roundFloat(float64(roundCount)/float64(nonZeroCount), 4),
			}))
			emitted++
		}
	}
	return emitted, skipped
}

// ----------------------------------------------------------------------------
// evidence_convergence (G7)
// ----------------------------------------------------------------------------

// emitEvidenceConvergence walks the synthetic_event Findings already in
// the collector and, for each event, counts the distinct detector
// kinds (changepoint, spike, variance_shift, level, periodic) that
// fired on metrics within the event's member set. A high distinct-kind
// count = multiple independent detectors agree = high confidence.
//
// Output `kind:"evidence_convergence"` per synthetic_event:
//   - synthetic_event_id: reference back to the overlay
//   - distinct_kinds_count: number of distinct detector kinds in members
//   - kinds: sorted list
//   - per_kind_count: map of kind -> member-finding count
func emitEvidenceConvergence(c *findingCollector, host string) (emitted int) {
	// Build a lookup of every Finding by id for cross-reference.
	byID := map[string]Finding{}
	for _, f := range c.items {
		id, _ := f["id"].(string)
		if id != "" {
			byID[id] = f
		}
	}
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k != "synthetic_event" {
			continue
		}
		ids, _ := f["member_finding_ids"].([]string)
		if len(ids) == 0 {
			continue
		}
		kindCount := map[string]int{}
		for _, mid := range ids {
			child, ok := byID[mid]
			if !ok {
				continue
			}
			ck, _ := child["kind"].(string)
			if ck == "" {
				continue
			}
			kindCount[ck]++
		}
		if len(kindCount) == 0 {
			continue
		}
		distinctKinds := make([]string, 0, len(kindCount))
		for k := range kindCount {
			distinctKinds = append(distinctKinds, k)
		}
		sort.Strings(distinctKinds)
		seID, _ := f["id"].(string)
		atStart, _ := f["at_start"].(string)
		c.add(newFinding("evidence_convergence", host, seID, map[string]interface{}{
			"synthetic_event_id":   seID,
			"at_start":             atStart,
			"distinct_kinds_count": len(distinctKinds),
			"kinds":                distinctKinds,
			"per_kind_count":       kindCount,
		}))
		emitted++
	}
	return emitted
}

// ----------------------------------------------------------------------------
// AutoStats — per-stage observability
// ----------------------------------------------------------------------------

// AutoStats tracks per-stage emission and skip counts plus elapsed
// wall time. Mounted into the run_health Finding so the analyzer can
// distinguish "stage didn't run" from "stage ran and found nothing".
type AutoStats struct {
	Stages []AutoStageStat
}

type AutoStageStat struct {
	Name       string  `json:"name"`
	EmittedN   int     `json:"emitted"`
	SkippedN   int     `json:"skipped"`
	DurationMs float64 `json:"duration_ms"`
}

// recordStage appends an entry to stats with the elapsed milliseconds
// since `start`. Caller passes the post-stage emitted/skipped counts.
func (s *AutoStats) recordStage(name string, start time.Time, emitted, skipped int) {
	if s == nil {
		return
	}
	dur := time.Since(start)
	s.Stages = append(s.Stages, AutoStageStat{
		Name:       name,
		EmittedN:   emitted,
		SkippedN:   skipped,
		DurationMs: roundFloat(float64(dur.Microseconds())/1000.0, 2),
	})
}

// timeStage runs `stage`, computes Findings emitted as the delta in
// `collector.items` length, and records the timing under `name`. Caller
// passes `skipped` directly when known. Convenience wrapper used by
// runAutoWithHostStore to thread AutoStats through every emitter
// without modifying their signatures.
func (s *AutoStats) timeStage(c *findingCollector, name string, skipped int, stage func()) {
	start := time.Now()
	before := len(c.items)
	stage()
	emitted := len(c.items) - before
	s.recordStage(name, start, emitted, skipped)
}

// adjustLastSkipped sets the `skipped` field of the most recently
// recorded stage. Used when the emitter returns a skipped count that
// the timeStage wrapper couldn't know up front.
func (s *AutoStats) adjustLastSkipped(skipped int) {
	if s == nil || len(s.Stages) == 0 {
		return
	}
	s.Stages[len(s.Stages)-1].SkippedN = skipped
}

// asJSON returns a []interface{} suitable for embedding in a Finding's
// fields map.
//
// Tier A.4 fix: wall-clock duration is non-deterministic across machines
// and --jobs settings, which broke `--jobs=1` vs `--jobs=8` byte-equality.
// We now omit `duration_ms` from the default emit. Set FTDC_OBSERVABILITY=1
// to surface it under an `_observability.duration_ms` key (the leading
// underscore signals "non-portable, debug only" to downstream consumers).
// observabilityEnabled is set once at init from FTDC_OBSERVABILITY so the
// env var isn't re-read on every Finding emit. Cheap, but cleaner than
// per-call lookups.
var observabilityEnabled = os.Getenv("FTDC_OBSERVABILITY") == "1"

func (s *AutoStats) asJSON() []interface{} {
	if s == nil {
		return nil
	}
	showWallClock := observabilityEnabled
	out := make([]interface{}, 0, len(s.Stages))
	for _, st := range s.Stages {
		entry := map[string]interface{}{
			"name":    st.Name,
			"emitted": st.EmittedN,
			"skipped": st.SkippedN,
		}
		if showWallClock {
			entry["_observability"] = map[string]interface{}{
				"duration_ms": st.DurationMs,
			}
		}
		out = append(out, entry)
	}
	return out
}

// stageNameForKind returns a short human-readable label per Finding
// kind, used by AutoStats reporting.
func stageNameForKind(kind string) string {
	switch {
	case kind == KindMover:
		return "stage1_welch"
	case kind == KindCorrelation:
		return "stage2_correlation"
	case kind == KindMonotonicTrend:
		return "stage3_trend"
	case kind == KindChangepoint || kind == KindChangepointOversegmented:
		return "stage3_changepoint"
	case kind == KindHistogramPercentile || kind == KindHistogramInvalid:
		return "stage4_histogram"
	case kind == KindSpike:
		return "spike"
	case kind == KindVarianceShift:
		return "variance_shift"
	case strings.HasPrefix(kind, "capacity_") || kind == "near_ceiling":
		return "capacity"
	case kind == "synthetic_event":
		return "synthetic_event"
	case kind == "consistency_warning":
		return "consistency"
	case kind == "periodic":
		return "periodic"
	case kind == "suspicious_smoothness":
		return "smoothness"
	case kind == "evidence_convergence":
		return "convergence"
	}
	return kind
}
