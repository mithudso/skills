<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-3-5-bayes-theorem` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-3-5-bayes-theorem
description: >-
  Bayes' theorem as a result of probability theory within Foundations & Theory of
  Data Analysis — the formula that inverts a conditional probability, P(A|B) =
  P(B|A)P(A)/P(B), its derivation from the definition of conditional probability,
  the law-of-total-probability denominator, the prior/likelihood/posterior/evidence
  vocabulary, odds form, and the base-rate-fallacy pitfall in diagnostic testing.
  TRIGGER: a question asks how to invert a conditional probability, update a belief
  with new evidence, compute P(cause|effect) or P(disease|positive test), explain
  prior/likelihood/posterior/evidence as terms in the theorem, work a base-rate /
  false-positive-paradox example, or derive/state the theorem itself. SKIP: building
  full Bayesian parameter-estimation or inference workflows comparing statistical
  paradigms (defer to da-1-4-3-frequentist-vs-bayesian-paradigms and
  da-1-4-statistical-inference-foundations); the underlying P(A|B) conditional and
  joint/marginal definitions themselves (defer to
  da-1-3-4-joint-marginal-conditional-probability); the naive Bayes classifier or
  Bayesian networks as ML models (defer to the relevant modeling skill).
---

# Bayes' Theorem (Probability Theory)

Scope: Bayes' theorem as a theorem of probability theory under
**Data Analysis > Foundations & Theory > Probability theory**. This skill covers
the statement, derivation, vocabulary, forms, and the interpretation pitfalls of
the theorem. It does not cover the broader Bayesian-vs-frequentist inference debate,
parameter estimation, or machine-learning classifiers built on top of it — see the
SKIP cues in the frontmatter.

## Statement

For two events A and B with P(B) > 0:

```
P(A | B) = P(B | A) · P(A) / P(B)
```

In words: the conditional probability of A given B equals the probability of B given
A, scaled by the unconditional probability of A and divided by the unconditional
probability of B. The theorem lets you "flip" a conditional — turn a known P(B|A)
into the often-wanted P(A|B) [Wikipedia: Bayes' theorem; GeeksforGeeks: Bayes
Theorem].

## Derivation

Start from the definition of conditional probability (this prerequisite lives in
da-1-3-4-joint-marginal-conditional-probability):

```
P(A | B) = P(A ∩ B) / P(B)        P(B | A) = P(A ∩ B) / P(A)
```

Both reference the same joint probability P(A ∩ B). Rearranging each:

```
P(A ∩ B) = P(A | B) · P(B) = P(B | A) · P(A)
```

Divide the last equality by P(B) to isolate P(A | B):

```
P(A | B) = P(B | A) · P(A) / P(B)
```

That is Bayes' theorem [Wikipedia: Bayes' theorem; GeeksforGeeks: Bayes Theorem]. The
theorem is named for Thomas Bayes; the modern formulation is due to Pierre-Simon
Laplace.

## The denominator: law of total probability

The marginal P(B) is usually not given directly. Expand it with the law of total
probability over a partition of the sample space. For a hypothesis A and its
complement Aᶜ:

```
P(B) = P(B | A) · P(A) + P(B | Aᶜ) · P(Aᶜ)
```

So the full single-hypothesis form is:

```
                 P(B | A) · P(A)
P(A | B) = ---------------------------------------
            P(B | A) · P(A) + P(B | Aᶜ) · P(Aᶜ)
```

For a partition into mutually exclusive, exhaustive events E₁, …, Eₙ:

```
                 P(B | Eᵢ) · P(Eᵢ)
P(Eᵢ | B) = ----------------------------
              Σⱼ P(B | Eⱼ) · P(Eⱼ)
```

[Lumen Learning: Bayes' Theorem; GeeksforGeeks: Bayes Theorem]. Computing the
denominator correctly is the step most people skip — it is what forces every
hypothesis's contribution into the normalization.

## Vocabulary (the four named pieces)

Write the theorem as P(H | D) = P(D | H) · P(H) / P(D), with H a hypothesis/belief
and D the observed data/evidence:

| Term | Symbol | Meaning |
|---|---|---|
| **Prior** | P(H) | belief in the hypothesis before seeing the data |
| **Likelihood** | P(D \| H) | probability of the data if the hypothesis is true |
| **Posterior** | P(H \| D) | updated belief after seeing the data |
| **Evidence** (marginal likelihood / normalizing constant) | P(D) | total probability of the data across all hypotheses |

[bookdown — Bayesian Statistics the Fun Way: Prior, Likelihood, Posterior;
ScienceDirect: Bayes Theorem overview]. The evidence P(D) is a normalizing constant:
it does not depend on which hypothesis you are evaluating, only on the data, and its
job is to make the posterior a valid probability in [0, 1].

## Useful forms

**Proportionality form.** Because P(D) is the same for every competing hypothesis,
it cancels when you only need to compare or rank hypotheses:

```
posterior ∝ likelihood × prior        P(H | D) ∝ P(D | H) · P(H)
```

You can drop the normalizing constant, compute the unnormalized scores, then divide
each by their sum to renormalize [bookdown — Bayesian Statistics the Fun Way].

**Odds form.** Dividing the theorem for hypothesis H by the theorem for its
complement Hᶜ cancels P(D):

```
P(H | D)       P(D | H)     P(H)
--------  =   ---------- ·  -----
P(Hᶜ | D)      P(D | Hᶜ)    P(Hᶜ)

posterior odds = likelihood ratio × prior odds
```

The likelihood ratio P(D|H)/P(D|Hᶜ) is the factor by which the data shift the odds
[Wikipedia: Bayes' theorem].

## Worked example 1 — diagnostic test (base rates matter)

Disease prevalence 2%, test sensitivity 90% (false-negative rate 10%), false-positive
rate 1%. Find P(disease | positive test) [Lumen Learning: Bayes' Theorem].

Natural-frequency method (imagine 10,000 people):

- Have disease: 10,000 × 0.02 = 200 → test positive: 200 × 0.90 = **180**
- No disease: 10,000 × 0.98 = 9,800 → test positive: 9,800 × 0.01 = **98**
- Total positives: 180 + 98 = 278
- Posterior: 180 / 278 ≈ **0.647 (about 65%)**

Formula check:

```
P(D⁺ | +) = (0.02)(0.90) / [(0.02)(0.90) + (0.98)(0.01)]
          = 0.018 / 0.0278 ≈ 0.647
```

## Worked example 2 — the credible-witness flip

A man tells the truth 3 of 4 times and reports a die showed a six. Find P(six |
he reports six) [GeeksforGeeks: Bayes Theorem].

```
P(six) = 1/6,  P(not six) = 5/6
P(reports six | six) = 3/4,  P(reports six | not six) = 1/4

P(six | reports) = (3/4 · 1/6) / [(3/4 · 1/6) + (1/4 · 5/6)]
                 = (1/8) / (1/8 + 5/24) = 3/8
```

Even an honest-3-of-4 witness only makes "six" 3/8 likely, because sixes are rare to
begin with — the prior pulls the posterior down.

## Pitfalls

- **Base-rate fallacy / false-positive paradox.** Ignoring the prior (prevalence)
  and reading the test's accuracy as the answer is the classic error. When the base
  rate is below the false-positive rate, most positives are false — *overall* — even
  if any single test is rarely wrong. Example: population of 1,000, 2% infected, 5%
  false-positive rate → 20 true positives vs 49 false positives, so
  P(infected | positive) ≈ 20/69 ≈ 29% despite a "95% accurate" test
  [Wikipedia: Base rate fallacy; MetricGate: Base Rate Fallacy].
- **Confusing P(A|B) with P(B|A) (the prosecutor's fallacy).** P(evidence|innocent)
  is not P(innocent|evidence). The whole point of the theorem is that these differ
  unless the priors happen to make them equal.
- **Forgetting to normalize.** If you compute only the numerator P(D|H)·P(H), that is
  an unnormalized score, not a probability. Either divide by P(D) or renormalize a
  full set of competing hypotheses so they sum to 1.
- **Wrong or missing complement term in the denominator.** P(D) must include the
  false-positive mass P(D|Hᶜ)·P(Hᶜ), not just the true-positive mass.
- **Treating the prior as "subjective so it does not matter."** In this probability-
  theory framing the prior is just the unconditional base rate; it is a fact about
  the population and it materially changes the answer.

## When to reach for it

- You know P(effect | cause) and need P(cause | effect): a test's sensitivity but you
  want the chance of disease given a positive result.
- You want to update a probability as new evidence arrives — and, repeatedly, to let
  one result inform the next (today's posterior is tomorrow's prior; the full
  iterative inference treatment is in da-1-4-statistical-inference-foundations).
- You need to rank competing explanations of the same data via the proportionality or
  odds form.

## Sources

1. Wikipedia — "Bayes' theorem": https://en.wikipedia.org/wiki/Bayes'_theorem
2. GeeksforGeeks — "Bayes Theorem | Statement, Formula, Derivation, and Examples":
   https://www.geeksforgeeks.org/maths/bayes-theorem/
3. Lumen Learning — "Bayes' Theorem" (Mathematics for the Liberal Arts):
   https://courses.lumenlearning.com/waymakermath4libarts/chapter/bayes-theorem/
4. bookdown — "Bayesian Statistics the Fun Way: The Prior, Likelihood, and Posterior
   of Bayes' Theorem":
   https://bookdown.org/pbaumgartner/bayesian-fun/08-prior-likelihood-posterior.html
5. Wikipedia — "Base rate fallacy": https://en.wikipedia.org/wiki/Base_rate_fallacy
6. MetricGate — "The Base Rate Fallacy: Why Intuition Fails":
   https://metricgate.com/blogs/base-rate-fallacy-explained/
