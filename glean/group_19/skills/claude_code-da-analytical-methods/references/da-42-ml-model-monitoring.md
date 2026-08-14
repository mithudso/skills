<!-- hub-reference-banner -->
> **Reference file — part of the `da-analytical-methods` hub.** Authored as a hub reference (no prior standalone skill); consolidates ML model-monitoring material previously scattered across `da-7-machine-learning` (drift/skew/MLflow/W&B sub-bullets) and `da-17-feature-engineering-and-feature-stores` (feature monitoring).
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-42-ml-model-monitoring
title: ML Model Monitoring & Production Model Observability
version: "1.0.0"
updated: "2026-05-31"
category: data-analysis
origin: local
description: >
  Monitoring machine-learning models in production — the consolidated
  discipline. Covers the drift taxonomy (data/covariate, concept, prediction/
  output, label/prior, feature drift) and detection tests (PSI, KL/JS
  divergence, KS, Chi-square, Wasserstein/EMD, L-infinity); sequential
  concept-drift detectors (DDM, EDDM, ADWIN, Page-Hinkley, CUSUM); model
  performance monitoring under delayed or absent ground truth (proxy metrics,
  two-loop monitoring, label lag); performance estimation without labels
  (NannyML CBPE and DLE, M-CBPE); training-serving skew detection;
  segment/slice-based performance monitoring and fairness drift; input
  outlier/adversarial detection (Alibi Detect); alerting, retraining triggers,
  and the monitoring->retraining loop; and the tooling landscape (Evidently,
  Arize, Fiddler, WhyLabs/whylogs, NannyML, Seldon/Alibi Detect, SageMaker
  Model Monitor, Vertex AI Model Monitoring, MLflow). Clarifies how model
  monitoring relates to but differs from data observability and LLM
  observability.
  TRIGGER: a deployed model is silently degrading; detecting data drift,
  concept drift, or prediction drift; choosing a drift test (PSI vs KS vs
  Wasserstein vs JS); estimating model accuracy/precision/recall when labels
  are delayed or unavailable; training-serving skew; slice/subgroup or fairness
  drift; setting alert thresholds and retraining triggers; wiring a
  monitoring->retraining loop; comparing model-monitoring tools (Evidently,
  Arize, Fiddler, WhyLabs, NannyML, Alibi Detect, SageMaker, Vertex).
  SKIP: detecting unusual individual records/events as the end goal
  (da-16-anomaly-detection); the engineering of features and feature stores,
  including offline/online feature monitoring and freshness
  (da-17-feature-engineering-and-feature-stores); data-pipeline / dataset
  health — freshness, volume, schema, lineage of the data itself
  (da-19-data-observability under da-data-engineering-platform); choosing the
  model or its evaluation metrics pre-deployment (da-7-machine-learning);
  forecasting time series (da-15-forecasting).
triggers:
  - model monitoring
  - model observability
  - data drift
  - concept drift
  - prediction drift
  - feature drift
  - train-serve skew
  - training-serving skew
  - performance estimation without labels
  - NannyML CBPE
  - NannyML DLE
  - PSI drift
  - KS test drift
  - Wasserstein drift
  - ADWIN
  - DDM
  - Page-Hinkley
  - retraining trigger
  - fairness drift
  - Evidently
  - Arize
  - Fiddler
  - WhyLabs
  - Alibi Detect
  - SageMaker Model Monitor
  - Vertex Model Monitoring
keywords:
  - data-drift
  - covariate-shift
  - concept-drift
  - prediction-drift
  - label-shift
  - prior-probability-shift
  - feature-drift
  - PSI
  - KL-divergence
  - Jensen-Shannon
  - KS-test
  - Chi-square
  - Wasserstein
  - earth-movers-distance
  - L-infinity
  - DDM
  - EDDM
  - ADWIN
  - Page-Hinkley
  - CUSUM
  - delayed-ground-truth
  - label-lag
  - proxy-metrics
  - CBPE
  - DLE
  - M-CBPE
  - isotonic-calibration
  - train-serve-skew
  - slice-based-monitoring
  - subgroup-drift
  - fairness-drift
  - retraining-trigger
  - monitoring-retraining-loop
  - Evidently
  - Arize
  - Fiddler
  - WhyLabs
  - whylogs
  - NannyML
  - Alibi-Detect
  - Seldon
  - SageMaker-Model-Monitor
  - Vertex-Model-Monitoring
  - MLflow
  - C2ST
  - MMD
when_to_use:
  - A deployed model is silently degrading and you need to know why
  - Detecting data drift, concept drift, or prediction/output drift
  - Choosing a drift test (PSI vs KS vs Wasserstein vs JS divergence)
  - Estimating live performance when ground-truth labels are delayed or absent
  - Detecting training-serving skew
  - Monitoring performance by slice/segment and watching for fairness drift
  - Setting alert thresholds and designing retraining triggers
  - Wiring the monitoring -> retraining loop into an MLOps pipeline
  - Comparing model-monitoring tools and deciding build-vs-buy
when_not_to_use:
  - Flagging unusual individual records/events as the goal — use da-16-anomaly-detection
  - Building features or a feature store / feature monitoring — use da-17-feature-engineering-and-feature-stores
  - Dataset/pipeline health (freshness, volume, schema, lineage) — use da-19-data-observability
  - Picking the model or its offline eval metrics pre-deployment — use da-7-machine-learning
  - Forecasting a time series — use da-15-forecasting
related_skills:
  - da-7-machine-learning
  - da-16-anomaly-detection
  - da-17-feature-engineering-and-feature-stores
  - da-19-data-observability
  - da-12-ab-testing-causal-inference
  - da-31-conformal-prediction-uq
---

# ML Model Monitoring & Production Model Observability

A model that passed every offline test can still rot in production. The world
moves, the inputs move with it, and the model keeps emitting confident
predictions against a reality it no longer matches — usually with **no error
raised**. Model monitoring is the discipline of catching that silent
degradation: tracking the inputs, the predictions, and (eventually) the outcomes
of a live model so you know *when* it has degraded, *where*, and *what to do
about it*.

This reference consolidates a topic that was previously scattered:
`da-7-machine-learning` listed "drift detection, train-serve skew, MLflow, W&B"
as MLOps sub-bullets, and `da-17-feature-engineering-and-feature-stores` covered
feature-level monitoring. Neither owns model monitoring as a discipline. This
file does.

## Contents

- [When to use](#when-to-use-this-reference) · [When NOT to use](#when-not-to-use-this-reference)
- [1. The drift taxonomy](#1-the-drift-taxonomy--name-the-failure-before-you-measure-it)
- [2. Drift detection tests](#2-drift-detection-tests--pick-by-data-type-and-sample-size)
- [3. Sequential / streaming concept-drift detectors](#3-sequential--streaming-concept-drift-detectors)
- [4. Performance monitoring with delayed/absent ground truth](#4-performance-monitoring-with-delayed-or-absent-ground-truth)
- [5. Estimating performance without labels (NannyML CBPE & DLE)](#5-estimating-performance-without-labels--nannyml-cbpe--dle)
- [6. Training-serving skew & slice/fairness monitoring](#6-training-serving-skew-and-slice-based--fairness-monitoring)
- [7. Alerting, retraining triggers & the retraining loop](#7-alerting-retraining-triggers-and-the-monitoring---retraining-loop)
- [8. Tooling landscape](#8-tooling-landscape)
- [9. Model monitoring vs data observability vs LLM observability](#9-model-monitoring-vs-data-observability-vs-llm-observability)
- [Anti-patterns](#anti-patterns) · [Troubleshooting](#troubleshooting) · [References](#references)

## When to use this reference

- A deployed model is silently degrading and you need to know why.
- You are detecting **data drift, concept drift, or prediction drift**.
- You are choosing a **drift test** (PSI vs KS vs Wasserstein vs JS).
- You need to estimate live performance with **delayed or absent labels**.
- You are detecting **training-serving skew**.
- You are monitoring by **slice/segment** and watching for **fairness drift**.
- You are setting **alert thresholds** and designing **retraining triggers**.
- You are comparing model-monitoring **tools**.

## When NOT to use this reference

- Flagging unusual individual records/events as the goal → `da-16-anomaly-detection`
  (model monitoring *uses* outlier detection on inputs, but its question is "is
  the model still good?", not "is this row weird?").
- Building features or a feature store, including feature freshness/lineage →
  `da-17-feature-engineering-and-feature-stores`.
- Health of the data/pipelines themselves (freshness, volume, schema, lineage) →
  `da-19-data-observability` (under `da-data-engineering-platform`).
- Selecting the model or its offline evaluation metrics before deployment →
  `da-7-machine-learning`.

---

## 1. The drift taxonomy — name the failure before you measure it

Drift is not one thing. Decompose the joint distribution `P(X, Y) = P(X)·P(Y|X)`
and ask *which factor moved*. Getting this wrong is the most common monitoring
mistake, because the detection method and the fix differ by type.

| Drift type | What changes | Formal statement | Hurts accuracy? | Detectable without labels? |
|---|---|---|---|---|
| **Data drift / covariate shift** | Input distribution | `P(X)` changes, `P(Y\|X)` constant | Sometimes (only if it moves inputs into worse-served regions) | Yes — compare input distributions |
| **Concept drift** | Input→output relationship | `P(Y\|X)` changes | Almost always | **No** — needs labels (or a proxy) |
| **Prediction / output drift** | Model output distribution | `P(Ŷ)` changes | Symptom, not cause | Yes — compare prediction distributions |
| **Label / prior-probability shift** | Target base rate | `P(Y)` changes, `P(X\|Y)` constant | Often (esp. calibration/thresholds) | Only once labels arrive |
| **Feature drift** | One input feature's distribution | `P(X_j)` changes | Depends on feature importance | Yes — per-feature tests |

Key consequences that trip teams up:

- **Data drift is not the goal — performance is.** A feature can drift hard
  with zero accuracy impact (and a stable feature set can hide concept drift).
  Drift on inputs is an *early-warning proxy*, not proof of harm. Weight drift
  alerts by feature importance; alerting on every drifting column produces noise.
- **Concept drift is the dangerous one and the hard one.** Because `P(Y|X)`
  changed, the *only* way to confirm it is with ground truth (or a labelled
  proxy). No input-only method can see it. This is the single most important
  asymmetry in the whole discipline (it dictates §4–§5).
- **Prediction drift is a symptom.** Shifting output distributions tell you
  *something* changed upstream; pair it with input drift and (delayed) outcomes
  to localize the cause.
- Drift also varies by **temporal pattern**: *sudden* (a deploy, a market
  shock), *gradual* (slow population change), *incremental*, *recurring/seasonal*.
  Pattern drives whether scheduled vs trigger-based retraining fits (§7).

---

## 2. Drift detection tests — pick by data type and sample size

Two families: **statistical hypothesis tests** (give a p-value; oversensitive on
large N) and **distance/divergence metrics** (give a magnitude; need a
threshold, but are not fooled by huge samples). On production-scale data, prefer
distance metrics or sample before testing — with millions of rows a KS or
Chi-square test flags everything as "significant".

| Test / metric | Data type | What it measures | Notes & thresholds |
|---|---|---|---|
| **PSI** (Population Stability Index) | Numeric (binned) & categorical | Binned distribution shift | The credit-industry default. `<0.1` no shift; `0.1–0.25` moderate; `≥0.25` significant. Stable on large N. **Blind to concept drift** — measures inputs only. |
| **KS** (Kolmogorov–Smirnov) | Numeric, continuous | Max gap between CDFs | Non-parametric, no binning. **Oversensitive on large samples** — sample down or pair with an effect size. |
| **Chi-square** | Categorical | Frequency-table divergence | Standard categorical test; same large-N oversensitivity. |
| **KL divergence** | Numeric/categorical (as distributions) | Relative entropy `D_KL(P‖Q)` | **Asymmetric** (`D_KL(P‖Q) ≠ D_KL(Q‖P)`), unbounded, undefined where `Q=0`. Good default on larger sets if you respect asymmetry. |
| **Jensen–Shannon (JS) divergence** | Numeric/categorical | Symmetrized, smoothed KL | Bounded `[0,1]` (log base 2), symmetric, always defined. Easier to threshold than KL. Used by Vertex AI. |
| **Wasserstein / Earth-Mover's Distance** | **Numeric only** | "Work" to morph one distribution into another | Handles **non-overlapping** distributions and outliers gracefully where KL/JS struggle; no natural universal threshold — calibrate on history. Not for categoricals (assumes bin distance). |
| **L-infinity distance** | Categorical (and numeric) | Max per-bucket probability gap | Simple, interpretable max-deviation; the default categorical metric in Vertex AI Model Monitoring. |
| **MMD** (Maximum Mean Discrepancy) | Multivariate / embeddings | Kernel two-sample distance | Multivariate drift in one shot; core to Alibi Detect; usable on embeddings for text/image. |
| **Classifier two-sample test (C2ST)** | Any (multivariate) | Train a classifier to tell reference from production | Drift signal when held-out **AUC ≳ 0.55–0.60**; naturally multivariate and interpretable via feature importance. |

Rules of thumb: numeric univariate → **PSI or Wasserstein** (KS if N is modest);
categorical → **PSI, L-infinity, or Chi-square**; high-dimensional / embeddings
→ **MMD or C2ST**. Always carry a *reference window* (training or a stable
production baseline) and compare a rolling *analysis window* against it.

---

## 3. Sequential / streaming concept-drift detectors

The §2 tests compare two *batches*. For *streams* — where you want to detect a
change point online with bounded memory — use sequential detectors. The
performance-based ones (DDM/EDDM) monitor the **error stream** and therefore
*do* catch concept drift (they require labels); the signal-based ones
(Page-Hinkley/CUSUM/ADWIN) monitor any univariate stream (error rate, a metric,
a feature mean).

| Detector | Monitors | Mechanism | Best at |
|---|---|---|---|
| **DDM** (Drift Detection Method) | Binary classifier error rate | Error treated as binomial; raises **warning** then **drift** when error + std exceed learned minimums | Sudden/abrupt drift on a labelled stream |
| **EDDM** (Early DDM) | Distance *between* errors | Tracks mean distance between misclassifications | **Gradual** drift (earlier than DDM); more noise-sensitive |
| **ADWIN** (ADaptive WINdowing) | Any real-valued stream | Adaptive window; splits and shrinks when two sub-windows' means differ beyond a Hoeffding bound | Distribution-free change detection; doubles as a forgetting mechanism for adaptive learners |
| **Page-Hinkley** | Univariate Gaussian-ish stream | Cumulative deviation from running mean with a sensitivity/`δ` and threshold `λ` | Detecting a mean shift in a metric/error signal |
| **CUSUM** | Univariate stream | Cumulative sum of deviations vs a slack; alarm past a threshold | Classic mean-shift change-point detection |

Implementations live in **River** (online ML) and the older `scikit-multiflow`.
(For change-point detection as an *anomaly* discipline — PELT, BOCPD — see
`da-16-anomaly-detection`; here the change point is a *retraining trigger*.)

---

## 4. Performance monitoring with delayed or absent ground truth

The hardest production reality: you usually cannot measure true accuracy *now*,
because labels arrive late or never. A credit-default label may take 12–18
months; a fraud label depends on chargebacks that land weeks later; some
predictions are never resolved at all. This **label/ground-truth lag** is a
blind spot — the model runs without knowing how well it is doing.

Standard playbook — **two monitoring loops**:

1. **Real-time loop (no labels):** monitor what you *can* see immediately — input
   data drift, **prediction/output drift**, data-quality and schema checks, and
   system health (latency, error rate). These are *proxy/leading* signals.
2. **Delayed loop (labels):** once ground truth lands, compute the *true*
   performance metrics (accuracy, precision/recall, AUC, MAE/RMSE) on the matched
   predictions, and reconcile against what the proxies predicted.

**Proxy metrics** are signals correlated with the outcome you cannot yet see:
prediction-score drift, share of low-confidence predictions, rate of out-of-range
inputs, business KPIs (approval rate, click rate), or partial/early labels. Treat
them as early warnings, never as truth — validate the proxy↔outcome correlation
when labels eventually arrive, and watch for proxies that themselves drift.

Practical tactics: **join predictions to outcomes** on a stable key so delayed
labels backfill the right rows; track a **label-arrival/coverage curve** (how
much ground truth is in yet) so you don't over-read a half-labelled window;
**time-shift** delayed metrics back to the prediction timestamp when plotting.

---

## 5. Estimating performance without labels — NannyML CBPE & DLE

Beyond proxies, you can *estimate* the metric value itself. **NannyML** is the
reference open-source library here, with two algorithms. Both estimate
performance under **covariate shift / data drift** and **both fail under concept
drift** — by construction, since `P(Y|X)` changing cannot be seen without labels.

- **CBPE — Confidence-Based Performance Estimation (classification).** Uses the
  model's **predicted probabilities**. Intuition: if the model says 0.9 and is
  *well-calibrated*, ~90% of such predictions are correct, so you can reconstruct
  a confusion matrix and estimate **accuracy, precision, recall, F1, ROC-AUC,
  specificity** (binary; macro-averaged + one-vs-rest AUC for multiclass).
  Requires **calibrated** probabilities — NannyML calibrates post-hoc with
  **isotonic regression** on reference data (not on training data, which would
  bias it). Needs enough observations; cannot rescue an uncalibrated or
  below-0.5 model. **M-CBPE** is a multi-calibration variant that better handles
  covariate shift.
- **DLE — Direct Loss Estimation (regression).** Regression has no built-in
  confidence score, so DLE trains a **"nanny" (child-loss) model** — by default
  **LightGBM** — that predicts the *loss* (absolute/squared error) of the
  monitored model per observation, using the monitored model's features and
  predictions as inputs. Averaging the estimated loss yields estimated
  **MAE/RMSE/etc.** The nanny need not beat the monitored model; it just learns
  *where* error is large. Works when noise is heteroscedastic (error depends on
  inputs); gives little signal under homoscedastic noise.

Hard limit to state plainly: **no input-only method detects concept drift.** If
`P(Y|X)` shifts, CBPE/DLE will report "all good" while the model is wrong. Always
pair estimated performance with the delayed *realized* performance from §4 once
labels arrive, and treat estimation as a bridge across the label lag, not a
replacement for ground truth.

---

## 6. Training-serving skew and slice-based / fairness monitoring

### Training-serving skew

Skew is drift's static cousin: a discrepancy between the data/code path at
**training** and at **serving** — same point in time, not over time. Causes:
different feature-computation code offline vs online, a stale serving pipeline,
different preprocessing, or a feature available at train time but missing/late at
serve time. Detect it by comparing the **serving feature distribution against the
training baseline** (Vertex AI does exactly this for skew, vs comparing recent to
older production traffic for drift). The durable fix is upstream — shared feature
transforms / a feature store with point-in-time correctness (see
`da-17-feature-engineering-and-feature-stores`) — not a monitor.

### Slice / segment monitoring and fairness drift

**Aggregate metrics lie.** A model can hold 95% overall accuracy while collapsing
to 80% on a key segment, because the good majority masks a degrading minority.
Always monitor performance **per slice** — by cohort, region, device, customer
tier, and protected attributes. Subgroup-level drift surfaces degradation that is
invisible globally; tools like DriftInspector identify and efficiently track
interpretable subgroups over a model's lifetime.

**Fairness drift** is the special case where degradation lands disproportionately
on protected groups: a model that was near-parity at launch drifts to materially
worse error for a minority group. Monitor group-fairness metrics over time —
**demographic/statistical parity, equalized odds, predictive parity** — not just
at launch. Fiddler in particular ships bias/fairness assessment as a first-class
monitor. Slice + fairness monitoring is also your defense against the failure
mode where overall numbers look fine to leadership while a subgroup quietly
breaks.

### Input outlier / anomaly & adversarial detection

A complementary input-side guardrail: flag individual inputs that are
out-of-distribution or adversarial *before* trusting the prediction. **Alibi
Detect** (Seldon) is the reference open-source library for **outlier, adversarial,
and drift detection** over tabular/text/image/time-series, with detectors
including **KS, Chi-square, MMD, learned-kernel MMD, classifier (C2ST),
context-aware**, and **online MMD**. This is the monitoring-side *use* of anomaly
detection; for the methods themselves (Isolation Forest, LOF, autoencoders) see
`da-16-anomaly-detection`.

---

## 7. Alerting, retraining triggers, and the monitoring -> retraining loop

Detection is worthless without a response. The loop:
**monitor → detect (drift / performance drop / skew) → alert → diagnose →
decide → retrain & validate → (canary/shadow) redeploy → monitor.**

**Retraining-trigger strategies** (most teams combine them):

| Trigger | Fires when | Strength | Weakness |
|---|---|---|---|
| **Scheduled / time-based** | Fixed cadence (weekly, monthly) | Simple, predictable, easy to automate | Blind to *sudden* drift; retrains needlessly when nothing changed |
| **Performance-based** | True or estimated metric drops below a threshold (e.g. accuracy < 0.85) | Tied directly to business impact; catches concept drift | Needs labels (or a trusted estimate); reacts *after* damage |
| **Drift-triggered** | An input/prediction drift metric crosses a threshold (KS, PSI, JS) | Acts *before* labels arrive; earliest warning | Can fire on harmless drift; needs importance weighting |
| **Data-volume** | Enough new (ideally labelled) data has accumulated | Pragmatic when labels are the bottleneck | Not tied to whether the model actually degraded |

Alerting hygiene: set thresholds from a **historical baseline**, not a guess; use
warning + alarm bands (as DDM does); weight by feature importance; deduplicate and
rate-limit so a noisy feature can't page on-call repeatedly. Before auto-promoting
a retrained model, **validate it beats the incumbent** on a holdout (and ideally
in shadow/canary) — automated retraining without a quality gate can ship a *worse*
model. In a mature MLOps setup the whole loop is a pipeline (e.g. drift alert →
CI/CD retraining job → eval gate → registry → canary), tying this reference back to
the MLOps material in `da-7-machine-learning`.

---

## 8. Tooling landscape

| Tool | Type | Strengths / what's distinctive |
|---|---|---|
| **Evidently AI** | OSS library + hosted | Reports, **Test Suites**, and monitoring **Dashboards**; rich drift/quality **Presets**; tabular, text, embeddings, ranking. Frames "ML monitoring is a *subset* of ML observability." Great default OSS starting point. |
| **NannyML** | OSS library + cloud | The **performance-estimation-without-labels** specialist (**CBPE**, **DLE**, **M-CBPE**); multivariate drift; ties drift to *estimated business impact*. |
| **Arize AI** | Commercial platform (+ OSS **Phoenix**) | Real-time monitoring, drift, **AI-assisted root-cause**; strong LLM support; large enterprise footprint. Phoenix is the open-source tracing/eval sibling. |
| **Fiddler AI** | Commercial platform | Monitoring + drift + **post-hoc explainability** + first-class **bias/fairness** assessment; enterprise compliance (SOC 2 Type 2, HIPAA). |
| **WhyLabs / whylogs** | OSS profiling lib + platform | **whylogs** computes lightweight, privacy-preserving statistical **profiles** of data/predictions; the platform monitors drift/quality on those profiles at scale. Platform open-sourced (Apache 2.0, Jan 2025). |
| **Seldon Alibi Detect** | OSS library | Outlier, **adversarial**, and **drift** detection (KS, Chi-square, MMD, learned-kernel MMD, C2ST, context-aware, online); tabular/text/image/time-series; pairs with Seldon serving. |
| **Amazon SageMaker Model Monitor** | Cloud (AWS) | Four monitor types: **data quality, model quality, bias drift, feature-attribution drift**; baseline + scheduled monitoring; customizable distance metrics (e.g. MMD). |
| **Vertex AI Model Monitoring** | Cloud (GCP) | Logs prediction requests to BigQuery; **training-serving skew** (vs training baseline) and **prediction drift** (recent vs earlier prod); **L-infinity** + **JS divergence** metrics; per-feature thresholds; feature-attribution drift. |
| **MLflow** | OSS lifecycle / registry | Not a monitor per se — the **tracking + Model Registry** backbone that records the baseline, versions, and metrics a monitor compares against; integrates with the retraining pipeline. |

Selection heuristic: OSS-first / batch → **Evidently** (+ **NannyML** when labels
lag); already on a cloud → that cloud's native monitor (**SageMaker** / **Vertex**)
to avoid plumbing; enterprise governance, explainability, and fairness →
**Fiddler** or **Arize**; lightweight at-scale profiling → **whylogs/WhyLabs**;
custom adversarial/OOD guardrails → **Alibi Detect**.

---

## 9. Model monitoring vs data observability vs LLM observability

These overlap and get conflated; the boundary matters for routing.

| | **Model monitoring** (this ref) | **Data observability** (`da-19`) | **LLM observability** |
|---|---|---|---|
| Subject | A deployed **model's** behavior | The **data/pipelines** feeding everything | An **LLM application** (prompts, retrieval, agents) |
| Core signals | Drift, performance, skew, slice/fairness | Freshness, volume, schema, distribution, lineage (the "five pillars") | Traces, token/cost/latency, retrieval context, hallucination/eval scores, prompt versions |
| Central question | "Is the model still good — and where isn't it?" | "Is the data correct, fresh, and complete?" | "Why did this generation/agent step go wrong?" |
| Ground-truth problem | Severe (label lag) → §4–§5 | Usually has a checkable expectation | Often no single ground truth → LLM-as-judge / eval suites |

**Monitoring vs observability** (the verbs): *monitoring* tracks predefined
metrics and tells you *that* something is wrong (known-unknowns, threshold
alerts); *observability* gives the traces/context to explain *why* (unknown-
unknowns, root cause). Evidently's framing — "**monitoring is a subset of
observability**" — applies across all three columns.

Practical dependency: model monitoring *consumes* data observability. If a
freshness/schema break upstream poisons the features, a model "drift" alert is
really a data-quality incident — check `da-19-data-observability` first.
LLM/agent observability is the generative-AI specialization (semantic eval,
tracing); classical model monitoring is the discipline for tabular/structured
predictive models. For the upstream feature side (offline/online feature
monitoring, freshness, lineage) see `da-17-feature-engineering-and-feature-stores`;
for selecting the model and its offline metrics see `da-7-machine-learning`.

---

## Anti-patterns

- **Alerting on drift instead of impact.** Drift ≠ degradation. Unweighted
  per-feature drift alerts bury the real signal. Weight by importance; confirm
  with performance (or its estimate) before paging.
- **Assuming input-only methods catch concept drift.** PSI/KS/CBPE/DLE are all
  blind to `P(Y|X)` change. Concept drift needs labels or a labelled proxy.
- **Trusting aggregate accuracy.** It masks subgroup collapse and fairness drift —
  always monitor slices.
- **No reference window / static thresholds plucked from air.** Baselines and
  thresholds must come from training or a stable production period.
- **Auto-retraining without a quality gate.** Retraining on drifted/poisoned data
  can ship a *worse* model; always validate the candidate beats the incumbent.
- **Ignoring training-serving skew.** A monitor catches it late; the fix is shared
  transforms / a feature store, upstream.
- **KS/Chi-square on millions of rows.** Everything reads "significant." Use
  distance metrics (PSI/Wasserstein) or sample first.
- **Reading a half-labelled window as final.** Track label coverage before
  declaring a true-performance verdict.

---

## Troubleshooting

| Symptom | Likely cause | What to do |
|---|---|---|
| Predictions drifted, inputs look stable | Concept drift, or an upstream code/skew change | Get labels/proxy; diff serving vs training feature code; check for a recent deploy |
| Every feature "drifts" every run | Oversensitive test on large N | Switch to PSI/Wasserstein or sample; raise thresholds off a baseline |
| Accuracy fine overall, complaints from one segment | Subgroup / fairness drift | Add slice monitors; compare group metrics over time |
| Estimated perf (CBPE) says fine, reality is bad | Concept drift (CBPE can't see it) or poorly-calibrated probabilities | Recalibrate (isotonic); wait for/obtain labels; add a labelled proxy |
| Drift alerts constantly, model still performs | Unweighted, harmless input drift | Weight by feature importance; gate alerts on performance/estimate |
| Big train↔serve gap at launch, no time drift | Training-serving skew | Compare serving vs training baseline; unify feature transforms / feature store |
| "Model drift" alert that's really bad data | Upstream freshness/schema/lineage break | Check `da-19-data-observability`; fix the data incident first |

---

## References

Drift taxonomy & detection tests:
- Evidently AI — "Which test is the best? We compared 5 methods to detect data drift on large datasets." https://www.evidentlyai.com/blog/data-drift-detection-large-datasets
- IBM — "What Is Model Drift?" https://www.ibm.com/think/topics/model-drift
- DataCamp — "Understanding Data Drift and Model Drift: Drift Detection in Python." https://www.datacamp.com/tutorial/understanding-data-drift-model-drift
- MachineLearningMastery — "Detecting & Handling Data Drift in Production." https://machinelearningmastery.com/detecting-handling-data-drift-in-production/
- StatsTest — "Drift Detection: KS Test, PSI, and Interpreting Signals." https://www.statstest.com/drift-detection-ks-test-psi-interpret-signals
- Arize AI — "Wasserstein Distance" (glossary). https://arize.com/glossary/wasserstein-distance/
- Springer — "Detecting drifts in data streams using Kullback–Leibler (KL) divergence." https://link.springer.com/article/10.1007/s42488-024-00119-y

Sequential / streaming detectors:
- India Lindsay — "Concept Drift Detection" (ADWIN, DDM, EDDM, Page-Hinkley). https://indialindsay1.medium.com/concept-drift-detection-2667a3360091
- ScienceDirect — "A comparative study on concept drift detectors." https://www.sciencedirect.com/science/article/abs/pii/S0957417414004175
- Iguazio — "Concept Drift Deep Dive: How to Build a Drift-Aware ML System." https://www.iguazio.com/blog/concept-drift-deep-dive-how-to-build-a-drift-aware-ml-system/

Delayed/absent ground truth & performance estimation:
- Evidently AI — "Model monitoring for ML in production: a comprehensive guide." https://www.evidentlyai.com/ml-in-production/model-monitoring
- NannyML — "Why is Machine Learning Monitoring in production hard?" (availability of ground truth). https://www.nannyml.com/blog/availability-of-ground-truth
- NannyML docs — "Estimation of Performance of the Monitored Model" (CBPE, DLE, M-CBPE). https://nannyml.readthedocs.io/en/stable/how_it_works/performance_estimation.html
- arXiv — "Confidence-based Estimators for Predictive Performance in Model Monitoring." https://arxiv.org/html/2407.08649v1
- Chris Zhang — "The Hidden Lag: Understanding Delayed Ground Truth in Machine Learning." https://zhanghaolin66.medium.com/the-hidden-lag-understanding-delayed-ground-truth-in-machine-learning-377726d4d739

Slice / fairness drift:
- arXiv — "Detecting Interpretable Subgroup Drifts" (DriftInspector). https://arxiv.org/html/2408.14682v1
- Superwise — "Measuring Subgroup Performance in Machine Learning." https://superwise.ai/blog/measuring-performance-sub-groups/
- arXiv — "Who experiences large model decay and why? A Hierarchical Framework for Diagnosing Heterogeneous Performance Drift." https://arxiv.org/html/2506.00756

Training-serving skew, retraining, & tooling:
- Google Cloud — "Monitor feature skew and drift | Vertex AI." https://cloud.google.com/vertex-ai/docs/model-monitoring/using-model-monitoring
- AWS — "Detect NLP data drift using custom Amazon SageMaker Model Monitor." https://aws.amazon.com/blogs/machine-learning/detect-nlp-data-drift-using-custom-amazon-sagemaker-model-monitor/
- Seldon — "alibi-detect: Algorithms for outlier, adversarial and drift detection." https://github.com/SeldonIO/alibi-detect
- EnhancedMLOps — "Automatic Model Retraining: When and How to Do It?" https://enhancedmlops.com/automatic-model-retraining-when-and-how-to-do-it/
- WWT — "MLOps and Drift: Reducing Risk and Ensuring Robust ML Models." https://www.wwt.com/article/mlops-and-drift-reducing-risk-and-ensuring-robust-ml-models

Tool comparison & observability boundary:
- Medium (T. Kandivlikar) — "Comprehensive Comparison of ML Model Monitoring Tools: Evidently AI, Alibi Detect, NannyML, WhyLabs, and Fiddler AI." https://medium.com/@tanish.kandivlikar1412/comprehensive-comparison-of-ml-model-monitoring-tools-evidently-ai-alibi-detect-nannyml-a016d7dd8219
- Uplatz — "A Technical Leader's Comparative Analysis of AI Observability Platforms: Evidently AI, Arize AI, and Fiddler AI." https://uplatz.com/blog/a-technical-leaders-comparative-analysis-of-ai-observability-platforms-evidently-ai-arize-ai-and-fiddler-ai/
- InsightFinder — "ML vs LLM Observability: A Complete Guide." https://insightfinder.com/blog/ml-vs-llm-observability-guide/
- Elastic — "What Is LLM Observability? A Comprehensive Guide." https://www.elastic.co/what-is/llm-observability
