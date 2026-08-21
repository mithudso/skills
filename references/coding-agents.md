<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Formerly the standalone `coding-agents` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: coding-agents
title: AI Coding-Agent Design
description: >
  The design discipline of code agents — how autonomous/semi-autonomous LLM
  systems index code, manage context, apply edits, run control loops, design
  tools, and get evaluated (2024-2026). Covers codebase indexing & retrieval
  (Aider tree-sitter repo-map + PageRank ranking; Cursor precomputed embeddings +
  Merkle-tree incremental sync + AST-aware chunking/cAST; hybrid BM25+dense) and
  the central agentic-grep-vs-precomputed-RAG debate (Anthropic dropped embeddings
  for Claude Code; Cursor/Continue/Milvus defend them); edit/diff application
  formats & reliability (whole-file vs unified-diff vs search-replace — Aider's
  finding that edit format swings pass rates 30+ points, lint-gated edits);
  the Agent-Computer Interface (ACI — SWE-agent's compact commands, edit guardrails,
  paged file viewer) and CodeAct executable-action space; test-driven control
  loops & self-repair (plan→edit→test→repair, SBFL, dynamic-analysis debuggers,
  near-miss/SEIDR); Agentless structured localize→repair→validate pipelines as the
  challenge to free-form agents; context management for long-horizon code agents
  (repo-map-as-context, JIT file reads, subagent isolation, compaction, context
  rot); the architecture landscape (Aider, SWE-agent, OpenHands, Cursor, Claude
  Code, Devin, Cline, Continue); and benchmarks (SWE-bench Verified/Lite/Multimodal/
  Multilingual/Live, Aider polyglot + edit-format accuracy, Terminal-Bench;
  contamination via SWE-bench+). TRIGGER: designing or improving a coding agent;
  "repo map / codebase indexing", code retrieval / embeddings vs grep, edit/diff
  format reliability, the Agent-Computer Interface / ACI, plan→edit→test→repair /
  self-repair, Agentless, Aider / SWE-agent / OpenHands / Cursor / Claude Code /
  Cline design, SWE-bench / Aider polyglot / Terminal-Bench. SKIP: the general
  agent action-space/tool-def/observation-format theory (use
  agent-harness-construction — coding ACI specializes it); training a code agent
  with RL / SWE-bench as a reward signal (use agentic-rl); generic test-driven
  development with no agent (use software-engineering-patterns); code retrieval
  loop mechanics in the abstract (use iterative-retrieval).
origin: local
category: developer
version: "1.0"
updated: "2026-05-31"
tags:
  - coding-agents
  - swe-agent
  - aider
  - code-retrieval
  - aci
  - swe-bench
  - llm
  - agent
whenToUse:
  - "designing or debugging an AI coding agent"
  - "codebase indexing / code retrieval (repo map, embeddings, AST chunking, grep-vs-RAG)"
  - "choosing an edit/diff application format and improving its reliability"
  - "designing the Agent-Computer Interface (ACI) / tools / action space for code"
  - "building a plan→edit→test→repair loop or self-repair / Agentless pipeline"
  - "context management for long-horizon code agents (compaction, subagent isolation)"
  - "evaluating a coding agent (SWE-bench family, Aider polyglot, Terminal-Bench)"
whenNotToUse:
  - "general agent action-space/tool/observation theory — use agent-harness-construction (coding ACI specializes it)"
  - "RL-training a code agent / SWE-bench as a reward signal — use agentic-rl"
  - "generic TDD with no agent — use software-engineering-patterns"
related_skills:
  - ai-agents-orchestration
  - ai-llm-model-layer
  - ai-rag-retrieval
  - llm-context-engineering
---

# AI Coding-Agent Design

How LLM systems index code, manage context, apply edits, loop, design tools, and
get evaluated. Five findings dominate the 2024-2026 discipline:

1. **Edit-format reliability is a first-order variable** — *how* an agent expresses
   a change (whole-file vs unified-diff vs search-replace) swings pass rates by
   tens of points, independent of the model.
2. **The Agent-Computer Interface (ACI) is the field's most-imitated idea** —
   agents need purpose-built interfaces (compact commands, linted edits, paged
   viewers), not raw human tooling.
3. **Indexing converged on hybrid** structural (tree-sitter/AST + graph ranking) +
   semantic (embeddings) retrieval — but whether to *precompute an index at all*
   is contested.
4. **Test-based verification is the spine** of the control loop (plan→edit→test→
   repair); "Agentless" showed much of this needs no free-form agent.
5. **Benchmarks consolidated** on SWE-bench (Verified/Lite/Multimodal/Multilingual)
   + Aider polyglot + Terminal-Bench, with test-pass as ground truth.

## Codebase indexing & retrieval

- **Structural — repo map (Aider):** parse every file with **tree-sitter**, build
  a graph (files=nodes, references=edges), run **PageRank** to pick the most
  important identifiers within a token budget. Captures *transitive* importance;
  130+ languages; a compact ranked summary of code not in the chat.
- **Semantic — precomputed embeddings (Cursor):** scan → **Merkle tree of file
  hashes** → AST-aware chunk → embed → store vectors + metadata (obfuscated paths)
  in a remote vector DB; **incremental re-embed only changed files** every ~10 min
  via Merkle diff. Query → embed → vector search → return paths+line ranges →
  client reads local source.
- **AST-aware chunking** keeps functions/classes intact (line-based splitters
  degrade retrieval); **cAST** recursively splits large AST nodes and merges
  siblings within a budget. Embedding models are NL-trained, so raw code embeds
  poorly without context — hence **hybrid BM25+dense** and agentic query rewriting.

**The central debate — agentic grep vs precomputed RAG.** *Precomputed embedding
RAG* (Cursor, Continue, Cody, Milvus MCP): fast semantic recall over huge repos.
*Agentic just-in-time search* (Claude Code, SWE-agent, OpenHands): `grep`/`glob`/
read on demand. **Anthropic deliberately dropped embeddings for Claude Code** —
found agentic search beat it on precision (exact vs fuzzy), simplicity (no index
to maintain), freshness (a prebuilt index drifts during editing), and privacy
(nothing leaves the machine). Counter (Milvus): grep "burns too many tokens" on
large repos and misses conceptual matches → hybrid as the middle path.
*Contested / low-confidence:* no neutral large-scale head-to-head; outcome is
workload- and context-window-dependent; vendor benchmarks are self-interested.

## Edit / diff application formats

| Format | Reliability finding |
|---|---|
| Whole-file | Safe to apply, token-expensive, encourages laziness/truncation on big files |
| **Unified diff (udiff)** | Aider: raised GPT-4 Turbo 20%→61%, "3X less lazy" — models "write data for a program" |
| Search/Replace block | Aider default for many models; brittle if the model can't reproduce source exactly |
| diff-fenced | SEARCH/REPLACE with git-merge markers (e.g. Gemini) |

**Edit format is model-dependent** — the same model swings 30+ points by format.
Failure modes: fuzzy-match misses, line-number drift, lazy/truncated edits.
Mitigations: per-model format choice (Aider), **lint-gated edits** that reject
syntactically broken changes before they land (SWE-agent), reflect-and-repair on
failed application. Diff-XYZ benchmarks *diff understanding* as its own skill.

## The Agent-Computer Interface (ACI) & tool design

SWE-agent's thesis (NeurIPS 2024): agents need a **purpose-built interface**, and
ACI design matters as much as prompt engineering. Principles:
- **Compact specialized commands** (a `search`/`find_file` returning terse
  results) instead of noisy full shells.
- **Guardrails on edits** — the edit command runs a **linter and rejects broken
  edits**, so the agent can't accumulate broken state.
- **Paged file viewer** (a window of numbered lines) — bounds context, gives
  stable edit coordinates.
- **Concise, informative observations** — environment responses formatted for the
  model, noise suppressed.

Tool surface: shell, file read/write/edit, search (grep/glob), **LSP**
(go-to-def, references, diagnostics), test runners, linters, increasingly a
browser. **CodeAct (OpenHands)** uses *executable Python/bash as the action*
rather than a fixed JSON tool menu — a more expressive action space.
*(This specializes the general action-space/observation theory in
`agent-harness-construction`.)*

## Control loops & self-repair

- **The canonical loop: plan → edit → test → repair** — generate edit, run
  tests/compiler, feed the **error traceback** back on failure, repeat for a
  budget. **Cline's Plan/Act** productizes the split (architect mode = design,
  no writes; act mode = implement + approval gates).
- **Self-repair research:** iterative self-repair shows **diminishing returns**
  and weak models loop unproductively; **near-miss syndrome** (almost-correct
  code) → SEIDR-style synthesize-execute-instruct-debug-repair; **dynamic-analysis
  repair** couples the LLM to live debuggers (PDB/GDB) instead of pass/fail only.
  Superficial pass/fail causes inefficient cycles; **structured feedback** (failing
  + passing tests, **SBFL** suspicious-location ranking) works better.
- **Agentless** (FSE 2025) removes autonomous decision-making: a fixed two-phase
  **localize → repair** pipeline (hierarchical file→function→edit localization,
  multi-sample patch + filter/rank). Was the best open-source SWE-bench-Lite
  approach at the time (~27.3%, ~$0.34/issue) and argues much agent complexity is
  unnecessary — a real, contested challenge to free-form-agent orthodoxy.

## Context management for code agents

Repo-map-as-context (token-budgeted ranked summary always present); **just-in-time
file reads** (folder/file structure *is* context engineering; slice big files via
grep/tail, not whole-file dumps); **subagent context isolation** (subagents run in
their own windows, return only distilled results); **compaction** (auto-summarize
as the limit nears for long loops); event-stream/event-log context (OpenHands).
**Context rot** is acute for code: even at 128-200K, agents silently drop info,
recommend inconsistent patterns for the same problem, and hallucinate APIs/paths —
no error raised. (See `llm-context-engineering` for the general discipline.)

## Architecture landscape

| Agent | Indexing | Distinctive choice |
|---|---|---|
| **Aider** | tree-sitter repo map + PageRank | edit-format research; repo-map context; human-in-loop |
| **SWE-agent** | agentic search | the ACI concept; linted edits, paged viewer |
| **OpenHands** | agentic + browser | CodeAct executable-action; event-stream; `AgentDelegateAction` |
| **Cursor** | precomputed embeddings + Merkle sync | Turbopuffer index; obfuscated paths |
| **Claude Code** | **no index — agentic grep/glob/read** | RAG dropped; folder-structure-as-context; subagents + compaction |
| **Devin** | proprietary | hosted long-horizon autonomous SWE |
| **Cline** | agentic + browser | explicit Plan/Act modes + approval gates |
| **Continue** | embeddings | pivoted to CI-first PR enforcement ("Continuous AI") |

**Single vs multi/sub-agent:** single-agent dominates the SWE-bench leaderboard
(SWE-agent, Agentless, CodeAct 2.1 is explicitly a strong *single* agent).
Sub-agents are used more for **context isolation** than parallel problem-solving;
evidence that multi-agent *raises capability* on coding is weak (contested).

## Benchmarks & evaluation

- **SWE-bench family:** Verified (500, human-filtered), Lite (300), Multimodal
  (JS+UI), Multi-SWE-bench/Multilingual (Java/TS/Go/Rust/C/C++), SWE-bench-Live /
  SWE-rebench (anti-contamination), SWE-bench Pro (enterprise).
- **Aider polyglot:** 225 hard Exercism tasks across 6 languages, two attempts
  (measures generation *and* edit-from-feedback), reports **edit-format accuracy**
  as a first-class metric. An antidote to SWE-bench's "Python monoculture."
- **Terminal-Bench:** shell-agent tasks in Docker with programmatic success
  functions; frontier ~55-65%, hardest tier ~25-35%.
- **Test-based verification = ground truth** everywhere. **Contamination
  critiques** (SWE-bench+: solution leakage, weak tests) → live/anti-contamination
  benchmarks; read leaderboard numbers with caution. (SWE-bench as an RL *training
  signal* → `ai-llm-model-layer` (references/agentic-rl.md).)

## Anti-patterns

- **Edit-format brittleness** (search/replace fails on inexact source; diff drift)
  → per-model format, lint-gated edits, repair on failed apply.
- **Context overflow / rot** → inconsistent patterns, vague responses, no error;
  scope tightly, ground via repo map/spec.
- **Context-blindness → hallucinated APIs / paths / convention violations**, often
  surfacing late and triggering destabilizing backtracking.
- **No test verification** — declaring success without running tests.
- Mitigations that recur: tightly-scoped prompts ("only touch auth.ts"), pinned
  expected patterns, repo-map/spec grounding, lint-gated edits, mandatory test run
  before "done."

## Cross-references

- **General action-space / tool-def / observation-format theory** →
  `ai-agents-orchestration` (references/agent-harness-construction.md) (coding ACI is its specialization).
- **RL-training a code agent / SWE-bench as a reward signal** → `ai-llm-model-layer` (references/agentic-rl.md).
- **Abstract code-retrieval loop / cold-start context** → `ai-rag-retrieval` (references/iterative-retrieval.md);
  **advanced retrieval** → `ai-rag-retrieval` (references/advanced-rag-patterns.md).
- **Context compaction / rot / the general discipline** → `llm-context-engineering`.
- **Generic test-driven development** → `software-engineering-patterns`.

## References

SWE-agent (arXiv:2405.15793, NeurIPS 2024) + ACI docs; Aider repo-map / unified-
diffs / edit-formats / polyglot writeups; Agentless (2407.01489, FSE 2025);
OpenHands (2407.16741, ICLR 2025) + CodeAct 2.1; Cursor codebase-indexing blog +
TDS writeup; Anthropic "building agents with the Claude Agent SDK" + the
Claude-Code-no-indexing analysis; Milvus grep-vs-RAG counterpoint; Cline Plan/Act;
cAST (2506.15655, EMNLP 2025); Diff-XYZ (2510.12487); SWE-bench Verified +
SWE-bench+ (2410.06992); Terminal-Bench (2601.11868). *(32 sources, 2024–2026;
full URLs in the source research report.)*
