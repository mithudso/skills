<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-3-probability-theory` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-3-probability-theory
description: >-
  Foundational probability theory as the mathematical basis for data analysis and
  statistical inference: the Kolmogorov axioms, sample spaces and events, conditional
  probability, independence, the law of total probability, Bayes' theorem, random
  variables, expectation and variance, the law of large numbers, and the central
  limit theorem. TRIGGER: questions about what probability means, how to set up a
  probability model, computing or interpreting conditional probabilities, applying
  Bayes' theorem (including base-rate / diagnostic-test reasoning), why sample means
  behave the way they do, or the formal justification that underlies estimation,
  confidence intervals, and hypothesis tests. SKIP: specifying or fitting a named
  distribution family for modeling (defer to a distributions skill); designing or
  analyzing a statistical test, p-value, or confidence interval procedure (defer to
  an inference skill); random variables treated as their own deep topic with PMFs,
  PDFs, and CDFs (defer to da-1-3-1-random-variables); Bayesian modeling workflows,
  priors, and MCMC as a methodology (defer to a Bayesian-methods skill); combinatorics
  or counting as a standalone topic.
---

# Probability Theory (Foundations for Data Analysis)

This skill covers the probability layer that sits *underneath* statistics. In the
taxonomy this is `Data Analysis > Foundations & Theory > Probability theory`. The
goal is to model uncertainty rigorously enough that inference (estimation, testing,
prediction) rests on solid ground. Scope here is the **core machinery**: the axioms,
how to reason about events, and the two limit theorems that make sampling work.

Where a topic has its own node — random variables and their distributions as a deep
subject, or the design of inference procedures — this skill states the foundational
idea and defers the depth there.

## 1. The axiomatic foundation (Kolmogorov, 1933)

A probability model has three parts: a **sample space** Ω (the set of all possible
outcomes of an experiment), a collection of **events** (subsets of Ω, technically a
σ-algebra), and a **probability function** P that assigns a number to each event
([Probability axioms — Wikipedia](https://en.wikipedia.org/wiki/Probability_axioms)).

Kolmogorov's three axioms ([Probability axioms — Wikipedia](https://en.wikipedia.org/wiki/Probability_axioms)):

1. **Non-negativity:** P(A) ≥ 0 for every event A.
2. **Normalization:** P(Ω) = 1 — some outcome in the sample space always occurs.
3. **Countable additivity:** for any countable sequence of mutually exclusive
   (disjoint) events E₁, E₂, E₃, …, P(E₁ ∪ E₂ ∪ …) = Σ P(Eᵢ).

Everything else is derived, not assumed. Useful consequences:

- P(∅) = 0; P(Aᶜ) = 1 − P(A); 0 ≤ P(A) ≤ 1.
- Monotonicity: if A ⊆ B then P(A) ≤ P(B).
- Inclusion–exclusion: P(A ∪ B) = P(A) + P(B) − P(A ∩ B).

**Why a data analyst cares:** every downstream statistical claim ("the probability
this estimate is off by more than x is below 5%") is meaningful only because P obeys
these rules. The axioms are what stop probability statements from being arbitrary.

## 2. Conditional probability

The probability of A *given that* B has occurred
([MIT 18.05 Class 3 — Conditional Probability, Independence, Bayes](https://ocw.mit.edu/courses/18-05-introduction-to-probability-and-statistics-spring-2022/mit18_05_s22_class03-prep.pdf)):

```
P(A | B) = P(A ∩ B) / P(B),   provided P(B) > 0
```

Conditioning on B effectively shrinks the sample space to B and re-normalizes. The
**multiplication rule** is the same identity rearranged:

```
P(A ∩ B) = P(A | B) · P(B) = P(B | A) · P(A)
```

This is the workhorse for building joint probabilities from a sequence of dependent
steps (e.g., drawing without replacement).

## 3. Independence

A and B are **independent** when knowing one tells you nothing about the other
([MIT 18.05 Class 3](https://ocw.mit.edu/courses/18-05-introduction-to-probability-and-statistics-spring-2022/mit18_05_s22_class03-prep.pdf)):

```
P(A ∩ B) = P(A) · P(B)   ⇔   P(A | B) = P(A)
```

Independence is an *assumption you must justify*, not a default. The i.i.d.
(independent and identically distributed) assumption underpins most introductory
inference, and it is also the most common place a real analysis quietly breaks
(clustered data, time series, repeated measures on the same subject). Note the
distinction: **mutually exclusive** events (P(A ∩ B) = 0) are the opposite of
independent — if one happens the other cannot.

## 4. Law of total probability

If events B₁, …, Bₙ **partition** the sample space (mutually exclusive, and together
cover Ω), then for any event A
([Statistics LibreTexts — Conditional Probability and Bayes' Rule](https://stats.libretexts.org/Courses/Saint_Mary's_College_Notre_Dame/MATH_345__-_Probability_(Kuter)/2:_Computing_Probabilities/2.2:_Conditional_Probability_and_Bayes'_Rule)):

```
P(A) = Σᵢ P(A | Bᵢ) · P(Bᵢ)
```

This lets you compute an unconditional probability by splitting the world into cases
and averaging the conditional probabilities, weighted by how likely each case is. It
is the denominator that makes Bayes' theorem computable.

## 5. Bayes' theorem

The rule for inverting a conditional — updating a belief about A after observing B
([MIT 18.05 Class 3](https://ocw.mit.edu/courses/18-05-introduction-to-probability-and-statistics-spring-2022/mit18_05_s22_class03-prep.pdf)):

```
P(A | B) = P(B | A) · P(A) / P(B)
```

Expanded with the law of total probability in the denominator:

```
P(A | B) = P(B | A) · P(A) / Σᵢ P(B | Aᵢ) · P(Aᵢ)
```

Read it as: **posterior ∝ likelihood × prior**. P(A) is the prior, P(B | A) the
likelihood, P(A | B) the posterior.

### Worked example — diagnostic test (base-rate sensitivity)

Disease prevalence P(D) = 0.001. Test sensitivity P(+ | D) = 0.99. False-positive
rate P(+ | Dᶜ) = 0.05.

```
P(+) = P(+|D)P(D) + P(+|Dᶜ)P(Dᶜ)
     = 0.99·0.001 + 0.05·0.999 = 0.00099 + 0.04995 = 0.05094
P(D | +) = 0.00099 / 0.05094 ≈ 0.0194  (about 1.9%)
```

Even with a 99%-sensitive test, a positive result implies under 2% chance of disease
because the disease is rare. This is the canonical demonstration of why the prior
cannot be discarded.

## 6. Common pitfalls

- **Confusing P(A | B) with P(B | A)** — the "prosecutor's fallacy." These are
  generally different and describe opposite causal directions
  ([MIT 18.05 Class 3](https://ocw.mit.edu/courses/18-05-introduction-to-probability-and-statistics-spring-2022/mit18_05_s22_class03-prep.pdf)).
- **Base-rate neglect** — interpreting P(B | A) while ignoring the prior P(A); the
  diagnostic example above shows how badly this misleads
  ([MIT 18.05 Class 3](https://ocw.mit.edu/courses/18-05-introduction-to-probability-and-statistics-spring-2022/mit18_05_s22_class03-prep.pdf)).
- **Assuming independence to multiply** — multiplying marginal probabilities for
  events that are actually correlated inflates or deflates joint risk.
- **Conditioning on a zero-probability event** — P(A | B) is undefined when P(B) = 0;
  watch for this in continuous models.

## 7. Random variables and expectation (foundation only)

A **random variable** maps outcomes in Ω to numbers, letting you do arithmetic with
uncertainty. The deep treatment — discrete vs. continuous, PMFs, PDFs, CDFs — belongs
to `da-1-3-1-random-variables`; here is the part the limit theorems need.

**Expectation** (the probability-weighted average / "center of mass")
([Expected value — Wikipedia](https://en.wikipedia.org/wiki/Expected_value)):

```
discrete:    E[X] = Σ xᵢ · pᵢ
continuous:  E[X] = ∫ x · f(x) dx
```

**Linearity of expectation** holds with or without independence
([Expected value — Wikipedia](https://en.wikipedia.org/wiki/Expected_value)):

```
E[X + Y] = E[X] + E[Y]      E[aX] = a · E[X]
```

**Variance** measures spread around the mean: Var[X] = E[(X − E[X])²]
([Expected value — Wikipedia](https://en.wikipedia.org/wiki/Expected_value)).

## 8. The two limit theorems that make statistics work

These are the bridge from "probability model" to "I can learn from a sample."

### Law of Large Numbers (LLN)

As the sample size n grows, the sample mean X̄ₙ converges to the true expected value μ
([Expected value — Wikipedia](https://en.wikipedia.org/wiki/Expected_value); [Foundations of Statistics — CLT & LLN](https://bookdown.org/peter_neal/math4081_notes/Sec_CLT.html)).
The strong law states this convergence happens *almost surely*. This is the formal
reason an estimate based on more data is, on average, closer to the truth — and the
justification for using a sample mean to estimate a population mean.

### Central Limit Theorem (CLT) — Lindeberg–Lévy form

For i.i.d. X₁, X₂, … with mean μ and **finite** variance σ², the standardized sample
mean converges *in distribution* to a normal
([Central limit theorem — Wikipedia](https://en.wikipedia.org/wiki/Central_limit_theorem); [Foundations of Statistics — CLT & LLN](https://bookdown.org/peter_neal/math4081_notes/Sec_CLT.html)):

```
√n · (X̄ₙ − μ)  →  N(0, σ²)   as n → ∞
```

This is why so many inference procedures (z-tests, t-tests, normal-approximation
confidence intervals) work even when the underlying data are not normal: it is the
*distribution of the mean* that becomes normal, not the data.

Convergence rate is governed by the **Berry–Esseen theorem** — error in the normal
approximation shrinks at a rate of about 1/√n when the third moment is finite
([Central limit theorem — Wikipedia](https://en.wikipedia.org/wiki/Central_limit_theorem)).

### CLT misconceptions to avoid

([Central limit theorem — Wikipedia](https://en.wikipedia.org/wiki/Central_limit_theorem))

- It applies to the distribution of **sample means**, not to a single random sample
  of raw values.
- It does **not** make the population distribution normal; large n does not normalize
  the data themselves.
- The **"n ≥ 30" rule of thumb has no general mathematical basis** — heavily skewed
  or heavy-tailed populations can need far larger n, and infinite-variance
  populations break the theorem entirely.

## How this connects downstream

- LLN → consistency of estimators.
- CLT → the sampling distributions behind standard errors, confidence intervals, and
  test statistics — handled in the inference skills, not here.
- Bayes' theorem → the formal core of Bayesian methods — the modeling *workflow*
  (priors, posteriors, MCMC) is a separate skill.

## Sources

1. [Probability axioms — Wikipedia](https://en.wikipedia.org/wiki/Probability_axioms) — Kolmogorov's three axioms, sample space, events, derived consequences.
2. [MIT OCW 18.05, Class 3 — Conditional Probability, Independence, Bayes' Theorem](https://ocw.mit.edu/courses/18-05-introduction-to-probability-and-statistics-spring-2022/mit18_05_s22_class03-prep.pdf) — conditional probability, multiplication rule, independence, Bayes, base-rate / inverse-conditional pitfalls.
3. [Statistics LibreTexts — Conditional Probability and Bayes' Rule](https://stats.libretexts.org/Courses/Saint_Mary's_College_Notre_Dame/MATH_345__-_Probability_(Kuter)/2:_Computing_Probabilities/2.2:_Conditional_Probability_and_Bayes'_Rule) — law of total probability and partition reasoning.
4. [Expected value — Wikipedia](https://en.wikipedia.org/wiki/Expected_value) — expectation (discrete/continuous), linearity, variance, link to the law of large numbers.
5. [Central limit theorem — Wikipedia](https://en.wikipedia.org/wiki/Central_limit_theorem) — Lindeberg–Lévy statement, Berry–Esseen rate, common misconceptions.
6. [Foundations of Statistics — Central Limit Theorem and Law of Large Numbers (P. Neal)](https://bookdown.org/peter_neal/math4081_notes/Sec_CLT.html) — formal LLN and CLT in a foundations-of-statistics text.
