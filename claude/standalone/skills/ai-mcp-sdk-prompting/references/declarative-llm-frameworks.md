<!-- hub-reference-banner -->
> **Reference file — part of the `ai-mcp-sdk-prompting` hub.** Formerly the standalone `declarative-llm-frameworks` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: declarative-llm-frameworks
title: Declarative & Programmatic LLM Frameworks
description: >
  The "prompt-as-program" framework layer — declare the interface (inputs,
  outputs, types, intent) and let a compiler, type system, or constrained
  decoder produce and enforce the prompt and output, instead of hand-writing
  brittle prompt strings. Three families: (1) compiled/optimized prompting —
  DSPy (signatures → modules like Predict/ChainOfThought/ReAct → optimizer
  compile step); (2) schema-first typed LLM functions — BAML (.baml DSL +
  Schema-Aligned Parsing + test-driven prompts + multi-language codegen),
  Instructor, Pydantic AI, Marvin, Mirascope, Ell; (3) constrained generation /
  decoding — Outlines (FSM), Guidance (token healing), LMQL (constraint query
  language) enforcing regex/CFG/JSON-schema at the token level. Plus the app-level
  structured-output substrate (native strict Structured Outputs, the now-obsolete
  plain JSON mode, function-calling-as-extraction, validate-and-reask loops) and
  test-driven/versioned prompting. TRIGGER: "prompt as program", DSPy
  signatures/modules/compile, BAML / schema-first typed prompts, constrained
  decoding / Outlines / Guidance / LMQL, Instructor / Pydantic AI / Mirascope,
  structured output / JSON schema enforcement / function-calling extraction,
  validate-and-retry, "should I use DSPy vs LangChain", testable/portable prompts.
  SKIP: the prompt-optimization ALGORITHMS themselves — MIPROv2 / BootstrapFewShot
  / COPRO / GEPA internals (use prompt-deep-optimizer / prompt-helper-optimizer);
  hand-writing a single prompt's wording (use prompt-engineering / phe); agent
  orchestration / routing / memory (use agent-ecosystem); serving-level constrained
  decoding XGrammar/Outlines runtime tuning (use llm-inference-serving).
origin: local
category: developer
version: "1.0"
updated: "2026-05-31"
tags:
  - llm
  - dspy
  - baml
  - structured-output
  - constrained-decoding
  - prompt-as-program
  - ai
whenToUse:
  - "treating prompts as compiled/typed programs rather than hand-written strings"
  - "using DSPy (signatures, modules, the optimize/compile step) as a framework"
  - "schema-first typed LLM functions (BAML, Instructor, Pydantic AI, Marvin, Mirascope)"
  - "constrained generation / decoding (Outlines FSM, Guidance token healing, LMQL)"
  - "reliable structured output (native strict mode, function-calling extraction, validate+reask)"
  - "test-driven / versioned prompting"
  - "deciding when declarative/compiled prompting beats hand-prompting"
whenNotToUse:
  - "the prompt-optimization ALGORITHM internals (MIPROv2/BootstrapFewShot/COPRO/GEPA) — use prompt-deep-optimizer / prompt-helper-optimizer"
  - "wording a single one-off prompt — use prompt-engineering / phe"
  - "agent orchestration / routing / memory — use agent-ecosystem"
  - "serving-level constrained decoding runtime (XGrammar/Outlines in vLLM) — use llm-inference-serving"
related_skills:
  - prompt-engineering
  - prompt-deep-optimizer
  - prompt-helper-optimizer
  - llm-inference-serving
---

# Declarative & Programmatic LLM Frameworks

The paradigm: **stop writing prompt strings; declare the interface and let a
compiler / type system / constrained decoder produce and enforce the prompt and
output.** DSPy's slogan — "programming, not prompting" — captures it. Hand-written
templates are "brittle and unscalable"; treating prompts as compiled or generated
artifacts unlocks four properties strings lack: **automatic optimization,
portability across models, testability, maintainability.**

The space splits into three loosely-coupled families (often combined):

| Family | Idea | Tools |
|---|---|---|
| **Compiled / optimized prompting** | Declare signatures + modules; an optimizer compiles them into prompts from training data + a metric | DSPy |
| **Schema-first typed LLM functions** | A prompt is a typed function: declare the return type, get validated structured output (+ retries) | BAML, Instructor, Pydantic AI, Marvin, Mirascope, Ell |
| **Constrained generation / decoding** | Enforce output structure at the token level via regex/CFG/FSM so invalid output is *impossible* (not retried) | Outlines, Guidance, LMQL |

The boundary with app-level "structured output" is porous — native provider
Structured Outputs *are* constrained decoding exposed as an API, and most typed-
function libraries route to it when available.

## 1. Compiled prompting — DSPy

Three abstractions: **signatures** (a natural-language typed I/O spec, e.g.
`question -> answer`, describing the task not the prompt), **modules** (wrap a
signature with a prompting strategy — `Predict`, `ChainOfThought`,
`ProgramOfThought`, `ReAct`, `MultiChainComparison` — so the reasoning strategy
is a swappable code object), and **optimizers** (formerly "teleprompters"). The
**compile step** takes the program + training examples + a metric and runs an
optimizer that generates and evaluates prompt variants — selecting few-shot
demos, rewriting instructions, tuning structure. DSPy `Assert`/`Suggest` add
constraint-and-self-correct.

> **Optimizer algorithms (owned elsewhere):** BootstrapFewShot, MIPROv2
> (joint instruction + demo search), COPRO, and **GEPA** (reflective Pareto
> evolution) are documented in `prompt-deep-optimizer` / `prompt-helper-optimizer`
> — this reference covers DSPy as a *programming framework*; defer the algorithm
> internals there.

**Best fit:** structured tasks (QA, classification, extraction, multi-hop)
**where you have eval data**; reported 10–40% gains on such tasks. *Not* for
one-shot tasks or pure orchestration.

## 2. Schema-first typed LLM functions

"Every prompt is a function that takes parameters and returns a type." Declaring
the return type yields type safety, IDE autocomplete, partial-type streaming, and
reliable structured output even from models without native tool-calling. Reframes
prompt engineering as **schema engineering**.

- **BAML (BoundaryML)** — a Rust-built `.baml` DSL: declare typed functions,
  schemas, prompts, model choice, retry policy; the compiler generates type-safe
  clients for Python/TS/Ruby/Go/Java/C#/Rust. Key innovation **Schema-Aligned
  Parsing (SAP)** tolerantly parses real LLM output (markdown-in-JSON,
  chain-of-thought before the answer, trailing commas) into the declared type
  without native tool-calling, so "structured outputs work on Day 1 of a model
  release." **Test-driven prompts**: `test` blocks with `check`/`assert` run
  against live APIs in a playground or `baml-cli test`. Prompts + model + retries
  live declaratively in version-controlled `.baml` files (vs library
  post-processing in app code).
- **Instructor** — minimal abstraction; stays close to the provider client, wraps
  validate-retry, routes to native structured output then tool-calling. The "safe
  80% default."
- **Pydantic AI** — model-agnostic agent framework from the Pydantic team; good
  default for new agentic projects.
- **Marvin / Mirascope / Ell** — function/decorator prompting. Mirascope
  ("anti-framework") turns plain functions into LLM calls via `@llm.call` +
  `@prompt_template`, provider-agnostic across 20+ providers; Lilypad auto-versions
  every prompt-bearing function. Ell shares the prompts-as-versioned-functions
  lineage.

## 3. Constrained generation / decoding

Rather than ask-then-validate, **make invalid output impossible** by masking
disallowed tokens during decoding. Autoregressive generation is reformulated as
transitions in a **finite-state machine** compiled from a regex/JSON-schema/CFG;
invalid next-token logits are set to `-inf`. Constraint forms: regex, CFG, JSON
Schema.

- **Outlines** — originated the FSM reformulation; compiles to index structures
  for ~O(1) valid-token lookup per step; the **fastest** option for simple
  high-volume tasks (e.g., ticket classification); integrated into serving stacks.
- **Guidance (Microsoft)** — token-by-token control + interleaved program flow;
  **token healing** inserts grammar-determined tokens (e.g., fills `h1>` after
  `</`) to skip forward passes and cut GPU cost.
- **LMQL (ETH)** — SQL-like constraint query language; constraints evaluated
  eagerly per token and compiled to token masks. (Momentum has slowed relative to
  Outlines/Guidance — low/medium confidence.)

> **Serving boundary:** the *runtime* integration of constrained decoding
> (XGrammar/Outlines/llguidance inside vLLM/SGLang) is in `llm-inference-serving`;
> the *framework/API* view is here.

## 4. App-level structured output (the substrate)

Three mechanisms: **(1)** native provider Structured Outputs with `strict: true`
(constrained-decoding guarantees); **(2)** constrained decoding (logit masking);
**(3)** prompt + validate + reask (Instructor's loop — re-prompt with the
validation error attached until valid or a retry cap). Plain **JSON mode** (valid
JSON, no schema enforcement) is considered obsolete in production since mid-2025,
superseded by schema-enforcing Structured Outputs. **Function-calling-as-
extraction** = provide a tool schema and use the model's filled arguments as your
typed data.

## When declarative wins (decision)

- **Declarative/compiled wins when you have eval data + repeated invocation** —
  the optimization payoff scales with how often the prompt runs and how measurable
  quality is.
- **Hybrid is the mature default** — an orchestrator (LangGraph) for
  routing/memory/tools, with DSPy *inside* a classifier/reranker/planner that
  benefits from optimization. "Keep each framework in its lane."
- **Constrained decoding for hard guarantees** (you control the decoder);
  **validate+retry for portability** (closed APIs without decoder access).
- **Typed-function libraries for portability + testability** — switch providers
  by one parameter; BAML/SAP makes new models work Day 1.

## Anti-patterns

- **Opaque compiled prompts** — auto-generated; understanding *why* a demo was
  chosen needs the compilation trace; DSPy lacks native tool-call logs.
- **Hidden token/compile cost** — a MIPRO run on 200 examples can cost $5–10 /
  thousands of API calls; justified only when it saves more manual prompt-fiddling.
- **Over-abstraction / "magic"** — frameworks that hide the exact levers (prompt
  shape, tool wiring, latency) teams need.
- **Framework lock-in / DSL learning cost** — BAML/LMQL add a non-Python language.
- **Reaching for compilation on one-shot tasks** — plain prompting wins.
- **Plain JSON mode in production** when strict Structured Outputs exist.

## Cross-references

- **Prompt-optimization algorithms** (MIPROv2 / BootstrapFewShot / COPRO / GEPA /
  APE / OPRO / TextGrad) → `prompt-deep-optimizer`, `prompt-helper-optimizer`.
- **Hand-writing a single prompt's wording / technique** → `prompt-engineering`.
- **Serving-level constrained decoding runtime** (XGrammar/Outlines in
  vLLM/SGLang) → `llm-inference-serving`.
- **Agent orchestration / routing / memory** → `agent-ecosystem`,
  `multi-agent-orchestration`.

## References

DSPy site + paper (arXiv:2310.03714) + optimizer docs; GEPA (2507.19457); "Prompts
Are Programs Too!" (2409.12447); PDL (2410.19135). BAML (BoundaryML GitHub + docs,
SAP, testing). LMQL (lmql.ai); Guidance (Microsoft Research + GitHub); Outlines
(dottxt; LMSYS compressed-FSM); "Guiding LLMs the Right Way" (2403.06988).
Instructor / Pydantic AI / Mirascope / Marvin docs; Simmering structured-output
comparison; Agenta + BuildMVPFast structured-output/function-calling guides.
*(27 sources, 2024–2026; full URLs in the source research report.)*
