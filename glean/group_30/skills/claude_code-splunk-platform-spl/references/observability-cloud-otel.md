<!-- Provenance: reference under the `splunk-platform-spl` standalone skill. Created 2026-06-18 via /dr deep-research (multi-source web research, ≥3 independent sources/concept). Overview level — places Observability Cloud + OTel relative to the Splunk platform. Volatile claims stamped verified-as-of: 2026-06-18. -->

# Splunk Observability Cloud & OpenTelemetry Convergence

`verified-as-of: 2026-06-18` for product naming, the OTel migration timeline, and the Cisco/Data-Fabric direction (all VOLATILE — verify). Overview level: the goal is to place Observability Cloud and OTel **correctly relative to the Splunk platform**, not to teach APM administration. Sourced to Splunk Docs, Splunk OTel Collector GitHub, OpenTelemetry.io/CNCF, and Splunk Community.

## Contents
- The one distinction that matters (platform vs Observability Cloud)
- Components (IM, APM, RUM, Synthetics, Log Observer Connect)
- Log Observer Connect — federation, not re-ingestion
- SignalFlow vs SPL
- The Splunk Distribution of the OpenTelemetry Collector
- OpenTelemetry as the vendor-neutral standard + Splunk's bet
- Cross-cutting (data-shaping tier, Cisco/Data Fabric, HEC)
- Disconfirming findings / gotchas

## The one distinction that matters
**Splunk Observability Cloud is a different product from "the Splunk platform."** Splunk's own style guide is explicit: **"Splunk Observability"** is a category (Splunk Observability Cloud + Splunk AppDynamics) that is **"separate from the Splunk platform product offerings,"** and **"the Splunk platform"** refers specifically to Splunk Enterprise + Splunk Cloud Platform.[^5] The cleanest contrast:[^1][^2][^6]
- **Splunk Enterprise / Splunk Cloud Platform** = the engine that **collects, indexes, searches, and visualizes machine data** (primarily logs/events), index-driven, queried with **SPL.**
- **Splunk Observability Cloud** = a **metrics-first, cloud-only SaaS** with its own ingestion/analysis backends for metrics, traces, and logs, **OpenTelemetry-native**, queried with **SignalFlow.**

Lineage: Observability Cloud was **built on SignalFx**, the streaming-metrics company Splunk acquired in 2019.[^8][^9] **Don't oversimplify to "platform = logs, observability = metrics"** — the platform has a metrics store too; the real split is **separate backends, separate languages (SPL vs SignalFlow), a separate cloud-only product built on SignalFx + OTel.**[^12][^27]

## Components (high level)
Splunk's service description enumerates the suite:[^4][^1]
- **Infrastructure Monitoring (IM)** — real-time **metrics** analytics across hybrid/multicloud infra, hundreds of prepackaged integrations (the SignalFx core).[^1][^13]
- **APM (Application Performance Monitoring)** — distributed **tracing**: spans→traces, service maps, and — distinctively — **full-fidelity ("NoSample") trace ingestion** that by default collects and analyzes **100% of traces** with no head/tail sampling, to catch rare/intermittent errors (tail-based sampling is still *optionally* configurable via the OTel Collector to control cost).[^1][^14][^15][^16]
- **RUM (Real User Monitoring)** — front-end/end-user experience (web vitals, errors), connecting browser sessions to back-end services.[^1][^13]
- **Synthetic Monitoring** — proactive uptime/performance testing of APIs, endpoints, and user flows.[^1][^13]
- **Log Observer Connect** — see below.[^1][^4]

## Log Observer Connect — federation, not re-ingestion
Log Observer Connect lets Observability Cloud users **query logs that already live in Splunk Cloud Platform / Splunk Enterprise** in context with metrics and traces, via a codeless interface — **without re-ingesting, storing, or indexing** those logs in Observability Cloud. Splunk docs: "your logs **remain in your Splunk Cloud Platform instance** … Log Observer Connect **does not store or index your logs data. There is no additional charge.**"[^19][^20] Mechanically it is federated search using a **service account** and the Splunk platform **search head**, honoring the user's existing index/role permissions; one click can "Open in Splunk platform" to continue in SPL.[^19][^21] (Common error to correct: LOC is *not* "sending logs to Observability Cloud" — the HEC path does that; LOC reads logs where they live.)

## SignalFlow vs SPL — different languages for different data
**SignalFlow** is "the analytics engine **at the heart of Splunk Observability Cloud**" — a **Python-like streaming language** powering all charts and detectors over **metric time series (MTS)** and trace-derived metrics; it operates on **streams** (raw MTS → computations → new streams), with `detect()` blocks firing alerts in real time.[^22][^23][^24] **SPL** is the search-centric, event-first language of the **Splunk platform** for logs/events.[^27] These are **two different languages for two different data models in two different products.** Splunk itself reinforces this by building **two separate AI assistants** — "AI Assistant for SPL" (platform) and a separate **SignalFlow-generation** specialist (Observability Cloud).[^28] The only bridge is the Splunk Infrastructure Monitoring Add-on's `| sim flow query=...` SPL operator, which *wraps* a SignalFlow string — proving the languages are distinct, not unified.[^26]

## The Splunk Distribution of the OpenTelemetry Collector
The **Splunk Distribution of the OpenTelemetry Collector** is Splunk's **officially supported build** of the upstream OTel Collector (based on the Contrib distribution).[^29][^30] Key facts:
- It is the **unified agent** to receive, process, and export **metrics, traces, and logs** to Observability Cloud.[^29][^30]
- **Support boundary (load-bearing):** "Splunk officially supports the Splunk Distribution … Splunk only provides best-effort support for the upstream OpenTelemetry Collector." You *can* send to Observability Cloud from any OTel distro, but only Splunk's is supported.[^30]
- It can **also route to the Splunk platform** (Enterprise / Cloud Platform) via the **`splunk_hec` exporter** — the Helm chart accepts `splunkObservability` (realm + access token), `splunkPlatform` (HEC endpoint + token), or **both at once.**[^29][^31][^32]
- **Two deployment modes** (mirroring upstream OTel): **host monitoring (agent) mode** — runs on/with the app, sends directly to Observability Cloud; and **data forwarding (gateway) mode** — per cluster/datacenter/region, collects from agent Collectors and forwards on. Gateways are used for centralized processing, **tail-based sampling** (the sampler must see all spans of a trace), or when hosts can't reach the SaaS backend.[^34][^36]

## OpenTelemetry as the vendor-neutral standard + Splunk's bet
**OpenTelemetry (OTel)** is a **CNCF** observability framework for generating/collecting/exporting **traces, metrics, and logs**, born from the 2019 merger of OpenTracing + OpenCensus; its guiding principle is **"you own the data you generate — no vendor lock-in."**[^37] **OTLP (OpenTelemetry Protocol)** is the standard wire protocol (gRPC/HTTP, Protobuf), now stable for traces, metrics, and logs.[^39] **VOLATILE milestone:** on **May 21, 2026, CNCF announced OpenTelemetry's GRADUATION** (its highest maturity level) — the second-highest-velocity CNCF project after Kubernetes; OTel co-creator Morgan McLean is "senior director of product management, Splunk, a Cisco company," signaling how institutional Splunk's bet is.[^40][^42]

**Splunk's strategic bet = deprecate proprietary SignalFx agents in favor of OTel:**
- **SignalFx Smart Agent: End of Support June 30, 2023** — customers "must transition to the Collector." The Smart Agent *monitors* live on as the **`smartagent` receiver** inside the Collector for backward compatibility, with Splunk progressively porting each to **native OTel receivers.**[^44][^45][^47]
- **SignalFx tracing libraries** → replaced by **Splunk Distributions of OpenTelemetry** language libraries (default OTLP exporter).[^9][^49]
- Convergence runs both ways: Splunk now supports **native log ingestion into the Splunk *platform* via OTLP/the OTel Collector** — "one agent for logs, metrics, and traces," explicitly to adopt an open, vendor-neutral standard.[^51]

**Distinct from the Splunk platform / SPL:** Observability Cloud is a different product with a different data model than the log/event side (indexers, SPL). OTel is the *ingestion standard* that feeds Observability Cloud by default and can *also* feed the platform via HEC — but the two analytical worlds remain separate backends/languages/UIs. (See `devops-observability` for OTel SDK app instrumentation; that hub owns *producing* telemetry.)

## Cross-cutting (brief, overview)
**The modern data-shaping/routing tier** sits **between sources and indexes** to filter/transform/route data before indexing, increasingly built on **SPL2:**[^52]
- **Edge Processor** — deployed **at the edge** (near sources) with a Splunk-Cloud control plane; receives from UFs/HEC/syslog and filters/transforms/routes into Splunk indexes or **Amazon S3** using **SPL2 pipelines.**[^52]
- **Ingest Processor** — a **Splunk-managed cloud service** between the forwarding (HF) layer and the Cloud indexers; transforms/routes via **SPL2** `$pipeline` modules (`from $source … into $destination`, with `route`/`thru`/`branch`).[^54][^55]
- **Ingest Actions** — the older equivalent: `props.conf`/`transforms.conf` rules at the HF/indexer level (mask/filter/route before indexing).[^57]
- **Federated Search** (e.g., Federated Search for Amazon S3) — **search data in place** (S3, via AWS Glue) without ingesting it into Splunk indexes.[^52]

**The Cisco acquisition & Cisco Data Fabric (VOLATILE):** **Cisco completed its acquisition of Splunk on March 18, 2024** (~$28B), making Splunk the data/analytics hub unifying Cisco's observability (AppDynamics, ThousandEyes) and security (Talos). Under Cisco, Splunk Observability Cloud and **AppDynamics** coexist as the two products under "Splunk Observability," with active integration (Log Observer Connect for AppDynamics, ITSI, SSO, a unified AI Assistant). The **Cisco Data Fabric** (announced Sept 8, 2025) is "a new architecture … powered by the Splunk platform" to make machine data AI-ready by **federating/searching data where it resides** (S3, Iceberg, Delta Lake, Snowflake, Azure) rather than forcing central ingestion. Much is "available through 2026" — treat feature scope as evolving.[^58][^60][^61][^64]

**HEC (HTTP Event Collector)** is the token-based HTTP(S) ingestion API connecting many of these pieces: clients POST JSON or raw events with a token (`Authorization: Splunk <token>`); HEC indexes per the source/sourcetype/index configured on the token; default port **8088**; endpoints `/services/collector/event` and `/services/collector/raw`. HEC is the bridge the OTel Collector's `splunk_hec` exporter uses to send to the Splunk platform, and a source type Edge Processor can receive.[^67][^68][^69]

## Disconfirming findings / gotchas
1. **The platform-vs-observability confusion is real and common** — users ask "isn't the Splunk platform already doing logs/metrics/traces?" The honest answer: there is genuine *capability overlap*; the cleanest differentiator is **native OTLP ingestion + the purpose-built streaming-metrics/full-fidelity-trace model**, not "metrics vs logs."[^12][^11]
2. **SignalFlow ≠ SPL** — different languages, reinforced by Splunk shipping two separate AI assistants.[^28]
3. **OTel migration pain is documented** — the Smart Agent EoS slipped (Dec 2022 → June 2023) to give migration time; you **cannot run Smart Agent and Collector simultaneously on the same host**; the backward-compat `smartagent` receiver exists precisely because not everything ported cleanly.[^44][^45]
4. **"OTel Collector = supported" is conditional** — Splunk supports **its distribution**, not the upstream/Contrib Collector (best-effort only).[^30]
5. **Log Observer Connect is NOT log ingestion into Observability Cloud** — it is federated, in-place search; logs stay in the Splunk platform, are not indexed/stored in Observability Cloud, and incur no extra charge.[^19][^20]

## Adjacent / frontier concepts
AppDynamics vs Splunk Observability Cloud under Cisco (active convergence, not a finished merge); OTel Collector gateway vs agent mode + tail-based sampling (deeper than overview — see `devops-observability`); Splunk Machine Data Lake (MDL), Time Series Foundation Model (TSFM), Cisco AI Canvas (the AI-data substrate under Data Fabric); AI observability (LLM/agent monitoring) inside Observability Cloud; SPL2 as the convergence language; OTel Profiles (continuous profiling, a new signal type); Federated Search for Amazon S3 / Apache Iceberg (query-in-place vs ingest — the economic counter-trend to "index everything").

## References
[^1]: Splunk Observability Cloud overview — Splunk Docs — https://help.splunk.com/en/splunk-observability-cloud/get-started/splunk-observability-cloud-overview/splunk-observability-cloud-overview — components, OTel-native architecture, platform contrast.
[^2]: Splunk Observability Cloud product page — Splunk (vendor) — https://www.splunk.com/en_us/products/observability-cloud.html — OTel-native, NoSample tracing, Log Observer Connect, AppDynamics. VOLATILE.
[^4]: Splunk Observability Cloud service description — Splunk Docs — https://help.splunk.com/en/splunk-observability-cloud/get-started/service-description/splunk-observability-cloud-service-description — definitive component list (IM/APM/RUM/Synthetics/LOC), LOC included free.
[^5]: Splunk product terminology (style guide) — Splunk Docs — https://help.splunk.com/en/splunk-style-guide/splunk-terminology-and-trademarks/splunk-product-terminology — "Splunk Observability" vs "The Splunk platform" as separate categories.
[^6]: Splunk Enterprise vs Observability Cloud comparison — SourceForge — https://sourceforge.net/software/compare/Splunk-vs-Splunk-Observability-Cloud/ — index-driven platform vs real-time observability.
[^8]: Splunk Observability Cloud (formerly SignalFx) — Cortex docs — https://docs.cortex.io/ingesting-data-into-cortex/integrations/splunk-observability — "formerly known as SignalFx."
[^9]: Announcing Native OpenTelemetry Support in Splunk APM — Splunk blog — https://www.splunk.com/en_us/blog/conf-splunklive/announcing-native-opentelemetry-support-in-splunk-apm.html — OTLP proposed by Splunk; SignalFx libs deprecated; HEC into Splunk Enterprise.
[^11]: Splunk Cloud migration question — Splunk Community — https://community.splunk.com/t5/Splunk-Cloud-Platform/Splunk-Cloud-migration-question/m-p/745971 — Observability Cloud "distinct … cloud-only"; RUM/Synthetics separate.
[^12]: What can Splunk Observability do that Splunk Platform can't? — Splunk Community — https://community.splunk.com/t5/Splunk-Cloud-Platform/What-can-Splunk-Observability-do-that-Splunk-Platform-Cloud/m-p/621744 — disconfirming/confusion source; OTLP is the key difference.
[^13]: Get-started Phase 2 / set up data collection — Splunk Docs — https://help.splunk.com/en/splunk-observability-cloud/get-started/get-started-as-an-admin/get-started-guide-for-admins/phase-2-set-up-data-collection — IM/APM/RUM/Synthetics setup.
[^14]: Manage services, spans, traces in Splunk APM — Splunk Docs — https://help.splunk.com/en/splunk-observability-cloud/monitor-application-performance/manage-services-spans-and-traces-in-splunk-apm — full-fidelity tracing definition.
[^15]: APM Trace Analyzer (Splunk workshop) — Splunk — https://splunk.github.io/o11y-onboarding/en/s4r/6-apm/5-apm-trace-analyzer/index.html — "NoSample … captures every trace."
[^16]: New Relic vs Splunk vs CubeAPM (2026) — CubeAPM (blog) — https://cubeapm.com/blog/new-relic-vs-splunk-observability-cloud-vs-cubeapm/ — independent confirmation of full-fidelity/no-default-sampling, optional sampling.
[^19]: Set up Log Observer Connect for Splunk Cloud Platform — Splunk Docs — https://help.splunk.com/en/splunk-observability-cloud/manage-data/view-splunk-platform-logs/set-up-log-observer-connect-for-splunk-cloud-platform — "logs remain … does not store or index … no additional charge."
[^20]: Set up Log Observer Connect for Splunk Enterprise — Splunk Docs — https://help.splunk.com/en/splunk-observability-cloud/manage-data/view-splunk-platform-logs/set-up-log-observer-connect-for-splunk-enterprise — service account + search head; data stays in Enterprise.
[^21]: Query logs in Log Observer Connect — Splunk Docs — https://help.splunk.com/en/splunk-observability-cloud/manage-data/view-splunk-platform-logs/query-logs-in-log-observer-connect — role/index permission inheritance; SPL escape hatch.
[^22]: SignalFlow Analytics — Splunk Docs — https://help.splunk.com/en/splunk-observability-cloud/signalflow-analytics/signalflow-analytics — "analytics engine at the heart"; Python-like language + engine.
[^23]: Analyze incoming data using SignalFlow — Splunk Docs — https://help.splunk.com/en/splunk-observability-cloud/signalflow-analytics/analyze-incoming-data-using-signalflow — stream objects, MTS, detectors.
[^24]: SignalFlow: What? Why? How? — Splunk Community blog — https://community.splunk.com/t5/Community-Blog/SignalFlow-What-Why-How/ba-p/706421 — powers all charts/detectors, Python-modeled.
[^26]: flow query syntax (Infrastructure Monitoring Add-on) — Splunk Docs — https://help.splunk.com/en/splunk-it-service-intelligence/administer-it-service-intelligence — `| sim flow query=` SPL wrapper around SignalFlow; proves languages distinct.
[^27]: Comparing Splunk vs Observe at Scale — Observe (competitor blog) — https://www.observeinc.com/resources/comparing-splunk-vs-observe-for-observability-at-scale — SPL event-first vs SignalFlow for metrics; separate backends/languages/UIs.
[^28]: Building an AI Assistant in Splunk Observability Cloud — Splunk blog — https://www.splunk.com/en_us/blog/artificial-intelligence/building-an-ai-assistant-in-splunk-observability-cloud.html — separate "AI Assistant for SPL" vs SignalFlow-generation specialist.
[^29]: signalfx/splunk-otel-collector — GitHub — https://github.com/signalfx/splunk-otel-collector — distribution definition; splunk_hec exporter to Enterprise/Cloud.
[^30]: Get started with the Splunk Distribution of the OTel Collector — Splunk Docs — https://help.splunk.com/en/splunk-observability-cloud/manage-data/splunk-distribution-of-the-opentelemetry-collector/get-started-with-the-splunk-distribution-of-the-opentelemetry-collector — based on Contrib; official-support-vs-best-effort boundary.
[^31]: Splunk HEC exporter — Splunk Docs — https://help.splunk.com/en/splunk-observability-cloud/manage-data/splunk-distribution-of-the-opentelemetry-collector/components/exporters/splunk-hec-exporter — exporter sends traces/logs/metrics to HEC.
[^32]: Install the Collector for Kubernetes using Helm — Splunk Docs — https://help.splunk.com/en/splunk-observability-cloud/manage-data/splunk-distribution-of-the-opentelemetry-collector/collector-for-kubernetes/install-the-collector-for-kubernetes-using-helm — splunkPlatform and/or splunkObservability; HEC token/endpoint.
[^34]: Deployment modes — Splunk Docs — https://help.splunk.com/en/splunk-observability-cloud/manage-data/splunk-distribution-of-the-opentelemetry-collector/get-started-understand-and-use-the-collector/deployment-modes — agent vs gateway.
[^36]: Agent-to-gateway deployment pattern — OpenTelemetry.io — https://opentelemetry.io/docs/collector/deployment/gateway/ — upstream pattern; tail-based sampling needs gateway.
[^37]: What is OpenTelemetry? — OpenTelemetry.io — https://opentelemetry.io/docs/what-is-opentelemetry/ — vendor-neutral framework; no-lock-in; OpenTracing+OpenCensus merger; CNCF.
[^39]: OTLP Specification — OpenTelemetry.io — https://opentelemetry.io/docs/specs/otlp/ — protocol details; stable for traces/metrics/logs.
[^40]: CNCF Announces OpenTelemetry's Graduation (May 21, 2026) — CNCF — https://www.cncf.io/announcements/2026/05/21/cloud-native-computing-foundation-announces-opentelemetrys-graduation/ — de facto standard; 2nd velocity; co-creator at "Splunk, a Cisco company". VOLATILE.
[^42]: OpenTelemetry: 2nd-Biggest CNCF Project — The New Stack — https://thenewstack.io/opentelemetry-whats-new-with-the-second-biggest-cncf-project/ — 2nd-most-active; McLean (Splunk) co-creator.
[^44]: Migrate from Smart Agent to the Collector — Splunk Docs — https://help.splunk.com/en/splunk-observability-cloud/manage-data/splunk-distribution-of-the-opentelemetry-collector/migrate-from-the-smart-agent-to-the-collector — Smart Agent EoS; transition mandatory; smartagent receiver.
[^45]: End of Support for SignalFx Smart Agent — Splunk blog (2022) — https://www.splunk.com/en_us/blog/devops/end-of-support-for-signalfx-smart-agent-moving-to-the-opentelemetry-collector.html — EoS June 30 2023; OTel vendor-agnostic; no feature gaps.
[^47]: Smart Agent receiver — Splunk Docs — https://help.splunk.com/en/splunk-observability-cloud/manage-data/splunk-distribution-of-the-opentelemetry-collector/components/receivers/smart-agent-receiver — receiver embeds legacy monitors in Collector.
[^49]: Migrate from the SFx Tracing Library — Splunk Docs — https://help.splunk.com/en/splunk-observability-cloud/manage-data/instrument-back-end-services/migrate-from-the-sfx-tracing-library — SFx tracing libs deprecated; default OTLP exporter.
[^51]: Your Data, Your Choice: Expanding Log Ingestion with OpenTelemetry — Splunk blog — https://www.splunk.com/en_us/blog/platform/expanding-log-ingestion-options-with-opentelemetry.html — native OTLP log ingestion into Splunk platform; one agent; anti-lock-in.
[^52]: Splunk Edge Processor and Federated Search — Splunk blog (2023) — https://www.splunk.com/en_us/blog/platform/splunk-edge-processor-and-federated-search.html — Edge Processor at edge → indexes/S3; Federated Search for S3; SPL2.
[^54]: About Ingest Processor — Splunk Docs — https://help.splunk.com/en/splunk-cloud-platform/process-data-at-ingest-time/use-ingest-processors — Splunk-managed cloud service + SPL2 pipelines, pre-index.
[^55]: Ingest Processor pipeline syntax — Splunk Docs — https://help.splunk.com/en/data-management/use-ingest-processors/ingest-processor-pipeline-syntax — SPL2 `$pipeline = from $source … into $destination`.
[^57]: Ingest Processor vs Ingest Actions — Splunk Community — https://community.splunk.com/t5/Getting-Data-In/Ingest-processor-vs-Ingest-actions/m-p/706299 — Ingest Actions = props/transforms at HF/indexer; Ingest Processor = cloud UI.
[^58]: Cisco Completes Acquisition of Splunk (Mar 18, 2024) — Cisco IR — https://investor.cisco.com/news/news-details/2024/Cisco-Completes-Acquisition-of-Splunk/default.aspx — close date, ~$157/share.
[^60]: Cisco makes Splunk the center of observability — TechTarget — https://www.techtarget.com/searchnetworking/news/366588252/Cisco-makes-Splunk-the-center-of-observability — $28B; Splunk as data hub; AppDynamics teams moved in.
[^61]: Cisco & Splunk Integrated Full-Stack Observability — Cisco newsroom (Jun 2024) — https://newsroom.cisco.com/c/r/newsroom/en/us/a/y2024/m06/cisco-and-splunk-launch-integrated-full-stack-observability.html — LOC for AppDynamics; ITSI integration; SSO.
[^64]: Cisco Data Fabric (Sept 8, 2025) — Cisco IR — https://investor.cisco.com/news/news-details/2025/Cisco-Data-Fabric-Transforms-Machine-Data-into-AI-Ready-Intelligence/default.aspx — Data Fabric powered by Splunk platform; federate S3/Iceberg/Delta/Snowflake/Azure; through 2026. VOLATILE.
[^67]: Set up and use HEC in Splunk Web — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/get-started/get-data-in/10.4/get-data-with-http-event-collector/set-up-and-use-http-event-collector-in-splunk-web — token-based model, no forwarder needed, port 8088, event/raw.
[^68]: Format events for HEC — Splunk Docs — https://help.splunk.com/en/data-management/collect-http-event-data/use-hec-in-splunk-cloud-platform/format-events-for-http-event-collector — token GUID, channel header for raw.
[^69]: HEC REST API endpoints — Splunk Docs — https://help.splunk.com/en/data-management/collect-http-event-data/use-hec-in-splunk-cloud-platform/http-event-collector-rest-api-endpoints — /services/collector/event|raw|health|ack.
