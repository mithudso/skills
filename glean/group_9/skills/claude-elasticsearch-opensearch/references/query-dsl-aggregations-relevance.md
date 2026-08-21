# Query DSL, Aggregations, Relevance Scoring, and ES|QL vs OpenSearch PPL/SQL

> Provenance: reference under the `elasticsearch-opensearch` skill. Researched 2026-06-18 via `/dr` deep research (Elastic Query DSL docs, OpenSearch query docs, Elastic ES|QL docs/blog, OpenSearch PPL/SQL docs, Lucene BM25 API, plus practitioner sources). Query DSL, aggregations, and BM25 behave **identically across Elasticsearch and OpenSearch** (shared Lucene); the **piped query languages differ** (ES|QL is Elastic-only). Version-specific claims stamped **verified-as-of 2026-06-18**.

## Contents

- The Query DSL and the `_search` API
- Query context vs filter context (the key distinction)
- Leaf and compound queries (match/term/range/bool) + the term-on-text gotcha
- Aggregations (bucket/metric/pipeline) + the terms-aggregation accuracy caveat
- Relevance scoring (BM25)
- ES|QL (Elastic) vs OpenSearch PPL/SQL

## The Query DSL and the `_search` API

The Query DSL is a **JSON-based DSL** that Elastic describes as an **AST (Abstract Syntax Tree) of queries** with two clause kinds: **leaf** queries (look for a value in a field — `match`, `term`, `range` — usable alone) and **compound** queries (wrap other queries to combine logically — `bool`, `dis_max` — or alter behavior — `constant_score`).[^1][^2] Submit it via the `_search` API's `query` request-body parameter; the body also takes `from`/`size` (pagination), `sort`, and `aggs`. By default hits sort by `_score`.[^3] OpenSearch's Query DSL is largely the same JSON interface with the same query categories.[^5]

## Query context vs filter context (the single most useful distinction)

- **Query context** answers *"How **well** does this document match?"* — it computes a relevance `_score`. Use it for full-text and anything that should affect ranking. Scoring queries' results are **not cacheable**.[^1][^4]
- **Filter context** answers the **binary** *"Does this document match? yes/no"* — **no score is computed.** Benefits: faster execution, lower CPU, and **automatic caching** of frequently-used filters.[^1][^4]
- **Rule of thumb:** put structured conditions (term, range, exists) in **filter context**; reserve query context for relevance-affecting full-text. This is one component set used in two contexts, not two languages.[^1][^4]

## Leaf and compound queries

- **`match`** — the standard full-text query; **analyzes the search text first** (same analyzer used at index time) and builds an internal boolean query from the tokens. `operator` (`or`/`and`, default `or`) and `minimum_should_match` control the clauses.[^6][^7]
- **`match_phrase`** — analyzes text then builds a *phrase* query requiring tokens in order within a configurable `slop` (default 0).[^8]
- **`multi_match`** — runs `match` across multiple fields; types include `best_fields`, `phrase`, `phrase_prefix`.[^9]
- **`term` / `terms`** — **exact, non-analyzed** matching (the term is NOT analyzed).[^10] **`range`** — numeric/date ranges via `gt`/`gte`/`lt`/`lte`. **`exists`** — field has a value (≈ SQL `IS NOT NULL`).[^7]

**The classic "`term` on a `text` field" gotcha (the most common Query DSL mistake):** Elastic's docs say plainly **"Avoid using the `term` query for `text` fields."** Because `text` fields are analyzed at index time (e.g., `"Quick Brown Foxes!"` → `[quick, brown, fox]`), a `term` search for the original string returns **no results** — the exact term no longer exists in the inverted index. The fix: use `match` (which analyzes the query term too), or run `term` against the un-analyzed `.keyword` sub-field, or map the field as `keyword`. Term-level queries are intended for `keyword` fields (filtering, sorting, aggregations).[^10][^12][^13]

**The `bool` query** combines clauses with typed occurrence:[^14]
- **`must`** — must match; **scored** (logical AND); not cached.
- **`should`** — optional; **scored** (boosts relevance / logical OR); not cached.
- **`must_not`** — must not match; **filter context**, scoring ignored, score 0, **cacheable** (logical NOT).
- **`filter`** — must match but **filter context**: score ignored, **cacheable** (logical AND).

**`minimum_should_match` behavior:** when a bool query has **only `should` clauses**, it defaults to **1** (≥1 should-clause must match). When `must` or `filter` is **also present**, `should` defaults to **0** and acts purely as a relevance booster — adding `minimum_should_match` then makes that many `should` clauses required.[^14] OpenSearch's `bool` query is functionally identical.[^5]

## Aggregations

Three categories:[^15]
- **Metric** — compute metrics from field values: `avg`, `sum`, `min`, `max`, `cardinality` (approximate distinct count via HyperLogLog++), `stats` (count/min/max/avg/sum together).
- **Bucket** — group docs into buckets: `terms` (one bucket per unique value, default top 10), `date_histogram` (time-interval buckets), `range`/`histogram`.[^15][^16]
- **Pipeline** — take input from *other aggregations* (not docs/fields) via `buckets_path`; run in the **reduce phase after all other aggregations complete**, so they **cannot be used for ordering**. Examples: `avg_bucket`, `derivative`, `bucket_selector`, `cumulative_sum`.[^17][^18]

Run them via the `aggs` parameter; the `query` parameter limits which docs are aggregated. **Sub-aggregations nest** with no depth limit (e.g., a `terms` agg with an `avg` sub-agg → per-bucket average).[^15] Aggregations read from **doc_values** (columnar), which is why they fail on un-`doc_values` `text` fields (see the fielddata gotcha in `engine-internals-lucene-mappings.md`).

**The terms-aggregation accuracy caveat (important):** terms-aggregation **counts can be inaccurate.** Each shard returns only its **own top-N terms**, then the coordinator merges them — so a term ranking below the cutoff on some shards can be **under-counted or missing** from the global top-N.[^18][^20][^21] Controls:
- **`shard_size`** — how many candidate terms each shard returns (larger = more accurate, more memory/CPU; default is a heuristic of `size`).
- **`show_term_doc_count_error: true`** surfaces **`doc_count_error_upper_bound`** (an upper bound on the count error). A known issue is that this bound can itself be under-estimated in some cases.[^22]
- Even with a large `shard_size`, `doc_count` values **may be approximate, and so may any sub-aggregations** on the terms agg. OpenSearch documents the identical behavior.[^20][^21]

## Relevance scoring (BM25)

- **BM25 (Okapi BM25) is the default similarity since Elasticsearch 5.0**, replacing the legacy `classic` TF/IDF similarity (a consequence of moving to Lucene 6; merged June 2016). Revert a field with `"similarity": "classic"`. BM25 lives in Lucene (`BM25Similarity`).[^23][^24][^25][^26]
- **Components:**[^26][^27][^28]
  - **Term frequency with saturation (`k1`)** — more occurrences raise the score but with **diminishing returns** toward an asymptote. **Default `k1 = 1.2`.**
  - **Document-length normalization (`b`)** — longer docs penalized; `b ∈ [0..1]` controls strength. **Default `b = 0.75`** (`b=0` disables length norm).
  - **IDF (inverse document frequency)** — rarer terms weigh more. Lucene's IDF = `log(1 + (docCount − docFreq + 0.5) / (docFreq + 0.5))`.
- **Tooling:** `_score` is returned per hit; the **explain API** (`_explain` or `"explain": true`) shows whether and *why* a doc matched and how the score was computed.[^30]
- **Boosting:** per-clause `boost` multipliers; the **`function_score` query** modifies scores via functions (`field_value_factor`, decay functions, `script_score`, `weight`), combined by `score_mode` and merged with the query score by `boost_mode` (default `multiply`); `max_boost` caps and `min_score` filters by threshold. Elastic's current guidance favors **multiplicative `function_score` boosting** to shape BM25 with business signals (margin, popularity) while preserving ranking geometry and remaining explainable.[^31][^32]
- **OpenSearch BM25 nuance (verified-as-of 2026-06):** OpenSearch also defaults to Okapi BM25 (`k1=1.2`/`b=0.75`). But **in OpenSearch 3.0 the default changed from `LegacyBM25Similarity` to Lucene's native `BM25Similarity`** — the legacy version included an extra `(k1 + 1)` constant factor in the numerator (a normalization trick that **does not change ranking order**); native `BM25Similarity` drops it, so OpenSearch 3.0+ raw scores are **lower by roughly ~2.2×** than before (ordering unchanged).[^33]

## ES|QL (Elastic) vs OpenSearch PPL/SQL

**Critical distinction:** **ES|QL is Elastic-only.** OpenSearch has **no ES|QL** — its piped language is **PPL** (Splunk-SPL-inspired), plus a separate SQL plugin. They are independent implementations, not the same feature under different names.[^34][^38][^41]

**ES|QL = Elasticsearch Query Language** (piped):
- **Introduced as technical preview in 8.11** (Nov 2023); **reached GA in 8.14** (June 5, 2024). **verified-as-of 2026-06-18.**[^34][^35]
- **Pipe syntax** chains operations with `|`: `FROM index | WHERE condition | STATS agg BY field | SORT field | LIMIT n`. Source command is `FROM`; aggregation is `STATS ... BY ...`; also `EVAL`, `KEEP`, `DROP`, `ENRICH`, `LOOKUP JOIN`.[^34][^36][^37]
- **Dedicated compute engine — no transpilation to Query DSL.** Each query is "parsed, analyzed, validated, and optimized into an execution plan executed in parallel," using a dedicated vectorized, multi-threaded engine (not the DSL's query/aggregation infrastructure), with a coordinating-node global plan plus per-node local re-planning ("hybrid planning"). This is the **opposite** of OpenSearch's PPL/SQL, which compile down to DSL.[^34][^37]
- **Limitations:** returns **up to 1,000 rows by default; `LIMIT` raises it but the hard ceiling is 10,000 output rows** (this caps output, not documents processed). Full-text `MATCH` must appear in `WHERE` directly after `FROM`; `_source` can't be disabled; cross-cluster ES|QL needs an Enterprise license.[^37]

**OpenSearch equivalents (one plugin, two languages, shipped in every distribution):**[^38][^39]
- **PPL (Piped Processing Language)** — pipe-chained, **modeled after Splunk SPL and Unix pipes**, well-suited to observability. A PPL query **starts with `source=<index>`** (or `search source=`) and chains commands with `|`: `source=logs | where status=500 | stats count() by host | sort -count`. Commands include `where`, `fields`, `stats ... by ...` (with `span()` for time/numeric buckets), `head`, `eval`, `rename`, `sort`, `dedup`, `top`, `timechart`. Run via the **`POST _plugins/_ppl`** endpoint or the Dashboards Query Workbench. Gotcha: like ES|QL, PPL **drops fields you didn't reference after `stats`** — project them with `fields` or include in `by`.[^40][^41][^42][^43]
- **SQL** — familiar `SELECT/WHERE/GROUP BY` via the **`_sql`** endpoint, for relational-SQL users.[^38][^39]
- Both can use **`_explain` to translate the query into OpenSearch Query DSL**, confirming PPL/SQL compile down to DSL (newer paths use a Calcite engine) — the design contrast with ES|QL.[^39]

## References

1. Elastic — Query and filter context: https://www.elastic.co/docs/reference/query-languages/query-dsl/query-filter-context — query vs filter, scoring, caching. docs
2. Elastic — Query DSL: https://www.elastic.co/docs/reference/query-languages/querydsl — AST/leaf/compound. docs
3. Elastic — The search API (8.11): https://www.elastic.co/guide/en/elasticsearch/reference/8.11/search-your-data.html — query/size/sort, _score default. docs
4. Elastic Definitive Guide — Queries vs filters: https://github.com/elastic/elasticsearch-definitive-guide/blob/master/054_Query_DSL/65_Queries_vs_filters.asciidoc — one component set, two contexts, cacheability. docs/book
5. logz.io — OpenSearch Queries: Query DSL and Beyond: https://logz.io/blog/opensearch-queries/ — OpenSearch DSL categories, match/term/bool, implicit match_all in aggs. blog
6. Elastic — Match query: https://www.elastic.co/docs/reference/query-languages/query-dsl/query-dsl-match-query — analyzed full-text, operator/min_should_match. docs
7. Elastic Definitive Guide — Important clauses: https://github.com/elastic/elasticsearch-definitive-guide/blob/master/054_Query_DSL/70_Important_clauses.asciidoc — match/multi_match/range/exists. docs/book
8. Elastic — Match phrase query (8.19): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/query-dsl-match-query-phrase.html — phrase query, slop default 0. docs
9. Elastic — Multi-match query: https://www.elastic.co/docs/reference/query-languages/query-dsl/query-dsl-multi-match-query — multi-field, types. docs
10. Elastic — Term query: https://www.elastic.co/docs/reference/query-languages/query-dsl/query-dsl-term-query — "avoid term on text fields," analysis example. docs
12. StackOverflow — find by term fails but match works: https://stackoverflow.com/questions/68829639/ — real-world term-on-text failure, `.keyword` fix. forum
13. StackOverflow — match vs term different results: https://stackoverflow.com/questions/58254201/ — term for keyword fields. forum
14. Elastic — Boolean query: https://www.elastic.co/docs/reference/query-languages/query-dsl/query-dsl-bool-query — must/should/must_not/filter, scoring/caching, minimum_should_match. docs
15. Elastic — Aggregations (8.19): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/search-aggregations.html — three categories, sub-aggregations. docs
16. Elastic — Date histogram aggregation: https://www.elastic.co/docs/reference/aggregations/search-aggregations-bucket-datehistogram-aggregation — date_histogram intervals. docs
17. Elastic — Pipeline aggregations (8.19): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/search-aggregations-pipeline.html — buckets_path, input from other aggs. docs
18. Elastic — Terms aggregation: https://www.elastic.co/docs/reference/aggregations/search-aggregations-bucket-terms-aggregation — top-10 default, inaccuracy, shard_size, doc_count_error, approximate sub-aggs. docs
20. Elastic — Terms aggregation accuracy (same reference): https://www.elastic.co/docs/reference/aggregations/search-aggregations-bucket-terms-aggregation. docs
21. OpenSearch — Terms aggregation: https://docs.opensearch.org/latest/aggregations/bucket/terms/ — shard_size accuracy/memory trade-off (mirrors ES). docs
22. elastic/elasticsearch — issue #35987 (doc_count_error_upper_bound under-estimated): https://github.com/elastic/elasticsearch/issues/35987. forum/issue
23. Elastic — Similarity settings: https://www.elastic.co/docs/reference/elasticsearch/index-settings/similarity — BM25 default, k1=1.2/b=0.75, classic = TF/IDF. docs
24. Elastic — Practical BM25 Part 1: https://www.elastic.co/blog/practical-bm25-part-1-how-shards-affect-relevance-scoring-in-elasticsearch — "In ES 5.0 we switched to Okapi BM25." blog
25. elastic/elasticsearch — PR #18948 (change default similarity to BM25, v5.0, Lucene 6): https://github.com/elastic/elasticsearch/pull/18948. docs/PR
26. Apache Lucene — BM25Similarity (9.12.1): https://lucene.apache.org/core/9_12_1/core/org/apache/lucene/search/similarities/BM25Similarity.html — k1/b definitions, IDF formula. docs
27. Elastic — Practical BM25 Part 2: https://www.elastic.co/blog/practical-bm25-part-2-the-bm25-algorithm-and-its-variables — tf saturation, k1 asymptote, default 1.2. blog
28. Elastic — Similarity in Elasticsearch (found): https://www.elastic.co/blog/found-similarity-in-elasticsearch — full Lucene BM25 weight, IDF, length norm, k1+1 trick. blog
29. Robertson & Zaragoza — The Probabilistic Relevance Framework: BM25 and Beyond: https://www.staff.city.ac.uk/~sbrp622/papers/foundations_bm25_review.pdf — canonical BM25. paper
30. Elastic — Explain API (8.19): https://www.elastic.co/guide/en/elasticsearch/reference/8.19/search-explain.html — score explanation per query+doc. docs
31. Elastic — Function score query: https://www.elastic.co/docs/reference/query-languages/query-dsl/query-dsl-function-score-query — functions, score_mode/boost_mode, max_boost, min_score. docs
32. Elastic Labs — BM25 ranking with multiplicative boosts: https://www.elastic.co/search-labs/blog/bm25-ranking-multiplicative-boosting-elasticsearch — multiplicative function_score preserves BM25 geometry. blog
33. OpenSearch — Keyword search: https://docs.opensearch.org/latest/search-plugins/keyword-search/ — Okapi BM25 default, k1/b; 3.0 switch Legacy→native BM25Similarity, ~2.2× lower scores. docs
34. Elastic Labs — ES|QL now GA: https://www.elastic.co/search-labs/blog/esql-piped-query-language-goes-ga — GA 8.14 (Jun 5 2024), introduced 8.11, no transpilation, vectorized dedicated engine, LIMIT/full-text/CCS notes. blog
35. Elastic — Elastic Stack 8.11 introduces ES|QL: https://www.elastic.co/blog/whats-new-elasticsearch-platform-8-11-0 — 8.11 tech-preview, dedicated compute engine. blog
36. Elastic — ES|QL product page: https://www.elastic.co/elasticsearch/piped-query-language — pipe syntax, concurrent processing. docs/product
37. Elastic — ES|QL limitations: https://www.elastic.co/docs/reference/query-languages/esql/limitations — 1,000 default / 10,000 hard cap, full-text WHERE placement, _source, cross-cluster. docs
38. OpenSearch — SQL and PPL: https://docs.opensearch.org/latest/sql-and-ppl/ — SQL plugin houses both SQL and PPL. docs
39. OpenSearch — SQL and PPL API: https://docs.opensearch.org/latest/sql-and-ppl/sql-and-ppl-api/ — `_sql`/`_ppl` endpoints, `_explain` translates to DSL, Calcite. docs
40. OpenSearch Blog — new PPL capabilities: https://opensearch.org/blog/better-observability-deeper-insights-opensearchs-new-piped-processing-language-capabilities/ — PPL for observability. blog
41. BigDataBoutique — OpenSearch PPL examples: https://bigdataboutique.com/blog/opensearch-ppl-examples-copy-paste-queries — PPL modeled after Splunk SPL/Unix pipes, `source=`, `_plugins/_ppl`, stats drops unreferenced fields. blog
42. OpenSearch — search command (PPL): https://docs.opensearch.org/latest/sql-and-ppl/ppl/commands/search/ — `search`/`source=` first command. docs
43. OpenSearch — stats command (PPL): https://docs.opensearch.org/latest/sql-and-ppl/ppl/commands/stats/ — stats aggregations, span() buckets. docs
