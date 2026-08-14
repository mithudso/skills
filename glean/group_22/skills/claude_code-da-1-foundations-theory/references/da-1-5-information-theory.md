<!-- hub-reference-banner -->
> **Reference file — part of the `da-1-foundations-theory` hub.** Formerly the standalone `da-1-5-information-theory` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-1-5-information-theory
description: >-
  Information-theoretic measures for data analysis foundations — Shannon entropy,
  joint/conditional entropy, mutual information, KL divergence (relative entropy),
  and cross-entropy. Covers definitions, formulas, units (bits/nats), core
  inequalities, and the practical pitfalls of estimating these quantities from
  finite samples (discretization bias, k-NN estimators, negative MI estimates).
  TRIGGER: questions about entropy/mutual information/KL divergence as a way to
  quantify uncertainty, dependence, or distribution difference in a dataset;
  "how much information does feature X carry about Y"; choosing or interpreting
  mutual_info_classif / mutual_info_regression; why an MI estimate is negative or
  unstable; cross-entropy vs KL divergence; bits vs nats. SKIP: decision-tree
  split criteria / information gain as an ML training mechanism (that is a model
  framing under predictive modeling, not this foundations node); cross-entropy
  *loss* as a neural-network objective (defer to ML/optimization skills);
  channel capacity, coding theory, and compression engineering (communications
  framing); correlation coefficients and linear association (defer to the
  statistics/association skills). Those framings get their own skills.
---

# Information Theory (entropy, mutual information, KL divergence)

Scope: information-theoretic quantities **as foundational data-analysis tools** —
measuring uncertainty in a variable, dependence between variables, and the gap
between two distributions. This is the "Foundations & Theory" framing. It is not
about communication channels, source coding, or using these quantities as a
training loss inside a specific model (those are separate taxonomy nodes; see
SKIP in the frontmatter).

## Why these measures exist

Variance and correlation describe *linear, second-moment* structure. Information
theory describes structure of **any** shape: entropy measures total uncertainty
in a distribution, and mutual information measures *any* statistical dependence
between two variables, linear or not. Entropy is the theoretical lower bound on
the average number of bits needed to encode symbols drawn from a source — rare
events carry more "surprise," hence more information [Toronto CSCC11 handout,
https://www.cs.toronto.edu/~fleet/courses/C11/Handouts/InformationTheory.pdf].

## Core definitions

All sums run over the support of the variable. Use log base 2 for **bits**,
natural log for **nats**. 1 nat ≈ 1.44 bits [machinelearningmastery.com,
https://machinelearningmastery.com/cross-entropy-for-machine-learning/].

### Shannon entropy

For a discrete random variable X with pmf p(x):

    H(X) = - Σ_x p(x) log p(x)

H(X) ≥ 0, is maximized by the uniform distribution (maximum uncertainty), and is
0 for a deterministic variable [Toronto handout]. Self-information of a single
outcome is −log p(x).

Worked example (base 2): a uniform distribution over 5 outcomes,
p = [0.2, 0.2, 0.2, 0.2, 0.2], gives
H(X) = −5·(0.2·log₂0.2) = 2.32 bits — the maximum for 5 categories
[eli.thegreenplace.net, https://eli.thegreenplace.net/2025/cross-entropy-and-kl-divergence/].

### Joint and conditional entropy

    H(X,Y) = - Σ_{x,y} p(x,y) log p(x,y)
    H(Y|X) = - Σ_{x,y} p(x,y) log p(y|x)

Chain rule: H(X,Y) = H(X) + H(Y|X). Conditioning never increases entropy:
H(Y|X) ≤ H(Y), with equality iff X and Y are independent [Toronto handout].

### Mutual information (MI)

The amount of information one variable carries about another — equivalently, how
much knowing X reduces uncertainty about Y:

    I(X;Y) = H(X) − H(X|Y) = H(Y) − H(Y|X) = H(X) + H(Y) − H(X,Y)

and, as a KL divergence between the joint and the product of marginals:

    I(X;Y) = D_KL( p(x,y) || p(x)p(y) )

[quantiki.org, https://www.quantiki.org/wiki/mutual-information; Toronto handout].

Key properties:
- I(X;Y) ≥ 0, and I(X;Y) = 0 **iff** X and Y are independent. This is the headline
  advantage over correlation: MI detects nonlinear dependence and is zero only
  under genuine independence.
- Symmetric: I(X;Y) = I(Y;X).
- Data-processing inequality: if X → Y → Z is a Markov chain, then
  I(X;Y) ≥ I(X;Z) — post-processing cannot create information about X
  [Toronto handout].

### KL divergence (relative entropy)

The expected extra cost (in bits/nats) of assuming distribution Q when the data
truly follow P:

    D_KL(P || Q) = Σ_x p(x) log( p(x) / q(x) )          (discrete)
    D_KL(P || Q) = ∫ p(x) log( p(x) / q(x) ) dx         (continuous)

Properties:
- D_KL(P||Q) ≥ 0 (Gibbs' inequality, a consequence of Jensen's inequality), with
  equality iff P = Q almost everywhere
  [tungmphung.com, https://tungmphung.com/information-theory-concepts-entropy-mutual-information-kl-divergence-and-more/].
- **Asymmetric**: D_KL(P||Q) ≠ D_KL(Q||P) in general. It is a *divergence*, not a
  distance/metric — no triangle inequality. Pick the direction deliberately
  (see "Pitfalls").
- Undefined / infinite when q(x) = 0 while p(x) > 0 (support mismatch).

### Cross-entropy and its relationship to KL

    H(P,Q) = - Σ_x p(x) log q(x)
    H(P,Q) = H(P) + D_KL(P || Q)

Because H(P) does not depend on Q, comparing models by cross-entropy and by KL
divergence to P gives the same ranking [eli.thegreenplace.net; machinelearningmastery.com].
(Cross-entropy *as a training loss* is the ML-objective framing — out of scope
here; this node is about cross-entropy as a distribution-comparison measure.)

## When to reach for which measure

- **Entropy H(X)** — quantify uncertainty/diversity of a single variable: class
  imbalance, the spread of a categorical column, "how predictable is this field."
- **Conditional entropy H(Y|X)** — remaining uncertainty about Y after observing X.
- **Mutual information I(X;Y)** — screen for *any* dependence between two
  variables during exploratory analysis or feature relevance ranking; detects
  nonlinear relationships that Pearson correlation misses.
- **KL divergence D_KL(P||Q)** — measure how far an observed distribution drifts
  from a reference (data-drift / population-stability monitoring), or how far an
  approximating distribution is from a target.
- **Cross-entropy** — compare candidate distributions against an empirical one.

## Estimation pitfalls (the part that bites in practice)

These measures are defined on *true* distributions. With finite samples you
estimate them, and the estimators have well-known failure modes
[scikit-learn mutual_info_classif docs,
https://scikit-learn.org/stable/modules/generated/sklearn.feature_selection.mutual_info_classif.html;
Bach MIFS, https://hoaibach.github.io/files/MIFS/MIFS.pdf].

1. **Discretization bias for continuous variables.** The entropy of a continuous
   variable differs from the entropy of its binned version; bin count and edges
   change the answer. Treating a continuous variable as discrete (or vice versa)
   usually gives wrong results. Prefer estimators built for continuous data
   rather than naive histogram binning [Bach MIFS].

2. **k-NN (Kraskov / Ross) estimators and the bias–variance knob.** scikit-learn's
   `mutual_info_classif` / `mutual_info_regression` estimate MI nonparametrically
   from k-nearest-neighbour distances. Larger `n_neighbors` lowers variance but
   raises bias; smaller `n_neighbors` does the reverse. Set `discrete_features`
   correctly per column, and note results are stochastic (set `random_state` for
   reproducibility) [scikit-learn docs].

3. **Negative MI estimates.** True MI cannot be negative; a negative *estimate* is
   an artifact and is typically clamped to 0. Do not interpret a small negative
   number as "negative dependence" [scikit-learn docs].

4. **KL support mismatch.** D_KL blows up to infinity when the reference Q assigns
   zero probability to an outcome that P gives positive probability. Smooth or
   add pseudocounts to Q, or use a symmetric alternative (Jensen–Shannon
   divergence) for drift comparisons.

5. **KL direction matters.** Forward KL D_KL(P||Q) is "mean-seeking" (Q must cover
   all of P's mass, penalizes Q=0 where P>0); reverse KL D_KL(Q||P) is
   "mode-seeking." Choose based on whether you care about coverage or
   concentration [tungmphung.com].

6. **Small-sample / high-dimensionality bias.** Plug-in (maximum-likelihood)
   entropy estimates are biased downward when the number of distinct values is
   large relative to sample size; bias corrections exist. Reserve enough data per
   cell before trusting an entropy or MI number.

## Quick reference

| Quantity | Formula | Range | Symmetric? | Zero when |
|---|---|---|---|---|
| Entropy H(X) | −Σ p log p | [0, log\|X\|] | n/a | X deterministic |
| Conditional H(Y\|X) | −Σ p(x,y) log p(y\|x) | [0, H(Y)] | no | Y determined by X |
| Mutual info I(X;Y) | H(X)+H(Y)−H(X,Y) | [0, min(H(X),H(Y))] | yes | X ⟂ Y |
| KL D_KL(P\|\|Q) | Σ p log(p/q) | [0, ∞) | no | P = Q |
| Cross-entropy H(P,Q) | −Σ p log q | [H(P), ∞) | no | P=Q (min = H(P)) |

## Sources

1. University of Toronto, CSCC11 Information Theory handout (entropy, joint/conditional
   entropy, mutual information, chain rule, data-processing inequality) —
   https://www.cs.toronto.edu/~fleet/courses/C11/Handouts/InformationTheory.pdf
2. Eli Bendersky, "Cross-entropy and KL divergence" (definitions, worked numeric
   examples, H(P,Q) = H(P) + D_KL) —
   https://eli.thegreenplace.net/2025/cross-entropy-and-kl-divergence/
3. Tung M. Phung, "Information Theory concepts: Entropy, Mutual Information,
   KL-Divergence" (KL non-negativity, asymmetry, direction) —
   https://tungmphung.com/information-theory-concepts-entropy-mutual-information-kl-divergence-and-more/
4. scikit-learn, mutual_info_classif documentation (k-NN MI estimation,
   discrete_features, n_neighbors bias/variance, non-negativity clamping) —
   https://scikit-learn.org/stable/modules/generated/sklearn.feature_selection.mutual_info_classif.html
5. Bach et al., "Mutual Information for Feature Selection: Estimation" (discretization
   bias, continuous vs discrete estimation) —
   https://hoaibach.github.io/files/MIFS/MIFS.pdf
6. Quantiki, "Mutual information" (I(X;Y) = D_KL(p(x,y)||p(x)p(y))) —
   https://www.quantiki.org/wiki/mutual-information
7. Jason Brownlee, "A Gentle Introduction to Cross-Entropy for Machine Learning"
   (bits vs nats, cross-entropy/KL relationship) —
   https://machinelearningmastery.com/cross-entropy-for-machine-learning/
