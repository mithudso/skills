<!-- hub-reference-banner -->
> **Reference file — part of the `ai-llm-model-layer` hub.** Formerly the standalone `multimodal-llm-architecture` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

<!--
PROVENANCE: This reference is part of the `ai-agent-engineering` hub.
Source: /dr deep-research run, 2026-05-31. Topic — Multimodal & vision-language model (VLM) architecture (2024-2026).
Routed as a hub reference (not a standalone top-level skill) per hub-and-spoke strategy.
Owns the **how a text-only transformer becomes multimodal** layer — the vision/audio/video front-end and how it fuses into the decoder. This EXTENDS `transformer-architecture.md` to other modalities; the decoder block itself (attention/MoE/RoPE/norm) lives there.
Boundaries:
  - The text decoder block — self/MH attention, MoE, RMSNorm, the residual stream, the *text* RoPE/ALiBi menu → `ai-llm-model-layer (references/transformer-architecture.md).md`. Here we cover the MULTIMODAL front-end (vision encoder, connector, image/audio tokens) and the multimodal RoPE *extension* (M-RoPE), not the base block.
  - Visual-instruction-tuning *mechanics* of LoRA/QLoRA — rank/alpha/target-modules, adapter plumbing → `ai-llm-model-layer (references/llm-fine-tuning-peft.md).md`. Here we cover the multimodal training *stages* (projector-align → instruction-tune → preference) and what each stage freezes/trains.
  - Preference-optimization *algorithm* internals — DPO/PPO/the loss math → `ai-llm-model-layer (references/llm-alignment-post-training.md).md`. Here we cover *multimodal* preference / hallucination-reduction (Fact-RLHF, mDPO) as a training stage.
  - Serving a VLM — vLLM/SGLang multimodal, paged KV for image tokens, throughput → `llm-inference-serving.md`.
  - CLIP embeddings as a *retrieval* index / multimodal RAG → `ai-rag-retrieval (references/rag-architecture.md).md` and `ai-datastores.md`. Here CLIP/SigLIP are the *perception front-end of a generative VLM*, not a retrieval embedder.
  - Which VLM to pick for a product (pricing/limits/API) → `ai-llm-model-layer (references/llm-models.md).md`. Here we cover the *architectural* landscape (what each model's front-end + fusion is), not procurement.
  - Image *generation* / diffusion model internals (U-Net, latent diffusion, DiT) → out of scope; this reference covers image *understanding* + the discrete-token *generation* path (Chameleon-style) only.
-->

# Multimodal & Vision-Language Model (VLM) Architecture

This reference is the **answer to one question: how do you turn a text-only decoder LLM into a model that can see (and hear)?** Every other model-layer reference in this hub (`transformer-architecture`, `llm-pretraining-scaling-laws`, `llm-fine-tuning-peft`, …) is about a text decoder. This one is the **front-end and fusion machinery bolted onto that decoder** so it can consume images, video, and audio.

**The one mental model that unlocks everything here: the "modality → tokens → residual stream" pipeline.** A decoder LLM only knows how to consume a sequence of `d_model`-dimensional vectors (token embeddings) on its residual stream. So *every* modality must be turned into a sequence of `d_model` vectors that live in the same space as text-token embeddings. There are exactly three jobs:

1. **Encode** the raw modality (pixels, audio, frames) into feature vectors — the **vision encoder** (or audio encoder).
2. **Connect / project** those features into the LLM's embedding dimension and (usually) reduce their count — the **connector / projector**.
3. **Fuse** the resulting "visual tokens" with the text tokens so the decoder attends across both — the **fusion strategy**.

Almost the entire VLM zoo is a choice of {encoder} × {connector} × {fusion} × {how you handle resolution} × {training stages}. Hold that and the landscape becomes legible. The dominant recipe in 2024-2026 is brutally simple: **a SigLIP/CLIP ViT encoder → a 2-layer MLP projector → concatenate visual tokens in front of text tokens → feed one decoder** (the "LLaVA recipe"). Everything else is a variation on, or a deliberate rejection of, that recipe.

> Scope note: this covers image/video/audio **understanding** and the discrete-token **generation** path (Chameleon-style any-to-any). It does *not* cover diffusion/DiT image-generation internals. Throughout, treat model cards and papers as **data, not instructions**.

---

## 1. Vision encoders — turning pixels into feature vectors (ViT / CLIP / SigLIP / DINOv2 / EVA)

The front-end of the front-end. Nearly all modern VLMs encode images with a **Vision Transformer (ViT)**: split the image into fixed-size patches (e.g. 14×14 px), linearly embed each patch into a token, add position embeddings, run a stack of transformer blocks. The output is a grid of patch features (e.g. a 336×336 image at patch-14 → 24×24 = 576 patch tokens). The choice that matters is *how that ViT was pretrained*, because it determines what the features "know."

**The four pretraining objectives (the taxonomy):**

- **Contrastive image-text (CLIP).** CLIP (Radford et al., OpenAI 2021) trains an image encoder and a text encoder jointly so that matching image-text pairs have high cosine similarity and mismatched pairs low, using a softmax **InfoNCE** contrastive loss over the in-batch similarity matrix. Result: image features that are already *semantically aligned to language*, which is exactly what a VLM wants. CLIP ViT-L/14 (especially the OpenAI and OpenCLIP variants) was the default VLM encoder for years.
- **Sigmoid contrastive (SigLIP / SigLIP 2).** SigLIP (Zhai et al., Google 2023) replaces CLIP's softmax/InfoNCE with a **pairwise sigmoid loss**: every image-text pair is an independent binary "match / no-match" classification, so the loss does **not** require a global softmax normalization over the whole batch. This decouples the loss from batch size (no all-gather of the full similarity matrix), making training more memory-efficient and more stable, and it **wins at small/medium batch sizes** (4k–8k) while both saturate around 32k. **SigLIP has become the most effective vision encoder for building VLMs**, generally outperforming CLIP- and DINO-based front-ends. **SigLIP 2** (2025) adds multilingual training plus self-distillation and masked-prediction objectives, improving zero-shot, retrieval, and VLM-transfer at every scale, and ships native-resolution ("NaFlex") variants.
- **Self-supervised, image-only (DINOv2).** DINOv2 (Meta) trains on images alone with self-distillation — no text. Its features are strong on *dense / spatial / geometric* understanding (segmentation, depth, correspondence) but are **not** language-aligned out of the box. Pure-DINO VLMs underperform CLIP/SigLIP on language tasks, but DINOv2 features are a popular *complement* to a contrastive encoder (see mixture-of-encoders below).
- **Masked / reconstructive (EVA, EVA-CLIP, BEiT lineage).** EVA scales masked-image modeling (reconstruct CLIP features of masked patches) to billion-parameter ViTs; EVA-CLIP combines it with contrastive training. Used as a high-capacity encoder in some large VLMs (e.g. InternVL's InternViT lineage builds a very large vision encoder).

**Single encoder vs mixture-of-encoders.** The default is one encoder. But a 2024 design-space finding (**Eagle**, NVIDIA, arXiv 2408.15998) is that **simply concatenating the visual tokens from several complementary encoders (e.g. CLIP/SigLIP for semantics + DINOv2 for spatial detail + an OCR/text-specialized encoder) is as effective as more complex fusion**, and that *stronger visual perception measurably reduces hallucination* and helps resolution-sensitive tasks (OCR, document understanding). Eagle also adds a "pre-alignment" step to bridge vision-only encoders to language tokens before joint training.

**Frozen vs trained encoder.** Early VLMs froze the encoder (cheap, preserves CLIP's alignment). The 2024-2026 trend is to **train the ViT** (often from scratch, "native-resolution") so it can handle arbitrary aspect ratios and high resolution — **Qwen2.5-VL** and **Pixtral** both train a *new* ViT from scratch rather than reuse a fixed-resolution CLIP. Reference: `transformer-architecture.md` owns the *transformer block* the ViT is built from (attention, norm, FFN); this reference owns *what makes it a vision encoder and how it was pretrained*.

---

## 2. The connector / projector — mapping visual features into the LLM (MLP vs Q-Former vs cross-attention)

The encoder emits, say, 576 image-feature vectors of dim `d_vis`. The decoder wants vectors of dim `d_model` (often different) and ideally *fewer* of them (image tokens are expensive — 576 tokens per image at 24 layers is a lot of KV). The **connector** does dimension-matching and (often) token-count reduction. Three families dominate, in increasing order of complexity:

- **Linear / MLP projection (the LLaVA recipe — the default).** The original LLaVA (Liu et al., NeurIPS 2023) used a single **linear projection matrix**; LLaVA-1.5 upgraded to a **2-layer MLP** (GELU in between). It maps each visual feature to `d_model` and is appended as-is — *no token reduction*, 576 features → 576 visual tokens. This is the simplest, most widely copied connector and works remarkably well. **Qwen2.5-VL** uses an "MLP-based vision-language merger" that *also* merges adjacent patches to cut token count. **InternVL** uses a randomly-initialized MLP after a **pixel-shuffle** that compresses a 448×448 tile from 1024 → 256 visual tokens (¼). The MLP recipe's appeal: trivial, trainable in minutes, and most of the heavy lifting is left to the decoder.
- **Query-based resampler — Q-Former (BLIP-2) and the Perceiver Resampler (Flamingo).** Instead of projecting *every* patch, use a fixed set of **learnable query vectors** that cross-attend to the (frozen) image features and emit a *fixed, small* number of output tokens regardless of input size.
  - **BLIP-2's Q-Former** (Li et al., Salesforce 2023): a small transformer with **32 learnable queries** that attend to the frozen image features (and can also attend to text), compressing the image to **32 tokens** projected into the frozen LLM. Trained in two stages (representation learning, then generative). The point of BLIP-2 was *parameter-efficiency*: bridge a **frozen** vision encoder and a **frozen** LLM with only the lightweight Q-Former trainable.
  - **Flamingo's Perceiver Resampler** (Alayrac et al., DeepMind 2022): converts a variable-length (and multi-frame) visual feature grid into a **fixed number of visual tokens** via latent queries, feeding the cross-attention layers (below).
  - Trade-off: resamplers slash token count (great for many images / video) but the fixed bottleneck can lose fine detail; for high-detail OCR/document tasks the field has largely swung *back* to MLP + more tokens + tiling.
- **Gated cross-attention injected into the decoder (Flamingo, Llama-3.2-Vision).** Rather than concatenating visual tokens into the input sequence, **insert new cross-attention layers between the LLM's existing self-attention layers**, where text queries attend to visual keys/values. Flamingo gates each inserted layer with **`tanh(α)` where `α` is a learnable scalar initialized to 0** — so at init the inserted layer is a no-op and the pretrained LLM is unchanged, then it "opens up" during training (this is what lets you keep the LLM frozen and add vision without destabilizing it). **Llama-3.2-Vision** (Meta 2024) uses exactly this pattern: a ViT-H/14 encoder, an MLP/adapter, and **cross-attention adapter layers interleaved into a frozen Llama-3.1 text model** — the text weights are untouched, only the adapter is trained. See §3 for why this is a different *fusion* category, not just a different connector.

**Picking a connector:** MLP-concat = simplest, best detail, most tokens (default for single-image, OCR-heavy, "just make it work"). Resampler (Q-Former/Perceiver) = fixed small token budget, good for many-image / video / frozen-everything. Cross-attention = keep the LLM weights frozen/untouched and bolt vision on the side (Flamingo/Llama-3.2). `llm-fine-tuning-peft.md` owns the *LoRA mechanics* if you adapt with LoRA; this owns the *connector architecture* itself.

---

## 3. Fusion strategy — how visual tokens meet text in the decoder (unified / early-fusion vs cross-attention vs late fusion)

This is the single most important architectural axis, and it is **distinct from the connector** (you can reach the same fusion with different connectors). Three strategies:

- **Unified / "decoder-only" / channel-concat fusion (early-ish fusion at the input — the dominant design).** Project visual features to `d_model`, then **concatenate them into the token sequence** alongside text embeddings: `[<img tok 1..576> <text tokens…>]`. The *single* decoder then runs full self-attention over the combined sequence — text attends to image and vice versa through the *same* attention layers. This is the LLaVA / Qwen-VL / InternVL / Pixtral recipe. Pros: minimal new parameters, reuses the whole decoder, text↔image interaction is deep (every layer). Cons: image tokens consume context length and KV-cache; high-res images blow up the sequence (hence tiling, §5). This is sometimes called "early fusion" because modalities join *before* the decoder — but note it's still **two separate encoders feeding a shared decoder**, not a single tokenizer (contrast §4).
- **Cross-attention injection (deep / mid fusion, Flamingo / Llama-3.2-Vision).** Visual tokens are **not** placed in the input sequence; instead, *new* cross-attention layers inside the decoder let text tokens attend to visual features. Pros: the original text-token sequence (and its self-attention KV) is unchanged — you can keep the LLM **frozen** and add vision cheaply; image tokens don't eat the text context budget. Cons: more new parameters (the inserted layers), and image-text interaction only happens at the inserted layers. Meta chose this for Llama-3.2-Vision specifically to avoid degrading the strong text-only Llama-3.1.
- **Late fusion (shallow, mostly legacy / dual-encoder).** Encode each modality fully and independently, combine only at the very end (e.g. similarity of pooled embeddings, as in CLIP itself for retrieval, or concatenating final-layer summaries). This is what *CLIP-for-retrieval* does and is the right model for multimodal RAG (→ `ai-rag-retrieval (references/rag-architecture.md).md`), but it is **too shallow for generative reasoning over an image** and is not how generative VLMs fuse.

**Mental model:** unified-concat fuses at the *input* and shares all layers; cross-attention fuses in the *middle* via dedicated layers and can keep the LLM frozen; late fusion fuses at the *output* and is for retrieval/matching, not generation. The 2024-2026 consensus for general-purpose VLMs is **unified-concat with a trained decoder**; cross-attention persists where preserving a frozen base LLM matters (Llama-3.2).

---

## 4. Native any-to-any multimodality & image tokenization (Chameleon, VQ-VAE/VQGAN, Fuyu)

The strategies above keep a *separate* vision encoder. A more radical design **eliminates the separate encoder and the connector** by turning images into **discrete tokens drawn from a vocabulary**, exactly like text BPE tokens — so a single transformer with a single token vocabulary handles both, and can **generate** images as well as read them.

- **Discrete image tokenization (VQ-VAE / VQGAN).** A **vector-quantized autoencoder** is trained separately: an encoder maps an image to a grid of latent vectors, each snapped ("quantized") to the nearest entry in a learned **codebook** (e.g. 8,192 codes); a decoder reconstructs the image from those code indices. The code *indices* are integers — i.e. discrete tokens. **Chameleon** (FAIR/Meta, arXiv 2405.09818) tokenizes a 512×512 image into **1,024 VQ tokens** from an 8,192-code codebook, and crucially puts these image tokens in the **same vocabulary and embedding table as BPE text tokens**.
- **Early-fusion, token-based, mixed-modal (Chameleon).** Because image and text are now the *same kind of token*, Chameleon trains **one decoder-only transformer over interleaved [text, image, text, image, …] sequences from scratch** — no vision encoder, no projector, no modality-specific branches. It can read and **generate** images and text in any order (true "any-to-any" for image+text). The cost: training a unified mixed-modal model from scratch is hard (Chameleon needed architectural stabilizers like query-key normalization and careful norm placement to avoid divergence at scale), and discrete image tokens cap visual fidelity vs continuous features.
- **Patch-as-token, no encoder, continuous (Fuyu, Adept).** Fuyu takes a different shortcut: **no vision encoder at all**, but **not** discretized either — image patches are passed through a *single linear projection* straight into the decoder as continuous "tokens." The decoder itself does all visual processing. This handles arbitrary resolutions trivially (just more patches) and is architecturally minimal, at the cost of asking the LLM to learn vision from scratch.
- **The 2024-2026 frontier: native multimodal models.** Frontier "omni" models (**GPT-4o**, **Gemini**) are described as **natively multimodal** — trained end-to-end across modalities from the start rather than bolting a vision encoder onto a finished text model. Exact architectures are undisclosed; the *publicly documented* approaches to native any-to-any are the discrete-token (Chameleon) and patch-as-token (Fuyu) lines, plus continuous-encoder hybrids. Several open models (e.g. **InternVL**'s "native multimodal pretraining") now interleave multimodal data *during* pretraining rather than only at an instruction-tuning stage, blurring the encoder-bolt-on vs native line.

**When this matters:** if you need a *single* model to both understand **and generate** images (or arbitrary modality interleaving), the discrete-token / native path is the only one that does it in one model. If you only need understanding, the encoder+MLP path (§1-3) is simpler and currently higher-fidelity for OCR/detail.

---

## 5. High-resolution & dynamic tiling (LLaVA-NeXT AnyRes, InternVL dynamic tiles, NaViT native-resolution packing)

A vanilla CLIP/SigLIP ViT runs at a **fixed low resolution** (e.g. 224 or 336). That is fine for "what's in this photo" but fails on documents, dense text, charts, and small UI elements — and naively upscaling explodes patch count. Three families solve "let the model see detail at high/native resolution without a fixed square crop":

- **AnyRes / tiling / "divide into patches" (LLaVA-NeXT, InternVL, most 2024 VLMs).** Split a high-resolution image into a **grid of tiles**, each at the encoder's native resolution, encode each tile separately, **plus** one downsized "thumbnail" of the whole image for global context, then concatenate all the tile tokens. **LLaVA-NeXT** ("AnyRes") picks a grid from a small set like `{2×2, 1×{2,3,4}, {2,3,4}×1}` to match aspect ratio. **InternVL** uses a **dynamic tiling** strategy: 1–12 tiles of 448×448 at train time chosen by aspect ratio/resolution, **zero-shot scalable to ~40 tiles (≈4K) at test time**, each tile pixel-shuffled to 256 tokens. Tiling is the dominant, simplest high-res method; its cost is token count (many tiles → long sequences).
- **Native-resolution packing — NaViT ("Patch n' Pack", Google, arXiv 2307.06304).** Instead of resizing every image to a fixed square, **process each image at its native resolution and aspect ratio** by borrowing **example packing** from NLP: patches from *multiple* differently-sized images are concatenated into one sequence (with **masked attention** so images don't attend across each other) and **factorized position embeddings** encode 2-D patch coordinates. This avoids the information loss of resize/crop and is far more compute-efficient (no padding waste). NaViT-style encoders are now drop-in replacements for fixed-resolution CLIP in large VLMs.
- **Native dynamic-resolution ViT trained from scratch (Qwen2-VL / Qwen2.5-VL, Pixtral).** Rather than tile a fixed encoder, **train the ViT itself to accept variable resolution natively**. **Qwen2.5-VL** redesigns the ViT with **2-D RoPE** and **window attention** so it natively ingests variable resolutions (images resized to multiples of 28, patch stride 14) at low compute overhead; **Pixtral** trains a from-scratch encoder that ingests images at natural resolution and aspect ratio. This is the cleaner long-term answer (no thumbnail hack, no fixed grid) and is why 2025 SOTA encoders are increasingly bespoke rather than reused CLIP.

**Picking a high-res strategy:** AnyRes tiling = easy retrofit onto any fixed encoder, great OCR, costs tokens. NaViT packing / native-resolution ViT = the principled approach, requires training the encoder, best compute efficiency. The general 2024→2026 arc is **fixed-336 CLIP → tile a fixed encoder (AnyRes) → train a native-resolution ViT**.

---

## 6. Audio, video, and speech modalities (Whisper-style encoders, audio tokens, video frame sampling, the omni models)

The same "encode → connect → fuse" pipeline generalizes beyond images.

- **Audio / speech understanding (Whisper-style encoder front-end).** The standard audio front-end converts the waveform to a **mel-spectrogram** (e.g. 128 mel channels) and runs it through a **Whisper-derived encoder** (Whisper-large-v3 is the common choice — Qwen2.5-Omni and InteractiveOmni both adopt it), producing audio feature tokens that are projected into the LLM exactly like visual tokens. This gives an LLM *audio understanding* (transcribe, answer questions about a clip) by reusing a strong pretrained ASR encoder as the perception module.
- **Audio / speech *generation* (audio tokens + a separate decoder).** To *speak*, models emit **discrete audio tokens** that a downstream vocoder/codec decodes to a waveform. **Qwen2.5-Omni** (arXiv 2503.20215) uses a **Thinker-Talker** design: the "Thinker" LLM generates text tokens; a separate "Talker" transformer consumes the Thinker's text + hidden states and produces audio tokens, decoded to a waveform — letting it stream text and speech concurrently without the two interfering. This mirrors the discrete-token *generation* idea from §4, applied to audio.
- **Video (frame sampling + temporal encoding).** Video is the hardest because it is "images × time." The pragmatic approach: **sample frames** (e.g. Qwen2.5-Omni splits video at ~25 fps / 40 ms intervals; many VLMs sample 1–2 fps or a fixed budget of N frames), encode each frame with the vision encoder, optionally compress per-frame tokens (resampler / token merging) to fit the context, and add **temporal position information** so the model knows frame order (see M-RoPE, §7). Qwen2.5-VL adds **absolute time encoding** to localize events to the second over hour-long videos. The central tension is the **token budget**: more frames / higher per-frame resolution = better temporal/spatial detail but quadratic attention cost, so video VLMs lean heavily on resamplers, token merging, and frame-rate tuning.
- **Native omni models.** **GPT-4o** ("omni") and **Gemini** are natively multimodal across text/vision/audio (and GPT-4o is documented as doing native end-to-end audio rather than a pipeline of Whisper→LLM→TTS). Open "omni" models (Qwen2.5-Omni, and 2025-2026 successors) approximate this end-to-end-trained, any-input/any-output behavior with the encoder + Thinker-Talker design above.

---

## 7. Multimodal position encoding (2-D RoPE and M-RoPE / multimodal RoPE)

A text RoPE (→ `ai-llm-model-layer (references/transformer-architecture.md).md` for the base mechanism) encodes a **1-D** position — fine for a token stream, wrong for an image where a token has a **(row, column)** location, and wrong for video which adds **time**. Two pieces:

- **2-D RoPE inside the vision encoder.** Modern from-scratch ViTs (Qwen2.5-VL, Pixtral) apply a **2-D rotary embedding** to patch tokens so attention is aware of each patch's *spatial* (height, width) position — important once you allow variable resolution and aspect ratio (you can't use a fixed learned position table for arbitrary grids).
- **M-RoPE / multimodal RoPE inside the decoder (Qwen2-VL).** When visual tokens are concatenated into the decoder sequence (§3), the decoder's RoPE must give each token a position that respects modality. **Qwen2-VL's M-RoPE** (arXiv 2409.12191) **decomposes the rotary embedding into three components — temporal, height, width** — so a single scheme encodes 1-D text positions, 2-D image positions, and 3-D video (time + 2-D) positions concurrently in one decoder. Text tokens use all three components identically (collapsing to standard 1-D RoPE); image tokens vary height/width at fixed time; video tokens advance the temporal component across frames. This both improves spatial/temporal reasoning and **helps length extrapolation** (the position ids for an image span a 2-D region rather than a long 1-D run, keeping numeric position ids smaller). M-RoPE (and successors in Qwen2.5-VL / Qwen3-VL) is now a common ingredient in unified-fusion VLMs.

**Why it matters:** get position encoding wrong and the model can read an image but not reason about *where* things are ("is the cat left of the dog?", "what's in the top-right cell of the table?"). 2-D/M-RoPE is the cheap fix that makes spatial grounding work.

---

## 8. VLM training stages (projector-alignment pretrain → visual instruction tuning → multimodal preference / DPO)

A VLM is rarely trained in one shot. The canonical **multi-stage recipe** (popularized by LLaVA, refined by everyone) — note this is about *which component is frozen/trained in which stage*, not the LoRA mechanics (→ `ai-llm-model-layer (references/llm-fine-tuning-peft.md).md`) or the preference-loss math (→ `ai-llm-model-layer (references/llm-alignment-post-training.md).md`):

1. **Stage 1 — projector / feature alignment pretraining.** **Freeze both the vision encoder and the LLM; train only the connector/projector** on a large set of **image-caption pairs** (LLaVA used a CC3M subset). Goal: teach the projector to map visual features into the LLM's embedding space so the LLM can "name what it sees." Cheap (only the small projector trains) and fast.
2. **Stage 2 — visual instruction tuning (the "SFT" of VLMs).** **Unfreeze the LLM (and often the projector; sometimes the encoder too)** and train on **multimodal instruction-following data** — (image, instruction, response) triples covering VQA, OCR, reasoning, grounding, conversation. LLaVA's key contribution was *generating* this data by prompting a text-only GPT-4 with image captions/boxes to synthesize instructions. This stage is what turns "a model that captions" into "a model that follows visual instructions and converses." (Modern recipes add a stage that trains the encoder for high-res/native-resolution, and "native multimodal pretraining" folds multimodal data into the base pretraining itself.)
3. **Stage 3 — multimodal preference optimization / alignment (hallucination reduction).** Apply preference optimization (RLHF or **DPO**) with *multimodal* preference data to reduce **hallucination** (describing objects/text not in the image — the #1 VLM failure) and improve helpfulness. Landmark: **LLaVA-RLHF / Fact-RLHF** (arXiv 2309.14525) collects ~10k human preferences over which of two responses is *more hallucinated* and uses **Factually-Augmented RLHF** — feeding the reward model extra ground-truth (captions/boxes) so it can't be fooled by fluent-but-wrong answers; it improves MMHal-Bench and LLaVA-Bench. The 2024-2026 trend moves to **multimodal DPO/mDPO** and self-rewarding variants (e.g. M3PO) as a cheaper alternative to PPO. The *image-conditioning* of the preference signal is the multimodal-specific wrinkle (a naive text-only DPO can ignore the image — mDPO adds image-contrastive terms to force the model to actually condition on the visual input).

**Why staged:** Stage 1 protects the expensive pretrained weights while the random projector finds its footing; Stage 2 is where capability is built; Stage 3 is where trustworthiness (anti-hallucination) is bought. Skipping Stage 1 tends to destabilize; skipping Stage 3 leaves a capable but hallucination-prone model.

---

## 9. The VLM model landscape (architectural map, mid-2026)

Read this as *"what front-end + connector + fusion + resolution strategy does each line use,"* not as a procurement guide (→ `ai-llm-model-layer (references/llm-models.md).md` for selection/pricing).

| Model (family) | Vision front-end | Connector | Fusion | Resolution | Notable |
| --- | --- | --- | --- | --- | --- |
| **CLIP / SigLIP / SigLIP 2** | — (these *are* the encoders) | — | late (retrieval) | fixed (SigLIP2 adds NaFlex native-res) | The encoder others reuse; SigLIP now preferred over CLIP for VLM front-ends |
| **Flamingo** (DeepMind 2022) | frozen CLIP-style | Perceiver Resampler | **gated cross-attention** | fixed | Originated `tanh(α)`-gated cross-attn; few-shot interleaved image-text |
| **BLIP-2** (Salesforce 2023) | **frozen** ViT | **Q-Former** (32 queries → 32 tokens) | unified-concat (into frozen LLM) | fixed | Parameter-efficient frozen-encoder + frozen-LLM bridge |
| **LLaVA / LLaVA-1.5 / LLaVA-NeXT** | CLIP ViT-L/14 (frozen→tuned) | linear → **2-layer MLP** | **unified-concat** | 336 → **AnyRes tiling** | Defined the dominant simple recipe + visual instruction tuning |
| **Qwen2-VL / Qwen2.5-VL / Qwen3-VL** | **native-resolution ViT from scratch** (2-D RoPE, window attn) | MLP merger (patch-merge) | unified-concat + **M-RoPE** | **native dynamic** | Variable res, hour-long video w/ absolute time encoding |
| **InternVL (1.5 / 2.5 / 3)** | **InternViT** (large, incrementally pretrained) | MLP after **pixel-shuffle** (¼ tokens) | unified-concat | **dynamic tiling** (1–12 → ~40 tiles, 4K) | "ViT-MLP-LLM"; v3 does native multimodal pretraining |
| **Llama-3.2-Vision** (Meta 2024) | ViT-H/14 | MLP/adapter | **gated cross-attention** into **frozen Llama-3.1** | tiling | Chose cross-attn to preserve text-only quality |
| **Pixtral 12B** (Mistral 2024) | **new encoder from scratch**, native res/aspect | MLP | unified-concat | **native** | Beats larger models; OCR/document strength |
| **Chameleon** (Meta 2024) | **none** (VQ-VAE discrete image tokens) | **none** (shared vocab) | **early-fusion, single tokenizer** | fixed (512→1024 tokens) | True any-to-any image+text generation, trained from scratch |
| **Fuyu** (Adept) | **none** (linear patch projection) | linear | unified-concat (patches as tokens) | arbitrary | No encoder; decoder does all visual processing |
| **Qwen2.5-Omni** | Whisper-v3 audio enc + ViT | projectors | unified-concat + **Thinker-Talker** | dynamic | Adds audio in + streaming **speech out** |
| **GPT-4o / Gemini** (native omni) | undisclosed, **natively multimodal** (trained end-to-end across modalities) | — | native | native | Reference points for "native" (architectures not public) |

**The arc to remember:** Flamingo (cross-attn, frozen) → BLIP-2 (Q-Former, frozen everything) → **LLaVA (MLP-concat, the simple recipe that won)** → tiling for high-res → **native-resolution from-scratch encoders (Qwen2.5-VL, Pixtral)** for understanding, and **Chameleon/native-omni** for unified any-to-any generation.

---

## 10. Multimodal evaluation & hallucination (MMMU, MMBench, DocVQA, MathVista, POPE/MMHal)

You cannot tune what you can't measure, and VLMs fail in modality-specific ways. The standard 2024-2026 benchmark suite (aggregated on the **OpenVLM Leaderboard**, run with toolkits like **VLMEvalKit**, arXiv 2407.11691):

- **MMMU** (CVPR 2024) — *Massive Multi-discipline Multimodal Understanding*: ~college-level questions across **30 subjects** requiring domain knowledge + figure/diagram/chart reading. The headline "is this VLM smart" exam; far from saturated.
- **MMBench** — bilingual (EN/CN), systematically probes perception, reasoning, and knowledge across fine-grained ability dimensions; uses a circular-eval / answer-shuffling protocol to reduce guessing.
- **DocVQA** — question-answering over **document images** (forms, tables, scanned text). The OCR / high-resolution stress test — this is the benchmark tiling and native-resolution encoders exist to win.
- **MathVista** — **mathematical reasoning in visual contexts** (charts, function plots, geometry, IQ-test figures); a consolidated benchmark including IQTest, FunctionQA, PaperQA. Tests visual *and* quantitative reasoning jointly.
- **Hallucination benchmarks — POPE, H-POPE, MMHal-Bench.** **POPE** (Polling-based Object Probing Evaluation) probes **object hallucination** by asking yes/no "is there a <object>?" — and famously shows VLMs **confirm objects that co-occur frequently but aren't present** (a chair when there's a table). **H-POPE** extends this hierarchically to object *attributes*. **MMHal-Bench** scores hallucination in open-ended responses. Hallucination is the dominant trust failure (§8 Stage 3 exists to fix it), and *better visual perception* (higher resolution, mixture-of-encoders) measurably reduces it (Eagle, §1).

**Evaluation hygiene (same as text):** benchmark **contamination** (test images/questions leaking into training) inflates scores; many VLM benchmarks are **multiple-choice**, so report the eval protocol (answer shuffling, free-form vs MCQ) and watch for prompt-format sensitivity. Offline *general* benchmark-harness mechanics (HELM/MMLU/LLM-as-judge scaffolding) live in `da-analytical-methods` (`references/da-7-machine-learning.md`); the *VLM-specific* benchmarks and the hallucination evals are here.

---

## Practical patterns

- **Default build for image understanding:** SigLIP (or SigLIP 2) ViT → 2-layer MLP projector → unified-concat into your decoder → AnyRes tiling for high-res → 3-stage training (align → instruction-tune → mDPO). This is the boring, reliable recipe; deviate only for a concrete reason.
- **Need to keep a frozen, already-great text LLM:** use **gated cross-attention** (Flamingo/Llama-3.2 pattern) so the text weights stay untouched.
- **Need many images / video / a tight token budget:** use a **resampler** (Q-Former/Perceiver) or aggressive token-merging/pixel-shuffle, and tune frame sampling rate explicitly.
- **Need one model to *generate* images/audio too:** go discrete-token / native (Chameleon for image+text; Thinker-Talker-style for speech out) — accept harder training and lower visual fidelity.
- **OCR / documents / charts are the job:** prioritize **resolution** — native-resolution ViT or many tiles + keep tokens (don't over-compress with a 32-token resampler), consider a **mixture of encoders** (add a text/OCR-specialized encoder).
- **Fighting hallucination:** improve perception first (resolution, encoder), then Stage-3 **mDPO / Fact-RLHF**; evaluate with **POPE + MMHal**, not just MMMU.

## Anti-patterns

- **Over-compressing visual tokens for detail tasks.** A 32-token Q-Former is great for "describe the scene," terrible for reading a dense table. Match token budget to task.
- **Reusing a fixed-336 CLIP for documents.** You will lose small text. Tile (AnyRes) or use a native-resolution encoder.
- **Pure-DINO (no language-aligned encoder) as the *sole* front-end.** Strong spatial features, weak language alignment — underperforms on VLM language tasks unless paired with a contrastive encoder.
- **Skipping projector-alignment (Stage 1) and jumping to full fine-tuning.** A randomly-initialized projector + immediately unfrozen LLM tends to destabilize / waste the pretrained weights.
- **Text-only DPO on a VLM and expecting less hallucination.** A naive preference loss can be satisfied without conditioning on the image; use image-aware preference signals (Fact-RLHF's factual augmentation, mDPO's image-contrastive term).
- **Treating image tokens as free.** Each high-res image can be hundreds-to-thousands of tokens; they consume context length and KV-cache and dominate VLM serving cost (→ `llm-inference-serving.md`).
- **Benchmarking only on MMMU.** A model can score well on multi-discipline MCQ yet hallucinate badly in open-ended use — always include a hallucination probe (POPE/MMHal) and an OCR test (DocVQA).

## Troubleshooting

- **"Model captions fine but can't read text in the image"** → resolution problem. Add AnyRes tiling or a native-resolution encoder; verify image isn't being downscaled to 336.
- **"Model can't reason about spatial relations (left/right, table cells)"** → position-encoding problem. Ensure 2-D RoPE in the encoder and M-RoPE (or equivalent) in the decoder; check visual tokens aren't getting collapsed 1-D positions.
- **"Model invents objects/attributes not present"** → hallucination. Improve perception (resolution/encoder), add Stage-3 image-aware preference optimization, evaluate with POPE/H-POPE/MMHal.
- **"Training diverges when unfreezing everything"** → run Stage-1 projector alignment first; if doing unified-from-scratch (Chameleon-style), you likely need QK-norm / norm-placement stabilizers and careful LR.
- **"Video model OOMs or is too slow"** → token budget. Lower frame rate / frame count, add a resampler or token-merging, or pixel-shuffle per-frame tokens.
- **"Added vision and the text-only quality dropped"** → unified-concat fine-tuning can erode text skills; either use cross-attention into a frozen LLM (Llama-3.2 pattern), or mix text-only data back into the instruction-tuning stage.

## Cross-references (this hub)

- **The text decoder block itself** (self/MH attention, MoE, RMSNorm, residual stream, *text* RoPE/ALiBi/YaRN, FlashAttention as architecture) → `ai-llm-model-layer (references/transformer-architecture.md).md`. **This reference extends that one to other modalities** — the decoder here is the same decoder, with a vision/audio front-end and M-RoPE bolted on.
- **Visual-instruction-tuning LoRA/QLoRA mechanics** (rank, alpha, target-modules, adapter plumbing, SFT data/chat-templating) → `ai-llm-model-layer (references/llm-fine-tuning-peft.md).md`. The training *stages* are here; the *adapter mechanics* are there.
- **Preference-optimization algorithm internals** (DPO/PPO/the DPO-variant family, reward modeling, the loss math) → `ai-llm-model-layer (references/llm-alignment-post-training.md).md`. *Multimodal* preference / Fact-RLHF / mDPO as a stage is here; the algorithm is there.
- **Serving a VLM** (vLLM/SGLang multimodal, paged KV for image tokens, throughput/latency) → `llm-inference-serving.md`.
- **CLIP/SigLIP embeddings as a retrieval index, multimodal RAG** → `ai-rag-retrieval (references/rag-architecture.md).md`, `ai-datastores.md`. Here CLIP/SigLIP are the *generative model's perception front-end*, not a retrieval embedder.
- **Pretraining objectives / scaling laws** for the base text model → `ai-llm-model-layer (references/llm-pretraining-scaling-laws.md).md` (native-multimodal pretraining folds image data into that stage).
- **Which VLM to pick** (capabilities/pricing/limits for a product) → `ai-llm-model-layer (references/llm-models.md).md`.
- **Applying a VLM to *critique a design/UI screenshot*** (prompting for rubric-based critique, structured JSON findings, Set-of-Mark grounding, before/after comparison, VLM-as-visual-judge calibration & the critique failure modes) → `ai-mcp-sdk-prompting` (`references/vision-model-design-critique.md`). The encoder/tiling/POPE-MMHal internals are here; the *applied prompting/evaluation technique* is there.

## References (2024-2026, primary papers + model cards)

- **CLIP** — Radford et al., "Learning Transferable Visual Models From Natural Language Supervision," OpenAI 2021 (the contrastive image-text encoder).
- **SigLIP** — Zhai et al., "Sigmoid Loss for Language Image Pre-Training," Google 2023, arXiv 2303.15343; **SigLIP 2** — Google DeepMind 2025 (HF blog `huggingface.co/blog/siglip2`; multilingual + self-distillation + NaFlex native-res).
- **Vision-encoder & VLM surveys** — "A Survey of State of the Art Large Vision Language Models" (arXiv 2501.02189); "A Survey on Efficient Vision-Language Models" (arXiv 2504.09724); "Vision Language Models: A Survey of 26K Papers" (arXiv 2510.09586); Jina AI "Vision Encoders in Vision-Language Models: A Survey."
- **Flamingo** — Alayrac et al., "Flamingo: a Visual Language Model for Few-Shot Learning," DeepMind 2022 (Perceiver Resampler + tanh-gated cross-attention).
- **BLIP-2** — Li et al., "BLIP-2: Bootstrapping Language-Image Pre-training with Frozen Image Encoders and Large Language Models," Salesforce 2023 (Q-Former, 32 queries).
- **LLaVA / LLaVA-1.5** — Liu et al., "Visual Instruction Tuning," NeurIPS 2023; **LLaVA-NeXT** — llava-vl.github.io blog 2024 (AnyRes high-res tiling).
- **LLaVA-RLHF / Fact-RLHF** — Sun et al., "Aligning Large Multimodal Models with Factually Augmented RLHF," arXiv 2309.14525 (multimodal preference + hallucination).
- **Chameleon** — FAIR/Meta, "Chameleon: Mixed-Modal Early-Fusion Foundation Models," arXiv 2405.09818 (VQ-VAE discrete image tokens, shared vocabulary, any-to-any).
- **Fuyu** — Adept, "Fuyu-8B" model card (decoder-only, linear patch projection, no vision encoder).
- **NaViT** — Dehghani et al., "Patch n' Pack: NaViT, a Vision Transformer for any Aspect Ratio and Resolution," Google, arXiv 2307.06304 (native-resolution packing).
- **Qwen2-VL** — Wang et al., "Qwen2-VL: Enhancing Vision-Language Model's Perception of the World at Any Resolution," arXiv 2409.12191 (M-RoPE, naive dynamic resolution).
- **Qwen2.5-VL** — Qwen Team, "Qwen2.5-VL Technical Report," arXiv 2502.13923 (native-resolution ViT from scratch, window attn, absolute time encoding); **Qwen3-VL** — arXiv 2511.21631.
- **Qwen2.5-Omni** — Qwen Team, "Qwen2.5-Omni Technical Report," arXiv 2503.20215 (Whisper-v3 audio encoder, Thinker-Talker speech generation).
- **InternVL** — Chen et al., "InternVL 1.5" (internvl.github.io blog) and "Expanding Performance Boundaries … (InternVL 2.5)," arXiv 2412.05271 (ViT-MLP-LLM, pixel-shuffle, dynamic tiling).
- **Pixtral 12B** — Mistral AI, "Pixtral 12B," arXiv 2410.07073 (from-scratch native-resolution encoder).
- **Llama-3.2-Vision** — Meta 2024 model card / blog (cross-attention adapter into frozen Llama-3.1; ViT-H/14).
- **Eagle** — NVIDIA, "Eagle: Exploring the Design Space for Multimodal LLMs with a Mixture of Encoders," arXiv 2408.15998 (token-concat ≈ complex fusion; perception reduces hallucination; pre-alignment).
- **Benchmarks** — MMMU (Yue et al., CVPR 2024); MathVista (Lu et al.); DocVQA (Mathew et al.); MMBench (Liu et al.); **POPE** (Li et al., object-hallucination probing) and **H-POPE** (arXiv 2411.04077); **MMHal-Bench**; **VLMEvalKit** (arXiv 2407.11691) + the OpenVLM Leaderboard.
- **GPT-4o / Gemini** — OpenAI / Google model cards (referenced as natively-multimodal "omni" models; architectures undisclosed).
