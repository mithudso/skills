package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestWriteFingerprint_ShapeDigestPresent: writeFingerprint must persist
// the SHA-256 tail digest of ShapeSet into the on-disk JSON so it can be
// read back as shape_digest. Before the fix the digest was computed and
// immediately discarded with `_ = ...`.
func TestWriteFingerprint_ShapeDigestPresent(t *testing.T) {
	dir := t.TempDir()
	entry := fingerprintEntry{
		Provenance: fingerprintProvenance{
			WrittenAt:           time.Now(),
			ParserVersion:       "test",
			PatternCatalogHash:  "abc123",
			InputSHA256:         "deadbeef00000000000000000000000000000000000000000000000000000001",
			MaxConfidenceAtWrite: 0.7,
		},
		ShapeSet: []string{"kind:spike|dir:up|mag:high", "kind:mover|dir:up|mag:mid"},
	}

	if err := writeFingerprint(dir, entry); err != nil {
		t.Fatalf("writeFingerprint: %v", err)
	}

	path := filepath.Join(dir, entry.Provenance.InputSHA256+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading written file: %v", err)
	}

	var got fingerprintEntry
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.ShapeDigest == "" {
		t.Error("shape_digest is empty after round-trip; expected a non-empty hex string")
	}
	// Digest must be stable: calling writeFingerprint again with the same
	// input must produce the same digest.
	entry2 := entry
	entry2.Provenance.InputSHA256 = "deadbeef00000000000000000000000000000000000000000000000000000002"
	if err := writeFingerprint(dir, entry2); err != nil {
		t.Fatalf("writeFingerprint second call: %v", err)
	}
	data2, _ := os.ReadFile(filepath.Join(dir, entry2.Provenance.InputSHA256+".json"))
	var got2 fingerprintEntry
	_ = json.Unmarshal(data2, &got2)
	if got.ShapeDigest != got2.ShapeDigest {
		t.Errorf("digest not stable across identical ShapeSets: %q vs %q", got.ShapeDigest, got2.ShapeDigest)
	}
}

// TestWriteFingerprint_EmptyShapeSet: when ShapeSet is empty the digest
// must remain empty (no error, no spurious digest).
func TestWriteFingerprint_EmptyShapeSet(t *testing.T) {
	dir := t.TempDir()
	entry := fingerprintEntry{
		Provenance: fingerprintProvenance{
			WrittenAt:   time.Now(),
			InputSHA256: "deadbeef00000000000000000000000000000000000000000000000000000003",
		},
		ShapeSet: nil,
	}
	if err := writeFingerprint(dir, entry); err != nil {
		t.Fatalf("writeFingerprint: %v", err)
	}
	data, _ := os.ReadFile(filepath.Join(dir, entry.Provenance.InputSHA256+".json"))
	var got fingerprintEntry
	_ = json.Unmarshal(data, &got)
	if got.ShapeDigest != "" {
		t.Errorf("expected empty shape_digest for empty ShapeSet, got %q", got.ShapeDigest)
	}
}
