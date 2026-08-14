# Subagent prompts — fix-and-commit (4) and submit-pr (5)

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

The model is set on the Agent tool call itself, not inside the subagent prompt body — the coordinator passes `model: "opus"` for investigation (3a, 3b) or `model: "sonnet"` for code editing (4) as an argument when dispatching. The "think harder" / "think hard" cues for templates 3b and 4 live inside the prompt bodies (see `subagent-investigate.md` for 3b; see below for 4).

If a future subagent type is added, default to `inherit` unless it requires open-ended classification or planning — those tip into `opus`.

## 4. fix-and-commit (one per branch — sequential or parallel)

Sequential for linear stacks (fix the earliest suspect, then `gt restack` propagates). Parallel via worktrees only for sibling branches with independent failures — see `references/parallel-fixes.md` for the worktree variant + cleanup invariants.

**Dispatch with `model: "sonnet"`.** This is the implementation tier — applying the plan that template 3a (opus) already produced. The "think hard" trigger word is included in the prompt body to elevate reasoning effort within the sonnet tier; do not also pass a higher model. Investigation already happened; if the plan is wrong, that's a bug in template 3a (or in the 3b workers it aggregated), not a reason to upgrade template 4.

```
You are a worker subagent for the evergreen-stack-ci skill.

Skill:  ~/.claude/skills/evergreen-stack-ci/SKILL.md
State:  <state-file-path>
Stack-root: <stack-root>
Branch:     <branch-name>
Patch-id:   <patch-id>

Think hard about the fix before applying it. The investigate-coordinator subagent
(template 3a, opus) has already produced an implementation plan inside the patch's
findings.notes — your job is to execute that plan precisely, not to re-investigate. If the plan is
ambiguous or appears wrong on inspection of the actual code, STOP and return
"<branch>: FAILED — plan-mismatch: <one-line reason>" so the coordinator can
re-dispatch investigation. Do NOT improvise a different fix.

Task — execute the implementation plan from the investigation findings:
1. Read SKILL.md
2. Read the state file. Pull the patch's `findings.notes` — it contains a "Plan:"
   section authored by the investigate-coordinator subagent (template 3a). Also pull the
   patch's `failed_tests` list — you'll pass those keys as --target-tests in step 7.
3. `git checkout <branch>` (assume working tree is clean — coordinator verified)
4. Walk the Plan steps in order. For each step, edit the named file at the named
   symbol/line, applying the change described. Do NOT add unrelated changes,
   refactors, or "while I'm here" cleanups.
5. Run the Plan's `Verify (build):` step. If it fails, surface as
   "<branch>: FAILED — verify-build: <reason>" and do NOT commit. The mms
   Evergreen cycle is ~30+ min per patch; local verification catches
   regressions in seconds (see SKILL.md Hard Rule 12).
5b. Run the Plan's `Verify (test):` step. If it fails, surface as
   "<branch>: FAILED — verify-test: <reason>" and do NOT commit — the
   coordinator will re-dispatch investigation with the local failure as
   evidence. If the plan declares `Verify (test): N/A — <reason>`, skip this
   step and proceed to step 6 (this is the documented graceful-degradation
   path for diffs with no covering bazel test target — e.g. BUILD-file-only,
   generated code, or docs). If the plan omits `Verify (test):` entirely or
   specifies it without a concrete target / N/A justification, the plan is
   defective — return
   "<branch>: FAILED — plan-mismatch: missing Verify (test) step" so the
   coordinator re-dispatches investigation.
6. `git add -A && git commit -m "<jira>: <terse fix message>"` — capture the new
   SHA from `git rev-parse HEAD`.
7. Record the fix so it shows on the dashboard:
     python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py record-fix \
       --stack-root <stack-root> --branch <branch> \
       --commit-sha <new-sha> --summary "<commit subject line>" \
       --target-tests "<branch>::<task1>::<test1>,<branch>::<task2>::<test2>"
   The --target-tests CSV is the set of failing-test keys this fix is meant to address
   (one per failed_tests entry on the patch). Omit --target-tests if none apply.

Return ONE LINE: "<branch>: fix committed <sha-short> recorded" or
"<branch>: FAILED — <reason>". Do NOT include diff output.
```

## 5. submit-pr (one per branch with existing PR)

Phase 4 only. Sequential — `gt submit` walks the working tree.

```
You are a worker subagent for the evergreen-stack-ci skill.

Skill:  ~/.claude/skills/evergreen-stack-ci/SKILL.md
Branch: <branch-name>

Task — push update to the existing PR for this branch (NEVER create a PR):
1. `gt checkout <branch>`
2. `gt submit --no-stack --update-only`
   (--update-only ensures we don't accidentally create a PR if one doesn't exist;
    --no-stack limits the submit to just this branch)
3. If gt reports no PR exists for this branch, return failure — do NOT create one

Return ONE LINE: "<branch>: PR updated" or "<branch>: SKIPPED — <reason>".
```

## Why this split

A single "do the whole phase" subagent would hide failures and make debugging hard. Per-task subagents fail individually, can be retried individually, and write their results to the state file — so a partial failure leaves the state file consistent and the coordinator can decide what to redo. The state file is the only memory; it is the entire reason long stack-CI sessions stay tractable across context limits.

---

Back to [SKILL.md](../SKILL.md).
