<!-- hub-reference-banner -->
> **Reference file — part of the `da-applied-and-communication` hub.** Formerly the standalone `da-40-pricing-and-revenue-analytics` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-40-pricing-and-revenue-analytics
title: Pricing and Revenue Analytics
version: "1.0.0"
updated: "2026-05-30"
category: data-analysis
origin: local
related_skills:
  - da-15-forecasting
  - da-22-marketing-mix-modeling
  - da-33-prescriptive-analytics
  - da-12-ab-testing-causal-inference
  - da-23-customer-lifetime-value
  - da-34-cohort-retention-analytics
  - da-6-statistical-modeling
description: >
  Pricing and revenue analytics for data analysts — estimating how price drives
  demand, revenue, and margin, and turning that into pricing decisions. Covers
  price elasticity of demand (own-price and cross-price, log-log / constant-
  elasticity vs linear/semi-log demand, arc vs point elasticity), estimating
  elasticity from observational data and the endogeneity-of-price problem
  (instrumental variables — Hausman, BLP, cost-shifter, and Nielsen-style
  instruments — plus panel FE and Gaussian-copula control-function corrections);
  demand-curve estimation and discrete-choice demand (multinomial logit, IIA and
  the red-bus/blue-bus problem, nested logit/GEV, mixed/random-coefficients logit,
  BLP for differentiated products); willingness-to-pay measurement (Van Westendorp
  PSM, Gabor-Granger, choice-based conjoint CBC, MaxDiff/best-worst scaling);
  price optimization and revenue management (revenue- vs profit-maximizing price,
  Lerner/inverse-elasticity rule, yield/revenue management, dynamic and surge
  pricing, markdown optimization, price laddering, constrained price optimization);
  promotion and discount analytics (baseline estimation, promo lift and
  incrementality, cannibalization, halo, pantry-loading/forward-buy, trade-promotion
  ROI); subscription/SaaS pricing analytics (packaging and tiering, price-volume-mix,
  expansion/contraction, NRR/GRR levers, usage-based pricing, WTP segmentation);
  price-volume-mix bridge / margin-variance walks; and price A/B testing including
  its fairness, ethics, and legal pitfalls plus geo-based price tests. Tooling:
  statsmodels / linearmodels (IV2SLS, PanelOLS), PyMC and PyMC-Marketing for
  Bayesian elasticity and price optimization, pyblp for BLP, biogeme / pylogit /
  xlogit for choice models, Sawtooth/Conjointly/Qualtrics for conjoint & MaxDiff,
  SciPy/CVXPY/OR-Tools for constrained price optimization.
  TRIGGER: estimate price elasticity (own- or cross-price) or a demand curve; fix
  price endogeneity / find an instrument for price; fit a logit / nested logit /
  mixed logit / BLP demand model; measure willingness to pay or run a Van Westendorp,
  Gabor-Granger, conjoint, or MaxDiff study; find a revenue- or profit-maximizing
  price; do yield/revenue management, dynamic pricing, or markdown optimization;
  measure promo lift, baseline, cannibalization, halo, or pantry-loading; analyze
  SaaS/subscription pricing, packaging, price-volume-mix, or expansion/contraction;
  build a price-volume-mix or margin bridge; run or critique a price A/B or geo
  price test. SKIP: forecasting demand/revenue over time as a time series with no
  price lever (da-15-forecasting); marketing-mix / media-spend ROI and ad
  attribution (da-22-marketing-mix-modeling); generic mathematical optimization,
  solver choice, or decision science with no pricing framing (da-33-prescriptive-
  analytics); A/B test statistics and causal-inference theory in general
  (da-12-ab-testing-causal-inference — this skill covers only the price-specific
  application and its pitfalls); probabilistic CLV / BG-NBD (da-23); cohort/retention
  and NRR-curve math (da-34).
triggers:
  - pricing analytics
  - revenue analytics
  - price elasticity
  - elasticity of demand
  - own-price elasticity
  - cross-price elasticity
  - log-log demand
  - constant elasticity
  - demand curve estimation
  - price endogeneity
  - instrument for price
  - instrumental variables price
  - hausman instrument
  - cost-shifter instrument
  - discrete choice demand
  - multinomial logit demand
  - nested logit
  - mixed logit
  - random coefficients logit
  - BLP
  - pyblp
  - biogeme
  - willingness to pay
  - van westendorp
  - price sensitivity meter
  - gabor-granger
  - conjoint analysis
  - choice-based conjoint
  - maxdiff
  - best-worst scaling
  - price optimization
  - profit-maximizing price
  - revenue-maximizing price
  - lerner rule
  - inverse elasticity rule
  - yield management
  - revenue management
  - dynamic pricing
  - surge pricing
  - markdown optimization
  - price laddering
  - promo lift
  - promotion incrementality
  - baseline estimation
  - cannibalization
  - halo effect
  - pantry loading
  - forward buying
  - trade promotion roi
  - saas pricing
  - subscription pricing
  - usage-based pricing
  - packaging and tiering
  - price-volume-mix
  - margin bridge
  - margin variance
  - price ab test
  - geo price test
  - price experiment
---

# Pricing and Revenue Analytics

## Overview

Pricing and revenue analytics is the analytical discipline of estimating how **price**
moves **demand, revenue, and margin**, and converting those estimates into pricing
decisions. It sits at the intersection of microeconomics (demand theory, elasticity),
econometrics (causal estimation under price endogeneity), survey/choice methodology
(stated- and revealed-preference WTP), and operations research (constrained price
optimization). The defining question is always *"what happens to quantity, revenue,
and margin if we change the price?"* — distinct from forecasting a time series, from
attributing marketing spend, or from generic optimization.

**Scope boundary (what this skill is NOT):**
- **da-15-forecasting** — predicting demand/revenue forward in time with no price
  decision lever. This skill uses demand *models* where price is the causal driver.
- **da-22-marketing-mix-modeling** — decomposing sales into media/promo/base with
  adstock and saturation for *budget* allocation. Promotion analytics here is the
  *price-discount* slice (lift, cannibalization, pantry-loading), not media ROI.
- **da-33-prescriptive-analytics** — the LP/MILP/convex solver machinery and decision
  science generally. This skill *applies* that machinery to the pricing objective
  (profit/revenue subject to elasticity and business constraints) and frames it.
- **da-12-ab-testing-causal-inference** — experiment design and causal theory in
  general. Here we cover only the *price-specific* experiment and its pitfalls.

## Core Concepts

### 1. Price elasticity of demand
- **Own-price elasticity** ε = (%ΔQ)/(%ΔP); demand is **elastic** (|ε|>1), **unit-elastic**
  (|ε|=1), or **inelastic** (|ε|<1). Revenue is maximized at |ε|=1; **profit**-maximizing
  price sits where |ε|>1 (you never price in the inelastic region with positive marginal cost).
- **Cross-price elasticity** ε_AB = (%ΔQ_A)/(%ΔP_B): positive ⇒ substitutes, negative ⇒
  complements. The full **own/cross elasticity matrix** drives portfolio and cannibalization analysis.
- **Functional forms.** *Linear* demand Q = a − bP gives elasticity that varies along the
  curve. *Log-log / constant-elasticity* ln Q = α + β ln P makes β the (constant) elasticity
  directly — the workhorse spec. *Semi-log* ln Q = α + βP gives a constant *semi-elasticity*.
  Add `ln P_competitor`, promo flags, seasonality, and `ln income` as controls.
- **Arc vs point elasticity:** arc (midpoint) elasticity for two discrete price points;
  point elasticity = derivative-based, read off a fitted curve.
- **The Lerner / inverse-elasticity rule:** at the optimum, (P − MC)/P = 1/|ε|. Markup is
  the inverse of elasticity — the bridge from a fitted elasticity to an optimal price.

### 2. Endogeneity of price and identification
- **The core problem:** price is **not exogenous**. Firms set high prices when they expect
  high demand (demand shocks correlate with price), and OLS on ln Q ~ ln P is biased —
  typically *toward zero / upward-sloping*, understating true elasticity. This is the central
  technical pitfall in observational pricing work.
- **Instrumental variables (IV / 2SLS):** find a variable that shifts price but is
  uncorrelated with the demand shock. Standard instruments:
  - **Cost shifters** — input costs, wages, exchange rates, fuel, freight.
  - **Hausman instruments** — prices of the same product in *other markets* (common cost
    shock, independent local demand shock); criticized when demand shocks are correlated
    across markets (e.g. national advertising).
  - **BLP instruments** — characteristics of *rival* products in the same market.
  - Wholesale/list price as instrument for retail price; promotion calendars set in advance.
- **Diagnostics:** first-stage F (weak-instrument rule of thumb F>10; use effective F /
  Montiel-Olea–Pflueger for robustness), over-identification (Sargan/Hansen J),
  endogeneity test (Durbin-Wu-Hausman). Weak instruments are worse than OLS.
- **Panel fixed effects** (store × week, product, time FE via `linearmodels.PanelOLS`)
  absorb confounders and are often combined with IV.
- **Gaussian-copula control function** (Park & Gupta) — corrects price endogeneity *without*
  an external instrument by exploiting non-normality of the endogenous regressor; convenient
  but assumes non-normal price and normal errors, and is fragile in small samples.

### 3. Demand-curve and discrete-choice demand estimation
- **Aggregate demand curves** — fit Q(P) (linear, log-log, exponential/decay, logistic) to
  observed price–quantity points, then read elasticity and optimal price off the curve.
- **Discrete-choice (random utility) demand** — model the probability a consumer picks a
  product as a function of its attributes and price:
  - **Multinomial logit (MNL):** closed-form shares; suffers **IIA** (independence of
    irrelevant alternatives) → the **red-bus/blue-bus** problem and unrealistic substitution.
  - **Nested logit (GEV):** groups alternatives into nests, relaxing IIA across nests.
  - **Mixed / random-coefficients logit:** random taste coefficients → flexible, realistic
    substitution; no closed form, simulated likelihood.
  - **BLP (Berry-Levsohn-Pakes):** random-coefficients logit for **differentiated products**
    using *aggregate market-share* data, with a demand inversion and **GMM** using BLP/cost
    instruments to handle the endogenous price inside utility. The standard for IO-style
    market demand and merger/price-change simulation.
- WTP and elasticity fall out of the estimated utility (price coefficient → marginal utility
  of income → WTP for attributes; simulate share changes for elasticities).

### 4. Willingness-to-pay (WTP) measurement
- **Van Westendorp Price Sensitivity Meter (PSM):** 4 questions (too cheap / cheap-bargain /
  expensive / too expensive); intersections give Point of Marginal Cheapness (PMC), Point of
  Marginal Expensiveness (PME), Optimal Price Point (OPP), Indifference Price Point (IPP),
  and the acceptable **range of acceptable prices**. Best **early**, for new-to-world products;
  directional only — no demand/volume.
- **Gabor-Granger:** show each respondent a sequence of specific prices, record purchase
  intent at each → builds a **demand curve** and a revenue-maximizing point. Needs a known
  price range (often follows Van Westendorp). Prone to demand artifacts / anchoring.
- **Conjoint analysis (choice-based, CBC):** respondents choose among product profiles where
  price is one attribute among many; estimate part-worth utilities (HB/logit) → derive WTP,
  share-of-preference simulators, and price elasticities. The gold standard for trade-offs
  and feature-vs-price decisions.
- **MaxDiff (best-worst scaling):** forces respondents to pick best/worst from sets → a stable
  ranked importance of features (no scale-use bias). Used to *prioritize* features feeding a
  conjoint/Gabor-Granger, not to set price directly.
- **Stated vs revealed preference:** surveys (above) are stated-preference and overstate WTP;
  transaction/experiment data is revealed-preference. Triangulate.

### 5. Price optimization & revenue management
- **Objective:** revenue-maximizing price (|ε|=1) vs **profit-maximizing** price (Lerner rule,
  requires marginal cost). Build a profit function π(P) = (P − MC)·Q(P) from the fitted demand
  curve and maximize.
- **Constrained price optimization:** maximize profit/revenue subject to constraints — price
  bounds, margin floors, price-ladder/gap rules across a line, MAP (minimum advertised price),
  cross-elasticity/cannibalization terms, inventory. Solve with `scipy.optimize` (nonlinear),
  **CVXPY** (convex formulations), or OR-Tools/Gurobi for MILP price-point selection. (Solver
  mechanics → da-33; here it's framed as the pricing objective.)
- **Yield vs revenue management:** *yield management* = price/allocate a **fixed, perishable**
  capacity (airline seats, hotel rooms) by segment and time; *revenue management* is the
  broader discipline (also assortment, overbooking, distribution). EMSR heuristics, booking
  limits, protection levels.
- **Dynamic / surge pricing:** prices adjust to real-time demand/inventory; increasingly
  ML-driven (demand prediction at each candidate price). Watch **perceived-fairness** backlash.
- **Markdown optimization:** for seasonal/perishable goods, choose the markdown depth and
  timing that maximizes sell-through revenue before end-of-life. UPPMO = unified pricing,
  promotion & markdown optimization across the lifecycle.
- **Price laddering / line pricing:** coherent price steps across good-better-best tiers and
  pack sizes; preserve sensible per-unit ladders to avoid arbitrage and trading-down.

### 6. Promotion & discount analytics
- **Baseline estimation:** the counterfactual non-promoted sales level. Must strip promo weeks,
  stock-outs, seasonality, and trend; a contaminated baseline misstates lift. Methods: moving-
  average/regression baselines, structural time-series, or causal-impact-style counterfactuals.
- **Lift & incrementality:** `Lift = Actual − Baseline`. Decompose total lift into **true
  incremental** volume, **pantry-loading/forward-buy** (pulled-forward demand → post-promo dip),
  **cannibalization** (own promoted SKU steals from sibling SKUs — cross-elasticities), and
  **halo** (lifts adjacent non-promoted items — market-basket/affinity).
- **Trade-promotion ROI / TPO:** net incremental margin vs promo cost; many promos are
  ROI-negative once cannibalization and forward-buy are netted out.
- This is the **price-discount** slice of promotion. Media/advertising ROI and adstock →
  da-22-marketing-mix-modeling.

### 7. Subscription / SaaS pricing analytics
- **Packaging & tiering:** good-better-best, feature gating, seat vs usage vs hybrid pricing;
  2025 trend toward **usage-based / consumption** pricing (correlates with higher NRR and lower
  churn).
- **Price-volume-mix for ARR:** decompose ARR/MRR growth into new, expansion, contraction, and
  churn; **NRR = expansion − contraction − gross churn** relative to starting base. Expansion is
  the dominant lever at scale (often >50% of growth at NRR ≥100%).
- **WTP segmentation:** estimate WTP by segment (conjoint, surveys, behavioral/usage signals)
  and align tiers/fences (feature, usage, identity) so each segment self-selects — price
  discrimination via versioning. (NRR/GRR curve math and retention cohorts → da-34.)

### 8. Price-volume-mix (PVM) bridge & margin analytics
- **PVM bridge / sales bridge:** decompose the change in revenue or gross margin between two
  periods into **price**, **volume**, and **mix** effects (plus FX and cost for margin), as a
  signed **waterfall** from prior-period to current-period margin where every dollar is accounted.
- Standard decomposition: Price effect = ΔP × Q (at a reference); Volume effect = ΔQ × P;
  Mix effect = shift toward higher/lower-margin products at constant total volume. Margin walk
  adds the COGS side (margin volume effect = sales volume effect − COGS volume effect).
- Used in FP&A and commercial reviews to explain *why* margin moved and assign accountability.

### 9. Price A/B testing — and its pitfalls
- **Why it's hard/fraught:** charging different customers different prices for the same product
  raises **fairness, trust, and legal** issues; visible price tests damage trust and, in the EU,
  create exposure (consumer-protection, hidden-test, and location-manipulation rules). Geo or
  segment splits that correlate with **protected characteristics** risk discriminatory pricing.
- **Design pitfalls:** price is a high-variance, low-frequency conversion outcome → low power,
  long runtimes; novelty/anchoring effects; contamination across billing/sales/finance systems;
  honoring the lower price for everyone after the test.
- **Geo-based tests / matched markets:** randomize or match at the **market** level (designated
  market areas) instead of user level to avoid within-customer price discrimination and arbitrage;
  use synthetic-control / matched-market analysis for the readout.
- **Switchback designs:** alternate price A/B over time windows for the whole market when
  user-level randomization is unethical or infeasible (marketplaces).

## Tools & Frameworks (2025-2026)
- **Elasticity / regression:** `statsmodels` (OLS, semi/log-log), `linearmodels` (`IV2SLS`,
  `PanelOLS` for FE + IV) — the standard Python stack for IV elasticity.
- **Bayesian elasticity & price optimization:** `PyMC` (priors on elasticity, full posterior of
  optimal price/profit), `PyMC-Marketing` (also used for promo/MMM-adjacent work).
- **Differentiated-products / BLP:** **`pyblp`** (Conlon & Gortmaker; v1.x, micro-moments
  framework, nested/mixed logit tutorials on the Nevo cereal data).
- **Choice modeling:** **`biogeme`** (MNL/nested/mixed logit, hybrid choice), `pylogit`,
  `xlogit` (GPU-accelerated mixed logit), `apollo` (R).
- **Conjoint / MaxDiff / WTP surveys:** Sawtooth Software (Lighthouse/CBC), Conjointly,
  Qualtrics; HB estimation for part-worths.
- **Constrained price optimization:** `scipy.optimize`, **`CVXPY`**, OR-Tools / Gurobi (see da-33).
- **Revenue management / markdown:** commercial suites (o9, Blue Yonder, Revionics, Vendavo,
  Pricefx, PROS); UPPMO patterns.
- **Causal/promo readout:** `CausalImpact`/structural time series for baseline counterfactuals.

## Methodology (end-to-end pricing study)
1. **Frame the decision** — what price lever, what objective (revenue vs profit vs share), what
   constraints (margin, ladder, MAP, fairness/legal).
2. **Choose data regime** — observational transactions, panel/scanner, survey/stated-preference,
   or experiment. Decide revealed vs stated preference up front.
3. **Identify causal price effect** — never trust raw OLS; use IV, panel FE, copula CF, or a
   designed experiment. Validate instrument strength and exogeneity.
4. **Estimate demand** — pick functional form / choice model matching the data and substitution
   structure; report elasticity matrix with uncertainty.
5. **Optimize** — build profit/revenue function, apply Lerner rule or constrained optimizer;
   simulate scenarios and sensitivity to elasticity uncertainty.
6. **Account for promo/portfolio effects** — net out cannibalization, halo, pantry-loading.
7. **Validate** — out-of-sample, holdout markets, or a controlled price/geo test before rollout.
8. **Communicate** — PVM/margin bridge to explain expected vs realized impact.

## Practical Patterns
- **Always log-log first** for a quick, interpretable elasticity, then check robustness with a
  flexible form. A coefficient that comes out *positive* is the classic endogeneity tell.
- **Pair IV with panel FE** (store×week FE + cost-shifter IV) — FE absorbs persistent confounders,
  IV handles the simultaneity that FE can't.
- **Sequence WTP methods:** MaxDiff (what matters) → Van Westendorp (acceptable range) →
  Gabor-Granger or CBC (price/demand curve and optimum).
- **Use BLP only when you have aggregate share data + differentiated products + endogenous price
  and need realistic substitution** (merger/price-change simulation); otherwise a mixed logit on
  individual choice data is simpler.
- **Decompose every promo** into incremental / pantry / cannibalization / halo before claiming ROI.
- **Test price at the market (geo) level**, not the user level, to dodge fairness/legal landmines.
- **Express optimal price as a posterior/range**, not a point — elasticity uncertainty dominates.

## Anti-Patterns
- **Running OLS of ln Q on ln P and reporting the coefficient as "the elasticity"** — ignores
  price endogeneity; the single most common error in pricing analytics.
- **Using a weak instrument** (first-stage F < 10) — biased *worse* than OLS; check, don't assume.
- **Trusting stated-preference WTP as absolute** — surveys overstate; calibrate against behavior.
- **Plain MNL where substitution matters** — IIA gives the red-bus/blue-bus absurdity; use nested
  or mixed logit.
- **Claiming promo lift = total uplift** without netting pantry-loading and cannibalization —
  inflates ROI and hides portfolio cannibalization.
- **User-level price A/B tests** — fairness, trust, and legal exposure; honoring-the-low-price tax.
- **Optimizing price in the inelastic region with positive MC** — always raise price there.
- **Mixing up yield vs revenue management**, or treating dynamic pricing as pure profit max while
  ignoring perceived-fairness backlash.

## Troubleshooting
- **Positive or near-zero elasticity coefficient** → endogeneity / reverse causality; instrument
  price or use a designed test; check for stockout-driven and promo-contaminated weeks.
- **Wrong-signed cross-elasticities** → omitted seasonality/promo confounders or collinear prices.
- **Weak first stage** → instrument too weak; find a stronger cost shifter, pool markets, or run
  an experiment instead.
- **Choice model gives implausible substitution** → relax IIA (nested/mixed logit); add random
  coefficients on price.
- **Van Westendorp range too wide / no clear OPP** → respondents don't understand the product;
  add anchoring/context or switch to CBC.
- **Promo "worked" but margin fell** → forward-buy/pantry-loading and cannibalization eating the
  lift; rebuild baseline and decompose.
- **Price test inconclusive** → underpowered (price is low-frequency, high-variance); lengthen,
  use geo/matched-market design, or model elasticity from history instead.

## References
- IV / endogeneity for elasticity: arXiv 2306.12863 *Price elasticity of electricity demand:
  using instrumental variable regressions to address endogeneity and autocorrelation*
  (https://arxiv.org/abs/2306.12863); Springer JAMS 2025 *Utilizing managerial beliefs for set
  identification of price elasticities of demand* (https://link.springer.com/article/10.1007/s11747-025-01090-9);
  UC Riverside *Estimating the Price Elasticity of Gasoline Demand*
  (https://economics.ucr.edu/repec/ucr/wpaper/202021R.pdf).
- BLP / discrete choice: PyBLP docs & v1.2 (https://pyblp.readthedocs.io/) and repo
  (https://github.com/jeffgortmaker/pyblp); Conlon & Gortmaker (2025) micro-moments; arXiv
  2501.02381 *Estimating Discrete Choice Demand Models with Sparse Market-Product Shocks*
  (https://arxiv.org/pdf/2501.02381); arXiv 2602.05137 *Nested Pseudo-GMM Estimation of Demand
  for Differentiated Products* (https://arxiv.org/pdf/2602.05137).
- WTP methods: Conjointly Gabor-Granger (https://conjointly.com/products/gabor-granger/) and
  Gabor-Granger vs Van Westendorp (https://conjointly.com/blog/gabor-granger-or-van-westendorp/);
  Marketbridge survey pricing methodologies (https://marketbridge.com/article/survey-pricing-methodologies/);
  SurveyKing Van Westendorp (https://www.surveyking.com/help/van-westendorp-analysis).
- Revenue management / markdown / dynamic pricing: Stripe *What is yield management*
  (https://stripe.com/resources/more/yield-management); o9 Price Planning & Optimization
  (https://o9solutions.com/solutions/pricing-yield-markdown-management/); Retalon UPPMO retail
  pricing strategy 2025 (https://retalon.com/blog/retail-pricing-strategy); GMInsights dynamic
  pricing & yield management market (https://www.gminsights.com/industry-analysis/dynamic-pricing-and-yield-management-market).
- Promotion analytics: Tredence uplift & halo (https://www.tredence.com/blog/decoding-the-metrics-a-deep-dive-into-calculating-promotion-effectiveness);
  Crosscap retail promotion lift (https://www.crosscap.com/guide-to-analyzing-the-overall-lift-of-a-retail-promotion/);
  SoftServe trade-promotion analytics formulas (https://softservebs.com/en/resources/trade-promotion-analysis/).
- SaaS pricing & NRR: Baremetrics value-based pricing metrics (https://baremetrics.com/blog/key-metrics-value-based-pricing-saas);
  Monetizely SaaS pricing benchmarks 2025 (https://www.getmonetizely.com/articles/saas-pricing-benchmarks-2025-how-do-your-monetization-metrics-stack-up);
  ProductQuant NRR benchmarks (https://productquant.dev/blog/nrr-benchmarks-saas/).
- PVM / margin bridge: Vendavo practical guide to PVM (https://www.vendavo.com/practical-guide-to-pvm-analysis/);
  Business Intelligist PVM for gross-margin variance (https://businessintelligist.com/2020/04/26/price-volume-mix-pvm-for-gross-margin-variance-analysis/);
  Under Controlling sales bridge (https://undercontrolling.com/sales-bridge-volume-price-mix-analysis/).
- Price A/B testing & ethics: Monetizely ethics of SaaS A/B pricing tests
  (https://www.getmonetizely.com/articles/the-ethics-of-saas-ab-pricing-tests-balancing-business-growth-and-customer-trust);
  Statsig A/B testing for pricing (https://www.statsig.com/perspectives/ab-testing-pricing-tips);
  Orb pricing experiments (https://www.withorb.com/blog/pricing-experiments).
- Python implementation: PyMC Bayesian price optimization
  (https://towardsdatascience.com/bayesian-price-optimization-with-pymc3-d1264beb38ee/);
  TDS elasticity with statsmodels (https://medium.com/data-science/calculating-price-elasticity-of-demand-statistical-modeling-with-python-6adb2fa7824d);
  ChenDataBytes linear & non-linear price elasticity (https://medium.com/@chenycy/unlock-price-optimization-potential-with-python-modelling-linear-and-non-linear-price-elasticity-563773e5ba53).
