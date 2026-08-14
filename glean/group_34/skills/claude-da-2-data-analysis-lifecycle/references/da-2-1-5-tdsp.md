<!-- hub-reference-banner -->
> **Reference file — part of the `da-2-data-analysis-lifecycle` hub.** Formerly the standalone `da-2-1-5-tdsp` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-2-1-5-tdsp
description: |
  Expert knowledge of TDSP (Team Data Science Process) as a process framework within the
  Data Analysis Lifecycle. Covers the five-phase iterative lifecycle, team roles and
  responsibilities, standardized project structure, infrastructure guidance, tooling,
  and how TDSP relates to CRISP-DM and Agile practices.

  TRIGGER: Use this skill when the user asks about TDSP, Team Data Science Process,
  Microsoft's data science process framework, TDSP lifecycle phases (Business Understanding,
  Data Acquisition and Understanding, Modeling, Deployment, Customer Acceptance), TDSP
  team roles (Group Manager, Team Lead, Project Lead, Individual Contributor), TDSP project
  templates, or TDSP adoption in enterprise data science teams.

  SKIP: Do NOT use this skill for CRISP-DM (see da-2-1-1-crisp-dm), KDD (see da-2-1-2-kdd),
  SEMMA (see da-2-1-3-semma), OSEMN (see da-2-1-4-osemn), or general Agile/Scrum practices
  outside the context of data science process frameworks. Also skip for MLOps or DevOps
  questions that are not specifically about TDSP as a lifecycle framework.
---

# TDSP (Team Data Science Process)

## What Is TDSP?

The Team Data Science Process (TDSP) is an agile, iterative data science methodology
designed to deliver predictive analytics solutions and intelligent applications in
production environments. Microsoft released TDSP publicly in 2016 as a distillation of
internal best practices and external industry patterns for enterprise data science teams.

TDSP occupies a specific niche in the process-framework taxonomy: it shares the five-phase
lifecycle shape of CRISP-DM, but adds explicit team structure, standardized tooling,
version control integration, and infrastructure guidance that CRISP-DM deliberately leaves
open. Where CRISP-DM describes *what* phases to run, TDSP also describes *who* does each
task and *how* to organize the work artefacts.

Source: Microsoft Azure Architecture Center — "Team Data Science Process overview"
(https://learn.microsoft.com/en-us/azure/architecture/data-science-process/overview) [1]

---

## Four Core Components

TDSP is not only a lifecycle; it is a four-component system [1][2]:

| Component | What it provides |
|-----------|-----------------|
| **Data science lifecycle** | Five iterative phases from business problem to production model |
| **Standardized project structure** | Shared directory layout and document templates for every project |
| **Recommended infrastructure** | Cloud storage, compute, version control, and collaboration platforms |
| **Tools and utilities** | Scripts that automate common tasks such as data exploration and baseline modeling |

Each component reinforces the others. The directory structure only gains value when teams
adopt version control; the lifecycle only flows predictably when roles are assigned before
sprints begin.

---

## The Five-Phase Lifecycle

The TDSP lifecycle is iterative: teams may cycle back to any earlier phase when findings
warrant it. The phases are not strictly sequential.

### Phase 1 — Business Understanding

**Goal:** Define the business objective as a concrete, measurable prediction or
classification question and identify the data sources needed to answer it.

Key tasks:
- Write a project charter that states the primary objective, success metric, and
  stakeholder sign-off criteria.
- Identify raw data sources and assess access, volume, and freshness.
- Define the analytic approach (regression, classification, clustering, recommendation,
  NLP, etc.).

Outputs: Project charter, preliminary data sources list, analytic approach document.

### Phase 2 — Data Acquisition and Understanding

**Goal:** Ingest all necessary data, build the initial data pipeline, and confirm that
the data can answer the business question from Phase 1.

Key tasks:
- Ingest raw data into the project storage environment.
- Perform exploratory data analysis (EDA): distributions, correlations, missing values,
  anomalies.
- Assess data quality and document known issues.
- Establish automated data pipelines where ongoing ingestion is required.

Outputs: Data pipeline code, data quality report, EDA notebook, updated data dictionary.

This phase corresponds roughly to CRISP-DM's Data Understanding + Data Preparation phases
merged into one, with a stronger emphasis on pipeline automation [2].

### Phase 3 — Modeling

**Goal:** Engineer features and train, evaluate, and select models that meet the
acceptance criteria defined in Phase 1.

Key tasks:
- Feature engineering: construct derived features from raw data; document transformations.
- Model training: train multiple algorithms; track experiments with parameters and metrics.
- Model evaluation: compare candidates on held-out validation data; perform error analysis.
- Select the best model and document its limitations.

Outputs: Feature engineering code, trained model artifacts, evaluation report, model
selection rationale.

### Phase 4 — Deployment

**Goal:** Move the selected model and its data pipeline into a production or
production-like environment where it can serve predictions.

Key tasks:
- Package the model as a callable service (REST API, batch job, embedded library, etc.).
- Deploy the data pipeline alongside the model.
- Set up monitoring for data drift, prediction quality, and system health.
- Write operational runbook.

Outputs: Deployed model endpoint, monitoring dashboard, operational runbook.

### Phase 5 — Customer Acceptance

**Goal:** Validate with the business stakeholder that the deployed system meets the
original success criteria defined in Phase 1 and formally hand off ownership.

Key tasks:
- Run user acceptance testing with the customer or domain expert.
- Document final performance against success metrics.
- Capture lessons learned and archive the project.
- Transfer ownership to the operations or product team.

Outputs: Acceptance sign-off document, lessons-learned report, archived project repository.

This phase is TDSP's most distinctive departure from CRISP-DM: it adds an explicit
external-facing validation gate before the project is considered complete [2][3].

---

## Team Roles and Responsibilities

TDSP defines a hierarchical four-role structure that maps directly onto the lifecycle
phases [1][4]:

### Group Manager
- Oversees the entire data science unit across multiple teams and business verticals.
- Establishes group-wide standards: code hosting accounts, directory templates, utility
  repositories, and security policies.
- Ensures TDSP adoption is consistent across teams.

### Team Lead
- Manages a single data science team, typically 3–8 people.
- Creates team-level repositories from the group template.
- Sets up shared compute, storage, and collaboration access for team members.
- Runs sprint planning and retrospectives.

### Project Lead
- Directs daily work on a single project.
- Creates the project repository from the team template.
- Manages project-specific data assets and access controls.
- Tracks work items (stories, bugs) against the project charter.

### Individual Contributor (Project IC)
- Data scientists, data engineers, analysts, or ML engineers executing the work.
- Clones the project repository; develops and commits code and notebooks.
- Participates in code review through pull requests.
- Produces the artifacts documented in the lifecycle phases above.

A notable TDSP design decision: roles are administrative and collaborative, not
purely technical. The framework explicitly acknowledges that data science and software
engineering are different disciplines and that the Team Lead role bridges them [2].

---

## Standardized Project Structure

Every TDSP project follows the same repository layout, available as a GitHub template
(https://github.com/Azure/Azure-TDSP-ProjectTemplate) [4]:

```
ProjectName/
  Code/               # All executable code: ingestion, modeling, scoring
  Docs/
    Project/          # Project charter, final report
    Data/             # Data dictionary, data quality report
    Model/            # Model performance report, selection rationale
    Deliverable/      # Customer-facing outputs
  Sample_Data/        # Small anonymized data samples for testing
  Notebooks/          # Exploratory and experimental Jupyter notebooks
```

Document templates are provided for:
- Project charter
- Data report
- Model performance report
- Final project report
- Exit report (lessons learned)

The rationale: when every project uses the same layout and every document uses the same
template, new team members can orient themselves within minutes, and organizations build
institutional knowledge that survives personnel turnover [1][4].

---

## Infrastructure and Tooling

TDSP recommends but does not mandate an Azure-centric stack; the principles apply to
any cloud or on-premise environment [1][2]:

**Storage and compute:**
- Cloud file systems (Azure Data Lake, AWS S3, GCS) for raw and processed data.
- Big Data clusters (Hadoop, Spark) for large-scale transformations.
- Machine learning platforms (Azure ML, SageMaker, Vertex AI) for experiment tracking.

**Version control and collaboration:**
- Git (GitHub, Azure DevOps, GitLab) for all code and documents.
- One repository per project, branching model aligned to sprint cycles.
- Pull requests as the integration point for code review.

**Work tracking:**
- Agile boards (Azure Boards, Jira, GitHub Projects) linked to VCS commits for
  traceability between work items and code changes.

**Automation utilities:**
- Microsoft published Python and R utility scripts in the TDSP GitHub organization
  (https://github.com/Azure/Microsoft-TDSP) for data exploration and baseline modeling
  to accelerate onboarding [4].

---

## TDSP and Agile

TDSP deliberately incorporates Agile principles:
- Work proceeds in fixed-length sprints (typically 2–4 weeks).
- Each sprint produces a demonstrable increment (EDA report, baseline model, deployed
  endpoint).
- Requirements and analytic approach are refined iteratively based on findings.
- Retrospectives capture process improvements.

The hybrid posture — Agile sprint cadence wrapped around a CRISP-DM-shaped lifecycle —
means TDSP is more prescriptive than CRISP-DM but more flexible than pure waterfall
project management [3].

---

## Comparison to Peer Frameworks

| Dimension | TDSP | CRISP-DM |
|-----------|------|----------|
| Origin | Microsoft, 2016 | Consortium (SPSS, NCR, etc.), 1999 |
| Lifecycle phases | 5 | 6 |
| Team roles defined | Yes (4 roles) | No |
| Project structure | Standardized template | Not specified |
| Infrastructure guidance | Yes (Azure-oriented) | No |
| Agile integration | Explicit | Not addressed |
| Customer Acceptance phase | Yes (explicit gate) | Embedded in Deployment |
| Tool support | GitHub templates, utility scripts | None official |

TDSP is best thought of as CRISP-DM + team structure + version control discipline +
production deployment focus. Teams already using CRISP-DM can adopt TDSP incrementally by
adding the role definitions and project template without changing their phase sequence [1][2].

---

## When to Use TDSP

TDSP fits well when:
- The team has multiple roles — data scientists, data engineers, and application
  developers working in parallel.
- The project goal is a production-deployed model, not a one-off analysis.
- The organization wants to build institutional knowledge and reuse across projects.
- The team is familiar with or transitioning to Agile/Scrum practices.
- Stakeholder sign-off at project close is required (Customer Acceptance phase).

TDSP is a poor fit when:
- A solo analyst is doing ad-hoc exploratory work without production goals.
- The project is purely statistical research with no deployment component.
- The team is very small (1–2 people) and the overhead of role definitions and structured
  templates exceeds the coordination benefit.
- The organization is locked into a non-Git version control system (TDSP's tooling
  assumes Git).

---

## Common Pitfalls

**Treating phases as strictly sequential.** TDSP is iterative. EDA findings in Phase 2
routinely force a rewrite of the Phase 1 charter. Budget for backtracking.

**Skipping the project structure.** Teams under deadline pressure often skip the
standardized directory and documentation templates, which erodes the institutional
knowledge benefit over time.

**Sprint rigidity vs. exploratory uncertainty.** Fixed-length Agile sprints can create
pressure to declare a model "good enough" before it is. Sprint goals for research phases
should be stated as questions answered, not models delivered.

**Azure lock-in assumption.** TDSP infrastructure guidance is Azure-centric in Microsoft's
documentation. On non-Azure stacks, teams must consciously map the guidance to their
actual tools.

**Conflating Team Lead with technical lead.** The Team Lead role in TDSP is primarily
managerial (sprint planning, repository setup, access control). Technical architecture
decisions typically belong to a senior Individual Contributor or a separately defined
Solution Architect role.

**Ignoring the Customer Acceptance phase.** This phase is TDSP's clearest
differentiator from CRISP-DM. Skipping it means projects often end with a deployed
model but no formal sign-off, leaving stakeholders uncertain whether success criteria
were met.

---

## Worked Example: Churn Prediction Project

**Phase 1 — Business Understanding:**
Project charter states: "Predict monthly subscriber churn 30 days in advance with
precision ≥ 0.70 at 0.25 recall threshold. Use last 24 months of CRM event data."

**Phase 2 — Data Acquisition and Understanding:**
Ingest CRM tables to Azure Data Lake. EDA reveals 18% missing values in the "last
contact date" feature. Data quality report flags this; Phase 1 charter updated to note
the limitation.

**Phase 3 — Modeling:**
Feature engineering creates recency/frequency/monetary (RFM) features. Three candidate
models trained: logistic regression (baseline), gradient boosted trees, neural network.
Gradient boosted trees wins at precision 0.74 / recall 0.28 on validation set. Model
report documents feature importances and calibration curve.

**Phase 4 — Deployment:**
Model packaged as REST API on Azure ML. Batch scoring pipeline runs nightly. Monitoring
dashboard tracks data drift on the top-5 features and prediction volume.

**Phase 5 — Customer Acceptance:**
Marketing team runs A/B test: intervention group (model-targeted) vs. control. After 45
days, churn rate in intervention group 12% lower. Sign-off obtained. Lessons-learned
document archived.

---

## References

1. Microsoft Azure Architecture Center — "Team Data Science Process overview"
   https://learn.microsoft.com/en-us/azure/architecture/data-science-process/overview

2. Data Science PM — "What is TDSP?"
   https://www.datascience-pm.com/tdsp/

3. Towards Data Science — "TDSP: When Agile Meets Data Science"
   https://towardsdatascience.com/tdsp-when-agile-meets-data-science-15ccb5bf8f87/

4. Azure/Microsoft-TDSP GitHub Repository — README and roles-tasks documentation
   https://github.com/Azure/Microsoft-TDSP

5. Azure/Azure-TDSP-ProjectTemplate GitHub Repository
   https://github.com/Azure/Azure-TDSP-ProjectTemplate
