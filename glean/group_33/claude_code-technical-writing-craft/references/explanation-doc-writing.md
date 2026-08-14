<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `explanation-doc-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: explanation-doc-writing
description: Understanding-oriented documentation craft — the Diátaxis explanation quadrant. Write the "discussion" docs that build mental models, give background, explain why decisions were made, surface alternatives considered, and let the reader see the system in context. Covers the discussion-doc form, mental-model construction, history paragraphs, decisions-log patterns, "Why we chose X" pages, "How X works under the hood" pages, and the discipline that explanation is for understanding — not for proposing change (RFCs) and not for executing procedures (runbooks). TRIGGER when the user asks to write a concept doc, design rationale, "Why X" doc, architecture overview, mental-model doc, theory of operation, "How it works" doc, philosophy/principles page, ADR-style explanation, or any doc whose primary purpose is to leave the reader understanding rather than doing or looking up. SKIP: proposing a new design or decision (use rfc-and-design-docs); newcomer "build your first X" lesson (use tutorial-writing); goal-directed recipes (use howto-writing); API/CLI/schema lookup pages (use reference-doc-writing); on-call execution under pressure (use runbook-craft); sentence-level prose review (use technical-writing-craft).
---

# Explanation Doc Writing — Diátaxis Explanation Quadrant

## Overview

An explanation doc is **a discussion**. Its purpose is not to instruct, not to enumerate, and not to walk a reader through a goal. Its job is to leave the reader with a clearer **mental model** — of why the system is shaped this way, what alternatives existed, what trade-offs were made, and what conceptual landscape the user is now operating in.

Explanation docs are understanding-oriented and sit in the theoretical + studying quadrant of Diátaxis. The reader is not at work; the reader is reflecting. They closed the IDE, opened the doc, and are reading to **think**.

The single most-violated rule: **explanation is not a proposal**. It describes the world as it is — choices already made, designs already shipped, reasoning already settled. If you are arguing for a change, you are writing an **RFC**, not an explanation. If you are walking through a live incident response, you are writing a **runbook**.

## Core Concepts

### 1. The reader's question is "why", not "how" or "what"

*"Why does the cache invalidate on write rather than on read?"*
*"Why is the default replication factor 3?"*
*"Why do we use append-only logs here when most systems would update in place?"*

These are the prompts that bring readers to explanation docs. If the reader's question is "how do I configure replication factor", they need a how-to. If the question is "what is the syntax", they need reference. Explanation answers the **shape-of-the-world** questions.

### 2. Build the mental model, then layer the detail

Lead with the analogy or the diagram or the one-sentence essence. *"Think of the write-ahead log as a journal: every change is written there first, and only later applied to the data files."* Then add the next layer: *"This means a crash in the middle of an update never leaves the data file half-written — recovery just replays the journal."* Then the next. Readers leave the doc with a tree they can hang specifics on later.

### 3. Show alternatives and why they were not chosen

A decision without alternatives is not a decision; it's a proclamation. The explanation doc names the roads not taken: *"We considered an LSM-tree here, but the workload is read-heavy and the write-amplification penalty was not worth the write throughput gain."* This is the "decisions log" pattern — and it answers the question that future readers (including future maintainers) will ask in three years: *"why didn't we just do X?"*.

### 4. History earns trust

A short history paragraph — *"v1 used Redis for the queue; we hit head-of-line blocking under load in 2024 and moved to Kafka in v2"* — gives the reader context that no amount of current-state description can replicate. Systems carry their scars; explanation docs make those scars visible and legible.

### 5. Discuss, don't prescribe

Explanation uses words like *because*, *however*, *the trade-off is*, *one consequence is*. It avoids *do this*, *use this*, *configure this*. Prescriptions belong in how-tos. The explanation reader is forming opinions, not executing them.

### 6. Cite the boundary of the explanation

Explanation docs accumulate scope creep ferociously — every concept connects to every other. Name what's in and what's out: *"This doc covers our write-path; for read-path discussion see [the read-path doc]; for the on-the-wire protocol see [reference]"*. Without a stated boundary the doc becomes infinite.

### 7. Use multiple perspectives

Explanations benefit from showing the same thing from multiple angles: a diagram, a metaphor, a code-level example, a real-world consequence. Each angle catches a different reader. Reference can't do this (it would lose neutrality); tutorials can't (it would lose linearity); explanation thrives on it.

### 8. Distinguish from RFCs and ADRs

- **RFC** = proposing a change before it ships. Argumentative voice. Open to debate. Has a status (Draft / Accepted / Rejected). Belongs in `rfc-and-design-docs`.
- **ADR** (Architecture Decision Record) = a small, dated note of a decision made, the alternatives considered, and the rationale. ADRs are explanation **artifacts** — they fit naturally into the explanation quadrant as a structured genre.
- **Explanation doc** = the broader discussion of how and why the system is the way it is. Synthesizes many ADRs into a readable narrative. Stable, evergreen.

### 9. Distinguish from runbooks

A runbook says *"if disk is at 90%, page the on-call, run script X, escalate if Y."* That's execution under pressure. An explanation says *"we trigger paging at 90% disk because below that threshold a normal compaction cycle will recover space; above it the LSM merges back-pressure can stall writes."* That's the **why** behind the runbook. Both can exist; they are not the same doc.

### 10. Stay evergreen

Explanation docs should age slowly. If yours needs revision every sprint, you're documenting current implementation detail — that belongs in reference. Explanation captures the durable reasoning: invariants, trade-offs, philosophy. Re-read your explanation docs once a year and prune anything that has decayed into stale state.

## Templates

### Template — "How <X> works" / concept doc

```markdown
# How <X> works

<One-paragraph essence: the analogy or the one-sentence mental model.>

## Why <X> exists

<The problem <X> was designed to solve. Briefly — the reader will go to a how-to
or reference if they want operational detail.>

## The mental model

<2–4 paragraphs building the model. Diagram if useful. Layer from coarse to fine.>

## How it fits with the rest of the system

<Where <X> sits relative to neighboring components. What it depends on. What
depends on it.>

## Trade-offs and alternatives considered

We chose <X> over <Y> because <reason>. <Y> would have given us <advantage>
but cost us <disadvantage>, and at our workload shape the disadvantage
dominates.

We also looked at <Z>; <Z> is appropriate when <condition>, which is not
typically our condition.

## History

<Short paragraph: how <X> evolved. What we used before. What forced the change.>

## Where this discussion ends

This doc covers <scope>. For:

- *How to <task> with <X>* — see [the how-to](…).
- *The full <X> API* — see [reference](…).
- *A first walkthrough* — see [the tutorial](…).
```

### Template — ADR (lightweight decision record)

```markdown
# ADR-0042 — Use Kafka for the event bus

**Status:** Accepted
**Date:** 2025-11-12
**Deciders:** <team / individuals>

## Context

<What forced this decision. The problem state.>

## Decision

We will use Kafka as the primary event-bus transport between <system A> and
<system B>.

## Alternatives considered

- **Redis Streams** — simpler ops, but consumer-group semantics under partition
  failure did not meet our durability bar.
- **RabbitMQ** — strong routing model, but the throughput envelope did not
  match the projected load.
- **NATS JetStream** — promising, but our ops team had no Kafka-equivalent
  operational experience yet.

## Consequences

- We accept Kafka's operational complexity and pay it via the managed offering.
- Downstream consumers must adopt the at-least-once delivery model.
- Replays become possible (positive); cross-partition ordering is not
  guaranteed (negative).
```

## Anti-Patterns

### AP-1 — The "explanation" that is secretly a tutorial

> "To understand caching, let's build a simple cache. First, create a file…"

If the reader is creating files, you are running a tutorial. Cut the build-along and replace it with a diagram and prose. Explanation is a reading mode, not a typing mode.

### AP-2 — The "explanation" that is secretly reference

Twenty parameter tables and no narrative. The reader looking for the mental model finds a parts catalog. Move tables to the reference doc and link in.

### AP-3 — Prescriptive sneak-in

> "You should always set `replication=3` because…"

That's a how-to recommendation, not an explanation. Recast: *"Replication factor 3 trades disk and write-bandwidth for the ability to tolerate single-node failure; setups with weaker durability needs commonly use 1."* The reader can now decide; the doc has not crossed into instructing.

### AP-4 — No alternatives named

A doc that explains a choice without ever naming what was rejected reads as a sales pitch. The discussion mode requires both sides on the table — including the ones we chose against.

### AP-5 — Drift into RFC territory

Argumentative voice, open questions, "we are considering moving to…". Stop. That's an RFC. Move it to `rfc-and-design-docs` with proper status tracking. Explanation describes the world that exists; RFC argues for a world that might.

### AP-6 — Infinite scope

The doc started on the cache layer and is now also covering serialization, the wire protocol, the deployment topology, and the team's hiring philosophy. Name the boundary in the first paragraph and link out for everything else.

### AP-7 — No diagram where one would carry the model

Some concepts (data flow, state machines, partition layouts, dependency graphs) collapse into prose painfully. A 200-pixel diagram replaces 400 words and is faster to skim. Include them.

## Decision Heuristics — is this actually explanation?

1. **Is the reader's question "why" or "how does this fit"?** If yes → explanation. If "how do I do X" → how-to. If "what is the parameter list" → reference.
2. **Is the reader reading to think, or reading to act?** Thinking-mode → explanation. Acting-mode → tutorial / how-to / runbook.
3. **Is the decision already made and shipped?** If yes → explanation. If you're proposing the decision → `rfc-and-design-docs`.
4. **Does the doc include alternatives considered and trade-offs named?** If not, it's likely a marketing page or a one-sided proclamation, not an explanation.
5. **Does the doc age slowly?** If you'd need to rewrite half of it next month, you are documenting implementation detail (→ reference) instead of durable reasoning.

When it's not explanation, switch quadrants:

- Newcomer hands-on lesson → `technical-writing-craft` (references/tutorial-writing.md)
- Goal-directed recipe → `technical-writing-craft` (references/howto-writing.md)
- Surface-area lookup → `technical-writing-craft` (references/reference-doc-writing.md)
- Production execution under pressure → `technical-writing-craft` (references/runbook-craft.md)
- Proposing a design / asking for review → `rfc-and-design-docs`
- Sentence-level prose review → `technical-writing-craft`

## Cross-pollination notes

- **rfc-and-design-docs** is the sibling for *proposing* a design; explanation describes a shipped design. The two often pair: ship the RFC, then distill an explanation doc from it for newcomers.
- **technical-writing-craft** owns prose-level review and the Diátaxis theory primer.
- **postmortem-writing** is a specialized explanation genre — explanation of a past incident, focused on causal chains and learnings.
- **software-architect** and **rfc-and-design-docs** both produce inputs that feed the explanation layer.
- **mongodb-expert** and similar concept-heavy skills implicitly produce explanation-style content — link this skill from their authoring workflows.

## References

1. [Explanation — Diátaxis](https://diataxis.fr/explanation/) — Procida's canonical specification: explanation is discussion-oriented, deepens understanding, and sits outside the user's immediate work.
2. [Diátaxis — Start here](https://diataxis.fr/start-here/) — the theoretical+studying quadrant placement and the "reflecting reader" model.
3. [Explanation — Divio Documentation](https://docs.divio.com/documentation-system/explanation/) — Divio's "discussion" formulation, including the multiple-perspectives principle and the alternative-considered convention.
4. [Decision process and principles — The Rust Language Design Team](https://lang-team.rust-lang.org/decision_process.html) — how a mature open-source project records reasoning and "dissents", a working model for the decisions-log pattern.
5. [Architecture documentation — Rust Project Primer](https://rustprojectprimer.com/documentation/architecture.html) — practitioner take on documenting *why* the architecture is as it is, not only *what* it is.
6. [Documentation Quadrants — The Grand Unified Theory of Documentation (Dunn)](https://dunnhq.com/posts/2023/documentation-quadrants/) — modern practitioner walkthrough including the explanation-vs-reference separation pitfalls.
