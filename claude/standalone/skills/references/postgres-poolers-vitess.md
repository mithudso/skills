# PostgreSQL Poolers (PgBouncer, PgCat, Pgpool-II) & Vitess

> Provenance: reference under the standalone `database-proxies-query-middleware` skill. Created
> 2026-06-18 via `/dr` deep research. `verified-as-of: 2026-06-18` — version-pinned facts carry
> the version they reflect; re-verify before relying on them.

## Contents
- [Why PostgreSQL needs an external pooler](#why-postgresql-needs-an-external-pooler)
- [PgBouncer — the pure pooler and what breaks](#pgbouncer--the-pure-pooler-and-what-breaks)
- [PgCat — the Rust alternative](#pgcat--the-rust-alternative)
- [Pgpool-II — the everything-middleware](#pgpool-ii--the-everything-middleware)
- [Vitess / vtgate — the sharding-proxy exemplar](#vitess--vtgate--the-sharding-proxy-exemplar)
- [Pool-sizing math & cross-cutting trade-offs](#pool-sizing-math--cross-cutting-trade-offs)
- [Sources](#sources)

---

## Why PostgreSQL needs an external pooler

- **Process-per-connection.** PostgreSQL uses a multi-process (not multi-threaded) model: the
  postmaster `fork()`s a dedicated OS backend process per client connection, which dies on
  disconnect. Each backend has private memory (~5–10 MB base + query memory); `work_mem` is
  allocated **per-operation**, so worst-case memory ≈ `work_mem × operations × parallel_workers ×
  connections`.[^1][^2][^3][^4]
- **Idle connections still cost** memory/CPU and that memory isn't fully freed even after
  `DISCARD ALL`; some internal structures scale super-linearly with `max_connections`.[^1][^5][^6]
- **`max_connections` is rigid** (default 100; changing it needs a server restart).[^3][^7]
- **What a pooler does:** multiplex many client connections onto a small set of warm backends
  (typical 1000 client conns → 20–50 backends).[^2][^3]

## PgBouncer — the pure pooler and what breaks

PgBouncer is **a pure pooler — no query routing, load balancing, or sharding** (single-threaded,
libevent, C).[^8][^9][^10]

| Mode | Backend held for… | Multiplexing | Multi-stmt txns |
|---|---|---|---|
| **Session** (default) | whole client connection | none (reuse only) | OK |
| **Transaction** (web-app default) | one `BEGIN…COMMIT` | high | OK |
| **Statement** | one statement | highest | forbidden (forces autocommit) |

**What breaks in transaction mode** (the load-bearing detail) — anything stored in **session
state** is lost or **leaks to the next client** reusing that backend: session `SET`/GUCs (e.g.
`search_path`), `LISTEN`, SQL-level `PREPARE`/`DEALLOCATE`, session-level advisory locks,
`WITH HOLD` cursors, session temp tables, `LOAD`. From PgBouncer's own feature matrix these are
marked **"Never"** in transaction mode.[^8][^9][^11][^13]

**Mitigations:** use `SET LOCAL` (transaction-scoped) not `SET`; `pg_advisory_xact_lock()` not
`pg_advisory_lock()`; run migrations on a direct connection; `server_reset_query = DISCARD ALL`.[^11][^13][^14]

**Prepared statements — version-stamped:**
- SQL-level `PREPARE`/`EXECUTE` is **always broken** in transaction mode.[^8][^9][^14]
- **Protocol-level** prepared statements (what ORMs use via Parse/Bind/Execute) were broken until
  **PgBouncer 1.21.0 (2023)**, which re-prepares on the routed backend when
  `max_prepared_statements > 0`.[^9][^13][^15]
- **PgBouncer 1.24.0 (Jan 10, 2025)** turned prepared-statement support **on by default**
  (`max_prepared_statements = 200`).[^15][^16] Footgun: PHP/PDO needs PHP 8.4+ **and** libpq 17;
  `pg_dump`'s `SET search_path` can poison a shared backend.[^13][^14][^17]

**Sizing keys:** `max_client_conn` (default 100, max client→PgBouncer) vs `default_pool_size`
(default 20, max PgBouncer→Postgres **per user/db pair**). FD ceiling ≈ `max_client_conn + (max
pool_size × databases × users)`. The headline value: `max_client_conn = 5000` served by
`default_pool_size = 50`.[^18][^19][^20]

## PgCat — the Rust alternative

A PostgreSQL pooler/proxy "like PgBouncer," **written in Rust, MIT-licensed**, started at
PostgresML (2022): "PostgreSQL pooler with sharding, load balancing and failover support."[^21][^22]

What it adds beyond PgBouncer:
1. **Multi-threaded** — uses all cores (PgBouncer is single-threaded, ~one-core ceiling).[^22][^23][^24]
2. **Load balancing** across replicas (random / least-outstanding-connections).[^21][^22][^25]
3. **Read/write query routing** via a built-in SQL parser (`SELECT`→replica, writes/explicit txns→primary).[^21][^22][^26]
4. **Failover** — health checks; unreachable servers "banned" for a duration; the **primary is
   never banned**; if all are banned the ban list clears.[^21][^22]
5. **Sharding** (`PARTITION BY HASH`, same `pg_bigint_hash` seed; or explicit `SET SHARD TO`).[^21][^22][^26]
6. **Mirroring.**[^21][^22]

Production: Instacart adopted/contributes to PgCat (a cross-language replacement for their Ruby
`Makara` app-side library).[^21][^25]

**DISCONFIRMING — the query parser is the weak point:** open issues (#865/#884/#887, plus an SO
report) show that when the `sqlparser` crate **fails to parse a query** (multi-statement strings,
`DO` blocks, `SELECT FOR NO KEY UPDATE`), PgCat **falls back to role "any" instead of primary**, so
writes get routed to a read replica and fail with `cannot execute INSERT in a read-only
transaction` — **intermittently**. Sharding is tied to the Postgres partition hash seed, making
**resharding costly**.[^26][^28][^29][^30]

## Pgpool-II — the everything-middleware

Middleware (since ~2003) that is a **comprehensive scalability/HA suite**, not just a pooler;
multi-process (a forked child per client connection enables SQL parsing/routing).[^31][^32][^33]

- **Connection pooling** — reuses backend connections with matching properties, but **only
  session-level pooling, and each child keeps its own pool** (a limitation vs PgBouncer's
  transaction pooling).[^31][^33]
- **Load balancing** — distributes read (`SELECT`) queries across replicas; node chosen at session
  start unless `statement_level_load_balance`.[^31][^35]
- **Automatic failover** + online recovery (add a standby without downtime).[^31][^36]
- **Watchdog HA** — HA for Pgpool-II *itself* (eliminates the pooler SPOF) with a VIP controller
  and a **quorum to prevent split-brain**.[^31][^32][^33]
- **In-memory query result cache** — caches `SELECT` results (shared mem or memcached);
  auto-invalidates a table's cache on update; **registered at commit, which can break row-visibility
  under REPEATABLE READ**; docs say disable if hit ratio < 70%.[^31][^37]
- **Replication** — six clustering modes (streaming, logical, native, snapshot-isolation,
  main-replica, raw).[^39][^33]

**vs PgBouncer:** Pgpool does routing/LB/failover/caching/replication; PgBouncer does none of those
(but adds the least latency, <100µs).[^27][^8]

**DISCONFIRMING — operational complexity (heavily attested):** "hundreds of configuration
parameters"; "complex… non-trivial to set up correctly." For a 100%-write workload Pgpool "just
adds overhead"; SSL "eats CPU." A long-standing report (and a 2025 geo-distributed follow-up)
shows Pgpool sync-messages **all** backends per statement — in a ~20 ms cross-region setup this
caused ~20× degradation; a maintainer confirms the fix is targeted for **v4.8 (~end 2026)**. Cannot
scale writes / no sharding.[^27][^40][^41][^42][^43]

## Vitess / vtgate — the sharding-proxy exemplar

The proxy-based **horizontal-scaling / sharding** exemplar. Three core components + a metadata
store:

- **vtgate** — a **stateless proxy** that accepts app requests and routes to the right tablet(s).
  Speaks **both the MySQL protocol and gRPC**, so apps connect to it **as if it were a single MySQL
  server** ("a fat client provided as a service" — routing logic moved out of the app).[^44][^45][^46]
- **vttablet** — sits in front of each MySQL instance; manages the **connection pool**, query
  rewriting, the query consolidator, and safeguards (query blacklisting, row-count limits).[^47][^48][^49]
- **Topology Service** — consistent metadata store (etcd/ZooKeeper/Consul) holding
  Keyspaces/Shards/Tablets + the serving graph; provides cluster locks & primary election.[^50]
- **VSchema** — the abstraction presenting a unified single-MySQL view; holds the **Vindex** per
  sharded table (maps sharding-key value → keyspace ID → shard). Every VSchema needs a Primary
  Vindex.[^52]

How vtgate presents one endpoint over many MySQLs: it maintains health-check connections to all
tablets + reads the shard map, then uses VSchema/Vindex + the query's `WHERE` clause to compute the
sharding key and route. With V3 the app sends queries "as if it's a single database." A sharded
query **without a usable vindex "scatters"** to all shards; results are merged (gather); aggregations
push down per-shard and combine at vtgate. vtgate caches execution plans in an LRU cache.[^50][^54][^55][^56]

**vttablet pooling + consolidator:** connection pooling was Vitess's first problem solved (YouTube,
2011) — forcing apps to connect only to vttablet made the max MySQL connection count centrally
configurable instead of growing with app-server count. The **query consolidator** holds a duplicate
of an in-flight identical query and serves the first result to all waiters (protecting MySQL from
hot-query QPS spikes) — but this can **lose read-after-write consistency** (disable on primary via
`--enable-consolidator=false --enable-consolidator-replicas=true`).[^48][^49][^57]

**Origin:** created 2010 for YouTube's MySQL scaling; serves YouTube since 2011; a graduated CNCF
project (Nov 2019). "Pushes MySQL closer to NoSQL by supporting sharding and replication, with
trade-offs on transactions and consistency."[^49][^57][^58]

**DISCONFIRMING — Vitess limitations:** unsupported/restricted MySQL features (stored procs
returning results, triggers, events, `LOCK TABLES`, `CREATE/DROP DATABASE`; window functions/temp
tables restricted; correlated cross-shard subqueries may fail). Cross-shard joins use expensive
nested-loop joins (co-locate on the same vindex); cross-shard writes need 2PC
(`transaction_mode=TWOPC`, ~an order of magnitude slower) and **even TWOPC does not guarantee
cross-shard isolation**. It's a fleet of processes + topology store — far heavier than a single
pooler; recommended shard size ≤ ~250 GB.[^59][^62][^63][^64][^65]

## Pool-sizing math & cross-cutting trade-offs

- **Authoritative formula (PostgreSQL wiki):** optimal *active* connections ≈ **`(core_count × 2)
  + effective_spindle_count`** (~1 for SSD/NVMe, ~0 if fully cached) → typically ~10–20. The TPS
  curve climbs to a saturation "knee" then **falls** due to contention.[^7][^66][^67]
- **Little's Law cross-check:** `Pool Size ≈ Requests/sec × Avg Query Duration` lands on the same
  ~10–20 range.[^68]
- **Why bigger hurts:** more connections → super-linear LWLock contention + context-switch cost +
  `max_connections`-proportional overhead. HikariCP benchmark: pool 10 was the sweet spot on a
  4-core box; pool 200 was a **44% throughput reduction with 10× tail latency**.[^6][^7][^67]
- **Anti-pattern by name:** "500 users so I need 500 connections" — the point of pooling is that
  500 users share ~10 backends (each holds a connection 2–5 ms). The fix for a saturated pool is
  **faster queries, not more connections.**[^67][^68][^69]
- **Transaction/statement pooling is only safe if the app treats every transaction as stateless** —
  the same rule that drives PgBouncer's "Never" matrix, PgCat's routing fragility, and Vitess's
  consolidation trade-off: **proxies buy scale by relaxing per-session/per-read guarantees.**[^11][^8][^14]
- **Operational cost (disconfirming the "always add a pooler" reflex):** every tier adds latency, a
  new failure domain, and auth/TLS complexity; **don't deploy a pooler at all if connection count
  is naturally low and stable.** Escalation path: **PgBouncer (transaction mode, correct sizing) →
  PgCat (in-pooler routing/failover/sharding) → platform-scale tools only when needed.**[^23][^7][^27]

## Sources

[^1]: AWS Database Blog — Resources consumed by idle PostgreSQL connections. https://aws.amazon.com/blogs/database/resources-consumed-by-idle-postgresql-connections/
[^2]: Citus Data — Analyzing the Limits of Connection Scalability in Postgres. https://www.citusdata.com/blog/2020/10/08/analyzing-connection-scalability
[^3]: PlanetScale database-skills — process-architecture. https://raw.githubusercontent.com/planetscale/database-skills/main/skills/postgres/references/process-architecture.md
[^4]: Quang Chien Tran — Deep Inside PostgreSQL: Processes, Forking, Memory. https://quangchientran.substack.com/p/deep-inside-postgres-processes-forking
[^5]: Timescale pg-guide — process architecture. https://github.com/timescale/pg-guide/blob/main/01-process-architecture.md
[^6]: dba.stackexchange — pool recommendation vs 100s of connections. https://dba.stackexchange.com/questions/82558/
[^7]: PostgreSQL wiki — Number Of Database Connections. https://wiki.postgresql.org/wiki/Number_Of_Database_Connections
[^8]: PgBouncer — Features. https://www.pgbouncer.org/features.html
[^9]: augusteo.com — How PgBouncer Works. https://www.augusteo.com/blog/how-pgbouncer-works
[^10]: pgbouncer — doc/usage.md. https://github.com/pgbouncer/pgbouncer/blob/master/doc/usage.md
[^11]: Heroku Dev Center — PgBouncer Configuration and Best Practices. https://devcenter.heroku.com/articles/best-practices-pgbouncer-configuration
[^13]: PgBouncer — FAQ. https://www.pgbouncer.org/faq.html
[^14]: Stack Harbor KB — PgBouncer transaction-pool mode. https://stackharbor.com/en/knowledge-base/pgbouncer-transaction-pool-advanced/
[^15]: PgBouncer — Changelog. https://www.pgbouncer.org/changelog.html
[^16]: PgBouncer 1.24.0 release notes. https://www.pgbouncer.org/2025/01/pgbouncer-1.24.0
[^17]: pgbouncer issues #976 / #1313 (warn on unsupported; postgres_fdw search_path). https://github.com/pgbouncer/pgbouncer/issues/976
[^18]: PgBouncer — config. https://www.pgbouncer.org/config.html
[^19]: Crunchy Data — pgBouncer config docs. https://access.crunchydata.com/documentation/pgbouncer/1.25.2/config/
[^20]: pgbouncer issue #643 — max_client_conn and default_pool_size. https://github.com/pgbouncer/pgbouncer/issues/643
[^21]: postgresml/pgcat — README. https://github.com/postgresml/pgcat
[^22]: postgresml/docs — pgcat.md. https://github.com/postgresml/docs/blob/main/pgcat.md
[^23]: Pavan Rangani — PgBouncer vs PgCat vs Supavisor. https://pavanrangani.com/blog/postgres-connection-pooling-pgbouncer-pgcat-supavisor
[^24]: Gold Lapel — Poolers Compared: PgBouncer vs Pgpool-II vs PgCat vs Odyssey. https://goldlapel.com/compare/postgresql-poolers
[^25]: Instacart Tech — Adopting PgCat: A Nextgen Postgres Proxy. https://tech.instacart.com/adopting-pgcat-a-nextgen-postgres-proxy-3cf284e68c2f
[^26]: 3manuek — Evaluating PGCat's Sharding by hash. https://tr3s.ma/posts/2025-01/pgcat/
[^27]: Julian Markwort — Comparing Connection Poolers for PostgreSQL (PGConf.EU 2024). https://www.postgresql.eu/events/pgconfeu2024/sessions/session/5846/slides/547/comparing_poolers.pdf
[^28]: pgcat issue #865 — sqlparser parse failures misroute. https://github.com/postgresml/pgcat/issues/865
[^29]: pgcat issue #884 — parse error routes to "any" not default_role. https://github.com/postgresml/pgcat/issues/884
[^30]: SO — PgCat sometimes routes write queries to replicas. https://stackoverflow.com/questions/79804089/
[^31]: pgpool.net — What is Pgpool-II? https://www.pgpool.net/docs/latest/en/html/intro-whatis.html
[^32]: Percona — Understanding Pgpool-II Use Cases and Benefits. https://www.percona.com/blog/postgresql-pgpool-ii-use-cases-and-benefits/
[^33]: PGConf.EU 2023 — Pgpool: What, Why, Where (slides). https://www.postgresql.eu/events/pgconfeu2023/sessions/session/4808/slides/449/
[^35]: pgpool.net — Load Balancing config. https://pgpool.net/docs/42/en/html/runtime-config-load-balancing.html
[^36]: pgpool.net — failover/online recovery (intro-whatis). https://www.pgpool.net/docs/latest/en/html/intro-whatis.html
[^37]: Tatsuo Ishii — In Memory Query Cache. https://tatsuo-ishii.github.io/pgpool-II/current/runtime-in-memory-query-cache.html
[^39]: pgpool.net — Configuring Pgpool-II (clustering modes). https://www.pgpool.net/docs/pgpool-II-4.4.2/en/html/configuring-pgpool.html
[^40]: pgpool-general mailing list — poor performance / LB bottlenecked by master. https://www.pgpool.net/pipermail/pgpool-general/2013-May/001702.html
[^41]: ADHDecode — Pooler Alternatives: PgBouncer vs PgPool vs pgcat. https://adhdecode.com/databases/connection-pooling-and-application-integration/pgcat-pgpool-ii-alternatives/
[^42]: Pgpool-II FAQ (SSL/OpenSSL CPU). https://pgpool.github.io/faq/
[^43]: pgpool2 issue #130 — degraded perf geo-distributed (fix targeted v4.8). https://github.com/pgpool/pgpool2/issues/130
[^44]: Vitess Docs — vtgate. https://vitess.io/docs/25.0/reference/programs/vtgate/
[^45]: Vitess Docs — Architecture. https://vitess.io/docs/25.0/overview/architecture/
[^46]: vitessio/vitess — VTGateV3Features.md. https://github.com/vitessio/vitess/blob/main/doc/design-docs/VTGateV3Features.md
[^47]: Vitess Docs — VTTablet Connection Pools and Sizing. https://vitess.io/docs/25.0/reference/features/connection-pools/
[^48]: Vitess Docs — Query Consolidation. https://vitess.io/docs/23.0/user-guides/configuration-advanced/query-consolidation/
[^49]: PlanetScale — Connection pooling in Vitess. https://planetscale.com/blog/connection-pooling
[^50]: Vitess Docs — Topology Service. https://vitess.io/docs/25.0/reference/features/topology-service/
[^52]: Vitess Docs — VSchema. https://vitess.io/docs/24.0/reference/features/vschema/
[^54]: Vitess Docs — Sharding. https://vitess.io/docs/23.0/reference/features/sharding/
[^55]: Vitess Docs — Execution Plans. https://vitess.io/docs/23.0/concepts/execution-plans/
[^56]: PlanetScale database-skills — query-serving. https://raw.githubusercontent.com/planetscale/database-skills/main/skills/vitess/references/query-serving.md
[^57]: Vitess Docs — History. https://vitess.io/docs/24.0/overview/history/
[^58]: CNCF — Vitess graduation (via Vitess history). https://vitess.io/docs/24.0/overview/history/
[^59]: Vitess Docs — MySQL Compatibility. https://vitess.io/docs/25.0/reference/compatibility/mysql-compatibility/
[^62]: Vitess Docs — Sharding Guidelines. https://vitess.io/docs/23.0/user-guides/vschema-guide/sharding-guidelines/
[^63]: Vitess Docs — FAQ: cross-shard JOINs or Transactions. https://vitess.io/docs/faq/sharding/advanced/can-i-use-vitess-to-do-cross-shard-joins-or/
[^64]: Vitess Docs — Shard Isolation and Atomicity Model. https://vitess.io/docs/25.0/user-guides/configuration-advanced/shard-isolation-atomicity/
[^65]: vitessio/vitess issue #16515 — cross-keyspace joins broke v16→v18. https://github.com/vitessio/vitess/issues/16515
[^66]: HikariCP wiki — About Pool Sizing. https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing
[^67]: Gold Lapel — HikariCP Pool Sizing for Postgres. https://goldlapel.com/grounds/connection-pooling/hikaricp-pool-sizing-postgres
[^68]: Michal Drozd — Connection Pool Sizing with Little's Law. https://www.michal-drozd.com/en/blog/connection-pool-littles-law/
[^69]: Gold Lapel — PostgreSQL max_connections: Why 100 Is Not Enough. https://goldlapel.com/grounds/connection-pooling/postgresql-max-connections
