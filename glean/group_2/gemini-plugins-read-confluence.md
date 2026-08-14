# read-confluence

**Category:** Science, Biology & Medicine
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/tooling/read-confluence/skills/read-confluence

## Description
Use when reading MongoDB Confluence pages via MCP tools. Covers fetch by URL, fetch by ID, CQL search, and Glean fallback search patterns.

---

# Reading Corporate Confluence Pages

> Check `references/` for your repo's Confluence space keys and frequently referenced pages before searching.

## When to Use

- Reading MongoDB internal wiki pages
- Searching for documentation across Confluence spaces
- Fetching runbooks, guides, or technical documentation

## MCP Servers

| Server                           | Tools                                                                      | Best For                                   |
| -------------------------------- | -------------------------------------------------------------------------- | ------------------------------------------ |
| Atlassian MCP                    | `confluence_get_page`, `confluence_search`, `confluence_get_page_children` | CQL queries, full Confluence API           |
| Glean MCP (`user-glean_default`) | `search`, `read_document`, `chat`                                          | Keyword search, reading URLs, AI synthesis |

**Note:** The Atlassian MCP server is commonly installed via Docker Desktop MCP Toolkit, but may be installed differently. Check available MCP servers if the expected name isn't found.

Read tool schemas at `/mcps/<server>/tools/<tool>.json` for full parameter details.

---

## URL Parsing

Extract page ID from Confluence URLs:

```
https://wiki.corp.mongodb.com/spaces/<SPACE>/pages/427462539/Page+Title
                                              ↑
                                         Page ID: 427462539
```

Alternative format:

```
https://wiki.corp.mongodb.com/pages/viewpage.action?pageId=427462539
```

---

## Common Confluence Spaces

# see references/ for your space keys

| Space Key | Content                                             |
| --------- | --------------------------------------------------- |
| <SPACE>   | Check your repo's `references/` file for space keys |

---

## Search Syntax

### Atlassian MCP (CQL)

Use `confluence_search` with Confluence Query Language:

- `space = <SPACE> AND text ~ "migration"` - Full-text search in space # see references/ for your space keys
- `title ~ "runbook"` - Search in titles
- `label = "alerts"` - Search by label
- `type = page` - Only pages (not attachments)

### Glean MCP

Use `search` with keyword matching:

- Use SHORT, targeted keywords (not full sentences)
- Filter to Confluence with `app: "confluence"`
- Filter by contributor with `from: "person name"`
- Filter by date with `updated: "past_week"`
- Do NOT use boolean logic (OR/AND) in queries

---

## Tips

- **Glean `read_document`** can batch multiple URLs in one call
- **Glean `chat`** synthesizes across sources - use for complex questions
- **Page not found?** Try searching by title instead of ID
- **Rate limited?** Add delays between requests or use search instead of crawling

---

## Error Handling

If MCP tools are unavailable:

1. Check `.cursor/mcp.json` for configured servers
2. Look for Atlassian/Confluence tools under any available MCP server
3. Fall back to Glean if Atlassian MCP isn't configured
4. Setup instructions: https://wiki.corp.mongodb.com/spaces/MMS/pages/384998696/MCP+Servers

---

## Repo-Specific Context

Before searching, check if a `references/` directory exists alongside this skill.
Look for a file matching your repo name (e.g., `references/mms.md`, `references/server.md`).
Load it to get the correct Confluence space keys and frequently referenced pages for your repo.
Personal overrides in your local `references/` take precedence over central ones.