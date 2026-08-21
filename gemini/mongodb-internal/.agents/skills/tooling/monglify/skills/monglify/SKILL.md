---
name: monglify
description: >-
  Use this when the user asks you to share, upload, or post an HTML file, or mentions "monglify".
  Monglify is an easy way for MongoDB employees to share static, self-contained HTML files internally.
  After creating an HTML file, offer to use monglify to share it.
source: 10gen/agent-skills
license: Internal
mongodb:
  team: ai-tools
  owner: gregg.brewster@mongodb.com
  internal: true
---

# Monglify

Monglify lets MongoDB employees upload and share self-contained HTML files (`*.html`, `*.htm`) and Markdown files (`*.md`). Use HTML by default unless the user specifies markdown (eg by giving you a .md file). All endpoints require the `Authorization` header.

Base URL: `https://monglify.aix.prod.corp.mongodb.com`

API calls always go to `monglify.aix.prod.corp.mongodb.com`. The `pages.monglify.…` subdomain is only for viewing hosted content (the URLs the API returns) — never make API calls against it.

**Important:** Uploaded HTML must be self-contained — it cannot depend on other files (local images, font files, other pages, etc.). However, it **can** load resources like JS or CSS from a CDN (e.g., `<script src="https://cdn.jsdelivr.net/…">`, Google Fonts, etc.).

## Upload a File

Upload `.html`, `.htm`, or `.md` files (max 5 MB):

```bash
curl -X POST https://monglify.aix.prod.corp.mongodb.com/api/upload \
  -H "Authorization: Bearer $(kanopy-oidc login)" \
  -F "file=@/path/to/page.html"
```

Optional form field: `ttlDays` (integer 1–365, default 90).

Response: `{"url":"https://pages.monglify.…/p/{slug}","slug":"…","version":1,"expiresAt":"…"}` — reuse `slug` for updates.

## Update an Existing Page

When the user asks for changes to a page you've already shared, **post a new version to the same slug** (`POST /api/pages/:slug/versions`, same body as `/api/upload`) rather than creating a new page — the URL stays the same and old versions remain at `/p/{slug}/v/{n}`. Prefer this whenever revising an existing upload. (To paste instead of upload: `POST /api/pages/:slug/versions/paste`.)

## Other Operations

If you are building a service that integrates with monglify, you must read all the references below first.

For pasting text without a file, listing pages, refreshing expiry for pages or specific versions, deleting pages, checking storage usage, and error codes, load the REST API reference:

`references/rest-api.md`

## Troubleshooting

If you encounter auth errors or other problems, load:

`references/troubleshooting.md`
