# skills

Clone-and-go Claude skills, commands, and a local semantic skill index.

## Quickstart

```bash
git clone https://github.com/mithudso/skills.git ~/.claude/skills
cd ~/.claude/skills && ./setup.sh
```

`setup.sh` is idempotent and safe to re-run. It installs dependencies
(Homebrew/Git/Node), prompts for API keys (written to `~/.claude/.env`, never
committed), lays down `~/.claude` config non-destructively, starts Ollama +
the embedding model, builds and embeds the skill index, bridges the legacy
`~/.claude/skill-consolidation` path to the vendored copy, merges MCP servers
(keyless servers stay disabled), and installs a nightly `git pull` + re-embed
cron job. Flags: `--skip-ollama --skip-service --skip-mcp --skip-keys
--skip-config --skip-cron`.

## Useful commands

- `/cfe <subject>` — Concept Family Explorer (map gaps, loop `/dr`, saturate)
- `/dr <topic>` — deep research that ends in an installed skill
- `/cdo`, `/ddo`, `/pdo`, `/sko`, `/dqo`, `/deso` — the deep-optimizer family
- `node skill-consolidation/gen-skills-index.mjs --search "..."` — semantic search

## Docs

Deeper technical docs live in [`skill-consolidation/`](skill-consolidation/):
`README.md` (index engine) and `INDEX-ARCHITECTURE.md` (full spec).
