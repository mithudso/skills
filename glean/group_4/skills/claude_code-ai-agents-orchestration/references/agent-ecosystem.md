<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Formerly the standalone `agent-ecosystem` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: agent-ecosystem
version: 1.1.0
updated: 2026-05-29
description: >
  AI agent ecosystem expert — frameworks (Claude SDK, OpenAI, LangGraph, CrewAI,
  ADK, Mastra, AutoGen/AG2, Agno, Strands), orchestration patterns (supervisor,
  swarm, pipeline, DAG, hierarchical, consensus), lifecycle management,
  MCP/A2A protocols, memory architectures, cost optimization, distribution,
  TAM agent infrastructure, and security. TRIGGER: user asks which agent
  framework to use, wants to design multi-agent orchestration, needs to compare
  frameworks, or asks about agent security/cost/observability.
  SKIP: A2A protocol spec details (use a2a-interop); AI language/framework
  selection without agent orchestration (use ai-languages).
origin: local
updated: 2026-05-29
tags: [agents, frameworks, orchestration, MCP, A2A, multi-agent, production, security, memory, cost-optimization]
related_skills: [a2a-interop, agent-plan-writing, agent-council, ai-languages, mcp-servers]
---

# Agent Ecosystem Expert

Comprehensive reference for AI agent development, orchestration, infrastructure, and security. Deep context at `references/agent-ecosystem-context.md`.

## When to use this skill

Activate when the user:

- asks which agent framework to use
- needs to design multi-agent orchestration (supervisor, swarm, pipeline, DAG)
- wants agent lifecycle management (spawning, monitoring, stopping, error recovery)
- asks about MCP, A2A, or agent communication protocols
- needs agent memory architecture guidance
- wants to package, distribute, or deploy agents
- asks about agent security (sandboxing, permissions, HITL, prompt injection, audit trails)
- is building TAM-focused agent workflows (account monitoring, case triage, meeting prep)
- needs to compare frameworks or select an orchestration pattern
- asks about agent observability, cost control, or resource budgets
- needs production deployment guidance, failure analysis, or reliability engineering

## When NOT to use this skill

- A2A protocol spec, Agent Cards, 8-state task lifecycle details — use `a2a-interop`
- Language/framework selection without agent orchestration context — use `ai-languages`
- Multi-agent council/consensus patterns (debate, adversarial review) — use `agent-council`
- Writing an agent execution plan — use `agent-plan-writing`

## Quick framework selection

| Scenario | Framework | Why |
| --- | --- | --- |
| TypeScript teams | Mastra | TS-native, built-in memory, workflow engine |
| Complex stateful workflows | LangGraph | Graph-based state machines, checkpointing, LangSmith observability |
| Fastest prototyping | CrewAI or Agno | Role-based DSL (~20 lines to start) |
| Enterprise .NET/Azure | Semantic Kernel | Microsoft-backed, Azure-native |
| OpenAI-native with safety rails | OpenAI Agents SDK | Lightweight, strong voice support, handoff model |
| Google Cloud / multi-framework | Google ADK | Python/TS/Java/Go SDKs, A2A Agent Cards, Vertex AI |
| Deep MCP integration / OS access | Claude Agent SDK | 37 pre-built tools, file/shell access, in-process MCP servers |
| Open-ended/deliberative patterns | AG2 (AutoGen successor) | Event-driven, universal framework interop via AgentOS |
| Minimal overhead, fastest instantiation | Agno | High-throughput swarms, RBAC, tracing, scheduling built in |
| AWS-native, model-driven simplicity | Strands Agents SDK | LLM controls the agent loop, minimal abstraction |

## Framework tiers (2026)

### Tier 1 — Production proven at scale

**LangGraph** — Default for stateful production workflows. Models agents as nodes in a directed cyclic graph, enabling loops, retries, and iterative reasoning. Verified deployments: Klarna, Uber, LinkedIn, BlackRock, Cisco, JPMorgan, Replit. Ecosystem: LangSmith (observability), LangServe (deployment), LangGraph Cloud (managed hosting).

**Claude Agent SDK** — Deepest OS access. Treats MCP as infrastructure: in-process, remote, or standard MCP servers plug in without custom integration. 37 pre-built tools. Managed Agents added Q2 2026. Best for "give the agent a computer" use cases.

**OpenAI Agents SDK** — Lightweight with strong voice support. Handoff model maps directly to triage → specialist → escalation flows. Model-flexible with documented paths to non-OpenAI models.

### Tier 2 — Strong adoption, niche strengths

**Google ADK** — Hierarchical agent trees with native A2A support. Python/TS/Java/Go SDKs, 200+ models via Model Garden, Vertex AI Agent Engine. Best for enterprise systems on Google Cloud or multi-language teams.

**CrewAI** — Lowest learning curve with role-based DSL. Teams often prototype with CrewAI then migrate to LangGraph for production-grade state management. Best for hackathons, MVPs, proof-of-concept.

**AG2 (AutoGen successor)** — 2026 redesign: streaming, event-driven, multi-provider, dependency injection, typed tools, first-class testing. AgentOS provides universal framework interoperability (AG2, ADK, OpenAI, LangChain in one team).

### Tier 3 — Emerging / specialized

**Mastra** — One of three frameworks with genuine built-in memory. TS-native with workflow engine.

**Strands Agents SDK (AWS)** — Open-sourced May 2025, 14M+ downloads. Model-driven: the LLM controls the entire agent loop. Simpler code, less deterministic execution.

**Agno** — SDK for building agent platforms with tracing, scheduling, and RBAC. Designed for high-throughput swarms.

## Orchestration patterns

| Task shape | Pattern | Key trait |
| --- | --- | --- |
| Sequential dependencies | Pipeline | Defined, repeatable; easiest to reason about |
| Quality control needed | Supervisor | Single control flow; 2026 production default |
| Parallelizable independent work | Swarm | Fully decentralized, peer-to-peer, shared task queue |
| Mixed dependencies + parallelism | DAG | Batch independent steps, sequence dependent ones |
| Large teams (50+ agents) | Hierarchical | Tree of supervisors; scales with sub-team decomposition |
| High-stakes decisions | Consensus | Multiple agents deliberate; reduces single-point-of-failure bias |
| Most production systems | Hybrid | Combine patterns; supervisor at top, pipeline within workers |

### What survived production in 2026

Peer-collaboration multi-agent systems failed production. Only 3 patterns survived at scale: agent-flow (pipeline), orchestration (supervisor), and bounded collaboration (constrained swarm with explicit coordination rules).

### Supervisor (the 2026 production default)

Easiest to debug: single control flow to trace. Scales horizontally by adding workers. Maps cleanly to audit trails and rollback points. Claude Code subagents, LangGraph Supervisor, and OpenAI Agents SDK handoffs all implement this pattern.

### Swarm

Fully decentralized: agents communicate peer-to-peer, claim tasks from a shared queue. Kimi K2.6 runs 300 sub-agents executing 4,000 coordinated steps. Best for embarrassingly parallel workloads where coordination overhead must be minimized.

## Production failure analysis

### MAST study — 14 failure modes across 1,600 traces

- **41–87% failure rates** across every framework tested
- **79% of failures are specification/coordination issues**, not technical bugs
- **Top 3 failure modes**: step repetition, reasoning-action mismatch, unaware of termination conditions
- **Coordination breakdowns**: 36.9% of all failures
- Many failures are structural — not fixable with better prompts

### Production survival rules

1. Build a single-agent baseline first — most teams adopt multi-agent before single-agent reaches its ceiling
2. Prefer supervisor over peer-collaboration in production
3. Add explicit termination conditions to every agent
4. Use structured output (not free-form text) for inter-agent communication
5. Set per-agent token budgets and circuit breakers
6. Log every agent turn with full context for debugging

## Protocols — MCP and A2A

### MCP (Model Context Protocol) — agent-to-tool connectivity

- **97M+ monthly SDK downloads** (Python + TypeScript) as of April 2026
- **78% of enterprise AI teams** have ≥1 MCP-backed agent in production
- **9,400+ public MCP servers**, +18% month-over-month
- JSON-RPC 2.0, three primitives: Tools, Resources, Prompts

### A2A (Agent-to-Agent) — agent discovery and delegation

- **23% adoption** among competing protocols (growing)
- **150+ partner organizations** including Atlassian, Salesforce, SAP, ServiceNow
- Agent Cards at `/.well-known/agent.json` for cross-team discovery
- 8-state task lifecycle over JSON-RPC 2.0

**Combined strategy:** MCP for each agent's tool access (vertical); A2A for inter-agent coordination (horizontal). Organizations using both achieve 40–60% faster workflow development. 42% of CTOs expect to combine both in production.

## Agent memory architectures (2026)

| Need | Choose |
| --- | --- |
| General-purpose, any framework | Mem0 (48K+ stars, $24M Series A, hybrid vector+graph+KV) |
| Long-running production agents | Letta (OS-inspired: core memory as RAM, archival as disk) |
| Temporal reasoning / "when did X happen" | Zep (Graphiti): temporal KG with validity windows |
| Already using LangGraph | LangMem |
| Built-in, no external dependency | Mastra |

## Cost optimization

LLM API calls account for 70–85% of total agent operating costs. Most teams overpay by 40–85% by defaulting to frontier models everywhere. Unoptimized agents cost $10–$100+ per session; agents burn 50x more tokens than chat.

**Key strategies:**

1. **Model routing (saves 40–75%):** Route each step to the cheapest model that meets quality. Example: routing 4 steps to mid-tier, 2 to small, kept only final step on frontier → cost reduced from $1.40 to $0.34/document.
2. **Prompt caching (saves 45–80%):** Reduces costs 45–80%; improves time-to-first-token 13–31%.
3. **Context management (saves ~72%):** Switching from naive full-context injection to retrieval-based memory cuts 594 tokens to 166 per call with identical quality.
4. **Token budgets:** Per-request `max_tokens`, per-task budgets, per-day/month caps with alerts at 50% and 80%.

By Q1 2027, the standard metric shifts from "cost per token" to **"cost per successful task"**.

## Agent security (2026)

### Prompt injection — top open risk

Every stream entering an agent (web pages, PDFs, support tickets, emails, tool outputs) is a potential instruction channel.

### Critical: RCE via prompt injection

Microsoft discovered a path in Semantic Kernel where a single prompt could turn prompt injection into host-level RCE (remote code execution). Agent security failures can escalate far beyond prompt-level misbehavior.

### Defense-in-depth

1. **Input validation** — Filter and sanitize all external content before it enters agent context
2. **Sandboxed tool execution** — MicroVMs, gVisor, or container isolation for code execution
3. **Context-layer governance** — Least privilege; limit ungoverned access and unverified context
4. **Runtime guardrails** — LlamaFirewall (Meta, open-source) for real-time guardrails

### Security checklist

- [ ] All tool calls sandboxed (no direct shell access without isolation)
- [ ] Per-agent permission scoping (least privilege)
- [ ] Human-in-the-loop gates for destructive actions
- [ ] Token budget limits to prevent infinite loops
- [ ] Audit trail for every agent action
- [ ] Input sanitization on all external data sources
- [ ] Output validation against expected schemas

## Vendor SDK comparison

| Feature | Claude Agent SDK | OpenAI Agents SDK | Google ADK | LangGraph | CrewAI |
| --- | --- | --- | --- | --- | --- |
| Language | Python | Python | Py/TS/Java/Go | Python/JS | Python |
| Model lock-in | Claude only | OpenAI default, swappable | 200+ via Model Garden | Any LLM | Any LLM |
| MCP support | Native (infrastructure) | Supported | Supported | Supported | Supported |
| A2A support | Via integration | Via integration | Native | Via integration | Via integration |
| Built-in tools | 37 pre-built | Moderate | Moderate | Via LangChain | Role-based |
| Memory | Via MCP | Custom | Via Vertex | LangMem / custom | Basic |
| Managed hosting | Managed Agents (Q2 2026) | API-native | Vertex AI Agent Engine | LangGraph Cloud | CrewAI Enterprise |
| Best for | Computer-use agents | Triage/handoff flows | Multi-language enterprise | Stateful workflows | Rapid prototyping |

## Anti-patterns

1. **Start with multi-agent.** Build a single-agent baseline first. Most teams adopt multi-agent before single-agent reaches its ceiling.
2. **Peer-collaboration in production.** Failed in 2026. Use supervisor or bounded collaboration.
3. **Same model for every agent step.** Overpays by 40–85%. Route each step to the cheapest model that meets quality.
4. **Skip termination conditions.** "Unaware of termination conditions" is a top-3 failure mode. Every agent must have an explicit exit criterion.
5. **Free-form text between agents.** Use structured output (JSON schemas) to prevent reasoning-action mismatches.
6. **No token budgets.** Agents burn 50x more tokens than chat. Set per-request, per-task, and per-day caps.
7. **Trust external content.** Every text stream entering an agent is a potential prompt injection vector. Sanitize all inputs.

## Guidance

- 41–87% of multi-agent LLM systems fail in production; 79% of failures are specification/coordination issues, not bugs.
- For TAM workflows: orchestrator-specialist pattern dominates. Central orchestrator dispatches to case, metrics, document, and communication agents in parallel.
- The boundaries between frameworks are blurring — future agent systems will be multi-framework compositions.
- Always implement token budgets, circuit breakers, and audit trails before going to production.
- Measure cost per successful task, not cost per token.

## Cited sources

1. [Best Multi-Agent Frameworks in 2026 — GuruSup](https://gurusup.com/blog/best-multi-agent-frameworks-2026)
2. [Claude Agent SDK vs OpenAI Agents SDK vs Google ADK — Composio](https://composio.dev/content/claude-agents-sdk-vs-openai-agents-sdk-vs-google-adk)
3. [Multi-Agent Orchestration: 5 Patterns That Work](https://www.digitalapplied.com/blog/multi-agent-orchestration-5-patterns-that-work)
4. [Multi-Agent in Production 2026: 3 Patterns That Survived](https://niteagent.com/blog/multi-agent-production-2026/)
5. [MCP Adoption Statistics 2026](https://www.digitalapplied.com/blog/mcp-adoption-statistics-2026-model-context-protocol)
6. [AI Agent Cost Optimization in 2026](https://niteagent.com/blog/ai-agent-cost-optimization-2026/)
7. [When Prompts Become Shells: RCE Vulnerabilities](https://www.microsoft.com/en-us/security/blog/2026/05/07/prompts-become-shells-rce-vulnerabilities-ai-agent-frameworks/)
8. [State of AI Agent Memory 2026 — Mem0](https://mem0.ai/blog/state-of-ai-agent-memory-2026)
