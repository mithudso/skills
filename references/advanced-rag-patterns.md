<!-- hub-reference-banner -->
> **Reference file — part of the `ai-rag-retrieval` hub.** Formerly the standalone `advanced-rag-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: advanced-rag-patterns
title: Agentic & Advanced RAG Patterns
description: >
  The advanced and agentic patterns layered on top of naive top-k RAG — the deep
  companion to rag-architecture.md (base pipeline) and ai-datastores /
  mongodb-search-ai (the vector-DB substrate). Covers the naive-RAG failure
  taxonomy (Barnett's 7 failure points, lost-in-the-middle, context rot); query
  transformation (rewrite, decomposition, HyDE, step-back, multi-query /
  RAG-Fusion + Reciprocal Rank Fusion); hybrid search & re-ranking
  (BM25+dense+RRF, cross-encoders, ColBERT late interaction, Cohere Rerank, MMR);
  advanced indexing/chunking (parent-document / small-to-big, sentence-window,
  auto-merging, semantic chunking, Anthropic Contextual Retrieval); agentic RAG
  control (retrieval-as-a-tool, Self-RAG reflection tokens, Corrective RAG/CRAG,
  Adaptive-RAG complexity routing, iterative/multi-hop, the 7-pattern agentic
  survey taxonomy); GraphRAG (Microsoft GraphRAG, Leiden community summarization,
  hybrid graph+vector); RAG evaluation (Ragas retrieval-vs-generation split,
  golden-context sets); multimodal RAG (ColPali vision-native page retrieval);
  and the long-context-vs-RAG routing debate. TRIGGER: a RAG pipeline that
  underperforms naive top-k; "advanced RAG", "agentic RAG", Self-RAG / CRAG /
  Adaptive-RAG, GraphRAG, contextual retrieval, query rewriting/decomposition,
  HyDE, RAG-Fusion, re-ranking / ColBERT, RAG evaluation / Ragas, multimodal /
  ColPali RAG, "is RAG dead / long-context vs RAG". SKIP: base RAG pipeline,
  chunking/embedding/vector-store selection (use rag-architecture.md); vector-DB
  internals HNSW/IVF, embedding models (use ai-datastores); Atlas $vectorSearch /
  $search mechanics (use mongodb-search-ai); context assembly / caching strategy
  (use llm-context-engineering).
origin: local
category: developer
version: "1.0"
updated: "2026-05-31"
tags:
  - rag
  - agentic-rag
  - graphrag
  - reranking
  - query-transformation
  - contextual-retrieval
  - rag-evaluation
  - llm
  - ai
whenToUse:
  - "a RAG system underperforms naive top-k and you need the advanced/agentic patterns"
  - "choosing query transformation: rewrite, decomposition, HyDE, step-back, multi-query/RAG-Fusion"
  - "adding hybrid search + re-ranking (cross-encoders, ColBERT, Cohere Rerank, RRF, MMR)"
  - "advanced indexing: parent-document, sentence-window, auto-merging, semantic, Contextual Retrieval"
  - "agentic RAG: retrieval-as-tool, Self-RAG, Corrective RAG, Adaptive-RAG, iterative/multi-hop, routing"
  - "GraphRAG for global/sense-making questions; hybrid graph+vector"
  - "evaluating RAG with Ragas (retrieval vs generation), golden-context sets"
  - "multimodal RAG (ColPali) and the long-context-vs-RAG routing decision"
whenNotToUse:
  - "base RAG pipeline, chunking, embeddings, vector-store selection — use rag-architecture.md"
  - "vector-DB internals (HNSW/IVF), embedding-model details — use ai-datastores"
  - "MongoDB Atlas Vector Search / $search specifics — use mongodb-search-ai"
related_skills:
  - rag-architecture
  - ai-datastores
  - mongodb-search-ai
  - llm-context-engineering
---

# Agentic & Advanced RAG Patterns

The advanced and agentic layer on top of naive RAG. `rag-architecture.md` owns
the base pipeline (chunking, embeddings, vector-store choice, the
hybrid→rerank→assemble default); `ai-datastores` / `mongodb-search-ai` own the
storage substrate. **This reference goes deep on the mechanisms those summaries
only name**: the failure taxonomy, the named agentic systems, GraphRAG
internals, the contextual-retrieval numbers, multimodal page retrieval, the
eval split, and the long-context-vs-RAG decision.

## The baseline problem — Barnett's 7 failure points

Barnett et al., *Seven Failure Points When Engineering a RAG System* (CAIN 2024,
arXiv:2401.05856) is the canonical failure taxonomy: (1) missing content, (2)
missed top-ranked docs, (3) not-in-context (consolidation/rerank failure), (4)
not-extracted (answer present, LLM misses it), (5) wrong format, (6) incorrect
specificity, (7) incomplete answer. Central lesson: **RAG quality is dominated
by retrieval, and retrieval quality can only be validated *operationally*, not
at design time.** Chroma's 2025 "context rot" work adds that supplying *more*
retrieved context past ~8K tokens often *degrades* answers — precision beats
volume, and **lost-in-the-middle** positional bias means buried evidence is
ignored even when retrieved.

The maturity ladder (Modular RAG / agentic-RAG survey framing): **Naive** (top-k,
one pass) → **Advanced** (pre/post-retrieval optimizations on a static pipeline)
→ **Modular** (composable, routed) → **Agentic** (an agent treats retrieval as a
tool it invokes iteratively, critiques, and loops).

## Query transformation (pre-retrieval)

The literal user query is rarely the optimal retrieval query.

- **Query rewriting** — LLM rephrases for retrieval (typos, acronym expansion,
  resolve chat-history references).
- **Decomposition / sub-questions** — split a multi-part question into atomic
  sub-queries, retrieve per sub-query, synthesize. Tackles *structural*
  complexity.
- **Step-back prompting** (DeepMind, Zheng et al.) — generate a more abstract
  question first to retrieve broad principles, then answer the specific. Tackles
  *informational* complexity; complementary to decomposition.
- **HyDE** (Hypothetical Document Embeddings) — LLM hallucinates an *answer*
  document, embed *that*, retrieve neighbors; a hypothetical answer sits closer
  in embedding space to real answer-docs than the question does. Strong
  zero-shot / no-labels.
- **Multi-query / RAG-Fusion** — generate 3–5 query variants, retrieve each,
  merge with **Reciprocal Rank Fusion (RRF)**: `score(d)=Σ 1/(k+rank_i(d))`,
  k≈60. RRF needs no score normalization and is also the standard hybrid-search
  merge.

## Hybrid search & re-ranking (post-retrieval)

- **Hybrid (BM25 + dense) + RRF** — the single highest-leverage upgrade over
  dense-only. BM25 rescues recall on exact terms, IDs, codes, and rare names the
  dense model ranked too low. (rag-architecture.md: hybrid recall@10 0.91 vs 0.72
  dense-only.)
- **Re-ranking** — two-stage: retrieve a wide candidate set (top 20–100), then
  rescore and keep the top few.
  - **Cross-encoders** (BGE-reranker, Cohere Rerank, Jina) jointly encode
    (query, doc) in one pass — highest quality, slow, small candidate sets only;
    +15–30% retrieval accuracy.
  - **ColBERT / late interaction** stores per-token embeddings and computes
    MaxSim — between bi- and cross-encoders on speed/quality, but inflates index
    size 1–2 orders of magnitude.
- **MMR (Maximal Marginal Relevance)** — diversity-aware selection: pick chunks
  relevant *and* dissimilar to those already chosen, cutting redundancy.

## Advanced indexing / chunking

Core insight: **the chunk you *match on* need not be the chunk you *feed the
LLM*.** Decouple them.

- **Parent-document / small-to-big** — embed small children for precise
  matching, return the larger parent for context.
- **Sentence-window** — match one sentence, return a window of N around it.
- **Auto-merging / hierarchical** — leaf→parent tree; match leaves; if enough
  siblings under a parent hit, merge up and return the parent.
- **Semantic chunking** — split where consecutive-sentence embedding distance
  spikes, aligning chunks to topic boundaries.
- **Anthropic Contextual Retrieval** (Sept 2024) — the standout technique:
  prepend a 50–100-token LLM-generated blurb situating each chunk in its whole
  document *before* embedding ("Contextual Embeddings"), plus a parallel
  **Contextual BM25** index. Generated cheaply with Haiku + prompt caching
  (~$1.02 / M doc tokens). Top-20 retrieval-failure reduction: contextual
  embeddings **−35%**; + contextual BM25 **−49%**; + reranking **−67%**.
  Anthropic's own rule: under ~200K tokens, just prompt-stuff the KB with
  caching; use Contextual Retrieval above that.

## Agentic RAG (retrieval as a tool, iterative control)

The shift from a fixed pipeline to an agent that *decides whether / what / when*
to retrieve, critiques results, and loops. Singh et al. (arXiv:2501.09136) ground
it in four primitives — reflection, planning, tool use, multi-agent — and a
7-pattern taxonomy (single-agent router, multi-agent, hierarchical, corrective,
adaptive, graph-based, agentic document workflows). The load-bearing systems:

- **Self-RAG** (Asai et al., ICLR 2024) — *trains* the LM to emit **reflection
  tokens**: a `Retrieve` decision token plus critique tokens `IsREL` (passage
  relevant?), `IsSUP` (generation supported by it?), `IsUSE` (output useful?). It
  learns *when* to retrieve and to self-critique on demand.
- **Corrective RAG / CRAG** (Yan et al., 2024) — a lightweight **retrieval
  evaluator** grades retrieved docs and routes to one of three actions:
  *correct* (use, with knowledge refinement / decompose-recompose), *incorrect*
  (discard → fall back to **web search**), or *ambiguous* (blend). Self-healing,
  model-agnostic, bolts onto any RAG.
- **Adaptive-RAG** (Jeong et al., NAACL 2024) — a trained **complexity
  classifier** routes each query to *no retrieval* / *single-step* / *iterative
  multi-step*, optimizing the cost/accuracy trade-off by not over-spending on
  easy queries.
- **Iterative / multi-hop** (IRCoT, FLARE) — interleave reasoning and retrieval,
  retrieving again as the chain-of-thought reveals new sub-questions.
- **Routing** — classify at the front door and send to the right index/tool/
  pipeline (vector vs graph vs SQL vs web).

## GraphRAG

**Microsoft GraphRAG** (Edge et al., arXiv:2404.16130) addresses naive RAG's
blindness to *global* "sense-making" questions ("main themes across the
corpus?"). Pipeline: LLM extracts an entity-relationship **knowledge graph** →
**Leiden community detection** builds a cluster hierarchy → LLM generates
**community summaries** bottom-up. At query time, **global search** map-reduces
over community summaries; **local search** traverses entity neighborhoods.
Reported ~50–70% comprehensiveness gains over vector RAG *on global questions
specifically* (Microsoft's own win-rate metric — not a general accuracy number).
Dynamic community selection later cut global-search token cost. **Hybrid
graph+vector** (Agent-G, GeAR) is the practical sweet spot: graph for
multi-hop/global, vector for fuzzy local lookups.

## RAG evaluation

Converts RAG from vibes to engineering. **Ragas** is the de-facto framework,
split by stage:

- **Retrieval:** Context Precision (relevant + well-ranked), Context Recall (did
  we retrieve everything needed — needs ground truth), Context Entities Recall,
  Noise Sensitivity.
- **Generation:** Faithfulness (every claim grounded — the primary hallucination
  guard), Answer/Response Relevancy, Factual Correctness, Semantic Similarity.

Most metrics are LLM-as-judge (no labels except recall/reference-based). Start
with **Faithfulness** (catches hallucination) + **Context Recall** (tells you if
the retrieval architecture is fundamentally sound). Build a **golden-context
set** (queries + ideal contexts + answers) to regression-test pipeline changes.
**Separate retrieval eval from generation eval** — it localizes whether a bad
answer is a retrieval miss or a generation failure. (See `llm-observability` for
the production tracing/online-eval side; this is the offline RAG-quality side.)

## Multimodal RAG

Real documents mix text, tables, figures, layout. Two families:

- **Caption-and-embed** — VLM produces text descriptions of images/tables, embed
  those (simple, lossy).
- **Vision-native / late interaction** — **ColPali** (Faysse et al.,
  arXiv:2407.01449) embeds *rendered page images* directly via a VLM into
  ColBERT-style multi-vectors, skipping OCR/layout parsing. Strong on
  visually-rich docs; cost is ~100× more vectors per page (index size / retrieval
  speed is the open concern). Successors: ColQwen2, ColSmol.

## Long-context vs RAG (the live debate)

With 1M-token windows, "RAG is dead" was a 2024 talking point. The 2024–2026
evidence settled it as *nuanced*, and the figures are **directional, not
precise** (they vary widely by benchmark):

- Long-context LLMs *can* beat RAG on accuracy when the whole corpus fits — but
  at roughly ~30–60× latency and up to ~1,000× per-query cost, with degradation
  past a model-specific length (lost-in-the-middle; multi-fact recall ~60% in
  some studies).
- RAG stays indispensable for large/dynamic corpora, cost/latency-sensitive
  serving, freshness, and source attribution.
- **2025 consensus: route.** Simple/local → RAG; complex/global/multi-hop →
  long-context or GraphRAG. "RAG → context engineering" is the late-2025
  reframing — RAG is one tool in a broader context-assembly discipline (see
  `llm-context-engineering`).

## Practical patterns

1. **Default production stack:** hybrid (BM25+dense) → RRF → cross-encoder rerank
   → top 3–5 chunks under ~8K assembled tokens. Fixes most of Barnett's points
   2–3.
2. **Decouple match-chunk from context-chunk** (small→big / windows / auto-merge).
3. **Add context before you embed** (Contextual Retrieval + Contextual BM25:
   −49% failures, −67% with rerank).
4. **Transform the query** (rewrite + decompose + HyDE for zero-shot / step-back
   for abstraction; multi-query+RRF when recall is the bottleneck).
5. **Make retrieval a tool and let the agent loop** (CRAG grade-and-fallback,
   Adaptive-RAG complexity routing, Self-RAG faithfulness self-critique) — with
   iteration caps.
6. **GraphRAG for global/sense-making; vector RAG for local; hybridize.**
7. **Route long-context vs RAG by query type and corpus size**; under ~200K
   tokens consider prompt-stuffing with caching.
8. **Evaluate continuously** (golden set + Ragas Faithfulness + Context Recall;
   measure retrieval and generation separately).
9. **Order context deliberately** — strongest evidence first and last to dodge
   lost-in-the-middle.

## Anti-patterns

- **Eval-free RAG** — no retrieval/generation harness; the field's #1 mistake.
- **Dense-only, no lexical** — silently under-recalls exact terms/IDs/codes.
- **No re-ranking** — raw top-k noise dilutes signal (rerank +15–30%).
- **Over-/fixed-size chunking** — severs semantic units; ~80% of failures start
  in ingestion/chunking.
- **Context stuffing** — dumping 50K tokens; degrades past ~8K (context rot) +
  lost-in-the-middle.
- **Treating naive top-k as sufficient** — blaming the LLM for retrieval defects.
- **Assuming long context kills RAG** — paying ~1000× to stuff a corpus a routed
  RAG path answers cheaper.
- **Unbounded agentic loops** — reflection/multi-hop without iteration or cost
  caps.

## Cross-references

- **Base RAG pipeline** (chunking, embeddings, vector-store choice, the
  hybrid→rerank default, quick-start code) → `rag-architecture.md` (this is its
  deep companion).
- **Vector-DB internals, embedding models, agent memory** → `ai-datastores`.
- **MongoDB Atlas `$vectorSearch` / `$search` / `$rankFusion`** →
  `mongodb-search-ai`.
- **Context assembly, prompt caching, "RAG → context engineering"** →
  `llm-context-engineering`.
- **Production tracing & online eval (RAGAS in prod)** → `llm-observability`.

## References

Barnett et al. *7 Failure Points* (arXiv:2401.05856); Singh et al. *Agentic RAG
Survey* (2501.09136); Asai et al. *Self-RAG* (2310.11511, ICLR 2024); Yan et al.
*CRAG* (2401.15884); Jeong et al. *Adaptive-RAG* (2403.14403, NAACL 2024);
Rackauckas *RAG-Fusion* (2402.03367); Edge et al. *GraphRAG* (2404.16130);
Faysse et al. *ColPali* (2407.01449); *Long Context vs RAG* (2501.01880).
Anthropic *Contextual Retrieval*; Microsoft Research *GraphRAG*; Ragas metrics
docs; Databricks *Long-Context RAG Performance*; Weaviate / Qdrant / RAGFlow
practitioner guides. *(28 sources, 2024–2026; full URLs in the source research
report.)*
