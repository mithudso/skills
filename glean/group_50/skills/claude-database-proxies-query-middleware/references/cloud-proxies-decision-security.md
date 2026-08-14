# Cloud/Commercial Proxies, Decision Framework & Security

> Provenance: reference under the standalone `database-proxies-query-middleware` skill. Created
> 2026-06-18 via `/dr` deep research. `verified-as-of: 2026-06-18` — pricing, support, and
> version facts are volatile; re-verify against current vendor docs before relying on them.

## Contents
- [AWS RDS Proxy](#aws-rds-proxy)
- [Heimdall Data](#heimdall-data)
- [Decision framework: proxy vs driver vs application](#decision-framework-proxy-vs-driver-vs-application)
- [Security: SQL firewall & credential brokering](#security-sql-firewall--credential-brokering)
- [Anti-patterns / when NOT to use a proxy](#anti-patterns--when-not-to-use-a-proxy)
- [Sources](#sources)

---

## AWS RDS Proxy

A fully managed, highly available proxy for RDS/Aurora that pools and multiplexes connections,
hardens auth, and accelerates failover.[^1][^2][^4]

- **What it solves.** Holds a warm pool of long-lived backend connections and multiplexes many
  client connections onto a smaller backend set ("transaction-level reuse is called
  multiplexing"). Purpose-built for the **serverless/Lambda connection explosion**: Lambda
  micro-VMs can't share a traditional pool, so thousands of concurrent invocations exhaust the DB;
  RDS Proxy externalizes the pool. It also absorbs bursts (queue/throttle rather than "too many
  connections"). No app code change — repoint the endpoint.[^1][^2][^3][^6]
- **Pinning.** When the proxy "can't be sure it's safe to reuse a database connection outside of
  the current session," it **pins** the client to one backend for the session — reverting to 1:1
  and "reducing the effectiveness of connection reuse." Heavy pinning makes the proxy "essentially
  not proxying."[^5][^6][^13] **Documented pinning conditions (load-bearing):**
  - *All engines:* statement text > **16 KB** pins.[^5]
  - *PostgreSQL:* `SET`; `PREPARE`/`DISCARD`/`DEALLOCATE`/`EXECUTE`; temp sequences/tables/views;
    cursors; `LISTEN`; loading a module (e.g. `auto_explain`); sequence functions; **session**-level
    advisory locks (`pg_advisory_lock`) — but NOT transaction-scoped (`pg_advisory_xact_lock`).[^5][^11]
  - *MySQL/MariaDB:* `LOCK TABLE(S)`, `FLUSH TABLES WITH READ LOCK`; named locks (`GET_LOCK`); temp
    tables; prepared statements; executable comments (`/*! */`). `SET LOCAL` and stored-proc calls
    do NOT pin; some `SET` vars are *tracked* (reused for matching sessions) rather than pinned.[^5][^11]
  - *SQL Server:* MARS; DTC; temp tables/cursors/prepared statements; a list of `SET` statements.[^5]
  - **Detect/avoid:** watch CloudWatch `DatabaseConnectionsCurrentlySessionPinned` + proxy logs
    (each pin logged with a reason); push common `SET`s to the **initialization query**; for
    **MySQL only** use the `EXCLUDE_VARIABLE_SETS` session pinning filter (**not for PostgreSQL** —
    set vars DB-side there); disable ORM prepared statements; keep a direct non-proxy connection for
    session-state-dependent workloads.[^5][^7][^12][^14][^16]
- **Failover acceleration.** During failover the proxy holds client connections open, accepts at
  the same IP, queues in-flight requests, and reconnects to the new writer — bypassing DNS/CNAME
  caching; only mid-transaction connections are canceled. Headline **"up to 66%"** reduction
  (product page/FAQ/docs). Blog-specific (2020): up to 79% Aurora MySQL, 32% RDS MySQL Multi-AZ.
  **Qualifier (don't overstate):** AWS itself says "during failovers the application may experience
  increased latencies and ongoing transactions may have to be re-tried." Klaviyo (2025) found it
  "deficient" at failover in their environment yet kept it for reduced ops burden.[^1][^2][^8][^10][^13][^22]
- **Security/auth.** DB credentials, standard IAM, or end-to-end IAM. **Standard IAM:** client
  presents an IAM token; the proxy fetches the DB username/password from **Secrets Manager** ("RDS
  Proxy always connects to the database using password authentication through Secrets Manager" in
  this mode). **End-to-end IAM** (GA Sept 2025, MySQL/PostgreSQL): IAM on both hops, no Secrets
  Manager secret. **IAM auth requires TLS.** Up to **200 secrets per proxy**.[^24][^25][^26][^27]
- **Supported engines (as of 2026-06-18).** Aurora MySQL/PostgreSQL, RDS for MySQL, PostgreSQL,
  MariaDB, SQL Server. **Not** Oracle or Db2.[^19][^24] PostgreSQL protocol v3.0 only; SQL Server
  2016/2017/2019 (not 2014/2022).[^19][^30]
- **Limitations/costs.** Priced **per vCPU-hour of the underlying DB instance** (per-ACU-hour for
  Aurora Serverless), billed even when idle, no Reserved pricing. Latency overhead ~**1–5 ms/query**
  (negligible for 20 ms web queries; "2 ms on a 0.5 ms query = 30% increase"). Same-VPC only, not
  public; needs ≥2 AZ subnets and up to 215 IPs (IP exhaustion blocks scaling/patching, event
  RDS-EVENT-0243); 20 proxies/account/Region (soft). Real "not worth it" reverts: a default
  `SET search_path` pinned everything (12-hr rollback); Prisma `COM_STMT_PREPARE` pinned hundreds of
  connections.[^14][^15][^28][^29][^30][^33]

## Heimdall Data

**Credibility framing up front:** almost all detailed material is **vendor** (heimdalldata.com),
**co-marketing** (AWS blogs co-authored by Heimdall's CTO), or a **sponsored** analyst paper. The
only near-independent signals are AWS Marketplace reviews. No neutral benchmark / consistency audit
was found.[^h1][^h3][^h4]

- **What it is.** A commercial, database-neutral proxy / "intelligent data access layer" between
  app and DB, positioned as a query-optimization/SQL-offload proxy; no app code changes; deployable
  as sidecar, separate tier, or on the DB. The "proxy between app and DB" categorization is
  corroborated by AWS and third-party directories.[^h1][^h3][^h4]
- **Key feature — auto query caching + auto invalidation (the headline differentiator).** Caches
  result sets to in-heap L1 or look-aside L2 (Redis/ElastiCache); claims to auto-decide what to
  cache and auto-invalidate on data change with no code changes; mechanism is **table-aware** (a
  detected DML puts the table in a ~2-second invalidation window).[^h1][^h5][^h6] **Credibility
  note:** the *capability class* (table-level auto-invalidation) is independently real (MaxScale
  documents the same), but Heimdall's specific **"guaranteed/strong (ACID) consistency, better than
  ProxySQL/RDS Proxy"** framing lives in the **sponsored** paper + vendor site — VENDOR-CLAIM, no
  independent corroboration. Heimdall's own docs undercut "fully automatic": caching silently won't
  apply inside transactions, on queries with no detectable table, on system/`pg_*`/temp tables, or
  on read/update table-name mismatches, and NTP clock skew causes cache back-off.[^h1][^h5][^h6][^h7]
- **Also:** read/write split with replication-lag detection; connection pooling + multiplexing;
  multi-DB (PostgreSQL, MySQL, SQL Server, Aurora/Redshift, RDS). Sold on AWS/Azure/GCP/OCI
  marketplaces as AMIs, commercial subscription.[^h3][^h10][^h16]

## Decision framework: proxy vs driver vs application

**Core principle:** these are **layers with distinct jobs, not substitutes** — "the right number of
pooling layers is the smallest number that makes your problem go away." Most teams need one (the
app/driver pool), many need two (+ proxy), some three.[^d1][^d3][^d5][^d6]

- **Step 0 — exhaust cheaper options first.** Right-size the **in-driver pool** (HikariCP
  `(cores×2)+spindles` → ~10–20) and **fix slow queries** before adding any tier. A bigger pool
  usually loses throughput (a pool of 200 = "44% throughput reduction with 10× tail latency").
  Raising `max_connections` to mask the problem "will crash your server faster."[^d2][^d10][^d11][^d14]
- **Stay driver-only when:** few app instances; total connections (pool × processes) well under the
  DB limit; no serverless/autoscaling; short queries. **Driver-level read/write split** (MySQL
  Connector/J `replication://`, the AWS JDBC driver) is fine when routing is **simple and static** —
  but the policy "lives in application code… no central audit trail," and the driver won't
  auto-split (the app must mark a connection read-only). **Hard limit:** in serverless the driver
  pool structurally **cannot** work (no persistent process across invocations).[^d3][^d5][^d6][^d9]
- **Own it in the application when:** the concern is read offload with tolerable staleness →
  **cache-aside** Redis/Memcached is the AWS-recommended default (tolerates cache unavailability by
  falling through to the DB). NOT for payments/balances/read-your-writes. App-side sharding is
  defensible only for stable/few-shard cases; the weight of evidence says move routing **out** of
  the app into a transparent proxy (Vitess) once it grows.[^d21][^d22][^d23][^d17][^d25]
- **Add a PROXY/MIDDLEWARE tier when any trigger fires:** serverless/Lambda (near-mandatory);
  horizontal scale-out past the connection budget / "too many clients" / connection storms; need
  code-free centralized read/write split, sharding, query rules, mirroring, or caching
  (ProxySQL/Vitess/MaxScale for *query-aware*, RDS Proxy/connector-side for *basic*); tight failover
  RTO + IAM/rotation as first-class.[^d3][^d5][^d6][^d9][^d12][^d16]
- **What the proxy hop COSTS:**
  1. **Latency** ~1–5 ms/query in-VPC (sidecar/loopback ~0; cross-AZ 1–2 ms; separate host
     10–15 ms). Negligible vs 10–50 ms OLTP; **pure overhead on long analytical queries** — "a
     reporting query that already takes 8 seconds does not benefit from pooling; direct connections
     are correct."[^d5][^d6][^d15][^d17]
  2. **An extra HA tier** to size/scale/patch/monitor; managed proxies shift but don't remove this
     (RDS Proxy IP exhaustion can block AWS patching; 20-proxy limit; pinning).[^d6][^d15][^d16][^d26]
  3. **SPOF if not HA** — "a single instance of HAProxy or PgBouncer simply moves your single point
     of failure from the database to the proxy layer."[^d28][^d29]
  4. **Semantic gotchas** — transaction-mode pooling (the mode that pays) breaks session `SET`,
     server-side prepared statements, `LISTEN`/`NOTIFY`, advisory locks, temp tables, cursors.[^d1][^d3][^d5]
- **Managed vs self-hosted:** RDS Proxy for low-ops basic pooling/failover/IAM; self-hosted
  ProxySQL/PgBouncer/PgCat/MaxScale for fine tuning, query rules, mirroring, or non-RDS DBs. RDS
  Proxy bills per-vCPU regardless of throughput — "2× cost for a benefit you don't need" on a small
  fleet, but "the only thing that keeps the database alive" for a Lambda fleet.[^d1][^d5][^d6][^d16]

## Security: SQL firewall & credential brokering

### Pattern A — SQL firewall (allow/deny query digests)
An inline proxy firewall decides per-query whether to forward/block/warn, resting on **query
normalization/fingerprinting** (strip literals → a stable "digest"). Record the finite set of
digests the app legitimately emits; block anything else (an injected `OR 1=1`/`UNION SELECT`/stacked
`DROP` produces a new digest and is rejected). Lifecycle: **Train → Detect (log only) → Protect
(enforce)**.[^s4][^s11][^s13][^s14][^s22] **OWASP framing:** parameterized statements are the
*primary* SQLi defense; allow-listing is defense-in-depth/secondary, not a replacement.[^s18][^s19][^s20]
- **ProxySQL** `mysql_firewall_whitelist_*` (OFF/DETECTING/PROTECTING; per-user digest allow-lists;
  a bundled `libsqlinjection` engine to suppress false positives). Prefer precomputed `digest` over
  regex (far more efficient).[^s1][^s3][^s5][^s6][^s7][^s21]
- **MaxScale `dbfwfilter`** ("similar to iptables"; `allow`/`block`/`ignore`; column/where/function
  rules; learn/supervise/enforce modes persisting allowed canonical statements; no multi-statements).[^s2][^s24][^s25][^s26]
- **DISCONFIRMING:** parser-differential bypasses exist (wrap a payload in `CREATE PROCEDURE` or
  versioned executable comments so the proxy parser and backend disagree); digest bugs (PostgreSQL
  digest truncating complex queries); false positives / app breakage (any new/altered query is a new
  digest → blocked in enforce mode; enabling before training "blocks all client connections"); the
  bundled SQLi engine is "far from perfect… quite a lot of false positives."[^s12][^s14][^s16][^s17][^s22][^s26][^s27]

### Pattern B — Credential brokering
The app authenticates to the proxy (e.g. a short-lived IAM token); the proxy authenticates to the
DB with credentials it holds (RDS Proxy fetches the DB password from Secrets Manager, or uses
end-to-end IAM). **Benefits:** no DB passwords in app code/config; centralized rotation; auditable
IAM-governed access (`rds-db:connect`); offloaded auth/connection-setup cost. The topology
generalizes beyond AWS (ProxySQL also holds DB creds and relays the auth handshake).[^s9][^s28][^s30][^s31][^s33][^s34]
- **DISCONFIRMING:** a **concentrated trust point** — in standard mode the proxy can fetch real DB
  passwords from Secrets Manager, so compromise of the proxy's IAM role ≈ DB-credential compromise
  (end-to-end IAM exists precisely to remove the stored password); standard mode **cannot fully
  eliminate the password**; an IAM-auth rate ceiling (~200 new conns/sec, then fall back to password
  mode); end-to-end IAM is MySQL/PostgreSQL only.[^s28][^s29][^s32][^s9]
- **Cross-cutting:** Patterns A and B are orthogonal/composable — ProxySQL does both; RDS Proxy is
  primarily a Pattern-B broker/pooler and ships **no** query-digest firewall. Neither replaces
  in-app defenses.[^s18][^s28][^s32]

## Anti-patterns / when NOT to use a proxy

The anti-pattern is usually **misuse, premature adoption, or un-mitigated deployment** — not the
proxy concept.

- **5.1 Hidden SPOF.** "If the proxy goes down, every application behind it loses database access…
  the central risk of the architecture." New failure modes even when HA: failover **thundering
  herd**; pool **saturation while the DB is healthy** (PlanetScale incident); a proxy outage
  **poisoning the client pool** so it can't self-recover until restart (n8n #30612). **Pushback:** a
  properly-HA proxy *improves* failover resilience.[^a1][^a2][^a5][^a6][^a7]
- **5.2 Transaction pooling breaks session semantics.** Marked "Never" in transaction mode:
  `SET`/`RESET`, `LISTEN`, `WITH HOLD` cursors, `PREPARE`/`DEALLOCATE`, temp tables, `LOAD`, session
  advisory locks. Concrete breakage: `SET statement_timeout` "swapped out from my web request,
  picked up by my job"; `search_path` poisoning via `pg_dump`; multi-tenant `SET app.tenant_id` (RLS)
  colliding across requests; `prepared statement "…" does not exist`; advisory locks held forever on
  the wrong connection. **The failure is silent and non-deterministic.** Same problem drives RDS
  Proxy pinning. Workarounds: `SET LOCAL`; proxy init query; protocol prepared-statement support in
  PgBouncer ≥1.21 / Supavisor / pgcat; or session mode (gives up the win).[^a9][^a10][^a11][^a12][^a14][^a15][^a20]
- **5.3 Caching staleness / invalidation hazards.** Meta's "Cache made consistent": "a dynamic cache
  produces race conditions beyond your wildest imagination." Bug classes: forgotten write path,
  mismatched serialized representations (`"42"` vs `42`), silent fire-and-forget invalidation
  failures; thundering herd / cache stampede on key expiry; null-key cache penetration as a DB DDoS.
  **Pushback:** consistency at scale is achievable (Meta's Polaris versioned cache objects against
  the MySQL WAL LSN) — but it took a purpose-built system; don't cache data that must be
  fresh.[^a22][^a23][^a24][^a25][^a27][^a28]
- **5.4 Masking a real scaling/schema problem.** "A slow query holds a pooled connection for its
  whole duration… fix the slow query and the pool pressure disappears without changing any pool
  config." N+1: `Promise.all` "doesn't eliminate them — more concurrent queries means more pool
  pressure." A **caching** proxy masks even more potently — the canonical story: an unindexed query
  fronted with Redis (800ms→5ms), six months later a migration made it 3s "but nobody noticed
  because the cache absorbed 99.8% of traffic… the cache had masked the real problem so effectively
  the team lost awareness of it." Prescribed order: measure → fix queries/indexes/N+1 → fix schema →
  tune → *then* cache.[^a29][^a31][^a32][^a33][^a24][^a30]
- **5.5 Operational/observability burden.** "Another process to deploy, configure, monitor, and
  upgrade… can become a bottleneck" (PgBouncer's single-thread = one-core ceiling; ProxySQL per-node
  config drift). Latency: sub-1 ms typical on-network, but an AWS re:Post case reported
  Lambda→RDS Proxy→RDS "8–12× slower" for some, with AWS advising low-frequency Lambdas to consider
  direct connections. **Observability distortion:** a multiplexing proxy makes `pg_stat_activity` no
  longer map to real clients, breaking standard connection-storm triage. **Pushback:** a proxy "sees
  every query… a natural point to collect query stats" — it *relocates* observability rather than
  losing it.[^a1][^a11][^a2][^a35][^a37][^a39]
- **Real incidents worth citing:** OpenAI Feb-2023 DB outage (PgBouncer pools clogged by slow
  queries during recovery; new replicas bypassed PgBouncer and exceeded connection limits → a second
  outage); PlanetScale shared-pool saturation; n8n proxy-outage pool congestion.[^a-openai][^a5][^a6]

## Sources

### RDS Proxy
[^1]: AWS — Amazon RDS Proxy product page. https://aws.amazon.com/rds/proxy/
[^2]: AWS — Using Amazon RDS Proxy (UserGuide). https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
[^3]: AWS Compute Blog — Using RDS Proxy with Lambda. https://aws.amazon.com/blogs/compute/using-amazon-rds-proxy-with-aws-lambda/
[^4]: AWS — RDS Proxy concepts (UserGuide). https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
[^5]: AWS — Managing connections with RDS Proxy (pinning). https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-managing.html
[^6]: AWS — RDS Proxy overview / multiplexing. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
[^7]: AWS — RDS Proxy session pinning filters. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-managing.html
[^8]: AWS — RDS Proxy failover (UserGuide). https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
[^10]: AWS — RDS Proxy FAQ (failover reduction). https://aws.amazon.com/rds/proxy/faqs/
[^11]: AWS — RDS Proxy pinning conditions per engine. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-managing.html
[^12]: AWS — Avoiding pinning. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-managing.html
[^13]: Klaviyo Engineering — RDS Proxy evaluation (2025). https://klaviyo.tech/
[^14]: Jeremy Daly / practitioner — RDS Proxy pinning & cost notes. https://www.jeremydaly.com/
[^15]: Prisma issue — RDS Proxy COM_STMT_PREPARE pinning. https://github.com/prisma/prisma/issues
[^16]: AWS — RDS Proxy setup & initialization query. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-setup.html
[^19]: AWS — RDS Proxy supported engines/versions. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
[^22]: AWS Database Blog — RDS Proxy failover benchmarks (2020). https://aws.amazon.com/blogs/database/
[^24]: AWS — IAM/Secrets Manager with RDS Proxy. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-setup.html
[^25]: AWS — RDS Proxy Secrets Manager. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-setup.html
[^26]: AWS — RDS Proxy end-to-end IAM (2025). https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-setup.html
[^27]: AWS — RDS Proxy IAM auth (TLS required). https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
[^28]: AWS — RDS Proxy pricing. https://aws.amazon.com/rds/proxy/pricing/
[^29]: AWS — RDS Proxy pricing (Aurora Serverless per-ACU). https://aws.amazon.com/rds/proxy/pricing/
[^30]: AWS — RDS Proxy quotas/limits. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
[^33]: AWS — RDS Proxy IP/subnet exhaustion (RDS-EVENT-0243). https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_Events.Messages.html

### Heimdall
[^h1]: Heimdall Data — product / data sheet. https://www.heimdalldata.com/
[^h3]: AWS Database Blog — Heimdall Data co-marketing. https://aws.amazon.com/blogs/database/
[^h4]: Tracxn / company directory — Heimdall Data. https://tracxn.com/
[^h5]: Heimdall Data — caching docs. https://www.heimdalldata.com/docs/
[^h6]: Heimdall Data — auto-invalidation window docs. https://www.heimdalldata.com/docs/
[^h7]: Heimdall Data — caching limitations docs. https://www.heimdalldata.com/docs/
[^h10]: AWS Marketplace — Heimdall Data listing. https://aws.amazon.com/marketplace/
[^h16]: Heimdall Data — supported databases. https://www.heimdalldata.com/

### Decision framework
[^d1]: PGConf.EU 2024 — Comparing Connection Poolers (layers, the smallest number). https://www.postgresql.eu/events/pgconfeu2024/sessions/session/5846/slides/547/comparing_poolers.pdf
[^d2]: HikariCP wiki — About Pool Sizing. https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing
[^d3]: AWS — When to use RDS Proxy / driver vs proxy. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
[^d5]: Brandur — Connection pooling layers (transaction mode, long queries). https://brandur.org/
[^d6]: AWS Compute Blog — RDS Proxy with Lambda (serverless requires external pool). https://aws.amazon.com/blogs/compute/using-amazon-rds-proxy-with-aws-lambda/
[^d9]: AWS JDBC driver / MySQL Connector/J read-write splitting. https://dev.mysql.com/doc/connector-j/en/connector-j-master-slave-replication-connection.html
[^d10]: HikariCP wiki — pool-sizing benchmark. https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing
[^d11]: Gold Lapel — HikariCP Pool Sizing for Postgres. https://goldlapel.com/grounds/connection-pooling/hikaricp-pool-sizing-postgres
[^d12]: AWS — RDS Proxy use cases. https://aws.amazon.com/rds/proxy/
[^d14]: PostgreSQL wiki — Number Of Database Connections. https://wiki.postgresql.org/wiki/Number_Of_Database_Connections
[^d15]: Percona — ProxySQL Overhead Explained and Measured. https://www.percona.com/blog/proxysql-overhead-explained-and-measured/
[^d16]: AWS — RDS Proxy limits/cost considerations. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
[^d17]: PlanetScale — Sharding & moving routing out of the app (Vitess). https://planetscale.com/blog/connection-pooling
[^d21]: AWS — Caching strategies (cache-aside / lazy loading). https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/Strategies.html
[^d22]: AWS — Caching best practices. https://aws.amazon.com/caching/best-practices/
[^d23]: AWS — Lazy loading vs write-through. https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/Strategies.html
[^d25]: Shopify Engineering — sharding logic out of app code. https://shopify.engineering/
[^d26]: AWS — RDS Proxy IP exhaustion impact on patching. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
[^d28]: PgBouncer/HAProxy SPOF discussion (PGConf.EU 2024 slides). https://www.postgresql.eu/events/pgconfeu2024/sessions/session/5846/slides/547/comparing_poolers.pdf
[^d29]: Severalnines — HAProxy/PgBouncer HA. https://severalnines.com/blog/

### Security
[^s1]: ProxySQL — Firewall Whitelist. https://proxysql.com/documentation/firewall-whitelist/
[^s2]: MariaDB — MaxScale Firewall Filter (dbfwfilter). https://mariadb.com/docs/maxscale/reference/maxscale-filters/maxscale-firewall-filter
[^s3]: ProxySQL — Security. https://proxysql.com/documentation/security/
[^s4]: Percona — query digests / firewall rules. https://www.percona.com/blog/proxysql-2-0-9-introduces-firewall-whitelist-capabilities/
[^s5]: ProxySQL — mysql_firewall_whitelist_rules (digests). https://proxysql.com/documentation/firewall-whitelist/
[^s6]: ProxySQL — firewall modes. https://proxysql.com/documentation/firewall-whitelist/
[^s7]: ProxySQL — SQLi fingerprints table. https://proxysql.com/documentation/firewall-whitelist/
[^s9]: AWS — IAM DB auth with RDS Proxy (credential brokering). https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-setup.html
[^s11]: Generic SQL-firewall fingerprinting (digest normalization). https://www.percona.com/blog/proxysql-2-0-9-introduces-firewall-whitelist-capabilities/
[^s12]: ProxySQL issue — PostgreSQL digest truncation. https://github.com/sysown/proxysql/issues
[^s13]: ProxySQL — train/detect/protect lifecycle. https://proxysql.com/documentation/firewall-whitelist/
[^s14]: Academic — automatic whitelist generation for ORM apps (burden). https://dl.acm.org/
[^s16]: Severalnines — SQLi detection false positives. https://severalnines.com/blog/how-protect-your-mysql-or-mariadb-database-sql-injection-part-two/
[^s17]: OWASP — WAF/normalization bypasses (inline comments). https://owasp.org/
[^s18]: OWASP — SQL Injection Prevention Cheat Sheet. https://cheatsheetseries.owasp.org/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.html
[^s19]: OWASP — parameterized statements primary. https://cheatsheetseries.owasp.org/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.html
[^s20]: OWASP — allow-list as defense-in-depth. https://cheatsheetseries.owasp.org/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.html
[^s21]: Percona — prefer digest over regex. https://www.percona.com/blog/proxysql-2-0-9-introduces-firewall-whitelist-capabilities/
[^s22]: Mydbops — Building a MySQL Firewall with ProxySQL (train first). https://www.mydbops.com/blog/building-a-mysql-firewall-with-proxysql
[^s24]: MariaDB — dbfwfilter action types. https://mariadb.com/docs/maxscale/reference/maxscale-filters/maxscale-firewall-filter
[^s25]: MariaDB — dbfwfilter rule syntax. https://mariadb.com/docs/maxscale/reference/maxscale-filters/maxscale-firewall-filter
[^s26]: MaxScale dbfwfilter — parser-differential bypass / multi-statement limit. https://mariadb.com/docs/maxscale/reference/maxscale-filters/maxscale-firewall-filter
[^s27]: ProxySQL — firewall false positives in app code paths. https://proxysql.com/documentation/firewall-whitelist/
[^s28]: AWS — RDS Proxy end-to-end IAM (removes stored password). https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-setup.html
[^s29]: AWS — RDS Proxy standard IAM mode (Secrets Manager required). https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-setup.html
[^s30]: AWS — RDS Proxy credential management. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
[^s31]: AWS — rds-db:connect IAM policy. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html
[^s32]: AWS — IAM-auth rate guidance (~200 conns/sec). https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html
[^s33]: ProxySQL — backend credentials. https://proxysql.com/documentation/security/
[^s34]: AWS — centralized rotation benefits. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html

### Anti-patterns
[^a1]: PGConf.EU 2024 — pooler SPOF / operational burden. https://www.postgresql.eu/events/pgconfeu2024/sessions/session/5846/slides/547/comparing_poolers.pdf
[^a2]: nuffing.com — ProxySQL config drift / placement. https://nuffing.coutinho.net/2026/05/proxysql-in-front-of-aws-rds-aurora-mysql-part-1-why-and-where-to-place-it/
[^a5]: PlanetScale — shared-pool saturation incident. https://planetscale.com/blog/connection-pooling
[^a6]: n8n issue #30612 — proxy outage poisons client pool. https://github.com/n8n-io/n8n/issues/30612
[^a7]: AWS — RDS Proxy failover resilience (pushback). https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
[^a9]: PgBouncer — Features (transaction-mode "Never" matrix). https://www.pgbouncer.org/features.html
[^a10]: PgBouncer — FAQ. https://www.pgbouncer.org/faq.html
[^a11]: Stack Harbor KB — silent non-deterministic breakage. https://stackharbor.com/en/knowledge-base/pgbouncer-transaction-pool-advanced/
[^a12]: Heroku Dev Center — PgBouncer best practices. https://devcenter.heroku.com/articles/best-practices-pgbouncer-configuration
[^a14]: practitioner — search_path poisoning / RLS tenant collision. https://devcenter.heroku.com/articles/best-practices-pgbouncer-configuration
[^a15]: AWS — RDS Proxy pinning as the same problem. https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-managing.html
[^a20]: PgBouncer changelog — protocol prepared statements 1.21+. https://www.pgbouncer.org/changelog.html
[^a22]: Meta Engineering — Cache made consistent (Polaris). https://engineering.fb.com/2022/06/08/core-infra/cache-invalidation/
[^a23]: cache invalidation bug classes (serialized representation mismatch). https://engineering.fb.com/2022/06/08/core-infra/cache-invalidation/
[^a24]: AWS — caching best practices (don't cache fresh-required data). https://aws.amazon.com/caching/best-practices/
[^a25]: GreptimeDB / Consul stale-cache bug reports. https://github.com/
[^a27]: cache stampede / thundering herd on expiry. https://aws.amazon.com/caching/best-practices/
[^a28]: cache penetration (null-key DDoS). https://aws.amazon.com/caching/best-practices/
[^a29]: practitioner — slow query holds pooled connection (cl_waiting). https://www.percona.com/blog/proxysql-overhead-explained-and-measured/
[^a30]: caching masks unindexed query (canonical story). https://aws.amazon.com/caching/best-practices/
[^a31]: N+1 + Promise.all doesn't reduce pool pressure. https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing
[^a32]: raising max_connections as the worse band-aid. https://wiki.postgresql.org/wiki/Number_Of_Database_Connections
[^a33]: prescribed order: fix queries before caching. https://aws.amazon.com/caching/best-practices/
[^a35]: AWS re:Post — Lambda→RDS Proxy 8–12× slower case. https://repost.aws/
[^a37]: AWS — low-frequency Lambda direct-connection guidance. https://aws.amazon.com/blogs/compute/using-amazon-rds-proxy-with-aws-lambda/
[^a39]: multiplexing breaks pg_stat_activity 1:1 mapping. https://www.citusdata.com/blog/2020/10/08/analyzing-connection-scalability
[^a-openai]: OpenAI — Feb 2023 incident (PgBouncer pools clogged). https://status.openai.com/
