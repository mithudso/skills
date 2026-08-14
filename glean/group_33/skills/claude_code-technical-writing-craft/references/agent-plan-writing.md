<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `agent-plan-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: agent-plan-writing
version: 1.1.0
updated: 2026-05-29
description: >
  Write execution plans for AI agent workflows — multi-agent orchestration,
  subagent prompt crafting, context budgeting, tool/action-space design, memory
  management, safety guardrails, evaluation, and non-code agent tasks.
  Complements code-plan-writing with agent-specific planning patterns.
  TRIGGER: user needs to plan a workflow with 2+ coordinating agents, design
  subagent prompts, budget token usage across agents, plan tool selection
  strategy, add safety guardrails or output validation to an agent plan, plan
  agent evaluation, or structure non-code agent tasks (research, analysis,
  document generation, data pipelines).
  SKIP: single-agent code implementation plan (use code-plan-writing); building
  the agent framework itself (use agent-ecosystem or mcp-servers); one-shot
  prompt with no multi-step workflow.
origin: local
tags: [agent-planning, multi-agent, orchestration, context-budget, subagent, guardrails, evaluation]
related_skills: [agent-ecosystem, agent-council, code-plan-writing, mcp-servers]
---

# Agent Plan Writing

## Overview

Agent plan writing is the discipline of designing execution plans for AI agent workflows that go beyond single-agent code implementation. While `code-plan-writing` covers translating specs into implementable coding tasks, this skill covers agent-specific planning: decomposing work across agents, crafting subagent prompts, budgeting context windows, scoping tools, managing state, adding safety guardrails, planning evaluation, and structuring non-code tasks like research, analysis, and document generation.

**Core insight:** The harness matters more than the model. Agent completion rates depend more on action-space design, context engineering, and orchestration patterns than on which frontier model is used. Planning these elements deliberately — rather than leaving them to runtime discovery — separates production agents from prototypes.

## Output format

When this skill activates, produce a markdown agent plan containing:

1. **Workflow Overview** — what the system does, which orchestration pattern, and why
2. **Agent Roster** — each agent's role, model, tools, and context scope
3. **Orchestration Graph** — how agents coordinate (sequential, parallel, supervised, etc.)
4. **Context Budget** — token allocation per agent, what each receives and doesn't
5. **Safety Constraints** — permission boundaries, output validation, human-in-the-loop gates
6. **Evaluation Plan** — what to trace, quality metrics, where to set gates
7. **Failure Handling** — per-pattern failure modes and recovery strategies

**Condensed example (supervisor pattern):**
```markdown
# Customer Research Report — Agent Plan

## Workflow Overview
Supervisor pattern: coordinator dispatches to 3 specialist agents, then synthesizes.

## Agent Roster
| Agent | Role | Model | Tools | Context Scope |
|-------|------|-------|-------|--------------|
| coordinator | Decompose query, synthesize | opus | None | Full user query |
| web-researcher | Search and extract web sources | sonnet | WebSearch, WebFetch | Sub-question only |
| doc-analyzer | Read and summarize internal docs | haiku | FileRead | Document paths only |
| report-writer | Draft final report | sonnet | None | Aggregated findings only |

## Context Budget
coordinator: 2k system + 4k working | web-researcher: 1k system + 8k results
doc-analyzer: 500 system + 6k docs | report-writer: 1k system + 10k input

## Safety Constraints
- web-researcher: no file writes, no code execution
- report-writer: output max 3000 words, must cite sources

## Evaluation Plan
Trace: tool calls per agent, total tokens, source count. Gate: report must cite ≥3 sources.

## Failure Handling
web-researcher timeout → retry once, then return partial results to coordinator
```

## Plan completeness checklist

Before delivering an agent plan, verify:

- [ ] Every agent has a defined role, model, tool list, and context scope
- [ ] Orchestration pattern named and justified
- [ ] Token budget allocated per agent (system prompt + working memory + output headroom)
- [ ] Each agent's tools scoped to its task (no symmetric tool loading)
- [ ] Output contracts between agents specify exact format, required fields, and max length
- [ ] Safety constraints defined at tool level (what each agent cannot do)
- [ ] At least one evaluation metric is measurable and automated
- [ ] Failure handling covers the primary failure mode for the chosen pattern
- [ ] Handoff payloads defined (what's included, what's excluded)

---

## Orchestration patterns

Five patterns dominate production agent systems. Choose based on task structure, not framework preference.

### Fan-Out
Parallel execution of independent subtasks. Coordinator dispatches to N agents simultaneously, then aggregates.

**When to use:** Research across multiple sources, parallel file review, independent data extraction.

**Plan design:** Identify genuinely independent subtasks. Design aggregation logic upfront — decide how to handle partial failures before execution.

**Failure mode:** Silent partial results. An agent returns empty data instead of erroring, and the aggregator treats it as "nothing found." Plan explicit empty-result detection.

### Pipeline
Sequential chain where each stage requires the prior stage's output.

**When to use:** Any workflow with a natural dependency chain where output quality compounds (research → draft → critique → revise).

**Plan design:** Map the dependency chain before implementation. Insert validation gates between stages to prevent cascade contamination.

**Failure mode:** Cascade failure — a bad mid-stage output poisons every downstream stage. Validate aggressively at each gate.

**Contract-first principle:** Before building any multi-agent pipeline, define the output contract between every agent pair: exact fields, types, and validation rules. This prevents the most common multi-agent failure — agents producing outputs the next agent can't parse.

### Supervisor
A supervisor agent decomposes the task, delegates to specialists, and synthesizes results. The 2026 production default.

**When to use:** Cross-domain tasks (coder + researcher + reviewer), tasks requiring heterogeneous expertise.

**Plan design:** Assign the supervisor a frontier model (needs reasoning to route correctly); assign workers cheaper models (need tool-calling reliability). Set iteration ceilings per subagent (Claude Agent SDK default: ~25 turns).

**Failure mode:** Over-delegation loops. Subagents return partial results; supervisors re-delegate instead of synthesizing. Treat iteration ceilings as architectural boundaries.

### Debate
Send identical prompts to 2+ diverse models, route responses to a judge for arbitration.

**When to use:** High-stakes decisions where divergent perspectives are genuinely valuable. Reserve for externally visible decisions.

**Plan design:** Confirm stakes justify the ~2.5x cost multiplier. Select models with genuinely different perspectives. Set hard maximum debate rounds.

**Failure mode:** Judge bias favoring style over accuracy. Plan explicit judging rubrics.

### Swarm
Dynamic spawning of agents based on workload, with shared-memory coordination.

**When to use:** Only when you have 50+ parallel independent tasks. Significant engineering overhead.

**Plan design:** Confirm genuine scale requirement (most teams should use supervisor + fan-out instead). Set hard population-size caps.

**Failure mode:** Runaway agent spawning. Always set population caps.

**Decision tree:** Start with supervisor (widest native support). Add fan-out branches when subtasks are independent. Use debate only if stakes justify 2.5x cost. Deploy swarm only at 50+ task scale.

---

## Subagent prompt crafting

Each subagent starts with a clean context window pre-loaded with its own system prompt and only its assigned tools. The parent agent only sees the final result.

**Construction principles:**
- **Role declaration:** State what the agent is ("You are a security reviewer") not what the system is
- **Context scoping:** Include only what this agent needs — exclude parent conversation, sibling outputs, and unrelated tool schemas
- **Output contract:** Specify exact format, required fields, and maximum length
- **Constraint injection:** State what the agent must NOT do (prevents scope creep more effectively than positive instructions alone)
- **Tool restrictions:** List only the tools this agent should use, even if more are available

**Example:**
```
Agent: research-analyst
Model: sonnet (cost-effective for retrieval tasks)
Tools: WebSearch, WebFetch (no file editing tools)
Prompt: "You are a research analyst. Search for {topic}. Return exactly:
  - 5 key findings, each with a cited source URL
  - 1 paragraph synthesis
  Do NOT make recommendations. Do NOT write code. Return findings only."
```

**Context firewall principle:** Each subagent gets only the context it needs — never the parent's full conversation. This prevents context pollution, reduces token cost, and enables independent verification.

---

## Context window budget planning

Agents consume ~7x more tokens than standard chat sessions. Plan token budgets explicitly.

**Budget allocation:**
1. **System prompt:** 500–2,000 tokens. Cached input costs 10–25% of normal.
2. **Tool schemas:** Each MCP tool adds 100–500 tokens to context. Load only tools the agent will actually use.
3. **Working memory:** Reserve 30–50% of context for conversation/reasoning accumulation.
4. **Output headroom:** Reserve 15–25% for the agent's response generation.

**Cost optimization:**
- Route simple subagent work to cheaper models (Haiku for retrieval, Sonnet for analysis, Opus for reasoning/routing)
- Enable prompt caching (saves $2,000+/day for 5,000 agent loops)
- Use deferred tool loading — don't load all tool schemas upfront
- Compact conversation history aggressively between phases

**Warning:** Opus 4.7's tokenizer adds up to 35% more tokens than Opus 4.6 for code-heavy inputs. Validate budgets when upgrading models.

**Model-tier pattern:** Assign by cognitive demand — frontier (Opus) for supervisors/judges needing reasoning, mid-tier (Sonnet) for analysis/synthesis, fast/cheap (Haiku) for retrieval and classification. Cuts costs 3–5x vs. using frontier everywhere.

---

## Tool and action-space design

In 2026, the MCP ecosystem offers ~177,000 distinct tools. Planned tool selection is a critical design decision.

**Selection strategy:**
1. **Minimal viable toolset:** Give each agent only the tools it needs. Extra tools increase context cost and decision complexity.
2. **Tool priority ordering:** When multiple tools could accomplish a task, specify preferred order (e.g., "Use deep-research first, fall back to WebSearch, last resort WebFetch").
3. **Fallback chains:** For every critical tool, plan what happens if it fails (timeout, error, rate limit).

**Action-space sizing rule:** Keep each agent's action space under 20 tools. Beyond that, tool selection accuracy degrades measurably.

---

## Agent memory and state management

| Tier | Scope | Persistence | Use Case |
| --- | --- | --- | --- |
| Working memory | Single agent turn | In-context | Reasoning, scratch work |
| Thread memory | Single task/conversation | Checkpointed | Multi-step workflows, HITL |
| Cross-session memory | Across tasks | Database-backed | User preferences, accumulated knowledge |

**Checkpointing strategy:**
- Checkpoint after every agent node execution, not just at phase boundaries
- Use thread IDs to keep different users/tasks separate
- Never overwrite previous state — maintain checkpoint history for rollback
- For human-in-the-loop: compile graph with `interrupt_before` at approval nodes

**Handoff payload design:** Include task description, relevant findings (not raw data), constraints, and expected output format. Exclude conversation history, tool call logs, and reasoning traces.

---

## Safety guardrails in agent plans

| Layer | What it catches | Implementation |
| --- | --- | --- |
| Model-level | Content policy violations | Built into the LLM (no plan action needed) |
| Application-level | Domain errors, hallucination | Output validators, LLM-as-judge scoring |
| Tool-level | Unauthorized actions | Permission boundaries per agent, tool allowlists |
| Human oversight | Judgment calls, high-stakes decisions | Interrupt gates, approval workflows |

**Planning guardrails:**
1. **Permission boundaries:** Write-capable tools (file edit, API calls, database writes) require explicit scoping.
2. **Output validation:** For every agent output that feeds another agent or reaches a user, specify validation criteria.
3. **Escalation paths:** Define when an agent should stop and ask for human input vs. continue autonomously.
4. **Cost caps:** Set per-agent and per-workflow token budgets with hard stops.

**Regulatory note:** EU AI Act high-risk obligations take effect August 2, 2026. Plan compliance requirements (logging, human oversight, risk documentation) now.

---

## Agent evaluation planning

Build evaluation into the plan, not after deployment.

**Three evaluation layers:**
1. **Unit evals** — Test individual agent steps (tool selection accuracy, output format compliance, constraint adherence). Run in CI.
2. **LLM-as-judge regression suites** — Score subjective output quality against golden examples. Catch quality drift.
3. **Production trace sampling** — Continuously score live agent behavior. Convert failures into eval cases.

**Minimum viable observability (what to trace):**
- Tool-call spans: tool name, arguments, output, duration, retry count, error state
- Reasoning spans: model plans, selected actions, decision branches
- State transition spans: working memory before/after each step

**Staged rollout:**
- Day 1: Basic trace capture (LLM calls, tool invocations, errors)
- Week 1: Reasoning steps, cost tracking, failure-to-eval conversion
- Month 1: Online scoring on trace samples with quality regression alerts

**Eval-from-failures principle:** When a production agent fails, convert the trace into an eval test case. The suite grows from real failures, catching regressions automatically.

**Tools:** LangSmith (best inside LangGraph), Braintrust (eval-science focus, free tier: 1M spans), Langfuse (open-source), MLflow (model-centric tracking).

---

## Non-code agent task planning

### Research agent plans
1. Define the research question and scope boundaries
2. Specify source priorities (academic > official docs > blogs > forums)
3. Set minimum source count per sub-question (3+)
4. Plan synthesis format (executive summary → themed sections → sources)
5. Add saturation signal: stop when 2 consecutive searches yield nothing new

### Document generation agent plans
1. Define the document type and audience
2. Specify structure template upfront (don't let the agent choose)
3. Plan review gates: outline review → draft review → final review
4. Add quality criteria: factual accuracy, tone consistency, format compliance

### Data pipeline agent plans
1. Map data flow: source → transform → validate → destination
2. Assign specialist agents per stage (extraction, quality, enrichment)
3. Plan data contract validation between stages
4. Add anomaly detection gates with automatic rollback triggers

**Common failure:** Treating the agent like a human assistant ("research this and write it up"). Instead, specify: what sources to check, what format to return, what counts as "done," and what to do when stuck.

---

## Methodology: writing a plan from scratch

1. **Choose orchestration pattern** — based on task structure (see decision tree above)
2. **Define agent roster** — each agent: role, model, tools, context scope
3. **Design context budgets** — token allocation per agent, what's included/excluded
4. **Map coordination flow** — how agents communicate, handoff payloads, synchronization points
5. **Add safety constraints** — permission boundaries, output validation, escalation paths, cost caps
6. **Plan evaluation** — what to trace, quality metrics, where to set gates
7. **Plan failure handling** — per-pattern failure modes with recovery strategies
8. **Document non-obvious decisions** — why this pattern, why this model assignment, why these tools

---

## Anti-patterns

**Orchestrating when you should prompt:** Multi-agent adds latency, cost, and failure modes. Start with one agent; split only when a single context window can't hold the task.

**Symmetric tool loading:** Giving every agent access to every tool wastes context tokens on unused schemas and increases tool-selection errors.

**Implicit handoffs:** "The next agent will figure it out" is the agent-planning equivalent of undocumented APIs. Define handoff payloads explicitly.

**Evaluating after deployment:** If you can't define "good output" before building, you can't build a good agent.

**Unbudgeted context:** Agents silently hit context limits, outputs degrade, and nobody notices until production.

---

## Troubleshooting

| Symptom | Fix |
| --- | --- |
| Inconsistent outputs across runs | Add output contracts with validation; retry with tighter constraints on failure |
| Workflow takes too long | Check if pipeline should be fan-out; parallelize independent stages; downgrade model tiers |
| Subagent ignores its constraints | Move constraints to top of prompt; use "Do NOT" phrasing |
| Token costs unexpectedly high | Audit tool schema loading (often the largest hidden cost); enable prompt caching |
| Agent fails silently instead of erroring | Add explicit error reporting to every agent's output contract |
| Evaluation catches nothing | Your evals are testing format, not quality; add LLM-as-judge scoring on semantic correctness |

---

## References

1. [Multi-Agent Orchestration: 5 Patterns That Work](https://www.digitalapplied.com/blog/multi-agent-orchestration-5-patterns-that-work)
2. [Agent Observability: Complete Guide 2026](https://www.braintrust.dev/articles/agent-observability-complete-guide-2026)
3. [Context Engineering: Why More Tokens Makes Agents Worse](https://www.morphllm.com/context-engineering)
4. [Subagents — Agentic Engineering Patterns](https://simonwillison.net/guides/agentic-engineering-patterns/subagents/)
5. [AI Agent Guardrails & Output Validation 2026](https://toolhalla.ai/blog/ai-agent-guardrails-io-validation-2026)
6. [LangGraph Persistence Guide](https://fast.io/resources/langgraph-persistence/)
7. [Tools in 2026: The New Planning Problem](https://medium.com/@Micheal-Lanham/tools-in-2026-why-picking-the-right-action-is-the-new-planning-problem-d28d8443bf3f)
8. [AI Agents Burn 50x More Tokens](https://leanopstech.com/blog/agentic-ai-cost-runaway-token-budget-2026/)
