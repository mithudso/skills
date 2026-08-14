# Numerical confidence rubric

Used to fill the `Confidence` field in the report's
`Severity classification` and `Root cause hypothesis` sections, and
the `Confidence` slot in the Run summary block.

Express confidence as a single integer in the range **0–100** with a
**≤ 15-word rationale**. The number is calibrated against the
evidence-source rubric below; the rationale names the specific
evidence that produced the number.

## Source-quality bands

| Band | Score | What it means |
| --- | --- | --- |
| **Very high** | 85–100 | Reproducible-and-reproduced locally, OR exact failing assertion line cited from a non-expired raw task log AND a known resolved BF/SERVER with the same fault signature in the past 90 days, OR change-point detector named a `fixed by <sha>` that `git_show` confirms cancels the caused-by. |
| **High** | 70–84 | Exact failing assertion line cited from the raw task log OR a Build-Baron extracted fault category that matches the severity-type identifier verbatim. No corroborating resolved BF. |
| **Medium** | 50–69 | Symptom matches a magic string from `log_patterns.md` (e.g. `ERROR REPORT`, `Heartbeat timed out`) but only via `bb_get_bfg` extracted faults; raw logs unavailable or only summary-level. |
| **Low** | 25–49 | Indirect evidence only — BF Jira ticket prose, sibling-BF inference, or `(local fallback)` git history without log corroboration. |
| **Very low** | 0–24 | No fault marker found; classification rests on the BF title and reporter alone. Pair with `disposition = keep_open_investigating`. |

## Adjustments

- **−10** if any of: raw task log expired (`evg_get_raw_task_logs`
  returned 404 / empty), `bb_*` retry was needed, or the run
  encountered a `(local fallback)` git citation that the report
  could not corroborate with a live Build-Baron fault.
- **−15** if running in Limited-evidence mode (gateway `bb_*` or
  MCP-only `evg_*` down with explicit user approval to continue).
- **+5** if the BF has ≥ 3 prior BFGs in the past 30 days with the
  identical extracted fault category (high recurrence reinforces
  classification certainty).
- **+10** if a sibling BF closed `Gone away` in the past 90 days
  with the same fault signature.

The adjustments are advisory; the band score already encodes most
of the calibration. Cap the final number at 100.

## Reporting

- Severity classification line: `Confidence: <score>/100 — <one-line rationale>`
- Root cause hypothesis line: same shape.
- Run summary block: copy the Root-cause-hypothesis confidence
  verbatim (single source of truth for the report).

If multiple severity types apply, score each separately. If a single
hypothesis combines two evidence kinds (e.g. log line + sibling BF),
pick the lower band and add the `+5` / `+10` adjustment from the
lower-evidence kind.
