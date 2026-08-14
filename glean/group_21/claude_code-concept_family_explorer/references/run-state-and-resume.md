# Run-state checkpoint & resume (Step write points + re-entry algorithm)

cfe is the family's longest-running orchestrator (a multi-batch saturation
loop spawning many `/dr` runs); without a checkpoint, an interrupted Step 8
loop restarts at Step 1 and repays the full Steps 1–4 scoring cost.

**Schema is defined once** in
`~/.claude/skill-consolidation/run-state/README.md` (`cfe-run-state/v1`:
subject, budget as an opaque object, `round`, `scored`, `researched` with
per-entry `persisted` flags, `skipped`, `failed`, `queue`, `status`) — cite
it, never restate it. The checkpoint is mutable working state; the JSONL
telemetry ledgers are the immutable audit. Plain file writes, no scripts.

## Write points

| When | Write |
| --- | --- |
| End of Step 4 | initial `scored` table |
| End of Step 5 | `queue` (above-threshold selection) |
| After each `/dr` return + its Step 6b persistence | move the concept `queue` → `researched`; set `persisted: true` only after persistence succeeds (honors the budget rule: never stop mid-persist) |
| End of each Step 7 round | `round`++, merge new `scored` rows |
| Step 10 | `status: "complete:<verdict>"` — the file is then evidence, not a lock |

**dryRun exception:** a plan-only `dryRun` (stops at Step 5; SKILL.md Inputs)
writes **no** checkpoint — the Step 4/5 write points above apply only to a real
run — so a preview is never mislabelled `in-progress` and offered for resume.

Every write also refreshes `updated`. Writes are fail-safe per the canonical
rule (`~/.claude/skill-consolidation/convergence-and-severity.md` §Telemetry:
a write error never blocks or fails a run) — log the failure in the report and
continue.

## Re-entry algorithm

On invocation, check for
`~/.claude/skill-consolidation/run-state/cfe-<subject-slug>.json`:

1. **No file, or `status` starts with `complete:`** → fresh run; a complete
   checkpoint is evidence for the report, not a lock.
2. **`status: "in-progress"`** → ask **resume-vs-fresh ONCE**, showing file
   age, stored budget, and queue size. Never auto-resume — a deliberately
   abandoned run must not silently restart.
3. **On resume:**
   - Skip Steps 1–5; use the stored `scored` table.
   - Re-run the Step 2 inventory ONLY for `queue` items, to catch coverage
     added since the checkpoint (a queued concept another run filled is
     re-tagged HAVE and dropped from the queue, logged as a skip).
   - `researched` entries with `persisted: false` re-run **Step 6b only** —
     persist-spoke is idempotent (dr.md Phase 3: re-running is safe). This
     closes the silent-wipe window where a researched-but-unpinned skill
     vanishes on the next sync yet reads as HAVE.
   - Enter Step 6 at the queue head.
   - Remaining budget = counts only: `maxConcepts` − (researched + failed),
     `maxRounds` − `round`. If the stored budget carries a time cap
     (`budgetMinutes`), surface it to the user — do not reconcile elapsed
     wall-clock across sessions.
4. **On fresh (user chose fresh):** overwrite the checkpoint at the first
   Step 4 write.
