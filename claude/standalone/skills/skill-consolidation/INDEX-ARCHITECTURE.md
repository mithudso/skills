# Skill-Library Indexing System — Technical Architecture

## In brief

This system is a routing index for a large library of AI *skills* — small instruction
files an agent picks from to carry out a task. The library holds ~578 skills (125
standalone skills plus 453 "folded" spokes, grouped into 30 hub families). Without an
index, an agent has to open dozens of skill files just to decide which one to use. The
generator (`gen-skills-index.mjs`) reads them all once and writes a single compact
routing table (`SKILLS-INDEX.md`) plus a full machine-readable index
(`SKILLS-INDEX.json`). The agent scans the one table, picks the right skill, and only
then opens that skill's file for the full detail. The index is always derived from the
skills themselves and is regenerated automatically when skills change, so it never
becomes a second copy to maintain by hand.

On top of the plain table, there is an **optional semantic-search layer**. The default
router matches on literal trigger keywords, which misses requests phrased in different
words ("the borrow checker keeps rejecting my code" never says *Rust*). The semantic
layer embeds every skill with a local, free, no-API-key model (`qwen3-embedding:4b` via
Ollama) and answers a free-text query (`--search "..."`) by meaning, blended with
keyword overlap (a *hybrid* score). It is purely additive — opt-in, never required, and
on any failure it falls back to the keyword router — and the embeddings are refreshed
automatically by the same background agent that rebuilds the index.

**What it saves**

- **Cost (tokens).** Choosing a skill from the one-page table instead of reading the
  whole library cuts roughly **96% of the tokens** (~577K saved on a full scan). A
  typical disambiguation saves **~20K–40K input tokens** — about **$0.09–$0.45 per
  selection**, or roughly **$450–$2,250 a year** at 5,000 selections.
- **Time.** About **10–30 seconds saved per non-trivial selection**, from fewer file
  reads plus a much smaller prompt for the model to process.
- **Accuracy.** An optional local semantic-search layer catches paraphrased requests
  the plain keyword router misses — **70% top-1 / 87% top-3** hit rate, versus
  **47% / 50%** for keyword matching alone.

The rest of this document is the technical architecture behind those numbers.

---

> Status: living document. Describes `gen-skills-index.mjs` and the consolidated
> `SKILLS-INDEX.{json,md}` artifacts, their integration hooks, and operations.
> Last verified against the generator on 2026-06-30 (578 skills, 30 families, schemaVersion 1).
> Optional semantic layer verified the same day: 578 vectors, `qwen3-embedding:4b`, 2560-d, hybrid
> (cosine + keyword) search by default and **fail-open to keyword-only when the server is down**;
> auto-refreshed by the `com.mitch.skills-embed` launchd agent, covered by
> `gen-skills-index.test.mjs` (14/14) and A/B-benchmarked across five routing methods by
> `bench-search.mjs`.

## Contents

- [How-To Guide](#how-to-guide)
1. [Design rationale](#1-design-rationale)
2. [Implementation architecture](#2-implementation-architecture)
3. [Synchronization & drift control](#3-synchronization--drift-control)
4. [Performance & token metrics](#4-performance--token-metrics)
5. [Operational guidance](#5-operational-guidance)
6. [Semantic retrieval architecture](#6-semantic-retrieval-architecture)
7. [Cluster → agent overlay analysis](#7-cluster--agent-overlay-analysis)
8. [Appendix: file map & sources](#8-appendix-file-map--sources)
9. [System Bootstrap Prompt](#9-system-bootstrap-prompt)

---

## How-To Guide

A practical walkthrough of the whole system: what the pieces are, how to use them day to
day, and how to keep them fresh. Everything lives in `~/.claude/skill-consolidation/` and runs
on Node.js with **zero npm dependencies**.

### The mental model

| Artifact | What it is | Who reads it |
| --- | --- | --- |
| `~/.claude/skills/<id>/SKILL.md` | The skills themselves (source of truth) | you / the agent, for full detail |
| `SKILLS-INDEX.md` | One-page routing table, grouped by family | the agent, to pick **one** skill |
| `SKILLS-INDEX.json` | Full machine index (every field) | tooling, `--search`, CI |
| `SKILLS-EMBEDDINGS.json` | One vector per skill (optional semantic layer) | `--search` only |
| `gen-skills-index.mjs` | The generator + search CLI | you, via `node` |

The index is **always derived** from the skills — never hand-edited. Regenerate it; don't patch it.

### 1. Everyday use (no server, no setup)

Pick a skill by scanning the routing table:

```bash
# regenerate the table from the live library, then read it
node gen-skills-index.mjs
open SKILLS-INDEX.md          # or just read it in your editor
```

`SKILLS-INDEX.md` is grouped by family; each row carries the skill's summary, top triggers,
and peer links. That keyword router needs nothing else — no server, no model.

### 2. Semantic search (optional layer)

When a request is phrased in words the triggers don't contain ("the borrow checker keeps
rejecting my code"), rank by *meaning*:

```bash
node gen-skills-index.mjs --search "the borrow checker rejects my lifetimes"
# 0.695  lang-rust  —  Rust language sub-hub ...
```

Useful flags: `--top=N` (default 10), `--threshold=T` (default 0.30), `--alpha=A`
(hybrid blend, default 0.92), `--no-hybrid` (pure cosine), `--json` (machine output).

**It fails open.** If the embedding server is down or the corpus is missing, `--search`
prints a one-line `degraded to keyword-only ranking` note to stderr and returns
keyword-coverage results with **exit 0** — it never blocks you. Semantic ranking resumes
automatically once the server is back.

### 3. Enabling the semantic layer

The layer needs a local Ollama server and the pinned embedding model (free, no API key):

```bash
brew install ollama
brew services start ollama          # server that persists across reboot/login
ollama pull qwen3-embedding:4b      # 2560-dim model, pulled once
node gen-skills-index.mjs --embed   # build SKILLS-EMBEDDINGS.json (incremental)
```

`--embed` is incremental: it re-embeds only skills whose source hash changed and drops
vectors for deleted skills, so routine runs are cheap.

### 4. Keeping it fresh automatically (persistence)

Two macOS services keep everything current with zero touch:

| Service | Role | Check it |
| --- | --- | --- |
| `ollama` (Homebrew) | embedding server, owns port 11434 | `brew services list \| grep ollama` |
| `com.mitch.skills-embed` (launchd) | rebuilds index + embeddings | `launchctl list \| grep skills-embed` |

The launchd agent (`~/Library/LaunchAgents/com.mitch.skills-embed.plist`) runs
`embed-refresh.sh` on a **6-hour interval** and **immediately whenever `~/.claude/skills`
changes** (WatchPaths). Its log is `embed-refresh.log`. Net effect: add or edit a skill, and
the index + vectors reconcile themselves within seconds.

### 5. After you change the library (manual reconcile)

If you don't want to wait for the agent — or you're on a machine without it:

```bash
node gen-skills-index.mjs           # rebuild the index from the current skills
node gen-skills-index.mjs --check   # exit 0 = up to date, exit 1 = stale (for CI)
node gen-skills-index.mjs --embed   # re-sync vectors (needs the server up)
```

### 6. Verify it works

```bash
node --test gen-skills-index.test.mjs   # 14/14; live tests self-skip if the server is down
```

### Exit codes & environment

| Exit | Meaning |
| --- | --- |
| `0` | success (results returned, or index up to date) |
| `1` | stale index (`--check`), or no match for a query |
| `2` | usage error (empty query) or a model/dimension mismatch |

| Env var | Default | Purpose |
| --- | --- | --- |
| `OLLAMA_HOST` | `http://localhost:11434` | embedding server URL |
| `SKILLS_EMBED_MODEL` | `qwen3-embedding:4b` | pinned model (the corpus is locked to it) |
| `SKILLS_ROOT_DIR` | `~/.claude/skills` | skills source directory |
| `SKILLS_OUT_DIR` | the repo dir | where artifacts are written |

### Troubleshooting

| Symptom | Cause | Fix |
| --- | --- | --- |
| `--check` exits 1 | a skill changed since the last build | `node gen-skills-index.mjs` |
| `--search` prints "degraded to keyword-only" | embedding server down | `brew services start ollama` (search still works meanwhile) |
| `dim mismatch` on `--search` | corpus built with a different model | `node gen-skills-index.mjs --embed` to rebuild vectors |
| brew `ollama` shows `error` | a manual `ollama serve` is squatting port 11434 | kill it, then `brew services restart ollama` |
| no results for a query | nothing matched above the threshold | lower `--threshold`, or read `SKILLS-INDEX.md` directly |

---

## Overview

The skill library is ~578 skills: 125 top-level skills (each a `~/.claude/skills/<id>/SKILL.md`)
plus 453 "folded" spokes whose verbatim content lives under a hub's `references/<spoke>.md`.
Before this system, routing metadata was fragmented across **per-family `*-manifest.json`** files
that stored only the pre-`TRIGGER` `routingLine` — discarding the structured `TRIGGER`/`SKIP`
clauses, peer edges, version, and usage signal an agent needs to select **one** skill without
reading 117+ files. `gen-skills-index.mjs` is a zero-dependency consolidator that **joins the
existing sources of truth** into a single cross-family index in two forms:

- `SKILLS-INDEX.json` — full machine index (every field, all 578 skills).
- `SKILLS-INDEX.md` — compact, agent-ingestible routing table grouped by family.

It does **not** replace the manifests, `hub-registry.mjs`, or the `tiering/` ranker — it reads them.

---

## 1. Design rationale

### 1.1 Per-family manifests → consolidated hub-and-spoke index

The pre-existing pipeline keyed everything by family: one `*-manifest.json` per family, each
parsed independently. That shape is correct for *building* a family but wrong for *selection*,
where the agent's question is cross-cutting ("which of all 583 skills fits this task?"). Three
concrete gaps drove consolidation:

- **No cross-family view.** Picking between, say, `mongodb-expert` and `da-data-engineering-platform`
  required opening two family manifests. The consolidated index puts every skill in one table.
- **Lossy routing.** Manifests persisted only `routingLine` (the summary *before* `TRIGGER:`).
  The structured trigger/skip routing that 109+ `SKILL.md` files carry was dropped on the floor.
- **Hubs grouped with their spokes.** A spoke inherits its hub's family; a hub belongs to its own
  family; only true standalones are family-less. This keeps each hub adjacent to the spokes it
  routes to, which is exactly the adjacency a selector needs.

The index is **derived, never authoritative** — `SKILL.md` frontmatter and `hub-registry.mjs`
remain the source of truth, so consolidation adds a read-optimized projection without a second
system to keep correct by hand.

### 1.2 Usage-based ranking (`tiering/access-log.jsonl`)

Selection ties are broken by *recent demonstrated use*, not guesswork. `loadRank()` reads
`tiering/access-log.jsonl` (one JSON line per skill access: `{ts, spoke, hub, via}`), counts
accesses within a **30-day window** (`RANK_WINDOW_DAYS`), and emits `rankAccess30d` per skill.

- Skills sort by `rankAccess30d` **descending**, then `id` **ascending** for stable diffs.
- Within each family table in the `.md`, rows are ordered by this rank, so the most-used skill
  in a family surfaces first.
- The signal is **soft**: a `0` rank (never-accessed in window) is common and simply sorts last;
  it never hides a skill. This avoids a cold-start trap where new skills become invisible.

### 1.3 Trigger-based routing extraction (`TRIGGER` / `SKIP`)

`splitDesc()` parses each `description` into three routing-grade parts:

- **`summary`** — text before the first `TRIGGER:`/`SKIP:` marker (the one-line "what is this").
- **`triggers[]`** — clauses after `TRIGGER:`, split on `;` and `·`. These are the positive
  match phrases an agent scans.
- **`skip[]`** — clauses after `SKIP:`, each typically a deferral ("X → other-skill").

Crucially, `SKIP` clauses are mined for **peer edges**: the regex
`/(?:→|->|use)\s+([a-z0-9][a-z0-9-]+[a-z0-9])\b/i`
harvests the target skill id, so "→ mongodb-atlas-expert" becomes a discoverable peer. The final
`peers[]` set is `related_skills ∪ SKIP-routed targets ∪ owning hub` (minus self) — turning prose
deferrals into a navigable graph for peer discovery.

### 1.4 Path portability (home-relative `~/` resolution)

All persisted paths are stored **home-relative** via `tilde(p)` (`p.startsWith(HOME) ? "~"+… : p`)
and re-resolved at read time against `os.homedir()`. This is not cosmetic: the legacy
`tiering/tier-config.json` still bakes in a stale absolute root (`~/.claude/…`)
from a prior machine identity. By never persisting an absolute home path, the index stays valid
across machines and user-name changes, and every workspace under `$HOME` can consume the same
artifact. The generator itself resolves `SKILLS_ROOT` from `os.homedir()` at runtime for the same
reason.

### 1.5 Applicability boundary — when a generated index fits (and when it doesn't)

This index earns its cost only where **two conditions hold together**: (a) there is **no
retrieval engine** in front of the corpus, so an agent would otherwise load metadata to choose;
and (b) the corpus is **many independent, separately-selectable units**. The skill library is the
exact case — 583 discrete skills and a selector with no semantic search, so a single compact
routing table replaces reading 100+ files (§4).

**Code repositories generally fail condition (a) and should not get a generated index:**

- **Augment** auto-indexes the working tree every run with a real semantic context engine
  (embeddings, relationship-aware, history-aware). A generated `REPO-INDEX.md` duplicates what the
  engine already does and becomes a drift liability — the opposite of the skill case.
- **Claude Code** has no semantic engine (it navigates by on-demand `glob`/`grep`/`read`), so a
  curated map *can* help — **but only if it is in the always-loaded layer.** A standalone index
  file is just another artifact the agent must first discover, so the map belongs **inside (or
  `@`-imported from) the auto-loaded `CLAUDE.md`/`AGENTS.md`**, never as a separate generated file.
- For a small, flat repo the map is near-worthless (the agent orients in one or two globs); it pays
  off only for large trees, and even then as an `@`-imported doc, not a standalone index.
- Caveat for cross-tool repos: Augment auto-loads **both** `CLAUDE.md` and `AGENTS.md`, so keeping
  identical copies double-loads the same guidance — prefer one file (symlink the other) unless a
  tool reads only `AGENTS.md`.

Net: the always-loaded guidance file is the **code-repo analog** of this index; the standalone,
generated `SKILLS-INDEX.*` shape is specific to a retrieval-less, many-unit corpus.

---

## 2. Implementation architecture

`~/.claude/skill-consolidation/gen-skills-index.mjs` — 455 lines, **zero runtime dependencies**
(only `node:fs`, `node:fs/promises`, `node:os`, `node:path`, `node:crypto`, `node:url`, and the
local `hub-registry.mjs`). It is an ESM module run directly with `node`.

### 2.1 Pipeline

```text
loadHubRegistry()   →  spoke↔hub↔family relationships (canonical manifest parser)
loadRank()          →  Map<id, rankAccess30d> from tiering/access-log.jsonl (30d window)
resolveSources(reg) →  Map<id, {path, kind}>  (top-level dirs + folded references/<spoke>.md)
readAll(paths, 32)  →  parallel bounded read of every SKILL.md  (← also the cache-warm)
build()             →  per-skill record (frontmatter + splitDesc + peers + rank + srcBytes)
renderMd(idx)       →  compact family-grouped markdown table
main()              →  write both outputs | --check | --stdout
```

### 2.2 Source resolution

`resolveSources()` first walks `~/.claude/skills/*/SKILL.md` for **top-level** skills (symlinks
into `.agents/skills` are followed transparently). For every registry spoke without a top-level
dir, it falls back to the **folded** copy at `<hub>/references/<spoke>.md` that `build.mjs`
maintains. Each entry is tagged `kind: "top-level" | "folded-spoke"`.

### 2.3 Parallel read pool (concurrency 32) = cache-warm

`readAll(paths, concurrency = 32)` launches 32 async workers that drain a shared cursor over the
path list, each `await fsp.readFile(...)`. On a high-latency / scanner-throttled filesystem (the
managed-Mac on-access scanner makes the *first* read of each file slow), this bounded fan-out pays
the cold-read cost **once, in parallel**, in the main process. A failed read stores `null` rather
than throwing, so one unreadable file never aborts the build. This mirrors the cache-warming
strategy in the Sniffies `test/global-setup.mjs`.

### 2.4 Dual-output schema

**`SKILLS-INDEX.json`** (schemaVersion 1) — top-level keys: `generatedAt`, `schemaVersion`,
`skillsRoot`, `counts`, `families`, `skills[]`. Each `skills[]` record:

| field | type | meaning |
| --- | --- | --- |
| `id` | string | skill id (directory or spoke name) |
| `name` | string | frontmatter `name`, falls back to `id` |
| `kind` | enum | `"top-level"` \| `"folded-spoke"` |
| `isHub` | bool | this skill is itself a router/hub |
| `hub` | string\|null | hub this skill is a spoke of |
| `family` | string\|null | spoke inherits hub's family; hub uses own; else null |
| `version` | string\|null | frontmatter `version` |
| `updated` | string\|null | frontmatter `updated` date |
| `origin` | string\|null | frontmatter `origin` |
| `summary` | string | description text before `TRIGGER:`/`SKIP:` |
| `triggers[]` | string[] | positive match clauses |
| `skip[]` | string[] | deferral clauses |
| `peers[]` | string[] | `related_skills ∪ SKIP-targets ∪ hub` (self removed) |
| `rankAccess30d` | int | accesses in trailing 30 days |
| `srcBytes` | int | byte size of the source `SKILL.md` |
| `path` | string | **home-relative** (`~/…`) source path |

`counts` example: `{"total":583,"topLevel":130,"foldedSpokes":453,"families":30}`. `families` is an
id-sorted map of `family → [skill ids]`.

**`SKILLS-INDEX.md`** — an `AUTO-GENERATED` header, a one-line corpus summary, an agent-usage note
("scan to route, then READ the chosen `SKILL.md` for depth"), then one `## <family>` section each
containing a table: `| skill | v | summary | top triggers | peers |`. Hubs are tagged `*(hub)*`,
folded spokes `*(folded)*`. Cells are truncated (summary ≤140 chars, ≤4 triggers, ≤6 peers) and
pipe-escaped to keep the table compact and render-safe.

> **Consumption guidance:** an agent should read the **`.md`** (~21.6K tokens) for routing; the
> **`.json`** (~103K tokens) is for tooling that needs full fields. Reading the `.json` into an
> agent context forfeits most of the token win in §4.

---

## 3. Synchronization & drift control

The index is derived, so any mutation of a `SKILL.md` (description, triggers, version) or of the
tree shape (folds, re-files, new hubs) makes it stale. Freshness is enforced at **three
write-points** plus an **idempotent verifier**. The contract everywhere is the same: regenerate
where skills change, then `--check` to gate.

### 3.1 `--check` — idempotent verification

`node gen-skills-index.mjs --check` rebuilds the index in memory and compares it against the
on-disk `SKILLS-INDEX.json` with `generatedAt` normalized out (so the timestamp alone never
trips it). Identical → exit 0 + "up to date"; different → stderr `SKILLS-INDEX.json is STALE` +
**exit 1**. It writes nothing, so it is safe in CI and pre-commit gates.

### 3.2 `skill-optimizer` — Step 7.7 + Step 8 (v2.12.0)

`skill-optimizer` mutates a skill's description/triggers/version (and Step 7.5 compresses byte
size), so it owns the **post-completion regeneration hook**:

- **Step 7.7 — Refresh the index.** After the Step 7.5 compress, run
  `node ~/.claude/skill-consolidation/gen-skills-index.mjs` to rewrite both outputs, then
  `gen-skills-index.mjs --check` which must exit 0; a non-zero means the regen didn't land →
  re-run once and re-gate. Runs even under `--meta` and `--no-sync` (the index is a local
  artifact, and `--meta` still changes frontmatter). **Non-blocking:** on a missing/erroring
  generator, record `index: skipped (<reason>)` and continue — it never affects Step 7 hub sync.
- **Step 8 — Report.** One line: `index: refreshed (N skills)` or `index: skipped (<reason>)`.

### 3.3 `skill-tree-architect` — Phase 3 step 9 + Phase 4 gate (v1.3.0)

A tree rebalance changes folds, re-files, and hub/family membership wholesale:

- **Phase 3 step 9 — Regenerate.** After the tree mutates, run `gen-skills-index.mjs` to rebuild
  from the new tree state. Non-blocking (`index: skipped (<reason>)` on error).
- **Phase 4 step 1 — Gate.** In VERIFY, `gen-skills-index.mjs --check` must exit 0. A `STALE`
  exit means Phase 3 step 9 didn't run — re-run it without `--check` and re-gate.

### 3.4 Sniffies `test/global-setup.mjs` — non-fatal drift gate

A separate but identical pattern guards the Sniffies userscript's generated `INDEX.md`. After the
existing `node_modules` cache-warm, `checkIndexFreshness()` runs `regen-index.mjs --check` via
`execFile`. It is deliberately **non-fatal**: a stale index emits
`console.warn("[global-setup] INDEX.md drift check — …")` but never fails the suite, so test
success is not coupled to doc freshness during active editing. A one-line comment documents how
to flip the `console.warn` to a `throw` for hard-fail. (This gate already caught one real drift:
the source had advanced to `…-0.8.2.txt` while `INDEX.md` lagged.)

---

## 4. Performance & token metrics

Measured against the live corpus (token figures use the ~4-chars/token heuristic; Claude
tokenization runs slightly denser, so treat as ±15%). The absolute byte/token/count figures
below are a **point-in-time snapshot** — their scan populations (136 physical `SKILL.md`, 118
carrying an always-on description) predate the current 583-skill / 130-top-level headline and
will lag it; the durable result is the **ratios** (§4.1–§4.2), not the raw totals. The win is
**avoiding full-file reads during skill selection and peer discovery**, not shrinking the
always-on description layer.

| Artifact | Bytes | ≈ Tokens | Coverage |
| --- | --- | --- | --- |
| All 136 `SKILL.md` (full) | 2,396,637 | ~599K | 136 skills |
| `description:` frontmatter only (always-on) | 116,710 | ~29.2K | 118 skills |
| **`SKILLS-INDEX.md`** | 86,571 | **~21.6K** | **583 skills** |
| `SKILLS-INDEX.json` | 412,355 | ~103K | 583 skills |
| Avg single `SKILL.md` | 17,622 | ~4.4K | 1 skill |

### 4.1 Library scans (worst case)

Routing by reading the whole corpus (599K tokens) vs. reading the index once (21.6K tokens):
**~96% reduction (~577K tokens saved)** — and the index covers **5× more skills** (583 vs the 118
with an always-on description), because folded spokes have no standalone description at all.

### 4.2 Per-selection events (realistic)

Disambiguating a handful of overlapping skills: without the index, an agent reads ~5–10 candidate
full `SKILL.md` (≈22K–44K tokens); with it, it reads the relevant family slice or the whole 21.6K
table once (≈1K–5K tokens). **Savings ≈ 20K–40K input tokens per non-trivial selection.**

### 4.3 Cost

Input-token cost scales linearly with the per-event delta (actual rates governed by
`~/.llm-cache-a/prices.json`):

| ~30K saved/event | $3/Mtok | $5/Mtok | $15/Mtok |
| --- | --- | --- | --- |
| per selection | $0.09 | $0.15 | $0.45 |
| 5,000 selections/yr (~150M tok) | ~$450 | ~$750 | ~$2,250 |

### 4.4 Latency

- **Cold I/O** (scanner-throttled FS, ~1–2 s/file first read): 1 index read vs 5–10 file reads ≈
  **5–18 s saved** per cold selection. The concurrency-32 pool collapses the generator's own cold
  pass into a single parallel burst.
- **LLM prefill**: ~30K fewer input tokens at ~2–5K tok/s ≈ **6–15 s saved** + lower
  latency-to-first-token.
- **Combined: ~10–30 s per non-trivial selection**; warm-cache runs drive the I/O term to ~0.

> Caveat: vs. the always-on description layer (29K tok) the index (21.6K) is roughly size-neutral —
> its edge there is coverage + structure (TRIGGER/SKIP + `rankAccess30d`), not raw token savings.

---

## 5. Operational guidance

### 5.1 Manual regeneration

```bash
# Regenerate both outputs (re-reads every SKILL.md via the parallel pool):
node ~/.claude/skill-consolidation/gen-skills-index.mjs

# Verify freshness without writing (CI / pre-commit):
node ~/.claude/skill-consolidation/gen-skills-index.mjs --check

# Print the JSON to stdout without touching disk (piping / inspection):
node ~/.claude/skill-consolidation/gen-skills-index.mjs --stdout

# Suppress the informational stderr line:
node ~/.claude/skill-consolidation/gen-skills-index.mjs --quiet
```

Regenerate after **any** of: editing a `SKILL.md` description/triggers/version; adding, folding,
or re-filing a skill; a `skill-tree-architect` rebalance. The optimizer/architect hooks (§3) do
this automatically; the manual command is for ad-hoc edits and recovery.

### 5.2 `--check` exit codes

| Exit | Meaning | Action |
| --- | --- | --- |
| `0` | Index matches the live tree (or `--stdout`/normal write succeeded) | none |
| `1` | **STALE** — on-disk `SKILLS-INDEX.json` differs from a fresh build | run the generator (no flag), then re-`--check` |
| `2` | Unhandled error (e.g. registry unreadable) — thrown by `main().catch` | inspect stderr stack; fix the source/registry, re-run |

The `--check` comparison normalizes out `generatedAt`, so a timestamp difference alone never
reports stale. If `--check` keeps failing after a regen, suspect a non-deterministic input (an
access-log write mid-run) and re-run once.

### 5.3 Recovery & gotchas

- **Missing generator / registry** → hooks record `index: skipped (<reason>)` and continue; the
  index is never on the critical path of a sync.
- **Stale `tier-config.json` absolute paths** are expected and harmless — the generator resolves
  `SKILLS_ROOT` from `os.homedir()` and persists only `~/…` paths (§1.4).
- **Do not hand-edit** `SKILLS-INDEX.{json,md}` — both carry an `AUTO-GENERATED` header and are
  overwritten on every run; edit the upstream `SKILL.md` or the generator instead.

---

## 6. Semantic retrieval architecture

> Optional layer. The regex `TRIGGER`/`SKIP` routing in §2 is the default and only
> always-available router; the vector layer below augments it and is never required for the
> index to build, verify, or be read.

### 6.1 What it adds

`--search "<query>"` answers "which skill best fits this free-text task?" by a **hybrid** score
(cosine similarity over per-skill embeddings, blended with literal keyword overlap — §6.5),
catching matches the literal `TRIGGER` keywords miss (paraphrase, synonym, cross-domain) while
still rewarding exact jargon hits. It **complements** — does not replace — the regex layer:

- **Regex `TRIGGER`/`SKIP`** (§2): deterministic, offline, present in every artifact; the
  routing an agent uses by reading `SKILLS-INDEX.md`. Always available.
- **Vector `--search`**: opt-in, needs a running embedding server and a generated corpus; used
  for recall on fuzzy queries. On any failure it **falls back** to the regex layer.

### 6.2 Model & the same-model invariant

Embeddings come from a **local** server (Ollama, OpenAI-compatible `POST /api/embed`), so the
zero-npm-dependency property holds (built-in `fetch` only) and there is no API key, cost, or
outbound network on the default path.

Pinned model: **`qwen3-embedding:4b`** (Qwen3-Embedding-4B) — Apache-2.0 / open-weight, **2560-d**,
~2.5 GB. Chosen over the smaller `qwen3-embedding:0.6b` (1024-d, ~639 MB) on the §6.9 A/B, where it
won every metric. Override via `SKILLS_EMBED_MODEL` and `OLLAMA_HOST`.

**Invariant:** cosine similarity is only meaningful within one model + dimension. The corpus
records its `model` and `dim`; `--search` re-embeds the query with the same model and **errors
(exit 2)** if the query vector's dimension ≠ the stored `dim`. Changing `SKILLS_EMBED_MODEL`
invalidates every cache entry (the source hash mixes in the model name) and forces a full
re-embed.

### 6.3 Artifact: `SKILLS-EMBEDDINGS.json`

Written to the same directory as the other artifacts via `tilde(p)` pathing; kept **separate**
from `SKILLS-INDEX.json` so the `--check` idempotency gate is unaffected.

```jsonc
{
  "schemaVersion": 1,
  "model": "qwen3-embedding:4b",
  "dim": 2560,
  "generatedAt": "<ISO-8601>",
  "vectors": {
    "<skill-id>": { "hash": "<sha256(model+source)[:16]>", "embedding": [/* 2560 floats */] }
  }
}
```

Source text per skill = `name` + `summary` + `triggers` (newline-joined). Current corpus:
**583 vectors, 2560-d**.

### 6.4 Incremental re-embed

`--embed` is incremental. For each skill it computes `hash = sha256(model + "\n" + source)[:16]`;
a stored vector whose hash matches is reused, otherwise the skill is re-embedded (a changed
`model` re-embeds all). Embedding requests run through the same bounded fan-out pattern as the
read pool, sized by `SKILLS_EMBED_CONCURRENCY` (default **8** — lower than the FS pool's 32
because a single local model server, not disk latency, is the bottleneck). Cold build of 583
skills ≈ 21 s; a no-change re-run reuses all 583 and embeds 0 — and is **no-churn**: when the
model and id set are unchanged it skips the write entirely, so the automation in §6.7 has no side
effects on an unchanged library. The skills root and output dir honor `SKILLS_ROOT_DIR` /
`SKILLS_OUT_DIR` overrides, which the tests (§6.8) and benchmark (§6.9) use to run the whole
pipeline against a throwaway corpus.

### 6.5 Search contract

```bash
node gen-skills-index.mjs --embed                               # (re)build the corpus
node gen-skills-index.mjs --search "atlas slow query high CPU"  # ranked table (hybrid by default)
node gen-skills-index.mjs --search "<q>" --alpha=0.92 --top=5 --threshold=0.4 --json
node gen-skills-index.mjs --search "<q>" --no-hybrid            # pure cosine (no keyword blend)
```

- **Hybrid scoring (default):** `score = alpha·cosine + (1 − alpha)·coverage`, where `coverage`
  is the fraction of content query tokens (lowercase alnum runs ≥3, minus a small stop list)
  present in the skill's `name` + `summary` + `triggers`. `alpha` defaults to **0.92**
  (`--alpha=A`); `--no-hybrid` reverts to pure cosine. The keyword bag is read from the offline
  `SKILLS-INDEX.json`, so hybrid adds **no extra embedding call** — if that index is absent,
  coverage is 0 and the score degrades gracefully to cosine.
- Ranks by score **descending**, then `rankAccess30d` **descending**, then `id` ascending —
  the same tie-break order as `build()`'s sort.
- Returns the top **N** (default 10, `--top=N`) scoring **≥ threshold** (default **0.30**,
  `--threshold=T`). Default output is a `score · id · summary` table; `--json` for machine use
  (in hybrid mode each hit also reports its `cosine` and `coverage` components).
- **Exit codes:** `0` = at least one hit; `1` = nothing cleared the threshold, or the corpus is
  missing/empty (with a message to fall back to regex routing); `2` = error (empty query, dim
  mismatch, or query-embed failure).

### 6.6 Secrets & fail-open

- A local Ollama needs no key. If `OLLAMA_HOST` points at an authenticated gateway, the key is
  read from `OLLAMA_API_KEY` / `OPENAI_API_KEY`, sent as a bearer header, and **never logged,
  printed, or written into any artifact**.
- Fail-open wherever the network is involved: `--embed` preflights the server and exits 2
  **without touching** the existing corpus if it is unreachable; a per-skill embed failure keeps
  the prior vector. `--search` degrades to the regex `TRIGGER`/`SKIP` routing on a missing/empty
  corpus or any embed failure. The default run and `--check` never call the network and stay
  byte-deterministic.
- Corpus writes are **atomic** (write `.tmp`, then `rename`), so a manual `--embed` racing the
  WatchPaths agent (§6.7) can never leave a half-written `SKILLS-EMBEDDINGS.json`.

### 6.7 Automation — `embed-refresh.sh` + the `com.mitch.skills-embed` launchd agent

The corpus stays in sync with the library automatically. A user LaunchAgent
(`~/Library/LaunchAgents/com.mitch.skills-embed.plist`) runs `embed-refresh.sh`, which
regenerates the index and then runs the incremental `--embed`:

- **On skill add / remove / rename** — `WatchPaths` fires on `~/.claude/skills`. (launchd watches
  a directory's direct children, so a new top-level skill dir triggers it immediately.)
- **On edits inside an existing skill** — a 6-hour `StartInterval` catches content changes that
  don't alter the watched directory's own mtime; the source-hash check then re-embeds only what
  actually moved.
- **At load / login** — `RunAtLoad` catches anything changed while logged out.

It is fail-open and cheap: a down server leaves the index fresh and the agent still exits 0, and
the no-churn `--embed` (§6.4) skips the write when nothing changed, so the periodic run has no
side effects. Output is appended to `embed-refresh.log`. The embedding server is the
Homebrew-managed `ollama` service (`brew services start ollama`), so both halves survive
reboot/login.

```bash
launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/com.mitch.skills-embed.plist  # install
launchctl bootout   gui/$(id -u)/com.mitch.skills-embed                               # uninstall
```

### 6.8 Tests — `gen-skills-index.test.mjs`

```bash
node --test gen-skills-index.test.mjs
```

`main()` auto-runs only as the process entrypoint, so the module is importable and the suite
exercises the real routines. **Offline unit tests** (always run) cover `cosine`, `embedSource`,
`srcHash`, `splitDesc`, and a populated `build()`. **Live integration tests self-skip when the
embedding server is down**: `embedText` dimension, the `--search` ranking + `0/1/2` exit
contract, and an end-to-end incremental proof (new skill → embedded, then reuse, then no-churn)
that runs against a throwaway `SKILLS_ROOT_DIR` + `SKILLS_OUT_DIR` so it never touches the live
corpus or races the agent. Offline unit tests also cover the keyword helpers (`tokens`,
`kwCoverage`); a live test asserts `--search` is hybrid by default and `--no-hybrid` is pure
cosine. Current: **13/13 pass**.

### 6.9 Benchmark — `bench-search.mjs` (five-method A/B)

```bash
node bench-search.mjs            # comparison table + per-query first-hit ranks
node bench-search.mjs --alpha=0.92
node bench-search.mjs --json
```

Over 30 labeled, deliberately **paraphrased** queries (synonyms / described symptoms that often
share no keyword with the target's `TRIGGER` list), it A/Bs five routing methods by P@1, P@3, and
MRR: a **keyword** baseline (the literal token overlap the regex router depends on), **semantic**
cosine over each corpus (`4b`, `0.6b`), and **hybrid** (cosine + keyword coverage) over each. A
query is a hit if any acceptable target (the hub **or** its precise spoke) lands in top-k. The
primary corpus is the live production `SKILLS-EMBEDDINGS.json` (now 4b; override via
`SKILLS_4B_CORPUS`) and is **required**; the smaller-model comparison corpus is read from
`/tmp/skills-0.6b/SKILLS-EMBEDDINGS.json` (override via `SKILLS_06_CORPUS`) and, when absent, its
two rows are skipped. Read-only.

| Method | P@1 | P@3 | MRR | mean query latency |
| --- | --- | --- | --- | --- |
| **hybrid-4b** (`qwen3-embedding:4b`, 2560-d, α=0.92) | **70%** | **87%** | **0.775** | ~220 ms (+kw) |
| semantic-4b | 67% | 80% | 0.744 | ~220 ms |
| hybrid-0.6b (α=0.92) | 60% | 77% | 0.699 | ~120 ms (+kw) |
| semantic-0.6b | 53% | 77% | 0.644 | ~120 ms |
| keyword baseline | 47% | 50% | 0.525 | ~0 ms |

Two findings drove the production defaults (§6.2, §6.5):

- **The 4b model wins outright.** Bigger embeddings rescue the jargon queries 0.6b mis-routed —
  e.g. "the borrow checker keeps rejecting my code" → `lang-rust` jumps from rank ~35
  (cosine-0.6b) to rank 1 (cosine-4b).
- **Hybrid beats pure cosine at α=0.92.** Blending in keyword coverage lifts several borderline
  queries into top-3 (e.g. "build an MCP server" 6→2, "tighten this system prompt" 2→1) **without**
  hurting paraphrase recall — hybrid-4b dominates semantic-4b on all three metrics. An α sweep
  (0.6→0.95) put the optimum at **0.92**; below it the keyword term over-weights and P@1 drops.

A residual hard case — "prepare a quarterly business review for a key account" still ranks
`executive-comms` over `tam-operations`, even under hybrid-4b (the keyword boost narrows but
doesn't close the gap) — is why the layer stays a **complement** to regex routing (§6.1), not a
replacement.

---

## 7. Cluster → agent overlay analysis

Once every skill carries a `qwen3-embedding:4b` vector (§6), the 17 agent definitions in
`~/.claude/agents/` can be projected into the same 2560-d space and measured against the 30 family
centroids. This answers a structural question raised during design: are the naturally-emerging skill
**families** "nascent agents," or do agents follow a different organizing force? The pass is
**read-only** and runs entirely against the cached corpus — no Ollama call, no network.

### 7.1 Method

- **Family centroid.** L2-normalize all 583 skill vectors; for each of the 30 families, mean its
  members and renormalize → 30 unit centroids.
- **Agent position.** Harvest the skill IDs each agent `.md` cites, average those skills' normalized
  vectors → an agent centroid in the same space.
- **Alignment.** Cosine of the agent centroid to each family centroid, plus **membership share**
  (fraction of cited skills falling in the top family) and **normalized entropy** `Hn` of the
  family distribution (0 = one family, 1 = maximally spread).
- **Same-model invariant (§6.2) holds.** Agents are placed with the identical `qwen3-embedding:4b` /
  2560-d vectors the skills use; there is no cross-model comparison.

### 7.2 Two corrections (both moved the result)

1. **Grab-bag confound.** `(standalone)` (65 skills) and `misc-catch-all` sit near the global mean,
   so every diffuse agent ranked "nearest" them and spuriously read as cross-cutting. Re-ranked
   against **real topical families only**.
2. **Harvest contamination.** Single-word skill IDs (`check`, `build`, `review`, `research`) matched
   as prose verbs in the agent text. Restricted bare tokens to **backtick-cited** mentions;
   hyphenated IDs (which can't be plain English) match anywhere.

After both corrections, mean real-family cosine for the 12 placed agents = **0.89**.

### 7.3 Result — a four-way structure, not a binary

| Class | n | Agents (nearest real family · cos · margin) | Reading |
| --- | --- | --- | --- |
| **Nascent agent** (≈ one family) | 3 | `mongodb-claim-validator` (mongodb · 0.96 · 0.22), `uber-mongodb-diagnostician` (mongodb · 0.94 · 0.21), `customer-comms-psychologist` (writing · 0.89 · 0.07) | Cluster closed under a control loop; sits on a hub centroid |
| **Single-skill core** (thin) | 4 | `incident-postmortem-drafter` + `tam-weekly-update-builder` (both = `document-critique`), `repo-health-auditor` (= `repo-bootstrapper`), `error-monitor-remediator` (= `debugging`) | One skill + workflow shell |
| **Cross-cutting** | 4 | `tam-assistant` (7 families), `multi-lens-doc-reviewer` (9), `skill-gap-filler` (3), `convergence-loop-runner` (2) | Merged by workflow co-occurrence, not topic |
| **Orchestration** | 5 | `account-data-collector`, `account-state-delta-watcher`, `firedrill-scenario-runner`, `harsh-reviewer`, `tam-doc-validator` | Cite **zero** skills; pure control/permission |
| *Mixed* | 1 | `prompt-optimizer-loop` (optimizers · 0.95 · margin 0.03 from `ai-agent`) | Near-fused neighbor families |

Headline: `meanRealCos 0.888` · nascent 3 · cross-cutting 4 · thin 4 · orchestration 5 · mixed 1 ·
12/17 placed.

### 7.4 Interpretation

- **The "clusters are nascent agents" thesis is confirmed where knowledge dominates.** The two
  MongoDB agents sit almost *on* the `mongodb` hub centroid (cos 0.94–0.96) with a large separation
  from the next family (margin > 0.2) — the cleanest possible "cluster closed under a control loop,"
  and it confirms the **hub-as-router** reading: the nascent agent is the hub node, not the family
  blob.
- **The strongest micro-evidence is the `thin` class.** Two separate agents are each just
  `document-critique` wrapped in a different workflow shell. The nascent unit isn't always a
  *family* — it can be a **single skill** promoted by adding a trigger + loop. That is the
  cluster→agent transition at its smallest granularity, visible in the roster.
- **The counter-evidence is real and localized.** The 5 `orchestration` agents cite zero skills —
  topically invisible because they exist for **side-effect / coordination** reasons (collect data,
  watch deltas, run a drill, review, validate). The 4 `cross-cutting` agents are the other force:
  **workflow co-occurrence** fusing 7–9 families into one job. Neither is a nascent cluster; both
  name a force *other* than topical similarity.
- **Adjacency is itself a signal.** `prompt-optimizer-loop` is topically a clean optimizers agent
  (100% membership, cos 0.95) but lands `mixed` only because `optimizers` and `ai-agent` are 0.03
  apart in the manifold — near-fused neighbors, which is correct: prompt/skill optimization *is*
  applied agent-engineering.

**Bottom line.** Skill clusters are nascent agents only in the read-only / single-domain corner of
the roster (MongoDB diagnosis, doc-critique, comms). The moment an agent acquires write side-effects
or a multi-source workflow, its boundary detaches from the knowledge manifold — and the overlay
localizes precisely where (`orchestration` = permission force, `cross-cutting` = workflow force).

### 7.5 Caveats & reproduction

- **N = 17**, and placement is **citation-based**: an agent that *uses* a skill without naming it
  reads thinner than it is. This most affects the `orchestration` agents, which call MCP tools and
  other agents by design — hence 5 of 17 are unplaceable by citation.
- A complementary pass — embedding each agent's **full description prose** through
  `qwen3-embedding:4b` — places all 17 agents by *meaning* rather than cited skills, rescuing the 5
  orchestration agents. It has now been **run** (§7.6): the Ollama runtime was reinstalled and the
  model re-pulled, and the same-model invariant (§6.2) is preserved (identical model, 2560-d).
- Reproduction: the citation pass is **deterministic** against the cached corpus via the scratch
  scripts `overlay.mjs` (harvest + overlay) and `render.mjs` (table); the prose pass uses
  `overlay-prose.mjs`, which embeds each agent's frontmatter `description` against the live
  `qwen3-embedding:4b` server at query time (deterministic for a fixed model build).

### 7.6 Complementary pass — placement by full-description prose

Instead of harvesting cited skill IDs, this pass embeds each agent's frontmatter `description`
(`name` + `description` — the direct analog of the skill `embedSource` text in §6.3) with the same
`qwen3-embedding:4b` / 2560-d vectors and cosines it to the same 30 family centroids. It places
**all 17** agents — including the 5 orchestration agents that cite zero skills — by what they *say
they do* rather than what they name.

| Agent | Citation (§7.3) family · cos | Prose family · cos | Verdict |
| --- | --- | --- | --- |
| `mongodb-claim-validator` | mongodb · 0.96 | mongodb · 0.92 | **stable** |
| `uber-mongodb-diagnostician` | mongodb · 0.94 | mongodb · 0.92 | **stable** |
| `prompt-optimizer-loop` | optimizers · 0.95 | optimizers · 0.91 | **stable** |
| `tam-assistant` | tam-operations · 0.94 | tam-operations · 0.88 | **stable** |
| `convergence-loop-runner` | optimizers · 0.95 | optimizers · 0.88 (mgn 0.00) | **stable** (seam) |
| `harsh-reviewer` | *(0 skills)* | writing · 0.87 | **rescued** |
| `tam-weekly-update-builder` | content-ingestion · 0.82 | tam-operations · 0.86 | **flip** |
| `skill-gap-filler` | ai-agent · 0.94 | claude-code · 0.86 (mgn 0.00) | **flip** (seam) |
| `customer-comms-psychologist` | writing · 0.89 | writing · 0.86 | **stable** |
| `tam-doc-validator` | *(0 skills)* | tam-operations · 0.85 | **rescued** |
| `multi-lens-doc-reviewer` | software · 0.91 | optimizers · 0.85 | **flip** |
| `account-data-collector` | *(0 skills)* | tam-operations · 0.84 | **rescued** |
| `repo-health-auditor` | optimizers · 0.83 | optimizers · 0.84 | **stable** (seam) |
| `incident-postmortem-drafter` | content-ingestion · 0.82 | tam-operations · 0.83 | **flip** |
| `account-state-delta-watcher` | *(0 skills)* | tam-operations · 0.82 | **rescued** |
| `error-monitor-remediator` | software · 0.71 | mongodb · 0.81 (mgn 0.00) | **flip** (seam) |
| `firedrill-scenario-runner` | *(0 skills)* | ai-agent · 0.73 | **rescued** |

Mean real-family cosine: prose **0.853** across all 17 vs citation **0.888** across the 12 it could
place — prose is slightly lower and *should* be: it is a noisier but **complete** signal that now
includes the diffuse coordinators. Tally: **7 stable · 5 rescued · 5 flip**.

What the prose pass reveals:

- **The orchestration agents have a topical home after all.** All 5 resolve, and **3 of 5 land in
  `tam-operations`** (`account-data-collector` 0.84, `tam-doc-validator` 0.85,
  `account-state-delta-watcher` 0.82); `harsh-reviewer` → `writing` 0.87, `firedrill-scenario-runner`
  → `ai-agent` 0.73. They were topically *invisible* to citation (zero named skills), but by meaning
  they are TAM-domain workflow coordinators — their prose is saturated with the TAM operating
  vocabulary even though they name no skill.
- **Citation places by *mechanism*; prose places by *purpose* — and the flips prove it.** The
  cleanest pair: `tam-weekly-update-builder` and `incident-postmortem-drafter` both flip
  `content-ingestion-extraction → tam-operations`. Citation saw their one tool (`document-critique`,
  bucketed under content-ingestion); prose saw the *job* (a TAM deliverable). Likewise
  `multi-lens-doc-reviewer` flips `software → optimizers`: by its toolbox it is a software-knowledge
  agent, by its purpose a review/optimization agent.
- **Prose is more honest about seams.** Agents that straddle near-fused families collapse to
  margin ≈ 0 under prose (`convergence-loop-runner` 0.00, `skill-gap-filler` 0.00,
  `error-monitor-remediator` 0.00, `repo-health-auditor` 0.01, `firedrill-scenario-runner` 0.01).
  Sparse citation can manufacture a clean win from a single lucky skill; prose, drawing on the whole
  description, refuses to force the call where the manifold does not separate.

**The 7 stable placements are the method-robust core.** Both passes agree on the two MongoDB agents
(`mongodb` wins under *both* methods with cos ≥ 0.92), the two optimizer-loop agents (`optimizers`),
`tam-assistant` (`tam-operations`), and `customer-comms-psychologist` (`writing`). Where the two
methods agree the nascent-agent reading is strongest; where they disagree, the disagreement *names
the force* — tool vs job, or an unseparated seam.

**Refined bottom line.** §7.4 said the orchestration agents detach from the knowledge manifold. The
prose pass sharpens that: they don't detach from the *manifold*, they detach from the *citation
graph*. By meaning they sit squarely in `tam-operations`. The deepest reading is that the TAM
operating domain is itself a (coarse) family, and the coordinators are its nascent agents — bound by
workflow and permission rather than by any single knowledge hub.

## 8. Appendix: file map & sources

| Path | Role |
| --- | --- |
| `~/.claude/skill-consolidation/gen-skills-index.mjs` | the generator (this doc's subject) |
| `~/.claude/skill-consolidation/SKILLS-INDEX.json` | machine index output |
| `~/.claude/skill-consolidation/SKILLS-INDEX.md` | agent routing-table output |
| `~/.claude/skill-consolidation/SKILLS-EMBEDDINGS.json` | semantic corpus (§6); separate from `--check` |
| `~/.claude/skill-consolidation/gen-skills-index.test.mjs` | test suite (§6.8) — `node --test` |
| `~/.claude/skill-consolidation/bench-search.mjs` | semantic-vs-keyword benchmark (§6.9) |
| `~/.claude/skill-consolidation/embed-refresh.sh` | launchd refresh wrapper (§6.7) |
| `~/Library/LaunchAgents/com.mitch.skills-embed.plist` | WatchPaths/interval embed agent (§6.7) |
| `~/.claude/skill-consolidation/hub-registry.mjs` | canonical spoke↔hub↔family parser (consumed) |
| `~/.claude/skill-consolidation/build.mjs` | maintains folded `references/<spoke>.md` copies |
| `~/.claude/skill-consolidation/*-manifest.json` | per-family manifests (family lookup source) |
| `~/.claude/skill-consolidation/tiering/access-log.jsonl` | usage signal for `rankAccess30d` |
| `~/.claude/skill-consolidation/tiering/tier-config.json` | legacy config (stale abs paths; see §1.4) |
| `~/.claude/skills/<id>/SKILL.md` | top-level skill source |
| `~/.claude/skills/<hub>/references/<spoke>.md` | folded-spoke source |
| `~/.claude/skills/skill-optimizer/SKILL.md` | hosts Step 7.7 + Step 8 hook (§3.2) |
| `~/.claude/skills/skill-tree-architect/SKILL.md` | hosts Phase 3 step 9 + Phase 4 gate (§3.3) |
| `~/Documents/Sniffies Userscript/test/global-setup.mjs` | sibling non-fatal drift gate (§3.4) |
| `~/Documents/Sniffies Userscript/regen-index.mjs` | the Sniffies `INDEX.md` generator (pattern origin) |

**Related repos:** the live skill registry is synced to `mdb-context-hub`
(`~/Documents/GitHub/mdb-context-hub/scripts/skill-pack.config.mjs`, `npm run sync:skills`); that
sync is review-required and committed by the user, separate from index regeneration.

---

## 9. System Bootstrap Prompt

The block below is a **self-contained prompt**. Handed to a capable coding LLM with no other
context, it reconstructs the indexing system — `gen-skills-index.mjs` and its three artifacts — to
match this architecture, including the semantic retrieval layer (§6) as a core part of the same
utility. It is written to be paste-ready; everything it needs is inline.

> **PROMPT — Build the skill-library indexing utility**
>
> Build a single Node.js ESM script, `gen-skills-index.mjs`, with **no third-party dependencies**,
> that consolidates a hub-and-spoke skill library into one machine index (`SKILLS-INDEX.json`), one
> agent-ingestible routing table (`SKILLS-INDEX.md`), and one semantic-search corpus
> (`SKILLS-EMBEDDINGS.json`). The index is **derived, never authoritative** — it only reads existing
> sources of truth and joins them. Obey every requirement below exactly.
>
> **R0 — Runtime & dependencies.** ESM (`.mjs`), `#!/usr/bin/env node` shebang. Import **only**
> Node built-ins: `node:fs`, `node:fs/promises`, `node:os`, `node:path`, `node:url`, and
> `node:crypto` (for the SHA-256 source hashing in R11). The one local import is a sibling
> `hub-registry.mjs` exposing `loadHubRegistry()` → `{ spokes:Set, hubs:Set, spokeHub:Map<spoke,hub>
> }`. **No third-party packages and no build step** — this zero-dependency invariant is absolute. The
> semantic layer (R10–R12) reaches a local embedding server over HTTP using the **built-in `fetch`**
> only (no HTTP-client library); the default index build, `--check`, and `--stdout` paths make **no
> network call**.
>
> **R1 — Constants & paths (all home-relative for portability).** Resolve `HOME = os.homedir()` at
> runtime; never hardcode an absolute home. Set `SKILLS_ROOT = ~/.claude/skills`. Write all three
> outputs **next to the script** (`path.dirname(fileURLToPath(import.meta.url))`). Read the usage log
> from `<scriptDir>/tiering/access-log.jsonl`. Define `RANK_WINDOW_DAYS = 30`. For the semantic
> layer, read the embedding model from `SKILLS_EMBED_MODEL` (default **`qwen3-embedding:4b`**), the
> server base URL from `OLLAMA_HOST` (default `http://127.0.0.1:11434`), and the embed fan-out from
> `SKILLS_EMBED_CONCURRENCY` (default **8**); honor `SKILLS_ROOT_DIR`/`SKILLS_OUT_DIR` overrides so
> the pipeline can run against a throwaway corpus. Provide a `tilde(p)` helper that rewrites a path
> under `$HOME` to `~/…`, and persist **only** tilde paths so the artifacts are valid across machines
> and usernames. Parse CLI flags from `process.argv.slice(2)`: support `--quiet` (suppress the
> informational stderr line; route all logs to `console.error`, never stdout), `--check`, `--stdout`,
> `--embed`, `--search` (whose value is the query string), `--json`, `--no-hybrid`, and the
> `key=value` options `--alpha=`, `--top=`, `--threshold=`.
>
> **R2 — Frontmatter parsing.** Extract the leading YAML frontmatter (`/^---\n([\s\S]*?)\n---/`).
> Implement `field(fm,name)` that reads a scalar `name: value`, tolerating block scalars (`|`/`>`)
> and folded multi-line values (continuation lines indented under the key), stripping wrapping
> quotes. Implement `listField(fm,name)` for YAML block lists (`name:` then `  - item` lines). Use
> these to read `name`, `description`, `version`, `updated`, `origin`, and `related_skills[]`.
>
> **R3 — Description splitter + TRIGGER/SKIP → peer edges.** `splitDesc(desc)`: collapse
> whitespace; locate `TRIGGER:` and `SKIP:` markers. `summary` = text before the first marker
> (trim trailing punctuation). `triggers[]` = the TRIGGER block split on `;` **and** `·`, trimmed,
> empties dropped. `skip[]` = the SKIP block split the same way. From each skip clause, harvest a
> peer skill id with the regex `/(?:→|->|use)\s+([a-z0-9][a-z0-9-]+[a-z0-9])\b/i` (lowercased) into
> `skipTo[]`. A skill's final `peers[]` = `related_skills ∪ skipTo ∪ {owning hub}`, with self
> removed and de-duplicated. This turns prose deferrals into a navigable graph.
>
> **R4 — Usage rank (30-day window).** `loadRank()`: read `access-log.jsonl` if present (missing →
> empty Map, never throw). Each line is JSON `{ts, spoke, hub, via}`. Count, per `spoke`, the
> entries whose `Date.parse(ts)` is within the trailing `RANK_WINDOW_DAYS`. Return
> `Map<id, count>`; absent → `0`. This rank is a **soft tie-breaker only** — a `0` sorts last but is
> never hidden (no cold-start trap).
>
> **R5 — Source resolution.** `resolveSources(reg)` → `Map<id,{path,kind}>`. First, every
> `<SKILLS_ROOT>/<dir>/SKILL.md` (skip dotfiles; symlinks followed transparently) is a
> `kind:"top-level"` skill. Then, for every registry spoke **without** a top-level dir, fall back to
> the folded copy `<SKILLS_ROOT>/<hub>/references/<spoke>.md` tagged `kind:"folded-spoke"`.
>
> **R6 — Parallel read pool = cache-warm (concurrency 32).** `readAll(paths, concurrency=32)`:
> launch `concurrency` async workers draining a shared cursor, each `await fsp.readFile(p,"utf8")`.
> A failed read stores `null` (never throws) so one bad file can't abort the build. This bounded
> fan-out doubles as the cold-FS cache-warm — pay the first-read cost once, in parallel.
>
> **R7 — Build records.** For each resolved skill with non-null content: compute `family` —
> a spoke inherits its hub's family, a hub uses its own family, a true standalone is `null` (derive
> hub→family by scanning sibling `*-manifest.json` files once and caching; each manifest exposes
> `family` and `hubs`). Emit a record with fields: `id`, `name` (fallback `id`), `kind`, `isHub`
> (`reg.hubs.has(id)`), `hub` (`reg.spokeHub.get(id)||null`), `family`, `version`, `updated`,
> `origin`, `summary`, `triggers[]`, `skip[]`, `peers[]`, `rankAccess30d`, `srcBytes`
> (`Buffer.byteLength(md)`), `path` (tilde). Sort skills by `rankAccess30d` **descending**, then
> `id.localeCompare` **ascending** (stable diffs). Group ids by family into `families` (standalones
> under `"(standalone)"`), and emit an id-sorted `families` map.
>
> **R8 — JSON schema (`SKILLS-INDEX.json`).** Object with `schemaVersion: 1`, `skillsRoot` (tilde),
> `counts` `{total, topLevel, foldedSpokes, families}`, `families` (sorted map family→ids[]), and
> `skills[]` (the records above). On a real write (not `--check`/`--stdout`), prepend a
> `generatedAt` ISO timestamp as the first key. Pretty-print with 2-space indent + trailing newline.
>
> **R9 — Markdown table (`SKILLS-INDEX.md`).** First line:
> `<!-- AUTO-GENERATED by gen-skills-index.mjs — do not edit by hand. -->`. Then an `# Skill Library
> Index` H1, a one-line corpus summary (`N skills (T top-level + F folded spokes) across K families`
> + skills root), and an **agent-usage note**: scan the table to pick the single best skill by its
> triggers, then READ that skill's `SKILL.md` for depth — this table is for routing, not depth; rows
> are ordered by recent-use rank within each family. For each family (sorted): an `## <family>`
> heading and a table `| skill | v | summary | top triggers | peers |`. Per row: tag hubs `*(hub)*`
> and folded spokes `*(folded)*`; truncate summary to 140 chars, triggers to the first 4 (joined
> `"; "`), peers to the first 6 (each wrapped in backticks); pipe-escape (`\|`) and newline-flatten
> every cell; use `—` for empties.
>
> **R10 — Embeddings artifact (`SKILLS-EMBEDDINGS.json`).** The semantic corpus is written **next to
> the script** but kept **separate from `SKILLS-INDEX.json`** so the R13 `--check` idempotency gate is
> never affected by embedding churn. Schema: `{ schemaVersion: 1, model, dim, generatedAt, vectors }`,
> where `vectors` maps `<skill-id> → { hash, embedding:[…floats] }`. The per-skill source text is
> `name` + `summary` + `triggers` (newline-joined) — the same routing-grade text the index already
> holds. Write **atomically** (write a `.tmp` sibling, then `rename`) so a manual `--embed` racing a
> background refresh can never leave a half-written file. **No-churn write:** when the model and the
> id set are unchanged, skip the write entirely.
>
> **R11 — `--embed` (incremental embedding).** Build/refresh the corpus from a local, OpenAI-
> compatible embedding server (Ollama: `POST {OLLAMA_HOST}/api/embed`) using the **built-in `fetch`**
> only. For each skill compute `hash = sha256(model + "\n" + source).slice(0,16)` (via `node:crypto`);
> a stored vector whose `hash` matches is **reused**, otherwise the skill is **re-embedded**. A
> changed `model` invalidates every entry (the model name is mixed into the hash) and forces a full
> re-embed. Run embed requests through the same bounded fan-out as R6, sized by
> `SKILLS_EMBED_CONCURRENCY` (default 8 — a single model server, not disk, is the bottleneck).
> **Fail-open:** preflight the server and exit `2` **without touching** the existing corpus if it is
> unreachable; a per-skill embed failure keeps the prior vector. If `OLLAMA_HOST` is an authenticated
> gateway, read the key from `OLLAMA_API_KEY`/`OPENAI_API_KEY`, send it as a bearer header, and
> **never log, print, or persist it**.
>
> **R12 — `--search "<query>"` (hybrid ranked search).** Re-embed the query with the **same model**
> and **error (exit 2)** if its dimension ≠ the corpus `dim` (cosine is only meaningful within one
> model+dim). Score each skill `score = alpha·cosine + (1 − alpha)·coverage`, where `cosine` is the
> cosine similarity of the query and skill vectors and `coverage` is the fraction of content query
> tokens (lowercase alnum runs ≥3, minus a small stop list) present in the skill's `name` + `summary`
> + `triggers`. `alpha` defaults to **0.92** (`--alpha=`); `--no-hybrid` forces pure cosine
> (`alpha = 1`). Read the keyword bag from the offline `SKILLS-INDEX.json` so hybrid adds **no extra
> embedding call**. **Graceful fallbacks:** if `SKILLS-INDEX.json` is missing, `coverage = 0` and the
> score degrades to cosine-only; if `SKILLS-EMBEDDINGS.json` is missing/empty or any embed fails, fall
> back to the regex `TRIGGER`/`SKIP` routing. Rank by `score` desc, then `rankAccess30d` desc, then
> `id` asc; return the top **N** (default 10, `--top=`) scoring **≥ threshold** (default **0.30**,
> `--threshold=`). Default output is a `score · id · summary` table; `--json` emits machine output
> (each hit also reporting its `cosine` and `coverage` components). Exit `0` = ≥1 hit; `1` = nothing
> cleared the threshold or the corpus is missing/empty (advise the regex fallback); `2` = empty query,
> dim mismatch, or query-embed failure.
>
> **R13 — `main()` modes & exit codes.** Branch on flags (resolve `--search`/`--embed` **before**
> building the full index — they need only the corpus or the keyword bag):
> - `--search "<q>"`: run R12 and exit with its `0/1/2` code; writes nothing to disk.
> - `--embed`: run R11 (build the source text in memory, refresh the corpus); exit 0 on success or
>   `2` if the server is unreachable.
> - `--stdout`: write the JSON (no `generatedAt`) to **stdout**, write nothing to disk, exit 0.
> - `--check`: read the on-disk `SKILLS-INDEX.json`; compare it to a fresh build with `generatedAt`
>   **normalized out** of both sides (timestamp alone must never trip it). Identical → log "up to
>   date", exit 0. Different → `console.error("SKILLS-INDEX.json is STALE — run: node
>   gen-skills-index.mjs")` and `process.exit(1)`. Writes nothing (safe in CI/pre-commit) and **never
>   touches `SKILLS-EMBEDDINGS.json`**.
> - default: write both stamped index outputs, log a one-line summary (suppressed by `--quiet`),
>   exit 0. The corpus is refreshed only by `--embed`, never by a default run.
>
>   Wrap the entry point in `main().catch(e => { console.error(e); process.exit(2); })` so any
>   unhandled error (e.g. an unreadable registry) is exit code **2**, distinct from a clean STALE
>   (exit 1) and success (exit 0).
>
> **R14 — Invariants to preserve.** No third-party packages and no build step; the only network is
> the local embedding server reached via the built-in `fetch`, and only on the `--embed`/`--search`
> paths (the default build, `--check`, and `--stdout` stay offline and byte-deterministic). Never
> persist an absolute `$HOME` path; never throw on a missing/unreadable optional input (access-log, a
> single SKILL.md, an absent corpus). The index is reproducible (same tree ⇒ byte-identical output
> modulo `generatedAt`), and `SKILLS-EMBEDDINGS.json` stays **separate** from it so `--check`
> idempotency is preserved. Nothing is hand-editable (the `AUTO-GENERATED` header declares each output
> overwritten on every run). Do **not** mutate or replace the manifests, `hub-registry.mjs`, or the
> tiering ranker — read them only.
