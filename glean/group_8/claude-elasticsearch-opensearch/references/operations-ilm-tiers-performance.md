# Operations & Performance — ILM/ISM, Data Tiers, Shard Sizing, Tuning, Anti-Patterns

> Provenance: reference under the `elasticsearch-opensearch` skill. Researched 2026-06-18 via `/dr` deep research (Elastic ILM/data-tiers/size-your-shards/tune-for-* docs, OpenSearch ISM/Data Prepper docs, Elastic & AWS blogs). **verified-as-of 2026-06-18.** Vendor figures (frozen-tier density, cost ratios) are flagged. **Read the "20 shards/GB is deprecated" flag below before quoting any shards-per-heap rule.**

## Contents

- Index Lifecycle Management (ILM) and OpenSearch ISM
- Data streams (for time-series)
- Data tiers (hot-warm-cold-frozen) and searchable snapshots
- Shard sizing (and the deprecated "20 shards/GB" rule)
- Indexing performance tuning
- Query performance tuning
- Anti-patterns

## Index Lifecycle Management (ILM) and OpenSearch ISM

**ILM (Elasticsearch)** defines five phases — **hot** (actively updated + queried), **warm** (rarely updated, still queried), **cold** (infrequently queried, slower OK), **frozen** (queried rarely, extremely slow OK), **delete**. Indices transition by **`min_age`**, which must increase monotonically across phases; for rolled-over indices, `min_age` is relative to the rollover time, not creation.[^15][^16]

**Actions by phase** (selected):[^15]
- **Rollover** — *hot only*; creates a new write index when the current one hits a `max_*` threshold.
- **Shrink** / **Force merge** — hot + warm (reduce primary shard count / merge segments).
- **Searchable snapshot** — hot, cold, frozen (snapshot to a repository and mount).
- **Downsample** — TSDS only (pre-compute statistical summaries).
- **Allocate / Migrate** — warm + cold (tier movement / replica changes).
- **Delete / Wait for snapshot** — delete phase only.

Note: there is **no standalone "freeze index" action** in modern ILM — the frozen *tier* is reached via the **searchable_snapshot** action with a partially-mounted index (the old `freeze` index API is deprecated/removed in 8.x).[^15]

**Rollover conditions:** a rollover action specifies ≥1 `max_*` condition; it fires when *any* `max_*` is satisfied **and** all `min_*` are satisfied. Conditions: `max_age`, `max_docs`, `max_size`, `max_primary_shard_size`, `max_primary_shard_docs` (+ `min_*`). There is a built-in implicit rollover at **200,000,000 docs per shard**. Canonical policy: roll at `max_primary_shard_size: 50gb` or 30 days; shrink + force-merge in warm; move to cheaper hardware in cold after 7 days; delete at retention.[^16][^17][^19][^20]

**OpenSearch ISM (Index State Management)** — same goal, different execution model. Instead of fixed phases, you define **states**, each with **actions** (run sequentially on entry) and **transitions** (conditions checked after actions complete — index age/size/doc count).[^25][^26] Key contrasts vs ILM:[^27]
- ILM ships hot/warm/cold/frozen/delete phases out of the box; **in ISM you name and build the states yourself** (more flexible, more burden).
- ISM supports per-transition **notifications** (Slack/Email/Chime — risk of alert fatigue); ILM instead writes phase/action outcomes to an index you can alert on.
- The ISM job runs **every 5 minutes by default** (ILM's poll interval is 10 minutes).
- On AWS-managed OpenSearch, the warm tier is **UltraWarm** and there's a **cold** tier, orchestrated via ISM states + snapshot actions.[^28]

## Data streams (for time-series / append-only data)

A **data stream** is the recommended abstraction for append-only time-series data (logs, events, metrics): a single named resource backed by multiple hidden, auto-generated **backing indices**, requiring a matching index template and a `@timestamp` field. Writes route to the current write index; reads fan out to all backing indices. It is **append-only (mostly)** — you cannot update/delete existing docs directly via the stream (target the backing index). ILM automates backing-index management; Elastic recommends data streams over manually managing indices.[^21][^22][^23] A **Time Series Data Stream (TSDS)** (`index.mode: time_series`) is a metrics-optimized variant with time-bound backing indices, `_tsid`+`@timestamp` index sorting for compression, and `routing_path`.[^21][^24]

## Data tiers and searchable snapshots

Tier node roles map to `node.roles`: **`data_hot`, `data_warm`, `data_cold`, `data_frozen`** (plus `data_content`). These were **formalized as roles in Elasticsearch 7.10**, replacing the older `node.attr.node_type` custom-attribute approach; built-in roles let ILM auto-migrate indices across tiers.[^29][^30][^31][^32]

The cost/latency tradeoff:[^29][^30][^33]
- **Hot** — most recent, most-frequently-accessed; handles indexing load; local SSD.
- **Warm** — less-frequently accessed, rarely updated; cheaper hardware.
- **Cold** — infrequently queried; can hold **fully-mounted searchable snapshots** needing **no replicas** (recovery from the snapshot), cutting required local disk by **~50%**; fully-mounted indices are read-only.
- **Frozen** — queried rarely; holds **partially-mounted** searchable snapshots **exclusively**, storing essentially no full local copy and fetching from object storage (S3/blob) on demand with a local cache. Extends storage capacity "by up to 20×" vs the warm tier. Frozen searches are slower because data is sometimes fetched from the object store. Frozen was introduced as tech preview in 7.12; Elastic recommends **dedicated nodes** for it.

**Frozen-tier economics (Elastic's headline, vendor figures — verified-as-of 2026-06):** Elastic states frozen storage cost goes down **up to 90% vs hot/warm and up to 80% vs cold**, by searching object storage directly via searchable snapshots with a local cache for repeat queries.[^33][^30] *Treat the density/cost ratios as Elastic vendor figures.*

## Shard sizing (and the deprecated "20 shards/GB" rule)

- **Target shard size 10–50 GB**, and **keep docs-per-shard below 200 million.** With ILM, set rollover `max_primary_shard_size: 50gb` (ceiling) and `min_primary_shard_size: 10gb` (floor). Time-based data: 20–40 GB shards are commonly cited.[^34][^35][^36]
- **CRITICAL FLAG — the "≤20 shards per GB of heap" rule is DEPRECATED.** It was retired in **Elasticsearch 8.3** after per-shard overhead optimizations (compact metadata serialization, off-heap structures, compressed cluster state), replaced by **field-density-based sizing** plus the per-node shard cap. The old rule (a 30 GB-heap node → ≤600 shards) still appears in older blogs but **should NOT be quoted for 8.3+.** When you see "20 shards/GB," flag it as stale.[^35][^37]
- **Replacement per-node limits:** **≤1,000 non-frozen shards per node** and **≤3,000 frozen shards per dedicated frozen node** (`cluster.max_shards_per_node`).[^34][^37]
- **Master-node heap rule:** **<3,000 indices per GB of heap** on master-eligible nodes (4 GB heap → <12,000 indices).[^34][^36]
- **Hard Lucene ceiling:** each shard (a Lucene index) caps at **~2.1 billion (2³¹−129) docs** — far above the 200M recommendation.[^34]
- **Default since Elasticsearch 7.0 is 1 primary + 1 replica** per index (was 5+1 pre-7.0). The default 1 shard is fine for <50 GB indices but inappropriate for large volumes.[^36][^38]
- Searches run **single-threaded per shard**: too many shards depletes the search thread pool (low throughput); too few/too-large shards slow recovery. The sweet spot balances both.[^34]

## Indexing performance tuning

- **Bulk API** — use bulk requests, not single-doc; find the optimal size by doubling (100 → 200 → 400…); stay under a couple tens of MB per request.[^39][^40]
- **`refresh_interval`** — refresh (making docs searchable) is costly; default is 1s (only on indices searched in the last 30s). For bulk loads, **set `refresh_interval: -1`** to disable refreshes; for steady ingest with search, raise to e.g. `30s`.[^39][^40]
- **Replicas during bulk load** — set **`index.number_of_replicas: 0`** for large initial loads (docs indexed once), then restore afterward. Caveat: no replicas means a single node loss = data loss, so the source must be retriable.[^39][^40]

## Query performance tuning

- **Filter context + caching** — put structured conditions (term, range) in `filter`, not `must`/`should`: filter context skips scoring (faster, less CPU) and is **automatically cached** in the node query cache.[^42][^43][^44]
- **Node query cache** — one per node, shared by all shards, LRU, holds up to 10,000 queries in up to **10% of heap** (`indices.queries.cache.size`). Only **filter-context** results are cacheable; caching is per-segment (segment must have ≥10,000 docs and ≥3% of shard docs); merges invalidate cached entries.[^45]
- **Filesystem cache** — Elasticsearch relies heavily on the OS filesystem cache; **leave at least half of available RAM to it** so hot index regions stay in physical memory.[^41]
- **Caching caveat** — node-level caches help only if repeat requests hit the same shard copy; with replicas + default round-robin routing, two identical requests can land on different copies and miss the cache. *(qualify)*[^41]
- **Avoid expensive queries / scripts** — scripts for filtering/scoring often force a full scan and aren't cacheable. **`search.allow_expensive_queries: false`** blocks wildcard, regexp, fuzzy, range-on-text, etc.[^41][^42][^46]
- **Avoid high open search contexts** — each search opens a per-shard context; scroll contexts held for their timeout cause high JVM memory pressure.[^41]

## Anti-patterns

- **Oversharding** — too many shards degrades search and can make the cluster unstable (GC pressure, slow cluster-state updates, depleted search thread pool). Fixes: wider rollover intervals (weekly/monthly vs daily), fewer primary shards per index, shrink/force-merge/reindex to consolidate, add nodes, and **delete** (don't just stop writing to) unneeded indices.[^34][^36][^37]
- **Mapping explosion / too many fields** — uncontrolled dynamic mapping causes runaway field growth, degrading search, raising JVM memory pressure, and prolonging startup. Default **1,000-field limit** per index (`index.mapping.total_fields.limit`); exceeding it rejects docs ("Limit of total fields [X] has been exceeded"). Overriding to 10,000 is usually fine; 100,000 has strong performance implications. Mitigations: (1) **disable dynamic mapping** (`dynamic: false` ignores new fields; `strict` rejects them); (2) **runtime fields** (`dynamic: runtime` — schema-on-read, no upfront RAM cost, slower queries); (3) the **`flattened`** field type for arbitrary key sets. **Mappings cannot be reduced** once created.[^47][^48][^49][^50]
- **Default 1 shard used inappropriately** — fine for <50 GB; causes giant shards for high-volume indices. Conversely, over-provisioning shards on small data wastes overhead. Match shard count to volume (~1 shard per 10–30 GB).[^36][^38] *(qualify)*
- **Deep pagination (`from`+`size`)** — capped at **10,000 hits** by `index.max_result_window`; deep pages load all preceding hits into memory on every shard, degrading performance or failing nodes. Correct tools, in order:[^51][^52][^54]
  - **`search_after`** — cursor-based, uses the prior page's sort values; the recommended deep-pagination method.
  - **Point In Time (PIT)** — a lightweight frozen view of index state; combine with `search_after` for *consistent* paging across refreshes (open PITs keep old segments alive → heap cost under heavy delete/update).
  - **Scroll** — **no longer recommended** for deep pagination; Elastic steers users to `search_after` + PIT.
- **Wildcard / leading-wildcard queries** — `*` queries are resource-intensive, especially **leading** wildcards: a leading `*`/`?` is "the equivalent of a full table scan." Mitigations: the **`wildcard` field type** (much faster for leading-wildcard/regexp on high-cardinality fields), **n-gram/edge-ngram** tokenizers (faster matching, larger index), pre-narrowing with `match`/`bool`/time-range filters, and `search.allow_expensive_queries: false` to block them.[^46][^55][^56]
- **Unbounded / heavy aggregations** — putting heavy clauses in `must` instead of `filter` and requesting very large `size` can make a cluster unresponsive even over a narrow time range; mitigate with filter-context + bounded result sizes. *(tentative — strong forum example + filter-context docs)*[^56]
- **Single huge index vs time-based indices/data streams** — best practice is **data streams + ILM** for time-series with automatic rollover, and **"delete indices, not documents"** (deleted docs linger until segment merge).[^34][^22]

## References

15. Elastic — ILM phases and actions: https://www.elastic.co/docs/manage-data/lifecycle/index-lifecycle-management — 5 phases, min_age, phase/action table. docs
16. Elastic — ILM hub: https://www.elastic.co/guide/en/elasticsearch/reference/8.19/index-lifecycle-management.html — ILM vs data streams, rollover/shrink/force-merge/delete, example policy. docs
17. Elastic — Rollover (ILM 8.19): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/ilm-rollover.html — hot-only, max_*/min_* conditions. docs
19. Elastic — Configure a lifecycle policy (8.19): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/set-up-lifecycle-policy.html — example, implicit 200M-doc rollover. docs
20. Elastic — Automate rollover with ILM (8.19): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/getting-started-index-lifecycle-management.html — hot-warm-cold architecture, 50GB/30d example. docs
21. Elastic — Time series data stream (TSDS): https://www.elastic.co/docs/manage-data/data-store/data-streams/time-series-data-stream-tsds — index.mode time_series, _tsid sorting. docs
22. Elastic — Data streams (8.19): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/data-streams.html — append-only, backing indices, @timestamp. docs
23. Elastic — Data lifecycle hub: https://www.elastic.co/docs/manage-data/lifecycle — data streams recommended over manual indices. docs
24. Elastic — Time series data stream (8.19 tsds): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/tsds.html — time-bound indices, routing_path, index sorting. docs
25. OpenSearch — Index State Management: https://docs.opensearch.org/latest/im-plugin/ism/index/ — states/actions/transitions, ism_template, 5-min job. docs
26. OpenSearch — ISM policies: https://docs.opensearch.org/latest/im-plugin/ism/policies/ — state/transition definitions, ordered transition evaluation. docs
27. Opster — Elasticsearch ILM vs OpenSearch ISM: https://opster.com/guides/opensearch/opensearch-data-management/opensearch-ism-vs-elasticsearch-ilm/ — phases-built-in vs build-your-own, notifications, outcome storage. blog
28. AWS Big Data Blog — Automating ISM for Amazon OpenSearch Service: https://aws.amazon.com/blogs/big-data/automating-index-state-management-for-amazon-opensearch-service/ — hot-snapshot-warm-cold-delete, UltraWarm. blog
29. Elastic — Data tiers: https://www.elastic.co/docs/manage-data/lifecycle/data-tiers — tier roles, cold ~50% savings, frozen ~20× density, partial vs full mounts. docs
30. Elastic — Node roles: https://www.elastic.co/docs/deploy-manage/distributed-architecture/clusters-nodes-shards/node-roles — data_hot/warm/cold/frozen, cold no-replica savings, frozen needs repo. docs
31. Elastic Blog — Data lifecycle management with data tiers (2021): https://www.elastic.co/blog/elasticsearch-data-lifecycle-management-with-data-tiers — roles formalized in 7.10, searchable snapshots in 7.10. blog
32. Elastic — Configure data tiers (ECK/self-managed): https://www.elastic.co/guide/en/cloud-on-k8s/current/k8s-data-tiers.html — node.roles config, dedicated frozen nodes recommended. docs
33. Elastic Blog — Directly search S3 with the frozen tier (2021): https://www.elastic.co/blog/introducing-elasticsearch-frozen-tier-searchbox-on-s3 — frozen 90% vs hot/warm, 80% vs cold, tech preview 7.12. blog
34. Elastic — Size your shards: https://www.elastic.co/docs/deploy-manage/production-guidance/optimize-performance/size-shards — 10–50GB, <200M docs, 1000/3000 shard cap, 3000 indices/GB master heap, Lucene 2³¹ cap, single-thread-per-shard, fixes. docs
35. Elastic Blog — How many shards should I have? (updated): https://www.elastic.co/blog/how-many-shards-should-i-have-in-my-elasticsearch-cluster — editor's note: 20-shards/GB deprecated in 8.3; 50GB; 20–40GB time-based. blog
36. Elasticsearch Labs — shards & replicas guide (2025): https://www.elastic.co/search-labs/blog/elasticsearch-shards-and-replicas-guide — oversharding, 10–50GB/<200M, default 1+1 since 7.0. blog
37. Elasticsearch Labs — shard & node size best practices (2026): https://www.elastic.co/search-labs/blog/elasticsearch-shard-size-best-practices — 20-shards/GB retired in 8.3 → field-density sizing, 1000 shards/node. blog
38. Elasticsearch Labs — how to reduce shard count (2025): https://www.elastic.co/search-labs/blog/reduce-elasticsearch-shard-count — 1 shard/10–30GB, default 1 since 7.0, weekly/monthly rollover. blog
39. Elastic — Tune for indexing speed: https://www.elastic.co/docs/deploy-manage/production-guidance/optimize-performance/indexing-speed — bulk sizing, refresh_interval -1/30s, replicas 0 for initial load. docs
40. Elastic — Tune for indexing speed (8.19): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/tune-for-indexing-speed.html — same, version-pinned. docs
41. Elastic — Tune for search speed: https://www.elastic.co/docs/deploy-manage/production-guidance/optimize-performance/search-speed — filesystem cache ≥50% RAM, round-robin defeats caches, constant_keyword, avoid high open contexts. docs
42. Elastic — Tune for search speed (8.19): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/tune-for-search-speed.html — caches, expensive queries. docs
43. Elastic — Query and filter context: https://www.elastic.co/docs/reference/query-languages/query-dsl/query-filter-context — filter context binary, faster, cached. docs
44. (filter-context caching — corroborated by [42][43][45]). docs
45. Elastic — Node query cache settings (8.19): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/query-cache.html — per-node, LRU, 10k queries/10% heap, filter-context only, per-segment rules. docs
46. Elastic — Wildcard query: https://www.elastic.co/docs/reference/query-languages/query-dsl/query-dsl-wildcard-query — avoid leading `*`/`?`, n-gram alternative, allow_expensive_queries gate. docs
47. Elastic — Mapping explosion troubleshooting: https://www.elastic.co/docs/troubleshoot/elasticsearch/mapping-explosion — 1000-field default, rejection error, 10k OK/100k bad. docs
48. Elastic — Mapping limit settings: https://www.elastic.co/docs/reference/elasticsearch/index-settings/mapping-limit — total_fields.limit default 1000, flattened type. docs
49. Elastic Blog — 3 ways to prevent mapping explosion (2022): https://www.elastic.co/blog/found-crash-elasticsearch — dynamic false/runtime/flattened, schema-on-read tradeoff. blog
50. Elastic — Dynamic field mapping (8.19): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/dynamic-field-mapping.html — true/runtime/false/strict behavior. docs
51. Elastic — Paginate search results: https://www.elastic.co/docs/reference/elasticsearch/rest-apis/paginate-search-results — from+size 10k cap, search_after, PIT, scroll not recommended. docs
52. Elastic — Point in time API (8.19): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/point-in-time-api.html — PIT lightweight view, search_after consistency, segment retention. docs
54. Elastic — Scroll API (8.19): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/scroll-api.html — "no longer recommend scroll for deep pagination." docs
55. Elastic Blog — Find strings within strings faster (wildcard field, 2020): https://www.elastic.co/blog/find-strings-within-strings-faster-with-the-new-elasticsearch-wildcard-field — keyword vs wildcard, leading-wildcard high-cardinality much faster. blog
56. Elastic Discuss — leading wildcard = full table scan; regexp + size 10000 makes cluster unresponsive: https://discuss.elastic.co/ — use filter clause. forum
