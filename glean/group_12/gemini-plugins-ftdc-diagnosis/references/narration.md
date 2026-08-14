# Narration rules

How to turn the parser's Findings into a diagnosis without overstepping what the data supports. The governing idea: the parser already did the pattern matching and deliberately withholds verdicts; your job is to narrate its shapes faithfully, not to re-derive or escalate them.

## Contents
- [Core principles](#core-principles)
- [Narrating `diagnosis_candidate`](#narrating-diagnosis_candidate)
- [Narrating `hypothesis_test`](#narrating-hypothesis_test)
- [Narrating `seal_event`](#narrating-seal_event)
- [Narrating `role_asymmetry`](#narrating-role_asymmetry)
- [Statistical contract](#statistical-contract)
- [Bias-control](#bias-control)
- [Contract notes (field renames)](#contract-notes-field-renames)
- [Fallback diagnosis recipes](#fallback-diagnosis-recipes)
- [Analytical thinking guidelines](#analytical-thinking-guidelines)

## Core principles

1. **Read `diagnosis_candidate` first.** The parser already pattern-matched; don't re-derive.
2. **Don't invent verdicts the parser refused to make.** `confidence_rank < 0.5` → "weak signal." Don't escalate.
3. **Quote `hypothesis_test.recipe` verbatim** when the user needs to act on a candidate.
4. **`regression` Findings are advisory** — never gate other diagnoses on them.
5. **Capture quality first.** Below 70, warn the user before any other narration.
6. **The bias-vocab ban applies to you too.** Don't introduce `role`/`severity`/`urgency` in prose the parser doesn't.
7. **Cite Finding IDs, not metric names.** Findings have stable IDs; the user can grep the raw output for them.

> **Tier-R contract:** the parser emits `diagnosis_candidate` and `hypothesis_test` Findings that name the pattern shape directly. On those kinds your job is to **narrate, not re-derive**. The fallback recipes below are only for shapes the parser's catalog doesn't cover.

## Narrating `diagnosis_candidate`

```json
{"kind":"diagnosis_candidate","pattern_name":"primary_pali_log_pressure",
 "supporting_finding_ids":["f-aaa","f-bbb"],"contradicting_finding_ids":[],
 "discriminating_kinds_needed":["counter_wrap","regression"],
 "confidence_rank":0.72,"rationale":"PALI log pressure shape: ..."}
```

1. **Lead with `pattern_name`** translated to a one-liner (`primary_pali_log_pressure` → "primary PALI log pressure suspected").
2. **Cite `supporting_finding_ids` by id**, not by metric — the collector keys those ids.
3. **Treat `confidence_rank` as a ranking score**, not a probability: ≥0.9 → "strong match"; 0.5–0.9 → "plausible"; <0.5 → "weak signal, mention but defer."
4. **If `contradicting_finding_ids` is non-empty**, mention the contradiction — plausible but not certain.
5. **If `discriminating_kinds_needed` is non-empty**, find the paired `hypothesis_test` Findings and surface their `recipe` verbatim.

## Narrating `hypothesis_test`

```json
{"kind":"hypothesis_test",
 "parent_finding_id":"f-ab12cd34",
 "parent_pattern_name":"primary_pali_log_pressure",
 "discriminating_kind":"counter_wrap",
 "recipe":"ftdc_parser /path/to/next/capture/ --auto --filter 'disagg.phylog.*'",
 "expected_outcome_if_confirmed":"counter_wrap on phylog rate, or materialization_lag at higher confidence_rank",
 "expected_outcome_if_ruled_out":"no counter_wrap; phylog counters flat or slowly increasing"}
```

1. **Quote the `recipe` verbatim** — it's pre-engineered to be copy-pasteable.
2. **Quote both expected outcomes** (confirmation AND rule-out) so the user knows what to look for.
3. **Pair each hypothesis_test with its parent** via `parent_finding_id` (the candidate's `id`); `parent_pattern_name` is for readability.

## Narrating `seal_event`

```json
{"kind":"seal_event","at":"2026-03-10T05:50:00.000Z",
 "seal_window":{"from":"2026-03-10T05:49:30.000Z","to":"2026-03-10T05:50:30.000Z"},
 "supporting_finding_ids":["f-aaa","f-bbb","f-ccc"],"confidence_rank":0.85,"conditions_met":4,
 "evidence_summary":{"term_change":true,"local_spike_metrics":["serverStatus.metrics.disagg.logServer.appendLogSuccess"],"phylog_gap_count":1,"discontinuous_checkpoint_install_delta":2,"counter_inversion_blockers":0}}
```

1. **Lead with the `seal_window`** — most downstream PALI patterns compose ON it. Quote the time range.
2. **Cite `supporting_finding_ids` by id.**
3. **If kind is `seal_event_candidate`** (≤2/4 conditions), call it "weak / structurally consistent" — never assert a seal occurred.
4. **If `contradicting_finding_ids` is non-empty**, surface that a `counter_inversion` in the prior 30s could be a process restart masquerading as a seal — say the analyzer "cannot rule out restart-shape."
5. **Don't add severity / cause / "this is bug X."** The Finding declares structural shape only.

## Narrating `role_asymmetry`

```json
{"kind":"role_asymmetry","metric":"serverStatus.metrics.disagg.logServer.appendLogSuccess",
 "hosts":["mongod.0","mongod.1"],"cohens_d":1.7,"direction":"primary_higher",
 "mean_by_host":{"mongod.0":1024.0,"mongod.1":3.2},"sd_by_host":{"mongod.0":120.0,"mongod.1":1.1},
 "roles_by_host":{"mongod.0":"primary","mongod.1":"standby"}}
```

1. **Quote `direction` verbatim**; never re-introduce a top-level `role` field when narrating (bias-control, Tier A.6).
2. **`cohens_d ≥ 0.8`** is "large effect" — call it a real divergence; do not hedge.
3. **`roles_by_host`** is the parser's inference (paliRole helper), not a verdict. Cite it but make clear it's derived from FTDC alone.
4. **`node_pair_asymmetry`** = parser declined to infer role. Narrate as "host A vs host B differ but role attribution is ambiguous."
5. **Combine with `seal_event`** when the seal window is known: an asymmetry concentrated in that window is the cluster-A shape; one that persists outside it is something else.

## Statistical contract

- **`mover`** carries `welch_t`, `welch_df` (Welch-Satterthwaite), `p_value` (two-sided Student-t approx), `cohens_d`. Trust `p_value` for significance, `cohens_d` for effect size. Bare `welch_t` is a back-compat shim.
- **Cohen's d interpretation:** `|d| < 0.2` ignorable noise (don't surface unless corroborated); `0.2 ≤ |d| < 0.5` informational (mention only if relevant); `0.5 ≤ |d|` significant (surface as a key observation).

## Bias-control

The parser ships no `role` / `class` / `severity` / `verdict` / `recommendation` payload fields (enforced by `tools/bias_check.sh` against `BIAS_CONTROL.md`). **The same rule applies to narration.** Don't add `role` back when narrating `host_metric_set_asymmetry`; derive it from the evidence list. Findings declare counter/gauge/histogram shape, units when known, and scores — never cause or named-bug citations.

## Contract notes (field renames)

Additive back-compat was broken in these specific places:
- `changepoint.confidence` → `changepoint.confidence_rank` (Tier B.13): signals "ranking score, NOT calibrated probability."
- `missing_signal.class` → `missing_signal.category` (Tier C bias-vocab rename).
- `cross_host_role_inferred` → `host_metric_set_asymmetry` (Tier A.6): new kind has only `evidence`, no `role` field.

## Fallback diagnosis recipes

Use these only when NO `diagnosis_candidate` fires, or one fires with `confidence_rank < 0.3`, or the pattern is absent from the parser catalog.

### Cache sawtooth
Look for: sustained `changepoint` Findings on `wiredTiger.cache.bytes currently in the cache` with alternating directions, OR a `periodic` Finding on the same metric with `period_seconds` ∈ [60, 600].
Confirm: `ftdc_parser /path/ --rates --filter 'wiredTiger\.cache\.(bytes|evict)' --from T1 --to T2`.

### Burstiness vs sustained load
Look for: high `mover.welch_t` (≥10) with `cohens_d` ≥ 0.5 on `opcounters.*` OR `opLatencies.*`; AND `monotonic_trend.spearman_p > 0.05` (NOT monotonic).
Confirms: bursty user-driven activity; not a gradual ramp.

### Regular ~30-second oplog stalls
Look for: regular `gap` Findings at ~30s spacing in oplog-fetch metric series.
Recipe: `ftdc_parser /path/ --rates --filter 'oplog' | grep -E 'duration_seconds":(2[5-9]|3[0-5])'`.
(If the user asks "is this $TICKET?", consult `jira-index.md` separately — it's not loaded by default, by design.)

### Ticket exhaustion
Look for: `near_ceiling` on `queues.execution.{read,write}.out` AND `correlation` with `metrics.queryExecutor.scannedObjects` as leader.
Confirms: query plans run long enough to hold tickets past arrival rate.

### Checkpoint duration > syncdelay
The parser ships `kind:"checkpoint_stall"`. If it didn't fire, manually check:
`ftdc_parser /path/ --stats --filter 'wiredTiger\.transaction\.transaction checkpoint'` and compare `most recent time (msecs)` vs `interval`.

### FTDC gaps + their causes
Gaps > 5s mean the process couldn't write its periodic doc. Common causes (priority-ordered):
1. Stop-the-world GC / fork at startup — gap at the very beginning.
2. Heavy checkpoint — gap correlated with `kind:"checkpoint_stall"`.
3. Disk I/O stall — gap correlated with high `systemMetrics.disks.io_in_progress`.
4. OOM killer victims — gap followed by `missing_signal` on previously-present metrics.
5. Process restart — gap followed by `counter_inversion`.

### Schema drift mid-capture
The parser emits `kind:"parameter_change"` when Type-2 PeriodicMetadata diverges. Surface verbatim; usually `setParameter`, profiler toggle, or featureFlag flip mid-capture.

## Analytical thinking guidelines

1. **Verify mechanism applicability to topology.** Don't suggest "ticket exhaustion" on a topology that doesn't use the legacy ticket model. Check `metric_context.mongo_version` first.
2. **Distinguish interval vs duration.** `checkpoint interval` (60s default) ≠ `checkpoint duration` (variable); the parser surfaces both separately.
3. **Check current state before proposing changes.** If `connections.current` is 100 of 1000, "increase max connections" is unhelpful.
4. **Check existing metrics before adding new ones.** The parser may already emit what you'd compute by hand.
5. **Quantify expected impact.** "Increase X by Y%" without quantification is a guess.
6. **Don't hedge — investigate.** If a `confidence_rank` is 0.7, surface what would push it to 0.9 (via `discriminating_kinds_needed`).
