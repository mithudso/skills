<!-- hub-reference-banner -->
> **Reference file — part of the `integration-clients` hub.** Formerly the standalone `monday-api` skill.
> Sibling topics in this family are now reference files under the hubs (`integration-clients`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: monday-api
version: 2.1.0
last_updated: 2026-05-29
description: >-
  Monday.com GraphQL API v2 — queries, mutations, column value JSON formats,
  cursor pagination, complexity budgeting, webhooks, file uploads, OAuth 2.0,
  and MCP server configuration.
  TRIGGER: user writes a GraphQL query/mutation against monday.com; formats
  column value JSON; implements pagination with next_items_page; handles
  COMPLEXITY_BUDGET_EXHAUSTED errors; configures monday MCP servers; builds
  Chrome extension → monday.com API bridge.
  SKIP: monday platform architecture, Apps Framework CLI (mapps), Vibe Design
  System, or board views → use monday-dev instead.
category: developer
tags: [monday, graphql, api, mutations, pagination, webhooks, oauth, mcp, chrome-extension]
whenToUse:
  - writing a GraphQL query or mutation against api.monday.com/v2
  - formatting column_values JSON for status, date, people, link, or dropdown columns
  - implementing cursor-based pagination with items_page and next_items_page
  - handling COMPLEXITY_BUDGET_EXHAUSTED errors or optimizing query complexity
  - creating or consuming monday.com webhooks
  - configuring the monday MCP server (@mondaydotcomorg/monday-api-mcp)
  - routing monday API calls through a Chrome MV3 service worker
  - implementing OAuth 2.0 for a monday app
  - uploading files to monday item columns
whenNotToUse:
  - monday platform architecture, board types, or workspace hierarchy → monday-dev
  - monday CLI (mapps), app versions, or code deployment → monday-dev
  - Vibe Design System components → monday-dev
  - monday automations or workflow builder → monday-dev
related_skills:
  - monday-dev
  - monday-board-audit
  - chrome-dev
---

# Monday.com GraphQL API and MCP Server

## When to use this skill

Use when writing GraphQL queries/mutations against the monday.com API, formatting column value JSON, implementing pagination, handling webhooks, managing rate limits/complexity, configuring monday MCP servers, or integrating monday from a Chrome extension or backend service.

For platform architecture, CLI commands, Apps Framework, or AI features, use the companion `monday-dev` skill instead.

---

## 1. Authentication

### 1.1 Token types

| Token Type | Header Format | Lifetime | Use case |
|---|---|---|---|
| Personal API token | `Authorization: <token>` (no Bearer) | Never expires | Scripts, personal tools, Chrome extensions |
| OAuth 2.0 access token | `Authorization: Bearer <token>` | Until user uninstalls app | Apps with delegated user access |
| shortLivedToken (JWT) | In request body from monday | Per-request | Integration/automation recipes |

**Get personal token:** Profile picture > Developers > My access tokens.

**Security:** Personal tokens grant full access scoped to the user. Never log, hardcode, or store unencrypted. In Chrome extensions, encrypt via `chrome.storage.local` with `secret-vault.js` or equivalent.

### 1.2 OAuth 2.0 flow

```
1. Redirect user to:
   https://auth.monday.com/oauth2/authorize?client_id=CLIENT_ID&redirect_uri=REDIRECT_URI&state=STATE

2. User approves permission scopes

3. Monday redirects to your redirect_uri with ?code=AUTH_CODE&state=STATE
   (authorization code valid for 10 minutes)

4. Exchange code for token (server-side POST):
   POST https://auth.monday.com/oauth2/token
   Body: { client_id, client_secret, code, redirect_uri }
   Response: { access_token, token_type: "Bearer" }
```

**Friction-free install (2025+):** Add `force_install_if_needed=true` to the authorize URL to auto-install uninstalled apps before consent.

### 1.3 App request verification

Every server-to-app request includes a JWT in the Authorization header signed with your app's Signing Secret. Verify before trusting the payload:

```js
import jwt from 'jsonwebtoken';
const decoded = jwt.verify(token, process.env.MONDAY_SIGNING_SECRET);
// decoded contains: accountId, userId, backToUrl, etc.
```

---

## 2. GraphQL API Fundamentals

### 2.1 Endpoint and headers

```
POST https://api.monday.com/v2         (queries and mutations)
POST https://api.monday.com/v2/file    (file uploads only)
```

**Required headers:**

| Header | Value | Notes |
|---|---|---|
| `Authorization` | `<token>` or `Bearer <token>` | Personal = no Bearer; OAuth = Bearer |
| `Content-Type` | `application/json` | Omit for file uploads |
| `API-Version` | `2025-07` | Pin to avoid breaking changes |

### 2.2 API versioning

Always pin the `API-Version` header. Versions older than 2025-04 are no longer supported.

| Version | Status | Key changes |
|---|---|---|
| `2025-04` | Supported | Variables must be JSON objects (not strings); column validations enforced |
| `2025-07` | Stable (recommended) | Complexity error format changed; empty board creation; managed columns |
| `2025-10` | Supported | Multi-level boards; doc mutations; `change_item_position` |
| `2026-01` | Supported | `aggregate` object; status/dropdown column creation; Resource Directory APIs |
| `2026-04` | Supported | Multi-level boards default; `create_project`; search query; notetaker API |
| `2026-07` | Current | User entity overhaul (breaking); validation rules; `create_doc_blocks`; async `job_status` |
| `2026-10` | Preview | User migration completes; legacy fields removed |

### 2.3 Basic request pattern

```js
const res = await fetch('https://api.monday.com/v2', {
  method: 'POST',
  headers: {
    Authorization: MY_TOKEN,
    'Content-Type': 'application/json',
    'API-Version': '2025-07'
  },
  body: JSON.stringify({ query, variables })
});
const { data, errors, account_id } = await res.json();

// Monday returns HTTP 200 even on errors — always check errors array
if (errors?.length) {
  const code = errors[0].extensions?.code;
  const msg = errors[0].message;
  throw new Error(`${code}: ${msg}`);
}
```

**Key rules:**
- Variables must be a proper JSON object (not a JSON string — changed in 2025-04)
- All responses include a unique `request_id` (2025-05+)
- Add `Idempotency-Key` header to mutations for safe retries

---

## 3. Board Operations

```graphql
# List boards
query { boards(limit: 10, order_by: created_at) { id name description state board_kind } }

# Board with groups and first page of items
query {
  boards(ids: [<BOARD_ID>]) {
    groups { id title color }
    columns { id title type settings_str }
    items_page(limit: 50) {
      cursor
      items { id name group { id } column_values { id text value type } }
    }
  }
}

# Create board
mutation {
  create_board(board_name: "My Board", board_kind: public, workspace_id: <WS_ID>) { id }
}

# Create item
mutation {
  create_item(
    board_id: <BOARD_ID>,
    group_id: "topics",
    item_name: "New task",
    column_values: "{\"status\": {\"label\": \"Working on it\"}, \"date4\": {\"date\": \"2025-12-31\"}}"
  ) { id }
}

# Update multiple column values
mutation {
  change_multiple_column_values(
    board_id: <BOARD_ID>,
    item_id: <ITEM_ID>,
    column_values: "{\"status\": {\"label\": \"Done\"}, \"numbers\": \"42\"}"
  ) { id }
}

# Move item to another group
mutation { move_item_to_group(item_id: <ITEM_ID>, group_id: "<GROUP_ID>") { id } }

# Archive / delete item
mutation { archive_item(item_id: <ITEM_ID>) { id } }
mutation { delete_item(item_id: <ITEM_ID>) { id } }

# Change item position (2025-10+)
mutation { change_item_position(item_id: <ITEM_ID>, relative_to: <OTHER_ITEM_ID>, position: before_at) { id } }
```

---

## 4. Column Values — Write Format Reference

The `column_values` argument must be a **JSON-stringified string**. Always use `JSON.stringify()`.

| Column Type | API Type | Write JSON | Example |
|---|---|---|---|
| Text | `text` | `"plain string"` | `{"text_col": "Hello"}` |
| Numbers | `numbers` | `"42"` or `42` | `{"numbers": "42"}` |
| Status | `status` | `{"label": "Done"}` or `{"index": 1}` | `{"status": {"index": 1}}` |
| Date | `date` | `{"date": "YYYY-MM-DD"}` | `{"date4": {"date": "2025-12-31"}}` |
| Date + Time | `date` | `{"date": "YYYY-MM-DD", "time": "HH:MM:SS"}` | `{"date4": {"date": "2025-12-31", "time": "14:30:00"}}` |
| Timeline | `timeline` | `{"from": "YYYY-MM-DD", "to": "YYYY-MM-DD"}` | `{"timeline": {"from": "2025-06-01", "to": "2025-06-30"}}` |
| People | `people` | `{"personsAndTeams": [{"id": N, "kind": "person"}]}` | See §4.1 |
| Dropdown | `dropdown` | `{"labels": ["Label 1"]}` | `{"dropdown": {"labels": ["Option A"]}}` |
| Checkbox | `checkbox` | `{"checked": "true"}` | `{"checkbox": {"checked": "true"}}` |
| Email | `email` | `{"email": "a@b.com", "text": "display"}` | |
| Link | `link` | `{"url": "https://...", "text": "display"}` | |
| Phone | `phone` | `{"phone": "+1234567890", "countryShortName": "US"}` | |
| Location | `location` | `{"lat": "40.69", "lng": "-74.04", "address": "NYC"}` | |
| Country | `country` | `{"countryCode": "US", "countryName": "United States"}` | |
| Tags | `tags` | `{"tag_ids": [295026, 295064]}` | Tag IDs, not label text |
| Rating | `rating` | `3` | Integer 1-5 |
| Hour | `hour` | `{"hour": 14, "minute": 30}` | 24-hour format |
| Week | `week` | `{"week": {"startDate": "YYYY-MM-DD", "endDate": "YYYY-MM-DD"}}` | |
| Long Text | `long_text` | `{"text": "content"}` | Markdown supported |
| World Clock | `timezone` | `{"timezone": "America/New_York"}` | IANA timezone |

### 4.1 People column examples

```json
{"personsAndTeams": [{"id": 12345, "kind": "person"}]}
{"personsAndTeams": [{"id": 12345, "kind": "person"}, {"id": 111, "kind": "team"}]}
```

### 4.2 Status column — use index for stability

Labels can be renamed by board admins, breaking `{"label": "Done"}`. Use `{"index": N}` for programmatic access:

```graphql
# Discover indices
query {
  boards(ids: [<BOARD_ID>]) {
    columns { id title settings_str }
  }
}
# settings_str: {"labels": {"0": "Working on it", "1": "Done", "2": "Stuck"}}
```

### 4.3 Reading column values

```graphql
query {
  items(ids: [<ITEM_ID>]) {
    column_values {
      id text value type
      ... on StatusValue { label label_style { color } index }
      ... on DateValue { date time }
      ... on NumbersValue { number }
    }
  }
}
```

### 4.4 Clearing a column value

```json
{"status": null, "date4": {}, "people": {"personsAndTeams": []}}
```

---

## 5. Subitems and Multi-Level Boards

```graphql
# Read subitems
query {
  items(ids: [<ITEM_ID>]) {
    subitems {
      id name
      column_values { id text value }
      parent_item { id name }
      board { id }
    }
  }
}

# Create subitem
mutation {
  create_subitem(parent_item_id: <ITEM_ID>, item_name: "Sub task") {
    id board { id }
  }
}

# Multi-level boards (2025-10+): query children by parent_item_id filter
query {
  items_page_by_column_values(
    board_id: <BOARD_ID>,
    limit: 50,
    columns: [{column_id: "parent_item", column_values: ["<PARENT_ITEM_ID>"]}]
  ) { cursor items { id name parent_item { id name } } }
}

# Rollup columns (multi-level boards)
query {
  items(ids: [<ITEM_ID>]) {
    column_values(capabilities: [CALCULATED]) { id text value }
  }
}
```

---

## 6. Pagination

Max 500 items per page (1,000 for `items_page_by_column_values`). Cursors expire after 60 minutes — no resume on expiry, restart from page 1.

```js
let cursor = null;
const allItems = [];

do {
  const { query, variables } = cursor
    ? {
        query: `query($c: String!) { next_items_page(limit: 500, cursor: $c) { cursor items { id name } } }`,
        variables: { c: cursor }
      }
    : {
        query: `query($bid: ID!) { boards(ids: [$bid]) { items_page(limit: 500) { cursor items { id name } } } }`,
        variables: { bid: BOARD_ID }
      };

  const res = await mondayApi({ query, variables });
  const page = cursor ? res.data.next_items_page : res.data.boards[0].items_page;
  allItems.push(...page.items);
  cursor = page.cursor;  // null on final page
} while (cursor);
```

**Key rules:**
- Always pass cursor as a GraphQL variable (never interpolate into query string)
- `next_items_page` at query root has lower complexity than re-querying `boards.items_page`
- Updates query returns max 100 per page (2025-04+)

### Filtered pagination

```graphql
query {
  boards(ids: [<BOARD_ID>]) {
    items_page(
      limit: 500,
      query_params: {
        rules: [{ column_id: "status", compare_value: ["Done"] }],
        operator: and
      }
    ) { cursor items { id name } }
  }
}
```

---

## 7. Webhooks

### Event types (common)

| Event | Config needed? |
|---|---|
| `create_item` | No |
| `change_column_value` | No |
| `change_specific_column_value` | Yes (`columnId`) |
| `item_moved_to_specific_group` | Yes (`groupId`) |
| `create_update` / `edit_update` / `delete_update` | No |
| `create_subitem` / `change_subitem_column_value` | No |
| `create_column` | No |

### Creating webhooks

```graphql
mutation {
  create_webhook(
    board_id: <BOARD_ID>,
    url: "https://yourapp.com/webhook",
    event: change_specific_column_value,
    config: "{\"columnId\": \"status\"}"
  ) { id board_id }
}

mutation { delete_webhook(id: <WEBHOOK_ID>) { id board_id } }
```

### Challenge verification

```js
app.post('/webhook', (req, res) => {
  if (req.body.challenge) {
    return res.json({ challenge: req.body.challenge });
  }
  const { event } = req.body;
  handleEvent(event);
  res.sendStatus(200);
});
```

### Event payload structure

```json
{
  "event": {
    "type": "change_column_value",
    "boardId": 1234567890,
    "pulseId": 9876543210,
    "itemId": 9876543210,
    "columnId": "status",
    "value": { "label": { "index": 1, "text": "Done" } },
    "previousValue": { "label": null },
    "userId": 12345,
    "triggerTime": "2025-06-15T10:30:00.000Z"
  }
}
```

**Security:** Plain API-created webhooks have no per-request auth — secure with a shared secret in the URL path or validate sender IP. App integration webhooks include a signed JWT — verify with your Signing Secret.

---

## 8. File Uploads

```js
const formData = new FormData();
formData.append('query', `mutation ($file: File!) {
  add_file_to_column(item_id: ${ITEM_ID}, column_id: "${FILE_COLUMN_ID}", file: $file) { id }
}`);
formData.append('variables[file]', fileBlob, 'report.pdf');

await fetch('https://api.monday.com/v2/file', {
  method: 'POST',
  headers: { Authorization: MY_TOKEN },
  // Do NOT set Content-Type — let fetch set the multipart boundary
  body: formData
});
```

**Common mistake:** Sending to `/v2` instead of `/v2/file` returns 400. Setting `Content-Type` manually breaks the multipart boundary.

---

## 9. Rate Limiting and Complexity

| Limit | Value |
|---|---|
| Max complexity per single query | 5,000,000 points |
| Per-minute budget (paid accounts) | 10,000,000 points |
| Per-minute budget (free/trial) | 1,000,000 points |
| Budget scope | Per account (rolling 60 seconds) |

### Check complexity in a query

```graphql
query {
  boards(ids: [<BOARD_ID>]) { items_page(limit: 50) { items { id name } } }
  complexity { before after query reset_in_x_seconds }
}
```

`reset_in_x_seconds` is null during normal operation; populated only in error responses when the budget is exhausted.

### Retry helper

```js
async function mondayApiWithRetry(query, variables, maxRetries = 3) {
  for (let attempt = 0; attempt < maxRetries; attempt++) {
    const res = await fetch('https://api.monday.com/v2', {
      method: 'POST',
      headers: {
        Authorization: TOKEN,
        'Content-Type': 'application/json',
        'API-Version': '2025-07',
        'Idempotency-Key': `${query.slice(0, 20)}-${Date.now()}`
      },
      body: JSON.stringify({ query, variables })
    });

    const json = await res.json();

    if (json.errors?.[0]?.extensions?.code === 'COMPLEXITY_BUDGET_EXHAUSTED') {
      const waitSec = json.errors[0].extensions.reset_in_x_seconds ?? 60;
      await new Promise(r => setTimeout(r, waitSec * 1000));
      continue;
    }

    const retryAfter = res.headers.get('Retry-After');
    if (retryAfter) {
      await new Promise(r => setTimeout(r, parseInt(retryAfter) * 1000));
      continue;
    }

    return json;
  }
  throw new Error('Max retries exceeded');
}
```

### Optimization strategies

1. Request only needed fields — complexity grows with field count and nesting depth
2. Use `next_items_page` at query root instead of re-nesting inside `boards`
3. Batch column updates with `change_multiple_column_values`
4. Cache repeated reads where live data is not critical
5. Use webhooks instead of polling for event-driven updates
6. Avoid deeply nested queries — complexity grows exponentially with depth

---

## 10. MCP Server Configuration

Package: `@mondaydotcomorg/monday-api-mcp` (v3.1.2)

### Hosted MCP (recommended)

```json
{
  "mcpServers": {
    "monday-mcp": { "url": "https://mcp.monday.com/mcp" }
  }
}
```

### Local MCP — default mode

```json
{
  "mcpServers": {
    "monday-api-mcp": {
      "command": "npx",
      "args": ["@mondaydotcomorg/monday-api-mcp@latest"],
      "env": { "MONDAY_TOKEN": "your_token" }
    }
  }
}
```

### CLI flags

| Flag | Description | Default |
|---|---|---|
| `--token` / `-t` | API token | Required (or env `MONDAY_TOKEN`) |
| `--version` / `-v` | API version | `current` |
| `--mode` / `-m` | `default` or `apps` | `default` |
| `--read-only` / `-ro` | Disable write operations | `false` |
| `--enable-dynamic-api-tools` / `-edat` | Full GraphQL access (beta) | `false` |

### Platform MCP tool inventory

**Board and item operations:**
`get_board_info`, `get_board_schema`, `get_board_items_page`, `create_board`, `create_item`, `delete_item`, `change_item_column_values`, `move_item_to_group`, `create_group`, `create_column`, `delete_column`, `board_insights`, `get_board_activity`

**Users and context:** `get_user_context`, `list_users_and_teams`, `search`, `create_notification`

**Workspaces and folders:** `list_workspaces`, `workspace_info`, `create_workspace`, `update_workspace`, `create_folder`, `update_folder`, `move_object`

**Documents:** `create_doc`, `read_docs`, `add_content_to_doc`, `update_doc`, `get_updates`, `create_update`, `get_assets`, `update_assets_on_item`

**Dashboards and views:** `create_dashboard`, `create_view`, `create_widget`, `all_widgets_schema`

**Forms:** `create_form`, `get_form`, `update_form`, `form_questions_editor`, `create_form_submission`

### MCP precondition rules

1. Before reading/writing items: call `get_board_info` to discover column IDs, types, status label indices, and group IDs
2. Before creating widgets: call `all_widgets_schema` for the JSON schema
3. Before creating columns: call `get_column_type_info` for type-specific defaults
4. Before editing docs: call `read_docs` with `include_blocks: true`
5. Before form submissions: call `get_form` for question IDs and validation rules
6. When `get_board_items_page` returns `has_more: true`: continue with `cursor = nextCursor` until exhausted

---

## 11. Chrome Extension Integration

Chrome extensions cannot be embedded as monday App views (iframes). Route all API calls through the MV3 service worker — monday's API does not allowlist Chrome extension origins, so direct fetch from popups, content scripts, or offscreen documents is blocked by CORS.

```js
// service-worker.js
chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
  if (msg.type === 'MONDAY_API') {
    (async () => {
      const token = await getDecryptedToken();
      const res = await fetch('https://api.monday.com/v2', {
        method: 'POST',
        headers: {
          Authorization: token,
          'Content-Type': 'application/json',
          'API-Version': '2025-07'
        },
        body: JSON.stringify({ query: msg.query, variables: msg.variables })
      });
      sendResponse(await res.json());
    })();
    return true;  // keep channel open for async response
  }
});
```

**Token storage:** Personal tokens — encrypt before storing in `chrome.storage.local`. OAuth access tokens — use `chrome.storage.session` (short-lived). See `chrome-storage-patterns` skill for vault patterns.

---

## 12. Error Codes

| Code | Cause | Resolution |
|---|---|---|
| `COMPLEXITY_BUDGET_EXHAUSTED` | Query too heavy or budget depleted | Wait `reset_in_x_seconds`, simplify query |
| `UserUnauthorizedException` | Invalid token or missing OAuth scope | Verify token, check app scopes |
| `ResourceNotFoundException` | Board/item ID does not exist or user lacks access | Verify ID, check permissions |
| `InvalidArgumentException` | Malformed query arguments | Check argument types and values |
| `ColumnValueException` | Invalid column value JSON | Verify JSON format matches column type |
| `CorrectedValueException` | Value was auto-corrected (warning) | Review correction in response |
| `ItemsLimitationException` | Item limit reached (free plan) | Upgrade plan or delete items |

---

## 13. Common Mistakes

| Mistake | Fix |
|---|---|
| `Bearer` prefix on personal token | Personal tokens: no Bearer. OAuth tokens: include Bearer |
| `column_values` sent as JS object | Always `JSON.stringify()` before sending |
| File upload to `/v2` | Use `/v2/file` for multipart file uploads |
| Setting `Content-Type` on file upload | Omit it — let fetch set the multipart boundary |
| Fetching all items without pagination | Use `items_page` + cursor loop |
| Ignoring `API-Version` header | Pin to `2025-07` to avoid breaking changes |
| CORS error from extension context | Route calls through service worker or backend |
| Cursor expired mid-pagination | Restart from page 1; no resume capability |
| Complexity exhausted silently | Add `complexity { before after reset_in_x_seconds }` |
| HTTP 200 treated as success | Always check `errors` array in response body |
| Variables sent as JSON string | Must be proper JSON object (2025-04+) |
| Querying multi-level boards before 2026-04 | Pass `hierarchy_types` argument |
| Rollup column returning empty | Pass `capabilities: [CALCULATED]` |

---

## 14. Quick Reference

```
API endpoint:           POST https://api.monday.com/v2
File upload endpoint:   POST https://api.monday.com/v2/file
OAuth authorize URL:    https://auth.monday.com/oauth2/authorize
OAuth token URL:        https://auth.monday.com/oauth2/token
API version (stable):   2025-07
API version (current):  2026-07
Max items per page:     500 (1,000 for items_page_by_column_values)
Cursor TTL:             60 minutes
Auth code validity:     10 minutes

Rate limits (per account):
  Single call max:      5,000,000 complexity points
  Per-minute (paid):    10,000,000 complexity points
  Per-minute (free):    1,000,000 complexity points

MCP package:            @mondaydotcomorg/monday-api-mcp (v3.1.2)
MCP hosted URL:         https://mcp.monday.com/mcp
```

---

## 15. Sources

1. [Authentication](https://developer.monday.com/api-reference/docs/authentication)
2. [API Versioning](https://developer.monday.com/api-reference/docs/api-versioning)
3. [GraphQL overview](https://developer.monday.com/api-reference/docs/introduction-to-graphql)
4. [Rate limits](https://developer.monday.com/api-reference/docs/rate-limits)
5. [Column Types Reference](https://developer.monday.com/api-reference/reference/column-types-reference)
6. [Items page](https://developer.monday.com/api-reference/reference/items-page)
7. [Webhooks Reference](https://developer.monday.com/api-reference/reference/webhooks)
8. [Subitems](https://developer.monday.com/api-reference/reference/subitems)
9. [Assets/Files](https://developer.monday.com/api-reference/reference/assets-1)
10. [OAuth and Permissions](https://developer.monday.com/apps/docs/oauth)
11. [monday.com MCP GitHub](https://github.com/mondaycom/mcp)
12. [Build on monday with AI](https://developer.monday.com/api-reference/docs/build-on-monday-with-ai)
