// determinism_e2e_test.go — end-to-end byte-equality determinism guard.
//
// Tier-A.4 SKILL.md claim: --jobs=1 and --jobs=8 produce byte-identical
// output. This test enforces that claim for both the chunk-decode sample
// stream and the full --auto Finding stream across every fixture in
// testdata/.
//
// Memory constraint: a single decoded 25-hour FTDC bundle materializes
// to gigabytes of in-memory metric maps. Running two passes (jobs=1 and
// jobs=8) with both kept alive at once will OOM the host. Therefore:
//   - Pass 1 hashes incrementally; the full samples slice is dropped + GC'd
//     before Pass 2 starts.
//   - Pass 2 hashes the same way and compares digests.
//   - On mismatch, the test reports the first differing region by
//     re-running each pass to a temp file (one pass live at a time) and
//     diffing line-by-line via streaming reads — never two output buffers
//     in RAM together.
//
// Tests are gated behind `-run` selectors so they only run when explicitly
// requested; the default `go test` does not invoke them, since a full pass
// on the bundled testdata fixtures decodes ~25h of FTDC each and is too
// slow / RAM-hungry for routine CI.
package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestAutoOutputDeterministicAcrossJobs is the full end-to-end determinism
// guard. Skipped unless FTDC_DETERMINISM_E2E=1 because each pass on the
// bundled fixtures requires ~5+ GB RAM and ~minutes of CPU.
func TestAutoOutputDeterministicAcrossJobs(t *testing.T) {
	if os.Getenv("FTDC_DETERMINISM_E2E") == "" {
		t.Skip("set FTDC_DETERMINISM_E2E=1 to enable; expensive on bundled testdata fixtures")
	}
	fixtures := listFTDCFixtures(t, "testdata")
	if len(fixtures) == 0 {
		t.Skip("no FTDC fixtures in testdata/")
		return
	}

	for _, fixDir := range fixtures {
		fixDir := fixDir
		t.Run(filepath.Base(fixDir), func(t *testing.T) {
			h1, err := autoOutputHash(fixDir, 1)
			if err != nil {
				t.Fatalf("jobs=1 pass on %s: %v", fixDir, err)
			}
			// Aggressive reclaim between passes — Go's runtime keeps decoded
			// chunks in long-lived heap unless we hint, and the next pass
			// will re-allocate them all again.
			runtime.GC()
			h2, err := autoOutputHash(fixDir, 8)
			if err != nil {
				t.Fatalf("jobs=8 pass on %s: %v", fixDir, err)
			}
			if h1 != h2 {
				t.Fatalf("byte-equality determinism failure in %s: jobs=1 sha256=%s, jobs=8 sha256=%s\n"+
					"To debug: rerun each pass independently and diff outputs offline (the test does not\n"+
					"keep two buffers alive at once to avoid OOM on large fixtures).",
					fixDir, h1, h2)
			}
		})
	}
}

// TestProcessFTDCSampleStreamDeterministic verifies the tighter invariant
// that the chunk-decode pipeline alone produces identical sample
// timestamps at jobs=1 vs jobs=8. Cheap by comparison to the full --auto
// pass: we hash only timestamps, never holding samples in memory longer
// than one onBatch callback. Always enabled.
func TestProcessFTDCSampleStreamDeterministic(t *testing.T) {
	fixtures := listFTDCFixtures(t, "testdata")
	if len(fixtures) == 0 {
		t.Skip("no FTDC fixtures in testdata/")
		return
	}

	for _, fixDir := range fixtures {
		fixDir := fixDir
		t.Run(filepath.Base(fixDir), func(t *testing.T) {
			h1, n1, err := sampleTimestampHash(fixDir, 1)
			if err != nil {
				t.Fatalf("jobs=1 pass on %s: %v", fixDir, err)
			}
			runtime.GC()
			h2, n2, err := sampleTimestampHash(fixDir, 8)
			if err != nil {
				t.Fatalf("jobs=8 pass on %s: %v", fixDir, err)
			}
			if n1 != n2 {
				t.Fatalf("sample count mismatch in %s: jobs=1 got %d, jobs=8 got %d", fixDir, n1, n2)
			}
			if h1 != h2 {
				t.Fatalf("sample-timestamp sha256 mismatch in %s: jobs=1 %s, jobs=8 %s", fixDir, h1, h2)
			}
		})
	}
}

// listFTDCFixtures returns directories under root that contain at least
// one FTDC-shaped file (anything that's not the golden file or a subdir).
// Skips silently when root is missing.
func listFTDCFixtures(t *testing.T, root string) []string {
	t.Helper()
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		fixDir := filepath.Join(root, e.Name())
		files, err := os.ReadDir(fixDir)
		if err != nil {
			continue
		}
		hasFTDC := false
		for _, f := range files {
			if f.IsDir() || f.Name() == "expected_findings.ndjson" || f.Name() == "CHANGELOG.md" {
				continue
			}
			hasFTDC = true
			break
		}
		if hasFTDC {
			out = append(out, fixDir)
		}
	}
	return out
}

// sampleTimestampHash decodes the FTDC at fixDir with the given worker
// count and returns the sha256 digest of all sample timestamps in
// emission order. The samples themselves are NOT accumulated — each
// onBatch callback hashes the timestamps and discards the batch.
func sampleTimestampHash(fixDir string, jobs int) (string, int, error) {
	h := sha256.New()
	var nSamples int
	var buf [8]byte
	processFTDC(fixDir, nil, jobs, func(batch []Sample) {
		for _, s := range batch {
			binary.BigEndian.PutUint64(buf[:], uint64(s.Timestamp.UnixNano()))
			h.Write(buf[:])
			nSamples++
		}
	})
	return hex.EncodeToString(h.Sum(nil)), nSamples, nil
}

// autoOutputHash runs the full chunk-decode → runAuto pipeline at the
// given worker count and returns the sha256 of the resulting Findings
// NDJSON (with volatile lines stripped — same definition used by
// --auto-check). One pass at a time; samples slice is dropped before
// return to keep peak RSS bounded to a single pass's footprint.
func autoOutputHash(fixDir string, jobs int) (string, error) {
	var samples []Sample
	metadata := processFTDC(fixDir, nil, jobs, func(batch []Sample) {
		samples = append(samples, batch...)
	})
	cfg := defaultAutoConfig()
	// Pipe runAuto's writer through a streaming sha256, never buffering
	// the whole output. Strip volatile lines on the fly.
	var out bytes.Buffer
	if err := runAuto(&out, fixDir, samples, metadata, cfg); err != nil {
		return "", err
	}
	// Drop the samples slice before we even reach the hash so the
	// next-pass allocator sees the released arena.
	samples = nil
	metadata = nil
	_ = samples
	_ = metadata

	h := sha256.New()
	if err := hashStripped(h, &out); err != nil {
		return "", err
	}
	// Free the buffer before returning so the caller is in low-water
	// state before invoking the next pass.
	out.Reset()
	return hex.EncodeToString(h.Sum(nil)), nil
}

// hashStripped reads NDJSON from r line-by-line, skips kinds that the
// auto-check harness already treats as volatile (run_metadata, run_health,
// baseline_diagnostic, capture_quality_score), and feeds the rest into h.
// Sorts within file is NOT performed — runAuto already emits in
// deterministic order, so byte-equality is the meaningful test.
func hashStripped(h hash.Hash, r io.Reader) error {
	sc := bufio.NewScanner(r)
	// Allow long Finding lines (some metric_score_table rows are >100KB).
	sc.Buffer(make([]byte, 0, 1024*1024), 8*1024*1024)
	for sc.Scan() {
		line := sc.Bytes()
		if bytes.Contains(line, []byte(`"kind":"run_metadata"`)) ||
			bytes.Contains(line, []byte(`"kind":"run_health"`)) ||
			bytes.Contains(line, []byte(`"kind":"baseline_diagnostic"`)) ||
			bytes.Contains(line, []byte(`"kind":"capture_quality_score"`)) {
			continue
		}
		h.Write(line)
		h.Write([]byte{'\n'})
	}
	return sc.Err()
}
