<!-- hub-reference-banner -->
> **Reference file — part of the `writing-expert` hub.** Formerly the standalone `realtime-writing-under-pressure` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: realtime-writing-under-pressure
description: "Writing craft for real-time, time-pressured channels: Slack-thread updates, breaking-news comms, 60-second answers, live blogs, status-page posts, holding statements, time-stamped update threads, twitter-thread structure under embargo lift, 'I don't know yet' answers, the two-pass-minimum-even-in-5-minutes discipline, hot-take vs cold-take, accuracy-vs-speed tradeoffs, 'say less, link more' rule, the first-draft-is-the-publish pattern, and the journalistic verification-before-publication standard adapted for engineering, support, and ops. Drawn from AP Stylebook breaking-news guidance, Reuters real-time newsroom standards, Slack/incident.io/PagerDuty incident-thread patterns, and live-blog conventions. TRIGGER: user is writing under time pressure and asks for a holding statement, Slack incident-thread update, breaking-news note, customer-facing status update during an outage, time-stamped live thread, 'we don't know yet' answer, two-pass quick edit, or any 'I have 5 minutes' draft. SKIP: structured incident communications with a formal template (use incident-comms); structured executive briefings (use executive-comms); post-incident write-ups, post-mortems, RCAs (use postmortem-writing); calm-time reports, QBRs, summaries (use writing-expert); legal-adjacent breach-notification drafts (use legal-adjacent-writing)."
origin: local
version: "1.0.0"
updated: "2026-05-29"
keywords:
  - real-time writing
  - Slack thread
  - breaking news
  - holding statement
  - status update
  - time-stamped update
  - live blog
  - two-pass edit
  - accuracy vs speed
  - hot take
  - cold take
  - first draft is the publish
  - I don't know yet
  - say less link more
  - incident channel writing
  - status page writing
  - 60-second answer
tags:
  - writing
  - realtime-writing
  - slack
  - status-page
  - breaking-news
  - live-blog
whenToUse:
  - User is writing a Slack thread update during an incident or live event
  - User must publish a holding statement before facts are confirmed
  - User has 5 minutes or less and asks for a quick draft
  - User is running a live blog, live-tweet thread, or time-stamped update channel
  - User must answer a stakeholder while the situation is still unfolding
  - User is posting to a status page mid-incident
  - User must say "we don't know yet" without sounding evasive
  - User is writing a breaking-news note under embargo lift, press release, or earnings-call moment
  - User must publish an update knowing it will be read in 30 seconds, possibly on mobile
  - User is on hour 3 of an incident and needs a fresh stakeholder update
whenNotToUse:
  - User has a structured incident-communications template (use incident-comms)
  - User is writing a structured executive briefing or board memo (use executive-comms)
  - User is writing a post-mortem, RCA, or after-action report (use postmortem-writing)
  - User is writing a planned QBR, status report, or routine summary (use writing-expert)
  - User is drafting a regulator-facing breach notification under GDPR/SEC rules (use legal-adjacent-writing)
related_skills:
  - incident-comms
  - executive-comms
  - postmortem-writing
  - writing-expert
  - support-ticket-writing
  - storytelling-and-narrative
  - public-speaking-and-presentations
---

# Real-time Writing Under Pressure

Reference for writing under a clock — Slack threads during incidents, holding statements before facts are confirmed, status-page posts at hour 3 of an outage, live-blog updates, "the customer is on a call right now" notes, "I have five minutes" drafts. The discipline is journalistic: accuracy first, speed second, transparency about uncertainty, and a structure that survives being read on a phone with 30 seconds of attention.

This is the partner skill to `incident-comms` (which provides structured templates) and `executive-comms` (which provides structured briefings). This skill kicks in when there is no time for either.

## When to use this skill

- Writing a Slack thread update during an active incident or live event
- Writing a holding statement before facts are confirmed
- Composing a status-page post mid-outage
- Running a live blog or time-stamped update thread
- Drafting a 60-second answer to a stakeholder while a situation is unfolding
- Writing a "we don't know yet" message that doesn't sound evasive
- Posting an update under embargo lift, breaking-news moment, or earnings event

## When NOT to use this skill

- You have a formal incident-comms template available → `incident-comms`
- You have time for a structured executive briefing → `executive-comms`
- You're writing post-incident — RCA, post-mortem, after-action → `postmortem-writing`
- You have time and a calm room → `writing-expert`
- You're drafting a regulator notification with legal-adjacent risk → `legal-adjacent-writing`

## The five rules of real-time writing

These are the rules that survive all situations. When you have one minute, follow these. When you have five, follow these and also run the two-pass edit.

1. **Accuracy beats speed.** Reuters' standard. The Washington Post's standard. Reasonable people accept a 90-second delay; nobody forgives a wrong update sent in 30. If the fact isn't confirmed, say "investigating," not the unconfirmed fact.
2. **Timestamp every update.** "16:42 UTC — Investigating elevated error rate." Without timestamps, a reader cannot tell what is current. Use one timezone (UTC for engineering, your headquarters timezone for customer-facing).
3. **Lead with status, not narrative.** First word should be one of: Investigating, Identified, Monitoring, Resolved, Update, or a status verb. Not "We have observed…" or "It appears that…"
4. **Say what is known, what is unknown, and what is next.** Three sentences. "Known: error rate is up 4x since 16:30. Unknown: root cause is still under investigation. Next: update in 15 minutes or sooner."
5. **One thought per message.** A Slack thread is read line-by-line on mobile. Walls of text get scrolled past. If you have three updates, post three messages with three timestamps.

## Core concepts

### 1. The two-pass-minimum-even-in-five-minutes rule

Speed-vs-accuracy is a real tradeoff. Two passes — even tiny ones — catch the most embarrassing errors:

- **Pass 1: Write fast.** Type the message. Don't self-edit yet.
- **Pass 2: Read it back.** Take 15 seconds. Check three things: (a) Are the facts that you stated actually confirmed? (b) Is there a stale time, status, or number from the previous update that should be updated? (c) Is the most important sentence the first sentence?

Five minutes is enough for this. One minute is not. If you have one minute, write the holding statement; do not invent facts to fill space.

The discipline is journalistic. Reuters: "accuracy, as well as balance, always takes precedence over speed." The Washington Post runs multi-step verification even at speed. Even an internal Slack thread benefits from the same instinct — every wrong update creates a remediation loop that costs more time than the verification would have.

### 2. Hot take vs cold take

A **hot take** is a fast reaction to live developments. It is appropriate when:
- The audience knows you are reacting in real time and accepts uncertainty
- The cost of being wrong is recoverable (an internal Slack message can be corrected; an SEC filing cannot)
- The alternative — silence — is worse than a flagged-as-tentative take

A **cold take** is the considered version that comes later — with more facts, after rereading, possibly after one good night's sleep. It is appropriate when:
- The audience expects accuracy (customers, regulators, the press)
- The cost of being wrong is high (legal exposure, brand damage, contract risk)
- You have time to wait for facts

The discipline is to know which one you're writing and to label it. "Quick reaction, not a final assessment" is a fair caveat. So is "I'll update once we have confirmation."

### 3. The "I don't know yet" sentence

The biggest temptation under pressure is to fill silence with speculation. Don't. The fix is a sentence pattern that says "I don't know" without sounding evasive:

- "We're investigating and don't have a confirmed cause yet."
- "We don't have that data yet. We will have it by [time]."
- "I want to give you an accurate answer; let me confirm and come back in [N] minutes."
- "What I can tell you now: [confirmed facts]. What we're still confirming: [list]."
- "Two scenarios are possible based on early signals. I won't speculate until we narrow it."

These all communicate care, competence, and that an update is coming. They are far better than "It looks like a database issue, probably" — which becomes a quoted commitment.

### 4. The holding statement template

A holding statement is the first communication after a situation is detected, before facts are confirmed. Its job is to acknowledge, set expectations for the next update, and avoid commitments.

Structure:

1. **Acknowledgment.** "We are aware of [observable symptom]."
2. **Status.** "Investigating." Not "isolated incident." Not "minor issue." Not "we're sure it's fine."
3. **Action being taken.** "Our team is actively working on identifying the cause."
4. **Next update commitment.** "We'll provide an update by [specific time]."
5. **(Optional) Where to follow.** "Watch this thread / our status page at [URL] for updates."

A holding statement does NOT include:
- Root-cause speculation
- Eta for resolution (unless you genuinely know)
- Words like "minor," "isolated," "small subset" before the data supports them
- Apologies that imply more than is known ("we apologize for the major outage" — when you don't yet know it's major)

### 5. The "say less, link more" rule

In real-time channels, the message you can fit on a phone screen is the message that gets read. Everything longer than four lines gets scrolled. Patterns:

- **Link to the dashboard, status page, or runbook** instead of describing it.
- **Link to the previous detailed update** instead of restating it.
- **Link to the incident channel** for stakeholders watching from outside.
- **One detailed update per hour, plus brief heartbeat updates** in between.

The mistake to avoid: posting a comprehensive message that takes 90 seconds to scan, when a one-line "status quo, next update at 17:00" would have served better.

### 6. Time-stamped update thread structure

A live-blog or incident-thread update sequence follows a pattern. Each post is self-contained — a reader joining at 16:45 should understand without scrolling to 16:00.

```
16:42 UTC — Investigating
We are investigating reports of elevated error rates on the API.

16:55 UTC — Identified
We have identified the cause as a degraded connection pool on the
authentication service. We are rolling back the most recent
deployment. ETA to restoration: 17:15 UTC. Status page: [link]

17:18 UTC — Monitoring
The rollback is complete and error rates have returned to baseline.
We are monitoring before declaring fully resolved.

17:42 UTC — Resolved
The incident is resolved. Total impact: 16:30–17:18 UTC. We will
post a full post-incident review within 5 business days.
```

Each post starts with a timestamp and a status verb. Each post can stand alone. Each post commits to the next update or marks resolution.

### 7. Verbs of status

A common-vocabulary list for status channels — used by Atlassian Statuspage, PagerDuty, incident.io, and most incident-management systems. Pick from this list, not from your imagination:

- **Investigating** — we see a symptom, cause is unknown
- **Identified** — we know what is wrong
- **Monitoring** — we believe a fix is applied, watching to confirm
- **Resolved** — we are confident the incident is over
- **Update** — neutral status carrier for additional information

For maintenance:
- **Scheduled** — announced in advance
- **In Progress** — happening now
- **Completed** — done

Mixing these (using "ongoing" or "still happening" or "we think it's better") confuses status-page readers who scan for verbs.

### 8. The first-draft-is-the-publish pattern

In real-time channels there is no editor between you and the audience. The first draft is also the final draft. Three implications:

1. **No revisions to land jokes or improve flow.** Get the facts right first; humor is risky under pressure.
2. **No long sentences with embedded clauses.** Short, declarative, scannable.
3. **No "Here's what we think happened" prose.** Replace with "Confirmed: X. Investigating: Y."

The exception is the cold-take follow-up: when you have an hour to write the longer-form post-incident summary, you may revise freely. But the live-thread updates remain as they were posted, timestamped, even if later updates supersede them.

### 9. Mobile-first formatting

Most real-time updates are read on phones. Formatting moves that survive mobile:

- **First sentence carries the message.** Not the third.
- **Bold the most important word or number.** Skim-readers see bold first.
- **Bulleted lists beat paragraph prose** when listing facts.
- **Avoid wide tables.** They wrap on mobile and become unreadable.
- **Hyperlinks live on their own line** if they're long; otherwise inline.
- **No more than 5–6 lines per message.** Slack collapses anything longer.

### 10. The "next update by" commitment

Every real-time update should end with one of three things:

1. **A specific next-update time.** "Next update by 17:00 UTC."
2. **A trigger condition.** "Next update when we confirm the rollback is clean."
3. **A resolution marker.** "Resolved. No further updates needed."

The reader's anxiety comes from not knowing when to check again. The next-update commitment relieves it. Missing the commitment (no update by the time you promised) is worse than promising in the first place — so promise conservatively. "Update by 17:00" is better than "update in 5 minutes" if you might slip.

### 11. Stakeholder vs operator channels

Most real-time situations have two audiences and should have two channels:

- **Operator channel** (internal Slack #inc-YYYYMMDD-X, war-room voice call): technical detail, debugging, log excerpts, hypotheses. Noisy. Limited participation.
- **Stakeholder channel** (Slack #incident-updates, status page, customer comms): summary, status, impact, ETA. Quiet. Broad audience.

The writing differs:

| | Operator channel | Stakeholder channel |
|---|---|---|
| Tone | Technical, terse | Plain, calm |
| Frequency | High; whenever something changes | Low; planned cadence |
| Jargon | Welcome | Avoided |
| Hypotheses | Welcome | Avoided until confirmed |
| Log excerpts | Welcome | Linked only |
| Length | Whatever the situation needs | Five lines or fewer |

The skill of real-time stakeholder writing is filtering operator-channel detail down to the 10% that stakeholders need.

### 12. "Two-by-two-by-two" — known/unknown × confirmed/suspected

A useful matrix for sorting what to put in an update:

| | Confirmed | Suspected |
|---|---|---|
| **Known** | Always say this | Maybe, with caveat |
| **Unknown** | "Still investigating Y" | Never speculate |

The publishable cells are top-left always, top-right with explicit hedging, and bottom-left as transparency. Bottom-right is the speculation trap — don't write what you suspect about a thing you don't yet know.

### 13. Reuters/AP discipline adapted

Adapted from journalistic real-time standards:

- **Reuters:** "Accuracy, as well as balance, always takes precedence over speed."
- **AP Stylebook (breaking news):** Prioritize clarity and accuracy over strict style. Get it right; style can be cleaned later.
- **NPR:** Verification is not optional even under pressure. When information is uncertain, the responsible approach is to say so.
- **The Washington Post:** Multi-step verification even at speed.

The engineering and support adaptation: when in doubt, verify before publishing. When you can't verify in the time you have, publish the uncertainty explicitly. "We do not yet know X" is a complete and publishable sentence.

## Templates

### Template 1: Holding statement (first 5 minutes of an incident)

```
16:42 UTC — Investigating
We are aware of [symptom] on [service]. Our team is investigating.

Impact: [what users are seeing, in plain terms — be conservative].
Status page: [URL]
Next update by 17:00 UTC or sooner if conditions change.
```

### Template 2: Brief status-page update (mid-incident, no new facts)

```
17:00 UTC — Update
We are continuing to investigate. We have ruled out [thing] and are
examining [other thing]. Impact remains [same as before / changed to X].
Next update by 17:30 UTC.
```

### Template 3: Identified-cause update

```
17:18 UTC — Identified
We have identified the cause as [brief, plain-English description].
Mitigation in progress: [what is happening now]. Expected restoration
window: [time range].
Next update by 17:45 UTC.
```

### Template 4: Resolved update

```
18:02 UTC — Resolved
The incident is resolved as of 17:54 UTC.

Total impact window: 16:30–17:54 UTC (84 minutes).
Affected: [scope — services, customers, regions].
Cause: [brief, accurate, non-speculative].
Next steps: A full post-incident review will be published within
[N business days].
```

### Template 5: Customer-facing breaking note ("we know about it")

```
We know some customers are seeing [symptom] right now. Our team is
on it. We will update you here within the next 30 minutes. You can
also watch [status page URL] for live updates.

Thank you for your patience.
```

### Template 6: Slack stakeholder update at hour 3 of incident

```
:rotating_light: Hour 3 update — 19:30 UTC

*Status:* Monitoring (rollback applied at 18:55, error rates normal)
*Impact window:* 16:30–18:55 UTC
*Customer comms posted:* Status page + email to affected accounts
*What we're watching:* P95 latency, replication lag (both green)
*Decision needed by stakeholders:* None right now
*Next update:* 20:30 UTC unless conditions change
```

### Template 7: "We don't know yet" answer to a stakeholder DM

```
Honest answer: we don't have a confirmed cause yet. Here's where we are.

What is confirmed:
- [Fact 1]
- [Fact 2]

What we're investigating:
- [Open question 1]
- [Open question 2]

I'll have a more concrete answer by [time]. If anything changes
before then I'll ping you.
```

### Template 8: Live-tweet / live-thread structure (embargo lift, launch, breaking news)

```
1/ [Lead with the news in one sentence. Most important fact. No hype.]

2/ [The one-line context — why this matters. Still no hype.]

3/ [The first detail or quote. Concrete.]

4/ [The second detail.]

5/ [Caveats, limitations, what we don't know yet.]

6/ [Source, link, where to read the full story.]

/end
```

Each tweet stands alone. Each is timestamped by the platform. Each adds one fact, not five.

### Template 9: Quick two-pass edit checklist

```
PASS 1 (write fast)
[ ] Draft the message in the channel's draft area (don't post yet)

PASS 2 (read back, 15 seconds)
[ ] Are all stated facts confirmed?
[ ] Is the timestamp current?
[ ] Is the status verb correct (Investigating / Identified / Monitoring / Resolved)?
[ ] Does the first sentence carry the message?
[ ] Is there a next-update commitment?
[ ] Is anything from the previous update now stale?
[ ] Reading on mobile, does it fit on one screen?

POST.
```

### Template 10: "I'll find out and come back" answer (when you're cornered)

```
I don't want to give you an answer I'm not sure about. Let me
[check the dashboard / talk to the on-call engineer / verify with the
team] and come back to you in [N] minutes with a confirmed answer.
```

## Anti-patterns

1. **Posting without a timestamp.** Without a time, a reader can't tell if the message is from 5 minutes ago or 50.
2. **Speculating to fill silence.** "It's probably the database." Becomes the quoted commitment. Wait for facts.
3. **"Minor issue" / "isolated" / "small subset" before data supports those words.** These get screen-capped and shared.
4. **No next-update commitment.** Reader doesn't know when to come back; tickets and pings increase.
5. **Burying status under prose.** "We have observed elevated error rates that began approximately at 16:30 UTC and have been working to identify the root cause…" — start with the verb. "Investigating: elevated error rates since 16:30 UTC."
6. **Posting the same comprehensive update at every cadence.** Stakeholder fatigue. Use brief heartbeat updates ("status quo, next update at 17:30") between substantive updates.
7. **Mixing operator and stakeholder voice in the same channel.** Either split channels or pick the voice deliberately.
8. **Apologizing prematurely.** "We deeply apologize for the major outage" — when you don't yet know it's major or how long. Apologize specifically once you know what you're apologizing for.
9. **Posting a Slack message that's a wall of text.** Use bullets, short paragraphs, and links to detail.
10. **Promising an aggressive next-update time you can't hit.** "Update in 5 minutes" and then no update for 25. Promise conservatively.
11. **Editing the message after posting** (in channels where edits are silent). Reply with a correction; don't pretend the original was always right.
12. **Treating an internal Slack message as low-stakes.** Slack messages get screenshotted and forwarded. Write every public-ish message as if it could appear in a regulator's exhibit.
13. **"Working on it" with no detail.** Without scope, the reader assumes the worst. "Working on the rollback, ETA 15 minutes" is far better.
14. **Not separating cold-take time from hot-take time.** The post-incident summary written in the same 5-minute window as a live update will repeat live-update mistakes. Write the cold version later.
15. **Saying "all clear" before the monitoring window has confirmed it.** A premature resolved-marker forces a retraction. Stay in Monitoring until you're sure.

## Decision heuristics

| Situation | Choice |
|---|---|
| You have 1 minute, no facts confirmed | Holding statement. Acknowledge + investigate + next-update time. |
| You have 5 minutes | Holding statement + the two-pass edit. |
| You have 30 minutes | Run a structured update via `incident-comms` if a template exists. |
| You're asked "what happened?" 20 minutes in | Time-stamped status: confirmed facts, open questions, ETA. |
| You're tempted to speculate | Don't. Use one of the "I don't know yet" patterns. |
| You're tempted to call it "minor" | Don't, unless you have confirmed scope. Wait for data. |
| Stakeholder asks for an ETA you don't have | "I don't have a confirmed ETA. I'll update you when I do, by [time]." |
| You posted something wrong | Post a correction reply ASAP with the corrected fact and a timestamp. Don't quietly edit. |
| Your next-update time is approaching and you have nothing new | Post a heartbeat anyway: "No change since [time]. Next update at [time]." |
| Reading the update on mobile feels long | Cut it. 4–5 lines max. |
| You're at hour 3 and fatigued | Slow down. Two-pass everything. Fatigue is when mistakes happen. |
| You want to add humor or personality | Wait for the cold-take version. |
| The incident is over | Resolved marker, impact window, cause, post-incident plan. Then write the cold-take version separately. |

## Cross-skill notes

- **Use `incident-comms` when you have a template.** That skill is structured incident communications (the form, the cadence, the roles). This skill is what to do when there is no template and the clock is running.
- **Use `executive-comms` for the executive briefing once the situation stabilizes.** Real-time writing is for the live thread; the executive briefing comes after.
- **Use `postmortem-writing` for the cold-take after the incident.** Same incident, different document, different time pressure.
- **Use `writing-expert` for the post-incident summary, customer-facing recap, or QBR mention.** Calm-time prose.
- **Use `support-ticket-writing` for the per-customer follow-up** once the incident is resolved and individual customers ask for specifics.
- **Use `storytelling-and-narrative`** for the eventual retrospective ("how the team handled it" story). Not during.
- **Use `legal-adjacent-writing`** for any regulator-facing breach notification triggered by the incident. The 72-hour clock and the 4-business-day SEC clock are separate from the real-time stakeholder thread.

## References

1. Poynter, *How breaking news is changing the AP Stylebook in real time*: https://www.poynter.org/reporting-editing/2025/ap-stylebook-breaking-news-updates/
2. Reuters, *Handbook of Journalism* (real-time newsroom standards and accuracy-over-speed): https://mediakar.wordpress.com/wp-content/uploads/2012/10/handbook-of-journalism-reuters.pdf
3. NPR, *Accuracy*: https://www.npr.org/about-npr/688139552/accuracy
4. Slack, *Stuff happens: using Slack for incident management*: https://slack.com/blog/transformation/incident-management-slack
5. incident.io, *Incident communication best practices: Keep stakeholders informed*: https://incident.io/blog/incident-communication-best-practices
6. PagerDuty, *Why Dedicated Incident Channels are the Modern Standard for Slack-Based Incident Response*: https://www.pagerduty.com/blog/insights/why-dedicated-incident-channels-are-the-modern-standard-for-slack-based-incident-response/
7. Instatus, *Comprehensive Guide to Slack Incident Management*: https://instatus.com/blog/slack-incident-management
8. Digital Content Next, *Speed vs. accuracy: Journalism's ethical balancing act*: https://digitalcontentnext.org/blog/2026/03/16/speed-vs-accuracy-journalisms-ethical-balancing-act/
9. Fiveable, *Balancing speed and accuracy in breaking news*: https://library.fiveable.me/newsroom/unit-5/balancing-speed-accuracy-breaking-news/study-guide/HTLOthO831dCRD6f
10. Atlassian Statuspage status-verb vocabulary (Investigating / Identified / Monitoring / Resolved) — industry-standard convention.
