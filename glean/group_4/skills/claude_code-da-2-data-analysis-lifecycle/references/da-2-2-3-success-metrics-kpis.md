<!-- hub-reference-banner -->
> **Reference file — part of the `da-2-data-analysis-lifecycle` hub.** Formerly the standalone `da-2-2-3-success-metrics-kpis` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-2-2-3-success-metrics-kpis
version: "1.1.0"
updated: "2026-05-30"
description: |
  Defines how to select, specify, and validate success metrics and KPIs during the
  problem framing stage of a data analysis project — before data collection begins.

  TRIGGER:
  - User is framing or scoping a data analysis project and needs to define what "success" looks like
  - User asks how to choose KPIs, metrics, or success criteria for an analysis
  - User asks about the difference between metrics, KPIs, OKRs, or goals in an analytics context
  - User asks how to avoid vanity metrics or Goodhart's Law traps in their measurement design
  - User asks about leading vs lagging indicators, proxy metrics, or North Star metrics in problem framing
  - User asks how to decompose a high-level business goal into measurable sub-metrics
  - User asks what makes a metric "good" for evaluating a data analysis project
  - User asks how many KPIs to define for a project

  SKIP:
  - Metrics that arise after analysis is complete and the question is about communicating findings
    (that belongs to da-9-reporting-communication)
  - Statistical inference or hypothesis testing about whether an observed metric value is significant,
    including significance testing of a KPI change
    (that belongs to da-2-2-2-hypothesis-formulation or da-1-4-statistical-inference-foundations)
  - Evaluating whether a metric changed significantly due to an intervention or A/B test
    (that belongs to da-12-ab-testing-causal-inference)
  - Model evaluation metrics (accuracy, F1, AUC) chosen to tune a machine learning model
    (that belongs to da-7-machine-learning)
  - KPI dashboarding, visualization design, or BI tooling
    (that belongs to da-8-data-visualization)
  - Ongoing performance monitoring once a solution is live and in production
    (that belongs to da-9-reporting-communication)
  - OKR writing as an organizational goal-setting practice
    (that belongs to the okr-writing skill)

related_skills:
  - da-2-2-problem-framing
  - da-2-2-1-business-research-question-definition
  - da-2-2-2-hypothesis-formulation
  - da-12-ab-testing-causal-inference
  - okr-writing
tags:
  - data-analysis
  - problem-framing
  - kpis
  - metrics
  - success-criteria
  - measurement
---

# Success Metrics & KPIs — Problem Framing

## 0. Position in the taxonomy

This skill covers concept **2.2.3** in the data analysis lifecycle:

```
Data Analysis
  └── Data Analysis Lifecycle (Process)      [2]
        └── Problem framing                  [2.2]
              ├── Business research question definition  [2.2.1]
              ├── Hypothesis formulation                 [2.2.2]
              └── Success metrics & KPIs                 [2.2.3]  ← here
```

Defining success metrics happens *before* data collection and *before* modeling. The goal is to agree, while the question is still open, on exactly what evidence would constitute a good answer.

---

## 1. Why success metrics must be set during problem framing

The most common analytical failure mode is not poor modeling — it is answering the wrong question. Success metrics serve as a contract: they prevent scope creep, give stakeholders a shared language, and make it possible to say "done."

Setting metrics after the analysis is complete is a form of HARKing (Hypothesizing After Results are Known). Once you have seen the data you will unconsciously favor metrics on which your results look strong. Agreement at the start guards against this. ([productleadership.com](https://www.productleadership.com/blog/business-problem-framing-in-the-age-of-data-analytics/))

---

## 2. Definitions: goals, metrics, KPIs, OKRs

These terms are often used interchangeably but carry distinct meanings.

| Term | What it is | Timeframe |
|------|-----------|-----------|
| **Goal / objective** | A qualitative direction ("improve retention") | Open-ended |
| **Metric** | Any quantifiable measurement ("monthly active users") | Continuous |
| **KPI** (Key Performance Indicator) | A metric chosen because it indicates progress toward a strategic objective | Periodic review |
| **OKR** (Objective + Key Results) | A goal paired with 2–5 measurable outcomes that confirm it was achieved | Quarterly |

KPIs are a **subset** of metrics — not every metric is a KPI. Choosing which metrics to elevate to KPI status is the central judgment in this stage. ([atlan.com](https://atlan.com/kpis-for-data-team/), [atlassian.com](https://www.atlassian.com/work-management/project-management/project-planning/kpi))

---

## 3. The output / outcome / impact hierarchy

Before selecting specific metrics, place each candidate in the right layer of this hierarchy:

```
Activities  →  Outputs  →  Outcomes  →  Impact
(what we do)   (what we    (what changes  (long-term
                produce)    because of it)  consequence)
```

- **Activities** are the work done: running queries, building models, conducting interviews.
- **Outputs** are direct, tangible products of the work: a report published, a model deployed, a dataset cleaned.
- **Outcomes** are the changes that outputs cause: churn rate dropped, decision latency shortened.
- **Impact** is the long-term, broader consequence: competitive position improved, revenue trajectory shifted.

([spiderstrategies.com — outcome vs output](https://www.spiderstrategies.com/blog/outcome-vs-output-metrics/))

**Common mistake:** framing success purely in output terms ("we will deliver a churn model") rather than outcome terms ("churn will fall 8 percentage points within one quarter"). Output metrics show that work happened; outcome metrics show that the work mattered.

---

## 4. The SMART test for each metric candidate

Every proposed success metric should pass all five tests:

| Letter | Test | Example failure | Fixed version |
|--------|------|-----------------|---------------|
| **S**pecific | Names exactly what is measured | "improve sales" | "increase weekly order volume in the North region" |
| **M**easurable | Names a unit and a data source | "customer happiness" | "NPS score from post-purchase survey (Qualtrics)" |
| **A**chievable | Realistic given available data, tools, and time | "100% accuracy" | "≥ 90% precision at a threshold of 0.6" |
| **R**elevant | Connects to the business question being asked | page views on an operations dashboard | cost per unit processed |
| **T**ime-bound | States when the target should be met | "eventually reduce churn" | "reduce churn 5 pp within Q3" |

([datacalculus.com](https://datacalculus.com/en/knowledge-hub/data-analytics/pre-data-collection-phase/setting-smart-objectives-in-data-analytics/), [simplekpi.com](https://www.simplekpi.com/Blog/smart-and-smarter-kpis-explained))

---

## 5. Leading vs lagging indicators

| Type | Definition | Example | Advantage | Limitation |
|------|-----------|---------|-----------|------------|
| **Lagging** | Measures what already happened; confirms whether the goal was achieved | Revenue, churn, NPS | Definitive; hard to dispute | Available only after the fact; can't course-correct mid-flight |
| **Leading** | Predicts whether you are on track to hit the lagging indicator | Sales calls made, trial-to-paid conversion | Actionable early; enables pivots | Can be gamed; correlation with lagging metric may be weaker than assumed |

A well-designed success framework contains both. The lagging indicator is the *objective* (the thing that actually matters); leading indicators are the *key results* used to monitor progress in real time. ([humanperf.com](https://www.humanperf.com/en/blog/project-management/articles/performance-measurement-indicators), [xodiac.ca](https://xodiac.ca/blog/articles/leading-vs-lagging-indicators-to-measure-okrs))

---

## 6. Proxy metrics

Sometimes the true success metric is too slow, too expensive, or too difficult to measure directly. A **proxy metric** stands in for the target metric during the project.

**Example:** A recommendation algorithm wants to measure "revenue attributable to recommendations over 12 months." That is too slow for A/B testing cycles. A proxy: "average order value within 7 days of a recommendation click."

### Conditions a proxy must satisfy

1. **Strong historical correlation** with the true metric — backed by data, not assumption.
2. **Causal plausibility** — moving the proxy should plausibly move the outcome.
3. **Timely sensitivity** — it moves fast enough to be informative within the project timeline.

### Proxy pitfalls

- **Proxy paradox:** the proxy can move in the opposite direction of the true metric. Optimizing click-through rate can hurt revenue if it attracts low-intent traffic.
- **Relationship drift:** the proxy-to-target correlation established in historical data may not hold experimentally. Meta's Analytics team found three systematic violation types: unconfoundedness (confounding variables, not the intervention, cause the association), exclusivity (the experiment bypasses the proxy pathway entirely), and comparability (observational and experimental contexts differ structurally). ([Analytics at Meta, Medium](https://medium.com/@AnalyticsAtMeta/dont-be-seduced-by-the-allure-a-guide-for-how-not-to-use-proxy-metrics-in-experiments-9530caa0eb7c))
- **Use proxies as directional signals ("fire alarms"), not precise measurements ("air quality monitors").** Meta's recommendation: use proxies for causal exploration, not causal estimation. Validate proxy decisions against ground truth on a rolling basis.

---

## 7. Metric trees and decomposition

A single North Star metric is useful for alignment but not for diagnosis. A **metric tree** decomposes the North Star into sub-metrics so you can identify *which lever to pull*.

```
North Star: Monthly Recurring Revenue
├── New MRR
│   ├── New customers acquired
│   └── Average contract value (new)
├── Expansion MRR
│   ├── Upgrade conversion rate
│   └── Average expansion per upgrade
└── Churned MRR
    ├── Churn rate
    └── Average contract value (churned)
```

Each child metric is linked to its parent by a causal relationship (multiplication, subtraction, or a defined process). ([mixpanel.com — metric trees](https://mixpanel.com/blog/metric-trees-benefits-guide/), [kpitree.co](https://kpitree.co/guides/core-concepts/north-star-metric))

**Why this matters during problem framing:** the tree makes explicit which sub-metric the analysis will move. If the analysis scope is only "churn prediction," success should be measured on churn rate and churned MRR, not on total MRR — which can rise or fall for reasons entirely outside the scope.

---

## 8. Vanity metrics vs actionable metrics

**Vanity metrics** look impressive but cannot inform decisions. The test: "Can this number, taken alone, tell us what to do next?" If the answer is no, it is a vanity metric.

| Vanity | Actionable equivalent |
|--------|-----------------------|
| Total registered users | Monthly active users / 30-day retention rate |
| Page views | Pages per session + exit rate |
| Number of reports produced | Reports acted upon / decisions attributably changed |
| Model training accuracy | Precision / recall at production threshold on held-out data |

From Tableau's definition: vanity metrics "may look impressive on the outside but don't reflect actual business performance" and "can't be used to make informed strategies." ([Tableau](https://www.tableau.com/learn/articles/vanity-metrics))

---

## 9. Goodhart's Law and gaming

"When a measure becomes a target, it ceases to be a good measure." — Charles Goodhart, 1975.

Three mechanisms drive this:

1. **Self-determination theory:** extrinsic targets crowd out intrinsic motivation, shifting behavior from doing good work to hitting the number.
2. **Campbell's Law:** high-stakes metrics attract gaming; the higher the stakes, the more distorted the process.
3. **Principal-agent problems:** information asymmetry lets measured parties exploit gaps between what is measured and what is valued.

Real examples: Soviet nail factories produced unusably heavy nails when targets were by weight; unusably tiny nails when by count. Wells Fargo employees opened unauthorized accounts to hit account-creation targets. ([kpitree.co — Goodhart's Law](https://kpitree.co/guides/frameworks/goodharts-law))

### Design countermeasures

| Countermeasure | How it works |
|----------------|-------------|
| Pair quantity with quality | "Leads generated" + "Lead-to-opportunity rate" — inflating volume exposes as quality decline |
| Combine leading and lagging | Divergence signals gaming |
| Measure outcomes, not outputs | Gaming an outcome requires actually changing the thing that matters |
| Rotate metrics periodically | Prevents ossification of gaming strategies |
| Use metric trees | Causal connections make local gaming immediately visible in neighbor metrics |
| Add qualitative checks | Case studies and retrospectives expose what numbers cannot |

---

## 10. How many metrics to define

More metrics is not better. A useful starting point:

| Project type | North Star | Supporting KPIs | Leading indicators |
|---|---|---|---|
| Exploratory / descriptive | 1 outcome metric | 2–3 | 2–3 |
| Predictive model | 1 primary model + business metric | 3–5 | 2–3 |
| Decision-support analysis | 1 decision metric | 2–4 | 1–2 |

The upper bound is **7–10 KPIs total** across all levels. Beyond that, attention fragments and accountability diffuses. ([spiderstrategies.com](https://www.spiderstrategies.com/blog/strategy-led-kpi-management/))

When a stakeholder asks to add a metric mid-project, ask: "Does this change the definition of success, or is it a diagnostic slice we can add to the appendix?" If the former, restart the alignment process. If the latter, add it as a secondary view without elevating it to KPI status.

### Revising metrics without HARKing

Sometimes a metric defined at project start proves unmeasurable once data collection begins (missing data, wrong granularity, prohibitive cost). Acceptable revision requires:

1. **Document the original metric and the reason it cannot be used** before looking at any results.
2. **Propose the replacement metric** based on data availability, not on which direction the data happens to be moving.
3. **Get stakeholder sign-off on the revision** before proceeding, using the same alignment process as section 11 below.

If you look at preliminary results first and then propose an alternative metric, that is HARKing — even if unintentional.

---

## 11. Stakeholder alignment process

Success metrics are not purely technical choices — they encode whose definition of success gets used. Disagreements that surface at review time should have been surfaced at framing time.

Recommended process at the start of a project:

1. **Convene stakeholders early** — include business owners, domain experts, and data consumers.
2. **Build a success matrix** — map each stakeholder group to their primary concern and the metric that captures it.
3. **Expose conflicts explicitly** — two groups may genuinely want metrics that trade off against each other (e.g., precision vs. recall in a risk model).
4. **Document agreed metrics in writing** — a one-paragraph "success definition" that stakeholders sign off on before data collection begins.
5. **Assign metric ownership** — for each KPI, one person is accountable for data availability and reporting.

([plane.so](https://plane.so/blog/how-to-measure-project-success-metrics-kpis-and-methods), [monday.com — project success](https://monday.com/blog/project-management/how-to-define-project-success/))

---

## 12. Worked example: churn analysis project

**Business question:** Why are enterprise customers churning at higher rates than last year?

### Step 1 — Identify the lagging indicator
Primary outcome metric: Enterprise customer churn rate (% of enterprise accounts that did not renew, measured quarterly).

### Step 2 — Select leading indicators
- Support ticket volume per account (leading: high volume predicts churn)
- Days-since-last-login for key personas (leading: declining engagement predicts churn)
- NPS score from QBR surveys (leading: low NPS predicts churn)

### Step 3 — Decompose with a metric tree
```
Enterprise churn rate
├── Churn by customer segment (industry, ARR band)
├── Churn by cohort (renewal quarter)
└── Churn by health score tier
    ├── Support volume tier
    ├── Engagement tier
    └── NPS tier
```

### Step 4 — SMART the primary metric
"Reduce enterprise churn rate from 14% to under 10% by end of Q4 of the current fiscal year, measured against the enterprise account list as of January 1."

### Step 5 — Set proxy guardrails
The proxy "NPS improvement" is acceptable as an early signal, but must be validated against actual renewal outcomes at end of Q2.

### Step 6 — Document and sign off
The above definition is written into the project brief, reviewed with the VP of Customer Success and the VP of Revenue Operations, and stored in the project wiki before analysis begins.

---

## 13. Common pitfalls checklist

- [ ] Metrics defined after seeing preliminary results (post-hoc metric selection)
- [ ] No lagging indicator — only leading or proxy metrics
- [ ] Metric not tied to a specific data source (unmeasurable in practice)
- [ ] Too many KPIs — more than 7-10 causes analysis paralysis ([spiderstrategies.com](https://www.spiderstrategies.com/blog/strategy-led-kpi-management/))
- [ ] KPI scope broader than the analysis scope (e.g., measuring total revenue when only one product is being analyzed)
- [ ] No time bound on target
- [ ] Vanity metric dressed as a KPI (total users, total reports)
- [ ] Single metric that can be gamed without corresponding quality check
- [ ] Stakeholders never aligned on the definition before data collection started
- [ ] Proxy metric accepted without validation against ground truth

---

## 14. Sources

1. [Business Problem Framing in the Age of Data Analytics — productleadership.com](https://www.productleadership.com/blog/business-problem-framing-in-the-age-of-data-analytics/)
2. [Setting SMART Objectives in Data Analytics — datacalculus.com](https://datacalculus.com/en/knowledge-hub/data-analytics/pre-data-collection-phase/setting-smart-objectives-in-data-analytics/)
3. [KPIs for Data Teams: A Comprehensive Guide — atlan.com](https://atlan.com/kpis-for-data-team/)
4. [KPIs: Understanding Metrics for Business Success — atlassian.com](https://www.atlassian.com/work-management/project-management/project-planning/kpi)
5. [Performance measurement: Leading vs. Lagging Indicators — humanperf.com](https://www.humanperf.com/en/blog/project-management/articles/performance-measurement-indicators)
6. [Leading vs. Lagging Indicators: Measuring with OKRs — xodiac.ca](https://xodiac.ca/blog/articles/leading-vs-lagging-indicators-to-measure-okrs)
7. [Don't Be Seduced by the Allure: Proxy Metrics in Experiments — Analytics at Meta, Medium](https://medium.com/@AnalyticsAtMeta/dont-be-seduced-by-the-allure-a-guide-for-how-not-to-use-proxy-metrics-in-experiments-9530caa0eb7c)
8. [Metric Trees 101: Benefits and Process — mixpanel.com](https://mixpanel.com/blog/metric-trees-benefits-guide/)
9. [North Star Metric: What It Is and How to Find Yours — kpitree.co](https://kpitree.co/guides/core-concepts/north-star-metric)
10. [Goodhart's Law and Metric Design — kpitree.co](https://kpitree.co/guides/frameworks/goodharts-law)
11. [Vanity Metrics: Definition, Identification, Examples — Tableau](https://www.tableau.com/learn/articles/vanity-metrics)
12. [Outcome vs. Output Metrics — spiderstrategies.com](https://www.spiderstrategies.com/blog/outcome-vs-output-metrics/)
13. [How to Measure Project Success: Metrics, KPIs, Methods — plane.so](https://plane.so/blog/how-to-measure-project-success-metrics-kpis-and-methods)
14. [Project Success Criteria: How to Measure Real Outcomes — monday.com](https://monday.com/blog/project-management/how-to-define-project-success/)
15. [Strategy-Led KPI Management: Align Metrics With Goals — spiderstrategies.com](https://www.spiderstrategies.com/blog/strategy-led-kpi-management/)
