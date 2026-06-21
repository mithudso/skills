---
description: >-
  Data acquisition, collection & sampling hub (da family, stage 3). TRIGGER: data-source taxonomy (primary/secondary, internal/external, structured/semi/unstructured); collection methods (web scraping/crawling, robots.txt/legality, APIs & data feeds, REST/GraphQL/gRPC, pagination/rate-limit/auth, web/app analytics instrumentation, event tracking); data-collection overview; survey & sampling-frame design, response bias; sampling methods (simple random, stratified, cluster, systematic, convenience, quota, snowball); sample size/power; acquisition plumbing — API ingestion, DB extraction/CDC, streaming ingest (Kafka/Kinesis/Pub-Sub), ETL/ELT & modern data stack. SKIP: theory → da-1-foundations-theory; lifecycle → da-2-data-analysis-lifecycle; methods/ML → da-analytical-methods; pipelines/platform → da-data-engineering-platform; viz/comms → da-applied-and-communication.
name: da-3-data-acquisition-sampling
version: 2.0.1
updated: "2026-06-21"
model: claude-opus-4-8
effort: high
whenNotToUse: >
  Use da-1-foundations-theory for measurement/probability/inference theory; da-2-data-analysis-lifecycle
  for process frameworks, problem framing, and lifecycle/process work; da-analytical-methods for
  cleaning, EDA, modeling, ML, causal, and time-series methods; da-data-engineering-platform for
  production pipeline/streaming/ingestion engineering, OLAP, warehouse modeling, governance and
  observability; da-applied-and-communication for visualization, reporting, communication, and ethics.
related_skills:
  - da-1-foundations-theory
  - da-2-data-analysis-lifecycle
  - da-analytical-methods
  - da-data-engineering-platform
  - da-applied-and-communication
---
# Data Acquisition and Sampling

Skill covers **third stage** of data analysis curriculum: getting data into form analysis can operate on, constructing sample that supports valid inference about target population. Sits between lifecycle/problem framing (da-2-data-analysis-lifecycle) and downstream cleaning, EDA, modeling (da-analytical-methods). Activities intertwined — source choice constrains sampling design; sampling plan determines which sources acceptable.

Common mistake: treat acquisition as logistics problem — "just pull data" — then discover at analysis stage that population was wrong, sample frame had coverage gaps, schema drifted mid-pull, or file format broke planned query. Acquisition/sampling stage owns catching these failures *before* they contaminate downstream work.

---

## Sub-skill routing table

Hub absorbs 9 former standalone skills as on-demand reference files. Task matches row → **Read listed `references/` file** before answering — don't rely on table alone for depth.

| Sub-topic | When to load | Reference file |
| --- | --- | --- |
| `da-3-1-data-sources` | Expert knowledge on data sources as a category within data collection & acquisition. | `references/da-3-1-data-sources.md` |
| `da-3-1-1-primary-vs-secondary` | Primary vs. secondary data sources — provenance, fit, and trade-offs. | `references/da-3-1-1-primary-vs-secondary.md` |
| `da-3-1-2-internal-vs-external` | The internal-vs-external dimension of data sources in analysis projects. | `references/da-3-1-2-internal-vs-external.md` |
| `da-3-1-3-structured-semi-structured-unstructured` | Structured, semi-structured, and unstructured data as data-source types. | `references/da-3-1-3-structured-semi-structured-unstructured.md` |
| `da-3-2-collection-methods` | Categorical overview and comparison of data collection methods across the full taxonomy. | `references/da-3-2-collection-methods.md` |
| `da-3-2-5-web-scraping-crawling` | Web scraping and crawling as a data collection method — tooling, legality, ethics. | `references/da-3-2-5-web-scraping-crawling.md` |
| `da-3-2-6-apis-data-feeds` | Collecting data via APIs and structured data feeds for analysis purposes. | `references/da-3-2-6-apis-data-feeds.md` |
| `da-3-2-7-web-app-analytics-instrumentation` | Event instrumentation for web and mobile apps, end to end. | `references/da-3-2-7-web-app-analytics-instrumentation.md` |
| `da-3-data-collection-acquisition` | Overview and orientation for Data Collection & Acquisition as a curriculum stage. | `references/da-3-data-collection-acquisition.md` |

---

## 1. Data Source Taxonomy

First decision: where data comes from. Three orthogonal axes characterize any source.

### 1.1 Primary vs secondary

**Primary data** collected directly to answer current research question. Instrument designed for this purpose, participants selected for this study, format matches analysis plan. Examples: surveys you designed/fielded, interviews you conducted, sensor data your instrumentation captured, A/B test exposure logs, custom user research sessions.

**Secondary data** collected by someone else for different purpose, being reused. Instrument may not fit your question, sampling frame may differ, coding may not align with constructs. Examples: government statistics (Census, BLS), commercial panels (Nielsen, Comscore), academic data archives (ICPSR), internal datasets gathered for other projects, third-party API exports.

Strong analyses combine both: primary characterizes participants of interest with right instrument; secondary provides counterfactual/population-level context primary collection can't afford. Trade-off: purpose-fit (primary) vs speed/scale/cost (secondary).

### 1.2 Structured vs semi-structured vs unstructured

- **Structured**: row/column tabular data, fixed schema, typed columns — relational databases, data warehouses, CSV/Parquet. Examples: transaction records, sales line items, test scores, sensor readings.
- **Semi-structured**: hierarchical/self-describing — discoverable structure, doesn't fit rows/columns cleanly. Examples: JSON API responses, XML, log files, NoSQL documents, Avro/Protobuf messages.
- **Unstructured**: free-form text, images, audio, video, raw documents, no inherent schema. Requires feature extraction (OCR, ASR, embeddings) or model-based interface (LLM, CV model) to participate in tabular analysis.

2020s blurred this taxonomy — vector embeddings/LLMs reduce cost of operating on unstructured data — but physics still apply: closer source is to structured form, cheaper and more deterministic analysis. Schema-on-read (data lakes) defers structure to query time; schema-on-write (warehouses) enforces at load time.

### 1.3 Internal vs external

Internal sources (production databases, event logs, CRM, application telemetry) typically more reliable, more granular, more legally defensible — but reflect your org's operations and may not generalize. External sources (APIs, public datasets, scraped pages, panels) broaden population but introduce coverage uncertainty, licensing risk, schema drift you don't control.

---

## 2. API Ingestion

APIs most common acquisition path for external structured data. API style choice shapes everything downstream — pagination, error handling, schema evolution, auth all inherit from this decision.

### 2.1 REST, GraphQL, and gRPC

- **REST** (HTTP + JSON, resource-oriented): broadest compatibility, mature HTTP caching, well-defined contract via OpenAPI. Stateless, cache-friendly, easy to debug with `curl`. Default for public APIs and most acquisition pipelines. Downside: over-fetching (whole resource) and under-fetching (multiple round-trips to assemble what one screen needs).
- **GraphQL** (typed query language): client asks for exactly needed fields in one request. Single endpoint, strongly typed schema, introspectable. Eliminates over- and under-fetching. Downside: harder to cache at HTTP layer, rate limiting requires query-cost analysis (GitHub's GraphQL API uses point budgets, not request counts), auth becomes per-field.
- **gRPC** (HTTP/2 + Protocol Buffers, binary): code-generated clients from `.proto`, multiplexed bidirectional streaming, payloads ~3-10x smaller than JSON. Best for service-to-service internal traffic. Not directly browser-accessible without gRPC-Web. Use when throughput/latency matter more than human readability.

Common 2026 pattern: REST for public API surface, GraphQL for BFF/frontend layer, gRPC between internal services. From acquisition perspective, gRPC rare for external data — mostly REST and GraphQL.

### 2.2 Authentication

- **API key**: shared secret in header (`X-API-Key:`) or query string. Simplest. Identifies app, not user. Use only over TLS. Rotate on schedule. Never embed in client-side code or commit to source control.
- **OAuth 2.0**: delegated auth with access tokens (short-lived) and refresh tokens (long-lived). Authorization Code with PKCE for user-context apps; Client Credentials for service-to-service. Always validate scopes server-side.
- **JWT**: signed bearer token (HS256, RS256, ES256). Carries claims (subject, expiry, scopes, custom). Stateless verification — no DB lookup if signature checks out. Keep tokens short-lived (5-60 min), use refresh tokens for renewal, verify `alg` header to avoid `alg: none` attacks.
- **mTLS / certificate auth**: client presents cert signed by trusted CA. Used for high-trust internal traffic and some financial/healthcare APIs.

For acquisition pipelines, treat refresh tokens as most sensitive secret in system. Typically months-long lifetime, full account access. Store encrypted, rotate on suspicion, log every refresh.

### 2.3 Pagination

- **Offset/limit** (`?offset=200&limit=50`): simple but breaks when underlying data changes during paging (rows shift across boundaries). Acceptable for read-mostly small datasets.
- **Page number** (`?page=3&per_page=50`): same trade-offs as offset/limit, more readable.
- **Cursor-based** (`?after=<opaque_token>`): server returns token pointing to next slice. Stable under concurrent writes, efficient indexed lookups — preferred for high-volume APIs (Stripe, GitHub, Slack, Twitter/X).
- **Keyset / seek pagination** (`?after_id=12345`): cursor-based but cursor is domain value (sort key plus tiebreaker). Cheap on indexed columns.

When acquiring large datasets, persist cursor after every page so partial failure can resume. Never assume API returns same total count on re-pull.

### 2.4 Rate limiting

- **Fixed window**: N requests per minute, counter resets at boundary. Gameable by bursting at window edges.
- **Sliding window**: rolling N requests over last 60 seconds.
- **Token bucket**: refill rate plus bucket size, allows bursts up to bucket capacity.
- **Cost-based** (GraphQL): each query assigned point cost based on complexity; budget per hour. Monitor via `Authorization` and `X-RateLimit-*` response headers.

Build acquisition clients with exponential backoff on `429 Too Many Requests` and `503 Service Unavailable`, respecting `Retry-After` header. Cap total retries; surface sustained failures rather than retrying forever.

### 2.5 Webhooks

Webhooks = inverse of polling — provider POSTs to your endpoint on event. Always verify signature header (HMAC-SHA256 over raw body with shared secret) before trusting payload. Webhooks can be replayed by attackers without verification; treat unsigned webhook as untrusted input.

---

## 3. Web Scraping

Web scraping = acquisition without contract. Works until page structure changes or site owner objects. Use only when no API exists and legal/ethical posture is sound.

### 3.1 Tooling

- **`requests` + BeautifulSoup**: minimal stack for static HTML. Fast, simple, no JS execution.
- **Scrapy**: full crawling framework with built-in concurrency, throttling, retries, item pipelines, middleware. Best for multi-page crawls and structured extraction at scale.
- **Playwright / Selenium**: headless browser automation. Required when page renders via JavaScript or hides data behind auth flows. Slower (~0.5-1 page/sec per worker) and more resource-intensive than static scraping.
- **`curl_cffi` and similar TLS-impersonation tooling**: needed when target site fingerprints TLS handshakes to block bots. Treat as signal site doesn't want to be scraped.

### 3.2 Legality

US law as of 2026 turns on hiQ Labs v. LinkedIn line: scraping publicly accessible data generally doesn't violate CFAA. But that's only the federal anti-hacking question. Full risk surface:

- **CFAA**: bypassing auth, ignoring access controls, or manipulating URLs to reach gated content can still create CFAA exposure.
- **Terms of Service**: violating ToS rarely criminal but can trigger civil contract claims and account termination.
- **Copyright**: bulk reproduction of copyrighted content can infringe even if publicly accessible.
- **GDPR / CCPA**: scraping personal data triggers data-protection obligations regardless of whether data was public. EU position much stricter than US.
- **Database rights** (EU): Database Directive grants sui generis rights to maker of database; bulk extraction can infringe without copyright.

### 3.3 Ethics and good-faith practice

- Respect `robots.txt` — public machine-readable signal of site owner expectations. Ignoring it doesn't make scraping illegal but undermines good-faith argument.
- Set descriptive `User-Agent` with contact info.
- Throttle requests (target 1 request/sec or slower for small sites; faster only with explicit permission).
- Cache aggressively to avoid re-fetching unchanged pages.
- Avoid scraping PII without lawful basis.
- Don't bypass auth, paywalls, or rate-limiting controls.

---

## 4. Database Extraction

For internal data, usually direct database access. Choice between bulk-extract, incremental query, and log-based CDC.

### 4.1 Bulk export

`SELECT * FROM <table>` into file (CSV, Parquet, Avro). Simplest and most reliable for cold historical data. Anti-pattern for large operational tables — long-running queries pin transactions, hold read locks, can cause replica lag. Use database-native bulk-export utilities:
`mongoexport`, `pg_dump --data-only --format=custom`, `mysqldump --single-transaction`,
`bq extract`, `aws redshift unload`.

### 4.2 Incremental polling via JDBC/ODBC

Connector (Kafka Connect JDBC Source, Airbyte, Fivetran, Meltano) polls source on schedule, identifying new/changed rows by monotonically increasing column — typically `updated_at` or auto-incrementing `id`. Cheap to deploy but known gaps:

- Hard deletes invisible (row simply disappears).
- Rows updated with backdated `updated_at` missed.
- High-frequency polling stresses source database.
- Schema changes detected on next poll, can break connector.

### 4.3 Log-based change data capture (CDC)

Source database emits transaction log (MySQL binlog, PostgreSQL WAL, MongoDB oplog/change streams, SQL Server CDC tables). CDC connector reads log directly, publishes change events to downstream stream (Kafka, Kinesis, EventBridge). Debezium is dominant open-source CDC platform built on Kafka Connect.

Advantages over polling:

- Captures inserts, updates, **and deletes** with same fidelity.
- Captures every intermediate row state, not just latest.
- Minimal source load — reads log, not table.
- Low latency (sub-second feasible).
- Stable ordering per row.

Typical CDC pipeline: initial snapshot under low-impact isolation level (connector reads table once to establish baseline), then streams log from snapshot's LSN/position. Resume tokens/LSN offsets persisted so connector recovers after restart.

For MongoDB specifically, change streams provide oplog-backed CDC API; see mongodb-expert for details (resume tokens, pre/post images, split events).

---

## 5. Streaming Ingest

When acquisition needs continuous rather than batch, need streaming platform.

### 5.1 Kafka, Kinesis, Pub/Sub — selection

- **Apache Kafka** (and managed variants — MSK, Confluent Cloud, Aiven, Redpanda): de facto standard. Open protocol, runs everywhere, strongest ecosystem (Kafka Connect, Kafka Streams, ksqlDB), highest raw throughput (>1M events/sec achievable). Best for multi-cloud, hybrid, complex stream processing, or exactly-once semantics across full pipeline. Requires operational expertise unless managed.
- **Amazon Kinesis Data Streams**: AWS-native. Tightly integrated with Lambda, Firehose, Glue, EMR. Shard-based capacity model. Best AWS-only when you want IAM-based auth and CloudWatch metrics out of box. 2026 trend: new AWS deployments default to MSK over Kinesis unless pipeline is small and serverless.
- **Google Pub/Sub**: GCP-native, fully serverless, push or pull delivery. Supports exactly-once delivery within region (introduced 2024) via acknowledgement IDs and dedup windows. Best on GCP for zero operational burden.

### 5.2 Exactly-once semantics

Three delivery guarantees:

- **At-most-once**: messages may be lost, never duplicated.
- **At-least-once**: messages never lost, may be duplicated (default in most systems).
- **Exactly-once**: each message processed exactly once end-to-end. Most expensive guarantee.

Kafka supports exactly-once via idempotent producers (`enable.idempotence=true`) and transactions API (read-process-write atomicity across topics). Kinesis exactly-once via KCL checkpoint mechanism plus idempotent downstream processing. Pub/Sub exactly-once is regional, requires subscriber to use exactly-once delivery API.

Practical guidance: design downstream consumers idempotent (handle same event twice without side effects), use at-least-once as baseline, invoke exactly-once only when duplicate cost genuinely exceeds coordination cost.

### 5.3 Order, partitioning, and back-pressure

- **Partitioning** determines parallelism and ordering. Same partition key → same partition → read in order. Choose key that distributes load evenly *and* groups messages whose order matters (typically entity ID — user ID, account ID, order ID).
- **Hot partitions** = primary streaming failure mode. Single key receiving disproportionate traffic backs up one consumer while others sit idle.
- **Back-pressure** must be handled at consumer. Consumer falls behind → choose: buffering (memory pressure), dropping (data loss), or pausing source (back-pressure propagation).

---

## 6. ETL vs ELT and the Modern Data Stack

### 6.1 ETL vs ELT

- **ETL** (Extract → Transform → Load): transformations before loading into warehouse. Historically necessary — warehouse was expensive compute, couldn't handle raw data scale. Tools: Informatica, Talend, SSIS, AWS Glue.
- **ELT** (Extract → Load → Transform): transformations after loading, inside warehouse. Possible because cloud warehouses (Snowflake, BigQuery, Redshift, Databricks) made compute cheap enough to transform at query time. Default 2026 pattern.

ELT decouples ingest from transformation, keeps raw history, friendlier to multiple downstream consumers with different transformation needs.

### 6.2 The modern data stack

Typical 2026 stack:

- **Ingest / EL**: Fivetran (fully managed, 500+ connectors), Airbyte (open-source, 600+ connectors, on track for 1000 by end of 2026), Meltano (Singer-based, code-first, Git-friendly), Stitch (Talend-owned). Hand-rolled extractors when SaaS connector list doesn't cover source.
- **Transform / T**: dbt (SQL- and Python-based, dominant ELT transformation framework). dbt doesn't handle extract or load — expects data already in warehouse, needs orchestrator to trigger it.
- **Orchestration**: Airflow (Python DAGs, mature, large ecosystem), Dagster (asset-oriented, type-aware), Prefect (Python-native), Mage. Runs schedule, manages dependencies, retries on failure, alerts on SLA breach.
- **Warehouse**: Snowflake, BigQuery, Databricks, Redshift, ClickHouse (real-time analytics), Motherduck (DuckDB at warehouse scale).
- **Reverse ETL**: Hightouch, Census — push transformed warehouse data back to operational systems.
- **Catalog / governance**: DataHub, OpenMetadata, Atlan, Collibra.
- **Observability**: Monte Carlo, Bigeye, Lightup, Soda.

### 6.3 Selection guidance

- **Small team, all-SaaS**: Fivetran + dbt Cloud + Snowflake + Hightouch. Pay for time saved. Fivetran's per-MAR/connection pricing (March 2025 change) can spike costs on data-rich stacks — monitor active-connection count.
- **Mid-size team wanting open-source control**: Airbyte + dbt Core + Snowflake/BigQuery + Airflow.
- **Code-first engineering team**: Meltano + dbt Core + Airflow + Snowflake/BigQuery. CLI workflows, Git-tracked everything, CI/CD for pipeline itself.

Singer specification (Taps and Targets) underpins Meltano, Stitch, and parts of Airbyte — understanding it pays off when writing custom connectors.

---

## 7. Surveys and Primary Collection

When data doesn't exist anywhere, collect it. Surveys are workhorse — Qualtrics, SurveyMonkey, Typeform, Google Forms commercially; LimeSurvey, Formr open-source. Interviews, focus groups, observational studies, and sensor instrumentation are other primary-collection paths.

### 7.1 Sampling frame

**Sampling frame** = operational list of population members actually reachable. Rarely identical to target population — this gap creates **coverage bias**. Examples:

- Target population = "all our customers". Frame = "customers we have email address for who haven't opted out of marketing". Frame systematically excludes anyone who opted out (often privacy-conscious customers or churn risks).
- Target population = "voters in upcoming election". Frame = "registered voters with phone number on file". Excludes non-registered eligible voters and voters who don't answer unknown numbers.

Document frame explicitly at design time. If frame doesn't match population, no sample size or sampling design can fix resulting bias.

### 7.2 Response bias

Bias = systematic distortion. Unlike random error, does **not** shrink with sample size — larger biased sample is more confidently wrong. Common forms:

- **Nonresponse bias**: non-respondents differ systematically from respondents. Always track response rate and benchmark respondents against population demographics.
- **Acquiescence bias** (yea-saying): respondents agree regardless of content. Mitigate with reverse-coded items and forced-choice formats.
- **Social desirability bias**: respondents answer how they think they should rather than how they feel. Mitigate with anonymous administration, indirect questioning, Bayesian truth-serum-style designs.
- **Recall bias**: respondents misremember past events, typically toward present and emotionally salient.
- **Order effects**: earlier questions prime answers to later ones. Randomize question and option order where possible.
- **Mode effects**: phone, web, mail, in-person surveys produce systematically different answers to same question. Mode-effects meta-analyses across decades show consistent, non-zero gaps.
- **Selection bias** in frame itself (see 7.1).

### 7.3 Practical design

- Pilot every survey with at least 10-20 respondents before fielding at scale.
- Use vertical scales for mobile administration to avoid side-to-side scrolling artifacts.
- Cap survey length aggressively — completion rate falls sharply past 5-7 minutes.
- Include attention-check items in long surveys but use sparingly to avoid annoying real respondents.
- Pre-register analysis plan when survey is for publication or high-stakes decision; single best defense against post-hoc rationalization.

---

## 8. Sampling Methodology

Population = full set of units you want to draw conclusions about. Sample = subset selected from sampling frame. Design choice: probability vs non-probability sampling.

### 8.1 Probability sampling

Every unit has known, non-zero probability of selection. Only probability designs support valid statistical inference (margin of error, confidence interval, hypothesis test) in standard frequentist sense.

- **Simple random sampling (SRS)**: each unit has equal probability `1/N`. Reference design. Easy to analyze, often impractical at scale (needs enumerated frame).
- **Stratified sampling**: partition frame into mutually exclusive strata (gender, geography, account tier), sample independently within each. Guarantees representation of every stratum and reduces variance for stratum-correlated outcomes. Two allocation rules:
  - **Proportional allocation**: stratum sample size proportional to stratum population size. Simple, balanced.
  - **Neyman optimal allocation**: stratum sample size proportional to stratum size times stratum standard deviation. Minimizes overall variance for fixed total `n`. Gains over proportional large when stratum variances differ substantially; when variances equal, Neyman collapses to proportional.
- **Cluster sampling**: partition population into clusters (schools, geographic blocks, branches), randomly select clusters, then sample all units within (one-stage) or sample within (two-stage). Practical when enumerated frame of individuals is impossible but frame of clusters exists. Loses precision relative to SRS — units within cluster correlated (**design effect** captures this; effective sample size = `n / DEFF`).
- **Systematic sampling**: pick every `k`th element after random start, where `k = N/n`. Works well unless frame has periodicity matching `k`, in which case catastrophically biased.

### 8.2 Non-probability sampling

Selection probability unknown or zero for some units. Statistical inference requires assumptions that often don't hold; treat as exploratory unless you can model selection.

- **Convenience sampling**: whoever is at hand. Cheap, fast, biased toward most accessible group.
- **Quota sampling**: hit target counts for specific subgroups (e.g., 50% female, 30% under 30). Looks representative on quota dimensions but can be arbitrarily biased on every other dimension.
- **Snowball sampling**: ask respondents to refer further respondents. Useful for hidden populations (substance users, undocumented migrants, niche communities) but inherits social networks of seeds.
- **Purposive / judgment sampling**: researcher selects units considered informative. Appropriate for case studies and qualitative work; never appropriate for population estimates.

Common modern hybrid: **online panel sampling with post-stratification weighting** — recruit non-probability online panel, then weight responses to match population marginals on demographic dimensions. Weighting reduces but doesn't eliminate selection bias on dimensions correlated with outcome but not weighting variables.

<!-- cross-hub-map -->
## Cross-hub map — where every data-analytics topic lives

Family split across these hubs. Deep material **not** in this hub's Sub-skill routing table → reference file under sibling hub below — **activate that hub or `Read` its `references/<name>.md` directly**. Every former standalone skill in this family is now a reference under one of these hubs (nothing deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `da-1-foundations-theory` | Data Analysis Foundations & Theory (hub) | `references/da-1-1-definitions-scope.md`, `references/da-1-1-1-data-analysis-vs-analytics-vs-data.md`, `references/da-1-1-2-analysis-vs-synthesis.md`, `references/da-1-1-3-quantitative-vs-qualitative-analysis.md`, … |
| `da-2-data-analysis-lifecycle` | Data Analysis Lifecycle & Process (hub) | `references/da-2-1-process-frameworks.md`, `references/da-2-1-1-crisp-dm.md`, `references/da-2-1-2-kdd.md`, `references/da-2-1-3-semma.md`, … |
| `da-3-data-acquisition-sampling` | Data Acquisition, Collection & Sampling (hub) | `references/da-3-1-data-sources.md`, `references/da-3-1-1-primary-vs-secondary.md`, `references/da-3-1-2-internal-vs-external.md`, `references/da-3-1-3-structured-semi-structured-unstructured.md`, … |
| `da-analytical-methods` | Data Analytical Methods (cleaning, EDA, modeling, ML, causal, time-series) | `references/da-4-data-cleaning-preparation.md`, `references/da-5-exploratory-data-analysis.md`, `references/da-6-statistical-modeling.md`, `references/da-7-machine-learning.md`, … |
| `da-data-engineering-platform` | Data Engineering & Analytics Platform (pipelines, OLAP, modeling, governance) | `references/da-10-tools-and-languages.md`, `references/da-13-data-engineering-and-pipelines.md`, `references/da-14-streaming-analytics.md`, `references/da-18-semantic-layer-headless-bi.md`, … |
| `da-applied-and-communication` | Applied Analytics, Visualization, Communication & Ethics | `references/da-8-data-visualization.md`, `references/da-9-reporting-communication.md`, `references/da-11-ethics-and-privacy.md`, `references/da-21-product-analytics.md`, … |
| `data-analytics` | Family ROUTER — entry point for all da-* sub-hubs | (this file's parent hub) |