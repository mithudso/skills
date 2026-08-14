#!/usr/bin/env bash
# list_active_bfs.sh — Mode B helper for the bf-triage skill.
#
# This script is a prompt-builder, NOT a Jira client. Shells cannot call
# MCP tools, so this script:
#   1. Resolves Mode-B defaults (team, limit, statuses, extra-jql) from
#      env vars, the skill config file, and CLI flags.
#   2. Prints the exact JQL the skill will send to
#      `devprod-mcp-gateway.jira_search_issues` in Mode B.
#   3. Prints a ready-to-paste agent prompt that invokes the skill in
#      Mode B with the resolved params.
#
# The user pastes the prompt into their agent chat (Cursor / Claude
# Code); the actual `jira_search_issues` call happens inside the skill
# (which has MCP access).
#
# Priority order (highest wins):
#   1. CLI flags (--team / --limit / --statuses / --extra-jql)
#   2. Environment variables (BF_TRIAGE_TEAM / BF_TRIAGE_TEAM_LIMIT /
#      BF_TRIAGE_TEAM_STATUSES / BF_TRIAGE_TEAM_JQL_EXTRA)
#   3. Config file `<skill-dir>/config.yaml` — where <skill-dir> is the
#      directory containing this script's parent (auto-detected from $0,
#      or overridden via BF_TRIAGE_SKILL_DIR env var).
#   4. Built-in defaults (limit=5, statuses="Needs Triage,Open")

set -euo pipefail

# Self-locate the skill directory: scripts/ is one level under the root.
# Caller may override via BF_TRIAGE_SKILL_DIR.
if [[ -z "${BF_TRIAGE_SKILL_DIR:-}" ]]; then
  BF_TRIAGE_SKILL_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
fi

usage() {
  cat <<'USAGE'
Usage: list_active_bfs.sh --team "Team Name" [--limit N] [--statuses "..."] [--extra-jql "..."]

Prints the JQL the bf-triage skill will use in Mode B (team batch),
plus a ready-to-paste agent prompt.

Required:
  --team NAME         Team name as it appears in Jira's "Assigned Teams"
                      field (e.g. "Workload Resilience"). Can also be
                      supplied via env var BF_TRIAGE_TEAM.

Optional overrides:
  --limit N           Max BFs to triage. Default 5 (env BF_TRIAGE_TEAM_LIMIT).
                      Hard cap 50 (the jira_search_issues max).
  --statuses LIST     Comma-separated active statuses. Default
                      "Needs Triage,Open" (env BF_TRIAGE_TEAM_STATUSES).
  --extra-jql CLAUSE  Extra JQL ANDed onto the base query. Default empty
                      (env BF_TRIAGE_TEAM_JQL_EXTRA).
  -h, --help          Show this help.

Examples:
  ./list_active_bfs.sh --team "Workload Resilience"
  ./list_active_bfs.sh --team "Workload Resilience" --limit 3
  ./list_active_bfs.sh --team "DevProd Performance Infrastructure" \
                       --extra-jql 'Temperature ~ "hot"'

  BF_TRIAGE_TEAM_LIMIT=10 ./list_active_bfs.sh --team "Query Execution"

This script never modifies state. It only prints to stdout.
USAGE
}

CONFIG_FILE="${BF_TRIAGE_SKILL_DIR}/config.yaml"
DEFAULT_LIMIT=5
DEFAULT_STATUSES="Needs Triage,Open"
DEFAULT_EXTRA_JQL=""
HARD_CAP=50

TEAM=""
LIMIT=""
STATUSES=""
EXTRA_JQL=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help) usage; exit 0 ;;
    --team) TEAM="${2:?--team needs a value}"; shift 2 ;;
    --limit) LIMIT="${2:?--limit needs a value}"; shift 2 ;;
    --statuses) STATUSES="${2:?--statuses needs a value}"; shift 2 ;;
    --extra-jql) EXTRA_JQL="${2:?--extra-jql needs a value}"; shift 2 ;;
    *) echo "ERROR: unknown argument: $1" >&2; usage; exit 3 ;;
  esac
done

read_config_key() {
  local key="$1"
  [[ -f "$CONFIG_FILE" ]] || return 0
  awk -v k="$key" '
    $0 ~ "^"k":" {
      sub("^"k":[[:space:]]*", "")
      sub("^\"", ""); sub("\"$", "")
      sub("^'\''", ""); sub("'\''$", "")
      print
      exit
    }
  ' "$CONFIG_FILE"
}

[[ -z "$TEAM" ]] && TEAM="${BF_TRIAGE_TEAM:-$(read_config_key team || true)}"
[[ -z "$LIMIT" ]] && LIMIT="${BF_TRIAGE_TEAM_LIMIT:-$(read_config_key team_limit || true)}"
[[ -z "$STATUSES" ]] && STATUSES="${BF_TRIAGE_TEAM_STATUSES:-$(read_config_key team_statuses || true)}"
[[ -z "$EXTRA_JQL" ]] && EXTRA_JQL="${BF_TRIAGE_TEAM_JQL_EXTRA:-$(read_config_key team_jql_extra || true)}"

[[ -z "$LIMIT" ]] && LIMIT="$DEFAULT_LIMIT"
[[ -z "$STATUSES" ]] && STATUSES="$DEFAULT_STATUSES"
[[ -z "$EXTRA_JQL" ]] && EXTRA_JQL="$DEFAULT_EXTRA_JQL"

if [[ -z "$TEAM" ]]; then
  echo "ERROR: --team is required (or set BF_TRIAGE_TEAM env var, or 'team:' in $CONFIG_FILE)" >&2
  usage
  exit 3
fi

if ! [[ "$LIMIT" =~ ^[0-9]+$ ]] || (( LIMIT < 1 )); then
  echo "ERROR: --limit must be a positive integer, got: $LIMIT" >&2
  exit 3
fi

if (( LIMIT > HARD_CAP )); then
  echo "ERROR: --limit $LIMIT exceeds the jira_search_issues hard cap ($HARD_CAP)." >&2
  echo "       Re-run with --limit $HARD_CAP or fewer." >&2
  exit 3
fi

QUOTED_STATUSES=$(echo "$STATUSES" \
  | awk -F',' '{
      for (i=1;i<=NF;i++){
        s=$i; gsub(/^[[:space:]]+|[[:space:]]+$/,"",s);
        printf "%s\"%s\"", (i>1?", ":""), s
      }
    }')

JQL="project = BF AND \"Assigned Teams\" = \"${TEAM}\" AND status in (${QUOTED_STATUSES})"
[[ -n "$EXTRA_JQL" ]] && JQL="${JQL} AND ${EXTRA_JQL}"
JQL="${JQL} ORDER BY created DESC"

cat <<EOF
# bf-triage Mode B — resolved invocation

Team:           ${TEAM}
Limit:          ${LIMIT}
Statuses:       ${STATUSES}
Extra JQL:      ${EXTRA_JQL:-<none>}
Config source:  $([[ -f "$CONFIG_FILE" ]] && echo "$CONFIG_FILE (+ env + CLI)" || echo "env + CLI + built-in defaults")

## JQL the skill will send to devprod-mcp-gateway.jira_search_issues

${JQL}

## Ready-to-paste agent prompt

triage open BFs for "${TEAM}" limit=${LIMIT} statuses="${STATUSES}"$([[ -n "$EXTRA_JQL" ]] && echo " extra-jql=\"${EXTRA_JQL}\"")

EOF
