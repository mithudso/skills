# Log Analytics, Observability & SIEM — and Positioning vs Splunk

> Provenance: reference under the `elasticsearch-opensearch` skill. Researched 2026-06-18 via `/dr` deep research (Elastic/OpenSearch docs + multiple independent Splunk-vs-Elastic comparisons). **verified-as-of 2026-06-18.** Cost *direction* is well-corroborated; specific dollar figures are third-party estimates with wide variance — flagged. This reference's commercial-SIEM counterpart is the sibling skill **`splunk-platform-spl`**, which reciprocally forward-references here.

## Contents

- Elasticsearch/OpenSearch as a log/observability/SIEM backend
- Architecture contrast: schema-on-write vs schema-on-read
- The cost model (the central positioning point)
- The honest tradeoff
- When each wins

## Elasticsearch/OpenSearch as a log/observability/SIEM backend

Both engines ingest machine data at scale and serve **log management, search, analytics, SIEM, and observability** — the same use cases as Splunk, which is why the two overlap directly.[^57][^58] Elastic ships a **Security (SIEM)** solution and an **Observability** solution inside Kibana; OpenSearch ships equivalent **observability** + **security-analytics** plugins. ES/OpenSearch is the **open-source counterpart** to Splunk; the Splunk platform itself (SPL, knowledge objects, Splunk ES, Splunk Observability Cloud) is owned by the sibling `splunk-platform-spl` skill.

## Architecture contrast: schema-on-write vs schema-on-read

- **Elasticsearch ≈ schema-on-write.** It is a document store on a Lucene inverted index that **maps and indexes fields at ingest** — fast search/aggregations (especially on high-cardinality fields), but more upfront mapping work.[^57][^58][^59]
- **Splunk = schema-on-read.** It indexes **raw events** and extracts fields **at search time** — more tolerant of evolving/unknown log formats, but slower for complex queries.[^57][^59]
- Elasticsearch supports a **hybrid** schema-on-read path via **runtime fields**; Splunk normalizes via CIM field aliases at search time.[^59][^60]
- **SIEM consequence:** Splunk is more forgiving of imperfect normalization, whereas Elastic's prebuilt detection rules depend on correct **ECS** (Elastic Common Schema) field mappings *at index time* — bad/missing ECS mappings cause rules to **silently miss events.**[^60]
- **Query languages:** Splunk uses **SPL**; Elastic uses **Query DSL** (+ **KQL** in Kibana, and **ES|QL** in newer versions); OpenSearch uses **Query DSL** + **PPL/SQL** (PPL is itself Splunk-SPL-inspired — the closest cross-over).[^57][^60]
- **Visualization:** Kibana / OpenSearch Dashboards (highly customizable, more setup) vs Splunk dashboards / Dashboard Studio (more pre-built out of the box).[^57][^61]

## The cost model (the central positioning point)

- **Splunk = ingest/volume-based licensing (per GB/day)** — "becomes very expensive at scale" and complicates budgeting during traffic spikes.[^57][^58][^62] Indicative third-party figures (**verified-as-of 2025–2026, ballpark, NOT vendor list prices — wide variance from negotiation/volume/Cloud-vs-self-managed**): free tier capped at 500 MB/day; production roughly $1,800/yr per GB/day and up; security-focused estimates put Splunk Enterprise + Enterprise Security in the low-hundreds-of-dollars per GB/day, so 100 GB/day lands in the multi-million-dollars/yr range in license alone.[^61][^62][^63][^64]
- **Elastic = resource/capacity-based**, not raw-ingest. Self-managed **Basic tier is free** (pay only infrastructure); Elastic Cloud bills by compute/consumption units; paid subscriptions add per-node cost for security/ML/CCR.[^59][^61][^64]
- **OpenSearch = $0 software at any scale** (Apache 2.0; infrastructure-only on Amazon OpenSearch Service or self-managed).[^28][^30]
- Multiple independent sources estimate Elastic/OpenSearch deliver **meaningfully lower licensing cost** at equivalent ingest, though the **TCO gap narrows** once you count the **extra engineering/FTE** to build and maintain detection content and run the cluster.[^60][^64]

## The honest tradeoff (recurring across sources)

- **Elasticsearch/OpenSearch:** lower software cost, open-source, no vendor lock-in, horizontal scale, stronger vector/semantic search (k-NN, ELSER) — **but** a steeper learning curve and meaningfully more operational/engineering effort (cluster sizing, shard management, detection-content authoring).[^57][^58][^60][^64]
- **Splunk:** easier onboarding, a more mature/complete SIEM with extensive prebuilt detections, better out-of-box experience — **but** higher cost and acute ingest-pricing pain at scale.[^57][^60][^64]
- The decision hinges on **ingest volume, team skills, and whether you already run Elastic/OpenSearch for observability** (a shared cluster tips the balance toward the open-source stack).[^64]

## When each wins (synthesized)

- **Lean toward Elasticsearch/OpenSearch** when: cost-per-GB at high ingest is the binding constraint; you want open-source / no lock-in / Apache-2.0 procurement (OpenSearch); you already operate the stack for search or observability; you want semantic/vector search alongside logs; you have the engineering capacity to own a cluster and detection content.
- **Lean toward Splunk** when: you need a turnkey, analyst-firm-leading SIEM with extensive prebuilt detections; your team lacks cluster-operations capacity; schema-on-read tolerance for messy/varied log formats matters more than per-GB cost; time-to-value outweighs licensing spend.

> For Splunk-internal mechanics (SPL command types, indexer/search-head clustering, SmartStore, CIM, Enterprise Security / RBA, Splunk Observability Cloud + OTel Collector), route to **`splunk-platform-spl`**. This reference covers only the head-to-head positioning from the Elastic/OpenSearch side.

## References

57. unanswered.io — Splunk vs Elastic (2026): https://unanswered.io/blog/splunk-vs-elastic — schema-on-read vs schema-on-write, ingest vs resource pricing, SPL vs Query DSL. blog
58. kloudfuse.com — Splunk vs Elastic (2025): https://www.kloudfuse.com/blog/splunk-vs-elastic — schema-on-ingest CPU/disk cost, raw+index duplication, Lucene document model. blog
59. passitexams.com — Splunk vs Elastic (2026): https://www.passitexams.com/blog/splunk-vs-elastic — schema-on-write mappings, resource pricing, free Basic tier. blog
60. decryptiondigest.com — Splunk vs Elastic SIEM (2026): https://decryptiondigest.com/splunk-vs-elastic-siem — ECS-at-index vs CIM-at-search silent-miss risk, TCO table. blog
61. modern-datatools.com — Elasticsearch vs Splunk (2026): https://modern-datatools.com/elasticsearch-vs-splunk — Elastic Cloud tiers, Splunk 500MB free / ~$1,800/yr per GB/day, Splunk ES maturity. blog
62. last9.io / decryptiondigest — Splunk ingest pricing prohibitive, per-GB/day + ES add-on, 100GB/day multi-million: https://last9.io/blog/splunk-pricing — vendor lock-in. blog
63. uptrace.dev — Splunk Pricing (2025): https://uptrace.dev/blog/splunk-pricing — ~$1,800/GB/yr, volume-based pain. blog
64. opsiocloud.com — ELK vs Splunk cost/migration (2026): https://www.opsiocloud.com/blog/elk-vs-splunk — Splunk per-GB/day; ELK no per-GB cost, per-node subscriptions; TCO + FTE counts; Elastic k-NN/ELSER advantage. blog
