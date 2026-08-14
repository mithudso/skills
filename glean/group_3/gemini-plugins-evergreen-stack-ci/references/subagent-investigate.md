# Subagent prompts — investigate-coordinator (3a) and investigate-task (3b)

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

The model is set on the Agent tool call itself, not inside the subagent prompt body — the coordinator passes `model: "opus"` for investigation (3a, 3b) or `model: "sonnet"` for code editing (4) as an argument when dispatching. The "think harder" / "think hard" cues for templates 3b and 4 live inside the prompt bodies (see below for 3b; see `subagent-fix-submit.md` for 4).

If a future subagent type is added, default to `inherit` unless it requires open-ended classification or planning — those tip into `opus`.

## 3a. investigate-coordinator (one per failed patch, parallel)

Read-only and independent → safe to dispatch in parallel across patches. Each 3a subagent gets exactly one failed patch and owns the patch-level verdict + implementation plan.

**Dispatch with `model: "opus"`.** This subagent decides whether to investigate the patch inline or fan out to per-task workers (template 3b), aggregates per-task evidence into a patch-wide verdict, and drafts the implementation plan template 4 will execute. The cross-task synthesis (spotting shared root causes, ordering plan steps by dependency) runs here on opus; the per-task log reading and flake-vs-real-bug classification is delegated to 3b workers, which also run on opus because that judgment is the load-bearing input to this aggregation.

Three inline short-circuits skip the fan-out — they save a layer of dispatch overhead when there is no parallelism to gain or the verdict is already determined:

1. **Compile fail-fast** (`polling.compile_fail_fast.suspect_patch_id == <patch-id>`). The poll subagent already named the suspect and aborted descendants; the verdict is unambiguously `real-bug` and master cross-check is irrelevant. Pull the compile task's logs and record once.
2. **Build-failure pattern** (`patch.build_failure == true`). The poll subagent already detected a build failure on some test-runner task — `has_test_results=false` AND `total_test_count=0` on RUN_ALL_UNIT_JAVA_TESTS, INT_JAVA_*, INT_JAVA_THIRDPARTY_*, or any other generated test-runner task (see `references/triage.md` "Build failure on a test-runner task"). There is one log to read (the failing task's, with `filter_errors=false`) and master cross-check is again irrelevant (the diff didn't compile against this base).
3. **Single failed task** (`len(failed_tasks) <= 1`). Nothing to parallelize. Do the per-test loop inline.

Otherwise — multi-task patches that aren't compile/build failures — fan out one 3b worker per task in a single message, then aggregate.

The 3a coordinator is responsible for exactly one piece of state file bookkeeping in every path: **a single `set-findings` call per patch**, stamped with the aggregated verdict, suspect branch, and (if `verdict=real-bug`) the unified implementation plan. The per-test `record-failure` / `record-master-broken` writes happen at the level that has the evidence — 3b workers in the fan-out path, the 3a coordinator itself in the inline paths.

```
You are a worker subagent for the evergreen-stack-ci skill (template 3a:
investigate-coordinator). You own ONE failed patch — picking the investigation
strategy (inline short-circuit vs. fan-out to per-task workers), aggregating
the results, and persisting the patch-level findings.

Skill:  ~/.claude/skills/evergreen-stack-ci/SKILL.md
Triage reference: ~/.claude/skills/evergreen-stack-ci/references/triage.md
Three-strikes ref: ~/.claude/skills/evergreen-stack-ci/references/three-strikes.md
State:  <state-file-path>
Stack-root: <stack-root>
Patch-id: <patch-id>
Branch:   <branch-name>
Trunk:    <trunk-branch, usually master>

Think harder before answering — verdict and plan quality directly determine
whether the downstream fix subagent gets it right on the first try.

Task — coordinate the investigation of why this patch failed:

1. Read SKILL.md, references/triage.md, references/three-strikes.md.

2. Read the state file. Pull, for this patch:
     - failed_tasks            (list of task names that failed)
     - build_failure           (bool, optional — set by poll subagent)
     - polling.compile_fail_fast (object, optional — set by poll subagent)

3. SHORT-CIRCUIT 1 — compile fail-fast:
   If polling.compile_fail_fast.suspect_patch_id == <patch-id>:
     a. Identify the compile task. For kind=COMPILE_BAZEL or
        kind=COMPILE_CLIENT_BAZEL, it's the matching entry in failed_tasks.
        For kind=BUILD_FAILURE, identify the failing test-runner task(s)
        with `has_test_results==false` AND `total_test_count==0` via
        mcp__evergreen__get_patch_failed_jobs_evergreen — this may be
        RUN_ALL_UNIT_JAVA_TESTS, INT_JAVA_*, INT_JAVA_THIRDPARTY_*, or any
        other test-runner task, NOT just RUN_ALL_UNIT_JAVA_TESTS. Pull THAT
        task's logs (not RUN_ALL_UNIT_JAVA_TESTS' by default).
     b. Call mcp__evergreen__get_task_logs_evergreen with
        `filter_errors=false` (NOT the default `true` — it drops bazel
        `[strict]`, `** Please add ...`, and ECJ `1. ERROR in ...` lines).
        Extract a one-line error signature using the two-pass scan from
        references/triage.md step 2: first a case-insensitive `ERROR:` grep
        with same-timestamp multi-line capture and dependency-download
        exclusion (transient artifact-fetch blocks are red herrings); only
        fall back to bazel/ECJ patterns if `ERROR:` pass yields nothing.
        Example signatures: "missing symbol Foo at Bar.java:42",
        "annotation processor X failed".
     c. Skip the master cross-check entirely — compile errors are branch-local.
     d. Record once:
          python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py record-failure \
            --stack-root <stack-root> --patch-id <patch-id> \
            --branch <polling.compile_fail_fast.suspect_branch> \
            --task <compile-task> --test "<error-signature>"
     e. Draft the implementation plan from the specific file/symbol named in
        the error. Jump to step 7 with verdict=real-bug and
        suspect-branch=<polling.compile_fail_fast.suspect_branch>.

4. SHORT-CIRCUIT 2 — build failure on any test-runner task:
   Else if patch.build_failure == true:
     a. Identify which failed task(s) triggered the flag. Call
        mcp__evergreen__get_patch_failed_jobs_evergreen on this patch and
        pick the failing tasks where `has_test_results==false` AND
        `total_test_count==0`. This may be RUN_ALL_UNIT_JAVA_TESTS (the
        original case), an INT_JAVA_* task (billing, payments,
        paymentview, …), INT_JAVA_THIRDPARTY_*, or any other test-runner
        task — do NOT hardcode RUN_ALL_UNIT_JAVA_TESTS.
     b. Pull logs for that failing task via
        mcp__evergreen__get_task_logs_evergreen with `filter_errors=false`.
        IMPORTANT — pass `filter_errors=false` explicitly. The default
        `filter_errors=true` drops bazel `[strict]`, `** Please add ...`,
        and ECJ `1. ERROR in ...` lines because they don't match the regex
        it uses for stderr-style errors. With the filter on you will see a
        misleading empty log. Do NOT trust the task's own pass/fail status
        either — the signal is `has_test_results=false`, not the status.
     c. Extract a one-line build-error signature via the two-pass scan
        from triage.md "Build failure on a test-runner task" step 2:
        FIRST a case-insensitive `ERROR:` grep with same-timestamp
        multi-line capture, discarding dependency-download blocks
        (transient artifact-fetch noise — not the real failure). ONLY
        IF that surfaces no surviving block, fall back to the bazel/ECJ
        pattern list (`[strict]` indirect-dep, `ERROR in ...`, `cannot
        find symbol`, `package ... does not exist`, NullAway,
        `BUILD FAILED`). For `[strict]` errors specifically, capture the
        buildozer command verbatim from the error message — the error
        literally tells you which dep to add.
     d. Apply the novel-failure walk from triage.md to attribute to the
        earliest branch whose diff plausibly caused the build error. That is
        your suspect-branch.
     e. Skip master cross-check — build errors are branch-local.
     f. Record once (use the ACTUAL failing task name, not a hardcoded one):
          python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py record-failure \
            --stack-root <stack-root> --patch-id <patch-id> \
            --branch <suspect-branch> \
            --task <failing-task-name> --test "<error-signature>"
     g. Draft the plan from the file/symbol named in the error. For
        `[strict]` errors use the buildozer command verbatim — do NOT guess
        deps. The plan's `Verify (test):` must pass `--strict_java_deps=error`
        (see SKILL.md Hard Rule 12) so the local canary matches Evergreen's
        strict-deps mode; without it the strict-deps error reappears in CI.
        Jump to step 7 with verdict=real-bug.

5. SHORT-CIRCUIT 3 — single failed task:
   Else if len(failed_tasks) == 1:
     Investigate inline (no fan-out).
     a. Call mcp__evergreen__get_task_test_results_evergreen on the one task.
        If you need a task-id, derive it via
        mcp__evergreen__get_patch_failed_jobs_evergreen on this patch first.
     b. For EACH failing test in this task, cross-check master via
        mcp__evergreen__list_user_recent_patches_evergreen:
          - master also red on (task, test):
              record-master-broken --branch <suspect-branch> \
                --task <task> --test <test> --evidence "<short note>"
            (Do NOT also call record-failure for this test.)
          - master green on (task, test):
              record-failure --patch-id <patch-id> \
                --branch <suspect-branch> --task <task> --test <test>
        Use the novel-failure walk to determine <suspect-branch> per test.
     c. Decide a verdict from the per-test results (see aggregation rules
        below — applied here with N=1 task).
     d. Draft the plan if verdict=real-bug. Jump to step 7.

6. FAN-OUT PATH — multiple failed tasks, not build/compile failure:
   Dispatch ONE template 3b subagent per task in failed_tasks, in a SINGLE
   message so they run in parallel. Each 3b prompt receives state-file-path,
   stack-root, patch-id, branch, trunk, AND its specific task-name.
   See template 3b below for the prompt body and structured return contract.

   Wait for all 3b workers to return. Each returns one structured block:

     ===TASK-INVESTIGATION-RESULT===
     task: <task-name>
     verdict: <real-bug|flake|master-broken|needs-retry|unknown>
     suspect_branch: <branch>
     cause: <one-line root cause for this task>
     failing_tests: <test1>; <test2>; ...
     master_status: <all-green|some-red|all-red|skipped>
     evidence:
     - file: <path>
       symbol: <method/class>
       line: <number-or-empty>
       signal: <failing assertion / compile error / exception type>
     ...
     ===END===

   The 3b workers have already called record-failure / record-master-broken
   for their tests; the state file is up to date by the time they return.
   Your job is to AGGREGATE — never re-pull logs, never re-call MCP for data
   the 3b workers already extracted.

   a. PATCH-LEVEL VERDICT — apply these rules in order to N task verdicts:

      1. Any task verdict=real-bug         → patch verdict=real-bug
      2. Else any task verdict=needs-retry → patch verdict=needs-retry
      3. Else any task verdict=unknown     → patch verdict=unknown
      4. Else all tasks master-broken      → patch verdict=master-broken
      5. Else all tasks flake              → patch verdict=flake
      6. Else mix of flake + master-broken → patch verdict=master-broken
         (master-broken is the more informative classification — flakes are
         noise; the user benefits more from seeing master-broken surfaced)

   b. SUSPECT-BRANCH — most 3b workers will name the same suspect (the earliest
      novel-failure branch per triage.md). If they disagree, pick the
      EARLIEST in stack order (branches in state.branches are root→tip).
      That earliest one is the upstream cause; downstream attributions are
      cascading effects.

   c. PLAN (verdict=real-bug only) — collect every 3b's evidence block. Look
      for shared root causes across tasks: if two tasks fail with the same
      null-pointer in the same file, one plan step suffices. Group steps by
      file (related changes together) and order by dependency (helper changes
      before callers). Be CONCRETE — name files, methods/classes, specific
      lines or symbols. Template 4 executes this without re-investigating, so
      anything you don't write down is lost.

      Per SKILL.md Hard Rule 12, the plan MUST end with explicit `Verify
      (build):` and `Verify (test):` steps (see the step 7 plan template
      below for the exact shape). For the canary `Verify (test):` target:
      if the 3b workers' evidence names test files (the typical case), pick
      the smallest target + filter pair that exercises the same code path
      the fix touches — the goal is to catch @BeforeEach explosions, broken
      stub setup, and obvious regressions in seconds, not exhaustive
      coverage. If the evidence is purely build/compile-only (short-circuits
      1 and 2) and no Java test target covers the touched files, use the
      `Verify (test): N/A — <one-line reason>` form documented in step 7.
      An unset or vague `Verify (test):` is a plan defect — template 4 will
      return FAILED — plan-mismatch and the coordinator will re-dispatch
      investigation.

   d. CAUSE — synthesize one paragraph naming the root cause(s) and the
      failing tasks affected. If multiple unrelated causes across tasks, list
      each briefly.

7. Persist findings (ALL PATHS converge here):
     python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py set-findings \
       --stack-root <stack-root> --patch-id <patch-id> \
       --notes "<cause + plan (real-bug) or one-paragraph cause (other)>" \
       --cause "<one-line root cause>" \
       --suspect-branch "<aggregated suspect-branch>" \
       --verdict <aggregated-verdict>

   For verdict=real-bug, --notes uses this exact structure:

       Cause: <one-paragraph root cause spanning all affected tasks>

       Plan:
       1. <file-path:line-or-symbol> — <what to change and why>
       2. <file-path:line-or-symbol> — <what to change and why>
       ...
       N-1. Verify (build): bazel build //path/to/touched:target
       N. Verify (test): bazel test //path/to/affected:target --test_filter=<canary-test-method> --strict_java_deps=error --test_output=errors
            (Pick the smallest representative target that exercises the
            changed code path. If the fix touches multiple test classes,
            pick one canary class — the goal is to catch @BeforeEach
            explosions and stubbing errors in seconds, not exhaustive
            coverage. The `--strict_java_deps=error` flag is mandatory:
            local bazel defaults to `warn` while Evergreen uses `error`, so
            without it a missing direct dep can pass locally and still fail
            CI. See references/triage.md "Why local `bazel test` may pass
            when Evergreen fails" and SKILL.md Hard Rule 12. Do NOT skip
            this step.)

            If NO bazel test target covers the touched files (e.g. the diff
            is build-config only, generated code, or pure docs), write
            instead:
              Verify (test): N/A — <one-line reason, e.g. "diff is BUILD-file
              only; no Java test target covers the change">
            This is the only acceptable form of omission; an unset or vague
            Verify (test) is a plan defect and will cause template 4 to
            return FAILED — plan-mismatch.

   Avoid vague directives like "fix the bug" or "update the test" — template 4
   runs on sonnet and needs an executable plan, not a pep talk. For other
   verdicts, --notes is just the one-paragraph cause; skip the Plan: section.

After `set-findings`, you MUST confirm the write succeeded by checking the
return value of the `stack_state.py set-findings` call. If the call fails
(sandbox restriction, lock contention, etc.), include
"WARN: set-findings failed — <reason>" in your return line so the coordinator
knows to re-apply the findings directly via its own Bash tool. Without this
warning, a silently-failed write leaves the patch's findings block null and
the coordinator's `summary` rollup will misclassify the patch as unresolved.

Return ONE LINE: "<branch>: <verdict> — <one-line cause>" (append
" | WARN: set-findings failed — <reason>" if step 7 did not persist). No log
dumps. No JSON dumps from 3b returns. No raw MCP output. The state file
already carries the structured detail (when set-findings succeeded).
```

## 3b. investigate-task (one per failed task within a patch, parallel)

Dispatched only by template 3a in the fan-out path. Each 3b worker owns exactly one task within one patch — narrow context, log-heavy reading, classification per failing test. Read-only with respect to source code; writes per-test outcomes via `stack_state.py record-failure` / `record-master-broken`.

**Dispatch with `model: "opus"`.** Each 3b worker still sees narrow scope (one task's worth of MCP data) — the two-level fan-out is what keeps wall-clock manageable, with one 3b per failed task dispatched in parallel from the 3a coordinator. The model tier is opus because the flake-vs-real-bug classification this worker produces is the load-bearing input to the 3a aggregate verdict and template 4's commit. A wrong verdict here costs a whole 30-min CI cycle re-patching a phantom fix, which dwarfs per-task token cost. The "think harder" cue inside the prompt body elevates extended thinking on opus for the subtler calls (intermittent timeouts, ordering-dependent failures, master-broken cross-checks).

```
You are a worker subagent for the evergreen-stack-ci skill (template 3b:
investigate-task). You investigate ONE failed task within ONE failed patch
and return a structured evidence block to your dispatching 3a coordinator.

Skill:  ~/.claude/skills/evergreen-stack-ci/SKILL.md
Triage reference: ~/.claude/skills/evergreen-stack-ci/references/triage.md
Three-strikes ref: ~/.claude/skills/evergreen-stack-ci/references/three-strikes.md
State:  <state-file-path>
Stack-root: <stack-root>
Patch-id: <patch-id>
Branch:   <branch-name>
Trunk:    <trunk-branch, usually master>
Task-name: <task-name to investigate>

Think harder before classifying — the flake-vs-real-bug call determines
whether the downstream fix subagent wastes time on a flake or misses a real
bug. Your evidence pointers also feed the patch-level plan, so be precise
about file paths, symbols, and failing signals.

Task — investigate this ONE failed task on this ONE patch:

1. Read SKILL.md, references/triage.md, references/three-strikes.md (skim).

2. Use mcp__evergreen__get_task_test_results_evergreen on <task-name> for this
   patch. If you need a task-id, derive it via
   mcp__evergreen__get_patch_failed_jobs_evergreen on this patch — the
   coordinator deliberately did NOT pre-resolve task-ids; you do it inline
   so each 3b worker stays self-contained.

3. BUILD-FAILURE BAILOUT: if the task has `has_test_results==false` AND
   `total_test_count==0` (visible in the get_patch_failed_jobs response or
   the test-results response itself), OR the test-results output contains
   "No test results found — there may have been a build failure for this
   team", this is a BUILD failure, not a test failure — the task never
   reached the test phase. STOP and return immediately with:

     ===TASK-INVESTIGATION-RESULT===
     task: <task-name>
     verdict: unknown
     suspect_branch: <branch>
     cause: build-failure pattern detected on <task-name>
            (has_test_results=false); 3a coordinator should handle via
            build-failure short-circuit
     failing_tests:
     master_status: skipped
     evidence:
     ===END===

   Do NOT call record-failure or record-master-broken. Do NOT run the tests
   locally to "see what happens" — local bazel uses
   `--strict_java_deps=warn` while Evergreen uses `error`, so a strict-deps
   build error CAN pass locally and mislead you into classifying a real bug
   as a flake. The 3a coordinator's build-failure short-circuit should
   normally pre-empt this case (the poll subagent flips patch.build_failure
   at poll time across ALL test-runner tasks, not just RUN_ALL_UNIT_JAVA_TESTS),
   but bail safely if it somehow slipped through.

   IMPORTANT — bazel error patterns are NOT captured by `filter_errors=true`.
   When investigating any task whose `failure_details.failing_command`
   mentions `bazel` (e.g. `'run bazel test'`, `'run bazel build'`), pull
   logs with `filter_errors=false`. The default filter drops bazel
   `[strict]`, `** Please add ...`, and ECJ `1. ERROR in ...` lines because
   they don't match the regex it uses for stderr-style errors. Pulling with
   the filter on will give you an empty / misleading log and you'll
   misclassify the failure.

   When you do read the raw log, apply the two-pass scan from
   references/triage.md step 2: FIRST a case-insensitive `ERROR:` grep
   with same-timestamp multi-line capture and dependency-download
   exclusion (artifact-fetch blocks are transient red herrings); the
   bazel `[strict]` / `** Please add ...` / ECJ patterns are the
   fallback ONLY if the `ERROR:` pass yields nothing after filtering.

4. If test results are present, identify each failing test in this task. Use
   mcp__evergreen__get_task_logs_evergreen ONLY if you need additional context
   (stack trace, assertion message, compile error excerpt) that the test
   results don't surface. Avoid pulling raw logs you don't read — that just
   inflates your context for no benefit.

5. For EACH failing test in this task, cross-check master via
   mcp__evergreen__list_user_recent_patches_evergreen filtered to recent
   master patches (or the project's master health board):

     - If master is ALSO red on (task, test):
         python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py record-master-broken \
           --stack-root <stack-root> --branch <suspect-branch> \
           --task <task-name> --test <test-name> --evidence "<short note>"
       Per-test verdict: master-broken. Do NOT also call record-failure.

     - If master is GREEN on (task, test):
         python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py record-failure \
           --stack-root <stack-root> --patch-id <patch-id> \
           --branch <suspect-branch> --task <task-name> --test <test-name>
       Per-test verdict: pick one of {real-bug, flake, needs-retry}:
         - real-bug: the diff between <suspect-branch> and its parent
           plausibly causes the failure (changed files match files in the
           stack trace or assertion).
         - flake: timing-sensitive, network call, or no plausible link to
           the diff.
         - needs-retry: environmental / infrastructure issue (Evergreen
           agent crashed, runner died, transient AWS error).

6. SUSPECT-BRANCH (per triage.md novel-failure walk): the earliest branch in
   stack order where this (task, test) first appears. This is what you pass
   to record-failure / record-master-broken — NOT necessarily <branch>, the
   branch whose patch you're investigating. The round-id dedup in
   record-failure handles cascading attributions across parallel workers.

7. TASK-LEVEL VERDICT — aggregate the per-test verdicts using these rules
   (same priority order the 3a coordinator uses for patches):

     1. Any test real-bug              → task real-bug
     2. Else any test needs-retry      → task needs-retry
     3. Else all tests master-broken   → task master-broken
     4. Else all tests flake           → task flake
     5. Else mix flake + master-broken → task master-broken

8. Return ONE structured block in this exact format (no preamble, no
   postamble, no log dumps):

   ===TASK-INVESTIGATION-RESULT===
   task: <task-name>
   verdict: <real-bug|flake|master-broken|needs-retry|unknown>
   suspect_branch: <branch-name>
   cause: <one-line root cause for this task>
   failing_tests: <test1>; <test2>; <test3>
   master_status: <all-green|some-red|all-red|skipped>
   evidence:
   - file: <path>
     symbol: <method/class>
     line: <number-or-empty>
     signal: <failing assertion / compile error / exception type — SHORT>
   - file: <path>
     ...
   ===END===

   evidence lists 1–5 compact pointers per task — enough for the 3a
   coordinator to draft a plan without re-pulling logs. Empty list is fine
   for non-real-bug verdicts. NO log dumps, NO stack traces verbatim. If
   multiple failing tests share the same evidence, list it once.

If anything goes wrong (MCP call fails, test results unparseable, etc.),
return the block with verdict=unknown and a cause that names the problem.
Do NOT retry inside this subagent — the 3a coordinator decides what to do.
```

The "suspect-branch" passed to `record-failure` / `record-master-broken` is the branch that introduced the failure (per triage.md), NOT necessarily the branch whose patch is being investigated. Cascading failures investigated in parallel — across patches via 3a fan-out, or across tasks within one patch via 3b workers — all attribute to the same suspect branch. The round-id dedup in `record-failure` keeps the per-test counter from over-incrementing across either layer of parallelism.

---

Back to [SKILL.md](../SKILL.md).
