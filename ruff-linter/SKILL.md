---
name: ruff-linter
description: >-
  Uses ruff for blazing fast Python linting and formatting.
---

# Ruff Instructions

When writing or modifying Python code, you MUST use `ruff` to verify your changes and format the code perfectly before concluding your task. Ruff is an extremely fast Python linter and formatter written in Rust.

## Formatting Code
Always run the formatter after editing a Python file:
```bash
ruff format /path/to/file.py
```

## Linting and Auto-Fixing
Run the linter to catch syntax errors, unused imports, or bad practices. Use the `--fix` flag to automatically resolve simple issues.
```bash
ruff check /path/to/file.py --fix
```

If the linter outputs any remaining errors, fix them in the source code before completing your turn.
