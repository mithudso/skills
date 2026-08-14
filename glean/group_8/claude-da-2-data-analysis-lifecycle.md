# da-2-data-analysis-lifecycle

**Category:** Databases, Data Engineering & Analytics
**Platform:** Claude
**Original Path:** claude/standalone/da-2-data-analysis-lifecycle

## Description
Data-analysis lifecycle & process hub (da family) — the structured workflow from raw data to insight. TRIGGER: overall process/workflow/methodology of an analysis project; structuring/planning start-to-finish; phases/stages; iterate vs advance; process frameworks (CRISP-DM, KDD, SEMMA, OSEMN, TDSP); problem framing (business/research question, hypotheses, success metrics/KPIs, scoping/feasibility); iteration & experimentation loops; documentation & reproducibility (notebooks, provenance, environment capture); stakeholder communication & handoff. SKIP: theory → da-1-foundations-theory; acquisition/sampling → da-3-data-acquisition-sampling; methods/ML → da-analytical-methods; pipelines/platform → da-data-engineering-platform; viz/comms → da-applied-and-communication.

---

# Data Analysis Lifecycle (Process)

**Taxonomy context:** Data Analysis > Data Analysis Lifecycle (Process)

Data Analysis Lifecycle = structured, iterative process. Takes project from business question through data acquisition, preparation, analysis, interpretation, to communicated insight. No single canonical standard — multiple frameworks describe same phases with different emphasis. Lifecycle helps analysts know current phase, what must be true before advancing, when to loop back.

---

## Sub-skill routing table

Hub consolidates 13 lifecycle sub-skills as on-demand reference files. When task matches row, **Read listed `references/<name>.md` before answering** — table alone insufficient for deep answers.

| Sub-topic | When to load | Reference file |
| --- | --- | --- |
| `da-2-1-process-frameworks` | Choosing or comparing structured process frameworks (CRISP-DM, KDD, SEMMA, OSEMN, TDSP) at the family level; deciding which fits a project type. | `references/da-2-1-process-frameworks.md` |
| `da-2-1-1-crisp-dm` | Deep guidance on CRISP-DM phases, bidirectional arrows, phase deliverables, and applying CRISP-DM to a specific project. | `references/da-2-1-1-crisp-dm.md` |
| `da-2-1-2-kdd` | Deep guidance on the KDD process, its 9-step model, and differences from CRISP-DM. | `references/da-2-1-2-kdd.md` |
| `da-2-1-3-semma` | Deep guidance on SEMMA (Sample, Explore, Modify, Model, Assess) and its SAS origins and use cases. | `references/da-2-1-3-semma.md` |
| `da-2-1-4-osemn` | Deep guidance on OSEMN (Obtain, Scrub, Explore, Model, iNterpret) and its academic context. | `references/da-2-1-4-osemn.md` |
| `da-2-1-5-tdsp` | Deep guidance on TDSP (Team Data Science Process): team roles, deliverable templates, sprint cadence, and Azure DevOps integration. | `references/da-2-1-5-tdsp.md` |
| `da-2-2-problem-framing` | Translating a vague business goal into a concrete, measurable analysis objective; the full problem-framing practice. | `references/da-2-2-problem-framing.md` |
| `da-2-2-1-business-research-question-definition` | Writing a well-formed business or research question; distinction between exploratory, descriptive, and causal questions. | `references/da-2-2-1-business-research-question-definition.md` |
| `da-2-2-2-hypothesis-formulation` | Formulating null and alternative hypotheses; operationalizing a business question into testable statistical claims. | `references/da-2-2-2-hypothesis-formulation.md` |
| `da-2-2-3-success-metrics-kpis` | Selecting, specifying, and validating success metrics and KPIs at project inception; metric ownership and threshold-setting. | `references/da-2-2-3-success-metrics-kpis.md` |
| `da-2-2-4-scoping-feasibility` | Scoping an analysis project: data availability assessment, resource and timeline constraints, risk flags, go/no-go criteria. | `references/da-2-2-4-scoping-feasibility.md` |
| `da-2-4-documentation-reproducibility` | Computational notebooks (Jupyter, R Markdown, Quarto), data provenance, environment capture, and reproducibility standards. | `references/da-2-4-documentation-reproducibility.md` |
| `da-2-5-stakeholder-communication-handoff` | Planning and executing stakeholder communication at phase gates; handoff packages for models, findings, and deployed results. | `references/da-2-5-stakeholder-communication-handoff.md` |

---

## 1. Why a lifecycle matters

Raw data not auto-answer questions. Each phase performs distinct transformation:

- **Reduces ambiguity** — vague questions become measurable objectives.
- **Ensures fitness of data** — catches quality problems before corrupting findings.
- **Separates concerns** — keeps exploratory from confirmatory work, analysis from deployment.
- **Creates checkpoints** — natural gates for stakeholder alignment before investing further.

Without explicit lifecycle, projects suffer scope creep, premature modeling on dirty data, findings that can't be reproduced or deployed [Source 1, Source 2].

---

## 2. Major frameworks compared

### 2.1 CRISP-DM (Cross-Industry Standard Process for Data Mining)

Developed late 1990s by Daimler-Chrysler, SPSS, NCR. Most widely cited open standard for data mining and analytics [Source 3].

Six phases in cycle (outer ring can restart after deployment):

| Phase | Core question |
|---|---|
| Business Understanding | What problem are we solving, and how will success be measured? |
| Data Understanding | What data exists, and is it adequate? |
| Data Preparation | How do we transform raw data into a modeling-ready dataset? |
| Modeling | Which technique fits the problem, and how do we build/tune the model? |
| Evaluation | Does the model actually meet the business objective? |
| Deployment | How do stakeholders access and act on the results? |

CRISP-DM arrows bidirectional: unsatisfactory evaluation sends team back to modeling or data preparation; business understanding may revise when data understanding reveals original question is unmeasurable [Source 3].

### 2.2 EMC / Big Data Analytics Lifecycle

Popularized by EMC's *Data Science and Big Data Analytics* book and Wiley companion. Six phases, heavier emphasis on analytic sandboxes and operationalization [Source 4]:

1. Discovery
2. Data Preparation (ELT/ETL into sandbox)
3. Model Planning
4. Model Building
5. Communicate Results
6. Operationalize

Differs from CRISP-DM: explicitly names sandbox as phase 2 prerequisite; distinguishes "Model Planning" (choosing techniques) from "Model Building" (executing).

### 2.3 OSEMN

Minimalist five-step mnemonic popular in academic data science:

- **O**btain
- **S**crub
- **E**xplore
- **M**odel
- i**N**terpret

Strengths: simple, memorable. Weaknesses: omits business framing (starts at "Obtain"), ignores deployment, treats process as linear [Source 5].

### 2.4 TDSP (Team Data Science Process)

Published Microsoft 2017. Closest to CRISP-DM but adds explicit team roles, deliverable templates, agile sprint cadence. Five stages: Business Understanding → Data Acquisition and Understanding → Modeling → Deployment → Customer Acceptance [Source 5].

---

## 3. Canonical phase descriptions

Synthesizes phases common across frameworks into single reference. `references/` files in Sub-skill routing table cover each phase and framework in detail.

### Phase 1 — Problem Definition / Business Understanding

**Input:** stakeholder intent, domain knowledge, prior analyses.
**Output:** written problem statement, success criteria (KPIs or evaluation metrics), initial hypotheses.

Team works with business owners to translate vague goal ("improve customer retention") into concrete measurable objective ("predict 30-day churn with precision ≥ 0.75 at recall ≥ 0.60"). Resources, timeline, risks assessed here.

**Why it matters:** ill-defined question can't be answered with data. Changing question mid-project wastes preparation and modeling effort.

**Pitfall:** treating phase as formality. Teams that rush it often discover midway that available data can't answer question they care about [Source 1].

### Phase 2 — Data Acquisition and Understanding

**Input:** problem statement, knowledge of available data sources.
**Output:** data inventory, quality assessment report, initial summary statistics, decision on whether data sufficient to proceed.

Team collects initial data, examines structure and provenance, documents quality issues (nulls, duplicates, encoding errors, date range gaps), explores distributions and inter-variable relationships.

**Pitfall:** trusting data labeled "clean" is clean. Source systems commonly have undocumented conventions (e.g., sentinel values like -9999 for missing) that only domain knowledge or careful profiling reveals [Source 1, Source 2].

### Phase 3 — Data Preparation

**Input:** raw or semi-structured data, quality assessment.
**Output:** analysis-ready dataset (feature matrix + target variable, or cleaned tabular data for descriptive work).

Typically consumes 60–80% of total project time. Includes:

- **Cleaning:** removing or imputing nulls; correcting format inconsistencies; deduplication.
- **Transformation:** normalization, encoding categorical variables, date parsing, log transforms.
- **Integration:** joining tables across systems, resolving entity mismatches.
- **Feature engineering:** constructing derived columns that encode domain knowledge.

**Analytic sandbox** — compute environment with sufficient CPU, RAM, storage for working data copies — often set up at phase start [Source 4].

**Pitfall — data leakage:** features encoding future information (relative to prediction point) inflate model performance metrics while producing models that fail in production. Any transformation aggregating across full dataset (e.g., z-score mean on both training and test rows) must be fit on training data only and applied to test data [Source 6].

**Pitfall — aggressive outlier removal:** deleting extreme values simplifies modeling but can remove most informative signals, especially in anomaly detection or fraud contexts [Source 6].

### Phase 4 — Analysis / Modeling

**Input:** analysis-ready dataset, modeling plan (technique selection, validation strategy).
**Output:** trained model(s) or analytical findings with performance metrics.

For **descriptive** and **exploratory** analysis: summary statistics, visualizations, identified patterns. For **predictive** analysis: fitted models with cross-validated performance estimates.

Model planning (choosing technique family and validation design) logically distinct from model building (running training and tuning loops). Conflating leads to technique choices driven by familiarity rather than problem fit [Source 4].

**Pitfall — overfitting through hyperparameter tuning:** testing many parameter combinations without held-out test set causes model to fit validation noise, producing strong validation scores that don't transfer to new data [Source 6].

### Phase 5 — Evaluation

**Input:** model or analysis output, success criteria from Phase 1.
**Output:** judgment of whether findings meet original objective; recommendation to proceed or iterate.

Team compares model performance against Phase 1 thresholds, assesses practical and statistical significance, checks model behavior makes sense to domain experts (sanity check catching leakage and labeling errors metrics miss).

Evaluation fails → loop back, usually to Phase 3 (more features, different cleaning) or Phase 2 (additional data sources).

**Pitfall — confusing statistical and practical significance:** result can be statistically significant yet operationally irrelevant. 0.1% improvement in click-through may not justify implementation cost [Source 6].

### Phase 6 — Communication of Results

**Input:** evaluated findings, audience domain knowledge.
**Output:** narrative report, dashboard, or presentation conveying key findings and recommended actions to decision-makers.

Effective communication translates technical outputs into business terms. Team quantifies business value (revenue impact, cost savings, risk reduction), documents key assumptions, acknowledges limitations, prepares supporting materials (code, data dictionaries, reproducibility documentation).

**Pitfall — model explainability missteps:** presenting SHAP plots or feature importances without business context confuses stakeholders. Explanation tools most useful when tied to specific decision audience must make [Source 6].

### Phase 7 — Operationalization / Deployment

**Input:** approved findings or model, deployment environment specifications.
**Output:** running system (scheduled report, API endpoint, embedded model), monitoring plan.

Team deploys model or analysis process for regular stakeholder access. Pilot deployments precede full rollout. Monitoring tracks performance degradation as data distributions shift.

**Pitfall — ignoring concept drift:** model trained on historical data may fail silently as real-world behavior changes. Without monitoring plan and retraining schedule, staleness goes undetected [Source 6].

---

## 4. The iterative nature of the lifecycle

All frameworks represent lifecycle as cyclic or iterative, not linear. Common feedback loops:

- **Evaluation → Data Preparation:** model fails threshold; team engineers features or acquires more data.
- **Modeling → Business Understanding:** most predictive variables are ones business can't act on; problem definition revised.
- **Communication → Problem Definition:** stakeholders raise follow-up question outside original scope; new iteration begins.
- **Operationalization → Data Understanding:** production data differs from training data in distribution; team re-examines source systems.

Treating lifecycle as strictly sequential is recognized anti-pattern. Teams refusing to revisit earlier phases when evidence demands it produce analyses that are technically complete but practically useless [Source 1, Source 5].

---

## 5. Cross-cutting concerns

Apply across all phases, not single ones:

### Documentation and provenance
Every data transformation should be recorded for reproducibility and audit. Missing metadata about data origins is top cause of unreproducible analyses [Source 6].

### Stakeholder alignment
Phase-gate checkpoints — presenting outputs to stakeholders before proceeding — catch misalignment early. Rework cost grows with each phase completed before mismatch surfaces.

### Team roles
In TDSP, roles explicitly assigned: project lead, data scientist, data engineer, solution architect, business analyst each own specific deliverables. Smaller teams cover multiple roles per person, but responsibilities remain distinct [Source 5].

### Governance and ethics
Data collection and use must comply with applicable regulations (GDPR, HIPAA, etc.) and internal data governance policies. Easiest to apply at phase transitions, not after deployment.

---

## 6. Practical worked example

**Scenario:** retail company wants to reduce stockouts.

| Phase | What happens |
|---|---|
| Problem Definition | KPI: reduce stockout events by 20% in 90 days without increasing inventory cost. |
| Data Understanding | Inventory system exports 3 years of daily stock levels; POS system has daily sales by SKU. 8% of SKU-days have NULL stock values. |
| Data Preparation | Impute NULLs using category-level median; create lag features (stock 7 and 14 days prior); join weather data as external feature. |
| Modeling | Train gradient-boosted classifier to predict stockout 7 days in advance; 5-fold time-series cross-validation. |
| Evaluation | Precision 0.78, Recall 0.65 on held-out last-6-months data; meets threshold. Domain review confirms top features make supply-chain sense. |
| Communication | Present to supply chain VP: model will flag ~200 SKUs per week for reorder; estimated 18% reduction in stockout events. |
| Operationalization | Weekly batch job; email alert to purchasing team; monitoring dashboard tracks weekly precision on resolved predictions. |

At evaluation, team discovers weather features add noise rather than signal — loops back to Phase 3, drops those columns, re-evaluates. Normal iteration, not failure.

---

## Sources

1. "Understanding the data analytics lifecycle from end-to-end," Quadratic HQ.
   https://www.quadratichq.com/blog/understanding-the-data-analytics-lifecycle-from-end-to-end

2. "Data Analytics Lifecycle: Phases And Importance," TechCanvass Business Analyst Blog.
   https://businessanalyst.techcanvass.com/data-analytics-lifecycle-phases/

3. "CRISP-DM Methodology: Industry Standard for Data Mining Processes," Medium / Learning Data.
   https://medium.com/learning-data/crisp-dm-methodology-industry-standard-for-data-mining-processes-f896b33dc5ce

4. "6 Phases of Data Analytics Lifecycle Every Data Analyst Should Know," DEV Community / BPB Online.
   https://dev.to/bpb_online/6-phases-of-data-analytics-lifecycle-every-data-analyst-should-know-1k

5. "Data Science Life Cycle: CRISP-DM and OSEMN frameworks," Data Rundown.
   https://datarundown.com/data-science-life-cycle/

6. "Common Pitfalls to Avoid When Analyzing and Modeling Data," freeCodeCamp.
   https://www.freecodecamp.org/news/common-pitfalls-to-avoid-when-analyzing-and-modeling-data/

---

## Outbound routing — when this hub's job ends

Problem framed, analysis ready to begin — hand off to next-phase skill:

| Situation | Route to |
| --- | --- |
| Problem is framed, frameworks chosen — now execute analytical work (EDA, modeling, ML) | `da-analytical-methods` |
| Analysis is complete — now write a KB article, runbook, or technical spec about it | `technical-writing-craft` |
| Analysis is complete — now build an exec presentation, board memo, or decision briefing | `executive-comms` |
| Analysis is complete — now produce a report, visualization, or stakeholder narrative | `da-applied-and-communication` (references/da-9-reporting-communication.md) |

<!-- cross-hub-map -->
## Cross-hub map — where every data-analytics topic lives

Family split across these hubs. If task's deep material **not** in this hub's Sub-skill routing table, it's a reference file under sibling hub below — **activate that hub or `Read` its `references/<name>.md` directly**. Every former standalone skill in this family is now a reference under one of these hubs (nothing deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `da-1-foundations-theory` | Data Analysis Foundations & Theory (hub) | `references/da-1-1-definitions-scope.md`, `references/da-1-1-1-data-analysis-vs-analytics-vs-data.md`, `references/da-1-1-2-analysis-vs-synthesis.md`, `references/da-1-1-3-quantitative-vs-qualitative-analysis.md`, … |
| `da-2-data-analysis-lifecycle` | Data Analysis Lifecycle & Process (hub) | `references/da-2-1-process-frameworks.md`, `references/da-2-1-1-crisp-dm.md`, `references/da-2-1-2-kdd.md`, `references/da-2-1-3-semma.md`, … |
| `da-3-data-acquisition-sampling` | Data Acquisition, Collection & Sampling (hub) | `references/da-3-1-data-sources.md`, `references/da-3-1-1-primary-vs-secondary.md`, `references/da-3-1-2-internal-vs-external.md`, `references/da-3-1-3-structured-semi-structured-unstructured.md`, … |
| `da-analytical-methods` | Data Analytical Methods (cleaning, EDA, modeling, ML, causal, time-series) | `references/da-4-data-cleaning-preparation.md`, `references/da-5-exploratory-data-analysis.md`, `references/da-6-statistical-modeling.md`, `references/da-7-machine-learning.md`, … |
| `da-data-engineering-platform` | Data Engineering & Analytics Platform (pipelines, OLAP, modeling, governance) | `references/da-10-tools-and-languages.md`, `references/da-13-data-engineering-and-pipelines.md`, `references/da-14-streaming-analytics.md`, `references/da-18-semantic-layer-headless-bi.md`, … |
| `da-applied-and-communication` | Applied Analytics, Visualization, Communication & Ethics | `references/da-8-data-visualization.md`, `references/da-9-reporting-communication.md`, `references/da-11-ethics-and-privacy.md`, `references/da-21-product-analytics.md`, … |
| `data-analytics` | Family ROUTER — entry point for all da-* sub-hubs | (this file's parent hub) |