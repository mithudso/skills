# Diagnostic playbook — symptom → skill/tool routing + evidence checklist

Loaded on demand by `solve-case` Phase 1 (assemble expertise) and Phase 3 (diagnosis).
The point is **selective loading**: pick the spoke(s) that match the symptom, not the whole tree.

## Symptom → primary skill → first evidence

| Presenting symptom | Primary skill(s) | First evidence to pull |
|---|---|---|
| Slow query / high latency | `mongodb-expert`, `atlas-diagnostics-expert` | `explain("executionStats")`, index list, Performance Advisor, slow-query log |
| High CPU / cache pressure | `atlas-diagnostics-expert` | FTDC, WiredTiger cache stats, metrics (cpu, cache dirty/used) |
| Connection errors (timeout, ECONNREFUSED, pool) | `mongodb-expert` (drivers/CMAP), `mongodb-kb` | driver logs, connection-string, IP access list, `serverStatus.connections` |
| Replication lag / elections / rollback | `mongodb-operations-expert` (replication), `mongodb-kb` | `rs.status()`, oplog window, `rs.printSecondaryReplicationInfo()` |
| Sharding: balancer, jumbo chunks, StaleConfig | `mongodb-expert` (sharding), `mongodb-kb` | `sh.status()`, chunk distribution, balancer state |
| Index build / missing index / wrong plan | `mongodb-expert` | `collection-indexes`, `explain`, `$indexStats` |
| Schema / data-model pain (unbounded arrays, 16MB) | `mongodb-expert` (schema) | `collection-schema`, doc-size sampling |
| Backup / restore / PITR / DR | `mongodb-operations-expert` | snapshot schedule, oplog retention, restore logs |
| Migration (mongosync, Relational Migrator, CDC) | `mongodb-operations-expert` | migration logs, lag metrics |
| Auth / TLS / network (SCRAM, x.509, LDAP, peering) | `mongodb-atlas-expert`, `mongodb-kb` | auth logs, TLS handshake error, network/peering config |
| Atlas platform (tiers, limits, Search, Stream) | `mongodb-atlas-expert` | Atlas Admin API state, alert config, limits |
| Encryption / CSFLE / Queryable Encryption | `mongodb-operations-expert` | key-vault config, KMS, schema map |
| Need a 10gen diagnostic tool (FTDC parse, etc.) | `10gen` | which repo/tool + install/run guidance |
| Error code lookup (11000, 64, 112, 286, …) | `mongodb-kb`, `mongodb-docset-lookup` | KB article, Manual page |

If the symptom spans **three or more** of these rows, prefer dispatching the
`uber-mongodb-diagnostician` agent rather than loading many spokes inline.

## Evidence checklist (Phase 3)

Aim for evidence, not assertion. Collect what is relevant; record what you could not get as a blocker.

- **Identity:** server version + FCV, Atlas tier, topology (replica set / sharded), driver + version.
- **Timeline:** when it started, what changed (deploy, config, traffic, version), what the customer tried.
- **Reproduction:** is it constant, intermittent, load-correlated, time-correlated?
- **Metrics (FTDC / Atlas):** CPU, memory, cache used/dirty, opcounters, connections, queues, page faults, IOPS, disk latency.
- **Query layer:** `explain("executionStats")` for the suspect op — keys vs docs examined, stage, plan, in-memory sort.
- **Indexes:** `collection-indexes`, `$indexStats` (used vs unused), build state.
- **Logs:** slow-query lines, errors, elections, assertions (`mongodb-logs`).
- **KB match:** does an existing Public KB article describe this exact symptom/error? Capture the URL.

## Deep-dive rule (Phase 4)

When you hit a term, default, error, or behavior you are not certain of:
1. `mongodb-docset-lookup` for exact Manual semantics/version behavior.
2. The matching `mongodb-*` spoke for distilled guidance.
3. `mongodb-kb` for known-issue / error-code context.
4. `context7` for driver/library API specifics.
5. Web (firecrawl/exa) only as a last resort, and mark any web-sourced claim as lower confidence.

Resolve the unknown before it enters the Phase-5 analysis. An unverified mechanism is a hypothesis,
not a fact, and must be labeled as such.
