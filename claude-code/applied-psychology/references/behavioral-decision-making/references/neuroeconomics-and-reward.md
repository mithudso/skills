# Neuroeconomics & Reward — Mechanism Behind the Behavior

The **neural-mechanism counterpoint** to the behavioral effects in the main skill. Where SKILL.md
describes *what* people do (loss aversion, discounting, framing), this file sketches the brain systems
that *implement* valuation and learning — dopamine prediction errors, value coding in vmPFC and
striatum, the common-currency idea, and emotion's role in choice (somatic markers).

> **Honesty rule for this file (read first).** Keep claims at the **mechanism level** and resist
> over-interpretation. fMRI is correlational and coarse; "region X lit up" does **not** license
> **reverse inference** ("therefore the subject was computing value / feeling loss") — most regions
> serve many functions (Poldrack, 2006). A behavioral effect being real does **not** require a clean
> one-region "center," and finding a neural correlate does **not** add normative weight to a decision
> ("my striatum made me do it" explains nothing the behavior didn't already show). Neuroscience here
> *enriches the model*; it does not *prove* the behavioral story or hand you a lever. The operator
> payoff is conceptual fluency, not a brain-based sales tactic.

---

## 1. Dopamine reward prediction error (Schultz)

The best-validated link between a neurotransmitter and a computational quantity in decision
neuroscience. Wolfram Schultz recorded midbrain dopamine neurons and found they don't signal reward
itself — they signal the **reward prediction error (RPE)**: *received minus predicted* reward.

- **Positive RPE** (better than predicted): phasic burst, dopamine fires above baseline.
- **Fully predicted** reward: **no** response — the neuron stays at baseline (the prediction already
  accounted for it).
- **Negative RPE** (worse than predicted, or an expected reward omitted): activity **dips below**
  baseline at the time the reward was due.
- **Transfer to the predictor:** with learning, the burst migrates from the reward to the earliest
  *cue* that predicts it — exactly the signature of a learned value estimate.

**Why it mattered: it maps onto reinforcement learning.** This RPE signal is a near-exact biological
match to the **teaching signal (δ) in temporal-difference (TD) learning** — itself a descendant of the
Rescorla–Wagner model. The dopamine signal is plausibly a **teaching signal** that drives synaptic
plasticity in the striatum, frontal cortex, and amygdala, updating value estimates so future
predictions improve (Schultz, Dayan & Montague, 1997). This is the headline result of the field: a
unit-recording finding and a machine-learning equation describing the same computation.

- *Caveats / honest scope:* dopamine does more than RPE — it also tracks **reward uncertainty**,
  **salience/novelty**, and movement/motivation (Berke, 2018; "Dopamine, Updated," Watabe-Uchida &
  Uchida, 2020). The RPE account is robust but **not the whole story**; don't reduce dopamine to a
  "pleasure" or "reward" molecule (a common pop-science error — it's a *learning/teaching* signal, not
  a hedonic one).

---

## 2. Neural valuation: subjective value in vmPFC and striatum; the common-currency hypothesis

To choose between unlike things (a discount vs a feature, money vs time, food vs status), the brain
appears to put them on **one comparable scale** — a **"common currency"** of subjective value.

- **Where:** the **ventromedial prefrontal cortex (vmPFC) / medial orbitofrontal cortex** and the
  **ventral striatum** carry signals that track **subjective value** — domain-generally. BOLD activity
  in these regions correlates parametrically with how much *this subject* values a reward, whether the
  reward is money, food, or a consumer good, and whether the outcome is delayed or probabilistic
  (Kable & Glimcher, 2007; Bartra, McGuire & Kable, 2013, coordinate-based meta-analysis; Levy &
  Glimcher, 2012, "The root of all value").
- **Common-currency hypothesis:** because the *same* circuitry scales with value across reward types,
  these regions are read as implementing a shared valuation metric that makes otherwise
  incommensurable options comparable — the neural precondition for choice.

**Critical caveat — "an embarrassment of riches."** O'Doherty (2014, *PLOS Biology*, "Representations
of Value in the Brain: An Embarrassment of Riches?") and Bartra et al. (2013) note that value-correlated
signals appear in *many* regions, not a single "value center," and that a correlation with value can
also reflect attention, arousal, salience, motor preparation, or outcome anticipation. So:
- "Common currency" is a **leading hypothesis**, not a settled fact; whether vmPFC/striatum compute a
  *single integrated* value or several overlapping value-like signals is **actively debated**.
- Do not over-localize. The honest statement is "vmPFC and ventral striatum *reliably track* subjective
  value across domains," **not** "the vmPFC *is* the brain's value calculator."

---

## 3. Neural basis of loss aversion, risk, and temporal/delay discounting

These connect the behavioral effects in SKILL.md §3 and §5 to the valuation system above — as
*correlates and mechanisms*, with the same correlational caveats.

- **Loss aversion (Tom, Fox, Trepel & Poldrack, 2007, *Science*).** For 50/50 gain/loss gambles, a set
  of value-sensitive regions (ventral striatum, vmPFC, among others) *increased* activity as potential
  **gains** grew and *decreased* more steeply as potential **losses** grew — "**neural loss aversion**":
  losses suppressed the value signal faster than equivalent gains raised it, and the steepness of that
  asymmetry **predicted individuals' behavioral** loss aversion. Notably this study found **no special
  role for the amygdala / fear circuitry** at the choice stage, complicating a purely emotional "fear of
  loss" account. (Other work implicates amygdala/striatum more at *outcome receipt* than at choice, and
  loss aversion's status as a universal constant is itself contested — see SKILL.md §3 and
  `references/replication-status.md`.)
- **Risk / probability.** Expected value, risk (variance), and probability-weighting correlates appear
  across striatal and prefrontal/parietal regions; risk and expected value are partly dissociable
  signals. The behavioral overweighting of small probabilities (prospect theory) has neural correlates
  but no single tidy "probability-distortion center."
- **Temporal / delay discounting.** vmPFC and ventral striatum track the **discounted (present)
  subjective value** of delayed rewards — i.e., value *after* discounting, consistent with present bias
  / hyperbolic discounting (SKILL.md §5). McClure et al. (2004) proposed a dual-system split (limbic
  "β" valuing immediacy vs prefrontal "δ" for the long view); Kable & Glimcher (2007) argued instead for
  a **single** valuation system whose output is simply discounted by delay. **This dual-vs-single debate
  is unresolved** — present it as two competing models, not a fact.

---

## 4. Somatic marker hypothesis (Damasio) — emotion's role in decision quality

Antonio Damasio's claim that **emotion is not the enemy of good decisions but a necessary input**.
Bodily/affective states ("somatic markers") become associated with the anticipated outcomes of options
and bias choice rapidly, before (or alongside) explicit reasoning — a fast affective shortcut through an
otherwise unmanageable option space (Damasio, *Descartes' Error*, 1994).

- **Substrate:** somatic markers are integrated in the **vmPFC**, drawing on the amygdala and
  body-state signaling. **vmPFC-lesion patients** retain IQ and logic yet make disastrous
  personal/financial decisions — the evidence that "cold" reasoning without affective signal is
  *impaired*, not purer.
- **Iowa Gambling Task (Bechara, Damasio et al., 1994/1997):** subjects learn to avoid high-payout but
  net-losing decks. Healthy subjects develop **anticipatory skin-conductance responses** before
  reaching toward a bad deck — an emotional "warning" — *before* they can verbally explain which decks
  are bad. vmPFC-lesion patients fail to develop these markers and keep choosing badly.
- *Caveats / honest scope:* the SMH has drawn serious criticism (Dunn, Dalgleish & Lawrence, 2006, a
  critical evaluation). The "**before conscious knowledge**" claim and the SCR-as-cause interpretation
  are contested; the IGT confounds learning, working memory, and reversal, so poor IGT performance is
  **not** clean proof of an emotion-specific deficit. Treat SMH as an **influential, partially supported
  framework** — the safe, defensible takeaway is the broad one: *emotion is integral to normal
  decision-making, and decisions stripped of affect are impaired*, not the strong specific mechanism.

---

## 5. Operator note — use the mechanism level, stay honest

- **What this buys you:** a richer mental model of why the behavioral effects in SKILL.md are *robust*
  — they're grounded in learning and valuation systems, not arbitrary quirks. Prediction-error framing
  is genuinely useful for thinking about expectations management: a delivered win that was fully
  *predicted* generates **no** "reward signal" (zero RPE) — so a result you over-promised lands flat
  even when objectively good, while the same result against a modest expectation registers as a
  positive surprise. That's an expectations lesson, stated honestly as analogy.
- **What it does NOT buy you:**
  - **No reverse inference, no brain-scan sales tactics.** Never claim a neural finding *proves* a
    customer's motive or that you can "target their striatum." That's pseudoscience and it's
    reputationally radioactive. (See the over-interpretation critiques throughout this file.)
  - **No normative upgrade.** A neural correlate doesn't make a biased decision rational or an
    intuition trustworthy — that question is settled by the *environment*, per Kahneman–Klein (see
    `references/naturalistic-decision-making.md`), not by neuroimaging.
  - **Don't reduce people to chemicals.** "Dopamine hit," "lizard brain," "amygdala hijack" are
    pop-neuro shorthand that overstate localization and mislead. Keep the behavioral level as your
    primary working layer; reach for mechanism only to *explain robustness*, and flag the caveats.
- **Bottom line:** cite neuroeconomics to deepen understanding and to resist the "people are just
  irrational" framing — never as evidence for a manipulation, and never past what correlational fMRI
  and lesion data actually support.

---

## Sources

- Schultz, W., Dayan, P., & Montague, P. R. (1997). "A Neural Substrate of Prediction and Reward."
  *Science* 275(5306), 1593–1599. (dopamine RPE ↔ temporal-difference learning)
- Schultz, W. (2016). "Dopamine reward prediction error coding." *Dialogues in Clinical Neuroscience*
  18(1), 23–32. (review/update)
- Watabe-Uchida, M., & Uchida, N. (2020). "Dopamine, Updated: Reward Prediction Error and Beyond."
  *Current Opinion in Neurobiology.* (beyond-RPE: uncertainty, salience — honest scope)
- Kable, J. W., & Glimcher, P. W. (2007). "The neural correlates of subjective value during
  intertemporal choice." *Nature Neuroscience* 10(12), 1625–1633. (single-system discounted value in
  vmPFC/striatum)
- Levy, D. J., & Glimcher, P. W. (2012). "The root of all value: a neural common currency for choice."
  *Current Opinion in Neurobiology* 22(6), 1027–1038.
- Bartra, O., McGuire, J. T., & Kable, J. W. (2013). "The valuation system: a coordinate-based
  meta-analysis of BOLD fMRI experiments…" *NeuroImage* 76, 412–427. (where value tracks reliably)
- O'Doherty, J. P. (2014). "Representations of Value in the Brain: An Embarrassment of Riches?"
  *PLOS Biology* 12(6), e1002174. (over-localization / common-currency caveat)
- Poldrack, R. A. (2006). "Can cognitive processes be inferred from neuroimaging data?" *Trends in
  Cognitive Sciences* 10(2), 59–63. (the reverse-inference problem — the key methodological caveat)
- Tom, S. M., Fox, C. R., Trepel, C., & Poldrack, R. A. (2007). "The Neural Basis of Loss Aversion in
  Decision-Making Under Risk." *Science* 315(5811), 515–518.
- McClure, S. M., Laibson, D. I., Loewenstein, G., & Cohen, J. D. (2004). "Separate neural systems
  value immediate and delayed monetary rewards." *Science* 306(5695), 503–507. (dual-system view —
  contested)
- Damasio, A. R. (1994). *Descartes' Error: Emotion, Reason, and the Human Brain.* (somatic marker
  hypothesis)
- Bechara, A., Damasio, A. R., Damasio, H., & Anderson, S. W. (1994); Bechara et al. (1997, *Science*).
  Iowa Gambling Task and anticipatory SCRs.
- Dunn, B. D., Dalgleish, T., & Lawrence, A. D. (2006). "The somatic marker hypothesis: A critical
  evaluation." *Neuroscience & Biobehavioral Reviews* 30(2), 239–271. (the principal critique)
