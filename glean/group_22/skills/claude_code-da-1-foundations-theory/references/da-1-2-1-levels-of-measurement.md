<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-2-1-levels-of-measurement` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-2-1-levels-of-measurement
description: >-
  Stevens' levels (scales) of measurement — nominal, ordinal, interval, ratio —
  as the measurement-theory foundation for data analysis: classifying a
  variable's scale type, deciding which arithmetic operations and summary
  statistics are permissible, and avoiding scale-mismatch errors before
  analysis. TRIGGER: "what level/scale of measurement is this variable",
  "nominal vs ordinal vs interval vs ratio", "can I take the mean of this
  variable", "which statistic is valid for ordinal/Likert data", "true zero",
  "permissible statistics for a scale", "is temperature interval or ratio",
  classifying columns by measurement level before choosing a summary or test.
  SKIP: choosing a specific statistical hypothesis test or model (defer to the
  statistical-inference / test-selection skills); database/programming column
  data types such as INT/VARCHAR/categorical dtype (a storage concern, not a
  measurement-theory concern); reliability, validity, or measurement error
  (a separate measurement-theory sub-topic); feature encoding for ML
  (one-hot/ordinal encoding belongs to the modeling skills).
---

# Levels of Measurement (Nominal, Ordinal, Interval, Ratio)

Scope: this skill covers measurement-scale classification as the foundational
theory step in **Data Analysis** — specifically *Foundations & Theory >
Measurement theory > Levels of measurement*. It tells you what kind of variable
you are holding and, from that, what arithmetic and which summary statistics are
legitimate. It does **not** select hypothesis tests, design measurement
instruments, or handle storage data types — those are separate concerns (see
SKIP cues above).

## What it is

Psychologist Stanley Smith Stevens introduced the four-level typology in his
1946 *Science* article "On the theory of scales of measurement." Each level is
defined by which mathematical operations remain meaningful, and that in turn
fixes which statistics you may compute ([Wikipedia: Level of measurement](https://en.wikipedia.org/wiki/Level_of_measurement)).
The levels are cumulative: each higher level supports every operation of the
levels below it, plus one more.

| Level | Distinguishing property | Valid operations | Example |
|-------|------------------------|------------------|---------|
| Nominal | Labels only, no order | Equality (=, ≠) | Blood type, eye color, country |
| Ordinal | Order, but unequal/unknown gaps | + ordering (<, >) | Satisfaction rating, education level, pain 1–10 |
| Interval | Equal gaps, no true zero | + addition/subtraction | Celsius/Fahrenheit temperature, SAT score, calendar year |
| Ratio | Equal gaps + true zero | + multiplication/division | Weight, height, duration, Kelvin, count |

Sources agree on this hierarchy of valid operations ([Statology](https://www.statology.org/levels-of-measurement-nominal-ordinal-interval-and-ratio/); [Statistics By Jim](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/)).

## The four levels in detail

### Nominal
Distinct categories with no inherent ordering; categories are mutually
exclusive. You can only test whether two values are the same or different
([Statology](https://www.statology.org/levels-of-measurement-nominal-ordinal-interval-and-ratio/)).
Examples: gender, nationality, blood type, marital status ([Statistics By Jim](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/)).
- Central tendency: **mode only**.
- Summaries: counts, frequencies, proportions; bar/pie charts.
- You cannot compute a mean, median, or standard deviation, because you only
  have category membership ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).

A binary/dichotomous variable (yes/no) is a special nominal case.

### Ordinal
Values have a natural rank order, but the distance between adjacent ranks is
not constant or not known. You know one value is higher or lower than another,
not by how much ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).
Examples: Likert survey responses, education level, socioeconomic status,
pain-intensity ratings ([Statistics By Jim](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/)).
- Central tendency: **median and mode**.
- Spread: interquartile range, percentiles.
- The mean is contested at this level — see the parametric debate below.

### Interval
Ordered values with meaningful, equal differences between them, but **no true
zero** — zero is an arbitrary point on the scale, not the absence of the
quantity. So differences are meaningful but ratios are not: 20 °C is not "twice
as hot" as 10 °C ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).
Examples: Celsius/Fahrenheit temperature, calendar dates, SAT scores, credit
scores ([Statology](https://www.statology.org/levels-of-measurement-nominal-ordinal-interval-and-ratio/)).
- Valid operations: equality, ordering, addition, subtraction.
- Central tendency and spread: mean, median, mode, standard deviation, range.

### Ratio
Everything interval scales have, plus a **true (absolute) zero** that represents
the genuine absence of the quantity. This makes ratios meaningful: 200 lbs is
twice 100 lbs ([Statology](https://www.statology.org/levels-of-measurement-nominal-ordinal-interval-and-ratio/)).
Examples: weight, height, length, duration, counts, energy, Kelvin temperature
([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).
- Valid operations: all four arithmetic operations.
- All statistics apply, including geometric mean, harmonic mean, and the
  coefficient of variation, which require a true zero ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).

## Permissible statistics — quick reference

| Statistic | Nominal | Ordinal | Interval | Ratio |
|-----------|:---:|:---:|:---:|:---:|
| Mode | ✓ | ✓ | ✓ | ✓ |
| Median | ✗ | ✓ | ✓ | ✓ |
| Mean (arithmetic) | ✗ | contested | ✓ | ✓ |
| Standard deviation / range | ✗ | ✗ | ✓ | ✓ |
| Geometric/harmonic mean, coefficient of variation | ✗ | ✗ | ✗ | ✓ |

(Compiled from [Statistics By Jim](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/) and [Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement).)

## How to classify a variable (decision procedure)

Ask, in order:
1. Do the values just name categories with no order? → **Nominal**.
2. Is there a meaningful order, but are the gaps between values uneven or
   unknown? → **Ordinal**.
3. Are the gaps equal and meaningful, but is zero arbitrary (negatives possible,
   or zero does not mean "none")? → **Interval**.
4. Are the gaps equal *and* is zero a true absence (so ratios make sense)? →
   **Ratio**.

The fastest discriminators: ordering separates nominal from the rest; a true
zero separates ratio from interval ([Statology](https://www.statology.org/levels-of-measurement-nominal-ordinal-interval-and-ratio/)).

## Pitfalls

- **Numeric codes for categories.** Encoding "1 = red, 2 = blue, 3 = green"
  does not make the variable ordinal or interval; it is still nominal. Averaging
  such codes is meaningless ([Statistics By Jim](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/)).
- **Interval mistaken for ratio.** Temperature in °C/°F and calendar years have
  no true zero — do not take ratios of them ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).
- **Collecting at too coarse a level.** You cannot recover precision later: if
  you store only income *bands*, you cannot reconstruct exact income. Capture
  the highest level the situation allows ([Statistics By Jim](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/)).
- **Treating the typology as absolute law.** Stevens' rules are widely taught
  but not universally accepted (see next section); use them as a default guide,
  not an unbreakable constraint.

## The parametric-on-ordinal debate

Many statisticians hold that the mean is inappropriate for ordinal data because
the gaps between ranks are not equal. In practice, behavioral researchers
routinely apply means and standard deviations to ordinal (e.g., Likert) data,
arguing that when the unknown gaps between ranks are roughly the same order of
magnitude, interval-style statistics are usable ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).
Tools such as SPSS prompt you to declare each variable's measurement level to
steer you away from inappropriate analyses, but the restriction "lacks universal
acceptance among statisticians" ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).
Report this as a genuine controversy rather than picking a side: state which
convention you are following and why.

## Criticism and alternative typologies

Stevens' four-level scheme has been challenged repeatedly ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)):
- **Velleman & Wilkinson (1993)** argued the framework can mislead, especially
  for the nominal and ordinal categories, and that scale type should not
  mechanically dictate analysis.
- **Mosteller & Tukey (1977)** proposed seven categories instead of four, adding
  types such as "grades" (ordered labels like beginner/advanced) and "counted
  fractions" (bounded between 0 and 1).
- **Chrisman (1998)** expanded to ten levels, adding cyclical-ratio (angles,
  times of day), log-interval (e.g., stock charts), and graded membership
  (fuzzy sets).

These do not replace Stevens in everyday data analysis, but knowing them keeps
you from over-applying the four-level rules to data they fit poorly (cyclic
data, bounded proportions, fuzzy membership).

## Worked micro-examples

- *"Customer tier: bronze/silver/gold"* → ordinal (ordered, gaps not numeric);
  report mode/median, not mean.
- *"Response time in milliseconds"* → ratio (true zero, ratios meaningful);
  any statistic, including coefficient of variation.
- *"Survey rating 1–5"* → ordinal by theory; if you summarize with a mean,
  declare you are following the common applied convention and note the caveat.
- *"Year a record was created (2019, 2024…)"* → interval (differences valid,
  but year 0 is arbitrary; do not say one year is "twice" another).

## Sources

1. [Level of measurement — Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement) — Stevens' typology, valid operations per level, permissible statistics, the parametric-on-ordinal debate, and the Velleman & Wilkinson / Mosteller & Tukey / Chrisman criticisms.
2. [Nominal, Ordinal, Interval, and Ratio Scales — Statistics By Jim](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/) — practical identification, central-tendency/spread per scale, and the "collect at the highest level" pitfall.
3. [Levels of Measurement: Nominal, Ordinal, Interval and Ratio — Statology](https://www.statology.org/levels-of-measurement-nominal-ordinal-interval-and-ratio/) — definitions, valid-operation table (equality/ordering/add-subtract/multiply-divide), and the true-zero interval-vs-ratio distinction.
4. [4 Levels of Measurement — CareerFoundry](https://careerfoundry.com/en/blog/data-analytics/data-levels-of-measurement/) — data-analytics framing of the four levels and their analysis implications.
5. [Stevens' levels of measurement — PSYCTC.org](https://www.psyctc.org/psyctc/glossary2/ratio-scaling/) — origin of the typology and its use/limits in behavioral measurement.
