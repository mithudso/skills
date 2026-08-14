<!-- hub-reference-banner -->
> **Reference file — part of the `da-data-engineering-platform` hub.** Formerly the standalone `da-41-knowledge-graphs-and-semantic-analytics` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-41-knowledge-graphs-and-semantic-analytics
description: >-
  Knowledge graphs and semantic analytics as a data-analysis discipline — the
  SEMANTIC / ONTOLOGY / meaning angle of graph data, distinct from graph
  algorithms. Covers KG fundamentals (nodes/edges/properties, labeled property
  graph (LPG) vs RDF triple store, RDF-star, when a KG beats relational/
  dimensional modeling); the semantic web stack (RDF, RDFS, OWL 2 ontologies,
  SPARQL, SHACL validation, named graphs); ontology & taxonomy engineering
  (classes/relations, upper ontologies, SKOS controlled vocabularies, ontology
  design patterns); KG construction (entity extraction, entity resolution &
  linking, relation extraction, schema mapping with R2RML/RML, ontology-based
  data access / virtual KGs, and LLM-assisted KG construction 2024-2026);
  querying & analytics (SPARQL vs Cypher vs ISO GQL, semantic queries,
  reasoning / inference / materialisation, graph analytics over a KG);
  enterprise KGs (data fabric, metadata knowledge graph, semantic layer,
  data catalogs as KGs); GraphRAG & semantic retrieval for grounded,
  entity-centric analytics; and tooling (Neo4j + neosemantics/n10s, RDFLib,
  Apache Jena/Fuseki, Ontotext GraphDB, Stardog, Amazon Neptune, TigerGraph,
  Virtuoso, Wikidata, schema.org). TRIGGER: designing or modeling a knowledge
  graph; choosing LPG vs RDF/triple store or RDF-star; writing an ontology or
  taxonomy (OWL/RDFS/SKOS); SPARQL vs Cypher vs GQL choice; SHACL validation;
  reasoning/inference over a KG; constructing a KG from text/relational sources;
  entity resolution/linking into a KG; LLM-assisted KG construction; GraphRAG
  or semantic retrieval for analytics/Q&A; enterprise/metadata knowledge graph,
  data fabric, or data catalog as a graph; picking a triple store or KG
  platform; "when should I use a knowledge graph". SKIP: graph ALGORITHMS —
  centrality, community detection, link prediction, GNNs, node embeddings for
  analytics (use da-27-network-graph-analytics); vector databases, agent
  memory, or KG-as-AI-memory storage angle (use ai-datastores); the
  metrics/headless-BI semantic layer of governed business metrics with no
  ontology/graph angle (use da-18-semantic-layer-headless-bi); generic data
  catalogs/governance with no graph modeling (use da-30-data-governance-catalogs);
  Kimball/dimensional warehouse modeling (use da-29-dimensional-data-modeling).
---

# Knowledge Graphs & Semantic Analytics

## Overview

A **knowledge graph (KG)** represents entities (nodes) and the typed,
meaning-bearing relationships between them (edges), with attributes (properties)
on both, plus a **schema/ontology** that says what the types *mean*. The point
is not just to store connections (that is da-27's graph-algorithms angle) but to
encode **semantics** — shared, machine-interpretable meaning — so that data from
many sources can be integrated, queried by meaning, validated against a model,
and reasoned over to infer new facts.

Reach for a KG when the **connections and their meaning carry the signal and
must be queried, integrated, or reasoned about**: multi-hop questions ("which
suppliers two tiers down depend on a sanctioned entity?"), heterogeneous data
integration under one vocabulary, provenance/lineage, regulatory traceability,
and grounding LLMs. A KG beats relational/dimensional modeling when traversal
depth is variable and deep: relationships are stored as explicit, first-class
connections, so following one costs roughly the same on 1k or 100M nodes,
whereas a relational model pays an exponential JOIN cost as hop depth grows and
the schema must be known up front. If the question is answerable with a
`GROUP BY` or a couple of JOINs over a stable schema, you do **not** need a KG —
use a warehouse/dimensional model (da-29).

This is the semantic/ontology node of the data-analytics curriculum (da-1 onward).

## Scope boundary (read this first)

- **da-27-network-graph-analytics** owns graph *algorithms*: centrality,
  community detection, link prediction, shortest paths, node embeddings, GNNs.
  This skill owns *meaning*: ontologies, RDF/OWL/SHACL, semantic queries,
  reasoning, KG construction, GraphRAG. Running PageRank *over* a KG is da-27;
  defining what the nodes mean is here.
- **ai-datastores ("Knowledge Graphs for AI")** owns the vector-DB / agent-memory
  / KG-as-storage angle. This skill owns the analytics/semantic-integration angle.
- **da-18-semantic-layer-headless-bi** owns the *metrics* layer (dbt SL, Cube,
  MetricFlow) — governed business measures. A KG semantic layer is about
  *entities and their meaning*; the two are complementary, not the same.
- **da-30-data-governance-catalogs** owns governance/catalog policy generally;
  this skill covers the specific pattern of *modeling the catalog itself as a
  knowledge graph*.

## Core Concepts

### 1. Two graph data models: LPG vs RDF triple store
- **Labeled Property Graph (LPG)**: nodes/edges carry labels and arbitrary
  key/value **properties**; both nodes and individual relationships can hold
  attributes. Optimized for fast traversal, expressive analytics, embeddings.
  Native to Neo4j, TigerGraph, Memgraph. Query with Cypher / ISO GQL.
- **RDF triple store**: everything is a `(subject, predicate, object)` **triple**
  with global URIs/IRIs; emphasizes shared meaning, W3C standards, external
  vocabulary reuse, and logical inference. Query with SPARQL. Native to GraphDB,
  Jena, Stardog, Virtuoso; Neptune supports both models.
- **Key limitation of plain RDF**: you cannot attach attributes to a single
  relationship instance (multiple same-type edges collapse to one triple).
  **RDF-star (RDF\*)** fixes this by letting a triple be the subject/object of
  another triple — bringing LPG-style edge properties to RDF; **SPARQL-star**
  queries it. RDF-star reached W3C Recommendation track maturity (RDF 1.2).
- **Choosing**: RDF for cross-org interoperability, regulatory/semantic
  reasoning, vocabulary reuse, and provenance; LPG for speed, deep traversal,
  graph algorithms, and AI/embedding workflows. Hybrid is common (RDF as the
  governed model of record; project into an LPG for fast traversal/analytics).

### 2. The semantic web stack (W3C)
- **RDF**: the triple data model; serializations: Turtle, N-Triples, JSON-LD,
  RDF/XML. **IRIs** give every thing a global identifier.
- **RDFS**: lightweight schema: `rdfs:Class`, `rdfs:subClassOf`,
  `rdfs:domain`/`rdfs:range`, property hierarchies.
- **OWL 2**: rich ontology language with formal Description-Logic semantics
  (classes, restrictions, cardinality, transitivity, disjointness, equivalence).
  Profiles **OWL 2 EL / QL / RL** trade expressivity for tractable reasoning.
- **SPARQL**: W3C query language: graph-pattern matching, `OPTIONAL`, property
  paths (variable-length traversal), federated queries (`SERVICE`), `CONSTRUCT`
  to build new graphs, aggregation.
- **SHACL** (Shapes Constraint Language): W3C *validation* standard: declare
  shapes (e.g. every `Person` has exactly one `birthDate`) and validate a graph.
  Closed-world validation, complementary to OWL's open-world inference.
- **Named graphs**: group a set of triples under a graph IRI (a "quad") to
  manage provenance, trust, versioning, and access control on subsets.

### 3. OWL vs SHACL (don't confuse them)
- **OWL = inference (open-world)**: derives *new* facts. "X is a Manager and
  Manager ⊑ Employee, therefore X is an Employee." Absence of a fact ≠ false.
- **SHACL = validation (closed-world)**: checks whether the data *as it stands*
  meets constraints and reports violations. Absence ≠ allowed.
- Modern practice uses **OWL for modeling + SHACL for validation** together; a
  2025 KCAP paper documents lessons from combining them on one project.

### 4. Ontology & taxonomy engineering
- **Taxonomy** = hierarchical classification (broader/narrower). **Ontology** =
  taxonomy + non-hierarchical typed relations + axioms/constraints.
- **SKOS** (Simple Knowledge Organization System): W3C standard for controlled
  vocabularies/thesauri/taxonomies: `skos:Concept`, `broader`/`narrower`,
  `related`, `prefLabel`/`altLabel`, concept schemes. Use SKOS for vocabularies,
  OWL for formal domain models.
- **Upper (top-level) ontologies**: domain-independent category systems (BFO,
  DOLCE, SUMO, gist, schema.org as a light upper vocab) that align separate
  domain ontologies for integration.
- **Ontology Design Patterns (ODPs)**: reusable modeling templates (like
  software design patterns) for recurring problems; improve quality, reuse,
  maintainability.
- **schema.org**: pragmatic web vocabulary (JSON-LD/Microdata/RDFa) behind
  Google's Knowledge Graph and SEO structured data.

### 5. KG construction (the hard part)
- **Schema-first vs data-first**: design the ontology up front (ontology-oriented,
  e.g. OntoKG) vs extract then align. LLM-assisted pipelines increasingly do both.
- **Entity extraction / NER** → identify entities in unstructured text.
- **Relation extraction** → identify typed links between entities.
- **Entity resolution / deduplication** → merge records referring to the same
  real entity (correlation clustering, embedding similarity, density partitioning).
- **Entity linking** → map a mention to a canonical KG identifier (e.g. a
  Wikidata Q-number); zero-shot LLM linkers re-rank candidates by log-prob.
- **Schema mapping from relational sources**: **R2RML** / **RML** declaratively
  map RDBMS rows or CSV/JSON to RDF triples. **Ontology-Based Data Access
  (OBDA) / virtual KGs** expose relational data *as* a KG at query time without
  materializing (e.g. Ontop), rewriting SPARQL to SQL.
- **LLM-assisted construction (2024-2026)**: LLMs (DeepSeek-R1, Qwen3, GLM-4.5,
  Claude, GPT) perform NER, relation extraction, schema/ontology generation,
  and entity disambiguation; agentic pipelines orchestrate ontology creation →
  population with minimal manual schema work. Validate LLM output via OWL/Turtle
  formalization + SHACL/reasoning, since LLM extraction is approximate.

### 6. Querying, reasoning & analytics
- **SPARQL** (RDF) vs **Cypher** (LPG) vs **GQL** (ISO/IEC 39075:2024, the first
  new ISO query-language standard since SQL in 1987 — a standardized property-graph
  language inspired by Cypher; `INSERT` vs Cypher's `CREATE`, `FOR` vs `UNWIND`).
- **Reasoning / inference / materialisation**: an OWL reasoner derives entailed
  triples (e.g. subclass, transitive, inverse) and can **materialise** them as
  explicit triples — *logically sound*, unlike approximate embedding-based KG
  completion. Choose **batch/forward-chaining materialisation** vs **on-the-fly
  backward-chaining** by write/read pattern.
- **Graph analytics over a KG**: you *can* run centrality/community/path
  algorithms over a KG, but that math lives in **da-27**. Here the analytics are
  *semantic*: entity-centric aggregation, multi-hop semantic joins, lineage
  traversal, "explain the path" queries.

### 7. Enterprise KGs, data fabric & catalogs-as-graphs
- **Data fabric** provides connectivity + governance plumbing; the **knowledge
  graph** provides the semantic/relational intelligence — complementary, not
  competing. Gartner elevated the semantic layer to essential BI infrastructure
  in the 2025 Hype Cycle.
- **Metadata knowledge graph**: model technical/business/operational metadata
  (lineage, ownership, classifications) as a KG → a unified, queryable catalog
  for governance, compliance, discovery (Informatica, SAP Business Data Cloud /
  Datasphere knowledge graph, 2025).
- **Data catalog as a KG**: represent assets + their relationships as a graph so
  discovery becomes traversal/semantic search rather than keyword lookup.

### 8. GraphRAG & semantic retrieval for analytics
- **GraphRAG** = RAG over a knowledge graph: retrieve via graph structure (and
  often hybrid vector + graph), not just nearest-neighbor chunks. Microsoft
  open-sourced GraphRAG (2024); vendors (ServiceNow, Workday) integrated it.
- **Why it matters for analytics**: vector-only RAG can hit ~0% on schema-bound
  KPI/forecast questions; optimized GraphRAG reaches 90%+ — it handles multi-hop
  and "global" questions and gives **explainable, grounded, entity-centric**
  answers with citations/lineage.
- **Cost/tradeoff**: KG extraction can cost 3–5× baseline RAG and needs
  domain tuning. The 2025-2026 production pattern is **hybrid routing**: send each
  query to vector / KG-traversal / structured-SQL backend by query type; agentic
  GraphRAG (StepChain, GraphSearch, Youtu-GraphRAG) plans multi-hop traversals.

## Tools / Frameworks

| Tool | Model | Notes |
| --- | --- | --- |
| **Neo4j + neosemantics (n10s)** | LPG (+RDF import/export) | Dominant LPG; Cypher 5 / GQL; n10s bridges RDF/OWL into the property graph for hybrid projects |
| **RDFLib** | RDF | Python library for parsing/serializing/querying RDF + SPARQL; prototyping/ETL |
| **Apache Jena / Fuseki** | RDF | Java triple store + SPARQL server; reasoners; reference open-source stack |
| **Ontotext GraphDB** | RDF | De-facto standard for heavy semantic reasoning + SHACL validation at scale (v10.6, 2025) |
| **Stardog** | RDF | Enterprise KG platform: reasoning, ontology management, virtual graphs/OBDA |
| **Amazon Neptune** | RDF **and** LPG | Managed; Gremlin/openCypher/SPARQL; Serverless v2 (2025) cuts idle cost |
| **TigerGraph** | LPG | High-scale parallel graph analytics; GSQL |
| **Virtuoso (OpenLink)** | RDF | Fast SPARQL 1.1 + SQL federation in one binary (v8, 2025) |
| **Ontop** | Virtual RDF (OBDA) | SPARQL-to-SQL over existing relational DBs without materializing |
| **Wikidata** | RDF | World's largest open KG (CC0); ~1.65B statements; Q/P identifiers; entity-linking target |
| **schema.org** | Vocabulary | JSON-LD/Microdata web vocabulary; SEO + light upper ontology |

## Methodology

1. **Decide if you even need a KG**: variable-depth traversal, multi-source
   integration, reasoning, provenance, or LLM grounding? If a few stable JOINs
   suffice, use da-29 instead.
2. **Pick the model**: RDF (interoperability/reasoning/standards) vs LPG
   (speed/analytics/AI); consider RDF-star or a hybrid RDF-of-record + LPG-projection.
3. **Design the ontology/taxonomy**: reuse SKOS/schema.org/upper ontologies and
   ODPs before inventing terms; keep OWL (modeling) and SHACL (validation) distinct.
4. **Construct**: extract entities/relations, resolve & link to canonical IDs,
   map relational sources via R2RML/RML or expose virtually via OBDA; use
   LLM-assisted extraction but validate output with SHACL/reasoning.
5. **Validate & reason**: run SHACL for data quality; materialise or query-time
   reason for entailed facts.
6. **Query & serve analytics**: SPARQL/Cypher/GQL; layer GraphRAG for grounded
   Q&A; route hybrid (vector/KG/SQL) by query type.

## Practical Patterns

- **RDF of record, LPG for speed**: govern + reason in RDF, project to Neo4j for
  fast traversal and to run da-27 algorithms.
- **Validate LLM-extracted triples** with SHACL before committing: treat LLM
  output as approximate candidates, not ground truth.
- **Reuse before you model**: pull terms from schema.org, SKOS thesauri, and
  domain/upper ontologies; bespoke vocabularies kill interoperability.
- **Materialise hot inferences, backward-chain the rest** to balance write cost
  vs query latency.
- **Hybrid GraphRAG routing**: classify the query, then pick vector / KG / SQL;
  return the traversal path as the citation for explainability.

## Anti-Patterns

- **Using a KG as a slow relational DB**: if every query is a 1–2 JOIN over a
  fixed schema, a graph adds cost without payoff.
- **Confusing OWL and SHACL**: expecting SHACL to infer, or OWL to "reject" data.
  OWL infers (open-world); SHACL validates (closed-world).
- **Plain RDF for edge attributes**: needing per-relationship properties in RDF
  and not using RDF-star (or an LPG) leads to ugly reification.
- **Skipping entity resolution**: duplicate, unlinked entities make traversal
  and analytics silently wrong.
- **Ontology over-engineering**: a 500-class OWL model nobody populates; start
  small, pattern-driven, and grow.
- **Trusting unvalidated LLM extractions**: no SHACL/reasoning gate → confident
  but wrong triples poisoning downstream analytics and GraphRAG.

## Troubleshooting

- **SPARQL query times out on deep traversal** → use property paths sparingly;
  consider materialising inferred closure; or project to an LPG for the traversal.
- **Reasoner explodes / never finishes** → you're outside a tractable OWL profile;
  drop to OWL 2 RL/QL/EL or move constraints to SHACL.
- **GraphRAG answers are vague/global-only** → entity resolution/linking is weak,
  or you're retrieving subgraphs too coarsely; tighten linking + multi-hop planning.
- **Integrated data is "connected" but meaningless** → no shared ontology; you
  joined IDs without aligning semantics — introduce a common vocabulary.
- **Same real-world entity appears many times** → entity-resolution step missing
  or under-tuned; add embedding/clustering-based dedup before load.

## References

- Neo4j — RDF Triple Stores vs. Property Graphs (knowledge-graph blog), 2025: https://neo4j.com/blog/knowledge-graph/rdf-vs-property-graphs-knowledge-graphs/
- TigerGraph — RDF vs Property Graph: Choosing the Right Foundation, 2025: https://www.tigergraph.com/blog/rdf-vs-property-graph-choosing-the-right-foundation-for-knowledge-graphs/
- Enterprise Knowledge — Intro to RDF & LPG Graphs: https://enterprise-knowledge.com/cutting-through-the-noise-an-introduction-to-rdf-lpg-graphs/
- AI Business — Ending the RDF vs Property Graph debate with RDF-star: https://aibusiness.com/data/ending-the-rdf-vs-property-graph-debate-with-rdf-
- GraphDB 11.x — W3C specifications (RDF/RDFS/OWL/SPARQL/SHACL): https://graphdb.ontotext.com/documentation/11.2/w3c-specs.html
- Atlan — RDF vs OWL: key differences & use cases: https://atlan.com/know/rdf-vs-owl/
- KCAP 2025 — Lessons Learned from Combined Development of OWL and SHACL: https://dl.acm.org/doi/full/10.1145/3731443.3771340
- Modern Data 101 — Demystifying SKOS for Practitioners: https://moderndata101.substack.com/p/demystifying-skos-for-practitioners
- Hedden Information Management — SKOS Taxonomies: https://www.hedden-information.com/skos-taxonomies/
- The New Stack — GQL: A New ISO Standard for Querying Graph Databases: https://thenewstack.io/gql-a-new-iso-standard-for-querying-graph-databases/
- AWS — GQL: the ISO standard for graphs has arrived: https://aws.amazon.com/blogs/database/gql-the-iso-standard-for-graphs-has-arrived/
- Wikipedia — Graph Query Language (GQL): https://en.wikipedia.org/wiki/Graph_Query_Language
- arXiv 2510.20345 — LLM-empowered knowledge graph construction: A survey (2025): https://arxiv.org/html/2510.20345v1
- SiliconFlow — Best Open Source LLMs for Knowledge Graph Construction in 2026: https://www.siliconflow.com/articles/en/best-open-source-LLM-for-Knowledge-Graph-Construction
- LLM-TEXT2KG 2026 Workshop (Text2KG): https://aiisc.ai/text2kg2026/
- arXiv 2604.02618 — OntoKG: Ontology-Oriented Knowledge Graph Construction: https://arxiv.org/html/2604.02618v1
- Springer — Ontology-Based Update in Virtual Knowledge Graphs via Schema Mapping Recovery (2024): https://link.springer.com/chapter/10.1007/978-3-031-72407-7_6
- Emergent Mind — Knowledge Graph Construction Techniques: https://www.emergentmind.com/topics/knowledge-graph-construction
- Sage Journals — Logic & Reasoning in Neurosymbolic Systems Using OWL-Based KGs (2025): https://journals.sagepub.com/doi/10.1177/29498732251320043
- VLDB vol.17 — Efficient Validation of SHACL Shapes with Reasoning: https://www.vldb.org/pvldb/vol17/p3589-acosta.pdf
- IBM — What is GraphRAG?: https://www.ibm.com/think/topics/graphrag
- Meilisearch — What is GraphRAG: Complete guide [2025]: https://www.meilisearch.com/blog/graph-rag
- Salfati Group — Graph RAG Guide 2025: Architecture, Implementation & ROI: https://salfati.group/topics/graph-rag
- arXiv 2510.02827 — StepChain GraphRAG: Multi-Hop QA over KGs: https://arxiv.org/html/2510.02827v1
- arXiv 2509.22009 — GraphSearch: Agentic Deep Searching for Graph RAG: https://arxiv.org/pdf/2509.22009
- NStarX — The Next Frontier of RAG: Enterprise Knowledge Systems 2026-2030: https://nstarxinc.com/blog/the-next-frontier-of-rag-how-enterprise-knowledge-systems-will-evolve-2026-2030/
- Informatica — Why a Metadata Knowledge Graph Is Essential to a Data Fabric: https://www.informatica.com/blogs/why-a-metadata-knowledge-graph-is-essential-to-an-intelligent-data-fabric.html
- Enterprise Knowledge — Graph Analytics in the Semantic Layer: https://enterprise-knowledge.com/graph-analytics-in-the-semantic-layer-architectural-framework-for-knowledge-intelligence/
- Coalesce — Semantic Layers in 2025: Catalog Owner & Data Leader Playbook: https://coalesce.io/data-insights/semantic-layers-2025-catalog-owner-data-leader-playbook/
- AI CERTs — Enterprise Data Fabric Powers SAP Datasphere Knowledge Graph: https://www.aicerts.ai/news/enterprise-data-fabric-powers-sap-datasphere-knowledge-graph/
- TechTarget — 5 knowledge graph use cases in data fabric architecture: https://www.techtarget.com/searchdatamanagement/tip/5-knowledge-graph-use-cases-in-data-fabric-architecture
- Kyvos — Survey of the Knowledge Graph Ecosystem: https://www.kyvosinsights.com/blog/survey-of-the-knowledge-graph-ecosystem/
- Galaxy — Top 10 Stardog Alternatives for Enterprise KGs in 2025: https://www.getgalaxy.io/resources/stardog-alternatives
- Medium (Vishal Mysore) — Graph Databases & Query Languages in 2025: A Practical Guide: https://medium.com/@visrow/graph-databases-query-languages-in-2025-a-practical-guide-39cb7a767aed
- Wikipedia — Wikidata: https://en.wikipedia.org/wiki/Wikidata
- Kloia — Knowledge Base vs Knowledge Graph for LLM Systems (2026 Guide): https://www.kloia.com/blog/knowledge-base-vs-knowledge-graph-llm
- PuppyGraph — Property graph vs RDF: key differences, working, use cases: https://www.puppygraph.com/blog/property-graph-vs-rdf
