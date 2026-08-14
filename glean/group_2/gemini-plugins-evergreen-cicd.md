# evergreen-cicd

**Category:** Cloud, DevOps & Infrastructure
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/tooling/evergreen-cicd/skills/evergreen-cicd

## Description
Use when creating Evergreen patches, running tests on CI, or diagnosing CI/CD failures. Covers patch workflow, failure triage, and MCP tool patterns.

---

# Evergreen CI/CD Guide

> Before running commands, replace `<PROJECT>` with your project name (e.g., `mongodb-mongo-master`, `mms`). Check `references/` directory for your repo-specific project flag, task aliases, and local test commands.

This guide covers two workflows: **creating patches** (submitting code for CI validation) and **diagnosing failures** (investigating and fixing CI issues).

## Creating Patches with the Evergreen CLI

### Quick Reference — Common Patch Commands

**Run with an alias:**

```bash
evergreen patch -p <PROJECT> -a <alias> -f -y -u
# see references/ for aliases
```

**Run specific variants and tasks:**

```bash
evergreen patch -p <PROJECT> -v <variant> -t <task> -f -y -u
# see references/ for variants and tasks
```

### Essential Flags

| Flag           | Purpose                                                                                 |
| -------------- | --------------------------------------------------------------------------------------- |
| `-p <PROJECT>` | Target your project — see `references/` for the correct project name                    |
| `-v <variant>` | Select build variant(s) — repeatable or comma-separated                                 |
| `-t <task>`    | Select task(s) — repeatable or comma-separated. Use `all` for every task in the variant |
| `-a <alias>`   | Use a predefined alias instead of specifying variants/tasks manually                    |
| `-f`           | Finalize immediately (start running tasks right away)                                   |
| `-y`           | Skip confirmation prompts                                                               |
| `-u`           | Include uncommitted local changes                                                       |
| `-d "msg"`     | Set patch description                                                                   |
| `--ad`         | Auto-describe from last commit message                                                  |
| `--browse`     | Open patch URL in browser                                                               |

### Repeat/Retry Patterns

```bash
# Rerun the same tasks/variants as the latest patch
evergreen patch -p <PROJECT> --repeat -f -y -u

# Rerun only the FAILED tasks from the latest patch
evergreen patch -p <PROJECT> --repeat-failed -f -y -u

# Rerun tasks from a specific patch
evergreen patch -p <PROJECT> --repeat-patch <patch_id> -f -y -u
```

### Using Regex for Variant/Task Selection

**Caution**: `--rv`/`--rt` silently create empty patches if the regex doesn't match any variant+task intersection. Prefer explicit `-v`/`-t` flags. Use `evergreen list --variants -p <PROJECT>` and `evergreen list --tasks -p <PROJECT>` to discover valid names first.

```bash
# Regex variant matching (omit -f to inspect before finalizing)
evergreen patch -p <PROJECT> --rv "^pattern" --rt ".*" -y -u

# Then finalize after confirming tasks were matched
evergreen finalize-patch -i <patch_id>
```

### Patch Management

```bash
evergreen list-patches                         # List recent patches
evergreen list-patches -n 10                   # List 10 most recent
evergreen finalize-patch -i <patch_id>         # Finalize an unfinalized patch
evergreen cancel-patch -i <patch_id>           # Cancel a running patch
```

For complete variant/task/alias mapping and all CLI flags, consult **`.agents/skills/evergreen-cicd/references/patch-creation.md`**.

---

## Diagnosing CI/CD Failures

### MCP Tool Strategy

**ALWAYS use MCP tools first** for failure analysis:

1. `list_user_recent_patches_evergreen` — get patch status
2. `get_patch_failed_jobs_evergreen` — analyze failures
3. `get_task_test_results_evergreen` — get specific failed tests
4. `get_task_logs_evergreen` — only when job name and general failure are insufficient

IMPORTANT: If a job is failing on master, do not attempt to fix it.

### Failure Prioritization

1. **Fix compilation first** — cascading failures make other results unreliable
2. **Lint/format next** — quick wins that unblock commits
3. **Unit tests** — most likely related to code changes
4. **Integration tests** — check test logic, service dependencies, and configuration
5. **Other tasks** — investigate only after higher priorities are resolved

### Exclusions — DO NOT FIX

- **Known flaky tests** — if a test is flaky and unrelated to your changes, skip it
- **Infrastructure failures** (host allocation, archive failures) — escalate to Build team
- **Master branch failures** — do not attempt to fix; wait for upstream resolution
- **Nightly/cron tests** — not patch-blocking

### Deduplication

- Fix compilation tasks before investigating test failures (cascading)
- Multiple test failures in the same suite or module likely share a single root cause
- Failures unrelated to patch changes — wait for master fixes

For detailed per-task-type diagnosis guidance (commands, MCP strategies, debugging steps), consult **`.agents/skills/evergreen-cicd/references/failure-diagnosis.md`**.

For the full phase-by-phase investigation workflow (Phase 1-8: retrieve, analyze, approve, investigate, fix, verify, commit), consult **`.agents/skills/evergreen-cicd/references/investigation-workflow.md`**.

## Additional Resources

### Reference Files

- **`.agents/skills/evergreen-cicd/references/patch-creation.md`** — Complete CLI flag reference, variant/task mapping, alias definitions, patch management commands
- **`.agents/skills/evergreen-cicd/references/failure-diagnosis.md`** — Per-task-type diagnosis: compilation, unit tests, integration tests, lint/format. Includes local test commands and debugging steps.
- **`.agents/skills/evergreen-cicd/references/investigation-workflow.md`** — Phase-by-phase workflow for systematic failure investigation. Analysis format template, approval gates, verification checklist.

## Repo-Specific Context

Before running commands, check if a `references/` directory exists alongside this skill.
Look for a file matching your repo name (e.g., `references/mms.md`, `references/server.md`).
Load it to get the correct project flag, task aliases, and local test commands for your repo.
Personal overrides in your local `references/` take precedence over central ones.