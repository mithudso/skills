<!-- hub-reference-banner -->
> **Reference file — part of the `da-2-data-analysis-lifecycle` hub.** Formerly the standalone `da-2-2-2-hypothesis-formulation` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-2-2-2-hypothesis-formulation
description: |
  Expert knowledge on hypothesis formulation as practiced in the data analysis lifecycle
  under the Problem Framing phase. Covers how to construct null/alternative hypotheses,
  choose directional vs. non-directional forms, ground hypotheses in prior knowledge,
  and distinguish exploratory (EDA) hypothesis generation from confirmatory (CDA)
  hypothesis testing.

  TRIGGER: Use when the user asks how to write, formulate, or structure a hypothesis
  before running a statistical test or data analysis; when distinguishing H₀ from H₁;
  when choosing one-tailed vs. two-tailed tests; when translating a business or
  research question into a testable statistical claim; or when reviewing whether a
  hypothesis is well-formed, specific, and testable within the Problem Framing step
  of a data analysis lifecycle (CRISP-DM, TDSP, OSEMN, etc.).

  SKIP: If the user is asking about executing a hypothesis test (computing p-values,
  effect sizes, power calculations, or interpreting test results) — that belongs to
  the hypothesis-testing execution step, not formulation. Skip also if the question is
  about causal model discovery (Bayesian networks, DAG construction) or about
  predictive/ML model selection — those are downstream. Skip if the question is
  about business/research question definition from scratch (da-2-2-1) or the broader
  problem-framing container (da-2-2). Defer general statistical inference foundations
  to da-1-4 skills.
---

# Hypothesis Formulation

**Taxonomy position:** Data Analysis > Data Analysis Lifecycle (Process) > Problem Framing > Hypothesis formulation

Hypothesis formulation is the step in which an analyst translates a research or business question into one or more concrete, testable statistical claims before any confirmatory analysis begins. Done well, it anchors every downstream decision — test selection, data collection scope, sample size, significance threshold, and reporting language. Done poorly, it introduces circular reasoning, inflated false-positive rates, and findings that cannot be reproduced.

---

## 1. What a hypothesis is (and is not)

A hypothesis is **a declarative statement predicting the relationship between variables, expressed in a form that can be evaluated against data and that can in principle be found false** [Statsig, 2024]. It is not a question, not a vague intention, and not a post-hoc rationalization.

Two forms are always paired:

| Component | Symbol | Meaning |
|---|---|---|
| **Null hypothesis** | H₀ | The default position: no effect, no difference, no relationship in the population. Always expressed with an equality sign (=, ≤, ≥). |
| **Alternative hypothesis** | H₁ or Hₐ | The claim the analyst expects or wishes to detect. Always expressed with an inequality (≠, <, >). |

Together, H₀ and H₁ must be **mutually exclusive and collectively exhaustive** — they cover every possible outcome for the parameter of interest [Scribbr, 2024; OpenLearn/Open University, 2024].

**Example — salary analysis:**
- H₀: µ = £26,000 (population mean salary equals £26,000)
- H₁: µ ≠ £26,000 (population mean salary differs from £26,000)

**Example — product A/B test:**
- H₀: Conversion rate of variant B = conversion rate of variant A (µ_B − µ_A = 0)
- H₁: Conversion rate of variant B > conversion rate of variant A (µ_B − µ_A > 0)

---

## 2. Components of a well-formed hypothesis

A testable hypothesis identifies three elements [Statsig, 2024]:

1. **Independent variable** — what is changed or compared (the treatment, grouping factor, or predictor).
2. **Dependent variable** — what is measured as the outcome.
3. **Predicted relationship** — the direction and, where possible, the magnitude of the expected effect.

A useful structural template is the **if/then** form:

> "If [change to independent variable], then [expected change in dependent variable]."

For statistical expression, the equivalent form is:

> H₀: [population parameter] [= / ≤ / ≥] [reference value]
> H₁: [population parameter] [≠ / > / <] [reference value]

---

## 3. Directional vs. non-directional hypotheses

The shape of the alternative hypothesis determines which tail(s) of the sampling distribution are used to evaluate evidence:

| Type | Alternative | Test | When to use |
|---|---|---|---|
| **Non-directional (two-tailed)** | H₁: µ₁ ≠ µ₂ | Two-tailed | No strong prior on the direction of the effect; the conservative default for most scientific work |
| **Directional (one-tailed, upper)** | H₁: µ₁ > µ₂ | One-tailed (right) | Strong prior or theoretical reason the effect can only be positive |
| **Directional (one-tailed, lower)** | H₁: µ₁ < µ₂ | One-tailed (left) | Strong prior or theoretical reason the effect can only be negative |

A one-tailed test concentrates the entire significance level (α) in one tail, making it easier to reach significance in the predicted direction — but it has **zero statistical power to detect effects in the opposite direction** [Statistics by Jim, 2024]. This asymmetry means the choice must be made before data collection, justified by prior theory, and documented. Choosing one-tailed post-hoc to obtain a significant p-value is a form of p-hacking [PMC/4317225, 2015].

---

## 4. Grounding hypotheses in prior knowledge

Strong hypothesis formulation draws on existing theory, domain expertise, and prior data rather than pure speculation. Sources include:

- **Literature / prior studies** — what effect sizes and directions have been replicated?
- **Domain expertise** — practitioners with deep subject-matter knowledge can narrow the plausible hypothesis space, especially when data collection is expensive [ScienceDirect, 2024; Copernicus BG, 2022].
- **EDA outputs** — exploratory analysis can surface candidate patterns, but hypotheses derived from EDA on the same dataset must be tested on *independent* data or treated as exploratory, not confirmatory [Scribbr, 2024; PMC/4317225, 2015].

The discipline here is to form hypotheses **before** looking at the data you intend to analyze — or to explicitly label them as exploratory if formed after.

---

## 5. The EDA → hypothesis formulation → CDA workflow

In a full data analysis lifecycle, hypothesis formulation sits between exploratory and confirmatory phases:

| Phase | Purpose | Hypothesis role |
|---|---|---|
| **EDA** (Exploratory Data Analysis) | Inductive: find patterns, anomalies, candidate signals | Generates candidate hypotheses |
| **Hypothesis formulation** | Translate candidates into precise, testable statistical claims | The output of this step |
| **CDA** (Confirmatory Data Analysis) | Deductive: evaluate pre-specified hypotheses against data | Consumes the formulated hypotheses |

EDA and CDA apply different epistemic standards [ResearchGate, 2015]:
- EDA is flexible, open, and discovery-oriented; findings are provisional.
- CDA is structured, hypothesis-driven, and probabilistic; results can be stated with a quantified error rate.

Using the same dataset for both EDA and CDA without correction inflates the false-positive rate. Remedies include data splitting (hold-out set for confirmatory analysis), pre-registration, or explicitly reporting results as exploratory [PMC/4317225, 2015].

---

## 6. SMART criteria for hypothesis quality

A practical self-check before locking in a hypothesis [Nudge Now / LinkedIn, 2024]:

| Criterion | Check |
|---|---|
| **Specific** | Does it identify exactly which variables and which population? |
| **Measurable** | Can the dependent variable be quantified with available data? |
| **Achievable** | Is the sample size realistically obtainable to power the test? |
| **Relevant** | Does it address the underlying business or research question? |
| **Time-bound** | Is the time window for data collection defined? |

A hypothesis that fails any criterion should be revised before analysis proceeds.

---

## 7. Statistical forms by test type

Different tests express the null and alternative in different notation [Scribbr, 2024]:

| Analysis type | H₀ | H₁ |
|---|---|---|
| Two-sample t-test | µ₁ = µ₂ | µ₁ ≠ µ₂ (or > or <) |
| One-sample t-test | µ = µ₀ | µ ≠ µ₀ |
| Correlation | ρ = 0 | ρ ≠ 0 |
| Linear regression (slope) | β₁ = 0 | β₁ ≠ 0 |
| Chi-square (independence) | Variables are independent | Variables are not independent |
| ANOVA | µ₁ = µ₂ = … = µₖ | At least one µᵢ differs |

Writing the hypothesis in the correct mathematical form before choosing a test prevents mismatches between the stated claim and the procedure used to evaluate it.

---

## 8. Pitfalls and how to avoid them

### 8.1 HARKing (Hypothesizing After Results are Known)
Analyzing data, finding patterns, then reporting those patterns as if the hypothesis pre-existed the data. This creates **circular reasoning**: the same data are used both to generate and to evaluate the hypothesis, so statistical p-values lose their calibration and observed relationships are more likely noise than signal [Embassy.science/PMC, 2015].

**Remedy:** Pre-register hypotheses (time-stamp them before data collection), or clearly label any hypothesis developed from the data as exploratory.

### 8.2 Vague or compound hypotheses
"Sales will improve" is not testable. "Customers who receive the email campaign will have higher average order value than those who do not, within the 30-day window" is testable. Compound hypotheses ("A will increase and B will decrease") need separate H₀/H₁ pairs to avoid ambiguity in interpretation [Statsig, 2024; Scribbr, 2024].

### 8.3 Choosing one-tailed tests post-hoc
Switching from two-tailed to one-tailed after seeing that results are significant in one direction roughly halves the p-value. This is undisclosed researcher flexibility (p-hacking) [PMC/4317225, 2015; Statistics by Jim, 2024].

### 8.4 Confusing statistical significance with practical significance
Even a correctly formulated and tested hypothesis can yield a significant p-value that reflects a trivially small effect. Effect size (Cohen's d, r², odds ratio) should be specified as part of the hypothesis when practically important thresholds are known.

### 8.5 Ignoring the null's implications
Failing to reject H₀ does not mean H₀ is true — it means the data did not provide sufficient evidence against it. Reporting "there is no effect" when a test was underpowered misrepresents the finding [PMC/4317225, 2015].

### 8.6 Overly specific hypotheses without domain justification
Hypothesizing a precise effect size (e.g., "the mean will increase by exactly 5 units") without a theoretical basis creates a brittle test. Unless theory mandates a specific value, a directional or inequality form is more appropriate.

---

## 9. Worked example end-to-end

**Business question:** "Does adding a live-chat widget to the checkout page increase purchase completion rates?"

**Step 1 — Identify variables:**
- Independent variable: presence/absence of live-chat widget (treatment vs. control)
- Dependent variable: checkout completion rate (proportion of sessions ending in purchase)

**Step 2 — State hypotheses:**
- H₀: p_treatment = p_control (no difference in completion rates)
- H₁: p_treatment > p_control (directional; prior A/B tests of friction-reduction interventions have shown positive effects)

**Step 3 — Check SMART:**
- Specific: yes — single feature, single metric, two-group comparison
- Measurable: yes — completion rate is directly observable in the analytics platform
- Achievable: power analysis → need N=2,000 per group at α=0.05, 80% power, for a 2% lift; achievable in two weeks of traffic
- Relevant: yes — aligns with conversion-rate optimization goal
- Time-bound: yes — two-week experiment window

**Step 4 — Document before touching data** (pre-register or record in analysis plan).

---

## References

1. Scribbr. "Null & Alternative Hypotheses | Definitions, Templates & Examples." https://www.scribbr.com/statistics/null-and-alternative-hypotheses/ (accessed 2024).
2. Statsig. "The essential parts of a hypothesis: How to formulate one." https://www.statsig.com/perspectives/essential-parts-hypothesis-formulation (accessed 2024).
3. Open University / OpenLearn. "Data analysis: hypothesis testing — 1.3 Hypothesis formulation." https://www.open.edu/openlearn/science-maths-technology/data-analysis-hypothesis-testing/content-section-3.3 (accessed 2024).
4. Statistics by Jim. "One-Tailed and Two-Tailed Hypothesis Tests Explained." https://statisticsbyjim.com/hypothesis-testing/one-tailed-two-tailed-hypothesis-tests/ (accessed 2024).
5. Motulsky H.J. & Lorber A. "Common misconceptions about data analysis and statistics." *PMC* 4317225 (2015). https://pmc.ncbi.nlm.nih.gov/articles/PMC4317225/
6. Embassy.science / PMC. "Hypothesizing after the results are known (HARKing)." https://embassy.science/wiki/Theme:Cc742a7b-826d-4201-b33e-457f2ef79fb9
7. ResearchGate. "The logic of exploratory and confirmatory data analysis." https://www.researchgate.net/publication/283437564_The_logic_of_exploratory_and_confirmatory_data_analysis (2015).
8. ScienceDirect. "A hypotheses-driven framework for human–machine expertise process." https://www.sciencedirect.com/science/article/pii/S1389041724000494 (2024).
