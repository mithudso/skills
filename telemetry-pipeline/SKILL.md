---
name: telemetry-pipeline
description: >-
  Telemetry / observability pipelines — the vendor-neutral processing layer (telemetry router) between telemetry SOURCES and DESTINATIONS (Splunk, Elasticsearch/OpenSearch, SIEMs, S3, data lakes) that collects, parses, enriches, routes, reduces, samples, and redacts logs/metrics/traces IN FLIGHT. TRIGGER: choosing/designing an observability pipeline or telemetry router; Cribl Stream (Sources/Routes/Pipelines/Functions/Packs, Leader+Workers, replay); the OpenTelemetry Collector AS A PIPELINE (receivers/processors/exporters/connectors, OTTL, agent vs gateway, tail_sampling, routing connector); Vector (VRL), Fluentd, Fluent Bit, Logstash (grok) and which to pick; SIEM/Splunk per-GB license-cost reduction (drop/trim/dedupe/sample/logs-to-metrics); PII redaction/masking in flight; schema normalization (OTel SemConv, ECS, Splunk CIM, OCSF); fan-out routing to many destinations; the observability data lake / tiered-telemetry pattern (full fidelity to cheap object storage, reduced subset to the SIEM) and replay/rehydration; whether you need a pipeline layer at all. SKIP: the Splunk platform/SPL that INGESTS this telemetry → splunk-platform-spl; Elasticsearch/OpenSearch engine & ELK stack → elasticsearch-opensearch; app-level OTel SDK instrumentation (generating spans/metrics in code), Pino, Sentry, eBPF → devops-observability hub; DATA-pipeline / data-QUALITY observability (freshness, volume, schema, lineage of analytical data — a DIFFERENT meaning of "observability") → da-data-engineering-platform.
version: 1.0.0
updated: 2026-06-18
category: observability
whenToUse: Use when the task is about the telemetry-routing / observability-pipeline layer that sits between telemetry sources and analytics/SIEM destinations — selecting or designing a pipeline (Cribl Stream, OpenTelemetry Collector, Vector, Fluentd, Fluent Bit, Logstash), reducing SIEM/Splunk ingest cost, redacting PII in flight, normalizing/enriching/routing telemetry to multiple destinations, or building the observability-data-lake / tiered-telemetry pattern with replay. Not for the Splunk or Elasticsearch destinations themselves, app-code OTel-SDK instrumentation, or analytical data-quality observability.
keywords:
  - telemetry pipeline
  - observability pipeline
  - telemetry router
  - Cribl Stream
  - OpenTelemetry Collector
  - Vector
  - Fluent Bit
  - Fluentd
  - Logstash
  - log routing
  - SIEM cost reduction
  - PII redaction
  - observability data lake
  - log rehydration
  - OTTL
tags:
  - observability
  - devops
  - telemetry
  - logging
  - siem
  - opentelemetry
---

# Telemetry / Observability Pipelines

The vendor-neutral **processing layer between telemetry sources and destinations** — a "telemetry router" that collects, parses, enriches, routes, reduces, samples, and redacts logs/metrics/traces **in flight**, before they land in expensive analytics/SIEM backends or cheap object storage.[^datadog-tp][^chrono-guide][^solarwinds] This reference is tool-comparative and pattern-focused; it cites the destinations (`splunk-platform-spl`, `elasticsearch-opensearch`) and the instrumentation layer (the `devops-observability` hub) rather than re-covering them.

**verified-as-of: 2026-06-18.** Vendor landscape, pricing, product names, version numbers, and CNCF maturity below are volatile — re-verify before quoting figures to a customer. Vendor-published reduction percentages and cost-savings dollar figures are marketing claims unless an independent source corroborates; they are flagged inline.

## What an observability pipeline is, and why it exists

An **observability pipeline** (a.k.a. **telemetry pipeline**, **observability data router**, "first-mile" layer) is a real-time stream-processing tier inserted between telemetry **sources** (agents, OTel SDKs, syslog, cloud log feeds) and **destinations** (observability/APM platforms, SIEMs, object storage / data lakes).[^datadog-tp][^chrono-guide][^solarwinds] Instead of shipping data agent-to-tool, it handles MELT data (metrics, events, logs, traces) and **reduces, transforms, enriches, redacts, and routes it in flight**.[^datadog-tp][^chrono-guide]

Architecturally every product is the same shape: **input/collection → processing (parse, enrich, filter, redact, sample) → routing → destination/output**, plus self-monitoring and governance.[^chrono-guide][^solarwinds][^otel-arch]

**Why it emerged** (corroborated across Datadog, Chronosphere, Cribl): telemetry volume grew exponentially with cloud-native/hybrid adoption; much of it is high-volume, low-access (write-once, rarely-queried — CDN/VPC-flow/firewall logs) and noisy/redundant; per-GB-ingest pricing of destinations (especially Splunk) made that economically painful. Pipelines arose to (a) cut cost before data hits expensive tools, (b) enforce compliance/PII governance close to the source, (c) de-risk/accelerate vendor migrations, and (d) break vendor lock-in.[^datadog-tp][^chrono-guide][^sacra]

**Analyst framing.** Gartner published a **Market Guide for Telemetry Pipelines** (doc 6906466) on **2025-09-02** — a *Market Guide* (emerging market), not a Magic Quadrant.[^gartner] A widely-reprinted Gartner projection: **by 2027 ~40% of log telemetry will pass through a telemetry-pipeline product, up from <20% in 2024**.[^gartner] *(Qualified: the primary doc is paywalled; the 40%/<20% figures come through three independent reprint summaries and are consistent. Volatile — Gartner refreshes annually.)*

The category was popularized commercially by **Cribl** (~2017–2018), then joined by the CNCF **OpenTelemetry Collector**, **Vector**, **Fluentd/Fluent Bit**, and others (Edge Delta, Mezmo, Chronosphere, Datadog Observability Pipelines).[^sacra][^otel-arch]

---

## The #1 driver: SIEM / Splunk per-GB ingest cost

This is the dominant reason teams adopt a pipeline, so it leads.

**The economics.** Splunk and most SIEM/observability vendors charge by **GB/day ingested**, measured at the indexing pipeline on *raw* (uncompressed) bytes, *before* compression to disk.[^splunk-license] The cardinal rule, near-verbatim across vendors and analysts: **filter during a search and you already paid for it; filter before the indexer and you never pay at all.**[^redhound][^tekstream] Splunk's own docs confirm data dropped pre-index, summary indexes, and metric rollups do **not** count against quota; internal `_internal`/`_introspection` indexes are free.[^splunk-license]

**The waste premise.** A widely-repeated industry rule-of-thumb: **~70–80% of data sent to Splunk is never queried, never triggers a detection, and never appears in a dashboard.**[^redhound] *(Qualified: traces to vendor/practitioner write-ups, not an independent audited study.)*

**Pricing data points** *(stamp: 2025–2026; vendor/analyst-stated and volatile — get a live quote before quoting a customer):*
- **Splunk ingest:** ~**$100–180/GB/day** (Cloud vs self-managed, by volume/term), rule-of-thumb **~$150/GB/day**, **+$20–45/GB/day** for Enterprise Security/SIEM.[^expanso-splunk] (One vendor frames this as **$1,800–2,500 per GB/day annually**, falling toward ~$1,000 at very large volume — the *annual* framing of the same unit economics; don't confuse it with the daily rate.[^last9-costs])
- **Datadog Logs (official, 2026):** **$0.10 per ingested GB** for processing, with **indexing billed separately** per million log events by retention (≈$1.70–3.75 per 1M events at 15–30-day tiers). Datadog explicitly **decouples ingestion from indexing**.[^dd-pricing][^dd-flex]

**Combined real-world results** *(independent or semi-independent):*
- Cribl + Palo Alto firewall logs: **160 GB/day → 60 GB/day (62.5%)** in about an hour.[^cribl-powerhour]
- A bank, upstream filtering: **14.3 TB/day → 5.2 TB/day (64%)**, annual Splunk bill **$3.7M → $1.4M**, and searches ~**4× faster** (fewer indexed bytes = less compute).[^expanso-splunk]
- Analyst/vendor routine ranges: **40–60%** (Bindplane), **≥30%** (Chronosphere), **10–30%** (TekStream).[^bindplane-reduce][^chrono-platform][^tekstream]

A crucial corollary: **fewer indexed bytes also means faster searches and lower compute** — not just a smaller license bill.[^expanso-splunk][^tekstream]

---

## Capabilities & the tiered-telemetry pattern

See **`references/capabilities-and-tiered-pattern.md`** for tool-neutral pipeline capabilities (collect, parse, enrich, route, reduce, sample, redact) and the observability data-lake / tiered-telemetry + replay pattern.

## Agent vs gateway / aggregator topology

A general, tool-neutral pattern (clearest in OTel docs, identical in Cribl Edge+Stream, Fluent Bit, etc.):[^otel-agent2gw][^elastic-otel-arch]
- **Agent** — runs per-host (DaemonSet) or as a **sidecar**; stays **thin: receive-and-forward**, plus host-local collection (`filelog`, `hostmetrics`).[^otel-agent][^newrelic-modes]
- **Gateway / aggregator** — a standalone pool doing the **heavy centralized work**: filtering, **tail-based sampling**, transformation, **fan-out routing**, and **credential isolation**.[^otel-gateway][^elastic-otel-arch]
- **Combined (agent → gateway)** is the **de facto production standard** at scale.[^otel-agent2gw][^sematext]

**Why centralize:**
- **Credential isolation** — often the primary reason: only the gateway holds backend API keys.[^elastic-otel-arch]
- **Tail-sampling correctness** — the sampler must see **all spans of a trace**, so the agent tier must use a **trace-ID-aware load-balancing exporter** to pin a trace to one gateway instance. This forces a **two-layer topology**: stateless load-balancers → stateful tail-sampling collectors (StatefulSet + headless service).[^otel-gateway][^otel-scaling]

**Cautions:** the gateway is itself a **single point of failure** — run **≥2 instances behind a load balancer**, mind the single-writer principle for metrics, and size carefully. A gateway adds latency and resource cost. Most orgs aggregate metrics/logs locally but forward **traces** centrally for cross-service correlation (an egress-vs-correlation trade-off).[^otel-scaling][^sematext][^dd-otel-deploy]

---

## The tools and where each fits

See **`references/tools-and-where-each-fits.md`** for the per-tool breakdown (Cribl Stream, OpenTelemetry Collector, Vector, Fluentd, Fluent Bit, Logstash) — capabilities, config model, and when to pick each.

## Decision guidance: do you even need a pipeline?

**You can probably skip it when:** small scale (a few dozen servers, one monitoring tool, telemetry bill **< ~$50K/yr** → direct collection is fine); a single OTLP-ready backend needing no pre-processing (OTel's "no-collector" direct export); no compliance/residency obligations and no cost pressure. Multiple independent voices warn a full pipeline stack on day one is **overkill** — the inflection is "a second/third host, 3+ services with internal traffic, ad-hoc SSH stops scaling," not a fixed volume.[^expanso-tp][^productimpossible][^trailonix][^dd-otel-deploy]

**You need one when:** **multiple destinations** (fan-out / one config instead of N); **cost pressure** (the math works even with a single backend above ~$200K/yr); **compliance** (in-flight redaction/residency); **vendor-neutrality / backend migration** (decouple instrumentation from destination — otherwise switching = a multi-month re-instrumentation); **resilience** (buffer through a vendor outage); **distributed/edge** sources on expensive/unreliable links.[^datadog-tp][^expanso-tp]

**OTel Collector vs commercial Cribl trade-off** *(stamp 2024–2026):* OTel is YAML/code-driven, OTLP-native, "free" but with **TCO = engineering time + infra**; Cribl is GUI + code, any-to-any (500+ connectors), full ETL (regex/lookup/aggregation/masking), built-in autoscale/fleet UI, throughput-priced. **Break-even** where Cribl's TCO beats hand-maintained OTel YAML is a practitioner estimate of roughly **50–100 GB/day**; below that, careful OTel YAML wins.[^cloudrps][^edgedelta-top3] **Hybrid is common and recommended:** OTel Collectors at the edge feeding Cribl/commercial Stream as the central hub — open-source instrumentation + commercial-grade processing.[^cribl-combine][^edgedelta-top3]

---

## The boundary with the Splunk Edge / Ingest Processor

Splunk itself now ships processing that sits **exactly at this telemetry-pipeline seam** — the boundary is worth drawing explicitly because it overlaps both this skill and `splunk-platform-spl`:

- **Splunk Edge Processor** — a Splunk-managed, edge-deployed processing tier that filters, masks, transforms, and routes data **before it reaches the indexers** (or to S3), configured with **SPL2** pipelines.[^splunk-ingest-agg]
- **Splunk Ingest Processor** — a cloud-side equivalent that processes data at ingest (drop, mask, aggregate partial `stats`, route), also SPL2-driven.[^splunk-ingest-agg]
- **Ingest Actions** — a lighter-weight, built-in rules feature (filter/mask/route to S3) configurable in Splunk Web.

These are **Splunk's first-party answer to the same job** a vendor-neutral pipeline does — with the trade-off that they are **Splunk-specific** (SPL2, Splunk-managed) rather than any-to-any. A team already all-in on Splunk may use Edge/Ingest Processor instead of Cribl/Vector; a multi-destination or vendor-neutral shop reaches for the tools in this reference. The Splunk-platform mechanics of Edge/Ingest Processor (SPL2, deployment, HEC) live in `splunk-platform-spl`; the routing/reduction **patterns** they implement are the ones in this reference.

---

## Anti-patterns and pitfalls

1. **Over-reduction / dropping data you later need.** Sampling discards permanently (vs dedup, which preserves counts); aggressive sampling **drops the very events a detection rule needs → its true-positive rate falls**, and over-short retention → **audit failure** (PCI/HIPAA mandate minimums). Mitigate: keep errors/outliers unsampled, tune sampling to preserve the tail, verify legal retention before cutting, and never drop a field a downstream dashboard/alert/detection depends on.[^otel-dedup][^syshard-siem][^bindplane-reduce]
2. **Pipeline as a single point of failure (the reliability paradox).** "Your apps have HA and circuit breakers — your monitoring pipeline probably doesn't." When the pipeline/Prometheus is down, `absent()` alerts can't fire and nobody gets paged — so put **meta-monitoring (Dead-Man's-Switch / Watchdog) on an independent path**, and run gateways as **≥2 replicas behind a load balancer**.[^reliability-paradox][^airbnb-monitoring][^obs-antipatterns]
3. **No buffer between collection and storage.** Telemetry going straight from collection to storage chokes on a spike; a load balancer alone doesn't help ("distributing an overwhelming load just gives more overwhelmed instances"). Volumes **spike 10× during incidents — exactly when you need them** — so use real buffering (disk WAL / Kafka).[^signoz-hailmary][^oneuptime-shipping]
4. **Parsing brittleness.** A Grok mismatch or a new log format breaks the pipeline; route failures to a **dead-letter / failure store** rather than dropping them, and iterate against the actual failed samples (note that fixes typically apply to *new* data only — already-failed records need reprocessing).[^elastic-failure][^automq]
5. **The pipeline becomes its own ops burden.** "Half the job is keeping the monitoring alive." The OTel YAML surface area is large and the misconfig blast radius real; letting every team write its own collector config yields "as many configs as teams." Fix: **treat the collector config as a platform/product** (a templated base config, teams extend within guardrails) — the root problem is diffuse ownership. "Pipeline-as-code" (git-reviewed config) is the emerging mitigation.[^daninja][^mezmo-pac]
6. **Double-billing / paying twice.** Two flavors: the **"observability tax"** (paying once for poor source telemetry — over-instrumentation/high-cardinality — and again for backend workarounds), and the **rehydration/egress double-charge** (store once, pay again to access, plus cold-retrieval and cross-cloud egress). A commercial pipeline you pay for **and** a per-GB destination can *stack* costs unless the pipeline's reduction is real — the ROI hinges on actual reduction.[^platform-tax][^snare][^itbroker]
7. **"Log everything and figure it out later."** High noise-to-signal raises storage, **degrades query performance, and causes alert fatigue**; fix instrumentation **at the source** rather than only filtering downstream. Best practice across nearly every source: **start simple (basic filter + route), add complexity only when a specific requirement demands it.**[^platform-tax][^trailonix][^datadog-tp]

---

## Cross-references

- **`splunk-platform-spl`** — the Splunk platform & SPL: the primary **destination** this pipeline feeds, and the home of Splunk's first-party Edge/Ingest Processor mechanics (SPL2, HEC, indexers). This reference reciprocates that skill's forward-pointer.
- **`elasticsearch-opensearch`** — Elasticsearch/OpenSearch & the ELK/Elastic Stack: another **destination**, and the home of Logstash/Beats/Elastic Agent engine-and-stack depth and the 2021 licensing/fork history.
- **`devops-observability` hub** (this hub) — app-level **OTel SDK instrumentation** (generating spans/metrics/logs in code), `pino-structured-logging`, `sentry-monitoring`, `ebpf-observability`, `linux-perf-tracing`. The pipeline routes the telemetry those produce.
- **`da-19-data-observability`** and the **`da-data-engineering-platform`** hub — a **different meaning of "observability"**: the freshness/volume/schema/distribution/lineage health of *analytical* data and pipelines (Monte Carlo, OpenLineage). That is data-QUALITY observability of a data warehouse, not telemetry routing to a SIEM. Do not conflate the two.
- For Kafka as a buffering tier between pipeline stages → `devops-containers-cicd` (`references/terraform-kafka-infra.md`).

---

## References

[^datadog-tp]: Datadog — "What are Telemetry Pipelines?" https://www.datadoghq.com/knowledge-center/telemetry-pipelines/ (tier: docs/vendor)
[^chrono-guide]: Chronosphere — Complete Guide to Observability Pipelines. https://chronosphere.io/learn/the-complete-guide-to-observability-pipelines-transform-your-telemetry-strategy/ (tier: vendor)
[^solarwinds]: SolarWinds — "What Is an Observability Pipeline?" IT glossary. https://www.solarwinds.com/resources/it-glossary/observability-pipeline (tier: vendor)
[^otel-arch]: OpenTelemetry — Collector Architecture. https://opentelemetry.io/docs/collector/architecture/ (tier: docs)
[^sacra]: Sacra — Cribl revenue/valuation & category history. https://sacra.com/c/cribl/ (tier: analyst)
[^gartner]: Gartner — Market Guide for Telemetry Pipelines (doc 6906466, 2025-09-02). https://www.gartner.com/en/documents/6906466 ; reprint summaries via Bindplane https://bindplane.com/gartner and Mezmo https://www.mezmo.com/resources/gartner-market-guide-for-telemetry-pipelines (tier: analyst; primary paywalled, figures via reprints — qualified)
[^splunk-license]: Splunk — How Splunk Enterprise licensing works. https://help.splunk.com/en/splunk-enterprise/administer/admin-manual/10.4/configure-splunk-licenses/how-splunk-enterprise-licensing-works (tier: docs)
[^redhound]: Red Hound — Splunk log volume optimization. https://redhound.us/blog/splunk-log-volume-optimization (tier: practitioner)
[^tekstream]: TekStream — Splunk Cloud cost optimization strategies. https://www.tekstream.com/blog/splunk-cloud-cost-optimization-enterprise-strategies/ (tier: practitioner)
[^expanso-splunk]: Expanso — Splunk pricing guide (bank case study). https://expanso.io/blog/splunk-pricing-guide/ (tier: vendor/analyst, 2026)
[^last9-costs]: Last9 — Breaking down Splunk costs. https://last9.io/blog/breaking-down-splunk-costs/ (tier: vendor, 2025)
[^dd-pricing]: Datadog — Log Management pricing. https://www.datadoghq.com/pricing/?product=log-management (tier: docs/pricing, 2026)
[^dd-flex]: Datadog — Flex Logs. https://docs.datadoghq.com/logs/log_configuration/flex_logs/ (tier: docs)
[^cribl-powerhour]: Cribl — Reducing Splunk log data volume (Palo Alto example). https://cribl.io/blog/logstream-power-hour-reducing-splunk-log-data-volume/ (tier: vendor)
[^bindplane-reduce]: Bindplane — Reduce telemetry ingestion costs. https://docs.bindplane.com/how-to-guides/data-collection-and-processing/reduce-telemetry-ingestion-costs (tier: vendor)
[^chrono-platform]: Chronosphere — Telemetry Pipeline / Control Plane. https://chronosphere.io/platform/telemetry-pipeline/ (tier: vendor)
[^otel-dedup]: OpenTelemetry — Log deduplication processor. https://opentelemetry.io/blog/2026/log-deduplication-processor/ (tier: docs)
[^otel-sampling]: OpenTelemetry — Sampling (head vs tail). https://opentelemetry.io/docs/concepts/sampling/ (tier: docs)
[^honeycomb-sample]: Honeycomb — Sampling. https://docs.honeycomb.io/manage-data-volume/sample (tier: vendor)
[^splunk-ingest-agg]: Splunk — Aggregate event data using Ingest Processor (Edge/Ingest Processor, SPL2). https://help.splunk.com/en/data-management/transform-and-route-data/process-data-at-ingest-time/process-data-using-pipelines/aggregate-event-data-using-ingest-processor (tier: docs)
[^dd-logs-metrics]: Datadog — Generate metrics from high-volume logs. https://www.datadoghq.com/blog/observability-pipelines-generate-metrics-from-high-volume-logs/ (tier: vendor)
[^oo-logs-metrics]: OpenObserve — Logs to metrics. https://openobserve.ai/blog/logs-to-metrics/ (tier: vendor)
[^syshard-siem]: Systems Hardening — SIEM cost optimization (sampling/retention risk table). https://www.systemshardening.com/articles/observability/siem-cost-optimization/ (tier: practitioner)
[^oneuptime-mask]: OneUptime — Data masking pipeline / PII redaction. https://oneuptime.com/blog/post/2026-02-06-data-masking-pipeline-pii-redaction/view (tier: vendor)
[^netenrich-pii]: Netenrich Praxis — PII Guard processor. https://praxis.netenrich.com/docs/praxis/Processors/pii_guard (tier: docs)
[^dd-redaction]: Datadog — Observability Pipelines sensitive data redaction. https://www.datadoghq.com/blog/observability-pipelines-sensitive-data-redaction/ (tier: vendor)
[^oneuptime-redact]: OneUptime — Redaction processor for the OpenTelemetry Collector (allow/block, hashing/HMAC). https://oneuptime.com/blog/post/2026-02-06-redaction-processor-opentelemetry-collector/view (tier: vendor)
[^streamkap-mask]: Streamkap — Data masking in streaming (tokenization, SOC 2 audit trail). https://streamkap.com/resources-and-guides/data-masking-in-streaming (tier: vendor)
[^dd-sds]: Datadog — Sensitive Data Scanner. https://docs.datadoghq.com/observability_pipelines/processors/sensitive_data_scanner.md (tier: docs)
[^oo-redaction]: OpenObserve — Sensitive data redaction (ingestion- vs query-time). https://openobserve.ai/blog/sensitive-data-redaction-openobserve/ (tier: vendor)
[^otel-schemas]: OpenTelemetry — Telemetry schemas. https://opentelemetry.io/docs/specs/otel/schemas/ (tier: docs)
[^otel-schemaproc]: OpenTelemetry Collector Contrib — schema processor. https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/schemaprocessor (tier: docs)
[^elastic-ecs-otel]: Elastic — ECS and OpenTelemetry (ECS donated to OTel, Apr 2023). https://www.elastic.co/docs/reference/ecs/ecs-opentelemetry (tier: docs)
[^cribl-central]: Cribl — Building central logging architectures for scale (CIM, tiering, destination split). https://cribl.io/resources/sb/building-central-logging-architectures-for-scale/ (tier: vendor)
[^ocsf]: OCSF — Understanding OCSF. https://github.com/ocsf/ocsf-docs/blob/main/overview/understanding-ocsf.md (tier: docs/standard)
[^apto-playbook]: Apto Solutions — Telemetry Pipeline Playbook (OTel→OCSF). https://www.aptosolutions.co.uk/lp/the-telemetry-pipeline-playbook-whitepaper/ (tier: vendor)
[^signoz-geoip]: SigNoz — OTel Collector GeoIP processor. https://signoz.io/docs/opentelemetry-collection-agents/opentelemetry-collector/geoip-processor/ (tier: vendor/docs)
[^tell-routing]: Tell — Pipeline routing (fan-out independence). https://docs.tell.rs/pipeline/routing (tier: docs)
[^fluentbit-router]: Fluent Bit — Router (tag-based vs conditional routing). https://docs.fluentbit.io/manual/data-pipeline/router (tier: docs)
[^otel-routing]: OpenTelemetry Collector Contrib — routing connector. https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/connector/routingconnector/README.md (tier: docs)
[^dataprepper]: OpenSearch — Data Prepper pipelines (route). https://docs.opensearch.org/latest/data-prepper/pipelines/pipelines/ (tier: docs)
[^cribl-replay]: Cribl — Master telemetry replay with Cribl Stream and Cribl Lake. https://cribl.io/blog/master-telemetry-replay-with-cribl-stream-and-cribl-lake/ (tier: vendor)
[^itbroker]: ITBroker — Why SIEM becomes expensive (hot/cold cost ratio, egress). https://www.itbroker.com/blog/why-siem-becomes-expensive (tier: analyst, 2026)
[^solide]: Solide — Enterprise log retention strategies (Parquet savings). https://solideinfo.com/enterprise-log-retention-strategies/ (tier: practitioner)
[^dd-rehydrate]: Datadog — Rehydrate archived logs with Observability Pipelines. https://www.datadoghq.com/blog/rehydrate-archived-logs-with-observability-pipelines/ (tier: vendor)
[^cribl-replay-s3]: Cribl — Replay from S3 (path-first filtering). https://docs.cribl.io/stream/4.6/usecase-replay-s3/ (tier: docs)
[^honeycomb-s3]: Honeycomb — Collector export to S3 (field-based indexing). https://docs.honeycomb.io/send-data/telemetry-pipeline/enhance/collector-export-s3 (tier: docs)
[^edgedelta-rehydrate]: Edge Delta — What is log rehydration (Glacier latency/fees). https://edgedelta.com/company/knowledge-center/what-is-log-rehydration (tier: vendor)
[^hydrolix]: Hydrolix — "Tiered storage is costing you" (contrarian always-warm view). https://hydrolix.io/blog/tiered-storage-is-costing-you/ (tier: vendor/competitor — biased)
[^automq]: AutoMQ — Observability pipelines on Kafka (replay objectives, failure store). https://www.automq.com/blog/observability-pipelines-on-kafka-logs-metrics-traces (tier: vendor)
[^otel-agent2gw]: OpenTelemetry — Agent-to-gateway deployment. https://opentelemetry.io/docs/collector/deploy/other/agent-to-gateway/ (tier: docs)
[^elastic-otel-arch]: Elastic — OpenTelemetry Collector reference architectures (credential isolation). https://www.elastic.co/observability-labs/blog/opentelemetry-collector-reference-architectures (tier: vendor)
[^otel-agent]: OpenTelemetry — Agent deployment. https://opentelemetry.io/docs/collector/deploy/agent/ (tier: docs)
[^newrelic-modes]: New Relic — OTel Collector deployment modes in Kubernetes. https://newrelic.com/blog/infrastructure-monitoring/opentelemetry-collector-deployment-modes-in-kubernetes (tier: vendor)
[^otel-gateway]: OpenTelemetry — Gateway deployment (tail sampling, load-balancing). https://opentelemetry.io/docs/collector/deploy/gateway/ (tier: docs)
[^otel-scaling]: OpenTelemetry — Scaling the Collector (loadbalancingexporter, StatefulSet, queue metrics). https://opentelemetry.io/docs/collector/scaling/ (tier: docs)
[^sematext]: Sematext — Running OpenTelemetry at scale. https://sematext.com/blog/running-opentelemetry-at-scale-architecture-patterns-for-100s-of-services/ (tier: vendor)
[^dd-otel-deploy]: Datadog — OTel deployment patterns (no-collector, HA). https://www.datadoghq.com/blog/otel-deployments/ (tier: vendor)
[^cncf-migrate]: CNCF / The New Stack — Fluentd → Fluent Bit migration guide (2025-10-01; 10–40× efficiency claim). https://www.cncf.io/blog/2025/10/01/ (tier: standard/foundation; efficiency claim qualified)
[^victoria-bench]: VictoriaMetrics — log-collector benchmark (Mar 2026: Fluent Bit/Vector/OTel/Fluentd/Filebeat throughput+memory). https://victoriametrics.com/blog/ (tier: vendor benchmark; relative ranking corroborated, absolute figures soft)
[^otel-config]: OpenTelemetry — Collector configuration (service block, pipelines). https://opentelemetry.io/docs/collector/configuration/ (tier: docs)
[^otel-internal]: OpenTelemetry — Internal telemetry (self-observability, :8888). https://opentelemetry.io/docs/collector/internal-telemetry/ (tier: docs)
[^otel-processors]: OpenTelemetry Collector — processor README (core = batch + memory_limiter only; ordering). https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/README.md (tier: docs)
[^otel-exporters]: OpenTelemetry Collector Contrib — exporters (debug replaced logging; splunk_hec, awss3, elasticsearch). https://github.com/open-telemetry/opentelemetry-collector-contrib (tier: docs)
[^signoz-processors]: SigNoz — OTel Collector processors (recommended ordering). https://signoz.io/blog/opentelemetry-collector-processors/ (tier: vendor)
[^otel-ottl]: OpenTelemetry Collector Contrib — OTTL README (grammar, contexts, Beta). https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/ottl/README.md (tier: docs)
[^otel-transform]: OpenTelemetry — Transforming telemetry (transform processor). https://opentelemetry.io/docs/collector/transforming-telemetry/ (tier: docs)
[^otel-otlp]: OpenTelemetry — OTLP specification (vendor-neutral wire protocol). https://opentelemetry.io/docs/specs/otlp/ (tier: docs)
[^cncf-graduate]: CNCF — OpenTelemetry graduation announcement (2026-05-21). https://www.cncf.io/announcements/2026/05/21/cloud-native-computing-foundation-announces-opentelemetrys-graduation-solidifying-status-as-the-de-facto-observability-standard/ (tier: standard/foundation)
[^prnews-graduate]: PR Newswire — OpenTelemetry graduation (velocity, 12,000+ contributors). https://www.prnewswire.com/news-releases/cloud-native-computing-foundation-announces-opentelemetrys-graduation-solidifying-status-as-the-de-facto-observability-standard-302778233.html (tier: press)
[^otel-resiliency]: OpenTelemetry — Collector resiliency (in-memory queue, file_storage, data loss). https://opentelemetry.io/docs/collector/resiliency/ (tier: docs)
[^otel-mgmt]: OpenTelemetry — Collector management (OpAMP). https://opentelemetry.io/docs/collector/management/ (tier: docs)
[^otel-memlimiter]: OpenTelemetry Collector — memory_limiter processor README. https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/memorylimiterprocessor/README.md (tier: docs)
[^otel-survey]: OpenTelemetry — Collector follow-up survey analysis (YAML complexity, managed-service demand). https://opentelemetry.io/blog/2026/otel-collector-follow-up-survey-analysis/ (tier: docs/survey)
[^cribl-arch]: Cribl — Core architecture concepts (three-plane model, Leader/Workers). https://docs.cribl.io/reference-architectures/core-arch-concepts/ (tier: docs)
[^cribl-basics]: Cribl — Stream basic concepts (Sources/Routes/Pipelines/Functions/Packs). https://docs.cribl.io/stream/basic-concepts/ (tier: docs)
[^cribl-functions]: Cribl — Functions library + shared-nothing/state caveat. https://docs.cribl.io/stream/functions/ (tier: docs)
[^cribl-deploy]: Cribl — Deploy architecture / architectural considerations. https://docs.cribl.io/stream/deploy-architecture/ (tier: docs)
[^cribl-geoip]: Cribl — GeoIP Function. https://docs.cribl.io/stream/geoip-function/ (tier: docs)
[^cribl-pii]: Cribl — PII redaction use case (Mask/Eval). https://docs.cribl.io/use-cases/usecase-pii/ (tier: docs)
[^cribl-edge]: Cribl — Cribl Edge product. https://cribl.io/products/edge/ (tier: vendor)
[^cribl-lake]: Cribl — Cribl Lake product. https://cribl.io/products/lake/ (tier: vendor)
[^hydrolix-cribl]: Hydrolix — Reducing Splunk costs: Cribl vs Hydrolix (ingest-on-discarded-data critique). https://hydrolix.io/blog/reducing-splunk-costs-cribl-vs-hydrolix/ (tier: vendor/competitor — biased)
[^chrono-vs-cribl]: Chronosphere — Telemetry Pipeline vs Cribl (Node.js per-vCPU throughput claim). https://chronosphere.io/chronosphere-telemetry-pipeline-vs-cribl/ (tier: vendor/competitor — biased)
[^vector-about]: Vector — vector.dev docs (sources/transforms/sinks, VRL, buffers, acks) + Datadog acquisition (2021-02-11) + governance. https://vector.dev/docs/ (tier: docs/vendor)
[^fluentd-about]: Fluentd — docs.fluentd.org (plugin model, tag routing, buffering, fluent-package v6) + CNCF graduation (2019-04-11). https://docs.fluentd.org/ (tier: docs)
[^fluentbit-about]: Fluent Bit — docs.fluentbit.io (pipeline model, processors, OTLP/metrics/traces, v5.x, kubernetes filter) + CNCF graduated. https://docs.fluentbit.io/manual (tier: docs)
[^logstash-about]: Logstash — elastic.co Logstash docs (input/filter/output, grok/dissect, PQ/DLQ, pipeline-to-pipeline, Apache-2.0 core, Beats/Elastic Agent). https://www.elastic.co/guide/en/logstash/current/index.html (tier: docs)
[^expanso-tp]: Expanso — What is an observability pipeline (when to skip / need it). https://expanso.io/blog/what-is-observability-pipeline/ (tier: vendor/analyst)
[^productimpossible]: Product Impossible — You don't need an observability stack yet. https://productimpossible.com/articles/you-dont-need-observability-stack-yet/ (tier: practitioner)
[^trailonix]: Trailonix — Observability stack overkill. https://trailonix.com/blog/posts/observability-stack-overkill/ (tier: practitioner)
[^cloudrps]: CloudRPS — Observability pipeline: OTel Collector vs Cribl vs Vector (break-even). https://cloudrps.com/blog/observability-pipeline-otel-collector-cribl-vector/ (tier: practitioner)
[^edgedelta-top3]: Edge Delta — Cribl vs Edge Delta vs DIY OpenTelemetry. https://edgedelta.com/company/blog/top-3-telemetry-pipelines-cribl-vs-edge-delta-vs-diy-opentelemetry-choosing-the-right-approach-for-observability-and-security-data (tier: vendor — biased)
[^cribl-combine]: Cribl — Combining logging, metrics, and monitoring (hybrid OTel+Cribl). https://cribl.io/blog/combining-logging-metrics-and-monitoring/ (tier: vendor)
[^reliability-paradox]: Rafael Muller — The Observability Reliability Paradox (monitoring as SPOF). https://medium.com/@rafael_muller/the-observability-reliability-paradox-when-your-monitoring-becomes-a-single-point-of-failure-cfdccb1b35f6 (tier: practitioner)
[^airbnb-monitoring]: Airbnb Engineering — Monitoring reliably at scale (independent watchdog Prometheus). https://medium.com/airbnb-engineering/monitoring-reliably-at-scale-ca6483040930 (tier: practitioner)
[^obs-antipatterns]: Observability Antipatterns. https://observability-antipatterns.github.io/ (tier: practitioner)
[^signoz-hailmary]: SigNoz — Observability setup behind an observability tool (buffering). https://signoz.io/blog/project-hail-mary-observability-setup-behind-an-observability-tool/ (tier: vendor)
[^oneuptime-shipping]: OneUptime — Log shipping strategies (10× incident spike). https://oneuptime.com/blog/post/2026-01-30-log-shipping-strategies/view (tier: vendor)
[^elastic-failure]: Elastic — Streams failure store processing. https://www.elastic.co/observability-labs/blog/elastic-streams-failure-store-processing (tier: vendor)
[^daninja]: dev.to (itzdaninja) — Why setting up observability takes forever (config ownership). https://dev.to/itzdaninja/why-setting-up-observability-takes-forever-and-what-to-do-about-it-2o0i (tier: practitioner)
[^mezmo-pac]: Mezmo — Pipeline-as-code. https://www.mezmo.com/blog/modernize-telemetry-pipeline-management-with-mezmo-pipeline-as-code (tier: vendor)
[^platform-tax]: Platform Engineering — The observability tax (paying twice). https://platformengineering.com/contributed-content/the-observability-tax-why-platform-teams-pay-twice-and-how-to-stop-it/ (tier: practitioner)
[^snare]: Snare — Log visibility without the hidden costs (rehydration double-charge). https://www.snaresolutions.com/log-visibility-without-the-hidden-costs/ (tier: vendor — biased)
