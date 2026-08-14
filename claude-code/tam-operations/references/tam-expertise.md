<!-- hub-reference-banner -->
> **Reference file — part of the `tam-operations` hub.** Formerly the standalone `tam-expertise` skill.
> Sibling topics in this family are now reference files under the hubs (`tam-operations`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: tam-expertise
description: >-
  Structuring and reviewing TAM account deliverables — which framework to apply,
  operationalizing health scores, and common review mistakes. Use when structuring
  or reviewing an EBR, QBR, account review, post-mortem, or runbook; choosing a framework
  (BLUF, STAR, Pyramid, SCQA, Diataxis); operationalizing RYG health scores;
  sanity-checking churn-risk and sentiment framing. Framework selection and review
  heuristics, not metric formulas. TRIGGER: "which framework for this EBR", "structure my
  QBR", "review this account deliverable", "is my RYG score defensible", "BLUF vs
  Pyramid". SKIP: MongoDB Premium Services operating procedures → tam-operations
  (tam-reference); case management and TS Tools → tam-operations (case-tracker); Atlas
  cluster diagnostics → atlas-diagnostics-expert; code/implementation work →
  domain-specific skills; behavior-change models (Self-Determination Theory, Fogg B=MAP)
  → applied-psychology (behavior-change-psychology).
origin: local
version: 1.2.0
category: custom
updated: "2026-06-01"
whenToUse:
  - "choose the right communication framework for an EBR or QBR"
  - "structure or review an account deliverable for document soundness"
  - "operationalize a RYG health score so each state has a mandatory action"
  - "sanity-check churn-risk or sentiment framing in a deliverable"
  - "structure a stakeholder communication (BLUF, Pyramid, STAR, SCQA)"
  - "apply Diataxis discipline to a runbook or reference doc"
  - "layer a deliverable for both engineer and executive audiences"
related_skills:
  - tam-operations
  - atlas-diagnostics-expert
  - tam-account-reports
  - operator-report-generator
---

# TAM Expertise — Deliverable Structure, Framework Selection & Review Heuristics

Framework-selection and review heuristics for TAM account deliverables — which communication framework fits a document, how to operationalize health scores, and the common mistakes to catch in review.

> **Content scope:** this reference carries framework selection, structural discipline, and review heuristics only. It does not contain step-by-step metric formulas (NRR, MEDDPICC scoring, named SaaS benchmarks); a fuller `tam-expertise-context.md` companion once held that depth but is not present. For MongoDB Premium Services operating-procedure depth, load `tam-operations (references/tam-reference.md)`.

## When NOT to use

- MongoDB Premium Services operating procedures → tam-operations (references/tam-reference.md)
- Active case management and TS Tools API → tam-operations (references/case-tracker.md)
- Atlas cluster diagnostics and troubleshooting → use `atlas-diagnostics-expert`
- Code review, frontend design, or implementation tasks → use domain-specific skills

## When to use

- Structuring or reviewing account deliverables (EBRs, QBRs, architecture reviews, post-mortems, runbooks, migration guides, case notes) for document soundness
- Choosing a communication framework (BLUF, STAR, Pyramid, SCQA, Diataxis) for a given deliverable
- Operationalizing RYG health scores and sanity-checking churn-risk or sentiment framing
- Layering a deliverable so it serves both engineer and executive audiences

## Skill guidance

- Apply the Diataxis framework for document structure and BLUF/Pyramid/STAR for communication.
- Use named frameworks and checklists before improvising structure.
- Cross-reference `tam-operations (references/tam-reference.md)` for MongoDB Premium Services specifics (roles, SLAs, JIMP, lifecycle).
- Cross-reference `tam-operations (references/case-tracker.md)` for active case-management specifics.
- SaaS benchmarks are 2025-2026 vintage; verify current figures for customer-facing deliverables.

## Common mistakes

- **Mixing Diataxis types:** a runbook (how-to) that drifts into explanation, or a reference doc that tries to teach via tutorial.
- **Unoperationalized RYG scores:** every score state needs a mandatory action, not just a color.
- **Benchmarks without recommendations:** presenting a number without a prescriptive next step.
- **Score inflation:** defaulting accounts to Green without data-backed justification.
- **Single-audience framing:** writing for engineers when the deliverable serves both engineers and executives. Use layered structure: exec summary → findings → technical appendix.
