# harness-streamliner

**Category:** AI, Agents & Prompt Engineering
**Platform:** Claude Code
**Original Path:** claude-code/harness-streamliner

## Description
Streamline and optimize the Codex harness config — settings.json, hooks, MCP servers, permissions, env, model, statusline — to cut redundancy, startup latency, and per-turn token overhead. Runs a read-only scan, then AUTO-APPLIES every Medium-or-higher config fix in place (timestamped backup first, jq-validated after, auto-restore on breakage) and reports what changed. Use this whenever the user wants to "streamline / optimize / clean up / trim / tune / speed up my Codex setup or config", mentions a bloated or slow settings.json, duplicate or redundant permission rules, dead or broken hooks, too many MCP servers eating context, high startup cost, or asks to audit ~/.Codex (and project .Codex/) for cruft. Covers user-level ~/.Codex AND project-level .Codex/settings*.json when a repo is in scope. Triggers on: "/hso", "streamline my Codex setup", "optimize my settings", "my config is bloated", "clean up my permissions", "audit my hooks", "reduce startup overhead", "too many mcp servers". SKIP: skill body/triggering quality -> skill-optimizer; skill tree, count, placement, or tiering -> skill-tree-architect, /skill-tier, /sync-skills; prompt or agent-instruction quality -> prompt-deep-optimizer; a single precise setting change the user already named -> update-config.

---

# harness-streamliner

Optimize the **Codex harness config** for redundancy, startup latency, and
token overhead. This skill owns the *config/perf* layer of the setup — not skill
content, not the skill tree, not prompts. Hand those to their owners (see
**Boundaries**).

Slash alias: **`/hso`**.

## Why this exists

Harness config accretes. Permission lists grow until broad rules silently
subsume narrow ones. Hooks pile up on every event, each adding latency or
injecting context that eats the token budget every turn. MCP servers get added
and never removed, each loading tools into context. Two `defaultMode` keys end
up disagreeing. None of this is visible day-to-day, but it taxes every session's
startup and every turn's budget. This skill makes the cost visible and trims it
safely.

## Operating mode (decided up front)

**Auto-apply Medium+ with backup.** After scanning, apply every Medium-or-higher
fix in place — but only after a timestamped backup, and re-validate JSON after
every write, auto-restoring on any breakage. Report what changed. Skip Low /
cosmetic. A small set of findings are **flag-only** and are *never* auto-edited
(secrets, mode conflicts) — see the rubric.

## Workflow

Create a TodoWrite item per step; keep one in_progress.

### 1 — Scope

Determine which config files are in play:

- Always: `~/.Codex/settings.json`, `~/.Codex/settings.local.json` (if present), `~/.Codex/mcp.json` (if present).
- If a repo is in scope (the user is working in one, or names one): also that repo's `./.Codex/settings.json` and `./.Codex/settings.local.json`.

Pass the repo path (or `none`) to the scanner.

### 2 — Scan (read-only)

Run the bundled fact-gatherer. It never mutates anything:

```bash
~/.Codex/skills/harness-streamliner/scripts/scan.sh [PROJECT_DIR|none]
```

It emits findings grouped by pass: permissions (exact dupes, clean
directory-prefix subsumption, autoMode↔permissions overlap), hooks (broken
script paths, per-event/startup/per-turn load), MCP (count, duplicate servers,
plaintext-secret heuristic), and env/model/statusline/misc (mode conflicts,
missing statusline script, budget fraction). Read its output as **candidates**,
not verdicts.

The five passes — what each looks for and the fix patterns — are detailed in
`references/passes.md`. Read it before classifying findings. The scanner is
deliberately conservative on glob subsumption; for fuzzier overlaps (bash
command globs, `**/`-suffix matchers) reason about the full rule list yourself,
since glob semantics are a judgment call.

### 3 — Classify by severity

Apply the rubric below to every candidate. Be honest: when in doubt, drop a
notch. The point is a leaner, correct config — not a high finding count.

### 4 — Apply (Medium+ only) — SAFELY

Follow `references/safety.md` exactly. The non-negotiable loop per file:

1. **Back up** → `cp settings.json settings.json.hso-bak.<UTC-timestamp>`.
2. **Edit** with `jq`-based transforms or surgical Edits (prefer `jq` so you
   can't corrupt structure).
3. **Validate** → `jq -e . <file>` must succeed. If it fails, **restore the
   backup** and report the failure; do not leave a broken config.
4. Re-run `scan.sh` to confirm the finding is gone and no new one appeared.

Never auto-edit flag-only findings. Never touch a file you didn't back up.

### 5 — Report

Produce a concise report:

```
# harness-streamliner — <UTC timestamp>
Scope: <files>
Backups: <paths>

## Applied (Medium+)
- [Pass N] <what> — <why> — <before> → <after>

## Flagged (you decide — not auto-applied)
- [secret] <where> — recommended: move to env / secret manager
- [conflict] <where> — recommended: <resolution>

## Skipped (Low / cosmetic)
- <one-line each>

Startup/per-turn cost: SessionStart=<n> hooks, UserPromptSubmit=<n> hooks. <note>
```

State plainly what was changed and verified. If a restore happened, lead with it.

## Severity rubric

| Severity | Action | Examples |
|---|---|---|
| **High** | auto-apply | Broken hook/statusline script path (silently failing every session); duplicate MCP server with identical command+args; exact-duplicate permission rules. |
| **Medium** | auto-apply | Permission rule cleanly subsumed by a broader same-tool rule; `autoMode.allow` entries already in `permissions.allow`; a hook whose script no longer exists; obviously-unused `env` var referencing a removed tool. |
| **Low** | skip (note only) | Stylistic ordering, key sort, comment-level nits, `skillListingBudgetFraction` within normal range. |
| **Flag-only** | NEVER auto-edit — report + recommend | Plaintext secrets/tokens in `mcpServers` args/env; conflicting `defaultMode` vs `permissions.defaultMode`; anything that changes the *security posture* (loosening `deny`/`ask`, broadening `allow`) or could break live auth. |

Two hard rules behind the rubric:

- **Never loosen security to "simplify."** Collapsing two redundant `deny`/`ask`
  rules is fine *only* if coverage is provably identical. If removing a rule
  could expose a path, keep it — redundant-but-safe beats lean-but-exposed.
- **Secrets and auth are the user's call.** Flag plaintext secrets and propose a
  fix (env indirection, secret manager), but don't rewrite a working MCP auth
  block — a bad edit there breaks the user's tooling mid-session.

## Boundaries — what this skill does NOT do

Defer; don't duplicate:

- **Skill content quality** (a skill's body, triggering, frontmatter) → `skill-optimizer` / `/sko`.
- **Skill tree / count / placement / cap balance / tiering** → `skill-tree-architect`, `/skill-tier`, `/sync-skills`.
- **Prompt quality** (system prompts, agent instructions) → `prompt-deep-optimizer` / `/pdo`.
- **The mechanics of a single setting** the user asks to change directly (add a permission, set an env var, configure one hook) → `update-config`. This skill is the *audit-and-trim sweep*; `update-config` is the precise single-change tool. Use `update-config`'s patterns for the actual writes when convenient.

If a scan surfaces a non-config problem (e.g., a skill that should be folded, a
bloated skill body), note it in the report and point at the right tool — don't
fix it here.

## Idempotence

Running twice should be a no-op the second time: a clean config yields an empty
"Applied" section. If a second run still finds Medium+ items, either the first
run failed to apply (check for a restore) or new cruft appeared — investigate
rather than re-applying blindly.