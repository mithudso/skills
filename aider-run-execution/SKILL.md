---
name: aider-run-execution
description: >-
  Configures Aider to autonomously emit and execute shell commands (linters, scrapers, tools) and auto-test loops. TRIGGER: "make aider run commands", "auto-lint in aider", "aider autonomous execution", "aider auto-test". SKIP: general aider usage → aider-expert.
version: 1.1.0
updated: 2026-08-18
category: developer
model: claude-sonnet-5
effort: medium
triggers:
  - make aider run commands
  - auto-lint in aider
  - aider autonomous execution
  - aider auto-test
whenToUse:
  - "I want aider to run shell commands automatically"
  - "how do I configure aider auto-linting"
  - "make aider run my test suite automatically"
related_skills:
  - aider-expert
---

# Aider Shell Execution & Auto-Run Strategies

Aider operates with a strict "human-in-the-loop" security philosophy. By default, it **does not** allow the LLM to autonomously execute shell commands in the background. However, you can prompt the LLM to *propose* commands and configure Aider to seamlessly execute them.

## 1. How Aider Reads Shell Commands from the LLM

Aider's default system prompt explicitly tells the LLM:
> *"Concisely suggest any shell commands the user might want to run in \`\`\`bash blocks."*

When the LLM outputs a \`\`\`bash block (e.g., \`\`\`bash\npytest\n\`\`\`), Aider intercepts this block and prompts the user in the terminal:
`Run shell command? (Y)es/(N)o/(D)on't ask again`

## 2. Enabling "Auto-Execute" for the Session

To allow the LLM to run commands without interrupting you every time:
When Aider prompts you with `Run shell command?`, **choose `D` (Don't ask again)**.
This persists your "Yes" preference for the remainder of that session. Any subsequent \`\`\`bash blocks emitted by the LLM will be automatically executed, and the `stdout`/`stderr` output will be fed directly back into the LLM's context window.

## 3. Prompting the LLM to Proactively Emit Commands

To make the LLM use tools like `ast-grep`, web scrapers, or specific linters, you need to add behavioral instructions to Aider (via `.aider.conf.yml`, `CONVENTIONS.md`, or a custom system prompt).

**Example Prompt Injection:**
\`\`\`markdown
You have access to a terminal. When you need to gather information or verify your code, you MUST output the exact command inside a standard markdown bash block. 
- To search the codebase structurally, output a bash block with: `sg -p 'PATTERN' -l LANG`
- To fetch documentation, output a bash block with: `curl -s 'URL' | pandoc -f html -t markdown`
- To check syntax, output a bash block with: `ruff check .`
I have configured my environment to automatically execute your bash blocks and return the stdout/stderr to you. Do not guess; run the command to verify.
\`\`\`

## 4. Alternative: The `--auto-lint` and `--auto-test` Workflow

If your primary goal is to run linters or test suites automatically after every LLM edit, you do not need to prompt the LLM to emit the commands. Aider has built-in execution loops:

- **Auto-Linting:** 
  Run `aider --auto-lint --lint-cmd "ruff check --fix ."` (or `eslint --fix`). Aider will run this command after every file save. If it returns a non-zero exit code, Aider feeds the error back to the LLM to automatically fix it.
- **Auto-Testing:** 
  Run `aider --auto-test --test-cmd "pytest"`. After making changes, Aider runs the test suite. Failing tests trigger an automatic "fix-and-retest" loop by the LLM.

## Summary Checklist for Autonomous Tool Usage in Aider:
1. Provide a `CONVENTIONS.md` instructing the LLM to use \`\`\`bash blocks for specific CLI tools.
2. Trigger the first command and reply with `D` (Don't ask again).
3. The LLM will now act as a semi-autonomous agent, running scripts and analyzing output within that session.
