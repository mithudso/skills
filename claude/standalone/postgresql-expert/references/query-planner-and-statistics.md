# Query Planner, EXPLAIN & Statistics

> Reference for **postgresql-expert**. How the cost-based optimizer plans queries, how to read EXPLAIN, and how statistics drive (and misdrive) plan choice. Verified-as-of 2026-06-30. For cross-dialect rewrite-and-verify loops route to `deep-query-optimizer`.

## The optimizer is cost-based and there are no hints

PostgreSQL parses → rewrites (views, RLS) → **plans** → executes. The planner enumerates access paths and join orders, prices each with a **cost model**, and runs the cheapest. There is **no query-hint syntax** (by design). You influence plans by: fixing statistics, adding/removing indexes, rewriting SQL, or adjusting planner GUCs — not by pinning a plan. (The `pg_hint_plan` extension exists but is not core and is discouraged for general use.)

### The cost model

Costs are in abstract units anchored to `seq_page_cost = 1.0`:

| GUC | Default | Meaning |
|---|---|---|
| `seq_page_cost` | 1.0 | Sequential page fetch |
| `random_page_cost` | 4.0 | Random page fetch — **lower to ~1.1 on SSD/cloud storage** so index scans are costed fairly |
| `cpu_tuple_cost` | 0.01 | Processing one row |
| `cpu_index_tuple_cost` | 0.005 | Processing one index entry |
| `cpu_operator_cost` | 0.0025 | One operator/function eval |
| `effective_cache_size` | 4GB | *Hint* of OS+PG cache available — set to ~50–75% of RAM; raises the appeal of index scans |

`random_page_cost=4` assumes spinning disk. On SSD/EBS/NVMe set it to **1.1–2.0** — this is the single most common config fix for "why is it seq-scanning?".

### Scan node types

- **Seq Scan** — read the whole heap. Cheapest when a large fraction of rows match or the table is tiny.
- **Index Scan** — walk an index, then fetch each matching heap row (random I/O).
- **Index Only Scan** — answer entirely from the index; needs a **covering** index and an **all-visible** page (visibility map → vacuum-dependent). `Heap Fetches:` in EXPLAIN ANALYZE shows VM misses.
- **Bitmap Heap Scan** (+ **Bitmap Index Scan**) — build a bitmap of matching pages, then read them in physical order. Bridges the gap when too many rows for a plain index scan but too few for a seq scan; can combine multiple indexes (`BitmapAnd`/`BitmapOr`).

### Join node types

- **Nested Loop** — for each outer row, probe inner (ideally via index). Best when outer is small. Catastrophic if the planner underestimates outer rows.
- **Hash Join** — build a hash of the smaller side, probe with the larger. Best for big unsorted equijoins; needs `work_mem` (else spills to disk in batches).
- **Merge Join** — sort/scan both inputs in order and merge. Best when inputs are already sorted (e.g. both on indexed keys).

## Reading EXPLAIN

```sql
EXPLAIN (ANALYZE, BUFFERS, VERBOSE, FORMAT TEXT) <query>;
```

- `EXPLAIN` alone = estimates only (no execution). `ANALYZE` actually runs it (wrap DML in a rolled-back transaction). Always add **`BUFFERS`** — shared hit/read tells you what came from cache vs disk.
- **`cost=startup..total`** then **`rows=` / `width=`** (estimates). With `ANALYZE`: **`actual time=startup..total rows= loops=`**.
- **The killer signal: estimated `rows` vs `actual rows`.** A large divergence (e.g. est 1, actual 2M) means **bad statistics** and is the root of most bad plans — it cascades into the wrong join type/order.
- `loops` > 1 on a nested-loop inner node: multiply `actual time * loops` for true cost.
- Look for: `Rows Removed by Filter` (index not selective / missing index), `Heap Fetches` (VM stale → vacuum), `Sort Method: external merge Disk` (raise `work_mem`), `Batches:` >1 on Hash (spilled).

## Statistics — the foundation of good plans

The planner estimates row counts from per-column statistics gathered by **`ANALYZE`** (and autovacuum's analyze pass) into `pg_statistic` (readable via `pg_stats`):

- `null_frac`, `n_distinct`, **most-common values (MCVs)** + their frequencies, and a **histogram** of the rest. `default_statistics_target` (default 100) sets MCV/histogram resolution; raise per-column for skewed data: `ALTER TABLE t ALTER c SET STATISTICS 1000;`.
- **`n_distinct` errors** wreck estimates on big tables (sampling under-counts distinct values). Override: `ALTER TABLE t ALTER c SET (n_distinct = -0.5);` (negative = fraction of rows).

### Extended statistics — for correlated columns

The planner assumes columns are **independent**; correlated predicates (e.g. `city` and `zip`) get badly under/over-estimated. Fix with `CREATE STATISTICS`:

```sql
CREATE STATISTICS s_geo (dependencies, ndistinct, mcv)
  ON city, state, zip FROM addresses;
ANALYZE addresses;
```

- `dependencies` — functional dependencies (zip → state).
- `ndistinct` — multi-column distinct counts (for GROUP BY estimates).
- `mcv` — multi-column most-common-value lists.

## pg_stat_statements — find the queries that matter

The essential extension for production triage. Aggregates normalized query stats:

```sql
SELECT query, calls, total_exec_time, mean_exec_time, rows,
       100.0*shared_blks_hit/nullif(shared_blks_hit+shared_blks_read,0) AS hit_pct
FROM pg_stat_statements
ORDER BY total_exec_time DESC LIMIT 20;
```

Sort by **`total_exec_time`** (aggregate pain), not just `mean` — a fast query called millions of times often dominates. Pair with `auto_explain` (logs plans of slow statements) and `log_min_duration_statement` to capture offenders with their parameters.

## Common plan pathologies & fixes

| Symptom | Likely cause | Fix |
|---|---|---|
| Seq scan where an index exists | `random_page_cost` too high; index not selective; type mismatch / non-sargable predicate | Lower `random_page_cost`; check `WHERE` is sargable (no `func(col)` unless expression index) |
| Nested loop blows up | Underestimated outer rows (bad stats / correlation) | `ANALYZE`; extended statistics; raise statistics target |
| Sort/Hash spills to disk | `work_mem` too low | Raise `work_mem` (per-operation, per-connection — size carefully) |
| Plan flips after data growth | Crossed a cost threshold; stale stats | More aggressive autoanalyze; extended stats |
| `OR` across columns slow | Planner can't use one index | Rewrite as `UNION`/`UNION ALL`, or add a bitmap-combinable index per branch |
| Generic plan worse than custom (prepared stmts) | `plan_cache_mode` chose generic | `SET plan_cache_mode = force_custom_plan` for skewed params |

> **Mental model:** a "slow query" is usually a **bad estimate**, not a missing index. Always look at estimated-vs-actual rows first. Fix statistics before touching GUCs, and fix sargability/indexing before blaming the planner.
