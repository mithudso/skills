<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-4-3-frequentist-vs-bayesian-paradigms` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-4-3-frequentist-vs-bayesian-paradigms
description: >-
  The two competing schools of statistical inference — frequentist and Bayesian —
  contrasted as foundations of inference: what probability *means* in each, whether
  a parameter is a fixed unknown constant or a random variable with a distribution,
  the role (or absence) of a prior, and how each paradigm answers the same question
  with different objects (confidence interval vs credible interval, p-value vs
  posterior probability). Covers the philosophical split, the shared Bayes-theorem
  machinery one side uses and the other rejects for inference, when each is
  preferred, and the misinterpretations that arise from mixing them.
  TRIGGER: "frequentist vs Bayesian", "which statistical paradigm/approach should I
  use", "is probability a long-run frequency or a degree of belief", "are parameters
  fixed or random", "do I need a prior", "confidence interval vs credible interval
  interpretation", "p-value vs posterior probability", "why can't I say there's a
  95% chance the parameter is in this confidence interval", choosing an inference
  philosophy for an analysis or A/B test.
  SKIP: Bayes' theorem as a probability identity (deriving/applying P(A|B)) ->
  da-1-3-5-bayes-theorem; conditional/joint/marginal probability mechanics ->
  da-1-3-4; what a sampling distribution or standard error *is* ->
  da-1-4-2-sampling-distributions-standard-error; population vs sample as a
  conceptual contrast -> da-1-4-1; building/interpreting a specific confidence
  interval, hypothesis test, or p-value as a procedure (later inference nodes);
  the umbrella overview of statistical inference -> da-1-4-statistical-inference-foundations.
  Defer those framings rather than absorbing them here.
---

# Frequentist vs. Bayesian paradigms

Scope: Data Analysis > Foundations & Theory > Statistical inference foundations >
Frequentist vs. Bayesian paradigms. This skill explains the *foundational choice*
an analyst makes before running any inference: what probability means and whether
the unknown thing you want to learn is a fixed constant or a random variable. It
deliberately stops short of executing any single procedure (a specific confidence
interval, hypothesis test, or posterior computation) and of re-deriving Bayes'
theorem itself; those are separate nodes (see SKIP cues above).

## 1. The split in one sentence

The two paradigms disagree about **what probability is**, and that single
disagreement cascades into different definitions of parameters, intervals, and
evidence. Frequentists define probability as the long-run frequency of an event in
repeated trials; Bayesians define it as a degree of belief that can be assigned to
any proposition, including a hypothesis [Wikipedia, "Frequentist inference";
Probabilistic World, "Frequentist and Bayesian Approaches in Statistics"].

## 2. The frequentist paradigm

**Probability = long-run frequency.** "Only repeatable random events have
probabilities, and these probabilities are equal to the long-term frequency of
occurrence of the events in question" [Probabilistic World]. You cannot assign a
probability to a one-off proposition such as "this coin's bias is 0.6" — that
number either is or is not 0.6.

**Parameters are fixed but unknown constants.** "In frequentist statistics,
probability represents long-run frequencies in repeated experiments, and parameters
are fixed but unknown constants" [search synthesis of frequentist sources; see
Lovrić, "Statistical Paradigms"]. The *data* are random; the parameter is not.

**Inference is about the procedure, not the parameter.** Because the parameter is
fixed, frequentists make probability statements about the *method* — how it would
behave across many hypothetical repetitions of the experiment. "Inference is based
on the sampling distribution — how results would behave if we repeated the
experiment many times" [search synthesis; ScienceDirect, "Frequentist Approach"].

**Core tools:** point estimation (often maximum likelihood), confidence intervals,
hypothesis tests, and p-values, "with a focus on controlling error rates over many
hypothetical repetitions" [search synthesis].

**The p-value** is "the probability, under the assumption of the null hypothesis,
to obtain a result equal to or more extreme than observed" [Berger, "Could Fisher,
Jeffreys and Neyman Have Agreed on Testing?"]. It is a statement about data given a
hypothesis, **not** about the hypothesis given the data.

**Two founders, two flavours (a frequent confusion).** Frequentist testing was
built mostly by R.A. Fisher and by Jerzy Neyman with Egon Pearson, and they did not
agree. Fisher's significance tests treat the p-value as a continuous measure of
evidence against the null; the Neyman–Pearson framework treats testing as a binary
decision (accept/reject) with pre-set error rates (alpha, beta) controlled over the
long run [Berger]. "Modern practice awkwardly blends both Fisher and Neyman–Pearson
approaches" — the source of much p-value misuse [search synthesis].

## 3. The Bayesian paradigm

**Probability = degree of belief.** Bayesians "embrace subjective probability,
treating probabilities as degrees of belief rather than only long-run frequencies"
[Wikipedia, "Bayesian inference"]. Any proposition — including a hypothesis or a
parameter value — can carry a probability.

**Parameters are random variables.** A parameter is described by a probability
distribution that encodes your uncertainty about it, "rather than fixed unknowns"
[Wikipedia, "Bayesian inference"]. Here the *parameter* is treated as random and the
observed data are fixed (conditioned on).

**Inference updates belief via Bayes' theorem.** Bayesian inference "is a method of
statistical inference in which Bayes' theorem is used to calculate a probability of
a hypothesis, given prior evidence, and update it as more information becomes
available" [Wikipedia, "Bayesian inference"]. The machine:

```
posterior  ∝  likelihood × prior

            P(E | H) · P(H)
P(H | E) = -----------------
                 P(E)
```

- **Prior** P(H): "the estimate of the probability of the hypothesis H before the
  data E ... is observed" [Wikipedia].
- **Likelihood** P(E|H): "the probability of observing E given H" — the
  compatibility of the evidence with the hypothesis [Wikipedia].
- **Posterior** P(H|E): "the probability of H given E, i.e., after E is observed.
  This is what we want to know" [Wikipedia].
- **Marginal likelihood / evidence** P(E): the data distribution marginalized over
  the parameter; the normalizing constant [Wikipedia].

Note the relationship to node 1.3.5: Bayes' theorem is the same probability
identity in both worlds, but the frequentist uses it only for events with genuine
frequencies, while the Bayesian extends it to hypotheses and parameters. The
*identity* belongs to da-1-3-5; the *paradigm-level decision to apply it to
parameters* belongs here.

**Conjugate priors** are priors chosen so that "the corresponding posterior
distribution will be in the same family, and the calculation may be expressed in
closed form" [Wikipedia, "Bayesian inference"] — a practical convenience, not a
requirement.

## 4. The decisive contrast: intervals

This is where the philosophical split becomes concrete and where analysts most
often slip.

| | Frequentist confidence interval | Bayesian credible interval |
|---|---|---|
| What is random | the interval (it depends on the random sample) | the parameter |
| What is fixed | the parameter | the observed data |
| Uses a prior | no | yes |
| Correct reading | "if we repeated the experiment many times, 95% of the intervals built this way would contain the true value" | "given the data and prior, there is a 95% probability the parameter lies in this interval" |

"In frequentist terms, the parameter is fixed and the confidence interval is random
... whereas credible intervals have values with a probability density, representing
the plausibility that the parameter has those values" [Wikipedia, "Credible
interval"]. A frequentist "can only say that a particular realized interval either
does or does not contain the true parameter value" [search synthesis from
thestatsgeek / Wikipedia].

**The #1 misinterpretation:** reading a 95% *confidence* interval as "95% chance the
parameter is in here." That is the *credible*-interval (Bayesian) reading and is
invalid under frequentist semantics, because the frequentist parameter is not a
random variable and has no probability distribution. The two intervals can even
coincide numerically (e.g., a flat prior with a normal likelihood) yet still mean
different things.

## 5. When to prefer which

No universal winner; the choice follows the question and the available information
[Statsig, "Bayesian or Frequentist: Choosing your statistical approach"].

**Lean frequentist when:**
- You need objective, prior-free procedures with guaranteed long-run error rates
  (regulatory trials, fixed-sample experiments).
- There is no credible prior information, or you must avoid the appearance of
  injecting belief.
- Sample sizes are large and the questions are standard (means, proportions).

**Lean Bayesian when:**
- You have genuine prior information worth incorporating (Bayesian analysis
  "incorporates prior information into the analysis, whereas a frequentist analysis
  is purely driven by the data" [search synthesis]).
- You want a direct probability statement about the hypothesis or parameter (the
  intuitive reading people *want* from a confidence interval).
- You update sequentially as data arrive (e.g., continuous A/B-test monitoring),
  where frequentist fixed-sample error control is awkward [Amplitude, "Frequentist
  vs. Bayesian ... A/B Testing"].

## 6. Worked contrast: estimating a coin's bias

You flip a coin 10 times and see 7 heads. Question: what is the bias θ?

- **Frequentist.** θ is a fixed constant. The estimate is θ̂ = 7/10 = 0.7
  (maximum likelihood). A 95% confidence interval is built from the sampling
  distribution of the proportion; its correct reading is about the *procedure's*
  long-run coverage, not about θ having a 95% chance of being inside.
- **Bayesian.** θ is a random variable. Start with a prior (e.g., a uniform
  Beta(1,1) = "no idea"). The likelihood is Binomial. With a Beta prior the
  posterior is Beta (a conjugate pair): Beta(1+7, 1+3) = Beta(8,4), posterior mean
  8/12 ≈ 0.667. A 95% *credible* interval read off that posterior legitimately
  means "95% probability θ is in here, given the prior and data." A different prior
  (e.g., Beta(50,50) from strong belief the coin is fair) pulls the posterior toward
  0.5 — the prior visibly matters, which is the feature Bayesians want and
  frequentists distrust.

## 7. Pitfalls

- **Mixing the interpretations.** Using a frequentist tool but stating a Bayesian
  conclusion (the confidence-interval misreading in §4) is the most common error.
- **"Bayesian = subjective, so unscientific."** Overstated. Priors are explicit and
  can be weak/uninformative; with enough data the likelihood dominates and the two
  paradigms often converge numerically.
- **"Frequentist = objective, so prior-free is bias-free."** Also overstated. Model
  choice, stopping rules, and the choice of test statistic inject assumptions even
  without a named prior.
- **P-value as P(null is true).** The p-value is P(data this extreme | null), not
  P(null | data); only the Bayesian posterior gives the latter [Hubbard, "P Values
  are not Error Probabilities"].
- **Treating the paradigms as the only two.** Likelihood and fiducial paradigms also
  exist [Lovrić, "Statistical Paradigms"]; this node focuses on the dominant pair.
- **Optional stopping under frequentist tests.** Peeking at data and stopping when
  significant inflates the false-positive rate under frequentist error control — a
  reason teams sometimes switch to Bayesian monitoring for A/B tests.

## 8. Quick decision checklist

1. What does "probability" mean for this problem — a repeatable frequency, or a
   belief about a one-off quantity? (Frequency → frequentist; belief → Bayesian.)
2. Is the unknown a fixed constant or something you want a distribution over?
3. Do you have, and want to use, prior information?
4. Do you need a direct probability statement about the hypothesis (Bayesian) or a
   procedure with controlled long-run error rates (frequentist)?
5. Will data arrive sequentially with interim looks? (Favours Bayesian.)
6. Are there regulatory/convention constraints fixing the paradigm for you?

## Sources

- Wikipedia, "Bayesian inference" —
  https://en.wikipedia.org/wiki/Bayesian_inference
- Wikipedia, "Frequentist inference" —
  https://en.wikipedia.org/wiki/Frequentist_inference
- Wikipedia, "Credible interval" —
  https://en.wikipedia.org/wiki/Credible_interval
- Probabilistic World, "Frequentist and Bayesian Approaches in Inferential
  Statistics" —
  https://www.probabilisticworld.com/frequentist-bayesian-approaches-inferential-statistics/
- Statsig, "Bayesian or Frequentist: Choosing your statistical approach" —
  https://www.statsig.com/perspectives/bayesian-or-frequentist-choosing-your-statistical-approach
- M. M. Lovrić, "Statistical Paradigms: Frequentist, Bayesian, Likelihood &
  Fiducial" (Radford University) —
  https://sites.radford.edu/~mlovric/Statistical_Paradigms.html
- J. O. Berger, "Could Fisher, Jeffreys and Neyman Have Agreed on Testing?" (Duke) —
  https://www2.stat.duke.edu/~berger/papers/02-01.pdf
- R. Hubbard, "P Values are not Error Probabilities" —
  https://www.uv.es/sestio/TechRep/tr14-03.pdf
- Amplitude, "Frequentist vs. Bayesian: Comparing Statistics Methods for A/B
  Testing" — https://amplitude.com/blog/frequentist-vs-bayesian-statistics-methods
