---
name: claude-code-skills
description: >-
  Expert reference for Claude Code skills — anatomy (SKILL.md, frontmatter, references/), authoring, discovery, distribution (plugins, marketplaces, security), management (precedence, token budgets, skillOverrides), composition (hooks, agents), optimization (quality gates, CI/CD), and Claude Code workflows (plan mode, TDD/verification loops, hooks automation, headless/CI, worktrees, context management). Use when creating, editing, publishing, or debugging skills, or applying Claude Code workflow patterns. TRIGGER: "write a SKILL.md", "skill frontmatter fields", "skill not triggering", "install a skill", "skill marketplace", "skill token budget", "skill precedence", "skill security audit", "Claude Code workflow / plan mode / headless / worktrees". SKIP: plugin packaging and hook configuration (use claude-code-plugins); searching the prompts.chat registry to install a skill (use skill-lookup); creating a skill from scratch interactively (use skill-creator).
origin: local
version: 1.2.1
category: developer
updated: 2026-07-20
whenToUse:
  - "how do I write a SKILL.md"
  - "skill frontmatter fields and options"
  - "my skill isn't triggering"
  - "skill description best practices"
  - "how to publish or distribute a skill"
  - "skill marketplace comparison"
  - "skill token budget and skillOverrides"
  - "skill precedence order"
  - "skill security — ToxicSkills, mcp-scan"
  - "compose a skill with hooks or agents"
  - "optimize a skill with prompt-deep-optimizer"
  - "Claude Code workflow patterns — plan mode, TDD loops, hooks, headless, worktrees"
  - "migrate or re-optimize prompts/skills across model versions — prompt rot, regression"
related_skills:
  - misc-catch-all
  - skill-optimizer
  - prompt-deep-optimizer
  - skill-creator
---
# Claude Code Skills Expert

Ref: Claude Code skills ecosystem — anatomy, authoring, discovery, distribution, management, composition, optimization. Backed by `references/claude-code-skills-context.md`.

## When NOT to use

- Plugin packaging, hook config, agent definitions → use `claude-code-plugins`
- Search prompts.chat registry for existing skill → use `skill-lookup`
- Create new skill from description → use `skill-creator`
- Fix trigger accuracy, over-triggering, prose quality of existing skill → use `skill-optimizer`
- Whole-tree taxonomy rebalance, hub cap-balance, cross-hub placement, reshape for a new family → use `skill-tree-architect`

## Quick reference

### SKILL.md structure

```
---
name: kebab-case-name
description: Verb-first, 1-3 sentences, include trigger phrases and exclusions
origin: local
---

# Skill Title

## When to use this skill
- Specific trigger conditions

## Instructions / Reference content
```

### Key frontmatter fields

| Field | Purpose |
|---|---|
| `name` | Max 64 chars, kebab-case (required) |
| `description` | Max 1,024 chars — primary activation signal |
| `disable-model-invocation` | `true` prevents auto-trigger from description — only runs on explicit `/skill-name` |
| `user-invocable` | `false` hides from slash menu |
| `allowed-tools` | Pre-approve tools for skill execution |
| `context` | `fork` for isolated subagent execution |
| `agent` | Target agent type (e.g., `Explore`) |
| `hooks` | Skill-scoped lifecycle hooks |
| `paths` | Glob patterns limiting activation to matching files |

### Skill precedence (highest to lowest)

1. Enterprise managed settings
2. Personal (`~/.claude/skills/`)
3. Project (`.claude/skills/`)
4. Plugin (`plugin-name:skill-name`)

### Marketplace comparison

| Marketplace | Size | Security | Best for |
|---|---|---|---|
| Anthropic Official | ~20 | Verified by Anthropic | Trust-first installs |
| Agensi | 200+ | 8-point security scan | Vetted + creator payments |
| ClaudeSkills.info | 658+ | Community-vetted | Free, quality community skills |
| Skills.sh | ~2,000 | No audit | CLI-first, npm-style workflow |
| SkillsMP | 800K+ | No audit | Discovery only — audit before use |

### Security checklist

- Run `uvx mcp-scan@latest --skills` before installing unknown skills
- 36.8% scraped skills have security flaws (Snyk ToxicSkills, Feb 2026)
- 13.4% critical vulns; 76 confirmed malicious payloads found
- Use Anthropic Official or Agensi-vetted for production

## Sub-topic routing

| Sub-topic | When to load | Reference file |
| --- | --- | --- |
| Skill anatomy, authoring, distribution depth | Full frontmatter/token-budget/skillOverrides/CI-CD coverage | `references/claude-code-skills-context.md` |
| Plugin packaging, hooks config, agent definitions | Plugin-level work | `references/claude-code-plugins.md` |
| Claude Code workflows | Plan mode, verification/TDD loops, hooks automation, headless/CI, Agent SDK routines, worktree/subagent/team parallelism, context management, workflow failure modes | `references/claude-code-workflows.md` |
| Model-migration & prompt/skill portability | Migrating prompt/skill library across model versions — prompt rot, capability-shift regression, regression-eval/golden-transcript/canary gating, re-optimization order, pinning vs floating | `references/model-migration-prompt-portability.md` |
| Headless `-p` / stream-json **client** contract | Building a program on top of `claude`: NDJSON event types, `stream_event` SSE passthrough, the four terminal states, why `result` is neither guaranteed nor last, in-band `control_request` interrupt (cancel requires `--input-format stream-json`), session resume + the cwd-mismatch trap, SIGTERM→143, the 3s open-stdin tax, `CLAUDECODE` nested-hang, pipe-vs-PTY buffering | `references/claude-headless-streaming.md` |

## Full reference

Complete coverage: all frontmatter fields, token budget management, skillOverrides, CI/CD integration, skill composition patterns → `references/claude-code-skills-context.md`.