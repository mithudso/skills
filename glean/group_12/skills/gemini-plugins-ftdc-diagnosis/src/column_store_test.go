package main

import (
	"math"
	"math/rand"
	"testing"
)

// roundTrip appends vals to a fresh CompressedColumn, calls Finish, and
// returns ReadAllInt64 of the result. The returned slice should equal vals.
func roundTrip(t *testing.T, vals []int64) []int64 {
	t.Helper()
	c := NewCompressedColumn(len(vals))
	for _, v := range vals {
		c.Add(v, true)
	}
	c.Finish()
	got := c.ReadAllInt64(nil)
	if len(got) != len(vals) {
		t.Fatalf("len mismatch: got %d, want %d", len(got), len(vals))
	}
	for i, v := range vals {
		if got[i] != v {
			t.Fatalf("vals[%d]: got %d, want %d", i, got[i], v)
		}
	}
	return got
}

func TestColumnEmpty(t *testing.T) {
	c := NewCompressedColumn(0)
	c.Finish()
	if c.Len() != 0 {
		t.Fatalf("expected Len 0, got %d", c.Len())
	}
	got := c.ReadAllInt64(nil)
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestColumnSingleSample(t *testing.T) {
	roundTrip(t, []int64{42})
}

func TestColumnConstantBelow64(t *testing.T) {
	vals := make([]int64, 50)
	for i := range vals {
		vals[i] = 7
	}
	roundTrip(t, vals)
}

func TestColumnConstantExactlyOneBlock(t *testing.T) {
	vals := make([]int64, blockSize)
	for i := range vals {
		vals[i] = 100
	}
	roundTrip(t, vals)
}

func TestColumnConstantMultipleBlocks(t *testing.T) {
	vals := make([]int64, 3*blockSize+13)
	for i := range vals {
		vals[i] = -5
	}
	roundTrip(t, vals)
}

func TestColumnMonotonicCounter(t *testing.T) {
	vals := make([]int64, 200)
	v := int64(0)
	for i := range vals {
		v += int64(i % 7)
		vals[i] = v
	}
	roundTrip(t, vals)
}

func TestColumnSignedGauge(t *testing.T) {
	vals := []int64{0, -1, 2, -3, 4, -5, 100, -100, 0, 0, 0, 50, 50, 51, 49}
	roundTrip(t, vals)
}

func TestColumnLargeDeltas(t *testing.T) {
	// Force wide-fallback path: deltas near max int64.
	vals := []int64{
		0,
		math.MaxInt64 / 2,
		math.MinInt64 / 2,
		math.MaxInt64 / 3,
		1 << 60,
		-(1 << 60),
	}
	roundTrip(t, vals)
}

func TestColumnRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for trial := 0; trial < 50; trial++ {
		n := rng.Intn(500) + 1
		vals := make([]int64, n)
		v := int64(rng.Int63n(1000))
		for i := range vals {
			step := int64(rng.Intn(20)) - 10
			v += step
			vals[i] = v
		}
		roundTrip(t, vals)
	}
}

func TestColumnReadRange(t *testing.T) {
	vals := make([]int64, 500)
	for i := range vals {
		vals[i] = int64(i * 3)
	}
	c := NewCompressedColumn(500)
	for _, v := range vals {
		c.Add(v, true)
	}
	c.Finish()

	// Various ranges.
	cases := []struct{ start, end int }{
		{0, 1},
		{0, 64},
		{0, 65},
		{64, 128},
		{63, 65},
		{100, 200},
		{400, 500},
		{0, 500},
		{499, 500},
	}
	for _, tc := range cases {
		got := c.ReadRange(tc.start, tc.end, nil)
		if len(got) != tc.end-tc.start {
			t.Fatalf("range [%d,%d): len %d, want %d", tc.start, tc.end, len(got), tc.end-tc.start)
		}
		for i, v := range got {
			if v != vals[tc.start+i] {
				t.Fatalf("range [%d,%d)[%d]: got %d, want %d", tc.start, tc.end, i, v, vals[tc.start+i])
			}
		}
	}
}

func TestColumnReadRangeStraddlesPending(t *testing.T) {
	// nFlushed = 64, pendN = 5 → nTotal = 69.
	c := NewCompressedColumn(70)
	for i := 0; i < 69; i++ {
		c.Add(int64(i*2), true)
	}
	// Don't Finish — query while last block is still pending.

	got := c.ReadRange(60, 69, nil)
	if len(got) != 9 {
		t.Fatalf("got %d, want 9", len(got))
	}
	for i, v := range got {
		want := int64((60 + i) * 2)
		if v != want {
			t.Fatalf("got %d, want %d at idx %d", v, want, i)
		}
	}
}

func TestColumnReadAllFloat(t *testing.T) {
	vals := []int64{0, 1, 2, 100, -50, 1 << 30}
	c := NewCompressedColumn(len(vals))
	for _, v := range vals {
		c.Add(v, true)
	}
	c.Finish()
	got := c.ReadAll(nil)
	if len(got) != len(vals) {
		t.Fatalf("len mismatch")
	}
	for i, v := range vals {
		if got[i] != float64(v) {
			t.Fatalf("idx %d: got %v, want %v", i, got[i], float64(v))
		}
	}
}

func TestColumnOnlineStats(t *testing.T) {
	c := NewCompressedColumn(0)
	vals := []int64{10, 20, 5, 30, 15}
	for _, v := range vals {
		c.Add(v, true)
	}
	c.Finish()

	if c.ObservedMin() != 5 {
		t.Errorf("min: got %d, want 5", c.ObservedMin())
	}
	if c.ObservedMax() != 30 {
		t.Errorf("max: got %d, want 30", c.ObservedMax())
	}
	wantMean := (10 + 20 + 5 + 30 + 15) / 5.0
	if math.Abs(c.Mean()-wantMean) > 1e-9 {
		t.Errorf("mean: got %v, want %v", c.Mean(), wantMean)
	}
	if c.IsCounter() {
		t.Errorf("isCounter: got true, want false (20 -> 5 is a decrease)")
	}
	if c.IsDegenerate() {
		t.Errorf("isDegenerate: got true, want false (5 != 30)")
	}
}

func TestColumnIsCounterTrue(t *testing.T) {
	c := NewCompressedColumn(0)
	for v := int64(0); v < 100; v++ {
		c.Add(v, true)
	}
	c.Finish()
	if !c.IsCounter() {
		t.Errorf("expected isCounter=true for monotonic series")
	}
}

func TestColumnIsDegenerate(t *testing.T) {
	c := NewCompressedColumn(0)
	for i := 0; i < 50; i++ {
		c.Add(7, true)
	}
	c.Finish()
	if !c.IsDegenerate() {
		t.Errorf("expected isDegenerate=true for constant series")
	}
	if !c.IsCounter() {
		t.Errorf("constant series is also counter (no decrease)")
	}
}

func TestColumnRateGapReset(t *testing.T) {
	c := NewCompressedColumn(0)
	// Two consecutive ramps separated by a gap. Inside each ramp,
	// per-sample delta = 1, so rate Welford mean = 1.
	// Across the gap, no rate observation is recorded (the
	// 100->1000 jump must NOT inflate rateMean).
	for v := int64(0); v < 5; v++ {
		c.Add(v, v != 0) // first sample has no prev; rest consecutive
	}
	// Gap.
	for v := int64(1000); v < 1005; v++ {
		c.Add(v, v != 1000) // first sample after gap: consecutive=false
	}
	c.Finish()
	// Each ramp contributes 3 rate observations (samples 1-2, 2-3, 3-4
	// within a ramp; the first sample has no prev). The 100->1000 jump
	// is NOT recorded because the gap-sample's consecutive=false resets
	// hasPrev. So total = 3 + 3 = 6.
	if c.rateN != 6 {
		t.Errorf("rateN: got %d, want 6", c.rateN)
	}
	// All deltas == 1.
	if math.Abs(c.RateMean()-1.0) > 1e-9 {
		t.Errorf("rateMean: got %v, want 1.0", c.RateMean())
	}
}

func TestColumnCompressionRatio(t *testing.T) {
	// Constant series: should compress to ~1 byte of packed payload per
	// block (in fact ZERO with bitWidth=0).
	vals := make([]int64, 10000)
	for i := range vals {
		vals[i] = 42
	}
	c := NewCompressedColumn(len(vals))
	for _, v := range vals {
		c.Add(v, true)
	}
	c.Finish()
	if c.CompressedBytes() != 0 {
		t.Errorf("constant series should have 0 packed bytes, got %d", c.CompressedBytes())
	}

	// Monotonic counter with small deltas: should compress well.
	c2 := NewCompressedColumn(len(vals))
	v := int64(0)
	for i := 0; i < len(vals); i++ {
		v += int64(i % 5)
		c2.Add(v, true)
	}
	c2.Finish()
	rawBytes := 8 * len(vals)
	if c2.CompressedBytes() >= rawBytes/2 {
		t.Errorf("monotonic counter should compress ≥2×; got %d bytes for %d raw", c2.CompressedBytes(), rawBytes)
	}
}
