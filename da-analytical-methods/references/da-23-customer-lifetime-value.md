<!-- hub-reference-banner -->
> **Reference file — part of the `da-analytical-methods` hub.** Formerly the standalone `da-23-customer-lifetime-value` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-23-customer-lifetime-value
description: >-
  Probabilistic / statistical Customer Lifetime Value (CLV) modeling — the
  "buy-till-you-die" (BTYD) family. Covers Pareto/NBD, BG/NBD, MBG/NBD,
  Gamma-Gamma monetary model, sBG and BG/BB for contractual/discrete settings,
  RFM as model inputs, discounted expected residual transactions (DERT),
  predictive vs historical CLV, cohort-based CLV, the CAC:LTV ratio, and the
  lifetimes / CLVTools / PyMC-Marketing implementations. Part of the da-*
  data-analytics curriculum (extends da-1..da-22).
  TRIGGER: questions about predicting future customer value or purchases,
  "buy till you die", BG/NBD / Pareto/NBD / Gamma-Gamma / sBG / BG/BB,
  expected residual transactions / DERT, churn-probability from
  recency/frequency, fitting CLV with lifetimes / PyMC-Marketing / CLVTools,
  LTV:CAC ratio, contractual vs non-contractual / discrete vs continuous CLV,
  predictive vs historical CLV, cohort LTV curves.
  SKIP: cohort-retention math, retention-curve fitting, and cohort decay
  tables in isolation → da-34-cohort-retention-analytics; general survival
  analysis (Kaplan-Meier, Cox PH, hazard/time-to-event) not framed as CLV →
  da-24-survival-analysis; generic Bayesian inference, MCMC diagnostics, or
  PyMC how-to with no CLV model → da-25-bayesian-data-analysis; pure A/B-test
  causal lift → da-12; marketing-mix / channel-attribution spend modeling →
  da-22; generic ML regression with no BTYD/retention structure → da-7.
---

# Customer Lifetime Value Modeling (Probabilistic / BTYD)

## Overview

Customer Lifetime Value (CLV) is the present value of the future cash flows
attributed to a customer relationship. This skill covers the **probabilistic
"buy-till-you-die" (BTYD) family** — statistical models that decompose CLV into
(1) *how often* a customer transacts while active, (2) *whether/when* they
silently churn, and (3) *how much* they spend per transaction — then discount
the expected future stream to present value.

Two orthogonal axes define the model landscape (Fader/Hardie taxonomy):

| | **Non-contractual** (churn unobserved) | **Contractual** (churn observed at renewal) |
| --- | --- | --- |
| **Continuous time** | Pareto/NBD, BG/NBD, MBG/NBD (+ Gamma-Gamma for spend) | survival models → **da-24** |
| **Discrete time** | BG/BB (beta-geometric / beta-Bernoulli) | sBG (shifted-beta-geometric) |

Choosing the wrong quadrant is the #1 modeling error. Subscriptions/SaaS are
**contractual** (you see the cancellation) → sBG / survival. Retail, e-commerce,
donations are **non-contractual** (you infer churn) → Pareto/NBD family.

Authoritative source corpus: Bruce Hardie's notes (brucehardie.com), the
Fader/Hardie/Lee Marketing Science papers, and the three reference
implementations — `lifetimes` (Python, archived), `CLVTools` (R), and
`PyMC-Marketing` (Python, Bayesian, the active successor).

## Core Concepts

### 1. The buy-till-you-die (BTYD) framework
A customer is "alive" until an unobserved dropout, transacting stochastically
while alive. Models pair a **counting process** (transactions while alive) with
a **timing process** (lifetime/dropout), each with cross-customer heterogeneity.
First introduced by Schmittlein, Morrison & Colombo, "Counting Your Customers:
Who Are They and What Will They Do Next?", *Management Science* 33(1):1–24 (1987)
([INFORMS](https://pubsonline.informs.org/doi/10.1287/mnsc.33.1.1)). Lineage and
taxonomy: ([Retina.ai — History of BTYD](https://retina.ai/academy/lesson/history-of-buy-til-you-die-btyd-models/), 2023);
([Wikipedia — Buy Till you Die](https://en.wikipedia.org/wiki/Buy_Till_you_Die)).

### 2. Pareto/NBD
The original non-contractual continuous-time model. **NBD** (Poisson–gamma
mixture) for transaction counts while alive; **Pareto** (exponential–gamma
mixture) for the unobserved lifetime. Four parameters (r, α, s, β). Powerful but
numerically awkward (Gaussian hypergeometric functions), which motivated BG/NBD.
([Schmittlein et al. 1987](https://pubsonline.informs.org/doi/10.1287/mnsc.33.1.1));
([CLVTools `pnbd`](https://www.clvtools.com/reference/pnbd.html));
([PyMC-Marketing Pareto/NBD notebook](https://www.pymc-marketing.io/en/latest/notebooks/clv/pareto_nbd.html)).

### 3. BG/NBD ("Counting Your Customers the Easy Way")
The workhorse. Replaces Pareto's continuous dropout with a **beta-geometric**
story: a customer flips a coin to churn *immediately after each transaction*
(prob. p, beta-distributed across customers); active counts are NBD. Far easier
to fit (estimable in Excel), nearly identical predictive accuracy. Fader, Hardie
& Lee, *Marketing Science* 24(2):275–284 (2005)
([Hardie PDF](http://brucehardie.com/papers/018/fader_et_al_mksc_05.pdf);
[INFORMS](https://pubsonline.informs.org/doi/abs/10.1287/mksc.1040.0098));
step-by-step derivation ([Fader/Hardie note 039, 2019](http://www.brucehardie.com/notes/039/bgnbd_derivation__2019-11-06.pdf));
([PyMC-Marketing BG/NBD notebook](https://www.pymc-marketing.io/en/stable/notebooks/clv/bg_nbd.html)).

> **Quirk:** in BG/NBD a customer cannot churn until *after* their first repeat
> purchase, so it understates the share of one-and-done (zero-repeat) customers.
> That gap is exactly what MBG/NBD fixes.

### 4. MBG/NBD (Modified BG/NBD)
Adds a dropout opportunity at time zero (right after the first purchase), so
customers who never repeat can be "dead". Estimates of expected repeat purchases
are nearly identical to BG/NBD, but the alive/dead classification of zero-repeat
customers is more realistic. Batislam, Denizel & Filiztekin, "Empirical
validation and comparison of models for customer base analysis", *IJRM* 24(3)
(2007) ([ScienceDirect](https://www.sciencedirect.com/science/article/abs/pii/S0167811607000171);
[Erratum 2008](https://www.sciencedirect.com/science/article/abs/pii/S0167811608000372));
implemented as `ModifiedBetaGeoModel`
([PyMC-Marketing](https://www.pymc-marketing.io/en/stable/notebooks/clv/clv_quickstart.html)).

### 5. Gamma-Gamma monetary model
Separately models **spend per transaction** (the frequency models above only
predict *counts*). Three assumptions: (a) a transaction's value varies randomly
around the customer's mean; (b) mean spend varies across customers but not over
time; (c) spend is independent of the transaction process — **so you must verify
frequency and monetary value are roughly uncorrelated before trusting it.** Fit
only on repeat purchasers. Fader, Hardie & Lee, "RFM and CLV: Using Iso-Value
Curves for Customer Base Analysis", *JMR* 42(4):415–430 (2005)
([Hardie PDF](https://www.brucehardie.com/papers/rfm_clv_2005-02-16.pdf);
[SAGE](https://journals.sagepub.com/doi/10.1509/jmkr.2005.42.4.415));
derivation note ([Fader/Hardie note 025](https://www.brucehardie.com/notes/025/gamma_gamma.pdf));
([CLVTools `gg`](https://www.clvtools.com/reference/gg.html)).

### 6. RFM as model inputs (sufficient statistics)
BTYD models do **not** need the raw transaction log — only per-customer
**Recency, Frequency, and "T"**, because R and F are *sufficient statistics* for
the likelihood. Conventions (note the field-specific definitions — easy to get
wrong):
- **frequency** = number of *repeat* purchases (total transactions − 1).
- **recency** = time between the customer's **first and last** purchase (NOT
  time since last purchase, which is the marketing-RFM convention).
- **T** = customer "age" = time from first purchase to end of observation.
- **monetary_value** = average value of repeat transactions (Gamma-Gamma input).

([Fader/Hardie/Lee RFM-CLV 2005](https://www.brucehardie.com/papers/rfm_clv_2005-02-16.pdf));
([PyMC-Marketing `rfm_summary` quickstart](https://www.pymc-marketing.io/en/stable/notebooks/clv/clv_quickstart.html));
([Ruberts — Modelling CLV with lifetimes](https://antonsruberts.github.io/lifetimes-CLV/)).

### 7. Discounted Expected Residual Transactions (DERT) → CLV
CLV in the non-contractual setting = **(expected spend per transaction from
Gamma-Gamma) × DERT**, where DERT is the present value of all expected future
transactions discounted to the end of the calibration period (the integral from
T to ∞). Use a **continuously-compounded** discount rate (convert annual rates,
e.g. 15%/yr ≈ 0.0027/week). Fader/Hardie originally called this DET.
([Fader/Hardie/Lee RFM-CLV 2005](https://www.brucehardie.com/papers/rfm_clv_2005-02-16.pdf));
([CLVTools `pnbd_DERT`](https://www.clvtools.com/reference/pnbd_DERT.html));
([Fader/Hardie "What's Wrong With This CLV Formula?" note 033](http://www.brucehardie.com/notes/033/what_is_wrong_with_this_CLV_formula.pdf)).

### 8. sBG — shifted-beta-geometric (contractual / discrete churn)
For **subscriptions/contractual** settings: each period a customer renews with
prob. θ or cancels with prob. 1−θ; θ is fixed per customer and beta-distributed
across the base. Projects a few observed retention numbers into a full survival
curve and (crucially) explains the **observed rise in aggregate retention rate
over time** as a heterogeneity sorting effect, not behavior change. Fader &
Hardie, "How to Project Customer Retention", *Journal of Interactive Marketing*
21(1):76–90 (2007) ([ResearchGate](https://www.researchgate.net/publication/264466889_How_to_Project_Customer_Retention));
extended for multi-cohort valuation in "Customer-Base Valuation in a Contractual
Setting: The Perils of Ignoring Heterogeneity", *Marketing Science* 29(1):85–93
(2010) ([Hardie PDF](http://brucehardie.com/papers/022/fader_hardie_mksc_10.pdf));
([PyMC-Marketing `ShiftedBetaGeoModel`](https://www.pymc-marketing.io/en/stable/notebooks/clv/sBG.html)).

### 9. BG/BB — discrete-time non-contractual
The discrete-time analog of Pareto/NBD: transactions per period are Bernoulli
(buy/no-buy) instead of Poisson, paired with a beta-geometric dropout — used for
"transaction opportunities" data (e.g., annual donations, periodic catalog
buyers). Closed-form, spreadsheet-friendly. Fader, Hardie & Shang,
"Customer-Base Analysis in a Discrete-Time Noncontractual Setting", *Marketing
Science* 29(6):1086–1108 (2010)
([Hardie PDF](http://www.brucehardie.com/papers/020/fader_et_al_mksc_10.pdf);
[INFORMS](https://pubsonline.informs.org/doi/10.1287/mksc.1100.0580));
([lifetimes `BetaGeoBetaBinomFitter`](https://lifetimes.readthedocs.io/en/latest/lifetimes.fitters.html)).

### 10. Predictive vs. historical CLV
**Historical CLV** sums realized past margin (backward-looking; cheap but
ignores future behavior). **Predictive CLV** forecasts future value via models
(BTYD, ML, or the naive ARPU/churn formula). The naive `ARPU ÷ churn` shortcut
silently assumes a single constant retention rate — which Fader/Hardie show is
biased low when retention is heterogeneous. Prefer model-based predictive CLV
with uncertainty intervals.
([TDS — From Probabilistic to Predictive CLV](https://towardsdatascience.com/from-probabilistic-to-predictive-methods-for-mastering-customer-lifetime-value-72f090ebcde2/));
([Fader/Hardie 2010 — perils of ignoring heterogeneity](http://brucehardie.com/papers/022/fader_hardie_mksc_10.pdf));
([Improvado CLV guide 2026](https://improvado.io/blog/clv-guide)).

### 11. Cohort-based CLV
Group customers by acquisition period and track value accumulation per cohort.
Reveals retention dynamics and acquisition-quality drift that a single base-wide
average masks; pairs naturally with sBG fitted to multicohort data. (Keep the
*retention-curve fitting itself* in **da-34**; here it is a CLV input/segmentation
lens.) ([Macabacus — Forecasting CLV](https://macabacus.com/blog/forecasting-clv-modeling-customer-lifetime-value-changes));
([Morgan-Dibie — Forecasting CLV with cohort analysis](https://medium.com/@KingHenryMorgansDiary/forecasting-customer-lifetime-value-with-cohort-analysis-d84e9ab1cf8f));
([Fader/Hardie 2010](http://brucehardie.com/papers/022/fader_hardie_mksc_10.pdf)).

### 12. CAC:LTV ratio (unit economics)
LTV:CAC measures payback on acquisition spend. Rules of thumb: **~3:1 is the
healthy target** (B2C SaaS ≈ 2.5:1, B2B SaaS ≈ 4:1); below 2:1 = unsustainable
spend; above ~5:1 = likely *under*-investing in growth. CAC payback period:
healthy 6–12 months, elite < 3 months. Use a **margin-based, discounted**
predictive LTV — a gross-revenue LTV inflates the ratio.
([Phoenix Strategy Group — LTV:CAC benchmarks](https://www.phoenixstrategy.group/blog/ltvcac-ratio-saas-benchmarks-and-insights));
([SaaS Hero — B2B LTV:CAC benchmarks 2026](https://www.saashero.net/strategy/b2b-saas-ltv-cac-benchmarks/));
([Stripe — CAC in SaaS](https://stripe.com/resources/more/cac-in-saas)).

## Tools / Frameworks

| Tool | Lang | Notes |
| --- | --- | --- |
| **PyMC-Marketing** | Python | **Active successor**; Bayesian (MCMC), full uncertainty. Models: `BetaGeoModel`, `ParetoNBDModel`, `ModifiedBetaGeoModel`, `ShiftedBetaGeoModel`, `BetaGeoBetaBinomModel`, `GammaGammaModel`; `rfm_summary()` preprocessor. ([docs](https://www.pymc-marketing.io/en/stable/notebooks/clv/clv_quickstart.html)) |
| **lifetimes** | Python | Cam Davidson-Pilon; **archived / maintenance-only**, MLE fitters (`BetaGeoFitter`, `ParetoNBDFitter`, `GammaGammaFitter`, `BetaGeoBetaBinomFitter`). Still widely used; migrate new work to PyMC-Marketing. ([repo](https://github.com/CamDavidsonPilon/lifetimes); [successor issue #414](https://github.com/CamDavidsonPilon/lifetimes/issues/414)) |
| **CLVTools** | R | Bachmann et al.; S4 API, covariates, `pnbd`/`bgnbd`/`ggomnbd`/`gg`, built-in DERT/DECT. ([clvtools.com](https://www.clvtools.com/)) |
| **BTYD / BTYDplus** | R | Classic R packages; closed-form Pareto/NBD, BG/NBD, BG/BB. `btyd` (Python) is a transitional fork. ([btyd](https://github.com/ColtAllen/btyd)) |

## Methodology (non-contractual continuous: the common case)

1. **Confirm the quadrant.** Non-contractual + continuous → proceed. Contractual
   → sBG/survival instead. Discrete opportunities → BG/BB.
2. **Build RFM summary** from the transaction log (`rfm_summary()` / lifetimes
   `summary_data_from_transaction_data`). Watch the recency definition.
3. **Fit a frequency/dropout model** (BG/NBD default; MBG/NBD if many one-and-done
   customers; Pareto/NBD as benchmark).
4. **Check the model**: plot frequency/recency holdout calibration, the
   tracking plot (cumulative repeat transactions), and `P(alive)` distribution.
5. **Fit Gamma-Gamma on repeat purchasers**; first verify low corr(frequency,
   monetary).
6. **Compute discounted CLV** = E[spend] × DERT over a finite horizon, with a
   continuously-compounded discount rate.
7. **Validate** on a holdout window; compare predicted vs actual transactions
   and revenue by RFM decile.
8. **Segment / act**: rank by predicted CLV and `P(alive)`; feed CAC:LTV.

## Practical Patterns

- **Default to BG/NBD + Gamma-Gamma** for e-commerce/retail; it's the
  best-documented, fastest-to-fit baseline.
- **Use PyMC-Marketing for uncertainty** — point estimates of CLV are
  dangerous for budget decisions; posterior intervals expose how little signal
  thin-history customers carry.
- **Hold out a time window** (calibration/holdout split) and judge on holdout
  fit, never in-sample.
- **Keep horizon finite & discounted** for business reporting (e.g., 12–36 mo).
- **Tie CLV to CAC** at the cohort/channel level so acquisition spend is
  judged on predicted, margin-based, discounted value.

## Anti-Patterns

- **Wrong quadrant.** Applying Pareto/NBD to a subscription business (or sBG to
  e-commerce). Check whether churn is *observed* first.
- **Marketing-RFM recency.** Feeding "days since last purchase" where the model
  wants "first-to-last span." Silent, severe bias.
- **Gamma-Gamma without the independence check.** If spend correlates with
  frequency, the monetary estimates are invalid.
- **Naive ARPU ÷ churn as ground truth.** Assumes a single retention rate;
  biased low under heterogeneity (Fader/Hardie 2010).
- **Un-discounted / infinite-horizon CLV.** Inflates value and the LTV:CAC ratio.
- **Fitting Gamma-Gamma on all customers** instead of repeat purchasers only.
- **Trusting `lifetimes` for new long-lived projects** — it's archived; prefer
  PyMC-Marketing or CLVTools.

## Troubleshooting

- **`P(alive)` looks implausibly high for everyone** → likely BG/NBD with many
  one-and-done customers; switch to MBG/NBD.
- **Optimizer fails / NaN log-likelihood (Pareto/NBD)** → numerical instability
  in the hypergeometric terms; use the log-sum-exp-patched BTYD implementation or
  BG/NBD instead.
- **Gamma-Gamma returns absurd spend** → check you filtered to frequency > 0 and
  that monetary_value is the *average repeat* value, not total.
- **Holdout transactions systematically over-predicted** → calibration window
  caught a promo spike; re-split or model seasonality outside BTYD.
- **CLV explodes** → infinite horizon or zero discount rate; cap the horizon and
  set a continuously-compounded rate.

## References

1. Schmittlein, Morrison & Colombo, "Counting Your Customers…", *Management Science* 33(1):1–24 (1987) — https://pubsonline.informs.org/doi/10.1287/mnsc.33.1.1
2. Fader, Hardie & Lee, "'Counting Your Customers' the Easy Way (BG/NBD)", *Marketing Science* 24(2):275–284 (2005) — http://brucehardie.com/papers/018/fader_et_al_mksc_05.pdf
3. Fader, Hardie & Lee, "RFM and CLV: Iso-Value Curves" (Gamma-Gamma + DERT), *JMR* 42(4):415–430 (2005) — https://www.brucehardie.com/papers/rfm_clv_2005-02-16.pdf
4. Fader & Hardie, Gamma-Gamma derivation note 025 — https://www.brucehardie.com/notes/025/gamma_gamma.pdf
5. Fader & Hardie, "How to Project Customer Retention" (sBG), *J. Interactive Marketing* 21(1):76–90 (2007) — https://www.researchgate.net/publication/264466889_How_to_Project_Customer_Retention
6. Fader & Hardie, "Customer-Base Valuation in a Contractual Setting", *Marketing Science* 29(1):85–93 (2010) — http://brucehardie.com/papers/022/fader_hardie_mksc_10.pdf
7. Fader, Hardie & Shang, "Customer-Base Analysis in a Discrete-Time Noncontractual Setting" (BG/BB), *Marketing Science* 29(6):1086–1108 (2010) — http://www.brucehardie.com/papers/020/fader_et_al_mksc_10.pdf
8. Batislam, Denizel & Filiztekin, "Empirical validation and comparison…" (MBG/NBD), *IJRM* 24(3) (2007) — https://www.sciencedirect.com/science/article/abs/pii/S0167811607000171
9. Fader & Hardie, "What's Wrong With This CLV Formula?" note 033 — http://www.brucehardie.com/notes/033/what_is_wrong_with_this_CLV_formula.pdf
10. Fader & Hardie, BG/NBD step-by-step derivation, note 039 (2019) — http://www.brucehardie.com/notes/039/bgnbd_derivation__2019-11-06.pdf
11. PyMC-Marketing CLV docs (quickstart, BG/NBD, Pareto/NBD, sBG; v0.15.x, 2024–2025) — https://www.pymc-marketing.io/en/stable/notebooks/clv/clv_quickstart.html
12. CLVTools (R) reference — https://www.clvtools.com/
13. lifetimes (Python, archived) + successor issue #414 — https://github.com/CamDavidsonPilon/lifetimes
14. Retina.ai, "History of Buy-Till-You-Die Models" (2023) — https://retina.ai/academy/lesson/history-of-buy-til-you-die-btyd-models/
15. Phoenix Strategy Group, "LTV:CAC Ratio SaaS Benchmarks" — https://www.phoenixstrategy.group/blog/ltvcac-ratio-saas-benchmarks-and-insights
16. SaaS Hero, "Best LTV:CAC Benchmarks for B2B SaaS in 2026" — https://www.saashero.net/strategy/b2b-saas-ltv-cac-benchmarks/
17. "From Probabilistic to Predictive: Mastering CLV", Towards Data Science — https://towardsdatascience.com/from-probabilistic-to-predictive-methods-for-mastering-customer-lifetime-value-72f090ebcde2/
