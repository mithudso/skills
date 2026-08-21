---
id: local-llm-inference-serving
family: ai-llm-model-layer
name: Local LLM Inference Serving for Coding Agents
author: claude-code
created: 2026-07-10
updated: 2026-07-10
version: 1.0
model: claude-opus-4-8
effort: high
tags: ["LLM", "inference", "serving", "local", "coding-agents", "quantization", "optimization"]
description: |
  Setup, tuning, and troubleshooting for local LLM inference serving specifically for AI coding agents. Covers Ollama, vLLM, LocalAI, LM Studio. Focus: quantization (GGUF/AWQ/GPTQ), batching, KV-cache, latency/throughput tradeoffs, multi-GPU, context management, temperature/top-p tuning, and monitoring (tokens/s, TTFT).
whenToUse: |
  TRIGGER: "run qwen-coder locally", "set up local inference", "optimize llm serving", "configure ollama", "vllm setup", "quantization format choice", "batching strategy", "KV-cache tuning", "multi-GPU inference", "TTFT latency", "tokens/sec throughput", "temperature tuning for code", "monitoring local inference".
  SKIP: ai-llm-model-layer (for LLM theory/training), ai-mcp-sdk-prompting (for prompt engineering), ai-rag-retrieval (for RAG patterns).
---

# Local LLM Inference Serving for Coding Agents

Operator-level setup, tuning, and troubleshooting for local LLM inference specifically for AI coding agents.

---

## What This Covers

**In Scope** (operator-level):
- Framework selection (Ollama vs vLLM vs LocalAI vs LM Studio) for coding agents
- Quantization format tradeoffs (GGUF, AWQ, GPTQ, INT8/INT4)
- Setup and configuration for each framework
- Batching strategies (static, dynamic, token-level)
- KV-cache optimization (paged, prefix caching, context window sizing)
- Attention optimization (Flash Attention, GQA, MQA)
- Token generation optimization (speculative decoding, constrained generation)
- Multi-GPU setup (tensor parallelism, pipeline parallelism, sequence parallelism)
- Temperature/top-p tuning specific to code generation
- Monitoring metrics (tokens/sec, TTFT, VRAM, GPU utilization)
- Bottleneck identification and remediation
- Cost/performance tradeoffs (local vs cloud)

**Out of Scope** (deep LLM theory):
- Transformer internals (see ai-llm-model-layer for theory)
- Pretraining, fine-tuning, RLHF (see ai-llm-model-layer)
- General LLM architecture (see ai-llm-model-layer)
- Prompt engineering, RAG retrieval (see ai-mcp-sdk-prompting, ai-rag-retrieval)

---

## Framework Comparison (2026)

| Feature | Ollama | vLLM | LocalAI | LM Studio |
|---------|--------|------|---------|-----------|
| **Setup Complexity** | Very Low | Medium | Medium | Very Low (GUI) |
| **Best For** | Quick local dev | High-throughput serving | CPU + privacy | Interactive GUI testing |
| **Batching** | Static, auto-tuned | Dynamic (paged) | Static | Limited |
| **KV-Cache** | Optimized | Paged (SOTA) | Standard | Standard |
| **Multi-GPU** | Model sharding | Tensor+Pipeline | Limited | Limited |
| **Quantization** | GGUF native | GPTQ/AWQ/fp8 | GGUF/GPTQ | GGUF |
| **Context Window** | Scales well | Full support | Scales well | Good |
| **Production Ready** | Yes | Yes | Partial | No (GUI-only) |
| **Metrics/Monitoring** | Basic logs | Prometheus | Minimal | None |

**Quick Pick**:
- **Local dev + quick test**: Ollama
- **High-throughput production**: vLLM
- **Privacy + CPU fallback**: LocalAI
- **Interactive exploration**: LM Studio

---

## Quantization Formats & Model Selection

### GGUF (Industry Standard for Local)
**Best for**: Ollama, llama.cpp, consumer GPUs, edge deployment

- **Precision levels**: F32, F16, Q8_0, Q6_K, Q5_K, Q4_K, Q3_K, Q2_K
- **Code quality**: Good at Q4_K_M and above; acceptable Q5_K
- **Performance**: ~10-15 tokens/sec on RTX 3090 for 7B; ~5-8 for 13B
- **VRAM**: ~3.5GB (7B Q4), ~7GB (13B Q4)
- **Trade-off**: Maximum compatibility; mature tooling; slightly slower than AWQ on GPU

**Coding model availability**:
- Qwen-Coder-Next: Widely available on Ollama registry
- DeepSeek-Coder: Good availability via community conversions
- CodeLlama, Magicoder: Community GGUF conversions present

**When to use**: Development, edge deployment, CPU fallback needed, maximum model selection.

### AWQ (Better Quality-Speed Tradeoff)
**Best for**: vLLM high-performance serving, GPU-only inference

- **Quality**: ~0.5% loss vs FP16 at INT4 (best for coding)
- **Performance**: ~20% faster than GPTQ; slightly slower than GGUF on CPU
- **VRAM**: ~3-4GB (7B INT4), ~6-7GB (13B INT4)
- **Integration**: vLLM native; requires GPU (no CPU support)
- **Maturity**: Growing rapidly; Qwen-Coder has official AWQ quantizations

**When to use**: Production GPU inference, quality-sensitive (code review/reasoning), vLLM serving.

### GPTQ (Mature & Widely Available)
**Best for**: Balance of speed and quality; existing HuggingFace ecosystem

- **Quality**: Good for chat (1-2% loss); some loss for code (2-3%)
- **Performance**: ~15-20 tokens/sec on RTX 3090 for 7B
- **VRAM**: ~3-4GB (7B), ~6-7GB (13B)
- **Integration**: AutoGPTQ library, vLLM, ExLlama
- **Availability**: Abundant on HuggingFace

**When to use**: Mixed chat+code workloads, maximum format compatibility, balance needed.

### INT8 vs INT4 vs Mixed Precision
- **INT8**: Better quality (1-2% loss), ~15% slower inference; use only if quality critical
- **INT4**: Best speed/space, acceptable quality (3-5% loss); recommended for most cases
- **Mixed** (e.g., 8-bit activations, 4-bit weights): Emerging; better quality retention; not yet in most tools

**Recommendation for coding agents**: INT4 AWQ or GPTQ; INT8 only if hallucination unacceptable.

---

## Ollama Setup & Configuration

### Install & Basic Setup
```bash
# macOS/Linux
brew install ollama

# or Docker
docker run -d -v ollama:/root/.ollama -p 11434:11434 ollama/ollama

# Pull a model (auto-selects quantization)
ollama pull qwen:latest          # Auto 7B Q4
ollama pull qwen:32b             # 32B model
ollama pull qwen:32b-code-int4   # Code variant, INT4
```

### Run a Model
```bash
# Interactive chat
ollama run qwen:latest

# API server (listens on localhost:11434)
ollama serve

# In another terminal, test the API
curl http://localhost:11434/api/generate -d '{
  "model": "qwen:latest",
  "prompt": "Write a Python function to reverse a list",
  "stream": false,
  "temperature": 0.3,
  "top_p": 0.9
}'
```

### Batching & Performance Tuning
```bash
# Environment variables (set before running)
export OLLAMA_NUM_PARALLEL=4          # Concurrent requests (default 1)
export OLLAMA_NUM_GPU_LAYERS=35       # GPU layers to offload (auto if not set)
export OLLAMA_MAIN_GPU=0              # Primary GPU index (for multi-GPU)

ollama serve
```

**Tuning Strategy**:
- **Single GPU**: Leave defaults; Ollama auto-tunes
- **Multi-GPU**: Set `OLLAMA_MAIN_GPU=0` and `OLLAMA_NUM_PARALLEL=4` to start
- **VRAM-constrained**: Reduce `OLLAMA_NUM_GPU_LAYERS` to CPU-offload more weights

### Context Window & Memory Management
```bash
# Context window is typically auto-detected from model
# For explicit control via API:
curl http://localhost:11434/api/generate -d '{
  "model": "qwen:latest",
  "prompt": "...",
  "context_length": 8192    # Limit context to 8K tokens
}'
```

### Monitoring Ollama
- Check resource usage: `nvidia-smi` for GPU, `top` for CPU
- Logs: stderr output shows token/s rate; parse for TTFT from timing
- No native Prometheus export yet; custom monitoring via API polling recommended

---

## vLLM Setup & Configuration

### Install
```bash
pip install vllm

# If using AWQ quantization:
pip install autoawq

# For speculative decoding:
pip install outlines  # For grammar-based generation
```

### Start vLLM Server
```bash
python -m vllm.entrypoints.openai.api_server \
  --model=Qwen/Qwen1.5-32B-Chat \
  --tensor-parallel-size=1 \
  --gpu-memory-utilization=0.9 \
  --max-model-len=8192 \
  --dtype auto \
  --quantization awq \
  --enable-prefix-caching
```

**Key flags**:
- `--tensor-parallel-size=N`: Tensor parallelism across N GPUs
- `--gpu-memory-utilization=0.9`: Use 90% GPU VRAM (adjust for safety)
- `--max-model-len=`: Limit context window (e.g., 8192 for latency, 32768 for long context)
- `--quantization awq|gptq|fp8|bitsandbytes`: Choose quantization
- `--enable-prefix-caching`: Cache common prompt prefixes (code templates, system messages)

### Batching Configuration
```bash
python -m vllm.entrypoints.openai.api_server \
  --model=Qwen/Qwen1.5-32B-Chat \
  --max-num-batched-tokens=8192 \
  --max-num-seqs=16 \
  --max-seq-len-to-capture=2048
```

**Tuning**:
- `--max-num-batched-tokens`: Total tokens per batch (8192-16384 typical)
- `--max-num-seqs`: Max requests per batch (4-16 typical)
- `--max-seq-len-to-capture`: Precompute attention for short sequences (<2048 tokens)

### Using vLLM Programmatically
```python
from vllm import LLM, SamplingParams

llm = LLM(
    model="Qwen/Qwen1.5-32B-Chat",
    quantization="awq",
    tensor_parallel_size=1,
    gpu_memory_utilization=0.9,
    enable_prefix_caching=True,
    max_model_len=8192
)

# Generate
sampling_params = SamplingParams(
    temperature=0.3,
    top_p=0.9,
    max_tokens=1024,
    repetition_penalty=1.05  # Reduce repetition in code
)

prompts = [
    "Write a Python function to reverse a list",
    "Implement a binary search"
]

outputs = llm.generate(prompts, sampling_params)
for output in outputs:
    print(output.outputs[0].text)
```

### Monitoring vLLM
```bash
# Prometheus metrics (export to port 8000)
python -m vllm.entrypoints.openai.api_server \
  --model=Qwen/Qwen1.5-32B-Chat \
  --metric-log-dir=/tmp/vllm_metrics

# Query metrics
curl http://localhost:8000/metrics
```

**Key metrics to track**:
- `vllm:num_requests_running`: Concurrent requests
- `vllm:num_tokens_generated_total`: Total tokens generated
- `vllm:time_to_first_token_seconds`: TTFT histogram
- `vllm:request_duration_seconds`: End-to-end latency

---

## LocalAI Setup

### Install
```bash
# Binary download
wget https://github.com/mudler/LocalAI/releases/download/v2.x/local-ai-linux-amd64

# or Docker
docker run -p 8080:8080 localai/localai:latest
```

### Download & Run a Model
```bash
# LocalAI uses GGUF by default
docker run -v /models:/models -p 8080:8080 \
  localai/localai:latest \
  /app/local-ai --models-path /models start

# In another terminal, download a model
curl http://localhost:8080/v1/models -d '{
  "model": "qwen:7b-q4",
  "backend": "llama-cpp"
}'
```

### API Usage (OpenAI-compatible)
```bash
curl http://localhost:8080/v1/chat/completions -d '{
  "model": "qwen:7b-q4",
  "messages": [{"role": "user", "content": "Write Python code to reverse a list"}],
  "temperature": 0.3,
  "top_p": 0.9
}'
```

**LocalAI Strength**: CPU fallback (slow but works) + privacy-first design.
**LocalAI Limitation**: Limited batching, no native multi-GPU, minimal monitoring.

---

## LM Studio (GUI-Based)

### Setup
1. Download from https://lmstudio.ai
2. Launch GUI; built-in model browser
3. Select and download model (auto-quant selection)
4. Click "Load" to serve

### Local Server
- Default: `http://localhost:1234/v1`
- API-compatible with OpenAI client

**Use Case**: Interactive testing, debugging; not recommended for production automation.

---

## KV-Cache Optimization

### Paged KV-Cache (vLLM)
**Enable in vLLM** (automatic with recent versions):

```bash
python -m vllm.entrypoints.openai.api_server \
  --model=Qwen/Qwen1.5-32B-Chat \
  --enable-prefix-caching \
  --block-size=16  # Token block size for paging
```

**Effect**: ~2x batch size on same GPU VRAM; small latency overhead (<5%).

### Prefix Caching for Coding Agents
Store KV cache for reusable system prompts:

```python
# vLLM with prefix caching
prompts = [
    "Analyze this code:\n" + code_snippet_1,
    "Analyze this code:\n" + code_snippet_2,  # Shares "Analyze this code:" KV cache
]
```

**Benefit**: 2-5x speedup when prompt prefix is shared across requests (common in batch code review).

### Context Window Sizing for Code
```python
# Trade-off: more context = better understanding, but higher latency
# Practical guidelines:

# Short snippets (code review)
max_context = 8192  # 8K tokens typical

# File understanding
max_context = 16384  # 16K for small-medium files

# Deep analysis (full module)
max_context = 32768  # Full window only if latency acceptable (2-3x slower)
```

---

## Attention Optimization

### Flash Attention 2 / Flash Decoding
**Status**: Mature in vLLM, beta in Ollama; automatic in most models (2026).

- **Speedup**: 2-3x faster attention computation
- **VRAM**: Lower KV-cache memory
- **No manual configuration needed** — enabled by default on supported hardware (H100, L40S, RTX 40xx, RTX 3090+)

### Multi-Query/Grouped-Query Attention (MQA/GQA)
Modern coding models (Qwen-Coder) use GQA by default:

- **Benefit**: 4-8x smaller KV-cache; enables long context on modest GPUs
- **No tuning needed** — model architecture choice; just use modern models

---

## Multi-GPU Setup

### Tensor Parallelism (vLLM)
**Best for**: Large models (>20B params) on dual/quad GPUs.

```bash
python -m vllm.entrypoints.openai.api_server \
  --model=Qwen/Qwen1.5-32B-Chat \
  --tensor-parallel-size=2  # Distribute across 2 GPUs
```

**Performance**:
- Dual RTX 3090 (48GB total): Qwen-Coder-Next 32B with ~15-20 tokens/sec
- Dual RTX 4090 (48GB total): Same model with ~25-30 tokens/sec

**Communication overhead**: ~15-20% from AllReduce operations.

### Pipeline Parallelism (vLLM Advanced)
**For**: Very large models (70B+) or bandwidth-constrained.

```bash
python -m vllm.entrypoints.openai.api_server \
  --model=Qwen/Qwen1.5-70B-Chat \
  --pipeline-parallel-size=2  # Different layers on different GPUs
```

**Caveat**: Pipeline parallelism has bubble time; typically slower than tensor parallelism unless model >70B.

### Ollama Multi-GPU
**Automatic**: Set primary GPU via environment.

```bash
export OLLAMA_MAIN_GPU=0
ollama serve
```

Ollama handles model sharding automatically; no explicit config needed. Less flexible than vLLM for advanced users.

---

## Temperature & Top-P Tuning for Code Generation

### Temperature Guidelines
- **0.0-0.3**: Deterministic, high-quality code (best for code generation)
- **0.3-0.7**: Balanced creativity + correctness (good for suggestions)
- **0.7-1.0**: Creative but risky for code (likely hallucinations)

**Coding Agent Recommendations**:
- **Primary generation**: 0.2-0.4 (reliable, low hallucination)
- **Alternative suggestions**: 0.6 (some creativity without breaking)
- **Interactive exploration**: 0.5 (balance)

### Top-P (Nucleus Sampling)
- **0.8-0.95**: Typical for code (filters very unlikely tokens)
- **0.95+**: Allows more exploration; slightly increased hallucination

**Recommended combination for code**: `temperature=0.3, top_p=0.9`

### Min-P (Emerging Alternative, 2026)
```python
# In vLLM (if supported)
sampling_params = SamplingParams(
    temperature=0.3,
    min_p=0.05,  # Filter tokens below 5% of top token
)
```

**Benefit**: Better coherence than top-p alone for code.

---

## Monitoring & Metrics

### Key Metrics to Track

**Latency**:
- **TTFT (Time-to-First-Token)**: <100ms ideal for interactive code completion
- **TPS (Tokens-Per-Second)**: >5 for interactive, >20 for batch
- **p50/p95/p99 latencies**: Track percentiles for SLA compliance

**Resource Usage**:
- **GPU Utilization**: Monitor via `nvidia-smi`; target 80-95% for sustained loads
- **GPU Memory**: Track peak vs sustained; alert on OOM
- **GPU Power Draw**: Correlates with efficiency

**Quality Metrics**:
- **Perplexity**: Benchmark on code (HumanEval subset)
- **Hallucination Rate**: Monitor incorrect code generation
- **Context Window Efficiency**: Track quality vs context size used

### Setup Prometheus Monitoring (vLLM)

```bash
# Start vLLM with metrics export
python -m vllm.entrypoints.openai.api_server \
  --model=Qwen/Qwen1.5-32B-Chat \
  --metric-log-dir=/tmp/vllm_metrics

# Prometheus scrape config
cat > prometheus.yml <<EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'vllm'
    static_configs:
      - targets: ['localhost:8000']
EOF

# Run Prometheus
docker run -p 9090:9090 -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus
```

### Custom Monitoring Script
```python
import requests
import time

def monitor_vllm(endpoint="http://localhost:8000"):
    while True:
        try:
            resp = requests.get(f"{endpoint}/metrics")
            metrics = resp.text
            
            # Extract key metrics
            for line in metrics.split('\n'):
                if 'vllm:num_requests' in line:
                    print(line)
                elif 'vllm:time_to_first_token' in line:
                    print(line)
        except Exception as e:
            print(f"Error: {e}")
        
        time.sleep(10)

if __name__ == "__main__":
    monitor_vllm()
```

---

## Bottleneck Identification

### CPU-Bound
**Symptom**: GPU idle; high CPU utilization

**Cause**: Tokenization, quantization overhead, model loading

**Fix**:
- Async tokenization: Pre-tokenize offline
- Pre-load models: Keep in GPU memory
- Use GPU-side quantization (AWQ) instead of CPU quantization

### Compute-Bound
**Symptom**: GPU 90%+ utilization; low throughput

**Cause**: Underbatching or complex generation constraints

**Fix**:
- Increase batch size
- Enable speculative decoding (vLLM)
- Use constrained generation (grammar-based) only when necessary

### Memory-Bandwidth-Bound
**Symptom**: GPU 50-70% utilization; high VRAM throughput

**Cause**: Token-by-token generation (common in code agents)

**Fix**:
- Increase batch size (batch multiple requests together)
- Use longer sequences (but limits interactivity)
- Enable speculative decoding to skip low-confidence tokens

---

## Practical Examples

### Example 1: Quick Local Development (Ollama)
```bash
# Install
brew install ollama

# Run Qwen
ollama pull qwen:7b-code
ollama run qwen:7b-code

# In another terminal, test API
python3 << 'EOF'
import requests
import json

response = requests.post('http://localhost:11434/api/generate', json={
    'model': 'qwen:7b-code',
    'prompt': 'def reverse_list(lst):\n',
    'stream': False,
    'temperature': 0.3,
    'top_p': 0.9
})

result = response.json()
print(result['response'])
EOF
```

### Example 2: High-Throughput Production (vLLM)
```bash
# Install
pip install vllm autoawq

# Start server
python -m vllm.entrypoints.openai.api_server \
  --model=Qwen/Qwen1.5-32B-Chat \
  --quantization=awq \
  --tensor-parallel-size=2 \
  --max-num-batched-tokens=8192 \
  --gpu-memory-utilization=0.9 \
  --enable-prefix-caching

# Batch code review requests
python3 << 'EOF'
import openai

openai.api_base = "http://localhost:8000/v1"
openai.api_key = "EMPTY"

reviews = []
for code_file in code_files:
    reviews.append(openai.ChatCompletion.create(
        model="Qwen/Qwen1.5-32B-Chat",
        messages=[
            {"role": "system", "content": "You are an expert code reviewer."},
            {"role": "user", "content": f"Review this code:\n{code_file}"}
        ],
        temperature=0.3,
        max_tokens=1024
    ))
    
for review in reviews:
    print(review['choices'][0]['message']['content'])
EOF
```

### Example 3: Long-Context File Understanding
```python
from vllm import LLM, SamplingParams

llm = LLM(
    model="Qwen/Qwen1.5-32B-Chat",
    quantization="awq",
    max_model_len=32768,  # Full 32K context
    enable_prefix_caching=True,
    gpu_memory_utilization=0.9
)

# Read entire file
with open('large_module.py', 'r') as f:
    code = f.read()

prompt = f"""Analyze this Python module and identify:
1. Main functions and their purposes
2. Potential bugs or inefficiencies
3. Suggestions for improvement

Module:
{code}"""

outputs = llm.generate([prompt], SamplingParams(
    temperature=0.3,
    top_p=0.9,
    max_tokens=2048
))

print(outputs[0].outputs[0].text)
```

---

## Performance Expectations by Hardware

### RTX 3090 (24GB)
- **7B models**: 15-20 tokens/sec (GGUF), 18-25 tokens/sec (AWQ)
- **13B models**: 8-12 tokens/sec (GGUF), 12-15 tokens/sec (AWQ)
- **Context**: 8K tokens typical; 16K with paging

### RTX 4090 (24GB)
- **13B models**: 12-18 tokens/sec (GGUF), 15-22 tokens/sec (AWQ)
- **32B models (quantized)**: 5-10 tokens/sec (AWQ), 3-6 tokens/sec (GGUF)
- **Context**: 16K-32K tokens

### L40S (48GB)
- **32B models**: 15-25 tokens/sec (unquantized), 25-40 tokens/sec (AWQ)
- **70B models**: 8-12 tokens/sec (quantized)
- **Context**: Full 128K+ with paging

### Multi-GPU (Tensor Parallelism)
- **Dual RTX 3090**: ~2x single GPU throughput (minus 15-20% TP overhead)
- **Dual RTX 4090**: Similar; scales well for 32B+ models
- **Dual L40S**: ~1.8x single GPU (communication overhead on enterprise NICs)

---

## Cost/Performance: Local vs Cloud

| Approach | Setup Cost | Ongoing Cost | Break-Even | Best For |
|----------|------------|--------------|-----------|----------|
| **Local (Dual RTX 3090)** | $3-4k | $0/month | 6-8 months (24/7) | Dev + Testing |
| **Cloud (AWS g4dn.2xlarge)** | $0 | $540/month | — | Auto-scaling + Experimentation |
| **Local (Dual L40S)** | $8-10k | $0/month | 15-20 months | Production High-Throughput |
| **Cloud (AWS g5.24xlarge)** | $0 | $3-4k/month | — | Enterprise Multi-Tenancy |

**Decision**: Local for 24/7 dedicated inference; cloud for variable load or experimentation.

---

## Troubleshooting

### Out of Memory (OOM)
```bash
# Check current VRAM usage
nvidia-smi

# Strategies:
# 1. Reduce context window (vLLM --max-model-len=4096)
# 2. Use quantization (GGUF/AWQ)
# 3. Reduce batch size (--max-num-seqs=4)
# 4. Enable paging (--enable-prefix-caching)
# 5. Add GPU (tensor parallelism)
```

### Low Throughput
```bash
# Check GPU utilization
watch -n 1 nvidia-smi

# If <70% utilization: CPU-bound or compute-bound
# - Increase batch size
# - Check tokenization overhead (async if possible)
# - Enable speculative decoding (vLLM)

# If >90% utilization but low TPS: Memory-bandwidth-bound
# - Batch multiple requests
# - Use speculative decoding
```

### High Latency (TTFT)
```bash
# Check system state
nvidia-smi dmon  # Monitor throttling
top               # Check CPU load

# Strategies:
# 1. Preload model (avoid cold starts)
# 2. Reduce context window
# 3. Lower batch size (trade throughput for latency)
# 4. Enable speculative decoding
```

### Hallucinations in Code
```python
# Adjust sampling parameters
sampling_params = SamplingParams(
    temperature=0.1,  # Lower = more deterministic
    top_p=0.8,        # Stricter nucleus sampling
    repetition_penalty=1.1,  # Penalize repetition
    max_tokens=512    # Shorter outputs = fewer hallucinations
)
```

---

## References

- **vLLM Docs**: https://docs.vllm.ai
- **Ollama**: https://ollama.ai
- **LocalAI**: https://localai.io
- **llama.cpp**: https://github.com/ggerganov/llama.cpp
- **Paged Attention Paper**: https://arxiv.org/abs/2309.06393
- **AWQ Paper**: https://arxiv.org/abs/2306.00978
- **Speculative Decoding**: https://arxiv.org/abs/2211.17192
- **Flash Attention**: https://arxiv.org/abs/2205.14135

