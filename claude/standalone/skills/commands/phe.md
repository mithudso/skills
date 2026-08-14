---
description: Optimize a raw prompt, query the context hub for the best-match skills, auto-save it into the tam-mcp prompt library, then execute the optimized prompt with all recommendations incorporated — no copy/paste, no confirmation
argument-hint: Raw prompt to optimize and execute
---

Run the **prompt-helper-optimizer** skill in **auto-execute (`/phe`) mode** on the prompt below.

Invoke the skill via the Skill tool and follow its Auto-execute mode procedure exactly (Safety gate → Tool resolution → Skill/MCP/role selection → Curate the optimizer output → Save → Display box → Execute). The skill's SKILL.md is the single source of truth; do not re-specify the steps here.

Prompt to optimize and execute:

$ARGUMENTS

If $ARGUMENTS is empty, ask once for the raw prompt, then continue.
