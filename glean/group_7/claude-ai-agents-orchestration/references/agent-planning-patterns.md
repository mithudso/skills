<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Formerly the standalone `agent-planning-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

<!-- FOLDED REFERENCE under the ai-agent-engineering hub — agent planning & control-flow patterns (the inference-time control loop). Sibling refs: autonomous-loops (loop infrastructure), reasoning-models (single-turn CoT / test-time compute), agent-harness-construction (action space & tools), multi-agent-orchestration (multi-agent topologies), agent-state-and-durable-execution (state & runtime). Any "use the X" pointer refers to a sibling reference under this hub. -->

# Agent Planning & Control-Flow Patterns

How a **single** LLM agent decides its next action at inference time — the control loop. Three families: reasoning-interleaved acting (adaptive, costly), plan-first (efficient, less adaptive), and iterative self-improvement — plus the workflow-vs-agent altitude decision that sits above all of them.

> This is the control-flow PATTERN layer. For loop *infrastructure* see `autonomous-loops`; for single-turn CoT/test-time compute see `reasoning-models`; for the action/tool contract see `agent-harness-construction`.

## 1. ReAct — reason+act interleaving (the baseline)
Interleave a reasoning trace and a tool action each step (Yao et al. 2022, arXiv:2210.03629). Reasoning tracks/updates the plan and handles exceptions; actions ground each step in real tool output, which **reduces hallucination and error cascades** vs pure chain-of-thought.
- **Failure modes (treat as the norm, not edge cases):** infinite loops retrying the same failed action; no native mechanism to revert a wrong step; **error propagation** where one early mistake distorts all later reasoning (acute in long-horizon tasks); each hallucinated tool name burns retry budget.
- **Mandatory guards:** max-step cap, loop / no-progress detection, deterministic tool routing, grounding-validation before returning.
- **Cost:** highest per task — every step re-ingests the full accumulated trace, so context grows monotonically.

## 2. Plan-first — Plan-and-Execute, ReWOO, LLMCompiler
Decouple planning from execution to cut tokens/latency at the cost of adaptivity.
- **Plan-and-Solve / Plan-and-Execute** (Wang et al. 2023, arXiv:2305.04091): devise the full plan, then a separate executor runs each subtask; landed in LangChain as "Plan-and-Execute." Beats zero-shot CoT by fixing missing-step errors.
- **ReWOO** (Xu et al. 2023, arXiv:2305.18323): *Reasoning WithOut Observation* — Planner writes the whole plan with tool-call placeholders up front, Worker fills results, Solver composes. Reported **~5× token efficiency + ~4% accuracy on HotpotQA** and lets you offload planning to a smaller model.
- **LLMCompiler** (Kim et al. 2023, arXiv:2312.04511): Planner emits a **DAG of tasks with dependencies**, a fetching unit dispatches as deps resolve, executor runs independent calls **in parallel** — up to **3.7× latency / 6.7× cost** vs ReAct. Only helps when calls are independent.
- **Key tradeoff:** plan-first assumes the plan is right *without* seeing intermediate results. When an early observation should change course, you need **replanning** — which erodes the token savings. Single-study numbers above are vs ReAct; read as directional, not independently replicated.

## 3. Iterative self-improvement — Reflexion, Self-Refine
- **Self-Refine** (Madaan et al. 2023): one model generates → critiques → refines, no training; ~20% avg lift, most gains in the first iteration.
- **Reflexion** (Shinn et al. 2023): "verbal RL" — reflect on a feedback signal, store reflective text in episodic memory, retry. **Relies on an external/environmental signal** (unit tests, task success).
- **The over-hyped claim, corrected:** Huang et al. 2023 (arXiv:2310.01798), *"LLMs Cannot Self-Correct Reasoning Yet"* — with **intrinsic** self-correction (no external feedback), models often fail to improve and can **degrade**. **Operator rule: self-critique needs an external oracle** (tests, a verifier, tool errors, a stronger evaluator). This is the grounded form of the evaluator-optimizer workflow below.

## 4. Tree/Graph-of-Thoughts applied to acting
ToT (Yao 2023, arXiv:2305.10601) searches a tree of "thoughts" with lookahead + backtracking; GoT (Besta 2023) generalizes to a graph (aggregate/merge). Applied to ACTING (e.g., LATS), the search/evaluate/prune/backtrack machinery searches over **action plans**. **Steep token/latency multiplier** — reserve for high-value, expensive-to-get-wrong tasks with a strong state evaluator.

## 5. Workflows vs agents (Anthropic, "Building Effective Agents", Dec 2024)
The altitude decision above every pattern. **Augmented LLM** (retrieval + tools + memory) is the building block. **Workflows** orchestrate LLMs through *predefined code paths*; **agents** let the LLM *dynamically direct its own process*.
Five workflow patterns: **prompt chaining** (fixed sequence + gate checks) · **routing** (classify → specialized handler) · **parallelization** (sectioning + voting) · **orchestrator-workers** (central LLM decomposes/delegates/synthesizes when subtasks aren't predictable) · **evaluator-optimizer** (generate ↔ critique loop with clear criteria).
**Guidance:** start with the simplest solution; prefer a **workflow** when the task is decomposable in advance (cheaper, lower-latency, predictable); reach for an **autonomous agent** only when you need model-driven flexibility at scale — and don't add a framework that hides the prompts.

## Decision guide
- Uncertain path / exploratory → **ReAct** + step cap + loop detection + grounding check.
- Token/latency-bound and plan is largely knowable → **ReWOO / Plan-and-Execute** (+ replanning if observations can invalidate the plan).
- Multiple **independent** tool calls → **LLMCompiler** (parallel DAG).
- Refinement with a real signal → **evaluator-optimizer / Reflexion** (never trust intrinsic-only self-correction).
- Decomposable in advance → a **workflow**, not an agent.

## Anti-patterns
Unbounded ReAct loops; trusting intrinsic self-correction; ToT-over-actions on cheap tasks; reaching for an autonomous agent (or a heavy framework) when a fixed workflow would be cheaper and more predictable.

## References
ReAct (arXiv:2210.03629) · Plan-and-Solve (arXiv:2305.04091) · ReWOO (arXiv:2305.18323) · LLMCompiler (arXiv:2312.04511) · Reflexion (arXiv:2303.11366) · Self-Refine (arXiv:2303.17651) · "LLMs Cannot Self-Correct Reasoning Yet" (arXiv:2310.01798) · Tree of Thoughts (arXiv:2305.10601) · Graph of Thoughts (arXiv:2308.09687) · Anthropic, "Building Effective Agents" (anthropic.com/engineering/building-effective-agents, Dec 2024).
