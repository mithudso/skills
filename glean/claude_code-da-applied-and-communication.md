# da-applied-and-communication

**Category:** Databases, Data Engineering & Analytics
**Platform:** Claude Code
**Original Path:** claude-code/da-applied-and-communication

## Description
Applied & domain analytics + communication & data ethics hub (da family). TRIGGER: data visualization (chart selection, dashboards, perception); reporting & communication (analysis-to-narrative, exec summaries, data storytelling, stakeholder handoff); data ethics & privacy (bias/fairness, differential privacy, anonymization, consent, responsible AI); product analytics (funnels, activation, adoption); marketing mix modeling/incrementality (adstock, ROI); geospatial (spatial joins, choropleths); network/graph analytics (centrality, communities); cohort & retention (churn); synthetic data; text analytics & NLP (topics, sentiment, embeddings); augmented/LLM-assisted analytics (NL-BI); pricing & revenue (elasticity). SKIP: theory → da-1-foundations-theory; lifecycle → da-2-data-analysis-lifecycle; acquisition/sampling → da-3-data-acquisition-sampling; methods/ML → da-analytical-methods; pipelines/platform → da-data-engineering-platform.

---

# Data Analysis: Applied Analytics, Visualization, Communication & Ethics

The outward- and domain-facing end of data analysis. After a question has been
framed, data acquired, and methods run (the territory of the sibling hubs), this
hub covers what happens at the edges of an analysis: turning results into
**visualizations** people can read, shaping them into **reports and narratives**
stakeholders can act on, doing it all within **ethical and privacy** constraints,
and applying analytics inside specific **domains** (product, marketing,
geospatial, network, retention, text, pricing), each with its own conventions,
metrics, and pitfalls.

This skill is a hub. Its job is to (1) recognize which applied or communication
sub-topic a task belongs to, and (2) route you to the right reference file before
you give a deep answer. The hub-level body gives orientation; the references hold
the depth.

## How to use this skill

Twelve former standalone skills are absorbed here as reference files under
`references/`. The routing table below maps each sub-topic to its file. **For
anything beyond a one-line orientation, Read the matching `references/` file
before answering.** Several domains (MMM, geospatial, pricing) carry
method-specific assumptions and failure modes that you should ground in the
reference rather than recall from memory. If a task spans two sub-topics
(common: visualization + reporting, or product analytics + cohort/retention),
read both files.

## Sub-skill routing table

When a task matches a row, **Read the listed `references/` file** before
answering.

| Sub-topic | When to load | Reference file |
| --- | --- | --- |
| `da-8-data-visualization` | Data Visualization as a discipline — the principled translation of quantitative and qualitative data into visual encodings; chart selection, perception, dashboards | `references/da-8-data-visualization.md` |
| `da-9-reporting-communication` | Reporting and Communication as the final phase of the data analysis lifecycle — analysis-to-narrative, stakeholder handoff, data storytelling | `references/da-9-reporting-communication.md` |
| `da-11-ethics-and-privacy` | Data ethics and privacy for the working analyst — bias and fairness, differential privacy, anonymization, consent, responsible AI | `references/da-11-ethics-and-privacy.md` |
| `da-21-product-analytics` | Product analytics as an analytical discipline — turning event data into funnels, activation, adoption, and engagement insight | `references/da-21-product-analytics.md` |
| `da-22-marketing-mix-modeling` | Marketing Mix Modeling (MMM) and incrementality measurement — adstock/carryover, saturation, media ROI | `references/da-22-marketing-mix-modeling.md` |
| `da-26-geospatial-analytics` | Geospatial / spatial analytics for general (non-MongoDB) workflows — vector vs raster, spatial joins, choropleths | `references/da-26-geospatial-analytics.md` |
| `da-27-network-graph-analytics` | Network and graph analytics as a data-analysis discipline — representing data as graphs, centrality, community detection | `references/da-27-network-graph-analytics.md` |
| `da-34-cohort-retention-analytics` | Cohort and retention analytics — acquisition vs behavioral cohorts, retention curves, churn | `references/da-34-cohort-retention-analytics.md` |
| `da-35-synthetic-data-generation` | Synthetic data generation as a discipline — generating artificial records that preserve structure while protecting privacy | `references/da-35-synthetic-data-generation.md` |
| `da-36-text-analytics-nlp` | Applied text analytics and NLP for analysts — turning unstructured text into topics, sentiment, and embeddings | `references/da-36-text-analytics-nlp.md` |
| `da-39-augmented-analytics-llm-assisted` | Augmented analytics and LLM-assisted analysis — letting AI plan, query, and narrate analyses; natural-language BI | `references/da-39-augmented-analytics-llm-assisted.md` |
| `da-40-pricing-and-revenue-analytics` | Pricing and revenue analytics — estimating how price drives demand, willingness-to-pay, revenue optimization | `references/da-40-pricing-and-revenue-analytics.md` |

## Cross-hub note

This is one of five sibling hubs in the data-analytics family. Stay in this hub
for applied/domain analytics, visualization, communication, and ethics. Hand off
when a task is really about another layer:

- **Conceptual or mathematical theory** (definitions, measurement, probability,
  inference, epistemology) → `da-1-foundations-theory`.
- **Lifecycle and process** (CRISP-DM and friends, problem framing, success
  metrics, documentation, handoff mechanics) → `da-2-data-analysis-lifecycle`.
- **Data acquisition and sampling** (sources, collection methods, scraping, APIs,
  sampling design) → `da-3-data-acquisition-sampling`.
- **Core statistical and ML methods** (EDA, statistical modeling, machine
  learning, causal inference, forecasting, anomaly detection, feature
  engineering) → `da-analytical-methods`.
- **Pipelines and platform** (warehousing, OLAP, streaming, semantic layer,
  governance, catalogs) → `da-data-engineering-platform`.

Boundary cases worth naming: a *chart* of model output lives here, but the
*model* lives in `da-analytical-methods`; *communicating* a forecast lives here,
but *building* it does not; the *ethics* of a data pipeline lives here, but its
*governance plumbing* lives in `da-data-engineering-platform`. When a task
clearly straddles, read the relevant reference here and cite the sibling hub for
the other half.

**Lifecycle handoff edges** — when the analytical work is done and the output
needs professional prose treatment:

- Analysis results written up for customers, TAM ticket replies, or marketing
  analytics narratives → `content-and-marketing-writing`
- Data-science KB articles, runbooks, or technical specs built from an analysis
  → `technical-writing-craft`
- Executive dashboard briefings or decision memos grounded in analysis
  → `executive-comms`
- Behavior/retention analysis where psychological framing is needed
  (motivation, habit loops, SDT, stages-of-change) → `applied-psychology`

<!-- cross-hub-map -->
## Cross-hub map — where every data-analytics topic lives

This family is split across these hubs. If a task's deep material is **not** in this hub's Sub-skill
routing table, it is a reference file under a sibling hub below — **activate that hub or `Read` its
`references/<name>.md` directly**. Every former standalone skill in this family is now a reference under one
of these hubs (nothing was deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `da-1-foundations-theory` | Data Analysis Foundations & Theory (hub) | `references/da-1-1-definitions-scope.md`, `references/da-1-1-1-data-analysis-vs-analytics-vs-data.md`, `references/da-1-1-2-analysis-vs-synthesis.md`, `references/da-1-1-3-quantitative-vs-qualitative-analysis.md`, … |
| `da-2-data-analysis-lifecycle` | Data Analysis Lifecycle & Process (hub) | `references/da-2-1-process-frameworks.md`, `references/da-2-1-1-crisp-dm.md`, `references/da-2-1-2-kdd.md`, `references/da-2-1-3-semma.md`, … |
| `da-3-data-acquisition-sampling` | Data Acquisition, Collection & Sampling (hub) | `references/da-3-1-data-sources.md`, `references/da-3-1-1-primary-vs-secondary.md`, `references/da-3-1-2-internal-vs-external.md`, `references/da-3-1-3-structured-semi-structured-unstructured.md`, … |
| `da-analytical-methods` | Data Analytical Methods (cleaning, EDA, modeling, ML, causal, time-series) | `references/da-4-data-cleaning-preparation.md`, `references/da-5-exploratory-data-analysis.md`, `references/da-6-statistical-modeling.md`, `references/da-7-machine-learning.md`, … |
| `da-data-engineering-platform` | Data Engineering & Analytics Platform (pipelines, OLAP, modeling, governance) | `references/da-10-tools-and-languages.md`, `references/da-13-data-engineering-and-pipelines.md`, `references/da-14-streaming-analytics.md`, `references/da-18-semantic-layer-headless-bi.md`, … |
| `da-applied-and-communication` | Applied Analytics, Visualization, Communication & Ethics | `references/da-8-data-visualization.md`, `references/da-9-reporting-communication.md`, `references/da-11-ethics-and-privacy.md`, `references/da-21-product-analytics.md`, … |
| `data-analytics` | Family ROUTER — entry point for all da-* sub-hubs | (this file's parent hub) |