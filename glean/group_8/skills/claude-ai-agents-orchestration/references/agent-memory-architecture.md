<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Formerly the standalone `agent-memory-architecture` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

<!-- FOLDED REFERENCE under the ai-agent-engineering hub — the COGNITIVE memory architecture of an agent (memory TYPES + management policy), NOT the storage engine. Sibling refs: ai-datastores (the vector/graph/doc STORAGE substrate + memory frameworks at the storage level), rag-architecture (retrieval pipelines), agent-state-and-durable-execution (execution STATE/checkpointing != memory). Any "use the X" pointer = a sibling reference under this hub; load references/<name>.md, not a bare skill. -->

# Agent Memory Architecture

What an agent **remembers across turns/sessions**, how it decides what to keep vs forget, and how it writes/retrieves it. This is the memory-**type** taxonomy and **management policy** layer. The physical store (vector DB, graph, KV) lives in `ai-datastores`; retrieval-pipeline mechanics live in `rag-architecture`; run/execution state (checkpoints, durable replay) lives in `agent-state-and-durable-execution` — **state ≠ memory** (state is the data a run reasons over and resumes from; memory is the knowledge persisted *across* runs).

## Overview — the analogy and why memory exists
An LLM is **stateless**: each call sees only its context window. "Memory" is everything bolted on to give the *illusion* of continuity. The dominant framing (CoALA, Sumers et al. 2023) borrows cognitive science:
- **Sensory ≈ raw input** (the prompt/observation as it arrives).
- **Short-term / working memory ≈ the context window** — fast, tiny, volatile.
- **Long-term memory ≈ external store** — large, slow, persistent; must be *retrieved into* the window to be used.

The whole job of a memory system: **manage the bottleneck of a finite context window** — decide what enters working memory now, what is written to long-term store, and what is dropped.

## Memory types (CoALA taxonomy — the canonical four)
Adopted by Letta, Mem0, LangChain/LangGraph as their foundation.

| Type | Holds | LLM realization | Analogy |
| --- | --- | --- | --- |
| **Working / short-term** (in-context) | Current task: recent turns, observations, intermediate reasoning, partial plans, scratchpad | The context window itself | RAM |
| **Episodic** | Past *experiences/events* — "what happened when I tried X", specific past interactions, few-shot trajectories | Stored interaction logs / event records, retrieved by similarity+recency | Autobiographical memory |
| **Semantic** | *Facts/knowledge* — user preferences, world facts, durable entities ("user is on the Enterprise plan") | Text, embeddings, or **knowledge-graph triples** | Facts you know |
| **Procedural** | *Skills / how-to* — how to perform actions, tool-use routines | Tool/function defs, code snippets, the system prompt itself, or **implicit in model weights** | Muscle memory |

Practical split most frameworks ship: **short-term** (the in-context thread) vs **long-term** (episodic + semantic + procedural in an external store). Procedural memory is often **the agent editing its own system prompt / instructions** over time.

## Short-term vs long-term management (the core policy)
Four policy questions every memory system answers:

1. **What to write (extraction).** Not every turn is worth persisting. Common patterns: write **on every turn** (cheap, noisy), or **LLM-extracted salient facts** (an extraction prompt distills "memories worth keeping" from the conversation — Mem0's add-pipeline). Hot path (write inline, blocks latency) vs **background/subconscious** write (cheaper, eventually-consistent — LangMem's hot-path vs background distinction).
2. **What to retrieve (recall).** Pull relevant long-term items back into the window. Usually embedding similarity, often fused with **recency** and **importance** (see Generative Agents below). This *is* RAG mechanically — see `rag-architecture` — but scoped to the agent's own memories, not a document corpus.
3. **Context-window pressure / compaction.** When the thread overflows the window: **truncate** (drop oldest — lossy), **summarize/compact** (recursively summarize old turns into a running summary), or **page out** to long-term store (MemGPT). Summarization is the workhorse; it trades fidelity for room.
4. **What to forget (decay).** Unbounded memory degrades retrieval precision and cost. Forgetting strategies: **TTL/recency decay** (down-weight or expire old items), **relevance pruning** (drop low-importance memories), **dedup/merge** (collapse near-duplicate facts), and **invalidation** (mark a fact stale when contradicted — critical for *updates*, e.g. user changed jobs). Most production systems under-implement forgetting; it is an open problem (see Anti-patterns).

## MemGPT / Letta — virtual context management (OS-inspired paging)
**MemGPT** (Packer et al. 2023, "LLMs as Operating Systems," arXiv:2310.08560) reframes context management as **virtual memory paging**. The OS analogy:
- **Main context** = physical RAM = the context window (fast, finite).
- **External context** = disk = long-term store (archival + recall storage).
- An **"LLM OS"** moves data between the two via **function calls the model issues itself**.

**Self-editing memory:** the LLM is given tools to **page memory in/out of its own window** and to **rewrite a persistent "core memory" block** (e.g. `core_memory_append`, `core_memory_replace`, `archival_memory_insert`, `conversation_search`). When context nears the limit, the system issues a **"memory pressure" warning**; the agent summarizes/evicts and pages relevant archival data back in. This yields **unbounded effective context** on a fixed window without retraining.

**Letta** is the open-source framework/company that grew out of MemGPT. Its production primitive is the **memory block** — a labeled, editable, persistent segment of the context (e.g. `human`, `persona`) that the agent self-curates. Letta = the standout pick when you want the **agent to autonomously manage its own memory** rather than you hand-coding write/read rules.

## Reflection & memory consolidation
**Consolidation** = turning a pile of raw episodic records into durable, higher-level semantic knowledge — the agent analog of sleep-time memory replay.

**Generative Agents** (Park et al. 2023, "Interactive Simulacra," Smallville) is the canonical recipe — **observe → reflect → plan**:
- **Memory stream:** an append-only log of natural-language **memory objects** (description, creation time, last-access time).
- **Retrieval score = α·recency + β·importance + γ·relevance.** *Recency* = exponential decay since last access; *importance* = an LLM-scored 1–10 "poignancy" rating assigned at write time; *relevance* = embedding similarity to the current query. Top-scored memories enter the window.
- **Reflection:** periodically (when summed recent importance crosses a threshold) the agent asks itself *"what high-level questions can I answer about recent events?"*, retrieves evidence, and writes **reflection nodes** — abstractions that point back to their source memories, forming a tree. Reflections are themselves retrievable, so insight compounds.

This **episodic → semantic** rollup (summarize many events into a few stable facts/insights) is the consolidation pattern every serious memory framework now imitates.

## Memory frameworks (storage-level details → `ai-datastores`)
Brief orientation; pick the substrate in `ai-datastores`.

| Framework | Model | Best when |
| --- | --- | --- |
| **Mem0** | Drop-in REST API; **LLM-extracted** facts; three scopes (**user / session / agent**); hybrid vector + graph + KV store | Managed, framework-agnostic personalization memory with minimal code |
| **Zep / Graphiti** | **Temporal knowledge graph** — every fact (edge) carries `valid_from` / `valid_to` / `invalid_at`; old facts are **invalidated, not deleted**, so the agent can reason about *how facts changed over time* | Agents that must track evolving facts and answer temporal/knowledge-update questions |
| **LangMem** | LangChain/LangGraph sub-package; **hot-path vs background** memory writes; free/OSS | Teams already on LangGraph wanting native memory hooks |
| **Letta (MemGPT)** | Self-editing memory blocks + paging (above) | Agent autonomously curates its own memory |

Reported eval gaps illustrate the design tradeoff: on **LongMemEval (GPT-4o), Zep ≈ 63.8% vs Mem0 ≈ 49.0%** — the gap concentrated in **knowledge-update / temporal-reasoning** questions, where the temporal graph wins. On **DMR**, Zep ≈ 94.8% vs MemGPT 93.4%. (Treat vendor numbers skeptically; benchmarks below are weak — see Open problems.)

## Memory vs RAG — the decision
They share retrieval *mechanics* but answer different questions. **RAG retrieves from a (mostly static) external knowledge corpus; memory persists and retrieves the agent's own evolving experience.**

| Use **memory** when… | Use **RAG** when… |
| --- | --- |
| Content is **agent/user-specific and accumulates** over interactions (preferences, past decisions, "what we tried") | Content is a **shared, authored knowledge base** (docs, wiki, code) |
| Facts **change/update** and you must track the latest (or the history) | Corpus is read-mostly and curated externally |
| You need **personalization & continuity** across sessions | You need **grounding** of answers in source-of-truth documents |
| Write path matters (extraction, consolidation, forgetting) | Ingestion is a batch pipeline (chunk → embed → index) |

In practice production agents run **both**: RAG for the knowledge base, a memory layer for the user/agent state — and the memory layer *uses* RAG's retrieval stack under the hood. Pick the storage substrate (vector vs temporal-graph vs KV) in `ai-datastores`; build the retrieval pipeline with `rag-architecture`.

## Open problems (state honestly)
- **Hallucinated / false memories.** The extraction LLM can write a fact the user never stated; once persisted it is retrieved and trusted as ground truth, **compounding** over sessions. Provenance/source-linking (Generative Agents' citation back to source memories) and confidence scoring help but don't solve it.
- **Memory mutation is under-tested.** Real value is handling **updates/contradictions** (user changed jobs), yet **LoCoMo** has ~94% of questions answerable from ≤2 sessions and **LongMemEval** caps knowledge-updates at ≤2 sessions before eval — benchmarks place **minimal stress on mutation**, the hardest part. High benchmark scores overstate real-world robustness.
- **Eval is weak and synthetic.** LoCoMo conversations are short (~9k tokens); LongMemEval uses synthetic, low-diversity dialogue. No standard metric jointly scores **accuracy × latency × token-cost × forgetting quality**.
- **Forgetting is mostly unimplemented.** Most systems grow monotonically; principled decay/eviction (what's safe to forget without hurting future recall) is unsolved.
- **Consolidation cost/timing.** Reflection/summarization is extra LLM spend; *when* to consolidate, and avoiding **summary drift** (lossy summaries of summaries), lack good defaults.

## Anti-patterns
- **Writing every turn verbatim with no extraction or forgetting** → store bloats, retrieval precision collapses, cost climbs. Extract + consolidate + decay.
- **Treating long-term memory as append-only** → contradictory stale facts coexist; the agent answers with the wrong one. Implement **invalidation/update** (Zep's temporal edges, or replace-on-contradiction).
- **Confusing execution state with memory** → checkpointing a run (LangGraph `thread_id`, durable replay) is *not* cross-session memory; see `agent-state-and-durable-execution`. Resuming a crashed run ≠ remembering a user across runs.
- **Reaching for RAG when you need memory (or vice-versa)** → personalization needs a write/update/forget path RAG pipelines don't have; grounding in authored docs doesn't need per-user memory churn. See the decision table.
- **Trusting extracted memories as ground truth** → no provenance, no confidence → hallucinated memories poison future turns. Link memories to their source turn.
- **No context-pressure policy** → silent truncation drops the very fact the user just gave you. Choose summarize/page-out explicitly.

## References
- **CoALA** — Sumers, Yao, Narasimhan, Griffiths, "Cognitive Architectures for Language Agents," arXiv:2309.02427 (the four-type taxonomy).
- **MemGPT** — Packer et al., "MemGPT: Towards LLMs as Operating Systems," arXiv:2310.08560; Letta docs (docs.letta.com) — memory blocks, self-editing memory, virtual context.
- **Generative Agents** — Park et al., "Generative Agents: Interactive Simulacra of Human Behavior," arXiv:2304.03442 (memory stream, recency/importance/relevance retrieval, reflection).
- **Zep / Graphiti** — Rasmussen et al., "Zep: A Temporal Knowledge Graph Architecture for Agent Memory," arXiv:2501.13956 (bi-temporal edges, invalidation, DMR/LongMemEval results).
- **Mem0 / LangMem** — Mem0 docs (extraction, user/session/agent scopes); LangMem docs (hot-path vs background writes). Framework storage details → `ai-datastores`.
- **Memory surveys & eval** — "Memory for Autonomous LLM Agents" survey (arXiv:2603.07670); LoCoMo (Maharana et al., snap-research.github.io/locomo) and LongMemEval (arXiv:2410.10813) benchmarks and their documented limitations on memory mutation.

Verify fast-moving framework APIs and benchmark numbers against live official docs.
