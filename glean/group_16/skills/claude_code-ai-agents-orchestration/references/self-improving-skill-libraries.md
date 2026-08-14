<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** New `/dr`-researched reference (concept-family-explorer run, 2026-06-15).
> Sibling topics in this family are reference files under the hubs — **not** standalone skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling skill; load that topic's `references/<name>.md` from the owning hub.

---

---
name: self-improving-skill-libraries
title: Self-Improving Skill Libraries & Experience-to-Skill Distillation
description: >
  The connective methodology of compiling learned wins — trajectories, reflections, execution feedback — into durable, reusable agent skills/workflows/playbooks WITHOUT weight updates: the distill->admit->index/trigger->retrieve->refine/deprecate lifecycle; Voyager automatic executable-skill libraries; ACE Generator-Reflector-Curator delta-bullet playbooks (context-collapse & brevity-bias failure modes, grow-and-refine fix); reflection-to-procedural-memory loops (Reflexion->Memp/AWM/ReasoningBank/ExpeL); self-judged contrastive distillation from successes AND failures; trigger-quality + helpful/harmful utility feedback; library-lifecycle operators with admission/verification/lineage; and guardrails against skill bloat, negative transfer, and drift. TRIGGER: "turn our agent's wins into reusable skills", "build a self-improving skill/playbook library", "distill trajectories/reflections into procedures", "when should a learned win earn a durable skill", "curate/prune an evolving skill library", "why is my growing skill library hurting accuracy", "ACE / Voyager / ReasoningBank / Agent Workflow Memory / procedural memory". SKIP: instinct/confidence-based learning -> continuous-learning-v2; skill anatomy/authoring/quality-gates -> claude-code-skills / skill-optimizer; automated agent & compound-system search/optimization (ADAS/GEPA/Trace/DGM) -> automated-agent-system-design; context/memory architecture, CLAUDE.md, prompt caching, iterative self-refinement -> ai-mcp-sdk-prompting; plan->edit->test->repair self-repair -> coding-agents; TRAINING a repair agent (SWE-RL) -> agentic-rl.
category: developer
---

# Self-Improving Skill Libraries & Experience-to-Skill Distillation

`verified-as-of: 2026-06-15`

The **connective methodology** for turning an agent's *learned wins* — full task trajectories, verbal reflections, and execution feedback — into **durable, reusable skills / workflows / playbooks**, with the model's **weights frozen**. The unifying move: treat accumulated experience as an *optimizable artifact store* that is curated over time, rather than discarding each episode or hoarding raw transcripts. This is self-improvement *as an alternative to fine-tuning* — adaptation lives in the skill/memory layer, not the parameters.[^ace][^voyager]

The canonical reference loop this skill owns:

> **capture** (trajectory + reflection/execution feedback) → **distill** (to a procedure/insight/routine) → **admit** (gate on verification) → **index / trigger** (embed + write the `when_to_use`) → **retrieve / compose** (fetch top-k, synthesize) → **refine / deprecate** (update-in-place, merge, prune, version).

## Scope boundary (read first)

This reference owns **how learned experience BECOMES a durable skill artifact, and how the library is governed**. It defers to these peers:

- **`continuous-learning-v2`** — instinct/confidence-based learning and skill evolution from instincts. This reference is about distilling *trajectory/reflection/execution* evidence into procedures, not confidence-tracked instincts.
- **`claude-code-skills` / `skill-optimizer`** — skill *anatomy* (SKILL.md, frontmatter), authoring, quality gates, and optimization. This reference owns only how a distilled win *earns and lands in* such a skill; it does not teach SKILL.md structure or run the quality gate.
- **`ai-agents-orchestration` (automated-agent-system-design)** — **automated search/optimization OVER an agent system** (ADAS, AFlow, DSPy/MIPROv2, GEPA, TextGrad, Trace/OptoPrime, Darwin-Gödel). That family optimizes the *system around* a frozen model from a metric. This reference evolves the *experience/skill memory*; when the question is "search the prompts/workflow/architecture from a reward," go there.[^selfevolve]
- **`ai-mcp-sdk-prompting` (llm-context-engineering)** — agent-memory architecture, CLAUDE.md, prompt caching, context-window design. This reference borrows that substrate but owns the *distillation policy* on top of it.
- **`ai-mcp-sdk-prompting`** (self-refinement spoke) — Reflexion/Self-Refine/CRITIC as *iterative refinement loops* with stop-conditions. Here, Reflexion is only the **reflection SOURCE** that feeds procedural memory; the inner refine loop lives there.[^reflexion]
- **`coding-agents`** — the hand-designed plan→edit→test→repair self-repair loop. **`agentic-rl`** — *training* a repair model (SWE-RL). This reference keeps weights frozen and consumes those loops' feedback as a distillation signal only.

---

## Core Concepts

### 1. The experience-to-skill distillation pipeline (the reference workflow)
Every system in this space instantiates the same six stages above. The central design choice is the **representation**: store episodes *verbatim* vs distill to *abstractions*. The empirical answer is consistent — **distilled, abstract procedures transfer more reliably than raw trajectories**, and distilled-strategy beats stored-episode on web-browsing and SWE benchmarks.[^memp][^reasoningbank][^continual] Memp frames this contrast explicitly and makes procedural memory a *first-class optimization target* with a **Build / Retrieve / Update** operator triad plus explicit **deprecation** of stale procedures.[^memp]

### 2. Voyager — the automatic executable-skill library
The founding system: an LLM (GPT-4 via black-box queries, **no fine-tuning**) that **writes, stores, indexes, retrieves, and composes executable code skills** in Minecraft. Three parts — an automatic curriculum, an **ever-growing skill library** of code, and an iterative prompting loop using environment feedback + execution errors + self-verification to refine a program *before* it is committed as a skill. Each skill is **indexed by the embedding of its LLM-generated description**; retrieval queries with the embedding of the task plan + feedback to fetch top-5. **Complex skills are synthesized by composing simpler ones** — compositionality compounds capability and *alleviates catastrophic forgetting*. Results: 3.3× more unique items, tech-tree milestones up to 15.3× faster than prior SOTA, and the learned library **transfers to a fresh world**.[^voyager]

### 3. ACE — Generator–Reflector–Curator delta-bullet playbooks
Agentic Context Engineering treats context as an **evolving playbook** that accumulates and organizes strategies — self-improvement without weight updates. Three modular roles: the **Generator** produces reasoning trajectories; the **Reflector** distills concrete insights from successes/errors; the **Curator** merges insights as compact **delta bullets** via *deterministic, non-LLM* logic (so deltas batch/merge in parallel). Each bullet carries metadata — a unique id plus **helpful/harmful counters** (the utility signal). Reported: +10.6% on agents, +8.6% on finance, matching a top production agent (IBM-CUGA/GPT-4.1) on AppWorld with a *smaller* open model at **86.9% lower adaptation latency**, from *unlabeled* execution feedback.[^ace]

### 4. Context collapse & brevity bias — the two named failure modes
Naive self-improving context fails two specific ways. **Brevity bias**: summarization-style updates drop domain insight. **Context collapse**: iterative *full rewrites* erode accumulated detail. ACE's fix is **structured incremental delta updates + grow-and-refine** — append new bullets, update in place, **de-duplicate via semantic embeddings**, and prune lazily or proactively. Dynamic Cheatsheet established the proto-principle: keep memory toward *concise, transferable* snippets rather than full transcripts, avoiding context ballooning.[^ace][^dc]

### 5. Reflection-to-procedural-memory loop
Reflexion is the upstream signal: agents **verbally reflect** on task feedback (scalar or free-form, external or simulated) and store the reflective text in an **episodic buffer** to improve subsequent trials — reinforcement via *language*, not weights (91% pass@1 on HumanEval).[^reflexion] This skill consumes that reflective text and *compiles* it into reusable procedures (Memp/AWM/ReasoningBank). The reflection loop itself (stop-conditions, oscillation) is deferred to the self-refinement spoke.

### 6. Workflow induction — Agent Workflow Memory (AWM)
AWM **induces commonly-reused routines (workflows)** from past trajectories and selectively injects them into future generations. It works **offline** (induce from training examples beforehand) and **online** (induce from test queries on the fly: *induce → integrate → utilize* after each task — the trigger-quality feedback loop in action). Gains: +24.6% (Mind2Web), +51.1% (WebArena) relative success **while reducing steps taken** — reuse cuts cost, not just raises accuracy — and online AWM generalizes cross-task/site/domain by 8.9–14.0 absolute points as the train–test gap widens.[^awm]

### 7. Self-judged contrastive distillation from successes AND failures
ReasoningBank distills **generalizable reasoning strategies from BOTH successful and failed self-judged trajectories** — *failure reflection is a first-class memory item, not discarded*. Closed loop: retrieve → act → LLM-as-judge self-assesses → extract success insights or failure reflections → consolidate back; robust to noisy self-judgement. ExpeL established the **contrastive pattern** earlier — extracting natural-language insights by *contrasting* successful vs failed trajectories across a training-task pool (the offline "experience pool → insight bank" pipeline).[^reasoningbank][^expel]

### 8. Memory-aware test-time scaling (MaTTS)
ReasoningBank introduces a **bidirectional synergy**: allocate more compute per task to generate *diverse* trajectories whose self-contrast yields **higher-quality distilled memory**, and better memory then guides **more effective scaling**. This establishes *experience-driven memory* as a distinct scaling axis alongside parameters and inference compute.[^reasoningbank]

### 9. Trigger-quality & helpful/harmful utility feedback
What gates *when a skill loads* is the load-bearing curation signal. In Anthropic Agent Skills, the `description`/`when_to_use` frontmatter (truncated at **1,536 chars** in the listing) is the **model-invoked trigger** that decides load — making trigger quality the curation signal for any library, human or automated.[^skills] In ACE, the per-bullet **helpful/harmful counters** plus reranking are the utility signal that drives keep/prune decisions.[^ace] Trigger *engineering* of SKILL.md itself defers to `skill-optimizer`.

### 10. The skill-library lifecycle operators and the admission seven-tuple
The dynamic-skills survey formalizes library-level dynamics as a set of operators mapping `library_t → library_{t+1}`. This reference's connective synthesis names the working set as a **ten-operator algebra** — **ADD, REFINE, MERGE, SPLIT, PRUNE, DISTILL, ABSTRACT, COMPOSE, REWRITE, RERANK** (the explicit anti-bloat / curation toolkit; *a synthesis label, not a verbatim construct from a single source*). The survey also extends the options formalism toward a **seven-tuple** view (applicability, policy, termination, interface, **EDIT**, **VERIFICATION**, **LINEAGE**) — making admission gates, provenance, and verification *mandatory, not optional*.[^dynamic]

### 11. Negative transfer & skill-bloat guardrails
The headline empirical finding: **admission and repair matter more than raw skill count**; verifier quality is load-bearing; flat retrieval *degrades as the library grows* (all reported by the dynamic-skills survey).[^dynamic] The continual-learning study quantifies the harm mechanism: **negative transfer concentrates on hard cases** — cross-task retrieval disproportionately harms tasks the agent *cannot yet solve*, while solved tasks are robust; and finer organization is *not* universally good — designs with strong forward transfer can simultaneously induce severe forgetting. External memory **does not solve continual learning; it reshapes it** into a representation/retrieval design problem.[^continual]

### 12. Distilled-artifact portability and the productionized endpoint
Distilled procedures are **portable, not model-locked**: Memp shows procedural memory built from a *stronger* model transfers to a *weaker* model and still yields gains; abstract procedures transfer more reliably than raw trajectories.[^memp][^continual] The durable endpoint of the pipeline is a portable, governed, model-discoverable **Agent Skill** — composable, portable across surfaces, efficient via progressive disclosure, and able to ship executable code where deterministic code beats token generation; now an **open standard** (agentskills.io) with org-wide management and a `/v1/skills` API for versioning and governance.[^skills]

### 13. Offline vs online distillation loop selection
Two regimes: **offline** — build an insight/skill bank from a *training-task experience pool* before deployment (ExpeL, offline AWM); and **online** — induce-integrate-utilize *streaming* per task (online AWM, ReasoningBank). Online generalizes better as the train–test gap widens but demands the anti-bloat guardrails; offline is safer to verify but can be stale under drift.[^awm][^expel][^reasoningbank]

---

## References (sub-files a full skill would hold)

- `references/distillation-pipeline.md` — the capture→distill→admit→index→retrieve/compose→refine/deprecate reference workflow; representation choice (verbatim vs abstraction) and the Build/Retrieve/Update triad.
- `references/voyager-executable-skill-libraries.md` — automatic curriculum, embedding-indexed code skills, composition, transfer.
- `references/ace-playbook-curation.md` — Generator/Reflector/Curator, delta bullets, helpful/harmful counters, grow-and-refine, context-collapse & brevity-bias fixes; Dynamic Cheatsheet precursor.
- `references/reflection-and-contrastive-distillation.md` — Reflexion as source; ReasoningBank/ExpeL success-vs-failure contrast; MaTTS scaling.
- `references/workflow-induction-awm.md` — offline vs online/streaming routine induction; step-reduction economics.
- `references/lifecycle-algebra-and-governance.md` — the curation-operator set, the admission seven-tuple, admission gates, verification, lineage; fast/slow update clocks.
- `references/guardrails-negative-transfer-bloat.md` — hard-case harm, forgetting-vs-forward-transfer tradeoff, flat-retrieval degradation, drift/non-stationarity, versioned rollback.
- `references/trigger-quality-and-portability.md` — when_to_use/description triggers, reranking, distilled-artifact portability, landing in governed Agent Skills.

## Key landmarks

- **Founding artifact-library system:** Voyager (Wang et al., NeurIPS 2023 ALOE / TMLR 2024).[^voyager]
- **Frontier playbook method:** ACE (Stanford/SambaNova/UC Berkeley, Oct 2025, ICLR 2026).[^ace]
- **Experience-as-scaling-axis:** ReasoningBank + MaTTS (Google, Sep 2025, ICLR 2026).[^reasoningbank]
- **Governance backbone:** the dynamic-agent-skills lifecycle survey (OpenReview cjU3YbcRr8; titled "Dynamic Agent Skills: A Lifecycle Survey and Taxonomy of Evolving Skill Libraries", also circulated as "They Are Not Static").[^dynamic]
- **Productionization:** Anthropic Agent Skills as an open standard (agentskills.io).[^skills]

[^voyager]: Wang et al., "Voyager: An Open-Ended Embodied Agent with Large Language Models," arXiv:2305.16291 (NeurIPS 2023 ALOE workshop / TMLR 2024). https://arxiv.org/abs/2305.16291
[^ace]: Zhang, Hu et al., "Agentic Context Engineering: Evolving Contexts for Self-Improving Language Models," arXiv:2510.04618 (ICLR 2026). https://arxiv.org/abs/2510.04618
[^reflexion]: Shinn et al., "Reflexion: Language Agents with Verbal Reinforcement Learning," arXiv:2303.11366 (NeurIPS 2023). https://arxiv.org/abs/2303.11366
[^memp]: Fang et al., "Memp: Exploring Agent Procedural Memory," arXiv:2508.06433 (Aug 2025). https://arxiv.org/abs/2508.06433
[^awm]: Wang, Fried, Neubig, "Agent Workflow Memory," arXiv:2409.07429 (ICML 2025). https://arxiv.org/abs/2409.07429
[^reasoningbank]: Ouyang et al., "ReasoningBank: Scaling Agent Self-Evolving with Reasoning Memory," arXiv:2509.25140 (Google; ICLR 2026). https://arxiv.org/abs/2509.25140
[^skills]: Anthropic, "Equipping agents for the real world with Agent Skills," Oct 16 2025 (open standard, agentskills.io). https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills
[^dynamic]: "Dynamic Agent Skills: A Lifecycle Survey and Taxonomy of Evolving Skill Libraries" (also circulated as "They Are Not Static: A Survey of Dynamic Agent Skills"), OpenReview id=cjU3YbcRr8, 2025–2026. https://openreview.net/forum?id=cjU3YbcRr8
[^continual]: "When Continual Learning Moves to Memory: A Study of Experience Reuse in LLM Agents," arXiv 2026. https://arxiv.org/abs/2604.27003
[^dc]: Suzgun, Yuksekgonul et al., "Dynamic Cheatsheet: Test-Time Learning with Adaptive Memory," arXiv:2504.07952 (EACL 2026). https://arxiv.org/abs/2504.07952
[^selfevolve]: Gao et al., "A Survey of Self-Evolving Agents: What, When, How, and Where to Evolve," TMLR 2026, arXiv:2507.21046. https://arxiv.org/abs/2507.21046
[^expel]: Zhao et al., "ExpeL: LLM Agents Are Experiential Learners," AAAI 2024. https://ojs.aaai.org/index.php/AAAI/article/view/29936
