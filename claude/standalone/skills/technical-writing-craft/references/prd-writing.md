<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `prd-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: prd-writing
description: Product Requirements Document (PRD) craft — the PM-owned genre that defines what to build and why, before engineering proposes how. Covers the "problem / users / proposed solution / success metrics / open questions" template, MVP-vs-v1-vs-roadmap scoping discipline, the Marty Cagan "Inspired" lean approach, Lenny Rachitsky's 1-pager pattern, Basecamp Shape Up "pitch" as a fixed-appetite alternative, wireframes-vs-prose tradeoffs, stakeholder sign-off cadence, and anti-patterns (designed-by-committee, spec-as-discovery-substitute). TRIGGER when user says "write a PRD", "product requirements doc", "1-pager", "product spec for feature X", "Shape Up pitch", "scope this feature", "what's MVP vs v1", "stakeholder sign-off", "Cagan PRD", "Lenny template", "product brief", or asks to define what a feature should do from a product/user perspective. SKIP for engineering architecture proposals (use rfc-and-design-docs), implementation specs/contracts (use spec-writing), code-task plans (use code-plan-writing), agent task plans (use agent-plan-writing), software architecture decisions (use software-architect), user stories or acceptance criteria as standalone artifacts (use user-story-and-acceptance-criteria — PRDs reference but do not replicate this), API design (use api-design-patterns), executive memos and one-pagers (use executive-comms or writing-expert), and long-form whitepapers (use whitepaper-writing).
category: custom
tags: [writing, product-writing, prd, product-management]
---

# PRD Writing

## Overview

A Product Requirements Document (PRD) is a PM-owned artifact that defines **what** a product team will build and **why**, before engineering proposes **how**. PRDs sit upstream of RFCs, design docs, and implementation specs. They are read by engineers, designers, QA, marketing, support, sales, and leadership — so the genre serves a wider audience than any engineering document.

The modern PRD is short, validated, and lives close to the prototype. Marty Cagan himself now warns that the traditional heavy PRD has been "replaced by the high-fidelity prototype as the primary product spec," and cautions teams against "reverting to heavy artifacts as a substitute for product discovery." Lenny Rachitsky's 1-pager and Basecamp's Shape Up "pitch" represent the leaner end of the spectrum; classic Cagan and Atlassian/Confluence templates represent the structured-but-modern middle.

A PRD captures **validated decisions**. It does not perform validation. Validation happens through customer interviews, prototypes, user testing, and data analysis — the PRD records the conclusions.

## Core Concepts

### 1. Problem-first structure (Lenny's hierarchy)

Lenny Rachitsky: "Nailing the problem statement is the single most important step in solving any problem. It's deceptively easy to get wrong, and when done well it's a superpower of the best leaders."

A modern PRD opens with the problem, not the solution. Sections proceed in this order:

1. **Problem** — what user/business pain are we solving, with evidence
2. **Users** — who is affected, segment size, current workarounds
3. **Goals** — desired outcomes and what success looks like
4. **Proposed solution** — at a high level (often a prototype link)
5. **Success metrics** — measurable definition of done
6. **Open questions / risks** — what's unresolved
7. **Non-goals** — what we are explicitly NOT doing
8. **Rollout / milestones** — MVP → v1 → roadmap

The problem statement must stand alone. If a reviewer reads only that section and disagrees with the framing, the rest of the document is moot.

### 2. The Cagan four-section minimum

Marty Cagan's classic PRD structure has four sections: **Purpose**, **Features**, **Release Criteria**, **Rough Timing**. In modern Cagan ("Inspired", "Empowered", SVPG essays), this collapses further: the high-fidelity prototype IS the spec for features, and the PRD wraps it with purpose, release criteria, and validated insights.

The Cagan release criteria are non-negotiable acceptance bars (performance, accessibility, security, localization) — NOT user-facing features. Mixing them confuses the document.

### 3. MVP vs v1 vs roadmap scoping discipline

Three distinct scopes with three distinct purposes:

- **MVP** — minimum viable: smallest scope that lets us learn whether the proposed solution works. Often unshippable to all users; behind a flag or to a beta cohort.
- **v1 (GA)** — broadly shippable: meets release criteria, addresses the core user need, no known showstoppers.
- **Roadmap (vNext)** — follow-on work that the PRD acknowledges but does not commit to.

Anti-pattern: writing one undifferentiated feature list and labeling it "the PRD." Reviewers cannot tell what is required to launch versus aspirational. Always tag each feature with its scope.

### 4. Shape Up "pitch" as a fixed-appetite alternative

Basecamp's Shape Up framework replaces the PRD with a **pitch**: a short document containing problem, **appetite** (2 or 6 weeks of fixed budget), solution sketch (fat marker, not Figma), **rabbit holes** (risks to bound), and **no-gos** (explicit exclusions).

Key inversion: in Shape Up, **scope is the variable; time is fixed**. The pitch defines an appetite (not an estimate), and the team trims scope to fit. In a traditional PRD, scope is fixed and time flexes.

Use the pitch when: the team has authority to cut scope mid-build, the work is bounded (2–6 weeks), and discovery is already complete. Use a traditional PRD when: multiple teams coordinate, dates are externally committed, or scope cannot flex below a known floor.

### 5. Wireframes vs prose decision

When to use wireframes/prototype links:
- The solution is primarily a UI change
- The interaction sequence carries the design intent
- A reviewer can grasp the proposal faster from a picture

When to use prose:
- Backend/data/policy changes invisible to UI
- Conditional logic, edge cases, state transitions
- Cross-feature interactions
- Anything you want to be searchable/diff-able later

Modern PRDs lean on prototypes for UX and prose for behavior. A "complete" PRD usually has both — a prototype link AND a behavior section.

### 6. Success metrics: leading vs lagging

Every PRD declares how it will know it worked. Two metric classes:

- **Leading indicators** — observable within 1–4 weeks (adoption %, feature engagement, task completion rate, time-to-value)
- **Lagging indicators** — observable in 1–2 quarters (retention, revenue, NPS, churn)

A PRD without success metrics is a wish list. A PRD with only lagging metrics cannot be evaluated within the build cycle. Include at least one leading indicator that the team can act on during the first month post-launch.

### 7. Stakeholder sign-off pattern

Lenny's PRD process: share the draft with the entire project team, collect feedback via comments/email/sync, integrate, re-share, repeat. Continuously refer back to the problem statement during build to keep alignment.

Practical sign-off layers (named, not anonymous):

- **Product** — PM author and product lead
- **Engineering** — tech lead (validates feasibility within scope)
- **Design** — design lead (validates UX direction)
- **Cross-functional** — legal/security/support/sales as the feature requires
- **Leadership** — group PM or director (validates priority and scope)

Record sign-off in the doc itself with date and name. Re-trigger sign-off on material scope changes.

### 8. Non-goals as a first-class section

Engineers and reviewers consistently raise questions like "what about X?" Most are inside the solution's natural surface area but outside the current scope. A **Non-Goals** section answers these preemptively:

> Non-goals (v1):
> - Mobile app support — desktop only
> - Bulk import — single-record only
> - Real-time sync — daily batch acceptable
> - Multi-tenant admin — single-org assumption

Non-goals reduce review-cycle thrash and protect the team from scope creep during build.

### 9. PRD is not RFC, not spec, not plan

Hard genre boundaries:

- **PRD** (this skill) — PM-owned. WHAT to build, for WHOM, WHY now. User-facing outcomes.
- **RFC / Design doc** — engineering-owned. HOW we propose to build it. Architecture, tradeoffs, alternatives considered.
- **Spec** — engineering-owned. The contract. API shapes, behavior rules, data formats.
- **Plan / Task list** — engineering-owned. Sequenced units of work with owners and dates.

A PRD that drifts into architecture choices loses authority — it tells engineering how to do their job. An RFC that drifts into user-problem framing duplicates the PRD. Keep the genres separate; cross-link them.

### 10. The "designed-by-committee" failure mode

If a PRD accumulates feedback from N stakeholders and the author tries to honor every comment, the doc becomes incoherent: contradictory constraints, vestigial sections, hedged language, missing point of view. The PM is the **author**, not a scribe.

Healthy pattern: collect feedback, summarize disagreements explicitly, **make a call**, and record the call with one-sentence rationale. Stakeholders may continue to disagree but the doc has a clear position. "We considered X and chose Y because Z" is stronger than "we will support both X and Y."

## Templates and Examples

### Modern PRD skeleton (Cagan + Lenny hybrid)

```markdown
# [Feature Name] — PRD

**Owner:** [PM name] · **Status:** Draft / In Review / Approved
**Last updated:** YYYY-MM-DD · **Target ship:** Quarter/Year

## TL;DR (3 sentences max)
[One sentence: the problem. One: the solution. One: the metric of success.]

## Problem
[User pain or business gap, with evidence: support volume, NPS verbatims,
analytics, sales-loss reasons, user interview quotes. Cite sources.]

## Users
- Primary segment: [who, how many, current workaround]
- Secondary segment: [who, how many]
- Out of scope: [explicit segments NOT served by v1]

## Goals
- [Outcome 1 — user-visible]
- [Outcome 2 — business]
- [Outcome 3 — operational, if relevant]

## Non-goals (v1)
- [Capability we are explicitly NOT building this round]
- [Capability the user might expect but we are deferring]

## Proposed solution
[2–4 paragraphs at concept level. Link to high-fidelity prototype.
Reference user-story-and-acceptance-criteria doc rather than duplicating.]

Prototype: [link]
User stories: [link]

## Behavior (prose section, for non-UI logic)
- [Edge case 1 — what happens]
- [State transition — what triggers it]
- [Conditional path — when does X apply]

## Success metrics
**Leading (4-week window):**
- [Metric 1 — target threshold]
- [Metric 2 — target threshold]

**Lagging (quarter window):**
- [Metric 3 — target threshold]

## Rollout
- **MVP (beta, X users):** [smallest learnable scope]
- **v1 (GA):** [full launch scope]
- **vNext (roadmap, not committed):** [follow-on hooks]

## Release criteria (Cagan)
- Performance: [SLO]
- Accessibility: [WCAG level]
- Security: [review completed]
- Localization: [scope]

## Risks and open questions
- [Risk 1] — mitigation: [...]
- [Open Q] — owner: [...] — decide by: [date]

## Decisions log
- YYYY-MM-DD: Considered X, chose Y because Z. (signed: PM)
- YYYY-MM-DD: ...

## Sign-off
- [ ] Product (name, date)
- [ ] Engineering (name, date)
- [ ] Design (name, date)
- [ ] [Other stakeholder] (name, date)
```

### Lenny 1-pager skeleton (smaller projects)

```markdown
# [Project] — 1-Pager

**Problem:** [one paragraph, with evidence]
**Why now:** [one sentence]
**Users:** [who, segment size]
**Proposed solution:** [one paragraph + prototype link]
**Success looks like:** [one leading + one lagging metric]
**Won't do:** [3–5 explicit non-goals]
**Risks:** [top 2, with owners]
**Decision needed:** [what call, by when, from whom]
```

### Shape Up pitch skeleton (fixed-appetite alternative)

```markdown
# Pitch: [Feature]

## Problem
[Raw user pain, 1–2 paragraphs. Specific, not abstract.]

## Appetite
[Small batch (2 weeks) | Big batch (6 weeks)]
Time is fixed. Scope is the variable.

## Solution
[Fat-marker sketch / breadboard. Not high-fidelity. Enough to bound the work.]

## Rabbit holes
- [Specific risk 1] — bounded by: [decision X is out of scope]
- [Specific risk 2] — bounded by: [decision Y is out of scope]

## No-gos
- [Excluded capability 1]
- [Excluded capability 2]
- [Excluded edge case]
```

### Worked example (real PRD opening, redacted)

```markdown
# Saved Filters — PRD

**TL;DR:** Power users rebuild the same filter combinations multiple times
per week (avg 14 rebuilds/user/week per analytics). We will let users save
named filter sets and recall them in one click. Success = 30% of weekly
active users save at least one filter within 4 weeks of launch.

## Problem
Filtering on the main list view requires 4–7 clicks across 3 dropdowns.
Support sees ~80 tickets/month from users asking "how do I get back to
the same view I had yesterday." 12 user interviews confirmed this is the
#3 reported friction (after onboarding and search). Top quote: "I just
take a screenshot of my filters so I can rebuild them tomorrow."

## Non-goals (v1)
- Sharing filters across users (deferred to v2)
- Default filter on login (deferred)
- Filter as URL parameter (deferred)
- Mobile filter management (web only this round)
```

## Anti-patterns

1. **Solution-first opening** — leading with "we will build X" before establishing why. Reviewers can't evaluate fit if the problem isn't defined. *Fix:* mandate problem-first structure.

2. **Designed-by-committee text** — incorporating every comment without making a call. Doc becomes contradictory and toothless. *Fix:* PM owns the decision; record the decision and one-sentence rationale.

3. **PRD as discovery substitute** — writing a 12-page PRD to "figure out" what users want. Cagan's warning: validation happens in prototypes and user research, not in docs. *Fix:* do the research first; the PRD records conclusions.

4. **Spec-creep** — PRD drifts into API shapes, schema choices, or implementation tactics. Engineering loses agency; PM commits to decisions outside their domain. *Fix:* link to a separate RFC/spec; keep PRD at the WHAT level.

5. **No non-goals** — every reviewer asks "what about X?" and the PM either expands scope or hand-waves. *Fix:* maintain an explicit non-goals section; update it as questions arrive.

6. **Lagging-metric-only success** — PRD claims success will be measured by retention/revenue, observable only after the next quarter. Team has no signal during the launch window. *Fix:* always include at least one leading indicator.

7. **MVP = everything but slightly less** — labeling the full feature list as MVP. Real MVP is the smallest scope that produces learning. *Fix:* test the MVP definition by asking "what would we learn if ONLY this shipped?"

8. **Stale doc** — PRD is written once and never updated as decisions are made during build. The doc becomes a lie. *Fix:* maintain a Decisions Log section; update it on every material call.

9. **Anonymous sign-off** — "approved by stakeholders" with no names. Cannot escalate, cannot trace authority. *Fix:* name + date for every sign-off line.

10. **Prototype-or-prose, not both** — relying only on a Figma link (no searchable behavior text) or only on prose (no visual interaction). *Fix:* both for non-trivial features.

## Decision Heuristics

| Situation | Use |
|---|---|
| Single team, 2–6 week scope, can flex scope | Shape Up pitch |
| Cross-team work or external dates | Traditional PRD |
| Small bounded project, 1 PM, fast review | Lenny 1-pager |
| Major feature affecting >1 team and >1 quarter | Full Cagan/Atlassian PRD |
| Greenfield with high discovery uncertainty | Prototype first, then PRD captures conclusions |
| Architecture choice is the hard question | Pair PRD with engineering RFC (separate doc) |
| Compliance / regulated feature | Full PRD + legal sign-off line + audit trail |

**One-page vs multi-page test:** can you state problem + solution + success metric in 5 sentences? If yes, use a 1-pager. If no, the situation needs a full PRD.

**PRD vs pitch test:** is time fixed or flexible? Fixed time + flexible scope → Shape Up pitch. Fixed scope + flexible time → traditional PRD.

**Wireframe vs prose test:** is the user-visible behavior 80% UI, or 80% logic/data/policy? UI-heavy → lead with prototype, supplement with prose for edge cases. Logic-heavy → lead with prose, supplement with screenshots.

## Cross-references

- **Adjacent skills:** `user-story-and-acceptance-criteria` (PRD references; does not duplicate), `rfc-and-design-docs` (engineering's architecture proposal that follows PRD), `writing-expert` (prose craft for the prose sections), `kill-the-AI-ism` (genre is prone to vague hedge language).
- **Downstream artifacts produced from a PRD:** RFC, engineering implementation spec (`spec-writing`), code plan (`code-plan-writing`), release plan, GTM brief.
- **Not this skill:** code review, architecture decisions, engineering plans, agent plans, executive memos.

## References

1. Marty Cagan, "Revisiting the Product Spec," Silicon Valley Product Group — argues the high-fidelity prototype is the modern product spec and warns against PRD as discovery substitute. https://www.svpg.com/revisiting-the-product-spec/
2. Lenny Rachitsky, "Examples and templates of 1-Pagers and PRDs" — Lenny's Newsletter; problem-statement-first hierarchy and 1-pager template. https://www.lennysnewsletter.com/p/prds-1-pagers-examples
3. Lenny's Product Requirements template, Atlassian Confluence templates library. https://www.atlassian.com/software/confluence/templates/lennys-product-requirements
4. Ryan Singer, "Write the Pitch," *Shape Up* (Basecamp) — chapter on pitch as PRD alternative; problem, appetite, solution, rabbit holes, no-gos. https://basecamp.com/shapeup/1.5-chapter-06
5. Ryan Singer, "The Betting Table," *Shape Up* (Basecamp) — appetite-vs-estimate, six-week cycles, scope-hammering. https://basecamp.com/shapeup/2.2-chapter-08
6. Aakash Gupta, "Product Requirements Documents (PRDs): A Modern Guide" — modern PRD structure synthesis. https://www.news.aakashg.com/p/product-requirements-documents-prds

## When to Use This Skill

Use when the user asks to: write or critique a PRD, scope a feature into MVP/v1/roadmap, draft a 1-pager, write a Shape Up pitch, draft a problem statement, plan stakeholder review for a product doc, decide between prototype and prose, decide between pitch and PRD, or clarify the boundary between PRD/RFC/spec.

Skip when the user asks for: engineering architecture (use rfc-and-design-docs), implementation contracts (use spec-writing), an engineering task list (use code-plan-writing), an agent task plan (use agent-plan-writing), API behavior documentation (use api-design-patterns or api-docs-craft), a one-page executive memo unrelated to product (use executive-comms or writing-expert), or a long-form whitepaper (use whitepaper-writing).
