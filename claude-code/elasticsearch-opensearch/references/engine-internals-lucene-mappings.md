# Engine Internals — Lucene, Indexing Model, Mappings & Analyzers

> Provenance: reference under the `elasticsearch-opensearch` skill. Researched 2026-06-18 via `/dr` deep research (Elastic docs/blogs, OpenSearch docs, Apache Lucene API docs, the Elasticsearch Definitive Guide). Every mechanic below is **identical in Elasticsearch and OpenSearch** unless flagged — both wrap the same Apache Lucene core (OpenSearch forked ES 7.10.2). Version-specific numbers are stamped **verified-as-of 2026-06-18**.

## Contents

- Apache Lucene as the engine (segments, immutability, merging)
- Refresh vs flush vs commit + the translog
- Inverted index, doc_values, and BKD trees
- The document & index model (and the removal of mapping types)
- Mappings & analyzers (dynamic vs explicit, text vs keyword, the analysis chain)

## Apache Lucene as the engine

- Elasticsearch/OpenSearch are **distributed coordination layers on top of Apache Lucene**. Lucene is the actual search engine (inverted index + storage); the surrounding system adds sharding, replication, a JSON-over-HTTP API, node discovery, and the query DSL.[^4][^5][^6][^15]
- **Each shard is exactly one Apache Lucene index**, holding a subset of the index's documents and independently handling indexing and search.[^1][^2][^4][^16]
- A Lucene index is **composed of segments** — each segment is a complete, self-contained searchable mini-index over a subset of documents. A shard = a set of segments.[^2][^3][^5][^16]
- **Segments are immutable.** Once written, postings, term dictionaries, doc values, and stored fields are never modified in place. This is what enables lock-free concurrent reads, trivial replication, and crash safety (an incomplete segment is simply discarded).[^3][^4][^5][^6]
- **Updates/deletes don't mutate segments.** A delete marks the doc in a per-segment `liveDocs` bitmap (a tombstone / soft delete); an update is an atomic insert-of-new + delete-of-old. The bytes for deleted docs are *not* reclaimed until merge.[^3][^4][^6][^17]
- **Segment merging** combines smaller segments into larger ones in the background (Lucene's `TieredMergePolicy` by default) to keep search efficient and to **physically purge deleted documents** — deleted docs and "ghost terms" (terms occurring only in deleted docs) are reclaimed only at merge time.[^3][^5][^17][^18] *Operational consequence:* "delete indices, not documents" — deleted docs linger and cost resources until a merge runs; `force_merge` on read-only indices reclaims them deliberately.
- A Lucene index caps at **2,147,483,519 (2³¹ − 129) documents** (`docs.count` + `docs.deleted`); each doc has a 32-bit docid.[^1][^14] *(Exact figure is single-source; the 2³¹ order of magnitude is well established.)*

## Refresh vs flush vs commit + the translog

This is the most commonly confused area; the distinction is **searchability vs durability**.

**Write path:** an indexing op is written to (a) an in-memory buffer and (b) the **translog** (transaction log / write-ahead log) before the request is acknowledged.[^7][^8][^9]

- **Refresh** takes the in-memory buffer and writes it as a **new segment that becomes searchable** — the segment lands in the filesystem cache and is **not** yet fsync'd to durable storage. Mechanically it is a Lucene `reopen` of the reader/searcher over the new segment. This is what makes search **near-real-time (NRT)**: separating searchability from durability is precisely the ES/OpenSearch innovation over classical Lucene, where docs were searchable only after a full commit.[^4][^7][^9][^11]
  - **Default refresh interval = 1 second**, but **only on indices that received ≥1 search request in the last 30 seconds** ("search idle" optimization, introduced in **Elasticsearch 7.0**; OpenSearch inherits it). If no explicit `refresh_interval` is set and the shard is search-idle, scheduled refreshes are skipped until a search arrives. **verified-as-of 2026-06-18.**[^12][^13][^19]
- **Flush** performs a **Lucene commit** (fsyncs the in-memory segments to disk so they are durable) **and starts a fresh translog generation** (truncating the old one). Flushes happen automatically; a manual `_flush` API exists but is rarely needed.[^7][^8][^9][^10]
- **Lucene commit** is the expensive fsync-of-segments operation that makes data durable at the Lucene level — too costly to run per-operation, which is the entire reason the translog exists.[^7][^8][^10]
- **Translog** is a per-shard write-ahead log giving **durability between Lucene commits**. Operations are recorded there after processing but before acknowledgment; on crash, acknowledged ops not yet in the last commit are **replayed from the translog** during shard recovery.[^7][^8][^9][^10]
- **Orthogonality (key mental model):** a flush by itself does **not** make data searchable, and a refresh by itself does **not** make data durable. Data becomes searchable only after a *refresh*; durable only after a *flush*.[^9][^7][^11]

**Translog defaults (verified-as-of 2026-06-18):**
- `index.translog.durability` defaults to **`request`** — fsync + commit the translog on *every* request before acknowledging (on primary and every allocated replica). Set to `async` to fsync only every `sync_interval`, risking loss of the last interval's ops on crash.[^8][^20]
- `index.translog.sync_interval` defaults to **`5s`**.[^20]
- `index.translog.flush_threshold_size` defaults to **`10gb`** in current Elasticsearch — **but this default changed from `512mb` to `10gb` in Elasticsearch 8.8.** Many older docs/blogs still say 512mb (correct for pre-8.8). **Flag this when you see "512mb."**[^8][^20]
- *(OpenSearch's specific current numeric defaults were not independently pinned this pass; the mechanisms are identical, but do not assume OpenSearch's `flush_threshold_size`/`search.idle.after` numbers match ES without checking.)*

## Inverted index, doc_values, and BKD trees

Three on-disk structures power everything; mapping a field type is really choosing which of these gets built.

**Inverted index (full-text search):**
- Lucene's core is the **inverted index, built from postings.** Conceptually it is a **term dictionary → postings list** map: given a term (token), return the ordered list of documents containing it. It maps **content → documents** (the inverse of a "forward" index), enabling fast full-text search.[^14][^15][^21][^22]
- **Postings cannot retrieve terms given a document** — only the reverse — which is exactly why doc_values exist.[^14][^15]

**doc_values (the columnar inverse — sorting & aggregations):**
- `doc_values` is an **on-disk, column-oriented, per-document structure built at index time** — the inverse of the inverted index. It stores the same values as `_source` but in a columnar layout optimized for **sorting, aggregations, and field access in scripts**.[^15][^16]
- **Enabled by default on nearly all field types EXCEPT `text`** (and analyzed string fields).[^15][^16]
- doc_values are **not** good for searching/filtering (evaluating a range over them means a linear scan), so: sorting/aggregations use doc_values; searching uses the inverted index or points.[^15][^16]

**BKD trees / "points" (numeric, date, geo):**
- **Numeric, date, and geo_point fields are indexed with a BKD tree (block k-d tree)**, not the inverted index. Lucene calls these "points" (`PointValues`), optimized for **range, distance, nearest-neighbor, and point-in-polygon** queries; the tree partitions space so a range query can skip non-competitive cells. Stored off-heap.[^15][^24][^25]
- For numeric fields that are *both* indexed (points) and have doc_values (the default), Elasticsearch wraps queries in `IndexOrDocValuesQuery`, transparently choosing points-to-lead vs doc_values-to-verify based on the access pattern.[^15]

## The document & index model

- **Documents are JSON**, the basic unit of information, identified by **`_id`**; the original JSON is retained in **`_source`**; each doc also gets an internal 32-bit Lucene docid.[^14][^15][^16]
- An **index is a logical collection of (related) documents**, divided into shards distributed across nodes for horizontal scaling and redundancy.[^2][^16]

**Removal of mapping types (version timeline — verified-as-of 2026-06-18):**
Old Elasticsearch allowed **multiple "mapping types" per index** for multi-tenancy, but this was incompatible with Lucene (same-named fields across types share one Lucene field) and "caused more problems than it solved," so it was removed across four major versions:[^26][^27][^28][^31]
- **5.0** — began enforcing compatible mappings for same-named fields across types.
- **6.0** — new indices restricted to a **single mapping type**; `_default_` deprecated; `_uid` (type+id) replaced by `_id`.
- **6.7** — `include_type_name` query param introduced (default `true`) to ease transition.
- **7.0** — type-accepting APIs deprecated; typeless APIs introduced; `include_type_name` defaults to `false`; the dummy/implicit type name became **`_doc`**.
- **8.0** — APIs that accept types **removed**.

Net effect: **effectively one implicit type, `_doc`, per index** (`PUT {index}/_doc/{id}`, `POST {index}/_doc`). OpenSearch (forked from 7.10.2) inherits the post-removal typeless model.[^26][^27][^31][^16]

## Mappings & analyzers

**Mapping = the schema** (field names → data types and how each is stored/indexed).[^16]

- **Dynamic mapping:** ES/OpenSearch **auto-detects field data types from the first-seen values** and adds fields on the fly — great for getting started, but only a limited set of types is auto-detected (e.g., `geo_point`/`geo_shape` need explicit mapping). Settable to `true`, `runtime` (schema-on-read), `false` (ignore new fields), or `strict` (reject docs with unknown fields).[^32][^33][^34]
- **Explicit mapping:** you define types yourself; recommended for production for full control. The two can be combined. **Mappings cannot be reduced or have a field's type changed once created — that requires a reindex.**[^32][^34]

**text vs keyword (the central field-type decision):**
- **`text` fields are analyzed** — tokenized into terms before indexing, for full-text search. Searchable by default but **NOT used for sorting and seldom for aggregations**; doc_values/fielddata are **disabled** on them by default.[^15][^16]
- **`keyword` fields are stored verbatim** (not analyzed) for **exact match, sorting, aggregations, and term-level queries** — faster, lower storage.[^35]
- History: pre-5.0 there was a single `string` type with `analyzed`/`not_analyzed`; **Elasticsearch 5.0 split it into `text` and `keyword`** to remove that confusion.[^35]
- **Multi-field pattern:** the same string is commonly mapped as **`text` for search WITH a `.keyword` sub-field** (e.g., `city` + `city.keyword`) for aggs/sorting. Dynamic mapping of a string auto-creates exactly this — a `text` field with a `.keyword` sub-field (`ignore_above: 256`). Multi-fields are independent of the parent and don't change `_source`.[^35]
- **The fielddata gotcha (very common error):** sorting/aggregating/scripting on a `text` field fails with *"Fielddata is disabled on text fields by default… Alternatively use a keyword field instead."* `fielddata` is a query-time, **in-heap** structure built by uninverting the inverted index per segment — on high-cardinality text it can exhaust heap / OOM. **The fix is to aggregate on the `.keyword` sub-field, NOT to set `fielddata: true`.**[^36][^37][^38]

**The analysis chain:**
- An analyzer is exactly three building blocks **in order: character filters (0+) → tokenizer (exactly 1) → token filters (0+).** Char filters transform the raw character stream (strip HTML, map digits); the tokenizer breaks it into tokens recording positions/offsets; token filters add/remove/change tokens (lowercase, stop-words, synonyms, stemming) but **may not change positions or offsets.**[^39][^40][^41]
- **The `standard` analyzer** (the default) = Standard Tokenizer + Lower Case filter (+ Stop filter, disabled by default). It tokenizes on word boundaries per the Unicode Text Segmentation algorithm (UAX #29), removes most punctuation, and lowercases — i.e., "lowercases + splits on word boundaries" is exactly right.[^40][^41]
- **Index-time vs search-time analysis:** the field value is analyzed at index time; the **query string is analyzed at search time**; the **same analyzer should normally be used at both** so tokens match. A distinct search analyzer is the less-common exception (e.g., to stop an edge-ngram index analyzer from over-matching at query time).[^42][^41]
- OpenSearch uses the same char-filter → tokenizer → token-filter model and the same `standard` default.[^16]

## Where ES and OpenSearch are identical vs differ

- **Identical (shared Lucene core):** segments + immutability + merging, inverted index, doc_values, BKD/points, refresh→searchable / flush→durable / commit / translog, the document/`_id`/`_source` model, mappings (dynamic/explicit), text-vs-keyword, multi-fields, the fielddata behavior, the analysis chain.[^15][^16]
- **Potential drift to flag:** exact tunable *defaults* can diverge post-fork. ES's `flush_threshold_size = 10gb` (8.8+) is confirmed; OpenSearch's specific numeric defaults were not pinned here — the *mechanisms* are the same, the *numbers* may not be.

## References

1. Elastic — Elasticsearch shards & replicas guide: https://www.elastic.co/search-labs/blog/elasticsearch-shards-and-replicas-guide — shard = Lucene index; doc cap; merging rationale. blog
2. Elastic — Index basics: https://www.elastic.co/docs/manage-data/data-store/index-basics — shard is a self-contained Lucene index; immutable segments. docs
3. Apache Lucene — index package summary (10.3.1): https://lucene.apache.org/core/10_3_1/core/org/apache/lucene/index/package-summary.html — segments, immutability, deletes/updates → new segments, merging reclaims space. docs
4. systeminternals.dev — Elasticsearch Internals: https://systeminternals.dev/elasticsearch/ — Lucene-per-shard, segment immutability, liveDocs tombstones, translog. blog
5. letsbuildsolutions.com — How Elasticsearch Works Internally: https://letsbuildsolutions.com/blog/system-design/how-elasticsearch-works-internally — Lucene engine, immutable segments, TieredMergePolicy, refresh→filesystem cache. blog
6. abstractalgorithms.dev — How Apache Lucene Works: https://abstractalgorithms.dev/how-lucene-works — inverted index, immutable segments, deletes via tombstone, ES = Lucene distributed. blog
7. Opster — Flush, Translog & Refresh: https://opster.com/guides/elasticsearch/glossary/elasticsearch-flush-translog-and-refresh/ — flush = Lucene commit; refresh adds in-memory segment; translog for durability (cites pre-8.8 512mb). blog
8. Elastic — Translog (reference): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/index-modules-translog.html — commit cost, WAL, flush = commit + new translog gen, replay on crash, durability/sync defaults. docs
9. Elastic Discuss — What happens during refresh: https://discuss.elastic.co/t/what-happens-during-elasticsearch-refresh/356780 — experimental confirmation refresh (not flush) makes data searchable. forum
10. Elastic — Translog settings: https://www.elastic.co/docs/reference/elasticsearch/index-settings/translog — flush_threshold_size triggers a commit point. docs
11. StackOverflow — Refresh vs flush: https://stackoverflow.com/questions/19963406/refresh-vs-flush — refresh = Lucene reopen (NRT, no fsync); flush = commit + empties translog. forum
12. Elastic — Near real-time search (7.10): https://www.elastic.co/guide/en/elasticsearch/reference/7.10/near-real-time.html — refresh writes+opens new segment; 1s only if searched in last 30s. docs
13. Elastic — 7.0 release highlights (search-idle, PR #27500): https://www.elastic.co/guide/en/elasticsearch/reference/7.4/release-highlights-7.0.0.html — 7.0 introduced search-idle (30s). docs
14. Apache Lucene — index package summary (docids; postings can't get terms-from-doc): https://lucene.apache.org/core/10_3_1/core/org/apache/lucene/index/package-summary.html. docs
15. Elastic — doc_values reference + range-query planning + BKD blogs: https://www.elastic.co/docs/reference/elasticsearch/mapping-reference/doc-values ; https://www.elastic.co/blog/better-query-planning-for-range-queries-in-elasticsearch — doc_values columnar/default-on-except-text; points vs doc_values; IndexOrDocValuesQuery. docs/blog
16. OpenSearch — Concepts, Doc values, Index parameter: https://docs.opensearch.org/latest/getting-started/concepts/ ; https://docs.opensearch.org/latest/mappings/mapping-parameters/doc-values/ — Lucene engine, immutable segments, inverted index, doc_values inverse, text has no doc_values. docs
17. Elastic — Lucene's Handling of Deleted Documents: https://www.elastic.co/blog/lucenes-handling-of-deleted-documents — deleted bytes/ghost terms reclaimed only at merge. blog (vendor)
18. Apache Lucene — TieredMergePolicy (10.4): https://lucene.apache.org/core/10_4_0/core/org/apache/lucene/index/TieredMergePolicy.html — merge budget, delete-reclaim scoring. docs
19. Elastic — Refresh API + Tune for indexing speed (7.x): https://www.elastic.co/guide/en/elasticsearch/reference/7.5/indices-refresh.html — 1s only if searched in last 30s; refresh_interval=-1 for bulk. docs
20. Elastic GitHub — issue #101472 (flush_threshold_size 512mb → 10gb in 8.8): https://github.com/elastic/elasticsearch/issues/101472 — durability=request default, sync_interval=5s. docs/issue
21. Apache Lucene — Lucene103PostingsFormat: https://lucene.apache.org/core/10_3_2/core/org/apache/lucene/codecs/lucene103/Lucene103PostingsFormat.html — postings files .tim/.tip/.doc/.pos/.pay; term dictionary contents. docs
22. northcoder.com — Basic Lucene Index Structures: https://northcoder.com/post/some-basic-lucene-index-structures/ — inverted (content→docs) vs forward; TermsEnum/PostingsEnum. blog
24. Apache Lucene wiki — Point Value data structure + BKDWriter (10.4): https://lucene.apache.org/core/10_4_0/core/org/apache/lucene/util/bkd/BKDWriter.html — BKD = block KD-tree, off-heap, cells with min/max. docs
25. Apache Lucene — PointValues (9.4.1): https://lucene.apache.org/core/9_4_1/core/org/apache/lucene/index/PointValues.html — points indexed with KD-trees (not inverted index) for range/distance/NN/polygon. docs
26. Elastic — Moving from types to typeless APIs in 7.0: https://www.elastic.co/blog/moving-from-types-to-typeless-apis-in-elasticsearch-7-0 — 5.0/6.0/7.0/8.0 timeline; `_doc`; types incompatible with Lucene. blog (vendor)
27. Elastic — Removal of mapping types in 6.0: https://www.elastic.co/blog/removal-of-mapping-types-elasticsearch — 6.0 single type/index; `_uid`→`_id`. blog (vendor)
28. elastic/elasticsearch — removal_of_types.asciidoc (6.5): https://github.com/elastic/elasticsearch/blob/6.5/docs/reference/mapping/removal_of_types.asciidoc — per-version timeline in official docs source. docs
31. elastic/elasticsearch — issue #35190 (7.0 types deprecation): https://github.com/elastic/elasticsearch/issues/35190 — include_type_name lifecycle, `_doc` dummy. docs/issue
32. Elastic — Explicit mapping: https://www.elastic.co/docs/manage-data/data-store/mapping/explicit-mapping — explicit vs dynamic; can't change a field type → reindex. docs
33. Elastic — Dynamic field mapping (8.19): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/dynamic-field-mapping.html — auto-detect; dynamic true/runtime/false/strict. docs
34. Elastic — Mapping (8.19 + hub): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/mapping.html — dynamic auto-detects; explicit for control; mapping-explosion limits. docs
35. Elastic — Strings are dead, long live strings: https://www.elastic.co/blog/strings-are-dead-long-live-strings — 5.0 split string → text + keyword; auto `.keyword` sub-field. blog (vendor)
36. Elastic — fielddata (7.x): https://www.elastic.co/guide/en/elasticsearch/reference/7.4/fielddata.html — fielddata = on-demand in-heap, built by uninverting per segment; disabled by default; exact error text. docs
37. StackOverflow — "Fielddata is disabled…" error: https://stackoverflow.com/questions/60707430/ — fix = `.keyword`. forum
38. deadends.dev — Fix fielddata illegal_argument_exception: https://deadends.dev/elasticsearch/illegal-argument-exception-fielddata/ — enabling fielddata=true is a dead end (OOM); use `.keyword`. blog
39. Elastic — Anatomy of an analyzer (8.19): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/analyzer-anatomy.html — char filters (0+) → tokenizer (1) → token filters (0+). docs
40. Elastic — Standard analyzer: https://www.elastic.co/docs/reference/text-analysis/analysis-standard-analyzer — default; UAX #29; lowercases; stop disabled by default. docs
41. Elastic — Analyzer reference + Index/search analysis: https://www.elastic.co/docs/reference/text-analysis/analyzer-reference — built-in analyzers; same analyzer at index + search normally. docs
42. Elastic — Index and search analysis: https://www.elastic.co/docs/manage-data/data-store/text-analysis/index-search-analysis — query string analyzed at search time; when to diverge. docs
