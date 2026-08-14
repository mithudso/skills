# ftdc-diagnosis

**Category:** Databases, Data Engineering & Analytics
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/tooling/ftdc-analysis/skills/ftdc-diagnosis

## Description
Parse and interpret MongoDB FTDC (diagnostic.data / metrics.* files) end to end — ships the high-performance `ftdc_parser` Go binary (smart summaries, filtering, rates, per-metric stats with percentiles, metadata, issue detection, and an `--auto` mode emitting structured Findings NDJSON) plus the domain knowledge to turn those Findings into a diagnosis. Use this skill whenever the user has FTDC/diagnostic data and wants to summarize it, extract or graph a metric, detect known performance bottleneck patterns (cache/eviction/checkpoint/ticket/replication/disagg-PALI), reconstruct what a mongod or WiredTiger was doing in a window, or compare captures — even if they just say "analyze this FTDC", "what was the cache doing at 14:12", "why did this node stall", or point at a `metrics.YYYY-...` file.

---

# FTDC Analyzer

End-to-end MongoDB FTDC (Full Time Diagnostic Data Capture) analysis: a Go parser that extracts metrics and emits structured Findings, plus the domain knowledge to interpret them. The parser does the pattern matching and deliberately withholds verdicts (bias-control); narration turns its shapes into a diagnosis without overstepping the evidence.

This SKILL.md is the workflow. Detail lives in `references/`, read on demand:

| File | When to read it |
|---|---|
| [references/parser-cli.md](references/parser-cli.md) | Exact CLI flags, output modes, filtering, `--jobs`, t2-bundle + fixture workflows, performance, building the binary. |
| [references/finding-kinds.md](references/finding-kinds.md) | What each `--auto` Finding `kind` means + the parser-side pattern catalog. The single source of truth. |
| [references/narration.md](references/narration.md) | How to narrate `diagnosis_candidate` / `hypothesis_test` / `seal_event` / `role_asymmetry`, the statistical + bias-control contract, and fallback recipes when no candidate fires. |
| [references/pali-diagnostics.md](references/pali-diagnostics.md) | Disaggregated-storage (PALI) metric groups and cross-node diagnosis. |

The binary lives at `ftdc_parser` in this skill directory. If it's missing or stale, build it (see parser-cli.md): `cd src && go build -o ../ftdc_parser .`. Examples below write `ftdc_parser`; use `./ftdc_parser` from the skill dir.

## Recommended workflow

### Pass 0 — confirm architecture
```bash
ftdc_parser /path/to/capture/ --metadata
```
Read `version`, `process` (mongod vs mongos), `hostInfo.system.numCores`, `storage.engine`, `repl.setName`, `sharding.clusterRole`. Without this you'll diagnose against the wrong topology.

### Pass 1 — orient + analyzer feed
```bash
# Quick human-readable summary: KPIs, issue detection, top active metrics
ftdc_parser /path/to/capture/

# Structured Findings (single host or multi-host cluster dir), with a shareable t2 bundle
ftdc_parser /path/to/capture/ --auto \
  --auto-emit-t2 ~/claude-dump/ftdc-contexts/<input_sha256>/
```

`--auto` emits one NDJSON `Finding` per line. **Read it top-to-bottom in the order documented in [finding-kinds.md](references/finding-kinds.md)**: `run_metadata` → `schema_fingerprint`/`metric_context` → `capture_quality_score` (warn if < 70) → `diagnosis_candidate` (narrate first) → `hypothesis_test` → the rest → `run_health`.

`--auto-emit-t2 <dir>` additionally writes a bundle directory + `<dir>.tgz` (opens in t2 with all hosts on a shared time axis, `analysis.md` narrative, `findings.ndjson`). On Mac, open with `bash open-in-t2.sh`. See parser-cli.md for the bundle layout and the re-seal workflow (regenerating wipes hand-edited narrative).

### Pass 1.5 — system triage when `--auto` is inconclusive
```bash
ftdc_parser /path/ --stats --nonzero --filter "systemMetrics\.(disks|cpu|mem)"
```
Rules out I/O / CPU / kernel-memory bottlenecks before blaming the engine.

### Pass 2 — drill into a finding
```bash
ftdc_parser /path/ --auto --auto-around-event 2026-03-10T05:50:00Z --auto-window 30m
ftdc_parser /path/ --rates --filter "opcounters\.(insert|update)" --from T1 --to T2
ftdc_parser /path/ --json  --filter "queues\.execution\.write"   --from T1 --to T2
ftdc_parser /path/ --stats --nonzero    # full per-metric stats, one NDJSON line each
```

### Pass 3 — cross-capture comparison (when a baseline is available)
```bash
ftdc_parser /path/ --auto --auto-baseline-file ~/.local/share/ftdc-parser/baseline.jsonl
```
Emits `regression` Findings against a rolling EWMA. **Advisory only** — never gates other detectors. A fresh baseline (`baseline_n < 5`) produces false positives.

## Key narration principles

The full rules are in [narration.md](references/narration.md); the essentials:

1. **Read `diagnosis_candidate` first** and narrate it — the parser already pattern-matched; don't re-derive.
2. **Don't invent verdicts the parser refused to make.** `confidence_rank` is a ranking score, not a probability: ≥0.9 strong, 0.5–0.9 plausible, <0.5 weak.
3. **Quote `hypothesis_test.recipe` verbatim** when the user needs to confirm/rule-out a candidate.
4. **Capture quality first** — below 70, warn before any other narration.
5. **`regression` is advisory** — never gate other diagnoses on it.
6. **Bias-vocab ban applies to you too.** Don't introduce `role`/`severity`/`urgency` the parser doesn't. Cite Finding IDs, not metric names.

When NO `diagnosis_candidate` fires and the user describes a specific symptom, fall back to the recipes in narration.md (cache sawtooth, ticket exhaustion, oplog stalls, checkpoint-over-syncdelay, FTDC-gap causes, schema drift).

## Quick reference: key metric paths

```text
serverStatus.connections.{current,available,totalCreated,rejected}
serverStatus.opcounters.{insert,query,update,delete,command}
serverStatus.opLatencies.{reads,writes,commands,transactions}.{ops,latency}
serverStatus.queues.execution.{read,write}.{out,available}
serverStatus.wiredTiger.cache.{bytes currently in the cache, maximum bytes configured, tracked dirty bytes in the cache}
serverStatus.wiredTiger.transaction.transaction checkpoint {most recent time (msecs), currently running, max time (msecs)}
serverStatus.mem.{resident, virtual}
serverStatus.tcmalloc.{generic.current_allocated_bytes, tcmalloc.metadata_bytes, tcmalloc.pageheap_*}
serverStatus.repl.electionId
serverStatus.metrics.repl.{buffer,apply,network}.*
serverStatus.metrics.disagg.{logServer,phylog,pageServer}.*   # PALI — see references/pali-diagnostics.md
serverStatus.shardingStatistics.*
systemMetrics.{cpu,memory,disks,vmstat}.*
```

For full one-line semantics the parser ships `src/metric_docs.go` and emits `kind:"metric_glossary"` for any cited metric.

## Scope and limitations

- Diagnosis-grade narration assumes `--auto` output. For raw `--stats`/`--json` interpretation, use the fallback recipes only.
- Covers mongod/mongos process behavior — NOT Atlas control-plane issues, NOT query-shape regressions (that's the plan-stability skill).
- `jira-index.md` maps structural shapes to known JIRAs. It is intentionally NOT auto-loaded (loading it would prime verdict bias). Consult it yourself only after the user asks "is this $TICKET?", and cite the JIRA as a candidate matched on shape, never as the verdict.

## Related resources

- `ftdc-metric-query` skill — `alexandria` + PromQL exploration of the same FTDC data.
- `sys-perf-ftdc-analyzer` skill — cross-workload extraction from Evergreen patches.
- `mongod-incident-investigation` skill — full guided incident root-cause workflow built on this tool layer.
- T2 statistics visualization — use `--auto-emit-t2` for a shareable bundle.
- `BIAS_CONTROL.md` + `tools/bias_check.sh` — the bias-control invariants the parser (and narrator) must honor.