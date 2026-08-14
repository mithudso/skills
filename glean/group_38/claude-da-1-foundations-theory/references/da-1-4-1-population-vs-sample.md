<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-4-1-population-vs-sample` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-4-1-population-vs-sample
description: >-
  Foundational statistical-inference distinction between a population (the whole
  group you want conclusions about) and a sample (the subset you actually
  measure), plus the parameter-vs-statistic pairing, target population vs
  sampling frame, sampling error, and what makes a sample representative.
  TRIGGER: a user asks "population vs sample", what a parameter or statistic is,
  whether a number describes a population or a sample, how to define a target
  population or sampling frame, why we sample instead of taking a census, when a
  sample is representative, or how sample statistics estimate population
  parameters in inferential statistics. SKIP: choosing a concrete sampling
  technique like stratified/cluster/systematic design (defer to a sampling-
  methods skill); the shape and variability of a statistic across many samples
  i.e. sampling distributions and standard error (defer to
  da-1-4-2-sampling-distributions-standard-error); descriptive summary of one
  dataset with no inference goal (defer to descriptive-statistics); experimental
  vs observational design or causal inference.
---

# Population vs. Sample

Taxonomy: Data Analysis > Foundations & Theory > Statistical inference foundations > Population vs. sample.

This is the entry concept of statistical inference. Everything in inference rests on one move: you cannot measure the whole group, so you measure a subset and reason backward to the whole group. This skill covers that distinction precisely, the vocabulary that goes with it (parameter vs. statistic, target population vs. sampling frame), and the conditions under which the backward reasoning is valid.

## Core definitions

- **Population** — the complete set of every item, person, or unit you want to draw conclusions about. It can be finite (every registered voter in a country) or effectively a theoretical construct (every car a factory will ever produce). Its size is written **N**. [Statistics By Jim — Populations, Parameters, and Samples](https://statisticsbyjim.com/basics/populations-parameters-samples-inferential-statistics/)
- **Sample** — a subset of the population that you actually observe and measure, because measuring the entire population is impractical, too expensive, destructive, or impossible. Its size is written **n**, with n < N. [Scribbr — Population vs. Sample](https://www.scribbr.com/methodology/population-vs-sample/)
- **Census** — measuring every member of the population (n = N). A sample is what you use when a census is out of reach.
- **Parameter** — a number that describes a *population*. Almost always unknown in practice. [LibreTexts — Parameters vs. Statistics](https://stats.libretexts.org/Courses/Lumen_Learning/Concepts_in_Statistics_(Lumen)/07:_Linking_Probability_to_Statistical_Inference/7.03:_Parameters_vs._Statistics)
- **Statistic** — a number that you calculate from a *sample*. Known, because you computed it. Used as an estimate of the corresponding parameter. [LibreTexts](https://stats.libretexts.org/Courses/Lumen_Learning/Concepts_in_Statistics_(Lumen)/07:_Linking_Probability_to_Statistical_Inference/7.03:_Parameters_vs._Statistics)

Mnemonic: **p**arameter goes with **p**opulation; **s**tatistic goes with **s**ample.

## Notation (parameter ↔ statistic)

| Measure | Population parameter | Sample statistic |
|---|---|---|
| Mean | μ (mu) | x̄ (x-bar) |
| Standard deviation | σ (sigma) | s |
| Proportion | p | p̂ (p-hat) |
| Size | N | n |

Greek letters and capitals generally denote the (unknown) population side; Latin letters and lowercase generally denote the (computed) sample side. [LibreTexts — Parameters vs. Statistics](https://stats.libretexts.org/Courses/Lumen_Learning/Concepts_in_Statistics_(Lumen)/07:_Linking_Probability_to_Statistical_Inference/7.03:_Parameters_vs._Statistics)

## Why this matters for inference

Inferential statistics is the act of using a known sample statistic to estimate an unknown population parameter. You compute x̄ from your sample and use it as your best estimate of μ; you compute p̂ and use it to estimate p. The whole apparatus of confidence intervals and hypothesis tests exists to quantify how good that estimate is. [Statistics By Jim](https://statisticsbyjim.com/basics/populations-parameters-samples-inferential-statistics/)

The estimate is almost never exactly right. The gap between a sample statistic and the true population parameter is **sampling error** — the unavoidable variation that comes from observing only part of the whole. Sampling error is not a mistake; it is inherent to sampling and shrinks as n grows. (Note: *how* that error is distributed across repeated samples is the next concept, sampling distributions and standard error — defer that detail there.) [Scribbr — Population vs. Sample](https://www.scribbr.com/methodology/population-vs-sample/)

## Refining "population": target population vs. sampling frame

"Population" is not a single thing in practice. Three layers matter, and the gaps between them are where inference quietly breaks.

- **Target population** — the group about which you actually want to make inferences and estimates. The conceptual answer to "who is this study about?" [The Analysis Factor — Target Population and Sampling Frame](https://www.theanalysisfactor.com/target-population-sampling-frame/)
- **Sampling frame** — the concrete, operational list of units you can actually draw from (a voter roll, a customer database, a list of phone numbers). Ideally it coincides exactly with the target population; in practice it rarely does. [The Analysis Factor](https://www.theanalysisfactor.com/target-population-sampling-frame/)
- **Sample** — the units actually selected from the frame and measured.

Two mismatch types between target population and frame:
- **Out of coverage / undercoverage** — target-population members missing from the frame (e.g., adults with no phone are absent from a phone-number frame).
- **Out of scope** — frame entries that are not in the target population (e.g., children who own phone numbers when the target population is adults).

When the people in these gaps differ systematically from those in the frame on the outcome you measure, the result is **coverage bias** — a sub-type of selection bias. Your sample then represents the *sampled population* (whoever the frame could reach), not the target population you meant to study. [The Analysis Factor — Target Population and Sampling Frame](https://www.theanalysisfactor.com/target-population-sampling-frame/)

## What makes a sample usable: representativeness

A sample supports inference only if it is **representative** — its relevant characteristics mirror the population's. The standard route to representativeness is a **probability sampling method**, in which every population member has a known, non-zero chance of selection. Random selection is what licenses generalizing the sample result back to the population and what makes sampling error behave predictably. [Statistics By Jim](https://statisticsbyjim.com/basics/populations-parameters-samples-inferential-statistics/)

Without random selection you get **sampling bias**: some members are systematically more likely to be chosen than others, so the statistic is a biased estimate of the parameter no matter how large n is. A bigger non-random sample does not fix bias — it just gives you a more precise estimate of the wrong number. [Scribbr — Population vs. Sample](https://www.scribbr.com/methodology/population-vs-sample/)

(Picking among specific probability designs — simple random, stratified, cluster, systematic — is a separate concept; defer the design choice to a sampling-methods skill.)

## Why sample instead of census

- **Cost and time** — a census of a large population is usually infeasible.
- **Practicality / access** — many populations cannot be fully enumerated.
- **Destructive measurement** — testing every unit (e.g., crash-testing every car, tasting every batch) would destroy the population.
- **Sufficiency** — a well-drawn sample estimates the parameter closely enough that a census adds little. [Scribbr — Population vs. Sample](https://www.scribbr.com/methodology/population-vs-sample/)

## Worked example

Goal: estimate the mean annual income (μ) of all 200,000 adult residents of a city (the **target population**, N = 200,000).

1. **Frame.** You obtain the city tax-filer list. Adults who filed no return are out of coverage; deceased filers still listed are out of scope — both gaps to note.
2. **Sample.** You draw a simple random sample of n = 500 from the frame. Random selection makes the sample representative and gives every filer a known selection chance.
3. **Statistic.** You compute the sample mean income x̄ = \$58,200. This is a *statistic* — known, computed from the sample.
4. **Inference.** You use x̄ = \$58,200 as your point estimate of the *parameter* μ, the true mean income of all adult residents, which you can never observe directly.
5. **Sampling error.** The true μ is almost certainly not exactly \$58,200; the difference is sampling error, inherent because you measured 500 of 200,000.
6. **Bias check.** Because non-filers (often low or no income) are out of coverage, x̄ may systematically overstate μ for *all* adults — a coverage-bias warning, distinct from sampling error.

## Common pitfalls

- **Calling a sample number a parameter.** Anything you computed from observed data is a statistic, not a parameter. Parameters are the unknowns you are estimating. [LibreTexts](https://stats.libretexts.org/Courses/Lumen_Learning/Concepts_in_Statistics_(Lumen)/07:_Linking_Probability_to_Statistical_Inference/7.03:_Parameters_vs._Statistics)
- **Treating the sampling frame as the target population.** Inferences are only valid for whoever the frame could actually reach. State the gap. [The Analysis Factor](https://www.theanalysisfactor.com/target-population-sampling-frame/)
- **Confusing sampling error with bias.** Sampling error is random and shrinks with n; bias is systematic and does not. A larger biased sample is still biased. [Scribbr](https://www.scribbr.com/methodology/population-vs-sample/)
- **Assuming size guarantees representativeness.** Representativeness comes from random selection and frame coverage, not from a large n. [Statistics By Jim](https://statisticsbyjim.com/basics/populations-parameters-samples-inferential-statistics/)
- **Vague population definition.** If you cannot state the target population precisely, you cannot say what your statistic estimates or to whom results generalize.

## Sources

1. Scribbr — Population vs. Sample: Definitions, Differences & Examples. https://www.scribbr.com/methodology/population-vs-sample/
2. Statistics By Jim — Populations, Parameters, and Samples in Inferential Statistics. https://statisticsbyjim.com/basics/populations-parameters-samples-inferential-statistics/
3. LibreTexts (Lumen, Concepts in Statistics) — 7.3 Parameters vs. Statistics. https://stats.libretexts.org/Courses/Lumen_Learning/Concepts_in_Statistics_(Lumen)/07:_Linking_Probability_to_Statistical_Inference/7.03:_Parameters_vs._Statistics
4. The Analysis Factor — Target Population and Sampling Frame in Survey Sampling. https://www.theanalysisfactor.com/target-population-sampling-frame/
