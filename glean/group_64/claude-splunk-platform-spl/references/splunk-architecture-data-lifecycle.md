<!-- Provenance: reference under the `splunk-platform-spl` standalone skill. Created 2026-06-18 via /dr deep-research (multi-source web research, ≥3 independent sources/concept). Volatile claims stamped verified-as-of: 2026-06-18. -->

# Splunk Architecture & Data Lifecycle

`verified-as-of: 2026-06-18` for all licensing, SmartStore-constraint, and enforcement-threshold claims (version-dependent — verify against the running Splunk version and current pricing/Order terms). Splunk docs current to the 9.x–10.4 line. Naming: docs use **manager node / cluster manager** (was "master node") and **license manager** (was "license master"); legacy terms persist in community sources.

## Contents
- Component roles (forwarders, indexer, search head, management tier)
- Distributed topologies (indexer clustering, search head clustering, multisite)
- License model
- Data lifecycle: index-time vs search-time (schema-on-read)
- Indexes & buckets (states, rolling, internals)
- SmartStore (S2)
- Disconfirming findings / gotchas

## Component roles

**Forwarders.** Three types; the light forwarder is deprecated (since Splunk 6.0) and superseded by the universal forwarder.[^1][^2]
- **Universal forwarder (UF):** a separate, dedicated lightweight executable (not a full Splunk instance), smallest footprint, **does not bundle Python**, and **cannot index, search, or alert**.[^2][^3] It sends *unparsed* "cooked" data by default — it tags the stream with source/sourcetype/host metadata, divides it into ~64 KB blocks, and does rudimentary timestamping, but does **not** break individual events or do per-event filtering/routing.[^4][^5] It can parse only limited structured inputs (CSV/JSON).[^3][^5]
- **Heavy forwarder (HF):** a *full Splunk Enterprise instance* with some features disabled — you "enable forwarding" on a full install; there is no separate HF installer.[^6] It *parses* data into events before forwarding (sends *parsed* cooked data), can do **per-event filtering and conditional routing** on field values, and can optionally index locally. It cannot perform distributed searches.[^6][^5][^1]
- **Cooked vs raw:** both parsed (HF) and unparsed (UF) formats are "cooked." Set `sendCookedData=false` to send raw (e.g., to a non-Splunk system).[^5]

**Indexer.** Indexes data: accepts streams from forwarders, parses, writes buckets, maintains buckets, and acts as a **search peer** — accepting requests from search heads, searching its buckets, and streaming results back.[^7][^8] Indexer work is primarily **I/O-bound**; the canonical scaling move is to add indexers (spreading both index and search load).[^7] In clustered deployments indexers must be on dedicated servers.[^8]

**Search head (SH).** Provides the UI, handles user searches, distributes them across indexers (search peers), consolidates results, coordinates scheduled searches, and serves dashboards/search tools.[^1][^8]

**Management tier:**
- **Deployment server:** distributes configuration/app updates ("deployment apps") to **deployment clients** (forwarders, *non-clustered* indexers, search heads) grouped into **server classes**; clients "phone home" on an interval to pull updates.[^9][^10][^11] It does config/content only — not initial installs/upgrades.[^12] Maps clients→server classes→apps via `serverclass.conf` (filters on DNS/IP/build/platform/clientName). Default phone-home 60s; raise to 300s for >2000 clients. **>50 clients ⇒ dedicated host** (and not also an MC or SHC-deployer); >10,000 clients ⇒ tune `dedicatedIoThreads`.[^9][^19] Clustered peers get config from the cluster manager, not the deployment server.[^12]
- **Cluster manager / manager node (CM):** required coordinator for an indexer cluster — tracks buckets, enforces RF/SF, tells search heads where data is, orchestrates fixups.[^13][^14]
- **Search head cluster deployer (SHC-D):** distributes apps/baseline config to SHC members; **stands outside the cluster** and is not a runtime component.[^1][^15]
- **License manager (LM):** central license repository; other instances become **license peers** reporting ingest volume.[^16][^17]
- **Monitoring console (MC, formerly DMC):** centralized monitoring of the whole deployment (6.2+); essentially a search head.[^1][^18]
- **Colocation:** LM, MC, deployment server, and CM can be colocated in various combinations under load thresholds.[^16]

## Distributed topologies

**Basic distributed search:** forwarders → indexers (search peers) ← search head.[^1][^8]

### Indexer clustering (index replication)
A group of indexers replicating each other's data for HA. Node types: one **manager node**, several **peer nodes** (index + search; each can be source and target peer), one or more **search heads**.[^8][^14]
- **Replication factor (RF):** number of *copies* of each bucket on separate peers; cluster tolerates loss of **(RF − 1)** peers. **Default RF = 3.**[^20][^14] Replication is low-CPU; the cost is storage + more indexers.
- **Search factor (SF):** number of *searchable* copies (including the tsidx index files, not just rawdata). **Default SF = 2; SF ≤ RF.**[^21] With SF ≥ 2 a searchable copy takes over immediately if a peer dies; at SF = 1 recovery lags because tsidx must be rebuilt from rawdata.[^21][^14]
- **Cluster states:** a **valid** cluster has exactly one *primary* copy of each bucket (can serve all data); a **complete** cluster has RF copies and SF searchable copies. Complete ⇒ valid, not vice-versa.[^22]
- **Bucket fixup:** when a peer leaves, the CM directs remedial work to restore valid/complete state, iterating internal `to_fix` lists on a ~1s service loop and scheduling async jobs (replicate bucket, make a searchable copy, roll hot bucket, change primaries); a failed job goes back on the list.[^23] **Excess copies** beyond RF/SF (e.g., from a rejoining peer) can be removed to save disk.[^24]
- **Generation ID:** the CM hands search heads a generation ID identifying the current primary-copy set, ensuring each search hits exactly one copy of each bucket.[^14][^22]

### Search head clustering (SHC)
A group of search heads sharing configs, job scheduling, and search artifacts.[^1][^15]
- **Captain:** one member coordinates job scheduling and replication, *and* serves as a normal search head; captaincy can shift.[^15][^1]
- **Election uses RAFT:** randomized election timers (~60s `election_timeout_ms`); the first to time out requests votes; a **quorum of floor(N/2)+1** is required — so an **odd number of members** is recommended (an even cluster fails sooner under multiple losses). A "voting-only" member supplies quorum without running searches.[^25]
- **Knowledge-object replication:** the captain replicates *runtime* changes (saved searches, lookups, dashboards made via Splunk Web) to all members — a member saves locally, sends to the captain, and every ~5s each member pulls and applies new changes. Search artifacts are replicated to RF members.[^26][^1]
- **Deployer vs replication (key split):** runtime KO changes are auto-replicated by the captain; **non-runtime** changes (apps, and configs edited directly in `.conf` files rather than via Web) are pushed from the external **deployer** as a "configuration bundle" — the admin initiates each push. A change made by directly editing `savedsearches.conf` on one member is **not** replicated.[^27][^26]
- **Out-of-sync protection:** the cluster prevents an out-of-sync member from becoming captain (its stale baseline would overwrite everyone's configs); out-of-sync state outranks "preferred captain."[^28][^26]

### Multisite clustering
Site-aware indexer clustering across data centers for DR + search affinity.[^30][^31]
- **`site_replication_factor`** and **`site_search_factor`** replace single-site RF/SF on the CM; syntax like `origin:1, site1:2, total:3` (`origin` = minimum copies on the site where data entered). SF sites must be a subset of RF sites.[^32][^33]
- **Search affinity:** if each site has a full searchable copy set and a local SH, that SH's searches are answered only by *local* peers (still distributed cluster-wide, but only local peers respond), cutting WAN traffic while preserving full-data access. Only **explicitly named** sites get affinity; setting an SH's site to **`site0`** disables affinity.[^31][^32][^34]
- **Gotcha:** the single-site `replication_factor` (default 3) still exists alongside the site factors and can break startup if any site has fewer than 3 peers — lower it to ≤ the smallest site's peer count.[^33]

## License model — VOLATILE
Two coexisting pricing models; the long-standing one is volume/ingest-based.[^35][^36]
- **Ingest (volume) pricing:** priced on **GB/day ingested**, measured as the **raw data entering the indexing pipeline** (not compressed on-disk size); data filtered/dropped *before* indexing does **not** count. Users/searches are unlimited. Per-GB price decreases at higher tiers. **Metrics** events are measured but capped at **150 bytes/event** and draw from the same quota; `_internal`/`_introspection`, summary indexing, and metric rollups don't count.[^37]
- **Workload pricing (newer):** priced on **compute capacity** — measured in **SVC (Splunk Virtual Compute)** units on **Splunk Cloud Platform**, and in **vCPUs** for **Splunk Enterprise** (on-prem). For vCPU licensing, total vCPUs **across all search heads and indexers** count. Best for data you search infrequently.[^35][^36][^38] You **cannot stack volume and vCPU licenses**, and premium products (ES, ITSI) must **match** the core license type.[^39] *(single-source from the Licensed Capacity terms — verify against your Order)*
- **License management (Enterprise, volume):** LM groups licenses into **stacks**, carves **pools** from stacks, assigns **license peers** to pools; daily volume tracked at stack and pool level.[^40][^17]
- **Enforcement (version-dependent — VOLATILE):** for stacks **< 100 GB/day on Splunk 8.1.0+**, **45 warnings in a rolling 60-day window disables search** (per offending pool). For stacks **≥ 100 GB/day on 6.5+**, search is **NOT disabled** on violation. Indexing continues during violation; `_internal` stays searchable; resetting requires a reset license from Splunk Sales. (Legacy docs describe "5 warnings in 30 days" — the older framing.)[^43][^44]

## Data lifecycle: index-time vs search-time (schema-on-read)
Splunk's defining choice: minimal fixed schema at write, structure applied at read — **schema-on-read**.[^45] Most field extraction should be deferred to search time.

**At index time** (between ingestion and disk write): default field extraction (`host`, `source`, `sourcetype`, `_time`), host/sourcetype assignment, event line-breaking, timestamping, segmentation, structured-data field extraction, and any *custom* index-time fields.[^46][^45] Index-time custom fields are **discouraged**: they enlarge the index (slowing indexing and search), are inflexible (changing them requires re-indexing the whole dataset), and become permanent in the bucket.[^46][^47][^48]

**At search time** (while a search runs): search-time field extraction (automatic `key=value`, plus custom/`rex`/calculated/multivalue), field aliasing, lookups, event-type matching, sourcetype renaming, tagging, segmentation.[^46][^49] This is non-persistent — a "virtual overlay" computed against rawdata as it is read, so rules can be iterated endlessly without touching stored data.[^45] **Index-time operations always precede all search-time operations.**[^49]

**The indexing pipeline.** Data enters the **parsing pipeline** in large ~10,000-byte chunks, is broken into events, then handed to the **indexing pipeline** (segmentation, building index structures, writing rawdata + index files with compression).[^51] Docs note the parsing pipeline is itself **three pipelines — parsing, merging, typing**: parsing (UTF-8 decode, line breaking, header handling — a UF does *not* run parsing-pipeline work), merging (line-merge for multi-line events, per-event time extraction), typing (regex replacement/SEDCMD, `punct` extraction); then the index pipeline does output, indexing, byte-quota/license metering, and throughput metrics.[^51][^52] Pipelines are threads connected by queues; a full set is a parallelizable **pipeline set**. *(The per-processor breakdown is well-supported but partly community-sourced via the widely-cited "Masa diagram" thread.[^52])* There is also **ingest-time eval (`INGEST_EVAL`)**: index-time transforms computing fields/lookups before indexing.[^53]

## Indexes & buckets
Indexed data lives in **buckets** — directories holding a **rawdata journal** plus **tsidx** and metadata files; each bucket is bounded by a time range.[^54][^55] Each index is a directory under `$SPLUNK_HOME/var/lib/splunk/<index>`.

**Bucket states & rolling (default local storage):**

| State | Searchable? | Default path | Roll trigger |
|---|---|---|---|
| **Hot** | Yes | `…/<index>/db/*` | New data lands here; actively written; several can be open. Rolls to warm when `maxDataSize` hit, `maxHotBuckets` exceeded, or splunkd restarts.[^54][^56] |
| **Warm** | Yes | `…/<index>/db/*` (renamed) | Not written to; many exist. Rolls to cold when `maxWarmDBCount` (default 300) reached — oldest first.[^54][^57] |
| **Cold** | Yes | `…/<index>/colddb/*` (moved) | Can live on cheaper storage. Rolls to frozen at `frozenTimePeriodInSecs` or `maxTotalDataSizeMB`.[^54][^57] |
| **Frozen** | No | deleted by default, or archived via `coldToFrozenDir`/script | Default `frozenTimePeriodInSecs` ≈ 188,697,600s (~6 yrs). Removed from the index.[^54][^58] |
| **Thawed** | Yes (after rebuild) | `…/<index>/thaweddb/*` | Restored archive; **not** subject to the aging scheme.[^59][^54] |

Rolling is driven by **whichever limit is hit first**; freezing fires when the bucket's *newest* event exceeds the retention age, OR the index exceeds its size, OR (with volumes) the volume exceeds its size — in the volume case Splunk freezes the oldest bucket *across all indexes on that volume*.[^57]

**Bucket internals (conceptual):**
- **rawdata journal (`journal.gz`):** processed events in compressed form, stored as ~128 KB uncompressed **slices**. "Rawdata" isn't literally raw — it's post-parse event data, and since Splunk 4.2 it contains everything needed to **rebuild the tsidx files** if lost (so frozen archiving can keep only rawdata).[^60][^55][^58]
- **tsidx (time-series index) files:** the searchable index pointing into rawdata. Each has a **lexicon** (unique terms/segments found at index time) and a **posting list** (seek offsets into the journal for events containing each term) — like a book index but complete, so tsidx + metadata often **exceed** the compressed rawdata in size. `tstats`/accelerated searches read tsidx **only** (no journal decompression) — orders of magnitude faster.[^62][^63][^61]
- **Bloom filter:** built when a bucket rolls hot→warm by hashing each lexicon term; a search builds its own bloom filter and compares, letting Splunk skip entire buckets that can't contain a match (false positives possible, false negatives not).[^61][^64]
- **Metadata:** `hosts.data`, `sources.data`, `sourcetypes.data`, `bucketManifest`. Segmentation granularity (major/minor segmenters) is a tuning lever: coarser segmentation lowers lexicon cardinality and tsidx size at the cost of search flexibility.[^62][^64]

## SmartStore (S2)
SmartStore decouples compute from storage on the indexing tier: the **remote object store** (Amazon S3, GCS, or Azure Blob) holds the **master copies of warm buckets**, while indexer local storage becomes a **cache**.[^65][^66][^67] Motivation: storage demand outpaces compute demand, so scale them independently — add compute for search load without adding storage, or extend retention without paying for more compute.[^66][^68]

**How it changes the bucket lifecycle:**
- **Hot buckets** are still built in local storage and behave like classic hot buckets (still RF/SF-replicated among indexers). When a hot bucket **rolls to warm, a copy is uploaded to remote storage**, which becomes the master.[^65][^66]
- Once uploaded and acknowledged, the remote store owns resiliency, so **local copies become eligible for eviction**. SmartStore indexes generally have **no cold buckets** — warm → (remote) → frozen.[^67][^69][^58]
- **Cache manager** (one per indexer): on a search, fetches needed warm-bucket files from remote into local cache; evicts least-likely-to-be-searched buckets when full (no data lost — masters persist remotely). It favors **recently created/accessed** buckets and tries to download **only the files a search needs**. **Splunk cannot search the object store directly** — buckets must be localized to cache first.[^70][^71][^72]
- **Optimized for the common case:** Splunk states **97% of searches look back ≤24 hours** with spatial/temporal locality, so the working set stays cached and performance is "similar to local storage" for most searches.[^72][^66] (This is Splunk's stated design assumption, not an independently audited figure.)
- **Cluster operation:** RF/SF apply only to **hot** buckets; once a bucket is uploaded, replicating warm buckets requires only **metadata** replication (far faster/cheaper fixup), and SmartStore can recover warm buckets even when **≥ RF peers** are lost. Cluster-wide freezing: the CM runs a search every ~15 min to find buckets to freeze.[^73][^66][^69]
- **Sizing:** provision local cache for ~30 days of indexed data (min 7–10 days); indexed data ≈ ~50% of ingested (compression); **10 Gbps** indexer↔store link recommended.[^74]

## Disconfirming findings / gotchas
1. **"Use a heavy forwarder to filter before indexing" is often wrong.** Because an HF sends *parsed/cooked* data with all index-time fields + metadata, using HFs as intermediate filters can *increase* indexer network I/O (and sometimes CPU/memory). Splunk's guidance: do parsing/filtering on the indexers and use UFs, **unless** you drop a large fraction of data at source, need add-ons requiring a full instance (DBConnect, Checkpoint OPSEC LEA), or need complex per-event routing. A UF **cannot** filter on regex.[^81][^82]
2. **UFs are unsafe for advanced intermediate routing** — `useAck`/event-level processing on a UF risks duplicate/missing data, runaway queuing, and large-event truncation; use UFs as intermediate forwarders only for basic passthrough.[^83][^5]
3. **SmartStore punishes rare/long-lookback searches.** Searches over long timespans or infrequently-searched data force large remote→local downloads; Splunk explicitly says SmartStore "might not be appropriate" for frequent rare/long-lookback patterns. A .conf benchmark: a single 750 MB bucket fetch takes ~1.5s at 0ms latency but **~25s at 100ms**; a cache-miss search can jump from ~2s to >100s, and cache churn from one bad long search can evict other indexes' buckets and hurt ingestion.[^75][^76][^78]
4. **SmartStore config constraints (VOLATILE — verify per version):** with indexer clusters **RF must equal SF**; `homePath` and `coldPath` must be on the **same partition**; **tsidx reduction is unsupported**; each remote volume serves **only one cluster/standalone indexer**; one remote storage **type** per indexer/cluster; all peers must share SmartStore settings. For multisite, SmartStore indexes using report/data-model acceleration require disabling search affinity (all SHs → site0). A full-cluster rolling restart stampedes the object store re-downloading buckets.[^75][^79][^80]
5. **The cluster manager is a single coordinator (not auto-HA by default).** If it goes down the cluster keeps running "for a while" — peers keep ingesting/replicating/searching — but **no new fixups happen**: if a peer with primaries dies while the CM is down, primacy can't transfer and those buckets drop out of search. HA requires **cluster manager redundancy** (active/standby CMs + a load balancer/DNS so peers follow the active CM via `manager_uri`); multisite base topology has no automatic management-plane failover.[^84][^85][^86][^88]
6. **SHC split-brain:** an even-numbered SHC, or a two-site SHC where the majority side fails, **cannot elect a captain** → remediation is a temporary **static captain**.[^25][^29]

## Adjacent / frontier concepts
Splunk Edge Processor; Ingest Processor (tiered by daily processing volume); Ingest Actions (in-product filter/mask/route at index time — must not share a SmartStore S3 bucket); Federated Search (query data in place); Workload Management (admission/workload rules, mitigates SmartStore cache churn); indexer discovery (forwarders learn the peer list from the CM); DDAA / Dynamic Data Active Archive (Splunk Cloud retention tiers, the Cloud analog of frozen/thawed).

## References
[^1]: Deployment topologies — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/inherit-a-splunk-deployment/10.4/inherited-deployment-tasks/deployment-topologies — component roster, SHC, forwarder roles, HF can't do distributed search.
[^2]: Types of forwarders — Splunk Docs — https://help.splunk.com/en/data-management/transform-and-route-data/perform-basic-data-processing/process-data-with-forwarders/types-of-forwarders — UF/HF/light definitions, light deprecated 6.0.
[^3]: The universal forwarder — Splunk Docs — https://help.splunk.com/en/data-management/transform-and-route-data/perform-basic-data-processing/process-data-with-forwarders/the-universal-forwarder — UF limitations (no index/search/alert, no Python).
[^4]: Types of forwarder data — Splunk Docs — https://help.splunk.com/en/data-management/transform-and-route-data/perform-basic-data-processing/process-data-with-forwarders/types-of-forwarder-data — unparsed/parsed/raw, 64 KB blocks, cooked-data definition.
[^5]: Forwarder comparison — Splunk Docs — https://help.splunk.com/en/data-management/transform-and-route-data/perform-basic-data-processing/process-data-with-forwarders/forwarder-comparison — capability matrix per forwarder type.
[^6]: Deploy a heavy forwarder — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/forward-and-process-data — HF is a full instance with forwarding enabled; no HF installer.
[^7]: Distribute indexing and searching — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/get-started/deployment-capacity-manual/9.4/scale-your-splunk-enterprise-deployment/distribute-indexing-and-searching — indexer responsibilities, I/O-bound, scale by adding indexers.
[^8]: Indexes, indexers, and indexer clusters — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers/10.2/indexing-overview/indexes-indexers-and-indexer-clusters — indexer + cluster node types.
[^9]: Deployment server architecture — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/update-your-deployment/9.3/deployment-server-and-forwarder-management/deployment-server-architecture — clients/apps/server classes.
[^10]: Forwarder management overview — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/update-your-deployment — serverclass.conf, forwarder management UI.
[^11]: Use forwarder management to define server classes — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/update-your-deployment — server-class→app mapping.
[^12]: About deployment server and forwarder management — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/update-your-deployment — config/content only, non-clustered targets.
[^13]: Topology components (SVA) — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/get-started/splunk-validated-architectures/introduction-to-splunk-validated-architectures/topology-components — CM as required coordinator; indexers dedicated.
[^14]: The basics of indexer cluster architecture — Splunk Docs — https://help.splunk.com/en?resourceId=Splunk_Indexer_Basicclusterarchitecture — manager role, RF, search coordination, generation ID.
[^15]: Search head clustering architecture — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/distributed-search/10.4/overview-of-search-head-clustering/search-head-clustering-architecture — captain role, deployer outside cluster, artifact replication.
[^16]: Configure a license manager — Splunk Docs — https://help.splunk.com/en/data-management/splunk-enterprise-admin-manual/10.4/configure-splunk-licenses/configure-a-license-manager — LM as repository, license peers, colocation.
[^17]: Allocate license volume — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/admin-manual/10.4/configure-splunk-licenses/allocate-license-volume — stacks/pools/peers model.
[^18]: About the license usage report view — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/admin-manual — MC=search head, LURV.
[^19]: Estimate deployment server performance — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/update-your-deployment — phone-home interval, >50/>2000/>10000 client thresholds.
[^20]: Replication factor — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — RF default 3, tolerate RF-1.
[^21]: Search factor — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — SF default 2, SF≤RF, searchable=tsidx, rebuild lag at SF=1.
[^22]: Indexer cluster states — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — valid vs complete, one primary per bucket.
[^23]: Indexer Clustering Fixups (.conf 2017) — Splunk .conf talk — https://conf.splunk.com/files/2017/slides/indexer-clustering-fixups-how-a-cluster-recovers-from-failures.pdf — to_fix lists, service loop, fixup job types.
[^24]: Remove excess bucket copies — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — excess copies from rejoining peers.
[^25]: SHC captain election — Splunk Community — https://community.splunk.com/t5/Deployment-Architecture/question-related-to-search-head-captain/m-p/757149 — RAFT, floor(N/2)+1 quorum, odd-member rationale, two-site majority-loss → static captain.
[^26]: Configuration updates that the cluster replicates — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/distributed-search — runtime vs deployer; 5s pull; direct .conf edits not replicated; out-of-sync protection.
[^27]: Use the deployer to distribute apps — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/distributed-search — configuration bundle, manual push.
[^28]: Control captaincy — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/distributed-search — preferred captains, out-of-sync outranks preferred.
[^29]: Runtime considerations (SHC) — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/distributed-search — coordination latency; election window; two-site majority-loss → static captain.
[^30]: Distributed Clustered Deployment - Multisite (SVA) — Splunk Docs — https://help.splunk.com/en/data-management/splunk-validated-architectures/splunk-platform-indexing-and-search/distributed-clustered-deployment---multisite-m2--m12 — DR, single CM, latency limits, no auto management failover.
[^31]: Multisite indexer clusters — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers/10.0/overview-of-indexer-clusters-and-index-replication/multisite-indexer-clusters — site-awareness, search affinity.
[^32]: Implement search affinity in a multisite indexer cluster — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — affinity mechanics, origin, explicit-site requirement.
[^33]: Configure the site replication factor — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — site_replication_factor syntax/origin; single-site RF gotcha.
[^34]: How search works in an indexer cluster — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers/10.4/how-indexer-clusters-work/how-search-works-in-an-indexer-cluster — affinity search flow, generation ID, site0.
[^35]: Pricing Models — Splunk (vendor) — https://www.splunk.com/en_us/products/pricing/pricing-models.html — workload vs ingest; SVC (Cloud) vs vCPU (Enterprise). VOLATILE.
[^36]: Splunk Platform Pricing FAQ — Splunk (vendor) — https://www.splunk.com/en_us/products/pricing/faqs/enterprise-and-cloud.html — SVC/vCPU, GB/day, volume-discount curve. VOLATILE.
[^37]: How Splunk Enterprise licensing works — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/admin-manual — raw-pipeline measurement, dropped data not counted, metrics 150B cap, vCPU=all SH+indexers.
[^38]: What is Splunk Virtual Compute (SVC)? — Splunk blog — https://www.splunk.com/en_us/blog/platform/what-is-splunk-virtual-compute-svc.html — SVC definition. VOLATILE.
[^39]: Licensed Capacity — Splunk legal — https://www.splunk.com/en_us/legal/licensed-capacity.html — can't stack volume+vCPU; premium must match core. VOLATILE.
[^40]: Create or edit a license pool — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/admin-manual — auto_generated_pool_enterprise default.
[^43]: Splunk Enterprise License Enforcement FAQ — Splunk (vendor) — https://www.splunk.com/en_us/resources/splunk-enterprise-license-enforcement-faq.html — <100 GB & 8.1+ → 45/60 disables search; ≥100 GB & 6.5+ → not disabled. VOLATILE.
[^44]: About license violations — Splunk Docs — https://help.splunk.com/en/data-management/splunk-enterprise-admin-manual/9.3/manage-splunk-licenses/about-license-violations — 45/60 rule, indexing continues, _internal searchable, reset from Sales.
[^45]: Schema-on-Read and Field Extraction — Blasetti (blog, 2026) — https://blog.blasetti.cloud/splunk-schema-on-read-and-field-extraction-8461f18b04f5 — index-time "in stone" vs search-time "overlay".
[^46]: Index time versus search time — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers/10.4/indexing-overview/index-time-versus-search-time — definitive list of index-time vs search-time processes.
[^47]: About indexed field extraction — Splunk Docs — https://help.splunk.com/en/data-management/get-data-in — indexed vs search-time, re-index requirement.
[^48]: When Splunk software extracts fields — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/fields-and-field-extractions/when-splunk-software-extracts-fields — default fields at index time, caution on custom indexed fields.
[^49]: The sequence of search-time operations — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects — index-time precedes all search-time ops.
[^51]: How indexing works — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers/10.4/indexing-overview/how-indexing-works — 10,000-byte chunks, parsing vs indexing pipeline, parsing=parsing+merging+typing, pipeline set.
[^52]: Diagrams of how indexing works ("Masa") — Splunk Community — https://community.splunk.com/t5/Getting-Data-In/Diagrams-of-how-indexing-works-in-the-Splunk-platform-the-Masa/m-p/590774 — per-processor breakdown; UF skips parsing.
[^53]: Process events with ingest-time eval — Splunk Docs — https://help.splunk.com/en/splunk-cloud-platform — INGEST_EVAL transforms.
[^54]: How the indexer stores indexes — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — rawdata+tsidx+metadata, bucket states/paths, rolling triggers.
[^55]: Buckets and indexer clusters — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — rawdata=processed-not-raw, journal can regenerate index files.
[^56]: How the indexer stores indexes (hot roll triggers) — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — maxDataSize, restart roll.
[^57]: Index rolling off data before retention — Splunk Community — https://community.splunk.com/t5/Deployment-Architecture/Index-rolling-off-data-before-retention-age/m-p/684799 — indexes.conf knobs, volume-level freezing.
[^58]: Archive indexed data — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers/9.1/back-up-and-archive-your-indexes/archive-indexed-data — frozen=delete-or-archive, post-4.2 keeps only rawdata.
[^59]: Restore archived indexed data — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — thaweddb not subject to aging, splunk rebuild.
[^60]: Index Data Search — Mata (blog, 2022) — https://dante-cyber.medium.com/splunk-enterprise-index-data-search-78e32d73a296 — journal slices ~128 KB, tsidx→journal pointers, bloom filter build.
[^61]: SplunkConf – How search works — blog/talk notes — https://xmlisse.wordpress.com/2020/05/20/splunkconf-how-search-works/ — bucket contents, slices, bloom filter false-positive-only.
[^62]: Splunk bucket lexicons and segmentation — Waddle (blog) — https://www.duanewaddle.com/splunk-bucket-lexicons-and-segmentation/ — tsidx = values list + lexicon + posting list; segmentation→cardinality→size.
[^63]: What exactly is a tsidx file? — Splunk Community — https://community.splunk.com/t5/Splunk-Search/What-exactly-is-a-tsdix-file/m-p/269709 — tsidx terms+offsets, often > rawdata, tstats reads tsidx only.
[^64]: Splexicon:Tsidxfile — Splunk Docs — https://docs.splunk.com/Splexicon:Tsidxfile — canonical tsidx definition + metadata files.
[^65]: SmartStore architecture overview — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — remote=master of warm buckets, local cache, hot→warm upload.
[^66]: About SmartStore — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — decouple compute/storage, single permanent warm copy, fast fixup, recover when ≥RF peers lost, limitations.
[^67]: SmartStore for Splunk platform (SVA) — Splunk Docs — https://help.splunk.com/en/data-management/splunk-validated-architectures/splunk-platform-indexing-and-search/smartstore-for-splunk-platform — elasticity rationale, can't search object store directly, eviction order.
[^68]: About SmartStore (scale independently) — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — scale compute/storage independently.
[^69]: Configure data retention for SmartStore indexes — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — cluster-wide freezing, CM 15-min freeze search, no cold buckets.
[^70]: Configure the SmartStore cache manager — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers/9.3/manage-smartstore/configure-the-smartstore-cache-manager — [cachemanager] stanza, per-peer, global eviction.
[^71]: The SmartStore cache manager — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — upload on roll, fetch on search, evict.
[^72]: How search works in SmartStore — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — 97% searches ≤24h, locality, fetch-only-needed-files.
[^73]: Indexer cluster operations and SmartStore — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — SF for hot vs warm, metadata-only warm replication.
[^74]: SmartStore system requirements — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — cache ≈30 days, indexed≈50% ingested, 10 Gbps.
[^75]: About SmartStore (limitations) — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — rare/long-lookback unsuitable, RF=SF, tsidx-reduction unsupported, multisite accel→site0.
[^76]: SmartStore – Spend Less (.conf 2019) — Splunk .conf talk — https://conf.splunk.com/files/2019/slides/FN1435.pdf — cache-miss latency/network benchmarks (750 MB at 0/30/100ms).
[^78]: Reducing SmartStore cache churn with WLM — Splunk Lantern — https://lantern.splunk.com — long-lookback → cache miss → cross-index eviction → ingestion impact.
[^79]: Choose the storage location for each index — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — one remote volume per cluster/index, one storage type, all peers same settings.
[^80]: SmartStore Behaviors — Splunk Community — https://community.splunk.com/t5/Knowledge-Management/SmartStore-Behaviors/m-p/402790 — restart flushes cache, rolling-restart stampede, single-primary-searchable enforcement.
[^81]: Universal or Heavy, that is the question? — Splunk blog — https://www.splunk.com/en_us/blog/tips-and-tricks/universal-or-heavy-that-is-the-question.html — HF-as-filter increases network/CPU/memory; filter on indexers; when to use HF.
[^82]: Universal or Heavy (UF can't regex-filter) — Splunk blog — https://www.splunk.com/en_us/blog/tips-and-tricks/universal-or-heavy-that-is-the-question.html — filter on indexer unless dropping most data at source.
[^83]: Intermediate data routing using UF/HF (SVA) — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/splunk-validated-architectures/getting-data-in-forwarding-and-preprocessing/intermediate-data-routing-using-universal-and-heavy-forwarders — UF intermediate only for basic; advanced→dup/missing data.
[^84]: What happens when the manager node goes down — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — cluster runs "for a while," no fixups, primacy can't transfer.
[^85]: Why multisite requires 1 manager node — Splunk Community — https://community.splunk.com/t5/Deployment-Architecture/Why-does-multisite-clustering-require-1-manager-node/m-p/644976 — single site-less CM per cluster.
[^86]: Implement cluster manager redundancy — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers/10.2/configure-the-manager-node/implement-cluster-manager-redundancy — active/standby CMs, LB/DNS, manager_uri.
[^88]: Handle manager site failure — Splunk Docs — https://help.splunk.com/en/data-management/manage-splunk-enterprise-indexers — manual standby-manager bring-up on site loss.
