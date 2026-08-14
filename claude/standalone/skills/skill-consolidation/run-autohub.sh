#!/usr/bin/env bash
# run-autohub.sh — durable weekly entrypoint for the skill-hub-autodetect scheduler.
#
# Wired into the user's crontab (see README "Schedule it"):
#   13 8 * * 1  $HOME/.claude/skill-consolidation/run-autohub.sh >> $HOME/.claude/skill-consolidation/autohub.log 2>&1
#
# Step 1 (always, read-only): run the detector. If there is no ready candidate, exit quietly —
# nothing is mutated. Only when a real >=8 unhubbed candidate exists do we hand off to Claude
# headless with the standing prompt to perform the consolidation end-to-end.
#
# NOTE: this launcher does NOT bypass permission gates. When it hands off to `claude -p`, the
# normal permission prompts apply (the consolidation mutates the skills tree). If you want the
# job to run fully unattended you must opt in explicitly via your own Claude settings allowlist —
# this script will not do it for you.
set -euo pipefail

# cron's PATH (/usr/bin:/bin) has neither homebrew node nor ~/.local/bin claude — fix lookup for
# ALL invocations below (two bare `node` calls + the `claude` consolidation handoff) in one line.
export PATH="/opt/homebrew/bin:$HOME/.local/bin:$PATH"
command -v node >/dev/null || { echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] FATAL: node not found on PATH"; exit 1; }

CONS_DIR="$HOME/.claude/skill-consolidation"
PROMPT_FILE="$CONS_DIR/auto-hub-agent-prompt.md"
DETECTOR="$CONS_DIR/detect-candidates.mjs"

cd "$CONS_DIR"

JSON="$(node "$DETECTOR" --json)"

# Count ready candidates without needing jq.
READY="$(printf '%s' "$JSON" | node -e 'let d="";process.stdin.on("data",c=>d+=c).on("end",()=>{try{const j=JSON.parse(d);process.stdout.write(String((j.readyCandidates||[]).length))}catch{process.stdout.write("0")}})')"

TS="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
if [ "$READY" = "0" ]; then
  echo "[$TS] skill-hub-autodetect: no ready hub candidate this week. Exiting quietly."
  exit 0
fi

echo "[$TS] skill-hub-autodetect: $READY ready candidate(s) detected — handing off to consolidation agent."
printf '%s\n' "$JSON"

# Hand off to Claude headless with the standing prompt. Permission gates remain in force.
if command -v claude >/dev/null 2>&1; then
  claude -p "Follow the standing prompt at $PROMPT_FILE exactly. A ready candidate was just detected. Perform Step 1b (judgment filter) through Step 3 (report) for the highest-count ready candidate. Detection is read-only; you are the only mutator and only when a real >=8 candidate survives the filter." \
    || echo "[$TS] claude headless invocation returned non-zero (likely awaiting permission). Run the standing prompt manually if needed: $PROMPT_FILE"
else
  echo "[$TS] 'claude' CLI not on PATH — cannot auto-consolidate. Run the standing prompt manually: $PROMPT_FILE"
fi
