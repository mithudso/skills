#!/usr/bin/env python3
"""Normalize devprod-mcp-gateway Jira responses into legacy shape for slicing.

The v2 verification slicer (`_slice_helper.py`) expects a single JSON blob
shaped like the classic Jira REST API:

    {
      "key": "...",
      "summary": "...",
      "description": "...",
      "status": {"name": "Closed"},
      "resolution": {"name": "Won't Fix"},
      "resolutiondate": "...",
      "assignee": {"display_name": "..."},
      "created": "...",
      "comments": [
        {"author": {"display_name": "..."}, "created": "...", "body": "..."}
      ],
      "changelogs": [
        {"created": "...", "author": {"display_name": "..."},
         "items": [{"field": "...", "from_string": "...", "to_string": "..."}]}
      ]
    }

The `devprod-mcp-gateway` exposes the data across TWO tools:

  1. `jira_get_issue` (returns issue snapshot with `custom_fields` dict,
     no changelog, scalar status/resolution/assignee strings, nonce-wrapped
     `untrusted_fields`).
  2. `jira_get_issue_comments` (returns a separate `{comments: [...]}` blob,
     also nonce-wrapped).

There is NO changelog tool in the gateway today (verified by
introspecting the gateway tool list — no *changelog* tool is present).
For v2 replay we synthesize:

  - `changelogs: []` — empty, because the gateway can't supply it. The
    slicer falls back to current custom-field values, and the bash
    harness uses the `Team Assigned (Effective Date)` custom field
    as the cutoff timestamp.
  - Scalar status / resolution / assignee strings → wrapped in the
    legacy `{"name": "..."}` / `{"display_name": "..."}` shape.
  - Nonce-wrapped `--- BEGIN UNTRUSTED ... --- ... --- END UNTRUSTED ... ---`
    blocks → unwrapped to the plain body. This is purely for
    triage-internal use; output is never echoed back to a tool.

Caveat: because we have no changelog, the "held-in" snapshot is an
approximation — for BFs whose custom fields changed between routing
and close (e.g. Severity Type reclassified mid-investigation), the
snapshot will reflect the CURRENT value, not the value-at-cutoff.
The sliced.md output gets a warning banner saying so.

This stitcher is a stop-gap. The clean structural fix is for the
gateway to expose a `jira_get_issue_changelog` tool; once that lands,
delete this script and rely on the legacy shape directly.
"""

from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path
from typing import Any

_NONCE_RE = re.compile(
    r"--- BEGIN UNTRUSTED [^\n]*?nonce=[0-9a-f]+ ---\n?"
    r"(?P<body>.*?)"
    r"\n?--- END UNTRUSTED nonce=[0-9a-f]+ ---",
    re.DOTALL,
)


def _strip_nonce(value: Any) -> Any:
    """Unwrap `--- BEGIN UNTRUSTED ... --- ... --- END UNTRUSTED ... ---`
    nonce blocks recursively. Idempotent on values without a wrapper."""
    if isinstance(value, str):
        out = _NONCE_RE.sub(lambda m: m.group("body"), value)
        return out
    if isinstance(value, list):
        return [_strip_nonce(v) for v in value]
    if isinstance(value, dict):
        return {k: _strip_nonce(v) for k, v in value.items()}
    return value


def _coerce_assignee(value: Any) -> dict | None:
    if value is None:
        return None
    if isinstance(value, str):
        s = _strip_nonce(value).strip()
        return {"display_name": s} if s else None
    if isinstance(value, dict):
        return {"display_name": _strip_nonce(value.get("display_name") or value.get("name") or "")}
    return None


def _coerce_name_dict(value: Any) -> dict | None:
    if value is None:
        return None
    if isinstance(value, str):
        return {"name": value}
    if isinstance(value, dict):
        return {"name": value.get("name", "")}
    return None


def stitch(issue: dict, comments_blob: dict | None) -> dict:
    """Return a legacy-shape dict for `_slice_helper.py`."""
    issue_clean = _strip_nonce(issue)
    out: dict[str, Any] = {
        "key": issue_clean.get("key"),
        "summary": issue_clean.get("summary"),
        "description": issue_clean.get("description"),
        "created": issue_clean.get("created"),
        "updated": issue_clean.get("updated"),
        "labels": issue_clean.get("labels") or [],
        "status": _coerce_name_dict(issue_clean.get("status")),
        "resolution": _coerce_name_dict(issue_clean.get("resolution")),
        "resolutiondate": (
            issue_clean.get("resolutiondate")
            or issue_clean.get("updated")
        ),
        "assignee": _coerce_assignee(issue_clean.get("assignee")),
        "reporter": _coerce_assignee(issue_clean.get("reporter")),
        "custom_fields": issue_clean.get("custom_fields") or {},
        "changelogs": [],
        "comments": [],
        "_normalized_from_gateway": True,
    }

    if comments_blob is not None:
        comments_clean = _strip_nonce(comments_blob)
        raw = comments_clean.get("comments") or []
        for c in raw:
            author = c.get("author")
            if isinstance(author, str):
                author = {"display_name": author}
            elif isinstance(author, dict):
                author = {"display_name": author.get("display_name") or author.get("name") or ""}
            else:
                author = {"display_name": ""}
            out["comments"].append({
                "author": author,
                "created": c.get("created"),
                "updated": c.get("updated"),
                "body": c.get("body") or "",
            })

    return out


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--issue-json", required=True,
                    help="Output of `jira_get_issue`")
    ap.add_argument("--comments-json", default=None,
                    help="Output of `jira_get_issue_comments` (optional)")
    ap.add_argument("--out", required=True,
                    help="Path to write the stitched legacy-shape JSON")
    args = ap.parse_args()

    issue = json.loads(Path(args.issue_json).read_text())
    comments_blob = (
        json.loads(Path(args.comments_json).read_text())
        if args.comments_json
        else None
    )
    stitched = stitch(issue, comments_blob)
    Path(args.out).write_text(json.dumps(stitched, indent=2, ensure_ascii=False))
    print(json.dumps({
        "key": stitched["key"],
        "comments": len(stitched["comments"]),
        "changelogs": len(stitched["changelogs"]),
        "team_assigned_effective_date": (
            stitched["custom_fields"].get("Team Assigned (Effective Date)")
        ),
    }))
    return 0


if __name__ == "__main__":
    sys.exit(main())
