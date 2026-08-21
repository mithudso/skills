---
name: skill-refresher
description: >-
  Freshness/currency pass for ONE existing skill — validate its claims against the current world, hunt post-authoring updates, apply dated fixes in place.
  Inventories atomic claims (versions, APIs, numbers, links, frontier statements), risk-ranks, verifies them against live primary sources with per-claim verdicts (CONFIRMED/STALE/WRONG/UNVERIFIABLE), diffs coverage vs releases/deprecations since AS-OF, edits in place, re-syncs hub copy.
  TRIGGER: "is this skill still accurate/current", "refresh/revalidate this skill", "audit for stale claims", "check against latest versions", "what changed since written", "/skill-refresher", periodic freshness maintenance.
  SKIP: structural/quality, not currency → skill-optimizer; MISSING concepts or new skills → concept-family-explorer; whole-tree rebalance → skill-tree-architect; research a topic, not a skill → /dr; registry timestamp/URL staleness → tam_staleness_scan; fact-check a customer doc → mongodb-technical-document-review.
whenToUse:
  - "is the <skill-name> skill still accurate?"
  - "refresh the <skill-name> skill"
  - "revalidate this skill's claims"
  - "audit this skill for stale or outdated content"
  - "check this skill against the latest versions/releases"
  - "what changed in this domain since the skill was written?"
  - "run /skill-refresher on <skill-name>"
  - "go through my installed skills and flag any that are out of date"
  - "this skill cites an old version, bring it current"
  - "verify the frontier claims in this skill"
version: 1.2.0
category: meta
updated: 2026-07-20
model: claude-opus-4-8
effort: xhigh
triggers:
  - skill-refresh
  - refresh skill
  - revalidate skill
  - skill still accurate
  - stale skill claims
  - skill freshness audit
  - skill out of date
  - check skill currency
keywords:
  - skill-refresher
  - freshness audit
  - claim validation
  - stale claims
  - currency check
  - coverage diff
  - update hunt
  - per-claim verdicts
related_skills:
  - skill-optimizer
  - concept-family-explorer
  - skill-tree-architect
  - deep-research
  - document-distiller
metadata:
  changelog:
    - "2026-07-20 v1.1.1->v1.2.0: post-gate remediation of blind-audit-2 dissent findings — /skill-refresh command token corrected to /skill-refresher everywhere; sweep-mode note added (per-target runs, tam_staleness_scan candidate pick) and batch whenToUse recast; hub repo named (mdb-context-hub) in Step 5.2; --max-claims cut line clarified as the only verification boundary. Applied after the blind re-audit gate closed (BLIND-AUDIT-DISSENT); deterministic verification only"
    - "2026-07-20 v1.1.0->v1.1.1: sko iteration 2 (blind re-audit) — Step 4 version bump now conditional on edits applied; all-CONFIRMED runs report verdicts without bumping or re-stamping updated"
    - "2026-07-20 v1.0.0->v1.1.0: sko iteration 1 — description trimmed 1147->975 chars (Glean cap); keyword concept-tree-diff->coverage-diff (overclaim); em-dash density 2.3->0.7/100w; --no-sync wired into Step 5; meta-validate exit-1 gate; source-disagreement precedence rule; evidence rule extended to Step 3 additions; self-pointer rewrite example; budget section deduped; effort high->xhigh (convergence-loop agentic tier)"
    - "2026-07-20 v1.0.0: initial release — claim inventory + risk ranking, live-source verification with CONFIRMED/STALE/WRONG/UNVERIFIABLE verdicts, post-AS-OF update hunt with concept-scale routing, evidence-gated in-place edits, hub local-sources re-sync, per-claim report"
---

# Skill Refresher

Validate one existing skill's factual content against the current state of the world and bring it current. Reads `SKILL.md` plus `references/*.md`, inventories checkable claims, verifies the risky ones against live primary sources, hunts for what appeared in the domain after the skill was written, applies evidence-backed edits in place, and re-syncs the hub copy. Owns factual **currency** only; structural quality (triggers, length, AI-isms, placement) stays with `skill-optimizer`.

## When not to use

- The skill file does not exist or cannot be located; report the failure and stop.
- The target has no world-facing claims (pure process/workflow skill with nothing version-pinned or dated); route to `skill-optimizer`; there is nothing for a freshness pass to verify.
- The target is a customer-facing or account document, not a skill; route to `mongodb-technical-document-review` / `tam-doc-validator`.
- The ask is for NEW coverage (missing sibling concepts, a new domain); route to `concept-family-explorer` or `/dr`.
- The caller wants only mechanical registry staleness (timestamps, dead URLs across catalogs); run `tam_staleness_scan` directly.
- No network access for verification; a refresh without live sources would be guesswork, so report and stop.

## Invocation

`/skill-refresher <skill-id-or-path>` or conversational ("is the kubernetes skill still current?").

Flags: `--max-claims <n>` (default 40), `--max-iter <n>` (default 3), `--report-only` (verdicts, no edits), `--no-sync` (skip hub re-sync).

Sweep mode: a multi-skill ask ("flag any installed skills that are out of date") runs this skill once per target; pick candidates via `tam_staleness_scan` (skills catalog) or the oldest frontmatter `updated` dates, then refresh each selected skill individually.

## Step 0 — Resolve target and AS-OF date

1. Resolve `~/.claude/skills/<id>/SKILL.md`; fallback to `tam_get_skill` `originalPath`. Load the body and every `references/*.md`.
2. Establish **AS-OF**: frontmatter `updated`, else the newest `metadata.changelog` date, else file mtime. Everything in the world newer than AS-OF is candidate drift.
3. Abort per "When not to use" if the target fails the gate.

## Step 1 — Claim inventory

Parse the body and references into atomic, independently checkable claims. Type each:

| Type | Meaning | Examples | Default risk |
|---|---|---|---|
| V | Version-pinned fact | "Vite 8 on Rolldown", model IDs, "GIPS 2020" | High |
| F | Frontier / state-of-world | "as of 2026", "X deprecated Y", "current default is Z" | High |
| T | Tool/API surface | flags, endpoints, function names, config keys | Medium |
| N | Numbers, limits, defaults | caps, quotas, tier limits, pricing | Medium |
| L | Links / URLs | doc links, repo links | Low (mechanical) |
| E | Evergreen | math, definitions, stable concepts | Skip |

Output a claim table: `id, claim, type, anchor (file:line), risk`. Verify at most `--max-claims` claims, highest risk first; note the cut line in the report.

## Step 2 — Verification

For every High/Medium claim, check against **current primary sources**: official docs, changelogs, release notes, vendor announcements. Batch related claims into shared searches (one changelog page usually settles many claims). Prefer one authoritative page over several secondary posts. When authoritative sources disagree, dated release notes and changelogs take precedence over undated or lagging doc pages; record both sources in the claim row.

Verdict per claim:

- **CONFIRMED** — still true; record source + source date.
- **STALE** — was true, superseded; record the replacement fact + source.
- **WRONG** — was never true; record the correction + source.
- **UNVERIFIABLE** — no authoritative source found; flag it, do not edit it.

L claims: HTTP reachability plus a spot check that the page still covers what the skill cites it for.

**Evidence rule (hard):** no STALE/WRONG verdict (and therefore no edit) without a citable source. The same bar covers Step 3: no queued addition without a citable source. Model memory alone is not evidence; the refresher must not inject its own staleness.

## Step 3 — Update hunt (post-AS-OF diff)

Existing claims aside, look for what's NEW since AS-OF: search `<topic> <current year>`, release notes and deprecation notices of every tool/product the skill names, renames of key terms. Route each finding by scale:

- **In-scope addition** — a new version, flag, sub-concept, or behavior that belongs inside this skill → queue an edit.
- **Sibling-scale gap** — a whole new adjacent family → do NOT build it here; report it as a handoff to `concept-family-explorer`.
- **Taxonomy impact** — a rename/merge/split that moves concepts across hubs → report as a note for `skill-tree-architect`.

## Step 4 — Apply edits

- Fix every WRONG; replace every STALE with the current fact, adding "as of <date>" on volatile values.
- Add in-scope additions as tight body or reference updates; match the skill's existing structure and voice.
- Touch the description only where it contains a stale claim; leave the TRIGGER/SKIP contract alone unless a renamed concept forces a term change.
- Only when edits were applied: bump the version (minor), append a `metadata.changelog` entry with verdict counts and what changed, set `updated` to today. An all-CONFIRMED run with no in-scope additions changes nothing — skip the bump and report the verdicts.
- Re-verify only the edited claims, then loop to Step 2 if edits raised new questions. Converge when zero STALE/WRONG remain; hard cap `--max-iter` iterations.

`--report-only` stops before this step and emits verdicts + proposed edits.

## Step 5 — Validate and sync

1. Lint: `node ~/.claude/skill-consolidation/meta-validate.mjs <id>`. A non-zero exit (any High lint finding) stops the run: report the lint findings and skip the sync.
2. Unless `--no-sync`: hub durability (sync direction is one-way FROM `local-sources/`): regenerate `local-sources/<id>/context.md` in the mdb-context-hub repo checkout (locate it via the skill's registry `sourceRepo` or the workspace's dev directory) (flat file, `references/*.md` embedded as appendices, self-pointers rewritten: a body link to `references/x.md` becomes "the x appendix below"), then patch `manifest.yaml` (integer version bump; category stays within the registry enum), then run `node scripts/sync-skill-pack.mjs` from the hub repo root. Registry-only writes get rolled back by the next pack sync; never stop at `tam_update_skill`.
3. If structural problems surfaced along the way (over-cap description, trigger collisions, bloat), recommend `/sko <id>` in the report; do not do sko's job here.

## Step 6 — Report

- Per-claim table: `id, claim, verdict, source, edit applied`.
- Additions applied; sibling-scale gaps handed to `concept-family-explorer`; taxonomy notes for `skill-tree-architect`.
- Files touched, versions bumped, sync outcome, claims left below the `--max-claims` cut line.

## Budget guardrails

- The `--max-claims` cut line is the only verification boundary: every claim above it gets the full Step 2 treatment (L-claims included); claims below it are reported as unverified, never silently skipped.
- The whole run should cost less than re-researching the domain: if more than half the verified claims come back STALE/WRONG, recommend a full `/dr` rebuild instead of patching.
