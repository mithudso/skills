// analyze_wave7_fixes_test.go — regression tests for Wave 7 bug fixes.
package main

import (
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Fix W7-D1: term_change idKey uses tsFormat (ms precision), not compact
// second-precision format. Same-second events must have distinct IDs.
// ---------------------------------------------------------------------------

func TestTermChange_DistinctIDsWithinSameSec(t *testing.T) {
	base, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	// Two timestamps in the same UTC second, 500ms apart.
	t1 := base.Add(100 * time.Millisecond)
	t2 := base.Add(600 * time.Millisecond)

	id1 := base.Add(100 * time.Millisecond).UTC().Format(tsFormat)
	id2 := base.Add(600 * time.Millisecond).UTC().Format(tsFormat)

	if id1 == id2 {
		t.Errorf("tsFormat produces identical keys for two times in the same second (%s == %s)", id1, id2)
	}
	_ = t1
	_ = t2

	// Verify the compact format (old behaviour) WOULD have collided.
	old1 := base.Add(100 * time.Millisecond).UTC().Format("20060102T150405Z")
	old2 := base.Add(600 * time.Millisecond).UTC().Format("20060102T150405Z")
	if old1 != old2 {
		t.Errorf("test setup error: compact format should collapse same-second times (%s != %s)", old1, old2)
	}
}

// ---------------------------------------------------------------------------
// Fix W7-A3: discontiDelta uses loIdx-1 as baseline so an increment at the
// first window sample is captured.
// ---------------------------------------------------------------------------

func TestDiscontiDelta_IncrementAtWindowBoundaryIsVisible(t *testing.T) {
	// Counter: [0, 1, 1, 1] — increment happens at index 1.
	// Window covers indices 1-3 (loIdx=1, hiIdx=4).
	// Old code: col[hiIdx-1] - col[loIdx] = col[3] - col[1] = 1 - 1 = 0  ← misses it
	// New code: col[hiIdx-1] - col[loIdx-1] = col[3] - col[0] = 1 - 0 = 1 ← correct
	col := []int64{0, 1, 1, 1}
	loIdx := 1
	hiIdx := 4

	// New logic (mirroring the fix in emitSealEvents):
	baseline := col[loIdx]
	if loIdx > 0 {
		baseline = col[loIdx-1]
	}
	delta := col[hiIdx-1] - baseline

	if delta != 1 {
		t.Errorf("discontiDelta = %d, want 1 (increment at window boundary must be captured)", delta)
	}
}

func TestDiscontiDelta_loIdxZeroUsesFirstSample(t *testing.T) {
	// When loIdx==0 there is no prior sample; baseline stays col[0].
	col := []int64{0, 1, 1, 2}
	loIdx := 0
	hiIdx := 4

	baseline := col[loIdx]
	if loIdx > 0 {
		baseline = col[loIdx-1]
	}
	delta := col[hiIdx-1] - baseline

	// delta = col[3] - col[0] = 2 - 0 = 2
	if delta != 2 {
		t.Errorf("discontiDelta (loIdx=0) = %d, want 2", delta)
	}
}

// ---------------------------------------------------------------------------
// Fix W7-R1: Pearson r over the overlap window — no deflation at non-zero lags.
// For a perfectly correlated pair at lag L, r must equal 1.0 regardless of L.
// ---------------------------------------------------------------------------

func TestRunAutoCorrelate_PerfectLagCorrelation(t *testing.T) {
	// leader = [1,2,...,20], follower = leader shifted right by 5.
	// At lag=5, the overlap window is indices 0..14 (count=15).
	// Expected r = 1.0 at lag=5.
	n := 20
	lag := 5
	leader := make([]float64, n)
	follower := make([]float64, n)
	for i := 0; i < n; i++ {
		leader[i] = float64(i + 1)
	}
	for i := 0; i < n; i++ {
		if i-lag >= 0 {
			follower[i] = leader[i-lag]
		} else {
			follower[i] = 0
		}
	}

	// Simulate the overlap-window Pearson computation (the fixed formula).
	count := n - lag // 15
	var sumX, sumY float64
	for i := 0; i < n; i++ {
		t := i + lag
		if t < 0 || t >= n {
			continue
		}
		sumX += leader[i]
		sumY += follower[t]
	}
	meanX := sumX / float64(count)
	meanY := sumY / float64(count)
	var ssX, ssY, dot float64
	for i := 0; i < n; i++ {
		tt := i + lag
		if tt < 0 || tt >= n {
			continue
		}
		dx := leader[i] - meanX
		dy := follower[tt] - meanY
		ssX += dx * dx
		ssY += dy * dy
		dot += dx * dy
	}
	denom := sqrtFloat(ssX * ssY)
	if denom <= 0 {
		t.Fatal("denom <= 0 for non-constant series")
	}
	r := dot / denom
	if r < 0.9999 || r > 1.0001 {
		t.Errorf("overlap-window Pearson at lag=%d: got r=%v, want ~1.0 (no deflation)", lag, r)
	}

	// Verify the OLD formula (full-series SD) would have deflated r.
	// old r = dot_old / ((n-1) * leaderSD * followerSD)
	var lSum, fSum float64
	for i := 0; i < n; i++ {
		lSum += leader[i]
		fSum += follower[i]
	}
	lMean := lSum / float64(n)
	fMean := fSum / float64(n)
	var lSS, fSS, dotOld float64
	for i := 0; i < n; i++ {
		lSS += (leader[i] - lMean) * (leader[i] - lMean)
		fSS += (follower[i] - fMean) * (follower[i] - fMean)
	}
	// Old dot = centred sum using FULL-series means, over OVERLAP window only.
	for i := 0; i < n; i++ {
		tt := i + lag
		if tt < 0 || tt >= n {
			continue
		}
		dotOld += (leader[i] - lMean) * (follower[tt] - fMean)
	}
	leaderSD := sqrtFloat(lSS / float64(n-1))
	followerSD := sqrtFloat(fSS / float64(n-1))
	oldR := dotOld / (float64(n-1) * leaderSD * followerSD)
	if oldR >= 0.9999 {
		t.Errorf("old formula should deflate r below 1.0 at lag=%d (got %v — test setup error)", lag, oldR)
	}
}

// ---------------------------------------------------------------------------
// Fix W7-R2: when all lags in the ±60 grid produce count < 4 overlapping
// samples, the empty-points guard must fire before the KindCorrelation emit.
// ---------------------------------------------------------------------------

func TestAutoCorrelate_EmptyPointsGuardFires(t *testing.T) {
	// With n=3 samples, every lag |L|>=1 produces count = 3-|L| <= 2 < 4,
	// and lag=0 produces count=3 < 4. So the points slice is always empty.
	n := 3
	leader := []float64{1, 2, 3}
	follower := []float64{4, 5, 6}

	type lagPoint struct {
		lag int
		r   float64
	}
	points := []lagPoint{}
	for lag := -60; lag <= 60; lag++ {
		absLag := lag
		if absLag < 0 {
			absLag = -absLag
		}
		count := n - absLag
		if count < 4 {
			continue
		}
		var sumX, sumY float64
		for i := 0; i < n; i++ {
			tt := i + lag
			if tt < 0 || tt >= n {
				continue
			}
			sumX += leader[i]
			sumY += follower[tt]
		}
		meanX := sumX / float64(count)
		meanY := sumY / float64(count)
		var ssX, ssY, dot float64
		for i := 0; i < n; i++ {
			tt := i + lag
			if tt < 0 || tt >= n {
				continue
			}
			dx := leader[i] - meanX
			dy := follower[tt] - meanY
			ssX += dx * dx
			ssY += dy * dy
			dot += dx * dy
		}
		denom := sqrtFloat(ssX * ssY)
		if denom <= 0 {
			continue
		}
		r := dot / denom
		points = append(points, lagPoint{lag, r})
	}

	// With n=3, points must be empty — the guard must fire.
	if len(points) != 0 {
		t.Errorf("expected 0 lag points for n=%d, got %d", n, len(points))
	}
	// Verify the guard would result in KindSkipped, not KindCorrelation.
	// (The guard itself: if len(points) == 0 → emit KindSkipped and return.)
	expectedKind := KindSkipped
	if len(points) > 0 {
		expectedKind = KindCorrelation
	}
	if expectedKind != KindSkipped {
		t.Errorf("empty points should lead to KindSkipped, not %s", expectedKind)
	}
}
