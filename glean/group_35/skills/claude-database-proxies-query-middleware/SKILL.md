---
name: database-proxies-query-middleware
description: >-
  Database proxies & query-optimization middleware: the data-access control plane that pools,
  routes, caches, rewrites, and firewalls queries. SQL/engine sibling of llm-ai-gateways.
  TRIGGER: choosing/designing a DB proxy (ProxySQL, MariaDB MaxScale, Vitess/vtgate, Pgpool-II,
  PgBouncer, PgCat, AWS RDS Proxy, Heimdall); connection pooling/multiplexing; read/write
  splitting and the read-after-write staleness hazard (causal_reads, GTID); proxy result caching;
  SQL firewall; proxy-based sharding; failover acceleration and proxy-tier HA; credential
  brokering; the proxy-vs-driver-vs-app decision; RDS Proxy pinning; PgBouncer transaction-mode
  pooling; why MongoDB uses mongos natively.
  SKIP: LLM/AI request proxying -> llm-ai-gateways; driver CMAP/SDAM -> mongodb-driver-internals;
  compute-storage disaggregation -> mongodb-operations-expert; cross-engine federation -> adjacent
  only; generic API gateways -> software-engineering-patterns; learned/ML query optimization ->
  da-analytical-methods.
category: developer
version: "1.0.0"
updated: "2026-06-18"
whenToUse:
  - "deciding whether to put a database proxy (ProxySQL, MaxScale, PgBouncer, PgCat, Pgpool-II, RDS Proxy, Vitess) between the app and the database, and which one"
  - "designing connection pooling & multiplexing, read/write splitting, query routing, or proxy-based sharding"
  - "diagnosing the read-after-write / replication-lag staleness hazard when reads are routed to replicas"
  - "configuring proxy-side query result caching, query rewriting, or a SQL firewall / query allow-list"
  - "planning proxy-tier high availability, failover acceleration, or credential brokering (IAM / Secrets Manager)"
  - "weighing proxy/middleware vs driver-level features vs application-layer handling, and sizing a connection pool"
  - "explaining why MongoDB uses mongos natively and generally does not adopt third-party query proxies"
whenNotToUse:
  - "LLM/AI request proxying, model routing, token budgets, or AI gateways — use llm-ai-gateways"
  - "MongoDB driver connection-pool (CMAP) / SDAM internals — use mongodb-driver-internals"
  - "compute–storage disaggregation (Aurora / Neon / PolarDB / object-storage OLTP) — use mongodb-operations-expert"
  - "cross-engine query federation (Trino / Presto / Starburst) — adjacent only, not this skill"
  - "generic API gateway / BFF design — use software-engineering-patterns"
  - "learned / ML-based query optimization (learned cardinality, Bao, OtterTune) — use da-analytical-methods"
keywords:
  - database proxy
  - query optimization middleware
  - connection pooling
  - read write splitting
  - ProxySQL
  - PgBouncer
  - Vitess vtgate
  - RDS Proxy
  - SQL firewall
  - query routing
tags:
  - developer
  - database
  - middleware
  - infrastructure
  - architecture
related_skills:
  - llm-ai-gateways
  - mongodb-driver-internals
  - software-engineering-patterns
  - da-analytical-methods
---

# Database Proxies & Query-Optimization Middleware

> **Scope.** The **database-proxy / data-access-middleware control plane**: a layer that sits
> between the application and the database server(s) and pools, routes, caches, rewrites,
> firewalls, and load-balances queries, transparently, so application code keeps speaking the
> native database protocol. This is the **SQL/database-engine sibling** of the LLM world's
> request gateways; for AI/model request proxying see **`llm-ai-gateways`** ("AI Gateways & LLM
> Proxy Infrastructure"), which is the analog pattern for LLM traffic and owns that surface.
>
> `verified-as-of: 2026-06-18`. Version-pinned and vendor-landscape claims below carry inline
> footnotes to primary sources; re-verify volatile facts (product versions, pricing, support
> dates) before relying on them.

## When this skill applies

Reach for this skill when the question is *"should there be a proxy/middleware tier in front of
my database, which one, and how do I run it?"*: connection pooling, read/write splitting,
query caching/rewriting, a SQL firewall, proxy-based sharding, failover acceleration, or
credential brokering. It covers the **MySQL/PostgreSQL ecosystem** (where third-party proxies
are common) and explains the **MongoDB** position (native `mongos` router; the driver already
does pooling/routing, so third-party proxies are rarely used).

The deep per-product material lives in `references/` (load on demand):

| Reference | Covers |
|---|---|
| `references/proxysql-maxscale.md` | ProxySQL (3-layer config, query rules, multiplexing, firewall) and MariaDB MaxScale (routers/filters, `causal_reads`, `dbfwfilter`, `maxctrl`) |
| `references/postgres-poolers-vitess.md` | PgBouncer (3 modes + what breaks), PgCat, Pgpool-II, and Vitess/vtgate as the sharding-proxy exemplar |
| `references/cloud-proxies-decision-security.md` | AWS RDS Proxy (pinning, failover, IAM), Heimdall Data, the proxy-vs-driver-vs-app decision framework, SQL firewall + credential brokering, anti-patterns |
| `references/mongodb-angle.md` | `mongos` native router; why MongoDB rarely uses third-party query proxies (SDAM/CMAP/read-preference in the driver); BI Connector & Atlas Data Federation as adjacent query surfaces |

---

## Core concept: the data-access control plane

A database proxy is a **protocol-aware intermediary**. The application connects to the proxy as
if it were the database; the proxy parses (or at least inspects) the wire protocol, decides what
to do with each connection and query, and forwards to one or more backends. Because it speaks the
native protocol (MySQL, PostgreSQL, etc.), it is **transparent to application code** (no driver
change beyond repointing the connection string).[^proxysql-arch][^rdsproxy-overview]

The capabilities a proxy can provide (any given product implements a subset):

1. **Connection pooling & multiplexing**: keep a small set of warm backend connections and
   share them across many client connections, decoupling client-connection count from
   backend-connection count.
2. **Read/write splitting & load balancing**: route writes to the primary and reads across
   replicas, spreading read load.
3. **Query routing**: direct queries to specific backends by rule, shard key, or content.
4. **Failover / HA**: detect a failed backend, re-route, and (for some) preserve client
   connections across a database failover.
5. **Query result caching**: cache result sets at the proxy with a TTL and/or invalidation.
6. **Query rewriting & SQL firewall**: modify queries in flight (add hints, fix ORM SQL) and
   block queries that don't match an allow-list (defense-in-depth against injection).
7. **Sharding / horizontal scale-out**: present one logical endpoint over many shards.
8. **Observability**: a central vantage point that sees every query (per-digest stats, latency,
   pool state).
9. **Security / credential brokering**: the app authenticates to the proxy; the proxy holds the
   real database credentials (or uses IAM), so app code never carries DB passwords.

> **The sibling framing.** Almost every capability above has a one-to-one analog in an LLM/AI
> gateway: connection pooling ↔ provider connection reuse; read/write split & load balancing ↔
> model routing & fallback chains; query result caching ↔ exact/semantic prompt caching; SQL
> firewall ↔ guardrail/PII enforcement; credential brokering ↔ virtual keys / BYOK key-vault;
> per-query observability ↔ per-request token/cost tracing. The two patterns solve the same
> "control plane in front of a stateful backend" problem in different domains. Defer to
> **`llm-ai-gateways`** for the LLM side; do not duplicate it here.

---

## The product landscape (where each fits)

| Product | Ecosystem | Primary role | Notable capability |
|---|---|---|---|
| **PgBouncer** | PostgreSQL | **Pure connection pooler** | Lightest possible pooling; no routing.[^pgb-features] |
| **PgCat** | PostgreSQL | Pooler + routing | Rust; adds load balancing, read/write routing, sharding, failover.[^pgcat] |
| **Pgpool-II** | PostgreSQL | Pooler + LB + HA suite | Load balancing, failover, watchdog HA, query cache, replication.[^pgpool] |
| **ProxySQL** | MySQL (+ PgSQL) | Query-aware proxy | Multiplexing, rule-based routing, result cache, SQL firewall.[^proxysql-arch] |
| **MariaDB MaxScale** | MySQL/MariaDB | Query-aware proxy | Auto read/write split via query classifier, `causal_reads`, sharding, binlog relay.[^maxscale-arch] |
| **Vitess (vtgate)** | MySQL | **Sharding proxy** | Presents one MySQL endpoint over many sharded MySQLs; the horizontal-scaling exemplar.[^vitess-arch] |
| **AWS RDS Proxy** | RDS/Aurora (MySQL/PG/MariaDB/SQL Server) | Managed pooler | Multiplexing for serverless, failover acceleration, IAM + Secrets Manager.[^rdsproxy-overview] |
| **Heimdall Data** | Multi-DB (commercial) | Query-optimization proxy | Auto query caching with auto-invalidation, read/write split (vendor-heavy claims).[^heimdall] |
| **MongoDB `mongos`** | MongoDB | **Native** query router | Shard routing/merge; first-party, not a third-party bolt-on.[^mongos] |

**Two distinct centers of gravity:**

- **Pure poolers** (PgBouncer) just multiplex connections. They do *not* parse queries or route.
- **Query-aware proxies** (ProxySQL, MaxScale, Vitess, PgCat, Pgpool-II) inspect SQL and can
  route, split, cache, rewrite, firewall, and shard, at the cost of more CPU (they are
  Layer-7, content-aware) and more operational surface.[^percona-compare]

**Vitess is the sharding-proxy exemplar:** `vtgate` is a stateless proxy speaking the MySQL
protocol that routes to the right shard(s) using a `VSchema`/Vindex, scatter-gathers when no
shard key is present, and merges results, moving sharding logic out of the application. It came
out of YouTube's MySQL scaling and is a graduated CNCF project.[^vitess-arch][^vitess-history]

---

## Capability deep-dives (summary; details in references)

### Connection pooling & multiplexing
The foundational job. Databases are expensive per connection: PostgreSQL `fork()`s an OS
process per connection; MySQL is thread-per-connection. Idle connections still cost memory and
add contention that can scale **super-linearly**.[^citus-conns][^pg-process] A pooler keeps a
small warm set (often tens) and serves thousands of clients from it.

**Multiplexing** (ProxySQL, RDS Proxy) goes further: a single backend connection is reused across
*different* client connections when it is safe to do so. It becomes unsafe — and the proxy
**pins** the client to one backend (reverting to 1:1) — whenever **session state** would leak:
open transactions, table/named locks, session variables (`SET`), temp tables, prepared
statements, cursors, `LISTEN`. Heavy pinning silently defeats the benefit; minimizing it is a
core tuning task.[^proxysql-mux][^rdsproxy-pinning]

**Pool sizing is counter-intuitive:** a *small* pool usually beats a large one. The PostgreSQL
community formula is roughly `(core_count × 2) + effective_spindle_count` → typically ~10–20
connections; throughput climbs to a saturation "knee" then **falls** under contention. The fix
for a saturated pool is faster queries, not more connections.[^pg-poolsize][^hikari-poolsize]

### Read/write splitting & the read-after-write hazard
Routing reads to replicas multiplies read capacity, but asynchronous replicas lag, so a read
issued right after a write can hit a replica that hasn't applied that write yet → **stale /
non-read-your-writes** results. This is the central correctness hazard of read/write splitting.
A coarse "max replication lag" filter does **not** guarantee causality.[^maxscale-causal]

Mitigations:
- **MaxScale `causal_reads`** — uses GTID tracking + `MASTER_GTID_WAIT()` so a read waits until
  the chosen replica has applied the writing session's latest GTID (modes `local`/`global`/
  `fast`/`universal`, trading latency for scope).[^maxscale-causal]
- **MaxScale CCR filter** — after a write, force that connection's reads to the primary for N
  seconds/statements.[^maxscale-causal]
- **ProxySQL** — GTID-aware routing via `gtid_from_hostgroup` (route a read only to a replica
  that reached the needed GTID); operationally, keep reads-after-writes and transactions on the
  primary. ProxySQL has no single named "causal read" mode; naive `^SELECT → replica` rules are
  explicitly warned against because they break transactions and consistency.[^proxysql-rwsplit]

### Query result caching
Some proxies cache result sets. The hazard is **invalidation**: ProxySQL's cache is **TTL-only**
(no invalidation, no LRU, breaks on prepared statements) — stale until TTL expiry.[^proxysql-cache]
MaxScale's `cache` filter has soft/hard TTL plus optional change-based invalidation and pluggable
storage (in-memory/Redis/Memcached).[^maxscale-cache] Heimdall markets automatic table-aware
invalidation, but its strong-consistency framing is vendor/sponsored-sourced; treat with
skepticism and note that invalidation at scale is a genuinely hard problem.[^heimdall]

### Query rewriting & SQL firewall
**Rewriting:** match a query by regex/digest and rewrite it before execution (add index hints,
fix ORM-generated SQL, inject `SQL_NO_CACHE`) without touching app code.[^proxysql-rwsplit]

**SQL firewall** (defense-in-depth, *not* a replacement for parameterized queries): record the
finite set of query **digests** (normalized fingerprints with literals stripped) the app
legitimately emits, then block anything else. An injected query produces a new digest and is
rejected. Standard lifecycle: **train → detect (log only) → protect (enforce)**.
- **ProxySQL** `mysql_firewall_whitelist_*` (modes OFF / DETECTING / PROTECTING; per-user digest
  allow-lists; a bundled SQLi-fingerprint engine that is imperfect / false-positive-prone).[^proxysql-firewall]
- **MaxScale** `dbfwfilter` (`allow`/`block`/`ignore`; column/where-clause/function rules; learn/
  supervise/enforce modes).[^maxscale-firewall]
- **Caveat:** parser-differential bypasses exist (wrap payloads so the proxy parser and the
  backend disagree), and enabling enforce before fully training the allow-list **blocks all
  legitimate traffic**. OWASP treats allow-listing as secondary to parameterized statements.[^owasp-sqli][^maxscale-firewall]

### Failover, HA, and the proxy as a SPOF
A proxy can **accelerate failover**: AWS RDS Proxy holds client connections open during a DB
failover, reconnects to the new writer behind the scenes, and bypasses DNS/CNAME propagation —
AWS cites up to a **66%** reduction in failover time (magnitude is scenario-dependent; AWS still
warns latency rises and in-flight transactions may need retry).[^rdsproxy-failover]

But **a single proxy is a single point of failure** — it moves the SPOF from the database to the
proxy tier.[^spof] Run it HA:
- **Sidecar / co-located** (a proxy per app host, app connects to `127.0.0.1`): zero extra
  network hop, no shared SPOF — but N instances drift in config and can dilute the
  connection-reduction benefit as they autoscale.[^proxysql-arch]
- **Centralized cluster behind an L4 load balancer**: shared endpoint, config synced (ProxySQL
  Cluster); small added hop.[^proxysql-arch]
- **Keepalived + VRRP floating VIP** (active/passive): the most-cited HA pattern for ProxySQL and
  MaxScale; with MaxScale, only the active node may manipulate the cluster (`maxctrl alter
  maxscale passive true` on the others).[^maxscale-arch]

### Credential brokering
The app authenticates to the proxy (e.g., with a short-lived IAM token); the proxy authenticates
to the database using credentials it holds (RDS Proxy fetches the DB password from **Secrets
Manager**, or uses end-to-end IAM). App config carries no DB password; rotation is centralized
and access is auditable. The trade-off is a **concentrated trust point** — compromise of the
proxy's role/secret access approaches DB-credential compromise; IAM auth requires TLS.[^rdsproxy-iam]

---

## Decision framework: proxy vs driver-level vs application-layer

These are **layers with distinct jobs, not substitutes.** The guiding rule: *add the fewest tiers
that make your actual problem go away.* Most teams need one (the in-driver pool); many need two
(+ a proxy); some need three.[^decision]

**Step 0 — exhaust cheaper options first.**
- **Right-size the in-driver pool** (HikariCP, the driver's own pool) and **fix slow queries**
  *before* adding any tier. A bigger pool usually loses throughput. Raising `max_connections` to
  mask the problem makes the eventual collapse worse.[^pg-poolsize][^hikari-poolsize]

**Stay driver-only when:** few app instances; total connections (pool × processes) well under the
DB limit; no serverless/autoscaling; short queries. Driver-level read/write split (e.g. MySQL
Connector/J `replication://`, the AWS JDBC driver) is fine when routing is **simple and static** —
but the policy lives in app code (no central audit, updated per-codebase).[^decision]

**Handle it in the application when:** the need is read offload with tolerable staleness →
**cache-aside** (Redis/Memcached) is the standard answer (it tolerates cache unavailability by
falling through to the DB). Not for balances/read-your-writes data.[^decision]

**Add a PROXY/MIDDLEWARE tier when any trigger fires:**
- **Serverless/Lambda** — near-mandatory; a per-invocation function has no persistent process to
  pool, so thousands of concurrent invocations exhaust the DB. An external pool (RDS Proxy) is the
  correct fix.[^rdsproxy-overview][^decision]
- **Horizontal scale-out past the connection budget** / "too many clients" / connection storms.
- **Code-free centralized** read/write split, sharding, query rules, mirroring, or caching →
  query-aware proxy (ProxySQL / Vitess / MaxScale).
- **Tight failover RTO + IAM/rotation** as first-class → RDS Proxy.

**What the proxy hop COSTS (weigh against the benefit):**
- **Latency:** ~1–5 ms/query in-VPC (sidecar/loopback ≈ 0; cross-AZ 1–2 ms; separate host
  10–15 ms). Negligible for 10–50 ms OLTP queries; **pure overhead on long analytical queries** —
  a reporting query that already takes seconds does not benefit from pooling; use a direct
  connection.[^percona-overhead][^decision]
- **An extra HA tier** to size, scale, patch, monitor (managed proxies shift but don't remove
  this — e.g., RDS Proxy IP-subnet exhaustion can block patching; per-vCPU billing even when
  idle).[^rdsproxy-overview][^decision]
- **A SPOF if not HA.**[^spof]
- **Semantic gotchas** — transaction-mode pooling (the mode that gives the payoff) breaks session
  state (see anti-patterns).

**Managed vs self-hosted:** RDS Proxy for low-ops basic pooling/failover/IAM on RDS/Aurora;
self-hosted ProxySQL/PgBouncer/PgCat/MaxScale for fine tuning, query rules, mirroring, or non-RDS
databases. RDS Proxy bills per-vCPU regardless of throughput — "2× cost for a benefit you don't
need" on a small fleet, but "the only thing keeping the database alive" for a Lambda fleet.[^decision]

> **Learned/ML query optimization is a different axis.** Steering *which plan the database picks*
> via learned cardinality estimation, Bao, Neo, or knob-tuning (OtterTune) is **not** a proxy
> concern — see `da-analytical-methods`. A proxy routes/pools/caches; it does not (generally)
> replan the optimizer.

---

## The MongoDB angle

**MongoDB generally does NOT use third-party query proxies the way MySQL/PostgreSQL shops do**, and
understanding *why* is the useful insight. Full detail: `references/mongodb-angle.md`.

- **`mongos` is MongoDB's native query router**, the closest first-party equivalent of a database
  proxy. For a **sharded** cluster it is the only interface the application uses (clients connect to
  `mongos`, never directly to shards); it routes by shard key (targeted) or scatter-gathers to all
  shards when no shard key is present, then merges results. It holds no persistent state, caching
  the routing table from the **config servers** and reloading lazily on a stale-version
  error.[^mongos][^mongos-config]
- **The proxy responsibilities are built into the driver, not bolted on.** Every official MongoDB
  driver implements standardized specs that cover what a MySQL/PG shop adds via an external proxy:
  **SDAM** (server discovery/monitoring, failover reaction); **CMAP** (client-side connection
  pooling, `maxPoolSize` default 100); **read preference + server selection** (read/write splitting
  by routing reads to secondaries); and **retryable reads/writes** (on by default).[^mongo-retry] A
  MongoDB engineer states it directly: because drivers consistently provide connection monitoring
  and pooling, "external connection pooling solutions aren't required (ex: Pgpool,
  PgBouncer)."[^mongo-sdam][^mongo-cmap][^mongo-readpref][^mongo-bevilacqua]
- **Adjacent query-access layers (query surfaces, not proxies):** the **BI Connector / `mongosqld`**
  (a relational/SQL view over the MySQL wire protocol; **end-of-life after September 2026 per the
  docs**, volatile; successor is Atlas SQL / MongoSQL) and **Atlas Data Federation** (a distributed
  query engine across Atlas clusters + object stores through one MQL/Atlas-SQL endpoint). See the
  reference for both.[^mongo-bic][^mongo-federation]
- **Honest nuance:** `mongos` itself *is* arguably a proxy, and protocol-compatible poolers exist
  for extreme serverless fan-out (AWS's `mongobetween`, positioned for DocumentDB, not Atlas) — but
  MongoDB's own answers to connection exhaustion are native (a module-scope `MongoClient` with a
  small `maxPoolSize`; Atlas App Services; the driver's `loadBalanced=true` mode).[^mongo-serverless]

> For the **internals** of MongoDB driver pooling/SDAM (CMAP state machine, server selection
> algorithm), defer to **`mongodb-driver-internals`**; this skill uses them only to explain the
> "why no third-party proxy" position.

---

## Anti-patterns

1. **Transaction-mode pooling breaking session state.** The pooling mode that delivers the
   multiplexing payoff (PgBouncer transaction mode; ProxySQL/RDS Proxy multiplexing) hands
   successive transactions to *different* backends, so session state leaks or is lost:
   `SET`/GUCs, `LISTEN`/`NOTIFY`, SQL-level prepared statements, session advisory locks, temp
   tables, `WITH HOLD` cursors. **The failure is silent and non-deterministic** ("I execute a
   statement in Process A and Process B errors out"). Use `SET LOCAL`, transaction-scoped advisory
   locks, a proxy init query, or session mode (giving up the win).[^pgb-features][^txn-break]
2. **The proxy as a hidden SPOF.** A single un-HA proxy means every app behind it loses the DB
   when it dies; failover thundering-herds and pool saturation while the DB is healthy are real
   incident classes (OpenAI's 2023 outage involved PgBouncer pools clogged by slow queries). Run
   the proxy HA and size headroom.[^spof]
3. **Cache staleness / invalidation hazards.** Proxy result caches go stale; TTL-only caches
   (ProxySQL) serve stale data until expiry; "automatic" invalidation has sharp edges
   (transactions, untracked tables, clock skew). Don't cache data that must be read-your-writes
   fresh.[^proxysql-cache][^heimdall]
4. **Masking a real scaling/schema problem.** A caching proxy can hide an unindexed query so well
   that a later regression goes unnoticed until the cache fails — then the blast radius is huge.
   A slow query holds a pooled connection for its whole duration; **fix the query, not the pool
   size.** Order: measure → fix queries/indexes/N+1 → fix schema → tune → *then* cache.[^mask]
5. **Naive read/write split rules.** Generic `^SELECT → replica` routing breaks transactions,
   session state, and read-after-write consistency. Analyze by digest and target specific
   queries; keep transactions and post-write reads on the primary.[^proxysql-rwsplit]
6. **RDS Proxy pinning that nullifies the benefit.** A default `SET search_path`, ORM-driven
   prepared statements, or session locks can pin most connections, so you pay for the proxy
   without multiplexing. Watch `DatabaseConnectionsCurrentlySessionPinned`; push common `SET`s to
   the init query.[^rdsproxy-pinning]
7. **Enabling the SQL firewall before fully training it** blocks all legitimate traffic. Train
   long enough to cover every code path (new features/ORM variants emit new digests).[^proxysql-firewall]
8. **Observability distortion.** A multiplexing proxy makes `pg_stat_activity` / processlist no
   longer map 1:1 to clients, breaking standard connection-storm triage — use the proxy's own
   stats (per-digest, pool state) as the new vantage point.[^spof]

---

## Adjacent concepts (cross-references, not owned here)

- **AI Gateways & LLM Proxy Infrastructure** → `llm-ai-gateways` — the **sibling pattern** for LLM
  request traffic (the analog of every capability above). Cite it; don't duplicate it.
- **MongoDB driver internals (CMAP/SDAM)** → `mongodb-driver-internals`.
- **Compute–storage disaggregation** (Aurora / Neon / PolarDB / object-storage OLTP) →
  `mongodb-operations-expert` (cloud-native-database-disaggregation).
- **Cross-engine query federation** (Trino / Presto / Starburst) — *adjacent only*; a federation
  engine joins across heterogeneous sources, which overlaps the "single endpoint over many
  backends" idea but is a query engine, not a transparent protocol proxy.
- **Learned / ML-based query optimization** (learned cardinality estimation, Bao steering, Neo,
  OtterTune) → `da-analytical-methods`.
- **Generic API gateway / BFF design** → `software-engineering-patterns`.

---

## References

[^proxysql-arch]: ProxySQL — Architecture Overview & Multi-Layer Configuration. https://proxysql.com/documentation/architecture/ ; https://proxysql.com/documentation/main-runtime/multi-layer-configuration/
[^maxscale-arch]: MariaDB — MaxScale Architecture & Configuration Guide. https://mariadb.com/docs/maxscale/maxscale-architecture ; https://mariadb.com/docs/maxscale/maxscale-management/deployment/installation-and-configuration/maxscale-configuration-guide
[^vitess-arch]: Vitess — Architecture & vtgate reference. https://vitess.io/docs/25.0/overview/architecture/ ; https://vitess.io/docs/25.0/reference/programs/vtgate/
[^vitess-history]: Vitess — History (YouTube origin, CNCF graduation Nov 2019). https://vitess.io/docs/24.0/overview/history/
[^pgb-features]: PgBouncer — Features (pooling modes; transaction-mode feature matrix). https://www.pgbouncer.org/features.html
[^pgcat]: postgresml/pgcat — README & docs (Rust pooler: LB, routing, sharding, failover). https://github.com/postgresml/pgcat ; https://github.com/postgresml/docs/blob/main/pgcat.md
[^pgpool]: pgpool.net — What is Pgpool-II? (pooling, LB, failover, watchdog, query cache, replication). https://www.pgpool.net/docs/latest/en/html/intro-whatis.html
[^rdsproxy-overview]: AWS — Amazon RDS Proxy (managed pooling/multiplexing for serverless). https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
[^heimdall]: Heimdall Data — product docs & AWS co-marketing (vendor-heavy; auto query caching/invalidation claims flagged). https://www.heimdalldata.com/ (and AWS Database Blog co-authored posts)
[^mongos]: MongoDB Manual — Routing with mongos; Distributed Queries. https://www.mongodb.com/docs/manual/core/sharded-cluster-query-router/ ; https://www.mongodb.com/docs/v8.2/core/distributed-queries/
[^percona-compare]: Percona — Comparisons of Proxies for MySQL (ProxySQL L7 CPU cost vs HAProxy). https://www.percona.com/blog/comparisons-of-proxies-for-mysql/
[^citus-conns]: Citus Data — Analyzing the Limits of Connection Scalability in Postgres. https://www.citusdata.com/blog/2020/10/08/analyzing-connection-scalability
[^pg-process]: AWS Database Blog — Resources consumed by idle PostgreSQL connections. https://aws.amazon.com/blogs/database/resources-consumed-by-idle-postgresql-connections/
[^proxysql-mux]: ProxySQL — Multiplexing (conditions that disable multiplexing/pin a backend). https://proxysql.com/documentation/multiplexing/
[^rdsproxy-pinning]: AWS — Managing connections with RDS Proxy: pinning conditions & avoidance. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-managing.html
[^pg-poolsize]: PostgreSQL wiki — Number Of Database Connections ((cores×2)+spindles; saturation knee). https://wiki.postgresql.org/wiki/Number_Of_Database_Connections
[^hikari-poolsize]: HikariCP wiki — About Pool Sizing (small pool beats large; benchmark). https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing
[^maxscale-causal]: MariaDB — Ensuring causal consistency with MaxScale readwritesplit (causal_reads modes); CCR filter. https://mariadb.com/docs/maxscale/maxscale-use-cases/readwrite-split-router-usage/ensuring-causal-consistency-with-maxscales-readwrite-split-router
[^proxysql-rwsplit]: ProxySQL — How to Set Up Read/Write Split (digest-based routing; warns against naive ^SELECT rules). https://proxysql.com/documentation/proxysql-read-write-split-howto/
[^proxysql-cache]: ProxySQL — Query Cache (TTL-only, no invalidation/LRU, breaks on prepared statements). https://proxysql.com/documentation/query-cache/
[^maxscale-cache]: MariaDB — MaxScale Cache filter (soft/hard TTL, invalidation, storage backends). https://mariadb.com/docs/maxscale/reference/maxscale-filters/maxscale-cache
[^proxysql-firewall]: ProxySQL — Firewall Whitelist (2.0.9+; OFF/DETECTING/PROTECTING; SQLi engine). https://proxysql.com/documentation/firewall-whitelist/
[^maxscale-firewall]: MariaDB — MaxScale Firewall Filter (dbfwfilter; allow/block; learn/supervise/enforce). https://mariadb.com/docs/maxscale/reference/maxscale-filters/maxscale-firewall-filter
[^owasp-sqli]: OWASP — SQL Injection Prevention Cheat Sheet (parameterized queries primary; allow-list secondary). https://cheatsheetseries.owasp.org/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.html
[^rdsproxy-failover]: AWS — RDS Proxy (failover handling; up to ~66% reduction, scenario-dependent). https://aws.amazon.com/rds/proxy/ ; https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
[^spof]: Synthesis — proxy SPOF & HA patterns (Keepalived/VRRP; sidecar vs centralized; OpenAI 2023 PgBouncer incident). Percona/Mydbops/nuffing.com + OpenAI postmortem.
[^rdsproxy-iam]: AWS — Setting up IAM authentication / Secrets Manager with RDS Proxy (credential brokering; TLS required). https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-setup.html
[^decision]: Synthesis — proxy-vs-driver-vs-app decision & cost framing (PgBouncer/HikariCP sizing, cache-aside, RDS Proxy serverless). PostgreSQL wiki, AWS docs, PGConf.EU 2024 pooler comparison.
[^percona-overhead]: Percona — ProxySQL Overhead Explained and Measured (~25–45µs added latency; config-dependent). https://www.percona.com/blog/proxysql-overhead-explained-and-measured/
[^txn-break]: Heroku Dev Center — PgBouncer Configuration and Best Practices (transaction-mode session-state breakage). https://devcenter.heroku.com/articles/best-practices-pgbouncer-configuration
[^mask]: Synthesis — proxy/cache masking a scaling/schema problem; fix the query not the pool. HikariCP wiki + practitioner writeups.
[^mongos-config]: MongoDB Manual — Config Servers; Sharded Cluster Components; FAQ Sharding (lazy routing-cache reload). https://www.mongodb.com/docs/manual/core/sharded-cluster-config-servers/ ; https://www.mongodb.com/docs/v7.0/faq/sharding/
[^mongo-sdam]: MongoDB — Server Discovery and Monitoring (SDAM) spec. https://github.com/mongodb/specifications/blob/master/source/server-discovery-and-monitoring/server-discovery-and-monitoring.md
[^mongo-cmap]: MongoDB — Connection Pool Overview & CMAP spec (maxPoolSize 100, maxConnecting 2). https://www.mongodb.com/docs/manual/administration/connection-pool-overview/ ; https://github.com/mongodb/specifications/blob/master/source/connection-monitoring-and-pooling/connection-monitoring-and-pooling.md
[^mongo-readpref]: MongoDB Manual — Read Preference & Server Selection (read routing; mongos load-balancing within latency window). https://www.mongodb.com/docs/manual/core/read-preference/ ; https://www.mongodb.com/docs/manual/core/read-preference-mechanics/
[^mongo-retry]: MongoDB Manual — Retryable Writes & Retryable Reads (on by default). https://www.mongodb.com/docs/manual/core/retryable-writes/ ; https://www.mongodb.com/docs/manual/core/retryable-reads/
[^mongo-bevilacqua]: Alex Bevilacqua (MongoDB) — MongoDB and Load Balancer Support (external poolers not required). https://www.alexbevi.com/blog/2024/03/08/mongodb-and-load-balancer-support/
[^mongo-bic]: MongoDB — mongosqld / What is the BI Connector / transition to Atlas SQL (EOL after Sept 2026). https://www.mongodb.com/docs/bi-connector/current/reference/mongosqld/ ; https://www.mongodb.com/docs/sql-interface/transition-bic-to-atlas-sql/
[^mongo-federation]: MongoDB — Atlas Data Federation (distributed query engine across clusters + object stores). https://www.mongodb.com/docs/atlas/data-federation/
[^mongo-serverless]: AWS — Scale connections with DocumentDB using mongobetween; MongoDB Load Balancers spec (loadBalanced=true). https://aws.amazon.com/blogs/database/scale-your-connections-with-amazon-documentdb-using-mongobetween/ ; https://github.com/mongodb/specifications/blob/master/source/load-balancers/load-balancers.md
