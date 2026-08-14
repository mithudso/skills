# Indexes & Access Methods

> Reference for **postgresql-expert**. Choosing the right index type and shape, and the operational rules for building them safely. Verified-as-of 2026-06-30 (PG18 line).

## Pick the access method first

The biggest indexing mistake is reaching for a B-tree on everything. Match the method to the query and data:

| Method | Best for | Operators |
|---|---|---|
| **B-tree** (default) | Equality, range, `ORDER BY`, uniqueness; the all-rounder | `= < <= > >= BETWEEN IN`, `LIKE 'prefix%'`, `IS NULL` |
| **Hash** | Equality only; smaller than B-tree for long keys (crash-safe & replicated since PG10) | `=` only |
| **GIN** | "Multiple values in one row": JSONB, arrays, full-text `tsvector`, `pg_trgm` | `@> ? ?& ?|`, `@@`, trigram `%` |
| **GiST** | Geometry (PostGIS), ranges, nearest-neighbor (`<->`), exclusion constraints, full-text | `&& @> <@ <->`, range overlap |
| **SP-GiST** | Non-balanced / partitioned structures: quadtrees, radix trees, IP/`inet`, points | `<< >> <@`, prefix |
| **BRIN** | **Huge** tables with strong physical-order correlation (append-only time/sequence) — tiny index, lossy | range `< > BETWEEN` on correlated column |
| **HNSW / IVFFlat** (pgvector) | Approximate nearest-neighbor vector search | `<-> <=> <#>` (see types-extensions reference) |

Rules of thumb:
- **JSONB containment / arrays / full-text → GIN.** Use `jsonb_path_ops` GIN for `@>` only (smaller/faster) vs default `jsonb_ops` for key-existence too.
- **Range types, geometry, "find the nearest" → GiST.**
- **Append-only fact table ordered by time, billions of rows → BRIN** (a BRIN index can be megabytes where a B-tree would be hundreds of GB), but only if physical order tracks the column (cluster/insert order).
- **`inet`/text-prefix/phone trees → SP-GiST.**

## Index shapes that change everything

### Multicolumn (composite) indexes
Column order follows the **equality-first, then range/sort** rule. `(tenant_id, created_at)` serves `WHERE tenant_id = ? ORDER BY created_at` and `WHERE tenant_id = ? AND created_at > ?`. A leading-column-only predicate uses it; a trailing-column-only predicate generally cannot (until **skip scan**, below).

### Partial indexes
Index only the rows you query: `CREATE INDEX ON orders (created_at) WHERE status = 'open';`. Smaller, cheaper to maintain, and great for "hot subset of a big table" (soft-deletes, open tickets, unprocessed rows).

### Expression (functional) indexes
Index a computed value so a non-sargable predicate becomes indexable: `CREATE INDEX ON users (lower(email));` then query `WHERE lower(email) = $1`. Required whenever you wrap the column in a function.

### Covering indexes & index-only scans
`CREATE INDEX ON t (a, b) INCLUDE (c);` stores `c` in the leaf without making it part of the key. If a query reads only `a, b, c`, the planner can do an **index-only scan** — *provided the heap pages are all-visible* (visibility map set by VACUUM). Watch `Heap Fetches:` in `EXPLAIN ANALYZE`; high fetches mean vacuum isn't keeping the VM current.

### Unique, exclusion & deferrable
- `UNIQUE` indexes enforce uniqueness (incl. `NULLS NOT DISTINCT` in PG15+).
- **Exclusion constraints** (GiST) enforce "no two rows overlap" — e.g. no double-booked time ranges: `EXCLUDE USING gist (room WITH =, during WITH &&)`.

### PG18: skip scan
PostgreSQL 18 adds **skip scan** for multicolumn B-trees: a query filtering only on a *non-leading* column can still use the index by "skipping" through distinct leading-column values when the leading column has low cardinality. Reduces the need for redundant single-column indexes — but verify with EXPLAIN; it isn't a universal win.

## Building & maintaining indexes safely

- **`CREATE INDEX CONCURRENTLY` (CIC)** builds without an `ACCESS EXCLUSIVE` lock — essential in production. Caveats: it's slower, can't run in a transaction block, and on failure leaves an **`INVALID`** index you must drop (`DROP INDEX CONCURRENTLY`). Check `pg_index.indisvalid`.
- **`REINDEX INDEX CONCURRENTLY`** (PG12+) rebuilds a bloated index online.
- On **partitioned tables**, create the index on the parent `ON ONLY` then build each partition's index concurrently and `ATTACH` — or use partitionwise creation.

## Finding missing & unused indexes

- **Unused indexes** waste write bandwidth and disk: `pg_stat_user_indexes` where `idx_scan = 0` (over a representative window) — drop them (write amplification + bloat).
- **Missing-index signals**: high `Rows Removed by Filter`, seq scans on large tables in `pg_stat_user_tables` (`seq_scan` >> `idx_scan`), or repeated slow queries in `pg_stat_statements` that filter on an unindexed column.
- **Duplicate/overlapping indexes**: a single-column index is redundant if it's the leading column of an existing composite (post-skip-scan, even more so).

## Why an index isn't used (checklist)

1. **Non-sargable predicate** — `WHERE func(col) = x`, `col::text = …`, `col + 1 = y`, leading-wildcard `LIKE '%x'`. Fix with an expression index or rewrite.
2. **Type mismatch** — comparing `bigint` column to a `text`/`numeric` literal forces a cast that defeats the index.
3. **Low selectivity** — predicate matches most rows; seq scan is genuinely cheaper.
4. **`random_page_cost` too high** (see planner reference) — index scans mispriced on SSD.
5. **Stale stats** — planner thinks the predicate is unselective; `ANALYZE`.
6. **Collation / opclass mismatch** — `LIKE 'x%'` needs `text_pattern_ops` (or a `C`-collation index) to use a B-tree for prefix matching on non-`C` collations.

> **Mental model:** the column tells you *what* to index; the **query shape and access method** tell you *how*. Build with `CONCURRENTLY`, cover what you read, prune to partial where you can, and delete indexes nothing scans.
