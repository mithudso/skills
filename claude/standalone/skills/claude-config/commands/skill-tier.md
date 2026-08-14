---
description: Run or inspect the skill auto-tiering engine (promote hot skills into the index, demote idle ones back under their hub)
argument-hint: "status | apply | dry (default dry-run)"
---

Run the skill auto-tiering engine. Argument `$ARGUMENTS`:

- `status` (or empty + you want a view): `node ~/.claude/skill-consolidation/tiering/tier.mjs --status`
- `apply`: `node ~/.claude/skill-consolidation/tiering/tier.mjs --apply` — actually promote/demote
- anything else / `dry`: `node ~/.claude/skill-consolidation/tiering/tier.mjs` — dry-run preview

Tiering is S3-Intelligent-Tiering for skills:
- A cold spoke (a `references/<spoke>.md` under a hub, OUT of the index) is **promoted** to a standalone `~/.claude/skills/<spoke>/SKILL.md` (back IN the index) once it has been accessed `promoteThreshold` times within `windowDays`.
- A hot standalone is **demoted** back to a reference after `demoteAfterDays` idle, or evicted LRU when the hot set exceeds `maxHot`.
- Access is logged automatically by the PostToolUse hook (`log-access.mjs`) whenever a hub reference is Read or a tiered spoke is invoked via the Skill tool, and the engine runs at SessionStart.

Config + state: `~/.claude/skill-consolidation/tiering/tier-config.json`, `tier-state.json`, `access-log.jsonl`. Full rules: `~/.claude/skill-consolidation/tiering/README.md`.

Run the command matching `$ARGUMENTS`, then summarize the resulting tier table.
