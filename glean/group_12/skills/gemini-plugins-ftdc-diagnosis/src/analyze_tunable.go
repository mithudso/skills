// analyze_tunable.go — P1.NEW-A tunable_observation stage.
//
// Opt-in via `--auto-recommend-tunables`. When enabled, emits one
// `tunable_observation` Finding per applicable tunable whose supporting
// Finding kinds are all present AND no anti-recommendation gate fires
// against it. NEVER emits a recommended numeric value — only direction,
// supporting evidence, verification recipe, and guardrails.
//
// Bias contract (Tier-7 narration scaffold, allowlisted in BIAS_CONTROL.md):
//   - tunable_observation NAMES a tunable. The kind itself is editorial.
//   - The Finding refuses to commit to a numeric value.
//   - It cites supporting_finding_ids by id.
//   - It carries an explicit `disclaimer` field acknowledging it's editorial.
//   - The catalog of tunables (and anti-recommendations) is data, not
//     code logic; new entries go in tunable_guardrails.go.
package main

import (
	"sort"
)

// emitTunableObservations is the entry point invoked from the --auto
// pipeline when cfg.RecommendTunables is set. Pure post-stage; reads
// already-emitted Findings + the store's topology cues.
func emitTunableObservations(c *findingCollector, host string, store *seriesStore, cfg AutoConfig) (emitted int, skipped int) {
	if !cfg.RecommendTunables {
		return 0, 0
	}

	// Index emitted Findings by kind for cheap lookup. Carries the
	// Finding ID so we can cite supporting_finding_ids verbatim.
	kindIndex := map[string][]Finding{}
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		hf, _ := f["host"].(string)
		if hf != "" && hf != host {
			continue
		}
		kindIndex[k] = append(kindIndex[k], f)
	}

	// Topology inference reuses the existing gates + the new paliRole
	// helper (D6). We pre-compute once.
	topo := inferTunableTopology(store, kindIndex)

	// Sample-count gate (E.9: capture_too_short). Pull from the very
	// run_metadata Finding the parser already emitted.
	sampleCount := 0
	for _, f := range kindIndex["run_metadata"] {
		if n, ok := toIntField(f["sample_count"]); ok {
			sampleCount = n
			break
		}
	}

	// Determine which anti-recommendation gates are firing this run.
	firedGates := map[string]string{} // code → narration
	for _, gate := range antiRecommendationCatalog {
		fired := antiRecGateFires(gate, kindIndex, sampleCount)
		if fired {
			firedGates[gate.code] = gate.rationale
		}
	}

	// Evaluate each tunable in the catalog.
	for _, rule := range tunableCatalog {
		if !topologyApplicable(rule.topologyApplicability, topo) {
			skipped++
			continue
		}
		// Supporting Finding kinds all present?
		supportingIDs := []string{}
		allRequiredPresent := true
		for _, kind := range rule.requiredFindingKinds {
			fs, ok := kindIndex[kind]
			if !ok || len(fs) == 0 {
				allRequiredPresent = false
				break
			}
			// Cite the FIRST instance per kind for determinism. Many
			// fires per kind would explode supporting_finding_ids; one
			// per kind is the bias-control-friendly choice (downstream
			// can grep for more).
			for _, f := range fs {
				if fid, ok := f["id"].(string); ok && fid != "" {
					supportingIDs = append(supportingIDs, fid)
					break
				}
			}
		}
		if !allRequiredPresent {
			skipped++
			continue
		}

		// Confirm metric paths the verification recipe references exist
		// in the current capture. D3 verification gate.
		unresolvable := []string{}
		for _, mp := range rule.requiredMetricNames {
			if _, ok := store.nameToID[mp]; !ok {
				unresolvable = append(unresolvable, mp)
			}
		}

		// Apply any firing anti-recommendation gates that target this tunable.
		antiCodes := []string{}
		for code := range firedGates {
			gate := findAntiGate(code)
			if gate == nil {
				continue
			}
			if len(gate.tunablesBlocked) == 0 || stringInSlice(rule.tunable, gate.tunablesBlocked) {
				antiCodes = append(antiCodes, code)
			}
		}
		sort.Strings(antiCodes)

		direction := rule.directionWhenObserved
		var refusalRationales []string
		if len(antiCodes) > 0 {
			direction = "refuse"
			for _, code := range antiCodes {
				refusalRationales = append(refusalRationales, firedGates[code])
			}
		}
		validity := "complete"
		if len(unresolvable) > 0 && direction != "refuse" {
			validity = "partial"
		}

		fields := map[string]interface{}{
			"tunable":                rule.tunable,
			"tunable_kind":           rule.tunableKind,
			"topology_applicability": rule.topologyApplicability,
			"direction":              direction,
			"supporting_finding_ids": supportingIDs,
			"guardrails":             rule.guardrails,
			"verification_recipe":    rule.verificationRecipe,
			"verification_expected_if_helped":     rule.verificationHelped,
			"verification_expected_if_not_helped": rule.verificationNotHelped,
			"verification_validity":  validity,
			"anti_recommendation_codes": antiCodes,
			"disclaimer": "tunable_observation Findings are editorial. The parser identifies the structural signal; the operator decides the value. No recommended numeric value is emitted by design.",
		}
		if len(unresolvable) > 0 {
			fields["unresolvable_metric_references"] = unresolvable
		}
		if len(refusalRationales) > 0 {
			fields["refusal_rationales"] = refusalRationales
		}
		c.add(newFinding("tunable_observation", host, rule.tunable, fields))
		emitted++
	}
	return emitted, skipped
}

// inferTunableTopology returns a topology label this run applies to. Combines
// the existing isReplicaSet / isPALI / isSharded gates with paliRole
// (D6) for primary/standby distinction.
func inferTunableTopology(store *seriesStore, kindIndex map[string][]Finding) string {
	if isPALICapture(store) {
		role, _ := paliRole(store)
		switch role {
		case PALIRolePrimary:
			return "pali_primary"
		case PALIRoleStandby:
			return "pali_standby"
		default:
			return "pali_primary" // default PALI to primary if role indeterminate
		}
	}
	if isShardedCapture(store) {
		return "config_server" // shard servers also classify under this for catalog purposes
	}
	if isReplicaSetCapture(store) {
		// Without leader_transition / replSetGetStatus details we can't
		// disambiguate primary vs secondary cheaply here; default to
		// replica_set_primary which is the most permissive choice.
		_ = kindIndex
		return "replica_set_primary"
	}
	return "standalone"
}

// topologyApplicable returns true when topo is in the allowed list.
func topologyApplicable(allowed []string, topo string) bool {
	for _, a := range allowed {
		if a == topo {
			return true
		}
	}
	return false
}

// antiRecGateFires evaluates whether an anti-recommendation gate's
// structural conditions are met by the current run.
func antiRecGateFires(gate antiRecommendationGate, kindIndex map[string][]Finding, sampleCount int) bool {
	// Special-case the data-only gates.
	switch gate.code {
	case "signal_below_noise_floor":
		// Any Finding with |cohens_d| ≥ 0.2 (the noise floor per BIAS_CONTROL §B.5)?
		// If yes, signal is NOT below the noise floor.
		for _, f := range kindIndex["mover"] {
			if v, ok := f["cohens_d"].(float64); ok && absFloat(v) >= 0.2 {
				return false
			}
		}
		// If we have no movers at all OR none above the floor, the gate
		// fires. Captures with zero mover Findings inherently have no
		// significant signal.
		return true
	case "capture_too_short":
		return sampleCount > 0 && sampleCount < 600
	}

	// Standard present-AND / absent-AND rule.
	for _, k := range gate.requiredKindsPresent {
		if len(kindIndex[k]) == 0 {
			return false
		}
	}
	// requiredKindsAbsent in the structure means "this Finding kind
	// must NOT fire for the gate to apply." (We chose this awkward
	// formulation to keep the data declarative.)
	for _, k := range gate.requiredKindsAbsent {
		if len(kindIndex[k]) == 0 {
			// kind is absent → satisfies the "must not fire" check
			continue
		}
		return false
	}
	return true
}

// findAntiGate is a small lookup helper.
func findAntiGate(code string) *antiRecommendationGate {
	for i := range antiRecommendationCatalog {
		if antiRecommendationCatalog[i].code == code {
			return &antiRecommendationCatalog[i]
		}
	}
	return nil
}

// stringInSlice is the smallest possible "contains" check. We define it
// here rather than pull a util because the package keeps deps minimal.
func stringInSlice(s string, set []string) bool {
	for _, x := range set {
		if x == s {
			return true
		}
	}
	return false
}

