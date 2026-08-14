package main

import (
	"strings"
	"testing"
)

// TestNormalizeFindings filters run_metadata, run_health, baseline_diagnostic,
// and capture_quality_score from NDJSON input.
func TestNormalizeFindings(t *testing.T) {
	input := `{"kind":"run_metadata","host":"h1","window_start":"2026-05-21T00:00:00Z"}
{"kind":"mover","metric":"test.metric","score":1.5}
{"kind":"run_health","host":"h1","findings_count":3}
{"kind":"baseline_diagnostic","host_token":"xyz","metric":"test.metric"}
{"kind":"mover","metric":"test.metric2","score":2.0}`

	got := normalizeFindings(input)

	// Should only have 2 mover lines
	if len(got) != 2 {
		t.Errorf("normalizeFindings: expected 2 lines, got %d", len(got))
	}

	// Verify no filtered kinds remain
	for _, line := range got {
		if strings.Contains(line, `"kind":"run_metadata"`) {
			t.Errorf("run_metadata not filtered")
		}
		if strings.Contains(line, `"kind":"run_health"`) {
			t.Errorf("run_health not filtered")
		}
		if strings.Contains(line, `"kind":"baseline_diagnostic"`) {
			t.Errorf("baseline_diagnostic not filtered")
		}
		if strings.Contains(line, `"kind":"capture_quality_score"`) {
			t.Errorf("capture_quality_score not filtered")
		}
	}
}

// TestNormalizeFindingsSorted verifies output is sorted.
func TestNormalizeFindingsSorted(t *testing.T) {
	input := `{"metric":"zebra"}
{"metric":"apple"}
{"metric":"mango"}`

	got := normalizeFindings(input)

	// Verify sorted
	for i := 1; i < len(got); i++ {
		if got[i] < got[i-1] {
			t.Errorf("output not sorted: %q >= %q", got[i-1], got[i])
		}
	}
}

// TestDiffLinesEmpty checks that identical inputs produce no diff.
func TestDiffLinesEmpty(t *testing.T) {
	want := []string{"a", "b", "c"}
	got := []string{"a", "b", "c"}

	miss, extra := diffLines(want, got)
	if len(miss) != 0 || len(extra) != 0 {
		t.Errorf("identical inputs: expected no diff, got miss=%d extra=%d", len(miss), len(extra))
	}
}

// TestDiffLinesAddition checks extra lines are detected.
func TestDiffLinesAddition(t *testing.T) {
	want := []string{"a", "b"}
	got := []string{"a", "b", "c"}

	miss, extra := diffLines(want, got)
	if len(miss) != 0 {
		t.Errorf("expected 0 missing, got %d", len(miss))
	}
	if len(extra) != 1 {
		t.Errorf("expected 1 extra, got %d", len(extra))
	}
	if len(extra) > 0 && extra[0] != "c" {
		t.Errorf("expected extra[0]='c', got %q", extra[0])
	}
}

// TestDiffLinesMissing checks missing lines are detected.
func TestDiffLinesMissing(t *testing.T) {
	want := []string{"a", "b", "c"}
	got := []string{"a", "c"}

	miss, extra := diffLines(want, got)
	if len(miss) != 1 {
		t.Errorf("expected 1 missing, got %d", len(miss))
	}
	if len(extra) != 0 {
		t.Errorf("expected 0 extra, got %d", len(extra))
	}
	if len(miss) > 0 && miss[0] != "b" {
		t.Errorf("expected miss[0]='b', got %q", miss[0])
	}
}

// TestDiffLinesMixed checks both missing and extra.
func TestDiffLinesMixed(t *testing.T) {
	want := []string{"a", "b", "c", "d"}
	got := []string{"a", "c", "d", "e"}

	miss, extra := diffLines(want, got)
	if len(miss) != 1 || miss[0] != "b" {
		t.Errorf("expected miss=[b], got %v", miss)
	}
	if len(extra) != 1 || extra[0] != "e" {
		t.Errorf("expected extra=[e], got %v", extra)
	}
}

// TestTrimForReport truncates long lines.
func TestTrimForReport(t *testing.T) {
	short := "short finding line"
	long := strings.Repeat("x", 300)

	if trimForReport(short) != short {
		t.Errorf("short line should not be trimmed")
	}

	trimmed := trimForReport(long)
	if len(trimmed) >= len(long) {
		t.Errorf("long line should be trimmed")
	}
	if !strings.HasSuffix(trimmed, "...") {
		t.Errorf("trimmed output should end with '...'")
	}
}
