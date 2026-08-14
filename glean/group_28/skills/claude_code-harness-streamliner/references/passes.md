# Audit passes — checklist & fix patterns

Five passes, config/perf only. The scanner (`scripts/scan.sh`) surfaces
candidates mechanically; this file tells you what each means and how to fix it.
Classify with the SKILL.md rubric before applying. Prefer `jq` transforms for
writes — they preserve structure and can't half-corrupt a file the way a stray
text Edit can.

Throughout: `$F` is the settings file being fixed. Always work on a backup-first
basis (see `safety.md`).

---

## Pass 1 — Permissions

Targets `permissions.allow`, `permissions.deny`, `permissions.ask`, and
`autoMode.allow`.

**1a. Exact duplicates** — same string twice in one bucket. Always remove the
dupe (High). Dedupe a bucket while preserving order:

```bash
jq '.permissions.allow |= (reduce .[] as $x ([]; if index($x) then . else . + [$x] end))' "$F"
```

**1b. Clean directory subsumption** — a narrow rule fully covered by a broader
same-tool rule with a trailing `/**` or `/*`. The scanner reports these as
`SUBSUMED: <child> ⊂ <parent>`. Remove the child (Medium). Example from a real
config: `Write(~/.claude/skills/**)` and `Write(~/.claude/agents/*)` are both
covered by `Write(~/.claude/**)` → drop the two narrow ones.

Remove a specific set of rules by exact string. **Bind the element to a
variable** — `index()` after a `|` rebinds `.`, so `index(.)` silently matches
the wrong thing and can wipe the whole list. Always `. as $x`:

```bash
jq '.permissions.allow |= map(select(. as $x |
      ["Write(~/.claude/skills/**)","Edit(~/.claude/skills/**)"] | index($x) | not))' "$F"
```

After any allow/deny/ask edit, assert the bucket length dropped by exactly the
number you intended (`jq ".permissions.allow|length"`) before validating JSON —
a structurally-valid file with the wrong rule count is the dangerous failure
mode (valid JSON, broken access). Restore if the count is off.

Only act on subsumptions the scanner reports OR ones you can prove. Do **not**
collapse `**/`-suffix matchers (e.g. `Read(**/.env)` vs `Read(**/.env.*)`) on a
hunch — they match different sets. When unsure, keep both.

**1c. autoMode.allow ↔ permissions.allow overlap** — entries in `autoMode.allow`
that already appear in `permissions.allow` are redundant (auto mode already
inherits the base allowlist). Remove the duplicated entries from `autoMode.allow`
(Medium). If that empties `autoMode.allow`, drop the `autoMode` key.

**1d. Security direction (HARD RULE)** — `deny` and `ask` are safety nets.
Removing a "redundant" deny/ask rule is allowed **only** when another rule
provably covers the exact same paths. If there's any doubt the path stays
protected, keep the rule. Redundant-but-safe > lean-but-exposed. Never move a
rule from `deny`→`ask` or `ask`→`allow` as "cleanup" — that's a posture change,
which is flag-only.

---

## Pass 2 — Hooks

Targets `hooks.<Event>[].hooks[].command`.

**2a. Broken script paths (High)** — the scanner marks each command `ok`,
`MISSING`, or `inline/no-path`. A `MISSING` script means that hook fails silently
every time its event fires. Decide with the user's intent: if the script was
moved, fix the path; if it's dead, remove that hook entry. Don't guess a new
path — confirm the file exists first.

**2b. Per-event / startup / per-turn load** — count hooks per event:

- `SessionStart` hooks run on every session open → **startup latency**. Anything
  slow here (big `timeout`, heavy `node` process) delays first response.
- `UserPromptSubmit` hooks run on **every turn** → recurring latency *and* often
  inject context into the prompt, spending token budget each turn. A hook that
  appends a large block every turn is the most expensive kind.
- Duplicate commands across events, or the same command registered twice in one
  event, are removable (Medium for exact dupes).

You generally **cannot** measure wall-clock cost from config alone, and these
hooks are usually intentional (status integrations, tiering, context injectors).
So: **report** the load and flag obviously-wasteful entries; auto-remove only
exact duplicates and dead-script hooks. Don't delete a working hook because it
"looks heavy" — recommend, let the user decide.

**2c. Missing timeout on an external command** — a `command` hook that shells out
(node/npx/network) with no `timeout` can hang a session. Adding a sane `timeout`
(e.g. 10–30s) is a safe Medium fix:

```bash
# add timeout:10 to a specific hook lacking one — do surgically, validate after
```

---

## Pass 3 — MCP servers

Targets `mcpServers` in `settings.json` and `~/.claude/mcp.json`.

**3a. Server count (context cost)** — every configured MCP server can load its
tool list into context. Report the count. If a server is clearly unused/unauthed
(e.g. the user says they don't use it, or it points at a dead local port), note
it as a candidate for removal or deferral — but confirm before removing; a server
the user forgot about may still be wanted.

**3b. Duplicate / near-identical servers (High)** — two entries with the same
`command`+`args` signature (the scanner flags these). Real example: two
`@mondaydotcomorg/monday-api-mcp` servers with the same token differing only in
`--mode`. Keep the one that's actually used; remove the other. Confirm which
before deleting.

**3c. Plaintext secrets (FLAG-ONLY — never auto-edit)** — the scanner flags
JWT/`sk_`/`api_key`-shaped strings in `args` and `env` values. These are a
security smell, but rewriting a working auth block risks breaking the user's
tooling mid-session, and the "fix" (env var / secret manager indirection) depends
on their environment. So: **report each, recommend env indirection, and stop.**
Do not edit the secret in place.

---

## Pass 4 — env / model / statusline / misc

**4a. Conflicting defaultMode (FLAG-ONLY)** — top-level `defaultMode` and
`permissions.defaultMode` set to different values is almost always unintended,
but resolving it changes permission behavior → flag, recommend, let the user pick.

**4b. Missing statusline script (High if broken)** — if `statusLine.command`
points at a script that doesn't exist, the statusline silently breaks. Fix the
path or remove the block, per user intent.

**4c. Stale env vars (Medium)** — an `env` entry referencing a tool/path that no
longer exists is removable. A var you can't classify → leave it; env vars are
cheap and removing a needed one is worse than keeping a dead one.

**4d. skillListingBudgetFraction (Low)** — note if outside ~0.05–0.20. Higher
puts more skill metadata in context (better triggering, more tokens); lower saves
tokens but skills may under-trigger. This is a tuning preference, not a defect —
report, don't change without being asked.

**4e. model (note only)** — report the configured model; never change it.

---

## What NOT to touch here

- Skill files, the skill tree, prompts → defer (see SKILL.md Boundaries).
- Anything that changes security posture → flag-only.
- Plugin enable/disable state (`enabledPlugins`) — report obviously-broken
  marketplace refs, but enabling/disabling plugins is a user preference, not cruft.
