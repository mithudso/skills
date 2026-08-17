#!/bin/bash
#
# Automates SpiderMonkey fork preparation for a MozJS ESR upgrade.
# Covers Steps 1-5 of the MozJS Upgrade Guide:
#   1. Identify target ESR version
#   2. Clone mongodb-forks/spidermonkey
#   3. Fetch ESR sources from Mozilla
#   4. Cherry-pick base + MongoDB-specific commits
#   5. Build and test
#
# Usage:
#   ./upgrade.sh                                       # interactive prompts
#   ./upgrade.sh <target_major> <target_minor>         # auto-detect current version
#   ./upgrade.sh <cur_major> <cur_minor> <tgt_major> <tgt_minor>  # fully explicit

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MONGO_REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"

# ---------------------------------------------------------------------------
# Parse arguments or prompt interactively
# ---------------------------------------------------------------------------
if [ $# -eq 4 ]; then
    CURRENT_ESR_MAJOR="$1"
    CURRENT_ESR_MINOR="$2"
    TARGET_ESR_MAJOR="$3"
    TARGET_ESR_MINOR="$4"
elif [ $# -eq 2 ]; then
    TARGET_ESR_MAJOR="$1"
    TARGET_ESR_MINOR="$2"
    # Auto-detect current version from import.sh
    GET_SOURCES="${MONGO_REPO_ROOT}/src/third_party/mozjs/scripts/import.sh"
    if [ -f "$GET_SOURCES" ]; then
        CURRENT_VERSION=$(grep '^VERSION=' "$GET_SOURCES" | sed 's/VERSION="//' | sed 's/".*//')
        CURRENT_ESR_MAJOR=$(echo "$CURRENT_VERSION" | cut -d. -f1)
        CURRENT_ESR_MINOR=$(echo "$CURRENT_VERSION" | cut -d. -f2)
        echo "Auto-detected current ESR version: ${CURRENT_ESR_MAJOR}.${CURRENT_ESR_MINOR}"
    else
        echo "ERROR: Cannot auto-detect current version -- import.sh not found at $GET_SOURCES"
        exit 1
    fi
elif [ $# -eq 0 ]; then
    # Try auto-detect first
    GET_SOURCES="${MONGO_REPO_ROOT}/src/third_party/mozjs/scripts/import.sh"
    if [ -f "$GET_SOURCES" ]; then
        CURRENT_VERSION=$(grep '^VERSION=' "$GET_SOURCES" | sed 's/VERSION="//' | sed 's/".*//')
        CURRENT_ESR_MAJOR=$(echo "$CURRENT_VERSION" | cut -d. -f1)
        CURRENT_ESR_MINOR=$(echo "$CURRENT_VERSION" | cut -d. -f2)
        echo "Auto-detected current ESR version: ${CURRENT_ESR_MAJOR}.${CURRENT_ESR_MINOR}"
    else
        echo "What is the current ESR major version in the mongo repo?"
        read -r CURRENT_ESR_MAJOR
        echo "What is the current ESR minor version in the mongo repo?"
        read -r CURRENT_ESR_MINOR
    fi
    echo "What is the ESR major version you are targeting to upgrade to?"
    read -r TARGET_ESR_MAJOR
    echo "What is the ESR minor version you are targeting to upgrade to?"
    read -r TARGET_ESR_MINOR
else
    echo "Usage: $0 [<current_major> <current_minor>] <target_major> <target_minor>"
    exit 1
fi

echo "=== Upgrade: ESR ${CURRENT_ESR_MAJOR}.${CURRENT_ESR_MINOR} -> ESR ${TARGET_ESR_MAJOR}.${TARGET_ESR_MINOR} ==="

# ---------------------------------------------------------------------------
# Step 2: Clone the SpiderMonkey fork
# ---------------------------------------------------------------------------
echo "---------- CLONING THE REPO -----------------"
cd "$HOME"
if [ -d "spidermonkey" ]; then
    echo "spidermonkey/ already exists in $HOME -- reusing it"
    cd spidermonkey
else
    git clone git@github.com:mongodb-forks/spidermonkey.git
    cd spidermonkey
fi

# ---------------------------------------------------------------------------
# Step 3: Fetch ESR sources from Mozilla
# ---------------------------------------------------------------------------
echo "---------- FETCHING MOZILLA SOURCES ----------"

# NOTE (ESR 140.3+): Mozilla deprecated github.com/mozilla/gecko-dev and now
# uses github.com/mozilla-firefox/firefox. The repos contain the same code but
# different commit hashes.
if ! git remote get-url firefox &>/dev/null; then
    git remote add firefox git@github.com:mozilla-firefox/firefox.git
fi
# --depth 1: shallow fetch — we only need the tip commit at each tag, not full history.
# --force: without it, any rejected tag causes git to atomically roll back ALL fetched
# refs (including remote-tracking branches), so firefox/esr140 never gets written.
git fetch firefox --tags --depth 1 --force

FIREFOX_ESR_LOCAL_BRANCH="firefox_esr${TARGET_ESR_MAJOR}"
FIREFOX_ESR_REMOTE_BRANCH="firefox/esr${TARGET_ESR_MAJOR}"

if git show-ref --verify --quiet "refs/heads/${FIREFOX_ESR_LOCAL_BRANCH}"; then
    git checkout "${FIREFOX_ESR_LOCAL_BRANCH}"
    git pull
else
    git branch "${FIREFOX_ESR_LOCAL_BRANCH}" --track "${FIREFOX_ESR_REMOTE_BRANCH}"
    git checkout "${FIREFOX_ESR_LOCAL_BRANCH}"
fi

# ---------------------------------------------------------------------------
# Find the ESR release commit via git tags
# ---------------------------------------------------------------------------
echo "---------- FINDING ESR RELEASE TAG -----------"
TARGET_TAG=$(git tag -l "FIREFOX_${TARGET_ESR_MAJOR}_${TARGET_ESR_MINOR}_*esr_RELEASE" |
    sort -t_ -k4,4n -k5,5n |
    tail -n1)

if [ -z "$TARGET_TAG" ]; then
    echo "ERROR: No tag found matching FIREFOX_${TARGET_ESR_MAJOR}_${TARGET_ESR_MINOR}_*esr_RELEASE"
    echo "Available tags for ESR ${TARGET_ESR_MAJOR}:"
    git tag -l "FIREFOX_${TARGET_ESR_MAJOR}_*" | head -20
    exit 1
fi

TARGET_ESR_COMMIT=$(git rev-list -n1 "$TARGET_TAG")
echo "Found tag: ${TARGET_TAG}"
echo "Target commit: ${TARGET_ESR_COMMIT}"

# Reset to the ESR commit
git reset --hard "${TARGET_ESR_COMMIT}"

# ---------------------------------------------------------------------------
# Create branches (without LDAP prefix per current convention)
# ---------------------------------------------------------------------------
echo "---------- CREATING UPGRADE BRANCHES ---------"
ESR_BRANCH="esr${TARGET_ESR_MAJOR}.${TARGET_ESR_MINOR}"
SPIDERMONKEY_ONLY_BRANCH="spidermonkey-esr${TARGET_ESR_MAJOR}.${TARGET_ESR_MINOR}-cpp-only"

# Create the ESR tracking branch
if ! git show-ref --verify --quiet "refs/heads/${ESR_BRANCH}"; then
    git checkout -b "${ESR_BRANCH}"
    git push --set-upstream origin "${ESR_BRANCH}"
else
    echo "Branch ${ESR_BRANCH} already exists, skipping creation"
fi

# Create the working branch for SpiderMonkey extraction
git checkout -b "${SPIDERMONKEY_ONLY_BRANCH}"
git push --set-upstream origin "${SPIDERMONKEY_ONLY_BRANCH}"

# ---------------------------------------------------------------------------
# Step 4: Cherry-pick base commits from the current SpiderMonkey branch
# ---------------------------------------------------------------------------
echo "---------- FINDING BASE COMMITS --------------"
CURRENT_SPIDERMONKEY_REPO="spidermonkey-esr${CURRENT_ESR_MAJOR}.${CURRENT_ESR_MINOR}-cpp-only"
# Ensure the source branch exists locally (fetch from origin if needed)
if ! git show-ref --verify --quiet "refs/heads/${CURRENT_SPIDERMONKEY_REPO}"; then
    echo "Branch ${CURRENT_SPIDERMONKEY_REPO} not found locally -- fetching from origin..."
    git fetch origin "${CURRENT_SPIDERMONKEY_REPO}:${CURRENT_SPIDERMONKEY_REPO}"
fi
git switch "${CURRENT_SPIDERMONKEY_REPO}"

# Use --format="%h" --max-count=1 instead of | head -1 | cut to avoid SIGPIPE.
# bash 5.1 with set -euo pipefail: head -1 closes the pipe, causing SIGPIPE in the
# git log subshell inside $(), which propagates as a fatal error.
# Note: COMMIT_ONE grep is intentionally short — the actual message is
# "Removed all unrequired Non-SpiderMonkey Files"; the full phrase doesn't match.
COMMIT_ONE=$(git log --grep="Non-SpiderMonkey Files" -i --format="%h" --max-count=1)
COMMIT_TWO=$(git log --grep="Removed Rust from Spidermonkey Repo" -i --format="%h" --max-count=1)
COMMIT_THREE=$(git log --grep="Replaced Rust dependencies in SpiderMonkey with C++ implementation" -i --format="%h" --max-count=1)
COMMIT_FOUR=$(git log --grep="Add TestLatin1.cpp" -i --format="%h" --max-count=1)

# Validate we found all base commits
for var_name in COMMIT_ONE COMMIT_TWO COMMIT_THREE COMMIT_FOUR; do
    if [ -z "${!var_name}" ]; then
        echo "ERROR: Could not find base commit for ${var_name}"
        echo "Check that branch ${CURRENT_SPIDERMONKEY_REPO} has the expected commit messages."
        exit 1
    fi
done

echo "Base commits found:"
echo "  1 (Remove non-SM files): ${COMMIT_ONE}"
echo "  2 (Remove Rust):         ${COMMIT_TWO}"
echo "  3 (C++ replacement):     ${COMMIT_THREE}"
echo "  4 (TestLatin1):          ${COMMIT_FOUR}"

# Find the remaining MongoDB-specific commits (after commit 4)
LAST_COMMIT=$(git rev-parse HEAD)
# tail -1 reads all input (no SIGPIPE); git log without --reverse shows newest first,
# so tail -1 gives the oldest = first MongoDB commit after COMMIT_FOUR.
COMMIT_FIVE=$(git log --format="%h" "${COMMIT_FOUR}..HEAD" | tail -1)

if [ -n "$COMMIT_FIVE" ]; then
    echo "  5+ (MongoDB commits):    ${COMMIT_FOUR}..${LAST_COMMIT}"
    echo "  Total MongoDB commits:   $(git log --oneline "${COMMIT_FOUR}..HEAD" | wc -l | tr -d ' ')"
else
    echo "  No additional MongoDB-specific commits found after commit 4"
fi

# ---------------------------------------------------------------------------
# Cherry-pick the base commits onto the new branch
# ---------------------------------------------------------------------------
git switch "${SPIDERMONKEY_ONLY_BRANCH}"

echo "---------- CHERRY-PICKING BASE COMMITS -------"
echo "Cherry-picking commit 1 (remove non-SpiderMonkey files)..."
git cherry-pick "${COMMIT_ONE}" >my_output_file.txt 2>&1 || true

CHERRY_PICK_SCRIPT="${SCRIPT_DIR}/resolve_delete_non_spidermonkey_files_cherry_pick.sh"
bash "${CHERRY_PICK_SCRIPT}" my_output_file.txt deleted_files.txt
rm -f my_output_file.txt deleted_files.txt

echo "Cherry-picking commit 2 (remove Rust)..."
git cherry-pick "${COMMIT_TWO}"

echo "Cherry-picking commit 3 (C++ replacement)..."
git cherry-pick "${COMMIT_THREE}"

echo "Cherry-picking commit 4 (TestLatin1)..."
git cherry-pick "${COMMIT_FOUR}"

# ---------------------------------------------------------------------------
# Dependency check
# ---------------------------------------------------------------------------
echo ""
echo "=========================================="
echo "  Base commits cherry-picked successfully."
echo "=========================================="
echo ""
echo "Before building, ensure you have Rust and cbindgen installed:"
echo "  cargo install cbindgen"
echo ""
echo "Press ENTER to continue to the build step, or Ctrl-C to abort."
read -r

# ---------------------------------------------------------------------------
# Step 5: Build and test
# ---------------------------------------------------------------------------
build_and_test() {
    local build_label="$1"
    echo "---------- BUILD AND TEST (${build_label}) ----------"

    # Download bootstrap.py if not present
    if [ ! -f bootstrap.py ]; then
        curl -L https://hg.mozilla.org/mozilla-central/raw-file/default/python/mozboot/bin/bootstrap.py -O
    fi

    cd js/src
    rm -rf _build
    mkdir _build
    cd _build

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
    make -j"$(nproc)"

    # Run tests, filtering for unexpected failures.
    echo "Running make check..."
    make check 2>&1 |
        sed '/^test-[^(pass|known)]/!d' |
        sed '/.*check_vanilla_allocations.*/d' \
            >make_check_failures.txt || true

    echo "Running check-jstests..."
    EXCLUDE_FILE="${SCRIPT_DIR}/exclude.txt"
    make check-jstests JSTESTS_EXTRA_ARGS="-x=${EXCLUDE_FILE}" 2>&1 |
        sed '/^test-[^(pass|known)]/!d' \
            >jstests_failures.txt || true

    echo "Running jsapi-tests..."
    ./dist/bin/jsapi-tests 2>&1 |
        sed '/^test-[^(pass|known)]/!d' |
        sed '/.*DeflateStringToUTF8Buffer.*/d' \
            >jsapi_tests_failures.txt || true

    # Auto-check results
    local has_failures=0
    for log in make_check_failures.txt jstests_failures.txt jsapi_tests_failures.txt; do
        if [ -s "$log" ]; then
            echo "UNEXPECTED FAILURES in ${log}:"
            cat "$log"
            has_failures=1
        else
            echo "${log}: PASS (no unexpected failures)"
        fi
    done

    cd ../../../

    if [ "$has_failures" -eq 1 ]; then
        echo ""
        echo "WARNING: Unexpected test failures detected in ${build_label} build."
        echo "Review the failures above. Press ENTER to continue anyway, or Ctrl-C to abort."
        read -r
    else
        echo ""
        echo "${build_label} build: All tests passed."
    fi
}

# First build: verify base commits
build_and_test "AFTER BASE COMMITS"

# ---------------------------------------------------------------------------
# Cherry-pick remaining MongoDB-specific commits
# ---------------------------------------------------------------------------
if [ -n "$COMMIT_FIVE" ]; then
    echo "---------- CHERRY-PICKING REMAINING MONGODB COMMITS ----------"
    echo "Cherry-picking ${COMMIT_FOUR}..${LAST_COMMIT}"
    echo ""
    echo "TIP: If cherry-picks fail with conflicts, resolve them manually,"
    echo "     then run: git cherry-pick --continue"
    echo "     Compare each diff against the original branch to confirm correctness."
    echo "     If the source branch contains merge commits, use --skip to skip them"
    echo "     (their code changes are already in the preceding regular commits)."
    echo ""
    # COMMIT_FOUR..LAST_COMMIT: includes all commits AFTER COMMIT_FOUR up to HEAD.
    # Do NOT use COMMIT_FIVE..LAST_COMMIT — that range excludes COMMIT_FIVE itself.
    git cherry-pick "${COMMIT_FOUR}..${LAST_COMMIT}"

    # Second build: verify everything together
    build_and_test "AFTER ALL COMMITS"
else
    echo "No additional MongoDB commits to cherry-pick."
fi

# ---------------------------------------------------------------------------
# Done
# ---------------------------------------------------------------------------
echo ""
echo "=========================================="
echo "  SpiderMonkey fork preparation complete!"
echo "=========================================="
echo ""
echo "Branch: ${SPIDERMONKEY_ONLY_BRANCH}"
echo "Commit: $(git rev-parse HEAD)"
echo ""
echo "Next steps:"
echo "  1. Push the branch if not already pushed:"
echo "     git push origin ${SPIDERMONKEY_ONLY_BRANCH}"
echo "  2. Proceed to integration into the mongo repo."
echo "     Update src/third_party/mozjs/scripts/import.sh with:"
echo "       VERSION=\"${TARGET_ESR_MAJOR}.${TARGET_ESR_MINOR}.0esr\""
echo "       LIB_GIT_BRANCH=${SPIDERMONKEY_ONLY_BRANCH}"
echo "       LIB_GIT_REVISION=$(git rev-parse HEAD)"
