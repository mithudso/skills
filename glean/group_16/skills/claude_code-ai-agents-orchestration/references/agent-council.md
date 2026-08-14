<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Formerly the standalone `agent-council` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: agent-council
version: 1.1.0
updated: 2026-05-29
description: >
  Multi-agent council orchestration — consensus algorithms, debate patterns,
  voting mechanisms, adversarial review, LLM-as-judge, constitutional AI
  alignment, and critique-revision loops. TRIGGER: user says "summon the
  council", "ask other AIs", wants multiple AI perspectives, needs adversarial
  review of agent output, or is designing multi-agent decision systems.
  SKIP: single-agent orchestration without deliberation; general LLM evaluation
  without multi-agent setup.
origin: local
tags: [multi-agent, consensus, debate, adversarial-review, council, llm-as-judge]
related_skills: [agent-ecosystem, agent-plan-writing, agent-workflow-builder_ai_toolkit]
---

# Agent Council

Collect multiple AI opinions and synthesize one answer through structured consensus.

## When NOT to use this skill

- Single-agent tasks where deliberation adds no value — use `agent-ecosystem` for framework selection
- Latency-critical interactive UX (council averages 8.4s vs 2.8–5.6s for individual models)
- Simple classification tasks where one model already exceeds 95% accuracy
- When API budget cannot sustain the N-agent multiplier (a 5-agent council with 3 debate rounds costs 15x a single call)

## Running a council job

```bash
# One-shot
./skills/agent-council/scripts/council.sh "your question here"

# Multi-step
JOB_DIR=$(./skills/agent-council/scripts/council.sh start "your question here")
./skills/agent-council/scripts/council.sh wait "$JOB_DIR"
./skills/agent-council/scripts/council.sh results "$JOB_DIR"
./skills/agent-council/scripts/council.sh clean "$JOB_DIR"
```

References: `references/overview.md`, `references/examples.md`, `references/config.md`, `references/requirements.md`, `references/safety.md`.

---

## Why councils work

Assuming independent expert errors, the probability all N models hallucinate identically on the same claim is the product of individual rates. With three models at rates 0.18, 0.16, and 0.19, simultaneous hallucination drops to ~0.005 — a 97% reduction vs the worst individual model (Council Mode, arXiv 2604.02923).

**Key stats:**
- Hallucination reduction: 35.9% over best single model (HaluEval)
- Bias variance: council achieves 0.003 vs individual 0.021–0.028 (85–89% reduction)
- Complex reasoning (10 steps): 71.2% vs 50.8% for best single model
- Minimum viable quorum: 3 heterogeneous agents

---

## Consensus algorithms

### Plurality voting
Each agent independently answers; most common wins. Use for discrete classification where speed matters.

### Weighted voting
Agents weighted by historical domain accuracy. Use when agents have known heterogeneous strengths.

### Structured synthesis (Council Mode, arXiv 2604.02923)
Organizes expert outputs into four categories: consensus points, disagreements, unique findings, comprehensive analysis. A dedicated synthesis model applies a protocol with mandatory rules. Outperforms simple majority voting — replacing structured synthesis with majority voting increases hallucination by 32.7%.

### Oracle verification
When agents disagree on facts, trigger external tool calls (search, database) to ground consensus in objective reality.

---

## Debate patterns

### Round-robin debate
```
Phase 1: All agents generate independent proposals
Phase 2: Each agent critiques all other proposals
Phase 3: Each agent revises based on critiques
Phase 4: Repeat phases 2–3 until convergence or round limit
Phase 5: Synthesis or vote
```
**Termination:** true consensus, round limit with fallback voting, or escalation to supervisor.

### Judge-Critic loop
A dedicated Judge agent evaluates Debater arguments without generating solutions itself. Debaters refine based on judge feedback. Judge never generates — only evaluates.

### Builder-Critic (refinement loop)
For code generation, legal analysis, document drafting:
```
Builder -> Draft v1
Critic  -> [Bug #1, Bug #2]
Builder -> Draft v2 (fixes applied)
Critic  -> [No new issues]
Done.
```

### SWE-Debate (competitive)
Three-round debate among specialized agents with distinct reasoning perspectives on software problems. Competitive (not cooperative) — surfaces edge cases collaborative approaches miss.

---

## Adversarial review

### The echo-chamber problem
LLM-based self-review has systematic leniency bias. Reviewer and generator share blind spots and fail in correlated ways. This is the core motivation for adversarial review.

### Builder-Critic implementation
1. Builder Agent generates code in Session A
2. Output passed to Critic Agent in Session B (different model, different prompt)
3. Critic has an explicit **adversarial kill mandate** — its job is to find problems, not approve
4. **Context asymmetry**: critic receives spec + output but NOT the builder's chain-of-thought
5. Cross-model critics (Claude critiques GPT) share fewer blind spots

### Stage-gated adversarial review (arXiv 2604.19049)
```
Stage 1: Generation (multiple candidates)
Stage 2: Adversarial screening (kill mandate)
Stage 3: Cross-model verification (different architecture)
Stage 4: Empirical validation (test execution, tool verification)
Stage 5: Promotion (only survivors reach output)
```

### Key principles
- Heterogeneous models: 3 instances of same model → 18.3% hallucination improvement; 3 heterogeneous → 35.9%
- Enforce context asymmetry: critic must not see the generator's reasoning chain
- Explicit kill mandate: the critic's job is to break things
- Cross-domain rotation: periodically swap which model plays builder vs. critic

---

## LLM-as-Judge evaluation

| Variant | Description |
| --- | --- |
| Single-judge | One LLM evaluates another's output; prone to position bias and self-preference |
| Multi-agent panel | Multiple judges debate; closer to human panel |
| Agent-as-a-Judge | Judge has tool use, memory, multi-step reasoning; evaluates the full action chain |
| Meta-Judge (MAJ-EVAL) | Multi-dimensional weighted rubric + consensus aggregation; better alignment with human expert ratings |

**Evaluation dimensions:** correctness, completeness, reasoning quality, novelty, harmfulness, consistency.

---

## Constitutional AI in council settings

A shared constitution serves as a normative framework all agents must respect.

**Constitutional council pattern:**
```
1. Define shared constitution (principles, red lines, priorities)
2. Each agent receives the constitution in its system prompt
3. After each round, a Constitutional Auditor checks outputs
4. Violations trigger revision: offending agent must regenerate
5. Constitution is versioned and can evolve
```

**Warning:** Principles that produce ethical behavior in isolation may not scale when agents interact — strategic incentives can amplify goal conflicts.

**Automated constitutional optimization** (arXiv 2602.00755): Optimized constitutions can reduce agent communication by 98.6% while increasing productivity by 203%.

---

## Critique-revision loops

### Convergence detection
Track new issues found per round. Declare convergence when critique delta drops below threshold (no new HIGH or MEDIUM issues):
```python
round_issues = [12, 5, 2, 0]
if round_issues[-1] == 0 or delta < threshold:
    converged = True
```

### Multi-critic revision
Simultaneous critics with distinct specializations: Accuracy, Style, Safety, Completeness. Generator receives all critiques and prioritizes fixes.

### Escalating critique
- Round 1: Surface-level (typos, formatting, obvious errors)
- Round 2: Logical (argument structure, evidence quality)
- Round 3: Adversarial (actively try to break claims)
- Round 4: Expert-level (domain-specific deep analysis)

---

## Implementation architecture

### Three-phase pipeline (Council Mode)

**Phase 1 — Triage:** Lightweight classifier evaluates query complexity. Simple queries bypass the council entirely (reduces latency 30.6% without quality loss).

**Phase 2 — Parallel expert generation:** Architecturally diverse models run independently and concurrently. Total latency = slowest expert + synthesis time, not the sum.

**Phase 3 — Consensus synthesis:** Dedicated synthesis model organizes outputs into consensus points, disagreements, unique findings, and comprehensive analysis.

### Shared ledger pattern
Agents write proposals to structured JSON files. A consensus service reads all proposals and determines the final decision.

### Context window diversity
Use models with different context windows (128K, 200K, 1M tokens). Different experts access varying historical context, contributing to complementary knowledge retrieval.

---

## Design decisions

### Pattern selection

| Scenario | Pattern |
| --- | --- |
| Quick classification | Plurality vote (fast, cheap, 3 agents) |
| High-stakes decision | Multi-round debate |
| Code generation | Builder-Critic loop |
| Safety-critical output | Constitutional council |
| Factual research | Oracle-verified consensus |
| Complex reasoning | Structured synthesis |

### Minimum viable council
Three heterogeneous agents is the practical minimum. Beyond 5–7 agents, returns diminish unless using stochastic sampling.

### Latency budget
Council Mode averages 8.4s vs 2.8–5.6s for individual models. Quality improvement (91.7% vs 71.4–81.5%) justifies the overhead for accuracy-prioritized applications only.

---

## Anti-patterns

1. **Homogeneous echo chamber**: Three instances of the same model improves hallucination only 18.3% vs 35.9% for heterogeneous models.
2. **Majority pressure suppression**: Agents under majority pressure abandon correct minority positions. Diversity drives debate success, not vote counts.
3. **Unbounded debate rounds**: Always set `max_rounds` and a convergence threshold.
4. **Judge-generator conflation**: Using the same model as both generator and judge loses the independence needed for effective review.
5. **Ignoring cost scaling**: Budget accordingly and use triage to skip the council for simple queries.
6. **Over-weighting confident agents**: Calibration varies by model and domain. Validate confidence against historical accuracy before weighted voting.
7. **Symmetric prompting**: Identical system prompts defeat the purpose of diversity. Assign distinct roles, perspectives, or domain expertise to each agent.

---

## Framework support

| Framework | Council-relevant features |
| --- | --- |
| AutoGen / AG2 | GroupChat, debate patterns, supervisor-worker, swarm |
| LangGraph | State machines, conditional routing, parallel branches |
| CrewAI | Role-based agents, delegation, sequential/parallel tasks |
| Claude SDK | Subagent spawning, tool use, managed agents |
| MALLM | Purpose-built for multi-agent collaboration |

---

## References

- [Council Mode: Mitigating Hallucination](https://arxiv.org/abs/2604.02923)
- [Evolving Interpretable Constitutions](https://arxiv.org/abs/2602.00755)
- [Institutional AI: A Governance Framework](https://arxiv.org/abs/2601.10599)
- [Refute-or-Promote: Adversarial Stage-Gated Review](https://arxiv.org/abs/2604.19049)
- [SWE-Debate: Competitive Multi-Agent Debate](https://arxiv.org/abs/2507.23348)
- [Multi-Agent-as-Judge (MAJ-EVAL)](https://arxiv.org/abs/2507.21028)
- [Can LLM Agents Really Debate?](https://arxiv.org/abs/2511.07784)
- [Multi-Agent Collaboration Mechanisms Survey](https://arxiv.org/abs/2501.06322)
