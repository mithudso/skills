// analyze_stratified.go — P3.3 dimension-stratified mover scoring.
//
// Many MongoDB metrics live in high-cardinality families with a shared
// prefix and a varying middle "dimension" token: e.g.
//   serverStatus.queues.execution.<pool>.read.out
//   serverStatus.metrics.repl.network.replSetUpdatePosition.<member_id>.numOps
// Stage 1 scores each instance independently. When the regression is
// concentrated in a single dimension value (e.g. shard 5, pool foo),
// the family looks like ~50 noisy metrics + 1 outlier, which Stage 1
// surfaces only as a mid-rank single mover. This stage groups by family
// and surfaces the asymmetry explicitly: "of N dim values, k stand out;
// dimension explains the variance."
//
// Algorithm:
//   1. For each metric path, generate candidate family signatures by
//      replacing each token (path segment) in turn with the placeholder
//      `*`. A token is a candidate dimension if it varies across the
//      family — at least 3 distinct values seen with the same family
//      signature.
//   2. Score each candidate family by Welch t-statistic on baseline vs
//      stress means (computed from store's fullMean / Welford accumulators
//      — already in the store).
//   3. If only `k ≤ ceil(N/3)` dim values have |Welch t| ≥ family-mean-t
//      + 2*family-stdev-t, and at least 1 such outlier exists, emit
//      `kind:"mover_stratified"` with the family, the hot dim values,
//      and a chi-squared statistic for "dimension explains variance."
package main

import (
	"math"
	"sort"
	"strings"
)

// stratifiedMinFamily is the minimum number of distinct dimension
// values required to call something a "family worth stratifying."
const stratifiedMinFamily = 3

// stratifiedMaxFamily bounds the per-family work; very wide families
// (e.g. >200 dim values) are uninteresting for the "single hot dim"
// pattern and can dominate runtime.
const stratifiedMaxFamily = 200

// stratifiedOutlierZ is the number of standard deviations a per-dim
// Welch t-statistic must exceed the family-mean t to be flagged as a
// hot dimension.
const stratifiedOutlierZ = 2.0

// emitMoverStratified scans the store's name set for variable-token
// families and emits mover_stratified Findings.
//
// Recomputes per-window Welford accumulators via the same primitives
// Stage 1 uses (parallelWelfordWindowCol). The recomputation is bounded
// by M ≈ 6500 metrics × baseline+stress window sizes — modest relative
// to Stage 1 itself.
func emitMoverStratified(c *findingCollector, host string, store *seriesStore, bStart, bEnd, sStart, sEnd int) (emitted int, skipped int) {
	baseN, baseMean, baseM2, _ := parallelWelfordWindowCol(store, bStart, bEnd)
	stressN, stressMean, stressM2, _ := parallelWelfordWindowCol(store, sStart, sEnd)
	M := len(store.names)
	if M < stratifiedMinFamily {
		return 0, 1
	}

	// Family signature → list of (idx, dim_value).
	type famMember struct {
		id  int
		dim string
	}
	families := map[string][]famMember{}

	for id, name := range store.names {
		tokens := strings.Split(name, ".")
		if len(tokens) < 3 {
			// Need at least three tokens for a meaningful family.
			continue
		}
		// Try replacing each non-trivial token in turn. Skip the first
		// token (always "serverStatus" / "systemMetrics" — too coarse)
		// and the last (the "leaf" metric name — that's typically the
		// thing we want to keep stable).
		for pos := 1; pos < len(tokens)-1; pos++ {
			t := tokens[pos]
			// Avoid trivially over-grouping — single-character or
			// well-known terminal tokens (e.g. "read", "write") rarely
			// vary AS a dimension.
			if len(t) < 2 {
				continue
			}
			if commonFixedToken(t) {
				continue
			}
			// Build the signature: replace position pos with "*".
			sigTokens := make([]string, len(tokens))
			copy(sigTokens, tokens)
			sigTokens[pos] = "*"
			sig := strings.Join(sigTokens, ".")
			families[sig] = append(families[sig], famMember{id: id, dim: t})
		}
	}

	// Filter and process families of interest.
	famKeys := make([]string, 0, len(families))
	for sig := range families {
		if len(families[sig]) < stratifiedMinFamily || len(families[sig]) > stratifiedMaxFamily {
			continue
		}
		famKeys = append(famKeys, sig)
	}
	sort.Strings(famKeys)

	// To avoid double-emit (a single metric can land in multiple family
	// signatures via different positions), track which metric IDs have
	// already been emitted in some family and prefer the family with
	// the most members.
	sort.SliceStable(famKeys, func(i, j int) bool {
		return len(families[famKeys[i]]) > len(families[famKeys[j]])
	})
	emittedIDs := map[int]bool{}

	for _, sig := range famKeys {
		members := families[sig]
		// Skip if every member already emitted in a wider family.
		fresh := 0
		for _, m := range members {
			if !emittedIDs[m.id] {
				fresh++
			}
		}
		if fresh < stratifiedMinFamily {
			continue
		}

		// Compute per-member Welch t.
		type scored struct {
			dim string
			id  int
			t   float64
		}
		scoredMembers := make([]scored, 0, len(members))
		// Single-pass Welford accumulation for sample mean and M2.
		// Bessel-corrected variance = M2/(n-1) is consistent at small
		// family sizes; population variance would be 22% too small at n=3.
		var wMean, wM2 float64
		var tCount int
		for _, m := range members {
			bSD := welfordSD(baseN[m.id], baseM2[m.id])
			sSD := welfordSD(stressN[m.id], stressM2[m.id])
			t := welchT(baseMean[m.id], bSD, baseN[m.id], stressMean[m.id], sSD, stressN[m.id])
			if math.IsNaN(t) || math.IsInf(t, 0) {
				continue
			}
			scoredMembers = append(scoredMembers, scored{dim: m.dim, id: m.id, t: t})
			tCount++
			delta := t - wMean
			wMean += delta / float64(tCount)
			wM2 += delta * (t - wMean)
		}
		if tCount < stratifiedMinFamily {
			continue
		}
		mean := wMean
		variance := 0.0
		if tCount > 1 {
			variance = wM2 / float64(tCount-1)
		}
		if variance < 0 {
			variance = 0
		}
		stdev := math.Sqrt(variance)

		// Identify outliers: members whose |t - mean| exceeds Z * stdev.
		// Also require absolute |t| above 2 so we don't flag a family
		// of all-quiet metrics where the "outlier" is still uninteresting.
		var outliers []scored
		for _, s := range scoredMembers {
			if math.Abs(s.t-mean) >= stratifiedOutlierZ*stdev && math.Abs(s.t) >= 2.0 {
				outliers = append(outliers, s)
			}
		}
		if len(outliers) == 0 {
			continue
		}
		// Cap outlier count to ⌈N/3⌉ — if more than a third of the
		// family is "hot," there's no asymmetry; the whole family is
		// stressed and Stage 1 already surfaces it via individual movers.
		cap := (len(scoredMembers) + 2) / 3
		if len(outliers) > cap {
			continue
		}

		// Chi-squared sanity: divide observed t^2 values into bins
		// (high vs low), expected uniform across dim values, chi-sq
		// statistic = sum((O - E)^2 / E). High value = dimension explains
		// variance.
		tSqVals := make([]float64, len(scoredMembers))
		for i, s := range scoredMembers {
			tSqVals[i] = s.t * s.t
		}
		chi2 := chiSquaredDimUniform(tSqVals)

		sort.Slice(outliers, func(i, j int) bool {
			if math.Abs(outliers[i].t) != math.Abs(outliers[j].t) {
				return math.Abs(outliers[i].t) > math.Abs(outliers[j].t)
			}
			return outliers[i].dim < outliers[j].dim
		})
		hotDimList := make([]map[string]interface{}, 0, len(outliers))
		for _, o := range outliers {
			hotDimList = append(hotDimList, map[string]interface{}{
				"dimension_value": o.dim,
				"metric":          store.names[o.id],
				"welch_t":         roundFloat(o.t, 4),
			})
			emittedIDs[o.id] = true
		}
		// Mark family members as emitted so a position-overlapping family
		// doesn't double-fire on the same metrics.
		for _, s := range scoredMembers {
			emittedIDs[s.id] = true
		}

		c.add(newFinding("mover_stratified", host, sig, map[string]interface{}{
			"family_signature":    sig,
			"n_dimensions":        len(scoredMembers),
			"n_outlier_dimensions": len(outliers),
			"family_mean_welch_t": roundFloat(mean, 4),
			"family_stdev_welch_t": roundFloat(stdev, 4),
			"chi_squared":         roundFloat(chi2, 4),
			"hot_dimensions":      hotDimList,
		}))
		emitted++
	}
	if emitted == 0 {
		return 0, 1
	}
	return emitted, 0
}

// commonFixedToken returns true for tokens that are too generic to
// usefully treat as a "dimension." These appear repeatedly across
// otherwise-distinct metric trees and would explode the family count
// with meaningless groupings.
func commonFixedToken(t string) bool {
	switch t {
	case "read", "write", "total", "ops", "count", "latency",
		"current", "available", "ms", "us", "ns", "secs", "seconds",
		"in", "out", "rate", "max", "min", "size":
		return true
	}
	return false
}

// chiSquaredDimUniform computes a chi-squared statistic against the
// null hypothesis "all dim values contribute equally to the family's
// variance." Inputs are per-dim squared-t observations (variance proxy).
// Higher = stronger evidence that the dimension explains variance.
func chiSquaredDimUniform(tSquared []float64) float64 {
	if len(tSquared) < 2 {
		return 0
	}
	total := 0.0
	for _, v := range tSquared {
		total += v
	}
	if total == 0 {
		return 0
	}
	expected := total / float64(len(tSquared))
	chi2 := 0.0
	for _, v := range tSquared {
		diff := v - expected
		chi2 += (diff * diff) / expected
	}
	return chi2
}
