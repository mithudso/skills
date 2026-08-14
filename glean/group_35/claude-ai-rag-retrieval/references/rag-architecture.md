<!-- hub-reference-banner -->
> **Reference file — part of the `ai-rag-retrieval` hub.** Formerly the standalone `rag-architecture` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: rag-architecture
title: "RAG Architecture Expert"
version: "1.1.0"
updated: "2026-05-29"
description: >
  Retrieval-Augmented Generation architecture: chunking strategies, embedding pipelines,
  retrieval patterns (dense, sparse, hybrid BM25+vector+RRF), reranking (cross-encoders,
  ColBERT), vector store selection, query transformation (HyDE, multi-query, step-back),
  evaluation (RAGAS, DeepEval), advanced patterns (self-RAG, corrective RAG, GraphRAG,
  agentic RAG, adaptive RAG), context window optimization, and citation/attribution.
  TRIGGER: user is designing, building, evaluating, or debugging a RAG pipeline; user
  asks about chunking, embeddings, vector databases, hybrid search, reranking, or
  retrieval failures; user wants to compare RAG vs fine-tuning vs long-context.
  SKIP: general LLM selection without RAG (→ llm-models); MongoDB-specific vector search
  configuration (→ mongodb-search-ai); agent orchestration without RAG (→ agent-ecosystem).
category: developer
tags:
  - rag
  - retrieval-augmented-generation
  - vector-search
  - embeddings
  - chunking
  - reranking
  - hybrid-search
  - evaluation
  - ragas
  - deepeval
  - llm
  - ai
whenToUse:
  - "Designing a RAG pipeline or retrieval-augmented generation system"
  - "Choosing chunking strategies (fixed, semantic, recursive, agentic)"
  - "Selecting embedding models (OpenAI, Cohere, Voyage, BGE, open-source)"
  - "Choosing a vector database (Pinecone, Qdrant, Weaviate, Milvus, pgvector, MongoDB Atlas)"
  - "Implementing hybrid search (BM25 + dense vectors + reciprocal rank fusion)"
  - "Adding reranking (Cohere Rerank, cross-encoders, ColBERT)"
  - "Implementing query transformation (HyDE, multi-query, step-back prompting)"
  - "Evaluating RAG quality with RAGAS, DeepEval, or custom metrics"
  - "Debugging retrieval failures, poor faithfulness, or production RAG issues"
  - "Comparing RAG vs fine-tuning vs long-context approaches"
---

# RAG Architecture Expert

RAG augments LLM generation with external knowledge retrieved at query time, grounding responses in factual, up-to-date source material. As of 2026, RAG has matured from a simple retrieve-and-generate prototype into a modular production architecture with hybrid retrieval, cross-encoder reranking, agentic self-correction loops, and rigorous evaluation.

**When to use RAG vs alternatives:**

| Approach | Best For | Limitations |
|----------|----------|-------------|
| RAG | Dynamic data, attribution needs, large corpora | Retrieval quality ceiling, latency |
| Fine-tuning | Domain style/format, behavioral changes | Stale data, no source attribution, expensive |
| Long context | Small corpora (<200K tokens, <100 docs) | Cost at scale, attention degradation |
| Knowledge graphs | Multi-hop relational queries | Complex setup, maintenance overhead |

The recommended hybrid: fine-tune for domain language + RAG for factual freshness.

---

## RAG Pipeline Architecture

### Ingestion Pipeline (Offline)

1. **Document Loading** — Ingest from PDFs, HTML, databases, APIs, Confluence, Notion
2. **Preprocessing** — Clean, normalize, extract tables/images for multimodal RAG
3. **Chunking** — Split documents into retrievable segments (see Section 2)
4. **Embedding** — Convert chunks to dense vectors (see Section 3)
5. **Indexing** — Store vectors + metadata in a vector database (see Section 4)
6. **Metadata Enrichment** — Attach source URL, date, section title, document type for filtering

### Retrieval Pipeline (Online)

1. **Query Processing** — Optional query transformation (HyDE, multi-query, decomposition)
2. **Retrieval** — ANN search over vectors, optionally combined with BM25 (hybrid)
3. **Reranking** — Cross-encoder rescores top candidates
4. **Context Assembly** — Pack selected chunks into the LLM prompt with consistent formatting
5. **Generation** — LLM produces a response grounded in retrieved context
6. **Citation** — Attach source references to generated claims

### Architecture Taxonomy

| Generation | Description | Cost | Latency |
|-----------|-------------|------|---------|
| Naive RAG | Embed → vector search → generate | $0.001–0.01/query | 100–500ms |
| Advanced RAG | + hybrid retrieval, reranking, query transformation | $0.003–0.01/query | 500ms–2s |
| Modular RAG | Independently replaceable pipeline components | Varies | Varies |
| Agentic RAG | Autonomous agent decides retrieval strategy | 2–5x base | 2–10s+ |
| Adaptive RAG | Router selects optimal pipeline per query | Optimal per query | Optimal |

---

## Chunking Strategies

Chunking quality is the single most impactful factor in RAG retrieval accuracy. Semantic chunking achieves 0.79–0.82 faithfulness vs 0.47–0.51 for naive fixed-size — a 70% improvement.

| Strategy | Config | Accuracy | Use when |
|----------|--------|----------|----------|
| **Fixed-size** | 256–512 tokens, 10–20% overlap | Baseline | Prototyping, uniform document structure |
| **Recursive** | 512 tokens, 15–20% overlap | 69% (best balance) | Default starting point for most applications |
| **Semantic** | SemanticChunker, 90th percentile threshold | Up to 70% lift | High-accuracy requirements, varied topics |
| **Agentic** | LLM makes segmentation decisions | Highest | Complex documents with implicit structure |
| **Parent-Child** | Small chunks for retrieval, expand to parent | High precision + recall | Long documents where both matter |
| **Late Chunking (2026)** | Embed full doc first, then split embeddings | High context-awareness | Documents where local depends on global meaning |
| **Contextual Retrieval** | Prepend document context summary to each chunk | 49% fewer retrieval failures | Near-universal benefit in production |

**Practical defaults:** start with recursive at 512 tokens + 15–20% overlap + contextual retrieval headers. Graduate to semantic chunking when accuracy demands it.

---

## Embedding Models

Documents and queries must be embedded with the same model. The embedding model and generation LLM are completely independent.

| Model | MTEB | Cost | Dims | Best For |
|-------|------|------|------|----------|
| text-embedding-3-small | 64.6 | $0.02/1M | 1536 | Default starting point (90% of projects) |
| text-embedding-3-large | 66.4 | $0.13/1M | 3072 | When quality > cost |
| Cohere embed-v4 | ~67 | $0.10/1M | 1024 | Multilingual (100+ languages) |
| Voyage voyage-3-large | Top retrieval | $0.12/1M | 1024 | Retrieval-focused workloads |
| Qwen3-Embedding-8B | 70.6 | Self-host | 4096 | Surpasses all API models |
| BGE-M3 | Competitive | Free (OSS) | 1024 | Self-hosting, cost control |

**Best practices:**
- 1024-dim models offer best quality-to-cost ratio for most RAG applications
- Matryoshka embeddings (text-embedding-3) support truncating dimensions with graceful quality degradation
- Batch in groups of 100–2000 during ingestion
- Some models (e5, BGE) require task-specific prefixes ("query: " vs "passage: ")
- Test on your corpus with 50–200 real queries before committing

---

## Vector Stores

| Database | Type | Best For | Hybrid Search | Scale |
|----------|------|----------|---------------|-------|
| **Pinecone** | Managed | Ease of use, production default | Via metadata | Billions |
| **Qdrant** | OSS/Managed | Lowest latency, advanced filtering | Native | Billions |
| **Weaviate** | OSS/Managed | Native hybrid search champion | BM25 + vector + metadata | Billions |
| **Milvus/Zilliz** | OSS/Managed | Extreme scale, GPU acceleration | Via separate pipeline | Billions+ |
| **pgvector** | Postgres ext | Teams with existing Postgres | Via pg_search/ParadeDB | Millions |
| **MongoDB Atlas Vector Search** | Managed | Unified data + vector, auto-embedding | Native `$rankFusion` | Billions |
| **Chroma** | OSS | Local dev, prototyping | Limited | Millions |

**MongoDB Atlas Vector Search** (2026): auto-embedding in Public Preview (May 2026) — MongoDB generates and manages embeddings automatically. Native `$rankFusion` combines vector search with full-text search.

**Selection:**
1. Already have Postgres? → pgvector
2. Already have MongoDB? → Atlas Vector Search
3. Need managed + simple? → Pinecone
4. Need hybrid search natively? → Weaviate or MongoDB Atlas
5. Need lowest latency + best filtering? → Qdrant
6. Prototyping locally? → Chroma

---

## Retrieval Patterns

### Dense vs Sparse

| Method | Strengths | Weaknesses | Recall@10 standalone |
|--------|-----------|------------|---------------------|
| Dense (ANN) | Captures semantic meaning, synonyms | Misses exact matches | 0.72 |
| Sparse (BM25) | Exact token matching, interpretable | Misses synonyms | Lower than dense |

### Hybrid Search (Dense + Sparse + RRF)

Run BM25 and vector search in parallel; fuse results with Reciprocal Rank Fusion:

```
RRF_score(d) = sum( 1 / (k + rank_i(d)) ) for each retriever i
```
Where k = 60. Hybrid (BM25 + dense + RRF) achieves recall@10 of **0.91** vs 0.72 for dense alone — a 26% improvement.

**Config:** 60% semantic / 40% lexical; retrieve k=20–50 candidates per retriever before fusion.

---

## Reranking

First-stage retrieval optimizes recall. Reranking optimizes precision by scoring top candidates with a more expensive model.

### Production Pipeline (2026 default)

```
BM25 → top 50 candidates
Dense bi-encoder → top 50 candidates
Reciprocal Rank Fusion → merge to top 100
Cross-encoder rerank → top 5–10
```

### Cross-Encoder Models

| Model | Performance | Use when |
|-------|------------|----------|
| Cohere Rerank 4 Pro | +170 ELO over v3.5 | Production default |
| BGE-reranker | Good, self-hostable | Self-hosting |
| FlashRank | Milliseconds, no GPU | Latency-critical, edge |

Cross-encoder reranking adds 10–25% additional precision on top of hybrid retrieval; 40% accuracy improvement in controlled studies.

---

## Query Transformation

| Technique | Improvement | Cons | Best for |
|-----------|------------|------|----------|
| **HyDE** | 10–30% retrieval improvement | Hallucinated hypothetical can misdirect | Abstract/conceptual queries |
| **Multi-Query** | Covers different semantic interpretations | 3–5x retrieval cost | Complex multi-part questions |
| **Step-Back** | Provides foundational context | May retrieve overly broad content | Questions requiring domain background |
| **Query Routing** | Optimal cost per query type | Classifier overhead | Mixed query complexity |

Techniques are combinable: decompose a query, generate multiple variants of each sub-query, apply HyDE to each.

---

## Evaluation

### RAGAS Core Metrics

| Metric | Measures | Target |
|--------|----------|--------|
| Faithfulness | Is the answer grounded in retrieved context? | >0.90 |
| Answer Relevancy | Does the answer address the question? | >0.85 |
| Context Precision | Are relevant items ranked high? | >0.80 |
| Context Recall | Is all necessary information retrieved? | >0.80 |
| Noise Sensitivity | Robustness against irrelevant chunks? | Low |

### RAGAS vs DeepEval

- **RAGAS:** best for deep RAG-specific evaluation, academic-grade metrics, observability integration
- **DeepEval:** best for broad LLM evaluation in CI/CD, pytest integration, non-RAG metrics (bias, toxicity, agents)
- **Production pattern:** use RAGAS to define the metric suite, DeepEval to enforce thresholds in CI/CD

---

## Context Window Optimization

- Place most relevant chunks first (attention is strongest at start)
- Use stable delimiters and explicit source labels: `[Source: <title> | <url> | <date>]`
- Deduplicate near-duplicate chunks before packing
- **Optimal chunk count:** 4–5 targeted chunks consistently outperform 10 untargeted ones
- Token budget: `retrieved_chunks = total_context - system_prompt - history - generation_buffer`

---

## Citations

Source attribution drives a disproportionate share of perceived-quality lift. A day spent on citation UI beats a week of prompt tuning.

```
You are a helpful assistant. Answer using ONLY the provided context.
For every factual claim, add an inline citation [n] referencing the source.
If the context does not contain enough information, say "I don't have enough
information to answer that."

Context:
{formatted_chunks_with_source_labels}

Question: {user_query}
```

Track **citation coverage** (% of claims with sources) as a production metric — target >90%.

---

## Quick-Start Implementation

```python
# Minimal production RAG (LangChain LCEL + MongoDB Atlas)
from langchain_openai import OpenAIEmbeddings, ChatOpenAI
from langchain_mongodb import MongoDBAtlasVectorSearch
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.chains import create_retrieval_chain
from langchain.chains.combine_documents import create_stuff_documents_chain
from langchain_core.prompts import ChatPromptTemplate

splitter = RecursiveCharacterTextSplitter(chunk_size=512, chunk_overlap=80)
chunks = splitter.split_documents(documents)

embeddings = OpenAIEmbeddings(model="text-embedding-3-small")
vectorstore = MongoDBAtlasVectorSearch.from_documents(
    chunks, embeddings, collection=collection, index_name="vector_index"
)

retriever = vectorstore.as_retriever(search_type="similarity", search_kwargs={"k": 5})
prompt = ChatPromptTemplate.from_messages([
    ("system", "Answer using only the context. Cite sources with [n].\n\n{context}"),
    ("human", "{input}")
])
chain = create_retrieval_chain(retriever, create_stuff_documents_chain(ChatOpenAI(model="gpt-4o"), prompt))
result = chain.invoke({"input": "What is the refund policy?"})
```

---

## Decision Flowchart

```
Do I need RAG?
  Static data fits in context? → Long-context prompting
  Need behavioral/style changes only? → Fine-tuning
  Dynamic data, attribution needed, large corpus? → RAG
    |
    Chunking: Recursive (default) → Semantic (high accuracy)
    Embeddings: text-embedding-3-small (default) → domain-specific
    Vector store: based on existing stack (see Section 4)
    Retrieval: Hybrid BM25 + Dense + RRF (proven default)
    Reranking: Cross-encoder on top 20–100 candidates
    Evaluation: RAGAS + DeepEval in CI/CD
    Scale complexity: Naive → Advanced → Agentic as needed
```

For advanced patterns (Self-RAG, Corrective RAG, GraphRAG, Adaptive RAG), anti-patterns, and production monitoring checklists, see `references/advanced-rag-patterns.md`.

---

## Cross-References

- **ai-datastores**: vector database deep dives, knowledge graph storage, agent memory
- **llm-context-engineering**: context window design, prompt caching, chunk vs system prompt decisions
- **mongodb-search-ai**: MongoDB Atlas Vector Search specifics, `$vectorSearch`, analyzers, search nodes
- **llm-models**: embedding model details, generation model selection, API patterns
- **agent-ecosystem**: agentic RAG orchestration, LangGraph, CrewAI integration
