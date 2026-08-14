# Foundations & Measurement: Trust, Reliance, and Calibration

> Provenance: authored for the `human-ai-interaction-psychology` skill via /dr
> (deep-research workflow). Primary sources: Lee & See (2004); Hoff & Bashir
> (2015); the trust-measurement reviews (Frontiers, 2021); and the XAI
> trust-vs-reliance distinction (Scharowski et al., arXiv:2203.12318, 2022).

## Trust as attitude vs. reliance as behavior

The single most important conceptual move in this literature, and the one most
often botched in practice:

- **Trust** = an *attitude* — an internal, latent psychological state. Lee &
  See (2004): "the attitude that an agent will help achieve an individual's
  goals in a situation characterized by uncertainty and vulnerability." Measured
  **subjectively** (self-report scales: Jian et al.'s trust-in-automation scale,
  Madsen & Gregor's HCT, etc.).
- **Reliance** = a *behavior* that follows from the attitude — switching the
  automation on, accepting its suggestion, not double-checking. **Directly
  observable**, so it is measured **objectively** (did the user follow the
  advice? did they override? response times; monitoring/cross-check rate).

Why keep them apart (Scharowski et al., 2022): researchers frequently use a
*behavioral* measure (reliance) when they claim to be measuring the *attitude*
(trust), or vice versa. They dissociate — you can trust a system and still not
rely on it (e.g., policy forbids it; you lack authority), and you can rely on a
system you distrust (e.g., no alternative, time pressure). Conflating them
produces wrong conclusions about what an intervention did.

**Operator consequence:** when someone says "our users trust the AI," ask what
was actually measured. A survey of attitude tells you little about whether their
*reliance behavior* matches the tool's reliability — which is what determines
decision quality.

## Calibration, resolution, and the trust-calibration curve

- **Calibration** = correspondence between the *level* of trust and the *true
  capability* of the automation. Perfect calibration: trust tracks reliability
  one-for-one.
  - **Over-trust** (over-calibration): trust > capability → **over-reliance /
    misuse**.
  - **Under-trust** (under-calibration): trust < capability → **under-reliance /
    disuse**.
- **Resolution** = how well trust *discriminates between situations*. High
  resolution: the user trusts the system in exactly the cases it handles well
  and distrusts it in the cases it doesn't. A system can be well-calibrated *on
  average* yet have poor resolution (a flat, undifferentiated level of trust
  across cases it should be distinguishing). For real decisions, **resolution is
  what you want** — case-by-case appropriate reliance, not a global trust set-point.
- **Specificity** (Hoff & Bashir): trust can be *functionally specific* (per
  feature/sub-function) and *temporally specific* (right now vs. in general).
  Good calibration is specific, not blanket.

**The trust-calibration curve:** plot trust (y) against true automation
reliability (x). The 45° diagonal is perfect calibration. Points above the line
= over-trust; below = under-trust. The design goal is to move users *onto* the
diagonal and to sharpen resolution so they sit on the right point of it for each
case.

## What "appropriate reliance" means and how it is measured

Appropriate reliance is the behavioral target: **rely on the AI when it is
correct, and override it when it is wrong.** It decomposes into the cells of a
2×2 (AI correct/incorrect × human follows/overrides):

| | AI correct | AI incorrect |
| --- | --- | --- |
| **Human follows AI** | appropriate acceptance | **over-reliance** |
| **Human overrides AI** | **under-reliance** (self-bias / disuse) | appropriate rejection |

Common metrics in the HCI literature (Schemmer et al., 2023; Bansal et al.,
2021):

- **Over-reliance rate** — proportion of AI-incorrect cases where the human
  followed anyway.
- **Under-reliance rate** — proportion of AI-correct cases where the human
  overrode (often driven by *self-bias*: over-weighting one's own judgment).
- **Relative / appropriate reliance scores** — the degree to which switching
  behavior tracks AI correctness, ideally normalized against a baseline of
  always-follow and never-follow.
- **Team accuracy vs. AI-alone and human-alone** — the complementarity test
  (see the complementarity section in `reliance-interventions-and-design.md`).

Why it matters: an intervention can raise overall *agreement* with the AI while
making reliance *less* appropriate (more over-reliance). You cannot see that
without the 2×2 — agreement alone is a misleading headline metric.

## The closed-loop model of trust dynamics

Lee & See model trust as evolving in a closed loop, not a fixed trait:

1. **Dispositional trust** — the user's baseline propensity to trust machines
   (stable individual difference; varies with personality, age, culture).
2. **Situational / learned trust** — updates from *performance feedback*: the
   system's observed reliability, predictability, and the **visibility of its
   errors** (a single salient failure can swing trust far more than its base
   rate warrants — the hook for algorithm aversion).
3. **Context** — workload, stakes, time pressure, and organizational/cultural
   norms about deferring to machines all modulate how attitude becomes reliance
   behavior.

**Operator consequence:** trust is *built and lost continuously*. The first
visible error lands hardest (set expectations before it). Organizational norms
("we always sanity-check the model" vs. "the model is the source of truth")
shape reliance as much as the individual's attitude — so adoption is partly a
norms problem, not just a UX problem.

## Quick-reference definitions

- **Calibrated trust** — trust level matches true capability.
- **Over-trust → over-reliance / misuse** — deferring when you shouldn't.
- **Under-trust → under-reliance / disuse** — rejecting help that would have
  worked.
- **Resolution** — trust that discriminates good cases from bad ones.
- **Reliance** — the observable follow/override behavior; the thing to instrument.
- **Self-bias** — over-weighting one's own judgment over correct AI advice (a
  driver of under-reliance).
