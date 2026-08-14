package main

import (
	"math"
	"testing"
)

// TestTheilSenSlope_EvenPairCount verifies the even-length median fix.
//
// For a 4-element series the pairCount = 4*3/2 = 6 slopes (even).
// The correct median is the average of slopes[2] and slopes[3] (0-indexed
// middle two). The old code used slopes[6/2] = slopes[3], which is the upper
// middle element — off by one for even-length slices.
//
// Construct a series where the true underlying slope is exactly 1.5, and
// arrange the 6 pairwise slopes so the average of the two middle elements
// equals 1.5 exactly. Then confirm the fixed code returns 1.5 and that the
// old code (using upper-middle index) would NOT.
//
// Series: [0, 1, 3, 6]
// Pairwise slopes (i < j):
//   (0,1): (1-0)/(1-0)   = 1.0
//   (0,2): (3-0)/(2-0)   = 1.5
//   (0,3): (6-0)/(3-0)   = 2.0
//   (1,2): (3-1)/(2-1)   = 2.0
//   (1,3): (6-1)/(3-1)   = 2.5
//   (2,3): (6-3)/(3-2)   = 3.0
// Sorted: [1.0, 1.5, 2.0, 2.0, 2.5, 3.0]
// Even median (average of [2] and [3]): (2.0 + 2.0) / 2 = 2.0
// Old code median (upper middle [3]):    2.0
// For this particular series both happen to agree. Use a different series:
//
// Series: [0, 2, 3, 6]
// Pairwise slopes:
//   (0,1): 2/1 = 2.0
//   (0,2): 3/2 = 1.5
//   (0,3): 6/3 = 2.0
//   (1,2): 1/1 = 1.0
//   (1,3): 4/2 = 2.0
//   (2,3): 3/1 = 3.0
// Sorted: [1.0, 1.5, 2.0, 2.0, 2.0, 3.0]
// Even median: (2.0 + 2.0) / 2 = 2.0
// Old upper-middle: slopes[3] = 2.0 — still the same.
//
// Use a series where the two middle slopes differ:
// Series: [0, 1, 4, 6]
// Pairwise slopes:
//   (0,1): 1/1 = 1.0
//   (0,2): 4/2 = 2.0
//   (0,3): 6/3 = 2.0
//   (1,2): 3/1 = 3.0
//   (1,3): 5/2 = 2.5
//   (2,3): 2/1 = 2.0
// Sorted: [1.0, 2.0, 2.0, 2.0, 2.5, 3.0]
// Even median: (2.0+2.0)/2 = 2.0; old code: slopes[3] = 2.0 — same again.
//
// We need slopes[mid-1] != slopes[mid] to distinguish the two. Use:
// Series: [0, 1, 5, 6]
// Pairwise slopes:
//   (0,1): 1/1 = 1.0
//   (0,2): 5/2 = 2.5
//   (0,3): 6/3 = 2.0
//   (1,2): 4/1 = 4.0
//   (1,3): 5/2 = 2.5
//   (2,3): 1/1 = 1.0
// Sorted: [1.0, 1.0, 2.0, 2.5, 2.5, 4.0]
// Even median: (2.0 + 2.5) / 2 = 2.25
// Old code (upper-middle): slopes[3] = 2.5  ← differs from correct 2.25
func TestTheilSenSlope_EvenPairCount(t *testing.T) {
	series := []int64{0, 1, 5, 6}
	slope, _, _ := theilSenSlope(series)

	const want = 2.25
	if math.Abs(slope-want) > 1e-9 {
		t.Errorf("theilSenSlope even-median: got %v, want %v (average of two middle slopes)", slope, want)
	}
}

// TestTheilSenSlope_OddPairCount verifies that odd-length slope lists are
// unaffected by the even-median fix (upper-middle == true median for odd).
//
// Series: [0, 2, 6] (n=3, pairCount=3 — odd)
// Pairwise slopes:
//   (0,1): 2/1 = 2.0
//   (0,2): 6/2 = 3.0
//   (1,2): 4/1 = 4.0
// Sorted: [2.0, 3.0, 4.0]
// Median (index 1): 3.0
func TestTheilSenSlope_OddPairCount(t *testing.T) {
	series := []int64{0, 2, 6}
	slope, _, _ := theilSenSlope(series)
	const want = 3.0
	if math.Abs(slope-want) > 1e-9 {
		t.Errorf("theilSenSlope odd-median: got %v, want %v", slope, want)
	}
}

// TestTheilSenSlope_Intercept_EvenResid verifies the even-median fix is also
// applied to the intercept (residuals slice).
//
// Using series [0, 1, 5, 6] with computed slope=2.25:
// Residuals y_i - slope*i:
//   i=0: 0  - 2.25*0 = 0.0
//   i=1: 1  - 2.25*1 = -1.25
//   i=2: 5  - 2.25*2 = 0.5
//   i=3: 6  - 2.25*3 = -0.75
// Sorted: [-1.25, -0.75, 0.0, 0.5]
// Even median: (-0.75 + 0.0) / 2 = -0.375
// Old code (upper-middle index 2): 0.0  ← differs from correct -0.375
func TestTheilSenSlope_Intercept_EvenResid(t *testing.T) {
	series := []int64{0, 1, 5, 6}
	_, intercept, _ := theilSenSlope(series)
	const want = -0.375
	if math.Abs(intercept-want) > 1e-9 {
		t.Errorf("theilSenSlope intercept even-median: got %v, want %v", intercept, want)
	}
}

// TestTheilSenSlope_ConstantSeries verifies that a flat series produces
// slope=0 and signFraction=0.5.
func TestTheilSenSlope_ConstantSeries(t *testing.T) {
	series := make([]int64, 50)
	for i := range series {
		series[i] = 42
	}
	slope, _, signFrac := theilSenSlope(series)
	if slope != 0 {
		t.Errorf("constant series: slope = %v, want 0", slope)
	}
	if math.Abs(signFrac-0.5) > 1e-9 {
		t.Errorf("constant series: signFraction = %v, want 0.5", signFrac)
	}
}

// TestTheilSenSlope_PureRamp verifies a pure linear ramp gives slope=1.
func TestTheilSenSlope_PureRamp(t *testing.T) {
	const N = 100
	series := make([]int64, N)
	for i := range series {
		series[i] = int64(i)
	}
	slope, _, signFrac := theilSenSlope(series)
	if math.Abs(slope-1.0) > 1e-9 {
		t.Errorf("pure ramp: slope = %v, want 1.0", slope)
	}
	if signFrac < 0.99 {
		t.Errorf("pure ramp: signFraction = %v, want ≥0.99", signFrac)
	}
}
