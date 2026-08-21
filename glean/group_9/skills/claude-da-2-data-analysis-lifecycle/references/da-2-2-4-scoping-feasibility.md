<!-- hub-reference-banner -->
> **Reference file — part of the `da-2-data-analysis-lifecycle` hub.** Formerly the standalone `da-2-2-4-scoping-feasibility` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-2-2-4-scoping-feasibility
description: |
  Guides scoping and feasibility assessment for a data analysis or analytics project
  during the problem-framing stage of the data analysis lifecycle. Covers how to
  bound the analytical question, assess data availability, evaluate resource and
  technical constraints, and produce a go/no-go judgment before committing to
  full analysis.

  TRIGGER: User is in the early framing/definition stage of an analytics project
  and needs to determine what is in scope, whether the project can succeed given
  available data and resources, or how to structure a feasibility check. Also
  triggers when user asks about scope creep prevention, SMART objectives for
  analytics, data availability assessment, or pre-analysis feasibility criteria.
  Keywords: "scope the project", "is this feasible", "do we have the data",
  "define scope", "project constraints", "scope creep", "can we answer this
  question with our data", "feasibility study", "data requirements", "SMART goals
  for analysis".

  SKIP: Feasibility in the sense of statistical power calculations or sample-size
  planning (that is covered by da-1-4 / da-12). Skip for engineering feasibility
  of a data pipeline (da-13) or ML model feasibility/architecture choices (da-7).
  Skip when the user is already past scoping and is inside the data-collection
  phase (da-3) or EDA phase (da-5). Skip for business-case or ROI analysis that
  does not involve an analytics deliverable. Skip when the user's focus is on
  writing hypotheses or designing experiments — that is da-2-2-2 (hypothesis
  formulation).
---

# da-2-2-4 — Scoping & Feasibility

**Taxonomy location:** Data Analysis > Data Analysis Lifecycle > Problem Framing (2.2) > Scoping & Feasibility (2.2.4)

---

## 1. What it is

Scoping fixes the outer boundary of an analytics project: what questions will be answered, what data will be touched, what outputs will be produced, and what will be deliberately left out. Feasibility asks whether the project can succeed given that boundary — whether data, time, skills, and organizational will are sufficient.

The two activities are tightly coupled. A scope that exceeds available data is not feasible; a feasibility assessment that finds gaps should feed back to narrow or reshape the scope. In practice they run in short iterations at the start of the project, and often again when a new data source arrives or a stakeholder adds a requirement.

This concept sits inside *problem framing* (2.2), which means it happens **before** data collection or EDA. The output of scoping and feasibility is a documented scope statement plus a go/no-go signal — not yet an analytical result.

---

## 2. Why it matters

Without explicit scoping:

- **Scope creep** occurs when stakeholders add new requirements mid-project, overloading the team or forcing re-work ([Render Analytics](https://www.renderanalytics.net/post/why-it-is-important-to-properly-frame-and-scope-a-data-analytics-project)).
- **Misaligned expectations** cause deliverable disputes: one stakeholder expects a predictive model; another expected a dashboard.
- **Wasted effort** flows toward data or questions the organization cannot act on.

Without feasibility assessment:

- Projects invest weeks before discovering the key data source does not exist, is too sparse, or is inaccessible due to legal/privacy constraints ([DSSG Scoping Guide](https://datasciencepublicpolicy.org/our-work/tools-guides/data-science-project-scoping-guide/)).
- Technical and skill gaps surface late, when course-correction is expensive.

The Data Science for Social Good (DSSG) scoping guide frames this starkly: before a project starts, it must meet three threshold criteria — Impact (the problem is real and important), Solvability (data can address it), and Actionability (the organization will act on findings and has implementation capacity). Failing any one is a stop signal.

---

## 3. Dimensions of scope

A complete scope statement answers five questions:

| Dimension | Typical formulation |
|---|---|
| **Analytical question** | What specific question(s) will analysis answer? |
| **Population / unit of analysis** | Which entities (customers, transactions, time windows) are in scope? |
| **Data sources** | Which datasets, tables, or APIs are authorized and relevant? |
| **Analytical method type** | Description, detection, prediction, optimization, or behavior change? |
| **Outputs and consumers** | What deliverable (model, dashboard, report) goes to whom? |

Each dimension also has an explicit **out-of-scope** note. Writing "we will not analyze subsidiary accounts" is as important as writing what is included.

### SMART objectives in analytics scoping

SMART (Specific, Measurable, Achievable, Relevant, Time-bound) is a standard tool for converting vague aspirations into scope-constraining objectives ([LinkedIn Advice: best methods to define analytics scope](https://www.linkedin.com/advice/1/what-best-methods-define-data-analytics-project/)):

- **Specific** — names the metric and the decision it informs, not the technical artefact ("reduce 30-day readmission rate" not "build a model").
- **Measurable** — states the threshold that constitutes success and how it will be measured.
- **Achievable** — acknowledges realistic constraints on data, budget, and timeline.
- **Relevant** — ties the objective to a strategic priority the organization is actually pursuing.
- **Time-bound** — sets a deadline by which the analysis must be complete to be useful.

The DSSG guide adds a warning: "the technical solution (model, dashboard, map) is not itself the goal." Confusing the artefact with the outcome is a common scoping error that produces technically correct deliverables that no one uses.

### Scope tradeoffs to document explicitly

Three pairs of tradeoffs arise repeatedly and should be named in the scope document rather than resolved implicitly:

- **Efficiency vs. equity** — optimizing aggregate outcomes may concentrate harm on subgroups.
- **False positives vs. false negatives** — in detection or prediction tasks, which error type is more costly?
- **Coverage vs. depth** — broad descriptive analysis of all accounts vs. deep investigation of a high-risk segment.

---

## 4. Feasibility assessment

Feasibility has four dimensions that should each be evaluated explicitly:

### 4.1 Data feasibility

The core question: does available data contain the variables needed to answer the scoped question, at sufficient granularity, coverage, and recency?

The SPIFD (Structured Process to Identify Fit-For-Purpose Data) framework from health analytics ([PMC / SPIFD paper](https://pmc.ncbi.nlm.nih.gov/articles/PMC9299818/)) offers a transferable three-step process:

1. **Operationalize and rank minimal criteria.** Define the exact variables needed (exposure, outcome, covariates in research terms; features, target, controls in analytics terms). Rank them as must-have vs. nice-to-have.
2. **Identify and narrow candidate sources.** Screen available datasets against the ranked criteria; reduce to 3–5 candidates for detailed assessment.
3. **Conduct detailed comparison.** For each candidate document: availability of required fields, record count and outcome frequency, completeness and recency, confounders captured, and logistical factors (access speed, contracts, cost, PII constraints).

For general analytics work, the data feasibility checklist reduces to:

| Check | Pass condition |
|---|---|
| Variables present | Every must-have field exists in at least one accessible source |
| Granularity | Records are at the right unit of analysis (e.g., daily transactions, not monthly summaries) |
| Volume | Enough records to support the planned method (especially for ML or significance testing) |
| Recency | Data covers the relevant time window without large gaps |
| Linkability | If multiple sources needed, they share a join key or can be matched |
| Access / legal | Data can be accessed within project timeline; no blocking privacy or contractual barrier |

### 4.2 Technical feasibility

Can the team execute the analysis with available tools and compute? ([LinkedIn: assessing feasibility and viability of a data management project](https://www.linkedin.com/advice/3/how-do-you-assess-feasibility-viability-data-management))

- Skills: does anyone on the team know how to apply the chosen method?
- Infrastructure: is there sufficient compute, storage, and tooling?
- Integration: if the output must connect to a production system, does that path exist?

### 4.3 Organizational feasibility

Will the organization act on the results?

- Is there a named decision-maker who will use the output?
- Does the organization have capacity to implement the recommended action?
- Are there political or cultural factors that would block adoption?

DSSG names this the "actionability" gate: organizations that will not commit resources to act on findings are not good candidates for the project, regardless of data quality.

### 4.4 Time and resource feasibility

Is the timeline realistic given scope and constraints?

- Break the project into phases and estimate each.
- Check whether the deadline precedes a decision the analysis must inform (if the board meeting is in 6 weeks, a 10-week analysis is not feasible regardless of data quality).
- Budget for data access lead time (legal review, vendor contracts). In practice this frequently runs 2–6 weeks — longer than the analysis itself — so initiate access requests before the analytical work begins.

---

## 5. Process: how to run a scoping and feasibility session

A practical sequence draws from both the DSSG guide and the Render Analytics framework:

**Step 1 — Problem understanding (with stakeholders)**
Document the challenge, affected population, current approaches, and what success looks like operationally. Identify all stakeholders who will use or be affected by the analysis.

**Step 2 — Draft scope statement**
Produce an initial scope document covering all five dimensions (analytical question, population, data sources, method type, outputs). Include explicit out-of-scope items. Use SMART criteria to sharpen the analytical objective.

**Step 3 — Data inventory**
List internal data sources with field inventory, granularity, historical depth, update frequency, linkage keys, and access/privacy status. Then identify external sources that fill gaps.
([DSSG Scoping Guide](https://datasciencepublicpolicy.org/our-work/tools-guides/data-science-project-scoping-guide/))

**Step 4 — Run the four-dimension feasibility check**
For each dimension (data, technical, organizational, time/resource) produce a pass/flag/fail rating. Any "fail" is a stop signal; any "flag" needs a mitigation plan before proceeding.

**Step 5 — Document tradeoffs and risks**
Name the efficiency/equity, FP/FN, and coverage/depth tradeoffs. Record assumptions. List open risks with owners and mitigations.

**Step 6 — Go/no-go decision**
Produce an explicit recommendation: proceed, narrow scope, or stop. If scope changes, return to Step 2. When the answer is "narrow scope," prioritize cuts in this order: (1) remove optional outputs or secondary consumers first, (2) reduce the population or time window, (3) downgrade the method type (e.g., prediction → description), (4) defer a data source that is unavailable. Cut in this order to preserve the core analytical question as long as possible.

The whole cycle is typically a 1–2 week sprint for a medium-complexity project. It is meant to be iterative: scope will be refined as data reality becomes clearer.

---

## 6. Common pitfalls

**Skipping written scope** — verbal agreements dissolve. A one-page scope document prevents "that's not what I asked for" disputes.

**Confusing the artefact with the goal** — scoping toward "a predictive model" rather than "reduce 30-day readmissions by 5 points" will produce something technically correct but organizationally unused.

**Optimistic data assumptions** — assuming data exists before inventorying it. The most frequent cause of project failure is discovering late that a key field is missing, sparse, or legally inaccessible.

**Ignoring organizational feasibility** — a technically perfect analysis that no decision-maker is positioned to act on produces no value.

**Scope creep through informal additions** — every new requirement should trigger a re-evaluation of feasibility, not just a "sure, we'll add that."

**Conflating statistical feasibility with data feasibility** — data can be plentiful and still not answer the question (wrong variables, wrong granularity, survivorship bias). Volume alone does not equal feasibility.

---

## 7. Worked example

**Scenario:** A retail operations team wants to "use data to reduce stockouts."

**Initial scope draft:**
- Analytical question: *Which SKUs at which stores are at elevated stockout risk in the next 7 days?*
- Population: All SKUs carried at 120 stores in the western region.
- Data sources: POS transaction system (daily), inventory management system (real-time), supplier lead-time table (quarterly update).
- Method type: Prediction (binary: stockout within 7 days, yes/no).
- Outputs: Weekly risk report for store managers; API feed for automated reorder trigger.

**SMART check on objective:**
- Specific: yes — named metric (stockout), unit (SKU × store), horizon (7 days).
- Measurable: yes — stockout is a binary observable outcome.
- Achievable: *flag* — need to confirm POS data has daily granularity and 18-month history for seasonal patterns.
- Relevant: yes — stockouts directly reduce revenue.
- Time-bound: yes — weekly cadence, model needs to be live before Q4 holiday season (12 weeks away).

**Data feasibility check:**
| Check | Status | Note |
|---|---|---|
| Variables present | Pass | POS, inventory, lead time all exist |
| Granularity | Pass | Daily POS available |
| Volume | Flag | New SKUs (<3 months history) will need cold-start handling |
| Recency | Pass | Systems updated daily |
| Linkability | Pass | SKU + store code is common key |
| Access / legal | Pass | All internal systems, no PII |

**Scope tradeoffs documented:**
- False negatives (missed stockouts) are more costly than false positives (unnecessary reorders); model should optimize recall.
- New-SKU cold-start is out of scope for v1.

**Outcome:** Proceed with narrowed scope (exclude new SKUs in v1). Timeline is tight; data access for POS historical extract needs to be initiated immediately.

---

## 8. Relation to adjacent concepts

| Concept | Relationship |
|---|---|
| 2.2.1 Business/Research Question Definition | Scoping refines and bounds the question produced in 2.2.1; feasibility tests whether it is answerable. |
| 2.2.2 Hypothesis Formulation | Hypotheses are formulated after scope is stable; a scope that cannot support hypothesis testing needs revision. |
| 2.2.3 Success Metrics / KPIs | Success metrics are an output of scoping (they operationalize the objective); their measurability is a data feasibility question. |
| 3 Data Collection | Scoping defines which data to collect; feasibility determines what is available before collection begins. |
| 5 EDA | EDA can surface data quality issues that require scope re-evaluation, feeding back to this step. |

---

## Sources

1. [Data Science Project Scoping Guide — Data Science for Social Good (DSSG/CMU)](https://datasciencepublicpolicy.org/our-work/tools-guides/data-science-project-scoping-guide/) — Comprehensive scoping methodology covering goals, actions, data inventory, and analysis types with ethical framing.

2. [The Structured Process to Identify Fit-For-Purpose Data (SPIFD) — PMC / NCBI](https://pmc.ncbi.nlm.nih.gov/articles/PMC9299818/) — Peer-reviewed three-step framework for data feasibility assessment; reliability and relevancy dimensions.

3. [Why It Is Important to Properly Frame and Scope a Data Analytics Project — Render Analytics](https://www.renderanalytics.net/post/why-it-is-important-to-properly-frame-and-scope-a-data-analytics-project) — Practitioner guide covering 11 scoping dimensions and the case for written scope statements.

4. [What Are the Best Methods to Define Data Analytics Project Scope and Objectives — LinkedIn Advice](https://www.linkedin.com/advice/1/what-best-methods-define-data-analytics-project/) — SMART criteria applied to analytics scoping; stakeholder alignment.

5. [Defining the Scope of Data Analytics Projects — DataCalculus](https://datacalculus.com/en/knowledge-hub/data-analytics/pre-data-collection-phase/defining-the-scope-of-data-analytics-projects/) — Pre-collection phase scope definition including data source and method constraints.

6. [How to Assess Feasibility and Viability of a Data Management Project — LinkedIn Advice](https://www.linkedin.com/advice/3/how-do-you-assess-feasibility-viability-data-management) — Technical, organizational, and resource feasibility dimensions.
