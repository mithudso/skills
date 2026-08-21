<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-3-3-2-binomial` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-3-3-2-binomial
description: >-
  The binomial distribution as a named probability distribution in data-analysis
  foundations — the discrete distribution of the number of successes in a fixed
  number n of independent, identical Bernoulli trials, each with success
  probability p, with mean np and variance np(1-p). Covers the PMF, the four
  binomial assumptions (BINS: Binary, Independent, fixed N, Same p), the
  Bernoulli-sum relationship, mean/variance/mode/skewness, the MLE of p and
  confidence intervals for a proportion (Wald vs Wilson vs Clopper-Pearson), the
  normal approximation with continuity correction, the Poisson limit, and the
  sampling-without-replacement (hypergeometric) pitfall. TRIGGER when a
  data-analysis task asks: what is the binomial distribution, its PMF/mean/
  variance, when count-of-successes data is binomial, how to compute P(X=k) for
  k successes in n trials, how to build a confidence interval for a proportion,
  or whether n large + p small/large justifies a Poisson or normal approximation.
  SKIP when the work is logistic regression / GLM modeling of binary outcomes (a
  modeling technique, not this distribution-definition node), the single-trial
  Bernoulli framing or the geometric/negative-binomial "trials until success"
  count, the Poisson distribution (use da-1-3-3-3-poisson), the normal
  distribution (use da-1-3-3-1-normal), the general "what is a probability
  distribution" parent (use da-1-3-3-probability-distributions), or PMF/PDF
  mechanics in the abstract (use da-1-3-2-probability-mass-density-functions).
---

# Binomial distribution

Taxonomy location: Data Analysis > Foundations & Theory > Probability theory >
Probability distributions > Binomial.

This skill is scoped to the binomial distribution **as a named member of the
probability-distribution catalog** — its definition, PMF, parameters, moments,
defining assumptions, and the estimation/approximation relationships an analyst
needs to decide whether a count of successes is binomial and how to reason about
it. It is not a guide to logistic regression (a GLM modeling technique for binary
outcomes that lives under regression/modeling, not under this distribution-
definition node), nor to the single-trial Bernoulli framing or the "number of
trials until the first success" geometric/negative-binomial family.

## Definition

The binomial distribution is a **discrete** probability distribution that gives
the probability of observing exactly k successes in a fixed number n of
independent Bernoulli trials, where each trial has the same probability p of
success and 1 − p of failure [Wikipedia; LibreTexts; Statistics Fundamentals]. A
**Bernoulli trial** is a single experiment with exactly two outcomes, conventionally
labeled "success" and "failure" [LibreTexts; Statistics Fundamentals].

A binomial random variable is the **sum of n independent and identically
distributed Bernoulli random variables**: if X1, …, Xn each take the value 1 with
probability p and 0 with probability 1 − p, then X = X1 + … + Xn ~ Binomial(n, p)
[Wikipedia; LibreTexts]. The Bernoulli distribution is the special case n = 1
[Wikipedia].

Notation: X ~ B(n, p) or Binomial(n, p), with n ∈ {0, 1, 2, …} a fixed integer
and p ∈ [0, 1].

## Probability mass function

For k = 0, 1, 2, …, n:

  P(X = k) = C(n, k) · p^k · (1 − p)^(n − k)

where C(n, k) = n! / (k! · (n − k)!) is the binomial coefficient ("n choose k")
[Wikipedia; LibreTexts; Number Analytics]. The three pieces read intuitively:

- **C(n, k)** counts the number of distinct arrangements of k successes among the
  n trials,
- **p^k** is the probability of the k successes, and
- **(1 − p)^(n − k)** is the probability of the remaining n − k failures
  [Number Analytics; Statistics Fundamentals].

It satisfies the two PMF axioms: P(X = k) ≥ 0 for every k, and the sum over
k = 0…n equals 1 — this follows from the binomial theorem, since
Σ C(n, k) p^k (1−p)^(n−k) = (p + (1−p))^n = 1 [Probabilistic World]. The
cumulative distribution function P(X ≤ k) is the partial sum of the PMF and has no
simpler closed form; it is computed by summation or via the regularized incomplete
beta function [Wikipedia].

## The four assumptions (BINS)

A count is binomial only when **all four** conditions hold [LibreTexts; Statistics
By Jim; VT Pressbooks]:

1. **B — Binary outcomes.** Each trial results in one of exactly two outcomes
   (success / failure).
2. **I — Independent trials.** The outcome of one trial does not affect any other.
3. **N — fixed Number of trials.** n is set in advance, not determined by the data.
4. **S — Same probability.** The success probability p is identical on every trial.

The mnemonic is **BINS**. Each failed assumption points to a different distribution:
trials-until-first-success → geometric; counts-until-r-successes → negative
binomial; trials drawn without replacement so p changes → hypergeometric (see
Pitfalls); n itself random → not binomial.

## Parameters, mean, variance, and shape

For X ~ B(n, p) [Wikipedia; Statistics By Jim]:

| Quantity | Value |
|---|---|
| Mean | E[X] = np |
| Variance | Var(X) = np(1 − p) |
| Std. deviation | √(np(1 − p)) |
| Mode | ⌊(n + 1)p⌋ (two modes if (n+1)p is an integer) |
| Median | within [⌊np⌋, ⌈np⌉]; equals np when np is an integer |
| Skewness | (1 − 2p) / √(np(1 − p)) |
| MGF | M_X(t) = (1 − p + p·e^t)^n |

Shape notes [Wikipedia; Probabilistic World]:

- The distribution is **symmetric when p = 0.5**, right-skewed when p < 0.5, and
  left-skewed when p > 0.5. Skewness shrinks toward 0 as n grows.
- Variance is **maximized at p = 0.5** (where np(1−p) = n/4) and shrinks to 0 as
  p approaches 0 or 1.
- Unlike the Poisson, the binomial is **underdispersed relative to its mean**:
  Var(X) = np(1−p) < np = E[X] whenever 0 < p < 1.

**Sum of independent binomials with the same p:** if X ~ B(n, p) and Y ~ B(m, p)
are independent, then X + Y ~ B(n + m, p) [Wikipedia]. (This fails if the two p
differ; that case is the Poisson binomial.)

## Estimating p from data

The **maximum-likelihood estimator** of p, given x successes in n trials, is the
sample proportion p̂ = x / n. It is unbiased and is the minimum-variance unbiased
estimator [Wikipedia].

**Confidence intervals for a proportion p.** Several methods exist, and the choice
matters because the simplest one is the least reliable [Wikipedia "Binomial
proportion confidence interval"; Towards Data Science; DescTools]:

- **Wald (normal-approximation):** p̂ ± z · √(p̂(1 − p̂)/n). Simple, but its
  coverage is **erratically poor** for small n and for p near 0 or 1 — avoid it as
  a default.
- **Wilson score:** centers the interval using the score statistic and adds z²
  terms to the denominator. Accurate across most p and does not over-cover. A good
  general-purpose default.
- **Agresti–Coull:** an approximation that adds z²/2 pseudo-successes and failures
  (p̃ = (x + z²/2)/(n + z²)); simpler than Wilson, similar behavior for moderate n.
- **Clopper–Pearson (exact):** inverts the binomial CDF; coverage is **never below**
  the nominal level but is usually **conservative** (intervals too wide), so it
  needs larger samples for the same precision.

## Approximations

The binomial is the parent of two classic limits an analyst should recognize:

**Normal approximation (de Moivre–Laplace).** When n is large and p not too extreme,
B(n, p) ≈ Normal(μ = np, σ² = np(1 − p)). Common rule of thumb: use it when
**np ≥ 5 and n(1 − p) ≥ 5** (some texts require ≥ 10 for tighter accuracy)
[STAT 350; Wikipedia; Statology]. Because the binomial is discrete and the normal
continuous, apply a **continuity correction** of ±0.5 [Statology; Sparkl]:

- P(X ≤ k) ≈ Φ((k + 0.5 − μ)/σ)
- P(X ≥ k) ≈ 1 − Φ((k − 0.5 − μ)/σ)
- P(X = k) ≈ Φ((k + 0.5 − μ)/σ) − Φ((k − 0.5 − μ)/σ)

The correction matters most when n is moderate and p is far from 0.5.

**Poisson approximation.** When n is large and p is small so that λ = np stays
moderate, B(n, p) ≈ Poisson(np) [Wikipedia; STAT 350]. Common rule of thumb:
**n ≥ 20 with p ≤ 0.05**, or **n ≥ 50 with p ≤ 0.1**. This is the "rare events in
many trials" regime (see da-1-3-3-3-poisson for the target distribution).

## Worked example

A quality check pulls n = 20 items from a line where each is independently
defective with probability p = 0.1. Let X = number of defective items, so
X ~ B(20, 0.1).

- Mean: E[X] = np = 20 × 0.1 = 2 defects expected.
- Variance: np(1 − p) = 20 × 0.1 × 0.9 = 1.8; SD ≈ 1.34.
- Exactly 3 defects: P(X = 3) = C(20, 3) · 0.1³ · 0.9¹⁷
  = 1140 × 0.001 × 0.16677 ≈ **0.190**.
- At least 1 defect: P(X ≥ 1) = 1 − P(X = 0) = 1 − 0.9²⁰ ≈ 1 − 0.1216 = **0.878**.
- Normal approximation is **not** appropriate here: np = 2 < 5. A Poisson(2)
  approximation is reasonable (n large, p small): P(X = 3) ≈ e⁻² · 2³ / 3! ≈ 0.180,
  close to the exact 0.190.

## Pitfalls

- **Sampling without replacement.** Drawing from a finite population without
  replacement makes trials dependent and changes p as items are removed, so the
  count follows the **hypergeometric** distribution, not the binomial [Wikipedia].
  The binomial is the with-replacement (or infinite-population) idealization; the
  approximation is acceptable when the sample is a small fraction (commonly < 5–10%)
  of the population.
- **Don't default to the Wald interval.** Its coverage collapses for small samples
  and extreme p; prefer Wilson or Clopper–Pearson [Wikipedia; Towards Data Science].
- **Check all four BINS assumptions, not just "two outcomes."** Streaks, learning
  effects, or a drifting p (e.g., a machine degrading over a shift) break
  independence or constant-p and invalidate the binomial.
- **Fixed n is part of the model.** If you keep sampling until you hit r successes,
  the count of failures is negative-binomial, not binomial.
- **np vs n choose k confusion.** The mean is np; C(n, k) is the coefficient inside
  the PMF. Keep them separate.
- **Wrong approximation regime.** Using the normal approximation when np < 5 (or
  n(1−p) < 5) gives poor tail probabilities; that regime calls for the exact PMF or
  the Poisson approximation instead.

## Sources

- Wikipedia, "Binomial distribution" —
  https://en.wikipedia.org/wiki/Binomial_distribution
- Wikipedia, "Binomial proportion confidence interval" —
  https://en.wikipedia.org/wiki/Binomial_proportion_confidence_interval
- Statistics LibreTexts, "3.3: Bernoulli and Binomial Distributions" —
  https://stats.libretexts.org/Courses/Saint_Mary's_College_Notre_Dame/MATH_345__-_Probability_(Kuter)/3:_Discrete_Random_Variables/3.3:_Bernoulli_and_Binomial_Distributions
- VT Pressbooks, "4.3 The Binomial Distribution" (Significant Statistics) —
  https://pressbooks.lib.vt.edu/introstatistics/chapter/binomial-distribution/
- Statistics Fundamentals, "Binomial Distribution: Formula, Examples & Step-by-Step Guide" —
  https://statisticsfundamentals.com/binomial-distribution/
- Number Analytics, "The Ultimate Binomial Distribution Guide" —
  https://www.numberanalytics.com/blog/ultimate-binomial-distribution-guide
- Probabilistic World, "The Binomial Distribution (and Theorem)" —
  https://www.probabilisticworld.com/binomial-distribution-and-theorem-intuitive-understanding/
- STAT 350, "7.4 Understanding Binomial and Poisson Distributions through CLT" —
  https://treese41528.github.io/STAT350/Website/chapter7/lectures/7-4-discret-rvs-and-clt.html
- Statology, "A Simple Explanation of Continuity Correction" —
  https://www.statology.org/continuity-correction/
- Towards Data Science, "Five Confidence Intervals for Proportions That You Should Know About" —
  https://towardsdatascience.com/five-confidence-intervals-for-proportions-that-you-should-know-about-7ff5484c024f/
