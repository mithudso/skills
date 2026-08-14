The error: compressed has `chrome-extension-expert` once (in banner line 1) but original has it twice — also in banner line 2 as `(\`chrome-extension-expert\`)`.

Fix: restore the inline code in line 2 of the hub-reference-banner.

The compressed file content is in the prompt. I need to output the fixed version. The only change is in the second line of the banner:

Original compressed: `> Sibling topics now reference files under hub — **not** standalone skills.`
Fixed: `> Sibling topics now reference files under hub (\`chrome-extension-expert\`) — **not** standalone skills.`

`★ Insight ─────────────────────────────────────`
Compressing inline (no file path given) — applying rules manually: drop articles/hedging/verbose connectives, compress table cells, preserve all code blocks, backticks, URLs, headings exactly.
`─────────────────────────────────────────────────`

<!-- hub-reference-banner -->
> **Reference file — part of `chrome-extension-expert` hub.** Formerly standalone `chrome-offscreen-documents` skill.
> Sibling topics now reference files under hub (`chrome-extension-expert`) — **not** standalone skills. Ignore "use the X skill" / `related_skills` / SKIP pointers naming bare sibling skills; load from `references/<name>.md` in owning hub (see hub's "Cross-hub map").

---

---
name: chrome-offscreen-documents
description: >
  Chrome MV3 offscreen docs — creating, closing, silent audio keepalive,
  DOM access from service workers, SSE/WebSocket hosting, debugging,
  migrating from MV2 background pages.
  TRIGGER: SW needs DOM API (DOMParser, clipboard, canvas, localStorage,
  audio), SW suspends mid-stream, persistent SSE or WebSocket in MV3,
  LLM streaming fetch in Chrome extension,
  migrating MV2 background page DOM code to MV3.
  SKIP: basic chrome.tabs, chrome.storage, chrome.runtime messaging (use
  chrome-dev), native messaging depth (use chrome-native-messaging), browser
  automation/page inspection (use chrome-devtools-mcp).
  See also: chrome-mv3-advanced (full advanced API incl. offscreen
  as Section 2+6), mv3-service-worker-expert (SW lifecycle), extension-e2e-testing
  (integration testing offscreen docs).
version: 1.1.0
category: developer
tags: [chrome-extension, mv3, offscreen, service-worker, dom, sse, websocket, audio, keepalive, streaming]
related_skills: [mv3-service-worker-expert, chrome-mv3-advanced, chrome-storage-patterns, extension-e2e-testing]
updated: 2026-05-29
---

# Chrome MV3 Offscreen Documents

Hidden DOM-enabled docs Chrome extensions use to access browser APIs unavailable in service workers.

## When NOT to Use

- Basic `chrome.tabs`, `chrome.storage`, `chrome.runtime` messaging or SW fundamentals → `chrome-dev`.
- Full advanced MV3 API surface (native messaging, identity/OAuth, DNR, side panel, userScripts, tabCapture) → `chrome-mv3-advanced` (embeds this as Section 2+6).
- Browser automation/page inspection (Playwright/Puppeteer) → `chrome-devtools-mcp`.

## Overview

`chrome.offscreen` (Chrome 109+, MV3 only) creates hidden HTML doc with **full DOM access**. SWs can't touch DOM APIs — offscreen docs bridge gap without visible UI.

**Permission required:** `"offscreen"` in `manifest.json`
**One instance limit:** One offscreen doc per extension per profile. Incognito gets separate instance.
**Only `chrome.runtime` APIs** available inside offscreen doc — no tabs, storage, notifications, etc.

## Manifest Setup

```json
{
  "manifest_version": 3,
  "permissions": ["offscreen"],
  "background": { "service_worker": "background.js" }
}
```

Offscreen HTML must be **bundled with extension** (no external URLs):

```html
<!-- offscreen.html -->
<!DOCTYPE html>
<script src="offscreen.js"></script>
```

## Reason Enum — Full Reference

| Reason | Use for | Lifetime |
|---|---|---|
| `AUDIO_PLAYBACK` | `<audio>` element, audio playback | **Auto-closes after 30s without audio** |
| `DOM_PARSER` | `DOMParser` API, parsing HTML/XML strings | Unlimited |
| `DOM_SCRAPING` | Embed iframe, scrape its DOM | Unlimited |
| `IFRAME_SCRIPTING` | Embed iframe, modify its content | Unlimited |
| `BLOBS` | `URL.createObjectURL()`, Blob operations | Unlimited |
| `CLIPBOARD` | `navigator.clipboard` read/write | Unlimited |
| `LOCAL_STORAGE` | `localStorage` access | Unlimited |
| `USER_MEDIA` | `getUserMedia()` — mic/camera | Unlimited |
| `DISPLAY_MEDIA` | `getDisplayMedia()` — screen capture | Unlimited |
| `WEB_RTC` | `RTCPeerConnection`, WebRTC APIs | Unlimited |
| `WORKERS` | Spawn `Worker` threads | Unlimited |
| `BATTERY_STATUS` | `navigator.getBattery()` | Unlimited |
| `MATCH_MEDIA` | `window.matchMedia()` | Unlimited |
| `GEOLOCATION` | `navigator.geolocation` | Unlimited |
| `TESTING` | Tests only | Unlimited |

**Multiple reasons:** `reasons` accepts array — declare every applicable reason:
```js
reasons: ['DOM_PARSER', 'CLIPBOARD', 'LOCAL_STORAGE']
```

## Creating and Closing

### Canonical Create Pattern (Chrome 116+)

`chrome.runtime.getContexts()` requires Chrome 116+. Earlier Chrome → SW clients fallback below.

```js
// service-worker.js
let _creatingOffscreen = null; // prevents concurrent creation race within one SW lifetime

async function ensureOffscreen(path = 'offscreen.html') {
  const url = chrome.runtime.getURL(path); // must be full chrome-extension:// URL

  const existing = await chrome.runtime.getContexts({
    contextTypes: ['OFFSCREEN_DOCUMENT'],
    documentUrls: [url],
  });
  if (existing.length > 0) return;

  // Guard against concurrent creation calls within the same SW activation
  if (_creatingOffscreen) {
    await _creatingOffscreen;
    return;
  }

  _creatingOffscreen = chrome.offscreen.createDocument({
    url,
    reasons: ['DOM_PARSER', 'CLIPBOARD'],
    justification: 'Parse HTML content and write to clipboard',
  });
  await _creatingOffscreen;
  _creatingOffscreen = null;
}

async function closeOffscreen(path = 'offscreen.html') {
  const url = chrome.runtime.getURL(path);
  const existing = await chrome.runtime.getContexts({
    contextTypes: ['OFFSCREEN_DOCUMENT'],
    documentUrls: [url],
  });
  if (existing.length > 0) {
    await chrome.offscreen.closeDocument();
  }
}
```

> `_creatingOffscreen` — module-level state, resets on SW restart (~30s idle). Prevents concurrent `createDocument` calls *within one SW lifetime*. `getContexts()` handles across-restart case.

### Fallback for Chrome < 116 (SW clients API)

```js
async function offscreenExistsFallback(path) {
  const url = chrome.runtime.getURL(path);
  const clients = await self.clients.matchAll();
  return clients.some(c => c.url === url);
}
// Replace the getContexts block with:
// if (await offscreenExistsFallback(path)) return;
```

## Communication with the Service Worker

Only `chrome.runtime` messaging available. Use `target` field to route messages.

```js
// service-worker.js
await ensureOffscreen('offscreen.html');
const result = await chrome.runtime.sendMessage({
  target: 'offscreen',
  type: 'parseHTML',
  payload: rawHtml,
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
      return true; // async — must return true to keep channel open
    }
  }
  return true;
});
```

**Ports (long-lived channel):** Use `chrome.runtime.connect()` from offscreen doc for bidirectional streaming without request/response polling.

## Audio Keepalive Pattern

`AUDIO_PLAYBACK` docs auto-close after 30s silence. Silent audio loop → persist indefinitely. Canonical pattern for SSE/LLM streams.

```html
<!-- offscreen.html -->
<audio id="keepalive" loop
  src="data:audio/wav;base64,UklGRiQAAABXQVZFZm10IBAAAAABAAEAQB8AAEAfAAABAAgAZGF0YQAAAAA=">
</audio>
<script src="offscreen.js"></script>
```

```js
// offscreen.js — start keepalive on load
document.getElementById('keepalive').play().catch(() => {
  // Autoplay may be blocked; it will start on first user interaction if needed.
  // For extensions, autoplay is typically allowed for AUDIO_PLAYBACK offscreen docs.
});
```

**Why this works:** Chrome's AUDIO_PLAYBACK timer resets while audio plays. Silent WAV loop counts as playing. Keeps offscreen doc (and any SSE/fetch streams) alive indefinitely.

## SSE and WebSocket Hosting

SWs suspend after ~30s idle. Offscreen docs persist. Host persistent connections in offscreen, forward events to SW.

### SSE Pattern

```js
// offscreen.js
let sseController = null;

chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
  if (msg.target !== 'offscreen') return false;
  if (msg.type === 'startSSE') {
    startSSE(msg.url, msg.token);
    sendResponse({ ok: true });
  } else if (msg.type === 'stopSSE') {
    sseController?.abort();
    sendResponse({ ok: true });
  }
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
        target: 'service-worker',
        type: 'sseChunk',
        chunk: decoder.decode(value, { stream: true }),
      });
    }
    chrome.runtime.sendMessage({ target: 'service-worker', type: 'sseDone' });
  } catch (e) {
    if (e.name !== 'AbortError')
      chrome.runtime.sendMessage({ target: 'service-worker', type: 'sseError', error: e.message });
  }
}
```

```js
// service-worker.js — broadcast chunks to connected popups/panels
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

// popup.js or panel.js
const port = chrome.runtime.connect({ name: 'stream-listener' });
port.onMessage.addListener(msg => {
  if (msg.type === 'sseChunk') appendChunk(msg.chunk);
  if (msg.type === 'sseDone') finalize();
});
```

### WebSocket Pattern

```js
// offscreen.js
let ws = null;

function startWebSocket(url) {
  ws = new WebSocket(url);
  ws.onopen = () => chrome.runtime.sendMessage({ target: 'service-worker', type: 'wsOpen' });
  ws.onmessage = ({ data }) => chrome.runtime.sendMessage({ target: 'service-worker', type: 'wsMessage', data });
  ws.onclose = ({ code, reason }) => {
    chrome.runtime.sendMessage({ target: 'service-worker', type: 'wsClose', code, reason });
    setTimeout(() => startWebSocket(url), 3000); // auto-reconnect
  };
}
```

**Note (Chrome 116+):** Active WebSocket in SW itself extends SW lifetime. Offscreen needed only when SW must sleep but connection must persist.

## Multipurpose Offscreen Document

One doc open at a time. Route all cases through single doc via message dispatch:

```js
// offscreen.js — single doc handling multiple concerns
chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
  if (msg.target !== 'offscreen') return false;
  switch (msg.type) {
    case 'parseHTML':     return handleParseHTML(msg, sendResponse);
    case 'clipboardRead': return handleClipboard(msg, sendResponse);
    case 'startSSE':      return handleSSE(msg, sendResponse);
    default:              return false;
  }
});
```

## Testing Offscreen Documents

Offscreen docs can't be tested directly in Node/Vitest — no DOM. Three strategies:

**1. Export logic as pure functions (unit test separately):**
```js
// offscreen-logic.js — pure, no chrome.* dependencies
export function parseHTML(html) {
  const doc = new DOMParser().parseFromString(html, 'text/html');
  return { title: doc.title, text: doc.body.innerText.trim() };
}
```

**2. Mock `chrome.offscreen` in service-worker unit tests:**
```js
globalThis.chrome = {
  offscreen: {
    createDocument: vi.fn().mockResolvedValue(undefined),
    closeDocument: vi.fn().mockResolvedValue(undefined),
  },
  runtime: {
    getContexts: vi.fn().mockResolvedValue([]),
    getURL: (p) => `chrome-extension://fakeext/${p}`,
    sendMessage: vi.fn(),
  },
};
```

**3. Integration test via Playwright extension loading:** Load unpacked extension in headless Chrome, send `chrome.runtime.sendMessage` via `page.evaluate`, assert on response. See `extension-e2e-testing` skill.

## Debugging

1. Open `chrome://extensions`, find extension card, click `offscreen.html` link — opens dedicated DevTools window with Elements, Console, Network, Sources.
2. `chrome://inspect` → "Other" → find and inspect offscreen doc.
3. **DevTools not showing offscreen:** Bug in Chrome < 116 — update Chrome. Doc closes before inspect → add `setTimeout` delay before close, or relay logs via `chrome.runtime.sendMessage`.

**Console relay pattern for hard-to-catch bugs:**
```js
// offscreen.js
const origLog = console.log.bind(console);
console.log = (...args) => {
  origLog(...args);
  chrome.runtime.sendMessage({ type: 'offscreenLog', args: args.map(String) });
};
```

## Migration from MV2 Background Pages

| MV2 Background page | MV3 replacement |
|---|---|
| `window`, `document`, DOM access | Offscreen doc with appropriate reason |
| `localStorage` read/write | Offscreen (`LOCAL_STORAGE`) or migrate to `chrome.storage.local` |
| `<audio>` playback | Offscreen (`AUDIO_PLAYBACK`) |
| Persistent WebSocket/SSE | Offscreen (`WEB_RTC` or `WORKERS` + keepalive) |
| `DOMParser` / innerHTML scraping | Offscreen (`DOM_PARSER`) |
| Clipboard read/write | Offscreen (`CLIPBOARD`) |
| Long-running fetch streams | Offscreen (`AUDIO_PLAYBACK` keepalive) |
| Global in-memory state | `chrome.storage.session` (cleared on Chrome restart) |

## Common Mistakes

| Mistake | Fix |
|---|---|
| Calling `createDocument` without checking if one exists | Always call `ensureOffscreen()` with `getContexts` check first |
| Relative path in `createDocument` URL | Use `chrome.runtime.getURL('offscreen.html')` |
| `createDocument` called concurrently from multiple event handlers | Use `_creatingOffscreen` guard promise |
| `AUDIO_PLAYBACK` doc closes unexpectedly | Silent audio loop must play; check autoplay not blocked |
| Calling `chrome.storage` / `chrome.tabs` from offscreen | Only `chrome.runtime` available; route via SW |
| `return false` from async message handler in offscreen | Return `true` to keep channel open for async `sendResponse` |
| `chrome.offscreen.hasDocument()` not found | Removed after early experimental; use `runtime.getContexts()` |
| Trying to host `tabCapture` stream in offscreen | Only `tabCapture.getMediaStreamId` streams work; `desktopCapture` streams do not |
| Calling `closeDocument()` when no doc open | Throws — guard with `isOffscreenAlive()` or wrap in try/catch |

## Quick Decision Guide

| Need | Reason(s) | Notes |
|---|---|---|
| Parse HTML string | `DOM_PARSER` | `new DOMParser()` |
| Write to clipboard | `CLIPBOARD` | `navigator.clipboard.writeText()` |
| Play audio silently | `AUDIO_PLAYBACK` | 30s auto-close; silent loop extends |
| Persistent SSE connection | `AUDIO_PLAYBACK` (with keepalive) | SW sleeps; offscreen persists |
| Persistent WebSocket | `WEB_RTC` or `WORKERS` | SW if Chrome 116+ WS keepalive suffices |
| LLM streaming fetch | `AUDIO_PLAYBACK` (with keepalive) | Offscreen owns stream; SW routes chunks |
| Scrape embedded iframe | `DOM_SCRAPING` | Embed iframe, read its DOM |
| `localStorage` (MV2 migration) | `LOCAL_STORAGE` | Migrate data then switch to `chrome.storage.local` |
| Blob/object URLs | `BLOBS` | `URL.createObjectURL()` |
| Spawn Worker threads | `WORKERS` | Background compute without blocking |

## Sources

- [chrome.offscreen API Reference](https://developer.chrome.com/docs/extensions/reference/api/offscreen)
- [Offscreen Documents in Manifest V3](https://developer.chrome.com/blog/Offscreen-Documents-in-Manifest-v3)
- [Migrate to a Service Worker](https://developer.chrome.com/docs/extensions/develop/migrate/to-service-workers)
- [Service Worker Lifecycle](https://developer.chrome.com/docs/extensions/develop/concepts/service-workers/lifecycle)