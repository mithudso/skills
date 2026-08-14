# Agent Skills Repository

A unified, comprehensive collection of **skills**, **workflows**, **prompts**, and **domain intelligence** for modern AI agent environments: **Google Antigravity**, **Gemini**, **Claude Code**, and **Claude**.

> **Repository Status**: 357 Production Skills Cataloged & Indexed

## Repository Overview

```
skills/
├── antigravity/               # Google Antigravity core & built-in skills (AGY, SDK, rules)
├── gemini/                    # Gemini ecosystem & official plugin skills (Science, Firebase, BigQuery, DevTools, etc.)
├── claude-code/               # Claude Code & Agents skills library (Engineering, DevOps, RAG, Firecrawl, etc.)
├── claude/                    # Claude standalone workspace skills & specialized architectural trees
├── custom/                    # Dedicated domain skills (MongoDB Atlas, MQL Optimizer, Account Intel)
├── catalog.json               # Machine-readable JSON metadata catalog for programmatic skill discovery
├── INDEX.md                   # Full categorized and alphabetical searchable index
└── README.md                  # Main entry point and usage guide
```

## Distribution Summary

### By Platform
| Platform | Total Skills | Primary Focus |
| :--- | :---: | :--- |
| **Claude Code** | `156` | Deep engineering, infrastructure, AI orchestration, Firecrawl automation, psychology & analysis |
| **Gemini / Plugins** | `146` | Domain plugins: Science/Bio, Data Agent Kit (BigQuery/dbt), Firebase, Chrome DevTools, Modern Web |
| **Claude** | `33` | Standalone Claude workspace intelligence, security access trees, and prompt optimization |
| **Claude (Archived)** | `9` | Specialized RFID, NFC, physical security, and trading strategy reference trees |
| **Custom / MongoDB** | `7` | MongoDB Atlas connection, query optimization, stream processing, search & AI |
| **Antigravity** | `3` | Antigravity CLI, IDE customizations, SDK orchestration, and permissioning |
| **Custom / MQL Optimizer** | `2` | Deep MQL query optimization, execution plan analysis, and indexing |
| **Custom / Customer Intelligence** | `1` | Enterprise account intelligence and incident response playbooks |

### By Domain Category
| Category | Skills Count |
| :--- | :---: |
| **Science, Biology & Medicine** | `120` |
| **Databases, Data Engineering & Analytics** | `78` |
| **Cloud, DevOps & Infrastructure** | `46` |
| **Frontend & Web Development** | `45` |
| **Security, Auth & Diagnostics** | `24` |
| **AI, Agents & Prompt Engineering** | `23` |
| **General & Specialized Utilities** | `8` |
| **Programming & Languages** | `5` |
| **Web Scraping & Firecrawl** | `4` |
| **Writing, Communication & Marketing** | `2` |
| **Finance, Banking & Business** | `1` |
| **Psychology & Behavioral Science** | `1` |

---

## Quick Navigation

- 📖 **[Full Searchable Master Index (INDEX.md)](INDEX.md)**
- ⚡ **[Antigravity Skills Directory (antigravity/)](antigravity/README.md)**
- 💎 **[Gemini & Plugins Directory (gemini/)](gemini/README.md)**
- 🤖 **[Claude Code Skills Directory (claude-code/)](claude-code/README.md)**
- 🧠 **[Claude Workspace Skills Directory (claude/)](claude/README.md)**
- 🛠️ **[Custom & Specialized Skills (custom/)](custom/README.md)**

---

## What is an Agent Skill?

An **Agent Skill** is a self-contained capability package that equips an AI agent with specialized operational procedures, domain knowledge, and execution workflows. Each skill directory contains:

- `SKILL.md`: The core definition file with YAML frontmatter (`name`, `description`), activation triggers, and step-by-step procedures.
- `references/` *(optional)*: In-depth technical specifications, cheat sheets, and domain reference guides.
- `scripts/` *(optional)*: Automated scripts and CLI helpers executed by the agent.
- `examples/` *(optional)*: Input/output examples and golden reference templates.

## How to Use These Skills

### 1. In Antigravity / Gemini
Place skill directories in:
- System-wide: `~/.gemini/antigravity/builtin/skills/<skill-name>/` or `~/.gemini/config/plugins/<plugin-name>/skills/<skill-name>/`
- Project-specific: `.gemini/skills/<skill-name>/` or `.antigravity/skills/<skill-name>/`

### 2. In Claude Code / Anthropic Agent System
Place skill directories in:
- System-wide: `~/.agents/skills/<skill-name>/` or `~/.claude/skills/<skill-name>/`
- Project-specific: `.claude/skills/<skill-name>/` or `.agents/skills/<skill-name>/`

### 3. Programmatic Access
Use `catalog.json` to programmatically search, filter, and load skills into any custom LLM workflow or MCP server.

---

## License & Attribution

Maintained by Mitchell Hudson. Compiled and structured for unified multi-agent development.
