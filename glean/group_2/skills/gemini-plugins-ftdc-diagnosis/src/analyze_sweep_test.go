package main

import (
	"math"
	"testing"
)

// TestSilhouette1D_SingleMemberCluster: when one cluster has only a
// single member, that member hits `continue` (ownCount==0). Before the
// fix, the skipped point still contributed to the denominator
// `float64(len(values))`, inflating it and deflating the score. After
// the fix, only actually-computed points are counted.
func TestSilhouette1D_SingleMemberCluster(t *testing.T) {
	// 3 points: two in cluster 0, one alone in cluster 1.
	// The lone point (index 2) hits `otherCount == 0` and is skipped.
	// Wait — in silhouette1D the lone-cluster member has ownCount==0
	// (no OTHER members in its own cluster), so it hits `continue`.
	// Only 2 of 3 points contribute to sum; computed should be 2.
	values := []float64{1.0, 2.0, 100.0}
	labels := []int{0, 0, 1} // cluster 0: {1,2}, cluster 1: {100}

	score := silhouette1D(values, labels)

	// The lone member in cluster 1 (value=100) has ownCount=0 → skipped.
	// Point 0 (v=1): own={2}, other={100}; a=1, b=99, s=(99-1)/99≈0.9899
	// Point 1 (v=2): own={1}, other={100}; a=1, b=98, s=(98-1)/98≈0.9898
	// Expected score ≈ (0.9899 + 0.9898) / 2 ≈ 0.9899 (computed=2, not 3)
	//
	// With the old bug (divide by len(values)=3):
	// score = (0.9899 + 0.9898) / 3 ≈ 0.6599 — wrong.
	want := (99.0/99.0 - 1.0/99.0 + 98.0/98.0 - 1.0/98.0) / 2.0
	// Simplify: point0 s=(99-1)/max(1,99)=98/99; point1 s=(98-1)/max(1,98)=97/98
	// Recalculate:
	// point0: a=|1-2|/1=1, b=|1-100|/1=99, denom=max(1,99)=99, s=(99-1)/99=98/99
	// point1: a=|2-1|/1=1, b=|2-100|/1=98, denom=max(1,98)=98, s=(98-1)/98=97/98
	wantExact := (98.0/99.0 + 97.0/98.0) / 2.0

	if math.Abs(score-wantExact) > 1e-9 {
		t.Errorf("silhouette1D = %v, want %v (single-member cluster inflated denominator)", score, wantExact)
	}

	// Additionally verify the score is NOT what the buggy formula would
	// produce (dividing by 3 instead of 2).
	buggyScore := (98.0/99.0 + 97.0/98.0) / 3.0
	if math.Abs(score-buggyScore) < 1e-9 {
		t.Errorf("silhouette1D returned buggy value %v (divided by total count 3 instead of computed count 2)", score)
	}
	_ = want
}

// TestSilhouette1D_AllSingletonClusters: if every point is in its own
// cluster, ownCount==0 for all → computed=0 → return 0, not NaN.
func TestSilhouette1D_AllSingletonClusters(t *testing.T) {
	values := []float64{1.0, 2.0, 3.0}
	labels := []int{0, 1, 2}

	score := silhouette1D(values, labels)
	if math.IsNaN(score) || math.IsInf(score, 0) {
		t.Errorf("expected 0 for all-singleton clusters, got %v", score)
	}
	if score != 0 {
		t.Errorf("expected score=0 for all-singleton clusters, got %v", score)
	}
}

// TestSilhouette1D_PerfectSeparation: two tight clusters far apart
// should produce a score close to 1.
func TestSilhouette1D_PerfectSeparation(t *testing.T) {
	// Cluster 0: {0,1,2}, Cluster 1: {100,101,102}
	values := []float64{0, 1, 2, 100, 101, 102}
	labels := []int{0, 0, 0, 1, 1, 1}

	score := silhouette1D(values, labels)
	if score < 0.9 {
		t.Errorf("expected near-1 silhouette for perfectly separated clusters, got %v", score)
	}
}
