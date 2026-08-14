<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-4-2-sampling-distributions-standard-error` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-4-2-sampling-distributions-standard-error
description: >-
  The sampling distribution of a statistic (mean or proportion) and the standard
  error that measures its spread, as a foundation of statistical inference. Covers
  the three-distribution mental model (population, single sample, sampling
  distribution), the mean and standard error of x-bar and p-hat, why standard
  error shrinks as 1/sqrt(n), the finite population correction, and the practical
  shift from a known-parameter formula to an estimated one.
  TRIGGER: questions about what a sampling distribution is; the standard error of
  the mean (SEM = s/sqrt(n)) or of a proportion (sqrt(p(1-p)/n)); why standard
  error differs from standard deviation; how sample size changes the precision of
  an estimate; "standard deviation of the sample mean"; the distribution of x-bar
  or p-hat; setting up the spread term before a confidence interval or test.
  SKIP: the Central Limit Theorem as a theorem in its own right (limiting behavior,
  conditions, proof sketch) -> da-1-3-7-central-limit-theorem; the Law of Large
  Numbers convergence of the mean -> da-1-3-6-law-of-large-numbers; expectation,
  variance and covariance of a single random variable -> da-1-3-8; population vs
  sample as a conceptual contrast -> da-1-4-1; building or interpreting confidence
  intervals or hypothesis tests that consume the standard error (later inference
  nodes). Defer those framings rather than absorbing them here.
---

# Sampling distributions & standard error

Scope: Data Analysis > Foundations & Theory > Statistical inference foundations >
Sampling distributions & standard error. This skill explains the object that makes
inference possible — the distribution of a sample statistic — and the single number,
the standard error, that summarizes how much that statistic moves from sample to
sample. It deliberately stops at the door of the Central Limit Theorem, the Law of
Large Numbers, and confidence-interval construction; those are separate nodes (see
SKIP cues above).

## 1. The core idea in one sentence

A **sampling distribution** is the probability distribution of a *statistic*
(such as the sample mean x-bar or the sample proportion p-hat) computed over all
possible samples of a fixed size n drawn from the same population. The **standard
error** is the standard deviation of that sampling distribution — it measures how
far a typical sample statistic falls from the true population parameter
[Statistics By Jim; LibreTexts].

The key move is conceptual: you stop treating the sample mean as a fixed number and
start treating it as a *random variable* with its own distribution. Everything in
inference — confidence intervals, p-values, margins of error — is built on the
spread of that distribution.

## 2. The three-distribution model (do not confuse these)

Most errors with this topic come from conflating three different distributions.
Keep them separate:

| Distribution | What varies | Center | Spread |
|---|---|---|---|
| **Population distribution** | individual units in the population | mu | sigma (population SD) |
| **Sample distribution** | the n observed data points in one sample | x-bar | s (sample SD) |
| **Sampling distribution** | the statistic across all possible samples of size n | mu (for x-bar) | **standard error** |

- The **population** and a **single sample** are distributions of *individual
  values*. Their spread is a standard deviation.
- The **sampling distribution** is a distribution of a *summary statistic*. Its
  spread is the **standard error** — which is itself a standard deviation, but a
  standard deviation *of a statistic*, not of raw data [Statistics By Jim].

This is why "the standard deviation of the sample mean" and "the standard error of
the mean" name the same quantity.

## 3. Sampling distribution of the mean (x-bar)

For samples of size n from a population with mean mu and standard deviation sigma:

- **Mean (center):** E[x-bar] = mu. The sampling distribution of the mean is
  *centered on the population mean* — the sample mean is an unbiased estimator
  [PSU STAT; 365 Data Science].
- **Standard error (spread):**

  ```
  SE(x-bar) = sigma / sqrt(n)        (sigma known)
  SEM       = s     / sqrt(n)        (sigma unknown — estimated from the sample)
  ```

  where s is the sample standard deviation [Statistics By Jim].

In practice sigma is almost never known, so analysts use **s / sqrt(n)**. This is
the "standard error of the mean" (SEM) reported by software. Using s instead of
sigma is what later motivates the t-distribution in interval/test nodes (out of
scope here).

### Worked detail
A population has sigma = 20. With n = 25, SE(x-bar) = 20 / sqrt(25) = 20/5 = 4.
With n = 100, SE(x-bar) = 20/10 = 2. Quadrupling n halved the standard error.
With n = 400, SE = 20/20 = 1 — another quadrupling, another halving.

## 4. Sampling distribution of a proportion (p-hat)

For a binary outcome with true population proportion p and sample size n:

- **Mean (center):** E[p-hat] = p — the sample proportion is unbiased for p
  [LibreTexts; runestone].
- **Standard error (spread):**

  ```
  SE(p-hat) = sqrt( p(1-p) / n )         (p known / assumed)
  se(p-hat) = sqrt( p-hat(1-p-hat) / n ) (p unknown — estimated from the sample)
  ```

  [LibreTexts; SJSU StatPrimer].

- **Maximum spread at p = 0.5:** the product p(1-p) is largest at p = 0.5, so the
  standard error of a proportion is greatest when the true proportion is one half,
  for a fixed n [statisticshowto].

### Worked detail
With p = 0.5 and n = 100: SE = sqrt(0.5 * 0.5 / 100) = sqrt(0.0025) = 0.05.
With p = 0.1 and n = 100: SE = sqrt(0.1 * 0.9 / 100) = sqrt(0.0009) = 0.03.
Same n, smaller spread away from 0.5.

## 5. Why standard error shrinks as 1 / sqrt(n)

The sqrt(n) in the denominator is the single most important property:

- As n grows, the standard error **decreases**, so sample statistics cluster more
  tightly around the true parameter — estimates get more *precise* [Statistics By
  Jim].
- The relationship is **inverse-square-root, not linear**: to halve the standard
  error you must **quadruple** the sample size. This is the law of diminishing
  returns for sample size — going from n=100 to n=400 buys the same precision gain
  as going from n=400 to n=1600 [Statistics By Jim].

This square-root scaling is a property of the *sampling distribution itself*; it
holds regardless of whether you invoke the Central Limit Theorem for shape.

## 6. Finite population correction (FPC)

The standard formulas assume sampling **with replacement** or from an effectively
infinite population. When you sample a meaningful fraction of a *finite* population
without replacement, multiply the standard error by the FPC factor:

```
FPC = sqrt( (N - n) / (N - 1) )      N = population size, n = sample size
```

- Rule of thumb: apply the FPC when the sampling fraction n/N exceeds about 5%.
- When n is tiny relative to N, FPC -> 1 and the correction is negligible — which
  is why introductory treatments omit it.
- The correction *reduces* the standard error: sampling a large share of a finite
  population leaves less uncertainty than the uncorrected formula implies.

## 7. Common pitfalls

1. **Standard error is not standard deviation.** SD describes the spread of *raw
   data*; SE describes the spread of a *statistic*. Reporting SE when you mean SD
   (or vice versa) misstates variability — SE is almost always the smaller number
   because of the sqrt(n) division [Statistics By Jim].
2. **The sampling distribution is a thought construct, not a dataset.** You almost
   never draw "all possible samples." You compute SE from one sample using the
   formula; the sampling distribution is the theoretical object the formula
   describes.
3. **Known vs estimated parameter.** SE(x-bar) = sigma/sqrt(n) needs the *true*
   sigma; in practice you use s/sqrt(n). For proportions, you swap p for p-hat. The
   estimated version introduces extra uncertainty that downstream methods (t, score
   intervals) account for — not here [SJSU StatPrimer].
4. **A bigger sample shrinks the standard error, not the population SD.** Increasing
   n does nothing to sigma or s (the data spread); it only tightens the sampling
   distribution. Confusing the two leads people to expect the histogram of their
   data to narrow with more sampling — it does not.
5. **Proportion SE peaks at p = 0.5.** Designing a survey assuming the worst case
   (p = 0.5) gives a conservative, safe sample size — useful when p is unknown in
   advance [statisticshowto].
6. **Do not import CLT claims here.** Whether the sampling distribution is
   *approximately normal* (and under what n / np conditions) is the Central Limit
   Theorem's job — defer shape questions to da-1-3-7. This skill covers the
   distribution's *center and spread* only.

## 8. Quick reference

```
Sampling distribution = distribution of a statistic over all samples of size n
Standard error (SE)    = standard deviation of that sampling distribution

Mean:        center = mu,  SE = sigma/sqrt(n)  (or s/sqrt(n) when sigma unknown)
Proportion:  center = p,   SE = sqrt(p(1-p)/n) (or sqrt(p-hat(1-p-hat)/n))

Precision scales as 1/sqrt(n): 4x the sample size -> half the SE
FPC (finite pop, no replacement): SE * sqrt((N-n)/(N-1)), apply when n/N > ~5%
SE < SD almost always; they answer different questions
```

## Sources

- Statistics By Jim — *Standard Error of the Mean (SEM)*:
  https://www.statisticsbyjim.com/hypothesis-testing/standard-error-mean/
  (SEM definition, s/sqrt(n) formula, SE vs SD table, sqrt(n) scaling).
- Statistics LibreTexts — *The Sampling Distribution of Sample Proportions*:
  https://stats.libretexts.org/Bookshelves/Introductory_Statistics/Introductory_Statistics_(Hannah_Seidler-Wright)/06%3A_Inference_Involving_a_Single_Population_Proportion/6.01%3A_The_Sampling_Distribution_of_Sample_Proportions
  (mean = p, SE = sqrt(p(1-p)/n), normality conditions, worked example).
- San Jose State University StatPrimer (Gerstman) — *Standard Error of a Proportion*:
  https://www2.sjsu.edu/faculty/gerstman/StatPrimer/conf-prop.htm
  (known-p vs estimated p-hat standard error formulas).
- 365 Data Science — *Central Limit Theorem and Standard Error*:
  https://365datascience.com/tutorials/statistics-tutorials/central-limit-theorem/
  (SE as the standard deviation of the sampling distribution; centering on mu).
- Runestone Academy / Introductory Statistics — *Sampling distribution of a sample
  proportion*: https://runestone.academy/ns/books/published/ahss3rd/distributionphat.html
  (sample proportion as unbiased estimator of p).
- Statistics How To — *Sampling Distribution of the Sample Proportion*:
  https://www.statisticshowto.com/sampling-distribution-of-the-sample-proportion/
  (variability largest at p = 0.5).
