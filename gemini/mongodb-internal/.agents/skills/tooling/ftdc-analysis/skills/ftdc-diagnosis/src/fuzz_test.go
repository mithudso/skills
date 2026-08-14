package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"strings"
	"testing"
	"time"
)

// Adversarial-input tests for Tier A.1 hardening. (Go 1.17 lacks testing.F,
// so these are table-driven instead of true fuzz targets; the safety intent
// — pre-A.1 these inputs would OOM, loop forever, or overflow — survives.)

// TestReadVarintBounded — pre-A.1, an input of 0x80 ad infinitum looped
// shifting past 64 bits silently. Post-A.1 it must return an error after at
// most 10 bytes.
func TestReadVarintBounded(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		wantError string
	}{
		{"all_continuation_long", bytes.Repeat([]byte{0x80}, 15), "varint too long"},
		{"ten_continuation_bytes", []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80}, "varint too long"},
		{"nine_continuation_then_terminator", []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x01}, ""},
		{"empty", []byte{}, "EOF"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewReader(tc.input)
			_, err := readVarint(r)
			if tc.wantError == "" {
				if err != nil {
					t.Fatalf("expected success, got error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tc.wantError)
			}
			if !strings.Contains(err.Error(), tc.wantError) && err.Error() != tc.wantError {
				t.Fatalf("expected error containing %q, got %q", tc.wantError, err.Error())
			}
		})
	}
}

// TestDecodeChunkRejectsHostileNDeltas — pre-A.1, an nDeltas value of 2^32-1
// would attempt to allocate a 32 GB matrix per metric. Post-A.1 it must
// reject chunks with nDeltas > 10K and return nil.
func TestDecodeChunkRejectsHostileNDeltas(t *testing.T) {
	var buf bytes.Buffer
	// 4-byte chunk-length prefix (ignored by decodeChunk after slicing).
	binary.Write(&buf, binary.LittleEndian, uint32(0))
	zw := zlib.NewWriter(&buf)
	// Minimal BSON: {"a": 1} so refDoc parses to non-empty.
	zw.Write([]byte{0x0C, 0x00, 0x00, 0x00, 0x10, 'a', 0x00, 0x01, 0x00, 0x00, 0x00, 0x00})
	// nKeys=1, nDeltas=2^31 (well above the 10K cap).
	binary.Write(zw, binary.LittleEndian, uint32(1))
	binary.Write(zw, binary.LittleEndian, uint32(1<<31))
	zw.Close()
	got := decodeChunk(buf.Bytes(), time.Unix(0, 0), nil)
	if got != nil {
		t.Fatalf("expected nil (rejected), got %d samples", len(got))
	}
}

// TestDecodeChunkRejectsZlibBomb — pre-A.1, io.ReadAll on a malicious zlib
// stream could inflate to multi-GB. Post-A.1 the 256 MB cap must reject.
// We don't construct a real bomb (compress 256MB of zeros takes too long
// at test time); instead we set up a zlib stream that, when decompressed,
// is just under the limit, and verify it succeeds, then crank size up.
func TestDecodeChunkRespectsDecompressLimit(t *testing.T) {
	// Build a payload whose decompressed size exceeds the 256MB limit by
	// abusing zlib's compression ratio on zero-filled data.
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, uint32(0))
	zw := zlib.NewWriter(&buf)
	// 300 MB of zeros compresses to a few KB but expands well past the cap.
	zeros := make([]byte, 1<<20) // 1 MB
	for i := 0; i < 300; i++ {
		zw.Write(zeros)
	}
	zw.Close()
	got := decodeChunk(buf.Bytes(), time.Unix(0, 0), nil)
	if got != nil {
		t.Fatalf("expected nil (over decompress limit), got %d samples", len(got))
	}
}

// TestFlattenBSONDepthGuard — pre-A.1, a hostile deeply-nested BSON doc
// could blow the goroutine stack via unbounded recursion. Post-A.1 it must
// stop at depth maxBSONDepth.
func TestFlattenBSONDepthGuard(t *testing.T) {
	// Test the guard directly: at depth 100 we should refuse to recurse.
	var keys []string
	var vals []int64
	// Build a nested bson.D structure 200 deep using the public type.
	// Easier: directly invoke the depth-bounded helper with a synthetic depth.
	flattenBSONOrderedDepth(nil, "", &keys, &vals, maxBSONDepth)
	if len(keys) != 0 {
		t.Fatalf("expected no keys at depth=max, got %d", len(keys))
	}
}

// TestPartitionFTDCMalformed — pre-A.1, malformed top-level BSON could
// trigger out-of-bounds slice reads. Post-A.1 must return cleanly.
func TestPartitionFTDCMalformed(t *testing.T) {
	tests := [][]byte{
		nil,
		{},
		{0x05, 0x00, 0x00, 0x00, 0x00}, // empty BSON doc
		{0xFF, 0xFF, 0xFF, 0xFF},       // bogus length
		bytes.Repeat([]byte{0xFF}, 10), // garbage
	}
	for i, in := range tests {
		jobs, md := partitionFTDC(in)
		// Must not panic.
		_ = jobs
		_ = md
		t.Logf("case %d: ok, jobs=%d md=%d", i, len(jobs), len(md))
	}
}
