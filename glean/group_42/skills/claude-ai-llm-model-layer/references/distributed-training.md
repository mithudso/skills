<!-- hub-reference-banner -->
> **Reference file — part of the `ai-llm-model-layer` hub.** Formerly the standalone `distributed-training` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

<!--
PROVENANCE: This reference is part of the `ai-agent-engineering` hub.
Source: /dr deep-research run, 2026-05-31. Topic — Distributed training & training infrastructure for LLMs (2024-2026).
Routed as a hub reference (not a standalone top-level skill) per hub-and-spoke strategy.
Owns the LLM **training-infrastructure** layer: how you split a model + optimizer + activations across many GPUs to TRAIN it. Boundaries:
  - INFERENCE parallelism / serving runtime (vLLM TP/PP, paged KV, continuous batching) → `llm-inference-serving`. Training and serving share the words "tensor/pipeline parallel" but the constraints differ (no backward pass, no optimizer state, no grad sync at serving time).
  - GPU KERNELS / CUDA / Triton / roofline / occupancy → a GPU-kernels sibling (pointer only here). This file treats kernels as a black box and reasons at the collective/parallelism layer.
  - Single-GPU PEFT / QLoRA / LoRA fine-tuning depth → `llm-fine-tuning-peft`. FSDP+LoRA composition is noted here; LoRA mechanics live there.
  - The ARCHITECTURE itself (attention, MoE routing, RMSNorm, RoPE) → `transformer-architecture`. We consume the architecture; we do not derive it. (MoE *expert parallelism* placement is noted here as a parallelism axis; MoE routing math is there.)
  - Pretraining OBJECTIVES, data mixtures, and scaling laws (Chinchilla, compute-optimal N/D) → a pretraining sibling (pointer only). We cover the *systems* of training, not what loss to minimize or how much data.
-->

# Distributed Training & Training Infrastructure

Training a modern LLM does not fit on one GPU. A 70B model in BF16 is 140 GB of weights alone; add Adam optimizer states (2 FP32 moments + an FP32 master copy ≈ 12 bytes/param → 840 GB), gradients (another 140 GB), and activations, and you are far past any single accelerator's 80–192 GB. **Distributed training is the discipline of splitting the four things that consume GPU memory — parameters, gradients, optimizer states, and activations — across tens to tens of thousands of GPUs, while keeping the math identical to single-device training and keeping the expensive accelerators busy.**

**The one mental model that unlocks everything here:** every parallelism strategy is a different answer to "what do we split, and what must we therefore communicate?" There are only a few axes, and real training jobs *compose* them:

- **Data parallel**: split the *batch*; replicate the model; communicate *gradients* (all-reduce). Cheap, but each GPU holds a full model copy.
- **Sharded data parallel (FSDP / ZeRO)**: split the *model states* too; communicate *parameters* (all-gather) on demand plus *gradients* (reduce-scatter). Removes the full-copy cost.
- **Tensor parallel**: split *individual layers / matmuls*; communicate *activations* every layer. High bandwidth need → keep inside one node (NVLink).
- **Pipeline parallel**: split *layers into stages*; communicate *activations* at stage boundaries (point-to-point). Cheap comm, but introduces *bubble* (idle time).
- **Sequence / context parallel**: split the *sequence dimension*; communicate *attention* partials. Unlocks long context.
- **Expert parallel**: split *MoE experts*; communicate *tokens* (all-to-all).

Hold "split-what → communicate-what" in mind and the whole field — FSDP, ZeRO, Megatron, 3D parallelism — becomes legible. Everything else (mixed precision, gradient checkpointing, NCCL tuning, MFU) is about making those splits *fit in memory* and *run fast*.

> Scope guard: this file is the **training** runtime. For *serving* an already-trained model (vLLM/SGLang TP+PP, paged KV cache, continuous batching), see `llm-inference-serving` — the words overlap but serving has no backward pass, no optimizer state, and no gradient sync. For the *kernels* underneath (CUDA/Triton, roofline, occupancy), see the GPU-kernels reference; here kernels are a black box. For *LoRA/QLoRA* single-GPU fine-tuning, see `llm-fine-tuning-peft`. For *attention/MoE/norm* architecture, see `transformer-architecture`. For *what loss / how much data* (scaling laws, Chinchilla), see the pretraining reference. For *distributed-systems consensus* (the unrelated "distributed" homograph — Paxos/Raft/BFT, CAP/FLP, fault-tolerant agreement across nodes), see `distributed-systems-consensus`; this file is GPU model-sharding, not node agreement.

---

## 1. Data parallelism: DDP and the all-reduce baseline

**Data Parallel (DP)** is the floor everything else builds on. Replicate the full model on every GPU, give each a different *micro-batch*, run forward+backward locally, then **all-reduce** the gradients so every replica averages to the same gradient and steps identically. After the step all replicas are bit-identical again.

PyTorch's **`DistributedDataParallel` (DDP)** is the canonical implementation. Its two performance tricks are essential vocabulary:

- **Gradient bucketing**: gradients are grouped into buckets (default ~25 MB); the all-reduce fires per-bucket rather than per-parameter, amortizing collective launch overhead.
- **Backward/comm overlap**: because gradients become ready in roughly reverse layer order during backprop, DDP launches a bucket's all-reduce *as soon as that bucket fills*, overlapping communication with the still-running backward compute (an `async-collective-then-wait` pattern). This is why DDP scales near-linearly when interconnect is healthy. ([PyTorch DDP internals](https://medium.com/@arjunsrinivasan.a/demystifying-pytorch-distributed-data-parallel-ddp-an-inside-look-6d0d42a645ff))

**The fatal limitation:** DDP replicates *everything* — params + grads + optimizer states — on every GPU. It scales *throughput* (more batch in parallel) but not *model size*. A model that doesn't fit on one GPU doesn't fit under DDP either. That single fact is the reason FSDP and ZeRO exist.

**Communication cost:** classic DDP all-reduce moves `2·(N-1)/N · |params|` bytes per step regardless of GPU count (ring all-reduce is bandwidth-optimal — see §9). It is the cheapest data-parallel comm pattern; FSDP trades ~1.5× more comm for the ability to shard. ([PyTorch FSDP paper, VLDB 2023](https://www.vldb.org/pvldb/vol16/p3848-huang.pdf))

---

## 2. ZeRO: shard the optimizer/grad/param redundancy out of data parallelism

DDP's waste is **redundancy**: 8 GPUs each store identical copies of the optimizer states. **ZeRO (Zero Redundancy Optimizer)**, from DeepSpeed, removes that redundancy by *partitioning* the training states across the data-parallel group while keeping data-parallel *semantics*. It comes in three incremental stages — each adds to the previous: ([DeepSpeed ZeRO tutorial](https://www.deepspeed.ai/tutorials/zero/), [ZeRO blog, Microsoft Research](https://www.microsoft.com/en-us/research/blog/zero-deepspeed-new-system-optimizations-enable-training-models-with-over-100-billion-parameters/))

| Stage | What it partitions | Memory reduction | Comm vs DDP |
| --- | --- | --- | --- |
| **ZeRO-1** | Optimizer states (Adam moments + FP32 master) | ~4× | Same as DDP |
| **ZeRO-2** | + Gradients | ~8× | Same as DDP |
| **ZeRO-3** | + Parameters (gathered on demand in fwd/bwd) | Linear in DP degree (64 GPUs → ~64×) | ~1.5× DDP |

ZeRO-1/2 are "free" memory wins (same communication volume as plain DDP, just reduce-scatter+all-gather instead of all-reduce). **ZeRO-3** is the one that lets a model *bigger than one GPU* train under pure data parallelism: parameters live sharded and are **all-gathered just-in-time** for each layer's forward and backward, then immediately freed. The cost is the extra parameter all-gather traffic (~1.5× DDP volume).

**Offload tiers (trading bandwidth for capacity):**

- **ZeRO-Offload**: pushes optimizer states + gradients to **CPU RAM** (built on ZeRO-2). The Adam step runs on CPU. Lets a single GPU train ~10B+ models. ([DeepSpeed ZeRO docs](https://github.com/microsoft/DeepSpeed/blob/master/docs/_tutorials/zero.md))
- **ZeRO-Infinity**: extends ZeRO-3 with an "infinity offload engine" spilling to **CPU *and* NVMe SSD**, enabling trillion-parameter training on limited GPUs by streaming states off disk. Bandwidth-bound, so it is a capacity escape hatch, not a speed play. ([ZeRO-Infinity / zero3 docs](https://deepspeed.readthedocs.io/en/latest/zero3.html))

Mental model: **ZeRO-3 ≈ FSDP** (next section). They are two implementations of the same idea — shard params/grads/opt-states across the DP group and all-gather params on demand. ZeRO came first (DeepSpeed); FSDP is PyTorch-native.

---

## 3. FSDP and FSDP2: PyTorch-native sharded data parallelism

**Fully Sharded Data Parallel (FSDP)** is PyTorch's in-tree implementation of the ZeRO-3 idea: parameters, gradients, and optimizer states are sharded across data-parallel workers; before a layer runs, its full parameters are **all-gathered**; after backward, gradients are **reduce-scattered** back to shards. Peak memory drops roughly linearly with the shard count.

**FSDP1 vs FSDP2 — a real architectural change you must know (FSDP1 is deprecated as of 2025):** ([HF FSDP1 vs FSDP2](https://huggingface.co/docs/accelerate/concept_guides/fsdp1_vs_fsdp2), [PyTorch `fully_shard` docs 2.9](https://docs.pytorch.org/docs/2.9/distributed.fsdp.fully_shard.html), [torchtitan FSDP notes](https://github.com/pytorch/torchtitan/blob/main/docs/fsdp.md))

- **FSDP1 (`FullyShardedDataParallel`)** represents each wrapped module as a single **`FlatParameter`** — all the module's parameters flattened, concatenated, and chunked together as one 1D tensor. This makes per-parameter reasoning (frozen params, mixed dtypes, resharding to other parallelisms) awkward, and sharded state-dicts require extra all-gathers.
- **FSDP2 (`fully_shard`)** represents each parameter as a separate **`DTensor`** sharded on dim-0 (`torch.chunk(dim=0)`). Consequences that matter in practice:
  - **Communication-free sharded state dicts** (no all-gather to save/load) → faster, simpler distributed checkpointing.
  - **Mixed dtypes in one model out of the box** (e.g., some params FP8): impossible under one FlatParameter.
  - **Partial parameter freezing just works** → LoRA/PEFT composes with FSDP2 cleanly (see `llm-fine-tuning-peft`).

**Hybrid Sharded Data Parallel (HSDP / `HYBRID_SHARD`)** is the production sweet spot at multi-node scale: **shard within a node** (across the 8 fast-NVLink GPUs) but **replicate across nodes** (slower inter-node network), then all-reduce between replicas. This bounds the expensive all-gather traffic to intra-node NVLink while still scaling out. ([FairScale FSDP docs](https://fairscale.readthedocs.io/en/stable/api/nn/fsdp.html))

PyTorch's stated direction is a **unified data-parallel API** where DDP, ZeRO-1, ZeRO-2, and FSDP become configuration flavors of one construct. ([Introducing FSDP, PyTorch blog](https://pytorch.org/blog/introducing-pytorch-fully-sharded-data-parallel-api/))

---

## 4. Tensor (model) parallelism: split the matmul

When even one transformer *layer*'s weights or activations are too big — or when you simply want more aggregate compute on each token — **tensor parallelism (TP)**, introduced by **Megatron-LM**, splits individual operations across GPUs. The classic Megatron pattern shards the two big linear layers of each block:

- **Attention / first FFN matmul**: split by *columns* (each GPU computes a slice of the output features), needs **no communication** in forward.
- **Output projection / second FFN matmul**: split by *rows*, results summed with an **all-reduce** (one in forward, one in backward per block).

So TP costs **two all-reduces of the full activation per transformer block per direction** — enormous bandwidth. That is why **TP is kept *inside a single node*** where GPUs share NVLink/NVSwitch; running TP across the data-center network destroys throughput. TP degree rarely exceeds 8 (one node). ([Megatron-LM SC'21 paper](https://people.eecs.berkeley.edu/~matei/papers/2021/sc_megatron_lm.pdf), [Megatron-Core parallelism guide](https://docs.nvidia.com/megatron-core/developer-guide/latest/user-guide/parallelism-guide.html))

**Sequence parallelism (the Megatron companion to TP)** closes a memory gap TP leaves open: the LayerNorm and dropout regions *between* the TP matmuls are *not* split by TP, so their activations are fully replicated. Sequence parallelism shards *those* regions along the **sequence dimension**, converting the TP all-reduce into an all-gather + reduce-scatter pair (same total volume) while cutting activation memory. It is nearly always enabled alongside TP. ([Reducing Activation Recomputation in Large Transformer Models, arXiv 2205.05198](https://arxiv.org/pdf/2205.05198))

> Don't confuse this with §6 *context* parallelism. Megatron "sequence parallelism" shards only the non-attention norm/dropout regions and pairs with TP; *context* parallelism shards the attention computation itself across the full sequence for long context.

---

## 5. Pipeline parallelism: split layers into stages — and fight the bubble

**Pipeline parallelism (PP)** assigns *contiguous groups of layers* to different GPUs ("stages"). A micro-batch flows stage 0 → 1 → 2 → …, with only **point-to-point** sends of activations at stage boundaries — far cheaper communication than TP. The problem PP creates is the **pipeline bubble**: while stage 0 processes the first micro-batch, stages 1..N sit idle; the bubble is the fraction of time GPUs wait for the pipeline to fill and drain.

The history of PP is the history of *shrinking the bubble*: ([Megatron-LM SC'21 paper](https://people.eecs.berkeley.edu/~matei/papers/2021/sc_megatron_lm.pdf), [Megatron-Core pipeline schedules](https://github.com/NVIDIA/Megatron-LM/blob/main/megatron/core/pipeline_parallel/schedules.py), [Megatron-Core pipeline API](https://docs.nvidia.com/megatron-core/developer-guide/latest/api-guide/pipeline_parallel.html))

- **GPipe**: split the global batch into many **micro-batches** and push them through; bubble fraction ≈ `(P-1)/m` for `P` stages and `m` micro-batches. More micro-batches → smaller bubble, but all activations are stashed until backward → high activation memory.
- **PipeDream / 1F1B** (one-forward-one-backward): interleave a backward right after each forward once the pipeline is full. Same bubble fraction as GPipe but **caps in-flight activations to the pipeline depth**, slashing activation memory. This is the standard non-interleaved schedule.
- **Interleaved 1F1B (virtual pipeline stages)**: give each GPU *multiple non-contiguous* layer chunks. Bubble shrinks by the number of chunks `v` (≈ `(P-1)/(m·v)`) at the cost of *more* point-to-point communication. Megatron's default for large jobs.
- **Newer schedules**: research like **Seq1F1B** decomposes batch-level units into finer *sequence-level* units to shrink both bubble and activation memory further; "zero-bubble" and **DualPipe** (used in DeepSeek-V3) overlap forward/backward across directions to nearly eliminate the bubble. ([Seq1F1B, arXiv 2406.03488](https://arxiv.org/pdf/2406.03488))

**Rule of thumb:** PP is the *outer*, cross-node axis (cheap comm tolerates slow links); the more micro-batches you can fit, the smaller the bubble.

---

## 6. Sequence / context parallelism: training at 100K–1M+ tokens

Attention is `O(sequence²)`. As context grows from 4K to 128K to 1M tokens, the attention activations and the KV tensors blow past a single GPU even if the *weights* fit. **Context parallelism (CP)** shards the **sequence dimension across GPUs** so each holds only a slice of the tokens, and the attention computation is done collaboratively. The two landmark 2023–2024 approaches still define the field: ([A Unified Sequence Parallelism Approach, arXiv 2405.07719](https://arxiv.org/html/2405.07719v3), [Megatron-Bridge parallelisms guide](https://docs.nvidia.com/nemo/megatron-bridge/latest/parallelisms.html))

- **Ring Attention**: arrange GPUs in a ring; each holds a query block and *rotates* K/V blocks around the ring, overlapping the P2P transfer with the local attention compute so communication is hidden. Memory per GPU becomes `O(seq/P)`.
- **DeepSpeed-Ulysses**: uses **all-to-all** to redistribute so each GPU computes full attention for a *subset of heads* over the full sequence; communication volume stays constant as you scale sequence and devices proportionally.
- **Unified / hybrid SP (USP, "2D" Ulysses+Ring)**: combine both axes; production frameworks (Megatron-Core CP, NeMo) expose CP as a first-class degree. Reported results: CP cuts attention-layer intermediate memory by up to ~87.5% for a 32B model and enables ~5M-token training of Llama3-8B on a single 8×H100 node. ([context-parallelism survey results](https://arxiv.org/html/2405.07719v3))

CP composes with the other axes (it is one of the "D"s in 4D parallelism). Use it when *sequence length*, not parameter count, is the thing that won't fit.

---

## 7. 3D / ND parallelism: composing the axes (the real recipe)

No frontier model uses a single axis. **3D parallelism** composes **data × tensor × pipeline** (and at frontier scale **4D/5D** adds **context** and **expert** parallelism). The axes are orthogonal and chosen by their communication cost vs. the interconnect tier:

| Axis | Comm pattern | Volume | Placement |
| --- | --- | --- | --- |
| **Tensor (TP)** | all-reduce / all-gather of activations, **every block** | Highest | Inside a node (NVLink) |
| **Context (CP)** | P2P / all-to-all of attention partials | High | Inside / across a few nodes |
| **Pipeline (PP)** | P2P of stage-boundary activations | Low | Across nodes |
| **Data / FSDP (DP)** | all-reduce / all-gather of grads+params, **once per step** | Low–medium | Outermost, across the fleet |
| **Expert (EP)** | all-to-all of tokens (MoE only) | High (bursty) | Across nodes for many experts |

**The placement heuristic (memorize this):** **TP innermost (NVLink), then CP, then PP across nodes, then DP/FSDP outermost.** Map the highest-bandwidth axis to the fastest links. The global batch size factorizes as `micro_batch × grad_accum × DP_degree`, while `TP × PP × CP × EP × DP = total GPUs`.

**Worked examples from real 2024–2026 training:**
- **Llama 3.1 8B**: FSDP2 only (fits with sharding; maximize per-device throughput). ([torchtitan, ICLR 2025](https://proceedings.iclr.cc/paper_files/paper/2025/file/e6231c5f46598cfd09ff1970524e0436-Paper-Conference.pdf))
- **Llama 3.1 70B / 405B**: progressively stack **TP → then PP** on top of FSDP2; 405B training uses **TP × PP × CP × DP** 4D parallelism. ([torchtitan ICLR paper](https://proceedings.iclr.cc/paper_files/paper/2025/file/e6231c5f46598cfd09ff1970524e0436-Paper-Conference.pdf), [NeMo/Megatron 405B stack](https://developer.nvidia.com/blog/boosting-llama-3-1-405b-throughput-by-another-1-5x-on-nvidia-h200-tensor-core-gpus-and-nvlink-switch/))
- **DeepSeek-V3 (671B MoE, 37B active)**: **expert parallelism** (e.g., EP=64) × **pipeline parallelism** (PP=8) × DP, with DualPipe to hide the bubble; MoE makes **expert parallelism** a first-class axis. ([NeMo DeepSeek-V3 recipe](https://github.com/NVIDIA-NeMo/NeMo/blob/main/nemo/collections/llm/recipes/deepseek_v3.py), [Megatron-Bridge DeepSeek V3](https://docs.nvidia.com/nemo-framework/user-guide/latest/llms/deepseek_v3.html))

> Expert parallelism *placement* belongs here as a parallelism axis; the MoE *routing/gating/load-balancing* math belongs to `transformer-architecture`.

---

## 8. Mixed precision: BF16, FP8, and loss scaling

Training in full FP32 wastes memory and bandwidth. **Mixed precision** keeps a high-precision *master copy* of weights and an FP32 accumulation for the optimizer, but runs the heavy matmuls in a low-precision format. The format choice is the crux: ([Mixed-precision for LLMs](https://ehga.org/mixed-precision-training-for-large-language-models-fp16-bf16-and-beyond), [FP8-LM, arXiv 2310.18313](https://arxiv.org/pdf/2310.18313))

- **FP16**: 5 exponent bits, 10 mantissa. Narrow dynamic range → gradients underflow to zero. **Requires loss scaling**: multiply the loss by a factor `S` before backward so small gradients land in FP16's representable range, then unscale before the optimizer step. **Dynamic loss scaling** auto-adjusts `S` (back off on overflow/NaN, grow when stable).
- **BF16**: 8 exponent bits (same range as FP32), 7 mantissa. The dynamic range means **no loss scaling needed** and far better training stability. **BF16 is the default for LLM pretraining** on modern hardware for exactly this reason; FP16's instability is why BF16 won.
- **FP8 (E4M3 / E5M2)**: Hopper/Blackwell + NVIDIA **Transformer Engine** run matmuls in FP8 with per-tensor scaling. Reported: training GPT-175B in FP8 cut memory ~42% and ran ~64% faster than BF16. The challenge is *stability* — FP8's tiny range needs **per-tensor scaling factors** (Transformer Engine's `DelayedScaling` tracks an amax history) or analytic init schemes (e.g., μnit Scaling, 2025) that remove dynamic scaling entirely. ([FP8 vs BF16 trade-offs, arXiv 2411.08719](https://arxiv.org/html/2411.08719v1), [COAT FP8 training, arXiv 2410.19313](https://arxiv.org/pdf/2410.19313))

**Practical posture (mid-2026):** BF16 is the safe default; FP8 (selective — keep sensitive ops like the final logits and norms in higher precision) is the throughput frontier on H100/H200/B200, used when you can validate stability. Always keep master weights + optimizer moments in FP32.

---

## 9. Gradient checkpointing (activation recomputation) and gradient accumulation

Two complementary memory levers that have nothing to do with sharding — they reduce *per-GPU* memory directly.

**Gradient checkpointing / activation recomputation.** Backprop normally stores every layer's forward activations to reuse in the backward pass; for long sequences and deep models this dominates memory. Checkpointing instead **discards intermediate activations and recomputes them during backward**, trading compute for memory (classic uniform checkpointing recomputes the whole layer for ~`O(√L)` memory). ([Activation checkpointing explained](https://mbrenndoerfer.com/writing/activation-checkpointing-gradient-memory-selective-recomputation), [NeMo activation recomputation docs](https://docs.nvidia.com/nemo-framework/user-guide/24.09/nemotoolkit/features/optimizations/activation_recomputation.html))

- **Selective activation recomputation** (Megatron) is the key refinement: checkpoint and recompute *only* the operations that are **memory-heavy but cheap to recompute** (e.g., the attention softmax/dropout region), keeping the expensive matmul activations in memory. This recovers most of the memory at a fraction of the recompute penalty vs. full-layer checkpointing. ([Reducing Activation Recomputation, arXiv 2205.05198](https://arxiv.org/pdf/2205.05198))
- **Finer granularity**: checkpointing attention and FFN sublayers separately can cut peak activation memory another ~30–40% for proportionally more recompute.

**Gradient accumulation.** Run several micro-batches forward+backward, *accumulate* their gradients locally, and step the optimizer once. This decouples the **global (effective) batch size** from what fits in memory: `effective_batch = micro_batch × grad_accum_steps × DP_degree`. It is how you hit a target large batch on limited GPUs, and it pairs naturally with checkpointing (checkpointing caps per-micro-batch activation memory; accumulation simulates the big batch without raising peak memory). ([Gradient accumulation & checkpointing](https://mljourney.com/gradient-accumulation-and-gradient-checkpointing-explained/))

> Under DDP/FSDP, wrap all-but-the-last accumulation micro-step in `no_sync()` so gradients aren't all-reduced every micro-batch — only on the final accumulation step. Forgetting this is a common silent throughput killer.

---

## 10. Collective communication: NCCL and compute–comm overlap

Every parallelism axis above bottoms out in a handful of **collective operations**, and on NVIDIA hardware these are implemented by **NCCL (NVIDIA Collective Communications Library)**. Knowing the primitives and their cost models is what lets you reason about scaling: ([NCCL collectives docs](https://docs.nvidia.com/deeplearning/nccl/user-guide/docs/usage/collectives.html), [NCCL deep dive, NVIDIA blog](https://developer.nvidia.com/blog/nccl-deep-dive-cross-data-center-communication-and-network-topology-awareness/))

- **AllReduce**: sum a tensor across all ranks, result on all ranks. The DDP gradient-sync primitive. NCCL implements it as **reduce-scatter + all-gather**.
- **ReduceScatter**: sum, but each rank keeps only its slice (FSDP gradient sync).
- **AllGather**: each rank contributes its slice; all ranks end with the full tensor (FSDP parameter gather; the TP+SP forward).
- **All-to-All**: each rank sends a distinct chunk to every other rank (MoE expert routing, Ulysses CP).
- **Broadcast / Reduce / P2P Send-Recv**: P2P send/recv is the pipeline-parallel stage-boundary primitive.

**Algorithms and their cost:** ([NCCL algorithms overview](https://medium.com/@nitin966/unpacking-nccl-a-deep-dive-into-multi-gpu-communication-2b667e77d96d))
- **Ring**: bandwidth-optimal, latency grows with rank count; data is chunked and pipelined around a logical ring. NCCL's default for large messages. AllReduce traverses the ring twice (reduce, then broadcast).
- **Tree (double-binary)**: `O(log N)` latency; better for **small/medium** messages at large scale (NCCL 2.4+).
- NCCL auto-selects ring vs tree (vs CollNet) by message size and topology; channels map a collective across the GPU's SMs and the available NVLink/IB links.

**The single most important systems lever: compute–communication overlap.** Idle GPUs waiting on a collective is the main scaling loss. Techniques:
- DDP/FSDP **prefetch** the next layer's all-gather while the current layer computes, and overlap gradient reduce-scatter with backward.
- **Async tensor parallelism** (torchtitan) overlaps the TP all-gather/reduce-scatter with the matmul.
- Pin TP to NVLink so its heavy traffic never touches the slower inter-node fabric.

When overlap breaks down (e.g., a too-small message, a straggler GPU, an oversubscribed NIC), MFU (§12) collapses — overlap is what keeps the FLOP units fed.

---

## 11. Training stability at scale: loss spikes, z-loss, warmup, init

Large runs cost millions of dollars and run for weeks; an unrecovered **loss spike** (a sudden divergence where loss shoots up and may NaN) can waste days. Stability engineering is a first-class concern: ([Methods of Improving LLM Training Stability, arXiv 2410.16682](https://arxiv.org/pdf/2410.16682), [ZClip adaptive spike mitigation, arXiv 2504.02507](https://arxiv.org/pdf/2504.02507))

- **Learning-rate warmup**: start the LR near zero and ramp it linearly over the first few thousand steps. A large LR before the model has calibrated produces a burst of huge gradients; warmup lets training start in a stable regime. It is standard in essentially every LLM recipe and contributes as much stability as normalization or good init. ([Why warmup the LR, NeurIPS 2024](https://proceedings.neurips.cc/paper_files/paper/2024/file/ca98452d4e9ecbc18c40da2aa0da8b98-Paper-Conference.pdf))
- **z-loss**: a regularizer that pushes the softmax normalizer (`log Z`) toward 0, preventing the output **logits from drifting to large magnitudes** late in training (a common spike cause, used in PaLM/Chinchilla-class runs). **QK-LayerNorm** (normalizing queries/keys before attention) addresses the related *attention*-logit growth. Newer work notes z-loss and logit soft-capping treat symptoms; output-embedding centering targets the cause.
- **Initialization & residual scaling**: small/scaled init (e.g., scaling residual-branch weights by `1/√(2·n_layers)`) keeps activation variance bounded in deep stacks, a precondition for stable deep-transformer training.
- **Gradient clipping & spike detection**: global-norm gradient clipping is standard; adaptive methods like **ZClip** track the gradient-norm distribution and clip anomalous spikes before they corrupt the weights. Practical recovery: checkpoint often, and on a spike, roll back to the last good checkpoint and skip/curate the offending data batch. ([noisy-data divergence study, arXiv 2602.02400](https://arxiv.org/html/2602.02400))

---

## 12. Distributed checkpointing, frameworks, and MFU

**Distributed checkpointing — save/resume without melting the cluster.** When state is sharded across thousands of GPUs, you cannot gather it all to rank-0 to `torch.save`. **PyTorch Distributed Checkpoint (DCP)** has each rank save **only its local shards** (handling `DTensor`/`ShardedTensor`) to multiple files — reducing memory peaks and parallelizing the write. **Asynchronous checkpointing** offloads the state-dict copy to a background thread so training resumes while the write drains to storage; torchtitan reports DCP async writes shrink checkpoint overhead 5–15×. DCP also supports **resharding on load** (save on N GPUs, resume on M). ([Distributed Checkpoint, PyTorch blog](https://pytorch.org/blog/distributed-checkpoint-efficient-checkpointing-in-large-scale-jobs/), [DCP async recipe](https://docs.pytorch.org/tutorials/recipes/distributed_async_checkpoint_recipe.html)) FSDP2's per-parameter DTensor sharding is what makes these checkpoints *communication-free* (§3).

**The training-framework landscape (2024–2026):** ([torchtitan repo](https://github.com/pytorch/torchtitan), [torchtitan ICLR 2025](https://proceedings.iclr.cc/paper_files/paper/2025/file/e6231c5f46598cfd09ff1970524e0436-Paper-Conference.pdf))

| Framework | Origin | Niche |
| --- | --- | --- |
| **PyTorch FSDP2 + `torch.distributed`** | Meta / PyTorch | The native building blocks; everything below wraps these |
| **torchtitan** | PyTorch | Minimal clean-room PyTorch-native **3D/4D parallelism** (FSDP2 + TP/async-TP + PP + CP) + DCP + `torch.compile`; reference for new techniques |
| **Megatron-LM / Megatron-Core** | NVIDIA | The canonical TP/PP/SP/CP/EP library; powers most frontier GPU training |
| **DeepSpeed** | Microsoft | ZeRO/ZeRO-Offload/ZeRO-Infinity; 3D parallelism; strong CPU/NVMe offload |
| **NVIDIA NeMo** | NVIDIA | Production stack on Megatron-Core + Transformer Engine + Lightning; recipes for Llama/DeepSeek; MLPerf reference |
| **MosaicML Composer / LLM Foundry** | Databricks | FSDP-based training with throughput tooling; migrated to FSDP2 |

**MFU — the efficiency yardstick.** **Model FLOPs Utilization (MFU)** is observed throughput ÷ the hardware's theoretical peak FLOPs — a hardware-agnostic measure of how well your training uses the silicon (introduced in Google's PaLM paper). The standard transformer training-FLOPs estimate is **≈ 6N per token** (N = parameters: 2N forward + 4N backward), so `MFU = 6 · N · tokens_per_sec / (num_GPUs · peak_FLOPs_per_GPU)`. ([MFU/FLOPs in LLM training](https://debjitpaul.github.io/blog/2025/compute/), [MegaScale, arXiv 2402.15627](https://arxiv.org/abs/2402.15627))

- **Reference values:** PaLM hit ~46% MFU; well-tuned large runs land **35–55%**; MegaScale reported **55.2% MFU training a 175B model on 12,288 GPUs**. Anything in that band is healthy; <30% usually means a comm-overlap, data-loading, or pipeline-bubble problem.
- **MFU vs HFU:** **Hardware** FLOPs Utilization (HFU) *counts* recomputed activations (gradient checkpointing) as useful work, so HFU > MFU when checkpointing is on. MFU measures *useful model* work; HFU measures *silicon busy-ness*. Track both — a high HFU but low MFU means you're burning FLOPs on recomputation.
- **Scaling efficiency** = throughput at N GPUs ÷ (N × single-GPU throughput). The job of all the overlap/placement work above is to keep this near 1.0 as the fleet grows to 10K+ GPUs.

---

## Practical patterns

- **Pick the smallest parallelism that fits, then add axes.** Single GPU → FSDP2/HSDP → +TP (inside node) → +PP (across nodes) → +CP (long context) / +EP (MoE). Don't reach for 3D parallelism before sharded data parallel is exhausted.
- **TP stays inside the node; PP/DP go across nodes.** Map the bandwidth-hungriest axis to NVLink. Crossing the node boundary with TP is the most common self-inflicted throughput wound.
- **BF16 by default; FP8 only after validating stability** on your model, keeping logits/norms higher-precision.
- **Tune the global batch via `grad_accum`**, and remember `no_sync()` on the non-final micro-steps under DDP/FSDP.
- **Checkpoint async and often** (DCP). A 30-minute checkpoint interval bounds the blast radius of a loss spike to 30 minutes.
- **Instrument MFU and tokens/sec from day one.** Treat an MFU regression as a P1 — it usually means overlap broke, not that the math changed.
- **Warmup + global-norm grad clip + z-loss/QK-norm** are the cheap insurance that keeps a long run alive.

## Anti-patterns

- **Using DDP for a model that needs sharding.** DDP scales batch, not model size; if it OOMs, you need FSDP/ZeRO, not a smaller batch.
- **Tensor parallelism across the inter-node network.** TP's per-block all-reduce will saturate the slow link and tank MFU. Keep TP ≤ GPUs-per-node.
- **Too few pipeline micro-batches.** Bubble ≈ `(P-1)/m`; a deep pipeline with few micro-batches can idle most GPUs. Raise `m` or use interleaved 1F1B.
- **FP16 for LLM pretraining in 2026.** Use BF16 — FP16's narrow range causes the loss-scaling instability BF16 was adopted to avoid.
- **Gathering sharded state to rank-0 to save.** OOMs at scale and serializes the write; use DCP sharded + async.
- **Confusing training and inference parallelism.** Inference TP/PP (`llm-inference-serving`) has no optimizer, no backward, no gradient sync — different constraints; don't copy a serving config into a training job.
- **Forgetting `no_sync()` during gradient accumulation**: turns every micro-step into a full gradient all-reduce and silently halves throughput.
- **Treating MoE expert count as just a parameter knob** without provisioning expert-parallel all-to-all bandwidth — the all-to-all becomes the bottleneck.

## Troubleshooting

- **OOM at step 0** → reduce micro-batch; enable activation checkpointing; raise FSDP shard degree / move ZeRO-2→3; offload optimizer (ZeRO-Offload) as a last resort.
- **OOM only on long sequences** → the problem is activation/attention memory, not weights → add context parallelism (§6) and selective recompute, not more DP.
- **Low MFU (<30%) with healthy GPUs** → comm not overlapping: check TP isn't crossing nodes, raise pipeline micro-batches, enable async-TP/prefetch, profile the collective wait time.
- **Loss spikes / NaN** → confirm BF16 (not FP16); add/verify warmup, global-norm clip, z-loss/QK-norm; roll back to last checkpoint and curate the offending batch; consider ZClip.
- **Slow / OOMing checkpoint** → switch to DCP sharded + async writes; don't `all_gather` to rank-0.
- **Throughput halves after adding grad accumulation** → missing `no_sync()` on non-final micro-steps.
- **NCCL hang / timeout** → mismatched collective call across ranks, a straggler/dead GPU, or NIC issue; set `NCCL_DEBUG=INFO`, check the topology, and look for a rank that never entered the collective.
- **FP8 run diverges where BF16 was fine** → tighten per-tensor scaling (Transformer Engine `DelayedScaling` amax history), keep sensitive layers (final projection, norms) in BF16/FP32.

## References

Primary docs & frameworks:
- PyTorch FSDP2 — [`fully_shard` API (2.9)](https://docs.pytorch.org/docs/2.9/distributed.fsdp.fully_shard.html), [FSDP2 tutorial](https://docs.pytorch.org/tutorials/intermediate/FSDP_tutorial.html), [FSDP1 vs FSDP2 (HF)](https://huggingface.co/docs/accelerate/concept_guides/fsdp1_vs_fsdp2), [Introducing FSDP (blog)](https://pytorch.org/blog/introducing-pytorch-fully-sharded-data-parallel-api/), [PyTorch FSDP paper, VLDB 2023](https://www.vldb.org/pvldb/vol16/p3848-huang.pdf)
- DeepSpeed ZeRO — [ZeRO tutorial](https://www.deepspeed.ai/tutorials/zero/), [ZeRO docs (GitHub)](https://github.com/microsoft/DeepSpeed/blob/master/docs/_tutorials/zero.md), [ZeRO-Infinity/zero3 readthedocs](https://deepspeed.readthedocs.io/en/latest/zero3.html), [ZeRO blog (MSR)](https://www.microsoft.com/en-us/research/blog/zero-deepspeed-new-system-optimizations-enable-training-models-with-over-100-billion-parameters/)
- Megatron — [Megatron-Core parallelism guide](https://docs.nvidia.com/megatron-core/developer-guide/latest/user-guide/parallelism-guide.html), [pipeline schedules source](https://github.com/NVIDIA/Megatron-LM/blob/main/megatron/core/pipeline_parallel/schedules.py), [Megatron-LM SC'21 paper](https://people.eecs.berkeley.edu/~matei/papers/2021/sc_megatron_lm.pdf), [Megatron-Bridge parallelisms](https://docs.nvidia.com/nemo/megatron-bridge/latest/parallelisms.html)
- torchtitan — [GitHub](https://github.com/pytorch/torchtitan), [FSDP doc](https://github.com/pytorch/torchtitan/blob/main/docs/fsdp.md), [ICLR 2025 paper](https://proceedings.iclr.cc/paper_files/paper/2025/file/e6231c5f46598cfd09ff1970524e0436-Paper-Conference.pdf)
- NVIDIA NeMo — [DeepSeek-V3 recipe](https://github.com/NVIDIA-NeMo/NeMo/blob/main/nemo/collections/llm/recipes/deepseek_v3.py), [Megatron-Bridge DeepSeek V3](https://docs.nvidia.com/nemo-framework/user-guide/latest/llms/deepseek_v3.html), [activation recomputation docs](https://docs.nvidia.com/nemo-framework/user-guide/24.09/nemotoolkit/features/optimizations/activation_recomputation.html)
- NCCL — [collectives docs](https://docs.nvidia.com/deeplearning/nccl/user-guide/docs/usage/collectives.html), [NCCL deep dive (NVIDIA blog)](https://developer.nvidia.com/blog/nccl-deep-dive-cross-data-center-communication-and-network-topology-awareness/)
- Distributed checkpoint — [PyTorch DCP blog](https://pytorch.org/blog/distributed-checkpoint-efficient-checkpointing-in-large-scale-jobs/), [DCP async recipe](https://docs.pytorch.org/tutorials/recipes/distributed_async_checkpoint_recipe.html)

Papers & analysis:
- Activation recomputation — [Reducing Activation Recomputation in Large Transformer Models, arXiv 2205.05198](https://arxiv.org/pdf/2205.05198)
- Context/sequence parallelism — [A Unified Sequence Parallelism Approach, arXiv 2405.07719](https://arxiv.org/html/2405.07719v3), [Seq1F1B, arXiv 2406.03488](https://arxiv.org/pdf/2406.03488)
- Mixed precision / FP8 — [FP8-LM, arXiv 2310.18313](https://arxiv.org/pdf/2310.18313), [FP8 vs BF16 trade-offs, arXiv 2411.08719](https://arxiv.org/html/2411.08719v1), [COAT FP8 training, arXiv 2410.19313](https://arxiv.org/pdf/2410.19313)
- Stability — [Methods of Improving LLM Training Stability, arXiv 2410.16682](https://arxiv.org/pdf/2410.16682), [ZClip, arXiv 2504.02507](https://arxiv.org/pdf/2504.02507), [Why warmup the LR, NeurIPS 2024](https://proceedings.neurips.cc/paper_files/paper/2024/file/ca98452d4e9ecbc18c40da2aa0da8b98-Paper-Conference.pdf)
- Scale & MFU — [MegaScale (>10K GPUs), arXiv 2402.15627](https://arxiv.org/abs/2402.15627), [MFU/FLOPs in LLM training](https://debjitpaul.github.io/blog/2025/compute/)

---

*Boundaries recap — what this reference deliberately defers:* inference/serving parallelism → `llm-inference-serving`; GPU kernels/CUDA/Triton/roofline → GPU-kernels reference; LoRA/QLoRA single-GPU fine-tuning → `llm-fine-tuning-peft`; attention/MoE/norm architecture → `transformer-architecture`; pretraining objectives, data mixtures, and scaling laws → the pretraining reference.
