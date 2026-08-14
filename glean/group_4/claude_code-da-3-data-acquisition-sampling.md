# da-3-data-acquisition-sampling

**Category:** Science, Biology & Medicine
**Platform:** Claude Code
**Original Path:** claude-code/da-3-data-acquisition-sampling

## Description
Data acquisition, collection & sampling hub (da family, stage 3). TRIGGER: data-source taxonomy (primary/secondary, internal/external, structured/semi/unstructured); collection methods (web scraping/crawling, robots.txt/legality, APIs & data feeds, REST/GraphQL/gRPC, pagination/rate-limit/auth, web/app analytics instrumentation, event tracking); data-collection overview; survey & sampling-frame design, response bias; sampling methods (simple random, stratified, cluster, systematic, convenience, quota, snowball); sample size/power; acquisition plumbing — API ingestion, DB extraction/CDC, streaming ingest (Kafka/Kinesis/Pub-Sub), ETL/ELT & modern data stack. SKIP: theory → da-1-foundations-theory; lifecycle → da-2-data-analysis-lifecycle; methods/ML → da-analytical-methods; pipelines/platform → da-data-engineering-platform; viz/comms → da-applied-and-communication.

---

# Data Acquisition and Sampling

This skill covers the **third stage** of the data analysis curriculum: getting data into a form the
analysis can operate on, and constructing a sample that supports valid inference about the target
population. It sits between the analysis lifecycle and problem framing (da-2-data-analysis-lifecycle)
and downstream cleaning, EDA, and modeling (da-analytical-methods). The two activities are
intertwined because the choice of source
constrains what sampling design is possible, and the sampling plan determines which sources are
acceptable.

A common operational mistake is to treat acquisition as a logistics problem — "just pull the data" —
and to discover only at the analysis stage that the population was wrong, the sample frame had
coverage gaps, the schema drifted mid-pull, or the file format made the planned query infeasible. The
acquisition and sampling stage owns the responsibility for catching these failures *before* they
contaminate downstream work.

---

## Sub-skill routing table

This hub absorbs 9 former standalone skills as on-demand reference files. When a task matches a row, **Read the listed `references/` file** before answering — do not rely on this table alone for depth.

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

The first decision is where the data will come from. Three orthogonal axes characterize any source.

### 1.1 Primary vs secondary

**Primary data** is collected directly to answer the current research question. The instrument was
designed for this purpose, the participants were selected for this study, and the format matches the
analysis plan. Examples: surveys you designed and fielded, interviews you conducted, sensor data your
instrumentation captured, A/B test exposure logs you wired into the product, custom user research
sessions.

**Secondary data** was collected by someone else for a different purpose and is being reused. The
instrument may not exactly fit your question, the sampling frame may differ, the coding may not be
aligned with your constructs. Examples: government statistics (Census, BLS), commercial panels
(Nielsen, Comscore), academic data archives (ICPSR), internal company datasets gathered for other
projects, third-party API exports.

Strong analyses combine both: primary data characterizes the participants of interest with the right
instrument; secondary data provides counterfactual or population-level context that primary collection
cannot afford. The trade-off is purpose-fit (primary) versus speed/scale/cost (secondary).

### 1.2 Structured vs semi-structured vs unstructured

- **Structured**: row/column tabular data with a fixed schema, typed columns, typically in relational
  databases, data warehouses, or CSV/Parquet files. Examples: transaction records, sales line items,
  test scores, sensor readings.
- **Semi-structured**: hierarchical or self-describing — has discoverable structure but does not fit
  rows/columns cleanly. Examples: JSON API responses, XML documents, log files, NoSQL documents,
  Avro/Protobuf-encoded messages.
- **Unstructured**: free-form text, images, audio, video, raw documents with no inherent schema.
  Requires either feature extraction (OCR, ASR, embeddings) or a model-based interface (LLM, CV
  model) to participate in tabular analysis.

The 2020s blurred this taxonomy because vector embeddings and LLMs reduce the cost of operating on
unstructured data — but the underlying physics still apply: the closer the source is to structured
form, the cheaper and more deterministic the analysis. Schema-on-read (data lakes) defers structure
to query time; schema-on-write (warehouses) enforces it at load time.

### 1.3 Internal vs external

Internal sources (production databases, event logs, CRM, application telemetry) are typically more
reliable, more granular, and more legally defensible to use, but reflect your own organization's
operations and may not generalize. External sources (APIs, public datasets, scraped pages, panels)
broaden the population but introduce coverage uncertainty, licensing risk, and schema drift you do
not control.

---

## 2. API Ingestion

APIs are the most common acquisition path for external structured data. The choice of API style
shapes everything downstream — pagination, error handling, schema evolution, and authentication all
inherit from this decision.

### 2.1 REST, GraphQL, and gRPC

- **REST** (HTTP + JSON, resource-oriented): broadest compatibility, mature HTTP caching, well-defined
  contract via OpenAPI. Stateless, cache-friendly, easy to debug with `curl`. Default choice for
  public APIs and most acquisition pipelines. Downside: over-fetching (you get the whole resource)
  and under-fetching (you need multiple round-trips to assemble what one screen needs).
- **GraphQL** (typed query language): the client asks for exactly the fields it needs in one request.
  Single endpoint, strongly typed schema, introspectable. Eliminates over- and under-fetching.
  Downside: harder to cache at the HTTP layer, rate limiting requires query-cost analysis (GitHub's
  GraphQL API uses point budgets, not request counts), and authorization becomes per-field.
- **gRPC** (HTTP/2 + Protocol Buffers, binary): code-generated clients from `.proto`, multiplexed
  bidirectional streaming, payloads roughly 3-10x smaller than JSON. Best for service-to-service
  internal traffic. Not directly browser-accessible without gRPC-Web. Use when throughput and latency
  matter more than human readability.

A common 2026 pattern: REST for the public API surface, GraphQL for the BFF/frontend layer, gRPC
between internal services. From an acquisition perspective, gRPC is rare for external data; you will
mostly meet REST and GraphQL.

### 2.2 Authentication

- **API key**: a shared secret in a header (`X-API-Key:`) or query string. Simplest. Cannot identify a
  user, only the calling app. Use only over TLS. Rotate on a schedule. Never embed in client-side
  code or commit to source control.
- **OAuth 2.0**: delegated authorization with access tokens (short-lived) and refresh tokens
  (long-lived). Authorization Code with PKCE is the flow for user-context apps; Client Credentials is
  the flow for service-to-service. Always validate scopes server-side, never trust the client.
- **JWT (JSON Web Token)**: a signed bearer token (HS256, RS256, ES256). Carries claims (subject,
  expiry, scopes, custom). Stateless verification — no DB lookup needed if signature checks out.
  Keep tokens short-lived (5-60 minutes), use refresh tokens for renewal, and verify the `alg`
  header to avoid `alg: none` attacks.
- **mTLS / certificate auth**: client presents a certificate signed by a trusted CA. Used for
  high-trust internal traffic and some financial / healthcare APIs.

For acquisition pipelines, treat refresh tokens as the most sensitive secret in the system. They
typically have months-long lifetimes and full account access. Store them encrypted, rotate on
suspicion, and log every refresh.

### 2.3 Pagination

- **Offset/limit** (`?offset=200&limit=50`): simple but breaks when the underlying data changes
  during paging (rows shift across page boundaries). Acceptable for read-mostly small datasets.
- **Page number** (`?page=3&per_page=50`): same trade-offs as offset/limit, more readable.
- **Cursor-based** (`?after=<opaque_token>`): the server returns a token pointing to the next slice.
  Stable under concurrent writes, supports efficient indexed lookups, and is the preferred design for
  high-volume APIs (Stripe, GitHub, Slack, Twitter/X).
- **Keyset / seek pagination** (`?after_id=12345`): cursor-based but where the cursor is a domain
  value (typically a sort key plus tiebreaker). Cheap to implement on indexed columns.

When acquiring a large dataset, persist the cursor after every page so a partial failure can resume.
Never assume the API will return the same total count on a re-pull.

### 2.4 Rate limiting

- **Fixed window**: N requests per minute, counter resets at window boundary. Easy to game by
  bursting at window boundaries.
- **Sliding window**: rolling N requests over the last 60 seconds.
- **Token bucket**: refill rate plus a bucket size, allows bursts up to the bucket capacity.
- **Cost-based** (GraphQL): each query is assigned a point cost based on its complexity; you have a
  budget per hour. Use the `Authorization` and `X-RateLimit-*` response headers to monitor.

Build acquisition clients with exponential backoff on `429 Too Many Requests` and `503 Service
Unavailable`, respecting the `Retry-After` header when present. Cap total retries and surface
sustained failures rather than retrying forever.

### 2.5 Webhooks

Webhooks are the inverse of polling — the provider POSTs to your endpoint when an event occurs.
Always verify the signature header (HMAC-SHA256 over the raw body with a shared secret) before
trusting the payload. Webhooks can be replayed by attackers if you do not verify; treat any
unsigned webhook as untrusted input.

---

## 3. Web Scraping

Web scraping is acquisition without a contract. It works until the page structure changes or the
site owner objects. Use it only when there is no API and the legal/ethical posture is sound.

### 3.1 Tooling

- **`requests` + BeautifulSoup**: minimal stack for static HTML. Fast, simple, no JS execution.
- **Scrapy**: full crawling framework with built-in concurrency, throttling, retries, item pipelines,
  and middleware. Best for multi-page crawls and structured extraction at scale.
- **Playwright / Selenium**: headless browser automation. Required when the page renders content via
  JavaScript or hides data behind authentication flows. Slower (roughly 0.5-1 page per second per
  worker) and more resource-intensive than static scraping.
- **`curl_cffi` and similar TLS-impersonation tooling**: needed when the target site fingerprints
  TLS handshakes to block bots. Treat as a signal that the site does not want to be scraped.

### 3.2 Legality

US law as of 2026 turns on the hiQ Labs v. LinkedIn line of cases: scraping publicly accessible data
generally does not violate the Computer Fraud and Abuse Act. But that is only the federal anti-hacking
question. The full risk surface includes:

- **CFAA**: bypassing authentication, ignoring access controls, or manipulating URLs to reach gated
  content can still create CFAA exposure.
- **Terms of Service**: violating ToS is rarely criminal but can trigger civil contract claims and
  account termination.
- **Copyright**: bulk reproduction of copyrighted content can be infringing even if the content was
  publicly accessible.
- **GDPR / CCPA**: scraping personal data triggers data-protection obligations regardless of whether
  the data was public. The EU position is much stricter than the US.
- **Database rights** (EU): the Database Directive grants sui generis rights to the maker of a
  database; bulk extraction can infringe even without copyright.

### 3.3 Ethics and good-faith practice

- Respect `robots.txt` — it is the public, machine-readable signal of what the site owner expects.
  Ignoring it does not by itself make scraping illegal but undermines any good-faith argument.
- Set a descriptive `User-Agent` with contact information.
- Throttle requests (target 1 request/second or slower for small sites; faster only with explicit
  permission).
- Cache aggressively to avoid re-fetching unchanged pages.
- Avoid scraping personally identifying data unless you have a lawful basis.
- Do not bypass authentication, paywalls, or rate-limiting controls.

---

## 4. Database Extraction

For internal data, you usually have direct database access. The choice is between bulk-extract,
incremental query, and log-based change data capture.

### 4.1 Bulk export

`SELECT * FROM <table>` into a file (CSV, Parquet, Avro). Simplest and most reliable for cold
historical data. Anti-pattern for large operational tables — long-running queries pin transactions,
hold read locks, and can cause replica lag. Use database-native bulk-export utilities:
`mongoexport`, `pg_dump --data-only --format=custom`, `mysqldump --single-transaction`,
`bq extract`, `aws redshift unload`.

### 4.2 Incremental polling via JDBC/ODBC

Use a connector (Kafka Connect JDBC Source, Airbyte, Fivetran, Meltano) that polls the source on a
schedule, identifying new/changed rows by a monotonically increasing column — typically
`updated_at` or an auto-incrementing `id`. Cheap to deploy but has well-known gaps:

- Hard deletes are invisible (the row simply disappears).
- Rows updated with a backdated `updated_at` are missed.
- High-frequency polling stresses the source database.
- Schema changes are detected on the next poll and can break the connector.

### 4.3 Log-based change data capture (CDC)

The source database emits a transaction log (MySQL binlog, PostgreSQL WAL, MongoDB oplog/change
streams, SQL Server CDC tables). A CDC connector reads the log directly and publishes change events
to a downstream stream (Kafka, Kinesis, EventBridge). Debezium is the dominant open-source CDC
platform built on Kafka Connect.

Advantages over polling:

- Captures inserts, updates, **and deletes** with the same fidelity.
- Captures every intermediate state of a row, not just the latest.
- Minimal load on the source — reads the log, not the table.
- Low latency (sub-second feasible).
- Stable ordering per row.

Typical CDC pipeline: initial snapshot under a low-impact isolation level (the connector reads the
table once to establish baseline state), then switches to streaming the log from the snapshot's
LSN/position. Resume tokens or LSN offsets are persisted so the connector can recover after restart.

For MongoDB specifically, change streams provide an oplog-backed CDC API; see mongodb-change-streams
for the details (resume tokens, pre/post images, split events).

---

## 5. Streaming Ingest

When acquisition needs to be continuous rather than batch, you need a streaming platform.

### 5.1 Kafka, Kinesis, Pub/Sub — selection

- **Apache Kafka** (and managed variants — MSK, Confluent Cloud, Aiven, Redpanda): the de facto
  standard. Open protocol, runs everywhere, strongest ecosystem (Kafka Connect, Kafka Streams, ksqlDB),
  highest raw throughput (>1M events/sec achievable). Best for multi-cloud, hybrid, complex stream
  processing, or when you need exactly-once semantics across the entire pipeline. Requires operational
  expertise unless managed.
- **Amazon Kinesis Data Streams**: AWS-native. Tightly integrated with Lambda, Firehose, Glue, EMR.
  Shard-based capacity model. Best when you are AWS-only and want IAM-based auth and CloudWatch
  metrics out of the box. The 2026 trend is for new AWS deployments to default to MSK over Kinesis
  Data Streams unless the pipeline is small and serverless.
- **Google Pub/Sub**: GCP-native, fully serverless, push or pull delivery. Supports exactly-once
  delivery within a region (introduced in 2024) via acknowledgement IDs and dedup windows. Best
  when you are on GCP and want zero operational burden.

### 5.2 Exactly-once semantics

Three delivery guarantees:

- **At-most-once**: messages may be lost but never duplicated.
- **At-least-once**: messages are never lost but may be duplicated (the default in most systems).
- **Exactly-once**: each message is processed exactly once end-to-end. The most expensive guarantee.

Kafka supports exactly-once via idempotent producers (`enable.idempotence=true`) and the transactions
API (read-process-write atomicity across topics). Kinesis exactly-once is provided by the Kinesis
Client Library (KCL) checkpoint mechanism plus idempotent downstream processing. Pub/Sub exactly-once
is regional and requires the subscriber to use the exactly-once delivery API.

For acquisition, the practical guidance is: design downstream consumers to be idempotent (handle the
same event twice without side effects), use at-least-once delivery as your baseline, and only invoke
exactly-once when the cost of duplicates is genuinely higher than the cost of the additional
coordination.

### 5.3 Order, partitioning, and back-pressure

- **Partitioning** determines parallelism and ordering. Messages with the same partition key go to
  the same partition and are read in order. Choose a key that distributes load evenly *and* groups
  messages whose order matters together (typically the entity ID — user ID, account ID, order ID).
- **Hot partitions** are the primary streaming failure mode. A single key receiving disproportionate
  traffic backs up one consumer while others sit idle.
- **Back-pressure** must be handled at the consumer. If the consumer falls behind, you choose between
  buffering (memory pressure), dropping (data loss), or pausing the source (back-pressure
  propagation).

---

## 6. ETL vs ELT and the Modern Data Stack

### 6.1 ETL vs ELT

- **ETL** (Extract → Transform → Load): transformations happen before loading into the warehouse,
  historically because the warehouse was expensive compute and could not handle raw data scale.
  Tools: Informatica, Talend, SSIS, AWS Glue.
- **ELT** (Extract → Load → Transform): transformations happen after loading, inside the warehouse.
  Possible because cloud warehouses (Snowflake, BigQuery, Redshift, Databricks) made compute cheap
  enough to transform at query time. Default 2026 pattern.

ELT decouples ingest from transformation, lets you keep raw history, and is friendlier to multiple
downstream consumers with different transformation needs.

### 6.2 The modern data stack

A typical 2026 stack:

- **Ingest / EL**: Fivetran (fully managed, 500+ connectors), Airbyte (open-source, 600+ connectors,
  on track for 1000 by end of 2026), Meltano (Singer-based, code-first, Git-friendly), Stitch
  (Talend-owned). Hand-rolled extractors when the SaaS connector list does not cover your source.
- **Transform / T**: dbt (SQL- and Python-based, the dominant ELT transformation framework). dbt does
  not handle extract or load — it expects data already in the warehouse and needs an orchestrator
  to trigger it.
- **Orchestration**: Airflow (Python DAGs, mature, large ecosystem), Dagster (asset-oriented,
  type-aware), Prefect (Python-native), Mage. Run the schedule, manage dependencies, retry on
  failure, alert on SLA breach.
- **Warehouse**: Snowflake, BigQuery, Databricks, Redshift, ClickHouse (real-time analytics),
  Motherduck (DuckDB at warehouse scale).
- **Reverse ETL**: Hightouch, Census — push transformed warehouse data back to operational systems.
- **Catalog / governance**: DataHub, OpenMetadata, Atlan, Collibra.
- **Observability**: Monte Carlo, Bigeye, Lightup, Soda.

### 6.3 Selection guidance

- **Small team, all-SaaS**: Fivetran + dbt Cloud + Snowflake + Hightouch. Pay for the time you save.
  Fivetran's per-MAR/connection pricing (March 2025 change) can spike costs on data-rich stacks —
  monitor the active-connection count.
- **Mid-size team wanting open-source control**: Airbyte + dbt Core + Snowflake/BigQuery + Airflow.
- **Code-first engineering team**: Meltano + dbt Core + Airflow + Snowflake/BigQuery. CLI workflows,
  Git-tracked everything, CI/CD for the pipeline itself.

The Singer specification (Taps and Targets) underpins Meltano, Stitch, and parts of Airbyte —
understanding it pays off when writing a custom connector.

---

## 7. Surveys and Primary Collection

When the data does not exist anywhere, you collect it. Surveys are the workhorse — Qualtrics,
SurveyMonkey, Typeform, Google Forms commercially; LimeSurvey, Formr open-source. Interviews, focus
groups, observational studies, and sensor instrumentation are the other primary-collection paths.

### 7.1 Sampling frame

The **sampling frame** is the operational list of population members that can actually be reached.
It is rarely identical to the target population — this is the gap that creates **coverage bias**.
Examples:

- Target population = "all our customers". Frame = "customers we have an email address for who have
  not opted out of marketing". The frame systematically excludes anyone who opted out (often
  privacy-conscious customers or churn risks).
- Target population = "voters in the upcoming election". Frame = "registered voters with a phone
  number on file". Excludes non-registered eligible voters and voters who do not answer unknown
  numbers.

Document the frame explicitly at design time. If the frame does not match the population, no sample
size or sampling design can fix the resulting bias.

### 7.2 Response bias

Bias is systematic distortion. Unlike random error, it does **not** shrink with sample size — a
larger biased sample is more confidently wrong. Common forms:

- **Nonresponse bias**: people who do not respond differ systematically from those who do. Always
  track response rate and benchmark respondents against population demographics.
- **Acquiescence bias** (yea-saying): respondents agree with statements regardless of content.
  Mitigate with reverse-coded items and forced-choice formats.
- **Social desirability bias**: respondents answer the way they think they should rather than how
  they actually feel. Mitigate with anonymous administration, indirect questioning, and Bayesian
  truth-serum-style designs.
- **Recall bias**: respondents misremember past events, typically toward the present and the
  emotionally salient.
- **Order effects**: earlier questions prime answers to later ones. Randomize question and option
  order where possible.
- **Mode effects**: phone, web, mail, and in-person surveys produce systematically different answers
  to the same question. Mode-effects meta-analyses across decades of social-survey waves show
  consistent, non-zero gaps.
- **Selection bias** in the frame itself (see 7.1).

### 7.3 Practical design

- Pilot every survey with at least 10-20 respondents before fielding at scale.
- Use vertical scales for mobile administration to avoid side-to-side scrolling artifacts.
- Cap survey length aggressively — completion rate falls sharply past 5-7 minutes.
- Include attention-check items in long surveys but use them sparingly to avoid annoying real
  respondents.
- Pre-register the analysis plan when the survey is for a publication or high-stakes decision; it
  is the single best defense against post-hoc rationalization.

---

## 8. Sampling Methodology

A population is the full set of units you want to draw conclusions about. A sample is a subset
selected from the sampling frame. The design choice is between probability and non-probability
sampling.

### 8.1 Probability sampling

Every unit has a known, non-zero probability of selection. Only probability designs support valid
statistical inference (margin of error, confidence interval, hypothesis test) in the standard
frequentist sense.

- **Simple random sampling (SRS)**: each unit has equal probability `1/N`. The reference design.
  Easy to analyze, often impractical to implement at scale (you need an enumerated frame).
- **Stratified sampling**: partition the frame into mutually exclusive strata (gender, geography,
  account tier) and sample independently within each. Guarantees representation of every stratum
  and reduces variance for stratum-correlated outcomes. Two allocation rules:
  - **Proportional allocation**: stratum sample size proportional to stratum population size.
    Simple, balanced.
  - **Neyman optimal allocation**: stratum sample size proportional to stratum size times stratum
    standard deviation. Minimizes overall variance for a fixed total `n`. Gains over proportional
    are large when stratum variances differ substantially; when variances are equal Neyman collapses
    to proportional.
- **Cluster sampling**: partition the population into clusters (schools, geographic blocks, branches),
  randomly select clusters, then sample all units within (one-stage) or sample within (two-stage).
  Practical when an enumerated frame of individuals is impossible but a frame of clusters exists.
  Loses precision relative to SRS because units within a cluster are correlated (the **design
  effect** captures this; the effective sample size is `n / DEFF`).
- **Systematic sampling**: pick every `k`th element after a random start, where `k = N/n`. Works
  well unless the frame has periodicity matching `k`, in which case it can be catastrophically
  biased.

### 8.2 Non-probability sampling

Selection probability is unknown or zero for some units. Statistical inference requires assumptions
that often do not hold; treat these designs as exploratory unless you can model the selection.

- **Convenience sampling**: whoever is at hand. Cheap, fast, biased toward whatever group is most
  accessible to the researcher.
- **Quota sampling**: hit target counts for specific subgroups (e.g., 50% female, 30% under 30). Looks
  representative on the quota dimensions but can be arbitrarily biased on every other dimension.
- **Snowball sampling**: ask respondents to refer further respondents. Useful for hidden populations
  (substance users, undocumented migrants, niche communities) but inherits the social networks of
  the seeds.
- **Purposive / judgment sampling**: the researcher selects units they consider informative.
  Appropriate for case studies and qualitative work; never appropriate for population estimates.

A common modern hybrid is **online panel sampling with post-stratification weighting**: recruit a
non-probability online panel, then weight responses to match population marginals on demographic
dimensions. The weighting reduces but does not eliminate selection bias on dimensions correlated
with the outcome but not the weighting variables.

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