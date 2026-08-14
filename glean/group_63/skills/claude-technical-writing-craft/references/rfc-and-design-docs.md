<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `rfc-and-design-docs` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: rfc-and-design-docs
description: >
  RFC and design document craft — Stripe RFCs, Google design docs, Squarespace RFCs,
  ADRs (Architecture Decision Records), Tenets documents, and Premortem-style decision
  writing. Covers structure conventions, Goals/Non-Goals discipline, Alternatives
  Considered, Open Questions, Decision Log, ADR numbering and supersession, Amazon
  Tenets format, Premortem framing (Gary Klein), Stripe vs Google voice, one-pager vs
  multi-page scoping, two-pizza review, living document discipline, and common
  anti-patterns. TRIGGER: "write an RFC", "design doc", "ADR", "architecture decision
  record", "tenets document", "premortem", "alternatives considered", "decision
  document", "write a proposal", "write an architecture proposal", "decision record",
  "non-goals section", "goals and non-goals", "write a one-pager proposal".
  SKIP: code review (use code-reviewer); general docs (use writing-expert); runbooks
  (use technical-writing-craft + incident-response); incident post-mortems (use
  incident-response); sentence-level prose mechanics for RFC text (use
  technical-writing-craft — this skill covers RFC content and structure, not sentence
  style); customer-facing prose (use writing-expert + executive-comms).
origin: local
version: "1.2.0"
updated: "2026-05-29"
related_skills:
  - writing-expert
  - technical-writing-craft
  - executive-comms
  - document-critique
whenToUse:
  - "Write an RFC or design doc for a technical decision"
  - "Create an ADR (Architecture Decision Record) for a team"
  - "Draft a tenets document to guide team decisions"
  - "Apply premortem framing to a project plan or RFC"
  - "Structure the Alternatives Considered section of a proposal"
  - "Review an RFC draft for missing non-goals, buried decisions, or absent decision-by date"
  - "Choose between Stripe vs Google RFC voice for a document"
  - "Set up a two-pizza RFC review process"
  - "Write the Goals and Non-Goals section of a proposal"
  - "Add a Decision Log to an existing RFC"
---

# RFC and Design Document Craft

## When to use this skill

Activate when the user:

- asks to write, draft, or structure an RFC, design doc, or architecture proposal
- needs to write or review an ADR (Architecture Decision Record)
- is creating a tenets document for a team, product, or organization
- wants to apply premortem framing to a decision or project plan
- asks about Alternatives Considered, Goals/Non-Goals, or Open Questions discipline
- is scoping a review process (who reviews, how long, when to decide)
- wants to understand Stripe vs Google RFC voice and when to use each
- has an RFC draft that needs structural critique (missing non-goals, buried decisions, no decision-by date)

Skip this skill for: general technical writing (use `technical-writing-craft`); executive summaries and board memos (use `executive-comms`); runbooks and incident playbooks (use `incident-response`).

---

## How to respond when invoked

When a user asks to write or draft an RFC, ADR, tenets doc, or design doc:

1. **Clarify if ambiguous.** If the decision topic, affected teams, or document type are not clear, ask exactly one targeted question before proceeding. Example: "Is this a cross-team decision (full RFC) or a local architectural choice (ADR)?"
2. **Use the §10 template** as the starting skeleton for RFCs. Fill each section with the user's context. Mark any section where real data is missing with `[placeholder — replace with actual data]` rather than inventing numbers.
3. **Declare your voice choice.** State whether you are using Stripe voice (direct, opinionated) or Google voice (exhaustive, alternative-rich) and why, based on §5 guidance.
4. **Self-check before delivering.** Re-read the draft against the canonical section order in §1. Confirm Goals are falsifiable, Non-Goals are present, at least two real Alternatives are given, and a Decision-by date is in the header.
5. **Output shape:** deliver the filled document first, then a brief "Reviewer notes" block flagging any placeholders or sections that need real data from the user.

---

## 1. RFC Structure Conventions

A well-formed RFC has these sections, in order. Every section is mandatory except Shipping (full RFCs only) and Decision Log (required once the RFC is accepted).

**Canonical section order:** Title → TL;DR → Background → Goals → Non-Goals → Proposal → Alternatives Considered → Risks → Shipping → Open Questions → Decision Log

### Title
One-line. Describe what is being decided, not what is being discussed. Bad: "Caching Discussion." Good: "Use Redis as the distributed cache for session state."

### TL;DR
Three sentences maximum. State the problem, the chosen solution, and the key tradeoff. Readers who stop here should be able to vote.

### Background
What exists today? What is broken, costly, or limited? This section is factual, not argumentative. Do not embed the recommendation here. Cite specific numbers (latency, cost, error rate) whenever they exist.

### Goals
Concrete, verifiable outcomes. Each goal should be falsifiable — you should be able to look at the shipped system and confirm it was met. Prefer measurable outcomes over capability statements.

- Good: "Reduce P99 session-read latency from 120ms to under 20ms."
- Bad: "Improve performance."

Three to five goals is the upper limit. More than five is a scope problem.

### Non-Goals
**Non-Goals are often more useful than Goals.** They prevent scope creep, end review debates before they start, and give future readers a boundary map. A non-goal is a thing that a reasonable reader would assume is in scope — but is not.

- Good: "Multi-region replication is out of scope for this RFC; see RFC-0047."
- Bad: (omitting the section because everything seems obvious)

Write non-goals as affirmative sentences with explicit scope boundaries. Never write "we won't do X" without explaining what RFC or future work handles it, or why it is permanently out of scope.

### Proposal
The recommended approach. Be direct and opinionated — see §5 on Stripe vs Google voice. Include:

- The core mechanism (what gets built, changed, or removed)
- Data model or API shape sketches (enough to evaluate, not a full spec)
- Migration path if replacing something existing
- Dependency callouts (teams, systems, infra)

This section should convince a skeptical reader. Hedge sparingly and only when uncertainty is real.

### Alternatives Considered
Present at least two alternatives. Each must be a real option, not a strawman. For each:

1. State the option in one sentence
2. Explain why a reasonable engineer would choose it
3. Explain why it was rejected (cost, complexity, risk, incompatibility)

A strawman alternative ("we could do nothing") is acceptable only if doing nothing is a genuinely viable option that required deliberate rejection.

**The test:** if a reviewer says "why didn't you consider X?", X should appear in this section with a real answer. If it does not, the RFC is incomplete.

### Risks
Enumerate failure modes. For each risk:

- Likelihood: low / medium / high
- Impact: low / medium / high
- Mitigation: what reduces the likelihood or blast radius

Do not conflate risks with open questions. A risk is a known failure mode with a known owner. An open question is a thing that needs resolution before the risk can be assessed.

### Shipping *(full RFC only; omit from one-pagers)*
See §11 for required fields. Position: after Risks, before Open Questions.

### Open Questions
Questions that are parked for later — not buried in the Proposal section. Each open question should have:

- The question, precisely stated
- Who is responsible for resolving it
- The deadline or decision gate (e.g., "must be resolved before RFC is accepted" vs "can be deferred to implementation")

Open questions that block acceptance go at the top. Deferred questions go at the bottom.

### Decision Log
Append-only. Every entry records: the date, the decision made, who made it (or the forum), and a short rationale. Never edit a prior entry — mark it superseded and add a new entry.

```
2024-03-15 — Accepted. Reviewed by Platform Arch team. Redis Cluster chosen over
             single-node Redis; see rationale in Alternatives Considered §3.
2024-04-02 — Amended. Persistence mode changed from AOF to RDB after load test
             results (PR #1847). Decision by @eng-lead + @platform-arch.
```

---

## 2. ADR: Architecture Decision Records (Nygard 2011)

An ADR (Michael Nygard, "Documenting Architecture Decisions," 2011) captures a single architectural decision in a short, durable file.

### Required sections

**Title**: `ADR-NNNN: <imperative short description>`
Example: `ADR-0042: Use event sourcing for the audit log`

**Status**: One of:
- `Proposed` — under review
- `Accepted` — in force
- `Deprecated` — no longer recommended but not replaced
- `Superseded by ADR-NNNN` — replaced by a later decision

**Context**: The forces at play — technical, organizational, timeline. Neutral. Do not smuggle in the recommendation.

**Decision**: The active-voice statement of what was decided. One to three sentences.
- Good: "We will use PostgreSQL with JSONB columns for semi-structured audit events."
- Bad: "PostgreSQL was considered and it was decided that it would be used."

**Consequences**: What becomes easier, what becomes harder, what is now constrained. Include both positive and negative consequences. This is the section future architects will read when wondering why the system is the way it is.

### Numbering and supersession

- Number sequentially from 0001. Do not reuse numbers.
- Deprecated records stay in the repo — they are historical evidence.
- Supersession: set the old ADR status to `Superseded by ADR-NNNN`. In the new ADR, add `Supersedes ADR-MMMM` below the status line.
- Store ADRs in `docs/decisions/` or `adr/`. File name: `ADR-0042-use-event-sourcing-for-audit-log.md`.

### When to write an ADR vs an RFC

| Situation | Use |
|-----------|-----|
| Decision already made, document for posterity | ADR |
| Decision under active review, need stakeholder sign-off | RFC |
| Small, local architectural choice (one team, one service) | ADR |
| Cross-team or cross-system decision | RFC with Decision Log |
| Need to capture why an option was rejected | Both support Alternatives Considered |

---

## 3. Tenets (Amazon-style)

A tenets document establishes ordered priorities for a team, product, or principle. Tenets are used when a decision-maker is absent and the team needs to break ties.

### Format per tenet

```
N. <One-line tenet statement, phrased as a value claim.>

   <2–3 sentences: the reasoning behind this tenet, and the specific tradeoff
   it resolves. What does this tenet help you choose when two good options
   conflict? What would the world look like if this tenet were absent?>
```

### Rules for well-formed tenets

- **Order matters.** Tenet 1 beats Tenet 2 when they conflict. If two tenets never conflict, one is redundant.
- **Each tenet resolves a real tradeoff.** A tenet that everyone agrees with and that never forces a hard choice is a platitude, not a tenet.
- **The test** (Amazon-style): "Except when [opposite tenet]." A real tenet has a counterpart that would also be reasonable — and it wins precisely because you chose it over that counterpart.
- **Limit to 5–7 tenets.** More than 7 signals the document is a values list, not a decision framework.
- **Write in the affirmative.** Not "we do not sacrifice reliability" but "we treat reliability as a prerequisite, not a feature."

### Example (abbreviated)

```
1. Customer trust over developer convenience.

   When a UX decision benefits the operator but creates friction or opacity
   for the end customer, we choose the customer. Developer workflows can be
   improved iteratively; broken customer trust is expensive to rebuild.

2. Correctness over completeness.

   A partial result that is accurate is better than a full result that is
   approximate. We surface gaps explicitly rather than fill them silently.
```

---

## 4. Premortem Framing (Gary Klein, HBR 2007)

A premortem is a structured imagination exercise: assume the project has already failed catastrophically, then write the post-mortem in advance.

### When to use it

Use a premortem when:
- A decision has been made but not announced (surfacing last-minute objections in a psychologically safe way)
- A project plan exists but risks feel underspecified
- A team has strong social pressure toward optimism

### The structure

1. **Setup**: "It is [date 6–12 months out]. The project shipped and it was a disaster. The launch was rolled back / the migration corrupted data / the product was abandoned. Write the post-mortem."

2. **Individual generation** (5–10 minutes, silent): Each participant writes their top 3 failure causes independently. This prevents anchoring.

3. **Round-robin share**: Facilitator collects one cause per person per round, no discussion yet. Write on a shared surface.

4. **Prioritize**: Vote on the 3–5 most likely and most damaging failures.

5. **Convert to risk mitigations**: Each top failure becomes either a plan change, a monitoring addition, or an explicit accepted risk with a named owner.

### Writing a premortem in an RFC

Add a `Premortem` subsection inside Risks. Frame it: "If this RFC fails in production, the most likely causes are..." Then list the top 3 as if writing a retrospective from the future. This forces specificity that "Risks: medium likelihood" does not.

---

## 5. Stripe RFC Voice vs Google Design Doc Voice

These are distinct writing registers. Use the right one for the audience and culture.

### Stripe RFC voice: confident, opinionated, concrete

- State the recommendation in the first paragraph.
- Use first person plural ("We will migrate to…", "We reject X because…").
- Do not hedge unless uncertainty is real and material.
- Short sections. Dense paragraphs are acceptable. Tables preferred over prose lists.
- Tone: "Here is what we are doing and why. Convince us otherwise or consent."

### Google design doc voice: humble, exhaustive, alternative-rich

- Alternatives Considered is as long as the Proposal section.
- Open Questions are numerous and explicit.
- Use hedged language ("one approach is…", "we considered…").
- Long documents (10–20 pages) are acceptable.
- Tone: "We have done the analysis. Help us find what we missed."

### Choosing between them

Use Stripe voice when: decisions need to move fast, the author has authority, the audience is small (2–5 reviewers), or the decision is reversible.

Use Google voice when: the decision is irreversible, cross-team dependencies are high, the author needs to build consensus, or the system is large and the failure blast radius is wide.

Most engineering teams outside Stripe/Google benefit from a hybrid: Stripe's directness in the Proposal, Google's thoroughness in Alternatives Considered and Risks.

---

## 6. Scoping: One-Pager vs Multi-Page RFC

| Signal | One-pager (< 2 pages) | Full RFC (3–10+ pages) |
|--------|-----------------------|------------------------|
| Decision reversible in < 1 sprint | Yes | No |
| Single team affected | Yes | No |
| Alternatives are obvious | Yes | No |
| Irreversible or cross-cutting | No | Yes |
| Requires infra provisioning or data migration | No | Yes |
| Audit trail required for compliance | No | Yes |

When in doubt, start with a one-pager. If reviewers ask "but what about X?" three or more times, promote it to a full RFC.

---

## 7. Two-Pizza RFC Review

Coined informally at Amazon; the rule is that a review meeting should be small enough to feed with two pizzas (5–8 people).

### Who attends

- The RFC author (1 person or a pair)
- Decision-maker or DRI (1 person)
- Affected team leads (1–2 people)
- Domain expert who can catch technical errors (1 person)
- Optional: security, legal, or data reviewer if in scope

Do not invite everyone who might be curious. They can read the doc asynchronously.

### How to run it

1. Require async pre-read — no reading the RFC live.
2. Start by asking: "Does anyone have a blocking objection?" Surface hard stops first.
3. Time-box discussion of each section (Alternatives Considered gets the most time).
4. Close with an explicit decision or a precise list of what must change before re-review.
5. Record the outcome in the Decision Log within 24 hours.

### Decision-by date

Every RFC must have a `Decision by:` date field in the header. Without it, reviews drift. A good default is 5–10 business days from first distribution.

---

## 8. Living Document Discipline

An RFC is not frozen at acceptance. It is a living record.

### What to update after acceptance

- **Decision Log**: always append, never edit prior entries.
- **Status field**: change from `Proposed` to `Accepted`, `Deprecated`, or `Superseded`.
- **Open Questions**: mark each as Resolved with the answer and date when answered.
- **Risks**: mark mitigations as complete with a reference (PR number, postmortem link).

### What not to change

Do not retroactively edit the Proposal or Alternatives Considered to match what was actually built. The gap between the RFC and the implementation is historically valuable — it tells future readers what changed and why. Capture divergence in the Decision Log instead.

---

## 9. Anti-Patterns

| Anti-pattern | Why it fails | Fix |
|--------------|-------------|-----|
| Passive recommendations ("we could consider…") | Signals lack of conviction; invites endless discussion | State a recommendation; defend it |
| Missing Non-Goals | Scope creep, repeated review debates | Write non-goals before goals; show what is explicitly excluded |
| Decisions disguised as discussion | Reader cannot tell what was decided | Add a Decision Log; make the chosen option explicit in Proposal §1 |
| No decision-by date | Reviews drift indefinitely | Add `Decision by:` to the RFC header |
| Strawman alternatives | Looks like box-checking; erodes trust | Each alternative must be one a reasonable engineer would actually choose |
| Goals without metrics | Unverifiable; "success" is always debatable | Make each goal falsifiable (latency target, error budget, cost ceiling) |
| Open questions buried in Proposal text | Hard to track, never get answered | Extract to a dedicated Open Questions section with owners and deadlines |
| RFC updated to hide original proposal | Destroys audit trail | Append changes to Decision Log; leave original sections intact |

---

## 10. RFC Starter Template

Use this skeleton when asked to "write an RFC" from scratch. Fill in each section; delete sections marked "(one-pager: omit)" when producing a one-pager.

```markdown
# RFC: <Title — what is being decided, not discussed>

**Author:** <name>
**Decision by:** <date>
**Status:** Proposed

## TL;DR
<3 sentences: problem, solution, key tradeoff>

## Background
<What exists today. What is broken, costly, or limited. Specific numbers.>

## Goals
- <Falsifiable outcome 1>
- <Falsifiable outcome 2>

## Non-Goals
- <Thing a reader would assume is in scope — but is not. Reference the RFC or work that covers it.>

## Proposal
<Recommended approach. State it in the first paragraph. Include core mechanism,
data model sketch, migration path, and dependency callouts.>

## Alternatives Considered
### Option A: <name>
<Why a reasonable engineer would choose it. Why rejected.>

### Option B: <name>
<Why a reasonable engineer would choose it. Why rejected.>

## Risks
| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|-----------|
| <risk> | low/med/high | low/med/high | <mitigation> |

## Shipping  *(one-pager: omit)*
- **Timeline:** <milestones with dates>
- **Owner:** <DRI>
- **Rollout:** <dark launch / feature flag / staged / full cutover>
- **Rollback:** <how to undo; decision gate>
- **Success metrics:** <ties back to Goals>

## Open Questions
- [ ] <Question> — Owner: <name> — Must resolve before: <acceptance / implementation>

## Decision Log
*(append-only after acceptance)*
```

---

## 11. Shipping Section

Every full RFC should include a Shipping section between Risks and Open Questions (see canonical order in §1).

### Required fields

- **Timeline**: milestones with dates or sprint numbers
- **Owner**: the DRI (Directly Responsible Individual) for each milestone
- **Rollout strategy**: dark launch / feature flag / staged rollout / full cutover
- **Rollback plan**: how to undo the change if metrics degrade; what is the decision gate for rollback
- **Success metrics**: what signals confirm the RFC achieved its goals (ties back to Goals section)
- **Deprecation plan** (if replacing something): when the old path is removed and how migration is communicated

---

## 12. RFC 2119 / RFC 8174 normative keywords

When an RFC or design doc imposes requirements that another team, contributor, or
implementer must follow, the document should adopt the normative-keyword
vocabulary defined in IETF RFC 2119 (1997) and clarified by RFC 8174 (2017).
This is what gives "MUST" vs "SHOULD" their unambiguous, conformance-testable
meaning in standards work.

### The keywords and what they encode

| Keyword | Meaning |
|---------|---------|
| **MUST**, **REQUIRED**, **SHALL** | Absolute requirement. Non-conformance is a defect. |
| **MUST NOT**, **SHALL NOT** | Absolute prohibition. |
| **SHOULD**, **RECOMMENDED** | The behavior is strongly preferred; deviation requires the implementer to understand and document the consequence. |
| **SHOULD NOT**, **NOT RECOMMENDED** | The behavior is strongly discouraged; same justification burden as SHOULD. |
| **MAY**, **OPTIONAL** | Truly optional. Implementers can choose either behavior, and other implementations must remain interoperable with either choice. |

### RFC 8174 clarification

RFC 8174 amended RFC 2119 to fix a long-standing ambiguity: only the **uppercase**
forms carry the normative meaning. Lowercase "should" or "must" in the same
document is ordinary prose, not a conformance requirement. State this
explicitly in the RFC, typically in a "Conventions" or "Notation" section:

> The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
> "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" in this
> document are to be interpreted as described in RFC 2119 and RFC 8174 when,
> and only when, they appear in all capitals, as shown here.

Include this boilerplate verbatim in any RFC that ships conformance language to
external implementers.

### Worked example

Bad (semantics unclear — is "should" a requirement?):
> The client should retry on 503 responses and may give up after three attempts.

Good (RFC 2119 normative):
> The client **SHOULD** retry on 503 responses using exponential backoff.
> The client **MAY** stop retrying after three attempts.
> The client **MUST NOT** retry on 4xx responses other than 408 and 429.

A reviewer can now write a conformance test against each clause. Without
uppercase keywords, the same sentence reads as advice.

### When to break the rule

- **Internal-only ADRs and design docs** with a single implementing team often
  do not need RFC 2119 keywords — the team's review process resolves
  ambiguity. Use plain English: "We will retry on 503" is fine.
- **One-pagers** scoped to a single sprint rarely benefit from normative
  vocabulary; the overhead exceeds the value.
- **Cross-team RFCs, public protocol specs, API contracts, and SDK behavior
  documents** should use RFC 2119 keywords — these are the artifacts where
  silent disagreement on "should" vs "must" causes the most expensive bugs.

### References

- `https://www.rfc-editor.org/rfc/rfc2119` — Bradner, S. "Key words for use in
  RFCs to Indicate Requirement Levels." BCP 14, RFC 2119, March 1997.
- `https://www.rfc-editor.org/rfc/rfc8174` — Leiba, B. "Ambiguity of Uppercase
  vs Lowercase in RFC 2119 Key Words." BCP 14, RFC 8174, May 2017.

---

## 13. "Show your work" in design-doc reasoning

A design doc's job is not only to record the chosen design — it is to give
reviewers and future readers enough of the reasoning trail that they can
challenge or extend the conclusion. The phrase "show your work" comes from
math education (Polya, *How to Solve It*, 1945) and applies directly to
proposal writing: bare conclusions invite disagreement; visible reasoning
invites engagement.

### The rule

For every load-bearing claim in a Proposal or Alternatives Considered section,
make at least one of the following visible:

1. **The data the claim rests on** — the number, the source, the date.
2. **The constraint or principle that ranks the options** — cost cap, latency
   target, blast-radius limit, team tenet.
3. **The reasoning step from data + constraint to conclusion** — even a single
   sentence is enough if it names the inference.

### Worked example

Hidden reasoning (bad):
> We will use Redis Cluster.

Shown reasoning (good):
> We will use Redis Cluster. The session-cache working set is 240 GB
> (measured: 7-day average from staging, see Appendix A). Single-node Redis
> tops out at ~128 GB per instance on our standard SKU, so a sharded
> deployment is required. Redis Cluster is the operationally proven sharding
> mode in our infra (3 teams already running it); the alternative — Redis
> Sentinel with application-side sharding — would require new tooling we do
> not have.

The second version makes three inference steps visible: data → constraint →
alternative comparison. A reviewer who disagrees with the conclusion can now
challenge a specific step rather than just "voting no."

### When to break the rule

- **TL;DR sections** should be bare conclusions — that is their job. Show
  the work in the body, not in the summary.
- **Decisions already made and ratified** can be documented in the Decision Log
  without re-litigating reasoning. The original RFC body retains the work.
- **Reasoning chains longer than 4–5 steps** belong in an appendix or a linked
  analysis doc, not the body of a 5-page RFC.

### Anti-pattern: reverse-engineered justification

"Showing your work" is not writing a plausible-sounding chain after the
decision was made for other reasons (org politics, prior commitment, founder
preference). A reviewer can usually tell — and the Decision Log will record
the real reason eventually. Be honest about what actually drove the choice.

### References

- Polya, G. *How to Solve It: A New Aspect of Mathematical Method*. Princeton
  University Press, 1945. (Originating "show your work" pedagogy.)
- Industrial Empathy. "Design Docs at Google." `https://www.industrialempathy.com/posts/design-docs-at-google/` — emphasizes visible reasoning in Alternatives Considered.

---

## Sources

- Nygard, Michael. "Documenting Architecture Decisions." 2011. `https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions`
- Stripe RFC process. Referenced in public engineering blog posts and ex-Stripe engineer writeups (Bret Taylor, Will Larson "An Elegant Puzzle").
- Google. "Design Docs at Google." `https://www.industrialempathy.com/posts/design-docs-at-google/` (Louis Brandl, 2020); also Pluralsight "How to Write a Good Design Document" (2021).
- Klein, Gary. "Performing a Project Premortem." Harvard Business Review, September 2007.
- Bezos, Jeff. Amazon Shareholder Letters, 1997–2021 (Tenets / Leadership Principles framing).
- Larson, Will. "An Elegant Puzzle: Systems of Engineering Management." Stripe Press, 2019. (RFC process, one-pager vs full RFC scoping.)
- Bradner, S. "Key words for use in RFCs to Indicate Requirement Levels." RFC 2119, March 1997.
- Leiba, B. "Ambiguity of Uppercase vs Lowercase in RFC 2119 Key Words." RFC 8174, May 2017.
- Polya, G. *How to Solve It*. Princeton University Press, 1945.
