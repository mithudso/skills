<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly standalone `extension-message-bridge` skill.
> Sibling topics now reference files under hubs (`chrome-extension-expert`) — **not** standalone skills. Ignore "use X skill" / `related_skills` / SKIP pointers naming bare sibling skill; load `references/<name>.md` from owning hub (see hub's "Cross-hub map").

---

---
name: extension-message-bridge
version: 1.1.0
category: developer
tags: [chrome-extension, messaging, postMessage, runtime, mv3, service-worker, content-script, iframe, ports]
description: >-
  Chrome MV3 extension cross-context messaging — chrome.runtime.sendMessage, chrome.tabs.sendMessage,
  chrome.runtime.connect (ports), window.postMessage bridging, CustomEvent dispatch, and MV3 service
  worker message resilience patterns.
  TRIGGER: implementing cross-context communication in a Chrome extension, debugging "Receiving end does
  not exist" or "Extension context invalidated" errors, designing iframe↔parent postMessage bridges,
  long-lived port patterns, or MAIN-world↔ISOLATED-world message bridging.
  SKIP: extension E2E test setup (use extension-e2e-testing); non-Chrome messaging (WebSockets, SSE,
  fetch) not crossing extension context boundaries.
triggers:
  - "chrome.runtime.sendMessage not receiving a response"
  - "Receiving end does not exist error in Chrome extension"
  - "Extension context invalidated after update"
  - "How do I send a message from popup to service worker?"
  - "Long-lived port disconnects when service worker terminates"
  - "postMessage bridge between iframe panel and content script"
  - "MAIN world script relay to content script"
  - "Broadcast a message to all open extension pages"
  - "Service worker message channel stays open async"
  - "return true not working in onMessage listener"
related_skills:
  - mv3-service-worker-expert
  - chrome-storage-patterns
  - chrome-dev
  - extension-e2e-testing
  - websocket-extension-patterns
---

# Extension Message Bridge Patterns

## Channel quick reference

| Channel | API | Use when | Wakes suspended SW? |
|---------|-----|----------|---------------------|
| One-shot (ext-internal) | `chrome.runtime.sendMessage` | Request/response to SW | Yes |
| One-shot (to tab) | `chrome.tabs.sendMessage(tabId)` | SW → tab content script | N/A |
| Long-lived port | `chrome.runtime.connect` | Streaming / multi-step | Yes (on connect) |
| Cross-frame | `window.postMessage` | iframe↔parent, MAIN↔ISOLATED | N/A |
| Same-document | `CustomEvent` on `document` | CS module↔CS module | N/A |
| Storage broadcast | `chrome.storage.onChanged` | Settings propagation | N/A |

**Golden rules:**
1. Prefer one-shot `sendMessage` over ports — wakes terminated SW automatically.
2. Never use `async` directly on `onMessage` listener — wrap in sync outer + async IIFE.
3. Always validate `event.data.source` on `postMessage` listeners.
4. Always catch "Receiving end does not exist" on broadcasts.
5. Handle `port.onDisconnect` — ports die when SW terminates.

## Context topology

```
Host page (MAIN world)
     | window.postMessage
     v
Content Script (ISOLATED world)
     |                              |
     | chrome.runtime.sendMessage   | CustomEvent on document
     v                              v
Service Worker  <-- chrome.runtime.sendMessage --  Popup / Options / Dashboard
     |
     | chrome.tabs.sendMessage(tabId)
     v
Content Script (specific tab)
     |
     | window.postMessage (parent ↔ iframe)
     v
Panel iframe (chrome-extension:// origin)
     |
     | chrome.runtime.sendMessage
     v
Service Worker
```

## Pattern 1: One-shot request/response (any context → SW)

```js
// --- In popup.js, options.js, dashboard.js, content script ---

async function sendMessage(type, extra = {}) {
  if (!globalThis.chrome?.runtime?.sendMessage) return null;
  try {
    return await globalThis.chrome.runtime.sendMessage({ type, ...extra });
  } catch (err) {
    console.warn('[sendMessage] failed:', type, err?.message);
    return null;
  }
}

const options = await sendMessage('MCA_GET_OPTIONS');
const state   = await sendMessage('MCA_GET_TRACKING_STATE');
```

## Pattern 2: SW dispatch table (async handler)

**Only correct pattern** for async `onMessage` handlers in MV3:

```js
// --- In service-worker.js ---

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  const type = String(message?.type || '');
  if (!type.startsWith('MCA_')) return undefined;  // ignore foreign messages

  // CORRECT: sync outer returns true; async work in IIFE
  (async () => {
    switch (type) {
      case 'MCA_GET_OPTIONS':        return getOptions();
      case 'MCA_SAVE_OPTIONS':       return saveOptions(message.options || {});
      case 'MCA_GET_TRACKING_STATE': return trackingStateForUi(await getTrackingState());
      default:                       return { ok: false, error: `unknown: ${type}` };
    }
  })()
    .then(result => sendResponse({ ok: true, result }))
    .catch(err  => sendResponse({ ok: false, error: err?.message || String(err) }));

  return true;  // CRITICAL: keep channel open for async sendResponse
});
```

**Why `async` listeners break silently:**
```js
// WRONG — async returns a Promise, not true; channel closes immediately
chrome.runtime.onMessage.addListener(async (msg, sender, sendResponse) => {
  const result = await doWork();
  sendResponse(result);  // channel already closed — this is a no-op
  return true;           // ignored — Promise was returned, not true
});
```

## Pattern 3: SW → content script (tabs.sendMessage)

```js
async function sendOverlayCommand(tabId, action) {
  try {
    await chrome.tabs.sendMessage(tabId, { type: 'MCA_OVERLAY_COMMAND', action });
  } catch (err) {
    if (/Receiving end does not exist/i.test(String(err?.message))) {
      // Content script not injected — re-inject and retry once
      await chrome.scripting.executeScript({
        target: { tabId },
        files: ['src/content/case-overlay.js'],
      });
      await chrome.tabs.sendMessage(tabId, { type: 'MCA_OVERLAY_COMMAND', action });
    }
  }
}
```

## Pattern 4: Best-effort broadcast (SW → all open extension pages)

```js
function broadcastMessage(message) {
  const type = message?.type || 'unknown';
  const p = chrome.runtime.sendMessage(message);
  if (p && typeof p.catch === 'function') {
    p.catch(err => {
      if (/Receiving end does not exist/i.test(String(err?.message))) {
        console.debug('[broadcast] no listener for', type);  // expected when no pages open
      } else {
        console.warn('[broadcast] error for', type, err);
      }
    });
  }
}
```

## Pattern 5: CustomEvent (same-document CS modules)

```js
// Publisher (hub-extractor.js)
document.dispatchEvent(new CustomEvent('mca:context-updated', {
  detail: { caseNumber: context.case_number || '', pageUrl: context.page_url },
}));

// Subscriber (case-overlay.js)
document.addEventListener('mca:context-updated', (e) => {
  updateOverlayHeader(e.detail.caseNumber);
});
```

No serialization overhead; `detail` holds live object refs. Cannot cross context boundaries.

## Pattern 6: iframe ↔ parent postMessage bridge

```js
// --- case-overlay.js (parent content script) ---

function postContextToPanel() {
  overlay.iframe.contentWindow?.postMessage({
    source: 'mdb-case-assistant-overlay',
    type: 'MCA_PANEL_CONTEXT',
    payload: buildContextPayload(),
  }, '*');  // '*' acceptable — both sides validate source tag
}

window.addEventListener('message', (event) => {
  const data = event.data || {};
  if (data.source !== 'mdb-case-assistant-panel') return;
  switch (data.type) {
    case 'MCA_PANEL_READY':           postContextToPanel(); break;
    case 'MCA_PANEL_HIDE_OVERLAY':    hideOverlay(); break;
    case 'MCA_PANEL_UPDATE_SEVERITY': updateSeverityBadge(data.severity); break;
  }
});
```

```js
// --- panel.js (chrome-extension:// iframe) ---

window.parent?.postMessage({ source: 'mdb-case-assistant-panel', type: 'MCA_PANEL_READY' }, '*');

window.addEventListener('message', (event) => {
  const data = event.data || {};
  if (data.source !== 'mdb-case-assistant-overlay') return;
  if (data.type === 'MCA_PANEL_CONTEXT') renderCaseData(data.payload);
});
```

**Security note on `'*'` targetOrigin:** OK for non-sensitive UI coordination when both sides validate `data.source`. For tokens/PII, check `event.origin` against known `chrome-extension://<id>` origin.

## Pattern 7: Long-lived port (chrome.runtime.connect)

```js
// --- Content script (initiator) ---
const port = chrome.runtime.connect({ name: 'realtime-sync' });
port.onMessage.addListener(msg => { if (msg.type === 'CASE_UPDATE') updateLocalState(msg.data); });
port.onDisconnect.addListener(() => scheduleReconnect());  // SW may have terminated
port.postMessage({ type: 'SUBSCRIBE', caseNumbers: ['01234567'] });

// --- Service worker (acceptor) ---
const activePorts = new Map();
chrome.runtime.onConnect.addListener((port) => {
  if (port.name !== 'realtime-sync') return;
  const id = `${port.sender?.tab?.id}-${Date.now()}`;
  activePorts.set(id, port);
  port.onMessage.addListener(msg => { if (msg.type === 'SUBSCRIBE') registerSub(id, msg.caseNumbers); });
  port.onDisconnect.addListener(() => { activePorts.delete(id); removeSub(id); });
});
```

## Pattern 8: MAIN-world → content script bridge

```js
// --- MAIN-world injected script ---
const origSend = WebSocket.prototype.send;
WebSocket.prototype.send = function(data) {
  window.postMessage({ source: 'mca-main-world', type: 'WS_FRAME', payload: typeof data === 'string' ? data : '[binary]' }, '*');
  return origSend.call(this, data);
};

// --- Content script (ISOLATED world) ---
window.addEventListener('message', (event) => {
  if (event.source !== window) return;              // same-frame only
  if (event.data?.source !== 'mca-main-world') return;
  if (event.data.type === 'WS_FRAME') {
    chrome.runtime.sendMessage({ type: 'MCA_REALTIME_BRIDGE_EVENT', payload: event.data.payload });
  }
});
```

## MV3 service worker considerations

**One-shot messages wake SW; ports don't** — SW terminates while port open → port disconnects. Always implement `onDisconnect` + reconnect.

**Keepalive strategies:**
```js
// Strategy 1: chrome.alarms (min 30s period)
chrome.alarms.create('keepalive', { periodInMinutes: 0.5 });
chrome.alarms.onAlarm.addListener(alarm => { if (alarm.name === 'keepalive') { /* no-op */ } });

// Strategy 2: Offscreen document heartbeat (Chrome 109+)
// offscreen.html sends chrome.runtime.sendMessage every 25s

// Strategy 3: Active WebSocket (Chrome 116+) resets idle timer on each frame
```

**Re-injection after extension update:**
```js
chrome.runtime.onInstalled.addListener(async ({ reason }) => {
  if (reason === 'update' || reason === 'install') {
    const tabs = await chrome.tabs.query({ url: ['*://support.mongodb.com/*'] });
    for (const tab of tabs) {
      try {
        await chrome.scripting.executeScript({ target: { tabId: tab.id }, files: ['src/content/hub-extractor.js'] });
      } catch { /* tab may be discarded */ }
    }
  }
});
```

**Guard against "Extension context invalidated":**
```js
let contextAlive = true;
async function sendRuntimeMessage(message) {
  if (!contextAlive) return null;
  try {
    return await chrome.runtime?.sendMessage?.(message);
  } catch (err) {
    if (/Extension context invalidated/i.test(String(err?.message))) {
      contextAlive = false;
      return null;
    }
    throw err;
  }
}
```

## Anti-patterns

| Anti-pattern | Fix |
|-------------|-----|
| `async` listener function | Sync outer + async IIFE inside; `return true` |
| Missing `return true` for async response | Add `return true` before the async IIFE |
| Broadcasting without catching errors | Always `.catch()` broadcast `sendMessage` |
| `postMessage` without source validation | Check `event.data.source` and/or `event.origin` |
| Port without `onDisconnect` handler | Add reconnect logic in `onDisconnect` |
| `tabs.sendMessage` without injection fallback | Catch and re-inject with `scripting.executeScript` |
| Posting to iframe before it sends `READY` | Wait for iframe's `MCA_PANEL_READY` message first |

## Troubleshooting

**"Could not establish connection. Receiving end does not exist."**
- Broadcast: catch + log at debug level — expected when no pages open.
- Tab-targeted: content script not injected. Catch and re-inject.
- Check `content_scripts.matches` covers target URL.

**"Extension context invalidated"** — Extension updated/reloaded. Set flag on first occurrence; stop further `chrome.runtime` calls. Optionally show "please refresh" banner.

**`sendResponse` ignored / response is `undefined`**
- Listener used `async` function (returns Promise, not `true`).
- `return true` placed after `await` (too late — channel already closed).
- Fix: sync outer + async IIFE (Pattern 2).

**Port disconnects immediately** — `onConnect` handler didn't store port ref. Store in Map before handler returns.

**Messages arrive out of order** — Multiple rapid `sendMessage` calls independent. Add sequence numbers or use long-lived port (ordering guaranteed).