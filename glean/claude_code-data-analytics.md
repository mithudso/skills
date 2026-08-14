# data-analytics

**Category:** Science, Biology & Medicine
**Platform:** Claude Code
**Original Path:** claude-code/data-analytics

## Description
Data analytics family ROUTER (da-* hubs). Split into 6 sub-domain hubs: da-1-foundations-theory (probability, measurement, epistemology, inference); da-2-data-analysis-lifecycle (CRISP-DM, problem framing, process workflow); da-3-data-acquisition-sampling (data collection, APIs, web scraping, survey design, ETL); da-analytical-methods (regression, ML, A/B testing, forecasting, causal inference — 16 sub-skills); da-applied-and-communication (visualization, dashboards, ethics, NLP, product analytics); da-data-engineering-platform (dbt, Spark, Airflow, DuckDB, lakehouse, governance). Route to the matching sub-hub.

---

# data-analytics — da-* family ROUTER

Six sub-domain hubs covering the full data analytics stack from theory to platform.

| Sub-hub | Owns | Example topics |
|---|---|---|
| `da-1-foundations-theory` | Probability, measurement, statistical inference, information theory, epistemology | CLT, Bayes, ANOVA theory, levels of measurement |
| `da-2-data-analysis-lifecycle` | End-to-end process: CRISP-DM, KDD, OSEMN, problem framing, iteration | Project scoping, CRISP-DM phases, stakeholder handoff |
| `da-3-data-acquisition-sampling` | Data collection, sampling theory, APIs, scraping, ETL/ELT ingestion | Stratified sampling, Kafka ingest, REST pagination |
| `da-analytical-methods` | Applied techniques — 16 sub-skills on demand | Regression, ML eval, A/B testing, causal inference, CLV, survival |
| `da-applied-and-communication` | Viz, dashboards, data ethics, NLP, product analytics, geo, MMM | Chart selection, D3, cohort retention, differential privacy |
| `da-data-engineering-platform` | Tools, pipelines, warehousing, modeling, governance | dbt, Airflow, DuckDB, ClickHouse, star schema, data vault |

## Routing

Match the primary question:

- **"What does X mean / how does probability/stats work"** → `da-1-foundations-theory`
- **"How do I structure/plan a data project"** → `da-2-data-analysis-lifecycle`
- **"How do I collect/sample/ingest data"** → `da-3-data-acquisition-sampling`
- **"Apply a technique: regression, ML, A/B test, forecast"** → `da-analytical-methods`
- **"Visualize, communicate, ethics, product metrics, NLP"** → `da-applied-and-communication`
- **"Pipelines, tooling, warehousing, dbt, Spark, governance"** → `da-data-engineering-platform`

## Cross-hub map

| Hub | Owns |
|---|---|
| `data-analytics` | This router — all da-* sub-hubs |
| `mongodb-expert` | MongoDB-specific analytics queries and aggregation pipelines |
| `da-analytical-methods` | 16 sub-skills: cleaning, EDA, ML, causal, forecasting, anomaly, NLP, etc. |