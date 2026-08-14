# WAL, Replication & Recovery

> Reference for **postgresql-expert**. Durability (WAL + checkpoints), point-in-time recovery, and the replication topologies built on WAL. Verified-as-of 2026-06-30 (PG18 line).

## Write-Ahead Log (WAL)

PostgreSQL's durability contract: **every change is written to the WAL and flushed to disk before the corresponding data pages are flushed.** On crash, replaying WAL from the last checkpoint reconstructs committed work. WAL is the foundation of both crash recovery *and* all replication and PITR.

- WAL is a sequence of records in 16 MB segment files under `pg_wal/`, addressed by **LSN** (log sequence number, a byte offset).
- `wal_level`: `replica` (default — enough for streaming + PITR) or `logical` (adds row-level decoding for logical replication). `minimal` disables archiving/replication.
- **Full-page writes** (`full_page_writes=on`) guard against torn pages: the first change to a page after a checkpoint logs the whole page. This is why WAL volume spikes right after a checkpoint.
- `synchronous_commit` controls the durability/latency trade-off **per transaction**:
  - `on` (default) — commit waits for local WAL flush (+ sync standbys if configured).
  - `off` — commit returns before flush; risks losing the **last few hundred ms** of *committed* transactions on crash, but **never corrupts** the database. A legitimate per-session choice for bulk loads / non-critical writes.
  - `local`, `remote_write`, `remote_apply` — tune the standby acknowledgement point.

## Checkpoints

A **checkpoint** flushes all dirty shared-buffer pages to data files and records a safe WAL replay-start point. Recovery only replays WAL *after* the last checkpoint.

- Triggered by `checkpoint_timeout` (default 5 min) or when WAL since the last checkpoint hits `max_wal_size` (default 1 GB).
- `checkpoint_completion_target` (default 0.9) spreads checkpoint I/O across the interval to avoid write storms.
- **Tuning:** frequent checkpoints = more full-page writes (WAL bloat) and I/O spikes; raise `max_wal_size` and `checkpoint_timeout` for write-heavy systems to amortize. Watch `pg_stat_bgwriter` / `pg_stat_checkpointer` (PG17+ split) for `checkpoints_req` (forced) vs `checkpoints_timed` — many forced checkpoints means `max_wal_size` is too small.

## Backups & Point-in-Time Recovery (PITR)

Two backup families:

1. **Logical** — `pg_dump`/`pg_dumpall` (portable SQL/custom archive). Good for single DBs, version migration, selective restore; **not** for large-DB fast recovery or PITR.
2. **Physical base backup + WAL archive** — `pg_basebackup` (or `pgBackRest`/`Barman`) takes a filesystem-level snapshot; continuous WAL archiving (`archive_command`/`archive_library`, or streaming to the backup tool) lets you **replay to any point in time**.

PITR flow: restore the base backup → set `recovery_target_time`/`_lsn`/`_name` → PostgreSQL replays archived WAL up to the target → promotes. Use **pgBackRest** or **Barman** in production rather than hand-rolled `archive_command` (parallelism, compression, retention, verification, incremental/block backups). Always **test restores** — an unverified backup is a hope, not a backup.

## Physical (streaming) replication

The primary streams WAL to read-only **hot standbys** that continuously replay it.

- Set up a standby with `pg_basebackup -R` (writes `standby.signal` + `primary_conninfo`); it connects and streams WAL via the **walsender/walreceiver** protocol.
- **Replication slots** (`pg_create_physical_replication_slot`) make the primary retain WAL until the standby has consumed it — preventing "requested WAL segment already removed" errors. **Risk:** a dead/lagging standby with a slot pins WAL and **fills the primary's disk**; bound it with **`max_slot_wal_keep_size`** (PG13+) and alert on `pg_replication_slots`.
- **Synchronous replication**: `synchronous_standby_names` + `synchronous_commit=remote_apply/on` makes commits wait for standby acknowledgement (zero data loss, higher latency). Use a quorum set (`ANY 1 (s1,s2,s3)`) so one slow standby can't stall the primary.
- **`hot_standby_feedback=on`** stops the standby's queries from being cancelled by primary vacuum — at the cost of holding back the **primary's xmin horizon** (bloat trade-off; see MVCC reference).
- Cascading replication and **delayed standbys** (`recovery_min_apply_delay`) protect against erroneous writes propagating instantly.

## Logical replication

Decodes WAL into **row-level changes** and replicates them by `PUBLICATION`/`SUBSCRIPTION` — selective (per-table), cross-version, and writable on the subscriber. The mechanism behind upgrades and CDC.

```sql
-- primary (wal_level=logical):
CREATE PUBLICATION p_orders FOR TABLE orders, order_items;
-- subscriber:
CREATE SUBSCRIPTION s_orders CONNECTION 'host=… dbname=…'
  PUBLICATION p_orders;
```

- Use cases: **near-zero-downtime major-version upgrades**, selective replication, multi-source consolidation, CDC into a pipeline (via `pgoutput`/Debezium).
- PG16+ allows logical replication **from a standby** and bidirectional patterns; PG18 improves conflict logging.
- Caveats: doesn't replicate DDL (schema changes must be applied separately), sequences need manual sync at cutover, large transactions can lag, and a **logical slot also pins WAL** if the subscriber stalls.

## High availability & failover

PostgreSQL ships the replication primitives but **not** automatic failover orchestration — you add a manager:

- **Patroni** (with etcd/Consul/Kubernetes) — the de-facto standard for automated leader election, failover, and config management.
- **`pg_rewind`** resynchronizes a former primary as a standby after failover without a full base backup (needs `wal_log_hints` or data checksums).
- Cloud-managed (RDS/Aurora) abstracts this; the *engine* concepts still apply but Aurora's storage layer replaces streaming replication (see `aws-cloud` → `references/aws-aurora.md`).
- Connection-side failover (target_session_attrs=read-write, pooler-managed) → route to `database-proxies-query-middleware`.

> **Mental model:** WAL is the source of truth for durability *and* replication. Checkpoints bound recovery time; archived WAL enables PITR; streaming/logical replication are just "ship the WAL (or its decoded rows) somewhere else." A neglected replication slot is the most common way to fill a primary's disk.
