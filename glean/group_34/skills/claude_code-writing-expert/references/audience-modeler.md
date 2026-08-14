<!-- hub-reference-banner -->
> **Reference file — part of the `writing-expert` hub.** Formerly the standalone `audience-modeler` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: audience-modeler
description: >
  Infer and profile the audience for any document or communication, then re-evaluate
  the draft against that profile. Produces a structured audience profile table covering
  knowledge level, jargon tolerance, attention budget, motivations, objections, power/role,
  and emotional state. Scores the draft on vocabulary, length, motivation alignment, and
  objection coverage.
triggers:
  - "who is the audience for this"
  - "audience profile"
  - "model the reader"
  - "what does my exec want to see"
  - "audience analysis before writing"
  - "audience-modeler"
  - "who should I be writing this for"
  - "what level of detail does my audience need"
  - "is this too technical for my audience"
  - "reader persona for this"
  - "I need to understand my readers before writing"
  - "who reads this doc"
  - "tailor this for my audience"
skip: drafting (use writing-expert); critiquing (use document-critique); pure persona research (use deep-research)
related: writing-expert, document-critique, executive-comms, draft-review-revise-loop
version: 1.2.1
updated: 2026-05-29
---

# Audience Modeler

Meta-skill: given a document or topic, infer the audience, build a structured profile, then re-evaluate the doc against that profile. Every dimension leads to a concrete adjustment.

---

## Minimum deliverable per invocation

A complete invocation produces all of the following, in order:

1. **8-question answers** — work through the interrogation silently; surface any answer that reveals a significant inference or risk.
2. **Audience profile table** — all 10 rows filled, including the "Implication for this doc" column.
3. **Re-evaluation score** — if a document was provided, score all four dimensions and produce a finding for every dimension below 4.
4. **Lazy-generic check** — one sentence confirming the reader description is specific enough, or flagging a generic label.

If ambiguous input prevents filling a dimension, ask exactly one targeted question before proceeding. Do not produce a partial profile without flagging which values are assumed.

## Input handling

| Input provided | Action |
|---|---|
| Full document | Infer all dimensions from content; mark inferred values with `*` and note them below the table |
| Topic or description only | Infer from context; flag every assumed value; invite correction |
| URL, code file, slide deck | Ask the user to paste the text content, or describe the target reader explicitly |
| Nothing (no document, no topic) | Ask: "What document or topic should I profile the audience for?" |

---

## When not to use

- The user wants a draft written — use **writing-expert** (compose after profiling if needed).
- The user wants prose or structure critiqued — use **document-critique** (pass the audience profile as evaluation frame).
- The user wants market research or demographic data about a buyer persona — use **deep-research**.
- The user asks "is this too long?" without a document to profile — route to **document-critique**.

---

## Why audience modeling before writing

The Heath brothers ("Made to Stick") named this the **curse of knowledge**: once you know something, you cannot unknow it, so you write for people who already share your mental model. Hugh Rank ("The Pitch") frames persuasion as a matching problem — your message must meet the audience's existing hopes and fears, not the ones you wish they had. Aristotle's three appeals (ethos, logos, pathos) only work when calibrated to who the audience trusts, what they find logical, and what moves them emotionally. Get the audience wrong and all three collapse.

---

## The 8-question audience interrogation

Work through these before filling the profile table. If you cannot answer a question, that gap is a risk — surface it.

1. **Who reads this?** — Job title, team, organizational role, relationship to your work.
2. **What do they already know?** — Their domain background, familiarity with your project, vocabulary they use day-to-day.
3. **What do they want to do after reading?** — The outcome they seek: make a decision, understand something, be convinced to act, stay informed.
4. **What would make them say no?** — Predictable objections, prior failed pitches, institutional skepticism, competing priorities.
5. **How much time will they spend?** — Realistic reading time, not aspirational. Calibrate length accordingly.
6. **What's their power level — can they act on this?** — Decision-maker, influencer, executor, or informed party. Affects how much you ask them to do.
7. **What's their emotional state when reading?** — Anxious about a deadline, curious, skeptical, time-pressed, or actively hostile to your proposal.
8. **What else is competing for their attention right now?** — Their inbox load, their active crises, the other docs they're reviewing. This sets the real attention budget.

Surface any answer where the inference is uncertain or where the gap changes what the document should say. Silent reasoning is fine for clear cases; surface uncertain ones.

---

## Audience modeling dimensions

### 1. Knowledge level

| Level | Meaning | Writing implication |
|---|---|---|
| Novice | First encounter with domain | Define every term; use analogies; avoid jargon entirely |
| Intermediate | Working familiarity, not mastery | Light definitions on first use; acronyms OK if spelled out |
| Expert | Practitioner; daily use | No definitions; skip background; go straight to the problem |
| Domain specialist | Research-depth command | Match their vocabulary exactly; precision over clarity |

### 2. Jargon tolerance

- **High** — They use the jargon themselves; not using it signals you don't belong.
- **Medium** — Comfortable with industry terms; loses patience with internal abbreviations.
- **Low** — Cross-functional reader; jargon triggers confusion and distrust.
- **Zero** — External audience (customer, regulator, press); translate everything.

### 3. Attention budget

| Budget | Realistic read | Structural implication |
|---|---|---|
| 90 seconds | Exec in meeting prep | One sentence TL;DR, three bullets max, no appendices |
| 5 minutes | Manager scanning before sync | Headline + summary paragraph + key asks visible above fold |
| 30 minutes | IC or analyst doing their job | Full doc with sections; moderate depth OK |
| Hours | Auditor, regulator, technical reviewer | Footnotes, references, exhaustive evidence acceptable |

### 4. Motivations

What they want from this document:

- **Make a decision** — Give them a clear recommendation and the minimum evidence to trust it.
- **Understand** — Explain the model, not just the conclusion.
- **Be convinced** — Address their prior skepticism directly; don't just assert.
- **Act** — Precise steps; no ambiguity about who does what by when.

One document can have multiple motivations, but they should be ordered. State the primary motivation explicitly in the opening paragraph.

### 5. Objections

Pre-mortems on your own document. Common categories:

- "This isn't my problem to solve" — clarity of ownership needed.
- "We tried this before" — acknowledge prior history; explain what's different.
- "This costs too much / takes too long" — quantify and compare to cost of inaction.
- "I don't trust the data" — cite sources; show methodology.
- "This contradicts what I was told by X" — address the conflict head-on.

List objections explicitly in your audience profile. Each one should appear somewhere in the doc, even if only briefly.

### 6. Power / role

| Role | What they do with this doc | Implication |
|---|---|---|
| Decision-maker | Approve, reject, fund, block | Lead with recommendation; defer evidence to appendix |
| Influencer | Shapes the decision-maker's view | Give them language they can repeat; help them represent you |
| Executor | Implements the outcome | Precise specs, clear acceptance criteria, escalation paths |
| Informed party | Needs awareness, not action | Short; passive reading expected; no asks |

### 7. Emotional state

- **Anxious** — Reassure first; lead with stability and control.
- **Curious** — Reward them early; deliver the insight they're searching for.
- **Skeptical** — Don't oversell; use evidence and understatement.
- **Time-pressed** — Front-load everything; no narrative warmup.
- **Hostile** — Acknowledge their concern before making your case; Aristotle's ethos before logos.

---

## Audience profile output template

Produce this table on every invocation. Fill the "Implication for this doc" column — that column is the deliverable. A table without implications is classification, not analysis.

Mark inferred values with `*`. If a value cannot be inferred, write `unknown — ask` and stop before producing the re-evaluation score.

| Dimension | Value | Implication for this doc |
|---|---|---|
| Primary audience | [Job title / role] | |
| Secondary audience | [Job title / role, or "none"] | |
| Knowledge level | novice / intermediate / expert / domain-specialist | |
| Jargon tolerance | high / medium / low / zero | |
| Attention budget | 90s / 5m / 30m / hours | |
| Primary motivation | make decision / understand / be convinced / act | |
| Key objections | [list 2–4] | |
| Power / role | decision-maker / influencer / executor / informed party | |
| Emotional state | anxious / curious / skeptical / time-pressed / hostile | |
| Competing attention | [active pressures on this reader right now] | |

---

## Multi-audience documents

Most business documents have a primary and a secondary audience. Write for the primary; accommodate the secondary.

- The **primary audience** determines vocabulary, length, and tone.
- The **secondary audience** is served by structure: section headers that let them skip to relevant parts, a glossary in the appendix, a separate executive summary.
- When primary and secondary audiences have conflicting attention budgets, use progressive disclosure: short up front, depth in subsections or appendices.
- If you cannot identify a primary audience, the document has a structural problem, not just a writing problem. Clarify purpose before proceeding.

---

## Persona archetypes for common business contexts

### The busy executive

- Knowledge: broad, shallow on your domain.
- Attention: 90 seconds to 2 minutes, reading between meetings.
- Motivation: make a decision or stay informed.
- Emotional state: time-pressed, mildly skeptical of requests.
- Writing implication: TL;DR sentence in the subject line or first line. Three things max: what happened, what it means, what you need from them. No background they didn't ask for.

### The skeptical engineer

- Knowledge: deep, high jargon tolerance, strong pattern-matching for BS.
- Attention: 30 minutes if the problem interests them; 2 minutes if it doesn't.
- Motivation: understand and evaluate; will implement if convinced.
- Emotional state: skeptical by training; prizes precision over enthusiasm.
- Writing implication: lead with the problem statement, not the solution. Show your constraints. Acknowledge tradeoffs you rejected. Never oversell.

### The time-pressed customer

- Knowledge: varies; assume none about your internal systems.
- Attention: 60–90 seconds; mobile context likely.
- Motivation: act or decide quickly.
- Emotional state: anxious or frustrated; reading because something went wrong or because they need something now.
- Writing implication: one clear action per message. Zero jargon. Confirm that they matter more than your explanation does.

### The new-hire onboarding

- Knowledge: novice on your systems, possibly expert in their domain.
- Attention: 30 minutes to hours; motivated to absorb.
- Motivation: understand; build mental model.
- Emotional state: curious but overwhelmed by volume.
- Writing implication: explicit structure and wayfinding. Define every internal term. Show how concepts connect. Don't assume orientation docs were read.

### The auditor or regulator

- Knowledge: process and compliance expert; high domain literacy, low system literacy.
- Attention: hours; will re-read.
- Motivation: verify claims; identify gaps and exceptions.
- Emotional state: neutral-to-skeptical; not adversarial by default unless you give them reason.
- Writing implication: precision over persuasion. Cite every claim. State explicitly what is in scope and what is not. Leave no inference gaps — they will fill them with the worst case.

---

## Re-evaluation pass

Given a draft and an audience profile, score on four dimensions. Use 1–5 where 1 = fails entirely, 5 = optimized.

| Dimension | Score (1–5) | Finding |
|---|---|---|
| Appropriate vocabulary | | Jargon mismatches; terms not defined; wrong register for knowledge level |
| Appropriate length | | Too long or too short for stated attention budget |
| Addresses motivations | | Does the doc give the audience what they came for? |
| Anticipates objections | | Are the predictable objections surfaced and answered? |

For every dimension scored below 4: name the specific passage or pattern that fails, then state the change needed. "Section X uses jargon Y; define it or replace with Z" — not "vocabulary could be clearer."

After scoring, re-read the output and confirm it addresses the goal stated in "Addresses motivations." If it does not, revise before delivering.

---

## Worked example

**Input:** A 600-word engineering proposal recommending migration from MongoDB Atlas Serverless to Atlas Flex, addressed to "the team."

**Step 1 — 8-question reasoning (surfaced inferences):**
- Q1: "the team" is generic — inferring primary audience as the engineering manager who owns the decision, secondary as the engineers who implement it. `*`
- Q4: Likely objection: "Flex has different cost behavior — do we understand it?" Surfaced as key objection.
- Q7: Probably skeptical — this is a reversal of a prior choice.

**Step 2 — Profile table:**

| Dimension | Value | Implication for this doc |
|---|---|---|
| Primary audience | Engineering manager `*` | Lead with the decision and its cost impact; don't bury the recommendation |
| Secondary audience | Senior engineers `*` | Include a technical rationale section they can validate |
| Knowledge level | Expert | No definitions needed; use Atlas product names directly |
| Jargon tolerance | High | Atlas-specific terms (Flex, Serverless, vCPU hours) are fine |
| Attention budget | 5 minutes | Recommendation + rationale in the first 200 words; detail below |
| Primary motivation | make a decision | State the recommendation in sentence 1 |
| Key objections | "We just moved to Serverless"; "Cost model is unclear" | Acknowledge the prior decision; include a cost comparison table |
| Power / role | Decision-maker `*` | Ask for an explicit approval; don't leave the ask implicit |
| Emotional state | Skeptical `*` | Lead with data, not enthusiasm; show the constraint that forces the change |
| Competing attention | Sprint planning, Q3 roadmap review `*` | Keep it under 600 words; use headers for scanning |

**Step 3 — Re-evaluation (4/5 vocabulary; 2/5 length — the proposal buries the recommendation at paragraph 4; 4/5 motivations; 3/5 objections — Serverless decision not acknowledged):**
- Length: move recommendation to sentence 1 of the abstract.
- Objections: add one paragraph acknowledging "we chose Serverless 8 months ago because X; that constraint has changed because Y."

---

## The "lazy generic" test

Before finalizing a profile, check:

- Did you write "engineers" when you mean senior backend engineers at a Series B startup with a Postgres legacy?
- Did you write "customers" when you mean enterprise procurement managers who have never touched the product?
- Did you write "users" when you mean shift supervisors who read on a tablet between tasks?

Generic audience labels produce generic documents. The test passes when you can describe the reader in one specific sentence that no one else would write.

---

## Anti-patterns

**Writing for the audience you wish you had.** The doc assumes technical depth the actual reader does not have, or assumes enthusiasm for the topic the reader has not yet formed. Common when the author is expert and forgets the reader is not.

**The audience-of-one trap.** Optimizing entirely for a single powerful stakeholder (a specific VP, a key customer) at the expense of the broader audience who also reads and acts on the document. The one person you're writing for may never read it; the fifteen people you ignored will.

**Assuming uniform expertise within an audience.** "The engineering team" has staff engineers, junior engineers, and a product manager in the distribution. "Our customers" span Fortune 500 procurement and solo startup founders. Pick the primary reader and acknowledge the spread in your profile.

**Conflating role with knowledge.** A VP of Sales may be a domain expert in deals but a novice on infrastructure security. Title does not determine knowledge level.

---

## Compose with adjacent skills

- **writing-expert** — use audience-modeler first to build the profile, then invoke writing-expert to draft against it.
- **document-critique** — use audience-modeler to build the profile, then invoke document-critique with the profile as the evaluation frame.
- **executive-comms** — when the primary audience is an executive, combine with executive-comms for tone and format conventions.
- **draft-review-revise-loop** — embed the re-evaluation pass as the review criterion in each loop iteration.

---

## Sources

- Hugh Rank, "The Pitch" (1982) — persuasion as audience-matching; hope/fear targeting.
- Nielsen Norman Group persona research — behavioral dimensions over demographic; goals and mental models over job titles.
- Aristotle, *Rhetoric* — ethos, logos, pathos as audience-calibrated appeals, not fixed techniques.
- Chip Heath and Dan Heath, "Made to Stick" (2007) — curse of knowledge as the core failure mode of expert-to-novice communication.
