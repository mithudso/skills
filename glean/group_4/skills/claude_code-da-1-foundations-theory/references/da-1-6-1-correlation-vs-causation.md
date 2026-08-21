<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-6-1-correlation-vs-causation` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-6-1-correlation-vs-causation
description: >-
  Epistemology of data analysis: why a statistical association between two
  variables is not evidence that one causes the other, and what it takes to move
  from "they move together" to "one drives the other." Covers confounding, the
  directionality problem, spurious/chance correlations, the Bradford Hill
  viewpoints, the confounder/collider/mediator distinction in causal diagrams,
  and the study designs (RCTs, instrumental variables, difference-in-differences,
  natural experiments) that license causal claims.
  TRIGGER: someone reasons "X correlates with Y, so X causes Y"; interpreting a
  dashboard/regression/A-B result as if it proves causation; reviewing a claim
  that an intervention "drove" an outcome; asking "does this data prove X causes
  Y?", "is this a confounder?", "is this correlation spurious?", "how do we know
  it's causal?"; deciding whether an observed effect justifies acting on it.
  SKIP: the mechanics of computing a correlation coefficient or its hypothesis
  test (defer to a correlation/regression statistics skill); building or
  estimating a full causal-inference model end to end such as propensity scores,
  do-calculus derivations, or IV estimation code (defer to a dedicated causal
  inference / econometrics skill); experimental design mechanics like power and
  randomization procedure (defer to an experiment-design skill); reproducibility
  vs replicability as a concept (defer to da-1-6-3-reproducibility-replicability);
  inductive vs deductive reasoning as a concept (defer to
  da-1-6-2-inductive-vs-deductive-reasoning).
---

# Correlation vs. Causation

Taxonomy slot: Data Analysis > Foundations & Theory > Epistemology of data >
Correlation vs. causation. This skill is about the *epistemic* gap — how we
come to know whether an association reflects a causal relationship — not about
the arithmetic of correlation or the full machinery of causal estimation.

## The core claim

A correlation (or any statistical association) measures that two variables tend
to move together. It says nothing, on its own, about *why*. "Correlation does
not imply causation" is the warning that an observed association is consistent
with several distinct underlying structures, only one of which is "X causes Y"
[Statistics By Jim, correlation-vs-causation; Wikipedia, "Correlation does not
imply causation"].

When you see corr(X, Y) ≠ 0, at least these explanations remain open:

1. **X causes Y** (the claim being tempted).
2. **Y causes X** (reverse causation / the directionality problem).
3. **A third variable Z causes both X and Y** (confounding — the classic
   "lurking variable").
4. **Chance / coincidence** (a spurious correlation with no mechanism at all).
5. **Selection or collider bias** — the association was manufactured by *how the
   data was filtered or conditioned*, not by any causal link.
6. **A causal chain through a mediator** (X → M → Y) — causal, but the direct
   path may be misread.

The whole discipline of moving from correlation to causation is the work of
ruling out 2–5 (and characterizing 6).

## Why the gap exists (the four failure modes)

### 1. Confounding (the third-variable problem)
A **confounder** is a variable that causally influences *both* the supposed
cause and the supposed effect, creating an association between them even when
neither causes the other [Statistics By Jim, correlation-vs-causation; Cambridge
/ PMC11588567, "Methods in causal inference. Part 1"].

Canonical example: ice cream sales correlate with shark attacks. Neither causes
the other; **hot weather / beach attendance** drives both. Control for the
weather and the association collapses [Statistics By Jim].

Confounding is *the* reason observational associations mislead. A confounder
"cannot be defined purely in terms of statistical associations alone" — you need
a causal story (which variable is a common cause) to know that Z is a confounder
rather than something else [PMC5841844, "Causal inference — so much more than
statistics"].

### 2. The directionality (reverse-causation) problem
Even if X and Y are causally linked, the arrow may point the other way. Stress
correlates with insomnia; does stress cause sleeplessness, does sleeplessness
cause stress, or both? Correlation is symmetric — corr(X,Y) = corr(Y,X) — so it
carries no directional information. Temporal order is the usual tie-breaker:
a cause must precede its effect [Wikipedia, "Correlation does not imply
causation"; Statistics By Jim, causation].

### 3. Spurious / chance correlation
A **spurious correlation** is an association that appears by pure chance with no
mechanism linking the variables. With enough variables and enough data series,
some will correlate strongly by coincidence — Tyler Vigen's "Spurious
Correlations" catalogs absurd high-correlation pairs (e.g., cheese consumption
vs. deaths by bedsheet entanglement). There is **no statistical test that
detects a spurious correlation**; the guard is a plausible mechanism plus
out-of-sample replication [Statistics By Jim, spurious-correlation; Tyler Vigen,
tylervigen.com]. Multiple-comparison / data-dredging settings manufacture these
at scale.

### 4. Selection and collider bias
A **collider** is a variable caused by *both* X and Y (X → C ← Y).
Counterintuitively, conditioning on a collider — filtering, stratifying, or
selecting your sample on it — *creates* a spurious association between X and Y
that does not exist in the full population [PMC11588567]. "If both A and Y are
positively associated with the collider L, then within a stratum of L, A and Y
become negatively associated" even with no causal link between them
[PMC11588567]. This is why sampling on the outcome, or studying only a
self-selected subgroup, can invert or invent relationships.

## Confounder vs. collider vs. mediator — and what to adjust for

These three roles look similar in data but demand *opposite* handling. Getting
this wrong is the most common technical error in correlation→causation work
[PMC11588567].

| Role | Structure | Effect on naive correlation | Adjust for it? |
|---|---|---|---|
| **Confounder** | Z → X and Z → Y (common cause) | Creates a spurious association | **Yes** — adjusting removes the bias |
| **Collider** | X → C ← Y (common effect) | No bias *until* you condition on it | **No** — conditioning *induces* bias |
| **Mediator** | X → M → Y (on the causal path) | Carries the real effect | **No** (for total effect) — adjusting blocks the real effect and understates causation |

Pearl's **backdoor criterion** formalizes the rule: to identify the causal
effect of X on Y, adjust for a set L that (a) contains no descendant of X and
(b) blocks every "backdoor path" — every non-causal path from X to Y with an
arrow pointing *into* X. A correctly chosen L removes confounding without
opening collider paths [PMC11588567]. The actionable discipline:

- Adjust for **pre-treatment common causes** (confounders).
- Do **not** adjust for variables that the treatment *caused* (mediators or
  colliders); conditioning on post-treatment variables is a frequent source of
  selection bias [PMC11588567].
- You cannot decide any of this from the correlation matrix alone — it requires
  a causal diagram (DAG) expressing your assumptions about the arrows.

## How to actually establish causation

### Gold standard: randomization
A **randomized controlled trial (RCT)** assigns the treatment by chance, so on
expectation the treated and untreated groups differ only in the treatment.
Randomization breaks the link between treatment and *all* confounders
(measured and unmeasured), which is why it licenses causal claims that
observational correlation cannot [Statistics By Jim, causation; IZA World of
Labour, "Using instrumental variables to establish causality"]. When an RCT is
unethical, impractical, or impossible, you fall back to quasi-experimental
designs.

### Quasi-experimental / observational tools
- **Instrumental variables (IV):** find a variable that affects the treatment
  but influences the outcome *only through* the treatment, and is unrelated to
  confounders. This exploits natural variation to recover a causal effect from
  non-experimental data [IZA World of Labour; PMC2905668, "Instrumental
  variables I"].
- **Natural experiments:** events outside the researcher's control assign
  treatment in a near-random way, approximating an RCT [Statistical Tools for
  Causal Inference, Ch. 4 "Natural Experiments"].
- **Difference-in-differences (DiD):** compare the before/after change in a
  treated group against the before/after change in an untreated group, netting
  out fixed group differences and common time trends [Statistical Tools for
  Causal Inference].

Each rests on assumptions that are *not testable from the data alone* (exclusion
restriction for IV, parallel trends for DiD, ignorable assignment for matching).
The causal claim is only as good as those assumptions.

### The Bradford Hill viewpoints (weighing observational evidence)
When experiments are unavailable, Sir Austin Bradford Hill's 1965 framework
offers nine *viewpoints* (he deliberately avoided the word "criteria") for
judging whether an association is causal [Wikipedia, "Bradford Hill criteria";
Statistics By Jim, causation; PMC8206235, "Assessing causality in epidemiology"]:

1. **Strength** — a larger effect is harder to explain away by bias.
2. **Consistency** — reproduced by different people, places, samples.
3. **Specificity** — a specific cause maps to a specific effect.
4. **Temporality** — the cause precedes the effect. *The one viewpoint
   universally treated as necessary* [Statistics By Jim, causation].
5. **Biological gradient (dose-response)** — more exposure, more effect.
6. **Plausibility** — a credible mechanism exists.
7. **Coherence** — fits other known facts / lab evidence.
8. **Experiment** — experimental or intervention evidence supports it.
9. **Analogy** — similar causes produce similar effects elsewhere.

Hill stressed that **none is required as a sine qua non** and the set is "not a
checklist": a true cause may miss several viewpoints, and a confounder can
satisfy all nine. They support judgment, not a mechanical verdict [Wikipedia,
"Bradford Hill criteria"; PMC8206235].

## Worked walk-through

Claim: "Customers who use Feature F churn 30% less, so we should push everyone to
Feature F."

1. **Name the association.** corr(uses F, retained) is strongly positive.
2. **Enumerate alternatives.**
   - *Reverse causation:* engaged customers (already unlikely to churn) adopt F;
     F doesn't cause retention, retention-prone behavior causes F adoption.
   - *Confounding:* "overall product engagement" is a common cause — it drives
     both F adoption and retention. Z → X and Z → Y.
   - *Collider/selection:* if you only analyzed customers who renewed at least
     once, you conditioned on a downstream variable and may have induced bias.
3. **Apply the rules.** Engagement is a pre-treatment common cause → adjust for
   it. Do not adjust for anything F itself caused.
4. **Seek a stronger design.** Randomly prompt a subset to adopt F (an RCT or
   A/B test). If randomization is impossible, look for an instrument (e.g., a UI
   change that nudged F adoption for reasons unrelated to engagement) or a DiD
   around a rollout date.
5. **Weigh viewpoints.** Is there a dose-response (more F use → less churn)? A
   plausible mechanism? Temporal precedence (F use *before* the retention
   window)? Consistency across segments?
6. **Verdict.** Until confounding by engagement is broken — ideally by
   randomization — the 30% figure is an association, not a license to claim F
   *causes* retention.

## Pitfalls and red flags

- **"Statistically significant, therefore causal."** Significance addresses
  chance (failure mode 4) only; it is silent on confounding, direction, and
  selection [Statistics By Jim, correlation-vs-causation].
- **Controlling for a mediator or collider** and reporting the shrunken/inflated
  coefficient as "the causal effect." Adjusting for post-treatment variables
  biases the estimate [PMC11588567].
- **Treating Bradford Hill as a scoring checklist** ("met 7 of 9, so causal").
  Hill explicitly rejected this [Wikipedia, "Bradford Hill criteria"].
- **Data-dredging.** Testing many pairs and reporting the ones that correlate
  guarantees spurious hits [Statistics By Jim, spurious-correlation].
- **Ignoring temporality.** If the "effect" was measured before or alongside the
  "cause," the causal story is dead on arrival.
- **Assuming "no confounder I can think of" = "no confounding."** Unmeasured
  confounders are the rule, not the exception, in observational data — which is
  precisely why randomization is valued [IZA World of Labour].
- **Stating a causal verb ("drove", "caused", "led to", "boosted") from
  observational correlation.** Match the verb to the evidence: "is associated
  with" for correlation, "causes/increases" only when a design licenses it.

## Quick reference

- Correlation = co-movement. Causation = intervening on X changes Y.
- Five things an association can be: cause, reverse-cause, confounded,
  coincidence, selection artifact (plus the causal-chain/mediator case).
- Adjust for **confounders**, never for **colliders** or **mediators**.
- Only **temporality** is non-negotiable among the Hill viewpoints.
- The strongest correlation→causation move is **randomization**; absent it, use
  IV / DiD / natural experiments and state the (untestable) assumptions.
- You cannot read causation off a correlation matrix; you need a causal model.

## Sources

1. Statistics By Jim — "Correlation vs Causation: Understanding the Differences."
   https://statisticsbyjim.com/basics/correlation-vs-causation/
2. Statistics By Jim — "Causation in Statistics: Hill's Criteria."
   https://statisticsbyjim.com/basics/causation/
3. Statistics By Jim — "Spurious Correlation: Definition, Examples & Detecting."
   https://statisticsbyjim.com/basics/spurious-correlation/
4. Wikipedia — "Correlation does not imply causation."
   https://en.wikipedia.org/wiki/Correlation_does_not_imply_causation
5. Wikipedia — "Bradford Hill criteria."
   https://en.wikipedia.org/wiki/Bradford_Hill_criteria
6. "Methods in causal inference. Part 1: causal diagrams and confounding,"
   Evolutionary Human Sciences (Cambridge), PMC11588567.
   https://pmc.ncbi.nlm.nih.gov/articles/PMC11588567/
7. "Causal inference — so much more than statistics," PMC5841844.
   https://www.ncbi.nlm.nih.gov/pmc/articles/PMC5841844/
8. "Assessing causality in epidemiology: revisiting Bradford Hill," European
   Journal of Epidemiology, PMC8206235.
   https://www.ncbi.nlm.nih.gov/pmc/articles/PMC8206235/
9. IZA World of Labour — "Using instrumental variables to establish causality."
   https://wol.iza.org/articles/using-instrumental-variables-to-establish-causality/long
10. "Instrumental variables I," PMC2905668.
    https://pmc.ncbi.nlm.nih.gov/articles/PMC2905668/
11. Statistical Tools for Causal Inference, Ch. 4 "Natural Experiments."
    https://chabefer.github.io/STCI/NE.html
12. Tyler Vigen — "Spurious Correlations." https://www.tylervigen.com/spurious-correlations
