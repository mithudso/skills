<!-- hub-reference-banner -->
> **Reference file — part of the `da-applied-and-communication` hub.** Formerly the standalone `da-9-reporting-communication` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-9-reporting-communication
description: >-
  Reporting and communication as the final data-analysis phase — turning findings into
  decisions via reports, executive summaries, technical write-ups, dashboards, decks, and
  notebooks. Covers data storytelling (Knaflic), BLUF, Minto Pyramid/SCQA, audience
  adaptation (exec/technical/operational), recommendation framing, uncertainty
  disclosure, reproducibility appendices, and Jupyter/Quarto/Observable/R Markdown
  notebooks. TRIGGER: communicate analysis results to non-analysts; structure an exec
  summary, one-pager, or memo; lead with the bottom line; translate one finding per
  reader; recommendations with confidence levels; disclose uncertainty honestly;
  reproducible notebooks; avoid burying the lede or AI-generated tells. SKIP:
  chart-only design → da-applied-and-communication (da-8-data-visualization); general
  prose → writing-expert; exec comms without analysis → executive-comms; slide design
  apart from analysis → technical-writing-craft.
---

# Reporting and Communication

**Taxonomy context:** Data Analysis > Reporting and Communication

Reporting and Communication is the phase that converts analytical work into decisions. Every
upstream phase — problem framing, data preparation, modeling, evaluation — is wasted if the
audience cannot understand the finding, cannot trust its rigor, and cannot act on the
recommendation. Reporting is not a write-up of what the analyst did; it is an argument constructed
for a specific decision-maker [Source 1, Source 2].

The mistake most analysts make is structuring the report the way the work happened (problem →
methods → results → conclusions). Decision-makers consume information top-down: they want the
answer first, then the supporting evidence, then the method, then the caveats. The conventions in
this skill all serve that inversion.

---

## 1. The bottom-line-first inversion: BLUF, Minto, SCQA

Three overlapping frameworks all enforce the same discipline: lead with the conclusion, not the
reasoning chain.

### 1.1 BLUF — Bottom Line Up Front

BLUF originated in U.S. military staff writing where readers triage many documents under time
pressure and must extract the operational decision in the first sentence [Source 3]. It has since
moved into executive business writing, technical incident reports, and analytical memos.

The pattern is rigid:

1. **Line 1** — the answer or recommendation, in one sentence.
2. **Line 2-3** — the single most important supporting fact and the action requested.
3. **Body** — evidence, method, alternatives, caveats, in decreasing order of importance.

A BLUF analytic memo looks like:

> Renewal probability for the SMB segment dropped 11 percentage points in Q1 (from 78% to 67%).
> The drop is concentrated in customers acquired through the partner channel in the last 18
> months. Recommend pausing partner-channel acquisition spend until cohort retention recovers;
> details and confidence in the body.

The opposite — opening with "This analysis examines Q1 renewal trends across segments…" — fails
BLUF because the reader has to read three paragraphs to learn whether action is needed.

### 1.2 Minto Pyramid Principle

Developed by Barbara Minto at McKinsey in the 1970s, the Pyramid Principle is the canonical
top-down structure for business writing [Source 4]. It says every document has a single governing
thought at the top, supported by 2-5 grouping ideas one level down, each in turn supported by
specific facts.

```
                  Governing Thought (the answer)
                /              |              \
        Supporting       Supporting        Supporting
        argument 1       argument 2        argument 3
        /     \          /     \           /     \
      fact   fact      fact   fact       fact   fact
```

Three rules govern the pyramid:

- **MECE** — supporting arguments at any level must be Mutually Exclusive and Collectively
  Exhaustive.
- **Same-kind grouping** — arguments at the same level must be of the same logical type (all
  causes, all consequences, all options).
- **Summarized upward** — the level above must summarize, not merely list, the level below.

Minto's structure is the underlying skeleton of nearly every consulting deliverable and most
formal analytical reports.

### 1.3 SCQA — the opening frame

SCQA (Situation, Complication, Question, Answer), also from Minto, structures the opening
paragraph of a report [Source 4, Source 5]:

- **Situation** — the stable, accepted-by-reader context.
- **Complication** — what changed or what new fact disrupts the situation.
- **Question** — the question that disruption raises.
- **Answer** — the governing thought of the document (the BLUF sentence).

A worked SCQA opening:

> *Situation:* SMB renewal rates have been stable at 76-79% for eight consecutive quarters.
> *Complication:* In Q1 2026 renewal dropped to 67%, the largest single-quarter drop in three
> years. *Question:* What is driving the drop, and what should the company do about it?
> *Answer:* The drop is partner-channel cohorts retaining poorly; pause partner spend pending
> cohort recovery.

SCQA is what BLUF looks like when you need a paragraph rather than a sentence to set up the
governing thought.

---

## 2. Data storytelling (Cole Nussbaumer Knaflic)

The most-cited modern framework for analytical communication is Cole Nussbaumer Knaflic's
*Storytelling with Data* [Source 6, Source 7]. Her thesis: data alone does not persuade; a
deliberately constructed narrative around the data does. Six principles structure the framework.

### 2.1 Understand the context

Before opening any tool, the analyst must answer three questions [Source 6]:

- **Who** is the audience? (specific person or role, not "stakeholders")
- **What** do you want them to know or do? (one action verb)
- **How** will data support that message?

If the analyst cannot finish "I want [audience] to [action] because [data]" in one sentence, the
report is not yet ready to be written.

### 2.2 Choose the right visual

Match the chart to the question. Knaflic emphasizes a small vocabulary used precisely (table,
simple text, bar, line, scatter, slope graph) over an exotic chart that requires explanation.
For chart selection details, see `da-8-data-visualization`.

### 2.3 Eliminate clutter

Every element on a page must earn its place. Remove gridlines that don't aid comparison, borders,
default chart titles, redundant legends, and ornamental color. Apply the Gestalt principles
(proximity, similarity, enclosure, closure, continuity, connection) to organize what remains
[Source 6].

### 2.4 Focus attention with preattentive attributes

Use color, size, and position to direct the eye to the one thing that matters. Preattentive
attributes are processed in milliseconds before conscious attention — they tell the reader
*where to look first*. A finding chart should have one element in saturated color and everything
else in grey.

### 2.5 Think like a designer

Affordances (cues that suggest how to interact), accessibility (color-safe palettes, sufficient
contrast), aesthetics (deliberate alignment and whitespace), and acceptance (anticipating
audience pushback) all shape whether the finding lands [Source 6].

### 2.6 Tell a story — the narrative arc

Borrow the three-act structure from drama:

- **Setup** — the world before. Why this analysis exists, what was at stake.
- **Conflict** — the tension, the surprising finding, the gap between expectation and reality.
- **Resolution** — what the data shows, what to do about it, what happens next.

The narrative arc binds the BLUF answer to a *reason the audience cares*. Without conflict, a
report is a status update. With it, a report is an argument.

### 2.7 Action titles

A critical Knaflic technique: every chart and every slide gets a title that states the *finding*,
not the *topic* [Source 6, Source 8].

- Topic title (wrong): "Q1 Renewal Rate by Segment"
- Action title (right): "SMB renewal dropped 11 points in Q1; enterprise and mid-market held"

The action title is the BLUF of the chart. If the audience reads only the titles, they should
get the argument. Conventional reports — including most academic and engineering write-ups —
violate this by labeling figures with topic titles.

---

## 3. Audience adaptation: same finding, three ways

Every non-trivial analysis is consumed by at least three distinct audiences, and each demands a
different framing of the *same* underlying finding [Source 2, Source 9]. Failing to adapt is the
most common cause of "good analysis, no impact."

### 3.1 The executive audience

**Goal:** make a decision (often a resource-allocation decision) in 60-90 seconds.

**Constraints:** very short attention, high stakes, low tolerance for jargon or method detail,
needs confidence-level signal more than methodology.

**Format:** one-pager or single deck slide. BLUF in the first sentence. One headline metric. One
recommendation with one cost / one risk. Method buried or moved to appendix.

**Example framing:** "Pause partner-channel spend. Cost: $400K paused MDF. Risk: 90-day partner
revenue may dip ~6%. Confidence: high (n=2,400 customers, p<0.01)."

### 3.2 The technical audience

**Goal:** verify rigor and reproduce the work.

**Constraints:** tolerant of method detail, intolerant of unsupported claims, expects assumptions,
confounders, and limitations explicitly named.

**Format:** structured report with Methods, Results, Discussion, Limitations, Reproducibility
Appendix. Charts include error bars / confidence intervals. Data sources, query versions, and
seed values are documented.

**Example framing:** "Cohort retention curves by acquisition channel (Kaplan-Meier, log-rank
p=0.003). Partner-channel cohorts in the 2024Q3+ window show statistically lower 12-month
retention than direct cohorts. Confounders considered: ICP fit (controlled), price point
(controlled), competitive pressure (not controlled — see Limitations §5.3)."

### 3.3 The operational audience

**Goal:** know what to do differently on Monday morning.

**Constraints:** wants threshold rules and triggers, not statistical nuance; wants the report's
finding translated into a workflow change.

**Format:** runbook-style. "If X then do Y." Dashboards with red/yellow/green status. Clear
ownership.

**Example framing:** "New rule: customer success owns any partner-channel account onboarded
2024Q3 or later. Weekly health-check call required until 12-month renewal closes. Dashboard tab
'partner-retention' tracks all accounts in scope."

### 3.4 The "same finding, three audiences" template

When writing a real report, produce all three in parallel from one source of truth — the source
analysis. Variations of this in the wild include "the 1-3-25 rule" (1-page summary, 3-page
findings, 25-page detail) used in some consultancies, and tiered notebook reports with executive
abstract → narrative section → reproducibility appendix [Source 9].

---

## 4. The executive summary one-pager

The executive summary is the most-read, least-respected artifact in analytical communication.
Most are written last and treated as a perfunctory abstract; the experienced practice is to draft
it first as the *spine* of the analysis [Source 2, Source 10].

### 4.1 The one-pager structure

A working executive one-pager has six fixed slots, top to bottom:

| Slot | Purpose | Length |
|---|---|---|
| **1. Title** | Action title — the finding as a sentence. | 1 line |
| **2. BLUF** | Recommendation + headline number + confidence. | 2-3 lines |
| **3. Hero chart** | The single chart that earns the finding. | ~⅓ page |
| **4. So-what** | Why this matters to the reader's objectives. | 3-5 lines |
| **5. Cost / Risk / Alternatives** | What action costs, what could go wrong, what was rejected. | 3 bullets |
| **6. Next step + owner + date** | The single ask, named owner, target date. | 1 line |

Method, sample sizes, confidence intervals, and reproducibility live in a second page or
appendix link.

### 4.2 The "so what" test

For every sentence in an executive summary, ask "so what?" If you cannot answer in terms of the
reader's goals (revenue, risk, cost, customer experience, regulatory exposure), the sentence
belongs in the appendix [Source 2].

Common failures the "so what" test catches:

- Reporting *the metric moved* without reporting *what it costs the business*.
- Reporting *what was analyzed* instead of *what was found*.
- Reporting *what data was used* instead of *what decision is now possible*.

### 4.3 Headline metric framing

The headline number must satisfy four conditions:

1. **Direction** — the reader can tell at a glance whether this is good, bad, or neutral.
2. **Magnitude** — paired with a comparator (vs prior period, vs target, vs peer cohort).
3. **Materiality** — translated into a unit the reader cares about (dollars, customers, hours,
   risk score), not raw statistical output.
4. **Precision matched to certainty** — `~18%` reads as honest; `17.43%` reads as false precision
   (see anti-patterns §10.4).

---

## 5. Technical report structure

When the audience is technical (data science, engineering, audit, peer review), the report's job
shifts from persuasion to verifiability [Source 11, Source 12]. The classical structure derived
from scientific writing — IMRaD (Introduction, Methods, Results, Discussion) — adapts for
analytical work as the following sections.

### 5.1 Canonical sections

1. **Abstract / executive summary** — BLUF version of the finding (still required; technical
   readers also triage).
2. **Background / problem statement** — what question is being answered and why now.
3. **Data** — source, time range, schema, sample size, exclusion criteria, joins, freshness.
4. **Methods** — analysis steps, tool versions, parameter choices, validation strategy.
5. **Results** — findings with effect sizes, confidence intervals, and statistical tests where
   applicable.
6. **Discussion / interpretation** — what the results mean in context, alternate explanations
   considered and rejected.
7. **Limitations** — what the data cannot support, what was out of scope, known confounders.
8. **Recommendations** — the action implied by the results, with confidence level.
9. **Reproducibility appendix** — query text, notebook / script links, random seeds, package
   versions, data snapshot location.

### 5.2 The reproducibility appendix

A technical report without a reproducibility appendix is unverifiable. Minimum contents
[Source 13]:

- **Source query or pipeline** — full SQL / aggregation pipeline / dbt model name with git SHA.
- **Snapshot identifier** — point-in-time data reference (table snapshot, S3 prefix, dataset
  version) so the analysis can be re-run on the *same* data.
- **Tool versions** — Python, R, key packages, notebook environment lockfile.
- **Random seeds** — any RNG seed used for sampling, splitting, or model training.
- **Compute notes** — any non-determinism (parallelism, GPU, external API calls), and whether
  the result is invariant to it.

For Atlas / MongoDB analytical workloads, this typically means a saved aggregation pipeline plus
a snapshot timestamp and the cluster's deployment metadata.

### 5.3 Limitations sections done well

The credibility of a technical report is set by the Limitations section more than the Results
section. Patterns that signal honesty [Source 11]:

- Name the confounders that *were not* controlled for, not only those that were.
- State the population the finding generalizes to, and the population it does *not*.
- Distinguish *statistical* uncertainty (confidence intervals) from *epistemic* uncertainty
  (data quality, model assumptions, time-stability).
- Disclose the analyses run and *not reported* — pre-registered or otherwise — to head off
  garden-of-forking-paths suspicions.

Weak Limitations sections say "more data would help"; strong ones say "if X cohort had been
larger, we could test Y; until then, the finding holds only for Z."

---

## 6. Recommendation framing

A recommendation section is the riskiest part of the report — it converts evidence into opinion.
The discipline is to make the conversion *explicit* [Source 14, Source 15].

### 6.1 Evidence → opinion separation

For every recommendation, separate:

- **What the data shows** (fact, with effect size and uncertainty).
- **What we infer it means** (interpretation, dependent on assumptions).
- **What we recommend doing about it** (opinion, dependent on objectives and risk tolerance).

A reader who disagrees with the recommendation should still agree with the facts and (mostly)
the interpretation.

### 6.2 Confidence levels

Borrow the intelligence-community convention of explicit confidence language [Source 14]:

| Term | When to use |
|---|---|
| **High confidence** | Multiple independent lines of evidence, well-controlled design, large n, replicable. |
| **Moderate confidence** | Single line of evidence, some confounders unaddressed, or moderate n. |
| **Low confidence** | Suggestive only, small n, observational with confounders, exploratory analysis. |

State the confidence explicitly: "We recommend X with moderate confidence" beats "We recommend
X" with no confidence signal at all.

### 6.3 Alternatives considered

A recommendation gains credibility when the reader can see the alternatives that were rejected,
and *why* [Source 15]. A recommendation block has four parts:

1. **The recommendation** — one sentence.
2. **Why this option** — the one or two strongest arguments for it.
3. **Alternatives considered** — at minimum: status quo, the next-best option, and the most
   conservative option. State each, then state why each was rejected.
4. **Decision criteria** — the conditions under which a different option would become preferred
   (so the reader knows what evidence to watch for).

### 6.4 Action verbs and ownership

Recommendations without an owner are not recommendations. Every recommendation includes a verb,
a date, and a single named owner: "*Pause* partner-channel MDF spend *by 15 June*; *owner*: VP
Channel Sales."

---

## 7. Notebook-as-report: Jupyter, Quarto, Observable, R Markdown

Computational notebooks have become a primary medium for analytical reports because they bind
narrative, code, and output in one artifact [Source 12, Source 13, Source 16].

### 7.1 The four mainstream platforms

| Platform | Strengths | Best for |
|---|---|---|
| **Jupyter** | Ubiquitous in Python; widely supported in JupyterHub, Colab, VS Code. | Internal analyses; mixed audiences with technical reviewers in the loop. |
| **Quarto** | Multi-language (Python, R, Julia, OJS); polished HTML/PDF/Docx output; cross-references, citations, figure numbering; Pandoc-grade typesetting. | External-grade reports; academic-style write-ups; documents that need to look like a report, not a notebook. |
| **Observable** | Reactive JavaScript notebooks with live, interactive D3 charts. | Interactive web reports; explainers where reader interaction is the point. |
| **R Markdown** | Mature R ecosystem; LaTeX-quality output via Pandoc; large library of templates (papaja, bookdown). | Statistical reports in R; reproducible academic writing. |

Quarto is in practice the modern successor to R Markdown and a polished sibling to Jupyter — it
reads `.qmd` files directly and renders Jupyter notebooks too [Source 16].

### 7.2 Patterns that turn a notebook into a report

A notebook full of cells in execution order is *not* a report. Conversion requires:

1. **Hide code by default** — render output cells inline and collapse source. Quarto's
   `echo: false` and Jupyter's `nbconvert` `--no-input` are the standard knobs.
2. **Lift narrative to the top** — the first cell is the BLUF / executive summary, not the
   imports.
3. **Action-titled cell headings** — every section header is a finding, not a topic.
4. **Reproducibility block at the bottom** — kernel info, package versions, data snapshot
   reference, seed values (`session_info()` in R, `watermark` in Python).
5. **Parameterize for re-runs** — Papermill (Jupyter), Quarto parameters, or knit parameters
   (R Markdown) make the same notebook re-render against new data.
6. **Pin and lock** — commit the rendered HTML/PDF alongside source so reviewers see the
   version that was decided on.

### 7.3 What notebooks should *not* be used for

Notebooks fail when the audience is purely executive (use a deck or one-pager) or purely
operational (use a dashboard or runbook). They work best for the technical / mixed audience
where the reader may want to drill into a chart or a calculation.

---

## 8. Dashboards and live reports

Dashboards are a different artifact from reports: a report answers a specific question at a
point in time; a dashboard monitors a metric continuously [Source 17].

### 8.1 When a dashboard is the right answer

- The audience checks the same metric repeatedly over time.
- The decision is operational (act on threshold breaches), not strategic.
- A human owner watches the dashboard and is accountable for taking action on changes.

If those conditions don't hold, a dashboard is a graveyard — built once, looked at twice,
never decommissioned.

### 8.2 Dashboard structural rules

- **One question per page.** Multi-question dashboards force users to interpret which chart
  applies to their question.
- **Top-left is the headline.** Western readers scan top-left to bottom-right; the most important
  metric goes there with an action title.
- **Each chart has a target / threshold visible.** A line chart with no reference value cannot
  trigger action.
- **Red / yellow / green is reserved for status, not aesthetics.** Color should match the
  organization's escalation conventions.
- **Filters degrade trust.** Dashboards with many filters allow each reader to construct a
  different conclusion from the same data. Constrain filter sets deliberately.

### 8.3 Dashboards vs. exploratory tools

A dashboard is *not* a self-serve analytics surface. Confusing the two is one of the most common
analytical-org failures: stakeholders complain the dashboard "doesn't answer my question," and
analysts respond by adding filters until the dashboard becomes an exploratory tool nobody
trusts. Keep the two surfaces separate.

---

## 9. Honest framing of uncertainty

Honest framing is the difference between a report that survives scrutiny and one that
collapses on first challenge [Source 11, Source 14, Source 18].

### 9.1 What the data does and does not say

For every claim, the report should make explicit:

- The **population** the claim covers ("SMB customers acquired through partner channel in the
  last 18 months").
- The **time window** ("Q1 2026").
- The **conditions** under which the claim holds ("during the current product mix; pre new
  pricing launch").
- The **claim's logical type** — descriptive ("rates dropped"), associative ("partner cohorts
  correlate with lower retention"), or causal ("partner-channel acquisition causes lower
  retention"). Causal claims demand causal-grade evidence (randomization or strong quasi-experimental
  identification); observational data alone supports only descriptive and associative claims.

The cross-skill reference for the correlation-causation distinction is `da-1-6-1-correlation-vs-causation`.

### 9.2 Uncertainty taxonomy

| Type | What it means | How to disclose |
|---|---|---|
| **Sampling** | Random variation across hypothetical resamples. | Confidence intervals, p-values, standard errors. |
| **Measurement** | Underlying variables are imperfectly measured. | Note instrument / definition limitations; bound the effect. |
| **Model / specification** | The choice of model or features could change the result. | Run alternative specifications; report whether they qualitatively agree. |
| **Data quality** | Missingness, lag, freshness, ETL bugs. | Quantify missingness rates; date-stamp the snapshot. |
| **Generalization** | Will this hold next quarter, in a different region, post-policy change? | State explicitly which boundaries the finding *cannot* be extended past. |

A report that mentions only sampling uncertainty understates its actual uncertainty.

### 9.3 The "what we did not find" paragraph

A strong report includes a deliberate paragraph stating *what the data did not show* — null
results, hypotheses tested and rejected, segments where no effect was detected. This protects
against narrative-shaped selection bias (only reporting findings that fit the story) and signals
analytical rigor.

---

## 10. Anti-patterns (what wrecks analytical reports)

| # | Anti-pattern | What it looks like | Fix |
|---|---|---|---|
| 10.1 | **Burying the lede** | Three pages of context before the finding. | If the first sentence does not state the finding, rewrite (§1 BLUF / Minto / SCQA). |
| 10.2 | **Decoration over information** | 3D bars, gradients, drop shadows, multi-color palettes when one color would serve. | Apply Tufte's data-ink ratio [Source 19]; see `da-8-data-visualization`. |
| 10.3 | **Topic titles** | "Q4 Renewal Rates" instead of the finding. | Action title (§2.7): "Q4 renewals beat target by 4 points, driven by mid-market." |
| 10.4 | **False precision** | `17.43%` when variability is ±3 points; "$1,247,392" with 30% model uncertainty. | Match precision to the underlying estimate. Reserve significant figures for numbers that earn them. |
| 10.5 | **Single-source confirmation** | One model, one query, one definition, no robustness check. | Run one alternative specification and disclose whether the finding survives. |
| 10.6 | **Reverse-engineering the conclusion** | Picking the conclusion first, assembling evidence backward. | Pre-register the analysis; log analyses that *didn't* support the narrative; write the "what we did not find" paragraph (§9.3). |
| 10.7 | **AI-generated tells** | "In today's data-driven landscape…"; "leverage / unlock / harness / unleash"; "First… Furthermore… In conclusion…"; "truly transformative"; emoji bullets; "I hope this helps." | Use `kill-the-AI-ism` to diagnose; use `writing-expert` to enforce during drafting. |
| 10.8 | **Method-section autobiography** | Chronological narrative of what the analyst did. | Reorganize backward from the finding; move process detail to the appendix. |
| 10.9 | **Confidence-by-omission** | "We recommend X" with no confidence signal. | Always name confidence (§6.2). |
| 10.10 | **No-owner recommendations** | "The team should consider…" | Verb + date + single named owner (§6.4). |

---

## 11. Worked example: same analysis, three audiences

**Source analysis:** A churn-driver investigation finds that partner-channel customers in the
2024Q3-2025Q4 window churn at 13% over 12 months vs 5% for direct-acquired cohorts (log-rank
p=0.003, n=2,400 partner / n=4,100 direct).

A skeleton sketch of the three derived deliverables:

- **Executive one-pager** — action-titled finding; one hero chart; one ask, one cost, one risk;
  confidence level; named owner and date.
- **Technical write-up** — Kaplan-Meier curves with full statistics; confounders controlled and
  *not* controlled named explicitly; robustness checks on alternative specifications.
- **Operational runbook** — trigger condition; required action; dashboard reference; escalation
  rule; named owner.

The three deliverables share only the headline metric and the finding sentence; everything else
diverges in line with the audience adaptation rules in §3.

The three deliverable types — executive one-pager, technical write-up, and operational runbook
— follow the conventions in §4, §5, and §3.3 respectively.

---

## 12. Quick reference

### 12.1 Checklist before sending an analytical report

- [ ] Title states the finding, not the topic.
- [ ] First sentence is the BLUF (recommendation + headline number + confidence).
- [ ] One audience identified explicitly; the document is tuned for them.
- [ ] Every chart has an action title.
- [ ] Every recommendation has an owner, a date, and an explicit confidence level.
- [ ] Limitations section names confounders that were *not* controlled for.
- [ ] Reproducibility appendix links to query / notebook / snapshot.
- [ ] "What we did not find" paragraph present where applicable.
- [ ] No false precision (precision matches uncertainty).
- [ ] No AI-generated tells (run `kill-the-AI-ism` or the in-line checklist in §10.7).

### 12.2 Format-by-audience cheat sheet

| Audience | Artifact | Length | Method visibility |
|---|---|---|---|
| Executive | One-pager / 1 slide | 1 page | Hidden / linked appendix |
| Mixed (steering committee) | Deck or 3-pager | 5-10 slides / 3 pages | Summary in body, detail in appendix |
| Technical (peer review) | Structured report or notebook | 5-25 pages | Full Methods + Reproducibility appendix |
| Operational | Runbook / dashboard | < 1 page per rule | Hidden — translated into thresholds |

### 12.3 Five-minute report-quality test

If a colleague who is *not* the author reads the document for five minutes, they should be able
to state:

1. The finding (one sentence).
2. The recommended action (one sentence).
3. The confidence level (high / moderate / low).
4. One thing the finding does *not* say.

If they cannot, the report needs revision before sending.

---

## Sources

Foundational references for the inline `[Source N]` citations:

- **Knaflic, Cole Nussbaumer.** *Storytelling with Data.* Wiley, 2015 (Sources 1, 6, 7, 8).
- **Minto, Barbara.** *The Pyramid Principle.* Pearson, 3rd ed., 2008 (Sources 4, 5).
- **Few, Stephen.** *Now You See It* (2009); *Information Dashboard Design* (2013) (Sources 2, 17).
- **Tufte, Edward R.** *The Visual Display of Quantitative Information.* Graphics Press,
  2nd ed., 2001 (Source 19).
- **U.S. ODNI.** *Analytic Standards (ICD 203)* — confidence-level conventions (Source 14).
- **U.S. Army.** BLUF doctrine (Source 3).
- **Quarto** / **Project Jupyter** / **Wickham & Grolemund, R for Data Science** —
  notebook-as-report practice (Sources 11, 12, 13, 16).
- **Doumont, Jean-luc.** *Trees, Maps, and Theorems* (Source 9).
- **Harvard Business Review**, **Heath & Heath** *Decisive*, **Silver** *The Signal and the
  Noise* (Sources 10, 15, 18).

These sources are cited inline throughout the document as `[Source N]`.
