<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `runbook-craft` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: runbook-craft
description: Runbook-specific writing craft — execution-safety constraints for documents that must work under pressure on a sleep-deprived on-call engineer at 03:00. Covers numbered atomic steps, "you are here" markers, decision-tree branching with measurable thresholds, prereq blocks, rollback steps, post-condition checks, the "test on a fresh machine" rule, ownership/review cadence, and common runbook anti-patterns. TRIGGER: "write a runbook", "draft a runbook for X", "this runbook is unclear", "DR runbook", "playbook for the on-call", "on-call procedure", "incident response runbook", "rollback procedure", "fix this runbook", "make this runbook execution-safe". SKIP: general technical writing tone or word-level craft (use technical-writing-craft), the incident process itself (use incident-response), status-page and customer-facing incident communication (use incident-comms), postmortem authoring (use postmortem-writing), plain-language translation for accessibility (use plain-language).
---

# Runbook Craft

## Overview

A runbook is not a piece of documentation. It is a procedural script someone must execute correctly while tired, under pressure, with paging alerts firing in the background. That distinction changes everything about how it should be written.

`technical-writing-craft` and `writing-expert` cover tone, sentence flow, paragraph cohesion, and the Williams Given/New contract. Those rules still apply, but a runbook adds an execution-safety constraint: every step must be unambiguous to a person who did not write it and may never have run it before. The reader is doing, not learning.

Use this skill when you are:

- Authoring a new operational procedure for an on-call engineer, SRE team, or support engineer.
- Writing a disaster-recovery (DR) runbook that may sit unused for 18 months and then be executed once at 02:00 under load.
- Documenting a deployment or rollback procedure that mutates production state.
- Auditing an existing runbook that "didn't work" during a real incident.
- Converting tribal knowledge ("Bob always restarts the broker first") into a written procedure.
- Building a runbook that automation will eventually consume (the wording must be machine-parseable).

Skip this skill when:

- You are writing reference documentation, architectural overviews, or tutorial content — those are different genres. Use `technical-writing-craft` or `writing-expert`.
- You are writing the process of incident response itself (escalation, severity, command structure). Use `incident-response`.
- You are writing customer-facing or status-page incident updates. Use `incident-comms`.
- You are writing a postmortem. Use `postmortem-writing`.

## Core Concepts

### 1. The "fresh machine" test

A runbook is only correct if a person who has never run it before, on a freshly provisioned environment, with no tribal context, can complete it successfully. Veeam's runbook documentation states this directly: a runbook nobody has executed is still theory. The Flux CD DR documentation makes the same point — full runbook execution in an isolated environment tells you whether the actual document is usable.

Operationally this means: schedule a quarterly drill where someone who did not author the runbook runs it end-to-end on a sandbox or staging clone. Capture every place they paused, asked a question, or guessed. Each pause is a defect in the runbook, not in the runner.

The fresh-machine test is the only honest QA gate for procedural docs. Reading them does not count. Executing them does.

### 2. Atomic, numbered steps with a verb-first imperative

Each step performs exactly one action that produces exactly one verifiable result. The verb comes first. The expected outcome follows immediately.

Bad: "Now we need to make sure that the broker is running and you may also want to check the lag, and if the lag is high then restart things."

Good:
1. Run `kafka-broker-api status --broker mdb-prod-1`. Expected output: `STATUS: HEALTHY`.
2. Run `kafka-consumer-groups --describe --group mdb-tam-consumer`. Record the `LAG` column.
3. If `LAG > 50000`, go to step 7 (broker restart). Otherwise continue to step 4.

Each step has one verb (Run, Record, If…go to). Each step has one outcome the runner can verify before moving on. The runbook author has done the conditional logic in advance — the runner does not have to.

The Nobl9 and incident.io runbook guides converge here: each step should be an actionable command, not a narrative description, and the words between the action and the expected outcome should be as few as possible.

### 3. "You are here" markers and progress anchoring

In any runbook longer than ~10 steps, the runner will lose their place. They will get a Slack ping, a phone call, or a coworker asking a question. They need to come back and know exactly where they were.

Three techniques anchor progress:

- **Section banners** at the top of each major phase: `=== PHASE 2 of 5: failover the primary ===`.
- **State-check steps** at the boundary of each phase: "Before proceeding to Phase 3, confirm the replica is in `SECONDARY` state."
- **Numbered top-level steps that never restart**: do not use 1-5 in three different sections; use 1-25 across the whole runbook so the runner can say "I'm on step 14" with no ambiguity.

The DR-runbook literature (Databarracks, AllConnected) emphasizes this because DR drills are long, multi-hour procedures where context loss is inevitable.

### 4. Prerequisites block at the top, before step 1

A runbook must declare, before any action step, exactly what the runner needs to have in hand. This block is non-negotiable. If the runner does not have these things, they should stop and acquire them, not improvise during execution.

A complete prerequisites block contains:

- **Access**: which SSO group, which IAM role, which secrets vault entry, which kubectl context, which VPN.
- **Tools and versions**: `mongosh >= 2.0`, `aws-cli >= 2.13`, `jq`, `kubectl`. Pin versions when behavior diverges.
- **Inputs**: cluster ID, account ID, ticket number, change-control window.
- **Approvals**: who must sign off in writing before step 1, where the approval lives.
- **Communication**: which Slack channel to post in, who to page if blocked.

The PagerDuty runbook structure guide is explicit: required access/permissions and tools needed must be enumerated at the top, not discovered mid-execution.

### 5. Rollback as a first-class section, defined before the change

Any step that mutates production state must have a paired rollback. The rollback is written **before** the change is made, not improvised when the change fails. Several runbook guides converge on this: "rollback criteria should be defined before deployment, not during an incident, and written into the deployment runbook so the on-call engineer does not have to make judgment calls under pressure."

A rollback section answers four questions:

1. **What signals trigger a rollback?** Quantitative thresholds, not feelings. ("Error rate > 2% sustained for 5 minutes." Not "if things look bad.")
2. **What is the rollback command?** Exact, copy-pasteable.
3. **What is the rollback verification?** How the runner confirms the rollback succeeded.
4. **What is the data-loss / state-loss implication?** If the rollback truncates a queue or drops in-flight requests, say so explicitly.

If a change cannot be rolled back, that fact must be stated in BIG LETTERS at the top of the change steps. The runner must understand they are crossing a one-way door.

### 6. Decision points with measurable thresholds, not vague assessments

Decision branches save runbooks from becoming a wall of "well it depends." But branches must trigger on numbers the runner can read off a screen, not on judgments.

Bad: "If memory looks high, restart the service."

Good: "If `mem_used_pct > 85` for 3 consecutive samples, restart the service (step 12)."

The Nobl9 and Rootly runbook guides converge: decision points should branch responders to different actions based on measurable thresholds, and if a step reads "if memory looks high," it should be rewritten with a number.

A corollary: do not let a runbook grow more than ~3 levels of nested branching. If it does, split it into separate runbooks per scenario and link them. The incident.io runbook automation guide warns that too many branches is a smell the runbook needs to be split.

### 7. Post-condition checks (assertions) at the end of each phase

After every state-changing step or phase, the runbook must specify what the runner checks to confirm the change took effect. Without this, the runner can be 12 steps past a silent failure before noticing.

A post-condition check has three parts:

- **The command to run** (or signal to observe).
- **The expected result** (exact string, numeric range, dashboard graph shape).
- **What to do if the result does not match** (rollback, escalate, retry).

Example:
> **Post-condition (after step 8)**: Run `mongosh --eval "rs.status().members[0].stateStr"`. Expected: `PRIMARY`. If `SECONDARY` or `RECOVERING`, page the database on-call (PagerDuty service `mdb-data`) and do **not** proceed to step 9.

This is the runbook equivalent of unit-test assertions. They make silent failures impossible to miss.

### 8. Ownership, review cadence, and metadata at the top

Every runbook needs:

- **Owner** (a team, not a person — people leave).
- **Last reviewed** (date + name).
- **Next review due** (a real calendar date).
- **Linked alert / page / dashboard** (the runbook is useless if the on-call cannot find it from the alert).
- **Estimated duration** (so the runner knows whether to start at 23:55 or wait until tomorrow).
- **Risk level** (read-only / mutates non-prod / mutates prod / irreversible).

Google's SRE on-call material is explicit: playbook entries should be updated when the corresponding page fires. The Atlassian incident management handbook and PagerDuty runbook guides converge: runbook reviews belong in sprint planning or quarterly retrospectives, with an owner.

Stale runbooks are worse than missing runbooks because they consume trust. A "last reviewed: 18 months ago" header lets the runner calibrate skepticism.

### 9. Plain copy-pasteable commands, no placeholders in prose

The single highest-defect-density issue in real runbooks is `<placeholder>` syntax buried in a paragraph of prose. The runner copy-pastes the literal string, the command fails, they lose 90 seconds.

Two techniques fix this:

- **Code-fenced blocks** for every command, with no surrounding prose inside the block.
- **A "set variables" step at the top** that declares all the substitutions once and uses environment variables thereafter. So step 1 says `export CLUSTER_ID=mdb-prod-01` and every later step uses `$CLUSTER_ID`. The runner sets it once.

This also makes the runbook copyable into a terminal as a shell script for dry-running or automation.

### 10. Common runbook anti-patterns

Synthesizing across Nobl9, Splunk, Rootly, incident.io, Tines, IncidentHub, and Wikipedia's runbook entry:

- **The narrative blob**: paragraphs where steps should be. Bury one tab-stop deeper in execution-time the cost of reading 200 words of prose under page-load.
- **Hardcoded secrets** in code blocks. Never. Reference the vault path.
- **"You should know" gaps**: the runbook assumes the runner has the same context as the author. The fresh-machine test catches these.
- **Ambiguous phrasing**: "investigate the issue", "check the dashboard", "if things look wrong". Replace every such phrase with a concrete command + expected output.
- **Outdated commands**: the API version changed, the flag was renamed. The "last reviewed" header lets the runner calibrate.
- **No rollback**: change steps with no inverse documented.
- **Privately-stored runbooks**: a runbook in someone's Notes app is not a runbook. Discoverability from the alert is mandatory.
- **Decision trees with too many branches**: split into multiple runbooks linked from a dispatch table.
- **Mixing reference and procedure**: "here is how Kafka works, and now restart the broker." The runner does not need a Kafka tutorial at 03:00. Link to it; do not inline it.

## Templates and Patterns

### Central artifact: a single runbook step

```markdown
### Step 12 — Failover the primary replica

**Action**: Run the failover command against the current primary.
```bash
mongosh --host $PRIMARY_HOST --eval "rs.stepDown(60)"
```

**Expected output**: command returns within ~5s with no error. The primary returns to `SECONDARY` state.

**Post-condition check**:
```bash
mongosh --host $PRIMARY_HOST --eval "rs.status().members[0].stateStr"
```
Expected: `SECONDARY`. If still `PRIMARY` after 30s, the stepDown failed silently — go to **Rollback A** (step 25).

**Time budget**: ≤ 2 minutes. If longer, page DB on-call.

**Mutates prod state**: YES. Rollback available (step 25). Data loss: none expected; in-flight writes < 60s old may need retry.
```

### Full runbook skeleton

```markdown
# Runbook: <one-line title that matches the alert name>

| Field | Value |
| --- | --- |
| Owner team | mdb-tam-platform |
| Last reviewed | 2026-05-15 by @mitch.hudson |
| Next review due | 2026-08-15 |
| Linked alert | `PD: tam-helper-relay-down` |
| Linked dashboard | `https://...` |
| Estimated duration | 15–25 min |
| Risk level | Mutates prod state — rollback available |

## Prerequisites
- Access: `tam-prod-readwrite` SSO group, `kubectl` context `tam-prod`, VPN connected.
- Tools: `kubectl >= 1.28`, `mongosh >= 2.0`, `jq`.
- Inputs: cluster ID (from alert payload), incident channel (auto-created by PagerDuty).
- Approvals: none for SEV1 mitigation; SEV2+ requires change-control ticket.
- Comms: post in `#inc-<id>` at start and end; page `tam-eng` if blocked > 10 min.

## Variables (set once)
```bash
export CLUSTER_ID=<from alert>
export INCIDENT_CHAN=<auto>
```

## Phase 1 — Triage (steps 1–4)
1. ...
2. ...

## Phase 2 — Mitigate (steps 5–12)
5. ...

## Phase 3 — Verify (steps 13–18)
13. ...

## Rollback procedures
- **Rollback A** (revert failover): step 25.
- **Rollback B** (restore from snapshot): step 30. WARNING — data loss up to 5 min.

## Post-incident
- Update the linked dashboard annotation.
- File a postmortem ticket if SEV1/SEV2 (`postmortem-writing` skill).
- Append any lessons learned to this runbook's "Known gotchas" section.

## Known gotchas
- The `rs.stepDown` call sometimes times out at exactly 30s when the secondary is lagging. See ticket TAM-1842.
```

### Dispatch table (when one runbook is too many branches)

```markdown
# Dispatch: TAM helper-relay alerts

| Symptom in alert payload | Run this runbook |
| --- | --- |
| `helper-relay:offline` AND `worker:online` | `runbook-relay-process-crash.md` |
| `helper-relay:offline` AND `worker:offline` | `runbook-relay-host-failure.md` |
| `helper-relay:online` AND `auth:expired` | `runbook-relay-auth-rotate.md` |
| (anything else) | Page `tam-eng` and stop. |
```

A dispatch table fixes the "too many branches" smell by separating routing from execution.

## Anti-Patterns

- **Runbook as tutorial**: explaining how the system works in line with the procedure. The runner is not a student. Link out; do not inline.
- **Single mega-runbook**: 80 steps with 12 branches. Split by scenario, dispatch from a table.
- **Verbose conditionals**: "If the memory situation has become concerning..." → "If `mem_used_pct > 85` for 3 consecutive samples..."
- **No expected outputs**: telling the runner to run a command without saying what success looks like.
- **Decorative steps**: "Take a moment to consider the impact." Cut. The runner is already considering.
- **Untested rollback**: the rollback section exists but has never been exercised. Drill it.
- **Vanity ownership**: "Owned by the SRE team" with no specific squad. Pin to a Slack handle and a paging service.

## Decision Heuristics

- **When to split a runbook**: more than 3 levels of branching, more than ~40 atomic steps, or two different audiences (operator vs. responder). Split, then dispatch.
- **When to automate vs. document**: a runbook executed > 1x/month and fully deterministic is automation-eligible. Below that, keep it human-readable and let the human stay in the loop.
- **When to mark a step "stop and escalate"**: any condition the runbook author did not anticipate, any post-condition mismatch that does not match a documented rollback, any prerequisite that turns out to be missing mid-execution.
- **When to retire a runbook**: the underlying alert hasn't fired in 12 months and the system has changed. Either retire or schedule a drill — do not let it sit stale.
- **When to call this skill from another skill**: any time `incident-response` or `tam-reference` produces a procedural artifact that must be executable, route to `runbook-craft` for the writing pass. Any time `writing-expert` is asked to write an "operational doc" or "on-call procedure," route here.

## References

- [Google SRE Workbook — On-Call](https://sre.google/workbook/on-call/) — Google's on-call playbook entries: severity, impact, metric, background, mitigation, discovery. Source for the "playbook should be updated when the corresponding page fires" rule.
- [Google SRE Book — Postmortem Culture](https://sre.google/sre-book/postmortem-culture/) — supplies context for why runbooks and postmortems coevolve.
- [PagerDuty Runbook Automation](https://www.pagerduty.com/platform/automation/runbook/) — modern runbook structure: title, ID, version, last updated, owner, Slack channel, PagerDuty service, estimated duration, risk level, approval requirements.
- [Atlassian Incident Management Handbook](https://www.atlassian.com/incident-management/handbook) — runbook discoverability from the alert, ownership and review cadence.
- [Nobl9 — Runbook Example: A Best Practices Guide](https://www.nobl9.com/it-incident-management/runbook-example) — numbered atomic steps, verb-first imperatives, anti-patterns.
- [Rootly — Incident Response Runbooks: Templates, Examples & Guide](https://rootly.com/incident-response/runbooks) — decision trees with measurable thresholds, dispatch patterns.
- [incident.io — Runbook automation tools 2026: the complete guide](https://incident.io/blog/runbook-automation-tools-2026-the-complete-guide) — context-aware intelligent runbook execution, human-in-the-loop.
- [Splunk — Introduction to Runbooks](https://www.splunk.com/en_us/blog/learn/runbooks.html) — runbook fundamentals, anti-patterns.
- [The No-Nonsense Guide to Runbook Best Practices](https://blog.incidenthub.cloud/The-No-Nonsense-Guide-to-Runbook-Best-Practices) — outdated instructions, ambiguous phrasing, low discoverability, lack of accountability.
- [Crispy Umbrella — How to Create a Bulletproof DR Runbook That Anyone Can Execute Under Pressure](https://crispyumbrella.ai/blog/how-to-create-a-bulletproof-dr-runbook-that-anyone-can-execute-under-pressure) — the fresh-machine test, DR drill cadence.
- [Databarracks — The Disaster Recovery Runbook](https://www.databarracks.com/resources/the-disaster-recovery-runbook/) — DR-specific runbook structure, prereq blocks.
- [AWS Well-Architected — Create DR runbooks and regularly test backup and restoration processes](https://docs.aws.amazon.com/wellarchitected/latest/video-streaming-advertising-lens/advrel05-bp02.html) — test cadence requirements for DR runbooks.
- [Veeam Community — v13: Disaster Recovery Runbooks and Documentation](https://community.veeam.com/blogs-and-podcasts-57/veeam-v13-disaster-recovery-runbooks-and-documentation-12849) — a runbook nobody has executed is still theory.
- [Wikipedia — Runbook](https://en.wikipedia.org/wiki/Runbook) — historical definitions, runbook taxonomy.
