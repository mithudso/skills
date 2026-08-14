#!/usr/bin/env bash
# setup_repos.sh — idempotent shallow-but-history-preserving clone helper for
# the bf-triage Cursor skill.
#
# Clones 10gen/mongo and/or 10gen/dsi into /tmp/bf-triage-workdir-<random>/
# with `--filter=blob:none --no-checkout` so `git log` and `git show` work
# across any time window but the on-disk size stays small. Files are fetched
# on demand the first time they are checked out.
#
# Usage:
#   setup_repos.sh                # clone both (or reuse existing) and print export lines
#   setup_repos.sh both           # explicit: clone both mongo + dsi
#   setup_repos.sh mongo          # mongo only (skip dsi — useful for non-sys-perf BFs)
#   setup_repos.sh dsi            # dsi only (rare, but symmetric)
#   setup_repos.sh cleanup        # delete the workdir referenced by .last_workdir
#   setup_repos.sh cleanup all    # delete every /tmp/bf-triage-workdir-* directory
#
# The default (no arg) is `both` to preserve backward compatibility with
# callers written before per-repo selection existed.
#
# The script is read-only against the target repos: it never pushes, never
# commits, never modifies remote configuration. When a workdir is reused, only
# the repos in the current selection are ensured-present; any repo previously
# cloned under that workdir is left untouched (the caller's selection only
# governs what gets exported, not what gets removed).

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
STATE_DIR="${HOME}/.cache/bf-triage"
STATE_FILE="${STATE_DIR}/last_workdir"

REPOS=(
    "git@github.com:10gen/mongo.git mongo"
    "git@github.com:10gen/dsi.git dsi"
)

usage() {
    cat <<EOF
setup_repos.sh — idempotent helper for the bf-triage Cursor skill.

Usage:
  setup_repos.sh [both|mongo|dsi]  clone (or reuse) and print export lines
                                   default: both
  setup_repos.sh cleanup           delete the workdir referenced by ${STATE_FILE}
  setup_repos.sh cleanup all       delete every /tmp/bf-triage-workdir-* directory
  setup_repos.sh -h | --help       show this help

The clones use 'git clone --filter=blob:none --no-checkout' so 'git log' and
'git show' work over time windows while keeping disk usage small.

Use 'mongo' (skip dsi) when the BF's Evergreen project is not sys-perf /
not a DSI failure — Step 6 bisect only consults 10gen/dsi for sys-perf-style
tickets, so cloning it is wasted network for unrelated BFs.
EOF
}

ensure_state_dir() {
    mkdir -p "${STATE_DIR}"
}

clone_one() {
    local url="$1"
    local target_path="$2"

    if [ -d "${target_path}/.git" ]; then
        echo "[bf-triage] reusing existing clone at ${target_path}" >&2
        return 0
    fi

    echo "[bf-triage] cloning ${url} -> ${target_path}" >&2
    git clone --filter=blob:none --no-checkout "${url}" "${target_path}"
}

# Builds an array of REPOS entries that match the requested selection.
# Echoes one entry per line on stdout (caller reads via mapfile).
filter_repos() {
    local select="$1"
    case "${select}" in
        both)
            printf '%s\n' "${REPOS[@]}"
            ;;
        mongo|dsi)
            local entry repo_dir
            for entry in "${REPOS[@]}"; do
                repo_dir="$(echo "${entry}" | awk '{print $2}')"
                if [ "${repo_dir}" = "${select}" ]; then
                    printf '%s\n' "${entry}"
                fi
            done
            ;;
        *)
            echo "[bf-triage] unknown selection: ${select} (expected: both | mongo | dsi)" >&2
            return 2
            ;;
    esac
}

print_export_lines() {
    local base="$1"
    shift
    local entry repo_dir upper
    for entry in "$@"; do
        repo_dir="$(echo "${entry}" | awk '{print $2}')"
        upper="$(echo "${repo_dir}" | tr '[:lower:]' '[:upper:]')"
        echo "export ${upper}_REPO_PATH=${base}/${repo_dir}"
    done
}

cmd_clone() {
    local select="${1:-both}"
    local selected_repos=()
    mapfile -t selected_repos < <(filter_repos "${select}")

    if [ "${#selected_repos[@]}" -eq 0 ]; then
        echo "[bf-triage] selection '${select}' matched no repos" >&2
        return 2
    fi

    local existing
    if [ -f "${STATE_FILE}" ] && [ -d "$(cat "${STATE_FILE}" 2>/dev/null || true)" ]; then
        existing="$(cat "${STATE_FILE}")"
        echo "[bf-triage] reusing existing workdir: ${existing}" >&2
        local entry url repo_dir target
        for entry in "${selected_repos[@]}"; do
            url="$(echo "${entry}" | awk '{print $1}')"
            repo_dir="$(echo "${entry}" | awk '{print $2}')"
            target="${existing}/${repo_dir}"
            clone_one "${url}" "${target}"
        done
        print_export_lines "${existing}" "${selected_repos[@]}"
        return 0
    fi

    ensure_state_dir
    local workdir
    workdir="$(mktemp -d /tmp/bf-triage-workdir-XXXXXXXX)"
    echo "${workdir}" > "${STATE_FILE}"
    echo "[bf-triage] new workdir: ${workdir}" >&2

    local entry url repo_dir target
    for entry in "${selected_repos[@]}"; do
        url="$(echo "${entry}" | awk '{print $1}')"
        repo_dir="$(echo "${entry}" | awk '{print $2}')"
        target="${workdir}/${repo_dir}"
        clone_one "${url}" "${target}"
    done

    print_export_lines "${workdir}" "${selected_repos[@]}"
}

cmd_cleanup_one() {
    if [ ! -f "${STATE_FILE}" ]; then
        echo "[bf-triage] no recorded workdir; nothing to clean up" >&2
        return 0
    fi
    local workdir
    workdir="$(cat "${STATE_FILE}")"
    if [ -d "${workdir}" ]; then
        echo "[bf-triage] removing ${workdir}" >&2
        rm -r -- "${workdir}"
    fi
    rm -- "${STATE_FILE}"
}

cmd_cleanup_all() {
    shopt -s nullglob
    local removed=0
    for d in /tmp/bf-triage-workdir-*; do
        if [ -d "${d}" ]; then
            echo "[bf-triage] removing ${d}" >&2
            rm -r -- "${d}"
            removed=$((removed + 1))
        fi
    done
    if [ -f "${STATE_FILE}" ]; then rm -- "${STATE_FILE}"; fi
    echo "[bf-triage] removed ${removed} workdir(s)" >&2
}

main() {
    case "${1:-clone}" in
        -h|--help|help)
            usage
            ;;
        cleanup)
            if [ "${2:-}" = "all" ]; then
                cmd_cleanup_all
            else
                cmd_cleanup_one
            fi
            ;;
        clone|""|both)
            cmd_clone both
            ;;
        mongo)
            cmd_clone mongo
            ;;
        dsi)
            cmd_clone dsi
            ;;
        *)
            echo "[bf-triage] unknown subcommand: $1" >&2
            usage
            exit 2
            ;;
    esac
}

main "$@"
