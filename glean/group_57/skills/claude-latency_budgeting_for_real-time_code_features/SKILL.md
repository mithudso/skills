---
title: "Latency Budgeting for Real-Time Code Features"
aliases: []
description: >-
  Design real-time code features under strict latency budgets. SLAs by feature type (IDE autocomplete <500ms, code review <2s, refactoring <5s, CI/CD <30s). Model selection rubric under latency constraints: which models fit <1s? Quantization/batching/caching/speculation to meet SLAs. Measuring end-to-end latency (prompt encoding, inference, post-processing). Fallback strategies, timeouts, and degradation patterns.
trigger: >-
  latency budget for code feature, IDE autocomplete latency SLA, model selection under latency constraint, how fast LLM inference, quantization impact on latency, measuring end-to-end latency code agent, fallback strategy when model too slow, <500ms <2s <5s latency, code feature response time, real-time code generation latency, monitoring latency code features, timeout strategy LLM
expertise: operator, TAM
references:
  - references/latency-sla-table.md
  - references/model-latency-benchmarks.md
  - references/optimization-strategies.md
  - references/monitoring-patterns.md
hidden_refs: []
flags: []
---

# Latency Budgeting for Real-Time Code Features

You are designing a code-feature feature that needs real-time response. The user is an operator or TAM specifying an SLA; your job is to (1) assess feasibility under latency constraints, (2) recommend models by latency tier, (3) explain quantization/batching/caching trade-offs, (4) design monitoring and fallback strategies.

This skill is **code-feature specific** + **latency SLAs** + **model selection under constraints**. It does NOT cover generic infrastructure performance tuning or software-engineering patterns.

## When to Use

- **Feature designer specifying SLA:** "I need <500ms response for IDE autocomplete. Which models work?"
- **Operator evaluating quantization:** "Will int8 quantization get me below 2s latency?"
- **Latency requirement → model selection:** "Build me a latency-constrained model selector for code review feedback (<2s SLA)."
- **Troubleshooting latency miss:** "My code agent is hitting 3s per request. What's the bottleneck?"
- **Monitoring real-time features:** "How do I measure end-to-end latency (encoding + inference + post-proc)?"

## When NOT to Use

- **Generic performance tuning** (caching, CDN, infra optimization) → `software-engineering-patterns`
- **Code benchmark scores** (HumanEval, MBPP, correctness) → `code-benchmark-interpretation`
- **Choosing a model by capability alone** (without latency constraint) → `coding-model-selection`
- **LLM inference serving infrastructure** (vLLM, TensorRT setup) → `ai-llm-model-layer`

---

## Core Concepts

### Feature Latency SLAs (User Experience)

Real-time code features have distinct acceptable latencies based on user expectations:

| Feature | SLA Target | User Expectation | Notes |
|---------|-----------|-----------------|-------|
| IDE Autocomplete | <500ms | Keystroke feels instant; <100ms adds delight | P95 latency (not mean) |
| Inline hints/refactor preview | <1s | See suggestion while typing | Visible delay (~500ms+) feels sluggish |
| Code review feedback | <2s | Fast async check; acceptable to wait briefly | User can context-switch |
| Refactoring generation | <5s | User initiates, expects generation | Longer generation is tolerable if streamed |
| CI/CD build step | <30s | Build step timeout; failure is acceptable | Retry loop, not critical path |

**Latency is not mean—it's P95/P99.** A 100ms mean with 2s tail latency fails the SLA.

---

## Model Selection Under Latency Constraints

### Latency Tier Framework

Given your latency budget, **which models are viable?**

#### Tier 1: Sub-500ms (IDE Autocomplete)
Requires quantized small models or speculative decoding on frontier models.

- **Best fit:**
  - Copilot (Codex-series, proprietary quantization) — <200ms for 50-token completions
  - Claude 3.5 Haiku (int8) — 150–350ms
  - Qwen2.5-Coder-7B (int8, batched) — 100–200ms
  - Ollama DeepSeek Coder 6.7B (local, quantized) — 50–150ms

- **Achievable:** fp16 model <1B + speculative decoding, aggressive batching (8–16 tokens prefill)
- **Bottleneck:** Prompt encoding often 30–50% of latency; post-processing <5%
- **Trade-off:** Token budget is small (~50–100 tokens); quality drops if context is cut

**Decision:** For <500ms, **small quantized models are mandatory**. Frontier models cannot meet this SLA without speculative decoding + aggressive quantization.

#### Tier 2: <1s (Streaming Inline Hints)
Single-model forward pass with quantization.

- **Best fit:**
  - Claude 3.5 Haiku (fp16) — 200–400ms
  - Qwen2.5-Coder-32B (int8) — 300–600ms
  - Code Llama 34B (int8) — 400–700ms
  - Codestral (optimized for latency) — 300–500ms

- **Achievable:** fp16 models <20B params, aggressive batching (4–8 requests parallel)
- **Trade-off:** Quantization (int8) costs ~2–5% accuracy; batching adds 20–30ms for every concurrent request

**Decision:** <1s is **attainable with quantization**. Start with int8, measure, then micro-optimize.

#### Tier 3: <2s (Code Review Feedback)
Larger models become viable; async is acceptable.

- **Best fit:**
  - Claude 3.5 Sonnet (int8) — 600ms–1.5s
  - Qwen2.5-Coder-72B (int8) — 800ms–1.8s
  - DeepSeek Coder 33B (int8) — 700ms–1.5s
  - Mistral Medium (int8) — 800ms–1.8s

- **Achievable:** fp16 <70B, int8 <100B, streaming output
- **Trade-off:** Higher quality; latency is still acceptable (async)

**Decision:** <2s permits **mid-range quantized models**. Streaming output trades latency perception for responsiveness.

#### Tier 4: <5s (Refactor/Generation)
Full models, streaming strongly preferred.

- **Best fit:**
  - Claude 3.5 Sonnet (fp16) — 1.5–3s
  - Claude 3 Opus (int8) — 2–4s
  - Qwen2.5-Coder-72B (fp16) — 2–4s
  - DeepSeek Coder 236B (int8) — 3–5s

- **Achievable:** Full-size frontier models, streaming
- **Trade-off:** Latency high enough that **streaming is critical** for UX (user sees output as it arrives)

**Decision:** <5s is **streaming-or-bust**. Latency misses are recoverable via retry; streaming salvages UX.

#### Tier 5: <30s (CI/CD Build Step)
Any model; timeout-and-fail is acceptable.

- **Best fit:** Any frontier model (Opus, GPT-4, etc.), batched, cached
- **Trade-off:** Failure mode is retry or fallback; latency miss is not a feature break
- **Pattern:** Cache results per commit hash; parallelize across multiple agents if possible

**Decision:** <30s is **not a hard constraint**. Focus on correctness + timeout strategy.

---

### Latency Budget Allocation (End-to-End)

**Total latency = Prompt Encoding + Inference + Post-Processing + Network**

Typical breakdown for a code feature (example: autocomplete):
- **Prompt Encoding (20–40%):** Tokenize context, format prompt → 50–100ms for 2KB context
- **Inference (50–70%):** Model forward pass → 100–300ms for 50 tokens
- **Post-Processing (5–10%):** Parse output, rank suggestions → 10–20ms
- **Network (2–5%):** API roundtrip → 10–50ms if not local

**To meet SLA, allocate budget:**

| SLA | Encoding | Inference | Post-Proc | Network |
|-----|----------|-----------|-----------|---------|
| <500ms | 80ms | 300ms | 20ms | 100ms |
| <1s | 200ms | 600ms | 50ms | 150ms |
| <2s | 400ms | 1200ms | 100ms | 300ms |
| <5s | 1s | 3s | 200ms | 800ms |

**Action:** If your encoding is 200ms, you have 300ms for inference in a 500ms SLA. Optimize or adjust.

---

## Optimization Strategies (Latency Reduction)

### 1. Quantization (Model Size → Latency)

| Format | Size Reduction | Latency Improvement | Accuracy Loss | Use Case |
|--------|----------------|-------------------|---------------|----------|
| fp32 | 1x | — | — | Baseline |
| fp16 | 2x | 1.5–2x | <0.5% | Standard for GPU |
| int8 | 4x | 2.5–3x | 1–2% | Latency-critical |
| int4 | 8x | 3.5–4x | 2–5% | Extreme latency |
| nf4 (QLoRA) | 4x | 2.5–3x | 0.5–1% | Fine-tuning friendly |

**Empirical:** Quantizing a 70B model from fp16 to int8 reduces inference latency from ~2s to ~600ms, typically with <1% quality drop on code tasks.

**Trade-off:** int4 and nf4 show larger accuracy drops (2–5%) on code understanding tasks. Use only if latency is critical.

**Implementation:** Use GPTQ, AWQ, or GGUF quantization schemes. Avoid naive int8 (static quantization) for code models; use dynamic quantization instead.

### 2. Batching (Throughput → Latency)

Batching multiple requests together reduces per-token overhead.

**Latency cost of batching:**
- Batch size 1: 600ms
- Batch size 4: 600ms + 50ms per extra request = 750ms (125ms added per request)
- Batch size 8: 600ms + 100ms = 700ms (90ms added per request)
- Batch size 16: 600ms + 200ms = 800ms (80ms added per request)

**Decision:**
- Batch size **2–4** is sweet spot for most latency-critical features (adds 25–50ms).
- Batch size **8+** only if you have sufficient QPS (requests per second) to justify the added latency.
- Unbatched is only acceptable if requests are sparse (<1 req/sec).

**Implementation:** Use vLLM's `LoRAX` or `TensorRT-LLM` batching strategies.

### 3. Speculative Decoding (Frontier Model → Latency)

Use a small fast model to draft tokens; large model verifies/refines.

**Latency gain:**
- Claude 3.5 Sonnet alone: 2s for 200 tokens
- Claude 3.5 Sonnet + Haiku speculative decoding: 1.2s for 200 tokens (40% faster)

**Trade-off:** Quality preserved; small drafting model may hallucinate, but large model corrects.

**Use case:** Code review feedback (<2s SLA) or refactoring (<5s with streaming).

**Implementation:** Difficult to retrofit; requires API support (e.g., Claude API may support this in future) or self-hosted vLLM + medusa heads.

### 4. Caching (Context → Latency)

Prompt cache reduces encoding overhead for repeated context.

**Latency gain:**
- First request (no cache): 200ms encoding + 600ms inference = 800ms
- Cached requests (same context): 10ms cache hit + 600ms inference = 610ms (23% faster)

**Best for:** IDE features where file context is stable across multiple calls.

**Implementation:**
- Claude API: `cache_control: "ephemeral"` on system prompt or frequently-repeated context
- vLLM: KV cache for prefix tokens
- Local: Use `prefix_caching` in Ollama or similar

### 5. Streaming (Latency Perception)

Streaming output does NOT reduce latency—it **improves perceived latency** by showing output incrementally.

**Perceived latency with streaming:**
- User sees first token in 150ms (vs. 2000ms all-at-once)
- Total latency still 2000ms, but UX feels 8–10x faster

**Best for:** Any feature > 1s latency. Mandatory for >2s.

**Implementation:** Use SSE, WebSocket, or callback-based streaming. Update UI as tokens arrive.

### 6. Prefix Pruning (Context → Encoding)

Remove irrelevant context before encoding to reduce prompt size.

**Latency gain:**
- Full context (4KB): 150ms encoding
- Pruned context (1KB, 75% trimmed): 40ms encoding (3.75x faster)

**Trade-off:** Accuracy if pruned content was relevant.

**Implementation:**
- Semantic similarity search (embed recent context, keep top-K)
- Recency pruning (keep last N lines/blocks)
- Keyword filtering (keep context containing relevant identifiers)

**Best for:** IDE autocomplete where context is large but sparse relevance.

---

## Monitoring & Measuring End-to-End Latency

### Instrumentation Points

Measure latency at these key points to identify bottlenecks:

```
[Start]
  ↓ (prompt_encode_ms)
[Prompt Encoded]
  ↓ (queue_wait_ms)
[Queued for Inference]
  ↓ (inference_ms)
[Inference Complete]
  ↓ (post_process_ms)
[Output Ready]
```

**Minimal logging:**
```python
import time
import logging

logger = logging.getLogger("latency")

start = time.time()
prompt = encode_prompt(context)
t_encode = time.time() - start

# Queue if batching
queued = time.time()
result = model.generate(prompt)
t_queue = queued - start - t_encode
t_infer = time.time() - queued

output = parse_output(result)
t_post = time.time() - (queued + t_infer)

total = time.time() - start

logger.info(f"encode={t_encode*1000:.0f}ms queue={t_queue*1000:.0f}ms infer={t_infer*1000:.0f}ms post={t_post*1000:.0f}ms total={total*1000:.0f}ms")
```

### Percentile Reporting (Not Mean)

Report **P50, P95, P99** latency—not mean.

- **P50 (median):** 50% of requests are faster
- **P95 (95th percentile):** 95% of requests are faster (common SLA target)
- **P99 (99th percentile):** Catches tail latency (outliers)

**Example:**
```
latency_p50: 450ms
latency_p95: 580ms ← This is your SLA target
latency_p99: 800ms
```

If SLA is <500ms P95, and you measure P95=580ms, you **fail the SLA**.

### Dashboarding

Track over time:
- **Latency percentiles** (P50, P95, P99) by feature
- **Latency by component** (encode%, infer%, post%)
- **QPS (requests/sec)** and **batch size** correlation
- **Inference throughput** (tokens/sec)
- **Cache hit rate** (if caching is enabled)

Example Prometheus metrics:
```
code_feature_latency_seconds{feature="autocomplete", percentile="p95"}
code_feature_latency_seconds{feature="review", percentile="p95"}
model_inference_seconds{model="haiku", quantization="int8", percentile="p95"}
model_queue_seconds{batch_size="4"}
```

---

## Fallback Strategies & Timeouts

### Timeout Pattern

All real-time features need a timeout. Design graceful degradation:

```python
import asyncio

async def code_feature(context, sla_ms):
    try:
        result = await asyncio.wait_for(
            model.generate(context),
            timeout=sla_ms / 1000.0
        )
        return {"status": "ok", "output": result}
    except asyncio.TimeoutError:
        # Fallback: return cached suggestion or empty result
        return fallback_strategy(context)

def fallback_strategy(context):
    # Option 1: Cached previous result for this context
    if cached := lookup_cache(context):
        return {"status": "cached", "output": cached}
    
    # Option 2: Heuristic/rule-based fallback
    heuristic_output = apply_heuristic_rules(context)
    return {"status": "heuristic", "output": heuristic_output}
    
    # Option 3: Return "please wait" / retry UI signal
    return {"status": "timeout", "output": None}
```

### Fallback Strategies by Feature

| Feature | Timeout | Fallback |
|---------|---------|----------|
| Autocomplete | 500ms | Last seen completion, rule-based snippet |
| Inline hints | 1s | Cached hint for file, or none (optional feature) |
| Code review | 2s | Lightweight linter (fast rules), queue for async |
| Refactoring | 5s | Show "generating..." and stream when ready |
| CI/CD | 30s | Timeout, queue for retry, log for review |

### Retry Logic

For failed requests (timeout or error):

- **IDE autocomplete:** No retry (user moved on); cache miss is OK.
- **Code review:** Queue for async; notify user when result ready.
- **Refactoring:** Retry up to 2x with backoff; timeout and show "try again later."

---

## Practical Operator Checklist

### Before Launching a Real-Time Code Feature

- [ ] **Define SLA:** <500ms? <2s? What's acceptable for users?
- [ ] **Allocate latency budget:** Encoding (ms), inference (ms), post-proc (ms), network (ms)
- [ ] **Select model:** Use latency-tier framework; consider quantization
- [ ] **Measure baseline:** Deploy, measure P95 latency under production load
- [ ] **Identify bottleneck:** If not meeting SLA, which component is slow? (Encode? Infer? Queue?)
- [ ] **Optimize:** Quantize, cache, batch, or switch model
- [ ] **Monitor:** Set up alerts for P95 > SLA threshold; monitor per-component latency
- [ ] **Fallback:** Implement timeout + graceful degradation; test failure mode
- [ ] **Load test:** Verify P95 latency holds at 2x expected QPS

### Questions to Ask the Feature Designer

1. **What's the SLA?** (<500ms? <2s? <30s?)
2. **What's the acceptable quality drop for quantization?** (1%? 5%?)
3. **How much context can we afford to encode?** (1KB? 4KB?)
4. **Is streaming acceptable?** (UX improvement for >1s latency)
5. **What's the fallback strategy when model times out?** (Cached? Heuristic? Error?)
6. **What's the expected QPS?** (Helps decide batching strategy)

---

## References & Further Reading

- **Inference Benchmarks:** `references/model-latency-benchmarks.md` — Measured latencies for Haiku, Sonnet, Opus, Qwen, DeepSeek under quantization.
- **SLA Table:** `references/latency-sla-table.md` — Feature type, SLA, typical implementations.
- **Optimization Deep-Dive:** `references/optimization-strategies.md` — Quantization, batching, caching, speculative decoding with code examples.
- **Monitoring Patterns:** `references/monitoring-patterns.md` — Instrumentation, dashboarding, alerting, tail-latency debugging.

---

## Example: Designing an IDE Autocomplete Feature

**Requirement:** IDE autocomplete suggestions in <500ms (P95).

**Process:**

1. **SLA breakdown:**
   - Encoding: 80ms (tokenize 2KB context)
   - Inference: 300ms (50-token completion, small quantized model)
   - Post-proc: 20ms (rank suggestions)
   - Network: 100ms (API roundtrip)
   - **Total budget: 500ms** ✓

2. **Model selection:**
   - Qwen2.5-Coder-7B (int8) achieves 200–250ms inference for 50 tokens
   - Alternative: Claude 3.5 Haiku (int8) achieves 150–200ms
   - Decision: Haiku + caching for repeated context

3. **Optimizations:**
   - Prefix cache for repeated file context (saves 80ms on cache hit)
   - Batch size 4 (4 concurrent users) adds ~25ms; acceptable
   - Prune context to last 50 lines + relevant identifiers (reduce encoding to 60ms)

4. **Monitoring:**
   - Alert if P95 > 500ms or P99 > 700ms
   - Track encode/infer/post-proc breakdown daily

5. **Fallback:**
   - Timeout at 450ms (reserve 50ms for network/jitter)
   - Fallback: last 3 completions cached locally

6. **Deployment:**
   - A/B test with 10% users; measure P95 under real usage
   - Scale to 100% if P95 stays <500ms

This framework ensures the feature ships on-time and on-SLA.
