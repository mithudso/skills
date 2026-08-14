<!-- hub-reference-banner -->
> **Reference file — part of the `da-analytical-methods` hub.** Formerly the standalone `da-24-survival-analysis` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-24-survival-analysis
title: Survival Analysis / Time-to-Event Analysis
version: "1.1.0"
updated: "2026-05-30"
category: data-analysis
origin: local
description: >
  Survival analysis / time-to-event modeling — its own discipline, distinct
  from the cross-sectional regression in da-6 and the time-series forecasting
  in da-15. Covers censoring (right/left/interval) and truncation; the
  survival, hazard, and cumulative-hazard functions; Kaplan-Meier and
  Nelson-Aalen estimators; the log-rank test; the Cox proportional-hazards
  model and its assumption diagnostics (Schoenfeld residuals, cox.zph);
  parametric models (Weibull, exponential) and accelerated failure time (AFT);
  competing risks (cause-specific hazard vs. Fine-Gray subdistribution / CIF);
  time-varying covariates; discrete-time survival; churn / retention / customer
  lifetime applications; and ML survival methods (random survival forests,
  gradient-boosted survival, DeepSurv). Tooling: lifelines, scikit-survival,
  R survival/survminer.
  TRIGGER: modeling time-until-an-event (death, failure, churn, conversion,
  default, readmission); data has censored or truncated observations; need a
  survival curve, hazard ratio, or median survival time; fitting Kaplan-Meier,
  Cox PH, Weibull/AFT, or a survival forest; checking the proportional-hazards
  assumption; competing risks / cumulative incidence; time-to-churn or
  retention-curve modeling (the survival half of a CLV analysis); "how long
  until", "time to event", "hazard ratio", "survival probability".
  SKIP: forecasting a numeric time series (da-15-forecasting); ordinary
  regression/classification with no time-to-event or censoring
  (da-6-statistical-modeling, da-7-machine-learning); A/B test or causal
  inference with no time component (da-12); anomaly/change-point detection
  (da-16); simple cohort retention table with no modeling (da-21-product-analytics);
  computing a CLV dollar value via BG/NBD, Pareto/NBD, or Gamma-Gamma spend
  models (da-23-customer-lifetime-value — use this skill only for the
  time-to-churn / survival-curve half of CLV).
triggers:
  - survival analysis
  - time-to-event
  - time to event analysis
  - Kaplan-Meier
  - hazard ratio
  - Cox proportional hazards
  - censoring
  - censored data
  - log-rank test
  - Weibull survival
  - accelerated failure time
  - competing risks
  - cumulative incidence
  - random survival forest
  - DeepSurv
  - churn survival model
  - customer lifetime survival
  - how long until event
keywords:
  - survival-function
  - hazard-function
  - cumulative-hazard
  - right-censoring
  - left-truncation
  - interval-censoring
  - Kaplan-Meier
  - Nelson-Aalen
  - log-rank
  - Cox-PH
  - partial-likelihood
  - proportional-hazards
  - Schoenfeld-residuals
  - cox.zph
  - Weibull
  - exponential
  - AFT
  - competing-risks
  - cause-specific-hazard
  - Fine-Gray
  - subdistribution-hazard
  - cumulative-incidence
  - time-varying-covariates
  - discrete-time-survival
  - random-survival-forest
  - gradient-boosted-survival
  - DeepSurv
  - concordance-index
  - lifelines
  - scikit-survival
when_to_use:
  - Modeling the time until an event (death, failure, churn, default, conversion)
  - Data contains censored or truncated observations
  - Need a survival curve, hazard ratio, or median survival estimate
  - Comparing survival between groups (treatment vs. control, cohorts)
  - Checking or repairing the proportional-hazards assumption
  - Competing risks / cumulative-incidence estimation
  - Customer lifetime value, retention curves, or churn timing
  - ML survival prediction with nonlinear covariate effects
when_not_to_use:
  - Forecasting a numeric time series — use da-15-forecasting
  - Plain regression/classification, no censoring — use da-6 or da-7
  - A/B test or causal inference, no time-to-event — use da-12
  - Anomaly or change-point detection — use da-16-anomaly-detection
  - A simple cohort retention table, no modeling — use da-21-product-analytics
  - Computing a CLV dollar value (BG/NBD, Pareto/NBD, Gamma-Gamma) — use da-23-customer-lifetime-value
related_skills:
  - da-1-3-probability-theory
  - da-1-3-3-4-exponential
  - da-1-4-statistical-inference-foundations
  - da-6-statistical-modeling
  - da-7-machine-learning
  - da-15-forecasting
  - da-21-product-analytics
  - da-23-customer-lifetime-value
  - da-25-bayesian-data-analysis
---

# Survival Analysis / Time-to-Event Analysis

Modeling the *time until an event happens* when some observations are **incomplete** (censored or truncated). This is its own discipline because ordinary regression cannot use a row that says "this customer had not churned yet when we stopped looking" — survival methods extract information from exactly those incomplete rows. Canonical textbooks: Klein & Moeschberger *Survival Analysis: Techniques for Censored and Truncated Data* (2nd ed, 2003); Therneau & Grambsch *Modeling Survival Data* (2000). Primary Python tooling: **lifelines** and **scikit-survival**; R: **survival** + **survminer**.

## When to use this skill

- The outcome is a *duration* until an event: death, machine failure, churn, loan default, conversion, hospital readmission.
- Some subjects have **not** experienced the event by end of observation (censoring), or only entered observation partway through (truncation).
- You need a survival curve, hazard ratio, median time-to-event, or cumulative incidence.

## When NOT to use this skill

- Forecasting a numeric series over calendar time → `da-15-forecasting`
- Regression/classification with fully observed outcomes → `da-6` / `da-7`
- Causal/experiment analysis with no time component → `da-12`
- A descriptive cohort retention table (no estimator, no model) → `da-21-product-analytics`
- Computing a CLV dollar figure with BG/NBD, Pareto/NBD, or Gamma-Gamma spend models → `da-23-customer-lifetime-value` (this skill covers only the time-to-churn / survival-curve half)

---

## 1. Censoring and truncation — the defining feature

The reason survival analysis exists. Get this wrong and every downstream estimate is biased.

| Mechanism | What it means | Handling |
|---|---|---|
| **Right censoring** | Event not yet observed at end of follow-up (most common case). You know `T > c`. | Standard; all estimators below assume it. |
| **Left censoring** | Event already happened before observation began, exact time unknown. `T < c`. | Use models that accept left-censored entries (lifelines `KaplanMeierFitter.fit_left_censoring`). |
| **Interval censoring** | Event happened between two inspection times. `a < T < b`. | Turnbull estimator / interval-censored regression. |
| **Left truncation** | Subjects who had the event *before entry* never appear at all (delayed entry, e.g. age-as-timescale). | Supply an `entry`/`lower_bound`; biases ignored if untreated. |
| **Right truncation** | Only subjects who *have* had the event are sampled (e.g. registry of completed events). | Specialized estimators; rare. |

Key distinction: **censoring keeps the subject but loses event-time detail; truncation removes the subject from the sample entirely** ([Stats Ox lecture notes](https://www.stats.ox.ac.uk/~mlunn/lecturenotes1.pdf), 2020; [NJIT Math 659 Ch.3](https://web.njit.edu/~wguo/Math%20659_2011/Math659_Chapter3.pdf), 2011; [GeeksforGeeks: Censoring and Truncation](https://www.geeksforgeeks.org/data-science/censoring-and-truncation/), 2024). The standard estimators assume censoring is **non-informative** (independent of the event process) — a customer who churns *because* they were about to be observed violates this.

## 2. The survival, hazard, and cumulative-hazard functions

These three are interchangeable views of the same distribution; pick whichever the audience reads best.

- **Survival function** `S(t) = P(T > t)` — probability of surviving past `t`. Monotone non-increasing from 1.
- **Hazard function** `h(t) = lim Δ→0 P(t ≤ T < t+Δ | T ≥ t)/Δ` — instantaneous event rate *given survival so far*. The risk "right now."
- **Cumulative hazard** `H(t) = ∫₀ᵗ h(u)du`, with the bridge identity `S(t) = exp(−H(t))`.

The hazard is the modeling target for most methods (Cox, parametric) because covariate effects are cleanest there ([lifelines Quickstart](https://lifelines.readthedocs.io/en/latest/Quickstart.html), v0.30, 2025; Klein & Moeschberger Ch. 2, 2003).

## 3. Kaplan-Meier & Nelson-Aalen (non-parametric estimators)

The first thing to compute on any survival dataset — assumption-free descriptive curves.

- **Kaplan-Meier (product-limit) estimator** of `S(t)`: at each event time, multiply by `(1 − dᵢ/nᵢ)` where `dᵢ` events occur among `nᵢ` at risk. Step function; censored subjects drop out of the risk set without causing a step. Report the **median survival** (where `S(t)=0.5`) and confidence bands.
- **Nelson-Aalen estimator** of the cumulative hazard `H(t)`: sum of `dᵢ/nᵢ`. Estimates cumulative hazard non-parametrically under independent right-censoring and left-truncation ([lifelines NelsonAalenFitter / Survival regression docs](https://lifelines.readthedocs.io/en/latest/Survival%20analysis%20with%20lifelines.html), 2025).

```python
from lifelines import KaplanMeierFitter
kmf = KaplanMeierFitter()
kmf.fit(durations=df["tenure"], event_observed=df["churned"], entry=df.get("entry"))
kmf.median_survival_time_; kmf.plot_survival_function()
```

Sources: [lifelines Quickstart](https://lifelines.readthedocs.io/en/latest/Quickstart.html) (2025); Klein & Moeschberger Ch. 4 (2003); [CPSC 330 Survival lecture](https://ubc-cs.github.io/cpsc330-2023W1/lectures/20_survival-analysis.html) (2023).

## 4. The log-rank test (comparing groups)

Compares two-or-more KM curves: null hypothesis is **equal survival across groups**. It is a chi-square test that accumulates observed-minus-expected events at each event time across groups; it weights all time points equally (Wilcoxon/Tarone-Ware variants weight early times more). Use it to compare treatment arms or cohorts, but note it gives only a p-value, not an effect size — for an effect size use Cox.

```python
from lifelines.statistics import logrank_test, multivariate_logrank_test
logrank_test(durA, durB, eventA, eventB).p_value
```

Sources: [lifelines.statistics](https://lifelines.readthedocs.io/en/latest/lifelines.statistics.html) (2025); [STHDA log-rank](https://www.sthda.com/english/wiki/cox-model-assumptions) (2018); Klein & Moeschberger Ch. 7 (2003).

## 5. Cox proportional-hazards model (the workhorse)

Semi-parametric regression: `h(t|x) = h₀(t) · exp(βᵀx)`. The baseline hazard `h₀(t)` is left **unspecified** (estimated non-parametrically); only the `β` are estimated, via **partial likelihood** (Cox 1972). `exp(βⱼ)` is the **hazard ratio** for covariate `j` — a multiplicative, time-constant effect.

```python
from lifelines import CoxPHFitter
cph = CoxPHFitter(penalizer=0.1)
cph.fit(df, duration_col="tenure", event_col="churned")
cph.print_summary()          # coef, exp(coef)=HR, p, CI
cph.predict_partial_hazard(df)
```

Tie handling: Efron (default, preferred) or Breslow. Report HRs with CIs, not raw coefficients. Sources: [lifelines CoxPHFitter](https://lifelines.readthedocs.io/en/latest/Quickstart.html) (2025); Therneau & Grambsch (2000); [Researchers' Guide: Cox PH in Python](https://medium.com/the-researchers-guide/survival-analysis-in-python-km-estimate-cox-ph-and-aft-model-5533843c5d5d) (2021).

## 6. The proportional-hazards assumption & diagnostics

Cox is only valid if hazard ratios are **constant over time**. Always check.

- **Scaled Schoenfeld residuals**: one residual per covariate per event time. Under PH, they have **zero slope against (a function of) time**.
- **Grambsch-Therneau test** (`cox.zph` in R, `cph.check_assumptions()` / `proportional_hazard_test` in lifelines): null = slope is zero = PH holds. A small p-value flags a violation. A non-significant p-value (>0.05) supports PH.
- **Graphical**: `ggcoxzph()` (survminer) — the LOESS smooth of scaled residuals should be flat/parallel to the x-axis.

**Fixes when violated**: stratify on the offending covariate (`strata=`), add a time-interaction term (covariate × `f(t)`), split follow-up into intervals with separate effects, or switch to an AFT model.

Sources: [Grambsch-Therneau / cox.zph, UCLA OARC](https://stats.oarc.ucla.edu/other/examples/asa2/testing-the-proportional-hazard-assumption-in-cox-models/) (2021); [Stata stcox PH-assumption tests](https://www.stata.com/manuals14/ststcoxph-assumptiontests.pdf) (2015); [STHDA Cox model assumptions](https://www.sthda.com/english/wiki/cox-model-assumptions) (2018).

## 7. Parametric models: exponential, Weibull, and AFT

When you want a smooth survival curve, extrapolation beyond observed follow-up, or a fully generative model.

- **Exponential**: constant hazard `h(t)=λ`. Memoryless; rarely realistic but a baseline.
- **Weibull**: `h(t)=λρ(λt)^{ρ−1}` — monotone increasing (`ρ>1`, wear-out) or decreasing (`ρ<1`, early-failure) hazard. The default parametric choice.
- **AFT (Accelerated Failure Time)**: models `log(T) = βᵀx + error`; covariates **accelerate or decelerate** time-to-event by a constant factor (`exp(β)` = time ratio), instead of multiplying the hazard. More interpretable for "this doubles the expected lifetime."
- **Weibull is the only distribution expressible as both PH and AFT** — its results can be read either way. Log-logistic / log-normal AFT allow non-monotone hazards.

```python
from lifelines import WeibullAFTFitter
aft = WeibullAFTFitter().fit(df, duration_col="tenure", event_col="churned")
```

Sources: [AFT model — Wikipedia](https://en.wikipedia.org/wiki/Accelerated_failure_time_model) (2025); [CRAN eha parametric survival](https://cran.r-project.org/web/packages/eha/vignettes/parametric.html) (2024); [AFT vs Cox comparison, PMC4645729](https://pmc.ncbi.nlm.nih.gov/articles/PMC4645729/) (2015).

## 8. Competing risks (cause-specific vs. Fine-Gray)

When a subject can fail from **mutually exclusive** causes (e.g. churn-to-competitor vs. account-closed-by-fraud), naïve KM/Cox on one cause **over-estimates its incidence** by treating competing events as censored. Two correct approaches:

- **Cause-specific hazard** (Cox per cause): models the rate of cause `k` among those still at risk. Best for **etiology** ("what drives this cause's rate"). Censor competing events.
- **Fine-Gray subdistribution hazard**: directly links covariates to the **cumulative incidence function (CIF)** — the actual probability of cause `k` over time, accounting for competing events. Best for **prediction / risk communication** (`sHR`). Subjects who experience a competing event stay in the risk set with decaying weights.

Caveats: the sum of CIFs across causes can exceed 1 for some covariate patterns if you fit a separate Fine-Gray per cause; **avoid multiple Fine-Gray models** — prefer cause-specific for multi-event scientific questions. For causal effects, Fine-Gray is discouraged.

Sources: [Austin & Fine, Statistics in Medicine](https://onlinelibrary.wiley.com/doi/10.1002/sim.7501) (2017); [Austin et al., Stat Med (CIF>1)](https://onlinelibrary.wiley.com/doi/full/10.1002/sim.9023) (2021); [Statistical Horizons: don't use Fine-Gray for causal analysis](https://statisticalhorizons.com/for-causal-analysis-of-competing-risks/) (2023).

## 9. Time-varying covariates

When a predictor changes during follow-up (subscription tier upgrade, evolving lab value), a single baseline value is wrong. Put the data in **long (counting-process) format**: one row per subject per interval with `(id, start, stop, event, covariates)`.

```python
from lifelines import CoxTimeVaryingFitter
ctv = CoxTimeVaryingFitter()
ctv.fit(long_df, id_col="id", start_col="start", stop_col="stop", event_col="event")
```

`CoxTimeVaryingFitter` implements Cox's time-varying PH model with Efron tie-handling and requires the long format. This is also the standard fix for a *time-varying coefficient* (a PH violation) — though that needs a covariate×time interaction, not just a time-varying value. Sources: [lifelines Time-varying regression](https://lifelines.readthedocs.io/en/latest/Time%20varying%20survival%20regression.html) (2025); [CoxTimeVaryingFitter docs](https://lifelines.readthedocs.io/en/latest/fitters/regression/CoxTimeVaryingFitter.html) (2025); Therneau & Grambsch Ch. 3 (2000).

## 10. Discrete-time survival & churn / retention / CLV

When time is naturally **binned** (months subscribed, billing cycles) and many events tie at the same bin, discrete-time survival is cleaner than continuous Cox.

- **Method**: expand to **person-period** rows (one row per subject per period at risk), then fit ordinary **logistic regression** with the period (or a flexible function of it) as a predictor. The fitted per-period probabilities are the **discrete hazards**; chain them into a survival/retention curve. Add covariates and interactions like any GLM.
- **Churn / retention**: tenure = duration, churn = event, still-active customers = right-censored. KM gives the retention curve; Cox/AFT give "what drives churn timing"; integrating `S(t)` over the margin gives expected lifetime, the backbone of **CLV** (`CLV ≈ Σ margin·S(t)·discount`).

Survival analysis beats a static churn classifier because it answers *when*, uses censored (still-active) customers correctly, and yields retention curves and CLV directly. Sources: [SAS: Survival Data Mining / discrete time](https://support.sas.com/resources/papers/proceedings12/132-2012.pdf) (2012); [SAS: Modeling CLV with survival analysis](https://support.sas.com/resources/papers/proceedings/proceedings/sugi28/120-28.pdf) (2003); [Springer: survival analysis in churn prediction](https://link.springer.com/article/10.1057/s41270-025-00450-2) (2025).

## 11. Machine-learning survival models

When effects are nonlinear, interacting, or high-dimensional and you care about predictive accuracy over interpretability.

- **Random Survival Forests (RSF)**: forest of survival trees; each split maximizes the **log-rank statistic**; aggregate to an ensemble cumulative-hazard estimate. Handles nonlinearities/interactions, right-censoring, and gives variable importance. Introduced by Ishwaran et al. (2008).
- **Gradient-boosted survival (Cox/component-wise loss)**: boosts weak learners against a survival loss; often the strongest tabular baseline. scikit-survival's `GradientBoostingSurvivalAnalysis` reports concordance ~0.75 on the standard example; performance improves with ensemble size then degrades if over-grown.
- **DeepSurv** (Katzman et al., 2018): a deep neural net optimizing the **Cox partial-likelihood** loss — nonlinear Cox for personalized risk/treatment recommendation; performs as well as or better than classical Cox.
- **Evaluation**: **Harrell's concordance index (C-index)** — probability a model ranks a riskier subject as failing earlier (0.5 = random, 1.0 = perfect); also time-dependent AUC and the integrated Brier score.

```python
from sksurv.ensemble import RandomSurvivalForest
from sksurv.metrics import concordance_index_censored
rsf = RandomSurvivalForest(n_estimators=200).fit(X, y_structured)  # y = (event_bool, time)
```

Sources: [Ishwaran et al., Random Survival Forests, Annals of Applied Statistics 2(3):841-860](https://ishwaran.org/papers/IKBL.AOAS.pdf) (2008); [scikit-survival RSF guide](https://scikit-survival.readthedocs.io/en/stable/user_guide/random-survival-forest.html) & [boosting guide](https://scikit-survival.readthedocs.io/en/stable/user_guide/boosting.html) (2025); [Katzman et al., DeepSurv, BMC Med Res Methodol](https://link.springer.com/article/10.1186/s12874-018-0482-1) / [arXiv 1606.00931](https://arxiv.org/abs/1606.00931) (2018).

---

## Methodology (a default workflow)

1. **Define the timeline**: pick `t=0` (origin), the event, and the censoring rule. Decide if left truncation / delayed entry applies.
2. **Describe**: Kaplan-Meier curve + median survival; Nelson-Aalen for cumulative hazard. Stratify by key groups.
3. **Compare groups**: log-rank test (effect size deferred to Cox).
4. **Model effects**: Cox PH first (interpretable HRs). Or parametric/AFT if you need extrapolation or a smooth curve.
5. **Check assumptions**: Schoenfeld residuals / `cox.zph`. Repair PH violations (stratify, time-interaction, AFT).
6. **Handle structure**: competing risks → cause-specific or Fine-Gray; changing covariates → time-varying long format; binned time → discrete-time logistic.
7. **Predict at scale**: RSF / gradient boosting / DeepSurv when accuracy beats interpretability.
8. **Validate**: C-index, time-dependent AUC, integrated Brier; calibration of predicted vs. observed survival; never evaluate with plain accuracy.

## Practical patterns

- **Always plot KM first.** It reveals crossing curves (PH violation), plateaus (cured fraction), and data problems before any model.
- **Encode the outcome as a pair**, never a single column: `(event_indicator, time)`. scikit-survival needs a structured array; lifelines takes two columns.
- **Report hazard ratios with CIs**, and translate: "HR 1.4 → 40% higher instantaneous churn rate," not a raw coefficient.
- **Use the right time origin**: calendar time, age, or time-since-enrollment change the answer. Left-truncate when entry is delayed.
- **For churn/CLV**, integrate the survival curve for expected lifetime instead of guessing an average tenure from completed customers (that ignores still-active = censored ones).

## Anti-patterns

- **Dropping censored rows** ("only keep customers who churned"). This is the cardinal sin — it discards the majority of information and severely biases estimates.
- **Treating time-to-event as a regression target** with OLS — censoring makes the target undefined for survivors.
- **Treating competing events as plain censoring** when estimating one cause's incidence — over-states it. Use CIF / Fine-Gray.
- **Fitting Cox without checking PH** — a violated assumption silently corrupts every hazard ratio.
- **Reporting only a log-rank p-value** with no effect size or curve.
- **Evaluating an ML survival model with accuracy/AUC on a binarized label** instead of the C-index / Brier score.
- **One Fine-Gray model per cause and reading all the CIFs together** — they can sum past 1; prefer cause-specific hazards for multi-event questions.

## Troubleshooting

| Symptom | Likely cause | Fix |
|---|---|---|
| KM curves cross | PH violated | Stratify, time-varying coefficient, or AFT; don't trust a single HR |
| `cox.zph` p-value tiny for a covariate | Non-proportional effect | `strata=` that covariate, or add covariate×time interaction |
| Median survival is `inf`/undefined | Curve never reaches 0.5 (heavy censoring) | Report RMST or a fixed-horizon survival probability instead |
| Cumulative incidence sums > 1 | Multiple Fine-Gray models combined | Use cause-specific hazards, or one Fine-Gray for the single cause of interest |
| C-index ≈ 0.5 | No predictive signal / wrong outcome encoding | Re-check `(event, time)` pairing and feature leakage |
| Cox fails to converge / huge CIs | Separation, collinearity, or too few events | Add `penalizer=`, drop/merge covariates, respect ~10 events-per-variable |
| Suspiciously optimistic effects | Informative censoring / immortal-time bias | Audit how follow-up starts and ends; align time origin with eligibility |

## References

- Klein, J.P. & Moeschberger, M.L. *Survival Analysis: Techniques for Censored and Truncated Data*, 2nd ed. Springer (2003).
- Therneau, T.M. & Grambsch, P.M. *Modeling Survival Data: Extending the Cox Model*. Springer (2000).
- lifelines documentation — Quickstart, statistics, time-varying regression — https://lifelines.readthedocs.io/en/latest/ (v0.30, 2025).
- scikit-survival user guide — Random Survival Forests & Gradient Boosting — https://scikit-survival.readthedocs.io/en/stable/ (2025).
- Censoring & truncation — https://www.stats.ox.ac.uk/~mlunn/lecturenotes1.pdf (2020); https://web.njit.edu/~wguo/Math%20659_2011/Math659_Chapter3.pdf (2011).
- Grambsch-Therneau PH test — https://stats.oarc.ucla.edu/other/examples/asa2/testing-the-proportional-hazard-assumption-in-cox-models/ (2021); Stata https://www.stata.com/manuals14/ststcoxph-assumptiontests.pdf (2015).
- AFT models — https://en.wikipedia.org/wiki/Accelerated_failure_time_model (2025); https://pmc.ncbi.nlm.nih.gov/articles/PMC4645729/ (2015).
- Fine-Gray / competing risks — Austin & Fine, Stat Med https://onlinelibrary.wiley.com/doi/10.1002/sim.7501 (2017); https://onlinelibrary.wiley.com/doi/full/10.1002/sim.9023 (2021); https://statisticalhorizons.com/for-causal-analysis-of-competing-risks/ (2023).
- Discrete-time / churn / CLV — SAS https://support.sas.com/resources/papers/proceedings12/132-2012.pdf (2012), https://support.sas.com/resources/papers/proceedings/proceedings/sugi28/120-28.pdf (2003); Springer https://link.springer.com/article/10.1057/s41270-025-00450-2 (2025).
- ML survival — Ishwaran et al. RSF, Ann. Appl. Stat. 2008 https://ishwaran.org/papers/IKBL.AOAS.pdf (2008); Katzman et al. DeepSurv, BMC Med Res Methodol 2018 https://link.springer.com/article/10.1186/s12874-018-0482-1 / https://arxiv.org/abs/1606.00931 (2018).
