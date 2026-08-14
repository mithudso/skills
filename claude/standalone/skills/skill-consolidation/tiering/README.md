# Skill Auto-Tiering (S3 Intelligent-Tiering for skills)

Keeps the always-on Claude Code skill index small by moving rarely-used skills *out* of it,
while auto-promoting any skill that gets actively used back *into* it — like S3 moving objects
between frequent- and infrequent-access tiers based on access patterns.

## Tiers

| Tier | Where it lives | In the index? |
| --- | --- | --- |
| **HOT** | `~/.claude/skills/<spoke>/SKILL.md` (standalone) | yes — name+description loaded every session |
| **COLD** | `~/.claude/skills/<hub>/references/<spoke>.md` only | no — loaded on demand when the hub routes to it |

Everything starts cold (folded under its hub by the consolidation). Content is never lost: the
reference copy, a dated tar in `../backups/`, and the mdb-context-hub registry all retain it.

## How access is detected

A `PostToolUse` hook (`log-access.mjs`, matcher `Read|Skill`) appends an entry to `access-log.jsonl` when:
- a hub `references/<spoke>.md` is **Read** (cold access — the hub routed to it), or
- a promoted spoke is invoked via the **Skill** tool (hot access).

The hook is fail-safe: any error is swallowed and it always exits 0, so it can never break a session.

## The engine (`tier.mjs`)

Runs at **SessionStart** (`--apply --quiet`) and on demand via **`/skill-tier`**:

1. **PROMOTE** a cold spoke → hot when its accesses within `windowDays` ≥ `promoteThreshold`
   (materializes the standalone `SKILL.md` from the reference, stripping the provenance banner so
   the original frontmatter is intact → it re-enters the index next session).
2. **DEMOTE** a hot spoke → cold after `demoteAfterDays` idle (refreshes the reference from the
   possibly-edited standalone to preserve drift, then removes the standalone dir).
3. **maxHot LRU eviction**: if the hot set exceeds `maxHot`, demote least-recently-accessed hot
   spokes until back at the cap — so promotion can never re-bloat the index past the budget.

A demote is **blocked** if the reference copy is missing (refuses to remove the only copy).

## Config (`tier-config.json`)

| Key | Meaning | Default |
| --- | --- | --- |
| `manifests` | which family manifests participate (spoke→hub→reference map) | `["data-analytics-manifest.json"]` |
| `promoteThreshold` | accesses within the window to promote | 2 |
| `windowDays` | rolling access window | 7 |
| `demoteAfterDays` | idle days before a hot spoke returns to cold | 14 |
| `maxHot` | LRU cap on promoted standalones (index-growth budget) | 10 |

Add `mongodb-manifest.json` / `writing-manifest.json` to `manifests` to tier those families too.

## Commands

```bash
node tier.mjs --status        # show HOT + recently-accessed spokes and the hot/maxHot budget
node tier.mjs                 # dry-run: preview promotions/demotions
node tier.mjs --apply         # perform tier changes
```

or `/skill-tier status | apply | dry`.

## Why SessionStart (not instant)

Claude Code reads the skill index once at startup, so a promotion materialized mid-session takes
effect at the **next** session — exactly like a storage tier change that settles asynchronously.
A truly hot skill is referenced again soon, so it re-promotes; the system self-corrects over sessions.

## Reversibility

Fully reversible. Promotion only creates a standalone from an existing reference; demotion only
removes a standalone whose reference is verified present. To freeze a spoke hot, raise its tier in
`tier-state.json` and exclude it from demotion, or just keep accessing it.
