package main

import (
	"testing"
)

// TestSortedIDsByAbsScoreWithNamesStable verifies that when scores are tied,
// the IDs are sorted by name (alphabetically) in a stable deterministic order.
func TestSortedIDsByAbsScoreWithNamesStable(t *testing.T) {
	// Create tied scores: all have abs(score) = 5.0
	scores := []float64{5.0, -5.0, 5.0, -5.0, 5.0}
	names := []string{"delta", "alpha", "charlie", "bravo", "echo"}

	result := sortedIDsByAbsScoreWithNames(scores, names)

	// All have the same absolute score, so secondary sort by name should apply.
	// Expected order: alpha (idx 1), bravo (idx 3), charlie (idx 2), delta (idx 0), echo (idx 4)
	expected := []int{1, 3, 2, 0, 4}

	if len(result) != len(expected) {
		t.Fatalf("result length %d != expected length %d", len(result), len(expected))
	}

	for i, exp := range expected {
		if result[i] != exp {
			t.Errorf("result[%d] = %d, expected %d (name=%q vs %q)",
				i, result[i], exp, names[result[i]], names[exp])
		}
	}
}

// TestSortedIDsByAbsScoreOrderByScore verifies that higher absolute
// scores come first.
func TestSortedIDsByAbsScoreOrderByScore(t *testing.T) {
	scores := []float64{3.0, -5.0, 2.0, 5.0, 1.0}
	names := []string{"a", "b", "c", "d", "e"}

	result := sortedIDsByAbsScoreWithNames(scores, names)

	// Expected order by abs score (descending):
	// abs(5.0)=5.0 [idx 3], abs(-5.0)=5.0 [idx 1], abs(3.0)=3.0 [idx 0],
	// abs(2.0)=2.0 [idx 2], abs(1.0)=1.0 [idx 4]
	// When abs scores are equal (5.0), sort by name: "b" < "d" so idx 1 before 3.
	// Result: [1, 3, 0, 2, 4]
	expected := []int{1, 3, 0, 2, 4}

	if len(result) != len(expected) {
		t.Fatalf("result length %d != expected length %d", len(result), len(expected))
	}

	for i, exp := range expected {
		if result[i] != exp {
			t.Errorf("result[%d] = %d, expected %d (name=%q, score=%v vs name=%q, score=%v)",
				i, result[i], exp, names[result[i]], scores[result[i]], names[exp], scores[exp])
		}
	}
}
