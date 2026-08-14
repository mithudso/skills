<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `code-plan-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: code-plan-writing
description: >
  Translates specifications, feature requests, and refactors into structured implementation plans for human developers and AI coding agents. Produces Implementation Plans, RFCs, ADRs, and ExecPlans (PLANS.md) with file maps, ordered tasks, checkpoint gates, and context handoff patterns.
  TRIGGER: "write a plan", "implementation plan", "break down this feature", "decompose this work", "create an exec plan", "plan this refactor", "write an RFC", "task breakdown", "how should I structure this work", "plan before coding", "spec to tasks", "write an ADR", "PLANS.md".
  SKIP: single-file fix under 50 lines (just do it); writing the specification itself (use a spec-writing skill); exploratory research with no deliverable; reviewing or critiquing an existing plan (use document-critique); purely operational work with no code changes; multi-agent orchestration or non-code agent task planning (use agent-plan-writing).
version: 1.1.0
updated: "2026-05-29"
category: developer
tags: [planning, implementation-plan, RFC, ADR, ExecPlan, task-decomposition, agentic-coding]
related_skills:
  - agent-plan-writing
  - document-critique
  - coding-patterns
  - coding-standards
  - superpowers:writing-plans
whenToUse:
  - coding task touches 4+ files and needs a structured execution sequence
  - translating a specification or feature request into ordered, implementable tasks
  - planning work for AI agent execution (ExecPlans, PLANS.md)
  - writing an RFC, ADR, or implementation plan document
  - refactor with a defined end state requiring coordinated multi-file changes
  - establishing checkpoints, handoff points, or review gates
  - task estimation or sizing (story points, t-shirt sizing)
  - multiple agents or developers executing different parts of the same plan
whenNotToUse:
  - single-file fix under 50 lines (just do it directly)
  - writing the specification itself (use a spec-writing skill)
  - exploratory research with no defined deliverable
  - reviewing or critiquing an existing plan (use document-critique)
  - purely operational work (deploy, config change) with no code changes
  - multi-agent orchestration or non-code agent tasks (use agent-plan-writing)
---

# Code Plan Writing

Translating a specification, requirement, or feature request into a structured sequence of implementable tasks.

**Break-even point:** any change touching 4+ files, any refactor with a coherent end state, or any task where "what should this do exactly?" is the hard question.

---

## Output Format

Every plan document must contain these sections in order:

1. **Header** — Feature name, one-sentence goal, architecture summary (2–3 sentences), tech stack
2. **File Map** — Table: Action (Create/Modify/Delete) | File path | Responsibility
3. **Tasks** — Numbered blocks, each containing:
   - Files affected (with line ranges for modifications)
   - Checkbox steps (each step = one action + one verification)
   - Exact code, commands, and expected outputs — no placeholders
4. **Validation** — Per-task done criteria + overall acceptance criteria
5. **Not In Scope** — Explicit list of excluded work

Optional (include when relevant): Estimation, Checkpoint gates, Context handoff notes, Recovery/rollback instructions.

**Delivery:** Save to `docs/plans/YYYY-MM-DD-<feature-slug>.md` unless the user specifies otherwise.

---

## Core Concepts

### Spec-to-Plan Translation

A spec answers "what and why." A plan answers "in what steps." Three stages:

1. **Scope confirmation** — Verify spec boundaries before planning. "Add authentication" could mean OAuth, API keys, or sessions.
2. **File mapping** — List every file to create or modify before writing any task. This locks decomposition decisions.
3. **Task sequencing** — Order tasks so each produces a self-contained, testable change. Dependencies flow forward only.

**Example:** Spec: "add user avatar uploads" → File map: `api/upload.ts`, `models/user.ts`, `tests/upload.test.ts`, `components/AvatarPicker.tsx` → Sequence: schema migration → API endpoint + tests → frontend component → integration test.

### Task Decomposition

**Vertical slicing** cuts through all architecture layers to deliver end-to-end functionality. Preferred for feature work.

**Horizontal slicing** groups work by technical layer (all DB changes, then all API changes, then all UI). Use for infrastructure work; risks late integration failures for features.

**Granularity rule:** each step = 2–5 minutes, one action:
- Write the failing test
- Run it — confirm it fails
- Implement the minimal code to pass
- Run tests — confirm they pass
- Commit

### Plan Document Formats

#### Implementation Plan (for execution)

```markdown
# [Feature] Implementation Plan

**Goal:** [One sentence]
**Architecture:** [2–3 sentences]
**Tech Stack:** [Key technologies]

## File Map
| Action | File | Responsibility |
|--------|------|---------------|
| Create | path/to/new.ts | Description |
| Modify | path/to/existing.ts:45-60 | What changes |

## Task 1: [Component Name]
**Files:** Create: `path/file.ts` | Test: `tests/file.test.ts`

- [ ] Write failing test
- [ ] Run test, confirm failure
- [ ] Implement minimal code
- [ ] Run test, confirm pass
- [ ] Commit
```

Every step must contain actual code, actual commands, and expected outputs. "Add appropriate error handling" is a plan failure.

#### RFC (Request for Comments)

For changes needing team input before implementation. Structure: Abstract, Motivation, Proposal, Alternatives Considered, Rollout Plan, Open Questions. Include an Approvers field.

#### ADR (Architecture Decision Record)

1-page max. Records a decision already made. Structure: Title, Status, Context, Decision, Consequences. Store in `docs/decisions/`. An accepted RFC often produces multiple ADRs.

#### ExecPlan (for AI agents)

Self-contained documents for autonomous multi-hour execution. Required sections: Purpose, Progress (timestamped checkboxes), Surprises & Discoveries, Decision Log, Context & Orientation (full repo-relative paths), Concrete Steps (exact commands + expected outputs), Validation & Acceptance, Idempotence & Recovery.

**Key principle:** the plan must enable an engineer unfamiliar with this codebase to implement end-to-end without oral context transfer.

### Estimation

| Technique | When to use | Format |
|-----------|-------------|--------|
| Story Points | Sprint planning | Fibonacci: 1, 2, 3, 5, 8, 13 |
| T-Shirt Sizing | Roadmap / release planning | XS/S/M/L/XL |
| Planning Poker | Prevent anchoring bias | Simultaneous reveal → discuss outliers |
| Time estimates | Hard deadlines | Hours/days |

Don't mix techniques within the same planning context.

### Checkpoint Design

Three-gate model:

| Gate | Review Focus | Catches |
|------|-------------|---------|
| Spec Review | Scope, design assumptions | Wrong problem framing |
| Plan Review | Task ordering, completeness | Half-finished implementations, conflicting patterns |
| Code Review | Test survival, codebase alignment | Plans that fail under execution |

**Within-plan checkpoints:** commit after each task to create a revert point. For AI agents, each commit is a rollback target.

**Progress tracking:**
```
- [x] (2026-05-25 14:00Z) Database schema migration
- [x] (2026-05-25 14:15Z) API endpoint with tests
- [ ] Frontend component (in progress: layout done, validation remaining)
```

### Context Handoff

**Document-based:** plan carries all context — full file paths, term definitions, embedded code. Preferred for AI agents.

**File-based:** shared state file (`PLANS.md`, `todo.md`) tracks progress and discoveries. Each agent reads before starting, writes before stopping.

**Layered:** `CLAUDE.md` (project-wide) → plan file (task-level) → test files (verification). Load on demand; don't cram everything into one file.

### Plan Review Checklist

**Pre-execution:**
- [ ] Every task has a clear "done" criterion
- [ ] File map accounts for all changes in all tasks
- [ ] No task requires undoing a previous task's work
- [ ] Test strategy is explicit (which tests, how to run)
- [ ] No placeholders, TBDs, or "add appropriate X"
- [ ] Types, function names, property names consistent across tasks
- [ ] Each task can be understood without reading other tasks

**Post-execution:**
- [ ] All tests pass
- [ ] No untracked files left behind
- [ ] Implementation matches spec (not more, not less)
- [ ] Commit history is clean and reviewable

### Agentic Plan Execution

**TDD as verification loop:** write tests first, confirm they fail, implement, confirm they pass. Pass/fail is binary — agents need that signal.

**Three execution frameworks:**

| Framework | Use when |
|-----------|---------|
| **Superpowers** | Solo work, single-repo |
| **GitHub Spec Kit** | Specs need review by non-agent stakeholders |
| **BMAD-METHOD** | Multi-role teams with defined agent roles |

**Agent-specific requirements:** each task must specify the target file, change type, sequence position, and verification method. Abstract instructions cause drift.

---

## Methodology

### Writing a Plan from Scratch

1. Read the spec — understand what and why before planning how
2. Confirm scope — resolve "is X in scope?" before writing tasks
3. Map files — list every file to create, modify, or delete
4. Sequence tasks — each task produces testable output; dependencies flow forward
5. Write steps — one action, one verification, one commit point per step; include actual code
6. Add validation — per-task and overall acceptance criteria
7. Self-review — check for placeholders, type inconsistencies, missing files
8. Save — `docs/plans/YYYY-MM-DD-<feature>.md`

### Translating a Spec into Tasks

1. Extract every requirement (explicit and implied)
2. For each requirement, identify which files change
3. Group file changes into tasks (files that change together = one task)
4. Order by dependency (data model → API → UI)
5. For each task: write test step first, then implementation step
6. Add commit step between tasks

---

## Practical Patterns

| Pattern | Rule |
|---------|------|
| **File-First** | Start every plan with a file map table — forces decomposition decisions upfront |
| **Test-First** | Every implementation step is preceded by a test step — defines "done" unambiguously |
| **Checkpoint Commit** | Commit after every task, not every step — each commit is a revert point |
| **Scope Fence** | Add a "Not in scope" section to every plan — prevents scope creep by humans and agents |
| **Living Plan** | Update the plan during execution — mark steps done, log discoveries, amend if reality diverges |

---

## Anti-Patterns

| Anti-Pattern | Problem | Fix |
|---|---|---|
| **Planning Without a Spec** | Plan solves the wrong problem | Agree on what to build before planning how |
| **Placeholder Steps** | Defers decisions to the implementer | Every step must contain actual content |
| **Monolith Tasks** | Tasks touching 10+ files or >30 min are decomposition failures | Split until each task is 2–5 min of focused work |
| **Plan-Then-Forget** | Plan becomes fiction when reality diverges | Treat plans as living documents |
| **Over-Planning** | More time planning than implementing | Plan should take 10–20% of total implementation time |
| **Context Cramming** | Encyclopedic plan instead of focused plan | Reference CLAUDE.md and other docs; don't duplicate |

---

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| Plan is too large (>8 tasks) | Scope is too wide | Break into sub-plans, each producing working software independently |
| Tasks keep expanding during execution | Spec was under-specified | Stop, update spec, re-plan remaining tasks |
| Agent drifts from the plan | Plan lacks specificity | Add exact code blocks, commands, expected outputs |
| Estimation consistently wrong | Wrong technique for context | Try t-shirt sizing first, then decompose large items |
| Plan review takes too long | Review not scaled to risk | Light review for styling/copy; thorough for migrations/auth |
| Multi-agent conflicts | No coordination mechanism | Use file-based handoff with a single state file |

---

## References

1. [Spec-Driven Development with Claude Code](https://www.datacamp.com/tutorial/spec-driven-development-with-claude-code) — Three execution frameworks
2. [Engineering Planning with RFCs, Design Documents and ADRs](https://newsletter.pragmaticengineer.com/p/rfcs-and-design-docs) — Planning document types
3. [Architecture Decision Records](https://github.com/joelparkerhenderson/architecture-decision-record) — ADR examples and templates
4. [Agile Estimation Techniques](https://www.easyagile.com/blog/agile-estimation-techniques) — T-shirt sizing, planning poker, story points
5. [TDD with Claude Code](https://alexop.dev/posts/custom-tdd-workflow-claude-code-vue/) — Red-green-refactor for agentic coding
