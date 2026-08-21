#!/bin/bash
#
# Generates the "git hack branch" workaround for oversized Evergreen patches.
#
# When a MozJS upgrade patch is too large for Evergreen (even with --large),
# this script creates a minimal branch containing only Evergreen YAML changes
# that, at runtime, check out the actual MozJS branch with the full changes.
#
# Usage: ./gen-hack-branch.sh <mozjs_branch_name>
#
# Prerequisites:
#   - Run from the mongo repo root (or mozjsupgrade/ directory)
#   - The MozJS branch must already be pushed to origin
#   - The current branch should be the MozJS branch

set -euo pipefail

if [ $# -ne 1 ]; then
    echo "Usage: $0 <mozjs_branch_name>"
    echo "  mozjs_branch_name: the branch containing the full MozJS changes"
    exit 1
fi

BRANCH_NAME="$1"
REPO_ROOT=$(git rev-parse --show-toplevel)
DEFINITIONS_FILE="${REPO_ROOT}/etc/evergreen_yml_components/definitions.yml"
HACK_SCRIPT="${REPO_ROOT}/evergreen/functions/git_hack_branch.sh"

# ---------------------------------------------------------------------------
# Validate
# ---------------------------------------------------------------------------
if [ ! -f "$DEFINITIONS_FILE" ]; then
    echo "ERROR: definitions.yml not found at ${DEFINITIONS_FILE}"
    echo "Run this script from within the mongo repo."
    exit 1
fi

# Check that the branch exists on the remote
if ! git ls-remote --heads origin "$BRANCH_NAME" | grep -q "$BRANCH_NAME"; then
    echo "WARNING: Branch '${BRANCH_NAME}' not found on origin."
    echo "Make sure to push it before submitting the Evergreen patch."
fi

# ---------------------------------------------------------------------------
# Step 1: Create the git_hack_branch.sh script
# ---------------------------------------------------------------------------
echo "Creating ${HACK_SCRIPT}..."
mkdir -p "$(dirname "$HACK_SCRIPT")"
cat >"$HACK_SCRIPT" <<'HACK_EOF'
cd src
set -o errexit
set -o verbose
git checkout "$branch_name"
HACK_EOF

# ---------------------------------------------------------------------------
# Step 2: Add the hack branch definition to definitions.yml
# ---------------------------------------------------------------------------
echo "Adding git hack branch definition to definitions.yml..."

# Check if the definition already exists
if grep -q 'git hack branch' "$DEFINITIONS_FILE"; then
    echo "WARNING: 'git hack branch' definition already exists in definitions.yml."
    echo "Skipping definition insertion."
else
    # Insert the definition after the "git get shallow streams project" block.
    # We insert it before the "git get project no modules" function since it's
    # a logical place among the git-related definitions.
    HACK_DEFINITION=$(
        cat <<DEFEOF

  "git hack branch": &git_hack_branch
    command: subprocess.exec
    display_name: "get the branch we actually want to test"
    params:
      binary: bash
      args:
        - "src/evergreen/functions/git_hack_branch.sh"
      env:
        branch_name: ${BRANCH_NAME}
DEFEOF
    )
    # Insert after the first occurrence of "git get project no modules"'s closing line
    # We use a different approach: insert before the "add git tag" definition
    # which comes after the git get definitions.
    sed -i "/^  \"add git tag\": &add_git_tag$/i\\
\\
  \"git hack branch\": \\&git_hack_branch\\
    command: subprocess.exec\\
    display_name: \"get the branch we actually want to test\"\\
    params:\\
      binary: bash\\
      args:\\
        - \"src/evergreen/functions/git_hack_branch.sh\"\\
      env:\\
        branch_name: ${BRANCH_NAME}" "$DEFINITIONS_FILE"
fi

# ---------------------------------------------------------------------------
# Step 3: Insert *git_hack_branch after every *git_get_shallow_project
# ---------------------------------------------------------------------------
echo "Inserting *git_hack_branch after every *git_get_shallow_project reference..."

# Count occurrences before
BEFORE_COUNT=$(grep -c 'git_get_shallow_project' "$DEFINITIONS_FILE" || true)

# Insert *git_hack_branch after each line containing *git_get_shallow_project
# (but not the definition line itself)
sed -i '/^\s*- \*git_get_shallow_project$/a\    - *git_hack_branch' "$DEFINITIONS_FILE"

# Also handle the inline git.get_project in "git get project no modules"
# which uses a direct command instead of the anchor
sed -i '/^\s*- \*restore_git_history_and_tags$/{
    # Only add after restore_git_history_and_tags in "git get project no modules"
    # This is a secondary insertion point
}' "$DEFINITIONS_FILE"

AFTER_COUNT=$(grep -c 'git_hack_branch' "$DEFINITIONS_FILE" || true)
echo "Inserted *git_hack_branch in $((AFTER_COUNT - 1)) locations (plus 1 definition)."

# ---------------------------------------------------------------------------
# Step 4: Report
# ---------------------------------------------------------------------------
echo ""
echo "=========================================="
echo "  Git hack branch setup complete!"
echo "=========================================="
echo ""
echo "Files modified:"
echo "  - ${DEFINITIONS_FILE}"
echo "  - ${HACK_SCRIPT} (created)"
echo ""
echo "Next steps:"
echo "  1. Create a new branch from the base commit (BEFORE MozJS changes):"
echo "     git stash"
echo "     git checkout -b mozjs-evergreen-hack <base_commit>"
echo "  2. Apply the Evergreen changes:"
echo "     git checkout stash -- etc/evergreen_yml_components/definitions.yml evergreen/functions/git_hack_branch.sh"
echo "     git stash drop"
echo "     git add etc/evergreen_yml_components/definitions.yml evergreen/functions/git_hack_branch.sh"
echo "     git commit -m 'Temporary: git hack branch for MozJS upgrade testing'"
echo "  3. Submit the small patch:"
echo "     evergreen patch -p mongodb-mongo-master -a required -f -y -u"
echo ""
echo "IMPORTANT: After testing, delete the hack branch and ensure these"
echo "Evergreen changes are NOT included in the final MozJS PR."
