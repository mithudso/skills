---
name: okta-incident-response
description: >-
  Guides MongoDB responders through Okta Tier-0 incident response, fire drills, and
  high-severity case handling. Use when triaging or coordinating Okta S1/S2 incidents,
  choosing Atlas/MongoDB diagnostic paths, following the Okta outage scenario catalog,
  applying drill safety rules, or deciding when to escalate to Cloud Ops/HELP.
  TRIGGER: /ir, run ir, activate okta incident response, okta is down, okta S1/S2,
  run an okta fire drill, okta incident postmortem or case analysis.
  SKIP: configuring or securing the Okta identity product itself (SSO/OIDC, MFA, tenant
  hardening) -> okta-expert (security-review); generic MongoDB design or query tuning
  with no incident context -> mongodb-expert; general Atlas architecture planning with
  no active case -> mongodb-atlas-expert; a normal low-severity support case with no
  incident posture -> solve-case.
origin: local
version: 1.1.0
updated: "2026-06-23"
category: custom
model: claude-opus-4-8
effort: high
tags: [okta, incident-response, tier-0, fire-drill, atlas, high-severity, escalation, tam-operations]
related_skills:
  - atlas-diagnostics-expert
  - mongodb-expert
  - mongodb-atlas-expert
  - mongodb-operations-expert
  - misc-catch-all
  - solve-case
  - security-review
whenToUse:
  - "/ir or activate okta incident response"
  - "triage or coordinate an Okta S1/S2 Atlas incident"
  - "run or facilitate an Okta incident fire drill"
  - "pick the closest Okta outage scenario and first evidence surface"
  - "decide when to escalate an Okta incident to HELP/COE or break-glass"
  - "post-incident analysis of an Okta high-severity case"
  - "Okta Tier-0 readiness or gap review"
whenNotToUse:
  - "configuring or securing the Okta identity product, SSO/OIDC, MFA, tenant hardening (use okta-expert via security-review)"
  - "generic MongoDB design or query tuning with no Okta incident context (use mongodb-expert)"
  - "general Atlas architecture planning with no active case or readiness workflow (use mongodb-atlas-expert)"
  - "a normal low-severity support ticket with no incident posture (use solve-case)"
  - "pure repo setup or internal tool installation with no responder workflow decision (use 10gen)"
metadata:
  changelog: |
    2026-06-23 sko v1.0.0->v1.1.0 — Pass H 10/10->10/10 pos, 0/10->0/10 neg (predicted); 2 Medium fixed (frontmatter metadata + model/effort; description SKIP clause + okta-* disambiguation)
---

# Okta Incident Response

Use this skill for Okta Tier-0 Atlas incidents, fire drills, readiness reviews, and high-severity case analysis. Slash alias: `/ir`.

## Core posture

- Treat a credible Okta S1 or severe S2 as an incident, not a normal support ticket.
- Optimize for restoration first. Deep root cause analysis follows stabilization.
- Bias toward action and coordination, not passive triage.
- The MongoDB IR owns technical command and restoration sequencing. TAM/IC owns customer communications unless explicitly delegated.
- Keep the case record, working timeline, and next checkpoint current. Run command in the open, not in private side conversations.
- For any backend tooling or support-system calls, operate one case at a time on explicit user action only. Do not fan out or auto-refresh.
- When drafting Okta-facing content, be terse, factual, and prescriptive. Avoid fluff, apologies, or generic support language.

## Modes

Determine the mode before doing deep work.

1. Live incident
   - Real customer-impacting event or credible production threat.
2. Fire drill
   - Simulation that should be treated as real until the reveal, but without production-impacting execution.
3. Readiness / onboarding
   - Training, process hardening, checklist review, or gap assessment.
4. Historical case analysis
   - Retrospective synthesis, pattern extraction, or action review.

If the mode is unclear, start in live-incident posture and narrow only when evidence permits.

## Severity and activation

### Live incident defaults

- S1 / P1
  - Production down, broad authentication failure, severe multi-tenant latency, or a failure that broadly prevents service use.
  - Mobilize immediately.
  - Target bridge join by minute 15.
  - Break-glass and HELP/COE escalation are in scope.
- S2 / P2
  - Serious degradation, at-risk state, or meaningful customer impact that still requires active ownership and fast coordination.
  - Aggressive telemetry review and live technical engagement are expected.
- If uncertain, start high and de-escalate after evidence improves.

### Time model

Use this as the default first-15-minute operating sequence.

- 0–5 minutes
  - Acknowledge the trigger.
  - Record the symptom, timestamp, and initial impact statement.
  - State the first severity hypothesis.
- 0–10 minutes
  - Notify TAM/IC and required internal responders.
  - Pick the closest-fit scenario from `references/scenario-catalog.md`.
  - Choose the first technical evidence surface.
- 0–15 minutes
  - Get on the bridge.
  - Log the bridge link and next checkpoint.
  - State ownership, first action plan, and open blockers.

Do not wait for perfect certainty before establishing command.

## Technical routing

Before giving deep technical guidance, route to the right sibling skill.

- `atlas-diagnostics-expert`
  - Live triage, metrics, FTDC, logs, Performance Advisor, capacity, alert interpretation, escalation packaging.
- `mongodb-expert`
  - Query/index/root-cause details, replication, sharding, WiredTiger, write conflicts, oplog behavior, engine internals.
- `mongodb-atlas-expert`
  - Atlas control-plane, networking, Admin API, private connectivity, Atlas Search/Vector Search, platform features, provider/topology posture.
- `mongodb-operations-expert`
  - Backup/restore, DR, migrations, upgrade paths, security architecture, encryption, KMS, compliance, cost, data movement.
- `mongodb-kb`
  - Symptom-to-KB mapping, known issues, public customer-shareable references.
- `10gen`
  - Which internal repo or tool to use: `ts-diag`, `alexandria`, `SearchPlanIQ`, `mongolyser`, `ts-ftdc-requestor`, `devprod-mcp-router`, and related support tooling.

When a case crosses multiple domains, start with the dominant live symptom, then cross-load the secondary skill.

## Incident workflow

Follow this sequence unless the user already gave a narrower task.

1. Establish the current incident statement
   - What is happening now?
   - What customer impact is visible now?
   - What is the working severity now?
2. Pick the scenario
   - Read `references/scenario-catalog.md`.
   - Choose the nearest dominant symptom within the first 10–15 minutes.
   - Rename the scenario later if evidence changes; do not stay attached to a bad early guess.
3. Choose the first evidence surface
   - Atlas metrics / alerts / Performance Advisor
   - Logs
   - FTDC
   - Explain plans
   - Atlas control-plane evidence
   - Support case history / case timeline
4. Separate “what” from “why”
   - Communicate the current symptom and impact immediately.
   - Treat root-cause explanation as a separate track.
5. Drive a mitigation loop
   - Current symptom
   - Evidence gathered
   - Candidate mitigations
   - Blockers or decisions required
   - Next checkpoint time
6. Escalate when platform action or cross-team coordination is required
   - Use HELP/COE fast when recovery may depend on Atlas or cloud-provider intervention.
7. Capture what must survive the incident
   - Timeline
   - Scenario selected and later changes
   - Commands or actions proposed
   - What was executed vs simulated
   - Remaining risks
   - Follow-up work

## Fire-drill rules

If the work is a drill, use real incident discipline with drill safety constraints.

- Treat the initial trigger as real until the reveal.
- Do not execute production-impacting actions during the drill.
  - No failover, restart, host swap, volume change, config push, drop index, restart pod, or other production-impacting mutation.
- Do not page Okta DRIs who have not pre-acknowledged the drill window.
- Do not send external customer communications.
- After the reveal, prefix drill communications clearly as drill traffic.
- Simulate mitigation decisions and record what would have been executed.
- If a real S1 interrupts the drill, stop the drill and switch fully to live-incident mode.

## Break-glass posture

For live S1s, do not let evidence collection delay restoration.

- If Okta is already in broad impact and requests immediate Atlas-side action, prioritize the action path needed to restore service.
- Use break-glass logic only for real incidents, not drills.
- Make explicit what action is being considered, what approval/escalation is needed, and what risk is being accepted.
- Keep the distinction between proposal, approval, and execution clear in the timeline.

## Standard outputs

When asked to solve or analyze an Okta case, default to this structure unless the user asked for something else.

### Live incident brief

- Incident statement
- Working severity
- Dominant scenario
- Current customer impact
- Evidence gathered
- Actions in progress
- Blockers / decisions needed
- Escalation status
- Next checkpoint

### Historical case analysis

- Case summary
- Severity and incident-worthiness assessment
- Dominant failure pattern
- What was done well
- What slowed restoration
- Recommended playbook or tooling changes
- Follow-up actions

### Readiness / process review

- Current-state assessment
- Gaps blocking reliable execution
- What is already durable
- What needs a source-of-truth owner
- Recommended next hardening steps

## Boundaries

This skill is for Okta-specific incident response and incident-adjacent support handling.

Use another skill instead when the task is primarily:
- Configuring or securing the Okta identity product itself: SSO/OIDC, MFA, tenant or session hardening (use `okta-expert` via `security-review`). This skill handles MongoDB-side response to an Okta-customer Atlas incident, not Okta-product configuration.
- Generic MongoDB design with no Okta incident context (use `mongodb-expert`).
- General Atlas architecture planning with no active case or readiness workflow (use `mongodb-atlas-expert`).
- A normal low-severity support ticket with no incident posture (use `solve-case`).
- Pure repo setup / internal tool installation with no responder workflow decision (use `10gen`).

## Failure modes to guard against

- Spending the first 15 minutes gathering perfect evidence instead of establishing command.
- Treating a live S1 like a normal support queue item.
- Mixing customer communications and technical command ownership.
- Locking onto the wrong scenario and failing to switch when evidence changes.
- Doing broad multi-case fan-out instead of one-case-at-a-time incident handling.
- Performing production-impacting drill actions.
- Giving high-confidence RCA language before stabilization evidence exists.

## References to load on demand

- `references/scenario-catalog.md`
  - Quick symptom-to-scenario routing and first evidence surfaces.
- `references/source-basis.md`
  - Corpus basis, extraction notes, and provisional validation summary.
