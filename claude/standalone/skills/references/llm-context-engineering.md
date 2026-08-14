<!-- hub-reference-banner -->
> **Reference file — part of the `ai-mcp-sdk-prompting` hub.** Formerly the standalone `llm-context-engineering` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: llm-context-engineering
title: LLM Context Engineering
description: >
  Master LLM context design, optimization, and security for AI applications
  and agents. TRIGGER: writing or refactoring a CLAUDE.md, AGENTS.md,
  .cursorrules, or skill file; deciding what goes in a system prompt vs. a
  RAG chunk vs. agent memory; implementing Anthropic prompt caching
  (cache_control breakpoints) or cutting API token costs; building or
  debugging a RAG pipeline or vector-search corpus injection; architecting
  agent memory (short-term, episodic, semantic, procedural); switching
  context safely between customer accounts; hardening LLM inputs against
  prompt injection or PII leakage; diagnosing context rot where model quality
  degrades as context grows; any question about context windows, token
  budgets, context staleness, cache TTLs, or static vs. dynamic vs. ephemeral
  context. SKIP: one-off prompt wording with no architecture concern;
  infrastructure or deployment questions; database schema or query design;
  model selection or provider benchmarking.
origin: local
category: developer
version: "1.5"
updated: "2026-05-31"
tags:
  - llm
  - context
  - prompt-engineering
  - rag
  - caching
  - security
  - agent-memory
  - token-optimization
whenToUse:
  - "writing or refactoring CLAUDE.md, AGENTS.md, .cursorrules, or a skill file"
  - "deciding what goes in a system prompt vs. RAG chunk vs. agent memory"
  - "implementing Anthropic prompt caching or cache_control breakpoints"
  - "reducing API token costs for a Claude application"
  - "building or debugging a RAG pipeline or vector-search corpus injection"
  - "architecting agent memory (short-term, episodic, semantic, procedural)"
  - "switching context safely between customer accounts in TAM or support workflows"
  - "hardening LLM inputs against prompt injection or PII leakage"
  - "diagnosing context rot — model quality degrading as context grows"
  - "questions about context windows, token budgets, or cache TTLs"
whenNotToUse:
  - "one-off prompt wording with no architecture concern — use prompt-deep-optimizer"
  - "infrastructure or deployment (servers, containers, CI) — use backend-patterns"
  - "database schema or query design — use a database expert skill"
  - "model selection or provider benchmarking — use llm-models"
related_skills: [prompt-deep-optimizer, ai-rag-retrieval, ai-agent-engineering]
---

# LLM Context Engineering

## When NOT to Use

- One-off prompt wording with no architecture concern — use `prompt-deep-optimizer`
- Infrastructure or deployment (servers, containers, CI) — use `backend-patterns`
- Database schema or query design — use a database expert skill
- Model selection or provider benchmarking — use `llm-models`

## Context Type Decision Table

Pick the right mechanism before writing any context:

| What you need to express | Context type | Primary mechanism |
|--------------------------|--------------|-------------------|
| Stable project rules, coding conventions | Static / persistent | CLAUDE.md, AGENTS.md |
| Account-specific facts for one session | Dynamic / session | Injected system prompt block |
| Documents retrieved per query | Dynamic / ephemeral | RAG chunk injection |
| Conversation history across turns | Episodic memory | Summarized turn history |
| Cross-session learned preferences | Semantic memory | Vector DB + retrieval |
| Reusable task procedures | Procedural context | Skill files / tool descriptions |
| Credentials or tokens valid only in this request | Ephemeral only | In-turn injection, never persisted |

## Token Budget Allocation

Scale proportionally to your actual context window. Cached layers (rows 1–3) should
stay under 10% of the window combined — anything beyond that is dead weight; prune or
move to RAG.

| Layer | % of window | Notes |
|-------|-------------|-------|
| System prompt (static) | 1–4% | Cache with `cache_control`; always first |
| Skill / tool context | 0.5–2% | Cache alongside system prompt |
| Few-shot examples | 0.5–1.5% | Cache; 2–5 examples is optimal |
| Retrieved context (RAG) | 5–20% | Dynamic; insert before conversation history |
| Conversation history | 2–10% | Prune aggressively; summarize old turns |
| Current user turn | 0.25–1% | Always last; never cache this layer |
| Reserved output budget | 2–8% | Set `max_tokens` explicitly; do not leave it unset |

## Caching Strategy Decision Table

| Scenario | TTL | Why |
|----------|-----|-----|
| System prompt, static rules | 1-hour | Stable; reused across users/sessions |
| Skill / reference docs | 1-hour | Low churn; high reuse value |
| Per-session user / account context | 5-min | Session-scoped; not safe to share across users |
| Per-request RAG chunks | None | Too variable; cache hit rate near zero |
| Few-shot examples | 1-hour | Fixed content; pure cache win |

The 5-min TTL pays back on the second request in a session. Use 1-hour only when you
know the cached prefix will be reused within a 10-minute window.

## How to Structure a CLAUDE.md or System Prompt

Order content from most-stable to least-stable. This ordering keeps cache breakpoints
efficient — the stable prefix is written once and read cheaply on every subsequent call.

1. Role and capabilities — changes almost never
2. Project or domain rules — changes with major releases
3. Few-shot examples — changes rarely; last cacheable layer
4. Dynamic account/session context — changes per session; do not cache
5. Current turn / user message — always last

Place an Anthropic `cache_control` breakpoint on the content block immediately after
layer 3. The correct API shape is:

```json
{
  "role": "user",
  "content": [
    {
      "type": "text",
      "text": "<your static prefix here>",
      "cache_control": {"type": "ephemeral"}
    }
  ]
}
```

Never place a `cache_control` breakpoint inside dynamic layers (4 or 5) — it defeats
the purpose and adds cache-write cost with no benefit.

## Security Checklist Before Any Context Write

Run through this before injecting any external data into a prompt:

| Check | Rule |
|---|---|
| Untrusted inputs | Wrap in `<untrusted_input>` tags and HTML-entity escape before injection |
| PII in shared layers | No names, emails, or account IDs in cached or shared context |
| Injection resistance | System prompt instructs the model to ignore instructions in retrieved documents |
| Output validation | Validate output before passing to downstream tools (shell, SQL, API calls) |
| Minimal access | Model receives only context needed for this request — no ambient cross-account data |

## Agentic Context Engineering (2025–2026)

For multi-step agents, context engineering is "effectively the #1 job" (Cognition).
The discipline framing: the LLM is an operating system, the **context window is
its RAM** under a finite **attention budget** (Karpathy); the goal is the
*smallest set of high-signal tokens* that yields the behavior, not the biggest
dump. LangChain's operational frame is **Write / Select / Compress / Isolate**.

**Why filling the window backfires (the empirical bedrock):**
- **Context rot** — every frontier model degrades as tokens accumulate, *before*
  the window overflows (a 200K model can degrade meaningfully by ~50K).
- **Lost-in-the-middle** — U-shaped attention; the model uses the start and end
  far better than the middle. Put load-bearing tokens at the edges.
- **Effective length ≪ advertised** — validate with **RULER / NoLiMa**-style
  evals (which strip literal lexical cues), not vanilla needle-in-a-haystack.

**The toolkit:**
- **Compaction** — summarize history near a threshold (Claude Code auto-compacts
  at ~95%); distinguish server-side compaction from client-side **context
  editing** (clearing stale tool calls/results — Claude `context-management`
  beta; ~84% token cut in Anthropic's internal eval).
- **Just-in-time retrieval** — hold lightweight references (paths, queries, links)
  and load at runtime via tools; prefer agentic search to pre-loaded RAG dumps.
  (Retrieval-loop mechanics → `ai-agent-engineering` (references/iterative-retrieval.md); advanced RAG →
  `ai-rag-retrieval` (references/advanced-rag-patterns.md).)
- **External memory & note-taking** — scratchpads, **filesystem-as-memory** with
  *recoverable compression* (store a pointer, drop the body, re-fetch on demand),
  the Claude memory tool, recitation (re-write the goal/todo into recent context).
  (Long-term memory-store design → `ai-agents-orchestration` (references/agent-memory-architecture.md).)
- **KV-cache-aware layout** — agent input:output runs ~100:1, so cache hit rate is
  the top cost lever (~10× token-cost swing). Rules: **stable prefix** (no
  timestamps/volatile tokens in the system prompt), **append-only** context, and
  **mask — don't remove — tools** mid-run (removing them invalidates the cache and
  can break schemas). Complements the Caching Strategy table above.

**Multi-agent context (contested).** Cognition: *share full agent traces, not just
messages* — isolated sub-agents make conflicting implicit decisions, so default
single-threaded. Anthropic: isolated sub-agent windows returning *condensed*
findings to an orchestrator beat single-agent by 90% on parallel research (at
~15× tokens). Reconciliation: **isolate for read-heavy parallel research; share
context for tasks needing shared design decisions.**

**The four failure modes (Breunig).** **Poisoning** (a hallucination enters
context and is repeatedly referenced), **distraction** (over-long context makes
the model over-focus on history vs training knowledge), **confusion** (superfluous
content / too many tools degrades output — curate the tool set), **clash** (newly
accrued info conflicts with earlier context; a single early wrong turn can cause
large multi-turn accuracy drops). Counter-tension: keep *informative* failures
(they aid recovery), purge *propagating* falsehoods.

## Reference

Read the reference file for code examples, anti-patterns, RAG injection strategies,
agent memory architecture, context freshness patterns, and cited sources.

Section guide:
- Context window mechanics and types → Section 1
- Writing CLAUDE.md / AGENTS.md / skill files → Section 2
- Token budgeting and prompt caching code examples → Section 3
- RAG pipeline and chunk injection → Section 4
- Agent memory systems and architecture → Section 5
- Staleness detection and context versioning → Section 6
- Prompt injection, PII, context poisoning defenses → Section 7
- Anti-patterns to avoid → Section 8
- Sources → Section 9

`~/.claude/skills/llm-context-engineering/references/llm-context-engineering-context.md`
