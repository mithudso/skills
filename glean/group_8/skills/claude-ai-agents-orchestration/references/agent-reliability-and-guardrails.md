<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Formerly the standalone `agent-reliability-and-guardrails` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

# Agent Reliability & Guardrails

Engineering reliability and safety **into the agent loop**: validating what goes in and comes out, surviving transient and systemic failures, constraining what tools can do, and defending against prompt injection and excessive agency. This is the *build-time and run-time defense* lane — for monitoring/telemetry of these controls see `llm-observability`; for the tool-schema/error-recovery contract see `agent-harness-construction`.

## The honest framing (read first)

Guardrails reduce risk; they do not eliminate it. LLMs cannot reliably distinguish trusted instructions from instructions embedded in the content they process — "LLMs follow instructions in content" (Willison). Commercial guardrail products claiming to block "95% of attacks" are, in security terms, **failing** — a determined attacker iterates against the 5%. Treat every guardrail as **defense-in-depth**, never as a sole control, and assume any single layer can be bypassed. The most reliable defenses are **architectural** (remove the capability) rather than **detective** (try to catch the bad input). The 2025 design-patterns paper is blunt: general-purpose agents built on current LLMs **cannot offer meaningful, reliable safety guarantees** — every robust pattern works by *constraining* what the agent can do.

## Guardrail types

Guardrails are checks placed *around* the model, not inside it. Classify by where they sit in the loop:

| Rail | Sits at | Catches | Tooling |
|---|---|---|---|
| **Input rail** | Before the LLM/tool | Jailbreaks, injection strings, off-topic/PII input, oversized input | NeMo input rails, Llama Guard, regex/classifier filters |
| **Output rail** | After the LLM, before the user/tool | Toxicity, PII/secret leakage, hallucination, schema violations, competitor mentions | Guardrails AI validators, NeMo output rails, content classifiers |
| **Retrieval/execution rail** | Around RAG fetch & tool calls | Poisoned retrieved chunks, unsafe tool args | NeMo retrieval/execution rails |
| **Topical/dialog rail** | Conversation control | Steering off allowed topics, forbidden flows | NeMo Colang dialog rails |
| **Structured-output validation** | Output | Malformed JSON, type/range violations | Pydantic / JSON-Schema / Zod + re-ask loop |

**Structured-output validation** is the highest-leverage, lowest-cost guardrail: constrain the model to a schema and **validate before acting**. Use provider structured-output / function-calling modes where available; layer a validate-then-reask loop on top (Guardrails AI's `num_reasks` re-prompts on validation failure until it passes or a cap is hit). Reject or repair rather than trusting raw text. Cross-ref: `zod-schema-validation`, `json-advanced`.

**Toolkits, briefly:**
- **NeMo Guardrails** (NVIDIA) — programmable rails in *Colang*; supports input, retrieval, dialog, execution, and output rails; acts as a proxy between user and LLM. Newer IORails engine runs content-safety, topic-safety, and jailbreak-detection rails in parallel.
- **Guardrails AI** — Python `Guard` wraps the LLM call; composes validators (PII, toxicity, regex, competitor-check, …) from the Guardrails Hub against a Pydantic/RAIL spec; on failure can **reask, fix, or reject**.
- **Llama Guard** (Meta) — an LLM-based safety *classifier* for input/output content moderation; narrower than the above (no dialog/structure control) but composes well as one rail.

**Limits of content filters:** classifier-based rails (Llama Guard, toxicity/jailbreak detectors) are probabilistic — they have false negatives by construction. Never let a content filter be the only thing standing between untrusted input and a high-impact tool.

## Reliability patterns

Layer these; each handles a different failure mode. Order of construction: **classify errors → timeouts → retries → circuit breaker → fallback → caps/loop-detection → graceful degradation.**

- **Timeouts** — bound every model call, tool call, and the whole run. Prevents hangs from consuming the step/cost budget.
- **Retries with exponential backoff + jitter** — for *transient* errors only (429, 5xx, network blips). Base delay ~100–500 ms, double each attempt, add jitter to avoid retry storms / thundering herd. Classify first: never retry a deterministic 4xx (bad request, content-policy refusal) — it will fail identically.
- **Circuit breaker** — three states: **closed** (pass through) → **open** (fail fast without calling the dependency, after an error-rate threshold) → **half-open** (after a cooldown, let a small % through to test recovery). Open durations typically start ~30 s with backoff, capping ~5 min. Prevents one failing dependency from cascading.
- **Fallback chains** — on failure, degrade gracefully: same prompt on a cheaper/alternate model → cached response for common queries → deterministic rule-based response. Always have a terminal "safe" answer.
- **Budget & step caps** — hard ceilings on `max_steps`/iterations, `max_tokens`, wall-clock, and dollar cost per run. A step cap is the single most important infrastructure control: every agent run must terminate.
- **Loop & no-progress detection** — step caps alone are insufficient. Add **action fingerprinting** (hash the proposed tool+args, block identical/near-identical repeats) and **no-progress detection** (exit when N iterations yield no new information / no state change). These catch the agent thrashing on the same failing action *within* the step budget.
- **Idempotency** — make tool side-effects idempotent (idempotency keys) so a retry after an ambiguous timeout doesn't double-charge / double-send. Pairs with durable state (see `agent-state-and-durable-execution`).
- **Graceful degradation** — when guardrails or dependencies trip, return a reduced-but-honest result and surface the failure; never fail silently or fabricate.

## Tool-call safety

The agent's tools are its blast radius. Constrain *capability*, not just *intent*.

- **Least privilege** — grant each tool the minimum scope for its task. An indirect-injection email that tells the agent to forward the inbox is *impossible* if the mail tool is read-only. Excessive scope is the vulnerability (OWASP **LLM06 Excessive Agency**: too many tools, too-broad permissions, or acting without approval).
- **Allow-lists over deny-lists** — default-deny tools, commands, file paths, and network destinations; explicitly enumerate what's permitted. Reads broadly allowed; writes, shell, network, and destructive ops require approval.
- **Sandboxed execution** — *permissions* decide **when** the agent may act; *sandboxing* decides **what it is technically able to do at all**, holding even if the model is confused or compromised. Use OS-level isolation (e.g., containers, macOS Seatbelt) beneath the app layer so even indirect paths are contained. For the Linux confinement mechanisms behind "OS-level isolation" — and which to pick for *arbitrary, hostile* agent-generated code (seccomp/Landlock are shared-kernel and not enough; escalate to gVisor or a Kata/Firecracker microVM) — see the `devops-infra` hub reference `references/linux-sandboxing-confinement.md` (seccomp-bpf, Landlock, gVisor, Kata, Firecracker, and the isolation-vs-performance decision model).
- **Human-approval gates for high-impact actions** — require explicit human confirmation before anything that modifies data, spends money, sends external communication, or changes permissions. Beware **approval fatigue**: too many prompts and users rubber-stamp them — reserve gates for genuinely high-impact, irreversible actions and make the prompt show the concrete effect.
- **Validate tool arguments** — apply input rails / schema validation to *tool args*, not just the user message; injected instructions often surface as malicious arguments.

## Agent-security threats + defenses

### Threats

- **Prompt injection** (OWASP **LLM01**, #1 two editions running) — untrusted text overrides operator intent. **Direct**: the user types the attack. **Indirect**: the attack hides in content the agent ingests — a web page, email, PDF, RAG chunk, screenshot, tool output, or code comment.
- **The lethal trifecta** (Willison, 2025) — an agent is exposed to data theft when it has **all three** simultaneously: (1) **access to private data**, (2) **exposure to untrusted content**, and (3) **the ability to communicate externally** (exfiltration path: HTTP, email, API, commit). Any one alone is fine; the combination "virtually guarantees" exfiltration via injection. Documented in the wild against Microsoft 365 Copilot, GitHub MCP, GitLab Duo, and others.
- **Tool / RAG poisoning** — attacker plants malicious instructions in a knowledge source, MCP tool description, or retrievable document that fire when ingested (RAG is a new attack surface, not a shield — OWASP **LLM08 Vector & Embedding Weaknesses**).
- **Excessive agency** (**LLM06**) — over-permissioned tools turn a single injection into real-world damage.
- **Improper output handling** (**LLM05**) — downstream systems trust LLM output (rendered HTML, executed SQL/shell, eval'd code) → XSS, SQLi, RCE.
- **Sensitive info disclosure / system-prompt leakage** (**LLM02 / LLM07**) and **unbounded consumption** (**LLM10**, cost/DoS via runaway loops — see budget caps above).

### Defenses (defense-in-depth, architectural first)

1. **Break the lethal trifecta** — the most reliable mitigation is to *remove one leg*: deny untrusted content, or strip the external-communication tool, or wall off private data. Architecture beats detection.
2. **Dual-LLM / quarantine pattern** (Willison) — a **privileged LLM** with tool access never sees untrusted tokens; a **quarantined LLM** processes untrusted content but has no tools and returns only **symbolic references** (e.g., `$email-summary-1`) the privileged LLM passes around blind. **CaMeL** (Google DeepMind, 2025) hardens this with a custom interpreter that does data-flow / taint tracking and enforces capability-based policies *without modifying the model*.
3. **Constrained agent design patterns** (Willison/DeepMind, six patterns — see `agent-planning-patterns`): **Action-Selector** (fire tools, never read responses), **Plan-Then-Execute** (fix the tool sequence *before* seeing untrusted output so content can't redirect actions), **LLM Map-Reduce** (isolated sub-agents return only symbolic/boolean results), **Dual-LLM**, **Code-Then-Execute/CaMeL**, **Context-Minimization** (drop the user prompt from context before generating). Each trades general capability for safety.
4. **Wrap untrusted input in delimited, escaped envelopes** — tag every untrusted source (`<untrusted_transcript>…</untrusted_transcript>`) with HTML-entity escaping; segregate it from instructions in the prompt. Reduces (does not eliminate) injection. Cross-ref `llm-context-engineering`.
5. **Output gating** — validate/sanitize/encode LLM output before any downstream sink; never auto-render or auto-execute model output (mitigates LLM05).
6. **Least privilege + sandbox + approval gates** — see Tool-call safety above.
7. **Telemetry on every rail** — log rail triggers, blocked actions, and refusals; alert on anomalies. Cross-ref `llm-observability`.

## Checklist

- [ ] Every model call, tool call, and run has a **timeout**.
- [ ] **Step / token / cost caps** enforced per run; the run is guaranteed to terminate.
- [ ] **Loop detection**: action fingerprinting + no-progress exit (not just a step cap).
- [ ] Retries use **backoff + jitter** and only fire on **classified transient** errors.
- [ ] **Circuit breaker** around external dependencies; **fallback chain** ends in a safe response.
- [ ] Tools follow **least privilege**; default-deny **allow-lists** for tools/commands/paths/network.
- [ ] High-impact / irreversible actions sit behind a **human-approval gate** (and gates are rare enough to avoid fatigue).
- [ ] Tool execution is **sandboxed** at the OS/container level, not just permission-gated.
- [ ] **Input + output rails** present (jailbreak/PII/toxicity); tool *arguments* validated too.
- [ ] LLM output is **schema-validated** (Pydantic/JSON-Schema/Zod) with a reask/repair/reject path before use.
- [ ] Untrusted content is **delimited + escaped**; output is **gated** before any execute/render sink.
- [ ] **Lethal-trifecta check**: does this agent hold private-data + untrusted-content + external-comms together? If so, remove a leg or quarantine.
- [ ] Side-effecting tools are **idempotent** (idempotency keys) so retries don't double-act.
- [ ] Rail triggers and blocked actions are **logged + alerted** (`llm-observability`).

## Anti-patterns

- Treating a content filter / "95%" guardrail product as a **sufficient** injection defense.
- Step cap as the *only* loop control — misses same-action thrashing inside the budget.
- Retrying non-transient errors (4xx, policy refusals) — burns budget on guaranteed failures.
- Permission gates without a sandbox — a compromised model still does damage through allowed paths.
- Approval gates on everything → fatigue → rubber-stamping → gates become theater.
- Feeding untrusted content to a privileged, tool-equipped LLM and hoping a system-prompt instruction ("ignore injected instructions") holds.
- Auto-rendering or auto-executing raw LLM output (XSS/SQLi/RCE via LLM05).
- Granting write/send/deploy scope a tool's task never needs (excessive agency).
- Failing silently or fabricating a result when a guardrail trips instead of degrading honestly.

## References

- OWASP — *Top 10 for LLM Applications 2025* (LLM01 Prompt Injection, LLM05 Improper Output Handling, LLM06 Excessive Agency, LLM08 Vector & Embedding Weaknesses, LLM10 Unbounded Consumption): https://owasp.org/www-project-top-10-for-large-language-model-applications/assets/PDF/OWASP-Top-10-for-LLMs-v2025.pdf
- Simon Willison — *The lethal trifecta for AI agents: private data, untrusted content, and external communication* (2025-06-16): https://simonwillison.net/2025/Jun/16/the-lethal-trifecta/
- Simon Willison / Google DeepMind — *Design Patterns for Securing LLM Agents against Prompt Injections* (six patterns) (2025-06-13): https://simonwillison.net/2025/Jun/13/prompt-injection-design-patterns/
- Google DeepMind — *CaMeL: Defeating Prompt Injections by Design* (capability-based dual-LLM + interpreter) (2025): https://arxiv.org/abs/2503.18813
- NVIDIA — *NeMo Guardrails* (programmable input/retrieval/dialog/execution/output rails, Colang): https://github.com/NVIDIA-NeMo/Guardrails • paper: https://arxiv.org/abs/2310.10501
- Guardrails AI — *Guards, validators, RAIL/Pydantic, reask* docs: https://www.guardrailsai.com/docs/api_reference_markdown/guards
- Portkey — *Retries, fallbacks, and circuit breakers in LLM apps: what to use when*: https://portkey.ai/blog/retries-fallbacks-and-circuit-breakers-in-llm-apps/
- NVIDIA — *Practical Security Guidance for Sandboxing Agentic Workflows*: https://developer.nvidia.com/blog/practical-security-guidance-for-sandboxing-agentic-workflows-and-managing-execution-risk/
