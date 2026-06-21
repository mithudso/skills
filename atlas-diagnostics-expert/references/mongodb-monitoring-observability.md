<!-- hub-reference-banner -->
> **Reference file — part of the `atlas-diagnostics-expert` hub.** Formerly the standalone `mongodb-monitoring-observability` skill.
> Sibling MongoDB sub-topics are now reference files under the four hubs (`mongodb-expert`,
> `mongodb-atlas-expert`, `atlas-diagnostics-expert`, `mongodb-operations-expert`) — **not**
> standalone skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that
> name a bare `mongodb-*`/`atlas-*` skill; instead load that topic's `references/<name>.md`
> from the owning hub (see the hub's "Cross-hub map").

---

---
name: mongodb-monitoring-observability
title: MongoDB Monitoring and Observability
version: "1.3.0"
last-updated: "2026-05-29"
category: mongodb
description: >
  Comprehensive reference for monitoring MongoDB Atlas and self-managed deployments — Atlas
  built-in dashboards, Ops Manager/Cloud Manager monitoring agents, alert configuration and
  tuning, Datadog/New Relic/Prometheus integrations, FTDC (Full-Time Diagnostic Data Capture)
  analysis, slow query monitoring, replication lag diagnosis, connection metrics, and Atlas
  maintenance windows.
  TRIGGER: questions about Atlas metrics or dashboards, configuring alerts (CPU, replication lag,
  oplog window, connections), setting up Datadog or Prometheus integration for MongoDB,
  analyzing FTDC diagnostic data, slow query profiling (system.profile, Atlas Performance Advisor),
  replication lag root causes, connection pool exhaustion signals, Atlas maintenance windows
  (configuring, deferring, understanding rolling restart behavior), planned maintenance
  communication templates, or mongostat/mongotop/db.currentOp usage.
  SKIP: Atlas Search index tuning (use mongodb-atlas-search), Atlas cost optimization
  (use mongodb-cost-optimization), backup and restore planning (use mongodb-backup-restore),
  WiredTiger cache internals (use mongodb-wiredtiger), full query optimization beyond
  profiler analysis (use mongodb-query-performance).
tags:
  - mongodb
  - atlas
  - monitoring
  - observability
  - datadog
  - prometheus
  - new-relic
  - alerts
  - ftdc
  - slow-queries
  - replication
  - maintenance
  - maintenance-window
  - planned-operations
keywords:
  - Atlas monitoring
  - Atlas alerts
  - Datadog MongoDB integration
  - Prometheus MongoDB
  - New Relic MongoDB
  - FTDC diagnostics
  - slow query analysis
  - replication lag monitoring
  - connection pool monitoring
  - Atlas maintenance window
  - maintenance deferral
  - mongostat
  - mongotop
  - db.currentOp
  - Atlas Performance Advisor
  - system.profile
  - oplog window
  - WiredTiger metrics
  - Atlas Real-Time Performance Panel
  - Ops Manager monitoring
  - Cloud Manager monitoring
whenToUse:
  - Answering questions about Atlas metrics, dashboards, or alert configuration
  - Configuring third-party monitoring integrations (Datadog, Prometheus, New Relic)
  - Analyzing FTDC diagnostic files to diagnose historical performance issues
  - Slow query analysis via Atlas Profiler, system.profile, or mtools
  - Diagnosing replication lag root causes and flow control
  - Investigating connection pool exhaustion or connection limit issues
  - Configuring or deferring Atlas maintenance windows
  - Communicating planned Atlas maintenance to customers or stakeholders
  - Choosing the right monitoring tool for a given diagnostic question
  - Setting up Ops Manager or Cloud Manager for self-managed deployments
whenNotToUse:
  - Atlas Search index performance tuning — use mongodb-atlas-search
  - Atlas cost optimization or billing — use mongodb-cost-optimization
  - Backup and restore planning — use mongodb-backup-restore
  - WiredTiger cache internals and eviction tuning — use mongodb-wiredtiger
  - Deep query optimization beyond profiler output — use mongodb-query-performance
  - Ops Manager deployment, backup daemon, and automation (monitoring only covered here) — use mongodb-ops-manager
related_skills:
  - mongodb-wiredtiger
  - mongodb-performance-troubleshooting
  - mongodb-query-performance
  - mongodb-replication
  - mongodb-atlas-expert
  - mongodb-ops-manager
  - mongodb-backup-restore
  - mongodb-cost-optimization
audience: MongoDB TAMs, DBAs, and developers responsible for operating or advising on MongoDB Atlas and self-managed deployments
---

# MongoDB Monitoring and Observability

Comprehensive reference for monitoring MongoDB deployments — from Atlas built-in dashboards through third-party integrations, CLI tools, and low-level FTDC diagnostics.

**When to use this skill:** When answering questions about Atlas metrics, alert configuration, third-party monitoring integrations (Datadog, New Relic, Prometheus), FTDC diagnostics, slow query analysis, replication lag, connection pool behavior, or Atlas maintenance windows and planned operations.

**When not to use:** For Atlas Search index tuning (use `mongodb-search-ai`), Atlas cost optimization (use `mongodb-cost-optimization`), or backup/restore planning (use `mongodb-backup-restore`).

**Required roles for most monitoring operations:** `clusterMonitor` role on the `admin` database (self-managed), or Atlas `Project Read Only` / `Project Data Access Read Only` (Atlas UI). Third-party integrations (Datadog, Prometheus, New Relic) require Atlas `Project Owner` or `Organization Owner` to configure.

**Jump to:** [Quick Reference Tool Matrix](#quick-reference-tool-selection-matrix)

---

## 1. Atlas Cloud Monitoring — Built-in Metrics and Dashboard Customization

Atlas provides real-time and historical metrics for every cluster tier M10 and above. Free/shared-tier clusters have reduced metric granularity (5-minute resolution vs. 1-minute for dedicated tiers).

### Key metric categories available in Atlas UI

- **Opcounters** — insert, query, update, delete, getmore, command rates (ops/sec)
- **CPU / System** — process CPU, system CPU, I/O wait broken by read/write
- **Memory** — resident, virtual, mapped, cache (WiredTiger block cache, dirty bytes)
- **Disk I/O** — IOPS read/write, I/O utilization, disk queue depth
- **Network** — bytes in/out, number of requests
- **Connections** — current, available, total created
- **Replication** — oplog window hours, replication headroom, replication lag per secondary
- **Query targeting** — scanned/returned ratio (key indicator of missing indexes)
- **Tickets** — WiredTiger concurrent read/write tickets in use vs. available

### Dashboard customization

Atlas dashboards are pre-built per cluster but allow:
- Pin metric charts to a custom "Metrics" view for side-by-side comparison across nodes
- Toggle between individual node view (per-host) and cluster aggregate view
- Adjust time range (1h, 8h, 24h, 48h, 1w, custom)
- Use the **Real-Time Performance Panel** (RTPP) for 1-second granularity on live traffic — available on M10+ in the Atlas UI under the cluster's **Real Time** tab

The RTPP shows: opcounters, read/write tickets, connections, network, logical size, and an interactive `currentOp` view showing the slowest in-flight operations per namespace.

---

## 2. Ops Manager and Cloud Manager — Self-Managed Deployments

> **Deep reference:** see `mongodb-ops-manager` for full coverage of App DB sizing/HA, Backup Daemon placement, automation goal-state, air-gap/Local Mode, Kubernetes Operator, federation, and Live Migration to Atlas. This section covers the monitoring agent surface only.

**MongoDB Ops Manager** is the on-premises deployment of MongoDB's management platform for teams running MongoDB in their own data centers or private clouds. **MongoDB Cloud Manager** is the hosted SaaS version of the same platform — it provides identical monitoring, automation, and backup capabilities without requiring you to host the Ops Manager application yourself. Both share the same agent architecture described below.

### Core agents

| Agent | Role |
|---|---|
| **Automation Agent** | Deploys, configures, upgrades, and scales MongoDB processes via Ops Manager directives |
| **Monitoring Agent** | Collects real-time metrics from every managed `mongod`/`mongos`, ships to Ops Manager every 10 seconds |
| **Backup Agent** | Coordinates snapshot-based and oplog-based continuous backup |

### Monitoring agent behavior

- Runs as a daemon alongside your MongoDB processes
- Polls `serverStatus`, `replSetGetStatus`, `dbStats`, `collStats`, `currentOp` (filtered) at configurable intervals
- Stores time-series data in Ops Manager's own MongoDB backing store (separate from your application data)
- Sends alerts through Ops Manager's alert notification system — same alert types as Atlas

### Ops Manager / Cloud Manager dashboards

Both platforms replicate Atlas-style metric dashboards inside the web UI. The topology view shows replica set health, node states (PRIMARY/SECONDARY/ARBITER), and replication lag per member. The **Hardware** tab surfaces CPU, disk IOPS, and memory at host level for correlation with MongoDB behavior.

---

## 3. Atlas Alerts — Types, Channels, and Tuning

### Alert scope levels

- **Project-level alerts** — apply to all clusters in a project (e.g., CPU > 80% on any node)
- **Cluster-level alerts** — scoped to a specific cluster
- **Billing alerts** — monthly spend thresholds, data transfer thresholds

### Alert condition categories

| Category | Examples |
|---|---|
| Host / Node | CPU %, memory %, disk utilization %, disk IOPS |
| Replication | Replication lag > N seconds, oplog window < N hours |
| Connections | Connections > N (absolute or % of max) |
| Query performance | Slow queries, query targeting ratio |
| Indexes | Index build failures |
| Backup | Last successful snapshot age, restore failures |
| Atlas Search | Search index build failures |
| Billing | Monthly spend threshold, data transfer threshold exceeded |

### Notification channels

| Channel | Configuration |
|---|---|
| **Email** | One or more email addresses; configurable delay before sending |
| **Slack** | OAuth or webhook URL; route to specific channels |
| **PagerDuty** | PagerDuty integration key; supports routing rules/escalation policies |
| **Webhook** | HTTP POST to any endpoint; payload is JSON with alert details |
| **Datadog** | Forwards Atlas alert events as Datadog events alongside metrics |
| **OpsGenie** | OpsGenie API key |
| **VictorOps (Splunk On-Call)** | Routing key |
| **SMS / Phone** (via Twilio-backed Atlas feature) | Limited to some plan tiers |

### Alert tuning best practices

- **Set delay intervals** (e.g., "notify if condition persists for 5 minutes") to suppress transient spikes — CPU can spike briefly during flushes without being actionable
- **CPU alert baseline**: M10–M30 should alert at 75%; M50+ with sustained IOPS-heavy workloads often benefit from 85% thresholds with short delay
- **Replication lag**: alert at 10–15 seconds for most OLTP workloads; 60 seconds for batch-heavy pipelines
- **Oplog window**: never let it drop below 4 hours; alert at 48 hours to give time to investigate before backup windows are at risk
- **Connection count**: alert at 80% of the cluster's `maxIncomingConnections`; calculate max from `db.adminCommand({getCmdLineOpts:1})` or Atlas connection string parameters

---

## 4. Custom Metrics

### Atlas Custom Metrics

Atlas supports custom metric alerts via the **Atlas Administration API** (`/api/atlas/v2/groups/{groupId}/alertConfigs`). The `metricName` field accepts any metric Atlas exposes — including metrics not shown by default in the UI. Full metric name catalog: `https://www.mongodb.com/docs/atlas/reference/alert-conditions/`

### $currentOp polling for application-level insight

For application-level custom metrics, poll `$currentOp` on a schedule. Note: run this query from an admin-context connection — the `$all` field was deprecated in MongoDB 4.0 and removed in favor of the admin-context `currentOp` command directly:

```javascript
// Poll every 30 seconds via a dedicated monitoring connection (admin auth required)
// MongoDB 4.0+: run currentOp from admin context without $all
const ops = await db.admin().command({ currentOp: 1 });

// Extract slow ops (> 1 second), skip replication and internal namespaces
const slowOps = ops.inprog.filter(op =>
  op.secs_running > 1 &&
  op.ns &&
  !op.ns.startsWith('local.') &&
  !op.ns.startsWith('admin.')
);

// Emit to your metrics pipeline (Prometheus pushgateway, Datadog StatsD, etc.)
slowOps.forEach(op => {
  metrics.gauge('mongodb.slow_op.seconds', op.secs_running, {
    ns: op.ns,
    op: op.op,
    plan: op.planSummary
  });
});
```

Key fields in `$currentOp`:
- `secs_running` — wall time since op start
- `op` — command type (query, update, insert, command, getmore)
- `ns` — namespace (db.collection)
- `planSummary` — index used or COLLSCAN
- `waitingForLock` — boolean; stuck ops block others
- `msg` — progress message for long-running ops (index builds, aggregations)
- `locks` — lock modes held per resource type

### Application-level metrics to track

Beyond MongoDB server metrics, instrument your application code to capture:
- Query latency percentiles (p50, p95, p99) per collection
- Error rates by MongoDB error code
- Connection pool `waitQueueSize` — rising queue = pool exhaustion signal
- Retry attempt counts — spike in retries indicates transient replica set elections or network partitions

---

## 5. Datadog Integration

### Setup

Atlas Datadog integration requires M10+ clusters and a Datadog API key. Configure via:
- Atlas UI: **Project → Integrations → Datadog**
- Atlas API: `POST /api/atlas/v2/groups/{groupId}/integrations` with `type: "DATADOG"`
- Select region (`US1`, `US3`, `US5`, `EU1`, `AP1`, `US1_FED`) to match your Datadog account region

### Metrics shipped to Datadog

Atlas ships 100+ metrics prefixed `mongodb.atlas.*`. Key ones:

| Metric | Description |
|---|---|
| `mongodb.atlas.connections.current` | Active connections |
| `mongodb.atlas.opcounters.cmd` | Command rate |
| `mongodb.atlas.system.cpu.norm.guest` | Normalized CPU |
| `mongodb.atlas.mem.resident` | Resident memory bytes |
| `mongodb.atlas.cache.usage.dirty` | WiredTiger dirty cache bytes |
| `mongodb.atlas.repl.headroom` | Replication headroom (oplog - lag) |
| `mongodb.atlas.disk.partition.iops.total` | Total disk IOPS |
| `mongodb.atlas.query.targeting.scannedObjectsPerReturned` | Scan ratio |

### Datadog Database Monitoring for Atlas

Beyond the metrics integration, Datadog offers **Database Monitoring (DBM) for MongoDB Atlas** (separate feature, requires Datadog Agent with MongoDB integration configured). DBM provides:
- Query-level explain plan capture
- Wait event analysis
- Query normalization and fingerprinting
- Query sample storage for historical analysis

Configure via `conf.d/mongo.d/conf.yaml` in the Datadog Agent. Use Atlas connection string with a monitoring user that has `clusterMonitor` role.

### Datadog dashboards

Datadog ships a pre-built **MongoDB Atlas Overview** dashboard in the integrations library. Clone it to customize. Recommended widgets to add:
- Query targeting ratio over time (composite graph: scanned / returned)
- WiredTiger cache fill % (dirty bytes / total cache)
- Election events as event overlay

---

## 6. New Relic Integration

### Atlas New Relic integration

Configure via Atlas UI: **Project → Integrations → New Relic**. Requires a New Relic license key and account ID. Metrics are shipped as custom events/metrics to your New Relic account under the `MongoDBAtlas.*` namespace.

### APM correlation

New Relic's primary value in MongoDB monitoring is **APM-to-database correlation**: when your application uses the New Relic APM agent, New Relic links slow transaction traces in application code directly to slow MongoDB operations. This correlation shows:
- Which application endpoint triggered a slow query
- Time breakdown: application processing vs. MongoDB query time
- Database query fingerprints alongside application traces

### New Relic on-host integration (self-managed)

For self-managed MongoDB, use the New Relic Infrastructure agent with the MongoDB on-host integration (`nri-mongodb`):

```yaml
# /etc/newrelic-infra/integrations.d/mongodb-config.yml
integrations:
  - name: nri-mongodb
    env:
      USERNAME: newrelic_monitor
      PASSWORD: ${MONGODB_PASSWORD}
      HOST: localhost
      PORT: 27017
      AUTH_SOURCE: admin
      SSL: true
      REPLICA_SET: rs0
      EXTENDED_METRICS: true
      COLLECTION_METRICS: true
```

Metrics collected: throughput, latency, connections, memory, replication lag, collection-level stats when `COLLECTION_METRICS: true`.

---

## 7. Prometheus Integration

### Atlas managed Prometheus endpoint

Atlas exposes a Prometheus-compatible scrape endpoint (M10+ only). Enable via:
- Atlas UI: **Project → Integrations → Prometheus**
- Generates a scrape URL in the form:
  `https://cloud.mongodb.com/prometheus/v1.0/groups/{groupId}/metrics`
- Authentication: HTTP Basic auth with Atlas programmatic API public/private key pair

### Scrape configuration

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'mongodb-atlas'
    scrape_interval: 60s
    scrape_timeout: 55s
    scheme: https
    basic_auth:
      username: '<atlas_public_api_key>'
      password: '<atlas_private_api_key>'
    static_configs:
      - targets:
          - 'cloud.mongodb.com'
    metrics_path: '/prometheus/v1.0/groups/<groupId>/metrics'
    # Optional: filter to specific cluster
    params:
      cluster: ['<clusterName>']
```

### Key Prometheus metrics from Atlas

| Metric name | Type | Description |
|---|---|---|
| `hardware_system_cpu_norm_guest_percent` | gauge | Normalized CPU % |
| `assert_regular` | counter | Regular assertions |
| `connections_current` | gauge | Open connections |
| `cache_bytes_currently_in_cache` | gauge | WiredTiger cache bytes |
| `cache_dirty_bytes_in_the_cache` | gauge | Dirty cache bytes |
| `mem_resident` | gauge | Resident memory MB |
| `opcounters_cmd_per_second` | gauge | Commands per second |
| `replication_lag` | gauge | Replication lag seconds |
| `query_targeting_scanned_per_returned` | gauge | Scan ratio |

### Federation and self-managed Prometheus

For self-managed MongoDB, use `mongodb_exporter` (Percona's open-source exporter) or the official MongoDB community exporter. Connect to Prometheus via standard scrape config on port 9216 (default). Metrics include everything from `serverStatus`, `replSetGetStatus`, and `dbStats`.

```yaml
scrape_configs:
  - job_name: 'mongodb-selfmanaged'
    static_configs:
      - targets: ['mongo-host-1:9216', 'mongo-host-2:9216', 'mongo-host-3:9216']
    relabel_configs:
      - source_labels: [__address__]
        target_label: instance
```

Search for "MongoDB Overview Percona" in the Grafana dashboard library for a production-ready self-managed Prometheus + Grafana starting point.

---

## 8. FTDC (Full Time Diagnostic Capture)

### What is FTDC

FTDC (Full Time Diagnostic Data Capture) is MongoDB's always-on internal diagnostic system, enabled by default in `mongod` and `mongos` since MongoDB 3.2. It is the first artifact MongoDB Support requests for any performance investigation.

### What FTDC captures

Every **second**, FTDC samples:
- Full `serverStatus` output (opcounters, memory, connections, WiredTiger stats, lock stats, network)
- `replSetGetStatus` (for replica set members)
- `local.oplog.rs` metadata (oplog window, latest/oldest oplog timestamps)
- System CPU and memory counters (via OS-level sampling)
- Storage engine internal stats (WiredTiger block cache, page eviction, checkpoint timing)

Every **200ms** (high-frequency samples), FTDC captures a lighter subset focused on CPU and I/O counters. This allows reconstruction of spikes shorter than 1 second.

### Storage and location

FTDC files live at `<dbPath>/diagnostic.data/`:
```
diagnostic.data/
  metrics.2024-01-15T00-00-00Z-00000   # compressed binary
  metrics.interim                        # current unsaved buffer (partial — last ~seconds may be missing)
  metrics.2024-01-16T00-00-00Z-00000
```

Files rotate when they reach ~10 MB (configurable via `diagnosticDataCollectionFileSizeMB`). Atlas retains FTDC data and surfaces it to MongoDB Support automatically — you do not need to retrieve files manually for Atlas clusters.

### Retrieval for self-managed deployments

For on-prem MongoDB, copy the entire `diagnostic.data/` directory from the server. Copying while `mongod` is live is safe — FTDC uses its own internal write path. Be aware that `metrics.interim` is a partial buffer containing data not yet flushed to a numbered file; the most recent few seconds of telemetry may be absent from it.

### Analysis tools

| Tool | Usage |
|---|---|
| **Keyhole** | `keyhole --ftdc diagnostic.data/` — produces human-readable reports and Grafana-compatible output |
| **mongodb/ftdc Go library** | Low-level BSON parsing of FTDC metrics chunks |
| **mtools `mloginfo`** | Parses mongod logs alongside FTDC for correlation |
| **Atlas Support** | Atlas Support has direct access to FTDC for all Atlas clusters |

### FTDC in troubleshooting

Use FTDC to answer:
- Was there a checkpoint stall? (WiredTiger checkpoint duration spike)
- Was CPU saturated? (system CPU counters at 100%)
- Did connection count spike before the incident? (connections.current time series)
- Was replication lag gradual or sudden? (replSetGetStatus.members[].optimeDate delta)
- Did the WiredTiger cache fill and start aggressive eviction? (cache dirty % over time)

---

## 9. mongotop / mongostat / db.currentOp

### mongostat

Streams a running summary of MongoDB server activity, similar to Unix `vmstat`. Output updates every second by default.

```bash
mongostat --uri "mongodb+srv://user:pass@cluster.mongodb.net" \
          --discover \
          --rowcount 60
```

Key columns:
| Column | Meaning |
|---|---|
| `insert/query/update/delete` | Operations per second |
| `getmore` | Cursor getmore ops/sec |
| `command` | Command ops/sec (includes replicated writes counted twice) |
| `dirty` | WiredTiger dirty cache % |
| `used` | WiredTiger cache in-use % |
| `flushes` | WiredTiger checkpoints in the interval |
| `vsize/res` | Virtual/resident memory |
| `qrw/arw` | Queue length and active count for reads/writes |
| `conn` | Current connections |
| `repl` | Replication state (PRI/SEC/REC) |

**When to use**: quick snapshot of overall server load; watching for sudden spikes in op rates or queue buildup; checking cache utilization trend in real time.

### mongotop

Shows per-collection time spent on reads and writes, updated every second.

```bash
mongotop --uri "mongodb+srv://user:pass@cluster.mongodb.net" 5
```

The `5` argument sets the polling interval in seconds. Output shows `ns` (namespace), `total` (ms per interval), `read` (ms), `write` (ms).

**When to use**: identify which collection is hottest during a performance issue; confirm that a problematic query is hitting the expected namespace; narrow down which collection to profile next.

### db.currentOp()

`db.currentOp()` is a mongosh helper alias for `db.adminCommand({currentOp: true, ...})`. Returns in-flight operations at the moment of invocation. More detailed than `mongostat`/`mongotop`.

```javascript
// All active operations longer than 2 seconds, excluding replication namespaces
db.adminCommand({
  currentOp: true,
  active: true,
  secs_running: { $gt: 2 },
  ns: { $not: /^local\./ }
})

// Kill a specific operation
db.adminCommand({ killOp: 1, op: <opid> })
```

**When to use**: investigate a specific slow operation in real time; find lock waiters (`waitingForLock: true`); identify which session holds a lock before `killOp`.

---

## 10. Slow Query Monitoring

### Atlas Profiler

Atlas exposes the MongoDB query profiler via **Cluster → Performance Advisor** and **Cluster → Profiler** tabs. The Profiler tab streams slow queries in near-real-time (subject to Atlas data pipeline latency of ~2 minutes).

- Set slow query threshold via Atlas UI (Cluster → **Configuration** → "Slow Query Threshold") or via `db.setProfilingLevel()`
- Default Atlas slow threshold: **100ms** (configurable down to 0ms for profiling all queries)
- Atlas Performance Advisor automatically surfaces index recommendations based on slow query patterns, ranking by average execution time × frequency

### system.profile collection

When profiling level is 1 or 2, MongoDB writes operation metadata to `db.system.profile` (a capped collection, 1 MB by default):

```javascript
// Enable profiling for ops > 100ms
db.setProfilingLevel(1, { slowms: 100 })

// Query the profile
db.system.profile.find({
  millis: { $gt: 500 },
  ns: { $ne: 'admin.system.profile' }
}).sort({ ts: -1 }).limit(20).pretty()
```

Key fields: `millis` (execution time), `ns` (namespace), `op` (operation type), `planSummary` (IXSCAN vs. COLLSCAN), `keysExamined`, `docsExamined`, `nreturned`, `locks`, `queryHash`, `planCacheKey`.

**Important**: profiling level 2 (log all operations) has measurable overhead on busy clusters — use level 1 with a tuned `slowms` threshold in production.

### Slow query log parsing with mtools

```bash
pip install mtools
mloginfo mongod.log --queries
# Groups by query shape, shows count / avg ms / max ms / 95th percentile
```

### Atlas slow query threshold tuning guidance

| Workload | Recommended threshold |
|---|---|
| OLTP (< 10ms target) | 20–50ms |
| Mixed OLTP/analytics | 100ms (default) |
| Analytics-heavy | 200–500ms |
| Bulk load / maintenance windows | Raise to 1000ms to avoid noise |

---

## 11. Replication Lag Monitoring

### Measuring replication lag

```javascript
// From mongosh on any replica set member
rs.printSecondaryReplicationInfo()
// Output: source, syncedTo, lag in seconds, heartbeat lag

// Programmatic approach — use optimeDate (a JS Date), not optime.ts (a BSON Timestamp)
const status = db.adminCommand({ replSetGetStatus: 1 })
const primary = status.members.find(m => m.stateStr === 'PRIMARY')
status.members.filter(m => m.stateStr === 'SECONDARY').forEach(sec => {
  const lagMs = primary.optimeDate.getTime() - sec.optimeDate.getTime()
  console.log(`${sec.name}: lag ${lagMs}ms`)
})
```

Note: use `optimeDate` (a JavaScript Date object) rather than `optime.ts` (a BSON Timestamp with `.t` and `.i` integer fields) — calling `.getTime()` on a BSON Timestamp will return `undefined`.

### Root causes of replication lag

1. **Secondary under-resourced** — secondary CPU/disk can't keep up with replication workload; typically fixed by upgrading node tier or distributing reads
2. **Replication flow control** (MongoDB 4.2+) — primary throttles write throughput when secondaries fall behind; visible in `replSetGetStatus.flowControl`
3. **Chained replication misconfiguration** — secondary syncing from another secondary rather than primary; check `rs.status().syncSourceHost`
4. **Long-running transactions on secondary** — secondary blocks replication to apply multi-document transactions atomically; check `maxTransactionLockRequestTimeoutMillis`
5. **Network partition or bandwidth saturation** — check `replication.network.bytes` in `serverStatus`

### Flow control diagnostics

```javascript
db.adminCommand({ serverStatus: 1 }).flowControl
// Fields: enabled, targetRateLimit, timeAcquiringMicros, isLagged
```

When `isLagged: true`, the primary is actively throttling writes. Address the root cause on the lagging secondary rather than disabling flow control (`enableFlowControl: false`), which risks oplog overflow.

### Lag alert thresholds

| Deployment type | Warning | Critical |
|---|---|---|
| OLTP, strict secondary reads | 5s | 15s |
| General purpose | 15s | 60s |
| Analytics/reporting secondaries | 60s | 300s |

Atlas alert: **Replica Set Replication Lag** condition in **Project Alerts → Add New Alert → Replication**.

---

## 12. Connection Metrics

### Key connection counters in serverStatus

```javascript
const ss = db.adminCommand({ serverStatus: 1 })

// Current connection state
ss.connections.current       // Active connections right now
ss.connections.available     // Remaining capacity (maxIncomingConnections - current)
ss.connections.totalCreated  // Cumulative connections since server start (monotonic)

// WiredTiger connection/session info
ss.wiredTiger.concurrentTransactions.read.out   // Active read tickets
ss.wiredTiger.concurrentTransactions.read.available
ss.wiredTiger.concurrentTransactions.write.out  // Active write tickets
ss.wiredTiger.concurrentTransactions.write.available
```

### Connection pool exhaustion signals

| Signal | What to look for |
|---|---|
| `connections.available` approaching 0 | Imminent connection refusal |
| `totalCreated` rate high | Connections closing and reopening rapidly (pool churn) |
| Driver-side `waitQueueSize` rising | Application waiting for pool slot |
| `ServerSelectionTimeoutError` in app logs | Pool exhausted before timeout |
| `Too many open files` in mongod log | OS file descriptor limit hit (ulimit -n) |

### Atlas connection limits by cluster tier

| Tier | Max connections |
|---|---|
| M10 | 1,500 |
| M20 | 3,000 |
| M30 | 3,000 |
| M40 | 6,000 |
| M50 | 16,000 |
| M60 | 32,000 |
| M80 | 64,000 |
| M200+ | 128,000 |

Connections are per-node, not per-cluster. A 3-node replica set at M30 has 3 × 3,000 = 9,000 total connections across all nodes, but the driver only connects to the primary (and any secondary for tagged read preference operations).

### Connection tuning recommendations

- **Driver `maxPoolSize`**: default is 100 per MongoClient instance. For serverless/Lambda deployments, reduce to 5–10 to avoid connection storms during cold-start bursts.
- **Use a single MongoClient per process** — the most common connection leak pattern is creating a new `MongoClient` per request.
- **Enable `waitQueueTimeoutMS`** (Node.js) or `waitQueueTimeout` (Python) to surface pool exhaustion quickly rather than hanging indefinitely.
- **`minPoolSize`**: set to a value that keeps connections warm for predictable traffic, but be aware each minimum connection consumes a server-side connection slot permanently.
- **Atlas connection string option `maxIdleTimeMS`**: controls how long idle connections are kept; reduce (e.g., 60000ms) for Lambda/serverless to match function lifecycle.

---

## Quick Reference: Tool Selection Matrix

| Question | Tool |
|---|---|
| What is the server doing right now? | `mongostat` + Atlas RTPP |
| Which collection is hottest? | `mongotop` |
| What specific operation is slow right now? | `db.currentOp()` / `db.adminCommand({currentOp:1})` |
| What slow queries ran in the past hour? | Atlas Profiler / `system.profile` |
| Why was the server slow at 2am? | FTDC + Keyhole |
| Is replication healthy? | `rs.printSecondaryReplicationInfo()` |
| Are connections running out? | `serverStatus.connections` + Atlas alerts |
| Correlate MongoDB to app performance? | Datadog DBM or New Relic APM |
| Long-term trending (weeks/months)? | Prometheus + Grafana or Datadog dashboards |
| Billing and cluster-level spend? | Atlas billing alerts |
| Atlas Search index health? | Atlas UI → Search tab → Index Metrics |
| Self-managed cluster automation + monitoring? | Ops Manager or Cloud Manager |
| When is maintenance scheduled / what window is configured? | `atlas maintenanceWindows describe` / Atlas UI Project Settings → [§13](#13-atlas-maintenance-windows-and-planned-operations) |
| How do I defer upcoming maintenance? | Atlas UI Defer button or `atlas maintenanceWindows defer` → [§13](#13-atlas-maintenance-windows-and-planned-operations) |

---

## 13. Atlas Maintenance Windows and Planned Operations

### Free and shared tier clusters (M0, M2, M5)

**M0, M2, and M5 clusters do not support configurable maintenance windows.** Atlas manages all maintenance for free and shared-tier clusters entirely, with no operator control over timing. These clusters may be restarted at any time. If maintenance window control is required, upgrade to M10 or higher.

This is a common point of confusion when customers first configure maintenance windows at the project level — the setting applies only to dedicated-tier clusters (M10+) in that project.

### Maintenance window configuration

Atlas maintenance windows are configured at the **project level** and apply to all dedicated-tier (M10+) clusters within that project. A single project-level window governs when Atlas schedules rolling restarts for patch upgrades and infrastructure updates.

**Location:** Atlas UI → **Project Settings** → **Maintenance Window**

**Default behavior:** When no custom window is configured, Atlas selects the window (commonly Tuesday 10:00–12:00 UTC for many regions, though this varies). Relying on the default is not recommended for production workloads — configure an explicit window aligned with your lowest-traffic period.

**Configuring a custom window:**
- Choose day of week (Sunday through Saturday; Sunday=1 in the API/CLI, matching the integer table below)
- Choose start hour in UTC (0–23); the window is exactly 1 hour
- Minimum window size: 1 hour (Atlas cannot currently support sub-hourly windows)
- Changes take effect immediately and persist until cleared

**Important scope limitation:** Maintenance windows are project-scoped, not per-cluster. There is no supported mechanism to set a different maintenance window per cluster within the same project. If clusters have different maintenance requirements (e.g., a dev cluster vs. a prod cluster), place them in separate Atlas projects.

**Atlas CLI commands:**
```bash
# View the current maintenance window for a project
atlas maintenanceWindows describe --projectId <projectId>

# Set maintenance window to Sundays at 02:00 UTC
atlas maintenanceWindows update \
  --dayOfWeek 1 \
  --hourOfDay 2 \
  --projectId <projectId>

# Clear a configured maintenance window (reverts to Atlas-chosen default)
atlas maintenanceWindows clear --projectId <projectId>
```

Day-of-week values: Sunday=1, Monday=2, Tuesday=3, Wednesday=4, Thursday=5, Friday=6, Saturday=7.

### What triggers maintenance

Not all cluster changes go through the configured maintenance window:

| Trigger | Follows Maintenance Window? |
|---|---|
| MongoDB patch version upgrade (e.g., 7.0.8 → 7.0.9) | Yes |
| Atlas infrastructure / hardware updates | Yes |
| Feature releases requiring restart | Yes |
| Critical security patch (CVE) | No — Atlas may override window |
| Major version upgrade (e.g., 6.0 → 7.0) | No — separately scheduled by operator |
| Cluster tier scaling (scale up/down) | No — operator-initiated, immediate rolling restart |
| Storage scaling | No — operator-initiated |
| Cluster pause / resume | No — operator-initiated |

Emergency security patches bypass the maintenance window entirely. Atlas notifies project and organization owners via email before performing emergency maintenance, but the window configuration does not constrain it.

### How Atlas performs rolling maintenance

Atlas uses a rolling restart procedure to minimize downtime:

1. **Secondaries first** — Atlas restarts one secondary at a time, waiting for each to rejoin the replica set and catch up on replication before proceeding to the next.
2. **Primary last** — After all secondaries are restarted, the primary is restarted, which triggers a replica set election.
3. **Election window** — Primary elections typically complete in 10–30 seconds. During this window, write operations are temporarily unavailable and read operations fall back to secondaries (if read preference allows).
4. **mongos nodes (sharded clusters only)** — After all shards complete their rolling restarts, the mongos routers are restarted last. Skip this step for replica-set-only deployments.

**Application impact:** Applications experience a brief connection interruption during the primary election step. Drivers with retryable writes handle this transparently — the driver detects the "not primary" error and retries on the newly elected primary. Applications that do not use retryable writes may see transient write failures.

**Atlas alert during maintenance:** The **"Primary election"** alert fires during this procedure. If you have this alert configured, expect it to fire during every scheduled maintenance window. Consider suppressing or acknowledging this alert type during known maintenance periods, or configuring a lower-urgency notification channel for it.

**Total maintenance duration:**
- Simple 3-node replica set: typically 5–15 minutes end-to-end
- Sharded cluster with multiple shards: multiply by number of shards; total duration can be 30–60 minutes for large topologies
- Clusters with secondaries that are significantly behind in replication before maintenance starts will take longer

### Deferring maintenance

When Atlas has scheduled upcoming maintenance, a **Defer** button appears in the Atlas UI next to the maintenance notification.

**Rules for deferral:**
- Deferral postpones maintenance by exactly **7 days**
- Each scheduled maintenance event may be deferred **once only**
- After one deferral, the maintenance cannot be deferred again and will execute at the rescheduled time
- **Critical security patches cannot be deferred** regardless of the deferral policy — attempting to defer will return an error

**Atlas UI:** The Defer button appears on the cluster detail page and in **Project Settings → Maintenance Window** when maintenance is actively scheduled.

**Atlas CLI:**
```bash
atlas maintenanceWindows defer --projectId <projectId>
```

**Atlas API:**
```bash
POST /api/atlas/v2/groups/{groupId}/maintenanceWindow/defer
# No request body required; authorization via API key
```

### Querying the maintenance window via API

To retrieve the current maintenance window configuration programmatically (useful for automation, dashboards, or verifying settings):

```bash
# GET current maintenance window — returns dayOfWeek, hourOfDay, startASAP, autoDeferOnceEnabled
curl -u "{publicKey}:{privateKey}" --digest \
  "https://cloud.mongodb.com/api/atlas/v2/groups/{groupId}/maintenanceWindow" \
  -H "Accept: application/vnd.atlas.2023-01-01+json"
```

Response fields:
- `dayOfWeek` — integer 1–7 (Sunday=1 through Saturday=7); absent if no custom window set
- `hourOfDay` — integer 0–23 UTC start hour
- `startASAP` — boolean; true if Atlas has queued maintenance to run at next opportunity
- `autoDeferOnceEnabled` — boolean; if true, Atlas automatically defers the next maintenance once

### Emergency and critical security patches

When Atlas identifies a critical CVE or security vulnerability requiring an immediate patch:
- Atlas may perform out-of-window maintenance with **same-day or next-day notice** for critical CVEs
- For lower-severity security updates, Atlas typically provides **24–48 hours notice** via email
- Email notifications are sent to **all Project Owners and Organization Owners** for the affected project
- Emergency patches **cannot be deferred**
- Atlas will still use rolling restart procedure to minimize impact

Monitor your Atlas project **Activity Feed** (Atlas UI → Project → Activity) for maintenance events. Each maintenance start and completion is logged here with timestamps.

### Minimizing application impact

The most effective way to make maintenance transparent to users is to ensure your driver and connection configuration supports retryable operations:

**Recommended connection string options:**

| Option | Recommended value | Reason |
|---|---|---|
| `retryWrites` | `true` (default since driver 4.2) | Retry writes transparently on primary election |
| `retryReads` | `true` | Retry reads on network errors and primary changes |
| `serverSelectionTimeoutMS` | `30000` | Elections take 10–30s; 3000ms causes premature `ServerSelectionTimeoutError` |
| `connectTimeoutMS` | `10000` | Standard connection establishment timeout |
| `socketTimeoutMS` | `0` (disabled) | Short socket timeouts interfere with long-running operations |

**Why `serverSelectionTimeoutMS` matters:** During a primary election, all nodes briefly report themselves as non-primary. If `serverSelectionTimeoutMS` is 3000ms and the election takes 12 seconds, the driver gives up before the new primary is elected. Setting 30000ms gives the driver enough headroom to wait out the election.

**Connection pool warm-up:** After a primary election, the first queries to the new primary may be slower as connections are re-established and the WiredTiger cache warms for the new primary's working set. This is normal and resolves within 30–60 seconds.

**Alert correlation:** Configure your monitoring to correlate "Primary election" alerts with maintenance window times. Alerts firing within the configured window are expected; alerts outside the window warrant investigation.

### Atlas maintenance for sharded clusters

Sharded cluster maintenance follows this expanded sequence:

1. **Config server replica set** — the CSRS is restarted first (rolling restart of config server members)
2. **Shard replica sets** — each shard is restarted in turn (rolling restart within each shard; shards are processed sequentially, not in parallel)
3. **mongos routers** — all mongos instances are restarted last (can be done in parallel since mongos is stateless)

**Balancer behavior:** The chunk balancer is suspended during maintenance. Migrations that were in progress before maintenance began are allowed to complete; no new migrations are initiated until maintenance finishes. This avoids partial migrations stranded across a restart boundary.

**Total time estimation for sharded clusters:** Multiply the per-shard restart time by the number of shards, then add config server and mongos restart time. A 4-shard cluster where each shard takes 10 minutes will take approximately 40–50 minutes for the shard restarts alone.

### Customer communication template

When advising customers on how to communicate planned Atlas maintenance to their own end-users or stakeholders, use the template below. Replace all `[bracketed]` placeholders before sending.

```
Subject: Planned database maintenance — [Day, Month DD YYYY]

We have scheduled routine database maintenance to apply a MongoDB patch update.

Maintenance window: [Day of week, YYYY-MM-DD] [HH:MM]–[HH:MM] UTC
Expected impact:    Brief connection interruption lasting less than 30 seconds
                    during the database primary election. No data loss will occur.
Action required:    None. Applications using current MongoDB drivers with retryable
                    writes enabled will handle the interruption automatically.
                    Applications with custom retry logic may see a single retried
                    operation.

If you observe issues after [HH:MM] UTC, please contact [support channel / Slack #channel].
```

**Placeholder guide:**
- `[Day, Month DD YYYY]` — e.g., "Tuesday, June 10 2025"
- `[HH:MM]–[HH:MM] UTC` — e.g., "02:00–03:00 UTC" (the configured 1-hour window)
- `[support channel / Slack #channel]` — your team's escalation path

Key points every maintenance communication must include:
- Specific UTC time window — never use vague language like "tonight" or "overnight"
- Impact duration scoped to the election (< 30 seconds), not the full rolling restart
- Explicit "no data loss" statement — customers frequently conflate restarts with data risk
- Retryable writes note — acknowledges that non-retryable operations may see one transient error
- Clear escalation path if issues persist after the window closes
