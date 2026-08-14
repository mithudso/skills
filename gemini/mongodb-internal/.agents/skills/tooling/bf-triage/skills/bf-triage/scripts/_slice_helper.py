#!/usr/bin/env python3
"""Render Jira BF JSON into held-in / held-out markdown slices.

Invoked by `run_held_in_test.sh`. Not meant to be called directly by the
triager subagent.

Two input shapes are supported:

  1. Legacy single-blob shape (classic Jira REST API). Has top-level
     `changelogs` and `comments` arrays. Custom fields (Severity Type,
     Failure Type, etc.) are reconstructed at cutoff by replaying every
     changelog event with `created < cutoff`. This is the precise mode.

  2. devprod-mcp-gateway shape, stitched by `_stitch_gateway_jira.py`
     (marked with `_normalized_from_gateway: true`). The gateway does
     NOT expose a changelog tool today, so `changelogs` is always empty.
     For the held-in snapshot we fall back to the CURRENT
     `custom_fields` values — an approximation that is accurate only
     when fields did not change between routing and close. sliced.md
     carries a prominent warning banner when running in this mode.

The clean structural fix is for `devprod-mcp-gateway` to expose a
`jira_get_issue_changelog` tool; once that lands, the gateway path
becomes equivalent to the legacy path and the warning banner goes away.
"""

from __future__ import annotations

import argparse
import json
import sys
from datetime import datetime
from pathlib import Path

TRACKED_FIELDS = (
    "summary",
    "description",
    "Severity Type",
    "Failure Type",
    "Bug Symptoms",
    "Assigned Teams",
    "Performance Change Type",
    "Evergreen Project",
    "Failing Buildvariants",
    "Failing Tasks",
    "Temperature",
    "Score",
    "Score Breakdown",
    "Investigation Summary",
    "labels",
    "status",
    "assignee",
)


_TZ_OFFSET_RE = __import__("re").compile(r"([+-]\d{2})(\d{2})$")


def parse_ts(value: str) -> datetime:
    """Parse Jira timestamps. Handles both `+00:00` (changelog) and `+0000`
    (comments) timezone shapes, with or without microseconds."""
    if value.endswith("Z"):
        value = value[:-1] + "+00:00"
    m = _TZ_OFFSET_RE.search(value)
    if m and ":" not in value[-6:]:
        value = value[: m.start()] + f"{m.group(1)}:{m.group(2)}"
    return datetime.fromisoformat(value)


def _stringify(value: object) -> str:
    """Render a custom-field value (str / list / None) as a flat string."""
    if value is None:
        return ""
    if isinstance(value, str):
        return value
    if isinstance(value, list):
        return ", ".join(_stringify(v) for v in value if v not in (None, ""))
    return str(value)


def fields_at_cutoff(issue: dict, cutoff: datetime) -> dict:
    snapshot: dict[str, str] = {f: "" for f in TRACKED_FIELDS}

    cl = sorted(issue.get("changelogs") or [], key=lambda c: parse_ts(c["created"]))
    description_set_pre_cutoff = False
    summary_set_pre_cutoff = False
    cutoff_event_from_strings: dict[str, str] = {}
    for entry in cl:
        ts = parse_ts(entry["created"])
        if ts >= cutoff:
            if ts == cutoff:
                for item in entry.get("items") or []:
                    field = item.get("field")
                    if field in snapshot and item.get("from_string") is not None:
                        cutoff_event_from_strings[field] = item["from_string"]
            continue
        for item in entry.get("items") or []:
            field = item.get("field")
            if field not in snapshot:
                continue
            value = item.get("to_string") or ""
            snapshot[field] = value
            if field == "description":
                description_set_pre_cutoff = True
            if field == "summary":
                summary_set_pre_cutoff = True

    for field, frm in cutoff_event_from_strings.items():
        if not snapshot.get(field):
            snapshot[field] = frm

    # Gateway-normalized fallback: when there is no changelog, the
    # replay above yields an all-empty snapshot. Fill from CURRENT
    # custom_fields as a best-effort. This is approximate (post-cutoff
    # value) but better than nothing; render_md() adds a warning banner.
    if issue.get("_normalized_from_gateway") and not cl:
        cf = issue.get("custom_fields") or {}
        for field in TRACKED_FIELDS:
            if snapshot.get(field):
                continue
            if field in cf:
                snapshot[field] = _stringify(cf[field])

    if not description_set_pre_cutoff:
        snapshot["description"] = issue.get("description") or ""
    if not summary_set_pre_cutoff:
        snapshot["summary"] = issue.get("summary") or ""
    return snapshot


def render_field_block(snapshot: dict) -> list[str]:
    lines: list[str] = []
    for field in TRACKED_FIELDS:
        value = snapshot.get(field, "")
        if field in {"description", "Investigation Summary", "Bug Symptoms"}:
            continue
        if not value:
            value = "_(empty at cutoff)_"
        lines.append(f"- **{field}**: {value}")
    return lines


def render_long_field(label: str, value: str) -> list[str]:
    out = [f"### {label}", ""]
    if not value:
        out.append("_(empty at cutoff)_")
    else:
        out.append("```text")
        out.extend(value.splitlines() or [""])
        out.append("```")
    out.append("")
    return out


def render_comment(comment: dict) -> list[str]:
    author = (comment.get("author") or {}).get("display_name") or "unknown"
    created = comment.get("created", "?")
    body = (comment.get("body") or "").strip()
    out = [f"### Comment by {author} at {created}", ""]
    if not body:
        body = "_(empty)_"
    out.append("```text")
    out.extend(body.splitlines())
    out.append("```")
    out.append("")
    return out


def render_changelog(entry: dict) -> list[str]:
    author = (entry.get("author") or {}).get("display_name") or "unknown"
    created = entry.get("created", "?")
    out = [f"### Changelog at {created} by {author}", ""]
    for item in entry.get("items") or []:
        field = item.get("field", "?")
        frm = item.get("from_string")
        to = item.get("to_string")
        frm_disp = frm if frm not in (None, "") else "_(empty)_"
        to_disp = to if to not in (None, "") else "_(empty)_"
        out.append(f"- `{field}`: {frm_disp} → {to_disp}")
    out.append("")
    return out


def render_md(
    *,
    title: str,
    issue: dict,
    snapshot: dict | None,
    comments: list,
    changelogs: list,
    cutoff_iso: str,
    mode: str,
) -> str:
    out: list[str] = [f"# {title}", ""]
    out.append(f"- BF Jira key: `{issue.get('key', '?')}`")
    out.append(f"- Created: {issue.get('created', '?')}")
    out.append(f"- Cutoff: `{cutoff_iso}` (first `Assigned Teams = Workload Resilience`)")
    side = "< cutoff" if mode == "held-in" else ">= cutoff"
    out.append(f"- Slice mode: **{mode}** (events with timestamp {side})")
    out.append("")

    gateway_mode = bool(issue.get("_normalized_from_gateway"))
    if gateway_mode and mode == "held-in":
        out.append(
            "> **APPROXIMATION WARNING — gateway-normalized snapshot.**"
        )
        out.append(
            "> This BF was fetched via `devprod-mcp-gateway`, which does "
            "not expose a Jira changelog tool. The held-in field values "
            "below come from the BF's CURRENT `custom_fields` and are "
            "therefore a best-effort approximation. For BFs whose fields "
            "changed between routing and close (e.g. Severity Type "
            "reclassified mid-investigation), the snapshot is NOT the "
            "value-at-cutoff. The cutoff timestamp itself is reliable; "
            "it is read from the `Team Assigned (Effective Date)` "
            "custom field."
        )
        out.append("")

    out.append("## Fields")
    out.append("")
    if mode == "held-in":
        if gateway_mode:
            out.append(
                "Field values taken from current `custom_fields` (gateway "
                "stitched mode; no changelog available). Empty values mean "
                "the field is not currently set on the BF."
            )
        else:
            out.append(
                "Field values reconstructed from the Jira changelog by replaying "
                "every event with `created < cutoff`. Empty values mean the field "
                "was not yet set when the BF was routed to Workload Resilience."
            )
        out.append("")
        if snapshot:
            out.extend(render_field_block(snapshot))
            out.append("")
            out.extend(render_long_field("Description (at cutoff)", snapshot.get("description", "")))
            out.extend(render_long_field("Bug Symptoms (at cutoff)", snapshot.get("Bug Symptoms", "")))
            out.extend(render_long_field("Investigation Summary (at cutoff)", snapshot.get("Investigation Summary", "")))
    else:
        out.append(
            "These are the **final** ticket-level fields (post-cutoff). "
            "Use them as ground truth when grading."
        )
        out.append("")
        out.append(f"- Final summary: {issue.get('summary', '?')}")
        out.append(f"- Final status: {(issue.get('status') or {}).get('name', '?')}")
        out.append(
            f"- Final resolution: {(issue.get('resolution') or {}).get('name', 'Unresolved')}"
        )
        out.append(f"- Resolution date: {issue.get('resolutiondate', '?')}")
        final_assignee = (issue.get("assignee") or {}).get("display_name", "Unassigned")
        out.append(f"- Final assignee: {final_assignee}")
        out.append("")

    out.append(f"## Comments ({len(comments)})")
    out.append("")
    if not comments:
        out.append("_(none in this slice)_")
        out.append("")
    for comment in comments:
        out.extend(render_comment(comment))

    out.append(f"## Changelog ({len(changelogs)})")
    out.append("")
    if not changelogs:
        out.append("_(none in this slice)_")
        out.append("")
    for entry in changelogs:
        out.extend(render_changelog(entry))

    return "\n".join(out) + "\n"


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--bf", required=True)
    ap.add_argument("--cutoff", required=True)
    ap.add_argument("--input", required=True)
    ap.add_argument("--sliced", required=True)
    ap.add_argument("--heldout", required=True)
    args = ap.parse_args()

    issue = json.loads(Path(args.input).read_text())
    cutoff_dt = parse_ts(args.cutoff)

    comments = issue.get("comments") or []
    if isinstance(comments, dict):
        comments = comments.get("comments") or []
    changelogs = issue.get("changelogs") or []

    held_in_comments = [c for c in comments if parse_ts(c["created"]) < cutoff_dt]
    held_out_comments = [c for c in comments if parse_ts(c["created"]) >= cutoff_dt]
    held_in_cl = [x for x in changelogs if parse_ts(x["created"]) < cutoff_dt]
    held_out_cl = [x for x in changelogs if parse_ts(x["created"]) >= cutoff_dt]

    snapshot = fields_at_cutoff(issue, cutoff_dt)

    sliced_md = render_md(
        title=f"{args.bf} — held-in slice (pre-cutoff)",
        issue=issue,
        snapshot=snapshot,
        comments=held_in_comments,
        changelogs=held_in_cl,
        cutoff_iso=args.cutoff,
        mode="held-in",
    )
    heldout_md = render_md(
        title=f"{args.bf} — held-out slice (post-cutoff, GRADER ONLY)",
        issue=issue,
        snapshot=None,
        comments=held_out_comments,
        changelogs=held_out_cl,
        cutoff_iso=args.cutoff,
        mode="held-out",
    )

    Path(args.sliced).write_text(sliced_md)
    Path(args.heldout).write_text(heldout_md)

    counts = {
        "held_in_comments": len(held_in_comments),
        "held_out_comments": len(held_out_comments),
        "held_in_changelogs": len(held_in_cl),
        "held_out_changelogs": len(held_out_cl),
    }
    print(json.dumps(counts))
    return 0


if __name__ == "__main__":
    sys.exit(main())
