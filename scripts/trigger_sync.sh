#!/usr/bin/env bash
# Trigger sync in background if not already running
set -e

REPO_DIR="/Users/mitch.hudson/dev/skills"
PYTHON="/usr/bin/python3"
SCRIPT="$REPO_DIR/scripts/sync_skills.py"

# Export GitHub token if available
if [[ -z "$GITHUB_PERSONAL_ACCESS_TOKEN" ]]; then
  export GITHUB_PERSONAL_ACCESS_TOKEN="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"
fi

# Run detached in background so caller is not delayed
nohup "$PYTHON" "$SCRIPT" > /dev/null 2>&1 &
