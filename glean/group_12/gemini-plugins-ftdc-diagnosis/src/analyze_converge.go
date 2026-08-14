// analyze_converge.go — P4.1 autonomous diagnosis-loop planning.
//
// Auto-triggers when the run ends with the top diagnosis_candidate at
// confidence_rank < 0.9 AND a sibling capture directory exists adjacent
// to the input path (parser walks `..` for peer dirs).
//
// Memory constraint (see `feedback_single_ftdc_parser_at_a_time.md`):
// running a second full --auto pipeline in-process on a sibling capture
// while the primary store is still alive doubles peak RAM, which OOMs
// the host on large fixtures. Until per-host store eviction can be
// driven by an outer loop, this stage emits a PLAN (kind:converge_plan)
// describing what the autonomous loop WOULD do — sibling captures
// discovered, recipes extracted from existing hypothesis_test Findings,
// expected outcomes — rather than executing the loop in-process.
//
// The plan is structurally complete (Gregory can run each recipe by
// hand and re-pipe to ftdc_parser; the recipes are already mechanically
// parseable). A future P4.1.b can wire a true outer-loop driver that
// invokes ftdc_parser as a subprocess on each sibling.
//
// kind:"convergence_failed" is emitted when conditions for the plan
// itself are unmet (no siblings, no low-confidence candidate, etc).
// Honest by design — the parser declares its own failure mode.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// convergeConfidenceCutoff matches uncertaintyConfidenceCutoff — at or
// above, the candidate is "strong enough" and the loop is unnecessary.
const convergeConfidenceCutoff = uncertaintyConfidenceCutoff

// convergeMaxSiblings caps the number of sibling capture directories we
// emit in the plan. Beyond this, the user is likely better served by
// the chaos-sweep stage (P4.3) which is purpose-built for many runs.
const convergeMaxSiblings = 5

// emitConvergePlan reads the current run's diagnosis_candidate +
// hypothesis_test Findings + the input path, walks the parent
// directory for sibling captures, and emits one converge_plan Finding
// summarizing what an outer-loop driver could do. Pure read-only on
// the filesystem — no subprocess execution.
func emitConvergePlan(c *findingCollector, host string, store *seriesStore, cfg AutoConfig) (emitted int, skipped int) {
	if cfg.InputPath == "" {
		return 0, 1
	}

	// Find the highest-confidence candidate below the strong-match threshold
	// (0.9) that's still above the "informational" floor (0.3) — that's the
	// most useful target for convergence work.
	var (
		targetCandidate Finding
		targetCR        float64
		targetID        string
	)
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
		if cr >= convergeConfidenceCutoff {
			continue
		}
		if cr < 0.3 {
			continue
		}
		if targetCandidate == nil || cr > targetCR {
			targetCandidate = f
			targetCR = cr
			if id, ok := f["id"].(string); ok {
				targetID = id
			}
		}
	}
	if targetCandidate == nil {
		c.add(newFinding("convergence_failed", host, "no_candidate", map[string]interface{}{
			"reason_code": "no_low_confidence_candidate",
			"reason":      "No diagnosis_candidate fired with confidence_rank in [0.3, 0.9). Either the top diagnosis is already strong (≥0.9) or no shape matched.",
		}))
		emitted++
		return emitted, 0
	}

	// Find the hypothesis_test attached to the chosen candidate.
	var matchingHypothesis Finding
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		if k != "hypothesis_test" {
			continue
		}
		if pid, ok := f["parent_finding_id"].(string); ok && pid == targetID {
			matchingHypothesis = f
			break
		}
	}

	// Discover sibling directories chronologically.
	siblings := discoverSiblingCaptures(cfg.InputPath)
	if len(siblings) == 0 {
		c.add(newFinding("convergence_failed", host, "no_siblings", map[string]interface{}{
			"reason_code": "no_sibling_captures",
			"reason":      "No peer capture directories found alongside the input path. Convergence requires multiple captures to iterate against.",
			"input_path":  cfg.InputPath,
		}))
		emitted++
		return emitted, 0
	}
	if len(siblings) > convergeMaxSiblings {
		siblings = siblings[:convergeMaxSiblings]
	}

	// Extract recipe from the hypothesis_test (if present) for each
	// sibling. The recipe template uses "<path>" or implicit "next
	// capture" — we substitute the sibling's actual path so the plan
	// emits ready-to-run commands.
	recipeTpl := ""
	expectedIfConfirmed := ""
	expectedIfRuledOut := ""
	if matchingHypothesis != nil {
		recipeTpl, _ = matchingHypothesis["recipe"].(string)
		expectedIfConfirmed, _ = matchingHypothesis["expected_outcome_if_confirmed"].(string)
		expectedIfRuledOut, _ = matchingHypothesis["expected_outcome_if_ruled_out"].(string)
	}

	plannedSteps := make([]map[string]string, 0, len(siblings))
	for _, sib := range siblings {
		cmd := substituteSiblingPath(recipeTpl, sib)
		plannedSteps = append(plannedSteps, map[string]string{
			"sibling_path": sib,
			"recipe":       cmd,
		})
	}

	pattern, _ := targetCandidate["pattern_name"].(string)
	c.add(newFinding("converge_plan", host, pattern, map[string]interface{}{
		"target_candidate_id":      targetID,
		"target_pattern_name":      pattern,
		"current_confidence_rank":  targetCR,
		"target_confidence_rank":   convergeConfidenceCutoff,
		"sibling_capture_count":    len(siblings),
		"planned_steps":            plannedSteps,
		"expected_if_confirmed":    expectedIfConfirmed,
		"expected_if_ruled_out":    expectedIfRuledOut,
		"loop_executed":            false,
		"loop_execution_note":      "v1 emits the plan only. Memory constraints prevent in-process execution of a second full --auto pass on each sibling. Run the planned recipes by hand and re-invoke ftdc_parser; a future driver subcommand can wire the loop end-to-end.",
		"hardcoded_self_doubt": []string{
			"refuse to recurse on captures with capture_quality_score < 70",
			"refuse to use Findings with q_value > 0.05 when re-scoring",
			"refuse to escalate confidence_rank by more than 0.2 per iteration",
		},
	}))
	emitted++
	return emitted, 0
}

// discoverSiblingCaptures returns peer capture directories adjacent to
// the input path, sorted by mtime ascending (oldest first — natural
// chronological order for an investigation timeline).
func discoverSiblingCaptures(inputPath string) []string {
	parent := filepath.Dir(inputPath)
	if parent == "" || parent == "." {
		return nil
	}
	entries, err := os.ReadDir(parent)
	if err != nil {
		return nil
	}
	primaryAbs, _ := filepath.Abs(inputPath)
	type cand struct {
		path  string
		mtime int64
	}
	var cands []cand
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		full := filepath.Join(parent, e.Name())
		abs, _ := filepath.Abs(full)
		if abs == primaryAbs {
			continue
		}
		// Cheap sanity check: directory must contain at least one file
		// that doesn't look like a golden / changelog.
		if !looksLikeCaptureDir(full) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		cands = append(cands, cand{path: full, mtime: info.ModTime().Unix()})
	}
	sort.Slice(cands, func(i, j int) bool {
		if cands[i].mtime != cands[j].mtime {
			return cands[i].mtime < cands[j].mtime
		}
		return cands[i].path < cands[j].path
	})
	out := make([]string, 0, len(cands))
	for _, c := range cands {
		out = append(out, c.path)
	}
	return out
}

// looksLikeCaptureDir cheaply heuristics whether a directory is a
// plausible FTDC capture. Same gates as listFTDCFixtures in the
// determinism test: at least one file that isn't a golden / changelog.
func looksLikeCaptureDir(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		n := e.Name()
		if n == "expected_findings.ndjson" || n == "CHANGELOG.md" {
			continue
		}
		return true
	}
	return false
}

// substituteSiblingPath rewrites the recipe template's path argument
// (typically "/path/to/next/capture/" or a bare placeholder) to point
// at the sibling. The template doesn't always have a clean substitution
// point; we fall back to appending the path when no obvious slot exists.
func substituteSiblingPath(tpl, siblingPath string) string {
	if tpl == "" {
		return fmt.Sprintf("ftdc_parser %s --auto", siblingPath)
	}
	// Common recipe shapes:
	//   "ftdc_parser /path/to/next/capture/ --auto --filter '...'"
	//   "ftdc_parser <path> --auto"
	// Substitute the FIRST path-like token after the binary name.
	out := tpl
	out = strings.Replace(out, "/path/to/next/capture/", siblingPath, 1)
	out = strings.Replace(out, "<path>", siblingPath, 1)
	out = strings.Replace(out, "<next-capture>", siblingPath, 1)
	return out
}
