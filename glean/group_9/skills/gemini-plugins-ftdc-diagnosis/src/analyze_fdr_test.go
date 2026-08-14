package main

import (
	"math"
	"testing"
)

// TestBHKnownReferenceValues anchors the implementation against the
// classic worked example. Inputs: 10 p-values from BH-1995 Table 1
// (Benjamini & Hochberg 1995, "Controlling the False Discovery Rate"),
// expected critical points known from the literature.
func TestBHKnownReferenceValues(t *testing.T) {
	p := []float64{0.0001, 0.0004, 0.0019, 0.0095, 0.0201, 0.0278, 0.0298, 0.0344, 0.0459, 0.3240}
	q := benjaminiHochberg(p)
	// Expected: q_(i) = p_(i) * 10 / i, monotonized.
	// p_(10) = 0.3240 → raw 0.3240
	// p_(9)  = 0.0459 → raw 0.051 → 0.051 (< 0.3240)
	// p_(8)  = 0.0344 → raw 0.043
	// p_(7)  = 0.0298 → raw 0.04257 → wait recompute: 0.0298 * 10/7 = 0.04257
	// p_(6)  = 0.0278 → raw 0.0463
	// p_(5)  = 0.0201 → raw 0.0402
	// p_(4)  = 0.0095 → raw 0.02375
	// p_(3)  = 0.0019 → raw 0.00633
	// p_(2)  = 0.0004 → raw 0.002
	// p_(1)  = 0.0001 → raw 0.001
	// Then monotonize right→left (each is min of itself and all to its right).
	// Raw values at indices 4..7 are: 0.0402, 0.04633, 0.04257, 0.043.
	// Walking right-to-left: q_7=0.043, q_6=min(0.04257,0.043)=0.04257,
	// q_5=min(0.04633,0.04257)=0.04257, q_4=min(0.0402,0.04257)=0.0402.
	// So q_5 = 0.04257 (NOT 0.0402 — the envelope pulls raw_5 DOWN to
	// match its right neighbor's smaller q).
	expected := []float64{0.001, 0.002, 0.00633, 0.02375, 0.0402, 0.04257, 0.04257, 0.043, 0.051, 0.324}
	if len(q) != len(expected) {
		t.Fatalf("q length %d != expected %d", len(q), len(expected))
	}
	for i := range expected {
		if math.Abs(q[i]-expected[i]) > 0.001 {
			t.Errorf("q[%d] = %.6f, expected ≈ %.6f (input p=%.4f)", i, q[i], expected[i], p[i])
		}
	}
}

func TestBHMonotonicityEnforced(t *testing.T) {
	// Adversarial input: raw BH ratios are non-monotone before the
	// envelope step. Confirm output is non-decreasing in sorted-p order.
	p := []float64{0.001, 0.05, 0.04, 0.03}
	q := benjaminiHochberg(p)
	// Sort by p, walk through, check non-decreasing.
	// Sorted order: p=0.001 (idx 0), p=0.03 (idx 3), p=0.04 (idx 2), p=0.05 (idx 1).
	// q values in sorted order should be non-decreasing.
	sortedQ := []float64{q[0], q[3], q[2], q[1]}
	for i := 1; i < len(sortedQ); i++ {
		if sortedQ[i] < sortedQ[i-1] {
			t.Errorf("non-monotone at sorted idx %d: %.6f < %.6f", i, sortedQ[i], sortedQ[i-1])
		}
	}
}

func TestBHNaNPreserved(t *testing.T) {
	p := []float64{0.01, math.NaN(), 0.5}
	q := benjaminiHochberg(p)
	if !math.IsNaN(q[1]) {
		t.Fatalf("q[1] should be NaN (input was NaN), got %.6f", q[1])
	}
	// q[0] and q[2] should be computed over the remaining m=2 tests.
	// p=0.01 → raw 0.01 * 2 / 1 = 0.02 → monotonized 0.02
	// p=0.5  → raw 0.5  * 2 / 2 = 0.5  → 0.5
	if math.Abs(q[0]-0.02) > 1e-9 {
		t.Errorf("q[0] = %.9f, expected 0.02", q[0])
	}
	if math.Abs(q[2]-0.5) > 1e-9 {
		t.Errorf("q[2] = %.9f, expected 0.5", q[2])
	}
}

func TestBHEmpty(t *testing.T) {
	q := benjaminiHochberg(nil)
	if len(q) != 0 {
		t.Errorf("empty input → empty output, got len %d", len(q))
	}
}

func TestBHClampsInputs(t *testing.T) {
	// Bad inputs should be clamped, not propagate as NaN or break the
	// monotone envelope.
	p := []float64{-0.1, 1.5, 0.5}
	q := benjaminiHochberg(p)
	for i, v := range q {
		if math.IsNaN(v) || v < 0 || v > 1 {
			t.Errorf("q[%d] = %.6f out of [0,1]", i, v)
		}
	}
}

func TestBHSignificanceThresholdSanity(t *testing.T) {
	// Construct a case where only the smallest two p-values pass FDR=0.05.
	// m=4, alpha=0.05. Cuts: 0.0125, 0.025, 0.0375, 0.05.
	// p sorted: [0.001, 0.01, 0.04, 0.5]
	// 0.001 <= 0.0125 ✓
	// 0.01 <= 0.025 ✓
	// 0.04 <= 0.0375 ✗
	// 0.5 <= 0.05 ✗
	// Best surviving p: 0.01.
	thr := bhSignificanceThreshold([]float64{0.001, 0.01, 0.04, 0.5}, 0.05)
	if math.Abs(thr-0.01) > 1e-9 {
		t.Errorf("threshold = %.6f, expected 0.01", thr)
	}
}

func TestCountSignificantBH(t *testing.T) {
	q := []float64{0.001, 0.04, 0.06, math.NaN(), 0.049}
	n := countSignificantBH(q, 0.05)
	if n != 3 {
		t.Errorf("countSignificantBH = %d, expected 3 (0.001, 0.04, 0.049)", n)
	}
}
