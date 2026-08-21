package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestHostTokenDeterministic verifies that the same hostname+salt produces
// the same token, but different salts produce different tokens.
func TestHostTokenDeterministic(t *testing.T) {
	hostname := "mongo-prod-01.example.com"
	salt1 := "abcdef0123456789abcdef0123456789"
	salt2 := "ffffffffffffffffffffffffffffffff"

	token1a := hostToken(hostname, salt1)
	token1b := hostToken(hostname, salt1)
	token2 := hostToken(hostname, salt2)

	if token1a != token1b {
		t.Errorf("determinism failure: same hostname+salt produced different tokens")
	}
	if token1a == token2 {
		t.Errorf("salt isolation failure: different salts produced same token")
	}
	// host_token format should be hex with no whitespace.
	if len(token1a) != 32 { // 16 bytes hex
		t.Errorf("token length wrong: got %d, want 32", len(token1a))
	}
}

// TestBaselineRowExpiryViaLoad exercises the ACTUAL loadBaseline TTL logic
// (not just the time math in isolation). Writes a baseline.jsonl with one
// fresh row and one expired row, calls loadBaseline, and asserts the expired
// row is dropped.
func TestBaselineRowExpiryViaLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.jsonl")

	now := time.Now().UTC()
	oldTS := now.Add(-(baselineTTLDays + 1) * 24 * time.Hour).Format(time.RFC3339Nano)
	freshTS := now.Add(-1 * 24 * time.Hour).Format(time.RFC3339Nano)

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	enc := json.NewEncoder(f)
	// Salt line first (required for the loader not to regenerate).
	if err := enc.Encode(baselineSaltLine{Kind: "salt", Salt: "0123456789abcdef0123456789abcdef"}); err != nil {
		t.Fatal(err)
	}
	if err := enc.Encode(baselineRow{
		Kind: "baseline", HostToken: "h1", Metric: "fresh.metric",
		EWMAMean: 100, N: 5, LastUpdateTS: freshTS, ParserVersion: ParserVersion,
	}); err != nil {
		t.Fatal(err)
	}
	if err := enc.Encode(baselineRow{
		Kind: "baseline", HostToken: "h1", Metric: "expired.metric",
		EWMAMean: 200, N: 5, LastUpdateTS: oldTS, ParserVersion: ParserVersion,
	}); err != nil {
		t.Fatal(err)
	}
	f.Close()

	_, baseRows, _, err := loadBaseline(path)
	if err != nil {
		t.Fatalf("loadBaseline error: %v", err)
	}
	if _, ok := baseRows["h1|fresh.metric"]; !ok {
		t.Errorf("expected fresh row to be loaded; baseRows=%v", baseRows)
	}
	if _, ok := baseRows["h1|expired.metric"]; ok {
		t.Errorf("expected expired row to be dropped on load; baseRows=%v", baseRows)
	}
}

// TestBaselineWriteReadRoundTrip exercises the atomic-rename write path and
// verifies the round-trip preserves field values.
func TestBaselineWriteReadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.jsonl")

	salt := "0123456789abcdef0123456789abcdef"
	rows := map[string]baselineRow{
		"h1|m1": {
			Kind: "baseline", HostToken: "h1", Metric: "m1",
			EWMAMean: 3.14, EWMAVar: 0.5, N: 10,
			LastUpdateTS:  time.Now().UTC().Format(time.RFC3339Nano),
			ParserVersion: ParserVersion,
		},
	}
	if err := writeBaseline(path, salt, rows, nil); err != nil {
		t.Fatalf("writeBaseline error: %v", err)
	}
	gotSalt, gotBase, _, err := loadBaseline(path)
	if err != nil {
		t.Fatalf("loadBaseline error: %v", err)
	}
	if gotSalt != salt {
		t.Errorf("salt mismatch: got %q want %q", gotSalt, salt)
	}
	r, ok := gotBase["h1|m1"]
	if !ok {
		t.Fatalf("row not loaded: %v", gotBase)
	}
	if r.EWMAMean != 3.14 || r.EWMAVar != 0.5 || r.N != 10 {
		t.Errorf("row fields mismatch: %+v", r)
	}
}
