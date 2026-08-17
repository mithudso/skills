<!-- FOLDED SPOKE of the applied-psychology hub. This is the full former standalone skill; any "use the <X>" / SKIP pointer below now refers to a sibling reference under applied-psychology/references/<X>/. -->

---
name: human-ai-interaction-psychology
description: >-
  Psychology of human-AI interaction — trust in automation and appropriate
  reliance applied to AI-assisted decisions. Covers calibrated trust vs.
  miscalibration (over-trust/over-reliance, under-trust/disuse), trust as an
  attitude vs. reliance as a behavior, automation bias and complacency
  (commission vs. omission errors), algorithm aversion (Dietvorst) and algorithm
  appreciation (Logg), why explanations/confidence displays often FAIL to
  calibrate reliance, cognitive forcing functions, human-AI complementarity, and
  anthropomorphism / uncanny-valley effects on trust. The operator angle: how
  TAMs and customers should adopt, trust, and decide WITH AI tools.
  TRIGGER: "should I trust this AI output", "over-relying on the AI / AI
  automation", "team rubber-stamps the AI", "why don't people trust our AI
  feature", "calibrate trust in an AI tool", "AI gave a confident wrong answer
  and we believed it", "automation bias", "automation complacency", "algorithm
  aversion", "algorithm appreciation", "do explanations make people trust AI
  more", "confidence scores for AI advice", "human-AI team / human-in-the-loop
  decision quality", "AI copilot adoption resistance", "uncanny valley", "AI
  persona / anthropomorphism and trust", "appropriate reliance", "cognitive
  forcing function", "when should a human override the model".
  SKIP: interpersonal/team trust, psychological safety, rapport, and influence
  between PEOPLE (use trust-and-psychological-safety); the underlying
  cognitive biases as general decision-making theory — anchoring, framing,
  availability (use a behavioral-decision-making skill; this skill only applies
  them to AI reliance); driving customer feature adoption via motivation/habit
  models with no trust-in-machine angle (use behavior-change-psychology); AI
  ethics, fairness, governance, and privacy (use da-11-ethics-and-privacy);
  building/evaluating the ML model or agent itself, RAG, guardrails, eval
  harnesses (use ai-agent-engineering or da-analytical-methods).
---

> **Folded spoke of the `applied-psychology` hub.** This file was the standalone `human-ai-interaction-psychology` skill; it now lives at `applied-psychology/references/human-ai-interaction-psychology/SKILL.md`.
> Any "use the <X> skill" / "(use ...)" pointer below to another psychology skill now refers to a **sibling reference** under `applied-psychology/references/<X>/` — not a top-level skill.

# Psychology of Human-AI Interaction: Trust & Appropriate Reliance

How humans decide whether to follow, override, or ignore an AI system, and how
to design and coach for the *right* amount of reliance. This is human-factors
and decision psychology applied to **trust in machines**, not interpersonal
trust. The central problem is not "more trust" or "less trust" but **calibrated
trust**: reliance that tracks the system's *actual* reliability in the specific
task at hand.

## When to use this skill

Reach for it whenever a human and an AI are jointly producing a decision and the
question is *how much the human should defer*:

- A TAM, customer, or team is **over-relying** (rubber-stamping AI output) or
  **under-relying** (ignoring a tool that outperforms them).
- An AI feature is hitting **adoption resistance** rooted in distrust, or
  hitting **dangerous over-adoption** where users stop checking.
- You are **designing an AI-assisted workflow** (copilot, recommender, triage
  assistant, autoremediation gate) and must decide what to surface (confidence,
  explanations, friction) to get appropriate reliance.
- A **confidently wrong** AI answer was believed, and you need the vocabulary
  (automation bias, plausible-but-wrong explanation) to diagnose why.
- You are coaching a customer on a **human-in-the-loop** policy: when must a
  person override the model.

## The one thing to get right

**Trust is an attitude; reliance is a behavior; appropriate reliance is the
goal.** They are routinely conflated and must be kept separate. You can *trust*
a system and still not *rely* on it (and vice-versa). Optimizing for "trust"
("users say they trust it") is the wrong target — optimize for **reliance that
matches reliability**: follow the AI when it is right, override it when it is
wrong. Most failures in AI-assisted decisions are *miscalibration*, not a global
trust deficit.

## Core concepts

### 1. Calibrated trust and the trust–reliance distinction (Lee & See, 2004)

Lee & See's foundational framing defines **trust** as "the attitude that an
agent will help achieve an individual's goals in a situation characterized by
uncertainty and vulnerability," and **reliance** as the *behavior* that follows
from it. **Calibration** is the correspondence between a person's trust and the
system's actual capability.

Two failure modes:

- **Over-trust → over-reliance / misuse.** Trust exceeds capability; the user
  defers to the automation when they shouldn't (accepts wrong advice, stops
  monitoring).
- **Under-trust → under-reliance / disuse.** Trust falls short of capability;
  the user rejects or abandons automation that would have helped.

The **trust-calibration curve** plots trust against true reliability: the
diagonal is perfect calibration; the region above it is over-trust, below it is
under-trust. **Resolution** is a related property — fine-grained trust that
discriminates *which* situations the system handles well from those it doesn't
(a system can be well-calibrated on average but have poor resolution, trusting
uniformly across cases it should distinguish). Calibration is a *closed loop*:
it evolves from performance feedback, the user's disposition to trust,
organizational norms, and culture.

> **Operator translation:** Don't ask "do you trust the tool?" Ask "for *which*
> decisions does it earn the follow, and which require your own check?" Good
> adoption coaching builds **resolution**, not blanket trust.

### 2. Automation bias & automation complacency (Parasuraman & Manzey, 2010)

When a decision aid is *imperfect*, two error types appear:

- **Errors of commission:** the user *follows* an incorrect automated directive
  without cross-checking other available information. (Acting on a wrong alert.)
- **Errors of omission:** the automation *fails to flag* a problem and the user
  misses it too, because they were not actively monitoring. (Not catching what
  the system didn't catch.)

**Automation complacency** is the attentional root: under multitasking load,
manual tasks compete with monitoring the automation, so the human reduces
vigilance and "lets the system handle it." Parasuraman & Manzey's *attentional
integration* argues bias and complacency share this mechanism.

Two findings that matter operationally:

- Automation bias appears in **experts as well as novices**, and is **not
  reliably eliminated by training or instructions** alone.
- It affects **teams as well as individuals**. Having "another set of eyes"
  does not guarantee someone catches the automation's error; social diffusion of
  monitoring responsibility can make it *worse*.

> **Operator translation:** "A human reviews it" is not, by itself, a control.
> If the human is loaded, complacent, or assumes the model caught everything,
> review degrades to rubber-stamping. Design the review step to *demand*
> engagement (see cognitive forcing functions, below).

### 3. Algorithm aversion (Dietvorst et al., 2015)

People **abandon algorithms faster than humans after seeing them err**, even
when the algorithm still outperforms the human. In Dietvorst, Simmons & Massey's
forecasting experiments, participants who *saw* the model make mistakes were
much less likely to choose it afterward, despite the model beating their own
judgment, because algorithmic errors trigger a sharper loss of confidence than
equivalent human errors. **Error visibility** is the key trigger: it is the
*seeing it fail*, not the failure rate, that drives aversion. A follow-up
("Overcoming Algorithm Aversion," 2016/2018) showed that letting people **adjust
the algorithm's output even slightly** substantially increases willingness to
use it: *control over the algorithm* restores reliance.

> **Operator translation:** A single visible AI miss can sink adoption of a tool
> that is net-better than the status quo. Counter it with (a) **adjustability**
> ("you can edit the suggestion"), (b) honest **expectation-setting** before the
> first error, and (c) framing errors as expected-and-bounded, not as proof the
> tool is broken.

### 4. Algorithm appreciation (Logg et al., 2019) — and reconciling the two

The opposite effect is just as real. In Logg, Minson & Moore's six experiments,
lay people **weighted identical advice MORE heavily when told it came from an
algorithm** than from a person (measured by Weight-On-Advice), across numeric
estimates and forecasts. Appreciation **waned** under two boundary conditions:
when people chose between the algorithm and **their own** judgment (vs. an
external human advisor), and when they had **domain expertise** in the task.

**Reconciling aversion vs. appreciation:** the literature is fragmented and
*context determines the sign*:

| Pulls toward APPRECIATION (over-weighting) | Pulls toward AVERSION (discounting) |
| --- | --- |
| No error observed yet; pre-feedback | A visible algorithmic error |
| Objective / numeric / quantifiable task | Subjective, moral, or "human" task |
| User is a non-expert | User is a domain expert |
| Advice attributed to "the algorithm/data" | Choice is algorithm vs. *the user's own* judgment |
| No control over the output | (Aversion eased when user can adjust output) |

> **Operator translation:** Whether your customer over- or under-relies is
> *predictable* from these levers. Experts in subjective domains who just watched
> a model whiff → expect aversion. Non-experts on a numeric task pre-error →
> expect over-reliance. Design for the failure mode you actually face.

### 5. Why explanations & confidence displays often FAIL to calibrate reliance

The intuitive fix ("show the reasoning / show a confidence score so people know
when to trust it") frequently **backfires**:

- **Plausible-but-wrong explanations increase over-reliance.** Bansal et al.
  (2021), *"Does the Whole Exceed Its Parts?"* found that AI explanations
  increased the chance people accepted the AI's recommendation **whether it was
  right or wrong** — explanations made advice *feel* credible without improving
  discrimination. Explanations raised agreement, not accuracy.
- **Confidence scores help only if the confidence is itself well-calibrated.**
  Miscalibrated confidence (a wrong answer shown as "95% sure") actively
  *degrades* trust calibration and decision quality. Recent work also finds that
  displaying AI confidence **shifts the user's own self-confidence** toward the
  model's, anchoring the human's uncertainty without improving their actual
  ability, which can *undermine* appropriate reliance.
- **The mechanism (dual-process):** People rarely engage System-2 analysis with
  each recommendation; they fall back on a heuristic ("the AI is usually right,
  and here's a reason, so accept"). Explanations feed the heuristic instead of
  interrupting it.

> **Operator translation:** "We added explanations / confidence scores" is **not
> evidence** the tool now produces appropriate reliance — and can be evidence of
> the opposite. Calibration must be *measured behaviorally* (does override-rate
> track error-rate?), not assumed from the presence of an explanation. Only ship
> confidence numbers you have *validated as calibrated*.

### 6. Cognitive forcing functions & appropriate-reliance interventions

Because passive transparency fails, the effective lever is **friction that
compels analytical engagement at decision time** — *cognitive forcing
functions* (Buçinca, Malaya & Gajos, 2021, *"To Trust or to Think"*):

- **Make the human commit first**, then reveal the AI suggestion (forces an
  independent judgment before anchoring).
- **Withhold the recommendation on demand / add a wait** so the user must
  request it deliberately rather than receive it passively.
- **Show the AI's reasoning *on request*, not by default**, and surface
  **uncertainty/disagreement** as a prompt to slow down.

These **reduced over-reliance on incorrect AI** more than explanation-only
designs. The catch: forcing functions impose effort, are often *disliked*, and
their benefit **interacts with the user** (e.g., people higher in Need for
Cognition engage; others may disengage or route around them). Other
appropriate-reliance levers in the recent literature: **onboarding/tutorials**
that teach the model's error boundaries, **selective explanations** tied to
likely-error cases, and **adjustable outputs** (the Dietvorst lever).

> **Operator translation:** To stop rubber-stamping, change the *workflow*, not
> the disclaimer. "Write your own assessment first, then open the AI draft" beats
> "please review carefully." Reserve heavy friction for high-stakes / likely-wrong
> cases — it has a cost.

### 7. Human-AI complementarity (the actual goal of the team)

**Complementary Team Performance (CTP)**: the human+AI together beat *both* the
best human alone and the AI alone — achieved only when their errors are
*different* and each defers on the cases the other handles better. Sobering
finding from Bansal et al. and the broader literature: **CTP is hard and rarely
achieved by default.** Adding explanations or confidence does *not* reliably
produce it; teams frequently perform *worse* than the AI alone because humans
over-rely on wrong answers and override right ones. Complementarity requires the
human to hold knowledge the model lacks (context, side information,
"unobservables") *and* to apply it through appropriate reliance.

> **Operator translation:** "Human + AI" is not automatically better than either
> — that has to be *engineered and verified*. Put the human where they have an
> information edge the model lacks (account history, political context, a
> just-happened event the model can't see), not as a generic reviewer of things
> the model already does better.

### 8. Anthropomorphism, persona & the uncanny valley

How human-like an AI *seems* shifts trust independently of how *capable* it is:

- **Anthropomorphic cues** (name, persona, conversational warmth, an avatar) can
  raise initial trust and likability — but the effect is **mediated by perceived
  empathy and interaction quality**, not the cosmetic cues alone, and explicit
  "I am an AI" identity cues can override surface human-likeness.
- **Uncanny valley** (Mori, 1970): as an agent approaches — but doesn't reach —
  human-likeness, affinity *drops sharply* into eeriness. An **"uncanny valley of
  trust"** arises when the bot is human enough to raise competence expectations
  it can't meet, producing skepticism and over-promise/under-deliver
  disappointment.
- **Persona inflates trust calibration risk:** a warm, fluent, confident persona
  is a powerful (and misleading) cue that pushes users toward over-reliance
  regardless of correctness — fluency reads as competence. This is the persona
  analog of the plausible-but-wrong-explanation problem.

> **Operator translation:** A charming, confident AI persona can *manufacture*
> over-trust the capability doesn't justify. For high-stakes tools, prefer a
> sober, capability-honest voice over a maximally human, maximally confident one;
> match the persona's confidence to the system's validated reliability.

## Design & coaching principles (synthesis)

A consolidated checklist for designing — or advising a customer on — an
AI-assisted decision workflow:

1. **Target appropriate reliance, measured behaviorally.** Success = override
   tracks error. Instrument it (agreement-on-correct vs. agreement-on-wrong;
   over-reliance and under-reliance rates). "Users report high trust" is not the
   metric.
2. **Set honest expectations before the first error.** State the error
   boundaries and likely-wrong cases up front to inoculate against algorithm
   aversion (and against blind over-trust).
3. **Show confidence only if it's calibrated; communicate uncertainty
   honestly.** Validate confidence before displaying it; epistemic markers ("I'm
   not certain…") help over time. Miscalibrated confidence is worse than none.
4. **Don't rely on explanations to create skepticism — they usually do the
   opposite.** If you ship explanations, pair them with friction and verify
   behavior; consider *selective* explanations on likely-error cases.
5. **Engineer friction where stakes are high** (commit-first, on-demand reveal,
   forced confrontation of disagreement). Reserve it — friction has a cost and
   is unevenly effective across users.
6. **Give users control / adjustability** over outputs to restore reliance after
   visible errors (the Dietvorst lever).
7. **Place the human where they have an information edge,** not as a generic
   reviewer — that's where complementarity actually lives.
8. **Match persona confidence to validated capability;** avoid manufacturing
   over-trust through warmth, fluency, or human-likeness.
9. **Treat "a human reviews it" as a design problem, not a safeguard** — guard
   against complacency under load and diffusion of responsibility in teams.

## Anti-patterns

- **Optimizing for "trust" as a survey number** instead of for calibrated
  reliance behavior.
- **Shipping explanations or confidence scores and declaring the over-reliance
  problem solved** — they often *increase* over-reliance; confidence helps only
  if calibrated.
- **Treating a human-in-the-loop step as a guaranteed control** while the
  reviewer is loaded, complacent, or assuming the model caught everything.
- **Letting one visible AI error kill adoption** of a net-better tool, with no
  adjustability or expectation-setting to absorb it (avoidable algorithm
  aversion).
- **Maxing out an anthropomorphic, confident persona** on a high-stakes tool and
  manufacturing over-trust the capability can't back.
- **Assuming "human + AI" beats either alone** — complementarity is rare and
  must be verified, not presumed.
- **Conflating trust and reliance** in analysis or instrumentation.

## Operator scenarios (TAM / AI-native workflow)

- **"The team just rubber-stamps the AI triage."** Diagnose: automation
  bias/complacency + over-reliance, likely amplified by confident output. Fix:
  commit-first workflow, surface disagreement/uncertainty, instrument
  agreement-on-wrong, and reserve the human for cases where they have context the
  model lacks.
- **"Our customer's analysts refuse to use the new recommender."** Diagnose:
  likely algorithm aversion — experts, possibly post a visible error, choosing
  model-vs-own-judgment. Fix: adjustability, expectation-setting on error
  bounds, frame as advisor not replacement, show the win-rate vs. their baseline.
- **"We added explanations and people trust it more — ship it?"** Caution: more
  *agreement* is not more *appropriate reliance*. Verify override tracks error;
  watch for plausible-but-wrong over-acceptance.
- **"Should we give the assistant a friendly human persona?"** It will lift
  initial likability but can manufacture over-trust and risk the uncanny valley;
  for high-stakes decisions keep the voice capability-honest and match its
  confidence to validated reliability.
- **"What's our human-in-the-loop / override policy?"** Define it by *resolution*
  — the specific case classes where humans must independently judge — not a blanket
  "review everything," which decays into rubber-stamping.

## References (depth-first)

`references/` holds the deeper material — load it when a question goes past this
overview:

- **`references/foundations-and-measurement.md`** — Lee & See trust model, the
  trust–reliance attitude/behavior distinction, calibration vs. resolution, the
  trust-calibration curve, and how appropriate reliance / over- & under-reliance
  are measured.
- **`references/algorithm-aversion-vs-appreciation.md`** — Dietvorst and Logg in
  detail, the reconciliation map of moderators (error visibility, expertise,
  task type, control), and operator playbooks for each failure mode.
- **`references/reliance-interventions-and-design.md`** — automation
  bias/complacency mechanisms, the explanation/confidence backfire evidence,
  cognitive forcing functions, complementarity (CTP), anthropomorphism/uncanny
  valley, and the consolidated design checklist with the recent (2023-2025) HCI
  evidence base.

### Key sources

1. Lee, J. D., & See, K. A. (2004). *Trust in Automation: Designing for
   Appropriate Reliance.* Human Factors, 46(1), 50-80.
2. Parasuraman, R., & Manzey, D. H. (2010). *Complacency and Bias in Human Use
   of Automation: An Attentional Integration.* Human Factors, 52(3), 381-410.
3. Dietvorst, B. J., Simmons, J. P., & Massey, C. (2015). *Algorithm Aversion:
   People Erroneously Avoid Algorithms After Seeing Them Err.* Journal of
   Experimental Psychology: General, 144(1), 114-126. (See also Dietvorst et
   al., 2018, *Overcoming Algorithm Aversion*, Management Science.)
4. Logg, J. M., Minson, J. A., & Moore, D. A. (2019). *Algorithm Appreciation:
   People Prefer Algorithmic to Human Judgment.* Organizational Behavior and
   Human Decision Processes, 151, 90-103.
5. Bansal, G., et al. (2021). *Does the Whole Exceed Its Parts? The Effect of AI
   Explanations on Complementary Team Performance.* CHI 2021.
6. Buçinca, Z., Malaya, M. B., & Gajos, K. Z. (2021). *To Trust or to Think:
   Cognitive Forcing Functions Can Reduce Overreliance on AI in AI-Assisted
   Decision-Making.* Proc. ACM HCI (CSCW1).
7. Mori, M. (1970/2012). *The Uncanny Valley.* IEEE Robotics & Automation
   Magazine.
8. Microsoft Research (2024). *Appropriate Reliance on Generative AI: Research
   Synthesis.* And recent CHI 2024-2025 work on miscalibrated AI confidence and
   confidence–self-confidence alignment.
