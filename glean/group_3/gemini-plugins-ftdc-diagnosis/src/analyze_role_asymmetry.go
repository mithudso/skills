// analyze_role_asymmetry.go — cross-host PALI role asymmetry detector.
//
// P0.NEW-PALI-2: cluster A and chaos bug #6 (per Gregory's project memory)
// are fundamentally asymmetric-state bugs: primary applies log records,
// standby discards them. The generic `cross_host_drift` Finding catches
// ratio extremes but does not say "this metric diverged BETWEEN
// primary and standby specifically." This stage closes that gap.
//
// Detection rule (D2 in the P0 plan):
//   For each unordered host pair (A, B), for each metric M present on
//   both hosts:
//     1. M must NOT be in host_only_metrics (we want the SAME metric on
//        both sides; cross-set asymmetry is a different Finding kind).
//     2. M is in the PALI namespace (serverStatus.metrics.disagg.*).
//        v1 scope narrowing — keeps false positives bounded.
//     3. Cohen's d on (fullMean, fullSD) between A and B has |d| >= 0.8
//        (large effect per BIAS_CONTROL §B.5 threshold).
//   Emit kind:"role_asymmetry" when paliRole identified both sides
//   AND they are different roles (primary↔standby).
//   Emit kind:"node_pair_asymmetry" when role identification was
//   ambiguous; the Finding carries only host pair + evidence, no role
//   field (preserves bias-control posture).
//
// fullMean and fullSD are pre-computed during chunk decode and survive
// the column-store eviction in runAutoMultiStreaming, so this detector
// works even on the streaming multi-host path.
package main

import (
	"math"
	"sort"
	"strings"
)

// roleAsymmetryPaliPrefix scopes the detector to PALI metrics in v1.
// Generalizing to all intersection metrics would require additional
// noise gating; PALI is the narrow scope that matters for Gregory's
// active bugs.
const roleAsymmetryPaliPrefix = "serverStatus.metrics.disagg."

// roleAsymmetryDThreshold is the |Cohen's d| above which a metric pair
// is considered a role asymmetry. 0.8 is the "large effect" threshold
// from BIAS_CONTROL Tier B.5; using a higher bar than the default 0.5
// keeps cross-host comparisons biased toward strong signals.
const roleAsymmetryDThreshold = 0.8

// emitRoleAsymmetry inspects every unordered host pair, computes
// Cohen's d on shared PALI metrics, and emits role_asymmetry /
// node_pair_asymmetry Findings. Pure function of hostStores +
// hostRoles; no globals.
//
// Returns (emitted, skipped). Skipped includes cases where no shared
// PALI metric was found or there's a single host (nothing to compare).
func emitRoleAsymmetry(c *findingCollector, hostStores map[string]*seriesStore, hostRoles map[string]PALIRole, cfg AutoConfig) (emitted int, skipped int) {
	if len(hostStores) < 2 {
		return 0, 1
	}

	// Sort host names so pair iteration is deterministic.
	hosts := make([]string, 0, len(hostStores))
	for h := range hostStores {
		hosts = append(hosts, h)
	}
	sort.Strings(hosts)

	for i := 0; i < len(hosts); i++ {
		for j := i + 1; j < len(hosts); j++ {
			hA, hB := hosts[i], hosts[j]
			sA, sB := hostStores[hA], hostStores[hB]
			if sA == nil || sB == nil {
				continue
			}

			// Determine the kind for this pair upfront — Findings emitted
			// in this iteration all share kind based on role inference.
			rA, hasA := hostRoles[hA]
			rB, hasB := hostRoles[hB]
			rolePairKnown := hasA && hasB && rA != PALIRoleUnknown && rB != PALIRoleUnknown && rA != rB
			pairKind := "role_asymmetry"
			if !rolePairKnown {
				pairKind = "node_pair_asymmetry"
			}

			// Walk metrics deterministically: sorted by name. Intersection
			// = metrics in both hosts. Skip non-PALI metrics in v1.
			names := make([]string, 0, len(sA.nameToID))
			for n := range sA.nameToID {
				if _, ok := sB.nameToID[n]; !ok {
					continue
				}
				if !strings.HasPrefix(n, roleAsymmetryPaliPrefix) {
					continue
				}
				names = append(names, n)
			}
			sort.Strings(names)

			pairHadEmission := false
			for _, name := range names {
				idA := sA.nameToID[name]
				idB := sB.nameToID[name]
				if sA.isDegenerate[idA] && sB.isDegenerate[idB] {
					continue
				}

				meanA := sA.fullMean[idA]
				meanB := sB.fullMean[idB]
				sdA := sA.fullSD[idA]
				sdB := sB.fullSD[idB]

				// Cohen's d with pooled SD. When either SD is zero we
				// can still compute a direction; we skip emission if
				// both SDs are zero (no spread to compare).
				if sdA == 0 && sdB == 0 {
					continue
				}
				nA := fullNOrDefault(sA, idA)
				nB := fullNOrDefault(sB, idB)
				pooledSD := pooledStandardDeviation(sdA, nA, sdB, nB)
				if pooledSD <= 0 {
					continue
				}
				d := (meanA - meanB) / pooledSD
				if math.IsNaN(d) || math.IsInf(d, 0) {
					continue
				}
				if math.Abs(d) < roleAsymmetryDThreshold {
					continue
				}

				// For role_asymmetry: derive direction from which host is
				// the primary, not from alphabetical ordering of hosts.
				// hA < hB lexicographically, but the primary could be either.
				var direction string
				if rolePairKnown {
					primaryIsA := rA == PALIRolePrimary
					if (d > 0) == primaryIsA {
						direction = "primary_higher"
					} else {
						direction = "standby_higher"
					}
				} else {
					// node_pair_asymmetry: use neutral host-name-based labels.
					if d > 0 {
						direction = hA + "_higher"
					} else {
						direction = hB + "_higher"
					}
				}

				fields := map[string]interface{}{
					"metric":    name,
					"hosts":     []string{hA, hB},
					"cohens_d":  roundFloat(d, 4),
					"direction": direction,
					"mean_by_host": map[string]float64{
						hA: roundFloat(meanA, 6),
						hB: roundFloat(meanB, 6),
					},
					"sd_by_host": map[string]float64{
						hA: roundFloat(sdA, 6),
						hB: roundFloat(sdB, 6),
					},
				}
				if rolePairKnown {
					// Carry the inferred role mapping so the analyzer
					// can narrate without re-deriving. Role is on a
					// per-host basis (under `roles_by_host`), not as a
					// top-level `role` field — keeps the bias-control
					// posture from Tier A.6.
					fields["roles_by_host"] = map[string]string{
						hA: string(rA),
						hB: string(rB),
					}
				}
				idKey := name + "|" + hA + "|" + hB
				c.add(newFinding(pairKind, hA+"|"+hB, idKey, fields))
				pairHadEmission = true
				emitted++
			}
			if !pairHadEmission {
				skipped++
			}
		}
	}
	return emitted, skipped
}
