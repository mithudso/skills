---
name: elasticsearch-opensearch
description: >-
  Elasticsearch / OpenSearch & the ELK/Elastic Stack — the open-source search + log-analytics engine. TRIGGER: Elasticsearch-vs-OpenSearch choice & the 2021 fork/licensing (SSPL, Elastic License, AGPL, Linux Foundation); the stack (Beats, Elastic Agent/Fleet, Logstash, Kibana / OpenSearch Dashboards, Data Prepper, ingest pipelines); engine internals (Lucene, inverted index, doc_values, refresh/flush/translog, shards & replicas, mappings, analyzers, text-vs-keyword); Query DSL (bool/term/match/range, query-vs-filter context, aggregations, BM25, ES|QL vs OpenSearch PPL/SQL); operations (ILM/ISM, hot-warm-cold-frozen tiers, shard sizing, anti-patterns); log-analytics/observability/SIEM positioning vs Splunk. SKIP: Atlas Search & Atlas vector search → mongodb-atlas-expert; RAG retrieval → ai-rag-retrieval; Splunk platform itself → splunk-platform-spl; telemetry routers (Cribl/Vector/Fluent Bit/OTel Collector) → telemetry-pipeline; OTel SDK instrumentation → devops-observability.
version: 1.0.0
updated: 2026-06-18
category: observability
whenToUse: Use when the task involves Elasticsearch, OpenSearch, or the ELK/Elastic Stack — choosing between Elasticsearch and OpenSearch or explaining the 2021 fork and the SSPL/ELv2/AGPL licensing timeline; the stack components and data flow (Beats, Elastic Agent/Fleet, Logstash, Kibana, OpenSearch Dashboards, Data Prepper, ingest pipelines); engine internals (Lucene, segments, inverted index, doc_values, refresh/flush/translog, shards/replicas, mappings, analyzers); writing or reasoning about the Query DSL, aggregations, BM25 relevance scoring, or ES|QL / OpenSearch PPL/SQL; operations and performance (ILM/ISM, data tiers, shard sizing, indexing/query tuning, anti-patterns); or positioning Elasticsearch/OpenSearch as the open-source log-analytics/observability/SIEM counterpart to Splunk. Not for MongoDB Atlas Search or Atlas Vector Search (use mongodb-atlas-expert), the Splunk platform itself, telemetry-routing pipelines, or OTel-SDK instrumentation.
keywords:
  - Elasticsearch
  - OpenSearch
  - ELK Stack
  - Elastic Stack
  - Apache Lucene
  - Query DSL
  - ES|QL
  - PPL
  - BM25
  - inverted index
  - shards and replicas
  - mappings and analyzers
  - text vs keyword
  - aggregations
  - Index Lifecycle Management
  - ILM
  - ISM
  - data tiers
  - hot-warm-cold-frozen
  - searchable snapshots
  - Logstash
  - Beats
  - Elastic Agent
  - Fleet
  - Kibana
  - OpenSearch Dashboards
  - Data Prepper
  - SSPL
  - Elastic License
  - AGPL
  - OpenSearch fork
  - log analytics
  - observability
  - SIEM
  - shard sizing
  - oversharding
  - deep pagination
  - refresh flush translog
tags:
  - observability
  - search
  - log-analytics
  - elasticsearch
  - opensearch
  - data-platform
---

# Elasticsearch / OpenSearch & the ELK/Elastic Stack

Expert reference for the open-source distributed search + log-analytics engine — **Elasticsearch** and its 2021 Apache-2.0 fork **OpenSearch** — and the surrounding stack (Beats/Elastic Agent, Logstash, Kibana / OpenSearch Dashboards). Both engines wrap **Apache Lucene**; below the 7.10 fork point they are essentially identical, and they have diverged meaningfully since. This skill is the **open-source counterpart** to the commercial-SIEM sibling `splunk-platform-spl`.

> Volatile claims (versions, licensing state, vendor-feature landscape, pricing) are stamped **verified-as-of 2026-06-18** and concentrated in the references; re-verify before quoting version numbers or license terms. A major underlying-model migration or a new ES/OpenSearch major line is a refresh trigger.

## When this skill applies

Use it to: pick between Elasticsearch and OpenSearch and explain the licensing fork; reason about the engine (Lucene, segments, inverted index, shards, mappings/analyzers); write or debug Query DSL, aggregations, BM25 relevance, ES|QL or OpenSearch PPL/SQL; design operations (ILM/ISM, data tiers, shard sizing, tuning); or position the stack for log analytics / observability / SIEM against Splunk.

## Routing table — load the reference for the task

| If the task is about… | Read this reference |
|---|---|
| Elasticsearch-vs-OpenSearch choice, the 2021 fork, SSPL / Elastic License v2 / AGPL timeline, OpenSearch → Linux Foundation, version landmarks | `references/licensing-fork-and-choosing.md` |
| Lucene, segments & merging, inverted index, doc_values, BKD, refresh/flush/commit/translog, shards & replicas, the document/index model, mappings, analyzers, text-vs-keyword | `references/engine-internals-lucene-mappings.md` |
| Query DSL (`bool`/`term`/`match`/`range`), query-vs-filter context, aggregations, BM25 relevance scoring, the `_search` API, ES\|QL vs OpenSearch PPL/SQL | `references/query-dsl-aggregations-relevance.md` |
| Stack components & data flow (Beats, Elastic Agent/Fleet, Logstash, ingest pipelines, Kibana, OpenSearch Dashboards, Data Prepper) | `references/stack-components-and-data-flow.md` |
| ILM / ISM, hot-warm-cold-frozen tiers, searchable snapshots, data streams, shard sizing, indexing/query tuning, anti-patterns | `references/operations-ilm-tiers-performance.md` |
| Log-analytics / observability / SIEM use cases and positioning vs Splunk (cost model, schema-on-write vs schema-on-read) | `references/log-analytics-observability-vs-splunk.md` |

## Core orientation (the 60-second model)

- **One engine, two distributions.** Apache Lucene is the actual search engine (inverted index + storage); Elasticsearch/OpenSearch add the distributed layer (sharding, replication, JSON-over-HTTP API, cluster coordination). A **shard = one Lucene index = a set of immutable segments**. OpenSearch forked Elasticsearch 7.10.2 in 2021 and keeps Lucene as its core, so the *mechanics* are shared even where tunable defaults and higher-level features have drifted. See `engine-internals-lucene-mappings.md`.
- **The licensing fork is the defining event.** Elastic relicensed Elasticsearch/Kibana from Apache-2.0 to SSPL + Elastic License v2 starting 7.11 (Jan 2021), targeting AWS's managed-service use; AWS forked the last Apache-2.0 code into OpenSearch; in 2024 Elastic added AGPLv3 (making Elasticsearch OSI-open again) and OpenSearch moved to the Linux Foundation. Licensing is usually the cleanest decision criterion. See `licensing-fork-and-choosing.md`.
- **Two index data structures power everything.** The **inverted index** (term → postings list) answers full-text search; **doc_values** (a columnar, on-disk, per-document store) powers sorting and aggregations; numeric/date/geo fields add a **BKD tree** for range queries. The `text`-vs-`keyword` field choice maps directly onto these. See `engine-internals-lucene-mappings.md` and `query-dsl-aggregations-relevance.md`.
- **Query context vs filter context** is the single most useful querying distinction: query context scores (relevance), filter context is binary yes/no, skips scoring, and is cacheable — put structured conditions in `filter`. **BM25** (not TF/IDF) is the default similarity since Elasticsearch 5.0. See `query-dsl-aggregations-relevance.md`.
- **Time-series data wants data streams + lifecycle management.** ILM (Elastic) / ISM (OpenSearch) roll indices over and migrate them across hot → warm → cold → frozen tiers, with the frozen tier backed by **searchable snapshots** on object storage for cheap long retention. Most operational pain is **oversharding**, **mapping explosion**, and **deep pagination** — all avoidable. See `operations-ilm-tiers-performance.md`.

## Cross-references (do not duplicate — route)

- **MongoDB Atlas Search** (Lucene-backed `$search` embedded in Atlas) and **Atlas Vector Search** (`$vectorSearch`) → `mongodb-atlas-expert`. General **RAG retrieval** architecture → `ai-rag-retrieval`. Atlas Search is a Lucene cousin, but this skill does not cover the Atlas surface.
- **The Splunk platform & SPL** (the commercial-SIEM / log-analytics counterpart) → `splunk-platform-spl`. This skill cross-references it for the head-to-head comparison; `splunk-platform-spl` reciprocally forward-references this skill. Splunk-internal SPL, knowledge objects, and Splunk ES are out of scope here.
- **Telemetry-routing pipelines** that feed Elasticsearch/OpenSearch (Cribl, Vector, Fluent Bit, OpenTelemetry Collector — and Logstash *beyond* its stack role) → a sibling `telemetry-pipeline` skill (route there once it exists). Logstash and Data Prepper are covered here only as ELK/OpenSearch stack components.
- **OpenTelemetry SDK instrumentation**, Pino/Sentry/eBPF, Node/Linux observability → `devops-observability`.
