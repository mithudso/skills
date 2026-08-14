# Three-strikes quarantine, master-broken filtering, and the polling decision

The skill keeps cycling polling → fix → re-patch indefinitely until ONE of these is true:

1. **All-clean**: every patch is `succeeded`. Hand off to Phase 4 (or stop).
2. **Excluded-only**: every remaining failure is either *quarantined* (3 consecutive failures) or *master-broken* (already failing on master). The skill won't try any more fixes — surface the list to the user.

There is **no overall time cap**. The 1-hour cap that exists is per-cycle (see "Polling cycle" below).

## Two failure-exclusion mechanisms

Both feed into the actionability calculation in `stack_state.py summary`. Excluded failures stay visible in the summary but do NOT trigger fix attempts.

### Quarantine (three-strikes)

After a test fails on **3 consecutive rounds** with the same `(suspect-branch, task, test)` key, the skill stops trying to fix it. The user can keep working on it manually; the skill moves on to other fixable failures.

**A "round"** is a single patch on the suspect branch. Round-id dedup makes cascading failures investigated in parallel across child patches count once, not N times.

Maintained by `stack_state.py record-failure` which:
- Increments `consecutive_failures` for the `(branch, task, test)` key (deduped by round).
- Sets `quarantined: true` once the counter hits 3.
- Resets via `record-success` when a previously failing test passes on a later round.

### Master-broken filtering

Tests that are **already failing on master** are not the stack's problem. The investigation subagents detect this by querying master patches via `mcp__evergreen__list_user_recent_patches_evergreen` (or the project's master health board) and call `record-master-broken` instead of `record-failure` — see "What the investigation subagents do" below for which template (3a inline vs. 3b fan-out) runs the check.

A master-broken flag is sticky — it never auto-clears. If master gets fixed and the test starts passing, the user should manually clear it (or just delete the state file and start a fresh stack-CI session).

## What the investigation subagents do (per failing test)

In the two-level fan-out introduced for Phase 3.2 (templates 3a / 3b in [`subagent-investigate.md`](subagent-investigate.md)), the per-test record loop runs at the lowest level that has the test-level evidence — never at the top coordinator.

- **3a coordinator (inline short-circuit paths only)** — compile-fail-fast, build-failure, single-task. There's only one task to investigate, so the 3a coordinator runs the loop itself.
- **3b task investigator (fan-out path)** — multi-task patches. Each 3b worker runs the loop for its one task's failing tests in parallel with sibling 3b workers.

```
for each failing test on the patch (or task, in 3b):
    if master is also failing on (task, test):
        record-master-broken --branch <suspect> --task <T> --test <test> --evidence <note>
    else:
        record-failure --patch-id <P> --branch <suspect> --task <T> --test <test>
```

Both calls also append to the patch's `failed_tests` array in the state file, so summary's per-patch actionability check works. The 3a coordinator's `set-findings` call comes AFTER all per-test writes have completed — it only stamps the patch-level `findings` block, never the per-test data, so there's no race between concurrent 3b writes and the later `set-findings`.

## Polling cycle and decision tree

A polling cycle = poll every **5 minutes** for up to **1 hour** (12 iterations). The coordinator drives the loop via `ScheduleWakeup`, but does **not** invoke the skill on each wakeup. Each iteration is just: dispatch the poll-status subagent, read its one-line return, decide. The poll subagent itself bumps the iteration counter (in the state file) and reads `stack_state.py summary` to compute the decision; the coordinator never calls `summary` directly during polling.

| `decision:` value (from subagent return) | Meaning | Coordinator action |
|---|---|---|
| `actionable_failure` | At least one failed patch with non-excluded failing tests | STOP polling; go to Phase 3 (`fix`) |
| `in_progress` | No actionable failures yet, some patches still running | Schedule next wakeup (5 min) — same canonical prompt |
| `excluded_only` | All failed patches have only quarantined / master-broken failures, nothing running | STOP polling; surface to user |
| `all_clean` | All patches succeeded | STOP polling; offer Phase 4 |
| `needs_attention` | Edge case (e.g. aborted, no patches) | Surface to user |

```
# Inside each wakeup turn (driven by ScheduleWakeup with the canonical prompt in SKILL.md):
ret = dispatch poll-status subagent     # state writes + counter bump + decision compute happen inside
parse "iteration <i>/<max> | decision: <D>" from ret

case D:
  actionable_failure -> stop polling, go to fix
  all_clean          -> stop polling, offer Phase 4
  excluded_only      -> stop polling, surface quarantined + master-broken
  needs_attention    -> surface to user
  in_progress        -> if i >= max: ask user (extend? stop?)
                        else: ScheduleWakeup(5 min, canonical prompt); end turn
```

After Phase 3 fixes and re-patches, the coordinator runs `stack_state.py reset-poll-cycle --stack-root <r>` before re-arming the wakeup loop, so each new cycle gets a **fresh** 12-iteration budget. The hour cap restarts every cycle. Indefinite cycling is fine — the only termination conditions are `all_clean` and `excluded_only`.

## Why fail-fast on actionable failures (not on every failure)

The coordinator stops polling early when there is something *fixable*. If the only failures are quarantined / master-broken, those patches are already excluded from the fix loop — there is no reason to interrupt other still-running patches. Continuing to poll lets a different test that ISN'T excluded surface, which the skill can then act on.

## When all remaining failures are excluded

When `poll-decision: excluded_only`, the coordinator surfaces this to the user with both lists:

```
python3 ~/.claude/skills/evergreen-stack-ci/scripts/stack_state.py quarantine \
  --stack-root <root> --repo-root <repo>
```

And reads the state's `test_failures` for entries with `master_broken: true`. Together these tell the user exactly which tests the skill gave up on and why. The user can then:
- Manually fix and re-run patches (use `record-success` to clear quarantine).
- Confirm master-broken is real and ignore it.
- Reduce scope (e.g. `--exclude=int`) and try again.

## Edge cases

- **A test passes one round, fails the next**: counter resets via `record-success` on the passing round, so it has to start fresh from 1.
- **Test renamed mid-stack**: the new key is treated as a new test (no shared history with the old one). Acceptable — renames are rare and surface as a new failure.
- **Master is intermittently red**: master-broken is recorded once; the skill won't re-check master automatically. If the user thinks master has recovered, they should clear the entry manually.
- **All branches are excluded but one is still running**: poll continues until that one finishes, then the decision flips to `excluded_only` (or `actionable_failure` / `all_clean` depending on the result).
