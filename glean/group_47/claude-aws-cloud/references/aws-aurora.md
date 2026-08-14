<!-- hub-reference-banner -->
> **Reference file — part of the `aws-cloud` hub.** Created 2026-06-30 as the Aurora **engine deep-
> dive**, upgrading the decision-guide-level Aurora coverage in `references/databases-aws-cockroach-
> indexeddb.md` (which remains the cross-engine "which AWS database?" guide). Sibling topics are
> reference files under `aws-cloud` — **not** standalone skills. For the PostgreSQL engine itself
> (MVCC, vacuum, planner, WAL) use the `postgresql-expert` skill; for DocumentDB-vs-Atlas and the
> RDS/DynamoDB/Neptune decision tree use `references/databases-aws-cockroach-indexeddb.md`; for
> MongoDB use the `mongodb-*` hubs.

---

---
name: aws-aurora
description: >
  Amazon Aurora engine deep-dive — the cloud-native MySQL/PostgreSQL-compatible database and its
  disaggregated storage architecture. TRIGGER: Aurora storage-compute separation; the distributed
  log-structured storage layer; 6-way quorum across 3 AZs (4-of-6 write, 3-of-6 read); "the log is
  the database" / redo-log offload & no double-write; 10 GB protection groups & auto-scaling storage;
  reader endpoints & up to 15 low-lag replicas sharing storage; failover mechanics; Aurora Serverless
  v2 (ACUs, in-place scaling, scale-to-zero); Global Database (storage-level replication, write
  forwarding, managed planned failover, RPO/RTO); Blue/Green deployments; I/O-Optimized vs standard
  pricing; Backtrack (MySQL); cloning; parallel query; Babelfish (SQL Server compat); RDS-vs-Aurora;
  and **Aurora DSQL** [frontier] — the active-active, multi-Region, PostgreSQL-compatible distributed
  database (GA 2025). SKIP: PostgreSQL engine internals (vacuum/planner/WAL/MVCC) → postgresql-expert;
  the cross-engine AWS database decision tree, DocumentDB-vs-Atlas, DynamoDB single-table → references/
  databases-aws-cockroach-indexeddb.md; MongoDB → mongodb-* hubs; connection pooling/RDS Proxy as a
  pooler → database-proxies-query-middleware.
version: "1.0.0"
updated: "2026-06-30"
category: developer
tags:
  - aws
  - aurora
  - database
  - aurora-serverless
  - aurora-dsql
  - global-database
  - storage-architecture
keywords:
  - amazon aurora
  - aurora storage architecture
  - storage compute separation
  - 6-way quorum
  - protection group
  - log is the database
  - aurora replica
  - reader endpoint
  - aurora serverless v2
  - aurora capacity unit
  - scale to zero
  - aurora global database
  - write forwarding
  - blue green deployment
  - aurora i/o-optimized
  - backtrack
  - babelfish
  - parallel query
  - aurora dsql
  - active-active database
  - optimistic concurrency control
---

# Amazon Aurora — Engine Deep-Dive

> Aurora is **wire-compatible** with MySQL and PostgreSQL but is *not* those engines — it replaces their storage layer with a purpose-built distributed system. Understanding that storage layer explains nearly every Aurora behavior (durability, replica lag, failover speed, pricing). Verified-as-of 2026-06-30; pricing/limits/feature gates move — re-verify before quoting numbers.

## The core idea: separate compute from storage

Traditional MySQL/PostgreSQL couple the SQL engine and storage on one host; replication ships the full WAL/binlog to replicas that re-run writes. Aurora **decouples** them:

- **Compute (DB instances)** run the MySQL/PostgreSQL-compatible engine but are largely **stateless** — they cache pages and process SQL.
- **Storage** is a separate, distributed, multi-tenant, auto-scaling fleet shared by all instances in the cluster. It grows automatically in **10 GB segments** up to 128 TiB; you never provision volume size.

### "The log is the database"

Aurora's defining optimization: the compute tier **only writes the redo log** to storage — not dirty data pages. The storage nodes **materialize pages from the log asynchronously**. This eliminates the classic write amplification (no double-write buffer, no full-page writes, no checkpoints flushing pages over the network). It's why Aurora sustains high write throughput and why a crashed instance recovers fast (no redo replay on the compute side — storage already has the log).

### 6-way quorum across 3 AZs

Each 10 GB **protection group** is replicated as **6 copies across 3 Availability Zones (2 per AZ)**. Quorum rules:

- **Writes** need **4 of 6** acknowledgements → tolerates the loss of an **entire AZ** (2 copies) and still commits.
- **Reads** need **3 of 6** → tolerates AZ+1 loss for reads.
- The 4/6 write + 3/6 read quorum (4+3 > 6) guarantees a read quorum always intersects the last write quorum. Storage self-heals failed segments from peers in the background.

The practical payoff: durability and AZ-failure tolerance are properties of **storage**, independent of how many compute instances you run.

## Replicas, endpoints & failover

- Up to **15 Aurora Replicas** share the **same storage volume** as the writer — so they don't replay writes; they just read the shared pages. Replica lag is typically **single-digit milliseconds** (vs seconds for classic replication).
- **Endpoints:** the **cluster (writer) endpoint** always points at the primary; the **reader endpoint** load-balances across replicas; **custom endpoints** group instances by capacity/role.
- **Failover:** because replicas share storage, promotion is fast — typically **~30 seconds or less** (often faster) since there's no data to copy or redo to replay. Set replica **failover priority (tiers)** to control promotion order. Use the cluster endpoint (or RDS Proxy) so apps follow the new writer; connection-side failover tuning → `database-proxies-query-middleware`.

## Aurora Serverless v2

On-demand, fine-grained autoscaling measured in **ACUs (Aurora Capacity Units ≈ 2 GiB RAM + proportional CPU/network)**:

- Scales **in place**, in fractional ACU steps, in **seconds**, without dropping connections — a major improvement over v1's pausing/abrupt scaling. Set a min and max ACU range.
- Supports **scale-to-zero** (pauses to 0 ACU when idle, on supported versions) for dev/intermittent workloads — verify your engine version supports it.
- Mix Serverless v2 and provisioned instances in one cluster (e.g. provisioned writer + serverless readers). Best for variable/spiky or unpredictable workloads; steady high-utilization workloads are often cheaper on provisioned.

## Global, deployment & cost features

- **Aurora Global Database** — one primary Region + up to 5 secondary Regions, replicated at the **storage layer** (dedicated replication infrastructure, not logical replication) with typical cross-Region lag ~1s. **Write forwarding** lets secondary-Region apps send writes to the primary transparently. **Managed planned failover** (low RPO/RTO) for DR drills and Region migration; unplanned promotion for regional outage. This is the standard pattern for low-latency global reads + cross-Region DR.
- **Blue/Green deployments** — create a synchronized green (staging) copy, test schema/engine-version changes, and switch over in ~a minute with safeguards. The safe path for major upgrades.
- **I/O-Optimized** — a pricing mode that removes per-I/O charges (predictable cost) and improves price/performance for I/O-heavy workloads (>~25% of spend on I/O); standard mode charges per I/O. Choose by I/O intensity.
- **Backtrack** (Aurora MySQL) — rewind the cluster to a prior second in place, without restoring a backup (great for "undo a bad migration"). **Fast cloning** — copy-on-write clones of multi-TB databases in minutes. **Parallel Query** (Aurora MySQL) pushes filtering/aggregation into the storage layer. **Babelfish** lets Aurora PostgreSQL accept **SQL Server (TDS/T-SQL)** traffic for migrations.

## Aurora DSQL [frontier]

**Aurora DSQL** (preview re:Invent Dec 2024; **GA May 27, 2025**) is a *separate* offering from Aurora PostgreSQL — a **serverless, distributed, active-active, PostgreSQL-compatible** database:

- **Active-active multi-Region writes** with strong consistency; **99.99% single-Region / 99.999% multi-Region** availability SLAs. No instances to manage — fully serverless, scales independently on read, write, compute, and storage.
- Uses **optimistic concurrency control (OCC)** with snapshot isolation rather than lock-based concurrency — transactions validate at commit and **retry on conflict** (your app must handle retryable commit failures). Disaggregated architecture separates the transaction log, storage, and compute.
- **Not** a drop-in for every Postgres workload at launch: PostgreSQL **compatibility is a subset** (some extensions, features, and behaviors are unsupported) — verify feature support for your workload before committing. Targets scale-out OLTP and multi-Region apps that previously needed sharding.
- Position it against: Aurora Global Database (single-writer + read replicas + write-forwarding) when you need true multi-Region *write* scale-out; and against Spanner/CockroachDB in the distributed-SQL category.

> **Mental model:** classic Aurora = "PostgreSQL/MySQL engine on top of a shared, log-driven, quorum-replicated cloud storage fabric" — most of its magic (fast failover, low replica lag, cheap clones, instant storage growth) falls out of that storage design. **Aurora DSQL** is a different animal: a ground-up distributed SQL database that happens to speak Postgres, for active-active multi-Region scale. Engine-internal PostgreSQL behavior (vacuum, WAL tuning) from `postgresql-expert` applies to Aurora PostgreSQL only loosely — Aurora reimplements the storage layer.
