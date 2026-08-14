#!/usr/bin/env bash
# harness-streamliner: read-only fact gatherer.
# Surfaces config/perf issues across user + project Claude settings.
# NEVER mutates anything. Output is plain text grouped by pass for the model to read.
#
# Usage:
#   scan.sh [PROJECT_DIR]
#     PROJECT_DIR  optional path to a repo whose ./.claude/* should also be scanned.
#                  Defaults to $PWD. Pass "none" to scan user-level only.
#
# Requires: jq. Degrades gracefully (warns) if jq is missing.

set -uo pipefail

PROJECT_DIR="${1:-$PWD}"
HOME_CLAUDE="$HOME/.claude"

have_jq=1
command -v jq >/dev/null 2>&1 || have_jq=0

hr() { printf '%s\n' "----------------------------------------------------------------"; }
section() { printf '\n=== %s ===\n' "$1"; }

if [ "$have_jq" -eq 0 ]; then
  echo "WARN: jq not found. Install jq for full analysis. Falling back to file listing only."
fi

# Resolve the set of settings files to scan.
# User-level: settings.json, settings.local.json
# Project-level (if PROJECT_DIR != none and dir exists): .claude/settings.json, .claude/settings.local.json
declare -a SETTINGS_FILES=()
for f in "$HOME_CLAUDE/settings.json" "$HOME_CLAUDE/settings.local.json"; do
  [ -f "$f" ] && SETTINGS_FILES+=("$f")
done
if [ "$PROJECT_DIR" != "none" ] && [ -d "$PROJECT_DIR/.claude" ]; then
  for f in "$PROJECT_DIR/.claude/settings.json" "$PROJECT_DIR/.claude/settings.local.json"; do
    [ -f "$f" ] && SETTINGS_FILES+=("$f")
  done
fi

section "CONFIG FILES IN SCOPE"
printf 'HOME resolves to: %s\n' "$HOME"
if [ "${#SETTINGS_FILES[@]}" -eq 0 ]; then
  echo "No settings files found."
  exit 0
fi
for f in "${SETTINGS_FILES[@]}"; do
  lines=$(wc -l < "$f" | tr -d ' ')
  if [ "$have_jq" -eq 1 ]; then
    if jq -e . "$f" >/dev/null 2>&1; then valid="valid-json"; else valid="INVALID-JSON"; fi
  else
    valid="unknown"
  fi
  printf '  %-60s %5s lines  [%s]\n' "$f" "$lines" "$valid"
done

[ "$have_jq" -eq 0 ] && exit 0

# Helper: emit a JSON array's string elements one per line for a given jq path.
arr() { jq -r "($1 // []) | .[]" "$2" 2>/dev/null; }

# ---------------------------------------------------------------------------
# PASS 1: PERMISSIONS
# ---------------------------------------------------------------------------
section "PASS 1 — PERMISSIONS"
for f in "${SETTINGS_FILES[@]}"; do
  jq -e '.permissions' "$f" >/dev/null 2>&1 || continue
  printf '\n# %s\n' "$f"
  for bucket in allow deny ask; do
    items=()
    while IFS= read -r line; do [ -n "$line" ] && items+=("$line"); done < <(arr ".permissions.$bucket" "$f")
    n=${#items[@]}
    printf '  %s: %d rule(s)\n' "$bucket" "$n"
    [ "$n" -eq 0 ] && continue
    # exact duplicates
    dupes=$(printf '%s\n' "${items[@]}" | sort | uniq -d)
    [ -n "$dupes" ] && printf '    DUPLICATE exact rules:\n%s\n' "$(printf '%s\n' "$dupes" | sed 's/^/      /')"
    # Subsumption — CONSERVATIVE. Only flag clean directory-prefix parents
    # of the form Tool(<prefix>/**) or Tool(<prefix>/*) where <prefix> contains
    # no glob char. This avoids the false positives that leading-** suffix
    # matchers (e.g. Read(**/.env)) would otherwise generate. Fuzzier overlaps
    # (bash command globs, **/ suffix matchers) are left to the model — it
    # reads the full rule list below and reasons about glob semantics directly.
    for ((i=0;i<n;i++)); do
      a="${items[$i]}"
      case "$a" in
        *'/**)') ap="${a%\*\*)}" ;;   # strip trailing **)
        *'/*)')  ap="${a%\*)}"   ;;   # strip trailing *)
        *) continue ;;                # not a clean dir-glob parent
      esac
      case "$ap" in *'*'*) continue ;; esac   # prefix itself has a glob → skip
      for ((j=0;j<n;j++)); do
        [ "$i" -eq "$j" ] && continue
        b="${items[$j]}"
        case "$b" in
          "$ap"*) [ "$a" != "$b" ] && printf '    SUBSUMED: %s  ⊂  %s\n' "$b" "$a" ;;
        esac
      done
    done
  done
done

# autoMode.allow vs permissions.allow duplication (user file typically)
section "PASS 1b — autoMode.allow vs permissions.allow"
for f in "${SETTINGS_FILES[@]}"; do
  jq -e '.autoMode.allow' "$f" >/dev/null 2>&1 || continue
  overlap=$(jq -r '
    ((.autoMode.allow // []) as $am | (.permissions.allow // []) as $pa
     | $am | map(select(. as $x | $pa | index($x))) | .[])' "$f" 2>/dev/null)
  if [ -n "$overlap" ]; then
    printf '# %s\n  autoMode.allow entries ALSO in permissions.allow (redundant):\n%s\n' \
      "$f" "$(printf '%s\n' "$overlap" | sed 's/^/    /')"
  fi
done

# ---------------------------------------------------------------------------
# PASS 2: HOOKS  (broken paths + per-event/startup load)
# ---------------------------------------------------------------------------
section "PASS 2 — HOOKS"
# Extract every hook command with its event name, then resolve the script path.
for f in "${SETTINGS_FILES[@]}"; do
  jq -e '.hooks' "$f" >/dev/null 2>&1 || continue
  printf '\n# %s\n' "$f"
  # event \t command
  jq -r '.hooks | to_entries[] | .key as $ev
         | (.value[]?.hooks[]?.command // empty) | "\($ev)\t\(.)"' "$f" 2>/dev/null \
  | while IFS=$'\t' read -r ev cmd; do
      # pull the most likely script path token: after node/npx/python3/bash/sh, else first token
      read -r -a toks <<< "$cmd"
      path=""
      case "${toks[0]}" in
        node|npx|python3|python|bash|sh|deno|bun|tsx) path="${toks[1]:-}";;
        *) path="${toks[0]}";;
      esac
      # strip args like "reset" / "|| true" by taking only the path-looking token
      case "$path" in
        /*|~*|./*|../*) :;;
        *) path="";;   # not a path (e.g. inline command); skip existence check
      esac
      status=""
      if [ -n "$path" ]; then
        expanded="${path/#\~/$HOME}"
        if [ -e "$expanded" ]; then status="ok"; else status="MISSING"; fi
      else
        status="inline/no-path"
      fi
      printf '  [%-16s] %-9s %s\n' "$ev" "$status" "$cmd"
    done
  # per-event load summary
  printf '  -- load summary --\n'
  jq -r '.hooks | to_entries[]
         | "\(.key): \([.value[]?.hooks[]?] | length) hook(s)"' "$f" 2>/dev/null \
    | sed 's/^/    /'
  echo "    (SessionStart hooks = startup cost; UserPromptSubmit hooks = per-turn cost — measure these.)"
done

# ---------------------------------------------------------------------------
# PASS 3: MCP SERVERS
# ---------------------------------------------------------------------------
section "PASS 3 — MCP SERVERS"
for f in "${SETTINGS_FILES[@]}" "$HOME_CLAUDE/mcp.json"; do
  [ -f "$f" ] || continue
  jq -e '.mcpServers' "$f" >/dev/null 2>&1 || continue
  count=$(jq -r '.mcpServers | length' "$f")
  printf '\n# %s  (%s server(s) — each consumes tool/context budget)\n' "$f" "$count"
  jq -r '.mcpServers | to_entries[]
         | "  \(.key): \((.value.command // .value.httpUrl // .value.url // "?")) \((.value.args // []) | join(" "))"' \
    "$f" 2>/dev/null | cut -c1-160
  # duplicate command+args signatures (near-identical servers)
  dup=$(jq -r '.mcpServers | to_entries
         | map({k:.key, sig:((.value.command // .value.httpUrl // .value.url // "")+" "+((.value.args // [])|join(" ")))})
         | group_by(.sig) | map(select(length>1)) | .[] | map(.k) | join(", ")' "$f" 2>/dev/null)
  [ -n "$dup" ] && printf '  DUPLICATE/near-identical servers: %s\n' "$dup"
  # plaintext secret heuristic in args/env values
  secrets=$(jq -r '.mcpServers | to_entries[] | .key as $k
            | ((.value.args // []) + ([.value.env // {} | to_entries[] | .value]))
            | .[]? | select(type=="string")
            | select(test("(?i)(secret|token|api[_-]?key|client[_-]?secret|password)") or test("^eyJ[A-Za-z0-9_-]{10,}\\.") or test("(?i)(sk|sa_sk|mdb_sa_sk)_[A-Za-z0-9]{12,}"))
            | "  PLAINTEXT-SECRET in \($k): " + (.[0:24] + "…")' "$f" 2>/dev/null)
  [ -n "$secrets" ] && printf '%s\n' "$secrets"
done

# ---------------------------------------------------------------------------
# PASS 4: ENV / MODEL / STATUSLINE / MISC
# ---------------------------------------------------------------------------
section "PASS 4 — ENV / MODEL / STATUSLINE / MISC"
for f in "${SETTINGS_FILES[@]}"; do
  printf '\n# %s\n' "$f"
  # two different defaultMode keys (top-level vs permissions.defaultMode)
  top=$(jq -r '.defaultMode // empty' "$f")
  pm=$(jq -r '.permissions.defaultMode // empty' "$f")
  if [ -n "$top" ] && [ -n "$pm" ] && [ "$top" != "$pm" ]; then
    printf '  CONFLICT: defaultMode="%s" but permissions.defaultMode="%s" — confirm intended.\n' "$top" "$pm"
  fi
  # statusline command path existence
  sl=$(jq -r '.statusLine.command // empty' "$f")
  if [ -n "$sl" ]; then
    slpath="${sl%% *}"; slpath="${slpath/#\~/$HOME}"
    if [ -e "$slpath" ]; then echo "  statusLine: ok ($sl)"; else echo "  statusLine: MISSING SCRIPT ($sl)"; fi
  fi
  # skillListingBudgetFraction sanity
  slbf=$(jq -r '.skillListingBudgetFraction // empty' "$f")
  [ -n "$slbf" ] && printf '  skillListingBudgetFraction=%s (typical 0.05–0.20; higher = more skill metadata in context)\n' "$slbf"
  # env keys
  envk=$(jq -r '(.env // {}) | keys | join(", ")' "$f")
  [ "$envk" != "" ] && printf '  env: %s\n' "$envk"
  # model
  m=$(jq -r '.model // empty' "$f")
  [ -n "$m" ] && printf '  model: %s\n' "$m"
done

section "SCAN COMPLETE"
echo "Findings above are CANDIDATES. Apply the severity rubric in SKILL.md before changing anything."
echo "Secrets and conflicting modes are FLAG-ONLY (never auto-edit). Back up before any write."
