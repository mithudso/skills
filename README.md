# Agent Skills Repository

A unified, comprehensive collection of **skills**, **workflows**, **prompts**, and **domain intelligence** for modern AI agent environments: **Google Antigravity**, **Gemini**, **Claude Code**, and **Claude**.

> **Repository Status**: 443 Production Skills Cataloged & Indexed (Auto-Synchronized)

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
| **Gemini / Plugins** | `174` | Domain plugins: Science/Bio, Data Agent Kit (BigQuery/dbt), Firebase, Chrome DevTools, Modern Web |
| **Claude Code** | `156` | Deep engineering, infrastructure, AI orchestration, Firecrawl automation, psychology & analysis |
| **Claude** | `91` | Standalone Claude workspace intelligence, security access trees, and prompt optimization |
| **Claude (Archived)** | `9` | Specialized RFID, NFC, physical security, and trading strategy reference trees |
| **Custom / MongoDB** | `7` | MongoDB Atlas connection, query optimization, stream processing, search & AI |
| **Antigravity** | `3` | Antigravity CLI, IDE customizations, SDK orchestration, and permissioning |
| **Custom / MQL Optimizer** | `2` | Deep MQL query optimization, execution plan analysis, and indexing |
| **Custom / Customer Intelligence** | `1` | Enterprise account intelligence and incident response playbooks |

### By Domain Category
| Category | Skills Count |
| :--- | :---: |
| **Science, Biology & Medicine** | `152` |
| **Databases, Data Engineering & Analytics** | `95` |
| **Frontend & Web Development** | `57` |
| **Cloud, DevOps & Infrastructure** | `51` |
| **AI, Agents & Prompt Engineering** | `30` |
| **Security, Auth & Diagnostics** | `26` |
| **General & Specialized Utilities** | `17` |
| **Programming & Languages** | `6` |
| **Web Scraping & Firecrawl** | `4` |
| **Writing, Communication & Marketing** | `2` |
| **Finance, Banking & Business** | `2` |
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

## Automated Synchronization Hook

This repository features real-time hooks and background daemons. Any skill created or edited in **Google Antigravity**, **Gemini**, **Claude Code**, or local workspaces is automatically consolidated, cataloged, and pushed to GitHub.

---

## License & Attribution

Maintained by Mitchell Hudson. Compiled and structured for unified multi-agent development.
