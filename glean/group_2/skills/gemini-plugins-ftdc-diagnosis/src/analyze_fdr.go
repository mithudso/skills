// analyze_fdr.go — Benjamini-Hochberg FDR correction for mover rankings.
//
// Why this exists (P0.NEW-7):
//   The parser runs Welch's t-test against ~6,500 metrics per capture.
//   At α=0.05, ~325 tests are expected to look significant by chance
//   alone — every solo user's "mover" list is mostly false positives
//   unless the threshold is adjusted. The Bonferroni-corrected alpha
//   (≈ 7.7e-6) is over-conservative; Benjamini-Hochberg controls the
//   false discovery rate (FDR) instead of family-wise error rate, which
//   is the right knob for "I want most of my Findings to be real, even
//   if a few false positives slip in."
//
// Precedent in this repo: `analyze_signal.go` already runs
// Bonferroni-corrected Schuster periodicity (Tier B.11). Same posture,
// different correction.
//
// Interface:
//
//	q := benjaminiHochberg(pValues)  // q[i] is the BH-adjusted q-value
//	                                 // for the i'th input p-value.
//	threshold := bhSignificanceThreshold(pValues, 0.05)
//	                                 // largest p in the input that
//	                                 // remains significant at FDR=0.05.
//
// Both helpers preserve input order; only an internal working slice is
// sorted. Complexity is O(N log N) for the sort + O(N) for the monotone
// envelope. Allocates two N-sized slices.
package main

import (
	"math"
	"sort"
)

// benjaminiHochberg returns a slice of q-values in the same order as the
// input p-values. NaN inputs (e.g. degenerate metrics with df ≤ 0) are
// preserved as NaN in the output and excluded from the BH computation.
// q-values are clamped to [0, 1].
func benjaminiHochberg(pValues []float64) []float64 {
	n := len(pValues)
	q := make([]float64, n)
	if n == 0 {
		return q
	}

	// Build a working slice of (original_index, p) pairs for non-NaN
	// inputs only. Tests that produced NaN p-values were not really
	// "tested" — including them inflates m and over-corrects everyone.
	type pp struct {
		idx int
		p   float64
	}
	working := make([]pp, 0, n)
	for i, p := range pValues {
		q[i] = math.NaN() // default for NaN inputs
		if math.IsNaN(p) {
			continue
		}
		// Clamp to [0,1] defensively; some t-survival approximations can
		// return tiny negatives or > 1 from rounding.
		if p < 0 {
			p = 0
		}
		if p > 1 {
			p = 1
		}
		working = append(working, pp{idx: i, p: p})
	}
	m := len(working)
	if m == 0 {
		return q
	}

	// Sort ascending by p-value (deterministic on ties via original idx).
	sort.SliceStable(working, func(i, j int) bool {
		if working[i].p != working[j].p {
			return working[i].p < working[j].p
		}
		return working[i].idx < working[j].idx
	})

	// Apply BH-1995: q_(i) = p_(i) * m / i, with monotone-non-decreasing
	// envelope from the right so q-values are themselves a valid q-rank.
	mf := float64(m)
	rawQ := make([]float64, m)
	for i := 0; i < m; i++ {
		raw := working[i].p * mf / float64(i+1)
		if raw > 1 {
			raw = 1
		}
		rawQ[i] = raw
	}
	// Right-to-left running minimum enforces monotonicity: q_(i) is
	// non-decreasing in i, so a more-significant p never gets a worse
	// q than a less-significant one.
	minSoFar := math.Inf(1)
	for i := m - 1; i >= 0; i-- {
		if rawQ[i] < minSoFar {
			minSoFar = rawQ[i]
		}
		q[working[i].idx] = minSoFar
	}
	return q
}

// bhSignificanceThreshold returns the largest p-value in the input set
// that would still be flagged significant at the given FDR level
// (typically 0.05). Return values:
//   -1.0 — no test passed BH correction (sentinel; no valid threshold exists)
//    0.0 — returned for invalid alpha or an all-NaN input slice
//   (0,1] — the effective threshold p-value
//
// Callers must treat -1.0 as "no significance found," not as a p-value.
// alpha must be in (0, 1).
func bhSignificanceThreshold(pValues []float64, alpha float64) float64 {
	if alpha <= 0 || alpha >= 1 {
		return 0
	}
	// Reuse the working list shape from benjaminiHochberg.
	working := make([]float64, 0, len(pValues))
	for _, p := range pValues {
		if math.IsNaN(p) {
			continue
		}
		if p < 0 {
			p = 0
		}
		if p > 1 {
			p = 1
		}
		working = append(working, p)
	}
	m := len(working)
	if m == 0 {
		return 0
	}
	sort.Float64s(working)
	mf := float64(m)
	// Largest i such that p_(i) <= (i/m) * alpha; return p_(i).
	// Returns -1.0 when no test passes (sentinel; valid p-values are in [0,1]).
	best := -1.0
	for i := 0; i < m; i++ {
		bhCut := float64(i+1) / mf * alpha
		if working[i] <= bhCut {
			best = working[i]
		}
	}
	return best
}

// countSignificantBH returns the number of q-values strictly below alpha.
// Useful for `run_health` accounting + analyzer narration ("8 of 6,432
// metrics survived FDR-BH at α=0.05").
func countSignificantBH(qValues []float64, alpha float64) int {
	n := 0
	for _, q := range qValues {
		if !math.IsNaN(q) && q < alpha {
			n++
		}
	}
	return n
}

// countNonNaN returns the number of non-NaN entries in a float64 slice.
// Used to report bh_total_tested in the score-table Finding so a reviewer
// can verify (n_tested, n_surviving) against alpha.
func countNonNaN(xs []float64) int {
	n := 0
	for _, x := range xs {
		if !math.IsNaN(x) {
			n++
		}
	}
	return n
}
