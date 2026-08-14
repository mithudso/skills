#!/usr/bin/env python3
"""Compose a Parsley deep-link URL for a failing task or test.

Usage (task-level URL):
  _make_parsley_url.py --task-id <id> --execution <N> [opts]

Usage (test-level URL, when a specific failing test was identified):
  _make_parsley_url.py --task-id <id> --execution <N> --test-id <test>

Options applicable to either form:
  --bookmark <line>      1-indexed line number from the on-disk task log
                         copy. Repeatable. Converted to 0-indexed for
                         Parsley's `bookmarks` query param.
  --filter <pattern>     Regex/literal pattern to seed Parsley's filter
                         box. Repeatable; joined with OR (filterLogic=or).
  --share-line <line>    1-indexed line to anchor the "share" hash.

Prints the URL on stdout. Exits 0 on success, 2 on argument error.

This script is read-only and offline — it does not call Evergreen
or Parsley.
"""
from __future__ import annotations

import argparse
import sys
import urllib.parse

PARSLEY_TEST = "https://parsley.corp.mongodb.com/test/{tid}/{exec}/{test}"
PARSLEY_TASK = "https://parsley.corp.mongodb.com/evergreen/{tid}/{exec}/task"


def _build(args: argparse.Namespace) -> str:
    if args.test_id:
        base = PARSLEY_TEST.format(
            tid=args.task_id, exec=args.execution, test=args.test_id
        )
    else:
        base = PARSLEY_TASK.format(tid=args.task_id, exec=args.execution)
    query: list[tuple[str, str]] = []
    if args.bookmark:
        # Parsley uses 0-indexed line numbers; the CLI takes 1-indexed
        # to match the on-disk task log copy (Step 5.5 Path A).
        zero_indexed = sorted({int(b) - 1 for b in args.bookmark if int(b) >= 1})
        if zero_indexed:
            query.append(("bookmarks", ",".join(str(n) for n in zero_indexed)))
    if args.filter:
        # Parsley double-URL-encodes the filter payload.
        joined = "|".join(args.filter)
        once = urllib.parse.quote(joined, safe="")
        twice = urllib.parse.quote(once, safe="")
        query.append(("filterLogic", "or"))
        query.append(("filters", twice))
    if args.share_line is not None:
        query.append(("shareLine", str(max(0, int(args.share_line) - 1))))
    if not query:
        return base
    suffix = "&".join(f"{k}={v}" for k, v in query)
    return f"{base}?{suffix}"


def main(argv: list[str] | None = None) -> int:
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--task-id", required=True)
    ap.add_argument("--execution", required=True, type=int)
    ap.add_argument("--test-id", default=None)
    ap.add_argument("--bookmark", action="append", default=[])
    ap.add_argument("--filter", action="append", default=[])
    ap.add_argument("--share-line", type=int, default=None)
    args = ap.parse_args(argv)
    print(_build(args))
    return 0


if __name__ == "__main__":
    sys.exit(main())
