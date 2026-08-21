# Updating an Existing Stack

Workflows for modifying a Graphite stack after it has been created. All operations auto-restack upstack branches unless noted otherwise.

---

## Before Making Changes

Confirm the current stack state:

```bash
gt log          # full stack with PR status
gt log short    # branch names only
```

If the stack is out of date with trunk, sync first:

```bash
gt sync
```

This fetches trunk, cleans up merged branches, and restacks the entire stack. Resolve any conflicts with `gt add <file>` + `gt continue`.

---

## Amending a Branch

Three strategies depending on where you are and where the fix belongs.

### Current branch

Make changes, then amend:

```bash
gt modify -a                        # stage all + amend
gt modify -a -m "updated message"   # stage all + amend + update commit message
gt modify -c -am "extra commit"     # add a new commit instead of amending
```

### Specific downstack branch (without switching)

Stage changes, then target the branch:

```bash
git add <files>
gt modify --into <target-branch>
```

This amends the staged changes into the target branch and restacks everything above it. Useful when you discover a fix for branch 1 while working on branch 3.

### Multiple branches at once

Stage all changes, then let Graphite route them:

```bash
git add <files>
gt absorb --dry-run    # preview which branch gets which changes
gt absorb              # apply
```

`gt absorb` matches each hunk to the branch where the surrounding code was last modified. Hunks that can't be deterministically matched are left unstaged. Use `gt absorb --patch` to interactively select which hunks to absorb.

---

## Adding a Branch

Navigate to the branch that should be the **parent** of the new branch, then create:

```bash
gt checkout <parent-branch>
# make changes
gt create <new-branch-name> -am "commit message"
```

The new branch is inserted between the parent and its previous child. All upstack branches are restacked automatically.

To add a branch at the **top** of the stack:

```bash
gt top
# make changes
gt create <new-branch-name> -am "commit message"
```

---

## Removing a Branch

Three options depending on what you want to happen to the changes:

### Fold into parent (keep changes, merge into parent branch)

```bash
gt checkout <branch-to-remove>
gt fold
```

The branch is deleted and its changes are squashed into the parent branch. Use `gt fold --close` to also close the associated PR.

### Pop (keep changes in working tree, delete branch)

```bash
gt checkout <branch-to-remove>
gt pop
```

The branch is deleted but its changes remain as uncommitted modifications. Useful when you want to redistribute the changes manually.

### Delete (discard changes entirely)

```bash
gt delete <branch-name>
```

Use `gt delete --close` to also close the associated PR.

---

## Reordering Branches

Navigate to the top of the stack, then reorder:

```bash
gt top
gt reorder
```

This opens an editor showing each branch as a line. Rearrange the lines to restructure the stack. After saving, Graphite rebases everything to match the new order.

Resolve any conflicts with `gt add <file>` + `gt continue`.

---

## Splitting a Branch

If a branch has grown too large, split it into multiple branches:

```bash
gt checkout <branch-to-split>
gt split --by-commit     # split at commit boundaries (interactive)
gt split --by-hunk       # split by selecting hunks (interactive)
gt split --by-file <pattern>  # split files matching pattern into new parent (non-interactive)
```

`--by-file` is the only non-interactive mode. It extracts files matching the pattern into a new parent branch. Can be repeated with multiple patterns:

```bash
gt split --by-file "src/api/**" --by-file "src/models/**"
```

---

## Pushing Updates

After making changes to the stack:

### Update all PRs

```bash
gt submit --stack --no-interactive
```

If remote branches have diverged (e.g., after amending):

```bash
gt submit --stack --no-interactive --force
```

### Update only existing PRs (skip branches without PRs)

```bash
gt submit --stack --update-only
```

Short form: `gt ss -u`.

---

## Build Verification After Updates

After amending branches, verify that each code branch still builds. The same process as initial creation applies — see `creating-stacks.md`, Step 5.

Key difference: when fixing build failures in an existing stack, prefer `gt modify --into` or `gt absorb` over checking out each branch individually. This keeps you in place and auto-restacks.

---

## Common Update Scenarios

### Addressing PR review feedback

1. Check out the branch with feedback: `gt checkout <branch>`
2. Make the requested changes
3. Amend: `gt modify -a`
4. Push: `gt submit --stack --no-interactive --force`

### Rebasing after trunk changes

1. Sync: `gt sync`
2. Resolve conflicts if any: `gt add <file>` + `gt continue`
3. Verify builds if needed
4. Push: `gt submit --stack --no-interactive --force`

### Moving changes between branches

1. Pop the source branch: `gt checkout <source>` + `gt pop`
2. Stage the changes you want to move: `git add <files>`
3. Absorb into the right branches: `gt absorb`
4. Discard remaining changes if any: `git checkout -- .`

### Adding a hotfix to the bottom of the stack

1. Go to the bottom: `gt bottom`
2. Move down to trunk: `gt down`
3. Create the hotfix branch: `gt create <hotfix-branch> -am "fix: ..."`
4. Restack: all upstack branches restack automatically
5. Push: `gt submit --stack --no-interactive --force`
