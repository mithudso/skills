# ProxySQL & MariaDB MaxScale (MySQL-ecosystem query-aware proxies)

> Provenance: reference under the standalone `database-proxies-query-middleware` skill. Created
> 2026-06-18 via `/dr` deep research. `verified-as-of: 2026-06-18` — version-pinned facts carry
> the version they reflect; re-verify before relying on them.

## Contents
- [ProxySQL architecture](#proxysql-architecture)
- [ProxySQL capabilities](#proxysql-capabilities)
- [MariaDB MaxScale architecture](#mariadb-maxscale-architecture)
- [Read/write splitting & the consistency hazard](#readwrite-splitting--the-consistency-hazard)
- [Operations: latency, HA, observability, anti-patterns](#operations-latency-ha-observability-anti-patterns)
- [ProxySQL vs MaxScale at a glance](#proxysql-vs-maxscale-at-a-glance)
- [Sources](#sources)

These are the two canonical **query-aware** proxies in the MySQL/MariaDB world: they parse SQL and
can route, split, cache, rewrite, and firewall. Both are an extra network hop (Layer-7,
content-aware), trading CPU for capability versus a dumb TCP load balancer.

---

## ProxySQL architecture

- **What it is.** A high-performance, protocol-aware proxy for MySQL (and now PostgreSQL), in C++,
  using an asynchronous multi-threaded event-driven model (`libev`). The app connects to ProxySQL
  (default data port **6033**), which parses every query and forwards it to a backend. The app
  sees one endpoint and is unaware of routing.[^1][^2]
- **Threads.** Worker threads (client traffic, parsing, relaying), an Admin thread (port **6032**
  MySQL-protocol), Monitor threads (backend health/replication-lag checks), Cluster threads (P2P
  config sync). Default `mysql-threads=4` (capped at 4 cores unless raised — a tuning trap on
  many-core hosts).[^1]
- **3-layer configuration system.** All config lives in an internal SQLite DB across three layers:
  - **RUNTIME** — in-memory structures the live workers use; the *only* layer that affects
    behavior; never edited directly.
  - **MEMORY** (`main`) — an in-memory SQLite exposed via a MySQL-compatible admin interface;
    edited with standard `INSERT`/`UPDATE`/`DELETE`.
  - **DISK** — persistent SQLite (default `proxysql.db`).
  - Movement: `LOAD MYSQL QUERY RULES TO RUNTIME` (activate), `SAVE … TO DISK` (persist), `LOAD …
    FROM DISK` (revert). Startup priority DISK→MEMORY→RUNTIME. Goal: dynamic reconfiguration with
    no restart and no dropped connections.[^1][^2][^3]
- **Hostgroups.** Backends are grouped in `mysql_servers` (`hostgroup_id, hostname, port`); each
  hostgroup has its own connection pool. `mysql_replication_hostgroups` maps writer↔reader groups;
  ProxySQL auto-moves a server between them based on its `read_only` flag.[^1][^2]
- **Query Processor + `mysql_query_rules`.** The "brain." Per query it decides routing, rewriting,
  caching, blocking. Key columns: `rule_id` (PK, evaluated ascending), `active`, `username`,
  `schemaname`, `client_addr`, `digest` (exact), `match_digest` (regex on normalized digest),
  `match_pattern`/`replace_pattern` (regex rewrite), `destination_hostgroup`, `cache_ttl`,
  `mirror_hostgroup`/`mirror_flagOUT`, `multiplex`, `flagIN`/`flagOUT` (rule chaining), and
  **`apply`** (when `apply=1`, evaluation stops for that query). Supports RE2 and PCRE; compiled
  regexes are cached. Docs advise `match_digest` over `match_pattern` (smaller, faster).[^2][^4]
- **Query digest / normalization.** A per-query fingerprint with literals normalized (`WHERE id=5`
  → `WHERE id=?`) powers `stats_mysql_query_digest` and rule matching on query *types*.[^2]

## ProxySQL capabilities

- **Connection pooling & multiplexing.** Frontend (client→ProxySQL) connections are decoupled from
  backend (ProxySQL→MySQL). **Multiplexing** lets multiple frontend connections reuse one backend
  connection — MySQL is thread-per-connection, so each connection costs even when idle; the pool +
  multiplexing sharply cut backend connection count (Instawork reported an ~80% drop after a
  clustered ProxySQL deployment).[^1][^2][^13]
- **Conditions that disable multiplexing (pin the backend):** active transaction (re-enabled on
  commit/rollback); `LOCK TABLES`/`FLUSH TABLES WITH READ LOCK`; `GET_LOCK()`; any query with `@`
  in the digest (user/session variables — never re-enabled, except hardcoded `SELECT @@tx_isolation`
  / `SELECT @@version`); `SQL_CALC_FOUND_ROWS`; `CREATE TEMPORARY TABLE`; `PREPARE`. The
  `multiplex` rule column overrides per-rule. `mysql-auto_increment_delay_multiplex` (default **5**,
  since 1.4.14) keeps a connection pinned for N queries after an auto-increment insert to keep
  `LAST_INSERT_ID()` correct — **a non-obvious trap:** in write-heavy/sharded workloads the default
  5 can effectively defeat multiplexing (set to 0–1 to restore it).[^8][^12]
- **Read/write splitting.** Via `mysql_query_rules`, not a dedicated module. Three approaches:
  port-based (different `proxy_port`→different hostgroup); regex on query
  (`^SELECT.*FOR UPDATE`→primary, `^SELECT`→replicas); and the **recommended** digest-based
  intelligent routing (analyze `stats_mysql_query_digest`, target expensive `SELECT`s by digest to
  replicas). Docs warn generic `^SELECT`→replica rules break transactions, session state, and
  consistency.[^2][^4][^9]
- **Query caching (result cache).** Per-rule `cache_ttl` (in **milliseconds**; was seconds pre-1.1).
  Caches result sets on the wire; identical re-executions return immediately. Memory via
  `mysql-query_cache_size_MB` (default 256 MB, soft). **Limits:** no invalidation beyond TTL
  expiry, no LRU, incompatible with prepared statements. Cached hits show `hostgroup = -1` in the
  digest stats.[^4][^10]
- **Query rewriting.** `match_pattern`+`replace_pattern` rewrite SQL before execution (index hints,
  ORM fixes, inject `SQL_NO_CACHE`) without app changes.[^4]
- **Firewall / `mysql_firewall_whitelist`.** Introduced **2.0.9**. `mysql_firewall_whitelist_users`
  (per-user mode) + `mysql_firewall_whitelist_rules` (approved digests). Modes: **OFF** (allow all),
  **DETECTING** (allow but log non-whitelisted), **PROTECTING** (block non-whitelisted). Globally
  `mysql-firewall_whitelist_enabled=1`; activate with `LOAD MYSQL FIREWALL TO RUNTIME`. A bundled
  `libsqlinjection` engine + `mysql_firewall_whitelist_sqli_fingerprints` exists but has "quite a
  lot of false positives." **Practitioner warning:** configure rules *before* enabling PROTECTING
  or ProxySQL blocks all client connections.[^18][^19][^20][^21][^22]
- **Query mirroring.** `mirror_hostgroup`/`mirror_flagOUT` on a rule spawns a session that replays
  the (rewritten) query against a mirror hostgroup — useful to test a new replica/version with
  production traffic.[^23]

## MariaDB MaxScale architecture

Plugin-based proxy with a six-object model. **Note: MaxScale is BSL (Business Source License)
since 2.x — not fully OSS** (background knowledge, not re-verified this run).

- **Server** — a backend instance; monitored by exactly one monitor.
- **Monitor** — observes health/role and exposes it to routers (`mariadbmon` primary-replica,
  `galeramon` Galera).
- **Router** — `readwritesplit` (statement routing, each request routed individually),
  `readconnroute` (connection routing, server chosen at session start), `schemarouter` (sharding),
  `binlogrouter` (binlog relay).
- **Filter** — sits in front of routers in the request chain; order matters (`qlafilter`, `cache`,
  `ccrfilter`, `dbfwfilter`, `regexfilter`, `tee`).
- **Service** — binds a router + servers (+ filters) behind one logical DB.
- **Listener** — a port/protocol bound to a service.[^28][^29][^33][^34]

- **Query classifier (pluggable parser).** A central feature — MaxScale parses SQL via a query
  classifier (default `qc_sqlite`, a modified SQLite parser since 2.0) to report statement *type*
  (read→replica, write→primary) and operation. This is what enables **automatic** read/write
  split (vs ProxySQL's manual rules).[^38][^52]
- **`readwritesplit` router.** Splits reads (spread across replicas) from writes (primary), auto-
  detecting primary changes. **To primary:** explicit transactions, stored procs/functions/UDFs,
  temp tables, `EXECUTE` of prepared statements. **To a replica:** read-only `SELECT`, `SHOW` of
  vars. **To all servers:** `SET`, `USE`, `PREPARE`, client commands.[^31][^46]
- **`schemarouter` (sharding).** One logical server from many; database-based sharding — queries to
  a unique DB route to that DB's server; session-modifying queries route to all nodes.[^35]
- **`binlogrouter` (replication relay).** Pulls binary logs from a primary, stores them locally,
  and serves them to replicas as if it were a master — crucially it does **not re-execute**
  statements (byte-for-byte binlog cache), so replicas get 1:1 correlation with the primary and
  parallelism is preserved; lets the primary change without disconnecting replicas.[^36][^37]
- **`mariadbmon` failover/switchover/rejoin.** Auto-promote a replica on primary failure
  (`auto_failover=true`, after `failcount` passes), `switchover` (swap a running primary with a
  replica), `auto_rejoin` (return an old primary as a replica). A "primary change" only re-targets
  routing inside MaxScale; failover/switchover modify the actual cluster.[^30]
- **Cache filter.** `soft_ttl` (refresh threshold — first client after expiry triggers refresh) vs
  `hard_ttl` (discard threshold); storage `storage_inmemory` (default, non-persistent) /
  `storage_redis` / `storage_memcached`; `invalidate` (`never`/`current`) controls change-based
  invalidation.[^48]
- **`dbfwfilter` (SQL firewall).** Rules file; `action=allow` (whitelist), `block` (blacklist),
  `ignore`. Modes: *learning* (`learn-clear`/`learn-append`), *supervising* (warn), *enforcing*,
  *idle*. Rule syntax: `on_queries {select|update|…}`, `no_where_clause`, column rules. Does not
  support multi-statements. (The older standalone "Database Firewall filter" was deprecated in
  MaxScale 6.)[^47]
- **`maxctrl` admin.** CLI talking to the MaxScale **REST API** (default port **8989**, default
  creds `admin:mariadb`); replaced the removed `MaxAdmin`. `list`/`show`/`create`/`alter`/`destroy`
  at runtime (e.g. `maxctrl alter filter MyFirewall mode=enforce`).[^49][^50][^51]

## Read/write splitting & the consistency hazard

When reads route to asynchronous replicas, replication lag means a read right after a write may
hit a replica that hasn't applied that write → **stale / non-read-your-writes** results. MaxScale
docs are explicit that the coarse `max_slave_replication_lag` filter "does not guarantee that
writes done on the primary are visible for reads done on the replica."[^31][^39]

- **MaxScale `causal_reads`** (the feature that *does* guarantee causality): uses the server's
  `MASTER_GTID_WAIT()` + GTID session tracking. After a write, readwritesplit prefixes the next
  read with `MASTER_GTID_WAIT(...)` so the replica serves it only after applying that GTID. Modes:
  - **`local`** — affects only the writing connection; gives read-your-writes; trades latency for
    read scalability ("behaves as if you were on a single node").
  - **`global`** — writes globally visible to all service connections; slower.
  - **`fast`** — routes the read only to a server already known to have replicated (falls back to
    primary); improves latency at the cost of scalability.
  - **`universal`/`fast_universal`** — extend causality across any number of MaxScale instances and
    cover writes done outside MaxScale; highest cost ("last resort").[^31][^39][^40][^43]
- **MaxScale CCR filter (`ccrfilter`)** — coarser: on a write, force that connection's reads to the
  **primary** for `time` seconds and/or `count` statements; `global=true` propagates across
  connections.[^41]
- **ProxySQL** has **no single named causal-read mode**. Its closest analog is GTID-aware routing
  via the `gtid_from_hostgroup` query-rule column + `Queries_GTID_sync` metric (route a read only
  to a replica that reached the needed GTID). Operationally: keep transactions and post-write reads
  on the primary; the `auto_increment_delay_multiplex` pin helps for `LAST_INSERT_ID()`.[^2][^9]

**Version-stamped hazards (re-verify):**
- MaxScale **≤2.5.14** could produce non-causal reads with multi-domain GTIDs (fixed **2.5.15**);
  stored GTID coordinates never reset on a replication reset until **6.4.11** added `reset-gtid`.[^31]
- Open bug (single Jira source, 2024–25): `INSERT … RETURNING` can break `causal_reads` (the
  RETURNING result makes readwritesplit miss the new GTID), reproduced on 24.02.4 with
  `causal_reads=global`.[^42] **Tentative.**
- For **Galera**, the recommended causality lever is server-side `wsrep_sync_wait=1`, not
  `causal_reads` (which targets async MariaDB replication).[^43][^44]

## Operations: latency, HA, observability, anti-patterns

- **Latency overhead.** Percona measured a 46µs direct query at ~71µs through ProxySQL (~25µs
  added), more over TCP/IP than a Unix socket. Cross-proxy: MySQL Router highest latency; HAProxy
  best at low connection counts; **ProxySQL wins at high connection counts** (multiplexing) but
  uses **more CPU** (~2.5 cores vs HAProxy ~1.5) as a Layer-7 proxy — Percona explicitly warns it
  is "not optimal" inside a resource-constrained K8s operator. Poor config (too many/complex rules,
  default `mysql-threads=4` on many cores) inflates overhead.[^11][^14][^16]
- **Proxy-tier HA** (a single proxy is a SPOF):
  - **Sidecar / co-located** — proxy per app host, app connects to `127.0.0.1:6033`: zero extra
    hop, no shared SPOF; but **config drift** across instances and **connection creep** as sidecars
    autoscale are the documented failure modes.[^15][^17]
  - **Centralized cluster behind an L4 LB/NLB** — ProxySQL Cluster syncs config P2P via
    checksum/epoch; avoids connection creep; small added hop.[^2][^17]
  - **Keepalived + VRRP floating VIP** (active/passive) — the most-cited HA pattern. **MaxScale
    gotcha:** only the Keepalived-MASTER node may manipulate the cluster — set the others
    `maxctrl alter maxscale passive true` (passive nodes still route; they skip
    failover/switchover/rejoin).[^53][^54][^55]
- **Observability.** ProxySQL `stats` schema (port 6032): `stats_mysql_connection_pool` (per
  hostgroup+server status/ConnUsed/ConnFree/Queries/Latency_us), `stats_mysql_query_digest` (the
  read/write-split analysis source), `stats_mysql_query_rules` (per-rule hits),
  `stats_mysql_processlist`, `stats_mysql_query_cache`. **ProxySQL 2.1+ has a built-in Prometheus
  exporter** (no separate process; scrape `/metrics`; docs advise against sub-second scrape).
  MaxScale: `maxctrl list/show`, the REST API, and `qlafilter`.[^24][^25][^26]
- **Anti-patterns (synthesis):** naive `^SELECT`→replica split; default `mysql-threads=4` on many
  cores; leaving `auto_increment_delay_multiplex=5` in write-heavy workloads; sidecar config drift
  / connection creep; enabling the firewall before populating rules; TTL-only cache staleness;
  ProxySQL in a resource-constrained K8s operator; MaxScale multi-instance without `passive`;
  MaxScale `transaction_replay` duplicate-commit risk if a connection drops during `COMMIT`.[^2][^9][^45]

## ProxySQL vs MaxScale at a glance

| Dimension | ProxySQL | MariaDB MaxScale |
|---|---|---|
| Config | 3-layer SQLite (runtime/memory/disk), edited as SQL via admin 6032 | `maxscale.cnf` + runtime via `maxctrl`/REST API (8989) |
| RW split | Rule-based (manual/analyzed by digest) | Automatic via query classifier (`qc_sqlite`) |
| Causal reads | GTID sync via `gtid_from_hostgroup` (no named mode) | `causal_reads` (local/global/fast/universal); `ccrfilter` |
| Multiplexing | Yes (core differentiator) | Connection pooling; readwritesplit restores session state |
| Result cache | `cache_ttl` per rule (TTL-only, no invalidation) | `cache` filter (soft/hard TTL, optional invalidation, redis/memcached) |
| Firewall | `mysql_firewall_whitelist` + SQLi engine | `dbfwfilter` (allow/block, learn/supervise/enforce) |
| Sharding | (not a core feature) | `schemarouter` (DB-based) |
| License | Open source (GPL) | BSL since 2.x |
| Prometheus | Built-in exporter (2.1+) | REST API + qlafilter |

## Sources

[^1]: ProxySQL — Architecture Overview. https://proxysql.com/documentation/architecture/
[^2]: ProxySQL — Multi-Layer Configuration System. https://proxysql.com/documentation/main-runtime/multi-layer-configuration/
[^3]: sysown/proxysql — doc/configuration_system.md. https://github.com/sysown/proxysql/blob/master/doc/configuration_system.md
[^4]: ProxySQL — Configure ProxySQL for MySQL (mysql_query_rules). https://proxysql.com/documentation/proxysql-configuration/
[^8]: ProxySQL — Multiplexing. https://proxysql.com/documentation/multiplexing/
[^9]: ProxySQL — How to Set Up Read/Write Split. https://proxysql.com/documentation/proxysql-read-write-split-howto/
[^10]: ProxySQL — Query Cache. https://proxysql.com/documentation/query-cache/
[^11]: Percona — ProxySQL Overhead Explained and Measured. https://www.percona.com/blog/proxysql-overhead-explained-and-measured/
[^12]: The ProxySQL multiplexing wild goose chase. https://mysqlquicksand.wordpress.com/2019/11/28/the-proxysql-multiplexing-wild-goose-chase/
[^13]: Instawork Engineering — Overcoming database bottlenecks with ProxySQL. https://engineering.instawork.com/overcoming-database-bottlenecks-a-journey-with-proxysql-4bf525256092
[^14]: Percona — Comparisons of Proxies for MySQL. https://www.percona.com/blog/comparisons-of-proxies-for-mysql/
[^15]: Egnyte — How Egnyte Achieves MySQL High Availability. https://www.egnyte.com/blog/post/how-egnyte-achieves-mysql-high-availability
[^16]: sysown/proxysql Issue #4652 — High latency increasing with load. https://github.com/sysown/proxysql/issues/4652
[^17]: nuffing.com — ProxySQL in Front of AWS RDS & Aurora MySQL (Part 1). https://nuffing.coutinho.net/2026/05/proxysql-in-front-of-aws-rds-aurora-mysql-part-1-why-and-where-to-place-it/
[^18]: ProxySQL — Firewall Whitelist. https://proxysql.com/documentation/firewall-whitelist/
[^19]: ProxySQL — Security. https://proxysql.com/documentation/security/
[^20]: Percona — ProxySQL 2.0.9 Introduces Firewall Whitelist. https://www.percona.com/blog/proxysql-2-0-9-introduces-firewall-whitelist-capabilities/
[^21]: Severalnines — Protect MySQL/MariaDB from SQL Injection Part Two. https://severalnines.com/blog/how-protect-your-mysql-or-mariadb-database-sql-injection-part-two/
[^22]: Mydbops — Building a MySQL Firewall with ProxySQL. https://www.mydbops.com/blog/building-a-mysql-firewall-with-proxysql
[^23]: sysown/proxysql — doc/mirroring.md. https://github.com/sysown/proxysql/blob/master/doc/mirroring.md
[^24]: ProxySQL — MySQL Stats Tables. https://proxysql.com/documentation/the-admin-schemas/stats/stats-mysql/
[^25]: ProxySQL — Prometheus Exporter. https://proxysql.com/documentation/prometheus-exporter/
[^26]: ProxySQL Blog — Observability Enhancements in 2.1 (Prometheus & Grafana). https://proxysql.com/blog/observability-enhancements-in-proxysql-2-1-with-prometheus-grafana/
[^28]: MariaDB — MaxScale Architecture. https://mariadb.com/docs/maxscale/maxscale-architecture
[^29]: MariaDB — MaxScale Configuration Guide. https://mariadb.com/docs/maxscale/maxscale-management/deployment/installation-and-configuration/maxscale-configuration-guide
[^30]: MariaDB — MariaDB Monitor (mariadbmon). https://mariadb.com/docs/maxscale/reference/maxscale-monitors/mariadb-monitor
[^31]: MariaDB — MaxScale Readwritesplit. https://mariadb.com/docs/maxscale/reference/maxscale-routers/maxscale-readwritesplit
[^33]: MariaDB — MaxScale Reference (routers/filters/monitors index). https://mariadb.com/docs/maxscale/reference
[^34]: MariaDB — MaxScale Servers. https://mariadb.com/docs/maxscale/reference/maxscale-servers.md
[^35]: MariaDB — MaxScale SchemaRouter. https://mariadb.com/docs/maxscale/reference/maxscale-routers/maxscale-schemarouter
[^36]: MariaDB — MaxScale Binlogrouter. https://mariadb.com/docs/maxscale/reference/maxscale-routers/maxscale-binlogrouter
[^37]: MaxScale-Documentation — Replication-Proxy-Binlog-Router-Tutorial. https://github.com/mariadb-corporation/MaxScale-Documentation/blob/2.0.0/Tutorials/Replication-Proxy-Binlog-Router-Tutorial.md
[^38]: MariaDB — Query Classification and Pluggable Parser. https://mariadb.com/resources/blog/query-classification-and-pluggable-parser/
[^39]: MariaDB — Ensuring Causal Consistency with MaxScale RW Split Router. https://mariadb.com/docs/maxscale/maxscale-use-cases/readwrite-split-router-usage/ensuring-causal-consistency-with-maxscales-readwrite-split-router
[^40]: MariaDB.org — Read Causality (Fest 2024 talk). https://mariadb.org/fest-2024-berlin/read-causality/
[^41]: MariaDB — MaxScale Consistent Critical Read (CCR) Filter. https://mariadb.com/docs/maxscale/reference/maxscale-filters/maxscale-consistent-critical-read-filter
[^42]: jira.mariadb.org MXS-5527 — INSERT…RETURNING breaks causal_reads. https://jira.mariadb.org/browse/MXS-5527
[^43]: DBA StackExchange — disable MaxScale RW router for specific queries (causal_reads/wsrep_sync_wait). https://dba.stackexchange.com/questions/343165/
[^44]: MaxScale commit a2e55e5 — MXS-5588 Fix GTID parsing/probing. https://github.com/mariadb-corporation/MaxScale/commit/a2e55e5040de7057358ea780439ae3ad51f859de
[^45]: MariaDB — MaxScale Limitations and Known Issues. https://mariadb.com/docs/maxscale/maxscale-management/mariadb-maxscale-limitations-guide
[^46]: MariaDB — Routing Statements with the RW Split Router. https://mariadb.com/docs/maxscale/maxscale-use-cases/readwrite-split-router-usage/routing-statements-with-maxscales-readwrite-split-router
[^47]: MariaDB — MaxScale Firewall Filter (dbfwfilter). https://mariadb.com/docs/maxscale/reference/maxscale-filters/maxscale-firewall-filter
[^48]: MariaDB — MaxScale Cache filter. https://mariadb.com/docs/maxscale/reference/maxscale-filters/maxscale-cache
[^49]: MariaDB — MaxCtrl reference. https://mariadb.com/docs/maxscale/reference/maxscale-maxctrl.md
[^50]: MariaDB — REST API Tutorial. https://mariadb.com/docs/maxscale/mariadb-maxscale-tutorials/rest-api-tutorial
[^51]: MaxScale — Administration-Tutorial.md (24.02). https://github.com/mariadb-corporation/MaxScale/blob/24.02/Documentation/Tutorials/Administration-Tutorial.md
[^52]: MaxScale — Configuration-Guide.md (passive, query_classifier). https://github.com/mariadb-corporation/MaxScale/blob/24.02/Documentation/Getting-Started/Configuration-Guide.md
[^53]: Mydbops — Making MaxScale/ProxySQL Highly Available (Keepalived). https://www.mydbops.com/blog/making-maxscale-proxysql-highly-available-2-1
[^54]: SolusiDB — Setup MaxScale HA using Keepalived and Maxctrl. https://www.solusidb.com/2019/03/08/setup-maxscale-ha-using-keepalived-and-maxctrl/
[^55]: ProxySQL — Advanced Query Routing (features). https://proxysql.com/features/query-routing/
