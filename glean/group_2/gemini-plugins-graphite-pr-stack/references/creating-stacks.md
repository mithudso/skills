# Creating a New Stack

Detailed workflow for splitting a feature branch into a series of stacked PRs.

---

## Step 0: Prerequisites

1. **Verify Graphite CLI** is installed and the repo is initialized:
   ```bash
   gt --version
   ```
2. **Check GPG signing config** — note whether it's enabled, don't require it:
   ```bash
   git config --get commit.gpgsign  # note: true or empty
   ```
   If GPG signing is enabled, all commits in the stack will be automatically signed. If the user explicitly requests GPG signing and it's not configured, stop and help them set it up. Otherwise, proceed without it.
3. **Verify clean working tree:**
   ```bash
   git status  # must be clean
   ```
4. **Identify trunk branch** (usually `master` or `main`):
   ```bash
   gt log  # confirm trunk is recognized
   ```
5. **Store the current branch name and HEAD** for backup:
   ```bash
   ORIG_BRANCH=$(git branch --show-current)
   ORIG_HEAD=$(git rev-parse HEAD)
   git tag backup-pr-stack-$(date +%s) HEAD
   ```
6. **Detect the build system** for later verification:
   - If `BUILD.bazel` or `BUILD` files exist: Bazel
   - If `package.json` exists: npm/pnpm/yarn
   - If `Cargo.toml` exists: Cargo
   - If `go.mod` exists: Go
   - If `pom.xml` exists: Maven
   - Note the build system for Step 5. If none detected or `--no-verify` is passed, skip build verification.

If Graphite CLI is missing or the working tree is dirty, stop and inform the user.

---

## Step 1: Analyze Commits

Determine the base commit and enumerate all commits on the branch.

1. **Determine base commit.** Use `--base` if provided, otherwise:
   ```bash
   git merge-base <trunk> HEAD
   ```
2. **List all commits** from base to HEAD:
   ```bash
   git log --oneline <base>..HEAD
   ```
3. **Classify each commit** by checking which files it touches:

   | Classification | Rule                                               |
   | -------------- | -------------------------------------------------- |
   | **CODE**       | All changed files are outside `--exclude` patterns |
   | **EXCLUDED**   | All changed files match `--exclude` patterns       |
   | **MIXED**      | Contains both code and excluded files              |

   Use `git diff-tree --no-commit-id --name-only -r <hash>` to get files per commit.

4. **Extract the ticket prefix** from the branch name (e.g., `SERVER-12345` from `SERVER-12345-rate-limiter-replacement`). This prefix will be used for PR titles and branch names. If no ticket prefix is found, use short descriptive names instead.

5. **Present the classification** to the user before proceeding.

---

## Step 2: Group into PRs

Group the classified commits into logical PRs. This is the most critical step — get it right before touching any branches.

### If the user provided explicit groupings

The user may have already told you how to split (e.g., "put the infrastructure in one PR, the integration in another"). Honor their groupings exactly. Map their descriptions to the commits from Step 1 and confirm your mapping is correct before proceeding.

### If auto-detecting groupings

1. **Read commit messages and diffs** to identify logical phases. Look for:

   - Conventional commit prefixes (`feat`, `fix`, `refactor`, `test`, `chore`, `docs`)
   - Phase numbering in messages (e.g., `01-01`, `phase-02`)
   - Dependency relationships (does commit B depend on files created in commit A?)

2. **Propose PR groupings** as a numbered list. Each group should:

   - Represent a single logical unit of work (feature, cleanup, refactor, etc.)
   - Be independently buildable on top of the previous PR
   - Have a clear, descriptive title with the ticket prefix

3. **Excluded files go in a final, separate PR** with a `[do not merge]` title suffix. This PR is always last in the stack.

4. **Present the proposed stack to the user** and wait for confirmation before proceeding:

   ```
   Proposed PR stack:
   1. TICKET-ID: <title> (N commits)
   2. TICKET-ID: <title> (N commits)
   ...
   N. TICKET-ID: Excluded artifacts [do not merge] (excluded files)
   ```

   The user may reorder, merge, or split groups. Incorporate feedback and re-present until approved.

---

## Step 3: Generate Diffs

For each approved PR group, generate a cumulative diff that captures exactly the code changes for that group.

### Diff boundary strategy

For each PR group, identify the **boundary commits** — the last commit in the previous group and the last commit in the current group. Then generate the diff:

```bash
# For code PRs: exclude the --exclude patterns
git diff <prev_boundary>..<curr_boundary> -- ':!.some-excluded-dir' [other exclusions]

# For the excluded-files PR: include ONLY the excluded patterns
git diff <trunk>..<HEAD> -- .some-excluded-dir [other inclusions]
```

**Why this works:** Since excluded-only commits don't touch code files, the code diff between boundary commits is clean even when excluded commits are interspersed. This is more reliable than cherry-picking interleaved commits, which fails when code and non-code commits alternate.

### Saving diffs

**IMPORTANT:** `git diff` output piped through `tee` is reliable; direct shell redirection (`>`) may silently produce empty files in some sandbox environments. Always use:

```bash
git diff <range> -- <pathspec> | tee /path/to/output.diff | wc -l
```

Verify each diff file has non-zero line count before proceeding.

Save diffs to `$TMPDIR/pr-stack/` (or `/private/tmp/claude-<uid>/pr-stack/`).

### Validation

After generating all diffs:

1. Verify each diff is non-empty (except the excluded PR which may be empty if no excluded files exist).
2. Verify no code file appears in more than one code diff.
3. Verify the union of all diffs accounts for every changed file on the branch.

---

## Step 4: Reconstruct Stack

Build the Graphite stack from the generated diffs. This is destructive — the backup tag from Step 0 is the safety net.

### Setup

1. **Untrack the current branch** from Graphite if tracked:
   ```bash
   gt untrack  # only if currently tracked
   ```
2. **Checkout trunk:**
   ```bash
   git checkout <trunk>
   ```
3. **Verify you are on trunk** (not detached HEAD):
   ```bash
   git rev-parse --abbrev-ref HEAD  # must equal trunk name, not "HEAD"
   ```

### Apply each PR

For each PR group (in stack order):

1. **Apply the diff:**

   ```bash
   git apply /path/to/phaseN.diff
   ```

   If `git apply` fails, retry with `--3way` for three-way merge:

   ```bash
   git apply --3way /path/to/phaseN.diff
   ```

   If conflicts occur, resolve them manually. Common causes:

   - Pre-commit hooks reformatted files in a previous branch (e.g., BUILD file dep sorting, import ordering)
   - A file was added in one phase and modified in another with overlapping context

2. **Stage the changes.** Stage only relevant files — use `--exclude` patterns to separate code from excluded files:

   ```bash
   git add -A -- ':!.some-excluded-dir'  # for code PRs
   git add -A -- .some-excluded-dir      # for excluded PR
   ```

3. **Create the Graphite branch:**

   ```bash
   gt create <branch-name> -m "<commit message>"
   ```

   The first line of the commit message becomes the PR title. Use the format:

   ```
   TICKET-ID: <descriptive title>

   <body with details of what this PR contains>
   ```

4. **Verify GPG signature** if signing is enabled:
   ```bash
   git log -1 --format="%G?"  # "G" = good signature
   ```

### Branch naming convention

Use the ticket prefix with a numbered suffix and short description:

```
TICKET-ID-01-short-description
TICKET-ID-02-short-description
TICKET-ID-99-excluded-artifacts
```

The excluded-files branch uses `99` to sort last.

---

## Step 5: Build Verification

**Skip this step if `--no-verify` was passed or no build system was detected in Step 0.**

For each code branch in the stack (not the excluded-files branch), verify the build passes. This catches issues that a cleanly-applied diff can hide: missing dependencies, visibility modifiers, and BUILD file exports that break when code is split across PRs.

### Process

1. **Checkout the branch:**

   ```bash
   gt checkout <branch-name>
   ```

2. **Run the build.** Determine targets from the changed files and the build system detected in Step 0:

   | Build system | Command pattern                                        |
   | ------------ | ------------------------------------------------------ |
   | Bazel        | `bazel build //<package>/...` for each changed package |
   | npm/pnpm     | `npm run build` or `pnpm build`                        |
   | Cargo        | `cargo build`                                          |
   | Go           | `go build ./...`                                       |
   | Maven        | `mvn compile -pl <module>`                             |

   Use a **10 minute timeout** — builds can be slow.

3. **If the build fails**, diagnose and fix. See `references/troubleshooting.md` for common failures by build system. Choose the appropriate fix strategy:

   **Fix belongs to the current branch:**

   ```bash
   git add <fixed-files>
   gt modify
   ```

   **Fix belongs to a downstack branch** (e.g., a missing dep in branch 1 discovered while building branch 3) — use `gt modify --into` to amend the fix into the correct branch without switching:

   ```bash
   git add <fixed-files>
   gt modify --into <target-branch>
   ```

   **Fixes belong to multiple branches** (e.g., a missing dep in branch 1 and a visibility fix in branch 2) — use `gt absorb` to automatically distribute staged changes to the correct branches based on which commits they apply to:

   ```bash
   git add <fixed-files>
   gt absorb --dry-run  # verify correct routing first
   gt absorb
   ```

   All three approaches auto-restack upstack branches after amending.

4. **If restacking produces conflicts**, resolve them:

   ```bash
   # resolve conflicts in flagged files
   gt add <resolved-file>
   gt continue
   ```

5. **Re-verify the build** after any fix before moving to the next branch.

### Post-verification

After all branches build, verify the full stack:

```bash
gt log                                           # confirm stack structure
git log <trunk>..HEAD --format="%h %G? %s"       # confirm all commits signed (if GPG enabled)
```

---

## Step 6: Submit and Describe

Push the stack to Graphite and add PR descriptions.

### Submit

```bash
gt submit --stack --no-interactive
```

This creates draft PRs for all branches. If remote branches have diverged:

```bash
gt submit --stack --no-interactive --force
```

### PR Descriptions

Update each PR with a proper description using `gh pr edit`:

```bash
gh pr edit <branch-name> --body "$(cat <<'EOF'
## Summary

- <bullet points describing changes in this PR>

## Test plan

- [x] <completed verification items>
- [ ] <pending verification items>

---
**Stack: N/M** — <short description of this PR's role>

Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

For the excluded-files PR, note explicitly that it should not be merged:

```
## Summary

- Planning/tooling artifacts used during development
- **Not intended for merge** — close without merging or merge last and revert

## Test plan

- N/A — no code changes

---
**Optional PR** — Excluded artifacts (do not merge)
```

### Final output

Present the completed stack to the user:

```
| # | PR | Title |
|---|---|---|
| 1 | #NNNNN | TICKET-ID: ... |
| 2 | #NNNNN | TICKET-ID: ... |
| ... | ... | ... |
```

Include the Graphite link to the first PR (bottom of the stack) — this is the entry point for reviewers.
