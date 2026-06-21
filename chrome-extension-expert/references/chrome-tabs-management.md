<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly the standalone `chrome-tabs-management` skill.
> Sibling topics in this family are now reference files under the hubs (`chrome-extension-expert`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: chrome-tabs-management
description: >
  Chrome extension tabs API patterns — chrome.tabs query/get/update/create/remove,
  tab messaging (sendMessage to specific tabs), tab lifecycle events
  (onUpdated/onCreated/onRemoved/onActivated), activeTab vs tabs permission,
  chrome.tabGroups API, tab discarding and freezing, programmatic injection
  via chrome.scripting with tab targets, per-tab badge state via chrome.action,
  and finding extension-owned pages.
  TRIGGER: building or reviewing tab management logic in MV3 Chrome extensions,
  querying tabs by URL or state, messaging content scripts in specific tabs,
  handling tab lifecycle events, grouping or discarding tabs, per-tab badge
  state, injecting scripts into specific tabs on demand.
  SKIP: basic MV3 architecture, storage, or SW lifecycle (use chrome-dev or
  mv3-service-worker-expert), offscreen documents or native messaging (use
  chrome-mv3-advanced), extension packaging or publishing (use
  chrome-extension-packaging), E2E testing of extension pages (use
  extension-e2e-testing).
version: 1.1.0
category: developer
tags: [chrome-extension, tabs, mv3, service-worker, messaging, lifecycle, tabGroups, scripting, badge, activeTab, discard, freeze]
related_skills: [mv3-service-worker-expert, chrome-storage-patterns, chrome-notifications-patterns, extension-e2e-testing]
updated: 2026-05-29
---

# Chrome Tabs Management — MV3 Developer Reference

Complete reference for Chrome extension tab operations in Manifest V3. Covers querying, creating, updating, removing, messaging, lifecycle events, tab groups, discarding, programmatic injection, per-tab badge state, and finding extension pages.

All code examples target Manifest V3 with promise-based APIs (no callbacks). Current as of Chrome 132+ / May 2026.

---

## 1. Permissions Model

### The `tabs` permission

Grants access to four sensitive Tab properties:

| Property | Requires `tabs` | Description |
|----------|-----------------|-------------|
| `url` | Yes | Full URL of the tab |
| `pendingUrl` | Yes | URL being navigated to |
| `title` | Yes | Page title |
| `favIconUrl` | Yes | Favicon URL |
| `id`, `windowId`, `active`, `pinned`, `status`, `index`, `groupId`, `discarded` | No | Always available |
| `frozen` | No | Whether tab is frozen (Chrome 132+) |

**Warning shown to user:** "Read your browsing history" — only request when you genuinely need URL/title access.

### The `activeTab` permission

Grants temporary access to the **current tab only** when the user invokes the extension (action click, keyboard shortcut, context menu, permissions prompt).

```json
{ "permissions": ["activeTab"] }
```

**What it grants temporarily:**
- `chrome.scripting.executeScript()` / `insertCSS()` / `removeCSS()` on that tab
- Access to `url`, `title`, `favIconUrl` for that tab
- Behaves like a temporary host permission for the tab's origin

**No warning text shown** — prefer `activeTab` over `"tabs"` when you only need the current tab on user gesture.

### Permission decision matrix

| Need | Permission |
|-------------------------------------------------|------------------------|
| Read URL/title of any tab at any time | `"tabs"` |
| Inject script into active tab on user click | `"activeTab"` + `"scripting"` |
| Query tabs by URL pattern | `"tabs"` |
| Query tabs by window/active/pinned (no URL) | None |
| Create/remove/update/move tabs | None |
| Send messages to a tab's content scripts | None |
| Group/ungroup tabs | None (+ `"tabGroups"` to modify group properties) |
| Discard a tab | None |

---

## 2. Tab CRUD Operations

### chrome.tabs.query(queryInfo)

```js
// Active tab in current window
const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });

// All tabs on a specific domain (requires "tabs" permission)
const supportTabs = await chrome.tabs.query({
  url: ['https://support.mongodb.com/*', 'https://jira.mongodb.org/*']
});

// All discarded tabs
const discardedTabs = await chrome.tabs.query({ discarded: true });

// All loading tabs in focused window
const loadingTabs = await chrome.tabs.query({ lastFocusedWindow: true, status: 'loading' });
```

**Key queryInfo properties:** `active`, `audible`, `currentWindow`, `lastFocusedWindow`, `discarded`, `frozen` (Chrome 132+), `groupId`, `pinned`, `status` (`"loading"` or `"complete"`), `title` (requires `"tabs"`), `url` (requires `"tabs"`), `windowId`, `windowType`.

### Tab CRUD — create, update, remove, move, duplicate, reload

```js
// Create next to current tab with back-button support
const [current] = await chrome.tabs.query({ active: true, currentWindow: true });
const newTab = await chrome.tabs.create({
  url: 'https://example.com',
  index: current.index + 1,
  openerTabId: current.id
});

// Navigate, focus, pin, mute
await chrome.tabs.update(tabId, { url: 'https://new-url.example.com' });
await chrome.tabs.update(tabId, { active: true });
await chrome.tabs.update(tabId, { pinned: true });

// Close multiple tabs
await chrome.tabs.remove([tabId1, tabId2, tabId3]);

// Move tab to different window
await chrome.tabs.move(tabId, { windowId: targetWindowId, index: -1 });

// Duplicate and reload
const clone = await chrome.tabs.duplicate(tabId);
await chrome.tabs.reload(tabId, { bypassCache: true });
```

---

## 3. Tab Messaging

### Service worker to specific tab

```js
// Send to content script in a specific tab
const response = await chrome.tabs.sendMessage(tabId, {
  type: 'GET_PAGE_DATA',
  selector: '.case-header'
});

// Send to specific frame (0 = main frame)
const resp = await chrome.tabs.sendMessage(tabId, { type: 'EXTRACT_FORM' }, { frameId: 0 });
```

### Content script receiver

```js
// content-script.js
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === 'GET_PAGE_DATA') {
    const el = document.querySelector(message.selector);
    sendResponse({ text: el?.textContent ?? null });
    return; // synchronous response
  }
  if (message.type === 'ASYNC_OPERATION') {
    doAsyncWork().then(result => sendResponse(result));
    return true; // REQUIRED for async sendResponse
  }
});
```

### Broadcast to all matching tabs

```js
async function broadcastToSupportTabs(message) {
  const tabs = await chrome.tabs.query({ url: 'https://support.mongodb.com/*' });
  const results = await Promise.allSettled(
    tabs.map(tab => chrome.tabs.sendMessage(tab.id, message))
  );
  // allSettled — some tabs may not have content scripts loaded
  return results;
}
```

### Messaging frozen/discarded tabs

Since Chrome 132, messages to **frozen** tabs are queued and delivered on unfreeze. Messages to **discarded** tabs fail — content script is unloaded.

```js
async function safeTabMessage(tabId, message) {
  const tab = await chrome.tabs.get(tabId);
  if (tab.discarded) throw new Error(`Tab ${tabId} is discarded; reload first`);
  if (tab.frozen) console.warn(`Tab ${tabId} is frozen; message queued`);
  return chrome.tabs.sendMessage(tabId, message);
}
```

---

## 4. Tab Lifecycle Events

All lifecycle events fire in the service worker and wake it if suspended.

```js
// onCreated — URL may be empty at this point
chrome.tabs.onCreated.addListener((tab) => {
  console.log('New tab:', tab.id); // use onUpdated for URL
});

// onUpdated — fires many times per navigation; always filter early
chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
  if (changeInfo.status !== 'complete') return;
  if (!tab.url?.includes('support.mongodb.com')) return;
  sendCaseContextUpdate(tabId);
});

// MV3 event filter (Chrome 88+) — reduces unnecessary SW wakeups
chrome.tabs.onUpdated.addListener(
  (tabId, changeInfo, tab) => { injectOverlayIfNeeded(tabId); },
  { urls: ['https://support.mongodb.com/*'] }  // requires "tabs" permission
);

// onRemoved
chrome.tabs.onRemoved.addListener((tabId, removeInfo) => {
  cleanupTabState(tabId);
  if (removeInfo.isWindowClosing) cleanupWindowState(removeInfo.windowId);
});

// onActivated
chrome.tabs.onActivated.addListener(async ({ tabId }) => {
  const tab = await chrome.tabs.get(tabId);
  updateBadgeForTab(tabId);
});

// onReplaced — transfer state from old tab ID to new (prerendering)
chrome.tabs.onReplaced.addListener((addedTabId, removedTabId) => {
  const state = tabStateMap.get(removedTabId);
  if (state) { tabStateMap.set(addedTabId, state); tabStateMap.delete(removedTabId); }
});
```

**Lifecycle event ordering** for a typical navigation in an existing tab:
1. `onUpdated` — `status: 'loading'`, `url: <new URL>`
2. `onUpdated` — `title: <new title>` (may fire multiple times)
3. `onUpdated` — `favIconUrl: <icon URL>`
4. `onUpdated` — `status: 'complete'`

---

## 5. Tab Groups

Requires `"tabGroups"` permission to **modify** group properties (title, color, collapsed). Grouping/ungrouping tabs requires no extra permission.

```js
// Group tabs into a new group
const groupId = await chrome.tabs.group({ tabIds: [tabId1, tabId2] });

// Add tab to existing group
await chrome.tabs.group({ tabIds: [tabId3], groupId });

// Ungroup
await chrome.tabs.ungroup([tabId1, tabId2]);

// Set title and color (requires "tabGroups")
await chrome.tabGroups.update(groupId, { title: 'Support Cases', color: 'blue' });

// Collapse/expand
await chrome.tabGroups.update(groupId, { collapsed: true });

// Query groups
const blueGroups = await chrome.tabGroups.query({ color: 'blue' });

// Move group to different window
await chrome.tabGroups.move(groupId, { windowId: targetWindowId, index: -1 });
```

**Known limitation:** Chrome's "saved tab groups" (Chrome 119+) cannot be updated via the extension API. Calling `tabGroups.update()` on a saved group throws. There is no API property to detect whether a group is saved.

---

## 6. Tab Discarding and Freezing

```js
// Discard a specific tab (frees memory; tab reloads when user clicks it)
await chrome.tabs.discard(tabId);
// Cannot discard: active tab, already-discarded tab, pinned tab

// Discard stale inactive tabs (tab.lastAccessed is Chrome 121+)
async function discardStaleTabs(maxAgeMs = 30 * 60 * 1000) {
  const tabs = await chrome.tabs.query({ active: false, discarded: false, pinned: false });
  for (const tab of tabs) {
    if (tab.lastAccessed && (Date.now() - tab.lastAccessed) > maxAgeMs) {
      try { await chrome.tabs.discard(tab.id); } catch { /* tab may have become active */ }
    }
  }
}
```

---

## 7. Programmatic Injection via chrome.scripting

Requires `"scripting"` permission plus host access (`activeTab` or host permissions).

```js
// Inject a function with arguments
const results = await chrome.scripting.executeScript({
  target: { tabId },
  func: (selector) => document.querySelector(selector)?.textContent ?? null,
  args: ['.case-number']
});
// results[0].result — the return value from the page

// Inject into all frames
await chrome.scripting.executeScript({
  target: { tabId, allFrames: true },
  func: () => document.title
});

// Inject in MAIN world (shares page's JS context, can access page globals)
await chrome.scripting.executeScript({
  target: { tabId },
  world: 'MAIN',
  func: () => window.__APP_VERSION__
});

// CSS injection
await chrome.scripting.insertCSS({
  target: { tabId },
  css: '.injected-overlay { position: fixed; z-index: 999999; }'
});
await chrome.scripting.removeCSS({
  target: { tabId },
  css: '.injected-overlay { position: fixed; z-index: 999999; }'
});
```

---

## 8. Per-Tab Badge State

`chrome.action` supports per-tab badge text, color, icon, title, and popup. Per-tab values override global defaults and auto-clear when the tab closes.

```js
// Set badge for a specific tab
await chrome.action.setBadgeText({ text: '3', tabId });
await chrome.action.setBadgeBackgroundColor({ color: '#FF0000', tabId });
await chrome.action.setBadgeTextColor({ color: '#FFFFFF', tabId }); // Chrome 110+

// Per-tab icon and popup
await chrome.action.setIcon({ path: { 16: 'icons/active-16.png' }, tabId });
await chrome.action.setPopup({ popup: 'src/popup/special-popup.html', tabId });

// Clear tab-specific badge (falls back to global)
await chrome.action.setBadgeText({ text: '', tabId });

// Update badge on tab switch
chrome.tabs.onActivated.addListener(async ({ tabId }) => {
  const count = await getCaseCountForTab(tabId);
  await chrome.action.setBadgeText({ text: count > 0 ? String(count) : '', tabId });
});
```

---

## 9. Finding Extension Pages

```js
// Find all tabs showing extension pages
const extTabs = await chrome.tabs.query({ url: chrome.runtime.getURL('*') });

// Open-or-focus pattern — if page is already open, focus it; otherwise create it
async function openOrFocusExtensionPage(relativePath) {
  const fullUrl = chrome.runtime.getURL(relativePath);
  const [existing] = await chrome.tabs.query({ url: fullUrl });
  if (existing) {
    await chrome.tabs.update(existing.id, { active: true });
    await chrome.windows.update(existing.windowId, { focused: true });
    return existing;
  }
  return chrome.tabs.create({ url: fullUrl });
}

// Chrome 116+ — query all extension contexts
const contexts = await chrome.runtime.getContexts({});
const tabContexts = contexts.filter(c => c.contextType === 'TAB');
const popupContexts = contexts.filter(c => c.contextType === 'POPUP');
```

---

## 10. Patterns and Recipes

### Wait for tab to finish loading

```js
function waitForTabLoad(tabId, timeoutMs = 30000) {
  return new Promise((resolve, reject) => {
    const timer = setTimeout(() => {
      chrome.tabs.onUpdated.removeListener(listener);
      reject(new Error(`Tab ${tabId} did not finish loading within ${timeoutMs}ms`));
    }, timeoutMs);

    function listener(updatedTabId, changeInfo, tab) {
      if (updatedTabId === tabId && changeInfo.status === 'complete') {
        clearTimeout(timer);
        chrome.tabs.onUpdated.removeListener(listener);
        resolve(tab);
      }
    }
    chrome.tabs.onUpdated.addListener(listener);
  });
}
```

### Track tab state across SW restarts

```js
// Use chrome.storage.session — Map is lost on SW restart
async function persistTabState(tabId, data) {
  const { tabState = {} } = await chrome.storage.session.get('tabState');
  tabState[tabId] = data;
  await chrome.storage.session.set({ tabState });
}

chrome.tabs.onRemoved.addListener(async (tabId) => {
  const { tabState = {} } = await chrome.storage.session.get('tabState');
  delete tabState[tabId];
  await chrome.storage.session.set({ tabState });
});
```

### Deduplicate tabs by URL

```js
async function deduplicateTabs() {
  const tabs = await chrome.tabs.query({ currentWindow: true });
  const seen = new Map();
  const duplicates = [];
  for (const tab of tabs) {
    if (!tab.url) continue;
    if (seen.has(tab.url)) { duplicates.push(tab.id); }
    else { seen.set(tab.url, tab.id); }
  }
  if (duplicates.length > 0) await chrome.tabs.remove(duplicates);
  return duplicates.length;
}
```

### Safe message-or-inject pattern

```js
async function ensureContentScriptAndMessage(tabId, message) {
  try {
    return await chrome.tabs.sendMessage(tabId, message);
  } catch (err) {
    if (err.message?.includes('Receiving end does not exist')) {
      await chrome.scripting.executeScript({
        target: { tabId },
        files: ['src/content/hub-extractor.js']
      });
      return chrome.tabs.sendMessage(tabId, message);
    }
    throw err;
  }
}
```

### Group tabs by domain

```js
async function groupTabsByDomain() {
  const tabs = await chrome.tabs.query({
    currentWindow: true, pinned: false,
    groupId: chrome.tabGroups.TAB_GROUP_ID_NONE
  });
  const byDomain = new Map();
  for (const tab of tabs) {
    try {
      const domain = new URL(tab.url).hostname;
      if (!byDomain.has(domain)) byDomain.set(domain, []);
      byDomain.get(domain).push(tab.id);
    } catch { /* skip tabs without valid URLs */ }
  }
  for (const [domain, tabIds] of byDomain) {
    if (tabIds.length < 2) continue;
    const groupId = await chrome.tabs.group({ tabIds });
    await chrome.tabGroups.update(groupId, { title: domain, color: 'blue' });
  }
}
```

---

## 11. Error Handling

| Error | Cause | Fix |
|-------|-------|-----|
| `No tab with id: N` | Tab closed between query and use | Re-query or wrap in try/catch |
| `Receiving end does not exist` | No content script listener in the tab | Inject content script first, then message |
| `Cannot access a chrome:// URL` | Injecting into a protected page | Check `tab.url` before injecting |
| `Cannot discard the active tab` | Tried to discard focused tab | Only discard inactive tabs |
| `Saved tab groups cannot be updated` | Tried to update a saved tab group | No workaround — API limitation |
| `The extensions gallery cannot be scripted` | Injecting into Chrome Web Store | Skip `chrome.google.com/webstore` URLs |

```js
// Defensive tab access
async function safeGetTab(tabId) {
  try { return await chrome.tabs.get(tabId); }
  catch { return null; } // tab was closed
}
```

---

## 12. Anti-Patterns

1. **Querying all tabs without filters.** `chrome.tabs.query({})` returns every tab. Always filter by `currentWindow`, `active`, `url`, or other criteria.
2. **Not filtering onUpdated events.** This event fires many times per navigation. Always check `changeInfo` properties before doing work.
3. **Assuming `tab.url` is always available.** Without `"tabs"` permission, `tab.url` is undefined. Check for the permission or use `activeTab`.
4. **Sending messages without error handling.** `chrome.tabs.sendMessage()` throws if no listener exists. Always catch `"Receiving end does not exist"`.
5. **Storing tab IDs long-term.** Tab IDs are not stable across browser restarts. Use URL patterns or session storage for persistence.
6. **Using `chrome.tabs.executeScript` (MV2).** In MV3, use `chrome.scripting.executeScript()` — the old API is removed.
7. **Calling `chrome.tabs.getCurrent()` from the service worker.** Returns `undefined` — only works from extension pages running in a tab.
8. **Not checking discarded/frozen state before messaging.** Discarded tabs have no content script. Frozen tabs queue messages. Check state first.
