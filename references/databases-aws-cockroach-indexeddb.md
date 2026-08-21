<!-- hub-reference-banner -->
> **Reference file — part of the `aws-cloud` hub.** Formerly the standalone `databases-aws-cockroach-indexeddb` skill.
> Sibling topics in this family are now reference files under the hubs (`aws-cloud`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: databases-aws-cockroach-indexeddb
title: "Databases: AWS, DocumentDB, CockroachDB, IndexedDB"
description: >
  AWS database ecosystem decision guide covering RDS, Aurora, DynamoDB, ElastiCache,
  DocumentDB, Neptune, Timestream, QLDB, and DMS. Also covers Amazon DocumentDB
  compatibility gaps vs MongoDB Atlas, CockroachDB distributed SQL patterns, and
  browser-side IndexedDB storage.
  TRIGGER: user asks which AWS database to use, DynamoDB single-table design,
  Aurora Serverless v2 scaling, DocumentDB vs MongoDB Atlas, DocumentDB retryWrites,
  CockroachDB multi-region or serializable retry errors, IndexedDB browser storage limits,
  Safari ITP eviction, or Dexie.js patterns.
  SKIP: for MongoDB Atlas features (search, vector search, Atlas Functions) use mongodb-atlas-expert;
  for core MongoDB MQL, aggregation, indexes, schema design, transactions, change streams, or
  WiredTiger internals use mongodb-expert; for reviewing Dexie/IndexedDB local-first code use
  dexie-indexeddb-local-first-reviewer; for migrating data between databases use database-migrations;
  for AWS Lambda + serverless architecture use aws-serverless.
category: developer
version: "1.2.1"
updated: "2026-05-31"
tags: [aws, databases, dynamodb, aurora, rds, documentdb, mongodb, cockroachdb, indexeddb, dexie]
keywords:
  - DynamoDB single table design
  - Aurora Serverless v2
  - Aurora DSQL
  - DocumentDB vs MongoDB Atlas
  - DocumentDB compatibility
  - CockroachDB multi-region
  - CockroachDB RETRY_SERIALIZABLE
  - IndexedDB browser storage
  - Dexie.js
  - ElastiCache Valkey
  - RDS Proxy Lambda
  - AWS DMS migration
  - DynamoDB hot partition
whenToUse:
  - Deciding which AWS database service to use for a given workload
  - Comparing DocumentDB vs MongoDB Atlas for an AWS-native MongoDB workload
  - Configuring Aurora Serverless v2 min/max ACU and scaling behavior
  - Diagnosing DocumentDB retryWrites errors or change stream stalls
  - Designing CockroachDB multi-region table locality (REGIONAL BY ROW, GLOBAL)
  - Handling RETRY_SERIALIZABLE errors in CockroachDB application code
  - Building offline-first browser apps with IndexedDB and Dexie.js
  - Debugging Safari ITP data eviction in a PWA
  - Choosing between localStorage, IndexedDB, and Cache API
  - Setting up connection pooling for Lambda + RDS/Aurora with RDS Proxy
  - Planning a database migration with AWS DMS
whenNotToUse:
  - MongoDB Atlas Search, Vector Search, or Atlas App Services (use mongodb-atlas-expert)
  - Reviewing Dexie/IndexedDB local-first code for correctness (use dexie-indexeddb-local-first-reviewer)
  - Migrating data between any databases (use database-migrations)
  - AWS serverless architecture beyond database questions (use aws-serverless)
related_skills:
  - software-engineering-patterns
  - database-migrations
  - mongodb-atlas-expert
  - aws-serverless
  - mongodb-expert
---

# Databases: AWS, DocumentDB, CockroachDB, IndexedDB

## Overview

**Audience**: Software engineers and architects choosing, configuring, or migrating databases.

**How to use this skill**: This is a reference — consult it for decision criteria, compatibility facts, and patterns. It does not prescribe a single workflow; navigate to the relevant section for your database or use the decision tree in "Choosing the Right Database."

**What you can do with this skill:**
- Choose the right AWS database service for a given workload
- Understand DocumentDB's real compatibility limits vs MongoDB Atlas
- Deploy CockroachDB in multi-region configurations with correct survivability expectations
- Build offline-first browser apps with IndexedDB and Dexie.js
- Avoid the most common migration and integration mistakes for each system

This skill covers the AWS database ecosystem, Amazon DocumentDB (and its divergence from MongoDB), CockroachDB distributed SQL, and browser-side IndexedDB storage, as of 2025–2026. A response from this skill is correct when it identifies the right service for the stated workload, accurately represents compatibility limits, and flags known anti-patterns before they are committed to.

> **Staleness note:** Version numbers, feature availability, and pricing figures were current as of May 2026. Verify against current AWS documentation and vendor release notes before making architecture decisions.

**Navigation by task:**
- Choose an AWS database service → AWS Database Landscape + Decision Heuristics
- DynamoDB key design / single-table → DynamoDB section
- Aurora Serverless v2 sizing → Aurora section
- Aurora DSQL trade-offs → Aurora DSQL section
- DocumentDB vs MongoDB Atlas → Amazon DocumentDB section
- CockroachDB multi-region patterns → CockroachDB section
- Browser-side storage (IndexedDB, Dexie, Safari ITP) → IndexedDB section
- Pick between storage options (localStorage, IndexedDB, Cache API) → IndexedDB → localStorage vs IndexedDB table
- Troubleshooting symptoms by service → `references/troubleshooting.md`

> **Compatibility note**: Percentage compatibility figures (DocumentDB ~34%, CockroachDB ~40%) are sourced from third-party analyses as noted inline. These numbers drift as services evolve — verify against current AWS documentation and the CockroachDB compatibility page before making architectural decisions.

---

## AWS Database Landscape

AWS offers 15+ purpose-built database services. The core principle: **no single database is best for everything**. Match the workload to the engine.

### Service Map by Data Model

| Model | Service | Best For |
|---|---|---|
| Relational | RDS (PostgreSQL, MySQL, MariaDB, Oracle, SQL Server) | Traditional OLTP, complex joins |
| Relational (enhanced) | Aurora (PostgreSQL/MySQL compatible) | High-throughput production relational |
| Distributed relational | Aurora DSQL | Active-active multi-region SQL (GA May 2025) |
| Key-value / Document | DynamoDB | Millions of req/s, known access patterns |
| In-memory | ElastiCache (Redis/Valkey, Memcached) | Sub-ms caching, sessions, rate limiting |
| Document | DocumentDB | Simple CRUD on AWS — check compatibility first |
| Graph | Neptune | Relationships, social graphs, fraud detection |
| Time-series | Timestream | IoT sensor data, metrics, telemetry |
| Ledger | QLDB | Immutable audit logs, financial records |
| Wide-column | Keyspaces (Cassandra compatible) | High-volume append-heavy workloads |
| Search | OpenSearch Service | Full-text search, log analytics |

### Decision Heuristics (2025–2026)

1. **Need complex ad-hoc queries or joins?** → RDS/Aurora
2. **High throughput, predictable key-based access at any scale?** → DynamoDB
3. **Active-active multi-region SQL?** → Aurora DSQL
4. **MongoDB workload but want managed AWS service?** → DocumentDB only if feature set is sufficient (~34% MongoDB API compatible); otherwise MongoDB Atlas
5. **Sub-millisecond read latency, caching layer?** → ElastiCache
6. **Graph traversal queries?** → Neptune
7. **Time-ordered sensor/metric data?** → Timestream
8. **Tamper-proof audit log?** → QLDB

> **2025–2026 Trend**: Databases are now active components in AI pipelines. ElastiCache provides semantic caching for vector search. DynamoDB and Aurora are used as vector stores alongside Aurora pgvector.

---

## RDS & Aurora

### RDS (Relational Database Service)

Managed relational database for MySQL, PostgreSQL, MariaDB, Oracle, SQL Server.

**When to use RDS over Aurora:**
- Development/staging environments (cost)
- Database engines Aurora doesn't support (Oracle, SQL Server)
- Low-traffic applications where Aurora's overhead isn't justified

**Key features:**
- Automated backups, point-in-time restore
- Multi-AZ for high availability (synchronous standby)
- Read replicas (up to 5 for MySQL/PostgreSQL; Aurora supports up to 15)
- Storage auto-scaling up to 64 TB

### Aurora

MySQL/PostgreSQL-compatible with distributed storage architecture. ~20% more expensive than RDS but operationally superior for production.

**Aurora advantages over RDS:**
- Storage auto-scales to 128 TB (no pre-provisioning)
- Faster failover (typically <30s vs RDS ~60–120s)
- Up to 15 read replicas (vs 5 for RDS)
- 6-way storage replication across 3 AZs
- Aurora Global Database for low-latency cross-region reads (~1ms for local replica reads)

> **Note on Global Database latency**: ~1ms applies to reads served by the local regional replica. Write-forwarding to the primary region still incurs cross-region round-trip latency.

### Aurora Serverless v2

Scales compute in fine-grained ACU (Aurora Capacity Unit) increments with sub-second response. GA 2022, significantly improved in 2025.

**2025 improvements:**
- Doubled default scaling rate; scales 45% faster (0.5 ACU to 256 ACU)
- Scales to zero capability added (latest platform versions)

**Best practices:**
```yaml
# Capacity configuration
min_capacity: 0.5  # ACU — scale to zero for dev environments
max_capacity: 128  # ACU — set based on peak workload + buffer

# Hybrid pattern: cost-optimized
# Writer:  Provisioned instance  (stable baseline, RI-eligible)
# Readers: Serverless v2        (elastic read scaling)
```

**Always pair with RDS Proxy** when used behind Lambda (connection pooling).

**Scaling caveats:**
- First connection after scale-to-zero has cold start latency — set `min_capacity > 0` for latency-sensitive apps
- Min ACU of 0.5 still incurs storage costs
- Don't use Serverless v2 for sustained high-throughput workloads (provisioned is cheaper at high baseline)

### Aurora DSQL

Fully serverless, distributed SQL database for active-active multi-region writes. Previewed at re:Invent (December 2024) and reached general availability **May 27, 2025**.

> For the Aurora engine deep-dive (storage architecture, 6-way quorum, Serverless v2 internals, Global Database, and Aurora DSQL's OCC/snapshot-isolation model) see the sibling reference `aws-aurora.md`. This file remains the cross-engine "which AWS database?" decision guide.

**Key properties:**
- Active-active multi-region writes with strong consistency
- PostgreSQL-compatible query surface
- No cluster management; scales to zero
- Designed for globally distributed applications needing SQL semantics

**When to use:**
- Multi-region active-active write workloads with SQL
- Applications requiring global low-latency writes without custom sharding

**Current limitations (2025–2026):**
- No foreign key constraints
- No DDL inside transactions (schema changes are auto-committed)
- No stored procedures or triggers
- No cross-region joins (data stays in its home region)
- Expanding regional availability through 2026
- Not suitable for complex OLAP workloads

### RDS Proxy

Connection pooler that sits between application and RDS/Aurora. Critical for Lambda and serverless architectures.

- Maintains persistent connection pool to the database
- Applications connect to the Proxy endpoint, not directly to the DB
- Reduces connection overhead by up to 67%
- Supports IAM authentication
- Failover handled transparently (~5s vs ~30s without Proxy)

```text
Lambda → RDS Proxy → Aurora
         (pool of N persistent connections)
```

---

## DynamoDB

Fully managed, serverless key-value and document database. Delivers single-digit millisecond performance at any scale.

### Architecture Concepts

- **Table**: Collection of items (no schema except partition key)
- **Partition Key (PK)**: Determines which partition stores the item. Must distribute evenly.
- **Sort Key (SK)**: Optional; enables range queries within a partition
- **GSI (Global Secondary Index)**: Alternate access pattern; asynchronous replication
- **LSI (Local Secondary Index)**: Alternate sort key; must be defined at table creation

### Single-Table Design

Store multiple entity types in one table. Motivated by: *data that is fetched together should be stored together*.

**Steps:**
1. List ALL access patterns before designing keys
2. Design PK/SK to support the most common patterns
3. Use entity prefixes to distinguish types: `USER#123`, `ORDER#456`
4. Use GSIs only for access patterns the base table can't serve
5. Embrace denormalization — duplicate data to avoid multiple round-trips

**Example pattern:**
```text
PK            SK             Data
USER#u1       PROFILE        {name, email}
USER#u1       ORDER#o1       {total, date}
USER#u1       ORDER#o2       {total, date}
ORDER#o1      META           {status, userId}
```

**GSI overloading:**
```text
GSI1_PK       GSI1_SK        (sparse index)
STATUS#active ORDER#o1       (query all active orders)
```

### Capacity Modes

- **On-Demand**: Pay per request. Best for unpredictable traffic.
- **Provisioned + Auto-Scaling**: Better cost at predictable throughput.

### DynamoDB Accelerator (DAX)

In-memory cache for DynamoDB. Delivers microsecond read latency (vs. single-digit ms for DynamoDB directly).

- Drop-in compatible with DynamoDB API
- Use for read-heavy workloads where even ms latency matters
- Not appropriate for strongly consistent reads (DAX returns eventually consistent results)
- Adds operational cost and complexity — add only when DynamoDB latency is the measured bottleneck

### DynamoDB Streams + Lambda

Event-driven pattern: DynamoDB Streams captures item-level changes and triggers Lambda.
```text
DynamoDB Table → Stream → Lambda → SNS/SQS/EventBridge
```

### Best Practices

- **Avoid hot partitions**: Distribute writes across many partition key values
- **Use sparse GSIs**: Only project items that need the index
- **Transactions**: Use `TransactWriteItems` for atomic multi-item updates (up to 100 unique items)
- **Condition expressions**: Optimistic locking with version attributes
- **TTL**: Use DynamoDB TTL for automatic expiration (no extra cost)
- **Pagination**: Always handle `LastEvaluatedKey` — never assume one query returns all results

### When NOT to Use DynamoDB

Evaluate alternatives if any of the following apply:
- Complex reporting or analytics across many attributes → use Redshift or Aurora
- Ad-hoc queries you can't predict at design time → use Aurora/RDS (SQL is more flexible)
- Many-to-many relationships requiring flexible joins → relational model is clearer

> **Migration warning**: Migrating from DynamoDB to a relational database is extremely painful — there is no schema to export, key design assumptions don't translate, and data must be re-modeled from scratch. Confirm all access patterns before committing to DynamoDB.

**Practical limits to know:**
- Item size: 400 KB maximum per item
- Partition throughput: 3,000 RCU and 1,000 WCU per partition key value (design keys to spread load)

---

## ElastiCache

Managed in-memory caching service. Two engines: **Redis/Valkey** and **Memcached**.

> **2025 Note**: AWS is transitioning from Redis OSS to **Valkey** (open-source Redis fork). New deployments should use Valkey unless Redis-specific features are required.

### Redis / Valkey Use Cases

- Session storage
- Rate limiting (token bucket with INCR + TTL)
- Pub/Sub messaging
- Sorted sets for leaderboards
- Distributed locks
- Semantic caching for LLM/AI workloads (2025 pattern)

### Memcached Use Cases

- Simple string caching
- Multi-threaded, horizontal scale
- No persistence needed

### ElastiCache Serverless vs Provisioned

| Factor | Serverless | Provisioned |
|---|---|---|
| Management | No cluster sizing | Must choose node type and count |
| Scaling | Auto, instant | Manual or scheduled |
| Cost | Higher per-request; minimum cost applies | Lower per-request at sustained load |
| Best for | Spiky or unpredictable traffic | Steady, high-volume workloads |
| When to avoid | Predictable high-throughput (cost inefficient) | Variable traffic with difficult capacity planning |

### Pattern: Query Result Caching

```text
App → ElastiCache (cache hit?) → Return cached result
                ↓ (cache miss)
            RDS/Aurora → Store in cache → Return result
```

---

## Amazon Neptune (Graph)

Property graph (Apache TinkerPop/Gremlin) and RDF (SPARQL) support.

**Best for:**
- Social networks (friend-of-friend queries)
- Fraud detection (relationship patterns)
- Knowledge graphs
- Recommendation engines

**Not for:** General-purpose storage, or data where relationships are not the primary query driver. Graph query performance degrades when used as a simple document store.

---

## Amazon Timestream

Serverless time-series database optimized for IoT and operational analytics.

- Automatically tiers data from memory to magnetic storage
- SQL-like query language with time-series functions
- Built-in interpolation, smoothing, approximation

**Best for:** Sensor data, application metrics, DevOps telemetry.

**Not for:** General OLTP, relational joins across non-time-ordered data, or workloads requiring arbitrary update/delete patterns.

---

## Amazon QLDB (Quantum Ledger Database)

Immutable, cryptographically verifiable ledger.

- Append-only journal — records cannot be altered or deleted
- Built-in verification using SHA-256 hash chaining
- SQL-like query language (PartiQL)

**Best for:** Financial transactions, audit trails, supply chain provenance.

**Not for:** General-purpose databases, high-write-throughput workloads, or use cases that require deleting records (QLDB cannot delete committed data). Niche service — evaluate whether Aurora with audit logging suffices before adopting.

---

## AWS Database Migration Service (DMS)

Managed service for migrating databases to AWS with minimal downtime.

### Migration Types

**Homogeneous**: Same engine source and target (PostgreSQL → Aurora PostgreSQL).
- Prefer native tools (pg_dump/pg_restore, mysqldump) for one-time migrations — lower overhead, faster
- Use DMS when you need CDC (ongoing replication) or when native tools can't handle the network topology

**Heterogeneous**: Different engines (Oracle → PostgreSQL, SQL Server → Aurora MySQL).
- Always run AWS Schema Conversion Tool (SCT) first — it converts schema and up to 90% of stored procedures automatically
- Review SCT output manually for complex PL/SQL, triggers, and DB-specific functions (these often need rewriting)
- Use DMS for the data migration and CDC after schema conversion is validated

### DMS Patterns

```text
Pattern 1: Full Load
  Source DB ──DMS full load──→ Target DB (downtime required)

Pattern 2: Full Load + CDC (Change Data Capture)
  Source DB ──full load──→ Target DB
           ──CDC ongoing──→ Target DB (near-zero downtime)
  App cutover when replication lag ≈ 0

Pattern 3: Ongoing Replication
  Source DB ──CDC──→ Target DB (continuous sync, cross-region/hybrid)
```

### DMS Best Practices

- Size replication instances appropriately — heterogeneous migrations are CPU-intensive
- Enable multi-AZ for the replication instance in production
- Test with DMS validation tasks to compare row counts
- Pre-create indexes on the target after initial load completes (faster bulk insert)
- Monitor replication lag in CloudWatch

---

## Amazon DocumentDB

### What DocumentDB Actually Is

DocumentDB is **NOT MongoDB**. It is a proprietary AWS database that emulates a subset of the MongoDB API. The underlying storage engine is the same Aurora distributed storage layer — it is not a MongoDB fork.

> AWS marketing: "MongoDB-compatible"
> Reality: ~34% compatible with the full MongoDB API (as of 2025, per MongoDB's compatibility testing — verify against current AWS docs)

### Compatibility Level

DocumentDB supports MongoDB API versions 4.0, 5.0, and 8.0 (8.0 added November 2025). Common MongoDB drivers work because DocumentDB speaks the MongoDB wire protocol — but many features are absent or behave differently.

### Supported Features

- CRUD operations
- Aggregation pipeline (subset of operators)
- Basic indexes (single-field, compound, text)
- Change streams (with limitations)
- ACID transactions (single-node, with limits)
- TTL indexes

### Missing / Limited Features (2025)

| Feature | Status |
|---|---|
| `$facet` aggregation | Not supported |
| `$unionWith` | Not supported |
| `$lookup` with correlated subqueries | Only equality joins + uncorrelated subqueries |
| Sharding | Limited; elastic clusters use hash sharding only (hot partition risk) |
| Full MongoDB aggregation operator set | Subset only |
| Atlas Search / Vector Search | Not available |
| Client-side field-level encryption | Not supported |
| Time-series collections | Not supported |
| `$text` search | English-only; no array field indexing |
| Change stream event ordering guarantees | Different from MongoDB |
| Transactions | Single-writer; 1-min timeout; 32MB log cap |

### Change Streams Specifics

- Must be **explicitly enabled** per collection
- Events available for up to 7 days (default: 3 hours; configurable)
- Higher latency than MongoDB native change streams
- Long-running `updateMany`/`deleteMany` can stall change stream event writing
- `updateDescription` output differs: DocumentDB omits fields where the new value equals the old value

### Decision: DocumentDB vs MongoDB Atlas

| Dimension | DocumentDB | MongoDB Atlas |
|---|---|---|
| API compatibility | ~34% of MongoDB API | 100% (it IS MongoDB) |
| Sharding | Hash only (elastic) | Range, zone, geo |
| Transactions | Multi-document, limited | Full ACID, multi-shard |
| Change streams | Limited, must enable | Full, real-time |
| Write scaling | Single primary only | Multiple via sharding |
| Multi-cloud | AWS only | AWS, GCP, Azure |
| Atlas Search | No | Yes (Lucene-based) |
| Vector search | No | Yes |
| Cost (production) | ~26.7% more expensive | Baseline |
| Backup restore | Must create new cluster | Restore to existing cluster |
| Multi-AZ | Extra cost per replica | Included |

**Use DocumentDB when:**
- Already on AWS, need a managed document database, and the required feature set is confirmed to be within DocumentDB's supported subset
- Primarily simple CRUD with basic aggregation
- Cannot use MongoDB Atlas for policy/compliance reasons

**Use MongoDB Atlas instead when:**
- Needing full MongoDB features ($facet, $unionWith, Atlas Search, Vector Search)
- Needing multi-region writes or proper horizontal sharding
- Using complex aggregation pipelines
- Needing client-side field-level encryption

### DocumentDB Connection Patterns

Use standard MongoDB drivers — DocumentDB supports the MongoDB wire protocol.

```python
# Python (pymongo)
import pymongo
client = pymongo.MongoClient(
    "mongodb://user:pass@docdb-cluster.cluster-xxx.us-east-1.docdb.amazonaws.com:27017/",
    tls=True,
    tlsCAFile="global-bundle.pem",  # Required: AWS CA bundle
    retryWrites=False               # Required: DocumentDB does not support retryable writes
)
```

```javascript
// Node.js (mongoose)
mongoose.connect(uri, {
  tls: true,
  tlsCAFile: './global-bundle.pem',
  retryWrites: false,               // Must disable
  directConnection: false
});
```

**Critical connection requirements:**
- `retryWrites: false` — DocumentDB does not support retryable writes
- TLS required; must include AWS CA bundle (`global-bundle.pem`) in deployment
- Connect to cluster endpoint for writes; reader endpoint for reads

---

## CockroachDB

### Architecture: Distributed SQL

CockroachDB implements SQL semantics on top of a distributed key-value storage engine inspired by Google Spanner. Three layers:

1. **SQL Layer**: Parses and plans SQL queries (reuses PostgreSQL parser)
2. **Distribution Layer**: Breaks data into 512 MB ranges; routes queries to correct nodes
3. **Storage Layer**: Each range is a Raft consensus group; writes require majority quorum

**Key architectural properties:**
- No master node — every node is equal (removes single points of failure)
- Leaseholder model: one replica per range serves reads without consensus overhead (low latency)
- Serializable isolation by default (strongest SQL isolation level)
- HLC (Hybrid Logical Clocks) for global time ordering without GPS hardware

### Consistency Model

CockroachDB defaults to **serializable isolation** — the strongest SQL isolation level. This means:
- No dirty reads, no phantom reads, no write skew
- Slightly higher latency than read-committed (default in PostgreSQL)
- Applications must handle transaction retries (`RETRY_SERIALIZABLE` errors under contention)
- Can be relaxed to `READ COMMITTED` per transaction if needed — this reduces contention and retry errors but permits non-repeatable reads (the same row may return different values within the same transaction if another transaction commits between reads). Only use where that trade-off is acceptable.

### PostgreSQL Compatibility

CockroachDB reuses PostgreSQL's syntax parser (with extensions). Most PostgreSQL SQL works, but not all.

**Third-party PostgreSQL Compatibility Index score: ~40% (as of 2025)** — meaning many PostgreSQL extensions, features, and behaviors differ. (Source: [Bytebase 2025 analysis](https://www.bytebase.com/blog/cockroachdb-vs-postgres/))

#### What Works

- Standard DDL (CREATE TABLE, ALTER TABLE, DROP)
- Most DML (SELECT, INSERT, UPDATE, DELETE)
- Joins (including complex joins)
- CTEs, window functions
- Standard aggregate functions
- JSON/JSONB operators
- Most PostgreSQL data types

#### Known Incompatibilities / Differences

| Feature | CockroachDB Behavior |
|---|---|
| `CREATE DOMAIN` | Not supported |
| PostgreSQL range types | Not supported |
| Unary `~` precedence | High precedence (same as `-`) vs PostgreSQL low precedence |
| `ENUM` types | Supported but limited mutability |
| Sequences | Supported; may have gaps under concurrent access |
| `RETURNING` with `ON CONFLICT` | Slightly different semantics |
| `pg_catalog` tables | Partial support |
| Extensions (PostGIS, pg_trgm, etc.) | Not supported |
| Stored procedures | Limited support |
| Triggers | Not supported |

#### CockroachDB-Specific SQL Extensions

```sql
-- Multi-region locality
ALTER TABLE orders SET LOCALITY REGIONAL BY ROW;

-- Table survivability
ALTER DATABASE mydb SURVIVE REGION FAILURE;

-- Follower reads (stale, low-latency)
SELECT * FROM orders AS OF SYSTEM TIME follower_read_timestamp();

-- Show ranges
SHOW RANGES FROM TABLE orders;
```

### Multi-Region Deployment Patterns

```text
Pattern 1: Single Region (basic HA)
  3 nodes across 3 AZs
  Survives 1 AZ failure (2/3 quorum maintained)

Pattern 2: Multi-Region (regional survival)
  5 nodes: 3 in primary region + 1 each in 2 other regions
  Survives failure of either secondary region (3/5 quorum maintained)
  Does NOT survive failure of the primary region (only 2/5 nodes remain)
  For primary-region failure survival: use 3 nodes per region (9 total)

Pattern 3: Multi-Region with REGIONAL BY ROW
  Rows pinned to specific regions via crdb_region column
  Low-latency reads/writes for region-local data
  Global tables replicated everywhere (config data, reference data)
```

#### Table Locality Types

| Locality | Behavior | Use Case |
|---|---|---|
| `REGIONAL BY TABLE` | Entire table in one region | Single-region apps |
| `REGIONAL BY ROW` | Each row pinned to a region | Geo-partitioned user data |
| `GLOBAL` | All replicas in all regions (optimized for reads) | Config/reference tables |

### Serverless vs Dedicated (2025)

**CockroachDB Serverless:**
- Disaggregated architecture: SQL layer and KV storage run as separate processes
- True scale-to-zero; sub-second scale-from-zero
- Multi-region serverless now available
- Shared infrastructure (multi-tenant)
- Better for: dev/test, bursty workloads, small customer tenants
- Not ideal for: high-security compliance requirements, sustained heavy workloads

**CockroachDB Dedicated (Advanced plan):**
- Dedicated hardware per cluster
- Stricter SLAs, compliance certifications (SOC 2, HIPAA)
- Customizable; better performance predictability
- Available on AWS, GCP, Azure
- Validated support for 300-node clusters, up to 1 PB per cluster (v25.4.4+)

### Licensing Note (2025)

CockroachDB moved to a **proprietary license** (Business Source License / CockroachDB Community License). The core is no longer Apache 2.0 open source. Evaluate licensing implications before committing to CockroachDB for production workloads.

---

## IndexedDB

### What It Is

IndexedDB is a browser-native structured storage API — a transactional, indexed, key-value object store in the browser. It is:
- Asynchronous (Promise/event-based)
- Transactional (ACID within the browser)
- Supports indexes for efficient querying
- Works on objects (not just strings like localStorage)

### Storage Architecture

- **Database**: Top-level container (versioned)
- **Object Store**: Like a table; stores JavaScript objects
- **Index**: Secondary lookup path on an object property
- **Transaction**: Groups reads/writes; auto-committed on completion

### Storage Limits (Browser-Specific, 2025–2026)

| Browser | Quota |
|---|---|
| Chrome / Edge | Up to 60% of total disk space per origin; browser ceiling at 80% across all origins |
| Firefox | Up to 50% of free disk space; ~2 GB cap per eTLD+1 group |
| Safari | Starts ~1 GB; prompts user for more in 200 MB increments; **evicts after 7 days of inactivity** (ITP) |
| Incognito / Private | Chrome: ~5% cap |

**Compare to localStorage: 5–10 MB max.** IndexedDB is suitable for large datasets (images, documents, thousands of records).

### localStorage vs IndexedDB vs Cache API

| Feature | localStorage | IndexedDB | Cache API |
|---|---|---|---|
| Storage limit | 5–10 MB | Up to 50–60% of disk | Up to 50–60% of disk |
| Data type | Strings only | Any structured data | HTTP Request/Response |
| Async | No (synchronous — blocks main thread) | Yes | Yes |
| Indexes | No | Yes | No |
| Transactions | No | Yes | No |
| Use case | Small config, flags | Structured app data | Network response caching |
| Service Worker | No | Yes | Yes |

**Rule of thumb:**
- `localStorage`: Tiny simple key-value (user prefs, auth tokens) — max 5–10 MB, synchronous
- `IndexedDB`: App data, offline records, large structured data
- `Cache API`: Network responses, assets (use with Workbox)

### Libraries

#### Dexie.js (Recommended Default)

Best overall IndexedDB library as of 2026. Provides:
- TypeScript-first, typed tables
- Declarative schema with versioning/migration
- `useLiveQuery` hook for React (live reactive queries)
- Compound queries with `.where().and().filter()`
- Built-in transaction support

```typescript
// Schema definition and versioning
const db = new Dexie('MyAppDB');
db.version(1).stores({
  friends: '++id, name, age',           // ++ = auto-increment PK
  orders: '++id, userId, [status+date]', // compound index
});

// Version migration
db.version(2).stores({
  friends: '++id, name, age, email',    // added email field
}).upgrade(tx => {
  return tx.table('friends').toCollection().modify(friend => {
    friend.email = '';
  });
});

// Query
const activeOrders = await db.orders
  .where('status').equals('active')
  .and(o => o.userId === currentUserId)
  .sortBy('date');

// React live query
const friends = useLiveQuery(() => db.friends.toArray());
```

#### idb (Google)

Thin Promise wrapper around raw IndexedDB. Lower-level than Dexie. Good when you want minimal abstraction.

```typescript
import { openDB } from 'idb';
const db = await openDB('MyDB', 1, {
  upgrade(db) {
    db.createObjectStore('store', { keyPath: 'id' });
  }
});
```

#### localForage

Simplest API; falls back from IndexedDB to localStorage for compatibility with very old browsers. Good for simple key-value use cases. Not suitable for complex queries or indexed lookups.

> **Note**: WebSQL (previously in the fallback chain) is removed in all modern browsers as of 2022. localForage in 2025 effectively falls back only to localStorage.

### Transactions

IndexedDB transactions are **auto-committed** when the last request completes and no new requests are pending.

```javascript
// Transaction spanning multiple stores
const tx = db.transaction(['users', 'orders'], 'readwrite');
const userStore = tx.objectStore('users');
const orderStore = tx.objectStore('orders');

userStore.put({ id: 1, name: 'Alice' });
orderStore.put({ id: 100, userId: 1, total: 50 });

await tx.done; // Wait for both to commit
```

**Key rules:**
- Keep transactions short — they auto-close when idle
- Don't use `await` on unrelated work inside a transaction (the idle gap closes the transaction)
- Use Dexie's `.transaction()` wrapper to handle this safely

### Versioning and Migrations

```javascript
// Raw IndexedDB
const request = indexedDB.open('MyDB', 2); // Bump version number

request.onupgradeneeded = (event) => {
  const db = event.target.result;
  const oldVersion = event.oldVersion;

  if (oldVersion < 1) {
    // Create initial schema
    const store = db.createObjectStore('users', { keyPath: 'id' });
    store.createIndex('email', 'email', { unique: true });
  }

  if (oldVersion < 2) {
    // Add new index in v2
    const store = event.target.transaction.objectStore('users');
    store.createIndex('name', 'name');
  }
};
```

### Offline-First Pattern with Dexie + Service Worker

```text
User Action
    ↓
Write to IndexedDB (immediate, offline-safe)
    ↓
useLiveQuery re-renders component (automatic)
    ↓
Background sync (service worker or useEffect)
    ↓
POST to API when online
    ↓
Update IndexedDB with server response
    ↓
useLiveQuery re-renders with canonical data
```

Best practice stack: **Dexie.js + Workbox Background Sync** for offline-first PWAs.

### Safari / ITP Caveat

Safari's Intelligent Tracking Prevention (ITP) **evicts IndexedDB data after 7 days of no user interaction**. For PWAs that need persistent local data on iOS/Safari:
- Request persistent storage: `navigator.storage.persist()`
- Users must add the PWA to Home Screen for full persistence
- Never rely on IndexedDB alone as a durable store on Safari — always implement server sync

---

## Choosing the Right Database

### Decision Tree

```text
Is your data relational with complex queries/joins?
  ├── Yes + production + need scale → Aurora (PostgreSQL or MySQL)
  ├── Yes + lower traffic or non-Aurora engine → RDS
  └── Yes + active-active multi-region → Aurora DSQL

Is your access pattern primarily key-based at high volume?
  └── DynamoDB

Do you need MongoDB API compatibility?
  ├── Need full MongoDB features → MongoDB Atlas
  └── Basic CRUD/aggregation on AWS, confirmed within DocumentDB limits → DocumentDB

Do you need distributed SQL with horizontal write scale?
  └── CockroachDB

Is it browser-side structured storage?
  └── IndexedDB (Dexie.js)

Do you need sub-millisecond caching?
  └── ElastiCache (Valkey/Redis)

Is it relationship/graph data?
  └── Neptune

Is it time-series sensor/metric data?
  └── Timestream

Do you need immutable audit log?
  └── QLDB
```

### RDS vs Aurora vs DynamoDB: Quick Reference

| Factor | RDS | Aurora | DynamoDB |
|---|---|---|---|
| Query model | SQL, flexible | SQL, flexible | Key-value, limited |
| Max storage | 64 TB | 128 TB (auto) | Unlimited |
| Read replicas | Up to 5 | Up to 15 | N/A (GSIs serve as alternate read paths) |
| Write scale | Single primary | Single primary | Unlimited (partitioned) |
| Latency | ms | ms | Single-digit ms |
| Cost model | Instance-based | Instance + storage | Per-request or provisioned |
| Best for | Dev/staging, SQL, Oracle/SQL Server | Production SQL at scale | Known access patterns, extreme throughput |
| Dangerous mistake | Underprovisioning | Over-engineering small apps | Choosing without confirming all access patterns |

---

## Common Patterns

> Cross-references: most patterns link back to full code in the relevant service section above. The DMS Full Load + CDC pattern and the CockroachDB Follower Reads pattern are fully specified here (not duplicated above).

### Pattern: Connection Pooling for Lambda + RDS/Aurora

```text
Lambda → RDS Proxy → Aurora
```
- RDS Proxy maintains a persistent connection pool; Lambda connects to Proxy, not the DB directly
- Prevents connection exhaustion when Lambda scales to thousands of concurrent instances
- See also: [RDS Proxy](#rds-proxy) section

### Pattern: DynamoDB + Lambda Event-Driven Pipeline

```text
Write to DynamoDB → DynamoDB Stream → Lambda → SNS/SQS/EventBridge
```
- Near-real-time item-level change processing without polling
- Lambda trigger fires for INSERT, MODIFY, REMOVE events

### Pattern: Caching Layer with ElastiCache

```text
App → Check ElastiCache → Hit: return cached value
                        → Miss: query RDS → cache result → return
```
- Set appropriate TTLs; invalidate on write for consistency
- Write-through: update cache synchronously on every write (consistent, higher write latency)
- Write-behind: update cache immediately, DB asynchronously (faster writes, risk of data loss)

### Pattern: DocumentDB Connection (retry disabled)

See the **Amazon DocumentDB → DocumentDB Connection Patterns** section above for full code examples.

Key requirement: always set `retryWrites: false` and include the AWS CA bundle (`global-bundle.pem`).

### Pattern: CockroachDB Follower Reads

```sql
-- Read slightly stale data from nearest replica (low latency)
SELECT * FROM orders
AS OF SYSTEM TIME follower_read_timestamp()
WHERE user_id = $1;
```
Use when stale reads are acceptable and cross-region read latency is a bottleneck.

### Pattern: IndexedDB as Single Source of Truth (Offline-First)

See [Offline-First Pattern with Dexie + Service Worker](#offline-first-pattern-with-dexie--service-worker) for full flow.

Key principle: write to IndexedDB first (always succeeds offline), sync to API in background, `useLiveQuery` re-renders automatically on both sides.

### Pattern: DMS Full Load + CDC Migration

```text
1. Spin up DMS replication instance (size for CPU-intensive heterogeneous migrations)
2. Run SCT (Schema Conversion Tool) for heterogeneous migrations
3. Start DMS full-load task → migrates existing data
4. Enable CDC on DMS task → captures ongoing changes
5. Monitor replication lag in CloudWatch
6. When lag ≈ 0 seconds, cut over app to target DB
7. Stop CDC; decommission source
```

---

## Anti-Patterns

### AWS Databases

- **DynamoDB for reporting**: Running multi-attribute scans/analytics on DynamoDB causes full table scans and high cost. Use Redshift or Aurora for analytics.
- **Skipping RDS Proxy with Lambda**: Direct Lambda → RDS connections exhaust the database's max connections at scale.
- **Hot partitions in DynamoDB**: Using a low-cardinality value (e.g., `status` = "active") as a partition key routes most traffic to one partition. Use composite keys.
- **Aurora Serverless v2 for predictable high-throughput**: Serverless auto-scaling overhead costs more than provisioned at sustained baseline load. Use provisioned writer.
- **Forgetting Aurora storage behavior**: Aurora's shared storage layer doesn't reclaim space immediately after large DELETEs (background vacuum). Monitor `FreeLocalStorage` separately.

### DocumentDB

- **Assuming MongoDB feature parity**: Testing on local MongoDB then deploying to DocumentDB. Always test against DocumentDB early. Key gaps: `$facet`, `$unionWith`, sharding, Atlas Search.
- **Using retryWrites: true**: DocumentDB does not support retryable writes; this causes runtime errors. Always set `retryWrites: false`.
- **Relying on change streams for real-time pipelines**: DocumentDB change streams have higher latency and limited retention. For critical event pipelines, use SQS/EventBridge instead.
- **Expecting horizontal write scaling**: DocumentDB has a single write primary. If you need sharded writes, use MongoDB Atlas.

### CockroachDB

- **Ignoring serializable retry errors**: CockroachDB may surface `RETRY_SERIALIZABLE` under contention. Application code must implement retry loops; unhandled retries surface as transaction failures to users.
- **Expecting full PostgreSQL extension support**: PostGIS, pg_trgm, and other popular extensions don't exist in CockroachDB. Validate extension dependencies before migrating.
- **Assuming all PostgreSQL tooling works**: ORM compatibility, migration tools (Flyway, Liquibase), and admin tools may have partial or broken support. Test thoroughly.
- **Using stored procedures / triggers**: Not supported. Application-level logic required.

### IndexedDB

- **Running multiple `await` calls inside a transaction without active requests**: IndexedDB auto-closes idle transactions. Use Dexie's `.transaction()` wrapper to avoid this.
- **Using synchronous localStorage as a drop-in**: localStorage blocks the main thread; IndexedDB is async. The APIs are not interchangeable.
- **Trusting Safari persistence**: ITP evicts IndexedDB data after 7 days of inactivity unless persistent storage is granted. Always implement server sync.
- **Storing large blobs without compression**: Storing uncompressed images or large files bloats IndexedDB quickly. Compress or use Cache API for binary assets.
- **No migration strategy**: Opening IndexedDB without versioned `onupgradeneeded` logic leads to schema corruption on updates. Always increment version and guard each migration block with `if (oldVersion < N)`.

---

## Troubleshooting

Per-service symptom/cause/fix tables for DynamoDB, Aurora/RDS, DocumentDB, CockroachDB, and IndexedDB: see `references/troubleshooting.md`.

---

## References

AWS, DocumentDB, CockroachDB, and IndexedDB reference links: see `references/references.md`.
