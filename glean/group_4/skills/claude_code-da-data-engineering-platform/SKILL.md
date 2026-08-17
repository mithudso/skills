---
description: >-
  Data-engineering & analytics-platform hub (da family) — tooling, pipelines, storage, modeling,
  governance beneath any analysis. TRIGGER: data tools/languages (Python, Polars, pandas, DuckDB,
  SQL, Spark, dbt); pipelines (ELT/ETL, orchestration, Airflow/Dagster, modern data stack,
  lakehouse); streaming analytics (Kafka, Flink, Spark Streaming, windowing, exactly-once);
  semantic layer / headless BI (metric stores, dbt, Cube); data observability
  (freshness/volume/schema/lineage); reverse ETL / operational analytics; real-time OLAP
  (ClickHouse, Druid, Pinot, StarRocks); dimensional modeling (Kimball, star schema, SCD, data
  vault); governance & catalogs (DAMA-DMBOK, DCAM, lineage); data FinOps / cost optimization;
  knowledge graphs & semantic analytics (RDF, ontologies). SKIP: cleaning/EDA/modeling/ML →
  da-analytical-methods; theory → da-1-foundations-theory; lifecycle →
  da-2-data-analysis-lifecycle; viz/comms → da-applied-and-communication.
name: da-data-engineering-platform
version: "1.1.0"
updated: "2026-05-31"
category: custom
tags: [data-analysis, data-engineering, pipelines, olap, platform, governance]
whenToUse:
  - "choosing the right tool or language for a data engineering task — Python, Polars, pandas, DuckDB, dbt, Spark"
  - "designing or debugging a data pipeline: ELT/ETL, orchestration, Airflow, Dagster, workflow scheduling"
  - "setting up or evaluating the modern data stack, lakehouse architecture, or warehouse loading"
  - "building streaming or real-time analytics — Kafka, Flink, Spark Structured Streaming, windowing"
  - "defining a semantic layer, governed metric store, or headless BI with dbt Semantic Layer or Cube"
  - "setting up data observability, data quality monitoring, or freshness/volume/schema-drift detection"
  - "selecting a real-time OLAP engine — ClickHouse, Druid, Pinot, StarRocks — or tuning analytical queries"
  - "modeling a warehouse with Kimball dimensional patterns, star schema, SCDs, or data vault"
  - "establishing data governance, data catalogs, lineage tracking, or data stewardship"
  - "controlling data-platform cost: Snowflake spend, BigQuery slots, Databricks DBUs, FinOps for data"
  - "building knowledge graphs or semantic analytics using RDF, ontologies, or graph-backed queries"
  - "deciding which part of the data/analytics platform owns a problem — routing across the data stack"
whenNotToUse:
  - "MongoDB-specific query, index, schema, or aggregation work — use mongodb-expert"
  - "MongoDB Atlas platform config, Atlas Search, Vector Search, App Services — use mongodb-atlas-expert"
  - "MongoDB operations: Kafka/Spark connectors, CDC, backup, DR, migration — use mongodb-operations-expert"
  - "Kubernetes, Docker, or CI/CD deployment of data pipelines where the infra question dominates — use devops-infra"
  - "AWS-managed data services as primary focus: Glue job config, Redshift cluster sizing, Kinesis shard management — use aws-core or aws-serverless"
  - "Regulatory compliance, GDPR/CCPA enforcement, privacy policy, security review — use security-compliance-auditor"
  - "Statistical/ML modeling methods, EDA, causal inference, forecasting — use da-analytical-methods"
  - "Data acquisition, web scraping, APIs, survey sampling, CDC ingestion design — use da-3-data-acquisition-sampling"
  - "Visualization, dashboarding, reporting, or domain-applied analytics — use da-applied-and-communication"
related_skills:
  - da-1-foundations-theory
  - da-2-data-analysis-lifecycle
  - da-3-data-acquisition-sampling
  - da-analytical-methods
  - da-applied-and-communication
  - mongodb-expert
  - mongodb-operations-expert
  - devops-infra
  - aws-cloud
---

# Data Analysis: Data Engineering & Analytics Platform

The platform substrate beneath data analysis. Where the foundations hub answers
*what an analysis can claim* and the methods hub answers *which technique to run*,
this hub answers *what moves, shapes, stores, serves, and watches the data* so
that analysis is possible at all. It spans the modern data stack end to end: the
languages and engines analysts reach for, the pipelines and orchestration that
land data, the streaming and real-time OLAP layers that serve it at low latency,
the dimensional models and semantic layers that make it reusable, and the
governance, observability, and cost controls that keep it trustworthy and
affordable.

This is a routing hub, not a deep reference in itself. It consolidates 11 former
standalone skills into on-demand reference files.

## How to use this hub

The detail lives in `references/`. This SKILL.md exists to *route* you to the
right file. When a task matches a row in the table below, **Read the listed
`references/` file before giving a deep or authoritative answer** — the routing
line is a one-line summary, not a substitute for the reference content. For broad
"which part of the platform handles this?" questions you can answer from the
table directly; for design choices, trade-offs, tool selection, or anything a
user will act on, load the file first. Multiple references may apply to a single
question (e.g., a pipeline-cost question spans pipelines, OLAP, and FinOps) — load
each relevant one.

## Sub-skill routing table

| Sub-topic | When to load | Reference file |
| --- | --- | --- |
| `da-10-tools-and-languages` | Languages and tools the data analyst reaches for — Python, Polars, pandas, DuckDB, SQL, Spark, dbt | `references/da-10-tools-and-languages.md` |
| `da-13-data-engineering-and-pipelines` | Data engineering and pipelines — ELT/ETL, orchestration (Airflow/Dagster), workflow scheduling, modern data stack, lakehouse, warehouse loading | `references/da-13-data-engineering-and-pipelines.md` |
| `da-14-streaming-analytics` | Streaming analytics — continuous, low-latency processing of unbounded event streams; Kafka, Flink, Spark Structured Streaming, windowing, exactly-once | `references/da-14-streaming-analytics.md` |
| `da-18-semantic-layer-headless-bi` | Semantic layer and headless BI — governed, reusable business metrics, metric stores, dbt Semantic Layer, Cube | `references/da-18-semantic-layer-headless-bi.md` |
| `da-19-data-observability` | Data observability — the five pillars (freshness, volume, schema, lineage, distribution), data quality monitoring, Monte Carlo-style detection | `references/da-19-data-observability.md` |
| `da-20-reverse-etl-operational-analytics` | Reverse ETL and operational analytics — activating modeled warehouse/lakehouse data into SaaS tools, data activation | `references/da-20-reverse-etl-operational-analytics.md` |
| `da-28-realtime-olap-databases` | Real-time OLAP and analytical database engines — the query/storage layer; ClickHouse, Druid, Pinot, StarRocks | `references/da-28-realtime-olap-databases.md` |
| `da-29-dimensional-data-modeling` | Analytical and warehouse data modeling — Kimball dimensional modeling (facts, dimensions), star schema, slowly changing dimensions, data vault | `references/da-29-dimensional-data-modeling.md` |
| `da-30-data-governance-catalogs` | Data governance, catalogs, and discovery — DAMA-DMBOK and DCAM, data catalogs, lineage, stewardship | `references/da-30-data-governance-catalogs.md` |
| `da-37-data-finops-cost-optimization` | FinOps and cost optimization for data/analytics platforms — warehouse/compute spend, cost attribution, query cost tuning | `references/da-37-data-finops-cost-optimization.md` |
| `da-41-knowledge-graphs-and-semantic-analytics` | Knowledge graphs and semantic analytics — RDF, ontologies, graph-backed analytical queries | `references/da-41-knowledge-graphs-and-semantic-analytics.md` |

## Cross-hub routing

This hub owns the data/analytics platform layer. Route out when the question is
not about *what the platform does with data* but about adjacent concerns:

**DA family siblings:**

- **Concepts, vocabulary, measurement, probability, inference, epistemology** →
  `da-1-foundations-theory`.
- **The analysis lifecycle, process frameworks, problem framing, stakeholder
  handoff** → `da-2-data-analysis-lifecycle`.
- **Getting the data in the first place — sources, collection methods, sampling**
  → `da-3-data-acquisition-sampling`.
- **Running the analysis — statistical and ML modeling, EDA, causal inference,
  forecasting** → `da-analytical-methods`.
- **Communicating it — visualization, reporting, and applied/domain analytics** →
  `da-applied-and-communication`.

**Cross-ecosystem deferrals:**

- **MongoDB-specific engine, query, index, aggregation, or schema work** →
  `mongodb-expert`. Atlas platform config → `mongodb-atlas-expert`. MongoDB
  operational pipelines (Kafka/Spark connectors, CDC architecture, backup, DR,
  migration) → `mongodb-operations-expert`. The modeling, governance,
  observability, and cost principles here still inform MongoDB work, but the
  mechanics live in those skills.
- **Kubernetes, Docker, CI/CD deployment of data pipelines** where the
  infrastructure setup (cluster config, networking, Helm charts, container
  builds) is the core question rather than the data semantics → `devops-infra`.
- **AWS-managed analytics services as the primary focus** — Glue job authoring
  and debugging, Redshift cluster configuration and resizing, Kinesis shard
  management, Athena query tuning — → `aws-cloud` (references/aws-core.md) or `aws-serverless`. When the
  data pipeline design question is primary and AWS is just the deployment target,
  stay here and use the pipeline references; when the AWS service mechanics are
  the question, hand off.
- **Governance that escalates to regulatory compliance** — GDPR/CCPA
  enforcement, privacy-by-design policy, security policy review, or compliance
  certification — → `security-review` (references/security-compliance-auditor.md). Data governance frameworks,
  catalogs, and stewardship design stay here.

<!-- cross-hub-map -->
## Cross-hub map — where every data-analytics topic lives

This family is split across these hubs. If a task's deep material is **not** in this hub's Sub-skill
routing table, it is a reference file under a sibling hub below — **activate that hub or `Read` its
`references/<name>.md` directly**. Every former standalone skill in this family is now a reference under one
of these hubs (nothing was deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `da-1-foundations-theory` | Data Analysis Foundations & Theory (hub) | `references/da-1-1-definitions-scope.md`, `references/da-1-1-1-data-analysis-vs-analytics-vs-data.md`, `references/da-1-1-2-analysis-vs-synthesis.md`, `references/da-1-1-3-quantitative-vs-qualitative-analysis.md`, … |
| `da-2-data-analysis-lifecycle` | Data Analysis Lifecycle & Process (hub) | `references/da-2-1-process-frameworks.md`, `references/da-2-1-1-crisp-dm.md`, `references/da-2-1-2-kdd.md`, `references/da-2-1-3-semma.md`, … |
| `da-3-data-acquisition-sampling` | Data Acquisition, Collection & Sampling (hub) | `references/da-3-1-data-sources.md`, `references/da-3-1-1-primary-vs-secondary.md`, `references/da-3-1-2-internal-vs-external.md`, `references/da-3-1-3-structured-semi-structured-unstructured.md`, … |
| `da-analytical-methods` | Data Analytical Methods (cleaning, EDA, modeling, ML, causal, time-series) | `references/da-4-data-cleaning-preparation.md`, `references/da-5-exploratory-data-analysis.md`, `references/da-6-statistical-modeling.md`, `references/da-7-machine-learning.md`, … |
| `da-data-engineering-platform` | Data Engineering & Analytics Platform (pipelines, OLAP, modeling, governance) | `references/da-10-tools-and-languages.md`, `references/da-13-data-engineering-and-pipelines.md`, `references/da-14-streaming-analytics.md`, `references/da-18-semantic-layer-headless-bi.md`, … |
| `da-applied-and-communication` | Applied Analytics, Visualization, Communication & Ethics | `references/da-8-data-visualization.md`, `references/da-9-reporting-communication.md`, `references/da-11-ethics-and-privacy.md`, `references/da-21-product-analytics.md`, … |
| `data-analytics` | Family ROUTER — entry point for all da-* sub-hubs | (this file's parent hub) |
