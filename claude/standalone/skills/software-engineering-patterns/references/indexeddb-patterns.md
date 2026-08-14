<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `indexeddb-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: indexeddb-patterns
title: "IndexedDB Patterns"
description: >
  IndexedDB usage patterns for Chrome MV3 extensions and web apps — object stores, indexes,
  transactions, cursor iteration, IDBKeyRange queries, versioning/migration, Dexie.js wrapper,
  service worker constraints, and large data storage.
  TRIGGER: user asks about IndexedDB object stores, IDBKeyRange queries, cursor iteration,
  Dexie.js schema or queries, IndexedDB versioning and migration (onupgradeneeded), storing
  blobs or large data in the browser, chrome.storage vs IndexedDB decision, IndexedDB in a
  Chrome MV3 service worker, or offline-first browser storage patterns.
  SKIP: browser-side structured storage questions in a Chrome extension with no IndexedDB
  specifics (use chrome-storage-patterns); Dexie.js live query React integration beyond
  basic usage (use dexie-indexeddb-local-first-reviewer); server-side database patterns
  (use mongodb-expert or database-migrations).
version: "1.1.0"
updated: "2026-05-29"
category: developer
tags: [indexeddb, chrome-extension, storage, offline-first, dexie, browser-storage, mv3, service-worker]
keywords:
  - IndexedDB object store
  - IDBKeyRange query
  - cursor iteration getAll
  - Dexie.js wrapper
  - IndexedDB versioning migration onupgradeneeded
  - chrome.storage vs IndexedDB
  - IndexedDB service worker MV3
  - blob large data storage
  - offline-first browser storage
  - fake-indexeddb testing
whenToUse:
  - Implementing or reviewing IndexedDB object stores and indexes in a web app or extension
  - Writing IDBKeyRange range queries or cursor-based iteration
  - Using Dexie.js for schema declaration, CRUD, bulk operations, or version migrations
  - Deciding between chrome.storage.local and IndexedDB for a Chrome extension
  - Handling IndexedDB in a Chrome MV3 service worker (reopen on wake, short transactions)
  - Storing blobs, ArrayBuffers, or large structured datasets in the browser
  - Implementing TTL expiry, chunked storage, or storage quota management
  - Testing IndexedDB code with fake-indexeddb in Node.js
whenNotToUse:
  - Chrome extension storage questions not requiring IndexedDB specifics (use chrome-storage-patterns)
  - Advanced Dexie.js + React live query patterns (use dexie-indexeddb-local-first-reviewer)
  - Server-side or Node.js database patterns (use mongodb-expert or database-migrations)
  - Safari ITP eviction in a PWA without extension context (use databases-aws-cockroach-indexeddb)
related_skills:
  - chrome-storage-patterns
  - dexie-indexeddb-local-first-reviewer
  - databases-aws-cockroach-indexeddb
  - mv3-service-worker-expert
  - database-migrations
triggers:
  - indexeddb
  - indexed db
  - object store
  - IDBKeyRange
  - cursor iteration
  - dexie
  - client-side database
  - browser storage large data
  - offline-first storage
  - chrome extension storage alternative
  - idb
---

# IndexedDB Patterns

Expert reference for IndexedDB in web apps and Chrome MV3 extensions. A response from this skill is correct when it identifies the right transaction mode, avoids the macrotask-boundary pitfall that auto-closes transactions, and uses the appropriate API (raw IndexedDB vs Dexie.js) for the stated constraints.

**Navigation by task:**
- When to use IndexedDB vs chrome.storage vs localStorage → §10 Decision Matrix
- Opening a database, schema design, object store options → §2
- Transaction modes, multi-store atomic operations, transaction lifetime pitfall → §3
- IDBKeyRange (exact, range, prefix, date range) → §4
- Cursor iteration, pagination, `getAll()` for bulk reads → §5
- Schema versioning, data migrations in onupgradeneeded → §6
- Promise wrapper (zero-dep) → §7
- Dexie.js schema, CRUD, queries, bulk ops, version migration → §8
- Service worker constraints (reopen on wake, short transactions) → §9
- Blob storage, chunked files, TTL expiry, storage quota → §11
- Error handling and common errors table → §12
- Performance tips → §13
- Testing with fake-indexeddb → §14
- Raw API vs Dexie quick-reference table → §15

## When to Use

- Structured client-side storage beyond key-value pairs
- Storing large datasets (blobs, files, structured records) in the browser
- Building offline-first or cache-heavy web applications
- Choosing between IndexedDB and chrome.storage in a Chrome extension
- Migrating database schemas across application versions
- Reviewing IndexedDB code for correctness (transaction lifetime, error handling)

## When NOT to Use

- Simple key-value config storage (use localStorage or chrome.storage.local)
- Server-side database design (use MongoDB, PostgreSQL, etc.)
- Session-only ephemeral state (use sessionStorage or chrome.storage.session)
- Small extension settings under 10 MB (chrome.storage.local is simpler)

## Skill Guidance

- Start with the decision matrix (section 10) when choosing between IndexedDB and chrome.storage.
- Use the service worker constraints (section 9) before writing any IndexedDB code in an MV3 background context.
- Prefer Dexie.js (section 8) for new projects unless zero-dependency is a hard requirement.
- Reference the error table (section 12) when debugging IndexedDB failures.

---

## 1. Core Concepts

IndexedDB is a transactional, asynchronous, object-oriented database in every modern browser.

```
Database (name + version)
  +-- Object Store "cases"       (like a table)
  |     +-- Index "by-severity"  (secondary index)
  |     +-- Index "by-date"
  +-- Object Store "attachments"
        +-- Index "by-caseId"
```

| Property | Detail |
|---|---|
| Storage model | Key-value with indexes; values can be JS objects, blobs, arrays |
| Transaction types | `readonly`, `readwrite`, `versionchange` (schema-only) |
| Async model | Event-based (onsuccess/onerror); promise-wrappable |
| Storage limit | ~80% of disk (Chrome); browser-managed eviction in best-effort mode |
| Thread safety | Available in Window, Worker, and Service Worker contexts |
| Structured clone | No functions, DOM nodes, or symbols stored |

---

## 2. Opening a Database and Schema Design

```js
function openDB(name, version) {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(name, version);
    request.onerror = () => reject(request.error);
    request.onsuccess = () => resolve(request.result);
    request.onupgradeneeded = (event) => {
      const db = request.result;
      if (event.oldVersion < 1) {
        const store = db.createObjectStore('cases', { keyPath: 'caseId' });
        store.createIndex('by-severity', 'severity', { unique: false });
        store.createIndex('by-created', 'createdAt', { unique: false });
      }
      if (event.oldVersion < 2) {
        const store = db.createObjectStore('attachments', { keyPath: 'id', autoIncrement: true });
        store.createIndex('by-caseId', 'caseId', { unique: false });
      }
    };
  });
}
```

### Object Store Options

```js
db.createObjectStore('logs', { autoIncrement: true });       // auto-increment key
db.createObjectStore('users', { keyPath: 'email' });         // inline key path
db.createObjectStore('metrics', { keyPath: ['date', 'name'] }); // compound key
db.createObjectStore('blobs');                                // out-of-line keys
```

### Index Options

```js
store.createIndex('by-name', 'name', { unique: false });           // standard
store.createIndex('by-email', 'email', { unique: true });          // unique
store.createIndex('by-tags', 'tags', { multiEntry: true });        // array elements
store.createIndex('by-status-date', ['status', 'updatedAt'], { unique: false }); // compound
```

---

## 3. Transactions

Everything happens inside a transaction. They auto-commit when all requests complete and no new requests are queued.

| Mode | Use case | Concurrency |
|---|---|---|
| `readonly` | Reads, queries, cursor scans | Multiple concurrent txns allowed |
| `readwrite` | Inserts, updates, deletes | Exclusive on same stores |
| `versionchange` | Schema changes (onupgradeneeded) | Exclusive across entire DB |

```js
// Read
const tx = db.transaction('cases', 'readonly');
const record = await idbRequest(tx.objectStore('cases').get(caseId));

// Write (put = upsert; add = insert-only)
const tx = db.transaction('cases', 'readwrite');
tx.objectStore('cases').put(caseRecord);
await idbTransaction(tx);

// Multi-store atomic operation
const tx = db.transaction(['cases', 'attachments'], 'readwrite');
tx.objectStore('cases').delete(caseId);
// ... delete related attachments in same tx
```

### Transaction Lifetime Pitfall

```js
// BROKEN: transaction dies during fetch (macrotask boundary)
const tx = db.transaction('cases', 'readwrite');
const data = await fetch('/api/case/123'); // tx auto-commits here!
tx.objectStore('cases').put(await data.json()); // ERROR: transaction inactive

// CORRECT: fetch first, then transact
const data = await (await fetch('/api/case/123')).json();
const tx = db.transaction('cases', 'readwrite');
tx.objectStore('cases').put(data);
```

---

## 4. IDBKeyRange Queries

```js
IDBKeyRange.only('case-12345');             // exact match
IDBKeyRange.lowerBound(3);                  // >= 3
IDBKeyRange.lowerBound(3, true);            // > 3
IDBKeyRange.upperBound(10);                 // <= 10
IDBKeyRange.upperBound(10, true);           // < 10
IDBKeyRange.bound(3, 10);                   // 3 <= x <= 10
IDBKeyRange.bound(3, 10, true, true);       // 3 < x < 10
```

### Date Range Query

```js
const range = IDBKeyRange.bound(startDate.getTime(), endDate.getTime());
const request = tx.objectStore('cases').index('by-created').openCursor(range);
```

### Prefix Query (String Range)

```js
// All keys starting with "CASE-"
const range = IDBKeyRange.bound(prefix, prefix + '￿', false, false);
```

---

## 5. Cursor Iteration

```js
// Basic forward cursor
const request = store.openCursor();
request.onsuccess = (e) => {
  const cursor = e.target.result;
  if (cursor) { process(cursor.value); cursor.continue(); }
};

// Direction options
store.openCursor(null, 'prev');         // reverse
index.openCursor(null, 'nextunique');   // skip index duplicates
index.openCursor(null, 'prevunique');   // reverse, skip duplicates

// Pagination with advance()
request.onsuccess = (e) => {
  const cursor = e.target.result;
  if (!skipped && offset > 0) { skipped = true; cursor.advance(offset); return; }
  if (cursor && results.length < limit) { results.push(cursor.value); cursor.continue(); }
};
```

### Prefer getAll() for Bulk Reads

```js
store.getAll();                              // all records
store.getAll(null, 100);                     // first 100
index.getAll(IDBKeyRange.only('P1'), 50);    // up to 50 matching records
```

`getAll()` is faster than cursors -- single IPC round-trip vs one per record.

---

## 6. Versioning and Schema Migration

Schema changes only happen inside `onupgradeneeded` when opening a higher version number.

```js
request.onupgradeneeded = (event) => {
  const db = request.result;
  const tx = request.transaction;
  const v = event.oldVersion;

  if (v < 1) { /* create initial stores */ }
  if (v < 2) { /* add new store or index */ }
  if (v < 3) {
    // Data migration via cursor
    const store = tx.objectStore('cases');
    store.openCursor().onsuccess = (e) => {
      const cursor = e.target.result;
      if (cursor) {
        cursor.update({ ...cursor.value, priority: cursor.value.severity === 'P1' ? 'critical' : 'normal' });
        cursor.continue();
      }
    };
  }
};
```

### Multi-Tab Coordination

```js
request.onblocked = () => console.warn('Upgrade blocked -- close other tabs');
db.onversionchange = () => { db.close(); location.reload(); };
```

---

## 7. Promise Wrapper (Minimal, Zero-Dep)

```js
function idbRequest(req) {
  return new Promise((resolve, reject) => {
    req.onsuccess = () => resolve(req.result);
    req.onerror = () => reject(req.error);
  });
}

function idbTransaction(tx) {
  return new Promise((resolve, reject) => {
    tx.oncomplete = () => resolve();
    tx.onerror = () => reject(tx.error);
    tx.onabort = () => reject(tx.error);
  });
}
```

---

## 8. Dexie.js Wrapper

Dexie.js (~5 KB min) -- promise-based, chainable API with bulk operations.

### Schema Declaration

```js
import Dexie from 'dexie';
const db = new Dexie('CaseDatabase');
db.version(1).stores({
  cases: 'caseId, severity, createdAt, owner, [owner+severity]',
  attachments: '++id, caseId',
});
```

Schema syntax: `&` unique, `++` auto-increment PK, `*` multi-entry, `[a+b]` compound.

### CRUD and Queries

```js
await db.cases.add({ caseId: 'CS-001', severity: 'P1', owner: 'alice' });
await db.cases.put({ caseId: 'CS-001', severity: 'P2', owner: 'bob' }); // upsert
await db.cases.update('CS-001', { severity: 'P2' });                     // partial
await db.cases.delete('CS-001');

const p1 = await db.cases.where('severity').equals('P1').toArray();
const recent = await db.cases.where('createdAt').above(Date.now() - 7*86400000).toArray();
const compound = await db.cases.where('[owner+severity]').equals(['alice', 'P1']).toArray();
const multi = await db.cases.where('severity').anyOf(['P1', 'P2']).toArray();
const page = await db.cases.orderBy('createdAt').reverse().offset(20).limit(10).toArray();
```

### Bulk Operations and Transactions

```js
await db.cases.bulkPut(arrayCases);   // fast upsert batch
await db.cases.bulkDelete(idArray);

await db.transaction('rw', db.cases, db.attachments, async () => {
  await db.cases.delete('CS-001');
  await db.attachments.where('caseId').equals('CS-001').delete();
});
```

### Version Migration with Data Transform

```js
db.version(2).stores({ cases: 'caseId, severity, priority' }).upgrade((tx) => {
  return tx.table('cases').toCollection().modify((c) => {
    c.priority = c.severity === 'P1' ? 'critical' : 'normal';
  });
});
```

---

## 9. IndexedDB in Chrome MV3 Extensions

### Availability by Context

| Context | IndexedDB | chrome.storage | localStorage |
|---|---|---|---|
| Service worker | Yes | Yes | No |
| Content script | Yes (page origin) | Yes (extension origin) | Yes (page origin) |
| Popup / Options / Dashboard | Yes (ext origin) | Yes | Yes |
| Offscreen document | Yes (ext origin) | Yes | Yes |

### Service Worker Constraints

The MV3 service worker terminates after ~30s idle / 5min active. Key rules:

1. **Reopen on every wake.** Do not cache `IDBDatabase` globally. Use a lazy singleton that reopens if closed.
2. **Write eagerly.** No guaranteed shutdown hook -- write on every state change.
3. **Keep transactions short.** Do not hold readwrite txns across async gaps.
4. **No long-running cursors.** Use `getAll` or paginated batches.

```js
let _dbPromise = null;
function getDB() {
  if (!_dbPromise) {
    _dbPromise = new Promise((resolve, reject) => {
      const req = indexedDB.open('ExtensionDB', 1);
      req.onupgradeneeded = (e) => {
        if (e.oldVersion < 1) req.result.createObjectStore('cache', { keyPath: 'key' });
      };
      req.onsuccess = () => resolve(req.result);
      req.onerror = () => { _dbPromise = null; reject(req.error); };
    });
  }
  return _dbPromise;
}
function resetDB() { _dbPromise = null; }
```

### Content Script Caveat

Content scripts use the **host page** origin for IndexedDB, not the extension origin. To store private data from a content script, relay it to the service worker via `chrome.runtime.sendMessage`.

---

## 10. chrome.storage vs IndexedDB -- Decision Matrix

| Factor | chrome.storage.local | IndexedDB |
|---|---|---|
| API style | Simple key-value (get/set/remove) | Object stores, indexes, cursors, transactions |
| Max size | 10 MB (unlimited with permission) | ~80% of disk |
| Data types | JSON-serializable only | Structured clone (Blob, ArrayBuffer, Date, Map, Set) |
| Cross-context | All extension contexts | Extension origin only (not content scripts) |
| Query capability | None (get by key) | Indexes, ranges, compound queries |
| Change events | `chrome.storage.onChanged` | None built-in |
| SW safe | Transparent across restarts | Must reopen DB handle |

### When to Use Which

- **Settings/flags/small caches (<10 MB)**: chrome.storage.local
- **Binary data (images, PDFs, audio)**: IndexedDB (chrome.storage cannot store Blob/ArrayBuffer)
- **Large structured datasets (>10 MB)**: IndexedDB
- **Indexed queries on thousands of records**: IndexedDB
- **Hybrid**: Settings in chrome.storage.local (change events, cross-context); large data in IndexedDB

---

## 11. Large Data Storage Patterns

### Blob Storage

```js
tx.objectStore('blobs').put({ key, blob, size: blob.size, storedAt: Date.now() });
```

### Chunked Storage (Very Large Files)

```js
const CHUNK_SIZE = 1024 * 1024; // 1 MB
for (let i = 0; i < Math.ceil(file.size / CHUNK_SIZE); i++) {
  chunkStore.put({ fileId, index: i, data: file.slice(i * CHUNK_SIZE, (i+1) * CHUNK_SIZE) });
}
```

### TTL / Expiry Pattern

```js
// Store with expiry
tx.objectStore(store).put({ ...record, _expiresAt: Date.now() + ttlMs });

// Purge expired via index
const range = IDBKeyRange.upperBound(Date.now());
index.openCursor(range).onsuccess = (e) => {
  const cursor = e.target.result;
  if (cursor) { cursor.delete(); cursor.continue(); }
};
```

### Storage Quota

```js
const { usage, quota } = await navigator.storage.estimate();
await navigator.storage.persist(); // request persistent (prevent eviction)
```

---

## 12. Error Handling

### Comprehensive Transaction Wrapper

```js
async function withIDB(db, stores, mode, fn) {
  return new Promise((resolve, reject) => {
    let tx;
    try { tx = db.transaction(stores, mode); }
    catch (err) { reject(new Error(`Transaction start failed: ${err.message}`)); return; }
    try { fn(tx); } catch (err) { try { tx.abort(); } catch (_) {} reject(err); return; }
    tx.oncomplete = () => resolve();
    tx.onerror = () => reject(tx.error);
    tx.onabort = () => reject(tx.error || new DOMException('Aborted', 'AbortError'));
  });
}
```

### Common Errors

| Error | Cause | Fix |
|---|---|---|
| `InvalidStateError` | DB handle closed (SW restarted) | Reopen DB; reset cached handle |
| `TransactionInactiveError` | Await crossed macrotask boundary | Collect async data before starting tx |
| `ConstraintError` | Duplicate key on `add()` or unique violation | Use `put()` for upsert |
| `QuotaExceededError` | Storage full | Purge old data or request persistence |
| `VersionError` | Requested version < current | Always increment version |
| `AbortError` | Transaction aborted | Check inner error; retry if transient |

---

## 13. Performance Tips

1. **Prefer `getAll()`/`getAllKeys()` over cursors** -- single IPC round-trip.
2. **Batch writes in one transaction** (Dexie `bulkPut`; raw API chained `put`).
3. **Index only what you query.** Each index adds write overhead.
4. **Use `readonly` transactions** when possible -- they run concurrently.
5. **Separate blob stores from metadata** so reads stay fast.
6. **Chrome Snappy compression** handles on-disk compression -- no need to pre-compress.
7. **Storage Buckets** (Chrome 122+): isolate heavy data to prevent evicting critical data.
8. **Batch deletes in chunks** to avoid long readwrite lock holds.

---

## 14. Testing IndexedDB Code

```js
// fake-indexeddb for Node.js unit tests
import 'fake-indexeddb/auto';

test('stores and retrieves', async () => {
  const db = await openDB('test', 1);
  await upsertCase(db, { caseId: 'CS-1', severity: 'P1' });
  expect((await getCase(db, 'CS-1')).severity).toBe('P1');
  db.close();
});

// Dexie unit testing
let db;
beforeEach(() => { db = new Dexie('testDB'); db.version(1).stores({ cases: 'caseId, severity' }); });
afterEach(() => db.delete());

test('query by severity', async () => {
  await db.cases.bulkAdd([{ caseId: 'A', severity: 'P1' }, { caseId: 'B', severity: 'P2' }]);
  expect(await db.cases.where('severity').equals('P1').count()).toBe(1);
});
```

---

## 15. Quick Reference -- Raw API vs Dexie

| Operation | Raw IndexedDB | Dexie.js |
|---|---|---|
| Open DB | `indexedDB.open(name, ver)` | `new Dexie(name).version().stores()` |
| Add | `store.add(obj)` | `table.add(obj)` |
| Upsert | `store.put(obj)` | `table.put(obj)` |
| Get by key | `store.get(key)` | `table.get(key)` |
| Get all | `store.getAll()` | `table.toArray()` |
| Query index | `index.getAll(range)` | `table.where(field).equals(val)` |
| Cursor | `store.openCursor()` | `table.each(fn)` |
| Delete | `store.delete(key)` | `table.delete(key)` |
| Transaction | `db.transaction(stores, mode)` | `db.transaction(mode, tables, fn)` |
| Bulk write | loop `put()` in one tx | `table.bulkPut(arr)` |
| Schema migration | `onupgradeneeded` switch | `.version(n).stores().upgrade()` |
