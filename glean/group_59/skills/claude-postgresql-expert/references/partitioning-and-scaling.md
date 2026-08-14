# Partitioning & Scaling

> Reference for **postgresql-expert**. Declarative partitioning for large tables and the engine-side limits of scaling a single PostgreSQL instance. Verified-as-of 2026-06-30. Connection poolers route to `database-proxies-query-middleware`; distributed PostgreSQL routes to Citus (types-extensions reference) or `aws-cloud` (Aurora DSQL).

## Declarative partitioning

Native partitioning (PG10+, mature from PG12) splits one logical table into child **partitions**, each a real table. The parent is empty and routes rows to partitions.

```sql
CREATE TABLE events (id bigint, ts timestamptz, ...) PARTITION BY RANGE (ts);
CREATE TABLE events_2026_06 PARTITION OF events
  FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');
```

Three strategies:
- **RANGE** — time-series, sequential ids (the most common; one partition per day/week/month).
- **LIST** — discrete categories (region, tenant tier, status).
- **HASH** — even distribution when there's no natural range key (`PARTITION BY HASH (customer_id)` with N modulus partitions) to spread write hotspots.

Sub-partitioning (range-then-hash) is supported but adds complexity — prefer a single level unless you have a clear two-dimensional access pattern.

## Why partition (and why not)

Partitioning is **not** a general performance button. It pays off for specific reasons:

1. **Cheap bulk deletion / retention** — `DETACH`/`DROP` a whole partition is instant and bloat-free, vs a massive `DELETE` that generates dead tuples and vacuum load. This is the #1 real reason to partition time-series data.
2. **Partition pruning** — the planner eliminates irrelevant partitions when the query filters on the partition key, scanning far less data.
3. **Smaller per-partition indexes** that fit in cache; faster, more targeted autovacuum per partition.
4. **Partitionwise join/aggregate** (enable `enable_partitionwise_join`, `enable_partitionwise_aggregate` — off by default) push work down per matching partition pair.

Costs/caveats:
- Queries that **don't** filter on the partition key touch *every* partition (worse than an unpartitioned table) and need the key in unique/PK constraints.
- Too many partitions (thousands) raise planning time and memory; keep partition counts sane (favor monthly over daily unless volume demands it).
- Foreign keys *referencing* a partitioned table were limited historically (improved in PG12+); test on your version.

## Pruning, ATTACH/DETACH, maintenance

- **Plan-time pruning** uses constant predicates; **run-time pruning** (PG11+) handles parameterized/`IN`/join-driven predicates and `PREPARE`d statements. Confirm with EXPLAIN (look for `Subplans Removed` / fewer scanned partitions).
- **`ATTACH PARTITION`** validates the new partition against the constraint (use a matching `CHECK` first to make it a metadata-only, non-blocking attach in modern versions); **`DETACH PARTITION CONCURRENTLY`** (PG14+) removes a partition without a long lock.
- Automate partition lifecycle with **`pg_partman`** (creates future partitions, retains/drops old ones) or TimescaleDB hypertables (which automate chunk management for time-series).
- Create indexes on the parent (`CREATE INDEX ON events (...)`) and PostgreSQL cascades them to all partitions, including future ones.

## Scaling a single instance — engine limits

PostgreSQL uses a **process-per-connection** model: each connection is a forked backend (~5–15 MB RSS plus `work_mem` per sort/hash). Implications:

- **`max_connections` is not a throughput dial.** A few hundred busy backends can thrash CPU and memory. The fix is **not** raising `max_connections` — it's a **connection pooler** (PgBouncer transaction mode, PgCat, RDS Proxy) so thousands of clients share a small backend pool. Route pooler design to `database-proxies-query-middleware`.
- **`shared_buffers`** ~25% of RAM (let the OS page cache hold the rest; `effective_cache_size` advertises the total). **`work_mem`** is *per operation per connection* — multiply by concurrency before raising it (a common OOM cause). **`maintenance_work_mem`** sizes vacuum/index builds.
- **Vertical scaling** (bigger box, faster NVMe) is the simplest lever and goes very far for OLTP.
- **Read scaling** = streaming-replica read replicas (route reads via the pooler/app); accept replica lag and the `hot_standby_feedback` bloat trade-off (WAL reference).
- **Write/horizontal scaling beyond one node** = application sharding, **Citus** (distributed tables, reference table replication, distributed query planner — see types-extensions reference), or **Aurora DSQL** (active-active distributed PostgreSQL-compatible — `aws-cloud` → `references/aws-aurora.md`). Vanilla PostgreSQL has a single write primary.

## Large-table operational playbook

- **Schema changes**: `ALTER TABLE … ADD COLUMN` with a non-volatile default is metadata-only (PG11+) — fast. Adding a `NOT NULL` constraint can be validated `NOT VALID` then `VALIDATE CONSTRAINT` to avoid a long full-table lock. Backfills should be chunked with vacuums between.
- **Big deletes**: prefer partition `DROP`/`DETACH`; otherwise batch and vacuum to control bloat.
- **`CLUSTER`** physically reorders a table by an index (one-time, takes `ACCESS EXCLUSIVE`); helps BRIN/range-scan locality but isn't maintained automatically.
- **`COPY`** for bulk load (far faster than row-by-row `INSERT`); drop/rebuild non-essential indexes around very large loads; raise `maintenance_work_mem` for the index build.

> **Mental model:** partition to make *retention and pruning* cheap, not for generic speed; scale *connections* with a pooler (not `max_connections`); scale *reads* with replicas; and reach for Citus/DSQL/sharding only when one write primary truly isn't enough.
