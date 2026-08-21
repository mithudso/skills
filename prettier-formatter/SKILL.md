---
name: prettier-formatter
description: >-
  Uses prettier to format web code (HTML, CSS, JS, TS, JSON).
---

# Prettier Instructions

When writing or modifying web frontend files (HTML, CSS, JavaScript, TypeScript, or JSON), you MUST use `prettier` to format the code perfectly before concluding your task.

## Formatting Code In-Place
Always run prettier with the `--write` flag to format files in place:
```bash
# Format a specific file
prettier --write /path/to/file.js

# Format an entire directory of web files
prettier --write "/path/to/src/**/*.{js,jsx,ts,tsx,css,html}"
```

Do not present unformatted web code to the user.
