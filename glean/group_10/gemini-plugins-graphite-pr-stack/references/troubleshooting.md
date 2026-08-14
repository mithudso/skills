# Troubleshooting: Graphite PR Stack

Common issues encountered when splitting branches into PR stacks, with root causes and fixes.

---

## Diff Generation

### Empty diff files from `git diff > file`

**Symptom:** `git diff A..B` produces output when piped (`| wc -l` shows hundreds of lines) but writing to a file with `>` produces a 0-byte file.

**Root cause:** Sandbox environments may intercept file descriptor redirection. The `git diff` process runs correctly, but the shell redirect is silently blocked.

**Fix:** Always use `tee` instead of `>`:

```bash
# BAD — may produce empty files
git diff A..B -- ':!.planning' > /tmp/phase1.diff

# GOOD — reliable in all environments
git diff A..B -- ':!.planning' | tee /tmp/phase1.diff | wc -l
```

### Empty diff output from `git diff --stat`

**Symptom:** `git diff A..B --stat` shows nothing, but the actual diff has content.

**Root cause:** The sandbox blocks reads to sensitive files (`.pem`, `.key`, `.p12`, etc.). When `git diff` encounters an `EPERM` on any file in the repo, `--stat` may silently return empty. Piped output still works because the sandbox handles pipe writes differently.

**Fix:** Run git commands with sandbox disabled when generating diffs:

```bash
# Use dangerouslyDisableSandbox: true for git diff operations
```

---

## Diff Application

### `patch does not apply` errors

**Symptom:** `git apply` fails with `error: patch failed` on a specific file.

**Root cause:** Pre-commit hooks (formatters, import sorters, dep sorters) modified files during a previous `gt create`, so the file content on disk no longer matches the diff's context lines.

**Fix:** Use three-way merge:

```bash
git apply --3way /path/to/phase.diff
```

This falls back to three-way merge for conflicting hunks. If conflicts remain, resolve manually — typically the difference is whitespace or line ordering.

### Conflicts from excluded-file removal in code diffs

**Symptom:** A code diff removes a line that was only present because it was added in a phase that also touched excluded files.

**Example:** Phase 1 adds an entry to a config. Phase 2 removes it. When you fix phase 1 to not add it, phase 2's diff has a removal that no longer applies.

**Fix:** After fixing the root cause in the earlier branch and running `gt modify`, the restack will either auto-resolve or produce a clean conflict. For the clean conflict, the resolution is usually to take the current branch's state (the entry was already absent).

---

## Build Failures

### General: Missing dependency after split

**Symptom:** Build fails with "cannot find symbol", "undefined reference", "module not found", or similar for a class/module/symbol that exists in another PR in the stack.

**Root cause:** The original branch had all changes in one commit, so deps were implicitly available. After splitting, a file in PR 2 may depend on a target that was only modified in PR 1.

**Fix:** Verify the dependency manifest (BUILD.bazel, Cargo.toml, etc.) includes all necessary targets for the files in each PR. Add missing deps to the PR that needs them.

### Bazel: Missing `deps` in BUILD.bazel after splitting C++ targets

**Symptom:** Bazel build fails with `undefined reference` or `use of undeclared identifier` for a symbol defined in a library that was split into a different PR.

**Root cause:** When a C++ library and its consumer were modified in the same branch, the dependency was implicit. After splitting into separate PRs, the consumer's `BUILD.bazel` may be missing the `deps` entry for the library target.

**Example:**

```
# PR 1 adds src/lib/util.{h,cpp} with a new BUILD.bazel target
# PR 2 modifies src/app/processor.cpp to #include "lib/util.h"
# PR 2's BUILD.bazel is missing: deps = ["//src/lib:util"]
```

**Fix:** Add the missing dependency to the consumer's `BUILD.bazel`:

```python
cc_library(
    name = "processor",
    srcs = ["processor.cpp"],
    hdrs = ["processor.h"],
    deps = [
        "//src/lib:util",  # added — split from original branch
        # ... other deps
    ],
)
```

Use `gt modify --into <branch>` to amend the fix into the correct branch if it belongs to a downstack PR.

---

## Graphite Operations

### `Cannot perform this operation without a branch checked out`

**Symptom:** `gt create` fails with this error after checking out trunk.

**Root cause:** You're in detached HEAD state, not on the trunk branch. This can happen if `gt untrack` or a failed operation left the repo in a bad state.

**Fix:** Explicitly checkout the trunk branch:

```bash
git checkout master  # or main
git rev-parse --abbrev-ref HEAD  # verify it says "master", not "HEAD"
```

### Duplicate commits after restack

**Symptom:** `gt log` shows two commits on a branch after `gt modify` + restack.

**Root cause:** When restacking fails partway and is retried, the rebase may create duplicate commits for branches that didn't conflict.

**Fix:** Squash the duplicate:

```bash
gt checkout <branch-with-duplicate>
git reset --soft HEAD~1
git commit --amend --no-edit
```

### `Branch has been updated remotely` on submit

**Symptom:** `gt submit --stack` fails because a branch was pushed in the initial submit and then locally amended.

**Fix:** Use `--force` to override:

```bash
gt submit --stack --no-interactive --force
```

---

## GPG Signing

### Commits not signed after `gt create`

**Symptom:** `git log --format="%G?"` shows `N` (no signature) instead of `G`.

**Root cause:** `commit.gpgsign` is not set globally or the GPG agent is not running.

**Fix:**

```bash
git config --global commit.gpgsign true
gpg --list-secret-keys  # verify key exists
```

If `gt create` already ran, amend the commit to re-sign:

```bash
gt modify  # will re-sign with current config
```
