<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `salesforce-scraping-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: salesforce-scraping-patterns
version: 1.1.0
category: developer
tags: [salesforce, scraping, chrome-extension, lightning, lwc, content-script, shadow-dom, rest-api, soql, session-cookie]
description: >-
  Salesforce data extraction patterns for Chrome extensions — Lightning DOM traversal (LWC synthetic vs
  native shadow roots), Classic UI selectors, API-first alternatives (REST/SOQL/Composite), session-based
  auth via sid cookie, OAuth Connected App flow, CSP handling, SPA navigation detection, and ethical
  scraping guidelines.
  TRIGGER: building a content script that extracts data from Salesforce pages, handling LWC shadow DOM
  in a content script, extracting case/account/opportunity data from Lightning UI, leveraging Salesforce
  session cookies for API calls from an extension, or configuring host_permissions for Salesforce domains.
  SKIP: Salesforce platform development (Apex, LWC authoring, Flow, CI/CD) — use salesforce-developer-expert;
  non-extension Salesforce integrations; Salesforce admin tasks.
triggers:
  - "Extract case data from Salesforce Lightning page with a content script"
  - "Shadow DOM traversal for LWC components in Chrome extension"
  - "Read Salesforce session cookie (sid) from extension service worker"
  - "Configure host_permissions for *.salesforce.com in MV3 manifest"
  - "Lightning SPA navigation detection in content script"
  - "SOQL query from Chrome extension using session cookie"
  - "kagekiri shadow-piercing querySelector for Salesforce"
  - "Detect Lightning vs Classic Salesforce UI in content script"
  - "Composite API batch queries from Chrome extension"
  - "Handle pushState navigation in Salesforce Lightning"
related_skills:
  - salesforce-developer-expert
  - chrome-dev
  - dom-scraping-resilience
  - extension-message-bridge
---

# Salesforce Scraping Patterns

Extracting data from Salesforce UI is uniquely challenging. Lightning Experience uses LWC with shadow DOM encapsulation, progressive loading, and framework-managed DOM. Classic uses server-rendered Visualforce with different selectors.

**Cardinal rule: prefer the Salesforce REST API over DOM scraping.** DOM scraping breaks on every Salesforce release (3x/year) and violates stable interfaces. Use scraping only when: (1) data is unavailable via API, (2) you need real-time UI-context awareness ("what record is the user viewing right now"), or (3) you are augmenting the Salesforce UI with an overlay that reacts to page context.

**When not to use:** Apex development, LWC authoring, Flow Builder, Salesforce DX/CI/CD — use `salesforce-developer-expert` for those.

## URL parsing (do this first — free context)

```js
function parseSalesforceUrl(url = location.href) {
  // Lightning: /lightning/r/<ObjectType>/<RecordId>/<view>
  const lightningMatch = url.match(/\/lightning\/r\/([A-Za-z_]+)\/([a-zA-Z0-9]{15,18})(?:\/(\w+))?/);
  if (lightningMatch) return { objectType: lightningMatch[1], recordId: lightningMatch[2], view: lightningMatch[3] || 'view' };

  // Lightning list
  const listMatch = url.match(/\/lightning\/o\/([A-Za-z_]+)\/list/);
  if (listMatch) return { objectType: listMatch[1], recordId: null, view: 'list' };

  // Classic: /<15-or-18-char-id>
  const classicMatch = url.match(/salesforce\.com\/([a-zA-Z0-9]{15,18})(?:\?|$|\/)/);
  if (classicMatch) return { objectType: inferObjectFromPrefix(classicMatch[1]), recordId: classicMatch[1], view: 'view' };

  return { objectType: null, recordId: null, view: null };
}

// Common 3-char ID prefixes
const SF_PREFIXES = { '500': 'Case', '001': 'Account', '006': 'Opportunity', '003': 'Contact', '00Q': 'Lead', '005': 'User', '570': 'EmailMessage' };
function inferObjectFromPrefix(id) { return SF_PREFIXES[id.substring(0, 3)] || null; }
```

## Lightning DOM architecture

### Synthetic vs native shadow

| Mode | Content script access | Selector behavior |
|------|----------------------|-------------------|
| Synthetic shadow (legacy) | Elements often accessible via `document.querySelector` — polyfill doesn't enforce true encapsulation | Works inconsistently |
| Native shadow (current migration) | Elements hidden behind real `shadowRoot` boundaries | Must traverse `.shadowRoot` explicitly |
| Light DOM (opt-in per component) | Elements directly in document tree | Standard selectors work |

Orgs may have a mix of synthetic and native shadow during migration. Your code must handle both.

### Stable vs unstable selectors

| Stable (prefer) | Unstable (avoid) |
|-----------------|-----------------|
| Custom element tag names: `lightning-formatted-text`, `records-record-layout-item` | Generated class names: `.slds-form-element__label` |
| `data-field-id`, `data-target-selection-name` attributes | Positional: `:nth-child(3)` |
| `data-aura-rendered-by` (Aura components) | Hashed class names: `.cCase_Detail` |
| `[field-label]` on record layout items | Generated `id` attributes |

## Shadow DOM traversal

```js
// Full-tree traversal (slower — use when structure is unknown)
function queryLightningShadow(selector, root = document, maxDepth = 12) {
  const light = root.querySelector(selector);
  if (light) return light;
  if (maxDepth <= 0) return null;
  for (const el of root.querySelectorAll('*')) {
    const shadow = el.shadowRoot;
    if (shadow) {
      const found = queryLightningShadow(selector, shadow, maxDepth - 1);
      if (found) return found;
    }
  }
  return null;
}

function queryAllLightningShadow(selector, root = document, maxDepth = 12) {
  const results = [...root.querySelectorAll(selector)];
  if (maxDepth <= 0) return results;
  for (const el of root.querySelectorAll('*')) {
    if (el.shadowRoot) results.push(...queryAllLightningShadow(selector, el.shadowRoot, maxDepth - 1));
  }
  return results;
}

// Targeted path traversal (faster — use when component nesting is known)
function walkShadowPath(path, finalSelector) {
  let current = document;
  for (const tag of path) {
    const el = current.querySelector(tag);
    if (!el) return null;
    current = el.shadowRoot || el;
  }
  return current.querySelector(finalSelector);
}

// Example: extract case number from Lightning record page
const caseNumber = walkShadowPath(
  ['one-record-home-flexipage2', 'records-record-layout-section', 'records-record-layout-item[data-field-id="RecordCaseNumberField"]'],
  'lightning-formatted-text'
)?.textContent?.trim();
```

### kagekiri (shadow-piercing querySelector)

```js
// npm install kagekiri — ~8KB gzipped
import { querySelector, querySelectorAll } from 'kagekiri';
const caseField = querySelector('lightning-formatted-text[data-field-id="CaseNumber"]');
```

Only works with **open** shadow roots. Use for unknown structures; use `walkShadowPath` for known ones.

## Manifest configuration

```jsonc
{
  "manifest_version": 3,
  "host_permissions": [
    "https://*.salesforce.com/*",
    "https://*.force.com/*",
    "https://*.my.salesforce.com/*",
    "https://*.visualforce.com/*"
  ],
  "content_scripts": [
    {
      "matches": ["https://*.lightning.force.com/*", "https://*.salesforce.com/*"],
      "js": ["src/content/sf-extractor.js"],
      "run_at": "document_idle",
      "world": "ISOLATED"
    },
    {
      "matches": ["https://*.lightning.force.com/*"],
      "js": ["src/content/sf-navigation-hook.js"],
      "run_at": "document_start",
      "world": "MAIN"
    }
  ],
  "permissions": ["cookies", "activeTab", "storage"]
}
```

Content scripts run in an isolated world and are **not** subject to the page's CSP. Injecting into the MAIN world IS subject to the page CSP. Fetching Salesforce APIs from the service worker bypasses page CSP entirely — the recommended approach.

## API-first: REST/SOQL from the service worker

### Why API beats DOM scraping

| Concern | DOM scraping | REST API |
|---------|-------------|---------|
| Stability | Breaks every SF release (3x/year) | Versioned, backward-compatible |
| Coverage | Only rendered fields | All fields, related records |
| Performance | Requires page load + shadow traversal | Direct HTTP, returns JSON |
| Auth | Implicit (user session) | Session ID or OAuth token |

### Session cookie extraction

```js
// service-worker.js
async function getSalesforceSession(instanceUrl) {
  const cookies = await chrome.cookies.getAll({
    domain: new URL(instanceUrl).hostname,
    name: 'sid',
  });
  return cookies.length ? { sessionId: cookies[0].value, instanceUrl } : null;
}
```

### SOQL query helper

```js
function validateSfId(id) {
  if (typeof id !== 'string' || !/^[a-zA-Z0-9]{15,18}$/.test(id)) throw new Error(`Invalid SF ID: ${id}`);
  return id;
}

async function sfQuery(instanceUrl, sessionId, soql) {
  const apiBase = instanceUrl.replace('.lightning.force.com', '.my.salesforce.com');
  const url = `${apiBase}/services/data/v67.0/query?q=${encodeURIComponent(soql)}`;
  const resp = await fetch(url, { headers: { 'Authorization': `Bearer ${sessionId}` } });
  if (!resp.ok) throw new Error(`SF query failed: ${resp.status}`);
  return resp.json();
}

// Always validate IDs before SOQL interpolation
async function fetchCaseDetails(instanceUrl, sessionId, caseId) {
  const soql = `SELECT Id, CaseNumber, Subject, Status, Priority FROM Case WHERE Id = '${validateSfId(caseId)}' LIMIT 1`;
  const result = await sfQuery(instanceUrl, sessionId, soql);
  return result.records?.[0] || null;
}

// Paginate through all records
async function sfQueryAll(instanceUrl, sessionId, soql) {
  const apiBase = instanceUrl.replace('.lightning.force.com', '.my.salesforce.com');
  const allRecords = [];
  let result = await sfQuery(instanceUrl, sessionId, soql);
  allRecords.push(...result.records);
  while (!result.done && result.nextRecordsUrl) {
    const resp = await fetch(`${apiBase}${result.nextRecordsUrl}`, { headers: { 'Authorization': `Bearer ${sessionId}` } });
    result = await resp.json();
    allRecords.push(...result.records);
  }
  return allRecords;
}
```

### Composite API (batch multiple queries in one round trip)

```js
async function sfCompositeQuery(instanceUrl, sessionId, queries) {
  const apiBase = instanceUrl.replace('.lightning.force.com', '.my.salesforce.com');
  const resp = await fetch(`${apiBase}/services/data/v67.0/composite`, {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${sessionId}`, 'Content-Type': 'application/json' },
    body: JSON.stringify({
      compositeRequest: queries.map((q, i) => ({
        method: 'GET',
        url: `/services/data/v67.0/query?q=${encodeURIComponent(q.soql)}`,
        referenceId: q.refId || `query_${i}`,
      })),
    }),
  });
  if (!resp.ok) throw new Error(`SF composite failed: ${resp.status}`);
  return resp.json();
}
```

## DOM extraction: Lightning record pages

```js
function extractRecordFields() {
  const fields = {};
  // Strategy 1: records-record-layout-item (most reliable)
  for (const item of queryAllLightningShadow('records-record-layout-item')) {
    const shadow = item.shadowRoot || item;
    const label = shadow.querySelector('.slds-form-element__label, [data-output-element-id="output-field-label"]');
    const value = shadow.querySelector('lightning-formatted-text, lightning-formatted-url, lightning-formatted-email');
    if (label && value) fields[label.textContent?.trim()] = value.textContent?.trim() || null;
  }
  // Strategy 2: data-target-selection-name (newer builds)
  for (const field of queryAllLightningShadow('[data-target-selection-name^="sfdc:RecordField"]')) {
    const apiName = field.getAttribute('data-target-selection-name')?.split('.').pop();
    const valueEl = (field.shadowRoot || field).querySelector('lightning-formatted-text, [data-output-element-id]');
    if (apiName && valueEl) fields[apiName] = valueEl.textContent?.trim() || null;
  }
  return fields;
}
```

## DOM extraction: Classic UI

```js
function extractClassicRecord() {
  const fields = {};
  for (const labelCell of document.querySelectorAll('td.labelCol, td.labelCol.last')) {
    const label = labelCell.textContent?.trim();
    const valueCell = labelCell.nextElementSibling;
    if (label && valueCell?.classList.contains('dataCol')) fields[label] = valueCell.textContent?.trim() || null;
  }
  return fields;
}

function detectSalesforceUI() {
  if (location.hostname.includes('.lightning.force.com') || document.querySelector('one-app-nav-bar')) return 'lightning';
  if (document.querySelector('#bodyTable, .bPageBlock')) return 'classic';
  return 'unknown';
}
```

## SPA navigation detection

Lightning is an SPA — content scripts do NOT re-run on `pushState` navigation.

```js
// sf-navigation-hook.js — MAIN world, document_start
(function () {
  const notify = (method) => window.dispatchEvent(new CustomEvent('sf:navigation', { detail: { method, url: location.href } }));
  const origPush = history.pushState.bind(history);
  const origReplace = history.replaceState.bind(history);
  history.pushState = (...a) => { origPush(...a); notify('pushState'); };
  history.replaceState = (...a) => { origReplace(...a); notify('replaceState'); };
  window.addEventListener('popstate', () => notify('popstate'));
})();
```

```js
// sf-extractor.js — ISOLATED world
let extractionAbort = new AbortController();

window.addEventListener('sf:navigation', async (e) => {
  extractionAbort.abort();
  extractionAbort = new AbortController();
  const { recordId } = parseSalesforceUrl(e.detail.url);
  if (!recordId) return;
  await waitForLightningReady(extractionAbort.signal);
  performExtraction();
});

async function waitForLightningReady(signal) {
  return new Promise((resolve, reject) => {
    if (signal?.aborted) return reject(new DOMException('Aborted', 'AbortError'));
    let attempts = 0;
    const timer = setInterval(() => {
      if (queryLightningShadow('records-record-layout-section') || ++attempts >= 20) { clearInterval(timer); resolve(); }
    }, 250);
    signal?.addEventListener('abort', () => { clearInterval(timer); reject(new DOMException('Aborted', 'AbortError')); });
  });
}
```

## Authentication

**Priority order:**
1. Check `sid` cookie (zero-config for the user)
2. Use stored OAuth token (cross-domain scenarios)
3. Prompt for OAuth login via `chrome.identity.launchWebAuthFlow`

**Never** store or transmit the `sid` cookie outside the extension — it is a full session token equivalent to the user's credentials. Keep it in-memory or in `chrome.storage.session` only.

## Anti-patterns

| Anti-pattern | Fix |
|-------------|-----|
| Hardcoded CSS class selectors | Use tag names, `data-` attributes, `data-target-selection-name` |
| `document.querySelector` without shadow traversal | `queryLightningShadow()` or `kagekiri` |
| Scraping when API data is available | REST API + SOQL from service worker |
| Storing `sid` in `chrome.storage.local` | Keep in `chrome.storage.session` (memory-only) |
| Polling with `setInterval` for element appearance | MutationObserver + `waitForElement` |
| Full shadow-tree traversal on every mutation | Debounce 300–500ms + `walkShadowPath` for known structures |
| SOQL string interpolation without ID validation | Always call `validateSfId()` before interpolating |

## Troubleshooting

**Shadow traversal returns null:**
- Check DevTools — is the component using native shadow (`#shadow-root (open)`) or synthetic (no real boundary)?
- Lightning loads progressively — add `waitForLightningReady()` before extraction.
- Increase `maxDepth` to 15–20 for flexipage-heavy orgs.

**API returns 401:** `sid` cookie may have expired (default session timeout: 2 hours). Try reading cookies from both `*.lightning.force.com` and `*.my.salesforce.com` — they may differ.

**Extraction works on first load but breaks on navigation:** Lightning is an SPA. Implement the `pushState` hook in MAIN world; re-extract after navigation + render wait.

**SOQL returns no records despite valid ID:** Salesforce IDs are case-sensitive in 15-char form. Use 18-char IDs from Lightning URLs. Check for `INSUFFICIENT_ACCESS` in the error response.
