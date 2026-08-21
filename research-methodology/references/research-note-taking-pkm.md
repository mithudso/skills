<!-- hub-reference-banner -->
> **Reference file — part of the `research-methodology` hub.** Formerly the standalone `research-note-taking-pkm` skill.
> Sibling topics in this family are now reference files under the hubs (`research-methodology`, `deep-research`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: research-note-taking-pkm
title: "Research Note-Taking & Personal Knowledge Management (PKM)"
description: >-
  Design and operate a personal note/knowledge-management (PKM) system for a large, ongoing body of technical research — methods, tool choice, and the capture->synthesis->writing pipeline that turn captured material into durable, linkable knowledge and output. Covers Zettelkasten (atomic notes, stable IDs, linking, MOCs), Building a Second Brain (CODE/PARA, progressive summarization), evergreen notes, spaced repetition, the PKM tool landscape (Obsidian/Logseq/Roam/Notion/Tana), and anti-patterns (collector's fallacy, over-structuring, tool-churn). TRIGGER: setting up or fixing a note/PKM system; "which note app should I use"; organizing, linking, or atomizing research notes; building a slip-box or second brain; "my notes aren't useful". SKIP: drafting prose -> writing-expert; citation styles -> career-and-formal-writing; web sources -> deep-research; coding/labeling notes for analysis -> research-data-annotation-coding; customer doc stores -> content-ingestion-extraction.
category: custom
version: 1.1.1
updated: 2026-07-04
whenToUse:
  - "setting up or fixing a personal note / PKM system"
  - "choosing a note-taking method or PKM tool (Obsidian, Logseq, Notion, Roam, Tana)"
  - "how to take, organize, or link research notes; Zettelkasten, second brain, MOCs"
  - "progressive summarization, atomic/permanent/evergreen notes, spaced repetition"
  - "my notes aren't useful / collector's fallacy / PKM feels like procrastination"
keywords:
  - zettelkasten
  - personal knowledge management
  - pkm
  - note-taking
  - second brain
  - progressive summarization
  - evergreen notes
  - maps of content
  - spaced repetition
  - obsidian
tags:
  - research
  - knowledge-management
  - notes
  - writing
metadata:
  changelog:
    - "2026-06-15 sko v1->v1.1.0 — Pass H 10/10->10/10 pos, 0/10->0/10 neg; 5 Medium fixed (desc<=1000, 2 routing-resolvability, robust->well-supported, em-dash 2.99->1.16/100)"
---

# Research Note-Taking & Personal Knowledge Management (PKM)

Build and run a personal system that captures a large, ongoing stream of technical
research and turns it into durable, linkable knowledge — and eventually into output.
This skill is about the *capture/knowledge system itself*: the note methods, the tool
choice, and the capture→synthesis→writing pipeline. It is **not** about drafting prose,
finding/synthesizing web sources, or standing up shared engagement doc stores (see SKIP
routes in the description).

> Routing detail: reorganizing an existing messy note dump right now (not designing or choosing a system) → note-organizer.

The single most important reframing, from Andy Matuschak: **"'Better note-taking' misses
the point; what matters is 'better thinking.'"** Every method below is a means to develop
insight over time, not an end in itself. If a practice isn't producing thinking or output,
it is overhead.

## Core Concepts

- **A note system is a thinking partner, not a file cabinet.** The goal is a system that
  *surfaces relevant notes when you need them — even ones you didn't search for* (Ahrens),
  and that "gets exponentially better the more you feed it." A passive archive does the
  opposite (see Collector's Fallacy).
- **Knowledge should accrete.** Notes are valuable when new thoughts *combine with old
  ones* into a stronger whole, rather than scattering across documents. This is the
  organizing principle behind atomic, concept-oriented, linked notes.
- **Writing is the engine of understanding.** Putting an idea into your own words in a
  self-contained note *is* the act of comprehension. "Elaboration is nothing more than
  connecting information to other information in a meaningful way" (Ahrens). Highlighting
  and clipping are not.
- **Separate capture friction from processing rigor.** Capture should be near-zero
  friction (anywhere, anytime). Processing — turning captures into permanent knowledge —
  is deliberate, effortful, and where the value is created. Conflating the two produces
  either lost ideas or a bloated archive.

## Method 1 — Zettelkasten (the slip-box)

Originated by sociologist Niklas Luhmann (~90,000 cards over ~40 years; credited for ~600
articles and ~50+ books) and popularized for the digital era by Sönke Ahrens, *How to Take
Smart Notes* (2017). Three note types form a pipeline:

| Type | Lifespan | Purpose |
| --- | --- | --- |
| **Fleeting** | Hours–2 days | Frictionless capture of a passing thought. *Disposable by design* — process or delete within a day or two. A fleeting note alive a week later is clutter. |
| **Literature** | Kept with the source | Your selective, in-your-own-words response to one source ("what did I take away?"), with the reference. Not a full summary. Lives in/with your reference manager. |
| **Permanent (Zettel)** | Forever | One fully-articulated idea, in complete sentences, written *as if for publication*, linked into the network. This is the Zettelkasten. |

Three rules define a Zettel:

1. **Atomicity**: one idea per note. Atomicity is about *clear edges*, not brevity: a
   note is "atomic" when it has one idea you can link to from many directions and cite
   without dragging surrounding clutter. (Conflating atomic with "short" is a fallacy —
   atomicity is relative to your point of reference.) Use **two-step compression**: a tight
   note body plus a sharp, claim-like title.
2. **Stable unique IDs**: every note has a fixed, permanent address so later notes can
   reference it precisely, *forever*. Luhmann used branching alphanumeric codes (1 → 1a →
   1a1). Digital systems use timestamps, ULIDs/UUIDs, or backlink-resolved wiki-names.
   **Critical failure mode:** apps that use the *title* as the identifier break every link
   the moment you rename a note. Real implementations keep identity stable across edits.
3. **Linking**: connect notes deliberately, ideally with a phrase explaining *why* the
   link exists. The network of explicit links — not the folder it sits in — carries the
   intellectual value.

**Why atomic + linked beats folders:** "Luhmann's slip-box provides combinatorial
possibilities which were never planned, never preconceived." A single atomic note can
participate in many contexts at once; a note buried in one folder cannot.

## Method 2 — Building a Second Brain (CODE + PARA)

Tiago Forte's *Building a Second Brain* is a more output- and project-oriented framing.
"Our brains are for having ideas, not storing them."

**CODE**: the macro-loop:
- **Capture**: keep only what *resonates* (an intuitive, System-1 filter, not exhaustive coverage).
- **Organize**: file for *actionability*, using PARA.
- **Distill**: find the essence via progressive summarization.
- **Express**: ship something. The whole point is creative output.

**PARA**: organize by actionability, not by topic (subject taxonomies are ambiguous;
actionability takes a moment to judge):
- **Projects**: short-term efforts with a goal and an end.
- **Areas**: ongoing responsibilities to maintain over time.
- **Resources**: topics/interests possibly useful later.
- **Archive**: inactive items from the other three.

**Progressive summarization**: "opportunistic compression," done in small spurts across
time, only as much as a note deserves. Layers:
- L0 source → **L1** captured excerpt → **L2** bold the best passages → **L3** highlight
  the essential phrases within the bold → **L4** a one-line executive summary in your own words.

Forte's four guidelines keep it from becoming busywork: (1) don't apply all layers to all
notes; (2) use **resonance** (fast, intuitive) as the criterion, not analytical rules; (3)
*design the system for the laziest version of yourself*; (4) keep notes **glanceable** —
reviewable in seconds.

**Zettelkasten vs BASB (when to lean which way):** Zettelkasten optimizes for *idea
development and original synthesis* over years (researchers, writers, theory-building).
BASB optimizes for *project throughput and resurfacing* (knowledge workers shipping
deliverables). They compose: PARA can hold project/literature material while a Zettelkasten
holds the durable concept notes.

## Method 3 — Evergreen notes (Matuschak)

A rigorous sharpening of the permanent-note idea. Evergreen notes are "written and organized
to evolve, contribute, and accumulate over time, across projects." Principles:

- **Atomic**: one concept, so it can be linked and reused independently.
- **Concept-oriented**: factor notes by *concept*, not by source/author/book/project. This
  is the key move: when you read a second source on the same concept, the new thought lands
  on the *existing* concept note and the two combine. Note-per-book scatters; note-per-concept accretes.
- **Densely linked**: pushing yourself to add links forces you to think about how ideas
  relate (and aids memory via elaborative encoding). Note titles act "like APIs" you compose.
- **Prefer associative ontologies to hierarchical taxonomies**: links cut across fields;
  rigid trees don't. Tags are a comparatively weak association structure.
- **Write for yourself by default**, disregarding audience.

Leading indicator for knowledge work: *evergreen notes written per day.*

## Linking, backlinks & emergent structure (MOCs)

Structure should **emerge bottom-up** from links, not be imposed top-down by folders.
Folders are "rigid and exclusionary" — a note in a folder is separated from the rest of the
collection, which discourages interdisciplinary connection as the corpus grows.

**Maps of Content (MOCs)**: Nick Milo (Linking Your Thinking) — are the emergent-structure
tool. An MOC is *a note that links to other notes* on a theme. It behaves like a tag (groups
links non-exclusively; the notes live elsewhere) and a folder (tightly-packed grouping) at
once, and lets you position notes *in proximity* to each other.

MOCs do three things — **gather, develop, navigate**:
1. **Gather**: when you hit a *mental squeeze point* (~10–20 scattered notes on a thing and
   a tickle of overwhelm), make an MOC and drop links to all of them in one place.
2. **Develop**: the concentrated space sparks ideation; notes "collide" and combine.
3. **Navigate**: the MOC becomes a durable index/launchpad. Often a finished MOC means a
   future essay is "already 80% complete."

Practical model: links + backlinks are the substrate; MOCs are the optional index layer you
add *when a cluster gets unwieldy*, not preemptively. (Milo's A.C.C.E.S.S. shows folders and
link-based MOCs can coexist — a hybrid, not a religious war.)

## Spaced repetition for retention

Distinct from the note network: spaced repetition is for the subset of facts/skills you must
*recall on demand* (APIs, error codes, definitions, syntax) — not for everything you note.

- **Two mechanisms** (both well-supported in the literature, both underused because they feel
  counter-intuitive — Nature Reviews Psychology, 2022): **retrieval practice** (actively
  recalling beats re-reading) and **spacing** (reviews separated in time beat massed
  cramming). Study-phase-retrieval theory: each spaced review re-activates and strengthens
  the trace; too short = no retrieval needed, too long = trace already lost.
- **Mechanics** (Anki/SuperMemo): you grade each recall; intervals expand on success (e.g.
  1 day → 3 → 15 → 45…). Leitner boxes are the analog version.
- **Optimal spacing scales with the test delay** (Cepeda et al.): longer retention horizons
  want longer gaps. For indefinite retention, longer expanding intervals.
- **Integration tip:** keep flashcards downstream of permanent notes. Tools like RemNote/Logseq
  embed spaced-repetition cards directly in notes; or export a few high-value cards to Anki.
  Don't card everything — that recreates the collector's fallacy in flashcard form.

## Tool landscape & selection

The methods are tool-agnostic; tools differ on **data model** and **linking model**.

| Tool | Storage / model | Linking | Local-first | Best for |
| --- | --- | --- | --- | --- |
| **Obsidian** | Local **markdown** files in a vault | `[[wikilinks]]` + backlinks + graph | Yes (sync optional) | Power-user default; long-term PKM; 1,900+ plugins (Dataview, Templater, spaced-rep) |
| **Logseq** | Local markdown/Org, **outliner**, block-level | Block-level `[[links]]` + backlinks | Yes (OSS/AGPL) | Outliner thinkers; daily-journal capture; privacy/longevity; built-in flashcards |
| **Roam Research** | Cloud, outliner, block-level | Bidirectional block links | No | Pioneered networked-thought UX; outliner-first |
| **Notion** | Cloud **relational database** + blocks | Page links / relations | No (offline cache only) | Teams, structured DBs, multiple views, all-in-one workspace |
| **Tana** | Cloud graph DB, **supertags** | Tag/field-driven graph | Partial (offline mode) | Structured "notes-as-objects," AI-augmented capture, meeting intelligence |

**Decision criteria (use these, not feature checklists):**
1. **Data ownership / longevity.** If you're building a 10-year research corpus, weight the
   *Long Now* heavily (Kleppmann, "Local-First Software"): plain text and markdown are
   "the digital equivalent of stone" — readable for decades, on any tool, even if the app
   dies. Cloud-only tools give you *access, not ownership*; their database is the schema and
   you can't pipe it through future tools (incl. local LLMs). **The plain-text test:** export
   everything, open it in a dumb text editor — if it's readable, you own it; if it's binary
   blobs or proprietary XML, you're locked in.
2. **Linking vs database.** Want emergent, cross-domain idea networks → markdown + backlinks
   (Obsidian/Logseq). Want typed records with multiple views (tables/boards/calendars) →
   database model (Notion/Tana). Most *research* PKM wants the former; *operational* tracking
   wants the latter.
3. **Outliner vs document.** Think in nested bullets and daily logs → Logseq/Roam/Tana.
   Think in prose notes/pages → Obsidian/Notion.
4. **Collaboration.** Real-time multi-user → Notion. Obsidian/Logseq are single-user
   (shareable via Sync/Git, not live co-editing).
5. **Offline / unreliable connectivity.** Only the local-first tools work fully offline.

**Default recommendation for a technical researcher building an ongoing corpus:** Obsidian
(local markdown, stable, huge plugin ecosystem, the tool Ahrens himself settled on after
Roam). Logseq if you prefer an open-source outliner with first-class daily notes and
flashcards. Reserve Notion/Tana for structured/operational data, not the idea network.

## The capture → synthesis → writing workflow

This is the throughline for a researcher. Ahrens's 8-step model (lightly generalized):

1. **Fleeting capture**: always have a frictionless inbox; jot, don't curate.
2. **Literature notes**: while reading a source, write selective, in-your-own-words notes
   tied to the reference. *Work as if writing is the only thing that matters* — that purpose
   makes you read more rigorously and capture only what raises open questions.
3. **Permanent/evergreen notes**: process literature + fleeting notes into atomic,
   concept-oriented notes. Ask the elaborating questions: *What does it mean? How does it
   connect to…? How does it differ from…? What's it similar to?* This is the comprehension step.
4. **Link on the way in**: file each new permanent note *next to* related notes and add
   links (and, where useful, into an MOC). Filing forces you to find where it fits, which
   surfaces unexpected connections.
5. **Develop topics bottom-up**: don't start from a blank outline. Walk the note network,
   see what has clustered, what's missing, what questions arise; let topics emerge from
   accumulated notes ("start from abundance").
6. **Pick a topic from within the system**: choose what's already well-developed.
7. **Sequence into a draft**: select a *linear path* through the network, then translate
   (don't copy) notes into coherent prose. A literature review or essay becomes
   *rearrangement + connective tissue*, not blank-page writing.
8. **Edit & proofread**: at this point you're rewriting a draft, not generating from nothing.

**Cadence to sustain it:** daily — clear the capture inbox (≈5 min); weekly — a **review**
that processes literature/fleeting notes into permanent notes, links them, and pulls
relevant notes toward an active project. The weekly review is the load-bearing ritual; skip
it and the system silently becomes an archive.

**Hand-off note:** this pipeline *ends* where prose drafting and citation formatting begin —
route prose to `writing-expert` and citation styles (APA/Chicago/MLA/IEEE) to
`career-and-formal-writing` (→ `references/academic-and-citation-writing.md`). It *starts*
downstream of source discovery — route source-finding to `deep-research`.

## Pitfalls & anti-patterns

- **The Collector's Fallacy**: the big one. *"To know about something" ≠ "knowing
  something."* Saving feels like progress but changes nothing; "kept isn't read." Collected
  archives **anti-compound**: every saved item adds search cost, clutter, and decision
  fatigue, so the system gets *less* useful as it grows — the inverse of the promise. **Fix:**
  short research→read→assimilate cycles (process before collecting more), and always write a
  note in your own words. The Zettelkasten should fill, not the inbox.
- **PKM as sophisticated procrastination**: "capturing constantly, reviewing rarely,"
  endless linking with no output. **Litmus test:** in the last 30 days, what concrete output
  came *directly* from your notes — and what recurring ritual converts notes into next
  actions/outlines/drafts? If neither, you have a graveyard, not a brain. **Fix is smaller,
  not bigger:** add one weekly review, pick one project, produce one small artifact — *not*
  another tool.
- **Over-structuring**: building elaborate folder hierarchies, deep tag taxonomies, and
  ornate templates before you have notes. Structure should *emerge* from links/MOCs at the
  mental-squeeze point. Premature taxonomy is rigid and stunts cross-domain connection.
- **Tool-churn / migration thrash**: perpetually switching apps in search of the perfect
  one. "Changing tools midway is always a hassle" (Ahrens). Churn destroys the link graph and
  resets the habit. **Mitigation:** pick a local-first, plain-text tool *once* so your corpus
  outlives any single app, then stop shopping. The graph outlives the tool.
- **Atomicity-as-brevity**: chopping notes to be short rather than to hold one clear idea.
  Atomicity is about edges and linkability, not word count.
- **Title-as-ID fragility**: relying on a tool that breaks links when you rename a note.
  Prefer stable IDs / backlink resolution.
- **Verbatim transcription**: copying quotes/highlights instead of restating in your own
  words. Longhand/own-words processing is where understanding happens (cf. Mueller &
  Oppenheimer); a wall of highlights bypasses the brain.
- **Flashcarding everything**: spaced repetition is for must-recall facts only; carding
  indiscriminately is the collector's fallacy in another costume.

## Quick decision aids

- *"Where does this go?"* → Disposable thought = fleeting (process/delete in ≤2 days).
  Response to a source = literature note (with reference). A finished idea worth keeping =
  permanent/evergreen note (atomic, linked).
- *"Should I make an MOC?"* → Only when a cluster (~10–20 notes) gets unwieldy or you feel
  the mental squeeze. Not preemptively.
- *"Which tool?"* → Building a long-lived idea corpus → local-first markdown (Obsidian/Logseq).
  Structured operational data / team views → database (Notion/Tana). When unsure, optimize for
  the plain-text test and data ownership.
- *"Is my system healthy?"* → Run the 30-day output test and check the create:collect ratio of
  your last ~20 notes — aim for majority *your own words*, not clips.

## References
- zettelkasten.de: Ahrens's note categories explained; the atomicity guide; the Collector's Fallacy.
- Sönke Ahrens, *How to Take Smart Notes* (2017) / soenkeahrens.de: the 8-step workflow; "writing is the only thing that matters."
- Tiago Forte, *Building a Second Brain* / fortelabs.com: CODE, PARA, progressive summarization (guidelines & principles).
- Andy Matuschak, notes.andymatuschak.org: evergreen notes (atomic, concept-oriented, densely linked).
- Nick Milo / Linking Your Thinking: Maps of Content; folders vs tags vs links vs MOCs; A.C.C.E.S.S.
- Nature Reviews Psychology (2022), "The science of effective learning with spacing and retrieval practice"; Anki manual (background): spaced repetition.
- Martin Kleppmann et al., "Local-First Software: You Own Your Data": longevity, ownership, archival formats.
- Tool-landscape syntheses: Obsidian/Logseq/Roam/Notion/Tana, local-first vs database.
- buildfirstbrain.com / aethel.click / "The Zettelkasten Trap": collector's fallacy and PKM-procrastination anti-patterns.
