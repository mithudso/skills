---
name: local-semantic-search
description: "Use this skill to query the local codebase semantically instead of reading full files, drastically reducing token usage."
---

# Local Semantic Search

You have access to semantic search over the local file corpus via the
**global_ai_hub** MCP server (Ollama embeddings; SQLite-backed file index in
`~/.global-ai-hub/hub.db`, ChromaDB-backed docset indexes).

## When to use this
- Do NOT use `view_file` to blindly explore large files.
- ALWAYS use `hub_search_codebase` first when trying to locate specific routing logic, architecture decisions, or code snippets.
- For mirrored documentation sets (web-text-mirror output), use `hub_query_docset` — see `hub_list_docsets`; index a new mirror with `hub_index_docset`.
- Use the semantic index to keep the context window small. Only read full files if you need to surgically edit them.

## Setup & Tooling
The MCP server `~/.global-ai-hub/mcp-server/hub_mcp_server.py` exposes
`hub_search_codebase`, `hub_query_docset`, `hub_list_docsets`,
`hub_index_docset`, `hub_distill_run`, `hub_memory_search`, `hub_memory_stats`
(wiring: `~/.global-ai-hub/.mcp.json`; docs: `~/.global-ai-hub/docs/MCP.md`).
If the MCP server is not connected, fall back to the CLI:
- file index: `python3 ~/.global-ai-hub/scripts/search.py "your query"`
- docsets: `~/.global-ai-hub/.venv/bin/python ~/.global-ai-hub/scripts/docset_indexer.py query <docset> "your query"`
