<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-2-1-1-nominal` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-2-1-1-nominal
description: >-
  The NOMINAL level of measurement (Stevens 1946) — the lowest scale type, where
  values are unordered category labels that name group membership and nothing
  more. Covers what makes a variable nominal, the only permissible operations
  (equality / set membership), the only valid summaries (mode, frequency counts,
  proportions), valid tests (chi-square, Cramer's V), the binary/dichotomous
  special case, and the errors that follow from treating a category code as a
  number. TRIGGER: "is this variable nominal", "nominal data examples", "can I
  average a category code / zip code / blood type", "what statistic for an
  unordered category", "mode vs mean for categorical labels", "chi-square for two
  categorical variables", "why can't I rank gender / nationality / color",
  "permissible statistics for a nominal variable", binary yes/no as a nominal
  scale. SKIP: the four-level framework as a whole or choosing among the levels
  (defer to da-1-2-1-levels-of-measurement); ORDERED categories like
  small/medium/large or Likert (that is the ordinal level — defer to its skill);
  interval or ratio scales; the broader measurement-theory parent (reliability,
  validity, scaling — defer to da-1-2-measurement-theory); one-hot/categorical
  ENCODING for machine-learning models (a modeling concern); database column
  types such as VARCHAR/ENUM (a storage concern).
---

# Nominal level of measurement

**Where this sits:** Data Analysis > Foundations & Theory > Measurement theory > Levels of measurement > **Nominal**.

This skill covers *only* the nominal level. For the four-level framework as a
whole, or for deciding which level a variable belongs to in the first place,
defer to `da-1-2-1-levels-of-measurement`. For ordered categories (the next
level up), defer to the ordinal skill.

## Definition

A **nominal scale** (from Latin *nomen*, "name") assigns observations to
**distinct, mutually exclusive, exhaustive categories that carry no inherent
order or magnitude**. The category label only names which group an observation
belongs to; it asserts nothing about more/less, bigger/smaller, or distance.
Measuring on a nominal scale is therefore equivalent to **classifying**
([Wikipedia, "Level of measurement"](https://en.wikipedia.org/wiki/Level_of_measurement);
[Scribbr, "Nominal Data"](https://www.scribbr.com/statistics/nominal-data/)).

It is the lowest, least precise, and qualitative pole of the four-level
typology that psychologist **Stanley Smith Stevens** introduced in his 1946
*Science* paper "On the Theory of Scales of Measurement"
([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).

> "Nominal data is labelled into mutually exclusive categories within a
> variable. These categories cannot be ordered in a meaningful way."
> — [Scribbr](https://www.scribbr.com/statistics/nominal-data/)

## The defining test (how to recognize a nominal variable)

Apply these checks; a variable is nominal when the first is yes and the rest are no:

1. **Are the values named categories / labels?** (yes for nominal)
2. **Is there a meaningful order among them?** If yes → it is *ordinal*, not
   nominal (defer to the ordinal skill).
3. **Are the gaps between values equal and meaningful?** If yes → interval/ratio.
4. **Is there a true, meaningful zero?** If yes → ratio.

Two practical tells:

- **The "shuffle test":** if you can reorder the categories in the legend and
  nothing about the data's meaning changes, the variable is nominal. Reordering
  {*car, bus, train*} loses no information; reordering {*low, medium, high*}
  does — so the latter is ordinal.
- **Numerals are just labels.** Jersey numbers, ZIP codes, route numbers,
  and category codes (1 = married, 2 = single) *look* numeric but are nominal:
  the digits identify a category, they do not quantify anything
  ([Fiveable, "Stevens' levels of measurement"](https://fiveable.me/key-terms/college-intro-stats/stevens-levels-measurement);
  [Scribbr](https://www.scribbr.com/statistics/nominal-data/)).

## Properties

- **Identity / distinguishability only.** The single relation the scale
  supports is *equality* (same category) and its complement, *inequality*
  (different category). Stevens' framework states that "equality and other
  operations that can be defined in terms of equality, such as inequality and
  set membership, are the only non-trivial operations that generically apply to
  objects of the nominal type"
  ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).
- **No order**, no rank, no greater-than/less-than.
- **No meaningful distance** between categories and **no true zero**
  ([Scribbr](https://www.scribbr.com/statistics/nominal-data/)).
- **No arithmetic.** Addition, subtraction, multiplication, and division are
  undefined on the codes ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).
- **Permissible transformation (group-theoretic view):** only **one-to-one
  (bijective) relabelings** preserve the scale's information — you may rename or
  recode the categories any way you like, as long as the mapping is invertible,
  and no statistical conclusion changes
  ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)). This
  "admissible transformation" is what formally distinguishes nominal from the
  higher scales (which admit only order-preserving, affine, or similarity
  transforms, respectively).

## Examples

Canonical nominal variables: sex/gender, nationality, ethnicity, biological
species, blood type, political party, religion, eye color, parts of speech
(noun/verb/preposition), transportation mode (car/bus/train), movie genre,
ZIP/postal code, and employment status (employed/unemployed)
([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement);
[Scribbr](https://www.scribbr.com/statistics/nominal-data/);
[Fiveable](https://fiveable.me/key-terms/college-intro-stats/stevens-levels-measurement)).

### Binary / dichotomous nominal

A nominal variable with exactly two categories is **dichotomous** (binary):
yes/no, true/false, success/failure, pass/fail, alive/dead
([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)). It is still
nominal — there is no inherent order between the two labels even when they are
coded 0/1. (Coding as 0/1 is convenient for computation and lets you compute a
*proportion*, but the proportion is a frequency summary, not an average of an
ordered quantity.)

## What you can and cannot do

### Permissible summaries and statistics

- **Central tendency: the mode only** — the most frequently occurring category.
  The mean and median are undefined because they require arithmetic or order
  ([Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement);
  [Scribbr](https://www.scribbr.com/statistics/nominal-data/)).
- **Frequency distributions / counts and proportions** per category.
- **Dispersion / diversity:** because the standard deviation needs a mean, use
  qualitative-variation measures instead — e.g. the **index of qualitative
  variation (IQV)**, Shannon **entropy**, or the **number of distinct
  categories** present.
- **Association between two nominal variables:** the **chi-square test**
  (goodness-of-fit for one variable against expected proportions; test of
  independence for a two-way contingency table), with effect-size measures such
  as **Cramer's V**, the **phi coefficient** (for 2x2), or **Goodman-Kruskal
  lambda** ([Scribbr](https://www.scribbr.com/statistics/nominal-data/);
  [Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).
- **Visualization:** bar chart and pie chart of category frequencies; a
  contingency table (cross-tab) for two variables
  ([Scribbr](https://www.scribbr.com/statistics/nominal-data/)).

### Not permissible

- **Do not average category codes.** The mean of {1 = car, 2 = bus, 3 = train}
  = 2.0 is meaningless; "average transport mode = bus" is nonsense.
- **Do not take a median or rank**, and do not report a standard deviation of
  the codes.
- **Do not apply order-based or parametric procedures** (t-test on the codes,
  Pearson correlation, linear regression treating the code as a continuous
  predictor without encoding, monotone transforms) — these all assume order or
  interval structure the scale does not have
  ([Scribbr](https://www.scribbr.com/statistics/nominal-data/)).

## Worked example

Survey question "Primary commute mode?" with responses coded
`1 = car, 2 = bus, 3 = bike, 4 = train` over 200 respondents:

| Mode  | Code | Count | Proportion |
|-------|------|-------|------------|
| car   | 1    | 90    | 0.45       |
| bus   | 2    | 50    | 0.25       |
| bike  | 3    | 40    | 0.20       |
| train | 4    | 20    | 0.10       |

- **Mode** (most common category) = *car*. This is the only valid measure of
  central tendency.
- Reporting `mean(code) = (90·1 + 50·2 + 40·3 + 20·4) / 200 = 1.95` is an
  **error** — 1.95 names no category and the codes carry no magnitude.
- To test whether commute mode is associated with another category (say, urban
  vs. rural), build a 4x2 contingency table and run a **chi-square test of
  independence**, then report **Cramer's V** for effect size.

## Pitfalls

1. **Numeric-code trap.** A column of integers is *not* evidence of a numeric
   scale. ZIP codes, jersey numbers, and 1/2/3 category codes are nominal; the
   storage type (INT) is a programming concern, not a measurement level. (For
   the storage-vs-measurement distinction, this is exactly the SKIP boundary —
   database column types are out of scope here.)
2. **Silent ordinal creep.** If categories actually have an order
   (small/medium/large, disagree→agree), the variable is **ordinal**, not
   nominal — handle it under the ordinal skill, not here.
3. **Averaging a 0/1 dummy and calling it a "mean."** The proportion of 1s is
   fine as a frequency, but do not interpret it as the mean of an interval
   quantity.
4. **Default-encoding a nominal predictor into a model as a single integer.**
   In analysis you must one-hot/dummy-encode rather than feed the raw code — but
   that *encoding* step itself belongs to the modeling skills (SKIP boundary);
   what matters here is recognizing that the raw code is unordered.

## A caveat from the literature (do not over-apply the rules)

Stevens' typology, and especially the strict "permissible statistics" rule, has
been criticized. **Velleman and Wilkinson (1993)**, "Nominal, Ordinal,
Interval, and Ratio Typologies are Misleading" (*The American Statistician*
47(1): 65-72), argue the typology "ignores important developments in data
analysis" and that mechanically forbidding operations by scale type can mislead
([Velleman & Wilkinson 1993, abstract](https://www.tandfonline.com/doi/abs/10.1080/00031305.1993.10475938)).
The nominal case is the least contested part of their critique — for genuinely
unordered categories there is broad agreement that order- and arithmetic-based
summaries are inappropriate. The lesson for the nominal level is narrow:
**classify the variable correctly first**, and once it is genuinely nominal,
restrict yourself to mode, counts/proportions, and frequency-based tests.

## Sources

1. [Wikipedia — "Level of measurement"](https://en.wikipedia.org/wiki/Level_of_measurement)
   — definition, Stevens 1946 origin, permissible operations (equality / set
   membership), bijective transformation, mode-only central tendency, no
   arithmetic, dichotomous case, examples.
2. [Scribbr — "Nominal Data: Definition, Examples, and Analysis"](https://www.scribbr.com/statistics/nominal-data/)
   — definition, mutually-exclusive unordered categories, examples, mode /
   frequency / chi-square permissible, mean and median impermissible.
3. [Fiveable — "Stevens' levels of measurement" (Intro to Statistics)](https://fiveable.me/key-terms/college-intro-stats/stevens-levels-measurement)
   — categories as labels with no order; numerals-as-labels (jersey numbers,
   blood type); never average nominal values.
4. [Velleman, P. F. & Wilkinson, L. (1993), "Nominal, Ordinal, Interval, and Ratio Typologies are Misleading," *The American Statistician* 47(1): 65-72](https://www.tandfonline.com/doi/abs/10.1080/00031305.1993.10475938)
   — critique of the typology and its rigid permissible-statistics rule.
