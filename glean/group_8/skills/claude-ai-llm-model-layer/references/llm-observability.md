<!-- hub-reference-banner -->
> **Reference file — part of the `ai-llm-model-layer` hub.** Formerly the standalone `llm-observability` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

<!-- Provenance: reference under the `ai-agent-engineering` hub. Created 2026-05-31 via /dr deep-research from OpenTelemetry GenAI / OpenInference specs and vendor docs (Langfuse, LangSmith, Arize Phoenix, Datadog, WhyLabs, Traceloop). Scope: PRODUCTION observability of LLM/RAG/agent systems. NOT offline benchmark eval (HELM/MMLU/LLM-as-Judge benchmarking → da-7-machine-learning) and NOT building LLM integrations (→ references/llm-integration-reviewer.md). -->

# LLM / LLMOps Production Observability

Observing, tracing, and monitoring LLM, RAG, and agent applications **running in
production**. This is the operational counterpart to the build-time references in
this hub: once an agent or LLM app ships, you need to see what it does on live
traffic, what it costs, when it degrades, and when it produces unsafe or
ungrounded output.

## Scope boundary (read first)

- **This reference** = production observability: live traces/spans, online
  evaluation on real traffic, drift, runtime guardrail/PII monitoring, cost, RAG
  and agent-trajectory observability, and the tooling that captures it.
- **Offline benchmark evaluation** (HELM, MMLU, BIG-bench, leaderboards,
  LLM-as-Judge benchmarking, model selection bake-offs against curated datasets)
  → defer to the `da-7-machine-learning` skill. Online eval *reuses* LLM-as-judge
  scorers but runs them reference-free on production traffic — that part lives
  here.
- **Building / reviewing an LLM integration** (failover, retries, caching, rate
  limits, structured output, provider abstraction) → `references/llm-integration-reviewer.md`.
- **Designing the RAG pipeline itself** (chunking, embedding choice, retrieval
  design) → `references/rag-architecture.md`. **Observing** that pipeline's
  retrieval quality in production is here.

---

## 1. Tracing and spans — the data model

Production LLM/agent observability is built on **distributed tracing**. A request
becomes a **trace**; each unit of work is a **span**, nested into a parent-child
tree. Unlike HTTP tracing, the span tree must capture LLM-specific work: the
prompt sent, the completion returned, the model name, retrieved documents, tool
arguments, and token counts. Agent-specific failures (wrong tool selection,
drifted retrieval, infinite loops) do not show up in HTTP-level metrics — you
need the semantic span tree to see them.

Common span/observation types across the field (names vary by vendor, but the
taxonomy is consistent):

- **LLM / generation** — one model inference call (prompt, completion, model,
  tokens, latency, finish reason).
- **Chain / workflow** — a predetermined sequence grouping LLM calls plus
  contextual operations.
- **Agent** — a dynamic, agent-driven sequence of steps.
- **Tool / tool-call** — one external tool or function invocation.
- **Retriever / retrieval** — one vector-store or search query (with the
  documents it returned).
- **Reranker** — a reordering pass over retrieved documents.
- **Embedding** — one embedding generation.
- **Guardrail** — an input/output safety or validation check.
- **Evaluator** — an inline scoring step (e.g., an LLM-as-judge call).

Spans are grouped into **sessions / threads** (by `session_id`) so multi-turn
conversations can be evaluated as a coherent unit, and trace-level attributes
(`user_id`, `session_id`, tags, metadata) propagate down to every child span.

---

## 2. The standards: OpenTelemetry GenAI and OpenInference

Two open semantic-convention standards dominate. Both sit **on top of
OpenTelemetry**, so traces flow into any OTel-compatible backend (Datadog,
Honeycomb, Grafana/Tempo, Langfuse, etc.).

### OpenTelemetry GenAI semantic conventions

The official OTel spec for generative-AI spans. Still marked **development /
experimental** status — it moves fast, so treat the live spec as source of truth.
Key span attributes (verified against the OTel semconv spec):

- `gen_ai.operation.name` (required) — operation type. Well-known values:
  `chat`, `embeddings`, `execute_tool`, `generate_content`, `text_completion`.
- `gen_ai.provider.name` (required) — `openai`, `anthropic`, `aws.bedrock`,
  `gcp.vertex_ai`, `azure.ai.openai`, etc.
- `gen_ai.request.model` — the requested model name.
- `gen_ai.usage.input_tokens` / `gen_ai.usage.output_tokens` — token counts
  (input should include cached tokens).
- `gen_ai.response.finish_reasons` — array, e.g. `["stop"]`, `["stop","length"]`.
- Structured message content: `gen_ai.system_instructions`,
  `gen_ai.input.messages`, `gen_ai.output.messages`.

**Important migration note:** the older `gen_ai.prompt` and `gen_ai.completion`
attributes were **deprecated and removed** (around semconv v1.38, late 2025) in
favor of the structured `gen_ai.input.messages` / `gen_ai.output.messages`
attributes. Instrumentation pinned to the old names will emit deprecated fields —
check the version. There are also metric instruments: `gen_ai.client.token.usage`
and `gen_ai.client.operation.duration`.

### OpenInference (Arize)

A complementary semantic-convention spec, also OTel-based, created by Arize and
natively consumed by **Arize Phoenix** (and Arize AX). Its distinguishing feature
is a required `openinference.span.kind` attribute with a fixed taxonomy:
**CHAIN, LLM, RETRIEVER, RERANKER, EMBEDDING, AGENT, TOOL, GUARDRAIL, EVALUATOR,
PROMPT**. Attributes carry the prompt, response, model name, retrieved documents,
tool arguments, and token counts (`llm.token_count.*`, `retrieval.documents.*`,
`embedding.*`, `tool.*`, `input.value`, `output.value`). OpenInference ships
auto-instrumentors for many frameworks (LangChain, LlamaIndex, OpenAI, etc.).

**OTel GenAI vs OpenInference:** GenAI is the vendor-neutral OTel standard;
OpenInference is Arize's broader AI-app taxonomy (richer span-kind enum, predates
parts of GenAI). Many tools emit/accept both. When choosing, follow the backend
you target — Phoenix prefers OpenInference; Datadog and others lean OTel GenAI.

### OpenLLMetry / Traceloop

**OpenLLMetry** (by Traceloop) is an open-source set of OTel extensions plus a
**Traceloop SDK** that auto-instruments LLM providers (OpenAI, Anthropic) and
vector DBs (Pinecone, Chroma, Weaviate) non-intrusively, emitting standard OTel
data. Because the output is plain OTel, it exports to Traceloop's backend **or**
any existing observability stack. SDKs exist for Python, TypeScript/JS, Go, Ruby.

---

## 3. Cost, token, latency, and throughput monitoring

The operational baseline, derived from span attributes:

- **Token usage** — input/output tokens per call, aggregated per user, session,
  feature, model, prompt version. Track cached vs uncached tokens.
- **Cost** — token counts × per-model pricing, attributed per request/agent step.
  Cost-per-step is the unit that surfaces runaway agent loops.
- **Latency** — per-span and end-to-end. For streaming, **time-to-first-token
  (TTFT)** matters as much as total duration.
- **Throughput / error rate** — requests/sec, provider error and rate-limit
  counts, retry counts, finish-reason distribution (a spike in `length` finishes
  signals truncation).

These are the metrics that map cleanly onto traditional APM dashboards and
alerting; the AI-specific quality signals (sections 4–8) layer on top.

---

## 4. Online evaluation in production (vs offline)

| | Offline eval | Online eval |
| --- | --- | --- |
| When | Pre-deployment, in CI | After launch, on live traffic |
| Data | Curated, labeled datasets | Raw production traffic (messy, novel) |
| Method | **Reference-based** (compare to known ideal answer) | **Reference-free** (no labeled ground truth) |
| Examples | Regression suites, model bake-offs | Live LLM-as-judge scoring, guardrail hits |

Production traffic almost never arrives with labeled expected outputs, so online
evaluation relies on **reference-free scorers**:

- **LLM-as-judge** — a strong model scores each live output against criteria
  (helpfulness, groundedness, tone, policy compliance). Best practice: calibrate
  the judge against human labels, and keep the *same* scorer logic across offline
  and online so scores are interpretable in both. (LLM-as-judge for *benchmark*
  scoring is `da-7-machine-learning` territory; running it reference-free on
  production traffic is observability.)
- **Deterministic / heuristic scorers** — regex, JSON-schema validity, length,
  refusal detection, format checks. Cheap, run on every request.
- **Sampling** — judges are expensive, so score a sampled %, or trigger judging
  on flagged spans (low confidence, guardrail hits, user thumbs-down).

Scores attach back to spans/traces so you can slice quality by model, prompt
version, user cohort, and time.

---

## 5. Hallucination, groundedness, faithfulness, answer-relevance

The most-watched quality dimensions for production LLM/RAG apps:

- **Faithfulness / groundedness** — is the answer factually supported by the
  context the model was given? Treated as the single most critical production
  metric for RAG. Detection pipeline: **claim extraction** (break the answer into
  atomic claims) → **verification** (check each claim against retrieved context)
  → **scoring** (proportion of supported claims).
- **Answer relevance** — does the answer actually address the user's question
  (independent of whether it is grounded)?
- **Hallucination detection** — typically the inverse of faithfulness, run as an
  LLM-as-judge scorer on live traffic. Several platforms (Datadog, Traceloop,
  Patronus, Galileo) ship turnkey hallucination evaluators that score and alert.

These run as online evaluators (section 4) and feed dashboards/alerts; sustained
faithfulness drops are an early warning of retrieval rot or model/prompt change.

---

## 6. Drift monitoring: prompt, output, and embedding drift

Drift = the production distribution diverging from what you validated against:

- **Input / prompt drift** — the distribution of incoming prompts shifts (new
  topics, new phrasing, adversarial patterns) so the app faces inputs it was
  never tested on.
- **Output / response drift** — output distribution shifts (length, sentiment,
  refusal rate, topic mix). Can be caused by a silent provider model update,
  a prompt edit, or upstream data change.
- **Embedding drift** — the embedding-space distribution of inputs or retrieved
  content moves over time; a classic signal that a RAG corpus or query mix has
  changed. Detected by comparing embedding distributions across time windows.

Tooling: **WhyLabs + LangKit** is the reference open-source approach —
**whylogs** computes statistical *profiles* locally and ships only the profile
(not raw prompts/responses) to a central service, enabling drift monitoring
without exfiltrating sensitive text. LangKit adds LLM-specific signals (text
quality, sentiment, toxicity, topic, PII, jailbreak detection) on top.

---

## 7. Runtime guardrails, safety, and PII monitoring

Guardrails are **runtime control layers** that validate/filter/redact LLM inputs
and outputs *before* they reach the model or the user — and every guardrail
decision is an observability signal worth tracing (guardrail span, section 1).

Common risk categories (the OWASP Top 10 for LLM Applications is the canonical
taxonomy): prompt injection / jailbreak, PII leakage, toxic/unsafe output,
off-topic / topic drift, and hallucination.

Reference tooling:

- **NVIDIA NeMo Guardrails** — open-source (Apache 2.0) runtime orchestration
  layer; policies written in **Colang**; intercepts requests, calls classifier
  models, and blocks/rewrites traffic. Handles content moderation, PII redaction,
  topic relevance, jailbreak detection, dialog state.
- **Guardrails AI** — library-style; a `Guard` wraps the LLM call and runs
  **validators** (from the Guardrails Hub) for data leakage, factuality, safety,
  format/structured-output validation.
- Platform-native scanners — e.g., Datadog LLM Observability auto-scans and
  redacts sensitive data and flags prompt injection inline.

PII monitoring specifically: detect-and-redact at the guardrail layer, and many
profile-based monitors (LangKit) flag PII presence as a tracked metric. Note the
regulatory backdrop — EU AI Act high-risk obligations phase in through 2026 —
which is pushing guardrail + audit-trail observability from nice-to-have to
compliance requirement.

---

## 8. RAG observability

Observing a *deployed* RAG pipeline's retrieval and generation quality on live
traffic. The reference metric set is **RAGAS** (Retrieval Augmented Generation
Assessment) — a reference-free framework, well-suited to production because it
needs no human-labeled ground truth:

- **Faithfulness** — answer grounded in retrieved context (section 5).
- **Answer relevance** — answer addresses the question.
- **Context precision** — are the retrieved chunks relevant, and ranked with the
  relevant ones on top?
- **Context recall** — does the retrieved context actually contain the
  information needed to answer? (Low recall = retrieval gap.)

RAGAS deliberately replaces classic NLP metrics (BLEU/ROUGE) that correlate
poorly with RAG quality. In production these run as online evaluators on sampled
traffic, attached to retrieval + LLM spans so you can see *which stage*
(retrieval vs generation) is failing. Retrieval-stage failures (bad context
precision/recall) and generation-stage failures (low faithfulness despite good
context) demand different fixes — the span tree separates them.

---

## 9. Agent trajectory and tool-call observability

Agents add a failure surface that single-call metrics miss. Observability here is
**trajectory-level**:

- **Trajectory evaluation** — an LLM-as-judge assesses the *entire sequence* of
  tool calls an agent took: did the task complete, where did it fail, and why.
- **Tool-call observability** — per-tool latency, error patterns, and **tool-
  selection accuracy** (did the agent pick the right tool?) across the trace tree.
- **Loop / recursion detection** — trajectory mapping surfaces recursive patterns
  (which tool the agent keeps returning to before looping) — a top agent failure
  mode and a cost sink.
- **Multi-agent span trees** — nested spans must preserve parent-child
  relationships across agent handoffs, so a multi-agent system reads as one
  coherent trace.
- **Multi-level evaluation** — evaluators configurable at **session, trace, or
  span** level; per-step cost/token tracking ties spend to specific agent steps.

(For *designing* the agent's action space and tool definitions, see
`references/agent-harness-construction.md`; this section is about watching that
harness behave in production.)

---

## 10. Datasets, experiment tracking, and prompt management

The "engineering loop" features that production observability platforms bundle so
that production findings feed back into improvement:

- **Prompt management** — central store + **versioning** of prompt templates,
  with branching, approval workflows, and client/server caching to iterate without
  adding latency. Prompts are treated as first-class versioned artifacts so you
  can correlate a production quality regression with a specific prompt version.
- **Datasets** — curated sets of inputs/outputs/eval-metadata/span-context, with
  version control; production traces (especially failures) can be promoted into
  datasets to grow regression coverage.
- **Experiment tracking** — run a new prompt/model against a pinned dataset
  version in a playground, with full span-level tracing of each experiment, then
  compare experiments. This is the bridge between online findings and offline
  re-validation.

---

## 11. Tooling landscape

All of these emit/consume OTel and overlap heavily; pick by deployment model,
framework affinity, and depth needed.

| Tool | Type | Notes |
| --- | --- | --- |
| **Langfuse** | OSS, self-hostable | Framework-agnostic; full platform (tracing, evals, prompt mgmt, playground, datasets). Self-host stack = Postgres + **ClickHouse** (OLAP for traces/observations/scores) + Redis + S3. Is both an **OTel exporter and backend** (`/api/public/otel` ingestion). Data model: observations → traces → sessions. |
| **LangSmith** | Commercial (LangChain) | Deepest **LangChain/LangGraph** integration — node-by-node state diffs, full agent execution graphs, replay against new model versions. Session/thread grouping for multi-turn eval. |
| **Arize Phoenix** | OSS | Open-source counterpart to commercial **Arize AX**; native **OpenInference**; tracing + evals + prompt mgmt + playground; runs locally (Postgres). Arize AX is the enterprise SaaS (ClickHouse). |
| **WhyLabs + LangKit** | OSS/commercial | Profile-based (whylogs) — ships statistical profiles, not raw text; strong on **drift** (embedding/output distribution) and PII/toxicity/jailbreak signals without exfiltrating prompts. |
| **Helicone** | Commercial, **proxy** | Drop-in: change the API base URL, get traces — simplest install, no SDK changes; trade-off is shallower depth than native instrumentation; ~50–80ms added latency. |
| **Datadog LLM Observability** | Commercial APM | Enterprise default for Datadog shops; SDK auto-instruments OpenAI/LangChain/Bedrock/Anthropic; out-of-the-box evals (quality, topic relevancy, toxicity, sentiment, hallucination, prompt-injection) + sensitive-data scan/redact; **natively supports OTel GenAI semantic conventions**. |
| **OpenLLMetry / Traceloop** | OSS instrumentation | OTel extensions + Traceloop SDK (Python/TS/Go/Ruby); auto-instruments LLM providers + vector DBs; exports to Traceloop or any OTel backend. |

Adjacent eval-leaning platforms you'll also encounter: **Braintrust**,
**Patronus**, **Galileo**, **Maxim**, **Comet/Opik**, **PromptLayer**, **Vellum**,
**Honeycomb** (event-based deep tracing) — most are best classified as the
eval/experiment tier that overlaps observability.

---

## Quick decision guide

- **Need observability in prod today, minimal code change** → Helicone (proxy) or
  OpenLLMetry SDK into your existing backend.
- **Want a full OSS platform you can self-host** → Langfuse (or Arize Phoenix for
  OpenInference / ML-grade eval rigor).
- **Already a Datadog shop** → Datadog LLM Observability.
- **Heavy LangChain/LangGraph agent** → LangSmith.
- **Drift / PII without exfiltrating prompts** → WhyLabs + LangKit (profile-based).
- **Standardize the wire format first** → emit OTel GenAI and/or OpenInference so
  you can switch backends later.
