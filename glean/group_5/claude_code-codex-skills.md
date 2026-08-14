# Codex-skills

**Category:** Cloud, DevOps & Infrastructure
**Platform:** Claude Code
**Original Path:** claude-code/claude-code-skills

## Description
Expert reference for Codex skills — anatomy (SKILL.md, frontmatter, references/), authoring, discovery, distribution (plugins, marketplaces, security), management (precedence, token budgets, skillOverrides), composition (hooks, agents), optimization (quality gates, CI/CD), and Codex workflows (plan mode, TDD/verification loops, hooks automation, headless/CI, worktrees, context management). Use when creating, editing, publishing, or debugging skills, or applying Codex workflow patterns. TRIGGER: "write a SKILL.md", "skill frontmatter fields", "skill not triggering", "install a skill", "skill marketplace", "skill token budget", "skill precedence", "skill security audit", "Codex workflow / plan mode / headless / worktrees". SKIP: plugin packaging and hook configuration (use Codex-plugins); searching the prompts.chat registry to install a skill (use skill-lookup); creating a skill from scratch interactively (use skill-creator).

---

# Codex Skills Expert

Ref: Codex skills ecosystem — anatomy, authoring, discovery, distribution, management, composition, optimization. Backed by `references/Codex-skills-context.md`.

## When NOT to use

- Plugin packaging, hook config, agent definitions → use `Codex-plugins`
- Search prompts.chat registry for existing skill → use `skill-lookup`
- Create new skill from description → use `skill-creator`
- Fix trigger accuracy, over-triggering, prose quality of existing skill → use `skill-optimizer`

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
2. Personal (`~/.Codex/skills/`)
3. Project (`.Codex/skills/`)
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
| Skill anatomy, authoring, distribution depth | Full frontmatter/token-budget/skillOverrides/CI-CD coverage | `references/Codex-skills-context.md` |
| Plugin packaging, hooks config, agent definitions | Plugin-level work | `references/Codex-plugins.md` |
| Codex workflows | Plan mode, verification/TDD loops, hooks automation, headless/CI, Agent SDK routines, worktree/subagent/team parallelism, context management, workflow failure modes | `references/Codex-workflows.md` |
| Model-migration & prompt/skill portability | Migrating prompt/skill library across model versions — prompt rot, capability-shift regression, regression-eval/golden-transcript/canary gating, re-optimization order, pinning vs floating | `references/model-migration-prompt-portability.md` |

## Full reference

Complete coverage: all frontmatter fields, token budget management, skillOverrides, CI/CD integration, skill composition patterns → `references/Codex-skills-context.md`.