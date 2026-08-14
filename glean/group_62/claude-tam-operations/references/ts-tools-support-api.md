<!-- hub-reference-banner -->
> **Reference file — part of the `tam-operations` hub.** Formerly the standalone `ts-tools-support-api` skill.
> Sibling topics in this family are now reference files under the hubs (`tam-operations`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: ts-tools-support-api
version: 1.1.1
updated: "2026-06-01"
description: >-
  MongoDB TS Tools Support API implementation patterns — auth cascade (Hub-tab cookie →
  bearer fallback), normalization for 6 comment envelope shapes, BFS account-cases
  extraction, case-number zero-padding, paginated comments, fetchWithTimeout retry/backoff,
  getCaseBundle parallel fetch, case enrichment (DOM+API merge), severity model, JIRA
  reference extraction, Socket.IO stale-signal pattern. TRIGGER: building or debugging code
  calling the TS Tools Support API; writing authenticatedFetch, getCaseBundle,
  searchSupport, normalizeCommentsEnvelope, or ts-tools-client; hub-tab cookie vs bearer
  fallback, silent auth redirects, comment visibility defaults, case-number normalization.
  SKIP: non-TS-Tools Chrome extension architecture → chrome-extension-expert; Jira API →
  integration-clients; MongoDB query/performance → mongodb-expert; Glean/LLM pipelines →
  integration-clients (glean-dev); endpoint shapes only → tam-operations
  (tstools-reference).
autoDetect:
  - ts-tools-client
  - ts-tools-normalizers
  - case-enricher
  - authenticatedFetch
  - getCaseBundle
  - getCaseByNumber
  - getCaseComments
  - searchSupport
  - probeSupportAuthentication
  - hub.corp.mongodb.com
  - support-api.ts-tools
category: developer
tags:
  - ts-tools
  - support-api
  - case-management
  - authentication
  - chrome-extension
  - response-normalization
  - pagination
related_skills:
  - tam-operations
  - case-mcp-server-guide
  - chrome-extension-expert
---

# TS Tools Support API — Implementation Patterns

Patterns for building, debugging, and reviewing integrations with the MongoDB TS Tools Support API inside Chrome extensions and Node services.

**For endpoint shapes and domain overview, see `tam-operations` (references/tstools-reference.md).** This skill covers the implementation layer: auth cascade, normalizers, pagination, error handling, and the case bundle pattern.

## 0. Scope and Selection

**Use this skill when:** building, debugging, or reviewing code that calls the TS Tools Support API — case fetches, account lookups, comment pagination, search, auth probing, response normalization, or severity-driven refresh logic.

**Do not use for:** Chrome extension architecture (use `chrome-extension-expert`, references/chrome-dev.md), Jira API (use `integration-clients`, references/jira-extension-client.md), MongoDB query/performance (use `mongodb-expert`), Glean/LLM pipelines (use `integration-clients`, references/glean-dev.md), endpoint shape reference only (use `tam-operations`, references/tstools-reference.md).

**Code examples below are simplified for clarity.** For exact implementations, refer to source files listed in Section 16. This file runs slightly long because it is the implementation-pattern guide: all reference material (endpoint catalog, schemas, enums, paging shapes) has been moved to `tam-operations` (references/tstools-reference.md), and what remains is load-bearing normalizer, fetch, and error-handling code.

## 1. Service Identity

| Property | Value |
|---|---|
| Production host | `https://support-api.ts-tools.prod.corp.mongodb.com` |
| Staging host | `https://support-api.ts-tools.staging.corp.mongodb.com` |
| Hub frontend | `https://hub.corp.mongodb.com` |
| Public REST base | `//api.support.mongodb.com/public/v1.0` |
| OpenAPI spec | `GET /public/v1.0/api-docs` (requires auth) |
| Realtime socket | `wss://support-api.ts-tools.prod.corp.mongodb.com/socket.io/?EIO=4&transport=websocket` |
| Contact | tstools@mongodb.com, Slack `#ts-supporttools`, Jira `TSTOOLS` |

**Legacy URL migration:** The old host `https://support-api.ts-tools.prod.corp.mongodb.com` is transparently remapped to `https://hub.corp.mongodb.com` by the normalizer. Always call `normalizeBaseUrl()` before constructing request URLs.

## 2. Authentication Cascade

The client uses a two-tier auth fallback with a Hub-only guard. Preserve this order in any integration.

### Tier 1: Hub-tab cookie-backed fetch

Execute the fetch inside the Hub tab via `chrome.scripting.executeScript`. The browser attaches session cookies automatically because the request originates from the Hub origin.

```js
// Executed inside the Hub tab context
const response = await fetch(fullUrl, {
  credentials: 'include',
  headers: { Accept: 'application/json' }
});
```

On 401/403, call `openHubForAuth()` to open or refresh the Hub tab, then retry once.

### Hub-only guard

When the API base URL resolves to the Hub origin (`hub.corp.mongodb.com`), return `null` immediately after Hub-tab fetch fails. Bearer fallback is intentionally skipped — Hub-origin endpoints only accept cookie sessions.

```js
if (hubApiBase) {
  return null;  // Do not attempt bearer fallback for Hub-origin APIs
}
```

### Tier 2: Bearer token fallback (non-Hub endpoints only)

```js
async function bearerFetch(url, token, options = {}) {
  if (!token) return null;
  return extensionFetch(url, {
    ...options,
    headers: { ...options.headers, Authorization: `Bearer ${token}` }
  });
}
```

### Accepted cookie names

```js
const COOKIE_CANDIDATES = [
  'support', 'stg_support', 'dev_support',
  'HUB', 'stg_HUB', 'dev_HUB',
  'internal_user',
  '__Secure-next-auth.session-token'
];
```

### Auth probe

```
GET /status/auth  →  200 when authenticated
```

```js
const probe = await probeSupportAuthentication({ tabId, pageUrl, baseUrl, token });
// probe.ok            -- boolean, true if any auth path works
// probe.hasCookieAuth -- has valid session cookies
// probe.hasHubTab     -- an open Hub tab exists
// probe.hasBearerToken -- bearer token is configured
```

### API key auth endpoint

```
POST /auth
Body: { "username": "test@email.com", "key": "EREPIJ..." }
Response: { "token": "..." }
```

For first-party integrations (Chrome extensions, portal UIs), always prefer cookie auth.

## 3. Core API Endpoints

**Path prefix convention:** The OpenAPI spec documents paths without a prefix (e.g. `/case/{number}`). The Hub frontend proxies these under `/api/` (e.g. `/api/case/{number}`). Use the `/api/` prefix when constructing URLs via `authenticatedFetch` against the Hub base URL.

The full 92-endpoint catalog grouped by domain (case lifecycle, search, account, attachment, project, AWS S3) lives in the endpoint reference — see `tam-operations` (references/tstools-reference.md). The endpoints this client's code calls directly:

| Method | Path | Used by |
|---|---|---|
| GET | `/case/{number}` | `getCaseByNumber`, case bundle |
| POST | `/case/{number}/comments` | `fetchAllCaseComments` (paginated) |
| GET | `/case/{number}/stage` | `fetchCaseStage` |
| GET | `/case/{number}/ai/nextAction` | `fetchCaseNextAction` |
| POST | `/account/v1/{id}/activities` | `searchCasesByAccount` (Coveo CQL) |
| POST | `/search` | `searchSupport` (global) |
| GET | `/status/auth` | `probeSupportAuthentication` |
| POST | `/aws/signedUpload` | diagnostic upload flow |

## 4. Case Schema

The `FullCase` field list, enums (`Severity`, `CaseStatus`, `CaseClosedReason`), and identifier patterns (`IdPattern`, `MongoIdPattern`, `CaseNumber`) live in the endpoint reference — see `tam-operations` (references/tstools-reference.md). Implementation-critical facts:

- **Visibility:** `partner` (MongoDB-internal only) is the default on comment writes; pass `public` for customer-visible comments
- **CaseNumber** — 8-digit zero-padded numeric string (e.g. `"01234567"`); normalize before any API call. The normalized case-context shape this client produces is in Section 11

## 5. Case Severity Model

### Severity ranking

```js
function getSeverityRank(severity = '') {
  const normalized = String(severity).trim().toLowerCase();
  if (!normalized) return 999;
  if (normalized.includes('sev1') || normalized.includes('critical')) return 1;
  if (normalized.includes('sev2') || normalized.includes('high')) return 2;
  if (normalized.includes('sev3') || normalized.includes('medium')) return 3;
  if (normalized.includes('sev4') || normalized.includes('low')) return 4;
  return 5; // unrecognized non-empty string
}
```

### Severity display metadata

| Rank | Label | Color | Notes |
|---|---|---|---|
| 1 | Sev1 | `#b42318` (red) | Critical — production down, no workaround |
| 2 | Sev2 | `#f79009` (orange) | High — reduced capacity, no easy workaround |
| 3 | Sev3 | `#f2c94c` (yellow) | Medium — workaround exists |
| 4 | Sev4 | `#2e90fa` (blue) | Low — non-critical, enhancements, questions |
| 5+ | Verify | `#667085` (gray) | Unrecognized severity token |
| 999 | Verify | `#667085` (gray) | Empty severity |

### Refresh targets by severity

| Severity | Tracker refresh target |
|---|---|
| Sev1, Sev2 (critical) | 1 minute |
| Sev3, Sev4 (non-critical) | 60 minutes |

Note: these are tracker refresh intervals, not contractual customer SLAs.

### Ownerless case alert logic

A case triggers an ownerless alert when all of:
- `authOk` is true
- status is NOT closed/resolved/done/duplicate/cancelled
- severity is Sev1 (or P1/S1 variants)
- owner is empty, `Unassigned`, `None`, `N/A`, `TBD`, or `Not Assigned`

## 6. Response Normalization

The API returns varying envelope shapes. Normalizers extract consistent shapes from any variant.

### Case envelope extraction

```js
function extractCaseEnvelope(payload) {
  if (!payload || typeof payload !== 'object') return {};
  return {
    case: payload.case || payload,
    comments: payload.comments?.comments
           || payload.comments?.combined
           || payload.comments || payload.combined || [],
    acl: payload.acl || {},
    attachments: payload.attachments || [],
    success: payload.success !== false
  };
}
```

### Comment envelope normalization

Comments arrive in six different shapes:

```js
function normalizeCommentsEnvelope(payload) {
  if (!payload || typeof payload !== 'object') return [];
  if (Array.isArray(payload.comments?.comments)) return payload.comments.comments;
  if (Array.isArray(payload.comments?.combined)) return payload.comments.combined;
  if (Array.isArray(payload.comments)) return payload.comments;
  if (Array.isArray(payload.combined)) return payload.combined;
  if (Array.isArray(payload.docs)) return payload.docs;
  if (Array.isArray(payload)) return payload;
  return [];
}
```

### Account cases normalization

Uses BFS traversal to find case records in any nesting:

```js
function normalizeAccountCasesEnvelope(payload) {
  const rows = [];
  const queue = [payload];
  while (queue.length) {
    const current = queue.shift();
    if (Array.isArray(current)) { queue.push(...current); continue; }
    if (typeof current !== 'object') continue;
    rows.push(current);
    queue.push(
      current.cases, current.openCases, current.rows,
      current.items, current.results, current.docs
    );
  }
  // Deduplicate by case number, keeping the record with most populated fields
}
```

### Search result normalization

```js
function normalizeSupportSearchResults(payload) {
  const rawResults = Array.isArray(payload) ? payload
    : payload?.results || payload?.items || payload?.docs || [];
  return rawResults.map(entry => ({
    searchType: entry.searchType || entry.search_type || entry.type || '',
    id: entry.Id || entry.id || '',
    caseNumber: normalizeCaseNumber(entry.CaseNumber || entry.caseNumber || ''),
    subject: entry.Subject || entry.subject || '',
    accountId: entry.Account?.Id || entry.account_id || '',
    accountName: entry.Account?.Name || entry.account_name || '',
    projectId: entry.Project__r?.Id || entry.project_id || '',
    projectName: entry.Project__r?.Name || entry.project_name || ''
  })).filter(e => e.id || e.caseNumber || e.name);
}
```

### Case number normalization

```js
function normalizeCaseNumber(value = '') {
  const match = String(value).match(/\b0?\d{7,8}\b/);
  return match ? match[0].padStart(8, '0').slice(-8) : '';
}
// Both '1234567' and '01234567' normalize to '01234567'
```

## 7. Pagination

`PagingObject` (list endpoints) and `PaginateObject` (comment endpoints) shapes live in the endpoint reference — see `tam-operations` (references/tstools-reference.md). The implementation pattern below is what matters here.

### Paginated comment fetching

```js
async function fetchAllCaseComments({ caseNumber, tabId, pageUrl, baseUrl,
                                       token, maxPages = 20 }) {
  const allComments = [];
  let page = 1;
  let totalPages = 1;

  while (page <= totalPages && page <= maxPages) {
    const payload = await authenticatedFetch({
      tabId, pageUrl, baseUrl,
      path: `/api/case/${encodeURIComponent(caseNumber)}/comments` +
            `?page=${page}&length=50&sortOrder=DESC`,
      token, allowEmpty: true
    });
    if (!payload) break;
    const docs = normalizeCommentsEnvelope(payload);
    allComments.push(...docs);
    const totalCount = Number(
      payload?.allCommentCount ?? payload?.count ?? payload?.totalCount ?? docs.length
    );
    totalPages = totalCount > 0 ? Math.ceil(totalCount / 50) : 1;
    if (docs.length < 50 || allComments.length >= totalCount) break;
    page += 1;
  }
  return allComments;
}
```

Key safeguards: page cap (`maxPages = 20`), early exit on short page, total-count convergence check.

## 8. Error Handling and Retry

### fetchWithTimeout wrapper

```js
const response = await fetchWithTimeout(url, init, {
  timeoutMs: 10_000,
  retries: 1,          // retries AFTER first attempt
  backoffBaseMs: 500,  // exponential: 500ms, 1000ms, 2000ms...
  logger,
  label: 'ts-tools'
});
```

**Retryable:** status 429, status >= 500, network TypeError, AbortError (timeout).
**Non-retryable:** 4xx client errors (except 429), successful non-JSON responses.

### Response validation

```js
function isAcceptableApiResponse(response, requestUrl, { allowEmpty } = {}) {
  if (!response?.ok) return false;
  if (didRedirectFromRequest(requestUrl, response.url, response.redirected))
    return false;  // Auth redirect to login page
  if (response.parsedAsJson) return true;
  if (allowEmpty && (response.body == null || response.body === '')) return true;
  if (looksLikeHtmlDocument(response.body)) return false;  // Got login page HTML
  return false;
}
```

Redirect detection catches silent auth redirects that return 200 with an HTML login page.

### Error classification for dashboard

| Code | Cause |
|---|---|
| `timeout` | AbortError or FetchTimeoutError |
| `network` | TypeError (DNS, connection refused) |
| `http-401` | TS Tools session expired |
| `http-403` | scope/permission rejected |
| `http-429` | rate limited |
| `http-4xx` | other client error |
| `http-5xx` | server error |
| `no-listener` | chrome.runtime messaging with no receiver |
| `threw` | unexpected runtime exception |

## 9. Search Patterns

### Global search

```js
const results = await searchSupport({
  query: 'replication lag Atlas',
  tabId, pageUrl, baseUrl, token,
  limit: 25
});
// Returns: [{ searchType, id, caseNumber, subject, accountId, accountName, ... }]
```

### Account case search (Coveo-backed)

```js
const cases = await searchCasesByAccount({
  accountId: '5003000000D8cuI',
  tabId, pageUrl, baseUrl, token,
  query: '(@documenttype==Case) AND @activitysortdate<=today'
});
```

Calls `POST /api/account/v1/{accountId}/activities` with CQL body. Default fields:

```js
fieldsToInclude: [
  'case', 'subject', 'sfstatus', 'sfseverity__c',
  'sflastmodifieddate', 'sfcreateddate', 'sfownername', 'activitysortdate'
]
```

## 10. The Case Bundle Pattern

`getCaseBundle` is the primary high-level fetch — runs 5 API calls in parallel:

```js
const bundle = await getCaseBundle({
  caseNumber: '01234567',
  tabId, pageUrl, baseUrl, token
});

// bundle.authOk               -- boolean
// bundle.case                 -- normalized case object
// bundle.case.stagePayload    -- from /case/{number}/stage
// bundle.case.nextActionPayload -- from /case/{number}/ai/nextAction
// bundle.acl                  -- access control list
// bundle.comments             -- all comments (paginated fetch)
// bundle.attachments          -- (reserved, currently [])
```

The five parallel calls are `probeSupportAuthentication`, `getCaseByNumber`, `fetchAllCaseComments`, `fetchCaseStage`, and `fetchCaseNextAction` (paths in the Section 3 table).

## 11. Case Enrichment Pipeline

Raw data from DOM scraping and API fetches merges through `case-enricher.js`:

```js
const merged = mergeCaseData(domData, apiData);
// API data wins over DOM data for matching fields
// merged.source = 'api+dom' when API contributed data
// merged.usedApiData = true/false
```

### Normalized case context fields

```
caseNumber, projectId, orgId, clusterName, clusterId, atlasHost,
accountId, accountName, severity, severityRank, severityLabel,
severityColor, severityWarning, status, subject, owner,
sentimentScore, sentimentLabel, createdDate, lastModified,
commentCount, description, errorStrings, pageUrl, pageTitle,
pageStatus, source, comments, jira, timeline,
refreshTargetMinutes, updatedAt
```

### Comment normalization

```js
{
  id: 'string',
  author: 'string',
  body: 'string (markdown)',
  createdAt: 'ISO8601',
  visibility: 'partner|public',
  authorKind: 'mongodb|customer'
}
```

### JIRA reference extraction

`HELP-NNNNN` references extracted from case descriptions, subjects, comments, stage payloads, next-action payloads, and page text. Found via regex patterns `HELP-\d+` and `jira.mongodb.org/browse/HELP-\d+`.

## 12. Realtime Socket Events

### Event flow

1. Connect with session cookie to `wss://support-api.ts-tools.prod.corp.mongodb.com/socket.io/?EIO=4&transport=websocket`
2. Emit `joinCase` (single) or `joinCases` (array of case numbers)
3. Listen for:
   - `Case_Updated` — case field changed (stale signal only — always REST-fetch after)
   - `notifyCaseViewers` — presence notification
   - `ALL_Updated` — broad update signal

**Critical:** Never update UI state directly from socket events. Always re-fetch from `GET /case/{number}`.

Upstream source: Salesforce platform event `/event/Case_Updated__e`.

## 13. Integration Checklist

1. Always call `normalizeBaseUrl()` on any stored TS Tools URL
2. Preserve the auth cascade order (Hub-tab cookie → bearer for non-Hub bases only)
3. Validate session with `GET /status/auth` before case-mutating calls
4. Use `GET /project/{id}/casecreation` for valid enum values — never hardcode
5. Set `partner_visibility` explicitly on every comment write (default `"partner"` = MongoDB-internal)
6. Wrap diagnostic uploads in the full S3 flow: `signedUpload` → PUT to S3 → `signedRequest` → `case/add/attachments`
7. Use `fetchWithTimeout` for all outbound HTTP from the service worker
8. Record every call attempt via `recordCallAttempt('ts-tools-api', {...})`
9. Check `looksLikeHtmlDocument()` and `didRedirectFromRequest()` before treating a 200 as success
10. Cap paginated fetches with `maxPages` to prevent runaway loops
11. Normalize case numbers to 8-digit zero-padded format before API calls

## 14. URL Construction

```js
// Case URL
function buildCaseUrl(caseNumber, baseUrl = 'https://hub.corp.mongodb.com') {
  return `${baseUrl}/case/${encodeURIComponent(caseNumber)}`;
}

// API URL
const fullUrl = `${normalizeBaseUrl(baseUrl)}${path}`;
// e.g. https://hub.corp.mongodb.com/api/case/01234567

// Extraction
extractCaseNumberFromUrl('https://hub.corp.mongodb.com/case/01234567') // '01234567'
extractAccountIdFromUrl('https://hub.corp.mongodb.com/account/5003000000D8cuI') // '5003000000D8cuI'
```

## 15. Common Pitfalls

| Pitfall | Impact | Fix |
|---|---|---|
| Silent auth redirect: 200 with HTML login page | Data appears missing | Check `isAcceptableApiResponse()` |
| Stale Hub tab with expired session | 401/403 on fetch | Call `openHubForAuth()`, then retry |
| Case number not zero-padded | API returns 404 | Run through `normalizeCaseNumber()` |
| Missing `partner_visibility` on comment | Defaults to `"partner"` (internal only) | Explicitly pass `"public"` for customer-visible comments |
| Mixed Salesforce/snake_case field names | Normalization gaps | Handle both in new code; normalizers cover existing paths |
| Socket treated as authoritative | Stale UI state | Always REST-fetch after `Case_Updated` event |
| Manual sleep/delay after 429 | Doubles backoff overhead | `fetchWithTimeout` handles 429 automatically |
| Empty body on stage/nextAction | Treated as error | Use `allowEmpty: true` in `authenticatedFetch` |

## 16. Related Skills and Files

Related references: endpoint domain reference → `tam-operations` (references/tstools-reference.md); tracked-case refresh consumer → `tam-operations` (references/case-tracker.md); MCP wrapper for these calls → `case-mcp-server-guide`; extension architecture hosting the client → `chrome-extension-expert` (references/chrome-dev.md).

| File | Purpose |
|---|---|
| `src/background/ts-tools-client.js` | Auth cascade, all fetch functions |
| `src/background/ts-tools-normalizers.js` | Response normalization helpers |
| `src/background/case-enricher.js` | DOM+API merge, severity model, JIRA extraction |
| `src/shared/fetch-with-timeout.js` | Timeout, retry, structured logging wrapper |
| `src/shared/external-calls-registry.js` | Call attempt recording for dashboard |
| `src/background/tracking-alerts.js` | Ownerless case alert logic |
| `docs/tstools-context.md` | Full 92-endpoint API reference |
| `docs/integrations-and-assumptions.md` | Auth order and behavioral assumptions |
| `docs/external-calls.md` | Audit-grade call inventory |
