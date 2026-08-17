<!-- hub-reference-banner -->
> **Reference file — part of the `da-2-data-analysis-lifecycle` hub.** Formerly the standalone `da-2-1-process-frameworks` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-2-1-process-frameworks
description: |
  Expert knowledge on structured process frameworks for data analysis — specifically CRISP-DM, KDD,
  SEMMA, TDSP, and the Explore-Refine-Produce (ERP) model — as phases/stages within the data
  analysis lifecycle. Covers when and how to apply each framework, how phases interact iteratively,
  and common failure modes.

  TRIGGER: Use when the user asks about how to structure or plan a data analysis or data mining
  project, which phases or stages a project should follow, how to organize analytical work from
  problem definition through deployment, how CRISP-DM or KDD phases work, what methodology to
  adopt for a data science project, or how analytical workflows are sequenced and iterated.
  Also trigger on questions about the overall shape of an analysis pipeline (business understanding,
  data preparation, modeling, evaluation, deployment) as a process question.

  SKIP: Questions about executing a specific phase in isolation (e.g. "how do I clean data?" or
  "how do I build a regression model?") without asking how it fits into a process framework.
  Skip if the question is about software engineering processes (Agile, Scrum, Kanban) that are
  not specifically combined with a data analysis lifecycle question. Skip if the question is
  about a higher-level concept in this taxonomy — the full lifecycle overview belongs to
  da-2-data-analysis-lifecycle; this skill focuses on the named process framework variants
  (CRISP-DM, KDD, SEMMA, TDSP, ERP) as distinct methodological choices within that lifecycle.
---

# Process Frameworks for Data Analysis

This skill covers named, reusable process frameworks that structure the data analysis lifecycle.
It is scoped to **Data Analysis > Data Analysis Lifecycle (Process) > Process frameworks** —
meaning the question is which methodology to adopt and how to apply it, not what happens during
any single phase.

---

## Why frameworks matter

Ad hoc analysis tends to skip problem definition, treat data preparation as a one-time step,
and omit evaluation against original business goals. Structured frameworks create shared
vocabulary, checkpoints that surface quality issues early, and explicit iteration paths so
analysts know when to go back rather than press forward on flawed assumptions.

A framework does not replace judgment; it provides scaffolding so judgment is applied at
the right moments. [datascience-pm.com; CRISP-DM guide]

---

## The major frameworks

### CRISP-DM (Cross-Industry Standard Process for Data Mining)

**Origin:** Developed in the late 1990s by a consortium including NCR, SPSS, Daimler-Benz,
and OHRA. Published 1999. Still the most widely used framework — KDnuggets polls 2002–2014
and a 2020 site poll (n=109) consistently place it at ~50% adoption. [datascience-pm.com]

**Six phases (iterative, not sequential):**

| Phase | Core question | Typical outputs |
|---|---|---|
| 1. Business Understanding | What problem are we solving, and what counts as success? | Project charter, success criteria, risk assessment |
| 2. Data Understanding | What data do we have, and is it fit for purpose? | Data inventory, quality report, initial findings |
| 3. Data Preparation | What transformations are needed to model-ready form? | Cleaned/integrated dataset (~80% of effort) |
| 4. Modeling | Which algorithms and parameters best fit the problem? | One or more candidate models, tuning records |
| 5. Evaluation | Do the models meet the original business criteria? | Evaluation report, go/no-go decision |
| 6. Deployment | How do results get into production and stay there? | Deployment plan, monitoring protocol, retrospective |

**Critical property — non-linearity:** "The sequence of the phases is not rigid. Moving back
and forth between different phases is always required." [datascience-pm.com] In practice,
Data Understanding frequently reveals that the business problem needs re-framing (loop back
to phase 1). Modeling often exposes data quality issues (loop back to phase 3).

**Strengths:**
- Covers the full project lifecycle, including the often-neglected business framing and
  deployment/maintenance bookends.
- Industry-neutral; applies to finance, healthcare, retail, manufacturing, etc.
- Well-documented: the original 1999 guide is freely available and widely cited.
- Flexible enough to be run in waterfall style (comprehensive horizontal slices) or agile
  style (thin vertical slices delivering incremental value).

**Weaknesses:**
- Documentation requirements are heavy; a strict reading can slow delivery.
- Does not prescribe team roles or project management mechanics, leaving coordination gaps
  on large teams.
- Predates cloud-native, MLOps, and streaming architectures — practitioners must augment it.
- Phases 2 and 3 consume the majority of calendar time; this is often under-scoped in plans.

**Common failure modes:**
- Skipping phase 1 and letting data availability drive problem definition rather than
  business need ("we have this data, so let's see what we can learn from it").
- Treating the framework as strictly sequential and not iterating — teams complete phase 3
  and then refuse to revisit data collection when modeling exposes data gaps.
- Over-documenting in phases 1–2 before having any data intuition, causing analysis paralysis.

---

### KDD (Knowledge Discovery in Databases)

**Origin:** Defined by Fayyad, Piatetsky-Shapiro, and Smyth in their 1996 paper "From Data
Mining to Knowledge Discovery: An Overview" (*AI Magazine* 17(3)). The foundational academic
framework. [U Regina KDD notes; kdnuggets.com]

**Key definition:** "The non-trivial process of identifying valid, novel, potentially useful,
and ultimately understandable patterns in data." KDD is the outer process; data mining is
one step within it.

**Nine steps:**

1. **Domain Understanding** — clarify objectives, prior knowledge, and user goals.
2. **Target Data Selection** — identify which dataset or variables to focus on.
3. **Data Cleaning & Preprocessing** — handle noise, outliers, missing values, and temporal
   order issues.
4. **Data Reduction & Projection** — dimensionality reduction, feature extraction, appropriate
   variable transformations.
5. **Task Selection** — choose the mining task type: classification, regression, clustering,
   anomaly detection, etc.
6. **Algorithm Selection** — pick specific methods and their parameters.
7. **Pattern Mining** — execute the chosen algorithm(s).
8. **Pattern Interpretation** — make sense of discovered patterns; assess validity and novelty.
9. **Knowledge Consolidation** — integrate findings into existing knowledge stores, document,
   and report.

**Pattern quality criteria (Fayyad et al.):** Valid (generalizes beyond the training set),
Novel (not previously known), Potentially useful (actionable), Understandable (interpretable
by domain experts or end users).

**Compared to CRISP-DM:**
- KDD is more granular in the technical discovery steps (steps 4–8 split what CRISP-DM
  calls "Modeling" into finer sub-tasks).
- KDD lacks explicit coverage of business deployment and monitoring (phases 5–6 of CRISP-DM).
- KDD predates CRISP-DM by three years and is more academic in tone.
- Both agree the process is iterative; KDD explicitly notes feedback loops between steps.

---

### SEMMA (Sample, Explore, Modify, Model, Assess)

**Origin:** Developed by the SAS Institute to guide users of SAS Enterprise Miner.
[starburst.io; datascience-pm.com]

**Five phases:**

| Phase | Description |
|---|---|
| Sample | Draw a representative subset from the data source |
| Explore | Visualize and summarize; identify patterns, anomalies, unexpected relationships |
| Modify | Transform, impute, encode, and engineer features |
| Model | Apply statistical or machine learning algorithms |
| Assess | Evaluate model performance with held-out data or business metrics |

**Scope:** SEMMA is deliberately narrower than CRISP-DM — it covers the technical modeling
workflow but omits business problem framing and production deployment. SAS positioned it as
a modeling companion, not a full project lifecycle.

**Strengths:** Provides a clear, memorably acronym-structured approach for technical
practitioners in statistics-heavy environments; well-supported within the SAS tooling ecosystem.

**Weaknesses:**
- Omitting business understanding and deployment means that projects following only SEMMA
  may produce accurate models that answer the wrong question or never reach production.
- Outside the SAS user base, SEMMA is less familiar and less adopted than CRISP-DM.
- "Sample" as a first step can be misleading when working with big data systems where full
  data passes are cheap.

**When to use SEMMA:** As a complement to a broader framework (CRISP-DM or TDSP) when the
analyst's focus is the modeling sub-task and business framing has already been handled upstream.

---

### TDSP (Team Data Science Process)

**Origin:** Launched by Microsoft in 2016. Designed to bring Agile practices to data science
teams working on production ML systems. Published and maintained on GitHub (Azure/Microsoft-TDSP).
[Microsoft Learn; datascience-pm.com]

**Five lifecycle stages:**

1. **Business Understanding** — define objectives, identify data sources, build project plan.
2. **Data Acquisition & Understanding** — ingest data; determine fitness for the stated question.
3. **Modeling** — feature engineering, model selection, training, and tuning.
4. **Deployment** — move model into a production or scoring environment.
5. **Customer Acceptance** — validate with stakeholders that the delivered system meets
   the original goals.

**Key differences from CRISP-DM:**
- TDSP adds explicit *team roles* (solution architect, data scientist, data engineer, project
  lead) and associated responsibilities at each stage.
- Evaluation is embedded within Customer Acceptance rather than a standalone phase.
- Provides a template repository structure (standardized project directory layout, document
  templates) to make adoption practical.
- Designed explicitly for teams shipping *intelligent applications* — not one-off analyses.

**When to prefer TDSP:** When the output is a deployed model or ML service that needs
ongoing maintenance; when team coordination across multiple roles is a problem; when working
in an Azure/Microsoft ecosystem where TDSP tooling integrates natively.

---

### ERP (Explore–Refine–Produce) — Principles for Reproducible Workflows

**Origin:** Lowndes et al. (2017), "Our path to better science in less time using open data
science tools," *Nature Ecology & Evolution*; further formalized in Lowndes & Horst (2021),
"Principles for data analysis workflows," *PLOS Computational Biology*. [PMC7971542]

**Three phases:**

| Phase | Primary audience | Key activities |
|---|---|---|
| Explore | Yourself | Data cleaning, gut-check tests, hypothesis generation; prioritize documentation over optimization |
| Refine | Collaborators | Narrow to promising approaches; increase code quality, testing standards, and modularity |
| Produce | Public/stakeholders | Polish for external critique; generate multiple research products (papers, datasets, tools) |

**Core contribution:** ERP explicitly ties each phase to its *communication audience*, so
analysts calibrate documentation, testing, and code quality to actual needs rather than
over-engineering exploration code or under-documenting production outputs.

**Reproducibility as a spectrum:** Rather than perfect reproducibility as a binary goal,
ERP argues for "good enough" practices — the rigor appropriate to the team's domain, tooling,
and timeline.

**Practical principles from the paper:**
- Separate mechanistic code (plotting, processing functions) from narrative analysis code
  to keep computational notebooks readable.
- Use "gut checks" (dimension checks, missing-value summaries, sanity assertions) during
  exploration instead of formal unit tests, which come in the Refine phase.
- Treat iteration between Explore and Refine as a balance: scope creep from cycling back
  too liberally can prevent finishing.

---

## Framework comparison summary

| Dimension | CRISP-DM | KDD | SEMMA | TDSP | ERP |
|---|---|---|---|---|---|
| Phases / steps | 6 | 9 | 5 | 5 | 3 |
| Covers business framing | Yes | Partially | No | Yes | No |
| Covers deployment | Yes | No | No | Yes | No |
| Audience | Cross-industry | Academic/research | SAS/modeling | Enterprise ML teams | Research/science |
| Agile-compatible | Yes (explicit) | Not specified | Not specified | Yes (explicit) | Yes (iterative) |
| Team role guidance | No | No | No | Yes | No |
| Origin year | 1999 | 1996 | ~1990s | 2016 | 2017–2021 |

---

## Choosing a framework

**CRISP-DM** is the default choice for most applied data analysis projects because it covers
the full lifecycle, is industry-neutral, and has the deepest community knowledge base.

**KDD** is appropriate when publishing academic findings or when the goal is knowledge
extraction and pattern validity matters more than deployment.

**SEMMA** applies when you are working within SAS tooling and the business framing is already
handled; do not use it alone for full project governance.

**TDSP** fits when the team is building a production ML service, especially in a Microsoft/Azure
environment, and when multi-role coordination is a recurring coordination problem.

**ERP** applies to research and scientific data analysis where reproducibility and
communication to different audiences (collaborators, public) are primary concerns.

Teams commonly *combine* frameworks — for example, CRISP-DM for project lifecycle governance
plus Scrum sprints for iteration cadence and team coordination.

---

## Common pitfalls across all frameworks

1. **No clear problem definition.** Starting with available data rather than a stated
   business or research question causes work to drift toward interesting-but-irrelevant findings.
   Every framework puts problem understanding first for this reason.

2. **Linear execution.** Treating phases as a strict sequence and refusing to loop back
   when later phases reveal earlier assumptions were wrong. All five frameworks explicitly
   allow or require iteration.

3. **Data preparation under-scoped.** Plans consistently underestimate preparation time.
   CRISP-DM practitioners report data preparation consumes ~80% of project effort [datascience-pm.com].
   Budgeting 20% for it causes chronic schedule slippage.

4. **Evaluation against the wrong criterion.** Optimizing a model's statistical metric
   (accuracy, AUC) without checking whether it meets the original business criterion leads
   to technically correct but operationally useless outputs. CRISP-DM's explicit Evaluation
   phase exists precisely to catch this.

5. **Premature optimization.** Over-engineering data pipelines or model code during
   exploration before knowing which approaches are worth keeping; recognized by ERP as a
   characteristic Explore-phase anti-pattern. [PMC7971542]

6. **Scope creep in iteration.** The freedom to loop back between phases becomes an excuse
   to keep exploring indefinitely. Setting phase-exit criteria upfront (time-boxes,
   evidence thresholds) prevents this.

---

## Sources

1. "What is CRISP DM?" — Data Science PM. https://www.datascience-pm.com/crisp-dm-2/
2. Fayyad, U., Piatetsky-Shapiro, G., & Smyth, P. (1996). "From Data Mining to Knowledge
   Discovery in Databases." *AI Magazine* 17(3): 37–54. https://www.kdnuggets.com/gpspubs/aimag-kdd-overview-1996-Fayyad.pdf
3. Overview of the KDD Process — University of Regina, CS 831 notes.
   https://www2.cs.uregina.ca/~dbd/cs831/notes/kdd/1_kdd.html
4. "SEMMA vs CRISP-DM: The differences" — Starburst.
   https://www.starburst.io/blog/semma-vs-crisp-dm/
5. "CRISP-DM is Still the Most Popular Framework" — Data Science PM.
   https://www.datascience-pm.com/crisp-dm-still-most-popular/
6. "What is TDSP?" — Data Science PM. https://www.datascience-pm.com/tdsp/
7. Lowndes, J.S.S. et al. (2021). "Principles for data analysis workflows."
   *PLOS Computational Biology*. https://pmc.ncbi.nlm.nih.gov/articles/PMC7971542/
