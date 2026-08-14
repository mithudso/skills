# ai-agents-orchestration

**Category:** Cloud, DevOps & Infrastructure
**Platform:** Claude
**Original Path:** claude/standalone/ai-agents-orchestration

## Description
AI agent building & orchestration sub-hub (ai-agent family). TRIGGER: agent ecosystems/frameworks (Claude SDK, LangGraph, CrewAI, OpenAI Agents); multi-agent orchestration & councils; A2A interop / Agent Cards; agent harness construction (action spaces, tool defs, observation format); agent memory architecture; planning patterns; reliability & guardrails; durable execution/state; agent workflow builders; autonomous loops; coding agents; vibe coding & AI-assisted dev (spec-driven development); computer-use/GUI agents; voice/realtime agents; eval-driven development; automated agent/compound-system optimization from a metric (ADAS, AFlow, DSPy/MIPROv2, GEPA, TextGrad, Trace/OptoPrime, Darwin-Gödel Machine, self-improving agents); algorithmic discovery & superoptimization (AlphaTensor/AlphaDev/FunSearch/AlphaEvolve, program synthesis, evolving faster algorithms); self-improving skill libraries (Voyager, ACE, experience-to-skill distillation, playbook curation). SKIP: RAG/retrieval → ai-rag-retrieval; model training/serving → ai-llm-model-layer; MCP/SDK/prompting → ai-mcp-sdk-prompting.

---

# ai-agents-orchestration

AI agent building & orchestration sub-hub (ai-agent family).

Hub routes to on-demand reference files under `references/`. See each spoke for depth.

<!-- cross-hub-map -->
## Cross-hub map — where every ai-agent topic lives

Family split across hubs. If task's deep material **not** in this hub's Sub-skill routing table, it's reference file under sibling hub below — **activate that hub or `Read` its `references/<name>.md` directly**. Every former standalone skill in family now reference under one of these hubs (nothing deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `ai-agents-orchestration` | ai-agents-orchestration | `references/a2a-interop.md`, `references/agent-council.md`, `references/agent-ecosystem.md`, `references/agent-harness-construction.md`, … |
| `ai-rag-retrieval` | ai-rag-retrieval | `references/advanced-rag-patterns.md`, `references/rag-architecture.md`, `references/iterative-retrieval.md`, `references/ai-datastores.md` |
| `ai-llm-model-layer` | ai-llm-model-layer | `references/agentic-rl.md`, `references/distributed-training.md`, `references/llm-alignment-post-training.md`, `references/llm-compression.md`, … |
| `ai-mcp-sdk-prompting` | ai-mcp-sdk-prompting | `references/mcp-builder.md`, `references/mcp-servers.md`, `references/anthropic-sdk.md`, `references/prompt-engineering.md`, … |