<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-3-2-probability-mass-density-functions` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-3-2-probability-mass-density-functions
description: >-
  Foundational probability-theory treatment of probability mass functions (PMF)
  and probability density functions (PDF): the formal definitions, the defining
  properties (non-negativity, normalization, support), why a PMF gives a
  probability directly while a PDF gives a density that must be integrated, the
  discrete-vs-continuous split, the relationship to the CDF, the change-of-
  variables formula with the Jacobian, expectation computed via PMF/PDF, and
  mixed distributions. TRIGGER: questions about what a PMF or PDF IS at the
  probability-foundations level; "is this a valid PMF/PDF", checking that a
  candidate function sums or integrates to 1, why a PDF can exceed 1, why
  P(X=x)=0 for a continuous variable, getting a PDF from a CDF (or vice versa),
  finding the density of a transformed variable Y=g(X), or distinguishing a mass
  function from a density function. SKIP: named distribution families and their
  specific PMFs/PDFs such as normal, binomial, Poisson, exponential, uniform
  (defer to the per-distribution skills); the random variable concept itself or
  expectation/variance machinery as a standalone topic (defer to
  da-1-3-1-random-variables); joint/marginal/conditional densities across
  several variables (defer to da-1-3-4-joint-marginal-conditional-probability);
  estimating a density FROM sample data — kernel density estimation, histograms,
  fitting (defer to estimation/inference skills); the broader survey of
  probability theory (defer to da-1-3-probability-theory). This skill covers the
  concept only as it sits under Data Analysis > Foundations & Theory >
  Probability theory > Probability mass / density functions.
---

# Probability mass / density functions (Probability Theory foundations)

A PMF and a PDF are the two ways to *describe* the distribution of a random
variable by attaching a number to each possible value. They are siblings, not
the same object: a PMF answers "what is the probability of exactly this value?"
while a PDF answers "what is the probability *per unit length* near this value?"
Conflating the two is the single most common foundational error, so this skill
keeps the contrast explicit throughout.

Scope note: this is the *general machinery* of mass and density functions. The
PMF/PDF of a specific named family (normal, binomial, Poisson, exponential,
uniform) lives in its own skill; the random-variable concept and expectation as
topics in themselves live in `da-1-3-1-random-variables`; multi-variable joint
densities live in `da-1-3-4-joint-marginal-conditional-probability`.

## The one idea to fix first

A **PMF** is a probability. A **PDF** is a density. `p_X(x) = P(X = x)` is a real
probability you can read off directly; `f_X(x)` is *not* a probability — it is a
rate, and you only get a probability by integrating it over an interval. Because
of this, a PMF value can never exceed 1, but a PDF value can be arbitrarily large
([Wikipedia: Probability density function](https://en.wikipedia.org/wiki/Probability_density_function);
[StatLect: Probability density function](https://www.statlect.com/glossary/probability-density-function)).

## Probability mass function (discrete)

A random variable `X` is **discrete** when its range is finite or countably
infinite. Its distribution is fully described by the **probability mass function**

    p_X(x) = P(X = x),    p_X: ℝ → [0, 1].

The PMF must satisfy two defining conditions; a function is a valid PMF if and
only if both hold:

1. **Non-negativity:** `p_X(x) ≥ 0` for all `x`. A probability cannot be negative.
2. **Normalization:** `Σ_x p_X(x) = 1`, summed over the support — total
   probability is 1 because some outcome must occur.

The **support** is the countable set of `x` where `p_X(x) > 0`; outside it the
PMF is zero ([Wikipedia: Probability mass function](https://en.wikipedia.org/wiki/Probability_mass_function);
[Statistics LibreTexts / STAT 350: Discrete RVs and PMFs](https://treese41528.github.io/STAT350/Website/chapter5/lectures/5-1-discrete-rvs-and-pmfs.html)).

For a discrete variable, `P(X = x) > 0` is meaningful — individual values carry
positive mass — and the probability of falling in a set `A` is just a sum:
`P(X ∈ A) = Σ_{x ∈ A} p_X(x)`. The "mass" metaphor is load-bearing: think of a
fixed amount of mass (total 1) distributed as lumps over the support, with the
physical intuition that mass is conserved exactly as total probability is
conserved ([Wikipedia: Probability mass function](https://en.wikipedia.org/wiki/Probability_mass_function)).

Measure-theoretically, the PMF is the Radon–Nikodym derivative of the
distribution of `X` (its pushforward measure) with respect to the **counting
measure** on the support ([Wikipedia: Probability mass function](https://en.wikipedia.org/wiki/Probability_mass_function)).

## Probability density function (continuous)

A random variable `X` is **(absolutely) continuous** when there is a non-negative
function `f_X` — the **probability density function** — such that the probability
of any interval is the integral of `f_X` over it:

    P(a ≤ X ≤ b) = ∫_a^b f_X(x) dx.

The defining conditions, again necessary and sufficient for a valid PDF:

1. **Non-negativity:** `f_X(x) ≥ 0` for all `x`.
2. **Normalization:** `∫_{-∞}^{∞} f_X(x) dx = 1` — total area under the curve is 1
   ([Wikipedia: Probability density function](https://en.wikipedia.org/wiki/Probability_density_function);
   [Statistics LibreTexts 4.1: PDFs and CDFs for continuous RVs](https://stats.libretexts.org/Courses/Saint_Mary's_College_Notre_Dame/MATH_345__-_Probability_(Kuter)/4:_Continuous_Random_Variables/4.1:_Probability_Density_Functions_(PDFs)_and_Cumulative_Distribution_Functions_(CDFs)_for_Continuous_Random_Variables)).

Three consequences that trip people up constantly:

- **`f_X(x)` is a density, not a probability.** Its units are probability *per unit
  of x*. It can be greater than 1 — e.g. `Uniform[0, 1/2]` has `f(x) = 2` on that
  interval, which is still a valid density because the area is `2 × (1/2) = 1`
  ([Wikipedia: Probability density function](https://en.wikipedia.org/wiki/Probability_density_function)).
- **`P(X = x) = 0` for every single point.** The integral over a single point is
  zero. Probability comes only from intervals (sets of positive measure).
- **Endpoints do not matter:** because point mass is zero,
  `P(a ≤ X ≤ b) = P(a < X < b)`. (For a PMF the endpoints very much matter.)

Measure-theoretically, the PDF is the Radon–Nikodym derivative of the
distribution of `X` with respect to **Lebesgue measure** — the continuous-case
analogue of the PMF-vs-counting-measure relationship
([Wikipedia: Probability density function](https://en.wikipedia.org/wiki/Probability_density_function)).

## Relationship to the CDF — the unifying object

Every real-valued random variable, discrete or continuous, has a **cumulative
distribution function** `F_X(x) = P(X ≤ x)`. The CDF is the bridge:

- **Continuous:** `F_X(x) = ∫_{-∞}^{x} f_X(t) dt`, and by the Fundamental Theorem
  of Calculus the PDF is the derivative of the CDF wherever `F_X` is
  differentiable: `f_X(x) = d/dx F_X(x)`
  ([The Book of Statistical Proofs: PDF is first derivative of CDF](https://statproofbook.github.io/P/pdf-cdf.html);
  [Wikipedia: Probability density function](https://en.wikipedia.org/wiki/Probability_density_function)).
- **Discrete:** `F_X(x) = Σ_{t ≤ x} p_X(t)` — a step function that jumps by
  `p_X(x)` at each support point. The PMF is recovered from the jump sizes.

This is why, when you are unsure whether a variable is discrete, continuous, or
mixed, you reason with the CDF: it always exists and always determines the
distribution.

## Getting a PMF or PDF from the CDF (worked direction)

If you are handed `F_X` and want the density: differentiate. If `F_X` is a step
function: read off the jump heights as the PMF. The CDF is non-decreasing,
right-continuous, with `F_X(-∞) = 0` and `F_X(+∞) = 1`; any candidate "CDF"
violating these cannot correspond to a valid PMF/PDF.

## Change of variables — the density of Y = g(X)

A PMF transforms trivially: if `Y = g(X)` with `g` defined on the support, then
`p_Y(y) = Σ_{x : g(x) = y} p_X(x)` — just collect the mass landing on each `y`.

A PDF does **not** transform that simply, because stretching or compressing the
axis changes how density is packed. For a **monotonic, differentiable,
one-to-one** `g` with inverse `g⁻¹`:

    f_Y(y) = f_X(g⁻¹(y)) · |d/dy g⁻¹(y)|.

The factor `|d/dy g⁻¹(y)|` (the one-dimensional **Jacobian**) rescales the density
so that `f_Y` still integrates to 1 — it accounts for how `g` locally stretches or
shrinks the axis ([Penn State STAT 414, 23.1: Change-of-variables technique](https://online.stat.psu.edu/stat414/lesson/23/23.1);
[Statistics LibreTexts 3.7: Transformations of random variables](https://stats.libretexts.org/Bookshelves/Probability_Theory/Probability_Mathematical_Statistics_and_Stochastic_Processes_(Siegrist)/03:_Distributions/3.07:_Transformations_of_Random_Variables)).

In `d` dimensions the derivative becomes the absolute value of the determinant of
the Jacobian matrix of `g⁻¹`:

    f_Y(y) = f_X(g⁻¹(y)) · |det J_{g⁻¹}(y)|

([Penn State STAT 414, 23.1](https://online.stat.psu.edu/stat414/lesson/23/23.1)).

If `g` is **not** one-to-one (e.g. `Y = X²`), split the domain into monotonic
pieces and sum the contributions, or — the safe fallback — derive `F_Y` directly
and differentiate (the "CDF method").

## Expectation via the PMF / PDF

The PMF/PDF is what you weight by to compute expectations:

- Discrete:  `E[X] = Σ_x x · p_X(x)`,  `E[g(X)] = Σ_x g(x) · p_X(x)`
- Continuous: `E[X] = ∫ x · f_X(x) dx`,  `E[g(X)] = ∫ g(x) · f_X(x) dx`

(The transformed forms are LOTUS; the expectation/variance machinery as a topic
belongs to `da-1-3-1-random-variables` — here the point is only that the mass or
density function is the weight.)

## Mixed distributions — not everything is one or the other

Not every random variable is purely discrete or purely continuous. By **Lebesgue's
decomposition theorem**, any probability distribution on ℝ splits uniquely into a
discrete part, an absolutely continuous part, and a singular-continuous part
([Wikipedia: Lebesgue's decomposition theorem](https://en.wikipedia.org/wiki/Lebesgue's_decomposition_theorem);
[Wikipedia: Random variable](https://en.wikipedia.org/wiki/Random_variable)).

A **mixed** variable has both jumps and a density — e.g. a sensor reading censored
at 0 (positive mass at exactly 0, a density above it), or the mixture
`½ N(0,1) + ½ δ₀`, which is half an absolutely continuous normal density and half
a point mass at 0 ([Wikipedia: Random variable](https://en.wikipedia.org/wiki/Random_variable)).
Such a variable has *neither* a clean PMF *nor* a clean PDF; describe it with the
CDF (which has both jumps and a continuous part). The **singular-continuous** case
(e.g. the Cantor distribution) is continuous yet has *no* density at all — its CDF
rises without any jumps but its derivative is 0 almost everywhere
([Wikipedia: Singular distribution](https://en.wikipedia.org/wiki/Singular_distribution)).

## Common pitfalls (checklist)

- **Reading a PDF value as a probability.** `f_X(x)` is a density; it can exceed 1.
  Only `∫_a^b f_X` is a probability.
- **Forgetting `P(X = x) = 0` for continuous `X`.** A density gives interval
  probabilities, never point probabilities; endpoint inclusion is irrelevant.
- **Skipping the validity check.** A candidate PMF must be `≥ 0` and *sum* to 1; a
  candidate PDF must be `≥ 0` and *integrate* to 1. Verify before using it.
- **Transforming a PDF without the Jacobian.** `f_Y(y) ≠ f_X(g⁻¹(y))`; you must
  multiply by `|d/dy g⁻¹(y)|` (or `|det J|` in higher dimensions).
- **Applying the change-of-variables formula to a non-monotonic `g`.** Split into
  monotonic pieces and sum, or use the CDF method.
- **Assuming every variable has a PMF or a PDF.** Mixed and singular distributions
  have neither in clean form — fall back to the CDF.
- **Differentiating a CDF blindly at a jump.** At a discrete jump the derivative
  does not exist; the jump height is the PMF value there.

## Worked examples

**Discrete (valid PMF check).** A loaded coin: `p_X(0) = 0.3`, `p_X(1) = 0.7`.
Both values `≥ 0` and `0.3 + 0.7 = 1`, so it is a valid PMF. `P(X = 1) = 0.7`
directly — no integration.

**Continuous (density > 1 is fine).** `X ~ Uniform[0, 1/2]` has `f_X(x) = 2` for
`0 ≤ x ≤ 1/2` and 0 elsewhere. Check: `∫_0^{1/2} 2 dx = 1`. Valid, even though the
density value is 2 ([Wikipedia: Probability density function](https://en.wikipedia.org/wiki/Probability_density_function)).

**CDF → PDF.** If `F_X(x) = 1 − e^{−λx}` for `x ≥ 0`, then
`f_X(x) = d/dx F_X(x) = λ e^{−λx}` for `x ≥ 0` — recovering a density by
differentiation ([statproofbook: PDF is first derivative of CDF](https://statproofbook.github.io/P/pdf-cdf.html)).

**Change of variables.** Let `X ~ Uniform[0,1]` (so `f_X(x) = 1` on `[0,1]`) and
`Y = −ln X`. Then `X = e^{−Y} = g⁻¹(Y)`, `|d/dy g⁻¹(y)| = e^{−y}`, so
`f_Y(y) = f_X(e^{−y}) · e^{−y} = e^{−y}` for `y ≥ 0` — the exponential density,
confirmed by the Jacobian factor ([Penn State STAT 414, 23.1](https://online.stat.psu.edu/stat414/lesson/23/23.1)).

## Sources

1. Wikipedia — Probability mass function. https://en.wikipedia.org/wiki/Probability_mass_function
2. Wikipedia — Probability density function. https://en.wikipedia.org/wiki/Probability_density_function
3. Wikipedia — Random variable (types; mixed distributions). https://en.wikipedia.org/wiki/Random_variable
4. Wikipedia — Lebesgue's decomposition theorem. https://en.wikipedia.org/wiki/Lebesgue's_decomposition_theorem
5. Wikipedia — Singular distribution. https://en.wikipedia.org/wiki/Singular_distribution
6. StatLect — Probability density function (interpretation, density vs probability). https://www.statlect.com/glossary/probability-density-function
7. Statistics LibreTexts 4.1 — PDFs and CDFs for Continuous Random Variables. https://stats.libretexts.org/Courses/Saint_Mary's_College_Notre_Dame/MATH_345__-_Probability_(Kuter)/4:_Continuous_Random_Variables/4.1:_Probability_Density_Functions_(PDFs)_and_Cumulative_Distribution_Functions_(CDFs)_for_Continuous_Random_Variables
8. Statistics LibreTexts 3.7 — Transformations of Random Variables (Siegrist). https://stats.libretexts.org/Bookshelves/Probability_Theory/Probability_Mathematical_Statistics_and_Stochastic_Processes_(Siegrist)/03:_Distributions/3.07:_Transformations_of_Random_Variables
9. The Book of Statistical Proofs — PDF is the first derivative of the CDF. https://statproofbook.github.io/P/pdf-cdf.html
10. Penn State STAT 414, Lesson 23.1 — Change-of-Variables Technique. https://online.stat.psu.edu/stat414/lesson/23/23.1
11. STAT 350 (Treese) — Discrete Random Variables and Probability Mass Functions. https://treese41528.github.io/STAT350/Website/chapter5/lectures/5-1-discrete-rvs-and-pmfs.html
