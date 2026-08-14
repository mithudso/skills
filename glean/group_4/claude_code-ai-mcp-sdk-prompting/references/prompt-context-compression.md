<!-- hub-reference-banner -->
> **Reference file — part of the `ai-mcp-sdk-prompting` hub.** This is the learned/ML **prompt-and-context compression technique** spoke.
> It pairs with `references/llm-context-engineering.md` (which owns *compaction/eviction strategy* and *KV-cache layout*) and is distinct from
> the model-weight compression covered under `ai-llm-model-layer` (`references/llm-compression.md`) and the serving-layer KV-cache /
> prefix-cache work under `ai-llm-model-layer` (`references/llm-inference-serving.md`). Load this file when the question is **"shrink the
> token count of a prompt/context while keeping task quality."**

---

---
name: prompt-context-compression
title: Prompt & Context Compression
description: >-
  ML techniques that shrink the token count of a prompt or retrieved context while
  preserving task quality. TRIGGER: cut prompt/context tokens without losing answer
  quality; LLMLingua / LongLLMLingua / LLMLingua-2 (perplexity scoring, budget
  controller, token-classification compressor); gist tokens / soft-prompt
  compression (prefix-tuning, ICAE, AutoCompressor, 500xCompressor);
  selective-context / self-information token pruning; retrieval-as-compression /
  query-aware summarization; choosing a compression ratio and measuring
  faithfulness; deciding whether to compress, cache, compact, or retrieve. SKIP:
  KV-cache / prefix caching / RadixAttention → ai-llm-model-layer
  (references/llm-inference-serving.md); compaction/eviction STRATEGY and
  what-goes-where design → ai-mcp-sdk-prompting (references/llm-context-engineering.md);
  MODEL-WEIGHT quant/distill/prune → ai-llm-model-layer (references/llm-compression.md);
  RAG chunking/reranking → ai-rag-retrieval (references/rag-architecture.md).
origin: local
category: developer
version: "1.1"
updated: "2026-06-15"
tags:
  - prompt-compression
  - context-compression
  - llmlingua
  - gist-tokens
  - soft-prompts
  - prefix-tuning
  - token-pruning
  - faithfulness
whenToUse:
  - "cut prompt or context tokens without losing answer quality"
  - "choose between LLMLingua-2 and LongLLMLingua for a compression job"
  - "compress a long context for an API-only model (Claude, GPT-4o)"
  - "compress reused context with gist tokens or ICAE when you own the weights"
  - "measure faithfulness or grounding of a compressed prompt, not just accuracy"
  - "decide whether to compress, cache, compact, or retrieve"
related_skills:
  - llm-context-engineering
  - llm-inference-serving
  - llm-compression
  - rag-architecture
metadata:
  changelog:
    - "2026-06-15 sko v1.0->v1.1 — Pass H N/A (folded spoke); 11 Medium fixed (desc 1365->992 chars, +whenToUse/related_skills, fixed self-information wrap, 500x ratio framing, rate-distortion attribution + forward ref, ratio-vs-rate split, dedup beats-full-prompt, +hard-prompt break-even row, em-dash padding trim, arXiv 2304.12102->2310.06201)"
---

# Prompt & Context Compression

The set of ML techniques that **reduce the token count of a prompt or retrieved
context while keeping downstream task quality**. The lever is *information density*,
not *cache reuse* and not *model size*. You compress the input the model reads; you
do not change the model, and you do not rely on a previously-cached prefix.

## When NOT to use this spoke (read the SKIP boundary first)

These four neighbors look adjacent but are different problems. Route correctly:

| If the goal is… | It is NOT this spoke — go to |
|---|---|
| Reuse a previously-encoded prefix so it isn't recomputed (KV-cache, prefix caching, RadixAttention, prompt caching) | `ai-llm-model-layer` (`references/llm-inference-serving.md`); Anthropic `cache_control` mechanics → `references/llm-context-engineering.md` |
| Decide what to keep/evict, summarize history at a threshold, clear stale tool calls — the *strategy* of managing a growing window | `references/llm-context-engineering.md` (compaction, context-editing, Write/Select/Compress/Isolate) |
| Make the *model* smaller/faster (quantization, weight pruning, knowledge distillation into a smaller model) | `ai-llm-model-layer` (`references/llm-compression.md`) |
| Pick chunks, embed, rerank, or transform queries for a corpus | `ai-rag-retrieval` (`references/rag-architecture.md`, `references/advanced-rag-patterns.md`) |

The clean test: **prompt compression takes a specific span of text and returns a
shorter span (or a few learned vectors) that the *same* model reads as a substitute.**
If you are reusing a cache, choosing what to drop as a policy, shrinking the model, or
selecting documents, you are in a neighbor's territory.

## The one taxonomy that organizes the field

Every method is either **hard-prompt** or **soft-prompt** compression (the split used
by the NAACL 2025 survey and the rate-distortion framing of Nagle et al., 2024 — see
"Evaluating compression" below):

| Family | Output | Readable? | Works on a black-box API model? | Needs training? | Representative methods |
|---|---|---|---|---|---|
| **Hard-prompt** (Text→Text, discrete) | shorter natural-language text | yes (often choppy) | **yes** — output is just tokens | usually no (training-free at inference) | Selective Context, LLMLingua / LongLLMLingua / LLMLingua-2 |
| **Soft-prompt** (Text→Vector, continuous) | a few learned embedding/KV vectors | no | **no** — needs access to model internals/weights | yes (train the compressor) | Gist tokens, Prefix/P-tuning, ICAE, AutoCompressor, 500xCompressor |

This decides everything downstream: if you only have an API (Claude, GPT-4o), you can
**only** use hard-prompt compression. Soft-prompt methods require an open-weights model
you can feed embeddings or KV values into.

## Hard-prompt compression — drop low-information tokens

The shared intuition: natural language is redundant, and a **small language model can
score which tokens carry information** so you can delete the rest. The compressed
prompt may read like a telegram to a human but is highly effective for an LLM.

### Selective Context (Li, EMNLP 2023) — the entropy-pruning baseline

- Computes **self-information** (Shannon surprisal, `-log p(token)`) for each token
  using a base causal LM (GPT-2 / OPT / LLaMA).
- Merges tokens into **lexical units** (token / phrase / sentence), then **percentile-
  filters**: keep units above the p-th percentile of self-information; drop the rest.
- Reported ~40%+ token savings on summarization/QA at comparable quality. Limitation:
  it scores each unit roughly independently, ignoring that dropping one token changes
  the surprisal of its neighbors.

### LLMLingua (Jiang et al., EMNLP 2023) — coarse-to-fine, up to 20x

The reference hard-prompt method. Three modules:

1. **Budget controller** — allocates *different* compression ratios to prompt
   components (instruction vs demonstrations vs question), and does coarse,
   demonstration-level pruning first to keep semantic integrity at high ratios.
2. **Iterative Token-level Prompt Compression (ITPC)** — fixes Selective Context's
   conditional-independence flaw by compressing in segments and re-conditioning, so the
   importance of a token accounts for which tokens were already kept.
3. **Distribution alignment** — instruction-tunes the small "compressor" LM so its
   perplexity signal matches the target black-box LLM.

Up to **20x compression with little performance loss** on GSM8K, BBH, ShareGPT,
Arxiv-March23. The small model (GPT2-small / LLaMA-7B) is the scorer; the big model
never sees the original prompt.

### LongLLMLingua (Jiang et al., ACL 2024) — make it question-aware

Adds, for long-context / RAG settings, four things on the LLMLingua backbone:

1. **Question-aware coarse-grained compression** — rank documents by how relevant they
   are *to the question* (conditional likelihood of the question given the doc), not by
   generic perplexity. This raises "key-information density."
2. **Document reordering** — counter the **lost-in-the-middle** position bias by moving
   the most relevant content to the edges.
3. **Dynamic compression ratio** — vary the per-document budget between coarse and fine
   stages.
4. **Subsequence recovery** — a post-step that repairs key entities/spans mangled by
   token-level dropping, restoring faithfulness.

Crucially, on long-context QA the compressed prompt can **beat the full prompt**:
removing noise improves the model, not just the bill.

### LLMLingua-2 (Pan et al., ACL 2024 Findings) — task-agnostic, distilled, bidirectional

A different design that fixes two flaws of perplexity scoring (it is *unidirectional*
and *not aligned to the compression objective*):

- **Data distillation** — prompt GPT-4 to compress text losslessly, building an
  extractive compression dataset.
- **Token classification** — frame compression as per-token *preserve/discard*, trained
  on that dataset. Guarantees the output is a faithful subsequence of the input.
- **Bidirectional encoder** — XLM-RoBERTa-large / mBERT (BERT-level), so each token's
  importance uses the *full* left-and-right context.
- Result: **3x–6x faster** than LLMLingua, smaller compressor, strong **out-of-domain**
  generalization, multilingual.

Mental model: **LLMLingua = perplexity from a causal SLM; LLMLingua-2 = a trained
preserve/discard classifier.** Reach for v2 first for task-agnostic compression and
latency; reach for LongLLMLingua when the task has a clear *question* to condition on.

## Soft-prompt compression — distill text into a few learned vectors

Instead of returning shorter text, train the model to absorb a span into a handful of
continuous vectors it can attend to. Higher ceilings on ratio, but the vectors only
work inside a model whose internals you control.

### The soft-prompt / PEFT lineage (the conceptual root)

- **Prefix-Tuning** (Li & Liang, ACL 2021) — freeze the LM, prepend a sequence of
  continuous **"virtual tokens"** (a prefix) that later tokens attend to; tune only
  ~0.1% of parameters. The first clean statement that a *learned continuous prefix can
  stand in for a long natural-language prompt*.
- **P-tuning** (Liu et al., 2021) and **Prompt Tuning** (Lester et al., 2021) — sibling
  "P*-tuning" methods; prompt tuning shows a single prepended soft prompt (no
  intermediate-layer prefixes) is enough at scale.
- **Wingate et al., 2022** — the explicit bridge: optimize a **soft prompt to reproduce
  the behavior of a longer context** (contrastive conditioning). This is "soft prompt =
  compressed context," but it needs a fresh optimization per context, so it does not
  generalize across contexts.

These are *task/context-specific* compression: great compression, but you re-train per
prompt or per task.

### Gist Tokens (Mu, Li & Goodman, NeurIPS 2023) — generalizable, near-free to train

The key advance over prefix-tuning: learn a model that **compresses arbitrary prompts
on the fly**, not one fixed prefix per task.

- Insert a few **gist tokens** after the prompt; modify the **Transformer attention
  mask** so tokens *after* the gist tokens cannot attend to tokens *before* them. This
  forces the prompt's information to funnel through the gist activations.
- Trained by **ordinary instruction finetuning** — "roughly 10 lines" of mask change,
  **no extra training cost.**
- Up to **26x prompt compression**, ~40% FLOPs reduction, and the gist activations can
  be **cached and reused**, with minimal output-quality loss.

### Context-autoencoder lineage — compress *documents*, not just instructions

- **AutoCompressor** (Chevalier et al., EMNLP 2023) — recursively produce **summary
  vectors** for document segments and feed them as soft prompts to later segments;
  unsupervised objective; extends effective context (tested to 30,720 tokens) and
  summary vectors substitute for plaintext few-shot demonstrations.
- **ICAE — In-Context Autoencoder** (Ge et al., ICLR 2024) — a LoRA-adapted encoder
  compresses a long context into **memory slots**; the frozen LLM decoder conditions on
  them. ~1% extra params, ~4x compression. Frames context compression as an
  autoencoding/working-memory problem.
- **500xCompressor** (Li et al., ACL 2025) — extreme end: compress a context of up to
  ~500 tokens into as few as **one** special token, with reported compression ratios of
  6x–480x ("500x" is the branded headline, not the measured maximum). Two findings worth
  remembering:
  **(a) caching the KV values of the compressed tokens preserves far more information
  than caching their embeddings**; (b) the compressed tokens generalize to unseen text
  and can be used **zero-shot** by the original LLM. Adds ~0.3% params; retains roughly
  70–84% of capability at the extreme ratio.

> Naming caution: "prompt distillation" / "context distillation" usually means *baking a
> prompt's behavior into model parameters or soft vectors* (the soft-prompt lineage
> above). Do **not** confuse it with **knowledge distillation** = training a smaller
> *model* (that is `ai-llm-model-layer` → `references/llm-compression.md`).

## Retrieval-as-compression and abstractive (summary) compression

Two more hard-prompt-adjacent levers, both API-safe:

- **Query-aware summarization** — replace a long context with an LLM-generated summary
  *conditioned on the query*. In head-to-head characterization (Jha et al., ES-FoMo
  2024), this can beat token-pruning by up to ~10 points on multi-doc QA at ~30x.
- **Extractive selection** (rerank-and-keep top spans) — often the strongest simple
  baseline: ">10x compression with minimal accuracy loss" in the same study.
- **Retrieval as compression** — instead of compressing the whole corpus, retrieve only
  the spans you need (the JIT-retrieval idea). The retrieval *mechanics* live in
  `ai-rag-retrieval`; treat it here only as the "don't compress what you can avoid
  loading" alternative on the decision tree below.

## Evaluating compression — the two axes and the faithfulness trap

Compression quality is always a **trade-off curve**, never a single number. Report
both axes:

1. **Compression ratio vs. rate** — the *ratio* is original ÷ compressed (e.g. 6x;
   higher = cheaper, the intuitive number); the *rate* is its inverse, compressed ÷
   original (e.g. 0.17; lower = cheaper, the number you compare against the
   rate-distortion frontier). Report the ratio for intuition and the rate when placing
   a method on that frontier.
2. **Task-quality retention** — downstream metric on the compressed vs. full prompt:
   Exact Match / accuracy for QA & reasoning, BERTScore / ROUGE for summarization.

The **rate-distortion framework** (Nagle et al., 2024) formalizes this as the optimal
distortion achievable at each rate and shows existing methods sit *well below* the
theoretical frontier. It also shows that **variable-rate, query-aware** compression
(e.g. the dynamic-rate LLMLingua-2 variant they evaluate) gets closest, consistent with
the LongLLMLingua "beats-the-full-prompt" effect noted above.

**The faithfulness trap.** Downstream accuracy alone hides whether the model *used the
context* or *guessed from parametric memory*. Two guards:

- **Grounding / faithfulness scoring** — extract claims from the answer and check each
  against the (original) evidence. The EMNLP 2025 information-preservation framework uses
  **FABLES**-style claim verification (an LLM judge rates each decontextualized claim);
  measure **downstream performance, response grounding, and information preservation**
  separately, because a method can score well on one and badly on another.
- **Train/test-overlap leakage** — if eval text overlaps the compressor's or target
  model's pretraining (e.g., Pile/Arxiv), high scores may reflect memorization, not
  compression fidelity. 500xCompressor was explicitly built/evaluated on *strictly
  unseen, cross-domain* QA to dodge this. Always check the eval corpus provenance before
  trusting a headline ratio.

## Decision boundary — compress vs. cache vs. compact vs. retrieve

These four are complementary, not competing; pick by *what changes between calls*:

| Situation | Best lever | Why |
|---|---|---|
| Same long prefix reused across many calls (system prompt, big tool spec, fixed docs) | **Cache** (prefix/KV cache, prompt caching) | Free after first call; zero quality loss; no compression needed. Stable-prefix layout matters → `references/llm-context-engineering.md`. |
| One long context, used once, must stay under budget, **API-only model** | **Hard-prompt compress** (LLMLingua-2 / LongLLMLingua) | Works on a black box; 6–20x; question-aware if there's a query. |
| You **own the model weights** and a context is reused a lot (RAG passages, long instructions) | **Soft-prompt compress** (gist / ICAE / 500xCompressor) | Highest ratios; vectors cache and reuse; needs internal access. |
| Conversation/agent history growing toward the window limit | **Compact** (summarize/evict — *strategy*) | This is window *management*, not span compression → `references/llm-context-engineering.md`. |
| The needed facts are a small slice of a big corpus | **Retrieve** (JIT / agentic search) | Don't compress what you never have to load → `ai-rag-retrieval`. |
| One short, novel context used a single time | **Often: do nothing** | Hard-prompt compression runs a forward pass of the compressor LM, so for a short context used once the compressor's own latency/cost can exceed the token savings. Compression earns out only when the context is large, or the compressed form is reused/cached (gist activations, LLMLingua-2's small fast encoder). |

Common stack in practice: **retrieve → (query-aware) compress the retrieved set →
cache the stable system prefix.** Compress when the context is *novel per call and you
can't cache it*; cache when it's *stable and reused*; compact when it's *history you
own and can summarize*; retrieve when *most of the corpus is irrelevant to this query*.

Heuristics that hold up: prefer **LLMLingua-2** as the task-agnostic default and for
latency; switch to **LongLLMLingua** when there's a clear question to condition on and
position bias to fight; reach for **soft-prompt** methods only when you control weights
*and* a context is reused enough to amortize training; never push the ratio past the
point where **grounding** (not just accuracy) starts to drop.

## Sources

Primary papers (preferred):
- Jiang et al., **LLMLingua: Compressing Prompts for Accelerated Inference of LLMs**,
  EMNLP 2023. https://aclanthology.org/2023.emnlp-main.825/ · code:
  https://github.com/microsoft/LLMLingua
- Jiang et al., **LongLLMLingua: Accelerating and Enhancing LLMs in Long Context
  Scenarios via Prompt Compression**, ACL 2024. https://aclanthology.org/2024.acl-long.91/
- Pan et al., **LLMLingua-2: Data Distillation for Efficient and Faithful Task-Agnostic
  Prompt Compression**, ACL 2024 Findings. https://aclanthology.org/2024.findings-acl.57/ ·
  project: https://llmlingua.com/llmlingua2.html
- Li (Yucheng), **Compressing Context to Enhance Inference Efficiency of LLMs**
  (Selective Context), EMNLP 2023. https://aclanthology.org/2023.emnlp-main.391/ ·
  arXiv:2310.06201
- Mu, Li & Goodman, **Learning to Compress Prompts with Gist Tokens**, NeurIPS 2023.
  https://arxiv.org/abs/2304.08467
- Li & Liang, **Prefix-Tuning: Optimizing Continuous Prompts for Generation**, ACL 2021.
  https://aclanthology.org/2021.acl-long.353/ (also: Lester et al. 2021 Prompt Tuning,
  arXiv:2104.08691; Liu et al. 2021 P-tuning; Wingate et al. 2022 soft-prompt compression)
- Chevalier et al., **Adapting Language Models to Compress Contexts** (AutoCompressor),
  EMNLP 2023. https://arxiv.org/abs/2305.14788
- Ge et al., **In-context Autoencoder for Context Compression** (ICAE), ICLR 2024.
  https://arxiv.org/abs/2307.06945
- Li et al., **500xCompressor: Generalized Prompt Compression for LLMs**, ACL 2025.
  https://aclanthology.org/2025.acl-long.1219/ · arXiv:2408.03094

Surveys / evaluation / theory:
- Li, Liu, Su & Collier, **Prompt Compression for Large Language Models: A Survey**,
  NAACL 2025 (hard vs soft taxonomy). https://aclanthology.org/2025.naacl-long.368/
- **Understanding and Improving Information Preservation in Prompt Compression for LLMs**,
  EMNLP 2025 Findings (FABLES grounding; performance/grounding/preservation axes).
  https://aclanthology.org/2025.findings-emnlp.949.pdf
- Nagle et al., **Fundamental Limits of Prompt Compression: A Rate-Distortion Framework
  for Black-Box Language Models**, 2024. https://arxiv.org/abs/2407.15504
- Jha et al., **Characterizing Prompt Compression Methods for Long Context Inference**,
  ES-FoMo 2024 (extractive vs abstractive vs token-pruning). https://openreview.net/forum?id=vs6CCDuK7l
- **Efficient Prompting Methods for Large Language Models: A Survey**, 2024 (T2V/T2T
  framing). https://arxiv.org/abs/2404.01077
