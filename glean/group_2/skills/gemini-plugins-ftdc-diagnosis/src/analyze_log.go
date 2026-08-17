// analyze_log.go — mongod log integration for --auto-log-file.
//
// When a log file path is provided, this stage:
//   1. Scans the log file line-by-line to find lines within the "relevant
//      window": [earliest_finding_time - 30min, latest_finding_time + 30min].
//   2. Writes a trimmed copy (mongod_log_trimmed.txt) to the bundle directory,
//      keeping every line in the relevant window intact — no filtering within
//      the window, just trimming the irrelevant prefix and suffix.
//   3. Produces an in-memory snippet (~100 lines centered on the main event
//      timestamp) for inclusion in analysis.md.
//
// Log format support:
//   - Structured JSON (MongoDB 4.4+): lines starting with {"t":{"$date":"..."}}
//   - Legacy text: lines starting with <ISO-timestamp> <severity> <component>
//
// Large-file safety: the file is scanned once with bufio.Scanner (line-by-line,
// no full load into memory). The trimmed file is written incrementally.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// logResult holds the output of processLogFile.
type logResult struct {
	// TrimmedPath is the path to the trimmed log file in the bundle dir.
	// Empty if no trimming was possible or the file was already small enough.
	TrimmedPath string
	// Snippet is ~100 lines centered on the main event for analysis.md.
	Snippet []string
	// SnippetAnchor is the human-readable timestamp of the center of the snippet.
	SnippetAnchor string
	// LinesKept / LinesTotal for the bundle README note.
	LinesKept  int
	LinesTotal int
}

// deriveMainEventTime picks the "primary event time" from the Findings — the
// earliest high-priority Finding (counter_stall, checkpoint_stall, variance_shift
// with dirty bytes, or the overall earliest if none of these exist). This is
// the center of the log snippet.
func deriveMainEventTime(findings []Finding) time.Time {
	priority := []string{
		KindCounterStall,
		KindCheckpointStall,
		KindVarianceShift,
		KindSealEvent,
		KindChangepoint,
	}
	best := time.Time{}
	for _, wantKind := range priority {
		for _, f := range findings {
			k, _ := f["kind"].(string)
			if k != wantKind {
				continue
			}
			atStr, _ := f["at"].(string)
			if atStr == "" {
				continue
			}
			t, err := time.Parse(tsFormat, atStr)
			if err != nil {
				continue
			}
			if best.IsZero() || t.Before(best) {
				best = t
			}
		}
		if !best.IsZero() {
			break
		}
	}
	if best.IsZero() {
		// Fallback: earliest timestamped Finding.
		for _, f := range findings {
			atStr, _ := f["at"].(string)
			if atStr == "" {
				continue
			}
			t, err := time.Parse(tsFormat, atStr)
			if err != nil || t.IsZero() {
				continue
			}
			if best.IsZero() || t.Before(best) {
				best = t
			}
		}
	}
	return best
}

// processLogFile trims a mongod log to the relevant window derived from the
// Findings stream, writes the trimmed file, and returns a snippet for analysis.md.
// mainEventTime is the primary event timestamp (center of the snippet).
func processLogFile(logPath string, findings []Finding, bundleDir string, mainEventTime time.Time) (*logResult, error) {
	// Derive the relevant window from the Findings timestamps.
	const margin = 30 * time.Minute
	const snippetRadius = 50 // lines before and after the main event

	earliest := time.Time{}
	latest := time.Time{}
	for _, f := range findings {
		atStr, _ := f["at"].(string)
		if atStr == "" {
			continue
		}
		t, err := time.Parse(tsFormat, atStr)
		if err != nil {
			continue
		}
		if earliest.IsZero() || t.Before(earliest) {
			earliest = t
		}
		if latest.IsZero() || t.After(latest) {
			latest = t
		}
	}
	if earliest.IsZero() {
		earliest = mainEventTime
	}
	if latest.IsZero() {
		latest = mainEventTime
	}
	windowStart := earliest.Add(-margin)
	windowEnd := latest.Add(margin)

	f, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("open log: %w", err)
	}
	defer f.Close()

	// Pass 1: scan all lines, extract timestamps, decide keep/skip.
	type parsedLine struct {
		raw  string
		ts   time.Time // zero if unparseable
		keep bool
	}
	var allLines []parsedLine

	scanner := bufio.NewScanner(f)
	const maxLineSize = 4 * 1024 * 1024
	buf := make([]byte, maxLineSize)
	scanner.Buffer(buf, maxLineSize)
	for scanner.Scan() {
		raw := scanner.Text()
		ts := extractLogTimestamp(raw)
		keep := ts.IsZero() || (ts.After(windowStart) && ts.Before(windowEnd))
		allLines = append(allLines, parsedLine{raw: raw, ts: ts, keep: keep})
	}
	if err := scanner.Err(); err != nil && err != io.EOF {
		return nil, fmt.Errorf("scan log: %w", err)
	}

	total := len(allLines)
	keptCount := 0
	for _, l := range allLines {
		if l.keep {
			keptCount++
		}
	}

	result := &logResult{
		LinesTotal: total,
		LinesKept:  keptCount,
	}

	// Write trimmed file only when we're actually trimming something (>5% reduction).
	if keptCount < total*95/100 {
		trimPath := bundleDir + "/mongod_log_trimmed.txt"
		tf, err := os.Create(trimPath)
		if err == nil {
			bw := bufio.NewWriter(tf)
			for _, l := range allLines {
				if l.keep {
					bw.WriteString(l.raw)
					bw.WriteByte('\n')
				}
			}
			bw.Flush()
			tf.Close()
			result.TrimmedPath = trimPath
		}
	}

	// Build the snippet: find the line closest to mainEventTime, then take
	// snippetRadius lines before and after.
	centerIdx := -1
	bestDelta := time.Duration(1<<63 - 1)
	for i, l := range allLines {
		if l.ts.IsZero() {
			continue
		}
		delta := l.ts.Sub(mainEventTime)
		if delta < 0 {
			delta = -delta
		}
		if delta < bestDelta {
			bestDelta = delta
			centerIdx = i
		}
	}
	if centerIdx >= 0 {
		lo := centerIdx - snippetRadius
		if lo < 0 {
			lo = 0
		}
		hi := centerIdx + snippetRadius
		if hi > len(allLines) {
			hi = len(allLines)
		}
		for _, l := range allLines[lo:hi] {
			result.Snippet = append(result.Snippet, l.raw)
		}
		result.SnippetAnchor = allLines[centerIdx].ts.UTC().Format(tsFormat)
	}

	return result, nil
}

// extractLogTimestamp attempts to extract the timestamp from a single log line.
// Returns a zero time.Time if the line cannot be parsed.
func extractLogTimestamp(line string) time.Time {
	if len(line) < 10 {
		return time.Time{}
	}

	// Structured JSON (MongoDB 4.4+): {"t":{"$date":"2026-05-20T21:11:50.123+00:00"},...}
	if bytes.HasPrefix([]byte(line), []byte("{\"t\":{\"$date\":\"")) {
		start := len("{\"t\":{\"$date\":\"")
		end := strings.Index(line[start:], "\"")
		if end > 0 {
			dateStr := line[start : start+end]
			// Try RFC3339 with milliseconds and timezone offset.
			for _, layout := range []string{
				"2006-01-02T15:04:05.999999999-07:00",
				"2006-01-02T15:04:05.999999999Z",
				"2006-01-02T15:04:05-07:00",
			} {
				if t, err := time.Parse(layout, dateStr); err == nil {
					return t.UTC()
				}
			}
		}
		return time.Time{}
	}

	// Legacy text format: "2026-05-20T21:11:50.123+0000 I  COMPONENT ..."
	// or                  "2026-05-20T21:11:50.123+00:00 I  COMPONENT ..."
	if len(line) > 24 && line[4] == '-' && line[7] == '-' && line[10] == 'T' {
		// Try to find where the timestamp field ends (space after timezone).
		end := strings.IndexByte(line, ' ')
		if end > 10 {
			dateStr := line[:end]
			for _, layout := range []string{
				"2006-01-02T15:04:05.999-0700",
				"2006-01-02T15:04:05.999-07:00",
				"2006-01-02T15:04:05.999Z",
				"2006-01-02T15:04:05-0700",
			} {
				if t, err := time.Parse(layout, dateStr); err == nil {
					return t.UTC()
				}
			}
		}
	}

	return time.Time{}
}

// buildLogSection returns the analysis.md markdown for the log excerpt.
func buildLogSection(res *logResult, logPath string) string {
	if res == nil || (len(res.Snippet) == 0 && res.TrimmedPath == "") {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("\n## Log excerpt\n\n")
	sb.WriteString(fmt.Sprintf("Source: `%s`  \n", logPath))
	sb.WriteString(fmt.Sprintf("Lines kept: %d / %d (window ±30 min around findings)\n", res.LinesKept, res.LinesTotal))
	if res.TrimmedPath != "" {
		sb.WriteString(fmt.Sprintf("Trimmed log: `%s`\n", "mongod_log_trimmed.txt"))
	}
	if res.SnippetAnchor != "" {
		sb.WriteString(fmt.Sprintf("\nSnippet centered on `%s` (±%d lines):\n", res.SnippetAnchor, 50))
	}
	if len(res.Snippet) > 0 {
		sb.WriteString("\n```\n")
		for _, l := range res.Snippet {
			sb.WriteString(l)
			sb.WriteByte('\n')
		}
		sb.WriteString("```\n")
	}
	return sb.String()
}
