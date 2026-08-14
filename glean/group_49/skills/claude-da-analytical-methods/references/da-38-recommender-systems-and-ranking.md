<!-- hub-reference-banner -->
> **Reference file — part of the `da-analytical-methods` hub.** Formerly the standalone `da-38-recommender-systems-and-ranking` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-38-recommender-systems-and-ranking
description: >-
  Recommender systems and learning-to-rank as a data-analysis discipline —
  predicting and ordering items a user will engage with, then evaluating those
  predictions offline and online. Covers problem framing (explicit vs implicit
  feedback, the feedback loop / popularity bias, user/item/system cold-start),
  collaborative filtering (user-user, item-item, neighborhood methods), matrix
  factorization (SVD, funkSVD, ALS, implicit-feedback WRMF, BPR pairwise),
  content-based and hybrid recommenders, factorization machines (FM, FFM,
  DeepFM), modern deep recommenders (two-tower retrieval, neural CF, sequential
  GRU4Rec/SASRec/BERT4Rec, transformers-for-rec) and LLM-augmented / generative
  recommenders with semantic IDs (2025-2026), learning-to-rank (pointwise,
  pairwise RankNet, listwise ListNet/LambdaMART) and the retrieval→ranking→
  re-ranking funnel, offline metrics (Precision@k, Recall@k, MAP, MRR, NDCG, hit
  rate, coverage, diversity, novelty, serendipity) and the offline/online gap,
  online evaluation (A/B testing, interleaving, counterfactual/off-policy IPS &
  doubly-robust), and production concerns (candidate generation vs ranking,
  feature stores, embedding/vector stores, exploration/exploitation bandits,
  fairness & filter bubbles). Tooling: implicit, LightFM, Surprise, TorchRec,
  NVIDIA Merlin/NVTabular, RecBole, Vespa, Feast, Amazon Personalize, Vertex AI.
  Part of the da-* data-analytics curriculum.
  TRIGGER: building/choosing a recommender or personalization model; collaborative
  filtering, matrix factorization, ALS, BPR, WRMF, funkSVD; factorization machines
  / DeepFM; two-tower retrieval, neural CF, SASRec/BERT4Rec/GRU4Rec sequential
  recsys, transformers for recommendation; LLM-augmented or generative
  recommenders, semantic IDs, RQ-VAE for recsys; learning-to-rank (RankNet,
  LambdaMART, ListNet, pairwise/listwise), the retrieve-then-rank funnel;
  ranking/recsys offline metrics (NDCG, MAP, MRR, Recall@k, coverage, diversity,
  serendipity) or the offline/online metric gap; A/B testing or interleaving a
  recommender; off-policy / counterfactual evaluation (IPS, doubly-robust) for
  ranking; candidate generation vs ranking architecture; exploration/exploitation
  or contextual bandits for recsys; popularity bias, filter bubbles, fairness in
  recommendation; picking implicit/LightFM/Surprise/TorchRec/Merlin/RecBole/Vespa/
  Personalize/Vertex.
  SKIP: generic supervised ML model training/tuning with no ranking or
  recommendation target → da-7-machine-learning; graph/network descriptive
  metrics, centrality, community detection, GNN-as-graph-analytics → da-27-network-
  graph-analytics (note: GNN-for-recsys link prediction lives here at the framing
  level but graph-structure analysis belongs to da-27); A/B test design and causal
  inference in the general case (DiD, RDD, IV, PSM) with no ranking policy →
  da-12-ab-testing-causal-inference; semantic/vector search or RAG retrieval with
  no personalization/ranking-quality objective → rag-architecture; feature-store
  engineering in the abstract → da-17-feature-engineering-and-feature-stores.
license: Internal TAM curriculum content.
---

# da-38 — Recommender Systems & Learning-to-Rank Analytics

## Overview

A **recommender system** predicts, for each user, which items from a (often
huge) catalog they are most likely to engage with, then **orders** a small slate
to show. As a data-analysis discipline it sits at the intersection of three
problems: (1) *modeling* preference from sparse interaction data, (2) *ranking*
candidates under a relevance objective, and (3) *evaluating* both offline and
online while fighting the bias the system itself creates.

**Scope boundaries (read first):**

- **vs da-7 (machine learning):** da-7 covers generic supervised/unsupervised
  model training, regularization, and hyperparameter tuning. da-38 covers the
  *recommendation- and ranking-specific* objectives (BPR/WRMF losses, NDCG-aware
  LambdaMART), the retrieve-then-rank funnel, and recsys evaluation. Use da-7 for
  "train a classifier"; use da-38 for "rank items for a user."
- **vs da-27 (network/graph analytics):** da-27 owns graph structure metrics
  (centrality, community detection, link prediction *as graph analysis*, GNNs as
  graph models). da-38 references GNN-for-recsys and bipartite user-item graphs
  only at the *recommendation* level; graph-structure questions belong to da-27.
- **vs da-12 (A/B testing & causal inference):** da-12 owns general experiment
  design and causal estimators. da-38 owns the *recsys-specific* online-eval
  wrinkles: interleaving, feedback-loop confounding, and off-policy/counterfactual
  evaluation (IPS, doubly-robust) of a ranking policy.
- **vs rag-architecture / vector search:** semantic retrieval with no
  personalization or ranking-quality objective is RAG; ranked personalization is
  da-38.

---

## Core Concepts

### 1. Problem framing: feedback, the feedback loop, and cold-start

- **Explicit vs implicit feedback.** *Explicit* = ratings/likes the user
  deliberately gives (sparse, low-volume, but signed). *Implicit* = clicks,
  views, dwell, purchases, plays (abundant but **positive-only and ambiguous** —
  a non-click is not a dislike). Implicit feedback dominates production systems
  and forces *positive-unlabeled* modeling: you model **confidence** that an
  observed interaction is a preference, not a rating value.
- **Popularity bias & the feedback loop.** A deployed recommender only logs
  feedback on items it *showed*, which it chose because they scored high — so the
  next training set over-represents popular/previously-recommended items
  (**exposure bias**). Left unchecked this is a self-reinforcing loop that
  narrows the catalog and creates **filter bubbles / echo chambers**. Breaking it
  requires exploration and/or de-biasing (propensity weighting).
- **Cold-start, three flavors.** *User cold-start* (new user, no history),
  *item cold-start* (new item, no interactions), *system cold-start* (brand-new
  product with no data). Mitigations: content/side features, hybrid models,
  factorization machines, popularity fallbacks, bandit exploration, and (2025-26)
  LLM/semantic-ID content grounding for long-tail and brand-new items.

### 2. Collaborative filtering (CF) — neighborhood methods

CF predicts preference purely from the user-item interaction matrix, no content.

- **User-user CF:** find users with similar taste (cosine / Pearson similarity on
  co-rated items), recommend what neighbors liked. Sensitive to sparsity; user
  base churns, so similarities go stale.
- **Item-item CF:** precompute item-item similarity ("users who interacted with X
  also interacted with Y"). More stable than user-user (item catalogs change
  slower), scales better, and was the Amazon production workhorse. Still the right
  baseline for many problems.
- **Limitations:** cold-start, sparsity, popularity skew, no use of side
  features — which motivates latent-factor models.

### 3. Matrix factorization (MF)

Factor the sparse interaction matrix R ≈ **P·Qᵀ** into low-rank user (P) and item
(Q) latent factors; a prediction is the dot product pᵤ·qᵢ.

- **SVD / funkSVD.** "SVD" in recsys (Simon Funk, Netflix Prize) is not literal
  SVD — it learns factors by **SGD on observed entries only** with L2
  regularization, plus user/item **bias terms**. SVD++ adds implicit signal.
- **ALS (Alternating Least Squares).** Fix P, solve Q in closed form, alternate.
  Embarrassingly parallel → the standard for large, distributed (Spark MLlib)
  training.
- **Implicit-feedback MF / WRMF (Hu-Koren-Volinsky 2008).** The canonical
  implicit model: treat all entries, weight observed interactions by a
  **confidence** cᵤᵢ = 1 + α·rᵤᵢ, fit preference (0/1) with weighted ALS. This is
  what the `implicit` library's `AlternatingLeastSquares` implements.
- **BPR — Bayesian Personalized Ranking (Rendle 2009).** A **pairwise ranking**
  objective for implicit feedback: maximize the probability that an observed item
  ranks above an unobserved one, σ(x̂ᵤᵢ − x̂ᵤⱼ), trained with SGD over sampled
  (user, pos, neg) triples. Optimizes ranking (AUC), not rating error — usually a
  better fit for top-N than pointwise MF.

### 4. Content-based, hybrid, and factorization machines

- **Content-based:** recommend items similar to ones the user liked, using item
  features (text, tags, embeddings) → great for item cold-start, but
  over-specializes (no serendipity).
- **Hybrid:** combine CF + content (weighted, switching, feature-augmented,
  cascade). Solves cold-start while keeping CF's collaborative signal.
- **Factorization Machines (FM, Rendle 2010):** model all pairwise feature
  interactions with low-rank factorized weights — generalizes MF to arbitrary
  side features (user, item, context), so it natively handles cold-start.
  **FFM (field-aware)** gives each feature a separate latent vector per *field*
  it interacts with — strong for CTR. **DeepFM** shares an embedding between a
  wide FM component (low-order interactions) and a deep DNN (high-order), needs
  no manual feature crosses, and is a standard CTR/ranking model.

### 5. Modern deep recommenders

- **Two-tower retrieval.** Separate **user/query tower** and **item tower** map to
  a shared embedding space; relevance = dot product. Item embeddings are
  precomputed and indexed in an ANN/vector store, so retrieval is a sub-linear
  nearest-neighbor lookup → the dominant **candidate-generation** architecture at
  scale (and TorchRec's headline target).
- **Neural CF (NCF):** replace the MF dot product with an MLP over user/item
  embeddings — more expressive interactions (though a well-tuned dot product is a
  surprisingly strong baseline).
- **Sequential / session-based.** Model the *order* of interactions:
  - **GRU4Rec** — RNN over session events.
  - **SASRec** — *causal* (left-to-right) self-attention; predicts the next item.
  - **BERT4Rec** — *bidirectional* self-attention with masked-item (cloze)
    training. **Caveat (current research):** when trained with the *same* loss,
    SASRec generally matches or beats BERT4Rec at lower cost — BERT4Rec's reported
    edge often came from its training objective, not bidirectionality. Don't
    assume "bidirectional = better."
- **LLM-augmented & generative recommenders (2025-2026).** The fast-moving
  frontier:
  - **Semantic IDs**: quantize an item's content embedding (RQ-VAE) into a short
    code; these IDs reflect content, so they generalize to **cold-start /
    long-tail** items and slot into an LLM vocabulary. YouTube reported gains
    swapping video-ID embeddings for semantic IDs; Spotify uses them to unify
    search + recommendation.
  - **Generative retrieval**: an LLM/seq2seq model *generates* the semantic ID of
    the next item instead of scoring a candidate set.
  - **LLMs as data augmenters / rerankers** and knowledge-guided RAG (e.g.
    ColdRAG) for cold-start. Watch for LLM-reranker **exposure/coverage** issues.

### 6. Learning-to-rank (LTR)

Directly optimize the *order* of a candidate list given relevance labels/features.

- **Pointwise** — predict each item's relevance independently (regression/
  classification); ignores list context. Simplest, weakest ranking fidelity.
- **Pairwise** — learn from item pairs. **RankNet** (neural, cross-entropy on
  pairwise order) is the archetype; its loss correlates only loosely with
  list-level NDCG.
- **Listwise** — optimize the whole list. **ListNet** uses a Plackett-Luce
  permutation probability; **LambdaRank/LambdaMART** weight pairwise gradients
  ("lambdas") by the **ΔNDCG** a swap would cause, directly targeting NDCG.
  **LambdaMART** (LambdaRank gradients + gradient-boosted trees / MART) is the
  long-standing workhorse — available in XGBoost (`rank:ndcg`), LightGBM, and
  RankLib. The pairwise-vs-listwise distinction is *how the loss treats the list*,
  not the model family.
- **The retrieval → ranking → re-ranking funnel.** Production recsys is
  multi-stage: (1) **candidate generation/retrieval** — cheap, high-recall, from
  millions to ~hundreds (two-tower + ANN, item-item, popularity); (2) **ranking**
  — expensive, high-precision scorer (DeepFM / LambdaMART / DLRM) on the few
  hundred; (3) **re-ranking** — apply business rules, diversity (MMR), freshness,
  fairness, and exploration on the top slate. Each stage trades recall for
  precision.

### 7. Offline evaluation

Evaluate on held-out interactions (time-based split is more honest than random).

- **Accuracy / relevance @k:** **Precision@k**, **Recall@k**, **Hit Rate@k**,
  **MAP** (mean average precision, order-aware), **MRR** (rank of the *first*
  relevant item — best when one correct answer is expected, e.g. next-item),
  **NDCG@k** (graded relevance, position-discounted — the headline ranking metric).
  Empirically these cluster: {Recall}, {MRR, NDCG, HR}, {Precision, MAP}.
- **Beyond-accuracy:** **Coverage** (fraction of catalog ever recommended),
  **Diversity** (intra-list dissimilarity), **Novelty** (how non-popular the recs
  are), **Serendipity** (relevant *and* surprising). Optimizing accuracy alone
  degrades these and can worsen long-term engagement.
- **The offline/online gap.** Offline metrics measure fit to *logged* behavior,
  which is itself biased by the old recommender. Higher offline NDCG does **not**
  guarantee higher online engagement; too much novelty can hurt online for novice
  users. Treat offline metrics as a *filter*, not the verdict — and beware
  well-documented offline-eval flaws (sampled negatives distort metrics; leakage
  from random splits).

### 8. Online evaluation & off-policy estimation

- **A/B testing** is the gold standard, but recommenders create **feedback loops**
  that contaminate the control/treatment populations over time — keep tests short,
  watch novelty effects, and guard against network/interference effects.
- **Interleaving** mixes results from two rankers into one list per user and
  attributes clicks → far more sensitive than A/B (detects differences with orders
  of magnitude less traffic); used for ranking comparisons (e.g. Airbnb search).
- **Counterfactual / off-policy evaluation (OPE)** estimates how a *new* policy
  would have performed using only logs from the *old* policy:
  - **IPS (inverse propensity scoring):** weight each logged outcome by 1/(logging
    propensity) — unbiased if propensities are known/positive, but **high variance**
    (capping/clipping helps).
  - **Direct Method (DM):** fit a reward model — low variance, biased if misspecified.
  - **Doubly Robust (DR):** combine DM + IPS — unbiased if *either* the propensity
    *or* the reward model is correct, with lower variance than IPS. Standard for
    large action spaces and offline ranker comparison.

### 9. Production concerns

- **Two-stage serving** (candidate gen vs ranking) as above; embeddings live in a
  **vector/ANN store**, features in a **feature store** (Feast) shared across
  train and serve to prevent train/serve skew.
- **Real-time features & embeddings:** session/recency features must be computed
  at request time; embedding tables can be huge (sharded — TorchRec).
- **Exploration vs exploitation:** **contextual bandits** (LinUCB, **Thompson
  sampling**) inject controlled exploration to break feedback loops and gather
  unbiased data; ε-greedy / random injection in re-ranking mitigates filter
  bubbles. Caveat (RecSys 2025): in pure offline eval, greedy models often
  *appear* to beat exploratory bandits — a structural eval bias, not proof
  exploration is useless.
- **Fairness & filter bubbles:** monitor provider-side exposure fairness,
  popularity bias, and diversity; randomization/fairness constraints in re-ranking.

---

## Tools & Frameworks

| Tool | Niche |
| --- | --- |
| **implicit** | Fast WRMF (ALS) + BPR for implicit feedback; Python, multithreaded/GPU. Production baseline. |
| **LightFM** | Hybrid FM blending CF + content/metadata features; great for cold-start; WARP/BPR losses. |
| **Surprise** | Classic explicit-rating CF (SVD, SVD++, KNN baselines); teaching/prototyping, not large-scale. |
| **TorchRec** | PyTorch lib for *large-scale* models — sharded embedding tables across GPUs, two-tower/DLRM. |
| **NVIDIA Merlin** (NVTabular, HugeCTR, Models) | End-to-end GPU recsys: feature prep → train retrieval+ranking → Triton multistage serving. |
| **RecBole** | Research framework, 100+ algorithms, unified benchmarking across CF/sequential/context-aware. |
| **Vespa** | Serving engine that fuses ANN retrieval + tensor ranking + business logic in one query. |
| **Feast** | Feature store for consistent train/serve features. |
| **Amazon Personalize / Google Vertex AI Recommendations** | Managed/turnkey — live in weeks, less control. |
| **XGBoost / LightGBM (`rank:*`), RankLib** | LambdaMART / LTR. |

---

## Methodology (recommended workflow)

1. **Frame the problem.** Implicit or explicit feedback? Top-N retrieval, CTR
   ranking, or next-item? This decides loss (BPR/WRMF vs pointwise vs LambdaMART)
   and metric (Recall@k/NDCG vs MRR).
2. **Baseline first.** Popularity + item-item CF + WRMF/BPR. Many "deep" wins
   vanish against a tuned MF baseline — establish the bar.
3. **Split honestly.** Time-based (leave-last-out per user) over random; avoid
   leakage; be wary of sampled-negative metrics.
4. **Add structure as needed.** Side features → FM/LightFM/DeepFM; sequence →
   SASRec; cold-start/long-tail → content + semantic IDs.
5. **Build the funnel** if scale demands it: two-tower retrieval → DeepFM/
   LambdaMART ranker → diversity/fairness/exploration re-rank.
6. **Evaluate in layers:** offline (NDCG/MAP/Recall + coverage/diversity) →
   off-policy (IPS/DR) → interleaving → A/B. Never ship on offline alone.
7. **Close the loop safely:** log propensities, add exploration, monitor
   popularity bias and exposure fairness.

---

## Practical Patterns

- **Start with item-item CF or WRMF** as the retrieval baseline; it's cheap,
  stable, and hard to beat per unit effort.
- **Log propensities at serving time** — without them, off-policy evaluation and
  de-biasing are impossible later.
- **Use NDCG-aware training (LambdaMART) for the ranker**, not pointwise MSE, when
  the objective is ordering.
- **Precompute item embeddings + ANN index** for two-tower; only the user tower
  runs at request time.
- **Pick metrics by intent:** MRR for "one right answer" (next-item), NDCG for
  graded multi-relevant lists, Recall@k for retrieval-stage recall.
- **Validate offline winners with interleaving before a full A/B** — far cheaper
  traffic-wise.

## Anti-Patterns

- **Trusting offline NDCG as ground truth.** It measures fit to a biased log, not
  online lift. Confirm online.
- **Treating implicit non-interactions as negatives.** They're unlabeled — use
  confidence weighting (WRMF) or sampled pairwise negatives (BPR), not 0-labels.
- **Random train/test split for sequential/temporal data.** Leaks the future →
  inflated metrics.
- **Assuming BERT4Rec > SASRec.** Loss function, not bidirectionality, drove much
  of the reported gap; SASRec is often better and cheaper.
- **Optimizing accuracy only.** Tanks coverage/diversity/serendipity and feeds the
  popularity feedback loop and filter bubbles.
- **Raw IPS with no clipping.** Variance explodes on rare actions; clip or use
  doubly-robust.
- **Deep model with no MF/CF baseline.** You can't claim a win without the bar.

## Troubleshooting

- **Great offline, flat online?** Offline/online gap — biased logs, novelty
  backlash, or a metric that doesn't match the business goal. Add interleaving/OPE.
- **Recommendations collapse to popular items?** Popularity/exposure bias + feedback
  loop; add exploration (bandits), de-bias with propensities, add diversity re-rank.
- **New items never surface?** Item cold-start — add content/FM/LightFM features or
  semantic IDs; reserve exploration slots.
- **OPE estimate wildly unstable?** IPS variance — clip propensities or switch to
  doubly-robust.
- **Two-tower retrieval recall poor?** Check negative-sampling strategy (in-batch vs
  hard negatives), embedding dim, and ANN recall vs exact.

---

## References

**Foundational methods**
- Hu, Koren, Volinsky (2008), *Collaborative Filtering for Implicit Feedback Datasets* (WRMF). https://dl.acm.org/doi/10.1145/1864708.1864726 (related fast-ALS)
- Rendle et al. (2009), *BPR: Bayesian Personalized Ranking from Implicit Feedback*. https://arxiv.org/pdf/1205.2618
- He et al. (2016), *Fast Matrix Factorization for Online Recommendation with Implicit Feedback*. https://dl.acm.org/doi/10.1145/2911451.2911489
- A. Bloch, *An Overview of Collaborative Filtering Algorithms for Implicit Feedback Data*. https://andbloch.github.io/An-Overview-of-Collaborative-Filtering-Algorithms/

**Factorization machines / hybrid**
- Guo et al. (2017), *DeepFM*. https://arxiv.org/pdf/1703.04247
- *FM²: Field-matrixed Factorization Machines*. https://arxiv.org/pdf/2102.12994
- LightFM deep dive (recommenders-team). https://github.com/recommenders-team/recommenders/blob/main/examples/02_model_collaborative_filtering/lightfm_deep_dive.ipynb

**Deep & sequential**
- *The Two-Tower Model for Recommendation Systems* (Shaped). https://www.shaped.ai/blog/the-two-tower-model-for-recommendation-systems-a-deep-dive
- *BERT4Rec* (Sun et al. 2019). https://arxiv.org/pdf/1904.06690
- Petrov & Macdonald (2023), *Turning Dross Into Gold Loss: is BERT4Rec really better than SASRec?* https://arxiv.org/pdf/2309.07602

**LLM-augmented / generative (2025-2026)**
- Spotify Research (2025), *Semantic IDs for Generative Search and Recommendation*. https://research.atspotify.com/2025/9/semantic-ids-for-generative-search-and-recommendation
- *Semantic IDs for Joint Generative Search and Recommendation* (RecSys 2025). https://dl.acm.org/doi/10.1145/3705328.3759300
- *Cold-Start Recommendation with Knowledge-Guided RAG (ColdRAG)*. https://arxiv.org/html/2505.20773v2
- *LLMs as Data Augmenters for Cold-Start Item Recommendation*. https://arxiv.org/pdf/2402.11724

**Learning-to-rank**
- *From RankNet to LambdaMART* (HeThink). https://en.heth.ink/Ranking/
- *LambdaMART Explained* (Shaped). https://www.shaped.ai/blog/lambdamart-explained-the-workhorse-of-learning-to-rank
- XGBoost, *Learning to Rank*. https://xgboost.readthedocs.io/en/latest/tutorials/learning_to_rank.html

**Evaluation (offline / online / off-policy)**
- *Evaluating Recommender Models: Offline vs. Online* (Shaped). https://www.shaped.ai/blog/evaluating-recommender-models-offline-vs-online-evaluation
- *10 metrics to evaluate recommender and ranking systems* (Evidently AI). https://www.evidentlyai.com/ranking-metrics/evaluating-recommender-systems
- *Widespread Flaws in Offline Evaluation of Recommender Systems* (2023). https://arxiv.org/pdf/2307.14951
- *Harnessing Interleaving and Counterfactual Evaluation for Airbnb Search Ranking* (KDD 2025). https://arxiv.org/html/2508.00751v1
- *Counterfactual Learning and Evaluation for Recommender Systems* (foundations). https://par.nsf.gov/biblio/10309941
- *Doubly Robust Estimator for Off-Policy Evaluation with Large Action Spaces*. https://www.researchgate.net/publication/372961616

**Production, bandits, fairness**
- *Introducing TorchRec* (PyTorch). https://pytorch.org/blog/introducing-torchrec/
- *Recommender Systems: Lessons From Building and Deployment* (Neptune). https://neptune.ai/blog/recommender-systems-lessons-from-building-and-deployment
- *Exploration vs. Exploitation in Recommendation Systems* (Shaped). https://www.shaped.ai/blog/explore-vs-exploit
- *Exploitation Over Exploration: Unmasking the Bias in Linear Bandit Recommender Offline Evaluation* (RecSys 2025). https://dl.acm.org/doi/10.1145/3705328.3748166
