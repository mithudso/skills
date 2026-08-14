#!/usr/bin/env python3
"""UserPromptSubmit hook: cheap, free, no-network model/effort recommendation
plus a price estimate for the prompt just submitted. Runs on every prompt.

Writes ~/.claude/cost-estimator/<session_id>.json for a status line to render,
and emits a one-line additionalContext note so the assistant sees the
recommendation immediately — no hook or tool can force-switch the live
session's model, so surfacing it as context (for the assistant or user to act
on via /model + /effort) is the closest real automation available.

Never blocks or fails the prompt: any error is caught, breadcrumbed to
~/.claude/cost-estimator/hook-errors.log (best-effort, capped at 50 lines),
and the hook exits with no output. No network calls — pure local heuristic,
same one `estimate_cost.py --recommend` uses.
"""
import json
import os
import re
import sys
import time
from pathlib import Path

SKILL_DIR = Path(__file__).resolve().parent
sys.path.insert(0, str(SKILL_DIR))

CLAUDE_DIR = Path(os.environ.get("CLAUDE_CONFIG_DIR") or (Path.home() / ".claude"))

_SAFE_ID = re.compile(r"[^A-Za-z0-9_-]")


def _sanitize_session_id(raw: str) -> str:
    """Strip anything but [A-Za-z0-9_-] so this can never escape state_dir
    via path separators or '..' segments, however session_id was derived."""
    cleaned = _SAFE_ID.sub("_", raw or "")
    return cleaned or "unknown"


def _log_error(context: str, exc: Exception) -> None:
    """Best-effort breadcrumb so a broken hook is diagnosable, not just
    silently dead forever. Never allowed to raise back into main()."""
    try:
        log_path = CLAUDE_DIR / "cost-estimator" / "hook-errors.log"
        log_path.parent.mkdir(parents=True, exist_ok=True)
        lines = []
        if log_path.exists():
            lines = log_path.read_text().splitlines()[-49:]
        lines.append(f"{time.strftime('%Y-%m-%dT%H:%M:%S')} {context}: {exc!r}")
        log_path.write_text("\n".join(lines) + "\n")
    except Exception:
        pass


def main():
    try:
        raw = sys.stdin.read()
        data = json.loads(raw) if raw.strip() else {}
    except Exception as e:
        _log_error("parse stdin", e)
        return

    prompt = (data.get("prompt") or "").strip()
    if not prompt:
        return

    session_id = data.get("session_id")
    if not session_id:
        transcript = data.get("transcript_path") or ""
        session_id = Path(transcript).stem if transcript else "unknown"
    session_id = _sanitize_session_id(str(session_id))

    try:
        from estimate_cost import (
            recommend, count_tokens_heuristic, build_estimates,
            DEFAULT_MULTIPLIERS, MODELS, fmt_usd,
        )
    except Exception as e:
        _log_error("import estimate_cost", e)
        return

    try:
        rec = recommend(prompt, list(MODELS))
        input_tokens = count_tokens_heuristic(prompt)
        rows = build_estimates(input_tokens, 2000, DEFAULT_MULTIPLIERS, [rec["model"]])
        row = next((r for r in rows if r["effort"] == rec["effort"]), rows[0])
        price = row["total_cost"]
    except Exception as e:
        _log_error("compute recommendation", e)
        return

    state = {
        "ts": time.time(),
        "tier": rec["tier"],
        "model": rec["model"],
        "model_display": MODELS[rec["model"]]["display"],
        "effort": rec["effort"],
        "price": price,
        "price_fmt": fmt_usd(price),
        "input_tokens_est": input_tokens,
        "source": "heuristic",
    }

    try:
        state_dir = CLAUDE_DIR / "cost-estimator"
        state_dir.mkdir(parents=True, exist_ok=True)
        tmp = state_dir / f".{session_id}.json.tmp"
        dest = state_dir / f"{session_id}.json"
        tmp.write_text(json.dumps(state))
        tmp.replace(dest)  # atomic, so the status line never reads a half-written file
    except Exception as e:
        _log_error("write state file", e)

    eff_str = "" if rec["effort"] in ("-", None) else f" @ {rec['effort']}"
    switch = f"/model {rec['model']}"
    if rec["effort"] not in ("-", None):
        switch += f" and /effort {rec['effort']}"
    note = (f"[cost-estimate] ~{fmt_usd(price)} · recommended: "
            f"{MODELS[rec['model']]['display']}{eff_str} (tier: {rec['tier']}) — "
            f"switch with {switch} if warranted")

    try:
        print(json.dumps({
            "hookSpecificOutput": {
                "hookEventName": "UserPromptSubmit",
                "additionalContext": note,
            }
        }))
    except Exception as e:
        _log_error("emit hook output", e)


if __name__ == "__main__":
    main()
