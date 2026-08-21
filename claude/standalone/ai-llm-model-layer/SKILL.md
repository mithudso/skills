---
name: ai-llm-model-layer
description: >-
  LLM model-layer sub-hub (ai-agent family) — training, alignment, serving, architecture. TRIGGER: pretraining & scaling laws; distributed training; fine-tuning/PEFT; alignment & post-training (SFT, RLHF/PPO, DPO family); RLHF infrastructure; agentic RL; LLM compression (quantization/distillation/pruning); GPU kernels; inference serving/runtime (PagedAttention, batching, speculative decoding, TTFT/TPOT); transformer & multimodal architecture; reasoning models; model landscape/selection; LLM observability; continuous learning; neural architecture search (NAS — DARTS/ENAS/weight-sharing supernets, zero-cost proxies, hardware-aware NAS). SKIP: agent orchestration → ai-agents-orchestration; RAG → ai-rag-retrieval; MCP/SDK/prompting → ai-mcp-sdk-prompting.
origin: local
---
# ai-llm-model-layer

LLM model-layer sub-hub (ai-agent family) — training, alignment, serving, architecture.

Hub routes to on-demand reference files under `references/`. See each spoke for depth.

<!-- cross-hub-map -->
## Cross-hub map — where every ai-agent topic lives

Family split across hubs. Deep material not in this hub's sub-skill routing table lives as reference file under sibling hub below — activate that hub or `Read` its `references/<name>.md` directly. Every former standalone skill now reference under one hub (nothing deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `ai-agents-orchestration` | ai-agents-orchestration | `references/a2a-interop.md`, `references/agent-council.md`, `references/agent-ecosystem.md`, `references/agent-harness-construction.md`, … |
| `ai-rag-retrieval` | ai-rag-retrieval | `references/advanced-rag-patterns.md`, `references/rag-architecture.md`, `references/iterative-retrieval.md`, `references/ai-datastores.md` |
| `ai-llm-model-layer` | ai-llm-model-layer | `references/agentic-rl.md`, `references/distributed-training.md`, `references/llm-alignment-post-training.md`, `references/llm-compression.md`, … |
| `ai-mcp-sdk-prompting` | ai-mcp-sdk-prompting | `references/mcp-builder.md`, `references/mcp-servers.md`, `references/anthropic-sdk.md`, `references/prompt-engineering.md`, … |