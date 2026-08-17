// input_fingerprint.go — provenance helpers for run_metadata.
//
// The bias-control invariant ("show the choice") requires every algorithmic
// threshold to be echoed in run_metadata. The same principle extends to
// the inputs: a reviewer must be able to confirm "same bytes, same invocation"
// when comparing two runs. P0.3 adds:
//   - input_sha256: content hash over the FTDC input(s).
//   - input_files: per-file hashes when the input is a directory, so a
//     reviewer can tell which file changed without re-hashing.
//   - cmd_line: argv with path-bearing positions reduced to basenames so
//     the same invocation from different cwd doesn't look like a different
//     run.
//
// Hashes are streaming SHA256; memory footprint is O(1) per file. No
// dependencies beyond the stdlib.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// fileSHA256 streams a file through SHA256 and returns hex(digest). Files
// larger than 100MB are streamed in 1MB chunks so memory stays bounded.
func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	buf := make([]byte, 1<<20)
	if _, err := io.CopyBuffer(h, f, buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// InputFingerprint captures everything emitRunMetadata needs to populate
// the input_sha256 / input_files fields. Caller computes once (typically
// from main() before runAuto is invoked) and threads it through AutoConfig.
type InputFingerprint struct {
	// TopSHA256 is the canonical content digest:
	//   - single file: sha256 of the file's content
	//   - directory:   sha256 of "<basename>:<file_sha256>\n" lines sorted
	//                  by basename, joined.
	// Same input bytes => same TopSHA256 regardless of cwd or file ordering
	// on disk.
	TopSHA256 string
	// PerFileSHA256 is populated only in directory mode; the keys are
	// file basenames (no leading paths). Hidden files and non-FTDC files
	// are excluded using the same rules processFTDC applies, so the
	// fingerprint reflects exactly what the parser will decode.
	PerFileSHA256 map[string]string
	// IsDirectory disambiguates between single-file and directory mode in
	// downstream consumers; some narrators care.
	IsDirectory bool
}

// computeInputFingerprint walks the given path and produces a streaming
// SHA256 fingerprint. Returns nil + err on filesystem errors so the caller
// can decide whether to abort or proceed without provenance.
func computeInputFingerprint(path string) (*InputFingerprint, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !st.IsDir() {
		h, err := fileSHA256(path)
		if err != nil {
			return nil, err
		}
		return &InputFingerprint{TopSHA256: h, IsDirectory: false}, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	perFile := make(map[string]string, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		// Skip the same set the parser skips: hidden files and the
		// per-fixture golden file (only present in testdata layouts).
		if strings.HasPrefix(name, ".") || name == "expected_findings.ndjson" {
			continue
		}
		h, err := fileSHA256(filepath.Join(path, name))
		if err != nil {
			return nil, err
		}
		perFile[name] = h
	}

	// Canonical aggregate: sort by basename, join "name:hash\n" lines.
	names := make([]string, 0, len(perFile))
	for n := range perFile {
		names = append(names, n)
	}
	sort.Strings(names)
	agg := sha256.New()
	for _, n := range names {
		agg.Write([]byte(n))
		agg.Write([]byte{':'})
		agg.Write([]byte(perFile[n]))
		agg.Write([]byte{'\n'})
	}
	return &InputFingerprint{
		TopSHA256:     hex.EncodeToString(agg.Sum(nil)),
		PerFileSHA256: perFile,
		IsDirectory:   true,
	}, nil
}

// normalizedCmdLine returns os.Args with path-bearing positions reduced
// to basenames. The first element (argv[0]) is reduced to its basename
// regardless of position; subsequent elements are only basenamed when
// they appear to be filesystem paths (contain a path separator and exist
// or look like one). Flag values starting with "-" are passed through
// unchanged.
//
// Why: a reviewer running ftdc_parser /home/user/cap1 and another running
// ftdc_parser /tmp/cap1 against the same bytes should see the same
// cmd_line in run_metadata. The InputSHA256 already confirms content
// equivalence; this confirms invocation equivalence.
func normalizedCmdLine(argv []string) []string {
	out := make([]string, len(argv))
	for i, a := range argv {
		if i == 0 {
			out[i] = filepath.Base(a)
			continue
		}
		if strings.HasPrefix(a, "-") {
			// flag; leave as-is. We don't try to parse flag values that
			// might be paths (e.g. --auto-baseline-file PATH) because
			// the parser's own flag-value handling normalizes those
			// elsewhere and over-aggressive basenaming risks losing
			// disambiguation (different paths with the same basename
			// are real).
			out[i] = a
			continue
		}
		if strings.ContainsAny(a, "/\\") {
			out[i] = filepath.Base(a)
			continue
		}
		out[i] = a
	}
	return out
}
