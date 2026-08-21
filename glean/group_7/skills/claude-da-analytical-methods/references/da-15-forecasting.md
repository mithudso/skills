<!-- hub-reference-banner -->
> **Reference file — part of the `da-analytical-methods` hub.** Formerly the standalone `da-15-forecasting` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-15-forecasting
title: Forecasting
version: "1.0.0"
updated: "2026-05-30"
category: data-analysis
origin: local
description: >
  Forecasting — going deeper than the time-series treatment in
  da-6-statistical-modeling. Covers ARIMA family in depth (SARIMA, SARIMAX,
  ARIMAX, auto-ARIMA, Box-Jenkins workflow), exponential smoothing
  (Holt-Winters, ETS state-space form), Prophet and NeuralProphet, modern ML
  forecasting (XGBoost/LightGBM, foundation models like TimeGPT, Chronos,
  Lag-Llama, Moirai), hierarchical forecasting and reconciliation,
  intermittent demand methods (Croston, SBA, ADIDA), and prediction intervals
  via bootstrap, conformal prediction, and Monte Carlo simulation.
  TRIGGER: building a time-series forecast (sales, demand, revenue, capacity,
  traffic, energy load); choosing between ARIMA / Prophet / ML / foundation
  models; needing prediction intervals or scenario simulation; doing
  hierarchical forecasting that must reconcile across product / region /
  channel; intermittent / sparse demand; backtesting a forecasting model.
  SKIP: cross-sectional regression or classification (da-6-statistical-modeling
  or da-7-machine-learning); causal inference from a time series
  (da-12-ab-testing-causal-inference); change-point or anomaly detection
  (da-16-anomaly-detection); EDA on a time series (da-5).
triggers:
  - forecasting
  - time series forecast
  - ARIMA
  - SARIMA
  - exponential smoothing
  - Holt-Winters
  - Prophet
  - NeuralProphet
  - TimeGPT
  - Chronos
  - Lag-Llama
  - Moirai
  - foundation model forecasting
  - hierarchical reconciliation
  - intermittent demand
  - Croston's method
  - conformal prediction intervals
  - Monte Carlo simulation
keywords:
  - ARIMA
  - SARIMA
  - SARIMAX
  - ARIMAX
  - auto-ARIMA
  - Box-Jenkins
  - ACF
  - PACF
  - ETS
  - Holt-Winters
  - exponential-smoothing
  - Prophet
  - NeuralProphet
  - TimeGPT
  - Chronos
  - Lag-Llama
  - Moirai
  - hierarchical-forecasting
  - MinT-reconciliation
  - intermittent-demand
  - Croston
  - SBA
  - ADIDA
  - conformal-prediction
  - Monte-Carlo
when_to_use:
  - Building a time-series forecast (sales, demand, revenue, traffic, energy)
  - Choosing between ARIMA, Prophet, ML, and foundation models
  - Needing prediction intervals or scenario simulation
  - Doing hierarchical forecasting that must reconcile across groups
  - Intermittent / sparse demand (most periods are zero)
  - Backtesting a forecasting model with rolling windows
  - Quantifying forecast uncertainty for downstream decisions
when_not_to_use:
  - Cross-sectional regression or classification — use da-6 or da-7
  - Causal inference from a time series — use da-12-ab-testing-causal-inference
  - Change-point or anomaly detection — use da-16-anomaly-detection
  - Exploratory time-series analysis — use da-5-exploratory-data-analysis
related_skills:
  - da-1-3-probability-theory
  - da-1-4-statistical-inference-foundations
  - da-6-statistical-modeling
  - da-7-machine-learning
  - da-12-ab-testing-causal-inference
  - da-16-anomaly-detection
---

# Forecasting

Predicting future values of a time-indexed series. The canonical textbook is Hyndman & Athanasopoulos *Forecasting: Principles and Practice* (3rd ed) — every section in this skill maps to a chapter there.

## When to use this skill

Activate when the user:
- is building a time-series forecast
- needs to choose between ARIMA / Prophet / ML / foundation models
- needs prediction intervals or scenario simulation
- has hierarchical structure (e.g., total → region → store) that must reconcile
- has intermittent / sparse demand
- is backtesting a forecasting model

## When NOT to use this skill

- Cross-sectional regression → `da-6-statistical-modeling`
- Causal inference → `da-12-ab-testing-causal-inference`
- Anomaly detection → `da-16-anomaly-detection`
- Exploratory time-series analysis → `da-5-exploratory-data-analysis`

---

## The ARIMA family

The classical workhorse. Five letters: AR, I, MA, S (seasonal), X (exogenous).

| Variant | What it adds | When to use |
|---|---|---|
| **AR(p)** | Autoregressive — value depends on prior p values | Stationary series with persistence |
| **MA(q)** | Moving average of past q shocks | Series with short-lived shocks |
| **ARMA(p,q)** | Both | Stationary; check ACF/PACF for orders |
| **ARIMA(p,d,q)** | Difference d times to make stationary | Non-stationary level / trend |
| **SARIMA(p,d,q)(P,D,Q)ₘ** | Add seasonal component with period m | Strong seasonality |
| **SARIMAX** | + exogenous regressors X | Want to use external drivers |
| **ARIMAX** | ARIMA + exogenous | Same idea, no seasonality |

### Box-Jenkins workflow

1. **Identify** — plot the series, check stationarity (ADF / KPSS), difference if needed
2. **Order selection** — read ACF and PACF, propose (p, d, q) and (P, D, Q, m)
3. **Estimate** — fit by MLE
4. **Diagnose** — residual ACF should be white noise; Ljung-Box test
5. **Iterate** — if residuals show structure, increase order or revisit step 1
6. **Forecast** — produce point + interval forecasts

### auto-ARIMA

`pmdarima.auto_arima` (Python) / `forecast::auto.arima` (R). Searches the (p, d, q, P, D, Q) grid by AIC/BIC. Good first pass; almost never the final model. Always check residuals.

---

## Exponential smoothing (ETS family)

Each new forecast is a weighted average of past values with weights that decay exponentially. Simple, robust, often beats fancy methods on short horizons.

| Method | Components | When |
|---|---|---|
| **Simple ES** | Level only | No trend, no seasonality |
| **Holt's** | Level + trend | Trend, no seasonality |
| **Holt-Winters additive** | Level + trend + seasonal | Constant seasonal amplitude |
| **Holt-Winters multiplicative** | Level + trend + seasonal | Seasonal amplitude scales with level |
| **ETS(error, trend, seasonal)** | Unified state-space framework | All of the above + automatic selection |

The ETS state-space form (Hyndman) is the modern unification: error type ∈ {A, M}, trend ∈ {N, A, Ad}, seasonal ∈ {N, A, M}. Software (statsforecast, R `fable`) picks the best combination via AICc.

ETS often outperforms ARIMA on M3 / M4 / M5 competition data with similar effort.

---

## Prophet and NeuralProphet

### Prophet (Meta, 2017)

An additive model: `y(t) = g(t) + s(t) + h(t) + ε`
- **g(t)** — trend (piecewise linear or logistic with changepoints)
- **s(t)** — seasonality (Fourier series)
- **h(t)** — holiday effects (user-specified calendar)
- **ε** — noise

Strengths: handles missing data, holidays, multiple seasonalities (daily + weekly + yearly), automatic changepoint detection. Robust to outliers.

Weaknesses: not great for short series, ignores autocorrelation in residuals, often beaten by simpler methods on competitions.

When to reach for Prophet: business time series with strong calendar/holiday effects (retail sales, web traffic) and you want a fast competent baseline without tuning.

### NeuralProphet (2020-2024)

Prophet + AR terms + neural network for nonlinear effects. PyTorch under the hood. More expressive than Prophet but harder to tune. Worth trying if Prophet underfits.

---

## Modern ML forecasting

### Tree boosting (XGBoost / LightGBM / CatBoost)

Treat the forecasting problem as a regression: features are lagged values, calendar variables, rolling aggregates. Fit a gradient-boosted tree per horizon (or one model with horizon as a feature).

Strengths: handles many series, exogenous features, interactions. Won the M5 competition (Walmart sales).

Pitfalls: stationarity matters less but train-test split must be **time-based**, not random. Use expanding-window or rolling-origin backtesting (never k-fold).

### Foundation models (2024-2026)

The 2024-2026 frontier: pre-trained transformer models that zero-shot forecast any series.

| Model | Origin | Notes |
|---|---|---|
| **TimeGPT** (Nixtla, 2023) | Commercial API | First production-grade time-series foundation model |
| **Chronos** (Amazon, 2024) | Open-source on Hugging Face | Tokenizes values into a vocabulary; LLM-style |
| **Lag-Llama** (ServiceNow, 2024) | Open-source | Llama-architecture for time series |
| **Moirai** (Salesforce, 2024) | Open-source | Universal forecasting model |
| **TimesFM** (Google, 2024) | Open-source | Decoder-only foundation model |
| **TabPFN-derived adaptations** | TabPFN (Hollmann et al, ICLR 2023) repurposed for time series via lag-feature framing | Strong on short series with few examples |

**The hype-vs-reality check (2026):** foundation models are not universally better. On the GIFT-Eval benchmark, well-tuned classical methods (ETS, ARIMA, Theta) still beat foundation models on many real datasets, especially short series. The foundation models win when:
- you have no time to tune per series
- you have many heterogeneous series (the model amortizes across them)
- the series has structure the foundation model recognizes from pretraining

For a single critical series, fitting ETS or ARIMA usually still wins. For 10,000 series, a foundation model is the right call.

---

## Hierarchical forecasting

When forecasts must roll up consistently: total = sum(regions) = sum(stores). The naive approach (forecast each level independently) produces inconsistent totals.

### Reconciliation methods

| Method | Approach |
|---|---|
| **Bottom-up** | Forecast the bottom level, sum to higher levels |
| **Top-down** | Forecast the top, allocate down by historical proportions |
| **Middle-out** | Forecast middle, sum up and allocate down |
| **MinT (Minimum Trace)** | Optimal reconciliation; minimizes forecast variance subject to coherence (Wickramasuriya, Athanasopoulos, Hyndman 2019) |

MinT is the modern default. Implemented in `hierarchicalforecast` (Nixtla), `hts` and `fable` (R). Forecast each level independently, then reconcile via MinT.

---

## Intermittent demand

Series where most periods are zero (spare parts, slow movers, rare events). Standard methods produce silly forecasts (often negative, often constant).

| Method | What it does |
|---|---|
| **Croston's method (1972)** | Forecast the inter-demand interval and the demand size separately |
| **SBA** (Syntetos-Boylan Approximation, 2005) | Bias correction on Croston |
| **TSB** (Teunter-Syntetos-Babai, 2011) | Updates probability of demand at each period |
| **ADIDA / IMAPA** | Aggregate to higher frequency, forecast, disaggregate |

The 2-by-2 classification (Syntetos-Boylan-Croston 2005): coefficient of variation² of demand sizes (CV²) × average inter-demand interval (ADI). Recommends a method per quadrant.

---

## Prediction intervals

Point forecasts alone are dangerous. Always quantify uncertainty.

### Sources of uncertainty

1. **Parameter uncertainty** — your estimated coefficients aren't the true ones
2. **Model uncertainty** — your model class isn't the right one
3. **Innovation uncertainty** — even with correct parameters, future shocks happen

### Methods

- **Analytical** — ARIMA / ETS produce closed-form intervals from the assumed error distribution. Often too narrow (assume normal errors).
- **Bootstrap** — resample residuals, simulate forward, repeat. Captures parameter + innovation uncertainty.
- **Monte Carlo simulation** — explicit simulation of the model going forward. Useful when business rules transform the forecast.
- **Conformal prediction** — distribution-free intervals with finite-sample coverage guarantee. The 2020s frontier; works with ANY model (including ML and foundation models). Implementations: `mapie`, `crepes`.

### What to report

- A point forecast alone is misleading
- An 80% prediction interval is the common business default
- The 95% interval is what statisticians want
- A fan chart (cone of intervals at increasing percentiles) is the executive-friendly version

---

## Backtesting

Never random-split a time series. The future is not interchangeable with the past.

- **Train-test split** — train on past, test on future. Minimum viable.
- **Expanding window** — fit on `[1..t]`, forecast `t+h`, evaluate, increment t. The standard.
- **Rolling origin** — fit on a fixed-width window, slide forward. Stronger if the relationship changes over time.

### Evaluation metrics

- **MAE** — mean absolute error; robust to outliers
- **RMSE** — penalizes large errors more
- **MAPE** — % error; undefined when actual is zero; biased toward over-forecasts
- **sMAPE** — symmetric MAPE; better-behaved
- **MASE** (Mean Absolute Scaled Error, Hyndman & Koehler 2006) — error scaled by the in-sample naive forecast error. Comparable across series. **The recommended default.**
- **CRPS** (Continuous Ranked Probability Score) — for probabilistic forecasts (intervals or full distributions)

---

## Anti-patterns

1. **Random k-fold cross-validation on time series** — leaks future into training.
2. **MAPE on series with zeros** — gives nonsense.
3. **Point forecasts without intervals** — executives anchor on the point and underestimate downside.
4. **ARIMA on non-stationary series without differencing** — wildly wrong.
5. **Holt-Winters multiplicative on a series that hits zero** — divides by zero.
6. **Foundation model as drop-in replacement for tuned classical** — often loses on small datasets.
7. **Independent forecasts of hierarchy levels that don't reconcile** — pick a reconciliation method.
8. **Forecasting horizons that ignore the calendar** — Q4 isn't like Q1 in retail; bake seasonality in.
9. **Single-model bake-off without ensemble** — averaging ETS + ARIMA + Prophet often beats any one.

---

## References

1. Hyndman, R. J. & Athanasopoulos, G. (2021). *Forecasting: Principles and Practice* (3rd ed). otexts.com/fpp3 — free online. The book.
2. Box, G. E. P., Jenkins, G. M., Reinsel, G. C., Ljung, G. M. (2015). *Time Series Analysis: Forecasting and Control* (5th ed).
3. Hyndman, R. J. & Koehler, A. B. (2006). "Another look at measures of forecast accuracy." *IJF*. (MASE paper.)
4. Wickramasuriya, S. L., Athanasopoulos, G., & Hyndman, R. J. (2019). "Optimal Forecast Reconciliation for Hierarchical and Grouped Time Series Through Trace Minimization." *JASA*. (MinT paper.)
5. Croston, J. D. (1972). "Forecasting and stock control for intermittent demands." *Op Res Quarterly*.
6. Syntetos, A. A., Boylan, J. E., & Croston, J. D. (2005). "On the categorization of demand patterns." *J Op Res Society*.
7. Nixtla — TimeGPT, statsforecast, hierarchicalforecast — https://nixtlaverse.nixtla.io/
8. Meta Prophet — https://facebook.github.io/prophet/
9. Salinas, D. et al. (2020). "DeepAR: Probabilistic forecasting with autoregressive RNNs." *IJF*.
10. Ansari, A. et al. (2024). "Chronos: Learning the Language of Time Series." arXiv:2403.07815.
11. Das, A. et al. (2024). "A decoder-only foundation model for time-series forecasting." (TimesFM) ICML.
12. Liang et al. (2025). "GIFT-Eval: A Benchmark for General Time Series Forecasting." arXiv. (The honest evaluation.)
13. Vovk, V., Gammerman, A., & Shafer, G. (2005). *Algorithmic Learning in a Random World*. (Conformal prediction.)
14. M-Competition results (M3, M4, M5) — https://www.researchgate.net/publication/345753797 (M5 paper).
