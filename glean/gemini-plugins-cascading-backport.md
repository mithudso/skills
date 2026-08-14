# cascading-backport

**Category:** Databases, Data Engineering & Analytics
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/query-integration/cascading-backport/skills/cascading-backport

## Description
Backport a CVE security fix from master as a single squashed change across a chain of release-staging branches, following the MongoDB CVE process, with conflict resolution and a consistent 3-commit audit trail per branch.

---

# CVE Cascading Squash Backport

Apply a CVE security fix from master as a single squashed change to a chain of release-staging
branches, with each resolved branch feeding the next, and a consistent 3-commit audit trail per
branch.

## CVE Process Context

**Before doing anything else**, load the current CVE process guide from the internal wiki:

```
mcp__glean_default__read_document(urls=["https://wiki.corp.mongodb.com/spaces/SEC/pages/546766874/Interact+with+CVE+process+for+Server+engineering+teams"])
```

If the page is unavailable, surface that as an error — do not proceed without it. The guide covers
the full CVE lifecycle (filing, embargo, coordinated disclosure, advisory publishing) and takes
precedence over any assumptions embedded here. The backport procedure below covers only the
technical branching and patch-application steps; everything else (CVE ID assignment, JIRA
security-issue tracking, reporter coordination, disclosure timeline) is governed by that guide.

## When to Use

Use this skill when:

- You are backporting a **security fix** for a CVE to release-staging branches.
- The fix lives on `master` or a developer's branch (possibly across multiple commits) and must land on multiple major versions (e.g., v8.3, v8.2, v8.0, v7.0).
- You want a squash (single diff) rather than per-commit cherry-picks.
- You want reviewer-visible conflict resolution: raw markers, then explained resolution, then clean.

## Prerequisites

- CVE process steps up to and including "prepare the fix" are complete — consult the wiki guide loaded above.
- CLI tools set up and authenticated — see the `backport` skill's Important Notes. This skill lives in the 10gen/mongo repository.
- All target backport branches already exist (created earlier via `gh pr create`).
- Each branch has a dedicated git worktree (see `git-workflows` skill).

---

## Flags

| Flag                 | Effect                                                                                                                                      |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| `--skip-review-gate` | Skip the user review gate between Commit 2 and Commit 3. Useful for straightforward, targeted diffs where conflict resolutions are obvious. |

Pass flags when invoking the skill, e.g. "Run cascading-backport --skip-review-gate".

---

## Procedure

### Step 0: Collect context

```bash
git log origin/master..HEAD --oneline    # commits to squash
git diff origin/master...HEAD --stat    # files changed
gh pr list --author "@me" --state open  # existing backport PRs with their numbers
git worktree list                       # confirm worktrees for each version
```

Record:

- **Server ticket**: e.g. `SERVER-129617`
- **PR title suffix**: rest of the title after the ticket prefix
- **Versions + PR numbers**: e.g. v8.3→#56606, v8.2→#56607, v8.0→#56608, v7.0→#56610
- **Worktree paths**: one per version

**IMPORTANT — squash boundary check**: Before generating the patch, confirm that your squash
captures only YOUR ticket's changes and not unrelated commits that landed on master after you
branched. Run `git log --oneline origin/master..HEAD -- <affected_files>` and verify all commits are from
your ticket. Any unrelated commits that touched the same files will be included in the squash diff —
decide whether to include or exclude them before proceeding.

### Step 1: Generate the squash patch

```bash
git diff origin/master...HEAD > /tmp/${TICKET}-squashed.patch
```

This captures every committed change on the branch relative to master as a single unified diff.

### Step 2: Fetch staging branches

```bash
git fetch origin v8.3-staging v8.2-staging v8.0-staging v7.0-staging
```

---

## Per-branch loop (first: v8.3, cascade from there)

For **v8.3**, apply `/tmp/${TICKET}-squashed.patch`. For **v8.2**, apply
`git diff origin/v8.3-staging...HEAD` (generated after v8.3 is resolved). For **v8.0**, apply
`git diff origin/v8.2-staging...HEAD` (generated after v8.2 is resolved). For **v7.0**, apply
`git diff origin/v8.0-staging...HEAD` (generated after v8.0 is resolved).

This cascade improves apply success rate: adjacent versions share more code than master and an old
branch.

### Commit 1 — squash with conflict markers

```bash
cd /path/to/worktree/${VERSION}
git reset --hard origin/${VERSION}-staging

# Try 3-way merge first (produces git conflict markers):
git apply --3way <patch_file>

# Stage conflicted files WITH markers:
git add <files_with_conflicts>
git add <cleanly_applied_files>

git commit -m "SERVER-XXXXX [squashed] <title>

Squashed from <source-branch>. Conflict markers left in-place in N files
for review transparency (see next commit for resolution):
  - <file1>
  - <file2>
"
```

**If `--3way` fails with "does not exist in index" or too many rejects**, use `--reject`:

```bash
git apply --reject <patch_file>
# Stage cleanly applied files:
git add <files_that_applied>
# For each .rej file, insert hand-crafted markers in the target file:
#   // <<< CONFLICT-START (vX.Y): <why it failed>
#   //   Current code: ...
#   //   Patch wants: ...
#   // >>> CONFLICT-END
git add <files_with_hand_crafted_markers>
find . -name '*.rej' -delete  # clean up reject artifacts
git commit -m "SERVER-XXXXX [squashed] <title>

Hand-crafted conflict markers (// <<< CONFLICT) left in N files ..."
```

### Commit 2 — resolve with TMP-MERGE explanations

For each conflicted file:

1. Replace git conflict markers (`<<<<<<<`/`=======`/`>>>>>>>`) or hand-crafted markers with the
   correct resolved code.
2. Before resolving, run `git log --oneline origin/${VERSION}-staging -- <file>` and
   `git blame <file>` on the staging branch to find which commit (and ticket) introduced the
   conflicting code. Use that context to inform the resolution.
3. Immediately above each resolved block, add:

   ```
   // TMP-MERGE: <one-sentence explanation of what conflicted and how you resolved it>
   // TMP-MERGE-REF: <clickable link to the ticket or commit that introduced the conflicting code>
   ```

   For YAML files use `#` comments. Format the reference as a clickable link:
   - Jira ticket: `https://jira.mongodb.org/browse/SERVER-XXXXX`
   - GitHub commit: `https://github.com/10gen/mongo/commit/<full-SHA>`

   Prefer the Jira link when a ticket is known — it carries more context than a raw SHA.

4. Stage and commit:

```bash
git add <resolved_files>
git commit -m "Resolve merge conflicts for SERVER-XXXXX ${VERSION} backport

<file1>: <explanation>
<file2>: <explanation>
...

See TMP-MERGE comments for inline context."
```

### User review gate (before Commit 3)

If `--skip-review-gate` was passed, proceed directly to Commit 3.

Otherwise, **pause here and ask the user** whether they want to review the conflict resolutions
before the TMP-MERGE comments are stripped. Present the diff of the resolution commit:

```bash
git show HEAD --stat
git diff HEAD~1 HEAD -- <resolved_files>
```

Offer to walk through each resolved block. Only proceed to Commit 3 once the user confirms the
resolutions look correct. This is the last point at which the reasoning behind each resolution is
readable inline.

### Commit 3 — remove TMP-MERGE comments

```python
import re, sys
for path in sys.argv[1:]:
    content = open(path).read()
    # Use separate passes — do NOT combine // and # into one alternation,
    # as that matches #include lines and silently deletes them.
    cleaned = re.sub(r'[ \t]*//[^\n]*TMP-MERGE[^\n]*\n(?:[ \t]*//[^\n]*\n)*', '', content)
    cleaned = re.sub(r'[ \t]*#[^\n]*TMP-MERGE[^\n]*\n(?:[ \t]*#[^\n]*\n)*', '', cleaned)
    open(path, 'w').write(cleaned)
```

**After stripping**, visually check that no surrounding comments were eaten. The regex consumes
every `//` line following a TMP-MERGE line, so a legitimate comment block immediately after a
TMP-MERGE block (e.g. a `//---` separator) will be silently deleted. Re-add any such lines before
committing.

```bash
git add <cleaned_files>
git commit -m "Remove TMP-MERGE comments"
```

### Format changed files

Formatting failures are a frequent source of CI failures on backport patches. After stripping
TMP-MERGE comments, run the format command for the target branch before pushing.

See `references/format-by-branch.md` for the full branch-by-branch reference. Quick lookup:

| Branch | Command |
|---|---|
| `v6.0-staging` | `python buildscripts/clang_format.py format-my origin/v6.0` |
| `v7.0-staging` | `python buildscripts/clang_format.py format-my origin/v7.0` |
| `v8.0-staging` | `bazel run //:format -- --origin-branch origin/v8.0` |
| `v8.2-staging` | `bazel run //:format -- --origin-branch origin/v8.2` |
| `v8.3-staging` | `bazel run //:format -- --origin-branch origin/v8.3` |

If the command produces changes, stage and commit them before pushing:

```bash
git add -u
git commit -m "Formatting"
```

### Push

```bash
git push --force origin <your-github-username>/SERVER-XXXXX-${VERSION}
```

Force push is safe: the old commits are superseded by the new squash approach.

### Generate cascade patch for the next version

```bash
git diff origin/${VERSION}-staging...HEAD > /tmp/${TICKET}-${VERSION}-resolved.patch
```

---

## Test file multiversion compatibility

Any test file modified by the backport that asserts new behavior will fail in multiversion CI suites
(e.g. `query_shape_hash_stability_last_continuous_new_old`) because burn-in picks up changed files
and runs them against the old binary. Make sure the added/updated tests are in
'backports_required_for_multiversion_tests.yml'. This is only relevant for tests in directories
which have such multiversion passthrough suites.

**The right tag is `backport_required_multiversion`** (not `multiversion_incompatible`):

- `backport_required_multiversion` — temporary; once all backports land and the tag is removed,
  tests run in multiversion automatically → better long-term coverage.
- `multiversion_incompatible` — permanent; tests never run in multiversion even after full backport.

Add to the test file's `@tags` block on all backport branches:

```js
 * @tags: [
 *   requires_fcv_NN,
 *   backport_required_multiversion,
 * ]
```

**Also check**: not all multiversion suites exclude `backport_required_multiversion`. The
`query_shape_hash_stability` override only excluded `multiversion_incompatible` (as of this
writing). If that suite picks up your test files, add `backport_required_multiversion` to the
relevant `exclude_with_any_tags` blocks in
`buildscripts/resmokeconfig/matrix_suites/overrides/query_shape_hash_stability.yml` **on both the
backport branches and master**. This is a one-line fix per block and is worth doing on master as a
general improvement.

---

## Running Evergreen patches

Follow the `backport` skill's steps 7–8 (patch build + PR comment). One difference: omit `-f` on
the first run so the user can inspect and finalize manually before CI starts. Add `-f` only when
confident.

---

## Structural differences between versions

Older branches often diverge structurally. Common adaptations:

| Situation                                                             | Resolution                                                                                       |
| --------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| `window_function_percentile.cpp` absent in old branch                 | Apply changes to `.h` file inline                                                                |
| `query/util/BUILD.bazel` absent or uses only `exports_files`          | Create it with `load() + mongo_cc_library`; put `load()` first                                   |
| `mongo_cc_library` dep path changed (e.g. `sort_pattern` moved)       | Keep old path; only add new dep                                                                  |
| YAML uses 2-space not 4-space indentation                             | Rewrite entry insertion logic for local indentation                                              |
| File uses `platform/basic.h` catch-all instead of individual includes | Don't remove the "obsolete" includes; just add the new one                                       |
| Class hierarchy differs (method only on base class in old branch)     | Note in TMP-MERGE; add free function to the appropriate header instead                           |
| Header added to master after staging branches forked                  | Cannot use it in backport — adapt the code to work without it, or backport the header separately |

**Header dependency trap**: When adding a new utility header (e.g. `represent_as_util.h`), check
whether every header it `#include`s exists on the staging branches. If a dependency (e.g.
`mongo/util/demangle.h`) was added to master after the branch forked, Bazel's `auto_header` system
will generate a dep on `//src/mongo/util/.auto_header:demangle_h` which won't exist, breaking
clang-tidy. Strip the dependency from the backport version of the header and simplify the
implementation accordingly (e.g. drop the type-name demangling from the error message).

In these cases, commit 1 uses **hand-crafted markers** with a descriptive comment:

```cpp
// <<< CONFLICT-START (v7.0): <what the patch tried to do and why it failed>
//   Current code: <original content summary>
//   Patch wants: <desired change summary>
// >>> CONFLICT-END
```

---

## Verification checklist

After each branch:

- [ ] `git log --oneline -4` shows exactly 3 new commits on top of staging tip (or 4 if formatting needed a separate commit)
- [ ] No `<<<<<<`/`>>>>>>>`/`<<< CONFLICT` markers remain in any file
- [ ] No TMP-MERGE comments remain after commit 3
- [ ] No surrounding legitimate comments were eaten by the TMP-MERGE regex
- [ ] Diff contains only changes from THIS ticket (no unrelated commits bundled in)
- [ ] Formatting applied and committed (see `references/format-by-branch.md`)
- [ ] Modified test files have `backport_required_multiversion` tag
- [ ] Force push succeeded
- [ ] Cascade patch generated for next version
- [ ] Evergreen patch submitted (finalized or left for user) and URL posted to PR

---

## Notes

- **Don't squash the commits on the branch** — `git diff origin/master...HEAD` captures them all
  without rewriting history on master.
- **The cascade direction** flows from newest (v8.3) to oldest (v7.0): each resolved diff is closer
  in code to the next target than the original master diff is.
- **`--3way` vs `--reject`**: prefer `--3way` (produces real git conflict markers). Fall back to
  `--reject` when files don't exist in the target index.
- **TMP-MERGE regex trap**: don't use `(?://|#)` in cleanup — it matches `#include` lines and can
  eat legitimate comment blocks. Use separate passes for `//` and `#`, and visually inspect after.
- **Cancelling stale patches**: `evergreen cancel-patch -id <version_id>`.