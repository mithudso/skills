<!-- FOLDED SPOKE of the applied-psychology hub. This is the full former standalone skill; any "use the <X>" / SKIP pointer below now refers to a sibling reference under applied-psychology/references/<X>/. -->

---
name: behavioral-decision-making
description: >-
  The descriptive science of how people actually judge and decide (cognitive
  biases, prospect theory, choice architecture), applied to customer/buyer
  decisions, negotiation framing, and the operator's own judgment. TRIGGER:
  "why did the customer decide X", "they're being irrational", "anchor the
  negotiation / pricing", "frame this as a gain or loss", "loss aversion",
  "sunk cost", "cognitive bias", "behavioral economics", "prospect theory",
  "System 1 / System 2", "nudge the buyer toward", "default option", "choice
  architecture", "too many options / choice overload", "de-bias my decision",
  "premortem", "reference-class forecast", "am I overconfident in this
  forecast", "is power posing / ego depletion / priming real", "status quo
  bias", "anchoring effect", "framing effect", "naturalistic decision making",
  "recognition-primed decision / RPD", "expert intuition", "when to trust your
  gut", "premortem origin", "dopamine prediction error", "reward prediction
  error", "neuroeconomics", "neural valuation / common currency", "somatic
  marker". Covers the heuristics-and-biases
  program (availability, representativeness, anchoring, confirmation, hindsight,
  overconfidence, status-quo/default, sunk-cost); prospect theory (loss
  aversion, reference dependence, probability weighting, fourfold pattern,
  gain/loss framing); dual-process System 1/2; bounded rationality and
  satisficing (Simon) vs ecological rationality and fast-and-frugal heuristics
  (Gigerenzer); mental accounting; present bias / hyperbolic discounting; nudges
  (defaults, EAST, MINDSPACE, sludge) vs boosts (Hertwig); debiasing; and the
  replication status of major decision/social-psych effects so you never cite a
  debunked one. SKIP: designing a habit/adoption/enablement program to change
  customer behavior over time, motivation, Fogg B=MAP, stages-of-change, habit
  loops (use behavior-change-psychology); normative/prescriptive decision
  optimization, LP/MILP, decision trees, EVPI, utility maximization, "the
  optimal action" (use da-33-prescriptive-analytics); persuasion/influence
  theory and the mechanism layer of attitude change — ELM/HSM, cognitive
  dissonance, reactance, attribution, social-influence/norm theory (use
  persuasion-and-influence-psychology); applied persuasion/negotiation tactics
  — Cialdini's 6, Voss, BATNA (use negotiation-and-persuasion); A/B-test
  statistics and causal lift measurement (use
  da-12-ab-testing-causal-inference / da-21-product-analytics); trust in
  automation, over-/under-reliance on AI, automation bias, algorithm
  aversion/appreciation, and calibrated reliance on AI tools (use
  human-ai-interaction-psychology); building expertise / skill-acquisition
  stages / deliberate practice as a learning topic (use
  learning-and-expertise-psychology); the neuroscience of emotion generation
  per se rather than its role as a decision input (use
  emotion-and-affect-psychology).
---

> **Folded spoke of the `applied-psychology` hub.** This file was the standalone `behavioral-decision-making` skill; it now lives at `applied-psychology/references/behavioral-decision-making/SKILL.md`.
> Any "use the <X> skill" / "(use ...)" pointer below to another psychology skill now refers to a **sibling reference** under `applied-psychology/references/<X>/` — not a top-level skill.

# Behavioral Decision-Making & Cognitive Biases

The **descriptive** account of judgment and decision-making: how people *actually*
decide, not how they *should*. Use it to read why a customer, buyer, or stakeholder
made a seemingly irrational choice, to ethically shape the decision environment, and
to catch bias in your own forecasts and recommendations.

**Descriptive vs normative: keep them separate.** This skill is descriptive (the
psychology of real decisions). For the *normative* side, computing the optimal action
under constraints (linear programming, decision trees, EVPI, expected-utility
maximization), use `da-33-prescriptive-analytics`. The gap between the two *is* the
subject matter here: people deviate from the normative optimum in systematic, predictable ways.

**The honesty rule (read first).** Decision/social psychology went through a replication
crisis. Several once-famous effects did **not** survive (power posing, social priming,
and ego depletion are contested/failed). The core judgment-and-decision-making findings
(anchoring, framing, the disposition effect, present bias, default effects) replicate well;
several adjacent social-psych effects do not. **Never present a debunked effect as
established fact.** See `references/replication-status.md` before citing any effect to a
customer or in a written recommendation.

---

## When to reach for this skill

- A customer/buyer made a choice that looks irrational → name the bias, then address the real driver.
- You're setting a price, an opening offer, or a contract renewal → anchoring and framing.
- A renewal/expansion stalls on "we already invested in X" → sunk-cost; on "let's keep things as they are" → status-quo/default bias.
- You're designing a signup, plan-selection, or opt-in/opt-out flow → choice architecture (and the sludge you should *remove*).
- You're writing a forecast, capacity plan, or project timeline → debias your *own* judgment (overconfidence, planning fallacy).
- Someone cites "power posing" / "priming" / "ego depletion" as fact → check replication status first.

If the task is changing a customer's *behavior over time* (adoption, habit, enablement), that's `behavior-change-psychology`, not this skill.

---

## Core concepts (the minimum working set)

### 1. Dual-process theory: System 1 / System 2
Two modes of cognition. **System 1** is fast, automatic, associative, affect-laden, and effortless; it produces most snap judgments and most biases. **System 2** is slow, deliberate, effortful, and lazy (it endorses System 1 unless prompted). Labels coined by **Stanovich & West**, popularized by **Kahneman** (*Thinking, Fast and Slow*, 2011). Biases are System 1 outputs that System 2 fails to catch.
- *Caveat:* treat "two systems" as a **useful metaphor, not literal brain architecture**. The strict two-box model is contested (it's better read as a continuum of automaticity). Don't oversell it.
- *Operator use:* high-stakes decisions (renewals, escalations, architecture calls) deserve a deliberate System-2 step — a checklist or premortem — precisely because the default is a System-1 gut call.

### 2. Heuristics & biases (Tversky & Kahneman, 1974, *Science*)
Mental shortcuts that are "highly economical and usually effective" but produce "systematic and predictable errors." The working set every operator should recognize:
- **Anchoring-and-adjustment** — the first number dominates; adjustments are insufficient. *The single most useful effect in negotiation.*
- **Availability** — judging probability by how easily examples come to mind (the loud outage feels more likely than the silent risk).
- **Representativeness** — judging by stereotype/similarity, ignoring base rates (the conjunction fallacy, "Linda the bank teller").
- **Confirmation bias** — seeking/weighting evidence that fits the prior; the engine behind most bad root-cause calls.
- **Hindsight bias** — "I knew it all along" after the fact; corrupts postmortems (see `postmortem-writing`).
- **Overconfidence** — calibration gap; 90%-confident estimates are right far less than 90% of the time. Poisons forecasts and timelines (**planning fallacy**).
- **Status-quo / default bias** — disproportionate preference for the current state / the pre-selected option.
- **Sunk-cost fallacy** — honoring past, unrecoverable spend ("we've already invested 18 months in this cluster design") instead of deciding on the margin.

Detail and operator scripts: `references/biases-catalog.md`.

### 3. Prospect theory (Kahneman & Tversky, 1979, *Econometrica*)
How people choose under risk — the descriptive replacement for expected-utility theory:
- **Reference dependence** — outcomes are judged as *changes* (gains/losses) from a reference point, not absolute states. **Whoever sets the reference point frames the decision.**
- **Loss aversion** — losses loom larger than equivalent gains. Classic estimate ≈ 2×. **CONTESTED — see below; do not state "2×" as universal law.**
- **Diminishing sensitivity** — the value function is concave for gains, convex for losses (the difference between $10 and $20 feels bigger than $1010 vs $1020).
- **Probability weighting** — small probabilities are *overweighted* (why people buy lottery tickets *and* insurance), moderate-to-high probabilities underweighted; the certainty effect makes "guaranteed" disproportionately attractive.
- **The fourfold pattern** — risk-averse for likely gains and unlikely losses; risk-seeking for unlikely gains and likely losses.
- **Framing effects** — logically identical options described as gains vs losses flip the choice (the Asian-disease problem; "90% uptime" vs "10% downtime").

> **Loss aversion is contested, not debunked.** Gal & Rucker (2018), "The Loss of Loss Aversion," argue it is far more context-dependent than the "universal 2× law" framing implies; some studies find *gains* loom larger at small magnitudes (Harinck et al., 2007), and predicted emotional pain of loss overstates the actual (Kermer et al., 2006). Loss aversion is real *in many settings* (especially higher-stakes, endowed goods) but is **not** a context-free constant. Use it as a *hypothesis to test for this customer*, not a guaranteed lever. Citable critique + defense in `references/replication-status.md`.

### 4. Bounded rationality & satisficing (Simon) vs ecological rationality (Gigerenzer)
- **Bounded rationality (Herbert Simon)** — real decision-makers have limited information, time, and compute, so they **satisfice**: pick the first option that clears an aspiration threshold rather than optimizing. Buyers rarely run an exhaustive vendor comparison; they stop at "good enough."
- **Ecological rationality / fast-and-frugal heuristics (Gerd Gigerenzer)** — the *counterpoint* to heuristics-and-biases. Simple heuristics (take-the-best, recognition, 1/N) are not defective; they are **adaptive** and often *more* accurate than complex models when information is scarce or uncertain (less-is-more). A heuristic is "rational" relative to its environment.
- *Why both matter:* heuristics-and-biases says shortcuts cause errors; ecological rationality says shortcuts are often the smart move. The truth is conditional — **match the diagnosis to the environment** before "fixing" a heuristic. These two research programs are the respective foundations of nudging and boosting (see §6).

### 5. Mental accounting & present bias
- **Mental accounting (Thaler)** — people sort money into non-fungible mental "buckets" (budget categories, "house money," the renewal-vs-new-purchase line) and violate fungibility. A spend framed against the "innovation budget" lands differently than the same spend against "BAU/maintenance."
- **Present bias / hyperbolic discounting** — near-term costs and rewards are discounted far more steeply than distant ones (quasi-hyperbolic β-δ; Laibson 1997). Produces **time-inconsistency**: plans made for the future ("we'll do the migration next quarter") get reversed when "next quarter" becomes "now." Upfront-cost / delayed-benefit work (migrations, upgrades, tech-debt paydown) is chronically under-chosen.

### 6. Choice architecture & nudges (Thaler & Sunstein, 2008) — and boosts (Hertwig)
- **Choice architecture** — every way options are presented (order, defaults, number, framing) influences choice. There is no neutral presentation, so design it deliberately.
- **Nudge** — "any aspect of the choice architecture that alters behavior predictably without forbidding options or significantly changing incentives" (libertarian paternalism). The most powerful nudge is the **default** (opt-out organ donation, auto-enrollment in 401(k)s, pre-checked plan tiers).
- **EAST** (UK Behavioural Insights Team) — make the desired action **Easy, Attractive, Social, Timely**. The most practical operator checklist.
- **MINDSPACE** — nine influences: **M**essenger, **I**ncentives, **N**orms, **D**efaults, **S**alience, **P**riming, **A**ffect, **C**ommitments, **E**go. (Note: the "Priming" element is on shaky empirical ground — see replication notes.)
- **Sludge** — choice architecture that adds *friction* against the person's own interest (cancellation mazes, hidden opt-outs). **Find and remove sludge in your own onboarding/renewal flows**; don't deploy it.
- **Boosts (Hertwig & Grüne-Yanoff, 2017)** — the *contrast* to nudges. Instead of steering the chooser, **build their competence** to choose well (teach a decision rule, give a fast-and-frugal tree, present risks as natural frequencies). Nudges rest on the heuristics-and-biases program; boosts rest on the simple-heuristics program. Boosts preserve agency and persist after the intervention; prefer them when the relationship is long-term and trust matters — which is most TAM work.
- *Ethics:* nudge toward the chooser's *own* interest, keep it transparent, and never sludge. A nudge a customer would resent if they saw it is a dark pattern. Full applied detail (incl. choice overload): `references/choice-architecture.md`.

### 7. Debiasing
You cannot will a bias away, but structured procedures help:
- **Consider-the-opposite** — actively list reasons the favored conclusion is wrong; the best-evidenced general debiaser, strong against anchoring and overconfidence.
- **Premortem (Gary Klein)** — *before* committing, imagine the project has failed and explain why; "prospective hindsight" surfaces risks that optimism suppresses. Cheap, fast, high-yield for plans and architecture decisions.
- **Reference-class forecasting (Flyvbjerg / Kahneman's "outside view")** — estimate from the distribution of *comparable past cases*, not from this case's inside details. The fix for the planning fallacy: replace assumptions with base rates.
- **Checklists** — force System 2 through a comprehensive, disciplined pass; only work with *consistent* adherence (partial use fails).

Procedures and when each applies: `references/debiasing-and-application.md`.

### 8. The two counterpoints: expert intuition (NDM) and the neural mechanism

The heuristics-and-biases core above is the *deficit* lens — shortcuts that misfire on lab tasks. Two
adjacent literatures push back and deepen it; both are loaded on demand from references.

- **Naturalistic Decision Making / Recognition-Primed Decision (Klein).** How *experts* decide under
  time pressure via pattern recognition + mental simulation (cue → pattern → action-script), evaluating
  one option serially rather than comparing many. This is §4's ecological rationality realized in the
  field: a fast gut call is trustworthy *only* in a **high-validity, well-practiced** environment with
  clear feedback (the **Kahneman–Klein 2009** truce) — and confidence is **not** a valid cue to
  accuracy. It also explains where the **premortem** (a Klein technique, §7) comes from and why it
  complements debiasing. Operator gate (trust the fast read vs slow down and debias):
  `references/naturalistic-decision-making.md`.
- **Neuroeconomics & reward.** The mechanism under the behavior: **dopamine reward-prediction error**
  (Schultz) as a learning signal that maps onto reinforcement learning; **subjective-value coding** in
  vmPFC and ventral striatum and the contested **common-currency** hypothesis; neural correlates of
  loss aversion, risk, and delay discounting; and **Damasio's somatic-marker** account of emotion as a
  decision *input*. Keep it mechanism-level and honest — **no reverse inference, no brain-scan sales
  tactics**, and a neural correlate never makes a biased decision rational:
  `references/neuroeconomics-and-reward.md`.

---

## Operator quick-map (bias → tell → move)

**How to use this table (it is a hypothesis generator, not a verdict).** Treat the
observed "tell" — what the customer said or did — as *data*, not as a confirmed
diagnosis. A tell suggests a *candidate* effect; confirm it against this specific person
and context before acting (people are heterogeneous, and several effects here are
context-dependent — see loss aversion in §3). If the tell is ambiguous or fits more than
one row, gather one more observation or ask a clarifying question before choosing a move.
Never state the bias label *to* the customer or imply they are irrational; the label is
your internal working hypothesis, the "move" is what you actually do.

| Situation / tell | Candidate effect (verify) | Operator move |
| --- | --- | --- |
| First price/number sets the whole conversation | Anchoring | Set the anchor first; if anchored against, re-anchor with your own reference before negotiating |
| "We've already sunk 18 months into this design" | Sunk-cost fallacy | Reframe to the *marginal* decision from today; make past spend explicitly irrelevant |
| "Let's just keep what we have / leave it as-is" | Status-quo / default bias | Make the better option the default; reduce switching friction; or set a decision deadline |
| Renewal framed only as new spend | Reference dependence / framing | Reframe against the reference point (cost of *losing* current capability, not net-new cost) |
| Team is 90% sure the timeline holds | Overconfidence / planning fallacy | Reference-class forecast + premortem; widen the interval |
| Customer fixates on the rare catastrophic risk | Availability + probability weighting | Provide base rates as **natural frequencies** ("3 in 1,000," not "0.3%") — a boost |
| Buyer stopped at the first "good enough" vendor | Satisficing (bounded rationality) | Don't assume full comparison happened; be the easy, salient option that clears the bar |
| You catch yourself collecting only confirming evidence | Confirmation bias | Consider-the-opposite; assign a devil's advocate |
| Post-incident "it was obviously going to fail" | Hindsight bias | In the postmortem, reconstruct what was *knowable at the time* (see `postmortem-writing`) |
| Signup/cancel flow has hidden friction | Sludge | Remove it; measure completion; opt-out only where it serves the user |

---

## Anti-patterns

- **Calling a customer "irrational."** They're predictably *boundedly* rational. Name the mechanism and design around it; never argue someone out of a bias by labeling it.
- **Citing a debunked effect.** Power posing, social priming, and ego depletion are contested/failed. The "2× loss-aversion constant" is over-stated. Check `references/replication-status.md` first.
- **Weaponizing nudges (sludge / dark patterns).** Steering a customer against their own interest is a trust-destroying short-term play; in a TAM relationship it's self-defeating. Prefer boosts.
- **Treating System 1/2 as literal neuroanatomy.** It's a model. Don't overclaim.
- **One-shot debiasing.** Awareness alone barely moves biases; only *structured procedures* (premortem, reference class, consider-the-opposite, checklists) reliably help — and only with disciplined use.
- **Over-applying loss aversion / "fixing" a heuristic that's actually ecologically rational.** Diagnose the environment first (Gigerenzer's caution).

---

## References (load on demand)

- `references/biases-catalog.md` — full heuristics-and-biases catalog with definitions, the canonical experiment, and a TAM application note each.
- `references/choice-architecture.md` — defaults, EAST, MINDSPACE, sludge vs nudge vs boost, choice overload, ethics of influence, applied to onboarding/pricing/renewal flows.
- `references/debiasing-and-application.md` — debiasing procedures (consider-the-opposite, premortem, reference-class forecasting, checklists) and worked operator scenarios (negotiation, forecasting, renewal, escalation).
- `references/replication-status.md` — what survived vs what's contested/failed (power posing, social priming, ego depletion, loss-aversion debate), with citable sources. **Read before citing any effect externally.**
- `references/naturalistic-decision-making.md` — the expert-intuition counterpoint: Klein's Recognition-Primed Decision model (cue→pattern→action-script, serial vs concurrent evaluation), NDM / macrocognition / sensemaking, the Kahneman–Klein (2009) conditions for trustworthy intuition (high-validity environment + practice with feedback), the premortem's origin, and an operator gate for fast read vs slow down.
- `references/neuroeconomics-and-reward.md` — the neural-mechanism counterpoint: dopamine reward-prediction error (Schultz) and RL, subjective-value coding in vmPFC/striatum and the common-currency debate, neural correlates of loss aversion / risk / delay discounting, Damasio's somatic-marker hypothesis, and the over-interpretation / reverse-inference caveats (mechanism-level only, no brain-scan tactics).

## Cross-references

- `behavior-change-psychology` — adjacent and complementary. *This* skill = the descriptive psychology of a **decision** (biases, framing, choice architecture). *That* skill = changing **behavior over time** (motivation, Fogg B=MAP, stages-of-change, habit loops, adoption/enablement). "Design an onboarding nudge to drive adoption" → that skill; "what default/framing will shape this purchase decision" → this skill.
- `da-33-prescriptive-analytics` — the **normative** counterpart (optimal action under constraints: LP/MILP, decision trees, EVPI, utility theory). Use it to compute the optimum; use this skill to understand why humans deviate from it.
- `executive-comms` — persuasion and decision-driving *communication* (board memos, negotiation prep, influence). Framing overlaps; for the rhetoric/persuasion craft go there, for the underlying decision psychology stay here.
- `postmortem-writing` — applies hindsight-bias control in incident reviews.
- `deep-research-methods` — covers confirmation bias / echo chambers as research anti-patterns.

## Sources

- Tversky, A. & Kahneman, D. (1974). "Judgment under Uncertainty: Heuristics and Biases." *Science* 185(4157), 1124–1131.
- Kahneman, D. & Tversky, A. (1979). "Prospect Theory: An Analysis of Decision under Risk." *Econometrica* 47(2), 263–291.
- Tversky, A. & Kahneman, D. (1991). "Loss Aversion in Riskless Choice: A Reference-Dependent Model." *QJE* 106(4), 1039–1061.
- Kahneman, D. (2011). *Thinking, Fast and Slow.*
- Simon, H. A. (1955/1956). Bounded rationality and satisficing.
- Gigerenzer, G. & colleagues — fast-and-frugal heuristics / ecological rationality (ABC Research Group).
- Thaler, R. & Sunstein, C. (2008/2021). *Nudge* (and *Nudge: The Final Edition*) — choice architecture, defaults, sludge.
- Laibson, D. (1997). "Golden Eggs and Hyperbolic Discounting." *QJE*. (present bias / quasi-hyperbolic β-δ)
- Hertwig, R. & Grüne-Yanoff, T. (2017). "Nudging and Boosting: Steering or Empowering Good Decisions." *Perspectives on Psychological Science* 12(6), 973–986.
- Dolan, P. et al. (2010). MINDSPACE; Behavioural Insights Team (2014). EAST.
- Gal, D. & Rucker, D. (2018). "The Loss of Loss Aversion: Will It Loom Larger Than Its Gain?" *Journal of Consumer Psychology*.
- Replication: Open Science Collaboration (2015) *Science*; Many Labs 2; Ranehill et al. (2015) and Simmons & Simonsohn (2017) on power posing.
