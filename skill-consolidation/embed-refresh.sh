#!/bin/sh
# embed-refresh.sh — keep SKILLS-INDEX.* and SKILLS-EMBEDDINGS.json in sync with the skill
# library. Invoked by the com.mitch.skills-embed launchd agent: on changes to ~/.claude/skills
# (a skill added/removed/renamed) and every 6h (catches edits inside an existing skill).
# Fail-open: a down embedding server never fails the index refresh, and the agent always exits 0.
set -u
DIR="$HOME/.claude/skill-consolidation"
LOG="$DIR/embed-refresh.log"
cd "$DIR" 2>/dev/null || exit 0

# Resolve a node binary (launchd starts with a minimal PATH; survive Homebrew version bumps).
NODE=""
for c in /opt/homebrew/bin/node /usr/local/bin/node "$(command -v node 2>/dev/null)"; do
  if [ -n "$c" ] && [ -x "$c" ]; then NODE="$c"; break; fi
done

ts() { date "+%Y-%m-%dT%H:%M:%S"; }
{
  echo "[$(ts)] refresh start (node=${NODE:-NONE})"
  if [ -z "$NODE" ]; then
    echo "[$(ts)] no node binary found — abort (fail-open)"
    echo "[$(ts)] refresh done"
    exit 0
  fi
  # 1) Keep the regex index fresh (offline, deterministic — never touches the network).
  "$NODE" gen-skills-index.mjs --quiet || echo "[$(ts)] index regen non-zero"
  # 2) Incrementally (re)embed. --embed is no-churn: it skips the write when nothing changed,
  #    and exits non-zero if the embedding server is unreachable — caught here, never fatal.
  if "$NODE" gen-skills-index.mjs --embed; then :; else
    echo "[$(ts)] embed non-zero — embedding server down? (fail-open, index still fresh)"
  fi
  echo "[$(ts)] refresh done"
} >> "$LOG" 2>&1
exit 0
