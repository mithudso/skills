<!-- hub-reference-banner -->
> **Reference file — part of the `frontend-ui` hub.** The named heuristic-evaluation method (Nielsen's 10 heuristics + the 0–4 severity process) and the Laws of UX (Yablonski / lawsofux.com) as a structured UI-critique instrument.
> Sibling topics in this family are reference files under the hubs (`frontend-ui`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").
>
> **Scope boundary — do not duplicate siblings.** General UX guidelines live in `references/ui-ux-pro-max.md`; WCAG / accessibility review lives in `references/accessibility-ux-reviewer.md`. This file owns only the *named* heuristic-evaluation method and the *named* Laws of UX, framed for detection → severity → fix.

---

---
name: usability-heuristics-laws-of-ux
title: Usability Heuristic Evaluation & Laws of UX
description: >-
  Named usability-inspection instruments for critiquing a UI or product screen:
  Nielsen's 10 usability heuristics (each with detection signals + fixes), the
  heuristic-evaluation method (evaluator process, the 0–4 severity scale, 3–5
  evaluators, the discovery formula, aggregation, heuristic evaluation vs
  usability testing), and the Laws of UX (Yablonski/lawsofux.com — Hick, Fitts,
  Jakob, Miller, Tesler, Postel, Doherty, Aesthetic-Usability, Von Restorff,
  Serial Position, Peak-End, Zeigarnik, Gestalt grouping laws) with each law's
  concrete UI implication. TRIGGER: running a heuristic evaluation; flagging
  usability violations by heuristic; rating severity 0–4; how many evaluators;
  citing a Law of UX to critique a UI screen; detection-and-fix for a usability
  smell. SKIP: general UX/visual-design guidance → ui-ux-pro-max; WCAG audit →
  accessibility-ux-reviewer; a UX law applied to a journey/retention/behavior
  goal, or deep decision/persuasion psychology → applied-psychology.
category: developer
tags:
  - developer
  - frontend
  - ux
  - usability
  - design-critique
  - heuristic-evaluation
keywords:
  - usability heuristics
  - heuristic evaluation
  - Nielsen 10 heuristics
  - severity rating scale
  - Laws of UX
  - Jakob's Law
  - Fitts's Law
  - Hick's Law
  - Gestalt principles
  - usability inspection
whenToUse:
  - "running a heuristic evaluation of a UI, screen, or flow against Nielsen's 10 heuristics"
  - "flagging a usability violation and naming which heuristic it breaks"
  - "assigning a 0–4 severity rating to a usability problem and prioritizing fixes"
  - "deciding how many evaluators a heuristic evaluation needs"
  - "citing a Law of UX (Hick, Fitts, Jakob, Miller, Tesler, Doherty, etc.) in a design critique"
  - "translating a usability smell into a concrete detection signal and fix"
whenNotToUse:
  - "general UX or visual-design guidance (styles, palettes, components) — use ui-ux-pro-max"
  - "WCAG 2.2 / WAI-ARIA accessibility audit — use accessibility-ux-reviewer"
  - "applying a UX law to a multi-step journey, retention, or behavior-change goal (not a screen critique) — use applied-psychology"
  - "deep psychology of decision-making, persuasion, or memory — use applied-psychology"
related_skills:
  - ui-ux-pro-max
  - accessibility-ux-reviewer
  - applied-psychology
  - deceptive-design-and-dark-patterns
version: "1.0.1"
updated: "2026-06-16"
---

# Usability Heuristic Evaluation & Laws of UX

A structured **critique instrument** for inspecting a UI or product screen. Two
complementary bodies of practice:

1. **Heuristic evaluation** — Nielsen & Molich's usability-inspection method: a
   few evaluators judge an interface against the **10 usability heuristics**,
   log each violation as a discrete problem, and rate it on a **0–4 severity**
   scale.
2. **Laws of UX** — Jon Yablonski's curation (lawsofux.com / O'Reilly book) of
   psychology and HCI principles, each with a concrete design implication you
   can check on a screen.

Use this as the **dimension list + scoring rubric** for a design/image critique:
for each finding, name the heuristic or law it violates, rate its severity, and
give a concrete fix. Heuristic hits are **hypotheses to verify**, not proof;
expect inter-rater variance and validate the severe ones with real users.[^evaleffect][^false]

## Contents

- [How to run a critique pass](#how-to-run-a-critique-pass)
- [Part A — Nielsen's 10 usability heuristics (detection + fix)](#part-a--nielsens-10-usability-heuristics-detection--fix)
- [Part B — The heuristic-evaluation method](#part-b--the-heuristic-evaluation-method)
  - [Evaluator process](#evaluator-process)
  - [The 0–4 severity rating scale](#the-04-severity-rating-scale)
  - [Number of evaluators & the discovery formula](#number-of-evaluators--the-discovery-formula)
  - [Aggregating & prioritizing findings](#aggregating--prioritizing-findings)
  - [Heuristic evaluation vs usability testing](#heuristic-evaluation-vs-usability-testing)
- [Part C — The Laws of UX (principle → UI implication)](#part-c--the-laws-of-ux-principle--ui-implication)
- [Anti-patterns & honest limits](#anti-patterns--honest-limits)
- [References](#references)

## How to run a critique pass

For a single-screen or whole-flow critique:

1. **Walk the interface twice.** First pass: learn the flow and scope without
   judging. Second pass: inspect each element against each heuristic/law.[^howto]
2. **Log each problem discretely.** One row per problem: *description · location ·
   heuristic or law violated · severity (0–4) · suggested fix.* List problems
   separately so each can be counted, rated, and fixed.[^howto][^severity]
3. **Rate severity** using the 0–4 scale (below), weighing frequency × impact ×
   persistence (+ market impact).[^severity]
4. **Prioritize** by severity vs fix cost; lead the report with 3s and 4s.[^severity]
5. **Name the law where it sharpens the "why."** A Law of UX often explains the
   mechanism behind a heuristic violation (e.g., a cluttered choice screen =
   Aesthetic/Minimalist violation, explained by Hick's Law).

A scannable shorthand per finding: **DETECT** (the observable smell) → **SEVERITY**
(0–4) → **FIX** (the remediation).

**Latency rule of thumb** (one rule, cited several ways below): **<0.1 s** feels
instant; **<1 s** keeps the user's flow (show a spinner/progress once a wait
crosses ~1 s); **<10 s** holds attention; treat the Doherty Threshold's **400 ms**
as the *productivity* target to design toward.[^responsetimes]

**Worked finding row** (the instrument applied once): *Checkout shows no spinner
after "Pay" — dead air past ~1 s* (**DETECT**) → **Severity 3** (frequent, high
impact, repeats) → optimistic spinner + disabled button while the request runs
(**FIX**, Heuristic #1 Visibility of system status).

---

## Part A — Nielsen's 10 usability heuristics (detection + fix)

The 10 heuristics were developed by **Jakob Nielsen & Rolf Molich in 1990** and
refined by Nielsen in **1994** via factor analysis of 249 usability problems.[^ten]
NN/g **updated the names, descriptions, and examples in 2020, but the 10
heuristics themselves are unchanged since 1994** — so a critique keyed to these
names is current; only the explanatory copy was modernized.[^ten] Names and order
below are the canonical NN/g list.[^ten]

> Confidence note: every heuristic *definition* below is canonical NN/g fact
> (3+ sources). Detection signals and fixes are drawn from the canonical article
> plus the per-heuristic NN/g deep-dive (footnoted per heuristic); individual
> product examples within a deep-dive are illustrative, single-source. Treat the
> definitions and the recurring signals as solid; treat any lone product example
> as illustrative rather than load-bearing.

### 1. Visibility of system status
**Means:** Keep users informed about what is going on, with appropriate feedback
in reasonable time.[^h1]
**DETECT:** no feedback after a click/submit; missing loading/progress indicator
("dead air"); silent state changes (cart/wishlist items vanish with no notice);
hidden current location / active state in nav.[^h1][^ten]
**FIX:** immediate visual confirmation (color change + checkmark); progress bar
for any operation over ~1s; status microcopy/badges ("Only a few left", "$12 to
free shipping", out-of-stock labels); a visible "you are here" in nav.[^h1]

### 2. Match between the system and the real world
**Means:** Speak the users' language; follow real-world conventions; present
information in a natural, logical order.[^h2]
**DETECT:** internal jargon, technical terms, unexplained acronyms; system-centric
labels (DB/engineering terms) instead of user terms; info out of expected order;
icons/controls that resemble no familiar object.[^h2][^ten]
**FIX:** plain headlines + relatable copy with simple definitions; natural
mappings (volume "+" above "−"); real-world metaphors (highlighter-style marking,
page-turn albums); sequence steps the way users think about the task.[^h2]

### 3. User control and freedom
**Means:** Users act by mistake, so give a clearly marked "emergency exit"
(undo/cancel/back) to leave the unwanted state without an extended process.[^h3]
**DETECT:** no undo/redo for reversible actions; Back restarts a process instead
of stepping back one level; close/cancel hidden, mislabeled, or non-standard;
ambiguous "X" (save or discard?); destructive actions with no confirmation; flows
that lose work if abandoned.[^h3][^ten]
**FIX:** visible Cancel + a confirmation that offers "discard" vs "save draft";
overlays that respond to both in-UI and browser back; inline Remove links;
discoverable Undo/Redo; a post-action snackbar with "Undo".[^h3]

### 4. Consistency and standards
**Means:** Users shouldn't wonder whether different words/situations/actions mean
the same thing. Two sub-types — **internal** (within the product family) and
**external** (match conventions of other sites, i.e. Jakob's Law).[^h4]
**DETECT:** the same concept named differently across screens (or different things
named alike); button positions/styles shifting page to page; inconsistent
icon/border/image treatment; non-standard placement of expected elements (home
top-left, cart, search); inconsistent input formats.[^h4][^ten]
**FIX:** reuse one component pattern across the product; lock button order,
placement, and labels; conform to web conventions (blue/underlined links, cart
icon, logo→home top-left).[^h4]

### 5. Error prevention
**Means:** The best designs prevent problems before they occur, eliminating
error-prone conditions or surfacing a confirmation before commit. Grounded in
**slips** (right intent, wrong action) vs **mistakes** (flawed plan).[^h5]
**DETECT:** inputs not constrained to valid ranges (return date before departure
selectable); rigid formatting that rejects valid variants; no sensible defaults
in repetitive tasks; no autocomplete where typos are likely; no safeguard before
a costly/irreversible commit.[^h5][^ten]
**FIX:** constraints (calendar widget blocking invalid ranges); good defaults
(preset options); forgiving formatting (auto-format/strip as the user types);
autocomplete; a confirmation step before irreversible commits.[^h5]

### 6. Recognition rather than recall
**Means:** Minimize memory load by making elements, actions, and options visible;
users shouldn't have to remember info from one part of the interface to
another.[^h6]
**DETECT:** users must carry codes/IDs/values between screens; hidden options that
must be remembered; unlabeled icons; gesture-only interactions with no affordance;
reliance on memorized commands; format rules stated once, not shown at the point
of entry.[^h6][^ten]
**FIX:** show recently viewed items / search history; autocomplete (turns
*generation* into *recognition*); visible menus and labeled icons; persistent
favorites; contextual tips tied to the page rather than one-time tutorials.[^h6]

### 7. Flexibility and efficiency of use
**Means:** Accelerators — unseen by novices — let experts speed up; allow users to
tailor frequent actions, serving both inexperienced and experienced users.[^h7]
**DETECT:** only one rigid path; experts forced through novice hand-holding; no
keyboard shortcuts (or undiscoverable ones); no batch/bulk operations; no
personalization; preferences don't persist across sessions; no way to save/reuse
frequent actions.[^h7][^ten]
**FIX:** keyboard shortcuts shown unobtrusively (on hover); saveable
workspace/layout presets; multi-select for bulk commands; macros/automation for
power users; accelerators that don't confuse novices.[^h7]

### 8. Aesthetic and minimalist design
**Means:** Interfaces shouldn't hold irrelevant or rarely needed info — every
extra unit competes with the relevant units. **This is signal-to-noise /
content prioritization, NOT visual sparseness:** a dense page can satisfy it; an
over-stripped page can violate it.[^h8]
**DETECT:** visual clutter; many competing CTAs of equal weight; no hierarchy;
decorative elements crowding functional ones; redundant content in prime space;
**over-minimalism** (so little info users can't tell what the product does or find
functions).[^h8][^ten]
**FIX:** keep only content that supports a task; give negative space a purpose;
progressive disclosure of advanced options; establish visual hierarchy
(scale/contrast) — info-dense-but-purposeful is fine.[^h8]

### 9. Help users recognize, diagnose, and recover from errors
**Means:** Error messages in plain language (no codes), precisely stating the
problem and constructively suggesting a solution.[^h9]
**DETECT:** raw error codes / jargon / stack traces; vague "An error occurred";
message far from its source or low-contrast; accusatory tone ("illegal",
"invalid"); no recovery path; user input lost on error; color-only error
indication; errors fired prematurely during exploration.[^h9][^ten]
**FIX:** inline message adjacent to the field (red text + icon + border), stating
the exact problem + at least one next step; non-blaming tone; preserve input /
offer save-as-draft; proactively warn before a predictable failure (e.g.,
missing-attachment send warning).[^h9]

### 10. Help and documentation
**Means:** Ideally no extra explanation is needed; when documentation is required
it should be searchable, task-focused, list concrete steps, and stay
concise.[^h10]
**DETECT:** no search in the help area; help as one undifferentiated scroll (no
titles/grouping); help not tied to the user's goal/context; docs that use terms
they never define; dense paragraphs instead of steps; launch-time tutorials that
interrupt before the user has a question, with no skip.[^h10][^ten]
**FIX:** contextual help triggered by the current task; searchable, categorized
docs with headings + numbered steps + keyword-highlighted results; multimodal
help (video + text); make pushed tutorials skippable.[^h10]

---

## Part B — The heuristic-evaluation method

Heuristic evaluation is a **usability-inspection method**: a small set of
evaluators examines an interface and judges its compliance with the heuristics,
producing a prioritized list of problems. It belongs to the inspection family
(alongside cognitive walkthrough), distinct from usability testing with real
users, and is the flagship of Nielsen's **"discount usability engineering."**[^howto][^inspect][^measuringu]

### Evaluator process

- **Inspect alone, then combine.** Each evaluator inspects **independently**;
  only after all evaluations are complete may they communicate. Independence is
  what makes the discovery model valid and yields unbiased problem lists.[^howto][^theory]
- **Two passes each.** First pass: go through the task once **just to learn the
  system**, without evaluating (flow & scope). Second pass: go back through and
  inspect each element, feature, and decision against the 10 heuristics. NN/g
  reserves ~1–2 hours per session.[^howto]
- **Record discrete, heuristic-tagged problems.** Each issue is a separate entry
  paired with the heuristic number it violates (and its location); fix
  recommendations go in a separate column.[^howto]
- **Support domain-heavy systems.** Evaluators with little domain knowledge are
  especially likely to miss problems; a facilitator may supply domain mechanics
  but not usability judgments.[^found]

### The 0–4 severity rating scale

Verbatim NN/g scale:[^severity]

| Rating | Meaning |
| --- | --- |
| **0** | I don't agree that this is a usability problem at all |
| **1** | Cosmetic problem only — need not be fixed unless extra time is available |
| **2** | Minor usability problem — fixing this should be given low priority |
| **3** | Major usability problem — important to fix, so should be given high priority |
| **4** | Usability catastrophe — imperative to fix before the product can be released |

**Three factors set severity** (plus a 4th consideration):[^severity]
1. **Frequency:** is it common or rare?
2. **Impact:** easy or hard for users to overcome?
3. **Persistence:** a one-time hurdle once learned, or a repeated bother?
4. **Market impact:** some problems devastate a product's popularity even if
   "objectively" easy to overcome.

**Collected after the eval, and averaged.** Severity is **not** assigned live —
during inspection evaluators are focused on *finding* problems and each sees only
a fraction of the set. Ratings are gathered afterward via a questionnaire listing
the *combined* problem set; because single-evaluator severity scores are
unreliable, NN/g averages them — "the mean of a set of ratings from three
evaluators is satisfactory for many practical purposes."[^severity]

### Number of evaluators & the discovery formula

- **One evaluator finds ~35% of problems** (mean across six of Nielsen's
  projects).[^theory]
- **Per-evaluator detection rate λ ≈ 0.31–0.34** (the six-study mean was ~34%;
  the canonical curve uses Nielsen & Landauer's pooled ≈31%).[^theory][^landauer]
- **Discovery formula:** `ProblemsFound(i) = N · (1 − (1 − λ)^i)`, where *i* =
  number of independent evaluators, *N* = total problems in the interface, *λ* =
  proportion found by one evaluator.[^theory][^landauer]
- **3–5 evaluators is the sweet spot.** Nielsen recommends "about five, but
  certainly at least three"; the cost-benefit ratio in his worked example peaks
  around **four**, then marginal gains flatten. Model-derived curve points:
  ~35% (1) → ~75% (3) → ~85% (5) (the 75/85 figures are model-derived, quoted by
  secondary sources rather than tabulated by NN/g).[^howto][^theory]
- **Use more for mission-critical/high-traffic systems**, where large payoffs
  justify it.[^theory]
- **Prefer "double experts."** Detection rises with expertise: novices ~22%,
  usability specialists ~41%, **double experts (usability + domain) ~60%** in
  Nielsen's 1992 study.[^doubleexpert]

### Aggregating & prioritizing findings

- **Merge the independent lists, consolidate duplicates** (NN/g uses affinity
  diagramming to cluster the same underlying problem); the distinct-problem count
  is what `ProblemsFound(i)` refers to.[^howto][^theory]
- **Each consolidated problem carries the mean 0–4 rating** from the post-session
  questionnaire round.[^severity]
- **Prioritize by averaged severity vs fix cost** — the scale wording itself
  encodes cost ("fix only if spare time" … "imperative"). Lead the report with
  3s and 4s.[^severity][^howto]
- **Deliverable:** a prioritized list — each problem with description, heuristic(s)
  violated, location, averaged severity, and an optional fix recommendation.[^howto][^severity]

### Heuristic evaluation vs usability testing

| | Heuristic evaluation | Usability testing |
| --- | --- | --- |
| Who judges | Experts inspecting against principles | Real users attempting tasks |
| Needs a working system? | No — runs on prototypes/specs | Usually yes |
| Cost/speed | Cheaper, faster | Costlier, slower |
| Catches | More total + more minor problems, cheaply | Fewer but more **severe**, behavior-dependent problems |
| Weakness | **False alarms** (~34% avg, range 9–46%) | Smaller raw count per session |

- **Complementary, not substitutes** — Nielsen frames HE as a complement to (not a
  replacement for) user research.[^measuringu][^howto]
- **Overlap is partial.** HE's overlap with usability-test findings averages ~36%
  (range 30–43%); it **misses ~49%** of problems users hit. "HE will typically
  find between 30% and 50% of problems found in a concurrent usability
  test."[^measuringu]
- **Severity skew:** more severe problems tend to surface through user testing;
  HE produces a far larger raw count (Jeffries et al. 1991).[^measuringu]
- **When to use which:** run HE early and cheaply to clear obvious violations on
  prototypes before spending costly user sessions; use usability testing to
  validate against real behavior and catch the severe, behavior-dependent
  problems HE misses — then iterate with both.[^measuringu][^howto]

---

## Part C — The Laws of UX (principle → UI implication)

Source spine: **lawsofux.com** and Jon Yablonski, *Laws of UX* (O'Reilly,
2020).[^lawsofux][^yablonskibook] lawsofux.com is one author's curation, so a claim
resting only on it is **one** evidentiary point — contested laws below are
cross-checked against IxDF, NN/g, and primary papers. For each law: the principle
(with origin) → what to **check on a screen**. A **★** marks the laws whose UX
application most outruns its evidence (apply the *practice*, drop the shaky
*rationale*). Apply these laws to critique a concrete UI screen or layout; for a
multi-step *journey*, retention, or behavior-change goal, the underlying
psychology belongs to `applied-psychology`, not this critique pass.

| Law (origin) | Principle | UI implication — what to check |
| --- | --- | --- |
| **Hick's Law** (Hick & Hyman, 1952)[^lawsofux] | Decision time grows with the number/complexity of choices. | Reduce, stage, or guide choices; highlight a recommended/default; use progressive disclosure. **★ Caveat:** the log-time mechanism is **often misused in HCI** — most on-screen selection is visual search, near-constant for familiar controls.[^hickcrit] |
| **Fitts's Law** (Fitts, 1954)[^lawsofux] | Time to acquire a target = f(distance, size). | Make tap/click targets large and well-spaced; put primary actions big and near attention. *Touch nuance:* the desktop "edges/corners are infinitely deep" corollary **inverts on touch** — top corners are the hardest one-thumb region; reward bottom-reachable nav.[^fittstouch][^asktog] |
| **Jakob's Law** (Jakob Nielsen)[^lawsofux] | Users spend most time on *other* sites and expect yours to work the same. | Follow established conventions (control patterns, layout, icon/placement); on a redesign, let users keep the familiar version. Flag novelty without payoff. |
| **Miller's Law** (Miller, 1956)[^lawsofux] | The average person holds ~7±2 items in **working memory**. | Use **chunking** to group related content/controls. **★ Misinterpretation:** the "max 7 items on screen" rule is a myth — Miller's limit is *recall*, not *visible* options (recognition). True WM capacity is ~4 (Cowan 2001). Never praise/penalize a screen for hitting "7."[^millermyth][^cowan] |
| **Tesler's Law** (Larry Tesler, Xerox PARC)[^lawsofux] | Conservation of complexity — every system has irreducible complexity that must live somewhere. | The design should absorb inherent complexity (smart defaults, pre-filled/derived values) rather than dump it on the user. Flag screens that offload computable work onto people. |
| **Postel's Law** (Jon Postel)[^lawsofux] | Be liberal in what you accept, conservative in what you send. | Inputs should tolerate variant formats (phone/date/card spacing) and normalize them, with clear feedback; output stays precise. Flag rigid fields that reject obviously-valid input. |
| **Doherty Threshold** (Doherty & Thadani, IBM, 1982)[^lawsofux] | Productivity soars when response stays **<400 ms**. | Keep feedback under ~400 ms; otherwise mask with skeletons, progress, optimistic UI. *Keep honest:* pair with the better-corroborated Nielsen/Card **0.1 s / 1 s / 10 s** response-time limits rather than treating 400 ms as the only number.[^responsetimes] |
| **Aesthetic-Usability Effect** (Kurosu & Kashimura, 1995)[^lawsofux] | Users *perceive* attractive design as more usable. | Polish raises satisfaction and buys tolerance for minor flaws — but can **mask** real usability problems in testing. *Bounded:* the effect is on *perceived* usability; effect on *objective performance* is small/heterogeneous (meta-analytic g≈0.12) and fades with experience. Don't let looks substitute for task success.[^aesthcrit] |
| **Von Restorff Effect** (von Restorff, 1933)[^lawsofux] | The item that differs most is best remembered (isolation effect). | Make the single most important action visually distinct (one clear focal CTA). Guardrails: don't over-emphasize multiple elements (they compete; can read as ads), don't rely on color alone. |
| **Serial Position Effect** (Ebbinghaus)[^lawsofux] | First (primacy) and last (recency) items are best remembered. | Put the most important nav/list items at the **start and end**; bury low-priority ones in the middle. |
| **Peak-End Rule** (Kahneman et al., 1993)[^lawsofux] | People judge an experience by its emotional **peak** and its **end**, not the average. | Design the journey's most intense moment and its final moment to delight; negative peaks dominate. *Caveat:* evidence is from short/clinical/aversive episodes — a heuristic for long product journeys, not a measurement law.[^peakendcrit] |
| **Zeigarnik Effect** (Zeigarnik, 1927)[^lawsofux] | Uncompleted tasks are remembered better than completed ones. | Use progress indicators and endowed-progress (show early progress to motivate completion); surface incomplete-task cues (drafts, "resume," completion meters). **★ Caveat:** the *memory* claim is **weakly evidenced** — a 2025 meta-analysis found ~zero recall advantage; what survives is the *pull to resume*, which still justifies progress UI.[^zeigarnikcrit] |
| **Gestalt — Proximity**[^lawsofux] | Objects near each other are perceived as a group. | Spacing/whitespace should bind related items and separate unrelated ones. |
| **Gestalt — Similarity**[^lawsofux] | Visually similar elements are perceived as related / same function. | All buttons of one type look alike; interactive elements are visually differentiated from plain text. |
| **Gestalt — Common Region**[^lawsofux] | Elements inside a shared bounded area are perceived as a group. | Use cards/borders/background fills to bound related content. |
| **Gestalt — Uniform Connectedness**[^lawsofux] | Visually connected elements (lines, frames, shared container) are perceived as *more* related than merely-near ones — the strongest grouping cue. | Tie related controls together with an explicit connector/container. |

The Gestalt grouping laws (origin: Wertheimer & the Gestalt psychologists, master
principle **Prägnanz** — perceive the simplest form) are the **best-grounded**
family here and are not seriously contested as UX guidance.[^lawsofux] Note that
"Jakob's Law," "Tesler's Law," and the "Doherty Threshold" are **named by their
popularizers** — useful mnemonics, not formal scientific laws.[^millermyth]

---

## Anti-patterns & honest limits

Surface these so a critique isn't over-trusted:

- **Heuristic hits are hypotheses, not proof.** ~34% of HE-flagged items are not
  real problems when validated with users; prefer concrete observable smells over
  abstract judgments, and verify the severe ones with testing.[^false][^measuringu]
- **The evaluator effect is real and irreducible.** Different evaluators on the
  same UI find markedly different problem sets (pairwise agreement 5–65%); it must
  be *managed* (multiple evaluators, double experts, averaged severities), not
  eliminated.[^evaleffect]
- **Single-evaluator severity scores are unreliable** — always average across
  evaluators.[^severity][^evaleffect]
- **Heuristics are abstract** and overlap (Recognition-vs-Recall ↔ Minimalist ↔
  Consistency all touch cognitive load), so single-heuristic tagging is fuzzy;
  the 10 were shaped for desktop and don't explicitly cover
  emotion/accessibility/voice/mobile.[^uxpa]
- **Don't weaponize a misread law.** The two most-misapplied: Miller's "7±2" as a
  cap on visible options (it isn't), and Hick's *formula* cited for visible menus
  (the mechanism doesn't transfer cleanly).[^millermyth][^hickcrit]
- **Don't let polish hide problems.** The Aesthetic-Usability Effect means an
  attractive screen can pass review while still failing users.[^aesthcrit]
- **A heuristic can be weaponized into a dark pattern.** These same laws power
  *deceptive* design: the **bait-and-switch "X"** that opens an offer violates
  **Jakob's Law**, **Von Restorff** salience builds a **false hierarchy**, and the
  **Aesthetic-Usability Effect** is exactly how a polished skin masks a
  manipulative flow. When a usability finding looks deliberately user-hostile
  (asymmetric reject, roach-motel cancel, drip pricing), switch to the
  design-ethics lens — `deceptive-design-and-dark-patterns` — which rates user
  harm **plus** legal exposure (DSA Art. 25, GDPR, FTC, CCPA/CPRA) and prescribes
  the honest fix.

---

## References

[^ten]: Nielsen Norman Group — "10 Usability Heuristics for User Interface Design" (canonical list, names/order, 1990/1994 origin, 2020 language update). https://www.nngroup.com/articles/ten-usability-heuristics/ — *docs/canonical*. verified-as-of: 2026-06-16.
[^h1]: NN/g — "Visibility of System Status (Usability Heuristic #1)". https://www.nngroup.com/articles/visibility-system-status/ — *docs/canonical*.
[^h2]: NN/g — "Match Between the System and the Real World (Usability Heuristic #2)". https://www.nngroup.com/articles/match-system-real-world/ — *docs/canonical*.
[^h3]: NN/g — "User Control and Freedom (Usability Heuristic #3)". https://www.nngroup.com/articles/user-control-and-freedom/ — *docs/canonical*.
[^h4]: NN/g — "Consistency and Standards (Usability Heuristic #4)". https://www.nngroup.com/articles/consistency-and-standards/ — *docs/canonical*.
[^h5]: NN/g — "Slips" / error-prevention (slips vs mistakes, Usability Heuristic #5). https://www.nngroup.com/articles/slips/ — *docs/canonical*.
[^h6]: NN/g — "Recognition vs. Recall in UX (Usability Heuristic #6)". https://www.nngroup.com/articles/recognition-and-recall/ — *docs/canonical*.
[^h7]: NN/g — "Flexibility and Efficiency of Use (Usability Heuristic #7)". https://www.nngroup.com/articles/flexibility-efficiency-heuristic/ — *docs/canonical*.
[^h8]: NN/g — "Aesthetic and Minimalist Design (Usability Heuristic #8)" (signal-to-noise, not sparseness). https://www.nngroup.com/articles/aesthetic-minimalist-design/ — *docs/canonical*.
[^h9]: NN/g — "Error-Message Guidelines" (Usability Heuristic #9). https://www.nngroup.com/articles/error-message-guidelines/ — *docs/canonical*.
[^h10]: NN/g — "Help and Documentation (Usability Heuristic #10)". https://www.nngroup.com/articles/help-and-documentation/ — *docs/canonical*.
[^howto]: NN/g — "How to Conduct a Heuristic Evaluation" (two passes, independent inspection, discrete heuristic-tagged problems, ~1–2 hr sessions, 3–5 evaluators). https://www.nngroup.com/articles/how-to-conduct-a-heuristic-evaluation/ — *docs/canonical*. verified-as-of: 2026-06-16.
[^theory]: NN/g (Nielsen) — "How to Conduct a Heuristic Evaluation: The Theory" / "How Many Test Users in a Usability Study"-adjacent theory page (35% single-evaluator finding, λ range/mean, discovery formula, cost-benefit ~4, independence). https://www.nngroup.com/articles/how-to-conduct-a-heuristic-evaluation/theory-heuristic-evaluations/ — *docs/canonical*.
[^severity]: NN/g (Nielsen) — "Severity Ratings for Usability Problems" (verbatim 0–4 scale, frequency/impact/persistence + market impact, post-session questionnaire, averaging across evaluators). https://www.nngroup.com/articles/how-to-rate-the-severity-of-usability-problems/ — *docs/canonical*. verified-as-of: 2026-06-16.
[^found]: NN/g (Nielsen) — "Usability Problems Found by Heuristic Evaluation" (domain expertise; 42% major / 32% minor per single evaluator). https://www.nngroup.com/articles/usability-problems-found-by-heuristic-evaluation/ — *docs/canonical*.
[^doubleexpert]: Nielsen, J. (1992) — "Finding Usability Problems Through Heuristic Evaluation," CHI '92 (novice 22% / usability specialist 41% / double expert 60%). https://dl.acm.org/doi/10.1145/142750.142834 — *paper* (via citing secondary sources; primary PDF not text-extracted this run — treat exact tiers as qualified).
[^landauer]: Nielsen, J. & Landauer, T. K. (1993) — "A mathematical model of the finding of usability problems," INTERCHI '93 (pooled λ≈31%, basis of the discovery curve). https://dl.acm.org/doi/10.1145/169059.169166 — *paper* (via citing sources).
[^inspect]: Wikipedia — "Usability inspection" (inspection-method definition and taxonomy). https://en.wikipedia.org/wiki/Usability_inspection — *encyclopedia*.
[^measuringu]: MeasuringU — "How Effective are Heuristic Evaluations?" (~36% hit, ~34% false alarm, ~49% miss; 30–50% range; Jeffries et al. 1991 severity skew; complementarity). https://measuringu.com/effective-he/ — *blog (empirical)*. verified-as-of: 2026-06-16.
[^evaleffect]: Hertzum, M. & Jacobsen, N. E. (2003) — "The Evaluator Effect: A Chilling Fact About Usability Evaluation Methods," *IJHCI* 15(1) (pairwise agreement 5–65%; effect must be managed, not eliminated). https://mortenhertzum.dk/publ/IJHCI2003.pdf — *paper*.
[^false]: MeasuringU — "The Value of Multiple Evaluators in Heuristic Evaluations" (false-positive trade-off; single-evaluator coverage). https://measuringu.com/he-multiple/ — *blog (empirical)*.
[^uxpa]: UXPA Magazine — "Nielsen's Heuristic Evaluation: Limitations in Principles and Practice" (per-heuristic vagueness, coverage gaps, never-validated critique). https://uxpamagazine.org/nielsens-heuristic-evaluation/ — *magazine/critique*.
[^lawsofux]: Yablonski, J. — Laws of UX (lawsofux.com), canonical per-law pages (principle, Takeaways, Origins) for Hick, Fitts, Jakob, Miller, Tesler, Postel, Doherty, Aesthetic-Usability, Von Restorff, Serial Position, Peak-End, Zeigarnik, and the Gestalt laws (Proximity, Similarity, Common Region, Uniform Connectedness, Prägnanz). https://lawsofux.com/ — *docs/canonical*. verified-as-of: 2026-06-16.
[^yablonskibook]: Yablonski, J. (2020) — *Laws of UX: Using Psychology to Design Better Products & Services*, O'Reilly (2nd ed. ISBN 9781492055303). https://www.oreilly.com/library/view/laws-of-ux/9781492055303/ — *book*.
[^millermyth]: UX Myths — "Myth #23: Choices should always be limited to 7±2" (Miller's limit is recall, not visible options; the number was a rhetorical device; chunking is the real contribution). https://uxmyths.com/post/931925744/myth-23-choices-should-always-be-limited-to-seven — *blog/industry*.
[^cowan]: Cowan, N. (2001) — "The magical number 4 in short-term memory: A reconsideration of mental storage capacity," *Behavioral and Brain Sciences* (true WM capacity ~4, revises Miller). https://memory.psych.missouri.edu/assets/doc/articles/2001/cowan-bbs-2001.pdf — *paper*.
[^hickcrit]: Liu, Gori, Rioul, Beaudouin-Lafon & Guiard (CHI 2020) — "How Relevant is Hick's Law for HCI?" (the stimulus-response paradigm Hick measured "is rarely relevant to HCI tasks"; most UI selection is visual search, not choice-reaction, and familiar controls give near-constant times). https://perso.telecom-paristech.fr/rioul/publis/202001liugoririoulbeaudouinlafonguiard.pdf — *paper*. Corroborated by IxDF "The Hick-Hyman Law" and Cockburn et al. (CHI 2007) "A Predictive Model of Menu Performance."
[^fittstouch]: Smashing Magazine / Hoober, S. (2022) — "Fitts' Law In The Touch Era" (edge/corner corollary inverts on touch; top corners hardest for one-thumb use). https://www.smashingmagazine.com/2022/02/fitts-law-touch-era/ — *blog/industry*.
[^asktog]: Tognazzini, B. — "First Principles of Interaction Design" (Fitts edges/corners "infinitely deep" on desktop). https://asktog.com/atc/principles-of-interaction-design/ — *blog/industry*.
[^responsetimes]: NN/g (Nielsen) — "Response Times: The 3 Important Limits" (0.1 s instant / 1 s flow / 10 s attention; from R. Miller 1968 & Card et al. 1991). https://www.nngroup.com/articles/response-times-3-important-limits/ — *docs/canonical*.
[^aesthcrit]: Tuch, A. N. et al. (2012) — "Is beautiful really usable? Toward understanding the relation between usability, aesthetics, and affect in HCI," *Computers in Human Behavior* (a single lab study: aesthetics did not move perceived usability when usability was strongly varied; the "what is beautiful is usable" arrow can reverse). https://www.sciencedirect.com/science/article/abs/pii/S0747563212000908 — *paper*. Separately, a meta-analysis of visual aesthetics → *objective* performance reports only a small effect (Hedges' g ≈ 0.12, "Visual Aesthetics and Performance: A First Meta-Analysis"); Grishin & Gillan (JUS 2019) find the effect "weak and diminishes with experience." The bounded-effect claim rests on these *objective-performance* studies, not on Tuch alone.
[^peakendcrit]: Kemp, S., Burt, C. D. B., Furneaux, L. (2008) — "A test of the peak-end rule with extended autobiographical events," *Memory & Cognition* (peak-end "not an outstandingly good predictor" for long real-world experiences); cf. OBHDP (2022) meta-analysis supporting the rule overall for short episodes. https://link.springer.com/article/10.3758/MC.36.1.132 — *paper*.
[^zeigarnikcrit]: Ghibellini, R. & Meier, B. (2025) — "Interruption, recall and resumption: a meta-analysis of the Zeigarnik and Ovsiankina effects," *Humanities and Social Sciences Communications* (recall ratio ≈0.99 excluding Zeigarnik's own 1927 data — the memory advantage does not reliably replicate; resumption/goal-intrusion survives). https://www.nature.com/articles/s41599-025-05000-w — *paper*.
