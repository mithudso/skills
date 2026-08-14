---
description: >-
  Foundations & theory of data analysis hub (da family) — conceptual & mathematical grounding beneath any technique. TRIGGER: definitions/scope (analysis vs analytics vs data science vs statistics, quantitative vs qualitative); measurement (levels nominal/ordinal/interval/ratio, reliability/validity, operationalization, constructs); probability (random variables, PMF/PDF, named distributions, joint/marginal/conditional, Bayes, LLN, CLT); statistical inference foundations (population vs sample, sampling distributions, standard error, frequentist vs Bayesian); information theory (entropy, mutual information); epistemology (correlation vs causation, inductive vs deductive, reproducibility/replicability). SKIP: lifecycle → da-2-data-analysis-lifecycle; acquisition/sampling → da-3-data-acquisition-sampling; methods/ML → da-analytical-methods; pipelines/platform → da-data-engineering-platform; viz/comms → da-applied-and-communication.
name: da-1-foundations-theory
version: 2.1.0
updated: "2026-05-31"
category: custom
tags:
  - data-analysis
  - statistics
  - probability
  - measurement-theory
  - statistical-inference
  - information-theory
  - epistemology
related_skills:
  - da-2-data-analysis-lifecycle
  - da-3-data-acquisition-sampling
  - da-analytical-methods
  - da-data-engineering-platform
  - da-applied-and-communication
  - ai-agent-engineering
  - programming-languages
  - mongodb-expert
whenToUse:
  - "What is data analysis and how does it differ from analytics or data science?"
  - "Explain the levels of measurement — nominal, ordinal, interval, ratio"
  - "What is a probability distribution and how do I choose one?"
  - "Explain Bayes' theorem / conditional probability from first principles"
  - "What is the central limit theorem and why does it matter?"
  - "What is Shannon entropy or mutual information?"
  - "Why can't I compute a mean on ordinal / Likert data?"
  - "Explain frequentist vs. Bayesian inference — conceptually"
  - "What is the difference between correlation and causation?"
  - "How is reproducibility different from replicability in data work?"
  - "What does 'operationalizing a construct' mean in measurement theory?"
  - "Frame or scope a data analysis — what kind of analysis is this, conceptually?"
whenNotToUse: >
  Defer to da-2-data-analysis-lifecycle when the question is about the PROCESS or
  workflow of running an analysis (CRISP-DM, KDD, problem framing, lifecycle phases,
  iteration strategy). Defer to da-3-data-acquisition-sampling for data COLLECTION
  methods, sampling design, and survey methods. Defer to da-analytical-methods when
  the question is about APPLYING a statistical or ML method (fitting a model, running
  an A/B test, EDA on a real dataset). Defer to da-applied-and-communication for
  visualization, reporting, and stakeholder communication. Defer to
  da-data-engineering-platform for pipeline, warehouse, or infrastructure work.
  Defer to ai-agent-engineering for neural-net / ML model architecture theory; defer
  to programming-languages for Python/R implementation of statistical concepts; defer
  to mongodb-expert for database-level sampling or statistics.
---
# Data Analysis: Foundations & Theory

Conceptual grounding for data analysis. Answers *before-you-pick-a-tool* questions: what is data analysis, how it relates to neighboring fields, what kind of analysis needed, what data supports given how measured, what assumptions underneath. No technique execution — scopes and frames only.

## When not to use this hub

| If the question is about… | Use instead |
| --- | --- |
| Process, workflow, CRISP-DM, problem framing, lifecycle phases | `da-2-data-analysis-lifecycle` |
| Data collection methods, sampling design, survey methods | `da-3-data-acquisition-sampling` |
| Applying a method (fitting a model, running EDA, A/B tests) | `da-analytical-methods` |
| Visualization, reporting, stakeholder communication | `da-applied-and-communication` |
| Pipelines, warehouses, platform, governance | `da-data-engineering-platform` |
| Neural-net or ML model architecture theory | `ai-agent-engineering` |
| Python/R implementation of statistical concepts | `programming-languages` |
| Database-level statistics or sampling | `mongodb-expert` |

## Sub-skill routing table

Hub consolidates 37 foundations sub-skills as on-demand references — match task to table and **Read listed `references/<name>.md` before answering deep questions**. Overview below enough for framing/scoping; load reference when question needs depth.

| Sub-topic | When to load | Reference file |
| --- | --- | --- |
| `da-1-1-definitions-scope` | Defines what "data analysis" means and where its boundaries sit | `references/da-1-1-definitions-scope.md` |
| `da-1-1-1-data-analysis-vs-analytics-vs-data` | Disambiguates four overlapping terms — data analysis, data analytics, data science, statistics | `references/da-1-1-1-data-analysis-vs-analytics-vs-data.md` |
| `da-1-1-2-analysis-vs-synthesis` | Distinguishes analysis (breaking a whole into parts) from synthesis | `references/da-1-1-2-analysis-vs-synthesis.md` |
| `da-1-1-3-quantitative-vs-qualitative-analysis` | Foundational distinction between quantitative and qualitative analysis | `references/da-1-1-3-quantitative-vs-qualitative-analysis.md` |
| `da-1-2-measurement-theory` | Measurement theory as the foundational layer — how numbers map to reality | `references/da-1-2-measurement-theory.md` |
| `da-1-2-1-levels-of-measurement` | Stevens' levels (scales) of measurement — nominal, ordinal, interval, ratio | `references/da-1-2-1-levels-of-measurement.md` |
| `da-1-2-1-1-nominal` | The NOMINAL level of measurement (Stevens 1946) — the lowest scale type | `references/da-1-2-1-1-nominal.md` |
| `da-1-2-1-2-ordinal` | The ordinal level of measurement — ranked categories | `references/da-1-2-1-2-ordinal.md` |
| `da-1-2-1-3-interval` | The interval level of measurement — equal differences, no true zero | `references/da-1-2-1-3-interval.md` |
| `da-1-2-1-4-ratio` | The ratio level of measurement — equal differences plus a true zero | `references/da-1-2-1-4-ratio.md` |
| `da-1-2-2-discrete-vs-continuous-variables` | Distinguishes discrete from continuous variables as a measurement-theory concept | `references/da-1-2-2-discrete-vs-continuous-variables.md` |
| `da-1-2-3-reliability-validity` | Reliability and validity: whether a measure is consistent and measures what it claims | `references/da-1-2-3-reliability-validity.md` |
| `da-1-2-4-operationalization-constructs` | Turning an abstract construct into something measurable | `references/da-1-2-4-operationalization-constructs.md` |
| `da-1-3-probability-theory` | Foundational probability theory as the mathematical basis for analysis | `references/da-1-3-probability-theory.md` |
| `da-1-3-1-random-variables` | The formal treatment of random variables | `references/da-1-3-1-random-variables.md` |
| `da-1-3-2-probability-mass-density-functions` | Probability mass functions (PMF) and probability density functions (PDF) | `references/da-1-3-2-probability-mass-density-functions.md` |
| `da-1-3-3-probability-distributions` | What a probability DISTRIBUTION is and how distributions are characterized | `references/da-1-3-3-probability-distributions.md` |
| `da-1-3-3-1-normal` | The Normal (Gaussian) distribution as a named continuous distribution | `references/da-1-3-3-1-normal.md` |
| `da-1-3-3-2-binomial` | The binomial distribution as a named probability distribution | `references/da-1-3-3-2-binomial.md` |
| `da-1-3-3-3-poisson` | The Poisson distribution as a named probability distribution | `references/da-1-3-3-3-poisson.md` |
| `da-1-3-3-4-exponential` | The exponential distribution as a named continuous distribution | `references/da-1-3-3-4-exponential.md` |
| `da-1-3-3-5-uniform` | The uniform distribution as a named member of the distribution family | `references/da-1-3-3-5-uniform.md` |
| `da-1-3-4-joint-marginal-conditional-probability` | Joint, marginal, and conditional probability and their relationships | `references/da-1-3-4-joint-marginal-conditional-probability.md` |
| `da-1-3-5-bayes-theorem` | Bayes' theorem as a result of probability theory | `references/da-1-3-5-bayes-theorem.md` |
| `da-1-3-6-law-of-large-numbers` | The Law of Large Numbers (LLN) | `references/da-1-3-6-law-of-large-numbers.md` |
| `da-1-3-7-central-limit-theorem` | The central limit theorem (CLT) | `references/da-1-3-7-central-limit-theorem.md` |
| `da-1-3-8-expectation-variance-covariance` | Expectation (expected value / mean), variance, and covariance | `references/da-1-3-8-expectation-variance-covariance.md` |
| `da-1-4-statistical-inference-foundations` | Foundations of statistical inference — how you reason from sample to population | `references/da-1-4-statistical-inference-foundations.md` |
| `da-1-4-1-population-vs-sample` | The distinction between a population and a sample | `references/da-1-4-1-population-vs-sample.md` |
| `da-1-4-2-sampling-distributions-standard-error` | The sampling distribution of a statistic and the standard error | `references/da-1-4-2-sampling-distributions-standard-error.md` |
| `da-1-4-3-frequentist-vs-bayesian-paradigms` | The two competing schools of statistical inference — frequentist and Bayesian | `references/da-1-4-3-frequentist-vs-bayesian-paradigms.md` |
| `da-1-4-4-estimation-theory` | Point-estimation theory and its three criteria | `references/da-1-4-4-estimation-theory.md` |
| `da-1-5-information-theory` | Information-theoretic measures — Shannon entropy, mutual information | `references/da-1-5-information-theory.md` |
| `da-1-6-epistemology-of-data` | The epistemology of data — how data come to count as knowledge | `references/da-1-6-epistemology-of-data.md` |
| `da-1-6-1-correlation-vs-causation` | Why a statistical association is not causation | `references/da-1-6-1-correlation-vs-causation.md` |
| `da-1-6-2-inductive-vs-deductive-reasoning` | How inductive, deductive, and abductive reasoning differ | `references/da-1-6-2-inductive-vs-deductive-reasoning.md` |
| `da-1-6-3-reproducibility-replicability` | Reproducibility and replicability as epistemic standards | `references/da-1-6-3-reproducibility-replicability.md` |

## 1. What data analysis is (and its scope)

Data analysis = systematic process of inspecting, cleaning, transforming, interpreting data to extract useful info, support conclusions, aid decisions. Bounded and goal-directed: analyst works specific question from decision-maker against cleaned dataset narrowed to relevant info ([Julius AI](https://julius.ai/articles/data-analysis-vs-data-science)).

Four neighboring terms — overlap but not synonyms:

- **Data analysis** — act of evaluating data to answer defined question. *Process verb* ([RudderStack](https://www.rudderstack.com/learn/data-analytics/data-analytics-vs-data-science/)).
- **Data analytics** — broader practice/field around analyzing data: capture, process, organize to surface actionable insights, typically past/present. Stats tend elementary (sums, counts, averages, percentiles) ([IBM](https://www.ibm.com/think/topics/data-science-vs-data-analytics)).
- **Data science** — wider, multidisciplinary field that *includes* analytics but extends into ML, forecasting, large-scale data engineering across structured/unstructured; oriented toward predictive models ([Qlik](https://www.qlik.com/us/data-analytics/data-science-vs-data-analytics)).
- **Statistics** — math discipline of collecting, describing, drawing inferences from data under uncertainty. Data analysis *uses* statistics as toolkit; not coextensive. Tukey (1961) defined data analysis as procedures for examining datasets *plus* statistics techniques — analysis larger than math statistics alone ([Wikipedia: EDA](https://en.wikipedia.org/wiki/Exploratory_data_analysis)).

Rule of thumb: *data analysis = activity; analytics = field around it; data science extends toward modeling/engineering; statistics supplies inferential math.* Sources disagree on exact boundaries — state which definition you use.

### Analysis vs. synthesis

Analysis breaks whole into parts to understand structure/relationships; synthesis recombines parts into integrated whole or recommendation. Analysis effort usually *analyzes* data, then *synthesizes* findings into narrative or decision. Naming which mode prevents presenting raw decomposition as conclusion.

### Quantitative vs. qualitative

- **Quantitative analysis** — operates on numeric measurements, supports statistical computation. Interval/ratio scales are quantitative/continuous ([Statistics By Jim](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/)).
- **Qualitative analysis** — operates on categorical/textual/observational data (themes, labels, narratives). Nominal/ordinal scales are categorical, align with qualitative framing ([Statistics By Jim](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/)).

Split not purely about data type — also about method (counting/modeling vs. interpreting meaning). Many real efforts mixed-method.

## 2. The four families of analysis

Widely cited progression (Gartner's "Analytic Ascendancy") sorts analysis by question answered and increasing difficulty/value ([Domo](https://www.domo.com/learn/article/data-analytics-types), [EAG Inc.](https://eaginc.com/understanding-data-analytics/)):

| Type | Question | Typical methods |
|---|---|---|
| **Descriptive** | What happened? | aggregation, reporting, summary statistics |
| **Diagnostic** | Why did it happen? | root-cause analysis, correlation, drill-down |
| **Predictive** | What is likely to happen? | forecasting, regression, ML, probability scores |
| **Prescriptive** | What should we do? | optimization, decision rules, simulation |

Practical notes:

- Four form maturity progression, not strictly sequential per project — pick family matching decision ([Learning Discourses](https://learningdiscourses.com/subdiscourse/analytics-maturity-gartner-analytic-ascendancy-model/)).
- Prescriptive work needs reliable predictive models plus domain constraints; few orgs reach it ([Dataforest](https://dataforest.ai/blog/analytics-maturity-model)).
- "Diagnostic" in Gartner table maps loosely to **exploratory** in statistical tradition (see §3). Don't conflate marketing taxonomy with Tukey's exploratory/confirmatory split — different questions.

## 3. Exploratory vs. confirmatory analysis

Foundational theoretical distinction, separate from Gartner taxonomy:

- **Exploratory Data Analysis (EDA)** — Tukey's approach (*Exploratory Data Analysis*, 1977) for summarizing dataset's main characteristics, often with graphics, to *generate* hypotheses and check assumptions before formal testing ([Wikipedia: EDA](https://en.wikipedia.org/wiki/Exploratory_data_analysis)). Signature techniques: box-and-whisker plots, stem-and-leaf displays, histograms, scatter plots ([Sage](https://methods.sagepub.com/ency/edvol/the-sage-encyclopedia-of-social-science-research-methods/chpt/exploratory-data-analysis)).
- **Confirmatory Data Analysis (CDA)** — classical hypothesis testing: pick model *before* examining data, then assess how precisely sample statistics let you infer population parameters ([Wikipedia: EDA](https://en.wikipedia.org/wiki/Exploratory_data_analysis)).

Tukey argued 20th-century statistics over-weighted confirmatory testing, under-weighted using data to suggest hypotheses, but held researchers need *both* — "cannot get along without confirmatory data analysis but need not start with it" ([Wikipedia: EDA](https://en.wikipedia.org/wiki/Exploratory_data_analysis)).

**Critical pitfall:** running exploratory *and* confirmatory on same data introduces bias — hypotheses from data then tested on same data inflate apparent significance (double-dipping) ([Wikipedia: EDA](https://en.wikipedia.org/wiki/Exploratory_data_analysis)). Reserve holdout or fresh sample for confirmation.

For EDA on real dataset, route to `da-analytical-methods` (`references/da-5-exploratory-data-analysis.md`). For process frameworks sequencing exploratory/confirmatory phases, route to `da-2-data-analysis-lifecycle`.

## 4. Measurement theory and levels of measurement

How variable was measured constrains which operations/statistics are valid. Stevens introduced four levels in *On the Theory of Scales of Measurement* (Science, 1946) ([Wikipedia: Level of measurement](https://en.wikipedia.org/wiki/Level_of_measurement)):

| Level | Distinguishes | Permissible central tendency | Example |
|---|---|---|---|
| **Nominal** | categories only (=, ≠) | mode | dog/cat/rabbit |
| **Ordinal** | rank order | median (no standard deviation under Stevens) | 1st/2nd/3rd |
| **Interval** | equal differences, no true zero | mean, median, mode | Celsius, calendar dates |
| **Ratio** | equal differences + true zero (ratios meaningful) | adds geometric/harmonic means | mass, length, duration |

Each level adds info on top of previous. Nominal/ordinal = categorical/qualitative; interval/ratio = continuous/quantitative ([Statistics By Jim](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/)). Payoff: variable's level tells which descriptive statistics and inferential tests are appropriate ([Statistics By Jim](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/)).

**Common pitfalls:**

- Computing mean of ordinal codes (e.g., averaging 1–5 Likert) not licensed by Stevens — median is safer.
- Treating arbitrary numeric label as quantitative. Lord's 1953 satire about football jersey numbers warns against arithmetic on nominal codes ([Wikipedia: Level of measurement](https://en.wikipedia.org/wiki/Level_of_measurement)).
- Forgetting interval scales lack true zero, so ratios ("twice as hot" in Celsius) meaningless.

**Know the controversy.** Stevens' typology foundational but contested. Velleman & Wilkinson (1993) called it "misleading" for choosing statistics; Luce (1997) noted working measurement theorists don't accept Stevens' broad definition ([Wikipedia: Level of measurement](https://en.wikipedia.org/wiki/Level_of_measurement)). Treat measurement level as useful first filter, not iron law — experts genuinely disagree on edges.

## 5. The role of theory and assumptions

Data doesn't interpret itself. Every analysis rides on assumptions: sample represents population of interest, measurements mean what labels claim, model preconditions (independence, distributional form, scale type) hold. Two practices follow:

- **State assumptions explicitly**, tie each to analysis family and measurement level in play. EDA is partly assumption-checking before confirmatory work ([Wikipedia: EDA](https://en.wikipedia.org/wiki/Exploratory_data_analysis)).
- **Match method to question and data.** Gartner family (§2), exploratory/confirmatory mode (§3), measurement level (§4) together determine what analysis can legitimately claim. Misalignment — e.g., prescriptive recommendation built on descriptive data, or mean on ordinal categories — most common foundational error.

## Quick decision checklist

1. **Term check** — analysis, analytics, data science, or statistics? State definition used.
2. **Family** — descriptive, diagnostic, predictive, or prescriptive? Match to decision.
3. **Mode** — exploratory (generate hypotheses) or confirmatory (test)? Don't mix on same data.
4. **Lifecycle phase** — which phase and what next? (Route to `da-2-data-analysis-lifecycle` for CRISP-DM, KDD, other frameworks.)
5. **Measurement** — what level is each variable; which statistics are licensed?
6. **Assumptions** — what must be true for conclusion to hold; have I checked?

## Sources

1. [Julius AI — Data Analysis vs. Data Science](https://julius.ai/articles/data-analysis-vs-data-science) — scope and goal-directed nature of data analysis.
2. [IBM — Data science vs data analytics](https://www.ibm.com/think/topics/data-science-vs-data-analytics) — field distinctions and statistical depth.
3. [Qlik — Data Science vs Data Analytics](https://www.qlik.com/us/data-analytics/data-science-vs-data-analytics) — predictive-modeling scope of data science.
4. [RudderStack — Data Analytics vs Data Science](https://www.rudderstack.com/learn/data-analytics/data-analytics-vs-data-science/) — analysis as process vs. analytics as practice.
5. [Domo — The 4 Types of Data Analytics](https://www.domo.com/learn/article/data-analytics-types) — descriptive/diagnostic/predictive/prescriptive.
6. [EAG Inc. — Understanding the 4 Types of Data Analytics](https://eaginc.com/understanding-data-analytics/) — questions each type answers.
7. [Learning Discourses — Gartner Analytic Ascendancy Model](https://learningdiscourses.com/subdiscourse/analytics-maturity-gartner-analytic-ascendancy-model/) — maturity framing.
8. [Dataforest — Analytics Maturity Model](https://dataforest.ai/blog/analytics-maturity-model) — prescriptive adoption rates (caveated).
9. [Wikipedia — Exploratory data analysis](https://en.wikipedia.org/wiki/Exploratory_data_analysis) — Tukey, EDA vs. CDA, techniques, bias pitfall.
10. [Sage — Exploratory Data Analysis (encyclopedia)](https://methods.sagepub.com/ency/edvol/the-sage-encyclopedia-of-social-science-research-methods/chpt/exploratory-data-analysis) — EDA techniques.
11. [Wikipedia — Level of measurement](https://en.wikipedia.org/wiki/Level_of_measurement) — Stevens 1946, four levels, permissible statistics, Lord/Velleman & Wilkinson/Luce criticisms.
12. [Statistics By Jim — Nominal/Ordinal/Interval/Ratio](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/) — categorical vs. continuous, statistic selection.

<!-- cross-hub-map -->
## Cross-hub map — where every data-analytics topic lives

Family split across hubs. If task's deep material **not** in this hub's Sub-skill routing table, it's reference file under sibling hub below — **activate that hub or `Read` its `references/<name>.md` directly**. Every former standalone skill now reference under one of these hubs (nothing deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `da-1-foundations-theory` | Data Analysis Foundations & Theory (hub) | `references/da-1-1-definitions-scope.md`, `references/da-1-1-1-data-analysis-vs-analytics-vs-data.md`, `references/da-1-1-2-analysis-vs-synthesis.md`, `references/da-1-1-3-quantitative-vs-qualitative-analysis.md`, … |
| `da-2-data-analysis-lifecycle` | Data Analysis Lifecycle & Process (hub) | `references/da-2-1-process-frameworks.md`, `references/da-2-1-1-crisp-dm.md`, `references/da-2-1-2-kdd.md`, `references/da-2-1-3-semma.md`, … |
| `da-3-data-acquisition-sampling` | Data Acquisition, Collection & Sampling (hub) | `references/da-3-1-data-sources.md`, `references/da-3-1-1-primary-vs-secondary.md`, `references/da-3-1-2-internal-vs-external.md`, `references/da-3-1-3-structured-semi-structured-unstructured.md`, … |
| `da-analytical-methods` | Data Analytical Methods (cleaning, EDA, modeling, ML, causal, time-series) | `references/da-4-data-cleaning-preparation.md`, `references/da-5-exploratory-data-analysis.md`, `references/da-6-statistical-modeling.md`, `references/da-7-machine-learning.md`, … |
| `da-data-engineering-platform` | Data Engineering & Analytics Platform (pipelines, OLAP, modeling, governance) | `references/da-10-tools-and-languages.md`, `references/da-13-data-engineering-and-pipelines.md`, `references/da-14-streaming-analytics.md`, `references/da-18-semantic-layer-headless-bi.md`, … |
| `da-applied-and-communication` | Applied Analytics, Visualization, Communication & Ethics | `references/da-8-data-visualization.md`, `references/da-9-reporting-communication.md`, `references/da-11-ethics-and-privacy.md`, `references/da-21-product-analytics.md`, … |
| `data-analytics` | Family ROUTER — entry point for all da-* sub-hubs | (this file's parent hub) |