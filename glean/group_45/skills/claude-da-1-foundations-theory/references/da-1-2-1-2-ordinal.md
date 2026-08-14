<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-2-1-2-ordinal` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-2-1-2-ordinal
description: >-
  The ordinal level of measurement in Stevens' measurement theory — ranked
  categories with order but unknown, unequal spacing between them. Covers the
  defining order axiom, permissible statistics (median, mode, percentiles),
  permissible monotonic transformations, the rank-based / nonparametric methods
  ordinal data licenses, the Likert-scale "ordinalist vs intervalist" debate,
  and the central pitfall of treating ordinal data as interval. TRIGGER: a
  variable is ranked but its gaps are not equal (satisfaction ratings, Likert
  items, education level, pain scales, competition placings, severity tiers);
  questions about which central-tendency or test is "allowed" for ranked
  categorical data; "is it ok to average a Likert scale"; choosing nominal vs
  ordinal vs interval vs ratio for a ranked field. SKIP: nominal/unordered
  categories (defer to da-1-2-1-1-nominal); equal-interval scales with arbitrary
  zero (defer to da-1-2-1-3-interval); true-zero ratio scales (defer to
  da-1-2-1-4-ratio); the levels-of-measurement overview itself (defer to
  da-1-2-1-levels-of-measurement); ordinal regression model mechanics, ranking
  algorithms, or database ORDER BY / sort-order topics outside measurement
  theory.
---

# Ordinal level of measurement

**Taxonomy location:** Data Analysis > Foundations & Theory > Measurement theory > Levels of measurement > Ordinal.

This skill covers the **ordinal scale** as one of S. S. Stevens' four levels of
measurement (nominal, ordinal, interval, ratio). Scope is the *measurement-theory*
meaning: what mathematical structure an ordinal scale carries, and therefore which
descriptive statistics and inferential tests are admissible. The word "ordinal"
shows up elsewhere (ordinal numbers in math, ordinal regression as a model family,
SQL sort order); those framings are out of scope here and are deferred in SKIP.

## Definition

An ordinal scale assigns observations to **ordered categories**. The categories
have a meaningful rank — one is "more" or "less" than another on the underlying
attribute — but the **distance between adjacent categories is unknown and not
assumed equal** ([Ordinal data, Wikipedia](https://en.wikipedia.org/wiki/Ordinal_data)).

Formally it is the second of Stevens' four scale types. Stevens (1946) defined a
scale by the group of transformations that leave its structure invariant: an
ordinal scale is invariant under any **monotonically (order-preserving) increasing
transformation** — you may relabel the ranks 1,2,3 as 10,20,300 or as
A<B<C and lose no information, because only the *order* carries meaning
([Level of measurement, Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement);
[Sage / Stevens "On the Theory of Scales of Measurement"](https://methods.sagepub.com/ency/edvol/encyc-of-research-design/chpt/on-theory-scales-measurement)).

What ordinal **adds** over nominal: an ordering relation (>, <). What it still
**lacks** versus interval: equal, meaningful differences. Because the gaps are not
defined, sums and differences of the codes are not interpretable
([Statistics By Jim — Nominal, Ordinal, Interval, Ratio](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/)).

## Worked intuition: why the gaps are not equal

Stevens by Jim's race example: in a foot race, 1st, 2nd, 3rd are evenly *ranked*,
but the actual time between 1st and 2nd may be 0.1 s while 2nd to 3rd is 8 s. The
rank spacing (1, 2, 3) tells you the order, not the magnitude of the gaps
([Statistics By Jim](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/)).

Likewise, a health rating poor < fair < good < excellent is ordered, but the
"distance" from good to excellent need not equal the distance from poor to fair
([Ordinal data, Wikipedia](https://en.wikipedia.org/wiki/Ordinal_data)).

Common ordinal variables: education level, income *brackets* (not income itself),
satisfaction / Likert ratings, pain scales (0–10 self-report), military rank,
severity tiers (P1/P2/P3), competition placings
([Statistics By Jim](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/)).

## Permissible statistics

Stevens tied each scale to the statistics whose value is invariant under the
scale's admissible transformations. For ordinal data:

| Quantity | Permitted? | Notes |
|---|---|---|
| Mode | Yes | Inherited from nominal. |
| **Median** | Yes | The canonical center for ordinal data — the middle-ranked value, well-defined under any order-preserving relabeling ([Statistics By Jim](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/)). |
| Percentiles, quartiles, deciles | Yes | Positional measures of spread ([Sage / Stevens](https://methods.sagepub.com/ency/edvol/encyc-of-research-design/chpt/on-theory-scales-measurement)). |
| Interquartile range, range | Yes | Position-based dispersion. |
| **Mean** | No (per Stevens) | Requires equal intervals; not invariant under monotonic transforms. |
| **Standard deviation / variance** | No (per Stevens) | Same reason as the mean. |
| Count / frequency table | Yes | Inherited from nominal. |

Stevens (1946) argued that, because equal spacing does not hold, means and standard
deviations — and any inference built on them — are not appropriate; use positional
measures (median, percentiles) plus the nominal-level statistics (counts, mode)
([Level of measurement, Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).

## Appropriate inferential methods

Ordinal structure licenses **rank-based (nonparametric)** procedures, which use
only the ordering of values:

- **Association / correlation:** Spearman's rank correlation (rho), Kendall's tau,
  gamma, Somers' D
  ([Ordinal data, Wikipedia](https://en.wikipedia.org/wiki/Ordinal_data)).
- **Two-group comparison:** Mann–Whitney U (Wilcoxon rank-sum); Wilcoxon
  signed-rank for paired data
  ([Ordinal data, Wikipedia](https://en.wikipedia.org/wiki/Ordinal_data)).
- **3+ groups:** Kruskal–Wallis H; Friedman test for repeated measures
  ([Ordinal data, Wikipedia](https://en.wikipedia.org/wiki/Ordinal_data)).
- **Modeling an ordinal outcome:** ordinal (ordered) logistic regression, e.g. the
  proportional-odds model, accommodates ranked outcomes directly
  ([Statistics By Jim](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/)).

These tests compare ranks/medians rather than means, so they stay valid when the
mean is undefined for the scale.

## The Likert-scale debate (ordinalist vs intervalist)

A single Likert *item* (e.g. 1 = strongly disagree … 5 = strongly agree) is the
textbook example of ordinal data ([Ordinal data, Wikipedia](https://en.wikipedia.org/wiki/Ordinal_data)).
Whether you may treat it as interval — and so compute means and run t-tests/ANOVA —
is contested:

- **Ordinalist position:** the points are not proven equidistant, so use
  nonparametric tests (Mann–Whitney, Kruskal–Wallis, Spearman) and report medians
  ([ResearchGate discussion](https://www.researchgate.net/post/What_is_the_most_suitable_statistical_test_for_ordinal_data_eg_Likert_scales)).
- **Intervalist / pragmatic position:** simulation studies report that for
  comparing two samples of 5-point Likert items, the t-test and Mann–Whitney U
  give similar power and similar Type I error, so parametric tests are often
  defensible — especially on **summed multi-item scales** rather than single items
  ([Can I use parametric analyses for my Likert scales? — Lindeløv](https://lindeloev.net/can-i-use-parametric-analyses-for-my-likert-scales-a-brief-reading-guide/);
  [Statistics By Jim — Analyze Likert Scale Data](https://statisticsbyjim.com/hypothesis-testing/analyze-likert-scale-data/)).

Practical convention: treat a single Likert *item* as ordinal (median + bar
chart + nonparametric test); a summed/averaged multi-item *Likert scale* is more
often analyzed as interval, but state the assumption explicitly. (Confidence:
Medium — the literature genuinely disagrees; report both framings.)

## Pitfalls

1. **Averaging ranks.** Reporting a "mean satisfaction of 3.7" assumes equal
   spacing the scale does not provide; the number is sensitive to how categories
   were coded ([Statistics By Jim](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/)).
2. **Treating ordinal as interval by default.** This violates measurement theory
   and risks invalid conclusions; do it only with explicit justification
   ([Ordinal data, Wikipedia](https://en.wikipedia.org/wiki/Ordinal_data)).
3. **Reading equal numeric codes as equal distances.** Codes 1–5 are arbitrary
   labels; only their order is informative
   ([Level of measurement, Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement)).
4. **Power trade-off.** Nonparametric tests can be less powerful than parametric
   ones, which is part of why intervalists push back — weigh strictness against
   detection power ([Statistics By Jim — nonparametric vs parametric](https://statisticsbyjim.com/hypothesis-testing/nonparametric-parametric-tests/)).
5. **Visualizing with a misleading center.** Prefer bar charts of frequencies /
   stacked-percent diverging plots over anything implying a continuous mean.

## Quick decision aid

- Categories ordered but gaps unknown/unequal → **ordinal** (this skill).
- Categories with no order → nominal (defer: `da-1-2-1-1-nominal`).
- Equal intervals, arbitrary zero (e.g. °C, calendar year) → interval
  (defer: `da-1-2-1-3-interval`).
- Equal intervals + true zero, ratios meaningful (e.g. mass, count, duration) →
  ratio (defer: `da-1-2-1-4-ratio`).

## Sources

1. [Ordinal data — Wikipedia](https://en.wikipedia.org/wiki/Ordinal_data) — definition, permissible measures, rank-based methods, Likert framing, interval pitfall.
2. [Level of measurement — Wikipedia](https://en.wikipedia.org/wiki/Level_of_measurement) — Stevens' four scales, monotonic-transform invariance, permissible statistics by scale.
3. [Stevens, "On the Theory of Scales of Measurement" (1946) — Sage Encyclopedia entry](https://methods.sagepub.com/ency/edvol/encyc-of-research-design/chpt/on-theory-scales-measurement) — original typology and permissible-statistics rule.
4. [Statistics By Jim — Nominal, Ordinal, Interval, Ratio Scales](https://statisticsbyjim.com/basics/nominal-ordinal-interval-ratio-scales/) — examples, why gaps are unequal, median as valid center, ordinal logistic regression.
5. [Statistics By Jim — How to Analyze Likert Scale Data](https://statisticsbyjim.com/hypothesis-testing/analyze-likert-scale-data/) — analysis guidance and the parametric/nonparametric choice.
6. [Lindeløv — Can I use parametric analyses for my Likert scales?](https://lindeloev.net/can-i-use-parametric-analyses-for-my-likert-scales-a-brief-reading-guide/) — evidence on item vs scale and the intervalist case.
7. [ResearchGate — Most suitable statistical test for ordinal (Likert) data](https://www.researchgate.net/post/What_is_the_most_suitable_statistical_test_for_ordinal_data_eg_Likert_scales) — ordinalist-vs-intervalist debate.
