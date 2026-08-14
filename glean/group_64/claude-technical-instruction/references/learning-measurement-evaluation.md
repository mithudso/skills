<!-- hub-reference-banner -->
> **Reference file — part of the `technical-instruction` hub.** Formerly the standalone `learning-measurement-evaluation` skill.
> Sibling topics in this family are now reference files under the hubs (`technical-instruction`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: learning-measurement-evaluation
description: >-
  Measure training and enablement effectiveness to prove business value. Covers Kirkpatrick Four
  Levels, New World Kirkpatrick Model, Phillips ROI Level 5 (isolate effect, monetize, ROI%/BCR),
  xAPI/SCORM/cmi5, LRS architecture, training transfer (Baldwin-Ford, LTSI), leading/lagging KPIs,
  L&D dashboards. TRIGGER: Kirkpatrick model / Level 3 behavior / Level 4 results / Required
  Drivers; Phillips ROI Level 5 / isolate program effect; xAPI vs SCORM / cmi5 / learning record
  store; LTSI / transfer climate; L&D dashboard / completion rate vanity. SKIP: course design →
  instructional-design-course-architecture; psychometrics/item writing →
  assessment-certification-design; training delivery → technical-training-delivery;
  product adoption/cohort → da-applied-and-communication; causal inference →
  da-analytical-methods; cognitive load/spacing → applied-psychology; adoption behavior →
  applied-psychology; GenAI tutors → genai-education-instructional-design. Family hub →
  technical-instruction.
version: "1.0.2"
updated: "2026-06-16"
category: education
whenToUse:
  - "Calculate ROI for a training program using the Phillips methodology"
  - "Isolate a training program's effect on a business KPI"
  - "Choose between SCORM, xAPI, and cmi5 for a new eLearning project"
  - "Diagnose why training didn't transfer to on-the-job behavior"
  - "Design a Kirkpatrick Level 3 behavior evaluation"
  - "Build an L&D dashboard that executives will trust"
  - "Use the LTSI to identify what's blocking training transfer"
  - "Prove training business value to finance or a C-suite sponsor"
  - "Design post-training surveys, observation checklists, or action plans"
  - "Report training results using the three-tier stakeholder structure"
  - "Understand Required Drivers and the New World Kirkpatrick Model"
related_skills:
  - instructional-design-course-architecture
  - assessment-certification-design
  - technical-training-delivery
  - applied-psychology
  - da-analytical-methods
keywords:
  - Kirkpatrick four levels
  - New World Kirkpatrick Model
  - Required Drivers
  - leading indicators training
  - Phillips ROI methodology
  - Level 5 ROI training
  - benefit-cost ratio
  - isolate program effect
  - fully loaded program cost
  - xAPI Experience API Tin Can
  - SCORM elearning standard
  - cmi5 LMS launch standard
  - learning record store LRS
  - learning analytics
  - training transfer measurement
  - Baldwin Ford transfer model
  - Learning Transfer System Inventory LTSI
  - transfer climate
  - enablement KPIs
  - L&D dashboard
  - completion rate vanity metric
  - time to proficiency
  - on-the-job behavior change
  - evaluation instrument design
  - ROI percent training
  - backward design evaluation
tags:
  - education
  - training
  - evaluation
  - l-and-d
  - enablement
  - measurement
  - kirkpatrick
  - phillips-roi
  - xapi
  - learning-analytics
---

# Learning Measurement & Training Evaluation

Reference for measuring training and enablement program effectiveness across the full evaluation stack — from post-workshop smile sheets to executive ROI reports and cross-system xAPI analytics.

<!-- Contents
1. Evaluation Framework Overview
2. Kirkpatrick's Four Levels
3. The New World Kirkpatrick Model
4. Instrument Design by Level
5. Reporting Kirkpatrick Data to Stakeholders
6. Phillips ROI Methodology (Level 5)
7. xAPI, SCORM, and cmi5 Standards
8. Learning Record Store (LRS)
9. Training Transfer Measurement
10. Leading vs Lagging Enablement Indicators
11. L&D Dashboard Design
12. Evaluation Edge Cases
13. Quick Reference: Evaluation Planning Decision Guide
14. References
-->

See `references/learning-measurement-context.md` for the full deep-reference with worked examples, decision tables, and vendor comparisons.

## Evaluation Framework Overview

Three frameworks dominate training program evaluation. They are complementary, not competing:

| Framework | Primary Question | Output |
|---|---|---|
| Kirkpatrick Four Levels | Did participants react, learn, change behavior, get results? | Evidence at each level |
| New World Kirkpatrick | Were the right conditions built in from the start? | Required Drivers + leading indicators |
| Phillips ROI (Level 5) | Did the financial benefits exceed the fully-loaded cost? | ROI % and BCR |

**When to use each:** Kirkpatrick levels 1–4 apply to any formal training; New World Kirkpatrick adds the pre-design planning discipline; Phillips Level 5 is appropriate for 5–10% of programs — those that are high-cost, strategically critical, and directly tied to measurable business KPIs.

## Kirkpatrick's Four Levels

First published by Donald Kirkpatrick (1959) and updated as the New World Kirkpatrick Model (2016). Each level is harder and more valuable to measure than the previous.

### Level 1: Reaction
Measures whether participants found training favorable, engaging, and relevant to their jobs. Relevance is the most predictive sub-component — if training is not perceived as relevant, application almost never occurs.

**Instrument:** Post-training survey delivered within 24 hours. Include forward-looking items: "What will you do differently?" and "What might prevent you from applying this?" Avoid satisfaction-only smile sheets — they produce data with no correlation to behavior change.

**Limitation:** A 2009 study (n=335) found no statistically significant correlation between Level 1 scores and Level 3 behavior change. Alliger & Janak (1989) found the Level 1→Level 2 (reaction→learning) causal link produced only r=.23 — not the strong predictive relationship the model implies. A high satisfaction score is evidence of a pleasant training event, not of future performance.

### Level 2: Learning
Measures acquisition of the intended knowledge, skills, attitude, confidence, and commitment (5 components per the New World Model). Confidence and commitment are early predictors of Level 3 transfer.

**Instruments:** Pre-test/post-test; role plays; scenario-based cases; confidence items on a 5-point scale.

### Level 3: Behavior
Measures whether participants apply what they learned when back on the job, sustained over time. This is where most organizational value is created — and where most measurement stops.

**Structural barriers:** Most organizations skip formal Level 3 evaluation because it requires sustained effort beyond the training event — Kennedy et al. (2013) found 40%+ report L3 evaluation is not required by management, and broader surveys suggest the majority never measure on-the-job application.

**Instruments:** Structured observation checklists; 30/60/90-day surveys; 360-degree feedback; work product analysis; action plans reviewed in manager 1:1s.

**Timing:** Begin within 2–4 weeks of training, not after 90 days.

### Level 4: Results
Measures whether targeted organizational outcomes occurred. Examples: revenue, error rates, customer CSAT, safety incidents, retention, time-to-proficiency.

**The attribution challenge:** Business results are influenced by dozens of factors. "Training contributed to the result, particularly where reinforcement was consistent" is defensible; "training caused this" requires control group evidence.

**Instruments:** Pre-established KPI baselines; trained vs. untrained comparison groups; Return on Expectations (ROE) check-ins; operational metrics from existing systems.

## The New World Kirkpatrick Model

Published in 2016 by Jim and Wendy Kirkpatrick. Three additions to the original model:

### Shared Accountability Model
The New World model reframes who owns each level: L1 and L2 are the L&D team's accountability (did we design and deliver well?); L3 and L4 are shared accountability between L&D, managers, and the business (did we build the conditions for application and results?). This framing is practical for securing manager buy-in — they own the on-the-job reinforcement half of the equation.

### Required Drivers (Level 3 Addition)
Post-training reinforcement systems that management must provide: coaching, job aids, work review, recognition, and pay-for-performance. Organizations with well-implemented Required Drivers can expect ~85% on-the-job application versus ~15% for training-event-only approaches (citing Brinkerhoff, 2006 — treat as directional, not precise). [TENTATIVE: this figure appears in Kirkpatrick-affiliated sources; independent RCT verification is limited — treat as a directional heuristic.]

### Leading Indicators (Level 4 Bridge)
Short-term observations that signal whether critical behaviors are on track to produce desired results. They allow course correction before the lagging outcome data arrives.

### Backward Design Mandate
Plan the evaluation in reverse: start at Level 4 (define the business result with the sponsor), then Level 3 (Required Drivers and critical behaviors), then Level 2 (learning objectives), then Level 1 (learning experience). This prevents the common trap of starting at Level 1 and hoping Level 4 will follow.

## Instrument Design by Level

See `references/learning-measurement-context.md` §Instrument Design for full checklists.

**Summary quick reference:**

| Level | Primary Instrument | Timing | Key Design Rule |
|---|---|---|---|
| L1 Reaction | Post-training survey, 8–10 items | Within 24 hours | Include relevance + forward-looking + confidence items |
| L2 Learning | Pre/post-test + confidence rating | Before and immediately after | Match format to Bloom's level — see `assessment-certification-design` for item writing |
| L3 Behavior | Structured observation + surveys | 2–4 wks, then 30/60/90 days | Brief managers before training; establish inter-rater reliability — see `assessment-certification-design` |
| L4 Results | KPI baselines + comparison group | 3–12 months post-training | Agree on KPI with sponsor before design; collect baseline before training |

## Reporting Kirkpatrick Data to Stakeholders

**Tiered reporting structure:**

- **Tier 1 — Executive (1 page):** Bottom-line outcome first: what business metric moved, by how much, what it is worth. ROI or ROE %, total investment, one human example. Never bury the headline.
- **Tier 2 — L&D buyers/managers:** Behavior change evidence, completion, assessment scores, quotes. Explicitly link Level 2 gains to Level 3 observed behaviors.
- **Tier 3 — Finance:** Calculation assumptions, data sources, isolation methodology, and sensitivity scenarios.

**Communication rules:** Lead with business impact, not training activity. Speak the KPI language of the business sponsor. Explicitly state attribution confidence: "Training contributed; reinforcement consistency was the mediating factor." Show leading indicators alongside lagging outcomes.

## Phillips ROI Methodology (Level 5)

Developed by Jack Phillips (1973, published 1983). Adds financial return calculation on top of Kirkpatrick's four levels.

### The 5-Level Framework

| Level | Focus |
|---|---|
| 1 | Reaction & Planned Action |
| 2 | Learning |
| 3 | Application & Implementation |
| 4 | Business Impact |
| 5 | Return on Investment (ROI %) |

Level 5 adds: monetize the Level 4 impact, subtract fully-loaded program costs, compute ROI % and BCR.

### ROI Formulas

```
BCR  = Program Benefits / Program Costs
ROI% = [(Benefits − Costs) / Costs] × 100
```

BCR = 1.96 means $1.96 in gross benefits per $1 invested. ROI% = 96% means after recovering the investment, an additional $0.96 per dollar is returned.

### Isolation Methods (in descending reliability)

| Method | Reliability | When to Use |
|---|---|---|
| Control group (trained vs. untrained) | Highest | Gold standard; use when feasible |
| Trend-line analysis | High | Clean historical data series; no control group possible |
| Mathematical/forecasting model | High if data-rich | Complex multi-variable environments |
| Manager estimates (confidence-adjusted) | Medium | Management/leadership programs |
| Participant estimates (confidence-adjusted) | Low–Medium | When no objective data exists |

**Key principle:** Without isolation, Level 5 produces correlation, not causation.

### Fully-Loaded Cost Categories

Include: needs assessment, design/development (prorated), materials, facilities, facilitator/coach salaries + prep time, participant salaries + time in training, travel, evaluation costs, and overhead allocation. Leave no cost hidden — conservative full-costing builds credibility.

### The 12 Guiding Principles

Abbreviated; full text in the references file:
1–4: Collect lower-level data, use credible sources, choose conservative alternatives, collect data at each level when evaluating a higher level.
5: Use at least one isolation method — non-negotiable.
6: Assume no improvement when participant data is missing.
7: Adjust estimates downward for potential error.
8–9: Exclude extreme data items; use first-year benefits only for short programs.
10: Fully load all costs.
11: Intangibles are measures purposely NOT converted to money — legitimate, not a failure.
12: Communicate results to all key stakeholders.

### When to Use Level 5

Apply Level 5 to ~5–10% of programs. Use it when: program is high-cost or high-audience; directly tied to a measurable business KPI; leadership is questioning the budget; and you have a viable isolation method. Skip it for mandated compliance training, low-cost programs, and culture/awareness initiatives.

### Critiques

- Isolation via estimates is subject to upward bias (~2× true effect per Sauermann & Stenberg, 2020 IZA field experiment); conservative-adjustment principle partially corrects but does not eliminate this.
- Methodology developed by a company that sells it — limited independent peer review.
- Monetary conversion of soft data is not GAAP-compliant; CFOs may not accept it as auditable.
- Too late for formative improvement (calculate after program ends).

## xAPI, SCORM, and cmi5 Standards

### Standards at a Glance

| Standard | Origin | Communication | Requires LMS | Mobile | Offline | Customizable Data |
|---|---|---|---|---|---|---|
| SCORM 1.2 | ADL, 2001 | JavaScript API (browser) | Yes | No | No | No (fixed CMI model) |
| SCORM 2004 | ADL, 2003 | JavaScript API (browser) | Yes | No | No | No |
| xAPI | ADL, 2013; IEEE, 2023; ISO, 2025 | REST/JSON over HTTP | No | Yes | Yes (queue+sync) | Yes (extensions) |
| cmi5 | AICC/ADL, 2016 | REST/JSON (xAPI profile) | Yes | Yes | Yes | Yes (within cmi5 rules) |

### SCORM
Still the dominant standard: 81.7% adoption in 2024 industry surveys. Tracks completion, score, session time, bookmark, and interactions. Hard ceiling: LMS-only, browser-only, fixed data model, siloed in the LMS database.

### xAPI (Experience API)
Actor–Verb–Object statement model sent via REST to an external Learning Record Store (LRS). Tracks anything, anywhere: mobile apps, offline activity, simulations, VR, social learning, on-the-job behavior. Ratified as IEEE 9274.1.1-2023 and ISO/IEC/IEEE 39274-1-1:2025 — significant for government/regulated-industry procurement.

**Honest adoption picture (2024–2026, verified-as-of: 2026-06-16):** xAPI adoption remains ~17% despite 10+ years of promotion. The primary barriers: authoring tools produce "SCORM in xAPI format" (same completion/score data, different transport); LRS costs are significant (enterprise tier starts ~$4,000/month, verified-as-of: 2026-06-16); data governance (verb vocabularies, identity schemes) requires developer expertise most L&D teams lack.

### cmi5
An xAPI profile that adds defined launch mechanics, authentication, session management, and completion/mastery criteria — the interoperability SCORM always had, now running on xAPI. "Plug and play" for LMS environments, xAPI's richness underneath. Still a minority of LMS platforms support it.

## Learning Record Store (LRS)

An LRS is a REST API server that accepts, stores, and returns xAPI statements. Analytics live in the BI tool or purpose-built analytics platform on top of the LRS.

**Key vendors (2025–2026):**

| Vendor | Tier | Notes |
|---|---|---|
| Watershed LRS | Enterprise | ~$4,083+/mo; relational DB for direct BI access; ADL conformant |
| Learning Locker / Learning Pool | Mid-market | Open-source core; multi-tenant; BI connectors |
| SCORM Cloud (Rustici) | Developer/Mid | Free tier; supports SCORM+xAPI+cmi5; popular for POC and OEM |
| Yet Analytics | Enterprise/Gov | Strong semantic governance; ADL/DoD contexts |
| GrassBlade | SMB/WordPress | Low cost; popular for WordPress LMS environments |

**Pipeline pattern:** LRS → data warehouse → BI tool (Tableau, Power BI, Looker) → join with HRIS, CRM, ERP.

**What xAPI enables beyond SCORM:** Behavioral skill transfer tracking (not just completion); blended learning attribution across multiple platforms; simulation/VR performance correlation with field outcomes; predictive indicators (quiz hesitation time → safety incident risk).

## Training Transfer Measurement

### The Transfer Problem
Most training investment fails to change on-the-job behavior. Multiple studies converge on: employees apply ~15% of training without reinforcement systems (Brinkerhoff, cited by Kirkpatrick Partners) [TENTATIVE — directional estimate from Kirkpatrick-affiliated sources, not RCT-derived]. ATD (2023): 87% of organizations find it difficult to isolate training's impact; 75% of TD professionals lack access to the data needed for high-level evaluation.

### The Baldwin & Ford (1988) Model
Three inputs determine transfer: **trainee characteristics** (ability, motivation), **training design**, and **work environment** (supervisor support, peer support, opportunity to use). Work environment variables consistently show stronger transfer relationships than training design variables.

### Learning Transfer System Inventory (LTSI)
Developed by Holton, Bates & Ruona (2000). The empirically validated instrument for diagnosing which system factors are blocking transfer. 16 factors across two domains:

**Training-Specific (11 factors):** Learner Readiness, Motivation to Transfer, Positive/Negative Personal Outcomes, Personal Capacity for Transfer, Peer Support, Supervisor Support, Supervisor Sanctions, Perceived Content Validity, Transfer Design, Opportunity to Use.

**Training-General (5 factors):** Transfer Effort-Performance Expectations, Performance Outcomes Expectations, Openness to Change, Performance Self-Efficacy, Performance Coaching.

**Use:** Survey trainees before and after a program to identify which factors to intervene on.

### Supervisor Support: The Single Strongest Driver
Meta-analysis (*Human Factors*, 2019): Supervisor support correlated with sustained transfer at ρ = .51. Action-related support (discussing concrete application opportunities) predicts transfer; attitude-related support (expressing interest) does not. Only involvement and accountability (not coaching or openness) predicted transfer in a 3-wave longitudinal study (KU Leuven, 2018).

## Leading vs Lagging Enablement Indicators

| Type | When Available | Examples |
|---|---|---|
| Leading (process/activity) | During/immediately after learning | Completion rate, pass rate, learner satisfaction, manager reinforcement score, content drop-off rate |
| Lagging (outcome) | After learning, often weeks/months | Behavior change rate, error rate, time-to-proficiency, sales conversion, retention |

**Completion rate is a compliance metric, not a performance metric.** It measures that training was attempted; it says nothing about behavior change. The Wells Fargo ethics training case is the canonical example: 100% completion, millions of fraudulent accounts.

**Manager reinforcement score** is the highest-leverage leading indicator for technical enablement: it captures pre-transfer environment quality before outcomes materialize, and is directly actionable if low.

## L&D Dashboard Design

Three-audience design (see `references/learning-measurement-context.md` §L&D Dashboard for source detail):
1. **Executive view:** Capability velocity, time-to-proficiency, ROI by program category. One page, dollar-denominated.
2. **Manager view:** Team completion, behavioral change indicators, skill progression.
3. **L&D ops view:** Engagement drop-offs, quiz item analysis, content effectiveness.

**Top 10 metrics for a lean technical-enablement team:**

| Priority | Metric | Type | Data Source |
|---|---|---|---|
| 1 | Time-to-productivity (new hires / role changes) | Lagging | Manager + HRIS |
| 2 | 30/60/90-day skill retention (retrieval quiz decay) | Bridge | LMS + custom assessments |
| 3 | Manager reinforcement score (LTSI Supervisor Support) | Leading | Survey at 2–4 weeks post-training |
| 4 | Application rate (% applying skill, manager-observed at 60 days) | Lagging | Manager survey |
| 5 | Completion rate by cohort | Leading/process | LMS |
| 6 | Assessment pass rate vs. re-take rate | Leading | LMS |
| 7 | Learner confidence-to-apply score | Leading | Post-training survey |
| 8 | Support ticket / escalation rate for trained topics | Lagging | CRM/support system |
| 9 | Content drop-off rate (where do learners disengage?) | Diagnostic | LMS analytics |
| 10 | Cost per learner vs. performance lift (by program) | ROI | Finance + LMS |

**Structural access problem:** The most useful metrics (error rates, time-to-proficiency, sales win rate) live in CRMs, ERPs, and HRIS, not the LMS. Solution: define data-access agreements, baselines, and measurement windows at program brief stage, not after.

## Evaluation Edge Cases

### When to Skip Evaluation Entirely

For programs under ~20 participants or a one-time delivery, even Level 1 surveys may not justify the administrative overhead. Document the rationale for skipping: cohort size, delivery frequency, and availability of proxy data (e.g., existing performance metrics that would surface any impact). Mandated compliance training warrants completion tracking only unless a regulator requires evidence of comprehension.

### Retrospective Evaluation (No Pre-Training Baseline)

When evaluation is commissioned after training has already run — a common scenario — the absence of a pre-training baseline does not eliminate all options:
1. Use historical trend-line data from CRM, HRIS, or support systems for the pre-training period.
2. Identify a comparison cohort not yet trained (staggered rollout is useful here).
3. Use subject-matter expert estimates of the pre-training performance level, adjusted conservatively.
4. Document the limitation explicitly and reduce the confidence rating of any Level 4 or Level 5 claim accordingly.

### ROI Report Structure (for Finance and Executives)

A Level 5 ROI deliverable should follow this structure to be credible with CFOs:
1. **Executive summary** — ROI %, BCR, payback period, one-sentence business impact statement.
2. **Methodology** — isolation method used, data sources, confidence level, conservative adjustments applied.
3. **Financial detail** — fully-loaded cost breakdown, benefit conversion method and value per unit, gross benefit calculation.
4. **Intangibles appendix** — measures intentionally not monetized, with qualitative evidence.
5. **Sensitivity analysis** — how ROI % changes if key assumptions shift ±20%.

### Data Privacy in xAPI/LRS Deployments

xAPI statements include actor identity (email in `actor.mbox` by default). Before deploying an LRS in production:
- **PII handling**: use anonymized agent IDs or team-level aggregation for reporting where individual tracking is not required.
- **EU/UK**: xAPI learner data is personal data under GDPR; ensure data processing agreements with LRS vendors and confirm data residency.
- **FERPA**: applies to educational institution learners; workforce L&D programs are typically exempt unless the employer is also an educational institution.
- **Data governance**: define verb vocabulary and actor identity scheme before ingesting data — retrofitting governance on 6 months of statements is operationally expensive.

## Quick Reference: Evaluation Planning Decision Guide

```
Is the program high-cost or strategically critical?
├── Yes → Plan Levels 1–4; consider Level 5 if isolation is feasible
│         Plan backward: define L4 KPI → L3 behaviors → L2 objectives → L1 design
└── No  → Plan Levels 1–2; add L3 check if behavior change is the goal

Can you isolate the program's effect?
├── Yes (control group) → Level 5 ROI is defensible
├── Yes (trend-line) → Level 5 is plausible with honest caveats
├── Estimates only → Use Level 4 + conservative estimate; state confidence level
└── No → Report Level 4 direction + intangibles only; no ROI claim

Do you have xAPI ambitions?
├── Authoring tool exports xAPI natively → cmi5 if LMS supports it, else SCORM
├── Need mobile/offline/cross-platform → xAPI + LRS required; plan data governance
├── Simple completion + score → SCORM 1.2 is pragmatic
└── Government/regulated procurement → xAPI (IEEE/ISO status matters for RFPs)
```

## References

See `references/learning-measurement-context.md` for the complete annotated reference list (30+ sources). Key primary sources:

1. Kirkpatrick, D. (1959–1960). Four articles in *Journal of the American Society of Training Directors*. The original four-level model.
2. Kirkpatrick, J. & Kirkpatrick, W. (2016). *Kirkpatrick's Four Levels of Training Evaluation*. ATD Press.
3. Kirkpatrick Partners. *Introduction to The New World Kirkpatrick Model* (2022 PDF). https://www.kirkpatrickpartners.com/
4. Phillips, J. (1983). *Handbook of Training Evaluation and Measurement Methods*. Gulf Publishing.
5. Phillips, J. *Introduction to the ROI Methodology* (ROI Institute PDF). https://roiinstitute.net/
6. Baldwin, T.T. & Ford, J.K. (1988). Transfer of Training: A Review and Directions for Future Research. *Personnel Psychology*, 41, 63–105.
7. Holton, E.F., Bates, R.A. & Ruona, W.E.A. (2000). Development of a Generalized Learning Transfer System Inventory. *HRDQ*, 11(4), 333–360.
8. ADL Network. xAPI Specification v1.0.3 / IEEE 9274.1.1-2023. https://adlnet.gov/
9. Rustici Software. SCORM.com — SCORM Versions overview. https://scorm.com/
10. xAPI.com — cmi5 Specification (AICC). https://xapi.com/cmi5/
