# Prompt-deep-optimizer audit passes (A–P)

Authoritative per-pass definitions for the 16 evaluation domains, grouped into the 5 parallel-dispatch bundles. Extracted from `SKILL.md` (Step 2) so the body stays under the Pass J token budget. The Step 2 dispatch reads this file to build each subagent bundle; each subagent receives **only its group's passes**. The grouping rationale, dispatch rules, small-artifact profile, subagent-budget rules, and the skip-protocol table stay inline in `SKILL.md` Step 2 — they are control flow, not reference.

Group → passes: **Group 1 — Intent & Output** (A, D, M); **Group 2 — Context & Inputs** (B, C, N); **Group 3 — Process & Tools** (E, F, G); **Group 4 — Safety & Robustness** (H, I, O); **Group 5 — Structure, Model & Algorithm** (J, K, L, P).

## Group 1 — Intent & Output (passes A, D, M)

**Pass A — Intent & framing**
- Goal clarity: is the desired outcome stated, or only implied?
- Success criteria: how would the model know it's done well?
- Audience/persona: stated, and at the right expertise level?
- Scope boundaries: what's explicitly out of scope?

**Pass D — Output contract**
- Format specified (markdown, JSON schema, code, table)
- Length constraints stated
- Required vs optional sections distinguished
- Tone, style, voice, reading level set
- "I don't know" / refusal behavior defined

**Pass M — Meta**
- Versioning / ownership / changelog marker present (comment or header)
- Reusable template variables extracted
- Composability: do the output format and slot names align with downstream prompt inputs? Flag as Medium if the output schema cannot be consumed by a sibling prompt without transformation.

## Group 2 — Context & Inputs (passes B, C, N)

**Pass B — Context & grounding**
- Required background, definitions, domain vocabulary present?
- Prior decisions, constraints, or invariants stated?
- Source-of-truth references (files, URLs, schemas, examples) attached?
- Explicit "do not assume / do not invent" guardrails?

**Pass C — Inputs**
- Variables/placeholders clearly marked and named consistently
- Input format defined (JSON, free text, structured fields)
- Edge cases addressed (empty, malformed, ambiguous, multilingual)
- Trust level of each input distinguished (user vs system vs untrusted)

**Pass N — Variable templating & composition**
- Dynamic slots use a consistent templating standard
- **Conditional sections must be marked in the templating idiom already in use by the surrounding prompt:** Jinja2 (`{% if %}` / `{% endif %}`), Handlebars (`{{#if}}` / `{{/if}}`), Python f-string (conditional must wrap in helper function), or XML comment marker (`<!-- COND: ... -->` then `<!-- /COND -->`). Flag as Medium if conditions are described in prose but no idiom is established.
- No hardcoded values that should be variables (model names, thresholds, collection names)

## Group 3 — Process & Tools (passes E, F, G)

**Pass E — Reasoning & process**
- Whether to think step-by-step or answer directly stated
- Decomposition strategy for multi-step tasks
- **Self-check / verification pass requested when the output type matches one of: runnable code, JSON / YAML, numeric computation, SQL, structured schema, regex.** (Removed the "OR success criteria" prong because Pass A adds success criteria, which would create a feedback loop.)
- Clarifying-question policy. **Default if none stated:** "If input is ambiguous, ask exactly one targeted question before proceeding." Applying this default does NOT trigger BLOCKED; it is an established fallback the auditor may insert verbatim.

**Pass F — Tools & capabilities** (agentic prompts only — mark `N/A` if no tools)
- Allowed/forbidden tools listed
- Tool selection heuristics provided
- Parallel vs sequential call guidance
- Budgets (tokens, calls, time) bounded
- Stop conditions explicit

**Pass G — Examples / few-shot**
- Positive examples cover the common path
- Negative examples ("don't do this") present where failure modes exist
- Edge-case examples included
- Examples diverse rather than redundant
- Examples placed *after* the instruction they illustrate

## Group 4 — Safety & Robustness (passes H, I, O)

**Pass H — Constraints & guardrails**
- Hard prohibitions stated (PII, secrets, destructive ops)
- Safety / policy compliance addressed
- **Citation / sourcing requirements specified — flag as Medium ONLY when the prompt's task type produces factual claims (research summaries, documentation, reports, analytical writeups). For code generation, data transformation, format conversion, and reformatting tasks, mark `N/A (task produces no factual claims)`.**
- Determinism vs creativity expectations set

**Pass I — Robustness**
- Prompt-injection resistance (treat untrusted input as data, not instructions)
- Behavior under conflicting instructions defined
- Behavior when context is missing or contradictory
- Handling of long/truncated input
- No contradictory rules within the prompt itself
- No references to undefined variables, tools, or files

**Pass O — Auto-healing & resilience**
- Output validation defined: schema validators or format checks for structural failure detection
- Retry strategy specified: what to do when output fails validation
- **Fallback chain defined.** Standard pattern when the audited prompt lacks one — the auditor may insert this verbatim as the proposed fix without it counting as inventing content:
  > `If output fails validation, respond with: "I was unable to produce a valid <format>. Here is my best partial attempt: <attempt>. Please review."`
- Circuit breaker defined: maximum retry count to prevent infinite loops

## Group 5 — Structure, Model & Algorithm (passes J, K, L, P)

**Pass J — Structure & ergonomics**
- Section ordering follows persona → context → task → constraints → output format
- Headings, delimiters, or XML tags used consistently
- Critical instructions repeated at start *and* end when prompt > ~2,000 tokens AND the repeated content is a hard constraint (not a stylistic note)
- Dead weight removed: filler, redundant restatements, mergeable sections (except where the repetition rule above applies)
- Bullets vs prose chosen correctly (bullets for parallel/scannable, prose for causal/sequential)

**Pass K — Model fit**
- Prompt accounts for the target model's known strengths/limits
- **Stable prefix** structured for prompt-caching (cacheable content first, volatile last) — actionable for Claude (Anthropic `cache_control`), partial for OpenAI (automatic system+messages caching since GPT-4-turbo+), N/A for models with no prompt cache
- **Thinking / extended-reasoning budgets** set when supported:
  - **Claude 3.7+:** `thinking: { type: "enabled", budget_tokens: N }` in the API request
  - **OpenAI o1 / o3 / o4-mini:** `reasoning_effort: "low" | "medium" | "high"`
  - **Gemini 2.0+ thinking models:** `thinking_config: { thinking_budget: N }`
  - **Other models:** mark `N/A unless model supports extended thinking`
- Cost vs quality trade-offs explicit

**Pass L — Evaluation hooks**
- Test cases or golden outputs referenced
- Known failure modes called out
- Metrics defined (accuracy, latency, refusal rate)
- A/B comparison anchors present when the prompt has variants

**Pass P — Algorithm & pipeline fit**
- Training data availability: none / <30 / 30–100 / >100 paired examples?
- Evaluation metric: defined or undefined?
- Pipeline position: standalone or multi-step pipeline?
- Few-shot example quality: bootstrapped from traces or hand-written?
- **Pipeline-level cache reuse** (distinct from Pass K's within-prompt stable-prefix check): if this prompt is called from multiple pipeline steps with the same stable prefix, is a shared cache key or prompt-caching strategy declared at the pipeline level?
