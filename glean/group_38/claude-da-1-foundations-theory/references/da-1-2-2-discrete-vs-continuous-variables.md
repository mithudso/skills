<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-2-2-discrete-vs-continuous-variables` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-2-2-discrete-vs-continuous-variables
description: >-
  Distinguishes discrete from continuous variables as a measurement-theory
  property of quantitative data (Data Analysis > Foundations & Theory >
  Measurement theory). Covers the countable-vs-uncountable definition, the
  count/measure test, PMF vs PDF consequences, and when count data may be
  treated as continuous. TRIGGER: "is this variable discrete or continuous",
  "discrete vs continuous", "countable vs measurable data", "can I treat count
  data as continuous", "why does my count variable break a linear model",
  "PMF or PDF for this variable", classifying a quantitative column before
  choosing a model or distribution. SKIP: nominal/ordinal/interval/ratio scale
  classification — that is the *levels of measurement* axis, defer to
  da-1-2-1-levels-of-measurement and its children (da-1-2-1-1-nominal,
  da-1-2-1-2-ordinal, da-1-2-1-3-interval, da-1-2-1-4-ratio); the broader
  parent overview (da-1-2-measurement-theory); probability-distribution
  mechanics or random-variable theory (da-1-3-probability-theory,
  da-1-3-1-random-variables); discreteness of a programming/database data type
  rather than a measured variable.
---

# Discrete vs. Continuous Variables

A measurement-theory distinction for **quantitative** variables: whether the
values are *countable* (discrete) or fill a *range of real numbers*
(continuous). This is a different axis from the level of measurement
(nominal/ordinal/interval/ratio); a quantitative variable has both a level
*and* a discrete/continuous character. See `da-1-2-1-levels-of-measurement` for
the orthogonal scale axis.

## Definitions

**Discrete variable.** Takes one of a set of separated values, with no
meaningful values in between. Formally, a variable is discrete when there is a
one-to-one correspondence between its possible values and a subset of the
natural numbers — i.e. the values are countable (finite or countably infinite)
and there is a positive minimum gap between adjacent permissible values
[Wikipedia, *Continuous or discrete variable*]. You **count** discrete data
[Statistics By Jim]. Example: number of Facebook friends — "it doesn't make
sense to say that one has 33.7 friends" [LibreTexts / Poldrack, *Statistical
Thinking for the 21st Century*, §2.2].

**Continuous variable.** Defined in terms of a real number and can fall
anywhere in a range; between any two values there exist further possible values
[Wikipedia]. You **measure** continuous data [Statistics By Jim]. Example:
weight, height, time, temperature. Measurement tools limit the recorded
precision, but the underlying quantity is still continuous [LibreTexts §2.2].

### The fast test
- Can you **count** it in whole units with nothing valid in between? → discrete.
- Can you **measure** it, with finer values always conceivable between two
  readings? → continuous.
[Statistics By Jim; Statistics How To]

## Why it matters

The discrete/continuous character drives which probability model and which
analysis are appropriate — it is not cosmetic.

- **Probability representation.** Discrete variables are described by a
  **probability mass function (PMF)** — probability sits on individual points.
  Continuous variables use a **probability density function (PDF)** — the
  probability of any single exact value is zero; probability lives in intervals
  [Wikipedia]. (Distribution mechanics themselves belong to
  `da-1-3-probability-theory` / `da-1-3-1-random-variables`.)
- **Mathematical tooling.** Continuous variables admit calculus (integration
  over a density, derivatives); discrete variables use sums and combinatorial /
  integer methods [Wikipedia].
- **Modeling.** Continuous outcomes fit naturally into ordinary linear
  regression; discrete counts often need Poisson or negative-binomial models,
  and binary outcomes need logistic regression [Wikipedia; The Analysis Factor,
  *When Can Count Data be Considered Continuous?*].

## Worked classification examples

| Variable | Discrete or continuous | Why |
|---|---|---|
| Number of cars in a lot | Discrete | Countable; no 20.5 cars [Statistics By Jim] |
| Number of traffic accidents in a day | Discrete | Count, non-negative integers [LibreTexts §2.2] |
| Weight of a package | Continuous | Measured; arbitrarily fine values exist [LibreTexts §2.2] |
| Time to complete a task | Continuous | Real-valued duration [Wikipedia] |
| Temperature reading | Continuous | Real number in a range [G2; Statistics By Jim] |
| Shoe size (half sizes only) | Discrete | Separated permissible values |

## Pitfalls

1. **"Numeric format = continuous" is false.** A numeric code can still be a
   label, and a numeric column can still be discrete. The format does not tell
   you the type — what the value *represents* does [LibreTexts §2.2 zip-code
   caution]. Deciding the scale (label vs. magnitude) is the
   *levels of measurement* question — defer to
   `da-1-2-1-levels-of-measurement`.

2. **Count data is discrete and bounded at zero.** Counts are non-negative
   integers, which breaks the assumptions of an ordinary linear model that
   expect an unbounded continuous response [The Analysis Factor]. Reach for a
   count model (Poisson / negative binomial) rather than forcing a continuous
   method.

3. **…but high-valued counts can be *treated* as continuous.** When counts are
   large enough that the distribution is effectively normal and the spacing
   between integers is negligible relative to the range, the loss from treating
   them as continuous is unnoticeable [The Analysis Factor]. This is a pragmatic
   modeling choice, not a change in the variable's true nature.

4. **Discrete-for-analysis vs. continuous-for-analysis is a separate, practical
   convention.** Analysts often treat variables with few categories (nominal,
   ordinal) as discrete and variables with many values (interval, ratio) as
   continuous, purely to pick statistical methods [Statistics How To]. Keep this
   *handling* convention distinct from the *intrinsic* count-vs-measure
   property defined above.

5. **Precision-limited continuous data still counts as continuous.** Recording
   weight only to the nearest pound does not make weight discrete; the
   instrument truncates an underlying continuous quantity [LibreTexts §2.2].

6. **Mixed variables exist.** Some quantities have both a discrete lump and a
   continuous part — e.g. queue wait time has a point mass at exactly zero plus
   a continuous distribution over positive waits [Wikipedia]. Don't assume every
   variable is purely one or the other.

## Scope boundary

This skill answers *is the variable countable or measurable, and what follows
from that*. It does **not** classify the measurement scale (nominal / ordinal /
interval / ratio) — that orthogonal axis belongs to
`da-1-2-1-levels-of-measurement` and its children. It also does not derive
distribution formulas (`da-1-3-probability-theory`,
`da-1-3-1-random-variables`).

## Sources

- LibreTexts — *Statistical Thinking for the 21st Century* (Poldrack), §2.2
  "Discrete Versus Continuous Measurements":
  https://stats.libretexts.org/Bookshelves/Introductory_Statistics/Statistical_Thinking_for_the_21st_Century_(Poldrack)/02:_Working_with_Data/2.02:_Discrete_Versus_Continuous_Measurements
- Wikipedia — *Continuous or discrete variable*:
  https://en.wikipedia.org/wiki/Continuous_or_discrete_variable
- Statistics By Jim — *Discrete vs. Continuous Data: Differences & Examples*:
  https://statisticsbyjim.com/basics/discrete-vs-continuous-data/
- The Analysis Factor — *When Can Count Data be Considered Continuous?*:
  https://www.theanalysisfactor.com/count-data-considered-continuous/
- Statistics How To — *Discrete vs Continuous variables*:
  https://www.statisticshowto.com/probability-and-statistics/statistics-definitions/discrete-vs-continuous-variables/
