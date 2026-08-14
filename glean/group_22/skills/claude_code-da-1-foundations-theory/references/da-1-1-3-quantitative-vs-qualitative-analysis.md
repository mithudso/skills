<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-1-3-quantitative-vs-qualitative-analysis` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-1-3-quantitative-vs-qualitative-analysis
description: >-
  Foundational distinction between quantitative and qualitative analysis as a
  definitional concept in data analysis — what each is, the kind of question
  each answers, the data each works on, the analytic procedures each uses
  (statistics vs. coding/thematic interpretation), their trade-offs, and how
  mixed methods combine them. TRIGGER: someone asks "what is the difference
  between quantitative and qualitative analysis", "is this data qualitative or
  quantitative", "should I use a qualitative or quantitative approach", "numbers
  vs words analysis", "deductive vs inductive analysis", "objective vs
  subjective data", or needs to classify an analysis approach before choosing a
  method. SKIP: choosing a specific statistical test or running regression
  (defer to a statistics/inferential-analysis skill); designing a survey
  instrument or interview guide (defer to research-design/survey skills);
  running a specific qualitative coding pass such as grounded theory or
  discourse analysis end to end (defer to a qualitative-methods skill); levels
  of measurement nominal/ordinal/interval/ratio (defer to
  da-1-2-1-levels-of-measurement); the broader analysis-vs-synthesis distinction
  (defer to da-1-1-2-analysis-vs-synthesis).
---

# Quantitative vs. Qualitative Analysis

Scope: Data Analysis > Foundations & Theory > Definitions & scope >
Quantitative vs. qualitative analysis.

This skill covers the *definitional* boundary between the two families of
analysis — the framing you use to decide which family a problem belongs to.
It does not teach a specific statistical test or a specific coding method;
those are downstream choices that follow once the family is settled.

## The core distinction

- **Quantitative analysis** works on numerical data — anything that can be
  counted or measured and assigned a numeric value. It uses mathematics and
  statistics to find patterns, test theories, and produce generalizable facts.
  It answers questions like *what*, *how many*, *how often*, and *to what
  degree* [Scribbr; SimplyPsychology].
- **Qualitative analysis** works on non-numerical data — words, transcripts,
  images, observations, diaries. It interprets how people perceive and assign
  meaning to their experience. It answers questions like *why* and *how*,
  exploring motivations and reasons [Scribbr; SimplyPsychology].

A one-line test: *if the answer you want is a number or a measured pattern, it
is quantitative; if the answer you want is an interpretation of meaning, it is
qualitative* [SMU Learning Sciences].

## When to use which

| Use quantitative when… | Use qualitative when… |
|---|---|
| Confirming or testing a hypothesis/theory | Exploring an idea, meaning, or experience |
| You need to generalize to a population | You need depth on a poorly understood topic |
| The question is "how many / how much / how often" | The question is "why / how" |
| You can collect large, structured samples | You have small samples and rich detail |

Reasoning direction differs too: quantitative work is typically **deductive**
(start from a theory, gather data to test it), while qualitative work is
typically **inductive** (start from observations, build theory from them)
[Scribbr; SimplyPsychology].

## Data and collection

- **Quantitative data / methods:** survey instruments with closed questions,
  controlled experiments, and observations recorded as numbers; test scores,
  rating scales, frequency counts, sensor/log metrics [Scribbr;
  SimplyPsychology; SMU].
- **Qualitative data / methods:** open-ended interviews, focus groups,
  ethnographic/participant observation, document and literature review;
  transcripts, field notes, diaries, visual materials [Scribbr;
  SimplyPsychology].

## How the analysis actually proceeds

**Quantitative analytic workflow** [SMU; SimplyPsychology]:
1. Define the purpose and the analysis parameters (what counts, consistency
   rules).
2. Organize data in spreadsheets or a database.
3. Clean the data — remove blanks, duplicates, and irrelevant responses.
4. Apply statistics:
   - *Descriptive statistics* summarize the data (means, distributions,
     frequencies).
   - *Inferential statistics* test significance and generalize from a sample to
     a population (correlation, regression, hypothesis tests).
5. Tools: spreadsheets, statistical software (e.g., SPSS), or general analytics
   tooling.

**Qualitative analytic workflow** [SMU; Scribbr; SimplyPsychology]:
1. Review the data and record initial observations.
2. Organize by category (date, location, source, collection method).
3. **Code** the data — assign labels to segments and group them into themes
   that address the research question.
4. Interpret the themes and draw conclusions about recurring patterns and
   meaning.

Named qualitative techniques you will encounter (each is its own downstream
method — this skill only names them so you can route correctly):
- **Thematic analysis** — identify repeating patterns/themes [Scribbr;
  SimplyPsychology].
- **Content analysis** — categorize text and track word/concept meaning
  [Scribbr; SimplyPsychology].
- **Grounded theory** — build theory inductively from the data
  [SimplyPsychology].
- **Discourse analysis** — examine language in its social context [Scribbr;
  SimplyPsychology].

## Trade-offs

**Quantitative** [SimplyPsychology; Scribbr]
- Strengths: scientific objectivity, efficient analysis, easy to replicate,
  supports broad generalization, theory validation.
- Limitations: needs large samples; can feel artificial/decontextualized;
  vulnerable to statistical misapplication, sampling error, and confirmation
  bias.

**Qualitative** [SimplyPsychology; Scribbr]
- Strengths: deep, contextual understanding; surfaces new relationships;
  captures complexity that numbers miss.
- Limitations: small samples limit generalization; labor-intensive and slow;
  vulnerable to researcher/observer bias, social-desirability effects, and
  limited replicability.

Reliability and validity are framed differently across the two: quantitative
work prioritizes measurable consistency and generalizability; qualitative work
emphasizes trustworthiness and context-specific insight [SimplyPsychology].

## Mixed methods

The two families are complementary, not exclusive. A common mixed-methods
sequence: run qualitative interviews first to surface themes, then build a
quantitative survey to test those themes at scale — or pair quantitative
performance metrics with qualitative context for a fuller picture [Scribbr;
SMU]. Use this when you need both the breadth of numbers and the depth of
meaning.

## Common pitfalls

- **Misclassifying the data.** Numeric-looking codes (e.g., ZIP codes, category
  IDs) are labels, not measured quantities; counting them as quantitative
  invites meaningless statistics. (Where a numeric value sits on the
  nominal/ordinal/interval/ratio scale is a separate concept — see
  da-1-2-1-levels-of-measurement.)
- **Forcing one family's standards onto the other** — judging a small,
  meaning-rich qualitative study for "small sample size", or judging a survey
  for "lacking rich narrative".
- **Treating qualitative analysis as quick.** Coding and interpretation are
  typically slower and more labor-intensive than running statistics on
  structured data [Scribbr].
- **Confirmation bias in both directions** — selecting quotes that fit a prior
  belief, or running tests until something is "significant".

## Sources

1. Scribbr — "Qualitative vs. Quantitative Research | Differences, Examples &
   Methods." https://www.scribbr.com/methodology/qualitative-quantitative-research/
2. SimplyPsychology — "Qualitative vs Quantitative Research: What's the
   Difference?" https://www.simplypsychology.org/qualitative-quantitative.html
3. SMU Learning Sciences — "Qualitative vs. quantitative data analysis: How do
   they differ?" https://learningsciences.smu.edu/blog/qualitative-vs-quantitative-data-analysis
