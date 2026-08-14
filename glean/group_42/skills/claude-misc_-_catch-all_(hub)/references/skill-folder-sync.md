<!-- hub-reference-banner -->
> **Reference file — part of the `misc-catch-all` hub.** Formerly the standalone `skill-folder-sync` skill.
> Sibling topics in this family are now reference files under the hubs (`misc-catch-all`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: skill-folder-sync
description: >-
  Compare two .claude/skills folders and one-way sync them, hub-and-spoke aware.
  Copies skills present in the source tree but missing from the target, and
  reconciles fold-state so the target matches the source: a folded spoke (a skill
  at <hub>/references/<spoke>/SKILL.md) goes INTO the right hub, never at top
  level; a missing hub is copied whole so its spokes ride along; and when a skill
  is folded on one side but loose on the other, the source wins.
  TRIGGER: sync, mirror, reconcile, merge, or copy-missing skills between two
  .claude/skills dirs (local vs network mirror, laptop vs desktop, backup vs
  live), "copy skills I'm missing from that folder", "make this folder match that
  one", "move folded skills into the right hubs", "fix skills folded here but
  loose there". Prefer over plain cp/rsync, which misplace spokes. SKIP: pushing
  skills to the context-hub registry (use /sync-skills); editing one skill's
  content (use skill-optimizer); generic file copy unrelated to a skills tree.
---

# Skill folder sync (hub-and-spoke aware)

One-way sync from a **source** skills tree to a **target** skills tree. Source is
authoritative. The target only ever **gains** skills or has **fold-state
conflicts fixed** — non-conflicting skills already in the target are never
deleted. This is copy-missing + fix-fold-state, never a destructive mirror.

## The two skill locations (why a plain copy is wrong)

A skill's identity is its **directory name**. The same skill can live in two shapes:

- **Top-level (unfolded):** `<root>/<name>/SKILL.md`
- **Folded spoke:** `<root>/<hub>/references/<name>/SKILL.md` — the full former
  standalone skill, tucked under a hub's `references/`, usually marked with a
  leading `<!-- FOLDED SPOKE of the <hub> hub ... -->` comment.

A naive `cp -r` or `rsync` of a source tree into a target will drop a folded
spoke at top level (wrong — it belongs inside its hub) or duplicate a skill that
exists in the other shape. The bundled script understands both shapes, so it
places each skill where the source says it belongs.

## What the script decides, per skill

For every skill in the source:

1. **Missing from target →** copy it.
   - Folded spoke → copied into `target/<hub>/references/<name>/`.
   - Top-level skill → copied into `target/<name>/`.
   - If the skill's **entire hub** is missing from the target, the hub is copied
     as one unit and its folded spokes ride along inside it (the script skips
     copying those spokes individually to avoid double work).
2. **Present but fold-state differs →** reconcile (source wins): delete the
   target's wrong-state copy, then place the source's version in the correct
   location.
   - Source folds it under a hub, target has it top-level → delete the loose
     top-level copy, fold it into the hub.
   - Source has it top-level, target folds it under a hub → delete the folded
     copy, place it at top level.
   - Both fold it but under **different hubs** → move it to the source's hub.
3. **Present, same fold-state →** skip (already in sync).

## How to run it

The script defaults to a **dry-run** so the user can see the plan before anything
changes. Always show the plan first. `scripts/sync_skills.py` lives in this
skill's own directory (the folder containing this SKILL.md) — use the absolute
path to it shown in the skill's base-directory load message, so it runs no matter
the current working directory.

```bash
# 1. Preview (dry-run) — copies nothing, prints the COPY / RECONCILE / SKIP plan
python3 <skill-dir>/scripts/sync_skills.py <SOURCE_skills_root> <TARGET_skills_root>

# 2. Execute once the user is happy with the plan
python3 <skill-dir>/scripts/sync_skills.py <SOURCE_skills_root> <TARGET_skills_root> --apply
```

`<...skills_root>` is the directory that *contains* the skill folders (e.g.
`/Users/me/.claude/skills`), not an individual skill. `~` is expanded.

**Direction matters.** The first path is the authoritative source you copy
*from*; the second is the target you copy *into*. To make a true union of two
trees (each gains what the other has), run it once each way — but be aware that
fold-state conflicts resolve toward whichever tree you named as source that run,
so pick the side whose hub organization you trust as source.

## Workflow when the user asks for a sync

1. Confirm which folder is the **source of truth** and which is the **target**.
   If the user only says "sync these two," ask, because fold-state conflicts are
   resolved in the source's favor and the direction is not symmetric.
2. Run the **dry-run** and relay the plan: how many copies, how many fold
   reconciles, and call out every RECONCILE line explicitly (those are the only
   ones that delete something in the target).
3. Get a go-ahead, then run with `--apply`.
4. Re-run the dry-run to confirm it now reports nothing to do (the operation is
   idempotent).

## Notes & edge cases

- Dot-directories (e.g. `.sync-backup-*`) are ignored on both sides.
- Copies are atomic-ish per skill: an existing dest dir is removed and replaced,
  so a copy always reflects the source exactly.
- This pairs with the existing rsync mirror habit (`rsync -a` to the network
  mirror, no `--delete`): use rsync for a flat byte-for-byte mirror, use this
  when the two trees disagree on **fold-state** and you need spokes placed inside
  the right hubs.
