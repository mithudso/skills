// Tier R.1: Finding-shape pattern catalog for diagnosis_candidate emission.
//
// Each pattern declares:
//   - requiredKinds: Finding kinds whose presence INCREASES confidence
//   - excludedKinds: Finding kinds whose presence DECREASES confidence
//     (contradiction signals)
//   - discriminatingKindsNeeded: additional Finding kinds that, if observed,
//     would let downstream confirm or rule out this pattern. Drives the
//     companion hypothesis_test emission (R.2).
//
// Bias-control constraints (Tier C.2 lint-enforced):
//   - pattern_name strings reference Finding KINDS and Tier-5 DOMAIN STAGES
//     (PALI, replication, sharded, host_os) — never literal metric paths.
//   - rationale is generated mechanically (count + presence) — never freeform.
//
// New patterns: add to diagnosisPatterns. Run `./tools/bias_check.sh` after.

package main

type diagnosisPattern struct {
	name                      string
	requiredKinds             []string // kinds that must ALL be present
	supportingKinds           []string // optional: kinds that boost confidence
	excludedKinds             []string // kinds that contradict
	discriminatingKindsNeeded []string // emit hypothesis_test for these
	domainStageGate           string   // "" = generic; "pali"/"replication"/"sharded" = restricted
	rationaleTemplate         string   // sprintf'd with: supporting-count, excluded-count
}

var diagnosisPatterns = []diagnosisPattern{
	{
		name:                      "primary_pali_log_pressure",
		requiredKinds:             []string{"materialization_lag", "memory_triplet"},
		supportingKinds:           []string{"checkpoint_install_anomaly", "monotonic_trend"},
		excludedKinds:             []string{KindTermChange}, // a term change suggests election storm, not log pressure
		discriminatingKindsNeeded: []string{"counter_wrap", "regression", "io_routing"},
		domainStageGate:           "pali",
		rationaleTemplate:         "PALI log pressure shape: materialization_lag + memory_triplet co-occur on primary; %d supporting Finding(s) observed; %d excluded.",
	},
	{
		name:                      "primary_step_down",
		requiredKinds:             []string{KindTermChange, "oplog_window_low"},
		supportingKinds:           []string{"election_storm", "replication_buffer_pressure"},
		excludedKinds:             []string{},
		discriminatingKindsNeeded: []string{"alert_threshold", "regression"},
		domainStageGate:           "replication",
		rationaleTemplate:         "Step-down shape: term_change + oplog_window_low; %d supporting Finding(s); %d excluded.",
	},
	{
		name:                      "secondary_apply_lag",
		requiredKinds:             []string{"apply_saturation", "replication_buffer_pressure"},
		supportingKinds:           []string{"monotonic_trend"},
		excludedKinds:             []string{KindTermChange},
		discriminatingKindsNeeded: []string{"checkpoint_stall", "regression"},
		domainStageGate:           "replication",
		rationaleTemplate:         "Secondary apply-saturation shape: apply_saturation + replication_buffer_pressure; %d supporting Finding(s); %d excluded.",
	},
	{
		name:                      "wt_eviction_pressure",
		requiredKinds:             []string{"eviction_divergence", "capacity_pressure"},
		supportingKinds:           []string{"checkpoint_stall", "near_ceiling"},
		excludedKinds:             []string{},
		discriminatingKindsNeeded: []string{"counter_wrap", "regression", "synthetic_event"},
		domainStageGate:           "",
		rationaleTemplate:         "WT eviction-pressure shape: eviction_divergence + capacity_pressure; %d supporting Finding(s); %d excluded.",
	},
	{
		name:                      "host_resource_starvation",
		requiredKinds:             []string{"utilization", "memory_triplet"},
		supportingKinds:           []string{"clock_skew", KindGap, "monotonic_trend"},
		excludedKinds:             []string{},
		discriminatingKindsNeeded: []string{"regression", "missing_signal"},
		domainStageGate:           "host_os",
		rationaleTemplate:         "Host-resource-starvation shape: utilization + memory_triplet; %d supporting Finding(s); %d excluded; clock skew or gaps may corroborate.",
	},
	{
		name:          "noisy_capture_low_quality",
		requiredKinds: []string{"capture_quality_score"},
		// chunk_decode_error is in the schema enum but the parser does
		// not currently emit it (the worker-panic recovery only logs to
		// stderr). Removed from supportingKinds until emission is wired
		// — otherwise the kind would never match and inflate skipped count.
		supportingKinds:           []string{"missing_signal", "suspicious_smoothness", KindGap},
		excludedKinds:             []string{},
		discriminatingKindsNeeded: []string{},
		domainStageGate:           "",
		rationaleTemplate:         "Capture-quality concern: %d supporting Finding(s) consistent with a degraded capture; %d excluded.",
	},
	{
		name:                      "sharded_balancer_thrash",
		requiredKinds:             []string{"chunk_migration_event", "balancer_round_stall"},
		supportingKinds:           []string{"monotonic_trend"},
		excludedKinds:             []string{},
		discriminatingKindsNeeded: []string{"regression"},
		domainStageGate:           "sharded",
		rationaleTemplate:         "Sharded balancer-thrash shape: chunk_migration_event + balancer_round_stall; %d supporting Finding(s); %d excluded.",
	},
	// P1.NEW-PALI-3 (post-P0.NEW-PALI-1/2): compose on the new
	// seal_event + role_asymmetry primitives. Each maps to a structural
	// shape observed in the disagg-storage chaos sweeps. The patterns
	// here are PALI-gated and assume the seal-event detector + role
	// asymmetry detector are wired into the pipeline.
	{
		name:                      "standby_seal_topology_skew",
		requiredKinds:             []string{KindSealEvent, "role_asymmetry"},
		supportingKinds:           []string{"checkpoint_install_anomaly", KindTermChange},
		excludedKinds:             []string{"counter_inversion"},
		discriminatingKindsNeeded: []string{"materialization_lag", "regression"},
		domainStageGate:           "pali",
		rationaleTemplate:         "PALI standby seal/topology skew: seal_event + role_asymmetry co-occur; %d supporting Finding(s); %d excluded; cluster-A shape if log-server next_log_consensus_group asymmetry is the divergent metric.",
	},
	{
		name:                      "standby_read_log_wedge",
		requiredKinds:             []string{KindSealEvent},
		// Supporting: gap on the standby phylog rate stream + missing_signal
		// (catches a standby that stopped advancing entirely post-seal).
		supportingKinds:           []string{KindGap, "missing_signal"},
		// election_storm would suggest the seal is part of a normal
		// election cascade, not a wedge.
		excludedKinds:             []string{"election_storm"},
		discriminatingKindsNeeded: []string{"counter_wrap", "regression"},
		domainStageGate:           "pali",
		rationaleTemplate:         "PALI standby read-log wedge: seal_event followed by stalled standby phylog; %d supporting Finding(s); %d excluded; matches chaos-bug-6 shape (logReader->Finish() hang).",
	},
	{
		name:                      "materialization_apply_split",
		// Distinguishes log-side vs apply-side stall: materialization_lag
		// fires AND role_asymmetry shows lag divergence between primary
		// and standby on the materialization metric.
		requiredKinds:             []string{"materialization_lag", "role_asymmetry"},
		supportingKinds:           []string{"monotonic_trend", "synthetic_event"},
		excludedKinds:             []string{},
		discriminatingKindsNeeded: []string{"counter_wrap", "regression"},
		domainStageGate:           "pali",
		rationaleTemplate:         "Materialization/apply split: materialization_lag + role_asymmetry on PALI lag metric; %d supporting Finding(s); %d excluded.",
	},
	{
		name:                      "pageserver_contention",
		// Both nodes' read rate elevated AND pageServer latency rising;
		// matches the "same page server serves both" shape called out in
		// the analyzer SKILL.md PALI section.
		requiredKinds:             []string{"capacity_pressure", "correlation"},
		supportingKinds:           []string{"monotonic_trend", "near_ceiling"},
		excludedKinds:             []string{"counter_inversion"},
		discriminatingKindsNeeded: []string{"io_routing", "regression"},
		domainStageGate:           "pali",
		rationaleTemplate:         "PALI pageserver contention: capacity_pressure + correlation on pageServer metrics; %d supporting Finding(s); %d excluded.",
	},
}
