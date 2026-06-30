# Skill Hub-and-Spoke Strategy

How large concept families are kept out of the always-on skill index without losing functionality.

## The problem

Claude Code injects **every installed skill's `name` + `description`** into the system prompt at the
start of *every* session. With 400+ skills that is large, always-on context cost plus decision noise —
and most of it is irrelevant to any given session (a frontend task does not need 73 MongoDB descriptions
in context).

## The pattern

A **hub** is one top-level skill that owns a concept family. Each former standalone "spoke" skill becomes
an on-demand **reference file** inside the hub:

```
<hub>/
  SKILL.md            # frontmatter description (absorbs all spoke triggers) + routing table + cross-hub map
  references/
    <spoke>.md        # the spoke's full content, verbatim, led by a provenance banner
```

- **Index cost:** one entry per hub instead of one per spoke.
- **Functionality:** zero loss — when the hub activates, its routing table points to the exact
  `references/<spoke>.md` file, which the model `Read`s on demand (progressive disclosure).
- **Discoverability across hubs:** every hub carries a **cross-hub map** so a topic owned by a sibling
  hub is reachable; reference files carry a **provenance banner** that neutralizes stale "use the X skill"
  pointers baked into the verbatim copy.

## When to create / route to a hub

| Situation | Action |
| --- | --- |
| Family has **≥ ~8 sibling skills** (or is growing toward it) | Consolidate into a hub (or a few sub-domain hubs) |
| New `/dr` topic belongs to an **existing hub's domain** | Install as `references/<name>.md` under that hub — never a new top-level skill |
| New topic **starts** a large family (≥8 expected siblings) | Create it AS a hub from the start |
| Genuinely standalone single-topic skill, no family | Install as a normal top-level skill |
| Sub-domain within a family is large/distinct | Give it its own hub (fold *by sub-domain*, not one mega-hub) |

Keep distinct-mechanism skills standalone even within a family (e.g. `mongodb-kb` is a KB-article index;
`10gen` is repo intelligence — neither is knowledge depth, so they stay top-level).

## Registered hubs

- **MongoDB** (77 → 6): `mongodb-expert` (core data plane & engine internals), `mongodb-atlas-expert`
  (Atlas cloud platform), `atlas-diagnostics-expert` (live diagnostics/perf/monitoring),
  `mongodb-operations-expert` (ops, reliability & data movement); `mongodb-kb` + `10gen` standalone.
- **Data analytics** (`da-*`): consolidated into sub-domain hubs — see `da-mapping.json`.
- **Writing**: consolidated into sub-domain hubs — see `writing-mapping.json`.

## How to build/maintain (tooling in this directory)

1. **Map**: write `<family>-mapping.json` — `{ hubs: { <hub>: { keepExisting, title, spokes: [...] } }, standalone: [...] }`.
2. **Validate + back up**: confirm every family dir is a hub, a mapped spoke, or standalone (nothing unmapped),
   then `tar` the family before any removal. (Content is triple-preserved: `references/` + tar + hub registry.)
3. **Build**: `node build.mjs <family>-mapping.json` — copies each spoke verbatim into `<hub>/references/`,
   emits a routing-table fragment per hub, writes `consolidation-manifest.json`.
4. **Finalize hub `SKILL.md`** (judgment): rewrite each hub `description` to absorb every spoke's TRIGGER
   vocabulary; point `SKIP`/`whenNotToUse`/`related_skills` only at *sibling hubs* (never at absorbed spokes);
   insert the routing table; author new hubs from scratch.
5. **Remove spokes**: delete the spoke dirs from `~/.claude/skills/` only after verifying each `references/`
   copy exists.
6. **Cross-hub wiring**: `node fix-crosshub.mjs` (generalize per family) — prepend a provenance banner to every
   reference file and append a cross-hub map to every hub `SKILL.md`.
7. **Sync**: `tam_create_skill` / `tam_update_skill` for new/changed hubs; add hub IDs to
   `SELECTED_SKILLS` in the mdb-context-hub so `npm run sync:skills` keeps them.
8. **Test (the gate)**: dispatch blind functional probes — a subagent given a realistic deep question with
   NO hint about hubs/references must reach the right `references/*.md` purely by following routing. Test a
   direct in-hub topic, a new-hub topic, and a cross-hub topic. Link-integrity checks are necessary but not
   sufficient — only a functional probe proves routing works.

## `/hub-detect` — auto-detection & the weekly scheduler

A scheduler watches for a **new** large family that has grown to hub-worthy size and consolidates it
automatically. Three pieces, all in this directory:

1. **`detect-candidates.mjs`** (read-only) — lists every top-level `~/.claude/skills/*/` dir that is
   NOT already a registered hub, a meta/infra tool (`EXCLUDE_LIST`), or a folded spoke under any
   `*-manifest.json`; clusters the remainder into theme buckets (security, document-formats, tam-ops,
   testing, data-integration-extraction, catch-all — the `CATEGORIES` table is easy to extend); and
   reports any bucket with **≥ 8** members that is not already a family/hub as a **ready candidate**.

   ```bash
   node ~/.claude/skill-consolidation/detect-candidates.mjs          # human summary (default)
   node ~/.claude/skill-consolidation/detect-candidates.mjs --json   # machine-readable
   ```

   This script **never mutates** anything. Run it any time to see candidates manually (`/hub-detect`).

2. **`auto-hub-agent-prompt.md`** — the standing prompt for the consolidation agent. On a ready
   candidate it applies a judgment filter (drop hubs / meta tools / distinct-mechanism standalones;
   abort if < 8 spokes survive), then runs the full pipeline: mapping → `build.mjs` → author hub
   `SKILL.md` → tar+fold spokes → `fix-crosshub-generic.mjs` → register the manifest in
   `tiering/tier-config.json` → `referents.mjs --repair --apply` (if present) → sync → verify.

3. **The schedule** — job name **`skill-hub-autodetect`**, **weekly (Mondays 08:13 local)**, via two
   registration paths:
   - **Durable (survives restarts):** a system crontab entry runs `run-autohub.sh`, which executes the
     detector read-only and only hands off to `claude -p` (with the standing prompt) when a real ≥8
     candidate exists. Permission gates are NOT bypassed. Edit/inspect with `crontab -l` / `crontab -e`.
     Log: `autohub.log` in this directory.
   - **In-session convenience:** a CronCreate job (same cadence) registered while a Claude session is
     live. CronCreate jobs are session-scoped and auto-expire after 7 days — the crontab entry is the
     durable one.

Detection is read-only; the scheduled job is the only thing that mutates, and only when a real ≥8
unhubbed candidate survives the filter.

## Referent integrity (keeping cross-skill pointers resolvable)

The `skill-optimizer` seeds two kinds of cross-skill routing referents into skills:

1. **Frontmatter `related_skills:`** — a YAML list of skill ids (block or inline `[a, b]` form).
2. **Inline deferral tokens** in the `description`/body of the form `→ <skill-id>` (often after `SKIP:`).

When a spoke is folded into a hub (becomes a cold `references/<spoke>.md`, its top-level dir removed) or
promoted/demoted by the tiering engine, any referent that names that spoke now points at a skill id that
is no longer top-level (indexed) → routing breaks. `referents.mjs` keeps these pointers resolvable.

**Tool:** `node referents.mjs --repair [--apply] [--quiet] | --status`

- Builds a `spoke → owningHub` map from every `*-manifest.json`. A spoke is **COLD** if
  `~/.claude/skills/<spoke>/SKILL.md` is absent (folded), **HOT** if present (promoted/standalone).
- `--repair` scans every **top-level** `SKILL.md` (hubs + standalones). It never touches files under
  `references/` — those carry the provenance banner that already neutralizes stale pointers. For each
  referent naming a tiered spoke:
  - `related_skills` entry → COLD spoke: replace the entry with its owning hub (dedup; drop self-refs/exact dups).
  - inline `→ <spoke>` (cold): rewrite to `→ <hub> (references/<spoke>.md)`.
  - **Reverse (idempotent):** if a referent points at a hub for a spoke that is now **HOT**, restore the
    bare `<spoke>` id using the ledger — so repeated repairs stay stable as spokes promote/demote.
- Every rewrite is recorded in `referents-ledger.json` (`{file, surface, from, to, spoke, hub}`) so it is
  reversible and so reverse-restore knows what to undo.
- `--status` counts dangling referents (cold spokes still named directly) without writing.
- Default is **dry-run**; pass `--apply` to write. Safe and idempotent: only referent tokens / list
  entries are rewritten — never reference files, hub routing tables, or non-referent content.

**Wired into tiering:** `tiering/tier.mjs` runs `referents.mjs --repair --apply --quiet` automatically
after any `--apply` promote/demote, so every tier change keeps referents pointing at a resolvable location.

## `/dr` integration

The hub-routing decision is now Phase 2 step 3 of the `/dr` workflow (both `~/.claude/commands/dr.md` and
the hub prompt `prompts/saved/create-deep-research-skill-from-web-lookups.md`), so new research lands as a
hub reference instead of growing the index.

## Reversibility

Every consolidation keeps: the `references/` copies, a dated `tar` in `backups/`, the
`consolidation-manifest.json` spoke→hub→file map, and the spokes' entries in the mdb-context-hub registry.
To restore a spoke as standalone: copy its `references/<name>.md` body back to `~/.claude/skills/<name>/SKILL.md`
(strip the banner), or extract from the tar.
