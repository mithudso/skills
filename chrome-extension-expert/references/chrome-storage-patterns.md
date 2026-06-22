<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly standalone `chrome-storage-patterns` skill.
> Sibling topics now reference files under hubs (`chrome-extension-expert`) — **not** standalone
> skills. Ignore "use the X skill" / `related_skills` / SKIP pointers naming bare sibling
> skills; load that topic's `references/<name>.md` from owning hub (see hub's "Cross-hub map").

---

---
name: chrome-storage-patterns
description: >
  Chrome extension storage patterns — chrome.storage.local vs session vs sync,
  quota management (QUOTA_BYTES, QUOTA_BYTES_PER_ITEM), eviction, onChanged
  listeners, cross-context sync, batch get/set operations, chunking large data,
  and MV3 service worker persistence strategies.
  TRIGGER: designing storage architecture for Chrome extensions, debugging
  "QUOTA_BYTES quota exceeded" errors, fixing data lost after service worker
  restart, choosing between storage areas, migrating from localStorage to
  chrome.storage, cross-context data synchronization via onChanged.
  SKIP: IndexedDB patterns without chrome.storage (use indexeddb-patterns),
  server-side storage or databases (not extension-specific), Web Storage API
  (localStorage/sessionStorage) in non-extension web pages.
version: 1.1.0
category: developer
tags: [chrome-extension, storage, mv3, service-worker, persistence, quota, onChanged, batch-operations, chunking, IndexedDB, setAccessLevel, unlimitedStorage]
related_skills: [mv3-service-worker-expert, chrome-offscreen-documents, chrome-identity-oauth, indexeddb-patterns]
updated: 2026-05-29
---

# Chrome Storage Patterns

## Overview

`chrome.storage` = canonical persistence layer for Chrome MV3 extensions. Unlike `localStorage` (synchronous, blocking, unavailable in service workers), `chrome.storage` provides async bulk read/write accessible from every extension context: service workers, content scripts, popups, options pages, extension-owned tabs.

MV3 ephemeral service workers make storage architecture critical. Service workers terminate after 30s idle, lose all in-memory state on restart. Any state that must survive restart needs writing to `chrome.storage.local` or `chrome.storage.session` proactively.

## Storage Areas Comparison

| Area | Persistence | Quota | Access | Use case |
|------|-------------|-------|--------|----------|
| `local` | Survives browser restarts; cleared on uninstall | 10 MB (5 MB in Chrome ≤113) | All extension contexts | Durable settings, tracking state, vault envelopes |
| `session` | In-memory only; cleared on browser restart or extension reload | 10 MB (1 MB in Chrome ≤111) | Trusted contexts by default | Auth tokens, unlock secrets, disposable cache |
| `sync` | Synced across signed-in Chrome instances | 100 KB total, 8 KB per item | All extension contexts | User preferences that follow the user |
| `managed` | Read-only, set by enterprise policy | N/A | All extension contexts | Admin-enforced configuration |

## Decision Tree: Which Storage Area?

```
Must the data survive browser restarts?
├── YES → Must it sync across devices?
│   ├── YES → Is it under 8 KB per item and 100 KB total?
│   │   ├── YES → chrome.storage.sync
│   │   └── NO  → chrome.storage.local
│   └── NO  → Is it over 10 MB?
│       ├── YES → chrome.storage.local + unlimitedStorage permission
│       │         (or IndexedDB for very large datasets)
│       └── NO  → chrome.storage.local
└── NO  → Is the data sensitive (tokens, secrets, keys)?
    ├── YES → chrome.storage.session (never written to disk)
    └── NO  → Does it need to survive SW restarts within a session?
        ├── YES → chrome.storage.session
        └── NO  → In-memory JS variable (cheapest; lost on SW kill)
```

**Quick rules:**
- Durable config/state → `local`
- Secrets, auth tokens → `session`
- User prefs that roam → `sync`
- Enterprise policy → `managed` (read-only)
- Hot cache, disposable → `session` (fast, no disk I/O)
- Massive blobs (>10 MB) → IndexedDB (no cross-context `onChanged`, but no quota cap with `unlimitedStorage`)

## Core API Usage

```javascript
// Write
await chrome.storage.local.set({
  trackingState: { cases: [], lastRefresh: Date.now() },
  options: { theme: 'dark', refreshInterval: 300000 }
});

// Read
const { trackingState, options } = await chrome.storage.local.get([
  'trackingState', 'options'
]);

// Remove
await chrome.storage.local.remove(['obsoleteKey']);

// Clear (destructive)
await chrome.storage.local.clear();
```

### Enabling Content Script Access to Session Storage

Default: `chrome.storage.session` only accessible from trusted contexts (service worker, extension pages). To allow content scripts:

```javascript
// Call once in service worker at startup
chrome.storage.session.setAccessLevel({
  accessLevel: 'TRUSTED_AND_UNTRUSTED_CONTEXTS'
});
```

## Practical Patterns

### Pattern 1: Service Worker State Recovery

Most critical MV3 pattern. Restore state on every SW activation:

```javascript
// service-worker.js
let cachedState = null;

async function getState() {
  if (cachedState) return cachedState;
  const { appState } = await chrome.storage.local.get('appState');
  cachedState = appState || getDefaultState();
  return cachedState;
}

async function setState(patch) {
  cachedState = { ...cachedState, ...patch };
  await chrome.storage.local.set({ appState: cachedState });
}

// Register listeners synchronously at top level
chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
  if (msg.type === 'GET_STATE') {
    getState().then(sendResponse);
    return true;
  }
  if (msg.type === 'UPDATE_STATE') {
    setState(msg.payload).then(() => sendResponse({ ok: true }));
    return true;
  }
});
```

### Pattern 2: Dual-Layer Cache (Session + Local)

Session storage = fast read cache; local = durable backup:

```javascript
async function getCaseData(caseId) {
  const sessionKey = `case_${caseId}`;

  // Layer 1: session cache (fast, survives SW restart within session)
  const sessionResult = await chrome.storage.session.get(sessionKey);
  if (sessionResult[sessionKey]) return sessionResult[sessionKey];

  // Layer 2: durable local storage (survives browser restart)
  const localResult = await chrome.storage.local.get(sessionKey);
  if (localResult[sessionKey]) {
    await chrome.storage.session.set({ [sessionKey]: localResult[sessionKey] });
    return localResult[sessionKey];
  }

  // Layer 3: network fetch, populate both layers
  const data = await fetchCaseFromAPI(caseId);
  await Promise.all([
    chrome.storage.session.set({ [sessionKey]: data }),
    chrome.storage.local.set({ [sessionKey]: data })
  ]);
  return data;
}
```

### Pattern 3: Batch Read/Write

Each `get()`/`set()` = IPC between extension process and browser process. Always batch:

```javascript
// BAD: 3 round-trips
const a = await chrome.storage.local.get('keyA');
const b = await chrome.storage.local.get('keyB');

// GOOD: 1 round-trip
const { keyA, keyB } = await chrome.storage.local.get(['keyA', 'keyB']);

// BAD: 3 writes
await chrome.storage.local.set({ keyA: valueA });
await chrome.storage.local.set({ keyB: valueB });

// GOOD: 1 atomic write
await chrome.storage.local.set({ keyA: valueA, keyB: valueB });
```

### Pattern 4: Cross-Context Synchronization with onChanged

`onChanged` fires in every extension context when storage mutates:

```javascript
// In service worker, popup, options page, or content script:
chrome.storage.onChanged.addListener((changes, areaName) => {
  if (areaName !== 'local') return;
  if (changes.trackingState) {
    refreshUI(changes.trackingState.newValue);
  }
  if (changes.options) {
    applyOptions(changes.options.newValue);
  }
});

// Area-specific listener (no need to check areaName):
chrome.storage.local.onChanged.addListener((changes) => {
  for (const [key, { oldValue, newValue }] of Object.entries(changes)) {
    handleLocalChange(key, oldValue, newValue);
  }
});
```

**Note:** Content scripts receive `onChanged` events for `local` and `sync` by default but NOT for `session` unless `setAccessLevel` called (see above).

### Pattern 5: Chunking Large Data

When single value approaches ~1 MB+:

```javascript
const CHUNK_SIZE = 1024 * 1024; // 1 MB per chunk

async function setLargeData(key, data) {
  const serialized = JSON.stringify(data);
  const chunks = [];
  for (let i = 0; i < serialized.length; i += CHUNK_SIZE) {
    chunks.push(serialized.slice(i, i + CHUNK_SIZE));
  }
  const storageObj = {
    [`${key}__meta`]: { chunks: chunks.length, timestamp: Date.now() }
  };
  chunks.forEach((chunk, i) => { storageObj[`${key}__chunk_${i}`] = chunk; });
  await chrome.storage.local.set(storageObj);
}

async function getLargeData(key) {
  const metaResult = await chrome.storage.local.get(`${key}__meta`);
  const meta = metaResult[`${key}__meta`];
  if (!meta) return null;
  const chunkKeys = Array.from({ length: meta.chunks }, (_, i) => `${key}__chunk_${i}`);
  const chunks = await chrome.storage.local.get(chunkKeys);
  return JSON.parse(chunkKeys.map(k => chunks[k]).join(''));
}
```

### Pattern 6: Quota-Aware Writing

```javascript
async function safeSet(area, data) {
  try {
    await chrome.storage[area].set(data);
    return { success: true };
  } catch (error) {
    if (error.message?.includes('QUOTA_BYTES')) {
      const bytesUsed = await chrome.storage[area].getBytesInUse(null);
      console.warn(`Storage quota hit. ${bytesUsed} bytes in use.`);
      await evictOldEntries(area);
      try {
        await chrome.storage[area].set(data);
        return { success: true, evicted: true };
      } catch (retryError) {
        return { success: false, error: retryError.message };
      }
    }
    return { success: false, error: error.message };
  }
}

async function evictOldEntries(area) {
  const all = await chrome.storage[area].get(null);
  const cacheEntries = Object.entries(all)
    .filter(([key]) => key.startsWith('cache_'))
    .sort((a, b) => (a[1].timestamp || 0) - (b[1].timestamp || 0));
  const toRemove = cacheEntries
    .slice(0, Math.ceil(cacheEntries.length * 0.25))
    .map(([key]) => key);
  await chrome.storage[area].remove(toRemove);
}
```

### Pattern 7: Debounced Writes

Coalesce rapid state changes — avoids excessive IPC and write throttling:

```javascript
let writeTimer = null;
let pendingWrites = {};

function debouncedSet(data, delay = 500) {
  Object.assign(pendingWrites, data);
  if (writeTimer) clearTimeout(writeTimer);
  writeTimer = setTimeout(async () => {
    const toWrite = { ...pendingWrites };
    pendingWrites = {};
    writeTimer = null;
    await chrome.storage.local.set(toWrite);
  }, delay);
}

// Flush on service worker suspend
chrome.runtime.onSuspend?.addListener(async () => {
  if (Object.keys(pendingWrites).length > 0) {
    await chrome.storage.local.set(pendingWrites);
    pendingWrites = {};
  }
});
```

### Pattern 8: Schema Versioning

Handle schema evolution across extension updates:

```javascript
const CURRENT_SCHEMA_VERSION = 3;

async function ensureSchema() {
  const { schemaVersion } = await chrome.storage.local.get('schemaVersion');
  if (!schemaVersion || schemaVersion < CURRENT_SCHEMA_VERSION) {
    await runMigrations(schemaVersion || 0);
    await chrome.storage.local.set({ schemaVersion: CURRENT_SCHEMA_VERSION });
  }
}

async function runMigrations(fromVersion) {
  const migrations = [
    null, // v0→v1: no-op
    async () => {
      // v1→v2: rename 'prefs' to 'options'
      const { prefs } = await chrome.storage.local.get('prefs');
      if (prefs) {
        await chrome.storage.local.set({ options: prefs });
        await chrome.storage.local.remove('prefs');
      }
    },
    async () => {
      // v2→v3: add timestamp to cached cases
      const all = await chrome.storage.local.get(null);
      const updates = {};
      for (const [key, value] of Object.entries(all)) {
        if (key.startsWith('case_') && !value.cachedAt) {
          updates[key] = { ...value, cachedAt: Date.now() };
        }
      }
      if (Object.keys(updates).length > 0) {
        await chrome.storage.local.set(updates);
      }
    }
  ];
  for (let v = fromVersion; v < CURRENT_SCHEMA_VERSION; v++) {
    if (migrations[v]) await migrations[v]();
  }
}

chrome.runtime.onInstalled.addListener(() => { ensureSchema(); });
```

### Pattern 9: Migration from localStorage

Service workers cannot access `localStorage`. Use offscreen document for migration:

```javascript
// service-worker.js
async function migrateFromLocalStorage() {
  const { migrated } = await chrome.storage.local.get('migrated');
  if (migrated) return;

  await chrome.offscreen.createDocument({
    url: chrome.runtime.getURL('offscreen.html'),
    reasons: ['LOCAL_STORAGE'],
    justification: 'Migrate localStorage data to chrome.storage.local'
  });

  const response = await chrome.runtime.sendMessage({ type: 'MIGRATE_STORAGE' });
  if (response?.data) {
    await chrome.storage.local.set({ ...response.data, migrated: true });
  }
  await chrome.offscreen.closeDocument();
}

// offscreen.js
chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
  if (msg.type === 'MIGRATE_STORAGE') {
    const data = {};
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      try { data[key] = JSON.parse(localStorage.getItem(key)); }
      catch { data[key] = localStorage.getItem(key); }
    }
    sendResponse({ data });
  }
});
```

## Quota Reference

| Storage Area | Total Quota | Per-Item Limit | Max Items | Write Limits |
|---|---|---|---|---|
| `local` | 10 MB (5 MB in Chrome ≤113) | None | No limit | None |
| `session` | 10 MB (1 MB in Chrome ≤111) | None | No limit | None |
| `sync` | 102,400 bytes | 8,192 bytes | 512 | 1,800/hr, 120/min |
| `managed` | N/A | N/A | N/A | Read-only |

### Checking Usage

```javascript
// Check total bytes in use
const totalBytes = await chrome.storage.local.getBytesInUse(null);
const LOCAL_QUOTA = 10 * 1024 * 1024; // 10 MB
console.log(`Storage: ${(totalBytes / LOCAL_QUOTA * 100).toFixed(1)}% used`);
```

### unlimitedStorage Permission

```json
{ "permissions": ["storage", "unlimitedStorage"] }
```

With `unlimitedStorage`, `chrome.storage.local` no longer capped at 10 MB.

## Anti-Patterns

### 1. Relying on Global Variables in Service Workers

```javascript
// WRONG: lost on SW termination (every 30s of idle)
let activeCase = null;
chrome.runtime.onMessage.addListener((msg) => {
  if (msg.type === 'GET_CASE') return activeCase; // null after restart!
});

// RIGHT: read from storage on every access
chrome.runtime.onMessage.addListener(async (msg, sender, sendResponse) => {
  if (msg.type === 'GET_CASE') {
    const { activeCase } = await chrome.storage.session.get('activeCase');
    sendResponse(activeCase || null);
    return true;
  }
});
```

### 2. Individual Key Reads in Loops

```javascript
// WRONG: N round-trips
for (const id of caseIds) {
  const result = await chrome.storage.local.get(`case_${id}`);
}

// RIGHT: 1 round-trip
const keys = caseIds.map(id => `case_${id}`);
const results = await chrome.storage.local.get(keys);
```

### 3. Non-Serializable Data

```javascript
// WRONG: functions, Maps, Sets, circular refs silently fail or corrupt
await chrome.storage.local.set({ myMap: new Map([['a', 1]]) }); // becomes {}

// RIGHT: convert to plain objects
await chrome.storage.local.set({
  myMap: Object.fromEntries(myMap),
  mySet: [...mySet],
  myDate: Date.now() // store timestamps as numbers
});
```

### 4. Registering onChanged Inside Async Functions

```javascript
// WRONG: may miss events after SW restart
async function init() {
  chrome.storage.onChanged.addListener(handler); // registered too late
}

// RIGHT: register synchronously at module top level
chrome.storage.onChanged.addListener((changes, area) => {
  handleStorageChange(changes, area);
});
```

### 5. Tight-Loop Writes

```javascript
// WRONG: hammers storage and may hit write throttling
items.forEach(async (item) => {
  await chrome.storage.local.set({ [`item_${item.id}`]: item });
});

// RIGHT: single atomic write
const batch = Object.fromEntries(items.map(i => [`item_${i.id}`, i]));
await chrome.storage.local.set(batch);
```

## Troubleshooting

| Symptom | Cause | Fix |
|---|---|---|
| "QUOTA_BYTES quota exceeded" | Write exceeds 10 MB (local/session) or 100 KB (sync) | Call `getBytesInUse(null)`, evict stale cache entries, or add `unlimitedStorage` |
| "QUOTA_BYTES_PER_ITEM" (sync only) | Single item exceeds 8 KB | Split across multiple keys or move to `local` |
| Data missing after SW restart | Stored in JS variable only | Persist all state to `chrome.storage.session` or `local` |
| `onChanged` not firing in content scripts for session | `setAccessLevel` not called | Call `chrome.storage.session.setAccessLevel(...)` in service worker |
| Slow reads/writes with large values | JSON serialization + IPC overhead | Batch keys, chunk values > 1 MB, use `session` for hot data |
| Race conditions with concurrent writes | Two contexts read-modify-write same key | Funnel all writes through service worker via message passing |

## Performance Tips

1. **Prefer session over local for hot data** — in-memory, no disk I/O
2. **Batch all reads and writes** — each call = full IPC round-trip
3. **Avoid `get(null)` on large stores** — retrieves all keys; expensive at scale
4. **Cache locally after read** — keep JS variable as read cache, invalidate on `onChanged`
5. **Debounce rapid writes** — coalesce mutations within 100–500ms windows
6. **Structure keys for partial reads** — prefer `{ case_123: {...} }` over `{ allCases: [...] }` for cheap individual fetches
7. **Monitor usage** — periodically call `getBytesInUse()`, alert on high utilization

## Cross-Browser Compatibility

| Feature | Chrome (Chromium) | Firefox | Safari |
|---------|-------------------|---------|--------|
| API namespace | `chrome.storage` | `browser.storage` (Promise) or `chrome.storage` (callback) | `browser.storage` |
| `storage.session` | Chrome 102+ | Firefox 115+ | Safari 16.4+ |
| `setAccessLevel` | Supported | Not supported (content scripts always have session access) | Not supported |
| `unlimitedStorage` | Supported | Unlimited local by default | Not supported |
| `local` quota | 10 MB (5 MB in Chrome ≤113) | Unlimited by default | ~10 MB |
| `session` quota | 10 MB (1 MB in Chrome ≤111) | 10 MB | ~10 MB |
| Offscreen documents | Chrome 109+ | Not supported | Not supported |

**Portability notes:**
- Use `webextension-polyfill` for cross-browser extensions — normalizes to Promise-based API.
- `setAccessLevel` Chrome-only. Firefox: content scripts access `session` by default, no extra call needed.
- For `localStorage` migration in Firefox, use background scripts instead of offscreen documents.

## References

- [Chrome Storage API](https://developer.chrome.com/docs/extensions/reference/api/storage)
- [Extension Service Worker Lifecycle](https://developer.chrome.com/docs/extensions/develop/concepts/service-workers/lifecycle)
- [Migrate to a Service Worker (MV3)](https://developer.chrome.com/docs/extensions/mv3/migrating_to_service_workers/)