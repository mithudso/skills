<!-- hub-reference-banner -->
> **Reference file — part of the `ai-llm-model-layer` hub.** Formerly the standalone `llm-routing-cascades` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

<!--
PROVENANCE: Reference under the `ai-agent-engineering` hub. Authored 2026-05-31 from a completed multi-source research pass (frontier re-expansion of the LLM model-layer family via concept-family-explorer). Routed as a hub reference (not a standalone top-level skill) per hub-and-spoke.
Owns the **multi-model serving-decision** layer — choosing/orchestrating WHICH model(s) answer each request, distinct from serving ONE model well.
Boundaries:
  - Serving ONE model efficiently (vLLM, PagedAttention, continuous batching, autoscaling) → `llm-inference-serving`. This file is WHICH model; that file is HOW to serve one.
  - Speculative DECODING (draft tokens within ONE model, loss-less) → `llm-inference-serving`. Speculative CASCADES (across DIFFERENT models, quality can change) are here.
  - Reasoning route-by-difficulty as a cost bullet → `reasoning-models`. The deep routing discipline is here.
  - Agent orchestration / multi-step tool loops → `agent-ecosystem` / `autonomous-loops`. This is model SELECTION, not the agent loop.
-->

# LLM Model Routing, Cascades & Mixture-of-Agents

Most production LLM traffic is **heterogeneous** — simple classification and Q&A next to hard multi-step reasoning — yet sending every request to one frontier model pays peak price for trivial work. **Routing, cascading, and ensembling** ride the cost / quality / latency **Pareto frontier**: spend big model only where it earns its keep.

Three families, distinguished by **when the decision happens**:

1. **Route** — pick a model *before* generation from query features (1 call).
2. **Cascade** — run a cheap model, *observe* its answer, *escalate* on low confidence (1–N sequential calls).
3. **Ensemble / mixture** — run several models and *fuse* their outputs (N parallel calls).

And the zeroth route: **cache** — serve a remembered answer (0 calls).

## Scope boundary (read first)

- **This reference = the multi-model decision.** Which model, whether to escalate, whether to fuse, whether to cache.
- Serving *one* model for throughput/latency (vLLM/SGLang, PagedAttention, batching) is **`llm-inference-serving`**.
- **Speculative decoding ≠ speculative cascades.** Speculative *decoding* drafts tokens within one model and is provably output-identical to the target (pure latency win) → `llm-inference-serving`. Speculative *cascades* route across two *different* models with a token-level deferral rule and a *controlled quality change* → here (§4). State this distinction explicitly; it is the most common confusion.

---

## 1. Predictive routing (RouteLLM and the router taxonomy)

**Predictive routing** picks a model *before* generation from the query alone. **RouteLLM** (Ong et al., LMSYS/UC-Berkeley/Anyscale, 2024, arXiv:2406.18665) is the canonical OSS framework; it trains four router architectures on human-preference data:

- **Similarity-weighted (SW) ranking** — a weighted-Elo computation where each preference vote is weighted by prompt similarity to the query.
- **Matrix factorization** — learns a latent score for how well a model answers a prompt (best performer on MT-Bench).
- **BERT classifier** — encoder predicting which of strong/weak wins.
- **Causal-LLM classifier** — an LLM fine-tuned to predict relative performance.

The decision is a **strong-vs-weak** binary with a **cost-quality threshold** knob (what fraction of traffic goes to the expensive model). Training data is **Chatbot Arena** preference pairs plus **augmentation** (LLM-judge labels, a small golden-label set — ~1,500 MMLU samples gave large gains). Headline: ~**85% cost reduction on MT-Bench at 95% GPT-4 quality** (matrix-factorization router routing only ~14% to GPT-4). Router quality is measured by **PGR** (Performance Gap Recovered — fraction of the weak→strong gap recovered) and **CPT** (Call-Performance Threshold — % strong calls to hit a target PGR). The taxonomy generalizes to **N-model** routing: encoder/classifier, embedding/kNN-similarity, matrix-factorization, and LLM-judge routers.

## 2. Route-by-difficulty / complexity (and "should I think at all")

A specialization: estimate query **difficulty** and route easy→cheap, hard→strong. In 2025-2026 this fuses with reasoning models — deciding *whether to spend test-time compute*:

- **Route-to-Reason** (arXiv:2505.19435) jointly allocates *both* a model and a reasoning strategy under a budget — higher accuracy than the best single model while cutting tokens >60%.
- **RADAR** (arXiv:2509.25426) frames model-config selection as **multi-objective optimization** for the cost/quality Pareto front.
- **Adaptive think/non-think** routing picks per-instance between reasoning and direct answering *before* generation, because reasoning models **overthink** simple queries (excess deliberation can *lower* easy-input accuracy).

This is the deep treatment of the "route-by-difficulty" idea that `reasoning-models` mentions only as a cost bullet.

## 3. Model cascades + deferral / abstention (FrugalGPT)

A **cascade** is *sequential*: query a cheap model, **score** the answer, **defer** (escalate) to a more expensive model only on low confidence. **FrugalGPT** (Chen, Zaharia, Zou, Stanford, 2023, arXiv:2305.05176) is canonical, with three cost-reduction families:

1. **Prompt adaptation** — cut prompt cost (query concatenation, fewer few-shot exemplars).
2. **LLM approximation** — completion cache + fine-tuning to mimic a stronger model cheaply.
3. **LLM cascade** — the core: a learned ordered list of LLMs.

A trained **generation scoring function** rates each answer; clearing that LLM's **threshold** returns it (abstaining from further calls), else it **defers** to the next. FrugalGPT *learns* which LLMs to include and their thresholds via constrained optimization. Result: matches GPT-4 with **up to 98% cost reduction**, or +4% accuracy at equal cost (HEADLINES/OVERRULING/COQA).

**Calibration is the whole game.** A miscalibrated deferral scorer either over-defers (no savings) or under-defers (quality loss). Self-verification cascades (AutoMix-style) use the model's own consistency as the deferral signal. **Routing vs cascade:** a router picks upfront from features (no feedback, risks misrouting); a cascade pays extra latency to *observe* the cheap answer before escalating.

## 4. Speculative cascades

**Speculative cascades** (Google Research 2025; "Faster Cascades via Speculative Decoding," arXiv:2405.19261, ICLR 2025) hybridize cascades and speculative decoding. Standard cascades are sequential (must finish + assess the small model first); standard speculative decoding verifies a drafter's tokens against the target **but rejects any token that diverges** (output stays identical to the large model). The combination replaces *strict* verification with a **flexible token-level deferral rule**: accept the small model's draft token *even when it differs* from the large model's pick, deferring only when the rule says the large model is needed. The paper characterizes the optimal deferral rule and shows better cost-quality *and* speed-quality trade-offs than either technique alone.

**The load-bearing distinction:** speculative *decoding* = loss-less latency optimization within one model (output identical). Speculative *cascade* = routes across two models, output quality can differ (a controlled trade for speed + cost). → the within-one-model decoding mechanics live in `llm-inference-serving`.

## 5. Mixture-of-Agents (MoA)

**MoA** (Wang et al., Together AI, 2024, arXiv:2406.04692, ICLR 2025 Spotlight) runs models in parallel and **fuses via a layered architecture**: layer-1 **proposers** answer independently; each later layer takes *all* previous outputs as auxiliary context (an Aggregate-and-Synthesize prompt) and refines; a final **aggregator** synthesizes. Reference config: **3 layers × 6 proposers** with a strong aggregator. The key finding — **collaborativeness**: a model produces better output when shown others' responses, *even weaker ones*. Result: **AlpacaEval 2.0 LC win-rate 65.1% with only OSS models vs GPT-4o's 57.5%**. Cost: many calls per query and high time-to-first-token (must wait for all proposers). **Self-MoA** critiques that for many tasks, repeatedly sampling one strong model beats mixing weaker ones — the diversity-vs-quality trade.

## 6. Output ensembling & fusion (LLM-Blender)

**LLM-Blender** (Jiang, Ren, Lin, ACL 2023, arXiv:2306.02561) fuses at the **output** level via two modules: **PairRanker** (pairwise comparison to rank candidates; highest correlation with ChatGPT-based ranking) and **GenFuser** (a generator that fuses the top-K ranked candidates into one improved answer). Trained/evaluated on **MixInstruct** (11 OSS LLMs). The **input-level vs output-level** split: routing/MoA decide before/iteratively; ensembling decides *after* generation (run N, fuse once). PairRanker doubles as a **selection router** (pick best of N); GenFuser **fuses**.

## 7. Semantic / prompt caching as a routing layer

The cheapest route is **no model call**. **GPTCache** (Zilliz) is the canonical **semantic cache**: it matches requests by **embedding similarity** (so paraphrases hit), with a swappable embedding model, vector store (Milvus/Faiss/Redis), threshold-based similarity evaluator, and LRU/FIFO + TTL eviction — up to **~68.8% API-call reduction**. **Exact-match / provider prompt caching** keys on identical text (or repeated prefixes); **semantic caching** keys on meaning, trading a tunable false-hit risk for far higher hit rates. As a routing layer the cache is stage zero: **route-to-cache on hit**, fall through to the router/cascade on miss.

## 8. Cost / quality / latency Pareto modeling

The unifying frame: every request has a cost/quality/latency trade-off and the goal is the **Pareto frontier** — no allocation gives more quality at the same cost. Production orchestrators maintain a frontier across models (updated with empirical quality scores + current prices) and per request **solve a constrained optimization: maximize expected quality subject to cost ≤ budget and latency ≤ SLO**. Two equivalent framings: quality-at-fixed-cost vs cost-at-fixed-quality. Routers and cascades **push the frontier outward** vs any single model (which is one point on the cost axis).

## 9. Router evaluation & benchmarks

**RouterBench** (Hu et al., Martian, 2024, arXiv:2403.12031) is the standard benchmark: **>405k pre-computed inference outcomes** (so routers evaluate offline without re-querying) across commonsense/knowledge/math + a RAG set, with a formal **efficiency-maximization + cost-minimization** framework. Compare routers on **cost-quality curves** and area-under-curve / AIQ (Average Improvement in Quality) aggregates — **evaluate on the curve, not a single operating point**. **RouterArena** (arXiv:2510.00202) extends this. Router training data comes from human-preference logs (Chatbot Arena) + LLM-judge/golden-label augmentation.

## 10. Tooling landscape

A key split: **gateway routing** (load-balancing, fallback, unified API — *not* quality-predictive) vs **quality-predictive routing** (pick the best model per prompt).

| Tool | Type | Notes |
| --- | --- | --- |
| **OpenRouter** | Gateway + Auto | One API to 200+ models/50+ providers; its **Auto Router is powered by NotDiamond** |
| **LiteLLM Router** | Gateway (OSS, self-host) | 100+ providers; routing = load-balancing + fallback (latency-based, least-busy, retries), not quality prediction |
| **NotDiamond** | Quality-predictive | Intelligent model selection by prompt analysis (powers OpenRouter Auto) |
| **Martian** | Quality-predictive | Model-router startup; authors of RouterBench |
| **RouteLLM (OSS)** | Quality-predictive | The academic preference-trained framework (§1), drop-in |
| **vLLM Semantic Router** | OSS, cost-aware | BERT intent/complexity classifier via an Envoy ExtProc gateway (Rust + HF Candle); **v0.1 "Iris"** (Jan 2026) first production release; complements **llm-d** |

## When routing/cascading beats a single frontier model — and when it doesn't

**Wins when:** query difficulty is **heterogeneous**; routing **overhead stays small** (prompt analysis/scoring <~40ms, <5% of a typical 500–3000ms completion); you have representative data to train/calibrate.

**Fails when:**
- **Routing collapse** — as the budget rises, routers default to the most expensive model even when a cheap one suffices (objective-decision mismatch: routers predict scalar scores, decisions are discrete).
- **Tail miscalibration** — the small % the router gets wrong is often the rare, high-stakes query it had least data for.
- **Added latency** — cascades add sequential escalation; MoA/ensembling add many parallel calls + high TTFT.
- **Maintenance** — every new model or price change needs the router re-fit and the Pareto frontier refreshed; non-stationary traffic needs adaptive budget pacing.

## Anti-patterns

- **Confusing speculative cascades with speculative decoding.** Decoding is loss-less within one model; cascades route across models with a quality change. (§4)
- **Routing on an uncalibrated scorer.** Cascades and difficulty routers live or die on calibration — miscalibration causes collapse or tail failures.
- **Evaluating at one operating point.** Use cost-quality curves / AIQ (RouterBench), not single-point accuracy.
- **MoA everywhere.** Many calls + high TTFT; for many tasks a single strong model sampled repeatedly (Self-MoA) wins.
- **Ignoring the cache.** Semantic caching is the cheapest "route" and often the biggest single cost lever.
- **Static router, drifting fleet.** New models/prices silently degrade a frozen router; re-fit on a cadence.

## References (primary sources)

- RouteLLM: Ong et al., "Learning to Route LLMs with Preference Data" — arXiv:2406.18665; LMSYS blog 2024-07-01; github.com/lm-sys/RouteLLM
- FrugalGPT: Chen, Zaharia, Zou — arXiv:2305.05176
- Speculative cascades: Google Research blog 2025; "Faster Cascades via Speculative Decoding" — arXiv:2405.19261 (ICLR 2025)
- Mixture-of-Agents: Wang et al., Together AI — arXiv:2406.04692 (ICLR 2025); github.com/togethercomputer/moa
- LLM-Blender: Jiang, Ren, Lin — arXiv:2306.02561 (ACL 2023); github.com/yuchenlin/LLM-Blender
- RouterBench: Hu et al. — arXiv:2403.12031; github.com/withmartian/routerbench
- Route-to-Reason — arXiv:2505.19435; RADAR — arXiv:2509.25426; RouterArena — arXiv:2510.00202
- GPTCache — github.com/zilliztech/GPTCache; zilliz.com/what-is-gptcache
- vLLM Semantic Router "Iris" v0.1 — blog.vllm.ai/2026/01/05/vllm-sr-iris.html; github.com/vllm-project/semantic-router
- "When Routing Collapses" — arXiv:2602.03478

*Authored via concept-family-explorer frontier re-expansion, 2026-05-31. Treat model IDs / leaderboard figures as point-in-time; anchor exact numbers to the primary source.*
