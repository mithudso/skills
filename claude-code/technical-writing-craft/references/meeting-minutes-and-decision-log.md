<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `meeting-minutes-and-decision-log` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: meeting-minutes-and-decision-log
description: Capture craft for meeting notes, action items, decision logs, attendee tracking, agenda-vs-minutes distinction, ADR (Architecture Decision Record) form for technical decisions, "what was decided / why / who-owns / when" template, async vs sync minutes, and verbatim-quote vs paraphrase judgment. TRIGGER when the user says "take minutes", "write meeting notes", "decision log", "record this decision", "ADR", "architecture decision record", "action items from the meeting", "what did we decide", "document this call", "post-meeting summary", "log this decision", "decision record template", "who owns this", "async standup notes", "board minutes", "retrospective notes", or asks how to capture a meeting in writing. SKIP for: pushing action items into Jira/Confluence (use atlassian:capture-tasks-from-meeting-notes); broader product/system design narratives (use software-architect or rfc-and-design-docs); executive-facing rollups and stakeholder briefings (use executive-comms); multi-step implementation plans (use code-plan-writing or agent-plan-writing); generic prose rewriting (use writing-expert); status reports that are not anchored to a specific meeting or decision event.
---

# meeting-minutes-and-decision-log

## Overview

This skill is for the moment a meeting just ended (or is happening async) and a written artifact has to come out the other side. The artifact has one job: let someone who was not there reconstruct what was decided, why, who owns the next step, and when it is due — without having to ask a participant.

There are three distinct artifacts and they should not be conflated:

1. **Agenda** — written *before* the meeting. Topics, owners, time boxes, desired outcomes. Forward-looking.
2. **Minutes** — written *during/after*. Outcomes, decisions, action items, dissent on the record. Backward-looking and durable.
3. **Decision Log / ADR** — written *when a decision is significant enough to outlive the meeting*. One decision per record, indexed, immutable once accepted. Cross-meeting.

The most common failure mode is producing a transcript instead of minutes — long discussion paragraphs that bury the four things that actually matter: what was decided, why, who owns the follow-up, and by when.

## Core Concepts

### 1. The four load-bearing fields

Every captured decision must answer:

- **What was decided** — a single declarative sentence, not a description of the debate.
- **Why** — the one or two reasons that tipped it. Trade-offs accepted. Alternatives rejected.
- **Who owns it** — exactly one named person per action item. Co-ownership is no ownership.
- **When** — a specific date, not "soon" or "next sprint". If it has no date, it is a wish, not an action item.

If any of the four is missing, the entry is incomplete and should be flagged before the minutes go out.

### 2. Agenda vs minutes — separate documents, different tense

Agenda items are intent ("discuss pricing tier rollout"). Minutes items are outcome ("decided to defer pricing tier rollout to Q3; revisit June 15 sync"). Do not paste the agenda into the minutes file and try to "fill it in" — the cognitive frame is different. Start minutes from a blank template with explicit "Decisions" and "Action Items" sections.

### 3. Sync vs async minutes — adjust the verbosity

**Sync minutes** (live meeting, scribed in real time): Capture only what is needed to reconstruct outcomes. Skip the back-and-forth. Aim for ~10% the length of the conversation. Bullet form.

**Async minutes** (Slack thread, Loom comment chain, doc-comment review): The thread *is* the discussion. The minutes job is to write the **synthesis at the top** — a 5-line summary of what was decided after N days of comments — so the next reader does not have to scroll through 80 messages. Pin the synthesis.

### 4. Verbatim quote vs paraphrase — when to use which

**Paraphrase by default.** Minutes are not a court transcript. Compress.

**Use a verbatim quote** in three cases only:

- **On-the-record dissent**: "Sarah objected: 'We are accepting a 6-month support liability we have not budgeted for.'" Verbatim protects the dissenter and creates accountability.
- **A specific commitment**: "Engineering committed: 'We will have the migration script in staging by Aug 14.'" Verbatim binds the speaker.
- **A regulatory / legal trigger**: anything that may need to be re-read later by counsel, audit, or a customer escalation.

Everywhere else, paraphrase. Verbatim quotes are expensive — they slow the reader and invite litigation of wording.

### 5. Decisions live in a log, not in the meeting file

A decision logged only inside one set of meeting minutes will be lost in six months. Significant decisions should be **extracted** into a separate decision log (or ADR file) with a stable ID, indexed, and linked back to the meeting where it was made. The meeting minutes then say: "Decided X — see DL-042."

This is the difference between *recording* a decision and *preserving* it.

### 6. The ADR form for technical decisions

For software/architecture decisions, the canonical form is Michael Nygard's 2011 template: **Title · Status · Context · Decision · Consequences**. (See Templates section.) It is intentionally short — 1 to 2 pages — and stored next to the source code in `docs/adr/NNNN-title.md`, source-controlled so it stays in sync with the code it justifies.

Status moves through: **proposed → accepted → deprecated → superseded by ADR-NNNN**. ADRs are never deleted or edited substantively after acceptance — they are superseded by a new ADR that links back.

### 7. Capture dissent on the record

If a participant disagreed with the decision, name them and one-sentence their objection. This is not gossip — it is provenance. Six months later when the decision is being questioned, future-you needs to know the dissent was heard and overridden (or not heard), not erased. Anonymous "concerns were raised" is useless and reads as cover.

### 8. Attendees, absentees, and the delegation chain

List attendees by name. List **invited-but-absent** separately — they did not get a vote. If a decision was made and a key stakeholder was absent, flag it: "Decision made without [Name]; [Name] to confirm by [date]." This prevents the predictable next-week pattern of "I never agreed to that."

### 9. Outcome-oriented language

Bad: "We talked about the latency issue for a while."
Good: "Decided to roll back the v2.3 caching change; Priya owns rollback PR by Wed."

Every bullet in the minutes should start with a verb of outcome (decided, agreed, deferred, rejected, assigned, scheduled, escalated) or a verb of action (review, ship, write, schedule, escalate). If a bullet starts with "discussed" or "talked about", it is probably not minute-worthy.

### 10. Distribute fast, freeze faster

Minutes that go out 5 days later are minutes nobody reads. Target: same-day or next-morning. Once distributed, treat them as frozen — corrections come in as a separate addendum email/thread, not as silent edits. Silent edits to minutes are how organizations lose institutional memory.

## Templates and Examples

### Template 1 — Sync meeting minutes

```markdown
# [Meeting name] — YYYY-MM-DD

**Attendees:** Name1, Name2, Name3
**Absent (invited):** Name4
**Scribe:** Name1
**Duration:** 30 min

## Decisions
- **[D1]** Decided to ship the dashboard refactor behind a feature flag in v1.4.
  - Why: avoids blocking the marketing launch on the 18th; allows rollback in <1 min.
  - Dissent: none.
- **[D2]** Decided NOT to backfill historical events older than 90 days.
  - Why: estimated 14 engineer-days for marginal user value.
  - Dissent: Jordan flagged "support team will get questions" — accepted as known risk.

## Action items
| # | Action | Owner | Due |
|---|---|---|---|
| A1 | Open feature-flag config PR for dashboard refactor | Priya | 2026-06-03 |
| A2 | Send customer-facing FAQ on 90-day backfill cutoff to Support | Jordan | 2026-06-05 |
| A3 | Confirm marketing launch date locks v1.4 ship window | Mitch | 2026-06-02 |

## Parked / not decided
- Pricing tier change — deferred to Q3 planning sync.

## Next meeting
2026-06-10, same time. Agenda owner: Priya.
```

### Template 2 — Async decision thread synthesis

Pin this as the first message in the thread once the decision lands:

```markdown
**SYNTHESIS (2026-05-29) — pinned**

**Decision:** Adopt OpenTelemetry for backend tracing, replace existing custom span library.
**Why:** Vendor-neutral, broad ecosystem, removes 1.2K lines of custom code.
**Owner:** Anand (migration), Lin (review).
**Timeline:** Migration spike by 2026-06-15; full cutover target 2026-08-01.
**Dissent on record:** Sam — "Migration cost is underestimated; we should pilot on one service first." Accepted; pilot will be the metrics service.
**Thread:** 47 messages above. Full ADR: docs/adr/0042-otel-adoption.md.
```

### Template 3 — ADR (Michael Nygard 2011 form)

```markdown
# ADR-0042: Adopt OpenTelemetry for backend tracing

- **Status:** Accepted — 2026-05-29
- **Deciders:** Anand, Lin, Mitch
- **Supersedes:** ADR-0017 (custom span library)

## Context

We currently emit traces through a custom span library written in 2022. It supports only HTTP and gRPC, and the maintenance burden has grown to ~0.5 FTE/quarter. Three new services need tracing in Q3 and would require extending the library. The vendor we use for trace storage (Honeycomb) now offers a first-class OTel ingest path.

## Decision

We will replace the custom span library with OpenTelemetry SDKs (Go and Python) across all backend services. Migration will run as a parallel emit phase (custom + OTel) for one sprint, then cut over to OTel-only.

## Consequences

**Positive**
- Removes ~1.2K lines of custom code and the half-FTE maintenance burden.
- Unlocks the OTel ecosystem (Jaeger, Tempo, Datadog if we ever swap vendors).
- Native support for the three Q3 services.

**Negative**
- Migration cost estimated at 6 engineer-weeks; some of this could slip.
- Short-term double-emit will increase trace storage cost by ~15% for one sprint.
- Team needs to learn OTel SDK conventions; previous library was bespoke.

**Neutral**
- Trace data format changes; existing saved queries in Honeycomb will need to be rewritten.
```

### Template 4 — Lightweight decision-log row

For a running decision log (one row per decision, in a doc or sheet):

```
ID    | Date       | Decision (1 line)                                          | Why (1 line)                            | Owner | Status   | Link
DL-042| 2026-05-29 | Adopt OpenTelemetry for backend tracing                    | Remove custom library; vendor-neutral   | Anand | Accepted | adr/0042.md
DL-043| 2026-05-29 | Defer pricing tier rollout to Q3                           | Marketing launch conflict in June       | Mitch | Deferred | minutes/2026-05-29.md
```

## Anti-Patterns

- **Transcript-as-minutes** — paragraphs of "Then Bob said... then Alice responded..." Minutes are an outcome record, not a play-by-play.
- **Anonymous dissent** ("concerns were raised") — if dissent is worth noting, the dissenter is worth naming.
- **Decisions buried in the discussion** — every decision should be visible in its own section, not embedded in a discussion bullet.
- **Action items without an owner or date** — "Team to follow up" is not an action item.
- **Co-owned action items** — "Priya and Jordan" means neither one. Pick one and let them delegate.
- **Silent edits after distribution** — once minutes are sent, fixes go in an addendum, not a quiet doc edit.
- **Conflating agenda and minutes** — copying the agenda forward and filling in blanks produces a discussion record, not an outcome record. Start minutes from a fresh outcomes-shaped template.
- **One-meeting decision graveyard** — significant decisions that live only inside one meeting's notes are lost. Extract to a decision log / ADR.
- **Verbose ADRs** — an ADR that is 8 pages long is not an ADR, it is a design doc. Keep ADRs to 1–2 pages; if you need a design doc, write one and link to it from the ADR.
- **Editing accepted ADRs in place** — accepted ADRs are immutable. Supersede them with a new ADR that links back.

## Decision Heuristics

**Should this go in minutes or a decision log / ADR?**
- One-meeting tactical call (who runs Tuesday's demo) → minutes only.
- Affects work for >1 sprint, or anyone outside the meeting needs to know → decision log row.
- Affects code architecture, system boundaries, or a long-lived tradeoff → ADR.

**Verbatim or paraphrase?**
- Dissent, commitment, or potential legal/regulatory readback → verbatim.
- Everything else → paraphrase, hard.

**Sync or async minutes?**
- Meeting happened on a call → sync minutes, distributed same day.
- "Meeting" was a Slack thread or doc-comment chain → async synthesis pinned at top of the thread; do not also write a separate minutes doc.

**Is this an action item or a wish?**
- Has a named single owner and a specific date → action item.
- Missing either → wish. Either complete it or drop it.

**Should I name the dissenter?**
- Yes, unless they explicitly asked to be off the record. Erased dissent looks like cover six months later.

**How long should this ADR be?**
- If it is >2 pages, it is probably trying to be a design doc. Split: short ADR + linked design doc.

## References

- Michael Nygard, "[Documenting Architecture Decisions](https://www.cognitect.com/blog/2011/11/15/documenting-architecture-decisions)" (2011) — the original ADR essay; defines the Status / Context / Decision / Consequences template still used today.
- [Architectural Decision Records (adr.github.io)](https://adr.github.io/) — community hub for ADR templates and patterns; ThoughtWorks Technology Radar grades ADRs as "Adopt".
- Joel Parker Henderson, [architecture-decision-record GitHub repo](https://github.com/joelparkerhenderson/architecture-decision-record) — broad collection of ADR templates and real-world examples in multiple formats.
- Martin Fowler, "[Architecture Decision Record](https://martinfowler.com/bliki/ArchitectureDecisionRecord.html)" (bliki) — short canonical write-up; describes ADR usage at ThoughtWorks projects.
- Wrike, "[Meeting minutes template with action items](https://www.wrike.com/blog/action-items-with-meeting-notes-template/)" — practical action-item structure (owner + due date + single verb).
- Atlassian Confluence, "[Meeting Notes Template](https://www.atlassian.com/software/confluence/templates/meeting-notes)" — widely used template covering attendees, decisions, action items, next meeting.
