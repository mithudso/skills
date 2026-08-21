<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly standalone `chrome-mv3-advanced` skill.
> Sibling topics now reference files under hubs (`chrome-extension-expert`) — **not** standalone
> skills. Ignore "use the X skill" / `related_skills` / SKIP pointers naming bare sibling
> skills; load that topic's `references/<name>.md` from owning hub (see hub's "Cross-hub map").

---

---
name: chrome-mv3-advanced
description: >
  Advanced Chrome MV3 extension APIs — deep reference for offscreen documents,
  native messaging, chrome.identity OAuth flows, declarativeNetRequest, side
  panel, userScripts, tabCapture, and offscreen+SW coordination patterns.
  TRIGGER: implementing or debugging chrome.offscreen (createDocument, Reason
  enum, single-instance constraint, SSE/WebSocket hosting, LLM stream hosting),
  native messaging (connectNative, sendNativeMessage, 4-byte framing, host
  manifest, install.sh, Python/Node.js host patterns), chrome.identity OAuth
  (getAuthToken, launchWebAuthFlow, PKCE, token refresh, removeCachedAuthToken),
  declarativeNetRequest (static/dynamic/session rules, modifyHeaders,
  testMatchOutcome), chrome.sidePanel (tab-specific panels, setPanelBehavior),
  chrome.userScripts (USER_SCRIPT world, configureWorld), chrome.tabCapture /
  desktopCapture (stream IDs, offscreen handoff), or offscreen+SW coordination
  (SSE forwarding, LLM stream hosting, keepalive patterns).
  SKIP: basic chrome.tabs, chrome.storage, chrome.runtime messaging, or SW
  fundamentals → use chrome-dev. Browser automation or Playwright/Puppeteer →
  use chrome-devtools-mcp. Simple MV2→MV3 migration without offscreen or native
  messaging → use chrome-dev. For standalone deep dives: chrome-offscreen-documents,
  chrome-native-messaging, chrome-identity-oauth each have their own skill.
version: 1.1.0
category: developer
tags: [chrome-extension, mv3, offscreen, native-messaging, identity, oauth, declarativeNetRequest, side-panel, userScripts, tabCapture, service-worker, pkce]
related_skills: [chrome-offscreen-documents, chrome-native-messaging, chrome-identity-oauth, mv3-service-worker-expert, chrome-tabs-management]
updated: 2026-05-29
---

# Chrome MV3 Advanced APIs — Developer Reference

Deep reference for advanced Chrome MV3 APIs beyond basic extension architecture: offscreen documents, native messaging, identity/OAuth, newer platform APIs (side panel, userScripts, tabCapture, DNR).

**When standalone depth needed:** `chrome-offscreen-documents`, `chrome-native-messaging`, `chrome-identity-oauth` each embed full standalone references. This skill covers all topics at useful depth plus cross-topic coordination patterns (Sections 5–6) standalone skills don't cover.

---

## 1. Decision Tables

### When to use Offscreen Document vs Service Worker vs Content Script

| Need | Best choice | Why |
|---|---|---|
| DOMParser, innerHTML parsing | Offscreen (`DOM_PARSER`) | SW has no DOM |
| Clipboard read/write | Offscreen (`CLIPBOARD`) | Clipboard API requires DOM |
| Playing audio silently | Offscreen (`AUDIO_PLAYBACK`) | `<audio>` needs DOM |
| Persistent SSE/WebSocket | Offscreen (`WEB_RTC` or `AUDIO_PLAYBACK` keepalive) | SW suspends; offscreen persists |
| LLM streaming fetch | Offscreen (`AUDIO_PLAYBACK` + keepalive) | SW may suspend mid-stream |
| Scrape a page iframe | Offscreen (`DOM_SCRAPING` or `IFRAME_SCRIPTING`) | Need live DOM with iframe |
| getUserMedia (mic/camera) | Offscreen (`USER_MEDIA`) | Media APIs need DOM |
| localStorage access | Offscreen (`LOCAL_STORAGE`) | localStorage not in SW |
| Simple fetch (no DOM) | Service Worker | No DOM needed |
| Interact with page content | Content Script | Direct DOM access |

### When to use Native Messaging vs Fetch/WebSocket

| Need | Best choice | Why |
|---|---|---|
| Call local CLI tool / binary | Native messaging | Only way to exec local processes |
| Read/write local file system | Native messaging or File System Access API | Browser sandbox blocks direct FS |
| Access OS keychain | Native messaging | No JS API for OS keychain |
| Call local HTTP server | Fetch to `localhost` | Simpler, no host install required |
| Real-time bidirectional with local app | `connectNative` (port) | Keeps process alive, streaming |
| One-shot query to local app | `sendNativeMessage` | Simpler, no port lifecycle |
| Large data (>1 MB) | Split chunks or use local HTTP | Native message from host capped at 1 MB |

### Which Identity API to use

| Scenario | API |
|---|---|
| Google APIs, simplest path | `chrome.identity.getAuthToken` |
| Google APIs but need refresh token | `launchWebAuthFlow` + PKCE |
| Non-Google OAuth (Auth0, GitHub, Okta, Azure) | `launchWebAuthFlow` + PKCE |
| Silent re-auth (no UI) | `getAuthToken({ interactive: false })` |
| Token expired / 401 received | `removeCachedAuthToken` then `getAuthToken` again |

---

## 2. Offscreen Documents

### Overview

`chrome.offscreen` (Chrome 109+, MV3 only) creates hidden HTML document with full DOM access. **One document per extension per profile** — always check before creating.

**Permission required:** `"offscreen"` in manifest.json

### Reason Enum (key entries)

| Reason | Lifetime | Use for |
|---|---|---|
| `AUDIO_PLAYBACK` | **30s auto-close** without audio | `<audio>`, keepalive for SSE/LLM streams |
| `DOM_PARSER` | Unlimited | `DOMParser` API |
| `CLIPBOARD` | Unlimited | `navigator.clipboard` |
| `LOCAL_STORAGE` | Unlimited | `localStorage` access |
| `WEB_RTC` | Unlimited | WebRTC, WebSocket hosting |
| `WORKERS` | Unlimited | Spawn `Worker` threads |
| `USER_MEDIA` | Unlimited | `getUserMedia()` |
| `BLOBS` | Unlimited | `URL.createObjectURL()` |

Use array for multiple reasons: `reasons: ['DOM_PARSER', 'CLIPBOARD']`

### Single-Instance Pattern (Chrome 116+)

```js
// service-worker.js
let _creatingOffscreen = null;

async function ensureOffscreen(path = 'src/offscreen/offscreen.html') {
  const url = chrome.runtime.getURL(path);
  const existing = await chrome.runtime.getContexts({
    contextTypes: ['OFFSCREEN_DOCUMENT'],
    documentUrls: [url],
  });
  if (existing.length > 0) return;

  if (_creatingOffscreen) { await _creatingOffscreen; return; }

  _creatingOffscreen = chrome.offscreen.createDocument({
    url,
    reasons: ['DOM_PARSER', 'CLIPBOARD'],
    justification: 'Parsing HTML content and writing to clipboard',
  });
  await _creatingOffscreen;
  _creatingOffscreen = null;
}
```

**Fallback for Chrome < 116:**
```js
async function offscreenExistsFallback(path) {
  const url = chrome.runtime.getURL(path);
  const clients = await self.clients.matchAll();
  return clients.some(c => c.url === url);
}
```

### Communication Pattern

```js
// service-worker.js
await ensureOffscreen('src/offscreen/offscreen.html');
const result = await chrome.runtime.sendMessage({
  target: 'offscreen', type: 'parseHTML', payload: rawHtml,
});

// offscreen.js
chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
  if (msg.target !== 'offscreen') return false;
  switch (msg.type) {
    case 'parseHTML': {
      const doc = new DOMParser().parseFromString(msg.payload, 'text/html');
      sendResponse({ title: doc.title, text: doc.body.innerText });
      break;
    }
    case 'clipboardWrite': {
      navigator.clipboard.writeText(msg.text)
        .then(() => sendResponse({ ok: true }))
        .catch(e => sendResponse({ error: e.message }));
      return true; // async
    }
  }
  return true;
});
```

### Audio Keepalive (for SSE/LLM streams)

```html
<!-- offscreen.html -->
<audio id="keepalive" loop
  src="data:audio/wav;base64,UklGRiQAAABXQVZFZm10IBAAAAABAAEAQB8AAEAfAAABAAgAZGF0YQAAAAA=">
</audio>
```
```js
// offscreen.js
document.getElementById('keepalive').play().catch(() => {});
```

---

## 3. Native Messaging

**Protocol:** 4-byte little-endian uint32 length prefix + UTF-8 JSON payload.
**Limits:** Chrome→host 64 MB; host→Chrome **1 MB max** (hard limit).
**Debug output must go to stderr** — any bytes on stdout corrupt framing.

**Permission:** `"nativeMessaging"` in manifest.json

### Host Manifest

```json
{
  "name": "com.example.myapp",
  "description": "My native messaging host",
  "path": "/absolute/path/to/host.py",
  "type": "stdio",
  "allowed_origins": ["chrome-extension://YOUR_EXTENSION_ID/"]
}
```

**OS registration paths:**
```
macOS user:   ~/Library/Application Support/Google/Chrome/NativeMessagingHosts/<name>.json
Linux user:   ~/.config/google-chrome/NativeMessagingHosts/<name>.json
Windows:      HKCU\Software\Google\Chrome\NativeMessagingHosts\<name> → path to JSON
```

### Python Host (minimal)

```python
#!/usr/bin/env python3
import sys, json, struct

def read_message():
    raw = sys.stdin.buffer.read(4)
    if not raw: sys.exit(0)
    return json.loads(sys.stdin.buffer.read(struct.unpack('<I', raw)[0]))

def send_message(data):
    msg = json.dumps(data, separators=(',',':')).encode('utf-8')
    sys.stdout.buffer.write(struct.pack('<I', len(msg)) + msg)
    sys.stdout.buffer.flush()

while True:
    try:
        msg = read_message()
        send_message({'echo': msg})
    except Exception as e:
        sys.stderr.write(f'Error: {e}\n')
        sys.exit(1)
```

### Chrome-Side Usage

```js
// Persistent connection (keeps SW alive while connected)
const port = chrome.runtime.connectNative('com.example.myapp');
port.onMessage.addListener(msg => console.log('From host:', msg));
port.onDisconnect.addListener(() => {
  if (chrome.runtime.lastError) console.error(chrome.runtime.lastError.message);
});
port.postMessage({ command: 'start' });

// One-shot (new process per call)
const response = await chrome.runtime.sendNativeMessage('com.example.myapp', { cmd: 'ping' });
```

**Note:** `connectNative` available in extension pages and service worker only — not content scripts.

---

## 4. Chrome Identity API

### getAuthToken — Google Only

```js
// Non-interactive: returns cached token or null
async function getToken(interactive = false) {
  try {
    const result = await chrome.identity.getAuthToken({ interactive });
    return result.token; // MV3: returns { token, grantedScopes }
  } catch { return null; }
}

// On 401: flush cache and retry
async function callGoogleAPI(endpoint) {
  const token = await getToken(true);
  const resp = await fetch(endpoint, { headers: { Authorization: `Bearer ${token}` } });
  if (resp.status === 401) {
    await chrome.identity.removeCachedAuthToken({ token });
    const fresh = await getToken(true);
    return fetch(endpoint, { headers: { Authorization: `Bearer ${fresh}` } });
  }
  return resp;
}
```

### launchWebAuthFlow — Any Provider (PKCE)

```js
async function generatePKCE() {
  const bytes = crypto.getRandomValues(new Uint8Array(32));
  // Use Array.from() to avoid call-stack limits on large typed arrays
  const verifier = btoa(Array.from(bytes, b => String.fromCharCode(b)).join(''))
    .replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '');
  const hash = await crypto.subtle.digest('SHA-256', new TextEncoder().encode(verifier));
  const challenge = btoa(Array.from(new Uint8Array(hash), b => String.fromCharCode(b)).join(''))
    .replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '');
  return { verifier, challenge };
}

async function launchOAuthFlow(provider) {
  const { verifier, challenge } = await generatePKCE();
  const redirectUrl = chrome.identity.getRedirectURL();
  const authUrl = new URL(provider.authEndpoint);
  authUrl.searchParams.set('client_id', provider.clientId);
  authUrl.searchParams.set('redirect_uri', redirectUrl);
  authUrl.searchParams.set('response_type', 'code');
  authUrl.searchParams.set('scope', provider.scopes.join(' '));
  authUrl.searchParams.set('code_challenge', challenge);
  authUrl.searchParams.set('code_challenge_method', 'S256');
  authUrl.searchParams.set('state', crypto.randomUUID());

  const responseUrl = await chrome.identity.launchWebAuthFlow({
    url: authUrl.toString(), interactive: true,
  });
  const code = new URL(responseUrl).searchParams.get('code');
  return { code, verifier };
}
```

---

## 5. Advanced MV3 APIs

### chrome.declarativeNetRequest (DNR)

MV3 replacement for blocking `webRequest`. Rules execute in browser without JS.

| Rule type | Max rules | Cleared on restart? |
|---|---|---|
| Static (manifest rulesets) | 30,000 across all rulesets | No |
| Dynamic | 5,000 | No |
| Session | 5,000 | Yes |

```js
// Block a tracker
await chrome.declarativeNetRequest.updateDynamicRules({
  addRules: [{
    id: 1, priority: 1,
    action: { type: 'block' },
    condition: {
      urlFilter: '||ads.tracker.com^',
      resourceTypes: ['script', 'xmlhttprequest', 'image'],
    }
  }],
  removeRuleIds: [],
});

// Remove CSP to allow iframe embedding
await chrome.declarativeNetRequest.updateDynamicRules({
  addRules: [{
    id: 10, priority: 1,
    action: {
      type: 'modifyHeaders',
      responseHeaders: [
        { header: 'content-security-policy', operation: 'remove' },
        { header: 'x-frame-options', operation: 'remove' },
      ],
    },
    condition: { urlFilter: 'https://target.example.com/*', resourceTypes: ['main_frame', 'sub_frame'] },
  }],
});

// Test if a URL would be matched
const outcome = await chrome.declarativeNetRequest.testMatchOutcome({
  request: { url: 'https://ads.tracker.com/pixel.gif', type: 'image', tabId: -1 }
});
```

### chrome.sidePanel (Chrome 114+)

```json
{ "side_panel": { "default_path": "sidepanel/sidepanel.html" } }
```

```js
// Open programmatically (must be called from user gesture)
chrome.action.onClicked.addListener(async (tab) => {
  await chrome.sidePanel.open({ windowId: tab.windowId });
});

// Per-tab panels
chrome.tabs.onActivated.addListener(async ({ tabId }) => {
  const tab = await chrome.tabs.get(tabId);
  await chrome.sidePanel.setOptions({
    tabId,
    path: tab.url.includes('github.com') ? 'sidepanel/github-panel.html' : 'sidepanel/default.html',
    enabled: true,
  });
});

// Auto-open on action click
await chrome.sidePanel.setPanelBehavior({ openPanelOnActionClick: true });
```

**Key behavior:** `sidePanel.open()` requires user gesture. Tab-specific panels override global panel for that tab only.

### chrome.userScripts (Chrome 120+)

Runs scripts in USER_SCRIPT world — separate from both MAIN (page) and ISOLATED (content script).

**Requires:** `"userScripts"` permission. Until Chrome 138, user must enable "Allow User Scripts" on extension details page.

```js
await chrome.userScripts.register([{
  id: 'my-user-script',
  matches: ['https://*.example.com/*'],
  js: [{ file: 'user-scripts/inject.js' }],
  runAt: 'document_end',
  world: 'USER_SCRIPT',
}]);

// Allow messaging from USER_SCRIPT world
await chrome.userScripts.configureWorld({ messaging: true });
```

| World | Shares JS with page | chrome.* APIs | CSP from page |
|---|---|---|---|
| MAIN | Yes | No | Yes |
| ISOLATED (content script) | No | Subset | No |
| USER_SCRIPT | No | `runtime.sendMessage` only (if messaging: true) | Configurable |

### chrome.tabCapture (Chrome 116+)

Stream IDs must be used within seconds of creation; consume in offscreen document.

```js
// service worker — get stream ID
chrome.action.onClicked.addListener(async (tab) => {
  const streamId = await chrome.tabCapture.getMediaStreamId({
    targetTabId: tab.id,
    consumerTabId: tab.id,
  });
  await ensureOffscreen('src/offscreen/offscreen.html');
  chrome.runtime.sendMessage({ target: 'offscreen', type: 'startCapture', streamId });
});

// offscreen.js — create MediaStream from stream ID
chrome.runtime.onMessage.addListener(async (msg, sender, sendResponse) => {
  if (msg.target !== 'offscreen' || msg.type !== 'startCapture') return;
  const stream = await navigator.mediaDevices.getUserMedia({
    audio: { mandatory: { chromeMediaSource: 'tab', chromeMediaSourceId: msg.streamId } },
    video: { mandatory: { chromeMediaSource: 'tab', chromeMediaSourceId: msg.streamId } }
  });
  const recorder = new MediaRecorder(stream);
  // ...
});
```

**Note:** `desktopCapture.chooseDesktopMedia` stream IDs NOT usable in offscreen documents. Only `tabCapture.getMediaStreamId` streams work offscreen. For desktop capture, use visible extension page (popup, side panel, or options page).

### chrome.omnibox

Register keyword activating when typed in address bar:

```json
{ "omnibox": { "keyword": "myext" } }
```

```js
chrome.omnibox.onInputChanged.addListener((text, suggest) => {
  suggest([
    { content: `search ${text}`, description: `Search for <match>${text}</match>` },
  ]);
});

chrome.omnibox.onInputEntered.addListener((text, disposition) => {
  const url = `https://example.com/search?q=${encodeURIComponent(text)}`;
  if (disposition === 'currentTab') chrome.tabs.update({ url });
  else chrome.tabs.create({ url, active: disposition === 'newForegroundTab' });
});
```

### File System Access API

Extensions can use web File System Access API for local file ops with user gesture. **Cannot call from service worker** — use from popup, options page, or side panel only.

```js
// Open a file (requires user gesture)
async function openFile() {
  const [fileHandle] = await window.showOpenFilePicker({
    types: [{ description: 'JSON files', accept: { 'application/json': ['.json'] } }],
  });
  return JSON.parse(await (await fileHandle.getFile()).text());
}

async function saveFile(content) {
  const handle = await window.showSaveFilePicker({
    suggestedName: 'export.json',
    types: [{ description: 'JSON', accept: { 'application/json': ['.json'] } }],
  });
  const writable = await handle.createWritable();
  await writable.write(JSON.stringify(content, null, 2));
  await writable.close();
}
```

For background file access without user gesture, use native messaging instead.

---

## 6. Offscreen + Service Worker Coordination

### Pattern: SSE/LLM Stream Host in Offscreen

Canonical pattern for persistent streaming in MV3. SW sleeps; offscreen (with `AUDIO_PLAYBACK` keepalive) persists indefinitely.

```
LLM API / SSE endpoint
       ↓ fetch (streaming)
offscreen.js  ──sendMessage(chunks)──→  service-worker.js  ──postMessage──→  popup / panel
```

```js
// offscreen.js — SSE host
let sseController = null;

chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
  if (msg.target !== 'offscreen') return false;
  if (msg.type === 'startSSE') { startSSE(msg.url, msg.token); sendResponse({ ok: true }); }
  else if (msg.type === 'stopSSE') { sseController?.abort(); sendResponse({ ok: true }); }
  return true;
});

async function startSSE(url, token) {
  sseController = new AbortController();
  try {
    const resp = await fetch(url, {
      headers: { Authorization: `Bearer ${token}` },
      signal: sseController.signal,
    });
    const reader = resp.body.getReader();
    const decoder = new TextDecoder();
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      chrome.runtime.sendMessage({
        target: 'service-worker', type: 'sseChunk',
        chunk: decoder.decode(value, { stream: true }),
      });
    }
    chrome.runtime.sendMessage({ target: 'service-worker', type: 'sseDone' });
  } catch (e) {
    if (e.name !== 'AbortError')
      chrome.runtime.sendMessage({ target: 'service-worker', type: 'sseError', error: e.message });
  }
}

// service-worker.js — router
const streamPorts = new Set();
chrome.runtime.onConnect.addListener(port => {
  if (port.name === 'stream-listener') {
    streamPorts.add(port);
    port.onDisconnect.addListener(() => streamPorts.delete(port));
  }
});
chrome.runtime.onMessage.addListener((msg) => {
  if (msg.target !== 'service-worker') return false;
  if (['sseChunk', 'sseDone', 'sseError'].includes(msg.type)) {
    streamPorts.forEach(p => { try { p.postMessage(msg); } catch { streamPorts.delete(p); } });
  }
  return false;
});

// popup.js or panel.js — listener
const port = chrome.runtime.connect({ name: 'stream-listener' });
port.onMessage.addListener(msg => {
  if (msg.type === 'sseChunk') appendToUI(msg.chunk);
  if (msg.type === 'sseDone') finishUI();
});
```

### Checking Context State from Service Worker

```js
// Check if offscreen is alive
async function isOffscreenAlive(path) {
  const url = chrome.runtime.getURL(path);
  const ctxs = await chrome.runtime.getContexts({
    contextTypes: ['OFFSCREEN_DOCUMENT'], documentUrls: [url],
  });
  return ctxs.length > 0;
}

// Check other contexts
const popups = await chrome.runtime.getContexts({ contextTypes: ['POPUP'] });
const panels = await chrome.runtime.getContexts({ contextTypes: ['SIDE_PANEL'] });
```

---

## Sources

1. [chrome.offscreen API Reference](https://developer.chrome.com/docs/extensions/reference/api/offscreen)
2. [Native Messaging Chrome Docs](https://developer.chrome.com/docs/extensions/develop/concepts/native-messaging)
3. [chrome.identity API Reference](https://developer.chrome.com/docs/extensions/reference/api/identity)
4. [chrome.sidePanel API Reference](https://developer.chrome.com/docs/extensions/reference/api/sidePanel)
5. [chrome.userScripts API Reference](https://developer.chrome.com/docs/extensions/reference/api/userScripts)
6. [chrome.tabCapture API Reference](https://developer.chrome.com/docs/extensions/reference/api/tabCapture)
7. [chrome.declarativeNetRequest API Reference](https://developer.chrome.com/docs/extensions/reference/api/declarativeNetRequest)
8. [Service Worker Lifecycle](https://developer.chrome.com/docs/extensions/develop/concepts/service-workers/lifecycle)