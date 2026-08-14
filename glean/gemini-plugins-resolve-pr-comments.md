# resolve-pr-comments

**Category:** Cloud, DevOps & Infrastructure
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/before-you-merge/resolve-pr-comments/skills/resolve-pr-comments

## Description
Workflow for fetching, addressing, and resolving GitHub PR review feedback — both inline review threads and top-level PR comments. Use when asked to resolve PR comments, work through open review threads, or clean up unresolved GitHub review threads before merging.

---

# Resolving PR Review Comments

## Overview

Full workflow: fetch open feedback → read context → fix code → commit → **push** → resolve threads.

**Do not resolve threads until the fixing commit has been pushed.** Resolving before pushing leaves the PR in a misleading state where comments appear addressed but no code change is visible.

## Step 1: Fetch Open Feedback

**Inline review threads:**

```bash
${CLAUDE_SKILL_DIR}/scripts/fetch-threads.sh <owner> <repo> <pr-number>
# e.g. ${CLAUDE_SKILL_DIR}/scripts/fetch-threads.sh 10gen fern 857
```

Outputs one JSON object per unresolved thread with `id`, `path`, `line`, `outdated`, `author`, and `comments` (array of all comment bodies in the thread).
The `id` (format: `PRRT_kwDO...`) is required for resolving.

**Top-level PR comments:**

```bash
gh pr view <pr-number> --comments
```

Top-level comments may contain requested changes that aren't captured as inline threads. Address any action items found — there is no resolve step for these.

## Step 2: Address the Feedback

For each thread: read the file at the indicated path/line, make the fix.

After all edits: commit the changes, then push:

```bash
git push
```

**Wait for the push to complete before resolving threads.**

## Step 3: Resolve Threads

```bash
${CLAUDE_SKILL_DIR}/scripts/resolve-thread.sh <thread-id>
# e.g. ${CLAUDE_SKILL_DIR}/scripts/resolve-thread.sh PRRT_kwDOA...
```

Prints `true` on success. Run once per thread ID.

## Step 4: Verify

```bash
${CLAUDE_SKILL_DIR}/scripts/verify-threads.sh <owner> <repo> <pr-number>
# Should print: 0
```

## Key Rules

- **Never reply to bot/Copilot threads** — fix the code and resolve, no reply
- **Outdated threads**: if `outdated: true`, the code moved under the comment — resolve without fixing unless the concern still applies
- **Human reviewer threads**: fix, push, then resolve
- **Thread ID format**: `PRRT_kwDO...` — from `fetch-threads.sh`, never from REST

## Common Mistakes

- Resolving before pushing — threads look addressed but no code change is visible
- Using REST to resolve — REST has no resolve endpoint; these scripts use GraphQL
- Resolving by comment ID — must use the thread `id` field (`PRRT_...`), not the comment ID