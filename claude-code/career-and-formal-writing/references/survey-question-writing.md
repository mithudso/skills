<!-- hub-reference-banner -->
> **Reference file — part of the `career-and-formal-writing` hub.** Formerly the standalone `survey-question-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: survey-question-writing
description: Bias-free survey question design — Likert scale construction (5-point vs 7-point, balanced/unbalanced anchors), NPS wording (Reichheld 2003), choosing among CSAT/CES/NPS, detecting double-barreled and leading questions, "Don't Know" vs "Neutral" vs forced-choice tradeoffs, anchor labels, ordering effects, the agree-disagree anti-pattern (Saris & Gallhofer), question-stem hygiene, Dillman's Tailored Design Method, open- vs closed-ended choice, and mobile question-length constraints. TRIGGER when the user is writing or auditing a survey question, picking a rating scale, debating CSAT/CES/NPS, suspects bias/leading wording, asks "is this a good question", "rewrite this survey question", "Likert scale", "NPS question", "double-barreled", "anchor labels", "survey response options", or is building a customer feedback / employee engagement / market research / UX research questionnaire. SKIP for: open-ended interview guides (use interview-and-conversational); marketing landing-page copy or CTAs (use sales-and-marketing-copy); support-ticket / case-write-up wording (use support-ticket-writing); executive emails or memos (use executive-comms); inclusive-language audits (use inclusive-language); general prose editing (use writing-expert).
---

# Survey Question Writing

## Overview

A survey question is a measurement instrument. Bad wording does not just irritate respondents — it injects measurement error that downstream statistics cannot fix. This skill applies the consensus rules of survey methodology (Dillman's Tailored Design Method, Saris & Gallhofer's question-quality framework, Pew Research and Reichheld) to write questions that produce reliable, comparable, low-bias data.

Default mental model: every question is a hypothesis about what the respondent will read. If two thoughtful readers could parse the stem differently, the question is broken — fix the stem before tuning the scale.

## When not to use

- Writing an open-ended interview discussion guide → use `interview-and-conversational`.
- Drafting marketing copy, landing pages, ads, or CTAs → use `sales-and-marketing-copy`.
- Writing a support-ticket reply or case note → use `support-ticket-writing`.
- Composing an executive email/memo → use `executive-comms`.
- Auditing prose for inclusive language → use `inclusive-language`.
- General prose editing without a measurement/instrument goal → use `writing-expert`.
- Choosing a survey *platform* (Qualtrics vs SurveyMonkey vs Google Forms) — out of scope; this skill covers question design, not tooling.

## Core Concepts

### 1. Question-stem hygiene (Dillman)

Stems must be **direct, concrete, mutually exclusive, and answerable in one read**. Test each stem against:

- **One concept per question** — no "and/or" coordination across distinct constructs.
- **Common vocabulary** — no jargon, acronyms, or terms a respondent might need to look up.
- **Concrete time window** — "in the last 30 days" not "recently".
- **Specified reference** — "your most recent purchase" not "purchases".
- **Symmetric framing** — avoid "do you agree that X is good" (loaded). Prefer "How would you rate X?" with a balanced scale.

### 2. Likert scales: 5-point vs 7-point

| Decision | Use 5-point | Use 7-point |
|---|---|---|
| Mobile / short pulse | Yes | Avoid |
| Need granularity (academic, MTMM) | No | Yes |
| Comparing to existing 5-pt benchmarks | Yes | No |
| Construct has fine gradations (attitude intensity) | Limited | Yes |
| Time-pressured respondents | Yes | No |

**Default: 5-point** for operational customer/employee surveys. **7-point** when you need discrimination for regression/factor analysis. Even-numbered (4/6-point) scales force a side and increase response error from undecided respondents; only use them when neutrality is meaningless for the construct.

### 3. Balanced anchors and label coverage

A balanced scale has the same number of positive and negative points around a neutral midpoint, with **labels on every point** (not just endpoints). Fully labeled scales are more reliable than endpoint-only scales (Saris & Gallhofer). Anchors must be roughly equidistant in intensity:

- Good: Strongly disagree / Disagree / Neither / Agree / Strongly agree
- Bad: Hate it / Dislike / Neutral / Like / Love (asymmetric intensity)

Avoid mixing frequency and agreement scales in the same matrix.

### 4. NPS, CSAT, and CES — pick one per question

| Metric | Question wording | Scale | Use when |
|---|---|---|---|
| **NPS** (Reichheld 2003) | "How likely is it that you would recommend [company/product] to a friend or colleague?" | 0–10 | Loyalty, top-line growth proxy, brand-level |
| **CSAT** | "How satisfied were you with [specific experience]?" | 1–5 or 1–7 | Post-interaction satisfaction, transactional |
| **CES** | "[Company] made it easy to handle my issue." (agree/disagree 1–7) | 1–7 | Service/support friction, predicting churn from effort |

Note: canonical CES uses an agree-disagree stem, which conflicts with the agree-disagree anti-pattern below. Keep the canonical wording if you need benchmark comparability; otherwise use an item-specific rewrite: "How much effort did it take to resolve your issue?" (Very low effort ... Very high effort).

Do not stack all three on one survey without a reason — they correlate and inflate length. NPS scoring: 0–6 = Detractors, 7–8 = Passives, 9–10 = Promoters. **Do not modify the NPS stem** if you want to compare to industry benchmarks.

### 5. The agree-disagree anti-pattern (Saris & Gallhofer 2014)

"Do you agree or disagree that [statement]?" is consistently lower quality than **item-specific response options** because it invites acquiescence bias (respondents lean toward "agree") and combines two cognitive steps (evaluate statement, then map to agreement).

- Anti-pattern: "Do you agree that the website is easy to use?" (5-pt agree/disagree)
- Better: "How easy or difficult is the website to use?" (Very difficult ... Very easy)

Item-specific scales reduce systematic error and increase reliability across constructs.

### 6. Double-barreled and leading questions

**Double-barreled** = one question, two concepts. Split it.

- Bad: "How satisfied are you with the price and quality of the product?"
- Fix: Two questions, one for price, one for quality.

**Leading** = stem prejudges the answer or uses charged language. Pew Research found "welfare" vs "assistance to the poor" shifted support by 20+ percentage points on otherwise identical questions.

- Bad: "How helpful was our amazing support team?"
- Fix: "How would you rate the support you received?" (Very poor ... Very good)

### 7. "Don't Know" vs "Neutral" vs forced-choice

- **Neutral midpoint** = respondent has an opinion but it is centered. Keep on attitude/agreement scales.
- **"Don't Know" / "Not applicable"** = respondent lacks the information to answer. Offer it when the question requires knowledge the respondent may not have; **place it visually offset from the scale** (Dillman) so satisficers do not select it by default.
- **Forced-choice** (no neutral, no DK) = use only when you genuinely need a side and respondents will plausibly have one (e.g., binary preference between two named options).

Conflating these inflates error: a "neutral" click from a respondent who actually doesn't know is unrecoverable noise.

### 8. Ordering effects

- **Question order**: early questions can prime later ones. Put sensitive/demographic questions last. Put the headline metric (NPS/CSAT) before drill-downs so it is not contaminated by reflection.
- **Response option order**: on visual scales, primacy bias (top) dominates; on aural, recency dominates. Randomize unordered option lists across respondents. Keep ordered scales (Likert, NPS) in their natural order — never randomize a Likert.
- **Matrix straight-lining**: long batteries of similar Likert items invite straight-lining. Break up with item-specific scales or shorter batteries.

### 9. Open- vs closed-ended

- **Closed-ended** = the default for quantitative analysis, benchmarking, large samples.
- **Open-ended** = use sparingly (1–2 per survey) to capture unanticipated themes or quotes. Place after the closed question that primed the topic ("You rated this 2/5 — what was the main reason?"). Mobile respondents abandon long open-text fields; cap at one short prompt for mobile.

### 10. Mobile constraints (Dillman 4th ed., tailored to mobile)

- Stem ≤ 20 words; aim for ≤ 12.
- Scale ≤ 7 points (5 preferred) so radio buttons fit horizontally on a phone.
- Avoid grids/matrices on mobile — collapse to single-question-per-screen.
- Total survey ≤ 5 minutes (≈ 10 questions) for transactional touchpoints; ≤ 12 minutes for relational/annual surveys.
- Never require pinch-to-zoom or horizontal scroll to read a question.

## Templates and Examples

### NPS (canonical, do not modify the stem)

> How likely is it that you would recommend [Company/Product] to a friend or colleague?
> 0 (Not at all likely) — 10 (Extremely likely)

Optional follow-up (open-ended, conditional on score):
> What is the primary reason for your score?

### CSAT (transactional)

> Thinking about your [recent support interaction / recent purchase / today's visit], how satisfied were you overall?
> 1 (Very dissatisfied) — 2 — 3 (Neither) — 4 — 5 (Very satisfied)

### CES (customer effort, support-flavored)

> [Company] made it easy for me to resolve my issue.
> Strongly disagree — Disagree — Somewhat disagree — Neither — Somewhat agree — Agree — Strongly agree

### Generic attribute rating (item-specific, preferred over agree-disagree)

> How would you rate the [speed / clarity / accuracy] of the response you received?
> Very poor — Poor — Fair — Good — Excellent

### Frequency

> In the last 30 days, how often have you used [feature]?
> Never — Once — 2–3 times — About once a week — Several times a week — Daily

### Rewrite examples

| Original (broken) | Rewrite | Why |
|---|---|---|
| "Do you agree that our new pricing is fair and competitive?" | "How would you rate our pricing?" (Very unfair ... Very fair) + separate "How would you rate our pricing vs. competitors?" | Agree-disagree + double-barreled |
| "How helpful was our amazing support team?" | "How would you rate the support you received?" (Very poor ... Very good) | Leading adjective |
| "Recently, how often have you used the dashboard?" | "In the last 30 days, how many times did you open the dashboard? (0, 1–3, 4–10, 11+)" | Vague time + vague frequency |
| "Rate your satisfaction with price and quality" | Two separate item-specific questions | Double-barreled |
| "Don't you think the new feature is useful?" | "How useful is the new feature?" (Not at all useful ... Extremely useful) | Leading + negation |

## Anti-Patterns

- **Agree-disagree everything.** Causes acquiescence bias. Replace with item-specific scales.
- **Double-barreled stems.** "Easy and fast", "price and quality", "domestic and foreign". Split.
- **Loaded adjectives.** "amazing", "innovative", "world-class", "fair". Strip them from stems.
- **Negation in stems.** "How much do you disagree with..." — confuses respondents. Reframe positively.
- **Endpoint-only labels on 7-pt scales.** Less reliable than full labels.
- **Random Likert order.** Never randomize ordered response options.
- **Hidden DK / forced midpoint.** Either causes data loss or noise. Decide explicitly per question.
- **Long matrix batteries on mobile.** Invites straight-lining; collapse to single-screen items.
- **Stacking NPS + CSAT + CES + 10 drill-downs.** Length kills response rates. Pick a primary metric.
- **Modifying the NPS stem** (e.g., "to a colleague" only, dropping "friend"). Breaks benchmark comparability.
- **Asking demographics first.** Primes responses, increases abandonment. Put them last.

## Decision Heuristics

1. **One sentence, one concept, one scale.** If you cannot say the stem in one breath, split it.
2. **Read the stem aloud.** If you stumble on a word or want to add a clause, the respondent will too.
3. **Default to 5-point fully-labeled item-specific.** Move off the default only with a stated reason (granularity, benchmark, construct nature).
4. **Pick a primary metric per touchpoint.** Transactional → CSAT or CES. Relational/brand → NPS. Effort-focused → CES.
5. **Demographics last, headline metric first** within the survey body.
6. **Mobile-test before launch.** Open the survey on a phone; if you have to scroll horizontally, redesign.
7. **Pre-test with 5 respondents using think-aloud.** Cheapest bug-fix in survey research; Dillman's recommended floor.
8. **If a stem includes "and", "or", "but" between concepts → split.**
9. **If a stem includes an adjective describing the thing being rated → strip the adjective.**
10. **Keep DK separated visually from the substantive scale.** Don't let satisficers grab it by accident.

## References

1. Dillman, D. A., Smyth, J. D., & Christian, L. M. (2014). *Internet, Phone, Mail, and Mixed-Mode Surveys: The Tailored Design Method* (4th ed.). Wiley. https://www.wiley.com/en-us/Internet,+Phone,+Mail,+and+Mixed+Mode+Surveys:+The+Tailored+Design+Method,+4th+Edition-p-9781118456149
2. Saris, W. E., & Gallhofer, I. N. (2014). *Design, Evaluation, and Analysis of Questionnaires for Survey Research* (2nd ed.). Wiley. Comparative work on agree-disagree vs item-specific scales: https://ojs.ub.uni-konstanz.de/srm/article/view/2682 and PMC review: https://pmc.ncbi.nlm.nih.gov/articles/PMC8692311/
3. Reichheld, F. F. (2003). "The One Number You Need to Grow." *Harvard Business Review*. https://hbr.org/2003/12/the-one-number-you-need-to-grow — and Reichheld, Darnell & Burns (2021), "Net Promoter 3.0." https://hbr.org/2021/11/net-promoter-3-0
4. Pew Research Center. "Writing Survey Questions." https://www.pewresearch.org/writing-survey-questions/
5. Taufique, K. M. R. (2024). "A Guide to Key Decision Criteria for Likert-Scale Use in Survey Research." *Global Business and Organizational Excellence*. https://onlinelibrary.wiley.com/doi/10.1002/joe.70032
