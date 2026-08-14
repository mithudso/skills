# source-command-sync-skills

**Category:** General & Specialized Utilities
**Platform:** Claude Code
**Original Path:** claude-code/source-command-sync-skills

## Description
Sync locally installed skills to the mdb-context-hub registry. Run after creating or updating skills.

---

# source-command-sync-skills

Use this skill when the user asks to run the migrated source command `sync-skills`.

## Command Template

You are executing the **context-hub-skill-sync** workflow.

## Instructions

1. List all locally installed skills: `ls ~/.Codex/skills/`
2. Call `tam_list_skills` (limit 100) to get the current hub registry.
3. Compare the two sets. For each skill installed locally but missing from the hub:
   a. Read `~/.Codex/skills/<id>/SKILL.md` for title and description.
   b. Read the largest `.md` file in `references/` (if it exists) for context markdown.
   c. Call `tam_create_skill` with the extracted data.
4. For each skill that exists in both, check if the local SKILL.md description differs from the hub entry. If so, call `tam_update_skill`.
5. Report a summary table of actions taken.

## Category inference
- Skills about MongoDB, Atlas, or database topics → `mongodb`
- Skills about development tools, languages, frameworks → `developer`  
- Skills about TAM workflows, prompts, or custom tooling → `custom`

## Constraints
- Never delete skills from the hub.
- Truncate context markdown to 50,000 characters.
- Skip skills whose SKILL.md can't be read.