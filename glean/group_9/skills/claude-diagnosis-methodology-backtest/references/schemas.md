# JSON schemas & report contract

The exact shapes the orchestration prompt reads and writes. All JSON is a single
array per file unless noted.

## Table of contents
- [Shared staging files](#shared-staging-files)
- [Prediction record](#prediction-record)
- [Grading queue (anonymized) + keymap (private)](#grading-queue--keymap)
- [Raw grade record](#raw-grade-record)
- [De-anonymized per-methodology grades](#de-anonymized-grades)
- [Skipped / blocked records](#skipped--blocked)
- [Synthesis report sections](#synthesis-report-sections)

## Shared staging files

`eval-case-list.json` — every case in scope, recency-sorted (dates parsed, not
lexically sorted):
```json
[{"case_id":"01XXXXXXX","title":"...","severity":"S2","last_updated":"MM/DD/YYYY","status":"Closed"}]
```

`eval-initial-prompts.json` — the blind input handed to every prediction agent. The
`initial_prompt` is the title, optionally plus a customer-authored opening report
(see the leakage guard). `title_only` marks low-signal cases:
```json
[{"case_id":"01XXXXXXX","title":"...","initial_prompt":"...","title_only":false}]
```

## Prediction record

One object per case per methodology. Written by a prediction agent; it never reads
the resolutions file and never grades.
```json
{
  "case_number": "01XXXXXXX",
  "phase": 1,
  "initial_prompt_used": "<verbatim copy from eval-initial-prompts.json>",
  "predicted_root_cause": "<one sentence>",
  "predicted_evidence_to_collect": ["<diagnostic 1>", "<diagnostic 2>"],
  "predicted_steps": ["<step 1>", "<step 2>", "<step 3>"],
  "citation": {"source_type": "skill|flowchart|expertise", "source_id": "<id>"},
  "confidence": "high|medium|low",
  "predicted_at": "<ISO timestamp set just before write>"
}
```
Validation gate before grading: array of objects, every required field present,
record count == case count minus skipped. Failing records move to skipped (below).

## Grading queue & keymap

Built by the orchestrator after predictions validate. The queue is what the grader
sees — **anonymized**: no `phase`, no `citation`, shuffled order, opaque `record_id`.
`eval-grading-queue.json`:
```json
[{"record_id":"r-0001","case_number":"01XXXXXXX","initial_prompt_used":"...","predicted_root_cause":"...","predicted_steps":["..."],"confidence":"high"}]
```
`eval-grading-keymap.json` — **private; the grader never reads this**:
```json
[{"record_id":"r-0001","phase":1,"case_number":"01XXXXXXX"}]
```

## Raw grade record

Written by the grader, keyed only by `record_id`. Every grade quotes the resolution.
`eval-grades-raw.json`:
```json
[{"record_id":"r-0001","case_number":"01XXXXXXX","grade":"Correct|Partial|Wrong|Unverifiable","justification_quote":"<verbatim from final_solution or troubleshooting_steps>"}]
```

## De-anonymized grades

The orchestrator joins raw grades to the keymap on `record_id` and splits per
methodology into `eval-phase-N-grades.json` (same fields as raw, plus restored
`phase`), then writes per-methodology `eval-phase-N-report.md` (per-case table +
scorecard + 5 worked grading examples drawn from that phase).

## Skipped / blocked

`eval-phase-N-skipped.json` — cases an agent could not honor a rule for, or records
that failed schema validation:
```json
[{"case_number":"01XXXXXXX","reason":"<one sentence>"}]
```
`eval-phase-N-blocked.md` — written instead of predictions when a methodology's
knowledge source can't be located; describes exactly what was searched, then STOP.

## Synthesis report sections

`eval-synthesis-report.md` — required sections, in order:
1. **Methodology recap** — what each phase used, what was blocked, predict/grade
   separation, parallel wall-clock.
2. **Per-case comparison table** — one row per case × one grade column per
   methodology; grouped by status for readability.
3. **Aggregate scorecard** — three scorecards: (a) **headline** = non-title-only AND
   strong-resolution cases; (b) title-only cases; (c) fallback/unavailable cases. The
   headline comparison is (a).
4. **Resolution-source caveat** — per case, strong vs. closed-fallback vs.
   unavailable; strong cases carry more signal weight.
5. **Disagreement analysis** — cases where methodologies disagreed; which was right
   and why.
6. **Methodology meta-analysis** — most accurate? most efficient? most explainable? —
   scoped to cohort (a), 2-3 sentences each.
7. **Limitations & threats to validity** — title-only inputs, resolution-quality
   variance, single-account bias, single-grader (inter-rater reliability unmeasured),
   no random/keyword baseline noise floor.
8. **Recommendation** — one methodology or a hybrid, explicitly caveated by §7.
