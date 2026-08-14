# Latency SLA Reference Table

Canonical SLAs for real-time code features, based on user perception studies and industry practice.

## Feature Type × SLA × Implementation

| Feature | SLA P95 | User Expectation | Implementation Pattern | Typical Models |
|---------|---------|-----------------|----------------------|-----------------|
| **IDE Autocomplete** | <500ms | Keystroke feels instant; <100ms ideal, <200ms good, >300ms noticeably slow | Quantized small model (int8, <10B), prefix cache, batch=1–2 | Haiku (int8), Qwen 7B (int8), Copilot Codex |
| **Inline Hints/Refactor Preview** | <1000ms | Suggestion appears while still typing; acceptable to wait 500–1000ms | Quantized model (int8, 7–20B) + streaming, batch=2–4 | Haiku (fp16), Qwen 32B (int8), Codestral |
| **Code Review/Feedback** | <2000ms | Fast async check; user doesn't need to see instant result, but <2s feels responsive | Mid-size quantized model (int8, 30–70B), batch=4–8, streaming | Sonnet (int8), Qwen 72B (int8), DeepSeek 33B (int8) |
| **Refactoring/Generation** | <5000ms | User initiates action; generation latency acceptable if streamed | Full-size model (fp16, >70B) or frontier, batch=1–4, **streaming mandatory** | Sonnet (fp16), Opus (int8), Claude 3.5 Sonnet, Qwen 72B (fp16) |
| **CI/CD Build Step** | <30000ms | Build step timeout; failure is acceptable (retry); latency is soft constraint | Any model, cached when possible, batch=10+ if parallel | Opus, GPT-4, Sonnet, any frontier |
| **Async Analysis (Nightly/PR)** | Unbounded | No user wait; result is eventual | Large model (any size), batch=100+ | Any frontier model, optimized for throughput not latency |

## Latency Perception Thresholds

Human-computer interaction studies (Nielsen, Card) establish latency perception:

- **<100ms:** Feels instantaneous; user perceives no delay
- **100–300ms:** Feels responsive; acceptable for "fast" UI
- **300–1000ms:** Noticeable delay; user feels they're waiting
- **1000–3000ms:** Slow but tolerable; user may context-switch
- **>3000ms:** Very slow; user assumes process is stuck (retry or abandon)

**Code feature implication:** IDE autocomplete must hit <300ms to feel "instant"; code review can tolerate <2s because it's async.

## P95 vs Mean Latency

**Critical:** SLAs are specified as P95 or P99, not mean.

- **P95:** 95% of requests faster than this latency; 5% slower
- **P99:** 99% of requests faster than this latency; 1% slower (tail latency)

A feature with mean 400ms but P95 1200ms **fails a <500ms SLA**, even though the mean suggests it should pass.

**Why P95?** A single slow request ruins the user experience. Mean latency hides tail latency.

## Latency Budget Allocation

For a given SLA, allocate latency across components:

### <500ms SLA (IDE Autocomplete)
```
Encoding:        80ms   (80% quality: prune to 1KB recent context)
Inference:      300ms   (small quantized model, 50 tokens)
Post-processing: 20ms   (rank, filter)
Network:        100ms   (API roundtrip or local latency)
─────────────────────
Total:          500ms
```

### <1000ms SLA (Inline Hints)
```
Encoding:       200ms   (tokenize 2–4KB context)
Inference:      600ms   (quantized <20B model, 100 tokens)
Post-processing: 50ms   (format, rank)
Network:        150ms
─────────────────────
Total:         1000ms
```

### <2000ms SLA (Code Review)
```
Encoding:       400ms   (4–8KB context, comment + diff)
Inference:     1200ms   (30–70B quantized model, 200+ tokens)
Post-processing:100ms   (parse, structure)
Network:        300ms
─────────────────────
Total:         2000ms
```

### <5000ms SLA (Refactoring)
```
Encoding:      1000ms   (8KB+ context, full file)
Inference:     3000ms   (frontier model, 500+ tokens, streaming)
Post-processing:200ms   (validation, formatting)
Network:        800ms
─────────────────────
Total:         5000ms
```

**Action:** If your measured encoding is 300ms but budget is 200ms, either:
1. Prune context size
2. Add caching to skip encoding on repeated context
3. Increase SLA
4. Add encoding to critical path (parallelize)

## Model Suitability by SLA

| SLA | Suitable Models | Constraints | Notes |
|-----|-----------------|-------------|-------|
| <500ms | Haiku (int8/fp16), Qwen 7B (int8), local quantized <6B | Quantization essential; small context | Speculative decoding cannot help (frontier models are baseline) |
| <1s | Haiku (fp16), Qwen 32B (int8), Codestral (int8) | int8 quantization, batch ≤4 | Streaming improves perception |
| <2s | Sonnet (int8), Qwen 72B (int8), DeepSeek 33B (int8), Code Llama 34B (int8) | int8 quantization, batch ≤8, streaming | Mid-size models + quantization + batching |
| <5s | Sonnet (fp16), Opus (int8), Qwen 72B (fp16), Claude 3.5 Sonnet | Streaming essential; batch ≤4 | Latency misses recoverable via streaming UX |
| <30s | Any frontier model (Opus, GPT-4, Claude 3.5 Sonnet) | Batch large, cache when possible | Timeout acceptable; retry on failure |

## Latency Levers (Ranked by Impact)

Ordered by typical latency reduction when applied:

| Lever | Latency Reduction | Effort | Risk | Best For |
|-------|-------------------|--------|------|----------|
| **Switch to smaller model** | 2–4x | Low | Medium (quality) | <500ms or <1s targets |
| **Quantization (int8)** | 2.5–3x | Low | Low (<1% quality) | All tiers with current model |
| **Prefix cache** | 1.5–2x on cached requests | Medium | Very low | Stable context (IDE files) |
| **Context pruning** | 1.2–2x | Low | Medium (accuracy) | Encoding-bound features |
| **Batching** | 1.1–1.2x per request (added latency) | Medium | Low | High QPS (>10 req/s) |
| **Speculative decoding** | 1.5–2x on inference | High | Low | Frontier models only |
| **Streaming output** | 0x (no latency reduction, UX only) | Medium | None | Perception improvement >1s |
| **Offload to smaller model (drafting)** | 1.5–2x | High | Medium | Expert use case |

**Decision rule:**
1. **First:** Check if quantization (int8) gets you there (2–3x speedup, low risk)
2. **Second:** If quantization insufficient, switch to smaller model (int8)
3. **Third:** Add caching or context pruning
4. **Fourth:** Only then consider batching or speculative decoding

---

## Example: SLA Verification Checklist

**Scenario:** You ship an IDE autocomplete feature. SLA is <500ms P95.

**Week 1 (Launch):**
- [ ] Measure production latency; P95 = 480ms ✓ **Pass**
- [ ] Log per-component: encode 90ms, infer 280ms, post 20ms, network 90ms
- [ ] Monitor alerting threshold: P95 > 520ms (5% margin)

**Week 3 (Scale):**
- [ ] QPS increases 5x; P95 latency rises to 580ms ✗ **Fail**
- [ ] Root cause: Batching (batch size grew from 2 to 8)
- [ ] Action: Enable prefix cache (skip encoding on repeated context)
- [ ] Result: P95 drops to 420ms ✓ **Pass**

**Week 8 (New model):**
- [ ] New Qwen Coder 3.5 released; claims 20% faster
- [ ] Measure: P95 = 350ms; quality unchanged ✓
- [ ] Migrate; save infrastructure costs

---

## SLA Trade-offs (Qualitative)

| Feature | Latency | Quality | Cost | Recommendation |
|---------|---------|---------|------|-----------------|
| Autocomplete | <500ms (tight) | Medium (quantized) | High (small model, real-time) | Necessary; trade-off for UX |
| Code review | <2s (comfortable) | High (larger model) | Medium | Sweet spot; good UX + quality |
| Refactoring | <5s (loose) | Very high (streaming) | Low (batched) | Streaming makes latency irrelevant |

---

## Further Reading

- **Bakshy et al.** (Facebook): "The Anatomy of a Large-Scale Hypertextual Web Search Engine" — establishes that page load delays reduce engagement (every 100ms delay = 1% drop in engagement)
- **Nielsen, "Response Times":** <100ms felt instantaneous; 100–300ms felt responsive; >1s lost user attention
- **GitHub Copilot case study:** Autocomplete at <300ms baseline; misses drop adoption significantly
- **Claude API latency best practices:** Quantization + caching can 3x throughput without quality loss
