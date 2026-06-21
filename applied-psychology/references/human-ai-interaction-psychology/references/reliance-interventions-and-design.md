# Reliance Interventions & Design for Appropriate Reliance

> Provenance: authored for the `human-ai-interaction-psychology` skill via /dr.
> Primary sources: Parasuraman & Manzey (2010); Bansal et al. (2021, CHI);
> Buçinca, Malaya & Gajos (2021, CSCW); Mori (1970); plus the 2023-2025 HCI
> evidence base on miscalibrated confidence, confidence/self-confidence
> alignment, and selective explanations (CHI 2024-2025; Microsoft Research 2024
> appropriate-reliance synthesis).

## Automation bias & complacency — the mechanisms to design against

**Automation bias** = the tendency to over-rely on automated cues as a heuristic
replacement for vigilant information-seeking. With an *imperfect* aid it produces:

- **Commission errors** — acting on a *wrong* automated directive without
  cross-checking other available evidence.
- **Omission errors** — missing an event because the automation didn't flag it
  *and* the human wasn't independently monitoring.

**Automation complacency** = reduced monitoring of automation under load.
Parasuraman & Manzey's *attentional integration* (2010): under multitask
conditions, attention is reallocated away from the automated task; the operator
"assumes it's handled." Bias and complacency share this attentional root.

**Findings that constrain design:**
- Appears in **experts and novices** alike.
- **Not reliably removed by training or warnings** alone — you must change the
  task structure, not just instruct harder.
- Occurs in **teams**, and *redundancy can backfire*: diffusion of monitoring
  responsibility means "two reviewers" may each assume the other (or the model)
  caught it.
- Worsens with **high automation reliability** (the better it usually is, the
  less people check — so very good systems have a complacency tail) and with
  **high workload**.

**Design implication:** a human-review step is a *control only if it forces
independent engagement.* Passive "please review" review decays into
rubber-stamping under exactly the conditions (load, high reliability) where
errors slip through.

## Why explanations & confidence displays backfire (the evidence)

The intuition that transparency yields calibration is largely **not supported**;
several lines of evidence show it can make reliance *worse*:

- **Explanations raise agreement, not discrimination (Bansal et al., 2021,
  "Does the Whole Exceed Its Parts?"):** AI explanations increased the
  probability that people accepted the AI recommendation **whether it was right
  or wrong**. Explanations made the advice *persuasive*, not *evaluable* — over-
  reliance rose; complementary team performance did not.
- **Plausible-but-wrong explanations are especially dangerous:** a coherent
  rationale for an incorrect answer is more convincing than no rationale,
  precisely because it satisfies the heuristic search for "a reason."
- **Confidence helps only if calibrated:** well-calibrated confidence can support
  appropriate reliance; **miscalibrated** confidence (high confidence on wrong
  answers) *degrades* trust calibration and decision quality (CHI 2024 work on
  miscalibrated AI confidence).
- **Confidence anchors the human's own self-confidence:** recent work (CHI 2025)
  finds people's *self-confidence shifts toward the AI's stated confidence*
  whenever they see it — even absent collaboration — which can miscalibrate the
  human's uncertainty without improving their actual ability, undermining
  appropriate reliance.

**Dual-process explanation (why):** people seldom engage System-2 analysis on
each recommendation; they form a heuristic ("usually right + has a reason →
accept"). Explanations and confidence *feed* the heuristic instead of
interrupting it.

**Takeaway:** "we added explanations / confidence scores" is **not** evidence of
appropriate reliance and may indicate the opposite. Verify behaviorally (does
override track error?). Ship confidence numbers only after validating they're
calibrated.

## Cognitive forcing functions (the interventions that work)

*Cognitive forcing functions* (Buçinca, Malaya & Gajos, 2021, "To Trust or to
Think") = design interventions that introduce **friction at decision time** to
trigger analytical (System-2) engagement:

- **Commit-first / decide-before-reveal:** require the human to make an
  independent judgment *before* the AI suggestion is shown (prevents anchoring).
- **On-demand AI:** withhold the recommendation until the user explicitly
  requests it (deliberate, not passive, consumption).
- **Slow-down / wait:** insert a brief forced delay before the AI answer appears.
- **Surface disagreement & uncertainty** as an explicit prompt to slow down on
  the cases most likely to be wrong.

**Result:** these reduced over-reliance on *incorrect* AI more than
explanation-only designs.

**Costs and caveats (do not over-apply):**
- Forcing functions impose **cognitive effort** and are frequently **disliked**;
  satisfaction can drop.
- Benefit **interacts with the person** — e.g., higher *Need for Cognition*
  users engage and benefit; others may disengage or route around the friction.
- Therefore **target high-stakes / likely-error cases**, not every interaction.

**Adjacent levers:** onboarding/tutorials that teach the model's *error
boundaries*; **selective explanations** shown on likely-error cases rather than
universally; **adjustable outputs** (the Dietvorst control lever — also reduces
aversion).

## Human-AI complementarity (Complementary Team Performance)

- **CTP** = human+AI together outperform *both* the best human alone *and* the AI
  alone.
- **Condition for it:** the two make **different errors** and each defers on the
  cases the other handles better — i.e., genuine appropriate reliance plus the
  human contributing information the model lacks ("unobservables," context, side
  knowledge).
- **Reality check (Bansal et al. and the broader literature):** CTP is **rare by
  default.** AI-assisted humans frequently do *worse than the AI alone* because
  they over-rely on wrong answers and override right ones. Adding explanations or
  confidence does **not** reliably yield CTP.

**Design implication:** put the human where they hold an **information edge** the
model can't have (recent events, account/political context, tacit knowledge),
not as a generic reviewer of things the model already does better. Verify CTP
empirically; never assume "human + AI > either."

## Anthropomorphism, persona & the uncanny valley

- **Anthropomorphic cues** (name, persona, warmth, avatar) can raise initial
  trust/likability — but the effect is **mediated by perceived empathy and
  interaction quality**, not the surface cues alone. Explicit "I am an AI"
  identity cues can override surface human-likeness.
- **Uncanny valley (Mori, 1970):** as an agent nears — but doesn't reach —
  human-likeness, affinity drops sharply into eeriness. An **"uncanny valley of
  trust"**: human enough to raise competence expectations it can't meet → user
  skepticism, over-promise/under-deliver disappointment.
- **Persona inflates over-trust risk:** a warm, fluent, confident persona is a
  strong (misleading) competence cue — fluency reads as competence — pushing
  users toward over-reliance regardless of correctness. This is the persona
  analog of the plausible-but-wrong explanation.

**Design implication:** match persona confidence to *validated* capability. For
high-stakes tools prefer a sober, capability-honest voice over a maximally
human, maximally confident one. Don't manufacture trust the system can't justify.

## Consolidated design & coaching checklist

1. **Target appropriate reliance, measured behaviorally** (override tracks
   error). Instrument over-reliance and under-reliance rates; don't optimize a
   "trust" survey number or raw agreement.
2. **Set honest expectations before the first error** (error boundaries,
   likely-wrong cases) — inoculates against both aversion and blind over-trust.
3. **Show confidence only if calibrated; communicate uncertainty honestly**
   (epistemic markers help over time; miscalibrated confidence is worse than
   none).
4. **Don't expect explanations to create skepticism** — they usually raise
   acceptance. If shipped, pair with friction; prefer *selective* explanations on
   likely-error cases; verify behavior.
5. **Engineer friction where stakes are high** (commit-first, on-demand reveal,
   forced confrontation of disagreement) — and reserve it; friction has a cost
   and is unevenly effective across users.
6. **Give users control/adjustability** over outputs (restores reliance after
   visible errors).
7. **Place the human where they have an information edge,** not as a generic
   reviewer — that's where complementarity lives.
8. **Match persona confidence to validated capability;** avoid manufacturing
   over-trust via warmth/fluency/human-likeness.
9. **Treat "a human reviews it" as a design problem, not a safeguard** — guard
   against complacency under load and responsibility diffusion in teams.

## Cross-references

- Trust/reliance definitions, calibration vs. resolution, the 2×2 reliance
  metrics → `foundations-and-measurement.md`.
- Which way a given user/task will lean (aversion vs. appreciation) and the
  per-mode playbooks → `algorithm-aversion-vs-appreciation.md`.
- Underlying cognitive biases as general theory (anchoring, framing,
  availability) → the `behavioral-decision-making` skill (this skill only
  applies them to AI reliance).
