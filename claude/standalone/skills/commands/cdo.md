---
description: Code Deep Optimizer — multi-pass review-and-fix over a source file or whole repo that applies every Medium+ fix in place and verifies via build/lint/tests
argument-hint: file or repo path [--read-only|--structural-only|--report|--dry-run|--no-verify|--no-sync|--max-files=N|--max-iter=N|--budget-minutes=N]
---

Read ~/.claude/skills/code-deep-optimizer/SKILL.md and execute it against $ARGUMENTS, flags included. The SKILL.md is the single source of truth; do not re-specify its steps here.

If $ARGUMENTS is empty, ask once for the file or repo path, then continue.
