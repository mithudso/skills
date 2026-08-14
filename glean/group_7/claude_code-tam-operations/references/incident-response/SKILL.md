<!-- hub-reference-banner -->
> **Reference file — part of the `tam-operations` hub.** Formerly the standalone `incident-response` skill.
> Sibling topics in this family are now reference files under the hubs (`tam-operations`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: incident-response
description: >
  Incident response lifecycle expert — severity classification (SEV0-SEV3/P0-P4), incident roles
  (IC, comms lead, scribe, SME, deputy), communication templates, blameless postmortem writing,
  SLO/SLI/SLA management and error budgets, on-call best practices, runbook design, incident
  metrics (MTTD/MTTR/MTTA/MTBF), tool comparison (PagerDuty/Grafana Cloud IRM/incident.io/Squadcast),
  chaos engineering for preparedness, anti-patterns, and maturity model.
  TRIGGER: designing incident response processes, writing or reviewing postmortems, setting up on-call
  rotations, defining severity levels, building runbooks, reviewing incident management practices,
  SLO/error budget policy, chaos engineering game days, MTTR/MTTD metrics.
  SKIP: general monitoring or alerting configuration without an incident process component;
  security incident response for compliance/forensics (use security-review,
  references/security-reviewer.md); CI/CD pipeline design without deployment incident considerations.
version: 1.1.2
updated: "2026-06-01"
category: developer
tags:
  - incident-response
  - sre
  - on-call
  - postmortem
  - runbook
  - severity
  - SLO
  - SLI
  - SLA
  - MTTR
  - MTTD
  - chaos-engineering
  - pagerduty
  - blameless
  - incident-commander
triggers:
  - incident response
  - incident management
  - on-call rotation
  - postmortem writing
  - blameless postmortem
  - severity levels SEV P0 P1
  - SLO SLI SLA error budget
  - MTTR MTTD MTTA metrics
  - runbook design
  - incident commander
  - on-call best practices
  - chaos engineering resilience
  - incident communication template
  - status page update
  - PagerDuty Opsgenie Grafana Cloud IRM
  - incident escalation
  - war room
globs:
  - "**/runbook*"
  - "**/postmortem*"
  - "**/incident*"
  - "**/on-call*"
  - "**/oncall*"
related_skills:
  - devops-infra
  - tam-reference
  - autoremediation
  - integration-clients
---

# Incident Response Expert

## When to use this skill

- Incident lifecycle design — detect-respond-recover-learn cycle
- Severity classification — SEV0-SEV3 or P0-P4 with response expectations
- Incident roles — IC, deputy, scribe, comms lead, SME, customer liaison
- Communication templates — stakeholder notifications, status page updates, executive briefs
- Postmortem writing — blameless postmortems, root cause analysis, action item tracking
- SLO/SLI/SLA management — indicators, objectives, error budget policy, burn-rate alerting
- On-call practices — rotation design, escalation policies, handoff, well-being
- Runbook design — operational playbooks, troubleshooting guides, automation levels
- Incident metrics — MTTD, MTTR, MTTA, MTBF, tracking and improvement
- Tool selection — PagerDuty, Grafana Cloud IRM, incident.io, Squadcast, Rootly, FireHydrant
- Chaos engineering — failure injection, game days, resilience validation
- Maturity assessment — auditing existing practices, identifying anti-patterns

## When NOT to use this skill

- General monitoring/alerting configuration without an incident process component
- Security incident response for compliance/forensics → use `security-review` (references/security-reviewer.md)
- CI/CD pipeline design without deployment incident considerations
- Generic project management or task tracking

---

## Incident lifecycle

```
DETECT --> RESPOND --> RECOVER --> LEARN
  ^                                  |
  +----------------------------------+
```

### Phase 1: Detection

**Detection sources (by reliability):**
1. Automated monitoring and alerting (preferred)
2. Synthetic probes and health checks
3. Customer-reported issues
4. Social media or status-page watchers
5. Internal team discovery

**Detection best practices:**
- Alert on symptoms (error rate, latency p99, throughput drop), not causes
- Monitor SLI burn rates: 1-hour 14.4x or 6-hour 6x burn rate should page
- Track MTTD — target under 5 minutes for SEV0/SEV1
- Alert quality: clear title, dashboard link, runbook link, appropriate severity, tested thresholds

### Phase 2: Response

1. Acknowledge the alert (stops re-paging, starts MTTA clock)
2. Assess severity using the classification matrix
3. Open an incident channel (#inc-YYYYMMDD-title)
4. Page the IC if SEV0/SEV1
5. Assemble the response team
6. Begin diagnostic investigation
7. Communicate initial status to stakeholders

**War room rules:** one voice at a time; IC moderates; state name and role when joining; updates not questions in main channel; use threads for side conversations.

### Phase 3: Recovery

Recovery strategies (ordered by speed):
1. **Rollback** — revert last deployment or configuration change
2. **Feature flag disable** — turn off the misbehaving feature
3. **Traffic shifting** — reroute to healthy instances or regions
4. **Scaling** — add capacity to absorb load
5. **Hotfix** — targeted fix (last resort during incident)
6. **Failover** — switch to DR environment

Do not close the incident until the SLI has returned to acceptable levels for at least 15 minutes (SEV0/SEV1).

### Phase 4: Learning

- **Within 2 hours:** IC writes brief incident summary
- **Within 24 hours:** Draft postmortem initiated
- **Within 48 hours:** Postmortem review meeting
- **Within 72 hours:** Final postmortem published with action items assigned
- **Within 30 days:** Follow-up review to verify action item completion

Full templates: [Communication templates, Postmortem template, Runbook template](./references/templates-and-checklists.md)

---

## Severity classification

### SEV-level framework

| Level | Impact | Response time | Update cadence | Example |
|-------|--------|---------------|----------------|---------|
| **SEV0** | Total outage; all users; revenue stopped | Immediate (24/7) | Every 15 min | Site down, payment failed, data corruption |
| **SEV1** | Core functionality degraded; large user segment | < 15 min (24/7) | Every 30 min | Checkout fails for subset, SSO broken |
| **SEV2** | Non-critical feature broken; workaround exists | < 1 hour (biz hours) | Every 2 hours | Reporting dashboard down, exports fail |
| **SEV3** | Cosmetic or minor; no user-facing impact | Next business day | Daily if needed | UI misalignment, log noise |

**Severity vs priority:** severity is technical (how broken?); priority is contextual (how urgently fix?). A SEV3 bug before a major demo may be P0 priority.

**Decision tree:**
```
Is the entire service unavailable? YES → SEV0
Is a core user journey broken for many users? YES → SEV1
Is a non-critical feature broken? YES → SEV2
Otherwise: SEV3
```

Any team member can escalate severity upward. Only the IC can de-escalate. When in doubt, escalate.

---

## Incident roles

### Incident Commander (IC)

The IC is the single source of truth. The IC does NOT fix the problem — they coordinate the people who do.

**Responsibilities:** declare the incident and set severity; manage the incident channel; assign all roles; drive investigation; approve repair actions; provide regular updates; decide when to escalate/de-escalate; declare resolved.

**IC anti-patterns:** debugging the issue yourself while commanding; allowing multiple people to give orders simultaneously; failing to delegate communications; silence ("silence is panic").

**IC training:** shadow 3-5 incidents before taking primary; start with SEV2 before SEV0/SEV1; run quarterly game days.

### Other roles

| Role | Responsibilities |
|------|-----------------|
| **Deputy** | Backs up IC; monitors channel; pages additional SMEs; tracks action items |
| **Scribe** | Logs all events with UTC timestamps; records decisions and reasoning; captures commands run |
| **Comms Lead** | Drafts/sends status page updates; communicates with support, executives, PR |
| **SME** | Investigates within their domain; reports findings to IC; proposes repair actions for IC approval |

**Scribe format:**
```
[14:32 UTC] @alice Deployed rollback of commit abc123
[14:35 UTC] @bob Error rate dropped from 45% to 8%
[14:37 UTC] @IC (carol) Severity downgraded from SEV0 to SEV1
```

---

## SLO / SLI / SLA management

| Concept | What it is | Example |
|---------|-----------|---------|
| **SLI** | Quantitative metric measuring service behavior | Request success rate, p99 latency |
| **SLO** | Internal reliability target for an SLI | 99.9% success rate over 30 days |
| **SLA** | Contractual commitment with consequences | 99.5% uptime or 10% credit |

**Rule:** SLO target > SLA floor. If SLA guarantees 99.5%, SLO should be 99.9% to give teams time to detect and fix drift before breach.

### Error budgets

```
Error Budget = 1 - SLO target
If SLO = 99.9%, Error Budget = 0.1% = 43.2 minutes/month
```

| Budget remaining | Action |
|-----------------|--------|
| > 50% | Normal velocity; deploy at will |
| 25-50% | Heightened awareness; review deploys carefully |
| 10-25% | Deployment freeze for non-critical changes |
| < 10% | All hands on reliability; no feature work |
| 0% | Full freeze; postmortem required; recovery plan before any deploys |

### Burn-rate alerting

| Window | Burn rate | Action |
|--------|-----------|--------|
| 1 hour | 14.4x | Page immediately (SEV0/SEV1) |
| 6 hours | 6x | Page (SEV1) |
| 1 day | 3x | Alert / ticket (SEV2) |
| 3 days | 1x | Warning (SEV3) |

---

## On-call best practices

**Rotation design:**
- Primary + Secondary: secondary pages after 5 minutes of no acknowledgment
- Follow-the-sun for global teams
- IC rotation separate from team on-call — never be IC and team on-call simultaneously
- Maximum 1 week continuous on-call per person per month; mandatory 1 week off after

**Escalation tiers:** primary (immediate) → secondary (5 min) → team lead (10 min) → engineering manager (15 min) → VP/Director (20 min, SEV0 only). Auto-escalate unacknowledged alerts — never rely on manual escalation alone.

**Burnout prevention:** compensate on-call duty; if paged more than 2 times per night treat as a reliability problem; aim for < 5% false positive alert rate; weekly on-call review to address top paging sources.

---

## Runbook best practices

1. Write for the 3 AM responder — assume tired, stressed, unfamiliar with the system
2. Include actual commands — exact CLI commands, queries, dashboard URLs
3. Test the runbook — have someone unfamiliar follow it step by step
4. Assign ownership to a team, not a person
5. Update after every incident — fix outdated runbooks immediately as postmortem action items
6. Link from every alert — every alert should link to its runbook
7. Keep runbooks short — target 1-2 pages
8. Version control — store in Git alongside the code they support
9. Include rollback instructions for every remediation step

**Automation levels:**
- Level 0: Fully manual documented steps
- Level 1: Manual with copy-paste commands
- Level 2: Semi-automated diagnostics with suggested actions
- Level 3: Fully automated remediation with notification
- Level 4: AI-assisted agent execution with human approval

Full runbook template: [references/templates-and-checklists.md](./references/templates-and-checklists.md)

---

## Incident metrics

| Metric | Measures | Target (SEV0/SEV1) |
|--------|----------|---------------------|
| **MTTD** | Detection speed | < 5 min |
| **MTTA** | Response readiness | < 5 min |
| **MTTR** | Recovery speed | < 60 min (SEV0), < 4h (SEV1) |
| **MTTC** | Containment speed | < 30 min |
| **MTBF** | System reliability | Increasing trend |

**Use metrics to identify systemic problems — never to evaluate individual performance.** Segment by severity; aggregate MTTR across SEV0-SEV3 is meaningless.

**Goodhart's Law:** if you reward teams for low MTTR, they will close incidents prematurely and downgrade severity to avoid tracking. Reward reducing incident count and improving detection quality instead.

---

## Tool comparison (2025-2026)

| Tool | Key strength | Pricing | Notes |
|------|-------------|---------|-------|
| **PagerDuty** | 700+ integrations, AI Ops, mature | $21-39/user/mo | Industry standard |
| **Grafana Cloud IRM** | Unified with Grafana stack | ~$230/user/yr (50-user bundle) | Replaced OSS Grafana OnCall (archived Mar 2026) |
| **incident.io** | Slack-native workflow, AI postmortems | Custom | Strong for Slack-heavy orgs |
| **Squadcast** | SLO tracking built in | $9-21/user/mo | Good for mid-size teams |
| **Rootly** | AI-powered, strong retrospectives | Custom | Focus on learning |
| **FireHydrant** | Runbook automation, service catalog | Custom | Strong automation |
| **Opsgenie** | Legacy Atlassian offering | Shutting down (Apr 2027) | Migrate to Jira SM or alternatives |

**Selection guide:**
- Grafana stack → Grafana Cloud IRM
- Slack-first team → incident.io or Rootly
- Need 700+ integrations + enterprise compliance → PagerDuty
- Cost-effective mid-market → Squadcast or Better Stack

---

## Chaos engineering for preparedness

**Principles:**
1. Start with a steady-state hypothesis
2. Inject real-world failures (network partitions, disk full, dependency timeouts)
3. Run in production — staging chaos doesn't prove production resilience
4. Minimize blast radius; always have a kill switch
5. Automate experiments

**Common scenarios:**

| Scenario | Tools | What it tests |
|----------|-------|---------------|
| Pod/instance termination | Chaos Monkey, Litmus | Auto-scaling, load balancing |
| Network latency injection | Toxiproxy, AWS FIS | Timeout handling, circuit breakers |
| Dependency failure | Toxiproxy, WireMock | Fallback paths, graceful degradation |
| Disk full | dd, fallocate | Log rotation, alerting |
| Region/AZ failure | AWS FIS, manual failover | DR procedures, replication |

**Game days:** brief team on scenario (not the specific failure) → inject and respond as if real → debrief. Run quarterly for most teams; monthly for critical services.

Full experiment template: [references/templates-and-checklists.md](./references/templates-and-checklists.md)

---

## Anti-patterns

| Anti-pattern | Problem | Fix |
|-------------|---------|-----|
| Hero culture | One person always saves the day; no documentation | Rotate IC; document everything in runbooks |
| Alert fatigue | Too many alerts; responders ignore them | Tune thresholds; target < 5% false positive rate |
| Blame-and-shame postmortems | People hide mistakes; learning stops | Enforce blameless culture; focus on systems |
| Postmortem graveyard | Postmortems written but action items never completed | Track in issue tracker; review weekly |
| Severity inflation | Everything is SEV0; real SEV0 lost in noise | Enforce definitions; only IC sets severity |
| War room without IC | Everyone talks; nobody decides | Always assign IC before opening a war room |
| Silent incidents | No updates for 30+ minutes | IC timer per severity cadence |
| Premature closure | Closed before full recovery verified | Require sustained metric recovery |
| Permanent on-call | Same person always on-call; burnout | Enforce fair rotation; track page distribution |
| Alert without runbook | Engineer paged with no guidance | Link runbook to every alert |

---

## Maturity model

| Level | Name | Characteristics |
|-------|------|-----------------|
| 1 | Ad Hoc | No process; hero culture; no postmortems; reactive only |
| 2 | Defined | Severity levels; IC role; postmortems sometimes; basic on-call |
| 3 | Measured | Metrics tracked; SLOs defined; action items tracked; runbooks for top alerts |
| 4 | Managed | Error budgets enforced; chaos engineering; automated runbooks; postmortem rate > 90% |
| 5 | Optimizing | AI-assisted detection; predictive alerting; self-healing; incident count decreasing QoQ |

---

## References

- [Communication templates, Postmortem template, Runbook template, Checklists](./references/templates-and-checklists.md)
- [PagerDuty Incident Response Guide](https://response.pagerduty.com/)
- [Google SRE Book](https://sre.google/sre-book/table-of-contents/)
- [incident.io SLO/SLA/SLI Guide](https://incident.io/blog/slo-sla-sli)
- *Site Reliability Engineering* — Beyer, Jones, Petoff, Murphy
- *Chaos Engineering* — Casey Rosenthal, Nora Jones
