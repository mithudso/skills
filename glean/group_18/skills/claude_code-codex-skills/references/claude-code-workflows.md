<!-- hub-reference-banner -->
> **Reference file — part of the `claude-code-skills` hub.** Installed by `/dr` research (2026-06-10).
> Sibling topics (skill anatomy/authoring → the hub's `claude-code-skills-context.md`; plugin packaging/hooks config → `references/claude-code-plugins.md`) are reference files under this hub — load those rather than re-deriving them here.

---

---
name: claude-code-workflows
title: Claude Code Workflows
description: >
  Repeatable practitioner patterns for driving Claude Code through real development work
  (2025–2026): the Explore→Plan→Code→Commit loop and plan mode scoping; verification-first
  workflows and the Stop-hook/goal escalation ladder; enforced TDD loops and test-gaming
  defenses; the dissolved think-keyword pattern and the effort knob; automation via hooks
  ("CLAUDE.md is advisory, hooks are deterministic"), headless claude -p, the Agent SDK and
  GitHub Actions; parallelism via worktrees, subagents, background agents and experimental
  agent teams; context management (CLAUDE.md discipline, /clear//compact hygiene, native
  auto memory, frequent intentional compaction at 40–60% utilization); and the documented
  workflow failure modes with guardrails. Per-claim confidence tags; sources cited.
origin: local
version: "1.0.0"
updated: "2026-06-10"
---

# Claude Code Workflows

Repeatable practitioner patterns for driving Claude Code through real development work (2025–2026). Covers the interactive loop, verification, automation (hooks/headless/SDK/CI), parallelism, context management, and workflow failure modes. Confidence tags: [HIGH] = 3+ independent sources agree; [MEDIUM] = 2 sources or quality-source-with-caveats; [LOW] = single-source/contested. Sources and access dates in References.

## Overview

Claude Code workflows converged on one governing principle, stated verbatim in Anthropic's best-practices doc: **"The context window is the most important resource to manage."** Nearly every documented pattern — plan-first loops, subagent fan-out, frequent compaction, session hygiene — derives from that constraint plus a second principle: **give the agent a way to verify its work** (claimed to improve output quality 2–3×, self-reported). The domain is fast-moving: several 2025 patterns (think-keyword budgets, plan-mode-for-everything) were obsoleted by late-2025/2026 releases; each case is flagged below.

## Core Concepts

### 1. The canonical interactive loop: Explore → Plan → Code → Commit [HIGH]

Enter plan mode and let Claude read files without changing anything; iterate on a written plan (editable via Ctrl+G) until it is right; switch out and implement while verifying against the plan; commit with a descriptive message and open a PR. Rationale: "Letting Claude jump straight to coding can produce code that solves the wrong problem." Scope rule: "If you could describe the diff in one sentence, skip the plan."

- Plan-quality lever: have a second Claude review the plan "as a staff engineer"; when implementation goes sideways, return to plan mode and re-plan rather than push forward.
- [MEDIUM — evolution, not reversal] Plan mode's importance is declining as models improve: its own creator now mostly uses auto mode for scoped work ("plan mode mattered for Opus 4 through ~4.5"). Planning effort should scale with task ambiguity; the threshold for "needs a plan" rose with model generation.
- Spec-first kickoff for larger features [HIGH]: "Interview me in detail using AskUserQuestion… write a complete spec to SPEC.md," then execute the spec in a fresh session.
- Onboarding workflow [HIGH]: for new users and new codebases, start with codebase Q&A ("ask Claude what you'd ask a senior engineer"), not code generation.
- Course correction [HIGH]: Esc to interrupt; checkpoints (`/rewind`, Esc+Esc) restore code/conversation/both; **two-strikes rule** — after two failed corrections on the same issue, `/clear` and rewrite the prompt; a clean session with a better prompt beats a long session with accumulated corrections.

### 2. Verification-first workflows [HIGH]

"Give Claude a way to verify its work" is the single most-cited quality lever — the difference between a session you watch and one you walk away from. Provide a pass/fail check (test suite, build exit code, linter, screenshot/browser test) and the loop closes itself. Real verification means *running the thing*, not just unit tests/lint.

Escalation ladder for how hard verification gates the stop (2026 docs): (1) in-prompt "run the check and iterate" → (2) `/goal` evaluator re-checks after every turn → (3) deterministic **Stop hook** blocks turn-end until the check passes (force-ends after 8 consecutive blocks) → (4) second-opinion verification subagent in fresh context. Anti-overengineering caveat: a reviewer asked to find gaps always finds some — instruct it to report only correctness-relevant gaps.

### 3. TDD with an agent — ordering must be enforced [HIGH]

Claude is implementation-first by default; TDD does not happen from CLAUDE.md prose alone (one experiment: test written *after* code 6/10 times despite explicit instructions — "the hook is the part that does the work"). Working enforcement patterns:
- PreToolUse hook requiring a failing test to exist before source files can be edited.
- Subagent-isolated red/green/refactor (test-writer cannot see implementation plans; implementer sees only the failing test); hook-injected phase gates raised compliance dramatically (self-reported ~20%→~84%).
- Human-written tests as contract: "make the tests pass without modifying the test file."
- **Test-gaming is the known failure mode** [HIGH]: agents adjust assertions to "cheat green." Defenses: never modify a test to make it pass; one test at a time; confirm red before green; Writer/Reviewer split (a fresh context "won't be biased toward code it just wrote").

### 4. Thinking budgets: a dissolved pattern [HIGH — historical]

The 2025 "think < think hard < ultrathink" keyword ladder is **obsolete**: v2.0 (Sept 2025) made thinking a binary Tab toggle; by Jan 2026 ultrathink was deprecated and thinking is on by default; at the API level adaptive thinking + an `effort` parameter replaced manual budgets. Current knobs: `/effort`, `MAX_THINKING_TOKENS`. [MEDIUM] Reasoning effort is an operational quality knob — Anthropic's April 2026 postmortem reverted a default-effort downgrade ("the wrong tradeoff"); practitioners pin effort high for complex work.

### 5. Automation: hooks, headless, Agent SDK, CI [HIGH]

- **Hooks doctrine:** "CLAUDE.md is advisory, hooks are deterministic." Canonical uses: format-after-edit (PostToolUse), test-on-write (TDD loops), ordering/verification gates, destructive-command blocks. Practitioner consensus: behavioral rules decay as context fills — anything that *must* hold moves from prose into hooks ("I stopped writing rules and started writing code").
- **Headless (`claude -p`):** Claude Code as a Unix tool in scripts/pre-commit/CI. Operational discipline from independent guides: hard wall-clock timeouts, `--max-turns` against runaway loops, exit-code mapping ("found issues" ≠ "failed to run"), `--allowedTools` scoping, model-tiering for high-frequency checks, "rigid infrastructure, iterable prompt." [MEDIUM, 2026] `--bare` skips auto-discovery of hooks/skills/MCP/CLAUDE.md for reproducible CI runs.
- **Agent SDK** (renamed from Claude Code SDK, Sept 2025): the same agent loop as a library (TS/Python/CLI). First widespread SDK workflow: *routines* — persistent loops watching tickets/CI that review code, babysit PRs, fix CI, rebase ("the interface moved from source code, to agent, to loop") [MEDIUM].
- **GitHub Actions** (`anthropics/claude-code-action@v1`): interactive mode answers `@claude` mentions; automation mode fires on events with a `prompt` input (e.g., auto-review every PR). Related pattern: agentic CI debugging — push a draft PR, read CI failures via `gh`, iterate against CI rather than local tests.

### 6. Parallelism: worktrees, subagents, background agents, teams [HIGH]

- Running multiple sessions in parallel is the most-cited productivity unlock; **git worktrees** are the standard isolation mechanism (3–5 worktrees, one session each). Equivalent alternatives practitioners actually use: separate full checkouts, fresh /tmp clones — the isolation matters, not the mechanism. [MEDIUM] Native support: `claude --worktree` / `-w` (+ `--tmux` panes); subagents can take throwaway worktrees via `isolation: worktree`.
- **The bottleneck shifts to human review**: practitioners report a hard attention ceiling (~1 significant change reviewable at a time; "wheels fall off past four sessions"); disciplines: rebase-don't-merge between worktrees, ~20-minute review ticks, notifications off ("it's my job to decide when to interrupt the agent"), parallel slots for low-cognitive-overhead side tasks; ~10–20% of remote sessions fail as accepted overhead at 10–15 concurrent.
- **Subagent fan-out** is the official context-preservation tool: subagents run in separate context windows and report summaries — fan out *reads* (research, review). [MEDIUM — contested] Write-heavy tasks parallelize poorly ("tasks mixing reads and writes create chaos"); serialize writes, or use fresh sessions + markdown handoff files. N-way implementation competition (3 agents, 3 worktrees, pick the winner) is a practitioner extension.
- **Background execution** is layered: background bash; background subagents (`run_in_background`); agent view (`claude agents`); cloud sessions; `--teleport` to move sessions between web and terminal.
- **Agent teams** [HIGH, experimental]: a team-lead session coordinates full Claude Code instances with a shared dependency-tracked task list + inter-agent mailbox (`CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS`). Teammates do NOT inherit the lead's history (task files + SendMessage are the only coordination channels) and are not worktree-isolated — partition work by file ownership; pattern: plan first in plan mode, then hand the approved plan to the team.

### 7. Context-management workflows [HIGH]

- **CLAUDE.md:** generate with `/init`, check into git; include only what Claude can't infer (build commands, deviating conventions, gotchas); per line ask "would removing this cause mistakes?" — bloated CLAUDE.md causes instruction-ignoring. Hierarchy: managed → `~/.claude/CLAUDE.md` → project → `CLAUDE.local.md` → parent/child dirs; `@path` imports (≤5 hops). Emphasis markers ("IMPORTANT", "YOU MUST") measurably improve adherence. 2026 division of labor: CLAUDE.md = always-true context; **skills** = on-demand domain knowledge; **hooks** = must-always-happen behavior.
- **Session hygiene:** `/clear` between unrelated tasks; `/compact [focus]` at natural boundaries; `/context` to see where the window goes; partial compaction via checkpoints; compaction-behavior customization in CLAUDE.md ("when compacting, preserve the modified-file list and test commands").
- **Native memory (v2.1.x, early 2026)** [HIGH]: two mechanisms — CLAUDE.md (human-written) + auto memory (Claude-written: `~/.claude/projects/<project>/memory/`, `MEMORY.md` index loaded each session, topic files on demand, `/memory` to manage). Replaced the cottage industry of MCP memory servers.
- **Frequent intentional compaction (HumanLayer methodology)** [HIGH — most influential independent methodology]: keep context utilization at 40–60% by designing the workflow as research → plan → implement, each phase starting near-fresh and compacting into a reviewed artifact (`research.md` ~200 lines, `plan.md` ~200 lines, `progress.md` ~100 lines); subagents do noisy grep/read work in isolated contexts. Converges with Anthropic's spec-session pattern.
- **Context rot is operationally real** [MEDIUM]: behavioral rules obeyed early are abandoned by ~message 30 as task content crowds them out; "temporarily skipped" tests are forgotten after compaction. Guardrails: shorter sessions, deliberate compaction at task boundaries, hooks over prose, artifacts on disk over conversation memory.

## Anti-Patterns

Officially named: **kitchen-sink session** (mixed tasks polluting context); **correcting over and over** (clear after two strikes); **over-specified CLAUDE.md**; **trust-then-verify gap** ("if you can't verify it, don't ship it"); **infinite exploration** (scope it or subagent it).

Practitioner-documented:
- **Verification theater / claimed-done-isn't-done** — the runtime takes the model's text as truth; the agent that writes the code must not be the invocation that certifies it done. Defense: evidence-based completion (hooks re-read edits; cross-check claims against tool logs).
- **Test-gaming** (see §3). **Scope creep/cutting mid-task** — agents cut scope silently as context fills; defense: plan files as contract + adversarial diff-vs-plan review.
- **Subtle correctness bugs that survive review** — the missing-`await` class passes tests, breaks under load; "reviewing is a different mental mode than writing."
- **Runaway sessions and cost blowouts** — recursive self-correction, stalled fleets multiplying API calls; defenses: `--max-turns`, timeouts, allowlists, fleet-retreat-to-single-session.
- **YOLO-mode risk** — `--dangerously-skip-permissions` spreads because approval prompts train click-through; credible answer is sandboxing + auto-mode classifier + allowlists.
- **Boundary/sequencing failures** — failures cluster at system boundaries (migration-vs-deploy ordering), not inside generated code; highest-value intervention is a human asking the boundary question.
- **Deskilling drift** — guardrail: form your own hypothesis before asking the agent (cross-ref: vibe-coding reference in `ai-agents-orchestration`).

## Troubleshooting

| Symptom | Likely cause | Fix |
|---|---|---|
| Quality degrades late in session | Context rot / rule decay | `/compact` at boundary or `/clear` + better prompt; move rules to hooks |
| Agent claims done, isn't | Verification theater | Stop hook / `/goal` check; fresh-context reviewer |
| Tests "pass" suspiciously | Test-gaming | Lock test files; red-before-green; Writer/Reviewer split |
| Repeated failed fixes | Context poisoned | Two-strikes rule: `/clear`, rewrite prompt |
| Parallel sessions colliding | Shared working tree | Worktrees/checkouts; partition by file ownership |
| CI run hangs/burns tokens | Runaway loop | `--max-turns`, wall-clock timeout, `--bare`, exit-code mapping |

## References

Access date 2026-06-10. Independence flags: Cherny material (b,c) collapses to ~2 origins; headless/CI tutorial cluster is doc-derivative (~2 independent points).

1. code.claude.com/docs/en/best-practices — official best practices (2026 revision). [docs — highest authority]
2. anthropic.com/engineering/claude-code-best-practices (~2025-04) — historical anchor. [docs]
3. "How Anthropic teams use Claude Code" (~2025-07). [docs]
4. code.claude.com/docs/en/hooks-guide; /headless; /github-actions; /agents.md; /agent-teams; /memory; /common-workflows. [docs]
5. anthropic.com/engineering/april-23-postmortem (2026-04) — effort default, thinking-clearing bug. [post-mortem]
6. anthropic.com/engineering/building-agents-with-the-claude-agent-sdk (2025-09). [docs]
7. Boris Cherny cluster: InfoQ 2026-01-10; paddo.dev 2026-01-05; team-tips thread 2026-02; The Neuron podcast digest 2026-06-09. [primary interviews/secondary — ~2 independent origins]
8. HumanLayer, "Advanced Context Engineering" (2025-08-29) — FIC methodology. [expert blog, original]
9. TDD experiments: dev.to/kenimo49 2026-05-25; alexop.dev 2025-11-30; thoughtbot 2026-01-12; loreai.dev 2026-03-10. [post-mortems/expert blogs — independent]
10. Parallelism: incident.io 2025-06-27; simonw.substack.com 2025-10-06; mitchellh.com 2026-02-05; codewithseb.com 2026-04-19; skeptrune.com 2025-05-26. [company report + expert blogs — independent]
11. Ronacher: lucumr.pocoo.org 2025-06-12 / 2025-07-30 / 2025-09-29 — subagent limits, 90% loop, CI debugging. [expert post-mortems]
12. Failure modes: dev.to/voxcore84 140-sessions audit 2026-03-12; antjanus.com 2026-02-27; seanfloyd.dev post-mortem 2026; github.com/anthropics/claude-code#42796 fleet telemetry 2026-04-02. [post-mortems]
13. Willison: simonwillison.net 2025-10-22 (sandboxing); substack 2025-10-06 (parallel agents). [expert]
14. Thinking evolution: claude-code issues #9072 (2025-10), #18072 (2026-01, "ultrathink deprecated"); platform.claude.com extended-thinking docs; changelog. [maintainer/primary]
15. Memory: code.claude.com/docs/en/memory; zenn.dev auto-memory guide 2026-02-19. [docs + expert]
