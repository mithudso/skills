// Tier R.1: emit kind:"diagnosis_candidate" Findings by matching the
// collector's Finding stream against the diagnosis pattern catalog.
//
// Runs AFTER all Tier 0-7 + Tier G stages emit, BEFORE applyOutputGates.
// Reads emitted Findings, NEVER consumes raw metric series — patterns
// key off Finding *shape* (kinds), not metric paths.
//
// Confidence ranking is mechanical:
//
//	confidence_rank = required_present / (required_total + excluded_present)
//	              + 0.05 * supporting_present
//
// Clamped to [0, 1]. Renamed (not "confidence") per Tier B.13 to prevent
// downstream consumers from treating it as a calibrated probability.

package main

import (
	"fmt"
	"sort"
	"strings"
)

func emitDiagnosisCandidates(c *findingCollector, host string, store *seriesStore) (emitted, skipped int) {
	// Index Findings by kind for fast presence checks.
	byKind := map[string][]string{} // kind → []finding_id
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		if k == "" {
			continue
		}
		fid, _ := f["id"].(string)
		byKind[k] = append(byKind[k], fid)
	}
	// Domain-stage gating: PALI / replication / sharded / host_os.
	domainHints := detectDomainHints(store)

	for _, p := range diagnosisPatterns {
		if p.domainStageGate != "" && !domainHints[p.domainStageGate] {
			skipped++
			continue
		}
		// Check required kinds are ALL present.
		requiredPresent := 0
		var requiredIDs []string
		for _, rk := range p.requiredKinds {
			if ids, ok := byKind[rk]; ok && len(ids) > 0 {
				requiredPresent++
				requiredIDs = append(requiredIDs, ids[0])
			}
		}
		if requiredPresent != len(p.requiredKinds) {
			skipped++
			continue
		}
		// Count supporting (optional) kinds.
		supportingCount := 0
		var supportingIDs []string
		for _, sk := range p.supportingKinds {
			if ids, ok := byKind[sk]; ok && len(ids) > 0 {
				supportingCount++
				supportingIDs = append(supportingIDs, ids[0])
			}
		}
		// Check excluded (contradicting) kinds.
		excludedCount := 0
		var contradictingIDs []string
		for _, ek := range p.excludedKinds {
			if ids, ok := byKind[ek]; ok && len(ids) > 0 {
				excludedCount++
				contradictingIDs = append(contradictingIDs, ids[0])
			}
		}
		// Confidence rank.
		denom := float64(len(p.requiredKinds) + excludedCount)
		if denom <= 0 {
			denom = 1
		}
		score := float64(requiredPresent)/denom + 0.05*float64(supportingCount)
		if score > 1 {
			score = 1
		}
		// P4.2 user-pattern confidence cap: patterns loaded from
		// ~/.config/ftdc-parser/patterns.d/ are user-authored editorial
		// scaffolds. We treat them as plausible-not-strong by capping
		// confidence_rank at userPatternConfidenceCap until fixtures
		// validate them. Built-in patterns are uncapped.
		patternSource := "builtin"
		if userPatternConfidenceCapApplies(p.name, userLoadedPatternNames) {
			patternSource = "user_yaml"
			if score > userPatternConfidenceCap {
				score = userPatternConfidenceCap
			}
		}
		// discriminating_kinds_needed: only those NOT already in byKind.
		var needed []string
		for _, dk := range p.discriminatingKindsNeeded {
			if _, ok := byKind[dk]; !ok {
				needed = append(needed, dk)
			}
		}
		sort.Strings(needed)
		sort.Strings(requiredIDs)
		sort.Strings(supportingIDs)
		sort.Strings(contradictingIDs)

		c.add(newFinding("diagnosis_candidate", host, p.name, map[string]interface{}{
			"pattern_name":                p.name,
			"required_finding_ids":        requiredIDs,
			"supporting_finding_ids":      supportingIDs,
			"contradicting_finding_ids":   contradictingIDs,
			"discriminating_kinds_needed": needed,
			"confidence_rank":             roundFloat(score, 4),
			"rationale":                   fmt.Sprintf(p.rationaleTemplate, supportingCount, excludedCount),
			"required_kinds_present":      requiredPresent,
			"required_kinds_total":        len(p.requiredKinds),
			"supporting_kinds_present":    supportingCount,
			"excluded_kinds_present":      excludedCount,
			"pattern_source":              patternSource,
		}))
		emitted++
	}
	return emitted, skipped
}

// detectDomainHints returns a presence map for known Tier-5 domain stages.
// Uses STRUCTURAL prefix tests on metric names (per BIAS_CONTROL.md
// allowance for Tier-5 domain stages).
//
// Implementation note: switched from hand-counted length-prefix slicing
// (which had a CRITICAL off-by-two bug — `"serverStatus.metrics.disagg"` is
// 27 chars not 25, so the original `n[:25]` check matched a different
// substring) to `strings.HasPrefix`, which is both correct and clearer.
func detectDomainHints(store *seriesStore) map[string]bool {
	out := map[string]bool{}
	if store == nil {
		return out
	}
	for _, n := range store.names {
		switch {
		case strings.HasPrefix(n, "serverStatus.metrics.disagg"):
			out["pali"] = true
		case strings.HasPrefix(n, "serverStatus.repl."):
			out["replication"] = true
		case strings.HasPrefix(n, "serverStatus.shardingStatistics."):
			out["sharded"] = true
		case strings.HasPrefix(n, "systemMetrics."):
			out["host_os"] = true
		}
	}
	return out
}
