---
name: ast-grep
description: >-
  Uses ast-grep (sg) for structural code search and refactoring based on AST patterns rather than regex.
---

# ast-grep Instructions

When the user asks you to find complex patterns in code, safely refactor functions, or perform structural searches that span multiple lines, you MUST use `ast-grep` (the `sg` command) instead of `grep_search` or regular `grep`.

## Usage
`ast-grep` parses code into an Abstract Syntax Tree (AST), meaning it ignores formatting, whitespace, and comments, and matches the actual logic.

### Structural Search
To search for a specific structure, use `-p` (pattern):
```bash
# Find all try-catch blocks
sg -p 'try { $$$STATEMENTS } catch ($E) { $$$CATCH }' -l js

# Find all calls to a specific function with any arguments
sg -p 'myFunction($$$ARGS)'
```

### Structural Replace
To refactor code, use `-p` (pattern) and `-r` (replacement):
```bash
# Replace console.log with a custom logger, preserving the arguments
sg -p 'console.log($$$ARGS)' -r 'logger.info($$$ARGS)' -i
```

## Metavariables
- `$VAR`: Matches a single AST node (e.g., an identifier, literal, or expression).
- `$$$VARS`: Matches multiple AST nodes (e.g., arguments in a function call, or multiple statements in a block).

Always run `sg --help` if you need to refresh your memory on the flags.
