# Test scope: profiles, exclusions, auto-detection

Each Evergreen patch can run a subset of variants/tasks. For stack CI, narrower is usually better — running `full` across 5 branches is slow and expensive, and frontend/E2E failures are often unrelated to the backend changes you're trying to validate.

The skill exposes named **profiles** (resolve to variant/task lists), **exclusions** (drop categories from any profile), and **thirdparty-team selection** (the most expensive class of tests, gated separately).

## Profiles

Resolved at `start` and recorded in the state file's `scope`. Reused for re-patching during `fix` so each iteration is comparable.

| Profile | What it runs | When to use |
|---|---|---|
| `backend` | `unit_java`, `int`, `code_health` (Java side: COMPILE_BAZEL, BLOCK_COMMIT_TASK, lint). **Thirdparty NOT included by default — opt in via `--thirdparty-teams`.** | Default for Java-only stacks. Fast, covers compile + unit + integration |
| `unit` | `unit_java`, `unit_python` | Smoke check; fastest. Use when you just need "does it compile and pass unit tests" |
| `compile` | `code_health -t COMPILE_BAZEL -t COMPILE_CLIENT_BAZEL` | Even faster — verify the stack still compiles after a rebase |
| `frontend` | `js`, `code_health -t COMPILE_CLIENT_BAZEL` | Frontend-only stacks |
| `full` | Evergreen alias `full` (JS + E2E local + server tests) | **Avoid for `mms`.** See "Never run the full Evergreen build for mms" below. Only use after the user has explicitly re-confirmed they really want full after being offered a narrower equivalent. |
| `custom` | Use `--variants` / `--tasks` directly | Anything off the menu |

## Never run the full Evergreen build for `mms`

`--profile=full` (and the equivalent bare `evergreen patch -a full` or `-a all`) fans out into the entire `mms` variant/task graph: tens of variants, hundreds of tasks, plus the entire E2E suite. Even on a single branch it routinely takes 30+ minutes; multiplied across a stack and across re-patch iterations, it burns hours of compute and human wait time per session for almost no signal beyond what `backend` + `--thirdparty-teams=auto` already gives.

**Rules:**

- **Default to the narrowest profile that covers the diff.** Java-only stacks → `backend`. Client-only stacks → `frontend`. Smoke checks → `compile` or `unit`.
- **Never auto-propose `full`.** The auto-detection heuristic below explicitly suggests narrower equivalents for mixed stacks (`full --exclude=cypress --exclude=e2e` is not a default — it's a last-resort suggestion only when the diff genuinely spans backend + frontend + E2E).
- **Push back on explicit `full` requests.** If the user asks for `--profile=full`, "everything", "the whole build", or "run all of CI", respond once with: (1) the narrower equivalent that covers the diff, (2) a reminder that full is ~30+ min/branch and multiplies across the stack, and (3) a one-question confirmation. Only proceed with `full` after the user re-confirms in that turn.
- **Same rule applies to `-a full` / `-a all` aliases.** Don't silently expand them into a patch — confirm first.
- **Never run `full` to "save time" on triage.** It does the opposite. If the user is uncertain about scope mid-session, suggest narrowing further (e.g. drop `int` or `code_health`), not widening to `full`.

This rule does *not* apply to `--alias=<custom-alias>` for narrow named scopes — only to `full` (and `all`, which is the same thing on mms).

## Exclusions

Layered on top of any profile/alias. Repeatable. Drops matching variants/tasks from the patch.

| Exclude key | Drops |
|---|---|
| `frontend` | `js`, `client` alias, `COMPILE_CLIENT_BAZEL` |
| `e2e` | All `e2e_*` variants and the `e2e` alias |
| `cypress` | `E2E_CYPRESS_*` tasks (already known-flaky per `evergreen-cicd` skill) |
| `python` | `unit_python`, `python` alias |
| `thirdparty` | The entire `thirdparty` variant (overrides any `--thirdparty-teams`) |
| `int` | `int` variant — useful when iterating on something the unit suite covers |
| `code_health` | `code_health` variant (lint/checkstyle) — when iterating on code that hasn't been formatted |

Examples:
- `--profile=full --exclude=cypress --exclude=e2e` → run everything except E2E
- `--profile=backend --thirdparty-teams=payments,billing` → backend with only the relevant thirdparty teams
- `--variants=unit_java --tasks=GENERATE_UNIT_TESTS_BAZEL` → fully manual

## Thirdparty teams (separate from `--exclude`)

Thirdparty integration tests are the **most expensive** class of tests in the repo: each team's generator expands into many `INT_JAVA_*` tasks. Running all 8 teams against every branch in a 5-branch stack can mean dozens of hours of compute.

The `thirdparty` Evergreen variant contains 8 team-specific generator tasks. Each generator emits N integration tasks based on a per-team YAML manifest. Select teams via `--thirdparty-teams=<spec>`:

| Spec | Behavior |
|---|---|
| (omitted) | **Skip thirdparty entirely.** Default for `backend` profile — don't run thirdparty unless asked. |
| `--thirdparty-teams=auto` | Detect which teams' code is touched in the stack diff (see "Auto-detection" below) and run only those generators. If nothing matches, **skip the variant entirely and surface a persistent reminder** — never fall back to running all. |
| `--thirdparty-teams=all` | Run every team's generator. **Only when the user explicitly asks for it.** Very expensive — never use as an auto fallback. |
| `--thirdparty-teams=payments,billing` | Explicit comma-separated team list. |

Known teams and their generator tasks + path globs:

| Team | Generator task | Path glob (matches → team is touched) |
|---|---|---|
| `billing` | `GENERATE_THIRDPARTY_BILLING_INT_TESTS_BAZEL` | `server/src/**/svc/mms/svc/billing/**`, `server/src/**/svc/mms/res/**`, `server/src/**/svc/nds/**` |
| `payments` | `GENERATE_THIRDPARTY_PAYMENTS_INT_TESTS_BAZEL` | `server/src/**/cloud/payments/**`, `server/src/**/cloud/partners/**`, `server/src/**/cloud/services/payments/**`, `server/src/**/cloud/common/jira/**` |
| `metering` | `GENERATE_THIRDPARTY_METERING_INT_TESTS_BAZEL` | `server/src/**/cloud/billingimport/**` |
| `intel` | `GENERATE_THIRDPARTY_INTEL_INT_TESTS_BAZEL` | `server/src/**/svc/mms/svc/ping/**` |
| `iam_authn` | `GENERATE_THIRDPARTY_IAM_AUTHN_INT_TESTS_BAZEL` | `server/src/**/cloud/services/authn/**`, `server/src/**/cloud/common/geolocation/**` |
| `iam_authz` | `GENERATE_THIRDPARTY_IAM_AUTHZ_INT_TESTS_BAZEL` | `server/src/**/cloud/federation/**`, `server/src/**/module/federation/**` |
| `iam_workid` | `GENERATE_THIRDPARTY_IAM_WORKID_INT_TESTS_BAZEL` | `server/src/**/cloud/account/runtime/**` |
| `iam_identity_security` | `GENERATE_THIRDPARTY_IAM_IDENTITY_SECURITY_INT_TESTS_BAZEL` | `server/src/**/cloud/user/**`, `server/src/**/cloud/azurenative/**`, `server/src/**/cloud/common/okta/**`, `server/src/**/module/account/**`, `server/src/**/svc/mms/res/admin/**` |

If a path matches multiple globs, include all matching teams. The mapping is derived from `server/src/test/thirdparty_<team>_test_tasks.yml` — when in doubt, grep those files for the exact bazel package being touched.

**Granularity not possible? Skip, never auto-fallback to `all`.**

Thirdparty tests are expensive enough that running all of them silently would burn hours of CI for no clear reason. Instead:

1. **Skip the entire `thirdparty` variant** for this run.
2. **Mark `scope.thirdparty_status = "skipped-no-mapping"`** in the state file (the field is sticky for this stack-CI session, so the warning persists across re-patches).
3. **Surface a non-blocking reminder banner** on every coordinator phase output (start, status, fix, update-prs) AND in the closing summary, telling the user:
   - Why thirdparty was skipped (couldn't map the diff to a team).
   - Exactly which `--thirdparty-teams=...` invocation they should use to re-run with thirdparty included (suggest `auto` won't help — they'll need an explicit team list or `all`).
   - That this is informational only — the cycle is not blocked.

Stick to skipping. Never run all 8 generators as a fallback unless the user explicitly asks for `--thirdparty-teams=all`.

## Resolved flags per profile

| Source | Resolved `evergreen patch` flags |
|---|---|
| `--profile=backend` | `-v unit_java -v int -v code_health -t GENERATE_UNIT_TESTS_BAZEL -t GENERATE_INT_TESTS_BAZEL -t all` |
| `--profile=unit` | `-v unit_java -v unit_python -t GENERATE_UNIT_TESTS_BAZEL -t all` |
| `--profile=compile` | `-v code_health -t COMPILE_BAZEL -t COMPILE_CLIENT_BAZEL` |
| `--profile=frontend` | `-v js -v code_health -t all -t COMPILE_CLIENT_BAZEL` |
| `--profile=full` | `-a full` |
| `--alias=<x>` | `-a <x>` |
| `--variants` + `--tasks` | as supplied |
| `--thirdparty-teams=<list>` | adds `-v thirdparty` and one `-t GENERATE_THIRDPARTY_<TEAM>_INT_TESTS_BAZEL` per resolved team |

Apply exclusions by *removing* matching `-v`/`-t` entries from the resolved list. `--exclude=thirdparty` drops the whole `thirdparty` variant regardless of what `--thirdparty-teams` selected. For alias-based runs (`-a <name>`), exclusions can't be subtracted from the alias — instead, expand the alias into explicit variants/tasks (consult the project's `evergreen-cicd` skill or `evergreen list --aliases -p mms`) and remove from there.

## Auto-detection of default profile

Before prompting the user, peek at what's changed in the stack:

```bash
git diff --name-only $(git merge-base HEAD master) HEAD
```

Heuristic for the **profile**:
- All paths under `server/` (and similar Java dirs) → suggest `backend`
- All paths under `client/` → suggest `frontend`
- Mixed → suggest running **two narrower patches in parallel** (one `--profile=backend`, one `--profile=frontend`) instead of `full`. Per the "Never run the full Evergreen build for `mms`" rule above, `full` is not an auto-suggestion — it's only used when the user explicitly asks for it and re-confirms.
- Only docs / .planning / yaml → ask user; CI may not be needed at all

Heuristic for **thirdparty teams** (only relevant when `--thirdparty-teams=auto` was requested, or when the user is choosing a profile that *could* include thirdparty):
- Match the changed paths against each team's path globs in the table above.
- Resolved teams = union of every team whose glob matched any changed path.
- If 1+ teams matched → record `scope.thirdparty_teams = [<list>]` and `scope.thirdparty_status = "auto-resolved"`. Run only those generators.
- If 0 teams matched **or the diff is purely shared-infra paths with no team mapping** → record `scope.thirdparty_status = "skipped-no-mapping"` and the reason text. Do NOT run any thirdparty generator. Surface the persistent reminder banner described above on every coordinator output.

Never propose `--thirdparty-teams=all` automatically. Suggesting `all` is only appropriate when the user explicitly says "run all thirdparty" or "I don't care about cost".

Always confirm both the profile **and** the thirdparty-team selection with the user once before `start`. Don't quietly default — running the wrong scope wastes minutes per branch × N branches, and running thirdparty when not needed wastes hours.

## Why scope is sticky

Once chosen at `start`, the scope (including the resolved thirdparty team list) is recorded in `state.scope` and reused for every re-patch. Don't change scope mid-stack unless the user explicitly says so — comparing patches across iterations only works if the inputs match.
