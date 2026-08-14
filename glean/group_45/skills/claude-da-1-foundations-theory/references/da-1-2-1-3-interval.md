<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-2-1-3-interval` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-2-1-3-interval
description: >-
  The interval level of measurement in Stevens' levels-of-measurement hierarchy
  (data-analysis foundations / measurement theory). Covers the definition,
  the affine structure (equal intervals, arbitrary zero), permissible
  transformations and statistics, worked examples (Celsius/Fahrenheit, calendar
  dates, IQ), and the pitfalls that come from an arbitrary zero (ratios are
  meaningless).
  TRIGGER: a question about what an interval scale is; whether a variable is
  interval vs ordinal vs ratio; why "20 degrees is not twice as hot as 10";
  which statistics or transformations are legal for interval data; whether the
  mean/SD/Pearson r are valid for a measure; classifying temperature, dates, or
  IQ scores by measurement level.
  SKIP: the ratio level with a true zero (defer to da-1-2-1-4-ratio); ordinal
  rank-only data (da-1-2-1-2-ordinal); nominal/categorical labels
  (da-1-2-1-1-nominal); the four-level taxonomy as a whole or how to choose
  among them (da-1-2-1-levels-of-measurement); discrete-vs-continuous as a
  separate axis (da-1-2-2-...); the unrelated "interval" senses elsewhere in the
  taxonomy — confidence intervals, prediction intervals, class intervals/bins in
  a histogram, or time intervals in scheduling — none of which this skill covers.
---

# Interval level of measurement

Scope: the **interval scale** as one of S. S. Stevens' four levels of
measurement (nominal → ordinal → interval → ratio), within data-analysis
foundations and measurement theory. This skill is about classifying a variable
as interval and knowing what you may legally do with it. It is **not** about
confidence/prediction intervals, histogram class intervals (bins), or
time-scheduling intervals — those are different uses of the word "interval".

## Definition

An interval scale assigns numbers such that the variable is ordered **and** the
distances (intervals) between adjacent values are equal and meaningful, but the
zero point is arbitrary rather than a true absence of the attribute
([Wikipedia, Level of measurement](https://en.wikipedia.org/wiki/Level_of_measurement)).
Stevens introduced this type in his 1946 *Science* article "On the Theory of
Scales of Measurement," where the interval type "permits defining the degree of
difference between measurements, but not the ratio between measurements"
([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).

In one line: interval = **ordinal + equal intervals, minus a true zero.**

## What makes it interval (the diagnostic test)

Ask, in order:

1. Are the values ordered (can you say one is greater)? If no → nominal
   (see `da-1-2-1-1-nominal`).
2. Are equal numeric differences equal in the attribute — is the step from 10
   to 20 the same amount as 20 to 30? If no → ordinal
   (see `da-1-2-1-2-ordinal`).
3. Does 0 mean "none of the attribute," so that ratios make sense (is 20 twice
   10)? If **yes** → ratio (see `da-1-2-1-4-ratio`). If **no** (zero is a
   convention) → **interval**.

The interval/ratio split turns entirely on whether the zero is true or
arbitrary
([UNSW Online, Types of Data and Scales of Measurement](https://studyonline.unsw.edu.au/blog/types-of-data)).

## Mathematical structure: the affine group

Interval scales are unique only up to a **positive affine (linear)
transformation** `y = a·x + b` with `a > 0`
([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)). You may freely
change the **unit** (`a`, the scale) and the **origin** (`b`, the zero) without
losing any information the scale carries.

- Celsius → Fahrenheit is exactly this: `F = 1.8·C + 32`. Both are valid
  interval scales of the same underlying temperature; neither zero is physical.
- Because the origin is free, **ratios of raw values are not preserved**
  (`a·x + b` does not scale proportionally), which is why "twice as hot" is
  meaningless.
- Because the unit and origin cancel out of differences,
  **ratios of differences ARE preserved and meaningful** — e.g. "the gap from
  10° to 30° is twice the gap from 10° to 20°" holds in both Celsius and
  Fahrenheit
  ([Rasch.org, Measurement Theory: Fallacies and Transformations](https://www.rasch.org/rmt/rmt81g.htm)).

Contrast: ratio scales admit only `y = a·x` (no additive `b`), which is what
preserves ratios — covered in `da-1-2-1-4-ratio`.

## Permissible statistics (per Stevens)

Because differences are meaningful, interval data supports more than ordinal
data does
([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)):

- **Central tendency:** mode, median, **and arithmetic mean** (the mean becomes
  legal at the interval level — it is not strictly legal at ordinal).
- **Dispersion:** range and **standard deviation**.
- **Association:** **Pearson correlation coefficient**.
- Only **central moments** (about the mean) are meaningful; moments about the
  **origin** are arbitrary because the origin itself is arbitrary
  ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).
- **Not** permitted: the **coefficient of variation** (SD ÷ mean) and any
  statistic that needs a true zero, because the mean's location depends on the
  arbitrary origin
  ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).

## Worked examples

- **Temperature in Celsius or Fahrenheit** — the textbook interval example.
  30°C − 20°C equals 20°C − 10°C (equal intervals), but 20°C is **not** twice as
  hot as 10°C because 0°C is the freezing point of water, not absence of heat
  ([Fiveable, Stevens' levels of measurement](https://fiveable.me/key-terms/college-intro-stats/stevens-levels-measurement);
  [ScienceDirect Topics, Interval Scale](https://www.sciencedirect.com/topics/computer-science/interval-scale)).
  (Kelvin, by contrast, has a true zero and is ratio — see
  `da-1-2-1-4-ratio`.)
- **Calendar years / dates** measured from an arbitrary epoch (e.g. AD). The
  year 2000 is not "twice" the year 1000; the zero is conventional
  ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).
- **IQ scores and many standardized psychometric scales** are commonly treated
  as interval: differences are taken as comparable, but there is no meaningful
  zero intelligence, so an IQ of 140 is not "twice as smart" as 70
  ([UNSW Online](https://studyonline.unsw.edu.au/blog/types-of-data)).
- **Cartesian location and compass direction (degrees)** — position and
  bearings have arbitrary origins
  ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).

## Common pitfalls

- **Computing ratios of raw values.** "Today (20°C) is twice as warm as
  yesterday (10°C)" is the canonical error. Ratios require a true zero; interval
  scales do not have one
  ([ScienceDirect Topics](https://www.sciencedirect.com/topics/computer-science/interval-scale)).
- **Reporting a coefficient of variation on interval data.** It depends on the
  mean, which moves when you shift the origin (e.g. switching Celsius to
  Fahrenheit changes the CV) — so it is not meaningful
  ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).
- **Confusing arbitrary zero with no zero.** Interval scales can have a zero
  value (0°C exists); it simply does not mark absence of the attribute.
- **Misreading the affine freedom.** A linear rescale (unit + origin) is
  information-preserving for interval data; a nonlinear rescale is not. Only
  monotone-order is preserved at the ordinal level below it.
- **Treating ordinal data as interval without justification.** Behavioural
  researchers routinely apply means and SDs to Likert-type ordinal scores,
  arguing real measures sit "between" ordinal and interval; this is a contested
  practice, not a free pass
  ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).

## Caveat on the typology

Stevens' four-level scheme is the standard teaching frame but has been
challenged. Measurement theorists (e.g. R. Duncan Luce, 1997) reject Stevens'
broad definition of measurement, and the rigid "permissible statistics" rule is
disputed — some argue the quality of an interval measure depends on theoretical
fit and additivity (Rasch modelling) rather than on a transformation-rule
checklist
([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement);
[Rasch.org](https://www.rasch.org/rmt/rmt81g.htm)). For the broader debate and
how the four levels relate, see `da-1-2-1-levels-of-measurement`.

## Sources

- [Wikipedia — Level of measurement](https://en.wikipedia.org/wiki/Level_of_measurement)
  — Stevens' 1946 typology, affine group `y = ax + b`, permissible statistics,
  central-moments rule, criticisms.
- [Rasch.org — Measurement Theory: Fallacies and Transformations](https://www.rasch.org/rmt/rmt81g.htm)
  — interval as additive/linear units, ratios-of-differences, transformation
  debate.
- [ScienceDirect Topics — Interval Scale](https://www.sciencedirect.com/topics/computer-science/interval-scale)
  — definition, arbitrary-zero, ratio-meaninglessness.
- [UNSW Online — Types of Data and Scales of Measurement](https://studyonline.unsw.edu.au/blog/types-of-data)
  — interval vs ratio split on true zero; IQ/temperature examples.
- [Fiveable — Stevens' levels of measurement](https://fiveable.me/key-terms/college-intro-stats/stevens-levels-measurement)
  — temperature worked example, equal-interval property.
