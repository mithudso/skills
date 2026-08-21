<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-1-definitions-scope` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-1-definitions-scope
description: >-
  Defines what "data analysis" means and where its boundaries sit, as the
  foundational/theory entry point for the Data Analysis taxonomy. Covers the
  canonical process definition (inspect, clean, transform, model), the scope of
  the field, Tukey's historical framing, the data-analysis-as-philosophy view
  (NIST/EDA), and the descriptive/diagnostic/predictive/prescriptive scope
  ladder used to bound an analysis effort.
  TRIGGER: someone asks "what is data analysis", "what counts as data analysis
  vs not", how to scope or frame a data-analysis task, the definition or
  boundaries of the field, what the stages/steps of analysis are at a
  definitional level, or where the four analytics types fit conceptually.
  SKIP: distinguishing data analysis from analytics/data science/statistics as
  disciplines (defer to da-1-1-1); analysis vs. synthesis (defer to da-1-1-2);
  quantitative vs. qualitative (defer to da-1-1-3); measurement theory and
  scales (defer to da-1-2-*); and any hands-on technique, tool, or statistical
  method (those live elsewhere in the taxonomy, not in this definitional node).
---

# Data Analysis — Definitions & Scope

Foundational, definitional node. This skill answers *what data analysis is* and
*what falls inside or outside its scope* — not how to perform any specific
technique. Use it to frame and bound a task before reaching for methods.

## The canonical definition

Data analysis is "the process of inspecting, cleansing, transforming, and
modeling data with the goal of discovering useful information, informing
conclusions, and supporting decision-making" (Wikipedia, "Data analysis",
https://en.wikipedia.org/wiki/Data_analysis). This is the most widely cited
working definition and it encodes four observations that bound the scope:

1. It is a **process**, not a single act — a sequence of stages, not one tool.
2. Its raw material is **data** (structured or unstructured, from one or many
   sources).
3. Its operations are **inspect → clean → transform → model**.
4. Its purpose is **decision support** — discovering information, drawing
   conclusions, informing a choice. An activity with data that has no such aim
   (e.g. pure storage or transport) is outside the scope.

The earliest formal scoping comes from John Tukey (1961), who defined data
analysis as "procedures for analyzing data, techniques for interpreting the
results of such procedures, ways of planning the gathering of data to make its
analysis easier, more precise or more accurate, and all the machinery and
results of (mathematical) statistics which apply to analyzing data"
(cited in Wikipedia, "Exploratory data analysis",
https://en.wikipedia.org/wiki/Exploratory_data_analysis). Note how broad
Tukey's scope is: it explicitly pulls *planning the data collection* and
*interpreting results* inside the boundary, not just the computation.

## Scope: what is in and what is out

| Inside the boundary | Outside (other nodes / not analysis) |
|---|---|
| Inspecting and profiling data | Raw data collection logistics (planning is in scope per Tukey; mechanics often treated separately) |
| Cleaning / cleansing (errors, duplicates, outliers) | Long-term data storage and transport |
| Transforming / restructuring into analyzable form | Building production data pipelines as engineering |
| Modeling to surface patterns and test hypotheses | The downstream business *decision* itself |
| Interpreting and communicating results | Defining "analytics" / "data science" as fields (→ da-1-1-1) |

Data analysis "has multiple facets and approaches, encompassing diverse
techniques under a variety of names, and is used in different business,
science, and social science domains" (Wikipedia, "Data analysis"). The scope is
therefore *domain-independent*: the definition holds whether the data is
financial, clinical, or experimental.

## Two complementary framings of "what data analysis is"

There are two long-standing views of the field, and a good answer names both.

**(a) Data analysis as a process (the procedural view).** The definition above.
Stages, roughly: define objective/question → collect → clean → transform →
explore → model/analyze → interpret/visualize → communicate
(DataCamp, "What is Data Analysis?",
https://www.datacamp.com/blog/what-is-data-analysis-expert-guide). The
boundaries of a task are set by deciding which of these stages it includes.

**(b) Data analysis as an attitude/philosophy (the EDA view).** NIST, following
Tukey, frames exploratory data analysis not as a set of techniques but as "an
approach/philosophy for data analysis" — "an attitude about how data analysis
should be conducted" that "postpones the usual assumptions about what kind of
model the data follow with the more direct approach of allowing the data itself
to reveal its underlying structure" (NIST/SEMATECH e-Handbook of Statistical
Methods, §1.1.1, https://www.itl.nist.gov/div898/handbook/eda/section1/eda11.htm).
NIST lists the goals that define the scope of this view: maximize insight into a
data set; uncover underlying structure; extract important variables; detect
outliers and anomalies; test underlying assumptions; develop parsimonious
models; and determine optimal factor settings (same source). This matters for
scoping because it tells you that "analysis" legitimately includes open-ended
exploration *before* any hypothesis or model is fixed.

## The scope ladder: four types of analysis

A standard way to bound how *far* an analysis goes is the four-type ladder, each
rung answering a different question and demanding more capability than the last:

- **Descriptive** — "What happened?" Summarizes past/current data to reveal
  patterns and trends.
- **Diagnostic** — "Why did it happen?" Investigates causes (root-cause analysis
  is one technique here).
- **Predictive** — "What is likely to happen?" Uses historical data plus
  statistical/ML models to forecast.
- **Prescriptive** — "What should we do about it?" Recommends actions, often via
  optimization or simulation.

(insightsoftware, "Comparing Descriptive, Predictive, Prescriptive, and
Diagnostic Analytics",
https://insightsoftware.com/blog/comparing-descriptive-predictive-prescriptive-and-diagnostic-analytics/;
Adobe Business, "Descriptive, predictive, diagnostic, and prescriptive analytics
explained",
https://business.adobe.com/blog/basics/descriptive-predictive-prescriptive-analytics-explained.)
Each rung builds on the previous one, moving from describing what happened to
recommending and influencing what will happen. When scoping a request, pin which
rung is being asked for — the answer changes the data, the methods, and the
effort required.

## When/why you reach for this node

- A stakeholder uses "data analysis" loosely and you need to set expectations
  about what the work will and will not include.
- You are writing a scope statement, proposal, or ticket and must state the
  boundary of an analysis effort.
- You need to decide whether a request is descriptive vs. predictive (etc.)
  before estimating effort.
- You are teaching or documenting the field from first principles.

## Pitfalls

- **Collapsing process and philosophy.** Treating data analysis as only a
  pipeline of steps misses the EDA/attitude view, where letting the data speak
  before assuming a model is itself part of the definition (NIST §1.1.1).
- **Confusing the type ladder with a maturity score.** Descriptive analysis is
  not "lesser"; many decisions only need "what happened." Choose the rung the
  question demands, not the highest rung available.
- **Scoping out planning.** Tukey explicitly folds *planning the data gathering*
  into data analysis; excluding it by default can leave a gap no one owns.
- **Equating "analysis" with "the decision."** The definition stops at informing
  and supporting the decision; making the decision is downstream and out of
  scope.
- **Treating the definition as domain-specific.** The canonical definition is
  domain-independent; do not narrow it to "business analytics" or "research
  statistics" alone.
- **Drifting into the sibling nodes.** Comparing data analysis to data science,
  statistics, or analytics-as-a-field belongs in da-1-1-1; analysis vs.
  synthesis in da-1-1-2; quantitative vs. qualitative in da-1-1-3. Keep this
  node about the definition and outer boundary of "data analysis" itself.

## Sources

- Wikipedia, "Data analysis" — canonical process definition and scope.
  https://en.wikipedia.org/wiki/Data_analysis
- Wikipedia, "Exploratory data analysis" — Tukey's 1961/1977 definitions and
  history. https://en.wikipedia.org/wiki/Exploratory_data_analysis
- NIST/SEMATECH e-Handbook of Statistical Methods, §1.1.1 "What is EDA?" —
  data-analysis-as-philosophy and the seven EDA goals.
  https://www.itl.nist.gov/div898/handbook/eda/section1/eda11.htm
- DataCamp, "What is Data Analysis? An Expert Guide With Examples" — process
  stages. https://www.datacamp.com/blog/what-is-data-analysis-expert-guide
- insightsoftware, "Comparing Descriptive, Predictive, Prescriptive, and
  Diagnostic Analytics" — the four-type scope ladder.
  https://insightsoftware.com/blog/comparing-descriptive-predictive-prescriptive-and-diagnostic-analytics/
- Adobe Business, "Descriptive, predictive, diagnostic, and prescriptive
  analytics explained" — type definitions and progression.
  https://business.adobe.com/blog/basics/descriptive-predictive-prescriptive-analytics-explained
