---
name: ai-rag-retrieval
description: >-
  RAG & retrieval sub-hub (ai-agent family). TRIGGER: Retrieval-Augmented Generation architecture (chunking, embeddings, hybrid search, reranking, query transformation, evaluation); advanced RAG patterns (self-RAG, corrective, GraphRAG, agentic, adaptive); iterative/multi-step retrieval; AI datastores (vector & graph databases). SKIP: agent orchestration → ai-agents-orchestration; model training/serving → ai-llm-model-layer; MCP/SDK/prompting → ai-mcp-sdk-prompting.
origin: local
---
# ai-rag-retrieval

RAG & retrieval sub-hub (ai-agent family).

Hub routes to on-demand reference files under `references/`. See each spoke for depth.

<!-- cross-hub-map -->
## Cross-hub map — where every ai-agent topic lives

Family split across these hubs. Deep material not in this hub's sub-skill routing table lives as reference file under sibling hub below — activate that hub or `Read` its `references/<name>.md` directly. Every former standalone skill now reference under one of these hubs (nothing deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `ai-agents-orchestration` | ai-agents-orchestration | `references/a2a-interop.md`, `references/agent-council.md`, `references/agent-ecosystem.md`, `references/agent-harness-construction.md`, … |
| `ai-rag-retrieval` | ai-rag-retrieval | `references/advanced-rag-patterns.md`, `references/rag-architecture.md`, `references/iterative-retrieval.md`, `references/ai-datastores.md` |
| `ai-llm-model-layer` | ai-llm-model-layer | `references/agentic-rl.md`, `references/distributed-training.md`, `references/llm-alignment-post-training.md`, `references/llm-compression.md`, … |
| `ai-mcp-sdk-prompting` | ai-mcp-sdk-prompting | `references/mcp-builder.md`, `references/mcp-servers.md`, `references/anthropic-sdk.md`, `references/prompt-engineering.md`, … |