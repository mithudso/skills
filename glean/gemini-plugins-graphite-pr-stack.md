# graphite-pr-stack

**Category:** General & Specialized Utilities
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/tooling/graphite-pr-stack/skills/graphite-pr-stack

## Description
Use when creating or updating a stacked PR series using Graphite (gt). Covers commit analysis, logical grouping, diff generation, branch reconstruction, and PR submission.

---

# Graphite PR Stack

> This skill covers the Graphite (gt) workflow for creating and updating stacked PRs. Graphite is nearly 100% generic across repositories. Check the `references/` directory for repo-specific build verification commands, pre-commit behavior, and PR title formats.

Create or update a series of stacked PRs using Graphite (`gt`). Each PR in the stack is independently buildable and represents a logical unit of work.

**CLI usage:** Always invoke the Graphite CLI as `gt`, never as a full path like `/opt/homebrew/bin/gt`. The short name avoids unnecessary permission prompts and is available on `$PATH` after standard installation.

**CLI over MCP:** If the `gt` CLI is installed (verify with `which gt`), always use it instead of any Graphite MCP tools. The CLI is the authoritative interface; MCP tools are a fallback only when the CLI is unavailable.

## Invocation

```
/graphite-pr-stack [create|update] [--base=<commit>] [--exclude=<glob>...] [--no-verify]
```

| Argument              | Description                                                                        |
| --------------------- | ---------------------------------------------------------------------------------- |
| `create`              | Create a new stack from the current branch (default if no stack exists)            |
| `update`              | Modify an existing Graphite stack                                                  |
| `--base=<commit>`     | Starting commit of the branch (default: merge-base with trunk)                     |
| `--exclude=<glob>...` | File patterns to split into a separate trailing PR (e.g., `.some-excluded-dir/**`) |
| `--no-verify`         | Skip build verification between branches                                           |

---

## Mode Selection

Determine whether to **create** or **update** based on context:

| Signal                                                                   | Mode       |
| ------------------------------------------------------------------------ | ---------- |
| Branch has no Graphite stack (`gt log` shows single branch or untracked) | **Create** |
| User says "stack this", "split into PRs", "create a PR stack"            | **Create** |
| `gt log` shows an existing multi-branch stack                            | **Update** |
| User says "update the stack", "add to the stack", "fix branch N"         | **Update** |
| User has new commits on an already-stacked branch                        | **Update** |

---

## Creating a New Stack

Follow the detailed workflow in `.agents/skills/graphite-pr-stack/references/creating-stacks.md` (from project root). High-level flow:

1. **Prerequisites** — Verify `gt` CLI, clean working tree, detect build system, create backup tag
2. **Analyze commits** — Classify commits as CODE, EXCLUDED, or MIXED; extract ticket prefix
3. **Group into PRs** — Propose logical groupings, get user approval
4. **Generate diffs** — Cumulative diffs per group using boundary commits
5. **Reconstruct stack** — Apply diffs onto trunk as Graphite branches via `gt create`
6. **Build verification** — Verify each code branch builds independently (see references/ for your repo's build command)
7. **Submit and describe** — Push with `gt submit --stack` and add PR descriptions

---

## Updating an Existing Stack

Follow the detailed workflow in `.agents/skills/graphite-pr-stack/references/updating-stacks.md` (from project root). Common operations:

- **Amend a branch** — `gt modify`, `gt modify --into`, or `gt absorb`
- **Add a branch** — `gt create` from the desired parent
- **Remove a branch** — `gt fold`, `gt pop`, or `gt delete`
- **Reorder branches** — `gt reorder`
- **Split a branch** — `gt split --by-commit`, `--by-file`, or `--by-hunk`
- **Sync with trunk** — `gt sync`
- **Restack before submitting** — After any branch modifications, always run `gt restack` before `gt submit` to ensure parent/child relationships are up to date
- **Push updates** — `gt submit --stack` or `gt ss -u` (update-only)

---

## Recovery

### Quick recovery: `gt undo`

If the most recent Graphite operation went wrong (e.g., a bad `gt create` or `gt modify`), undo it:

```bash
gt undo
```

This reverses the last Graphite mutation. Try this first — it's the simplest recovery path.

### Full recovery: backup tag

If multiple operations went wrong or `gt undo` isn't sufficient, restore from the backup tag created during stack creation:

```bash
git checkout <ORIG_BRANCH>
git reset --hard backup-pr-stack-<timestamp>
```

If the original branch was deleted:

```bash
git checkout -b <ORIG_BRANCH> backup-pr-stack-<timestamp>
```

Clean up partially-created stack branches:

```bash
gt log  # identify stack branches
# for each stack branch:
git branch -D <branch-name>
```

---

## Graphite Documentation

- **Index:** https://graphite.com/docs/llms.txt
- **Full docs:** https://graphite.com/docs/llms-full.txt

Fetch these when you need to look up unfamiliar commands or verify behavior not covered in this skill.

---

## Troubleshooting

Refer to `.agents/skills/graphite-pr-stack/references/troubleshooting.md` (from project root) for detailed guidance on common issues including:

- Empty diff files from shell redirection in sandbox environments
- Patch application failures from pre-commit hook reformatting
- Build failures from split dependencies and visibility modifiers
- Graphite state issues (detached HEAD, duplicate commits, remote divergence)

---

## Important Notes

- **Backup first.** When creating a stack, the backup tag is non-negotiable. The reconstruction is destructive.
- **Diffs, not cherry-picks.** Generating cumulative diffs between boundary commits is more reliable than cherry-picking interleaved commits.
- **Use `tee` for diff files.** Direct `>` redirection of `git diff` output may silently fail in sandboxed environments. Always pipe through `tee` and verify line counts.
- **Pre-commit hooks change files.** Hooks (formatters, import sorters, etc.) may modify files during `gt create`, making subsequent diffs stale. Always use `--3way` as a fallback for `git apply`.
- **Build each branch.** A diff that applies cleanly is not the same as code that builds. Verify every code branch independently (see references/ for your repo's build command).
- **The excluded-files PR is optional.** It exists to keep non-code artifacts out of the code PRs. It can be closed without merging.
- **Commits must be authored by the user.** Do not add `Co-Authored-By: Claude Code` trailers or set Claude as the commit author. The user is the author of all commits in the stack.

## Repo-Specific Context

Before running build verification or submitting stacks, check if a `references/` directory exists alongside this skill.
Look for a file matching your repo name (e.g., `references/mms.md`, `references/server.md`).
Load it to get the correct build verification command, pre-commit behavior, and PR title format for your repo.
Personal overrides in your local `references/` take precedence over central ones.