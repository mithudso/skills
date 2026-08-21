## Capabilities (tool-neutral)

### 1. Volume reduction

| Technique | What it does | Cited reduction | Loss profile |
|---|---|---|---|
| **Drop useless events** (debug/health-check/low-severity) | Never pay license for noise | 20–50% on chatty streams; severity filter often 30–40%[^redhound][^bindplane-reduce] | Permanent |
| **Trim/drop fields** (high-cardinality req IDs, tokens, UIDs) | Shrinks every event's byte size | 20–40% avg log-size reduction[^bindplane-reduce] | Per-field |
| **Deduplication** (e.g. OTel `logdedupprocessor`) | Collapse identical logs in a window into one event + `log_count` | "dramatically reduce storage"[^otel-dedup] | **Preserves frequency** (unlike sampling) |
| **Sampling — head vs tail** | Head: decide at creation (cheap, can't keep "all errors"). Tail: decide after seeing the whole trace (keep 100% of errors, sample the baseline) | 1% or lower common on high-volume traces[^otel-sampling][^honeycomb-sample] | **Permanently discards** |
| **Aggregation / in-stream stats** | Roll many events into batched aggregates | Splunk Ingest Processor batches partial `stats`[^splunk-ingest-agg] | Lossy by design |
| **Logs→metrics conversion** | Convert repetitive high-volume logs (CDN, VPC-flow, firewall) into count/gauge/distribution metrics, then drop/cold-store the logs | Infra telemetry is commonly 20–30% of license volume; Splunk caps each metric event at **150 bytes**[^redhound][^splunk-license][^dd-logs-metrics][^oo-logs-metrics] | Loses per-event detail |

**Key distinction:** deduplication preserves counts; **sampling discards permanently**. Choose dedup/aggregation when you need frequency; reserve aggressive sampling for data whose individual events you will never need.[^otel-dedup][^syshard-siem]

### 2. PII / sensitive-data redaction in flight

**Why do it in the pipeline:** it is a single, auditable enforcement point so that no matter which service emits telemetry, PII is caught **before data lands in a downstream store or leaves a data-residency boundary**.[^oneuptime-mask][^netenrich-pii][^dd-redaction]

**Techniques:**
- **Regex value masking** — full or **partial** (preserve last-4 of a card / TLD of an email for support workflows).[^dd-redaction][^oneuptime-redact]
- **Field deletion / allow-listing** — drop attributes not on an allow list (OTel `redaction` processor removes non-allowed keys).[^oneuptime-redact]
- **Hashing (SHA-256/SHA-3/HMAC)** — pseudonymization that **preserves correlation**: same input → same hash, so you can still join on a user without storing identity; HMAC is recommended for low-entropy values like IPs. Hashing/redaction are **irreversible**.[^oneuptime-redact][^streamkap-mask][^dd-sds]
- **Tokenization / field-level encryption** — reversible masking via a vault lookup.[^streamkap-mask]

**Ingestion-time vs query-time:** ingestion-time redaction is irreversible and best for data minimization (the original never enters storage); query-time masking stores the original and masks dynamically. Many orgs use both.[^oo-redaction] Compliance drivers named consistently: **GDPR, PCI-DSS, HIPAA, CCPA, SOC 2** — and SOC 2 expects a demonstrable audit trail of masking (which fields, which technique, when), which is why pipelines emit redaction-count audit attributes.[^oneuptime-redact][^dd-redaction][^streamkap-mask]

### 3. Parsing, normalization, enrichment

**Parsing** turns raw unstructured log lines into structured events (Grok/regex, JSON parse, event breakers); brittleness is a known failure mode (see anti-patterns).[^automq]

**Normalization to a common schema** — the schemas that matter:
- **OpenTelemetry Semantic Conventions (SemConv)** — standard names for attributes/metrics/log fields; OTel even ships telemetry schema files + a schema processor to translate data *between SemConv versions* in flight.[^otel-schemas][^otel-schemaproc]
- **Elastic Common Schema (ECS)** — **donated to OpenTelemetry in April 2023** with intent to converge with SemConv (directional, not a one-time merge — some areas like metrics data models differ).[^elastic-ecs-otel]
- **Splunk CIM** (Common Information Model) — the Splunk normalization target (see `splunk-platform-spl`).[^cribl-central]
- **OCSF (Open Cybersecurity Schema Framework)** — open, vendor-neutral, **format-agnostic** (JSON/Parquet/Avro) schema for **security** events (categories → classes → objects → attributes), explicitly designed for SIEM/SOAR/EDR **and pipelines doing ETL normalization**. Mapping **OTel → OCSF** is a headline transform for SIEM feeds.[^ocsf][^apto-playbook]

**Enrichment via lookups** — add context before storage: **GeoIP** (local MaxMind DB → `geo.*` fields), **threat-intel / IOC** lookups (OCSF models this as the `enrichments` object — happens during processing before immutable storage, with a latency gap from event time), and **CMDB / business dimensions** (actor, action, resource, outcome) so cross-source events join without constant re-mapping.[^signoz-geoip][^ocsf][^cribl-central]

### 4. Fan-out / routing to multiple destinations

A single ingest is copied to many sinks without duplicating agents — e.g. metrics → Prometheus, a subset → Datadog, logs → Elasticsearch **and** S3 simultaneously; fan-out is **independent per sink** so one slow/down sink doesn't block the others.[^datadog-tp][^tell-routing] **Conditional routing** evaluates each record's content (e.g. split by severity), with a default catch-all. Concrete mechanisms: Fluent Bit tag-based vs conditional routing;[^fluentbit-router] OTel's **routing connector** (`move` vs `copy` semantics — copy = replicate to multiple pipelines);[^otel-routing] OpenSearch Data Prepper `route`.[^dataprepper] The canonical destination split — normalized security events → SIEM, operational logs/metrics → observability, full-fidelity/compressed → object-storage lake — is what makes the data-lake pattern below work, and "adding a backend becomes a routing-rule change, not a re-instrumentation project."[^cribl-central][^datadog-tp]

---

## The observability data lake / tiered-telemetry pattern + replay

**The pattern.** Route **full-fidelity** telemetry to cheap object storage (S3/GCS/Azure Blob, or a purpose-built telemetry lake) as the **system of record**, and send only a **reduced, high-value subset** to the expensive SIEM/analytics hot tier; **replay (rehydrate)** from cheap storage into the SIEM on demand for investigations/audits.[^datadog-tp][^cribl-central][^cribl-replay] "Hot data lives in the SIEM for real-time workflows; warm/cold lands in object storage where it can be queried directly and replayed as needed — you keep more detail for longer but only pay hot-tier prices for what truly needs to be hot."[^cribl-central]

**Why it decouples retention cost from analytics cost.** The hot/searchable tier is typically **10–20× more expensive per GB** than cold object storage.[^itbroker] Datadog's **Flex Logs** and the broader "decouple ingestion from indexing" move are the productized expression.[^dd-flex][^dd-pricing] Converting raw JSON to **compressed Parquet** before archival is cited at **60–80% storage savings and up to 10× faster queries**.[^solide]

**Replay / rehydration mechanics.** The full-fidelity copy is written with a **partitioning scheme** (host/date/sourcetype/index/time). On demand the pipeline reads back **only the slice** matching a time range/filter and re-ingests it through the *same parsing/enrichment/routing logic* used for live data.[^cribl-replay][^dd-rehydrate] Critical efficiency technique: **use the path/key (not the contents) as the first filter** so you exclude most objects *without downloading them* — "avoid opening files to check for matches."[^cribl-replay-s3][^honeycomb-s3] **Schema-on-read:** because the lake stores raw/open formats (JSON/Parquet/Avro; OCSF is format-agnostic), structure is applied at query/replay time — which lets you re-interpret old data with new parsing or pull it into a different backend later.[^cribl-central][^ocsf]

**Disconfirming view — rehydration cost & friction is real:**
- Retrieving from cold storage **"can take hours or even days,"** stalling investigations (acknowledged in Datadog's own rehydration blog).[^dd-rehydrate]
- **S3 Glacier standard retrieval = 3–5 hours**, plus retrieval/request fees and the re-indexing cost; downstream systems can be overwhelmed by the inflow.[^edgedelta-rehydrate]
- **Cross-cloud egress** can exceed **$50,000/yr** for 2 TB/day moved between clouds; broad cold queries produce unpredictable bills.[^itbroker]
- A contrarian vendor (Hydrolix) argues tiering's economics can invert: one team's cold queries averaged **8 minutes**, blocking engineers for weeks during an audit and producing >$450K total cost — and markets a single always-warm tier on object storage as the alternative.[^hydrolix]

**Net:** tiering + replay is the dominant cost architecture, but **scope retrieval windows tightly, partition for path-level filtering, store columnar/compressed, and don't over-archive data you will predictably need fast.** Define per-data-class replay objectives up front (how far back must each class be replayable, who can trigger it) — a common failure is a raw topic retaining 6 hours when an investigation needs 2 days.[^automq]

---

