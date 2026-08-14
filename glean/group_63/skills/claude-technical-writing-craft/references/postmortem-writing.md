<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `postmortem-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: postmortem-writing
description: The writing-specific craft for postmortems — disciplined Five Whys in prose, blameless framing language, timeline reconstruction in UTC, contributing factors vs root cause distinction, action-items with owners and dates, the "what went well" balance, and the voice that lets an organization learn instead of finger-point. TRIGGER: "write a postmortem", "draft the RCA", "blameless postmortem template", "what went well section", "five whys for this outage", "post-incident review writeup", "PIR document", "retrospective writeup", "this postmortem reads as blamey — rewrite it", "structure the timeline", "action items section". SKIP: the incident response process itself or severity classification (use incident-response), live-incident customer communication (use incident-comms), runbook authoring or procedural execution docs (use runbook-craft), executive-only retro briefings stripped of operational detail (use executive-comms), general prose tone (use writing-expert).
---

# Postmortem Writing

## Overview

A postmortem is a learning artifact disguised as an incident report. It must satisfy three audiences simultaneously: the engineers who need to understand what failed, the leadership who need to evaluate organizational risk, and the people who lived through the incident and need to feel that the document represents their experience fairly. Get the writing wrong and the document becomes either a hunt for someone to blame, a sanitized PR statement, or an unread artifact filed away to satisfy a process gate.

This skill governs the writing craft of postmortems specifically. The incident-response process — declaring incidents, assigning severities, running the response — lives in `incident-response`. The customer-facing comms that go out during the incident live in `incident-comms`. This skill is what happens 24-72 hours after resolution, when someone has to sit down and write the document that will be read months or years later by people who were not there.

The defining writing constraint of a postmortem: every sentence must be defensible as factual, blameless, and pedagogically useful. Speculation, blame, and self-justification all corrode the document. The Google SRE Book Chapter 15, Etsy's Debriefing Facilitation Guide (Allspaw, Evans, Schauenberg), the Atlassian postmortem handbook, and Honeycomb's incident-retrospective practices all converge on the same core craft: blameless language, timeline-first structure, contributing factors over single root causes, and action items with owners and dates.

Use this skill when:

- Authoring a postmortem document after a SEV1, SEV2, or learning-worthy lower-severity incident.
- Reviewing a draft postmortem for blameless language, structural completeness, or action-item discipline.
- Writing the "what went well" section without it reading as patronizing.
- Performing a Five Whys analysis in prose.
- Reconstructing a multi-hour incident timeline with timestamps and actors.
- Distinguishing a root cause from a contributing factor in writing.
- Translating raw Slack-channel content into a learning document.

Skip this skill when:

- You are writing the incident response process. Use `incident-response`.
- You are writing live status-page or customer comms during the incident. Use `incident-comms`.
- You are writing a runbook or procedural doc. Use `runbook-craft`.
- You are writing an exec-only briefing with no operational detail. Use `executive-comms`.
- You are doing general technical writing with no incident context. Use `writing-expert` or `technical-writing-craft`.

## Core Concepts

### 1. Blameless framing in prose — the system, not the human

Blameless does not mean accountability-free. The Google SRE Book Chapter 15 is explicit: blameless postmortems "focus on identifying the contributing causes of incidents without indicting any individual or team for bad or inappropriate behavior." Etsy's Debriefing Facilitation Guide goes further: the human is not the root cause. Complex systems must be designed to absorb, detect, and recover from failures because mistakes are inevitable.

In writing, blameless framing is a discrete craft. Three substitutions do most of the work:

- **Names → roles.** "Alice deployed the bad change" → "The release engineer deployed change #4821."
- **Judgments → actions.** "Bob failed to notice the alert" → "The on-call engineer did not see the alert because it was routed to a paused channel."
- **Causal verbs → enabling conditions.** "X caused Y" → "X created conditions under which Y became possible."

The single most damaging phrase in postmortem prose is "should have." It encodes hindsight bias into the document. Replace it with "the system did not surface the information that would have enabled X." This is not euphemism. It is precision: the writer is naming an information gap rather than a moral failure.

When the document refers to specific people by name (which is sometimes unavoidable for timeline accuracy), it names them as actors carrying out their role, never as the cause of the incident.

### 2. Timeline reconstruction in UTC, source-of-truth first

The timeline is the spine of the postmortem. Every later section refers back to it. It must be reconstructible by a future reader who was not present.

Conventions that travel across the Atlassian, FireHydrant, Rootly, and Google SRE postmortem templates:

- **UTC timestamps**, always. Local times in parentheses if your team is mostly in one zone.
- **Source for each event**: which dashboard, which log line, which Slack message, which page. This is what lets a future reader audit the timeline.
- **Actor + action + observable result**, in that order. "14:32 UTC — Alerting system fired `mdb-prod-api-5xx-elevated` (PagerDuty incident #4821). Observable: error rate jumped from 0.2% to 14% in the previous 90 seconds (Datadog dashboard X)."
- **Decision points called out explicitly**, with the information the decider had at the time. This is the antidote to hindsight bias: the timeline shows what was knowable at each moment.

A good timeline table:

| UTC time | Actor | Action | Observable / source |
|---|---|---|---|
| 14:30 | Deploy bot | Released change #4821 to prod | GitHub Actions log |
| 14:32 | Alerting | Fired `api-5xx-elevated` SEV2 | PagerDuty #4821 |
| 14:33 | On-call (release eng) | Acknowledged page | PagerDuty |
| 14:35 | On-call | Joined #inc-0421 | Slack |
| 14:38 | On-call | Initiated rollback of change #4821 | GitHub Actions log |
| 14:46 | Rollback complete | Error rate returned to 0.2% | Datadog |

The timeline is not narrative. It is a forensic record. Save narrative for the analysis section.

### 3. Contributing factors vs root cause — the linguistic and conceptual distinction

The most-cited mistake in postmortem writing is naming a single root cause. Real incidents have a root cause plus contributing factors. The Atlassian postmortem handbook and the Allspaw/Etsy debriefing guide both reject single-cause framing: "contributing factors are not root causes; they're the context that makes the system fragile."

In prose, structure the analysis section as:

- **Triggering event** (the proximate change that made the latent problem visible).
- **Root cause** (the latent defect that the trigger exposed).
- **Contributing factors**, categorized:
  - Technical (missing monitoring, single points of failure, fragile retry logic, missing safeguards).
  - Process (insufficient testing, communication gaps, change-control bypass).
  - Environmental (time pressure, on-call fatigue, recent reorg, vendor incident).

Each contributing factor gets its own paragraph and its own action item. The Five Whys analysis (next concept) feeds this section; it does not replace it.

A litmus test: if removing any single contributing factor would have prevented the incident from happening or from being as severe, it is a contributing factor and it earns an action item.

### 4. Five Whys discipline — not literally five, never one

The Five Whys technique is widely cited and often misapplied. The Hootsuite engineering practice and Konrad Reiche's "Blameless postmortem by design" essay both make the same point: the 5 in Five Whys is not a magic number. It is a heuristic for going up the chain past the proximate cause until you reach an organizational or design-level factor.

The discipline in writing:

- **Start with the most proximate observable failure** ("the API returned 503s") and ask why.
- **Each answer must be a factual observation, not a feeling**. "Why did the API return 503s? Because the database connection pool was exhausted." Not "because the system was under stress."
- **Stop when the next "why" exits engineering control and becomes a business or organizational reality** ("Why is the architecture single-region? Because the cost-benefit analysis from 2023 prioritized single-region cost over multi-region resilience"). That is a legitimate terminus, even if it is only the third why.
- **Allow branching**. A single chain of whys is rare. A real incident has multiple parallel chains that converge on multiple contributing factors. Document each chain.

In prose, the Five Whys section often reads best as nested bullet points or a short chain-of-questions paragraph, not as a numbered list of five whys forced into a template.

### 5. Action items with owners, dates, severity, and traceability

The action-items section is where postmortems most often die. Atlassian's postmortem handbook and the Sherlocks.ai analysis of real postmortems converge: "if an action item doesn't trace back to a root cause or contributing factor, it should be removed." Action items that drift into "we should also do X" wishlist territory dilute the ones that matter.

A defensible action item has:

- **Owner**: a single named person (not a team).
- **Due date**: a real calendar date, not "next quarter."
- **Severity / priority**: P0/P1/P2 calibrated to actual risk reduction.
- **Traceability**: which contributing factor it addresses. ("Addresses contributing factor 3: missing monitoring on connection-pool saturation.")
- **Definition of done**: how the writer will know the action is complete. ("Done when alert fires in staging at 80% pool utilization, validated by drill.")

The Atlassian postmortem guide notes: "Not everything is equally urgent, and teams that mark everything as P1 have no P1s." Discipline in the action-items section is discipline in priority-setting.

Action items also need a follow-up mechanism. The postmortem should commit to a 30-day review where someone confirms each action item is on track, and a 90-day review where each P0/P1 must be closed. This is the difference between a postmortem that improves the system and one that performs improvement.

### 6. "What went well" without performative positivity

The "what went well" section is the most awkwardly written part of most postmortems. Done poorly, it reads as forced sunshine ("the team rallied together!") or as faint praise that undercuts the analysis. Done well, it is a serious accounting of what the response system did right, and it produces action items just like the contributing-factors section does.

A useful frame from the FireHydrant and OneUptime postmortem templates:

- **What worked**: the alert fired correctly, the rollback procedure was executable, the comms cadence was met, a junior engineer caught the symptom first.
- **What we got lucky on**: things that worked but only by accident, that should be made non-accidental. ("The on-call happened to be online; if they had been on a flight, MTTR would have been 4x longer.")
- **What we want to preserve**: practices, automation, or instincts that should be carried into future incidents. These deserve their own preservation action items ("Keep the existing 30-minute heartbeat cadence as a documented standard").

The "what went well" section is not a counterweight to "what went badly." It is a parallel source of action items. Treat it with the same writing discipline.

### 7. Hindsight bias — naming it and writing around it

Hindsight bias is the fact that everything looks obvious in retrospect. The Etsy debriefing guide is explicit on this: organizational learning can only happen when objective data about the event is placed into context with the subjective data of multiple perspectives in the room. The job of the writer is to reconstruct what was knowable at the time, not to grade the responders against what is knowable now.

Linguistic markers of hindsight bias to delete in revision:

- "Clearly..." (it was not clear at the time).
- "Obviously..." (it was not obvious to anyone in the moment).
- "Should have noticed..." (the information environment is what is on trial, not the responder).
- "Failed to recognize..." (replace with "the available signals did not surface X").

When the writer catches themselves typing "should have," they should stop and ask: what would have had to be different about the information available to the responder for the right action to have been the obvious action? That question reframes the sentence as a system-design observation, which is what belongs in a postmortem.

### 8. Customer impact section — quantified, time-bounded, honest

The customer-impact section is read by leadership, by sales, sometimes by legal, and sometimes by regulators. It must be quantified, time-bounded, and honest. It is not the place for the marketing voice.

A complete customer-impact paragraph contains:

- **Duration of impact** with start and end timestamps in UTC.
- **Scope**: which customer segments, which regions, which product surfaces.
- **Quantified impact**: requests failed, customers affected, revenue at risk, SLA credits owed.
- **Workarounds available during the incident** (if any).
- **Residual effects** post-resolution (data to re-ingest, retries needed, etc.).

This section bridges from `incident-comms` (which talked to customers in real time) to the postmortem (which retrospectively accounts for what they experienced). The numbers should match the resolution message verbatim, or the document should explain why they differ.

### 9. The publication and review ritual

A postmortem that lives in a Google Doc that nobody reads is not a postmortem. The Google SRE Book Chapter 15 emphasizes that senior management's active participation in the review and collaboration process is part of what makes the postmortem culture work.

The writing implication: the postmortem must be written for a review meeting, not for a filing cabinet. Specifically:

- **The opening summary must be readable in under 60 seconds** by an executive who will not read the full document. State the impact, the cause, and the top three action items. This is the BLUF (bottom-line-up-front) frame from `writing-expert`.
- **The action items must be the most carefully copy-edited section** because they are the one thing the review meeting will literally walk through line by line.
- **The document should be linkable, versioned, and findable** by a future engineer searching for "what happened the last time the connection pool exhausted?"

### 10. Sensitive incidents — security, privacy, legal review

Some incidents involve data exposure, security breach, or regulatory notification. The writing of those postmortems acquires an additional constraint: every claim must be reviewed by security and legal before publication.

The writing discipline:

- **Two versions, explicitly labeled**: the internal full-detail postmortem and the redacted external version. The redactions must be reviewable by the postmortem authors (so they know what was cut and why).
- **Confirmed claims only** in the external version. "We have evidence that X happened" is acceptable; "we believe Y" is not, in an externally-distributed regulated incident report.
- **Customer notification language coordinated with legal**, never improvised in the postmortem draft.

This is one of the few places where postmortem writing intersects with non-engineering processes. The blameless framing rules still apply, but the legal-defensibility rules also apply, and the writer must accommodate both.

## Templates and Patterns

### Central artifact: a postmortem section (contributing factor + action item)

```markdown
### Contributing factor 3 — Missing monitoring on connection-pool saturation

**Description.** The `mdb-prod-api` connection pool reached 100% utilization
at 14:30 UTC, eight minutes before the first customer-visible error. No alert
fired on pool saturation; only the downstream 5xx alert fired, at 14:32 UTC.

**Why it contributed.** Had a pool-saturation alert existed, the on-call
release engineer would have had 8 minutes of leading indicator before
customer impact. The 5xx alert is a lagging indicator: by the time it fires,
customers have already been affected.

**Five Whys chain.**
- Why did the pool saturate? Change #4821 introduced a query pattern that
  held connections longer than the previous pattern.
- Why was the change deployed without saturation testing? Load-test gating
  is not part of the standard release pipeline; it is opt-in.
- Why is it opt-in? The load-test framework was built in 2024 and never
  promoted to mandatory because of false-positive rates.
- Why are the false-positive rates not addressed? The load-test ownership
  was reorged out of the platform team in 2025 and has no current owner.

**Action item 3.1.** Add a `pool_utilization > 80%` alert to all
`mdb-prod-api-*` services.
- Owner: @mitch.hudson
- Due: 2026-06-15
- Priority: P0
- Done when: alert fires in staging at 80% pool utilization, validated by
  a chaos drill that ramps query latency.

**Action item 3.2.** Assign an owner to the load-test framework and place
load-test gating on the mandatory release-pipeline checklist for any change
that modifies query patterns.
- Owner: @release-eng-lead (to be named by 2026-06-08)
- Due: 2026-07-31
- Priority: P1
- Done when: a change with a regression query pattern is blocked by the
  pipeline gate in staging.
```

### Full postmortem skeleton

```markdown
# Postmortem: <one-line description matching the incident title>

| Field | Value |
| --- | --- |
| Incident ID | INC-0421 |
| Severity | SEV1 |
| Date | 2026-05-22 |
| Duration | 14:32 – 14:46 UTC (14 min) |
| Authors | @mitch.hudson, @alice |
| Reviewers | @release-eng-lead, @vp-eng |
| Status | Draft / In review / Published |
| Action-item review date | 2026-06-22 (30-day) |

## Executive summary (read this first — 60 seconds)
A change deployed at 14:30 UTC introduced a query pattern that saturated
the API connection pool. From 14:32 to 14:46 UTC, approximately 14% of API
requests in us-east-1 returned 503 errors. Rollback restored service at
14:46 UTC. Three contributing factors have been identified, five action
items have been opened, two are P0.

## Customer impact
[quantified, time-bounded, with workarounds and residuals]

## Timeline (UTC)
[table; actor + action + observable + source]

## What happened — narrative
[2-4 paragraphs of factual reconstruction, no judgments]

## Root cause and contributing factors
- Triggering event: [...]
- Root cause: [...]
- Contributing factor 1: [...]
- Contributing factor 2: [...]
- Contributing factor 3: [...]

## What went well
- [what worked]
- [what we got lucky on — needs action item]
- [what we want to preserve]

## Action items
| ID | Description | Owner | Due | Priority | Addresses |
|---|---|---|---|---|---|
| 1.1 | ... | @x | 2026-06-15 | P0 | CF1 |
| ... | ... | ... | ... | ... | ... |

## Follow-up
- 30-day review: 2026-06-22 (status check on all action items)
- 90-day review: 2026-08-22 (all P0/P1 must be closed)
- Linked tickets: [...]

## Appendix
- Linked runbooks: [...]
- Linked dashboards: [...]
- Slack channel archive: [...]
```

### Blameless rewrite cheat sheet

| Before (blamey) | After (blameless) |
| --- | --- |
| Alice deployed a bad change | Change #4821 was deployed at 14:30 UTC |
| Bob failed to notice the alert | The pool-saturation alert did not exist; only the lagging 5xx alert fired |
| The team should have caught this in review | The review checklist did not include load-test sign-off |
| Bad judgment by the release engineer | The release pipeline did not gate the change on pool-saturation testing |
| Communication broke down | The incident channel was created at 14:35 UTC, three minutes after the page; this is within the documented SLA |

## Anti-Patterns

- **Single-root-cause syndrome**: "the root cause was X" with no contributing factors. Real incidents have layered causes; pretending otherwise hides the real lessons.
- **Wishlist action items**: "we should also rewrite the deployment system." Not traceable to a contributing factor → cut.
- **Hindsight prose**: "obviously the team should have noticed..." — never appears in a competent postmortem.
- **Performative positivity**: "what went well: the team showed great resilience!" — not actionable, not learning-bearing, cut.
- **Marketing voice in the customer-impact section**: "Customers may have experienced a brief inconvenience" when the data says 14% of requests failed for 14 minutes. Honesty here protects future credibility.
- **The pile-on**: long, detailed enumeration of one team's mistakes. Almost always a sign the writer is grinding an axe.
- **The publish-and-forget**: no review date, no owner for action items, no follow-up commit. The document becomes culture-theater.

## Decision Heuristics

- **When to write a postmortem**: any SEV1 or SEV2 always; SEV3 if it is novel, repeatable, or surfaced a learning the team did not already have. Skip for SEV4 unless someone explicitly requests it.
- **When to share externally**: SEV1 customer-facing incidents always get an external version. SEV2 case-by-case. SEV3 and below stay internal. The external version is a redaction of the internal version, never a separate draft.
- **When to escalate the action-item review**: a P0 action item that has not closed by 90 days. Escalates to eng leadership for re-prioritization.
- **When to call this skill from `incident-response`**: the moment the incident is resolved and the resolution message has been sent, route the writeup to `postmortem-writing`. The two skills handle different phases of the same lifecycle.
- **When to involve security/legal review**: any incident with customer-data exposure, regulatory implications (SOC2, HIPAA, GDPR notification thresholds), or potential litigation exposure. Two-version pattern.
- **When to hand off to `executive-comms`**: when leadership wants a 1-page summary stripped of the operational detail. Produce a separate exec doc; do not let it dilute the engineering postmortem.
- **When the Five Whys is the wrong tool**: when the incident is genuinely a one-off (e.g., a vendor outage outside your control). In that case, the document is shorter, focuses on detection/response/recovery, and explicitly notes the cause is exogenous.

## References

- [Google SRE Book — Chapter 15: Postmortem Culture: Learning from Failure](https://sre.google/sre-book/postmortem-culture/) — the canonical reference for blameless postmortem culture, senior-management participation, template structure.
- [Google SRE Book — Example Postmortem](https://sre.google/sre-book/example-postmortem/) — a worked example, with timeline, root cause, contributing factors, and action items.
- [Google SRE Workbook — Postmortem Culture](https://sre.google/workbook/postmortem-culture/) — the practitioner's companion, including coordination across teams.
- [Google SRE Workbook — Postmortem Analysis](https://sre.google/workbook/postmortem-analysis/) — trend analysis across postmortems, working-group coordination.
- [Etsy — Debriefing Facilitation Guide (Allspaw, Evans, Schauenberg)](https://extfiles.etsy.com/DebriefingFacilitationGuide.pdf) — the seminal guide on blameless debrief facilitation; PDF, open-sourced by Etsy.
- [Etsy Code as Craft — Debriefing Facilitation Guide for Blameless Postmortems](https://www.etsy.com/codeascraft/debriefing-facilitation-guide/) — blog announcement and context for the guide.
- [Atlassian — Postmortems: Enhance Incident Management Processes](https://www.atlassian.com/incident-management/handbook/postmortems) — Atlassian's postmortem section of their incident management handbook.
- [Atlassian — The power of 5 Whys: analysis and defense](https://www.atlassian.com/incident-management/postmortem/5-whys) — the Five Whys technique applied to postmortems.
- [Atlassian — The importance of an incident postmortem process](https://www.atlassian.com/incident-management/postmortem) — overview of postmortem mechanics and contributing-factor framing.
- [Konrad Reiche — Blameless postmortem by design: In praise of the Five Whys](https://konradreiche.com/blog/blameless-postmortem-by-design-in-praise-of-the-five-whys/) — disciplined application of Five Whys in real engineering postmortems.
- [Hootsuite Engineering — 5 Whys: how we conduct blameless post-mortems](https://medium.com/hootsuite-engineering/5-whys-how-we-conduct-blameless-post-mortems-after-something-goes-wrong-a47687baeacc) — practical Five Whys execution in industry.
- [FireHydrant — The Ultimate, Incident Retrospective (Postmortem) Template](https://firehydrant.com/blog/incident-retrospective-postmortem-template/) — concrete template with "what went well" / "what went badly" / "where we got lucky."
- [Code for America / Lou Moore — From accident to investment: How to run better blameless postmortems](https://medium.com/code-for-america/from-accident-to-investment-how-to-run-better-blameless-postmortems-43553d421838) — facilitation craft and the framing of postmortems as investments.
- [Google Cloud — Google's tips on how to hold a fearless shared postmortem](https://cloud.google.com/blog/products/gcp/fearless-shared-postmortems-cre-life-lessons) — cross-organization postmortem sharing, customer-facing postmortems.
- [O'Reilly — Site Reliability Engineering, Chapter 15 (book version)](https://www.oreilly.com/library/view/site-reliability-engineering/9781491929117/ch15.html) — the book edition of the postmortem chapter.

---

## Appendix: The hourglass structure for postmortems

Most postmortems pick the wrong structural shape and pay for it. The two most common defaults are wrong in opposite directions: a pure inverted pyramid loses the timeline, and a pure chronological narrative buries the verdict. The right structure for almost every incident postmortem is the *hourglass* — a hybrid named by Roy Peter Clark (Poynter, 1983) that combines both.

**The three parts of an hourglass postmortem.**

1. **The top (inverted pyramid summary, 4–6 paragraphs).** The verdict first. What broke, when, who was affected, severity, current status, owner, time to detection, time to resolution, root cause class. A reader who reads only the top must come away knowing everything they need to make a decision.
2. **The turn (one sentence).** A transitional line that signals the document is now switching modes: "Here is how the incident unfolded, in chronological order." That sentence is load-bearing — without it, readers feel jolted from summary to timeline.
3. **The bottom (chronological narrative).** The timeline, starting before the trigger event (last known good state), through detection, escalation, mitigation, full resolution, and follow-up actions. This is where the *learning* lives — the top tells you what happened; the bottom teaches you why.

**Why the hourglass beats a pure inverted pyramid for postmortems.** An inverted pyramid forces every fact into "is this more or less important than the previous fact?" That works for news, but postmortems teach the *sequence* — what was true at 14:02, what changed at 14:07, what the on-call saw at 14:11. The hourglass lets execs read only the top (for the verdict) while engineers read the bottom (for the learning) — both audiences served by one document.

**Why the hourglass beats a pure narrative.** A pure timeline assumes the reader has 20 minutes. Most executives and adjacent teams have 90 seconds. They'll bail before the lesson lands. The top is what makes the postmortem skim-resilient.

**Worked example — postmortem opening (the top + the turn).**

> **Incident:** Checkout cache eviction storm.
> **Date:** 2026-05-08, 14:02–14:24 UTC (22 minutes).
> **Severity:** SEV-1 (full checkout outage).
> **Customers affected:** ~340,000 active sessions; 12% of checkout attempts during the window failed.
> **Root cause class:** Misconfigured TTL on cache cluster after deploy.
> **Status:** Resolved; rollback complete. Permanent fix shipped 2026-05-12.
>
> *Here is how the incident unfolded, in order.*

That's the top and the turn. The bottom is the timeline. The document is structurally clear before the reader has scrolled.

**Anti-pattern — interleaving timeline and summary.** Writers often try to be helpful by inserting summary lines into the middle of the timeline ("By 14:07, the situation was severe…"). This breaks both halves: the summary loses its skim-resilience because it's diffused, and the timeline loses its chronological clarity because narrative tension is repeatedly interrupted. Keep them separate. The top summarizes; the bottom narrates.

**Diagnostic — the 90-second test.** Hand the draft to someone unfamiliar with the incident. Ask them to read for 90 seconds and then summarize. If they can state what broke, when, who was affected, and current status, the top is doing its job. If they can't, restructure — the verdict isn't reaching them.

**When to break it.** Two cases warrant a different structure:

- *Near-miss postmortems with no customer impact.* The timeline matters less than the systemic lesson. Lead with the lesson; the timeline can live in an appendix.
- *Postmortems for incidents with no clear timeline* (e.g., a long-running performance regression discovered retroactively). There is no chronology to narrate. Use a pure inverted pyramid.

**References.**

- Roy Peter Clark, Poynter Institute — [The hourglass: serving the news, serving the reader](https://www.poynter.org/reporting-editing/2003/the-hourglass-serving-the-news-serving-the-reader/), the original 1983 framework.
- Pressbooks / CWI — [Newswriting Structures: The Inverted Pyramid and Beyond](https://cwi.pressbooks.pub/introductiontojournalismandnewswriting/chapter/chapter-5-newswriting-structures-the-inverted-pyramid-and-beyond/), comparative treatment of inverted pyramid, hourglass, and narrative.
- See also: `writing-expert` skill — *Hourglass structure — the third way between inverted pyramid and narrative* for the general-purpose treatment.
