<!-- hub-reference-banner -->
> **Reference file — part of the `claude-code-skills` hub.** Formerly the standalone `claude-code-plugins` skill.
> Sibling topics in this family are now reference files under the hubs (`claude-code-skills`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: claude-code-plugins
description: >-
  Claude Code plugin system expert — plugin anatomy (plugin.json, directory structure),
  hooks (29 events, 5 handler types, matchers, exit codes), commands/skills (unified system,
  frontmatter, argument parsing), agent definitions (.md with frontmatter, worktree isolation),
  plugin distribution (marketplaces, validation, submission), and settings (userConfig, secrets,
  env vars). Use when creating, debugging, publishing, or configuring plugins, hooks, commands,
  or custom agents.
  TRIGGER: "create a Claude Code plugin", "configure a hook", "PreToolUse hook",
  "plugin.json manifest", "slash command in a plugin", "custom subagent", "publish a plugin",
  "plugin validate", "hook exit codes", "worktree isolation for agents".
  SKIP: writing a standalone skill not packaged in a plugin (use claude-code-skills);
  general MCP server development (use mcp-server-dev:build-mcp-server).
origin: local
version: 1.1.1
category: developer
updated: 2026-05-31
whenToUse:
  - "create a Claude Code plugin"
  - "how do I configure a PreToolUse hook"
  - "plugin.json manifest fields"
  - "add a slash command to my plugin"
  - "define a custom subagent in a plugin"
  - "publish my plugin to a marketplace"
  - "hook handler types explained"
  - "hook exit codes and blocking behavior"
  - "worktree isolation for plugin agents"
  - "validate a plugin with claude plugin validate"
  - "plugin settings userConfig and secrets"
related_skills:
  - claude-code-skills
  - ai-agent-engineering
  - hookify:hookify
---

# Claude Code Plugins, Hooks & Commands

Reference for the Claude Code extension system: plugins, hooks, commands, and agents. Backed by `references/claude-code-plugins-context.md`.

## When NOT to use

- Writing a standalone skill not packaged in a plugin → use `claude-code-skills`
- General MCP server development → use `mcp-server-dev:build-mcp-server`
- Hook configuration in project settings only (no plugin packaging) → use the `hookify` skill

## Quick reference

### Plugin directory structure

```
my-plugin/
  .claude-plugin/plugin.json    # Manifest (ONLY this goes here)
  skills/<name>/SKILL.md        # Skills (auto-discovered)
  agents/<name>.md              # Subagent definitions
  hooks/hooks.json              # Hook configurations
  .mcp.json                     # MCP server definitions
  bin/                          # Executables added to PATH
  settings.json                 # Defaults (agent, subagentStatusLine only)
```

### Minimal plugin.json manifest

```json
{
  "name": "my-plugin",
  "version": "1.0.0",
  "description": "What this plugin does",
  "author": "Your Name",
  "license": "MIT"
}
```

Required fields: `name`, `version`. All others are recommended for marketplace submission.

### Hook events (key subset of 29 total)

| Event | When | Blocking? |
|---|---|---|
| `PreToolUse` | Before tool executes | Yes (exit 2 = deny) |
| `PostToolUse` | After tool completes | No |
| `UserPromptSubmit` | User sends message | Yes (exit 2 = block) |
| `SessionStart` | Session begins | No |
| `Stop` | Claude finishes turn | Yes (exit 2 = prevent stop) |

### Hook handler types

| Type | Mechanism | Use when |
|---|---|---|
| `command` | Shell script, JSON on stdin | Linting, formatting, validation |
| `http` | POST to endpoint | External service integration |
| `mcp_tool` | Call MCP server tool | Memory, logging, analytics |
| `prompt` | Single-turn LLM yes/no | Security review, semantic checks |
| `agent` | Subagent with tools | Cross-file analysis, compliance |

### Hook exit codes

| Code | Meaning |
|---|---|
| `0` | Success — continue |
| `1` | Non-blocking error — logs and continues |
| `2` | **Blocking** — denies tool call / blocks prompt / prevents stop |

### Agent frontmatter fields

| Field | Purpose |
|---|---|
| `name` | Unique identifier (required) |
| `description` | When to delegate (required) |
| `model` | `sonnet` / `opus` / `haiku` / `inherit` |
| `isolation` | `worktree` for git worktree isolation |
| `tools` / `disallowedTools` | Tool allowlist/denylist |
| `skills` | Skills preloaded at startup |
| `maxTurns` | Turn budget |
| `permissionMode` | `default` / `acceptEdits` / `auto` / `plan` |

## Full reference

For complete coverage of all 29 hook events, hook matchers, plugin.json schema, marketplace submission, and settings configuration, read `references/claude-code-plugins-context.md` in this directory.

## 2026-06 Delta (added 2026-06-10)

- **Marketplace-less distribution (2.1.157, 2026-05-29) [HIGH]:** plugins in `.claude/skills/` auto-load with no marketplace and no `--plugin-dir`; `claude plugin init <name>` (`--with skills,agents,hooks,mcp`) scaffolds; a skill folder containing `.claude-plugin/plugin.json` loads as plugin `<name>@skills-dir` and may bundle agents/hooks/MCP (project-level requires workspace trust). Committing `.claude/skills/` to a repo IS plugin distribution now — which widens the repo-borne attack surface (Datadog 2026-05-11: dynamic-context shell commands execute before the model sees the skill; mitigation `"disableSkillShellExecution": true`).
- **plugin.json / marketplace capabilities [HIGH]:** `defaultEnabled: false` (2.1.154); `skipLfs` on github/git sources (2.1.153); `pluginSuggestionMarketplaces` managed allowlist (2.1.152); directory-aware Discover pins (2.1.154); `/plugin list --enabled/--disabled` (2.1.163); marketplace search bar (2.1.172); `marketplace remove --scope user|project|local` (2.1.152).
- **Anthropic's two-marketplace structure [MEDIUM]:** `claude-plugins-official` (curated, auto-available) + `claude-plugins-community` (review pipeline, commit-SHA-pinned, CI auto-bumps, nightly catalog sync).
- **Ecosystem scale (mostly community-scraped) [MEDIUM]:** skills.sh ~90.7k skills / 43.6M installs; ~32k plugins / 282k components with skills = 57% of components, hooks 2.4%; Snyk ToxicSkills baseline: 36.8% of scanned skills have ≥1 flaw, 13.4% critical.
- **Subagents 5 levels deep (2.1.172); self-hosted runner `post-session` hook (2.1.169).** [HIGH]
