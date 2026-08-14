<!-- hub-reference-banner -->
> **Reference file — part of the `frontend-ui` hub.** Sibling topics in this family are
> reference files under the hubs (`frontend-ui`) — **not** standalone skills. Ignore any
> "use the X skill" / `related_skills` / SKIP pointers that name a bare sibling skill; load
> that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---
name: visual-design-principles-and-critique
description: >-
  Visual-design critique rubric for graphic/brand and UI/UX product screens. Named composition
  principles (visual hierarchy, gestalt, C.R.A.P., balance, emphasis, white space, scale/proportion,
  unity/rhythm, Z/F/Gutenberg scanning) each as DETECT a violation → RATE severity → FIX; the
  how-to-run-a-critique methodology is the companion file visual-design-critique-methodology.
  TRIGGER: critique/review a visual design, mockup, UI screen, landing page, or brand/marketing
  asset; "is this design good", "what's wrong with this layout"; judge visual
  hierarchy/balance/spacing/typography; rate a design issue's severity. SKIP: photographic or
  AI-generated-image critique; producing a NEW design (use web-design / ui-ux-pro-max);
  accessibility-only WCAG audit (use accessibility-ux-reviewer); a named Nielsen heuristic
  evaluation or a Law of UX as the lens (use usability-heuristics-laws-of-ux) — this file owns
  visual composition and only borrows Nielsen's 0–4 scale to rate findings.
category: developer
version: "1.0.2"
updated: "2026-06-16"
keywords:
  - visual design critique
  - design critique methodology
  - visual hierarchy
  - gestalt principles
  - CRAP principles
  - design severity rating
  - white space
  - scanning patterns
  - design QA
  - typographic scale
tags:
  - developer
  - frontend
  - design
  - ux
  - critique
  - review
whenToUse:
  - "critiquing or reviewing a visual design, UI mockup, screen, landing page, or brand asset"
  - "answering 'is this design good?' or 'what's wrong with this layout?'"
  - "judging visual hierarchy, balance, spacing, alignment, or typography on a screen"
  - "running or facilitating a structured design critique session"
  - "rating the severity of design issues or running a design QA / spec-parity pass"
  - "giving actionable, non-prescriptive design feedback to a designer"
whenNotToUse:
  - "critiquing photographs or AI-generated images — out of scope"
  - "producing a NEW design from scratch — use web-design / ui-ux-pro-max"
  - "accessibility-only audit against WCAG — use accessibility-ux-reviewer"
  - "running a named Nielsen heuristic evaluation or citing a Law of UX as the lens — use usability-heuristics-laws-of-ux"
related_skills:
  - frontend-ui
  - accessibility-ux-reviewer
  - web-design
  - ui-ux-pro-max
---

# Visual Design Principles as a Critique Rubric + Design-Critique Methodology

A working rubric for evaluating **graphic/brand design and UI/UX product screens**, plus the
**methodology** for running and giving design critiques. Built for critique, not creation: every
principle is expressed as **DETECT a violation → RATE its severity → FIX it**, and the methodology
section covers how a professional crit is structured, scoped, and delivered.

Scope: graphic/brand design and UI/UX product-screen critique. **Out of scope:** photographic and
AI-generated-image critique. For *producing* designs see `references/web-design.md` /
`references/ui-ux-pro-max.md`; for accessibility audits see `references/accessibility-ux-reviewer.md`;
for usability-heuristic + Laws-of-UX evaluation see the `frontend-ui` hub.

## Contents

- **Part A — How to use this rubric** (the critique pass + the severity scale)
- **Part B — The principles, as a rubric** (Detect / Rate / Fix per principle)
  1. Visual hierarchy
  2. Emphasis & focal point
  3. Gestalt principles
  4. The C.R.A.P. principles
  5. Balance
  6. White space / negative space
  7. Scale & proportion
  8. Unity, variety & visual rhythm
  9. Visual flow & scanning patterns
  10. Consistency (visual treatment)
- **Part C — Design "smells"** (flawed-screen anti-patterns, mapped to the principles)
- **Part D — Design-critique methodology** → companion file `visual-design-critique-methodology.md`
  (crit formats, frameworks, giving feedback, severity & design QA, running a crit, critique smells)
- **Part E — Quick reference** (the squint-test pass, severity cheat-sheet)
- **References**

---

# Part A — How to use this rubric

## The critique is a pass, not a vibe

Run the screen through the principles in Part B **in order**, and for each one ask the three rubric
questions: *Can I DETECT a violation? How severe is it? What's the FIX?* Lead with the **squint
test** (Part E) — it surfaces hierarchy, emphasis, balance, and grouping failures in one look.

**Separate observation from judgment.** This is the spine of every credible critique method
(see Part D / the methodology companion). State what you *see* before you state whether it's *good* — "the heading and the CTA are
the same size and weight" (observation) precedes "the user has no entry point" (interpretation)
precedes any verdict. Tie every finding to a **user goal or design objective**, never to personal
taste.[^conn][^feldman][^nngcrit]

## Severity: use Nielsen's 0–4 scale for everything

Rate every finding on **Nielsen's usability-severity scale**, which is a function of three
factors — **frequency** (how many users hit it), **impact** (how hard it is to overcome), and
**persistence** (one-time vs. repeated):[^nngsev]

| Rating | Meaning | Action |
| --- | --- | --- |
| **0** | Not a problem | — |
| **1** | **Cosmetic** | Fix only if spare time (e.g. ~2px misalignment, slight weight inconsistency) |
| **2** | **Minor** | Low priority |
| **3** | **Major** | High priority — impairs the task / comprehension |
| **4** | **Catastrophe** | Must fix before release |

Two cautions from the source:
- **Severity ≠ priority.** Severity is experiential damage; *priority* is business urgency, and they
  are independent. A **misspelled brand name or off-brand color on the homepage is cosmetic severity
  but high priority** because every visitor sees it on day one. Track the two separately so a
  high-visibility cosmetic issue can still jump the queue.[^betterqa][^uxsl]
- **Single-rater severity is unreliable.** Nielsen's own finding: ratings from one evaluator are too
  noisy to trust; the mean of ~3 evaluators is satisfactory. Don't over-engineer the number — a
  100-point formula exists but Nielsen explicitly *does not* recommend it for real projects.[^nngsev]

---

# Part B — The principles, as a rubric

## 1. Visual hierarchy

**Principle.** The deliberate arrangement of elements so the eye consumes them in order of
importance — what's seen first, second, third. NN/g: hierarchy "lets users know where to focus their
attention." It's built from *combinations* of **size, weight, color/contrast, position, and
spacing/density** — not one lever alone. The core failure: treating every element as equally
important is the same as treating every element as unimportant.[^nnghier][^digpolo]

**The levers (each raises an element's perceived importance):** larger size; heavier weight; higher
*value/saturation contrast* against neighbors (it's the contrast, not the hue); earlier reading
position (top, top-left in LTR); and more surrounding whitespace.[^nnghier][^artstyle]

**DETECT**
- Squint or blur the screen (5–10px). If a secondary label, a decorative image, a stock photo, or the
  nav dominates instead of the headline/primary CTA, hierarchy is broken.[^nnghier][^artstyle]
- "Everything is loud": many elements at similar size/weight/saturation, lots of ALL-CAPS, nothing
  recedes.[^nnghier][^digpolo]
- More than **~3 type sizes** or **>3 contrast levels** in play; arbitrary weight mix.[^nnghier]
- No single entry point — the eye has to pause and decide where to look first.[^bradley]

**RATE**
- **Catastrophe/Major (4–3):** No discernible entry point, or the *wrong* element dominates (an ad/
  decoration out-shouts the primary task). Users can't form a scan order.[^nnghier][^bradley]
- **Major (3):** Hierarchy exists but is muddy — too many sizes/contrast levels blur adjacent tiers.
- **Minor/Cosmetic (2–1):** Reads correctly overall; one secondary element slightly mis-weighted.

**FIX**
- Rank the content's importance *first*, then assign visual weight to match.[^nnghier][^bradley]
- Collapse to ≤3 type sizes and ≤3 contrast levels; reserve one accent color for the single most
  important element.[^nnghier]
- Widen the size gap between levels and add whitespace between sections.
- Re-run the squint test **on the real, populated screen** — a strongly colored content image can
  hijack a hierarchy the template never intended.[^nnghier]

## 2. Emphasis & focal point

**Principle.** Emphasis makes one element the most important; that element becomes the **focal
point** — the entry point answering "where do I look first?" Emphasis creates a *single* focal point;
hierarchy then sequences everything else. Emphasis is **relative** — something must recede so the
focal point can advance. The reliable ranking has **3 levels: dominant → sub-dominant →
subordinate** (more levels reduce the contrast between neighbors).[^designyourway][^madegood][^smashdom]

**Three techniques:** **dominance** (make it bigger/bolder/more saturated), **subordination** (quiet
*everything else* down), and **isolation** (separate it with whitespace — the Von Restorff/isolation
effect). Strong compositions usually combine ≥2.[^designyourway][^vonrestorff]

**DETECT**
- **No focal point:** squint and nothing emerges first; the eye has nowhere to land.[^vonrestorff][^smashdom]
- **Too many focal points:** several elements use full-strength emphasis at once — "if everything is
  highlighted, nothing is."[^vonrestorff]
- The textbook UI instance: **two equally-weighted filled buttons** side by side (e.g. "Start Free
  Trial" and "Schedule Demo," both solid, same size/color).[^vonrestorff][^cta]
- A primary CTA styled the same as feature icons or "Read more" links — present but camouflaged.

**RATE**
- **Catastrophe/Major (4–3):** Two+ elements compete for *primary* emphasis, splitting attention; in
  conversion UI this is measurable harm (choice paralysis on dual equal CTAs). A page with *no* focal
  point at all is also critical.[^madegood][^cta]
- **Major (3):** A single primary exists but a sub-dominant element occasionally steals first
  attention.
- **Minor (2):** Focal point reads first but is weaker than ideal.

**FIX**
- One dominant focal point per view/section; build everything else around it, then verify viewers see
  it first.[^designyourway][^madegood]
- Competing CTAs: keep one filled/high-contrast primary; demote the secondary to an **outline/ghost
  button or text link**, or move it below the fold.[^cta]
- Isolate the primary with whitespace and a reserved accent color; pair color with size/shape so the
  cue survives for color-blind users.[^vonrestorff]
- If you can't reduce competition, **subordinate** — quiet everything else rather than shouting
  louder.[^designyourway][^vonrestorff]

## 3. Gestalt principles

**Principle.** Gestalt psychology describes how the eye groups visual elements into wholes. In design
these are the laws behind *perceived grouping and structure*; violating them makes a layout read as
the wrong groups.[^ixdfgestalt][^nnggestalt] The load-bearing ones for critique:

- **Proximity** — elements close together are perceived as a group; this is the single strongest
  grouping cue (it can override similarity).[^ixdfgestalt][^nngproximity]
- **Similarity** — elements sharing color/shape/size/orientation are seen as related.[^ixdfgestalt]
- **Common region** — elements inside a shared boundary (a card, a box, a shaded panel) are grouped,
  even against proximity/similarity.[^ixdfgestalt]
- **Closure** — the eye completes incomplete shapes (logos, icon outlines).[^ixdfgestalt]
- **Continuity (good continuation)** — the eye follows lines/curves and aligned elements as a
  connected path.[^ixdfgestalt]
- **Figure-ground** — the eye separates a subject (figure) from its background; ambiguous figure-
  ground is a defect in UI (which is the modal? which is the page?).[^ixdfgestalt][^lidwellfg]
- **Common fate** — elements moving/animating together are grouped (carousels, parallax).[^ixdfgestalt]

**DETECT**
- **Proximity violation:** related items (a label and its field; a heading and its paragraph) spaced
  as far apart as unrelated ones — or unrelated items packed as tightly as related ones. This is the
  most common real-world gestalt defect.[^ixdfgestalt][^nngproximity]
- **Similarity miscue:** items that *aren't* related share a style (two unrelated things both blue and
  underlined → both read as links); or related items styled differently so they don't read as a set.
- **Common-region conflict:** a card boundary that encloses items belonging to two different groups,
  or a border that splits one group.
- **Figure-ground ambiguity:** insufficient contrast/elevation between a modal/overlay and the page;
  the user can't tell what's active.[^ixdfgestalt][^lidwellfg]

**RATE**
- **Major (3):** A grouping miscue that makes users misread *what belongs together* — e.g. a form
  label floating equidistant between two fields, or a "destructive" item visually grouped with safe
  ones. Figure-ground ambiguity on an interactive overlay is Major+.[^nngproximity][^lidwellfg]
- **Minor/Cosmetic (2–1):** Grouping reads correctly; spacing or styling is slightly off-rhythm.

**FIX**
- Apply proximity deliberately: tighten space *within* a group, widen space *between* groups — prefer
  a whitespace gap over a divider line.[^ixdfgestalt][^smashprinc]
- Use similarity to signal "same kind"; reserve a distinct style for distinct roles.
- Reach for **common region** (a card/panel) only when proximity alone can't carry the grouping —
  it's a strong cue and overusing boxes fragments the page.[^ixdfgestalt]
- Restore figure-ground with contrast, shadow/elevation, or a scrim behind overlays.[^lidwellfg]

## 4. The C.R.A.P. principles

**Principle.** Robin Williams's four foundational layout principles (*The Non-Designer's Design
Book*): **Contrast, Repetition, Alignment, Proximity.** They are the fastest first-pass checklist for
amateur-looking layouts.[^crapttd][^crapsmash]

- **Contrast** — if two elements are not the same, make them *very* different. Timid contrast (a
  heading only slightly bigger than body; near-identical greys) is the #1 amateur tell.[^crapttd]
- **Repetition** — repeat visual elements (color, type treatment, spacing, bullet style) throughout to
  create cohesion and reinforce identity.[^crapttd]
- **Alignment** — every element should have a *visual connection* to something else on the page;
  nothing arbitrary. Prefer one strong edge over centering everything.[^crapttd]
- **Proximity** — group related items; this is the same cue as gestalt-proximity (Part B§3) — Williams
  arrived at it independently as a layout rule.[^crapttd][^ixdfgestalt]

**DETECT**
- **Contrast:** elements that are *similar but not the same* (slightly different sizes/weights/colors)
  — the eye reads it as a mistake, not a choice.[^crapttd]
- **Repetition:** inconsistent treatment of the same role across the page (three bullet styles; two
  heading fonts; spacing that drifts).[^crapttd]
- **Alignment:** text/elements with no shared edge; a centered title over left-aligned body over a
  right-floated image — multiple competing axes; "almost aligned" (off by a few px).[^crapttd]
- **Proximity:** see gestalt-proximity above.

**RATE**
- **Major (3):** Broken alignment that fragments the layout, or contrast so weak the hierarchy
  collapses. **Minor/Cosmetic (2–1):** one mildly inconsistent repetition or a single off-axis
  element.

**FIX**
- Push contrast further than feels comfortable (size *and* weight *and* color, not one).[^crapttd]
- Pick one alignment edge per region and hold it; align to a grid.[^crapttd]
- Codify the repeated elements (type scale, spacing, bullet/icon style) so repetition is automatic —
  this is where design tokens earn their keep (Part B§10).

## 5. Balance

**Principle.** Balance is the distribution of **visual weight** across a composition; balanced reads
as stable, unbalanced reads as "tipping over," creating unintended tension. Visual weight comes from
size, color intensity/saturation, value (darkness), contrast, detail/density, and distance from
center.[^balpolo][^balramotion]

**Types:** **symmetrical (formal)** — mirrored across an axis; reads as order/stability/authority
(corporate, finance, government); self-balancing and easiest. **Asymmetrical (informal)** — different
elements whose *weights* still balance; reads as dynamic/modern; harder, and it *tips* when it fails.
**Radial** — elements radiate from a center (sunbursts, emblems, pie charts). **Mosaic/
crystallographic** — even all-over distribution, no focal point.[^balpolo][^balramotion][^balkittl]

**DETECT**
- Squint and check weight distribution: heavy blocks bunched on one side with no counterweight → it
  visibly leans.[^balpolo]
- An entire side left blank while the other is crammed (reads unfinished, not as intentional negative
  space).[^balpolo]
- **Forced symmetry:** asymmetric content jammed into a symmetric grid, leaving awkward dead
  quadrants.[^balpolo]
- An asymmetric desktop layout that loses its weighting when it stacks on mobile.[^balkose]

**RATE**
- **Major (3):** The composition visibly tips or looks unfinished — one side carries nearly all weight
  with no counterbalance, undermining credibility.[^balpolo][^balkose]
- **Minor (2):** A slight lean that *could* read as intentional energy — flag only if it conflicts
  with the intended tone (e.g. a bank that should feel stable).
- **Rater note:** weight *asymmetric* layouts more critically than symmetric ones — symmetry self-
  balances, asymmetry fails by tipping. And confirm intent: imbalance is sometimes a deliberate
  tension device.[^balsmash][^balpolo]

**FIX**
- Counterbalance the heavy side — a stronger headline, a color block, a larger image, or a
  *constellation of several small elements* opposite one large one.[^balpolo][^balkose]
- Use negative space as an *active* counterweight, not leftover emptiness.[^balkose]
- Match the balance type to the message (symmetry = authority; asymmetry = modern/energetic; radial =
  ceremonial center).[^balpolo][^balramotion]

## 6. White space / negative space

**Principle.** The empty space between and within elements. **Macro** whitespace separates major
groups and paces the page; **micro** whitespace (line-height, paragraph spacing, letter-spacing,
padding) drives legibility. **Active** negative space is placed to shape reading order and focus;
**passive** exists for breathing room. Generous whitespace also reliably signals quality/premium;
dense "every-pixel-used" layouts read as mass-market.[^nngmin][^ixdfspace][^smashprinc]

**Micro-whitespace numerics worth knowing for critique:**
- Body **line-height ≥ 1.5** (WCAG 2.2 baseline); raise to **1.6–1.7** for long lines.[^lineheight]
- Line length (measure) **~50–75 characters** (~66 ideal); outside that, reading degrades.[^measure]
- Leave body **letter-spacing** alone; add **+5–12%** tracking only to ALL-CAPS.[^typography]

> **Deep type-craft defects → `references/micro-typography-craft.md`** (frontend-ui hub). This file
> owns the high-level typographic rubric (scale, measure, line-height above). The expert,
> screenshot-detectable *micro*-typography defects — widows/orphans/runts, rag/rivers/bad breaks,
> kerning ("keming") & display tracking, baseline-grid rhythm, justified-text gaps & H&J, hanging
> punctuation & optical sizing, glyph crimes (curly vs straight quotes, en/em dashes, faux small caps,
> old-style/tabular figures, ligatures, true vs faux bold/italic), and "type crimes" (stretched/
> condensed faux styles, ransom-note font mixing) — live there as DETECT → RATE → FIX entries. Load it
> for the typography pass of a detailed critique.

**DETECT**
- Major sections butt together and read as one mass; "wall of content"; edge-to-edge content with
  near-zero outer margin (cramped, low-quality feel).[^nngmin][^smashprinc]
- **Clutter:** high element density, low signal-to-noise; decorative/redundant elements crowding what
  users need (NN/g: every extra unit competes with the relevant ones).[^nngclutter]
- **Over-application:** related items pushed so far apart they stop grouping; so much vertical space
  the page looks finished above the fold and users never scroll ("illusion of
  completeness").[^nngmin][^portent]
- Micro: line-height < 1.5 on long-form; measure far outside 50–75; tight/negative tracking on body;
  controls jammed against container edges with no padding.[^lineheight][^measure]

**RATE**
- **Major (3, usability):** Missing macro space breaks grouping so users misread what belongs together;
  *or* excess space hides essential content/buries the primary action; *or* micro misuse measurably
  impairs reading (line-height <1.5 on long-form, zero padding crowding tap targets); *or* clutter
  obscures necessary elements.[^nngmin][^nngclutter][^lineheight]
- **Minor (2, aesthetic):** Grouping is legible and content findable; the layout merely feels tight or
  sparse. Perceived-quality/premium-feel issues are usually Minor unless they undermine a trust moment
  (cluttered checkout/pricing).[^ixdfspace]

**FIX**
- Chunk into scannable groups with consistent section padding and outer margins; prefer a gap over a
  divider.[^nngmin][^smashprinc]
- Spend *active* space on the one or two elements that matter most — surround the primary CTA/hero with
  more space than its neighbors.[^ixdfspace]
- Fix clutter by **removing** irrelevant elements (not just shrinking them); maximize signal, minimize
  noise.[^nngclutter]
- Body line-height 1.5 (1.6–1.7 long lines); constrain measure to 50–75 chars; give every block and
  control real padding.[^lineheight][^measure]

## 7. Scale & proportion

**Principle.** A **modular scale** is a sequence derived from a base value × a consistent ratio, used
to size type, spacing, and components so everything relates proportionally (popularized for the web by
Tim Brown's "More Meaningful Typography"). Formula: `size = base × ratio^step`.[^scalemic][^modscale][^ala]

**Common typographic ratios:** 1.125 Major Second / 1.2 Minor Third (dense data UI); **1.25 Major
Third** and **1.333 Perfect Fourth** (the default for most digital products — Tailwind/Material/Stripe
cluster here); 1.5 Perfect Fifth and 1.618 Golden Ratio (display/hero, *use sparingly* with few
levels).[^scalemic][^brainy][^scalecalc]

**The 8-point grid.** Constrain margins, padding, gaps, and component sizes to multiples of 8 (8, 16,
24, 32, 48, 64), with a 4px half-step; Material aligns components to an 8dp grid, type/icons to 4dp.
8 divides cleanly and scales without fractional pixels across densities.[^scale8pt][^android] **Tension
to flag:** modular type ratios produce non-integer pixels (1.25³ = 31.25px) that fight the 8px grid;
the documented reconciliation is to round scale outputs to multiples of 4/8 for UI.[^scalemic]

**DETECT**
- **Arbitrary sizing off any scale:** type sizes plucked out of the air; inconsistent intervals between
  heading levels; spacing values not on a 4/8 grid (13px, 17px, 22px, 35px with no
  relationship).[^scale8pt][^modscale]
- **Over-applied scale:** a "plethora of similar sizes" — dozens of generated values, near-duplicates
  (18px and 19px), implementers unsure which to use.[^krycho]
- **Density mismatch:** a 1.618 scale in a dense dashboard (body feels tiny next to headings), or a
  1.125 scale on a marketing hero with no drama.[^brainy]

**RATE**
- **Major (3):** No consistent sizing system at all — every value ad hoc, visibly inconsistent screen-
  to-screen. Also Major: a fixed-px (non-rem) scale that **clips/overflows at 200% zoom** (accessibility-
  adjacent).[^scale8pt][^codexical]
- **Minor (2):** Mostly on-scale with a few pragmatic exceptions (rounding a value to fit the 8px grid).
  This is *endorsed* practice — the scale is a design aid, not a constraint. A *single* off-grid value
  is not critical.[^fontfyi]

**FIX**
- Pick one base (commonly 16px) and one ratio matched to product density (1.25–1.333 default); generate
  the scale, then **use a limited subset** (≤6 heading levels).[^scalemic][^krycho]
- Reconcile with the 8-point grid (round outputs to 4/8); express the scale as CSS custom
  properties/design tokens so one change cascades.[^scale8pt][^scalemic]
- Use **rem-based** sizes + flexible containers so the scale survives 200% zoom; consider fluid
  `clamp()` scales.[^codexical]

> **⚠️ Golden-ratio caveat (do not weaponize in critique).** The claim that φ ≈ 1.618 is a universal
> law of beauty or was deliberately used in the Parthenon/Mona Lisa/pyramids is **largely debunked** —
> a 2024 peer-reviewed review finds the evidence "without foundation," and the Parthenon's façade
> does not fit a golden rectangle (Markowsky measures a width-to-height ratio of ~1.71, stricter
> readings ~2.25 — none of them φ).[^goldenpmc][^goldendevlin] Treat the golden ratio and
> the **rule of thirds** as *optional* compositional aids, never as quality criteria: **do not flag a
> design as defective for "not following the golden ratio," and do not praise one because a φ overlay
> fits.** φ remains a legitimate *typographic ratio* (one option among many, for display contexts) —
> that narrow use is uncontested.[^modscale][^goldenmade]

## 8. Unity, variety & visual rhythm

**Principle.** **Unity** is the sense that elements belong to one coherent whole; **variety** is the
interest that prevents monotony — they're a deliberate tension. **Visual rhythm** is the repetition of
elements at intervals (regular, flowing, or progressive) that paces the eye, like beats in
music.[^lidwellunity][^smashprinc]

**DETECT**
- **Too much unity → monotony:** everything the same weight/spacing/color; nothing to anchor or relieve
  the eye; the page reads as undifferentiated.
- **Too much variety → chaos:** many fonts/colors/styles/spacings competing; no through-line; reads as
  several designs collided.
- **Broken rhythm:** spacing/sizing intervals that should repeat but jump unpredictably (a card grid
  with uneven gutters; a vertical rhythm that stutters).

**RATE**
- **Major (3):** Chaos that defeats scanning, or monotony that buries the hierarchy. **Minor (2):** a
  single rhythm break or one element slightly off the system.

**FIX**
- Establish a repeating system (type scale, spacing scale, color roles) for unity; introduce variety
  *intentionally and sparingly* at the points you want to emphasize.[^smashprinc]
- Restore rhythm by snapping intervals to the spacing scale (Part B§7).

## 9. Visual flow & scanning patterns

**Principle.** Where the eye goes on a screen, and which layout fights or follows it. Match the pattern
to the *content type* — the patterns are diagnostics, not templates to design toward.[^nngscan]

- **F-pattern** — on **text-heavy** content (articles, search results, long copy), users scan a top
  horizontal sweep, a shorter second sweep, then down the left edge. Note: the F appears **when
  designers have done nothing to guide the eye** — it's a *symptom* of unformatted text, not a goal.
  Grounded in NN/g eye-tracking (original 2006 study, 232 users).[^nngf][^nngfdisc]
- **Layer-cake** — eyes scan headings/subheads and drop into the matching section. NN/g calls this "by
  far the most effective" scan short of reading everything — it's the *desired* outcome.[^nngscan]
- **Spotted** — users hunt for one specific thing (a link, a number, a keyword); scattered fixations on
  standout elements.[^nngscan]
- **Commitment** — near word-by-word reading; driven by user *motivation*, not layout — don't expect it
  on low-motivation pages.[^nngscan]
- **Bypassing** — when stacked list items start with the same word, users skip the repeated leading
  words; they see only ~2 words per item, so a wasted first word hides the differentiator.[^nngbypass]
- **Z-pattern** — a *heuristic* (not eye-tracking) for **sparse, image-led** layouts (hero/landing/
  posters): top-left → top-right → diagonal → bottom-right.[^zpattern]
- **Gutenberg diagram** — for **homogeneous, evenly-distributed** text: reading gravity sweeps top-left
  (primary optical area) → bottom-right (terminal area); top-right and bottom-left are weak "fallow"
  corners. A 1950s model, explicitly "a framework, not a law," overridden by any real visual
  hierarchy — the **weakest-evidenced** pattern here.[^lidwellgut][^gutituon]

**DETECT**
- A wall of text with no subheads/bullets/bold (invites the lossy F-scan); key info or CTAs parked on
  the right rail or far down.[^nngf]
- List/nav items that all start with the same generic word ("Introducing…", "The…") → bypassing hides
  the meaningful word.[^nngbypass]
- A sparse landing page whose primary CTA sits at a low-attention point instead of the Z-terminus.[^zpattern]
- A homogeneous text layout with its CTA in a fallow corner.[^lidwellgut]

**RATE**
- **Major (3):** F-pattern present as a *symptom* over body content (real comprehension loss + missed
  content); bypassing on conversion/nav links (users can't distinguish options).[^nngf][^nngbypass]
- **Minor–Major (2–3):** Misplaced CTA on a true sparse/Z layout (costs clicks). Gutenberg-corner
  placement is **Low–Minor** and weakly evidenced — apply only to genuinely homogeneous layouts.[^lidwellgut]

**FIX**
- Break the wall: descriptive headings/subheads, bullets, short paragraphs, front-loaded keywords, most
  important content in the first two paragraphs — converts an F-scan into the better layer-cake.[^nngf][^nngscan]
- Front-load the distinguishing keyword in list items; strip repeated leading words.[^nngbypass]
- On sparse layouts, anchor the Z: logo top-left, primary **CTA at bottom-right** (or echoed top-
  right). If copy is heavy, abandon Z and design for F/layer-cake.[^zpattern]
- Rely on explicit hierarchy over Gutenberg for anything non-homogeneous.[^lidwellgut]

## 10. Consistency (visual treatment)

**Principle.** Similar parts should look and behave similarly so users transfer knowledge instead of
relearning. NN/g heuristic #4. Two axes: **internal** (within your product/family) and **external**
(**Jakob's Law** — users expect your site to work like the others they already know). Two layers:
**aesthetic** (style/brand) and **functional** (a signifier always triggers the same
action).[^nngconsist][^jakobslaw]

**DETECT**
- Same action styled differently across screens; more than one "primary" button style; drifting type
  scale/spacing/color between screens; a repeated element (nav/CTA/field) that moves position across a
  flow; mixed icon families (solid + outline) for peers.[^nngconsist][^iconincon]
- **External:** logo not top-left / not linking home; search not a magnifier or in an odd place;
  non-links styled as links; standard icons given idiosyncratic meanings.[^nngconsist][^jakobslaw]

**RATE**
- **Major (3):** Functional inconsistency (same icon → two meanings; a color that means "destructive"
  in one place and "neutral" in another) — it induces errors. Broken *external* convention (logo
  doesn't go home, search unfindable) trends **3–4** — it fails recognition on first contact.[^nngconsist][^jakobslaw]
- **Cosmetic–Minor (1–2):** Pure aesthetic drift (off-brand shade, minor weight mismatch) unless it
  erodes brand recognition at scale.
- **Re-frame, don't auto-flag:** an inconsistency that *improves* distinguishability or efficiency is
  not a defect (Google's 2020 over-uniform app icons became too similar to tell apart — "foolish
  consistency"). The defect is *unintentional* drift or convention-breaking with no user
  payoff.[^nngconsist][^tyranny]

**FIX**
- Centralize visual decisions in **design tokens** (global → semantic → component) so values are
  referenced, not hard-coded; back with a documented component library (one component per role, defined
  states); run periodic consistency audits — drift is ongoing maintenance.[^nngconsist][^tokens]
- Align core patterns to platform/web conventions; reserve novelty for high-value moments and document
  the exception so the next person doesn't read it as drift.[^jakobslaw][^nngconsist]

---

# Part C — Design "smells" (flawed-screen anti-patterns)

What a flawed screen looks like — each maps back to a principle in Part B:

- **Everything is bold / no hierarchy** — flat emphasis, no entry point (Part B§1).
- **Two equal primary CTAs** — split focus, choice paralysis (Part B§2).
- **Proximity mismatch** — related items spaced apart, unrelated items packed together (Part B§3).
- **Timid contrast / "almost aligned"** — the amateur tells (Part B§4).
- **Tipping layout** — one heavy side, no counterweight (Part B§5).
- **Wall of content / cramped, no whitespace** — or *too much* whitespace hiding the fold (Part B§6).
- **Arbitrary sizing / size soup** — values off any scale, or a scale over-generated into near-
  duplicates (Part B§7).
- **Inconsistent icon/button/color treatment across screens** — functional inconsistency (Part B§10).

---

# Part D — Design-critique methodology → companion reference

**How to *run* the critique and *deliver* the findings lives in the companion file
`references/visual-design-critique-methodology.md`.** Load it when the task shifts from *what's wrong
with this screen* (this file) to *how do I run/give the critique*. It covers:

- **Critique vs. feedback vs. review** — and why combining critique + review in one meeting fails.
- **Structured crit formats** — I-Like-I-Wish-What-If (d.school/IDEO), Describe/Interpret/Evaluate
  (Feldman — the canonical observation-before-judgment sequence), Rose/Bud/Thorn, design studio/charrette.
- **Connor & Irizarry, *Discussing Design*** — critique≠reaction; frame around objectives not taste;
  the facilitator reformulates directive comments back to goals; don't problem-solve mid-critique.
- **Actionable-but-non-prescriptive feedback** — flag the *problem* ("X feels buried"), not the
  *solution* ("make X red"); the translation habit; ask questions, not directives; giver/receiver dynamics.
- **Severity & design QA** — Nielsen's 0–4 scale (Part A) + the Blocker→Cosmetic ladder; the design-QA /
  spec-parity checklist; design-QA as a release gate.
- **Running a critique** — objectives up front, who attends, sync vs. async, time-boxing, closing the loop.
- **Critique anti-patterns** — the compliment sandwich (and why it's discouraged), HiPPO, bikeshedding,
  design-by-committee, vague praise.

**Affective / trust axis → companion reference.** When the question is *does this design feel right,
trustworthy, and on-tone at first glance* (rather than the composition mechanics this file owns), load
`references/emotional-design-and-visual-trust.md`: Norman's visceral/behavioral/reflective levels, the
first-impression/halo read, Stanford/Fogg web-credibility and trust signals, and brand personality
(Aaker, warmth×competence) as a detect→rate→fix pass.

---

# Part E — Quick reference

## The squint-test pass (do this first)

Squint or blur the screen 5–10px and ask, in order:
1. **What stands out first?** → If it's not the intended primary element, **hierarchy/emphasis** fail
   (B§1–2).
2. **Does it tip?** → weight bunched on one side = **balance** fail (B§5).
3. **What groups together?** → wrong groupings = **gestalt/proximity** fail (B§3–4).
4. **Is it cramped or empty?** → **white space** fail (B§6).
Run it on the **real, populated screen**, not the empty template.[^nnghier]

## Severity cheat-sheet

| Band | Visual examples | Typical Nielsen rating |
| --- | --- | --- |
| **Blocker/Catastrophe** | No entry point; figure-ground ambiguity on a modal; layout clips at 200% zoom | 4 |
| **Major** | Two equal primary CTAs; broken alignment fragmenting the page; functional inconsistency; F-symptom over body copy; tipping layout | 3 |
| **Minor** | Mild rhythm break; one off-scale value; slight intentional-looking lean | 2 |
| **Cosmetic** | ~2px misalignment; slight weight inconsistency; off-brand shade (but watch *priority* if homepage) | 1 |

Rate on **frequency × impact × persistence**; track **priority** (business urgency) separately; trust
the **mean of ~3 raters**, not one.[^nngsev][^betterqa]

---

## References

[^nngsev]: Nielsen Norman Group, "How to Rate the Severity of Usability Problems" — the 0–4 scale, frequency/impact/persistence, single-rater unreliability, ~3-rater mean, "don't use fancy scales." nngroup.com/articles/how-to-rate-the-severity-of-usability-problems/ (docs)
[^conn]: Connor & Irizarry, *Discussing Design* (book preview PDF) — feedback is a nonspecific reaction; critique is goal-framed analysis; requester provides scope/goals. api.pageplace.de preview-9781491902387 (book)
[^feldman]: Feldman model of art criticism (Describe/Analyze/Interpret/Judge); judgment withheld until the end. gisd.org / humankinetics.com excerpt (docs/book)
[^nngcrit]: Nielsen Norman Group, "Design Critiques: Encourage a Positive Culture to Improve Products" — standalone critique vs design review; facilitation definition; reformulate opinionated/directive comments. nngroup.com/articles/design-critiques/ (docs)
[^betterqa]: BetterQA, "Bug Priority vs Severity" — severity (impact) vs priority (urgency) are independent; cosmetic-but-high-priority brand example. betterqa.co/bug-priority-vs-severity-levels/ (blog)
[^uxsl]: bugnet.io, "Organize Bug Reports by Severity" — severity vs priority axes; S3 visual glitch / S4 cosmetic. bugnet.io/blog (blog)
[^nnghier]: Nielsen Norman Group, "Visual Hierarchy in UX: Definition" — levers, ≤3 sizes/contrast levels, squint test, real-content caveat. nngroup.com/articles/visual-hierarchy-ux-definition/ (docs)
[^digpolo]: DigitalPolo, "Visual Hierarchy in Design" — "everything equal = everything unimportant"; contrast as umbrella; squint test. digitalpolo.com/visual-hierarchy-in-design/ (blog)
[^artstyle]: Art of Styleframe, "Visual Hierarchy UI Design Principles" — 7 levers, squint test, ~80% above-fold. artofstyleframe.com (blog)
[^bradley]: Steven Bradley, Smashing Magazine, "Design Principles: Dominance, Focal Points and Hierarchy" / vanseodesign — 3-level ranking; entry point reduces cognitive load. smashingmagazine.com/2015/02/ (blog)
[^designyourway]: Designyourway, "What Is Emphasis in Graphic Design" — emphasis vs hierarchy; dominance/subordination/isolation; "emphasis is the first beat of hierarchy." designyourway.net (blog)
[^madegood]: MadeGood Designs, "Emphasis in Graphic Design" — focal point = entry point; one primary; split attention if two compete. madegooddesigns.com/emphasis-in-graphic-design/ (blog)
[^smashdom]: Smashing Magazine (Bradley), dominance/focal-points — emphasis relative; "ideally one dominant element"; co-dominance distracts. smashingmagazine.com/2015/02/ (blog)
[^vonrestorff]: UXScan, "Von Restorff Effect" — isolation effect; "if everything is highlighted nothing is"; pair color with shape for a11y. uxscan.ai/learn/von-restorff-effect (blog)
[^cta]: CTA-failure cluster (auditbuffet.com, levri.ai, homepagedream, bivisee.com) — two equal filled buttons cause choice paralysis; ghost/outline subordination fix; A/B +34% clicks/−12% conversions when secondary not subordinated. (blog)
[^ixdfgestalt]: Interaction Design Foundation, gestalt principles in design — proximity (strongest, can override similarity), similarity, common region, closure, continuity, figure-ground, common fate. interaction-design.org/literature/topics/gestalt-principles (docs)
[^nnggestalt]: Nielsen Norman Group on gestalt/proximity in UI grouping. nngroup.com (docs)
[^nngproximity]: NN/g, proximity principle — related items grouped by spacing; the most common grouping defect. nngroup.com (docs)
[^lidwellfg]: Lidwell, Holden & Butler, *Universal Principles of Design*, "Figure-Ground Relationship" — figure-ground stability underlies emphasis/overlays. oreilly.com/library/view/universal-principles (book)
[^smashprinc]: Smashing Magazine, "10 Principles of Effective Web Design" — macro vs micro whitespace; prefer space to dividers; rhythm/unity. smashingmagazine.com/2008/01/ (blog)
[^crapttd]: Robin Williams, *The Non-Designer's Design Book* — Contrast/Repetition/Alignment/Proximity definitions and amateur tells. (book)
[^crapsmash]: CRAP principles summaries (design-education sources) corroborating Williams. (blog)
[^balpolo]: DigitalPolo, "Balance in Design" — symmetric/asymmetric/radial; asymmetry harder & tips; squint for weight bunching; failure modes. digitalpolo.com/balance-in-design/ (blog)
[^balramotion]: Ramotion + MadeGood + NN/g balance video — symmetric=formal/trust, asymmetric=modern, radial=center focus, mosaic=no focal point; visual-weight factors. ramotion.com/blog/principles-of-design-balance/ (blog/docs)
[^balkittl]: Kittl, symmetrical/asymmetrical/crystallographic balance — 4 types; asymmetry = different elements, equal weight. kittl.com (blog)
[^balkose]: Samet Köseoğlu, asymmetrical balance guide — weight by texture/distance; mobile-reflow risk; counterbalance with multiple small elements/negative space. sametkoseoglu.com (blog)
[^balsmash]: Smashing Magazine, "Compositional Balance: Symmetry & Asymmetry" — asymmetry "harder to achieve," more dynamic; combine. smashingmagazine.com/2015/06/ (blog)
[^nngmin]: Nielsen Norman Group, "Characteristics of Minimalism" — negative space directs attention; too-much-space fold/hierarchy caution. nngroup.com/articles/characteristics-minimalism/ (docs)
[^ixdfspace]: Interaction Design Foundation, negative space — macro/micro, active/passive; whitespace → perceived luxury. interaction-design.org/literature/topics/negative-space (docs)
[^nngclutter]: NN/g, "Aesthetic-Minimalist Design" (Heuristic #8) — signal-to-noise; every extra unit competes; remove, don't shrink. nngroup.com/articles/aesthetic-minimalist-design/ (docs)
[^portent]: Portent, "Too Much White Space" — disconnects elements; "illusion of completeness." portent.com/blog (blog)
[^lineheight]: WCAG 2.2 (W3C) + readability guidance — body line-height ≥1.5 baseline, 1.6–1.7 long lines. w3.org/WAI/WCAG22 (docs/standard)
[^measure]: UXPin, "Optimal Line Length for Readability" — 50–75 char measure (~66 ideal). uxpin.com/studio/blog/optimal-line-length-for-readability/ (blog)
[^typography]: Practical Typography (via summaries) — body tracking ±0.01em; ALL-CAPS +5–12%; tight tracking hurts paragraphs. practicaltypography.com (docs)
[^scalemic]: MakeItClear, "Typography Systems" — modular scale + 8pt/4pt baseline; round to 4/8; contrast bands. makeitclear.com/typography-systems/ (blog)
[^modscale]: Modular Scale (modularscale.com) — ratio list incl. golden section; "use a scale like a ruler"; Tim Brown origin. modularscale.com (tool)
[^ala]: Tim Brown, "More Meaningful Typography" (A List Apart) — modular scale definition; works with grids/responsive. alistapart.com/article/more-meaningful-typography/ (blog/authority)
[^scalecalc]: Calculatey type-scale calculator — `base×ratio^step`; Major Third/Perfect Fourth most common. calculatey.com/type-scale-calculator/ (tool)
[^brainy]: Brainy.ink, modular type-scale guide — ratio→density table (Tailwind/Stripe/Material); pick ratio for density. brainy.ink (blog)
[^scale8pt]: GridMakerPro, "12-column / baseline / 8pt" — 8pt grid values; three grids must coordinate as 8px multiples; "approximately right" drift. gridmakerpro.com (blog)
[^android]: Android Developers / Material Design (m1) — 8dp grid for layout, 4dp for icons/type; margins/gutters. developer.android.com/design / m1.material.io (docs)
[^krycho]: Chris Krycho, "Modular Scales: Fantastic, But Don't Overdo It" — 36-value bloat; limit the subset. v3.chriskrycho.com (blog, negation)
[^codexical]: Codexical, "Design System Typography Failures" — fixed-px scales clip at 200% zoom; rem + fluid fix. codexical.com (blog, negation)
[^fontfyi]: FontFyi glossary, modular scale — when to abandon; scale as aid not constraint. fontfyi.com/glossary/modular-scale/ (blog)
[^goldenpmc]: Peer-reviewed review (PMC, 2024), golden ratio — evidence "without foundation"; no Fibonacci–beauty link. pmc.ncbi.nlm.nih.gov/articles/PMC10792139/ (paper)
[^goldendevlin]: Keith Devlin / MAA + Markowsky 1992 ("Misconceptions About the Golden Ratio") — "not a shred of evidence"; Markowsky measures the Parthenon width-to-height at ~1.71, stricter readings ~2.25 — neither is φ. (authority)
[^goldenmade]: MadeGood Designs, golden ratio — balanced view; strongest critique = retrofitting; φ as legit typographic ratio. madegooddesigns.com/golden-ratio-design/ (blog)
[^lidwellunity]: Lidwell et al., *Universal Principles of Design* — unity/variety, rhythm. (book)
[^nngscan]: Nielsen Norman Group, "Text Scanning Patterns: Eyetracking" — F/spotted/layer-cake/commitment taxonomy; layer-cake "most effective." nngroup.com/articles/text-scanning-patterns-eyetracking/ (docs)
[^nngf]: NN/g, "F-Shaped Pattern for Reading Web Content" — F applies to content blocks not whole pages; F is a *symptom* of unguided text, not a best practice. nngroup.com/articles/f-shaped-pattern-reading-web-content/ (docs)
[^nngfdisc]: NN/g, "F-Shaped Pattern… Discovered" — original 2006 study (232 users); three F components; first-two-paragraphs implication. nngroup.com (docs)
[^nngbypass]: NN/g, "The First 2 Words: A Signal for Scanning" — bypassing; ~2 words per item; front-load keywords; "Introducing" ~85% failure. nngroup.com/articles/first-2-words-a-signal-for-scanning/ (docs)
[^zpattern]: Creative Bloq + Instapage — Z-pattern definition/placement; heuristic (not eye-tracking), sparse-page use, CTA at terminal. creativebloq.com / instapage.com (blog)
[^lidwellgut]: Lidwell et al., *Universal Principles of Design*, "Gutenberg Diagram" — four quadrants, reading gravity, homogeneous-info scope, terminal-area placement. oreilly.com/library/view/universal-principles (book)
[^gutituon]: ITU Online, Gutenberg principle — "a framework, not a law"; overridden by content/culture/device. ituonline.com (blog, negation)
[^nngconsist]: NN/g, "Consistency and Standards in UX" — Heuristic #4; internal vs external; four visual layers; when-to-break caveat; "foolish consistency." nngroup.com/articles/consistency-and-standards/ (docs)
[^jakobslaw]: Laws of UX, "Jakob's Law" — users transfer mental models from other sites; minimize discord when changing. lawsofux.com/jakobs-law/ (docs)
[^iconincon]: Prakeerth, icon inconsistency — mixed icon styles; multiple icons for one action. blog.prakeerth.com (blog)
[^tokens]: Material Design design tokens (m3) — global/semantic/component token hierarchy to prevent hard-coded drift. m3.material.io/foundations/design-tokens/overview (docs)
[^tyranny]: UX Mag "Tyranny of Consistency" + TNW — "foolish consistency" (Google icons), US Airways logo collapse; bend-not-break. uxmag.com (blog, negation)