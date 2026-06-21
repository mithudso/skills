---
name: deep-optimizer
description: >-
  Convergence-loop optimizer family ROUTER. Multi-pass audit-and-fix loops for code, prose, prompts, skills, SQL queries, and UI/UX designs — each loops to zero Medium+ findings with build/verify gates. Routes to: code-deep-optimizer (source files & repos, 16-pass audit, CDO alias /cdo); document-critique (prose docs, passes 0–14, DDO alias /ddo); prompt-deep-optimizer (production prompts in code, PDO alias /pdo); skill-optimizer (SKILL.md files, 15 passes, SKO alias /sko); deep-query-optimizer (SQL queries, EXPLAIN-verified, DQO alias /dqo); design-deep-optimizer (graphic/UI/UX screens, 11-pass critique, DDeSO alias /deso). Route to the matching sub-hub. SKIP: MongoDB MQL/aggregation → mongodb-expert (references/deep-mongodb-mql-query-optimizer.md).
origin: local
model: claude-opus-4-8
effort: high
version: "1.0.0"
updated: "2026-06-21"
---

> **Output rules:** Skip preamble and recaps. When delivering fixes, output diffs/edits directly. Required structured outputs (convergence tables, findings tables) are not preamble; keep them.

# deep-optimizer — convergence-loop optimizer family ROUTER

Six siblings. Each runs a domain-specific multi-pass audit, applies every Medium+ fix in place, verifies the result, and loops to convergence.

| Sub-hub | Domain | Alias | Key gates |
|---|---|---|---|
| `code-deep-optimizer` | Source files & repos (any language/framework) | `/cdo` | build/lint/tests verify gate; 16 passes |
| `document-critique` | Prose documents (specs, RFCs, runbooks, KBs) | `/ddo` | passes 0–14; blind re-audit; human-voice pass |
| `prompt-deep-optimizer` | Production prompts shipped in code | `/pdo` | 16 passes; injection guard; algorithm pick |
| `skill-optimizer` | `SKILL.md` files | `/sko` | 15 passes A–O; Pass H trigger eval; hub sync |
| `deep-query-optimizer` | SQL queries (Postgres/MySQL/SQLite/SQL Server) | `/dqo` | EXPLAIN/EXPLAIN ANALYZE verify; index DDL |
| `design-deep-optimizer` | Graphic/brand & UI/UX screens | `/deso` | 11-pass critique; WCAG verify; code-backed fixes |

## Routing

Identify the artifact type and route:

- **Code (`.js`, `.ts`, `.py`, `.go`, `.rs`, `.java`, any source file or repo)** → `code-deep-optimizer`
- **Prose document (spec, RFC, README, runbook, KB article, weekly update)** → `document-critique`
- **Production prompt (system prompt, agent instruction block, tool template in codebase)** → `prompt-deep-optimizer`
- **Skill file (`SKILL.md`, Claude Code skill)** → `skill-optimizer`
- **SQL query** → `deep-query-optimizer`
- **MongoDB MQL / aggregation pipeline** → `mongodb-expert` (references/deep-mongodb-mql-query-optimizer.md) (mongodb family)
- **UI/UX screen, mockup, design asset, or HTML/CSS** → `design-deep-optimizer`

## Cross-hub map

| Hub | Owns |
|---|---|
| `deep-optimizer` | This router — all optimizer siblings |
| `mongodb-expert` | `deep-mongodb-mql-query-optimizer` (MQL/aggregation optimizer) |
| `software-engineering-patterns` | `/code-review` (one-shot diff review, not a convergence loop) |
| `skill-tree-architect` | Whole-tree taxonomy rebalance (not a per-artifact optimizer) |
