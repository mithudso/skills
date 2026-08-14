# The ELK / Elastic Stack — Components and Data Flow

> Provenance: reference under the `elasticsearch-opensearch` skill. Researched 2026-06-18 via `/dr` deep research (Elastic Stack docs, Logstash docs, Fleet/Elastic Agent docs, OpenSearch Data Prepper docs). **verified-as-of 2026-06-18.** Naming precision matters: Elastic's components vs OpenSearch's fork-equivalents.

## Contents

- The classic data flow
- Beats (lightweight shippers)
- Elastic Agent + Fleet (the unified agent)
- Logstash (the heavy ingest/transform engine)
- Elasticsearch ingest pipelines (the lighter alternative)
- Kibana
- OpenSearch equivalents (Dashboards, Data Prepper)
- Boundary: telemetry routers are out of scope

## The classic data flow

**source → Beats / Elastic Agent → (Logstash *or* Elasticsearch ingest pipeline) → Elasticsearch → Kibana.**[^1][^2][^3]

Elastic frames it: "Collect and ship… with Elastic Agent or Beats. Manage your Elastic Agents with Fleet… If you want to transform or enrich data before it's stored, you can use Elasticsearch ingest pipelines or Logstash."[^3] Three documented log-architecture patterns:[^4]
1. **Direct Ingest** — Beats → Elasticsearch, processing via ingest pipelines (simpler environments, basic processing).
2. **Logstash Pipeline** — adds Logstash for advanced transformation.
3. **Buffered Pipeline** — adds a broker/queue (e.g., Kafka) in front for spikes.

## Beats (lightweight, single-purpose shippers)

Beats are lightweight agents installed on servers, each collecting **one** data type:[^1][^4]
- **Filebeat** (logs/files), **Metricbeat** (metrics), **Packetbeat** (network), **Winlogbeat** (Windows events), **Heartbeat** (uptime), **Auditbeat** (audit data), **Functionbeat** (serverless).

Some Beats (Filebeat, Metricbeat) ship **modules** that install default configs, Elasticsearch ingest-pipeline definitions, and Kibana dashboards out of the box.[^2]

## Elastic Agent + Fleet (the newer unified agent)

**Elastic Agent** is a single-binary agent that replaces running multiple Beats — one agent collects logs, metrics, and more, and can also do endpoint security / osquery.[^5][^6] Key points:
- **"The Elastic Agent runs Beats under the covers"** — Agent is a management/deployment layer on top of Beats, not a from-scratch rewrite.[^7]
- **Fleet** (a Kibana app) centrally manages Agent policies, config, and binary upgrades — Beats have no native central management (Beats Central Management was deprecated).[^5][^6]
- Agent config is driven by **integrations** added via UI/API, vs Beats' per-file YAML.[^5]
- **Fleet Server** is a separate process — the communication channel between Fleet and enrolled Agents — scaling to thousands of hosts.[^6]
- Agent data lands in **data streams**, giving finer-grained volume visibility and lifecycle/permission control than default Beats indices.[^6]

## Logstash (the heavy ingest/transform engine)

Logstash's pipeline has three stages: **inputs → filters → outputs**, with optional **codecs** on inputs/outputs.[^8]
- **Inputs:** file, syslog, beats, redis, kafka, http, etc.
- **Filters:** **grok** (parses unstructured text into structured/queryable fields via pattern matching) is the workhorse; also `mutate`, `date`, `geoip`, `dissect`.[^8][^9]
- **Outputs:** elasticsearch, file, graphite, statsd, etc.

When to choose Logstash over Agent/ingest pipelines: if neither Elastic Agent nor Beats supports your source; for persistence during ingestion spikes; or to fan out to multiple destinations.[^2] It can also act as a **proxy/aggregator** — receiving from many Beats/Agents across networks and forwarding through a single firewall rule.[^2]

**Resiliency:** Logstash defaults to **in-memory queues** but offers opt-in **Persistent Queues (PQ)** — an on-disk buffer between input and filter stages (inputs with acknowledgement like beats/http are well-protected) — and **Dead Letter Queues (DLQ)** — on-disk storage for events that can't be processed (e.g., Elasticsearch mapping errors), reprocessable via the `dead_letter_queue` input plugin. Both are **disabled by default**.[^10][^11]

> Scope: Logstash is covered here only as an ELK stack component. Logstash **beyond a mention** — and other telemetry routers (Cribl, Vector, Fluent Bit, OTel Collector) — belong to the sibling `telemetry-pipeline` skill.

## Elasticsearch ingest pipelines (the lighter alternative)

The **lighter** alternative to Logstash, running *inside* Elasticsearch on ingest nodes. You configure one or more **processor** tasks that run sequentially to transform documents before indexing (e.g., `grok`, `set`, `rename`, `geoip`, `date`).[^1][^3] Reference-architecture framing: ingest-node processing suits "simpler environments / basic processing"; Logstash is for "more complex processing requirements."[^4]

## Kibana

The visualization/management UI: dashboards, visualizations, Discover, and the home for Elastic's **Observability** and **Security (SIEM)** solutions; also where you manage Stack components (ILM policies, Fleet, integrations).[^1][^3] Kibana uses **KQL** (Kibana Query Language) for its search bar, distinct from the Query DSL.

## OpenSearch equivalents (precise naming matters)

- **Kibana → OpenSearch Dashboards** — the fork's visualization layer.[^12]
- **Logstash / ingest pipelines → OpenSearch Data Prepper** — "a server-side data collector capable of filtering, enriching, transforming, normalizing, and aggregating data," and "the preferred data ingestion tool for OpenSearch."[^13] Its pipeline model is **source → buffer → processor → sink** (source and sink required; default `bounded_blocking` buffer and a no-op processor if omitted); pipelines form a DAG and support DLQ pipelines and OTLP/trace ingestion.[^13][^14] Note: Data Prepper is the Logstash-analog but is **heavily positioned around OpenTelemetry** trace/log/metric ingestion.[^14]
- **Beats → no first-party OpenSearch equivalent.** OpenSearch leans on OpenTelemetry collectors, Fluent Bit/Fluentd, and Data Prepper for collection; some users still run OSS Beats pointed at OpenSearch, but compatibility is not guaranteed across versions (verify before relying on it). *(qualify)*

## Boundary note (telemetry routers)

Vector, Fluent Bit, Cribl, and the OpenTelemetry Collector exist as *feeders* into both stacks but are owned by the sibling **`telemetry-pipeline`** skill and are deliberately not covered here. Reference them; don't duplicate.

## References

1. Elastic — The Elastic Stack: https://www.elastic.co/docs/get-started/the-stack — components + data-flow overview. docs
2. Elastic — Adding data to Elasticsearch (Cloud Enterprise): https://www.elastic.co/guide/en/cloud-enterprise/current/ece-add-data.html — Beats/Agent/Logstash choice, Logstash as proxy, PQ. docs
3. Elastic — Overview of the Elastic Stack (8.19 stack-components): https://www.elastic.co/guide/en/elastic-stack/8.19/overview.html — component definitions, ingest pipelines vs Logstash. docs
4. Elastic — Reference architectures: Logging and Log Analytics: https://www.elastic.co/docs/reference/architectures — 3 architecture patterns, ingest-node vs Logstash. docs
5. Elastic — Beats and Elastic Agent capabilities (Fleet 8.19): https://www.elastic.co/guide/en/fleet/8.19/beats-agent-comparison.html — Agent vs Beats comparison, Fleet management. docs
6. Elastic — Fleet and Elastic Agent overview (8.19): https://www.elastic.co/guide/en/fleet/8.19/fleet-overview.html — Agent unified, Fleet Server, integrations, data streams. docs
7. Elastic Blog — Easier data onboarding with Elastic Agent: https://www.elastic.co/blog/elastic-agent-and-ingest-manager — "Agent runs Beats under the covers." blog
8. Elastic — How Logstash Works (8.19): https://www.elastic.co/guide/en/logstash/8.19/pipeline.html — inputs→filters→outputs, codecs. docs
9. logstash-plugins/logstash-filter-grok: https://github.com/logstash-plugins/logstash-filter-grok — grok match/patterns. docs (repo)
10. Elastic — Queues and data resiliency (Logstash 8.19): https://www.elastic.co/guide/en/logstash/8.19/persistent-queues.html — PQ + DLQ definitions, disabled by default. docs
11. Elastic — Dead letter queues (Logstash): https://www.elastic.co/guide/en/logstash/8.19/dead-letter-queues.html — mapping-error reprocessing via dead_letter_queue input. docs
12. OpenSearch — Dashboards: https://docs.opensearch.org/latest/dashboards/ — visualization layer (Kibana fork). docs
13. OpenSearch — Data Prepper: https://docs.opensearch.org/latest/data-prepper/ — server-side collector, source/buffer/processor/sink, preferred ingestion tool. docs
14. OpenSearch — Data Prepper pipelines: https://docs.opensearch.org/latest/data-prepper/pipelines/pipelines/ — DAG, DLQ pipeline, OTLP ingestion. docs
