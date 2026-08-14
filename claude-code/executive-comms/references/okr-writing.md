<!-- hub-reference-banner -->
> **Reference file — part of the `executive-comms` hub.** Formerly the standalone `okr-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: okr-writing
description: Objectives and Key Results craft. The ambitious-but-not-impossible rule, measurable-KR test (numerator + denominator + deadline), OKR vs KPI distinction, the input-vs-output-vs-outcome hierarchy, the stretch OKR 70% achievement convention, the rollup pattern across team/org levels, mid-quarter check-in language, retrospective scoring, and common bad-OKR patterns ("ship the X feature" is a task not an OKR). TRIGGER when the user says "write OKRs", "draft an OKR", "review my OKRs", "objective and key results", "Q3 OKRs", "team OKRs", "stretch goal", "measure what matters", "is this a good KR", "OKR vs KPI", "OKR scoring", "OKR retrospective", "key result", "rollup OKRs", "OKR examples for engineering / product / marketing", or asks for help turning a goal into a measurable target. SKIP for: KPI dashboards and operational metric design that are not part of a goal-setting cycle; executive memos / board updates (use executive-comms); product roadmap or PRD work (use prd-writing); generic strategy narratives (use writing-expert); long-form planning of an implementation (use code-plan-writing or agent-plan-writing).
---

# okr-writing

## Overview

OKRs (Objectives and Key Results) are a goal-setting framework popularized by Andy Grove at Intel, taught to John Doerr, and adopted at Google and beyond. The form is deceptively simple — one inspirational Objective, three to five measurable Key Results — but most OKRs in the wild are broken. They are project plans wearing OKR clothing.

The job of this skill: turn a wish ("we want to grow") into an OKR you can score at quarter end without arguing about whether you hit it. That requires three disciplines:

1. The Objective is **qualitative, inspirational, time-bound** — it tells the team why they should care.
2. Each Key Result is **measurable** — a numerator over a denominator, by a date. If you cannot put a number on it, it is not a KR.
3. The KR measures an **outcome** (value delivered), not an **output** (work done) and never an **input** (effort spent).

This skill covers writing, reviewing, rolling up across the org, and scoring at the end of the cycle.

## Core Concepts

### 1. The OKR equation: "I will [Objective] as measured by [Key Results]"

John Doerr's compact form. The Objective says *what* and *why it matters*. The Key Results say *how we will know we did it*.

- **Objective**: memorable, qualitative, inspirational, has an end date.
- **Key Result**: quantitative, time-bound, measures outcome.

Three to five KRs per Objective. Fewer than three usually means the Objective is too narrow. More than five means you have not picked.

### 2. The measurable-KR test: numerator, denominator, deadline

Every KR must answer:

- **Numerator** — what is being counted? ("active paid teams", "median p95 latency", "NPS responses scored 9–10")
- **Denominator** — relative to what? ("out of all signups in Q3", "across the top 50 enterprise accounts")
- **Deadline** — by when? ("by 2026-09-30")

If a draft KR is missing any of the three, it fails the test. Examples:

- **Fails:** "Improve onboarding." (no numerator, no denominator, no deadline)
- **Fails:** "Reach 1000 users." (no denominator, no deadline)
- **Passes:** "Grow week-1 activation rate from 28% to 45% by 2026-09-30." (numerator = activated accounts, denominator = signups, deadline = Sept 30)

### 3. OKR vs KPI — different jobs, do not confuse them

- **KPI** (Key Performance Indicator) — a steady-state health metric you watch *all the time*. ("p95 latency under 200ms", "monthly churn under 2%"). KPIs run forever.
- **OKR** — a change goal for a *bounded period* (usually a quarter). It says "we are choosing to push this number from X to Y this quarter."

KPI = the speedometer. OKR = the destination this leg of the trip.

A KPI **can become** a KR for one quarter if you decide to push it — but the KR must specify the from-to delta and the deadline. "Maintain p95 < 200ms" is a KPI. "Drive p95 from 350ms to under 200ms by Sept 30" is a KR.

### 4. Input → Output → Outcome hierarchy

Three layers, in increasing order of value:

- **Input** — effort, headcount, money spent. ("hired 5 engineers", "ran 12 workshops") — never a KR.
- **Output** — work produced, artifacts shipped. ("shipped the new dashboard", "published 8 blog posts") — almost never a KR.
- **Outcome** — change in the world, value delivered. ("dashboard adoption reached 60% of paying teams", "blog drove 15K qualified signups") — this is what KRs measure.

Felipe Castro's "so what?" test: read the KR and ask "so what?" If the answer is another metric (the outcome), the original was an output. Replace it.

Outputs and inputs belong in a project plan, not an OKR.

### 5. Ambitious-but-not-impossible: the 70% rule

For **aspirational / stretch** OKRs (Wodtke / Google convention), you should:

- Pick a target where your confidence of hitting 100% is about **5 out of 10** (50/50 at draft time).
- Score around **0.7** on average at the end of the quarter is healthy.
- Consistently scoring 1.0 means you sandbagged. Scoring < 0.4 means you set fantasy targets.

For **committed** OKRs (operational must-haves — SLAs, compliance deadlines), the target is 1.0 and anything less is a problem.

Most teams should distinguish the two: a small set of committed OKRs (1.0 expected) plus a small set of aspirational ones (0.7 expected). Mixing them under one rubric breaks scoring.

### 6. The rollup pattern across org levels

OKRs at each level **align with** but do not **directly inherit from** the level above.

- Company sets 3–5 Objectives for the quarter.
- Each team picks ~3 Objectives that contribute to the company set. A team Objective may serve more than one company Objective, or one company Objective may be served by several team Objectives.
- An individual contributor may have 1–2 personal Objectives that map to a team Objective.

**Key insight (Doerr):** roughly half of OKRs should be set bottom-up. Pure top-down OKRs kill engagement. Pure bottom-up OKRs lose alignment. Aim for a negotiation.

A team KR can become a higher-level KR — *if* you accept that the higher level then depends on this team to hit. This is how organizational rollups work and it is intentional.

### 7. Mid-quarter check-in: traffic-light language, not narrative

Halfway through the quarter, score each KR with one of:

- **Green** — on track, no help needed.
- **Yellow** — at risk, here is what would unblock me.
- **Red** — will not hit current target. Either re-plan or formally revise the KR.

This is the only check-in language you need. Long narratives at mid-quarter usually hide red status. The "what would unblock me" sentence is the load-bearing part — that is where leadership earns their pay.

### 8. End-of-cycle scoring (0.0 to 1.0)

For each KR, compute the actual / target on its native scale:

- Started at 28%, target 45%, ended at 38%. Progress = (38 − 28) / (45 − 28) = 10/17 ≈ **0.59**.
- Started at 200ms, target 100ms, ended at 130ms. Progress = (200 − 130) / (200 − 100) = 70/100 = **0.70**.

The Objective score is the (often weighted) average of its KRs.

A retrospective should answer three questions per Objective:

1. What was the score, and why that score?
2. What changed in our understanding of the problem?
3. Are we keeping, dropping, or revising this Objective for next quarter?

### 9. Common bad-OKR patterns

- **"Ship feature X by date Y"** — that is a task, not a KR. The KR is the *outcome* the feature is supposed to produce: adoption %, conversion lift, ticket-deflection rate. Shipping the feature is the work; the OKR measures whether the work paid off.
- **"Do our normal job well"** — operational baselines are KPIs, not OKRs. OKRs are about deliberate change.
- **Activity counts** — "publish 12 blog posts" measures effort, not outcome. Replace with "drive 15K qualified signups via content channel."
- **Vanity metrics** — "1M pageviews" with no connection to retention or revenue. Pick metrics tied to value.
- **Sandbagged targets** — KRs you are 95% confident in. You will hit them and learn nothing.
- **Too many KRs** — more than 5 per Objective is a wishlist. Pick.
- **One-team-blocking-another KRs** — a KR that depends entirely on a team you do not control is a wish about that team. Either pull the dependency inside your team or rewrite as your team's contribution.
- **Set-and-forget** — OKRs scored only at quarter end. The mid-quarter check-in is non-negotiable.

### 10. The "is this even an OKR?" gate

Before drafting KRs, ask: **does this Objective require the team to do something differently than last quarter?** If the answer is "no, just keep doing what we're doing," it is a KPI to monitor, not an OKR to set. OKRs are about deliberate change in a bounded window. They are expensive — only spend them on the changes that matter.

## Templates and Examples

### Template — Single OKR

```markdown
## Objective
[Qualitative, inspirational, time-bound. One sentence. Why this matters.]

## Type
[ ] Committed (target = 1.0)
[ ] Aspirational / stretch (expected score ≈ 0.7)

## Key Results
1. [Move metric X from A to B by DATE]
2. [Move metric Y from A to B by DATE]
3. [Move metric Z from A to B by DATE]

## How we'll know we got it wrong
[One sentence — the leading indicator that says "stop, replan".]

## Dependencies / risks
- [Team or system we depend on]
- [Known risk]
```

### Example 1 — Product team, aspirational

```markdown
## Objective
Make new-team onboarding so fast that day-1 users feel the product's value before standup ends.

## Type
Aspirational (expected ≈ 0.7)

## Key Results
1. Increase week-1 activation rate from 28% to 45% by 2026-09-30.
2. Reduce median time-to-first-successful-action from 18 min to under 6 min by 2026-09-30.
3. Lift NPS among accounts ≤ 14 days old from 22 to 40 by 2026-09-30.

## How we'll know we got it wrong
Activation drops below 25% in any week — onboarding redesign actively regressing.

## Dependencies / risks
- Marketing landing-page revamp (Marketing team, due Aug 1).
- iOS app gate on review process for in-app onboarding tour.
```

### Example 2 — Engineering team, committed

```markdown
## Objective
Bring backend reliability up to the SLA we sold enterprise customers in Q1.

## Type
Committed (target = 1.0)

## Key Results
1. Reduce p95 API latency from 350ms to under 200ms by 2026-09-30, sustained for 14 consecutive days.
2. Cut Sev-2 incident count from 6/quarter to ≤ 2/quarter, measured 2026-07-01 to 2026-09-30.
3. Reach 99.95% measured uptime on the customer-facing API over the rolling 90-day window ending 2026-09-30.
```

### Example 3 — Bad → Good rewrites

| Bad | Why bad | Good |
|---|---|---|
| "Launch the new dashboard." | Output, not outcome; no metric, no deadline. | "Reach 60% weekly active usage of the new dashboard among paying teams by 2026-09-30." |
| "Improve customer happiness." | No numerator/denominator/deadline. | "Lift CSAT among top-50 enterprise accounts from 7.4 to 8.5 by 2026-09-30." |
| "Hire 10 engineers." | Input, not outcome. | "Reduce p95 API latency from 350ms to <200ms by 2026-09-30." (Hiring is *how*, not the OKR.) |
| "Publish 12 thought-leadership blog posts." | Activity count. | "Drive 15,000 qualified signups via the content channel by 2026-09-30." |
| "Migrate to OpenTelemetry." | Project plan. | "Complete OTel migration with <0.5% trace-data loss vs. baseline, by 2026-09-30." (Even here, the better OKR is the outcome the migration enables — debugging speed, vendor flexibility.) |

### Example 4 — Mid-quarter check-in line

```markdown
**KR1** (Activation 28% → 45%): YELLOW.
At 34% as of week 7. New onboarding flow ships Aug 2; expect lift to follow by ~2 weeks.
Unblock ask: need design review slot from Lin this week to finalize the empty-state copy.
```

### Example 5 — Retrospective scoring

```markdown
## Objective: Make new-team onboarding fast (final score: 0.65)

- KR1 (28% → 45%): ended at 38%. Score 0.59.
- KR2 (18min → 6min): ended at 7min. Score 0.92.
- KR3 (NPS 22 → 40): ended at 28. Score 0.33.

**What changed in our understanding:**
Time-to-first-action was the easy win; NPS is a lagging indicator that won't move until renewals.

**Next quarter:**
Keep the Objective. Drop KR3 (NPS lags). Add a retention KR (week-4 retention).
```

## Anti-Patterns

- **Tasks-as-OKRs** — "Ship X by Y" is a project milestone, not an outcome. Replace with the metric the ship is meant to move.
- **Outputs-as-KRs** — "Published 12 posts", "Ran 8 trainings", "Closed 100 tickets". Count things customers care about, not things you did.
- **Inputs-as-KRs** — "Hire 5 engineers", "Spend $200K on ads". Effort, not outcome.
- **KPI-as-OKR** — putting a steady-state metric into the OKR slot without specifying the from-to delta. Not a change goal, so not an OKR.
- **No baseline** — "Increase activation to 45%" without saying it is at 28% today. The from-to delta *is* the goal.
- **Too many OKRs** — more than 3–5 Objectives per cycle means no priorities. Pick.
- **Sandbagged targets** — pre-met or trivially-met KRs. Aspirational OKRs should feel uncomfortable at draft time.
- **All top-down** — kills engagement. Aim for roughly half team-originated.
- **Set and forget** — no mid-quarter check-in. By the time you score at end-of-quarter it is too late to course-correct.
- **Silent revision** — quietly moving the target down in week 9 because you will not hit it. Either re-plan publicly or score against the original.
- **Conflated committed and aspirational** — one is a 1.0 target, the other is a 0.7 target. Pick which type each OKR is and score accordingly.

## Decision Heuristics

**Is this a KPI or an OKR?**
- Steady-state metric I want to keep healthy → KPI.
- Specific delta I want to drive this quarter → OKR.

**Is this KR an outcome, output, or input?**
- Apply Castro's "so what?" test. If "so what?" yields another metric, that metric is the real KR.
- If it counts effort or artifacts, downgrade it.

**Committed or aspirational?**
- Must-hit (SLA, compliance, contractual) → committed, target 1.0.
- Stretch (50/50 confidence at draft time) → aspirational, expected ≈ 0.7.

**Is the target ambitious enough?**
- If you are >90% confident at draft time, raise the target.
- If you are <20% confident, lower it or break it into a smaller bet.

**Should this be one OKR or two?**
- If two KRs are measuring genuinely different outcomes — split.
- If five KRs are all flavors of the same outcome — collapse.

**Mid-quarter: revise or push through?**
- Red because of an environment change (market, dependency) → revise publicly with a note.
- Red because we underestimated the work → push through, score honestly at end of quarter, learn.

**Should this team's KR also be a higher-level KR?**
- Yes, if the team owns the outcome end-to-end and leadership wants visibility.
- No, if multiple teams contribute — roll up at the parent Objective layer, not by stacking KRs.

## References

- John Doerr, [whatmatters.com — What is an OKR?](https://www.whatmatters.com/faqs/okr-meaning-definition-example) — canonical "I will X as measured by Y" form; committed vs aspirational distinction.
- Christina Wodtke, [Radical Focus 2.0 (Mind the Product excerpt)](https://www.mindtheproduct.com/radical-focus-2-0-an-excerpt-from-christina-wodtke/) and [The Art of the OKR](https://cwodtke.com/the-art-of-the-okr/) — 50/50 confidence rule; stretch-goal mechanics; weekly cadence.
- Felipe Castro, ["An OKR should measure the outcome, not the work"](https://read.felipecastro.com/p/measure-outcome-not-work) and ["Key Results Should Pass the 'So What?' Test"](https://medium.com/@meetfelipe/key-results-should-pass-the-so-what-test-be6f9d36bac6) — outcome vs output discipline; bad-OKR diagnosis.
- Atlassian, [OKRs: The Ultimate Guide to Objectives and Key Results](https://www.atlassian.com/agile/agile-at-scale/okr) — practical rollup pattern, examples, common pitfalls.
- Workpath, [Criteria & Tips for Drafting Professional OKRs](https://www.workpath.com/en/magazine/okr-quality-check) — quality checklist drawing on Castro/Wodtke, including the task-vs-outcome pitfall.
