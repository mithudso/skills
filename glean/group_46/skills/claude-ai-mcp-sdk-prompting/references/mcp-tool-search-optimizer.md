<!-- hub-reference-banner -->
> **Reference file — part of the `ai-mcp-sdk-prompting` hub.** Formerly the standalone `mcp-tool-search-optimizer` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: mcp-tool-search-optimizer
description: Use when an MCP server or large tool catalog exposes 10+ tools that must be discovered by Anthropic Tool Search / deferred-tool loading — tools not surfacing on the right query, the wrong tool picked from a big set, deferred-tool discovery misses, or before shipping/expanding a large toolset. Audits tool names, descriptions, and argument descriptions for discovery quality and defer_loading posture. SKIP: building the MCP server itself (use mcp-builder / mcp-servers under ai-agent-engineering); writing or auditing a Claude Code skill file (use claude-code-skills / skill-optimizer); optimizing a prompt (use prompt-deep-optimizer).
---

# MCP Tool Search Optimizer

## Overview

Anthropic **Tool Search** discovers tools on demand by matching a query against four fields: **tool name, description, argument names, and argument descriptions**. Deferred tools (`defer_loading: true`) are invisible until a search surfaces them. So for a large toolset, *discoverability is a property of those four text fields* — a perfectly implemented tool that no query matches is a dead tool.

Two failure modes this skill prevents:
1. **Discovery miss** — the user's phrasing never matches the tool's text, so it's never loaded.
2. **Wrong-tool selection** — two tools have near-identical text, so the model can't choose (the "high semantic overlap" anti-pattern; accuracy degrades past 30-50 simultaneously-loaded tools).

The regex variant (`tool_search_tool_regex`) matches Python `re.search()` patterns and is **case-sensitive by default**; the BM25 variant takes natural language. Consistent name prefixes favor regex; heterogeneous names favor BM25.

## When to use

- A server exposes 10+ tools, or tool definitions exceed ~10K tokens.
- Users report "it didn't find the tool" or "it called the wrong one."
- Before adding tools to an already-large catalog.
- Auditing a catalog for `defer_loading` posture (which tools stay hot).

**When NOT to use:** fewer than ~10 tools (load them all, skip search); authoring the tool's *implementation* logic (use the `mcp-builder` reference under `ai-agent-engineering`); writing a Claude Code *skill* file (use `claude-code-skills`); auditing a Claude Code *skill file's* quality (use `skill-optimizer` — it audits SKILL.md prose; this skill audits MCP tool definitions).

## The audit rubric (7 checks)

| # | Check | Fail signal | Fix |
|---|---|---|---|
| R1 | **Namespace consistency** | mixed/ad-hoc prefixes | Prefix by service/resource (`tam_`, `slack_`); enables clean regex grouping |
| R2 | **Keyword coverage** | description omits the words users actually type | Add task-language synonyms ("stale/outdated/expired", "find/search/lookup") |
| R3 | **Discriminator-first descriptions** | two tools' descriptions are interchangeable | Lead each with what makes it DIFFERENT from its siblings |
| R4 | **Name quality** | vague/noun-only (`query`, `data`) | Verb-first and specific (`search_slack_messages` > `query_slack`) |
| R5 | **Argument searchability** | args named `q`, `x`; no arg descriptions | Name + describe args; search matches these too |
| R6 | **Hot-set sizing** | everything deferred, or >~10 non-deferred | Keep the 3-10 highest-traffic tools non-deferred; defer the rest; never defer the search tool |
| R7 | **Length/token budget** | bloated or one-word descriptions | Tight but keyword-rich; every sentence earns discovery value |

R3 (semantic-overlap) is the highest-leverage and hardest check — see `references/audit-procedure.md` for the clustering method and a full worked example.

## Quick method

1. Pull the full tool list (names + descriptions + arg schemas).
2. Run R1, R4, R5, R7 per-tool (mechanical).
3. Run R3 across the set: cluster tools by shared description keywords; every cluster of 2+ must have discriminator-first descriptions.
4. Run R2 by simulating real user queries and checking each intended tool is matched.
5. Run R6 on the catalog: pick the hot set from usage data (or obvious hot path); recommend the rest deferred.
6. Emit findings with severity (blocking / major / minor) and a concrete rewrite per finding.

See `references/audit-procedure.md` for the detailed procedure, the overlap-clustering algorithm, and the scoring/output contract.

## Common mistakes

- **Describing what the tool IS, not what task it serves.** "Query Slack" tells search nothing; "Search Slack messages by keyword, channel, or date range" matches real queries.
- **Optimizing names but ignoring argument descriptions** — search matches all four fields; thin arg text loses matches.
- **Leaving a whole family interchangeable** (`search_X`, `recommend_X`, `resolve_X` with the same description). The model coin-flips. Discriminator-first fixes it.
- **Deferring the hot path.** If a tool fires on most sessions, deferring it adds a search round-trip every time — keep it non-deferred.
- **Trusting `auto` blindly on a huge context window.** `ENABLE_TOOL_SEARCH=auto` triggers at 10% of context; on a 1M window that's a 100K-token bar — past the 30-50-tool accuracy cliff. Prefer `auto:5` or lower for large catalogs.

## Related

- Tool-definition design (action space, observation shape): `agent-harness-construction` reference under `ai-agent-engineering`.
- Building/maintaining the MCP server itself: `mcp-builder` / `mcp-servers` references under `ai-agent-engineering`.
- Official docs: Agent SDK "tool search" and API "tool search tool" (`tool_search_tool_regex_20251119` / `_bm25_`, `defer_loading`).
