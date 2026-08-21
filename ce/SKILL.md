---
name: ce
description: Slash-command alias for prompt-cost-estimator — run the cost estimator without spelling out the full skill name. TRIGGER when the user types "/ce" followed by a prompt, file path, or estimator flags. SKIP for anything that isn't invoking the estimator directly; for full flag documentation and guidance on presenting results, defer to prompt-cost-estimator.
version: 1.0.0
updated: 2026-07-05
model: claude-haiku-4-5
effort: low
---

# /ce — quick alias for prompt-cost-estimator

Thin wrapper. All logic lives in [[prompt-cost-estimator]] (`../prompt-cost-estimator/`) — this skill only exists so `/ce` is a one-token way to invoke it. Do not duplicate estimator logic here; if the underlying script changes, this file does not need to.

## What to do

1. Resolve the script path relative to this skill's own base directory (given when the skill loads): `<this-skill-dir>/../prompt-cost-estimator/scripts/estimate_cost.py`.
2. Parse the user's `/ce` arguments:
   - Free text after `/ce` with no recognized flags → treat as the prompt itself: `--text "<args>"`.
   - A path that exists on disk → `--file <path>`.
   - Recognized estimator flags (`--quick`, `--recommend`, `--exec`, `--confirm-exec`, `--target`, `--models`, `--json`, `--calls`, `--cached`, `--output-tokens`, `--multipliers`, `--offline`, `--tokens`) → pass through as-is.
   - No arguments at all → ask what prompt or file to estimate; don't guess.
3. Run the script with the resolved arguments.
4. Present results per prompt-cost-estimator's own "Presenting results" section — same rules apply here (lead with the answer, only pass `--exec` when explicitly requested, use `--confirm-exec` when unsure).

## Examples

- `/ce Summarize this repo's architecture` → `--text "Summarize this repo's architecture"`
- `/ce prompt.txt --quick` → `--file prompt.txt --quick`
- `/ce prompt.txt --recommend` → `--file prompt.txt --recommend`
- `/ce prompt.txt --exec` → `--file prompt.txt --exec` (spend-incurring — only when the user explicitly asked to run it, not just price it)

For anything beyond argument-passing — interpreting the matrix, cache/batch guidance, what `--recommend`'s tiers mean, exec safety — read `../prompt-cost-estimator/SKILL.md` rather than guessing from this file.
