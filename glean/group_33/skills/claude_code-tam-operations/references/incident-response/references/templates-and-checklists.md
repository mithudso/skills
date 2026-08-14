# Incident Response — Templates and Checklists

## Communication Templates

### Initial Internal Notification

```
INCIDENT DECLARED: [Title]
Severity: [SEV0/SEV1/SEV2/SEV3]
Time detected: [YYYY-MM-DD HH:MM UTC]
Impact: [What users are experiencing]
Affected services: [List]
Incident channel: [#inc-YYYYMMDD-title]
IC: [Name]
Status: Investigating
Next update: [Time]
```

### Status Page — Investigating

```
[Service Name] - Investigating
We are currently investigating reports of [symptom description].
Users may experience [specific user-facing impact].
We will provide an update within [15/30] minutes.
```

### Status Page — Identified

```
[Service Name] - Identified
We have identified the cause of [symptom].
[Brief root cause without internal jargon.]
We are implementing a fix and expect resolution by [estimated time or "TBD"].
Next update in [15/30] minutes.
```

### Status Page — Resolved

```
[Service Name] - Resolved
The issue affecting [symptom] has been resolved as of [HH:MM UTC].
[Brief description of what was done.]
We will publish a full postmortem within [48/72] hours.
We apologize for the disruption.
```

### Executive Stakeholder Update

```
Subject: [SEV0/SEV1] Incident Update - [Title] - [Time]

Current Status: [Investigating / Identified / Mitigating / Resolved]
Business Impact: [Revenue impact, affected customer count, SLA implications]
Customer Impact: [What customers are experiencing]
Root Cause: [Known / Under investigation]
ETA to Resolution: [Estimate or "investigating"]
Actions Taken: [Bulleted list of key actions]
Next Update: [Time]
IC: [Name] | Comms Lead: [Name]
```

### Customer-Facing Email (Post-Resolution)

```
Subject: Service Disruption on [Date] - Resolved

Dear [Customer],

On [date] between [start time] and [end time] UTC, [service]
experienced [symptom description]. During this period, you may
have been unable to [specific impact].

Root cause: [Plain language explanation without internal jargon.]

What we did:
- [Action 1]
- [Action 2]

What we are doing to prevent recurrence:
- [Prevention action 1]
- [Prevention action 2]

We sincerely apologize for the disruption. If you have questions,
please contact [support channel].

Sincerely,
[Team/Company]
```

### Communication Cadence

| Severity | Internal | External | Executive |
|----------|----------|----------|-----------|
| SEV0 | Every 15 min | Every 15 min | Every 15 min |
| SEV1 | Every 30 min | Every 30 min | Every 30 min |
| SEV2 | Every 2 hours | Every 2 hours | On request |
| SEV3 | Daily | N/A | N/A |

**Communication principles:**
- Describe symptoms, not internal systems
- "Unknown at this time" beats a fabricated ETA
- Always include the next update time
- Never blame a specific person or team in external communications

---

## Postmortem Template

```markdown
# Postmortem: [Incident Title]

**Incident ID:** INC-YYYY-NNNN
**Date:** YYYY-MM-DD
**Duration:** HH:MM (start to resolution)
**Severity:** SEV0 / SEV1 / SEV2 / SEV3
**Author:** [Name]
**Reviewers:** [Names]
**Status:** Draft / In Review / Final

## Summary
[2-3 sentences: what happened, how long, what was the impact.]

## Impact
- **Users affected:** [Number or percentage]
- **Revenue impact:** [Estimated or "none"]
- **SLA impact:** [Error budget consumed, SLA breach risk]
- **Duration of user-visible impact:** [HH:MM]

## Timeline (UTC)

| Time | Event |
|------|-------|
| HH:MM | [Detection / first alert fired] |
| HH:MM | [Alert acknowledged] |
| HH:MM | [Incident declared at SEV-X] |
| HH:MM | [Key diagnostic finding] |
| HH:MM | [Mitigation action taken] |
| HH:MM | [Service restored] |
| HH:MM | [Incident closed] |

## Root Cause
[Specific technical root cause. Example: "A configuration change to the connection
pool reduced max connections from 100 to 10, causing request queueing under normal load."]

## Contributing Factors
- [Factor 1: e.g., "No automated validation of config changes"]
- [Factor 2: e.g., "Monitoring did not alert on connection pool exhaustion"]

## Detection
- **How detected:** [Monitoring / customer report / internal discovery]
- **Could we detect faster?** [Yes/No, explain]
- **MTTD:** [Minutes]

## Response
- **MTTA:** [Minutes]
- **Were the right people engaged?** [Yes/No]
- **What slowed down the response?** [If applicable]

## Recovery
- **Mitigation strategy:** [Rollback / hotfix / scaling / feature flag / failover]
- **MTTR:** [Minutes]

## Lessons Learned

### What went well
- [Thing 1]

### What went poorly
- [Thing 1]

### Where we got lucky
- [Thing 1 — surfaces hidden risks]

## Action Items

| ID | Action | Owner | Priority | Due Date | Status |
|----|--------|-------|----------|----------|--------|
| 1 | [Specific, measurable action] | @name | P0/P1/P2 | YYYY-MM-DD | Open |
```

---

## Runbook Template

```markdown
# Runbook: [Alert/Scenario Name]

**Owning Team:** [Team name]
**Last Reviewed:** YYYY-MM-DD
**Severity:** SEV0 / SEV1 / SEV2 / SEV3
**Estimated Resolution Time:** [Minutes]

## Trigger
This runbook is invoked when alert "[Alert Name]" fires.
[Description of the symptom or condition]

## Prerequisites
- Access to [system/dashboard/tool]
- Permissions: [required role or credentials]

## Diagnostic Steps

### Step 1: [Check name]
```shell
kubectl get pods -n production | grep -v Running
```
**Expected:** All pods in Running state.
**If unexpected:** Proceed to Step 2.

## Remediation

### Option A: Rollback
```shell
# Rollback commands
```
**Rollback:** [How to undo this action if it makes things worse]

## Escalation
If not resolved within [N minutes]:
1. Page [team/person] via [paging tool]
2. Escalate severity to [SEV-X]

## Verification
- [ ] [Metric] has returned to normal range
- [ ] No new alerts firing for [N minutes]
```

---

## Quick-Reference Checklists

### Incident Declaration
- [ ] Alert acknowledged
- [ ] Severity assessed using the decision tree
- [ ] Incident channel created (#inc-YYYYMMDD-title)
- [ ] IC assigned
- [ ] Deputy, Scribe, Comms Lead assigned
- [ ] Initial stakeholder notification sent
- [ ] Status page updated (if customer-facing)

### Incident Resolution
- [ ] SLI within acceptable range for 15+ minutes
- [ ] All temporary fixes documented
- [ ] Incident channel summary posted
- [ ] Status page updated to "Resolved"
- [ ] Postmortem owner assigned
- [ ] Postmortem deadline set (48-72 hours)

### On-Call Handoff
- [ ] Review active incidents and open issues
- [ ] Check recent deployments (last 24 hours)
- [ ] Review alert history and ongoing investigations
- [ ] Verify access to all required systems and dashboards
- [ ] Confirm contact information and communication channels

### Weekly On-Call Review Agenda
1. (5 min) Page count and MTTA for the week
2. (10 min) Walk through each incident or notable alert
3. (10 min) Top paging sources — which alerts need tuning?
4. (5 min) Runbook gaps — any alerts without runbooks?
5. (5 min) Action item review from previous postmortems
6. (5 min) Handoff to next on-call rotation

---

## Chaos Experiment Template

```markdown
## Chaos Experiment: [Name]

**Hypothesis:** When [failure condition], the system will [expected behavior]
  and [SLI] will remain within [SLO threshold].

**Steady State:** [Metric] = [baseline value] (+/- [tolerance])

**Method:**
1. Record steady-state metrics
2. Inject: [specific failure]
3. Observe for [duration]
4. Remove injection / verify auto-recovery

**Abort Conditions:**
- [SLI] drops below [hard threshold]
- Customer-reported errors increase

**Blast Radius:** [What is affected; what is NOT affected]

**Results:** [Pass/Fail + details]
**Action Items:** [If fail, what changes are needed]
```
