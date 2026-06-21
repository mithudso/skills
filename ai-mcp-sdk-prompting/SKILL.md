---
name: ai-mcp-sdk-prompting
description: >-
  MCP, SDK, prompting & context-engineering sub-hub (ai-agent family). TRIGGER: MCP servers & the MCP builder; Anthropic SDK; prompt engineering; LLM context engineering (CLAUDE.md, caching, agent memory, prompt-injection hardening); declarative LLM frameworks; multi-provider LLM integration review; AI programming languages; AI red-teaming tooling; MCP tool-search/discovery optimization; prompt discovery/lookup & registries; iterative self-refinement & convergence loops (Self-Refine, Reflexion, CRITIC, stop-condition/oscillation/budget); prompt & context compression (LLMLingua, gist tokens, soft-prompt distillation); applied vision-model UI-critique technique (VLM critiques a screenshot — structured findings, Set-of-Mark grounding, before/after, VLM-as-judge calibration & failure modes). SKIP: agent orchestration → ai-agents-orchestration; RAG → ai-rag-retrieval; model training/serving (incl. VLM encoder/architecture internals) → ai-llm-model-layer.
origin: local
---
# ai-mcp-sdk-prompting

MCP/SDK/prompting/context-engineering sub-hub (ai-agent family).

Routes to on-demand refs under `references/`. See each spoke for depth.

<!-- cross-hub-map -->
## Cross-hub map — where every ai-agent topic lives

Family split across hubs. Deep material **not** in Sub-skill routing table → ref file in sibling hub — **activate that hub or `Read` its `references/<name>.md` directly**. All former standalone skills now refs under one of these hubs (nothing deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `ai-agents-orchestration` | ai-agents-orchestration | `references/a2a-interop.md`, `references/agent-council.md`, `references/agent-ecosystem.md`, `references/agent-harness-construction.md`, … |
| `ai-rag-retrieval` | ai-rag-retrieval | `references/advanced-rag-patterns.md`, `references/rag-architecture.md`, `references/iterative-retrieval.md`, `references/ai-datastores.md` |
| `ai-llm-model-layer` | ai-llm-model-layer | `references/agentic-rl.md`, `references/distributed-training.md`, `references/llm-alignment-post-training.md`, `references/llm-compression.md`, … |
| `ai-mcp-sdk-prompting` | ai-mcp-sdk-prompting | `references/mcp-builder.md`, `references/mcp-servers.md`, `references/anthropic-sdk.md`, `references/prompt-engineering.md`, `references/vision-model-design-critique.md`, … |