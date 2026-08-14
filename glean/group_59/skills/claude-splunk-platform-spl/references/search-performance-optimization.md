<!-- Provenance: reference under the `splunk-platform-spl` standalone skill. Created 2026-06-18 via /dr deep-research (multi-source web research, ≥3 independent sources/concept). The highest-value reference for a query-optimization-minded reader. -->

# Splunk Search Performance & Optimization

Command semantics and optimization principles here are stable (sourced to the official Splunk Search Manual, Lantern, .conf talks, and Splunk Community). This is the most actionable reference in the skill. Bucket lifecycle and the full SPL command reference live in sibling references — buckets and `tstats`/`stats` appear here only as performance levers.

## Contents
- Time-bounding (the biggest lever)
- Filter early / search-time best practices
- `tstats` over raw search
- Report vs data-model acceleration vs summary indexing
- The Search Job Inspector
- Streaming vs transforming command placement (map-reduce)
- Search modes (Fast / Smart / Verbose)
- Other levers (indexed fields, lookups vs joins, KV Store, concurrency, WLM, metrics)
- Optimization checklist / anti-patterns
- Disconfirming findings / gotchas

## Time-bounding — the single biggest lever
When a search runs, Splunk uses the tsidx files to identify which events to retrieve from disk; "the smaller the number of events to retrieve from disk, the faster the search runs."[^1][^7] The single most effective way to limit data pulled from disk is to **limit the time range.**[^2][^3] At index time every timestamp is converted to UNIX epoch and stored in `_time`; data is written into age-organized **buckets** carrying earliest/latest metadata, so a tight time range lets Splunk **eliminate whole buckets before reading them** — the Search Performance Evaluator literally reports "percentage of buckets eliminated from a search (bigger is better)."[^3][^6] The Search & Reporting app defaults to **Last 24 hours** precisely "to avoid running searches with overly-broad time ranges that waste system resources."[^3]

"All time searches consume a lot of resources by scanning unnecessary data"[^4]; an `index=* ... All time` search is the textbook target for admission rules that **block** searches outright. Gotchas: with a relative `earliest`, also set `latest=now` — omitting it "might produce unpredictable results or impact performance"; `latest=<large_number>` returns *future-dated* events; time boundaries use the *indexer's* timestamp/timezone, not the searching user's.[^5][^3]

## Filter early / search-time best practices
The governing rule: "Filter the data as early as possible in the search, so that processing is done on the minimum amount of data necessary."[^2][^1] Filtering before any transforming/eval work means every downstream command sees fewer events.

**The canonical "tale of two searches"** (Splunk's own load-bearing illustration):[^7]

Inefficient:
```
sourcetype=my_source | lookup my_lookup_file D OUTPUTNEW L | eval E=L/T | search A=25 L>100 E>50
```
→ index returns **1,000,000** events; `lookup` and `eval` run on all 1M; the final `search` cuts to 50,000.

Optimized:
```
sourcetype=my_source A=25 | lookup my_lookup_file D OUTPUTNEW L | search L>100 | eval E=L/T | search E>50
```
→ moving `A=25` into the base search drops the index extract to **300,000**; moving `L>100` right after the lookup drops to **200,000** before `eval`; `E>50` *cannot* move because it depends on the eval result. Same 50,000 results, far less cost.[^7]

**Concrete tactics (doc-grounded):**
- **Always pick an index** in the first line. `index=*` "slows down a search significantly"; partition dissimilar data (web vs firewall) into separate indexes and constrain to the relevant one.[^6][^2][^7]
- **Use the most specific terms possible**; specificity beats wildcards.[^2]
- **Use index-time metadata fields** (`index`, `source`, `sourcetype`, `host`) to filter — they are tagged on every event and cheap.[^2]
- **Avoid leading/mid-word wildcards.** A bare `*` matches everything up to the limit. A *trailing* wildcard (`word2*`) is fairly fast (Splunk uses the reverse-word index), but a **leading** wildcard forces examining *every* event because bloom filters and metadata can no longer prune. `access denied` is "always better than" `access*`. Mid-word wildcards can also silently return *wrong* results — tokens segment on punctuation (`SC-MG-G01` → `SC`,`MG`,`G01`, so `S*G01` matches nothing).[^9][^10]
- **Minimize `OR` and large OR-lists.** Each `OR` adds complexity; huge OR-lists, subsearches (which expand into OR-lists), and phrase searches are a top cause of slow base searches — prefer lookups/eventtypes/regex.[^8][^2]
- **Put non-streaming commands (`stats`, `sort`, `dedup`) as late as possible** (they need the full set; see map-reduce below).[^2][^3]
- **Combine evals** into one comma-delimited `eval`; do evals/lookups after the dataset is as small as possible.[^16][^1]

**The `NOT`/sparse anti-pattern.** Frequent `NOT field=value` (or `foo!=bar`) is slow when the field nearly always equals that value, because Splunk retrieves the events and *then* discards them after extraction.[^23] Splunk classifies searches by density: **dense** (≥10% match, CPU-bound), **sparse** (.01–1%, CPU-bound), **super-sparse** (I/O-bound, "up to 2 sec per bucket" — must scan all buckets), and **rare** (I/O-bound but rescued by **bloom filters** that eliminate non-matching buckets, 20–100× faster than super-sparse).[^21]

## `tstats` over raw search
`tstats` performs statistical queries directly on **indexed fields in the `.tsidx` files** (and accelerated data-model summaries) "instead of scanning raw events," so it "is faster than the `stats` command."[^12][^13] The `.tsidx` files act like a columnar store of `field::value` pairs plus locations.[^13][^14]

**What `tstats` can see:** only index-time fields by default (`sourcetype`, `host`, `source`, `_time`). To expose custom fields, use an **accelerated data model**, `tscollect`, index-time fields via `fields.conf`/`props.conf`/`transforms.conf`, or `INDEXED_EXTRACTIONS` for structured data.[^13]

**Accelerated data models + `tstats`:** data-model acceleration builds the **High Performance Analytics Store (HPAS)** — `.tsidx` summaries on indexers, parallel to the buckets. `tstats` and Pivot read the full HPAS even when distributed.[^15][^20] Important: a **plain regular search and the `datamodel` command are NOT accelerated** — only `tstats` and Pivot read the HPAS; `from datamodel=…` fetches the *original raw events* and is slower.[^19][^20]
- `summariesonly=true` → results come *only* from the summary, even if the range exceeds it; "the search runs fast, but no unsummarized data is included."[^12][^18]
- `summariesonly=false` (default) → uses accelerated data where available, then searches raw for the rest; past the summary range it is "only partially accelerated."[^13][^18]
- `prestats=true` → chainable map-reduce mode that must be followed by `stats`/`chart`/`timechart`; `tstats` must be the *first* command.[^17][^19]

**Disconfirming — when `tstats` is *slower*:** a **high-cardinality field in the `WHERE` clause** makes `tstats` evaluate many distinct tsidx values, and that overhead can outweigh the benefit — a regular search using exact term lookups can be faster.[^12] `tstats` is also "not very good at returning many, many results."[^13]

## Report vs data-model acceleration vs summary indexing
Three "summary-based search acceleration" methods:[^26]

| Method | Accelerates | Where stored | When to use |
|---|---|---|---|
| **Report acceleration** | One transforming report/search | `.tsidx` summaries in `…/index/summary` | Easiest; transforming search slow over large volume.[^26][^27] |
| **Data model acceleration (DMA)** | The whole *set of fields* in a data model, via HPAS | `.tsidx` in `datamodel_summary` dirs on indexers | Fields heavily used in Pivot/dashboards; can beat report acceleration for complex searches.[^26][^20] |
| **Summary indexing (SI)** | Precomputed results into a *separate index* via scheduled `si*`/`collect` searches | A dedicated summary index | Data must outlive raw retention; transforming search that doesn't qualify for report acceleration.[^26][^28] |

**Mechanics:** report acceleration runs a background search (~10 min refresh) that builds a summary; the next run runs against the smaller summary.[^27] DMA uses scheduled summarization searches (~5 min) to populate the HPAS, with periodic maintenance purges.[^33] SI uses `sistats`/`sitimechart`/`sichart` (or `collect`) to write summarized rows; you then report off the SI with `stats` substituted for `sistats`.[^31][^26]

**Tradeoffs (Splunk's "When should I use…?"):** use SI when (a) the report has non-streamable commands *before* the transforming command (disqualifies report acceleration), or (b) raw data rolls more frequently than the reporting window — because **report-acceleration summaries point to raw buckets and expire alongside them**, whereas **SI is a separate copy that can outlive the raw data.** SI also supports multi-layer rollups (raw→daily→monthly→yearly). Report acceleration is **self-repairing**; SI is "entirely reliant on the scheduled search performing as expected" and requires manual backfill.[^32][^33] SI data itself **cannot be report-accelerated.** The often-correct option is **"NOT ACCELERATED."**[^30]

## The Search Job Inspector
Job dropdown → **Inspect Job** (also via REST). Two sections: **Execution Costs** and **Search Job Properties.**[^34][^36]

**The benchmarking quartet:**[^34][^37]
- **`scanCount`** = events *scanned/read off disk*.
- **`eventCount`** = the subset of scanned events that match the search terms.
- **`resultCount`** = results after the *last* SPL command.
- **`runDuration`** = cumulative wall-clock across SH + indexers.

**The headline metric: scanCount, not resultCount.** "Search performance should not so much be measured using the `resultCount`/time rate but `scanCount`/time instead."[^37][^34] A healthy `scanCount`/sec is *roughly* ~10k–20k events/sec/indexer — but Splunk's own guidance stresses this is **environment-dependent** ("10k eps/indexer … well oiled machine, however you might also find that it is completely insufficient"); treat it as a heuristic, not an SLO.[^34][^37] Dense → `scanCount ≈ resultCount`; sparse → `scanCount >> resultCount` (you're scanning far more than you return — tighten the base filter).[^37][^35]

**Execution-cost components a query-tuner watches:**[^34][^35][^36]
- **`command.search.index`** = time crawling the **tsidx** files to locate which raw events to read. Invocations expose SmartStore: `command.search.index.bucketcache.hit`/`.miss` (S3 bucket downloads).
- **`command.search.rawdata`** = time reading the actual events from rawdata files.
- `command.search.kv` (search-time extraction), `command.search.lookups` (automatic-lookup cost), `command.search.typer` (eventtypes), `command.search.tags`, `command.search.filter` (drops non-matching events).
- **`dispatch.evaluate`** = parsing the search and building data structures, *including subsearch evaluation*, broken down per command (`dispatch.evaluate.join` etc.). Large values point at parse/subsearch overhead — a real 2,100-sec job was dominated by `dispatch.evaluate.join`.[^34]
- **`dispatch.fetch`** = SH time waiting for peers to return data.
- **`dispatch.stream.remote`** = bytes from indexers to SH; *lower = more efficient*; big variance between indexers signals **data imbalance**.
- `startup.handoff` = time shipping the search to indexers; large values vs indexer count → overloaded indexers/network.

**Where time goes → what to fix:** if `command.search.index`/`rawdata` dominate → tighten time range and base filter (too much disk scan); if `command.<transform>`/`dispatch.evaluate` dominate → the transforming/eval/regex/subsearch work is the bottleneck.[^38][^36] The **search.log** (visible only via the Job Inspector, never indexed) holds errors that don't reach the top bar (e.g., regex extraction errors) — grep it for `ERROR`/`WARN`.[^34]

## Streaming vs transforming command placement (map-reduce)
Splunk distributed search is **map-reduce**: the SH sends the search to each indexer (search peer), each indexer runs it over data it owns (**map**), and the SH aggregates (**reduce**).[^41][^5]
- **Distributable streaming** (`eval`, `rex`, `rename`, `fields`, `where`): per-event, order-independent, can run **on the indexers** in parallel — *but only until the pipeline moves to the SH*.[^40][^42]
- **Centralized (stateful) streaming** (`streamstats`, `head`): order matters → SH only.[^40][^42]
- **Transforming** (`stats`, `chart`, `timechart`, `top`) and most **dataset-processing** commands: **non-streaming** — need the *entire* event set before producing output.[^42][^43]

**Why transforming early kills parallelism:** a non-streaming command "requires the events from all of the indexers before [it] can operate," forcing the entire event set to the SH = "a lot of data movement and a loss of parallelism." Critically, **"once search processing moves to the search head, it can't be moved back to the indexer"** — so any distributable-streaming command placed *after* a transforming command is stranded single-threaded on the SH (e.g., in `… | stats … | rename …` the `rename` runs on the SH instead of in parallel).[^42][^44] **Rule:** put transforming/non-streaming commands **as late as possible**; keep the long streaming prefix on the indexers.[^42][^43] Lantern's pipeline ordering: base query → minimize (`fields`) → combine/summarize (`stats`, replacing `join`/`append`) → calculations (`eval`/`lookup`) → format (`stats`/`chart`/`timechart`).[^16] **Adjacent lever — parallel reduce / `redistribute`:** an intermediate-reduce phase (map-reduce-*reduce*) where a subset of indexers act as intermediate reducers cuts SH load for high-cardinality `BY`-clause reduces; it adds indexer load and can *slow* an already-overloaded cluster.[^41]

## Search modes (Fast / Smart / Verbose)
Three modes trade speed vs field richness:[^45][^48]
- **Fast** — prioritizes performance; **field discovery OFF**; no search-time extraction except fields you name in the SPL; returns only default + index-time fields.[^45][^46]
- **Verbose** — returns *all* fields: defaults + automatic + all user-defined extractions, *and* back-maps eventtypes and tags; slowest; best for short ranges / exploration.[^45][^48]
- **Smart** (default) — toggles: for **transforming** searches it behaves like Fast (straight to the table/viz); for **non-transforming** event searches it behaves like Verbose. "The best choice for most searches."[^45][^48]

Result *count* is identical across modes; the difference is extracted-field richness and cost.[^45] Levers: use the **`fields`** command to extract only what you need — meaningful because **reports always run in Smart mode**; disabling field discovery speeds the search (you lose discovered fields beyond `host`/`source`/`sourcetype`/`_time`).[^48][^7][^46]

## Other levers
**Search head vs indexer workload.** Keep work on indexers (parallel) and off the SH — the entire streaming/transforming discipline is about this. Dividing a 60-sec search across 6 indexers can finish in ~10 sec. Store apps on **fast local disk, not NFS** (a known bottleneck).[^6][^2][^16]

**Indexed extractions / indexed fields — double-edged.** Searching index-time fields is fast and unlocks `tstats`,[^13] but Splunk strongly cautions against custom indexed fields: each one **enlarges the searchable index** (slowing both indexing and search — "bigger is slower"), removes schema-on-read flexibility, and requires **re-indexing the entire dataset** to change. A case study saw indexed extractions blow an index to **~2.6–3× larger.**[^49][^51] **When index-time extraction *is* justified** (Splunk's criteria): you frequently search `NOT field=value`; the value frequently appears *outside* the field (e.g., small integers); or the value is only *part* of a token. Otherwise extract at search time.[^50][^23]

**Lookups vs joins for enrichment.** `join` is slow (effectively a Cartesian-style product in virtual memory) and carries the **50,000-row / 60-second** subsearch caps that **silently truncate** — "which often goes unnoticed and leads to incorrect answers."[^r1][^16] **Preferred:** replace `join`/`append` with a single base search over both datasets plus **`stats … BY <key>`** (using `values()`/`latest()`), keeping it one trip to the indexers; for reference/enrichment tables use a **lookup** (CSV or KV Store).[^16][^r2][^r4]

**KV Store vs CSV lookup.** CSV lookups are replicated to indexers (run *on* the indexers, in parallel) and perform well for small/static tables; KV Store collections live **on the search head.** KV Store wins for **large or frequently-updated** tables (per-record update, REST access, **field accelerations**).[^KV1][^KV3] **Gotcha:** a KV Store replicated to indexers ships as a **CSV** — so KV accelerations are *lost* on the indexer side; accelerations only help when the lookup runs on the SH. A 137k-row KV lookup added 165 sec until accelerations were correctly applied (→ ~20 sec). Use the lookup definition's `filter` attribute to prefilter large collections.[^KV-acc][^KV1]

**Scheduled-search concurrency & skipped searches.** Concurrency is a primary cause of slow/incomplete searches — watch for search concurrency above the indexer CPU count. Scheduled searches pile up at regular cadences (top of hour, midnight). Remedies: **flatten the load** by moving scheduled searches to quieter times; use **schedule windowing/skewing** to spread start times. A resource bottleneck causes the scheduler to **skip** a run.[^22][^34]

**Workload management (WLM).** A rule-based framework allocating **CPU + memory**:[^WLM-overview][^WLM-how]
- **Workload pools** = named CPU/memory slices; **workload rules** place searches into pools by predicate (`app`, `role`, `user`, `index`, `search_type`, `search_mode`, `search_time_range`), re-evaluated every 10 sec.[^WLM-rules]
- **Monitoring rules** act on running searches by `runtime` (abort, move pool, or surface a message when `runtime>5m`).[^WLM-rules]
- **Admission rules** (`admission_rules_enabled=1`) **prefilter searches before they start** — the canonical use is to *block* `index=*`/All-time rogue searches, or cap ad-hoc concurrency.[^WLM-admission]
- First matching rule wins (ordering matters); `workload_rules.conf` is **search-head-only**; `workload_pools.conf` must be on SH *and* indexers.[^WLM-examples][^WLM-dist]

**Metrics indexes + `mstats`.** A purpose-built TSDB: ~50% less storage than an events index and order-of-magnitude faster aggregations than `tstats`-on-events; `mstats` performs best with one or few `metric_name`s.[^25]

## Optimization checklist / anti-patterns
**Do:**
1. Set the **tightest time range** that answers the question; never default to All time.[^2][^3]
2. **Specify an index** (+ `sourcetype`/`source`/`host`); partition dissimilar data.[^1][^6]
3. **Filter most-restrictive terms first**, before the first pipe (the "tale of two searches").[^7]
4. Use **`tstats` on accelerated data models** (with `summariesonly=true` when staleness is acceptable).[^12][^18]
5. **Replace `join`/`append` with `stats … BY key`**; use **lookups** for enrichment.[^16][^r1]
6. Put **`stats`/`sort`/`dedup` last**; keep the streaming prefix on indexers.[^42][^43]
7. Use **Fast mode / the `fields` command** to cut field-discovery cost, especially in scheduled reports.[^45][^48]
8. **Skew/window** scheduled searches; use **WLM admission rules** to block rogue searches.[^22][^WLM-admission]
9. Optimize by **`scanCount`**, reading `command.search.index`/`rawdata` vs `dispatch.evaluate`/transform costs in the Job Inspector.[^37][^34]
10. For metrics, use a **metrics index + `mstats`**.[^25]

**Don't (anti-patterns):** leading/mid-word **wildcards**[^9][^10]; frequent **`NOT`/`!=`** on near-constant fields[^23]; large **OR-lists** and **subsearches** that silently truncate at 50k/60s[^8][^16]; **transforming commands early**[^42]; **over-indexing fields**[^49]; treating **DMA/report-acceleration as free**.[^27][^20]

## Disconfirming findings / gotchas
- **`tstats` can be slower than raw `stats`** on high-cardinality `WHERE` fields, and is weak at returning very many rows.[^12][^13]
- **DMA/report-acceleration disk cost is real and unbounded by default** — both default to *unlimited* summary disk; you must set size-based retention or risk freezing index buckets early and **losing data.** Rough sizing ~3.4× average daily indexed volume for one year of DMA; DMA can create *many* `.tsidx` files.[^DMA-disk][^DMA-caution][^30]
- **DMA `summariesonly=true` silently misses recent/lagged data** until acceleration propagates over the backfill window; a summarization search that exceeds **Max Summarization Search Time** is stopped mid-interval with **no auto-retry → gaps** that `summariesonly=true` then misses.[^DMA-backfill][^DMA-maxtime]
- **Summary indexing gaps & overlaps:** SI has no history before its start date (needs backfill); if the populating search runs **longer than its interval**, the next run won't start → **gaps**; misaligned schedules cause **duplicates or gaps** (use non-overlapping snapped lagged windows). **Backfill with `fill_summary_index.py`, not `overlap`** (which only *finds* gaps and can't read `_raw`); `-dedup true` *skips* timespans that already have data, so it won't repair *partial* days. Use **durable search** on SI populators to auto-backfill failed runs.[^SI-gaps][^overlap][^SI-durable]
- **Report-acceleration summaries expire with raw data** (they point at hot/warm buckets) — only SI can outlive raw retention.[^32]
- **KV Store acceleration is lost on indexers** (shipped as CSV); large KV lookups without accelerated fields can add minutes-to-hours.[^KV-acc]

## Adjacent / frontier concepts
Metrics indexes + `mstats`/`mcatalog`/`mcollect`/`mrollup`; parallel reduce / `redistribute`; SmartStore cache effects on search (`bucketcache.hit/miss`); `tsidxWritingLevel` and tsidx reduction; `tstats`/`sistats`/`mstats` `chunk_size` memory tuning; durable search; search head pooling (deprecated) → search head clustering.

## References
[^1]: About search optimization — Splunk Docs (Search Manual) — https://help.splunk.com/en/splunk-enterprise/search/search-manual/10.2/optimize-searches/about-search-optimization — core principles, index/bucket retrieval, worked example.
[^2]: Quick tips for optimization — Splunk Docs (Search Manual) — https://help.splunk.com/en/splunk-cloud-platform/search/search-manual/10.3.2512/optimizing-searches/quick-tips-for-optimization — narrow time window, filter early, specificity, OR overuse, NFS, late non-streaming.
[^3]: About searching with time — Splunk Docs (Search Manual 10.4) — https://help.splunk.com/en/splunk-enterprise/search/search-manual/10.4/specify-time-ranges/about-searching-with-time — `_time`/epoch, narrow ranges, defaults, timezone-by-indexer.
[^4]: Clara-fication: Search Best Practices — Splunk blog (Merriman) — https://www.splunk.com/en_us/blog/customers/splunk-clara-fication-search-best-practices.html — All-time cost, short prototyping ranges, scanCount as metric.
[^5]: Specify time modifiers — Splunk Docs (Search Manual) — https://help.splunk.com/en/splunk-cloud-platform/search/search-manual/10.4.2604/specify-time-ranges/specify-time-modifiers-in-your-search — relative time, latest=now requirement, future-event gotcha.
[^6]: Optimizing search — Splunk Lantern — https://lantern.splunk.com/Platform_Data_Management/Transform_Data/Optimizing_search — index-first, buckets-eliminated metric, horizontal scaling, tsidxWritingLevel.
[^7]: Write better searches — Splunk Docs (Search Manual) — https://help.splunk.com/en/splunk-cloud-platform/search/search-manual/10.4.2604/optimizing-searches/write-better-searches — tale of two searches, non-streaming-on-SH, field discovery off, dense/sparse, OR/subsearch complexity.
[^8]: Write better searches (dense/sparse, OR/subsearch) — Splunk Docs — https://help.splunk.com/en/splunk-cloud-platform/search/search-manual/10.4.2604/optimizing-searches/write-better-searches — subsearch/OR complexity.
[^9]: Wildcards — Splunk Docs (Search Manual 10.4) — https://help.splunk.com/en/splunk-enterprise/search/search-manual/10.4/search-primer/wildcards — bare `*` cost, prefix-wildcard cost, mid-word segmentation failures.
[^10]: Search performance of raw without wildcards — Splunk Community — https://community.splunk.com/t5/Splunk-Search/Search-performance-of-raw-without-wildcards/m-p/603037 — reverse-word index, bloom filters, leading-wildcard full scan.
[^12]: tstats — Splunk Docs (SPL Search Reference 10.4) — https://help.splunk.com/en/splunk-enterprise/search/spl-search-reference/10.4/search-commands/tstats — tstats on indexed fields/DMA, faster than stats, summariesonly, high-cardinality-WHERE-slower note, chunk_size.
[^13]: What is tstats and why is it faster than stats — Splunk Community — https://community.splunk.com/t5/Splunk-Search/What-is-tstats-and-why-is-so-much-faster-than-stats/m-p/116960 — tsidx metadata, index-time fields only, columnar analogy, weak at many results.
[^14]: What is tstats (columnar reading) — Splunk Community — https://community.splunk.com/t5/Splunk-Search/What-is-tstats-and-why-is-so-much-faster-than-stats/m-p/116960 — KEY::VALUE columnar reading.
[^15]: Accelerate data models — Splunk Docs (Knowledge Mgmt Manual 10.4) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/use-data-summaries-to-accelerate-searches/accelerate-data-models — HPAS, .tsidx summaries on indexers, partial acceleration, Size on Disk, unlimited default, Max Summarization Time, Backfill Range.
[^16]: Writing better queries in SPL — Splunk Lantern — https://lantern.splunk.com/Platform_Data_Management/Transform_Data/Writing_better_queries_in_Splunk_Search_Processing_Language — avoid join/append, stats-instead-of-join, pipeline ordering, combine evals.
[^17]: Rules/requirements for tstats prestats — Splunk Community — https://community.splunk.com/t5/Splunk-Search/What-exactly-are-the-rules-requirements-for-using-quot-tstats/m-p/319801 — prestats=t needs tsidx-backed data; summariesonly behavior.
[^18]: Accelerate data models (summariesonly detail) — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/use-data-summaries-to-accelerate-searches/accelerate-data-models — summariesonly true/false behavior.
[^19]: Is datamodel command as fast as tstats — Splunk Community — https://community.splunk.com/t5/Splunk-Search/Is-datamodel-command-as-fast-as-tstats-command-on-accelerated/m-p/759071 — datamodel/from read raw events; tstats first command.
[^20]: Splunk Accelerated Data Models Part 3 — Helge Klein (blog/conf2015) — https://helgeklein.com/blog/splunk-accelerated-data-models-part-3/ — HPAS only for Pivot/tstats, regular search + datamodel NOT accelerated, 5-min build/30-min purge.
[^21]: How search types affect performance — Splunk Docs (Deployment/Capacity Manual) — https://help.splunk.com/en/splunk-enterprise/get-started/deployment-capacity-manual/10.4/hardware-capacity-planning/how-search-types-affect-splunk-enterprise-performance — dense/sparse/super-sparse/rare, bloom filters, CPU vs I/O bound.
[^22]: Performance Insights troubleshooting — Splunk Lantern — https://lantern.splunk.com/Manage_Performance_and_Health/Using_the_Performance_Insights_for_Splunk_app — search concurrency > CPU, schedule windowing/skewing, subsearch/join limit truncation.
[^23]: Should I use an index-time field extraction — Splunk Community (Müller) — https://community.splunk.com/t5/Splunk-Search/Should-I-use-an-index-time-field-extraction/m-p/220352 — when to index (NOT field=value, value outside field, partial token), retrieve-then-discard.
[^25]: mstats / metric index performance — Splunk Docs + .conf2019 FN2268 — https://help.splunk.com/en/splunk-enterprise/search/spl-search-reference/10.4/search-commands/mstats — mstats vs tstats latency, ~50% less storage, single-metric best perf.
[^26]: Overview of summary-based search acceleration — Splunk Docs (Knowledge Mgmt Manual 10.4) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/use-data-summaries-to-accelerate-searches/overview-of-summary-based-search-acceleration — the three methods + storage + when to use each.
[^27]: Accelerate reports / Manage report acceleration — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.2/use-data-summaries-to-accelerate-searches/manage-report-acceleration — RA builds summary, 10-min refresh, qualification rules, disk impact.
[^28]: Different data acceleration methods — Splunk Community — https://community.splunk.com/t5/Security/Different-data-acceleration-methods/m-p/343319 — RA ~2-5× (single-source estimate), SI "astronomical," DMA not guaranteed, "NOT ACCELERATED" often right.
[^30]: Different acceleration methods (OLAP-cube analogy / NOT ACCELERATED) — Splunk Community — https://community.splunk.com/t5/Security/Different-data-acceleration-methods/m-p/343319 — SI as OLAP-cube cousin.
[^31]: sistats — Splunk Docs (SPL Search Reference 10.2) — https://help.splunk.com/en/splunk-enterprise/spl-search-reference/10.2/search-commands/sistats — sistats populates SI then report with stats.
[^32]: Should we move from SI to report acceleration — Splunk Community — https://community.splunk.com/t5/Reporting/Should-we-move-from-summary-indexes-to-report-accelerations/m-p/356454 — SI when nonstreamable-before-transform or raw rolls faster; RA expires with raw; SI multi-layer.
[^33]: Summary Index vs Report Acceleration — Splunk Community — https://community.splunk.com/t5/Getting-Data-In/Summary-Index-vs-Report-Acceleration/m-p/161909 — RA self-repairing, SI manual backfill/fragile, storage locations.
[^34]: Clara-fication: Job Inspector — Splunk blog (Merriman) — https://www.splunk.com/en_us/blog/tips-and-tricks/splunk-clara-fication-job-inspector.html — execution-cost & job-property walkthrough, scanCount/eventCount/resultCount/runDuration, bucketcache, dispatch.evaluate.join, search.log.
[^35]: Troubleshooting Splunk Search Performance by Job Inspector — djangocas.dev (blog) — https://djangocas.dev/blog/splunk/splunk-timechart-examples/ — command.search.index/filter/rawdata, scanCount vs resultCount.
[^36]: View search job properties — Splunk Docs (Search Manual 10.4) — https://help.splunk.com/en/splunk-enterprise/search/search-manual/10.4/manage-jobs/view-search-job-properties — command.search.index/rawdata, dispatch.evaluate, scanCount/eventCount.
[^37]: View search job properties (scanCount focus) — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/search/search-manual/10.4/manage-jobs/view-search-job-properties — execution-cost definitions, dense/sparse note.
[^38]: Splunk Search Job Inspector Analysis — ThinkCloudly (blog) — https://thinkcloudly.com/blog/search-job-inspector-analysis/ — pipeline stages, event-retrieval-dominant signal, transformation-heavy command cost.
[^40]: Command types — Splunk Docs (SPL Search Reference) — https://help.splunk.com/en/splunk-enterprise/search/spl-search-reference/10.4/quick-reference/command-types — six command types, indexer-vs-SH placement.
[^41]: Types of commands + Parallel reduce overview — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/distributed-search/10.4/manage-parallel-reduce-search-processing/supported-commands-for-parallel-reduce-search-processing — non-streaming forces data to SH, map-reduce-reduce, forward-SH-data best practice.
[^42]: Types of commands — Splunk Docs (Search Manual 10.4) — https://help.splunk.com/en/splunk-enterprise/search/search-overview/types-of-commands — non-streaming definition, "can't move back to indexer", placement.
[^43]: Types of commands (non-streaming) — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/search/search-overview/types-of-commands — non-streaming requires full dataset.
[^44]: Learn SPL Command Types — Splunk blog (Kuranami) — https://www.splunk.com/en_us/blog/tips-and-tricks/learn-spl-command-types-efficient-search-execution-order-and-how-to-investigate-them.html — distributable streaming on indexers, stats-then-rename single-threaded example.
[^45]: Search modes — Splunk Docs (Search Manual 10.4) — https://help.splunk.com/en/splunk-enterprise/search/search-manual/10.4/use-the-search-app/search-modes — Fast/Smart/Verbose, field discovery off in Fast, identical result counts, Smart toggles.
[^46]: Extracted field in fast mode — Splunk Community — https://community.splunk.com/t5/Splunk-Search/Extracted-field-in-fast-mode/m-p/754841 — Verbose back-maps eventtypes/tags, `fields` lever, reports run in Smart.
[^48]: Difference between verbose and fast mode — Splunk Community — https://community.splunk.com/t5/Splunk-Dev/What-is-the-difference-of-results-between-verbose-and-fast-mode/m-p/415877 — Smart=Fast for transforming / Verbose for events.
[^49]: When Splunk extracts fields + Create indexed field extraction — Splunk Docs — https://help.splunk.com/en/data-management/get-data-in/get-data-into-splunk-enterprise/10.4/configure-indexed-field-extraction/create-custom-fields-at-index-time — indexed fields enlarge index, re-index requirement, when justified.
[^50]: Create custom fields at index time — Splunk Docs — https://help.splunk.com/en/data-management/get-data-in/get-data-into-splunk-enterprise/10.4/configure-indexed-field-extraction/create-custom-fields-at-index-time — index-time-extraction criteria, drop duplicated value.
[^51]: Indexed Extractions vs Search-Time Extractions case study — Hurricane Labs (blog) — https://hurricanelabs.com/splunk-tutorials/the-indexed-extractions-vs-search-time-extractions-splunk-case-study/ — index ~2.6-3× raw size, worse CPU/storage.
[^r1]: join — Splunk Docs (SPL Search Reference 10.4) — https://help.splunk.com/en/splunk-enterprise/spl-search-reference/10.4/search-commands/join — 50k/60s limits, max=1.
[^r2]: Why is join slow — Splunk Community — https://community.splunk.com/t5/Splunk-Search/Why-is-join-slow/td-p/434152 — Cartesian product, not a RDB.
[^r4]: Alternative method to using Join — Splunk Community — https://community.splunk.com/t5/Splunk-Search/Alternative-method-to-using-Join-command/m-p/532978 — stats-BY replacement pattern.
[^KV1]: Configure KV Store lookups — Splunk Docs (Knowledge Mgmt Manual 10.4) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/use-the-configuration-files-to-configure-lookups/configure-kv-store-lookups — filter attribute, KV-on-SH.
[^KV3]: Define a KV Store lookup in Splunk Web — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/use-lookups-in-splunk-web/define-a-kv-store-lookup-in-splunk-web — CSV-on-indexers vs KV-on-SH, accel fields.
[^KV-acc]: KVStore acceleration not working as expected — Splunk Community — https://community.splunk.com/t5/Splunk-Dev/KVStore-acceleration-not-working-as-expected/m-p/680438 — 165s→20s; KV→CSV on indexers loses accel.
[^SI-gaps]: Manage summary index gaps — Splunk Docs (Knowledge Mgmt Manual 10.4) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/use-data-summaries-to-accelerate-searches/manage-summary-index-gaps — no pre-start history, long-run-than-interval gap, fill_summary_index.py, -dedup.
[^overlap]: overlap — Splunk Docs (SPL Search Reference) — https://help.splunk.com/en/splunk-cloud-platform/spl-search-reference/10.1.2507/search-commands/overlap — don't use overlap to backfill; can't read _raw.
[^SI-durable]: Make scheduled reports durable — Splunk Docs (Reporting Manual) — https://help.splunk.com/en/splunk-cloud-platform/create-dashboards-and-reports/reporting-manual/10.3.2512/report-management/make-scheduled-reports-durable-to-prevent-event-loss — durable search backfills failed runs.
[^DMA-disk]: Accelerate data models (Size on Disk, unlimited default) — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/use-data-summaries-to-accelerate-searches/accelerate-data-models — unlimited summary disk default, avoid overuse.
[^DMA-caution]: Caution on Retention Impact of Accelerated Data Model — Splunk Community — https://community.splunk.com/t5/Reporting/Caution-on-Retention-Impact-of-Accelerated-Data-Model-and-Report/m-p/155808 — summaries age with hot/warm buckets; separate volumes.
[^DMA-backfill]: Datamodel Tuning - Backfill question — Splunk Community — https://community.splunk.com/t5/Splunk-Enterprise-Security/Datamodel-Tuning-Backfill-question/m-p/759213 — backfill range; summariesonly misses until propagated.
[^DMA-maxtime]: Datamodel Acceleration Reaching Max Summarization — Splunk Community — https://community.splunk.com/t5/Knowledge-Management/Datamodel-Acceleration-Consistently-Reaching-Max-Summarization/m-p/747888 — timeout → no retry → gaps → summariesonly misses.
[^WLM-overview]: About workload management — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/manage-workloads/9.3/workload-management-overview/about-workload-management — WLM purpose.
[^WLM-how]: How workload management works — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/manage-workloads/10.4/workload-management-overview/how-workload-management-works — pools/rules mechanics.
[^WLM-rules]: Configure workload rules — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/manage-workloads/10.4/configure-workload-management/configure-workload-rules — predicates, 10-sec re-eval, monitoring runtime actions.
[^WLM-admission]: Configure admission rules to prefilter searches — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/manage-workloads/10.4/configure-workload-management/configure-admission-rules-to-prefilter-searches — block alltime/index=*, adhoc_search_percentage.
[^WLM-examples]: Workload management examples — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/manage-workloads/10.4/workload-management-examples/workload-management-examples — rule ordering, first match wins.
[^WLM-dist]: Configure WLM on distributed deployments — Splunk Docs — https://help.splunk.com/en/splunk-enterprise/administer/manage-workloads/10.4/configure-workload-management/configure-workload-management-on-distributed-deployments — rules SH-only, pools SH+indexers.
