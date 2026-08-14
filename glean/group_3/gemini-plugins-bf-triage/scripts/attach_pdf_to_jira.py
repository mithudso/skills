#!/usr/bin/env python3
"""Attach a file (typically a triage PDF) to a Jira issue via REST.

Usage:
  attach_pdf_to_jira.py <ISSUE-KEY> <file-path>

Token resolution (first hit wins, checked regardless of which agent
runtime invokes the script — a Claude Code run will happily use a
token found in the Cursor config, and vice-versa):
  1. $JIRA_PERSONAL_TOKEN environment variable
  2. The first `JIRA_PERSONAL_TOKEN` value found, scanning in order:
       ~/.cursor/mcp.json        (Cursor MCP config)
       ~/.claude.json            (Claude Code MCP config)
       ~/.claude/settings.json   (Claude Code settings)
     The key may be nested anywhere (e.g. under mcpServers.<srv>.env).

Jira base URL comes from $JIRA_URL (default https://jira.mongodb.org).

Exit codes (so callers can degrade gracefully — an attach failure must
never fail the triage run):
  0  uploaded (HTTP 200)
  1  bad CLI usage
  2  no token found in env or any MCP config
  3  input file not readable
  4  upload failed (non-200 response or network error)

Security: the token is read in-process and sent only as an
Authorization header. It is NEVER printed, NEVER passed on a command
line, and NEVER written to disk.
"""
from __future__ import annotations

import json
import mimetypes
import os
import sys
import urllib.request
import uuid

JIRA_BASE = os.environ.get("JIRA_URL", "https://jira.mongodb.org").rstrip("/")

# Header required by Jira's XSRF guard; the upload is rejected without it.
_XSRF_HEADER = ("X-Atlassian-Token", "no-check")


def _find_token() -> str | None:
    env = os.environ.get("JIRA_PERSONAL_TOKEN")
    if env:
        return env
    home = os.path.expanduser("~")
    candidates = (
        os.path.join(home, ".cursor", "mcp.json"),
        os.path.join(home, ".claude.json"),
        os.path.join(home, ".claude", "settings.json"),
    )
    for path in candidates:
        try:
            with open(path) as fh:
                doc = json.load(fh)
        except Exception:
            continue
        stack = [doc]
        while stack:
            node = stack.pop()
            if isinstance(node, dict):
                for key, val in node.items():
                    if key == "JIRA_PERSONAL_TOKEN" and isinstance(val, str) and val:
                        return val
                    stack.append(val)
            elif isinstance(node, list):
                stack.extend(node)
    return None


def _upload(issue: str, path: str, token: str) -> int:
    with open(path, "rb") as fh:
        payload = fh.read()
    filename = os.path.basename(path)
    ctype = mimetypes.guess_type(filename)[0] or "application/octet-stream"
    boundary = uuid.uuid4().hex
    body = b"".join(
        [
            f"--{boundary}\r\n".encode(),
            (
                f'Content-Disposition: form-data; name="file"; '
                f'filename="{filename}"\r\n'
            ).encode(),
            f"Content-Type: {ctype}\r\n\r\n".encode(),
            payload,
            f"\r\n--{boundary}--\r\n".encode(),
        ]
    )
    req = urllib.request.Request(
        f"{JIRA_BASE}/rest/api/2/issue/{issue}/attachments",
        data=body,
        method="POST",
        headers={
            "Authorization": f"Bearer {token}",
            _XSRF_HEADER[0]: _XSRF_HEADER[1],
            "Content-Type": f"multipart/form-data; boundary={boundary}",
        },
    )
    try:
        with urllib.request.urlopen(req) as resp:
            raw = resp.read() or b"[]"
        parsed = json.loads(raw)
        summary = (
            ", ".join(
                f"{a.get('filename')} (id={a.get('id')}, {a.get('size')}B)"
                for a in parsed
            )
            if isinstance(parsed, list)
            else str(parsed)[:200]
        )
        print(f"[attach] uploaded to {issue}: {summary}")
        return 0
    except Exception as exc:  # noqa: BLE001 - report verbatim, never raise
        print(f"[attach] upload failed: {exc}", file=sys.stderr)
        return 4


def main(argv: list[str]) -> int:
    if len(argv) != 3:
        print(
            "usage: attach_pdf_to_jira.py <ISSUE-KEY> <file-path>",
            file=sys.stderr,
        )
        return 1
    issue, path = argv[1], argv[2]
    if not os.path.isfile(path):
        print(f"[attach] cannot read input file: {path}", file=sys.stderr)
        return 3
    token = _find_token()
    if not token:
        print(
            "[attach] no JIRA_PERSONAL_TOKEN in env or MCP config "
            "(~/.cursor/mcp.json, ~/.claude.json, ~/.claude/settings.json) "
            "— skipping upload",
            file=sys.stderr,
        )
        return 2
    return _upload(issue, path, token)


if __name__ == "__main__":
    sys.exit(main(sys.argv))
