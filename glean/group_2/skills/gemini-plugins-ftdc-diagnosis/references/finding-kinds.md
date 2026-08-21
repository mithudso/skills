# Finding kinds

The single source of truth for every `kind` the parser's `--auto` mode emits. The parser is the authority; this catalog is for narration. When you narrate a Finding, cite it by `id` (stable per run) so the user can grep the raw stream — never by metric name alone.

For *how* to narrate the composed/diagnostic kinds (`diagnosis_candidate`, `hypothesis_test`, `seal_event`, `role_asymmetry`), see [narration.md](narration.md).

## Read order

`--auto` output is sorted by kind then a per-kind stable key. Read it top-to-bottom:

1. `run_metadata` — every algorithmic choice. Spot capture-window misalignments (e.g. baseline window contains the regression).
2. `schema_fingerprint`, `metric_context` — what kind of capture this is. Topology, version, host token.
3. `capture_quality_score` — composite 0–100. Below 70, warn the user before any other diagnosis.
4. `diagnosis_candidate` — parser-emitted shape matches. Narrate these first.
5. `hypothesis_test` — copy-pasteable recipes to confirm/rule-out a candidate. Surface verbatim when `confidence_rank < 0.9`.
6. The rest — `mover`, `correlation`, `monotonic_trend`, `changepoint`, etc. Flesh out the narrative.
7. `run_health` — per-stage emit + skip counts. If any stage shows ≥40% cv_failure, drop confidence in that stage's Findings.

## Tier 0 — core scoring + ranking
- **`metric_score_table`** — ALL ~6500 metrics ranked by Welch t-stat + saturation deviation, nothing pre-filtered. Opt-in (`--auto-emit-score-table`); ~190K tokens.
- **`mover`** — top-K up + top-K down by Stage-1 score. Carries `welch_t`, `welch_df`, `p_value`, `cohens_d`, plus baseline/stress mean, stddev, n.
- **`schema_fingerprint`** — sha256 over `(name, inferred_kind)` tuples. Detects gauge↔counter reclassification or metric-set drift across runs/hosts.

## Tier 1 — context
- **`run_metadata`** — every algorithmic choice (baseline/stress windows, thresholds, lag grid, top-K, recognizer table). Read first.
- **`run_health`** — end-of-run counts by kind, CV-failure %, per-stage emit/skip.
- **`metric_context`** — version, storage engine, tcmalloc variant, topology.
- **`metric_glossary`** — one-line catalog descriptions for cited metrics.

## Tier 2 — structural + correlation
- **`capacity_pressure`** — metric sustained near a structural ceiling (e.g. `cache.bytes >= 0.95 * max_bytes`).
- **`near_ceiling`** — point-in-time ceiling proximity.
- **`consistency_warning`** — cross-Finding contradictions on the same metric.
- **`synthetic_event`** — multi-metric co-occurrence cluster (OVERLAY: read FIRST when present).
- **`evidence_convergence`** — meta-OVERLAY on `synthetic_event`. Lists how many *distinct* Finding kinds co-occur in the cluster; a high distinct-kind count is the parser's structural way of saying "many independent detectors agree."
- **`correlation`** — lagged Pearson against a follower vector (writes p99, reads p99, ticket queues, connections). Multi-follower; lag grid {0,1,5,10,30,60,120,300}s.

## Tier 3 — temporal + spectral
- **`monotonic_trend`** — Spearman rank corr + OLS slope. Carries `spearman_p`, `p_value`. Catches slow ramps (tcmalloc creep, connection leaks).
- **`trend_inconclusive`** — short capture, refused to claim a trend.
- **`changepoint`** — ED-PELT breakpoint. Carries `confidence_rank` (NOT `confidence`) and `cohens_d`. Includes `cross_validated` (CUSUM agreement) — when false, confidence is halved. Breakpoints accurate to ±1 sample via local refinement.
- **`changepoint_oversegmented`** — metric is structurally non-stationary; single summary instead of many breakpoints.
- **`variance_shift`** — CUSUM on rolling-median deviation `|x − median(x)|`. Detects sustained variance shifts even when the mean stays put. The v1 substitute for E-divisive (which fabricates spurious breakpoints on monotonic drift).
- **`spike`** — global MAD robust-z on full-resolution series. Threshold 5σ-equivalent. Catches momentary spikes ED-PELT and variance-shift miss.
- **`local_spike`** — Hampel rolling-window outlier. Complementary to `spike`. Requires ≥110 samples (50+1+50); shorter series silently skipped.
- **`periodic`** — Hann-windowed + AR-prewhitened periodogram, Schuster test. SNR threshold gated by Bonferroni-corrected α.
- **`suspicious_smoothness`** — degenerate variance / constant-zero / impossibly steady rate.

## Tier 4 — observability
- **`utilization`** — USE/RED-style ratios on saturable resources.
- **`alert_threshold`** — Atlas alert classes (election, restart, etc.).
- **`clock_skew`** — sample-interval drift suggesting clock instability.
- **`missing_signal`** — expected metric absent from the capture. `category` field (NOT `class`) buckets by source (core/repl/wt/sharding/host_os).

## Tier 5 — domain stages
- **`term_change`** — replica-set term change (election).
- **`leader_transition`** — primary↔secondary role change.
- **`replication_buffer_pressure`** — secondary oplog buffer near `maxCount`.
- **`oplog_window_low`** — server-side metric absent in current MongoDB; detector dormant.
- **`checkpoint_install_anomaly`** — PALI `discontinuousCheckpointInstallCount` increase.
- **`materialization_lag`** — PALI `disagg.phylog.committedToMaterializedLagMillis` p95 > 1s.
- **`io_routing`** — PALI local-vs-non-local read ratio shift.
- **`chunk_migration_event`** — sharded migrations started.
- **`balancer_round_stall`** — failed balancer rounds.
- **`host_pathology`** — OS-level signal (numa, fragmentation, etc.).
- **`host_metric_set_asymmetry`** — (multi-host) host exposes metrics not in the cluster intersection. Replaces `cross_host_role_inferred`; carries an `evidence` array, no `role` field.
- **`seal_event`** / **`seal_event_candidate`** — composed from co-occurring `term_change` + `local_spike` on `disagg.logServer.*` + (gap on `disagg.phylog.*` OR `discontinuousCheckpointInstallCount` delta) inside a 60s window. Carries `seal_window` (from/to), `supporting_finding_ids`, `confidence_rank` (0.85 for 4/4 conditions, 0.60 for 3/4; `seal_event_candidate` for ≤2/4), and `contradicting_finding_ids` if a `counter_inversion` in the prior 30s suggests a process-restart artifact.
- **`role_asymmetry`** — per-metric divergence between two PALI nodes where roles were inferred. Carries `metric`, `hosts`, `cohens_d` (≥0.8 to emit), `direction`, `mean_by_host`, `sd_by_host`, `roles_by_host` (NOT a top-level `role` field).
- **`node_pair_asymmetry`** — fallback when PALI role inference is ambiguous. Same fields as `role_asymmetry` minus `roles_by_host`.

## Tier 6/7 — reasoning + workflow
- **`workload_spec`** — inferred workload signature.
- **`precursor_candidates`** — Findings that lead a target by t seconds.
- **`memory_triplet`** — joint shape on `mem.resident`, `tcmalloc.generic.current_allocated_bytes`, `tcmalloc.pageheap_free_bytes` — distinguishes fragmentation vs leak.
- **`recurrence_fingerprint`** — 8-hex shape hash; persisted in the Tier-D baseline.

## Tier B/C/D/G/R additions (May 2026)
- **`counter_wrap`** — uint32 wraparound suspected; do NOT reset accumulator.
- **`counter_inversion`** — real counter decrease (process restart / shard movement); reset accumulator.
- **`glossary_gap`** — top-50 cited-but-uncatalogued metrics.
- **`regression`** — current capture deviates from rolling baseline by ≥3σ. **Advisory only** — never gates other detectors. Fresh baseline (`baseline_n < 5`) → expect false positives.
- **`baseline_diagnostic`** — inspectable per-host EWMA state.
- **`apply_saturation`** — secondary repl apply ms-per-batch sustained near saturation.
- **`election_storm`** — ≥3 `term_change` Findings within 5 minutes.
- **`checkpoint_stall`** — WT checkpoint > 80% of syncdelay.
- **`eviction_divergence`** — WT cache read-into / written-from > 2.0.
- **`parameter_change`** — Type-2 PeriodicMetadata diverges mid-capture (setParameter, profiler toggle, featureFlag flip).
- **`capture_quality_score`** — 0-100 composite over `missing_signal` count, `gap` count, `clock_skew`, schema drift, sample-rate stability. Emitted once per host. Below 70: warn before any other narration.
- **`chunk_decode_error`** — declared in the JSON schema enum but **not currently emitted** (worker-panic recovery only writes to stderr today). Reserved.
- **`diagnosis_candidate`** / **`hypothesis_test`** — Tier R; see [narration.md](narration.md).

## Misc
- **`gap`** — consecutive samples > 5s apart (process stall, FTDC pause).
- **`histogram_percentile`** / **`histogram_invalid`** — derived p50/p95/p99 + monotone-violation flag. Recognizers are STRUCTURAL (detect `<prefix>.histogram.<N>.{count,micros}` and `<prefix>.[lo, hi).count` patterns, not hardcoded names); new families auto-detected.
- **`cross_host_drift`** — (multi-host) pairwise window-aggregate ratio between any two hosts exceeds threshold (default 5×, `--auto-drift-ratio`).
- **`host_only_metrics`** — (multi-host) metrics present on a single host (capability fingerprint).
- **`truncated`** — output exceeded budget; list of elided kinds.
- **`skipped`** — stage skipped or detector refused to emit.

## Parser-side pattern catalog (Tier R.1)

Shapes the parser knows. When the capture matches one, the parser emits `diagnosis_candidate`; you narrate. When NO `diagnosis_candidate` fires, use the fallback recipes in [narration.md](narration.md).

| pattern_name | Domain | Required kinds | Discriminators |
|---|---|---|---|
| `primary_pali_log_pressure` | PALI | `materialization_lag` + `memory_triplet` | `counter_wrap`, `regression`, `io_routing` |
| `primary_step_down` | replication | `term_change` + `oplog_window_low` | `alert_threshold`, `regression` |
| `secondary_apply_lag` | replication | `apply_saturation` + `replication_buffer_pressure` | `checkpoint_stall`, `regression` |
| `wt_eviction_pressure` | generic | `eviction_divergence` + `capacity_pressure` | `counter_wrap`, `regression`, `synthetic_event` |
| `host_resource_starvation` | host_os | `utilization` + `memory_triplet` | `regression`, `missing_signal` |
| `noisy_capture_low_quality` | generic | `capture_quality_score` | (none — capture-level concern) |
| `sharded_balancer_thrash` | sharded | `chunk_migration_event` + `balancer_round_stall` | `regression` |
