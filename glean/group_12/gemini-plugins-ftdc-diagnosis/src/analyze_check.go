// analyze_check.go - --auto-check test harness.
//
// Iterates testdata/* fixture subdirectories, runs --auto on each, normalizes
// the output (strips wall-clock timestamps that vary per run), and diffs
// against expected_findings.ndjson if present. Prints PASS/FAIL/MISSING/EXTRA
// per fixture and exits non-zero on any failure. Fixture-addition workflow
// is documented in SKILL.md.
//
// Implementation is intentionally minimal in v1: the diff is line-set
// equality (after normalization) and report formatting is plain text.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// autoCheckRun is invoked from main() when --auto-check DIR is set. It walks
// DIR, runs --auto on each subdirectory that contains FTDC files, compares
// Findings against the per-fixture golden file, and prints a summary.
func autoCheckRun(w io.Writer, dir string, jobs int) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("--auto-check: cannot read %s: %v", dir, err)
	}
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	var (
		passed int
		failed int
	)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		fixDir := filepath.Join(dir, e.Name())
		goldenPath := filepath.Join(fixDir, "expected_findings.ndjson")
		if _, err := os.Stat(goldenPath); err != nil {
			fmt.Fprintf(bw, "[SKIP] %s — no expected_findings.ndjson\n", e.Name())
			continue
		}
		// Run --auto on the fixture and capture Findings.
		var got bytes.Buffer
		var samples []Sample
		metadata := processFTDC(fixDir, nil, jobs, func(batch []Sample) {
			samples = append(samples, batch...)
		})
		cfg := defaultAutoConfig()
		if err := runAuto(&got, fixDir, samples, metadata, cfg); err != nil {
			fmt.Fprintf(bw, "[FAIL] %s — runAuto error: %v\n", e.Name(), err)
			failed++
			continue
		}
		gotLines := normalizeFindings(got.String())
		expBytes, err := os.ReadFile(goldenPath)
		if err != nil {
			fmt.Fprintf(bw, "[FAIL] %s — read golden: %v\n", e.Name(), err)
			failed++
			continue
		}
		expLines := normalizeFindings(string(expBytes))
		miss, extra := diffLines(expLines, gotLines)
		if len(miss) == 0 && len(extra) == 0 {
			fmt.Fprintf(bw, "[PASS] %s — %d findings match\n", e.Name(), len(gotLines))
			passed++
		} else {
			fmt.Fprintf(bw, "[FAIL] %s — missing=%d, extra=%d\n", e.Name(), len(miss), len(extra))
			// Cap report lines so a totally broken run doesn't explode
			// the output. The full diff is still observable via direct
			// file comparison; this just bounds the harness output.
			const maxReport = 100
			missCap := len(miss)
			if missCap > maxReport {
				missCap = maxReport
			}
			for _, m := range miss[:missCap] {
				fmt.Fprintf(bw, "   MISSING: %s\n", trimForReport(m))
			}
			if len(miss) > maxReport {
				fmt.Fprintf(bw, "   ... %d more missing\n", len(miss)-maxReport)
			}
			extraCap := len(extra)
			if extraCap > maxReport {
				extraCap = maxReport
			}
			for _, x := range extra[:extraCap] {
				fmt.Fprintf(bw, "   EXTRA:   %s\n", trimForReport(x))
			}
			if len(extra) > maxReport {
				fmt.Fprintf(bw, "   ... %d more extra\n", len(extra)-maxReport)
			}
			failed++
		}
	}
	fmt.Fprintf(bw, "\n%d passed, %d failed\n", passed, failed)
	if failed > 0 {
		return fmt.Errorf("--auto-check: %d fixture(s) failed", failed)
	}
	return nil
}

// normalizeFindings parses NDJSON lines, drops fields that vary across runs
// (wall-clock timestamps inside Finding ids and run_metadata's
// window_start/window_end), and returns a sorted slice for diff. We do this
// purely textually: split on newlines, drop empties, strip well-known
// volatile fields by regex-like substring replacement. A future iteration
// can move to a JSON-aware normalizer for stricter comparison.
func normalizeFindings(s string) []string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		// Drop run_metadata entirely (window bounds, sample counts vary).
		if strings.Contains(l, `"kind":"run_metadata"`) {
			continue
		}
		// Drop run_health (counts depend on what fired).
		if strings.Contains(l, `"kind":"run_health"`) {
			continue
		}
		// Drop baseline_diagnostic (Tier D) — host_token + timestamps
		// + row counts vary per run, breaking golden-file diffs even
		// when nothing meaningful changed.
		if strings.Contains(l, `"kind":"baseline_diagnostic"`) {
			continue
		}
		// Drop capture_quality_score (Tier G.3) — composite score
		// depends on run-variable Finding counts; not stable across
		// detection-algorithm tweaks.
		if strings.Contains(l, `"kind":"capture_quality_score"`) {
			continue
		}
		// P2.1 / P4.4 / P2.4 / P4.1 / P4.3: stateful Findings whose
		// emission depends on the local filesystem (~/claude-dump/...).
		// These vary between back-to-back runs as files are written,
		// breaking golden-file equality even when the parser logic
		// hasn't changed. Auto-check is logic equivalence, not state
		// equivalence — strip these kinds the same way run_metadata
		// is stripped.
		if strings.Contains(l, `"kind":"recap"`) ||
			strings.Contains(l, `"kind":"similar_prior_capture"`) ||
			strings.Contains(l, `"kind":"stored_inference_review"`) ||
			strings.Contains(l, `"kind":"followups_backlog_snapshot"`) ||
			strings.Contains(l, `"kind":"converge_plan"`) ||
			strings.Contains(l, `"kind":"convergence_failed"`) ||
			strings.Contains(l, `"kind":"sweep_discriminator"`) ||
			strings.Contains(l, `"kind":"sweep_label_audit"`) ||
			strings.Contains(l, `"kind":"sweep_inconclusive"`) {
			continue
		}
		out = append(out, l)
	}
	sort.Strings(out)
	return out
}

// diffLines returns (missing, extra) — sorted-set diff between expected and
// got. Both inputs must be pre-sorted.
func diffLines(want, got []string) (missing, extra []string) {
	i, j := 0, 0
	for i < len(want) && j < len(got) {
		switch {
		case want[i] < got[j]:
			missing = append(missing, want[i])
			i++
		case want[i] > got[j]:
			extra = append(extra, got[j])
			j++
		default:
			i++
			j++
		}
	}
	missing = append(missing, want[i:]...)
	extra = append(extra, got[j:]...)
	return
}

// trimForReport shortens a Finding line for display in the failure report.
func trimForReport(s string) string {
	if len(s) > 240 {
		return s[:237] + "..."
	}
	return s
}
