<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `howto-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: howto-writing
description: Problem-oriented step-by-step documentation craft — the Diátaxis how-to quadrant. Write goal-directed recipes for competent users who already know the product and need "tell me how to do X" without learning scaffolding. Covers "How to <verb> a <noun>" titles, prerequisite blocks, single-goal-per-document discipline, numbered steps with imperative verbs, cookbook conventions, and the rule that how-tos assume competence (no teaching, no theory). TRIGGER when the user asks to write a how-to guide, recipe, cookbook entry, task guide, "How to X" doc, or any doc whose audience is a user with a specific goal who is not a newcomer. Also TRIGGER when the user is unsure whether to write a tutorial or how-to, or asks how to title or scope a task guide. SKIP: newcomer-onboarding lessons (use tutorial-writing); ops/incident execution with safety gates and rollback (use runbook-craft); API/SDK lookup pages (use reference-doc-writing); background/concept/"why" docs (use explanation-doc-writing); sentence-level prose review (use technical-writing-craft); proposing a design or decision (use rfc-and-design-docs).
---

# How-To Writing — Diátaxis How-To Quadrant

## Overview

A how-to guide is **a recipe**. It serves a competent user who has arrived at the page with a specific goal already formed in their head — *"How do I add OAuth to my app?"*, *"How do I back up this database?"*, *"How do I deploy to staging?"* — and who needs an efficient series of steps to get from where they are to where they want to be.

How-tos are problem-oriented and sit in the practical + working quadrant of Diátaxis. They are not lessons. They do not teach. They do not coddle. They assume the reader knows the vocabulary, owns the tools, and only needs the **sequence**.

The single most-violated rule: **how-tos answer a question only a competent user could ask**. If your reader cannot even formulate the goal, they need a tutorial — write that instead.

## Core Concepts

### 1. Title starts with "How to <verb>"

The title is the goal. Concrete verb, concrete object. *"How to rotate the signing key"*. *"How to add a custom domain"*. *"How to migrate from v1 to v2"*. If you cannot fit your doc's purpose into that pattern, you do not have a how-to — you have something else (probably an explanation or a tutorial in disguise).

### 2. One goal per document

A how-to with three goals is three how-tos that have not yet been separated. Resist the merge. Each goal gets its own URL, its own title, its own search hit. "How to deploy and roll back and configure secrets" is a category, not a document.

### 3. Prerequisite block up front

The first block under the title says what the reader must already have, know, or have done. This is the assumption surface. *"You have a v2 cluster running. You have admin credentials. You have `mongosh` installed."* If the reader can't tick every box, they bounce — which is correct, because the doc is not for them yet.

### 4. Numbered, imperative steps

Steps are numbered. Each step begins with an imperative verb. *"Create…"*, *"Run…"*, *"Set…"*, *"Verify…"*. Not *"Now you might want to…"*, not *"Let's…"*. The voice is direct because the reader is direct. Compare with tutorial narrator-voice (*"we'll create…"*) — how-tos drop the chaperone.

### 5. Assume competence, omit teaching

A how-to does **not** explain what a Kubernetes namespace is. It does **not** define OAuth. It does **not** walk the reader through what JSON looks like. If they need that knowledge, they're at the wrong door. Link out to an explanation or reference doc if useful; do not inline the lesson.

### 6. Branches are allowed (unlike tutorials)

How-tos can — and often must — branch. *"If you use Atlas, run X; if you self-host, run Y."* *"For production, also do Z."* The competent reader can hold the fork. Keep branches shallow (1–2 levels) and label them clearly with a marker like **If…** or **For self-hosted clusters:**.

### 7. Verifications, not promises

Tutorials promise *"you'll see exactly this"*. How-tos verify: *"Confirm that `kubectl get pods` shows the new pod in `Running` state."* The reader is checking their own work, not following a guided demo. Verifications come at decision points, not after every line.

### 8. The "how to choose" meta-guide is sometimes the right doc

Sometimes the reader's real question is not "how do I do X" but "which of A/B/C do I want?". A short decision matrix at the top of a how-to (or a dedicated "How to choose between X and Y" page) is legitimate and often more useful than three separate how-tos.

### 9. Failure modes are part of the recipe

When a step can plausibly fail, name the failure inline. *"If you see `Permission denied`, your role does not have `clusterAdmin` — see [granting roles](…)."* This is the cookbook convention: the user is at the stove, the food is burning, and the recipe is the only friend they have right now.

### 10. End at the goal, not past it

When the reader has done the thing, stop. Do not add *"now you might also want to…"* — that's a related-links section, not a step. The how-to is done when the goal is done. Optional follow-up belongs in a **See also** block, not in the numbered list.

## Template — minimum viable how-to

```markdown
# How to <verb> <noun>

<One sentence stating the outcome and when you'd want this.>

## Before you start

- You have <prerequisite 1>.
- You have <prerequisite 2>.
- You have <permission / role / credential>.

If you don't yet have <X>, see [the relevant doc].

## Steps

1. <Imperative verb>… <action>.

   ```
   <command>
   ```

2. <Imperative verb>… <action>.

   <If a decision is needed:>
   - **If <condition A>**, <do X>.
   - **If <condition B>**, <do Y>.

3. Verify the result:

   ```
   <verification command>
   ```

   You should see <observable indicator>.

## If something goes wrong

- **<Error 1>** — <cause and fix or link>.
- **<Error 2>** — <cause and fix or link>.

## See also

- [Related how-to]
- [Background / explanation doc]
- [API reference for <X>]
```

## Anti-Patterns

### AP-1 — Teaching inside the how-to

> "Before we deploy, let's understand what a deployment is…"

Stop. The reader knows what a deployment is — that's why they're here. If they don't, write or link an explanation doc. Inlining lessons triples the length and insults the competent reader.

### AP-2 — Multi-goal mega-guide

> "How to configure, deploy, and monitor your service"

That is three guides. Split them. A guide that needs five top-level sections to cover the goal is hiding three siblings.

### AP-3 — Tutorial drift

Narrator voice ("we'll now create…"), promise language ("you'll see exactly…"), copy-pasteable hand-holding for every micro-step. This is a tutorial wearing a how-to label. Either commit to the tutorial form (and accept the constraints) or trust the reader.

### AP-4 — Untitled or vague title

> "Database setup"

Setup for what? In what context? "How to set up MongoDB Atlas as a read replica" is the title. The reader is searching for that exact phrase.

### AP-5 — Missing prerequisites

The reader hits step 3 and discovers they needed `kubectl` configured against a specific cluster. They rage-quit. The prerequisite block is the contract — list everything the reader must have *before* they start typing.

### AP-6 — Steps that branch into theory

Step 5 is two sentences. Step 6 is three paragraphs explaining the algorithm being invoked. The asymmetry is a smell — that paragraph is an explanation doc trying to escape. Link to it, don't inline it.

### AP-7 — No verification

Twelve steps, no checks. The reader does all twelve, then discovers at the end that step 4 silently failed. Insert verification at every decision point and at the end. The cost is two lines per check; the savings are an hour of debugging.

## Decision Heuristics — is this actually a how-to?

1. **Could the reader phrase their question as "How do I <verb>?" before reading?** If they couldn't even ask, they need a tutorial.
2. **Is there exactly one goal?** If not, split.
3. **Does the reader already own the vocabulary?** If you need to define basic terms, you are drifting into tutorial or explanation.
4. **Is the goal a daily task, or a once-ever event?** Once-ever production-critical events with rollback procedures are **runbooks**, not how-tos. Use `runbook-craft`.
5. **Is the answer "click these buttons"?** That's still a how-to — UI-driven recipes are a valid form.

When the answer points elsewhere, consult the siblings:

- Newcomer who can't yet phrase the goal → `tutorial-writing`
- "Why does this exist / what is this for" → `explanation-doc-writing`
- "What are the parameters of `foo()`" → `reference-doc-writing`
- "How do I respond to this incident at 3am" → `runbook-craft`
- "Should we adopt X" → `rfc-and-design-docs`

## Cross-pollination notes

- **runbook-craft** is a specialized how-to genre with safety gates, rollback steps, on-call context, and escalation paths. Default to runbook-craft for ops-critical procedures.
- **api-docs-craft** uses how-to patterns for the "Guides" or "Use cases" section of API documentation.
- **technical-writing-craft** owns prose-level review and the Diátaxis theory primer.
- **support-ticket-writing** borrows the imperative-step structure for customer-facing reproduction steps.

## References

1. [How-to guides — Diátaxis](https://diataxis.fr/how-to-guides/) — Procida's canonical definition: how-to guides serve users at work, who already have context and a goal in mind.
2. [Diátaxis — Start here](https://diataxis.fr/start-here/) — the working-vs-studying and practical-vs-theoretical axes that place how-to in its quadrant.
3. [How-to guides — Divio Documentation](https://docs.divio.com/documentation-system/how-to-guides/) — Divio's earlier formulation: a how-to is an answer to a question only a user with some experience could even formulate.
4. [A technical guide to the Diátaxis framework for modern documentation (Ekline)](https://ekline.io/blog/a-technical-guide-to-the-diataxis-framework-for-modern-documentation) — modern application notes, including the how-to vs tutorial separation in real-world doc sets.
5. [Documentation Quadrants — The Grand Unified Theory of Documentation (Dunn)](https://dunnhq.com/posts/2023/documentation-quadrants/) — a practitioner's walk-through of when to pick how-to over tutorial vs reference.
