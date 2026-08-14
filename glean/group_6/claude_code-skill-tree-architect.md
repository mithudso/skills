# skill-tree-architect

**Category:** Science, Biology & Medicine
**Platform:** Claude Code
**Original Path:** claude-code/skill-tree-architect

## Description
Whole-tree architect for the ~/.Codex/skills hub-and-spoke taxonomy — audits the whole tree against the two-tier description cap (>1000 soft / >1536 hard), hub balance, and cross-hub concept placement, then orchestrates the existing skill-consolidation toolchain to rebalance and reshape it when a large new concept family lands (reuses the tools, never reimplements them). TRIGGER: rebalance / restructure / optimize the skill tree; is my taxonomy optimal; reshape for a new family; concept in the wrong hub or family; hubs over the cap / too big / should split; audit hub-and-spoke placement; find misplaced spokes; hub the homeless skills. SKIP: ONE skill's content or peer-seeds → skill-optimizer; a single NEW ≥8 homeless family → detect-candidates.mjs; what to build → concept-family-explorer; lint one skill → meta-validate.mjs; a prompt → prompt-deep-optimizer / phe.

---

# Skill Tree Architect

You operate at the **whole-tree** altitude: the shape and balance of the entire
`~/.Codex/skills` hub-and-spoke taxonomy. Per-skill quality belongs to `skill-optimizer`;
single-new-family detection belongs to `detect-candidates.mjs`. You **orchestrate the existing
toolchain in `~/.Codex/skill-consolidation/` — never reimplement it.**

**Read `~/.Codex/skill-consolidation/HUB-STRATEGY.md` first.** It is the canonical strategy
(the ≥8-sibling hub threshold, the build/maintain pipeline, reversibility, referent integrity).
This skill adds the one analytic that strategy doc's tools lacked — a tree-level placement +
cap-balance audit — and sequences every tool around it. Mutating phases follow the canonical
contract at `~/.Codex/skill-consolidation/convergence-and-severity.md` (snapshots, telemetry,
exit statuses, budget — cite it, never restate it).

## Phase 1 — ANALYZE (read-only; always safe)

Run all three and read them together; none mutates anything:

```bash
node ~/.Codex/skill-consolidation/audit-placement.mjs --desc-cap 1000   # cap headroom (two-tier), hub balance, MISPLACED spokes, homeless best-fit
node ~/.Codex/skill-consolidation/detect-candidates.mjs      # NEW ≥8 homeless families (dir-name clusters)
for d in ~/.Codex/skills/*/; do [ -d "$d/references" ] && node ~/.Codex/skill-consolidation/meta-validate.mjs "$(basename "$d")" --json; done
```

- `audit-placement.mjs` is the analytic this skill owns: per-hub description length vs the
  **two-tier cap** — >1000 chars = Medium (Glean export hard cap; single definition: sko Pass M),
  >1536 = High (harness truncation) — spoke counts, **misplaced spokes** (a spoke whose trigger
  vocabulary fits a sibling hub better than its current owner — heuristic), and homeless best-fit.
  `--json` for machine output.
- `meta-validate.mjs` is the **hard gate** (exit 1 on any High finding): spoke-copy-exists, dangling
  routing rows, circular SKIP, description cap, tier-config presence.
- `detect-candidates.mjs` only sees UNHUBBED skills by dir-name — it cannot see a misplaced *spoke*.
  `audit-placement.mjs` covers that blind spot; run both.
- **Concept-tree crosswalk (targeted, never full-tree):** for each spoke audit-placement flags
  MISPLACED — and any spoke you plan to re-file — pull its tree node via
  `tam_concept_tree_list(skillId:<id>)` or `tam_concept_tree_search` (a handful of calls, not a
  page-through). Treat tree-vs-hub disagreement as "one surface is stale — adjudicate", NOT as
  corroboration: the tree is currently the less-maintained surface (its parents still name
  pre-consolidation skills), so tree parentage must not boost a re-file's rank. A full
  dangling-skillId sweep of the tree is a one-time scripted reconciliation (offline join of the
  tree list vs manifests + installed dirs), not an every-run input.

End Phase 1 by capturing the **tree-health baseline**: `audit-placement.mjs --json` counts
(overCapHubs, misplacedSpokes, capPct distribution, homeless), `referents.mjs --status` dangling
count, and the meta-validate High/Medium aggregates. Phase 4 re-captures the same numbers — the
delta is the rebalance's evidence.

## Phase 2 — PLAN

Produce a ranked plan. Tag every item **safe-auto** or **review-required**:

| Finding | Action | Risk |
| --- | --- | --- |
| Dangling referent (`referents.mjs --status`) | `referents.mjs --repair --apply` | safe-auto (idempotent, ledgered) |
| Hub description **>1536** (High — harness truncation) | split by sub-domain, OR compress prose **keeping every spoke keyword** | review-required |
| Hub description **>1000** (Medium — Glean export hard cap; single definition: sko Pass M) | compress prose **keeping every spoke keyword** | review-required |
| Misplaced spoke, margin **≥ 0.25** (`audit-placement.mjs` suggestedScore − ownerScore) | re-file via the full `crossroute.mjs` sequence (Phase 3 step 2) | review-required |
| Misplaced spoke, margin 0.15–0.25 (the tool's default flag threshold) | watch | informational |
| Tree↔hub parentage disagreement | adjudicate; usual case is a stale TREE node → re-upsert `parentConcept` + `tam_concept_tree_link` to the actual hub; re-file via `crossroute.mjs` only when audit-placement independently flags the spoke | review-required |
| New ≥8 homeless family | full consolidation per `auto-hub-agent-prompt.md` | review-required |
| Hub with very high spoke count, cap OK | watch / split later | informational |

**Never drop spoke trigger enumeration to save characters — the enumeration IS the routing signal**
(only the hub name + description are in-context until the hub is chosen). Fix over-cap by splitting
or compressing prose, never by deleting keywords.

## Phase 3 — APPLY (only on explicit go; `~/.Codex/skills` is NOT git-backed)

Drive the existing tools in order; do not hand-roll any step. Mutating tools default to dry-run —
pass `--apply` explicitly (`referents.mjs --repair`, `tier.mjs`); snapshot before any write per
the canonical contract:

1. **Snapshot (canonical backups convention):** copy EVERY file the run will modify into
   `~/.Codex/skill-consolidation/backups/<run>/` — `tar` the affected family PLUS the
   out-of-family top-level SKILL.mds a `referents.mjs --repair` dry-run says it would touch.
   (`skill-pack.config.mjs` is git-backed in mdb-context-hub — step 8 covers it; not in the
   snapshot set.)
2. Edit `<family>-mapping.json` → `node build.mjs <family>-mapping.json` (folds spokes into
   `<hub>/references/`, writes routing + manifest). For a single re-file, `crossroute.mjs` adds
   the spoke to the NEW hub **only** — always complete the sequence: remove the spoke's row from
   the old hub's family manifest, delete the old routing-table row, move the old
   `references/<spoke>.md` aside (or banner it superseded), then `node referents.mjs --repair
   --apply`. Unresolved double-ownership makes spoke→owningHub resolution ambiguous —
   hub-registry.mjs already has to special-case consolidation-manifest.json as "a known
   duplicate". After any re-file, re-upsert + `tam_concept_tree_link` the spoke's tree node to
   the new hub.
3. Author/adjust the hub `SKILL.md` (absorb spoke TRIGGER vocab; SKIP → sibling hubs only).
4. Fold/remove a spoke dir **only after verifying its `references/<spoke>.md` copy exists** — and
   if `tiering/tier-state.json` says the spoke is HOT, route the removal through
   `tiering/tier.mjs --demote <spoke> --apply` (targeted demote; dry-run first by omitting
   `--apply`), NEVER raw `rm`: tier.mjs does the drift-preserving
   one-way sync (standalone → `references/<spoke>.md`) before removal and chains
   `referents.mjs --repair`. A manual reconcile (diff banner-stripped content via
   `tiering/lib.mjs` stripBanner) is needed only on two-sided divergence — the hub reference was
   ALSO edited since promotion, which tier.mjs would blindly overwrite. Live anchor: as of
   2026-06-11 deep-research-methods is HOT with 27 diverged lines, standalone newer — a
   copy-exists-only check plus raw `rm` would have destroyed them.
5. `node fix-crosshub-generic.mjs <family>-manifest.json` (provenance banners + cross-hub map).
6. `node referents.mjs --repair --apply` (re-point `related_skills:` and inline `→ <id>` peer seeds).
7. Register the family in `tiering/tier-config.json` (`meta-validate.mjs --register-tier --apply`).
8. Add new hub ids to `SELECTED_SKILLS` in
   `~/Documents/GitHub/mdb-context-hub/scripts/skill-pack.config.mjs`, then `npm run sync:skills`
   (regenerates the live registry; review-required — different repo, user commits).

## Phase 4 — VERIFY

1. Re-run the Phase-1 trio; confirm **zero High** `meta-validate.mjs` findings and that the issue
   you targeted is gone (re-filed spoke now scores to its new hub; over-cap hub now under 1536).
   Re-capture the Phase-1 tree-health numbers and report the delta — that delta is the rebalance's
   evidence. Append the run's tree-health rows to `optimizer-telemetry.jsonl`
   (`artifact_type: "skill-tree"`, `pass: "tree-health"`) per the canonical telemetry schema.
2. **Scored routing-probe replay** (the documented gate link-checks alone don't satisfy): replay
   the touched family's persisted CQ corpus (`~/.Codex/skill-consolidation/evals/<family>.cq.md`)
   plus 1 new probe this run — dispatch each as a fresh-context subagent question with NO hint
   about hubs/references; it must reach the right `references/*.md` by routing alone. If no
   `.cq.md` exists for the family, invent the usual 3 probes (an in-hub topic, a moved topic, a
   cross-hub topic) AND persist them as that family's seed corpus (see `evals/README.md`).
   **Gate:** reached-correct-reference ≥ 80%; record no-hub-backtrack as an informational
   directness column, not a gate. Bounded fail path: failed probe → one fix iteration → one
   re-probe → exit `PROBE-DISSENT` listing the failures. Append probe outcomes to the family's
   evals rows.
3. **Exit status (canonical contract):** `CLEAN` only after the probe gate passes; `BLOCKED` for
   findings that cannot be fixed without inventing content; `PROBE-DISSENT` per step 2.
4. End the Phase 4 report with the **literal restore command** for this run's snapshot (e.g.
   `tar -xzf ~/.Codex/skill-consolidation/backups/<run>/<family>.tar.gz -C ~/.Codex/skills` plus
   copying back the snapshotted out-of-family SKILL.mds) and the rollback procedure: snapshot
   restore, then reverse this run's `referents-ledger.json` entries.

## Hand-offs (peer deferral)

- One skill's content/triggers/length/peer-seeds → **skill-optimizer** (it runs `referents.mjs` after tier changes).
- Which new skills a domain needs → **concept-family-explorer** — including any tree concepts found skill-less while adjudicating.
- Skill anatomy / frontmatter / description-cap rules → **Codex-skills**.
- This skill is in the shared EXCLUDE_LIST (`~/.Codex/skill-consolidation/exclude-list.mjs`, imported by both `detect-candidates.mjs` and `audit-placement.mjs`) — a meta tool is never itself hubbed.