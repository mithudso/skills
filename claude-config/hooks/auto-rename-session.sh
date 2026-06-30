#!/bin/bash
# auto-rename-session.sh
#
# SessionStart hook: auto-names an UNNAMED session from the first human prompt
# in its transcript, with the same effect as /rename.
#
# Why SessionStart (and not UserPromptSubmit)?
#   Only the SessionStart event can emit hookSpecificOutput.sessionTitle to rename
#   a session (verified against the Claude Code hooks reference). UserPromptSubmit
#   carries the prompt text but CANNOT set the title. SessionStart carries no prompt
#   text but DOES carry transcript_path — so we read the first prompt from there.
#
# "New task, not follow-ups":
#   - Follow-up prompts fire UserPromptSubmit, which this hook does not touch.
#   - We only act when session_title is empty, so once a session is named (by us,
#     by `claude -n`, or by `/rename`) later SessionStarts no-op. No re-renaming.
#   - The description is always taken from the FIRST prompt = the originating task.
#
# Caveat (a real Claude Code limitation, not a bug):
#   sessionTitle only applies on source "startup"/"resume". On a brand-new session
#   (source=startup) the transcript has no prompt yet, so the rename lands on the
#   next start of that session (e.g. `claude --continue`/`--resume`). An unnamed
#   session is already shown by its first prompt in the picker until then.

set -uo pipefail

JQ=/usr/bin/jq
command -v "$JQ" >/dev/null 2>&1 || JQ=jq
command -v "$JQ" >/dev/null 2>&1 || exit 0   # no jq -> do nothing, never block

input=$(cat)

source=$(printf '%s' "$input"        | "$JQ" -r '.source // ""')
existing=$(printf '%s' "$input"      | "$JQ" -r '.session_title // ""')
transcript=$(printf '%s' "$input"    | "$JQ" -r '.transcript_path // ""')

# sessionTitle is honored only on startup/resume; ignored on clear/compact.
case "$source" in
  startup|resume) ;;
  *) exit 0 ;;
esac

# Respect any title the user already set (and stop re-firing on later resumes).
[ -n "$existing" ] && exit 0

# No transcript yet (fresh startup) -> nothing to describe; rename on next start.
[ -z "$transcript" ] || [ ! -s "$transcript" ] && exit 0

# Pull the first genuine human prompt out of the transcript JSONL.
#  - user-role entries only, skipping injected/meta entries
#  - content may be a string or an array of blocks; take text blocks
#  - skip slash-command expansions and the local-command caveat banner
raw=$(
  "$JQ" -rn --slurpfile _ /dev/null '
    first(
      inputs
      | select(.type == "user")
      | select((.isMeta // false) | not)
      | (.message.content) as $c
      | (
          if   ($c | type) == "string" then $c
          elif ($c | type) == "array"  then ([ $c[]? | select(.type == "text") | .text ] | join(" "))
          else "" end
        )
      | select(test("\\S"))
      | select(test("^[[:space:]]*<command-")        | not)
      | select(test("^[[:space:]]*<local-command")   | not)
      | select(test("^[[:space:]]*Caveat:")          | not)
    ) // ""
  ' < "$transcript" 2>/dev/null
)

[ -z "$raw" ] && exit 0

# Collapse ALL whitespace (newlines included) to single spaces, then trim.
title=$(printf '%s' "$raw" | tr '\n\r\t' '   ' | sed -E 's/[[:space:]]+/ /g; s/^[[:space:]]+//; s/[[:space:]]+$//')
# Strip a leading "/slash-command" token, then leading markdown/list/quote markers.
title=$(printf '%s' "$title" | sed -E 's#^/[A-Za-z0-9_-]+[[:space:]]+##; s/^[*_#>`~ -]+//; s/^[[:space:]]+//')

[ -z "$title" ] && exit 0

# Truncate to a short slug on a word boundary.
MAX=50
if [ "${#title}" -gt "$MAX" ]; then
  title=${title:0:$MAX}
  title=${title% *}
  title="${title}…"
fi

# Emit via jq so the title is correctly JSON-escaped.
"$JQ" -cn --arg t "$title" \
  '{hookSpecificOutput: {hookEventName: "SessionStart", sessionTitle: $t}}'
exit 0
