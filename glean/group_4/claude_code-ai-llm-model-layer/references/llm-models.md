<!-- hub-reference-banner -->
> **Reference file — part of the `ai-llm-model-layer` hub.** Formerly the standalone `llm-models` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: llm-models
version: "1.2.0"
updated: "2026-05-29"
description: >
  LLM landscape reference — Claude (Opus 4.7/Sonnet 4.6/Haiku 4.5), GPT (4.1/4o/o3/o4-mini),
  Gemini (2.5 Pro/Flash), open-source (Llama 4, DeepSeek, Qwen, Mistral, Phi), API patterns
  (structured output, tool calling, streaming, batch), model selection, fine-tuning
  (LoRA/QLoRA/DPO), local inference (Ollama, vLLM, llama.cpp, MLX), and embeddings.
  TRIGGER: user asks which LLM to use, needs model pricing or context window comparison,
  wants to implement structured output or tool calling, asks about fine-tuning or local
  inference, needs embedding model selection, or asks about prompt caching costs.
  SKIP: prompt engineering techniques (→ prompt-engineering); MCP server development
  (→ mcp-servers); agent orchestration patterns (→ agent-ecosystem).
origin: local
---

# LLM Models, APIs & Capabilities

Reference for the May 2026 LLM landscape. Verify pricing against provider docs if this information is more than 3 months old — model pricing and availability change frequently.

## When to Use

- Choosing which LLM to use for a task
- Comparing model pricing, context windows, or capabilities
- Implementing structured output, tool calling, or streaming
- Fine-tuning guidance (when to fine-tune, LoRA, QLoRA, DPO)
- Local inference setup (Ollama, vLLM, llama.cpp, MLX)
- Embedding model selection for RAG/vector search
- Comparing Responses API vs Chat Completions
- Understanding prompt caching cost savings

## When NOT to Use

- Prompt engineering techniques → `prompt-deep-optimizer` or `prompt-engineering`
- MCP server development → `mcp-servers` or `mcp-builder`
- Agent orchestration patterns → `agent-ecosystem`

---

## Model Selection Decision Tree

```
What is the primary constraint?
|
+-- COST --> Is quality critical?
|            +-- No  --> Gemini 2.5 Flash ($0.30/$2.50) or GPT-4.1 nano ($0.10/$0.40)
|            +-- Yes --> GPT-4.1 mini ($0.40/$1.60) or DeepSeek V4 (self-host, MIT)
|
+-- QUALITY --> What domain?
|              +-- Reasoning/math --> o3 ($0.40/$1.60) or o4-mini ($1.10/$4.40)
|              +-- Coding --> Sonnet 4.6 (79.6% SWE-bench) or Opus 4.7
|              +-- Long documents --> Gemini 2.5 Pro (1M) or Claude (1M)
|              +-- General --> Opus 4.7 or GPT-4.1
|
+-- SPEED --> Haiku 4.5 / Gemini Flash / GPT-4o-mini
|
+-- CONTEXT --> Llama 4 Scout (10M) > GPT-4.1 (1M) = Claude (1M) = Gemini (1M)
|
+-- PRIVACY --> Local: DeepSeek V4 (MIT) or Qwen 3.5 (Apache 2.0) via Ollama/vLLM
```

## Quick Selection Table

| Need | Budget | Balanced | Premium |
|------|--------|----------|---------|
| Cost | Gemini Flash ($0.30/$2.50) | GPT-4.1 mini ($0.40/$1.60) | Claude Opus ($5/$25) |
| Reasoning | o4-mini ($1.10/$4.40) | Sonnet 4.6 ($3/$15) | Opus 4.7 / o3 |
| Speed | Haiku 4.5 / Flash | GPT-4o-mini | Sonnet 4.6 |
| Context | Llama 4 Scout (10M) | GPT-4.1 / Claude (1M) | Gemini 2.5 Pro (1M) |
| Coding | DeepSeek V4 Pro (self-host) | Sonnet 4.6 (79.6% SWE) | Opus 4.7 |
| Multilingual | Qwen 3.5 (201 langs) | Mistral Medium 3.5 (80+) | Gemini 2.5 Pro |
| Agentic | GPT-4.1 mini (cheap loops) | Sonnet 4.6 | Opus 4.7 / o3 |

---

## Key Pricing (per MTok, May 2026)

| Model | Input | Output | Context | Notes |
|-------|-------|--------|---------|-------|
| Claude Opus 4.7 | $5.00 | $25.00 | 1M | 128K max output, adaptive thinking |
| Claude Sonnet 4.6 | $3.00 | $15.00 | 1M | 64K max output, extended + adaptive thinking |
| Claude Haiku 4.5 | $1.00 | $5.00 | 200K | Fastest latency, extended thinking |
| GPT-4.1 | $2.00 | $8.00 | 1M | 54.6% SWE-bench |
| GPT-4.1 mini | $0.40 | $1.60 | 1M | Best value general-purpose |
| GPT-4.1 nano | $0.10 | $0.40 | 1M | Classification/extraction |
| GPT-4o | $2.50 | $10.00 | 128K | Multimodal (text/image/audio) |
| o3 | $0.40 | $1.60 | 200K | 87.7% GPQA Diamond |
| o4-mini | $1.10 | $4.40 | 200K | 99.5% AIME 2025, fast reasoning |
| Gemini 2.5 Pro | $1.25 | $10.00 | 1M | $2.50/$15 above 200K tokens |
| Gemini 2.5 Flash | $0.30 | $2.50 | 1M | Free tier available for prototyping |
| DeepSeek V4 Pro | Self-host | Self-host | 1M | MIT license, 80.6 SWE-bench |

### Prompt Caching Comparison

| Provider | Cache hit discount | Write cost | TTL | Best practice |
|----------|-------------------|------------|-----|---------------|
| Anthropic | 90% (0.1x) | 1.25x (5 min) / 2x (1 hr) | 5 min or 1 hr | Stable system prompt first, volatile content last |
| OpenAI | 50% (0.5x) | Free | Auto | Automatic on repeated prefixes |
| Google | ~75% | Varies | Configurable | Explicit context caching API |

A 30K-token system prompt on Sonnet drops from $0.09 to $0.009/request with caching.

### The 2–3 Model Production Pattern

LLM prices dropped ~80% from 2025 to 2026. Standard production architecture:

1. **Fast/cheap tier** (80–95% of requests): Gemini Flash, GPT-4.1 mini, Haiku — classification, summarization, simple extraction
2. **Strong tier** (5–15% of requests): Sonnet 4.6, GPT-4.1 — complex generation, code, analysis
3. **Deep reasoning tier** (1–5% of requests): Opus 4.7, o3 — multi-step reasoning, agentic workflows

---

## Claude (Anthropic)

**Opus 4.7** — Flagship. 1M context, 128K max output. Jan 2026 knowledge cutoff. Adaptive thinking. Same agent loop as Claude Code.

**Sonnet 4.6** — Best production value. 1M context, 64K max output. 79.6% SWE-bench Verified. Supports both extended and adaptive thinking.

**Haiku 4.5** — Speed tier. 200K context. Fastest latency. Extended thinking. Use for classification, summarization, high-concurrency chat.

**Managed Agents** — Public beta April 8, 2026. $0.08/session-hour + standard token rates.

### Claude Structured Output (Python)

```python
from pydantic import BaseModel
from anthropic import Anthropic

class ContactInfo(BaseModel):
    name: str
    email: str
    plan_interest: str
    demo_requested: bool

client = Anthropic()
response = client.messages.parse(
    model="claude-sonnet-4-6-20250514",
    max_tokens=1024,
    messages=[{"role": "user", "content": "Extract: John Smith (john@co.com) wants Enterprise demo."}],
    output_format=ContactInfo,
)
print(response.parsed_output)
```

### Claude Tool Use (Python)

```python
import anthropic
client = anthropic.Anthropic()
response = client.messages.create(
    model="claude-sonnet-4-6-20250514",
    max_tokens=1024,
    tools=[{
        "name": "get_weather",
        "description": "Get current weather for a location",
        "input_schema": {
            "type": "object",
            "properties": {
                "location": {"type": "string", "description": "City name"},
                "unit": {"type": "string", "enum": ["celsius", "fahrenheit"]}
            },
            "required": ["location"]
        }
    }],
    messages=[{"role": "user", "content": "What's the weather in San Francisco?"}]
)
```

---

## OpenAI

**GPT-4.1 family** — April 2025. 1M context. Three sizes: GPT-4.1 ($2/$8), mini ($0.40/$1.60), nano ($0.10/$0.40). 26% better cost efficiency than GPT-4o.

**Reasoning models (o-series):**
- **o3**: 87.7% GPQA Diamond. $0.40/$1.60. Combines reasoning with tool use.
- **o4-mini**: Fast reasoning at $1.10/$4.40. 99.5% AIME 2025 with code interpreter.

**Responses API vs Chat Completions:** Responses API (March 2025) is the recommended new primitive — server-managed conversation state, built-in tools (web search, file search, code interpreter, remote MCPs), 3% better SWE-bench, 40–80% improved cache utilization.

---

## Google Gemini

**Gemini 2.5 Pro** — 1M context. $1.25/$10 under 200K, $2.50/$15 above. Grounding with Google Search. Supervised fine-tuning on Vertex AI.

**Gemini 2.5 Flash** — 1M context, 65K max output. $0.30/$2.50. Free tier on Gemini Developer API. Flash-Lite for extreme cost sensitivity.

---

## Open-Source Models

Almost every flagship 2026 open model is sparse MoE — frontier quality at a fraction of the compute:

| Model | Total params | Active params | License | Context |
|-------|-------------|---------------|---------|---------|
| DeepSeek V4 Pro | 1.6T | 49B | MIT | 1M |
| Llama 4 Maverick | 400B | 17B | Llama License | 1M |
| Qwen 3.5 | 397B | 17B | Apache 2.0 | 256K |
| Mistral Large 3 | 675B | 41B | Apache 2.0 | 128K |
| Llama 4 Scout | 109B | 17B | Llama License | 10M |

**DeepSeek V4 Pro** — 80.6 SWE-Bench, 90.1 GPQA Diamond. MIT license. Best overall open-source as of May 2026.

**Llama 4 Scout** — 10M context window. Ideal for massive document processing.

**Qwen 3.5** — 201 language support. Runs on a MacBook with 64GB RAM at quantized sizes. Apache 2.0.

### Licensing Guide

| License | Models | Commercial use |
|---------|--------|----------------|
| MIT | DeepSeek V4, GLM-5, Phi-4 | Unrestricted |
| Apache 2.0 | Qwen 3.5, Mistral, Gemma 4 | Unrestricted |
| Llama License | Llama 4 family | Yes with conditions (EU multimodal limits, MAU caps) |
| CC-BY-NC | Command R+ | Non-commercial only |

---

## API Patterns

### Structured Output (three approaches)

| Approach | Guarantee | OpenAI | Claude | Gemini |
|----------|-----------|--------|--------|--------|
| JSON Mode | Valid JSON, no schema | `response_format: {"type": "json_object"}` | System prompt instruction | `response_mime_type` |
| Structured Output | 100% schema compliance | Pydantic schemas via SDK | `messages.parse()` with Pydantic | Native with validation |
| Tool/Function Calling | Structured params | `tools` parameter | `tools` parameter | `function_declarations` |

### Streaming Pattern (SSE)

```python
# Claude
with client.messages.stream(model="claude-sonnet-4-6-20250514", max_tokens=1024,
    messages=[{"role": "user", "content": "Explain quantum computing"}]) as stream:
    for text in stream.text_stream:
        print(text, end="", flush=True)

# OpenAI
stream = openai_client.chat.completions.create(
    model="gpt-4.1", messages=[{"role": "user", "content": "Explain quantum computing"}],
    stream=True)
for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="")
```

### Batch APIs

All three major providers offer batch APIs at reduced cost:
- **Anthropic**: 50% discount, 300K max output tokens via beta header
- **OpenAI**: 50% discount, 24-hour completion window
- **Google**: Batch prediction API on Vertex AI

Use for: data labeling, bulk classification, document summarization, evaluation runs.

---

## Fine-Tuning

### Decision Framework

```
Is the problem behavior/style or knowledge?
|
+-- Knowledge --> Use RAG
|
+-- Behavior/style --> Does prompt engineering solve it?
                       +-- Yes --> Stay on prompts (70% of cases)
                       +-- No  --> Volume >10K requests/day?
                                   +-- Yes --> Fine-tune (ROI in 1–2 months)
                                   +-- No  --> Try few-shot examples first
```

**Key principle:** Fine-tuning is for form, not facts — shape behavior, style, and output format. Do not use fine-tuning to inject knowledge that changes frequently.

### Methods

| Method | What it does | Data needed | When to use |
|--------|-------------|-------------|-------------|
| LoRA | Low-rank adapter matrices (0.1–1% of params) | 500–5K examples | Default choice |
| QLoRA | LoRA + 4-bit NF4 quantization (~$12 for 8B/50K) | 500–5K examples | Budget-constrained, single-GPU |
| DPO | Direct preference optimization from pairs | Preference pairs | Alignment, tone, safety |
| ORPO | SFT + preference in one step | Small datasets | SFT + alignment together |
| RLHF | Reward model + RL | Complex reward shaping | Only when DPO is insufficient |

**QLoRA is the default in 2026:** single A100 80GB fine-tunes Llama 3 8B in 6 hours on 50K examples for ~$12.

**Data quality rule:** 1,000 hand-curated examples often beat 100,000 noisy ones.

---

## Local Inference

### Framework Comparison

| Framework | Best for | Throughput (Apple Silicon) | Key feature |
|-----------|----------|---------------------------|-------------|
| Ollama | Developer adoption, teams | ~150 tok/s | One-command install, OpenAI-compatible API |
| llama.cpp | Maximum control, portability | ~150 tok/s | Created GGUF format, runs on everything |
| vLLM | Production GPU serving | 35x+ throughput at load | PagedAttention, continuous batching |
| MLX | Apple Silicon optimization | ~230 tok/s | Unified memory architecture |
| TGI | Containerized serving | High (GPU) | HuggingFace ecosystem integration |

**Progression:** LM Studio (discovery) → Ollama (development) → vLLM (production)

### Quantization Guide

| Format | Quality retention | Best for |
|--------|-----------------|----------|
| Q4_K_M (GGUF) | ~92% | Sweet spot: 75% size reduction |
| Q5_K_M (GGUF) | ~95% | When quality matters more |
| AWQ | ~95% | Production GPU |
| GPTQ | ~90% | GPU inference |

**Minimum VRAM:** 7B models need 6GB at Q4 quantization.

### Hardware Recommendations

| Hardware | Price | Best for | Max model |
|----------|-------|----------|-----------|
| Mac Mini M4 Pro 48GB | ~$1,999 | Best value for 70B | 70B Q4 |
| Mac Studio M4 Ultra 192GB | ~$4,999 | 400B+ MoE | Maverick/Qwen 3.5 |
| RTX 3090 24GB (used) | ~$800–1K | Budget GPU for 13B | 13B FP16, 30B Q4 |
| RTX 4090 24GB | ~$1,600 | Single-GPU production | 30B Q4 |

---

## Embeddings

### Top Models (May 2026)

| Model | MTEB | Price/MTok | Dims | Max tokens | Best for |
|-------|------|-----------|------|------------|----------|
| Gemini Embedding 2 | 68.32 | Free/$0.006 | 768 | 8,192 | Multimodal (text+image+video+audio) |
| Voyage voyage-3-large | 67.8 | $0.18 | 1,024 | 32,000 | Max accuracy, code/legal/medical |
| Cohere embed-v4 | 65.2 | $0.10 | 1,024 | 512 | Enterprise multilingual (100+ langs) |
| OpenAI text-embedding-3-large | 64.6 | $0.13 | 3,072 | 8,191 | General + dimension reduction |
| OpenAI text-embedding-3-small | — | $0.02 | 1,536 | 8,191 | Cheap general default |
| BGE-M3 (OSS) | — | Free | 1,024 | 8,192 | 100+ languages, no API cost |

### Embedding Selection

```
General RAG (English) → OpenAI text-embedding-3-small ($0.02/MTok)
Maximum accuracy → Voyage voyage-3-large or Cohere embed-v4
High-volume budget → Google text-embedding-005 ($0.006)
Multilingual → Cohere embed-v4 or BGE-M3
Multimodal (text+images) → Gemini Embedding 2
Self-hosted / privacy → BGE-M3, Nomic, or GTE
Code search → Voyage voyage-code-3
```

---

## Model Version Pinning

Always pin model versions in production to avoid behavior drift:

| Provider | Pinned format | Alias (avoid in prod) |
|----------|--------------|----------------------|
| Anthropic | `claude-sonnet-4-6-20250514` | `claude-sonnet-4-6-latest` |
| OpenAI | `gpt-4.1-2025-04-14` | `gpt-4.1` |
| Google | `gemini-2.5-flash-preview-05-20` | `gemini-2.5-flash` |

Re-evaluate pinned versions quarterly or when the provider announces deprecation.
