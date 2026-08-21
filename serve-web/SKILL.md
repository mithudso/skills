---
name: serve-web
description: >-
  Uses serve to quickly spin up a local static HTTP server for web projects.
---

# Serve Instructions

When the user asks you to build a static web prototype (HTML/JS/CSS) and preview it, or asks you to serve a directory, you MUST use the `serve` CLI tool instead of writing custom Python `http.server` scripts.

## Usage
Navigate to the directory containing the `index.html` file and run:
```bash
serve .
```

To run it in the background so you can continue working, use the `run_command` tool with `IsDaemon: true` or manually send it to the background:
```bash
serve . &
```

Provide the user with the localhost URL (usually `http://localhost:3000`) so they can view the site.
