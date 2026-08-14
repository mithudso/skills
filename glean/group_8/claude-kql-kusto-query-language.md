# kql-kusto-query-language

**Category:** Databases, Data Engineering & Analytics
**Platform:** Claude
**Original Path:** claude/standalone/kql-kusto-query-language

## Description
Kusto Query Language (KQL) — the read-only, pipe-based query language across Azure Data Explorer (ADX), Azure Monitor / Log Analytics, Microsoft Sentinel, Defender XDR advanced hunting, Azure Resource Graph, and Microsoft Fabric Real-Time Intelligence (Eventhouse). Sibling of splunk-platform-spl and elasticsearch-opensearch under "Log & event query languages". TRIGGER: writing/debugging a KQL query; the tabular dataflow / pipe (|) model; operators (where, project, extend, summarize, join, union, parse, mv-expand, make-series, render, let, evaluate); time filtering (ago, bin, between, datetime/timespan); aggregations (count, dcount, percentile, arg_max, make_list); join kinds and the innerunique default surprise; KQL-vs-SQL / -vs-SPL translation; performance (time-bound first, has vs contains, column pushdown, materialize, broadcast/shuffle hints); per-product schema/scope differences (Sentinel/Defender/Log Analytics/ADX/Fabric); user-defined & stored functions; query vs management (. control) commands. SKIP: Splunk SPL → splunk-platform-spl; Elasticsearch ES|QL / OpenSearch PPL / Query DSL → elasticsearch-opensearch; telemetry routing/ingest (Cribl/Vector/Fluent Bit/OTel Collector) → telemetry-pipeline; OTel SDK, Pino, Sentry, eBPF → devops-observability; LogQL / Grafana Loki → logql-grafana-loki; ANSI SQL query tuning → deep-query-optimizer; general Azure platform/ARM/networking → Azure docs directly.

---

# Kusto Query Language (KQL)

KQL is a **read-only**, declarative query language built for fast exploration of large volumes of structured, semi-structured, and free-text data — logs, telemetry, and time-series. It was created for **Azure Data Explorer (ADX)** and is now the shared query language across the Microsoft observability/security stack. Reads use KQL; administrative changes use separate **management/control commands** (prefixed with `.`), not KQL.

> Verified-as-of 2026-06-29 against Microsoft Learn (`learn.microsoft.com/kusto/query`, KQL Quick Reference last updated 2025-09-15). KQL evolves; reconfirm operator syntax and per-product schema names before relying on them with a customer.

## Where KQL runs (one language, many products)

The **language** is the same everywhere; the **table schemas, scoping, and available plugins differ** per product. Always confirm which surface the user is on before answering about specific tables.

| Product | What you query | Notable differences |
|---|---|---|
| **Azure Data Explorer (ADX)** | Tables in a Kusto database on a cluster | Full language + management commands; the reference implementation |
| **Azure Monitor / Log Analytics** | Workspace tables (e.g. `AzureDiagnostics`, `Perf`, `Heartbeat`) | Built on Log Analytics store; adds Azure Monitor functions; workspace/`union workspace()` scoping |
| **Microsoft Sentinel** | Same Log Analytics workspace + security tables (`SecurityEvent`, `SigninLogs`) | KQL underlies analytics rules, hunting, workbooks |
| **Defender XDR (advanced hunting)** | Specialized security schema (`DeviceProcessEvents`, `EmailEvents`, `IdentityLogonEvents`) | Subset of operators; its own schema; 30-day default window |
| **Azure Resource Graph** | Azure resource inventory | KQL subset; resource-oriented |
| **Microsoft Fabric Real-Time Intelligence** | Eventhouse / KQL database | Same engine as ADX, surfaced in Fabric |

**Tip for SQL users:** in ADX/Fabric you can preface a SQL statement with `--` and `explain` to get the KQL translation — a fast learning and migration aid.

## The mental model: a left-to-right tabular dataflow

A KQL query is a **source table followed by operators chained with the pipe `|`**. Each operator takes the tabular result of the preceding stage and emits a new table. The result of a query is **always a table**. `T` in the docs denotes "the table from the preceding pipe."

```kql
StormEvents                          // source table
| where StartTime > ago(7d)          // filter (push time-bound FIRST)
| where State == "FLORIDA"
| summarize Events = count() by EventType   // aggregate
| top 10 by Events desc              // order + limit
| render barchart                    // optional visualization
```

**Operator order matters** — `top` before `where` produces different results than after, because each stage transforms the data in sequence.

## Core operator vocabulary (the 80%)

- **Filter / search:** `where` (preferred, predicate-based); `search` (scans all columns — slower); `has` (whole-word/token match, **indexed and fast**) vs `contains` (substring, slower); `take`/`limit` (sample, non-deterministic order).
- **Shape columns:** `project` (select + order), `project-away` / `project-keep`, `project-rename`, `project-reorder`, `extend` (calculated column), `print` (single-row scalar output).
- **Aggregate:** `summarize Agg by Group` is the workhorse. Common aggregation functions: `count()`, `countif()`, `dcount()` (approx distinct), `sum()`, `avg()`, `min()`/`max()`, `percentile()`/`percentiles()`, `arg_max(col, *)` / `arg_min` (row of the extreme), `make_list()` / `make_set()`.
- **Sort / limit:** `sort by col [asc|desc] [nulls first|last]`, `top N by expr`.
- **Combine tables:** `join` (see join kinds below), `union T1, T2`, `lookup` (dimension enrichment, optimized leftouter/inner).
- **Reshape:** `mv-expand` (dynamic array → rows), `mv-apply`, `parse` (structure unstructured strings; `kind=regex|simple|relaxed`), `make-series` (gap-filled time series along an axis), `range` (generate a series), `bin()` (bucket values, esp. time).
- **Reuse / extend:** `let` (bind a name to a scalar, a table, or a lambda — i.e. a query-defined function), `invoke` (run a function on the piped table), `evaluate PluginName(...)` (plugins like `bag_unpack`, `pivot`, `narrow`).
- **Visualize:** `render` (`timechart`, `barchart`, `piechart`, `columnchart`, ...).

## Time is first-class

- `ago(7d)`, `now()`, `datetime(2026-01-01)`, `timespan` literals (`1h`, `30m`, `7d`).
- `bin(Timestamp, 1h)` buckets time for `summarize ... by bin(Timestamp, 1h)`.
- `between (datetime(...) .. datetime(...))` for ranges; `startofday()`, `endofweek()`, `format_datetime()`.
- **Always filter on the time column as early as possible** — it is the single biggest performance lever and lets the engine prune data extents.

## Join kinds — and the `innerunique` gotcha

`join` defaults to **`kind=innerunique`**, not `inner`. `innerunique` deduplicates the **left** side's join keys before joining, which can silently drop expected rows for SQL users who assume a standard inner join. Specify the kind explicitly. Full set: `inner`, `innerunique`, `leftouter`, `rightouter`, `fullouter`, `leftanti` / `leftantisemi`, `rightanti` / `rightantisemi`, `leftsemi`, `rightsemi`. Put the **smaller table on the left** and use `hint.strategy=broadcast` (small) or `hint.strategy=shuffle` (large, high-cardinality key) to tune distribution.

## Performance checklist

1. **Time-bound first** (`where Timestamp > ago(...)`) before any other operator.
2. Prefer **`has`** over `contains`; prefer `where` over `search`.
3. **`project` early** to drop unneeded columns (column pushdown reduces scanned data).
4. Filter **before** `join`/`summarize`, not after.
5. Use `summarize ... by bin()` instead of post-aggregation bucketing.
6. `materialize()` a subquery reused multiple times in a `let`.
7. Use `dcount()` (approximate) unless exact distinct counts are required (`count_distinct()` is exact but heavier).
8. Watch result-set and query limits (e.g. ADX `set` statement limits, Log Analytics row/size caps); use `take` while iterating.

## Query vs management commands

- **Query** (KQL): reads data — everything above.
- **Management / control commands**: start with a dot (`.show tables`, `.create table`, `.alter`, `.ingest`). These are **not KQL queries** and are restricted by role. Don't mix them into a query pipeline.

## How KQL differs from SPL and SQL (for migrators)

- **vs SQL:** KQL is dataflow (pipe) not set-clause (SELECT/FROM/WHERE). `project`≈SELECT, `where`≈WHERE, `summarize`≈GROUP BY, `take`≈TOP/LIMIT. KQL is read-only and column-oriented. Use the `explain` trick for translations.
- **vs Splunk SPL:** both are pipe-based. SPL `search`→KQL `where`; SPL `stats`→`summarize`; SPL `eval`→`extend`; SPL `table`/`fields`→`project`; SPL `dedup`→`distinct`/`summarize arg_max`. Microsoft publishes an SPL→KQL cheat sheet.

## Common pitfalls

- Forgetting time-bounding → full-extent scans and slow/expensive queries.
- Relying on default `join` (it's `innerunique`, not `inner`).
- Using `search` or `contains` where `where`/`has` would be indexed and far faster.
- Assuming table/column names are portable across products — Defender, Sentinel, and ADX schemas differ.
- Treating `take` as deterministic — it is not ordered.
- Confusing management commands (`.`) with queries.

## Authoritative sources

- KQL overview & Sentinel overview — `learn.microsoft.com/kusto/query/`
- KQL Quick Reference (operators/functions) — `learn.microsoft.com/kusto/query/kql-quick-reference`
- Tabular expression statements — `learn.microsoft.com/kusto/query/tabular-expression-statements`
- SQL→KQL cheat sheet — `learn.microsoft.com/kusto/query/sql-cheat-sheet`
- Defender advanced hunting query language — `learn.microsoft.com/defender-xdr/advanced-hunting-query-language`
- Learn common operators tutorial — `learn.microsoft.com/kusto/query/tutorials/learn-common-operators`