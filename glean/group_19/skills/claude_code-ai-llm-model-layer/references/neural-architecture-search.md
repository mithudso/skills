<!-- hub-reference-banner -->
> **Reference file — part of the `ai-llm-model-layer` hub.** New `/dr`-researched reference (concept-family-explorer frontier run, 2026-06-15).
> Sibling topics in this family are reference files under the hubs — **not** standalone skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling skill; load that topic's `references/<name>.md` from the owning hub.

---

---
name: neural-architecture-search
title: Neural Architecture Search (NAS)
description: >
  Expert reference for automated search over neural-network ARCHITECTURES, organized by the three-axis framing (Elsken et al. 2019): search space, search strategy, performance-estimation strategy. Covers search spaces (chain, cell/NASNet micro vs macro, hierarchical motifs); search strategies (RL controller + REINFORCE — Zoph & Le; aging/regularized evolution — AmoebaNet; gradient-based continuous relaxation + bilevel optimization — DARTS; Bayesian opt); performance-estimation acceleration (weight-sharing one-shot SuperNet — ENAS; predictors; zero-cost proxies — synflow/jacov; lower-fidelity/early-stop); hardware-aware NAS (MnasNet latency reward, FBNet, ProxylessNAS, Once-for-All progressive shrinking); GPU-day efficiency landmarks; reproducibility benchmarks (NAS-Bench-101/201/301, random-search baselines). TRIGGER: choosing/comparing a NAS search strategy; cell vs macro search-space design; weight-sharing/one-shot supernets; zero-cost proxies or performance predictors; hardware/latency-aware NAS; NAS benchmarks/reproducibility. SKIP: transformer internals/training/scaling/alignment -> ai-llm-model-layer; classical HPO/non-LLM AutoML -> da-analytical-methods; algorithm & program discovery / superoptimization (AlphaTensor/AlphaDev/FunSearch/AlphaEvolve) and agent/compound-system search -> algorithmic-discovery-superoptimization (ai-agents-orchestration).
category: developer
---

# Neural Architecture Search (NAS)

NAS automates the design of neural-network architectures: instead of a human hand-tuning depth, width, connectivity, and operations, a search procedure proposes candidate architectures, a (cheap) evaluator scores them, and a search strategy steers toward better ones. It is the **architecture-search specialization of AutoML** and a direct **sibling of algorithmic/program discovery** under one shared meta-pattern — *a proposer generates candidates that a verifier/evaluator scores* (NAS searches architectures; FunSearch/AlphaEvolve search programs; ADAS/AFlow search agent systems → see the SKIP edges).

**The one mental model that unlocks the whole field (Elsken, Metzen & Hutter, JMLR 2019, "Neural Architecture Search: A Survey", https://jmlr.org/papers/v20/18-598.html).** Every NAS method is a choice along **three orthogonal axes**:

1. **Search space** — *which* architectures are even representable (the inductive bias / human prior).
2. **Search strategy** — *how* you explore that space (the classic exploration–exploitation tradeoff): RL, evolution, gradient descent, Bayesian optimization, or random search.
3. **Performance-estimation strategy** — *how cheaply* you estimate a candidate's quality (the dominant cost lever): full training → lower-fidelity proxies → weight sharing → zero-cost-at-init.

Hold those three axes and every paper below slots into one cell of the grid. The historical arc is a relentless push on **axis 3**: the original RL/evolution methods cost thousands of GPU-days because they trained every candidate to convergence; weight sharing and zero-cost proxies cut that by 3–4 orders of magnitude.

> Scope guard: this reference is **which architecture to search for and how**. It defers *training/serving/scaling* the winning net to `ai-llm-model-layer`, *classical HPO* to `da-analytical-methods`, and *program/algorithm/agent discovery* to `algorithmic-discovery-superoptimization`.

---

## Core Concepts

### 1. The three-axis framing (search space / search strategy / performance estimation)
The organizing taxonomy of the field (Elsken et al. 2019). Decoupling these axes lets you reason about a method's *cost* (axis 3) separately from its *optimizer* (axis 2) and its *expressiveness/bias* (axis 1). It also exposes the central tension: **incorporating prior knowledge shrinks the search space and simplifies search, but introduces human bias that may preclude discovering genuinely novel building blocks.** Published as Ch. 3 of the AutoML book, framing NAS as a sub-area of AutoML distinct from classical HPO/CASH.

### 2. Search spaces: chain → cell/micro vs macro → hierarchical
- **Chain-structured**: a sequence of layers (depth + per-layer op/hyperparams). Simple, low-dimensional, limited.
- **Multi-branch macro-architectures**: skip connections, branching (ResNet/DenseNet-like) at the whole-network level.
- **Cell/block-based (micro) search** — the workhorse since NASNet: search a small repeatable **cell** (a DAG of ops) on a cheap proxy dataset, then **stack copies** of it into the full network. Reduces dimensionality and gives **transferability** (search on CIFAR-10, deploy on ImageNet), at the cost of a **hand-fixed macro-architecture** (a human bias).
- **Hierarchical search space** (Liu et al. 2018): levels of motifs — primitives → motifs that wire primitives via a DAG → motifs-of-motifs. The cell-based space is the special **3-level** case with a hard-coded top level.

### 3. NASNet search space: Normal Cell + Reduction Cell, transfer by stacking
Zoph et al., CVPR 2018, "Learning Transferable Architectures for Scalable Image Recognition" (https://arxiv.org/abs/1707.07012) — the NASNet paper, distinct from the 2017 founding RL paper (§4). Defines two cell types: a **Normal Cell** (preserves spatial resolution) and a **Reduction Cell** (stride-2, halves resolution). Each cell consumes two prior hidden states; the controller predicts **B blocks of 5 choices** each. **ScheduledDropPath** regularization is key to generalization. NASNet hit 2.4% CIFAR-10 error and 82.7% ImageNet top-1 with ~28% fewer FLOPs than the best human designs — the proof that searched architectures can beat hand-design at scale. Canonical "expensive era" cost: **~1800–2000 GPU-days**.

### 4. RL-based search: RNN controller + REINFORCE over a config string
Zoph & Le, ICLR 2017, "Neural Architecture Search with Reinforcement Learning" (https://arxiv.org/abs/1611.01578). The founding NAS paper: an **RNN "controller"** emits a variable-length string describing an architecture (a flat macro space); the **child network** is trained to convergence and its validation accuracy is a **non-differentiable reward**; the controller is updated by the **REINFORCE policy gradient** (Williams 1992) to maximize expected accuracy. Landmark compute: **800 GPUs**, ~12,800 child models for CIFAR-10 — the original "thousands of GPU-days" baseline. Established the controller/child template and config-string encoding reused everywhere after. (Cells/ScheduledDropPath/transfer are the *2018* NASNet follow-up in §3, not this paper.)

### 5. Evolutionary search: aging / regularized evolution (AmoebaNet)
Real et al., AAAI 2019, "Regularized Evolution for Image Classifier Architecture Search" (https://arxiv.org/abs/1802.01548). Modifies **tournament selection** with **aging**: genotypes are discarded by *age* (oldest removed), favoring younger candidates and preventing lucky early models from dominating the population. **AmoebaNet-A** was the first evolved classifier to surpass hand-designs (83.9% ImageNet top-1). In a controlled head-to-head with RL on the *same* NASNet space and hardware, **evolution reaches good models faster early** — relevant under limited compute. Most-expensive landmark: **~3150 GPU-days**. Cements RL / evolution / gradient as directly comparable strategy families on a shared cell space.

### 6. Gradient-based search: continuous relaxation + bilevel optimization (DARTS)
Liu, Simonyan & Yang, ICLR 2019, "DARTS: Differentiable Architecture Search" (https://arxiv.org/abs/1806.09055). Reframes search as **continuous optimization**: relax the categorical choice of operation on each DAG edge into a **softmax over candidate ops** weighted by architecture parameters **α**, making the architecture differentiable. **Bilevel optimization** alternately updates network weights *w* (on train) and α (on validation); a **first-order approximation** drops the costly second-order term for speed. After search, **discretize** by argmax per edge — introducing the **discretization gap** between supernet and derived child. Drove cost from thousands of GPU-days to **~1–4 GPU-days** with no controller/predictor. Spawned PC-DARTS (~0.1 GPU-day), P-DARTS, GDAS, FairDARTS.

### 7. Weight sharing / one-shot SuperNet (ENAS)
Pham et al., ICML 2018, "Efficient Neural Architecture Search via Parameter Sharing" (https://arxiv.org/abs/1802.03268). The breakthrough that killed NAS's compute bottleneck: **all candidate architectures are subgraphs of one large computational supergraph and share a single set of weights**, so you train *one* supernet instead of training each child from scratch. An RL controller samples subgraphs; supernet weights are trained on them. **>1000× cheaper** than the original NAS (<16 hours / ~0.5 GPU-days on one GTX 1080Ti), reaching 2.89% CIFAR-10 error. This **one-shot model / SuperNet** idea is the foundation under DARTS, ProxylessNAS, FBNet, and Once-for-All.

### 8. Performance-estimation acceleration spectrum
The axis-3 cost lever, ordered cheapest-to-fund last: **lower-fidelity / early-stopping** (fewer epochs, smaller images, subset data) → **learning-curve extrapolation** → **weight inheritance / network morphisms** → **performance predictors** (regress final accuracy from an architecture encoding, e.g. GNN/graph encoders, BANANAS) → **weight sharing / one-shot** (§7) → **zero-cost proxies** (§9). Each step trades estimation fidelity for speed; the field's progress is essentially walking down this list.

### 9. Zero-cost proxies (training-free NAS)
Abdelfattah et al., ICLR 2021, "Zero-Cost Proxies for Lightweight NAS" (https://openreview.net/forum?id=0cmMMy8J5q). Adapts **pruning-at-initialization** saliency metrics — **snip, grasp, synflow**, plus **grad_norm** and **jacob_cov** — into network-level scores computed from a **single minibatch at initialization, no training**. **Synflow** needs *no data at all* (loss = product of all parameters) and best preserves ranking: Spearman ~0.82 with final accuracy on NAS-Bench-201 (vs 0.61 for the reduced-training EcoNAS proxy). Uses **~3 orders of magnitude less compute** than proxy tasks — the extreme end of axis 3. Plugged into random/RL/evolution/predictor search to improve sample efficiency; the standard baseline any new estimation method must beat. Open question: do proxies generalize *beyond* NAS-Bench spaces?

### 10. Hardware-aware NAS: latency in the objective (MnasNet)
Tan et al., CVPR 2019, "MnasNet: Platform-Aware NAS for Mobile" (https://arxiv.org/abs/1807.11626). The founding hardware-aware paper: a **multi-objective** search optimizing **accuracy *and* measured real-device inference latency** via a weighted/Pareto reward. Key insight: **FLOPs is an inaccurate latency proxy** (MobileNet vs NASNet have similar FLOPs, very different latency), so MnasNet **runs candidates on real phones (Pixel)**. Introduced a **factorized hierarchical search space** allowing per-block diversity. Hit 75.2% ImageNet top-1 at 78ms — 1.8× faster than MobileNetV2. Set the template: direct latency objective, real-device measurement, device-specific architectures.

### 11. Differentiable hardware-aware NAS (ProxylessNAS, FBNet)
Cai et al., ICLR 2019, "ProxylessNAS" (https://arxiv.org/abs/1812.00332); Wu et al., CVPR 2019, "FBNet" (https://arxiv.org/abs/1812.03443). **ProxylessNAS** searches **directly on the target task and target hardware** (no proxy dataset), **binarizes supernet paths** to cut one-shot training memory, and adds a **differentiable latency-regularization loss** — ~200× fewer GPU-hours than MnasNet (74.6% top-1 at 78ms mobile). **FBNet (DNAS)** uses **Gumbel-softmax** differentiable search over a layer-wise supernet with a **latency lookup-table** term, making latency differentiable *without* per-candidate device runs (FBNet-B: 74.1% top-1 at 23.1ms, ~420× lower search cost than MnasNet). Recipe: DARTS-style relaxation + latency penalty = the dominant practical-mobile NAS approach.

### 12. Once-for-All: decouple training from search (progressive shrinking)
Cai et al., ICLR 2020, "Once-for-All: Train One Network and Specialize It for Efficient Deployment" (https://arxiv.org/abs/1908.09791). Train a **single supernet once**, then extract **~10^19 specialized sub-networks** for different hardware/latency targets **with no retraining**. **Progressive Shrinking** trains the largest net first (max depth/width/kernel/resolution) then fine-tunes down to support smaller weight-sharing sub-nets without small/large interference — a generalized multi-dimension pruning. Sub-net selection uses **accuracy predictors + latency tables**, so "search" per device is near-instant. Reaches 80.0% ImageNet top-1 under the mobile <600M-MAC setting, beating MobileNetV3/EfficientNet on measured latency. Explicitly motivated by the **carbon/compute cost** of per-device NAS. The supernet itself becomes the deployable artifact.

### 13. Efficiency landmarks (GPU-days) and the cost collapse
The canonical timeline: **NASNet ~1800–2000** (RL) and **AmoebaNet ~3150** (evolution) → **DARTS ~1–4**, **ENAS ~0.5**, **PC-DARTS ~0.1** → **zero-cost proxies ~3 orders of magnitude cheaper still**. Memorize this arc: it *is* the story of NAS axis-3 progress and the basis for "is this method affordable?" reasoning. Green-NAS / CO₂ accounting now treats search cost as a first-class objective.

### 14. Reproducibility crisis and NAS benchmarks (101 / 201 / 301)
- **NAS-Bench-101** (Ying et al., ICML 2019, https://arxiv.org/abs/1902.09635): first **tabular** benchmark — a compact cell space (≤7 edges/9 vertices; ops 3×3 conv, 1×1 conv, 3×3 maxpool) = **423,624 unique architectures**, all exhaustively trained (3× at 4 budgets) → ~5M models, queryable in milliseconds. Motivated by the reproducibility crisis: methods used different pipelines/spaces, so gains couldn't be attributed to the *search algorithm* vs the *protocol*.
- **NAS-Bench-201 / NATS-Bench** (Dong & Yang, ICLR 2020, https://arxiv.org/abs/2001.00326): fixed cell (4 nodes, 5 ops/edge) = **15,625 architectures** trained on **3 datasets** (CIFAR-10/100, ImageNet-16-120), applicable to almost any algorithm including weight-sharing.
- **NAS-Bench-301** (Siems, Zela et al. 2020, https://arxiv.org/abs/2008.09777): the first **surrogate** benchmark — meta-learns a regression surrogate (deep ensembles for uncertainty) over the full **DARTS cell space (~10^18 architectures)** from ~60k evaluations, escaping the tabular cap that exhaustive training imposes.

### 15. NAS as specialized HPO; random search with weight sharing as a strong baseline
Li & Talwalkar's reproducibility critique (UAI 2019, https://arxiv.org/abs/1902.07638): **NAS is a specialized hyperparameter-optimization problem**, and **random search with weight sharing + early stopping is a strong baseline** that matches ENAS/DARTS — many published NAS gains were not robustly reproducible. Pair this with the NAS best-practice checklist (Lindauer & Hutter, JMLR 2020, https://arxiv.org/abs/1909.02453) and **always report a random-search baseline**.

### 16. The verifier/evaluator-gated search meta-pattern
NAS is one instance of a broader pattern: a **search strategy proposes candidates that a (cheap) evaluator scores**. NAS searches *architectures*; algorithm/program discovery (AlphaTensor/AlphaDev/FunSearch/AlphaEvolve) searches *programs*; automated agent design (ADAS/AFlow/DSPy) searches *compound LLM systems*. Same scaffolding (proposer + verifier + search loop), different artifact — which is exactly where the SKIP edges route.

---

## References outline (suggested `references/` files)

- `references/three-axis-framing.md` — Elsken et al. survey; search-space taxonomy (chain/cell/hierarchical); prior-knowledge-vs-novelty tension.
- `references/search-strategies.md` — RL+REINFORCE (Zoph & Le), aging evolution (AmoebaNet), gradient/DARTS bilevel + discretization gap, Bayesian opt, random-search baseline.
- `references/weight-sharing-one-shot.md` — ENAS supernet/subgraph; supernet training; the one-shot lineage to DARTS/Proxyless/FBNet/OFA.
- `references/performance-estimation.md` — the acceleration spectrum; predictors (BANANAS/GNN); zero-cost proxies (synflow/jacob_cov/snip/grasp/grad_norm) and Spearman ranking; generalization limits.
- `references/hardware-aware-nas.md` — MnasNet latency reward + real-device measurement; ProxylessNAS path binarization; FBNet Gumbel-softmax + latency LUT; OFA progressive shrinking; Pareto/multi-objective.
- `references/darts-robustness.md` — skip-connection collapse, discretization gap, R-DARTS / FairDARTS / DARTS- / perturbation-based pruning.
- `references/benchmarks-reproducibility.md` — NAS-Bench-101/201/301, tabular vs surrogate, Li & Talwalkar critique, Lindauer & Hutter best-practice checklist, green-NAS / CO₂.
- `references/nas-for-transformers.md` — searching attention variants, MoE routing, hybrid blocks; defer training mechanics to `ai-llm-model-layer`.

## Authoritative sources (inline-cited above)

- Elsken, Metzen, Hutter — *Neural Architecture Search: A Survey* (JMLR 2019) — the three-axis framing.
- Zoph & Le — *NAS with Reinforcement Learning* (ICLR 2017, arXiv:1611.01578) — RL controller + REINFORCE (flat macro space).
- Zoph et al. — *NASNet / Learning Transferable Architectures* (CVPR 2018, arXiv:1707.07012) — cell-based transferable search space, ScheduledDropPath, ~1800–2000 GPU-days.
- Real et al. — *Regularized Evolution / AmoebaNet* (AAAI 2019, arXiv:1802.01548) — aging evolution, ~3150 GPU-days.
- Liu, Simonyan, Yang — *DARTS* (ICLR 2019, arXiv:1806.09055) — differentiable continuous relaxation.
- Pham et al. — *ENAS* (ICML 2018, arXiv:1802.03268) — weight sharing / one-shot.
- Tan et al. — *MnasNet* (CVPR 2019, arXiv:1807.11626); Cai et al. — *ProxylessNAS* (ICLR 2019, arXiv:1812.00332) & *Once-for-All* (ICLR 2020, arXiv:1908.09791); Wu et al. — *FBNet* (CVPR 2019, arXiv:1812.03443) — hardware-aware NAS.
- Abdelfattah et al. — *Zero-Cost Proxies* (ICLR 2021, arXiv:2101.08134).
- Ying et al. — *NAS-Bench-101* (ICML 2019, arXiv:1902.09635); Dong & Yang — *NAS-Bench-201* (ICLR 2020, arXiv:2001.00326); Siems, Zela et al. — *NAS-Bench-301* surrogate (arXiv:2008.09777).
- Li & Talwalkar — *Random Search and Reproducibility in NAS* (UAI 2019, arXiv:1902.07638); Lindauer & Hutter — *Best Practices for Scientific Research on NAS* (JMLR 2020, arXiv:1909.02453).
