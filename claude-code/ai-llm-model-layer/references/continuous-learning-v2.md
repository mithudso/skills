<!-- hub-reference-banner -->
> **Reference file — part of the `ai-llm-model-layer` hub.** Formerly the standalone `continuous-learning-v2` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: continuous-learning-v2
description: >
  Hook-based learning system for Claude Code that observes sessions via PreToolUse/PostToolUse hooks, extracts atomic "instincts" with confidence scoring, and evolves them into skills, commands, and agents. Use when setting up persistent cross-session learning, managing instincts, or evolving observed patterns into reusable artifacts.
  TRIGGER: "continuous learning", "instinct", "observe sessions", "learn from my workflow", "persist behavior across sessions", "/instinct-status", "/evolve", "hook-based learning", "confidence scoring", "learn my patterns".
  SKIP: project-specific skills that don't need session observation (place in ~/.claude/skills/ directly); standard Claude Code use without cross-session persistence; Claude Code versions below v2.1 (PreToolUse/PostToolUse hooks unsupported).
version: 2.0.0
updated: "2026-05-29"
category: meta
tags: [learning, instincts, hooks, confidence-scoring, self-improvement, claude-code]
related_skills:
  - claude-code-skills
  - skill-optimizer
  - prompt-deep-optimizer
whenToUse:
  - setting up persistent learning that survives across Claude Code sessions
  - reviewing or managing learned instincts and their confidence scores
  - evolving clusters of instincts into skills, commands, or agents
  - importing or exporting instincts between machines or users
  - troubleshooting observation hooks or the homunculus pipeline
whenNotToUse:
  - no need for cross-session learning (standard Claude Code is sufficient)
  - only project-specific skills needed (place directly in ~/.claude/skills/)
  - Claude Code below v2.1 (hooks not supported)
---

# Continuous Learning v2 — Instinct-Based Architecture

Transforms Claude Code sessions into reusable knowledge through small learned behaviors called "instincts," each with confidence scoring, domain tagging, and an evolution path into skills/commands/agents.

## Requirements

- **Claude Code v2.1+** — PreToolUse/PostToolUse hook support required
- **Haiku API access** — Observer agent runs background analysis (`observer.model` in `~/.claude/homunculus/config.json`)

## v2 vs v1

| Feature | v1 | v2 |
|---------|----|----|
| Observation | Stop hook (session end) | PreToolUse/PostToolUse (100% reliability) |
| Analysis | Main context | Background Haiku agent |
| Granularity | Full skills | Atomic instincts |
| Confidence | None | 0.3–0.9 weighted |
| Evolution path | Direct to skill | Instinct → cluster → skill/command/agent |
| Sharing | None | Export/import instincts |

## Instinct Model

An instinct is a single learned behavior:

```yaml
---
id: prefer-functional-style
trigger: "when writing new functions"
confidence: 0.7
domain: "code-style"
source: "session-observation"
---

# Prefer functional style

## Action
Use functional patterns over classes when appropriate.

## Evidence
- Observed functional preference 5 times
- User corrected class-based approach to functional on 2025-01-15
```

Properties:
- **Atomic** — one trigger, one action
- **Confidence-weighted** — 0.3 = tentative, 0.9 = near-certain
- **Domain-tagged** — code-style, testing, git, debugging, workflow, etc.
- **Evidence-backed** — tracks the observations that created it

## Pipeline

```
Session activity
      │ Hooks capture prompt + tool use (100% reliability)
      ▼
observations.jsonl
      │ Observer agent reads (background, Haiku)
      ▼
Pattern detection
  • User corrections → instinct
  • Error resolutions → instinct
  • Repeated workflows → instinct
      │ create/update
      ▼
instincts/personal/
  • prefer-functional.md (0.7)
  • always-test-first.md (0.9)
      │ /evolve cluster
      ▼
evolved/
  • commands/new-feature.md
  • skills/testing-workflow.md
  • agents/refactor-specialist.md
```

## Quick Start

### 1. Enable observation hooks

**Installed as a plugin (recommended):** the plugin's `hooks/hooks.json` auto-registers — no manual settings.json changes needed. If you previously copied `observe.sh` into `~/.claude/settings.json`, remove those duplicate `PreToolUse`/`PostToolUse` blocks to prevent double-execution and `${CLAUDE_PLUGIN_ROOT}` resolution errors.

**Installed manually to `~/.claude/skills`:** add to `~/.claude/settings.json`:

```json
{
  "hooks": {
    "PreToolUse":  [{ "matcher": "*", "hooks": [{ "type": "command", "command": "~/.claude/skills/continuous-learning-v2/hooks/observe.sh" }] }],
    "PostToolUse": [{ "matcher": "*", "hooks": [{ "type": "command", "command": "~/.claude/skills/continuous-learning-v2/hooks/observe.sh" }] }]
  }
}
```

### 2. Initialize directories

```bash
mkdir -p ~/.claude/homunculus/{instincts/{personal,inherited},evolved/{agents,skills,commands}}
touch ~/.claude/homunculus/observations.jsonl
```

### 3. Use instinct commands

| Command | Action |
|---------|--------|
| `/instinct-status` | Show all learned instincts with confidence |
| `/evolve` | Cluster related instincts into skill/command/agent |
| `/instinct-export` | Export instincts for sharing |
| `/instinct-import <file>` | Import instincts from another user |

## Configuration

`~/.claude/homunculus/config.json`:

```json
{
  "version": "2.0",
  "observation": {
    "enabled": true,
    "store_path": "~/.claude/homunculus/observations.jsonl",
    "max_file_size_mb": 10,
    "archive_after_days": 7
  },
  "instincts": {
    "personal_path": "~/.claude/homunculus/instincts/personal/",
    "inherited_path": "~/.claude/homunculus/instincts/inherited/",
    "min_confidence": 0.3,
    "auto_approve_threshold": 0.7,
    "confidence_decay_rate": 0.05
  },
  "observer": {
    "enabled": true,
    "model": "haiku",
    "run_interval_minutes": 5,
    "patterns_to_detect": ["user_corrections", "error_resolutions", "repeated_workflows", "tool_preferences"]
  },
  "evolution": {
    "cluster_threshold": 3,
    "evolved_path": "~/.claude/homunculus/evolved/"
  }
}
```

## File Structure

```
~/.claude/homunculus/
├── identity.json           # Profile, tech level
├── observations.jsonl      # Current session observations
├── observations.archive/   # Processed observations
├── instincts/
│   ├── personal/           # Auto-learned instincts
│   └── inherited/          # Imported from others
└── evolved/
    ├── agents/             # Generated specialist agents
    ├── skills/             # Generated skills
    └── commands/           # Generated commands
```

## Confidence Scoring

| Score | Meaning | Behavior |
|-------|---------|----------|
| 0.3 | Tentative | Suggested but not enforced |
| 0.5 | Moderate | Applied when relevant |
| 0.7 | Strong | Auto-approved on application |
| 0.9 | Near-certain | Core behavior |

**Confidence rises** when a pattern repeats across sessions, when the user doesn't correct suggested behavior, or when similar instincts from other sources agree.

**Confidence falls** when the user explicitly corrects behavior, when the pattern hasn't been observed recently, or when contradicting evidence appears.

## Why hooks instead of skills for observation?

Skills fire probabilistically (~50–80%) based on Claude's judgment. Hooks fire deterministically — 100% of tool calls are observed, no patterns missed, learning is complete.

## Backward Compatibility

v2 is fully compatible with v1:
- Existing `~/.claude/skills/learned/` skills continue to work
- Stop hooks continue to run (and feed v2)
- Both systems can run in parallel during migration

## Privacy

- Observations stay **local** on your machine
- Only **instincts** (patterns) are exportable — not actual code or conversation content
- You control exactly what gets exported

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| Observations not recording | Hooks not registered or plugin not loaded | Check `~/.claude/settings.json` hooks config or reinstall plugin |
| `${CLAUDE_PLUGIN_ROOT}` error | Hooks registered in both `settings.json` and plugin `hooks.json` | Remove the `PreToolUse`/`PostToolUse` blocks from `settings.json` |
| `observations.jsonl` growing large | `archive_after_days` too long or archiving stopped | Reduce `archive_after_days` in `config.json` (default: 7) |
| Observer not running | Haiku API disabled or `observer.enabled: false` | Set `observer.enabled: true`; verify API key |
| Instinct confidence not rising | Pattern not re-observed enough | Repeat the pattern across multiple sessions or adjust confidence manually |

## Integration with Skill Creator

[Skill Creator](https://skill-creator.app) generates both traditional `SKILL.md` files (backward-compatible) and instinct collections (for v2). Repository-analysis instincts include `source: "repo-analysis"` with the source repository URL.
