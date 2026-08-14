<!-- hub-reference-banner -->
> **Reference file — part of the `ai-mcp-sdk-prompting` hub.** Formerly the standalone `ai-languages` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: ai-languages
version: 1.1.0
updated: 2026-05-29
description: >
  AI/LLM programming languages and frameworks reference — Python (LangChain,
  LlamaIndex, DSPy, Instructor, Pydantic AI, CrewAI, Marvin, Outlines),
  TypeScript (Vercel AI SDK, Mastra), Rust (Rig, Burn, Candle), Go (LangChainGo,
  Genkit), prompt engineering languages (LMQL, Guidance, SGLang), evaluation
  frameworks (DeepEval, PromptFoo, Braintrust, RAGAS), and language selection
  guidance. TRIGGER: user asks which language or framework to use for an AI
  project, needs to compare Python vs TypeScript vs Rust vs Go for AI work, wants
  help with a specific AI framework, or needs LLM evaluation setup.
  SKIP: general Python/TypeScript/Rust/Go programming unrelated to AI/LLM (use
  language-specific skills); model selection or API pricing (use llm-models);
  agent orchestration architecture (use agent-ecosystem); MCP server development
  (use mcp-servers or mcp-builder).
origin: local
tags: [python, typescript, rust, go, langchain, pydantic-ai, llamaindex, dspy, instructor, mastra, evaluation]
related_skills: [agent-ecosystem, llm-models, mcp-servers, ai-datastores]
---

# AI Programming Languages & Frameworks

Comprehensive reference for languages and frameworks used to build AI/LLM applications. Deep context at `references/ai-languages-context.md`.

## Quick language selection

| Scenario | Language | Framework |
| --- | --- | --- |
| Default for new projects | Python | Pydantic AI (or LangChain for a specific connector) |
| RAG pipeline | Python | LlamaIndex |
| Prompt optimization at scale | Python | DSPy |
| Structured data extraction | Python | Instructor |
| Web app with AI chat | TypeScript | Vercel AI SDK |
| TypeScript agent framework | TypeScript | Mastra |
| Performance-critical inference | Rust | Rig or Burn |
| Edge/IoT deployment | Rust | Candle or Burn |
| AI infrastructure / API gateway | Go | Genkit or LangChainGo |
| CI/CD eval gate | Python | DeepEval + RAGAS |
| Prompt security testing | Any | PromptFoo |

## Quick framework comparison

| Framework | Language | Strength | Use when |
| --- | --- | --- | --- |
| LangChain | Python | 700+ connectors, broadest ecosystem | Integration-heavy projects |
| LlamaIndex | Python | Document ingestion, RAG | Retrieval quality is critical |
| DSPy | Python | Automatic prompt optimization | You have metrics + labeled data |
| Instructor | Python | Structured output via Pydantic | Extracting typed data from LLMs |
| Pydantic AI | Python | Type-safe agents, FastAPI-style DI | Teams using Pydantic/FastAPI |
| Outlines | Python | Constrained generation via FSM | Guaranteed schema compliance (self-hosted) |
| Marvin | Python | AI decorators, minimal boilerplate | Sprinkling AI into existing code |
| Vercel AI SDK | TypeScript | Streaming UI, generative components | React/Next.js apps |
| Mastra | TypeScript | Zod schemas, agent workflows | TypeScript-first agent development |
| CrewAI | Python | Role-based multi-agent teams | Business process automation |
| Rig | Rust | Composable LLM traits | Performance-critical Rust agents |
| Genkit | Go | Batteries-included, Google Cloud | Go AI services |

---

## Python frameworks — deep dive

### LangChain (v0.3+, 2026)

**Philosophy:** Breadth of integration and a declarative chaining DSL (LCEL) for composing LLM pipelines.

**Production status:** 700+ third-party connectors as of May 2026. Ecosystem includes LangGraph (stateful multi-agent orchestration), LangSmith (observability/evaluation), and LangServe (deployment).

**Strengths:**
- Fastest time-to-prototype for integration-heavy projects
- LangGraph models workflows as directed graphs with checkpointing and HITL
- Largest community, most tutorials

**Weaknesses:**
- Frequent breaking changes between 0.1, 0.2, and 0.3 API versions
- LCEL learning curve; implicit type coercion inside chains is hard to debug
- Heavy abstraction overhead for simple use cases

**When to choose:** You need a specific connector it supports, your team already knows it, or you need LangGraph's stateful multi-agent orchestration.

**When to avoid:** Greenfield projects where type safety and maintainability outweigh integration breadth.

### Pydantic AI (v1.0+, stable since September 2025)

**Philosophy:** Python's type system is the best guardrail for LLM applications. Full Pydantic v2 validation, native async streaming, and built-in dependency injection.

**Production status:** V1 reached September 2025 — no breaking API changes until V2. Recommended default for new production projects in 2026.

**Strengths:**
- Type-safe agent definitions with full IDE autocomplete
- Dependency injection makes unit testing trivial (swap LLM clients in tests)
- Native async streaming
- Models workflows as typed state machines

**Weaknesses:**
- Smaller connector ecosystem than LangChain (growing rapidly)
- Fewer tutorials and community content vs LangChain

**When to choose:** New production project, teams using Pydantic/FastAPI, type safety is a priority, or you want stable APIs.

**When to avoid:** You need a niche LangChain connector with no Pydantic AI equivalent.

### LlamaIndex (2026)

**Philosophy:** A data framework for connecting LLMs to data sources, structuring indexes, and synthesizing responses.

**Key primitives:** Agents, Workflows, Tools, Memory, RAG pipelines, Evals.

**2025–2026 launches:** LlamaAgents (one-click deployment), LlamaParse v2 (up to 50% cost reduction), LlamaSheets, LlamaSplit.

**When to choose:** RAG is your core use case, document processing is complex, or you need production-grade retrieval pipelines.

**When to avoid:** Application does not need retrieval (pure chat, code generation), or you need broad non-RAG agent orchestration (use LangGraph or CrewAI).

### CrewAI (enterprise multi-agent, 2026)

**Philosophy:** Models workflows as teams — agents with explicit roles, goals, and backstories. Maps directly to business process automation where you think in terms of "who does what."

**When to choose:** Multi-agent business process automation, teams that think in roles/responsibilities, or when stakeholder communication about agent architecture matters.

**When to avoid:** Simple single-agent tasks, performance-critical pipelines, or when you need fine-grained graph control (use LangGraph).

### DSPy (v2.x, 2026)

**Philosophy:** Optimize prompts programmatically instead of manually. Define what the LLM should do (signatures), how it should reason (modules), and how to measure quality (metrics).

**MIPROv2 — the key optimizer:**
- Proposes new instructions and few-shot demonstrations, evaluates on mini-batches, keeps what moves your metric
- Outcomes: higher downstream accuracy, fewer hallucinations, reusable prompt assets

**Production deployment pattern:**
1. Start with automatic few-shot selection
2. Add MIPROv2 when output formats wobble
3. Fine-tune only if traffic volume or distribution drift demands it
4. Validate each module on 25–100 examples before touching production

**When to choose:** You have metrics and labeled data, prompts need systematic optimization, or you want reproducible prompt engineering.

**When to avoid:** You lack evaluation data, or your use case is simple enough that manual prompting suffices.

### Instructor (3M+ monthly downloads)

**Philosophy:** Structured data extraction from LLMs using Pydantic models as schemas, with automatic validation and retries.

**Key features:**
- Define Pydantic models specifying exactly what data you want
- Automatic retry logic when validation fails
- Works with 15+ LLM providers (OpenAI, Anthropic, Google, Mistral, Cohere, Ollama, DeepSeek)
- Multi-language: Python, TypeScript, Go, Ruby

**Pattern:**
```python
import instructor
from pydantic import BaseModel

client = instructor.from_openai(openai.OpenAI())

class UserInfo(BaseModel):
    name: str
    age: int

user = client.chat.completions.create(
    model="gpt-4o",
    response_model=UserInfo,
    messages=[{"role": "user", "content": "John is 25 years old"}]
)
```

**When to choose:** You need typed data extraction from LLMs with post-generation validation and retries.

### Outlines (constrained generation)

**Philosophy:** Use a finite state machine (FSM) to mask invalid tokens during generation. The model can only produce schema-compliant output — no post-hoc parsing or retry loops.

**Key difference from Instructor:** Instructor validates after generation and retries on failure. Outlines constrains during generation so invalid output is impossible.

**When to choose:** You control the inference runtime (self-hosted models), need guaranteed schema compliance on first pass, or retry latency is unacceptable.

**When to avoid:** Using hosted API providers (OpenAI, Anthropic) where you cannot modify the decoding process.

### Marvin (lightweight AI functions)

**Philosophy:** Expose LLM capabilities as simple Python functions and types via decorators (`@ai_fn`, `@ai_model`, `@ai_classifier`).

**When to choose:** You want to add AI capabilities to an existing Python codebase with minimal boilerplate.

**When to avoid:** Complex multi-step agent workflows, RAG pipelines, or when you need the orchestration power of LangChain/Pydantic AI/CrewAI.

---

## TypeScript frameworks — deep dive

### Vercel AI SDK (v6, 2026)

**Philosophy:** Lowest-friction path from a prompt to a production UI. Intentionally lower-level than Genkit or Mastra, favoring simplicity and composability.

**Strengths:**
- First-class React Server Components integration
- Streaming UI primitives (useChat, useCompletion, useAssistant)
- Generative UI components
- Model-agnostic provider system
- Edge runtime compatible

**When to choose:** React/Next.js apps with AI features, streaming chat UIs, generative components.

**When to avoid:** You need agent orchestration (use Mastra on top), or you are not in the React ecosystem.

### Mastra (v1.0, January 2026)

**Philosophy:** TypeScript-first agent framework built on Vercel AI SDK. Six primitives: agents, workflows, tools, memory, RAG, and evals.

**Production status:** 22,000+ GitHub stars within 15 months. 300,000+ weekly npm downloads by v1.0. Developer experience scored 9/10 vs 5/10 for LangChain in December 2025 benchmark.

**Strengths:**
- Zod-based schema validation throughout
- First-class workflow orchestration with state machines
- Built-in evaluation framework
- From the team behind Gatsby — strong DX pedigree

**When to choose:** TypeScript-first agent development, content or research workflows with multiple collaborating agents.

**When to avoid:** You need Python-specific ML library integrations, or your team prefers Python.

---

## Rust frameworks — deep dive

### Rig (production-viable, 2026)

**Philosophy:** Build modular and scalable LLM applications in Rust with a unified LLM interface and composable traits.

**Performance benchmarks vs Python:**
- 5x memory reduction (stays under 1.1GB peak vs 5+ GB for Python equivalents)
- 25–44% latency improvements
- Orders-of-magnitude better cold start times
- These are structural advantages, not tunable by configuration

**When to choose:** Performance-critical AI agents, high-throughput API services, edge deployments, teams with Rust expertise.

**When to avoid:** Rapid prototyping, teams without Rust experience, or when Python library ecosystem coverage is essential.

### Burn (backend-agnostic ML)

**Philosophy:** Write once, deploy anywhere. Backend-agnostic architecture supporting CPU, GPU, and WebAssembly execution via CubeCL integration.

**When to choose:** Custom model training in Rust, cross-platform deployment (GPU + WASM), research requiring low-level control.

### Candle (lightweight inference)

**Philosophy:** Minimalist ML framework focused on inference performance.

**When to choose:** Edge/IoT deployment, minimal binary size requirements, inference-only workloads.

---

## Go frameworks — deep dive

### Genkit Go (production-ready, 2026)

**Philosophy:** Batteries-included framework from Google giving Go developers structured Gen AI with streaming, evaluation, and tracing built in.

**Strengths:**
- Native streaming support
- Built-in evaluation and tracing
- Google Cloud integration
- Structured output support

**When to choose:** Go AI services, Google Cloud environments, teams needing built-in observability.

### LangChainGo

**Philosophy:** Go implementation of LangChain's composable components — chains, agents, tools, embeddings, and integrations with local and remote LLM providers.

**When to choose:** Teams familiar with LangChain patterns who want to stay in Go.

### Other Go options (2026)

- **Eino** — Go-native design, idiomatic patterns
- **Anyi** — Specialized AI agent solutions

---

## Prompt engineering languages

### SGLang (structured generation runtime)

Accelerates LLM serving using xgrammar to mask logits at each decoding step, ensuring the model can only emit tokens that keep output valid against your schema. No post-hoc parsing or retry loops. Best for high-throughput structured generation on self-hosted models.

### LMQL (query language for LLMs)

Introduces syntax for automatically transforming Python strings into LLM prompts. Supports constraint-based decoding. Best for complex prompt templates with conditional logic.

### Guidance (Microsoft)

Adds LLM reasoning, tool usage, and template guidance syntax to Python. Pioneers the template-as-program approach. Best for fine-grained control over generation and interleaving computation with generation.

---

## Evaluation frameworks — deep dive

| Tool | Best for | Integration style |
| --- | --- | --- |
| DeepEval | Python CI/CD quality gates | pytest plugin, 14+ metrics |
| PromptFoo | Red-teaming, adversarial testing, multi-model comparison | YAML-driven CLI, 40+ red-team plugins |
| Braintrust | Production tracing + human annotation + CI enforcement | Full platform, dataset management |
| RAGAS | RAG-specific retrieval/generation scoring | Research-backed metrics library |

**Recommended combinations (2026 consensus):**
- **DeepEval + Braintrust** — de facto standard for engineering-led AI product teams
- **PromptFoo + LangSmith** — security-focused teams with LangChain stack
- **RAGAS + DeepEval** — RAG applications needing both retrieval and generation metrics

You almost certainly need two tools: a lightweight framework for CI/CD gating paired with a platform for human annotation, regression tracking, and stakeholder dashboards.

---

## Language selection decision tree

**Start: What is your primary constraint?**

1. **Performance/memory critical?** — Rust: Rig for agents, Burn for training, Candle for lightweight inference
2. **Go team / Go infrastructure?** — Go: Genkit for Google Cloud, LangChainGo otherwise
3. **React/Next.js frontend?** — TypeScript: Vercel AI SDK for UI, add Mastra for agent orchestration
4. **Everything else?** — Python (continue below)

**Python sub-decision:**

1. Need 700+ integrations? — LangChain
2. New project, type safety matters? — Pydantic AI
3. RAG is the core feature? — LlamaIndex
4. Programmatic prompt optimization? — DSPy
5. Structured data extraction? — Instructor
6. Multi-agent with role-based teams? — CrewAI
7. Multi-agent with graph control? — LangGraph
8. Multiple constraints match equally? — Default to Pydantic AI; use LangChain only if you need a specific connector it has

## Migration paths (2026)

| From | To | Effort | Reason |
| --- | --- | --- | --- |
| LangChain 0.1/0.2 | LangChain 0.3 | Medium | Breaking changes require migration anyway |
| LangChain | Pydantic AI | Medium–High | Type safety, stability, simpler debugging |
| LangChain | LlamaIndex | Medium | Better RAG performance, cleaner data pipeline |
| Vercel AI SDK alone | Mastra | Low | Add agent orchestration on top of existing SDK |
| Custom Python | DSPy | Low–Medium | Systematic prompt optimization without rewriting logic |
| Any Python | Rig (Rust) | High | 5x memory reduction, but requires Rust expertise |

---

## Production readiness summary (May 2026)

| Framework | Stability | Weekly downloads | Recommendation |
| --- | --- | --- | --- |
| LangChain 0.3 | Stable | 2M+ | Use for integration breadth |
| Pydantic AI v1 | Stable (frozen until v2) | 200K+ | Default for new Python projects |
| LlamaIndex | Stable | 500K+ | Default for RAG |
| DSPy 2.x | Maturing | 100K+ | Use when you have eval data |
| CrewAI | Stable | 400K+ | Default for role-based multi-agent |
| Instructor | Stable | 3M+ | Default for structured extraction |
| Vercel AI SDK v6 | Stable | 1M+ | Default for React AI apps |
| Mastra v1.0 | Stable | 300K+ | Default for TS agent development |
| Rig | Production-viable | Early (crates.io) | Default for Rust AI agents |
| Genkit Go | Production-ready | Early (Go modules) | Default for Go AI services |

---

## Sources

All sources retrieved May 2026.

1. [Pydantic AI vs LangChain 2026](https://www.kunalganglani.com/blog/pydantic-ai-vs-langchain)
2. [Vstorm OSS Framework Comparison](https://oss.vstorm.co/blog/choosing-ai-framework/)
3. [LlamaIndex Complete Guide 2026](https://www.agentframeworkhub.com/blog/llamaindex-complete-guide-2026)
4. [DSPy Framework Guide 2026](https://myengineeringpath.dev/tools/dspy-guide/)
5. [Instructor Official Docs](https://python.useinstructor.com/)
6. [Mastra AI Complete Guide 2026](https://www.generative.inc/mastra-ai-the-complete-guide-to-the-typescript-agent-framework-2026)
7. [Rig Official Site](https://rig.rs/)
8. [Rust AI Agent Frameworks 2026](https://zylos.ai/research/2026-04-01-rust-native-ai-agent-frameworks-ecosystem-2026)
9. [AI and Go 2026](https://appliedgo.net/spotlight/ai-and-go/)
10. [SGLang: Fast LLM Inference 2026](https://markaicode.com/sglang-fast-llm-inference-structured-generation/)
11. [Promptfoo vs DeepEval 2026](https://genai.qa/blog/promptfoo-vs-deepeval/)
12. [Python vs Go vs Rust for AI Agents 2026](https://dev.to/thedailyagent/python-vs-go-vs-rust-for-ai-agents-in-2026-a-pragmatic-field-guide-5fda)
