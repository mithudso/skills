# mozjs-upgrade-fork

**Category:** Databases, Data Engineering & Analytics
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/query-integration/mozjs-upgrade-fork/skills/mozjs-upgrade-fork

## Description
Use when upgrading MozJS/SpiderMonkey to a new ESR version. Covers fork preparation: cloning the spidermonkey fork, fetching from Mozilla, cherry-picking base and MongoDB-specific commits, and building/testing. Use when asked to "upgrade mozjs", "upgrade spidermonkey", "prepare spidermonkey fork", or "begin ESR upgrade".

---

# MozJS Upgrade -- Fork Preparation (Steps 1-5)

This skill guides you through preparing the SpiderMonkey fork for a MozJS ESR upgrade. It covers
cloning the fork, fetching the new ESR sources from Mozilla, cherry-picking the base and
MongoDB-specific commits, and building/testing the result.

## Prerequisites

Before starting, ensure the following are in place:

1. **MANA access** to push to `mongodb-forks`:
   `mana.corp.mongodbgov.com/resources/61fae35b2f8fb870b9a97c1d`
2. **SSH key authorized** for the `mongodb-forks` org: go to GitHub **Settings > SSH and GPG keys >
   Configure SSO** and authorize your key for `mongodb-forks`.
3. **Sufficient disk space** -- `upgrade.sh` uses a shallow fetch (`--depth 1`) which only pulls
   the tip commit at each tag. Ensure at least 5 GB free. The fetch typically takes 1-3 minutes.

## Quick Start -- Using upgrade.sh

The preferred path is to run the automation script located at
`.agents/skills/mozjs-upgrade-fork/references/upgrade.sh`:

```bash
# Run from the mongo repo root:
UPGRADE=.agents/skills/mozjs-upgrade-fork/references/upgrade.sh

# Auto-detect current version from import.sh, specify target:
bash "$UPGRADE" <target_major> <target_minor>

# Fully explicit (current and target):
bash "$UPGRADE" <cur_major> <cur_minor> <tgt_major> <tgt_minor>

# Interactive prompts:
bash "$UPGRADE"
```

The script automates the entire fork preparation:

1. **Cloning** `mongodb-forks/spidermonkey` into `~/spidermonkey` (or reuses an existing clone).
2. **Fetching** from `mozilla-firefox/firefox` (the canonical Mozilla remote) with tags.
3. **Finding the ESR release tag** via `git tag -l "FIREFOX_{MAJOR}_{MINOR}_*esr_RELEASE"`.
4. **Creating branches:**
   - `esr{MAJOR}.{MINOR}` -- tracks the ESR release commit
   - `spidermonkey-esr{MAJOR}.{MINOR}-cpp-only` -- the working branch for extraction
5. **Cherry-picking the 4 base commits** (file removal, Rust removal, C++ replacement, TestLatin1),
   using `resolve_delete_non_spidermonkey_files_cherry_pick.sh` for the first commit.
6. **Building and testing** after the base commits to verify the baseline.
7. **Cherry-picking remaining MongoDB-specific commits** from the current branch.
8. **Final build and test** to verify everything together.

**Important:** The script has interactive pauses (`read` prompts) after cherry-picks and builds. It
cannot be run in the background. The first cherry-pick (removing non-SpiderMonkey files) is
especially slow (10-30+ minutes) due to the massive number of file deletions and conflict
resolution. Run the script in a foreground terminal and expect the full process to take 1-2+ hours.

If the script completes successfully, skip to the [Verification Gate](#verification-gate) section.

## Manual Steps (When the Script Fails)

If `upgrade.sh` encounters issues, the individual steps can be run manually.

### Step 1: Identify the Current and Target Versions

The current ESR version is recorded in `import.sh`:

```bash
grep 'VERSION=' src/third_party/mozjs/import.sh
# Output example: VERSION="140.9.0esr"
```

The target ESR version should be the latest release from Mozilla. Check the Firefox ESR
release channel at `support.mozilla.org/en-US/kb/choosing-firefox-update-channel`.

### Step 2: Clone the SpiderMonkey Fork

```bash
cd ~
git clone git@github.com:mongodb-forks/spidermonkey.git
cd spidermonkey
```

### Step 3: Fetch ESR Sources from Mozilla

> **Note (ESR 140.3+):** Mozilla deprecated `github.com/mozilla/gecko-dev` and now uses
> `github.com/mozilla-firefox/firefox`. The repos contain the same code but different commit hashes.

```bash
git remote add firefox git@github.com:mozilla-firefox/firefox.git
git fetch firefox --tags --depth 1 --force
```

Create a local branch tracking the ESR major release:

```bash
git branch firefox_esr${MAJOR} --track firefox/esr${MAJOR}
git checkout firefox_esr${MAJOR}
```

Find the ESR release commit via tags:

```bash
MAJOR=<target_major>
MINOR=<target_minor>

TARGET_TAG=$(git tag -l "FIREFOX_${MAJOR}_${MINOR}_*esr_RELEASE" \
  | sort -t_ -k4,4n -k5,5n \
  | tail -n1)

TARGET_COMMIT=$(echo "$TARGET_TAG" | xargs git rev-list -n1)
echo "Tag: ${TARGET_TAG}  Commit: ${TARGET_COMMIT}"
```

Reset to the release commit and create branches:

```bash
git reset --hard "${TARGET_COMMIT}"

# ESR tracking branch (no username prefix per convention)
git checkout -b "esr${MAJOR}.${MINOR}"
git push --set-upstream origin "esr${MAJOR}.${MINOR}"

# Working branch for SpiderMonkey extraction
git checkout -b "spidermonkey-esr${MAJOR}.${MINOR}-cpp-only"
git push --set-upstream origin "spidermonkey-esr${MAJOR}.${MINOR}-cpp-only"
```

### Step 4: Cherry-Pick Base and MongoDB-Specific Commits

#### 4a. Enumerate ALL commits on the current branch before touching anything

**CRITICAL: Before cherry-picking a single commit, enumerate the complete list of MongoDB-specific
commits from the current branch.** This list is your checklist. Every commit on it MUST appear in
the new branch's `git log` after the 4 base commits. No commit may be silently skipped — skipping
any commit without explicit justification has caused build failures on s390x/ppc64le (GCC) and
required extensive retrospective analysis in past upgrades.

```bash
CURRENT_BRANCH="spidermonkey-esr<CUR_MAJOR>.<CUR_MINOR>-cpp-only"
git switch "${CURRENT_BRANCH}"

# Identify the 4 base commits
COMMIT_ONE=$(git log --grep="Removed all Non-SpiderMonkey Files" -i --oneline | head -1 | cut -d' ' -f1)
COMMIT_TWO=$(git log --grep="Removed Rust from Spidermonkey Repo" -i --oneline | head -1 | cut -d' ' -f1)
COMMIT_THREE=$(git log --grep="Replaced Rust dependencies in SpiderMonkey with C++ implementation" -i --oneline | head -1 | cut -d' ' -f1)
COMMIT_FOUR=$(git log --grep="Add TestLatin1.cpp" -i --oneline | head -1 | cut -d' ' -f1)

# Print ALL MongoDB-specific commits (everything after COMMIT_FOUR)
echo "=== MongoDB-specific commits to cherry-pick (in order, oldest first) ==="
git log --oneline --reverse "${COMMIT_FOUR}..${CURRENT_BRANCH}"
echo ""
echo "Total: $(git log --oneline "${COMMIT_FOUR}..${CURRENT_BRANCH}" | wc -l) commits"
```

Save this list. You will verify against it after cherry-picking all commits.

#### 4b. Cherry-pick the 4 base commits

Switch to the new working branch and cherry-pick:

> **WARNING: The first cherry-pick (file removal) is extremely long-running.** It deletes thousands
> of files from the full Firefox tree, generating massive conflict output. It can take 10-30+
> minutes depending on the machine. **Do not run this in the background** — the resolver script and
> `git cherry-pick --continue` require a foreground terminal. Run each cherry-pick command
> individually and wait for completion before proceeding to the next.

```bash
git switch "spidermonkey-esr${MAJOR}.${MINOR}-cpp-only"

# Commit 1: Remove non-SpiderMonkey files (uses resolver script)
# THIS STEP IS VERY SLOW — can take 10-30+ minutes. Do NOT run in background.
git cherry-pick "${COMMIT_ONE}" > my_output_file.txt 2>&1 || true
bash .agents/skills/mozjs-upgrade-fork/references/resolve_delete_non_spidermonkey_files_cherry_pick.sh my_output_file.txt deleted_files.txt
rm -f my_output_file.txt deleted_files.txt

# Commit 2: Remove Rust
git cherry-pick "${COMMIT_TWO}"

# Commit 3: C++ replacement for Rust dependencies
git cherry-pick "${COMMIT_THREE}"

# Commit 4: Add TestLatin1.cpp
git cherry-pick "${COMMIT_FOUR}"
```

> **Stop here and build/test** (Step 5) to verify the baseline before continuing.

#### 4c. Cherry-pick ALL remaining MongoDB-specific commits — one by one

After the baseline passes, cherry-pick every MongoDB-specific commit individually. **Do not use a
range (`A..B`) and let git cherry-pick in bulk.** Cherry-pick each commit individually so you can
review the result before moving to the next.

```bash
COMMIT_FIVE=<first_commit_after_COMMIT_FOUR>
LAST_COMMIT=$(git rev-list "${CURRENT_BRANCH}" | head -1)
```

For each commit in the list you saved in step 4a, run:

```bash
git cherry-pick <SHA>
# Then compare the result against the original:
git diff HEAD~1 HEAD
git show <SHA_ON_OLD_BRANCH>
```

If a cherry-pick produces no diff (the upstream absorbed the change), do **not** silently skip it.
Instead, document it explicitly:

```bash
# Option A: The upstream change made this commit obsolete — skip with a note in the PR description.
#   Confirm by: git show <SHA> | grep MONGODB MODIFICATION
#   If no such annotation exists, the commit may be a pure upstream cherry-pick; skip is safe.
#   If an annotation exists, you must explain WHY the modification is no longer needed.

# Option B: The commit is still needed but conflicts — resolve and re-commit.
git cherry-pick --continue
```

#### 4d. Verify completeness after all cherry-picks

```bash
NEW_COMMIT_FOUR=$(git log --grep="Add TestLatin1.cpp" -i --oneline | head -1 | cut -d' ' -f1)

echo "=== Commits on new branch after base commits ==="
git log --oneline --reverse "${NEW_COMMIT_FOUR}..HEAD"
echo ""
echo "Total: $(git log --oneline "${NEW_COMMIT_FOUR}..HEAD" | wc -l) commits"
```

The total MUST match or exceed the count from step 4a. Any discrepancy means a commit was dropped
and must be tracked down and either re-applied or explicitly documented as intentionally skipped
with a justification in the PR description.

> **Lesson from 140.11 upgrade**: During the 140.9→140.11 cherry-pick, approximately 17
> MongoDB-specific commits were silently dropped. This caused GCC compilation failures on s390x and
> ppc64le that were only discovered in CI. Every commit must be accounted for.

## MongoDB Modification Annotation Rule

**This rule is IMPERATIVE and applies to every custom commit on the fork — not just conflict
resolutions.**

Any time you author a new commit that is not a direct cherry-pick from the Firefox repo (i.e., you
are writing new or modified code yourself), you MUST add an annotation comment immediately above
each changed line or block. This allows future upgraders to identify which changes are
MongoDB-specific and why they exist.

The comment style depends on the file type:

- C/C++ files: `// MONGODB MODIFICATION: <explanation>`
- Python/moz.build files: `# MONGODB MODIFICATION: <explanation>`
- Plain-text config files (e.g. `python/sites/build.txt`): `# MONGODB MODIFICATION: <explanation>`

**What counts as a custom commit:**

- Any fix to make the build work (e.g. adding libraries to `USE_LIBS`, restoring missing config)
- Any workaround for a build system issue specific to the MongoDB fork approach
- Any code you write to resolve a conflict in a non-obvious way

**What does NOT count (no annotation needed):**

- Cherry-picks from the Firefox repo (upstream commits, even if they conflict)
- Files restored verbatim from the Firefox ESR base (e.g.
  `git show <base-commit>:path/to/file > path/to/file`)

## Cherry-Pick Conflicts

- The first cherry-pick (file removal) is handled automatically by
  `.agents/skills/mozjs-upgrade-fork/references/resolve_delete_non_spidermonkey_files_cherry_pick.sh`.
- Remaining cherry-picks may have conflicts. Key guidance:
  - Always compare the visual diff against the original branch: `git diff HEAD~1 HEAD` vs. the same
    commit on the old branch.
  - All MongoDB modifications must be annotated with a comment immediately above the change. The
    comment style depends on the file type:
    - C/C++ files: `// MONGODB MODIFICATION: Justification for change.`
    - Python/moz.build files: `# MONGODB MODIFICATION: Justification for change.`
    - Plain-text config files (e.g. `python/sites/build.txt`):
      `# MONGODB MODIFICATION: Justification for change.`
  - If conflicts are non-trivial (semantic changes, API differences), stop and seek help.
- See `.agents/skills/mozjs-upgrade-fork/references/cherry-pick-troubleshooting.md` for detailed
  guidance on conflict resolution.

## Build and Test

After cherry-picking (both after base commits and after all commits), build and run the test suites.

Quick reference:

```bash
cd js/src
rm -rf _build && mkdir _build && cd _build

../configure \
    --disable-jemalloc \
    --disable-jit \
    --disable-wasm-moz-intgemm \
    --with-system-icu \
    --with-system-zlib \
    --without-intl-api \
    --enable-optimize \
    --enable-tests \
    --disable-bootstrap

make -j$(nproc)

# Test suite 1: Main tests
make check

# Test suite 2: JS tests with exclusions
make check-jstests JSTESTS_EXTRA_ARGS=-x=$(realpath <mongo-repo-root>/.agents/skills/mozjs-upgrade-fork/references/exclude.txt)

# Test suite 3: JSAPI tests
./dist/bin/jsapi-tests

# Test suite 4: JIT tests (many failures expected -- JIT is disabled)
make check-jit-test
```

See `.agents/skills/mozjs-upgrade-fork/references/build-and-test.md` for full details on
dependencies, configure flags, and how to interpret test results.

## Verification Gate

Before proceeding to integration (`mozjs-upgrade-integrate`), confirm:

1. **All three test suites pass** (modulo known exclusions in
   `.agents/skills/mozjs-upgrade-fork/references/exclude.txt`):
   - `make check` -- `check_vanilla_allocations.py` failures are only in non-MongoDB-modified code
   - `make check-jstests` -- no unexpected failures beyond the exclude list
   - `./dist/bin/jsapi-tests` -- only `test_DeflateStringToUTF8Buffer` fails (SERVER-99489)
2. **The branch is pushed to origin:**
   ```bash
   git push origin "spidermonkey-esr${MAJOR}.${MINOR}-cpp-only"
   ```

## Repo-Specific Context

Before starting, check for a repo-specific context file in `references/`. Use
[references/_template.md](references/_template.md) as a template if none exists for your repo.

## Output

After completion, record the branch name and commit hash. These are needed for
`mozjs-upgrade-integrate` to pull the sources into the mongo repo:

```
VERSION="{MAJOR}.{MINOR}.0esr"
LIB_GIT_BRANCH=spidermonkey-esr{MAJOR}.{MINOR}-cpp-only
LIB_GIT_REVISION=<commit_hash>
```

Update `src/third_party/mozjs/import.sh` with these values in the integration step.