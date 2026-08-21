<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-3-3-probability-distributions` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-3-3-probability-distributions
description: >-
  Foundational probability-theory treatment of a probability DISTRIBUTION: the
  formal object (a probability measure / the law of a random variable), how a
  distribution is fully described (CDF as the universal descriptor; PMF and PDF
  as type-specific descriptors), the classification of distributions into
  discrete, absolutely continuous, and singular continuous via the Lebesgue
  decomposition, support, parameters/families, and the descriptors that pin a
  distribution down (moments, MGF, characteristic function). TRIGGER: "what is a
  probability distribution", how a distribution is defined as a measure or law,
  what fully determines/identifies a distribution, distinguishing a distribution
  from the random variable that induces it, classifying a distribution as
  discrete/continuous/mixed/singular, reasoning with a CDF, what "support" or a
  "parametric family" means, or whether the MGF/characteristic function
  determines the distribution. SKIP: the random variable itself as a measurable
  function and its expectation/variance machinery (defer to
  da-1-3-1-random-variables); specific NAMED families and their parameters —
  Normal, Binomial, Poisson, Exponential, Uniform — each defer to their own
  da-1-3-3-N skill; PMF/PDF treated as the focal topic in their own right (defer
  to da-1-3-2-probability-mass-density-functions); joint, marginal, and
  conditional distributions across several variables (defer to
  da-1-3-4-joint-marginal-conditional-probability); estimating or fitting a
  distribution FROM sample data / inferential statistics (defer to
  estimation/inference skills); the broader survey of probability theory (defer
  to da-1-3-probability-theory). This skill covers the concept only as it sits
  under Data Analysis > Foundations & Theory > Probability theory > Probability
  distributions.
---

# Probability Distributions (Probability Theory foundations)

A probability distribution is the object that says *how probability is spread
over the possible values of a random quantity*. This skill scopes the concept as
it sits in the foundations of probability theory: what a distribution **is** as a
formal object, how it is **described** (CDF, and the type-specific PMF/PDF), how
distributions are **classified** (discrete / absolutely continuous / singular),
and what **pins one down** (support, parameters, moments, transforms). It does
not cover the random variable itself, any single named family, or fitting a
distribution to data — those are separate concepts elsewhere in the taxonomy.

## The one idea to fix first

A probability distribution is **not** the random variable. The random variable
`X` is a function from outcomes to values; the **distribution** is the
*probability measure that `X` induces on the value space*. Two random variables
on entirely different sample spaces can share the same distribution. When people
say "`X` has a Normal distribution," they are naming the induced measure, not the
function itself ([Wikipedia: Probability distribution](https://en.wikipedia.org/wiki/Probability_distribution)).

## Formal definition

A probability distribution is a **probability measure** — a function that assigns
a number in `[0, 1]` to events, is countably additive, and assigns `1` to the
whole space ([Wikipedia: Probability distribution](https://en.wikipedia.org/wiki/Probability_distribution)).

In the measure-theoretic setup, start from a probability space `(Ω, ℱ, P)` and a
measurable space `(E, ℰ)` (for a real-valued quantity, `E = ℝ` with the Borel
σ-algebra). A random variable `X: Ω → E` induces its distribution as the
**pushforward measure** (also called the **law** of `X`):

    P_X(B) = X₊(P)(B) = P(X⁻¹(B)) = P(X ∈ B),   for every B ∈ ℰ.

This carries all the probability mass off the abstract sample space `Ω` and onto
the value space `E`, where it can actually be computed with
([Wikipedia: Probability distribution](https://en.wikipedia.org/wiki/Probability_distribution);
[IIT Madras EE5110, "Types of Random Variables"](https://www.ee.iitm.ac.in/~krishnaj/EE5110_files/notes/lecture11_Part_Two_Types_Of_Random_Variables.pdf)).

Equality **in distribution** (`X =ᵈ Y`, meaning `P_X = P_Y`) is strictly weaker
than almost-sure equality (`P(X = Y) = 1`). Most of probability and statistics
only ever cares about the distribution, not which underlying `Ω` produced it.

## The CDF — the universal descriptor

Every real-valued distribution is fully and uniquely described by its
**cumulative distribution function (CDF)**:

    F_X(x) = P(X ≤ x).

A function `F` is the CDF of some distribution **iff** it is

1. **non-decreasing**,
2. **right-continuous**, and
3. has limits `lim_{x→−∞} F(x) = 0` and `lim_{x→+∞} F(x) = 1`

([Wikipedia: Cumulative distribution function](https://en.wikipedia.org/wiki/Cumulative_distribution_function)).

Two distributions are equal **iff** their CDFs agree everywhere. This is why the
CDF is the universal descriptor: it exists for *every* distribution (discrete,
continuous, mixed, or singular), whereas a PMF or PDF exists only for particular
types. When you are unsure which type you are dealing with, reason with the CDF.

Useful CDF identities:

- `P(a < X ≤ b) = F_X(b) − F_X(a)`.
- `P(X = x) = F_X(x) − F_X(x⁻)`, the size of the **jump** at `x` (zero where `F`
  is continuous).

## Classification: discrete, absolutely continuous, singular

By the **Lebesgue decomposition theorem**, any distribution on the real line
splits uniquely into a mixture of three mutually exclusive parts:

    F = α·F_discrete + β·F_ac + γ·F_singular,   with α + β + γ = 1, α,β,γ ∈ [0,1]

([MSU lecture 17, "Lebesgue Decomposition and Probability Theory"](https://hal.cse.msu.edu/teaching/2025-spring-computational-foundations-of-ai/scribe/lecture-17.pdf);
[Wikipedia: Probability distribution](https://en.wikipedia.org/wiki/Probability_distribution)).

### Discrete

Probability is concentrated on a finite or countable set of points (**atoms**).
The CDF is a **step function**; each value `x` with `P(X = x) > 0` is a jump. The
distribution is described by a **probability mass function (PMF)**
`p_X(x) = P(X = x)` with `Σ_x p_X(x) = 1`, and `F_X(x) = Σ_{x_i ≤ x} p_X(x_i)`
([Wikipedia: Cumulative distribution function](https://en.wikipedia.org/wiki/Cumulative_distribution_function)).

### Absolutely continuous

There is a **probability density function (PDF)** `f_X ≥ 0` with respect to
Lebesgue measure such that

    F_X(x) = ∫₋∞ˣ f_X(t) dt,   P(a ≤ X ≤ b) = ∫_a^b f_X(t) dt,   ∫₋∞^∞ f_X = 1.

The CDF is continuous and `f_X(x) = F_X′(x)` where the derivative exists. Critical
fact: every single point has probability zero, `P(X = x) = 0`, so the PDF value is
a **density, not a probability** (it may exceed 1), and endpoint inclusion does
not matter (`P(a ≤ X ≤ b) = P(a < X < b)`)
([Wikipedia: Probability distribution](https://en.wikipedia.org/wiki/Probability_distribution)).

### Singular continuous

The CDF is **continuous everywhere yet has derivative zero almost everywhere**, so
there is no PMF (no atoms) and no PDF (no density). The canonical example is the
**Cantor distribution**. Rare in applied work, but it is the reason "continuous"
and "has a density" are not synonyms: a distribution can be continuous without
being absolutely continuous
([Wikipedia: Probability distribution](https://en.wikipedia.org/wiki/Probability_distribution);
[IIT Madras EE5110](https://www.ee.iitm.ac.in/~krishnaj/EE5110_files/notes/lecture11_Part_Two_Types_Of_Random_Variables.pdf)).

### Mixed

Any nontrivial mixture (`α, β > 0`) is a **mixed** distribution: the CDF has both
jumps and a continuous, growing part. A measurement censored at a boundary (e.g.
sensor readings clamped at 0, with a continuous spread above 0) is the everyday
example: point mass at the boundary plus a density elsewhere.

## Support

The **support** of a distribution is the set of values `x` such that the random
variable has positive probability of falling in *every* open neighborhood of `x`
([Wikipedia: Probability distribution](https://en.wikipedia.org/wiki/Probability_distribution)).
Informally, the smallest closed set carrying all the probability. For a discrete
distribution it is the set of atoms; for a continuous one it is (the closure of)
the region where the density is positive. Watch the support boundaries: an
Exponential is supported on `[0, ∞)`, a Uniform`(a,b)` on `[a, b]`. Specific
families' supports belong to those families' own skills.

## Parameters and families

A **parametric family** is a set of distributions indexed by one or more
**parameters** `θ` (e.g. a Normal family indexed by `(μ, σ²)`). Fixing the
parameter picks out one member. Parameters typically split into

- **location** (shift the distribution along the axis),
- **scale** (stretch/compress its spread), and
- **shape** (change the form itself, e.g. skewness).

This skill covers the *concept* of a family and its parameter roles; the concrete
families (Normal, Binomial, Poisson, Exponential, Uniform) each have their own
`da-1-3-3-N` skill and are out of scope here.

## What pins a distribution down (descriptors and transforms)

- **CDF** — always determines the distribution uniquely (see above).
- **Moments.** The `k`-th moment is `E[X^k]`; the mean is the 1st moment and the
  variance is `E[X²] − (E[X])²`. Moments summarize but do **not** in general
  determine a distribution: distinct distributions can share all moments (the
  classic example is the log-normal and a family of distributions matching all its
  moments).
- **Moment generating function (MGF).** `M_X(t) = E[e^{tX}]`. When the MGF is
  finite on an open interval around `0`, it **uniquely determines** the
  distribution, and `M_X^{(k)}(0) = E[X^k]` recovers the moments. The MGF may fail
  to exist (heavy tails)
  ([Oxford Part A Probability, J. Martin](https://www.stats.ox.ac.uk/~laws/msc/Probability.pdf);
  [MIT 18.440 Lecture 27, "Moment generating functions and characteristic functions"](https://dspace.mit.edu/bitstream/handle/1721.1/97467/18-440-spring-2011/contents/lecture-notes/MIT18_440S11_Lecture27.pdf)).
- **Characteristic function.** `φ_X(t) = E[e^{itX}]` always exists (it is the
  Fourier transform of the distribution) and **always uniquely determines** the
  distribution — the reliable tool when the MGF does not exist
  ([Oxford Part A Probability](https://www.stats.ox.ac.uk/~laws/msc/Probability.pdf);
  [MIT 18.440 Lecture 27](https://dspace.mit.edu/bitstream/handle/1721.1/97467/18-440-spring-2011/contents/lecture-notes/MIT18_440S11_Lecture27.pdf)).

## Worked example: a mixed distribution from censoring

Let `Y` be the recorded output of a sensor that reads a true continuous signal but
clamps any negative reading to `0`. Suppose the true signal `Z` has density and
`P(Z ≤ 0) = 0.3`. Then `Y = max(Z, 0)` has distribution:

- **Atom at 0:** `P(Y = 0) = P(Z ≤ 0) = 0.3` — a discrete component (`α = 0.3`).
- **Continuous part on `(0, ∞)`:** density equal to that of `Z` there
  (`β = 0.7`, `γ = 0`).
- **CDF:** `F_Y(y) = 0` for `y < 0`; jumps to `0.3` at `y = 0`; then rises
  continuously toward `1`.

`Y` is neither discrete nor absolutely continuous — it is genuinely mixed, and the
CDF is the only descriptor that captures it cleanly. The PMF describes only the
jump; a "PDF" alone cannot account for the `0.3` point mass.

## Common pitfalls (checklist)

- **Conflating the distribution with the random variable.** The distribution is
  the induced measure; the variable is the function. Same law ≠ same variable.
- **Assuming "continuous" means "has a density."** Singular continuous
  distributions are continuous with no PDF. "Has a PDF" = *absolutely* continuous.
- **Treating every distribution as discrete-or-continuous.** Mixed distributions
  exist (censoring, truncation with a spike); use the CDF.
- **Reading a PDF value as a probability.** It is a density and can exceed 1; only
  integrals over intervals are probabilities, and points have probability zero.
- **Believing moments determine the distribution.** They do not in general; only
  the CDF, the characteristic function, or an MGF finite near 0 do.
- **Forgetting the support.** Integrating or summing outside the support (e.g.
  below 0 for an Exponential) silently corrupts results.
- **Ignoring CDF jumps for discrete/mixed parts.** `P(X = x)` is the jump size
  `F(x) − F(x⁻)`, not `0`, whenever there is an atom.

## Sources

1. Wikipedia — Probability distribution. https://en.wikipedia.org/wiki/Probability_distribution
2. Wikipedia — Cumulative distribution function. https://en.wikipedia.org/wiki/Cumulative_distribution_function
3. Michigan State University — Lecture 17: Lebesgue Decomposition and Probability Theory. https://hal.cse.msu.edu/teaching/2025-spring-computational-foundations-of-ai/scribe/lecture-17.pdf
4. IIT Madras EE5110 (Krishna Jagannathan) — Probability Foundations, Lecture 11: Types of Random Variables. https://www.ee.iitm.ac.in/~krishnaj/EE5110_files/notes/lecture11_Part_Two_Types_Of_Random_Variables.pdf
5. University of Oxford — Part A Probability lecture notes (James Martin, Michaelmas 2016). https://www.stats.ox.ac.uk/~laws/msc/Probability.pdf
6. MIT 18.440 — Lecture 27: Moment generating functions and characteristic functions. https://dspace.mit.edu/bitstream/handle/1721.1/97467/18-440-spring-2011/contents/lecture-notes/MIT18_440S11_Lecture27.pdf
