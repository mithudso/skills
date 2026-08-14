# splunk-platform-spl

**Category:** Databases, Data Engineering & Analytics
**Platform:** Claude
**Original Path:** claude/standalone/splunk-platform-spl

## Description
Splunk platform & SPL (Search Processing Language) expert for the log/event ingest-index-search engine and its query language. TRIGGER: Splunk architecture & data lifecycle (forwarders, indexers, search heads, clustering, buckets, SmartStore); writing or optimizing SPL (streaming vs transforming commands, subsearches, CIM, SPL2); search performance (slow searches, time-bounding, tstats, acceleration, summary indexing, Search Job Inspector); knowledge objects (field extraction, tags, event types, macros, lookups); Simple XML / Dashboard Studio dashboards and alerts; Splunk Enterprise Security / SIEM (correlation searches, notable events/findings, risk-based alerting); Splunk Observability Cloud & the OpenTelemetry Collector. SKIP: MongoDB/Atlas logs INTO Splunk → mongodb-operations-expert / mongodb-atlas-expert; telemetry routing (Cribl/Vector/Fluent Bit) → telemetry-pipeline skill; OTel SDK instrumentation, Pino/Sentry/eBPF → devops-observability; Elasticsearch → elasticsearch skill.

---

# Splunk Platform & SPL

Expert reference for the **Splunk platform** (Splunk Enterprise / Splunk Cloud Platform) and **SPL (Search Processing Language)** — the log/event ingest-index-search-analyze engine, its query language, search performance tuning, knowledge objects, dashboards/alerts, the Enterprise Security SIEM layer, and how Splunk Observability Cloud + OpenTelemetry sit relative to it.

This skill routes to on-demand reference files under `references/`. Each reference is self-contained and citation-anchored. Read the one matching the task; for cross-cutting questions, read the overview below first.

> Naming note: Splunk docs now use **manager node / cluster manager** (formerly "master node") and **license manager** (formerly "license master"); older community/blog sources use the legacy terms. Splunk Enterprise Security 8.x renamed **notable event → finding** and **correlation search → detection** (see the ES reference). Volatile claims (licensing, SmartStore constraints, SPL2 scope, ES 8.x terminology, the Cisco/portfolio direction) are stamped `verified-as-of: 2026-06-18` in the references and should be re-verified against current Splunk docs.

## What Splunk is (and is not)

The **Splunk platform** is a schema-on-read engine that ingests machine data (logs/events), indexes it into time-bounded **buckets**, and lets you query it with **SPL**. Its defining design choice is **minimal fixed schema at write, structure applied at read** — most field extraction happens at *search time*, not index time, so you can change how data is interpreted without re-indexing.

A search is a left-to-right pipeline of commands joined by the pipe `|`: terms retrieve events from the index, then commands filter, extract, evaluate, aggregate, and chart them. Performance is fundamentally a **map-reduce** story — work runs in parallel on the indexers until a non-streaming command forces it onto the search head — which is why command *type* and command *order* dominate search cost.

**Splunk Observability Cloud is a different product.** It is a metrics-first, cloud-only SaaS suite (built on SignalFx, OpenTelemetry-native) for infrastructure monitoring, APM/tracing, RUM, and synthetics, queried with **SignalFlow**, not SPL. The log/event platform and Observability Cloud have **separate backends, separate languages, and separate UIs**; OpenTelemetry is the ingestion standard that feeds Observability Cloud by default and can also feed the Splunk platform via HEC. Do not conflate "the Splunk platform" (SPL, indexers) with "Splunk Observability" (SignalFlow, metrics/traces).

## Sub-skill routing table

| If the task is about… | Read |
| --- | --- |
| Forwarders, indexers, search heads, deployment server, indexer/SH clustering, multisite, license model, index-time vs search-time, buckets, SmartStore | `references/splunk-architecture-data-lifecycle.md` |
| The SPL pipeline, streaming vs transforming command types, core commands (search/stats/eval/rex/lookup/join/transaction/tstats/timechart/top/dedup…), subsearches, CIM, SPL2 | `references/spl-language-and-commands.md` |
| Why a search is slow, time-bounding, filter-early, `tstats` over raw search, report/data-model acceleration vs summary indexing, the Search Job Inspector, search modes, workload management | `references/search-performance-optimization.md` |
| Field extraction (search-time vs index-time, IFX, props/transforms), tags, event types, macros, lookups (CSV/KV Store/external/automatic), permissions/knowledge bundle, dashboards (Simple XML, Dashboard Studio), alerts/scheduled searches | `references/knowledge-objects-dashboards-alerts.md` |
| Splunk Enterprise Security (SIEM), correlation searches, notable events/findings, the analyst triage workflow, risk-based alerting (RBA), threat intel, assets & identities, MITRE ATT&CK, the SOAR/UBA portfolio | `references/enterprise-security-siem.md` |
| Splunk Observability Cloud (IM/APM/RUM/Synthetics/Log Observer Connect), SignalFlow vs SPL, the Splunk Distribution of the OpenTelemetry Collector, OTel/OTLP, HEC, Edge/Ingest Processor, the Cisco Data Fabric direction | `references/observability-cloud-otel.md` |

## Cross-references to sibling skills

- **Shipping MongoDB/Atlas logs or metrics INTO Splunk** → `mongodb-operations-expert` and `mongodb-atlas-expert` (they cover the Atlas Log Push / PagerDuty / Datadog / Splunk integration from the MongoDB side; the Atlas Log Push → SIEM path lives under the Atlas/operations hubs). This skill covers the *Splunk* side of receiving and searching that data.
- **Telemetry-routing pipelines** (Cribl, Vector, Fluent Bit, OTel Collector as a routing fabric) → the dedicated telemetry-pipeline skill (if present). Splunk's own Edge Processor / Ingest Processor / Ingest Actions are covered here only as the *Splunk-native* pre-index data-shaping tier; treat external routing fabrics as the upstream pipeline layer that delivers data to Splunk's HEC/forwarders.
- **Generic application observability / OpenTelemetry SDK instrumentation**, Pino structured logging, Sentry, eBPF, Linux perf → `devops-observability`. That hub owns *producing/instrumenting* telemetry; this skill owns the Splunk *platform* that ingests, indexes, and searches it, plus Splunk's own OTel Collector distribution.
- **Elasticsearch / OpenSearch** as the alternative log-analytics platform → the dedicated Elasticsearch skill (if present). Reference it for the comparison; this skill does not cover the Elastic stack.

## Highest-value entry points

- **"Why is my Splunk search slow / how do I optimize it?"** → `references/search-performance-optimization.md` (time-bounding, the "tale of two searches", `tstats`, the Search Job Inspector). This is the densest, most actionable reference.
- **"What's the difference between `stats`, `eventstats`, `streamstats`, and `tstats`?" / "Why is `join` discouraged?"** → `references/spl-language-and-commands.md`.
- **"How does Splunk store data / what is SmartStore?"** → `references/splunk-architecture-data-lifecycle.md`.
- **"What is risk-based alerting and why does it exist?"** → `references/enterprise-security-siem.md`.