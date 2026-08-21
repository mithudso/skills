<!-- hub-reference-banner -->
> **Reference file — part of the `da-2-data-analysis-lifecycle` hub.** Formerly the standalone `da-2-1-1-crisp-dm` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-2-1-1-crisp-dm
description: |
  Expert knowledge of CRISP-DM (Cross-Industry Standard Process for Data Mining) as a process framework within the data analysis lifecycle. Covers the six-phase cycle, tasks and outputs within each phase, the four-level hierarchy of abstraction, iterative vs. sequential use, pitfalls, and positioning relative to other frameworks (KDD, SEMMA, TDSP).

  TRIGGER: Use when the user asks about CRISP-DM specifically — its phases, tasks, outputs, history, how to apply it to a project, or how to compare it to other data mining/science process frameworks. Also trigger when the user is scoping a data mining or analytics project and asks which process model to follow.

  SKIP: Do not use for general data analysis lifecycle questions that are not specifically about CRISP-DM (defer to da-2-data-analysis-lifecycle). Do not use for KDD-specific questions (defer to da-2-1-2-kdd), SEMMA-specific questions (defer to da-2-1-3-semma), OSEMN (defer to da-2-1-4-osemn), or TDSP (defer to da-2-1-5-tdsp). Do not use for questions about problem framing, hypothesis formulation, or other lifecycle sub-topics that have their own skills.
---

# CRISP-DM — Cross-Industry Standard Process for Data Mining

## What It Is

CRISP-DM is an open, vendor-neutral process model that describes the sequence of activities for conducting a data mining or data science project. It organizes work into six phases that loop back on each other until the project goals are met or refined. First published in 1999, it remains the most widely adopted methodology in the field. [1][2]

The name encodes its design intent: *cross-industry* (usable in retail, finance, healthcare, manufacturing, etc.) and *standard process* (a common vocabulary and workflow that a team can follow regardless of the software tools they choose).

## History and Origins

CRISP-DM was conceived in 1996 and became a formal European Union ESPRIT research project in 1997, led by five organizations: Integral Solutions Ltd (ISL), Teradata, Daimler AG, NCR Corporation, and the Dutch insurer OHRA. The first version was presented at the 4th CRISP-DM Special Interest Group Workshop in Brussels in March 1999 and published as a step-by-step guide later that year. [2]

KDnuggets polls conducted in 2002, 2004, 2007, and 2014 consistently showed CRISP-DM as the leading methodology in industry. A 2009 academic review termed it "the de facto standard for developing data mining and knowledge discovery projects." [1][3] Efforts to publish a CRISP-DM 2.0 were started between 2006 and 2008 but stalled; neither the original website nor the 2.0 site remains active. IBM's ASUM-DM (2015) is the most prominent formal extension. [2]

## Four-Level Hierarchy of Abstraction

The methodology is specified in terms of four nested levels [4]:

1. **Phase** — the six top-level stages of the project.
2. **Generic task** — second-level activities within each phase, written to be complete and stable across all possible data mining situations. For example, "build model" is a generic task.
3. **Specialized task** — third-level actions specific to a problem type or tool. For example, "build response model" is a specialized version of "build model."
4. **Process instance** — a record of the actual actions, decisions, and results taken in a specific engagement. This is what a practitioner fills in when running a real project.

The *Reference Model* covers levels 1 and 2. The *User Guide* supplements it with heuristics and examples at levels 3 and 4.

## The Six Phases

The phases are not strictly sequential. The methodology explicitly states: "The sequence of the phases is not rigid. Moving back and forth between different phases is always required." [1] The classic diagram shows arrows cycling through all six, with the outer loop representing that the entire engagement can be repeated as business needs evolve.

### Phase 1 — Business Understanding

The starting phase. Its purpose is to translate a business problem into a data mining goal before any data is touched.

**Key tasks:**
- Determine business objectives: document what the sponsor actually wants to achieve and define business success criteria (measurable where possible, otherwise subjective but agreed in writing).
- Assess situation: inventory available data, personnel, hardware/software, budget; document risks, assumptions, and constraints; conduct a cost-benefit analysis.
- Determine data mining goals: translate the business objective into technical terms. Example: business goal = "increase catalogue sales to existing customers"; data mining goal = "predict how many widgets a customer will buy given their purchase history and demographic attributes." [5]
- Produce a project plan: step-by-step schedule with required inputs, expected outputs, and dependencies for each phase.

**Primary outputs:** business objectives document, data mining goal definition, initial assessment of tools/techniques, project plan.

**Key pitfall:** Skipping or rushing this phase. As Smart Vision Europe notes, "Neglecting this step can mean that a great deal of effort is put into producing the right answers to the wrong questions." [6]

---

### Phase 2 — Data Understanding

Collect data and explore it enough to identify quality issues and initial hypotheses before investing in full preparation.

**Key tasks:**
- Collect initial data and document its source, format, and volume.
- Describe surface properties: schema, record count, attribute types.
- Explore relationships using visualization, queries, and aggregation; document preliminary hypotheses.
- Assess data quality: completeness, accuracy, consistency, timeliness; document any known issues and planned solutions.

**Primary outputs:** initial data collection report, data description report, exploration report with early hypotheses, data quality report.

---

### Phase 3 — Data Preparation

Transform raw data into the form required for modeling. This phase typically consumes the largest share of project effort — commonly cited as roughly 80% of total work. [1][6]

**Key tasks:**
- Select data: choose which tables, rows, and columns to include or exclude, with documented rationale.
- Clean data: handle missing values, outliers, duplicates, and inconsistencies.
- Construct derived attributes: create computed features (ratios, bins, flags, lags) that are not directly in the source.
- Integrate data: merge tables or aggregate records where necessary. Example: "create records for customers who made no purchase in the past year" to represent the negative class. [6]
- Format data: reshape or reorder attributes to match the expectations of the chosen modeling technique.

**Primary outputs:** cleaned dataset, derived attribute documentation, data cleaning report, integrated/aggregated data specifications.

---

### Phase 4 — Modeling

Build, assess, and iterate on analytical models.

**Key tasks:**
- Select modeling technique(s): decision trees, regression, neural networks, clustering, etc.; document assumptions of each.
- Design test procedure: define train/test/validation splits or cross-validation strategy.
- Build models: run the algorithm and record parameter settings.
- Assess models: evaluate against technical success criteria (accuracy, AUC, lift, etc.) and iterate — "iterate model building and assessment until you strongly believe that you have found the best model(s)." [6]

**Primary outputs:** modeling technique documentation, test design, trained models with parameter logs, model assessment report.

**Note:** Iteration commonly cycles back to Data Preparation when modeling reveals missing features or data quality gaps not caught earlier.

---

### Phase 5 — Evaluation

Step back from technical metrics and ask whether the model actually solves the business problem.

**Key tasks:**
- Evaluate results against *business* success criteria, not just technical ones: ask whether "there is some business reason why this model is deficient." [6]
- Review the entire process for overlooked factors, biases, or unconsidered scenarios.
- Determine next steps: deploy, run another iteration, or initiate a new project with a refined scope.

**Primary outputs:** assessment of results against business success criteria, list of approved models, process review notes, ranked list of possible next actions, a go/no-go deployment decision.

**Common failure mode here:** Over-indexing on model accuracy while neglecting deployment readiness — a model that passes evaluation but lacks a production schema match or monitoring plan will fail after launch.

---

### Phase 6 — Deployment

Put the model to use and ensure it keeps working. Deployment ranges from generating a one-time report to embedding a scoring pipeline into a production system. [1]

**Key tasks:**
- Plan deployment: implementation steps, responsible parties, timelines.
- Develop monitoring and maintenance plan: how to detect model drift, who reviews performance, on what schedule.
- Produce final report: document the entire engagement including decisions made, alternatives considered, and the final model.
- Conduct a retrospective: capture lessons learned for future projects.

**Primary outputs:** deployment plan, monitoring plan, final written report, final presentation, experience documentation.

**Critical note:** CRISP-DM treats deployment as a first-class concern, not an afterthought. Predictive value is only realized when the model operates in the environment where decisions are made. [6]

---

## Iterative and Cyclical Use

The outer ring of the CRISP-DM diagram is a circle, indicating that deployment is not the end — business conditions change, new data arrives, and the model's real-world performance feeds back into a new round beginning at Business Understanding. Project plans are explicitly described as "dynamic documents" to be updated at the end of each phase. [6]

In practice, teams execute CRISP-DM in one of two broad styles:

- **Waterfall / horizontal-slice style:** work comprehensively across all deliverables in a phase before moving to the next. Infrequent backtracking; single final delivery. Suits projects with well-defined requirements and stable data.
- **Agile / vertical-slice style:** deliver thin end-to-end slices (a minimal working model, then iteratively improve). Faster stakeholder feedback. Suits exploratory or research-oriented projects where requirements evolve. [1]

Blending CRISP-DM with Kanban, Scrum, or "Data-Driven Scrum" is a recommended pattern to address the framework's lack of built-in team coordination guidance. [1]

## When to Use CRISP-DM

**Good fit:**
- The project has a clear business sponsor and a defined success criterion.
- The team needs a shared vocabulary across business and technical stakeholders.
- The project type is new to the organization and a structured checklist reduces the risk of skipping important steps.
- There is no in-house methodology already in place.

**Less ideal fit:**
- Purely exploratory research without a business objective (CRISP-DM's Business Understanding phase presupposes one).
- Teams already embedded in CI/CD or MLOps workflows may find the CRISP-DM document cadence heavy; lighter overlays (e.g., just using the phase names as sprint labels) often work better in those contexts.

## Common Pitfalls

| Pitfall | What goes wrong | Mitigation |
|---------|----------------|------------|
| Skipping Business Understanding | Building a technically correct model that answers the wrong question | Require a signed objectives document before any data work starts |
| Treating phases as strictly sequential | Missing that modeling reveals data gaps that require re-preparation | Plan for backtracking budget in the project schedule |
| Over-documenting | CRISP-DM has a documentation step for nearly every task; teams get bogged down | Right-size: document decisions and rationale, not every action |
| Ignoring Deployment | Model sits on a laptop and is never used | Plan deployment approach in Phase 1; revisit in Evaluation |
| Data leakage and overfitting | Optimistic evaluation metrics that collapse in production | Enforce strict train/test separation before any feature engineering; use held-out test sets |
| No monitoring plan | Model degrades silently as data distribution shifts | Treat the monitoring plan in Deployment as a required output, not optional |

## Relationship to Other Process Frameworks

CRISP-DM, KDD, and SEMMA cover similar ground but with different emphasis [3][7]:

- **KDD (Knowledge Discovery in Databases):** academically oriented; its nine steps map closely to CRISP-DM's six phases, but it emphasizes the algorithmic discovery aspect over business framing.
- **SEMMA (Sample, Explore, Modify, Model, Assess):** developed by SAS Institute; tool-centric and more narrowly scoped; lacks an explicit business understanding or deployment phase. The five SEMMA steps correspond to CRISP-DM phases 2–5.
- **TDSP (Team Data Science Process):** Microsoft's 2016 extension; adds team roles, DevOps tooling hooks, and a project template structure that CRISP-DM leaves to the practitioner.
- **ASUM-DM:** IBM's 2015 refinement of CRISP-DM adding a "Client Engagement" pre-phase and more explicit agile guidance.

The consensus framing: use CRISP-DM to keep business purpose and end-to-end delivery in view; borrow SEMMA's structure when iterating on models with settled scope; consult TDSP when standing up a team-level data science practice with DevOps tooling.

## Worked Illustration

**Scenario:** A retailer wants to reduce customer churn.

- **Business Understanding:** Objective = reduce churn by 10% in 12 months. Data mining goal = build a binary classifier predicting 30-day churn with ≥ 0.75 AUC on a held-out test set.
- **Data Understanding:** Pull 24 months of transaction logs, returns, and support contacts. Discover 15% missing values in the "last contact date" field; note strong correlation between support ticket volume and churn.
- **Data Preparation:** Impute missing dates using median by customer segment. Construct features: recency of last purchase, frequency of contacts, monetary value in last 90 days (an RFM derivation). Split 70/15/15 for train/validation/test.
- **Modeling:** Compare logistic regression, gradient boosting, and random forest. Gradient boosting reaches 0.81 AUC on validation set.
- **Evaluation:** Model performs well technically, but the marketing team points out that the highest-risk segment has already been targeted by a manual outreach campaign — adding a business constraint that excludes that segment reduces actionable scope. Return to Business Understanding to revise the target population.
- **Deployment:** Score all remaining eligible customers weekly via a batch pipeline. Monitor AUC and churn rate monthly; schedule a model refresh if AUC drops below 0.72.

## Sources

[1] Data Science PM, "What is CRISP-DM?" — https://www.datascience-pm.com/crisp-dm-2/

[2] Wikipedia, "Cross-industry standard process for data mining" — https://en.wikipedia.org/wiki/Cross-industry_standard_process_for_data_mining

[3] Azevedo & Santos (2008), "KDD, SEMMA and CRISP-DM: A Parallel Overview," IADIS European Conference on Data Mining — https://www.researchgate.net/publication/220969845_KDD_semma_and_CRISP-DM_A_parallel_overview

[4] Wirth & Hipp (2000), "CRISP-DM: Towards a Standard Process Model for Data Mining," Proceedings of the 4th International Conference on the Practical Application of Knowledge Discovery and Data Mining — https://cs.unibo.it/~danilo.montesi/CBD/Beatriz/10.1.1.198.5133.pdf

[5] Dummies.com, "Phase 1 of the CRISP-DM Process Model: Business Understanding" — https://www.dummies.com/article/technology/information-technology/data-science/general-data-science/phase-1-of-the-crisp-dm-process-model-business-understanding-148161/

[6] Smart Vision Europe, "CRISP-DM Methodology" — https://www.sv-europe.com/crisp-dm-methodology/

[7] ScienceDirect, "A Systematic Literature Review on Applying CRISP-DM Process Model" — https://www.sciencedirect.com/science/article/pii/S1877050921002416
