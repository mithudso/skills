// analyze_uncertainty.go — P2.2 diagnosis_uncertainty Finding.
//
// A diagnosis_candidate Finding carries `discriminating_kinds_needed`
// — Finding kinds the parser would want to see (or not see) to push
// confidence higher or to rule the pattern out. That's already useful
// to the LLM analyzer, but it's implicit: the user has to know what each
// discriminator kind MEANS for the candidate.
//
// diagnosis_uncertainty makes it actionable. For each candidate whose
// confidence_rank < 0.9, emit a Finding that enumerates each pending
// discriminator with: kind, what its presence would do (push higher /
// rule out), and a hint at how the user could elicit it (re-capture,
// filter the existing capture, check a sibling host, etc.).
//
// Solo-user motivation: there's no colleague to bounce "what would
// settle this?" off. The parser surfaces the next-step list itself.
package main

import "strings"

// uncertaintyConfidenceCutoff is the confidence_rank above which a
// diagnosis_candidate is considered "strong enough" that surfacing its
// remaining discriminators would just add noise. Matches the analyzer
// SKILL.md ≥0.9 "strong match" tier.
const uncertaintyConfidenceCutoff = 0.9

// emitDiagnosisUncertainty reads diagnosis_candidate Findings from the
// collector and emits one diagnosis_uncertainty Finding per candidate
// below the cutoff. Pure post-stage; runs after diagnosis_candidate
// (which carries the discriminating_kinds_needed list).
func emitDiagnosisUncertainty(c *findingCollector, host string) (emitted int, skipped int) {
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
		if cr >= uncertaintyConfidenceCutoff {
			skipped++
			continue
		}
		pattern, _ := f["pattern_name"].(string)
		parentID, _ := f["id"].(string)
		// discriminating_kinds_needed is stored as []string by emitDiagnosisCandidates.
		// asserting to []interface{} always fails (Go type system is exact).
		discKinds, ok := f["discriminating_kinds_needed"].([]string)
		if !ok || len(discKinds) == 0 {
			skipped++
			continue
		}

		// Build per-discriminator action hints. The mapping is small and
		// static (parser-side editorial scaffolding, analogous to
		// hypothesis_test recipe templates) — analyzer narrates.
		actions := make([]map[string]string, 0, len(discKinds))
		for _, d := range discKinds {
			hint := actionHintForDiscriminator(d, pattern)
			actions = append(actions, map[string]string{
				"discriminator_kind": d,
				"effect_if_observed": hint.effect,
				"how_to_elicit":      hint.howToElicit,
			})
		}

		idKey := pattern + "|uncertainty"
		c.add(newFinding("diagnosis_uncertainty", host, idKey, map[string]interface{}{
			"parent_finding_id":         parentID,
			"parent_pattern_name":       pattern,
			"current_confidence_rank":   cr,
			"target_confidence_rank":    uncertaintyConfidenceCutoff,
			"pending_discriminators":    discKinds,
			"discriminator_actions":     actions,
			"rationale": "Pattern matched below the strong-match cutoff (0.9). Observing any of " +
				"the pending discriminators would either push the candidate higher or rule it out.",
		}))
		emitted++
	}
	return emitted, skipped
}

// discriminatorHint pairs the effect of seeing a discriminator with a
// brief hint about how the user could elicit it.
type discriminatorHint struct {
	effect       string
	howToElicit  string
}

// actionHintForDiscriminator returns a static editorial pair per
// (discriminator kind, parent pattern). Keep the strings short and
// neutral — they're scaffolding, not narration.
func actionHintForDiscriminator(disc, parentPattern string) discriminatorHint {
	// Parent-specific overrides come first; fall back to the generic
	// table.
	if strings.Contains(parentPattern, "pali") || strings.HasPrefix(parentPattern, "standby_") || strings.HasPrefix(parentPattern, "materialization_") || strings.HasPrefix(parentPattern, "pageserver_") {
		switch disc {
		case "materialization_lag":
			return discriminatorHint{
				effect:      "Observing materialization_lag would push the candidate toward primary_pali_log_pressure / materialization_apply_split.",
				howToElicit: "Re-capture during a window with high write throughput on primary; filter --filter 'disagg.phylog.*'.",
			}
		case "regression":
			return discriminatorHint{
				effect:      "Observing regression (Tier-D vs baseline) would confirm this capture diverges from prior PALI runs.",
				howToElicit: "Provide --auto-baseline-file with a recent good-run capture's NDJSON.",
			}
		case "io_routing":
			return discriminatorHint{
				effect:      "Observing io_routing would push toward pageserver_contention vs primary-side bottleneck.",
				howToElicit: "Re-capture while both nodes are reading from the same page server; filter --filter 'disagg.pageServer.*'.",
			}
		case "counter_wrap":
			return discriminatorHint{
				effect:      "Observing counter_wrap would push the candidate toward standby_read_log_wedge — the standby's counters stopped advancing.",
				howToElicit: "Re-capture immediately following a seal_event; filter --filter 'disagg.phylog.*'.",
			}
		}
	}
	switch disc {
	case "regression":
		return discriminatorHint{
			effect:      "Observing a regression Finding would corroborate that this capture diverges from baseline.",
			howToElicit: "Provide --auto-baseline-file pointing at a known-good prior capture's NDJSON.",
		}
	case "counter_wrap":
		return discriminatorHint{
			effect:      "Observing counter_wrap would explain abrupt counter resets as wraparound rather than process restart.",
			howToElicit: "Inspect the metric_score_table for monotone counters; re-capture if needed.",
		}
	case "alert_threshold":
		return discriminatorHint{
			effect:      "Observing an alert_threshold Finding would corroborate the candidate via Atlas-class event detection.",
			howToElicit: "Check Atlas alert history for the capture window.",
		}
	case "missing_signal":
		return discriminatorHint{
			effect:      "Observing missing_signal would corroborate degraded capture quality.",
			howToElicit: "Confirm FTDC was uninterrupted; check capture_quality_score.",
		}
	case "synthetic_event":
		return discriminatorHint{
			effect:      "Observing synthetic_event would confirm multiple detectors fire on the same timestamp cluster.",
			howToElicit: "Already inspected by the parser; if absent, the temporal co-occurrence wasn't strong enough.",
		}
	case "checkpoint_stall":
		return discriminatorHint{
			effect:      "Observing checkpoint_stall would corroborate WT-side back-pressure.",
			howToElicit: "Filter --filter 'wiredTiger.transaction.transaction checkpoint' over the suspected window.",
		}
	}
	return discriminatorHint{
		effect:      "Observing this kind would corroborate or contradict the current candidate.",
		howToElicit: "Re-capture or filter the existing capture for the relevant metric family.",
	}
}
