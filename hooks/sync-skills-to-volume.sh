#!/usr/bin/env bash
# PostToolUse hook: sync an updated skill directory to /Volumes/mitch/.claude/skills
#
# Fires on Edit|Write. Filters to ~/.claude/skills/** paths only.
# Silently no-ops if the volume isn't mounted — never blocks or errors.
# Syncs the specific skill directory (not the whole tree) for speed.

set -uo pipefail

TARGET="/Volumes/mitch/.claude/skills"
SKILLS_DIR="$HOME/.claude/skills"

payload="$(cat 2>/dev/null || true)"
[ -z "$payload" ] && exit 0

file_path="$(printf '%s' "$payload" \
  | jq -r '.tool_input.file_path // .tool_response.filePath // empty' 2>/dev/null || true)"
[ -z "$file_path" ] && exit 0

# Only act on files under ~/.claude/skills/
case "$file_path" in
  "$SKILLS_DIR"/*) ;;
  *) exit 0 ;;
esac

# Volume must be mounted; skip silently if not
[ -d "$TARGET" ] || exit 0

# Derive the skill name — first path component under skills/
relative="${file_path#"$SKILLS_DIR"/}"
skill_name="${relative%%/*}"

if [ -n "$skill_name" ] && [ -d "$SKILLS_DIR/$skill_name" ]; then
  rsync -a --delete \
    "$SKILLS_DIR/$skill_name/" \
    "$TARGET/$skill_name/" \
    2>/dev/null || true
fi

exit 0
