<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly the standalone `dexie-indexeddb-local-first-reviewer` skill.
> Sibling topics in this family are now reference files under the hubs (`chrome-extension-expert`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: dexie-indexeddb-local-first-reviewer
description: >
  Dexie.js and IndexedDB local-first architecture reviewer. Audits schema/index design,
  migration safety, transaction atomicity, cache invalidation, and Chrome MV3 persistence
  constraints. Covers compound index ordering, bulkAdd vs bulkPut semantics, quota/eviction
  handling, and service-worker lifecycle constraints on IndexedDB state.
  TRIGGER: user is reviewing, debugging, or designing Dexie or IndexedDB usage in a local-first
  or Chrome extension app; asking about schema migrations, index design, transaction boundaries,
  bulk writes, cache invalidation, or MV3 storage persistence.
  SKIP: pure MongoDB or server-side database questions; SQLite or relational DB questions;
  browser storage questions not involving IndexedDB/Dexie.
version: "1.1.0"
updated: "2026-05-29"
category: developer
origin: local
tags:
  - dexie
  - indexeddb
  - local-first
  - chrome-extension
  - mv3
  - offline
  - browser-storage
keywords:
  - dexie
  - IndexedDB
  - local-first
  - schema migration
  - compound index
  - transaction
  - bulkAdd
  - bulkPut
  - cache invalidation
  - service worker
  - MV3
  - quota eviction
  - upgrade
  - IDBKeyRange
whenToUse:
  - Reviewing Dexie schema definitions and index declarations
  - Auditing IndexedDB migration safety and upgrade() correctness
  - Checking transaction boundaries and atomicity in local-first writes
  - Evaluating bulkAdd vs bulkPut usage and BulkError handling
  - Diagnosing cache invalidation drift in derived views or thread summaries
  - Reviewing Chrome MV3 extension persistence: what survives SW suspend
  - Designing compound indexes that match multi-field access patterns
  - Checking quota/eviction assumptions for offline-capable apps
whenNotToUse:
  - Server-side or MongoDB database questions — use mongodb-expert
  - SQLite or relational schema design — use databases-aws-cockroach-indexeddb
  - General browser storage (localStorage, sessionStorage) with no Dexie/IndexedDB angle
related_skills:
  - chrome-extension-expert
  - databases-aws-cockroach-indexeddb
---

# Dexie / IndexedDB Local-first Reviewer

Practical review reference for auditing Dexie and IndexedDB usage in local-first applications.
Covers schema/index correctness, migration safety, transaction boundaries, cache invalidation,
and MV3 persistence constraints.

**Sources:** [Dexie docs](https://dexie.org/docs/), [MDN IndexedDB API](https://developer.mozilla.org/en-US/docs/Web/API/IndexedDB_API), [Chrome MV3 service worker lifecycle](https://developer.chrome.com/docs/extensions/develop/concepts/service-workers/lifecycle).

*Context verified: 2026-05-10. Repo context: `docs/ARCHITECTURE.md`, `packages/store/src/db.ts`.*

## Quick Review Rules

1. **Index only what you query.** If code sorts or filters on a field, that field must be in `stores()`. Over-indexing hurts performance; Dexie explicitly recommends avoiding large indexed values.
2. **Compound index order must match query shape.** Declared as `[a+b]`; query prefix must align or the index is skipped.
3. **Schema changes are migrations.** New indexes and renamed fields belong in versioned `upgrade()` calls, not ad hoc runtime fixups.
4. **Wrap related writes in one transaction.** Control-yielding async gaps can break intended atomicity even inside Dexie transactions.
5. **Prefer batch writes.** Review per-document write loops skeptically — `bulkAdd()`/`bulkPut()` exist for a reason.
6. **Every cache needs an explicit invalidation trigger** near the write paths that can stale it.
7. **MV3 globals are disposable.** Anything that must survive suspend belongs in IndexedDB or `chrome.storage.*`, not memory.

## Review Workflow

| Step | What to check | Key APIs |
|---|---|---|
| 1. Map access patterns | Identify hot reads, sorted/range queries, fan-out joins before judging schema | `stores()`, `where()` |
| 2. Verify index coverage | Each `where()`, `orderBy()`, compound selector → actual schema index exists and order matches | `where()`, `IDBIndex` |
| 3. Inspect write atomicity | Related writes, migration flags, derived-record updates in one transaction or partial-commit risk | `transaction()`, `IDBTransaction` |
| 4. Check migration safety | New indexes, renamed fields, backfills are versioned and idempotent; first-open cannot loop | `upgrade()`, `IDBDatabase` |
| 5. Audit cache invalidation | Thread summaries, unread counts, prompt caches cannot drift after writes | repo `docs/ARCHITECTURE.md` |
| 6. Check failure modes | Quota/eviction assumptions, re-open behavior, MV3 suspend/resume | MDN storage quotas, Chrome lifecycle |

## Method Reference

| Method / surface | Purpose | Review focus | Caveats |
|---|---|---|---|
| `version().stores()` | Declares tables and indexes | Index only queried fields; verify compound index order | Over-indexing hurts performance and stability |
| `version().upgrade()` | Versioned schema/data migration | Idempotence, field backfill safety, old-to-new mapping | Runs inside upgrade transaction; mistakes can strand users on open |
| `transaction()` | Atomic read/write group | Async gaps, nested scope, rollback behavior | IndexedDB auto-commits when control yields unexpectedly |
| `bulkAdd()` / `bulkPut()` | Batch insert/upsert | Per-doc loops, `BulkError` handling, partial-success assumptions | Insert-only vs upsert semantics differ — do not use `bulkAdd()` where idempotence requires `bulkPut()` |
| `where()` / `WhereClause` | Indexed lookup/range query | Selector-to-index alignment, fallback filtering | Misaligned queries silently degrade to slower JS-side filtering |
| `IDBKeyRange`-style range | Efficient lower/upper bounds | Inclusive/exclusive edge correctness | Easy off-by-one errors |
| `Dexie.waitFor()` | Keep transaction alive across async | Whether truly necessary | CPU-expensive in hot paths |
| Repo cache declarations | Fast-path derived reads | Explicit invalidation, write adjacency, bounded growth | TTL-only caches can hide stale state if writes don't bust them |

## Standards

- Use **compound indexes for real multi-field access patterns** instead of broad JS-side filtering on a weak seed.
- Keep **migration logic explicit, versioned, and one-way** — no live shape mutations without an upgrade path.
- Distinguish **insert-only** (`bulkAdd`) and **upsert** (`bulkPut`) semantics deliberately.
- Treat **quota and eviction as variable browser behavior**, not a guaranteed capacity budget.
- In MV3, persist **authoritative state** to IndexedDB/storage; treat in-memory caches as performance helpers only.
- Every cache must have an **explicit invalidation trigger** — TTL alone is insufficient for write-adjacent caches.

## Known Ambiguities

- Exact performance depends on actual data volume and query mix, not just schema correctness.
- Code can look correct while still relying on JS-side filtering after a partial indexed seed — treat that as a performance/correctness finding.
- Distinguish **authoritative persisted state**, **derived persisted state**, and **ephemeral cache state** — most bugs come from mixing these categories.
