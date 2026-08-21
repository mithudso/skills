<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-3-4-joint-marginal-conditional-probability` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-3-4-joint-marginal-conditional-probability
description: >-
  Joint, marginal, and conditional probability as the core relationships among
  multiple random variables in probability theory — the product (multiplication)
  rule, the sum rule (marginalization), conditioning, the chain rule, Bayes'
  rule, and conditional independence. Covers discrete (pmf / contingency tables)
  and continuous (pdf / integration) cases, worked numeric examples, and the
  reasoning fallacies analysts fall into (confusing P(A|B) with P(B|A), base-rate
  neglect, the prosecutor's fallacy).
  TRIGGER: a data-analysis question asks how to combine, decompose, or condition
  probabilities of two or more variables/events — "joint vs marginal vs
  conditional", "P(A and B)", "P(A given B)", "sum/marginalize out a variable",
  "build or read a contingency / cross-tab table of probabilities", "multiplication
  rule", "chain rule for a joint distribution", "is this conditional or marginal",
  "law of total probability", or a Bayes-rule update from a confusion matrix /
  diagnostic test.
  SKIP: defining a single distribution's pmf/pdf shape (defer to
  da-1-3-2-probability-mass-density-functions); naming/parameterizing a specific
  named distribution such as Normal/Binomial/Poisson (defer to the matching
  da-1-3-3-* skill); the axioms of probability or basic random-variable setup
  outside the joint/conditional framing (defer to da-1-3-* parents). Also defer
  any framing of "conditional probability" that belongs to a different taxonomy
  parent (e.g. Bayesian inference workflows, causal inference, ML feature
  engineering) — this skill is scoped to Foundations & Theory > Probability theory.
---

# Joint, Marginal & Conditional Probability

Scope: **Data Analysis > Foundations & Theory > Probability theory > Joint,
marginal & conditional probability.** This is the layer that connects the
probabilities of *several* variables. It sits above single-distribution shape
(pmf/pdf) and below applied Bayesian/causal workflows. Keep the focus on the
three quantities and the rules that move between them.

## The three quantities

For two events/variables A and B:

- **Joint** — both happen together. Written `P(A ∩ B)`, `P(A, B)`, or `P(A and B)`.
  "The probability of two (or more) events occurring simultaneously."
  ([LibreTexts 3.2](https://stats.libretexts.org/Workbench/ADAPT_Statistics_book/03:_Probability/3.02:_Marginal_Joint_and_Conditional_Probability),
  [MachineLearningMastery](https://machinelearningmastery.com/joint-marginal-and-conditional-probability-for-machine-learning/))
- **Marginal** — one event, *ignoring* the other. The probability of a single
  event "irrespective of the outcome of another variable," obtained by summing
  (or integrating) the joint over the other variable.
  ([Wikipedia: Marginal distribution](https://en.wikipedia.org/wiki/Marginal_distribution),
  [GeeksforGeeks](https://www.geeksforgeeks.org/maths/probability-joint-vs-marginal-vs-conditional/))
- **Conditional** — one event *given* another has occurred. Written `P(A | B)`:
  "the probability of one event occurring in the presence of a second event."
  ([LibreTexts 3.2](https://stats.libretexts.org/Workbench/ADAPT_Statistics_book/03:_Probability/3.02:_Marginal_Joint_and_Conditional_Probability))

The name "marginal" is literal: in a cross-tabulation (contingency table) the
single-variable totals are written in the **margins** of the table.

## The three rules that tie them together

**1. Definition of conditional probability**

```
P(A | B) = P(A, B) / P(B),   provided P(B) > 0
```
([Wikipedia: Conditional probability](https://en.wikipedia.org/wiki/Conditional_probability))

**2. Product / multiplication rule** (rearranged definition)

```
P(A, B) = P(A | B) · P(B) = P(B | A) · P(A)
```
This is "the fundamental rule of probability" / "product rule."
([MachineLearningMastery](https://machinelearningmastery.com/joint-marginal-and-conditional-probability-for-machine-learning/))
For **independent** events it simplifies to `P(A, B) = P(A) · P(B)`
([GeeksforGeeks](https://www.geeksforgeeks.org/maths/probability-joint-vs-marginal-vs-conditional/)).

**3. Sum rule / marginalization** — recover a marginal by summing the joint over
the nuisance variable:

```
Discrete:    P(X = x) = Σ_y  P(X = x, Y = y)
Continuous:  f_X(x)   = ∫    f(x, y) dy
```
([Wikipedia: Marginal distribution](https://en.wikipedia.org/wiki/Marginal_distribution),
[Stanford EE178 lect03](https://isl.stanford.edu/~abbas/ee178/lect03-2.pdf))

These three rules are the whole toolkit. Everything below is a consequence.

## Derived relationships you will actually use

- **Law of total probability** (sum rule + product rule). For a partition
  `{B_i}` of the sample space:
  ```
  P(A) = Σ_i P(A | B_i) · P(B_i)
  ```
  This is exactly marginalization with the joint expanded via the product rule.
  ([Brown AM165 ch.5](https://www.dam.brown.edu/people/huiwang/classes/am165/Prob_ch5_2007.pdf))

- **Bayes' rule** — invert a conditional by equating the two product-rule
  decompositions `P(A|B)P(B) = P(B|A)P(A)`:
  ```
  P(B | A) = P(A | B) · P(B) / P(A)
  ```
  "These are perhaps the two most fundamental equations in probability theory."
  ([Bayes is chain rule](https://eregis.github.io/blog/2024/12/21/bayes-theorem-chain-rule.html))

- **Chain rule** — factor any joint into a product of conditionals:
  ```
  P(A₁,…,A_k) = P(A₁) · P(A₂|A₁) · P(A₃|A₁,A₂) ⋯ P(A_k | A₁,…,A_{k-1})
  ```
  ([Wikipedia: Chain rule (probability)](https://en.wikipedia.org/wiki/Chain_rule_(probability)),
  [Berkeley CS188 probability](https://inst.eecs.berkeley.edu/~cs188/textbook/bayes-nets/probability.html))

- **Conditional independence** — `A ⟂ B | C` means `P(A,B|C) = P(A|C)·P(B|C)`.
  It is the assumption that lets the chain rule collapse into far fewer terms
  (the basis of Bayesian networks).
  ([Berkeley CS188 probability](https://inst.eecs.berkeley.edu/~cs188/textbook/bayes-nets/probability.html))

## Worked example — reading a contingency table

A survey of 100 people records device (Phone / Laptop) by OS preference:

|             | Likes A | Likes B | **Row total** |
|-------------|--------:|--------:|--------------:|
| **Phone**   |      30 |      20 |        **50** |
| **Laptop**  |      10 |      40 |        **50** |
| **Col total** | **40** | **60** |       **100** |

- **Joint:** `P(Phone, A) = 30/100 = 0.30` (a single inner cell ÷ grand total).
- **Marginal:** `P(Phone) = 50/100 = 0.50` (a margin total ÷ grand total) — note
  `0.30 + 0.20 = 0.50`, i.e. summing the joint over OS recovers the marginal.
- **Conditional:** `P(A | Phone) = P(A, Phone)/P(Phone) = 0.30/0.50 = 0.60`
  (restrict to the Phone *row*, then renormalize within it).

Check the product rule: `P(A|Phone)·P(Phone) = 0.60 · 0.50 = 0.30 = P(A, Phone)`. ✓
Independence check: `P(A)·P(Phone) = 0.40 · 0.50 = 0.20 ≠ 0.30`, so device and OS
preference are **dependent**. (Construction following
[LibreTexts 3.2](https://stats.libretexts.org/Workbench/ADAPT_Statistics_book/03:_Probability/3.02:_Marginal_Joint_and_Conditional_Probability).)

### How to decide which one a question is asking for

- "…and…" / "both" / a single inner cell → **joint**.
- "…overall" / "ignoring the other factor" / a margin total → **marginal**.
- "…given…" / "among the X" / "of the people who…" → **conditional** (you are
  renormalizing inside one row or column).

## Pitfalls (where analysts get it wrong)

1. **Confusing `P(A|B)` with `P(B|A)`** — the *confusion of the inverse* /
   inverse fallacy. They are generally unequal. A dengue test can have
   `P(positive | disease) = 90%` while `P(disease | positive) ≈ 15%` when the
   disease is rare. Use Bayes' rule to flip them — do not just swap the bar.
   ([Confusion of the inverse](https://en.wikipedia.org/wiki/Confusion_of_the_inverse))

2. **Base-rate neglect** — ignoring the marginal `P(B)` (the prior). With a
   1-in-10,000 disease and a 99%-accurate test, a positive result still implies
   only ~1% chance of disease, because the base rate dominates. The marginal
   term in Bayes' rule is not optional.
   ([Base rate fallacy](https://en.wikipedia.org/wiki/Base_rate_fallacy))

3. **Prosecutor's fallacy** — treating `P(evidence | innocent)` as if it equals
   `P(innocent | evidence)`. A small random-match probability is *not* the
   probability of innocence.
   ([CEBM Oxford](https://www.cebm.ox.ac.uk/news/views/the-prosecutors-fallacy))

4. **Assuming independence to multiply marginals.** `P(A,B) = P(A)·P(B)` holds
   *only* when A and B are independent. Defaulting to it on correlated variables
   (e.g. naively multiplying feature probabilities) understates or overstates the
   joint. Verify with `P(A,B) =? P(A)·P(B)` as above.
   ([GeeksforGeeks](https://www.geeksforgeeks.org/maths/probability-joint-vs-marginal-vs-conditional/))

5. **Conditioning on a zero-probability event.** `P(A|B)` is undefined when
   `P(B) = 0` (discrete). For continuous variables, condition via densities, not
   point probabilities (`P(Y = y) = 0` for continuous Y).
   ([Wikipedia: Conditional probability](https://en.wikipedia.org/wiki/Conditional_probability))

6. **Renormalization slip.** A conditional must sum to 1 over the conditioned
   variable. After restricting to a row, divide by that row's total — forgetting
   the divisor leaves you with a joint, not a conditional.

## Quick reference

| Want | From a joint table | Formula |
|------|--------------------|---------|
| Joint `P(A,B)` | inner cell ÷ grand total | — |
| Marginal `P(A)` | margin total ÷ grand total | `Σ_b P(A,b)` |
| Conditional `P(A|B)` | cell ÷ that row/col total | `P(A,B)/P(B)` |
| Invert a conditional | — | Bayes: `P(B|A)=P(A|B)P(B)/P(A)` |
| Factor a big joint | — | chain rule |

## Sources

1. [Marginal, Joint, and Conditional Probability — Statistics LibreTexts (3.2)](https://stats.libretexts.org/Workbench/ADAPT_Statistics_book/03:_Probability/3.02:_Marginal_Joint_and_Conditional_Probability) — definitions + contingency-table method.
2. [Marginal distribution — Wikipedia](https://en.wikipedia.org/wiki/Marginal_distribution) — discrete sum / continuous integral definitions; marginal vs conditional.
3. [Conditional probability — Wikipedia](https://en.wikipedia.org/wiki/Conditional_probability) — `P(A|B)=P(A,B)/P(B)`, zero-probability conditioning.
4. [Chain rule (probability) — Wikipedia](https://en.wikipedia.org/wiki/Chain_rule_(probability)) — joint factorization.
5. [A Gentle Introduction to Joint, Marginal, and Conditional Probability — MachineLearningMastery](https://machinelearningmastery.com/joint-marginal-and-conditional-probability-for-machine-learning/) — product rule, DA/ML framing.
6. [Joint, Marginal, and Conditional pmfs; Bayes Rule — Stanford EE178 lect03](https://isl.stanford.edu/~abbas/ee178/lect03-2.pdf) — discrete pmf treatment.
7. [Multivariate Probability Distributions — Brown AM165 ch.5](https://www.dam.brown.edu/people/huiwang/classes/am165/Prob_ch5_2007.pdf) — law of total probability, continuous case.
8. [Probability review — UC Berkeley CS188](https://inst.eecs.berkeley.edu/~cs188/textbook/bayes-nets/probability.html) — chain rule, conditional independence.
9. [Joint vs Marginal vs Conditional — GeeksforGeeks](https://www.geeksforgeeks.org/maths/probability-joint-vs-marginal-vs-conditional/) — independence simplification.
10. [Confusion of the inverse — Wikipedia](https://en.wikipedia.org/wiki/Confusion_of_the_inverse), [Base rate fallacy — Wikipedia](https://en.wikipedia.org/wiki/Base_rate_fallacy), [The Prosecutor's Fallacy — CEBM, Oxford](https://www.cebm.ox.ac.uk/news/views/the-prosecutors-fallacy) — reasoning pitfalls.
