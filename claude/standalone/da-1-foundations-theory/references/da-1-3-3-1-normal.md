<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-3-3-1-normal` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-3-3-1-normal
description: >-
  The Normal (Gaussian) distribution as a named continuous probability
  distribution in probability theory for data analysis: its PDF/CDF, the two
  parameters (mean mu, variance sigma^2), the standard normal, z-score
  standardization, the 68-95-99.7 empirical rule, symmetry/skewness/kurtosis,
  closure under linear combinations, the Central Limit Theorem connection, and
  how to check the normality assumption with a Q-Q plot.
  TRIGGER: a data-analysis task names the normal/Gaussian/bell-curve
  distribution; needs the normal PDF or CDF formula; converts a value to a
  z-score or back; applies the 68-95-99.7 rule; reads a standard normal table;
  asks whether data are "normally distributed" or how to test/visualize that;
  models measurement error or a sampling distribution as normal.
  SKIP: a different named distribution (binomial -> da-1-3-3-2-binomial,
  Poisson -> da-1-3-3-3-poisson, exponential -> da-1-3-3-4-exponential,
  uniform -> da-1-3-3-5-uniform); the general idea of a probability
  distribution -> da-1-3-3-probability-distributions; what PMFs/PDFs are in
  general -> da-1-3-2-probability-mass-density-functions; random variables as a
  concept -> da-1-3-1-random-variables; the normal distribution as it appears
  elsewhere in the taxonomy (e.g. inferential-statistics normality assumptions,
  regression error terms, or hypothesis-test sampling distributions) -> defer
  to that branch's own skill.
---

# Normal (Gaussian) distribution

The single most-used continuous distribution in statistics. This skill scopes
it as a **named probability distribution within probability theory** — what it
is, its formulas, its properties, and how to recognize and check it. Inferential
procedures that *rely* on normality (t-tests, regression error assumptions,
confidence intervals) live in other taxonomy branches; this skill supplies the
distribution itself, not the test.

## Definition

A random variable `X` is normally distributed with mean `mu` and variance
`sigma^2`, written `X ~ N(mu, sigma^2)`, when its probability density function
(PDF) is:

```
f(x) = (1 / (sigma * sqrt(2*pi))) * exp( -(x - mu)^2 / (2 * sigma^2) )
```

for `x` in `(-inf, +inf)`. [NIST e-Handbook 1.3.6.6.1; Wikipedia: Normal distribution]

Two parameters fully define it:

- **`mu`** — the **location** parameter: the center of the bell. It equals the
  mean, the median, and the mode (the distribution is unimodal and symmetric).
- **`sigma`** — the **scale** parameter (`sigma > 0`): the standard deviation,
  controlling spread/width. `sigma^2` is the variance. [NIST e-Handbook]

The curve is the familiar symmetric "bell": single peak at `x = mu`, tails that
approach but never touch zero in both directions.

## The standard normal

Set `mu = 0` and `sigma = 1` to get the **standard normal**, conventionally
written `Z`, with density:

```
phi(z) = (1 / sqrt(2*pi)) * exp( -z^2 / 2 )
```

[NIST e-Handbook; Wikipedia]

Any normal variable maps to the standard normal by **standardization**
(z-scoring):

```
z = (x - mu) / sigma
```

A z-score is "how many standard deviations `x` lies from the mean": positive =
above the mean, negative = below. [Statistics LibreTexts 5.05; OpenStax 6.1]
Standardization is what lets a single standard normal table (the "z-table")
serve every normal distribution — there is no need for a separate table per
`(mu, sigma)`. [NIST PMC 5.1; Standard normal table, Wikipedia]

## CDF and finding probabilities

The cumulative distribution function (CDF) gives `P(X <= x)`. For the standard
normal it is denoted `Phi(z)`. It has **no closed form in elementary
functions**; it is expressed through the error function:

```
Phi(z) = (1/2) * [ 1 + erf( z / sqrt(2) ) ]
```

[Wikipedia: Normal distribution]

To find a probability for a general normal, standardize first:

```
P(X <= v) = Phi( (v - mu) / sigma )
```

[NIST PMC 5.1]

Worked example. `Z ~ N(0,1)`. The probability of landing between `z = -1` and
`z = +1` is `Phi(1) - Phi(-1) = 0.8413 - 0.1587 = 0.6826`. That 68.26% is
exactly the first tier of the empirical rule below. [NIST PMC 5.1]

Worked example (general normal). Male 500 m speed-skating times are
`N(mu = 70.42 s, sigma = 0.34 s)`. To be faster than 95% of competitors a skater
must finish in the bottom 5% of times, at z about `-1.645`:
`x = mu + z*sigma = 70.42 + (-1.645)(0.34) = 69.86 s`. [worked example via
Statistics LibreTexts / Krista King Math]

## The 68-95-99.7 (empirical) rule

For any normal distribution, the proportion of values within `k` standard
deviations of the mean is fixed:

| Interval            | Proportion |
|---------------------|------------|
| `mu +/- 1*sigma`    | 68.27%     |
| `mu +/- 2*sigma`    | 95.45%     |
| `mu +/- 3*sigma`    | 99.73%     |

[NIST e-Handbook; Statistics LibreTexts 5.05]

Practical rounding ("68 / 95 / 99.7") is the everyday version. Common
critical-value shortcuts: about 90% of mass falls within `+/- 1.645*sigma`, 95%
within `+/- 1.960*sigma`, and 99% within `+/- 2.576*sigma` — the two-sided
quantiles used for percentiles and intervals.

## Moments and shape

| Statistic        | Value     |
|------------------|-----------|
| Mean             | `mu`      |
| Median           | `mu`      |
| Mode             | `mu`      |
| Std. deviation   | `sigma`   |
| Skewness         | `0`       |
| Kurtosis         | `3`       |
| Excess kurtosis  | `0`       |
| Range            | `(-inf, +inf)` |

[NIST e-Handbook; Wikipedia]

Skewness `0` reflects perfect symmetry. Kurtosis `3` is the reference point for
"mesokurtic"; excess kurtosis (kurtosis minus 3) is therefore `0`, which is why
the normal is the baseline against which other distributions are called heavy-
or light-tailed. The moment generating function is
`M(t) = exp(mu*t + sigma^2 * t^2 / 2)`. [Wikipedia]

## Key structural properties (why it dominates)

- **Closed under linear transformation.** If `X ~ N(mu, sigma^2)` then
  `aX + b ~ N(a*mu + b, a^2 * sigma^2)`. Standardization is the special case
  `a = 1/sigma`, `b = -mu/sigma`.
- **Closed under addition of independents.** If `X1 ~ N(mu1, sigma1^2)` and
  `X2 ~ N(mu2, sigma2^2)` are independent, then
  `X1 + X2 ~ N(mu1 + mu2, sigma1^2 + sigma2^2)`. Variances add; means add. This
  is what makes the normal the natural model for propagating uncertainty. [Wikipedia]
- **Central Limit Theorem (CLT).** The average (or sum) of many independent
  samples from *any* distribution with finite mean and variance converges to a
  normal distribution as sample size grows. This is the deep reason the normal
  appears everywhere — sampling distributions of means, measurement error that
  is the sum of many small effects, aggregate noise. [NIST e-Handbook; Wikipedia]

In data analysis the normal is the default model for measurement error and for
the error/residual term in classical regression and ANOVA, and it underpins the
significance levels of classical hypothesis tests. [NIST e-Handbook] (Those
procedures themselves belong to the inferential-statistics branch — see SKIP.)

## Checking the normality assumption

Do not assume normality; check it. The primary visual tool is the **normal
probability plot (Q-Q plot)**: plot the ordered data against the quantiles of a
theoretical normal. If the data are approximately normal the points fall on an
approximate straight line; **departures from the line indicate departures from
normality**. [NIST e-Handbook 1.3.3.21]

Reading the deviation tells you *how* it is non-normal:

- An **S-curve / systematic bend** -> skewness (right- or left-skewed data).
- Points **curving away at both ends, above the line on the right and below on
  the left** -> heavier tails than normal (fat tails / outliers).
- The **opposite curvature** -> shorter (lighter) tails than normal.

[NIST e-Handbook] The correlation coefficient of the plotted points can be
compared to critical-value tables for a formal normality test; complementary
formal tests include Shapiro-Wilk and Anderson-Darling.

## Common pitfalls

- **Assuming normality without checking.** Real data are often skewed or
  heavy-tailed (incomes, latencies, financial returns). Run a Q-Q plot or a
  formal test first; an exponential, log-normal, or Poisson model may fit better
  (see the sibling distribution skills).
- **Confusing the empirical rule with exact tails.** 68/95/99.7 is a quick
  mental check, not a substitute for `Phi`. For precise probabilities or
  percentiles, standardize and use the CDF / z-table.
- **Sign and direction errors in z-scores.** `z = (x - mu)/sigma`. A value below
  the mean gives a negative z. To go from a percentile back to a value, invert:
  `x = mu + z*sigma`.
- **`sigma` vs `sigma^2`.** The PDF and `N(mu, sigma^2)` notation use the
  *variance*; the empirical rule and z-scores use the *standard deviation*.
  Mixing them is a frequent and silent error.
- **Treating a large sample as automatically normal.** The CLT makes the
  *sampling distribution of the mean* approximately normal; it does not make the
  *raw data* normal. A histogram of skewed raw data stays skewed no matter how
  many points you collect.
- **Probability of an exact point is zero.** Because the normal is continuous,
  `P(X = c) = 0`; probabilities are areas over intervals, read from the CDF.
- **Tails never reach zero.** The support is all real numbers, so a normal model
  technically allows impossible values (e.g. negative heights). For
  strictly-positive quantities consider a truncated normal or a log-normal.

## Quick reference

```
PDF:            f(x) = (1/(sigma*sqrt(2*pi))) * exp(-(x-mu)^2 / (2*sigma^2))
Standard PDF:   phi(z) = (1/sqrt(2*pi)) * exp(-z^2/2)
Standardize:    z = (x - mu) / sigma          (invert: x = mu + z*sigma)
CDF:            P(X <= v) = Phi((v - mu)/sigma);  Phi(z) = 0.5*(1 + erf(z/sqrt(2)))
Empirical rule: 1sigma~68.27%  2sigma~95.45%  3sigma~99.73%
Critical z:     90% +/-1.645   95% +/-1.960   99% +/-2.576
Moments:        mean=median=mode=mu, var=sigma^2, skew=0, kurtosis=3
Linear combo:   aX+b ~ N(a*mu+b, a^2*sigma^2)
Sum (indep):    X1+X2 ~ N(mu1+mu2, sigma1^2+sigma2^2)
```

## Sources

- NIST/SEMATECH e-Handbook of Statistical Methods, 1.3.6.6.1 "Normal
  Distribution": https://www.itl.nist.gov/div898/handbook/eda/section3/eda3661.htm
- NIST/SEMATECH e-Handbook, 1.3.3.21 "Normal Probability Plot":
  https://www.itl.nist.gov/div898/handbook/eda/section3/normprpl.htm
- NIST/SEMATECH e-Handbook, PMC 5.1 (standardizing to a z-score, standard normal
  CDF): https://www.itl.nist.gov/div898/handbook/pmc/section5/pmc51.htm
- Wikipedia, "Normal distribution" (PDF/CDF, error function, moments, MGF,
  closure properties, CLT): https://en.wikipedia.org/wiki/Normal_distribution
- Statistics LibreTexts 5.05, "The Empirical Rule and Standard Normal (Z)
  Distribution": https://stats.libretexts.org/Courses/Red_Rocks_Community_College/Introduction_to_Statistics_(RRCC)/05:_Probability_Distributions/5.05:_The_Empirical_Rule_and_Standard_Normal_(Z)_Distribution
- OpenStax, Introductory Statistics 2e, 6.1 "The Standard Normal Distribution":
  https://openstax.org/books/introductory-statistics-2e/pages/6-1-the-standard-normal-distribution
