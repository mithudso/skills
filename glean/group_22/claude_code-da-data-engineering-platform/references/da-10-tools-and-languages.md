<!-- hub-reference-banner -->
> **Reference file — part of the `da-data-engineering-platform` hub.** Formerly the standalone `da-10-tools-and-languages` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-10-tools-and-languages
title: Data Analysis Tools and Languages
version: "1.0.0"
updated: "2026-05-30"
category: data-analysis
origin: local
description: >
  Languages and tools the data analyst actually reaches for in 2026 — Python
  (pandas, Polars, DuckDB, NumPy, SciPy, scikit-learn), R (tidyverse,
  data.table, arrow), SQL (modern dialects, window functions, CTEs),
  Spark/Databricks, dbt for analytics engineering, notebooks (Jupyter,
  Marimo, Quarto), and the performance comparison that drives the choice.
  TRIGGER: user is choosing a language or library for an analysis,
  benchmarking pandas vs Polars vs DuckDB, setting up a notebook environment,
  dealing with notebook anti-patterns (hidden state, out-of-order execution),
  asks "what should I use for X", or needs a reproducible analysis stack.
  SKIP: language-as-language questions (use python-patterns or javascript-nodejs);
  MongoDB query writing (use mongodb-developer); deep statistical theory
  (da-1 / da-6); machine-learning frameworks (da-7).
triggers:
  - data analysis tools
  - data analysis languages
  - pandas vs Polars
  - DuckDB
  - data.table vs tidyverse
  - which tool for data analysis
  - Jupyter notebook setup
  - reproducible analysis
  - notebook anti-patterns
keywords:
  - pandas
  - Polars
  - DuckDB
  - NumPy
  - SciPy
  - scikit-learn
  - tidyverse
  - data.table
  - arrow
  - SQL
  - window-functions
  - CTE
  - Spark
  - PySpark
  - Databricks
  - dbt
  - Jupyter
  - Marimo
  - Quarto
  - papermill
when_to_use:
  - Choosing a stack for a new analysis
  - Comparing pandas vs Polars vs DuckDB performance / API
  - Setting up a reproducible notebook pipeline
  - Migrating from pandas to Polars or DuckDB
  - Choosing SQL vs DataFrame approach for the same problem
  - Picking between tidyverse, data.table, and Python for R-friendly users
when_not_to_use:
  - General Python language questions — use python-patterns
  - MongoDB-specific query optimization — use mongodb-query-performance
  - Deep ML framework comparison (PyTorch vs JAX) — use da-7-machine-learning
  - Pure visualization library comparison — use da-8-data-visualization
related_skills:
  - python-patterns
  - da-3-data-acquisition-sampling
  - da-4-data-cleaning-preparation
  - da-5-exploratory-data-analysis
  - mongodb-bi-connector
---

# Data Analysis Tools and Languages

The choice of stack constrains everything else. This skill is a 2026-current map of the languages, libraries, and notebook environments data analysts actually use, plus the performance and ergonomic tradeoffs that drive the choice.

## When to use this skill

Activate when the user:
- is choosing a stack for a new analysis
- asks `pandas vs Polars`, `Python vs R`, `SQL vs DataFrame`
- needs to set up a reproducible notebook pipeline
- is hitting performance problems with their current stack
- has notebook anti-patterns (hidden state, out-of-order execution) to fix

## When NOT to use this skill

- General language questions → `python-patterns` / `javascript-nodejs`
- MongoDB query writing → `mongodb-developer`
- Deep statistical theory → `da-1-foundations-theory` and children
- ML framework choice (PyTorch / JAX / TF) → `da-7-machine-learning`
- Visualization library choice → `da-8-data-visualization`

---

## Python — the modern default stack

### pandas

The reference DataFrame library. Eager evaluation, single-threaded by default, in-memory. Pandas 2.x ships with Arrow-backed dtypes which closes part of the performance gap with Polars. Still the most common stack because:

- enormous ecosystem (every library produces or consumes a DataFrame)
- the most StackOverflow / tutorial coverage
- adequate for datasets under ~5 GB on a modern laptop

Pain points: SettingWithCopyWarning, axis ambiguity, NaN vs None vs NaT vs pd.NA, slow groupby on large data, memory overhead 5-10× the raw data size.

### Polars

Rust-cored, lazy by default, columnar (Arrow), multi-threaded. The 2024-2026 rising challenger. When to prefer Polars over pandas:

| Signal | Choose Polars |
|---|---|
| Dataset > 5 GB | Yes — Polars streams |
| Lots of groupby + aggregation | Yes — 5-50× faster |
| Fresh codebase, no pandas dependency | Yes — cleaner API |
| Mixed-team Python+Rust | Yes — interop |
| Heavy library ecosystem dependency | No — many libs assume pandas |
| Older Python (< 3.9) | No — Polars requires 3.9+ |

API style: method chaining, expression-based, lazy `.collect()` at the end. Materially different from pandas; not a drop-in replacement.

### DuckDB

In-process OLAP database with full SQL support. Reads Parquet, CSV, Arrow, pandas, Polars directly without ETL. The "SQL alternative" inside Python. When to reach for DuckDB:

- The analysis is naturally expressed in SQL (joins, window functions, CTEs)
- The data lives in Parquet / Arrow files
- You need to query a dataset that doesn't fit in pandas memory but fits on disk
- You want to share the analysis with SQL-native team members

DuckDB outperforms pandas on most aggregation workloads and often matches Polars. The mental model is "in-process Snowflake."

### NumPy, SciPy, scikit-learn

- **NumPy** — n-dimensional arrays, linear algebra. The foundation.
- **SciPy** — statistical tests, optimization, signal processing, sparse matrices.
- **scikit-learn** — the standard ML interface; estimator → fit → predict → score. Mature, well-tested.

### Performance comparison (2026 representative benchmarks)

| Operation | pandas | Polars | DuckDB |
|---|---|---|---|
| Group-by + aggregate, 100M rows | 60s | 4s | 3s |
| Join, 10M × 10M | 30s | 2s | 1.5s |
| Read 10 GB Parquet | 25s | 8s | 6s |
| API stability | Excellent | Improving | Excellent (SQL) |
| Ecosystem | Largest | Growing | Solid |

Numbers are rough — actual values vary by hardware, columns, and dtypes. The pattern (`pandas` slow, `Polars/DuckDB` fast) holds.

---

## R — the statistician's default

### tidyverse

The `dplyr` / `ggplot2` / `tidyr` / `purrr` cluster. Pipe-based grammar (`%>%` or `|>`), expressive, slower than data.table. Optimized for ergonomics and statistical analysis. Best for:
- moderate datasets (< 1 GB)
- exploratory analysis
- statistical modeling (especially with `broom`, `infer`, `tidymodels`)
- publication-quality static graphics via ggplot2

### data.table

Reference-semantics, in-place, single-line syntax (`DT[i, j, by]`). Fastest R option. Best for:
- large datasets in memory (10 GB+)
- production-grade pipelines
- repeated transforms on the same data

### arrow + duckdb (R bindings)

The same engines as Python. R analysts increasingly use `arrow` for Parquet I/O and `duckdb` for out-of-core SQL.

### Choosing Python vs R

| Need | Python | R |
|---|---|---|
| ML / deep learning | Python (PyTorch, scikit-learn) | R (less depth) |
| Statistical modeling (lme4, brms) | Python (statsmodels, PyMC) | R (gold standard) |
| Web app / dashboard | Python (Streamlit, Shiny-Python, Dash) | R (Shiny) |
| Publication-quality graphics | Python (matplotlib + seaborn, plotnine) | R (ggplot2 leads) |
| Production deployment | Python | R (via Plumber / Posit Connect) |
| Tidyverse-style grammar | Python (Polars approximates) | R (native) |

Either works. The wrong question is "Python or R?"; the right question is "what does the team already know and what does the org already deploy?"

---

## SQL

Still the most-portable, most-deployed data language. Modern dialects (BigQuery, Snowflake, Postgres, DuckDB, Trino) share enough common ground that a competent SQL writer can move between them.

### Modern SQL features the analyst should know

- **Window functions** — `OVER (PARTITION BY ... ORDER BY ...)`, ranking, lead/lag, rolling aggregates. The single biggest leap from "intermediate" to "advanced" SQL.
- **CTEs** — `WITH name AS (...)`. Use them to flatten nested subqueries. Recursive CTEs for hierarchies.
- **Lateral joins** — `CROSS JOIN LATERAL` (Postgres) / `CROSS APPLY` (SQL Server). Row-by-row subquery joins.
- **JSON operations** — `->`, `->>`, `JSON_VALUE`, `JSON_EXTRACT`. Every major dialect has them now.
- **`QUALIFY`** — BigQuery / DuckDB / Snowflake; filter on a window function without wrapping in a subquery.
- **Array / struct types** — Snowflake `ARRAY_AGG`, BigQuery `STRUCT`, DuckDB nested. Move beyond flat tables.

---

## Spark / Databricks

For datasets too large for a single machine. PySpark is the standard Python interface; Databricks runs the managed offering with the Photon engine (a Rust/C++ rewrite of the Spark execution layer that closed the gap with Snowflake on most workloads).

When to reach for Spark:
- Dataset > 100 GB and doesn't fit on one machine even with DuckDB
- Workload is fundamentally distributed (large joins, ML training on huge data)
- Team is already on Databricks for unrelated reasons

When NOT to use Spark:
- Dataset fits in DuckDB (< 1 TB on a beefy node)
- You're optimizing latency, not throughput
- Your team has no JVM operational experience

The classic mistake: reaching for Spark for a 10 GB dataset because "it's big data." Polars or DuckDB on a single node will be 10× faster and 100× cheaper.

---

## dbt — analytics engineering

dbt turned the "SQL transformations" layer from notebooks into engineered software: models, tests, docs, lineage, CI. If you have a warehouse (Snowflake, BigQuery, Redshift, Databricks) and more than one analyst, dbt is the standard.

- **Layers** — staging → intermediate → marts is the canonical layering
- **Tests** — `not_null`, `unique`, `accepted_values`, `relationships`; plus custom tests in SQL
- **Documentation** — auto-generated from model metadata
- **Lineage** — DAG of model dependencies

dbt Core is the open-source CLI. dbt Cloud adds a hosted scheduler + IDE.

---

## Notebooks

### Jupyter (classic)

The default. Massive ecosystem. Weak on reproducibility because:
- out-of-order execution is permitted (hidden state)
- diffs are noisy (cell outputs in JSON)
- no native dependency management
- async issues with long-running cells

### JupyterLab / VS Code notebooks

The modern hosts. Same notebook format, better editor, integrated debugger.

### Marimo

Reactive notebook (2024+). When you change cell A, cells that depend on A re-run automatically. No hidden state. Stored as pure Python files (clean diffs). When the team values reproducibility over interactive scratch-pad culture, switch to Marimo.

### Quarto

Successor to RMarkdown, polyglot (Python, R, Julia). Authoring system for HTML / PDF / docx / slides from a single source. The modern choice when the deliverable is a report, not an interactive notebook.

### papermill + nbconvert

For productionizing Jupyter — parameterize a notebook, run headless, export.

### Notebook anti-patterns to flag in code review

1. **Hidden state** — cell order ≠ execution order. Always restart-and-run-all before sharing.
2. **No version control** — output diffs make notebooks unreviewable. Use `jupytext` or `nbstripout` (or move to Marimo / Quarto).
3. **Monolithic notebook** — one notebook does ingest + clean + EDA + model + report. Split into separate notebooks per stage, glue with papermill or a Makefile.
4. **Hard-coded paths** — `/Users/alice/Desktop/data.csv` in someone else's notebook is useless. Use `pathlib` + project root + a single config file.
5. **Mixing prose for the analyst with prose for the reader** — exploratory scratch comments in the same notebook as the final report. Split.

---

## Reproducible analysis (recipe)

A minimal reproducible analysis stack:

```
project/
  data/
    raw/         # never modified; .gitignore'd if large
    processed/   # derived from raw via scripts
  notebooks/
    01-eda.ipynb
    02-model.ipynb
  src/
    __init__.py
    ingest.py
    clean.py
    model.py
  reports/
    final.qmd    # Quarto report, parameterized
  Makefile       # `make all` reproduces everything
  pyproject.toml # or environment.yml
  uv.lock        # or poetry.lock — pin dependencies
  README.md
```

Reproducibility checks:
- Restart Python, run from top — does it work?
- Clone to a fresh machine, `uv sync && make all` — does it work?
- Two months from now, does it still work? (Pin transitive dependencies.)

---

## Decision matrix — picking a stack

| Constraint | Stack |
|---|---|
| Solo analyst, < 5 GB data, exploratory | Python + pandas + Jupyter |
| Solo analyst, > 10 GB data | Python + Polars or DuckDB |
| SQL-native team, Snowflake / BigQuery warehouse | SQL + dbt + Looker/Tableau |
| Statistical modeling depth (mixed models, survival) | R + tidyverse + brms/lme4 |
| Need ML / deep learning in same pipeline | Python + scikit-learn / PyTorch |
| Need to share interactive analysis with non-coders | Streamlit / Shiny / Observable |
| Production pipeline, scheduled, monitored | Airflow / Dagster + dbt + Python |
| > 100 GB and won't fit on one machine | Spark / Databricks |

---

## References

1. McKinney, W. (2022). *Python for Data Analysis* (3rd ed.). O'Reilly.
2. Wickham, H. & Grolemund, G. (2023). *R for Data Science* (2nd ed.). O'Reilly.
3. Pola.rs documentation — https://pola.rs/
4. DuckDB documentation — https://duckdb.org/docs/
5. dbt docs — https://docs.getdbt.com/
6. Marimo announcement (2024) — https://marimo.io/
7. Quarto — https://quarto.org/
8. Databricks Photon engine paper — Behm et al., SIGMOD 2022
9. pandas-vs-polars benchmarks (TPCH H2O) — h2o.ai
10. Jupyter notebook anti-patterns — Joel Grus, "I Don't Like Notebooks" (JupyterCon 2018)
