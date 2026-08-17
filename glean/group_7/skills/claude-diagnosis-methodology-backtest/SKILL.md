---
name: diagnosis-methodology-backtest
description: >-
  Run a blind, parallel, multi-agent backtest comparing competing diagnosis
  methodologies against ground-truth case resolutions to score which predicts
  root causes most accurately. TRIGGER: "which approach diagnoses these cases
  best", "backtest my diagnosis strategies", "blind eval over our cases",
  "run the 3-phase eval", any work inside a strategy-backtest scoreboard repo,
  MongoDB/Atlas TAM support-case diagnosis methodology comparison. SKIP:
  diagnosing a single live case (use atlas-diagnostics-expert or mongodb-kb);
  general prompt optimization (use prompt-deep-optimizer); plain writing tasks.
version: "1.2.0"
updated: "2026-06-01"
category: developer
tags: [diagnosis, backtest, eval, methodology, support-cases, mongodb, tam, blind-eval, grading]
whenToUse:
  - "which approach diagnoses these cases best"
  - "backtest my diagnosis strategies"
  - "blind eval over our support cases"
  - "score predictions against the case resolutions"
  - "run the 3-phase eval"
  - "compare the flowchart vs expertise approaches"
  - "methodology comparison over historical cases"
  - "strategy-backtest scoreboard workflow"
whenNotToUse:
  - "diagnosing a single live case — use atlas-diagnostics-expert or mongodb-kb"
  - "general prompt optimization — use prompt-deep-optimizer"
  - "plain writing or report generation"
related_skills: [atlas-diagnostics-expert, misc-catch-all, prompt-deep-optimizer, technical-writing-craft]
---

# Diagnosis Methodology Backtest

## What this skill does

It measures **which of several diagnosis methodologies most accurately predicts the
root cause** of historical support cases — by having each methodology predict
*blind* (from only the case title + the customer's opening message), then grading
those predictions against the cases' real resolutions.

The whole point is a *fair, defensible comparison*. Three design choices make the
result trustworthy, and they are the reason this is a skill and not a one-line
prompt:

1. **Blind prediction.** A methodology sees only what the customer first reported —
   never the resolution. Otherwise you're measuring hindsight, not diagnosis.
2. **Predictor ≠ grader.** The agent that makes predictions never grades its own
   work. A separate grader, blind to *which* methodology produced each prediction,
   scores everything against one rubric. Self-grading lets each contestant mark its
   own exam — the numbers become indefensible.
3. **Parallel + isolated.** Each methodology runs in its own subagent context, all
   dispatched in one message. This is faster *and* makes "no cross-methodology
   leakage" physically true instead of a promise.

The output is a per-case grade for every methodology plus a synthesis report that
recommends one (or a hybrid) for production use, honestly caveated by the quality of
the ground truth.

## When to reach for it

Use it when someone wants to know *which way of diagnosing is best*, backed by data:
comparing troubleshooting playbooks, flowcharts, model prompts, or expert heuristics
against a corpus of resolved cases. It fits naturally in a "strategy backtest
scoreboard" workflow.

Don't use it to diagnose one live case (that's a single prediction, no comparison),
to optimize a prompt's wording (use a prompt optimizer), or for general writing.

## Before you start: confirm the inputs exist

This eval is **only as good as its ground truth**, and it must never fall back to
live data fetches mid-run (that would break blindness and reproducibility). Confirm
up front that you have, as static files:

- **A case-header source** — one row per case with at least `case_id` and `title`.
- **A resolution ground-truth source** — per case, a `final_solution` and/or
  `troubleshooting_steps`. This is what predictions are graded against.
- **One knowledge source per methodology** — e.g. a set of skills, a flowchart doc,
  a playbook corpus. If a methodology's source can't be located, that methodology is
  *blocked* (recorded, not silently skipped).

If any of these is missing, say so and stop — do not substitute a live API/MCP fetch.

## Parameters (collect these, then fill the template)

Gather these from the user / repo, then instantiate `references/orchestration-prompt.md`:

| Parameter | Meaning | MongoDB/TAM default |
|---|---|---|
| `SUMMARY_FILE` | Case headers (`case_id`, `title`, …) | the account's `*-cases-summary.json` |
| `RESOLUTIONS_FILE` | Ground-truth `final_solution` + `troubleshooting_steps` | the enriched `*-cases-*.json` |
| `OUTPUT_DIR` | Where all artifacts land | `customer-files/<account>/` |
| `METHODOLOGIES` | The list of approaches to compare (see below) | the 3 below |
| `BATCH_SIZE` | Cases per write batch (truncation safety) | 40 |
| `EXPECTED_CASE_COUNT` | Assert-after-load sanity check | (the dataset's count) |

`METHODOLOGIES` is a list — **this is the generalization knob**. Each entry is a
"phase":

```
- name: "MongoDB skill expertise"
  source_type: skill
  source: [list of skill ids the agent may cite, exactly one per prediction]
- name: "Authored flowcharts"
  source_type: flowchart
  source: <path or discovery procedure for the flowchart authoring>
- name: "Documented flowchart corpus"
  source_type: flowchart
  source: <path to the corpus doc>
```

Two methodologies or five — the machinery is identical. The MongoDB-flavored default
compares (1) raw MongoDB skill expertise, (2) an engineer's authored flowcharts, and
(3) a documented flowchart corpus; the canonical Phase-1 skill list lives in the
template.

## The workflow

Hand `references/orchestration-prompt.md` (instantiated with the parameters above) to
the orchestrating agent. It runs five stages:

1. **Stage shared inputs** (run once, serially):
   - Load `SUMMARY_FILE`; include **all** cases (assert the count); sort by recency —
     **parse dates to real Dates, never lexically sort `MM/DD/YYYY`** (it isn't
     chronological).
   - Derive each case's blind `initial_prompt` = title, optionally plus a
     *customer-authored* opening report. **Leakage guard:** only append a
     `troubleshooting_steps[0]` line if it reads as the customer's symptom report —
     never if it names a cause, a fix, or an engineer action. Those steps live in the
     *resolution* file; appending the wrong one leaks the answer into the "blind"
     input. When in doubt, title-only.

2. **Predict in parallel** — dispatch one prediction subagent per methodology, **all
   in a single message** with parallel `Agent` calls. Each agent reads only the
   shared case list + initial prompts, predicts a root cause per case (with a
   citation to its one source), writes records in batches, and **never opens the
   resolutions file and never grades.**

3. **Grade blind** — after predictions validate against the schema, the orchestrator
   pools every prediction into one shuffled, **anonymized** queue (phase + citation
   stripped, opaque `record_id`s, a private keymap the grader never sees), and
   dispatches **one** grader. The grader scores each record against the resolution
   using a fixed rubric, quoting the resolution to justify every grade. The
   orchestrator then de-anonymizes via the keymap and splits grades back per
   methodology.

4. **Synthesize** — aggregate scorecards (headline cohort = strong-resolution,
   non-title-only cases only; weak cohorts reported separately), disagreement
   analysis, threats to validity, and a caveated recommendation.

5. **Honesty checks** — predictions exist and validate before grading; predictor
   never touched the resolutions; grader never saw which methodology authored a
   record; mtimes order correctly.

See `references/schemas.md` for the exact JSON shapes (prediction record, grade
record, grading queue, keymap) and the required synthesis-report sections.

## The grading rubric (the measurement instrument)

Grades are **Correct / Partial / Wrong / Unverifiable**, decided on the failure
*mechanism* — the causal noun plus its qualifier — not on surface wording. The
template carries three worked calibration examples; reproduce them verbatim so the
grader applies one consistent standard. The two rules that keep it honest:

- Every grade quotes a verbatim line from `final_solution` or `troubleshooting_steps`.
- When torn between two grades, choose the lower one.

`Unverifiable` covers empty, unavailable, or autoclose-template resolutions — and
those cases are **kept out of the headline number**, because grading against a
non-resolution tells you nothing.

## Pitfalls this skill exists to prevent

- **Self-grading** — the single biggest validity killer. Keep predict and grade in
  separate agents.
- **Ground-truth leakage** into the blind prompt via resolution-derived fields.
- **Lexical date sorting** of `MM/DD/YYYY` (silently wrong ordering).
- **Silent truncation** — a single agent emitting hundreds of records can overflow
  context; batch with per-batch disk flushes.
- **Over-claiming** — this is a *relative* comparison under one grader against
  partly-synthetic ground truth, not an absolute accuracy measurement. Frame
  conclusions that way.

## Optional: land results on a scoreboard

If you're in a strategy-backtest scoreboard repo (schemas/runs/judges/scoreboard),
consider emitting predictions + grades as a `runs/` entry scored by the repo's
harness against a `judges/` artifact, so the result is reproducible, content-hashed,
and visible on the leaderboard rather than living as loose files. This is a larger
integration, not part of the core eval — offer it, don't assume it.
