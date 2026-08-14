<!-- hub-reference-banner -->
> **Reference file — part of the `misc-catch-all` hub.** Formerly the standalone `skill-lookup` skill.
> Sibling topics in this family are now reference files under the hubs (`misc-catch-all`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: skill-lookup
description: >-
  Search, retrieve, and install Agent Skills from the prompts.chat registry using MCP tools. Use when the user asks to find skills, browse skill catalogs, install a skill for Claude, or extend Claude's capabilities with reusable AI agent components. TRIGGER: "find a skill for X", "install a Claude skill", "browse skill catalog", "what skills exist for Y", "add a skill to my Claude". SKIP: creating a new skill from scratch (use skill-creator); optimizing an existing skill (use skill-optimizer).
origin: local
version: 1.1.0
category: developer
updated: 2026-05-29
whenToUse:
  - "find a skill for code review"
  - "install a Claude skill"
  - "browse the skill catalog"
  - "what skills are available for X"
  - "add a skill to my Claude setup"
  - "search for agent skills"
  - "extend Claude's capabilities with a skill"
related_skills:
  - skill-creator
  - skill-optimizer
  - claude-code-skills
---

# Skill Lookup

Search for and install Agent Skills from the prompts.chat registry.

**Prerequisite:** The `search_skills` and `get_skill` MCP tools must be available in the current session (provided by the prompts.chat MCP server). If they are not listed in available tools, tell the user: "The prompts.chat MCP server is not connected — connect it in your Claude settings to use this skill."

## When NOT to use

- Creating a new skill from scratch → use `skill-creator`
- Optimizing or auditing an existing local skill → use `skill-optimizer`
- The user already has the installed skill path and just needs to use it

## Workflow

1. Search for skills matching the user's request using `search_skills`
2. Present results and ask the user to pick one (if multiple match)
3. Retrieve the chosen skill with `get_skill` to get all files
4. Install by saving files to `.claude/skills/{slug}/` and verify `SKILL.md` exists
5. Confirm installation and explain what the skill does and when it activates

## Available tools

| Tool | Purpose |
|---|---|
| `search_skills` | Search by keyword, category, or tag |
| `get_skill` | Retrieve a specific skill by ID (returns all files) |

## How to search

Call `search_skills` with:

- `query`: Keywords from the user's request
- `limit`: Number of results (default 10, max 50)
- `category`: Filter by category slug (e.g., `"coding"`, `"automation"`)
- `tag`: Filter by tag slug

Present results showing: title, description, author name, file list, category and tags, and link to the skill.

**If results are returned:** list them and ask "Which of these would you like to install?" before proceeding.

**If no results are returned:** broaden the query (remove category filter, use synonyms), then report "No skills found for [query]. Try describing what you want the skill to do."

## How to get a skill

Call `get_skill` with:

- `id`: The skill ID from search results

Returns skill metadata and all file contents: `SKILL.md`, reference docs, helper scripts, and config files. The response includes a `slug` field (e.g., `"pr-review-toolkit"`) — use this as the directory name in the install step.

**If `get_skill` fails or returns an error:** report the error to the user and offer to search again or try a different skill ID.

## How to install

1. Call `get_skill` to retrieve all files
2. Use the `slug` from the response to create `.claude/skills/{slug}/`
3. Save each file to the correct location:
   - `SKILL.md` → `.claude/skills/{slug}/SKILL.md`
   - Other files → `.claude/skills/{slug}/{filename}`
4. Read back `SKILL.md` and confirm the frontmatter `name` field is present

## Example

```
search_skills({"query": "code review", "limit": 5, "category": "coding"})
# → presents results, user picks "pr-review-toolkit"
get_skill({"id": "abc123"})
# → saves to .claude/skills/pr-review-toolkit/SKILL.md
```

## Guidelines

- Always search before suggesting the user create their own skill
- Present search results in a readable format with file counts
- Ask the user to confirm which skill to install when multiple match
- Confirm successful save after installation
- Explain what the skill does and what phrases trigger it
