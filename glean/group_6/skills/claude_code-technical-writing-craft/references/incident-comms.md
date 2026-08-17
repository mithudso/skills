<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `incident-comms` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: incident-comms
description: The writing surface of incident communication — status-page entries, customer-facing incident updates, internal SEV channels, and the cadence/voice shifts across the incident timeline. Covers the "what we know / what we're doing / next update at" template, hour-zero vs hour-6 vs resolution voice shifts, the heartbeat rule, internal vs external audience separation, SEV1 stakeholder notifications, and the "blameless transparency without legal exposure" balance. TRIGGER: "write a status-page update", "draft a customer incident message", "we need to communicate this outage", "what do we tell customers", "SEV1 update template", "incident update for execs", "post a heartbeat update", "incident is still ongoing what do we say", "draft the all-clear message", "resolution communication". SKIP: the incident-response process itself (severity, command, escalation — use incident-response), runbook-style execution procedures (use runbook-craft), postmortem authoring (use postmortem-writing), exec-only audience strategy (use executive-comms), general prose tone (use writing-expert).
---

# Incident Comms

## Overview

Incident communication is the writing that happens while the fire is burning. It is distinct from runbooks (which are read by operators executing a procedure), from postmortems (which are read after the incident), and from the incident response process itself (severity classification, command structure, roles). This skill governs the prose that goes out the door — to customers, to status pages, to internal SEV channels, to executives, to the support queue — between minute zero and the all-clear.

The defining constraint: you are writing under uncertainty, on a clock, with multiple audiences who need different information at different times, and every word may end up in a regulatory filing or a customer's lawsuit. The Atlassian incident-communication handbook, PagerDuty's response documentation, and incident.io's stakeholder-communication guidance converge on the same architecture: predictable cadence, plain language, no speculation, separate internal vs external surfaces, and one designated voice (usually the incident commander or a comms lead).

Use this skill when:

- Drafting a status-page entry for a live incident.
- Writing customer-facing email or in-product banner copy about an active outage.
- Posting heartbeat updates in an internal SEV channel.
- Crafting executive briefings during a SEV1.
- Writing the "all-clear" / resolution message.
- Building reusable templates for future incidents.
- Reviewing draft incident communication for tone, accuracy, and legal exposure.

Skip this skill when:

- You are writing the incident response *process* — severity definitions, command roles, escalation paths. Use `incident-response`.
- You are writing the runbook the operators will execute. Use `runbook-craft`.
- You are writing the postmortem afterward. Use `postmortem-writing`.
- You are writing strictly to executive audiences with no operational layer. Use `executive-comms`.
- You are doing general prose work with no incident context. Use `writing-expert`.

## Core Concepts

### 1. The "what we know / what we're doing / next update at" triple

Every incident update — internal or external, hour zero or hour six — answers three questions in this order:

1. **What we know.** The observed impact and known scope. No speculation. No causes unless confirmed.
2. **What we're doing.** The current mitigation activity. Not the long-term fix; just the next action.
3. **Next update at.** A specific timestamp the reader can plan around.

This template appears, with minor wording shifts, across Atlassian's incident-communication tutorial, PagerDuty's response documentation, incident.io's stakeholder-communication guide, UptimeRobot's observability hub, and the Statuspage incident-tips documentation. It works because each clause maps to a distinct reader anxiety: *what is broken*, *who is on it*, *when will I hear more*. A message that omits any one of the three breeds escalation calls.

A good update is short. Atlassian's templates run 1-3 sentences. The "next update at" line is non-negotiable even when the answer is "we still don't know" — that line is the heartbeat.

### 2. The heartbeat rule and predictable cadence

For SEV1: post an update every 30 minutes. For SEV2: every 60 minutes. For lower severities: at meaningful state changes. The rule that travels across PagerDuty, Atlassian, incident.io, and Hyperping: never let more than the cadence interval pass in silence, even if the update is "we are still investigating, no new information, next update in 30 minutes."

Silence reads as either incompetence or hidden bad news. A predictable heartbeat — even saying nothing new — is reassurance. Customers who can rely on a 30-minute cadence stop opening tickets to ask if you noticed.

The heartbeat is also a forcing function on the responders: every 30 minutes the incident commander has to assemble what's actually known and write it down, which catches drift between what people think they know and what they can defend in print.

### 3. Voice shifts across the incident timeline (hour zero → resolution → postmortem)

The same incident requires three different tones over its lifecycle:

- **Hour zero (acknowledge):** terse, acknowledging, forward-leaning. "We are aware. We are investigating. Next update at HH:MM." No apologies yet (they read as performative). No technical detail (you do not have any yet that you can defend).
- **Mid-incident (heartbeat, hours 1 through resolution):** factual, slightly more detail as scope becomes clear. Add quantifiable impact ("affecting customers in us-east-1," "approximately 15% of API requests failing"). Keep "what we're doing" present-tense and active.
- **Resolution (all-clear):** acknowledge the resolution, give the time it was fixed, briefly state the immediate cause if confirmed, and explicitly commit to a postmortem with a date. This is the first time apology language is appropriate, and it should be specific to the impact, not generic.
- **Postmortem (separate skill):** detailed, retrospective, blameless, root-cause oriented. Goes out 24-72 hours later. Owned by `postmortem-writing`.

The Atlassian, PagerDuty, and incident.io guides all reinforce this tonal progression. The most common writing mistake is using mid-incident voice at hour zero (too much detail you cannot defend) or resolution voice mid-incident (apologizing while the system is still down reads as defeat).

### 4. Internal vs external surfaces — separate by default

A SEV1 generates communication across at least four surfaces, and they must be written separately:

- **External status page** (Atlassian Statuspage, Statuspal, etc.) — public, customer-readable, neutral and reassuring tone, minimal technical detail.
- **Customer email / in-product banner** — sent to subscribed or impacted customers, slightly more direct ownership of impact.
- **Internal SEV channel (Slack/Teams)** — frank, detailed, real-time, can include hypotheses and "we don't know yet" language because the audience is responders.
- **Executive briefing** — concise, business-impact focused, sent every 30-60 minutes to leadership separately from the operational channel.

PagerDuty's stakeholder-communications guide is explicit on this separation: messages routed to specific channels for specific audiences, with predefined SEV-1/SEV-2 stakeholder groups. The Atlassian incident-management handbook uses separate status pages for internal staff and external customers.

The cardinal mistake is reusing a Slack message verbatim on the status page. Slack tolerates speculation ("might be the load balancer"), the status page does not. The status page tolerates calm reassurance; Slack needs urgency.

### 5. No speculation, no causes until confirmed

The single most common litigated phrase in incident communication is "due to" or "caused by" in an early-incident update that turns out to be wrong. Atlassian's tutorial, PagerDuty's response docs, and the incident.io best-practices guide all warn: do not name a cause until it is confirmed beyond reasonable doubt, and even then, prefer "we have identified an issue with X" over "X caused Y."

The taxonomy of permissible language by phase:

- **Hour zero**: "We are investigating reports of …" / "We are aware of …" / "Symptoms include …"
- **Mid-incident, before root cause is confirmed**: "We have identified the affected component as …" / "Mitigation is in progress."
- **Mid-incident, after root cause is confirmed**: "We have identified the cause as …"
- **Resolution**: "The incident was caused by … and has been resolved by …" — only if the postmortem-grade root cause is established. Otherwise: "We have restored service. The full root-cause analysis will be published within 48 hours."

This is not weasel-wording. It is the discipline of not making claims under uncertainty. Speculation is what creates the awkward "earlier we said it was the database, it was actually the cache" follow-up, which destroys trust.

### 6. Quantified impact, not vague gestures

"Some customers may be experiencing issues" is uninformative. The reader cannot tell if they are affected or not.

Better: "Approximately 12% of API requests in us-east-1 are returning 503s as of HH:MM. Customers in eu-west-1 and ap-south-1 are unaffected. Read operations are succeeding; write operations are degraded."

Quantification gives the customer the data to decide whether to act. It also signals that you actually have telemetry on the incident. The Atlassian tutorial and incident.io guides both push toward concrete impact statements, and Statuspage's incident-tips documentation recommends naming the affected components by their public-status-page names so subscribers know whether their service is involved.

If you cannot quantify yet, say so: "We are investigating and will quantify impact in the next update at HH:MM."

### 7. Apology calibration — specific, late, owned

Apology placement is the single most-debated topic in incident-comms guides. The consensus across Atlassian, PagerDuty, and incident.io: apologies belong at resolution, not at acknowledgment. A premature apology reads as either performative or as conceding fault before the facts are in.

A good resolution apology is specific to the actual harm: "We are sorry for the disruption to customers running write workloads against us-east-1 between 14:32 and 15:18 UTC." Not "we apologize for any inconvenience this may have caused" (generic, hedged, untrusted).

Apologies in customer-facing email at resolution should also commit to two things: the postmortem deliverable and what the customer can expect going forward (credit policy, mitigation, etc., if applicable). This is the bridge from incident comms into the postmortem timeline.

### 8. Designated voice, single source of truth

During a SEV1, exactly one person writes the external comms. Two writers in parallel will diverge — different word choices, different impact estimates, different next-update times — and the divergence will surface to customers and create a credibility crisis.

The PagerDuty incident-response documentation assigns this to a comms role (sometimes the incident commander, sometimes a dedicated customer-liaison). The Atlassian incident commander role is responsible for setting up communication channels and inviting the right people, with stakeholder communication owned by the IC or delegated to a comms lead.

The internal SEV channel can have many voices. The external surfaces must have one.

### 9. The "subscribe" affordance and acknowledging the silent audience

Status pages exist partly to deflect inbound support load. The Atlassian Statuspage documentation emphasizes the subscribe path (email, SMS, webhook) so that customers can opt into push updates rather than refreshing the page. Every status-page entry should mention the subscribe path on the first major incident in a session, and update copy should be optimized for the email/SMS notification, not the rendered web page.

The silent audience — customers who don't open a ticket — is the majority. They read your status page, they decide whether you are competent and honest, and they renew or churn quietly. The status-page voice is calibrated to that silent reader, not to the loud one.

### 10. The all-clear has a checklist, not just a message

The resolution update is the most-read message of the incident lifecycle. It needs to do five things at once:

1. **State that service is restored** with the timestamp.
2. **Confirm no further customer action is required** (or specify what is required, e.g., "please retry any failed requests").
3. **Briefly characterize the cause** (component-level, not blame-level).
4. **Commit to a postmortem** with a date ("we will publish a full root-cause analysis within 72 hours").
5. **Apologize specifically** for the actual impact.

A common failure is publishing a resolution that says only "the incident has been resolved" — the reader has no idea if their data is intact, if they need to retry, if a postmortem is coming, or whether the team understands what happened.

## Templates and Patterns

### Central artifact: a status-page heartbeat update

```markdown
**[INVESTIGATING] 14:32 UTC — API errors in us-east-1**

We are investigating elevated error rates affecting approximately 12% of API
requests in our us-east-1 region. Read operations are completing normally;
write operations are returning 503 errors. Customers in eu-west-1 and
ap-south-1 are unaffected.

Our engineering team is actively mitigating. We will post the next update at
15:00 UTC.

Affected components: API (us-east-1), Write Operations (us-east-1)
```

This template embeds the three required clauses (what we know / what we're doing / next update at), names affected components by their status-page names, quantifies impact, and commits to a heartbeat.

### Hour-zero acknowledgment (terse, no speculation)

```markdown
**[INVESTIGATING] HH:MM UTC**

We are aware of reports of [observable symptom]. Our engineering team is
investigating. We will post an update within [cadence] minutes, by HH:MM UTC.
```

### Mid-incident heartbeat (no new news, still on it)

```markdown
**[INVESTIGATING] HH:MM UTC — update**

We continue to investigate [symptom]. Mitigation is in progress; no change in
customer impact since the previous update. Next update at HH:MM UTC.
```

The "no new information" heartbeat is not filler. It is the contract.

### Mitigation update (cause identified, working the fix)

```markdown
**[IDENTIFIED] HH:MM UTC — issue isolated**

We have identified the affected component as [component]. Mitigation is in
progress and we are seeing initial signs of recovery. Customer impact remains
[quantified]. Next update at HH:MM UTC.
```

### Resolution / all-clear

```markdown
**[RESOLVED] HH:MM UTC — service restored**

Service was restored at HH:MM UTC. All customer-facing operations have
returned to normal. No customer action is required; any previously failed
requests can now be retried successfully.

Between HH:MM and HH:MM UTC, approximately [quantified impact] customers
experienced [impact description]. The incident was triggered by [neutral
component-level description].

We will publish a full root-cause analysis within 72 hours. We are sorry for
the disruption.
```

### Internal SEV-channel heartbeat (different audience, different voice)

```markdown
[14:32 UTC] *IC update — SEV1 #inc-0421*
- Impact: 12% 5xx on us-east-1 writes, holding steady
- Working hypothesis: degraded primary in db-shard-3; not confirmed
- Current action: @alice running rs.stepDown, ETA 5 min
- Comms: external status page updated at 14:30, next external update 15:00
- Blockers: none
- Next IC update: 14:45
```

Internal channel tolerates the hypothesis, names the responder, and links the comms cadence to the external one.

### Executive briefing (SEV1, every 30-60 min)

```markdown
**SEV1 exec update — 14:32 UTC**

- Customer impact: ~12% of us-east-1 write traffic failing. EU/APAC unaffected.
- Customers most affected: [top 3 by traffic, if known]
- Time to mitigation: estimated 30-45 min from now (14:45-15:00 UTC).
- Public comms: status page updated; next public update 15:00 UTC.
- Risk: if mitigation fails, escalation path is a regional failover; ~2hr RTO.
- Next exec update: 15:00 UTC.
```

Exec format strips operational detail and surfaces business-impact decisions.

## Anti-Patterns

- **The vanishing update**: posting an incident, then going silent for 90 minutes. The heartbeat is the contract. Even "no news" is content.
- **Premature blame**: naming a vendor, a team, or a system as the cause in the first 15 minutes. Almost always wrong, almost always litigated.
- **Generic apology theater**: "We apologize for any inconvenience" with no acknowledgment of actual impact. Reads worse than no apology.
- **Marketing voice in an outage**: "Thank you for your patience as we work to deliver an even better experience." This is the cardinal sin. Customers' systems are down. They do not want to be thanked.
- **Reusing Slack copy on the status page**: speculation that is fine internally becomes a press quote externally.
- **Quantification refusal**: "some customers may be affected" when telemetry exists to say "12% in us-east-1." Lazy and erodes trust.
- **Missing the silent audience**: writing only to the loudest customer who paged the exec, not to the 10,000 silent readers of the status page.
- **No clock on the next update**: "we'll update as we learn more" — guarantees inbound support load every 5 minutes.

## Decision Heuristics

- **When to declare an external incident**: customer-visible impact, even small, lasting longer than your cadence interval. If you are debating whether to post, post.
- **When to upgrade voice from "investigating" to "identified"**: when the affected component is confirmed with telemetry, not when an engineer has a hunch.
- **When to apologize**: at resolution, specifically about the actual impact. Not at acknowledgment. Not generically.
- **When to commit to a postmortem in the resolution message**: always for SEV1/SEV2, never with a date later than 5 business days, never without naming the deliverable.
- **When to switch from public status-page to direct customer email**: when impact is concentrated on identifiable customers (regional, named-account, enterprise-tier) and the status page would under-serve them. Run direct email in parallel with the status page, not instead of.
- **When to hand off to `executive-comms`**: when the audience is exclusively C-suite and the operational layer is irrelevant. For mixed audiences, stay in `incident-comms` and produce a separate exec briefing.
- **When to hand off to `postmortem-writing`**: at resolution + 24 hours. The resolution message is the bridge; the postmortem is a different document.
- **When to fall back on a templated message**: always for hour-zero acknowledgment and heartbeat. Custom prose is for resolution and root-cause-confirmed updates. Templates plus a designated voice eliminate the "draft-by-committee in a Slack channel" failure mode.

## References

- [Atlassian — Incident communication best practices](https://www.atlassian.com/incident-management/incident-communication) — the canonical reference for status-page cadence, voice calibration, and the heartbeat rule.
- [Atlassian — Incident communication templates and examples](https://www.atlassian.com/incident-management/incident-communication/templates) — concrete templates for investigating / identified / monitoring / resolved phases.
- [Atlassian — Learn incident communication with Statuspage](https://www.atlassian.com/incident-management/tutorials/incident-communication) — tutorial covering internal vs external surface separation.
- [Atlassian — Incident Communication Templates (PDF)](https://www.atlassian.com/dam/jcr:fc1c3565-32f8-47a5-ac3d-bbf594068037/incident_communication_templates.pdf) — Atlassian's internal handbook templates, open-sourced.
- [Statuspage — Incident communication tips](https://support.atlassian.com/statuspage/docs/incident-communication-tips/) — affected-component naming, subscribe affordances, tone calibration.
- [Statuspage — Create an incident template](https://support.atlassian.com/statuspage/docs/create-an-incident-template/) — operational mechanics of reusable templates.
- [PagerDuty — Stakeholder Communication](https://www.pagerduty.com/platform/incident-management/stakeholder-communication/) — internal stakeholder cadence, predefined SEV-1/SEV-2 groups, exec-comms separation.
- [PagerDuty — During an Incident](https://response.pagerduty.com/during/during_an_incident/) — the role of designated comms voice; what to post when.
- [PagerDuty — Internal Stakeholder Communications Guide](https://stakeholders.pagerduty.com/) — full open-source playbook for internal comms during major incidents.
- [PagerDuty — Different Roles](https://response.pagerduty.com/before/different_roles/) — incident commander vs. comms lead vs. customer liaison.
- [incident.io — Incident communication best practices: Keep stakeholders informed](https://incident.io/blog/incident-communication-best-practices) — quantified impact, designated voice, cadence enforcement.
- [Hyperping — 7 Incident Communication Templates (+ Best Practices)](https://hyperping.com/blog/incident-communication-templates) — concrete template library.
- [UptimeRobot — Incident Communication: Strategy, Best Practices & Templates](https://uptimerobot.com/knowledge-hub/observability/incident-communication-guide/) — the "Golden Hour" framing, hour-zero discipline.
- [Atlassian — The role of the incident commander](https://www.atlassian.com/incident-management/incident-response/incident-commander) — IC ownership of stakeholder communication.
- [Atlassian — Incident Management Handbook](https://www.atlassian.com/incident-management/handbook) — full handbook covering process and comms together.
