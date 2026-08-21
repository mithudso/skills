# The MongoDB Angle: `mongos` & Why MongoDB Rarely Uses Third-Party Query Proxies

> Provenance: reference under the standalone `database-proxies-query-middleware` skill. Created
> 2026-06-18 via `/dr` deep research. `verified-as-of: 2026-06-18` — the BI Connector EOL date is
> volatile; re-verify before relying on it.

## Contents
- [`mongos`: MongoDB's native query router](#mongos-mongodbs-native-query-router)
- [Connection handling, HA, placement](#connection-handling-ha-placement)
- [Why no third-party proxy: the driver already does it](#why-no-third-party-proxy-the-driver-already-does-it)
- [Adjacent query-access layers](#adjacent-query-access-layers)
- [Disconfirming nuance](#disconfirming-nuance)
- [Sources](#sources)

**The thesis:** MongoDB generally does **not** use third-party query proxies the way MySQL/Postgres
shops do, because the proxy responsibilities (server discovery, connection pooling, read/write
splitting, failover retries) are built into every official MongoDB driver, and `mongos` is
MongoDB's own native shard router. For the driver-spec **internals** (CMAP state machine, server
selection algorithm), defer to `mongodb-driver-internals` — this reference uses them only to
support the "why no third-party proxy" argument.

---

## `mongos`: MongoDB's native query router

- `mongos` is the query router and the **only** interface to a sharded cluster from the
  application's perspective; applications **never connect directly to shards** (even unsharded
  collections are reached via `mongos`).[^1][^4][^5]
- **Routing procedure:** (1) determine the shards that must receive the query; (2) establish a
  cursor on all targeted shards; (3) **merge** results into one response. Sort/limit modifiers are
  pushed to each shard first.[^1]
- **Targeted operation:** when the query includes the shard key (or a prefix of a compound shard
  key), `mongos` locates the chunk covering the value and routes only to the shard(s) holding it —
  the fastest pattern.[^1][^4][^5]
- **Scatter/gather (broadcast):** queries without the shard key are broadcast to **all** shards;
  `mongos` waits, merges, returns. These are "unfeasible for routine operations" at scale.
  Multi-updates and `updateMany()`/`deleteMany()` without the full shard key always broadcast.[^1][^4]
- **Config servers & metadata caching:** config servers (a replica set, CSRS) hold the
  **authoritative** routing table (chunk ranges per shard). `mongos` **caches** it in memory, holds
  no persistent state, and reloads **lazily** on a stale-version error (`flushRouterConfig` forces
  a reload).[^2][^3][^7]

## Connection handling, HA, placement

- Each `mongos` maintains its **own pool of connections to every member** of every replica set
  (including secondaries); requests use a connection one at a time (not multiplexed). Pools do not
  shrink when client count drops (restart to release).[^7][^1]
- **HA/scaling:** run **multiple `mongos`** for HA and scale. An LB in front must use **client
  affinity** (each client's connection reaches the same `mongos`).[^3]
- **Placement:** "the most common practice is to run `mongos` on the same systems as your
  application servers"; they can also live on the shards or dedicated hosts.[^1][^3]
- **Atlas:** `mongos` is managed for you; with a `mongodb+srv://` DNS seed list the app "automatically
  connects to the `mongos`," and topology changes don't require connection-string edits.[^28][^29]
- **DISCONFIRMING (a `mongos`-scaling caveat):** more routers is not strictly better — "`mongos`
  routers communicate frequently with your config servers. As you increase the number of routers,
  performance may degrade. If performance degrades, reduce the number of routers."[^3]

## Why no third-party proxy: the driver already does it

The functions a MySQL/Postgres shop bolts on via an external proxy are **standardized in the
driver**, via specs every official driver implements:

- **(a) SDAM** (Server Discovery and Monitoring) → server discovery, topology tracking, failover
  reaction. `serverMonitoringMode` defaults to `auto` (streaming normally; polling in FaaS).[^10][^26]
- **(b) CMAP** (Connection Monitoring and Pooling) → client-side connection pooling. Each
  `MongoClient` has a built-in per-server pool (`maxPoolSize` default **100**; `maxConnecting`
  default **2** throttles connection storms); the pool ties into SDAM (a server's pool is marked
  ready only after a successful check).[^8][^9][^16]
- **(c) Read preference + server selection** → read/write splitting. Read preference (`primary`/
  `primaryPreferred`/`secondary`/`secondaryPreferred`/`nearest`) routes reads to secondaries — the
  read-splitting an external proxy would otherwise provide. Server selection filters by read
  preference, then to servers within a latency window (`localThresholdMS` default **15 ms**), then
  picks randomly; for sharded topologies it load-balances across `mongos` within that window.[^11][^12][^13]
- **(d) Retryable reads/writes** → the driver retries transient failures **once**. Retryable writes
  on by default with server 4.2+ (`retryWrites=true`); retryable reads on by default with server
  6.0+ (`retryReads=true`).[^14][^15]

A MongoDB engineer states it directly: because drivers consistently provide connection monitoring
and pooling, "external connection pooling solutions aren't required (ex: Pgpool, PgBouncer)."[^17]
(Single named-engineer source — strongly indicative, not multi-source proof; the composite is
supported by the four specs above.)

**Synthesis:** routing (SDAM + server selection), pooling (CMAP), read/write splitting (read
preference), and failover retries (retryable reads/writes) live in the **driver/topology** layer;
`mongos` adds the shard-routing tier as MongoDB's **own native** component, not a third-party bolt-on.

## Adjacent query-access layers

These are **query surfaces, not proxies** — they expose a different query interface, not a
transparent pooling/routing layer:

- **BI Connector / `mongosqld`** — "acts as a layer that translates queries and data between a
  `mongod` or `mongos` instance and your reporting tool. The BI Connector **stores no data**." It
  gives BI tools (Tableau, Power BI) a **relational/SQL view** over MongoDB; `mongosqld` accepts SQL
  over the **MySQL wire protocol** (default port 3307). **DEPRECATION (as of 2026-06-18, per the
  docs banner):** the BI Connector reaches end-of-life and "will no longer be supported after
  September 2026"; MongoDB recommends the newer **Atlas SQL / MongoSQL** instead. (Volatile — verify
  the EOL date.)[^19][^20][^21]
- **Atlas Data Federation** — a distributed query engine that queries/transforms/moves data across
  Atlas clusters, Online Archive, and object stores (AWS S3 / Azure Blob / GCS) through **one
  unified service**, data staying in place. Query with MQL; an **Atlas SQL** endpoint serves BI
  tools. **MongoSQL is a (read-only) feature of Atlas Data Federation** and the BI Connector's
  recommended successor.[^22][^23][^21]

## Disconfirming nuance

Are proxy-like layers EVER used in front of MongoDB? Honestly:

- **`mongos` itself is arguably a proxy** — it terminates client connections and proxies/merges to
  backend `mongod`s; the distinction is that it is MongoDB's **native** component.[^1]
- **Protocol-compatible poolers exist for extreme serverless fan-out.** AWS published
  **`mongobetween`**, "a lightweight MongoDB connection pooler… handle a large number of incoming
  connections and multiplex them across a smaller connection pool," for Lambda/Fargate/ECS/EKS —
  but AWS positions it for **Amazon DocumentDB** (hard 30K-connection limit), not MongoDB Atlas.
  Trade-off flagged by AWS: sending more connections to the pooler than it holds downstream
  increases application latency.[^24][^25]
- **Serverless connection exhaustion is recognized, but MongoDB's answer is native:** cache the
  `MongoClient` at module scope with a small `maxPoolSize`, or use **Atlas Functions / App
  Services** (a MongoDB employee describes these as "optimized to not exceed cluster connection with
  smart internal connection pooling").[^27]
- **Load balancers in front of `mongos` are supported but constrained.** Drivers' per-server
  monitoring connections make a generic LB hard; MongoDB's sanctioned answer is the driver's
  **`loadBalanced=true`** mode (pins cursors/transactions to one connection behind the LB, reducing
  total connections), not a PgBouncer-style transparent pooler.[^18][^3]

**Net:** a PgBouncer/Pgpool-equivalent transparent third-party query proxy is **not standard** in
front of MongoDB. The genuine exceptions are (1) `mongos` itself, (2) protocol-compatible poolers
like `mongobetween` for hard-connection-cap systems under extreme fan-out (AWS positions this for
DocumentDB), and (3) MongoDB's own `loadBalanced` mode / Atlas App Services — all native or
narrowly scoped.

## Sources

[^1]: MongoDB Manual — Routing with mongos. https://www.mongodb.com/docs/manual/core/sharded-cluster-query-router/
[^2]: MongoDB Manual — Config Servers. https://www.mongodb.com/docs/manual/core/sharded-cluster-config-servers/
[^3]: MongoDB Manual — Sharded Cluster Components. https://www.mongodb.com/docs/manual/core/sharded-cluster-components/
[^4]: MongoDB Manual — Distributed Queries. https://www.mongodb.com/docs/v8.2/core/distributed-queries/
[^5]: MongoDB Manual — Sharding overview. https://www.mongodb.com/docs/manual/sharding/
[^7]: MongoDB Manual — FAQ: Sharding (mongos connection use, lazy cache reload). https://www.mongodb.com/docs/v7.0/faq/sharding/
[^8]: MongoDB Manual — Connection Pool Overview. https://www.mongodb.com/docs/manual/administration/connection-pool-overview/
[^9]: MongoDB — Connection Monitoring and Pooling (CMAP) spec. https://github.com/mongodb/specifications/blob/master/source/connection-monitoring-and-pooling/connection-monitoring-and-pooling.md
[^10]: MongoDB — Server Discovery and Monitoring (SDAM) spec. https://github.com/mongodb/specifications/blob/master/source/server-discovery-and-monitoring/server-discovery-and-monitoring.md
[^11]: MongoDB Manual — Read Preference. https://www.mongodb.com/docs/manual/core/read-preference/
[^12]: MongoDB Manual — Read Preference Mechanics (server selection). https://www.mongodb.com/docs/manual/core/read-preference-mechanics/
[^13]: MongoDB — Server Selection spec. https://github.com/mongodb/specifications/blob/master/source/server-selection/server-selection.md
[^14]: MongoDB Manual — Retryable Writes. https://www.mongodb.com/docs/manual/core/retryable-writes/
[^15]: MongoDB Manual — Retryable Reads. https://www.mongodb.com/docs/manual/core/retryable-reads/
[^16]: MongoDB Node.js driver — Connection Pools. https://www.mongodb.com/docs/drivers/node/current/connect/connection-options/connection-pools/
[^17]: Alex Bevilacqua (MongoDB) — MongoDB and Load Balancer Support. https://www.alexbevi.com/blog/2024/03/08/mongodb-and-load-balancer-support/
[^18]: MongoDB — Load Balancers spec (loadBalanced=true). https://github.com/mongodb/specifications/blob/master/source/load-balancers/load-balancers.md
[^19]: MongoDB — mongosqld (BI Connector). https://www.mongodb.com/docs/bi-connector/current/reference/mongosqld/
[^20]: MongoDB — What is the BI Connector? https://www.mongodb.com/docs/bi-connector/current/what-is-the-bi-connector/
[^21]: MongoDB — Transition from Atlas BI Connector to MongoSQL. https://www.mongodb.com/docs/sql-interface/transition-bic-to-atlas-sql/
[^22]: MongoDB — Atlas Data Federation. https://www.mongodb.com/docs/atlas/data-federation/
[^23]: MongoDB — Atlas Data Federation product page. https://www.mongodb.com/products/platform/atlas-data-federation
[^24]: AWS Database Blog — Scale connections with DocumentDB using mongobetween. https://aws.amazon.com/blogs/database/scale-your-connections-with-amazon-documentdb-using-mongobetween/
[^25]: aws-samples — mongobetween sample README. https://github.com/aws-samples/amazon-documentdb-samples/blob/master/samples/mongobetween_sample/README.md
[^26]: MongoDB Node.js driver — connection options (serverMonitoringMode FaaS). https://www.mongodb.com/docs/drivers/node/current/connect/connection-options/
[^27]: MongoDB Community forum — serverless connection limits (MongoDB employee reply). https://www.mongodb.com/community/forums/t/what-optimizations-for-not-hitting-connection-limit-in-serverless-functions/11507
[^28]: MongoDB Atlas — Modify a Cluster (DNS seed list auto-connect to mongos). https://www.mongodb.com/docs/atlas/scale-cluster/
[^29]: MongoDB Manual — Connection String Formats. https://www.mongodb.com/docs/manual/reference/connection-string-formats/
