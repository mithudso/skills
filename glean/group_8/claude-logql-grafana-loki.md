# logql-grafana-loki

**Category:** Databases, Data Engineering & Analytics
**Platform:** Claude
**Original Path:** claude/standalone/logql-grafana-loki

## Description
LogQL — the query language for Grafana Loki, the label-indexed "schema-at-query" log system. PromQL-inspired; two query types: log queries (return lines) and metric queries (logs → time series). Sibling of splunk-platform-spl, kql-kusto-query-language, and elasticsearch-opensearch under "Log & event query languages". TRIGGER: writing/debugging a LogQL query; the required log stream selector {label="value"} and matchers (=, !=, =~, !~); the log pipeline (line filters |= != |~ !~, pattern |> !>, parsers json/logfmt/pattern/regexp/unpack, label filters, line_format, label_format); metric queries (rate, count_over_time, *_over_time, sum/avg/quantile by/without); unwrap for numeric values; the __error__ parse-failure label; left-to-right filter ordering; Loki's low-cardinality-label discipline and why high-cardinality labels destroy performance; LogQL-vs-PromQL / -vs-SPL / -vs-KQL. SKIP: Loki cluster operation, ingesters/queriers, TSDB index, retention/deployment → Grafana Loki operational docs directly; PromQL on Prometheus/Mimir → Prometheus docs; shipping logs into Loki (Promtail/Alloy/Fluent Bit/OTel Collector) → telemetry-pipeline; OTel SDK, Pino, Sentry, eBPF → devops-observability; Splunk SPL → splunk-platform-spl; KQL → kql-kusto-query-language; Elasticsearch ES|QL / OpenSearch PPL / Query DSL → elasticsearch-opensearch.

---

# LogQL (Grafana Loki)

LogQL is the query language for **Grafana Loki**, a horizontally scalable, multi-tenant log aggregation system inspired by Prometheus. Its defining design choice: Loki **indexes only a small set of labels**, not the full log content. A schema is **inferred at query time** ("schema-at-query"), not at ingest — so you store cheaply and parse on read.

> Verified-as-of 2026-06-29 against Grafana Loki docs (`grafana.com/docs/loki/latest/query/`). Loki/LogQL evolve across versions; reconfirm operator and parser syntax against the version the customer runs.

## The single most important concept: the stream

A **stream** is a unique combination of label key/value pairs (e.g. `{app="api", env="prod", level="error"}`). Loki indexes these labels; everything else lives unindexed inside the log line. **Every LogQL query must begin with a log stream selector** — it is not optional, because it tells Loki which streams' chunks to read.

**Cardinality discipline is the whole game.** Labels must be **low-cardinality** (app, env, namespace, level — small finite sets). Putting high-cardinality values (user IDs, request IDs, IPs, timestamps) **in labels** explodes the number of streams and destroys Loki performance and cost. High-cardinality data belongs **in the log line**, extracted at query time with a parser. This is the opposite instinct from Elasticsearch, and the most common LogQL mistake.

## Two query types

1. **Log queries** — return the actual log lines (filtered/parsed/reformatted). Output stays textual.
2. **Metric queries** — wrap a log query in a function to produce a **time series** (a number over time). This is how you build dashboards and alerts from logs.

## Query structure (log query)

```
<stream selector> <log pipeline>
{namespace="prod", app="api"} |= "error" | logfmt | status >= 500 | line_format "{{.status}} {{.msg}}"
```

A query is a **stream selector** followed by an optional **log pipeline** of stage expressions chained left to right. Each stage filters, parses, or mutates lines and their labels.

### 1. Log stream selector (required)
`{label="value", label2=~"regex"}` with matchers:
- `=` equals · `!=` not equals · `=~` regex matches · `!~` regex not-matches.
The more specific the selector, the fewer streams scanned — selector specificity is the **first** performance lever.

### 2. Line filters (a distributed grep over line content)
- `|=` line contains · `!=` line does not contain · `|~` line matches regex · `!~` line does not match regex (all case-sensitive by default).
- Pattern match filters: `|>` (line matches pattern, `<_>` = wildcard) · `!>` (does not). Faster and clearer than regex for fixed-shape lines.
Line filters run on raw line bytes and are **very fast** — apply them **before** parsing.

### 3. Parsers (extract labels from the line — "schema at query")
- `| json` — parse JSON (optionally `| json field="expr"` for selected paths).
- `| logfmt` — parse `key=value` logfmt lines.
- `| pattern "<ip> - <_> [<_>] \"<method> <uri> <_>\""` — extract by a simple positional pattern.
- `| regexp "(?P<name>re)"` — named-capture regex extraction.
- `| unpack` — unpack labels packed by Promtail/Alloy's pack stage.
Extracted fields become **query-time labels** usable by later stages and by metric aggregations.

### 4. Label filter expressions (filter on labels, original or extracted)
`| status >= 500`, `| duration > 1s`, `| level="error"`, with `and`/`or`/`,`. Comparisons are typed: numbers, durations (`1s`, `500ms`), bytes (`1MB`), and strings. **The only way to filter out parse errors** is a label filter on `__error__` (see below).

### 5. Line/label formatting (mutate output)
- `| line_format "{{.label}} ..."` — rewrite the log line using Go template syntax over labels.
- `| label_format new=old` / `| label_format x=`{{.a}}`` — rename/template labels.
- `| decolorize`, `| drop label,...`, `| keep label,...`.

## Metric queries (logs → time series)

Wrap a log query (a "log range") in a range-aggregation function with a duration:

- **Log range aggregations** (count lines/bytes): `rate(<logq>[5m])`, `count_over_time(<logq>[5m])`, `bytes_rate(...)`, `bytes_over_time(...)`, `absent_over_time(...)`.
- **Unwrapped range aggregations** (aggregate a numeric **value** extracted from the line): end the pipeline with `| unwrap <label>` then use `sum_over_time`, `avg_over_time`, `max_over_time`, `min_over_time`, `quantile_over_time(0.99, ...)`, `stddev_over_time`, `rate(... | unwrap bytes [5m])`.
- **Vector aggregations** (group across streams, like PromQL): `sum(...) by (label)`, `avg`, `max`, `min`, `topk`, `bottomk`, `count`, `stddev`, `stdvar`, plus `... without (label)`.

```logql
sum by (status) (rate({app="api"} | logfmt | __error__="" [5m]))
quantile_over_time(0.95, {app="api"} | json | unwrap duration_ms [5m]) by (route)
```

**Metric queries cannot contain errors** — if a parse error occurs during a metric query, Loki errors out. So filter errors (`| __error__=""`) before aggregating.

## The `__error__` label (parse-failure model)

When a parser/label-filter fails on a line, Loki does **not** silently drop it — it passes the line on with a system label `__error__` set (e.g. `JSONParserErr`). It cannot be renamed. To exclude failed lines: `| __error__=""`. To inspect them: `| __error__!=""`. For metric queries, filtering errors is effectively required.

## Performance checklist (order matters — left to right)

1. **Tightest stream selector first** — fewer streams = less data read.
2. **Line filters (`|=`, `|>`) before parsers** — grep on raw bytes is cheap; parse only what survives.
3. **Then label filters**, then formatting.
4. Keep labels low-cardinality; extract high-cardinality fields at query time, don't index them.
5. Prefer `pattern`/`logfmt` over `regexp` when the line shape allows.
6. Bound the time range; large `[range]` windows and wide selectors are the usual cost culprits.
7. In metric queries, `| __error__=""` before aggregating.

## How LogQL differs from its siblings

- **vs PromQL:** LogQL borrows PromQL's selector + vector-aggregation syntax and `by/without`, but adds the **log pipeline** and operates on logs; `rate()` over a log range counts entries, not a counter metric.
- **vs Splunk SPL:** both pipe-based. SPL `search`→ Loki stream selector + line filter; SPL `rex`/`spath`→ `regexp`/`json`; SPL `stats`→ vector aggregation / `*_over_time`; SPL `eval`→ `label_format`. Loki forces you to pick a stream selector up front; SPL does not.
- **vs KQL:** KQL `where`→ line/label filters; KQL `summarize`→ `sum(...) by`; KQL `extend`/`parse`→ `label_format`/parsers. KQL has no mandatory label-selector / cardinality model.

## Common pitfalls

- **High-cardinality labels** (the #1 mistake) → stream explosion, slow queries, high cost. Keep them in the line.
- Omitting or over-broadening the stream selector → scans everything.
- Parsing before line-filtering → wasted work on lines you'll discard.
- Forgetting `| __error__=""` in metric queries → query errors.
- Assuming labels exist before a parser runs — extracted labels only exist *after* the parser stage.
- Expecting full-text indexing like Elasticsearch — Loki greps content; it does not index it.

## Authoritative sources

- LogQL overview / Query Loki — `grafana.com/docs/loki/latest/query/`
- Log queries — `grafana.com/docs/loki/latest/query/log_queries/`
- Metric queries — `grafana.com/docs/loki/latest/query/metric_queries/`
- LogQL Reference (operators/functions) — `grafana.com/docs/loki/latest/query/query_reference/`
- Query examples — `grafana.com/docs/loki/latest/query/query_examples/`
- Labels & cardinality best practices — `grafana.com/docs/loki/latest/get-started/labels/`