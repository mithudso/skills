<!-- hub-reference-banner -->
> **Reference file — part of the `ai-llm-model-layer` hub.** Formerly the standalone `llm-compression` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

<!-- Provenance: reference under the `ai-agent-engineering` hub. Created 2026-05-31 via /dr deep-research from primary papers (GPTQ, AWQ, SmoothQuant, LLM.int8, QLoRA/NF4, BitNet b1.58, EfficientQAT, KIVI, KVQuant, OCP MX spec, GKD, MiniLLM, DistiLLM, on-policy-distillation survey 2026, SparseGPT, Wanda, task-arithmetic, TIES, DARE, model-soups) and framework docs (llama.cpp, bitsandbytes/HF, MergeKit, llm-compressor/compressed-tensors). Scope: the MODEL LAYER of LLM compression — making a trained model smaller/cheaper while preserving quality. NOT serving/inference-engine integration (vLLM batching, paged attention, speculative decoding → the sibling `references/llm-inference-serving.md`; reference only), NOT distributed training infra (separate skill), NOT PEFT/LoRA/QLoRA fine-tuning depth (→ fine-tuning/PEFT skill; QLoRA only noted as a *consumer* of NF4 quant here). -->

# LLM Compression: Quantization, Distillation & Merging

Techniques for taking a trained LLM and making it **smaller, faster, or cheaper to
serve** while losing as little quality as possible. Four families:

1. **Quantization** — store weights/activations/KV-cache in fewer bits.
2. **Knowledge distillation** — train a small student to mimic a large teacher.
3. **Pruning & sparsity** — zero out weights the model does not need.
4. **Model merging** — fuse several fine-tuned models into one without retraining.

Plus a cross-cutting fifth concern: **how to evaluate a compressed model** so you do
not ship a regression that benchmarks hide.

## Scope boundary (read first)

- **This reference = the model layer**: the algorithms that produce a compressed set
  of weights (or a smaller model), and how to measure the damage.
- **Serving / inference-engine integration** — running the compressed model
  efficiently (vLLM, TensorRT-LLM, continuous batching, paged attention, speculative
  decoding, tensor/pipeline parallelism) — is the sibling **`references/llm-inference-serving.md`**
  reference in this hub. The compression *format* must match what the engine can load
  (e.g. vLLM reads `compressed-tensors`; llama.cpp reads GGUF), so formats are covered
  here, but serving runtime tuning is not.
- **Distributed training infrastructure** (FSDP, DeepSpeed, megatron, multi-node
  orchestration) → a separate training skill. QAT and on-policy distillation *use*
  training, but the infra is out of scope.
- **PEFT / LoRA / QLoRA fine-tuning** → the **fine-tuning skill**. Note only the
  crossover: **QLoRA fine-tunes adapters on top of a frozen NF4-quantized base**, so
  NF4 (below) is the quantization half of QLoRA — but the adapter/training recipe
  lives in the fine-tuning skill.
- **Offline benchmark methodology** (HELM/MMLU harness mechanics, LLM-as-judge
  benchmarking) → `da-analytical-methods` (references/da-7-machine-learning.md) (`references/da-7-machine-learning.md`).
  The compression-specific *evaluation pitfalls* (perplexity vs KL vs flips) are here.

---

## Part 1 — Quantization

Quantization maps high-precision values (FP16/BF16/FP32) to a low-bit representation.
A single affine scheme stores `q = round(x/scale) + zero_point`; dequantization
reverses it. The key design axes:

- **What you quantize**: **weight-only** (W4A16 — 4-bit weights, 16-bit activations)
  vs **weight+activation** (W8A8 — both 8-bit). Weight-only shrinks memory and is
  bandwidth-bound-friendly for single-stream/decode; W8A8 also speeds up the matmul
  on INT8/FP8 tensor cores and helps at high batch (compute-bound) regimes.
- **Granularity**: per-tensor (one scale) → per-channel / per-row → per-group
  (e.g. group size 128) → per-block (microscaling, 32 elements). Finer granularity =
  better accuracy, more scale-storage overhead.
- **When**: **Post-Training Quantization (PTQ)** vs **Quantization-Aware Training
  (QAT)**.

### 1.1 PTQ vs QAT (the foundational split)

- **PTQ** quantizes an already-trained model, using at most a small **calibration
  set** (often 128–512 sequences) to fit scales — no gradient training. Cheap (minutes
  to a few GPU-hours), the default for inference. GPTQ, AWQ, SmoothQuant, GGUF
  k-quants, and bitsandbytes are all PTQ.
- **QAT** simulates quantization *during* training/fine-tuning with "fake-quant" ops
  in the forward pass, maintaining shadow FP weights and using the **straight-through
  estimator (STE)** to pass gradients through the non-differentiable rounding. More
  accurate at very low bit-widths (≤3-bit) but expensive. Use QAT only when PTQ has
  failed your accuracy bar and you can afford training.

### 1.2 GPTQ — one-shot layer-wise error compensation

GPTQ (Frantar et al., 2023, *arXiv:2210.17323*) is a one-shot weight-only PTQ method.
It quantizes weights **column-by-column** within each layer and, after rounding each
column, **updates the not-yet-quantized columns** to compensate for the error
introduced, using second-order (inverse-Hessian) information approximated from
calibration activations (an Optimal Brain Surgeon descendant). Result: accurate INT4
(and INT3) weights with 3–4× weight-memory reduction. Strong for **weight-only** at
batch-1 decode; needs calibration data; classic config is W4A16 group-128.

### 1.3 AWQ — protect the salient weights

Activation-aware Weight Quantization (Lin et al., 2023, *arXiv:2306.00978*) starts from
the observation that **a small fraction of weight channels (≈0.1–1%) are
disproportionately important**, identified by the *activation* magnitude flowing
through them. Rather than keep those in FP16 (mixed precision is hardware-hostile),
AWQ applies a **per-channel scaling** that mathematically protects the salient
channels before uniform INT4 quantization. No backprop, no regression to specific
calibration data → generalizes better across domains and is fast. AWQ W4A16 is the
common on-device/edge default and is widely supported (vLLM, TensorRT-LLM, HF).

> **GPTQ vs AWQ in practice**: both target W4A16 with similar size/quality. GPTQ's
> error-compensation can edge out accuracy but is more calibration-sensitive; AWQ is
> faster to produce and more robust to calibration choice. Try both and measure
> (see Part 5). Modern toolkits even **combine** them per-submodule.

### 1.4 SmoothQuant — make activations quantizable (W8A8)

The barrier to **activation** quantization is **outlier channels**: a few activation
dimensions have huge magnitudes that wreck per-tensor INT8 scales. SmoothQuant (Xiao
et al., ICML 2023, *arXiv:2211.10438*) applies an offline, mathematically-equivalent
**per-channel smoothing**: divide activations by a per-channel factor `s` and multiply
the corresponding weight columns by `s`, **migrating quantization difficulty from
activations into weights** (which are easy to quantize). This enables accurate **W8A8
INT8** for the whole matmul, giving up to ~1.5× speedup and ~2× memory cut with
negligible loss on large models (OPT-175B matches FP16). SmoothQuant is the canonical
recipe when you want activation quantization for throughput, not just memory.

### 1.5 GGUF / llama.cpp k-quants and i-quants

GGUF is the on-disk model format for **llama.cpp / ggml** (CPU + Apple Metal + GPU
offload), the dominant local-inference stack. Its quantization is block-wise and named
by a convention:

- **Legacy** `Q4_0`, `Q4_1`, `Q8_0` — simple per-block affine.
- **k-quants** `Q2_K … Q6_K` with `_S/_M/_L` variants — a **two-level (super-block)**
  scheme: small blocks get their own scale/min, grouped into a super-block with an
  extra scale, giving piecewise-affine fidelity at non-integer effective bits
  (e.g. **`Q4_K_M` ≈ 4.5 bits-per-weight**, the widely-cited "sweet spot": ~75% size
  cut, minimal quality loss).
- **i-quants** `IQ2_XXS … IQ4_NL` — lower-bit codebook/lookup quantization for
  aggressive compression.
- **Importance matrix (`--imatrix`)** — compute per-weight importance from real-prompt
  activations and quantize important weights more carefully. Strongly recommended
  **below Q5**, marginal at Q4_K_M and above. A 2026 unified evaluation
  (*arXiv:2601.14277*) on Llama-3.1-8B confirms the practical ordering of these
  formats by quality-vs-size.

### 1.6 bitsandbytes — LLM.int8() and NF4

The `bitsandbytes` library (HF-integrated) provides two influential schemes:

- **LLM.int8()** (Dettmers et al., 2022, *arXiv:2208.07339*) — 8-bit weight loading
  with **mixed-precision decomposition**: detect the rare **outlier feature
  dimensions** and compute those in FP16 while the rest run vector-wise INT8. Enables
  zero-degradation 8-bit inference for very large models.
- **NF4 (4-bit NormalFloat)** (from the QLoRA paper, Dettmers et al., 2023,
  *arXiv:2305.14314*) — a 4-bit data type whose 16 levels are the **quantiles of a
  standard normal**, information-theoretically optimal for the (roughly Gaussian)
  weight distribution, beating uniform INT4/FP4. Paired with **double quantization**
  (quantize the quantization constants too) to shave further memory.
  **Crossover note:** QLoRA loads the base model in NF4 and trains LoRA adapters on
  top — NF4 is the *quantization* contribution; the adapter/fine-tuning recipe belongs
  to the **fine-tuning skill**. bitsandbytes is convenient and training-friendly, but
  for inference throughput GPTQ/AWQ/compressed-tensors formats are usually faster.

### 1.7 KV-cache quantization

At long context / high batch, the **KV cache dominates memory**, not the weights.
KV-cache quantization compresses the stored keys/values:

- **KIVI** (Liu et al., ICML 2024, *arXiv:2402.02750*) — tuning-free **asymmetric
  2-bit**, quantizing the **key cache per-channel** and the **value cache per-token**
  (because key outliers are channel-aligned). ~2.6× less peak memory, enabling up to
  4× larger batch and 2.3–3.5× throughput, near-lossless.
- **KVQuant** (Hooper et al., NeurIPS 2024, *arXiv:2401.18079*) — per-channel key
  quant + **pre-RoPE** key quant + non-uniform datatypes; <0.1 perplexity degradation
  at 3-bit, enabling **up to 1M–10M token context** on limited GPUs.

Key insight: **keys and values have different statistics** and want different
quantization axes. This is partly a serving concern (the engine must support the
quantized cache), so coordinate with the `references/llm-inference-serving.md`
sibling — but the algorithm is here.

### 1.8 FP8, INT4, and MX / microscaling formats

- **FP8** (E4M3 for weights/activations, E5M2 for gradients) — 8-bit float with much
  better dynamic range than INT8 for activations; native on NVIDIA Hopper/Blackwell
  and AMD MI300. Often the easiest accuracy-preserving "free" speedup; a **data-free
  PTQ** path exists (no calibration needed for weight FP8).
- **INT4 weight-only** — the aggressive memory play (GPTQ/AWQ land here).
- **MX (Microscaling)** — the **OCP MX v1.0 standard** (2024): a block of (typically)
  **32 elements shares one 8-bit (E8M0) power-of-two scale**, with each element in a
  narrow low-bit type. Formats: **MXFP8, MXFP6, MXFP4, MXINT8**. **MXFP4** uses E2M1
  (1 sign, 2 exp, 1 mantissa) per element + the shared block scale. NVIDIA's
  **NVFP4** is a Blackwell-native compatible E2M1 variant. On Blackwell, MXFP4 matmuls
  are ~2× FP8 and ~4× FP16. Microscaling is the **direction of travel** for hardware
  4-bit: GPT-OSS shipped MXFP4, and `llm-compressor` 0.9 (2026) added MXFP4 support.
  Naive MXFP4 PTQ loses accuracy; recent work (*arXiv:2603.08713*) studies
  error-reduction strategies, and MXFP4 is increasingly a *training* format too.

### 1.9 QAT and extreme low-bit (BitNet, EfficientQAT)

When PTQ can't hold accuracy at very low bits:

- **BitNet b1.58** (Ma et al., 2024, *arXiv:2402.17764*) — QAT from scratch with
  **ternary weights {-1, 0, +1}** = **1.58 bits/param**, replacing matmul with
  add/subtract. Near-lossless vs FP16 at ≤3B params, but requires **pretraining from
  scratch** on the full corpus — a model-design choice, not a post-hoc compressor.
- **EfficientQAT** (Chen et al., 2024, *arXiv:2407.11062*) — makes QAT affordable
  post-hoc: block-wise training of all params then end-to-end scale tuning, producing
  a **2-bit Llama-2-70B in 41 h on one A100** with <3-point degradation.

---

## Part 2 — Knowledge Distillation (KD)

Train a small **student** to reproduce a large **teacher**'s behavior. The supervision
signal and *where the training data comes from* define the method.

### 2.1 Response / logit-based KD (the classic)

The teacher's **soft probability distribution** (logits, optionally temperature-
softened) is the target; the student minimizes a divergence (classically forward KL)
to it. Soft labels carry "dark knowledge" — relative probabilities of wrong answers —
far richer than a one-hot label. For LLMs this is **token-level** logit matching.
Requires **white-box** teacher access (you need the logits).

### 2.2 Feature-based KD

Match **intermediate representations** (hidden states, attention maps) between teacher
and student, not just output logits. Adds dense supervision but needs a projection
when dimensions differ and assumes architectural compatibility. Common in encoder
distillation (e.g. the DistilBERT/TinyBERT lineage); less dominant for modern
decoder-only LLMs than logit/sequence methods.

### 2.3 Sequence-level KD (off-policy / supervised)

The teacher **generates full sequences** and the student is trained (often plain
cross-entropy / SFT) on that synthetic corpus. Simple and scalable — this is what most
"distilled from GPT-4" datasets are. **Black-box compatible** (only needs teacher
*samples*, not logits). Weakness: **exposure bias** — the student only ever sees the
teacher's trajectory, so at inference its own mistakes push it into states it never
trained on, and errors compound autoregressively (O(T²) error growth).

### 2.4 On-policy / online distillation (the 2024–2026 frontier)

On-policy distillation has the **student generate its own trajectories** and the
teacher score them, grounding training in imitation-learning theory and cutting error
compounding to O(T). The *On-Policy Distillation* survey (2026, *arXiv:2604.00626*)
organizes the space along three axes: **feedback signal** (logit-based / outcome-based
/ self-play), **teacher access** (white-box / black-box / self), and **granularity**
(token / sequence / hybrid). Canonical methods:

- **GKD — Generalized KD** (Agarwal et al., 2024, *arXiv:2306.13649*): trains on
  **student-generated** sequences with a configurable divergence (forward KL, reverse
  KL, or JSD) and a mixture policy interpolating dataset and student samples —
  unifying off- and on-policy KD.
- **MiniLLM** (Gu et al., 2024, *arXiv:2306.08543*): uses **reverse KL** (mode-seeking)
  so the small student concentrates on the teacher's major modes instead of averaging
  over everything; framed as policy-gradient RL with teacher log-probs as reward.
- **DistiLLM** (Ko et al., 2024, *arXiv:2402.03898*): stabilizes training with a
  **skewed-KL** objective and an adaptive on-/off-policy mix for efficiency.

**Divergence choice is the central knob:**

- **Forward KL** = mode-**covering** (zero-avoiding) — student tries to cover all
  teacher modes; can smear probability into incoherent "average" text.
- **Reverse KL** = mode-**seeking** (zero-forcing) — student locks onto dominant
  modes; confident but can drop minor valid behaviors.
- **JSD** — symmetric, bounded middle ground. Newest methods (e.g. ToDi) **switch
  per-token** by teacher entropy.

> On-policy distillation is now a leading recipe for compressing reasoning models (the
> student learns from its *own* mistakes against a strong teacher). Distillation can be
> **stacked with quantization** (distill, then quantize the student).

---

## Part 3 — Pruning & Sparsity

Remove weights (set to zero) the model can do without.

### 3.1 Structured vs unstructured

- **Unstructured** — zero individual weights anywhere. Highest quality at a given
  sparsity, but the resulting irregular sparsity is **hard to accelerate** on GPUs
  without specialized kernels — often only a disk-size/theoretical-FLOP win.
- **Structured** — remove whole units (attention heads, neurons/MLP channels, layers).
  Directly shrinks the dense tensors → real speedup on any hardware, but blunter, so
  bigger quality hit per param removed.
- **Semi-structured (N:M)** — the practical middle. **2:4** means **2 of every 4
  contiguous weights are zero**; NVIDIA Ampere+ sparse tensor cores execute this
  pattern at **~2× throughput**. **4:8** is a looser variant.

### 3.2 SparseGPT — one-shot pruning with reconstruction

SparseGPT (Frantar & Alistarh, ICML 2023, *arXiv:2301.00774*) prunes massive LLMs to
**50%+ sparsity in one shot, no retraining**, using OBS-style layer-wise
**reconstruction**: it solves a sparse-regression problem per layer with efficient
inverse-Hessian updates, and like GPTQ updates remaining weights to compensate.
Supports **unstructured and 2:4/4:8 semi-structured**, and **composes with
quantization** (prune + quantize together).

### 3.3 Wanda — pruning by weights × activations, no updates

Wanda (Sun et al., 2023, *arXiv:2306.11695*) prunes by a dead-simple saliency score:
**|weight| × ‖input activation‖** per output, comparing within each output's inputs.
**No weight updates, no second-order info, no retraining** — far cheaper than
SparseGPT and surprisingly competitive, also supporting 2:4. Wanda is the strong,
cheap baseline; SparseGPT buys more accuracy via reconstruction at higher cost. (2024+
methods like Wanda++/OPTIMA add lightweight updates on top of these mask selectors.)

> Sparsity and quantization are complementary; many production stacks ship **2:4 sparse
> + INT4/FP8** together. Always verify the **target inference engine has kernels** for
> the chosen sparsity pattern (a serving concern) — otherwise unstructured pruning is
> just a smaller checkpoint, not a faster one.

---

## Part 4 — Model Merging

Combine **multiple models that share an architecture** (usually fine-tunes of a common
base) into **one** model — **no training, no GPU gradient passes, often just CPU
weight arithmetic** (MergeKit, Arcee). It is *compression* in the sense of collapsing N
deployed checkpoints into one multi-talented model. A **task vector** = `fine-tuned −
base` weights; merging is mostly arithmetic on task vectors.

| Method | Idea | Notes |
| --- | --- | --- |
| **Linear / Model Soups** (Wortsman et al., 2022, *arXiv:2203.05482*) | Average weights of multiple fine-tunes of the same base | "Uniform" or "greedy" soup; improves accuracy/robustness for free, no extra inference cost |
| **SLERP** | Spherical linear interpolation between **two** models' weights | Respects the geometry of weight space better than naive averaging; **only 2 models at a time** |
| **Task Arithmetic** (Ilharco et al., 2023, *arXiv:2212.04089*) | Add/subtract **task vectors**: `θ_base + Σ λ_i τ_i` | Add skills, or **negate** a task vector to *remove* a behavior; simplest multi-task merge |
| **TIES-Merging** (Yadav et al., 2023, *arXiv:2306.01708*) | **TrIm, Elect Sign, Merge**: trim small task-vector entries, resolve **sign conflicts** by majority, then average agreeing params | Fixes interference/redundancy that breaks naive averaging when merging many models |
| **DARE** (Yu et al., 2024, *arXiv:2311.03099*) | **D**rop **A**nd **RE**scale: randomly zero **90–99%** of task-vector deltas, rescale the rest | A pre-processing step; combine with TIES (`dare_ties`) or linear (`dare_linear`); exploits delta redundancy |
| **Model Stock / DELLA / Frankenmerge** | Newer/aggressive variants (Model Stock = fewer models; passthrough/frankenmerge stacks **layers** to grow depth) | Frankenmerges create novel param counts but are unpredictable |

**MergeKit** (Goddard et al., 2024, *arXiv:2403.13257*) is the open-source standard
implementing all of the above via a YAML config, runnable on modest hardware
(low-memory/lazy loading). It powered most top open-LLM-leaderboard merges. NVIDIA and
Arcee publish practitioner guides. **Constraints:** models must share architecture and
(usually) a common base/tokenizer; merging unrelated bases produces garbage. Merging is
near-free to try, so it is high-ROI experimentation — but **measure** (Part 5), since
merge quality is empirical and config-sensitive.

---

## Part 5 — Evaluating a Compressed Model (don't get fooled)

Compression that "passes the benchmark" can still be a real regression. Evaluate on
**three tiers**, and never on accuracy alone.

### 5.1 Perplexity (PPL) — fast, necessary, insufficient

Perplexity (exp of mean per-token NLL on held-out text, e.g. WikiText-2/C4) is the
cheap first screen and the field's default for quantization. **Pitfalls:** it measures
fluency/likelihood, not correctness or reasoning; it is **averaging-biased** and can
look fine while behavior shifts; absolute values are dataset- and tokenizer-dependent
(only compare same-data, same-tokenizer). Use it to *catch gross breakage*, not to
certify quality.

### 5.2 KL divergence / top-token agreement — the better distance metric

Compare the **compressed model's output distribution to the FP16 baseline's**, per
token, via **KL divergence** (or top-1/top-k agreement). This directly measures *how
different the model became*, sidesteps perplexity's averaging bias, and is the metric
llama.cpp's own community adopted over PPL for judging quants. The HF/`kld-guided`
work and Fireworks' evaluation both argue KL-to-baseline is the most reliable
single signal. **imatrix** quantization explicitly optimizes a KL-like importance
objective.

### 5.3 Flips and downstream task accuracy — what users feel

The paper **"Accuracy is Not All You Need"** (Dutta et al., 2024,
*arXiv:2407.09141*) is the key result: across many quantization schemes **accuracy
stays within ~1%**, yet **5–13.6% of individual answers "flip"** (correct↔incorrect)
because roughly equal numbers flip each way and **cancel in the aggregate**. So
matching accuracy ≠ matching behavior. They recommend reporting **flips** (% of
answers that change vs baseline) and **KL divergence** (Spearman ≈0.98 with flips on
MMLU), plus **generative evals (MT-Bench)** which showed **5–10% degradation despite
matching MCQ accuracy**. Practical rule: **(1)** PPL screen → **(2)** KL/flips vs the
FP16 baseline → **(3)** real downstream + generative evals (and, for agents, tool-call
/ trajectory correctness — see `references/llm-observability.md`).

> Bottom line: quantization noise behaves like symmetric noise on token margins —
> harmless to averages, visible to users. Gate releases on **flips + KL + a generative
> eval**, not benchmark accuracy alone.

---

## Tooling map (formats must match your serving engine)

| Tool / format | Covers | Target runtime |
| --- | --- | --- |
| **`llm-compressor`** + **`compressed-tensors`** (vLLM project) | PTQ pipeline: GPTQ, AWQ, SmoothQuant, INT8/FP8/INT4/**NVFP4/MXFP4**, SparseGPT, 2:4, multi-modifier & model-free PTQ (0.8–0.9, 2025–2026) | **vLLM** (serving = separate skill) |
| **llama.cpp / GGUF** (`llama-quantize`, `--imatrix`) | k-quants, i-quants, legacy | CPU / Metal / GPU-offload local inference |
| **AutoGPTQ / GPTQModel**, **AutoAWQ** | GPTQ, AWQ checkpoints | vLLM, TGI, TensorRT-LLM |
| **bitsandbytes** (HF `transformers`) | LLM.int8(), **NF4**/FP4, double-quant | HF inference; QLoRA training (→ fine-tuning skill) |
| **MergeKit** (Arcee) | All merge methods (linear/SLERP/TIES/DARE/task-arith/passthrough) | Produces a normal checkpoint |
| **TorchAO / NVIDIA TensorRT Model Optimizer / AMD Quark** | FP8/INT8/INT4/MX QAT+PTQ | PyTorch / TensorRT / ROCm |

Pick the compressor by **the engine you will serve on** — a GGUF quant won't load in
vLLM and an AWQ checkpoint won't run in llama.cpp. Runtime tuning of that engine is the
**`references/llm-inference-serving.md`** sibling's job.

---

## Decision guide

- **Just need it to fit in VRAM, single-stream**: weight-only **INT4** via **AWQ or
  GPTQ** (W4A16), or **GGUF Q4_K_M** for local/CPU. Start here.
- **Throughput at high batch / compute-bound**: **W8A8** via **SmoothQuant (INT8)** or
  **FP8** on Hopper/Blackwell; explore **MXFP4/NVFP4** on Blackwell.
- **Long context blowing up KV memory**: add **KV-cache quantization** (KIVI/KVQuant),
  independent of weight quant.
- **PTQ accuracy unacceptable at ≤3-bit and you can train**: **EfficientQAT**; for a
  ground-up extreme-low-bit model, **BitNet b1.58**.
- **Want a small model, have a strong teacher**: **distillation** — sequence-level for
  cheap/black-box, **on-policy (GKD/MiniLLM)** for reasoning quality; then quantize the
  student.
- **Have several fine-tunes to consolidate**: **merge** with **MergeKit** (TIES/DARE
  for many models, SLERP for two, task arithmetic to add/subtract skills).
- **Always**: gate on **flips + KL-to-baseline + a generative eval**, not accuracy
  alone.

---

## References

**Quantization**
- GPTQ — *GPTQ: Accurate Post-Training Quantization for Generative Pre-trained Transformers*, Frantar et al., 2023. arXiv:2210.17323
- AWQ — *AWQ: Activation-aware Weight Quantization for LLM Compression and Acceleration*, Lin et al., 2023. arXiv:2306.00978
- SmoothQuant — Xiao et al., ICML 2023. arXiv:2211.10438 · https://hanlab.mit.edu/projects/smoothquant
- LLM.int8() — Dettmers et al., 2022. arXiv:2208.07339
- QLoRA / NF4 — Dettmers et al., 2023. arXiv:2305.14314 · https://huggingface.co/blog/4bit-transformers-bitsandbytes
- bitsandbytes docs — https://huggingface.co/docs/transformers/quantization/bitsandbytes
- llama.cpp quantize README — https://github.com/ggml-org/llama.cpp/blob/master/tools/quantize/README.md · GGUF quant overview gist — https://gist.github.com/Artefact2/b5f810600771265fc1e39442288e8ec9
- *Which Quantization Should I Use? A Unified Evaluation of llama.cpp Quantization on Llama-3.1-8B*, 2026. arXiv:2601.14277
- KIVI — Liu et al., ICML 2024. arXiv:2402.02750
- KVQuant — Hooper et al., NeurIPS 2024. arXiv:2401.18079
- OCP Microscaling (MX) Formats v1.0 Spec, 2024 — https://www.opencompute.org/documents/ocp-microscaling-formats-mx-v1-0-spec-final-pdf · *Microscaling Data Formats for Deep Learning*, arXiv:2310.10537
- MXFP4 strategies — arXiv:2603.08713 · MXFP4 explainer — https://huggingface.co/blog/RakshitAralimatti/learn-ai-with-me
- BitNet b1.58 — Ma et al., 2024. arXiv:2402.17764
- EfficientQAT — Chen et al., 2024. arXiv:2407.11062

**Distillation**
- *A Survey of On-Policy Distillation for LLMs*, 2026. arXiv:2604.00626
- GKD — *On-Policy Distillation of Language Models*, Agarwal et al., 2024. arXiv:2306.13649
- MiniLLM — *Knowledge Distillation of LLMs*, Gu et al., 2024. arXiv:2306.08543
- DistiLLM — Ko et al., 2024. arXiv:2402.03898

**Pruning & Sparsity**
- SparseGPT — Frantar & Alistarh, ICML 2023. arXiv:2301.00774
- Wanda — *A Simple and Effective Pruning Approach for LLMs*, Sun et al., 2023. arXiv:2306.11695

**Model Merging**
- Model Soups — Wortsman et al., 2022. arXiv:2203.05482
- Task Arithmetic — *Editing Models with Task Arithmetic*, Ilharco et al., 2023. arXiv:2212.04089
- TIES-Merging — Yadav et al., 2023. arXiv:2306.01708
- DARE — *Language Models are Super Mario...*, Yu et al., 2024. arXiv:2311.03099
- MergeKit — Goddard et al., 2024. arXiv:2403.13257 · https://huggingface.co/blog/mlabonne/merge-models · https://developer.nvidia.com/blog/an-introduction-to-model-merging-for-llms/

**Evaluation**
- *Accuracy is Not All You Need*, Dutta et al., 2024. arXiv:2407.09141
- KL-guided quantization — https://huggingface.co/blog/rishiraj/kld-guided-quantization · llama.cpp PPL-vs-KL discussion #4110 · Fireworks quantization eval — https://fireworks.ai/blog/fireworks-quantization

**Tooling**
- llm-compressor — https://github.com/vllm-project/llm-compressor · compressed-tensors — https://github.com/vllm-project/compressed-tensors · LLM Compressor 0.9 (MXFP4, attention quant), 2026 — https://developers.redhat.com/articles/2026/01/16/llm-compressor-090-attention-quantization-mxfp4-support-and-more
