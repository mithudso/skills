<!-- hub-reference-banner -->
> **Reference file — part of the `da-analytical-methods` hub.** Formerly the standalone `da-6-statistical-modeling` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-6-statistical-modeling
description: Statistical modeling for data analysis — inference-first treatment of linear regression (OLS, Gauss-Markov, diagnostics), logistic regression (logit link, odds ratios, ROC/AUC, Hosmer-Lemeshow), GLMs (Poisson, negative binomial, gamma, link functions, overdispersion), regularization (Ridge, Lasso, Elastic Net), classification (LDA/QDA, naive Bayes, decision trees), clustering (k-means, hierarchical, DBSCAN, GMM, silhouette), time series (ARIMA/SARIMA, exponential smoothing, Prophet, state-space, VAR), survival analysis (Kaplan-Meier, Cox PH, censoring, hazard ratios), mixed-effects models (fixed vs random effects, lme4, ICC, REML), model selection (AIC, BIC, CV, nested CV), model interpretation (coefficients, marginal effects, partial dependence, SHAP, LIME), and Bayesian modeling (Stan/PyMC, MCMC, posterior predictive checks). TRIGGER when the user asks about fitting a model with interpretable coefficients, choosing between OLS and a GLM, diagnosing residuals, interpreting odds ratios or hazard ratios, picking between Ridge/Lasso/Elastic Net, checking Cox PH assumptions, comparing AIC/BIC across models, running posterior predictive checks, or any inference-driven modeling decision. SKIP when the task is pure predictive ML benchmarking with no inference component (use da-7-machine-learning), A/B testing or causal identification (use da-12-ab-testing-causal-inference), foundational sampling/estimation theory without a fitted model (use da-1-4-statistical-inference-foundations), or building a deep-learning architecture.
when_to_use:
  - Fitting a linear regression and validating OLS assumptions (linearity, homoscedasticity, normality, independence, no multicollinearity)
  - Interpreting an odds ratio, deviance, ROC/AUC, or Hosmer-Lemeshow result from logistic regression
  - Choosing between Poisson, quasi-Poisson, and negative binomial for count data
  - Choosing between Ridge, Lasso, and Elastic Net for a regularized regression
  - Running residual diagnostics (Q-Q plot, scale-location, Cook's distance, leverage, VIF)
  - Picking between LDA, QDA, naive Bayes, and a decision tree for a classification problem with interpretability goals
  - Selecting k for k-means, choosing DBSCAN epsilon, comparing GMM vs k-means, validating with silhouette
  - Fitting ARIMA/SARIMA, exponential smoothing, Prophet, state-space, or VAR for a time-series problem
  - Running Kaplan-Meier, building a Cox proportional hazards model, handling censored data, checking PH assumption with Schoenfeld residuals
  - Choosing fixed vs random effects, fitting lme4 / nlme / statsmodels mixed models, interpreting ICC and variance components
  - Comparing models with AIC, BIC, k-fold CV, or nested CV; avoiding nested-CV vs CV optimism bias
  - Explaining a model with marginal effects, partial dependence plots, SHAP, or LIME
  - Building a Bayesian regression in Stan / PyMC / brms, running MCMC, diagnosing R-hat and ESS, running posterior predictive checks
related_skills:
  - da-1-4-statistical-inference-foundations
  - da-7-machine-learning
  - da-12-ab-testing-causal-inference
  - da-1-3-probability-theory
  - da-2-2-2-hypothesis-formulation
tags:
  - data-analysis
  - statistical-modeling
  - regression
  - glm
  - bayesian
  - inference
keywords:
  - linear regression
  - OLS
  - Gauss-Markov
  - logistic regression
  - logit
  - odds ratio
  - GLM
  - Poisson regression
  - negative binomial
  - gamma regression
  - link function
  - overdispersion
  - Ridge
  - Lasso
  - Elastic Net
  - LDA
  - QDA
  - naive Bayes
  - decision tree
  - k-means
  - DBSCAN
  - hierarchical clustering
  - GMM
  - silhouette
  - ARIMA
  - SARIMA
  - exponential smoothing
  - Prophet
  - state-space
  - VAR
  - Kaplan-Meier
  - Cox proportional hazards
  - censoring
  - hazard ratio
  - mixed-effects
  - lme4
  - ICC
  - REML
  - AIC
  - BIC
  - cross-validation
  - nested CV
  - marginal effects
  - partial dependence
  - SHAP
  - LIME
  - Stan
  - PyMC
  - MCMC
  - posterior predictive check
---

# Statistical Modeling

Statistical modeling sits between descriptive statistics (what the data look like) and prediction-first machine learning (how well a model forecasts new data). The defining commitment is **inference**: the goal is to estimate parameters with calibrated uncertainty, test assumptions, and produce interpretable effects that other humans can argue with. Predictive accuracy matters, but it is not the only criterion — a model whose coefficients you cannot defend is not a statistical model in the sense used here.

This skill emphasizes the assumptions, diagnostics, and interpretation steps that distinguish careful inference from black-box prediction. For predictive performance benchmarking, ensembles, neural networks, and pure ML pipelines, use `da-7-machine-learning`. For randomized experiments and causal identification, use `da-12-ab-testing-causal-inference`. For the underlying probability and estimation theory, use `da-1-3-probability-theory` and `da-1-4-statistical-inference-foundations`.

---

## 1. Linear Regression — OLS, Assumptions, Diagnostics

Ordinary Least Squares (OLS) is the workhorse linear model: y = Xβ + ε, fit by minimizing the sum of squared residuals. Under the **Gauss-Markov assumptions**, the OLS estimator is BLUE — the **Best Linear Unbiased Estimator** — meaning it has the smallest variance among all linear unbiased estimators of β.

### The five Gauss-Markov assumptions

1. **Linearity in parameters** — the model is linear in β (not necessarily in x; x² is fine).
2. **Random sampling** — observations are an i.i.d. draw from the population (or at least exchangeable).
3. **No perfect multicollinearity** — the design matrix X has full column rank; no predictor is a perfect linear combination of others.
4. **Zero conditional mean / exogeneity** — E[ε | X] = 0; the error has mean zero given any value of the regressors.
5. **Homoscedasticity** — Var(ε | X) = σ², constant across the range of X.

A sixth assumption — **normality of errors** — is *not* required for BLUE, but is required for the usual t- and F-tests to follow exact small-sample distributions. With large n, the Central Limit Theorem rescues inference even when errors are non-normal.

### Residual diagnostics (the four essential plots)

| Plot | What it checks | Red flag |
| --- | --- | --- |
| Residuals vs Fitted | Linearity, mean-zero errors | Curvature, funneling |
| Q-Q plot of residuals | Normality | Heavy tails, skew at extremes |
| Scale-Location (√|standardized residual| vs fitted) | Homoscedasticity | Upward-sloping trend |
| Residuals vs Leverage | Influential observations | Points outside Cook's distance contours |

### Influence measures

- **Leverage (h_ii)** — diagonal of the hat matrix H = X(X'X)⁻¹X'. Measures how unusual an observation's *x* values are. High-leverage rule of thumb: h_ii > 2p/n or 3p/n where p is the number of parameters.
- **Standardized residual** — residual divided by its estimated standard error; |r_i| > 2 or 3 flags an outlier in y.
- **Cook's distance (D_i)** — combines leverage and residual size: D_i = (r_i² / p) · (h_ii / (1 - h_ii)). Combines how unusual the observation is in X with how much its residual is. Common thresholds: D_i > 4/n or D_i > 1.
- **DFBETAS** — per-coefficient influence; how much β̂_j changes when observation i is dropped.
- **DFFITS** — per-fit influence; how much the prediction at i changes when i is dropped.

### Multicollinearity

- **Variance Inflation Factor (VIF)** — VIF_j = 1 / (1 - R²_j) where R²_j is from regressing predictor j on the others. VIF > 5 is suspicious; VIF > 10 is usually addressed (drop, combine, regularize, or center).
- **Condition number** of X'X — large values (>30) signal numerical instability.

### When OLS breaks and the fix

| Problem | Symptom | Fix |
| --- | --- | --- |
| Non-linearity | Curved residuals vs fitted | Polynomial term, spline, GAM, transform y |
| Heteroscedasticity | Funnel in scale-location | Heteroscedasticity-consistent (HC) standard errors (Huber-White, "robust SEs"), weighted least squares, transform y, GLM with appropriate variance |
| Non-normal errors | Heavy Q-Q tails | Robust regression (Huber M-estimator, MM-estimator), bootstrap CIs, GLM with non-Gaussian family |
| Autocorrelation | Durbin-Watson ≠ 2; pattern in residuals over time | Time-series model (ARIMA), Newey-West HAC standard errors, mixed model with autoregressive structure |
| Multicollinearity | High VIF, unstable coefficients | Drop a predictor, combine, regularize (Ridge), or PCA on predictors |
| Influential points | High Cook's D | Investigate (data error vs real); refit with and without; report sensitivity; consider robust regression |
| Endogeneity | E[ε|X] ≠ 0 | Instrumental variables (2SLS), causal identification — see da-12 |

### Sources
- [Gauss-Markov Theorem & BLUE — Statistics By Jim](https://statisticsbyjim.com/regression/gauss-markov-theorem-ols-blue/)
- [Linear Regression Assumptions and Diagnostics in R — STHDA](https://www.sthda.com/english/articles/39-regression-model-diagnostics/161-linear-regression-assumptions-and-diagnostics-in-r-essentials/)
- [Gauss-Markov Assumptions — Michael Brenndoerfer](https://mbrenndoerfer.com/writing/gauss-markov-assumptions-linear-regression-ols-blue-estimator)
- [Gauss-Markov Theorem — Wikipedia](https://en.wikipedia.org/wiki/Gauss%E2%80%93Markov_theorem)
- [Regression Diagnostics — Thomas Love](https://thomaselove.github.io/431notes-2017/regression-diagnostics.html)

---

## 2. Logistic Regression — Link, Odds Ratios, ROC/AUC

When the outcome is binary (0/1, success/failure, churn/retain), OLS is inappropriate — predicted probabilities can go below 0 or above 1, residuals are necessarily heteroscedastic, and the conditional distribution is Bernoulli, not Gaussian. **Logistic regression** is the GLM whose link is the logit and whose family is the binomial.

### The model

logit(p) = log(p / (1 - p)) = β₀ + β₁x₁ + … + β_k x_k

The logit (log-odds) maps p ∈ (0, 1) onto the real line, so the right-hand side can be unconstrained linear in the parameters. The inverse — the **logistic function** σ(z) = 1 / (1 + e^(-z)) — recovers the probability.

### Odds ratio interpretation (the load-bearing skill)

For a one-unit increase in x_j, the **odds** of the outcome multiply by exp(β_j). That multiplier is the **odds ratio (OR)**.

- exp(β_j) > 1 → x_j increases the odds of the outcome
- exp(β_j) = 1 → no association
- exp(β_j) < 1 → x_j decreases the odds

A common error: confusing odds ratios with probability ratios. An OR of 2 does **not** mean the probability doubles — only the odds do. When the base rate is low (<10%), OR ≈ risk ratio; when the base rate is high, they diverge sharply.

For a continuous predictor, "a one-unit change" is often a poor scale. Report ORs per clinically/operationally meaningful unit (e.g., per SD, per decade of age).

### Fit and inference

- **Maximum likelihood estimation** — no closed form; solved by iteratively reweighted least squares (IRLS) / Newton-Raphson. Quasi-separation (a predictor that perfectly classifies the outcome) causes |β̂| → ∞ — regularize, drop the predictor, or use Firth's penalized likelihood.
- **Deviance** = -2 · log-likelihood; the GLM analog of residual sum of squares. The difference in deviance between nested models is χ²-distributed with df equal to the number of additional parameters.
- **Wald test** — β̂ / SE(β̂); asymptotically z-distributed. Be cautious with the Hauck-Donner phenomenon: as |β| → ∞, the Wald statistic can actually shrink.
- **Likelihood ratio test** — preferred over Wald for small samples; compares deviances.

### Calibration and discrimination

Two distinct concepts, both required:

- **Discrimination** — how well does the model order subjects from low to high risk? Measured by the **ROC curve** (sensitivity vs 1-specificity across thresholds) and its area, **AUC**. AUC = 0.5 is chance; AUC = 1.0 is perfect. AUC = 0.7-0.8 is acceptable, 0.8-0.9 is excellent, >0.9 is outstanding.
- **Calibration** — when the model says 30% risk, do ~30% of those cases actually occur? Measured by the **Hosmer-Lemeshow test** (which divides observations into deciles of predicted risk and compares observed vs expected counts), the **calibration plot**, or **Brier score** (mean squared error between predicted probabilities and outcomes).

Hosmer-Lemeshow caveats: low power with small n; over-powered to detect trivial miscalibration with very large n; sensitive to the number of groups (default 10 is conventional but arbitrary).

### Multinomial and ordinal extensions

- **Multinomial logistic regression** — outcome with K > 2 unordered categories; baseline-category logits.
- **Ordinal logistic regression** — outcome has natural ordering (e.g., Likert); the proportional odds model (cumulative logits with a common slope) is standard. Test the proportional odds assumption with the score test or Brant test.

### Sources
- [Logistic Regression — rinterested.github.io](https://rinterested.github.io/statistics/logistic_regression.html)
- [UCLA FAQ: How do I interpret odds ratios? — UCLA OARC](https://stats.oarc.ucla.edu/other/mult-pkg/faq/general/faq-how-do-i-interpret-odds-ratios-in-logistic-regression/)
- [Logistic Regression Multiple — StatsDirect](https://www.statsdirect.com/help/regression_and_correlation/logistic.htm)
- [Assessment of Logistic Regression Models — Shang/AmStat](https://ww2.amstat.org/meetings/proceedings/2019/data/assets/pdf/1199666.pdf)
- [Logistic Regression — MedCalc](https://www.medcalc.org/en/manual/logistic-regression.php)

---

## 3. Generalized Linear Models — Poisson, Negative Binomial, Gamma

GLMs (Nelder & Wedderburn, 1972) unify regression for any response from the exponential family by combining:

1. A **random component** — the distribution of Y (Gaussian, Bernoulli, binomial, Poisson, gamma, inverse Gaussian, negative binomial as approximation).
2. A **linear predictor** — η = Xβ.
3. A **link function** — g(μ) = η, where μ = E[Y].

| Family | Use for | Canonical link | Variance function |
| --- | --- | --- | --- |
| Gaussian | continuous, unbounded, symmetric | identity | σ² |
| Binomial | binary, proportion | logit | μ(1-μ) |
| Poisson | counts (rate-based) | log | μ |
| Quasi-Poisson | overdispersed counts (mean-variance proportional) | log | φ·μ |
| Negative binomial | overdispersed counts (quadratic mean-variance) | log | μ + μ²/θ |
| Gamma | positive continuous, right-skewed (cost, duration) | log or inverse | φμ² |
| Inverse Gaussian | positive continuous, heavy right tail | 1/μ² | φμ³ |

### Poisson regression

Models the rate of events given an exposure (population, time at risk). The log link is canonical: log(μ) = Xβ. To handle different exposure, include log(exposure) as an **offset** (a covariate with fixed coefficient 1).

Hard assumption: **mean equals variance** (E[Y] = Var(Y) = μ). Real data rarely obey this — see overdispersion.

### Overdispersion

When Var(Y) > E[Y], standard errors from a naive Poisson fit are too small, p-values are too small, and CIs are too tight. The Pearson dispersion statistic φ̂ = Σ(Pearson residuals)² / (n - p) should be ~1; values >1.5 are concerning, >2 require action.

Causes of overdispersion include unmodeled heterogeneity (omitted variables), clustering (use mixed model), zero inflation, and contagion (events are not independent).

Three standard remedies:

- **Quasi-Poisson** — inflate standard errors by √φ̂; coefficients are unchanged. Easy, but you lose a proper likelihood (no AIC, no LRT in the usual sense; QAIC is available).
- **Negative binomial** — adds a dispersion parameter θ; variance = μ + μ²/θ. Recovers a proper likelihood, allows quadratic mean-variance, and is preferred when AIC/LRT/BIC matter.
- **Zero-inflated** (ZIP/ZINB) or **hurdle** models — when excess zeros come from a separate process (e.g., "did not engage at all" vs "engaged k times").

The **Cameron-Trivedi overdispersion test** formally tests Var = μ vs Var = μ + α·g(μ).

### Negative binomial — two flavors

- **NB1** — variance = μ(1 + α); linear in μ.
- **NB2** — variance = μ + αμ²; quadratic in μ. This is the default in `MASS::glm.nb()` and `statsmodels`.

### Gamma regression

For strictly positive, right-skewed continuous outcomes — insurance claim size, hospital length of stay, time-to-event with no censoring, semiconductor yield. The log link is most interpretable; the inverse link is canonical but rarely interpretable. The shape parameter ν controls skewness; ν → ∞ approaches Gaussian.

Diagnostics for GLMs use **deviance residuals** and **Pearson residuals**; neither should show structure against fitted values or predictors. Half-normal plots with simulated envelopes (DHARMa in R, scipy/pingouin in Python) are the practical gold standard for GLM residual checking.

### Sources
- [GLM Practical Guide — MCP Analytics](https://mcpanalytics.ai/articles/general-linear-models-glm-practical-guide-for-data-driven-decisions)
- [Poisson vs Negative Binomial Count Models — Rezaee/UGA](https://babakrezaee.github.io/SU_POLS537/Count_04242019.html)
- [Overdispersion in GLMMs — ResearchGate Q&A thread](https://www.researchgate.net/post/How_to_deal_with_overdispersion_in_Generalized_linear_mixed_models_in_R)
- [Generalized Linear Models chapter — Fisher](https://tjfisher19.github.io/introStatModeling/generalized-linear-models.html)
- [Handling Overdispersion — CAS Actuarial Forum](https://www.casact.org/sites/default/files/database/forum_07wforum_07w109.pdf)

---

## 4. Regularization — Ridge, Lasso, Elastic Net

When p (predictors) is large relative to n, or predictors are correlated, OLS is unstable: variance explodes, coefficients flip signs across samples, and out-of-sample predictions degrade. **Regularization** trades a small amount of bias for a large reduction in variance by adding a penalty to the loss:

### Ridge regression (L2)

minimize: ‖y - Xβ‖² + λ‖β‖²₂   where ‖β‖²₂ = Σ β_j²

- **Shrinks coefficients toward zero but never exactly to zero** — no automatic variable selection.
- **Stable under multicollinearity** — handles correlated predictors gracefully, distributing the effect across them.
- **Closed-form solution**: β̂ = (X'X + λI)⁻¹ X'y. Always invertible even when X'X is singular.
- **Best when** all predictors are relevant and the issue is multicollinearity or p ≈ n.

### Lasso regression (L1)

minimize: ‖y - Xβ‖² + λ‖β‖₁   where ‖β‖₁ = Σ |β_j|

- **Drives some coefficients exactly to zero** — performs automatic variable selection.
- **Unstable under correlated predictors** — tends to pick one from a correlated group arbitrarily.
- **No closed form** — solved by coordinate descent or LARS. `glmnet` in R, `sklearn.linear_model.Lasso` in Python.
- **Best when** many predictors are irrelevant and you want a sparse model.

### Elastic Net (L1 + L2)

minimize: ‖y - Xβ‖² + λ[α‖β‖₁ + ((1-α)/2)‖β‖²₂]

- **Combines selection (Lasso) with stability under correlation (Ridge)**.
- **α controls the mix** — α=1 is Lasso, α=0 is Ridge.
- **Best when** you have many correlated predictors and want sparsity — Lasso would pick one arbitrarily; Elastic Net keeps the group.

### Choosing λ (and α for Elastic Net)

- **k-fold cross-validation** is standard — `cv.glmnet()` in R, `LassoCV/RidgeCV/ElasticNetCV` in sklearn.
- Two conventional choices: **λ_min** (minimum CV error) and **λ_1se** (largest λ within 1 SE of the minimum — sparser, more parsimonious).
- Standardize predictors before fitting — the penalty is scale-dependent.
- Don't penalize the intercept.

### Inference under regularization (the trap)

Regularized coefficients are **biased** by construction. Standard confidence intervals and p-values do not apply directly. For inference:

- **Bootstrap** — sample with replacement, refit, build percentile intervals (still biased; honest only with correction).
- **Selective inference** — `selectiveInference` R package implements post-selection p-values for Lasso (Lee, Sun, Sun, Taylor).
- **Debiased / desparsified Lasso** — Zhang & Zhang, van de Geer et al. — undoes the bias to enable CI construction.
- **Cross-validated R²** for predictive performance reporting, not inference.

Treat regularized regressions as predictive tools or feature-selection tools; if you need defensible coefficients, refit unpenalized OLS/GLM on the Lasso-selected variables ("relaxed Lasso") and use those standard errors — with the caveat that this ignores selection uncertainty.

### Sources
- [Regularization Lasso vs Ridge vs Elastic Net — GeeksforGeeks](https://www.geeksforgeeks.org/machine-learning/lasso-vs-ridge-vs-elastic-net-ml/)
- [Simplifying Regularization Ridge to Elastic Net — DigitalOcean](https://www.digitalocean.com/community/tutorials/regularization-in-machine-learning-lasso-ridge-elastic-net)
- [Elastic Net and Comparison of Regularization Methods — Fiveable](https://fiveable.me/modern-statistical-prediction-and-machine-learning/unit-7/elastic-net-comparison-regularization-methods/study-guide/wh7mOO8V82v2MWbg)
- [Regularization in R Tutorial — DataCamp](https://www.datacamp.com/tutorial/tutorial-ridge-lasso-elastic-net)
- [Lasso and Elastic Net Visual Guide — Towards Data Science](https://towardsdatascience.com/lasso-and-elastic-net-regressions-explained-a-visual-guide-with-code-examples-5fecf3e1432f/)

