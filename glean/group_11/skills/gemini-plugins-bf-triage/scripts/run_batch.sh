#!/usr/bin/env bash
# Parallel pre-slicer for the v2 verification harness.
#
# Usage: run_batch.sh [--output-dir DIR] [--input-dir DIR] BF-KEY...
#
# For each BF, this script invokes `run_held_in_test.sh` in the background
# and waits for all of them to finish. It does NOT itself fetch the BF
# JSON from Jira and it does NOT launch any triager / grader subagents —
# those are the coordinator's responsibility (see scripts/README.md).
#
# Inputs per BF (must be staged by the coordinator beforehand):
#   ${INPUT_DIR}/<BF>/raw.json    OR    ${OUTPUT_DIR}/<BF>/raw.json
#
# If --input-dir is given, this script will copy each <BF>/raw.json from
# ${INPUT_DIR}/<BF>/raw.json to ${OUTPUT_DIR}/<BF>/raw.json before slicing
# (so a single read-only fetch tree can be re-sliced in different output
# locations without re-fetching).
#
# Outputs per BF (written by the underlying run_held_in_test.sh):
#   ${OUTPUT_DIR}/<BF>/sliced.md
#   ${OUTPUT_DIR}/<BF>/heldout.md
#   ${OUTPUT_DIR}/<BF>/cutoff.txt
#
# Stdout: a final summary block listing one line per BF with cutoff +
# held-in / held-out counts. Per-BF stdout is captured into
# ${OUTPUT_DIR}/<BF>/slice.log and surfaced if the slicer fails.

set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: run_batch.sh [--output-dir DIR] [--input-dir DIR] BF-KEY...

Pre-slices N BFs in parallel using run_held_in_test.sh. The coordinator
must have already saved each BF's raw Jira JSON (see scripts/README.md).

Options:
  --output-dir DIR  Where slice outputs go. Default: /tmp/bf-triage-test
  --input-dir DIR   Where raw.json files live. Default: same as output-dir.
                    If different, raw.json is copied into the output tree.
  -h, --help        Show this help.

Example (after staging /tmp/bf-triage-test/<BF>/raw.json for each BF):
  run_batch.sh BF-43272 BF-43270 BF-41745 BF-41375 BF-39815 \
               BF-42258 BF-42203 BF-42067 BF-41872 BF-42845 BF-42797
USAGE
}

OUTPUT_DIR="/tmp/bf-triage-test"
INPUT_DIR=""
BF_KEYS=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help) usage; exit 0 ;;
    --output-dir) OUTPUT_DIR="${2:?--output-dir needs a value}"; shift 2 ;;
    --input-dir) INPUT_DIR="${2:?--input-dir needs a value}"; shift 2 ;;
    BF-*|bf-*) BF_KEYS+=("${1^^}"); shift ;;
    *) echo "ERROR: unknown argument: $1" >&2; usage; exit 3 ;;
  esac
done

if [[ ${#BF_KEYS[@]} -eq 0 ]]; then
  echo "ERROR: at least one BF key required" >&2
  usage
  exit 3
fi

[[ -z "$INPUT_DIR" ]] && INPUT_DIR="$OUTPUT_DIR"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SLICER="$SCRIPT_DIR/run_held_in_test.sh"

if [[ ! -x "$SLICER" ]]; then
  echo "ERROR: slicer not found or not executable: $SLICER" >&2
  exit 3
fi

mkdir -p "$OUTPUT_DIR"

slice_one() {
  local bf="$1"
  local work_dir="$OUTPUT_DIR/$bf"
  local input_json="$INPUT_DIR/$bf/raw.json"
  local log="$work_dir/slice.log"
  mkdir -p "$work_dir"

  if [[ ! -f "$input_json" ]]; then
    echo "MISSING_RAW_JSON $bf $input_json" > "$log"
    echo "MISSING_RAW_JSON: $bf (expected at $input_json)"
    return 4
  fi

  if [[ "$input_json" != "$work_dir/raw.json" ]]; then
    cp "$input_json" "$work_dir/raw.json"
  fi

  if "$SLICER" "$bf" --input-json "$work_dir/raw.json" --output-dir "$OUTPUT_DIR" > "$log" 2>&1; then
    grep -E '^(BF:|Cutoff:|Held-in:|Held-out:)' "$log" | sed "s/^/[$bf] /"
  else
    local rc=$?
    echo "[$bf] SLICE_FAILED rc=$rc — see $log" >&2
    sed "s/^/[$bf] /" "$log" >&2
    return "$rc"
  fi
}

export -f slice_one
export OUTPUT_DIR INPUT_DIR SLICER

declare -A PIDS=()
declare -A RC=()

for bf in "${BF_KEYS[@]}"; do
  ( slice_one "$bf" ) &
  PIDS["$bf"]=$!
done

FAIL=0
for bf in "${BF_KEYS[@]}"; do
  if ! wait "${PIDS[$bf]}"; then
    RC["$bf"]=$?
    FAIL=1
  else
    RC["$bf"]=0
  fi
done

echo
echo "===== Batch summary ====="
for bf in "${BF_KEYS[@]}"; do
  if [[ "${RC[$bf]}" -eq 0 ]]; then
    cutoff=$(grep '^Cutoff:' "$OUTPUT_DIR/$bf/cutoff.txt" 2>/dev/null | head -n 1 | sed 's/^Cutoff: //')
    counts=$(grep -E '^Held-(in|out):' "$OUTPUT_DIR/$bf/slice.log" 2>/dev/null | tr '\n' '|' | sed 's/|$//')
    printf "  %-12s OK   cutoff=%s  %s\n" "$bf" "$cutoff" "$counts"
  else
    printf "  %-12s FAIL rc=%s  log=%s\n" "$bf" "${RC[$bf]}" "$OUTPUT_DIR/$bf/slice.log"
  fi
done

exit "$FAIL"
