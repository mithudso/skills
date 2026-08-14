# Data Types, Extensions & pgvector

> Reference for **postgresql-expert**. The rich type system and extension ecosystem that distinguish PostgreSQL, plus pgvector as the AI/RAG crossover. Verified-as-of 2026-06-30. RAG retrieval architecture (chunking, reranking) routes to `ai-rag-retrieval`; this owns the storage/index layer.

## Rich data types

PostgreSQL's type system is a major reason teams choose it over other RDBMSs:

- **JSONB** — binary, indexed JSON. Prefer over `json` (which stores text verbatim). Query with `->`, `->>`, `#>`, and containment `@>`; index with **GIN** (`jsonb_path_ops` for `@>`-only, smaller). Supports `jsonpath` (`@@`, `jsonb_path_query`). Use JSONB for genuinely schemaless/variable attributes — but don't model relational data as one JSONB blob; you lose constraints, stats granularity, and join efficiency.
- **Arrays** — first-class typed arrays (`int[]`, `text[]`); GIN-indexable for containment/overlap (`@> && <@`).
- **Range & multirange types** (`int4range`, `tstzrange`, …) — represent intervals with inclusive/exclusive bounds; GiST-index for overlap; power **exclusion constraints** (no overlapping bookings).
- **`numeric`** — exact, arbitrary precision (use for money — never `float`). **`float4/8`** for scientific/approximate.
- **`uuid`** — 16-byte; PG18 adds native `uuidv7()` (timestamp-ordered → better index locality than random v4). **`hstore`** (key-value), **enum**, **`inet`/`cidr`/`macaddr`** (network), **`bytea`**, **`tsvector`/`tsquery`** (full-text), **`citext`** (case-insensitive), and geometric types round out the set.
- **Generated columns** — `GENERATED ALWAYS AS (...)`; PG18 makes **virtual** (compute-on-read) the default, with stored available and logically replicable.

## The extension model

Extensions add types, functions, index methods, and background workers via `CREATE EXTENSION`. They are the PostgreSQL superpower — but on **managed platforms the available set is curated** (RDS/Aurora/Cloud SQL publish allow-lists; verify before designing around one). Inspect with `pg_available_extensions` / `pg_extension`.

Extensions you should know:

| Extension | Purpose |
|---|---|
| **`pg_stat_statements`** | Aggregate query stats — the essential production observability extension (see planner reference) |
| **PostGIS** | The leading open-source spatial database — geometry/geography types, spatial indexes (GiST/SP-GiST), thousands of functions |
| **`pg_trgm`** | Trigram similarity / fuzzy `LIKE '%x%'` and `ILIKE` acceleration via GIN/GiST |
| **`postgres_fdw`** | Query remote PostgreSQL tables as if local (foreign data wrapper); also `file_fdw`, third-party FDWs to other engines |
| **`pgcrypto`** | Hashing/encryption functions |
| **`pg_partman`** | Automated partition lifecycle management (see partitioning reference) |
| **`auto_explain`** | Logs execution plans of slow statements automatically |
| **`pgaudit`** | Detailed session/object audit logging for compliance |
| **TimescaleDB** | Time-series: hypertables (auto-partitioned chunks), continuous aggregates, compression, retention |
| **Citus** | Distributed PostgreSQL: sharded distributed tables + replicated reference tables + a distributed planner for horizontal scale-out |
| **`pgvector`** | Vector similarity search (below) |

## pgvector — PostgreSQL as a vector database

`pgvector` adds a `vector` type and ANN indexes, letting one database hold relational data *and* embeddings — avoiding a separate vector store for many RAG/semantic-search workloads. (Current line **0.8.x**; verify the version your platform offers — feature gates move.)

### Types & distance operators

- Types: `vector` (float32, up to 16,000 dims indexable to ~2,000), **`halfvec`** (float16 — half the storage/memory, supports higher dims), `bit` (binary/Hamming), `sparsevec` (sparse).
- Distance operators (match the one your embedding model was trained with): **`<->`** L2/Euclidean, **`<=>`** cosine distance, **`<#>`** negative inner product, plus `<+>` L1 (and Hamming/Jaccard for bit/sparse).

```sql
CREATE TABLE docs (id bigserial PRIMARY KEY, embedding vector(1536), tenant_id int);
SELECT id FROM docs ORDER BY embedding <=> $1 LIMIT 10;   -- cosine k-NN
```

### Index choice: HNSW vs IVFFlat

- **HNSW** (graph-based) — best recall/latency, builds without seeing data first, but more memory and slower build. The default choice for most workloads. Tune `m`, `ef_construction` at build; `hnsw.ef_search` at query time (higher = better recall, slower).
- **IVFFlat** (inverted lists) — smaller/faster to build but **must be built after representative data is loaded** (it clusters existing rows); recall depends on `lists` (build) and `ivfflat.probes` (query).

```sql
CREATE INDEX ON docs USING hnsw (embedding vector_cosine_ops)
  WITH (m = 16, ef_construction = 64);
```

### Operational notes & RAG patterns

- **ANN is approximate** — the index trades recall for speed; raise `ef_search`/`probes` if results miss. Exact search (no index) is fine for small tables.
- **Filtered/hybrid search**: combine the vector order-by with relational `WHERE` (tenant, ACL, date) and/or full-text (`tsvector`) / `pg_trgm` for hybrid retrieval. pgvector 0.8 improved **iterative index scans** so post-filtering returns enough rows. **Multi-tenant isolation must be enforced by the `WHERE`/RLS predicate** — the index alone is not an access boundary (see `atlas-vector-search-pii-isolation` for the parallel MongoDB threat model and `embedding-inversion-threat-model` for why embeddings aren't anonymization).
- **`halfvec`** roughly halves index size/memory with minimal recall loss — often the better default for high-dimensional embeddings at scale.
- For chunking strategy, reranking, query transformation, and end-to-end RAG architecture, route to **`ai-rag-retrieval`** — this reference owns only the Postgres storage/index layer.

> **Mental model:** PostgreSQL's types and extensions let one engine serve relational, document (JSONB), spatial (PostGIS), time-series (TimescaleDB), and vector (pgvector) workloads — but on managed platforms always confirm the extension allow-list, and remember pgvector's ANN index is a *speed/recall* structure, never a security boundary.
