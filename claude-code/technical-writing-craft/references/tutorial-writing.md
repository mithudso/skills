<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `tutorial-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: tutorial-writing
description: Learning-oriented documentation craft — the Diátaxis tutorial quadrant. Write scaffolded "follow me" lessons that guarantee success at every step, build a tiny working artifact, and leave the learner exit-with-a-completed-thing in hand. Covers narrator voice ("we'll create…"), no-detour discipline, Carpentries-style backward design, cognitive-load budgeting (7±2), confidence-by-layer construction, and the rule that learning is the goal, not the artifact. TRIGGER when the user asks to write a tutorial, "Getting Started", "Your first X", a quickstart, an onboarding lesson, a workshop module, a Carpentries-style lesson, or any doc whose primary purpose is to teach a newcomer by walking them through a complete worked example. Also TRIGGER when the user asks how a Diátaxis tutorial differs from a how-to, or how to scaffold a lesson. SKIP: problem-oriented step-by-step for a competent user (use howto-writing); ops/incident procedures with safety gates (use runbook-craft); API/spec lookup tables (use reference-doc-writing); concept/why/background discussion (use explanation-doc-writing); sentence-level prose craft on an existing tutorial (use technical-writing-craft); API endpoint or SDK reference page (use api-docs-craft).
---

# Tutorial Writing — Diátaxis Learning Quadrant

## Overview

A tutorial is **a lesson**. Its only job is to take a complete newcomer through a meaningful, hand-held experience and leave them with two things: a tiny working artifact they built themselves, and the confidence that they can use this tool. Tutorials are learning-oriented — they sit in the practical + studying quadrant of Diátaxis.

The single most-violated rule in tutorial writing: **the artifact does not matter; the learning does**. A child learning to cook is not there to produce a meal — they are there to be in a kitchen, holding utensils, hearing the words. A tutorial reader is not there to ship the thing they build — they are there to encounter the tool, the vocabulary, and the shape of the workflow under your protection.

If you find yourself optimizing for "they could use this output in production" — stop. You are writing a how-to. Switch quadrants.

## Core Concepts

### 1. The learner's promise

When a learner runs your step, they must see the result you said they would see. This is non-negotiable. Confidence is built layer by layer, and one broken step shakes the whole stack. Re-test every step on a clean environment before you ship.

### 2. Narrator voice — "we" not "you alone"

Tutorials use the first-person plural: *"we'll create a file called `app.py`"*, *"now we'll run it"*. The instructor is present. The learner is not abandoned to figure it out. Contrast with how-to voice ("Create a file named `app.py`") which assumes competence and leaves the reader on their own.

### 3. No detours

There will be a hundred interesting tangents — "by the way, you could also…", "in production you'd usually…", "the underlying mechanism is…". **Cut all of them.** A tutorial is not the place to teach options, alternatives, or theory. Every sentence either moves the learner toward the artifact or it leaves the document. Diversions belong in explanation docs; alternatives belong in how-to guides.

### 4. Exit-with-a-completed-artifact

The learner must finish with **something visible**: a running web server on `localhost:8000`, a printed "Hello, world", a deployed function that responded to a curl. The artifact is the receipt for the learning. Without it the tutorial dissolves into reading.

### 5. Cognitive load budget (7±2)

Human working memory holds roughly 7 ± 2 items. A tutorial step that introduces a new tool, a new file format, a new command-line flag, a new concept, and a new error mode — in one step — has already broken the learner. Each step introduces **one** new thing. Earlier-introduced things are reused, not re-explained.

### 6. Backward design (Carpentries)

Start at the end. Write down — in one sentence — what the learner can do after the tutorial that they could not do before. Then work backwards: what is the last step? What must they have done immediately before that? Keep regressing until you reach an empty machine. This is the spine. Now write the steps forward.

### 7. Concrete, particular, robust

Tutorials are built around **specific** actions and **specific** outcomes. Not "create a database" — *"create a database called `tutorial_db`"*. Not "you'll see some output" — *"you'll see exactly this output: `{...}`"*. The specificity is the safety rail; vagueness is where learners fall off.

### 8. The instructor's safety contract

You — the writer — are responsible for everything the learner encounters. If the install command on macOS 14 prompts for a password and the learner does not expect it, that is your error, not theirs. If a step depends on Python 3.11 and the learner has 3.9, you must say so in step 0. The learner is a guest in your kitchen; you don't let them touch the hot pan.

### 9. Inspire confidence, not competence

The goal is *"I can do this"*, not *"I have mastered this"*. Mastery is the job of explanation docs, reference docs, and time. The tutorial's emotional outcome is **agency**.

### 10. Tutorials are not the place for "if" or "depending on"

Branching kills tutorials. *"If you're on Windows, do X; if Mac, do Y; if Linux, do Z"* triples the cognitive load and triples the test surface. Pick **one** environment, declare it in step 0, and keep one linear path. Cover the other platforms in separate tutorials.

## Template — minimum viable tutorial

```markdown
# Your First <Thing>

In this tutorial we'll build a <small concrete thing>. By the end, you'll have
<artifact> running on <where>, and you'll have used <2-3 core concepts> for
the first time.

This tutorial takes about <N> minutes. We assume you have <one prerequisite>
installed.

## What we'll build

<one or two sentences and ideally a screenshot or sample output>

## Step 0 — Set up

<copy-pasteable env setup. Pin versions. Test on a clean machine.>

You should see:
```
<exact expected output>
```

## Step 1 — <one new thing>

<narrator voice: "Let's create a file called …">

```<lang>
<code>
```

Run it:

```
<command>
```

You should see:
```
<exact expected output>
```

<one sentence on what just happened — no theory, no alternatives>

## Step 2 — <one more new thing, building on step 1>

…

## Step N — <the final visible artifact>

…

Congratulations. You've built <artifact>. You can find the finished code at
<link>.

## What's next

- To do <related task>, see the [How to <verb>](…) guide.
- To understand *why* <concept> works this way, see [<concept> explained](…).
- For the full <thing> reference, see [API reference](…).
```

## Anti-Patterns

### AP-1 — The "tutorial" that is secretly reference

> "This tutorial covers the `Client` class, which has the following methods…"

That is reference docs in a costume. Tutorials walk a learner through *doing one thing*; they do not enumerate surface area.

### AP-2 — The "kitchen sink" tutorial

Twelve features demonstrated. Three languages. Two installation paths. Five optional sections. Every reader gets lost. **One concrete artifact. One environment. One path.**

### AP-3 — Theory before action

Three paragraphs of background before the first command. The newcomer's attention is finite and the artifact is the bait. Get them to `Hello, world` in the first five minutes, then offer concept context later (or push it to an explanation doc).

### AP-4 — Hand-waved steps

*"Now set up your database."* How? Which database? On what port? With what credentials? Every imperative must be **copy-pasteable**. If the learner has to make a decision, you have failed to scaffold.

### AP-5 — Untested steps

Tutorials decay. A bumped dependency, an OS update, a deprecated flag, and step 3 silently breaks. **Re-run the entire tutorial on a clean VM before each release.** This is not optional.

### AP-6 — "You'll see something like…"

Either it's exact or you've broken the learner's promise. Show the exact output. If the output is variable (a UUID, a timestamp), call that out explicitly: *"you'll see a UUID — yours will be different from ours, that's expected"*.

### AP-7 — Choose-your-own-adventure

Branches, conditional steps, optional add-ons, "intermediate readers can skip this" — all of it. A tutorial is a single line through the territory. Branches go in how-to guides.

## Decision Heuristics — is this actually a tutorial?

Ask these in order. If any answer is "no", you are not writing a tutorial.

1. **Is the reader a complete newcomer to the product?** If they could already articulate a specific goal ("I want to add OAuth to my app"), they need a **how-to**, not a tutorial.
2. **Is the artifact small enough that you can guarantee every step?** If you can't re-run the whole thing in 30 minutes on a clean machine, it's too big — split it.
3. **Is there exactly one path through?** If you find yourself writing "if/depending on/optionally", you are drifting toward how-to territory.
4. **Does the reader end with a visible, working thing they built themselves?** If they end with "now you understand…", you're writing **explanation**, not tutorial.
5. **Could you teach this to a friend in person, live?** If not, the scope is wrong.

When the answer to any of these is "no", consult the sibling skills:

- Competent reader with a goal → `howto-writing`
- Need to look up an API or schema → `reference-doc-writing`
- Need to convey background or "why" → `explanation-doc-writing`
- Production execution with safety gates → `runbook-craft`
- API/SDK reference page specifically → `api-docs-craft`

## Cross-pollination notes

- **technical-writing-craft** holds sentence-level prose craft and the Diátaxis theory overview — pair it with this skill when reviewing tutorial prose.
- **api-docs-craft** references this skill for the "Getting Started" / "Quickstart" section of API documentation.
- **writing-expert** owns audience calibration; consult it when deciding the assumed reader level.

## References

1. [Tutorials — Diátaxis](https://diataxis.fr/tutorials/) — Daniele Procida's canonical specification of the tutorial quadrant, including the learner-promise principle and the "lesson, not goal" distinction.
2. [Diátaxis — Start here](https://diataxis.fr/start-here/) — the five-minute overview of the four quadrants and the axes that separate them.
3. [Tutorials — Divio Documentation](https://docs.divio.com/documentation-system/tutorials/) — the earlier Divio formulation: tutorials must be concrete, robust, and built around specific outcomes; what a beginner needs to know vs. what an experienced user asks.
4. [Collaborative Lesson Development Training — Lesson Design (The Carpentries)](https://carpentries.github.io/lesson-development-training/lesson-design.html) — backward design, learning objectives, and the practice-feedback model.
5. [Overview of Carpentries pedagogic model](https://github.com/orchid00/The_Carpentries_info/blob/master/overview-of-carpentries-pedagogic-model.md) — live coding, scaffolding, cognitive-load budgeting (7 ± 2), and continuous-feedback pedagogy.
6. [What is Diátaxis and should you be using it](https://idratherbewriting.com/blog/what-is-diataxis-documentation-framework) — Tom Johnson's working-writer take on applying the four quadrants in practice.
