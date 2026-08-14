#!/bin/bash
#
# install-launchd.sh
# ------------------
# Generate a machine-specific LaunchAgent from com.example.bf-triage.plist.template
# and load it into the per-user GUI domain (so the osascript banner can appear).
# Idempotent: re-running regenerates the plist and reloads it.
#
# Portable: derives all absolute paths from where this repo is cloned and from
# the current user's $HOME / installed tools — nothing is hardcoded.
#
# Usage:
#   ./install-launchd.sh
#   LABEL=com.acme.bf-triage HOUR=12 MINUTE=0 ./install-launchd.sh
#   NOTIFY_SLACK_WEBHOOK="https://hooks.slack.com/services/..." ./install-launchd.sh
#
# Uninstall:
#   launchctl bootout gui/$(id -u)/<label> && rm ~/Library/LaunchAgents/<label>.plist

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]:-$0}")" && pwd)"
TEMPLATE="$SCRIPT_DIR/com.example.bf-triage.plist.template"
TARGET_SCRIPT="$SCRIPT_DIR/run-bf-triage.sh"

LABEL="${LABEL:-com.example.bf-triage}"
HOUR="${HOUR:-12}"
MINUTE="${MINUTE:-0}"
NOTIFY_SLACK_WEBHOOK="${NOTIFY_SLACK_WEBHOOK:-}"

WORKDIR="$SCRIPT_DIR"
LOGDIR="$WORKDIR/bf-triage-logs"
AGENT_DIR="$HOME/Library/LaunchAgents"
PLIST="$AGENT_DIR/$LABEL.plist"

[ -f "$TEMPLATE" ]      || { echo "ERROR: template not found: $TEMPLATE" >&2; exit 1; }
[ -x "$TARGET_SCRIPT" ] || { echo "ERROR: run-bf-triage.sh missing/not executable: $TARGET_SCRIPT" >&2; exit 1; }

# Build a PATH from the dirs of the tools the job needs, plus standard dirs (deduped).
collect_path() {
  local dirs=() seen="" b d
  for b in claude python3 kanopy-oidc curl; do
    if d="$(command -v "$b" 2>/dev/null)"; then dirs+=("$(dirname "$d")"); fi
  done
  dirs+=("$HOME/.local/bin" /opt/homebrew/bin /usr/local/bin /usr/bin /bin /usr/sbin /sbin)
  local out="" x
  for x in "${dirs[@]}"; do
    case ":$seen:" in *":$x:"*) continue ;; esac
    seen="$seen:$x"; out="${out:+$out:}$x"
  done
  printf '%s' "$out"
}
PATH_VALUE="$(collect_path)"

mkdir -p "$AGENT_DIR" "$LOGDIR"

# Substitute placeholders. '|' is the sed delimiter since values contain '/'.
sed \
  -e "s|__LABEL__|$LABEL|g" \
  -e "s|__SCRIPT__|$TARGET_SCRIPT|g" \
  -e "s|__PATH__|$PATH_VALUE|g" \
  -e "s|__HOME__|$HOME|g" \
  -e "s|__WORKDIR__|$WORKDIR|g" \
  -e "s|__LOGDIR__|$LOGDIR|g" \
  -e "s|__HOUR__|$HOUR|g" \
  -e "s|__MINUTE__|$MINUTE|g" \
  -e "s|__SLACK_WEBHOOK__|$NOTIFY_SLACK_WEBHOOK|g" \
  "$TEMPLATE" > "$PLIST"

plutil -lint "$PLIST" >/dev/null || { echo "ERROR: generated plist failed lint: $PLIST" >&2; exit 1; }

# (Re)load into the GUI domain.
launchctl bootout  "gui/$(id -u)/$LABEL" 2>/dev/null || true
launchctl bootstrap "gui/$(id -u)" "$PLIST"

printf 'installed: %s\n' "$PLIST"
printf 'schedule:  Mon-Fri %02d:%02d local\n' "$HOUR" "$MINUTE"
printf 'slack:     %s\n' "${NOTIFY_SLACK_WEBHOOK:+configured}${NOTIFY_SLACK_WEBHOOK:-(none — macOS banner only)}"
printf 'test now:  launchctl kickstart -k gui/%s/%s\n' "$(id -u)" "$LABEL"
printf 'uninstall: launchctl bootout gui/%s/%s && rm %q\n' "$(id -u)" "$LABEL" "$PLIST"
