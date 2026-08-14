# Quantization Format Comparison: GGUF, AWQ, GPTQ, FP8, INT8

Deep technical comparison for operators choosing between formats.

## Format Internals & Mechanics

### GGUF (GGML Universal Format)

**Origin:** Developed by ggerganov for llama.cpp; now de-facto standard for portable quantized models.

**Quantization levels:**

1. **Q2_K (2-bit Quantization)**
   - Structure: 256-weight blocks, 2 bits per weight + K-quant scale & min values.
   - Formula: `w_quantized = round((w - min) / (max - min) * 3) # 4 levels`
   - Error: ~6% vs FP16 (high variance across different layers).
   - When: Extreme memory constraint (edge devices, CPU).

2. **Q3_K (3-bit)**
   - Structure: 256-weight blocks, 3 bits per weight (8 levels).
   - Error: ~4% vs FP16.
   - Trade: 86% size reduction vs ~3–4% HumanEval loss.

3. **Q4_0 (4-bit, symmetric)**
   - Structure: 32-weight blocks, symmetric quantization around zero.
   - Formula: `w_q = clamp(round(w / scale), -8, 7) # 4-bit signed int`
   - Error: ~2–3% vs FP16.
   - **Most common; good quality/speed/size balance.**

4. **Q4_K (4-bit, improved)**
   - Structure: 256-weight blocks, separate low-bit (5-bit) and high-bit (3-bit) parts.
   - Key insight: Use more bits for "important" weights; fewer for rest.
   - Error: ~2% vs FP16.
   - Trade: Slightly better quality than Q4_0; negligible speed difference.

5. **Q5_K (5-bit)**
   - Structure: 256-weight blocks, mixed 5–6 bits.
   - Error: <1% vs FP16.
   - Trade: 70% size reduction vs near-lossless quality.

6. **Q6_K (6-bit)**
   - Error: <0.5% vs FP16.
   - VRAM savings: ~50% (same as INT8).
   - Trade: Rarely justified; use FP8 or INT8 instead (faster hardware support).

7. **Q8_0 (8-bit, symmetric)**
   - Error: <0.3% vs FP16.
   - VRAM savings: ~50%.
   - Use case: Archive format; rarely used for inference (INT8 faster).

**Quantization process:**
1. Load FP16 model weights.
2. For each quantization block (32–256 weights):
   - Calculate `min, max` of weights.
   - Map to fixed bit range (4, 5, 8 bits).
   - Encode outliers via K-quant (higher-bit side values).
3. Serialize to single `.gguf` binary file.
4. **Result:** Single file, self-contained, no dependencies.

**Speed characteristics (RTX 4090, Qwen-Coder 32B):**
- FP16: 60 tok/s
- INT8: 70 tok/s (1.16× speedup; minimal; dequant overhead)
- Q5_K: 110 tok/s (1.83×)
- Q4_K: 140 tok/s (2.33×) **← Recommended for GPU inference**
- Q3_K: 180 tok/s (3×) but −4% accuracy
- Q2_K: 220 tok/s (3.66×) but −6% accuracy

**Speed on CPU (llama.cpp, macOS M3 Pro):**
- Q4_K: 8–12 tok/s (single-threaded).
- Parallelism: Scales to ~30 tok/s with 8-thread CPU workload (not linear due to memory bandwidth).

**Portability:** ✅ Single `.gguf` file; load in Ollama, LM Studio, llama.cpp, ExLlama, GPT4All.

---

### AWQ (Activation-aware Weight Quantization)

**Origin:** Paper by Lin et al. (2023); implemented in AutoGPTQ library.

**Key idea:** Weights matter less than activations. Find "important" weight channels (high activation variance) and preserve them; quantize others aggressively.

**Algorithm:**
1. Forward pass: Compute activation statistics (per-channel variance / max).
2. Identify top-K channels with highest activation magnitudes.
3. Preserve top-K channels in higher precision (or FP16); quantize rest to INT4.
4. Per-output-channel quantization scale (not global; highly adaptive).

**Accuracy:** INT4 only; claims near-FP16 quality (−1% vs −2–3% for uniform INT4).

**Why better than GPTQ:**
- GPTQ: Minimizes Hessian approximation error (layer-wise).
- AWQ: Minimizes end-to-end activation error (task-aware).
- AWQ more stable under different prompt/input distributions.

**Format:** `.safetensors` file or `.pt` (PyTorch checkpoint). Requires AutoGPTQ library to load/inference.

**Speed (vLLM on RTX 4090, Qwen-Coder 32B):**
- FP16: 60 tok/s.
- AWQ INT4: 120 t/s (2× speedup). **Slower than GGUF Q4_K on CPU; faster on GPU than GPTQ.**

**Inference framework dependencies:**
- ✅ vLLM (native kernel; highest priority in dev).
- ✅ AutoGPTQ (inference + quantization).
- ⚠️ Hugging Face Transformers (via AutoGPTQ; experimental).
- ❌ Ollama (no official support; community plugin exists).
- ❌ llama.cpp (no native AWQ; conversion to GGUF possible but lossy).

**Quantization cost:** ~1–2 hours CPU time for 32B model (matrix operations on CPU).

**Portability:** ❌ Not portable as single file; requires AutoGPTQ ecosystem. Heavier dependencies than GGUF.

---

### GPTQ (Quantization using Approximate Second-order Information)

**Origin:** Paper by Frantar et al. (2022); most widely used INT4 method historically.

**Algorithm:**
1. Compute Hessian (second-order) approximation of loss per layer.
2. Find weights that minimally increase loss (greedy per-channel).
3. Quantize; recompute Hessian; adapt scales.
4. Group quantization (typically 128–256 weights per group) for efficiency.

**Accuracy:** INT4 (primary); INT3/INT8 variants exist but less common.

**Why GPTQ vs AWQ:**
- GPTQ: Well-studied; many existing quantized models (TheBloke library).
- AWQ: Better per-output-channel precision; newer codebases prefer it.
- **Empirical:** AWQ ~1–2% better quality at INT4 vs GPTQ INT4.

**Format:** `.safetensors` (preferred) or `.pt`. Requires AutoGPTQ or ExLlama.

**Speed (RTX 4090, Qwen-Coder 32B):**
- GPTQ INT4: 90 tok/s (1.5× vs FP16). Slower than AWQ because:
  - Kernel optimizations less mature (AWQ prioritized by vLLM team).
  - Group dequantization CPU-intensive on older frameworks.

**Inference frameworks:**
- ✅ AutoGPTQ (native).
- ✅ ExLlama / ExLlamaV2 (highly optimized for GPTQ; 2–3× faster than AutoGPTQ).
- ✅ vLLM (via AutoGPTQ backend; not optimized; slower than native AWQ).
- ⚠️ Hugging Face Transformers (experimental).
- ❌ Ollama, llama.cpp (no native support).

**Quantization cost:** ~1–2 hours CPU time for 32B model (similar to AWQ).

**Portability:** ❌ Requires AutoGPTQ/ExLlama ecosystem.

---

### FP8 (8-bit Floating Point)

**Formats:**
1. **E4M3 (4 exponent, 3 mantissa):** Better precision; range ~1e-7 to 1e2.
2. **E5M2 (5 exponent, 2 mantissa):** Better range; less precision.

**How quantization works:**
- Map FP32 activations/weights → FP8 via scaling.
- `fp8_val = scale * fp32_val # With clipping to prevent overflow`
- Inference: Dequant to higher precision before compute (or use FP8 Tensor Cores if available).

**Accuracy (HumanEval):**
- E4M3: −0.5–1% (near-lossless).
- E5M2: −1–2% (more variance in extreme ranges).

**Speed characteristics:**
- **H100 / L40S (Ada Tensor Cores):** 2× speedup over FP16 (native FP8 support).
- **RTX 6000 Ada:** 1.7–2× speedup.
- **RTX 4090 (Hopper, no FP8 core):** No speedup; emulation via casting ~0.9–1.3× (slower). **Don't use.**

**VRAM savings:** 50% (same as INT8). FP32 → FP8: 4× reduction; quantized model: 2× vs FP16.

**Format:** Native PyTorch/Hugging Face format (no special quantized file). Loaded as `.safetensors` with dtype=fp8 metadata.

**Inference frameworks:**
- ✅ vLLM (0.4.1+, FP8 kernel for Ada).
- ✅ bitsandbytes (FP8 quantization + inference).
- ✅ Hugging Face Transformers (experimental).
- ❌ Ollama, llama.cpp (no FP8 support).

**Quantization cost:** ~30 minutes (much faster than INT4; just FP32→FP8 casting).

**Portability:** Requires explicit FP8 Tensor Core hardware or emulation fallback.

---

### INT8 (8-bit Integer)

**Formats:**
1. **Symmetric:** Weights quantized around zero; same scale for positive/negative.
2. **Asymmetric:** Separate scales for positive/negative; better for skewed distributions.

**Quantization:**
```python
# Symmetric INT8
w_min, w_max = min(w), max(w)
scale = max(abs(w_min), abs(w_max)) / 127
w_int8 = round(w / scale).astype(int8)

# Asymmetric INT8
scale = (w_max - w_min) / 255
zero_point = round(-w_min / scale)
w_int8 = round((w - w_min) / scale).astype(uint8)
```

**Accuracy:** <0.5–1% HumanEval loss (near-lossless).

**Speed:**
- Hugging Face `load_in_8bit=True`: 1.1–1.3× (slow; minimal optimization).
- bitsandbytes INT8: 1.2–1.5× (better kernels).
- vLLM INT8 (experimental): 1.1–1.4×.
- **Summary:** Speed gains marginal on consumer GPUs; no specialized hardware.

**VRAM savings:** 50% (same as FP8).

**Inference frameworks:**
- ✅ bitsandbytes (optimized).
- ✅ Hugging Face Transformers (naive `load_in_8bit`; slow).
- ✅ GGUF Q8_0 variant (portable).
- ⚠️ vLLM (experimental).
- ❌ Ollama, llama.cpp (no native INT8; use GGUF Q8_0 instead).

**Quantization cost:** ~10 minutes (simple casting).

**Portability:** Format-agnostic; can be INT8 Tensors in `.safetensors` or GGUF Q8_0.

---

## Comparison Matrix: All Dimensions

| Dimension | GGUF Q4_K | AWQ INT4 | GPTQ INT4 | FP8 (Ada) | INT8 |
|-----------|-----------|----------|-----------|-----------|------|
| **Accuracy (HumanEval)** | 78–79% | 78–79% | 77–78% | 79–80% (E4M3) | 79–80% |
| **Speed (RTX 4090)** | 140 t/s | 120 t/s | 90 t/s | 60 t/s (no gain) | 70 t/s |
| **Speed (H100)** | N/A | ~120 t/s | ~90 t/s | 120 t/s | 70 t/s |
| **VRAM (32B model)** | 18 GB | 16 GB | 16 GB | 32 GB | 32 GB |
| **Quantization time** | N/A (pre-made) | 1–2 hrs | 1–2 hrs | 30 min | 10 min |
| **Portability** | ✅ Single file | ❌ Framework-tied | ❌ Framework-tied | ❌ Hardware-tied | ⚠️ Format-agnostic |
| **Ollama support** | ✅ Native | ⚠️ Plugin | ❌ | ❌ | ❌ |
| **llama.cpp support** | ✅ Native | ❌ | ❌ | ❌ | ✅ Q8_0 |
| **vLLM support** | ⚠️ Experimental | ✅ Native | ✅ (AutoGPTQ) | ✅ (0.4.1+, Ada) | ⚠️ Experimental |
| **Batch inference** | ✅ Good | ✅ Good (via vLLM) | ⚠️ Slower | ✅ Good (Ada) | ⚠️ Marginal gain |
| **Production readiness** | ✅ Mature | ✅ Mature | ⚠️ Older standard | ❌ Ada-only | ⚠️ Niche use |
| **Best for** | CPU/consumer GPU | GPU batch (vLLM) | Legacy systems | H100 production | Research fallback |

---

## Decision Flowchart

```
Start: Choose quantization format for Qwen-Coder 32B

1. GPU Type?
   ├─ H100 / L40S
   │  └─ → FP8 E4M3 (native Tensor Core; 2× speed, −1% quality)
   │
   ├─ RTX 4090 / RTX 6000
   │  ├─ vLLM batch inference?
   │  │  └─ → AWQ INT4 (vLLM native; 2× speed, −2% quality)
   │  │
   │  ├─ CPU inference / Ollama?
   │  │  └─ → GGUF Q4_K (portability; 2.3× speed, −2% quality)
   │  │
   │  └─ Single-GPU serving (Python)?
   │     └─ → AWQ INT4 (best speed/quality)
   │
   ├─ RTX 3060 / RTX 3070 (12–16 GB)
   │  └─ → GGUF Q4_K (only fits; 2.3× speed, −2% quality)
   │
   ├─ Apple M-series Mac
   │  └─ → GGUF Q4_K (Ollama; Metal acceleration)
   │
   └─ A100 (already have VRAM)
      ├─ High quality required?
      │  └─ → No quantization (FP16 baseline)
      │
      └─ Cost-optimize?
         └─ → FP8 (50% VRAM; <−1% quality; native support)

2. If batch processing (overnight runs)?
   └─ Speed doesn't matter; use GGUF Q3_K (most aggressive; saves VRAM).

3. If research / reproducibility?
   └─ Use GPTQ INT4 (well-documented; many existing models).
```

---

## Known Interactions & Edge Cases

### Quantization + Context Length

**Issue:** Longer context (8K+ tokens) increases KV-cache memory. Quantization reduces weight VRAM but KV-cache unchanged.

**Example (Qwen-Coder 32B, RTX 4090, 8K context, batch size 4):**
- FP16: 65 GB (weights) + 24 GB (KV) = 89 GB total. Doesn't fit.
- INT4 GGUF: 18 GB (weights) + 24 GB (KV) = 42 GB total. Fits with headroom.

**Takeaway:** Quantization benefit shrinks as context grows. Use INT8/FP8 if long context (>4K tokens).

### Quantization + LoRA Fine-tuning

**Issue:** Standard quantized models (GGUF, AWQ, GPTQ) not trainable without full precision.

**Options:**
1. **QLoRA:** Load model in INT4, add FP16 LoRA adapters.
   - Framework: bitsandbytes + Hugging Face `transformers`.
   - Trade: Training slow (INT4 dequant overhead); VRAM savings worth it.

2. **Full-precision fine-tune:** Unquantize to FP16; fine-tune; re-quantize.
   - Trade: Expensive; loses quantization benefit during training.

**Recommendation:** Use QLoRA + INT4 for fine-tuning; standard quantized for inference-only.

### Quantization + Speculative Decoding

**Issue:** Speculative decoding uses a smaller draft model to predict next K tokens; main model verifies.

**Quantization strategy:**
- Draft model: Aggressive quantization (Q3_K, INT4) for speed.
- Main model: Conservative quantization (Q4_K, INT8, FP8) for accuracy.

**Example (vLLM):**
```python
# Draft model: Mistral 7B GGUF Q3_K
# Main model: Qwen-Coder 32B AWQ INT4
# Result: 1.5–2× total speedup vs main model alone
```

### Quantization + Batching

**Issue:** Batch inference increases KV-cache linearly; may offset quantization VRAM savings.

**Example (Qwen-Coder 32B, RTX 4090):**
- Batch size 1: INT4 GGUF = 18 GB (weights) + 3 GB (KV) = 21 GB. Fits.
- Batch size 8: INT4 GGUF = 18 GB + 24 GB (KV) = 42 GB. Still fits.
- Batch size 32: INT4 GGUF = 18 GB + 96 GB (KV) = 114 GB. OOM.

**Takeaway:** Monitor KV-cache; quantization VRAM savings per-model, not per-batch.

---

## Real-World Selection Examples

### Example 1: IDE Autocomplete on Mac
- Model: Mistral 7B or Qwen2.5-Coder 7B
- Format: **GGUF Q4_K**
- Speed: 5–6 sec latency (acceptable for autocomplete).
- Framework: Ollama.
- Quality loss: −2%, acceptable for suggestions user reviews.

### Example 2: Production Code Generation (RTX 4090)
- Model: Qwen-Coder 32B or DeepSeek-Coder 33B
- Format: **AWQ INT4**
- Speed: 2–3 sec (vLLM batch, paged attention).
- Quality loss: −2%, acceptable for automated testing.
- Framework: vLLM (OpenAI-compatible API).

### Example 3: One-off Refactoring (Batch)
- Model: Any 32B specialist
- Format: **GGUF Q3_K** (aggressive, saves VRAM, humans review).
- Speed: Doesn't matter (overnight run).
- Quality loss: −4%, acceptable for review-before-commit.
- Framework: llama.cpp (most portable).

### Example 4: High-bar Code (Security Audit)
- Model: DeepSeek-Coder 33B
- Format: **FP8 E4M3** (if H100); **INT8** (if A100); **AWQ INT4 or GGUF Q5_K** (consumer GPU, accept −1–2%).
- Speed: Secondary; accuracy primary.
- Framework: vLLM (H100, L40S) or llama.cpp (portable).

---

## References

1. [llama.cpp GGUF Format](https://github.com/ggerganov/llama.cpp/blob/master/gguf-py/README.md)
2. [AWQ Paper: Activation-Aware Weight Quantization](https://arxiv.org/abs/2306.00978)
3. [GPTQ Paper: Accurate Post-Training Quantization](https://arxiv.org/abs/2210.17323)
4. [FP8 Training & Inference (vLLM)](https://docs.vllm.ai/en/latest/quantization)
5. [bitsandbytes INT8 & FP8](https://github.com/TimDettmers/bitsandbytes)
6. [AutoGPTQ](https://github.com/AutoGPTQ/AutoGPTQ)
