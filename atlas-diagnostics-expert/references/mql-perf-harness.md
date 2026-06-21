<!-- hub-reference-banner -->
> **Reference file — part of the `atlas-diagnostics-expert` hub.** Formerly the standalone `mql-perf-harness` skill.
> Sibling topics in this family are now reference files under the hubs (`mongodb-expert`, `atlas-diagnostics-expert`, `mongodb-operations-expert`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: mql-perf-harness
description: >-
  Heuristic, index-aware performance scorer plus a 50-query anti-pattern benchmark corpus for
  MongoDB MQL find()/aggregation queries. Detects COLLSCANs, in-memory sorts, ESR violations,
  $where/regex/negation/$in/$size operator anti-patterns, late $match, unindexed $lookup, deep
  $skip, and missed covering indexes, then assigns a weighted cost score (a transparent proxy for
  docsExamined / blocking-stage risk, NOT measured latency) so you can quantify query cost before
  vs after optimization and prove a delta.
  TRIGGER: "score/benchmark this MQL query", "how bad is this aggregation", "before/after delta for
  query optimization", "run the mql perf harness", build an anti-pattern test corpus, measure /dmqo
  impact across many queries.
  SKIP: actually rewriting a query to convergence -> deep-mongodb-mql-query-optimizer (/dmqo); live
  cluster explain/latency profiling -> atlas-diagnostics-expert; index/aggregation theory ->
  mongodb-expert; SQL -> deep-query-optimizer (/dqo).
version: 1.0.1
category: developer
updated: 2026-06-18
metadata:
  changelog:
    - "2026-06-18 sko v1.0.0->1.0.1: iter1 fixed 1 Medium (Pass A dead reference -> created references/scoring-model.md); iter2 blind re-audit clean. Pass M desc 974/1000; Pass K em-dash 0.73/100w; Pass H predicted 10/10 pos, 1/10 neg."
keywords:
  - mql performance
  - query anti-patterns
  - heuristic scorer
  - COLLSCAN detection
  - ESR index
  - in-memory sort
  - aggregation pipeline cost
  - before/after benchmark
  - dmqo benchmark
related_skills:
  - deep-mongodb-mql-query-optimizer
  - mongodb-expert
  - mongodb-query-optimizer
  - atlas-diagnostics-expert
whenToUse:
  - "score or benchmark an MQL find() or aggregation pipeline for performance risk"
  - "produce a before/after cost delta to show an optimization paid off"
  - "generate a corpus of bad MQL queries to test an optimizer"
  - "measure the impact of /dmqo across many queries at once"
whenNotToUse:
  - "rewrite/optimize a single query to convergence (use deep-mongodb-mql-query-optimizer /dmqo)"
  - "measure real latency/explain on a live cluster (use atlas-diagnostics-expert + the MongoDB MCP)"
  - "learn index/aggregation theory (use mongodb-expert)"
  - "score SQL (use deep-query-optimizer /dqo)"
---

# MQL Performance Harness

A reproducible, **index-aware heuristic** for scoring MongoDB MQL queries by performance risk, plus
a bundled **50-query anti-pattern corpus** (each bad query paired with its optimized form). It exists
to answer one question with data: *did optimizing these queries actually help, and by how much?*

The score is a **transparent proxy for executor work** (docsExamined / blocking-stage risk), not a
measured millisecond latency. It is deterministic, needs no live cluster, and rewards exactly the
moves `/dmqo` makes: adding a usable index, ordering it Equality→Sort→Range, anchoring a regex,
inverting a negation, pushing `$match` early, fusing `$sort`+`$limit`, keyset-paginating, and
indexing a `$lookup` foreign key.

## Files

- `scripts/mql-perf-tester.mjs` — the scorer (Node, zero dependencies).
- `corpus/corpus.mjs` — the 50-query corpus; each entry has `bad`, `opt` (optimized query +
  recommended indexes), available `idx`, and anti-pattern tags.
- `references/scoring-model.md` — the full penalty table, detector rules, and honest limitations.

## Run it

```bash
T=~/.claude/skills/mql-perf-harness/scripts/mql-perf-tester.mjs

# baseline: score every bad query against its pre-existing indexes
node "$T" --mode baseline  --out baseline.json

# after: score every optimized query against pre-existing + recommended indexes
node "$T" --mode optimized --out after.json

# delta: per-query + aggregate before/after improvement
node "$T" --compare baseline.json after.json --out delta.json
```

Each scoring run prints a one-line summary and writes a JSON report plus a CSV. `--compare` writes a
per-query delta table. Point at a different corpus with `--corpus <path.mjs>`.

## How a query is scored

`scoreEntry(entry, mode)` picks the query (`bad` or `opt`) and the indexes in scope (baseline =
available; optimized = available ∪ recommended), runs the find or aggregation detectors, and sums
weighted penalties. Higher = worse; an indexed point query floors near the `BASE` cost. The detectors
implement the `/dmqo` pass taxonomy (P2 index/ESR, P3 selectivity, P4 operators, P5 stage ordering,
P6 sort/limit, P7 covered queries, P8 pagination). Full table in `references/scoring-model.md`.

## Limitations (read before quoting a number)

- **It is a heuristic, not a profiler.** The cost units are not milliseconds and not docsExamined;
  they are a weighted anti-pattern proxy. Use it for *relative* before/after comparison, not absolute
  capacity planning.
- **Index-awareness is structural.** It checks index-prefix/ESR alignment, not real cardinality or
  selectivity statistics, so it cannot see data skew.
- **Some optimizations assume schema support** (e.g., a precomputed boolean to replace `$where`/`$size`/
  `$mod`); those are flagged in the corpus `note`.
- **For ground truth, validate on a live cluster.** Run `/dmqo` with the MongoDB MCP connected to
  confirm COLLSCAN→IXSCAN and a real docsExamined drop via `explain()`.
