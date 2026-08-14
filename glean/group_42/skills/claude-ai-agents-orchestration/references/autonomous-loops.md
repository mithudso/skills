<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Formerly the standalone `autonomous-loops` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: autonomous-loops
version: 1.1.0
updated: 2026-05-29
description: >
  Patterns and architecture for autonomous Claude Code loops — from simple
  sequential pipelines to RFC-driven multi-agent DAG systems. TRIGGER: user
  wants to run Claude Code autonomously without human intervention, set up a
  CI/CD-style development pipeline, run parallel agents with merge coordination,
  persist context across loop iterations, or add quality gates and cleanup steps
  to an autonomous workflow. SKIP: interactive single-turn Claude Code usage;
  agent framework selection (use agent-ecosystem); writing an agent plan (use
  agent-plan-writing).
origin: ECC
tags: [autonomous-loops, claude-code, pipeline, dag, parallel-agents, ci-cd, worktrees]
related_skills: [agent-ecosystem, agent-plan-writing, agent-council, git-workflows]
---

# Autonomous Loops

Patterns, architecture, and reference implementations for running Claude Code
autonomously in a loop. Covers the full spectrum from simple `claude -p`
pipelines to RFC-driven multi-agent DAG orchestration.

## When to use this skill

- Building autonomous development workflows that run without human intervention
- Choosing the right loop architecture for your problem (simple vs complex)
- Building CI/CD-style continuous development pipelines
- Running parallel agents with merge coordination
- Implementing context persistence across loop iterations
- Adding quality gates and cleanup steps to autonomous workflows

## Compatibility note

`autonomous-loops` is retained for one release cycle. The canonical skill name
going forward is `continuous-agent-loop`. New loop guidance should be written
here; this skill remains available to avoid breaking existing workflows.

---

## Loop pattern spectrum

From simplest to most complex:

| Pattern | Complexity | Best for |
| --- | --- | --- |
| [Sequential pipeline](#1-sequential-pipeline-claude--p) | Low | Daily dev steps, scripted workflows |
| [NanoClaw REPL](#2-nanoclaw-repl) | Low | Interactive persistent sessions |
| [Infinite agent loop](#3-infinite-agent-loop) | Medium | Parallel content generation, spec-driven work |
| [Continuous Claude PR loop](#4-continuous-claude-pr-loop) | Medium | Multi-day iterative projects with CI gates |
| [De-sloppify pass](#5-de-sloppify-pass) | Add-on | Quality cleanup after any implementation step |
| [Ralphinho / RFC-driven DAG](#6-ralphinho--rfc-driven-dag) | High | Large features with parallel work units and merge queues |

---

## 1. Sequential pipeline (`claude -p`)

The simplest loop. Break daily development into a series of non-interactive
`claude -p` calls. Each call is a focused step with a clear prompt.

```bash
#!/bin/bash
# daily-dev.sh — Sequential pipeline for a feature branch
set -e

# Step 1: Implement the feature
claude -p "Read the spec in docs/auth-spec.md. Implement OAuth2 login in src/auth/. Write tests first (TDD). Do NOT create any new documentation files."

# Step 2: De-sloppify (cleanup pass)
claude -p "Review all files changed by the previous commit. Remove any unnecessary type tests, overly defensive checks, or testing of language features. Keep real business logic tests. Run the test suite after cleanup."

# Step 3: Verify
claude -p "Run the full build, lint, type check, and test suite. Fix any failures. Do not add new features."

# Step 4: Commit
claude -p "Create a conventional commit for all staged changes. Use 'feat: add OAuth2 login flow' as the message."
```

### Key design principles

1. **Each step is isolated** — every `claude -p` call is a new context window; no context leaks between steps.
2. **Order matters** — steps execute sequentially; each builds on the filesystem state left by the prior step.
3. **Negative instructions are dangerous** — don't say "don't test the type system"; add a separate cleanup step instead (see [de-sloppify](#5-de-sloppify-pass)).
4. **Exit codes propagate** — `set -e` stops the pipeline on failure.

### Variants

**Model routing:**
```bash
# Research with Opus (deep reasoning)
claude -p --model opus "Analyze the codebase architecture and write a plan for adding caching..."

# Implement with Sonnet (fast, capable)
claude -p "Implement the caching layer according to the plan in docs/caching-plan.md..."

# Review with Opus (thorough)
claude -p --model opus "Review all changes for security issues, race conditions, and edge cases..."
```

**Context via files (not prompt length):**
```bash
echo "Focus areas: auth module, API rate limiting" > .claude-context.md
claude -p "Read .claude-context.md for priorities. Work through them in order."
rm .claude-context.md
```

**Tool restriction:**
```bash
# Read-only analysis pass
claude -p --allowedTools "Read,Grep,Glob" "Audit this codebase for security vulnerabilities..."

# Write-only implementation pass
claude -p --allowedTools "Read,Write,Edit,Bash" "Implement the fixes from security-audit.md..."
```

---

## 2. NanoClaw REPL

A persistent loop built into ECC. Uses a REPL with session-aware full
conversation history synchronized with `claude -p`.

```bash
# Start the default session
node scripts/claw.js

# Named session with skill context
CLAW_SESSION=my-project CLAW_SKILLS=tdd-workflow,security-review node scripts/claw.js
```

**How it works:**
1. Loads conversation history from `~/.claude/claw/{session}.md`
2. Each user message is sent to `claude -p` with the full history as context
3. Response is appended to the session file (Markdown as database)
4. Session persists across restarts

| Use case | NanoClaw | Sequential pipeline |
| --- | --- | --- |
| Interactive exploration | Yes | No |
| Scripted automation | No | Yes |
| Session persistence | Built-in | Manual |
| Context accumulation | Grows per turn | New context each step |
| CI/CD integration | Poor | Excellent |

See the `/claw` command documentation for full details.

---

## 3. Infinite agent loop

A dual-prompt system for orchestrating parallel subagents for spec-driven
generation. Credit: @disler.

### Architecture

```
PROMPT 1 (Coordinator)               PROMPT 2 (Subagent)
┌─────────────────────┐              ┌──────────────────────┐
│ Parse spec file      │              │ Receive full context  │
│ Scan output dir      │  Deploy      │ Read assigned number  │
│ Plan iterations      │─────────────>│ Follow spec strictly  │
│ Assign directions    │  N agents    │ Generate unique output│
│ Manage batches       │              │ Save to output dir    │
└─────────────────────┘              └──────────────────────┘
```

### Pattern

1. **Spec analysis** — coordinator reads a spec file (Markdown) defining what to generate
2. **Directory scouting** — scans existing outputs to find the highest iteration number
3. **Parallel deployment** — launches N subagents, each with:
   - Full spec
   - Unique creative direction
   - Specific iteration number (no conflicts)
   - Snapshot of existing iterations (for uniqueness)
4. **Batch management** — for infinite mode, deploys waves of 3–5 agents until context is exhausted

### Via Claude Code command

Create `.claude/commands/infinite.md`:

```markdown
Parse from $ARGUMENTS:
1. spec_file — path to the spec Markdown file
2. output_dir — directory to save iterations
3. count — integer 1-N or "infinite"

Phase 1: Read and deeply understand the spec.
Phase 2: List output_dir, find the highest iteration number. Start from N+1.
Phase 3: Plan creative directions — each agent gets a DISTINCT theme/approach.
Phase 4: Deploy subagents in parallel (use Task tool). Each receives:
  - Full spec text
  - Current directory snapshot
  - Their assigned iteration number
  - Their unique creative direction
Phase 5 (infinite mode): Loop in waves of 3–5 until context is insufficient.
```

**Invocation:**
```bash
/project:infinite specs/component-spec.md src/ 5
/project:infinite specs/component-spec.md src/ infinite
```

### Batch strategy

| Count | Strategy |
| --- | --- |
| 1–5 | All agents simultaneously |
| 6–20 | Batches of 5 |
| Infinite | Waves of 3–5, progressively increasing complexity |

**Key insight:** The coordinator **assigns** each agent a specific creative direction and iteration number. Don't rely on agents self-differentiating — assignment prevents conceptual duplication across parallel agents.

---

## 4. Continuous Claude PR loop

A production-grade shell script that runs Claude Code in a continuous loop,
creates PRs, waits for CI, and auto-merges. Credit: @AnandChowdhary.

```
┌─────────────────────────────────────────────────────┐
│  CONTINUOUS CLAUDE ITERATION                        │
│                                                     │
│  1. Create branch (continuous-claude/iteration-N)  │
│  2. Run claude -p with enhanced prompt              │
│  3. (Optional) Reviewer pass — separate claude -p  │
│  4. Commit (claude generates commit message)        │
│  5. Push + create PR (gh pr create)                 │
│  6. Wait for CI checks (poll gh pr checks)          │
│  7. CI failed? → Auto-fix pass (claude -p)         │
│  8. Merge PR (squash/merge/rebase)                  │
│  9. Return to main → repeat                         │
│                                                     │
│  Limits: --max-runs N | --max-cost $X               │
│          --max-duration 2h | completion signal      │
└─────────────────────────────────────────────────────┘
```

**Installation:**
```bash
curl -fsSL https://raw.githubusercontent.com/AnandChowdhary/continuous-claude/HEAD/install.sh | bash
```

**Usage:**
```bash
# Basic: 10 iterations
continuous-claude --prompt "Add unit tests for all untested functions" --max-runs 10

# Cost-limited
continuous-claude --prompt "Fix all linter errors" --max-cost 5.00

# Time-boxed
continuous-claude --prompt "Improve test coverage" --max-duration 8h

# With review pass
continuous-claude \
  --prompt "Add authentication feature" \
  --max-runs 10 \
  --review-prompt "Run npm test && npm run lint, fix any failures"

# Parallel via worktrees
continuous-claude --prompt "Add tests" --max-runs 5 --worktree tests-worker &
continuous-claude --prompt "Refactor code" --max-runs 5 --worktree refactor-worker &
wait
```

### Cross-iteration context: SHARED_TASK_NOTES.md

A `SHARED_TASK_NOTES.md` file persists across iterations. Claude reads it at
iteration start and updates it at iteration end — bridging the context gap
between independent `claude -p` calls.

```markdown
## Progress
- [x] Added authentication module tests (iteration 1)
- [x] Fixed edge case in token refresh (iteration 2)
- [ ] Still needed: rate limiting tests, error boundary tests

## Next steps
- Focus next on the rate limiting module
- Mock setup in tests/helpers.ts can be reused
```

### CI failure recovery

When PR checks fail, continuous-claude automatically:
1. Fetches the failing run ID via `gh run list`
2. Spawns a new `claude -p` with CI fix context
3. Claude checks logs via `gh run view`, fixes code, commits, pushes
4. Re-waits for checks (up to `--ci-retry-max` attempts)

### Completion signal

```bash
continuous-claude \
  --prompt "Fix all bugs in the issue tracker" \
  --completion-signal "CONTINUOUS_CLAUDE_PROJECT_COMPLETE" \
  --completion-threshold 3  # Stop after 3 consecutive signals
```

### Configuration flags

| Flag | Purpose |
| --- | --- |
| `--max-runs N` | Stop after N successful iterations |
| `--max-cost $X` | Stop after spending $X |
| `--max-duration 2h` | Stop after elapsed time |
| `--merge-strategy squash` | squash, merge, or rebase |
| `--worktree <name>` | Parallel execution via git worktrees |
| `--disable-commits` | Dry-run mode (no git operations) |
| `--review-prompt "..."` | Add reviewer pass each iteration |
| `--ci-retry-max N` | Auto-fix CI failures (default: 1) |

---

## 5. De-sloppify pass

An add-on to any loop. After each implementation step, add a dedicated
cleanup/refactor step.

### The problem

When you ask an LLM to implement with TDD, its interpretation of "write tests"
is often overly literal:
- Tests that verify TypeScript's type system works
- Overly defensive runtime checks for things the type system already guarantees
- Tests of framework behavior rather than business logic
- Excessive error handling obscuring actual code

### Why not negative instructions?

Adding "don't test the type system" to the implementation prompt has downstream
effects: the model becomes hesitant about all tests, skips legitimate edge case
tests, and quality degrades unpredictably.

### Solution: separate step

Let the implementation step be thorough, then add a focused cleanup agent:

```bash
# Step 1: Implement (let it be thorough)
claude -p "Implement the feature with full TDD. Be thorough with tests."

# Step 2: De-sloppify (separate context, focused cleanup)
claude -p "Review all changes in the working tree. Remove:
- Tests that verify language/framework behavior rather than business logic
- Redundant type checks the type system already enforces
- Overly defensive error handling for impossible states
- console.log statements
- Commented-out code

Keep all business logic tests. Run the test suite after cleanup to ensure nothing breaks."
```

### In a loop

```bash
for feature in "${features[@]}"; do
  claude -p "Implement $feature with TDD."
  claude -p "Cleanup pass: review changes, remove test/code slop, run tests."
  claude -p "Run build + lint + tests. Fix any failures."
  claude -p "Commit with message: feat: add $feature"
done
```

**Key insight:** Two focused agents beat one constrained agent. Rather than adding negative instructions with downstream quality effects, add a separate de-sloppify step.

---

## 6. Ralphinho / RFC-driven DAG

The most complex pattern. An RFC-driven multi-agent pipeline that decomposes a
spec into a dependency DAG, runs each unit through a layered quality pipeline,
and lands changes through an agent-driven merge queue. Credit: @enitrat.

### Architecture overview

```
RFC/PRD Document
       │
       ▼
  Decompose (AI)
  Split RFC into work units with dependency DAG
       │
       ▼
┌──────────────────────────────────────────────────────┐
│  RALPH LOOP (up to 3 rounds)                         │
│                                                      │
│  For each DAG layer (in dependency order):           │
│                                                      │
│  ┌── Quality pipeline (each unit in parallel) ───┐  │
│  │  Each unit in its own worktree:               │  │
│  │  research → plan → implement → test → review  │  │
│  │  (depth varies by complexity tier)            │  │
│  └────────────────────────────────────────────┘  │
│                                                      │
│  ┌── Merge queue ─────────────────────────────┐     │
│  │  Rebase to main → run tests → merge/evict  │     │
│  │  Evicted units re-enter with conflict ctx  │     │
│  └────────────────────────────────────────────┘     │
│                                                      │
└──────────────────────────────────────────────────────┘
```

### RFC decomposition

AI reads the RFC and generates work units:

```typescript
interface WorkUnit {
  id: string;              // kebab-case identifier
  name: string;            // Human-readable name
  rfcSections: string[];   // Which RFC sections this addresses
  description: string;     // Detailed description
  deps: string[];          // Dependencies (other unit IDs)
  acceptance: string[];    // Concrete acceptance criteria
  tier: "trivial" | "small" | "medium" | "large";
}
```

**Decomposition rules:**
- Prefer fewer, cohesive units (minimize merge risk)
- Minimize cross-unit file overlap (avoid conflicts)
- Keep tests with implementation (never split "implement X" + "test X")
- Only set dependencies where real code dependencies exist

**Dependency DAG determines execution order:**
```
Layer 0: [unit-a, unit-b]     ← no dependencies, run in parallel
Layer 1: [unit-c]             ← depends on unit-a
Layer 2: [unit-d, unit-e]     ← depends on unit-c
```

### Complexity tiers

| Tier | Pipeline stages |
| --- | --- |
| **trivial** | implement → test |
| **small** | implement → test → code-review |
| **medium** | research → plan → implement → test → PRD-review + code-review → review-fix |
| **large** | research → plan → implement → test → PRD-review + code-review → review-fix → final-review |

Trivial changes skip expensive operations; large changes get maximum review depth.

### Independent context windows (eliminating author bias)

Each stage runs in its own agent process with its own context window:

| Stage | Model | Purpose |
| --- | --- | --- |
| Research | Sonnet | Read codebase + RFC, generate context doc |
| Plan | Opus | Design implementation steps |
| Implement | Codex | Write code per the plan |
| Test | Sonnet | Run build + test suite |
| PRD Review | Sonnet | Spec compliance check |
| Code Review | Opus | Quality + security check |
| Review Fix | Codex | Address review issues |
| Final Review | Opus | Quality gate (large tier only) |

**Key design:** The reviewer never wrote the code it reviews — this eliminates author bias, the most common cause of missed issues in self-review.

### Merge queue with eviction

```
Unit branch
    │
    ├─ Rebase to main
    │   └─ Conflict? → Evict (capture conflict context)
    │
    ├─ Run build + tests
    │   └─ Failed? → Evict (capture test output)
    │
    └─ Pass → fast-forward merge to main, push, delete branch
```

**File overlap intelligence:**
- Non-overlapping units: speculative parallel landing
- Overlapping units: sequential landing with rebase between each

**Eviction recovery:**
Evicted units receive full context (conflict files, diff, test output) for the next Ralph round:

```markdown
## Merge conflict — resolve before next push

Your previous implementation conflicted with another unit that landed first.
Refactor your changes to avoid the conflicting files/lines below.

{full eviction context with diff}
```

### Data flow between stages

```
research.contextFilePath ──────────────────> plan
plan.implementationSteps ──────────────────> implement
implement.{filesCreated, whatWasDone} ─────> test, review
test.failingSummary ───────────────────────> review, implement (next round)
reviews.{feedback, issues} ────────────────> review-fix → implement (next round)
evictionContext ───────────────────────────> implement (after merge conflict)
```

### Worktree isolation

Each unit runs in an isolated worktree:
```
/tmp/workflow-wt-{unit-id}/
```

Pipeline stages for the same unit **share** one worktree, preserving state (context files, plan files, code changes) across research → plan → implement → test → review.

### Key design principles

1. **Deterministic execution** — pre-decomposition locks parallelism and ordering
2. **Human review at the highest-value point** — the work plan is the single most impactful intervention point
3. **Separation of concerns** — each stage is an independent agent in an independent context window
4. **Context-aware conflict recovery** — full eviction context enables intelligent retry, not blind retry
5. **Tier-driven depth** — trivial changes skip research/review; large changes get maximum scrutiny
6. **Recoverable workflow** — full state persisted to SQLite; resumable from any point

### When to use Ralphinho vs simpler patterns

| Signal | Use Ralphinho | Use simpler pattern |
| --- | --- | --- |
| Multiple interdependent work units | Yes | No |
| Need parallel implementation | Yes | No |
| Merge conflicts likely | Yes | No (sequential is fine) |
| Single-file change | No | Yes (sequential pipeline) |
| Multi-day project | Yes | Maybe (continuous-claude) |
| Spec/RFC written | Yes | Maybe |
| Quick iteration on one thing | No | Yes (NanoClaw or pipeline) |

---

## Choosing the right pattern

```
Is the task a single, focused change?
├─ Yes → Sequential pipeline or NanoClaw
└─ No → Is there a written spec/RFC?
         ├─ Yes → Does it need parallel implementation?
         │        ├─ Yes → Ralphinho (DAG orchestration)
         │        └─ No → Continuous Claude (iterative PR loop)
         └─ No → Do you need multiple variants of the same thing?
                  ├─ Yes → Infinite agent loop (spec-driven generation)
                  └─ No → Sequential pipeline + de-sloppify
```

### Pattern combinations

1. **Sequential pipeline + de-sloppify** — most common. Cleanup pass after each implementation step.
2. **Continuous Claude + de-sloppify** — add a `--review-prompt` with de-sloppify instructions each iteration.
3. **Any loop + verification** — use the ECC `/verify` command or `verification-loop` skill as a quality gate before commit.
4. **Ralphinho-style model routing in simple loops** — route simple tasks to cheaper models even in sequential pipelines:
   ```bash
   claude -p --model haiku "Fix the import ordering in src/utils.ts"
   claude -p --model opus "Refactor the auth module to use the strategy pattern"
   ```

---

## Anti-patterns

1. **No exit condition** — always set `--max-runs`, `--max-cost`, `--max-duration`, or a completion signal.
2. **No context bridge between iterations** — each `claude -p` starts fresh. Use `SHARED_TASK_NOTES.md` or filesystem state to bridge context.
3. **Retrying the same failure** — if an iteration fails, capture error context and provide it to the next attempt.
4. **Negative instructions instead of a cleanup pass** — don't say "don't do X"; add a separate step to remove X.
5. **All agents in one context window** — for complex workflows, separate concerns into distinct agent processes. The reviewer should never be the author.
6. **Ignoring file overlap in parallel work** — if two parallel agents might edit the same file, plan a merge strategy (sequential landing, rebase, or conflict resolution).

---

## References

| Project | Author | Link |
| --- | --- | --- |
| Ralphinho | enitrat | credit: @enitrat |
| Infinite Agentic Loop | disler | credit: @disler |
| Continuous Claude | AnandChowdhary | credit: @AnandChowdhary |
| NanoClaw | ECC | `/claw` command in this repo |
| Verification Loop | ECC | `skills/verification-loop/` in this repo |
