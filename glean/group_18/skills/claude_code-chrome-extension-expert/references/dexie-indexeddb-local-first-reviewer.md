The file being fixed is not `SKILL.md` — the task wants the corrected compressed content returned directly. The fix is restoring the second `chrome-extension-expert` backtick occurrence in the banner (line 2: `under hubs` → `under hubs (\`chrome-extension-expert\`)`).

<!-- hub-reference-banner -->
> **Reference file — part of `chrome-extension-expert` hub.** Formerly standalone `dexie-indexeddb-local-first-reviewer` skill.
> Sibling topics now reference files under hubs (`chrome-extension-expert`) — **not** standalone skills. Ignore "use the X skill" / `related_skills` / SKIP pointers naming bare sibling skills; load topic's `references/<name>.md` from owning hub (see hub's "Cross-hub map").

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

Review reference for auditing Dexie and IndexedDB in local-first apps. Covers schema/index correctness, migration safety, transaction boundaries, cache invalidation, MV3 persistence.

**Sources:** [Dexie docs](https://dexie.org/docs/), [MDN IndexedDB API](https://developer.mozilla.org/en-US/docs/Web/API/IndexedDB_API), [Chrome MV3 service worker lifecycle](https://developer.chrome.com/docs/extensions/develop/concepts/service-workers/lifecycle).

*Context verified: 2026-05-10. Repo context: `docs/ARCHITECTURE.md`, `packages/store/src/db.ts`.*

## Quick Review Rules

1. **Index only what you query.** Field in `stores()` only if code sorts/filters on it. Over-indexing hurts performance; Dexie explicitly warns against large indexed values.
2. **Compound index order must match query shape.** Declared `[a+b]`; query prefix must align or index skipped.
3. **Schema changes are migrations.** New indexes and renamed fields go in versioned `upgrade()` calls, not ad hoc runtime fixups.
4. **Wrap related writes in one transaction.** Control-yielding async gaps break atomicity even inside Dexie transactions.
5. **Prefer batch writes.** Per-document write loops suspicious — `bulkAdd()`/`bulkPut()` exist for reason.
6. **Every cache needs explicit invalidation trigger** near write paths that can stale it.
7. **MV3 globals are disposable.** Anything surviving suspend goes in IndexedDB or `chrome.storage.*`, not memory.

## Review Workflow

| Step | What to check | Key APIs |
|---|---|---|
| 1. Map access patterns | Hot reads, sorted/range queries, fan-out joins before judging schema | `stores()`, `where()` |
| 2. Verify index coverage | Each `where()`, `orderBy()`, compound selector → actual schema index exists, order matches | `where()`, `IDBIndex` |
| 3. Inspect write atomicity | Related writes, migration flags, derived-record updates in one transaction or partial-commit risk | `transaction()`, `IDBTransaction` |
| 4. Check migration safety | New indexes, renamed fields, backfills versioned and idempotent; first-open cannot loop | `upgrade()`, `IDBDatabase` |
| 5. Audit cache invalidation | Thread summaries, unread counts, prompt caches cannot drift after writes | repo `docs/ARCHITECTURE.md` |
| 6. Check failure modes | Quota/eviction assumptions, re-open behavior, MV3 suspend/resume | MDN storage quotas, Chrome lifecycle |

## Method Reference

| Method / surface | Purpose | Review focus | Caveats |
|---|---|---|---|
| `version().stores()` | Declares tables and indexes | Index only queried fields; verify compound index order | Over-indexing hurts performance and stability |
| `version().upgrade()` | Versioned schema/data migration | Idempotence, field backfill safety, old-to-new mapping | Runs inside upgrade transaction; mistakes strand users on open |
| `transaction()` | Atomic read/write group | Async gaps, nested scope, rollback behavior | IndexedDB auto-commits when control yields unexpectedly |
| `bulkAdd()` / `bulkPut()` | Batch insert/upsert | Per-doc loops, `BulkError` handling, partial-success assumptions | Insert-only vs upsert semantics differ — don't use `bulkAdd()` where idempotence requires `bulkPut()` |
| `where()` / `WhereClause` | Indexed lookup/range query | Selector-to-index alignment, fallback filtering | Misaligned queries silently degrade to JS-side filtering |
| `IDBKeyRange`-style range | Efficient lower/upper bounds | Inclusive/exclusive edge correctness | Easy off-by-one errors |
| `Dexie.waitFor()` | Keep transaction alive across async | Whether truly necessary | CPU-expensive in hot paths |
| Repo cache declarations | Fast-path derived reads | Explicit invalidation, write adjacency, bounded growth | TTL-only caches hide stale state if writes don't bust them |

## Standards

- **Compound indexes for real multi-field access patterns** — not broad JS-side filtering on weak seed.
- **Migration logic explicit, versioned, one-way** — no live shape mutations without upgrade path.
- Distinguish **insert-only** (`bulkAdd`) vs **upsert** (`bulkPut`) deliberately.
- Treat **quota and eviction as variable browser behavior**, not guaranteed capacity.
- MV3: persist **authoritative state** to IndexedDB/storage; in-memory caches = performance helpers only.
- Every cache needs **explicit invalidation trigger** — TTL alone insufficient for write-adjacent caches.

## Known Ambiguities

- Actual performance depends on data volume and query mix, not just schema correctness.
- Code can look correct while relying on JS-side filtering after partial indexed seed — treat as performance/correctness finding.
- Distinguish **authoritative persisted state**, **derived persisted state**, **ephemeral cache state** — most bugs from mixing these.