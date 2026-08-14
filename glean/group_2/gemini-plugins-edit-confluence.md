# edit-confluence

**Category:** Science, Biology & Medicine
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/tooling/edit-confluence/skills/edit-confluence

## Description
Use when editing Confluence wiki pages via MCP tools. Covers safe editing workflow: fetch first, edit in storage format, handle CDATA blocks.

---

# Editing Confluence Pages

> Check `references/` for your repo's Confluence space keys and frequently edited page patterns before editing.

## Purpose

Provide a safe workflow for editing Confluence wiki pages via MCP tools. The update API replaces the full page body on every call, so edits must fetch the current content first and merge changes into it — otherwise unrelated content is destroyed.

## When to Use

- Adding, updating, or removing sections on a Confluence page
- Appending information (meeting notes, status updates, entries)
- Modifying content at a specific location in a page

## Critical Safety Rules

**NEVER use `content_format: "markdown"` or `"wiki"` for edits.** Both silently convert input to HTML and **replace the entire page body**, destroying all existing text and macro wrappers.

**ALWAYS use `content_format: "storage"`.** This is the only format that writes content exactly as provided, preserving page structure.

**ALWAYS fetch the page before editing.** The update API requires the complete page content — there is no partial/patch update.

These tools are provided by the Atlassian MCP server. See the `read-confluence` skill for MCP server details and fallback options.

---

## URL Parsing

Extract the page ID from Confluence URLs:

```
https://wiki.corp.mongodb.com/spaces/<SPACE>/pages/427462539/Page+Title
                                              ↑ Page ID
# see references/ for your space keys

https://wiki.corp.mongodb.com/pages/viewpage.action?pageId=427462539
                                                          ↑ Page ID
```

---

## Editing Workflow

### 1. Fetch Current Content

Fetch the page twice — raw storage for editing, markdown for readability:

```
confluence_get_page(page_id: "<id>", include_metadata: true, convert_to_markdown: false)
confluence_get_page(page_id: "<id>", include_metadata: false, convert_to_markdown: true)
```

From the raw response, save:

- **`title`** — required for the update call, must match exactly
- **`content.value`** — the raw storage format content to edit

### 2. Identify the Page Format

Inspect the raw `content.value` to determine the page type:

**Markdown macro page** (most MongoDB wiki pages) — contains:

```xml
<ac:structured-macro ac:name="markdown" ...>
  <ac:plain-text-body><![CDATA[ ...markdown... ]]></ac:plain-text-body>
</ac:structured-macro>
```

**Raw HTML page** — contains tags like `<h1>`, `<p>`, `<table>` directly.

### 3. Make Edits

**For markdown macro pages:**

1. Extract the markdown text from inside the `<![CDATA[ ... ]]>` block
2. Edit the markdown (append sections, modify text, add table rows)
3. Place the edited markdown back inside the CDATA block
4. Preserve the full `<ac:structured-macro>` wrapper and its `ac:macro-id` exactly
5. Ensure the edited markdown does not contain the literal sequence `]]>` — this terminates the CDATA block and corrupts the page XML. If needed, rephrase content to avoid it.

**For raw HTML pages:**

1. Edit the HTML directly — add new elements at the desired position

### 4. Submit the Update

```
confluence_update_page(
  page_id: "<id>",
  title: "<exact current title>",
  content: "<full page content with edits merged in>",
  content_format: "storage",
  is_minor_edit: true,
  version_comment: "Added <description of change>"
)
```

- `content` must be the **complete** page content, not just the diff
- Set `is_minor_edit: true` for small additions to reduce notification noise
- Always include a descriptive `version_comment`

### 5. Verify

Fetch the page again with `convert_to_markdown: true` to confirm the edit rendered correctly.

---

## Non-Destructive Alternatives

For changes that do not require modifying page content:

| Tool                                       | Use Case                                                 |
| ------------------------------------------ | -------------------------------------------------------- |
| `confluence_add_comment(page_id, content)` | Add discussion or feedback without touching page content |
| `confluence_add_label(page_id, name)`      | Tag pages for organization                               |

---

## Verified Content Format Behaviors

| `content_format` | Behavior                                              | Safe for edits?       |
| ---------------- | ----------------------------------------------------- | --------------------- |
| `"storage"`      | Writes raw Confluence storage XML exactly as provided | Yes                   |
| `"markdown"`     | Converts to HTML, **replaces entire page**            | No — destroys content |
| `"wiki"`         | Converts to HTML, **replaces entire page**            | No — destroys content |

---

## Error Recovery

| Error            | Cause                                | Resolution                                                        |
| ---------------- | ------------------------------------ | ----------------------------------------------------------------- |
| Version conflict | Page edited between fetch and update | Re-fetch the page and retry with merged content                   |
| Title mismatch   | Title in update differs from current | Use exact title from the fetch response                           |
| Mangled content  | Used `markdown` or `wiki` format     | Restore from Confluence page history, re-do with `storage` format |

---

## Repo-Specific Context

Before editing, check if a `references/` directory exists alongside this skill.
Look for a file matching your repo name (e.g., `references/mms.md`, `references/server.md`).
Load it to get the correct Confluence space keys and common page URL patterns for your repo.
Personal overrides in your local `references/` take precedence over central ones.