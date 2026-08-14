<!-- hub-reference-banner -->
> **Reference file — part of the `research-methodology` hub.** Formerly the standalone `systematic-review-research-integrity` skill.
> Sibling topics in this family are now reference files under the hubs (`research-methodology`, `deep-research`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: systematic-review-research-integrity
title: "Systematic Review & Research Integrity"
description: >-
  Formal evidence-synthesis methodology and open-science/research-integrity practice — the rigor layer for synthesizing a body of evidence and judging whether a study can be trusted. Covers the systematic-review pipeline (PRISMA 2020, PROSPERO, PICO, dual-reviewer screening), risk-of-bias appraisal (RoB 2, ROBINS-I), meta-analysis, GRADE certainty, scoping/living reviews, and research integrity (preregistration, registered reports, replication crisis, QRPs, reproducibility vs replicability, TOP/OSF). TRIGGER: planning or appraising a systematic review, meta-analysis, or scoping review; PRISMA/GRADE/RoB questions; judging whether evidence is trustworthy; preregistration, replication, QRPs. SKIP: ad-hoc web research/cited reports -> deep-research; research-thinking HOW -> deep-research-methods; AI research-tool landscape/provenance -> ai-assisted-literature-review; reproducibility as a foundational concept -> da-1-foundations-theory; stats/causal inference -> da-analytical-methods/da-12.
category: custom
version: 2
updated: "2026-06-15"
related_skills:
  - deep-research-methods
  - ai-assisted-literature-review
  - da-1-foundations-theory
  - da-analytical-methods
whenToUse:
  - "planning or appraising a systematic review, meta-analysis, or scoping review"
  - "PRISMA / PROSPERO / PICO / GRADE / risk-of-bias (RoB 2, ROBINS-I) questions"
  - "judging whether evidence or a single study is trustworthy"
  - "preregistration, registered reports, replication crisis, QRPs (p-hacking, HARKing)"
  - "reproducibility vs replicability, open data/code, TOP guidelines, OSF"
keywords:
  - systematic review
  - PRISMA
  - meta-analysis
  - GRADE
  - risk of bias
  - PROSPERO
  - preregistration
  - replication crisis
  - p-hacking
  - reproducibility
  - open science
tags:
  - research
  - evidence-synthesis
  - research-integrity
  - open-science
---

# Systematic Review & Research Integrity

This is the **rigor layer**. It answers two questions the hub's ad-hoc research skills do not:

1. **Evidence synthesis done properly** — when a question warrants a formal, protocol-driven
   review of an existing body of studies, and how to run that pipeline so the result is
   reproducible and defensible.
2. **Trustworthiness** — how to judge whether a study or a literature can be believed, and the
   open-science practices that make research credible in the first place.

> **Routing.** Use `deep-research` for ad-hoc, tool-driven web research that ends in a cited
> report. Use `deep-research-methods` for the general HOW of research thinking. Use the `da-*`
> skills for *doing* statistics, A/B tests, or causal inference. Come **here** when the
> deliverable is a formal evidence synthesis (systematic review / meta-analysis / scoping review)
> or when the question is about research *integrity* — preregistration, replication, QRPs,
> reproducibility, open data/code.
>
> Two close neighbours that share keywords but answer different questions: `ai-assisted-literature-review`
> owns the *AI-research-tool landscape and citation provenance* (which AI tool to use, hallucinated-citation
> and traceability checks) — this skill owns the *formal methodology* (PRISMA / meta-analysis / GRADE) those
> tools are meant to support. And `da-1-foundations-theory` owns reproducibility/replicability as a
> *foundational concept* in the epistemology of inference — this skill owns them as *integrity practice*
> (the replication crisis, preregistration, open data/code).

## Core concepts

- **Evidence synthesis** — methods for locating, appraising, and combining *all* studies that
  address a question, so the conclusion reflects the totality of evidence rather than a
  cherry-picked subset. A systematic review is the flagship; meta-analysis is the optional
  statistical combination step inside it.
- **Systematic review vs literature review** — a literature review is a narrative summary with no
  required method. A *systematic* review has a pre-specified protocol, an explicit reproducible
  search, transparent inclusion/exclusion, dual-reviewer screening, formal risk-of-bias appraisal,
  and (often) GRADE certainty. The method is the deliverable as much as the answer.
- **Reporting vs conduct** — PRISMA tells you how to *report* a review; the Cochrane Handbook
  tells you how to *conduct* one. They are complementary, not interchangeable.
- **Two halves of rigor** — (a) synthesizing existing evidence correctly, (b) producing/judging
  individual studies worth synthesizing. Garbage in, garbage out: a flawless review of p-hacked
  primary studies still misleads. That is why the integrity half lives in the same skill.

## The systematic-review pipeline (with PRISMA 2020)

A defensible review runs through these stages in order. Decisions are fixed in the protocol
*before* seeing results.

1. **Frame the question (PICO).** Population, Intervention, Comparator, Outcome. PICO drives
   eligibility criteria, search terms, and synthesis structure. For broad mapping questions use
   **PCC** (Population, Concept, Context) — that signals a scoping review.
2. **Register the protocol (PROSPERO).** Register *prospectively* — ideally before screening
   begins. A PROSPERO record captures key design/methods, the team, funder, and timelines, and is
   published immediately with a registration number. It is **not** peer-reviewed and **not** a full
   protocol. Registration deters outcome-driven changes to scope and lets others detect
   duplicate/abandoned reviews.
3. **Design the search.** Build sensitive searches across multiple databases (MEDLINE/PubMed,
   Embase, CENTRAL, etc.) plus registers and grey literature. Structure terms by concept:
   - **OR within a concept** → broad synonyms/MeSH → maximizes *sensitivity*.
   - **AND across concepts** → each concept must appear → adds *specificity*.
   - This is the **sensitivity vs precision** trade-off: reviews deliberately favor sensitivity and
     pay for it with more records to screen. Record the full strategy per database verbatim so it
     is reproducible.
4. **Screen in two stages, dual-reviewer.** Deduplicate, then (i) **title/abstract** screen, then
   (ii) **full-text** screen of survivors. Two independent reviewers at each stage, with a
   documented conflict-resolution rule. Record exclusion reasons at full-text stage — those numbers
   populate the flow diagram.
5. **Extract data.** Use a piloted, structured extraction form; extract in duplicate where
   feasible. Capture PICO details, study design, sample sizes, effect estimates and variances, and
   the items needed for risk-of-bias appraisal.
6. **Appraise risk of bias** (next section) — at the **result/outcome level**, not just study level.
7. **Synthesize** — narrative synthesis always; meta-analysis only when studies are similar enough.
   When studies are too heterogeneous to pool, structured narrative synthesis is the deliverable, not
   prose hand-waving: group studies (by population, intervention, or outcome), tabulate effects in a
   **harvest/summary-of-findings table**, and report direction and consistency of effects with the
   reasons they cannot be combined. Report it against **SWiM** (Synthesis Without Meta-analysis,
   Campbell et al., *BMJ* 2020) and follow the Popay et al. narrative-synthesis guidance; vote-counting
   by statistical significance is a known anti-pattern — count direction of effect, not p-value tallies.
8. **Rate certainty (GRADE)** and report against the PRISMA checklist.

### PRISMA 2020

PRISMA 2020 (Page et al., *BMJ* 2021;372:n71) is the current reporting standard and **replaces
PRISMA 2009**. It comprises:
- a **27-item checklist** organized into **7 sections** (Title; Abstract; Introduction; Methods;
  Results; Discussion; Other information);
- an **abstract checklist** (12 items);
- a **flow diagram** tracking records through **Identification → Screening → Included**, with counts
  of records identified, duplicates removed, records screened/excluded, reports sought/retrieved,
  reports assessed and excluded *with reasons*, and studies included. Templates are CC-BY.

Extensions: **PRISMA-ScR** (scoping), **PRISMA-LSR** (living), PRISMA-P (protocols), PRISMA-S
(search reporting), PRISMA-NMA (network meta-analysis).

## Risk of bias — RoB 2 and ROBINS-I

Modern appraisal is **domain-based and result-specific**: assess the risk of bias in a particular
*result* (an effect estimate for one outcome), via structured **signalling questions** per domain.

**RoB 2** — for **randomized trials** (Sterne et al., *BMJ* 2019). Five domains:
1. Bias arising from the **randomization process**
2. Bias due to **deviations from intended interventions**
3. Bias due to **missing outcome data**
4. Bias in **measurement of the outcome**
5. Bias in **selection of the reported result**

Judgements: **Low risk** / **Some concerns** / **High risk** (overall = High if any domain is High).

**ROBINS-I** — for **non-randomized studies of interventions** (Sterne et al., *BMJ* 2016; V2 2025
adds judgment algorithms). It frames each study as an attempt to emulate a target trial. Seven
domains: **confounding**; **selection** of participants; **classification of interventions**;
**deviations** from intended interventions; **missing data**; **measurement** of the outcome;
**selection of the reported result**. Judgements: **Low / Moderate / Serious / Critical**.
Confounding is the defining threat for non-randomized designs. Related tools: ROB-ME (non-reporting
bias), ROBIS (bias in a review), QUADAS-2 (diagnostic accuracy), AMSTAR 2 (quality of a review).

## Meta-analysis essentials

Meta-analysis statistically combines effect estimates across studies. Do it only when studies are
clinically and methodologically similar enough that a pooled number is meaningful.

- **Effect sizes.** Standardize before pooling: risk ratio (RR), odds ratio (OR), risk difference
  for binary; mean difference (MD) or standardized mean difference (SMD, Hedges' g) for continuous;
  hazard ratio (HR) for time-to-event. Each study contributes an estimate **and its variance**;
  weighting is by inverse variance.
- **Forest plot.** Each study is a block (area ∝ weight) at its point estimate with a 95% CI; the
  pooled estimate is the **diamond** at the bottom (width = pooled CI).
- **Heterogeneity.** Variation in true effects beyond chance:
  - **Q** (Cochran's chi-square) — tests *whether* heterogeneity exists; low power with few studies,
    so a non-significant Q does not prove homogeneity. **Don't pick the model from the Q-test.**
  - **τ²** — estimated between-study variance (the random-effects parameter).
  - **I²** — % of total variation due to heterogeneity rather than sampling error. Rough Cochrane
    bands: 0–40% might not be important; 30–60% moderate; 50–90% substantial; 75–100% considerable.
    Interpret alongside τ² and the size/direction of effects, not as a hard threshold.
- **Fixed-effect vs random-effects.**
  - **Fixed-effect** assumes *one* true effect; all variation is sampling error; weights favor large
    studies heavily.
  - **Random-effects** assumes true effects *vary* (drawn from a distribution); adds τ² to each
    weight, so small studies gain relative weight and the pooled CI is **wider** whenever I² > 0.
    Choose the model from the conceptual question (is one true effect plausible?), not a significance
    test — random-effects is usually the safer default when heterogeneity is expected.
- **Publication bias & funnel plots.** Studies with null results are less likely to be published,
  biasing the pooled estimate. A **funnel plot** charts effect (x) vs precision/SE (y); under no bias
  and low heterogeneity it is a symmetric inverted funnel. **Asymmetry** may indicate publication
  bias **or** small-study effects, true heterogeneity, or chance — suggestive, not diagnostic. Tests:
  Egger's regression, Begg's rank test, trim-and-fill (sensitivity). When asymmetry coexists with
  heterogeneity, random-effects pulls toward small (often biased) studies and is **not** automatically
  conservative (Sterne et al., *BMJ* 2011).

## GRADE — certainty of evidence

GRADE rates confidence in an estimate, **per outcome**, into four levels: **High / Moderate / Low /
Very low**. Certainty is separate from strength of recommendation.

- **Starting point.** RCT evidence **starts High**; non-randomized evidence **starts Low**.
- **Five domains that downgrade** (each −1, severe −2): **risk of bias**; **inconsistency**
  (unexplained heterogeneity); **indirectness**; **imprecision** (few events, wide CIs);
  **publication bias**.
- **Three domains that upgrade** (observational evidence, no downgrades): **large magnitude** (e.g.
  RR >2 or <0.5); **dose–response gradient**; **plausible residual confounding would shrink the
  observed effect**.

Certainty falls by at most three levels; you cannot rate below Very low. Present results + certainty
in a **Summary of Findings (SoF)** table (GRADEpro/GDT).

## Choosing the review type

| Type | When | Question | Reporting |
| --- | --- | --- | --- |
| **Systematic review (± meta-analysis)** | Defensible answer to a focused effectiveness/association question | PICO | PRISMA 2020 |
| **Scoping review** | Map breadth, clarify concepts, find gaps, decide if a SR is feasible | PCC (broad) | PRISMA-ScR |
| **Living systematic review** | Priority question, emerging evidence, low certainty — continual update | PICO | PRISMA-LSR |
| **Rapid review** | Decision deadline forces declared shortcuts | PICO | PRISMA + shortcuts noted |
| **Umbrella review** | Synthesize *existing* systematic reviews | — | PRIOR |

Heuristic (Munn et al. 2018): need to *answer* a meaningful question → systematic review. Need to
*map* what exists / whether a SR is warranted → scoping review (still needs an extensive documented
search, though it does not pool effects or usually appraise risk of bias).

## Open science & research integrity

The integrity half explains *why* primary studies may not be trustworthy and the practices that fix
it — the lens for appraising any single study you are about to cite or synthesize.

### The replication crisis

The **Open Science Collaboration (2015)** replicated 100 psychology studies: **97% of originals were
significant, but only 36% of replications were**, and replication effects were about **half** the
original magnitude. The takeaway is not "science is broken" (NASEM 2019 cautions against the crisis
framing) but that **a single significant published result is weak evidence**, and synthesis must
weight that.

### Questionable research practices (QRPs)

QRPs exploit **researcher degrees of freedom** (Simmons, Nelson & Simonsohn 2011): defensible-looking
analytic choices that, used opportunistically, make false positives "vastly more likely":
- **p-hacking** — trying analyses/exclusions/covariates/subgroups until p < .05, reporting only what
  "worked." *Fix:* pre-specify the analysis; report all variables, conditions, and exclusions.
- **HARKing** (Kerr 1998) — *Hypothesizing After Results are Known*; presenting a post-hoc finding as
  a priori. *Fix:* preregister hypotheses; label exploratory results as exploratory.
- **Optional stopping** — peeking and stopping collection once significant, inflating the false-positive
  rate. *Fix:* fix the stopping rule / sample size in advance (or use sequential designs with
  alpha-spending).
- **The garden of forking paths** (Gelman & Loken 2014) — even *one* analysis per dataset inflates
  error because you *would have* analyzed differently given different data. *Fix:* preregistration;
  multiverse/specification-curve analysis.
- **Selective/outcome reporting** — publishing only significant outcomes/studies (feeds publication
  bias). *Fix:* register outcomes; report all; share data.

### Preregistration & registered reports

- **Preregistration** time-stamps the hypotheses, design, and analysis plan *before* data are seen
  (OSF, AsPredicted). Its purpose is to **separate confirmatory from exploratory** work. It doesn't
  forbid exploration; it labels it. Deviations are allowed but must be disclosed.
- **Registered Reports** go further: the **introduction, methods, and analysis plan are peer-reviewed
  *before* data collection**; passing designs earn **in-principle acceptance** — publication is
  guaranteed regardless of whether results are positive or null. This structurally removes publication
  bias, p-hacking incentives, and HARKing.

### Reproducibility vs replicability (NASEM 2019)

- **Reproducibility** = **computational reproducibility**: same results from the *same data and code*.
  With transparent reporting and shared artifacts this *should* always hold; a failure is a defect.
- **Replicability** = consistent results from *new data* on the same question. Even a perfect study may
  fail to replicate; a single failed replication does not refute the original.
- **Generalizability** = whether results extend to other populations/settings.

### Open data, code, materials & the TOP guidelines

The **TOP Guidelines** (Center for Open Science; **TOP 2025**) are a policy framework with **seven
Research Practices** — Study Registration, Protocol Transparency, Analysis-Plan Transparency,
Materials Transparency, Analytic-Code Transparency, Data Transparency, Reporting Transparency — each at
**three escalating levels**: **Disclose** → **Share & Cite** → **Certify**. TOP 2025 adds
**Verification Practices** (computational reproducibility; comprehensive reporting). **Open Science
Badges** (open data / open materials / preregistration) reward and signal these practices and measurably
raise data-sharing rates. The **OSF** (osf.io) is the common infrastructure for preregistrations,
materials, data, code, and preprints.

## Pitfalls

- **Calling a narrative review "systematic."** No protocol, no reproducible search, no dual screening →
  it is a literature review. Don't claim PRISMA compliance you didn't earn.
- **PRISMA as conduct.** PRISMA is *reporting*; following the checklist doesn't make weak methods
  rigorous. Conduct lives in the Cochrane/JBI handbooks.
- **Skipping prospective registration.** Registering *after* screening defeats the purpose — scope can
  be retrofitted to findings.
- **Meta-analyzing apples and oranges.** A pooled number from clinically heterogeneous studies is
  precise nonsense. High I² is a stop-and-think signal.
- **Picking fixed vs random by the Q-test.** Choose from whether one true effect is plausible, not a
  low-powered significance test.
- **Treating funnel-plot symmetry as proof of no bias.** Asymmetry has many causes; symmetry with few
  studies proves little.
- **Study-level instead of result-level RoB.** A study can be low-risk for one outcome and high-risk for
  another.
- **Over-trusting a single significant study.** Given the replication record, weight one p < .05 finding
  lightly; prefer preregistered, replicated, well-powered evidence.
- **Conflating reproducibility and replicability.** Same data vs new data — different meanings; loose use
  muddies what a "failure" implies.
- **Preregistration theater.** Registering then silently deviating, or preregistering analyses so vague
  they constrain nothing. Specificity and disclosed deviations give it force.

## References

- PRISMA 2020 statement — Page MJ et al., *BMJ* 2021;372:n71 — https://www.prisma-statement.org/prisma-2020-statement · https://www.bmj.com/content/372/bmj.n71
- PRISMA 2020 flow diagram & checklist (CC-BY) — https://www.prisma-statement.org/prisma-2020-flow-diagram · https://www.prisma-statement.org/prisma-2020-checklist
- PRISMA-ScR — Tricco AC et al., *Ann Intern Med* 2018 — http://www.equator-network.org/reporting-guidelines/prisma-scr/
- PRISMA-LSR — Akl EA et al., *BMJ* 2024;387:e079183
- Cochrane Handbook for Systematic Reviews of Interventions v6.5 (2024) — https://www.cochrane.org/authors/handbooks-and-manuals/handbook/current
- RoB 2 & ROBINS-I — Sterne JAC et al., *BMJ* 2019 (RoB 2) / 2016;355:i4919 (ROBINS-I) — https://www.riskofbias.info/welcome/robins-i-v2
- GRADE Handbook — https://gradepro.org/handbook/ · Cochrane Ch.14 — https://www.cochrane.org/authors/handbooks-and-manuals/handbook/current/chapter-14
- Funnel-plot asymmetry — Sterne JAC et al., *BMJ* 2011;343:d4002 — https://www.bmj.com/content/343/bmj.d4002
- Scoping vs systematic review — Munn Z et al., *BMC Med Res Methodol* 2018;18:143 — https://link.springer.com/article/10.1186/s12874-018-0611-x
- PROSPERO — CRD, University of York — https://www.crd.york.ac.uk/PROSPERO/help/register
- Preregistration & Registered Reports — Center for Open Science — https://www.cos.io/initiatives/prereg · https://www.cos.io/initiatives/registered-reports
- TOP Guidelines 2025 & Open Science Badges — https://www.cos.io/initiatives/top-guidelines · https://www.cos.io/initiatives/badges
- Open Science Collaboration, "Estimating the reproducibility of psychological science" — *Science* 2015;349:aac4716 — https://www.science.org/doi/10.1126/science.aac4716
- Simmons, Nelson & Simonsohn, "False-Positive Psychology" — *Psychological Science* 2011;22:1359 — https://journals.sagepub.com/doi/full/10.1177/0956797611417632
- NASEM, "Reproducibility and Replicability in Science" (2019) — https://www.nationalacademies.org/read/25303/chapter/3
