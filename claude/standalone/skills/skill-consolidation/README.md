# Skill Consolidation — clone-and-go

A zero-dependency routing index for a large library of Claude **skills**. It joins every
skill's `TRIGGER`/`SKIP` metadata into one cross-family index so an agent can pick **one**
skill without reading hundreds of files, plus an optional local-embedding **semantic search**
layer. Runs on Node.js (≥18) with **no npm install**.

Full design + operations reference: [`INDEX-ARCHITECTURE.md`](./INDEX-ARCHITECTURE.md)
(see its **How-To Guide** section for day-to-day usage).

## Quickstart

**Standalone bootstrap** — no clone yet. Installs Homebrew + Git, clones the repo, then configures:

```bash
bash setup.sh --repo <git-url>            # clones into ~/.claude/skill-consolidation (override with --dir)
```

**Already cloned:**

```bash
git clone <git-url> ~/.claude/skill-consolidation
cd ~/.claude/skill-consolidation
./setup.sh                                # idempotent; safe to re-run
```

`setup.sh` runs ten steps (each skippable): install **Homebrew + Git** → **clone/update** the repo
→ Node check → **prompt for API keys** (Monday required; Firecrawl + Exa optional) and write `.env`
→ Ollama + `qwen3-embedding:4b` as a boot service → build `SKILLS-INDEX.*` → embed → install the
boot refresh agent → merge MCP servers into `~/.claude.json` → install a **nightly cron** that
`git pull`s and re-embeds. A server whose key you skip stays **disabled**. Restart Claude Code
afterward to pick up the MCP servers.

```bash
./setup.sh --repo <url>      # clone (standalone) / used by the nightly cron pull
./setup.sh --skip-keys       # don't prompt; edit .env by hand
./setup.sh --skip-ollama     # index only, no semantic layer
./setup.sh --skip-service    # don't install the boot agent
./setup.sh --skip-mcp        # don't touch ~/.claude.json (also skips key prompts)
./setup.sh --skip-cron       # don't install the nightly git-pull job
```

## What ships vs. what's generated locally

| Committed (clone gets it) | Generated locally / git-ignored |
| --- | --- |
| `gen-skills-index.mjs`, `gen-skills-index.test.mjs`, `bench-search.mjs` | `SKILLS-EMBEDDINGS.json` (18 MB — built by `--embed`) |
| `SKILLS-INDEX.json` + `SKILLS-INDEX.md` (instant keyword routing) | rendered `com.skills-embed.plist` / systemd units |
| `INDEX-ARCHITECTURE.md`, `README.md`, the `*.template` files | `.env` (your secrets), `*.log`, `*.jsonl`, `backups/`, `run-state/` |
| `setup.sh`, `.gitignore`, `.env.example`, `mcp-servers.template.json` | |

The index is **always re-derived** from `~/.claude/skills` — never hand-edited. The vector
corpus is locked to one embedding model, so it's regenerated per machine rather than shipped.

## Using it

```bash
node gen-skills-index.mjs                 # rebuild the routing table
open SKILLS-INDEX.md                       # keyword routing — no server needed
node gen-skills-index.mjs --search "..."   # semantic ranking (fails open to keyword-only)
node gen-skills-index.mjs --check          # CI drift gate: exit 1 if stale
node --test gen-skills-index.test.mjs      # 14/14
```

`--search` **never errors** if the embedding server is down — it prints a one-line
`degraded to keyword-only` note to stderr and returns ranked keyword results with exit 0.

## MCP servers

`setup.sh` reads `.env`, renders `mcp-servers.template.json`, and merges each server into your
`~/.claude.json`:

- **Non-destructive** — a server you already have is kept, never overwritten.
- **Backed up** — your `~/.claude.json` is copied to `~/.claude.json.bak-<ts>` before any write.
- **Credential-gated** — any server whose `${VAR}` is still unset in `.env` is **skipped**, so
  you only get the servers you have tokens/paths for.
- **Secret-safe** — real values live only in `.env` (git-ignored) and your private
  `~/.claude.json`; nothing secret is committed or printed (only server names are logged).

The template covers: `glean_default`, `granola`, `headroom`, `monday-access-mcp`,
`monday-apps-mcp`, `firecrawl`, `exa`, `tam_mcp`, `mdb_tam_account_context`,
`mdb_case_assistant`, and `skills-relay`.

## Persistence (auto-refresh)

The boot agent keeps the index + vectors in sync whenever a skill changes and every 6 hours.

- **macOS** — a launchd agent at `~/Library/LaunchAgents/com.skills-embed.plist`
  (`RunAtLoad` + `WatchPaths` on `~/.claude/skills` + 6h interval). Ollama runs as a Homebrew
  service so the embedding server is up after every reboot.
- **Linux** — systemd **user** units `skills-embed.{service,timer,path}` (timer at boot + 6h,
  path watch on `~/.claude/skills`). Start Ollama at boot yourself (login item or its own unit).

Both render from the committed `*.template` files, so no machine-specific paths are committed.

## Requirements

- **Node.js ≥ 18** (built-in `fetch` + test runner). No npm dependencies.
- **Ollama** + `qwen3-embedding:4b` (free, no API key) — only for `--search`. Everything else
  works without it.
- macOS (Homebrew) or Linux (systemd `--user`) for the boot agent; otherwise run the refresh
  from cron or by hand.

## Troubleshooting

See the **Troubleshooting** table in [`INDEX-ARCHITECTURE.md`](./INDEX-ARCHITECTURE.md).
Common one: `--search` says *degraded to keyword-only* → `brew services start ollama`
(search keeps working meanwhile).
