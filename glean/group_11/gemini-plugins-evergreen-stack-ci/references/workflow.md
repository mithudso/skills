# Workflow walk-through

A concrete example: a 4-branch Graphite stack on `mms`. The intent is to validate backend tests across the entire stack without burning hours on E2E noise.

```
master
 └── feat_01_aws_entitlement_svc        [PR #166551 OPEN]
      └── feat_02_extract_marketplace   [PR #166552 OPEN]
           └── feat_03_aws_entitlement  [no PR yet]
                └── feat_04_dryrun      [no PR yet]
```

## Setup

User: "Run evergreen across this whole stack, but skip the E2E and frontend stuff."

You:

1. **Detect scope** — `gt log short` → 4 branches in linear order. Scope = stack mode. Run `git worktree list` and verify no other worktrees have any of these branches checked out.
2. **Resolve test scope** — Confirm with user: `--profile=backend --exclude=thirdparty` (skips E2E, frontend, and slow third-party integrations). Profile resolves to `-v unit_java -v int -v code_health -t GENERATE_UNIT_TESTS_BAZEL -t GENERATE_INT_TESTS_BAZEL -t all`.
3. **Init state**:

   ```bash
   python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py init \
     --stack-root feat_01_aws_entitlement_svc \
     --branches feat_01_aws_entitlement_svc,feat_02_extract_marketplace,feat_03_aws_entitlement,feat_04_dryrun \
     --project-id mms \
     --repo-root "$PWD" \
     --profile backend \
     --variants unit_java,int,code_health \
     --tasks GENERATE_UNIT_TESTS_BAZEL,GENERATE_INT_TESTS_BAZEL,all \
     --exclude thirdparty
   ```

4. **Create patches in parallel via worktrees (default)** — coordinator sets up one worktree per branch (skipping the one already on the main checkout), then dispatches 4 create-patch subagents in a single message:
   ```bash
   git worktree add ~/worktrees/mms-feat_02_extract_marketplace feat_02_extract_marketplace
   git worktree add ~/worktrees/mms-feat_03_aws_entitlement     feat_03_aws_entitlement
   git worktree add ~/worktrees/mms-feat_04_dryrun              feat_04_dryrun
   # feat_01 is already checked out in the main repo — its subagent uses the main checkout
   ```
   Each subagent (template 1a) `cd`s into its worktree, runs `evergreen patch -p mms <resolved flags> -f -y -d "<branch> stack-ci"`, parses the patch URL, and calls `stack_state.py add-patch`. Concurrent `add-patch` writes are flock-serialized inside the script. All 4 Evergreen builds enter the queue at roughly the same instant.

   With `--sequential`, the coordinator instead runs `gt checkout` + `evergreen patch` + `add-patch` once per branch on the main checkout, restoring the original branch when done.
5. **Coordinator leaves the worktrees in place** — they're session-scoped and will be reused for any Phase 3 re-patching (see [SKILL.md "Worktree lifecycle"](../SKILL.md)). The coordinator tells the user the state file path, the dashboard URL (e.g. `file:///Users/.../evergreen-stack-ci/mms--feat_01_aws_entitlement_svc.dashboard.html` — get the path AND auto-open it in the user's default browser with `stack_state.py dashboard-path --open --stack-root feat_01_aws_entitlement_svc --repo-root "$PWD"`; set `STACK_STATE_NO_OPEN=1` to skip the open in headless sessions), and that there are 4 patches in flight. The dashboard auto-refreshes every 10 seconds and updates as Evergreen status changes propagate through `update-status`, `record-failure`, and `set-findings`.

## Polling

20 minutes later, user: "Status?"

```
mcp__evergreen__list_user_recent_patches_evergreen
```

Filter to the 4 tracked IDs. Update state:
- feat_01: succeeded
- feat_02: succeeded
- feat_03: failed (TASK: `RestApiV1IntTests` on variant `int`)
- feat_04: failed (TASK: `RestApiV1IntTests` AND `BillingControllerIntTests`)

Each `update-status` call writes the latest state.

## Triage

Earliest failing branch is feat_03. feat_04's `RestApiV1IntTests` failure is shared with feat_03 → likely the same root cause. feat_04's `BillingControllerIntTests` is novel — but only fix it after feat_03 is green and re-patched (feat_04's patch will be regenerated downstream of the fix anyway).

The coordinator dispatches two investigate-coordinator subagents in parallel (template 3a, `model: "opus"`) — one per failed patch:

- **feat_03's patch** has 1 failed task (`RestApiV1IntTests`) → 3a takes the **single-task short-circuit**: pulls test results, cross-checks master, calls `record-failure` (or `record-master-broken`) for each failing test, drafts a plan, and calls `set-findings`. No 3b fan-out — there's nothing to parallelize for a single task.
- **feat_04's patch** has 2 failed tasks (`RestApiV1IntTests`, `BillingControllerIntTests`) → 3a takes the **fan-out path**: dispatches 2 investigate-task subagents (template 3b, `model: "opus"` + "think harder" cue) in a single message. Each 3b worker pulls its one task's test results, cross-checks master per failing test, calls `record-failure` / `record-master-broken`, and returns a structured evidence block. The 3a coordinator aggregates: both tasks return `verdict=real-bug` → patch verdict=`real-bug`; the unified plan covers both files; the cause paragraph notes that `RestApiV1IntTests` shares its root cause with feat_03's failure (cascading) and `BillingControllerIntTests` is novel to feat_04.

```
mcp__evergreen__get_patch_failed_jobs_evergreen patch_id=<feat_03 patch>
mcp__evergreen__get_task_test_results_evergreen task_id=<RestApiV1IntTests task>
```

(These MCP calls happen inside the 3a coordinator for feat_03's single-task short-circuit, and inside each 3b worker for feat_04's fan-out — never inside the top coordinator.)

## Fix

Linear stack, single root cause → no parallelism needed.

```bash
gt checkout feat_03_aws_entitlement
# edit code to fix the test
```

Per [SKILL.md Hard Rule 12](../SKILL.md), run a local canary `bazel test` before committing. The plan that template 3a produced names a `Verify (build):` and `Verify (test):` target; the fix subagent runs both. Picking the smallest target + filter that exercises the changed code path is what catches @BeforeEach explosions and stubbing errors in seconds — mms Evergreen cycles are ~30+ min per patch, so the 1–3 min spent locally pays for itself many times over.

```bash
# Verify (build) — first sanity check that the changed files at least compile.
bazel build //server/integrations/aws-entitlement:aws-entitlement

# Verify (test) — canary run on the smallest target that exercises the fix.
# A whole class is fine when the bug shows up at @BeforeEach; otherwise use
# --test_filter to scope to one representative method.
# --strict_java_deps=error matches Evergreen's strict-deps mode (local default
# is warn), so a local pass means the CI deps graph also passes — without it
# an indirect-dep usage can pass locally and still fail Evergreen with
# `[strict] Using type ... from an indirect dependency`. See
# references/triage.md "Why local `bazel test` may pass when Evergreen fails".
bazel test //server/integrations/aws-entitlement:RestApiV1IntTests \
  --test_filter=RestApiV1IntTests.testHandlesNullEntitlementCode \
  --strict_java_deps=error \
  --test_output=errors
```

Only commit after `bazel test` is green. If it fails, fix the fix before committing and before restacking — never re-patch a known-broken local fix; the next CI cycle will burn another 30 minutes confirming what the local test already showed.

```bash
git add -A
git commit -m "fix: <jira>: handle null entitlement code in REST handler"
```

If no bazel test target covers the touched files (e.g. the diff is BUILD-file-only, generated code, or pure docs), the plan should have declared `Verify (test): N/A — <reason>`; in that case the fix subagent runs only `bazel build` and proceeds. See [SKILL.md Hard Rule 12](../SKILL.md) for the graceful-degradation contract.

## Restack and re-patch

The per-branch worktrees from Phase 1 are still alive and currently each hold one branch. `gt restack` would refuse to rebase those branches while they're held, so the coordinator brackets the restack with detach (release branch refs) and re-attach (move worktrees to the new tips). See [SKILL.md Phase 3 Step 6 + Hard Rule 2](../SKILL.md).

```bash
git worktree list   # verify no fix-worktrees were spawned (none in this case)

# 1. Detach: release branch refs while keeping the 77k-file working trees on disk.
#    feat_01 has no per-branch worktree (it's on the main checkout) — skip it.
git -C ~/worktrees/mms-feat_02_extract_marketplace checkout --detach
git -C ~/worktrees/mms-feat_03_aws_entitlement     checkout --detach
git -C ~/worktrees/mms-feat_04_dryrun              checkout --detach

# 2. Restack from the main checkout — propagates the fix into feat_04.
gt restack

# 3. Re-attach: move each worktree forward to its (possibly moved) branch tip and
#    discard any junk left by the previous evergreen patch run. -f handles both.
git -C ~/worktrees/mms-feat_02_extract_marketplace checkout -f feat_02_extract_marketplace
git -C ~/worktrees/mms-feat_03_aws_entitlement     checkout -f feat_03_aws_entitlement
git -C ~/worktrees/mms-feat_04_dryrun              checkout -f feat_04_dryrun

# 4. Re-patch from inside each affected worktree (suspect + descendants only —
#    feat_01 and feat_02 are upstream of the fix and don't need re-patching).
(cd ~/worktrees/mms-feat_03_aws_entitlement && \
   evergreen patch -p mms <same flags as before> -f -y -d "feat_03 stack-ci v2")
python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py add-patch \
  --stack-root feat_01_aws_entitlement_svc --branch feat_03_aws_entitlement \
  --patch-id <new id> --url <url>

(cd ~/worktrees/mms-feat_04_dryrun && \
   evergreen patch -p mms <same flags> -f -y -d "feat_04 stack-ci v2")
# add-patch for feat_04 with new id
```

feat_01 and feat_02 don't need re-patching — they're upstream of the fix. Their worktrees stay alive too (and were detach/re-attached along with the others), in case a later cycle surfaces an upstream failure.

## Loop

Poll again. If feat_04 still has the `BillingControllerIntTests` failure (now without the cascading RestApi failure), it's a real feat_04 bug — repeat the fix → restack → re-patch loop on feat_04.

## Tear down per-branch worktrees once the stack is green

When the next poll returns `decision: all_clean`, the coordinator removes every per-branch worktree it created in Phase 1 *before* offering Phase 4 (per [SKILL.md "Worktree lifecycle"](../SKILL.md)):

```bash
git worktree remove ~/worktrees/mms-feat_02_extract_marketplace
git worktree remove ~/worktrees/mms-feat_03_aws_entitlement
git worktree remove ~/worktrees/mms-feat_04_dryrun
# feat_01 had no worktree (it was on the main checkout) — nothing to remove there
```

Each `git worktree remove` is idempotent — skip silently if the path is already gone, and pass `--force` if removal complains about untracked files left by `evergreen patch`. Surface a one-line confirmation: `removed 3 stack-CI worktrees`.

The same teardown happens on `excluded_only` and on explicit user stop. It does NOT happen on `actionable_failure` — Phase 3 will reuse those worktrees.

## Push to existing PRs

Once everything is green and worktrees are torn down:

```bash
gt checkout feat_01_aws_entitlement_svc
gt submit --no-stack --update-only

gt checkout feat_02_extract_marketplace
gt submit --no-stack --update-only

# feat_03 and feat_04 do NOT have PRs — skip them
```

Confirm with the user before doing anything for branches without PRs. Default behavior: don't push, don't create.

## Archive

Once the stack lands:

```bash
python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py mark-completed-stack \
  --stack-root feat_01_aws_entitlement_svc --repo-root "$PWD"
```
