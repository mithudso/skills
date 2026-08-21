<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-3-8-expectation-variance-covariance` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-3-8-expectation-variance-covariance
description: >-
  Probability-theory foundations of expectation (expected value / mean), variance,
  and covariance for random variables — the first and second moments that summarize
  a distribution's center, spread, and linear co-movement. Covers exact definitions
  (discrete and continuous), linearity of expectation, LOTUS, variance of a sum,
  the covariance/correlation relationship, and the common pitfalls.
  TRIGGER: questions about E[X], Var(X), Cov(X,Y), expected value, mean of a random
  variable, variance/standard deviation of a random variable, covariance, correlation
  coefficient, linearity of expectation, variance of a sum, Var(aX+b), E[XY] vs
  E[X]E[Y], "uncorrelated vs independent", LOTUS / law of the unconscious statistician,
  computing moments of a distribution.
  SKIP: descriptive sample statistics computed from a dataset (sample mean/variance/SD,
  pandas/numpy .mean()/.var(), Bessel's n-1 correction) — those are descriptive-stats
  topics, not random-variable moments. SKIP probability distributions themselves
  (normal/binomial/Poisson — see the da-1-3-3-* distribution skills). SKIP the
  covariance/correlation MATRIX, PCA, or regression coefficients (multivariate /
  modeling topics). SKIP correlation-as-feature-selection or correlation-vs-causation
  framed as EDA. This skill is scoped to the probability-theory definitions of the
  moments of random variables.
---

# Expectation, Variance, Covariance (Probability Theory)

Scope: the **probability-theory** treatment of the moments of random variables — how
to define and manipulate `E[X]`, `Var(X)`, and `Cov(X,Y)` from a distribution. This is
distinct from *descriptive sample statistics* computed from observed data (those have an
`n-1` correction and are estimators of these quantities; defer to descriptive-stats
skills) and from the *distributions* themselves (defer to the `da-1-3-3-*` skills).

## 1. Expectation (expected value, mean, first moment)

The expected value is a probability-weighted average of all outcomes — the distribution's
center of mass.

- **Discrete:** `E[X] = Σ_x x · p_X(x)` over the support.
- **Continuous:** `E[X] = ∫_{-∞}^{∞} x · f_X(x) dx`.

The mean exists (is well-defined and finite) only if the sum/integral converges
absolutely, i.e. `E[|X|] < ∞`. Heavy-tailed distributions (e.g. Cauchy) have **no**
expected value — a frequent trap ([Expected value, Wikipedia](https://en.wikipedia.org/wiki/Expected_value)).

### LOTUS — Law of the Unconscious Statistician

To get the expectation of a *function* of `X`, you do **not** need the distribution of
`g(X)`; integrate/sum against the distribution of `X`:

- Discrete: `E[g(X)] = Σ_x g(x) · p_X(x)`
- Continuous: `E[g(X)] = ∫ g(x) · f_X(x) dx`

([LOTUS, Wikipedia](https://en.wikipedia.org/wiki/Law_of_the_unconscious_statistician);
[probabilitycourse 4.1.2](https://www.probabilitycourse.com/chapter4/4_1_2_expected_val_variance.php)).
LOTUS is what makes variance computable as `E[(X-μ)²]`.

### Properties of expectation

Let `a, b` be constants and `X, Y` random variables ([statlect: properties of expected value](https://www.statlect.com/fundamentals-of-probability/expected-value-properties)):

- **Constant:** `E[a] = a`.
- **Scaling:** `E[aX] = a·E[X]`.
- **Linearity (sums):** `E[X + Y] = E[X] + E[Y]` — and more generally
  `E[a₁X₁ + … + aₖXₖ] = a₁E[X₁] + … + aₖE[Xₖ]`.
- **Monotonicity:** if `X ≤ Y` almost surely, then `E[X] ≤ E[Y]`.
- **Product under independence:** if `X ⫫ Y`, then `E[XY] = E[X]·E[Y]`. (Without
  independence this generally fails.)

**Key fact about linearity:** `E[X+Y] = E[X]+E[Y]` holds **regardless of dependence** —
the variables need not be independent ([Linearity of expectation, Brilliant](https://brilliant.org/wiki/linearity-of-expectation/)).
This is what makes expectation so easy to compute on sums of dependent indicators.

## 2. Variance and standard deviation (second central moment)

Variance measures spread around the mean.

```
Var(X) = E[(X − μ)²]          where μ = E[X]
       = E[X²] − (E[X])²      (computational form)
```

([Covariance and Correlation, LibreTexts / Siegrist](https://stats.libretexts.org/Bookshelves/Probability_Theory/Probability_Mathematical_Statistics_and_Stochastic_Processes_(Siegrist)/04:_Expected_Value/4.05:_Covariance_and_Correlation)).

Standard deviation `sd(X) = √Var(X)` restores the original units. Variance is always `≥ 0`.

### Properties of variance

- **Constant:** `Var(a) = 0`.
- **Shift invariance + scaling:** `Var(aX + b) = a²·Var(X)`. The additive constant `b`
  drops out (it shifts the mean, not the spread); the multiplier is **squared**
  ([probabilitycourse 4.1.2](https://www.probabilitycourse.com/chapter4/4_1_2_expected_val_variance.php)).
- **Sum (general):** `Var(X + Y) = Var(X) + Var(Y) + 2·Cov(X, Y)`.
- **Sum (independent / uncorrelated):** `Var(X + Y) = Var(X) + Var(Y)`.

**Variance is NOT linear.** Unlike expectation, you cannot move it across a sum unless the
cross-covariance vanishes ([Linearity of expectation, Fiveable](https://fiveable.me/key-terms/introduction-probability/linearity-of-expectation)).

## 3. Covariance (joint second moment)

Covariance measures the *linear* co-movement of two random variables.

```
Cov(X, Y) = E[(X − E[X])(Y − E[Y])]
          = E[XY] − E[X]·E[Y]      (computational form)
```

Sign reading: positive → they tend to move together; negative → opposite directions;
near zero → little linear association ([Covariance, Britannica](https://www.britannica.com/topic/covariance)).

### Properties of covariance

([Covariance and Correlation, LibreTexts / Siegrist](https://stats.libretexts.org/Bookshelves/Probability_Theory/Probability_Mathematical_Statistics_and_Stochastic_Processes_(Siegrist)/04:_Expected_Value/4.05:_Covariance_and_Correlation)):

- **Self-covariance:** `Cov(X, X) = Var(X)`.
- **Symmetry:** `Cov(X, Y) = Cov(Y, X)`.
- **Bilinearity:** `Cov(aX + bY, Z) = a·Cov(X, Z) + b·Cov(Y, Z)`.
- **Constants don't covary:** `Cov(X, c) = 0`.
- **Affine:** `Cov(aX + b, cY + d) = ac·Cov(X, Y)`.

### Correlation coefficient

Covariance has the units of `X·Y`, so its magnitude is not directly interpretable.
Normalize to get a unitless measure in `[-1, 1]`:

```
ρ(X, Y) = Cov(X, Y) / ( sd(X) · sd(Y) )            with  −1 ≤ ρ ≤ 1
```

`ρ = ±1` exactly when `Y` is an affine function of `X` (perfect linear relationship)
([probabilitycourse: covariance & correlation](https://www.probabilitycourse.com/chapter5/5_3_1_covariance_correlation.php)).

## 4. The pitfalls (high-yield)

1. **Uncorrelated ≠ independent.** Independence implies `Cov = 0`; the converse is
   false. A nonlinear dependence (e.g. `Y = X²` with `X` symmetric about 0) gives
   `Cov(X, Y) = 0` while `X` and `Y` are fully dependent. Covariance/correlation only
   detect *linear* association ([LibreTexts / Siegrist](https://stats.libretexts.org/Bookshelves/Probability_Theory/Probability_Mathematical_Statistics_and_Stochastic_Processes_(Siegrist)/04:_Expected_Value/4.05:_Covariance_and_Correlation)).
2. **Don't distribute variance over a sum** unless the covariance term is zero — forgetting
   the `2·Cov(X,Y)` term is the most common error.
3. **`E[g(X)] ≠ g(E[X])` in general** (Jensen's inequality). For convex `g`,
   `E[g(X)] ≥ g(E[X])`; in particular `E[X²] ≥ (E[X])²`, which is why `Var(X) ≥ 0`.
4. **`a²`, not `a`, in `Var(aX)`** — a units error: variance scales with the square.
5. **Mean may not exist.** Always check absolute convergence before assuming `E[X]` is finite.
6. **Population vs sample.** These definitions are the *population* (true-distribution)
   moments. The sample variance estimator uses `n−1` (Bessel's correction); do not confuse
   the two (defer to descriptive-statistics skills).

## 5. Worked example

Two fair dice, `X` = first roll, `Y` = sum of both.

- `E[X] = (1+2+3+4+5+6)/6 = 3.5`. By linearity, `E[Y] = E[X₁]+E[X₂] = 7`.
- `Var(X) = E[X²] − (E[X])² = 91/6 − 12.25 ≈ 2.917`.
- For independent dice, `Cov(X₁, X₂) = 0`, so `Var(Y) = Var(X₁)+Var(X₂) ≈ 5.833`.
- `Cov(X, Y) = Cov(X₁, X₁+X₂) = Cov(X₁,X₁) + Cov(X₁,X₂) = Var(X₁) + 0 ≈ 2.917`
  (bilinearity).
- `ρ(X, Y) = 2.917 / (√2.917 · √5.833) ≈ 0.707` — the first roll is moderately,
  positively correlated with the total, as expected.

## Sources

1. [Expected value — Wikipedia](https://en.wikipedia.org/wiki/Expected_value) — definitions (discrete/continuous), existence conditions.
2. [Properties of the expected value — StatLect](https://www.statlect.com/fundamentals-of-probability/expected-value-properties) — linearity, monotonicity, product-under-independence.
3. [Covariance and Correlation — LibreTexts (Siegrist, *Probability, Mathematical Statistics and Stochastic Processes*)](https://stats.libretexts.org/Bookshelves/Probability_Theory/Probability_Mathematical_Statistics_and_Stochastic_Processes_(Siegrist)/04:_Expected_Value/4.05:_Covariance_and_Correlation) — covariance/variance formulas, symmetry, bilinearity, correlation bounds, uncorrelated≠independent.
4. [Covariance & Correlation — probabilitycourse.com](https://www.probabilitycourse.com/chapter5/5_3_1_covariance_correlation.php) — correlation coefficient, variance of a sum.
5. [Expected Value and Variance (4.1.2) — probabilitycourse.com](https://www.probabilitycourse.com/chapter4/4_1_2_expected_val_variance.php) — LOTUS, Var(aX+b)=a²Var(X).
6. [Law of the unconscious statistician — Wikipedia](https://en.wikipedia.org/wiki/Law_of_the_unconscious_statistician) — LOTUS statement.
7. [Linearity of Expectation — Brilliant](https://brilliant.org/wiki/linearity-of-expectation/) — linearity holds without independence.
8. [Covariance — Britannica](https://www.britannica.com/topic/covariance) — definition and sign interpretation.
