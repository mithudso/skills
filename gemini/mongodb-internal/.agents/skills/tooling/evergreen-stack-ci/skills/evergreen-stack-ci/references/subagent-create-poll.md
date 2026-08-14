# Subagent prompts — create-patch and poll-status

**Before dispatching:** set the Agent tool's `model:` argument per the Model selection table below, substitute the placeholders (`<state-file-path>`, `<stack-root>`, etc.), and pass the template's fenced block verbatim as the subagent prompt.

The coordinator never calls Evergreen, MCP, or git directly for any task that produces meaningful tool output. Every meaningful step runs in a subagent that reads the state file, does its job, writes results back via `stack_state.py`, and returns one short line. The coordinator only orchestrates from `stack_state.py summary` output. See SKILL.md "Context isolation" for the why.

## Common contract

Every subagent prompt includes:
1. **The state file path** (from `stack_state.py path --stack-root ...`).
2. **The skill path** so the subagent can read SKILL.md and the relevant reference.
3. **Its narrow scope** — exactly one branch / one patch / one task.
4. **Which `stack_state.py` subcommand to call when done**.
5. **A return contract:** ONE LINE of summary text. No log dumps. Verbose data goes into the state file.

## Model selection (REQUIRED)

The coordinator MUST pass the `model` argument when invoking the Agent tool for each subagent template. Investigation work — both cross-patch aggregation (3a) and per-task classification (3b) — gets `opus`. Code editing (template 4) gets `sonnet`. Pure I/O subagents inherit the coordinator's default.

| Template | Subagent | `model` | Rationale |
|---|---|---|---|
| 1a / 1b | create-patch | inherit (no override) | Mechanical: shell out to `evergreen patch`, parse one URL, call one `add-patch`. No reasoning required. |
| 2 | poll-status | inherit (no override) | Mechanical: MCP polls + state writes. The decision logic lives in `stack_state.py summary`, not the subagent. |
| 3a | investigate-coordinator (one per failed patch, parallel) | **`opus`** | Aggregates per-task evidence into a patch-level verdict and drafts the unified implementation plan template 4 will execute. Also handles the three inline short-circuits (compile-fail-fast, build-failure, single-task patches) without fanning out. The cross-task synthesis — spotting shared root causes across failed tasks, ordering plan steps by dependency — is the part that has to be on opus; the per-task log reading underneath is delegated. |
| 3b | investigate-task (one per failed task within a patch, parallel) | **`opus`** + "think harder" cue | Reads one task's test results / logs, cross-checks master per failing test, calls `record-failure` / `record-master-broken`, and returns a compact structured evidence block. Narrow scope by design — each 3b worker only sees its own task. The flake-vs-real-bug classification is load-bearing: a wrong verdict poisons the 3a aggregate and drives template 4 to commit a fix for the wrong thing, burning a 30-min CI cycle per mistake. That cost dominates per-task token cost, so 3b runs on opus and the "think harder" cue elevates extended thinking for the subtle calls. |
| 4 | fix-and-commit | **`sonnet`** + "think hard" | Implementation: apply the plan produced by template 3a, commit, record-fix. Sonnet is the right tier for code editing; include the `think hard` trigger word in the dispatch prompt so reasoning effort is elevated for tricky fixes. |
| 5 | submit-pr | inherit (no override) | Mechanical: `gt submit --no-stack --update-only`. No reasoning. |

The model is set on the Agent tool call itself, not inside the subagent prompt body — the coordinator passes `model: "opus"` for investigation (3a, 3b) or `model: "sonnet"` for code editing (4) as an argument when dispatching. The "think harder" / "think hard" cues for templates 3b and 4 live inside the prompt bodies (see `subagent-investigate.md` and `subagent-fix-submit.md`).

If a future subagent type is added, default to `inherit` unless it requires open-ended classification or planning — those tip into `opus`.

## 1. create-patch (one per branch)

Two variants. The default flow is **1a (parallel via worktree)** — all create-patch subagents are dispatched in a single message so every Evergreen build kicks off at roughly the same instant. Sequential **1b** is the `--sequential` opt-out (or for `--single-branch` runs where there's nothing to parallelize). Both run during Phase 1.5 and again after each fix in Phase 3.7. Both write through `stack_state.py add-patch`, which flock-serializes the writes so concurrent calls from 1a are race-free.

### 1a. Parallel via worktree (default)

```
You are a worker subagent for the evergreen-stack-ci skill.

Skill:    ~/.claude/skills/evergreen-stack-ci/SKILL.md
State:    <state-file-path>
Branch:   <branch-name>
Stack-root: <stack-root>
Worktree: <absolute path the coordinator already created via `git worktree add`>

Task — create one Evergreen patch for the given branch from its dedicated worktree:
1. `cd <worktree>` (the coordinator created the worktree on <branch>; do NOT git checkout — you'd contend with sibling subagents)
2. Read the state file's `scope` field to get the resolved variants/tasks/alias
3. Run: `evergreen patch -p <project_id> <resolved-flags> -f -y -d "<branch> stack-ci"`
4. Parse the patch URL from stdout (format: https://evergreen.../patches/<id>); extract <id>
5. Run (from anywhere — the script writes to the state file outside the repo):
     python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py add-patch \
       --stack-root <stack-root> --branch <branch> \
       --patch-id <id> --url <url> --description "<branch> stack-ci"
6. Do NOT remove the worktree — the coordinator owns cleanup after every sibling subagent has returned.

Return ONE LINE: "<branch>: patch <id> created" or "<branch>: FAILED — <reason>".
Do not include CLI output. If anything fails, return the failure reason — do not retry.
```

The coordinator dispatches every 1a subagent in a single message so they run concurrently. Worktrees are **session-scoped** — the coordinator does NOT remove them after Phase 1; they're reused for Phase 3 re-patching and only torn down once the cycle reaches `all_clean` / `excluded_only` or the user stops. See SKILL.md "Worktree lifecycle" for the full table. State writes are race-free because `stack_state.py` flocks the state file during `add-patch`.

### 1b. Sequential on the main checkout (`--sequential`)

```
You are a worker subagent for the evergreen-stack-ci skill.

Skill:  ~/.claude/skills/evergreen-stack-ci/SKILL.md
State:  <state-file-path>
Branch: <branch-name>
Stack-root: <stack-root>

Task — create one Evergreen patch for the given branch from the main checkout:
1. `gt checkout <branch>` (or `git checkout <branch>`)
2. Read the state file's `scope` field to get the resolved variants/tasks/alias
3. Run: `evergreen patch -p <project_id> <resolved-flags> -f -y -d "<branch> stack-ci"`
4. Parse the patch URL from stdout (format: https://evergreen.../patches/<id>); extract <id>
5. Run:
     python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py add-patch \
       --stack-root <stack-root> --branch <branch> \
       --patch-id <id> --url <url> --description "<branch> stack-ci"

Return ONE LINE: "<branch>: patch <id> created" or "<branch>: FAILED — <reason>".
Do not include CLI output. If anything fails, return the failure reason — do not retry.
```

The coordinator dispatches 1b subagents one at a time and restores the user's original branch after the last one completes.

## 2. poll-status (single subagent per polling iteration)

One subagent per poll cycle covers the whole stack. Pulls fresh status for every tracked patch via MCP, writes results back, **bumps the polling iteration counter**, **reads the poll-decision from `summary`**, and returns a one-line rollup that embeds both. The coordinator never calls MCP, never reads `summary`, and never reloads SKILL.md on a wakeup — the subagent's return line is sufficient to drive the decision tree in SKILL.md's Phase 2.

```
You are a worker subagent for the evergreen-stack-ci skill.

Skill:  ~/.claude/skills/evergreen-stack-ci/SKILL.md
State:  <state-file-path>
Stack-root: <stack-root>

Task — poll Evergreen for the latest status of every tracked patch, bump the polling
counter, and return a single line that carries the next-action decision:

1. Read the state file. Collect all unique patch IDs from each branch's *latest* patch entry
2. Call mcp__evergreen__list_user_recent_patches_evergreen to get fresh status for those IDs
3. For each tracked patch, determine: pending | started | succeeded | failed | aborted
4. For each patch in `failed` state, call mcp__evergreen__get_patch_failed_jobs_evergreen
   to enumerate failed task ids/names (do NOT pull task logs here — that's investigation)
5. For every patch, write status back via:
     python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py update-status \
       --stack-root <stack-root> --patch-id <id> --status <status> \
       --failed-tasks "<csv of failed task names if failed, empty otherwise>"
   Note: `stack_state.py summary` automatically detects the ≥10-failed-task
   early-triage condition on a still-`started` patch from the `failed_tasks`
   you just wrote (via `_detect_early_triage`, threshold lives in the
   module-level `EARLY_TRIAGE_TASK_THRESHOLD` constant). No extra subagent
   step is needed — the threshold check happens inside `cmd_summary`, so the
   `decision: <D>` you read back in step 7 already reflects it. Always pass
   `--failed-tasks` even for `started` patches with partial results so the
   early-triage check can fire mid-flight.
5b. For each patch that transitioned to `succeeded` this round, look at the **previous**
    latest patch on the same branch (from the state file you read in step 1) and replay each
    of its failed_tests entries through:
       python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py record-success \
         --stack-root <stack-root> --branch <suspect-branch> \
         --task <task> --test <test> --patch-id <new-green-patch-id>
    This stamps `fixed_at` / `fixed_in_patch` so the dashboard's "Recently fixed tests"
    panel can show the win. Skip this step for patches that didn't transition (still
    pending/started/failed/aborted). If the branch had no prior failed_tests, no calls.
5c. Build-failure inspection — for ANY failed task with
    `has_test_results=false` AND `total_test_count=0` (visible in the
    mcp__evergreen__get_patch_failed_jobs_evergreen response you fetched in
    step 4), call mcp__evergreen__get_task_logs_evergreen with
    `filter_errors=false` on the failing task and apply the two-pass scan
    from references/triage.md "Build failure on a test-runner task" step 2:
    FIRST a case-insensitive `ERROR:` grep with same-timestamp multi-line
    capture, dropping dependency-download blocks (transient artifact-fetch
    red herrings — not what broke the patch); ONLY IF the `ERROR:` pass
    yields nothing after filtering, fall back to the bazel/ECJ pattern list
    (`[strict]` indirect-dep, `ERROR in ...`, `cannot find symbol`,
    `package ... does not exist`, NullAway, `BUILD FAILED`).
    Pass `filter_errors=false` explicitly — the default `true` drops bazel-
    formatted error lines because they don't match the regex it uses for
    stderr-style errors, so you'd see a misleading empty log.
    If any of those patterns match:
      python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py mark-build-failure \
        --stack-root <stack-root> --patch-id <id>
    This covers RUN_ALL_UNIT_JAVA_TESTS (the original case), INT_JAVA_*
    (billing, payments, paymentview, …), INT_JAVA_THIRDPARTY_*, and any
    other test-runner task — NOT just RUN_ALL_UNIT_JAVA_TESTS. Heuristic
    per SKILL.md Hard Rule 13: `has_test_results=false` + `total_test_count=0`
    is ALWAYS a build failure first; only fall back to flake/master-broken
    if logs confirm tests actually ran.
    (COMPILE_BAZEL / COMPILE_CLIENT_BAZEL failures need no extra MCP call —
    stack_state.py detects them from the failed_tasks names you already wrote
    in step 5.)
5d. Compile fail-fast aborts — pick up any descendants that the latest
    summary computation flagged for abort:
      TO_ABORT=$(python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py \
                   get-fail-fast-aborts --stack-root <stack-root>)
      if [ -n "$TO_ABORT" ]; then
        for pid in $TO_ABORT; do
          evergreen abort -p "$pid"  # no-op if already terminal
        done
        python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py \
          mark-fail-fast-aborted --stack-root <stack-root>
      fi
    Idempotent: mark-fail-fast-aborted flips a flag so subsequent polls
    return nothing from get-fail-fast-aborts.
6. Bump the polling iteration counter (one call, captures `iteration: <i>/<max>`):
     python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py bump-poll-iteration \
       --stack-root <stack-root>
   Capture the printed `iteration: <i>/<max>` — keep <i> and <max> for the return line.
7. Read just the poll-decision (do NOT dump full summary into your output):
     python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py summary \
       --stack-root <stack-root> | rg "^poll-decision:"
   Capture <D> from "poll-decision: <D>". Valid values:
   actionable_failure | in_progress | excluded_only | all_clean | needs_attention.

Return ONE LINE in this exact format:
   "<succ>/<total> green, <failed> red, <running> running | iteration <i>/<max> | decision: <D>"

Do not enumerate per-patch results. Do not dump summary or MCP output. The coordinator
relies on the `decision: <D>` token to choose the next step — never omit it.
```

---

Back to [SKILL.md](../SKILL.md).
