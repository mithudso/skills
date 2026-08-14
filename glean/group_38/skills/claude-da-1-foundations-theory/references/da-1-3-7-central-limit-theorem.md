<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-3-7-central-limit-theorem` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-3-7-central-limit-theorem
description: >
  The central limit theorem (CLT) as a foundation of probability theory: why the
  sampling distribution of the sample mean (or sum) approaches a normal distribution
  as sample size grows, the i.i.d. and finite-variance conditions, the standard
  error sigma/sqrt(n), and the practical n>=30 guidance.
  TRIGGER: questions about "why is the sample mean normally distributed", "central
  limit theorem", "CLT", "sampling distribution of the mean approaches normal",
  "how large a sample for normality", "standard error of the mean", "n>=30 rule",
  or why normal-based methods work even on non-normal data.
  SKIP: the law of large numbers / convergence of the mean to mu itself (defer to
  da-1-3-6-law-of-large-numbers); the normal distribution's own properties as a
  distribution (defer to da-1-3-3-1-normal); the mechanics of sampling distributions
  and standard error as an inference topic (defer to
  da-1-4-2-sampling-distributions-standard-error); expectation/variance definitions
  (defer to da-1-3-8-expectation-variance-covariance); CLT framings under any
  non-probability-theory parent in the taxonomy.
---

# Central Limit Theorem (CLT)

Scope: this skill covers the CLT as a result *within probability theory*
(Data Analysis > Foundations & Theory > Probability theory). It explains the
theorem, its conditions, and its consequences for the sample mean. It defers the
*inference machinery* built on top of it (confidence intervals, hypothesis tests)
and the standalone properties of the normal distribution to their own skills.

## What the theorem says

Take a sequence of independent and identically distributed (i.i.d.) random
variables X1, X2, ..., Xn with finite mean E[Xi] = mu and finite, positive
variance Var(Xi) = sigma^2. The **classical Lindeberg-Levy CLT** states that the
standardized sample mean converges in distribution to a standard normal as
n grows without bound [Wikipedia, Central limit theorem]:

    sqrt(n) * (Xbar_n - mu)  -->d  N(0, sigma^2)

equivalently, the CDFs satisfy

    lim (n -> inf) P[ sqrt(n)(Xbar_n - mu) <= z ] = Phi(z / sigma)

where Xbar_n = (1/n) * sum(Xi), Phi is the standard normal CDF, and "-->d" means
convergence in distribution. The key point: this holds **even if the Xi
themselves are not normal** [Wikipedia; Statistics By Jim].

In applied notation, for a sample of size n the sampling distribution of the mean
is approximately [Statistics LibreTexts, 7.2]:

    Xbar  ~approx~  N( mu , sigma / sqrt(n) )

so its z-score is

    z = (xbar - mu) / (sigma / sqrt(n))

The denominator **sigma / sqrt(n)** is the **standard error of the mean** — the
standard deviation of the sampling distribution, not of the population.

CLT for sums: the sample sum behaves as Sum(X) ~approx~ N(n*mu, sqrt(n)*sigma)
[Statistics LibreTexts, 7.2].

## Why it matters

- It explains why normal-based methods (z-tests, t-tests, normal confidence
  intervals for a mean) give valid answers even when the underlying data are
  skewed or otherwise non-normal, *provided the sample is large enough*
  [Statistics By Jim].
- It is the reason the normal distribution shows up so often: any quantity that is
  effectively a sum or average of many small independent contributions tends
  toward normality.
- It refines the law of large numbers. The LLN says Xbar_n converges in
  probability to mu; the CLT describes the *size and shape of the fluctuations*
  around mu — they shrink like 1/sqrt(n) and are approximately normal
  [Wikipedia]. Convergence vs. fluctuation is the dividing line between these two
  theorems.

## Conditions (and what breaks them)

The classical i.i.d. form requires [Wikipedia; Statology; Scribbr]:

1. **Independent, identically distributed** observations — in practice, met by
   random sampling. Convenience or clustered samples violate this and the
   sampling distribution will not behave as the CLT predicts.
2. **Finite variance** (sigma^2 < inf). Heavy-tailed distributions with infinite
   variance (e.g. a Cauchy distribution) do **not** obey the classical CLT; the
   sample mean of Cauchy variates stays Cauchy regardless of n.
3. **Adequate sample size** — see the n>=30 discussion below.
4. **10% condition** for sampling without replacement: keep n <= 10% of the
   population so observations stay approximately independent [Statology].

Generalizations relax the "identically distributed" part:

- **Lindeberg condition** — for independent but non-identical Xi with
  s_n^2 = sum(sigma_i^2), the CLT holds if for every eps > 0,
  (1/s_n^2) * sum E[(Xi - mu_i)^2 * 1{|Xi - mu_i| > eps*s_n}] -> 0. No single term
  may dominate the total variance [Wikipedia].
- **Lyapunov CLT** — a simpler sufficient condition using a (2+delta) moment:
  (1/s_n^(2+delta)) * sum E[|Xi - mu_i|^(2+delta)] -> 0 for some delta > 0
  [Wikipedia].

## Rate of convergence

The **Berry-Esseen theorem** bounds how fast the approximation kicks in: for
i.i.d. variables with a finite third moment, the maximum gap between the true CDF
of the standardized mean and the normal CDF shrinks at rate O(1/sqrt(n))
[Wikipedia]. The constant depends on the skewness/third moment, which is why
skewed populations need larger n.

## The n >= 30 rule of thumb — use with judgment

A sample size of about 30 is *often* enough for the sampling distribution of the
mean to look normal, but this is a rule of thumb, not a law [Statistics By Jim]:

- **Moderately skewed** populations: n ~ 20-40 usually gives an approximately
  normal sampling distribution.
- **Severely skewed** populations: even n = 80 can still leave visible skew in the
  sampling distribution of the mean. The more the population departs from normal,
  the larger n must be.
- If the population is already normal, the sample mean is *exactly* normal for any
  n (no CLT needed).

## Worked example

Population mean mu = 90, population standard deviation sigma = 15, sample size
n = 25. Find P(85 < Xbar < 92) [Statistics LibreTexts, 7.2].

1. Sampling distribution of the mean: Xbar ~ N(90, 15/sqrt(25)) = N(90, 3). The
   standard error is 15/5 = 3.
2. Standardize the endpoints:
   z_low  = (85 - 90)/3 = -1.667
   z_high = (92 - 90)/3 =  0.667
3. Probability: P(-1.667 < Z < 0.667) = Phi(0.667) - Phi(-1.667)
   ~= 0.7475 - 0.0478 ~= 0.6997.

So there is roughly a 70% chance the sample mean of 25 observations lands between
85 and 92 — far tighter than the spread of individual values, because the
standard error (3) is much smaller than sigma (15).

## Pitfalls

- **Sampling distribution != sample distribution.** The CLT describes the
  distribution of the *mean across many samples*, not the histogram of one sample.
  A single sample of skewed data still looks skewed [Statistics By Jim].
- **The population does not change.** The CLT does not normalize individual
  observations or the population itself — only the distribution of the mean
  (or sum) [Statistics By Jim].
- **Do not use it for individual predictions.** Standard error sigma/sqrt(n)
  governs the mean; an individual observation still has spread sigma.
- **Independence is non-negotiable.** No sample size rescues dependent or
  non-random data [Statistics By Jim].
- **Finite variance required.** Heavy-tailed, infinite-variance populations break
  the classical CLT entirely.
- **n>=30 is not magic.** Always weigh it against the population's skew; check the
  shape of the data first.

## Sources

- Wikipedia, "Central limit theorem" — Lindeberg-Levy statement, Lindeberg and
  Lyapunov conditions, Berry-Esseen rate, relation to the law of large numbers.
  https://en.wikipedia.org/wiki/Central_limit_theorem
- Statistics By Jim, "Central Limit Theorem Explained" — applied guidance, n>=30
  caveats, skewness effects, standard error, common misconceptions.
  https://statisticsbyjim.com/basics/central-limit-theorem/
- Statistics LibreTexts, "7.2: The Central Limit Theorem for Sample Means" —
  Xbar ~ N(mu, sigma/sqrt(n)), z-score for the mean, worked numeric example, CLT
  for sums.
  https://stats.libretexts.org/Courses/Los_Angeles_City_College/Introductory_Statistics/07:_The_Central_Limit_Theorem/7.02:_The_Central_Limit_Theorem_for_Sample_Means_(Averages)
- Statology, "Central Limit Theorem: The Four Conditions to Meet" — i.i.d.,
  randomization, 10% condition, finite variance.
  https://www.statology.org/central-limit-theorem-conditions/
- Scribbr, "Central Limit Theorem | Formula, Definition & Examples" — definition
  and sampling-distribution framing.
  https://www.scribbr.com/statistics/central-limit-theorem/
