# Scenario Catalog

Use this file to pick the closest dominant symptom fast. Choose early, but switch if evidence contradicts the initial fit.

| Scenario | Dominant symptom | First evidence surfaces | Primary skill route |
| --- | --- | --- | --- |
| Inefficient Query | One or more expensive query shapes drive latency or resource contention | Performance Advisor, profiler, explain, query targeting | `atlas-diagnostics-expert` -> `mongodb-expert` |
| Operation Spike | Sudden surge in ops counters, docs examined, or data read | Metrics, alerts, recent deploy/change window | `atlas-diagnostics-expert` |
| Multi-Plan Storm | Replanning churn causes unstable latency | Profiler, explain, plan behavior, logs | `atlas-diagnostics-expert` -> `mongodb-expert` |
| WT Cache Pressure | Cache fill or eviction pressure is the main symptom | Metrics, FTDC, WiredTiger cache counters | `atlas-diagnostics-expert` -> `mongodb-expert` |
| Dirty Cache Full | Dirty cache growth blocks progress | Metrics, FTDC, storage signals | `atlas-diagnostics-expert` -> `mongodb-expert` |
| Update Cache Full | Update-heavy cache pressure dominates | Metrics, FTDC, write workload profile | `atlas-diagnostics-expert` -> `mongodb-expert` |
| Memory Fragmentation | Memory pressure is high without proportional working-set explanation | Metrics, FTDC, host/runtime context | `atlas-diagnostics-expert` -> `mongodb-expert` |
| Unbounded Cluster Memory Usage | Cluster memory use keeps climbing toward OOM or instability | Metrics, FTDC, workload shifts, connection counts | `atlas-diagnostics-expert` |
| Operation Priority Inversion — Index Build | Index work blocks or distorts primary workload | Activity feed, profiler, metrics, index build timing | `atlas-diagnostics-expert` -> `mongodb-expert` |
| Operation Priority Inversion — Sharding Ops | Chunk movement or range deletion competes with live traffic | Sharding status, metrics, logs, activity feed | `atlas-diagnostics-expert` -> `mongodb-expert` |
| Excessive dhandles | Open-handle or checkpoint pressure dominates | FTDC, storage metrics, logs | `atlas-diagnostics-expert` -> `mongodb-expert` |
| Connection Spike | Abrupt connection growth causes pressure or limit hits | Connection metrics, app rollout timing, pool behavior | `atlas-diagnostics-expert` |
| Write Conflict Storm | Write conflicts increase sharply and throughput degrades | Metrics, logs, workload pattern, hot-document behavior | `atlas-diagnostics-expert` -> `mongodb-expert` |
| Fell Off Oplog | Members are recovering or cannot catch up | Replica set state, oplog window, logs, FTDC | `atlas-diagnostics-expert` -> `mongodb-expert` |
| Disk Full | Storage exhaustion or hard floor is imminent | Disk metrics, alerts, activity feed | `atlas-diagnostics-expert` -> `mongodb-operations-expert` |
| Cloud Provider Incident — Capacity | Infra scarcity or capacity error blocks Atlas recovery actions | Atlas activity, provider signals, scaling/failover attempts | `mongodb-atlas-expert` -> `atlas-diagnostics-expert` |
| Cloud Provider Incident — Region / AZ Outage | Regional failure or zone disruption drives impact | Provider status, Atlas topology, failover state | `mongodb-atlas-expert` -> `mongodb-operations-expert` |
| Atlas Platform Incident | Atlas control-plane or platform behavior is the dominant blocker | Activity feed, Admin API symptoms, support-wide evidence | `mongodb-atlas-expert` -> `atlas-diagnostics-expert` |
| Networking — IP Access List / Peering / Private Endpoint | Reachability or routing is the dominant symptom | Network config, DNS, Atlas access settings, endpoint status | `mongodb-atlas-expert` |
| Security — KMS Permission Issue | Encryption or key-access failure blocks service | KMS errors, audit trail, recent IAM/key changes | `mongodb-operations-expert` -> `mongodb-atlas-expert` |

## Tooling defaults

- Start with `ts-diag` when the fastest path is Atlas project or cluster triage.
- Move to `alexandria` for FTDC-heavy diagnosis.
- Use SearchPlanIQ for Atlas Search or Vector Search explain-plan cases.
- Use `mongolyser` when a broader interactive diagnostics view is helpful.
- Use `mongodb-kb` when you need a known-issue check or a customer-shareable support article.

## Selection rules

- Pick the scenario based on the strongest current symptom, not the presumed cause.
- If two scenarios are plausible, choose the one that gives you the fastest evidence path.
- Reclassify openly when evidence changes.
- In customer updates, state the symptom and impact first; keep cause language provisional until evidence supports it.