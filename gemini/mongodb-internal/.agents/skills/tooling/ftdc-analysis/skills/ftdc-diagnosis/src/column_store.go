// column_store.go - Compressed columnar storage for FTDC metric series.
//
// CompressedColumn holds a metric's full 1Hz int64 series as a sequence of
// fixed-size blocks. Each block stores deltas relative to its first sample,
// zigzag-encoded and bit-packed to a per-block width. Constant runs are
// stored as a single base value (zero bytes of packed payload).
//
// Design goals:
//   1. Hold ~6500 metrics × ~170K samples (48h FTDC) in ~1.5 GB resident,
//      vs. ~8.8 GB raw int64.
//   2. Streaming append API for one-pass ingest from the chunk decoder.
//   3. Random-access ReadRange so analyzer stages can pull a window without
//      decoding the whole column.
//   4. Online min/max + Welford stats updated as values are appended; no
//      separate Phase 1 scan over decoded samples.
//
// Encoding (per block):
//   base      int64    value of the block's first sample
//   blockLen  uint16   number of samples in the block (≤ blockSize)
//   bitWidth  uint8    bits per zigzag-encoded delta (0..64). 0 means "all
//                      deltas zero — constant run; packed is empty"
//   flags     uint8    bit 0x01: wide-fallback (deltas stored as raw int64
//                      in packed, 8 bytes each). Used when zigzag bitWidth
//                      would round up to >57 and we'd rather just store raw.
//   packed    []byte   bit-packed zigzag deltas, MSB-first within bytes,
//                      OR raw int64 deltas when flags&0x01 set.
//   blockMin  int64    min value across the block
//   blockMax  int64    max value across the block
//
// All blocks except the last have len == blockSize. The last block may be
// shorter; its actual length is recorded in blockLen.
//
// Compression ratio expectations on real MongoDB FTDC:
//   - Constant/zero metrics:     50–200× (zero-run blocks)
//   - Monotonic counters:        6–10×  (small positive deltas)
//   - Gauges:                    3–6×   (small signed deltas)
//   - Bursty histograms:         4–8×
// Blended: ~5–7× for whole capture.

package main

import (
	"encoding/binary"
	"math"
	"math/bits"
)

// blockSize is the fixed sample count per block. 64 was chosen for:
//   - L1 cache friendliness (one block = 512B of decoded int64).
//   - Random access granularity finer than Welch baseline/stress windows
//     (~10800 samples = ~169 blocks).
//   - Cheap block-level min/max scan for window aggregates without full
//     decode when the window covers many full blocks.
const blockSize = 64

// columnBlock is one compressed run of samples. Allocated inline in
// CompressedColumn.blocks; one slice header + the packed []byte.
type columnBlock struct {
	base     int64
	blockLen uint16
	bitWidth uint8
	flags    uint8
	packed   []byte
	blockMin int64
	blockMax int64
}

const blockFlagWide = 0x01

// CompressedColumn is one metric's full series. Built by repeated Add() calls
// from the chunk decoder. Block flushes happen automatically every blockSize
// samples; Finish() flushes a final short block if needed.
//
// Online statistics (welfordMean, welfordM2, observedMin, observedMax,
// isCounter, isDegenerate) are updated on every Add() so the analyzer can
// read them at end-of-stream without scanning the column again.
type CompressedColumn struct {
	blocks []columnBlock
	nTotal int // total samples appended (across all blocks plus any in pendBuf)

	// Pending values not yet flushed to a block. Length is in [0, blockSize).
	pendBuf [blockSize]int64
	pendN   int

	// Online statistics over the full appended series so far.
	observedMin int64
	observedMax int64
	welfordMean float64
	welfordM2   float64
	welfordN    int
	isMonotonic bool // no decrease seen yet ⇒ candidate counter
	hasAny      bool // true once at least one value has been Add'd
	lastValue   int64

	// Rate-side Welford (counter delta interpretation). Updated only when
	// the caller signals consecutive samples (no FTDC gap). Used by
	// analysis stages that need counter rates.
	rateMean float64
	rateM2   float64
	rateN    int
	hasPrev  bool
}

// NewCompressedColumn allocates a new column. Initial capacity hint helps
// avoid reallocating the blocks slice early; pass an estimate of expected
// sample count to size the backing slice.
func NewCompressedColumn(capacityHint int) *CompressedColumn {
	blocksHint := capacityHint / blockSize
	if blocksHint < 4 {
		blocksHint = 4
	}
	return &CompressedColumn{
		blocks:      make([]columnBlock, 0, blocksHint),
		observedMin: math.MaxInt64,
		observedMax: math.MinInt64,
		isMonotonic: true,
	}
}

// Add appends one sample to the column. consecutive=true means this sample
// follows the previous one with no FTDC gap, so counter-rate accumulators
// may include the (v - prevValue) delta. consecutive=false (gap or first
// sample) skips the rate delta and resets hasPrev — mirrors the reset
// behavior at analyze_store.go:233–238.
func (c *CompressedColumn) Add(v int64, consecutive bool) {
	// Update value-side stats.
	if !c.hasAny {
		c.observedMin = v
		c.observedMax = v
		c.hasAny = true
	} else {
		if v < c.observedMin {
			c.observedMin = v
		}
		if v > c.observedMax {
			c.observedMax = v
		}
		if v < c.lastValue {
			c.isMonotonic = false
		}
	}

	// Welford on raw values (used for gauge interpretation).
	c.welfordN++
	delta := float64(v) - c.welfordMean
	c.welfordMean += delta / float64(c.welfordN)
	c.welfordM2 += delta * (float64(v) - c.welfordMean)

	// Welford on first-differences (counter rate interpretation).
	if consecutive && c.hasPrev {
		d := float64(v - c.lastValue)
		c.rateN++
		rd := d - c.rateMean
		c.rateMean += rd / float64(c.rateN)
		c.rateM2 += rd * (d - c.rateMean)
	}
	c.hasPrev = consecutive

	c.lastValue = v

	// Buffer the value and flush a block when full.
	c.pendBuf[c.pendN] = v
	c.pendN++
	c.nTotal++
	if c.pendN == blockSize {
		c.flushBlock()
	}
}

// AddPad appends v to the column WITHOUT updating Welford / min/max /
// rate / monotonic tracking. Used by streaming_store.go to zero-pad
// columns for samples where a metric was absent, while keeping the
// column's analytical statistics free of those synthetic values. Sets
// hasPrev to false so the next real Add doesn't include a delta against
// the pad value.
func (c *CompressedColumn) AddPad(v int64) {
	c.hasPrev = false
	c.lastValue = v
	c.pendBuf[c.pendN] = v
	c.pendN++
	c.nTotal++
	if c.pendN == blockSize {
		c.flushBlock()
	}
}

// Finish flushes any partial trailing block. Must be called once after the
// last Add(); subsequent Adds will start a fresh trailing block.
func (c *CompressedColumn) Finish() {
	if c.pendN > 0 {
		c.flushBlock()
	}
}

// Len returns the total number of samples appended so far (including any
// still in the pending buffer).
func (c *CompressedColumn) Len() int {
	return c.nTotal
}

// Mean returns the value-side Welford mean over the full appended series.
func (c *CompressedColumn) Mean() float64 {
	return c.welfordMean
}

// StdDev returns the value-side sample stddev (Bessel-corrected, N-1).
func (c *CompressedColumn) StdDev() float64 {
	return welfordSD(c.welfordN, c.welfordM2)
}

// RateMean returns the counter-rate Welford mean over consecutive sample
// pairs (skips pairs separated by FTDC gaps).
func (c *CompressedColumn) RateMean() float64 { return c.rateMean }

// RateStdDev returns the counter-rate sample stddev (Bessel-corrected, N-1).
func (c *CompressedColumn) RateStdDev() float64 {
	return welfordSD(c.rateN, c.rateM2)
}

// WelfordN returns the number of non-pad gauge samples accumulated.
func (c *CompressedColumn) WelfordN() int { return c.welfordN }

// RateN returns the number of rate samples (counter deltas) accumulated.
func (c *CompressedColumn) RateN() int { return c.rateN }

// ObservedMin / ObservedMax return the min/max value seen across the full
// appended series. Returns (0, 0) for empty columns.
func (c *CompressedColumn) ObservedMin() int64 {
	if !c.hasAny {
		return 0
	}
	return c.observedMin
}
func (c *CompressedColumn) ObservedMax() int64 {
	if !c.hasAny {
		return 0
	}
	return c.observedMax
}

// IsCounter returns true if no decrease was ever seen across the series
// (i.e. the series is monotonically non-decreasing). Matches the existing
// classification at analyze_store.go:162-169.
func (c *CompressedColumn) IsCounter() bool { return c.isMonotonic && c.hasAny }

// IsDegenerate returns true when min == max across the full series
// (constant counter or gauge). Mirrors the existing classification.
func (c *CompressedColumn) IsDegenerate() bool {
	return c.hasAny && c.observedMin == c.observedMax
}

// flushBlock encodes c.pendBuf[:c.pendN] as a new columnBlock and appends
// it to c.blocks. Called automatically when pendN reaches blockSize, or
// explicitly by Finish().
func (c *CompressedColumn) flushBlock() {
	if c.pendN == 0 {
		return
	}
	vals := c.pendBuf[:c.pendN]

	// Block min/max for window-aggregate fast paths.
	bMin := vals[0]
	bMax := vals[0]
	for _, v := range vals[1:] {
		if v < bMin {
			bMin = v
		}
		if v > bMax {
			bMax = v
		}
	}

	block := columnBlock{
		base:     vals[0],
		blockLen: uint16(c.pendN),
		blockMin: bMin,
		blockMax: bMax,
	}

	// One-sample block: nothing to pack.
	if c.pendN == 1 {
		c.blocks = append(c.blocks, block)
		c.pendN = 0
		return
	}

	// Compute zigzag deltas and find max bit width.
	var maxZigzag uint64
	allZero := true
	deltas := make([]uint64, c.pendN-1)
	prev := vals[0]
	for i := 1; i < c.pendN; i++ {
		d := vals[i] - prev
		prev = vals[i]
		z := uint64(d<<1) ^ uint64(d>>63)
		deltas[i-1] = z
		if z != 0 {
			allZero = false
		}
		if z > maxZigzag {
			maxZigzag = z
		}
	}

	if allZero {
		// Constant run: all deltas zero. Bit width 0 means "no packed
		// payload, every sample equals base".
		block.bitWidth = 0
		c.blocks = append(c.blocks, block)
		c.pendN = 0
		return
	}

	bw := uint8(bits.Len64(maxZigzag))
	if bw > 57 {
		// Wide-fallback: raw int64 deltas. Cheaper than bit-packing at
		// 58–64 bits, and avoids edge cases in the packer.
		block.flags |= blockFlagWide
		block.bitWidth = 64
		packed := make([]byte, 8*(c.pendN-1))
		prev = vals[0]
		for i := 1; i < c.pendN; i++ {
			d := uint64(vals[i] - prev)
			binary.LittleEndian.PutUint64(packed[(i-1)*8:], d)
			prev = vals[i]
		}
		block.packed = packed
		c.blocks = append(c.blocks, block)
		c.pendN = 0
		return
	}

	block.bitWidth = bw
	// Pack low bw bits of each delta, MSB-first within bytes.
	totalBits := int(bw) * (c.pendN - 1)
	packed := make([]byte, (totalBits+7)/8)
	bitOffset := 0
	for _, z := range deltas {
		// Write bw bits of z at bitOffset.
		for k := int(bw) - 1; k >= 0; k-- {
			if (z>>uint(k))&1 == 1 {
				byteIdx := bitOffset >> 3
				bitIdx := 7 - (bitOffset & 7)
				packed[byteIdx] |= 1 << uint(bitIdx)
			}
			bitOffset++
		}
	}
	block.packed = packed

	c.blocks = append(c.blocks, block)
	c.pendN = 0
}

// decodeBlock writes block's samples into dst. dst must have length >= blockLen.
// Returns the number of samples written.
func (c *CompressedColumn) decodeBlock(blockIdx int, dst []int64) int {
	b := &c.blocks[blockIdx]
	n := int(b.blockLen)
	if n == 0 {
		return 0
	}
	dst[0] = b.base
	if n == 1 {
		return 1
	}

	if b.bitWidth == 0 {
		// Constant run.
		for i := 1; i < n; i++ {
			dst[i] = b.base
		}
		return n
	}

	if b.flags&blockFlagWide != 0 {
		// Wide-fallback: raw int64 deltas.
		prev := b.base
		for i := 1; i < n; i++ {
			d := int64(binary.LittleEndian.Uint64(b.packed[(i-1)*8:]))
			prev += d
			dst[i] = prev
		}
		return n
	}

	bw := uint(b.bitWidth)
	mask := (uint64(1) << bw) - 1
	prev := b.base
	bitOffset := 0
	for i := 1; i < n; i++ {
		// Read bw bits at bitOffset, MSB-first.
		var z uint64
		for k := uint(0); k < bw; k++ {
			byteIdx := bitOffset >> 3
			bitIdx := 7 - (bitOffset & 7)
			bit := uint64((b.packed[byteIdx] >> uint(bitIdx)) & 1)
			z = (z << 1) | bit
			bitOffset++
		}
		z &= mask
		// Zigzag decode.
		d := int64(z>>1) ^ -int64(z&1)
		prev += d
		dst[i] = prev
	}
	return n
}

// ReadRange decodes samples in [start, end) into dst, growing dst as needed.
// Returns the (possibly grown) dst slice trimmed to the actual count.
//
// Pre-included pending values (samples in c.pendBuf that haven't been
// flushed yet) are read from the pending buffer.
func (c *CompressedColumn) ReadRange(start, end int, dst []int64) []int64 {
	if start < 0 {
		start = 0
	}
	if end > c.nTotal {
		end = c.nTotal
	}
	if start >= end {
		return dst[:0]
	}
	out := dst
	if cap(out) < end-start {
		out = make([]int64, end-start)
	} else {
		out = out[:end-start]
	}

	// Flushed-block region: samples [0, nFlushed).
	nFlushed := c.nTotal - c.pendN

	// Reader cursor: where the next output sample goes.
	writeIdx := 0

	if start < nFlushed {
		// Find the block containing 'start'.
		blockIdx := start / blockSize
		blockStartSample := blockIdx * blockSize
		// Decode whole blocks one at a time into a small scratch buffer.
		var scratch [blockSize]int64
		for blockIdx < len(c.blocks) && blockStartSample < end {
			n := c.decodeBlock(blockIdx, scratch[:])
			// Compute slice of this block to copy.
			localFrom := 0
			if blockStartSample < start {
				localFrom = start - blockStartSample
			}
			localTo := n
			if blockStartSample+n > end {
				localTo = end - blockStartSample
			}
			if localTo < localFrom {
				localTo = localFrom
			}
			copyN := localTo - localFrom
			if copyN > 0 {
				copy(out[writeIdx:writeIdx+copyN], scratch[localFrom:localTo])
				writeIdx += copyN
			}
			blockStartSample += n
			blockIdx++
			if blockStartSample >= end {
				break
			}
		}
	}

	// Pending region: samples [nFlushed, nTotal).
	if end > nFlushed {
		pendFrom := 0
		if start > nFlushed {
			pendFrom = start - nFlushed
		}
		pendTo := c.pendN
		if end-nFlushed < pendTo {
			pendTo = end - nFlushed
		}
		copyN := pendTo - pendFrom
		if copyN > 0 {
			copy(out[writeIdx:writeIdx+copyN], c.pendBuf[pendFrom:pendTo])
			writeIdx += copyN
		}
	}

	return out[:writeIdx]
}

// ReadAll decodes the entire column into dst as float64, growing dst as
// needed. Returned slice has length == Len().
func (c *CompressedColumn) ReadAll(dst []float64) []float64 {
	n := c.nTotal
	out := dst
	if cap(out) < n {
		out = make([]float64, n)
	} else {
		out = out[:n]
	}
	// Decode blocks directly into a scratch buffer; convert to float64 as we go.
	var scratch [blockSize]int64
	writeIdx := 0
	for bi := range c.blocks {
		nb := c.decodeBlock(bi, scratch[:])
		for j := 0; j < nb; j++ {
			out[writeIdx+j] = float64(scratch[j])
		}
		writeIdx += nb
	}
	// Pending tail.
	for j := 0; j < c.pendN; j++ {
		out[writeIdx+j] = float64(c.pendBuf[j])
	}
	return out
}

// ReadAllInt64 decodes the entire column into dst as int64.
func (c *CompressedColumn) ReadAllInt64(dst []int64) []int64 {
	return c.ReadRange(0, c.nTotal, dst)
}

// CompressedBytes returns an estimate of bytes held by the packed payload
// across all blocks. Excludes Go slice headers and struct fields. Used for
// the compression-ratio sanity check in the plan's verification step.
func (c *CompressedColumn) CompressedBytes() int {
	total := 0
	for i := range c.blocks {
		total += len(c.blocks[i].packed)
	}
	return total
}

// blockCount returns the number of flushed blocks. Useful for tests.
func (c *CompressedColumn) blockCount() int { return len(c.blocks) }
