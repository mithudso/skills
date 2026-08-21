<!-- hub-reference-banner -->
> **Reference file — part of the `frontend-ui` hub.** The methodology half of
> `visual-design-principles-and-critique` (the principles rubric is the companion file
> `references/visual-design-principles-and-critique.md`). Sibling topics in this family are
> reference files under the hubs (`frontend-ui`) — **not** standalone skills. Ignore any
> "use the X skill" pointers that name a bare sibling skill; load that topic's
> `references/<name>.md` from the owning hub.

---
name: visual-design-critique-methodology
description: >-
  Design-critique METHODOLOGY — how to run and give design critiques; companion to the
  visual-design-principles-and-critique rubric. Structured crit formats (I-Like-I-Wish-What-If,
  Describe/Interpret/Evaluate, Rose/Bud/Thorn, design studio/charrette), Connor & Irizarry
  "Discussing Design", critique-vs-feedback-vs-review, observation-before-judgment, giving
  actionable-but-non-prescriptive feedback, severity & design QA passes, running a crit, and
  critique anti-patterns (compliment sandwich, HiPPO, bikeshedding, design-by-committee, vague
  praise). TRIGGER: run/facilitate or structure a design critique; how to give design feedback;
  "design crit format"; separate observation from judgment; design QA / spec-parity pass; critique
  anti-patterns. SKIP: the visual-design principles themselves — hierarchy/gestalt/balance/scale
  (use visual-design-principles-and-critique); Nielsen heuristic eval or Laws of UX (use
  usability-heuristics-laws-of-ux); producing a design (use web-design / ui-ux-pro-max).
category: developer
version: "1.0.0"
updated: "2026-06-16"
keywords:
  - design critique methodology
  - design critique format
  - I like I wish what if
  - describe interpret evaluate
  - actionable feedback
  - design QA
  - critique anti-patterns
  - observation vs judgment
  - Discussing Design
  - design feedback
tags:
  - developer
  - frontend
  - design
  - ux
  - critique
  - feedback
whenToUse:
  - "running, facilitating, or structuring a design critique session"
  - "how to give actionable, non-prescriptive feedback to a designer"
  - "choosing a critique format (I-like-I-wish-what-if, describe/interpret/evaluate, rose/bud/thorn)"
  - "separating observation from judgment in design feedback"
  - "running a design QA / spec-parity pass or rating design-issue severity"
  - "diagnosing why a critique meeting goes wrong (sandwich, HiPPO, bikeshedding)"
whenNotToUse:
  - "the visual-design principles themselves (hierarchy/gestalt/balance/scale) — use visual-design-principles-and-critique"
  - "Nielsen heuristic evaluation or Laws of UX as the lens — use usability-heuristics-laws-of-ux"
  - "producing a NEW design from scratch — use web-design / ui-ux-pro-max"
related_skills:
  - frontend-ui
  - visual-design-principles-and-critique
  - usability-heuristics-laws-of-ux
---

# Design-Critique Methodology

How professional critiques are run and how feedback is delivered. The companion to
`references/visual-design-principles-and-critique.md` (the *principles* rubric — what to look for);
this file is the *methodology* — how to run the crit and deliver the findings. The whole discipline
rests on one move: **separate observation from judgment, and tie everything to objectives, not
taste.**[^conn][^feldman]

## Contents

- Critique vs. feedback vs. review
- Structured critique formats
- The frameworks (Connor & Irizarry, *Discussing Design*)
- Giving actionable-but-non-prescriptive feedback
- Severity & design QA passes
- Running a critique
- Critique anti-patterns
- References

## Critique vs. feedback vs. review (get this straight first)

- **Feedback** = a gut reaction; nonspecific; not measured against goals.[^uiebrain]
- **Critique** = deliberate *analysis* of a design against its objectives/personas/scenarios; its
  purpose is to make the work **better**, not to approve it. ("How can this be better?")[^uiebrain][^nngcrit]
- **Review** = an *evaluative gate* deciding readiness for the next stage (go/no-go). Different mindset,
  different (decision-maker) participants, different timing. ("Is this good enough to proceed?")[^crit-vs-rev]
- **Anti-pattern:** combining critique and review in one meeting → "the worst of both worlds"; a review
  where everyone nods until a senior leader weighs in has degraded into a politics-driven approval.[^crit-vs-rev]

## Structured critique formats

- **"I Like / I Wish / What If"** (Stanford d.school / IDEO). Three channels: *I like…* (specific
  things working, to preserve), *I wish…* (criticism reframed as aspiration), *what if…* (generative
  alternatives, **not** prescriptions). The "I" framing signals opinion-not-fact and lowers
  defensiveness. **Receiver protocol:** the receiver *just listens* in-session — "thank you" suffices;
  discussion is deferred. Best for group concept reviews and mixed-discipline psychological safety.[^ideo][^dschool]
- **"Describe / Interpret / Evaluate"** (Feldman's art-criticism method; full form Describe → Analyze →
  Interpret → Judge). Sequential and **order-enforcing**: inventory what's literally present →
  relationships among elements → meaning → informed verdict. The point: **judgment is withheld until the
  end** — premature judgment "is neither informed nor critical but simply an opinion." This is the
  canonical observation-before-judgment sequence; adapt to screens (describe layout → analyze hierarchy/
  relationships → interpret what it communicates → judge against goals). *Limit:* can be over-linear/
  formalist — add a "context" step to compensate.[^feldman][^feldmancrit]
- **"Rose / Bud / Thorn."** Rose = what's working; Thorn = friction/problem; **Bud = latent
  opportunity worth developing.** The separate Bud bucket distinguishes "needs fixing" from "worth
  growing" — useful for prioritizing.[^rtb]
- **Design studio / charrette.** Time-boxed collaborative sketching (≈5-min solo sketch → ≈2-min present
  → ≈1-min questions), then **Sketch → Present & Critique → Converge → Prioritize**. The critique step is
  the explicit evaluation method inside the studio; a skilled facilitator manages diverge/converge.[^charrette][^studio]

## The frameworks (Connor & Irizarry, *Discussing Design*)

The canonical text. Load-bearing principles for a critique-rubric:[^connbook][^connint]

- **Critique ≠ reaction.** Critique is intense, deliberate analysis — not a one-shot thumbs-up. You
  don't shoot for "approval and signoff."
- **You get the critique you ask for.** The requester must decide what they want — a *reaction*, a
  *direction*, or deep *analysis/critique* — and grant permission before anyone critiques.
- **Frame around objectives, not taste.** Agreed goals, design principles, personas, and scenarios "keep
  conversations from going off into personal preferences." **If goals aren't agreed, hold the critique**
  and define them first. (NN/g: "without agreed-upon design objectives, any feedback is subjective and
  baseless.")[^nngcrit]
- **The facilitator** "consciously, balancedly manages conversations toward a conclusion" — sets ground
  rules, balances voices, calls out dominators, and **reformulates opinionated/directive comments back
  to goals** (turn "This is too red!" into "does this color help the design meet its objective?").[^nngcrit]
- **Don't problem-solve mid-critique.** Critique is analysis; switching to solutions derails it. The
  designer gathers insights in-session, then decides *later* against goals/personas/scenarios.

## Giving actionable-but-non-prescriptive feedback

The three-property test: feedback must be **specific, actionable, and not prescriptive.**[^uifeed]

- **Specific** — references a real element/flow. ✗ "the onboarding feels long" → ✓ "steps 3 and 4 both
  ask for account details — could they combine?"
- **Actionable** — the designer can act on it.
- **Not prescriptive** — flag the *problem*, not the *solution*. The canonical contrast: ✗ "make X red" /
  "make the header smaller" (a solution) → ✓ "X feels buried" / "the header competes with the primary
  CTA" (a problem). Designers solve problems better than they execute mandates.[^uifeed][^steffen]

**The translation habit:** don't stop at the reaction ("this is boring") — name *what you noticed → why
it matters (the effect) → where to go next.* Connect observations to *outcomes* ("the blue reads too
corporate for a launch targeting a younger audience").[^steffen][^manypixels]

**When a solution IS appropriate,** mark it optional ("one option would be to reduce the header size,
but you might have another approach") and distinguish a **constraint** (must) from a **preference**
(describe the problem, let the designer solve).[^uifeed]

**Ask questions instead of directives** — "what led you to put the nav at the top?" surfaces the
rationale (maybe mobile-first or prior testing) and often changes what feedback is even
relevant.[^uifeed]

**Giver/receiver dynamics.** *Giver:* ensure feedback is invited, specific, goal-tied; lead with intent
to help. *Receiver:* "you are not your work" — pause before responding, **don't defend, get curious**;
the first job is to understand, not respond; take notes (signals it's captured, not to be defended live);
**collect first, don't live-edit.** Distinguish **directional** feedback ("something's off" — take the
signal) from **prescriptive** ("do this instead" — you may decline the fix while honoring the signal);
"advocacy is not defensiveness."[^zhuo][^dontdefend]

## Severity & design QA passes

Use **Nielsen's 0–4 severity scale** for every finding (0 not-a-problem, 1 cosmetic, 2 minor, 3 major,
4 catastrophe — a function of frequency × impact × persistence), tracking *priority* separately from
severity. The principles rubric's Part A carries the full scale; the single-rater-unreliability and
"trust the mean of ~3 raters" caveats apply here too.[^nngsev] For visual/UI defects, practitioner QA
splits the same idea into a severity ladder — **Blocker → Critical → Major → Minor → Cosmetic** — where
a ~2px misalignment or slight font-weight inconsistency is Cosmetic.[^betterqa][^qa-ladder]

A **design QA** (a.k.a. design-vs-build / spec-parity pass) compares the *live implementation* against
the *approved design*: typography (family/size/weight/line-height/tracking), color (fill/text/border/
icon), spacing (margins/padding/gaps), sizing, responsive behavior at breakpoints (e.g. 375/768/desktop;
44×44px touch targets; no horizontal scroll), **interactive states** (hover/focus/active/loading/error/
empty — missing hover & loading states are the most common gaps), cross-browser, and accessibility
(4.5:1 contrast, heading hierarchy, alt text, keyboard nav, form labels).[^designqa][^figmaqa] The
strongest practice makes a *small* set (3–5) of design checks **hard blockers** on merge (visual-
regression diff, WCAG-AA violation, token bypass) and everything else warn-only, complemented by
structured human heuristic review.[^qagate]

## Running a critique

- **Set context & objectives up front** (the biggest lever): distribute scope + agenda + goals + rules
  *before* the session; send materials ahead; **state the product goals but not how the design intends
  to meet them** (so attendees form independent insight); state what kind of feedback you
  need.[^nngderail][^nasdaq]
- **Who attends:** crit can/should be cross-disciplinary (designers, peers, engineers, PMs, SMEs); a
  *review's* participants differ (decision-makers). Watch for the wrong people or dominating
  personalities.[^nngcrit][^nngderail]
- **Sync vs. async:** **synchronous** for decision points, complex trade-offs, emotionally charged work,
  or converging divergent directions (risk: loudest voice, groupthink). **Asynchronous** for detail/
  polish, deep individual analysis without groupthink, time zones, or settled patterns (strengths:
  deeper considered feedback, self-documenting; risk: ignored, no dialogue). Neither is inherently
  better.[^syncasync][^alist]
- **Time-box** (inherited from the charrette rhythm) and use a **parking lot** for off-scope items.[^nngderail]
- **Close the loop** — the session only surfaces feedback; **follow-up builds trust:** show what changed
  and why, *and* be transparent about what did *not* change and the reasoning ("more credible than
  sharing only what did"). Attribute influence ("you flagged this last crit").[^afterderail]

## Critique anti-patterns (what a broken crit looks like)

- **Compliment / critique sandwich.** Discouraged: it manages the *giver's* discomfort, trains
  recipients to brace at any praise, and buries the actual message via primacy/recency. Empirical
  backing: Kluger & DeNisi's meta-analysis found ~one-third of feedback interventions *decreased*
  performance. Better: treat praise and criticism as separate, specific acts; address the hard thing
  directly; feedback about the task, not the person.[^sandwich][^grant]
- **HiPPO** (Highest-Paid Person's Opinion) — the most senior voice overrides data/research/expertise.
  Counter: de-personalize with data; "delay, gather data, respond via priority."[^hippo]
- **Bikeshedding** (Parkinson's Law of Triviality) — the group over-invests in trivial issues everyone
  feels qualified to opine on while complex ones get nodded through. Counter: a written pre-read
  equalizes expertise; **name the pattern out loud**.[^bikeshed]
- **Design-by-committee** — decisions by a large group with no clear owner; compromise dilutes the
  design. Counter: clear ownership; challenge ideas with evidence, not equal-weighting of all
  opinions.[^committee]
- **Vague praise** ("looks great," "make it pop") — teaches nothing; the receiver can't tell what to
  keep. "Praise is approval; feedback is information." If feedback could be pasted onto any work, it's
  too vague.[^vague]
- **Problem-solving mid-critique / live-editing** — switches analysis into premature solutioning and the
  giver isn't fully heard.[^connbook]

---

## References

[^conn]: Connor & Irizarry, *Discussing Design* (book preview PDF) — feedback is a nonspecific reaction; critique is goal-framed analysis; requester provides scope/goals. api.pageplace.de preview-9781491902387 (book)
[^feldman]: Feldman model of art criticism (Describe/Analyze/Interpret/Judge); judgment withheld until the end. gisd.org / humankinetics.com excerpt (docs/book)
[^nngsev]: Nielsen Norman Group, "How to Rate the Severity of Usability Problems" — the 0–4 scale, frequency/impact/persistence, single-rater unreliability, ~3-rater mean. nngroup.com/articles/how-to-rate-the-severity-of-usability-problems/ (docs)
[^nngcrit]: Nielsen Norman Group, "Design Critiques: Encourage a Positive Culture to Improve Products" — standalone critique vs design review; facilitation definition; reformulate opinionated/directive comments. nngroup.com/articles/design-critiques/ (docs)
[^uiebrain]: Connor & Irizarry interview (UIE/Brainsparks) — feedback=gut reaction; critique=specific/goal-framed/not signoff. archive.uie.com (authority)
[^crit-vs-rev]: TheCrit.co + uiguides — critique (improve) vs review (decide); "worst of both worlds"; political vs actionable. thecrit.co/resources/design-critique-vs-design-review (blog)
[^ideo]: IDEO, "Build Your Creative Confidence: I Like, I Wish" — framing; receiver-just-listens rule. ideo.com/blog (authority)
[^dschool]: Stanford d.school, "How to Give Feedback" — benevolent critique; separate person from work; focus on the goal. dschool.stanford.edu/tools/how-to-give-feedback (authority)
[^feldmancrit]: Jerwood / VCU scholarly critique of Feldman — formalist, over-linear; add context step. jerwoodvisualarts.org / scholarscompass.vcu.edu (paper, negation)
[^rtb]: Rose/Bud/Thorn (CU Boulder, Tamarack, OpenPracticeLibrary) — Rose/Thorn/Bud definitions; Bud = opportunity; design-thinking origin. colorado.edu/researchinnovation/rose-bud-thorn (docs)
[^charrette]: NN/g, "Design Charrettes" — 5-min sketch / 2-min present / 1-min questions; buy-in. nngroup.com/articles/design-charrettes/ (docs)
[^studio]: NN/g, "Facilitating a Design Studio Workshop" — Sketch→Present&Critique→Converge→Prioritize. nngroup.com/articles/facilitating-design-studio-workshop/ (docs)
[^connbook]: Connor & Irizarry, *Discussing Design* (book preview + author decks) — critique≠reaction; you get the critique you ask for; frame around objectives; don't problem-solve mid-critique. api.pageplace.de (book)
[^connint]: Connor & Irizarry interviews (Conversation Factory, Center Centre/Spool) — reaction/direction/analysis; hold critique if no agreed goals; facilitator balances voices. centercentre.com / theconversationfactory.com (authority)
[^uifeed]: uiguides, "How to Give Design Feedback" — specific/actionable/not-prescriptive; ask questions not directives; suggestion-not-directive exception; constraint vs preference. uiguides.com/guides/how-to-give-design-feedback (blog)
[^steffen]: Nicole Steffen, "Good Design Feedback" — emotional→actionable translation; "font size vs CTA prominence"; vague/taste/reactive trio. nicolesteffen.com (blog)
[^manypixels]: ManyPixels, design feedback — connect observation to outcome; name problem not solution; tell designer what's working. manypixels.co/blog (blog)
[^zhuo]: Julie Zhuo "Taking Feedback Impersonally" + summaries — "you are not your work"; growth mindset; collect first, don't live-edit. (authority)
[^dontdefend]: Austin Knight, "Don't Defend Your Work" + design-bootcamp — get curious not defensive; directional vs prescriptive; "advocacy is not defensiveness." austinknight.com (blog)
[^betterqa]: BetterQA, "Bug Priority vs Severity" — severity (impact) vs priority (urgency) are independent; cosmetic-but-high-priority brand example. betterqa.co/bug-priority-vs-severity-levels/ (blog)
[^qa-ladder]: QA severity ladder (code-note / Airbrake) — S1 Blocker → S5 Cosmetic with visual-defect examples. (blog)
[^designqa]: OverlayQA + Eleken, design-QA handoff checklist — live page vs design property-by-property; 375/768/desktop, 44×44 targets, 4.5:1; three-phase QA. overlayqa.com/blog/design-qa-agencies/ (blog)
[^figmaqa]: UIProbe, "Figma Design QA Checklist" — interactive states (hover/focus/loading/error/empty); self-check before review. uiprobe.io (blog)
[^qagate]: UX/UI Principles, "Design QA as Release Gate" — hard + automated + human heuristic; small block-on-fail set; cites Moran & Gordon. uxuiprinciples.com (blog)
[^nngderail]: NN/g, "Derailed Design Critiques" — set feedback expectations up front; derailers (no aligned goals, dominating personalities); parking lot. nngroup.com/articles/derailed-design-critiques/ (docs)
[^nasdaq]: Aaron Irizarry (Nasdaq Design) — don't problem-solve mid-critique; send materials + goals (not how-to-meet-them) ahead. medium.com/nasdaq-design (authority)
[^syncasync]: NN/g, "Synchronous vs Asynchronous Ideation" + thecrit.co matrix — sync=faster/team-building, async=anytime/fewer-timezone-issues; neither inherently better; risks. nngroup.com/articles/synchronous-asynchronous-ideation/ (docs)
[^alist]: A List Apart + Atlassian, async design critique — async lets you refine for clarity+actionability; self-documenting decisions. alistapart.com/article/async-design-critique-giving-feedback/ (blog/authority)
[^afterderail]: NN/g, "After the Design Critique" — closing the loop builds trust; be transparent about what didn't change; attribute influence. nngroup.com/articles/after-design-critique/ (docs)
[^sandwich]: Compliment-sandwich critique cluster (Psychology Today, Radical Candor, Cambridge/ScienceDirect peer-reviewed) — manages giver discomfort; trains discounting of praise; buries message; Kluger & DeNisi 1996 meta-analysis (~1/3 of feedback interventions decreased performance). (paper/blog, negation)
[^grant]: Adam Grant, "Stop Serving the Compliment Sandwich" — positives drown out negatives (primacy-recency); helps giver not receiver. adamgrant.substack.com (authority, negation)
[^hippo]: HiPPO anti-pattern (uxdesign.cc, auvik) — most senior voice overrides data; de-personalize with data; delay/gather/respond. uxdesign.cc (blog)
[^bikeshed]: Bikeshedding / Parkinson's Law of Triviality (tianpan.co) — over-invest in trivial; written pre-read equalizes expertise; name the pattern. tianpan.co/notes/19-bikeshedding (blog)
[^committee]: Design-by-committee anti-pattern — no clear owner; compromise dilutes; co-exists with HiPPO. architecture-evolution.blogspot.com (blog)
[^vague]: Vague-praise cluster (Psychology Today, julienflorkin) — "praise is approval, feedback is information"; copy-paste test for vagueness. (blog)
