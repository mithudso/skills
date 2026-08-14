<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-3-3-5-uniform` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-3-3-5-uniform
description: >-
  The uniform distribution as a named member of the probability-distribution
  family in Data Analysis foundations: the continuous Uniform(a,b) and its
  discrete counterpart, their PMF/PDF, CDF, mean, variance, entropy, and the
  reasons this distribution anchors random-number generation and the
  probability integral transform.
  TRIGGER: a question names or describes a "uniform distribution", "U(a,b)",
  "Uniform(0,1)", "rectangular distribution", "equally likely outcomes over a
  range", "constant density"; deriving or recalling the uniform PDF/PMF, CDF,
  mean (a+b)/2, or variance (b-a)^2/12; why a flat/no-information prior or
  baseline is uniform; the uniform as the maximum-entropy distribution on a
  bounded support; the standard uniform behind inverse-transform sampling and
  the probability integral transform.
  SKIP: random-number-generator engineering, PRNG algorithms, seeding, or
  Monte-Carlo tooling (sampling implementation, not the distribution's theory);
  uniform priors framed as a Bayesian-inference topic rather than a named
  distribution; other named distributions, each owned by its sibling skill
  (da-1-3-3-1-normal, da-1-3-3-2-binomial, da-1-3-3-3-poisson,
  da-1-3-3-4-exponential); the general taxonomy of "probability distributions"
  as a category (da-1-3-3-probability-distributions); the PMF-vs-PDF distinction
  as its own concept (da-1-3-2); random variables in the abstract (da-1-3-1).
---

# Uniform distribution (Data Analysis > Foundations > Probability theory > Probability distributions)

The uniform distribution is the formal statement of "no value in the range is
favored over any other." Every point (continuous) or outcome (discrete) on the
support carries equal probability weight. It is the simplest non-trivial
distribution and the reference point against which other distributions are
described as "more peaked," "skewed," or "heavy-tailed."

This skill covers the uniform **as a named distribution** in the
foundations layer: its functional form, its summary moments, and the two roles
that make it foundational — the maximum-entropy distribution on a bounded
support, and the standard uniform U(0,1) that sits underneath the probability
integral transform and inverse-transform sampling. It does not cover PRNG
implementation or Bayesian prior selection (see SKIP).

## Two forms: continuous and discrete

A uniform distribution exists in both flavors. Identify which one the data calls
for before writing any formula — the support type decides PDF vs PMF.

| | Continuous Uniform(a,b) | Discrete Uniform on {x_1,...,x_N} |
|---|---|---|
| Support | the interval [a,b] | N equally likely values, often {1,...,N} or {a,...,b} |
| Density/mass | f(x) = 1/(b-a) on [a,b], else 0 | P(X = x_i) = 1/N |
| Use case | any real value in a window is equally plausible | a fair die, a random index, a card draw |

Continuous: "a sensor reading lands anywhere between 2.0 and 2.5 V with no
preferred spot." Discrete: "a fair six-sided die." Mixing the two — for example
writing a 1/(b-a) density for a die — is the most common modeling error.

## Continuous Uniform(a,b) — the core object

Notation: X ~ U(a,b) or Uniform(a,b), with -inf < a < b < inf.
[Source: Wikipedia, "Continuous uniform distribution";
https://en.wikipedia.org/wiki/Continuous_uniform_distribution]

**Probability density function (PDF)**

    f(x) = 1/(b-a)   for a <= x <= b
    f(x) = 0         otherwise

The density is a flat rectangle of height 1/(b-a) and width (b-a); its area is 1,
which is the only constraint the normalization imposes. The flat shape is why
the distribution is also called the **rectangular distribution**.
[Source: probabilitycourse.com, "Uniform Distribution";
https://www.probabilitycourse.com/chapter4/4_2_1_uniform.php]

**Cumulative distribution function (CDF)**

    F(x) = 0              for x < a
    F(x) = (x - a)/(b - a) for a <= x <= b
    F(x) = 1              for x > b

The CDF rises in a straight line from 0 to 1 across [a,b] — a ramp, not a curve.
[Source: Wikipedia; https://en.wikipedia.org/wiki/Continuous_uniform_distribution]

**Quantile / inverse CDF**

    F^{-1}(p) = a + p*(b - a)   for 0 < p < 1

This linear inverse is exactly what makes the uniform the workhorse of
inverse-transform sampling (see "Why it is foundational").

**Summary measures**

| Quantity | Formula | Note |
|---|---|---|
| Mean E[X] | (a + b)/2 | the midpoint of the interval |
| Median | (a + b)/2 | equals the mean (symmetric) |
| Mode | any value in (a,b) | density is flat, so no unique mode |
| Variance Var(X) | (b - a)^2 / 12 | grows with the square of the width |
| E[X^2] | (a^2 + ab + b^2)/3 | used to derive the variance |
| Skewness | 0 | perfectly symmetric |
| Excess kurtosis | -6/5 = -1.2 | platykurtic (flatter than normal) |
| Differential entropy | ln(b - a) | maximal for fixed support |

[Sources: Wikipedia; probabilitycourse.com;
University College Dublin lecture notes, "Statistics: Uniform Distribution
(Continuous)", https://www.ucd.ie/msc/t4media/Uniform%20Distribution.pdf]

**Variance derivation** (worth keeping in mind — it explains the /12):

    Var(X) = E[X^2] - (E[X])^2
           = (a^2 + ab + b^2)/3 - ((a + b)/2)^2
           = (b - a)^2 / 12

[Source: probabilitycourse.com;
https://www.probabilitycourse.com/chapter4/4_2_1_uniform.php]

**Moment-generating function**

    M_X(t) = (e^{tb} - e^{ta}) / (t*(b - a))   for t != 0,  and 1 for t = 0

[Source: Wikipedia; https://en.wikipedia.org/wiki/Continuous_uniform_distribution]

## Discrete Uniform — the equally-likely-outcomes form

For N equally likely integer outcomes a, a+1, ..., b (so N = b - a + 1):

    PMF:  P(X = k) = 1/N         for k in {a, ..., b}
    Mean: E[X] = (a + b)/2
    Var:  Var(X) = (N^2 - 1)/12

For the canonical {1,...,N} case the mean is (N+1)/2 and the variance is
(N^2 - 1)/12. Note the parallel /12 with the continuous case but the (N^2 - 1)
numerator — discreteness changes the spread by a small correction.
[Sources: GeeksforGeeks, "Uniform Distribution Formula";
https://www.geeksforgeeks.org/maths/uniform-distribution-formula/ ;
VrcAcademy, "Uniform Distribution";
https://vrcacademy.com/tutorials/uniform-distribution/]

## Why it is foundational

**1. Maximum entropy on a bounded support.** Among all continuous distributions
confined to [a,b], the uniform has the largest differential entropy, ln(b-a).
In plain terms: if all you know is the range, the uniform is the least-assuming
description of the variable — it adds no information beyond the bounds. This is
why a "flat" or "uninformative" baseline over a finite range is uniform.
[Sources: Wikipedia, "Maximum entropy probability distribution",
https://en.wikipedia.org/wiki/Maximum_entropy_probability_distribution ;
The Book of Statistical Proofs, "Continuous uniform distribution maximizes
differential entropy for fixed range",
https://statproofbook.github.io/P/cuni-maxent.html]

**2. The probability integral transform (PIT).** For any continuous random
variable X with CDF F, the transformed variable Y = F(X) follows U(0,1). The
converse — feeding U(0,1) values through F^{-1} — produces samples from F. The
standard uniform U(0,1) is therefore the universal "raw material" of simulation:
draw u ~ U(0,1), return x = F^{-1}(u). The linear inverse CDF above is why U(0,1)
is the natural starting point.
[Source: Taylor & Francis, "Probability integral transform";
https://taylorandfrancis.com/knowledge/Engineering_and_technology/Engineering_support_and_special_topics/Probability_integral_transform/]

## Worked micro-examples

- **U(0,1) standard uniform.** Mean 1/2, variance 1/12 ~ 0.0833, entropy
  ln(1) = 0. P(X <= 0.3) = 0.3 directly from the CDF.
- **Sensor in [2.0, 2.5] V, continuous.** Density = 1/0.5 = 2 per volt.
  P(2.1 <= X <= 2.3) = (2.3 - 2.1)/(2.5 - 2.0) = 0.2/0.5 = 0.4.
  Mean = 2.25 V, variance = 0.5^2/12 ~ 0.0208.
- **Fair six-sided die, discrete.** P(X = k) = 1/6 for k in {1..6}.
  Mean = 3.5, variance = (6^2 - 1)/12 = 35/12 ~ 2.917. Do NOT use (b-a)^2/12 here.

## Pitfalls

- **Continuous/discrete mix-up.** A 1/(b-a) density on a finite set of outcomes,
  or a 1/N mass on an interval, is wrong. Check whether the support is an
  interval or a list before choosing PDF vs PMF.
- **Off-by-one in the discrete count.** N = b - a + 1, not b - a, when the
  support is {a,...,b}. The "+1" is the most frequent discrete-uniform slip.
- **Endpoints and zero-probability points.** For the continuous uniform,
  P(X = a) = 0 and the strict/inclusive endpoints (< vs <=) do not change any
  probability; do not treat the boundary as a special atom.
- **Reading the density as a probability.** The continuous density 1/(b-a) can
  exceed 1 (e.g., on [0,0.5] it is 2). A density is not a probability;
  probabilities come from integrating it over an interval.
- **Assuming uniformity for convenience.** "Equally likely" is a strong claim.
  Real arrival times, measurement errors, and counts are usually not uniform.
  The uniform's flat-ness should be justified by the range-only knowledge state,
  not chosen because it is the simplest formula.

## Sources

1. Wikipedia, "Continuous uniform distribution" —
   https://en.wikipedia.org/wiki/Continuous_uniform_distribution
2. probabilitycourse.com (Pishro-Nik, "Introduction to Probability"),
   "Uniform Distribution" —
   https://www.probabilitycourse.com/chapter4/4_2_1_uniform.php
3. University College Dublin, "Statistics: Uniform Distribution (Continuous)" —
   https://www.ucd.ie/msc/t4media/Uniform%20Distribution.pdf
4. The Book of Statistical Proofs, "Continuous uniform distribution maximizes
   differential entropy for fixed range" —
   https://statproofbook.github.io/P/cuni-maxent.html
5. Taylor & Francis, "Probability integral transform" —
   https://taylorandfrancis.com/knowledge/Engineering_and_technology/Engineering_support_and_special_topics/Probability_integral_transform/
6. GeeksforGeeks, "Uniform Distribution | Continuous and Discrete Formula" —
   https://www.geeksforgeeks.org/maths/uniform-distribution-formula/
