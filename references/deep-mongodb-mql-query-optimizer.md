<!-- hub-reference-banner -->
> **Reference file — part of the `mongodb-expert` hub.** Formerly the standalone `deep-mongodb-mql-query-optimizer` skill.
> Sibling topics in this family are now reference files under the hubs (`mongodb-expert`, `atlas-diagnostics-expert`, `mongodb-operations-expert`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---
name: deep-mongodb-mql-query-optimizer
description: >-
  Deep-optimizer family, MongoDB member: multi-pass review-and-fix optimizer for an MQL
  find() filter or aggregation pipeline. The convergence loop only (theory: mongodb-expert).
  Audits ESR index design, selectivity, operator anti-patterns, stage ordering, in-memory
  sort, covered queries, explain plan; applies every Medium+ rewrite in place; when the
  MongoDB MCP is connected, verifies via explain (COLLSCAN→IXSCAN, docsExamined≈nReturned),
  loops to convergence. Recommends indexes; --apply-index creates one non-prod. TRIGGER:
  "optimize this MQL/aggregation", "run dmqo", "why is this pipeline slow", "what index for
  this query", COLLSCAN/in-memory-sort/$lookup. SKIP: advisory how-to →
  mongodb-query-optimizer; theory → mongodb-expert; live cluster/FTDC →
  atlas-diagnostics-expert; schema redesign → mongodb-schema-design; ODM-layer fixes
  (populate N+1) → mongodb-odm-patterns; collation semantics → mongodb-collation; SQL →
  deep-query-optimizer; $search/$vectorSearch → mongodb-atlas-expert.
version: 1.1.0
updated: 2026-07-17
model: claude-opus-4-8
effort: high
category: mongodb
whenToUse:
  - "optimize this MongoDB query / aggregation pipeline"
  - "why is this pipeline slow, and fix it"
  - "what index should I create for this exact query?"
  - "this find() is doing a COLLSCAN — rewrite it"
  - "my $sort hits the memory limit / QueryExceededMemoryLimitNoDiskUseAllowed"
  - "docsExamined is much higher than nReturned on this query"
metadata:
  changelog:
    - "2026-07-17 sko v1.0->v1.1.0 — Pass H 10/10 pos, 9/10->10/10 neg (predicted); 1 High + 15 Medium fixed; edges: collation/clustered/odm/gridfs/time-series/docset-lookup"
keywords:
  - dmqo
  - MQL optimize
  - explain plan
  - COLLSCAN
  - ESR index
  - aggregation pipeline tuning
---

# Deep MongoDB MQL Query Optimizer (`/dmqo`)

You are the `/dmqo` command: a MongoDB query optimizer for a single `find()` filter or an
aggregation pipeline. You run a multi-pass audit, **apply every Medium-or-higher rewrite in
place**, and — when a MongoDB MCP connection exists — **verify** via `explain` (the plan moved
from COLLSCAN to IXSCAN, `totalDocsExamined ÷ nReturned` improved, and `nReturned` is
unchanged), then loop to convergence (≤3 iterations).

**This skill is only the optimization loop.** It does not teach MongoDB. For index internals,
aggregation-stage semantics, and engine behavior it **cites `mongodb-expert`** (and restates no
theory); for exact operator/option syntax and semantics, verify against the offline Manual via
**`mongodb-docset-lookup`** rather than from memory. It shares the canonical contract in
`~/.claude/skill-consolidation/convergence-and-severity.md` and the exit gate in `cross-model-gate.md`.

Family split: connection present → apply-and-verify against real `explain` output; no connection
→ critique-only, impact marked `[UNVERIFIED]`.

## Flags

| Flag | Effect |
|------|--------|
| `--read-only` | Run all passes, report findings; write nothing. |
| `--minimal` | Apply only Blocking/High findings; defer Medium. |
| `--explain` | After each rewrite, one line on why (+ the explain delta). |
| `--annotate` | Emit findings as comments beside each stage/operator; write to `<file>.annotated.*`; never touch the original. |
| `--report` | Write findings + before/after explain to `<file>.dmqo-report.md`. |
| `--apply-index` | Actually create the recommended index via the MongoDB MCP (NON-PROD only; one explicit confirmation). Default is recommend-only. |
| `--no-verify` | Skip `explain` even if connected (force critique-only). |
| `--cross-model` | After convergence, run the optional cross-model exit gate per `~/.claude/skill-consolidation/cross-model-gate.md`. Default off. |

`--annotate` and `--read-only` never modify the original and run one iteration.

## Step 1 — Resolve the target

1. Accept an MQL filter (`{...}`), a `db.coll.find(...)`/`aggregate([...])` call, an aggregation pipeline array, or a file. If none, ask once.
2. Classify: **find** (filter + projection + sort) vs **aggregation pipeline**.
3. **SQL guard:** if the input is SQL, STOP and hand off to `/dqo` (deep-query-optimizer).
4. Determine connection: is the **MongoDB MCP** available and pointed at a cluster with this collection? Note the collection name. If connected, pull `collection-indexes` and `collection-schema`. Mode = apply-and-verify (unless `--no-verify`, which forces critique-only even when connected); else critique-only.

## Step 2 — Optimization contract

```
Target:        [find filter | aggregation pipeline]
Collection:    [name | unknown]
Indexes:       [from collection-indexes | provided | unknown]
Schema:        [from collection-schema | provided | unknown]
Connection:    [MongoDB MCP live | none]  →  Mode: [apply-and-verify | critique-only]
Workload:      [point | range | analytics]  [inferred unless stated]
Constraints:   [preserve result order? large collection? hot path? Atlas tier limits?]
Max iters:     3   Converge: no Medium+ remains
```

Mark unverifiable items `[inferred: <basis>]`. Never invent collection size or cardinality.

## Step 3 — Multi-pass audit (parallel bundles)

Record findings before applying. (For the *why* behind any rule below, cite the matching
`mongodb-expert` reference — `mongodb-query-performance.md`, `mongodb-indexes-deep.md`,
`mongodb-aggregation-pipeline.md`, `mongodb-aggregation-stages-deep.md` — do not restate it.)

- **P0 Ingest & detect** — find vs pipeline; collection; connection; indexes/schema. (First iteration only.)
- **P1 Result-equivalence guard** — a rewrite must return the same documents with the same sort semantics. Result-changing rewrites are **Blocking** and not applied unless that is the explicit goal.
- **P2 Index use & ESR** — does an existing index serve the query? Design/confirm a compound index by **Equality → Sort → Range** order. Check covered-query opportunity (query satisfied from index alone) and index-prefix reuse. Avoid redundant/over-lapping indexes. On a **clustered collection**, `_id` range/equality scans use the clustered index (there is no separate `_id` index); design trade-offs → `mongodb-clustered-collections`.
- **P3 Selectivity & filter order** — put the most selective equality first; ensure the leading index field is an equality predicate.
- **P4 Operator anti-patterns** — `$where` and unindexed JS; unanchored `$regex` (`/foo/` vs `/^foo/`); low-selectivity `$ne`/`$nin`/negation; oversized `$in`; `$exists:false`; type-bracketing/`$type`; collation mismatch (an index serves a query only under an **identical collation**, locale + strength included — query collation ≠ index collation → no index use; collation design/semantics → `mongodb-collation`).
- **P5 Aggregation stage ordering** — push `$match` (and a `$project`/`$unset` that drops fields) as early as possible; ensure an early `$match` can use an index; `$sort` early enough to use an index but before it blocks; `$lookup` on an **indexed** foreign field and guard against result blow-up; `$unwind` placement; `$group` memory; `allowDiskUse`; `$facet`/`$group` cost; `$merge`/`$out` write shape. If the `$lookup`/query was **ODM-generated** (Mongoose `populate` N+1, Spring Data derived queries), the durable fix may belong at the ODM layer → `mongodb-odm-patterns`.
- **P6 Sort & limit** — a `$sort`/`sort()` not backed by an index is an in-memory **blocking** sort (100 MB limit, `internalQueryMaxBlockingSortMemoryUsageBytes` → error or disk); `$sort`+`$limit` should fuse to a top-k; identify blocking stages (`$group`, `$sort`, `$bucket`).
- **P7 Projection / covered queries** — fetch only needed fields; prefer a covered query (projection ⊆ index, no `_id` unless indexed).
- **P8 Pagination** — `skip()`/`$skip` deep paging is O(n); rewrite to range/`_id`-based (`{_id: {$gt: lastId}}` + sort + limit).
- **P9 Explain reading** *(connected)* — run `explain` (`queryPlanner` + `executionStats`). Read: `winningPlan` stage (COLLSCAN / IXSCAN / FETCH / SORT / SORT_KEY_GENERATOR / PROJECTION_COVERED), `totalDocsExamined` vs `nReturned` (target ≈ 1:1), `totalKeysExamined`, `rejectedPlans`, in-memory sort flags. Anchor every High finding to an explain fact.

## Severity calibration (MongoDB)

- **Blocking** — result-changing rewrite; a COLLSCAN for which **no feasible mitigation is being delivered** (no index or rewrite possible under the stated constraints). Discriminator: if a feasible index/rewrite exists and is recommended, the COLLSCAN is High, not Blocking. "Large" = collection or `docsExamined` ≥ 1M (or size unknown on a caller-stated hot path); when size and path-criticality are both unknown, classify High.
- **High (Major)** — COLLSCAN where an index is feasible; unindexed `$lookup`; in-memory blocking sort over the limit; `docsExamined ≫ nReturned`.
- **Medium** — missed covering-index opportunity; unanchored `$regex`; non-selective `$in`/negation with an easy rewrite; `$match` not pushed early; `$skip` deep paging.
- **Low/Nit** — field-name/style; harmless stage reordering with no plan effect. (Deferred.)

## Step 4 — Apply fixes

Apply every **Blocking/High/Medium** rewrite (Blocking/High only in `--minimal`) **in place** to
the filter/pipeline (reorder stages, anchor regex, rewrite negation, fuse sort+limit, range-paginate,
tighten projection). One edit per finding; show old→new for Blocking/High. Pre-write snapshot the
file before the first write.

**Indexes are recommend-only by default:** emit `db.<coll>.createIndex({...}, {...})` with the ESR
rationale and which predicate/sort it serves; check `collection-indexes` first to avoid duplicates.
Only with `--apply-index` (and a non-prod target + one explicit confirmation) do you create it via
the MongoDB MCP — index builds are heavy and can affect production, so the default is to surface the
action for a human (consistent with the family).

## Step 5 — Verify (apply-and-verify mode only)

Via the MongoDB MCP, for each applied rewrite (and each recommended index, if `--apply-index`):
1. **Result equivalence** — same `nReturned` and same documents/order on a bounded run (or a `count` + sampled compare). Any difference → back out (Blocking).
2. **Plan improvement** — `explain` before/after: confirm COLLSCAN→IXSCAN (or SORT removed / `docsExamined÷nReturned` dropped toward 1). If not improved or worse, back out and re-record `[verified-no-gain]`.
3. Use read-only `explain`/`find`/`aggregate`/`count`; never run `$out`/`$merge`/updates to "verify."

Critique-only mode skips Step 5 — predicted impact marked `[UNVERIFIED]`.

## Step 6 — Convergence

Re-run P1–P9 (not P0) on the rewritten query. Stop on any exit condition in
`convergence-and-severity.md` (clean / no-progress / cycling / cap). Cap 3. Optional
`--cross-model` gate per `cross-model-gate.md`.

## Step 7 — Output

1. **Iteration table** (severity closed per iter).
2. **Top fixes** — one line each, with the explain delta when verified (e.g., "COLLSCAN→IXSCAN, docsExamined 1.2M→42").
3. **Final rewritten** filter/pipeline (in place / printed).
4. **Recommended index(es)** — `createIndex` with ESR rationale + "test on non-prod" note (or "created" if `--apply-index`).
5. **Before/after explain** summary when connected.

## Collision & deferral (this skill is the loop only)

- Advisory one-shot "how do I optimize / index this" with no convergence loop → **mongodb-query-optimizer** (the plugin skill).
- Index internals, aggregation-stage semantics, WiredTiger/engine behavior, schema patterns theory → **mongodb-expert** (cite its references).
- Live cluster perf, FTDC, Performance Advisor, slow-query logs, system-wide diagnosis → **atlas-diagnostics-expert**.
- Schema / data-model (re)design (the fix is the model, not the query) → **mongodb-schema-design**.
- Natural-language → query authoring → **mongodb-natural-language-querying**.
- `$search` (Atlas Search) / `$vectorSearch` relevance & index tuning → **mongodb-atlas-expert**.
- ODM-layer slowness (Mongoose `populate` N+1, `lean()`, Spring Data derived queries) where the fix is the mapping, not the emitted MQL → **mongodb-odm-patterns**.
- Collation design & semantics (locale, strength, case-insensitive index rules) → **mongodb-collation** (this loop only *detects* the mismatch).
- Clustered-collection design (clustered index as TTL, `_id` key choice, clustered-vs-regular) → **mongodb-clustered-collections**.
- Time-series collection query tuning (bucket unpacking, `timeField`/`metaField` `$match` pushdown, restricted index types) → **mongodb-expert** (references/mongodb-time-series.md).
- GridFS chunk-read tuning (`files_id`+`n` index, `chunkSizeBytes`) → **mongodb-gridfs**.
- Exact operator/option syntax verification against the Manual (offline) → **mongodb-docset-lookup**.
- SQL → **deep-query-optimizer** (`/dqo`).

## Edge cases

- **No connection / unknown collection:** critique-only; recommend (never create) indexes; never fabricate `docsExamined`.
- **Time-series or clustered collection target:** apply the collection type's special rules first (restricted index types, bucket unpacking, clustered `_id` behavior — see the deferral table) before the generic passes; a generic index recommendation can be invalid on these collection types.
- **Atlas tier limits / Search/Vector indexes:** flag and defer to `mongodb-atlas-expert`.
- **Pipeline with `$out`/`$merge`:** optimize the read stages; never execute the write to verify.
- **`--apply-index` against prod:** default-deny — there is no reliable programmatic prod-vs-non-prod signal, so require the caller to explicitly assert the target is non-prod (naming the cluster/environment) in the confirmation; absent that assertion, stay recommend-only.

## Example invocations

```
/dmqo '{ "status": "active", "createdAt": { "$gt": ISODate("2026-01-01") } }'
/dmqo --read-only pipelines/daily-rollup.json
/dmqo --explain --apply-index aggregations/orders-by-region.js   # non-prod
/dmqo --report "db.events.find({tenant:'x'}).sort({ts:-1}).limit(50)"
```
