# evergreen-stack-ci

**Category:** Cloud, DevOps & Infrastructure
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/tooling/evergreen-stack-ci/skills/evergreen-stack-ci

## Description
Run Evergreen CI across every branch in a Graphite PR stack (or a single branch) by generating local patches per branch, polling for completion, and triaging failures back to the earliest-introducing branch. Evergreen only auto-builds master-based PRs — non-master-based PRs in a stack fail with permission errors, so this skill works around that by creating CLI patches per branch and tracking them in a state file outside the repo. Use whenever the user says "run evergreen on the stack", "patch every branch", "stack CI", "run CI on each branch in the stack", "find which branch broke evergreen", "evergreen the whole stack", "babysit stack patches", "fix stack CI failures", or any time CI has to run across multiple non-master-based branches because Evergreen won't build them automatically. Also coordinates worktree cleanup, gt restack after fixes, and gt submit --no-stack --update-only for branches with existing PRs (NEVER creates new PRs without explicit user request).

---

# Evergreen Stack CI

Run Evergreen CI across every branch in a Graphite PR stack (or a single branch) by creating per-branch CLI patches, persisting them outside the repo, polling for completion, walking failures back to the earliest branch that introduced them, fixing, restacking, and re-patching — until everything is green or every remaining failure is excluded.

**Why this skill exists.** Evergreen's auto-CI only builds branches/PRs based directly off `master`. Non-master-based PRs in a stack get permission errors. The workaround is creating local CLI patches per branch — but each patch contains all changes since master, so a failure on branch N could have been introduced at any of branches 1..N. This skill manages that complexity: tracking patches per branch, polling, attributing failures to their root branch, and coordinating fixes that may need `gt restack` and re-patching downstream branches.

---

## Invocation

```
/evergreen-stack-ci [start|status|fix|update-prs] [--single-branch] [--sequential] [--profile=<preset>] [--alias=<a>] [--variants=<v1,v2>] [--tasks=<t1,t2>] [--exclude=<frontend|e2e|cypress|...>] [--thirdparty-teams=<auto|all|csv>]
```

| Argument | Meaning |
|---|---|
| `start` | Phase 1 — create one patch per branch in scope and persist state |
| `status` | Phase 2 — poll all tracked patches and report completion / failures |
| `fix` | Phase 3 — triage failures to earliest branch, fix, restack, re-patch |
| `update-prs` | Phase 4 — `gt submit --no-stack --update-only` for branches with existing PRs |
| `--single-branch` | Treat current branch as standalone (skip stack detection, skip `gt restack`) |
| `--sequential` | Create patches one branch at a time from the main checkout (opt-out of the default parallel-via-worktrees flow). See "Patch creation: parallel vs. sequential" below |
| `--profile` / `--alias` / `--variants` / `--tasks` / `--exclude` | Test scope. See [`references/test-scope.md`](references/test-scope.md) |

If no argument is given: no state file → `start`; running patches → `status`; failures present → `fix`.

---

## Hard rules

These are non-negotiable. Each one prevents a real failure mode.

1. **Never `gt submit --stack` and never plain `gt submit`.** Both can create new PRs for branches that don't have one yet. Always use `gt submit --no-stack --update-only`. Skip the submit step entirely for branches without a PR. PR creation is routed to the `graphite-pr-stack` skill, not here.
2. **Restack requires zero conflicting worktree checkouts.** Before `gt restack`, no stack branch may be held by a worktree outside the main repo. Detect via `git worktree list --porcelain`. The skill's *own* per-branch CI worktrees DO hold their branches between phases (see "Worktree lifecycle"), so Phase 3 explicitly detaches them with `git -C <wt> checkout --detach` before restack and re-attaches with `git -C <wt> checkout -f <branch>` after — releasing the branch ref without destroying the 77k-file working tree. For *user-owned* worktrees that conflict, the skill MUST NOT detach, remove, or otherwise modify them — only the user knows whether their working tree state is safe to disturb. Instead, **BLOCK: stop the workflow, surface the exact `git -C <path> checkout --detach` commands the user must run, and wait for explicit confirmation** that the user has detached them. Do NOT proceed with `gt restack` while a user-owned worktree holds any branch that needs to be rebased — surfacing the conflict and continuing anyway will cause restack to fail with cryptic errors or silently corrupt stack state. See Phase 3 step 6.2b for the audit procedure (which runs *after* the skill has already detached its own CI worktrees in step 6.2 — i.e. only conflicts the skill cannot resolve on its own ever block the user).
3. **State file lives outside the repo.** All patch ↔ branch mapping persists to `~/.local/state/evergreen-stack-ci/<repo>--<stack-root>.json` (or `$XDG_STATE_HOME/...`). Never write into `.planning/`, the repo, or any tracked path. See [`references/state-schema.md`](references/state-schema.md).
4. **Earliest failing branch is the prime suspect, not the only one.** Don't blindly fix at the earliest red — verify against master, flake patterns, and unrelated diffs. See [`references/triage.md`](references/triage.md).
5. **Parallel patch creation uses one worktree per branch.** Default behavior: dispatch all create-patch subagents in a single message, each working in its own worktree (`git worktree add ~/worktrees/<repo>-<branch> <branch>`). Never run two `evergreen patch` commands from the *same* checkout — each needs a specific branch on HEAD. With `--sequential`, fall back to one branch at a time on the main checkout. Either way, the state file is the only shared resource — `stack_state.py` serializes concurrent `add-patch` calls via flock.
6. **All work happens in subagents; the state file is the only memory.** The coordinator does not hold patch IDs, log content, MCP responses, or failure details. Every meaningful step (patch creation, status poll, investigation, fix) is a subagent that reads/writes via `stack_state.py` and returns one line. The coordinator orchestrates from `stack_state.py summary`. See "Context isolation" below.
7. **Fail-fast on actionable failures, including in-flight compile failures.** As soon as `summary` shows `poll-decision: actionable_failure`, the polling loop ends and Phase 3 begins — even if other patches are still running. Three triggers:
    - Any patch reaches a terminal `failed` status with at least one non-quarantined, non-master-broken failed test (existing behavior), OR
    - Any patch in `started` *or* `failed` status has a `COMPILE_BAZEL` / `COMPILE_CLIENT_BAZEL` task failure or a build failure on any test-runner task (see Hard Rule 13 and [`references/triage.md` "Build failure on a test-runner task"](references/triage.md) — covers `RUN_ALL_UNIT_JAVA_TESTS`, `INT_JAVA_*`, `INT_JAVA_THIRDPARTY_*`, and any other test-runner task with `has_test_results=false`). When this fires, the suspect is the **earliest** branch in stack order with a compile failure; **descendants of the suspect get `evergreen abort`'d** (their builds inherit the same compile error and would be re-patched anyway). Independent upstream patches keep running.
    - Any patch in `started` status has accumulated **≥ 10 failed tasks** (the `failed_tasks` list written by the poll subagent via `update-status`; threshold lives in `stack_state.py` as `EARLY_TRIAGE_TASK_THRESHOLD`). The suspect is the **earliest** branch in stack order whose running patch has crossed the threshold. Descendants are **not** aborted automatically (unlike compile fail-fast) because the remaining tasks may still pass — the user should triage before deciding to abort. Unlike compile fail-fast, this trigger does **not** auto-transition to Phase 3; it surfaces as a normal test-failure `actionable_failure` requiring user confirmation, and it does **not** add a `compile-fail-fast:` line to `summary`.

    **Compile fail-fast auto-transitions to Phase 3 — no user confirmation.** When `summary` includes a `compile-fail-fast:` line, the coordinator MUST immediately enter Phase 3 (fix) instead of asking the user. Compile errors are real bugs by definition (the code doesn't compile against this base) — they are not flakes, not master-broken, and not ambiguous. Waiting for user input wastes 10–30 minutes per cycle while the rest of the patch grinds out doomed test tasks. The coordinator still surfaces a one-line notice ("Compile fail-fast on `<branch>` (kind=`<kind>`) — auto-fixing now.") so the user sees the transition, but does not pause for a yes/no. **Test-failure-driven `actionable_failure` (no `compile-fail-fast:` line) keeps the existing "tell user; offer Phase 3" behavior** — those failures may be flakes, master-broken, or otherwise ambiguous and warrant explicit confirmation.

    **The auto-fix path still obeys Hard Rule 12.** Skipping user confirmation does not skip local verification — template 4 must still run a local canary `bazel test` (or `bazel build` only when the plan declares `Verify (test): N/A`) before committing and re-patching. Compile-fail-fast usually lands on the `N/A` path because the diff doesn't compile and there's no test target to run yet, but the build-only verification is still mandatory.

    (Quarantined / master-broken failures still do NOT trigger fail-fast.)
8. **Polling never silently stops.** If the chosen polling mechanism (`ScheduleWakeup` or `CronCreate`) cannot continue, surface the reason to the user in one line and offer the alternative. Never end a turn while polling is `in_progress` without either (a) a scheduled wakeup confirmed alive, (b) a cron job confirmed alive (verified via `CronList`), or (c) an explicit user instruction to stop.
9. **Three-strikes quarantine + master-broken filtering.** A test that fails 3 consecutive rounds gets quarantined; tests that are also red on master are excluded from the start. The skill will keep cycling on non-excluded failures and only stop when everything is green or only excluded failures remain. See [`references/three-strikes.md`](references/three-strikes.md).
10. **Thirdparty tests are opt-in and never auto-fallback to "all".** Thirdparty integration tests are the most expensive class in the repo (8 team-specific generators, each emitting many INT_JAVA_* tasks). Default behavior: skip thirdparty unless `--thirdparty-teams=auto|all|<csv>` is provided. If `auto` can't map the diff to a team, **skip the variant and surface a non-blocking persistent reminder on every coordinator output** — never run all 8 generators silently. See [`references/test-scope.md`](references/test-scope.md) for the team→path mapping and reminder protocol.
11. **Never run the full Evergreen build for `mms`.** Full builds (`--profile=full`, `-a full`, or any unscoped patch that fans out into the entire project's variant/task graph) are extremely expensive and slow on mms — tens of variants, hundreds of tasks per branch, ~30+ minutes per branch even when nothing fails, and that cost multiplies across the stack and across re-patch iterations. Default to the narrowest scope that covers the diff (`backend` for Java-only stacks, `frontend` for client-only, `compile` or `unit` for smoke checks). If the user explicitly asks for "full", "everything", "the whole build", or `--profile=full`, push back once: name a narrower equivalent (e.g. `backend` + `--thirdparty-teams=auto`) and confirm they really want full before resolving scope. Only proceed with `full` after the user explicitly re-confirms. This rule applies to alias-based runs too — never silently expand `-a full` or `-a all` from auto-detection.
12. **Local bazel test before re-patch in fix loops.** When the fix-and-commit subagent (template 4) produces a fix that the next Phase 3 step 7 re-patch will validate on Evergreen, it MUST run at least one representative `bazel test` on a target that exercises the changed code BEFORE committing — not just `bazel build`. Surface failures as `FAILED — verify-test: <reason>` rather than committing. The mms Evergreen cycle is ~30+ min per patch; a local test pays 1–3 min and prevents whole cycles of wasted CI on trivially-catchable regressions (e.g. a `@BeforeEach` that throws because every test class shares the same broken stub setup). The investigate-coordinator subagent (template 3a) is required to name a concrete canary target in the plan's `Verify (test):` step. The canary `bazel test` invocation MUST pass `--strict_java_deps=error` so the local run matches Evergreen's strict-deps mode — without it, a missing direct dep can pass locally (warn) and still fail Evergreen (error). See [`references/triage.md` "Why local `bazel test` may pass when Evergreen fails"](references/triage.md) for the trap this closes. **Graceful degradation:** if no bazel test target covers the touched files (rare — usually means the diff is build-config only, generated code, or pure docs), template 3a writes `Verify (test): N/A — <one-line reason>` in the plan; template 4 then runs build-only and proceeds. The N/A path is intentional and not a defect — but it must be explicit in the plan, never silently omitted. See [`references/subagent-investigate.md`](references/subagent-investigate.md) template 3a (step 7 plan structure) and [`references/subagent-fix-submit.md`](references/subagent-fix-submit.md) template 4 (steps 5 and 5b) for the exact prompt language.
13. **`has_test_results=false` is always a build failure first.** Any failed task with `has_test_results==false` AND `total_test_count==0` is a bazel BUILD/COMPILE failure until proven otherwise — the task never reached the test phase. This applies to ALL test-runner tasks, not just `RUN_ALL_UNIT_JAVA_TESTS`: `INT_JAVA_*` (billing, payments, paymentview, etc.), `INT_JAVA_THIRDPARTY_*`, any `_BAZEL` task. The investigation subagent must verify via task logs (with `filter_errors=false` — the default `filter_errors=true` drops bazel `[strict]` and `** Please add ...` lines because they don't match the regex it uses for stderr-style errors) before classifying as flake, master-broken, or test failure. Misclassifying a build failure as a test failure burns entire CI cycles on phantom fixes — and local runs may mask the problem because `bazel test` defaults to `--strict_java_deps=warn` (see Hard Rule 12). See [`references/triage.md` "Build failure on a test-runner task"](references/triage.md) for the full procedure.
14. **Run `evergreen patch` from the coordinator's Bash tool, never inside subagents.** In the mms repo, Go-based CLI tools (including the `evergreen` binary) fail with TLS certificate errors (`x509: OSStatus -26276`) when run inside subagent processes, even when `curl` and other tools succeed from the same machine. The coordinator's own Bash tool does not have this problem.

    Applies to Phase 1 (start), Phase 3 step 7 (re-patch), and any single-branch retry. Change the create-patch workflow:
    - The coordinator runs `evergreen patch` directly via its own Bash tool (one call per branch, sequentially if needed, or in parallel by fanning out multiple Bash tool calls in a single message).
    - The coordinator then calls `stack_state.py add-patch` to record each patch ID.
    - Templates 1a and 1b (the create-patch subagent prompts in [`references/subagent-create-poll.md`](references/subagent-create-poll.md)) are retired for this repo. They may still apply to other repos where the subagent TLS issue does not occur.

    Exceptions: if the coordinator's own Bash tool also hits the TLS error, surface it to the user immediately with the manual command to run.

---

## Context isolation

A stack-CI session can run for hours across many polling/fix cycles. The state file is the durable memory; subagents do the work and persist results; the coordinator reads only the terse `stack_state.py summary` rollup. Without this discipline, MCP JSON blobs and CLI output overflow the coordinator's context within a couple of cycles.

| Coordinator does (sequential, irreversible) | Subagents do (forked context, return one line) |
|---|---|
| Talk to the user; pick phase; ask clarifying questions | Create one Evergreen patch per branch |
| Run `stack_state.py init` (records scope) | Poll all patch statuses + write status updates |
| Read `stack_state.py summary` and decide next move | Investigate one failed patch (read MCP logs, write findings) |
| Run `gt restack` (owns the working tree) | Apply one fix on one branch (in-place or via worktree) |
| Run `git worktree remove` for fix-worktrees | Run `gt submit --no-stack --update-only` for one branch |

For the actual subagent prompt templates — each file includes the common contract + model-selection table at the top — see [`references/subagent-create-poll.md`](references/subagent-create-poll.md) (create-patch 1a/1b, poll-status 2), [`references/subagent-investigate.md`](references/subagent-investigate.md) (investigate-coordinator 3a, investigate-task 3b), and [`references/subagent-fix-submit.md`](references/subagent-fix-submit.md) (fix-and-commit 4, submit-pr 5). For dispatching parallel vs. sequential, see [`references/parallel-fixes.md`](references/parallel-fixes.md).

### Model selection per subagent (REQUIRED on every dispatch)

The coordinator MUST set the Agent tool's `model` argument when dispatching subagents. This is not optional — defaults can drift across Claude Code versions, and the cost/quality trade-off only works if the right tier runs the right step.

| Subagent | Model | Why |
|---|---|---|
| create-patch (1a / 1b) | inherit | Mechanical: `evergreen patch` + one `add-patch` write. |
| poll-status (2) | inherit | Mechanical: MCP polls + state writes; decision logic is in `stack_state.py summary`. |
| **investigate-coordinator (3a)** | **`opus`** | One per failed patch (parallel across patches). Aggregates per-task evidence into a patch-level verdict and drafts the unified implementation plan template 4 executes. Also handles three inline short-circuits — compile-fail-fast, build-failure, and single-task patches — without fanning out. The cross-task synthesis is where opus earns its keep; log reading is delegated. |
| **investigate-task (3b)** | **`opus`** + "think harder" cue | One per failed task within a patch (parallel across tasks, within a 3a fan-out). Pulls one task's test results / logs, cross-checks master per failing test, calls `record-failure` / `record-master-broken`, and returns a compact structured evidence block. The per-task flake-vs-real-bug call is the load-bearing classification in the whole pipeline — a wrong verdict here propagates up into 3a's aggregate plan and downstream into template 4's commit, burning whole CI cycles. Parallelization across tasks (one 3b per failed task, dispatched in a single message) is what keeps wall-clock manageable, not a cheaper per-task tier. |
| **fix-and-commit (4)** | **`sonnet`** + "think hard" cue | Execute the plan template 3a produced. Sonnet is the right tier for code editing; the "think hard" trigger word inside the prompt body elevates reasoning within sonnet. |
| submit-pr (5) | inherit | Mechanical: `gt submit --no-stack --update-only`. |

If a fix subagent returns `FAILED — plan-mismatch:` it means template 3a's plan didn't survive contact with the code. Re-dispatch investigation on that branch (the 3a coordinator on opus, which may fan out to 3b workers again) — do NOT upgrade template 4 to opus to "muscle through" the bad plan. The investigation/implementation split is load-bearing for context isolation; collapsing it lets a single bad classification corrupt downstream work.

---

## Patch creation: parallel vs. sequential

**Default is parallel** so all Evergreen builds start at roughly the same instant. Each subagent runs `evergreen patch` from its own worktree (`git worktree add ~/worktrees/<repo>-<branch> <branch>`), and `stack_state.py` flock-serializes the resulting concurrent `add-patch` writes. Even though Evergreen will queue/process the resulting patches however it does, kicking them off together is what cuts the user's wall-clock wait — they don't pay for build N to finish before build N+1 even enters the queue, and once everything has run, an entire batch of failures can be triaged and fixed in one Phase-3 pass.

`--sequential` is the opt-out. Use it when the user explicitly asks for a serial flow, or when the environment can't afford N worktrees (low disk, weird filesystem, prior worktree conflicts the user doesn't want auto-resolved). The sequential flow uses the main checkout, `gt checkout <branch>` per branch, dispatches one create-patch subagent at a time, and restores the original branch when done.

The coordinator owns worktree lifecycle in both modes:

| Step | Parallel (default) | Sequential (`--sequential`) |
|---|---|---|
| Pre-flight | Confirm `git worktree list` shows no pre-existing checkouts of any stack branch outside the main repo. Conflicts → surface to user. | Same. |
| Setup | For each branch, `git worktree add ~/worktrees/<repo>-<branch> <branch>`. The branch currently checked out in the main repo can stay there — its subagent can use the main checkout as its "worktree". | Note the user's original branch so it can be restored at the end. |
| Dispatch | All create-patch subagents in a single message (template 1a). Each gets its assigned worktree path. | One subagent at a time (template 1b). |
| Cleanup | Worktrees stay alive across phases — they're reused for Phase 3 re-patching. The coordinator only removes them once the stack is terminal (`all_clean` / `excluded_only`) or the user explicitly stops. See "Worktree lifecycle" below. | `gt checkout <original-branch>`. |
| State writes | All concurrent `add-patch` calls are flock-serialized inside `stack_state.py`. | Naturally serial — no contention. |

If a parallel subagent fails (e.g. its `evergreen patch` errors out), the others still complete and are tracked. The coordinator surfaces the failure and offers a single-branch retry rather than aborting the whole batch.

---

## Worktree lifecycle

Per-branch worktrees at `~/worktrees/<repo>-<branch>` are **session-scoped**, not phase-scoped. Each one is a full checkout (77k+ files for `mms`); creating and destroying them per cycle wastes minutes per branch every time the stack re-patches. Keep them alive from Phase 1 through every fix iteration, and only remove them once the cycle is terminal.

| Trigger | Worktree action |
|---|---|
| Phase 1 (`start`) | `git worktree add ~/worktrees/<repo>-<branch> <branch>` for every stack branch (skip the branch already on the main checkout — its subagent uses the main checkout as its "worktree"). Surface pre-existing conflicting checkouts to the user; never auto-remove them. |
| Phase 3 (around `gt restack`) | **Detach before restack, re-attach after.** The per-branch worktrees hold branch refs that `gt restack` needs to rebase, so the coordinator brackets restack with: (a) `git -C ~/worktrees/<repo>-<branch> checkout --detach` for each worktree (releases the branch ref, preserves the 77k-file checkout), (b) `gt restack` from the main repo (now no descendant branch is held elsewhere), (c) `git -C ~/worktrees/<repo>-<branch> checkout -f <branch>` for each worktree (re-attaches to the moved branch tip and force-discards any junk left by the previous `evergreen patch`). Never `git worktree remove` followed by `git worktree add` — that's the wasteful 77k-file recreate. If a worktree directory is missing (user removed it, disk cleanup), recreate just that one. |
| Phase 2 reaches **`all_clean`** | **Coordinator removes every per-branch worktree it created.** Iterate the state file's branches and run `git worktree remove ~/worktrees/<repo>-<branch>` for each (idempotent — `--force` if removal complains about untracked files left from `evergreen patch`; skip silently if the path is already gone). Do this BEFORE offering Phase 4. Surface a one-line confirmation: `removed N stack-CI worktrees`. |
| Phase 2 reaches **`excluded_only`** | Same cleanup as `all_clean` — the cycle is terminal and no further re-patch is coming. |
| User explicitly stops (cron-stop, "stop polling", session end) | Same cleanup. Run it before ending the turn so the user doesn't have stale worktrees lingering after a session they consider done. |
| Phase 2 reaches `actionable_failure` / `needs_attention` | **Worktrees stay alive** — Phase 3 will reuse them. Do NOT clean up here. |
| `--single-branch` runs | The single per-branch worktree (if any was created) follows the same rules — kept across phases, removed only on terminal state or stop. |

The fix-worktrees described in [`references/parallel-fixes.md`](references/parallel-fixes.md) are a *separate* set, used only inside one Phase 3 step for parallel fix application. They are removed at the end of every fix step before `gt restack` (per Hard Rule 2 + invariant #2 in that reference). Do not conflate the two — per-branch CI worktrees persist; per-fix worktrees do not.

---

## Phase 1 — `start`: create patches across the stack

1. **Detect scope** — `gt log short` (>1 branch → stack mode); `git branch --show-current`; `git worktree list`. Single-branch if `--single-branch` or only one branch is found.
2. **Pre-flight worktrees** — for stack mode, ensure no other worktree has any stack branch checked out. Conflicts → stop and surface to user.
3. **Resolve test scope** — pick the profile / alias / explicit variants+tasks; apply exclusions. Confirm with the user. Record in `state.scope`. See [`references/test-scope.md`](references/test-scope.md).
4. **Init state file** — `stack_state.py init --stack-root ... --branches ... --project-id mms --repo-root "$PWD" --profile ... --variants ... --tasks ... --exclude ...`.
5. **Create patches in parallel from the coordinator's Bash tool (default)** — for each branch, `git worktree add ~/worktrees/<repo>-<branch> <branch>` so the desired branch is checked out somewhere the coordinator can `cd` into (or pass via `--ref`) without disturbing the main checkout. Per **Hard Rule 14**, the coordinator itself then runs `evergreen patch` directly via its own Bash tool — fan out one Bash tool call per branch in a single coordinator message so the builds kick off concurrently. After each `evergreen patch` returns a patch ID, the coordinator calls `stack_state.py add-patch` to record it (the script flock-serializes concurrent writes). **Do not remove the worktrees after patch creation** — they're session-scoped and reused for Phase 3 re-patching. The coordinator only tears them down on terminal states (`all_clean` / `excluded_only`) or explicit stop; see "Worktree lifecycle" below. With `--sequential`, run the coordinator Bash calls one at a time instead of in parallel. Templates 1a/1b in [`references/subagent-create-poll.md`](references/subagent-create-poll.md) are retired for mms — see Hard Rule 14. See "Patch creation: parallel vs. sequential" below for the rationale on parallel kickoff.

Confirm with `stack_state.py summary`.

---

## Phase 2 — `status`: poll patch completion

Polling is **periodic, not blocking**. The coordinator drives the loop via `ScheduleWakeup` with a thin self-contained prompt (the canonical wakeup prompt below) that **does NOT re-invoke `/evergreen-stack-ci status`** and does NOT reload [SKILL.md](SKILL.md). Each wakeup spawns one poll-status subagent (template 2) that pulls fresh status, bumps the polling counter, reads the poll-decision, and returns a single line carrying both. Never `bash sleep` or block the conversation.

### Phase 2 entry (or re-entry from Phase 3)

0. **Pre-flight `ScheduleWakeup`.** Before the first poll, issue one harmless probe — `ScheduleWakeup` with `delaySeconds=60` and a no-op prompt that just ends the turn. If it returns the "Wakeup not scheduled" / "loop has ended" error, the skill is running outside `/loop` dynamic mode and the canonical wakeup loop will not work. Switch to **Cron Polling Mode** (below) and tell the user **once** before polling:

   > "Background wakeups are unavailable (`/loop` not active). Falling back to scheduled tasks — polling will fire every 5 min on its own. Reply 'stop' to cancel."

   Probe once at Phase 2 entry only — pick the mode and stay in it for the rest of the cycle.
1. Run `python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py reset-poll-cycle --stack-root <r>` to zero the iteration counter so this cycle gets a fresh 12-wakeup budget.
2. Dispatch **one** poll-status subagent inline in the current turn (template 2). This gives the user an immediate first read without waiting 5 minutes.
3. Apply the decision branch table below to the subagent's return line.
4. **(Wakeup mode only.)** If the decision is `in_progress`, arm the wakeup loop:
   - Call `stack_state.py schedule-next-poll --stack-root <r> --in-seconds 300` so the dashboard knows when to expect the next poll.
   - Call `ScheduleWakeup` with a 5-minute delay and the canonical wakeup prompt below — then end the turn.
5. **(Cron mode only.)** If the decision is `in_progress`, arm the cron loop:
   - `CronList` → confirm no existing job whose prompt starts with `[evergreen-stack-ci-cron stack=<r>]`. If one is present (stale from a prior cycle), `CronDelete` it first.
   - `CronCreate` with `cron: "7-59/5 * * * *"`, `recurring: true`, `durable: false`, and `prompt:` the canonical cron prompt below.
   - Tell the user once that cron polling is armed — then end the turn.

### Canonical wakeup prompt

The coordinator passes this verbatim to `ScheduleWakeup`'s `prompt` argument (substituting `<state-file-path>` and `<stack-root>`). It is the **single source of truth** for wakeup behavior — references to wakeup logic elsewhere should point here.

```
You are resuming an evergreen-stack-ci polling cycle.

State:      <state-file-path>
Stack-root: <stack-root>
Skill:      do NOT load it — the poll subagent reads SKILL.md itself
Cap:        12 wakeups per cycle (the subagent's return line carries iteration <i>/12)

Action:
1. Dispatch ONE poll-status subagent (template 2 in
   ~/.claude/skills/evergreen-stack-ci/references/subagent-create-poll.md). It will return:
     "<n>/<m> green, <k> red, <r> running | iteration <i>/12 | decision: <D>"
2. Switch on D:
   - in_progress       -> if i < 12: run `stack_state.py schedule-next-poll --stack-root <stack-root> --in-seconds 300`,
                          then ScheduleWakeup in 5 min with THIS SAME PROMPT; end turn
                          if i >= 12: ask user (extend cycle? stop?)
   - actionable_failure -> Read `stack_state.py summary --stack-root <stack-root>` (this one
                          read is fine — needed to disambiguate compile vs. test failure).
                          If a `compile-fail-fast:` line is present:
                            * Print one-line notice: "Compile fail-fast on <suspect-branch>
                              (kind=<kind>) — auto-fixing now." (Hard rule 7.)
                            * Immediately enter Phase 3 — do NOT ask the user. The compile
                              error is, by definition, a real bug; descendants are already
                              aborted; there is nothing useful to wait for.
                          Otherwise (test-failure-driven actionable_failure):
                            * Tell user; offer Phase 3 (fix). DO NOT auto-fix — failures
                              may be flakes or master-broken and warrant confirmation.
   - all_clean         -> Tear down per-branch worktrees: for each branch in the state
                          file, run `git worktree remove ~/worktrees/<repo>-<branch>`
                          (idempotent — skip silently if the path is gone; pass --force
                          if removal complains about untracked files from `evergreen patch`).
                          Surface "removed N stack-CI worktrees" to user, then offer
                          Phase 4 (update-prs).
   - excluded_only     -> Tear down per-branch worktrees (same as all_clean). Surface
                          quarantined + master-broken lists AND "removed N stack-CI
                          worktrees" to user; stop polling.
   - needs_attention   -> surface edge case to user (do NOT tear down worktrees — the
                          user may want to re-enter Phase 2/3 after handling the edge case)
3. If the subagent returns "FAILED — ...": surface the failure and stop polling.

Fallback: If ScheduleWakeup returns "Wakeup not scheduled" / "loop has ended" mid-cycle
(the runtime gate flipped off after we entered Phase 2), do NOT silently stop. Hand off
to cron polling:

  1. CronCreate with cron="7-59/5 * * * *", recurring=true, durable=false,
     prompt=<canonical cron prompt below, substituting <stack-root> and <state-file-path>>.
  2. Tell the user once:
       "/loop ended; polling continued via cron — fires every 5 min until done.
        Reply 'stop' to cancel."
  3. End the turn.

The state file already has the latest iteration count, so cron-fired turns pick up
exactly where wakeup mode left off.

Do NOT invoke /evergreen-stack-ci status. Do NOT call the Skill tool.
Do NOT read stack_state.py summary yourself — the decision is in the subagent's return.
```

### Decision per iteration

The poll-status subagent returns a line of the form:
`<succ>/<total> green, <failed> red, <running> running | iteration <i>/<max> | decision: <D>`

| `decision:` | Coordinator action |
|---|---|
| `actionable_failure` **with `compile-fail-fast:` in summary** | STOP polling. **Auto-transition to Phase 3 immediately** — no user prompt. Print one-line notice naming the suspect branch + kind, then proceed to Phase 3 step 0. (Hard rule 7.) |
| `actionable_failure` (no `compile-fail-fast:` line) | STOP polling. Tell user; offer Phase 3. Wait for confirmation before fixing — any patch is terminal-failed with actionable tests, OR any running patch has ≥ 10 failed tasks (`EARLY_TRIAGE_TASK_THRESHOLD`); failures may be flakes or master-broken. |
| `in_progress` | **Wakeup mode** (Step 4): if `<i> < <max>`: call `schedule-next-poll --in-seconds 300`, then ScheduleWakeup in 5 min with the canonical prompt; end turn. If `<i> >= <max>`: ask user (extend? stop?). **Cron mode** (Step 5): if `<i> < <max>`: call `schedule-next-poll --in-seconds 300` so the dashboard countdown stays live, then end the turn — cron auto-fires again in ~5 min; do NOT call ScheduleWakeup. If `<i> >= <max>`: ask user (extend? stop?). The mode was committed at Phase 2 entry (Step 0) and does not change mid-cycle (mid-cycle wakeup failures hand off to cron via the canonical wakeup prompt's Fallback clause). |
| `excluded_only` | STOP polling. **Tear down per-branch worktrees** (see "Worktree lifecycle" — iterate state branches, `git worktree remove ~/worktrees/<repo>-<branch>` each, idempotent). Surface quarantined + master-broken lists to user along with the `removed N stack-CI worktrees` confirmation. |
| `all_clean` | STOP polling. **Tear down per-branch worktrees** (see "Worktree lifecycle" — iterate state branches, `git worktree remove ~/worktrees/<repo>-<branch>` each, idempotent). Surface `removed N stack-CI worktrees` to the user, then offer Phase 4. |
| `needs_attention` | Surface edge case to user (aborted, no-patch, etc.) |

### Cron Polling Mode (when `ScheduleWakeup` is unavailable)

When the Phase 2 entry probe fails (or when wakeup mode dies mid-cycle via the canonical wakeup prompt's Fallback clause), drive polling via `CronCreate` instead. Cron jobs fire on the REPL's idle ticks regardless of whether `/loop` is active, so the user gets the same hands-off cadence they'd get under wakeup mode.

**Arming the cron** (once, when entering this mode):

1. Build the cron prompt from the canonical cron prompt template below (substituting `<state-file-path>` and `<stack-root>`). The prompt MUST start with the literal marker `[evergreen-stack-ci-cron stack=<stack-root>]` so `CronList` can find it later for cleanup.
2. `CronCreate` with `cron: "7-59/5 * * * *"` (every 5 minutes, off-peak — minute 7, 12, 17, …, 57), `recurring: true`, `durable: false`, `prompt: <cron prompt>`. `durable: false` is intentional — `durable: true` would let a stale cron resurrect days later in an unrelated session.
3. Tell the user once: "Polling armed via cron — fires every 5 min until done. Reply 'stop' to cancel."

**Cleanup discipline** (every Phase-2-exit branch):

Unlike `ScheduleWakeup` (which auto-stops if you don't re-arm), `CronCreate` keeps firing until `CronDelete`. Before exiting Phase 2 for any terminal decision (`actionable_failure`, `all_clean`, `excluded_only`, `needs_attention`), or on user-instructed stop, or when the 12-iteration cap is reached and the user picks "stop":

- `CronList` → find the job whose prompt starts with `[evergreen-stack-ci-cron stack=<stack-root>]` → `CronDelete` it.
- Phase 3 entry MUST also call this cleanup before fixing/restacking, otherwise the cron will fire mid-fix and re-poll stale patches. See Phase 3 Step 0.

The 12-iteration cycle budget still applies — `stack_state.py bump-poll-iteration` is called by the poll subagent on every dispatch, and the coordinator surfaces an extend?/stop? question when `iteration >= max`. On stop, `CronDelete`; on extend, `reset-poll-cycle` and leave the cron armed.

### Canonical cron prompt

The coordinator (or the wakeup-mode Fallback clause) passes this verbatim to `CronCreate`'s `prompt` argument (substituting `<state-file-path>` and `<stack-root>`). It is the **single source of truth** for cron-fired polling behavior.

```
[evergreen-stack-ci-cron stack=<stack-root>]

You are firing a scheduled poll for an evergreen-stack-ci cycle (cron mode).
This prompt was armed via CronCreate; the marker on line 1 lets CronList find
this job for cleanup.

State:      <state-file-path>
Stack-root: <stack-root>
Skill:      do NOT load it — the poll subagent reads SKILL.md itself
Cap:        12 iterations per cycle (the subagent's return line carries iteration <i>/12)

Action:
1. Dispatch ONE poll-status subagent (template 2 in
   ~/.claude/skills/evergreen-stack-ci/references/subagent-create-poll.md). It will return:
     "<n>/<m> green, <k> red, <r> running | iteration <i>/12 | decision: <D>"
2. Switch on D:
   - in_progress       -> if i < 12:
                            run `stack_state.py schedule-next-poll
                              --stack-root <stack-root>
                              --repo-root <repo-root>
                              --in-seconds 300`
                            then end turn (cron auto-fires again in ~5 min)
                          if i >= 12: ask user (extend cycle? stop?)
                            on stop:   CronList -> CronDelete this job
                            on extend: stack_state.py reset-poll-cycle (leave cron armed)
   - actionable_failure -> CronList -> CronDelete this job FIRST (so cron can't refire
                          mid-fix and re-poll stale patches).
                          Read `stack_state.py summary --stack-root <stack-root>` to
                          disambiguate compile vs. test failure.
                          If a `compile-fail-fast:` line is present:
                            * Print one-line notice: "Compile fail-fast on <suspect-branch>
                              (kind=<kind>) — auto-fixing now." (Hard rule 7.)
                            * Immediately enter Phase 3 — do NOT ask the user. The compile
                              error is, by definition, a real bug; descendants are already
                              aborted; there is nothing useful to wait for.
                          Otherwise (test-failure-driven actionable_failure):
                            * Tell user; offer Phase 3 (fix). DO NOT auto-fix — failures
                              may be flakes or master-broken and warrant confirmation.
   - all_clean         -> CronList -> CronDelete this job. Tear down per-branch worktrees:
                          for each branch in the state file, run
                          `git worktree remove ~/worktrees/<repo>-<branch>` (idempotent —
                          skip silently if the path is gone; pass --force if removal
                          complains about untracked files from `evergreen patch`).
                          Surface "removed N stack-CI worktrees" to user, then offer
                          Phase 4 (update-prs).
   - excluded_only     -> CronList -> CronDelete this job. Tear down per-branch worktrees
                          (same as all_clean). Surface quarantined + master-broken lists
                          AND "removed N stack-CI worktrees" to user.
   - needs_attention   -> surface edge case to user; CronList -> CronDelete this job.
                          Do NOT tear down worktrees — the user may want to re-enter
                          Phase 2/3 after handling the edge case.
3. If the subagent returns "FAILED — ...": surface the failure;
   CronList -> CronDelete this job; stop polling. Do NOT tear down worktrees on
   subagent failure — the user may want to retry.

Do NOT invoke /evergreen-stack-ci status. Do NOT call the Skill tool.
Do NOT read stack_state.py summary yourself — the decision is in the subagent's return.
```

### Polling cycle budget

Each cycle gets 12 iterations (5 min × 12 = 1 hour). The counter lives in the state file's `polling.iteration_count` and is bumped by the subagent on every wakeup, so the budget survives across coordinator restarts. **The 1-hour cap is per cycle, not per session.** The session itself runs indefinitely until `all_clean` or `excluded_only`.

After Phase 3 fixes and re-patches, re-entry to Phase 2 calls `reset-poll-cycle` to start a fresh 12-iteration budget. See [`references/three-strikes.md`](references/three-strikes.md) for the full decision tree.

### Stale patches after a fix

When Phase 3 creates new patches, older patches on the same branches are stale. The state file's "latest patch per branch" semantics handle this automatically — `summary` and the next poll only consider each branch's latest patch entry. Do NOT cancel stale patches in Evergreen unless the user asks.

### Compile-failure fail-fast

The poll path detects in-flight compile failures so the cycle exits early instead of waiting 20–40 minutes for every downstream test task to fail with no useful signal. Triggers:

- **`COMPILE_BAZEL` / `COMPILE_CLIENT_BAZEL` task failures** — `_detect_compile_fail_fast` matches these (case-insensitive substring, so per-team variants like `compile_bazel_payments` count) directly off the `failed_tasks` field that the poll subagent already writes via `update-status`. No extra MCP call.
- **Build failure on any test-runner task** — the broader pattern from [`references/triage.md` "Build failure on a test-runner task"](references/triage.md): ANY failed task with `has_test_results=false` AND `total_test_count=0` is a bazel build/compile error that happened before any test could run. This covers `RUN_ALL_UNIT_JAVA_TESTS` (the original case — "No test results found — there may have been a build failure"), all `INT_JAVA_*` variants (billing, payments, paymentview, …), `INT_JAVA_THIRDPARTY_*`, and any other generated test-runner task. The poll subagent (template 2 step 5c) inspects every failed task with that signature, pulls its logs with `filter_errors=false` (so bazel `[strict]` / `** Please add ...` lines aren't dropped), and calls `mark-build-failure --patch-id <id>` if the pattern matches.

Detection runs every poll. When it fires, `summary` adds two lines:

```
compile-fail-fast: <suspect-branch> kind=<COMPILE_BAZEL|COMPILE_CLIENT_BAZEL|BUILD_FAILURE> patch=<id>
compile-fail-descendants: <branch>=<patch_id>,...
```

…and the poll-decision is forced to `actionable_failure`. The suspect is the earliest branch in stack order with a compile failure; descendants are later branches whose latest patch is still abortable.

**Aborting descendants.** Template 2 step 5d runs `evergreen abort -p <patch_id>` on each descendant once and then calls `mark-fail-fast-aborted` to flip an idempotency flag. Subsequent polls in the same cycle return nothing from `get-fail-fast-aborts` so we never double-abort. Independent upstream patches keep running because they may still surface real failures worth knowing about.

**Auto-transition to Phase 3.** Per Hard Rule 7, when the `compile-fail-fast:` line is present in `summary`, the coordinator skips the "tell user; offer Phase 3" handshake and enters Phase 3 immediately. The rationale: the code does not compile, so the verdict is unambiguously `real-bug` (no flake check, no master cross-check), and every minute spent waiting for user confirmation is a minute the rest of the patch keeps running doomed test tasks. The coordinator still prints a one-line notice naming the suspect branch and kind so the transition is visible. This auto-transition is **specific to compile fail-fast** — non-compile `actionable_failure` (terminal-failed patches with test failures) keeps the existing user-confirmation handshake.

**Relationship to the ≥10-task early-triage trigger.** The third Hard Rule 7 trigger (a running patch with ≥ `EARLY_TRIAGE_TASK_THRESHOLD` failed tasks; see `_detect_early_triage` in `stack_state.py`) shares the same `actionable_failure` decision value but does **not** emit a `compile-fail-fast:` line in `summary`. Because the line is absent, the coordinator keeps the existing "tell user; offer Phase 3" handshake (not auto-fix), and Phase 2 does not auto-abort descendants — the running patch's remaining tasks may still pass, so the user gets to make the call.

See [`references/triage.md` "Compile fail-fast triage"](references/triage.md) for what Phase 3 does when this fires.

---

## Phase 3 — `fix`: triage and repair actionable failures

0. **Tear down any active polling.** If Phase 2 was in cron mode, `CronList` → find the job whose prompt starts with `[evergreen-stack-ci-cron stack=<r>]` → `CronDelete` it. (Wakeup mode self-stops because no one re-arms it, so no cleanup needed there.) Skipping this risks the cron firing mid-fix and re-polling stale patches.
1. **Identify the earliest actionable branch** — the `earliest-actionable:` line in `stack_state.py summary`. (May differ from `earliest-red:` if the earliest red branch is excluded-only.) If `compile-fail-fast:` is also present, the suspect is that compile failure's branch, the kind is declared, and descendants have already been aborted in Phase 2 — Phase 3.7 will re-patch them along with the fix.
2. **Dispatch investigate-coordinator subagents in parallel — one per failed patch** (template 3a in [`references/subagent-investigate.md`](references/subagent-investigate.md), `model: "opus"`). Each 3a coordinator owns one patch end-to-end:
   - Reads `failed_tasks`, `build_failure`, and `polling.compile_fail_fast` from state.
   - **Short-circuits inline** (no further fan-out) when the patch is a compile-fail-fast suspect, has `build_failure=true`, or has only one failed task. In these paths the 3a coordinator itself reads the relevant logs, calls `record-failure` / `record-master-broken` per failing test (only for the single-task path — compile/build failures skip master cross-check entirely), and drafts the plan.
   - Otherwise **fans out to investigate-task subagents — one per failed task within the patch** (template 3b, `model: "opus"` + "think harder" cue), dispatched in a single message so they run in parallel. Each 3b worker pulls its one task's test results / logs, cross-checks master per failing test, calls `record-failure` / `record-master-broken`, and returns a structured `===TASK-INVESTIGATION-RESULT===` block of evidence pointers. The 3a coordinator aggregates the blocks (priority: real-bug > needs-retry > unknown > master-broken > flake) and drafts a unified implementation plan from the union of evidence.
   - Calls `set-findings` exactly once per patch with the aggregated verdict, suspect-branch, and plan.

   See [`references/subagent-investigate.md`](references/subagent-investigate.md) templates 3a and 3b for full prompt bodies, aggregation rules, and the 3b structured return contract.
3. **Read updated `summary`** — verdicts and counters are now visible.

    After dispatching all investigate-coordinator subagents and reading the updated summary, verify that every failed patch now has a non-null `findings` block. Do this by reading `stack_state.py show` and checking `branches[*].patches[-1].findings` for each RED branch.

    If any investigate-coordinator subagent completed but its patch's `findings` block is still null (the subagent hit sandbox restrictions or otherwise failed to call `set-findings`), the coordinator MUST call `set-findings` directly via its own Bash tool (`dangerouslyDisableSandbox=true`) using the verdict and cause from the subagent's return line. Do NOT proceed to step 4 (fix strategy) until all investigated patches have a persisted findings block.

    Why: without a findings block, `stack_state.py summary` cannot distinguish an unresolved failure from a classified flake, so flake patches appear as `actionable` and `earliest-actionable` points to the wrong branch, causing the coordinator to route fixes to the wrong place.
4. **Decide fix strategy** based on verdicts:
   - `flake` / `master-broken` → skip; surface to user.
   - `needs-retry` → re-run patch creation for that branch (Phase 1.5 subagent).
   - `real-bug` → dispatch a fix-and-commit subagent.
   - All remaining failures excluded → `summary` will already show `excluded_only`; loop will exit at next poll.
5. **Dispatch fix-and-commit subagent(s)** (template 4 in [`references/subagent-fix-submit.md`](references/subagent-fix-submit.md)). Linear stack with adjacent failures → fix the earliest only and let `gt restack` propagate. Sibling branches with independent failures → parallel via worktrees, see [`references/parallel-fixes.md`](references/parallel-fixes.md). Per Hard Rule 12, the fix subagent runs the plan's `Verify (build):` and `Verify (test):` steps before committing; the coordinator does not need to re-state this in the dispatch prompt. If a fix subagent returns `FAILED — verify-test: <reason>`, treat it as evidence of a wrong fix and re-dispatch investigation rather than re-running template 4 verbatim — the local failure is now part of the evidence the next 3a needs to see.
6. **Coordinator-only cleanup, detach, restack, re-attach.** Skip the whole step for `--single-branch` (no descendants → no restack needed).
    1. **Remove fix-worktrees only** — `git worktree list`; `git worktree remove <path>` for each *fix-worktree* (the per-fix scratch worktrees from [`references/parallel-fixes.md`](references/parallel-fixes.md)). Do NOT remove the per-branch CI worktrees here — they're session-scoped and we're about to detach-and-re-attach them around the restack.
    2. **Detach every per-branch CI worktree** so `gt restack` can rebase their branches. For each branch in `state.branches` whose worktree exists at `~/worktrees/<repo>-<branch>`:
       ```bash
       git -C ~/worktrees/<repo>-<branch> checkout --detach
       ```
       Skip the branch that's currently checked out in the main repo (it has no separate per-branch worktree). The detach releases the branch ref; the working tree stays put.
    2b. **PRE-RESTACK USER-WORKTREE AUDIT (blocking).** Step 6.2 has now detached every CI-created worktree, so the skill has done everything it can on its own side. Before `gt restack`, run `git worktree list --porcelain` and check whether any **user-owned** worktree (path NOT matching `~/worktrees/<repo>-<branch>`) still has one of the rebase-needed branches checked out. The rebase-needed set is the suspect branch + all its descendants in stack order.

       If any user-owned worktree still holds a rebase-needed branch, **STOP. Do NOT proceed to `gt restack`.** The skill must not detach or otherwise touch a user-owned worktree — the user may have uncommitted work, an in-progress rebase, or other state only they can safely resolve. Print the exact command(s) the user must run, one line per conflicting worktree:

       > "Before I can restack, please run:
       >   `git -C <worktree-path> checkout --detach`
       > (one line per conflicting user-owned worktree, in the order listed)
       > Then reply 'ready' and I will restack."

       Wait for explicit user confirmation. Do NOT continue until the user replies. Once they confirm, re-run `git worktree list --porcelain` to verify the conflicts are gone before proceeding to step 6.3. If a conflict remains, surface the still-conflicting paths and wait again.

       After step 6.3 completes successfully, remind the user how to re-attach (the skill will NOT re-attach user worktrees in step 6.4 — that set covers only the CI worktrees):

       > "Restack complete. To re-attach your worktree:
       >   `git -C <worktree-path> checkout -f <branch>`
       > (one line per worktree you detached above)"
    3. **`gt restack`** from the main checkout. Per Hard Rule 2 and the step 6.2b audit, no rebase-needed branch is held outside the main repo at this point — CI worktrees were detached in 6.2 and any conflicting user worktrees were resolved by the user in 6.2b. If restack still fails (e.g. a merge conflict in the rebase, or a user worktree was re-attached after the audit), surface the error to the user — do NOT silently re-attach. The detached CI worktrees can sit until the user resolves it.
    4. **Re-attach every detached worktree** to its (possibly moved) branch tip:
       ```bash
       git -C ~/worktrees/<repo>-<branch> checkout -f <branch>
       ```
       The `-f` discards any uncommitted/untracked junk left by the prior `evergreen patch` and moves HEAD + working tree to the branch's new tip in one shot. After this loop, every per-branch worktree is on its branch at the post-restack tip and ready for re-patching.
7. **Re-patch in parallel** — per Hard Rule 12, every fix that reaches this step has already passed a local canary `bazel test` (or explicitly declared `Verify (test): N/A`). If you arrived here with an un-verified fix — for example, because a fix subagent's return line was misread or the coordinator overrode the FAILED signal — stop and re-run template 4 with the canary target. Re-patching without local verification routinely burns full 30-min CI cycles on regressions that would have shown up in seconds locally. For the suspect branch + all descendants (NOT upstream), use the same default parallel-via-worktrees flow as Phase 1.5 (template 1a) so all replacement Evergreen builds kick off simultaneously. The per-branch worktrees were already re-attached to their post-restack tips in Step 6.4, so the create-patch subagents can `cd` straight into them. If a worktree directory is missing for any reason, recreate just that one with `git worktree add ~/worktrees/<repo>-<branch> <branch>` before dispatching its subagent. Never `git worktree remove`-then-`git worktree add` an existing worktree — that's the wasteful 77k-file recreate this skill is designed to avoid (see "Worktree lifecycle"). Honor `--sequential` if it was passed at start. Reuse `state.scope` so iterations are comparable.
8. **Loop back to Phase 2** — call `stack_state.py reset-poll-cycle --stack-root <r>` to start a fresh 12-iteration budget, then follow the Phase 2 entry steps (inline first poll + canonical wakeup prompt).

Repeat until `summary` reports `all_clean` or `excluded_only`.

---

## Phase 4 — `update-prs`: push updates to existing PRs only

After all patches are green and the user wants to update existing PRs:

1. **Enumerate PRs** —
   ```bash
   gh pr list --author "@me" --json number,headRefName --jq '.[] | .headRefName'
   ```
   Filter to branches in the state file's stack.
2. **Dispatch one submit-pr subagent per matching branch (sequential)** — template 5 in [`references/subagent-fix-submit.md`](references/subagent-fix-submit.md). Uses `gt submit --no-stack --update-only`. If a branch lacks a PR, the subagent returns SKIPPED — it does NOT create one.

If the user says "push the whole stack" or "create PRs for new branches", that's a separate explicit request — route to the `graphite-pr-stack` skill.

---

## Subagent coordination invariants

For investigation parallelism and worktree-based fix parallelism, see [`references/parallel-fixes.md`](references/parallel-fixes.md). Two strict invariants:

1. **All fix subagents commit before any worktree is removed.** Otherwise the commit may exist only in a removed worktree's reflog.
2. **All fix worktrees are removed before `gt restack`.** Restack walks all stack branches; conflicting checkouts make it fail with cryptic errors.

The main agent is the only restacker. Don't restack from a subagent.

---

## State file

Location: `$XDG_STATE_HOME/evergreen-stack-ci/<repo>--<stack-root>.json` (default `~/.local/state/...`). All writes go through `scripts/stack_state.py`. See [`references/state-schema.md`](references/state-schema.md) for the schema and `scripts/stack_state.py --help` for subcommands.

Key subcommands subagents use:
- `add-patch`, `update-status`, `set-findings`
- `record-failure`, `record-success` (pass `--patch-id` so the dashboard can show the green patch), `record-master-broken`
- `record-fix` (fix-and-commit subagent only — after every commit)
- `bump-poll-iteration` (poll-status subagent only — every wakeup)

Coordinator-only:
- `init`, `summary`, `quarantine`, `mark-completed-stack`, `rm`
- `reset-poll-cycle` (Phase 2 entry / re-entry from Phase 3)
- `schedule-next-poll --in-seconds N` (called immediately before each `ScheduleWakeup` so the dashboard can show a countdown to the next poll)
- `dashboard`, `dashboard-path` (HTML dashboard — see below)

## HTML dashboard

`stack_state.py` writes a self-contained HTML dashboard alongside the state file at `<repo>--<stack-root>.dashboard.html`. It uses meta-refresh (10s) so the user can open it once at session start and watch progress without any server, port, or daemon. Theme follows the MongoDB design system (DESIGN.md) — deep-teal dark surface with a darker header band and brand-green success accents.

The dashboard auto-regenerates after every mutating command (`init`, `add-patch`, `update-status`, `set-findings`, `record-failure`, `record-master-broken`, `record-success`) — there's nothing the coordinator or subagents need to do to keep it current. It's deleted by `rm` and `mark-completed-stack` so stale dashboards don't linger.

**Coordinator obligations:**
- After `init` in Phase 1, run `python3 scripts/stack_state.py dashboard-path --open --stack-root <root> --repo-root "$PWD"` — this prints the `file://` URL AND opens it in the user's default browser. The open step is best-effort (`webbrowser.open`, cross-platform) and never errors the command; set `STACK_STATE_NO_OPEN=1` in the environment to skip it for headless / SSH / CI sessions. Drop `--open` for read-only path lookups in later phases — the dashboard auto-refreshes every 10s, so the originally-opened tab stays current and re-opening on every poll would spawn dozens of tabs.
- If a render hiccups (warning on stderr from a mutating call), run `stack_state.py dashboard --stack-root <root>` to force-regenerate.

**Subagent obligations:** none. Subagents just call the same `stack_state.py` mutators they already use; the dashboard rides along.

**Opt-out:** set `STACK_STATE_NO_DASHBOARD=1` in the environment to skip auto-regeneration (mainly for tests).

---

## CLI conventions

- `evergreen` — invoke as just `evergreen`, never with a full path. Available on `$PATH`.
- `gt` — same. The Graphite CLI is `gt`, never a full path.
- `gh` — for checking PR existence.

For complete `evergreen patch` flag reference, see the project's `evergreen-cicd` skill. This skill assumes individual patch creation is understood — it's about orchestrating across a stack.

---

## Important notes

- **No new PRs without explicit user request.** Allergic to creating PRs. Route to `graphite-pr-stack` if asked.
- **Master-broken tests are excluded automatically.** Investigation cross-checks master; failures already red on master never trigger fixes. See [`references/three-strikes.md`](references/three-strikes.md).
- **Surface the thirdparty notice on every coordinator output.** When `state.scope.thirdparty_status` is `skipped-no-mapping` or `omitted` (and the profile is one that *could* have included thirdparty — `backend` or `full`), the coordinator MUST include a one-paragraph reminder banner in its user-facing output for *every* phase (start, status, fix, update-prs) AND a closing summary line at the end of the run. The banner names the suggested re-run command (e.g. `--thirdparty-teams=payments,billing` or `--thirdparty-teams=all`) and is informational only — never block the cycle on it. The banner is read from `stack_state.py summary`'s `thirdparty-notice:` line.
- **Cypress tests are flaky.** The `evergreen-cicd` skill flags `E2E_CYPRESS_*` as known-flaky. Default `--exclude=cypress` for `--profile=full` runs (in the rare case the user has explicitly opted in to `full` per Hard Rule 10).
- **State files accumulate.** Old state for completed stacks isn't auto-cleaned. Use `stack_state.py mark-completed-stack` (archive) or `rm` (delete) when done.
- **Patch creation is uncommitted-aware.** `evergreen patch -u` includes uncommitted changes. The default flow assumes branches are committed — surface uncommitted changes to the user (commit, stash, or include via `-u`).

---

## References

- [`references/test-scope.md`](references/test-scope.md) — profiles, exclusions, auto-detection
- [`references/subagent-create-poll.md`](references/subagent-create-poll.md) — subagent prompts: create-patch (1a/1b) + poll-status (2). Common contract + model-selection table inlined. **READ BEFORE dispatching a create-patch or poll-status subagent.**
- [`references/subagent-investigate.md`](references/subagent-investigate.md) — subagent prompts: investigate-coordinator (3a) + investigate-task (3b) + aggregation rules. Common contract + model-selection table inlined. **READ BEFORE Phase 3 investigation dispatch.**
- [`references/subagent-fix-submit.md`](references/subagent-fix-submit.md) — subagent prompts: fix-and-commit (4) + submit-pr (5). Common contract + model-selection table inlined. **READ BEFORE dispatching a fix or PR-update subagent.**
- [`references/three-strikes.md`](references/three-strikes.md) — quarantine, master-broken filtering, polling decision tree
- [`references/state-schema.md`](references/state-schema.md) — JSON schema and field reference
- [`references/triage.md`](references/triage.md) — identifying the earliest *causally* failing branch
- [`references/parallel-fixes.md`](references/parallel-fixes.md) — worktree + restack invariants for parallel fix dispatch
- [`references/workflow.md`](references/workflow.md) — concrete 4-branch walkthrough