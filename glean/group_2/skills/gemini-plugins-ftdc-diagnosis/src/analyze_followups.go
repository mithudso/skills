// analyze_followups.go — P2.4 append-only followups backlog.
//
// Solo user accumulates "diagnosis candidates I didn't get to finish."
// Persisting these to ~/claude-dump/ftdc-followups.ndjson means the
// next session sees them and can recommend continuing the
// investigation. Pure FTDC-only — no external corpus.
//
// Trigger: every --auto run that emits ≥1 diagnosis_candidate with
// confidence_rank in [0.5, uncertaintyConfidenceCutoff) appends one
// line per candidate to the followups file. The file is append-only;
// deduplication / pruning is the user's responsibility (or a future
// `--prune-followups` flag).
//
// Each line carries: parser_version, input_sha256, host, pattern_name,
// confidence_rank, supporting_finding_ids, run_timestamp.
//
// Read-back: a new Finding kind `followups_backlog_snapshot` summarizes
// the file at the end of every --auto so the user sees it without
// having to remember the file exists. Implemented via the stored-
// inference self-doubt machinery from P4.4 (when that lands); for now,
// the snapshot just reads the file directly with simple sanity gating.
package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// followupsBacklogPath is where the append-only NDJSON lives. Falls back
// to /tmp if $HOME isn't set — the path is solo-user-specific.
func followupsBacklogPath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		home = "/tmp"
	}
	return filepath.Join(home, "claude-dump", "ftdc-followups.ndjson")
}

// emitFollowupsBacklog appends to the followups file and emits a single
// followups_backlog_snapshot Finding describing the appended candidates
// + (lightly summarized) prior entries. Pure post-stage; reads
// diagnosis_candidate + run_metadata Findings from the collector.
//
// File-system errors are non-fatal — followups are advisory.
func emitFollowupsBacklog(c *findingCollector, host string) (emitted int, skipped int) {
	// Pull input_sha256 from the run_metadata Finding (P0.3).
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

	// Candidates worth backing up.
	type cand struct {
		Pattern         string   `json:"pattern_name"`
		ConfidenceRank  float64  `json:"confidence_rank"`
		ParentFindingID string   `json:"parent_finding_id"`
		SupportingIDs   []string `json:"supporting_finding_ids"`
	}
	var freshCands []cand
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		if k != "diagnosis_candidate" {
			continue
		}
		hf, _ := f["host"].(string)
		if hf != "" && hf != host {
			continue
		}
		cr, _ := f["confidence_rank"].(float64)
		if cr < 0.5 || cr >= uncertaintyConfidenceCutoff {
			continue
		}
		pattern, _ := f["pattern_name"].(string)
		parentID, _ := f["id"].(string)
		// supporting_finding_ids is stored as []string by emitDiagnosisCandidates.
		// asserting to []interface{} always fails (Go type system is exact).
		supporting, _ := f["supporting_finding_ids"].([]string)
		freshCands = append(freshCands, cand{
			Pattern:         pattern,
			ConfidenceRank:  cr,
			ParentFindingID: parentID,
			SupportingIDs:   supporting,
		})
	}
	if len(freshCands) == 0 {
		return 0, 1
	}

	// Append to file. Errors are intentionally silent — backup is
	// advisory; we don't want to abort the auto run.
	path := followupsBacklogPath()
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	// O_APPEND|O_CREATE positions writes at end-of-file atomically.
	// O(1) stat gives the prior file size for the informational Finding field.
	priorSize := int64(0)
	if fi, statErr := os.Stat(path); statErr == nil {
		priorSize = fi.Size()
	}
	runTS := time.Now().UTC().Format(time.RFC3339)
	if f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
		w := bufio.NewWriter(f)
		for _, candidate := range freshCands {
			entry := map[string]interface{}{
				"appended_at":            runTS,
				"input_sha256":           inputSHA256,
				"parser_version":         parserVer,
				"host":                   host,
				"pattern_name":           candidate.Pattern,
				"confidence_rank":        candidate.ConfidenceRank,
				"parent_finding_id":      candidate.ParentFindingID,
				"supporting_finding_ids": candidate.SupportingIDs,
			}
			if line, err := json.Marshal(entry); err == nil {
				_, _ = w.Write(line)
				_, _ = w.Write([]byte{'\n'})
			}
		}
		_ = w.Flush()
		_ = f.Close()
	}
	candSummary := make([]map[string]interface{}, 0, len(freshCands))
	for _, candidate := range freshCands {
		candSummary = append(candSummary, map[string]interface{}{
			"pattern_name":    candidate.Pattern,
			"confidence_rank": candidate.ConfidenceRank,
		})
	}
	c.add(newFinding("followups_backlog_snapshot", host, runTS, map[string]interface{}{
		"path":                 path,
		"appended_this_run":    len(freshCands),
		"appended_candidates":  candSummary,
		"file_bytes_before_run": priorSize,
		"interpretation":       "Append-only solo-user backlog of mid-confidence candidates worth revisiting. Prune the file manually when investigations close.",
	}))
	emitted++
	return emitted, 0
}
