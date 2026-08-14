<!-- hub-reference-banner -->
> **Reference file — part of the `da-analytical-methods` hub.** Formerly the standalone `da-25-bayesian-data-analysis` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-25-bayesian-data-analysis
description: >-
  Applied Bayesian data analysis and probabilistic programming workflow — the
  end-to-end practice of building, fitting, checking, and comparing Bayesian
  models with modern PPLs. Covers the Gelman/Vehtari Bayesian workflow, prior
  choice (weakly-informative priors, prior predictive checks), MCMC (NUTS/HMC),
  variational inference (ADVI), hierarchical/multilevel models (partial pooling,
  non-centered reparameterization, Neal's funnel), posterior predictive checks,
  model comparison (LOO-CV/PSIS, WAIC), convergence diagnostics (rank-normalized
  R-hat, ESS, divergences, BFMI), PPLs (PyMC, Stan/CmdStanPy, NumPyro, Bambi),
  Bayesian regression/GLMs, and ArviZ for diagnostics and plotting.
  TRIGGER: user wants to build/fit/debug a Bayesian or probabilistic model;
  mentions PyMC, Stan, NumPyro, Bambi, brms, ArviZ, MCMC, NUTS, HMC, priors,
  posterior, LOO, WAIC, R-hat, divergences, hierarchical/multilevel/partial
  pooling, prior or posterior predictive checks, variational inference; asks how
  to choose priors, diagnose sampling problems, or compare Bayesian models.
  SKIP: pure probability theory or Bayes' theorem derivations with no modeling
  (use da-1-3-5-bayes-theorem); conceptual frequentist-vs-Bayesian comparison
  with no applied workflow (use da-1-4-3-frequentist-vs-bayesian-paradigms);
  generic ML / point-estimate regression with no posterior (use da-7-machine-learning
  or da-6-statistical-modeling).
---

# Bayesian Data Analysis & Probabilistic Programming

Applied Bayesian modeling: specify a generative model, fit the posterior with a
probabilistic programming language (PPL), interrogate it with predictive checks
and diagnostics, and compare alternatives. This skill is the *workflow and
tooling* layer — it assumes Bayes' theorem and the frequentist/Bayesian contrast
are already understood (see da-1-3-5, da-1-4-3).

## Overview

A Bayesian model combines a **prior** `p(θ)` and a **likelihood** `p(y|θ)` into a
**posterior** `p(θ|y) ∝ p(y|θ) p(θ)`. For all but trivial models the posterior is
intractable analytically, so we approximate it by sampling (MCMC) or optimization
(variational inference). The discipline is not "run the sampler and read the mean."
It is an iterative loop — model, fit, check, expand — formalized as the **Bayesian
workflow** ([Gelman, Vehtari, Simpson et al. 2020, arXiv:2011.01808](https://arxiv.org/abs/2011.01808)).

Use Bayesian methods when you want: full uncertainty quantification (posteriors,
not just point estimates), principled regularization via priors, partial pooling
across groups (hierarchical models), the ability to incorporate domain knowledge,
and propagation of uncertainty into predictions/decisions.

## Core Concepts

### 1. The Bayesian workflow (iterative, not linear)
Build a model → simulate from priors → fit → diagnose computation → check the model
against data → compare candidate models → expand or simplify → repeat. The workflow
explicitly treats *computational failure* (divergences, bad R-hat) and *model
misfit* (failed posterior predictive checks) as distinct problems with distinct
fixes ([Gelman et al. 2020](https://arxiv.org/abs/2011.01808); [Betancourt, "Towards a
Principled Bayesian Workflow" 2020](https://betanalpha.github.io/assets/case_studies/principled_bayesian_workflow.html);
[Gabry, Simpson, Vehtari et al., "Visualization in Bayesian Workflow," JRSS-A 2019](https://academic.oup.com/jrsssa/article/182/2/389/7070184)).

### 2. Priors and prior predictive checks
Priors encode pre-data information and regularize. Modern practice favors
**weakly-informative priors** (e.g. `Normal(0, 1)` on standardized predictors,
`HalfNormal`/`Exponential` on scales, `LKJ` on correlation matrices) over flat or
diffuse priors, which can cause heavy tails, poor geometry, and prior-data conflict.
Validate priors with a **prior predictive check**: draw θ from the prior, simulate
y, and confirm simulated data is physically plausible — before touching the
likelihood ([Stan Prior Choice Recommendations wiki, updated 2024](https://github.com/stan-dev/stan/wiki/Prior-Choice-Recommendations);
[PyMC "Prior and Posterior Predictive Checks" v5 docs, 2024](https://www.pymc.io/projects/docs/en/stable/learn/core_notebooks/posterior_predictive.html);
[Gelman et al. 2020](https://arxiv.org/abs/2011.01808)).

### 3. MCMC: HMC and NUTS
Hamiltonian Monte Carlo (HMC) uses gradient information to propose distant,
high-acceptance moves; the **No-U-Turn Sampler (NUTS)** auto-tunes trajectory
length and step size, and is the default in every modern PPL. NUTS is the right
default for continuous parameters; it cannot directly sample discrete parameters
(marginalize them out instead) ([Hoffman & Gelman, "The No-U-Turn Sampler," JMLR 2014](https://jmlr.org/papers/v15/hoffman14a.html);
[Betancourt, "A Conceptual Introduction to HMC," arXiv:1701.02434, 2017](https://arxiv.org/abs/1701.02434);
[Stan Reference Manual — MCMC, 2024](https://mc-stan.org/docs/reference-manual/mcmc.html)).

### 4. Variational inference (VI / ADVI)
VI approximates the posterior with a simpler parametric family by maximizing the
ELBO — fast and scalable, but it typically *underestimates* posterior variance and
can miss multimodality. **ADVI** automates this for differentiable models; normalizing
flows give richer approximations. Use VI for large data / prototyping; validate
against MCMC before trusting it ([Kucukelbir et al., "Automatic Differentiation
Variational Inference," JMLR 2017](https://jmlr.org/papers/v18/16-107.html);
[PyMC variational API docs, 2024](https://www.pymc.io/projects/docs/en/stable/api/vi.html);
[NumPyro SVI docs, 2024](https://num.pyro.ai/en/stable/svi.html)).

### 5. Hierarchical / multilevel models
Model grouped data with group-level parameters drawn from a shared population
distribution. This produces **partial pooling** — group estimates shrink toward the
global mean by an amount the data determines, balancing complete pooling (ignore
groups) and no pooling (independent per group). Hierarchical models are the single
biggest practical reason to go Bayesian ([Gelman & Hill, *Data Analysis Using
Regression and Multilevel/Hierarchical Models*, 2006](http://www.stat.columbia.edu/~gelman/arm/);
[McElreath, *Statistical Rethinking* 2nd ed., 2020](https://xcelab.net/rm/);
[PyMC "A Primer on Bayesian Methods for Multilevel Modeling," 2024](https://www.pymc.io/projects/docs/en/stable/learn/core_notebooks/GLM_hierarchical.html)).

### 6. Posterior predictive checks (PPC)
After fitting, simulate replicated data `y_rep` from the posterior predictive and
compare to observed `y` — overlaid densities, test statistics (min/max/quantiles),
and **LOO-PIT** for calibration. A model that cannot reproduce key features of the
data it was fit on is misspecified ([Gabry et al., JRSS-A 2019](https://academic.oup.com/jrsssa/article/182/2/389/7070184);
[ArviZ `plot_ppc` / `plot_loo_pit` API, 2024](https://python.arviz.org/en/stable/api/index.html);
[PyMC predictive-checks notebook, 2024](https://www.pymc.io/projects/docs/en/stable/learn/core_notebooks/posterior_predictive.html)).

### 7. Model comparison: LOO-CV (PSIS), WAIC
Compare predictive accuracy with **PSIS-LOO** — leave-one-out cross-validation made
cheap via Pareto-smoothed importance sampling — reported as `elpd_loo`. The
**Pareto-k** diagnostic flags observations where the importance-sampling estimate is
unreliable (`k > 0.7` = problematic). WAIC is a related estimator; LOO is generally
preferred because of its self-diagnostics. Bayes factors are sensitive to priors and
hard to compute — prefer LOO for most predictive-comparison work ([Vehtari, Gelman,
Gabry, "Practical Bayesian model evaluation using LOO-CV and WAIC," Stat. & Computing
2017, arXiv:1507.04544](https://arxiv.org/abs/1507.04544);
[Vehtari et al., "Pareto Smoothed Importance Sampling," JMLR 2024, arXiv:1507.02646](https://arxiv.org/abs/1507.02646);
[ArviZ `compare` / `loo` API, 2024](https://python.arviz.org/en/stable/api/index.html)).

### 8. Convergence diagnostics
Trust nothing until the sampler is verified. Use the **rank-normalized split-R-hat**
(target `< 1.01`, not the old 1.1) and **bulk-ESS / tail-ESS** (want hundreds+).
HMC adds geometry-specific diagnostics: **divergent transitions** (signal the
sampler can't resolve posterior curvature — almost always a reparameterization or
prior problem), **max-treedepth** hits (efficiency, not correctness), and **E-BFMI**
< 0.3 (poor energy exploration / heavy tails) ([Vehtari, Gelman, Simpson, Carpenter,
Bürkner, "Rank-normalization, folding, and localization: An improved R-hat,"
Bayesian Analysis 2021, arXiv:1903.08008](https://arxiv.org/abs/1903.08008);
[Betancourt, "A Conceptual Introduction to HMC" 2017](https://arxiv.org/abs/1701.02434);
[Stan "Runtime Warnings and Convergence Problems," 2024](https://mc-stan.org/misc/warnings.html)).

### 9. Probabilistic programming languages (PPLs)
- **PyMC** (Python, v5+): expressive `with pm.Model():` API on a PyTensor backend;
  pluggable NUTS samplers (default, `nutpie`, `numpyro`, `blackjax`).
- **Stan** (own language; via **CmdStanPy**/**cmdstanr**): the reference HMC/NUTS
  implementation, gold standard for hard models.
- **NumPyro** (JAX): fastest NUTS via JIT + GPU, great for large/hierarchical models.
- **Bambi** (Python on PyMC) / **brms** (R on Stan): formula interface (`y ~ x + (1|g)`)
  for regression/GLMMs — the fast path for standard models.
- **ArviZ**: backend-agnostic diagnostics/plots on the `InferenceData` xarray format —
  the shared output layer for all of the above.
([PyMC docs, 2024](https://www.pymc.io/); [CmdStanPy docs, 2024](https://mc-stan.org/cmdstanpy/);
[NumPyro docs, 2024](https://num.pyro.ai/); [Bambi docs, 2024](https://bambinos.github.io/bambi/);
[ArviZ docs, 2024](https://python.arviz.org/)).

### 10. Bayesian regression & GLMs
The link-function GLM family carries over directly: Gaussian (linear), Bernoulli/Binomial
(logistic), Poisson/Negative-Binomial (counts), with priors on coefficients providing
regularization (Bayesian ridge ≈ Normal prior, LASSO ≈ Laplace, horseshoe for sparse
selection). Robust regression swaps Normal for Student-t likelihood. Standardize
predictors so default priors are sensible ([McElreath, *Statistical Rethinking* 2020](https://xcelab.net/rm/);
[Bambi GLM examples, 2024](https://bambinos.github.io/bambi/notebooks/);
[Stan Prior Choice wiki, 2024](https://github.com/stan-dev/stan/wiki/Prior-Choice-Recommendations)).

### 11. ArviZ — the diagnostics & plotting hub
`InferenceData` groups posterior, prior, predictive, log-likelihood, and sample
stats. Key functions: `az.summary` (R-hat/ESS/HDI), `az.plot_trace`, `az.plot_posterior`,
`az.plot_ppc`, `az.loo` / `az.compare` / `az.plot_compare`, `az.plot_forest` (shrinkage),
`az.plot_energy` (BFMI), `az.plot_parallel` / `plot_pair` (divergences)
([ArviZ API reference, 2024](https://python.arviz.org/en/stable/api/index.html);
[ArviZ example gallery, 2024](https://python.arviz.org/en/stable/examples/index.html)).

## Methodology — the loop in practice

1. **Scope & generative story.** Write the model as a data-generating process. Decide
   likelihood family from the outcome type before anything else.
2. **Priors + prior predictive check.** Set weakly-informative priors on standardized
   variables; `sample_prior_predictive`; reject priors that generate absurd data.
3. **Fit.** Run NUTS with ≥4 chains. Default `target_accept=0.8`; raise to 0.9–0.99 if
   you see divergences. For large/continuous models try the NumPyro/nutpie backend.
4. **Diagnose computation.** Check R-hat < 1.01, bulk/tail-ESS, **zero divergences**,
   E-BFMI > 0.3, no mass treedepth saturation. Fix *here* before interpreting anything.
5. **Posterior predictive check.** `sample_posterior_predictive` → `az.plot_ppc`,
   test statistics, LOO-PIT calibration.
6. **Compare.** `az.loo` per model, `az.compare` for the table; inspect Pareto-k.
7. **Iterate.** Expand (add hierarchy/interactions/robust likelihood) or simplify based
   on checks; re-run the loop.

## Practical Patterns

- **Standardize continuous predictors** (mean 0, sd 1) so `Normal(0, 1)`-ish priors and
  default sampler settings behave; back-transform for interpretation.
- **Non-centered parameterization for hierarchical models.** Replace
  `θ_g ~ Normal(μ, σ)` with `θ_g = μ + σ · z_g, z_g ~ Normal(0, 1)` (the "Matt trick").
  This decouples scale from location and removes the funnel that causes divergences
  when groups have little data ([Stan User's Guide — Reparameterization, 2024](https://mc-stan.org/docs/stan-users-guide/efficiency-tuning.html#reparameterization);
  [Betancourt & Girolami 2015, arXiv:1312.0906](https://arxiv.org/abs/1312.0906)).
- **≥4 chains, multiple random seeds**, and inspect trace plots — R-hat alone can miss
  problems a visual catches.
- **Marginalize discrete parameters** out of the model so NUTS can run (Stan can't sample
  them; PyMC/NumPyro handle some via mixtures/enumeration but marginalization is cleaner).
- **Use a formula library (Bambi/brms) for standard GLMMs**; drop to raw PyMC/Stan only
  when the model needs custom structure.
- **Save the full `InferenceData`** (posterior + predictive + log-likelihood + sample
  stats) so LOO and PPCs are reproducible without refitting.

## Anti-Patterns

- **Ignoring divergences.** Divergent transitions mean the reported posterior is biased.
  Never report results from a run with divergences "because R-hat looks fine."
  Reparameterize, raise `target_accept`, or tighten priors.
- **Flat/diffuse "uninformative" priors as a default.** They are rarely uninformative on
  quantities of interest, induce bad geometry, and break weak-prior LOO. Prefer
  weakly-informative ([Stan Prior Choice wiki, 2024](https://github.com/stan-dev/stan/wiki/Prior-Choice-Recommendations)).
- **Trusting VI/ADVI without checking against MCMC.** VI understates variance; calibrate it.
- **Old R-hat < 1.1 as the bar.** Use rank-normalized R-hat < 1.01 and check ESS too
  ([Vehtari et al. 2021](https://arxiv.org/abs/1903.08008)).
- **Comparing models on in-sample fit / DIC / raw likelihood.** Use LOO/WAIC for
  out-of-sample predictive comparison; read Pareto-k before trusting LOO numbers.
- **Centered hierarchical parameterization with sparse groups.** This is the classic
  divergence/funnel generator — go non-centered.
- **Reading the posterior mean only.** Report intervals (HDI), and propagate the full
  posterior into downstream decisions.

## Troubleshooting

| Symptom | Likely cause | Fix |
| --- | --- | --- |
| Divergent transitions | Funnel / sharp curvature, often hierarchical | Non-centered param; raise `target_accept` to 0.9–0.99; tighten priors |
| R-hat > 1.01, low ESS | Non-convergence, multimodality, too few iters | More tuning/draws; better inits; reparameterize; check for label-switching |
| max-treedepth hits | Inefficient (highly correlated posterior), not wrong | Raise `max_treedepth`; reparameterize; standardize/decorrelate |
| E-BFMI < 0.3 | Poor energy exploration, heavy tails | Reparameterize; switch/regularize priors; consider Student-t |
| Pareto-k > 0.7 for some obs | Influential points; LOO importance sampling unreliable | Use `loo` moment-matching or exact refit for those points; inspect outliers |
| PPC density misses data | Wrong likelihood / missing structure | Change family (e.g. Negative-Binomial for overdispersed counts); add hierarchy/interactions |
| Prior predictive absurd | Priors too wide/wrong scale | Tighten; standardize predictors; rethink units |
| VI posterior too narrow vs MCMC | Mean-field underestimates variance | Full-rank ADVI / normalizing flows, or just use NUTS |

## References

1. Gelman, Vehtari, Simpson, Margossian, Carpenter, et al. — *Bayesian Workflow* (2020). https://arxiv.org/abs/2011.01808
2. Vehtari, Gelman, Gabry — *Practical Bayesian model evaluation using LOO-CV and WAIC* (Stat. & Computing 2017). https://arxiv.org/abs/1507.04544
3. Vehtari, Simpson, Gelman, Yao, Gabry — *Pareto Smoothed Importance Sampling* (JMLR 2024). https://arxiv.org/abs/1507.02646
4. Vehtari, Gelman, Simpson, Carpenter, Bürkner — *Rank-normalization, folding, and localization: An improved R-hat* (Bayesian Analysis 2021). https://arxiv.org/abs/1903.08008
5. Hoffman & Gelman — *The No-U-Turn Sampler* (JMLR 2014). https://jmlr.org/papers/v15/hoffman14a.html
6. Betancourt — *A Conceptual Introduction to Hamiltonian Monte Carlo* (2017). https://arxiv.org/abs/1701.02434
7. Betancourt — *Towards a Principled Bayesian Workflow* (case study, 2020). https://betanalpha.github.io/assets/case_studies/principled_bayesian_workflow.html
8. Betancourt & Girolami — *Hamiltonian Monte Carlo for Hierarchical Models* (2015). https://arxiv.org/abs/1312.0906
9. Kucukelbir, Tran, Ranganath, Gelman, Blei — *Automatic Differentiation Variational Inference* (JMLR 2017). https://jmlr.org/papers/v18/16-107.html
10. Gabry, Simpson, Vehtari, Betancourt, Gelman — *Visualization in Bayesian Workflow* (JRSS-A 2019). https://academic.oup.com/jrsssa/article/182/2/389/7070184
11. Gelman, Carlin, Stern, Dunson, Vehtari, Rubin — *Bayesian Data Analysis*, 3rd ed. (BDA3, 2013). http://www.stat.columbia.edu/~gelman/book/
12. McElreath — *Statistical Rethinking*, 2nd ed. (2020). https://xcelab.net/rm/
13. Gelman & Hill — *Data Analysis Using Regression and Multilevel/Hierarchical Models* (2006). http://www.stat.columbia.edu/~gelman/arm/
14. Stan — *Prior Choice Recommendations* wiki (2024). https://github.com/stan-dev/stan/wiki/Prior-Choice-Recommendations
15. Stan — *Reference Manual (MCMC)* & *Runtime Warnings* (2024). https://mc-stan.org/docs/reference-manual/mcmc.html · https://mc-stan.org/misc/warnings.html
16. PyMC — *Docs & predictive-checks / hierarchical notebooks* v5 (2024). https://www.pymc.io/
17. NumPyro — *Docs, SVI, Neal's Funnel example* (2024). https://num.pyro.ai/
18. Bambi — *Docs & GLM examples* (2024). https://bambinos.github.io/bambi/
19. CmdStanPy — *Docs* (2024). https://mc-stan.org/cmdstanpy/
20. ArviZ — *API reference & example gallery* (2024). https://python.arviz.org/
