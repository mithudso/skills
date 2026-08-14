<!-- hub-reference-banner -->
> **Reference file — part of the `ai-mcp-sdk-prompting` hub.** Formerly the standalone `llm-integration-reviewer` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: llm-integration-reviewer
title: LLM Integration Reviewer
description: |
  Practical review reference for multi-provider LLM integrations covering failover, rate limiting, prompt caching, structured-output robustness, privacy/data boundaries, and prompt-injection resilience.
  TRIGGER: reviewing LLM provider failover or retry logic; prompt caching strategy audit (Anthropic or Gemini); structured output / JSON parsing robustness; provider rate-limit handling; prompt-injection boundary review in an LLM pipeline; reviewing what user/page content enters prompts; reviewing model lifecycle staleness; auditing LLM keys/config in extension contexts; reviewing service-worker-based multi-provider LLM code.
  SKIP: single-provider LLM integration with no failover concern (use anthropic-sdk or claude-api skills for pure Anthropic work); ML model training or fine-tuning (use aws-ai-ml); general Chrome extension security audit unrelated to LLM (use chrome-extension-security-reviewer).
category: developer
version: "1.1.0"
updated: "2026-05-29"
keywords:
  - LLM integration
  - multi-provider
  - failover
  - rate limiting
  - prompt caching
  - structured output
  - JSON parsing
  - prompt injection
  - Anthropic
  - Gemini
  - OpenAI
  - provider failover
  - model lifecycle
  - privacy boundary
  - extension LLM
when_to_use:
  - "review LLM provider failover logic"
  - "prompt caching audit"
  - "structured output parsing robustness"
  - "rate limit handling in LLM integration"
  - "prompt injection boundary review"
  - "model lifecycle staleness check"
  - "LLM keys in Chrome extension"
  - "multi-provider LLM pipeline review"
  - "Anthropic vs Gemini failover"
  - "JSON mode parser fallback"
related_skills:
  - anthropic-sdk
  - chrome-extension-security-reviewer
  - llm-context-engineering
  - rag-architecture
origin: local
---

# LLM Integration Reviewer

Practical review reference for multi-provider LLM integrations: failover, rate limiting, prompt caching, structured-output robustness, privacy/data boundaries, and prompt-injection resilience.

## How to use this skill

Start from the bundled context below. Defer to cited official documentation for exact APIs and edge-case behavior. If the request falls outside LLM integration review, choose a more appropriate skill.

**Sources of truth:**
- **Provider docs** — caching, rate limits, retention, structured output
- **OWASP prompt-injection guidance** — LLM-specific security threats
- **Chrome extension security/privacy docs** — extension-side data handling and context boundaries

**Version note:** based on official pages accessed 2026-05-10, framed for this repo's service-worker-based multi-provider LLM pipeline.

---

## Source scope

- **Anthropic caching, rate limits, data-retention:** [Anthropic prompt caching](https://docs.anthropic.com/en/docs/build-with-claude/prompt-caching), [Anthropic rate limits](https://docs.anthropic.com/en/api/rate-limits), [Anthropic API and data retention](https://docs.anthropic.com/en/docs/build-with-claude/api-and-data-retention)
- **Gemini caching, structured outputs, rate limits, model lifecycle:** [Gemini caching](https://ai.google.dev/gemini-api/docs/caching), [Gemini structured output](https://ai.google.dev/gemini-api/docs/structured-output), [Gemini rate limits](https://ai.google.dev/gemini-api/docs/rate-limits), [Gemini deprecations](https://ai.google.dev/gemini-api/docs/deprecations)
- **Prompt-injection and tool-abuse threats:** [OWASP LLM prompt injection prevention](https://cheatsheetseries.owasp.org/cheatsheets/LLM_Prompt_Injection_Prevention_Cheat_Sheet.html)
- **Extension-side privacy and security:** [Chrome user privacy](https://developer.chrome.com/docs/extensions/develop/security-privacy/user-privacy), [Chrome stay secure](https://developer.chrome.com/docs/extensions/develop/security-privacy/stay-secure)
- **OpenAI reliability guidance:** [OpenAI Cookbook reliability techniques](https://github.com/openai/openai-cookbook/blob/65290c6d72dfa8a8493bda0d3fa945ae0f0bdfef/articles/techniques_to_improve_reliability.md)
- **Repo-specific architecture:** `src/background/` LLM pipeline files and `docs/ARCHITECTURE.md`

## Quick review rules

1. **Separate reliability from security.** Failover and retry logic reduce outages; they do not mitigate prompt injection or privacy leakage.
2. **Keep static prompt prefixes cacheable and dynamic content late.** Both Anthropic and Gemini reward stable shared prefixes.
3. **Treat model output as untrusted until validated.** JSON mode improves reliability but callers still need parsing and failure handling.
4. **Review all fallback paths.** A tolerant parser or silent default can preserve UX while masking provider regressions or unsafe output drift.
5. **Bound provider-specific operational assumptions.** Hardcoded RPMs, model IDs, and deprecation dates are operational policies, not guarantees.
6. **Keep sensitive work in trusted extension contexts.** Never leak prompts, provider keys, or privileged output-handling logic into page-facing contexts.

## Review workflow

1. **Map the LLM pipeline.** Identify provider selection, prompt assembly, request dispatch, retry/failover, parsing, caching, and storage/writeback.
2. **Review provider failover and backoff.** Check whether failover is limited to the intended failure classes and whether retries respect provider-specific rate-limit signals.
3. **Review prompt caching and prompt assembly.** Confirm static/shared instructions are front-loaded and cacheable; verify volatile content does not pollute cache efficiency.
4. **Review output contract robustness.** Check whether structured output features are used where possible and whether parser fallbacks are explicit, observable, and safe.
5. **Review privacy and prompt-injection boundaries.** Track what user/page-derived text enters prompts, what secrets/config enter requests, and what output can trigger actions or storage writes.
6. **Review model lifecycle and operational freshness.** Check whether stale model lists, static fallbacks, and deprecation maps can silently drift from provider reality.

## Review surfaces and checks

| Surface | Purpose | Review focus | Caveats |
|---|---|---|---|
| Provider failover policy | Keep features alive under throttling/outage | Failure classes, provider ordering, retry-after handling | Over-broad failover can hide real bugs or duplicate risky requests |
| Provider rate-limit tracking | Avoid 429s and burst collapse | Freshness of hardcoded assumptions, headroom logic, bounded state | Provider quotas are moving targets |
| Prompt caching | Reduce latency/cost for repeated static prefixes | Prefix stability, TTL choice, cache key scope, privacy tradeoffs | Long-lived caches change cost/retention posture |
| Structured output / JSON parsing | Make downstream automation reliable | Schema strength, parser fallback visibility, coercion bugs | "Looks like JSON" is weaker than schema-backed output |
| Prompt assembly | Compose instructions and context | Untrusted text boundaries, prefix ordering, context minimization | More context is not always better — it can hurt safety and caching |
| Output-handling path | Convert model output into app behavior | Safe defaults, observability, no hidden success-shaped fallbacks | Silent fallback can conceal real regressions |
| Provider model updater / static lists | Keep models current | Staleness, API availability assumptions, deprecation behavior | CORS or provider limits can force static fallbacks |

## Standards and best practices

- Prefer **structured output or strong response contracts** over free-form parsing when downstream logic depends on shape.
- Treat **prompt injection as a system-boundary problem**, not just a prompt-writing problem. Keep tools/actions narrow and isolate untrusted content in clearly scoped prompt regions.
- Keep **prompt caching deliberate**: cache large shared prefixes, not volatile user/page text; review retention/cost implications before extending TTLs.
- Keep **provider configuration and keys in trusted extension contexts** and review storage choices against Chrome's privacy guidance.
- For this repo: preserve the documented architecture of prompt-module composition, provider cycling, and deterministic-cache-only response caching (`docs/ARCHITECTURE.md`).

## Known ambiguities

- Provider docs and quotas change faster than most codebases. Treat this skill as review guidance, not a frozen contract.
- Reliability guidance can conflict with privacy/minimization goals: extra retries, wider failover, and longer-lived caches all carry security/privacy costs.
- OpenAI platform docs were not directly fetchable from this environment during source gathering, so OpenAI coverage is intentionally lighter and grounded in official Cookbook material plus repo code.
