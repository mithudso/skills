<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-2-3-reliability-validity` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-2-3-reliability-validity
description: >-
  Measurement-theory treatment of reliability and validity: whether a measure is
  consistent (reliability) and whether it measures what it claims to (validity),
  grounded in classical test theory (observed = true + error). Covers the four
  reliability types (test-retest, parallel/alternate forms, internal consistency,
  inter-rater) and their coefficients (Cronbach's alpha, McDonald's omega, ICC,
  Cohen's/Fleiss's kappa), and the validity types (face, content, criterion with
  concurrent/predictive, construct with convergent/discriminant), plus the
  reliability-is-necessary-but-not-sufficient-for-validity relationship and the
  standard error of measurement.
  TRIGGER: questions about whether a test/scale/survey/instrument is reliable or
  valid; what kind of reliability or validity to report; interpreting Cronbach's
  alpha / omega / ICC / kappa; test-retest, inter-rater, internal consistency;
  content/criterion/construct/convergent/discriminant/face validity; "is my
  measure consistent" or "does my measure capture the construct"; classical test
  theory true-score / measurement-error framing.
  SKIP: defining a construct or turning it into measurable indicators
  (-> da-1-2-4-operationalization-constructs); choosing the measurement scale
  level nominal/ordinal/interval/ratio (-> da-1-2-1-levels-of-measurement);
  internal/external/statistical-conclusion validity of a STUDY DESIGN or causal
  inference (an experimental-design topic, not measurement-instrument reliability/
  validity); software test reliability, system uptime, or data-pipeline validation
  (engineering, not measurement theory); inter-annotator agreement framed purely
  as an ML-labeling QA metric rather than instrument reliability.
---

# Reliability & Validity (Measurement Theory)

This skill covers reliability and validity **as properties of a measurement
instrument** — a test, scale, survey, rubric, or rating procedure that turns a
construct into numbers. It sits in the taxonomy under Data Analysis > Foundations
& Theory > Measurement theory. The question it answers is: *can I trust the
numbers this instrument produces, and do they mean what I think they mean?*

## Scope note (what this is and is not)

The words "reliability" and "validity" are overloaded. This skill is about the
**psychometric / measurement-theory** sense.

- **In scope:** consistency and accuracy of an *instrument* — Cronbach's alpha on
  a questionnaire, test-retest correlation of a scale, inter-rater agreement on a
  coding rubric, whether a depression inventory actually measures depression.
- **Out of scope (defer):**
  - **Operationalization / construct definition** — deciding *what* to measure and
    *which indicators* represent the construct is the upstream step. Defer to
    `da-1-2-4-operationalization-constructs`.
  - **Levels of measurement** — whether the data are nominal/ordinal/interval/
    ratio. Defer to `da-1-2-1-levels-of-measurement`.
  - **Study-design validity** — internal validity, external validity, and
    statistical-conclusion validity describe whether a *study's causal claim* is
    sound. That is a research-design topic, not an instrument property; do not
    absorb it here.
  - **Engineering "reliability"** — uptime, MTBF, software test stability, data
    validation rules. Unrelated sense of the word.

## 1. The classical test theory (CTT) foundation

Reliability and validity both rest on one model. Classical test theory says every
**observed score** decomposes into a **true score** plus **random measurement
error**:

```
X (observed) = T (true) + E (error)
```

The true score is the value you would get with a perfect, error-free instrument
(conceptually, the average over infinitely many independent administrations).
Error is random noise — fatigue, ambiguous wording, a distracted rater
[MTU psychometrics, pages.mtu.edu/~shanem/psy5220].

- **Reliability** is about the size of `E` relative to `X`. A reliable measure has
  small random error, so repeated measurement reproduces the same number. The
  reliability coefficient is conceptually the proportion of observed-score variance
  that is true-score variance: `reliability = Var(T) / Var(X)`, ranging 0 to 1.
- **Validity** is about whether `T` (the thing you reliably measure) is actually
  the construct you intended, not some other construct or systematic bias.

This is why **reliability is necessary but not sufficient for validity**: a
bathroom scale that always reads 5 kg heavy is perfectly reliable (consistent) but
not valid (wrong). You can be consistently wrong, but you cannot be accurate
without first being consistent
[Scribbr methodology; Test Partnership; SimplyPsychology].

### Standard error of measurement (SEM)

The SEM expresses reliability on the score's own scale — the expected spread of an
individual's observed scores around their true score:

```
SEM = SD * sqrt(1 - reliability)
```

where `SD` is the standard deviation of observed scores. Higher reliability gives a
smaller SEM and tighter confidence bands around an individual score. Report SEM
when interpreting a single person's score, not just the group-level reliability
coefficient [MTU psychometrics].

## 2. Reliability — the four types

Reliability asks: *does the instrument produce consistent results?* Each type holds
a different thing constant and checks for agreement across the rest
[Scribbr; Arctic Shores; Test Partnership].

| Type | What varies | What it detects | Typical statistic |
|---|---|---|---|
| Test-retest | Time | Stability across occasions | Pearson/Spearman r, ICC |
| Parallel / alternate forms | Item set (equivalent forms) | Equivalence of two versions | Correlation between forms |
| Internal consistency | Items within one form | Whether items cohere on one construct | Cronbach's alpha, McDonald's omega |
| Inter-rater | Raters/observers | Agreement between scorers | Cohen's/Fleiss's kappa, ICC |

### 2.1 Test-retest reliability (stability)

Administer the same instrument to the same people on two occasions and correlate
the scores. High correlation means scores are stable over time. **Pitfall:** the
interval matters. Too short and people remember answers (inflates the estimate);
too long and the true construct may genuinely change (deflates it). Only
appropriate for traits expected to be stable, not for transient states like mood
[MTU psychometrics; SimplyPsychology].

### 2.2 Parallel (alternate) forms reliability

Build two versions drawn from the same content domain, give both, and correlate.
Useful when re-using the same items would cause practice effects (e.g., repeated
cognitive testing). **Pitfall:** the two forms must be genuinely equivalent in
difficulty and content, which is hard to construct
[Arctic Shores; Test Partnership].

### 2.3 Internal consistency

The most-reported reliability type because it needs only a single administration.
It asks whether items that are supposed to measure one construct actually correlate
with each other and with the total score.

- **Cronbach's alpha** is the standard coefficient. Rough conventions: > 0.9
  excellent, 0.8-0.9 good, 0.7-0.8 acceptable, < 0.7 questionable. Alpha rises with
  more items and with higher inter-item correlation.
- **Critical caution:** a high alpha does **not** prove the scale is
  unidimensional. Alpha *assumes* a single underlying factor rather than testing
  for it, and it inflates with item count. Use mean inter-item correlation,
  item-total correlations, and factor analysis to check structure; consider
  **McDonald's omega** (which models single-factor saturation) or the Greatest
  Lower Bound as alternatives [MTU psychometrics].

### 2.4 Inter-rater (inter-observer) reliability

When humans score or code, check that different raters produce the same scores on
the same targets.

- For **categorical** judgments, use **Cohen's kappa** (two raters), which corrects
  for chance agreement; **weighted kappa** for ordered categories; **Fleiss's** or
  **Light's kappa** for more than two raters.
- For **continuous** ratings on the same scale, use the **intraclass correlation
  coefficient (ICC)**. ICC has several variants (one-way vs two-way, random vs
  fixed raters, agreement vs consistency); it is easy to specify the model wrong,
  so state which ICC form you used [MTU psychometrics].

## 3. Validity — the main types

Validity asks: *does the instrument measure what it claims to, and do the score
interpretations hold up?* Modern measurement theory treats validity as a unified
judgment about score interpretation supported by multiple lines of evidence, but
the classic categories are still the practical vocabulary
[Scribbr; American Journal of Medicine 2006; TAO Testing].

### 3.1 Face validity

The weakest, most superficial check: does the instrument *look* like it measures
the construct, on the surface, to a lay reader? It is not real evidence of validity
— a test can look relevant and still be invalid — but low face validity can hurt
respondent buy-in [SimplyPsychology].

### 3.2 Content validity

Do the items cover the full breadth of the construct's domain, without gaps or
irrelevant content? Established by expert review against a content specification /
blueprint (e.g., a math test must sample the whole syllabus, not just one topic).
This is judgment-based, not a single coefficient, though indices like the Content
Validity Index (CVI) quantify expert agreement [Scribbr; TAO Testing].

### 3.3 Criterion validity

Do scores correspond to an external **criterion** — a gold-standard measure of the
same thing? Two sub-types by timing:

- **Concurrent validity:** instrument and criterion measured at the *same time*
  (new short depression screen vs full clinical interview today).
- **Predictive validity:** instrument predicts a *future* criterion (entrance exam
  predicting later grades). Quantified by the correlation between predictor and
  criterion [Scribbr; American Journal of Medicine 2006].

### 3.4 Construct validity

The central, overarching form: do the scores behave the way the construct's theory
says they should? Evidence comes in two complementary directions:

- **Convergent validity:** scores correlate strongly with measures of the *same* or
  theoretically related constructs.
- **Discriminant (divergent) validity:** scores correlate *weakly* with measures of
  *unrelated* constructs. An aggression questionnaire should track other aggression
  measures (convergent) but not be confused with assertiveness or social dominance
  (discriminant) [SimplyPsychology; arcticshores aggression example].

A common tool is the multitrait-multimethod matrix (multiple traits x multiple
methods) to show convergent and discriminant patterns simultaneously. Factor
analysis also supports construct validity by confirming the expected dimensional
structure [TAO Testing; MTU psychometrics].

## 4. How reliability and validity relate

- **Reliability bounds validity.** The validity coefficient cannot exceed the
  square root of the reliability — noisy scores attenuate every correlation, so an
  unreliable instrument *cannot* be highly valid. Fix reliability first.
- **A measure can be reliable but not valid** (consistently measuring the wrong
  thing — systematic bias).
- **A measure cannot be valid but unreliable** in practice, because random error
  would prevent consistent, accurate measurement.

The standard mental picture is the dartboard: tight cluster off-center = reliable
but invalid; tight cluster on the bullseye = reliable and valid; scattered shots =
unreliable (and therefore not valid) [Scribbr; SimplyPsychology].

## 5. Practical workflow for evaluating a measure

1. Confirm the construct and its operationalization are settled first
   (`da-1-2-4-operationalization-constructs`) and the scale level is known
   (`da-1-2-1-levels-of-measurement`).
2. **Establish reliability** appropriate to the instrument:
   - Single-administration self-report -> internal consistency (alpha *and* omega,
     plus a dimensionality check).
   - Stable trait re-measured -> test-retest with a defended interval.
   - Human-coded data -> inter-rater (kappa for categories, ICC for continuous).
3. **Accumulate validity evidence**, weakest to strongest: face -> content (expert
   blueprint) -> criterion (concurrent/predictive) -> construct (convergent +
   discriminant, factor structure).
4. Report the **coefficient, its CI, the sample, and the conditions** — reliability
   and validity are properties of *scores in a context*, not fixed properties of a
   test forever. Re-evaluate when the population or use changes.
5. Report **SEM** when interpreting individual scores.

## Common pitfalls

- Treating a high Cronbach's alpha as proof of unidimensionality or of validity.
  It is neither.
- Reporting reliability but never establishing validity (consistent ≠ correct).
- Picking the wrong reliability type — e.g., internal consistency on a
  deliberately heterogeneous index, or test-retest on a transient state.
- Using face validity as if it were evidence.
- Confusing **measurement** validity (this skill) with **study-design** internal/
  external validity (defer — different topic).
- Forgetting that reliability and validity are sample- and context-specific, so
  inherited coefficients from another population may not transfer.

## Sources

- Scribbr, "Reliability vs. Validity in Research" — definitions, four reliability
  types, four validity types, the relationship.
  https://www.scribbr.com/methodology/reliability-vs-validity/
- SimplyPsychology, "Reliability vs Validity in Research" — types, the
  reliable-but-not-valid distinction, construct/convergent/discriminant examples.
  https://www.simplypsychology.org/reliability-or-validity.html
- MTU PSY5220, "Psychometrics: Reliability and Validity" (pages.mtu.edu) —
  classical test theory, ICC variants, Cronbach's alpha cautions, omega/GLB, kappa,
  standard error of measurement.
  https://pages.mtu.edu/~shanem/psy5220/daily/Day05/psychometrics.html
- Arctic Shores glossary, "Reliability in Psychometric Tests" — the four reliability
  types and the extraversion reliable-not-valid example.
  https://www.arcticshores.com/glossary/reliability-in-psychometric-tests
- Test Partnership, "Reliability and Validity" — equivalence/parallel forms, alpha
  conventions, reliability-bounds-validity framing.
  https://www.testpartnership.com/academy/reliability-validity.html
- TAO Testing, "Validity and Reliability: The Core Concepts of Psychometrics" —
  content/criterion/construct validity, CVI, factor analysis as construct evidence.
  https://www.taotesting.com/blog/validity-and-reliability-the-core-concepts-of-psychometrics-in-assessment/
- Cook & Beckman (2006), "Current Concepts in Validity and Reliability for
  Psychometric Instruments," American Journal of Medicine — unified validity,
  criterion sub-types, reliability coefficients.
  https://www.amjmed.com/article/S0002-9343(05)01037-5/fulltext
