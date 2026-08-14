<!-- hub-reference-banner -->
> **Reference file — part of the `da-data-engineering-platform` hub.** Formerly the standalone `da-13-data-engineering-and-pipelines` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-13-data-engineering-and-pipelines
title: Data Engineering and Pipelines
version: "1.0.0"
updated: "2026-05-30"
category: data-analysis
origin: local
description: >
  Data engineering for analysts and data engineers in 2026 — workflow
  orchestration (Airflow, Dagster, Prefect), dbt for analytics engineering,
  lakehouse architectures (Delta, Iceberg, Hudi), batch vs streaming
  (Lambda / Kappa), data contracts as runtime enforcement, data quality
  testing (Great Expectations, Soda, dbt tests), and pipeline observability
  (OpenLineage, DataHub, Marquez).
  TRIGGER: building or operating a data pipeline; choosing between Airflow,
  Dagster, and Prefect; setting up dbt models, tests, or exposures; picking
  Delta/Iceberg/Hudi for a lakehouse; designing batch vs streaming for an
  analytics workload; introducing data contracts; setting up data quality
  tests; instrumenting pipelines with OpenLineage or DataHub; debugging a
  failing DAG or model run.
  SKIP: pure ML pipelines (use da-7-machine-learning + da-17 feature stores);
  application-level CDC from a transactional database (use mongodb-change-streams
  or mongodb-kafka-connector); ETL inside a Spark notebook only (use
  da-10-tools-and-languages); BI dashboard layer (da-8-data-visualization).
triggers:
  - data engineering
  - data pipelines
  - Airflow
  - Dagster
  - Prefect
  - dbt analytics engineering
  - lakehouse
  - Delta Lake
  - Apache Iceberg
  - Apache Hudi
  - data contracts
  - data quality testing
  - Great Expectations
  - OpenLineage
  - DataHub
  - Marquez
  - Lambda architecture
  - Kappa architecture
keywords:
  - Airflow
  - Dagster
  - Prefect
  - dbt
  - Delta-Lake
  - Iceberg
  - Hudi
  - lakehouse
  - data-warehouse
  - data-lake
  - data-contracts
  - schema-registry
  - Great-Expectations
  - Soda
  - OpenLineage
  - DataHub
  - Marquez
  - Lambda-architecture
  - Kappa-architecture
  - analytics-engineering
  - semantic-layer
when_to_use:
  - Choosing an orchestrator (Airflow / Dagster / Prefect) for a new project
  - Designing dbt model layering, tests, exposures, semantic layer
  - Picking between Delta, Iceberg, and Hudi for a lakehouse
  - Choosing between batch and streaming for an analytics workload
  - Introducing data contracts between producer and consumer teams
  - Setting up data quality tests at ingest, transform, or consumption layers
  - Instrumenting pipelines with OpenLineage / DataHub / Marquez
  - Debugging a failing DAG or dbt model run
when_not_to_use:
  - Pure ML training pipelines — use da-7-machine-learning + da-17
  - Application CDC from MongoDB — use mongodb-change-streams / mongodb-kafka-connector
  - Streaming-engine internals (Kafka / Flink) — use da-14-streaming-analytics
  - BI dashboard layer — use da-8-data-visualization
related_skills:
  - da-3-data-acquisition-sampling
  - da-4-data-cleaning-preparation
  - da-10-tools-and-languages
  - da-14-streaming-analytics
  - da-17-feature-engineering-and-feature-stores
  - mongodb-kafka-connector
  - mongodb-atlas-stream-processing
---

# Data Engineering and Pipelines

The layer between raw data and the analyst's notebook. This skill covers the orchestration, transformation, storage, and observability stack that moves data from sources into analysis-ready form on a schedule.

## When to use this skill

Activate when the user:
- is building or operating a data pipeline
- needs to pick an orchestrator (Airflow vs Dagster vs Prefect)
- is structuring dbt models, tests, or exposures
- is choosing a lakehouse table format (Delta vs Iceberg vs Hudi)
- is designing batch vs streaming for an analytics workload
- is introducing data contracts between teams
- is instrumenting pipelines with OpenLineage / DataHub

## When NOT to use this skill

- Pure ML training pipelines → `da-7-machine-learning` + `da-17-feature-engineering-and-feature-stores`
- App-level CDC from MongoDB → `mongodb-change-streams` / `mongodb-kafka-connector`
- Streaming engine internals → `da-14-streaming-analytics`
- BI dashboard layer → `da-8-data-visualization`

---

## Workflow orchestration: Airflow vs Dagster vs Prefect

The 2026 landscape has three serious contenders. They share the goal — schedule + retry + monitor DAGs of work — but disagree on the unit of abstraction.

| Tool | Abstraction | Strengths | Weaknesses |
|---|---|---|---|
| **Airflow** | Task-based (TaskFlow API in 2.x; asset-aware scheduling in 3.x) | Huge ecosystem, every tool has an Airflow provider, mature at scale | Task-DAG model awkward for asset-oriented pipelines, slow scheduling at 1000+ DAGs without optimization |
| **Dagster** | **Asset-based** (software-defined assets) | Lineage is first-class, asset materialization is the unit, excellent dev ergonomics, modern UI | Smaller ecosystem than Airflow; some teams find the asset model novel |
| **Prefect** | Task + flow with dynamic discovery | Pythonic, runs locally trivially, hybrid cloud control plane | Smaller production deployments; orchestration semantics differ from Airflow enough to retrain |

**The decision rule**: if your pipeline naturally maps to *"these tables / models / files should exist and be fresh"*, Dagster's asset model fits. If you have many existing Airflow DAGs or need every conceivable integration, stay on Airflow. If you want lightweight, Pythonic, with cloud-managed control plane, Prefect is good.

Airflow 3.0 added asset-aware scheduling (building on the Datasets concept introduced in Airflow 2.4), narrowing the gap with Dagster, but the DAG-of-tasks remains the underlying model.

---

## dbt — analytics engineering

dbt turned SQL transformations from notebooks into engineered software with models, tests, docs, and lineage. By 2026 it's the de-facto standard for the transformation layer in a modern data stack.

### Layering convention

```
sources (raw)
   ↓
staging       — rename, type-cast, light cleaning; 1:1 with sources
   ↓
intermediate  — joins, aggregations, business logic that's not yet a mart
   ↓
marts         — business-facing tables; one per business domain
```

Discipline: every model lives in one layer; cross-layer references are linear (never staging → marts → staging).

### Tests

- `not_null`, `unique`, `accepted_values`, `relationships` — built-ins
- `dbt_expectations` package — Great-Expectations-style assertions inside dbt
- `dbt_utils` package — equal_rowcount, recency, fewer_rows_than, expression_is_true
- Custom singular tests for project-specific invariants

### Exposures and the semantic layer

- **Exposures** — declarative dependencies from dbt models to downstream consumers (dashboards, ML models, reports). Make lineage end-to-end.
- **Semantic layer** (dbt Cloud / MetricFlow) — define metrics once, query them from any BI tool with consistent semantics. Replaces "what does ARR mean here?" arguments.

### dbt Core vs dbt Cloud vs dbt Fusion (2025)

- **dbt Core** — open-source CLI. The base.
- **dbt Cloud** — hosted scheduler + IDE + semantic layer + observability.
- **dbt Fusion (2025+)** — a Rust-based reimplementation of the dbt engine, much faster parsing/compilation for large projects. Watch this — it's the future of dbt internals.

---

## Lakehouse: Delta vs Iceberg vs Hudi

A lakehouse stores data in open file formats on object storage (S3, ADLS, GCS) but with the transactional and schema guarantees of a warehouse. Three table formats compete:

| Format | Origin | Strengths | Pick when |
|---|---|---|---|
| **Delta Lake** | Databricks (open-sourced 2019) | Most mature, deep Spark integration, Photon, Unity Catalog | Already on Databricks; or want the most mature option |
| **Apache Iceberg** | Netflix (2018) | Multi-engine first-class (Spark, Trino, Flink, Snowflake, Athena), better partition evolution | Multi-engine reads/writes; vendor-neutral |
| **Apache Hudi** | Uber (2016) | Best for streaming-heavy upserts, incremental processing | High-frequency CDC writes into the lake |

By 2026 the **Iceberg vs Delta** race is the headline. Snowflake added native Iceberg support; Databricks added Iceberg interop on Delta tables (Uniform). The walls between formats are coming down.

### Lakehouse vs warehouse vs lake

| | Lake | Warehouse | Lakehouse |
|---|---|---|---|
| Storage | Object storage | Proprietary | Object storage |
| Format | Parquet/JSON loose | Internal | Delta/Iceberg/Hudi |
| Transactions | No | Yes | Yes |
| Schema enforcement | No | Yes | Yes |
| Compute coupling | None | Tight | Loose |
| Cost | Cheapest | Highest | Cheap-ish |

The lakehouse exists because warehouses got expensive and lakes were unreliable. Modern teams put cheap historical / semi-structured / wide-table data in a lakehouse and use a warehouse only when query latency demands it.

---

## Batch vs streaming: Lambda vs Kappa

The architectural question for any near-real-time analytics workload.

### Lambda architecture

Two parallel paths:
- **Batch layer** — slow, comprehensive, source of truth (Spark, Hive, dbt)
- **Speed layer** — fast, approximate, recent (Flink, Kafka Streams)
- **Serving layer** — combines both for queries

Pros: separates concerns; batch is the immutable record of truth.
Cons: two codebases, two systems to maintain, reconciliation drift.

### Kappa architecture

One streaming pipeline that handles both real-time and reprocessing. Reprocessing = replay the same stream from offset 0. Simplifies the codebase but pushes complexity into the streaming engine (state management, exactly-once semantics).

### When each fits

- **Lambda** when batch correctness matters and the speed layer is a "best-effort preview" (financial reporting, regulated industries)
- **Kappa** when the streaming engine is mature enough and you'd rather maintain one pipeline (modern Kafka + Flink shops)
- **Neither — pure batch** when freshness > 1 hour is fine (the vast majority of analytics)

The 2025-2026 trend: most teams started with Lambda, got tired of the double-maintenance, and migrated to Kappa once Flink / Kafka Streams stabilized.

---

## Data contracts as runtime enforcement

A data contract is a producer's commitment about schema, semantics, freshness, and quality, enforced at the boundary between systems.

Why they exist: silent schema drift breaks downstream models hours later. By the time the report shows wrong numbers, the analyst can't tell whether the source changed, the transform broke, or the dashboard cache is stale.

### Enforcement mechanisms (2026)

- **Schema Registry** (Confluent, Apicurio) — Avro / Protobuf / JSON Schema versioning at the topic level. Producer can't publish an incompatible schema.
- **dbt model contracts** (dbt 1.5+) — `model.contract: enforced: true`. dbt fails the build if column names/types drift.
- **Pact** — consumer-driven contract tests, originally for APIs, applies to data too.
- **dbt-checkpoint / SqlFluff** — pre-commit linting catches contract violations before merge.

### Contract metadata (a working baseline)

A useful contract specifies:
- **Schema** — column names, types, nullability
- **Semantics** — what each field means (units, enum values, dimensions)
- **Freshness SLA** — max staleness before alert
- **Volume SLA** — expected row count with tolerance
- **Quality SLA** — % null, % unique, value ranges
- **Owner** — who responds when it breaks
- **Versioning rule** — breaking change requires a major version bump and consumer migration window

---

## Data quality testing

Three layers of tests every pipeline should have:

1. **Source tests** (ingest layer) — assert raw data matches expectations. Use Great Expectations, Soda, or dbt source tests.
2. **Transform tests** (transformation layer) — assert dbt models produce expected shapes. Built-in dbt tests + dbt_expectations.
3. **Consumption tests** (BI / report layer) — assert the consumer's invariants. Often custom SQL.

### Tool comparison

| Tool | Strengths | Watch out for |
|---|---|---|
| **Great Expectations** | Most comprehensive, batteries-included | Heavy; learning curve |
| **Soda Core / Soda Cloud** | Lightweight; SodaCL declarative YAML | Less expressive than GE for custom logic |
| **dbt tests** | Native to dbt projects | Limited to what dbt sees |
| **Monte Carlo / Bigeye / Lightup** | Data observability platforms with ML anomaly detection | Cost; vendor lock-in |

### Critical tests every pipeline needs

- **Row count delta vs prior run** — sudden 10× or 0.1× = upstream broke
- **Null rate per column vs baseline** — sudden spike usually means schema drift
- **Distribution drift on key columns** — KS test or PSI on numerical columns
- **Referential integrity** — every FK has a matching PK
- **Freshness** — `max(updated_at) >= now() - threshold`

---

## Pipeline observability

The "is my pipeline doing the right thing?" layer.

- **OpenLineage** — open standard for emitting lineage events. Airflow, Dagster, Spark, dbt, and Flink all emit OL events natively in 2026. Stores in Marquez or DataHub.
- **DataHub (LinkedIn / Acryl)** — open-source data catalog with lineage, ownership, schema, and ML feature metadata. Heavy but comprehensive.
- **Marquez** — the reference OpenLineage backend. Lighter than DataHub; good for "just show me lineage."
- **Datadog / Monte Carlo / Bigeye / Lightup** — commercial data-observability with anomaly detection on freshness, volume, schema, distribution.

CI/CD for data pipelines (2026 baseline):
- Every PR runs `dbt build --select state:modified+` on a clone of prod
- Schema changes are detected by `dbt list --resource-type:model --output-keys columns`
- Production deploys via a Slim CI pattern (only run downstream of changed models)

---

## Anti-patterns

1. **Two-codebase Lambda** that the team can't keep in sync — pick Kappa or pure-batch.
2. **dbt models that aren't tested** — unique + not_null on every primary key, recency on every fact table, minimum baseline.
3. **Adding columns to a published table without a contract bump** — silent breakage downstream.
4. **Schedule the DAG, hope for the best** — no freshness alert, no SLA, no on-call. The pipeline is now incident-prone by design.
5. **Putting business logic in BI layer instead of dbt** — the same metric defined 3 different ways in 3 dashboards.
6. **Pipeline runs on prod data but isn't reproducible** — no `dbt seed` of canonical test data, no clone for dev.
7. **Notebook is the production pipeline** — the analyst's notebook should not be on the cron. Promote it to dbt + orchestrator.

---

## References

1. *Data Pipelines Pocket Reference*, Densmore (O'Reilly, 2021).
2. *Fundamentals of Data Engineering*, Reis & Housley (O'Reilly, 2022).
3. dbt docs — https://docs.getdbt.com/
4. Airflow 3.0 announcement and AIP-77 — https://airflow.apache.org/
5. Dagster software-defined assets — https://docs.dagster.io/
6. Prefect 3 — https://docs.prefect.io/
7. Delta Lake — https://delta.io/
8. Apache Iceberg — https://iceberg.apache.org/
9. Apache Hudi — https://hudi.apache.org/
10. OpenLineage — https://openlineage.io/
11. DataHub — https://datahubproject.io/
12. Marquez — https://marquezproject.ai/
13. Great Expectations — https://docs.greatexpectations.io/
14. *Data Contracts*, Andrew Jones (O'Reilly, 2023). The book.
