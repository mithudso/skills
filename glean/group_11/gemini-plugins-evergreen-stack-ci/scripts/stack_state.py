#!/usr/bin/env python3
"""State manager for the evergreen-stack-ci skill.

Persists branch <-> patch mapping for a Graphite stack (or single branch) to a
JSON file outside the project directory so it survives across worktrees, repo
moves, and Claude Code sessions.

State file location: $XDG_STATE_HOME/evergreen-stack-ci/<repo>--<stack-root>.json
                     (default: ~/.local/state/evergreen-stack-ci/...)

Usage:
    stack_state.py init        --stack-root R --branches a,b,c [--mode stack|single] \\
                               --project-id P --repo-root /path [--trunk master] \\
                               [--profile NAME] [--alias A] [--variants v1,v2] \\
                               [--tasks t1,t2] [--exclude e1,e2]
    stack_state.py path        --stack-root R [--repo-root /path]
    stack_state.py show        --stack-root R [--repo-root /path]
    stack_state.py list        # all known state files
    stack_state.py add-patch   --stack-root R --branch B --patch-id ID --url URL \\
                               [--description D]
    stack_state.py update-status --stack-root R --patch-id ID --status STATUS \\
                                 [--failed-tasks t1,t2]
    stack_state.py set-findings --stack-root R --patch-id ID --notes "..." \\
                                [--cause STRING] [--suspect-branch B] [--verdict V]
    stack_state.py record-failure --stack-root R --patch-id ID \\
                                  --branch SUSPECT --task T --test TEST
    stack_state.py record-success --stack-root R --branch SUSPECT --task T --test TEST
    stack_state.py record-master-broken --stack-root R --branch SUSPECT \\
                                        --task T --test TEST [--evidence STR]
    stack_state.py quarantine  --stack-root R       # list quarantined tests (>=3 consecutive failures)
    stack_state.py summary     --stack-root R [--repo-root /path]   # coordinator-friendly rollup
    stack_state.py reset-poll-cycle --stack-root R   # zero polling counter at start of a polling cycle
    stack_state.py bump-poll-iteration --stack-root R  # increment polling counter; prints iteration: <i>/<max>
    stack_state.py schedule-next-poll --stack-root R --in-seconds N  # record planned next-wakeup time
    stack_state.py mark-build-failure --stack-root R --patch-id ID   # flag patch as RUN_ALL_UNIT_JAVA_TESTS-mediated build failure
    stack_state.py get-fail-fast-aborts --stack-root R   # print descendant patch_ids to abort (one per line)
    stack_state.py mark-fail-fast-aborted --stack-root R # set aborts_dispatched=true (idempotent)
    stack_state.py record-fix  --stack-root R --branch B --commit-sha SHA --summary "..." \\
                               [--target-tests k1,k2]
    stack_state.py mark-completed-stack --stack-root R   # archive when done
    stack_state.py rm          --stack-root R            # delete state file

All writes go through this script (never hand-edit the JSON).
"""

from __future__ import annotations

import argparse
import contextlib
import fcntl
import html
import json
import os
import re
import sys
import webbrowser
from datetime import datetime, timezone
from pathlib import Path
from string import Template

SCHEMA_VERSION = 2
TERMINAL_STATUSES = {"succeeded", "failed", "aborted"}
VALID_STATUSES = {"pending", "started", "succeeded", "failed", "aborted"}
VALID_VERDICTS = {"real-bug", "flake", "master-broken", "needs-retry", "unknown"}
QUARANTINE_THRESHOLD = 3
DASHBOARD_REFRESH_SECONDS = 10
MAX_POLL_ITERATIONS = 12
EARLY_TRIAGE_TASK_THRESHOLD = 10

# Compile/build task names whose failure invalidates the rest of a patch.
# Match is case-insensitive substring so per-team variants
# (e.g. "compile_bazel_payments") still trigger fail-fast.
COMPILE_TASK_TOKENS = ("compile_client_bazel", "compile_bazel")


def _is_compile_task(name: str) -> bool:
    n = (name or "").lower()
    return any(tok in n for tok in COMPILE_TASK_TOKENS)


def _compile_task_kind(name: str) -> str | None:
    """Map a failed-task name to its fail-fast kind.

    More specific token first ("compile_client_bazel") so per-team variants
    don't get mis-classified.
    """
    n = (name or "").lower()
    if "compile_client_bazel" in n:
        return "COMPILE_CLIENT_BAZEL"
    if "compile_bazel" in n:
        return "COMPILE_BAZEL"
    return None


def state_dir() -> Path:
    base = os.environ.get("XDG_STATE_HOME") or str(Path.home() / ".local" / "state")
    return Path(base) / "evergreen-stack-ci"


def slugify(value: str) -> str:
    return re.sub(r"[^A-Za-z0-9._-]+", "-", value).strip("-")


def state_path(stack_root: str, repo_root: str | None = None) -> Path:
    repo_name = Path(repo_root).name if repo_root else "repo"
    return state_dir() / f"{slugify(repo_name)}--{slugify(stack_root)}.json"


def dashboard_path_for(state_file: Path) -> Path:
    return state_file.with_name(state_file.stem + ".dashboard.html")


def now_iso() -> str:
    return datetime.now(timezone.utc).isoformat(timespec="seconds").replace("+00:00", "Z")


def load(path: Path) -> dict:
    with path.open() as f:
        return json.load(f)


def save(path: Path, data: dict) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    tmp = path.with_suffix(path.suffix + ".tmp")
    with tmp.open("w") as f:
        json.dump(data, f, indent=2, sort_keys=False)
        f.write("\n")
    tmp.replace(path)


@contextlib.contextmanager
def _locked(path: Path):
    """Hold an exclusive flock on a sidecar .lock file across a load → mutate → save block.

    Parallel subagents (e.g. multiple `add-patch` calls during parallel patch creation, or
    investigation `record-failure` calls) would otherwise race on read-modify-write and lose
    updates. The lock lives on a separate file so the atomic rename in save() doesn't drop it.
    """
    path.parent.mkdir(parents=True, exist_ok=True)
    lock_path = path.with_suffix(path.suffix + ".lock")
    with lock_path.open("w") as fp:
        fcntl.flock(fp.fileno(), fcntl.LOCK_EX)
        try:
            yield
        finally:
            fcntl.flock(fp.fileno(), fcntl.LOCK_UN)


def split_csv(value: str | None) -> list[str]:
    if not value:
        return []
    return [v.strip() for v in value.split(",") if v.strip()]


# ---------------------------------------------------------------------------
# Dashboard rendering
# ---------------------------------------------------------------------------
#
# A self-contained HTML snapshot of the state file, regenerated after every
# mutating command and viewable in any browser. Uses meta-refresh polling
# (no server, no JS) so the user can open it once and watch progress without
# extra plumbing. Best-effort: a render failure must never break the calling
# mutation, since the state file is the source of truth.

_DASHBOARD_STATUS_BADGE = {
    "succeeded": ("OK",      "badge-ok"),
    "failed":    ("RED",     "badge-red"),
    "started":   ("running", "badge-run"),
    "pending":   ("pending", "badge-pend"),
    "aborted":   ("aborted", "badge-abt"),
}

_DASHBOARD_DECISION_BADGE = {
    "actionable_failure": "badge-red",
    "in_progress":        "badge-run",
    "excluded_only":      "badge-amber",
    "all_clean":          "badge-ok",
    "needs_attention":    "badge-amber",
}

_DASHBOARD_TEMPLATE = Template("""<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<meta http-equiv="refresh" content="$refresh">
<meta name="color-scheme" content="dark">
<title>$title</title>
<style>
  :root {
    /* MongoDB design system (DESIGN.md) — dark surface */
    --bg: #001e2b;              /* canvas-dark / brand-teal-deep */
    --fg: #ffffff;              /* on-dark */
    --muted: #a8b3bc;           /* on-dark-muted */
    --border: #1c2d38;          /* hairline-dark / charcoal */
    --border-strong: #3d4f5b;   /* slate */
    --card-bg: #003d4f;         /* brand-teal — elevated surface */
    --code-bg: #1c2d38;         /* charcoal */
    --th-bg: rgba(255,255,255,0.04);
    --row-tint: rgba(255,255,255,0.03);
    --link-fg: #00ed64;         /* brand-green — high-contrast link on dark */
    --band-bg: #001018;         /* deeper-than-canvas band for separation */
    --band-fg: #ffffff;         /* on-dark */
    --band-muted: #a8b3bc;      /* on-dark-muted */
    --brand-green: #00ed64;
    --brand-green-soft: #c3f0d2;
    /* status badges — pill, accent-encoded per DESIGN.md category palette.
       On dark, the green-success badge uses bright brand-green on deep-teal
       ink (matches the hero CTA recipe); other accents already contrast. */
    --ok-bg: #00ed64;             --ok-fg: #001e2b;  /* brand-green / on-primary */
    --red-bg: #f06bb8;            --red-fg: #ffffff; /* accent-pink */
    --amber-bg: #fa6e39;          --amber-fg: #ffffff; /* accent-orange */
    --run-bg: #3d4f9f;            --run-fg: #ffffff; /* accent-blue */
    --pend-bg: #1c2d38;           --pend-fg: #a8b3bc; /* charcoal / on-dark-muted */
    --abt-bg: #7b3ff2;            --abt-fg: #ffffff; /* accent-purple */
    /* row tints — bumped alpha so they remain visible on the dark canvas */
    --row-failed:    rgba(240,107,184,0.14);
    --row-succeeded: rgba(0,237,100,0.10);
    --row-warning:   rgba(250,110,57,0.12);
    /* verdict text */
    --v-red-fg:   #f06bb8;
    --v-amber-fg: #fa6e39;
    --v-run-fg:   #8a9bff;        /* lighter accent-blue for legibility on dark */
  }
  * { box-sizing: border-box; }
  body {
    font: 14px/1.55 "Euclid Circular A", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, system-ui, sans-serif;
    margin: 0; color: var(--fg); background: var(--bg);
  }
  .band { background: var(--band-bg); color: var(--band-fg); padding: 0; border-bottom: 1px solid var(--band-bg); }
  .band .wrap { padding: 20px 24px; max-width: 1200px; margin: 0 auto; }
  .band .header { display: flex; justify-content: space-between; align-items: baseline; gap: 16px; }
  .band h1 { color: var(--band-fg); font-size: 28px; font-weight: 500; line-height: 1.30; letter-spacing: 0; margin: 0 0 6px; }
  .band code { background: rgba(255,255,255,0.08); color: var(--band-fg); }
  .band .muted, .band .meta { color: var(--band-muted); }
  .band .meta { font-size: 13px; text-align: right; }
  .band .meta-line span + span { margin-left: 12px; }
  .wrap { max-width: 1200px; margin: 0 auto; padding: 24px; }
  h1 { font-size: 28px; font-weight: 500; line-height: 1.30; margin: 0 0 12px; color: var(--fg); }
  h2 { font-size: 18px; font-weight: 600; line-height: 1.40; margin: 32px 0 12px; color: var(--fg); }
  code { font-family: "Source Code Pro", ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-size: 12.5px; background: var(--code-bg); color: var(--fg); padding: 1px 6px; border-radius: 6px; }
  .muted { color: var(--muted); }
  .scope { color: var(--band-muted); font-size: 13px; margin: 8px 0 0; }
  .cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(140px, 1fr)); gap: 12px; margin-bottom: 16px; }
  .card { background: var(--card-bg); border: 1px solid var(--border); border-radius: 12px; padding: 14px 16px; }
  .card .num { font-size: 22px; font-weight: 600; line-height: 1.20; }
  .card .label { color: var(--muted); font-size: 11px; font-weight: 600; text-transform: uppercase; letter-spacing: 1px; margin-top: 4px; }
  .badge { display: inline-block; padding: 3px 10px; border-radius: 9999px; font-size: 12px; font-weight: 600; line-height: 1.30; }
  .badge-ok    { background: var(--ok-bg);    color: var(--ok-fg); }
  .badge-red   { background: var(--red-bg);   color: var(--red-fg); }
  .badge-amber { background: var(--amber-bg); color: var(--amber-fg); }
  .badge-run   { background: var(--run-bg);   color: var(--run-fg); }
  .badge-pend  { background: var(--pend-bg);  color: var(--pend-fg); border: 1px solid var(--border); }
  .badge-abt   { background: var(--abt-bg);   color: var(--abt-fg); }
  .badge-none  { background: var(--code-bg);  color: var(--muted); border: 1px solid var(--border); }
  table { width: 100%; border-collapse: separate; border-spacing: 0; background: var(--card-bg); border: 1px solid var(--border); border-radius: 12px; overflow: hidden; }
  th, td { text-align: left; padding: 10px 14px; border-bottom: 1px solid var(--border); vertical-align: top; }
  th { background: var(--th-bg); font-size: 11px; font-weight: 600; text-transform: uppercase; letter-spacing: 1px; color: var(--muted); }
  tr:last-child td { border-bottom: none; }
  td.branch { font-family: "Source Code Pro", ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-size: 13px; }
  td.num { text-align: right; font-variant-numeric: tabular-nums; }
  .hist { color: var(--muted); font-size: 12px; }
  .row-failed td { background: var(--row-failed); }
  .row-succeeded td { background: var(--row-succeeded); }
  .tasks { font-family: "Source Code Pro", ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-size: 12px; }
  .counts { color: var(--muted); font-size: 12px; margin-top: 2px; }
  .verdict { font-size: 12px; }
  .v-real-bug { color: var(--v-red-fg); font-weight: 600; }
  .v-flake, .v-master-broken { color: var(--v-amber-fg); }
  .v-needs-retry { color: var(--v-run-fg); }
  .v-unknown { color: var(--muted); }
  details.finding { background: var(--card-bg); border: 1px solid var(--border); border-radius: 12px; margin-bottom: 10px; padding: 12px 16px; }
  details.finding summary { cursor: pointer; }
  details.finding p { margin: 6px 0; }
  details.finding .finding-field { margin: 8px 0; }
  details.finding .finding-body { white-space: pre-wrap; margin: 4px 0 0; }
  .earliest { margin: 0 0 16px; color: var(--muted); font-size: 13px; }
  tr.quar td, tr.mb td { background: var(--row-warning); }
  footer { color: var(--muted); font-size: 12px; margin-top: 32px; text-align: right; }
  a { color: var(--link-fg); text-decoration: none; font-weight: 500; }
  a:hover { text-decoration: underline; }
  ::selection { background: var(--brand-green-soft); color: var(--fg); }
</style>
</head>
<body>
<div class="band">
  <div class="wrap">
    <div class="header">
      <div>
        <h1>stack-ci · $stack_root</h1>
        <div class="meta-line"><span class="muted">repo:</span> <code>$repo</code> <span class="muted">mode:</span> <code>$mode</code> <span class="muted">decision:</span> <span class="badge $decision_class">$decision</span> $polling_pill</div>
      </div>
      <div class="meta">generated $generated_at<br>auto-refresh ${refresh}s</div>
    </div>
    <div class="scope">scope: profile=<code>$profile</code> · variants=<code>$variants</code> · tasks=<code>$tasks</code> · excluded=<code>$excluded</code></div>
  </div>
</div>

<div class="wrap">

<div class="cards">
  <div class="card"><div class="num">$n_branches</div><div class="label">branches</div></div>
  <div class="card"><div class="num">$n_patches</div><div class="label">total patches</div></div>
  <div class="card"><div class="num">$n_succeeded</div><div class="label">succeeded</div></div>
  <div class="card"><div class="num">$n_failed_actionable</div><div class="label">failed (actionable)</div></div>
  <div class="card"><div class="num">$n_failed_excluded</div><div class="label">failed (excluded)</div></div>
  <div class="card"><div class="num">$n_running</div><div class="label">running</div></div>
  <div class="card"><div class="num">$n_aborted</div><div class="label">aborted</div></div>
  <div class="card"><div class="num">$n_no_patch</div><div class="label">no patch</div></div>
  <div class="card" id="next-check-card" data-next-wakeup-at="$next_wakeup_at_iso"><div class="num" id="next-check-num">$next_check_in</div><div class="label">next check in</div></div>
  <div class="card"><div class="num">$n_fixes</div><div class="label">fixes applied</div></div>
  <div class="card"><div class="num">$n_fixed</div><div class="label">tests fixed</div></div>
</div>

$earliest_line

<h2>Branches</h2>
<table>
  <thead><tr><th>#</th><th>Branch</th><th>Status</th><th>Latest patch</th><th>Failed tasks / tests</th><th>Verdict</th><th>Suspect</th></tr></thead>
  <tbody>$branch_rows</tbody>
</table>

<h2>Findings</h2>
$findings_section

<h2>Recently fixed tests <span class="muted" style="font-weight:400;font-size:13px;">(count=$fixed_count)</span></h2>
$fixed_section

<h2>Fixes applied <span class="muted" style="font-weight:400;font-size:13px;">(count=$n_fixes)</span></h2>
$fixes_section

<h2>Test failures <span class="muted" style="font-weight:400;font-size:13px;">(tracked=$tracked_count, quarantined=$quarantined_count, master-broken=$master_broken_count, fixed=$fixed_count)</span></h2>
<table>
  <thead><tr><th>Branch</th><th>Task</th><th>Test</th><th class="num">Consecutive</th><th>Flags</th><th>Last patch</th></tr></thead>
  <tbody>$failures_section</tbody>
</table>

$thirdparty_notice

<footer>State file: <code>$state_file</code></footer>

</div>
<script>
// Live countdown for "next check in". The Python render bakes in a snapshot
// of the duration at the time the dashboard was generated, but state-file
// regenerations happen only on mutating commands (~every 5 min during a
// polling cycle), so without this script the number would freeze between
// mutations. We read the absolute next_wakeup_at timestamp from the DOM
// and update the value in place once per second. Mirrors _format_duration
// and the pill formatting in stack_state.py so JS and Python stay in sync.
(function () {
  try {
    var card = document.getElementById("next-check-card");
    var iso = card && card.getAttribute("data-next-wakeup-at");
    if (!iso) return;
    var target = Date.parse(iso);
    if (isNaN(target)) return;
    var numEl = document.getElementById("next-check-num");
    var suffixEl = document.getElementById("polling-suffix");
    function fmt(s) {
      s = Math.max(0, Math.round(s));
      if (s >= 3600) {
        var h = Math.floor(s / 3600);
        var rem = s - h * 3600;
        var m = Math.floor(rem / 60);
        return h + "h " + m + "m " + (rem - m * 60) + "s";
      }
      if (s >= 60) {
        var m2 = Math.floor(s / 60);
        return m2 + "m " + (s - m2 * 60) + "s";
      }
      return s + "s";
    }
    function update() {
      var delta = (target - Date.now()) / 1000;
      var dur, cardText, suffixText;
      if (delta >= 0) {
        dur = fmt(delta);
        cardText = dur;
        suffixText = "next ~" + iso + " (in " + dur + ")";
      } else {
        dur = fmt(-delta);
        cardText = "overdue " + dur;
        suffixText = "overdue (" + dur + ")";
      }
      if (numEl) numEl.textContent = cardText;
      if (suffixEl) suffixEl.textContent = suffixText;
    }
    update();
    setInterval(update, 1000);
  } catch (e) { /* never let dashboard JS break the page */ }
})();
</script>
</body>
</html>
""")


def _esc(value) -> str:
    if value is None:
        return ""
    return html.escape(str(value), quote=True)


_BACKTICK_RE = re.compile(r"`([^`\n]+)`")


def _render_notes(value) -> str:
    """Render an LLM-produced multi-line notes/cause field as readable HTML.

    The investigation subagent emits prose with embedded bullets, numbered
    plan steps, and indented code snippets. Rather than parse as markdown,
    preserve newlines/indentation (via white-space: pre-wrap on the wrapper)
    and convert `backtick` spans to <code>.
    """
    if value is None or value == "":
        return "—"
    escaped = html.escape(str(value), quote=True)
    return _BACKTICK_RE.sub(r"<code>\1</code>", escaped)


def _parse_iso(value) -> datetime | None:
    if not value:
        return None
    try:
        return datetime.fromisoformat(str(value).replace("Z", "+00:00"))
    except ValueError:
        return None


def _format_duration(seconds: float | int | None) -> str:
    """Render a non-negative duration as 'Xh Ym Zs' / 'Xm Ys' / 'Xs'."""
    if seconds is None:
        return "—"
    s = int(round(seconds))
    if s < 0:
        s = 0
    if s >= 3600:
        h, rem = divmod(s, 3600)
        m, sec = divmod(rem, 60)
        return f"{h}h {m}m {sec}s"
    if s >= 60:
        m, sec = divmod(s, 60)
        return f"{m}m {sec}s"
    return f"{s}s"


def _render_dashboard(data: dict, state_file: Path) -> str:
    branches = data.get("branches", [])
    failures = data.get("test_failures") or {}
    scope = data.get("scope") or {}

    counts = {"none": 0, "pending": 0, "started": 0, "succeeded": 0, "failed": 0, "aborted": 0}
    actionable_failed = 0
    excluded_only_failed = 0
    earliest_failed: str | None = None
    earliest_actionable: str | None = None

    branch_rows: list[str] = []
    findings_panels: list[str] = []
    for b in branches:
        order = b.get("order", "")
        name = b.get("name", "")
        patches = b.get("patches", [])
        latest = _latest_patch(b)
        n_patches = len(patches)
        if latest is None:
            counts["none"] += 1
            branch_rows.append(
                f'<tr class="row-no-patch"><td>{order}</td><td class="branch">{_esc(name)}</td>'
                f'<td><span class="badge badge-none">no patch</span></td><td colspan="4" class="muted">—</td></tr>'
            )
            continue

        status = latest.get("status", "pending")
        counts[status] = counts.get(status, 0) + 1
        label, css = _DASHBOARD_STATUS_BADGE.get(status, ("?", "badge-none"))

        actionability = _patch_actionability(latest, failures)
        if status == "failed":
            if earliest_failed is None:
                earliest_failed = name
            if actionability == "actionable":
                actionable_failed += 1
                if earliest_actionable is None:
                    earliest_actionable = name
            elif actionability == "excluded":
                excluded_only_failed += 1

        patch_id = latest.get("patch_id", "")
        url = latest.get("url") or ""
        patch_cell = (
            f'<a href="{_esc(url)}" target="_blank" rel="noopener">{_esc(patch_id)}</a>'
            if url else _esc(patch_id)
        )

        failed_tasks = latest.get("failed_tasks") or []
        failed_tests = latest.get("failed_tests") or []
        tasks_cell = ""
        if status == "failed":
            preview = ", ".join(failed_tasks[:3])
            if len(failed_tasks) > 3:
                preview += f", +{len(failed_tasks) - 3} more"
            action_badge = ""
            if actionability == "actionable":
                action_badge = ' <span class="badge badge-red">actionable</span>'
            elif actionability == "excluded":
                action_badge = ' <span class="badge badge-amber">excluded</span>'
            tasks_cell = (
                f'<div class="tasks">{_esc(preview) or "&mdash;"}</div>'
                f'<div class="counts">{len(failed_tasks)} task(s), {len(failed_tests)} test(s){action_badge}</div>'
            )

        verdict = ""
        suspect = ""
        if latest.get("findings"):
            verdict = (latest["findings"].get("verdict") or "")
            suspect = (latest["findings"].get("suspect_branch") or "")

        verdict_cell = f'<span class="verdict v-{_esc(verdict)}">{_esc(verdict)}</span>' if verdict else ""

        history = f' <span class="hist">(#{n_patches})</span>' if n_patches > 1 else ""

        branch_rows.append(
            f'<tr class="row-{status}"><td>{order}</td>'
            f'<td class="branch">{_esc(name)}{history}</td>'
            f'<td><span class="badge {css}">{label}</span></td>'
            f'<td>{patch_cell}</td>'
            f'<td>{tasks_cell or "<span class=\'muted\'>&mdash;</span>"}</td>'
            f'<td>{verdict_cell}</td>'
            f'<td class="branch">{_esc(suspect)}</td></tr>'
        )

        if latest.get("findings"):
            f = latest["findings"]
            findings_panels.append(
                f'<details class="finding" open>'
                f'<summary><strong>{_esc(name)}</strong> — patch <code>{_esc(patch_id)}</code> · verdict <em>{_esc(f.get("verdict") or "")}</em></summary>'
                f'<div class="finding-field"><strong>Cause:</strong>'
                f'<div class="finding-body">{_render_notes(f.get("cause"))}</div></div>'
                f'<p><strong>Suspect branch:</strong> <code>{_esc(f.get("suspect_branch") or "—")}</code></p>'
                f'<div class="finding-field"><strong>Notes:</strong>'
                f'<div class="finding-body">{_render_notes(f.get("notes"))}</div></div>'
                f'<p class="muted">recorded {_esc(f.get("recorded_at") or "")}</p>'
                f'</details>'
            )

    running = counts["pending"] + counts["started"]
    if actionable_failed > 0:
        decision = "actionable_failure"
    elif running > 0:
        decision = "in_progress"
    elif excluded_only_failed > 0:
        decision = "excluded_only"
    elif counts["succeeded"] > 0 and counts["failed"] == 0 and counts["none"] == 0 and counts["aborted"] == 0:
        decision = "all_clean"
    else:
        decision = "needs_attention"

    # Test failures section
    failure_rows: list[str] = []
    for key, e in sorted(failures.items()):
        flags = []
        row_class = []
        if e.get("quarantined"):
            flags.append('<span class="badge badge-amber">quarantined</span>')
            row_class.append("quar")
        if e.get("master_broken"):
            flags.append('<span class="badge badge-amber">master-broken</span>')
            row_class.append("mb")
        if e.get("fixed_at") and (e.get("consecutive_failures") or 0) == 0:
            flags.append('<span class="badge badge-ok">fixed</span>')
        elif e.get("fixed_at") and (e.get("consecutive_failures") or 0) > 0:
            flags.append('<span class="badge badge-amber">regressed</span>')
            row_class.append("mb")
        cls_attr = f' class="{" ".join(row_class)}"' if row_class else ""
        failure_rows.append(
            f'<tr{cls_attr}>'
            f'<td class="branch">{_esc(e.get("branch"))}</td>'
            f'<td class="branch">{_esc(e.get("task"))}</td>'
            f'<td class="branch">{_esc(e.get("test"))}</td>'
            f'<td class="num">{e.get("consecutive_failures", 0)}</td>'
            f'<td>{" ".join(flags)}</td>'
            f'<td><code>{_esc(e.get("last_failed_patch") or "")}</code></td>'
            f'</tr>'
        )

    # Recently fixed tests section
    fixed_entries = [e for e in failures.values() if e.get("fixed_at")]
    fixed_entries.sort(key=lambda e: e.get("fixed_at") or "", reverse=True)
    fixed_count = len(fixed_entries)
    if fixed_entries:
        fixed_rows: list[str] = []
        for e in fixed_entries[:20]:
            ttf = _format_duration(e.get("time_to_fix_seconds"))
            regressed = (e.get("consecutive_failures") or 0) > 0
            extra = ' <span class="badge badge-amber">regressed</span>' if regressed else ""
            fixed_rows.append(
                '<tr>'
                f'<td class="branch">{_esc(e.get("branch"))}</td>'
                f'<td class="branch">{_esc(e.get("task"))}</td>'
                f'<td class="branch">{_esc(e.get("test"))}{extra}</td>'
                f'<td class="num">{_esc(ttf)}</td>'
                f'<td><code>{_esc(e.get("fixed_in_patch") or "")}</code></td>'
                f'<td class="muted">{_esc(e.get("fixed_at") or "")}</td>'
                '</tr>'
            )
        more = ""
        if len(fixed_entries) > 20:
            more = f'<tr><td colspan="6" class="muted">+{len(fixed_entries) - 20} earlier</td></tr>'
        fixed_section = (
            '<table>'
            '<thead><tr><th>Branch</th><th>Task</th><th>Test</th>'
            '<th class="num">Time to fix</th><th>Fixed in patch</th><th>Fixed at</th></tr></thead>'
            f'<tbody>{"".join(fixed_rows)}{more}</tbody></table>'
        )
    else:
        fixed_section = '<p class="muted">No tests have been fixed yet.</p>'

    # Fixes applied section (commits applied across all branches)
    all_fixes: list[tuple[str, dict]] = []
    for b in branches:
        for fx in (b.get("fixes") or []):
            all_fixes.append((b.get("name", ""), fx))
    all_fixes.sort(key=lambda pair: pair[1].get("applied_at") or "", reverse=True)
    n_fixes = len(all_fixes)
    if all_fixes:
        fix_rows: list[str] = []
        for branch_name, fx in all_fixes[:25]:
            sha = fx.get("commit_sha") or ""
            sha_short = sha[:8] if sha else ""
            targets = fx.get("target_keys") or []
            targets_cell = (
                ", ".join(_esc(t) for t in targets[:3])
                + (f' <span class="muted">+{len(targets) - 3}</span>' if len(targets) > 3 else "")
            ) if targets else '<span class="muted">—</span>'
            fix_rows.append(
                '<tr>'
                f'<td class="branch">{_esc(branch_name)}</td>'
                f'<td class="branch"><code>{_esc(sha_short)}</code></td>'
                f'<td>{_esc(fx.get("summary") or "")}</td>'
                f'<td class="muted">{_esc(fx.get("applied_at") or "")}</td>'
                f'<td class="tasks">{targets_cell}</td>'
                '</tr>'
            )
        more = ""
        if n_fixes > 25:
            more = f'<tr><td colspan="5" class="muted">+{n_fixes - 25} earlier</td></tr>'
        fixes_section = (
            '<table>'
            '<thead><tr><th>Branch</th><th>Commit</th><th>Summary</th><th>Applied</th><th>Targeted tests</th></tr></thead>'
            f'<tbody>{"".join(fix_rows)}{more}</tbody></table>'
        )
    else:
        fixes_section = '<p class="muted">No fix commits recorded yet.</p>'

    # Polling pill + next-check-in card. The pill body is split into a stable
    # prefix (`polling: i/max · `) and a dynamic suffix (`next ~iso (in dur)` /
    # `overdue (dur)` / `paused`) so the inline <script> at the bottom of the
    # template can update the suffix once per second from the data attribute.
    # next_check_in renders the same way: the Python value is the initial paint
    # (works without JS); JS upgrades it to a live counter.
    polling = data.get("polling") or {}
    iter_count = polling.get("iteration_count") or 0
    max_iter = polling.get("max_iterations") or MAX_POLL_ITERATIONS
    next_at_iso = polling.get("next_wakeup_at")
    next_at_dt = _parse_iso(next_at_iso)
    now_dt = datetime.now(timezone.utc)
    pill_prefix = f'polling: {iter_count}/{max_iter} · '
    if next_at_dt is not None:
        delta = (next_at_dt - now_dt).total_seconds()
        if delta >= 0:
            pill_suffix = f'next ~{next_at_iso} (in {_format_duration(delta)})'
            pill_cls = "badge-run"
            next_check_in = _format_duration(delta)
        else:
            pill_suffix = f'overdue ({_format_duration(-delta)})'
            pill_cls = "badge-amber"
            next_check_in = f"overdue {_format_duration(-delta)}"
    elif iter_count > 0 or polling.get("cycle_started_at"):
        pill_suffix = 'paused'
        pill_cls = "badge-pend"
        next_check_in = "—"
    else:
        pill_prefix = ""
        pill_suffix = ""
        pill_cls = ""
        next_check_in = "—"
    if pill_prefix or pill_suffix:
        polling_pill = (
            f'<span class="badge {pill_cls}" data-next-wakeup-at="{_esc(next_at_iso or "")}">'
            f'{_esc(pill_prefix)}<span id="polling-suffix">{_esc(pill_suffix)}</span></span>'
        )
    else:
        polling_pill = ""
    next_wakeup_at_iso = next_at_iso or ""

    earliest_line = ""
    if earliest_actionable:
        earliest_line = (
            f'<div class="earliest">earliest actionable: <code>{_esc(earliest_actionable)}</code></div>'
        )
    elif earliest_failed:
        earliest_line = (
            f'<div class="earliest">earliest red: <code>{_esc(earliest_failed)}</code></div>'
        )

    # Thirdparty banner
    tp_status = scope.get("thirdparty_status") or "omitted"
    tp_teams = scope.get("thirdparty_teams") or []
    tp_reason = scope.get("thirdparty_skipped_reason") or ""
    profile_or_alias = scope.get("profile") or scope.get("alias") or ""
    could_have_tp = profile_or_alias in {"backend", "full"} or "thirdparty" in (scope.get("variants") or [])
    thirdparty_notice = ""
    if tp_status == "skipped-no-mapping" and could_have_tp:
        suffix = f" Reason: <em>{_esc(tp_reason)}</em>" if tp_reason else ""
        thirdparty_notice = (
            '<div class="earliest" style="color:var(--amber-fg);">thirdparty: SKIPPED — could not map stack diff to a team. '
            'Re-run with <code>--thirdparty-teams=&lt;csv&gt;</code> or <code>--thirdparty-teams=all</code>.' + suffix + '</div>'
        )
    elif tp_status == "omitted" and could_have_tp:
        thirdparty_notice = (
            '<div class="earliest">thirdparty: OMITTED — re-run with <code>--thirdparty-teams=auto|all|&lt;csv&gt;</code> to include.</div>'
        )
    elif tp_status in {"auto-resolved", "included"} and tp_teams:
        thirdparty_notice = (
            f'<div class="earliest">thirdparty: included teams=<code>{_esc(",".join(tp_teams))}</code></div>'
        )
    elif tp_status == "all":
        thirdparty_notice = '<div class="earliest">thirdparty: included ALL teams (expensive)</div>'

    profile = scope.get("profile") or scope.get("alias") or "custom"
    excluded = ", ".join(scope.get("excluded") or []) or "none"
    variants = ", ".join(scope.get("variants") or []) or "—"
    tasks = ", ".join(scope.get("tasks") or []) or "—"

    quarantined_count = sum(1 for e in failures.values() if e.get("quarantined"))
    master_broken_count = sum(1 for e in failures.values() if e.get("master_broken"))

    return _DASHBOARD_TEMPLATE.substitute(
        title=_esc(f"stack-ci: {data.get('stack_root') or ''}"),
        stack_root=_esc(data.get("stack_root") or ""),
        repo=_esc(data.get("repo") or ""),
        mode=_esc(data.get("mode") or ""),
        decision=_esc(decision),
        decision_class=_DASHBOARD_DECISION_BADGE.get(decision, "badge-none"),
        polling_pill=polling_pill,
        profile=_esc(profile),
        variants=_esc(variants),
        tasks=_esc(tasks),
        excluded=_esc(excluded),
        n_branches=len(branches),
        n_patches=sum(len(b.get("patches", [])) for b in branches),
        n_succeeded=counts["succeeded"],
        n_failed_actionable=actionable_failed,
        n_failed_excluded=excluded_only_failed,
        n_running=running,
        n_aborted=counts["aborted"],
        n_no_patch=counts["none"],
        next_check_in=_esc(next_check_in),
        next_wakeup_at_iso=_esc(next_wakeup_at_iso),
        n_fixes=n_fixes,
        n_fixed=fixed_count,
        tracked_count=len(failures),
        quarantined_count=quarantined_count,
        master_broken_count=master_broken_count,
        fixed_count=fixed_count,
        earliest_line=earliest_line,
        branch_rows="".join(branch_rows) or '<tr><td colspan="7" class="muted">No branches.</td></tr>',
        findings_section="".join(findings_panels) or '<p class="muted">No findings recorded yet.</p>',
        fixes_section=fixes_section,
        fixed_section=fixed_section,
        failures_section="".join(failure_rows) or '<tr><td colspan="6" class="muted">No test failures recorded.</td></tr>',
        thirdparty_notice=thirdparty_notice,
        generated_at=_esc(now_iso()),
        refresh=str(DASHBOARD_REFRESH_SECONDS),
        state_file=_esc(str(state_file)),
    )


def _maybe_render_dashboard(state_file: Path, data: dict) -> None:
    """Best-effort dashboard regeneration. Never raises into the calling mutation.

    Disable via STACK_STATE_NO_DASHBOARD=1 (useful for tests or when the caller
    only wants to write state without paying the render cost).
    """
    if os.environ.get("STACK_STATE_NO_DASHBOARD"):
        return
    try:
        out = dashboard_path_for(state_file)
        out.parent.mkdir(parents=True, exist_ok=True)
        tmp = out.with_suffix(out.suffix + ".tmp")
        with tmp.open("w") as f:
            f.write(_render_dashboard(data, out))
        tmp.replace(out)
    except Exception as exc:  # noqa: BLE001 — best-effort, don't break the mutation
        print(f"warning: dashboard regen failed: {exc}", file=sys.stderr)


def cmd_init(args: argparse.Namespace) -> int:
    branches = split_csv(args.branches)
    if not branches:
        print("error: --branches requires at least one branch", file=sys.stderr)
        return 2
    repo_root = os.path.abspath(args.repo_root)
    path = state_path(args.stack_root, repo_root)
    if path.exists() and not args.force:
        print(f"error: state file already exists: {path}", file=sys.stderr)
        print("  re-run with --force to overwrite, or use add-patch / update-status", file=sys.stderr)
        return 1

    data = {
        "version": SCHEMA_VERSION,
        "repo": Path(repo_root).name,
        "repo_root": repo_root,
        "project_id": args.project_id,
        "mode": args.mode,
        "stack_root": args.stack_root,
        "trunk": args.trunk,
        "created_at": now_iso(),
        "scope": {
            "profile": args.profile,
            "alias": args.alias,
            "variants": split_csv(args.variants),
            "tasks": split_csv(args.tasks),
            "excluded": split_csv(args.exclude),
            "thirdparty_status": args.thirdparty_status or "omitted",
            "thirdparty_teams": split_csv(args.thirdparty_teams),
            "thirdparty_skipped_reason": args.thirdparty_skipped_reason or "",
        },
        "branches": [
            {"name": b, "order": i, "patches": []}
            for i, b in enumerate(branches)
        ],
        "test_failures": {},
        "polling": {
            "iteration_count": 0,
            "cycle_started_at": None,
            "max_iterations": MAX_POLL_ITERATIONS,
            "last_poll_at": None,
            "next_wakeup_at": None,
            "next_wakeup_seconds": None,
        },
    }
    save(path, data)
    _maybe_render_dashboard(path, data)
    print(str(path))
    return 0


def cmd_path(args: argparse.Namespace) -> int:
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    print(str(state_path(args.stack_root, repo_root)))
    return 0


def cmd_show(args: argparse.Namespace) -> int:
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    print(json.dumps(load(path), indent=2))
    return 0


def cmd_list(_args: argparse.Namespace) -> int:
    d = state_dir()
    if not d.exists():
        return 0
    for p in sorted(d.glob("*.json")):
        try:
            data = load(p)
            n_branches = len(data.get("branches", []))
            n_patches = sum(len(b.get("patches", [])) for b in data.get("branches", []))
            print(f"{p.name}\tbranches={n_branches}\tpatches={n_patches}\troot={data.get('stack_root')}")
        except (OSError, json.JSONDecodeError):
            print(f"{p.name}\t<unreadable>")
    return 0


def find_branch(data: dict, branch: str) -> dict | None:
    for b in data.get("branches", []):
        if b["name"] == branch:
            return b
    return None


def cmd_add_patch(args: argparse.Namespace) -> int:
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path} (run init first)", file=sys.stderr)
        return 1
    with _locked(path):
        data = load(path)
        b = find_branch(data, args.branch)
        if b is None:
            print(f"error: branch {args.branch!r} not in state file", file=sys.stderr)
            return 1
        b["patches"].append({
            "patch_id": args.patch_id,
            "url": args.url,
            "description": args.description or "",
            "created_at": now_iso(),
            "status": "pending",
            "checked_at": None,
            "failed_tasks": [],
            "failed_tests": [],
            "findings": None,
        })
        save(path, data)
        _maybe_render_dashboard(path, data)
    print(f"ok: appended patch {args.patch_id} to branch {args.branch}")
    return 0


def cmd_update_status(args: argparse.Namespace) -> int:
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    if args.status not in VALID_STATUSES:
        print(f"error: invalid status {args.status!r} (expected one of {sorted(VALID_STATUSES)})", file=sys.stderr)
        return 2
    with _locked(path):
        data = load(path)
        found = False
        for b in data.get("branches", []):
            for p in b.get("patches", []):
                if p["patch_id"] == args.patch_id:
                    p["status"] = args.status
                    p["checked_at"] = now_iso()
                    if args.failed_tasks is not None:
                        p["failed_tasks"] = split_csv(args.failed_tasks)
                    found = True
                    break
            if found:
                break
        if not found:
            print(f"error: patch {args.patch_id} not found in {path}", file=sys.stderr)
            return 1
        save(path, data)
        _maybe_render_dashboard(path, data)
    print(f"ok: patch {args.patch_id} -> {args.status}")
    return 0


def cmd_set_findings(args: argparse.Namespace) -> int:
    """Subagent persists investigation findings for a patch."""
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    if args.verdict and args.verdict not in VALID_VERDICTS:
        print(f"error: invalid verdict {args.verdict!r} (expected one of {sorted(VALID_VERDICTS)})", file=sys.stderr)
        return 2
    with _locked(path):
        data = load(path)
        found = False
        for b in data.get("branches", []):
            for p in b.get("patches", []):
                if p["patch_id"] == args.patch_id:
                    p["findings"] = {
                        "notes": args.notes,
                        "cause": args.cause,
                        "suspect_branch": args.suspect_branch,
                        "verdict": args.verdict,
                        "recorded_at": now_iso(),
                    }
                    found = True
                    break
            if found:
                break
        if not found:
            print(f"error: patch {args.patch_id} not found in {path}", file=sys.stderr)
            return 1
        save(path, data)
        _maybe_render_dashboard(path, data)
    print(f"ok: findings recorded for patch {args.patch_id}")
    return 0


def _latest_patch(branch: dict) -> dict | None:
    return branch["patches"][-1] if branch.get("patches") else None


def _failure_key(branch: str, task: str, test: str) -> str:
    return f"{branch}::{task}::{test}"


def _find_patch(data: dict, patch_id: str) -> tuple[dict | None, dict | None]:
    """Returns (branch, patch) tuple for the patch_id, or (None, None)."""
    for b in data.get("branches", []):
        for p in b.get("patches", []):
            if p["patch_id"] == patch_id:
                return b, p
    return None, None


def cmd_record_failure(args: argparse.Namespace) -> int:
    """Investigation subagent records one failing test on a patch.

    Has two effects:
      1. Append the test to the patch's failed_tests array (so summary can compute actionability).
      2. Increment the consecutive-failure counter for (branch, task, test). Quarantines at >= 3.

    Dedup: counter only increments once per "round". A round is identified by the latest
    patch on the suspect branch — so cascading failures investigated in parallel across
    multiple child patches don't over-count.
    """
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    with _locked(path):
        data = load(path)

        _, patch = _find_patch(data, args.patch_id)
        if patch is None:
            print(f"error: patch {args.patch_id} not found", file=sys.stderr)
            return 1
        suspect_branch = find_branch(data, args.branch)
        if suspect_branch is None:
            print(f"error: branch {args.branch!r} not in state", file=sys.stderr)
            return 1

        failed_tests = patch.setdefault("failed_tests", [])
        test_entry = {"task": args.task, "test": args.test, "suspect_branch": args.branch}
        if test_entry not in failed_tests:
            failed_tests.append(test_entry)

        latest = _latest_patch(suspect_branch)
        round_id = latest["patch_id"] if latest else args.patch_id

        failures = data.setdefault("test_failures", {})
        key = _failure_key(args.branch, args.task, args.test)
        entry = failures.get(key)
        if entry is None:
            entry = {
                "branch": args.branch,
                "task": args.task,
                "test": args.test,
                "consecutive_failures": 1,
                "first_failed_at": now_iso(),
                "last_failed_at": now_iso(),
                "last_failed_patch": args.patch_id,
                "last_round_id": round_id,
                "quarantined": False,
            }
        elif entry.get("last_round_id") != round_id:
            entry["consecutive_failures"] += 1
            entry["last_round_id"] = round_id
            entry["last_failed_at"] = now_iso()
            entry["last_failed_patch"] = args.patch_id
        if entry["consecutive_failures"] >= QUARANTINE_THRESHOLD:
            entry["quarantined"] = True
        failures[key] = entry

        save(path, data)
        _maybe_render_dashboard(path, data)
    print(f"ok: {key} consecutive={entry['consecutive_failures']} quarantined={entry['quarantined']}")
    return 0


def cmd_record_master_broken(args: argparse.Namespace) -> int:
    """Mark a test as already failing on master (or whatever the trunk is).

    Once flagged, the test is excluded from actionable-failure counts in `summary` —
    it stays visible (so the user knows it's red) but the skill won't try to fix it,
    because the breakage isn't introduced by the stack.

    Investigation subagents should call this whenever they verify master is also red
    on the same task/test (via mcp__evergreen__list_user_recent_patches_evergreen
    against the project's master patches).
    """
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    with _locked(path):
        data = load(path)
        failures = data.setdefault("test_failures", {})
        key = _failure_key(args.branch, args.task, args.test)
        entry = failures.get(key)
        if entry is None:
            entry = {
                "branch": args.branch,
                "task": args.task,
                "test": args.test,
                "consecutive_failures": 0,
                "first_failed_at": now_iso(),
                "last_failed_at": now_iso(),
                "last_failed_patch": args.patch_id,
                "last_round_id": None,
                "quarantined": False,
            }
        entry["master_broken"] = True
        entry["master_broken_evidence"] = args.evidence or ""
        entry["master_broken_recorded_at"] = now_iso()
        failures[key] = entry
        save(path, data)
        _maybe_render_dashboard(path, data)
    print(f"ok: {key} flagged master_broken=true")
    return 0


def cmd_record_success(args: argparse.Namespace) -> int:
    """Reset the consecutive-failure counter for a test that's now passing.

    Investigation/poll subagents call this when they observe a previously-failing test
    pass on a new patch. Resetting un-quarantines the test if it had hit the threshold.

    Also stamps fixed_at, fixed_in_patch, and time_to_fix_seconds so the dashboard can
    surface a "Recently fixed tests" panel. The original first_failed_at is preserved.
    """
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    with _locked(path):
        data = load(path)
        failures = data.get("test_failures") or {}
        key = _failure_key(args.branch, args.task, args.test)
        entry = failures.get(key)
        if entry is None:
            print(f"ok: {key} (no prior failures recorded)")
            return 0
        entry["consecutive_failures"] = 0
        entry["quarantined"] = False
        entry["last_round_id"] = None
        fixed_at = now_iso()
        entry["fixed_at"] = fixed_at
        if args.patch_id:
            entry["fixed_in_patch"] = args.patch_id
        first = entry.get("first_failed_at")
        if first:
            try:
                first_dt = datetime.fromisoformat(first.replace("Z", "+00:00"))
                fixed_dt = datetime.fromisoformat(fixed_at.replace("Z", "+00:00"))
                entry["time_to_fix_seconds"] = int((fixed_dt - first_dt).total_seconds())
            except ValueError:
                pass
        save(path, data)
        _maybe_render_dashboard(path, data)
    print(f"ok: {key} reset (was quarantined? -> false)")
    return 0


def cmd_record_fix(args: argparse.Namespace) -> int:
    """Append a fix-commit entry to a branch.

    The fix-and-commit subagent calls this immediately after `git commit` so the
    dashboard can show what code was applied to fix a failure. Multiple fixes per
    branch (e.g. across cycles) are kept in append order.
    """
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    with _locked(path):
        data = load(path)
        b = find_branch(data, args.branch)
        if b is None:
            print(f"error: branch {args.branch!r} not in state file", file=sys.stderr)
            return 1
        fixes = b.setdefault("fixes", [])
        fixes.append({
            "commit_sha": args.commit_sha,
            "summary": args.summary,
            "applied_at": now_iso(),
            "target_keys": split_csv(args.target_tests),
        })
        save(path, data)
        _maybe_render_dashboard(path, data)
    print(f"ok: recorded fix {args.commit_sha[:8]} on branch {args.branch}")
    return 0


def cmd_quarantine(args: argparse.Namespace) -> int:
    """List quarantined tests (consecutive_failures >= QUARANTINE_THRESHOLD)."""
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    data = load(path)
    failures = data.get("test_failures") or {}
    rows = [e for e in failures.values() if e.get("quarantined")]
    if not rows:
        print("no quarantined tests")
        return 0
    rows.sort(key=lambda e: (e["branch"], e["task"], e["test"]))
    for e in rows:
        print(f"{e['branch']}\t{e['task']}\t{e['test']}\t"
              f"consecutive={e['consecutive_failures']}\tlast_patch={e.get('last_failed_patch')}")
    return 0


def _ensure_polling(data: dict) -> dict:
    """Auto-inject the polling block on legacy state files that pre-date the field.

    Returns the polling block (a reference into `data`) so callers can mutate it directly.
    """
    polling = data.get("polling")
    if not isinstance(polling, dict):
        polling = {
            "iteration_count": 0,
            "cycle_started_at": None,
            "max_iterations": MAX_POLL_ITERATIONS,
            "last_poll_at": None,
            "next_wakeup_at": None,
            "next_wakeup_seconds": None,
        }
        data["polling"] = polling
    polling.setdefault("iteration_count", 0)
    polling.setdefault("cycle_started_at", None)
    polling.setdefault("max_iterations", MAX_POLL_ITERATIONS)
    polling.setdefault("last_poll_at", None)
    polling.setdefault("next_wakeup_at", None)
    polling.setdefault("next_wakeup_seconds", None)
    return polling


def cmd_bump_poll_iteration(args: argparse.Namespace) -> int:
    """Increment the polling iteration counter. Called by the poll-status subagent
    every wakeup so the coordinator can read iteration <i>/<max> from the subagent's
    return line and decide whether to schedule another wakeup or pause for input.
    """
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    with _locked(path):
        data = load(path)
        polling = _ensure_polling(data)
        polling["iteration_count"] = int(polling.get("iteration_count") or 0) + 1
        if not polling.get("cycle_started_at"):
            polling["cycle_started_at"] = now_iso()
        polling["last_poll_at"] = now_iso()
        polling["next_wakeup_at"] = None
        polling["next_wakeup_seconds"] = None
        save(path, data)
        _maybe_render_dashboard(path, data)
    print(f"iteration: {polling['iteration_count']}/{polling['max_iterations']}")
    return 0


def cmd_schedule_next_poll(args: argparse.Namespace) -> int:
    """Record when the next ScheduleWakeup-driven poll is expected to fire.

    The coordinator calls this immediately before invoking ScheduleWakeup so the
    dashboard can display a countdown to the next poll. The field is cleared by the
    next bump-poll-iteration call (the wakeup fired) or by reset-poll-cycle.
    """
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    if args.in_seconds < 0:
        print(f"error: --in-seconds must be >= 0 (got {args.in_seconds})", file=sys.stderr)
        return 2
    with _locked(path):
        data = load(path)
        polling = _ensure_polling(data)
        next_at = datetime.now(timezone.utc).timestamp() + args.in_seconds
        next_iso = (
            datetime.fromtimestamp(next_at, tz=timezone.utc)
            .isoformat(timespec="seconds")
            .replace("+00:00", "Z")
        )
        polling["next_wakeup_at"] = next_iso
        polling["next_wakeup_seconds"] = int(args.in_seconds)
        save(path, data)
        _maybe_render_dashboard(path, data)
    print(f"next-wakeup: {next_iso} (in {int(args.in_seconds)}s)")
    return 0


def cmd_reset_poll_cycle(args: argparse.Namespace) -> int:
    """Reset the polling iteration counter. Coordinator calls this on entry to Phase 2
    (and again on Phase 3 -> Phase 2 hand-back after fixes/re-patches), so each polling
    cycle gets a fresh 12-iteration budget.
    """
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    with _locked(path):
        data = load(path)
        polling = _ensure_polling(data)
        polling["iteration_count"] = 0
        polling["cycle_started_at"] = now_iso()
        polling["last_poll_at"] = None
        polling["next_wakeup_at"] = None
        polling["next_wakeup_seconds"] = None
        # Phase 2 re-entry (after a fix) starts with a fresh fail-fast slate so
        # a stale event from the previous cycle doesn't suppress get-fail-fast-aborts.
        polling.pop("compile_fail_fast", None)
        save(path, data)
        _maybe_render_dashboard(path, data)
    print(f"iteration: 0/{polling['max_iterations']}")
    return 0


def _detect_compile_fail_fast(state: dict) -> dict | None:
    """Walk branches root → tip and return the earliest compile/build-failing patch.

    A branch is a candidate if its latest patch is in {started, failed} (still
    in flight or just finished) AND either:
      - any task name in failed_tasks matches _is_compile_task, OR
      - the patch has build_failure == true (set by mark-build-failure when the
        poll subagent detected the RUN_ALL_UNIT_JAVA_TESTS "No test results
        found" pattern).

    Descendants are later branches in stack order whose latest patch is still
    abortable (not succeeded / aborted / missing).

    Returns None when no candidate is found.
    """
    branches = state.get("branches", [])
    suspect_idx: int | None = None
    suspect: dict | None = None
    for i, b in enumerate(branches):
        latest = _latest_patch(b)
        if latest is None:
            continue
        if latest.get("status") not in {"started", "failed"}:
            continue
        kind: str | None = None
        for t in latest.get("failed_tasks") or []:
            kind = _compile_task_kind(t)
            if kind is not None:
                break
        if kind is None and latest.get("build_failure"):
            kind = "BUILD_FAILURE"
        if kind is None:
            continue
        suspect_idx = i
        suspect = {
            "suspect_branch": b.get("name"),
            "suspect_patch_id": latest.get("patch_id"),
            "kind": kind,
        }
        break
    if suspect is None or suspect_idx is None:
        return None

    descendants: list[dict] = []
    for j in range(suspect_idx + 1, len(branches)):
        b = branches[j]
        latest = _latest_patch(b)
        if latest is None:
            continue
        if latest.get("status") in {"succeeded", "aborted"}:
            continue
        descendants.append({"branch": b.get("name"), "patch_id": latest.get("patch_id")})
    suspect["descendants"] = descendants
    return suspect


def _detect_early_triage(branches: list[dict]) -> str | None:
    """Return the earliest branch (by stack order) whose latest patch is still
    `started` but has already accumulated >= EARLY_TRIAGE_TASK_THRESHOLD failed
    tasks. Returns None if no branch qualifies.
    """
    for b in branches:
        latest = _latest_patch(b)
        if latest is None:
            continue
        if latest.get("status") != "started":
            continue
        if len(latest.get("failed_tasks") or []) >= EARLY_TRIAGE_TASK_THRESHOLD:
            return b.get("name")
    return None


def _persist_compile_fail_fast(data: dict, ff: dict, mark_dispatched: bool = False) -> bool:
    """Persist the compile fail-fast event under polling.compile_fail_fast.

    Idempotent. If a block already exists, leaves the original
    triggered_at/suspect/kind/descendants alone. Optionally flips
    aborts_dispatched to true.

    Returns True if `data` was mutated (caller should save + re-render).
    """
    polling = _ensure_polling(data)
    block = polling.get("compile_fail_fast")
    changed = False
    if not isinstance(block, dict):
        block = {
            "triggered_at": now_iso(),
            "suspect_branch": ff.get("suspect_branch"),
            "suspect_patch_id": ff.get("suspect_patch_id"),
            "kind": ff.get("kind"),
            "descendants": list(ff.get("descendants") or []),
            "aborts_dispatched": False,
        }
        polling["compile_fail_fast"] = block
        changed = True
    if mark_dispatched and not block.get("aborts_dispatched"):
        block["aborts_dispatched"] = True
        block.setdefault("triggered_at", now_iso())
        changed = True
    return changed


def _is_excluded(entry: dict | None) -> bool:
    """A failure entry is excluded from actionability if it's quarantined OR master-broken."""
    if entry is None:
        return False
    return bool(entry.get("quarantined") or entry.get("master_broken"))


def _patch_actionability(patch: dict, failures: dict) -> str:
    """Returns one of:
      'actionable'  - has at least one failed test that is NOT quarantined and NOT master-broken
                      (or has no test-level data — default to actionable so a human looks)
      'excluded'    - every failed test is quarantined or master-broken (skill won't fix)
      'n/a'         - patch is not in failed state
    """
    if patch["status"] != "failed":
        return "n/a"
    failed_tests = patch.get("failed_tests") or []
    if not failed_tests:
        return "actionable"
    for ft in failed_tests:
        suspect = ft.get("suspect_branch") or ft.get("branch")
        key = _failure_key(suspect, ft["task"], ft["test"])
        entry = failures.get(key)
        if not _is_excluded(entry):
            return "actionable"
    return "excluded"


def cmd_summary(args: argparse.Namespace) -> int:
    """Coordinator-friendly terse summary. Designed to fit on a screen.

    Emits a `poll-decision:` line the coordinator consumes to choose its next move:
      actionable_failure  -> stop polling, fix
      in_progress         -> keep polling
      quarantine_only     -> stop polling, surface to user (no fixable failures left)
      all_clean           -> all patches green; offer Phase 4
    """
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    data = load(path)
    branches = data.get("branches", [])
    scope = data.get("scope") or {}
    failures = data.get("test_failures") or {}
    profile = scope.get("profile") or scope.get("alias") or "custom"
    excluded = ",".join(scope.get("excluded") or []) or "none"

    counts = {"none": 0, "pending": 0, "started": 0, "succeeded": 0, "failed": 0, "aborted": 0}
    actionable_failed = 0
    excluded_only_failed = 0
    earliest_failed: str | None = None
    earliest_actionable: str | None = None
    lines: list[str] = []
    for b in branches:
        p = _latest_patch(b)
        if p is None:
            counts["none"] += 1
            lines.append(f"  [{b['order']}] {b['name']:<50} -- no patch")
            continue
        counts[p["status"]] = counts.get(p["status"], 0) + 1
        marker = {
            "succeeded": "OK ",
            "failed":    "RED",
            "pending":   "...",
            "started":   ">> ",
            "aborted":   "ABT",
        }.get(p["status"], "?  ")
        actionability = _patch_actionability(p, failures)
        if p["status"] == "failed":
            if earliest_failed is None:
                earliest_failed = b["name"]
            if actionability == "actionable":
                actionable_failed += 1
                if earliest_actionable is None:
                    earliest_actionable = b["name"]
            elif actionability == "excluded":
                excluded_only_failed += 1
        failed_summary = ""
        if p["status"] == "failed":
            n_tasks = len(p.get("failed_tasks") or [])
            n_tests = len(p.get("failed_tests") or [])
            failed_summary = f" [{n_tasks} task(s), {n_tests} test(s), {actionability}]"
        verdict = ""
        if p.get("findings") and p["findings"].get("verdict"):
            verdict = f" verdict={p['findings']['verdict']}"
        lines.append(f"  [{b['order']}] {marker} {b['name']:<50} {p['patch_id']}{failed_summary}{verdict}")

    running = counts["pending"] + counts["started"]
    extra_lines: list[str] = []
    ff = _detect_compile_fail_fast(data)
    if ff is not None:
        # Fail-fast overrides the normal decision tree: any in-flight compile
        # failure forces actionable_failure even if other patches are still
        # running, because their builds will inherit the same compile error.
        decision = "actionable_failure"
        earliest_actionable = ff["suspect_branch"]
        descendants_csv = ",".join(
            f"{d['branch']}={d['patch_id']}" for d in ff["descendants"]
        ) or "(none)"
        extra_lines.append(
            f"compile-fail-fast: {ff['suspect_branch']} "
            f"kind={ff['kind']} patch={ff['suspect_patch_id']}"
        )
        extra_lines.append(f"compile-fail-descendants: {descendants_csv}")
        # Persist the event the first time it's seen so dashboards + later
        # phases see the same suspect across polls. Idempotent — re-load under
        # lock so we don't clobber a concurrent mutation.
        try:
            with _locked(path):
                fresh = load(path)
                if _persist_compile_fail_fast(fresh, ff):
                    save(path, fresh)
                    _maybe_render_dashboard(path, fresh)
        except OSError as exc:  # noqa: BLE001 — best-effort
            print(f"warning: compile-fail-fast persist failed: {exc}", file=sys.stderr)
    elif actionable_failed > 0:
        decision = "actionable_failure"
    elif (early_triage_branch := _detect_early_triage(branches)) is not None:
        # A still-running patch has already accumulated
        # EARLY_TRIAGE_TASK_THRESHOLD failed tasks. Surface as
        # actionable_failure (user-confirmation handshake) without setting a
        # compile-fail-fast line — the remaining tasks may still pass, so
        # descendants are not aborted automatically.
        decision = "actionable_failure"
        if earliest_actionable is None:
            earliest_actionable = early_triage_branch
    elif running > 0:
        decision = "in_progress"
    elif excluded_only_failed > 0:
        decision = "excluded_only"
    elif counts["succeeded"] > 0 and counts["failed"] == 0 and counts["none"] == 0 and counts["aborted"] == 0:
        decision = "all_clean"
    else:
        decision = "needs_attention"

    quarantined_count = sum(1 for e in failures.values() if e.get("quarantined"))
    master_broken_count = sum(1 for e in failures.values() if e.get("master_broken"))
    tracked_count = len(failures)

    print(f"stack-root:    {data.get('stack_root')}")
    print(f"mode:          {data.get('mode')}")
    print(f"scope:         profile={profile} excluded={excluded}")
    print(f"branches:      {len(branches)}  patches: {sum(len(b.get('patches', [])) for b in branches)}")
    print(f"latest-status: succeeded={counts['succeeded']} failed={counts['failed']} "
          f"started={counts['started']} pending={counts['pending']} aborted={counts['aborted']} no-patch={counts['none']}")
    print(f"failed-split:  actionable={actionable_failed} excluded-only={excluded_only_failed}")
    print(f"test-failures: tracked={tracked_count} quarantined={quarantined_count} master-broken={master_broken_count}")
    print(f"poll-decision: {decision}")
    for line in extra_lines:
        print(line)
    if earliest_failed:
        print(f"earliest-red:  {earliest_failed}")
    if earliest_actionable:
        print(f"earliest-actionable: {earliest_actionable}")

    tp_status = scope.get("thirdparty_status") or "omitted"
    tp_teams = scope.get("thirdparty_teams") or []
    tp_reason = scope.get("thirdparty_skipped_reason") or ""
    profile_or_alias = scope.get("profile") or scope.get("alias") or ""
    could_have_thirdparty = profile_or_alias in {"backend", "full"} or "thirdparty" in (scope.get("variants") or [])
    if tp_status == "skipped-no-mapping" and could_have_thirdparty:
        suffix = f" reason='{tp_reason}'" if tp_reason else ""
        print(f"thirdparty-notice: SKIPPED — could not map stack diff to a thirdparty team. "
              f"Re-run with --thirdparty-teams=<csv> (e.g. payments,billing) or --thirdparty-teams=all to include thirdparty.{suffix}")
    elif tp_status == "omitted" and could_have_thirdparty:
        print(f"thirdparty-notice: OMITTED — thirdparty tests not requested. "
              f"Re-run with --thirdparty-teams=auto|all|<csv> to include them.")
    elif tp_status in {"auto-resolved", "included"} and tp_teams:
        print(f"thirdparty-status: included teams={','.join(tp_teams)}")
    elif tp_status == "all":
        print(f"thirdparty-status: included ALL teams (expensive)")
    print()
    print("\n".join(lines))
    return 0


def cmd_mark_build_failure(args: argparse.Namespace) -> int:
    """Flag a patch as having an upstream build failure.

    Set by the poll-status subagent when a RUN_ALL_UNIT_JAVA_TESTS-mediated
    failure is detected via the "No test results found — there may have been
    a build failure" pattern. _detect_compile_fail_fast picks this flag up
    and surfaces it as kind=BUILD_FAILURE.

    COMPILE_BAZEL / COMPILE_CLIENT_BAZEL failures don't need this command —
    they're detected by name match on failed_tasks.
    """
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    with _locked(path):
        data = load(path)
        _, patch = _find_patch(data, args.patch_id)
        if patch is None:
            print(f"error: patch {args.patch_id} not found", file=sys.stderr)
            return 1
        if patch.get("build_failure"):
            print(f"ok: patch {args.patch_id} build_failure already set")
            return 0
        patch["build_failure"] = True
        save(path, data)
        _maybe_render_dashboard(path, data)
    print(f"ok: patch {args.patch_id} flagged build_failure=true")
    return 0


def cmd_get_fail_fast_aborts(args: argparse.Namespace) -> int:
    """Print descendant patch_ids that need to be aborted (one per line).

    Exit 0 with empty output means either:
      - No compile fail-fast detected, OR
      - aborts already dispatched (compile_fail_fast.aborts_dispatched == true).
    """
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    data = load(path)
    ff = _detect_compile_fail_fast(data)
    if ff is None:
        return 0
    polling = data.get("polling") or {}
    block = polling.get("compile_fail_fast")
    if isinstance(block, dict) and block.get("aborts_dispatched"):
        return 0
    for d in ff["descendants"]:
        pid = d.get("patch_id")
        if pid:
            print(pid)
    return 0


def cmd_mark_fail_fast_aborted(args: argparse.Namespace) -> int:
    """Set compile_fail_fast.aborts_dispatched = true (idempotent).

    First call also persists the full event block (suspect, kind, descendants,
    triggered_at) if cmd_summary hasn't already done so. Subsequent calls just
    re-affirm the dispatched flag.
    """
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    with _locked(path):
        data = load(path)
        ff = _detect_compile_fail_fast(data)
        polling = _ensure_polling(data)
        block = polling.get("compile_fail_fast")
        if ff is None and not isinstance(block, dict):
            print("ok: no compile fail-fast event to mark (nothing to do)")
            return 0
        if ff is not None:
            changed = _persist_compile_fail_fast(data, ff, mark_dispatched=True)
        else:
            # Detection no longer fires (e.g. patches finished) but block exists —
            # just flip aborts_dispatched if not already set.
            changed = False
            if not block.get("aborts_dispatched"):
                block["aborts_dispatched"] = True
                block.setdefault("triggered_at", now_iso())
                changed = True
        if changed:
            save(path, data)
            _maybe_render_dashboard(path, data)
        suspect = (ff or block or {}).get("suspect_branch") or "?"
    print(f"ok: compile fail-fast aborts marked dispatched (suspect={suspect})")
    return 0


def cmd_mark_completed_stack(args: argparse.Namespace) -> int:
    """Archive a state file once all PRs merge / work is done."""
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    archive_dir = state_dir() / "archive"
    archive_dir.mkdir(parents=True, exist_ok=True)
    target = archive_dir / f"{path.stem}--{now_iso().replace(':', '-')}.json"
    path.rename(target)
    dash = dashboard_path_for(path)
    if dash.exists():
        dash.unlink()
    print(f"archived: {target}")
    return 0


def cmd_rm(args: argparse.Namespace) -> int:
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    path.unlink()
    dash = dashboard_path_for(path)
    if dash.exists():
        dash.unlink()
    print(f"removed: {path}")
    return 0


def cmd_dashboard(args: argparse.Namespace) -> int:
    """Force-regenerate the HTML dashboard for a stack root.

    Auto-regen happens on every mutating command, so this is mainly useful for:
      - Recovering after a render failure (e.g. disk full mid-write).
      - Bumping the timestamp without otherwise touching state.
      - Re-rendering after STACK_STATE_NO_DASHBOARD was previously set.
    """
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    path = state_path(args.stack_root, repo_root)
    if not path.exists():
        print(f"error: no state file at {path}", file=sys.stderr)
        return 1
    data = load(path)
    out = dashboard_path_for(path)
    out.parent.mkdir(parents=True, exist_ok=True)
    tmp = out.with_suffix(out.suffix + ".tmp")
    with tmp.open("w") as f:
        f.write(_render_dashboard(data, out))
    tmp.replace(out)
    print(str(out))
    return 0


def cmd_dashboard_path(args: argparse.Namespace) -> int:
    """Print the dashboard HTML path for a stack root (for opening in a browser).

    With --open, also launch the user's default browser to the file:// URL
    (best-effort; opt out via STACK_STATE_NO_OPEN=1).
    """
    repo_root = os.path.abspath(args.repo_root) if args.repo_root else None
    dash = dashboard_path_for(state_path(args.stack_root, repo_root))
    print(str(dash))
    if getattr(args, "open", False) and not os.environ.get("STACK_STATE_NO_OPEN"):
        if not dash.exists():
            print(
                f"warning: dashboard not yet rendered at {dash}; skipping --open",
                file=sys.stderr,
            )
        else:
            try:
                webbrowser.open(f"file://{dash}", new=2)
            except Exception as exc:
                print(f"warning: failed to open dashboard in browser: {exc}", file=sys.stderr)
    return 0


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description="State manager for evergreen-stack-ci")
    sub = parser.add_subparsers(dest="cmd", required=True)

    p_init = sub.add_parser("init", help="Initialize a new state file")
    p_init.add_argument("--stack-root", required=True, help="Branch closest to trunk; uniquely identifies the stack")
    p_init.add_argument("--branches", required=True, help="CSV of branches in stack order (root -> tip)")
    p_init.add_argument("--mode", choices=["stack", "single"], default="stack")
    p_init.add_argument("--project-id", required=True, help="Evergreen project id, e.g. mms")
    p_init.add_argument("--repo-root", required=True, help="Absolute path to repo root")
    p_init.add_argument("--trunk", default="master")
    p_init.add_argument("--profile", default=None, help="Test scope preset (backend, unit, full, ...)")
    p_init.add_argument("--alias", default=None, help="Evergreen alias if used (overrides profile)")
    p_init.add_argument("--variants", default=None, help="CSV of variants")
    p_init.add_argument("--tasks", default=None, help="CSV of tasks")
    p_init.add_argument("--exclude", default=None, help="CSV of excluded categories")
    p_init.add_argument("--thirdparty-status", default=None,
                        choices=["included", "auto-resolved", "all", "skipped-no-mapping", "omitted", "excluded"],
                        help="Thirdparty variant status; defaults to 'omitted'")
    p_init.add_argument("--thirdparty-teams", default=None, help="CSV of resolved thirdparty teams (e.g. payments,billing)")
    p_init.add_argument("--thirdparty-skipped-reason", default=None,
                        help="Free-text reason when thirdparty-status=skipped-no-mapping (drives reminder banner)")
    p_init.add_argument("--force", action="store_true", help="Overwrite existing state file")
    p_init.set_defaults(func=cmd_init)

    p_path = sub.add_parser("path", help="Print state file path for a stack root")
    p_path.add_argument("--stack-root", required=True)
    p_path.add_argument("--repo-root", default=None)
    p_path.set_defaults(func=cmd_path)

    p_show = sub.add_parser("show", help="Pretty-print state JSON")
    p_show.add_argument("--stack-root", required=True)
    p_show.add_argument("--repo-root", default=None)
    p_show.set_defaults(func=cmd_show)

    p_list = sub.add_parser("list", help="List all state files")
    p_list.set_defaults(func=cmd_list)

    p_add = sub.add_parser("add-patch", help="Append a patch entry to a branch")
    p_add.add_argument("--stack-root", required=True)
    p_add.add_argument("--branch", required=True)
    p_add.add_argument("--patch-id", required=True)
    p_add.add_argument("--url", required=True)
    p_add.add_argument("--description", default=None)
    p_add.add_argument("--repo-root", default=None)
    p_add.set_defaults(func=cmd_add_patch)

    p_upd = sub.add_parser("update-status", help="Set a patch's status + failed tasks")
    p_upd.add_argument("--stack-root", required=True)
    p_upd.add_argument("--patch-id", required=True)
    p_upd.add_argument("--status", required=True, help=f"One of: {sorted(VALID_STATUSES)}")
    p_upd.add_argument("--failed-tasks", default=None, help="CSV of failed task ids")
    p_upd.add_argument("--repo-root", default=None)
    p_upd.set_defaults(func=cmd_update_status)

    p_find = sub.add_parser("set-findings", help="Persist investigation notes for a patch")
    p_find.add_argument("--stack-root", required=True)
    p_find.add_argument("--patch-id", required=True)
    p_find.add_argument("--notes", required=True, help="Free-text investigation summary")
    p_find.add_argument("--cause", default=None, help="One-line root cause if known")
    p_find.add_argument("--suspect-branch", default=None, help="Branch suspected to have introduced the failure")
    p_find.add_argument("--verdict", default=None, help=f"One of: {sorted(VALID_VERDICTS)}")
    p_find.add_argument("--repo-root", default=None)
    p_find.set_defaults(func=cmd_set_findings)

    p_rec_fail = sub.add_parser("record-failure", help="Record a single failing test on a patch (increments per-test counter, dedup by round)")
    p_rec_fail.add_argument("--stack-root", required=True)
    p_rec_fail.add_argument("--patch-id", required=True, help="Patch this test failed on")
    p_rec_fail.add_argument("--branch", required=True, help="Suspect branch (where the bug was introduced — used as the dedup key)")
    p_rec_fail.add_argument("--task", required=True)
    p_rec_fail.add_argument("--test", required=True)
    p_rec_fail.add_argument("--repo-root", default=None)
    p_rec_fail.set_defaults(func=cmd_record_failure)

    p_rec_succ = sub.add_parser("record-success", help="Reset consecutive-failure counter for a test that's now passing")
    p_rec_succ.add_argument("--stack-root", required=True)
    p_rec_succ.add_argument("--branch", required=True)
    p_rec_succ.add_argument("--task", required=True)
    p_rec_succ.add_argument("--test", required=True)
    p_rec_succ.add_argument("--patch-id", default=None, help="Patch where this test was observed green (optional)")
    p_rec_succ.add_argument("--repo-root", default=None)
    p_rec_succ.set_defaults(func=cmd_record_success)

    p_rec_fix = sub.add_parser("record-fix", help="Record a fix commit applied to a branch (called by fix-and-commit subagent)")
    p_rec_fix.add_argument("--stack-root", required=True)
    p_rec_fix.add_argument("--branch", required=True)
    p_rec_fix.add_argument("--commit-sha", required=True, help="Full or short SHA of the new commit")
    p_rec_fix.add_argument("--summary", required=True, help="One-line description of what the fix changed")
    p_rec_fix.add_argument("--target-tests", default=None, help="CSV of failing-test keys this fix targets (branch::task::test)")
    p_rec_fix.add_argument("--repo-root", default=None)
    p_rec_fix.set_defaults(func=cmd_record_fix)

    p_master = sub.add_parser("record-master-broken", help="Flag a test as already failing on master (excluded from actionability)")
    p_master.add_argument("--stack-root", required=True)
    p_master.add_argument("--branch", required=True, help="Suspect branch (used as the dedup key)")
    p_master.add_argument("--task", required=True)
    p_master.add_argument("--test", required=True)
    p_master.add_argument("--patch-id", default=None, help="Patch where this was first observed (optional)")
    p_master.add_argument("--evidence", default=None, help="Free-text evidence (e.g. 'failing on master patch abc123 since 2026-04-29')")
    p_master.add_argument("--repo-root", default=None)
    p_master.set_defaults(func=cmd_record_master_broken)

    p_quar = sub.add_parser("quarantine", help="List quarantined tests (>=3 consecutive failures)")
    p_quar.add_argument("--stack-root", required=True)
    p_quar.add_argument("--repo-root", default=None)
    p_quar.set_defaults(func=cmd_quarantine)

    p_sum = sub.add_parser("summary", help="Coordinator-friendly terse summary")
    p_sum.add_argument("--stack-root", required=True)
    p_sum.add_argument("--repo-root", default=None)
    p_sum.set_defaults(func=cmd_summary)

    p_bump = sub.add_parser(
        "bump-poll-iteration",
        help="Increment polling iteration counter (called by poll-status subagent each wakeup)",
    )
    p_bump.add_argument("--stack-root", required=True)
    p_bump.add_argument("--repo-root", default=None)
    p_bump.set_defaults(func=cmd_bump_poll_iteration)

    p_reset = sub.add_parser(
        "reset-poll-cycle",
        help="Zero the polling counter at the start of a polling cycle (Phase 2 entry)",
    )
    p_reset.add_argument("--stack-root", required=True)
    p_reset.add_argument("--repo-root", default=None)
    p_reset.set_defaults(func=cmd_reset_poll_cycle)

    p_sched = sub.add_parser(
        "schedule-next-poll",
        help="Record the planned next-wakeup time so the dashboard can show a countdown",
    )
    p_sched.add_argument("--stack-root", required=True)
    p_sched.add_argument("--in-seconds", type=int, required=True, help="Delay before next ScheduleWakeup fires (e.g. 300)")
    p_sched.add_argument("--repo-root", default=None)
    p_sched.set_defaults(func=cmd_schedule_next_poll)

    p_mbf = sub.add_parser(
        "mark-build-failure",
        help="Flag a patch as having an upstream build failure (RUN_ALL_UNIT_JAVA_TESTS pattern)",
    )
    p_mbf.add_argument("--stack-root", required=True)
    p_mbf.add_argument("--patch-id", required=True)
    p_mbf.add_argument("--repo-root", default=None)
    p_mbf.set_defaults(func=cmd_mark_build_failure)

    p_gffa = sub.add_parser(
        "get-fail-fast-aborts",
        help="Print descendant patch_ids needing abort (one per line); empty if not triggered or already dispatched",
    )
    p_gffa.add_argument("--stack-root", required=True)
    p_gffa.add_argument("--repo-root", default=None)
    p_gffa.set_defaults(func=cmd_get_fail_fast_aborts)

    p_mffa = sub.add_parser(
        "mark-fail-fast-aborted",
        help="Mark compile fail-fast aborts as dispatched (idempotent)",
    )
    p_mffa.add_argument("--stack-root", required=True)
    p_mffa.add_argument("--repo-root", default=None)
    p_mffa.set_defaults(func=cmd_mark_fail_fast_aborted)

    p_archive = sub.add_parser("mark-completed-stack", help="Archive state file under archive/")
    p_archive.add_argument("--stack-root", required=True)
    p_archive.add_argument("--repo-root", default=None)
    p_archive.set_defaults(func=cmd_mark_completed_stack)

    p_rm = sub.add_parser("rm", help="Delete state file")
    p_rm.add_argument("--stack-root", required=True)
    p_rm.add_argument("--repo-root", default=None)
    p_rm.set_defaults(func=cmd_rm)

    p_dash = sub.add_parser("dashboard", help="Force-regenerate the HTML dashboard for a stack root")
    p_dash.add_argument("--stack-root", required=True)
    p_dash.add_argument("--repo-root", default=None)
    p_dash.set_defaults(func=cmd_dashboard)

    p_dash_path = sub.add_parser("dashboard-path", help="Print the dashboard HTML path for a stack root")
    p_dash_path.add_argument("--stack-root", required=True)
    p_dash_path.add_argument("--repo-root", default=None)
    p_dash_path.add_argument(
        "--open",
        action="store_true",
        help="Also open the dashboard in the default browser (opt out via STACK_STATE_NO_OPEN=1)",
    )
    p_dash_path.set_defaults(func=cmd_dashboard_path)

    return parser


def main(argv: list[str] | None = None) -> int:
    args = build_parser().parse_args(argv)
    return args.func(args)


if __name__ == "__main__":
    sys.exit(main())
