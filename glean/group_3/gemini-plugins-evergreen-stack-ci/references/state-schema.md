# State file schema (v2)

## Location

`$XDG_STATE_HOME/evergreen-stack-ci/<repo>--<stack-root>.json`
(default: `~/.local/state/evergreen-stack-ci/...`)

State file is OUTSIDE the repo so it survives across worktrees, repo moves, and Claude Code sessions.

## Schema

```json
{
  "version": 2,
  "repo": "mms",
  "repo_root": "/Users/tom/worktrees/mms-1",
  "project_id": "mms",
  "mode": "stack",
  "stack_root": "feat_01_aws_entitlement_svc",
  "trunk": "master",
  "created_at": "2026-05-06T12:00:00Z",
  "scope": {
    "profile": "backend",
    "alias": null,
    "variants": ["unit_java", "int", "code_health"],
    "tasks": ["GENERATE_UNIT_TESTS_BAZEL", "GENERATE_INT_TESTS_BAZEL", "all"],
    "excluded": ["cypress"],
    "thirdparty_status": "skipped-no-mapping",
    "thirdparty_teams": [],
    "thirdparty_skipped_reason": "auto-detect found no team-specific path matches in stack diff"
  },
  "branches": [
    {
      "name": "feat_01_aws_entitlement_svc",
      "order": 0,
      "patches": [
        {
          "patch_id": "abc123",
          "url": "https://evergreen.mongodb.com/patch/abc123",
          "description": "feat_01 stack-ci",
          "created_at": "2026-05-06T12:00:00Z",
          "status": "succeeded",
          "checked_at": "2026-05-06T12:30:00Z",
          "failed_tasks": [],
          "failed_tests": [],
          "findings": null
        }
      ],
      "fixes": [
        {
          "commit_sha": "deadbeef1234",
          "summary": "fix null check in EntitlementController",
          "applied_at": "2026-05-06T12:45:00Z",
          "target_keys": ["feat_03::RestApiV1IntTests::testFoo"]
        }
      ]
    }
  ],
  "test_failures": {
    "feat_03::RestApiV1IntTests::testFoo": {
      "branch": "feat_03",
      "task": "RestApiV1IntTests",
      "test": "testFoo",
      "consecutive_failures": 0,
      "first_failed_at": "2026-05-06T12:20:00Z",
      "last_failed_at": "2026-05-06T13:05:00Z",
      "last_failed_patch": "def456",
      "last_round_id": null,
      "quarantined": false,
      "master_broken": false,
      "fixed_at": "2026-05-06T13:30:00Z",
      "fixed_in_patch": "ghi789",
      "time_to_fix_seconds": 4200
    }
  },
  "polling": {
    "iteration_count": 0,
    "cycle_started_at": null,
    "max_iterations": 12,
    "last_poll_at": null,
    "next_wakeup_at": null,
    "next_wakeup_seconds": null
  }
}
```

## Field reference

### Top level

| Field | Notes |
|---|---|
| `version` | Schema version. Currently `2`. Bumped when test_failures was added. |
| `mode` | `"stack"` or `"single"`. |
| `stack_root` | Branch closest to trunk; uniquely identifies the stack within a repo. |
| `scope` | See [`test-scope.md`](test-scope.md). Sticky — reused across re-patches. |
| `branches` | Ordered list (root → tip). |
| `test_failures` | Map keyed by `<branch>::<task>::<test>`. See [`three-strikes.md`](three-strikes.md). |
| `polling` | Per-cycle counter consumed by the wakeup-driven Phase 2 loop. See `polling.*` below. Auto-created by `bump-poll-iteration` / `reset-poll-cycle` if absent on legacy state files (no migration needed). |

### scope.thirdparty_*

| Field | Notes |
|---|---|
| `thirdparty_status` | One of: `included` (running with explicit teams), `auto-resolved` (auto-detect found teams), `all` (user requested everything), `skipped-no-mapping` (auto-detect found nothing — variant skipped, surface reminder), `omitted` (user didn't ask), `excluded` (`--exclude=thirdparty`). |
| `thirdparty_teams` | Resolved team list (e.g., `["payments", "billing"]`). Empty when status is `skipped-no-mapping`, `omitted`, or `excluded`. |
| `thirdparty_skipped_reason` | Free-text explanation when `thirdparty_status = skipped-no-mapping`. Used to render the persistent reminder banner. |

The `thirdparty-notice:` line printed by `summary` is a single-line, user-facing string derived from these three fields. The coordinator should include it verbatim in every phase output and in the final closing summary. When `thirdparty_status` is `auto-resolved`, `included`, or `all`, the notice line is absent (no warning needed).

### polling.*

| Field | Notes |
|---|---|
| `iteration_count` | Number of poll subagents dispatched in the current polling cycle. Incremented by `bump-poll-iteration`; zeroed by `reset-poll-cycle` (called at Phase 2 entry and on Phase 3 → Phase 2 hand-back). |
| `cycle_started_at` | ISO timestamp of the most recent `reset-poll-cycle`, or first `bump-poll-iteration` if reset was never called (legacy state). |
| `max_iterations` | Cycle cap (default `12`). The poll subagent's return line carries `iteration <i>/<max>` so the coordinator can compare against this without reading state. |
| `last_poll_at` | ISO timestamp of the most recent `bump-poll-iteration`. Cleared on `reset-poll-cycle`. Used by the dashboard to show "last polled <Xm ago>". |
| `next_wakeup_at` | ISO timestamp the coordinator expects the next ScheduleWakeup-driven poll to fire. Set by `schedule-next-poll`; cleared by the next `bump-poll-iteration` (the wakeup fired) or `reset-poll-cycle`. Used by the dashboard countdown. |
| `next_wakeup_seconds` | The delay value passed to the most recent `schedule-next-poll` (`300` for the default 5-minute cadence). Diagnostic only. |
| `compile_fail_fast` | Object, optional. Set by `cmd_summary` (first detection) and updated by `mark-fail-fast-aborted`. Cleared by `reset-poll-cycle`. Shape below. |

`compile_fail_fast` shape:

```json
{
  "triggered_at": "2026-05-08T14:32:01Z",
  "suspect_branch": "feature/foo",
  "suspect_patch_id": "abc123",
  "kind": "COMPILE_BAZEL",
  "descendants": [{"branch": "feature/bar", "patch_id": "def456"}],
  "aborts_dispatched": true
}
```

`kind` is one of `COMPILE_BAZEL`, `COMPILE_CLIENT_BAZEL`, `BUILD_FAILURE`. `aborts_dispatched` is the idempotency flag the poll subagent uses to avoid double-aborting descendants on subsequent polls. Cleared on `reset-poll-cycle` so Phase 2 re-entry after a fix starts with a fresh slate.

The polling block is deliberately separate from `scope`: scope is sticky across cycles, but the polling counter must reset between cycles so each fix → re-patch round gets a fresh 12-iteration budget.

### Patch entry

| Field | Notes |
|---|---|
| `patch_id` | Evergreen patch id (the segment after `/patch/` in the URL). |
| `status` | One of: `pending`, `started`, `succeeded`, `failed`, `aborted`. |
| `failed_tasks` | CSV of failed Evergreen task names — populated by `update-status`. |
| `failed_tests` | List of `{task, test, suspect_branch}` triples — populated by `record-failure` / `record-master-broken`. Used by `summary` to compute per-patch actionability. |
| `findings` | Investigation result: `{notes, cause, suspect_branch, verdict, recorded_at}`. Verdict is one of: `real-bug`, `flake`, `master-broken`, `needs-retry`, `unknown`. |
| `build_failure` | Boolean, optional. Set by `mark-build-failure` when the poll subagent detects a build failure on ANY test-runner task — `has_test_results=false` AND `total_test_count=0`. Covers `RUN_ALL_UNIT_JAVA_TESTS` (the original "No test results found — there may have been a build failure" case), `INT_JAVA_*`, `INT_JAVA_THIRDPARTY_*`, and any other generated test-runner task. See [`triage.md` "Build failure on a test-runner task"](triage.md). Distinct from `failed_tasks` containing `COMPILE_BAZEL`, which `_detect_compile_fail_fast` picks up by name match. Both flag the patch as a fail-fast candidate. |

### test_failures entry

| Field | Notes |
|---|---|
| `branch` | Suspect branch (where the bug was introduced — used as the dedup key). |
| `consecutive_failures` | Bumped by `record-failure`, reset by `record-success`. |
| `last_round_id` | Latest patch on the suspect branch at the time of the last increment. Dedups cascading failures investigated in parallel. |
| `quarantined` | `true` once `consecutive_failures >= 3`. Excluded from actionability. |
| `master_broken` | `true` if the test is also failing on master. Excluded from actionability. Set by `record-master-broken`. |
| `master_broken_evidence` | Free-text note about how it was confirmed (e.g. master patch id + date). |
| `fixed_at` | ISO timestamp set by `record-success` when the test was observed passing again. Survives subsequent failures so the dashboard can render a `regressed` badge (entry has `fixed_at` AND `consecutive_failures > 0`). |
| `fixed_in_patch` | Optional patch id where the test was first observed green (set when `record-success --patch-id` is passed). |
| `time_to_fix_seconds` | Computed at `record-success` time: `fixed_at - first_failed_at` in seconds. Used by the "Recently fixed tests" dashboard panel. |

### branches[].fixes entry

| Field | Notes |
|---|---|
| `commit_sha` | SHA recorded by the fix-and-commit subagent right after `git commit`. |
| `summary` | One-line description of the fix (typically the commit subject). |
| `applied_at` | ISO timestamp the entry was recorded. |
| `target_keys` | Optional list of `<branch>::<task>::<test>` keys this fix targets. Used to cross-reference fixes with the failures they address. |

## Latest-patch semantics

`summary` and the polling subagent only consult each branch's *latest* patch (the last entry in `patches[]`). Older patches are kept for audit/history but are not re-polled. After a fix, when a new patch is appended for a branch, the old entry is automatically superseded.

## What writes this file

`scripts/stack_state.py` is the only writer. Use it via subprocess from any subagent — never hand-edit the JSON. All writes are atomic via `tmp + os.rename`.

Subcommands that mutate state:
- `init` — create
- `add-patch` — append a patch entry
- `update-status` — set patch status + failed_tasks
- `set-findings` — investigation findings on a patch
- `record-failure` — append failed test + bump per-test counter (deduped by round)
- `record-success` — reset per-test counter and stamp `fixed_at` / `fixed_in_patch` / `time_to_fix_seconds` (pass `--patch-id` for the green-patch link)
- `record-master-broken` — flag a test as already-failing-on-master
- `record-fix` — append a fix-commit entry to `branches[].fixes[]` (fix-and-commit subagent calls after each commit)
- `bump-poll-iteration` — increment `polling.iteration_count` and stamp `polling.last_poll_at`; also clears `polling.next_wakeup_at` (subagent calls per wakeup)
- `reset-poll-cycle` — zero `polling.iteration_count`, set `cycle_started_at`, clear `last_poll_at` / `next_wakeup_at`, and **clear `polling.compile_fail_fast`** (coordinator calls at Phase 2 entry / re-entry)
- `schedule-next-poll --in-seconds N` — record `polling.next_wakeup_at = now + N` so the dashboard can show a countdown (coordinator calls immediately before each `ScheduleWakeup`)
- `mark-build-failure` — set `branches[].patches[].build_failure = true` for the build-failure-on-test-runner-task pattern: any failed task with `has_test_results=false` AND `total_test_count=0` (RUN_ALL_UNIT_JAVA_TESTS, INT_JAVA_*, INT_JAVA_THIRDPARTY_*, or any `_BAZEL` test-runner task). Poll subagent calls in template 2 step 5c.
- `mark-fail-fast-aborted` — set `polling.compile_fail_fast.aborts_dispatched = true` (idempotent; first call also persists the full event block) (poll subagent calls in template 2 step 5d after running `evergreen abort`)
- `summary` — best-effort writes `polling.compile_fail_fast` the first time fail-fast is detected (otherwise read-only)
- `mark-completed-stack` — archive
- `rm` — delete

Subcommands that only read:
- `path`, `show`, `list`, `quarantine`, `get-fail-fast-aborts`
