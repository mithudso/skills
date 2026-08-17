<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** New `/dr`-researched reference (concept-family-explorer run, 2026-06-15).
> Sibling topics in this family are reference files under the hubs — **not** standalone skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling skill; load that topic's `references/<name>.md` from the owning hub.

---

---
name: algorithmic-discovery-superoptimization
title: Algorithmic Discovery & Superoptimization
description: >
  Discovering and optimizing the ALGORITHM itself — from classical superoptimization and program synthesis to AI-driven algorithm discovery — where a candidate proposer (enumeration / MCMC / SMT / RL / LLM) is gated by a correctness oracle and a fitness/latency evaluator inside a search or evolutionary loop. TRIGGER: superoptimize an instruction sequence or find a missed peephole optimization (STOKE, Souper, Alive2, LPO); reduce an algorithm's operation count or asymptotic cost (matrix-multiply, sorting-network latency); AlphaTensor / AlphaDev / FunSearch; run or build an AlphaEvolve-style evolutionary coding agent (OpenEvolve, ShinkaEvolve) over a verifiable metric; program-synthesis substrate (CEGIS, SyGuS, sketching, GenProg); ML-guided compiler heuristics (MLGO, IR2Vec) or autotuning (OpenTuner, Halide, ATLAS/FFTW); "discover a faster algorithm", "evolve a better heuristic", "verifier-gated program search". SKIP: optimizing an AGENT architecture / multi-step LLM program from a metric (ADAS/AFlow/DSPy/GEPA/TextGrad/Trace/Darwin-Godel) -> automated-agent-system-design; agent plan->edit->test->repair self-repair, Agentless, SBFL -> coding-agents; TRAINING a repair agent / SWE-RL -> agentic-rl; debugging/code-review/perf-profiling-as-process -> software-engineering-patterns; browser/runtime profiling -> performance-profiling-expert; GPU kernel tuning -> llm-gpu-kernels; skill authoring -> claude-code-skills/skill-optimizer.
category: developer
---

# Algorithmic Discovery & Superoptimization

`verified-as-of: 2026-06-15`

The discipline of **discovering or optimizing the algorithm itself** — not tuning a runtime knob, not redesigning an agent, but changing what the program *computes* so it uses fewer operations, lower latency, or a better heuristic, with a **machine-checked guarantee or measured improvement**. Every method here is one instance of a single recurring loop:

> **candidate proposer** (exhaustive enumeration / MCMC / SMT solver / reinforcement learning / LLM) → **correctness oracle** (SMT equivalence proof, test suite, or theorem) → **fitness/latency evaluator** → **search or evolutionary loop** that resamples the best candidates.

The **evaluator/verifier gate is the load-bearing idea**: it is what turns an unreliable proposer (a random mutator or a hallucinating LLM) into a source of *provably or empirically correct* improvements. FunSearch and AlphaEvolve are explicit that the systematic evaluator is the guard against LLM hallucination.[^funsearch][^alphaevolve]

## Scope boundary (read first)

This reference owns **search/optimization OVER an algorithm or program**. It defers to these peers (load each from its owning hub's `references/`):

- **`automated-agent-system-design`** (this hub) — the closest neighbor and the sharpest line. That reference automates the design of an **agent architecture / multi-step LLM program** (ADAS, AFlow, DSPy/MIPROv2, GEPA, TextGrad, Trace/OptoPrime, Darwin-Gödel Machine). This one optimizes the **algorithm/code an ordinary program runs** (a matmul routine, a sort, a bin-packing heuristic, a compiler peephole). The mechanisms rhyme (both are verifier-gated evolutionary/search loops; FunSearch is the shared ancestor), but the *artifact being optimized* differs. When the question is "search over agent topology/prompts," go there; when it is "find a faster matmul / sort / heuristic / instruction sequence," stay here.
- **`coding-agents`** (this hub) — hand-designed plan→edit→test→repair self-repair loops, Agentless, SBFL fault localization. Those fix a *specific* bug in a *specific* codebase; this discovers a *general* better algorithm.
- **`agentic-rl`** (ai-llm-model-layer) — RL that updates **model weights** (SWE-RL, GRPO/PPO). AlphaTensor/AlphaDev use RL, but RL is the *proposer* inside a frozen-weights discovery loop here; this reference is not about training a policy to ship.
- **`open-endedness-quality-diversity-search`** (this hub) — the QD / illumination / open-endedness search *substrate* (MAP-Elites, Novelty Search/NSLC, CMA-ME, POET) and the *mechanics* of FunSearch / AlphaEvolve as evolutionary-search algorithms (proposer → evaluator → diversity-preserving archive). When the question is *how the QD/illumination search itself works*, go there; when it is *discover a faster concrete algorithm/program*, stay here. FunSearch/AlphaEvolve are the shared boundary objects — this reference owns them as algorithm-discovery applications, that one as QD algorithms.
- **`software-engineering-patterns`** — debugging, root-cause (5 Whys), code review, perf-profiling *as a process*, job scheduling.
- **`performance-profiling-expert`** — browser/runtime profiling (flame charts, Web Vitals), NOT algorithmic complexity reduction.
- **`llm-gpu-kernels`** — GPU kernel perf tuning. (AlphaEvolve's FlashAttention/matmul kernel speedups are *examples of its reach*; kernel-craft itself lives there.)
- **`devops-observability`** — telemetry/eBPF/OTel. **`claude-code-skills`/`skill-optimizer`** — skill authoring. **`ai-mcp-sdk-prompting`** — agent memory/context engineering. **`continuous-learning-v2`** — instinct-based learning.

---

## Core Concepts

### 1. Superoptimization & complexity reduction
Superoptimization (Massalin, 1987) searches for the *shortest provably equivalent* instruction sequence — the origin of the field. Three modern lineages:
- **Stochastic superoptimization (STOKE)** — Schkufza, Sharma & Aiken, ASPLOS 2013. Formulates loop-free x86-64 optimization as **MCMC / Metropolis–Hastings** search over program space, with a cost function blending correctness and performance; sacrifices completeness for scope. Starting from `llvm -O0` it matches or beats `gcc -O3` / `icc -O3` and sometimes expert hand-written assembly.[^stoke]
- **SMT-synthesizing superoptimization (Souper)** — Sasnauskas et al., arXiv 2017. Extracts DAGs of LLVM IR, poses **UNSAT queries to an SMT solver** to detect *missed* optimizations, then synthesizes shorter refining expressions via a **CEGIS** loop. Real impact: suggestions shipped in LLVM and Microsoft Visual C++; as an auto pass it shrank a Clang binary ~4.4%. Later work (Hydra, OOPSLA 2024) generalizes one-off finds into reusable rewrite rules verified with **Alive2**.[^souper] Souper's blind spots (memory, FP, vectors) motivate the LLM+verifier hybrids below.
- **Complexity reduction** — lowering the *operation count or asymptotic cost* of a primitive (matmul multiplications, sorting-network comparisons) — distinct from tuning a runtime implementation.

### 2. Program-synthesis & evolutionary substrate
The formal foundation that AI-driven discovery scales:
- **Sketching** (Solar-Lezama) — the programmer supplies a partial program with *holes*; a synthesizer fills them from a spec, reducing synthesis to **2QBF/SAT**.[^sketch]
- **CEGIS** (Counterexample-Guided Inductive Synthesis) — a synthesize↔verify loop where an SMT oracle returns counterexamples that constrain the next candidate; converges fast because a few corner-case inputs cover the space.[^sketch]
- **SyGuS** (Syntax-Guided Synthesis; Alur et al., 2013) — a *grammar* restricts the search space; the formal substrate under superoptimizers (Souper uses CEGIS).
- **Genetic programming / GenProg** — Le Goues et al., IEEE TSE 2012. Evolves **AST-level patches**; the test suite *is* the fitness function (positive tests = required behavior, negative test = the bug); mutation/crossover localized along a fault-localization-weighted path. The direct conceptual ancestor of LLM-as-mutation-operator agents.[^genprog]

### 3. AlphaTensor — TensorGame for matrix multiplication
Fawzi et al., *Nature* 610 (2022). Reframes algorithm discovery as a **single-player game (TensorGame)**: find low-rank decompositions of the matrix-multiplication tensor; built on **AlphaZero** with a transformer policy/value net. Found a **4×4 algorithm in GF(2) using 47 multiplications**, beating Strassen's two-level recursive 49 (naïve is 64) — **over GF(2) (the field with two elements); this does *not* beat Strassen over the reals**. Improved SOTA for >70 matrix sizes; can optimize for *real hardware runtime* (not just mult-count) — 10–20% faster on NVIDIA V100 and Google TPU; the repo publishes 14,236 nonequivalent algorithms with correctness-verification Colabs.[^alphatensor]

### 4. AlphaDev — AssemblyGame for sorting/hashing
Mankowitz et al., *Nature* 618:257–263 (2023). The **AssemblyGame**: state = current program + register/memory; each step appends one CPU instruction; reward = correctness + measured latency. AlphaDev (AlphaZero + a MultiQuery Transformer) discovered sort3/4/5 routines from scratch — including a branchless "AlphaDev swap" — that were reverse-engineered to C++ and **merged into LLVM libc++'s standard sort** (first change in over a decade): up to **70% faster for short sequences**. Proves correctness via the Assembly Zero-One Principle; also applied to hashing and protobuf deserialization.[^alphadev]

### 5. FunSearch — LLM proposer + systematic evaluator
Romera-Paredes et al., *Nature* (announced Dec 2023; published *Nature* vol. 625, 2024). Pairs a pretrained **LLM (code proposer)** with an **automated systematic evaluator** to guard against hallucination. It searches the **function/heuristic space** — evolving programs that *describe how to solve* a problem (interpretable heuristics), not the raw solution object — using an **island-based evolutionary algorithm** with periodic resets for diversity. First LLM-driven *new* discovery on an open problem: improved **cap-set** lower bounds (with Jordan Ellenberg) and found **online bin-packing** heuristics beating best-fit/first-fit (relevant to data-center job allocation). The direct precursor to AlphaEvolve.[^funsearch]

### 6. AlphaEvolve — the evolutionary coding-agent keystone
Novikov et al., DeepMind, arXiv:2506.13131 (2025). Evolves **entire codebases** (hundreds of lines), not single functions. Architecture: **prompt sampler → ensemble of LLMs (Gemini Flash for breadth + Gemini Pro for depth) → automated evaluators → programs database** implementing an evolutionary algorithm that selects parents for future prompts; LLMs emit **code diffs**. Production wins: a **Borg scheduling heuristic in production >1 year recovering ~0.7% of Google's global compute**; a Verilog TPU arithmetic-circuit simplification; a 23% Gemini matmul-kernel speedup; up to 32.5% FlashAttention speedup. Discovered a **4×4 complex-valued matmul in 48 scalar multiplications** — first improvement over Strassen's 49 in 56 years in that (complex-valued) setting. On 50+ open math problems: rediscovered SOTA ~75% of the time, improved best-known ~20% (e.g. a new 11-dimensional kissing-number lower bound: 593 spheres). **General requirement: any problem with an automatically verifiable/scorable solution.**[^alphaevolve]

### 7. Open descendants & the sample-efficiency frontier
Evaluator calls (thousands) are the binding cost constraint. **ShinkaEvolve** (Sakana AI, arXiv:2509.19349, 2025; Apache-2.0) hits a SOTA circle-packing (n=26) solution in **~150 samples** via three ideas: parent sampling balancing exploration/exploitation, **novelty-based rejection** via code-embedding similarity, and **bandit-based task-dependent LLM selection** across an ensemble. **OpenEvolve** (codelion/ASI Labs, 2025) is the reference open clone — prompt sampler + LLM ensemble + evaluator pool + **MAP-Elites / island-model** program database with migration; replicated circle-packing to 99.97% of DeepMind's result. Broader ecosystem: **LLM4AD**, **EoH** (Evolution of Heuristics), LLaMEA, MEoH (survey arXiv:2410.14716).[^shinka]

### 8. LLM + formal-verifier hybrids for classical optimization
**LPO** (arXiv:2508.16125, 2025) closes the loop on classical superoptimization: an **LLM creative proposer + a formal verifier** discovers *missed peephole optimizations that solver-based tools like Souper cannot* (the memory/FP/vector cases). The pattern — let the LLM be creative, let the verifier be the gate — is the same hallucination-guard as FunSearch/AlphaEvolve applied to compiler rewrites.[^souper]

### 9. ML-guided compiler optimization & autotuning (neighbors)
These search the *implementation* space, not the algorithm — the boundary this skill polices against `llm-gpu-kernels` and `performance-profiling-expert`:
- **MLGO** (Trofin et al., 2021; Google) — the first ML *policy* shipped in a production compiler (LLVM), replacing heuristics: **inlining-for-size** (up to ~7% vs `-Oz`) and **register-allocation eviction** (0.3–1.5% QPS). Trained via Policy Gradient / Evolution Strategies; ships **IR2Vec/MIR2Vec** code embeddings. It learns a *reusable policy*, distinct from per-program autotuning.[^mlgo]
- **Autotuning** — empirical search over a *parameterized implementation* space. **OpenTuner** (Ansel et al., PACT 2014) uses an **ensemble of search techniques with an AUC-bandit** budget allocator; lineage includes **ATLAS** (BLAS) and **FFTW** (FFT plans). **Halide** decouples algorithm from schedule; learned autoschedulers (Adams et al. 2019) use tree search + an ML cost model trained on random programs to beat expert-tuned code. Key difficulty: search spaces are huge, mostly-invalid, and have poor locality.[^opentuner]

### 10. Provable vs measured trust
A defining design axis: which discoveries are **formally verified** (AlphaTensor decompositions, Souper SMT proofs, AlphaEvolve math results) versus **empirically benchmarked** (kernel/scheduler heuristics, STOKE samples). Always state which guarantee a given output carries.

---

## References (sub-files a full skill would hold)

- `references/superoptimization-classical.md` — Massalin origins; STOKE MCMC cost-function search; Souper SMT + CEGIS; peephole superoptimization (Bansal–Aiken); Alive2/Hydra rule generalization; LPO LLM+verifier hybrid.
- `references/synthesis-evolutionary-substrate.md` — sketching → 2QBF/SAT; CEGIS loop; SyGuS grammars; GenProg AST genetic programming; how LLMs slot in as proposer and SMT/tests as oracle.
- `references/ai-algorithm-discovery.md` — AlphaTensor (TensorGame), AlphaDev (AssemblyGame, libc++), FunSearch (island evolution over heuristics); the proposer+oracle+evaluator+loop pattern.
- `references/alphaevolve-and-descendants.md` — AlphaEvolve architecture, diff-based whole-codebase evolution, Google-infra deployments, complex-matmul/kissing-number results; OpenEvolve, ShinkaEvolve sample-efficiency tricks, LLM4AD/EoH ecosystem.
- `references/compiler-ml-and-autotuning.md` — MLGO learned inlining/regalloc, IR2Vec; OpenTuner ensembles, ATLAS/FFTW, Halide schedule search + learned cost models; the algorithm-vs-implementation boundary.

[^alphatensor]: Fawzi et al., "Discovering faster matrix multiplication algorithms with reinforcement learning," *Nature* 610 (2022). DOI 10.1038/s41586-022-05172-4.
[^alphadev]: Mankowitz et al., "Faster sorting algorithms discovered using deep reinforcement learning," *Nature* 618:257–263 (2023). DOI 10.1038/s41586-023-06004-9.
[^funsearch]: Romera-Paredes et al., "Mathematical discoveries from program search with large language models," *Nature* vol. 625 (2024; announced Dec 2023). DOI 10.1038/s41586-023-06924-6.
[^alphaevolve]: Novikov et al., "AlphaEvolve: A coding agent for scientific and algorithmic discovery," arXiv:2506.13131 (2025); DeepMind blog 2025-05-14.
[^stoke]: Schkufza, Sharma & Aiken, "Stochastic Superoptimization," ASPLOS 2013.
[^souper]: Sasnauskas, Chen, Ketema, Taneja & Regehr, "Souper: A Synthesizing Superoptimizer," arXiv:1711.04422 (2017); Hydra OOPSLA 2024; LPO arXiv:2508.16125 (2025).
[^sketch]: Solar-Lezama, "The Sketching Approach to Program Synthesis"; Alur et al., "Syntax-Guided Synthesis," FMCAD 2013.
[^genprog]: Le Goues, Nguyen, Forrest & Weimer, "GenProg: A Generic Method for Automatic Software Repair," IEEE TSE 2012.
[^shinka]: ShinkaEvolve, Sakana AI, arXiv:2509.19349 (2025); OpenEvolve (codelion/ASI Labs, 2025); survey arXiv:2410.14716.
[^mlgo]: Trofin et al., "MLGO: a Machine Learning Guided Compiler Optimizations Framework" (2021); LLVM MLGO docs.
[^opentuner]: Ansel et al., "OpenTuner: An Extensible Framework for Program Autotuning," PACT 2014; Adams et al., Halide learned autoscheduler (2019).
