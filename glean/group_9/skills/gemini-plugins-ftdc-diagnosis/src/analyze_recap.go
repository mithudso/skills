// analyze_recap.go — P2.1 --auto-recap cross-session continuity.
//
// Re-opening the same FTDC capture days/weeks later, the solo user wants
// to know "what did I find last time?" Without persistence, the parser
// re-derives everything from scratch and the user re-reads the whole
// stream to spot what was important.
//
// Per-capture-digest recap files live at
//   ~/claude-dump/ftdc-recaps/<input_sha256>.ndjson
// Each line is one compact summary written by a past --auto run. On a
// fresh run, the recap stage reads the file (if present), validates
// each prior entry against the stored-inference self-doubt rules, and
// emits one of:
//
//   kind:"recap"                  — prior entries deemed still-valid
//   kind:"stored_inference_review" — entries quarantined for staleness,
//                                    version skew, or contradiction
//
// At the end of the run, the stage appends the current run's summary
// to the file (capped depth — old summaries past N entries roll off).
//
// Solo + FTDC-only: no external corpus. The file is a private rolling
// log per capture, keyed entirely by the input bytes' SHA256.
package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// recapStoreDir returns the directory where per-capture recap files live.
// $HOME/claude-dump/ftdc-recaps when $HOME is set; falls back to /tmp.
func recapStoreDir() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		home = "/tmp"
	}
	return filepath.Join(home, "claude-dump", "ftdc-recaps")
}

// recapMaxEntriesPerFile caps the per-digest recap file at this many
// lines. Older lines roll off (kept in a `.archive` sibling for manual
// inspection if the user wants longer history).
const recapMaxEntriesPerFile = 20

// recapStaleAge: recap entries older than this are quarantined unless
// they came from an exact-match parser_version + pattern_catalog_hash.
// 90 days mirrors the plan's stored-inference self-doubt window.
const recapStaleAge = 90 * 24 * time.Hour

// recapEntry is one line in the recap file. Mirrors the stored-inference
// self-doubt provenance contract.
type recapEntry struct {
	Provenance        recapProvenance     `json:"provenance"`
	TopPatterns       []recapPatternEntry `json:"top_patterns"`
	CaptureQualityScore *int              `json:"capture_quality_score,omitempty"`
	TopMovers         []string            `json:"top_movers,omitempty"` // metric paths only
	SignificantCount  int                 `json:"significant_count"`    // BH q<0.05
	WindowFrom        string              `json:"window_from,omitempty"`
	WindowTo          string              `json:"window_to,omitempty"`
}

type recapProvenance struct {
	WrittenAt           time.Time `json:"written_at"`
	ParserVersion       string    `json:"parser_version"`
	PatternCatalogHash  string    `json:"pattern_catalog_hash"`
	InputSHA256         string    `json:"input_sha256"`
	ConfidenceRankAtWrite float64 `json:"confidence_rank_at_write"`
}

type recapPatternEntry struct {
	PatternName     string  `json:"pattern_name"`
	ConfidenceRank  float64 `json:"confidence_rank"`
	ParentFindingID string  `json:"parent_finding_id"`
}

// emitRecap reads the prior recap file (keyed by current input_sha256)
// and emits a `recap` Finding summarizing still-valid entries; emits
// `stored_inference_review` for entries quarantined by self-doubt. At
// the end, appends the current run's summary to the file.
//
// File-system errors are non-fatal — recap is advisory only.
func emitRecap(c *findingCollector, host string) (emitted int, skipped int) {
	// Pull input_sha256 + parser_version from run_metadata.
	var inputSHA256, parserVer string
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		if k != "run_metadata" {
			continue
		}
		if s, ok := f["input_sha256"].(string); ok {
			inputSHA256 = s
		}
		if lv, ok := f["library_versions"].(map[string]interface{}); ok {
			if pv, ok := lv["parser_version"].(string); ok {
				parserVer = pv
			}
		}
		break
	}
	if inputSHA256 == "" {
		// Without a capture digest there's no recap key; bail silently.
		return 0, 1
	}

	catalogHash := patternCatalogHash()
	filePath := filepath.Join(recapStoreDir(), inputSHA256+".ndjson")
	prior, quarantined := readRecapFile(filePath, parserVer, catalogHash)

	// Emit a recap Finding summarizing prior entries.
	if len(prior) > 0 {
		summary := make([]map[string]interface{}, 0, len(prior))
		for _, e := range prior {
			row := map[string]interface{}{
				"written_at":        e.Provenance.WrittenAt.UTC().Format(time.RFC3339),
				"parser_version":    e.Provenance.ParserVersion,
				"top_patterns":      e.TopPatterns,
				"top_movers":        e.TopMovers,
				"significant_count": e.SignificantCount,
				"window_from":       e.WindowFrom,
				"window_to":         e.WindowTo,
			}
			if e.CaptureQualityScore != nil {
				row["capture_quality_score"] = *e.CaptureQualityScore
			}
			summary = append(summary, row)
		}
		c.add(newFinding("recap", host, inputSHA256, map[string]interface{}{
			"input_sha256":      inputSHA256,
			"prior_runs_count":  len(prior),
			"prior_runs":        summary,
			"interpretation":    "Prior --auto runs against the same FTDC bytes. Use to see whether the diagnosis shape has shifted across parser revisions or to remember earlier investigations.",
		}))
		emitted++
	}
	// Surface quarantined entries so self-doubt is visible (per the
	// stored-inference contract in the roadmap).
	for _, q := range quarantined {
		c.add(newFinding("stored_inference_review", host, inputSHA256+"|"+q.reason+"|"+q.writtenAt, map[string]interface{}{
			"action":              "quarantined_recap_entry",
			"quarantined_reason":  q.reason,
			"file":                filePath,
			"entry_written_at":    q.writtenAt,
			"entry_parser_version": q.parserVer,
		}))
		emitted++
	}

	// Append the current run's compact summary to the file.
	current := buildCurrentRecap(c, host, parserVer, catalogHash, inputSHA256)
	_ = appendRecapEntry(filePath, current)
	return emitted, 0
}

// patternCatalogHash returns a deterministic SHA256 prefix over the
// sorted list of diagnosisPattern names. Changes when the catalog
// shape changes (added/removed patterns), which is the trigger for
// recap entries from prior versions to be quarantined.
func patternCatalogHash() string {
	names := make([]string, 0, len(diagnosisPatterns))
	for _, p := range diagnosisPatterns {
		names = append(names, p.name)
	}
	sort.Strings(names)
	h := sha256.Sum256([]byte(strings.Join(names, "|")))
	return hex.EncodeToString(h[:8])
}

// readRecapFile returns the still-valid entries + any quarantined
// entries with their quarantine reason. Validation rules mirror the
// stored-inference self-doubt machinery in the roadmap.
type quarantinedEntry struct {
	reason     string
	writtenAt  string
	parserVer  string
}

func readRecapFile(path, currentParserVer, currentCatalogHash string) (valid []recapEntry, quarantined []quarantinedEntry) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	now := time.Now()
	for sc.Scan() {
		var e recapEntry
		if err := json.Unmarshal(sc.Bytes(), &e); err != nil {
			continue
		}
		// Reject entries with a missing or unparseable written_at (zero value):
		// now.Sub(zero) exceeds recapStaleAge, which would trigger the stale
		// path even for fresh entries and could pass the version check silently.
		if e.Provenance.WrittenAt.IsZero() {
			quarantined = append(quarantined, quarantinedEntry{
				reason:    "missing_written_at",
				parserVer: e.Provenance.ParserVersion,
			})
			continue
		}
		writtenAtStr := e.Provenance.WrittenAt.UTC().Format(time.RFC3339)
		// 90-day staleness + version-skew check (mirrors stored-inference
		// self-doubt rules in the plan).
		if now.Sub(e.Provenance.WrittenAt) > recapStaleAge {
			if e.Provenance.ParserVersion != currentParserVer || e.Provenance.PatternCatalogHash != currentCatalogHash {
				quarantined = append(quarantined, quarantinedEntry{
					reason:    "stale_with_version_mismatch",
					writtenAt: writtenAtStr,
					parserVer: e.Provenance.ParserVersion,
				})
				continue
			}
		}
		if e.Provenance.ConfidenceRankAtWrite > 0 && e.Provenance.ConfidenceRankAtWrite < 0.5 {
			quarantined = append(quarantined, quarantinedEntry{
				reason:    "low_confidence_at_write",
				writtenAt: writtenAtStr,
				parserVer: e.Provenance.ParserVersion,
			})
			continue
		}
		valid = append(valid, e)
	}
	return valid, quarantined
}

// buildCurrentRecap assembles a compact summary of the current run.
func buildCurrentRecap(c *findingCollector, host, parserVer, catalogHash, inputSHA256 string) recapEntry {
	var entry recapEntry
	entry.Provenance = recapProvenance{
		WrittenAt:          time.Now().UTC(),
		ParserVersion:      parserVer,
		PatternCatalogHash: catalogHash,
		InputSHA256:        inputSHA256,
	}
	maxConf := 0.0
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		hf, _ := f["host"].(string)
		if hf != "" && hf != host {
			continue
		}
		switch k {
		case "diagnosis_candidate":
			pattern, _ := f["pattern_name"].(string)
			cr, _ := f["confidence_rank"].(float64)
			parentID, _ := f["id"].(string)
			entry.TopPatterns = append(entry.TopPatterns, recapPatternEntry{
				PatternName:     pattern,
				ConfidenceRank:  cr,
				ParentFindingID: parentID,
			})
			if cr > maxConf {
				maxConf = cr
			}
		case "capture_quality_score":
			if v, ok := f["score"].(float64); ok {
				score := int(math.Round(v))
				entry.CaptureQualityScore = &score
			}
		case "mover":
			// Cap to 10 top movers (sorted by score in the collector
			// already, so we just take the first 10 we see).
			if len(entry.TopMovers) >= 10 {
				continue
			}
			metric, _ := f["metric"].(string)
			if metric != "" {
				entry.TopMovers = append(entry.TopMovers, metric)
			}
		case "metric_score_table":
			// Both int and float64 occur in the same binary: the in-process
			// collector stores native int (countSignificantBH return type), but
			// findings replayed from a prior NDJSON file deliver integers as
			// float64 after JSON unmarshal.
			switch v := f["bh_significant_count"].(type) {
			case int:
				entry.SignificantCount = v
			case float64:
				entry.SignificantCount = int(v)
			}
		case "run_metadata":
			if eb, ok := f["effective_baseline"].(map[string]interface{}); ok {
				if s, ok := eb["from"].(string); ok {
					entry.WindowFrom = s
				}
			}
			if es, ok := f["effective_stress"].(map[string]interface{}); ok {
				if s, ok := es["to"].(string); ok {
					entry.WindowTo = s
				}
			}
		}
	}
	entry.Provenance.ConfidenceRankAtWrite = maxConf
	return entry
}

// appendRecapEntry appends one recap entry to the file (creating dir +
// file as needed). After append, if the file has > recapMaxEntriesPerFile
// lines, rolls the oldest into a `.archive` sibling.
func appendRecapEntry(path string, entry recapEntry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	line, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.Write(append(line, '\n')); err != nil {
		_ = f.Close()
		return err
	}
	_ = f.Close()

	// Roll oldest lines if over cap.
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) > recapMaxEntriesPerFile {
		archive := path + ".archive"
		dropped := lines[:len(lines)-recapMaxEntriesPerFile]
		kept := lines[len(lines)-recapMaxEntriesPerFile:]
		// Append the rolled lines to .archive.
		if af, err := os.OpenFile(archive, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
			_, _ = af.WriteString(strings.Join(dropped, "\n") + "\n")
			_ = af.Close()
		}
		// Rewrite the main file with only the kept tail.
		_ = os.WriteFile(path, []byte(strings.Join(kept, "\n")+"\n"), 0o644)
	}
	return nil
}
