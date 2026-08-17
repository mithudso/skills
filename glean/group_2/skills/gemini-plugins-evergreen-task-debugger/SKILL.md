---
name: evergreen-task-debugger
description: >-
  Use when debugging a failed Evergreen CI task interactively using
  the task debugger CLI on a spawn host or locally.
argument-hint: ssh user@host OR task-id
source: 10gen/evergreen
license: Internal
mongodb:
  team: devprod-bv
  owner: malik.hadjri@mongodb.com
  internal: true
---

# Evergreen Task Debugger

Debug failed CI tasks interactively by stepping through commands, inspecting logs, fixing issues, and re-running steps.

**Target:** $ARGUMENTS

## Spawn Host Environment

- **Binary path**: The `evergreen` binary is not on `$PATH` on spawn hosts. Use the full path `~/evergreen` (i.e. `/home/<user>/evergreen`).
- **Project source directory**: The checked-out project code lives at `/data/mci/debug_project_config/evergreen/`, not `~/debug_project_config/`.

## Startup

Check daemon state, then load the task:

```bash
~/evergreen debug daemon status
```

- **If the daemon is not running**: restart it with `~/evergreen debug daemon start`, then check `~/.evergreen-local/daemon.log` if it fails to start.

Determine the loading mode:

- **On a spawn host** (`/data/mci/debug_project_config/` exists):
  ```bash
  ~/evergreen debug load /data/mci/debug_project_config/evergreen.yml
  ~/evergreen debug select <task_name>
  ```
  If `list-steps` shows steps already completed, the host was created with a starting step — prior steps ran during setup. Continue from the current position.
- **Locally** (task ID argument provided):
  ```bash
  evergreen debug daemon start
  evergreen debug load --task-id <task_id>
  ```
  `select` is not permitted in local mode — the task is auto-selected from the ID.

## Investigation

1. List steps and identify failures:
   ```bash
   evergreen debug list-steps        # → = current, checkmark = pass, X = fail
   ```
2. Pull logs for failed steps:
   ```bash
   evergreen debug logs --step <N>   # N can be 3, 2.1, pre:1, post:2
   ```
3. Gather environment context: `pwd`, `ls`, `env | grep <VAR>`, `which <tool>`, `uname -a`.
4. On spawn hosts, read `/data/mci/debug_project_config/evergreen.yml` to understand step intent.
5. Use `grep` and `find` to locate scripts/code referenced in the error.

## Debugger-Specific Failure Patterns

- **Step shows "skipped"** — it's a no-op command (e.g., `s3.put`, `attach.results`). Expected behavior, not an error.
- **Permission denied on AWS assume-role** — debug sessions use a different external ID format (`debug-[project_id]-[requester]`). The IAM trust policy may not allow it. Report to project admin.
- **GitHub token insufficient permissions** — debug tokens use the "Debug" requester type, which may be more restrictive than regular task execution. Project admin must configure debug permissions.
- **Variable is empty/missing** — admin-only project variables are excluded from debug sessions for non-admin users. Cannot be worked around.
- **Daemon not responding** — check `~/.evergreen-local/daemon.log`. Restart with `evergreen debug daemon start`.
- **Steps already completed on first load** — the spawn host was created with a "starting step". All prior steps ran during host setup. Continue from the current position shown in `list-steps`.

## Fix-and-Verify Loop

1. Edit the relevant file (config, script, source).
2. **Ask user before re-running** — `next`, `run-all`, and `run-until` mutate execution state.
3. Re-run:
   ```bash
   evergreen debug jump <N>
   evergreen debug next
   ```
4. If it passes, offer to continue remaining steps. If it fails, iterate with new analysis.
5. Hot-reload config changes (preserves position, expansions, history):
   ```bash
   ~/evergreen debug load /data/mci/debug_project_config/evergreen.yml
   ```

## Command Reference

| Command | Purpose |
|---|---|
| `daemon start [--port N]` | Start debugger (default 9090) |
| `daemon stop` | Stop debugger |
| `daemon status` | Check daemon/task state |
| `load <config.yml>` | Load project config (spawn host) |
| `load --task-id <id> [config.yml]` | Load task by ID (local) |
| `select <task> [-v V]` | Select task and variant (spawn host only) |
| `next` | Execute next step |
| `run-all` | Run all remaining (stops on failure) |
| `run-until <step>` | Run up to step (stops on failure) |
| `jump <step>` | Move position without executing |
| `set-var KEY=value` | Override expansion variable |
| `list-steps` | Show steps with status |
| `logs [--step N] [--tail N] [--setup]` | View logs |

All commands are prefixed with `evergreen debug`.

## No-op Commands

These are automatically skipped in debug mode (expected, not errors): `host.create`, `host.list`, `generate.tasks`, `downstream_expansions.set`, `attach.results`, `attach.xunit_results`, `gotest.parse_files`, `attach.artifacts`, `papertrail.trace`, `keyval.inc`, `perf.send`, `s3.put`, `s3Copy.copy`. They show as "skipped" in `list-steps` output and during `next` execution — do not investigate or attempt to fix.

## Boundaries

- Do NOT run `run-all` or `run-until` without user confirmation.
- Do NOT attempt fixes for infrastructure issues (host networking, service outages) — report them.
- Do NOT modify files outside the project directory without explaining why.
- If root cause is in a dependency or upstream repo, explain what needs to change rather than patching locally.

## After Resolving

Summarize: (1) root cause, (2) what changed, (3) whether the fix is local-only or needs an upstream commit, (4) remaining unverified steps.
