<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-2-1-4-ratio` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-2-1-4-ratio
description: >-
  The ratio level of measurement in Stevens' levels-of-measurement hierarchy
  (data-analysis foundations / measurement theory). Covers the definition, the
  absolute/true zero, the multiplicative (similarity) transformation group,
  permissible transformations and statistics (geometric mean, harmonic mean,
  coefficient of variation, all interval-level stats), worked examples (mass,
  length, age, income, Kelvin, counts), and the pitfalls of misclassifying or
  log-transforming ratio data.
  TRIGGER: a question about what a ratio scale is; whether a variable is ratio
  vs interval vs ordinal; why "10 kg is twice 5 kg" is valid but "20 degrees C
  is not twice 10"; which statistics or transformations are legal for ratio
  data; whether the geometric mean / coefficient of variation / multiplicative
  rescaling is valid for a measure; classifying mass, length, age, income,
  counts, duration, or Kelvin by measurement level.
  SKIP: the interval level with an arbitrary zero (defer to da-1-2-1-3-interval);
  ordinal rank-only data (da-1-2-1-2-ordinal); nominal/categorical labels
  (da-1-2-1-1-nominal); the four-level taxonomy as a whole or how to choose
  among them (da-1-2-1-levels-of-measurement); discrete-vs-continuous as a
  separate axis (da-1-2-2-...); the unrelated "ratio" senses elsewhere in the
  taxonomy — financial ratios, odds ratios / risk ratios in inference, the ratio
  of a proportion, aspect ratios, or ratio-scaled rates as a derived quantity —
  none of which this skill covers.
---

# Ratio level of measurement

Scope: the **ratio scale** as the top of S. S. Stevens' four levels of
measurement (nominal → ordinal → interval → ratio), within data-analysis
foundations and measurement theory. This skill is about classifying a variable
as ratio and knowing what you may legally do with it. It is **not** about
financial ratios, odds/risk ratios in statistical inference, or aspect ratios —
those are different uses of the word "ratio".

## Definition

A ratio scale assigns numbers such that the variable is ordered, the distances
between values are equal and meaningful, **and** the zero point is a true zero
that marks the complete absence of the attribute
([Wikipedia, Level of measurement](https://en.wikipedia.org/wiki/Level_of_measurement)).
Because the zero is absolute, the ratio scale lets you estimate "the ratio
between a magnitude of a continuous quantity and a unit of measurement of the
same kind"
([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)). Stevens
introduced this type as the highest of the four in his 1946 *Science* article
"On the Theory of Scales of Measurement."

In one line: ratio = **interval + a true (absolute) zero.**

The defining feature is that true zero. A weight of zero means there is no
weight; a count of zero means none of the thing
([Scribbr, Ratio Scales](https://www.scribbr.com/statistics/ratio-data/);
[CareerFoundry, Levels of Measurement](https://careerfoundry.com/en/blog/data-analytics/data-levels-of-measurement/)).

## What makes it ratio (the diagnostic test)

Ask, in order:

1. Are the values ordered (can you say one is greater)? If no → nominal
   (see `da-1-2-1-1-nominal`).
2. Are equal numeric differences equal in the attribute — is the step from 10
   to 20 the same amount as 20 to 30? If no → ordinal
   (see `da-1-2-1-2-ordinal`).
3. Does 0 mean "none of the attribute," so that ratios make sense (is 20 twice
   10)? If **no** (zero is a convention) → interval (see `da-1-2-1-3-interval`).
   If **yes** → **ratio**.

The interval/ratio split turns entirely on whether the zero is true or
arbitrary
([Scribbr](https://www.scribbr.com/statistics/ratio-data/)). The textbook test:
temperature in Celsius is interval (0 °C is the freezing point of water, an
arbitrary convention), but in Kelvin it is ratio (0 K is absolute zero, no
thermal energy), so 200 K really is twice 100 K
([Scribbr](https://www.scribbr.com/statistics/ratio-data/)).

## Mathematical structure: the similarity group

Ratio scales are unique only up to a **similarity transformation** `y = a·x`
with `a > 0` — multiplication by a positive constant, with **no additive term**
([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)). You may freely
change the **unit** (`a`, the scale) but **not** the **origin**: the zero is
fixed because it is physically meaningful.

- Converting kilograms to pounds (`lb = 2.205·kg`) or metres to feet is exactly
  this. The unit changes; zero stays at zero (no mass is no mass in either
  unit).
- Contrast the interval scale, which admits the full affine map `y = a·x + b`
  with a free additive `b` (Celsius ↔ Fahrenheit shifts the origin) — covered in
  `da-1-2-1-3-interval`. Dropping `b` is precisely what gives the ratio scale
  its extra power.
- Because there is no additive shift, **ratios of raw values are preserved**:
  if `x₂ = 2·x₁`, then `y₂ = a·x₂ = 2·(a·x₁) = 2·y₁` in any unit. This is what
  makes "twice as much" / "half as much" statements valid
  ([Scribbr](https://www.scribbr.com/statistics/ratio-data/)).

## Permissible statistics (per Stevens)

Ratio data supports **every operation legal at the interval level, plus those
that need a true zero**
([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)):

- **Central tendency:** mode, median, arithmetic mean — **and** the
  **geometric mean** and **harmonic mean**, which require a meaningful zero/origin
  and become legal only at the ratio level
  ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).
- **Dispersion:** range, standard deviation, variance — **and** the
  **coefficient of variation** (SD ÷ mean), which is meaningful here because the
  mean's location is fixed by the true zero (it is *not* legal at interval level)
  ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement);
  [Scribbr](https://www.scribbr.com/statistics/ratio-data/)).
- **Association and inference:** all interval-level methods carry over —
  Pearson correlation, t-tests, ANOVA, linear regression
  ([Scribbr](https://www.scribbr.com/statistics/ratio-data/)).
- **Logarithms** are meaningful because ratios are meaningful, so log transforms
  and multiplicative models are well defined on ratio data.
- Both **moments about the origin** and **central moments** are meaningful,
  because the origin itself is meaningful — unlike interval data, where only
  central moments are
  ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).

Practical note: collecting at the ratio level is the most informative choice,
because you can always step down to a coarser level (treat it as interval,
ordinal, or nominal) but never up
([Scribbr](https://www.scribbr.com/statistics/ratio-data/)).

## Worked examples

- **Mass / weight (kg, lb), length (m, ft), duration (seconds).** Canonical
  ratio quantities: 0 marks true absence, and 10 kg is genuinely twice 5 kg
  ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement);
  [CareerFoundry](https://careerfoundry.com/en/blog/data-analytics/data-levels-of-measurement/)).
- **Temperature in Kelvin** — ratio, because 0 K is absolute zero (no thermal
  energy); 400 K is twice 200 K. The same physical temperature in Celsius or
  Fahrenheit is only interval — see `da-1-2-1-3-interval`
  ([Scribbr](https://www.scribbr.com/statistics/ratio-data/)).
- **Age, income, height, population, energy, electric charge, plane angle.**
  All have a true zero and admit ratios
  ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement);
  [Scribbr](https://www.scribbr.com/statistics/ratio-data/)).
- **Counts (number of children, vehicles, clients, reaction time).** Zero
  means none, and "4 children is twice 2" is valid — counts are ratio (and also
  discrete; the discrete-vs-continuous distinction is a separate axis, see
  `da-1-2-2-...`)
  ([Scribbr](https://www.scribbr.com/statistics/ratio-data/);
  [CareerFoundry](https://careerfoundry.com/en/blog/data-analytics/data-levels-of-measurement/)).

## Common pitfalls

- **Confusing an arbitrary zero with a true zero.** This is the single most
  common misclassification: Celsius temperature, calendar years, IQ, and pH have
  a zero value but not a *true* zero, so they are interval, not ratio. The
  presence of a "0" on the scale is not enough; 0 must mean "none of the
  attribute"
  ([Scribbr](https://www.scribbr.com/statistics/ratio-data/)).
- **Treating a derived ratio as a ratio-scale measurement.** Percentages,
  rates, and indices computed *from* data are not automatically ratio-level in
  the Stevens sense; classify by whether the underlying scale has a true zero.
- **Forgetting that a multiplicative rescale must keep zero fixed.** A valid
  ratio-scale transform is `y = a·x` only; adding a constant (`y = a·x + b`)
  destroys the ratio property and demotes the data to interval.
- **Over-reaching with the geometric mean on data containing zeros or
  negatives.** The geometric mean and log transforms need strictly positive
  values; a true zero is the floor of a ratio scale, and an actual 0 makes the
  geometric mean 0 and the log undefined — handle zeros explicitly.
- **Assuming ratio status licenses any model.** The scale level tells you which
  statistics are *meaningful*, not whether a particular model fits; ratio data
  can still be skewed, bounded, or non-normal.

## Caveat on the typology

Stevens' four-level scheme is the standard teaching frame but has been
challenged. Measurement theorists reject Stevens' broad definition of
measurement — R. Duncan Luce noted that "no measurement theorist I know accepts
Stevens's broad definition" — and the rigid "permissible statistics" rule is
disputed
([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)). The ratio
level is the least contested of the four because physical ratio quantities have
firm measurement-theoretic footing, but the framework as a whole remains debated.
For the broader debate and how the four levels relate, see
`da-1-2-1-levels-of-measurement`.

## Sources

- [Wikipedia — Level of measurement](https://en.wikipedia.org/wiki/Level_of_measurement)
  — Stevens' 1946 typology, similarity group `y = ax`, true/absolute zero,
  permissible statistics (geometric/harmonic mean, coefficient of variation),
  moments-about-the-origin rule, criticisms (Luce).
- [Scribbr — Ratio Scales: Definition, Examples, & Data Analysis](https://www.scribbr.com/statistics/ratio-data/)
  — definition, true-zero vs interval, Kelvin/Celsius example, applicable
  descriptive and inferential statistics, "ratio is most precise" guidance.
- [CareerFoundry — 4 Levels of Measurement: Nominal, Ordinal, Interval & Ratio](https://careerfoundry.com/en/blog/data-analytics/data-levels-of-measurement/)
  — ratio = interval + true zero, weight example, supported operations.
- [Fiveable — Ratio Level of Measurement (AP Statistics)](https://fiveable.me/key-terms/ap-stats/ratio-level-of-measurement)
  — highest level inheriting all lower properties, true-zero definition,
  "twice as much" ratio statements, examples (height, weight, age, income).
