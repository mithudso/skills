# Latency Optimization Strategies

Detailed techniques to reduce latency for real-time code features, with code examples and trade-off analysis.

---

## 1. Quantization (Model Compression)

Quantization reduces model size and speeds up inference by using lower-precision arithmetic (int8, int4) instead of float32 or float16.

### Quantization Schemes

#### GPTQ (Static Quantization)
- **What:** Post-training quantization; uses calibration dataset to determine scale/zero-point per layer
- **Speed:** 2.5–3x faster (int8); 3.5–4x (int4)
- **Quality:** <1% drop (int8); 2–5% drop (int4)
- **Implementation:** HuggingFace Transformers + `auto_gptq` library
- **Best for:** Production; most stable

```python
from auto_gptq import AutoGPTQForCausalLM

model = AutoGPTQForCausalLM.from_pretrained(
    "Qwen/Qwen2.5-Coder-32B-Instruct-GPTQ-Int8",
    device_map="auto"
)
# ~350ms P95 latency for 100 tokens vs 700ms fp16
```

#### AWQ (Activation-Aware Quantization)
- **What:** Quantization-aware calibration; accounts for activation distributions
- **Speed:** 2.5–3x faster (similar to GPTQ)
- **Quality:** Marginally better than GPTQ (0.5–1% advantage)
- **Implementation:** `autoawq` library
- **Best for:** Higher quality when int8 is needed

```python
from awq import AutoAWQForCausalLM

model = AutoAWQForCausalLM.from_quantized(
    "Qwen/Qwen2.5-Coder-32B-Instruct-AWQ"
)
# ~330ms P95 vs 350ms GPTQ (marginal gain; GPTQ preferred for ease)
```

#### GGUF (Quantization Format for CPU)
- **What:** Portable quantization format; runs on CPU with minimal dependencies
- **Speed:** Slow on CPU (200–500ms for 7B); acceptable for local dev
- **Quality:** Similar to GPTQ
- **Implementation:** `ollama`, `llama.cpp`
- **Best for:** Local/privacy-sensitive; no GPU required

```bash
ollama pull mistral:7b-instruct-q4_K_M
# Runs locally; ~200ms P95 on M1/M2
```

#### NF4 (Normalized Float 4)
- **What:** QLoRA-compatible 4-bit quantization; preserves fine-tuning capability
- **Speed:** 2.5–3x faster
- **Quality:** 0.5–1% drop; better than GPTQ int4
- **Implementation:** `bitsandbytes` library
- **Best for:** Fine-tuning quantized models

```python
from bitsandbytes.nn import Linear4bit
from transformers import AutoModelForCausalLM

model = AutoModelForCausalLM.from_pretrained(
    "Qwen/Qwen2.5-Coder-32B",
    load_in_4bit=True,
    bnb_4bit_quant_type="nf4",
)
# ~380ms P95; allows fine-tuning
```

### Quantization Decision Tree

```
Is latency <500ms required?
  → YES: Use int8 (GPTQ/AWQ) on small model (<10B) OR int4 on smaller model
  → NO: Is latency <1s?
    → YES: Use int8 on mid-size (<32B)
    → NO: Is latency <2s?
      → YES: Use int8 on large (32–70B)
      → NO: Use fp16 or fp32 (no quantization)

For production code features:
  → Use GPTQ or AWQ (most stable)
For local/dev:
  → Use GGUF or NF4 (easier to manage)
For fine-tuning:
  → Use NF4 (preserves gradient computation)
```

### Quality Impact by Model Size

Quantization loss varies by model size:

| Model Size | int8 Loss | int4 Loss | Recommendation |
|------------|-----------|-----------|-----------------|
| <7B | <0.5% | 1–2% | int8 preferred; int4 acceptable |
| 7–32B | <1% | 2–3% | int8 for production; int4 for extreme latency |
| 32–70B | 1–2% | 3–5% | int8 nearly required; int4 risky |
| >70B | 2–3% | 5–10% | int8 mandatory; int4 not recommended |

**Action:** For code tasks, int8 is the default. Only use int4 if you've measured quality loss and it's acceptable.

---

## 2. Caching (Context Encoding)

Prompt caching skips re-encoding repeated context, saving 50–80% of encoding latency.

### Types of Caching

#### Prompt Caching (API)
- **Service:** Claude API supports `cache_control: "ephemeral"`
- **Latency savings:** 50–90% on encoding (if cached)
- **Cost:** Cached tokens cost ~25% of regular tokens (cheaper + faster)
- **TTL:** 5 minutes (ephemeral) or 24 hours (static)

```python
import anthropic

client = anthropic.Anthropic()

# First request: encodes full context (no cache)
message = client.messages.create(
    model="claude-3-5-sonnet-20241022",
    max_tokens=200,
    system=[
        {
            "type": "text",
            "text": "You are a code reviewer."
        },
        {
            "type": "text",
            "text": LARGE_FILE_CONTENT,  # 4KB file
            "cache_control": {"type": "ephemeral"}
        }
    ],
    messages=[{"role": "user", "content": "Review line 42."}]
)
# Latency: 200ms encoding + 600ms inference = 800ms

# Second request: uses cache (same file content)
message = client.messages.create(
    model="claude-3-5-sonnet-20241022",
    max_tokens=200,
    system=[
        {
            "type": "text",
            "text": "You are a code reviewer."
        },
        {
            "type": "text",
            "text": LARGE_FILE_CONTENT,  # Same 4KB file
            "cache_control": {"type": "ephemeral"}
        }
    ],
    messages=[{"role": "user", "content": "Review line 50."}]
)
# Latency: 10ms cache lookup + 600ms inference = 610ms
# 23% reduction in latency
```

**Limitations:**
- Minimum cache size: 1024 tokens (smaller context won't cache)
- Cache hits only if context is identical (byte-for-byte)
- Ephemeral TTL is 5 minutes; useful for IDE features in same session

#### KV Cache (Self-Hosted)
- **Service:** vLLM, TensorRT-LLM, Ollama
- **Mechanism:** Cache key/value tensors from encoder layers; skip re-computation
- **Latency savings:** 30–50% on encoding
- **Implementation:** Automatic in most inference servers

```python
# vLLM with KV cache
from vllm import LLM, SamplingParams

llm = LLM(
    model="Qwen/Qwen2.5-Coder-32B",
    tensor_parallel_size=1,
    dtype="float16",
    enable_prefix_caching=True,  # Enable prefix caching
)

# First prompt
prompt1 = "You are a code reviewer.\n" + LARGE_FILE + "\nReview line 42."
output1 = llm.generate(prompt1, sampling_params)
# Latency: 200ms encoding + 600ms inference

# Second prompt (same file, different query)
prompt2 = "You are a code reviewer.\n" + LARGE_FILE + "\nReview line 50."
output2 = llm.generate(prompt2, sampling_params)
# Latency: 10ms (cache hit) + 600ms inference = 610ms (23% reduction)
```

**Benefits:** Prefix caching works for any model; automatic in most servers.

#### Semantic Cache (Custom)
- **Mechanism:** Cache embeddings of similar prompts; reuse results
- **Latency savings:** Varies (5–90% depending on cache hit rate)
- **Trade-off:** Approximate (cached result may not be exact)
- **Implementation:** Custom with vector DB (FAISS, Milvus)

```python
from sentence_transformers import SentenceTransformer
import faiss

embedder = SentenceTransformer("all-MiniLM-L6-v2")
cache_embeddings = []
cache_results = []

def semantic_cache_get(query, threshold=0.95):
    query_emb = embedder.encode([query])[0]
    index = faiss.IndexFlatIP(384)
    index.add(np.array(cache_embeddings))
    distances, indices = index.search(np.array([query_emb]), k=1)
    if distances[0][0] > threshold:
        return cache_results[indices[0][0]]
    return None

def semantic_cache_set(query, result):
    query_emb = embedder.encode([query])[0]
    cache_embeddings.append(query_emb)
    cache_results.append(result)

# Usage
query1 = "Review the user authentication logic"
result1 = semantic_cache_get(query1)
if not result1:
    result1 = model.generate(query1)
    semantic_cache_set(query1, result1)

# Similar query hits cache
query2 = "Check the login code"
result2 = semantic_cache_get(query2, threshold=0.90)
# ~90% similarity; cache hit; latency ~10ms
```

**Best for:** Code review features where queries are similar but not identical.

### Caching Decision

| Cache Type | Latency Saving | Accuracy | Best For |
|-----------|-----------------|----------|----------|
| Prompt cache (API) | 50–90% encoding | 100% (identical context) | IDE with repeated file context |
| KV cache (self-hosted) | 30–50% encoding | 100% (identical prefix) | Same as prompt cache |
| Semantic cache | 90%+ (if hit) | 95–99% (approximate) | Code review with similar queries |
| No cache | 0% | — | One-off features |

**Action:** Implement prompt/KV caching for IDE features; semantic caching for review features if QPS is high.

---

## 3. Context Pruning (Reduce Prompt Size)

Shorten context before encoding to reduce latency.

### Pruning Strategies

#### Recency Pruning (Keep Recent N Lines)
- **Mechanism:** Keep only last N lines of file + surrounding functions
- **Latency saving:** 50–80% encoding reduction
- **Trade-off:** Accuracy if pruned context is relevant

```python
def prune_by_recency(file_content, max_lines=50):
    lines = file_content.split('\n')
    if len(lines) > max_lines:
        return '\n'.join(lines[-max_lines:])
    return file_content

# 4KB file → 1KB context → 80ms encoding vs 150ms
```

#### Semantic Pruning (Keep Relevant Context)
- **Mechanism:** Embed query, embed context chunks, keep top-K relevant
- **Latency saving:** 40–70% encoding reduction
- **Trade-off:** Adds embedding overhead (~50ms); net savings if context is large

```python
from sentence_transformers import SentenceTransformer

embedder = SentenceTransformer("all-MiniLM-L6-v2")

def prune_by_relevance(file_content, query, max_tokens=1024):
    lines = file_content.split('\n')
    chunks = [line for line in lines if line.strip()]
    
    query_emb = embedder.encode(query)
    chunk_embs = embedder.encode(chunks)
    
    scores = np.dot(chunk_embs, query_emb)
    top_indices = np.argsort(scores)[-10:]  # Keep top 10 chunks
    
    pruned = '\n'.join([chunks[i] for i in top_indices])
    return pruned[:max_tokens]

# 8KB file → 2KB pruned context
# Embedding: 50ms, Encoding: 80ms (vs 200ms full) = net saving 70ms
```

#### Identifier Filtering (Keep Relevant Names)
- **Mechanism:** Keep only lines containing relevant identifiers (function names, classes)
- **Latency saving:** 30–60% encoding reduction
- **Trade-off:** Heuristic; accuracy loss if identifier naming is poor

```python
def prune_by_identifiers(file_content, identifiers):
    lines = file_content.split('\n')
    identifier_set = set(identifiers)
    
    pruned = [
        line for line in lines
        if any(ident in line for ident in identifier_set)
    ]
    return '\n'.join(pruned)

# Keep only lines mentioning "auth", "user", "login"
# 8KB → 2KB → 80ms encoding (vs 200ms full)
```

### Pruning Decision

| Pruning Type | Latency Saving | Accuracy Loss | Complexity | Best For |
|--------------|-----------------|----------------|------------|----------|
| Recency (last N lines) | 50–80% | 1–3% | Very low | First pass; simple files |
| Semantic (relevance) | 40–70% | 1–2% | Medium | Large files; complex queries |
| Identifier filtering | 30–60% | 2–5% | Low | Known identifiers (function names) |
| No pruning | 0% | — | — | Small context (<2KB) |

**Action:** Start with recency pruning (simple, fast); add semantic pruning if accuracy drops.

---

## 4. Batching (Throughput → Latency Trade-off)

Batching multiple requests increases throughput but adds latency per-request.

### Batching Mechanics

```python
# Naive batching: wait for N requests before inference
batch_queue = []
batch_size = 4

async def enqueue_request(prompt):
    batch_queue.append(prompt)
    
    if len(batch_queue) == batch_size:
        # Ready to inference
        results = model.generate_batch(batch_queue)
        for i, result in enumerate(results):
            yield results[i]
        batch_queue = []

# Problem: Request N+1 waits 200ms for other requests to fill batch
# Solution: Use timeout + dynamic batching
```

#### Dynamic Batching with Timeout

```python
import asyncio
import time

class DynamicBatcher:
    def __init__(self, model, batch_size=4, timeout_ms=50):
        self.model = model
        self.batch_size = batch_size
        self.timeout_ms = timeout_ms / 1000.0
        self.queue = []
        self.lock = asyncio.Lock()
    
    async def process(self, prompt):
        async with self.lock:
            self.queue.append({"prompt": prompt, "future": asyncio.Future()})
        
        # Start batcher if this is first request
        if len(self.queue) == 1:
            asyncio.create_task(self._batch_process())
        
        # Wait for result (batcher will fulfill future)
        return await self.queue[-1]["future"]
    
    async def _batch_process(self):
        while self.queue:
            async with self.lock:
                # Wait for batch to fill or timeout
                start = time.time()
                while len(self.queue) < self.batch_size and (time.time() - start) < self.timeout_ms:
                    await asyncio.sleep(0.001)  # 1ms check interval
                
                # Process current batch
                current_batch = self.queue[:self.batch_size]
                self.queue = self.queue[self.batch_size:]
            
            # Inference (parallel for batch)
            prompts = [req["prompt"] for req in current_batch]
            results = self.model.generate_batch(prompts)
            
            # Fulfill futures
            for req, result in zip(current_batch, results):
                req["future"].set_result(result)

# Usage
batcher = DynamicBatcher(model, batch_size=4, timeout_ms=50)

async def handle_request(prompt):
    result = await batcher.process(prompt)
    return result

# Request 1: 50ms timeout, batch size 1 → latency 600ms (inference)
# Requests 2–4: added to batch while waiting
# Request 5: waits 50ms for batch to fill, then inference 600ms → total 650ms
```

**Latency cost per request in batch:**
- Request 1: 600ms (inference only)
- Request 2: 600ms + 0ms wait (arrives during inference)
- Request 3: 600ms + 0ms wait
- Request 4: 600ms + 0ms wait
- Request 5: 600ms + 50ms timeout = 650ms

**Decision:** Batch size 4–8 with 50ms timeout is sweet spot. Larger batches → higher per-request latency.

---

## 5. Speculative Decoding (Draft + Verify)

Use a small model to draft tokens; large model verifies and refines. Typical 1.5–2x speedup.

### How It Works

```
[Large Model (Sonnet)]
  ↓ (slow forward pass)
  [draft tokens from small model (Haiku)]
    ↓ (fast; 50 tokens/sec)
    [large model verifies in parallel]
      ↓ (batch verify is fast)
      [keep or reject speculated tokens]
        ↓
        [output; continue or redraft]
```

### Implementation (Conceptual)

```python
# Pseudocode; actual implementation requires custom kernel

def speculative_decode(prompt, draft_model, verify_model, max_tokens=100):
    output = []
    for _ in range(max_tokens):
        # Draft: small model generates 4–8 tokens
        draft = draft_model.generate(prompt + output, max_new_tokens=8)
        
        # Verify: large model checks draft in parallel
        logits = verify_model.forward(prompt + output + draft)
        
        # Accept speculated tokens if probabilities match
        for i, token in enumerate(draft):
            if token in top_k_from_logits(logits[-(len(draft)-i)]):
                output.append(token)
            else:
                # Reject; use verify model's choice
                output.append(top_1_from_logits(logits[-(len(draft)-i)]))
                break  # Restart draft from here
    
    return output

# Expected speedup: 1.5–2x (speculating 4–8 tokens per verify step)
```

**Availability:** Not yet in Claude API. Available in:
- vLLM + medusa heads (requires special model)
- Custom server + two-model setup

**Trade-off:** Quality preserved; setup complexity high.

---

## 6. Streaming Output (Latency Perception)

Streaming does NOT reduce latency—it improves perceived latency.

### Streaming Benefits

- **User sees first token in 150ms** (vs 2000ms all-at-once)
- **UI feels responsive** (output appears incrementally)
- **Total latency still 2000ms**, but UX perception is 8–10x faster

### Implementation (FastAPI + SSE)

```python
from fastapi import FastAPI, StreamingResponse
import anthropic
import json

app = FastAPI()
client = anthropic.Anthropic()

@app.post("/generate")
async def generate_stream(prompt: str):
    async def stream():
        with client.messages.stream(
            model="claude-3-5-sonnet-20241022",
            max_tokens=500,
            messages=[{"role": "user", "content": prompt}],
        ) as stream:
            for text in stream.text_stream:
                # Stream token as JSON event
                yield f"data: {json.dumps({'token': text})}\n\n"
    
    return StreamingResponse(stream(), media_type="text/event-stream")

# Client-side (JavaScript)
# fetch('/generate?prompt=...', {method: 'POST'})
#   .then(res => res.body.pipeThrough(new TextDecoderStream()))
#   .then(readable => {
#       const reader = readable.getReader();
#       while (true) {
#           const {done, value} = await reader.read();
#           if (done) break;
#           updateUI(value);  // Show tokens as they arrive
#       }
#   })
```

### Streaming Decision

| Latency | Streaming? | Perception | UX |
|---------|-----------|------------|-----|
| <500ms | No | Instant | Good (no streaming needed) |
| 500ms–1s | Optional | Fast | Good (maybe streamed) |
| 1–2s | Strongly recommended | Slower; streaming helps | Better with streaming |
| >2s | Mandatory | Very slow; must stream | Unusable without streaming |

**Action:** Always stream for >1s SLAs. <500ms features don't benefit.

---

## 7. Model Sharding (Distributed Inference)

Split model across multiple GPUs to reduce per-GPU memory; enables larger model on same hardware.

### Tensor Parallelism (TensorRT-LLM)

```python
from vllm import LLM

# Single GPU: runs out of memory on 70B model
# With tensor parallelism: split across 2 GPUs
llm = LLM(
    model="Qwen/Qwen2.5-Coder-72B",
    tensor_parallel_size=2,  # Split across 2 GPUs
    dtype="float16",
)

# Latency trade-off: ~2x throughput, ~1.2x latency increase
# P95: 800ms (single GPU) → 950ms (2-GPU)
# But can handle 2x concurrent requests without queue
```

**When to use:** High QPS (>10 req/s) where throughput is critical, and 1.2x latency increase is acceptable.

---

## Optimization Checklist

### Before Optimizing

- [ ] Measure baseline latency (P50, P95, P99)
- [ ] Identify bottleneck (encode? infer? post-proc?)
- [ ] Know SLA (what latency is acceptable?)

### Optimizations (Ranked by Impact)

1. **Quantization (int8):** 2.5–3x latency reduction; <1% quality loss → Do first
2. **Switch smaller model:** 1.5–4x latency reduction; quality drop ~1–5% → Do if quantization insufficient
3. **Prompt caching:** 1.5–2x encoding reduction; no quality loss → Do if context is repeated
4. **Context pruning:** 1.2–2x encoding reduction; 1–3% quality loss → Do if context is large
5. **Batching:** 1.1–1.2x throughput improvement; per-request latency +50ms → Do if QPS >5 req/s
6. **Speculative decoding:** 1.5–2x inference reduction; complex setup → Do only for frontier models and high QPS
7. **Streaming:** 0x latency reduction, UX improvement → Do for >1s SLAs
8. **Model sharding:** Throughput improvement; latency penalty → Do if QPS is very high

### Per-SLA Optimization Path

**<500ms SLA:**
1. Quantization (int8) on small model (<10B)
2. If insufficient: switch to even smaller model or add speculation
3. Caching if context is repeated
4. Pruning if context is large

**<1s SLA:**
1. Quantization (int8) on mid-size model (7–32B)
2. Caching + pruning
3. Streaming (optional; improves perception)

**<2s SLA:**
1. Quantization (int8) on larger model (32–70B)
2. Streaming (recommended)
3. Batching if QPS high

**<5s SLA:**
1. Streaming (mandatory)
2. Quantization (if fp16 too slow)
3. Batching + caching for high QPS

---

## Real-World Example: IDE Code Review

**SLA:** Code review feedback in <2s P95.

**Baseline (Qwen 72B fp16):**
- Encoding: 400ms
- Inference: 1500ms
- Post-processing: 100ms
- Total: 2000ms ✗ (P95 is ~2500ms; fails SLA)

**Optimization 1: Quantization to int8**
- Encoding: 400ms (same)
- Inference: 500ms (3x faster)
- Post-processing: 100ms
- Total: 1000ms ✓ (P95 ~1200ms; passes SLA)

**Cost:** <1% quality loss on code tasks; acceptable.

**Result:** Deployed with int8; monitoring shows P95 = 1100ms consistently. Success!

---

## References

- **GPTQ:** Frantar et al., "GPTQ: Accurate Post-Training Quantization for Generative Pre-Trained Transformers" (2023)
- **AWQ:** Lin et al., "AWQ: Activation-aware Weight Quantization for LLM Compression and Acceleration" (2023)
- **vLLM:** Kwon et al., "vLLM: Easy, Fast, and Cheap LLM Serving with PagedAttention" (2023)
- **Speculative Decoding:** Leviathan et al., "Fast Inference from Transformers via Speculative Decoding" (2023)
- **KV Cache:** Sheng et al., "Inference with Reference: Lossless Acceleration of Large Language Models" (2023)
