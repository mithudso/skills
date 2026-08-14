---
description: Optimize a prompt, recommend relevant skills and MCPs, critique weaknesses, and auto-save reusable prompts
argument-hint: Raw prompt to optimize
---

Run the **prompt-helper-optimizer** skill in **review (`/ph`) mode** on the prompt below.

Invoke the skill via the Skill tool and follow its Review mode procedure exactly (Safety gate → Tool resolution → Curate the optimizer output → Skill/MCP/role selection → return the six review sections). Review mode does not execute the prompt. The skill's SKILL.md is the single source of truth; do not re-specify the steps here.

Prompt to optimize:

$ARGUMENTS

If $ARGUMENTS is empty, ask once for the raw prompt, then continue.
