# Parallel work: patch creation, investigations, and fixes

**Reminder:** *every* meaningful task in this skill runs in a subagent — see [SKILL.md "Context isolation"](../SKILL.md). This document covers the three places multiple subagents run **in parallel**:

1. **Patch creation** (Phase 1.5 and Phase 3.7) — the default. One worktree per branch, all create-patch subagents dispatched in a single message so every Evergreen build kicks off at roughly the same instant. See [SKILL.md "Patch creation: parallel vs. sequential"](../SKILL.md) and [`subagent-create-poll.md`](subagent-create-poll.md) template 1a.
2. **Investigations** (Phase 3.2) — read-only and naturally parallel-safe, fanned out across **two** layers: one investigate-coordinator subagent per failed patch (outer), and one investigate-task subagent per failed task within a patch (inner, only for multi-task patches that aren't compile/build failures).
3. **Fixes on independent branches** (Phase 3.5) — only when the failing branches are not ancestors of one another. Linear cascades stay sequential because the fix lands on the earliest red branch and `gt restack` propagates.

For a linear stack with one cascading failure, fixes go one at a time (the fix has to land on the earliest red branch first, then `gt restack` propagates). Patch creation is *still* parallel, even in that case — the goal there is wall-clock concurrency for Evergreen, not avoiding sequencing.

## Patch creation (parallel by default)

See [SKILL.md "Patch creation: parallel vs. sequential"](../SKILL.md) for the coordinator flow. Key invariants:

- The coordinator creates one worktree per branch with `git worktree add ~/worktrees/<repo>-<branch> <branch>` *before* dispatching any subagent. Each subagent gets its worktree path in its prompt.
- All create-patch subagents are dispatched in a *single* message so they run concurrently — never one-by-one in separate messages.
- `stack_state.py add-patch` flock-serializes concurrent writes. The state file is the only shared resource between the parallel subagents; it is race-free by construction.
- After every sibling subagent has returned, the coordinator **leaves the per-branch worktrees in place** — they're session-scoped and reused for Phase 3 re-patching. They're only removed once the cycle reaches `all_clean` / `excluded_only` or the user stops. See [SKILL.md "Worktree lifecycle"](../SKILL.md).
- If one subagent fails, the others still complete. The coordinator surfaces the failure and offers a single-branch retry rather than aborting the whole batch.

## Investigations (read-only, two-level parallel fan-out)

Investigation is parallelized at two levels — outer per-patch and (when warranted) inner per-task. Both layers dispatch in a single message so their subagents run concurrently, and both are read-only with respect to source code.

### Outer: one investigate-coordinator per failed patch

The top coordinator dispatches one investigate-coordinator subagent (template 3a in [`subagent-investigate.md`](subagent-investigate.md), `model: "opus"`) per failed patch, all in a single message. Each 3a coordinator owns one patch end-to-end and decides which strategy to use inside its own scope:

1. **Compile fail-fast short-circuit** — `polling.compile_fail_fast.suspect_patch_id` matches this patch. 3a pulls the named compile task's logs, records one failure, drafts the plan, and persists `set-findings` with `verdict=real-bug`. No fan-out, no master cross-check (compile errors are branch-local).
2. **Build-failure short-circuit** — `patch.build_failure=true` (poll subagent already detected a build failure on a test-runner task: `has_test_results=false` AND `total_test_count=0` on RUN_ALL_UNIT_JAVA_TESTS, INT_JAVA_*, INT_JAVA_THIRDPARTY_*, or any other generated test-runner task — see [`triage.md`](triage.md)). Same shape as compile fail-fast: identify which failing task triggered the flag, pull THAT task's logs with `filter_errors=false`, record once, draft, persist. No fan-out, no master cross-check.
3. **Single-task short-circuit** — `len(failed_tasks) == 1`. 3a investigates inline: pull test results, cross-check master per failing test, record per-test outcomes, draft, persist. No fan-out — there's nothing to parallelize.
4. **Fan-out path** — anything else (≥2 failed tasks, not a compile/build failure). 3a dispatches the inner layer.

### Inner: one investigate-task per failed task (fan-out path only)

In the fan-out path, the 3a coordinator dispatches one investigate-task subagent (template 3b, `model: "opus"` + "think harder" cue) per task in `failed_tasks`, all in a single message. Each 3b worker reads its own task's MCP data, cross-checks master per failing test, calls `record-failure` / `record-master-broken`, and returns a compact structured `===TASK-INVESTIGATION-RESULT===` block (verdict, suspect-branch, cause, failing tests, master status, 1–5 evidence pointers — file/symbol/line/signal).

The 3a coordinator merges the blocks using a priority-ordered aggregation (real-bug > needs-retry > unknown > master-broken > flake), picks the earliest suspect-branch in stack order if 3b workers disagree, and drafts a unified plan from the union of evidence — grouping steps by file, ordering by dependency.

### State file as the only shared resource

Both layers only ever mutate the state file, via `stack_state.py` (which flock-serializes all writes through the `_locked()` context manager). N concurrent 3b workers within one patch may issue concurrent `record-failure` / `record-master-broken` calls — they serialize on the lockfile but never logically collide because each writes a different `(branch, task, test)` key. The 3a coordinator's later `set-findings` only stamps `patch.findings`, which doesn't touch the `failed_tests` / `test_failures` data written by 3b. The top coordinator collects results by reading `stack_state.py summary` once all 3a subagents return — it never sees raw 3b output.

The investigation prompts explicitly tell both 3a and 3b to **not** modify source files. Investigation is read-only by contract.

## Fixes on independent branches (worktrees required)

Only parallelize fixes when the branches don't have an ancestor/descendant relationship in `gt log`. For a linear stack with cascading failures, the fix goes on the earliest failing branch and `gt restack` propagates it — no parallelism gained.

For two siblings `feat_b` and `feat_c` (both children of `feat_a`) with independent failures, you can fix them in parallel:

### Per-subagent protocol (worktree variant)

Subagent prompt template:

```
You are a worker subagent for the evergreen-stack-ci skill.

Skill:  ~/.claude/skills/evergreen-stack-ci/SKILL.md
State:  <state-file-path>
Stack-root: <stack-root>
Branch:     <branch-name>
Patch-id:   <patch-id of the failed run on this branch>
Worktree:   ~/worktrees/<repo>-<branch>

Task — fix the bug for the given branch in an isolated worktree:
1. Read SKILL.md
2. Read the state file to see this patch's `findings` (notes, cause)
3. `git worktree add <worktree-path> <branch>`
4. `cd <worktree-path>` and apply the fix. Stay focused — only the diff for this bug.
5. `git add -A && git commit -m "<jira>: <terse fix message>"`
6. DO NOT remove the worktree — the coordinator owns cleanup

Return ONE LINE: "<branch>: fix committed <sha-short> in <worktree-path>" or
"<branch>: FAILED — <reason>". No diff output, no log.
```

### Coordinator protocol

The main agent (you) is the only thing that runs `git worktree remove`, `gt restack`, or any patch creation:

1. **Wait for all fix subagents** to return successfully. If any fails, abort the parallel branch and fix sequentially instead.
2. **Verify each commit exists** on the named branch (the worktree might have been on a detached HEAD if a subagent did something weird):
   ```bash
   git log <branch> -1 --oneline
   ```
3. **Remove every fix worktree**:
   ```bash
   git worktree list   # list paths
   git worktree remove ~/worktrees/<repo>-feat_b
   git worktree remove ~/worktrees/<repo>-feat_c
   ```
4. **Verify clean state**:
   ```bash
   git worktree list   # only the main checkout should remain (plus any unrelated worktrees the user has)
   ```
5. **Restack from the main checkout**:
   ```bash
   gt restack
   ```
6. **Re-create patches** for fixed branches and their descendants. Reuse `state.scope` from the state file.

## Why these invariants

- **Subagents commit before any worktree is removed**: a `git worktree remove` only removes the worktree directory; the branch in the main repo still has the commit. But if the subagent didn't actually commit (only staged), removing the worktree loses the work. Verify the commit.
- **Branch refs released before `gt restack`**: Graphite's restack walks every branch in the stack, switches to it, and rebases. If branch X is held by another worktree, the rebase fails with a "fatal: 'X' is already checked out at..." error and `gt`'s recovery messaging is not helpful. Two release strategies:
  - **Fix-worktrees** (this file's scratch worktrees): just `git worktree remove` them — they're disposable per-fix scratchpads.
  - **Per-branch CI worktrees** (the long-lived ones from [SKILL.md "Worktree lifecycle"](../SKILL.md)): `git -C <wt> checkout --detach` to release the branch ref *without* destroying the working tree, then `git -C <wt> checkout -f <branch>` after restack to re-attach to the moved tip. See [SKILL.md Phase 3 Step 6](../SKILL.md) for the exact bracket.
- **One coordinator**: if multiple agents try to write the state file or run `gt restack`, you'll get races. The main agent serializes both.

## When NOT to parallelize

- Linear stack, single failure cascading down: just fix the earliest red branch, restack, re-patch. One worktree. No subagents needed.
- Failures spread across multiple branches but all caused by the same root issue (e.g., a bad import that propagated): fix on the earliest, restack, re-patch — let the descendants pick up the fix automatically.
- Working tree is dirty: clean it up first. Subagents that try to create worktrees from a dirty main checkout will fail in confusing ways.
