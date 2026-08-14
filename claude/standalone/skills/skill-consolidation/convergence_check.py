#!/usr/bin/env python3
"""Deterministic convergence-exit checker for the optimizer family.

Computes the canonical exit verdict (convergence-and-severity.md, this directory) from the
iteration N-1 and N artifact files plus model-counted severity totals passed as CLI flags.
Edit distance is derived from difflib.SequenceMatcher opcodes (insertions + deletions +
substitutions over the longer artifact) — the one computation the model cannot do reliably.

Verdict precedence: CLEAN > STABLE-REWRITE > INSTABILITY > NO-PROGRESS > CONTINUE.
CYCLING (exit condition 3) is model-judged and is NOT computed here; iteration-cap and
budget exits are caller-side.
"""
import argparse
import difflib
import sys


def edit_ratio(prev: str, curr: str) -> float:
    sm = difflib.SequenceMatcher(None, prev, curr, autojunk=False)
    changed = 0
    for tag, i1, i2, j1, j2 in sm.get_opcodes():
        if tag != "equal":
            changed += max(i2 - i1, j2 - j1)
    return changed / (max(len(prev), len(curr)) or 1)


def main() -> int:
    p = argparse.ArgumentParser(
        description="Print the canonical convergence verdict for one loop iteration.",
        epilog="Example: convergence_check.py prev.md curr.md --prev-medium 5 --curr-medium 2",
    )
    p.add_argument("prev", help="path to the iteration N-1 artifact")
    p.add_argument("curr", help="path to the iteration N artifact")
    p.add_argument("--prev-medium", type=int, required=True, help="Medium+ count, iteration N-1")
    p.add_argument("--curr-medium", type=int, required=True, help="Medium+ count, iteration N")
    p.add_argument("--introduced", type=int, default=None,
                   help="Medium+ findings introduced this iteration (with --closed, enables INSTABILITY)")
    p.add_argument("--closed", type=int, default=None,
                   help="Medium+ findings closed this iteration (with --introduced, enables INSTABILITY)")
    p.add_argument("--threshold", type=float, default=0.02,
                   help="stable-rewrite edit-distance ratio (default: 0.02)")
    a = p.parse_args()

    prev = open(a.prev, encoding="utf-8").read()
    curr = open(a.curr, encoding="utf-8").read()
    ratio = edit_ratio(prev, curr)

    if a.curr_medium == 0:
        verdict = "CLEAN"
    elif ratio < a.threshold:
        verdict = "STABLE-REWRITE"
    elif (a.introduced is not None and a.closed is not None
          and (a.introduced or a.closed) and a.introduced >= a.closed):
        verdict = "INSTABILITY"
    elif a.curr_medium >= a.prev_medium:
        verdict = "NO-PROGRESS"
    else:
        verdict = "CONTINUE"

    print(f"verdict: {verdict}")
    print(f"edit-distance ratio: {ratio:.4f} (stable-rewrite threshold: {a.threshold})")
    print(f"medium-plus: prev={a.prev_medium} curr={a.curr_medium}")
    if a.introduced is not None and a.closed is not None:
        print(f"introduced={a.introduced} closed={a.closed}")
    print("note: CYCLING (exit condition 3) is model-judged — not computed here.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
