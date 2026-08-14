#!/bin/bash
#
# run-bf-triage.sh
# ----------------
# Unattended wrapper: batch-triage the "Workload Resilience" team's active BFs
# by invoking Claude Code headlessly with the bf-triage skill (Mode B).
#
# Design:
#   1. PRE-FLIGHT GUARD — verify the devprod-mcp-proxy token is present, valid,
#      and fresh enough. If the session looks dead, NOTIFY and EXIT instead of
#      letting Claude hang forever on a browser login.
#   2. HARD TIMEOUT — run Claude under a watchdog so an auth hang (or stuck run)
#      can't block indefinitely; on timeout we kill it and notify.
#   3. POST-RUN CHECK — scan output for auth/connection markers, then notify
#      success or failure.
#
# Run it:   ./run-bf-triage.sh         (manual, or from launchd/cron)
# Tune it:  override any UPPERCASE var via the environment.
#
# Scheduling note (macOS):
#   - While the macOS banner is the notifier, schedule via a launchd LaunchAgent,
#     NOT cron: cron runs outside your GUI session, so the osascript banner won't
#     appear. (Switch to the Slack webhook later and cron becomes viable.)
#   - The token must be fresh when this runs: do your morning device-flow login
#     (refresh-mcp-auth.sh) within ~12h before the scheduled time, OR rely on the
#     proxy's own silent refresh if Test A shows it survives the gap.
#
# Permission note:
#   Runs Claude with --permission-mode "${PERMISSION_MODE:-auto}". 'auto' is the
#   org default; 'bypassPermissions' is DISABLED by org policy, so don't rely on it.
#   Auto auto-approves the skill's read/CLI/file tools unattended; if a future op is
#   denied, add a targeted entry to .claude/settings.json 'permissions.allow'
#   rather than a blanket grant.

set -uo pipefail   # NOT -e: we handle errors explicitly so post-run checks run.

# ---------------------------------------------------------------------------
# Config (override via env)
# ---------------------------------------------------------------------------
TEAM="${TEAM:-Workload Resilience}"
MODEL="${MODEL:-claude-opus-4-8}"      # Opus 4.8. Use 'claude-opus-4-8[1m]' for 1M ctx, or alias 'opus'.
EFFORT="${EFFORT:-xhigh}"              # low | medium | high | xhigh | max
LIMIT="${LIMIT:-5}"                    # Mode B: number of BFs (skill default 5)
STATUSES="${BF_TRIAGE_TEAM_STATUSES:-Needs Triage,Open}"  # Mode B status filter (both active statuses)
JQL_EXTRA="${BF_TRIAGE_TEAM_JQL_EXTRA:-assignee is EMPTY}"  # extra Mode B JQL — default: only UNASSIGNED tickets
POST_COMMENTS="${POST_COMMENTS:-auto}" # auto (post dev-only summaries, DEFAULT) | never (report-only)

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]:-$0}")" && pwd)"
WORKDIR="${WORKDIR:-$SCRIPT_DIR}"      # defaults to this script's own dir (portable across clones)
OUTPUT_DIR="${BF_TRIAGE_OUTPUT_DIR:-$WORKDIR/bf-reports}"
LOG_DIR="${LOG_DIR:-$WORKDIR/bf-triage-logs}"

TOKEN_FILE="${TOKEN_FILE:-$HOME/.kanopy/token-devprod-mcp-proxy-prod.json}"  # shared prod token (same file your Claude Code sessions use). Seed once: ./refresh-mcp-auth.sh
MAX_TOKEN_AGE_HOURS="${MAX_TOKEN_AGE_HOURS:-48}"  # guard skip threshold; 48h passes weekday gaps, flags weekend (NOT a refresh requirement)
RUN_TIMEOUT_SECS="${RUN_TIMEOUT_SECS:-1800}"      # hard cap on the Claude run (30 min)

# Kill competing Claude Code processes before running (DEFAULT ON). One Claude Code
# window spawns its session + several --bg-spare helpers, EACH opening its own prod
# proxy on the shared refresh chain; those churn it and pop Okta tabs mid-run.
# Killing them first => this job is the ONLY proxy => no churn, no clicks.
#   *** DESTRUCTIVE: this terminates ALL your interactive Claude Code sessions. ***
#   Run this wrapper from a PLAIN TERMINAL, never from inside Claude Code's '!'.
#   To DISABLE (keep your sessions open and accept churn):
#       KILL_OTHER_CLAUDE=0 ./run-bf-triage.sh
KILL_OTHER_CLAUDE="${KILL_OTHER_CLAUDE:-1}"

# Notify channels. macOS banner is the ACTIVE channel for now; Slack is a
# placeholder until the webhook app is approved (set NOTIFY_SLACK_WEBHOOK then).
# NOTE: the osascript banner only appears from launchd or an at-desk GUI session
# — NOT from cron. So use a launchd LaunchAgent while the banner is the notifier.
NOTIFY_SLACK_WEBHOOK="${NOTIFY_SLACK_WEBHOOK:-}"  # placeholder: https://hooks.slack.com/services/... once approved
NOTIFY_OSASCRIPT="${NOTIFY_OSASCRIPT:-1}"         # 1 = macOS banner (active; launchd/at-desk only, not cron)

# ---------------------------------------------------------------------------
TS="$(date +%Y%m%d-%H%M%S)"
mkdir -p "$OUTPUT_DIR" "$LOG_DIR"
LOG="$LOG_DIR/bf-triage-$TS.log"

log() { printf '%s %s\n' "$(date '+%F %T')" "$*" | tee -a "$LOG" >&2; }

notify() {
  # $1 = subject, $2 = body
  local subject="$1" body="$2" text payload
  log "NOTIFY: $subject — $body"
  if [ -n "$NOTIFY_SLACK_WEBHOOK" ] && command -v python3 >/dev/null 2>&1; then
    text="*$subject*"$'\n'"$body"
    payload=$(python3 -c 'import json,sys; print(json.dumps({"text": sys.argv[1]}))' "$text")
    curl -fsS -X POST -H 'Content-type: application/json' --data "$payload" \
      "$NOTIFY_SLACK_WEBHOOK" >/dev/null 2>&1 || log "WARN: slack notify failed"
  fi
  if [ "$NOTIFY_OSASCRIPT" = "1" ] && command -v osascript >/dev/null 2>&1; then
    osascript -e "display notification \"${body//\"/\\\"}\" with title \"${subject//\"/\\\"}\"" 2>/dev/null || true
  fi
}

fail() { notify "BF-triage FAILED ($TEAM)" "$1  (log: $LOG)"; exit 1; }

# ---------------------------------------------------------------------------
# Clean slate: kill competing Claude Code processes + proxies (KILL_OTHER_CLAUDE)
# ---------------------------------------------------------------------------
if [ "$KILL_OTHER_CLAUDE" = "1" ]; then
  log "before kill: $(pgrep -x claude 2>/dev/null | wc -l | tr -d ' ') claude launcher(s), $(pgrep -f '/.local/share/claude/' 2>/dev/null | wc -l | tr -d ' ') engine proc(s), $(pgrep -f 'devprod-mcp-proxy' 2>/dev/null | wc -l | tr -d ' ') proxy(ies)"
  log "KILL_OTHER_CLAUDE=1: terminating other Claude Code processes + proxies for a single-proxy run"
  pkill -x claude 2>/dev/null || true                    # CLI launchers (claude -r/-c/daemon); our 'claude -p' isn't started yet
  pkill -f '/.local/share/claude/' 2>/dev/null || true   # app + version engine: sessions, --bg-spare, pty-hosts
  sleep 1
  pkill -f 'devprod-mcp-proxy' 2>/dev/null || true       # proxies last, so parents can't respawn them
  sleep 1
  left=$(pgrep -f 'devprod-mcp-proxy' 2>/dev/null | wc -l | tr -d ' ')
  log "after cleanup: ${left:-0} devprod-mcp-proxy process(es) remaining (want 0 before the job starts)"
else
  log "KILL_OTHER_CLAUDE=0: leaving other Claude processes alone — any running will churn the shared token (expect popups)"
fi

# ---------------------------------------------------------------------------
# Pre-flight guard
# ---------------------------------------------------------------------------
command -v claude  >/dev/null 2>&1 || fail "claude CLI not on PATH."
command -v python3 >/dev/null 2>&1 || fail "python3 required (token validation)."

[ -s "$TOKEN_FILE" ] || fail "proxy token missing ($TOKEN_FILE). Run: ./refresh-mcp-auth.sh and approve."
python3 -m json.tool "$TOKEN_FILE" >/dev/null 2>&1 \
  || fail "proxy token is not valid JSON ($TOKEN_FILE). Re-seed: ./refresh-mcp-auth.sh."
python3 -c "import json,sys; sys.exit(0 if json.load(open('$TOKEN_FILE')).get('refresh_token') else 1)" \
  || fail "proxy token has no refresh_token. Re-seed: ./refresh-mcp-auth.sh."

# Freshness heuristic: the proxy rewrites the token file on each refresh, so an
# old mtime means the session probably lapsed (e.g. overnight / weekend).
file_mtime=$(stat -f %m "$TOKEN_FILE" 2>/dev/null || echo 0)
age_hours=$(( ( $(date +%s) - file_mtime ) / 3600 ))
if [ "$age_hours" -ge "$MAX_TOKEN_AGE_HOURS" ]; then
  fail "proxy token is ${age_hours}h old (>= ${MAX_TOKEN_AGE_HOURS}h); session likely expired. Run: ./refresh-mcp-auth.sh before this job."
fi
log "pre-flight OK: token present, valid JSON, has refresh_token, ${age_hours}h old."

# ---------------------------------------------------------------------------
# Prompt: bf-triage Mode B for the team (report-only by default)
# ---------------------------------------------------------------------------
POST_CLAUSE="Do NOT post any Jira comments (report-only)."
[ "$POST_COMMENTS" = "auto" ] && POST_CLAUSE="Post the developers-only Jira comment summary for each BF and attach the pdf."

PROMPT="Use the bf-triage skill in Mode B to batch-triage the Build Failures assigned to the \"$TEAM\" team that are in status \"$STATUSES\" and currently UNASSIGNED (empty Assignee). Triage the top $LIMIT such BFs by creation date, writing one report per BF plus the team index under the configured output directory. This is an unattended scheduled run: do not ask interactive questions — proceed with sensible defaults. $POST_CLAUSE"

# ---------------------------------------------------------------------------
# Run Claude headlessly under a portable hard timeout
# ---------------------------------------------------------------------------
run_with_timeout() {
  local secs="$1"; shift
  "$@" </dev/null >>"$LOG" 2>&1 &
  local pid=$!
  ( sleep "$secs"; kill -0 "$pid" 2>/dev/null && { kill -TERM "$pid" 2>/dev/null; sleep 5; kill -KILL "$pid" 2>/dev/null; } ) &
  local watch=$!
  wait "$pid" 2>/dev/null; local rc=$?
  kill "$watch" 2>/dev/null; wait "$watch" 2>/dev/null
  return "$rc"
}

export BF_TRIAGE_OUTPUT_DIR="$OUTPUT_DIR"
export BF_TRIAGE_TEAM_STATUSES="$STATUSES"   # Mode B: restrict to these statuses (default: "Needs Triage,Open")
export BF_TRIAGE_TEAM_JQL_EXTRA="$JQL_EXTRA" # Mode B: extra JQL ANDed onto the base query (default: only unassigned)
# Make the post-comment decision deterministic for the skill (Step 0.5):
#   BF_TRIAGE_AUTO_POST_COMMENT=1 => posts a developers-only summary per BF
#   (set BF_TRIAGE_JIRA_COMMENT_PUBLIC=1 if you ever want public comments).
if [ "$POST_COMMENTS" = "auto" ]; then export BF_TRIAGE_AUTO_POST_COMMENT=1; else export BF_TRIAGE_AUTO_POST_COMMENT=0; fi
cd "$WORKDIR" || fail "cannot cd to WORKDIR ($WORKDIR)."

log "starting bf-triage (team='$TEAM' statuses='$STATUSES' model=$MODEL effort=$EFFORT limit=$LIMIT post=$POST_COMMENTS)"
log "output dir: $OUTPUT_DIR ; timeout: ${RUN_TIMEOUT_SECS}s"

RUN_MARKER="$LOG_DIR/.runstart-$TS"; : > "$RUN_MARKER"   # mtime = run start; lets us count only THIS run's reports

run_with_timeout "$RUN_TIMEOUT_SECS" \
  claude -p "$PROMPT" \
    --model "$MODEL" \
    --effort "$EFFORT" \
    --permission-mode "${PERMISSION_MODE:-auto}" \
    --add-dir "$OUTPUT_DIR"
rc=$?

# Timeout (TERM=143, KILL=137) almost always means an auth hang or stuck run.
if [ "$rc" -eq 143 ] || [ "$rc" -eq 137 ]; then
  fail "triage exceeded ${RUN_TIMEOUT_SECS}s and was killed — likely an auth hang (re-seed token) or a stuck run."
fi

# Even on a clean exit, surface auth / permission / connection trouble in the output.
if grep -qiE "opens a browser|browser for Okta|debug_connection_status|Not connected|MCP error -32000|reconnect .*devprod-mcp-gateway|haven't granted it yet|requested permission|permission-gated|could not start|tool family has been denied" "$LOG"; then
  fail "triage output shows an auth/permission/connection problem — check the log (re-seed the token, or grant the tool permissions)."
fi

[ "$rc" -ne 0 ] && fail "claude exited with code $rc."

# Success REQUIRES reports. Count ONLY this run's reports (newer than the run-start
# marker), including -vN versions — NOT same-day leftovers from earlier runs.
slug=$(printf '%s' "$TEAM" | tr '[:upper:] ' '[:lower:]-')
team_dir="$OUTPUT_DIR/team-$slug-$(date +%F)"
this_run=$(find "$team_dir" -name '*-triage*.md' -newer "$RUN_MARKER" 2>/dev/null)
rm -f "$RUN_MARKER"
count=$(printf '%s\n' "$this_run" | grep -c .)
if [ "${count:-0}" -eq 0 ]; then
  fail "run finished but produced 0 NEW reports — likely permission-gated, no matching BFs, or the skill didn't run. Check the log."
fi
bf_keys=$(printf '%s\n' "$this_run" | sed 's@.*/@@; s@-triage.*\.md$@@' | sort -u | paste -sd, -)
notify "BF-triage OK ($TEAM)" "Triaged $count BF(s): ${bf_keys:-(none)}"
log "DONE rc=0 ($count NEW reports this run: ${bf_keys:-none}) -> $team_dir"
