<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-3-6-law-of-large-numbers` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-3-6-law-of-large-numbers
description: >-
  The Law of Large Numbers (LLN) in probability theory — the result that the
  sample mean of independent draws converges to the population expected value as
  the sample size grows. Covers the weak law (convergence in probability,
  Chebyshev/Khinchin), the strong law (almost sure convergence, Kolmogorov),
  required assumptions (iid, finite mean, finite variance), the Chebyshev-inequality
  proof, convergence rate, and the misconceptions it spawns (gambler's fallacy,
  "law of averages", small-sample over-trust).
  TRIGGER: questions about why averages stabilize with more data, "law of large
  numbers", "law of averages", whether a sample mean approaches the true mean,
  how many observations are "enough", weak vs strong law, almost-sure vs
  in-probability convergence of a mean, Monte Carlo convergence justification, or
  gambler's-fallacy reasoning about independent trials.
  SKIP: the Central Limit Theorem and the distribution/spread of the sample mean
  (defer to da-1-3-7-central-limit-theorem); the definitions of expectation,
  variance, and covariance themselves (defer to da-1-3-8-expectation-variance-covariance);
  sampling distributions and standard error as an inference tool (defer to
  da-1-4-2-sampling-distributions-standard-error); the population-vs-sample
  distinction as a concept (defer to da-1-4-1-population-vs-sample); basic
  probability distributions (defer to da-1-3-3-probability-distributions).
---

# Law of Large Numbers (Probability Theory)

## What it is

The **Law of Large Numbers (LLN)** states that the average of results from a
large number of independent, identically distributed (iid) random draws converges
to the true expected value μ, provided that expected value exists. If you repeat
an experiment independently many times and average the outcomes, the average
approaches the population mean as the number of trials grows
([Statlect](https://www.statlect.com/asymptotic-theory/law-of-large-numbers);
[probabilitycourse.com 7.1.1](https://www.probabilitycourse.com/chapter7/7_1_1_law_of_large_numbers.php)).

Formally, let X₁, X₂, …, Xₙ be iid with E[Xᵢ] = μ, and let the sample mean be

    X̄ₙ = (X₁ + X₂ + … + Xₙ) / n.

The LLN says X̄ₙ → μ as n → ∞. The two versions of the law differ in *what kind*
of convergence they assert.

This is a statement about the **center** (the mean) settling down. It says nothing
about the shape or spread of X̄ₙ around μ — that is the job of the Central Limit
Theorem (out of scope here; see `da-1-3-7-central-limit-theorem`).

## Why it matters in data analysis

The LLN is the theoretical license for empirical estimation. It is the reason a
sample average is a sensible estimate of a population mean, the reason Monte Carlo
simulation works (averaging many simulated outcomes estimates an expectation), and
the reason relative frequency converges to probability (a special case where each
Xᵢ is the 0/1 indicator of an event). Without the LLN, there would be no guarantee
that collecting more data brings an estimate closer to the truth
([Wikipedia](https://en.wikipedia.org/wiki/Law_of_large_numbers);
[Britannica](https://www.britannica.com/science/probability-theory/The-strong-law-of-large-numbers)).

## The Weak Law (WLLN) — convergence in probability

The **weak law** asserts **convergence in probability**: for any ε > 0,

    lim (n→∞) P( |X̄ₙ − μ| ≥ ε ) = 0.

In words: for any tolerance ε, the probability that the sample mean is more than ε
away from μ shrinks to zero as n grows
([probabilitycourse.com 7.1.1](https://www.probabilitycourse.com/chapter7/7_1_1_law_of_large_numbers.php)).

Two classic forms:

- **Chebyshev's WLLN** — assumes finite variance σ² (and, in its general form,
  uncorrelated rather than fully independent terms). The proof is short
  (see below). It also extends to correlated sequences whose covariances average
  to zero ([Statlect](https://www.statlect.com/asymptotic-theory/law-of-large-numbers)).
- **Khinchin's WLLN** — for iid variables, finite mean alone is enough; finite
  variance is not required ([Wolfram MathWorld](https://mathworld.wolfram.com/WeakLawofLargeNumbers.html)).
  Historically the simplest case is **Bernoulli's theorem** (1713), the WLLN for
  binary trials, where the sample proportion converges to the true probability.

### Proof of the WLLN via Chebyshev's inequality (finite-variance case)

This is the proof every analyst should be able to reconstruct.

1. The sample mean has E[X̄ₙ] = μ and, for iid terms, Var(X̄ₙ) = σ²/n.
2. Chebyshev's inequality applied to X̄ₙ gives

       P( |X̄ₙ − μ| ≥ ε ) ≤ Var(X̄ₙ) / ε² = σ² / (n ε²).

3. As n → ∞, the right-hand side σ²/(nε²) → 0 for any fixed ε > 0. Hence the
   probability of a deviation of size ε vanishes
   ([probabilitycourse.com 7.1.1](https://www.probabilitycourse.com/chapter7/7_1_1_law_of_large_numbers.php)).

The key quantity is **Var(X̄ₙ) = σ²/n**: the variance of the mean falls off as 1/n,
which is also what sets the convergence rate (below).

## The Strong Law (SLLN) — almost sure convergence

The **strong law** asserts **almost sure (a.s.) convergence**:

    P( lim (n→∞) X̄ₙ = μ ) = 1.

That is, the *sequence* of sample means converges to μ along almost every infinite
realization — the set of outcomes where it fails has probability zero
([Statlect](https://www.statlect.com/asymptotic-theory/law-of-large-numbers);
[Britannica](https://www.britannica.com/science/probability-theory/The-strong-law-of-large-numbers)).

- **Kolmogorov's SLLN**: for iid variables, a **finite mean** is sufficient; finite
  variance is not required ([Statlect](https://www.statlect.com/asymptotic-theory/law-of-large-numbers)).
- The **ergodic theorem** generalizes the SLLN: it replaces iid with stationarity
  and ergodicity while keeping the finite-mean requirement
  ([Statlect](https://www.statlect.com/asymptotic-theory/law-of-large-numbers)).
- Borel's **normal number theorem** (1909) is a famous SLLN instance: for fair coin
  tosses, the long-run fraction of heads equals 1/2 with probability 1, proved with
  measure theory ([Britannica](https://www.britannica.com/science/probability-theory/The-strong-law-of-large-numbers)).

## Weak vs. Strong — the exact difference

The strong law **implies** the weak law, but not conversely; this is exactly what
the adjectives "strong" and "weak" encode
([Statlect](https://www.statlect.com/asymptotic-theory/law-of-large-numbers);
[Wikipedia](https://en.wikipedia.org/wiki/Law_of_large_numbers)).

| | Weak law (WLLN) | Strong law (SLLN) |
|---|---|---|
| Convergence mode | In probability | Almost sure |
| What converges | The probability P(\|X̄ₙ − μ\| ≥ ε) for each fixed n | The whole sample-mean sequence X̄ₙ |
| Intuition | At any large n, X̄ₙ is *probably* near μ; rare excursions can still recur | For almost every infinite run, X̄ₙ *eventually settles* at μ |
| Sufficient iid condition | Finite mean (Khinchin); finite variance (Chebyshev) | Finite mean (Kolmogorov) |

The practical distinction: the weak law allows large deviations to keep reappearing
infinitely often along a single trajectory (just with vanishing probability at each
fixed n), whereas the strong law guarantees that, for almost every trajectory, the
deviations eventually stop exceeding any fixed ε
([Britannica](https://www.britannica.com/science/probability-theory/The-strong-law-of-large-numbers)).

## Assumptions (and what breaks if they fail)

1. **A finite expected value must exist.** This is the load-bearing assumption.
   The textbook counterexample is the **Cauchy distribution**, whose mean is
   undefined: the sample mean of Cauchy draws does *not* converge — it stays
   Cauchy-distributed no matter how large n is. The LLN simply does not apply
   ([Wikipedia](https://en.wikipedia.org/wiki/Law_of_large_numbers)).
2. **Identically distributed** (or stationary, for the ergodic generalization).
   If the distribution drifts over time, X̄ₙ tracks a moving target, not a single μ.
3. **Independence** (or at least vanishing average covariance). The Chebyshev form
   tolerates uncorrelated or weakly correlated terms; strong positive dependence
   can prevent the variance of the mean from shrinking like σ²/n.
4. **Finite variance is needed for the simple Chebyshev proof**, but is *not*
   required for the law to hold under iid (Khinchin/Kolmogorov cover finite-mean,
   infinite-variance cases) ([Statlect](https://www.statlect.com/asymptotic-theory/law-of-large-numbers)).

## Convergence rate (how fast)

The variance of the sample mean is σ²/n, so its standard deviation is σ/√n. The
typical size of |X̄ₙ − μ| therefore shrinks like **1/√n**, not 1/n
([probabilitycourse.com 7.1.1](https://www.probabilitycourse.com/chapter7/7_1_1_law_of_large_numbers.php)).
Consequences for analysts:

- To **halve** the typical error you need **4×** the data.
- Convergence is slow in the tail: precise estimates of small probabilities or
  heavy-tailed means can require very large n.
- This 1/√n rate is the bridge to the Central Limit Theorem, which gives the full
  approximate distribution of the deviation — but that quantification of *spread*
  belongs to `da-1-3-7-central-limit-theorem`, not here.

## Common pitfalls and misconceptions

- **Gambler's fallacy / "law of averages."** The LLN does *not* say that past
  outcomes get "corrected." After a streak of heads, the next fair toss is still
  50/50. The long-run average converges because later trials *swamp* the early
  imbalance by sheer count, not because the process compensates for it
  ([Statistics How To](https://www.statisticshowto.com/law-large-numbers/);
  [Statistics By Jim](https://statisticsbyjim.com/basics/gamblers-fallacy/)).
- **"It guarantees I'll hit the expected value."** It does not. For any finite n
  the sample mean may never equal μ; the law is an asymptotic statement about the
  limit, not about any particular trial or finite run
  ([Statistics How To](https://www.statisticshowto.com/law-large-numbers/)).
- **Applying it to small samples.** The LLN is a *large-n* result. Over-trusting an
  average from a handful of observations is the "law of small numbers" error.
- **Confusing it with the CLT.** The LLN says *where* the average goes (to μ); the
  CLT says *how it is distributed* around μ on the way. Keep them separate.
- **Ignoring the existence-of-mean precondition.** With heavy tails (e.g. Cauchy,
  or some power-law data), the mean may not exist and averaging is meaningless —
  check before invoking the LLN.

## Worked sketch: relative frequency → probability

Let Xᵢ = 1 if trial i is a "success" (probability p) and 0 otherwise. Then
E[Xᵢ] = p and the sample mean X̄ₙ is the observed success fraction. By the LLN,
the observed fraction converges to the true probability p — the weak form is
Bernoulli's theorem, the strong form is covered by Borel's normal-number result
([Wikipedia](https://en.wikipedia.org/wiki/Law_of_large_numbers);
[Britannica](https://www.britannica.com/science/probability-theory/The-strong-law-of-large-numbers)).
This is precisely why simulating an event many times and taking the hit-rate gives
a good probability estimate, and why Monte Carlo averages estimate expectations.

## Sources

- Statlect — *Law of Large Numbers (strong and weak, with proofs)*:
  https://www.statlect.com/asymptotic-theory/law-of-large-numbers
- ProbabilityCourse.com — *7.1.1 Law of Large Numbers*:
  https://www.probabilitycourse.com/chapter7/7_1_1_law_of_large_numbers.php
- Encyclopaedia Britannica — *Probability theory: The strong law of large numbers*:
  https://www.britannica.com/science/probability-theory/The-strong-law-of-large-numbers
- Wikipedia — *Law of large numbers*: https://en.wikipedia.org/wiki/Law_of_large_numbers
- Wolfram MathWorld — *Weak Law of Large Numbers*:
  https://mathworld.wolfram.com/WeakLawofLargeNumbers.html
- Statistics How To — *Law of Large Numbers / Law of Averages*:
  https://www.statisticshowto.com/law-large-numbers/
- Statistics By Jim — *Gambler's Fallacy*: https://statisticsbyjim.com/basics/gamblers-fallacy/
