# Monglify REST API Reference

All `/api/*` endpoints require the `Authorization` header. The Kanopy mesh injects this automatically for browser requests; include it manually for direct API calls.

Base URL: `https://monglify.aix.prod.corp.mongodb.com`

Authentication:
```bash
Authorization: Bearer $(kanopy-oidc login)
```

## Upload a New File

**POST** `/api/upload`

Body: `multipart/form-data`

| Field | Type | Required | Description |
|---|---|---|---|
| `file` | File | yes | `.html`, `.htm`, or `.md`, max 5 MB |
| `ttlDays` | integer | no | Days until expiry, 1–365 (default: 90) |

Response:
```json
{
  "url": "https://pages.…/p/tiger-mango-ocean",
  "slug": "tiger-mango-ocean",
  "version": 1,
  "expiresAt": "2026-08-17T…"
}
```

You can also upload directly from a heredoc without creating a file first:

```bash
curl -X POST https://monglify.aix.prod.corp.mongodb.com/api/upload \
  -H "Authorization: Bearer $(kanopy-oidc login)" \
  -F "file=@-;filename=page.html;type=text/html" <<'EOF'
<!DOCTYPE html>
<html>
<body>hello</body>
</html>
EOF
```

## Paste Text

**POST** `/api/paste`

Paste HTML or Markdown text without a file. If `format` is omitted, the server auto-detects HTML (looks for `<!DOCTYPE` or `<html>` at the start) and falls back to Markdown.

Body: JSON

| Field | Type | Required | Description |
|---|---|---|---|
| `text` | string | yes | The content to publish |
| `format` | string | no | `"markdown"` or `"html"`, auto-detected if omitted |
| `title` | string | no | Page title, defaults to `paste.md` or `paste.html` |
| `ttlDays` | integer | no | Days until expiry, 1–365 (default: 90) |

Response: Same shape as upload.

## Add a New Version

**POST** `/api/pages/:slug/versions`

Body: same as `/api/upload`. Returns the same response shape. The URL stays the same; old versions remain accessible at `/p/{slug}/v/{n}`.

## Paste a New Version

**POST** `/api/pages/:slug/versions/paste`

Body: same as `/api/paste`. Adds a new version by pasting text.

## Refresh Page Expiry

**POST** `/api/pages/:slug/refresh`

Reset expiry on the latest version of a page. Body: optional JSON

```json
{ "ttlDays": 30 }
```

Response: `{"slug": "…", "version": 3, "expiresAt": "…"}`

## Refresh Version Expiry

**POST** `/api/pages/:slug/versions/:version/refresh`

Reset expiry on a specific version. Body: optional JSON (same as page refresh).

Response: `{"slug": "…", "version": 2, "expiresAt": "…"}`

## List Pages

**GET** `/api/pages`

Returns all pages owned by the authenticated user. No parameters.

## Get Page Metadata

**GET** `/api/pages/:slug`

Returns metadata and version history for a page. HTML content is excluded.

## Delete a Page

**DELETE** `/api/pages/:slug`

Permanently deletes a page and all its versions.

Response: `{"slug": "…", "deleted": true}`

## Storage Usage

**GET** `/api/usage`

Returns storage usage for the authenticated user.

Response:
```json
{ "usedBytes": 102400, "limitBytes": 52428800 }
```

## Error Responses

All errors return JSON: `{ "error": "Human-readable message" }`

| Status | Meaning |
|---|---|
| 400 | Validation error (missing file, wrong extension, invalid slug or ttlDays) |
| 401 | Missing or invalid authentication |
| 403 | Authenticated but not the page owner |
| 404 | Page not found |
| 409 | Version conflict (retry the request) |
| 413 | File too large or storage quota exceeded |
| 500 | Internal server error |

## Limits

- Max file size: 5 MB per upload
- Storage quota: 50 MB per user
- Expiry: 90 days by default, configurable 1–365 days
