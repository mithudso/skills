#!/usr/bin/env bash
# Slice a Jira BF ticket into held-in and held-out files for v2 verification.
#
# This script implements the "Verification Methodology v2 — Process Isolation"
# pattern: the coordinator pre-fetches each BF ticket, the slicer splits it
# into held-in / held-out files at a cutoff timestamp, then the coordinator
# spawns isolated triager and grader subagents in parallel.
#
# Responsibilities of this script (the harness):
#   - Read pre-fetched Jira JSON for the BF (either the legacy single-blob
#     shape or the devprod-mcp-gateway two-call shape: issue.json + comments.json).
#   - When given gateway shape, normalize to legacy via _stitch_gateway_jira.py.
#   - Compute the cutoff timestamp using the following resolution order:
#       1. First changelog event setting `Assigned Teams = OWNING_TEAM`
#          (legacy shape only; the gateway has no changelog tool).
#       2. `custom_fields["Team Assigned (Effective Date)"]` from the
#          gateway shape, when `custom_fields["Assigned Teams"]` includes
#          OWNING_TEAM. Used because the gateway does not expose a changelog.
#   - Write three files in ${OUTPUT_DIR}/<BF>/:
#       sliced.md   - everything with timestamp < cutoff (input for triager)
#       heldout.md  - everything with timestamp >= cutoff (input for grader)
#       cutoff.txt  - audit log: cutoff timestamp + event description
#
# Responsibilities of the COORDINATOR (the agent that runs this script):
#   - Pre-fetch the BF via `devprod-mcp-gateway`:
#       (a) `jira_get_issue`           -> save to ${WORK_DIR}/raw.json
#       (b) `jira_get_issue_comments`  -> save to ${WORK_DIR}/comments.json
#     The gateway does NOT expose a changelog tool today; we therefore
#     reconstruct the cutoff from the `Team Assigned (Effective Date)`
#     custom field. Held-in field snapshots become approximations
#     (current-state instead of value-at-cutoff); the produced sliced.md
#     carries a warning banner saying so.
#     If the gateway is down, follow SKILL.md Hard rule 5 (retry once,
#     then stop and ask the user to fix access). The skill does NOT fall
#     back to `user-atlassian` or any other Jira MCP.
#   - Invoke this script with `--issue-json` (and `--comments-json` when
#     using gateway shape), OR with `--input-json` for a legacy single-blob.
#
# Responsibilities of the TRIAGER subagent:
#   - Read ONLY ${OUTPUT_DIR}/<BF>/sliced.md.
#   - MUST NOT call any `jira_*` tool against the BF under test
#     (see templates/triager_prompt.md for the full forbidden-actions list).

set -euo pipefail

usage() {
  cat <<'USAGE'
Usage:
  run_held_in_test.sh BF-KEY [--issue-json PATH] [--comments-json PATH] \
                             [--input-json PATH] [--output-dir DIR] \
                             [--owning-team NAME]

Slices a Jira BF ticket into held-in and held-out files for v2 verification.

Required inputs (pick ONE of the two modes below):

  MODE 1: gateway two-call shape (recommended; matches what
  `devprod-mcp-gateway` actually returns):
    --issue-json     PATH   Output of jira_get_issue.
                            Defaults to ${OUTPUT_DIR}/<BF>/raw.json.
    --comments-json  PATH   Output of jira_get_issue_comments.
                            Defaults to ${OUTPUT_DIR}/<BF>/comments.json
                            (optional; missing file = empty comments).

  MODE 2: legacy single-blob shape (with top-level `changelogs` and
  `comments` arrays, as the classic Jira REST API would return):
    --input-json     PATH   Path to the single legacy-shape JSON file.
                            Mutually exclusive with --issue-json /
                            --comments-json.

Options:
  --output-dir   DIR    Output base directory. Default: /tmp/bf-triage-test
  --owning-team  NAME   Team whose first `Assigned Teams = NAME` event (or
                        `Team Assigned (Effective Date)` value, in MODE 1)
                        anchors the cutoff. Default: "Workload Resilience".

Outputs (idempotent — re-running overwrites):
  ${OUTPUT_DIR}/<BF>/sliced.md
  ${OUTPUT_DIR}/<BF>/heldout.md
  ${OUTPUT_DIR}/<BF>/cutoff.txt
  ${OUTPUT_DIR}/<BF>/normalized.json   (only in MODE 1)

Exit codes:
  0  Success.
  1  Cannot determine cutoff (neither changelog nor Team Assigned
     (Effective Date) yields a timestamp).
  2  Required input JSON not found.
  3  Bad arguments / dependency missing (jq, python3).
USAGE
}

BF_KEY=""
INPUT_JSON=""
ISSUE_JSON=""
COMMENTS_JSON=""
OUTPUT_DIR="/tmp/bf-triage-test"
OWNING_TEAM="Workload Resilience"

while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help)
      usage
      exit 0
      ;;
    --input-json)
      INPUT_JSON="${2:-}"
      [[ -z "$INPUT_JSON" ]] && { echo "ERROR: --input-json needs a path" >&2; exit 3; }
      shift 2
      ;;
    --issue-json)
      ISSUE_JSON="${2:-}"
      [[ -z "$ISSUE_JSON" ]] && { echo "ERROR: --issue-json needs a path" >&2; exit 3; }
      shift 2
      ;;
    --comments-json)
      COMMENTS_JSON="${2:-}"
      [[ -z "$COMMENTS_JSON" ]] && { echo "ERROR: --comments-json needs a path" >&2; exit 3; }
      shift 2
      ;;
    --output-dir)
      OUTPUT_DIR="${2:-}"
      [[ -z "$OUTPUT_DIR" ]] && { echo "ERROR: --output-dir needs a path" >&2; exit 3; }
      shift 2
      ;;
    --owning-team)
      OWNING_TEAM="${2:-}"
      [[ -z "$OWNING_TEAM" ]] && { echo "ERROR: --owning-team needs a value" >&2; exit 3; }
      shift 2
      ;;
    BF-*|bf-*)
      BF_KEY="$1"
      shift
      ;;
    *)
      echo "ERROR: unknown argument: $1" >&2
      usage
      exit 3
      ;;
  esac
done

if [[ -z "$BF_KEY" ]]; then
  echo "ERROR: BF key required (e.g. BF-43272)" >&2
  usage
  exit 3
fi
BF_KEY="${BF_KEY^^}"

if [[ -n "$INPUT_JSON" && ( -n "$ISSUE_JSON" || -n "$COMMENTS_JSON" ) ]]; then
  echo "ERROR: --input-json is mutually exclusive with --issue-json / --comments-json" >&2
  exit 3
fi

for cmd in jq python3; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "ERROR: required dependency '$cmd' not found in PATH" >&2
    exit 3
  fi
done

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SLICE_HELPER="${SCRIPT_DIR}/_slice_helper.py"
STITCH_HELPER="${SCRIPT_DIR}/_stitch_gateway_jira.py"

for helper in "$SLICE_HELPER" "$STITCH_HELPER"; do
  if [[ ! -f "$helper" ]]; then
    echo "ERROR: helper not found at $helper" >&2
    exit 3
  fi
done

WORK_DIR="${OUTPUT_DIR}/${BF_KEY}"
mkdir -p "$WORK_DIR"

# Default paths if not supplied.
if [[ -z "$INPUT_JSON" && -z "$ISSUE_JSON" ]]; then
  ISSUE_JSON="${WORK_DIR}/raw.json"
fi
if [[ -z "$INPUT_JSON" && -z "$COMMENTS_JSON" && -f "${WORK_DIR}/comments.json" ]]; then
  COMMENTS_JSON="${WORK_DIR}/comments.json"
fi

# Resolve which input to feed to the slicer.
SLICER_INPUT=""
NORMALIZED_FROM_GATEWAY="false"

if [[ -n "$INPUT_JSON" ]]; then
  if [[ ! -f "$INPUT_JSON" ]]; then
    cat >&2 <<EOF
ERROR: legacy-shape JSON not found at: $INPUT_JSON

Use --issue-json (and --comments-json) for the gateway two-call shape,
or pre-write a single legacy-shape blob to --input-json.
EOF
    exit 2
  fi
  if ! jq empty "$INPUT_JSON" >/dev/null 2>&1; then
    echo "ERROR: $INPUT_JSON is not valid JSON" >&2
    exit 3
  fi

  # Auto-detect: if it has no top-level `changelogs` AND has `custom_fields`,
  # it's actually the gateway issue shape misrouted to --input-json. Stitch it.
  if jq -e '(has("changelogs") | not) and has("custom_fields")' "$INPUT_JSON" >/dev/null; then
    ISSUE_JSON="$INPUT_JSON"
    INPUT_JSON=""
  else
    SLICER_INPUT="$INPUT_JSON"
  fi
fi

if [[ -z "$SLICER_INPUT" ]]; then
  if [[ ! -f "$ISSUE_JSON" ]]; then
    cat >&2 <<EOF
ERROR: issue JSON not found at: $ISSUE_JSON

The COORDINATOR must pre-fetch the BF via the gateway and save it:

  CallMcpTool devprod-mcp-gateway jira_get_issue {
    "issue_key": "$BF_KEY"
  }
  -> save response to: $ISSUE_JSON

  CallMcpTool devprod-mcp-gateway jira_get_issue_comments {
    "issue_key": "$BF_KEY",
    "limit": 200
  }
  -> save response to: ${COMMENTS_JSON:-${WORK_DIR}/comments.json}

If the gateway returns an error (unauthenticated / Not connected / TLS):
retry once, then STOP and ask the user to fix gateway access per
SKILL.md Hard rule 5. Do NOT fall back to user-atlassian.
EOF
    exit 2
  fi
  if ! jq empty "$ISSUE_JSON" >/dev/null 2>&1; then
    echo "ERROR: $ISSUE_JSON is not valid JSON" >&2
    exit 3
  fi
  if [[ -n "$COMMENTS_JSON" && ! -f "$COMMENTS_JSON" ]]; then
    echo "WARNING: comments JSON not found at $COMMENTS_JSON; held-in/held-out will have no comments" >&2
    COMMENTS_JSON=""
  fi
  if [[ -n "$COMMENTS_JSON" ]] && ! jq empty "$COMMENTS_JSON" >/dev/null 2>&1; then
    echo "ERROR: $COMMENTS_JSON is not valid JSON" >&2
    exit 3
  fi

  STITCHED_JSON="${WORK_DIR}/normalized.json"
  STITCH_ARGS=(--issue-json "$ISSUE_JSON" --out "$STITCHED_JSON")
  if [[ -n "$COMMENTS_JSON" ]]; then
    STITCH_ARGS+=(--comments-json "$COMMENTS_JSON")
  fi
  python3 "$STITCH_HELPER" "${STITCH_ARGS[@]}" >/dev/null
  SLICER_INPUT="$STITCHED_JSON"
  NORMALIZED_FROM_GATEWAY="true"
fi

# Cutoff resolution.
#   Path 1: legacy shape -> first changelog event setting Assigned Teams = team.
#   Path 2: gateway shape -> Team Assigned (Effective Date) custom field
#           (requires `Assigned Teams` to currently include OWNING_TEAM).
CUTOFF=""
CUTOFF_EVENT=""
CUTOFF_SOURCE=""

CUTOFF=$(jq -r --arg team "$OWNING_TEAM" '
  [ (.changelogs // [])[]
    | select(any(.items[]?;
        .field == "Assigned Teams"
        and ((.to_string // "") | contains($team))))
    | .created
  ]
  | sort
  | (first // "")
' "$SLICER_INPUT")

if [[ -n "$CUTOFF" && "$CUTOFF" != "null" ]]; then
  CUTOFF_EVENT=$(jq -r --arg ts "$CUTOFF" '
    .changelogs[]
    | select(.created == $ts)
    | "[\(.created)] author=\((.author.display_name // "?")) sets: " +
      ([.items[]? | "\(.field)=\((.to_string // "(empty)"))"] | join("; "))
  ' "$SLICER_INPUT" | head -n 1)
  CUTOFF_SOURCE="changelog"
else
  CUTOFF=""
fi

if [[ -z "$CUTOFF" ]]; then
  CUTOFF=$(jq -r --arg team "$OWNING_TEAM" '
    if (.custom_fields // {}) | has("Team Assigned (Effective Date)") then
      ( (.custom_fields["Assigned Teams"] // [])
        | (if type == "array" then . else [.] end)
        | map(tostring)
        | any(. == $team)
      ) as $matches
      | if $matches
        then (.custom_fields["Team Assigned (Effective Date)"] // "")
        else ""
        end
    else ""
    end
  ' "$SLICER_INPUT")

  if [[ -n "$CUTOFF" && "$CUTOFF" != "null" ]]; then
    CUTOFF_EVENT="[gateway custom field] Team Assigned (Effective Date) = $CUTOFF; current Assigned Teams includes '$OWNING_TEAM'"
    CUTOFF_SOURCE="custom_field"
  else
    CUTOFF=""
  fi
fi

if [[ -z "$CUTOFF" ]]; then
  cat >&2 <<EOF
ERROR: $BF_KEY: cannot determine cutoff.
  - No changelog event sets 'Assigned Teams' to '$OWNING_TEAM'.
  - No 'Team Assigned (Effective Date)' custom field with 'Assigned Teams'
    including '$OWNING_TEAM' either.
  Either the BF was never routed to '$OWNING_TEAM' or the gateway response
  is missing both signals. Try --owning-team with a sibling/former team name.
EOF
  exit 1
fi

cat > "$WORK_DIR/cutoff.txt" <<EOF
BF: $BF_KEY
Owning team (cutoff trigger): $OWNING_TEAM
Cutoff: $CUTOFF
Cutoff source: $CUTOFF_SOURCE
Event: $CUTOFF_EVENT
Source JSON: $SLICER_INPUT
Normalized from gateway: $NORMALIZED_FROM_GATEWAY
Generated by: run_held_in_test.sh (v2 harness)
EOF

COUNTS=$(python3 "$SLICE_HELPER" \
  --bf "$BF_KEY" \
  --cutoff "$CUTOFF" \
  --input "$SLICER_INPUT" \
  --sliced "$WORK_DIR/sliced.md" \
  --heldout "$WORK_DIR/heldout.md")

HI_C=$(echo "$COUNTS" | jq -r '.held_in_comments')
HO_C=$(echo "$COUNTS" | jq -r '.held_out_comments')
HI_CL=$(echo "$COUNTS" | jq -r '.held_in_changelogs')
HO_CL=$(echo "$COUNTS" | jq -r '.held_out_changelogs')

cat <<EOF
BF: $BF_KEY
Cutoff: $CUTOFF  (source: $CUTOFF_SOURCE)
Held-in:  $HI_C comments, $HI_CL changelog entries
Held-out: $HO_C comments, $HO_CL changelog entries
Files: $WORK_DIR/{sliced.md,heldout.md,cutoff.txt}
EOF
