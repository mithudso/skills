## The tools and where each fits

The market spans a spectrum from **lightweight edge collector** to **heavyweight aggregator/processor**, and on governance from **CNCF-neutral** to **single-vendor**.

| | Fluent Bit | OTel Collector | Vector | Fluentd | Logstash |
|---|---|---|---|---|---|
| Language/runtime | C | Go | Rust | Ruby + C ext | JVM (JRuby) |
| Footprint class | Lightest | Light–middle | Light–middle | Middle | Heaviest |
| Governance | CNCF **graduated** | CNCF **graduated** | Single-vendor (**Datadog**) | CNCF **graduated** | Single-vendor (**Elastic**) |
| Primary role | Edge/node collector (also aggregator) | Agent **and** gateway | Agent **and** aggregator | Aggregator / unified logging | Heavy central processor |
| Signals | Logs+metrics+traces (OTLP) | Traces+metrics+logs (profiles emerging) | Logs+metrics (traces experimental) | Logs (metrics/traces via plugins) | Logs (events) |
| Transform language | Lua + processors | **OTTL** | **VRL** | Ruby filters | **grok** / Ruby filters |

A recurring 2025–2026 theme: the market is consolidating toward **Fluent Bit at the edge + (Vector or the OTel Collector) as aggregator**, with **Fluentd and Logstash treated as incumbents being migrated away from**.[^cncf-migrate][^victoria-bench]

### OpenTelemetry Collector

The vendor-neutral CNCF pipeline. A single binary that **receives, processes, and exports** telemetry.[^otel-arch][^otel-config]

- **Pipeline model:** `receivers → processors → exporters`, typed per signal (traces/metrics/logs); every component in a pipeline must support that signal. Processors run **sequentially**; **fan-out** = listing multiple exporters (each gets a copy).[^otel-arch][^otel-config]
- **Connectors** join two pipelines (exporter of one + receiver of another) to summarize/replicate/route — e.g. `spanmetricsconnector` (traces → RED metrics), the **routing connector** (conditional fan-out), the core **forward connector** (zero-config fan-in/fan-out). **Extensions** add capability without touching data (`health_check`, `pprof`, `zpages`, `file_storage`, auth).[^otel-arch][^otel-config][^otel-routing]
- **Config:** six top-level blocks (`receivers/processors/exporters/connectors/extensions` + the `service:` block that *enables* them — a configured-but-unreferenced component is **not enabled**). Self-observability lives under `service::telemetry` (internal `otelcol_*` metrics, default Prometheus `:8888`).[^otel-config][^otel-internal]
- **Key components:** receivers `otlp` (core/stable), `filelog`, `hostmetrics`, `prometheus`, `kafka`, `fluentforward`, `splunk_hec`, `syslog`, `journald`; processors — core ships **only `batch` and `memory_limiter`** (no processors enabled by default), plus contrib `attributes`/`resource`/`resourcedetection`/`filter`/`transform`(OTTL)/`tail_sampling`/`probabilistic_sampler`/`k8sattributes`/`redaction`; exporters `otlp`/`otlphttp` (core/stable), `debug` (replaced the removed `logging` exporter), `splunk_hec`, `elasticsearch`, `kafka`, `awss3`, `prometheusremotewrite`, `clickhouse`.[^otel-processors][^otel-exporters][^signoz-processors]
- **Recommended processor order:** **`memory_limiter` first** → drop/sample early (`filter`, sampling) → enrichment that must precede batching (`k8sattributes`, `resourcedetection`) → transform → **`batch` last (before export)**. Specific trap: **don't batch before `tail_sampling`** — sample first, batch last.[^otel-processors][^signoz-processors]
- **OTTL (OpenTelemetry Transformation Language)** is the shared expression engine for the transform processor, filter processor, routing connector, and tail-sampling. Grammar = an **Editor** function (`set`, `keep_keys`, `delete_key`, `replace_pattern`, `truncate_all`) + optional `where` Condition; paths are context-scoped (`span.attributes[...]`, `resource.attributes`, `log.body`). It is **Beta** and does **not** support cross-signal interactions.[^otel-ottl][^otel-transform]

```yaml
# OTel Collector: redact a secret, drop noise, fan out to Splunk + S3 archive
processors:
  memory_limiter: { check_interval: 5s, limit_mib: 4000 }   # FIRST
  transform/redact:
    log_statements:
      - replace_pattern(log.attributes["cmd"], "password=\\S*", "password=***")
  filter/drop_health:
    logs: { log_record: ['attributes["http.target"] == "/healthz"'] }
  batch: {}                                                  # LAST
exporters:
  splunk_hec: { token: "${SPLUNK_HEC_TOKEN}", endpoint: "https://hec:8088" }
  awss3: { s3uploader: { region: us-east-1, s3_bucket: telemetry-lake } }
service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [memory_limiter, transform/redact, filter/drop_health, batch]
      exporters: [splunk_hec, awss3]    # fan-out: each exporter gets a copy
```

- **Strengths:** OTLP open-protocol **vendor-neutrality** ("instrument once, swap backends by config, not code"); **CNCF GRADUATED on 2026-05-21** (update any "incubating" language) with the second-highest project velocity after Kubernetes.[^otel-otlp][^cncf-graduate][^prnews-graduate]
- **Honest weaknesses (several self-admitted by the project):** tail-sampling **statefulness** forcing the two-tier load-balanced architecture; the **default exporter send queue is in-memory** (data lost on crash; the `file_storage` extension gives only best-effort disk durability, with a default ~5-min max-retry before the oldest data is dropped — the docs do **not** claim at-least-once delivery); **no native GUI/control plane** (the fleet standard is OpAMP; the polished control plane comes from commercial products like Bindplane); **memory/OOM fragility** (`memory_limiter` "is not a replacement for properly sizing the collector"); **contrib breaking-change churn**; and **operational/YAML complexity** that the project's own follow-up survey flags ("the YAML can get unwieldy fast," "configuration changes require full restarts," explicit demand for a managed service).[^otel-scaling][^otel-resiliency][^otel-mgmt][^otel-memlimiter][^otel-survey]

### Cribl Stream

The commercial market leader / category creator (founded 2018 by ex-Splunk engineers; Clint Sharp is credited with popularizing "observability pipeline").[^sacra][^cribl-arch]

- **Event-flow concepts:** **Sources** (receive/collect), **Destinations** (downstream receivers), **Routes** (evaluate filter expressions **in order** → pick one Pipeline + one Destination; a **`Final` toggle**, on by default, controls consumption — turning it **off clones** events for fan-out), **Pipelines** (an ordered series of Functions; events only move forward), **Functions** (pieces of **JavaScript** run per event — mask, hash, add field; can carry filters), **Packs** (portable bundles of routes/pipelines/sources/destinations from the Packs Dispensary), **QuickConnect** (drag-and-drop Source↔Destination that bypasses Routes).[^cribl-basics][^cribl-functions]
- **Topology:** **Leader** (control plane; **formerly "Master"** — older material still says Master) pushes config to **Workers** (data plane; each spreads work across **Worker Processes ≈ one per CPU**; designed for several TB/day each). **Worker Groups** share config. **Shared-nothing caveat:** state is **not shared across processes**, so stateful Functions (Aggregations, Sampling, Suppress, Rollup Metrics) operate **per Worker Process** unless you add an external tier like **Redis**.[^cribl-arch][^cribl-functions][^cribl-deploy] Deploys single-instance, distributed (Leader + Workers, with a load balancer in front of push Sources), or as **Cribl.Cloud (SaaS)**.[^cribl-deploy]
- **Capabilities:** routing/fan-out (Routes + `Final=off` cloning); volume reduction (**Drop**, **Sampling/Dynamic Sampling**, **Suppress** dedupe, **Aggregations/Rollup Metrics** incl. logs→metrics); enrichment (**Lookup**, **Redis**, **GeoIP**, **DNS Lookup**); **PII** (**Mask** Function regex match→replace, **Eval** for hashing/pseudonymization, field-level encryption; newer **Cribl Guard** AI-based PII detection requires a separate license); normalization to **OTLP** and security **OCSF**; **replay** (store full-fidelity raw in S3/Blob/GCS/**Cribl Lake**, then replay selected data back through a Pipeline on demand).[^cribl-functions][^cribl-geoip][^cribl-pii][^cribl-replay] Adjacent products: **Cribl Edge** (vendor-neutral endpoint agent managing 100,000s of agents via Fleets — a proprietary agent, noted as a lock-in critique), **Cribl Lake** (turnkey long-term lake), **Cribl Search** (federated search-in-place using KQL, no rehydration).[^cribl-edge][^cribl-lake]
- **Cost-control value prop:** "collect everything, filter/sample/aggregate/trim, forward only high-value data to the expensive SIEM" while routing full-fidelity copies to cheap storage; Cribl helps cut Splunk-type spend by a vendor-cited **30–90%**.[^cribl-powerhour][^sacra]
- **Disconfirming:** Cribl charges a **flat ingest fee per GB even on data it discards** (a competitor cites ~$0.32/GB Enterprise), and routing to Splunk still means **paying ingest twice** (into Cribl, then into Splunk) — so the ROI hinges on reduction being real.[^hydrolix-cribl] A competitor (Chronosphere, Fluent-Bit-based) claims Cribl's **Node.js** engine needs ~**200 GB/day per vCPU** vs their C engine ~2 TB/day, implying a large infra-cost gap.[^chrono-vs-cribl] *(Both figures are competitor-sourced and biased, but directionally consistent with Cribl's own sizing calculator that both cite.)*

### Vector

High-performance, vendor-neutral unified logs+metrics pipeline in **Rust** (single static binary). **Created by Timber; acquired by Datadog, announced 2021-02-11** (not 2022); MPL-2.0; powers Datadog Observability Pipelines; works fully without a Datadog account; **not a CNCF project**.[^vector-about]

- **Pipeline model:** a DAG of **sources → transforms → sinks**; backpressure propagates upstream.[^vector-about]
- **VRL (Vector Remap Language):** an expression DSL **compiled to bytecode at boot** with compile-time checks (no dead code, unhandled errors, or type mismatches). Used for parsing (`parse_json`, `parse_syslog`), filtering/routing, enrichment, and **PII redaction** (`redact()` with built-in SSN/credit-card/email filters; `del()`). VRL is the safe/sandboxed alternative to the in-process **Lua** transform (no OS/FS/network access).[^vector-about]
- **Roles:** thin per-node **agent** + dedicated **aggregator** (heavy processing, disk buffers); a common pattern is Fluent Bit DaemonSet → Vector aggregator (2–3 replicas).[^vector-about]
- **Reliability:** end-to-end acknowledgements on most sources/sinks (worst-status-wins on fan-out); built-in memory and disk buffers (markets these as removing the need for an external Kafka buffer) — but **disk buffers alone don't guarantee end-to-end delivery**, and not all sources can ack.[^vector-about]
- **Reduction transforms:** `sample`, `filter`, `dedupe`, `reduce`, `throttle`, `log_to_metric`; a cited field example cut a Loki cluster 500 GB/day → 50 GB/day by stripping high-cardinality request IDs.[^vector-about]
- **Disconfirming:** VRL has a **learning curve** and **drops events silently on runtime transform errors unless you explicitly route errors**; Vector is **heavier than Fluent Bit** on memory (a user reported 150 MiB vs Fluent Bit 25 MiB at the same load); **traces are experimental**; smaller plugin ecosystem than Fluentd/Logstash; **single-vendor (Datadog) governance** despite marketing vendor-neutrality.[^vector-about][^victoria-bench] The **OpenObserve fork** is low-activity; the meaningful relationship is OpenObserve consuming **VRL as a library**.[^vector-about]

### Fluentd

The CNCF **graduated** (2019-04-11) "unified logging layer" in **Ruby + C extensions**; Apache-2.0; the **F in the EFK stack** (Elasticsearch-Fluentd-Kibana).[^fluentd-about]

- **Plugin model:** input → filter → buffer → output (nine plugin types). Every event is **tag + time(ns) + record(JSON)**.[^fluentd-about]
- **Tag-based routing:** inputs assign tags; `<match>` routes by tag with wildcards (**first match wins** — broad patterns last); tags are mutable mid-pipeline (`rewrite_tag_filter`); `<label>` sections escape strict ordering.[^fluentd-about]
- **Buffering:** chunks (stage → queue → flush); **memory** (fast, lost on restart) vs **file** (durable; recommended for production); **at-least-once delivery** (configuration-dependent — memory buffers/`drop_oldest_chunk` can lose data); exponential-backoff retries; back-pressure via `overflow_action`.[^fluentd-about]
- **Reduction/enrichment/PII:** `grep`, `record_transformer`, `parser`, `geoip`; PII via `record_transformer`/`grep` masking or `fluent-plugin-sanitizer`; fan-out via `out_copy`.[^fluentd-about]
- **Packaging note:** the stable distribution was **td-agent**, renamed **`fluent-package`** in 2023; **td-agent reached EOL Dec 2023**; current line is **fluent-package v6 (LTS)** (Fluentd 1.19.x, Ruby 3.4).[^fluentd-about]
- **Disconfirming:** **Ruby GIL** bounds CPU-intensive work (needs multi-worker to use multiple cores); **~40–60MB+** footprint; **lowest throughput of the modern collectors** (a 2026 benchmark: ~5,100 logs/sec, starts losing logs before 10k); **declining engagement** (CNCF DevStats shows falling contributors/stars; "might correlate with transitioning towards maintenance mode"). Actively maintained but in a **stable-maintenance/declining** phase as cloud vendors favor Fluent Bit and Vector.[^fluentd-about][^victoria-bench][^cncf-migrate]

### Fluent Bit

CNCF **graduated**, a sub-project under the Fluentd umbrella; **C**, ~450KB baseline (tens of MB under load); Apache-2.0; the default lightweight collector for **containers/edge/Kubernetes/embedded**.[^fluentbit-about]

- **Pipeline model:** Input → (Parser) → (Processor) → Filter → Buffer → Routing → Output(s). Filters (`grep`, `modify`, `nest`, `kubernetes`, `lua`) run in the main event loop with tag/match routing. **Processors** (v2.12+/v3+) attach **directly to a specific input or output** (not global, no tag matching) and, with `threaded: true`, run in the plugin's thread — better performance than chained filters. Buffering: `memory`, `memrb` (ring), `filesystem` (hybrid mmap, backpressure control + data safety).[^fluentbit-about]
- **Now logs + metrics + traces:** native **OTLP** in/out; **Prometheus** scrape + remote write; **eBPF** input + profiling; **trace sampling** (head + tail) processor. **Version note:** current major line is **v5.x** (v5.0.0 = 2026-03-24), not v4 — v4 was the April-2025 10th-anniversary release.[^fluentbit-about]
- **Reduction/enrichment/PII:** `grep`/`modify`/`record_modifier`/`throttle` + trace sampling; the **`kubernetes` filter** for metadata enrichment (validated to 30,000 pods); `lua` + content-modifier processor (`hash` SHA-256) for transforms; **PII has no single built-in `redact` primitive** — assembled from `nightfall`/content-modifier `hash`/Lua masking, and it forwards PII unscrubbed unless explicitly configured.[^fluentbit-about]
- **Relationship to Fluentd:** both CNCF graduated, both Apache-2.0, sharing the **Forward protocol**; classic split is Fluent Bit = lightweight forwarder, Fluentd = aggregator with the bigger plugin ecosystem — but both can be forwarder OR aggregator (Fluent Bit ships an Aggregator Helm chart). CNCF published a **Fluentd → Fluent Bit migration guide (2025-10-01)** claiming Fluent Bit processes **10–40× more log volume per unit of resource**.[^cncf-migrate][^fluentbit-about] Adoption: **15+ billion deployments**; the OpenAI KubeCon-NA-2025 case study runs Fluent Bit at 9+ PB/day (a one-line `inotify: false` change cut CPU 50% fleet-wide).[^fluentbit-about]
- **Disconfirming:** **weaker transform language than Vector's VRL** (Lua + chained filters); **multiline parsing is a known pain point** (CRI multiline can OOM-kill in K8s); **heavy/stateful processing is not its strength** (use a heavier tier for complex transforms); memory-only buffering pauses inputs on `mem_buf_limit` (data loss) unless filesystem buffering is on; debugging rated "Hard."[^fluentbit-about][^victoria-bench]

### Logstash

The heavy central processor of the Elastic Stack; **JRuby on the JVM**; started 2009 (Jordan Sissel), acquired by Elastic 2013 (coined "ELK"). **Logstash core stayed Apache-2.0** — the 2021 SSPL/Elastic-License relicensing (which triggered the OpenSearch fork) hit **Elasticsearch and Kibana, not Logstash** (see `elasticsearch-opensearch`).[^logstash-about]

- **Pipeline model:** input → filter → output with **codecs**; each input runs in its own thread writing to a central queue; **pipeline worker threads** (default = CPU cores) pull batches through filters then outputs. No order guarantee unless single worker + `pipeline.ordered`.[^logstash-about]
- **grok** parses unstructured text into fields via regex (Oniguruma; ~120 built-in patterns). **`dissect`** is the faster, no-regex delimiter parser. grok's documented pitfall: **failed matches can be ~6× slower** than successful ones (explains "grok maxing the CPU"); mitigate with anchors and a `timeout`.[^logstash-about]
- **Footprint:** Elastic recommends **heap 4–8GB** for production ingestion; real-world minimal RSS ~390–500MB; PQ page size default 64MB → each PQ needs ≥128MB heap. **Heaviest of the four.**[^logstash-about]
- **Reliability:** in-memory bounded queues by default (lost on unsafe termination); **persistent queues** (PQ) give at-least-once but only protect inputs with a request-response protocol (**beats, http** — not tcp/udp/syslog); **DLQ** (Elasticsearch-output in practice) for unprocessable events; by default Logstash is **blocked when any single output is down** (use the output-isolator pipeline-to-pipeline pattern).[^logstash-about]
- **Enrichment is its strength:** `geoip`, `dns`, `useragent`, `translate` (dictionary), `jdbc_streaming`/`jdbc_static` (DB enrichment with caching), the `elasticsearch` filter (correlate prior events); PII via `mutate`/`gsub`/`fingerprint`/`anonymize`. Routing via conditionals + **pipeline-to-pipeline** (distributor / output-isolator / forked-path / collector patterns).[^logstash-about]
- **Relationship to Beats/Elastic Agent:** **Beats** = lightweight Go shippers; **Elastic Agent/Fleet** = the unified next-gen agent. Canonical pattern: lightweight Beats/Agent at the edge → Logstash heavy parsing → Elasticsearch → Kibana.[^logstash-about]
- **Disconfirming:** **heaviest footprint** (500MB–2GB RAM, "50× Fluent Bit"); JVM GC can pin a core; teams put **Beats + Kafka/Redis in front** to absorb spikes; Elastic itself is making Logstash **optional per-integration** (Elastic Agent 8.16) where Elastic Agent + ingest pipelines suffice — but it is **not phased out** (retained for DB enrichment, event-splitting, external lookups).[^logstash-about][^victoria-bench]

### Selection summary

Footprint, lightest → heaviest: **Fluent Bit < Vector ≈ OTel Collector < Fluentd ≪ Logstash** (robust across independent benchmarks; absolute numbers vary by config).[^victoria-bench] A 2026 corrective to Vector's "10× faster" marketing: that holds **vs Fluentd/Logstash**, but **vs Fluent Bit it's mixed** — Fluent Bit often wins raw node-level throughput/memory, while **Vector wins under heavy transformation**.[^victoria-bench]

- **Fluent Bit** — default for **Kubernetes/edge/IoT/embedded** and resource-constrained per-node collection; weak at heavy stateful transforms.
- **OTel Collector** — when **vendor-neutrality / open-protocol** matters most and you can own the YAML; the de-facto neutral standard, now CNCF-graduated.
- **Vector** — when you need **complex transformation/routing/cardinality control** (VRL) as a mid-tier aggregator and aren't wedded to CNCF neutrality.
- **Fluentd** — incumbent **EFK / unified logging layer** with a huge plugin ecosystem; increasingly a migrate-away target.
- **Logstash** — **Elasticsearch/ELK-first** shops needing mature **grok** parsing and the richest filter/enrichment set; heaviest footprint.
- **Cribl Stream** — commercial-grade **any-to-any** routing, a graphical UI, fleet management, and replay, when TCO beats hand-maintained open-source YAML (practitioner break-even ~50–100 GB/day).
- **Common modern architecture:** Fluent Bit DaemonSet (collect/forward) → Vector or OTel Collector aggregator (transform/route) → backend (Elasticsearch/Loki/ClickHouse/SIEM/S3).[^cncf-migrate][^victoria-bench]

---

