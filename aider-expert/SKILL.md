---
name: aider-expert
description: Deep integration knowledge and optimization strategies for the Aider CLI agent, bridging Antigravity's tools and MCP servers into Aider's workflow.
---

# Aider Expert & Optimization Guide

This skill documents the advanced setup constructed to bring Antigravity-level power into Aider. Trigger this skill whenever the user asks about Aider, its configuration, performance optimization, or bridging external tools into its workflow.

## 1. Global Configuration (`~/.aider.conf.yml`)
Aider is configured to load optimal settings by default:
- **`edit-format: whole`**: Used to ensure highly accurate replacements without diff-matching errors.
- **`auto-commits: false`**: Mimics Antigravity's `caveman-commit` style. Commits should be manual or requested via `/commit` to avoid history pollution.
- **`auto-lint: true`**: Automatically runs `ruff` (Python) and `prettier` (Web) after edits to guarantee syntax formatting without wasting model tokens.
- **Conventions Injection**: Aider automatically reads `~/.aider.conventions.md` on boot.

## 2. Antigravity Conventions (`~/.aider.conventions.md`)
We compressed 240+ Antigravity skills into a dense set of rules for Aider. These include:
- **Caveman Mode:** High signal-to-noise ratio in communication.
- **Functional Core, Imperative Shell:** Architecture patterns.
- **Antigravity Design Protocol:** UI aesthetics (HSL colors, clean typography, micro-animations, no cliches).
- **Security Protocols:** Strict environment variable and sanitization rules.

## 3. Tooling & MCP Bridging
Aider does NOT support native MCP clients. Instead, Aider interacts with Antigravity tools using its **`/run`** command.

### Structural Refactoring (`ast-grep`)
Aider can use `ast-grep` (sg) for deep AST searches instead of regex:
```text
/run sg -p 'function $A() { $$$B }' -l javascript
```

### JSON/YAML Slicing (`yq` / `jq`)
When reading massive config files, tell Aider to run:
```text
/run yq eval '.dependencies' package.json
```

### Future MCP Interaction
If you need Aider to fetch external data (like Firecrawl scraping), do NOT ask Aider to use MCP natively. Instead, have Aider run a Python or Node wrapper script that communicates with the MCP server, and output the result to stdout for Aider to ingest.

## 4. Hardware Integrations (eGPU)
Aider has a dedicated alias `aider-egpu` that runs it against the local `tinygrad` server connected to the NVIDIA 5080. When using this, API requests bypass Anthropic/OpenAI and route to `localhost:8000`.

## 5. Typical Workflows
- **Architecting:** Aider has a `/architect` mode. Instruct the user to use it for multi-file changes to prevent premature code generation.
- **Linting:** Instruct the user to rely on the auto-linters. If code needs formatting, just save it; Aider's hooks will run `ruff` or `prettier`.
