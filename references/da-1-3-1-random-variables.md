<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-3-1-random-variables` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-3-1-random-variables
description: >-
  Foundational probability-theory treatment of random variables: the formal
  definition as a measurable function on a probability space, the induced
  distribution (law / pushforward), discrete vs. continuous types, PMF/PDF/CDF,
  expectation and variance, functions of random variables (LOTUS), and
  independence. TRIGGER: questions about what a random variable IS at the
  probability-foundations level; "is X a valid random variable", measurability
  of a map on a sample space, deriving a distribution/law from outcomes,
  choosing between PMF and PDF, computing E[X]/Var(X)/E[g(X)] from a
  distribution, distinguishing the random variable from its distribution, or
  building intuition for discrete vs. continuous outcomes. SKIP: named
  distribution families and their parameters (defer to a distributions skill);
  estimating moments or distributions FROM sample data / inferential statistics
  (defer to estimation/inference skills); stochastic processes, time-indexed or
  random-vector/matrix modeling (defer to process skills); measure-theory
  integration as a topic in itself (defer to a measure-theory skill); the
  broader survey of probability theory or other probability sub-topics (defer to
  da-1-3-probability-theory). This skill covers the concept only as it sits under
  Data Analysis > Foundations & Theory > Probability theory > Random variables.
---

# Random Variables (Probability Theory foundations)

A random variable is the bridge between an abstract model of chance and the
numbers we actually compute with. This skill scopes the concept as it sits in
the foundations of probability theory: what a random variable *is*, what object
it induces (its distribution), and the core machinery (PMF/PDF/CDF, expectation,
variance, functions, independence). It does not cover specific named
distributions, estimation from data, or stochastic processes — those are
separate concepts elsewhere in the taxonomy.

## The one idea to fix first

A random variable is **not random and not a variable**. It is a deterministic
*function* from outcomes to values. The randomness lives in the underlying
probability measure on the sample space, not in the function. Once you internalize
this, almost every confusion about random variables dissolves
([Wikipedia: Random variable](https://en.wikipedia.org/wiki/Random_variable);
[ETH Zürich, "A random variable is neither random nor variable"](https://ethz.ch/content/dam/ethz/special-interest/mavt/dynamic-systems-n-control/idsc-dam/Lectures/Stochastic-Systems/Probability.pdf)).

## Formal definition

Let `(Ω, ℱ, P)` be a probability space:

- `Ω` — the **sample space**, the set of all possible outcomes `ω`.
- `ℱ` — a **σ-algebra** on `Ω`: the collection of subsets we are allowed to
  assign probability to (the *events*).
- `P` — a **probability measure**: `P: ℱ → [0,1]`, with `P(Ω) = 1`.

Let `(E, ℰ)` be a measurable space (for a real-valued random variable, `E = ℝ`
and `ℰ` is the Borel σ-algebra). A random variable is a **measurable function**

    X: Ω → E    such that for every B ∈ ℰ,   X⁻¹(B) = { ω ∈ Ω : X(ω) ∈ B } ∈ ℱ.

The measurability condition is the whole point: it guarantees that the preimage
of any measurable set of values is an event we can assign a probability to. If
`X⁻¹(B)` were not in `ℱ`, the expression `P(X ∈ B)` would be undefined
([Wikipedia: Random variable](https://en.wikipedia.org/wiki/Random_variable);
[Wikipedia: Measurable function](https://en.wikipedia.org/wiki/Measurable_function);
[Bernstein, "Demystifying measure-theoretic probability theory, part 2"](https://mbernste.github.io/posts/measure_theory_2/)).

For real-valued `X`, it suffices to check measurability on a generating class:
`X` is a random variable iff `{ω : X(ω) ≤ x} ∈ ℱ` for every `x ∈ ℝ`, because the
sets `(-∞, x]` generate the Borel σ-algebra
([UWaterloo lecture notes ch. 3, "Random Variables and Measurable Functions"](https://sas.uwaterloo.ca/~dlmcleis/s901/chapt3.pdf)).

## The distribution (law) — the object you usually care about

Most of the time you never touch `Ω` directly. You work with the **distribution**
(or **law**) of `X`: the pushforward measure `P_X` on `(E, ℰ)` defined by

    P_X(B) = P(X⁻¹(B)) = P(X ∈ B).

This moves all the probability mass from the abstract sample space onto the value
space. Two random variables can live on completely different sample spaces yet
have the *same distribution* — written `X =ᵈ Y`. Equality in distribution is
weaker than almost-sure equality (`X =ᵃ·ˢ· Y`, meaning `P(X = Y) = 1`)
([Wikipedia: Random variable](https://en.wikipedia.org/wiki/Random_variable)).

Keep the distinction sharp: **the random variable is the function; the distribution
is the measure it induces.** Conflating them is the most common foundational error.

## Describing the distribution: CDF, PMF, PDF

The **cumulative distribution function (CDF)** works for every real-valued random
variable:

    F_X(x) = P(X ≤ x).

`F_X` is non-decreasing, right-continuous, with `F_X(-∞)=0` and `F_X(+∞)=1`. The
CDF fully determines the distribution, which is why it is the universal descriptor.

### Discrete random variables

`X` is **discrete** when its range is finite or countably infinite. Its
distribution is captured by a **probability mass function (PMF)**:

    p_X(x) = P(X = x),   with   Σ_x p_X(x) = 1.

Here `P(X = x) > 0` is meaningful — individual values carry positive mass
([Colorado AMath ch.3, "Discrete Random Variables and Probability Distributions"](https://www.colorado.edu/amath/sites/default/files/attached-files/ch3_0.pdf);
[Cambridge IntroProb lecture 3](https://www.cl.cam.ac.uk/teaching/2324/IntroProb/slides/03-variance-discr-distr-answers-hidden-handout.pdf)).

### Continuous (absolutely continuous) random variables

`X` is **(absolutely) continuous** when there is a **probability density function
(PDF)** `f_X ≥ 0` with

    F_X(x) = ∫₋∞ˣ f_X(t) dt,   and   ∫₋∞^∞ f_X(t) dt = 1.

Critical pitfall: for a continuous `X`, `P(X = x) = 0` for every single point `x`.
The PDF value `f_X(x)` is a **density, not a probability** — it can exceed 1.
Probabilities come only from integrating over an interval: `P(a ≤ X ≤ b) =
∫_a^b f_X(t) dt`. Because point mass is zero, `P(a ≤ X ≤ b) = P(a < X < b)`; the
endpoints do not matter for continuous variables (they very much matter for
discrete ones)
([Wikipedia: Random variable](https://en.wikipedia.org/wiki/Random_variable)).

> Not every random variable is purely discrete or purely continuous. Mixed
> variables (e.g. a measurement censored at zero) have a CDF with both jumps and
> a continuous part. When in doubt, reason with the CDF.

## Expectation and variance

The **expected value** `E[X]` is the probability-weighted average — the "center of
mass" of the distribution:

- Discrete:  `E[X] = Σ_x x · p_X(x)`
- Continuous: `E[X] = ∫₋∞^∞ x · f_X(x) dx`

(provided the sum/integral converges absolutely; otherwise the expectation does
not exist). The **variance** measures spread about the mean:

    Var(X) = E[(X − μ)²] = E[X²] − (E[X])²,   where μ = E[X].

The second form is the usual computational shortcut. The standard deviation is
`√Var(X)`, restoring the original units
([Purdue Northwest "Mean and Variance" lecture notes 4](https://www.pnw.edu/wp-content/uploads/2020/03/lecturenotes4-10.pdf);
[UW CSE312 section 4, "Random Variables and Expectation"](https://courses.cs.washington.edu/courses/cse312/21sp/files/sections/section04.pdf)).

Linearity of expectation always holds, even without independence:
`E[aX + bY + c] = a·E[X] + b·E[Y] + c`. Variance does **not** add in general —
`Var(X + Y) = Var(X) + Var(Y)` requires `X` and `Y` to be uncorrelated (independence
suffices).

## Functions of random variables — LOTUS

If `Y = g(X)`, then `Y` is itself a random variable (the composition of a
measurable function with a measurable function is measurable). You do **not** need
to derive the distribution of `Y` to get its expectation. The **Law of the
Unconscious Statistician (LOTUS)** weights the transformed values by the
*original* distribution of `X`:

- Discrete:  `E[g(X)] = Σ_x g(x) · p_X(x)`
- Continuous: `E[g(X)] = ∫₋∞^∞ g(x) · f_X(x) dx`

The name reflects that people often treat this as the *definition* of `E[g(X)]`,
when formally it is a theorem derived from the true definition of expectation
([Wikipedia: Law of the unconscious statistician](https://en.wikipedia.org/wiki/Law_of_the_unconscious_statistician);
[probabilitycourse.com 3.2.3, "Functions of Random Variables"](https://www.probabilitycourse.com/chapter3/3_2_3_functions_random_var.php);
[Yale S&DS 241 lecture 6, "Functions of random variables. LOTUS rule."](http://www.stat.yale.edu/~yw562/teaching/241/lec06.pdf)).

LOTUS generalizes to several variables: `E[g(X, Y)] = ∫∫ g(x, y) f_{X,Y}(x, y) dx dy`.

## Independence

Random variables `X` and `Y` are **independent** when the joint distribution
factors:

    P(X ∈ A, Y ∈ B) = P(X ∈ A) · P(Y ∈ B)   for all measurable A, B,

equivalently `F_{X,Y}(x,y) = F_X(x) · F_Y(y)`, or in the density case
`f_{X,Y}(x,y) = f_X(x) · f_Y(y)`. A key consequence used constantly: if `X ⟂ Y`,
then `E[f(X) g(Y)] = E[f(X)] · E[g(Y)]` for any measurable `f, g`
([Yale S&DS 241 lecture 6](http://www.stat.yale.edu/~yw562/teaching/241/lec06.pdf)).

Independence implies zero correlation, but zero correlation does **not** imply
independence — a frequent trap.

## Worked example (discrete)

Roll one fair die; `Ω = {1,…,6}`, `P` uniform. Define `X(ω) = ω` (the face) — a
valid random variable since every `{X ≤ x}` is an event.

- PMF: `p_X(k) = 1/6` for `k = 1,…,6`.
- `E[X] = Σ k·(1/6) = 21/6 = 3.5`. (The mean need not be an attainable value.)
- `E[X²] = Σ k²·(1/6) = 91/6 ≈ 15.17`, so `Var(X) = 91/6 − (3.5)² ≈ 2.92`.
- Let `g(x) = (x − 3.5)²`. By LOTUS, `E[g(X)] = Σ (k − 3.5)²·(1/6) = Var(X) ≈ 2.92`
  — confirming LOTUS without ever forming the distribution of `g(X)`.

## Common pitfalls (checklist)

- **Treating `X` as the distribution.** `X` is a function; `P_X` is the measure it
  induces. Same law ≠ same variable.
- **Reading a PDF value as a probability.** `f_X(x)` is a density; it can exceed 1.
  Only integrals over intervals are probabilities.
- **Forgetting `P(X = x) = 0` for continuous `X`.** Endpoint inclusion changes
  nothing for continuous variables but matters for discrete ones.
- **Skipping the measurability check.** A map on `Ω` is only a random variable if
  preimages of measurable sets are events.
- **Assuming `Var` is linear or that uncorrelated implies independent.** Neither
  holds in general.
- **Assuming every distribution is moment-bearing.** `E[X]` or `Var(X)` may fail to
  exist if the defining sum/integral diverges.

## Sources

1. Wikipedia — Random variable. https://en.wikipedia.org/wiki/Random_variable
2. Wikipedia — Measurable function. https://en.wikipedia.org/wiki/Measurable_function
3. Wikipedia — Law of the unconscious statistician. https://en.wikipedia.org/wiki/Law_of_the_unconscious_statistician
4. M. N. Bernstein — Demystifying measure-theoretic probability theory (part 2: random variables). https://mbernste.github.io/posts/measure_theory_2/
5. ETH Zürich — Probability Theory: "A random variable is neither random nor variable." https://ethz.ch/content/dam/ethz/special-interest/mavt/dynamic-systems-n-control/idsc-dam/Lectures/Stochastic-Systems/Probability.pdf
6. University of Waterloo — Chapter 3: Random Variables and Measurable Functions. https://sas.uwaterloo.ca/~dlmcleis/s901/chapt3.pdf
7. University of Colorado AMath — Ch. 3: Discrete Random Variables and Probability Distributions. https://www.colorado.edu/amath/sites/default/files/attached-files/ch3_0.pdf
8. Purdue Northwest — Mean and Variance (lecture notes 4). https://www.pnw.edu/wp-content/uploads/2020/03/lecturenotes4-10.pdf
9. University of Washington CSE312 — Section 4: Random Variables and Expectation. https://courses.cs.washington.edu/courses/cse312/21sp/files/sections/section04.pdf
10. probabilitycourse.com — 3.2.3 Functions of Random Variables. https://www.probabilitycourse.com/chapter3/3_2_3_functions_random_var.php
11. Yale S&DS 241 — Lecture 6: Functions of random variables, LOTUS rule. http://www.stat.yale.edu/~yw562/teaching/241/lec06.pdf
