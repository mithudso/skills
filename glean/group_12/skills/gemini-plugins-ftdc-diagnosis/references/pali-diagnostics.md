# Disaggregated Storage (PALI) diagnostics

PALI captures emit a richer set of metric paths under `serverStatus.metrics.disagg.*`. The parser's Tier-5 domain stage covers the principal shapes; your role is to narrate them with PALI context. For the narration rules on the composed PALI kinds (`seal_event`, `role_asymmetry`), see [narration.md](narration.md).

## Architecture context
- **Primary** writes oplog entries to PALI's distributed log via `LogServer`.
- **Standby** reads oplog from the same log and applies via `OplogApplicationCoordinator`.
- **PageServer** holds pages; both nodes read from it.

## Key metric groups
| Group | Path prefix | Interpretation |
|---|---|---|
| Log server | `disagg.logServer.*` | Write side: `appendLog*`, `discontinuousCheckpointInstallCount` |
| PALI log (phylog) | `disagg.phylog.*` | LSN lag: `committedToMaterializedLagMillis`, `generatedToCommittedLagMillis` |
| Page server | `disagg.pageServer.*` | Read latency, cache hit ratio |
| getPage | `disagg.getPageRequest.*` | Queue depth, retry counts |

## Cross-node comparison (replica set)

These are **hypotheses to verify against the capture**, not narration defaults. Primary and standby exhibit different bottleneck shapes; surface the structural evidence and let the data decide.

1. Always inspect both nodes — the parser emits per-host Findings either way.
2. Primary bottlenecks **may** show as page-server read latency rising under combined user + eviction read load. Verify: `mover` on `disagg.pageServer.*` latency + correlation with eviction-read counters.
3. Standby bottlenecks **may** show as oplog application speed limited by WT stable-table lookups. Verify: `materialization_lag` + `apply_saturation` on the standby specifically.
4. Both nodes can be served by the same page server; if so, primary read rate and standby read latency may correlate. Verify: `correlation` between them, lag grid included.
5. Standby reads/s may exceed primary reads/s if the standby reads base pages plus delta chains for key checks. **Verify before assuming normal:** standby:primary ratio via `cross_host_drift` + the workload signature.
6. A standby:primary read ratio > 2× **may** indicate excessive standby reads (e.g. for key existence checks). Verify via `cross_host_drift`'s ratio field; check which `disagg.*` metrics dominate before naming a cause.

## Diagnosing PALI patterns the parser doesn't already match

If the parser fires `diagnosis_candidate("primary_pali_log_pressure")`: narrate per [narration.md](narration.md).

If NOT — look for these signals manually:
- `kind:"materialization_lag"` (Tier-5) — PALI log lag p95 > 1s.
- `kind:"checkpoint_install_anomaly"` (Tier-5) — `discontinuousCheckpointInstallCount` increased.
- `kind:"io_routing"` (Tier-5) — local-to-non-local read ratio shifted toward non-local.

## Standby seal / topology transitions

When the parser emits `kind:"seal_event"`, narrate it directly per [narration.md](narration.md). If it emits only the lower-confidence signals below (and not a composed `seal_event`), you may **hypothesize** a seal window from their co-occurrence — surface it as a candidate window, not an assertion:
- `term_change` in proximity to changes in `replication_buffer_pressure`
- `local_spike` on one of `disagg.logServer.appendLogTotal`, `disagg.logServer.appendLogSuccess`, `disagg.logServer.discontinuousCheckpointInstallCount`
- Optionally a `gap` on the standby host's `disagg.phylog.*` rate stream

Phrase narration as "the capture shows a likely seal window at [timestamp range]; verify via [recipe]." Do not state the window definitively.
