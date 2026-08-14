# postgresql-expert

**Category:** Databases, Data Engineering & Analytics
**Platform:** Claude
**Original Path:** claude/standalone/postgresql-expert

## Description
PostgreSQL data-plane & engine hub — the open-source RDBMS beneath RDS/Aurora PostgreSQL, Supabase, Neon, and Citus. TRIGGER: MVCC & tuple visibility (xmin/xmax, snapshots, HOT updates, dead tuples), VACUUM/autovacuum, table & index bloat, transaction-ID wraparound & freezing; the query planner & EXPLAIN (ANALYZE, BUFFERS) — seq/index/index-only/bitmap scans, nested-loop/hash/merge joins, cost model, ANALYZE & extended statistics, planner GUCs; index design & access methods (B-tree, GIN, GiST, SP-GiST, BRIN, Hash, partial/expression/covering, skip scan); WAL, checkpoints, PITR/base backup, streaming & logical replication, replication slots, sync commit, failover (Patroni); declarative partitioning (range/list/hash, pruning, partitionwise join/agg); rich types (JSONB, arrays, ranges), the extension model (PostGIS, pg_stat_statements, postgres_fdw, TimescaleDB, Citus) and pgvector for vector search; isolation levels (RC/RR/serializable SSI), locking, deadlocks; CTEs, window functions, MERGE; roles/privileges/RLS. SKIP: connection poolers (PgBouncer/PgCat/Pgpool) & read-write splitting middleware → database-proxies-query-middleware; AWS Aurora PostgreSQL platform/storage/Serverless v2/DSQL → aws-cloud (references/aws-aurora.md); MongoDB engine/MQL → mongodb-expert; generic cross-dialect SQL tuning → deep-query-optimizer; dbt/warehouse modeling → da-data-engineering-platform; RAG retrieval architecture → ai-rag-retrieval.

---

# PostgreSQL Expert — Engine & Data Plane

Expert reference for the **PostgreSQL** engine and data plane — the world's most widely deployed open-source RDBMS and the engine beneath managed offerings (RDS/Aurora PostgreSQL, Supabase, Neon, Crunchy Bridge, Citus, Timescale Cloud). Covers what happens *inside* a PostgreSQL backend: storage and visibility (MVCC), reclamation (VACUUM), planning (the optimizer + statistics), access methods (indexes), durability (WAL + replication), scaling (partitioning), and the type/extension system (JSONB, PostGIS, pgvector).

> Volatile claims (version numbers, feature-availability, extension state) are stamped **verified-as-of 2026-06-30**. **PostgreSQL 18** (released 2025-09-25) is the current major line; PG 13 reached EOL in November 2025. Re-verify version-gated features and EOL dates before quoting them. A new major release or a `pgvector`/extension major is a refresh trigger.

## When this skill applies

Use it to: diagnose MVCC/visibility and VACUUM/bloat/wraparound problems; read EXPLAIN output and reason about the planner, cost model, and statistics; design indexes and pick access methods; configure WAL, checkpoints, PITR, and streaming/logical replication; design declarative partitioning and large-table maintenance; model with JSONB/arrays/ranges and extensions; reason about isolation, locking, and deadlocks; and write advanced SQL.

## Routing table — load the reference for the task

| If the task is about… | Read this reference |
|---|---|
| MVCC, tuple headers (xmin/xmax/ctid), snapshots & visibility, HOT updates, dead tuples, VACUUM/autovacuum tuning, freezing & transaction-ID wraparound, bloat diagnosis & remediation (`pg_repack`, `VACUUM FULL`) | `references/mvcc-vacuum-and-bloat.md` |
| The query planner & optimizer, reading `EXPLAIN` / `EXPLAIN (ANALYZE, BUFFERS)`, scan & join node types, the cost model, `ANALYZE` & `pg_statistic`, extended statistics, `pg_stat_statements`, planner GUCs, common plan pathologies | `references/query-planner-and-statistics.md` |
| Index design & access methods (B-tree, Hash, GIN, GiST, SP-GiST, BRIN), partial / expression / covering / multicolumn indexes, index-only scans, operator classes, `CREATE INDEX CONCURRENTLY`, PG18 skip scan, index bloat | `references/indexes-and-access-methods.md` |
| WAL & checkpoints, `fsync`/durability, crash recovery, base backups & point-in-time recovery, streaming vs logical replication, replication slots, `synchronous_commit` levels, failover & HA (Patroni, `pg_rewind`), `pg_basebackup` | `references/wal-replication-and-recovery.md` |
| Declarative partitioning (range/list/hash), partition pruning (plan-time & run-time), partitionwise join/aggregate, attach/detach, large-table maintenance, retention, and engine-side connection-scaling guidance | `references/partitioning-and-scaling.md` |
| Data types (JSONB, arrays, ranges, hstore, enum, `numeric` vs float), the extension model, key extensions (PostGIS, `pg_stat_statements`, `pg_trgm`, `postgres_fdw`, TimescaleDB, Citus), and **pgvector** (HNSW/IVFFlat, `halfvec`, distance operators) for vector search | `references/types-extensions-and-pgvector.md` |
| Transaction isolation (Read Committed, Repeatable Read, Serializable/SSI), lock modes & the lock hierarchy, deadlocks & lock-wait diagnosis, advisory locks, `SELECT … FOR UPDATE/SKIP LOCKED`, roles/privileges/`SET ROLE`, and Row-Level Security (RLS) | `references/transactions-locking-and-security.md` |

## Core orientation (the 60-second model)

- **PostgreSQL is MVCC-first.** Every row write creates a new *tuple version* tagged with the writing transaction's id (`xmin`); the old version is marked dead (`xmax`). Readers see a *snapshot* and never block writers. The cost is **garbage**: dead tuples must be reclaimed by **VACUUM**, and unreclaimed dead tuples become **bloat**. Almost every "Postgres got slow/big" mystery traces back to MVCC + vacuum behavior. See `mvcc-vacuum-and-bloat.md`.
- **The planner is cost-based and statistics-driven.** It estimates row counts from `ANALYZE` statistics, prices each plan with the `*_cost` GUCs, and picks the cheapest. Bad plans almost always mean **bad estimates** (stale stats, correlated columns, bad `n_distinct`) — not a "missing hint" (PostgreSQL has no query hints). `EXPLAIN (ANALYZE, BUFFERS)` is the single most useful diagnostic. See `query-planner-and-statistics.md`.
- **The right index *type* matters as much as the column.** B-tree for equality/range/sort; **GIN** for JSONB/arrays/full-text; **GiST** for geometry/ranges/nearest-neighbor; **BRIN** for huge naturally-ordered tables; **HNSW** (pgvector) for vectors. Index-only scans need a covering index *and* a visible heap (vacuum-dependent). See `indexes-and-access-methods.md`.
- **Durability rides on the WAL.** Every change is written to the write-ahead log before the data files; checkpoints flush dirty pages; PITR replays WAL to any point; streaming & logical replication ship WAL to replicas. Replication slots prevent WAL removal — and can fill the disk if a replica falls behind. See `wal-replication-and-recovery.md`.
- **Scale tables with declarative partitioning, scale connections elsewhere.** Native range/list/hash partitioning with pruning handles large tables; PostgreSQL's process-per-connection model means heavy connection counts need a **pooler** (route to `database-proxies-query-middleware`). See `partitioning-and-scaling.md`.
- **The type & extension system is the differentiator.** First-class JSONB, arrays, and ranges; a deep extension ecosystem (PostGIS, TimescaleDB, Citus, `pg_stat_statements`); and `pgvector` turning Postgres into a vector DB for RAG. See `types-extensions-and-pgvector.md`.

## Cross-references (do not duplicate — route)

- **Connection poolers & query middleware** (PgBouncer, PgCat, Pgpool-II, read/write splitting, RDS Proxy as a pooler) → `database-proxies-query-middleware`. This skill covers the *engine's* connection/process model and `max_connections`; the pooler tier lives there.
- **AWS Aurora PostgreSQL** (the disaggregated storage engine, 6-way quorum, Serverless v2 ACUs, Global Database, **Aurora DSQL**) → `aws-cloud` → `references/aws-aurora.md`. Aurora is *PostgreSQL-compatible* but a different storage engine; engine-internal Postgres behavior (vacuum, WAL) differs on Aurora.
- **MongoDB engine, MQL, aggregation, WiredTiger** → `mongodb-expert`; **DocumentDB-vs-Atlas** → `aws-cloud`.
- **Generic cross-dialect SQL query optimization with EXPLAIN verification** → `deep-query-optimizer` (`/dqo`). This skill owns Postgres-specific planner reasoning; `/dqo` owns the multi-pass rewrite-and-verify loop.
- **Warehouse/analytics modeling, dbt, ELT** → `da-data-engineering-platform`. **RAG retrieval architecture** (chunking, reranking) → `ai-rag-retrieval`; this skill owns the `pgvector` storage/index layer only.