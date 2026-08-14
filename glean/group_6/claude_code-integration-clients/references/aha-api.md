<!-- hub-reference-banner -->
> **Reference file — part of the `integration-clients` hub.** Formerly the standalone `aha-api` skill.
> Sibling topics in this family are now reference files under the hubs (`integration-clients`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: aha-api
version: 1.1.0
updated: 2026-05-29
category: developer
tags: [aha, product-management, rest-api, oauth2, webhooks, chrome-extension, ember-spa]
related_skills: [jira-extension-client, salesforce-scraping-patterns, chrome-mv3-advanced, dom-scraping-resilience]
description: >
  Use when integrating with Aha! product management platform via REST API, webhooks, or Chrome
  extension DOM scraping — features, initiatives, goals, ideas portal, shared pages, custom fields,
  OAuth2/API key auth, or Ember.js SPA navigation patterns specific to Aha!.
  SKIP: general Chrome extension architecture (use chrome-mv3-advanced), generic OAuth2 flows
  not targeting Aha!, generic Ember.js apps unrelated to Aha!, or non-API product management advice.
---

# Aha! API

## When to Use This Skill

- Calling the Aha! REST API (features, initiatives, goals, ideas, custom fields, webhooks)
- Authenticating to Aha! with API keys or OAuth2
- Building Chrome extension content scripts that run on `*.aha.io` pages
- Parsing Aha! webhook payloads (activity or audit types)
- Extracting data from Aha! published notebooks (shared pages)

## When NOT to Use This Skill

- General Chrome extension architecture — use `chrome-mv3-advanced` instead
- Generic OAuth2 flows not targeting Aha! — use `web-auth-patterns` instead
- General Ember.js SPA development unrelated to Aha!
- Non-API product management advice (strategy, roadmap planning)
- Aha!'s native JavaScript extension SDK (separate product surface, not covered here)

---

## Overview

Aha! exposes a versioned REST API at `https://<subdomain>.aha.io/api/v1/` covering the full product management model: features, epics, releases, initiatives, goals, ideas, custom fields, webhooks, and more. All responses are JSON; all writes use JSON bodies.

`<subdomain>` is your Aha! account subdomain — the prefix before `.aha.io` in your browser URL (e.g. `acme` → `https://acme.aha.io/api/v1/`). Find it in **Settings → Account → General**.

**Strategic hierarchy:** Goals → Initiatives → Features/Epics → Requirements. Each level links up and down; audit webhook record types reflect this (`strategic_imperative` = goal, `strategic_initiative` = initiative).

**Aha! products:** This skill covers the REST API shared by Aha! Roadmaps and Aha! Ideas. Aha! Develop has a partially overlapping API (features, releases) but uses the same base URL and auth.

## Authentication

### API Key (service integrations)

Generate at **Settings → Personal → Developer → Generate API key** (shown once). Use as a bearer token:

```http
Authorization: Bearer <api_key>
```

### OAuth2 (user-facing apps)

| Grant | Endpoint | Use for |
|-------|----------|---------|
| Authorization Code | `GET https://<subdomain>.aha.io/oauth/authorize?response_type=code&client_id=…&redirect_uri=…` | Server-side apps |
| Implicit | Same URL with `response_type=token` | Browser/SPA apps |
| Token exchange | `POST https://<subdomain>.aha.io/oauth/token` with `code`, `client_id`, `client_secret`, `grant_type=authorization_code`, `redirect_uri` | Authorization code flow only |

Register apps at `https://secure.aha.io/oauth/applications`. Token response: `{"access_token": "…", "token_type": "bearer"}`.

**Scopes:** Aha! OAuth has no explicit scope parameter — the token inherits the permissions of the authorizing user. A read-only user's token cannot create or delete records.

Both API key and OAuth token use identical `Authorization: Bearer …` header.

## Base URL and Request Format

```http
https://<subdomain>.aha.io/api/v1/<resource>
Content-Type: application/json
Accept: application/json
Authorization: Bearer <token>
User-Agent: MyApp (contact@example.com)   # recommended
```

Rate limits: **300 req/min, 20 req/sec** per account. HTTP 429 on breach. Response headers: `X-Ratelimit-Limit`, `X-Ratelimit-Remaining`, `X-Ratelimit-Reset` (Unix timestamp).

### Rate limit backoff

```js
async function fetchAha(url, options, retries = 3) {
  for (let attempt = 0; attempt <= retries; attempt++) {
    const resp = await fetch(url, options);
    if (resp.status !== 429) return resp;

    const reset = parseInt(resp.headers.get('X-Ratelimit-Reset') || '0', 10);
    const waitMs = reset
      ? Math.max(0, reset * 1000 - Date.now())
      : Math.min(1000 * 2 ** attempt, 30000);  // exponential backoff, cap 30s

    if (attempt < retries) await new Promise(r => setTimeout(r, waitMs));
  }
  throw new Error('Aha! API rate limit exceeded after retries');
}
```

## Error Handling

| Status | Meaning | Action |
|--------|---------|--------|
| 400 | Bad request (malformed JSON, invalid field value) | Parse `errors` array in response body for field-level detail |
| 401 | Unauthorized (bad or expired token) | Refresh OAuth token or regenerate API key |
| 403 | Forbidden (valid token, insufficient permission) | Check user role — OAuth token inherits authorizing user's permissions |
| 404 | Not found (wrong ID or key) | Verify the `reference_num` or numeric ID; user may lack project access |
| 429 | Rate limit exceeded | See rate limit backoff above |
| 5xx | Server error | Retry with exponential backoff (max 3 attempts) |

Aha! error responses include an `errors` key: `{"errors": ["Name can't be blank"]}`.

## Pagination and Filtering

```http
GET /api/v1/features?page=1&per_page=50   # default 30, max 200
```

Response includes pagination metadata:

```json
{
  "features": [],
  "pagination": {
    "total_records": 342,
    "total_pages": 7,
    "current_page": 1
  }
}
```

### Field selection

```http
GET /api/v1/features/PRJ1-F-1?fields=name,workflow_status,custom_fields
GET /api/v1/features/PRJ1-F-1?fields=*    # all fields
```

### Common filters (features, ideas, initiatives)

| Param | Example | Notes |
|-------|---------|-------|
| `q` | `q=login` | Name search |
| `updated_since` | `updated_since=2025-01-01T00:00:00Z` | ISO8601 UTC |
| `tag` | `tag=backend` | Single tag |
| `assigned_to_user` | `assigned_to_user=user@example.com` | ID or email |

## Features API

```http
GET    /api/v1/releases/:release_id/features          # list in release
GET    /api/v1/products/:product_id/features          # list in product
GET    /api/v1/initiatives/:initiative_id/features    # list on initiative
GET    /api/v1/features/:id                           # get (key like PRJ1-F-1 or numeric)
POST   /api/v1/releases/:release_id/features          # create
PUT    /api/v1/features/:id                           # update
DELETE /api/v1/features/:id
```

Create/update body:

```json
{
  "feature": {
    "name": "Support SSO login",
    "description": "<p>HTML or plain text</p>",
    "workflow_status": { "name": "In progress" },
    "assigned_to_user": { "email": "dev@example.com" },
    "tags": "auth,backend",
    "custom_fields": {
      "priority": "P1",
      "t_shirt_size": "L",
      "target_date": "2025-09-30"
    }
  }
}
```

Key response fields: `id`, `reference_num` (e.g. `PRJ1-F-42`), `name`, `description`, `workflow_status`, `release`, `initiative`, `assigned_to_user`, `progress`, `score`, `tags`, `custom_fields`, `created_at`, `updated_at`.

## Initiatives API

Initiatives are the strategic layer above features. They link to goals above and features below.

```http
GET    /api/v1/products/:product_id/initiatives
GET    /api/v1/goals/:goal_id/initiatives             # filter by goal
GET    /api/v1/initiatives/:id
POST   /api/v1/products/:product_id/initiatives
PUT    /api/v1/initiatives/:id
DELETE /api/v1/initiatives/:id
PUT    /api/v1/initiatives/:id/custom_fields          # update custom fields only
PUT    /api/v1/initiatives/:id/progress               # update progress
```

Filter: `?workflow_status=In%20progress` or by custom field values.

## Goals API (strategic_imperative)

Goals sit above initiatives in the strategy hierarchy.

```http
GET    /api/v1/products/:product_id/goals
GET    /api/v1/goals/:id
POST   /api/v1/products/:product_id/goals
PUT    /api/v1/goals/:id
DELETE /api/v1/goals/:id
GET    /api/v1/initiatives/:initiative_id/goals       # goals for an initiative
```

Internal audit webhook record type for goals: `strategic_imperative`; for initiatives: `strategic_initiative`.

## Ideas API (Ideas Portal)

```http
GET    /api/v1/products/:product_id/ideas
GET    /api/v1/ideas/:id
POST   /api/v1/products/:product_id/ideas
PUT    /api/v1/ideas/:id
DELETE /api/v1/ideas/:id

# Portal
GET    /api/v1/products/:product_id/idea_portals
POST   /api/v1/ideas/:id/promote_to_feature
POST   /api/v1/ideas/:id/promote_to_initiative

# Votes
POST   /api/v1/ideas/:id/endorsements                # cast a vote
GET    /api/v1/ideas/:id/endorsements
DELETE /api/v1/endorsements/:id

# Comments
POST   /api/v1/ideas/:id/comments
GET    /api/v1/ideas/:id/comments
```

Create with portal assignment:

```json
{
  "idea": {
    "name": "Dark mode support",
    "description": "Users want dark mode",
    "submitted_idea_portal_id": "6881234567",
    "visibility": "public",
    "custom_fields": { "priority": "P3" }
  }
}
```

Portal users vs Aha! users are separate identity spaces. Set `"skip_portal": true` to create an idea without associating it with any portal (useful for internal/admin-created ideas).

## Custom Fields

Custom fields are defined account-wide by record type and must be added to a layout before appearing on records.

**Read (GET response):** returned as an **array** of objects:

```json
"custom_fields": [
  { "key": "priority", "name": "Priority", "value": "P1", "type": "string" },
  { "key": "target_date", "name": "Target Date", "value": "2025-06-30", "type": "date" }
]
```

**Write (POST/PUT request body):** use a **flat key-value object** — the format is asymmetric from the read format:

```json
{
  "feature": {
    "custom_fields": {
      "priority": "P1",
      "risk_level": "high",
      "target_date": "2025-06-30",
      "complexity_score": 8,
      "labels": "urgent,backend"
    }
  }
}
```

**Wrong — do not send the read-response array format in writes:**

```json
{
  "feature": {
    "custom_fields": [
      { "key": "priority", "value": "P1" }
    ]
  }
}
```

Field types: `string`, `note` (rich text), `integer`, `float`, `date` (YYYY-MM-DD), `predefined_list` (single select), `predefined_multiple_list` (multi-select, comma-separated string), `url`.

## Webhooks

### Activity Webhooks (outbound — Aha! notifies you)

Configured at **Settings → Integrations → Webhooks**. Fires on create/update/destroy with a ~5-minute delay. Only sends changed fields.

To inspect the exact payload shape for your account: set up a temporary endpoint at [webhook.site](https://webhook.site) or [ngrok](https://ngrok.io), then perform a test create/update in Aha! and examine the inbound request.

Payload structure:

```json
{
  "event": "create",
  "object": {
    "type": "feature",
    "id": "123456",
    "reference_num": "PRJ1-F-42",
    "name": "SSO login"
  },
  "changed_fields": ["name", "workflow_status"]
}
```

### Audit Webhooks (live stream, no delay)

Full create/update/destroy stream for compliance or real-time automation.

```json
{
  "type": "audit",
  "id": "9876543210987654",
  "action": "update",
  "created_at": "2025-06-01T12:00:00Z",
  "user_id": "123456",
  "user_name": "Jane Smith",
  "record": {
    "type": "feature",
    "id": "99887766",
    "reference_num": "PRJ1-F-42"
  },
  "changes": [
    { "field_name": "workflow_status", "value": "In progress" },
    { "field_name": "assigned_to_user_id", "value": "555000" }
  ]
}
```

Audit record types: `feature`, `release`, `requirement`, `strategic_imperative` (goal), `strategic_initiative` (initiative), `attachment`, `comment`, `note`.

### Inbound Webhooks (Aha! receives)

Aha! accepts webhook payloads from GitHub, Asana, and others via integration-specific endpoints:

```http
POST /api/v1/integrations/github/webhook?callback_token=<token>
```

GitHub commit hook sends data in `payload` param (not `webhook`).

## Shared Pages (Published Notebooks)

Published notebooks generate a public read-only HTML page — not an API endpoint:

```
https://<subdomain>.aha.io/published/<hash>
```

No authentication required for public notebooks. Enterprise+ can restrict to SSO or Aha! users only. To extract data from a shared page in a Chrome extension, use DOM parsing against the rendered HTML (see "Reading record IDs from Aha! page" below for selector patterns).

## Chrome Extension / SPA Scraping Patterns

Aha! is built on Ember.js — a client-side SPA where URL changes do not reload the page. Key patterns for Chrome extension content scripts:

### Route change detection (Ember SPA)

```js
// popstate fires on back/forward; pushState must be monkey-patched for link navigation.
// Guard against double-patching if the content script is injected more than once.
window.addEventListener('popstate', handleNavigation);

if (!history.__ahaPatchedByExtension) {
  history.__ahaPatchedByExtension = true;
  const origPush = history.pushState.bind(history);
  history.pushState = function(...args) {
    origPush(...args);
    handleNavigation();
  };
}
```

### Waiting for dynamic content

```js
// Always wait for a stable DOM selector — never rely on URL change alone.
// Guards against resolve() being called after the promise settles (race condition).
function waitForElement(selector, timeout = 8000) {
  return new Promise((resolve, reject) => {
    const el = document.querySelector(selector);
    if (el) return resolve(el);
    let settled = false;
    const mo = new MutationObserver(() => {
      const found = document.querySelector(selector);
      if (found && !settled) { settled = true; mo.disconnect(); resolve(found); }
    });
    mo.observe(document.body, { childList: true, subtree: true });
    setTimeout(() => {
      if (!settled) { settled = true; mo.disconnect(); reject(new Error('timeout: ' + selector)); }
    }, timeout);
  });
}
```

### Reading record IDs from Aha! page

Aha! embeds record reference numbers in DOM attributes and URLs:

```js
// Reference num from URL: /features/PRJ1-F-42
const refMatch = location.pathname.match(/\/(features|initiatives|ideas)\/([A-Z0-9\-]+)/);
const referenceNum = refMatch?.[2]; // "PRJ1-F-42"

// Numeric ID from data attribute (varies by page)
const id = document.querySelector('[data-id]')?.dataset.id;
```

For published notebook pages (`/published/<hash>`), target heading and table elements directly — the HTML is server-rendered and stable.

### Accessing Ember data store (advanced)

> **Warning:** This approach is brittle and breaks across Aha! version updates. Do not use in production integrations. Use the REST API via proxy instead (see below).

```js
const app = window.Ember?.Application?.INSTANCES?.values().next().value;
const store = app?.__container__?.lookup('service:store');
const record = store?.peekRecord('feature', id);
```

### Preferred approach: proxy through REST API

Content scripts read the current page URL/DOM to extract the record reference number, then post a message to the MV3 service worker, which calls the Aha! REST API with the stored bearer token. This is more reliable and version-stable than scraping Ember data models.

## Quick Reference

| Task | Endpoint |
|------|----------|
| List features in release | `GET /api/v1/releases/:id/features` |
| Get feature by key | `GET /api/v1/features/PRJ1-F-42` |
| Create feature | `POST /api/v1/releases/:id/features` |
| Update feature / custom fields | `PUT /api/v1/features/:id` |
| List ideas | `GET /api/v1/products/:id/ideas` |
| Promote idea to feature | `POST /api/v1/ideas/:id/promote_to_feature` |
| List initiatives | `GET /api/v1/products/:id/initiatives` |
| List goals | `GET /api/v1/products/:id/goals` |
| OAuth authorize | `GET https://<subdomain>.aha.io/oauth/authorize?response_type=code&client_id=…` |
| OAuth token exchange | `POST https://<subdomain>.aha.io/oauth/token` |

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Sending `custom_fields` as an array in write requests | Use a flat key-value object in the request body; only GET responses return an array |
| Using the wrong OAuth host | Register apps at `secure.aha.io`; all authorize/token calls use `<subdomain>.aha.io` |
| Expecting immediate activity webhook delivery | Activity webhooks have ~5-min delay; use audit webhooks for real-time needs |
| Not waiting for Ember render after route change | Use MutationObserver + stable selector — URL change is not a render signal |
| Token grants too much access | OAuth token inherits the authorizing user's role — use a least-privilege Aha! user for service integrations |
| Missing `submitted_idea_portal_id` for portal users | Always include when the idea creator is a portal (idea) user |
| Using Ember data store in production | Breaks on Aha! updates — proxy through REST API instead |
| Retrying on 403 | 403 means permission denied, not a transient error — check user role before retrying |

## Sources

- [Aha! REST API Overview](https://www.aha.io/api)
- [OAuth2 Authentication](https://www.aha.io/api/oauth2)
- [Features API](https://www.aha.io/api/resources/features/list_features_in_a_release)
- [Ideas API](https://www.aha.io/api/resources/ideas)
- [Initiatives API](https://www.aha.io/api/resources/initiatives)
- [Goals API](https://www.aha.io/api/resources/goals)
- [Custom Fields API](https://www.aha.io/api/resources/custom_fields/list_options_for_a_custom_field)
- [Webhooks API](https://www.aha.io/api/resources/webhooks/process_github_commit_hook_payload)
- [Activity Webhooks Support Article](https://www.aha.io/support/roadmaps/integrations/aha-api/activity-webhooks-integration)
- [Integrating with Aha!](https://www.aha.io/api/integrating-aha)
