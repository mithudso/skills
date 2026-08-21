<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-4-4-estimation-theory` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-4-4-estimation-theory
description: >-
  Point-estimation theory for statistical inference foundations — the three
  finite-sample criteria (bias, consistency, efficiency) plus mean-square-error
  decomposition, the Cramér-Rao lower bound, Fisher information, UMVUE, relative
  and asymptotic efficiency, and how to compare or judge a candidate estimator.
  TRIGGER: choosing or evaluating an estimator for an unknown parameter; "is
  this estimator unbiased / consistent / efficient"; bias-variance tradeoff of
  an estimator; mean square error of an estimator; Cramér-Rao lower bound or
  Fisher information; minimum-variance unbiased estimator (MVUE/UMVUE);
  relative efficiency of two estimators (e.g. mean vs median); why MLE is
  asymptotically efficient; "what makes a good estimator". SKIP: deriving a
  specific estimator by maximum likelihood or method of moments (estimation
  *method*, not estimator *properties*); confidence-interval or hypothesis-test
  construction (interval/test inference, not point-estimator quality);
  bias-variance tradeoff of a predictive ML *model* / regularization (that is a
  modeling-context skill, not foundations); sampling distribution and standard
  error mechanics (defer to da-1-4-2); frequentist-vs-Bayesian framing of
  estimation (defer to da-1-4-3); population-vs-sample definitions (defer to
  da-1-4-1).
---

# Estimation Theory: Bias, Consistency, Efficiency

Scope: how to judge the quality of a **point estimator** of an unknown
population parameter, within the *Statistical inference foundations* of Data
Analysis. This skill is about the *properties* an estimator can have, not about
how to *derive* one (maximum likelihood, method of moments) and not about
interval estimation or testing — those are separate concepts. The same
"bias/variance" words reappear for predictive ML models; that framing is out of
scope here (see SKIP).

## Setup and notation

- A **parameter** θ is an unknown fixed constant describing the population
  (e.g. a mean μ, a variance σ², a proportion p).
- An **estimator** θ̂ = θ̂(X₁,…,Xₙ) is a function of the sample — itself a
  random variable with its own distribution (the *sampling distribution*).
- An **estimate** is the realized number θ̂ takes on one observed sample.

We ask three finite-sample questions of θ̂ — is it right *on average* (bias), does
it *home in* on θ as data grows (consistency), and is its *spread* as small as
possible (efficiency)? [Properties of Point Estimators, datasciencezone.org;
AnalystPrep CFA L1, "Properties of an Estimator"]

## 1. Bias

Bias is the systematic error — the gap between what the estimator gives on
average and the truth.

```
bias(θ̂) = E[θ̂] − θ
```

- **Unbiased**: `E[θ̂] = θ` for every value of θ ∈ Θ. The estimator is correct
  on average across repeated samples. [Siegrist, LibreTexts 7.5]
- The sample mean X̄ is unbiased for μ. The sample variance with divisor
  (n−1), `S² = Σ(Xᵢ−X̄)²/(n−1)`, is unbiased for σ²; the divisor-n version is
  biased downward — this is the canonical reason for "n−1". [AnalystPrep]
- **Asymptotically unbiased**: bias(θ̂ₙ) → 0 as n → ∞, even if biased at finite
  n. This is weaker than unbiasedness and is one ingredient of consistency.

Unbiasedness alone is not enough: an unbiased estimator can have huge variance,
and a slightly biased estimator can be preferable if it is far less variable
(the motivation for MSE below).

## 2. Mean square error (the unifying criterion)

MSE folds bias and variance into one number and is the standard way to compare
estimators that are *not* both unbiased.

```
MSE(θ̂) = E[(θ̂ − θ)²] = var(θ̂) + [bias(θ̂)]²
```

This decomposition is exact and is the core identity of the topic. [Siegrist,
LibreTexts 7.5; Kozdron, U. Regina Stat 252 ch.3] When θ̂ is unbiased the second
term vanishes and **MSE = variance** — which is why, *among unbiased
estimators*, minimizing variance and minimizing MSE are the same goal.

Practical consequence: a biased estimator can have smaller MSE than an unbiased
one. Comparing only on unbiasedness, or only on variance, can pick the worse
estimator. Compare on MSE when bias differs.

## 3. Consistency

Consistency is the large-sample guarantee: the estimator converges to θ as the
sample grows.

```
θ̂ₙ → θ in probability:   P(|θ̂ₙ − θ| > ε) → 0  as n → ∞,  for every ε > 0
```

This is *(weak) consistency*. [Siegrist; CS109/Stanford 7.7; Fiveable,
"Unbiasedness and consistency"]

- **Sufficient condition**: if bias(θ̂ₙ) → 0 *and* var(θ̂ₙ) → 0, then
  MSE(θ̂ₙ) → 0, which implies convergence in probability (mean-square
  consistency ⇒ weak consistency). This is the easiest practical test. [CS109 7.7]
- Consistency and unbiasedness are independent properties: an estimator can be
  **biased but consistent** (the divisor-n sample variance), or **unbiased but
  inconsistent** (e.g. using only X₁ to estimate μ — unbiased, but its variance
  never shrinks). [datasciencezone.org]

## 4. Efficiency

Efficiency is about *spread*: among competing estimators, the efficient one has
the smallest variance, so it wrings the most information out of the sample.

### Relative efficiency (comparing two estimators)

For two unbiased estimators T₁, T₂ of the same parameter:

```
eff(T₁, T₂) = var(T₂) / var(T₁)
```

If this ratio > 1, T₁ is more efficient (smaller variance). [Wikipedia,
"Efficiency (statistics)"; Kozdron ch.3] Classic example: for a normal
population the sample **median** has efficiency ≈ 2/π ≈ 0.64 relative to the
sample **mean**, i.e. its variance is about π/2 ≈ 1.57× larger — so the mean is
the more efficient estimator of the center. [Wikipedia, "Efficiency (statistics)"]

### Absolute efficiency and the Cramér-Rao lower bound (CRLB)

There is a floor on how small the variance of *any* unbiased estimator can be.
Under regularity conditions (the support of the density does not depend on θ,
and log-likelihood is twice differentiable so the integral and derivative can
swap):

```
var(θ̂) ≥ 1 / I(θ)      (CRLB, for an unbiased estimator of θ)
```

where I(θ) is the **Fisher information**, the expected curvature (negative
second derivative) of the log-likelihood — for an i.i.d. sample of size n,
`I_n(θ) = n·I₁(θ)`, so the bound shrinks like 1/n. [Wikipedia, "Cramér–Rao
bound"; Siegrist 7.5] For a general unbiased estimator of a function λ(θ) the
numerator becomes `(dλ/dθ)²`.

The **efficiency** of an unbiased estimator T is the ratio of the bound to its
actual variance:

```
e(T) = [1 / I(θ)] / var(T) ,    0 < e(T) ≤ 1
```

[Wikipedia, "Efficiency (statistics)"]

- **e(T) = 1**: T attains the CRLB — it is an **efficient estimator**.
- An unbiased estimator that achieves the CRLB is automatically the
  **(uniformly) minimum-variance unbiased estimator (UMVUE/MVUE)** — no unbiased
  estimator can beat it. [Siegrist 7.5; Wikipedia, "Cramér–Rao bound"]
- Finite-sample-efficient estimators are rare; they exist mainly for natural
  parameters of exponential-family distributions. [Wikipedia, "Efficiency (statistics)"]

### Asymptotic efficiency

Most "good" estimators reach the bound only in the limit. An estimator θ̂ₙ is
**asymptotically efficient** if it is consistent, asymptotically normal, and its
asymptotic variance equals the CRLB:

```
√n (θ̂ₙ − θ)  →  N(0, 1/I₁(θ))   in distribution
```

The maximum-likelihood estimator satisfies this under regularity conditions —
which is the main reason MLE is the default estimator: it is asymptotically
unbiased, consistent, and asymptotically efficient. [Wikipedia, "Cramér–Rao
bound"; Fiveable, "Properties of Maximum Likelihood Estimators"] Note: this
skill covers *why MLE is efficient*; deriving the MLE itself is a separate
estimation-method concept (see SKIP).

## How to evaluate a candidate estimator (workflow)

1. **Bias** — compute E[θ̂] and subtract θ. Zero ⇒ unbiased; nonzero-but-→0 ⇒
   asymptotically unbiased.
2. **Variance** — compute var(θ̂). Combine with bias into MSE = var + bias².
3. **Consistency** — check bias → 0 and var → 0 as n → ∞ (sufficient condition).
4. **Efficiency** — if comparing two, take the variance ratio; if judging in
   absolute terms, compute the CRLB via Fisher information and form e(T).
5. **Decide** — prefer the smaller **MSE**. Among unbiased candidates this is the
   smaller-variance / more-efficient one; allow a little bias if it buys a large
   variance reduction.

## Common pitfalls

- **Treating unbiasedness as the goal.** A high-variance unbiased estimator can
  have larger MSE than a low-variance biased one. Optimize MSE. [Siegrist 7.5]
- **Confusing asymptotic unbiasedness with consistency.** Asymptotic
  unbiasedness (bias → 0) is necessary but not sufficient; the variance must
  also vanish. [CS109 7.7]
- **Assuming consistency implies unbiasedness (or vice versa).** They are
  independent. [datasciencezone.org]
- **Applying the CRLB outside its regularity conditions.** If the support
  depends on θ (e.g. Uniform(0, θ)), the bound need not hold and an estimator can
  appear to "beat" it. [Wikipedia, "Cramér–Rao bound"]
- **Comparing efficiency of estimators of different parameters** — relative
  efficiency is only meaningful for estimators of the *same* quantity.
- **Importing the ML-model bias-variance tradeoff vocabulary here.** In this
  foundations context bias/variance are properties of a *parameter estimator*,
  not prediction-error components of a fitted predictive model — those belong to
  the modeling/regularization part of the taxonomy (see SKIP).

## Sources

- Siegrist, *Probability, Mathematical Statistics, and Stochastic Processes*,
  §7.5 "Best Unbiased Estimators", Statistics LibreTexts —
  https://stats.libretexts.org/Bookshelves/Probability_Theory/Probability_Mathematical_Statistics_and_Stochastic_Processes_(Siegrist)/07:_Point_Estimation/7.05:_Best_Unbiased_Estimators
- "Cramér–Rao bound", Wikipedia —
  https://en.wikipedia.org/wiki/Cram%C3%A9r%E2%80%93Rao_bound
- "Efficiency (statistics)", Wikipedia —
  https://en.wikipedia.org/wiki/Efficiency_(statistics)
- "Properties of an Estimator", AnalystPrep CFA Level 1 —
  https://analystprep.com/cfa-level-1-exam/quantitative-methods/properties-of-an-estimator/
- "Properties of Point Estimators: Unbiasedness, Consistency, Efficiency, and
  Sufficiency", datasciencezone.org —
  https://datasciencezone.org/properties-of-point-estimators-unbiasedness-consistency-efficiency-and-sufficiency/
- Stanford CS109 §7.7 "Properties of Estimators II (Consistency)" —
  https://web.stanford.edu/class/archive/cs/cs109/cs109.1218/files/student_drive/7.7.pdf
- Kozdron, U. Regina Stat 252 ch.3 "Bias, Mean-Square Error, Relative
  Efficiency" —
  https://uregina.ca/~kozdron/Teaching/Regina/252Winter16/Handouts/ch3.pdf
