#!/usr/bin/env python3
"""Deterministic Concept Viability Score (CVS) checker for concept-family-explorer.

Recomputes each scored gap-table row from the canonical rubric
(concept-family-explorer references/scoring-rubric.md):

    CVS = 0.25*Relevance + 0.25*Usefulness + 0.20*Novelty
        + 0.15*Interest  + 0.15*Viability          (range 0-5)

and asserts |delta| < 0.005 against the row's stated CVS — the one computation
the model must never verify in its head. Also enforces, as decision checks:
  - hard gates: Viability <= 1 -> SKIP; Novelty = 0 -> SKIP
  - threshold decisions: CVS >= threshold -> RESEARCH (REFRESH when STALE),
    else SKIP

Near-threshold tie ORDERING (Via -> Rel -> Int) is printed as an ADVISORY
only: per the rubric, the tie-break governs queue order, not decisions, and
ordering stays a judgment check.

Scope: cfe's CVS arithmetic only. Not part of the optimizer-family convergence
contract (convergence-and-severity.md / convergence_check.py).

Input: a JSON file (or '-' for stdin) containing either a list of row objects
or {"threshold": <x>, "rows": [...]}. Row keys are case-insensitive; accepted
aliases: concept/name; rel/relevance; use/usefulness; nov/novelty;
int/interest; via/viability; cvs/score; decision; stale (bool) or
tag/status == "STALE". A missing/non-numeric CVS (e.g. "—" on a hard-gated
row) skips the arithmetic check for that row; gates and decision are still
checked.

Exit status: 0 = PASS, 1 = one or more arithmetic/gate/decision failures,
2 = unusable input.
"""
import argparse
import json
import sys

WEIGHTS = {"rel": 0.25, "use": 0.25, "nov": 0.20, "int": 0.15, "via": 0.15}
ALIASES = {
    "concept": ("concept", "name"),
    "rel": ("rel", "relevance"),
    "use": ("use", "usefulness"),
    "nov": ("nov", "novelty"),
    "int": ("int", "interest"),
    "via": ("via", "viability"),
    "cvs": ("cvs", "score"),
    "decision": ("decision",),
    "stale": ("stale",),
    "tag": ("tag", "status"),
}
TOL = 0.005
EPS = 1e-9


def field(row: dict, key: str):
    lowered = {str(k).lower(): v for k, v in row.items()}
    for alias in ALIASES[key]:
        if alias in lowered:
            return lowered[alias]
    return None


def as_float(value):
    if value is None:
        return None
    try:
        return float(value)
    except (TypeError, ValueError):
        return None


def expected_decisions(rel, nov, via, cvs, threshold, stale):
    """Return (set of acceptable decisions, reason)."""
    if via is not None and via <= 1:
        return {"SKIP"}, "hard gate: Via <= 1"
    if nov is not None and nov == 0:
        return {"SKIP"}, "hard gate: Nov = 0"
    if cvs is None:
        return None, "no CVS to check threshold against"
    if cvs + EPS >= threshold:
        if stale is True:
            return {"REFRESH"}, "STALE + above threshold"
        if stale is False:
            return {"RESEARCH"}, "above threshold"
        return {"RESEARCH", "REFRESH"}, "above threshold (staleness not given)"
    return {"SKIP"}, "below threshold"


def main() -> int:
    p = argparse.ArgumentParser(
        description="Verify a cfe scored gap table: CVS arithmetic, hard gates, threshold decisions.",
        epilog="Example: cvs_check.py table.json --threshold 3.2",
    )
    p.add_argument("table", help="path to the scored gap table as JSON, or '-' for stdin")
    p.add_argument("--threshold", type=float, default=None,
                   help="CVS research threshold (default: table's 'threshold' key, else 3.2)")
    p.add_argument("--tie-band", type=float, default=0.25,
                   help="advisory near-threshold band for tie-ordering output (default: 0.25)")
    a = p.parse_args()

    raw = sys.stdin.read() if a.table == "-" else open(a.table, encoding="utf-8").read()
    try:
        data = json.loads(raw)
    except json.JSONDecodeError as e:
        print(f"error: input is not valid JSON ({e})")
        return 2

    if isinstance(data, dict):
        rows = data.get("rows")
        threshold = a.threshold if a.threshold is not None else float(data.get("threshold", 3.2))
    else:
        rows = data
        threshold = a.threshold if a.threshold is not None else 3.2
    if not isinstance(rows, list) or not rows:
        print("error: no rows found (expected a JSON list or an object with a 'rows' list)")
        return 2

    failures = 0
    arith_failures = 0
    checked = []
    for i, row in enumerate(rows):
        name = field(row, "concept") or f"row {i + 1}"
        axes = {k: as_float(field(row, k)) for k in ("rel", "use", "nov", "int", "via")}
        problems = []

        missing = [k for k, v in axes.items() if v is None]
        computed = None
        if missing:
            problems.append(f"missing/non-numeric axis score(s): {', '.join(missing)}")
        else:
            computed = sum(WEIGHTS[k] * axes[k] for k in WEIGHTS)
            stated = as_float(field(row, "cvs"))
            if stated is None:
                print(f"  note  [{name}] no numeric CVS stated — computed {computed:.2f}; arithmetic check skipped")
            elif abs(stated - computed) >= TOL:
                problems.append(f"CVS mismatch: stated {stated} vs computed {computed:.4f} (|delta| >= {TOL})")
                arith_failures += 1

        stale = field(row, "stale")
        if stale is None:
            tag = field(row, "tag")
            stale = (str(tag).strip().upper() == "STALE") if tag is not None else None
        decision_raw = field(row, "decision")
        decision = str(decision_raw).strip().upper() if decision_raw is not None else None
        if decision == "SKIPPED-PRIOR":  # Step 4 warm-start inherits a prior SKIP
            decision = "SKIP"

        expected, reason = expected_decisions(axes["rel"], axes["nov"], axes["via"],
                                              computed, threshold, stale)
        if decision is None:
            problems.append("no decision stated")
        elif expected is not None and decision not in expected:
            problems.append(f"decision {decision} but expected {'/'.join(sorted(expected))} ({reason})")

        if problems:
            failures += 1
            for msg in problems:
                print(f"  FAIL  [{name}] {msg}")
        if computed is not None:
            checked.append((name, computed, axes))

    # Advisory: near-threshold tie ordering (Via -> Rel -> Int). Queue order only — never a failure.
    band = [(n, c, x) for n, c, x in checked if abs(c - threshold) <= a.tie_band]
    if len(band) > 1:
        ordered = sorted(band, key=lambda t: (-round(t[1], 2), -t[2]["via"], -t[2]["rel"], -t[2]["int"]))
        print(f"advisory: near-threshold rows (|CVS - {threshold}| <= {a.tie_band}), "
              "tie-break order Via -> Rel -> Int:")
        for n, c, x in ordered:
            print(f"    {c:.2f}  Via={x['via']:g} Rel={x['rel']:g} Int={x['int']:g}  {n}")
        print("advisory: tie ordering is a judgment check — it governs queue order, not decisions.")

    verdict = "PASS" if failures == 0 else "FAIL"
    print(f"result: {verdict} ({len(rows)} rows, {arith_failures} arithmetic mismatches, "
          f"{failures} rows with failures; threshold {threshold})")
    return 0 if failures == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
