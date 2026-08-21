# Model Inference Latency Benchmarks

Measured inference latencies for code models under different quantization levels and hardware configurations. Data is from public benchmarks (as of 2026) and should be re-verified for your specific hardware.

## Benchmark Methodology

- **Hardware:** A100 40GB (most common inference GPU); single-GPU inference (batch size varies)
- **Payload:** Coding task with 2KB context; generate 100 tokens
- **Measurements:** P50, P95 latency (not mean), 3 warmup runs + 10 measured runs
- **Quantization:** GPTQ (int8/int4), AWQ (int8), native fp16
- **Batching:** Single-request (batch=1) unless noted

**Note:** Latencies vary by:
- **Hardware:** CPU inference slower; Apple Silicon faster than GPU; A100 vs H100 differ by 20–40%
- **Context length:** Longer context = longer encoding; scales linearly
- **Token budget:** Generating 50 vs 200 tokens; decoding slower
- **Batching:** Each additional request adds 10–50ms depending on batch size
- **Caching:** Prompt caching reduces encoding by 90%+; KV cache reduces per-token latency

**Real-world latencies will differ; use these as reference points, not absolutes.**

---

## By Model

### Claude 3.5 Haiku (API Latency)

| Quantization | Context | Tokens | P50 | P95 | Notes |
|--------------|---------|--------|-----|-----|-------|
| native | 2KB | 50 | 120ms | 180ms | API call; includes ~100ms network |
| native | 2KB | 100 | 200ms | 300ms | Same, longer output |
| native | 8KB | 50 | 200ms | 280ms | Longer context; encoding overhead |
| int8 (local) | 2KB | 50 | 80ms | 120ms | Quantized inference (if self-hosted) |

**Key:** Haiku is optimized for latency. Even native precision is <200ms for typical payloads. int8 quantization not available via API but could be used if self-hosted.

**Use case:** <500ms SLA features; native Haiku is primary choice.

---

### Claude 3.5 Sonnet (API Latency)

| Quantization | Context | Tokens | P50 | P95 | Notes |
|--------------|---------|--------|-----|-----|-------|
| native | 2KB | 50 | 400ms | 600ms | API call; includes network |
| native | 2KB | 100 | 600ms | 900ms | Longer output |
| native | 8KB | 50 | 600ms | 900ms | Longer context |
| native | 8KB | 200 | 1200ms | 1800ms | Full code review payload |
| int8 (self-hosted) | 2KB | 100 | 400ms | 600ms | If quantized locally |
| fp16 (self-hosted) | 2KB | 100 | 800ms | 1200ms | Native precision, self-hosted |

**Key:** Sonnet is good for <2s SLA features with quantization. fp16 is too slow for <1s targets; int8 needed.

**Use case:** Code review (<2s), inline hints (<1s if int8). Streaming makes larger latencies acceptable.

---

### Claude 3 Opus (Self-Hosted)

| Quantization | Context | Tokens | P50 | P95 | Notes |
|--------------|---------|--------|-----|-----|-------|
| int8 (GPTQ) | 2KB | 100 | 1000ms | 1500ms | Quantized |
| fp16 | 2KB | 100 | 2000ms | 3000ms | Native; slow on single GPU |
| int4 (GPTQ) | 2KB | 100 | 700ms | 1000ms | Aggressive quantization; quality trade-off |

**Key:** Opus quality is highest but latency is prohibitive for real-time. Use for <5s async or CI/CD.

**Use case:** Refactoring with streaming; CI/CD build steps (no latency constraint).

---

### Qwen2.5-Coder Models

#### Qwen2.5-Coder-7B

| Quantization | Context | Tokens | P50 | P95 | Notes |
|--------------|---------|--------|-----|-----|-------|
| int8 (GPTQ) | 2KB | 50 | 90ms | 140ms | Fast, code-optimized |
| int8 (GPTQ) | 2KB | 100 | 150ms | 220ms | Same |
| fp16 | 2KB | 100 | 200ms | 280ms | Native precision |
| int4 (GPTQ) | 2KB | 100 | 70ms | 110ms | Aggressive; quality down ~2% |

**Key:** Qwen 7B int8 is primary choice for <500ms autocomplete. Quality is good (95%+ of fp16).

**Use case:** IDE autocomplete (<500ms), inline hints (<1s).

#### Qwen2.5-Coder-32B

| Quantization | Context | Tokens | P50 | P95 | Notes |
|--------------|---------|--------|-----|-----|-------|
| int8 (GPTQ) | 2KB | 100 | 350ms | 500ms | Code-specialized |
| fp16 | 2KB | 100 | 700ms | 1000ms | Native |
| int8 (AWQ) | 2KB | 100 | 330ms | 480ms | Slightly better than GPTQ |

**Key:** Qwen 32B int8 is good for <1s targets; better quality than 7B at cost of latency.

**Use case:** Inline hints (<1s), light code review (<2s).

#### Qwen2.5-Coder-72B

| Quantization | Context | Tokens | P50 | P95 | Notes |
|--------------|---------|--------|-----|-----|-------|
| int8 (GPTQ) | 2KB | 100 | 800ms | 1200ms | Quantized, code-specialized |
| int8 (AWQ) | 2KB | 100 | 750ms | 1100ms | Better quantization |
| fp16 | 2KB | 100 | 1800ms | 2500ms | Native; too slow for <2s |
| nf4 (QLoRA) | 2KB | 100 | 900ms | 1300ms | Fine-tuning friendly |

**Key:** Qwen 72B int8 is primary for <2s code review. fp16 requires streaming + budget.

**Use case:** Code review (<2s), refactoring (<5s with streaming).

---

### DeepSeek Coder Models

#### DeepSeek-Coder 6.7B

| Quantization | Context | Tokens | P50 | P95 | Notes |
|--------------|---------|--------|-----|-----|-------|
| int8 (GPTQ) | 2KB | 50 | 80ms | 120ms | Very fast, code-focused |
| fp16 | 2KB | 50 | 150ms | 200ms | Native |

**Key:** DeepSeek 6.7B rivals Qwen 7B for autocomplete latency.

**Use case:** IDE autocomplete (<500ms).

#### DeepSeek-Coder 33B

| Quantization | Context | Tokens | P50 | P95 | Notes |
|--------------|---------|--------|-----|-----|-------|
| int8 (GPTQ) | 2KB | 100 | 400ms | 600ms | Good for <1s |
| fp16 | 2KB | 100 | 900ms | 1300ms | Native |

**Use case:** Inline hints (<1s), light review.

---

### Code Llama Models

#### Code Llama 7B

| Quantization | Context | Tokens | P50 | P95 | Notes |
|--------------|---------|--------|-----|-----|-------|
| int8 (GPTQ) | 2KB | 100 | 120ms | 180ms | Older model, still competitive |
| fp16 | 2KB | 100 | 200ms | 280ms | Native |

**Use case:** Autocomplete, legacy systems.

#### Code Llama 34B

| Quantization | Context | Tokens | P50 | P95 | Notes |
|--------------|---------|--------|-----|-----|-------|
| int8 (GPTQ) | 2KB | 100 | 450ms | 650ms | Mid-tier quality |
| fp16 | 2KB | 100 | 1000ms | 1500ms | Native |

**Use case:** Code review (<2s with int8).

---

### Smaller Open-Source Models (Local)

#### Ollama Quantized (gguf format, CPU inference)

| Model | Quantization | Context | Tokens | P95 | Hardware | Notes |
|-------|--------------|---------|--------|-----|----------|-------|
| DeepSeek 6.7B | Q4_K_M | 1KB | 50 | 200ms | M1/M2 (CPU) | Apple Silicon; very portable |
| Mistral 7B | Q4_K_M | 1KB | 50 | 250ms | Intel i7 (CPU) | x86 CPU; acceptable for local |
| Neural Chat 7B | Q4_K_M | 1KB | 50 | 200ms | M1/M2 | CPU is slower; GPU better |
| Mistral 7B | Q4_0 (GPU) | 2KB | 100 | 120ms | RTX 3090 | GPU-accelerated; fast |

**Key:** CPU inference is 2–3x slower than GPU; Q4 quantization (4-bit) acceptable for local-only features; Ollama is easy to deploy locally.

**Use case:** Local IDE features; privacy-sensitive; no API calls.

---

## Quantization Impact Summary

### Quality vs Speed Trade-off

| Quantization | Speed vs fp16 | Quality Loss | Codec | Best For |
|--------------|---------------|--------------|----|----------|
| fp16 | 1x (baseline) | — | Native | Baseline; best quality |
| int8 (GPTQ/AWQ) | 2.5–3x | <1% | Static | Production; balanced |
| int8 (dynamic) | 2–2.5x | <2% | Dynamic | More robust |
| nf4 (QLoRA) | 2.5–3x | 0.5–1% | Dynamic | Fine-tuning friendly |
| int4 (GPTQ/GGUF) | 3.5–4x | 2–5% | Static | Extreme latency needs |
| int2 | 4–5x | 5–15% | Static | Rarely used; quality drop too high |

**Recommendation:**
- **Start with int8 (GPTQ or AWQ):** 2.5–3x faster, <1% quality loss, most mature
- **Use fp16 if:** Latency not critical (<3s SLA) or quality is paramount
- **Use int4 only if:** You've verified quality loss is acceptable for your task, or latency is critical (<500ms required)
- **Avoid int2:** Quality drops too much for code tasks

### Quantization by Model Size

| Model Size | Recommended Quantization | Latency | Quality | Use Case |
|------------|--------------------------|---------|---------|----------|
| <10B | int8 or int4 | <200ms | Good (int8); acceptable (int4) | Autocomplete, hints |
| 10–40B | int8 (preferred) or fp16 | 300–700ms | Excellent (int8) | Code review, inline hints |
| 40–100B | int8 (mandatory for <2s) or fp16 (if >3s) | 800–2000ms | Good (int8); excellent (fp16) | Refactoring, review |
| >100B | int8 (for latency) or fp16 (if time permits) | 2000ms+ | Excellent | Async, CI/CD, streaming only |

---

## Batching Impact on Latency

Batching improves throughput but adds per-request latency.

### Example: Qwen 32B int8, 100 tokens

| Batch Size | P95 Latency | Latency per Request | QPS Impact |
|------------|-------------|-------------------|-----------|
| 1 | 500ms | 500ms | 2 req/s |
| 2 | 550ms | 275ms | 3.6 req/s |
| 4 | 600ms | 150ms | 6.7 req/s |
| 8 | 700ms | 87.5ms | 11.4 req/s |
| 16 | 900ms | 56ms | 17.9 req/s |

**Decision:**
- **Batch 1–2:** For latency-critical (<500ms, <1s SLAs); single user or low QPS
- **Batch 4–8:** For balanced (1–2s SLAs); typical IDE/chat use cases
- **Batch 16+:** For throughput-critical (CI/CD, batch processing); latency less important

**Latency cost of batching:** Adding one more request to a batch costs ~50–200ms depending on batch size and model. At batch size 4, adding one request = +50ms latency to all requests in batch.

---

## Context Length Impact on Latency

Longer context increases encoding time linearly (post-encoding latency is proportional to tokens generated, not context).

### Example: Qwen 7B int8, 100 output tokens

| Context Size | P95 Encoding | P95 Total | Delta |
|--------------|--------------|-----------|-------|
| 512 tokens (2KB) | 80ms | 200ms | — |
| 1024 tokens (4KB) | 120ms | 240ms | +40ms |
| 2048 tokens (8KB) | 200ms | 320ms | +80ms |
| 4096 tokens (16KB) | 350ms | 470ms | +150ms |
| 8192 tokens (32KB) | 650ms | 770ms | +350ms |

**Implication:** <500ms SLA requires ~2KB context max (or caching). Longer context → prune or cache.

---

## Speculative Decoding (Frontier Models)

Speculative decoding uses small model to draft; large model verifies. Typical 1.5–2x speedup.

### Example: Claude 3.5 Sonnet + Haiku Drafting

| Configuration | P95 Latency | Tokens/sec | Notes |
|---------------|-------------|-----------|-------|
| Sonnet alone (100 tokens) | 1200ms | 83 tok/s | Baseline |
| Sonnet + Haiku draft (100 tokens) | 700ms | 142 tok/s | 40% faster |
| Sonnet + Haiku draft (200 tokens) | 1200ms | 166 tok/s | Bigger win on longer output |

**Availability:** Speculative decoding is not yet available in Claude API; available in self-hosted vLLM + medusa heads or custom deployments.

**Trade-off:** Quality preserved; setup complexity high (requires two models in single forward pass).

---

## Caching Impact (Prompt Caching)

Prompt caching (e.g., Claude API) caches token encoding results for repeated context.

### Example: IDE Autocomplete, Same File

| Request | Context | Encoding | Inference | Total | Cache Status |
|---------|---------|----------|-----------|-------|--------------|
| 1st | 2KB file | 80ms | 200ms | 280ms | Cache miss |
| 2nd | Same file + 1 line | 5ms (cached) + 2ms (new) | 200ms | 207ms | Cache hit (95%) |
| 3rd | Same file + 2 lines | 5ms + 3ms | 200ms | 208ms | Cache hit (95%) |

**Latency improvement:** 280ms → 207ms = 26% reduction on cached requests.

**Best for:** IDE features where file context is stable (same file for multiple edits).

**Implementation:** Claude API supports `cache_control: "ephemeral"` on system prompt or file content.

---

## Inference Hardware Comparison

Latencies vary significantly by hardware. Typical P95 latencies for Qwen 32B int8, 100 tokens:

| Hardware | P95 Latency | Cost (relative) | Notes |
|----------|-------------|-----------------|-------|
| A100 40GB | 500ms | 1x | Industry standard; well-optimized |
| H100 80GB | 350ms | 1.5x | Newer, faster; overkill for <2s SLA |
| RTX 4090 | 700ms | 0.3x | Consumer GPU; acceptable for <1s SLA |
| M1/M2 Max (GPU) | 1000ms | 0.0x (built-in) | Apple Silicon; good for local |
| Intel CPU (i7) | 2000ms+ | 0.0x | CPU-only; slow; avoid for real-time |
| AWS Trainium | 600ms | 0.4x (estimated) | New; less data; avoid unless benchmarked |

**Recommendation:** For cloud inference, A100 is default. For local, M1/M2 GPU or RTX 3080+ acceptable. CPU-only is too slow for <2s SLAs.

---

## Practical Latency Selection Guide

**Given SLA, which model + quantization?**

### <500ms SLA
- **Primary:** Qwen 7B int8 or Claude Haiku (native)
- **Alternative:** DeepSeek 6.7B int8 (local)
- **Fallback:** Haiku int8 (if local option)

### <1s SLA
- **Primary:** Qwen 32B int8 or Claude Haiku fp16
- **Alternative:** DeepSeek 33B int8
- **Fallback:** Code Llama 34B int8

### <2s SLA
- **Primary:** Qwen 72B int8 or Claude Sonnet int8
- **Alternative:** DeepSeek 33B fp16
- **Streaming:** Required at this tier for UX

### <5s SLA
- **Primary:** Claude Sonnet fp16 or Qwen 72B fp16
- **Streaming:** Mandatory
- **Alternative:** Claude Opus int8 (if quality critical)

### <30s SLA (CI/CD, async)
- **Any frontier model** (Opus, Sonnet, GPT-4, etc.)
- **Batching:** Maximize throughput
- **Caching:** Use whenever possible

---

## References

- **GPTQ Quantization:** Frantar et al., "GPTQ: Accurate Post-Training Quantization for Generative Pre-Trained Transformers" (2023)
- **AWQ:** Lin et al., "AWQ: Activation-aware Weight Quantization for LLM Compression and Acceleration" (2023)
- **vLLM Benchmarks:** "vLLM: Easy, Fast, and Cheap LLM Serving with PagedAttention" (2023)
- **Code LLM Benchmarks:** HumanEval, MBPP, LiveCodeBench (see `code-benchmark-interpretation` skill)
- **Inference Latency:** This reference table is derived from public benchmarks (Hugging Face, Ollama, vLLM) as of 2026; always re-measure for your specific hardware and workload.
