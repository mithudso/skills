<!-- hub-reference-banner -->
> **Reference file — part of the `ai-rag-retrieval` hub.** Formerly the standalone `ai-datastores` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: ai-datastores
version: 1.1.0
updated: 2026-05-29
description: >
  AI datastore expert — vector databases (Pinecone, Qdrant, Weaviate, Chroma,
  Milvus, pgvector, MongoDB Atlas Vector Search), knowledge graphs (Neo4j,
  Graphiti, FalkorDB), document stores for agent state (MongoDB, DynamoDB,
  Firestore), memory frameworks (Mem0, Letta/MemGPT, Graphiti), caching (Redis
  for agents), and hybrid storage architectures. TRIGGER: user needs to choose
  a datastore for an AI/agent application, implement vector search, design agent
  memory, compare storage options, or build a GraphRAG or RAG pipeline.
  SKIP: general database design unrelated to AI/ML; pure full-text search
  without vector components; data warehouse or analytics workloads; ML model
  training infrastructure; embedding model architecture or fine-tuning.
origin: local
updated: 2026-05-29
tags: [vector-database, agent-memory, knowledge-graph, RAG, GraphRAG, redis, mongodb, pgvector]
related_skills: [agent-ecosystem, agent-plan-writing, rag-architecture, mongodb-atlas-vector-search]
---

# AI Datastores

Reference for storage systems powering AI applications and agent memory. Deep benchmarks and code samples are in `references/`.

## Quick datastore selection

| Need | Recommended |
| --- | --- |
| Managed zero-ops at scale | Pinecone |
| Best price-performance (mid-scale) | Qdrant |
| Hybrid search (vector + keyword + metadata) | Weaviate |
| Rapid prototyping | Chroma |
| 100B+ vectors with GPU | Milvus/Zilliz |
| Already on PostgreSQL | pgvector |
| Document + vector in one collection | MongoDB Atlas Vector Search |
| Temporal fact tracking | Graphiti + Neo4j/FalkorDB |
| Agent memory (drop-in API) | Mem0 |
| Agent memory (autonomous curation) | Letta (MemGPT) |
| Sub-ms coordination + caching | Redis |
| GraphRAG with multi-hop reasoning | Neo4j + LangChain |

## When NOT to use this skill

- General database design unrelated to AI/ML — use a standard DB skill
- Pure full-text search without vector or semantic components
- Data warehouse or analytics workloads
- ML model training infrastructure
- Embedding model architecture or fine-tuning

---

## Vector databases: deep comparison (2026)

### Performance benchmarks

| Database | p50 latency (10M vectors) | p99 latency (10M vectors) | Max practical scale | Hybrid search |
| --- | --- | --- | --- | --- |
| Qdrant | <5ms | ~12ms | 100M+ (self-hosted) | Yes |
| Weaviate | ~8ms | ~16ms | 100M+ (managed) | Yes (best-in-class) |
| Milvus | ~10ms | ~18ms | Billions (GPU/sharded) | Yes |
| Pinecone Serverless | ~15ms | ~30ms | Unlimited (managed) | No (dense only) |
| pgvector | ~12ms | ~25ms | 50M comfortably | Via tsvector combo |
| MongoDB Atlas VS | <50ms (quantized) | ~80ms | 100M+ | Yes (Atlas Search) |
| Chroma | <5ms (small sets) | ~15ms | ~5M practical | No |

### Pricing (2026 managed cloud)

| Database | 10M vectors/month | 100M vectors/month | Self-hosted | Free tier |
| --- | --- | --- | --- | --- |
| Pinecone Serverless | ~$70 | ~$700+ | No | 2GB |
| Qdrant Cloud | ~$65 | ~$200 | Yes ($30/mo VPS) | 1GB |
| Weaviate Cloud | ~$135 | ~$400 | Yes | Sandbox |
| Milvus/Zilliz | ~$80 | ~$250 | Yes (open source) | Free tier |
| pgvector (RDS) | ~$45 | ~$100 | Yes (any Postgres) | N/A |
| MongoDB Atlas VS | Included in Atlas tier | Included in Atlas tier | Yes (self-managed) | M0 free |
| Chroma | Free (self-hosted) | Not recommended | Yes | Unlimited local |

Self-hosted Qdrant on a $30/month VPS handles 10M+ vectors — 10x cheaper than equivalent Pinecone. The crossover where self-hosting beats Pinecone is roughly $600/month in vector DB costs.

### Decision tree

```
How many vectors?
  |
  +-- <1M
  |     Prototyping? → Chroma (zero config, embedded mode)
  |     Already on Postgres? → pgvector (no new infra)
  |     Already on MongoDB? → Atlas Vector Search (unified)
  |     Need managed + simple? → Pinecone (free tier)
  |
  +-- 1M–50M
  |     Budget-sensitive? → pgvector or self-hosted Qdrant
  |     Need hybrid search? → Weaviate (best native hybrid)
  |     Need fastest filtered search? → Qdrant
  |     Want zero ops? → Pinecone Serverless
  |     Need document+vector unified? → MongoDB Atlas VS
  |
  +-- 50M–500M
  |     Self-host OK? → Qdrant or Milvus
  |     Need managed? → Weaviate Cloud or Zilliz
  |     On Postgres? → pgvector with HNSW tuning (watch rebuild times)
  |
  +-- 500M+
        GPU available? → Milvus with GPU indexing
        Managed only? → Pinecone (pod-based) or Zilliz
        Need sharding? → Milvus (most mature partitioning)
```

### Per-database notes

**Pinecone** — Zero-ops managed service. Serverless trades latency for convenience; pod-based (p2) is faster but significantly more expensive. No self-hosted option.

**Qdrant** — Rust-based, consistently fastest p50/p99 among open-source options. Best-in-class filtering performance. Supports scalar, product, and binary quantization. Strong self-hosted economics.

**Weaviate** — Best native hybrid search (keyword + vector in one query). Multi-tenancy improvements shipped in 2026.

**Chroma** — Best developer experience for local prototyping — zero to working RAG in under 10 minutes. Practical ceiling ~5M vectors; plan migration to Qdrant or pgvector for production.

**Milvus/Zilliz** — Most popular open-source vector DB by GitHub stars. Built for billion-scale with the most mature sharding and partitioning. GPU indexing support.

**pgvector** — For most teams starting new RAG projects in 2026, pgvector is the right default. Relational data and vector data live in the same database, queryable together with ACID transactions. Beyond 50–100M vectors, HNSW index rebuild times become a constraint.

**MongoDB Atlas Vector Search** — Unified document + vector store eliminates synchronization between two systems. Voyage AI Automated Embeddings (2026 preview) auto-generates embeddings on write/update. MongoDB 8.3 improves read performance up to 45% and write up to 35%. With scalar or binary quantization at 2048 dimensions, retains 90–95% accuracy with <50ms query latency.

---

## Knowledge graphs for AI agents

### When vector search is not enough

Use a knowledge graph when:
- You need multi-hop reasoning ("which customers of company X also use product Y?")
- Entity relationships and provenance matter more than semantic similarity
- You need temporal fact tracking ("what was true at time T?")
- You want to reduce hallucinations with deterministic graph traversal
- You need explainable retrieval paths (audit trails)

### GraphRAG architecture (2026 pattern)

```
User Query
    |
Planner Agent (intent classification)
    |
    +-- Structural query --> Graph Retriever --> Cypher on Neo4j --> Subgraph
    |                                                                    |
    +-- Semantic query --> Vector Retriever --> Embedding search -------+
    |                                                                    |
                         Result Merger + Reranker
                                |
                         LLM Generation (grounded response)
```

The LLM translates natural language into a deterministic graph query, retrieving factual ground truth before generating — reducing hallucinations significantly.

### Neo4j
Primary graph database for AI applications. Native graph storage with Cypher. Neo4j GraphRAG Context Provider integrates with Microsoft Agent Framework and Google ADK. Best for enterprise knowledge graphs where entity relationships and compliance auditing matter.

### FalkorDB
Redis-compatible graph database optimized for low-latency queries. Sub-10ms query performance. Supports the Graphiti framework as an alternative backend to Neo4j.

### Graphiti (by Zep)
Python framework for temporally-aware knowledge graphs for multi-agent systems. Every fact carries a validity window (`valid_from`, `valid_to`). Production-grade as of 2026.

Supported backends: Neo4j (primary), FalkorDB (sub-10ms multi-agent), AWS Neptune, Kuzu.

Key capabilities:
- Episodes auto-decompose into entities, relationships, and temporal metadata
- Hybrid search across graph structure, vector similarity, and full-text
- Multi-tenant isolation for SaaS agent deployments
- MCP server available for direct agent integration

---

## Agent memory frameworks (2026)

| Framework | Best for | Architecture | Key stat |
| --- | --- | --- | --- |
| Mem0 | Drop-in API for any framework | Hybrid vector + graph + KV | 48K+ stars, $24M Series A |
| Letta (MemGPT) | Autonomous memory curation | OS-inspired virtual context (RAM / disk) | Agent manages its own memory |
| Zep / Graphiti | Multi-agent, temporal reasoning | Temporal KG with episode decomposition | 15 points higher on LongMemEval for temporal queries |
| LangMem | Teams already in LangGraph | Pluggable backends via LangGraph checkpointers | Tightest LangGraph integration |

### Accuracy vs latency trade-offs

| Pattern | Accuracy | p95 Latency |
| --- | --- | --- |
| Simple buffer (recent N turns) | ~66.9% | 1.44s |
| Summarization + retrieval | ~70.2% | 4.8s |
| Vector-indexed episodic | ~72.1% | 8.2s |
| Graph-augmented retrieval | ~72.9% | 17.1s |
| Temporal KG (FalkorDB) | ~72.9% | 12.4s |

---

## Redis for AI agents (2026)

Redis has evolved from cache to agent infrastructure platform.

| Capability | Product | Latency | Use case |
| --- | --- | --- | --- |
| Session state | Redis data structures | <1ms | Current conversation, active tasks |
| Semantic caching | Redis LangCache | <5ms | Deduplicate similar queries (up to 70% cost reduction) |
| Vector search | Redis Search | <10ms | RAG retrieval |
| Agent memory | Agent Memory Server | <5ms | Dual-tier persistent memory |
| Rate limiting | Redis rate limiter | <1ms | API management for tool calls |
| Pub/Sub coordination | Redis Streams | <1ms | Multi-agent event bus |

**Redis Iris (2026):** Context and memory platform combining real-time data ingestion, a semantic interface that auto-generates MCP tools from business data models, and an agent memory server.

**When to use Redis vs alternatives:**
- Use Redis for sub-millisecond coordination, semantic caching, and agent memory that fits in memory (<100GB)
- Add a vector DB when you need billion-scale vector search exceeding Redis memory economics
- Add a graph DB when you need multi-hop entity reasoning

---

## Hybrid storage architecture (2026 best practice)

| Memory type | Storage | Tech | Retention | Access |
| --- | --- | --- | --- | --- |
| Working memory (session) | In-memory | Redis | 15min–2hr | <1ms |
| Episodic memory (history) | Document store | MongoDB, PostgreSQL | Weeks–months | 5–50ms |
| Semantic memory (embeddings) | Vector DB | Pinecone, Qdrant, Atlas VS | Indefinite | 10–50ms |
| Relational memory (entities) | Graph DB | Neo4j, FalkorDB/Graphiti | Indefinite | 5–20ms |
| Procedural memory (workflows) | Structured store | MongoDB, PostgreSQL | Indefinite | 10–50ms |
| Durable record (audit) | SQL / object storage | PostgreSQL, S3 | Indefinite | 50–200ms |

### Runtime flow

```
Agent receives request
1. Pull working context from Redis (<1ms)
2. Retrieve semantically similar past cases from vector DB (10–50ms)
3. Query entity relationships from knowledge graph if needed (5–20ms)
4. Generate response with LLM

Agent completes work
5. Update working memory in Redis
6. Append to episodic store (MongoDB/Postgres)
7. Update long-term embeddings in vector DB
8. Archive to durable record (S3/SQL)
```

### Unified vs specialized guidance

**Default to unified PostgreSQL stack** (pgvector + relational + tsvector) unless you have a specific requirement:
- Billion-scale vectors → add Milvus or Qdrant
- Extreme multi-tenant concurrency → add dedicated vector DB
- Temporal entity reasoning → add Graphiti + Neo4j/FalkorDB
- Sub-ms caching/coordination → add Redis

**Default to MongoDB Atlas** if already in the MongoDB ecosystem — Atlas Vector Search + document storage + Atlas Search provides a strong unified platform without additional infrastructure.

---

## Document stores for agent state

**MongoDB** — Best for flexible schemas that evolve as agent capabilities grow. Native change streams enable real-time agent coordination. Atlas Vector Search adds semantic retrieval without a second database.

**DynamoDB** — Best for serverless agent architectures on AWS. Single-digit millisecond reads at any scale. Lacks vector search natively (pair with Amazon OpenSearch or Bedrock Knowledge Bases).

**Firestore** — Best for real-time agent state in Firebase/GCP ecosystems. Real-time listeners enable instant state synchronization across agent instances.

---

## Integration patterns

### LangChain / LangGraph

```python
# pgvector with LangChain
from langchain_postgres import PGVector
vectorstore = PGVector(
    connection=connection_string,
    embeddings=embeddings,
    collection_name="agent_memory",
)

# MongoDB Atlas Vector Search with LangChain
from langchain_mongodb import MongoDBAtlasVectorSearch
vectorstore = MongoDBAtlasVectorSearch(
    collection=collection,
    embedding=embeddings,
    index_name="vector_index",
)

# Redis Agent Memory with LangGraph
from langgraph.checkpoint.redis import RedisSaver
checkpointer = RedisSaver(redis_url="redis://localhost:6379")
```

### Graphiti for temporal memory

```python
from graphiti_core import Graphiti
from graphiti_core.llm_client import AnthropicClient

graphiti = Graphiti(
    neo4j_uri="bolt://localhost:7687",
    neo4j_user="neo4j",
    neo4j_password="password",
    llm_client=AnthropicClient(),
)

# Add an episode (auto-decomposes into entities + relationships)
await graphiti.add_episode(
    name="customer_call",
    episode_body="Customer mentioned switching from Pinecone to Qdrant due to costs",
    source_description="support_ticket",
)

# Search with temporal awareness
results = await graphiti.search("customer vector database preferences", num_results=5)
```

### Production hardening (pgvector)

```python
from langchain_postgres import PGVector
from sqlalchemy import create_engine
from tenacity import retry, stop_after_attempt, wait_exponential

engine = create_engine(
    connection_string,
    pool_size=10,
    max_overflow=20,
    pool_timeout=30,
    pool_recycle=1800,
)

@retry(stop=stop_after_attempt(3), wait=wait_exponential(multiplier=1, max=10))
def search_with_retry(query: str, k: int = 5):
    return vectorstore.similarity_search(query, k=k)
```

---

## Embedding model selection

| Use case | Model | Dimensions | Notes |
| --- | --- | --- | --- |
| General-purpose RAG | OpenAI text-embedding-3-large | 3072 (or 1536 via shortening) | Best quality/cost balance |
| Cost-sensitive RAG | OpenAI text-embedding-3-small | 1536 | Good enough for most retrieval |
| Privacy-sensitive / on-prem | Sentence-Transformers (all-MiniLM-L6-v2) | 384 | Free, self-hosted, fast |
| Multilingual | Cohere embed-v3 | 1024 | 100+ languages, int8 quantization |
| MongoDB Atlas (auto) | Voyage AI (via Atlas integration) | 1024–2048 | Zero-config; auto-generates on write |
| Maximum quality | Voyage 3 Large | 2048 | Best MTEB scores; higher cost |

Rules:
- Never mix embedding models in the same index — dimensions and semantic spaces differ
- Higher dimensions are not always better; 1024-dim models often match 3072-dim at 1/3 the storage cost
- Always benchmark on your data; public leaderboard scores may not reflect your domain

---

## Common mistakes

1. **Premature specialization:** Adding a dedicated vector DB when pgvector or Atlas Vector Search handles your scale.
2. **Ignoring hybrid search:** Pure vector search underperforms hybrid (keyword + vector) on most production workloads.
3. **Mixing embedding models:** Storing vectors from different models in the same collection produces meaningless similarity scores.
4. **Skipping quantization:** Scalar or binary quantization cuts storage 4–8x with minimal recall loss (90–95% retained).
5. **No memory eviction strategy:** Agent memory grows unbounded without TTLs or summarization.
6. **Over-engineering graph:** Graphs pay off only when you need multi-hop or temporal reasoning; single-hop semantic lookups don't need them.
7. **Redis as sole persistence:** Redis is volatile by default. Always pair with a durable store.
8. **Ignoring cold start latency:** HNSW index loading on pgvector and Qdrant can take minutes for large indexes after a restart.

---

## Edge cases and failure modes

| Scenario | Impact | Mitigation |
| --- | --- | --- |
| Embedding dimension mismatch | Silent bad results | Validate dimensions on insert; reject mismatched vectors |
| HNSW index rebuild during writes | High write latency spikes on pgvector (minutes at 50M+) | Schedule rebuilds during low-traffic windows |
| Pinecone serverless cold start | First query after idle takes 2–5s | Use keep-warm pings for latency-sensitive apps |
| Graphiti episode explosion | Memory and latency grow with entity count | Set entity merge thresholds; prune stale relationships |
| Redis memory exhaustion | Agent crashes if memory exceeds allocation | Set maxmemory with eviction policy (allkeys-lru); monitor with alerts |
| MongoDB Atlas VS index sync lag | Newly inserted documents not immediately searchable | Allow 1–2s propagation delay; use change streams for confirmation |
| Chroma at scale | Performance degrades sharply past ~5M vectors | Plan migration to Qdrant or pgvector before hitting ceiling |

---

## Cited sources

1. [Vecstore — Vector Database Performance Compared](https://vecstore.app/blog/vector-database-performance-compared)
2. [MarkTechPost — Best Vector Databases 2026](https://www.marktechpost.com/2026/05/10/best-vector-databases-in-2026-pricing-scale-limits-and-architecture-tradeoffs-across-nine-leading-systems/)
3. [LeanOps — Vector DB Cost Comparison 2026](https://leanopstech.com/blog/vector-database-cost-comparison-2026/)
4. [Graphiti GitHub](https://github.com/getzep/graphiti)
5. [AgentMarketCap — Agent Memory Vendor Landscape 2026](https://agentmarketcap.ai/blog/2026/04/10/agent-memory-vendor-landscape-2026-letta-zep-mem0-langmem)
6. [Redis — AI Agent Architecture 2026](https://redis.io/blog/ai-agent-architecture/)
7. [MongoDB Atlas Vector Search Benchmarks](https://www.mongodb.com/docs/atlas/atlas-vector-search/benchmark/results/)
8. [Neo4j — GraphRAG and Agentic Architecture](https://neo4j.com/blog/developer/graphrag-and-agentic-architecture-with-neoconverse/)
9. [Zylos Research — Agent-Native Hybrid Storage 2026](https://zylos.ai/research/2026-04-14-agent-native-data-layer-hybrid-storage-architectures)
10. [Encore — You Probably Don't Need a Vector Database](https://encore.dev/blog/you-probably-dont-need-a-vector-database)
