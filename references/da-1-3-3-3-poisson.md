<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-3-3-3-poisson` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-3-3-3-poisson
description: >-
  The Poisson distribution as a named probability distribution in data-analysis
  foundations — the discrete distribution of the count of independent events in a
  fixed interval of time or space, governed by a single rate parameter lambda
  where mean = variance = lambda. Covers the PMF, the four Poisson-process
  assumptions, mean/variance/equidispersion, the Poisson limit (approximating the
  binomial), MLE and confidence intervals for lambda, and the overdispersion
  pitfall. TRIGGER when a data-analysis task asks: what is the Poisson
  distribution, its PMF/mean/variance, when count data is Poisson, how to compute
  P(X=k) for an event rate, whether n large + p small justifies a Poisson
  approximation to the binomial, or how to spot equidispersion vs overdispersion.
  SKIP when the work is Poisson *regression* / GLM model fitting (a modeling
  technique, not this distribution-definition node), the exponential distribution
  for inter-arrival times (use da-1-3-3-4-exponential), the binomial distribution
  itself (use da-1-3-3-2-binomial), the normal distribution (use
  da-1-3-3-1-normal), the general "what is a probability distribution" parent
  (use da-1-3-3-probability-distributions), or PMF/PDF mechanics in the abstract
  (use da-1-3-2-probability-mass-density-functions).
---

# Poisson distribution

Taxonomy location: Data Analysis > Foundations & Theory > Probability theory >
Probability distributions > Poisson.

This skill is scoped to the Poisson distribution **as a named member of the
probability-distribution catalog** — its definition, PMF, parameters, moments,
defining assumptions, and the limit/approximation relationships an analyst needs
to decide whether count data is Poisson. It is not a guide to Poisson regression
(a GLM modeling technique that lives under regression/modeling, not under this
distribution-definition node), nor to the exponential distribution that describes
the waiting time between Poisson events.

## Definition

The Poisson distribution is a **discrete** probability distribution that gives the
probability of a given number of events occurring in a fixed interval of time or
space, when those events happen at a known constant average rate and independently
of the time since the last event [Scribbr; DataCamp; Statistics By Jim]. It is
defined by a single parameter, **lambda (λ)** — the mean number of occurrences per
interval [GeeksforGeeks; DataCamp].

A random variable X taking values k = 0, 1, 2, ... is Poisson-distributed when:

  P(X = k) = (e^(−λ) · λ^k) / k!

This is the probability mass function (PMF) [GeeksforGeeks; DataCamp; Scribbr].
The three pieces have an intuitive reading: λ^k carries the likelihood driven by
the average rate, k! corrects for the number of arrangements of k indistinguishable
events, and e^(−λ) is the normalizing factor that makes the probabilities sum to 1
[DataCamp].

It satisfies the two PMF axioms: P(X = k) ≥ 0 for every k, and the sum over all
k ≥ 0 equals 1 [GeeksforGeeks].

Notation: X ~ Poisson(λ), with λ > 0.

## Parameter, mean, and variance

The single parameter λ is simultaneously the **mean** and the **variance**:

  E[X] = λ   and   Var(X) = λ

so the standard deviation is √λ [GeeksforGeeks; DataCamp; Statistics By Jim]. This
equality, **equidispersion**, is the signature property of the Poisson and the
single most useful diagnostic: if observed count data has sample variance roughly
equal to its sample mean, Poisson is plausible; if not, it is suspect (see
Pitfalls) [DataCamp; SAS KB 56549].

Because mean and variance both equal λ, variability grows with the expected count —
a process averaging 100 events per interval has standard deviation 10, while one
averaging 4 events has standard deviation 2 [DataCamp].

## The defining assumptions (the Poisson process)

A counting process produces Poisson counts when these hold [DataCamp; Scribbr]:

1. **Independence.** Each event occurs independently of others; one occurrence does
   not change the probability of another.
2. **Constant rate.** The average rate λ is constant over the interval (homogeneity).
3. **No simultaneity.** Two events cannot occur at exactly the same instant; events
   are counted singly.
4. **Fixed, well-defined interval.** The interval of time, area, or volume is fixed
   and specified in advance. A common error is letting the interval vary.

Under these assumptions the count in any interval is Poisson(λ) and counts in
disjoint intervals are independent [DataCamp; GeeksforGeeks]. This counting process
is the **Poisson process**; the times *between* successive events in it follow the
exponential distribution (covered by da-1-3-3-4-exponential) [DataCamp].

## When to use it

Reach for the Poisson when you are counting how many times something happens in a
fixed window and the four assumptions are credible [Statistics By Jim; DataCamp]:

- Queueing / arrivals: customers per hour, calls to a call center, packets on a
  link, requests to a server.
- Rare-event counts: defects per batch, typos per page, accidents per intersection
  per month, mutations per genome region.
- Epidemiology: occurrences of a rare disease, outbreak-detection counts.
- Spatial counts: organisms per quadrat, stars per sky region.

Contrast with the binomial (da-1-3-3-2-binomial): the binomial counts successes in
a **fixed number n of trials** with success probability p; the Poisson counts
events in a continuous interval with **no fixed upper bound** on the count
[DataCamp].

## Worked example

A help desk receives on average λ = 3 tickets per hour. Probability of exactly 5
tickets in the next hour:

  P(X = 5) = e^(−3) · 3^5 / 5! = 0.0498 · 243 / 120 ≈ 0.1008, about 10.1%.

Probability of **at most 1** ticket:

  P(X ≤ 1) = P(0) + P(1) = e^(−3) + e^(−3)·3 = e^(−3)(1 + 3) = 0.0498 · 4 ≈ 0.199.

Scaling the interval scales λ: over a 20-minute window the rate is λ = 3 · (1/3) = 1,
so P(0 in 20 min) = e^(−1) ≈ 0.368. Rate scales linearly with interval length only
because the process is assumed homogeneous (assumption 2).

## Key relationships

- **Poisson limit of the binomial (law of rare events).** As n → ∞ and p → 0 with
  np → λ fixed, Binomial(n, p) → Poisson(λ) [Poisson limit theorem, Wikipedia].
  Practical rule of thumb for using Poisson as a binomial approximation: large n,
  small p, with common thresholds n ≥ 20 and p ≤ 0.05 (some texts use n ≥ 100 and
  np ≤ 10) [Solon Karapanagiotis; SOGA-Py; BathMASH HELM]. This is the standard
  shortcut when the binomial's factorials are unwieldy and successes are rare.
- **Additivity.** If X1 ~ Poisson(λ1) and X2 ~ Poisson(λ2) are independent, then
  X1 + X2 ~ Poisson(λ1 + λ2). Summing independent Poisson counts keeps them Poisson.
- **Normal approximation.** For large λ (rule of thumb λ ≳ 10), Poisson(λ) is
  approximately Normal(λ, λ), useful for tail probabilities [Save My Exams].
- **Exponential link.** Inter-arrival times in a Poisson process are
  Exponential(λ); the two distributions describe the same process from the count
  side vs. the waiting-time side (see da-1-3-3-4-exponential) [DataCamp].

## Estimating lambda from data

- **MLE.** Given counts x1, ..., xn, the maximum-likelihood estimate of λ is the
  sample mean, λ̂ = x̄. This is also the method-of-moments estimate.
- **Confidence interval.** An exact interval uses the chi-square / gamma
  relationship; for large counts the Wald interval λ̂ ± z·√(λ̂/n) (or per-interval
  λ̂ ± z·√λ̂) is a serviceable approximation that follows from the normal
  approximation above.
- **Goodness of fit.** Compare observed count frequencies against the Poisson PMF
  with a chi-square goodness-of-fit test, and check the variance-to-mean ratio
  (the dispersion index, ≈ 1 under Poisson) [SAS KB 56549].

## Pitfalls

- **Overdispersion.** Real count data often has variance **greater** than the mean,
  violating equidispersion. Forcing a Poisson model then underestimates standard
  errors and overstates significance, yielding invalid conclusions; the usual fix
  is the negative binomial (a gamma-mixed Poisson) or a quasi-Poisson adjustment
  [Overdispersion, Wikipedia; SAS KB 56549; DataCamp]. Underdispersion (variance <
  mean) is rarer but also breaks the Poisson assumption.
- **Non-constant rate.** If λ drifts across the interval (rush hours, seasonality),
  the homogeneous Poisson is wrong; use a time-varying / non-homogeneous rate.
- **Dependence and clustering.** Events that arrive in bursts or trigger one another
  violate independence and inflate variance — another route to overdispersion.
- **Varying interval length.** Counts collected over unequal intervals are not
  identically distributed; normalize by exposure (an offset in modeling terms) or
  rescale λ proportionally to interval length.
- **Confusing count with rate.** λ is a per-interval expected count, not a
  probability; it can exceed 1. Do not treat it like a Bernoulli p.
- **Zero inflation.** Excess zeros beyond what Poisson predicts (structural zeros)
  call for zero-inflated models rather than a plain Poisson.

## Sources

1. Statistics By Jim — "Poisson Distribution: Definition & Uses."
   https://statisticsbyjim.com/probability/poisson-distribution/
2. DataCamp — "Poisson Distribution: A Comprehensive Guide."
   https://www.datacamp.com/tutorial/poisson-distribution
3. GeeksforGeeks — "Poisson Distribution | Formula, Table, Mean and Variance."
   https://www.geeksforgeeks.org/maths/poisson-distribution/
4. Scribbr — "Poisson Distributions | Definition, Formula & Examples."
   https://www.scribbr.com/statistics/poisson-distribution/
5. Poisson limit theorem (Poisson approximation to binomial) — Wikipedia.
   https://en.wikipedia.org/wiki/Poisson_limit_theorem
6. Overdispersion — Wikipedia.
   https://en.wikipedia.org/wiki/Overdispersion
7. SAS Knowledge Base 56549 — "Models for overdispersed and underdispersed count
   data." http://support.sas.com/kb/56/549.html
8. SOGA-Py, FU Berlin — "Poisson Approximation to the Binomial Distribution."
   https://www.geo.fu-berlin.de/en/v/soga-py/Basics-of-statistics/Discrete-Random-Variables/The-Poisson-Distribution/Poisson-Approximation-to-the-Binomial-Distribution/index.html
