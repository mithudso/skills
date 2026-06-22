This is the skill file itself, not the file to fix. The compressed file to fix is the content provided in the prompt. The only code block difference is in "Simulate SW termination" — line 3 changed "Trigger a new event" to "Trigger new event".

<!-- hub-reference-banner -->
> **Reference file — part of `chrome-extension-expert` hub.** Formerly standalone `chrome-dev` skill.
> Sibling topics now reference files under hubs (`chrome-extension-expert`) — **not** standalone skills. Ignore "use the X skill" / `related_skills` / SKIP pointers naming bare sibling skills; load `references/<name>.md` from owning hub (see hub's "Cross-hub map").

---

---
name: chrome-dev
version: 1.1.0
category: developer
tags: [chrome-extension, mv3, service-worker, chrome-apis, devtools-mcp, testing, debugging, puppeteer, playwright, jest-chrome, content-script, messaging, storage, permissions]
description: >-
  Chrome MV3 extension development reference — architecture fundamentals, full Chrome API index,
  deep-reference for core APIs (tabs, runtime, scripting, storage, action, offscreen, declarativeNetRequest,
  cookies, identity, alarms, notifications, sidePanel, devtools, debugger), messaging patterns,
  MV3 snippets, chrome-devtools-mcp server (44 tools), extension testing (Puppeteer/Playwright/jest-chrome),
  service worker debugging, common error patterns, performance profiling, and MV3 troubleshooting.
  TRIGGER: Chrome MV3 extension APIs, architecture, or messaging patterns; chrome-devtools-mcp automation;
  extension testing setup; service worker suspension/restart behavior; common extension errors (port
  disconnect, CSP, quota, auto-update race); extension performance profiling; MV2→MV3 migration.
  SKIP: extension E2E test fixtures (use extension-e2e-testing); badge/icon state machines (use
  chrome-badge-metrics); cross-context messaging patterns (use extension-message-bridge); Salesforce
  content-script scraping (use salesforce-scraping-patterns).
triggers:
  - "What Chrome APIs are available in MV3?"
  - "How do I keep a Chrome extension service worker alive?"
  - "chrome.scripting.executeScript example"
  - "Use chrome-devtools-mcp to automate a browser task"
  - "Set up Puppeteer to test a Chrome extension"
  - "Extension gives 'message port closed before response received'"
  - "MV2 to MV3 migration checklist"
  - "chrome.declarativeNetRequest block a URL"
  - "chrome.offscreen createDocument for DOM parsing"
  - "Debug Chrome extension service worker in DevTools"
related_skills:
  - extension-message-bridge
  - extension-e2e-testing
  - chrome-badge-metrics
  - mv3-service-worker-expert
  - chrome-storage-patterns
  - chrome-extension-security-reviewer
---

# Chrome Extension API & DevTools MCP

## MV3 architecture fundamentals

### Execution contexts (no shared memory)

| Context | Has DOM | Chrome APIs | Lifetime |
|---------|---------|-------------|---------|
| Service Worker | No | Almost all | Disposable — suspend/resume |
| Popup | Yes | Almost all | Until popup closed |
| Options page | Yes | Almost all | Until tab closed |
| Side panel | Yes | Almost all | Persistent |
| Offscreen document | Yes | `runtime` only | Explicit create/close |
| Content script | Yes (page DOM) | Limited subset | Page lifetime |
| DevTools page | Yes | `devtools.*` + limited | DevTools open |

### Key constraints

- **Static imports only** in MV3 service worker — `import()` disallowed at runtime.
- **Service worker disposable** — module-global state can vanish; persist via `chrome.storage.*` or IndexedDB.
- Contexts communicate via `chrome.runtime.sendMessage` / `chrome.storage` / IndexedDB.
- All `chrome.*` APIs available as `browser.*` (Firefox alias) since Chrome 146.
- All async APIs return Promises; old callback style still works.

### Minimal MV3 manifest

```json
{
  "manifest_version": 3,
  "name": "My Extension",
  "version": "1.0.0",
  "permissions": ["storage", "tabs", "scripting"],
  "host_permissions": ["https://*.example.com/*"],
  "background": { "service_worker": "src/background/service-worker.js", "type": "module" },
  "action": { "default_popup": "popup.html", "default_icon": "icons/icon48.png" },
  "content_scripts": [{ "matches": ["https://*.example.com/*"], "js": ["src/content/content.js"] }]
}
```

---

## Core APIs — quick reference

### chrome.tabs

```js
const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
await chrome.tabs.create({ url: 'https://example.com', active: false });
await chrome.tabs.update(tabId, { url: 'https://example.com', pinned: true });
await chrome.tabs.remove([tabId1, tabId2]);
const dataUrl = await chrome.tabs.captureVisibleTab(windowId, { format: 'png' });
const response = await chrome.tabs.sendMessage(tabId, { type: 'getData' });
chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
  if (changeInfo.status === 'complete') console.log('loaded', tab.url);
});
```

### chrome.runtime

```js
chrome.runtime.onInstalled.addListener(({ reason, previousVersion }) => { /* install/update */ });
const manifest = chrome.runtime.getManifest();
const url = chrome.runtime.getURL('src/offscreen/offscreen.html');

// One-shot messaging
const response = await chrome.runtime.sendMessage({ type: 'ping' });
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === 'ping') sendResponse({ pong: true });
  return true; // required for async sendResponse
});

// Long-lived port
const port = chrome.runtime.connect({ name: 'myChannel' });
port.postMessage({ data: 'hello' });
chrome.runtime.onConnect.addListener(port => {
  port.onMessage.addListener(msg => port.postMessage({ echo: msg }));
  port.onDisconnect.addListener(() => {});
});
```

### chrome.scripting

```js
// Inject function (preferred)
const results = await chrome.scripting.executeScript({
  target: { tabId, allFrames: false },
  func: (arg) => document.title + ' ' + arg,
  args: ['extra'],
});

// Inject file
await chrome.scripting.executeScript({ target: { tabId }, files: ['src/content/inject.js'] });

// CSS
await chrome.scripting.insertCSS({ target: { tabId }, css: 'body { background: red }' });
await chrome.scripting.removeCSS({ target: { tabId }, css: 'body { background: red }' });

// Register persistent content scripts
await chrome.scripting.registerContentScripts([{
  id: 'my-script', matches: ['https://*.example.com/*'],
  js: ['src/content/content.js'], runAt: 'document_idle',
}]);
```

### chrome.storage

| Area | Quota | Persists | Notes |
|------|-------|---------|-------|
| `local` | 10 MB (`unlimitedStorage`: unlimited) | Disk | Default choice |
| `sync` | 100 KB total / 8 KB per item | Cross-device | Requires signed-in Chrome |
| `session` | 10 MB | Memory only | Fast; clears on browser restart |
| `managed` | — | Enterprise policy | Read-only |

```js
await chrome.storage.local.set({ key: 'value' });
const data = await chrome.storage.local.get(['key']);
const { theme = 'dark' } = await chrome.storage.local.get({ theme: 'dark' });
await chrome.storage.local.remove(['key']);
chrome.storage.onChanged.addListener((changes, area) => {
  for (const [key, { oldValue, newValue }] of Object.entries(changes))
    console.log(`${area}/${key}: ${oldValue} → ${newValue}`);
});
```

### chrome.offscreen

Use when service worker needs DOM APIs (clipboard, DOMParser, audio, WebRTC).

```js
async function ensureOffscreen(path) {
  const url = chrome.runtime.getURL(path);
  const existing = await chrome.runtime.getContexts({ contextTypes: ['OFFSCREEN_DOCUMENT'], documentUrls: [url] });
  if (existing.length) return;
  await chrome.offscreen.createDocument({ url: path, reasons: ['DOM_PARSER'], justification: 'Parse HTML' });
}

// Offscreen JS
chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
  if (msg.target !== 'offscreen') return false;
  if (msg.type === 'parseHTML') {
    const doc = new DOMParser().parseFromString(msg.data, 'text/html');
    sendResponse({ title: doc.title });
  }
  return true;
});
```

Valid `reasons`: `AUDIO_PLAYBACK`, `IFRAME_SCRIPTING`, `DOM_SCRAPING`, `BLOBS`, `DOM_PARSER`, `USER_MEDIA`, `DISPLAY_MEDIA`, `WEB_RTC`, `CLIPBOARD`, `LOCAL_STORAGE`, `WORKERS`, `BATTERY_STATUS`, `MATCH_MEDIA`, `GEOLOCATION`.

### chrome.declarativeNetRequest

```js
await chrome.declarativeNetRequest.updateDynamicRules({
  addRules: [{
    id: 100, priority: 1,
    action: { type: 'block' },
    condition: { urlFilter: '||tracking.com^', resourceTypes: ['xmlhttprequest'] },
  }],
  removeRuleIds: [99],
});

// Debug: test a URL against rulesets
const outcome = await chrome.declarativeNetRequest.testMatchOutcome({
  request: { url: 'https://ads.com/pixel.gif', type: 'image', tabId: chrome.tabs.TAB_ID_NONE },
});
```

### chrome.cookies / chrome.identity / chrome.alarms

```js
// Cookies
const cookie = await chrome.cookies.get({ url: 'https://example.com', name: 'session_id' });
await chrome.cookies.set({ url: 'https://example.com', name: 'my_cookie', value: 'abc', secure: true });

// Identity (OAuth)
const token = await chrome.identity.getAuthToken({ interactive: true });
const redirectUrl = chrome.identity.getRedirectURL();
const responseUrl = await chrome.identity.launchWebAuthFlow({ url: oauthUrl, interactive: true });

// Alarms (survive SW sleep)
await chrome.alarms.create('refresh', { periodInMinutes: 5 });
chrome.alarms.onAlarm.addListener(alarm => { if (alarm.name === 'refresh') fetchData(); });
```

---

## MV3 snippets

### Keep service worker alive

```js
chrome.alarms.create('keepAlive', { periodInMinutes: 0.4 });
chrome.alarms.onAlarm.addListener(alarm => { if (alarm.name === 'keepAlive') { /* no-op */ } });
```

### Inject once per tab

```js
const injected = new Set(); // resets on SW restart — that's fine
async function injectOnce(tabId) {
  if (injected.has(tabId)) return;
  await chrome.scripting.executeScript({ target: { tabId }, files: ['content.js'] });
  injected.add(tabId);
}
chrome.tabs.onRemoved.addListener(tabId => injected.delete(tabId));
```

### Detect SW vs DOM context

```js
const isServiceWorker = typeof window === 'undefined';
```

---

## Chrome DevTools MCP server

`chrome-devtools-mcp` provides 44 MCP tools for AI agents: screenshots, DOM inspection, script execution, network monitoring, performance tracing, memory profiling, Lighthouse audits.

```bash
npx -y chrome-devtools-mcp@latest                          # latest
npx -y chrome-devtools-mcp@latest --slim --headless        # CI
npx -y chrome-devtools-mcp@latest --browser-url=http://127.0.0.1:9222  # existing Chrome
```

### Key flags

| Flag | Effect |
|------|--------|
| `--slim` | Minimal tool set (nav + input + snapshot) |
| `--headless` | Launch Chrome headless |
| `--categoryExtensions=true` | Enable install/reload/trigger extension tools |
| `--experimentalMemory` | Enable heap snapshot tools |
| `--browser-url=URL` | Connect to existing Chrome |

### Automation principles

1. Use `take_snapshot` (a11y tree) over `take_screenshot` — uid-tagged elements for reliable automation.
2. Use `fill_form` over sequential `fill` + `click` — fewer race conditions.
3. `evaluate_script` runs in page context — same as DevTools console.
4. Use `wait_for` when content loads async after navigation.

### Common MCP recipes

**Fill and submit form:**
```
1. navigate_page { url }
2. take_snapshot  → find input uids
3. fill_form { elements: [{uid, value}, ...] }
4. click { uid: "submit-button" }
5. wait_for { text: ["Success"] }
```

**Run Lighthouse:**
```
1. navigate_page { url }
2. lighthouse_audit { device: "mobile", mode: "navigation" }
```

**Profile memory:**
```
1. take_memory_snapshot { filePath: "/tmp/before.heapsnapshot" }
2. // perform actions
3. take_memory_snapshot { filePath: "/tmp/after.heapsnapshot" }
4. get_nodes_by_class  // compare retained objects
```

---

## Extension testing patterns

### Framework selection

| Framework | Best for | Extension support |
|-----------|---------|-----------------|
| Puppeteer | Deep Chrome automation, MV3 SW introspection | Native; `--load-extension` |
| Playwright | Cross-browser + extension E2E | Chrome only; see `extension-e2e-testing` skill |
| jest-chrome | Unit testing Chrome API calls | Mock `chrome.*` namespace; no browser needed |
| Vitest + chrome shims | Fast unit tests for background logic | Hand-rolled `chrome.storage`/`chrome.runtime` shims |

### Puppeteer extension setup

```js
import puppeteer from 'puppeteer';
import path from 'path';

const extensionPath = path.resolve('./src');
const browser = await puppeteer.launch({
  headless: false,  // extensions require headed or new headless
  args: [`--disable-extensions-except=${extensionPath}`, `--load-extension=${extensionPath}`],
});

const swTarget = await browser.waitForTarget(
  t => t.type() === 'service_worker' && t.url().startsWith('chrome-extension://')
);
const extensionId = new URL(swTarget.url()).hostname;
const popupPage = await browser.newPage();
await popupPage.goto(`chrome-extension://${extensionId}/popup.html`);
```

### jest-chrome unit testing

```js
import { chrome } from 'jest-chrome';
chrome.storage.local.get.mockImplementation((keys, callback) => callback({ theme: 'dark' }));

test('reads theme from storage', async () => {
  const result = await getTheme();
  expect(result).toBe('dark');
});
```

---

## Service worker debugging

### DevTools entry points

| Surface | How to open | What you get |
|---------|-------------|-------------|
| SW DevTools | `chrome://extensions` → "Inspect views: service worker" | Console, breakpoints, network |
| `chrome://serviceworker-internals` | Address bar | All SWs; start/stop/unregister |
| Application → Service Workers | DevTools Application tab | Status, update, push, sync controls |

### Key behaviors

- Inspecting SW keeps it alive — close DevTools to test realistic suspend/resume.
- SW logs go to own console, not tab console.
- Event listeners **must register synchronously** at module top level — not inside `await` or `setTimeout`.

### Simulate SW termination

```
1. Close "Inspect views: service worker" DevTools window
2. Wait ~30s, or click "Stop" in chrome://serviceworker-internals
3. Trigger a new event (alarm, toolbar click)
4. Re-inspect — fresh SW instance; all module-global state gone
```

---

## Common error patterns

### "Message port closed before response received"

Listener used `async` function or forgot `return true`. Fix:
```js
chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
  fetch('/api').then(r => r.json()).then(data => sendResponse(data));
  return true; // REQUIRED
});
```

### "Could not establish connection. Receiving end does not exist."

Content script not injected, or SW not started. Fix:
```js
async function sendToTab(tabId, message) {
  try {
    return await chrome.tabs.sendMessage(tabId, message);
  } catch {
    await chrome.scripting.executeScript({ target: { tabId }, files: ['content.js'] });
    return await chrome.tabs.sendMessage(tabId, message);
  }
}
```

### Port disconnected

SW terminated mid-stream. Fix with reconnect:
```js
function connect() {
  const port = chrome.runtime.connect({ name: 'stream' });
  port.onDisconnect.addListener(() => setTimeout(connect, 300));
  port.onMessage.addListener(handleMessage);
}
connect();
```

### "Refused to execute inline script" (CSP)

MV3 bans `unsafe-eval` and `unsafe-inline` in extension pages. Move inline scripts to `.js` files. Never use `innerHTML` to inject `<script>` tags.

### Storage quota exceeded

| Area | Limit | Fix |
|------|-------|-----|
| `storage.local` | 10 MB | Add `"unlimitedStorage"` permission |
| `storage.sync` | 100 KB / 8 KB per item | Split objects; use `local` |
| `storage.session` | 10 MB | Prune on SW startup |

---

## Extension performance profiling

### SW performance budget

- Keep static `import` graph shallow — all imports parsed synchronously on cold start.
- Move heavy parsing (large JSON, complex regex) to first-use lazy init, not module top level.
- Use `chrome.storage.session` for frequently read settings (memory-only, fast) over `storage.local` (disk I/O).

### Extension-specific memory leak patterns

| Pattern | Fix |
|---------|-----|
| `chrome.storage.onChanged` listeners added in popup, never removed | Remove in `window.unload` |
| Content scripts accumulating DOM references after navigation | Nullify in `beforeunload` |
| Offscreen document kept open indefinitely | Call `chrome.offscreen.closeDocument()` when idle |
| Module globals holding stale tab objects | Store only `tab.id`; re-query on use |

---

## MV2 → MV3 migration checklist

| MV2 pattern | MV3 replacement |
|-------------|----------------|
| `background.scripts` | `background.service_worker` (single entry file) |
| `background.persistent: true` | Remove — SW always non-persistent |
| `XMLHttpRequest` in background | `fetch()` |
| `chrome.browserAction` | `chrome.action` |
| `webRequest` blocking | `declarativeNetRequest` |
| Inline scripts in HTML | External `.js` files (MV3 CSP bans inline) |
| `chrome.extension.getBackgroundPage()` | `chrome.runtime.getContexts()` |
| `chrome.extension.sendRequest` | `chrome.runtime.sendMessage` |

### Extension error triage checklist

1. `chrome://extensions` → **Errors** button → read stack trace.
2. "Inspect views: service worker" → Console → filter `[ERROR]`.
3. Check `chrome.runtime.lastError` after every Chrome API call in callbacks.
4. Verify SW version: `chrome.runtime.getManifest().version`.
5. Check `chrome://serviceworker-internals` — running or stopped?
6. Application → Storage in DevTools — verify `chrome.storage.local` state.
7. Content script issues: open target page DevTools → Console (content script errors appear there, not in SW DevTools).