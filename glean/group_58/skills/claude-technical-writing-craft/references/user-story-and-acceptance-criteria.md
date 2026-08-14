<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `user-story-and-acceptance-criteria` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: user-story-and-acceptance-criteria
description: Author user stories and acceptance criteria that survive grooming and ship as working software — Mike Cohn's "As a / I want / So that" template (with the role-as-specific-actor rule), INVEST principles (Independent / Negotiable / Valuable / Estimable / Small / Testable, Bill Wake 2003), Given/When/Then acceptance criteria in Gherkin form, the vertical-slice rule (end-to-end value per story, not horizontal layers), SPIDR splitting (Spike / Path / Interface / Data / Rules, Mountain Goat Software), and the Definition-of-Done vs Acceptance-Criteria distinction. TRIGGER when user asks "write a user story", "format this as a user story", "split this story", "acceptance criteria for", "Given/When/Then", "is this story small enough", "INVEST check", "DoD vs AC", "vertical slice", "SPIDR", "groom the backlog", or pastes a feature description asking for backlog-ready stories. SKIP when authoring an architectural design doc that precedes the work (use rfc-and-design-docs or software-architect); writing the engineering implementation plan (use code-plan-writing or agent-plan-writing); creating Jira tickets or driving a spec-to-backlog flow in Atlassian tools (use jira-ticket-creation or atlassian:spec-to-backlog); writing PR descriptions or commit messages for completed work (use pr-description-craft or commit-message-craft); writing customer-facing release notes (use changelog-and-release-notes or changelogs-for-humans); writing prose-only documents that aren't backlog items (use writing-expert).
---

# User Story and Acceptance Criteria

## Overview

**What this skill covers.** Take a feature idea — a paragraph of intent from a PM, a customer request, a meeting note — and produce backlog items the team can groom, estimate, and ship. Covers Mike Cohn's user-story template ("As a *role*, I want *capability*, so that *outcome*"), the INVEST quality bar for backlog items (Bill Wake, 2003), Given/When/Then acceptance criteria written in Gherkin form, splitting strategies (vertical-slice principle and SPIDR), and the conceptual separation between *acceptance criteria* (per-story, varies) and *definition of done* (team-wide, stable).

**When to use.**
- The user wants a feature idea formatted into one or more user stories.
- The user has a "too big" story and needs splitting strategies.
- The user wants acceptance criteria written for an existing story.
- The user asks INVEST-check questions ("is this story Independent? Testable?").
- The user is grooming a backlog and needs help shaping items.
- The user is mixing up Definition of Done with Acceptance Criteria.

**When to skip.**
- Pre-implementation architecture or design — use `rfc-and-design-docs` or `software-architect`.
- The engineering implementation plan that translates the story into tasks — use `code-plan-writing` or `agent-plan-writing`.
- PR descriptions or commit messages for completed work — use `pr-description-craft` or `commit-message-craft`.
- Customer-facing release notes — use `changelog-and-release-notes`.
- Sentence-level prose review of a document that isn't a backlog item — use `writing-expert`.

## Core Concepts

### 1. The Mike Cohn user-story template — and why each clause matters

```
As a <role>, I want <capability> so that <benefit>.
```

Three clauses; each fails differently when written sloppily.

**`As a <role>`** — the user, not the system. "As a user" is the most common failure mode: it's so generic it adds no information. Cohn's recommendation: name a specific actor type. "As a job seeker", "As an employer", "As an on-call TAM", "As a paid subscriber on the Pro tier". A role is a *position-the-actor-is-in-when-they-want-this*, not a permission set.

**`I want <capability>`** — the *what*, written goal-first, agnostic to *how*. "I want to filter cases by severity" is goal-shaped. "I want a dropdown in the top-right of the case list that lets me pick S1/S2/S3" is solution-shaped — it pre-commits to a design decision.

**`so that <benefit>`** — the *why*. The benefit clause is the most-skipped and most-valuable. It anchors prioritization (is this benefit worth a sprint?) and disambiguation (when the implementation hits a fork, which path serves the stated benefit?).

**Does every story need a `so that`?** Cohn himself has said: most do, but a small number of obvious-benefit stories ("As an admin, I want to log in") don't need it spelled out. The test: *if a junior team member can't reconstruct the benefit from the role and capability, write the `so that`*.

### 2. INVEST — the quality bar for backlog items

Bill Wake coined INVEST in 2003 ("INVEST in Good Stories, and SMART Tasks", XP magazine) as a checklist for whether a backlog item is *ready* to estimate and pull.

**I — Independent.** The story can be built, demoed, and shipped without waiting on another story in the same iteration. Dependencies push to *between* iterations, not *within*. Test: can this story be the *only* story in the sprint and still deliver value?

**N — Negotiable.** The story is a placeholder for a conversation, not a contract. The team and product owner can adjust scope, swap details, or alter approach as they discover constraints during the sprint. A story crammed full of pre-locked specifications fails N.

**V — Valuable.** The story delivers value to a real user or stakeholder. If only the engineering team benefits, it's not a story — it's a task. Refactors and tech-debt items can still be Valuable if you can name the user-facing benefit (faster page load, fewer bugs in feature X).

**E — Estimable.** The team has enough context to size it. An un-Estimable story usually fails because (a) it needs a spike first, or (b) it's too big and needs splitting.

**S — Small.** Fits comfortably inside an iteration. Heuristic: ≤ 50% of one developer's iteration capacity. Stories bigger than that are "epics" that need splitting (see Concept 5).

**T — Testable.** A definite test exists for "done". If the team can't write the test before starting, the story's success condition is vague. Acceptance criteria are how Testable becomes concrete.

### 3. Given / When / Then — Gherkin acceptance criteria

Acceptance criteria in Gherkin form make INVEST's "T — Testable" concrete. Each criterion follows:

```
Given <some context>
When <some action>
Then <some observable outcome>
```

**Given** — the precondition. State the world before the action.
**When** — the trigger. The user does X, or the system receives Y.
**Then** — the observable consequence. *Observable* is load-bearing: it must be something a tester (human or automated) can check.

**`And` and `But`** chain steps in the same phase:

```
Given a user with role "admin"
  And the user is logged in
When the user clicks "Delete account" on user @bob
Then user @bob's account is soft-deleted
  And an audit-log entry is written
  But @bob's outgoing comments are preserved
```

**Rules:**
- 3–5 steps per scenario; longer scenarios lose specificity.
- 1–3 acceptance criteria per story; if you need 4+, the story is probably too big.
- Each criterion tests a *distinct* aspect, not minor variations of the same aspect.
- Concrete values, not generics: "Given a charge of $20.00" not "Given a valid charge".

**Non-Gherkin acceptance-criteria styles** (used when Gherkin feels heavy):
- **Checklist form:** plain bullet points. "Filter dropdown shows S1–S4. Default is 'all'. Filter persists across page reload." Faster to write; loses the Given/When/Then traceability.
- **Rule-oriented:** "Rule: refunds older than 90 days require manager approval." Pairs well with example tables.

### 4. The vertical slice rule

A user story must be a *thin vertical slice through the architecture* — a sliver that touches every layer (UI, API, business logic, data) and delivers end-to-end value — not a *horizontal layer* across the architecture.

**Horizontal (wrong):**
- Story 1: Build the UI for case filtering.
- Story 2: Build the API endpoint for case filtering.
- Story 3: Add the database index for case filtering.

None of those, individually, is shippable. The user gets value only when all three exist.

**Vertical (right):**
- Story 1: Filter cases by severity (S1 only, no UI persistence, no pagination, no index).
- Story 2: Filter cases by status, with severity already shipped.
- Story 3: Persist the last-used filter across sessions.

Each story is shippable on its own, even in degraded form. Each delivers some real-user value.

The deliberately-degraded first slice is the discipline: cut data, cut UI polish, cut performance — but *don't* cut the end-to-end path. A story that ships without database isn't a story; it's a mock.

### 5. SPIDR — five ways to split a story that's too big

When a story fails INVEST's `S` (Small) or `E` (Estimable), use SPIDR (Mike Cohn / Mountain Goat Software):

**S — Spike.** Time-box a research task to remove the uncertainty that's blocking estimation. The spike's deliverable is *knowledge*, not feature code. Result of the spike feeds the now-estimable real story.

**P — Path.** If a user can do the thing multiple ways, split by path. "Pay with credit card" / "Pay with Apple Pay" / "Pay with stored balance". Ship one path first; the others follow.

**I — Interface.** Split by client or platform. "Filter cases on the desktop dashboard" / "Filter cases in the mobile app". Or split by browser, OS, API version. Each slice is shippable to its own audience.

**D — Data.** Split by data scope. "Filter cases for active accounts only" first; "Filter cases including archived accounts" later. Or "Support the 5 top currencies" first; "support all 168" later.

**R — Rules.** Relax business rules in the first slice. "Refunds, with no approval workflow" first; "Refunds with the manager-approval workflow" later. Or "Allow uploads up to 10MB" first; "Allow uploads up to 5GB with chunked upload" later.

SPIDR techniques compose. A "too-big" story often splits via two letters at once: Path *and* Data, or Interface *and* Rules.

### 6. Acceptance criteria vs Definition of Done — different things, both required

The two terms are constantly conflated. They are *not* synonyms.

**Acceptance Criteria (AC):**
- Specific to *this* story.
- Authored by the product owner with the team.
- Vary widely between stories.
- Answer: "What must this story do for the user to accept it?"

**Definition of Done (DoD):**
- A team-wide standard that applies to *every* story.
- Authored once, evolves slowly.
- Stable across sprints.
- Answers: "What must any item meet to be called done?"

A typical DoD includes things like:
- Code reviewed by ≥1 other engineer.
- Unit tests written and passing.
- Integration tests passing.
- Documentation updated (changelog, API docs, user-facing docs as applicable).
- No new lint warnings.
- Deployed to staging.
- Accessibility audit passed (if UI).
- Performance budget honored.

A story is **done** when *both* its acceptance criteria are met *and* the team's definition of done is satisfied. Either alone is insufficient.

### 7. The 3 C's — Card, Conversation, Confirmation

Ron Jeffries' three-part anatomy of a story, often cited alongside Cohn:

- **Card** — the brief written artifact (the story sentence on the card, physical or digital).
- **Conversation** — the in-person discussion that fleshes out the details. The card is a *placeholder* for this conversation, not a substitute.
- **Confirmation** — the acceptance criteria that confirm the story is done.

Card → Conversation → Confirmation captures the rhythm: stories start sparse, get richer through discussion, and end with verifiable criteria.

The 3 C's matter because they remind teams that the *card text alone* isn't the story. A team that tries to write requirements into the card and skip the conversation produces overly-prescriptive cards that miss the actual user need.

### 8. Sizing — story points vs t-shirt sizes vs absolutes

Stories need a relative size for sprint planning. Three common approaches:

**Fibonacci story points (1, 2, 3, 5, 8, 13, 21):**
- The non-linear gaps (no 4, no 6, no 7) force differentiation.
- A 13-point story is "definitely much bigger than 8" — a signal to split.
- Velocity-based planning ("the team completes 30 points/sprint on average") falls out of this.

**T-shirt sizes (XS, S, M, L, XL):**
- Lower-stakes; fewer arguments about whether something is 5 or 8.
- Common in companies that want sizing without the velocity-tracking ritual.

**Story counts (no sizes, just "how many stories"):**
- "Sprint capacity is N stories." Works when stories are kept genuinely Small (the S in INVEST).
- Skips the estimation ceremony entirely. Popular with #NoEstimates advocates.

**Heuristic:** if a story is sized as the largest in the system (13 points, XL), treat it as a candidate for splitting. If the team consistently estimates everything as Medium, the team's not actually estimating.

### 9. Personas vs roles

Stories use *roles*, but the same backlog often references *personas* — fully-fleshed user archetypes (name, job, goals, frustrations). The distinction:

- A **role** is a position in a workflow ("approver", "submitter", "new user").
- A **persona** is an archetype ("Maria, 38, regional ops manager who reviews 12 reports a week and is impatient with slow loading").

Stories work with roles because roles are stable and the same actor can fill multiple roles. Personas live in product strategy docs and inform priority — they're a richer artifact than a story header can carry.

### 10. Anti-stories — what *isn't* a user story

Some backlog items shouldn't be force-fit into "As a / I want / so that":

- **Pure tech-debt items** ("Refactor the auth module"). Better as "Refactor card" or "Technical Story" — keep the form honest. If forced into user-story form, you get awkward fakes like "As a developer, I want to refactor the auth module so that future stories are easier" — fine in some teams, hollow in others.
- **Bugs.** Bugs have their own anatomy: reproduction steps, expected behavior, actual behavior, severity. Wrapping a bug in user-story syntax loses the structured fields.
- **Spikes.** A time-boxed research task with a knowledge deliverable. Cohn explicitly carves these out as the `S` in SPIDR.
- **Operational tasks** ("Rotate the staging cluster's certs"). These are tasks, not stories.

Tracking systems that force everything into user-story form (looking at you, certain Jira configurations) push teams to write hollow stories. Better: let bugs, spikes, and tasks have their own types.

## Templates and Examples

### Template — user story with acceptance criteria

```
**Title:** Filter case list by severity

**Story:**
As an on-call TAM,
I want to filter the case list by severity (S1 / S2 / S3 / S4),
so that during a busy on-call shift I can triage S1s first
without scrolling through lower-severity cases.

**Acceptance Criteria:**
1. Given the case list shows ≥ 1 case at each of S1–S4,
   When I select "S1 only" in the severity filter,
   Then only S1 cases are visible
     And the count badge shows the number of visible S1 cases.

2. Given I have applied a severity filter,
   When I reload the page,
   Then the same filter is reapplied
     And the URL contains the filter as a query parameter.

3. Given the filter is set to a value with zero matching cases,
   When the filtered list is rendered,
   Then an empty-state message is shown
     But no error is logged.

**Out of scope (not in this story):**
- Filtering by status or owner (see #1241, #1242).
- Multi-select severity (e.g., S1 + S2 together) — Path-split for v2.

**Definition of Done items called out:**
- a11y: filter is keyboard-navigable.
- Persistence: query-param-only for v1 (no chrome.storage).

**Estimate:** 5 points
```

### Template — story split using SPIDR

```
Original (too big, ~21 points):
"As an admin, I want to manage subscription plans, so that I can
adjust pricing tiers."

Split via Path + Rules:

Story A (5 pts) — Path: read-only plan list
  As an admin, I want to view all plans in a list,
  so that I can see the current pricing tier structure.

Story B (8 pts) — Path + Rules: create new plan
  As an admin, I want to create a new plan with a name and
  monthly price (USD only; no proration; no metadata),
  so that I can pilot new pricing tiers without engineering help.

Story C (5 pts) — Rules relaxed: edit plan name (only)
  As an admin, I want to rename a plan,
  so that I can correct typos and rebrand without engineering.

Story D (8 pts) — Rules: edit plan price with proration
  As an admin, I want to change a plan's price with proration
  rules applied to existing subscribers,
  so that price changes propagate fairly.

Story E (3 pts) — Path: deactivate a plan
  As an admin, I want to archive a plan,
  so that it stops appearing in the customer signup flow but
  existing subscribers keep their pricing.
```

### Example — INVEST check on a story

> "As an on-call TAM, I want to receive a Slack notification when any case I own goes 4+ hours without an update, so that I don't lose SLA cycles to forgotten cases."

- **I — Independent?** Yes — no upstream story is blocking this.
- **N — Negotiable?** Yes — Slack vs email vs in-app, threshold 4h vs 3h vs 6h, scope "cases I own" vs "all S1 cases in my queue" — all discussable.
- **V — Valuable?** Yes — concrete user (TAM) and concrete benefit (SLA preservation).
- **E — Estimable?** Maybe. The team needs to confirm: do we already have a Slack-notification pipeline? If no, this needs a spike.
- **S — Small?** If Slack pipeline exists, 5 points. If not, split: spike on Slack integration + a separate notification-content story.
- **T — Testable?** Yes — Given a case is owned by me and has no update for 4h, When the threshold elapses, Then I receive a Slack DM with the case link.

**Verdict:** ready to estimate IF the team confirms Slack pipeline exists. If not, split off a spike first.

### Example — Definition of Done (team-wide)

```markdown
# Team Definition of Done

A backlog item is **done** only when all of the following are true:

- [ ] Acceptance criteria met (verified by the PM or designate).
- [ ] Unit tests written for new logic; ≥ 80% coverage on changed files.
- [ ] Integration tests cover at least one happy-path and one error-path.
- [ ] Code reviewed by ≥ 1 other engineer; all blocking comments resolved.
- [ ] CHANGELOG.md updated (or marked N/A in the PR).
- [ ] User-facing docs updated (if UI or API surface changed).
- [ ] No new lint warnings; no new TypeScript `any`.
- [ ] Deployed to staging; smoke test passing.
- [ ] Accessibility audit clean (Lighthouse a11y ≥ 95) for UI changes.
- [ ] No new console errors in DevTools during the demo flow.
- [ ] Telemetry / observability added for new code paths.
```

### Example — Gherkin acceptance criteria for a payment flow

```gherkin
Feature: Process payment with card

  Scenario: Successful card payment with no fraud flags
    Given the user has selected a total amount of $120.00
      And the card details are valid and not expired
      And the fraud detection system has no flag for the transaction
    When the user clicks "Pay Now"
    Then the payment is processed successfully
      And the order status changes to "Paid"
      And the customer receives a confirmation email

  Scenario: Card declined for insufficient funds
    Given the user has selected a total amount of $120.00
      And the card details are valid but the card lacks funds
    When the user clicks "Pay Now"
    Then the payment fails with decline_code "insufficient_funds"
      And the order status remains "Pending"
      And the UI displays a generic "card declined" message
      But the UI does not display the decline reason to the user

  Scenario: Card declined as lost or stolen
    Given the user has selected a total amount of $120.00
      And the card is reported lost
    When the user clicks "Pay Now"
    Then the payment fails
      And the UI shows the same generic "card declined" message
      But the system silently flags the attempt to the fraud team
```

## Anti-Patterns

- **"As a user, I want..."** — every story starts the same way, every story tells you nothing about who the user actually is. Name the role specifically.
- **Solution-shaped capability clauses** — "I want a dropdown in the top-right corner that..." pre-commits to UI before the team has talked. Keep the `I want` goal-shaped.
- **No `so that` clause** — strips out the prioritization signal. Even obvious benefits earn a one-line spelling-out.
- **Horizontal-layer stories** ("Build the backend for X", "Build the UI for X") — each is independently unshippable. Slice vertically.
- **Acceptance criteria that restate the story** — "AC: the user can filter by severity" just paraphrases the `I want`. Acceptance criteria must be more *specific*, not more *general*.
- **20-criterion acceptance lists** — the story is too big. Split.
- **Conflating Acceptance Criteria with Definition of Done** — putting "tests pass" in every story's AC instead of in the team's DoD. AC is per-story; DoD is team-wide.
- **Stories with no `Then` clause that's actually observable** — "Then the system feels faster" isn't testable. "Then page-load completes within 1.5s at p95" is.
- **Skipping the Conversation in the 3 C's** — treating the card as the full requirement and writing it like a spec. The card is a prompt for talking.
- **Estimating an un-estimable story instead of spiking** — when nobody knows how big the thing is, the team eyes a number. Spike instead; estimate after.
- **Forcing bugs / tech-debt / spikes into user-story syntax** — produces hollow "As a developer, I want to refactor X, so that..." stories. Let those item types have their own shape.

## Decision Heuristics

- **"As a user" or a specific role?** Always specific. "As a free-tier subscriber", "As an on-call TAM", "As a first-time visitor". Generic "user" is a smell.
- **Story or epic?** If the story has > 1 sprint of work, > 5 acceptance criteria, or fails INVEST's `S`, it's an epic. Split via SPIDR.
- **Gherkin or checklist for AC?** Gherkin when the story has multi-step interactions, complex preconditions, or branching outcomes. Checklist when the criteria are independent observable facts ("filter persists", "empty state shown", "URL updates").
- **Put it in AC or in DoD?** Specific to this story → AC. Applies to every story → DoD. If you find yourself copy-pasting the same AC across stories, it belongs in DoD.
- **Spike or estimate?** If two senior engineers in the room give estimates more than 2 Fibonacci stops apart (3 vs 8), you don't have estimable consensus — spike.
- **Vertical or horizontal split?** Always vertical. If a split makes individual stories unshippable, you've split horizontally and need to re-cut.
- **Multiple roles in one story?** Usually a smell — split per role. Exception: a workflow where two roles interact and one slice involves both (e.g., "submitter creates a draft; approver receives it"). Even then, often two stories tied by a label.
- **Story or task?** If a real user is the beneficiary, story. If an engineer is the beneficiary, task. Tech debt can be either, depending on whether you can name the user-facing benefit.

## References

- [Mike Cohn — User Stories and User Story Examples (Mountain Goat Software)](https://www.mountaingoatsoftware.com/agile/user-stories) — the canonical statement of the "As a / I want / so that" template and the three essential elements.
- [Mike Cohn — User Stories Applied: For Agile Software Development (book PDF)](https://athena.ecs.csus.edu/~buckley/CSc191/User-Stories-Applied-Mike-Cohn.pdf) — long-form reference for stories, conversations, and acceptance.
- [Bill Wake — INVEST in Good Stories, and SMART Tasks (XP magazine, 2003)](https://xp123.com/articles/invest-in-good-stories-and-smart-tasks/) — original article coining INVEST.
- [Agile Alliance — INVEST glossary entry](https://agilealliance.org/glossary/invest/) — concise reference for each letter.
- [Mountain Goat Software — SPIDR: Five Simple but Powerful Ways to Split User Stories](https://www.mountaingoatsoftware.com/blog/five-simple-but-powerful-ways-to-split-user-stories) — Cohn's article on Spike / Path / Interface / Data / Rules.
- [Cucumber — Gherkin reference](https://cucumber.io/docs/gherkin/reference/) — canonical syntax for Given/When/Then, And, But, Background, Scenario Outline.
- [Scrum.org — Definition of Done vs Acceptance Criteria](https://www.scrum.org/resources/blog/what-difference-between-definition-done-and-acceptance-criteria) — official Scrum guidance on the distinction.
- [Parallel HQ — Given-When-Then Acceptance Criteria for Better User Stories](https://www.parallelhq.com/blog/given-when-then-acceptance-criteria) — practical patterns with examples.
- [Ron Jeffries — Essential XP: Card, Conversation, Confirmation](https://ronjeffries.com/xprog/articles/expcardconversationconfirmation/) — the 3 C's anatomy.
- [Mountain Goat Software — SPIDR Poster (PDF)](https://www.mountaingoatsoftware.com/uploads/blog/spidr-poster.pdf) — single-page splitting reference.
