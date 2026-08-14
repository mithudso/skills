# Claude Code Skills Ecosystem — Reference Context

Comprehensive reference for skill anatomy, authoring, discovery, distribution, management, composition, and optimization. Compiled from 20+ sources, May 2026.

---

## 1. Skill Anatomy

A skill is a directory containing at minimum `SKILL.md`. Directory name must be kebab-case.

### SKILL.md structure

Two parts: YAML frontmatter between `---` markers, and markdown body with instructions.

**Frontmatter fields (all optional except description recommended):**

| Field | Type | Purpose |
|-------|------|---------|
| `name` | string | Max 64 chars, kebab-case. Defaults to directory name |
| `description` | string | Max 1,024 chars. Primary activation signal |
| `disable-model-invocation` | bool | `true` = manual `/skill` only, Claude cannot auto-activate |
| `user-invocable` | bool | `false` = hidden from slash menu |
| `allowed-tools` | string[] | Pre-approve tools for frictionless execution |
| `context` | string | `fork` for isolated subagent execution |
| `agent` | string | Target agent type (e.g., `Explore`, `Plan`) |
| `model` | string | Override model for this skill |
| `effort` | string | Override effort level |
| `hooks` | object | Skill-scoped lifecycle hooks |
| `paths` | string[] | Glob patterns limiting activation to matching files |
| `arguments` | object[] | Named arguments with descriptions |
| `argument-hint` | string | Hint shown in autocomplete |
| `shell` | string | Shell for script execution |

### Directory structure

```
my-skill/
  SKILL.md           # Required entrypoint
  references/        # Reference material loaded on demand
  examples/          # Examples loaded on demand
  scripts/           # Executed, not loaded into context
```

### Progressive loading

At session start, Claude sees only name + description (~100 tokens per skill). Full SKILL.md body loads only when Claude decides the skill is relevant. Reference files load on demand. After auto-compaction, skills re-attach with first 5,000 tokens each, sharing a 25,000-token budget prioritized by most-recent invocation.

### String substitutions

| Variable | Expands to |
|----------|-----------|
| `$ARGUMENTS` | Full argument string |
| `$ARGUMENTS[N]` / `$N` | Nth positional argument |
| `$name` | Named argument from `arguments` frontmatter |
| `${CLAUDE_SESSION_ID}` | Current session ID |
| `${CLAUDE_EFFORT}` | Current effort level |
| `${CLAUDE_SKILL_DIR}` | Skill directory path |
| `` !`command` `` | Dynamic context — shell command output injected before Claude sees content |

---

## 2. Skill Authoring

### Description writing (most impactful decision)

- Write in third person: "Processes Excel files and generates reports"
- Never first/second person
- Include what the skill does AND specific trigger phrases
- End with explicit exclusions: "Do NOT use for blog articles, newsletters, emails"
- Be "a little pushy" — Claude has a measured tendency to under-trigger
- Optimized descriptions improve activation from 20% → 50%. Adding examples raises to 90%.

### Content guidance

- **Reference content** (conventions, patterns, style guides): inline in SKILL.md
- **Task content** (deployments, code generation): step-by-step instructions
- Use `disable-model-invocation: true` for destructive operations
- Target under 500 lines / 5,000 tokens for SKILL.md body
- Move reference material to `references/` for on-demand loading

### Degrees of freedom framework

| Freedom level | Format | Use when |
|--------------|--------|----------|
| High | Text-based heuristics | Context-dependent tasks |
| Medium | Pseudocode with parameters | Flexible but structured |
| Low | Exact scripts | Fragile operations, deployments |

### Evaluation-driven development

1. Identify gaps by running Claude on tasks without a skill
2. Build 3+ test scenarios
3. Establish baseline
4. Write minimal instructions
5. Iterate using "Claude A / Claude B" pattern (one writes, one tests)
6. Test across Haiku, Sonnet, and Opus

### Common anti-patterns

- Vague names (`helper`, `utils`)
- Deeply nested file references (keep one level deep)
- Offering too many library choices (provide one default + escape hatch)
- Time-sensitive information that rots
- Inconsistent terminology
- Windows-style backslash paths

---

## 3. Skill Discovery & Installation

### Marketplaces

| Marketplace | Size | Quality | Security | Install method |
|------------|------|---------|----------|---------------|
| **Anthropic Official** | ~20 | Highest | Verified by Anthropic | Ships with Claude Code |
| **Agensi** | 200+ | Curated | 8-point security scan | `curl` one-liner |
| **ClaudeSkills.info** | 658+ | Community | Community-vetted | Manual |
| **Skills.sh** | ~2,000 | Mixed | No audit | `npx skills add` |
| **SkillsMP** | 800K+ | Variable | No audit | Manual |
| **LobeHub** | 169K+ | Scraped | No review | Manual |
| **MCP Market** | ~500 | MCP-focused | Varies | Manual |

### CLI installation (Skills.sh)

```bash
npx skills find <query>     # Search skills
npx skills add <repo-url>   # Install from GitHub
npx skills list              # List installed
```

Works across Claude Code, Codex CLI, Cursor, and OpenClaw. Uses GitHub repos as registry — any public repo with SKILL.md at root is valid.

### Bundled skill-creator

The `skill-creator:skill-creator` skill interviews the user and generates properly structured SKILL.md files. Available in the Anthropic Official Directory.

---

## 4. Skill Distribution & Security

### Plugin packaging

Plugins bundle skills, agents, hooks, MCP servers into one installable unit:
```
.claude-plugin/
  plugin.json          # Manifest
skills/
  my-skill/
    SKILL.md
agents/
  my-agent.md
hooks/
  post-edit.sh
.mcp.json              # MCP server configs
```

### Marketplace system

A marketplace is a `marketplace.json` file in a GitHub repo indexing plugins. Register with `/plugin marketplace add`. Install plugins with `/plugin install`. Refresh with `/plugin marketplace update`.

### Versioning

`version` in plugin.json controls updates. If set, users get updates only on bumps. If omitted + git-backed marketplace, every commit = new version.

### Security landscape (Snyk ToxicSkills, February 2026)

Scan of 3,984 skills found:
- **36.82%** have security flaws of any severity
- **13.4%** contain critical vulnerabilities
- **76** confirmed malicious payloads (credential theft, backdoors, data exfiltration)
- **91%** of malicious skills use prompt injection
- **10.9%** have hardcoded credentials
- **17.7%** fetch untrusted external content

Attack methods: base64-encoded payloads, Unicode smuggling, password-protected ZIP downloads, instructions to disable safety mechanisms.

Two Claude Code CVEs: CVE-2025-54794 (path restriction bypass, CVSS 7.7), CVE-2025-54795 (command injection, CVSS 8.7).

### Remediation

```bash
uvx mcp-scan@latest --skills    # Audit installed skills
```

- Review agent memory files for unauthorized modifications
- Rotate credentials after exposure
- Prefer vetted marketplaces (Anthropic Official, Agensi)

---

## 5. Skill Management

### Precedence hierarchy (highest → lowest)

1. **Enterprise** managed settings
2. **Personal** (`~/.claude/skills/<name>/SKILL.md`)
3. **Project** (`.claude/skills/<name>/SKILL.md`)
4. **Plugin** (`plugin-name:skill-name` namespace)

Same-name skills at higher levels override lower levels. Plugin skills are namespaced and cannot conflict.

### skillOverrides

In `.claude/settings.local.json`:
```json
{
  "skillOverrides": {
    "my-skill": "off"
  }
}
```

Values: `"on"` (full description), `"name-only"` (saves tokens), `"user-invocable-only"` (hidden from Claude, visible in `/` menu), `"off"` (fully hidden). Visual toggle via `/skills` menu.

### Token budget

Descriptions share a budget at 1% of context window. When overflow occurs, least-used skills' descriptions drop first. Diagnostics: `/doctor`. Adjust with `skillListingBudgetFraction` (e.g., `0.02` for 2%). Combined description + when_to_use capped at 1,536 chars (configurable with `maxSkillDescriptionChars`).

### Live change detection

Editing, adding, or removing skills takes effect within the current session without restart. Exception: creating a top-level skills directory that didn't exist at session start requires restart.

### Monorepo support

Skills load from `.claude/skills/` in starting directory and every parent up to repo root. Nested `.claude/skills/` in subdirectories discovered on demand. `--add-dir` directories also auto-load.

---

## 6. Skill Composition

### Seven composable components

CLAUDE.md, Skills (subsume Commands), Subagents, Agent Teams, Plugins, Hooks, MCP Servers.

### Cross-skill invocation

Skills can invoke other skills as part of their workflow. `context: fork` + `agent: Explore` runs in isolated subagent, keeping main conversation clean.

### Subagent integration

Custom agents (`.claude/agents/*.md`) can preload skills via `skills` frontmatter field. Skill content injected at subagent startup rather than waiting for relevance discovery.

### Hook interaction

| Aspect | Hooks | Skills |
|--------|-------|--------|
| Execution | Deterministic (100%) | Probabilistic (Claude decides) |
| Trigger | Lifecycle events | Relevance matching |
| Handler types | Command, Prompt, Agent | Markdown instructions |
| Blocking | PreToolUse can deny (exit 2) | Cannot block |

12 lifecycle events. PreToolUse is the only blocking event. PostToolUse triggers formatting, linting, or skill invocation after edits.

### Skill-scoped hooks

`hooks` frontmatter field lets skills define hooks running only when that specific skill is active.

---

## 7. Skill Optimization

### Quality enforcement via hooks

PostToolUse hooks on `Edit|Write` run prettier, eslint, type checking after every code modification. PreToolUse hooks on `Edit` use prompt-type handlers for semantic security review. Agent-type hooks enable cross-file compliance checking.

### Skill-change-trigger pattern

PostToolUse hook watches SKILL.md changes → auto-runs validation or optimization. The hookify plugin (`/hookify`) automates creating these rules.

### CI/CD integration

Same hooks run in GitHub Actions / GitLab CI. PR events trigger Claude Code with quality gates, posting results as PR comments, blocking merge until gates pass.

### Token efficiency

`allowed-tools` in frontmatter pre-approves tools, reducing permission friction. Combined with hooks: skills define what to do, allowed-tools grants permissions, hooks enforce what must always happen.

### Optimization results

Prompt optimizer skills report 31% token reduction where clear prompts pass through with zero overhead while vague prompts get comprehensive restructuring.

---

## Sources

- [Extend Claude with skills - Claude Code Docs](https://code.claude.com/docs/en/skills)
- [SKILL.md Spec: Every Field - Agensi](https://www.agensi.io/learn/skill-md-format-reference)
- [Inside Claude Code Skills - Mikhail Shilkov](https://mikhail.io/2025/10/claude-code-skills/)
- [Skill authoring best practices - Claude API Docs](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices)
- [Skill Authoring Patterns - Generative Programmer](https://generativeprogrammer.com/p/skill-authoring-patterns-from-anthropics)
- [Claude Code Skills Guide - Nimbalyst](https://nimbalyst.com/blog/claude-code-skills-guide/)
- [Skills-CLI Guide - Medium](https://medium.com/@jacklandrin/skills-cli-guide-using-npx-skills-to-supercharge-your-ai-agents-38ddf3f0a826)
- [Best AI Agent Skills Marketplaces 2026 - Agensi](https://www.agensi.io/learn/best-ai-agent-skills-marketplaces-2026)
- [Anthropic Official Skills - GitHub](https://github.com/anthropics/skills)
- [Plugin Marketplaces - Claude Code Docs](https://code.claude.com/docs/en/plugin-marketplaces)
- [ToxicSkills Study - Snyk](https://snyk.io/blog/toxicskills-malicious-ai-agent-skills-clawhub/)
- [Claude Code Skill Security - Repello AI](https://repello.ai/blog/claude-code-skill-security)
- [Plugins reference - Claude Code Docs](https://code.claude.com/docs/en/plugins-reference)
- [Claude Code Full Stack - alexop.dev](https://alexop.dev/posts/understanding-claude-code-full-stack/)
- [Claude Code Plugins vs Skills - MorphLLM](https://www.morphllm.com/claude-code-plugins-vs-skills)
- [Claude Code Hooks Patterns - Pixelmojo](https://www.pixelmojo.io/blogs/claude-code-hooks-production-quality-ci-cd-patterns)
- [Agent Skills Spec - Claude API Docs](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview)
- [Provision Skills for Organizations - Claude Help Center](https://support.claude.com/en/articles/13119606-provision-and-manage-skills-for-your-organization)

## 2026-06 Delta (added 2026-06-10)

Window 2026-05-25 → 2026-06-10; Claude Code 2.1.152 → 2.1.172. Confidence tags as elsewhere; sources: changelog/npm timestamps, agentskills.io, claude.com/blog, Trail of Bits.

- **Spec stable; dialect drift is the story [HIGH]:** no new agentskills.io frontmatter since the 2026-05-16 clarification. Claude Code diverges from the portable spec: `description` optional with first-paragraph fallback; 1,536-char listing cap vs the spec's hard 1,024; new **`disallowed-tools`** frontmatter (2.1.152) — the first *subtractive* permission field; `context: fork`. Write to the portable core for cross-agent skills (8/10 portability in a 6-agent test).
- **Live reload [HIGH]:** `/reload-skills` command and SessionStart hooks returning `reloadSkills: true` make hook-installed skills available same-session (2.1.152). New `MessageDisplay` hook event; Stop/SubagentStop hooks can return `additionalContext` (2.1.163).
- **Kill-switches (2.1.169) [HIGH]:** `--safe-mode` / `CLAUDE_CODE_SAFE_MODE` starts with ALL customizations (CLAUDE.md/plugins/skills/hooks/MCP) disabled; `disableBundledSkills` hides Anthropic's bundled skills from the model. First-line debugging for "is a skill causing this?"
- **Anthropic's internal skills playbook (claude.com/blog, 2026-06-03) [HIGH]:** hundreds of internal skills in 9 categories (library reference, product verification, data fetching, business process, scaffolding, code quality, CI/CD, runbooks, infra ops); a good skill fits ONE category; **the Gotchas section is the highest-signal content** — grow it on every failure; skip what the model already knows; share via repo check-in → internal marketplace; **log skill usage with a hook to detect under-triggering**.
- **Trigger-eval methodology is now official (agentskills.io docs) [HIGH]:** ~20 labeled queries incl. near-miss negatives, 3 runs each headless, 0.5 trigger-rate threshold, **60/40 train/validation split against description overfitting**, ≤5 iterations, select by validation score. Directive phrasing ("Use this skill when…") hit 100% activation vs 77% passive in a 650-trial study; explicit anti-triggers are the best false-positive cut.
- **The recall ceiling holds [HIGH]:** models miss ~44–56% of conversational prompts lacking domain anchors regardless of description quality; remedies are forced-eval UserPromptSubmit hooks (84–100% activation), explicit system-prompt instructions (+26pp), or moving always-on content to CLAUDE.md/AGENTS.md.
- **Dormant skills cost real tokens [HIGH, N=1 mechanism-corroborated]:** a 7-hour/5-skill measurement attributed 18% of session tokens to skills, 11% to three that never fired; a separate report measured ~26% of a 249K-token session as repeated skill-list injection (issue #54978). Mitigate: `/skills` disable, `disable-model-invocation: true` for manual-only skills.
- **Skill scanners demonstrably bypassed (Trail of Bits, 2026-06-03) [HIGH]:** ClawHub, Cisco's scanner, and all three skills.sh scanners (Gen/Socket/Snyk) defeated in hours (one bypass: 100k newlines → scanner truncation). Static/LLM scanning cannot reliably catch malicious skills — prefer curated marketplaces (anthropics/skills, trailofbits/skills-curated) and org-managed plugins in sensitive contexts.
- **Skills in the API [HIGH]:** Managed Agents beta (`managed-agents-2026-04-01`) attaches ≤20 skills/session; 2026-06-09 added scheduled deployments + vault env-var credentials. Messages API unchanged (≤8 skills via `container`, `/v1/skills` beta).
