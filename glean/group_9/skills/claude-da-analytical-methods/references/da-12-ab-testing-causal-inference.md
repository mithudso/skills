<!-- hub-reference-banner -->
> **Reference file — part of the `da-analytical-methods` hub.** Formerly the standalone `da-12-ab-testing-causal-inference` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-12-ab-testing-causal-inference
title: A/B Testing and Causal Inference
version: "1.0.0"
updated: "2026-05-30"
category: data-analysis
origin: local
description: >
  A/B testing and causal inference — randomized experiments (sample size, power,
  multiple testing, sequential testing, CUPED, network effects, switchback,
  multi-armed bandits) and observational causal methods (Pearl's causal hierarchy,
  DAGs, backdoor / frontdoor criteria, propensity scores, RDD, instrumental
  variables, difference-in-differences, synthetic control, heterogeneous treatment
  effects via causal forests / uplift modeling).
  TRIGGER: user asks about A/B test design / analysis, sample-size calculation,
  power, multiple testing correction, sequential or always-valid tests,
  multi-armed bandits, CUPED, network effects in experiments, propensity score
  methods, regression discontinuity, instrumental variables, difference-in-
  differences, synthetic control, causal inference from observational data, or
  needs to disambiguate correlation from causation in a specific analysis.
  SKIP: pure descriptive statistics (use da-1-3-probability-theory or
  da-5-exploratory-data-analysis); ML prediction without causal claim (da-7);
  bias and fairness without experimental design (da-11-ethics-and-privacy).
triggers:
  - A/B test
  - AB testing
  - randomized experiment
  - sample size calculation
  - power analysis
  - sequential testing
  - multi-armed bandit
  - CUPED
  - network effects experiment
  - causal inference
  - propensity score
  - regression discontinuity
  - RDD
  - instrumental variable
  - difference-in-differences
  - synthetic control
  - causal forest
  - uplift modeling
  - DAG causal
  - Pearl backdoor frontdoor
keywords:
  - A/B-testing
  - MDE
  - sequential-test
  - mSPRT
  - O'Brien-Fleming
  - bandit
  - Thompson-sampling
  - UCB
  - CUPED
  - SUTVA
  - switchback
  - cluster-randomization
  - Pearl
  - DAG
  - d-separation
  - backdoor-criterion
  - frontdoor-criterion
  - propensity-score
  - IPTW
  - doubly-robust
  - RDD
  - bandwidth
  - instrumental-variables
  - LATE
  - exclusion-restriction
  - difference-in-differences
  - parallel-trends
  - synthetic-control
  - causal-forest
  - uplift
when_to_use:
  - Designing or analyzing a randomized experiment
  - Calculating sample size / minimum detectable effect / power
  - Choosing a multiple-testing correction or sequential design
  - Switching from A/B to a multi-armed bandit
  - Diagnosing or correcting for interference / network effects
  - Estimating causal effects from observational data
  - Building a DAG to justify an identification strategy
  - Picking between propensity matching, IPTW, doubly robust estimation
  - Setting up RDD, IV, DiD, or synthetic control
  - Estimating heterogeneous treatment effects (causal forests, uplift)
when_not_to_use:
  - Pure descriptive statistics — use da-1-3-probability-theory or da-5
  - ML prediction without a causal claim — use da-7-machine-learning
  - General hypothesis testing without an experimental design — use da-1-4-statistical-inference-foundations
related_skills:
  - da-1-3-probability-theory
  - da-1-4-statistical-inference-foundations
  - da-6-statistical-modeling
  - da-1-6-1-correlation-vs-causation
  - da-11-ethics-and-privacy
---

# A/B Testing and Causal Inference

The two pillars of "did X cause Y": randomized experiments when you can run one, and observational causal methods when you can't. This skill covers both.

## When to use this skill

Activate when the user:
- is designing or analyzing a randomized experiment
- needs sample size / power / MDE calculation
- is dealing with multiple testing, sequential tests, or bandits
- is diagnosing interference or network effects in an experiment
- needs to estimate a causal effect from observational data
- is building a DAG to justify an identification strategy
- is picking between propensity / IPTW / RDD / IV / DiD / synthetic control

## When NOT to use this skill

- Descriptive statistics → `da-1-3-probability-theory` / `da-5-exploratory-data-analysis`
- Prediction without causal claim → `da-7-machine-learning`
- General hypothesis testing without an experiment → `da-1-4-statistical-inference-foundations`

---

## Part 1 — Randomized experiments (A/B testing)

### Fundamentals

Randomly assign units (users, sessions, requests) to control or treatment. Under randomization, the average difference in outcomes between arms is an unbiased estimate of the average treatment effect.

The minimum viable experiment design records:
- **Unit of randomization** (user, session, request — has to be the unit at which interference is bounded)
- **Treatment definition** (must be deterministic given the unit)
- **Primary metric** (one, ideally; if more, declare them in advance)
- **Sample size** (computed from power analysis, not vibes)
- **Decision rule** (what conclusion follows from what test statistic)

### Sample size and power analysis

For a two-sample test of means at significance `α` and power `1 - β`:

```
n ≈ ((z_{α/2} + z_β)^2 · 2σ²) / Δ²
```

Where Δ is the **minimum detectable effect** (MDE) you care about. The headline tradeoff: halving the MDE quadruples the sample size.

Tools:
- **`statsmodels.stats.power`** (Python) — `tt_ind_solve_power`
- **G*Power** (desktop app) — clean UI, lots of test types
- **`pwr` package** (R)
- Online calculators — fine for the headline number; don't trust them for stratified designs

Defaults: `α = 0.05`, power = `0.80`. Walk in expecting roughly thousands of users per arm for percentage-point effects on conversion rates, tens of thousands for tenths-of-a-percent.

### Multiple testing

If you test 20 metrics at `α = 0.05`, you expect 1 false positive even if nothing's happening.

| Correction | Controls | Best for |
|---|---|---|
| **Bonferroni** | Family-wise error rate (FWER) | Few tests, conservative |
| **Holm-Bonferroni** | FWER | Slightly less conservative than Bonferroni |
| **Benjamini-Hochberg** | False discovery rate (FDR) | Many tests; you accept some false positives |
| **Knockoff filter** | FDR with mock variables | Modern alternative; works on selection problems |

Companies running 100+ tests/week typically use BH-FDR with `q = 0.10`. Bonferroni on a 50-metric dashboard is too conservative to detect anything.

### Sequential testing

Standard A/B test analysis assumes you peek **once**, at the end. Peeking early and stopping inflates the false-positive rate. Two fixes:

- **Group sequential designs (O'Brien-Fleming, Pocock)** — pre-specify N looks at fixed times; allocate `α` across looks via a spending function
- **Always-valid p-values (mSPRT)** — Optimizely / Statsig / GrowthBook style. Look as often as you want; the p-value is honest at every moment

mSPRT trades raw efficiency for the freedom to look. If you're going to peek anyway (you are), use always-valid.

### Multi-armed bandits

Instead of equal-split A/B, route traffic adaptively to the better-performing arm. The exploration / exploitation tradeoff.

- **ε-greedy** — explore with probability ε, exploit otherwise
- **UCB (Upper Confidence Bound)** — pick the arm with highest mean + uncertainty bonus
- **Thompson sampling** — Bayesian; sample from each arm's posterior, pick the highest

When to use a bandit:
- Cost of a bad arm during exploration is high (revenue, user experience)
- You're optimizing, not making a binary ship/no-ship decision

When NOT to use a bandit:
- You need a clean causal estimate (bandits produce biased estimates of arm value)
- The optimal arm changes over time and the algorithm hasn't been designed for that (use a non-stationary bandit)

### CUPED — variance reduction

CUPED (Controlled-experiment Using Pre-Experiment Data) uses pre-experiment covariates to soak up variance. Sample-size requirements drop by 30-50% on metrics with strong pre-period predictors. Microsoft published the method in 2013; nearly every mature experimentation platform now includes it.

Implementation in one line: regress the outcome on the covariate, use the residual as the new outcome.

### Network effects and interference (SUTVA violations)

Stable Unit Treatment Value Assumption (SUTVA): one unit's treatment doesn't affect another unit's outcome. Violated when:
- **Marketplace effects** (Airbnb, Uber) — treated drivers affect control riders
- **Social effects** (Meta, LinkedIn) — friends of treated users see the treatment indirectly
- **Switchover effects** — same user can see both arms

Mitigations:
- **Cluster randomization** — assign by group (city, social-graph community) so spillovers stay inside the cluster
- **Switchback testing** — flip the entire system between treatment and control on a schedule; analyze with a time-series model
- **Ego-network randomization** — assign by user + their immediate neighbors

---

## Part 2 — Observational causal inference

When you can't randomize.

### Pearl's causal hierarchy

| Rung | Question | Example | Method |
|---|---|---|---|
| 1. Association | What is? | `P(Y | X)` | Regression, ML |
| 2. Intervention | What if I do? | `P(Y | do(X))` | RCT, backdoor/frontdoor |
| 3. Counterfactual | What if I had done? | `P(Y_x | X', Y')` | Counterfactual models |

Most ML lives on rung 1. Most business decisions need rung 2. Confusing them is the source of "correlation isn't causation" failures.

### DAGs and identification

A directed acyclic graph encodes assumed causal relationships. Once you have a DAG you can mechanically check whether the causal effect of X on Y is **identifiable** from observed data.

- **d-separation** — set of variables that "block" a path between X and Y
- **Backdoor criterion** — a set Z that (a) blocks every non-causal path from X to Y, and (b) contains no descendant of X. Conditioning on Z gives the causal effect.
- **Frontdoor criterion** — when no admissible backdoor adjustment exists, but there's a mediator M such that X → M → Y and no other path

Tooling: DAGitty (browser), `dagitty` (R), `causalgraphicalmodels` (Python), `dowhy` (Python).

### Propensity score methods

The propensity score `e(X) = P(T=1 | X)` is the probability of receiving treatment given covariates. Three uses:

- **Matching** — pair each treated unit with a control unit of similar propensity
- **IPTW** (Inverse Probability of Treatment Weighting) — weight by `1/e(X)` for treated, `1/(1-e(X))` for control
- **Doubly robust estimation** — combine an outcome model and a propensity model; consistent if either is correct

Pitfalls:
- Extreme propensities (`e ≈ 0` or `1`) → unstable weights → wild estimates. Trim or use stabilized weights.
- Propensity scores don't fix unmeasured confounders. They handle confounders you've measured.

### Regression discontinuity (RDD)

Some treatment threshold exists (test score for scholarship, age for retirement benefit). Units just above and just below the threshold are otherwise similar. Compare them.

- **Sharp RDD** — treatment determined exactly by threshold
- **Fuzzy RDD** — threshold changes probability of treatment; use as an instrument

Bandwidth choice is the critical decision; modern practice uses optimal-bandwidth methods (Calonico-Cattaneo-Titiunik).

### Instrumental variables (IV)

An instrument `Z` affects `Y` only through `X`. The classic example: distance to college as an instrument for years of education, when estimating the returns to education.

Three requirements (only the first two are testable):
1. **Relevance** — `Z` predicts `X`
2. **Exogeneity** — `Z` is independent of unobserved confounders
3. **Exclusion restriction** — `Z` affects `Y` *only* through `X` (untestable)

Weak instruments (low Z-X correlation) produce wildly biased estimates. Stock-Yogo critical values for the first-stage F-statistic are the standard check; F > 10 is the rule of thumb.

LATE (Local Average Treatment Effect) is what IV identifies — the effect on the **compliers**, the units whose treatment status depends on `Z`. Not the average treatment effect.

### Difference-in-differences (DiD)

A treatment is rolled out to some units at some time. Compare:
- Treated units before vs after
- Untreated units before vs after
- DiD estimate = (treated post - treated pre) - (control post - control pre)

The critical assumption: **parallel trends**. The treated and control groups' outcomes would have moved in parallel absent the treatment. Test by plotting pre-period trends. If they diverge before the treatment, DiD is biased.

Recent advances (Callaway & Sant'Anna 2021, Goodman-Bacon 2021) handle staggered rollouts and heterogeneous treatment effects properly; the classic two-way fixed-effects estimator can be badly biased in these cases.

### Synthetic control

When you have one treated unit and many control units, construct a weighted combination of controls that matches the treated unit's pre-period trajectory. The synthetic counterfactual predicts what would have happened. Originated in Abadie, Diamond, Hainmueller (Basque Country, 2003; California Prop 99, 2010).

Modern variants: synthetic difference-in-differences (Arkhangelsky et al, 2021).

### Heterogeneous treatment effects

Average treatment effect (ATE) hides variation. Some sub-populations benefit hugely, others not at all, some are harmed.

- **Causal forests** (Wager & Athey 2018) — non-parametric estimation of conditional average treatment effects (CATE)
- **Uplift modeling** — predict the incremental effect of treatment per individual; route the treatment to those with positive uplift only
- **Meta-learners** — S-learner, T-learner, X-learner, R-learner; use any base ML model

Tooling: `econml` (Microsoft, Python), `grf` (Stanford, R), `causalml` (Uber).

---

## Online experimentation platforms (2026)

| Platform | Strength |
|---|---|
| **Optimizely** | Mature, marketing-oriented |
| **Statsig** | Sequential tests, generous free tier, modern UX |
| **GrowthBook** | Open-source option, feature-flag native |
| **Eppo** | Modern, strong CUPED + heterogeneous effects |
| **Internal (Meta / Google / Microsoft / Netflix)** | The biggest ones run their own; published papers describe the architectures |

The 2024-2026 frontier has been **CUPED-by-default**, **sequential tests by default**, **heterogeneous-effect estimation in-platform**, and integration with feature flags.

---

## Anti-patterns

1. **Stopping early without a sequential test.** Inflates false-positive rate dramatically.
2. **Running 20 metrics, picking the one that won, calling it a win.** Multiple-testing failure.
3. **Comparing pre-period to post-period without a control group.** Confounded with seasonality, secular trends.
4. **Treating an ML-derived treatment effect as causal.** Predictive accuracy ≠ causal validity.
5. **Ignoring SUTVA in marketplace / social products.** The control arm sees the treatment indirectly.
6. **Propensity matching with unobserved confounders.** Doesn't fix what you didn't measure.
7. **DiD without checking parallel trends.** Plot the pre-period.
8. **IV with a weak instrument.** F-statistic check.
9. **RDD with a manipulated running variable.** McCrary density test.
10. **Treating the LATE as the ATE.** IV identifies a specific subgroup's effect.

---

## References

1. Kohavi, R., Tang, D., & Xu, Y. (2020). *Trustworthy Online Controlled Experiments*. Cambridge University Press. The canonical reference.
2. Pearl, J. (2009). *Causality: Models, Reasoning, and Inference* (2nd ed.). Cambridge.
3. Pearl, J. & Mackenzie, D. (2018). *The Book of Why*. The accessible version.
4. Imbens, G. W. & Rubin, D. B. (2015). *Causal Inference for Statistics, Social, and Biomedical Sciences*.
5. Cunningham, S. (2021). *Causal Inference: The Mixtape*. mixtape.scunning.com (free online).
6. Hernán, M. A. & Robins, J. M. (2020). *Causal Inference: What If*. harvard.edu/...
7. Athey, S. & Imbens, G. (2019). "Machine Learning Methods for Estimating Heterogeneous Causal Effects." *Annual Review of Economics*.
8. Deng, A. et al. (2013). "Improving the Sensitivity of Online Controlled Experiments by Utilizing Pre-Experiment Data." WSDM. (CUPED paper.)
9. Howard, S. R. et al. (2021). "Time-uniform, nonparametric, nonasymptotic confidence sequences." (Always-valid inference.)
10. Calonico, S., Cattaneo, M. D., & Titiunik, R. (2014). "Robust Nonparametric Confidence Intervals for Regression-Discontinuity Designs." *Econometrica*.
11. Callaway, B. & Sant'Anna, P. (2021). "Difference-in-Differences with multiple time periods." *J. Econometrics*.
12. Wager, S. & Athey, S. (2018). "Estimation and Inference of Heterogeneous Treatment Effects using Random Forests." *JASA*.
13. Abadie, A., Diamond, A., & Hainmueller, J. (2010). "Synthetic Control Methods for Comparative Case Studies." *JASA*.
14. econml (Microsoft) — https://econml.azurewebsites.net/
15. DoWhy (PyWhy) — https://www.pywhy.org/dowhy/
