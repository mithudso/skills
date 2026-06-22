Compressing inline text directly (no file path provided, so running manually per skill rules).

<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly standalone `live-hub-toolkit` skill.
> Sibling topics now reference files under hubs (`chrome-extension-expert`) — **not** standalone skills. Ignore "use X skill" / `related_skills` / SKIP pointers naming bare sibling skills; load from `references/<name>.md` in owning hub (see hub's "Cross-hub map").

---

---
name: live-hub-toolkit
description: >
  Real-time hub data discovery patterns for Chrome extensions monitoring support portals —
  MutationObserver-driven extraction, debounced publish cycles, signature-based deduplication,
  multi-context case collection, incremental backfill, and account-level sync with archival detection.
  TRIGGER: building a content script that discovers cases from a support portal in real time,
  reviewing or extending hub-extractor or live-hub-case-discovery code, implementing
  MutationObserver-based extraction with deduplication, designing incremental data sync from
  a web dashboard to extension storage, adding account-level case collection or archival detection,
  normalizing heterogeneous case records from multiple page contexts.
  SKIP: general Chrome extension development without live-data extraction (use chrome-extension-expert);
  static DOM scraping without real-time sync (use dom-scraping-resilience);
  service worker message handling without content script extraction (use chrome-extension-expert).
version: 1.2.1
last_updated: 2026-05-31
updated: 2026-05-31
category: developer
tags:
  - chrome-extension
  - content-script
  - live-data
  - mutation-observer
  - case-discovery
  - deduplication
  - support-portal
  - real-time
  - hub-extractor
triggers:
  - building a content script that discovers cases from a support portal in real time
  - reviewing or extending hub-extractor or live-hub-case-discovery code
  - implementing MutationObserver-based extraction with deduplication
  - designing incremental data sync from a web dashboard to extension storage
  - adding account-level case collection or archival detection
  - choosing between polling and push for live dashboard monitoring
  - normalizing heterogeneous case records from multiple page contexts
globs:
  - "**/hub-extractor*"
  - "**/live-hub*"
  - "**/case-discovery*"
  - "**/case-enricher*"
  - "**/content/**extractor*"
  - "**/packages/live-hub-toolkit/**"
  - "**/case-overlay*"
related_skills:
  - content-ingestion-extraction
  - chrome-extension-expert
  - tam-operations
---

# Live Hub Toolkit

Real-time data discovery for Chrome extensions on support portals. Content script watches DOM, builds context snapshot on change, deduplicates against last snapshot, forwards deltas to service worker for enrichment, storage, rendering.

**Cardinal rules:**
1. **Extract raw, normalize elsewhere.** Content script emits raw scraped fields. Normalization, severity ranking, timestamp cleanup in background helpers.
2. **Deduplicate at signature level.** Hash meaningful fields; skip publish if unchanged since last cycle.
3. **Degrade gracefully.** Missing fields → partial records, not crashes.
4. **Respect extension lifecycle.** Detect invalidated contexts, disconnect observers, stop timers on reload.

## Pipeline

```
DOM mutation -> debounced extract -> signature check -> publish to service worker
                                                           |
                                            normalize + merge + store
                                                           |
                                            delta detection + backfill + archival
```

## Architecture: content script to service worker

```
Support Portal Page
  |
  +-- content script (hub-extractor.js)
  |     - MutationObserver on document root
  |     - history.pushState / replaceState patches
  |     - hashchange + popstate listeners
  |     - URL polling fallback (750ms interval)
  |     - debounced extractAndPublish() on every trigger
  |     - signature-based dedup before sending
  |     - MCA_UPSERT_CASE_CONTEXT message to service worker
  |
  +-- service worker (background)
        - receives MCA_UPSERT_CASE_CONTEXT
        - normalizes via case-enricher / live-hub-case-discovery
        - merges DOM data with API data
        - stores in chrome.storage.session (disposable cache)
        - triggers tracking prompts, alerts, dashboard refresh
```

### Key message types

| Message | Direction | Purpose |
|---------|-----------|---------|
| `MCA_UPSERT_CASE_CONTEXT` | content → background | Push extracted page context |
| `MCA_GET_CASE_CONTEXT` | background → content | Pull current snapshot on demand |
| `MCA_EVALUATE_TRACKING_PROMPT` | content → background | Ask if prompt operator to track |
| `MCA_SYNC_FAVORITED_ACCOUNTS` | content → background | Push portfolio favorites list |
| `MCA_REALTIME_BRIDGE_PAGE_READY` | content → background | Signal page bootstrap complete |

---

## MutationObserver-driven extraction

```js
// Observe full document subtree with targeted attribute filters.
// attributeFilter limits noise — only fire when test IDs, labels,
// links, or titles change, not on every style recalculation.
const observer = new MutationObserver(() => scheduleExtract(250));
observer.observe(document.documentElement || document.body, {
  childList: true, subtree: true, characterData: true,
  attributes: true,
  attributeFilter: ['data-testid', 'aria-label', 'href', 'title']
});
```

### Debounce tiers

| Event | Debounce | Rationale |
|-------|----------|-----------|
| MutationObserver callback | 250ms | Collapse SPA render batches |
| `hashchange` / `popstate` | 50ms | Navigation needs fast response |
| `history.pushState` / `replaceState` | 100ms | SPA route change, slight settle |
| `visibilitychange` (tab refocus) | 50ms | User returned, show fresh data |
| URL poll (fallback) | 750ms interval | Catches edge-case navigations |

```js
let extractTimer = null;
function scheduleExtract(delay = 150) {
  if (!extensionContextAlive) return;
  clearTimeout(extractTimer);
  extractTimer = window.setTimeout(() => extractAndPublish(), delay);
}
```

---

## Signature-based deduplication

```js
function buildContext() {
  const context = {
    source: 'hub-extractor',
    extracted_at: new Date().toISOString(),
    page_url: location.href,
    case_number: extractCaseNumber(),
    // ... all extracted fields
    account_cases: extractAccountCases(),
  };

  // extracted_at intentionally excluded — it changes every cycle
  context.signature = JSON.stringify({
    case_number: context.case_number,
    account_id: context.account_id,
    subject: context.subject,
    status: context.status,
    owner: context.owner,
    last_modified: context.last_modified,
    account_cases: context.account_cases,
    visible_errors: context.visible_errors,
    page_url: context.page_url
  });
  return context;
}

function extractAndPublish() {
  const context = buildContext();
  if (context.signature === lastSignature && lastUrl === location.href) return;
  lastSignature = context.signature;
  lastUrl = location.href;
  if (document.hidden) return;  // re-extract on visibilitychange
  publishContext(context);
}
```

**Include in signature:** case number, status, severity, owner, subject, account ID, last modified, visible errors, account cases list, page URL.
**Exclude:** `extracted_at`, `page_title`, transient UI state.

---

## Multi-source case record normalization

### Canonical case record shape

```js
{
  caseNumber: '01581027',     // 8-digit zero-padded
  subject: 'Atlas Search degraded',
  status: 'Open',
  severity: 'Sev2',
  owner: 'Mitch Hudson',
  accountId: 'acct-123',
  accountName: 'Acme Corp',
  lastModified: '2026-05-26T10:30:00.000Z'
}
```

### Field alias resolution

```js
function normalizeLiveCaseRecord(record = {}, { fallbackAccountId = '', fallbackAccountName = '' } = {}) {
  return {
    caseNumber: normalizeCaseNumber(
      record.caseNumber || record.case_number || record.number || ''
    ),
    subject: normalizeWhitespace(record.subject || record.title || ''),
    status:   normalizeWhitespace(record.status || ''),
    severity: normalizeWhitespace(record.severity || ''),
    owner:    normalizeWhitespace(record.owner || record.assignee || record.caseOwner || ''),
    accountId:   normalizeAccountId(record.accountId || record.account_id || fallbackAccountId),
    accountName: normalizeWhitespace(record.accountName || record.account_name || fallbackAccountName),
    lastModified: normalizeWhitespace(record.lastModified || record.last_modified || '')
  };
}
```

### Case number normalization

Always normalize to 8-digit zero-padded:

```js
const CASE_NUMBER_PATTERN = /\b0?\d{7,8}\b/;
function normalizeCaseNumber(value = '') {
  const match = String(value || '').match(CASE_NUMBER_PATTERN);
  return match ? match[0].padStart(8, '0').slice(-8) : '';
}
// '1581027'  -> '01581027'
// '01581027' -> '01581027'
```

---

## Collection merging and deduplication

### Record scoring

Score record by non-empty field count. Richer record wins:

```js
function scoreCaseRecord(record = {}) {
  return ['caseNumber','subject','status','severity','owner','accountId','accountName','lastModified']
    .filter(f => record[f]).length;
}
```

### Merge algorithm

```js
function mergeCaseCollections(...collections) {
  const deduped = new Map();
  for (const collection of collections) {
    for (const rawRecord of Array.isArray(collection) ? collection : []) {
      const record = normalizeLiveCaseRecord(rawRecord);
      if (!record.caseNumber) continue;
      const existing = deduped.get(record.caseNumber);
      if (!existing || scoreCaseRecord(record) >= scoreCaseRecord(existing)) {
        deduped.set(record.caseNumber, record);
      }
    }
  }
  return Array.from(deduped.values()).sort((a, b) => a.caseNumber.localeCompare(b.caseNumber));
}
```

Properties: **idempotent** (same collection twice → same result), **monotonic** (richer replaces sparser), **deterministic sort** (stable output order).

---

## Multi-context account case discovery

```js
function collectAccountCasesFromContexts(contexts = [], { accountId, accountName } = {}) {
  const discovered = [];
  let observedAccountPage = false;

  for (const context of contexts) {
    const isAccountPage = accountId
      && context.account_id === accountId
      && /\/accounts?\//i.test(context.page_url || '');

    if (isAccountPage) {
      observedAccountPage = true;
      discovered.push(...(context.account_cases || []).map(record =>
        normalizeLiveCaseRecord(record, { fallbackAccountId: accountId, fallbackAccountName: accountName })
      ));
    }

    // Also check if this is a case detail page for this account
    const directCase = normalizeLiveCaseRecord({
      caseNumber: context.case_number, subject: context.subject,
      status: context.status, severity: context.severity,
      owner: context.owner, accountId: context.account_id,
    }, { fallbackAccountId: accountId, fallbackAccountName: accountName });

    if (directCase.caseNumber && matchesTrackedAccount(directCase, accountId, accountName)) {
      discovered.push(directCase);
    }
  }

  return { cases: mergeCaseCollections(discovered), observedAccountPage };
}
```

`observedAccountPage` flag critical: determines if missing cases archive or merely not-yet-seen.

---

## Incremental backfill and delta detection

```js
const CORE_TRACKED_FIELDS = ['subject', 'status', 'severity'];

// Find which core fields are missing
function getIncompleteTrackedCaseFields(record = {}) {
  const normalized = normalizeLiveCaseRecord(record);
  return CORE_TRACKED_FIELDS.filter(field => !normalized[field]);
}

// Compute what actually changed before applying
function getTrackedCaseFieldDelta(existing = {}, candidate = {}) {
  const existingNorm = normalizeLiveCaseRecord(existing);
  const candidateNorm = normalizeLiveCaseRecord(candidate, {
    fallbackAccountId: existingNorm.accountId,
    fallbackAccountName: existingNorm.accountName
  });
  return ['subject','status','severity','owner','accountId','accountName']
    .filter(field => candidateNorm[field] && candidateNorm[field] !== existingNorm[field])
    .map(field => ({ field, from: existingNorm[field] || '', to: candidateNorm[field] || '' }));
}

// Fill missing fields without overwriting existing values
function applyTrackedCaseBackfill(base = {}, candidate = {}) {
  const baseNorm = normalizeLiveCaseRecord(base);
  const candidateNorm = normalizeLiveCaseRecord(candidate, {
    fallbackAccountId: baseNorm.accountId, fallbackAccountName: baseNorm.accountName
  });
  return {
    ...base,
    accountId:   candidateNorm.accountId   || base.accountId   || '',
    accountName: candidateNorm.accountName || base.accountName || '',
    subject:     candidateNorm.subject     || base.subject     || '',
    status:      candidateNorm.status      || base.status      || '',
    severity:    candidateNorm.severity    || base.severity    || '',
    owner:       candidateNorm.owner       || base.owner       || ''
  };
}
```

---

## Account-level sync and archival detection

```js
function shouldReplaceTrackedAccountCases({ existingCaseNumbers = [], mergedCases = [], observedAccountPage = false } = {}) {
  if (existingCaseNumbers.length === 0) return true;  // always accept new data
  if (mergedCases.length > 0) return true;             // new cases discovered
  return observedAccountPage;                          // saw authoritative list
}

function findCasesToArchiveFromAccountSync({ existingCaseNumbers = [], mergedCases = [], canReplace = false } = {}) {
  if (!canReplace) return [];
  const activeCaseNumbers = new Set(
    mergedCases.map(r => normalizeCaseNumber(r.caseNumber || '')).filter(Boolean)
  );
  return existingCaseNumbers
    .map(n => normalizeCaseNumber(n))
    .filter(n => n && !activeCaseNumbers.has(n))
    .sort();
}
```

Only archive when `observedAccountPage === true` — never archive cases extension hasn't seen yet.

---

## SPA navigation detection

Support portals typically SPAs. Layer all five mechanisms:

```js
// 1. History API patching
const _origPushState = history.pushState.bind(history);
history.pushState = function(...args) { _origPushState(...args); scheduleExtract(100); };
history.replaceState = function(...args) { _origReplaceState(...args); scheduleExtract(100); };

// 2. Event listeners
window.addEventListener('hashchange', () => scheduleExtract(50), { passive: true });
window.addEventListener('popstate',   () => scheduleExtract(50), { passive: true });
document.addEventListener('visibilitychange', () => { if (!document.hidden) scheduleExtract(50); });

// 3. URL polling fallback (catches frameworks that bypass history API)
let locationPollTimer = window.setInterval(() => {
  if (location.href !== lastUrl) scheduleExtract(50);
}, 750);
```

---

## Extension lifecycle safety

```js
let extensionContextAlive = true;

function isExtensionContextInvalidated(error) {
  return /Extension context invalidated/i.test(String(error?.message || error || ''));
}

function deactivateExtractor() {
  if (!extensionContextAlive) return;
  extensionContextAlive = false;
  clearTimeout(extractTimer);
  if (locationPollTimer) { window.clearInterval(locationPollTimer); locationPollTimer = null; }
  observer?.disconnect?.();
}

async function sendRuntimeMessage(message) {
  if (!extensionContextAlive) return null;
  try {
    return await chrome.runtime?.sendMessage?.(message);
  } catch (error) {
    if (isExtensionContextInvalidated(error)) { deactivateExtractor(); return null; }
    return null;
  }
}

// Boot guard — prevent double-initialization
if (window.__mdbCaseAssistantHubExtractorBooted) return;
window.__mdbCaseAssistantHubExtractorBooted = true;
```

---

## DOM extraction helpers

Core utilities: `normalizeWhitespace()`, `isElementVisible()`, `readElementValue()`, `extractLabeledValue()` (inline "Label: Value" + sibling strategies). Also `extractAccountCases()` (walks case links/rows) and `matchesTrackedAccount()` (ID-first, name-fallback).

Full implementations: [references/dom-extraction-helpers.md](./references/dom-extraction-helpers.md)

---

## Anti-patterns

| Anti-pattern | Why it fails | Fix |
|---|---|---|
| Publishing on every MutationObserver callback | Floods service worker with identical data | Signature-based dedup with debounced publish |
| Including `extracted_at` in signature | Signature changes every cycle, defeating dedup | Exclude volatile timestamps from signature |
| Normalizing in content script | Content script bloated, hard to test | Extract raw, normalize in background helpers |
| Single extraction strategy per field | Breaks when portal redesigns | Selector fallback chains + labeled value search |
| No boot guard | Content script runs twice after extension reload | Check `window.__booted` flag before init |
| Ignoring `document.hidden` | Publishes stale data for background tabs | Skip publish when hidden, re-extract on focus |
| Not disconnecting observers on invalidation | Console errors, leaked timers | `deactivateExtractor()` on first error |
| Trusting single tab for account case list | Misses cases open in other tabs | `collectAccountCasesFromContexts` across all |
| Archiving without observing account page | Incorrectly archives cases not yet seen | Only archive when `observedAccountPage === true` |

---

## Testing

`live-hub-toolkit` pure unit testing without browser APIs — all functions take/return plain objects.

```js
import test from 'node:test';
import assert from 'node:assert/strict';
import { normalizeLiveCaseRecord, mergeCaseCollections } from './index.js';

test('normalizes mixed field shapes', () => {
  assert.deepEqual(normalizeLiveCaseRecord({
    case_number: '1581027', title: ' Atlas Search degraded ', account_id: 'acct-1'
  }), {
    caseNumber: '01581027', subject: 'Atlas Search degraded',
    status: '', severity: '', owner: '', accountId: 'acct-1', accountName: '', lastModified: ''
  });
});

test('merge keeps richer duplicate', () => {
  const result = mergeCaseCollections(
    [{ caseNumber: '01581027', subject: 'Sparse' }],
    [{ caseNumber: '01581027', subject: 'Rich', status: 'Open', severity: 'Sev1', owner: 'Mitch' }]
  );
  assert.equal(result[0].subject, 'Rich');
  assert.equal(result[0].owner, 'Mitch');
});
```

---

## Quick reference

```
Extracting from a portal page?
  → MutationObserver (childList + subtree + attributeFilter)
  → Debounce: 250ms for DOM, 50ms for navigation
  → Signature from meaningful fields only; skip if unchanged

Normalizing heterogeneous records?
  → normalizeLiveCaseRecord() with field alias resolution
  → normalizeCaseNumber() for 8-digit zero-padded
  → normalizeWhitespace() for NBSP and multi-space cleanup

Merging from multiple tabs/contexts?
  → mergeCaseCollections() with score-based dedup
  → collectAccountCasesFromContexts() for cross-tab discovery
  → Track observedAccountPage flag for archival decisions

Tracking incremental changes?
  → getIncompleteTrackedCaseFields() to find gaps
  → getTrackedCaseFieldDelta() to compute what changed
  → applyTrackedCaseBackfill() to fill gaps without overwriting

Deciding what to archive?
  → shouldReplaceTrackedAccountCases() for replace decision
  → findCasesToArchiveFromAccountSync() for stale case detection
  → Only archive when the authoritative account page was observed
```

---

## Source files

- `packages/live-hub-toolkit/src/index.js` — pure normalization and merge functions
- `src/background/live-hub-case-discovery.js` — re-export bridge for service worker
- `src/content/hub-extractor.js` — production content script, full extraction pipeline
- `src/background/case-enricher.js` — background normalization + severity ranking