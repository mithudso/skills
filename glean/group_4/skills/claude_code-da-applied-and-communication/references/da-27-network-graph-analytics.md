<!-- hub-reference-banner -->
> **Reference file — part of the `da-applied-and-communication` hub.** Formerly the standalone `da-27-network-graph-analytics` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-27-network-graph-analytics
description: >-
  Network and graph analytics as a data-analysis discipline — representing data
  as nodes and edges, then measuring structure to answer analytical questions.
  Covers graph representations (adjacency matrix/list, directed/weighted,
  bipartite, ego networks), centrality (degree, betweenness, closeness,
  eigenvector, PageRank), community detection (Louvain, Leiden, label
  propagation, modularity and its resolution limit), connected components,
  shortest paths (BFS/Dijkstra/Bellman-Ford), link prediction (common neighbors,
  Jaccard, Adamic-Adar, preferential attachment), network motifs,
  bipartite projection, graph embeddings (node2vec, DeepWalk), GNN basics
  (GCN, GraphSAGE) for analytics, and tooling (NetworkX, igraph, graph-tool,
  cuGraph, Neo4j GDS). TRIGGER: modeling data as a graph/network; computing
  centrality or ranking influential nodes; detecting communities/clusters in a
  network; choosing Louvain vs Leiden; finding connected components or shortest
  paths; predicting missing edges/links; analyzing bipartite or ego networks;
  generating node embeddings (node2vec/DeepWalk) or using GNNs for analytics;
  picking a graph library (NetworkX/igraph/graph-tool/cuGraph/Neo4j GDS);
  questions about modularity, PageRank, betweenness, or graph scaling. SKIP:
  general relational/SQL joins with no graph structure (use da-13); graph
  database operational/Cypher schema design with no analytics (use Neo4j ops
  skills); pure ML model training unrelated to graphs (use da-7); time-series or
  forecasting (use da-15).
---

# Network & Graph Analytics

## Overview

Network (graph) analytics models data as **nodes (vertices)** connected by
**edges (links)** and measures the resulting structure to answer questions that
row/column tables cannot: who is influential, what clusters exist, what is the
shortest path, what links are likely to form. It is the analytics counterpart to
graph theory — the goal is **insight from relationships**, not just storing them.

Use a graph framing when the *connections* carry the signal: social networks,
fraud rings, supply chains, citation/co-authorship, recommendation, knowledge
graphs, dependency graphs, transaction flows. If the question is answerable with
a `GROUP BY`, you probably do not need a graph.

This skill is the network/graph node of the data-analytics curriculum (da-1 onward).

## Core Concepts

### 1. Graph representations
- **Directed vs undirected**: edges with vs without a direction (following vs friendship). **Weighted vs unweighted**: edges carry a cost/strength.
- **Adjacency matrix**: V×V matrix, O(V²) space, O(1) edge lookup — good for dense graphs and linear-algebra ops (PageRank, spectral methods).
- **Adjacency list**: per-node neighbor lists, O(V+E) space — the default for sparse real-world graphs; faster traversal.
- **Bipartite graph**: two disjoint node sets with edges only across sets (users↔products, authors↔papers).
- **Ego network**: the subgraph of one focal node ("ego"), its direct neighbors ("alters"), and edges among them — the unit of local social-structure analysis.
- **Multigraph / multi-relational**: parallel edges or typed edges (knowledge graphs).

### 2. Connectivity & paths
- **Connected components**: maximal sets of mutually reachable nodes. In directed graphs distinguish **weakly** (ignore direction) vs **strongly** connected components. NetworkX: `connected_components`, `strongly_connected_components`.
- **Shortest paths**: **BFS** for unweighted; **Dijkstra** for non-negative weights (O(E log V) with a heap on an adjacency list); **Bellman-Ford** when negative-weight edges exist (Dijkstra fails on negatives). All-pairs via repeated Dijkstra or Floyd-Warshall.
- **Diameter / eccentricity / average path length**: global reachability measures (expensive on large graphs — sample).

### 3. Centrality (who matters)
- **Degree centrality**: number of edges (in/out for directed) — local popularity, cheap.
- **Betweenness centrality**: fraction of shortest paths passing through a node — identifies bridges/brokers/bottlenecks. Expensive (Brandes' algorithm ≈ O(VE)); **approximate via sampling** on big graphs.
- **Closeness centrality**: inverse of mean shortest-path distance to all others — "how central/quick to reach everyone."
- **Eigenvector centrality**: importance is recursive — you matter if connected to nodes that matter. Can fail to converge on some directed graphs.
- **PageRank**: eigenvector centrality with a **damping factor** (~0.85) modeling a random surfer who teleports. Handles directed graphs reliably; the production default for influence ranking.

### 4. Community detection (what clusters)
- **Modularity (Q)**: scores a partition by edges-inside-communities vs expected at random, range roughly −1..1; higher = stronger community structure.
- **Louvain** (Blondel et al., 2008): fast greedy modularity maximization, two-phase (local move + aggregate). Ubiquitous but suffers the **resolution limit** (merges small real communities) and can produce **badly/dis-connected communities**.
- **Leiden** (Traag, Van Eck & Waltman, 2019): adds a refinement phase; **guarantees communities are connected and well-separated**, faster and higher-quality than Louvain — the recommended default.
- **Label propagation**: near-linear, no objective to optimize — fast but unstable/non-deterministic; good for a quick first pass.
- Modularity is not the only objective — CPM (constant Potts model) and resolution parameters address the resolution limit.

### 5. Link prediction (what edges will form)
Local proximity scores between non-adjacent node pairs x,y (Γ = neighbor set):
- **Common Neighbors**: |Γ(x) ∩ Γ(y)|.
- **Jaccard Coefficient**: |Γ(x) ∩ Γ(y)| / |Γ(x) ∪ Γ(y)| — normalizes by total reach.
- **Adamic-Adar**: sum of 1/log(degree) over shared neighbors — rare shared neighbors count more than hubs.
- **Preferential Attachment**: deg(x)·deg(y) — "rich get richer" growth model.
NetworkX exposes these in `networkx.algorithms.link_prediction`. Embedding/GNN methods (below) are the supervised upgrade.

### 6. Network motifs & bipartite projection
- **Motifs**: small, statistically over-represented subgraphs (e.g., feed-forward loops, triangles) treated as functional building blocks. Compare counts against a degree-preserving null model.
- **Bipartite (one-mode) projection**: collapse a two-set graph onto one set — two authors linked if they co-wrote a paper. Projection **loses information**, so weight edges by shared-neighbor count or hyperbolic/Newman weighting to avoid hub-dominated dense graphs.

### 7. Graph embeddings (nodes → vectors)
- **DeepWalk** (Perozzi et al., 2014): uniform random walks treated as "sentences," fed to skip-gram (Word2Vec) to learn node vectors.
- **node2vec** (Grover & Leskovec, 2016): **biased** random walks with return parameter **p** and in-out parameter **q** interpolating BFS-like (structural roles / homophily) vs DFS-like (community) exploration. Outperforms DeepWalk/LINE on classification and link prediction. Vectors then feed any downstream ML (clustering, classification, link prediction).

### 8. GNN basics for analytics
- **GCN** (Kipf & Welling, 2017): neighborhood feature aggregation via the normalized adjacency; **transductive** — needs the whole graph at train time, must retrain when nodes are added.
- **GraphSAGE** (Hamilton, Ying & Leskovec, NeurIPS 2017): learns *aggregator functions* over a **sampled** neighborhood → **inductive**, generalizes to unseen nodes and scales to large/dynamic graphs. The practical default when the graph grows. Use GNNs for analytics when you have rich node features and a supervised target (node classification, fraud scoring); use node2vec when you only have structure.

## Tools / Frameworks

| Tool | Backend | Best for | Notes |
| --- | --- | --- | --- |
| **NetworkX** | pure Python | prototyping, teaching, graphs up to ~10⁴–10⁵ nodes | richest algorithm coverage, simplest API; 40–250× slower than C libs |
| **igraph** | C (Python/R) | medium-large graphs, single machine | fast, mature, great community/centrality coverage |
| **graph-tool** | C++/Boost + OpenMP | large graphs, parallel centrality/PageRank | fastest CPU lib when OpenMP enabled; statistical inference (SBM) |
| **cuGraph (RAPIDS)** | GPU/CUDA | very large graphs, vertex-centric ops | up to ~870× over igraph; ~0.2s PageRank where igraph takes ~60s |
| **Neo4j GDS** | JVM, in-DB | enterprise graphs in a graph DB, Cypher-driven pipelines | 65+ algorithms (PageRank, Louvain, Leiden, node2vec, FastRP, link prediction); production/beta/alpha tiers |
| **PyG / DGL** | PyTorch | GNN training (GCN, GraphSAGE) | for embeddings and supervised graph ML |

Rule of thumb: prototype in NetworkX, move to igraph/graph-tool when it gets slow, cuGraph when it gets huge, Neo4j GDS when the graph already lives in Neo4j.

## Methodology

1. **Frame the question as a graph** — define what a node is, what an edge is, direction, and weight. Wrong node/edge definition dooms everything downstream.
2. **Build & sanity-check** — node/edge counts, degree distribution (expect heavy tails), number/size of connected components, density. Restrict analysis to the giant component when appropriate.
3. **Pick the analytic to match the question**: influence → centrality (PageRank default); clusters → community detection (Leiden default); reachability → components/shortest paths; missing/future links → link prediction or embeddings.
4. **Scale-match the tool** (see table) before running expensive O(VE) measures.
5. **Validate** — compare against a null model (configuration model / degree-preserving rewiring); for community detection check modularity *and* stability across runs/seeds; for link prediction use temporal train/test split and AUC/precision@k.
6. **Communicate** — layouts for small graphs only (<~1k nodes); for large graphs report metrics, ranked tables, and community summaries rather than hairball plots.

## Practical Patterns

- **PageRank as the default influence score** on directed graphs: degree is cheap but naive; betweenness is informative but slow; PageRank is the reliable middle ground.
- **Leiden over Louvain** for community detection unless you have a hard dependency on Louvain output; you get connectivity guarantees and speed for free.
- **node2vec for structure-only data, GraphSAGE for feature-rich + supervised** targets that need inductive generalization.
- **Work on the giant connected component** — isolates and tiny components distort global metrics and centrality.
- **Approximate expensive centralities** (sampled betweenness/closeness) on graphs over ~10⁵ nodes instead of exact computation.
- **Weight bipartite projections** rather than using raw co-occurrence, or hubs dominate.
- **Tune node2vec p/q deliberately**: low q → DFS/community-flavored embeddings; high q (low p) → BFS/structural-role embeddings.

## Anti-Patterns

- **Treating any join table as a graph.** If connectivity carries no signal and a `GROUP BY` answers the question, a graph adds cost, not insight.
- **Trusting Louvain communities as connected.** They can be internally disconnected; up to ~25% badly connected in the original study. Use Leiden or verify connectivity.
- **Ignoring the modularity resolution limit** — Louvain/modularity merges small genuine communities; do not over-interpret community count without a resolution sweep.
- **Exact betweenness on million-node graphs in NetworkX** — it will not finish; sample or switch to graph-tool/cuGraph.
- **Adjacency matrix for sparse graphs** — O(V²) memory blows up; use adjacency lists.
- **Dijkstra with negative weights** — silently wrong; use Bellman-Ford.
- **Plotting a 100k-node hairball.** Unreadable; summarize with metrics and community-level rollups.
- **Comparing motif/community counts without a null model** — high counts can be a byproduct of the degree distribution, not real structure.
- **Using transductive GCN on a growing graph** — you'll retrain constantly; use GraphSAGE for inductive settings.

## Troubleshooting

- **Eigenvector centrality won't converge** → directed graph with sinks/zero in-degree; use PageRank (damping handles it) or `eigenvector_centrality_numpy`.
- **Everything is one giant community** → resolution limit; lower the resolution parameter, switch to Leiden/CPM, or check the graph isn't accidentally near-complete.
- **Community results change every run** → expected for Louvain/label propagation; fix the seed, run multiple times and take consensus, or use Leiden.
- **Centrality job never finishes** → it's O(VE)-class; sample-approximate, restrict to the giant component, or move to a C/GPU backend.
- **Link prediction AUC ≈ 0.5** → likely no temporal split (leakage) or the graph is too sparse for local proximity; try embedding-based features.
- **node2vec embeddings look random** → walks too short / too few, or p/q untuned; increase walk length and number of walks, retune p/q.
- **Out of memory building the graph** → you used a dense matrix; switch to an edge list / sparse (CSR) representation or igraph/cuGraph.

## References

- NetworkX docs — centrality, components, shortest paths, link prediction. https://networkx.org/documentation/stable/ (2024)
- Brandes, U. "A Faster Algorithm for Betweenness Centrality." J. Math. Sociology (2001).
- Page, Brin et al. "The PageRank Citation Ranking." Stanford (1999).
- Blondel et al. "Fast unfolding of communities in large networks" (Louvain). J. Stat. Mech. (2008). https://arxiv.org/abs/0803.0476
- Fortunato & Barthélemy. "Resolution limit in community detection." PNAS (2007). https://arxiv.org/abs/physics/0607100
- Traag, Van Eck & Waltman. "From Louvain to Leiden: guaranteeing well-connected communities." Scientific Reports (2019). https://arxiv.org/abs/1810.08473
- Liben-Nowell & Kleinberg. "The Link Prediction Problem for Social Networks." (2007). https://www.cs.cornell.edu/home/kleinber/link-pred.pdf
- Bipartite network projection. Wikipedia / Arthur, "Modularity and Projection of Bipartite Networks" (2019). https://arxiv.org/pdf/1908.02520
- Perozzi, Al-Rfou & Skiena. "DeepWalk: Online Learning of Social Representations." KDD (2014).
- Grover & Leskovec. "node2vec: Scalable Feature Learning for Networks." KDD (2016). https://cs.stanford.edu/~jure/pubs/node2vec-kdd16.pdf
- Kipf & Welling. "Semi-Supervised Classification with Graph Convolutional Networks." ICLR (2017).
- Hamilton, Ying & Leskovec. "Inductive Representation Learning on Large Graphs" (GraphSAGE). NeurIPS (2017). https://cs.stanford.edu/people/jure/pubs/graphsage-nips17.pdf
- Neo4j Graph Data Science docs — centrality, community detection, embeddings, link prediction. https://neo4j.com/docs/graph-data-science/current/ (2024)
- igraph documentation. https://igraph.org/ (2024)
- graph-tool performance. https://graph-tool.skewed.de/performance.html (2024)
- RAPIDS cuGraph. https://docs.rapids.ai/api/cugraph/stable/ (2024)
- Benchmark of popular graph/network packages (NetworkX/igraph/graph-tool/cuGraph). https://www.timlrx.com/blog/benchmark-of-popular-graph-network-packages-v2/ (2020)
