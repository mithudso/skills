<!-- hub-reference-banner -->
> **Reference file — part of the `integration-clients` hub.** Formerly the standalone `jira-extension-client` skill.
> Sibling topics in this family are now reference files under the hubs (`integration-clients`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: jira-extension-client
version: 1.1.0
updated: 2026-05-29
description: >
  Jira REST API integration from Chrome MV3 extensions — cookie auth, OAuth 2.0
  3LO (public client, no secret), API tokens (Basic auth), PAT (Data Center),
  JQL search, issue CRUD, transitions, comments, ADF formatting, rate limiting,
  pagination (nextPageToken and startAt), CORS bypass via host_permissions, and
  service worker termination handling.
  TRIGGER: building a Chrome extension that reads or writes Jira issues; implementing
  Jira auth from extension context; writing JQL for extension dashboards; debugging
  CORS, 401, 403, or 429 errors from a Chrome extension calling Jira; handling Jira
  API pagination from a service worker; creating/updating Jira issues with ADF bodies.
  SKIP: general Jira administration or Jira CLI usage (use jira-cli or
  jira-developer-expert); Jira Forge or Connect app development (server-side context);
  non-Chrome extension browser platforms; Confluence or Bitbucket APIs.
category: developer
tags: [jira, chrome-extension, rest-api, authentication, jql, oauth, atlassian, adf, cors, service-worker]
related_skills: [chrome-mv3-advanced, mv3-service-worker-expert, jira-developer-expert, jira-cli]
whenToUse:
  - "Building a Chrome extension that reads or writes Jira issues"
  - "Jira OAuth 2.0 3LO from a Chrome extension"
  - "Jira API token Basic auth in service worker"
  - "JQL search from Chrome extension service worker"
  - "Handle Jira 429 rate limit in extension"
  - "Jira ADF body format for issue description or comment"
  - "Chrome extension CORS error calling Jira"
  - "Jira pagination nextPageToken vs startAt"
  - "Service worker killed mid-pagination"
  - "Jira Data Center PAT authentication"
whenNotToUse:
  - "Jira CLI, Jira admin, or Jira Forge/Connect app development: use jira-developer-expert"
  - "Confluence or Bitbucket APIs (different surfaces, different auth)"
  - "Non-Chrome browser extension platforms"
---

# Jira Extension Client

Expert reference for integrating Jira REST APIs from a Chrome MV3 extension context. Covers authentication strategies, CORS bypass via host_permissions, JQL search patterns, issue CRUD, transitions, comments, ADF body formatting, pagination, rate limiting, and error handling with re-auth flows.

## When to use this skill

- Building a Chrome extension that reads or writes Jira issues
- Implementing Jira authentication (cookie, OAuth 2.0 3LO, API token, PAT) from extension context
- Writing JQL search queries for extension-driven dashboards or panels
- Handling Jira API pagination (nextPageToken or startAt) from a service worker
- Debugging CORS, 401, 403, or 429 errors when calling Jira from an extension
- Creating or updating Jira issues programmatically, including ADF-formatted descriptions and comments
- Designing rate-limit-aware fetch wrappers for Jira API calls

## When NOT to use this skill

- General Jira administration or Jira CLI usage -- use `jira-cli` or `jira-developer-expert` instead
- Jira Forge app or Connect app development (server-side context, not extension context)
- Non-Chrome browser extension platforms (Firefox WebExtensions have different permission models)
- Atlassian Confluence or Bitbucket APIs -- different API surfaces despite shared auth

## 1. Core Concepts

### API versions

| Version | Base URL (Cloud) | Body format | Status |
|---------|-----------------|-------------|--------|
| v2 | `https://{site}.atlassian.net/rest/api/2/` | Wiki markup / plain text | Stable, widely used |
| v3 | `https://{site}.atlassian.net/rest/api/3/` | Atlassian Document Format (ADF) | Current, required for rich text |
| v2 (DC) | `https://{host}/rest/api/2/` | Wiki markup / plain text | Stable for Data Center |

**Rule of thumb:** Use v2 when you only need plain-text fields (summary, labels, status). Use v3 when you need rich-text bodies (description, comments) on Jira Cloud.

### Extension architecture for Jira calls

All Jira API calls should originate from the **service worker** (background script), never from content scripts. Content scripts are subject to the host page's CSP and cannot set Authorization headers on cross-origin requests.

```text
Content script / Panel UI
    |
    +-- chrome.runtime.sendMessage({ type: 'JIRA_SEARCH', jql: '...' })
    |
    v
Service worker (background.js)
    |
    +-- fetch('https://site.atlassian.net/rest/api/3/search', { headers, credentials })
    |
    v
Jira Cloud / Data Center
```

## 2. Manifest Configuration

### host_permissions

The extension must declare host_permissions for every Jira domain it will contact. This grants the service worker permission to make cross-origin fetch requests and automatically attach cookies for those domains.

```json
{
  "manifest_version": 3,
  "host_permissions": [
    "https://*.atlassian.net/*",
    "https://api.atlassian.com/*"
  ],
  "permissions": [
    "storage",
    "cookies"
  ],
  "optional_host_permissions": [
    "https://*/*"
  ]
}
```

**Key points:**

- `https://*.atlassian.net/*` covers all Jira Cloud tenants.
- `https://api.atlassian.com/*` is required if you use OAuth 2.0 3LO (token exchange and API calls route through this domain).
- For Jira Data Center, add the specific server URL: `"https://jira.company.com/*"`.
- Use `optional_host_permissions` when you want to let the user specify their Jira instance at runtime rather than requiring blanket access at install.
- The `cookies` permission is needed only if you explicitly read/write cookies via `chrome.cookies`; cookie-based fetch auth works without it as long as host_permissions covers the domain.

## 3. Authentication Strategies

### Strategy 1: Cookie-based auth (same-origin session)

Best for extensions where the user is already logged into Jira in their browser. The service worker piggybacks on existing session cookies.

```js
// service-worker.js
async function jiraFetch(siteUrl, path, options = {}) {
  const url = `${siteUrl}/rest/api/2/${path}`;
  const resp = await fetch(url, {
    ...options,
    credentials: 'include',           // attach cookies for siteUrl domain
    headers: {
      'Content-Type': 'application/json',
      'X-Atlassian-Token': 'no-check', // bypass XSRF check on mutating requests
      ...options.headers,
    },
  });
  return resp;
}
```

**When it works:** Jira Data Center (cookie auth still supported). Jira Cloud historically supported it but has deprecated cookie-based and basic-password auth. For Jira Cloud, prefer API tokens or OAuth 2.0.

**XSRF protection:** Jira Server/DC requires an `X-Atlassian-Token: no-check` header on POST/PUT/DELETE requests to bypass the XSRF token check when using cookie auth. Without it, mutating requests return 403.

### Strategy 2: API token (Basic auth) -- Jira Cloud

The user generates an API token at https://id.atlassian.com/manage-profile/security/api-tokens. The extension stores the email + token and sends them as a Basic Authorization header.

```js
// service-worker.js
function buildBasicAuthHeader(email, apiToken) {
  const encoded = btoa(`${email}:${apiToken}`);
  return `Basic ${encoded}`;
}

async function jiraCloudFetch(site, path, { email, apiToken, ...options } = {}) {
  const url = `https://${site}.atlassian.net/rest/api/3/${path}`;
  const resp = await fetch(url, {
    ...options,
    headers: {
      'Authorization': buildBasicAuthHeader(email, apiToken),
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });
  return resp;
}
```

**Storage:** Store the API token in `chrome.storage.session` (cleared on browser restart) or encrypted in `chrome.storage.local` using a vault pattern. Never store tokens in plain text in local storage.

### Strategy 3: Personal Access Token (PAT) -- Jira Data Center

PATs are the recommended auth method for Jira Data Center 8.14+. They are sent as a Bearer token.

```js
async function jiraDCFetch(baseUrl, path, { pat, ...options } = {}) {
  const url = `${baseUrl}/rest/api/2/${path}`;
  const resp = await fetch(url, {
    ...options,
    headers: {
      'Authorization': `Bearer ${pat}`,
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });
  return resp;
}
```

### Strategy 4: OAuth 2.0 (3LO) -- Jira Cloud

OAuth 2.0 with three-legged authorization is the most secure method for Jira Cloud. It uses rotating refresh tokens and scoped access.

**Setup flow:**

1. Register an OAuth 2.0 app at https://developer.atlassian.com/console/myapps/
2. Configure callback URL as your extension's redirect URI
3. Select required scopes (e.g., `read:jira-work`, `write:jira-work`, `read:jira-user`)

**Security: Chrome extensions are public clients.** Never bundle a `client_secret` in extension code -- it is extractable from the .crx package. Atlassian's OAuth 2.0 3LO for browser-based apps supports the authorization code flow without a client_secret (public client mode). If you see `client_secret` in your token exchange, you are misconfigured.

```js
// oauth-flow.js -- triggered from options page or popup
const CLIENT_ID = 'your-client-id';
const REDIRECT_URI = chrome.identity.getRedirectURL('jira');
const SCOPES = 'read:jira-work write:jira-work read:jira-user offline_access';

async function startOAuthFlow() {
  const authUrl = new URL('https://auth.atlassian.com/authorize');
  authUrl.searchParams.set('audience', 'api.atlassian.com');
  authUrl.searchParams.set('client_id', CLIENT_ID);
  authUrl.searchParams.set('scope', SCOPES);
  authUrl.searchParams.set('redirect_uri', REDIRECT_URI);
  authUrl.searchParams.set('response_type', 'code');
  authUrl.searchParams.set('prompt', 'consent');
  // state param for CSRF protection
  const state = crypto.randomUUID();
  authUrl.searchParams.set('state', state);

  const redirectResult = await chrome.identity.launchWebAuthFlow({
    url: authUrl.toString(),
    interactive: true,
  });

  const resultUrl = new URL(redirectResult);
  const code = resultUrl.searchParams.get('code');
  if (resultUrl.searchParams.get('state') !== state) {
    throw new Error('OAuth state mismatch -- possible CSRF');
  }

  return exchangeCodeForTokens(code);
}

async function exchangeCodeForTokens(code) {
  const resp = await fetch('https://auth.atlassian.com/oauth/token', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      grant_type: 'authorization_code',
      client_id: CLIENT_ID,
      code,
      redirect_uri: REDIRECT_URI,
    }),
  });
  const tokens = await resp.json();
  // tokens.access_token, tokens.refresh_token, tokens.expires_in
  await chrome.storage.session.set({ jiraTokens: tokens });
  return tokens;
}

async function refreshAccessToken(refreshToken) {
  const resp = await fetch('https://auth.atlassian.com/oauth/token', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      grant_type: 'refresh_token',
      client_id: CLIENT_ID,
      refresh_token: refreshToken,
    }),
  });
  return resp.json();  // new access_token + new refresh_token (rotating)
}
```

**Getting the Cloud ID:**
After OAuth, you must resolve the accessible Jira sites and their Cloud IDs:

```js
async function getAccessibleResources(accessToken) {
  const resp = await fetch('https://api.atlassian.com/oauth/token/accessible-resources', {
    headers: { 'Authorization': `Bearer ${accessToken}` },
  });
  return resp.json();
  // Returns: [{ id: 'cloud-id', url: 'https://site.atlassian.net', name: 'Site Name', ... }]
}
```

**API calls with OAuth use a different base URL:**

```js
async function jiraOAuthFetch(cloudId, path, accessToken, options = {}) {
  const url = `https://api.atlassian.com/ex/jira/${cloudId}/rest/api/3/${path}`;
  const resp = await fetch(url, {
    ...options,
    headers: {
      'Authorization': `Bearer ${accessToken}`,
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });
  return resp;
}
```

### Authentication decision matrix

| Scenario | Method | Jira target |
|----------|--------|-------------|
| User already logged into Jira in browser | Cookie + `credentials: 'include'` | DC (Cloud deprecated) |
| User provides email + API token in settings | Basic auth header | Cloud |
| User provides PAT in settings | Bearer token header | DC 8.14+ |
| Extension distributed publicly, needs user consent | OAuth 2.0 3LO | Cloud |
| Internal enterprise extension, SSO environment | Cookie (if DC) or PAT | DC |

## 4. Key API Endpoints

### Issue operations

```js
// Get a single issue
GET /rest/api/3/issue/{issueIdOrKey}
GET /rest/api/3/issue/{issueIdOrKey}?fields=summary,status,assignee,priority

// Create an issue
POST /rest/api/3/issue
{
  "fields": {
    "project": { "key": "PROJ" },
    "summary": "Issue title",
    "issuetype": { "name": "Bug" },
    "priority": { "name": "High" },
    "description": {                    // ADF format for v3
      "type": "doc",
      "version": 1,
      "content": [{
        "type": "paragraph",
        "content": [{ "type": "text", "text": "Description text here." }]
      }]
    }
  }
}

// Update an issue (partial update -- send only changed fields)
PUT /rest/api/3/issue/{issueIdOrKey}
{
  "fields": {
    "summary": "Updated title"
  }
}
```

### Comments

```js
// Add a comment (v3 requires ADF body)
POST /rest/api/3/issue/{issueIdOrKey}/comment
{
  "body": {
    "type": "doc",
    "version": 1,
    "content": [{
      "type": "paragraph",
      "content": [{ "type": "text", "text": "Comment from extension." }]
    }]
  }
}

// Add a comment (v2 uses plain text)
POST /rest/api/2/issue/{issueIdOrKey}/comment
{
  "body": "Comment from extension."
}

// List comments
GET /rest/api/3/issue/{issueIdOrKey}/comment?orderBy=-created&maxResults=10
```

### Transitions

```js
// Get available transitions
GET /rest/api/3/issue/{issueIdOrKey}/transitions

// Perform a transition
POST /rest/api/3/issue/{issueIdOrKey}/transitions
{
  "transition": { "id": "31" },
  "update": {
    "comment": [{
      "add": {
        "body": {
          "type": "doc",
          "version": 1,
          "content": [{
            "type": "paragraph",
            "content": [{ "type": "text", "text": "Moved via extension." }]
          }]
        }
      }
    }]
  }
}
```

### Custom fields

Custom fields are referenced by ID, not display name: `customfield_10001`. To discover field IDs:

```js
// List all fields (includes custom fields)
GET /rest/api/3/field

// Response includes:
// { "id": "customfield_10001", "name": "Story Points", "custom": true, "schema": { ... } }
```

## 5. JQL Search Patterns

### Endpoint selection

| Endpoint | Method | Pagination | Notes |
|----------|--------|-----------|-------|
| `/rest/api/3/search` | GET/POST | `startAt` + `maxResults` | Deprecated on Cloud but still functional |
| `/rest/api/3/search/jql` | GET/POST | `nextPageToken` + `maxResults` | Current Cloud endpoint |
| `/rest/api/2/search` | GET/POST | `startAt` + `maxResults` | Data Center / Cloud v2 |

### Common JQL queries

```js
// Issues assigned to current user, ordered by update time
const MY_ISSUES = 'assignee = currentUser() ORDER BY updated DESC';

// Issues in a specific project
const PROJECT_ISSUES = 'project = "PROJ" ORDER BY created DESC';

// Open issues assigned to me
const MY_OPEN = 'assignee = currentUser() AND status != Done ORDER BY priority DESC';

// Issues updated in the last 24 hours
const RECENTLY_UPDATED = 'updated >= -1d ORDER BY updated DESC';

// Issues by reporter
const BY_REPORTER = 'reporter = "user@example.com" ORDER BY created DESC';

// Text search across summary and description
const TEXT_SEARCH = 'text ~ "error timeout" ORDER BY updated DESC';

// Multiple projects with status filter
const CROSS_PROJECT = 'project in (PROJ, TEAM, OPS) AND status in ("In Progress", "Review") ORDER BY priority DESC';

// Issues with a specific label
const BY_LABEL = 'labels = "customer-reported" AND status != Done ORDER BY created DESC';

// Sprint-based (Jira Software)
const CURRENT_SPRINT = 'sprint in openSprints() AND assignee = currentUser()';
```

### Search implementation

> **Note:** Section 8 (Pagination Patterns) provides reusable async generator versions of these functions. The examples below are self-contained for quick reference.

```js
// Using the current /search/jql endpoint (Jira Cloud)
async function searchJql(jql, fields = ['summary', 'status', 'assignee', 'priority'], maxResults = 50) {
  const results = [];
  let nextPageToken = null;

  do {
    const params = new URLSearchParams({
      jql,
      fields: fields.join(','),
      maxResults: String(maxResults),
    });
    if (nextPageToken) params.set('nextPageToken', nextPageToken);

    const resp = await jiraFetch(`search/jql?${params}`);
    if (!resp.ok) {
      const err = await resp.json().catch(() => ({}));
      throw new JiraApiError(resp.status, err.errorMessages, err.errors);
    }

    const data = await resp.json();
    results.push(...data.issues);
    nextPageToken = data.nextPageToken ?? null;
  } while (nextPageToken);

  return results;
}

// Using the legacy /search endpoint (Data Center or Cloud v2 fallback)
async function searchLegacy(jql, fields = ['summary', 'status'], maxResults = 50) {
  const results = [];
  let startAt = 0;
  let total = Infinity;

  while (startAt < total) {
    const params = new URLSearchParams({
      jql,
      fields: fields.join(','),
      startAt: String(startAt),
      maxResults: String(maxResults),
    });

    const resp = await jiraFetch(`search?${params}`);
    if (!resp.ok) {
      const err = await resp.json().catch(() => ({}));
      throw new JiraApiError(resp.status, err.errorMessages, err.errors);
    }

    const data = await resp.json();
    results.push(...data.issues);
    total = data.total;
    startAt += data.maxResults;
  }

  return results;
}
```

## 6. Atlassian Document Format (ADF)

Jira Cloud v3 API requires ADF for rich-text fields (description, comments). ADF is a JSON tree of block and inline nodes.

### Minimal ADF helpers

```js
// Build a simple ADF document from plain text
function textToAdf(text) {
  const paragraphs = text.split('\n\n').filter(p => p.trim().length > 0);
  if (paragraphs.length === 0) {
    paragraphs.push(' '); // ADF requires at least one content node
  }
  return {
    type: 'doc',
    version: 1,
    content: paragraphs.map(paragraph => ({
      type: 'paragraph',
      content: [{ type: 'text', text: paragraph }],
    })),
  };
}

// Build ADF with a heading + paragraphs
function headingAndBody(heading, bodyText) {
  return {
    type: 'doc',
    version: 1,
    content: [
      {
        type: 'heading',
        attrs: { level: 3 },
        content: [{ type: 'text', text: heading }],
      },
      ...bodyText.split('\n\n').map(p => ({
        type: 'paragraph',
        content: [{ type: 'text', text: p }],
      })),
    ],
  };
}

// Build ADF with bold/italic marks
function markedText(text, marks = []) {
  // marks: array of { type: 'strong' } | { type: 'em' } | { type: 'code' }
  return { type: 'text', text, marks };
}

// Build ADF bullet list
function bulletList(items) {
  return {
    type: 'bulletList',
    content: items.map(item => ({
      type: 'listItem',
      content: [{
        type: 'paragraph',
        content: [{ type: 'text', text: item }],
      }],
    })),
  };
}

// Build ADF code block
function codeBlock(code, language = 'javascript') {
  return {
    type: 'codeBlock',
    attrs: { language },
    content: [{ type: 'text', text: code }],
  };
}
```

### v2 vs v3 body format decision

```js
function buildCommentBody(text, apiVersion) {
  if (apiVersion === 3) {
    return { body: textToAdf(text) };
  }
  return { body: text };  // v2 accepts plain string
}
```

## 7. Rate Limiting and Retry

### Rate limit headers from Jira Cloud

| Header | Meaning |
|--------|---------|
| `X-RateLimit-Limit` | Max requests allowed in the window |
| `X-RateLimit-Remaining` | Requests left in the current window |
| `X-RateLimit-Reset` | Unix timestamp when the window resets |
| `Retry-After` | Seconds to wait before retrying (on 429) |
| `RateLimit-Reason` | Which limit was hit (`jira-quota-global-based`, `jira-burst-based`, etc.) |

### Retry wrapper with exponential backoff

```js
class JiraClient {
  constructor({ baseFetch, maxRetries = 4, baseDelay = 1000 }) {
    this._fetch = baseFetch;
    this.maxRetries = maxRetries;
    this.baseDelay = baseDelay;
  }

  async request(path, options = {}) {
    let lastError;

    for (let attempt = 0; attempt <= this.maxRetries; attempt++) {
      const resp = await this._fetch(path, options);

      if (resp.status === 429) {
        const retryAfter = parseInt(resp.headers.get('Retry-After') || '0', 10);
        const backoff = retryAfter > 0
          ? retryAfter * 1000
          : this.baseDelay * Math.pow(2, attempt) + Math.random() * 500; // jitter
        console.warn(`[jira] 429 on ${path}, retrying in ${backoff}ms (attempt ${attempt + 1})`);
        await this._sleep(backoff);
        continue;
      }

      if (resp.status === 401) {
        // Attempt token refresh before retrying
        const refreshed = await this._refreshAuth();
        if (refreshed && attempt < this.maxRetries) continue;
      }

      if (!resp.ok) {
        lastError = await this._parseError(resp);
        if (resp.status >= 500 && attempt < this.maxRetries) {
          const backoff = this.baseDelay * Math.pow(2, attempt) + Math.random() * 500;
          await this._sleep(backoff);
          continue;
        }
        throw lastError;
      }

      return resp;
    }

    throw lastError || new Error('Max retries exceeded');
  }

  _sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  async _refreshAuth() {
    // Override in subclass -- return true if refresh succeeded
    return false;
  }

  async _parseError(resp) {
    try {
      const body = await resp.json();
      return new JiraApiError(resp.status, body.errorMessages, body.errors);
    } catch {
      return new JiraApiError(resp.status, [resp.statusText]);
    }
  }
}

class JiraApiError extends Error {
  constructor(status, messages = [], fieldErrors = {}) {
    super(`Jira API ${status}: ${messages.join('; ')}`);
    this.status = status;
    this.messages = messages;
    this.fieldErrors = fieldErrors;
  }
}
```

### Proactive rate-limit awareness

```js
// Track remaining quota from response headers
function checkRateLimitHeaders(resp) {
  const remaining = parseInt(resp.headers.get('X-RateLimit-Remaining') || '-1', 10);
  const limit = parseInt(resp.headers.get('X-RateLimit-Limit') || '-1', 10);

  if (remaining >= 0 && remaining < limit * 0.1) {
    console.warn(`[jira] Rate limit warning: ${remaining}/${limit} remaining`);
  }

  return { remaining, limit };
}
```

## 8. Pagination Patterns

### nextPageToken pagination (Jira Cloud current)

The `/search/jql` endpoint returns a `nextPageToken` field. Pass it on the next request to get the next page. There is no `total` field -- you paginate until `nextPageToken` is absent.

```js
async function* paginateJql(jql, fields, pageSize = 50) {
  let nextPageToken = null;

  do {
    const params = new URLSearchParams({ jql, maxResults: String(pageSize) });
    if (fields?.length) params.set('fields', fields.join(','));
    if (nextPageToken) params.set('nextPageToken', nextPageToken);

    const resp = await jiraClient.request(`search/jql?${params}`);
    const data = await resp.json();

    yield data.issues;
    nextPageToken = data.nextPageToken ?? null;
  } while (nextPageToken);
}

// Usage
for await (const page of paginateJql('project = PROJ', ['summary', 'status'])) {
  processBatch(page);
}
```

### startAt pagination (legacy / Data Center)

```js
async function* paginateLegacy(jql, fields, pageSize = 50) {
  let startAt = 0;
  let total = Infinity;

  while (startAt < total) {
    const params = new URLSearchParams({
      jql,
      startAt: String(startAt),
      maxResults: String(pageSize),
    });
    if (fields?.length) params.set('fields', fields.join(','));

    const resp = await jiraClient.request(`search?${params}`);
    const data = await resp.json();

    yield data.issues;
    total = data.total;
    startAt += data.maxResults;
  }
}
```

### Performance tip: limit fields

Requesting fewer fields allows Jira Cloud to return larger pages (up to 5000 results with minimal fields). Always specify the `fields` parameter rather than fetching all fields.

## 9. Error Handling and Re-auth

### Error taxonomy

| Status | Meaning | Action |
|--------|---------|--------|
| 400 | Bad request (invalid JQL, malformed ADF, missing required field) | Parse `errorMessages` and `errors` for field-level detail; fix and retry |
| 401 | Unauthorized (expired token, invalid credentials) | Trigger re-auth flow; refresh OAuth token or prompt user for new credentials |
| 403 | Forbidden (insufficient permissions, missing XSRF header) | Check scope/permissions; add `X-Atlassian-Token: no-check` for cookie auth |
| 404 | Issue or resource not found | Verify issue key; user may lack project access |
| 429 | Rate limited | Respect `Retry-After`; backoff with jitter |
| 5xx | Server error | Retry with exponential backoff; max 3-4 attempts |

### 401 re-auth flow

```js
class JiraOAuthClient extends JiraClient {
  async _refreshAuth() {
    const { jiraTokens } = await chrome.storage.session.get('jiraTokens');
    if (!jiraTokens?.refresh_token) return false;

    try {
      const newTokens = await refreshAccessToken(jiraTokens.refresh_token);
      await chrome.storage.session.set({ jiraTokens: newTokens });
      return true;
    } catch (err) {
      console.error('[jira] Token refresh failed, user must re-authenticate', err);
      // Notify the user via badge or notification
      chrome.action.setBadgeText({ text: '!' });
      chrome.action.setBadgeBackgroundColor({ color: '#e74c3c' });
      return false;
    }
  }
}
```

### Field validation error parsing

```js
// Jira returns structured errors for invalid fields
// { "errorMessages": ["..."], "errors": { "summary": "Summary is required", "priority": "..." } }

function formatFieldErrors(jiraError) {
  const parts = [];
  if (jiraError.messages?.length) {
    parts.push(...jiraError.messages);
  }
  if (jiraError.fieldErrors) {
    for (const [field, msg] of Object.entries(jiraError.fieldErrors)) {
      parts.push(`${field}: ${msg}`);
    }
  }
  return parts.join('\n');
}
```

### Handling ADF validation errors

Common ADF mistakes that cause 400 errors:

- Missing `version: 1` at the doc root
- Using `text` as a block node (it is inline-only; must be inside a `paragraph`)
- Empty `content` arrays (at least one child node is required)
- Invalid mark types (only `strong`, `em`, `code`, `underline`, `strike`, `textColor`, `link`, `subsup` are allowed on text nodes)

## 10. CORS Considerations

### Jira Cloud

Jira Cloud does **not** support CORS preflight requests from arbitrary origins. This means browser-initiated requests from content scripts or web pages will fail. The solution is to route all requests through the service worker, which is exempt from CORS when host_permissions are declared.

### Jira Data Center

Jira DC does not support preflighted requests for CORS (confirmed in JRASERVER-59101). Same solution: use the service worker.

### Extension pages

Extension-owned pages (popup, options, dashboard, panel in iframe) can also make direct fetch calls to Jira if host_permissions cover the domain. They share the extension's origin and bypass CORS.

### Content scripts

Content scripts inherit the host page's origin for network requests. Even with host_permissions, content script fetch calls may trigger CORS issues. Always relay through the service worker:

```js
// content-script.js
const response = await chrome.runtime.sendMessage({
  type: 'JIRA_REQUEST',
  method: 'GET',
  path: 'issue/PROJ-123',
  params: { fields: 'summary,status' },
});

// service-worker.js
chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
  if (msg.type === 'JIRA_REQUEST') {
    handleJiraRequest(msg).then(sendResponse);
    return true; // async response
  }
});

async function handleJiraRequest({ method, path, params, body }) {
  const url = new URL(`${JIRA_BASE}/${path}`);
  if (params) {
    for (const [k, v] of Object.entries(params)) url.searchParams.set(k, v);
  }
  const resp = await jiraClient.request(url.toString(), {
    method: method || 'GET',
    body: body ? JSON.stringify(body) : undefined,
  });
  return resp.json();
}
```

## 11. Anti-Patterns

1. **Calling Jira from content scripts.** Always route through the service worker. Content script fetch is subject to the host page CSP and CORS.

2. **Storing API tokens in plain text.** Use `chrome.storage.session` (ephemeral) or encrypt with a vault pattern in `chrome.storage.local`.

3. **Fetching all fields.** Always specify `fields=summary,status,...` to reduce payload and improve pagination limits. Fetching all fields drastically reduces maxResults caps.

4. **Ignoring Retry-After on 429.** Retrying immediately will extend the backoff period. Always respect the header.

5. **Using startAt on the /search/jql endpoint.** The current Cloud endpoint uses `nextPageToken`. Using `startAt` will either be ignored or cause errors.

6. **Hardcoding issue type names.** Issue type names vary by project and locale. Use the `/rest/api/3/issue/createmeta` endpoint (or its v2 equivalent) to discover available issue types for a project.

7. **Sending wiki markup to v3 endpoints.** The v3 API expects ADF, not wiki markup. Use v2 if you only have plain text and cannot build ADF.

8. **Not handling service worker termination.** MV3 service workers idle-terminate after 30 seconds. Long-running pagination loops must checkpoint state to `chrome.storage.session` and resume via `chrome.alarms` if needed.

9. **Ignoring OAuth token rotation.** Atlassian uses rotating refresh tokens. Each refresh response contains a new refresh token that replaces the old one. Failing to store the new refresh token will lock the user out.

10. **Using the deprecated `/rest/api/3/search` POST endpoint without checking.** Atlassian has announced deprecation of the legacy search endpoint on Cloud. Monitor deprecation notices and migrate to `/search/jql`.

## 12. Complete Client Example

```js
// jira-extension-client.js -- drop-in module for a Chrome MV3 extension

const JIRA_DEFAULTS = {
  maxRetries: 3,
  baseDelay: 1000,
  pageSize: 50,
  apiVersion: 3,
};

class JiraExtensionClient {
  constructor(config) {
    this.config = { ...JIRA_DEFAULTS, ...config };
    this.baseUrl = this._buildBaseUrl();
  }

  _buildBaseUrl() {
    const { site, cloudId, serverUrl, apiVersion } = this.config;
    if (cloudId) {
      return `https://api.atlassian.com/ex/jira/${cloudId}/rest/api/${apiVersion}`;
    }
    if (site) {
      return `https://${site}.atlassian.net/rest/api/${apiVersion}`;
    }
    if (serverUrl) {
      return `${serverUrl}/rest/api/2`; // DC always uses v2
    }
    throw new Error('JiraExtensionClient: provide site, cloudId, or serverUrl');
  }

  async getIssue(issueKey, fields) {
    const params = fields ? `?fields=${fields.join(',')}` : '';
    return this._request(`issue/${issueKey}${params}`);
  }

  async createIssue(fields) {
    return this._request('issue', { method: 'POST', body: JSON.stringify({ fields }) });
  }

  async updateIssue(issueKey, fields) {
    return this._request(`issue/${issueKey}`, {
      method: 'PUT',
      body: JSON.stringify({ fields }),
    });
  }

  async addComment(issueKey, text) {
    const body = this.config.apiVersion === 3
      ? { body: textToAdf(text) }
      : { body: text };
    return this._request(`issue/${issueKey}/comment`, {
      method: 'POST',
      body: JSON.stringify(body),
    });
  }

  async transitionIssue(issueKey, transitionId, comment) {
    const payload = { transition: { id: transitionId } };
    if (comment) {
      payload.update = {
        comment: [{ add: buildCommentBody(comment, this.config.apiVersion) }],
      };
    }
    return this._request(`issue/${issueKey}/transitions`, {
      method: 'POST',
      body: JSON.stringify(payload),
    });
  }

  async search(jql, fields, maxResults) {
    const size = maxResults || this.config.pageSize;
    const params = new URLSearchParams({ jql, maxResults: String(size) });
    if (fields?.length) params.set('fields', fields.join(','));
    return this._request(`search/jql?${params}`);
  }

  async *searchAll(jql, fields) {
    let nextPageToken = null;
    do {
      const params = new URLSearchParams({
        jql,
        maxResults: String(this.config.pageSize),
      });
      if (fields?.length) params.set('fields', fields.join(','));
      if (nextPageToken) params.set('nextPageToken', nextPageToken);

      const data = await this._request(`search/jql?${params}`);
      yield data.issues;
      nextPageToken = data.nextPageToken ?? null;
    } while (nextPageToken);
  }

  // --- internal ---

  async _request(path, options = {}) {
    // Implemented by subclass or composition with auth strategy
    throw new Error('_request must be implemented by auth-specific subclass');
  }
}

// Concrete subclass: API token (Basic auth) for Jira Cloud
class JiraBasicAuthClient extends JiraExtensionClient {
  constructor(config) {
    super(config);
    this._retrier = new JiraClient({
      baseFetch: (url, opts) => fetch(url, opts),
      maxRetries: config.maxRetries || 3,
    });
  }

  async _request(path, options = {}) {
    const { email, apiToken } = this.config;
    const url = path.startsWith('http') ? path : `${this.baseUrl}/${path}`;
    const resp = await this._retrier.request(url, {
      ...options,
      headers: {
        'Authorization': buildBasicAuthHeader(email, apiToken),
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });
    return resp.json();
  }
}

// Usage:
// const client = new JiraBasicAuthClient({ site: 'myteam', email, apiToken });
// const issue = await client.getIssue('PROJ-123', ['summary', 'status']);
```

## 13. Troubleshooting

| Symptom | Likely cause | Fix |
|---------|-------------|-----|
| 401 on every request | Expired token or wrong auth header format | Check `Authorization` header; refresh OAuth token; verify API token is base64-encoded correctly |
| 403 on POST/PUT with cookie auth | Missing XSRF header | Add `X-Atlassian-Token: no-check` header |
| 403 with OAuth | Missing scope | Add required scope in developer console and re-authorize |
| CORS error in content script | Fetching Jira directly from content script | Route through service worker |
| 400 "Unrecognized field" | Sending ADF to v2 endpoint or wiki markup to v3 | Match body format to API version |
| 400 on issue create | Missing required fields | Call `/issue/createmeta` to discover required fields for the project/issue type |
| Pagination loop never ends | `nextPageToken` keeps returning same value | Bug in Jira Cloud; add a max-pages safety limit and log the token values |
| Empty search results | JQL syntax error silently returns 0 results | Validate JQL via `/jql/autocomplete` or test in Jira UI first |
| Service worker killed mid-pagination | MV3 30-second idle timeout | Checkpoint progress to `chrome.storage.session`; use `chrome.alarms` to resume |
| `net::ERR_FAILED` on fetch | host_permissions missing for the Jira domain | Add the domain to manifest `host_permissions` |

## 14. References

- [Jira Cloud REST API v3 documentation](https://developer.atlassian.com/cloud/jira/platform/rest/v3/intro/)
- [Jira Cloud REST API v2 documentation](https://developer.atlassian.com/cloud/jira/platform/rest/v2/intro/)
- [Jira Data Center REST API examples](https://developer.atlassian.com/server/jira/platform/jira-rest-api-examples/)
- [OAuth 2.0 (3LO) for Jira Cloud](https://developer.atlassian.com/cloud/jira/platform/oauth-2-3lo-apps/)
- [Jira Cloud scopes for OAuth 2.0](https://developer.atlassian.com/cloud/jira/platform/scopes-for-oauth-2-3LO-and-forge-apps/)
- [Basic auth for Jira Cloud REST APIs](https://developer.atlassian.com/cloud/jira/platform/basic-auth-for-rest-apis/)
- [Personal access tokens for Data Center](https://confluence.atlassian.com/enterprise/using-personal-access-tokens-1026032365.html)
- [Jira Cloud rate limiting](https://developer.atlassian.com/cloud/jira/platform/rate-limiting/)
- [Atlassian Document Format structure](https://developer.atlassian.com/cloud/jira/platform/apis/document/structure/)
- [Cookie-based auth deprecation notice (Cloud)](https://developer.atlassian.com/cloud/jira/platform/deprecation-notice-basic-auth-and-cookie-based-auth/)
- [Chrome extension cross-origin network requests](https://developer.chrome.com/docs/extensions/develop/concepts/network-requests)
- [Chrome MV3 host_permissions](https://developer.chrome.com/docs/extensions/develop/concepts/declare-permissions)
- [Jira search/jql endpoint (Cloud)](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-search/)
- [JRASERVER-59101: CORS preflight not supported](https://jira.atlassian.com/browse/JRASERVER-59101)
