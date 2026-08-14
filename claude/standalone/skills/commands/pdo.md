---
description: Run the prompt-deep-optimizer multi-pass audit + convergence loop on a production prompt (system prompt, agent instructions, tool template)
argument-hint: Prompt text or file path to optimize
---

Run the **prompt-deep-optimizer** skill on the prompt below.

Invoke the skill via the Skill tool and follow its SKILL.md procedure exactly (Step 1 ingest + safety/fragment gates → Step 2 multi-pass audit → Step 3 triage → Step 4 fixes → Step 5 convergence loop → Step 6 verify + final output with algorithm recommendation). The skill's SKILL.md is the single source of truth; do not re-specify the steps here. One-off or exploratory prompts under ~600 tokens are better served by /ph or /phe — the skill's own handoff rule applies.

Prompt to optimize (text or file path):

$ARGUMENTS

If $ARGUMENTS is empty, ask once: "Paste the prompt to optimize, or give me a file path."
