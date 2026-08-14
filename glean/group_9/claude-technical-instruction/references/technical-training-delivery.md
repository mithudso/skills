<!-- hub-reference-banner -->
> **Reference file — part of the `technical-instruction` hub.** Formerly the standalone `technical-training-delivery` skill.
> Sibling topics in this family are now reference files under the hubs (`technical-instruction`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: technical-training-delivery
description: >-
  Deliver and facilitate technical training: hands-on labs, live workshops, participatory live coding, and developer/customer education programs. TRIGGER: workshop or lab design; live coding (pace, narrate, mistakes, sticky notes, minute cards); diagnostic MCQ assessment; faded examples or Parsons problems; teaching SQL (execution-order, WHERE vs. HAVING) or NoSQL; programming misconceptions (variable/reference, notional machines); mixed skill levels; peer instruction / Carpentries; DevRel or customer academy (onboarding, TTFV, role paths); partner/SI enablement; ILT vs. self-paced vs. blended; screencasts. SKIP: learning-science → applied-psychology; course architecture (ADDIE, Gagné, UbD) → instructional-design-course-architecture; certification psychometrics → assessment-certification-design; troubleshooting pedagogy → teaching-troubleshooting-diagnostic-reasoning; tutorial writing (Diátaxis) → technical-writing-craft; TAM deliverables → tam-operations. Family hub → technical-instruction.
version: "1.0.0"
updated: "2026-06-16"
category: education
whenToUse: >-
  Use this skill when planning or running a technical training event, workshop, or developer education program. It covers the hands-on delivery craft — what to do in the room (or virtual room), how to structure exercises, and how to build programs that drive adoption and retention. Reference it alongside instructional-design-course-architecture (for program-level sequence) and applied-psychology (for the science behind the techniques here).
keywords:
  - technical training delivery
  - hands-on labs
  - participatory live coding
  - Carpentries pedagogy
  - Greg Wilson Teaching Tech Together
  - formative assessment workshop
  - faded examples scaffolding
  - Parsons problems programming
  - SQL teaching misconceptions
  - notional machine databases
  - developer education program
  - customer academy
  - customer education onboarding
  - mixed skill levels cohort
  - sticky notes minute cards
  - peer instruction ConcepTest
  - NoSQL document model teaching
  - screencast video production
  - sandbox lab environment design
tags:
  - education
  - training
  - developer-education
  - workshop
  - carpentries
  - cs-education
  - customer-education
---

# Technical Training Delivery & Developer Education

Craft of live technical training, hands-on lab design, and developer/customer education programs. See `references/technical-training-delivery.md` for the full reference.

## Core Capabilities

- **Participatory live coding**: The ten-tip framework (Wilson et al. 2020 / Carpentries); pacing, narration, deliberate mistakes, sticky-note status system, helper ratio
- **Formative assessment**: Diagnostic MCQs targeting misconceptions, sticky notes, minute cards, one-up/one-down, peer instruction (productive disagreement zone: 30–70% correct)
- **Scaffolding progression**: Worked example → faded example → Parsons problem → open exercise; expertise-reversal caution for advanced learners
- **Teaching programming**: Variable/assignment/reference misconceptions; array misconceptions; notional machine as explicit teaching object; "predict, run, compare" exercises
- **Teaching databases**: SQL execution order misconception (the cascade failure source); four-category SQL misconception taxonomy (Miedema); WHERE vs. HAVING; ERD design misconceptions (Kallia 2025); NoSQL access-pattern-first framing vs. normalization transfer
- **Mixed skill levels**: Entry diagnostic, extension exercises for fast finishers, helper-based load balancing, avoiding "just/simply/obviously"
- **Program models**: DevRel education (content stack, attribution lag); customer academy (onboarding track, role-based paths, TTFV metrics, Forrester/Intellum 2024 benchmarks); partner/SI enablement (tiered certification gating deal registration)
- **Modality trade-offs**: ILT 70–90% completion vs. self-paced 20–30%; cohort-based for onboarding, self-paced library for reference; office hours as low-cost signal supplement
- **Lab and sandbox design**: Sandbox vs. guided lab vs. pre-built environment; setup failure mitigation; granular step completion tracking
- **Video production**: Script-first, one outcome per video, 6-minute maximum, accessibility (captions/508)

## Quick Reference: Technique Selection

| Situation | Technique |
|---|---|
| New concept for novices | Worked example + explicit mental model |
| Bridge observation → practice | 3-step faded example |
| Efficient in-workshop practice | Parsons problems |
| Real-time comprehension diagnostic | MCQ (30–70% → peer discussion) |
| Real-time status in the room | Sticky notes (green/orange) |
| End-of-session feedback | Minute cards |
| SQL execution order confusion | State order explicitly; "predict, run, compare" |
| Document data model teaching | Access-pattern-first; not normalization |
| Initial onboarding modality | ILT/cohort |
| Ongoing reference | Self-paced library |
| Program ROI proof | Trained vs. untrained cohort analysis |

## Key Adjacent Skills

- Learning science behind these techniques → `applied-psychology` (references/learning-and-expertise-psychology.md)
- Course architecture (ADDIE, Gagné, UbD) → `instructional-design-course-architecture`
- Assessment/certification psychometrics → `assessment-certification-design`
- Teaching troubleshooting and diagnosis → `teaching-troubleshooting-diagnostic-reasoning`
- Tutorial writing craft (Diátaxis) → `technical-writing-craft` (references/tutorial-writing.md), `technical-writing-craft`
- Measuring whether delivery produced behavior change or ROI (Kirkpatrick, transfer, dashboards) → `learning-measurement-evaluation`
