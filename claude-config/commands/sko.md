---
description: Audit and improve a skill to production quality via the skill-optimizer convergence-loop quality gate, then sync it to the mdb-context-hub
argument-hint: <skill-id-or-path> [--meta] [--max-iter=N] [--no-sync] [--budget-minutes=N] [--eval=measured] [--cross-model] [--sync-anyway]
---

Run the **skill-optimizer** skill on the target below.

Invoke the skill via the Skill tool and follow its SKILL.md procedure exactly (locate → baseline snapshot → convergence loop of analytical passes → triage → apply fixes → post-write verification → hub sync → report). The skill's SKILL.md is the single source of truth; do not re-specify the steps here.

Target and flags:

$ARGUMENTS

If $ARGUMENTS is empty, ask once: "Which skill should I optimize?", then continue.
