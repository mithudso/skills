<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-3-3-4-exponential` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-3-3-4-exponential
description: >-
  The exponential distribution as a named continuous probability distribution in
  probability theory — its PDF/CDF/survival function, rate (lambda) vs scale
  parameterization, mean 1/lambda and variance 1/lambda^2, the memoryless
  property, constant hazard rate, and its role as inter-arrival times of a
  Poisson process. TRIGGER: questions about modeling waiting times or times
  between events with a continuous distribution; "exponential distribution",
  "exp(lambda)", "memoryless waiting time", "time until next failure/arrival",
  "rate parameter lambda", "inter-arrival time distribution", computing
  P(X>t) / mean / variance / median / quantiles for an exponential, or fitting
  lambda by maximum likelihood. SKIP: the Poisson COUNT distribution for number
  of events (use da-1-3-3-3-poisson); the Normal/Binomial/Uniform distributions
  (use their da-1-3-3-* siblings); the general "what is a probability
  distribution" framing (use da-1-3-3-probability-distributions); exponential
  GROWTH/DECAY curves, exponential smoothing, or exponential-family GLMs (those
  are different topics, not this distribution); survival-analysis modeling
  beyond the bare distribution (Weibull/Cox).
---

# Exponential distribution

Scope: the exponential distribution as a single named continuous distribution
inside *Data Analysis > Foundations & Theory > Probability theory > Probability
distributions*. This skill covers the distribution's definition, moments,
defining properties, and the reasoning needed to recognize and apply it. It does
not cover the Poisson count distribution (its discrete sibling — see
`da-1-3-3-3-poisson`) or general distribution theory (see
`da-1-3-3-probability-distributions`).

## What it is

A continuous random variable `X` on `[0, ∞)` is **exponentially distributed**
with rate `λ > 0` if its probability density function is

    f(x; λ) = λ e^(−λx),   x ≥ 0   (and 0 for x < 0)

Written `X ~ Exp(λ)`. It models the **time until a single event** when that event
occurs at a constant average rate — time to next arrival, next failure, next
decay, next request. [Wikipedia: Exponential distribution; StatLect: Exponential
distribution]

### Core functions

| Function | Formula | Meaning |
|---|---|---|
| PDF | `f(x) = λe^(−λx)` | density at `x ≥ 0` |
| CDF | `F(x) = 1 − e^(−λx)` | `P(X ≤ x)` |
| Survival (CCDF) | `S(x) = e^(−λx)` | `P(X > x)` — "still waiting at time x" |
| Quantile | `F⁻¹(p) = −ln(1−p)/λ` | value at percentile `p` |
| Hazard rate | `h(x) = λ` (constant) | instantaneous event rate |

[Wikipedia: Exponential distribution]

### Moments and summary values

| Quantity | Value |
|---|---|
| Mean `E[X]` | `1/λ` |
| Variance `Var(X)` | `1/λ²` |
| Standard deviation | `1/λ` (equals the mean) |
| Median | `ln(2)/λ ≈ 0.693/λ` |
| Mode | `0` |
| Skewness | `2` (always right-skewed) |
| MGF | `M(t) = λ/(λ − t)` for `t < λ` |

The mean equals the standard deviation; the median is below the mean, reflecting
the right skew. [Wikipedia; StatLect: Exponential distribution]

## Two parameterizations — read carefully

This is the single most common source of error. The same distribution is written
two ways:

- **Rate** `λ`: `f(x) = λe^(−λx)`, mean `= 1/λ`. (Wikipedia, statistics, queueing.)
- **Scale** `β = 1/λ`: `f(x) = (1/β)e^(−x/β)`, mean `= β`.

`β` is itself the mean. SciPy's `scipy.stats.expon` uses **scale** (`scale = 1/λ`);
NumPy's `numpy.random.exponential` also takes **scale**; R's `rexp(n, rate)` uses
**rate**. Always confirm which one a function or formula expects before plugging
in a number. [Wikipedia: Exponential distribution]

## The defining property: memorylessness

The exponential is the **only continuous distribution that is memoryless**:

    P(X > s + t | X > s) = P(X > t),   for all s, t ≥ 0

Having already waited `s` units gives no information about the remaining wait —
the residual lifetime has the same `Exp(λ)` distribution as a fresh start. Proof
is one line using the survival function:

    P(X > s+t | X > s) = P(X > s+t)/P(X > s)
                       = e^(−λ(s+t)) / e^(−λs)
                       = e^(−λt) = P(X > t)

[StatLect: Exponential distribution; CS 547 Lecture 9, U. Wisconsin]

A direct consequence is the **constant hazard rate** `h(x) = f(x)/S(x) = λ`: the
instantaneous chance of the event is the same at every age. This makes the
exponential the natural "no aging / no wear-in" model — appropriate for
memoryless arrivals, often wrong for mechanical wear-out (where the Weibull
distribution is used instead). [Wikipedia: Exponential distribution]

## Relationship to the Poisson process (why it sits here)

If events occur in a homogeneous **Poisson process** with rate `λ`, then:

- the **number** of events in a fixed window is Poisson(`λt`) — see
  `da-1-3-3-3-poisson`;
- the **time between consecutive events** (inter-arrival time) is `Exp(λ)`;
- the time until the `k`-th event is Gamma(`k`, `λ`) — a sum of `k` independent
  `Exp(λ)` variables.

So the exponential is the continuous "waiting-time" face of the same process
whose counts are Poisson. The exponential is also a special case of the **gamma**
distribution (shape = 1) and the continuous analogue of the **geometric**
distribution (number of Bernoulli trials to first success). [LibreTexts 14.2: The
Exponential Distribution; randomservices.org: The Exponential Distribution;
Wikipedia]

## Estimating λ from data

Given an i.i.d. sample `x₁,…,xₙ`, the **maximum-likelihood estimator** is the
reciprocal of the sample mean:

    λ̂ = n / Σ xᵢ = 1 / x̄

`λ̂` is the MLE but is **biased** for `λ` (it overestimates on average); the
sample mean `x̄` is unbiased for the mean `1/λ`. For small samples the
bias-corrected estimator `(n−1)/Σxᵢ` is preferred. [Wikipedia: Exponential
distribution]

## Worked example

Support requests arrive at a constant average rate of 4 per hour, arrivals
independent (Poisson process). Model the time `X` (in hours) until the next
request as `Exp(λ)` with `λ = 4` per hour.

- Mean wait: `1/λ = 0.25 h = 15 min`.
- P(no request in the next 30 min): `S(0.5) = e^(−4·0.5) = e^(−2) ≈ 0.135`.
- P(next request within 10 min): `F(1/6) = 1 − e^(−4/6) ≈ 0.487`.
- Median wait: `ln(2)/4 ≈ 0.173 h ≈ 10.4 min` (below the 15-min mean — right skew).
- Memorylessness: if 20 min already passed with no request, the chance of waiting
  at least 30 more minutes is still `e^(−2) ≈ 0.135`, unchanged by the elapsed 20.

## Pitfalls

- **Rate vs scale mix-up.** Passing `λ` where a scale `1/λ` is expected (or vice
  versa) inverts the mean. Check the library/textbook convention every time.
  [Wikipedia]
- **Assuming a constant rate when it varies.** Real arrival rates change with time
  of day, load, seasonality. The exponential holds only over windows where `λ` is
  roughly constant; otherwise use a non-homogeneous process or piecewise rates.
  [Wikipedia: Exponential distribution]
- **Using it for wear-out / aging.** Constant hazard means "no aging." Components
  that degrade have an increasing hazard — model with Weibull or log-normal, not
  exponential.
- **Confusing it with the Poisson distribution.** Poisson counts *how many* events
  in a window (discrete); exponential measures *how long until* the next event
  (continuous). They describe the same process from different angles. See
  `da-1-3-3-3-poisson`.
- **Mean ≠ median.** Because of the right skew, reporting the mean as a "typical"
  wait understates how often waits are short; quote the median or a quantile when
  that matters.
- **Negative or zero support.** The distribution lives on `x ≥ 0`; it cannot model
  quantities that take negative values.

## Sources

1. Wikipedia, "Exponential distribution" — https://en.wikipedia.org/wiki/Exponential_distribution
2. StatLect, "Exponential distribution | Properties, proofs, exercises" — https://www.statlect.com/probability-distributions/exponential-distribution
3. Statistics LibreTexts, "14.2: The Exponential Distribution" (Siegrist, Probability, Mathematical Statistics, and Stochastic Processes) — https://stats.libretexts.org/Bookshelves/Probability_Theory/Probability_Mathematical_Statistics_and_Stochastic_Processes_(Siegrist)/14:_The_Poisson_Process/14.02:_The_Exponential_Distribution
4. Random Services (Siegrist), "The Exponential Distribution" — https://www.randomservices.org/random/poisson/Exponential.html
5. CS 547 Lecture 9, University of Wisconsin, "Conditional Probabilities and the Memoryless Property" — https://pages.cs.wisc.edu/~dsmyers/cs547/lecture_9_memoryless_property.pdf
