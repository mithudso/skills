#!/usr/bin/env bash
# PostToolUse hook: when a Write tool call CREATES a file, classify it by type
# and inject a system reminder telling Claude to run the matching deep optimizer
# on it, iterating until no medium-or-higher findings remain.
#
#   documents  (.md/.txt/.pdf/.docx/...)   -> /ddo  (Document Deep Optimizer)
#   scripts/apps (.js/.py/.sh/.go/.rs/...)  -> /cdo  (Code Deep Optimizer)
#   images     (.png/.svg/.jpg/...)         -> /deso (Design Deep Optimizer)
#   prompts    (prompts/ dir, *.prompt*)    -> /pdo  (Prompt Deep Optimizer)
#
# Creation semantics: this fires on the Write tool only (new file / overwrite),
# NOT on Edit/MultiEdit. That is deliberate — it (a) matches the "every time X is
# created" intent and (b) avoids a re-trigger loop, since the optimizers fix
# files in place via Edit, which this hook ignores.
#
# Loop caps / exit conditions are NOT defined here. They live canonically in
# ~/.claude/skill-consolidation/convergence-and-severity.md; this hook cites that
# file rather than inlining numbers, matching the sibling
# skill-change-trigger-optimizers.sh hook. SKILL.md files are skipped — that
# sibling hook already owns them.
#
# Triggered from ~/.claude/settings.json hooks.PostToolUse with matcher "Edit|Write".

set -u

# Read the hook payload from stdin once.
payload="$(cat 2>/dev/null || true)"
[ -z "$payload" ] && exit 0

# Creation only: ignore in-place edits so optimizer fixes never re-trigger.
tool_name="$(printf '%s' "$payload" | jq -r '.tool_name // empty' 2>/dev/null || true)"
case "$tool_name" in
  Edit|MultiEdit) exit 0 ;;
esac

# Extract the affected file path.
file_path="$(printf '%s' "$payload" | jq -r '.tool_input.file_path // .tool_response.filePath // empty' 2>/dev/null || true)"
[ -z "$file_path" ] && exit 0

lc_path="$(printf '%s' "$file_path" | tr '[:upper:]' '[:lower:]')"
ext="${lc_path##*.}"

# Precedence 1: SKILL.md is owned by skill-change-trigger-optimizers.sh.
case "$file_path" in
  */SKILL.md|*/skill.md|*-skill.md) exit 0 ;;
esac

cmd=""; kind=""; skill=""

# Precedence 2: prompt files (no standard extension -> path/name heuristic).
case "$lc_path" in
  *.prompt|*.prompt.md|*.prompt.txt|*/prompts/*|*/.prompts/*)
    cmd="/pdo"; kind="prompt"; skill="prompt-deep-optimizer" ;;
esac

# Precedence 3: classify by extension.
if [ -z "$cmd" ]; then
  case "$ext" in
    md|markdown|txt|text|rst|adoc|asciidoc|org|rtf|pdf|docx|odt|tex)
      cmd="/ddo"; kind="document"; skill="Document Deep Optimizer" ;;
    js|mjs|cjs|jsx|ts|tsx|py|sh|bash|zsh|go|rs|rb|php|java|kt|kts|scala|swift|c|h|cpp|cc|cxx|hpp|cs|lua|pl|pm|r|sql|vue|svelte)
      cmd="/cdo"; kind="script or application"; skill="code-deep-optimizer" ;;
    png|jpg|jpeg|gif|webp|svg|bmp|tif|tiff|avif|ico)
      cmd="/deso"; kind="image"; skill="design-deep-optimizer" ;;
  esac
fi

# Unclassified -> nothing to do.
[ -z "$cmd" ] && exit 0

# Debounce: skip if this exact path was flagged in the last 60s.
mark_dir="${TMPDIR:-/tmp}/claude-auto-optimize"
mkdir -p "$mark_dir" 2>/dev/null || true
mark="$mark_dir/$(printf '%s' "$file_path" | tr '/ ' '__')"
now="$(date +%s 2>/dev/null || echo 0)"
if [ -f "$mark" ]; then
  prev="$(cat "$mark" 2>/dev/null || echo 0)"
  case "$prev" in ''|*[!0-9]*) prev=0 ;; esac
  [ $((now - prev)) -lt 60 ] && exit 0
fi
printf '%s' "$now" > "$mark" 2>/dev/null || true

contract="${HOME}/.claude/skill-consolidation/convergence-and-severity.md"

# Emit the additionalContext payload back to Claude.
jq -n --arg f "$file_path" --arg c "$cmd" --arg k "$kind" --arg s "$skill" --arg contract "$contract" '{
  hookSpecificOutput: {
    hookEventName: "PostToolUse",
    additionalContext: (
      "New " + $k + " created: " + $f + "\n\n" +
      "Per the standing auto-optimizer instruction, run the " + $c + " (" + $s + ") on this file, iterating until no medium-or-higher findings remain. Iteration caps and exit conditions are defined canonically in " + $contract + " — cite that file; do not re-derive caps or instability rules here.\n\n" +
      "Procedure:\n" +
      "1. Run `" + $c + " " + $f + "` (or invoke the " + $s + " skill via the Skill tool with that path).\n" +
      "2. Apply every Medium-or-higher fix in place.\n" +
      "3. Re-run and stop when it reports no medium-or-higher findings, or on any canonical exit condition (iteration cap, no progress, instability) per " + $contract + ".\n" +
      "4. Report the per-iteration finding-severity distribution in your response.\n\n" +
      "If this file is an optimizer-generated report or scratch artifact, skip it — do not optimize your own output."
    )
  }
}' 2>/dev/null || true
