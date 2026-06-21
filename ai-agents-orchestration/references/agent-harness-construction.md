<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Formerly the standalone `agent-harness-construction` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: agent-harness-construction
description: Design and optimize AI agent action spaces, tool definitions, and observation formats to improve task completion rates. Use when building or debugging an agentic loop, defining tool schemas, improving error-recovery contracts, reducing context overload, or benchmarking agent performance (completion rate, pass@1, retries, cost-per-task). TRIGGER: "agent tool keeps failing", "my agent doesn't recover from errors", "tool schema too broad", "agent context budget blowing up", "design an agent harness", "improve agent completion rate". SKIP: general LLM prompt writing (use prompt-deep-optimizer); high-level agent architecture without tool design (use software-architect or agent-ecosystem).
origin: local
version: 1.1.0
category: developer
updated: 2026-05-29
whenToUse:
  - "my agent tool keeps failing and I don't know why"
  - "design tool schemas for an agentic loop"
  - "agent doesn't recover from errors"
  - "how do I structure tool responses so the agent can self-correct"
  - "agent context window keeps filling up"
  - "improve agent completion rate"
  - "benchmark my agent's pass@1"
  - "agent stuck in a retry loop"
  - "how to write a tool observation format"
related_skills:
  - ai-agent-engineering
  - technical-writing-craft
  - software-engineering-patterns
  - prompt-deep-optimizer
---

# Agent Harness Construction

Use this skill when improving how an agent plans, invokes tools, recovers from errors, and converges to a completed state.

## When NOT to use

- General LLM prompt writing without tool design → use `prompt-deep-optimizer`
- High-level multi-agent architecture without tool-level design → use `software-architect` or `agent-ecosystem`
- Autonomous loop orchestration patterns → use `autonomous-loops`

## Core model

Agent output quality is bounded by four factors — fix in order:

1. **Action space quality** — are tools named, scoped, and typed correctly?
2. **Observation quality** — does each tool response give the agent enough to reason forward?
3. **Recovery quality** — does every error path include a root-cause hint and safe retry?
4. **Context budget quality** — is irrelevant content excluded from the active context?

## Action space design

Principles (apply before choosing granularity):

- Use stable, unambiguous tool names.
- Keep input schemas narrow: required fields first, optional fields last.
- Return deterministic output shapes — same success path every call.
- Avoid omnibus tools unless the operations truly cannot be isolated.

### Granularity rules

| Operation type | Recommended granularity |
|---|---|
| High-risk (deploy, migrate, permission changes) | Micro-tools — one operation per tool |
| Common edit/read/search cycles | Medium tools — group by workflow step |
| High-roundtrip-overhead tasks | Macro tools — only when RTT is the dominant cost |

### Example: well-scoped tool schema

```json
{
  "name": "read_file",
  "description": "Read the contents of a single file. Returns an error if the path does not exist.",
  "parameters": {
    "type": "object",
    "properties": {
      "path": { "type": "string", "description": "Absolute file path" }
    },
    "required": ["path"]
  }
}
```

Compare with an anti-pattern: a `file_operations` tool with a `mode` parameter (`read|write|delete|move`) — the agent must guess the mode name and the schema cannot validate per-mode required fields.

## Observation design

Every tool response should include these fields:

| Field | Value |
|---|---|
| `status` | `success` \| `warning` \| `error` |
| `summary` | One-line result |
| `next_actions` | Actionable follow-up steps (empty array on error is wrong — always provide at least one) |
| `artifacts` | File paths or IDs produced (omit if none) |

## Error-recovery contract

For each error path, provide all three:

- **Root-cause hint** — what specifically failed and why
- **Safe retry instruction** — how to retry without making the situation worse
- **Stop condition** — when the agent must halt rather than retry (missing this causes retry loops that exhaust budgets)

## Context budget management

1. Keep the system prompt minimal and stable across turns.
2. Move large guidance blocks to on-demand skills (loaded only when relevant).
3. Prefer file-path references over inlining long documents.
4. **Summarize at phase boundaries** — defined as points where a distinct sub-goal is completed (e.g., after research, before writing). Do not summarize at arbitrary token counts.

## Architecture patterns

| Pattern | Best fit |
|---|---|
| ReAct | Exploratory tasks with uncertain paths |
| Function-calling | Structured, deterministic flows |
| Hybrid (recommended) | ReAct planning + typed-tool execution |

## Benchmarks to track

- **Completion rate** — % of tasks reaching a valid terminal state
- **Retries per task** — signals recovery quality; >2 avg retries indicates contract gaps
- **pass@1 and pass@3** — first-attempt and three-attempt success rates
- **Cost per successful task** — total tokens / successful completions

## Anti-patterns

- Tools with high semantic overlap — the agent cannot choose reliably between them.
- Opaque tool output with no recovery hint — the agent loops without progress.
- Error responses that contain only a message and no next-action.
- Context overloaded with references irrelevant to the current task phase.
