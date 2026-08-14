# Failure triage: finding the earliest *causal* branch

The naive heuristic — "earliest red branch caused it" — is right most of the time but not always. Each branch's patch contains all changes since master, so a failure can:

1. **Originate at branch N** and cascade to N+1, N+2, ...
2. **Be environmental** (Evergreen-side flake, infrastructure issue) and surface on a random subset of patches.
3. **Be flaky** (test passes sometimes, fails sometimes) and surface on patches that just happened to lose the dice roll.
4. **Be masked by a later branch** (rare — branch N introduces bug, branch N+2 happens to remove or work around it).

## The decision procedure

For each failing patch, get the failure set:

```python
failures(p) = { (variant, task, test) for each red item in patch p }
```

Walk branches root → tip:

```
For each branch B in order:
    If failures(B.patch) is empty: skip
    earliest_with_failure = B
    break

For each branch C downstream of earliest_with_failure:
    novel_failures(C) = failures(C.patch) - union(failures(D.patch) for D <= C in earlier branches)
```

The branch where a failure *first appears* is the suspect for that failure. If a test fails on branch 2 and on branches 3, 4, 5, branch 2 is the suspect — until proven otherwise.

## Build failure on a test-runner task ("No test results found" pattern)

When MCP `get_patch_failed_jobs` returns a failed task with:
- `has_test_results: false`
- `failed_test_count: 0`
- `total_test_count: 0`
- `failure_details.description` containing `"Build failed"`, OR
- `failure_details.failing_command` matching `'.*run bazel.*'` OR `'shell.exec'.*step \d+ of 10`

…it is a BUILD/COMPILE failure that happened **before any test could run**. The task simply never reached the test phase. This applies to:
- `RUN_ALL_UNIT_JAVA_TESTS` (the original case — "No test results found — there may have been a build failure")
- `INT_JAVA_*` (all INT_JAVA_* variants — billing, payments, paymentview, etc.)
- `INT_JAVA_THIRDPARTY_*`
- Any `_BAZEL` task or other generated test-runner task

**Do NOT classify as flake, master-broken, or test failure based on the absence of test results.** And do NOT trust a task's own pass/fail status — for `RUN_ALL_UNIT_JAVA_TESTS` in particular, the task can show as succeeded even when the build it owns failed for a specific team's outputs. The signal is `has_test_results=false` + `total_test_count=0`, not the task's status.

### When to suspect this BEFORE pulling logs

Heuristic: `total_test_count == 0` AND `has_test_results == false` → ALWAYS treat as build failure first; only fall back to flake/master-broken if logs confirm tests actually ran. Misclassifying a build failure as a test failure burns entire CI cycles on phantom fixes (see [SKILL.md Hard Rule 13](../SKILL.md)).

### Procedure

1. Call `mcp__evergreen__get_task_logs_evergreen` with `filter_errors=false` — **NOT the default `true`**. The `filter_errors=true` setting drops bazel-formatted error lines like `[strict]` and the `** Please add the following dependencies` blocks because they don't match the regex it uses for stderr-style errors. ECJ `1. ERROR in ...` lines are similarly missed. Without the override you will see an empty / misleading log and miss the actual error.
2. Scan the log in two passes — stop at the first pass that yields a real signal:

   **Pass 1 — `ERROR:` scan (try this first).**
   - Case-insensitively grep for the literal token `ERROR:` (e.g. `rg -in 'error:' <log>`).
   - For each match, capture the **full multi-line block**: the `ERROR:` line plus every subsequent line that carries the **same originating timestamp** as the `ERROR:` line. Evergreen / java loggers emit a single multi-line event (stacktrace, `Caused by:` chain, context) under one timestamp; that whole block is the error.
   - **Discard dependency-download red herrings.** Drop any block whose subject is fetching / resolving artifacts — common shapes: `Downloading from`, `Failed to download`, `Could not resolve dependencies`, `Could not transfer artifact`, `Failed to fetch ... from https://...`, bazel remote-cache fetch errors, gradle `Could not GET ...`. These are transient worker-network issues and are not what broke the patch. The rule is intent-based: drop blocks whose subject is fetching artifacts, not blocks whose text merely contains "download" — an `ERROR: Downloading job failed because user IAM grants prohibit ...` business-logic error is NOT dropped.
   - If any blocks survive the filter, those are the error signatures. Use the most relevant one (usually the earliest non-download block) for `record-failure` and stop — do not run Pass 2.

   **Pass 2 — bazel/ECJ pattern fallback (only if Pass 1 found nothing after filtering).**
   - `\[strict\] Using type .* from an indirect dependency` → missing bazel deps in a test BUILD.bazel
   - `ERROR in .*\.java \(at line \d+\)` → ECJ compile error (syntax, merge conflicts, missing import, broken signature)
   - `package .* does not exist` → missing import or missing dep
   - `cannot find symbol` → missing class, deleted method, refactor stale
   - `error: \[NullAway\]` → null-safety violation
   - `BUILD FAILED` / `Compilation failed` → upstream compile failure

   Why two passes and not one combined grep: the `ERROR:` scan is higher-recall — it catches application/integration test errors, gRPC stack traces, container exit messages, and anything else a logger emits — while the bazel patterns are higher-precision for bazel-side compile failures. The download-exclusion only matters for Pass 1 because the Pass 2 patterns are already narrow enough to never match download chatter.
3. The failing build target is named in the error (e.g. `//server/src/test/.../package:TestLibrary failed`). The fix is in **that target's BUILD.bazel**, not in the test or production .java file (unless the error itself is a compile error in a .java file).
4. For `[strict]` errors specifically: the error message names the exact `buildozer 'add deps ...'` command to run. Use it verbatim. Do NOT guess deps — the error tells you exactly what's missing.
5. Attribute the build error to its causal branch using the same novel-failure walk as for test failures (the offending diff is what broke compilation). Compile errors are **branch-local** — skip the master cross-check entirely.
6. Record via `record-failure` with `--task <failing-task-name> --test <short-error-signature>` so the three-strikes counter still applies if the build error recurs. Do NOT record as `master-broken` unless master itself fails to compile in the same way — a build error introduced by stack diffs is a real bug, not master rot.
7. Call `mark-build-failure --patch-id <id>` on the state file so the dashboard reflects this as a build failure, not a test failure. (Subcommand already exists in `scripts/stack_state.py`; the poll subagent calls it during template 2 step 5c whenever any failed task matches the signature above.)

### Where this runs in the 3a/3b split

The pattern is detected at **poll time** — template 2 step 5c flips `patch.build_failure=true` whenever **any** failed test-runner task matches the signature above (not just `RUN_ALL_UNIT_JAVA_TESTS`). At investigation time, the 3a coordinator reads that flag and handles the build failure **inline** (no 3b fan-out): identify which failed task(s) have `has_test_results=false`, pull THEIR logs with `filter_errors=false`, extract an error signature, attribute via the novel-failure walk, call `record-failure` once, and `set-findings --verdict real-bug`. If a 3b worker ever sees an `has_test_results=false` task — which should be pre-empted by the 3a short-circuit but isn't guaranteed — it bails out with `verdict=unknown` and a cause naming the pattern so the 3a coordinator can take over.

## Why local `bazel test` may pass when Evergreen fails

Local bazel runs with `--strict_java_deps=warn` by default on developer machines and benefits from a hot transitive-deps cache. Evergreen runs with `--strict_java_deps=error` (CI default). The direct consequences:

- A test that locally compiles via an **indirect** dep will compile clean locally but fail Evergreen with `[strict] Using type ... from an indirect dependency`. The fix is to add the direct dep to the test target's BUILD.bazel (the error names the exact `buildozer` command).
- A merge conflict marker that an IDE warned about may build locally if the file isn't touched in the current change, but the marker still breaks Evergreen's clean build.
- A newly-introduced test that pulls in a class only available transitively will pass `bazel test` locally and fail on Evergreen.

**Always reproduce with the strict flag before declaring a local pass** — this is mandated by [SKILL.md Hard Rule 12](../SKILL.md):

```
bazel test //path:target --test_filter=<canary-test-method> --strict_java_deps=error --test_output=errors
```

If the strict flag surfaces unrelated pre-existing errors in test infra (`KinesisClientConfig`, `SoftAssertionsExtension`, etc. — common in `mms`), document them in the plan's `Verify (test):` step and build the specific touched targets in isolation with the strict flag instead of the whole graph. Do NOT remove the flag to make the canary pass — that's the trap this section exists to close.

## Compile fail-fast triage

The fail-fast trigger lives on the **poll path**, not the triage path — by the time Phase 3 starts, the suspect is already named in `compile-fail-fast: <branch> kind=<kind>` from `summary`, descendants have been `evergreen abort`'d, and `polling.compile_fail_fast` is persisted in the state file. Triage doesn't re-discover any of that; it just pulls the right logs and records.

When the `compile-fail-fast:` line is present in `summary`:

1. **Pull the suspect patch's compile logs with `filter_errors=false`.** For `kind=COMPILE_BAZEL` or `kind=COMPILE_CLIENT_BAZEL`, fetch the failing task's logs via `mcp__evergreen__get_task_logs_evergreen` (the task name is in the patch's `failed_tasks` field). For `kind=BUILD_FAILURE`, pull logs from **whatever test-runner task the poll subagent flagged** — that's whichever failed task had `has_test_results=false` AND `total_test_count=0` (could be `RUN_ALL_UNIT_JAVA_TESTS`, `INT_JAVA_*`, `INT_JAVA_THIRDPARTY_*`, or any `_BAZEL` task — see "Build failure on a test-runner task" above). In every case pass `filter_errors=false` — the default `true` drops bazel `[strict]`, `** Please add ...`, and ECJ `1. ERROR in ...` lines.
2. **No master cross-check.** Compile errors are inherently branch-local — "the code in this commit doesn't compile against this base". Master health isn't the question; the diff on this branch is. Skip the `record-master-broken` path entirely.
3. **Record once.** A single `record-failure` call with `--task <failing-task> --test <short-error-signature>` (or `--test compile` / `--test build` if no concrete signature is available) so the failure shows up on the dashboard and the three-strikes counter still applies if the same compile error recurs across cycles.
4. **Set findings.** `set-findings --verdict real-bug --suspect-branch <suspect>` — fail-fast cases are real bugs by definition (the code doesn't compile). For `[strict]` errors, the plan should use the `buildozer 'add deps ...'` command verbatim from the error message — do NOT guess deps.

**In the 3a/3b split.** The 3a coordinator handles compile fail-fast inline — there's exactly one compile task to investigate, the verdict is unambiguously `real-bug`, and master cross-check is irrelevant. No 3b workers are dispatched in this case; the four steps above all run inside the 3a coordinator's own context.

Phase 3.7 re-patch covers the suspect plus descendants exactly as in the normal flow, because the descendants were aborted in Phase 2 and have no green latest patch.

## Cross-checking before fixing

Investigation subagents perform these checks per failing test, then either record-failure or record-master-broken (see [`three-strikes.md`](three-strikes.md)). Tests that fail on master are **structurally excluded** from the fix loop — the skill never tries to repair them.

1. **Does the failing task / test pass on master?**
   - Use `mcp__evergreen__list_user_recent_patches_evergreen` filtered to recent master patches, or check the project's master health board.
   - If master is also red on `(task, test)`: call `stack_state.py record-master-broken --branch <suspect> --task <T> --test <test> --evidence "<note>"`. Do NOT also call record-failure for that test.
   - This is the primary mechanism for "tests failing on the base branch should be ignored".
2. **Is this a known flake?**
   - `E2E_CYPRESS_*` tasks are pre-flagged as flaky in the project's `evergreen-cicd` skill. Don't auto-fix; ask first.
   - If the failure looks like a flake (timing assertion, network call, unrelated to the diff), retry once via `evergreen patch --repeat-failed -i <patch_id>` before assuming it's a bug.
3. **Does the diff between the suspect branch and its parent actually plausibly cause this failure?**
   - `git diff <parent>..<suspect> --stat` — does the changed area touch anything related to the failing test?
   - If unrelated → could be flake or env issue. Retry before fixing.
4. **Who does the cross-check in the 3a/3b split.** In the fan-out path (multi-task patches that aren't compile/build failures), each per-task investigator (template 3b) performs the loop above for its single task's failing tests and writes `record-failure` / `record-master-broken` itself. In the three short-circuit paths (compile-fail-fast, build-failure, single-task), the 3a coordinator runs the loop inline — there's only one task to investigate, so fan-out is wasteful. Either way the per-test writes happen at the lowest level that has the test-level evidence; the top coordinator never sees raw MCP output.

## Parallel-stack edge case

In a non-linear stack:

```
master
 └── feat_a
      ├── feat_b
      └── feat_c
```

`feat_b` and `feat_c` are siblings. If `feat_b` fails a test, `feat_c` may not — the test failure is on `feat_b` only because the diff is on `feat_b`. There's no causal cascade. Investigate them independently.

`gt log` shows the actual parent/child relationships — use that, not branch ordering, to determine ancestry.

## Output of the triage step

For each unique novel failure, produce a short triage line:

```
- task=<task>  test=<test>  earliest=<branch>  status=<NEW|FLAKE|MASTER_BROKEN>  fix-strategy=<commit on earliest | retry | skip + escalate>
```

This is what feeds into the fix phase.
