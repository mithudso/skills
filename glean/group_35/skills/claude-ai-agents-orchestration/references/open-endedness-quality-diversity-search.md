<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Authored as a hub spoke (not a standalone skill).
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

<!-- Provenance: reference under the `ai-agents-orchestration` hub. Created 2026-06-15 via /dr deep-research (exa neural + WebSearch, arXiv-direct verification) from primary sources — MAP-Elites (Mouret & Clune, arXiv:1504.04909, 2015), Novelty Search (Lehman & Stanley, "Abandoning Objectives", Evolutionary Computation 2011; GECCO-2010), Novelty Search with Local Competition (Lehman & Stanley, GECCO 2011), CMA-ME (Fontaine et al., GECCO 2020, 10.1145/3377930.3390232), Differentiable QD / CMA-MEGA (Fontaine & Nikolaidis, NeurIPS 2021), CMA-MAE/CMA-MAEGA (Fontaine & Nikolaidis 2023), POET (Wang, Lehman, Clune, Stanley, arXiv:1901.01753, GECCO 2019), Enhanced POET (Wang et al., ICML 2020), FunSearch (Romera-Paredes et al., Nature 625:468-475, 2023, 10.1038/s41586-023-06924-6), AlphaEvolve (Novikov et al., arXiv:2506.13131, 2025), "Why Greatness Cannot Be Planned" (Stanley & Lehman, Springer 2015), QD survey/framework (Pugh et al., Frontiers in Robotics & AI 2016; Cully & Demiris, IEEE TEVC 2018, "unifying modular framework"), pyribs/RIBS (Tjanaka et al., GECCO 2023). Scope: the EVOLUTIONARY / ILLUMINATION SEARCH-ALGORITHM SUBSTRATE — population-based divergent search that produces a diverse ARCHIVE of high-performing solutions (not one optimum), and the open-endedness theory beneath it. This is the ENGINE; `automated-agent-system-design` is the APPLICATION that runs this engine over agent/program design spaces. NOT gradient RL on model weights (→ agentic-rl); NOT test-time compute (→ reasoning-models); NOT classical single-objective Bayesian optimization / hyperparameter tuning (→ da-7-machine-learning). -->

---
name: open-endedness-quality-diversity-search
title: Open-Endedness & Quality-Diversity Search
description: >
  The EVOLUTIONARY / ILLUMINATION search-algorithm SUBSTRATE that powers automated
  agent & compound-system design — population-based divergent search that returns a
  diverse ARCHIVE of high-performing solutions, not one optimum. Covers
  Quality-Diversity (QD) algorithms: MAP-Elites (+ CVT-MAP-Elites, CMA-ME, CMA-MEGA,
  CMA-MAE) and Novelty Search / Novelty Search with Local Competition (NSLC); the
  illumination/archive structure (behavior descriptors / feature space, the elites
  grid, the novelty archive, selection + variation operators, the curiosity/
  improvement scores); open-endedness theory (Stanley & Lehman, "Why Greatness
  Cannot Be Planned" — the objective paradox, deception, divergent vs objective
  search, stepping stones) and POET / Enhanced POET (co-evolving environments +
  agents, minimal-criterion coevolution, transfer); LLM-driven evolutionary program
  search — FunSearch (LLM proposer + automated evaluator + programs database) and
  AlphaEvolve (Gemini evolutionary coding agent) — and HOW ADAS / AFlow / GEPA use
  QD/evolutionary search as their engine; plus objective/fitness + diversity-metric
  design, the QD-score, deception & stepping stones, and benchmarks/eval.
  TRIGGER: quality-diversity / illumination algorithms; MAP-Elites, CMA-ME/CMA-MEGA,
  Novelty Search, NSLC; behavior descriptors / feature space / elites archive; open-
  endedness, divergent search, the myth of the objective, deception, stepping
  stones; POET / Enhanced POET / co-evolving curricula; FunSearch / AlphaEvolve /
  LLM-driven evolutionary program search; designing a diversity metric or QD-score;
  "search for a diverse set of good solutions, not one optimum"; the search SUBSTRATE
  beneath ADAS/AFlow/GEPA.
  SKIP: searching/optimizing the agent or LLM-program design itself (ADAS, AFlow,
  DSPy/MIPROv2, GEPA, TextGrad, Trace/OptoPrime, Darwin-Gödel) → automated-agent-
  system-design (that is the APPLICATION; this is the ENGINE); using FunSearch/
  AlphaEvolve (or STOKE/Souper/AlphaTensor/AlphaDev) to discover a faster CONCRETE
  algorithm or program — matmul, sort, heuristic, superoptimized instruction
  sequence, verifier-gated program synthesis → algorithmic-discovery-superoptimization
  (this spoke owns their QD/illumination MECHANICS, that one the algorithm-discovery
  APPLICATION); gradient RL on model WEIGHTS (GRPO/PPO/DPO) → agentic-rl; test-time/
  inference compute scaling → reasoning-models; classical single-objective Bayesian
  optimization / hyperparameter tuning → da-7-machine-learning.
origin: local
category: developer
version: "1.0.1"
updated: "2026-06-15"
tags:
  - quality-diversity
  - map-elites
  - novelty-search
  - nslc
  - cma-me
  - cma-mega
  - open-endedness
  - poet
  - funsearch
  - alphaevolve
  - behavior-descriptor
  - illumination-algorithm
  - stepping-stones
  - divergent-search
when_to_use:
  - "quality-diversity / illumination algorithms (return an archive of diverse good solutions, not one optimum)"
  - "MAP-Elites and its variants (CVT-MAP-Elites, CMA-ME, CMA-MEGA, CMA-MAE)"
  - "Novelty Search and Novelty Search with Local Competition (NSLC)"
  - "designing behavior descriptors / a feature (measure) space / an elites archive"
  - "open-endedness theory — divergent vs objective search, deception, stepping stones, the myth of the objective"
  - "POET / Enhanced POET — co-evolving environments and agents, minimal-criterion coevolution, transfer"
  - "FunSearch / AlphaEvolve — LLM-driven evolutionary program search"
  - "understanding the search ENGINE that ADAS / AFlow / GEPA run on top of"
  - "designing an objective/fitness function, a diversity metric, or a QD-score; deception & stepping-stone benchmarks"
when_not_to_use:
  - "searching/optimizing the agent or multi-step LLM PROGRAM itself (ADAS, AFlow, DSPy/MIPROv2, GEPA, TextGrad, Trace, Darwin-Gödel) → automated-agent-system-design"
  - "using FunSearch/AlphaEvolve (or STOKE/Souper/AlphaTensor/AlphaDev) to discover a faster CONCRETE algorithm/program (matmul, sort, heuristic, instruction sequence) → algorithmic-discovery-superoptimization (this owns their QD/illumination mechanics; that owns the algorithm-discovery application)"
  - "gradient RL on model WEIGHTS (GRPO/PPO/DPO) → agentic-rl"
  - "test-time / inference compute scaling → reasoning-models"
  - "classical single-objective Bayesian optimization / hyperparameter tuning → da-7-machine-learning"
related_skills:
  - automated-agent-system-design
  - algorithmic-discovery-superoptimization
  - agentic-rl
  - reasoning-models
  - eval-driven-development
---

# Open-Endedness & Quality-Diversity Search

`verified-as-of: 2026-06-15`

The **search-algorithm substrate** beneath automated agent and compound-AI-system
design. Where most optimization asks *"what is the single best solution?"*, this
family asks two different questions:

- **Quality-Diversity (QD):** *"what is the best solution for **each kind** of
  solution?"* — return a whole **archive** of solutions that are individually
  high-performing yet **behaviorally diverse**.
- **Open-endedness:** *"how do we build a process that **keeps inventing** new,
  more complex, interesting challenges and solutions — forever, without a fixed
  target?"*

These are the engines that LLM-driven program search (FunSearch, AlphaEvolve) and
automated agent design (ADAS, AFlow, GEPA) actually run on top of. **This spoke is
the ENGINE; [`automated-agent-system-design`](./automated-agent-system-design.md)
is the APPLICATION that points that engine at agent/program design spaces.** When
the question is *how the search works* (archive, behavior descriptors, novelty,
stepping stones, the objective paradox), you are here. When it is *which agent/
program the optimizer should output from a metric*, defer there.

---

## 1. Why diverge? The objective paradox (Stanley & Lehman)

The foundational thesis — *Why Greatness Cannot Be Planned: The Myth of the
Objective* (Stanley & Lehman, 2015), growing out of the **Novelty Search** result
(Lehman & Stanley, 2011) — is that **for ambitious problems, optimizing the
objective directly is often the *worst* way to reach it.**

- **Deception.** A fitness/objective gradient can point *away* from the global
  optimum. The canonical demonstration: a maze where minimizing straight-line
  distance to the goal drives the agent into a wall and traps it; the solution
  requires temporarily moving *away* from the objective.
- **Stepping stones.** Great achievements are reached through a chain of stepping
  stones that **do not resemble the final product** (vacuum tubes → computers).
  You cannot know in advance which stepping stones lead onward, so a distant
  objective is a poor compass.
- **Objectives help only when the goal is ~one hop away.** For multi-hop /
  "great" goals, **pursue novelty or interestingness** instead of the objective —
  collect stepping stones, then exploit them. This is the philosophical root of
  every algorithm below.

**Operator takeaway:** when an automated-design loop *stalls or collapses to one
mode*, the cause is usually a **deceptive objective** with no diversity pressure.
The fix is structural (add a behavior-diversity axis or a novelty term), not "tune
the reward harder."

---

## 2. Novelty Search and NSLC

- **Novelty Search (NS)** (Lehman & Stanley, 2011, *"Abandoning Objectives:
  Evolution Through the Search for Novelty Alone"*): **drop the objective entirely.**
  Reward each individual purely by its **novelty** — its distance in *behavior
  space* to its *k* nearest neighbors among the current population **plus a
  persistent novelty archive** of past novel individuals. Surprisingly, on
  deceptive tasks pure NS often reaches the goal faster than objective-driven search.
- **Behavior characterization (BC) / behavior descriptor (BD):** the low-dimensional
  vector that summarizes *what a solution does* (e.g., a robot's final (x,y), gait
  duty factor, body orientation). **Choosing the BD is the single most consequential
  design decision** — it defines the space in which "novel" and "diverse" are
  measured. Getting it wrong is the #1 QD failure mode.
- **Novelty Search with Local Competition (NSLC)** (Lehman & Stanley, GECCO 2011):
  adds a **quality** signal *back in*, but **only locally** — an individual is
  rewarded for outperforming the individuals **most similar to it in behavior**
  (its niche), not the whole population. This is the **direct ancestor of QD**:
  diverge globally, compete locally, so you accumulate *many* good-and-different
  solutions in one run instead of converging to one morphology. First applied to
  evolving a diversity of locomoting virtual creatures.

---

## 3. MAP-Elites and the illumination archive

**MAP-Elites** — *Multi-dimensional Archive of Phenotypic Elites* (Mouret & Clune,
arXiv:1504.04909, 2015) — is the workhorse QD algorithm: simple, effective, easy to
implement.

**Structure (the canonical QD container):**

1. **Feature / measure space** = the BD axes the *user chooses* (e.g., molecule
   size × cost; robot height × weight). Discretized into a **grid of cells (niches)**.
2. **Elites archive** = each cell holds **at most one** solution: the
   highest-fitness solution discovered *with that behavior*. The population **is**
   the archive.
3. **Loop:** randomly select an occupied cell → **vary** (mutate/crossover) its
   elite → evaluate the child to get its (fitness, BD) → place it in the cell its
   BD maps to **iff** it beats the incumbent (or the cell is empty). Repeat.

**Why it's called an *illumination* algorithm:** it "lights up" the **fitness
potential of every region** of the feature space, exposing trade-offs between the
features of interest and performance — not just the single peak. A well-known
side-effect: because it explores far more of the space, **MAP-Elites can find a
better global optimum than a dedicated single-objective optimizer** — Mouret & Clune
(2015) report this on a deceptive domain (the always-true gain is coverage/diversity;
the better-global-optimum result is domain-dependent, strongest under deception, not
universal). It also yields a **repertoire** (e.g., a damaged robot falls back to a
different gait in the archive — the basis of *Robots that can adapt like animals*,
Cully et al., Nature 521:503-507, 2015).

**Key variants:**

- **CVT-MAP-Elites** (Vassiliades, Chatzilygeroudis, Mouret, IEEE TEVC 22(4):623-630,
  2018): use a **Centroidal Voronoi Tessellation** to partition the BD space into a
  fixed number of well-distributed cells, so MAP-Elites **scales to high-dimensional**
  behavior spaces where a regular grid would explode.
- **CMA-ME** (Fontaine et al., GECCO 2020): replace random mutation with the
  self-adaptation of **CMA-ES**. Emitters run covariance-matrix-adaptation searches
  that *steer toward filling and improving the archive*, **more than doubling**
  MAP-Elites' QD-score on continuous domains.
- **Differentiable QD (DQD) / CMA-MEGA** (Fontaine & Nikolaidis, NeurIPS 2021): when
  the objective and the measure functions are **differentiable**, exploit their
  **gradients** via a *gradient arborescence* (branch in measure space, ascend the
  QD objective). Far more sample-efficient; demonstrated illuminating a StyleGAN
  latent space.
- **CMA-MAE** (Covariance Matrix Adaptation MAP-**Annealing**; Fontaine & Nikolaidis,
  GECCO 2023): smooths the hard "elite replaces incumbent" rule with a soft, annealed
  acceptance threshold per cell, fixing CMA-ME's brittleness and unifying the
  optimization↔illumination spectrum. **CMA-MAEGA** is its separate *gradient* (DQD)
  variant (CMA-MAE + the CMA-MEGA gradient arborescence) — not a synonym for CMA-MAE.
- **Reference implementation:** **pyribs** (icaros-usc, GECCO 2023) — the RIBS
  framework factors every QD algorithm into **archive + emitter(s) + scheduler**;
  the cleanest mental model for "what are the moving parts."

---

## 4. Open-endedness in practice: POET

**POET** — *Paired Open-Ended Trailblazer* (Wang, Lehman, Clune, Stanley,
arXiv:1901.01753, GECCO 2019) — confronts open-endedness **directly**: instead of a
fixed task, it **co-evolves the environments and their solutions together.**

- Maintains a population of **(environment, agent) pairs** (e.g., 2-D BipedalWalker
  obstacle courses, each paired with a controller network).
- **Generate** new environments by mutating existing ones; keep a new environment
  only if it passes a **minimal criterion** — *not too hard, not too easy* for the
  current agents (inherited from **Minimal Criterion Coevolution**, MCC).
- **Optimize** each agent on its paired environment (ES), **and** (crucially)
  **attempt transfers**: periodically test every agent on every other environment.
  A solution to one environment is frequently the **stepping stone** that unlocks
  progress on another.
- Result: a single run radiates an **endless, branching curriculum** of increasingly
  complex challenges and skills — many of which **cannot be learned by direct
  optimization** or a hand-built straight-line curriculum.
- **Enhanced POET** (Wang et al., ICML 2020) adds a domain-general novelty measure
  for environments (so the environment generator keeps producing *meaningfully* new
  challenges) and improved transfer, removing hand-coded environment-encoding limits.

**Why this matters for agents:** POET is the conceptual template for **automated
curriculum generation** and **co-evolving the task distribution with the solver** —
the open-endedness pattern that newer self-improving-agent and environment-generation
work (e.g., automated agent design, self-play task generation) draws on.

---

## 5. LLM-driven evolutionary program search

The bridge from classical QD/EC to today's LLM systems: **use an LLM as the
variation operator** (the "mutation"/"crossover" that proposes new candidates),
keep an **automated evaluator** as the fitness function, and maintain an
**evolving database** of programs.

- **FunSearch** (Romera-Paredes et al., **Nature 625:468-475**; online Dec 2023,
  print 2024): pairs a **pretrained
  LLM** (proposes new *programs* — code, not raw answers) with a **systematic
  evaluator** (runs them; guards against hallucination). Best programs are stored in
  a **programs database** that is *islands-structured to maintain diversity* and
  fed back into prompts. Searches the space of **programs that describe *how* to
  solve a problem**, not the solution itself — so outputs are interpretable and
  reusable. Produced the **first new discoveries on open problems** via LLMs: larger
  **cap sets** (extremal combinatorics) and improved **online bin-packing** heuristics.
- **AlphaEvolve** (Novikov et al., **arXiv:2506.13131**, DeepMind; announced May
  2025, arXiv June 2025): a general-purpose **evolutionary coding agent**
  orchestrating an **ensemble of
  Gemini models** (fast Flash for breadth + Pro for quality) that **edit whole code
  files** (not just one function), scored by automated evaluators, evolved in a
  database. Notable results: a **rank-48 algorithm for multiplying two 4×4
  complex-valued matrices** — the first improvement over **Strassen (1969)** for that
  setting (fields of characteristic 0) in 56 years; the complex/characteristic-0
  scope is load-bearing, since AlphaTensor (2022) already reached rank-47 over GF(2) —
  plus real wins inside Google (data-center scheduling recovering
  ~0.7% of compute; a ~23% matmul-kernel speedup that cut Gemini training time ~1%).
  Generalizes FunSearch from single functions to full codebases. (Follow-up at
  scale: Georgiev, Gómez-Serrano, Tao, Wagner, "Mathematical exploration and
  discovery at scale", arXiv:2511.02864, 2025.)

**The common skeleton (memorize this):**

```
  init programs DB  ──►  sample parent(s) from DB (diversity-preserving)
        ▲                          │
        │                          ▼
   add to DB  ◄──  evaluate  ◄──  LLM proposes variant (the "mutation")
```

Everything else (islands vs MAP-Elites archive, single function vs whole file,
which model ensemble) is a knob on this loop.

---

## 6. How ADAS / AFlow / GEPA use this engine (the dedup boundary)

This spoke documents the **search machinery**; the **application** of that machinery
to *designing agents and multi-step LLM programs* lives in
[`automated-agent-system-design`](./automated-agent-system-design.md). The link is
explicit and was made by the same researchers:

- **ADAS / Meta Agent Search** (Hu, Lu, **Clune**, ICLR 2025): a meta-agent
  **iteratively programs new agents in code**, evaluates them, and adds them to an
  **ever-growing archive** used to condition the next design — the paper itself
  states it explores "novel and interesting designs **similar to existing
  open-endedness algorithms**." That archive-conditioned, novelty-seeking loop *is*
  the QD/open-endedness pattern of §1-§5, applied to the space of **agent
  architectures**. (Same Clune as MAP-Elites and POET — this is a deliberate
  lineage, not a coincidence.)
- **AFlow** (Zhang et al., ICLR 2025): **MCTS** over code-represented agentic
  workflow graphs — a tree-search variation/selection engine over the program space.
- **GEPA** (Agrawal et al., 2025): **reflective Genetic-Pareto** evolution of
  prompts/systems — an explicitly *evolutionary* (population + Pareto selection +
  LLM-reflection mutation) optimizer; the **Pareto front is its diversity-preserving
  archive**.

**Routing rule:**

| If the question is about… | Go to |
| --- | --- |
| the archive, behavior descriptors, novelty/QD-score, deception, stepping stones, MAP-Elites/CMA-ME/POET *as algorithms*; the QD/illumination *mechanics* of FunSearch/AlphaEvolve (proposer → evaluator → diversity-preserving programs DB) | **here** (this spoke) |
| using FunSearch/AlphaEvolve (or STOKE/Souper/AlphaTensor/AlphaDev) to **discover a faster concrete algorithm or program** — a matmul routine, a sort, a bin-packing heuristic, a superoptimized instruction sequence; verifier-gated program synthesis | **algorithmic-discovery-superoptimization** (this hub) |
| *which agent or LLM-program* the optimizer should output from a metric (ADAS/AFlow/DSPy/MIPROv2/GEPA/TextGrad/Trace/Darwin-Gödel) | **automated-agent-system-design** |
| gradient RL updating model **weights** (GRPO/PPO/DPO) | **agentic-rl** (ai-llm-model-layer) |
| test-time / inference compute scaling | **reasoning-models** (ai-llm-model-layer) |
| classical single-objective Bayesian optimization / hyperparameter tuning | **da-7-machine-learning** |

---

## 7. Designing the objective, the diversity metric, and the eval

QD/open-ended search lives or dies on three design choices:

1. **Behavior descriptor (BD) / measure space.** Pick axes that are *meaningful*,
   *cheap to compute*, and *aligned with the diversity you actually want*. Too few
   dims → premature convergence; too many → an archive too sparse to fill (use
   **CVT-MAP-Elites** to control cell count). Hand-chosen BDs are the classic
   approach; learned/auto BDs (e.g., AURORA-style) are the frontier.
2. **Quality / fitness vs novelty pressure.** Pure objective → deception &
   collapse. Pure novelty → wanders without improving. QD's insight is **both, with
   competition kept *local*** (NSLC) or **per-cell** (MAP-Elites). Variants add
   selection refinements — the **curiosity score** (Cully & Demiris, 2018) and
   CMA-ME's **improvement-driven emitters**.
3. **Eval / benchmark.** Standard health metrics:
   - **QD-score** = sum of fitnesses of all archive cells (rewards *both* coverage
     and quality); **coverage** = fraction of cells filled; **archive profile** =
     QD-score vs threshold.
   - Deceptive **maze** navigation tasks (Pugh et al., *Quality Diversity: A New
     Frontier*, Frontiers 2016) are the canonical QD benchmark; BipedalWalker
     obstacle courses for POET-style open-endedness.
   - For **LLM program search**, the evaluator *is* the benchmark and must
     **resist hallucinated / overfit programs** — FunSearch's separate verifier
     and AlphaEvolve's automated evaluators exist precisely to stop the LLM from
     "passing" with a wrong program. **Watch for eval-overfitting / reward-hacking**
     (deep treatment in `automated-agent-system-design`).

**Unifying frame (Cully & Demiris, IEEE TEVC 2018):** *every* QD algorithm =
**container** (how solutions are stored/ordered: grid archive vs unstructured
novelty archive) + **selection operator** (how parents are chosen) + **scores**
computed to drive both. If you can name those three for your system, you understand
its QD design.

---

## Quick-reference: the algorithm family

| Algorithm | Year / venue | What's new | One-line role |
| --- | --- | --- | --- |
| Novelty Search | Lehman & Stanley, EC 2011 | reward novelty in behavior space, **no objective** | escapes deception |
| NSLC | Lehman & Stanley, GECCO 2011 | novelty **+ local** competition | direct ancestor of QD |
| MAP-Elites | Mouret & Clune, 2015 | elites grid over chosen feature axes | the workhorse illumination algorithm |
| CVT-MAP-Elites | Vassiliades et al., 2018 | Voronoi cells | scale to high-dim BD spaces |
| CMA-ME | Fontaine et al., GECCO 2020 | CMA-ES emitters fill the archive | QD on continuous domains |
| CMA-MEGA (DQD) | Fontaine & Nikolaidis, NeurIPS 2021 | exploit objective+measure **gradients** | sample-efficient QD |
| CMA-MAE | Fontaine & Nikolaidis, GECCO 2023 | annealed acceptance threshold (CMA-MAEGA = its gradient variant) | stable optimization↔illumination |
| POET / Enhanced POET | Wang et al., 2019 / 2020 | **co-evolve environments + agents** + transfer | open-ended curriculum |
| FunSearch | Romera-Paredes et al., Nature 2023/24 | LLM proposer + evaluator + programs DB | LLM-driven program search |
| AlphaEvolve | Novikov et al., 2025 | Gemini ensemble edits whole files | general evolutionary coding agent |

---

## Sources (primary, verified via exa + WebSearch + arXiv-direct)

- Mouret, J.-B. & Clune, J. (2015). *Illuminating search spaces by mapping elites.*
  arXiv:1504.04909.
- Lehman, J. & Stanley, K. O. (2011). *Abandoning objectives: Evolution through the
  search for novelty alone.* Evolutionary Computation 19(2):189-223. (and GECCO-2010,
  *Efficiently evolving programs through the search for novelty*.)
- Lehman, J. & Stanley, K. O. (2011). *Evolving a diversity of virtual creatures
  through novelty search and local competition.* GECCO 2011 (NSLC).
- Fontaine, M. C. et al. (2020). *Covariance Matrix Adaptation for the rapid
  illumination of behavior space (CMA-ME).* GECCO 2020. 10.1145/3377930.3390232.
- Fontaine, M. C. & Nikolaidis, S. (2021). *Differentiable Quality Diversity
  (CMA-MEGA).* NeurIPS 2021.
- Fontaine, M. C. & Nikolaidis, S. (2023). *Covariance Matrix Adaptation MAP-Annealing
  (CMA-MAE).* GECCO 2023, pp. 456-465, 10.1145/3583131.3590389. (CMA-MAEGA, the
  gradient variant, is the same authors' separate DQD extension, not this title.)
- Vassiliades, V., Chatzilygeroudis, K. & Mouret, J.-B. (2018). *Using Centroidal
  Voronoi Tessellations to scale up the Multidimensional Archive of Phenotypic Elites
  algorithm.* IEEE TEVC 22(4):623-630. (online-first 2017; arXiv:1610.05729.)
- Wang, R., Lehman, J., Clune, J. & Stanley, K. O. (2019). *POET: Paired Open-Ended
  Trailblazer.* arXiv:1901.01753, GECCO 2019. Enhanced POET (Wang et al., ICML 2020).
- Romera-Paredes, B. et al. *Mathematical discoveries from program search with large
  language models (FunSearch).* Nature 625:468-475, 10.1038/s41586-023-06924-6
  (published online 14 Dec 2023; print issue 7995, 18 Jan 2024).
- Novikov, A. et al. (2025). *AlphaEvolve: A coding agent for scientific and
  algorithmic discovery.* arXiv:2506.13131. (Follow-up: arXiv:2511.02864.)
- Stanley, K. O. & Lehman, J. (2015). *Why Greatness Cannot Be Planned: The Myth of
  the Objective.* Springer.
- Pugh, J. K., Soros, L. B. & Stanley, K. O. (2016). *Quality Diversity: A New
  Frontier for Evolutionary Computation.* Frontiers in Robotics & AI 3:40. (QD-score,
  maze benchmark; term coined in Pugh et al., GECCO 2015.)
- Cully, A. & Demiris, Y. (2018). *Quality and Diversity Optimization: A unifying
  modular framework.* IEEE TEVC 22(2):245-259.
- Tjanaka, B. et al. (2023). *pyribs: A bare-bones Python library for Quality
  Diversity optimization.* GECCO 2023. (RIBS: archive + emitter + scheduler.)
- Hu, S., Lu, C. & Clune, J. (2024/2025). *Automated Design of Agentic Systems
  (ADAS / Meta Agent Search).* arXiv:2408.08435, ICLR 2025. — for the engine→
  application boundary; deep coverage in `automated-agent-system-design`.
