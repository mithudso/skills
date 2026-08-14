---
name: git-workflows
description: >-
  Use when working with git worktrees, GitHub CLI, branch management,
  or PR workflows. Covers gh CLI commands and git worktree patterns
  for parallel development.
source: 10gen/fern
license: Internal
compatibility: claude-code, cursor, gemini-cli
mongodb:
  team: devprod-bv
  owner: srdjan.pajic@mongodb.com
  internal: true
---

# Using Git Worktrees

## Overview

Git worktrees create isolated workspaces sharing the same repository, allowing work on multiple branches simultaneously without switching.

**Core principle:** Systematic directory selection + safety verification = reliable isolation.

**Announce at start:** "I'm using the git-workflows skill to set up an isolated workspace."

## Directory Selection Process

Follow this priority order:

### 1. Check Existing Directories

```bash
# Check in priority order
ls -d .worktrees 2>/dev/null     # Preferred (hidden)
ls -d worktrees 2>/dev/null      # Alternative
```

**If found:** Use that directory. If both exist, `.worktrees` wins.

### 2. Check CLAUDE.md

```bash
grep -i "worktree.*director" CLAUDE.md 2>/dev/null
```

**If preference specified:** Use it without asking.

### 3. Ask User

If no directory exists and no CLAUDE.md preference:

```
No worktree directory found. Where should I create worktrees?

1. .worktrees/ (project-local, hidden)
2. ~/.config/superpowers/worktrees/<project-name>/ (global location)

Which would you prefer?
```

## Safety Verification

### For Project-Local Directories (.worktrees or worktrees)

**MUST verify directory is ignored before creating worktree:**

```bash
# Check if directory is ignored (respects local, global, and system gitignore)
git check-ignore -q .worktrees 2>/dev/null || git check-ignore -q worktrees 2>/dev/null
```

**If NOT ignored:**

Per Jesse's rule "Fix broken things immediately":
1. Add appropriate line to .gitignore
2. Commit the change
3. Proceed with worktree creation

**Why critical:** Prevents accidentally committing worktree contents to repository.

### For Global Directory (~/.config/superpowers/worktrees)

No .gitignore verification needed - outside project entirely.

## Creation Steps

### 1. Detect Project Name

```bash
project=$(basename "$(git rev-parse --show-toplevel)")
```

### 2. Create Worktree

```bash
# Determine full path
case $LOCATION in
  .worktrees|worktrees)
    path="$LOCATION/$BRANCH_NAME"
    ;;
  ~/.config/superpowers/worktrees/*)
    path="~/.config/superpowers/worktrees/$project/$BRANCH_NAME"
    ;;
esac

# Create worktree with new branch
git worktree add "$path" -b "$BRANCH_NAME"
cd "$path"
```

### 3. Run Project Setup

Auto-detect and run appropriate setup:

```bash
# Node.js
if [ -f package.json ]; then npm install; fi

# Rust
if [ -f Cargo.toml ]; then cargo build; fi

# Python
if [ -f requirements.txt ]; then pip install -r requirements.txt; fi
if [ -f pyproject.toml ]; then poetry install; fi

# Go
if [ -f go.mod ]; then go mod download; fi
```

### 4. Verify Clean Baseline

Run tests to ensure worktree starts clean:

```bash
# Examples - use project-appropriate command
npm test
cargo test
pytest
go test ./...
```

**If tests fail:** Report failures, ask whether to proceed or investigate.

**If tests pass:** Report ready.

### 5. Report Location

```
Worktree ready at <full-path>
Tests passing (<N> tests, 0 failures)
Ready to implement <feature-name>
```

## Quick Reference

| Situation | Action |
|-----------|--------|
| `.worktrees/` exists | Use it (verify ignored) |
| `worktrees/` exists | Use it (verify ignored) |
| Both exist | Use `.worktrees/` |
| Neither exists | Check CLAUDE.md → Ask user |
| Directory not ignored | Add to .gitignore + commit |
| Tests fail during baseline | Report failures + ask |
| No package.json/Cargo.toml | Skip dependency install |

## Common Mistakes

### Skipping ignore verification

- **Problem:** Worktree contents get tracked, pollute git status
- **Fix:** Always use `git check-ignore` before creating project-local worktree

### Assuming directory location

- **Problem:** Creates inconsistency, violates project conventions
- **Fix:** Follow priority: existing > CLAUDE.md > ask

### Proceeding with failing tests

- **Problem:** Can't distinguish new bugs from pre-existing issues
- **Fix:** Report failures, get explicit permission to proceed

### Hardcoding setup commands

- **Problem:** Breaks on projects using different tools
- **Fix:** Auto-detect from project files (package.json, etc.)

## Example Workflow

```
You: I'm using the git-workflows skill to set up an isolated workspace.

[Check .worktrees/ - exists]
[Verify ignored - git check-ignore confirms .worktrees/ is ignored]
[Create worktree: git worktree add .worktrees/auth -b feature/auth]
[Run npm install]
[Run npm test - 47 passing]

Worktree ready at /Users/jesse/myproject/.worktrees/auth
Tests passing (47 tests, 0 failures)
Ready to implement auth feature
```

## Red Flags

**Never:**
- Create worktree without verifying it's ignored (project-local)
- Skip baseline test verification
- Proceed with failing tests without asking
- Assume directory location when ambiguous
- Skip CLAUDE.md check

**Always:**
- Follow directory priority: existing > CLAUDE.md > ask
- Verify directory is ignored for project-local
- Auto-detect and run project setup
- Verify clean test baseline

## Integration

**Called by:**
- **brainstorming** (Phase 4) - REQUIRED when design is approved and implementation follows
- **subagent-driven-development** - REQUIRED before executing any tasks
- **executing-plans** - REQUIRED before executing any tasks
- Any skill needing isolated workspace

**Pairs with:**
- **finishing-a-development-branch** - REQUIRED for cleanup after work complete

---

# GitHub CLI for PRs

## Create a PR

**Always create PRs in draft mode** unless the user explicitly asks for a ready-for-review PR. Use a heredoc for the body. See `commit-conventions` skill for the full PR template and format.

```bash
# Create draft PR (default — always use --draft)
gh pr create --draft --title "DEVPROD-####: Description" --body "$(cat <<'EOF'
<fill in PR template — see commit-conventions skill>
EOF
)"

# Push and create in one step (if not yet pushed)
git push -u origin HEAD && gh pr create --draft --title "DEVPROD-####: Description" --body "..."

# Mark ready for review when explicitly asked
gh pr ready
```

## View PR Details

```bash
# View current branch's PR
gh pr view

# View specific PR by number
gh pr view 857

# View in browser
gh pr view 857 --web

# View with comments
gh pr view 857 --comments

# JSON output for specific fields
gh pr view 857 --json title,body,state,reviews,comments
```

## Fetch Review Comments

```bash
# View PR review comments (inline code comments)
gh api repos/10gen/fern/pulls/857/comments --jq '.[] | {path: .path, line: .line, body: .body, user: .user.login}'

# View PR reviews (approve/request changes)
gh api repos/10gen/fern/pulls/857/reviews --jq '.[] | {state: .state, body: .body, user: .user.login}'

# View issue-level comments (conversation tab)
gh api repos/10gen/fern/pulls/857/comments --jq '.[] | {body: .body, user: .user.login, created_at: .created_at}'

# All of the above combined
gh pr view 857 --json reviews,comments --jq '{reviews: [.reviews[] | {state: .state, body: .body, author: .author.login}], comments: [.comments[] | {body: .body, author: .author.login}]}'
```

## Check CI Status

```bash
# View checks for current PR
gh pr checks

# View checks for specific PR
gh pr checks 857

# Watch checks until complete
gh pr checks 857 --watch
```

## Respond to Reviews

```bash
# Add a comment to a PR
gh pr comment 857 --body "Addressed feedback in latest commit"

# Request re-review
gh pr edit 857 --add-reviewer username
```

## List PRs

```bash
# List open PRs
gh pr list

# List your PRs
gh pr list --author @me

# List PRs needing review
gh pr list --search "review-requested:@me"
```
