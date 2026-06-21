<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly the standalone `websocket-extension-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`chrome-extension-expert`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: websocket-extension-patterns
version: 1.1.0
last_updated: 2026-05-29
description: >-
  WebSocket and Socket.IO patterns in Chrome MV3 extensions — service worker
  lifecycle challenges, 20-second keepalive strategy, offscreen document
  hosting, reconnection with exponential backoff, connection state machines
  persisted to chrome.storage.session, and outbound message buffering with
  inbound replay.
  TRIGGER: user implements a WebSocket or Socket.IO connection in a Chrome MV3
  extension service worker; debugs WebSocket disconnections caused by SW idle
  termination; chooses between WebSocket, SSE, polling, and offscreen document
  hosting; configures Socket.IO reconnection in an extension context; implements
  a keepalive strategy; buffers outbound messages while disconnected.
  SKIP: SSE or EventSource patterns → sse-streaming-patterns; general WebSocket
  server-side patterns without a Chrome extension context → backend-patterns.
category: developer
tags: [chrome-extension, websocket, socket-io, mv3, service-worker, real-time, offscreen-document, keepalive, reconnection]
whenToUse:
  - implementing WebSocket or Socket.IO in a Chrome MV3 extension service worker
  - debugging WebSocket connections that die every 30 seconds in an extension
  - choosing between WebSocket, offscreen document, and polling for a Chrome extension
  - configuring Socket.IO reconnection settings for an extension context
  - implementing keepalive pings to prevent service worker idle termination
  - buffering outbound messages while the WebSocket is disconnected
  - implementing inbound event replay after reconnection
  - using chrome.alarms for reconnection instead of setTimeout
whenNotToUse:
  - SSE / EventSource patterns → sse-streaming-patterns
  - WebSocket patterns outside a Chrome extension → backend-patterns or express-patterns
  - general Chrome extension architecture → chrome-dev or mv3-service-worker-expert
related_skills:
  - sse-streaming-patterns
  - mv3-service-worker-expert
  - chrome-offscreen-documents
  - chrome-storage-patterns
---

# WebSocket Extension Patterns

Expert reference for WebSocket and Socket.IO real-time connections inside Chrome MV3 extensions, covering service worker lifecycle challenges, keepalive strategies, offscreen document delegation, reconnection state machines, and message buffering.

## Overview

Chrome MV3 replaced persistent background pages with ephemeral service workers that terminate after 30 seconds of inactivity. This conflicts with WebSocket connections, which expect a long-lived runtime. Since Chrome 116, WebSocket message activity resets the SW idle timer, making in-SW WebSockets viable when combined with keepalive discipline. For use cases that cannot guarantee regular message flow, offscreen documents provide a longer-lived alternative host.

## Where to Host the Connection

| Hosting Context | Lifetime | Best For | Drawbacks |
|---|---|---|---|
| **Service worker** (Chrome 116+) | Active while messages flow every <30s | Server-push-heavy feeds, chat | Dies on idle; requires keepalive pings |
| **Offscreen document** | Until explicitly closed or Chrome reclaims it | Persistent subscriptions, audio/media streaming | One offscreen doc per extension; `chrome.runtime` messaging only |
| **Content script** | Tied to the page lifetime | Page-specific real-time features | Lost on navigation; no `chrome.storage.session` |
| **Popup / sidepanel** | Open while user interacts | Short-lived UI-driven connections | Closes when user clicks away |
| **chrome.alarms + polling** | Indefinite (alarm-driven) | Low-frequency updates (<1/min) | Not real-time; 30s minimum alarm period |

**Decision rule:** If the server sends data more than once per 25 seconds, host the WebSocket in the service worker with a keepalive. If messages are sparse or unpredictable, use an offscreen document or fall back to polling.

---

## Core Concepts

### SW Idle Timer and WebSocket Activity (Chrome 116+)

Starting Chrome 116, sending or receiving a WebSocket message resets the service worker's 30-second idle timer. Before Chrome 116, WebSocket activity did NOT prevent termination.

**Critical thresholds:**
- Idle timeout: 30 seconds of no events/API calls/WS messages
- Operation timeout: 5 minutes per single handler invocation
- Fetch timeout: 30 seconds for `fetch()` to receive the first byte

### The 5-Minute Port Reset

Chrome resets `chrome.runtime` long-lived ports every 5 minutes. If your architecture relays WebSocket data over a runtime port, handle `onDisconnect` and re-establish the port.

### Offscreen Document Constraints

- Requires `"offscreen"` permission in manifest.json
- Only ONE offscreen document per extension at a time
- Only `chrome.runtime` messaging APIs are available
- Must specify a `reason` from the allowed enum (e.g., `WORKERS`, `BLOBS`, `DOM_SCRAPING`)

---

## Service Worker WebSocket Pattern

### Manifest

```json
{
  "manifest_version": 3,
  "permissions": ["alarms", "storage"],
  "background": { "service_worker": "background.js", "type": "module" },
  "content_security_policy": {
    "extension_pages": "script-src 'self'; connect-src wss://your-server.example.com"
  }
}
```

**CSP rules:** `wss://` endpoints must be explicitly listed in `connect-src`. Wildcards work for subdomains: `wss://*.example.com`. `ws://` (unencrypted) is blocked in MV3.

### Service Worker Implementation

```javascript
// background.js
const WS_URL = 'wss://your-server.example.com/ws';
const KEEPALIVE_INTERVAL_MS = 20_000;  // safely under the 30s idle limit
const RECONNECT_ALARM = 'ws-reconnect';
const MAX_RECONNECT_ATTEMPTS = 10;

let ws = null;
let keepaliveTimer = null;

// All state persisted to chrome.storage.session — in-memory vars reset on SW restart
async function getReconnectAttempts() {
  const { wsReconnectAttempts } = await chrome.storage.session.get('wsReconnectAttempts');
  return wsReconnectAttempts || 0;
}
async function setReconnectAttempts(n) {
  await chrome.storage.session.set({ wsReconnectAttempts: n });
}
async function setState(state) {
  await chrome.storage.session.set({ wsState: state });
}

async function connect() {
  if (ws && ws.readyState === WebSocket.OPEN) return;
  await setState('connecting');

  try { ws = new WebSocket(WS_URL); }
  catch (err) { console.error('[ws] Constructor threw:', err); await scheduleReconnect(); return; }

  ws.onopen = async () => {
    await setReconnectAttempts(0);
    await setState('connected');
    startKeepalive();
    flushBuffer();
  };

  ws.onmessage = (event) => {
    // Receiving a message resets the SW idle timer (Chrome 116+)
    handleServerMessage(event.data);
  };

  ws.onclose = async (event) => {
    console.warn('[ws] Closed:', event.code, event.reason);
    cleanup();
    await scheduleReconnect();
  };

  ws.onerror = () => { /* onclose fires after onerror — reconnect handled there */ };
}

function startKeepalive() {
  clearInterval(keepaliveTimer);
  keepaliveTimer = setInterval(() => {
    if (ws?.readyState === WebSocket.OPEN) ws.send(JSON.stringify({ type: 'ping', ts: Date.now() }));
  }, KEEPALIVE_INTERVAL_MS);
}

function cleanup() { clearInterval(keepaliveTimer); keepaliveTimer = null; ws = null; }

async function scheduleReconnect() {
  const attempts = await getReconnectAttempts();
  if (attempts >= MAX_RECONNECT_ATTEMPTS) {
    await setState('disconnected');
    console.error('[ws] Max reconnect attempts reached');
    return;
  }
  await setState('reconnecting');
  await setReconnectAttempts(attempts + 1);
  // Exponential backoff: 1s, 2s, 4s ... capped at 60s
  const delaySec = Math.min(60, Math.pow(2, attempts));
  // Use chrome.alarms — setTimeout does NOT survive SW termination
  chrome.alarms.create(RECONNECT_ALARM, { delayInMinutes: delaySec / 60 });
}

chrome.alarms.onAlarm.addListener(async (alarm) => {
  if (alarm.name === RECONNECT_ALARM) await connect();
});

// Bootstrap on every SW wake
chrome.runtime.onStartup.addListener(connect);
chrome.runtime.onInstalled.addListener(connect);
(async () => {
  const { wsState } = await chrome.storage.session.get('wsState');
  if (wsState && wsState !== 'disconnected') await connect();
})();
```

### Outbound Message Buffer (volatile — in-memory)

```javascript
const messageBuffer = [];
const MAX_BUFFER_SIZE = 200;

function sendMessage(payload) {
  if (ws?.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(payload));
  } else {
    if (messageBuffer.length < MAX_BUFFER_SIZE) messageBuffer.push(payload);
    else console.warn('[ws] Buffer full; dropping message');
  }
}

function flushBuffer() {
  while (messageBuffer.length > 0 && ws?.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(messageBuffer.shift()));
  }
}
```

For a durable buffer that survives SW termination, persist to `chrome.storage.session` instead (see Outbound Buffer section below).

---

## Offscreen Document Pattern

Use when: Socket.IO client (requires DOM globals), sparse message patterns, audio/video streaming, or when you need `setTimeout`/`setInterval` without alarm workarounds.

### Manifest

```json
{
  "permissions": ["offscreen", "alarms", "storage"],
  "background": { "service_worker": "background.js", "type": "module" }
}
```

### Offscreen Document Script

```javascript
// offscreen.js — owns the WebSocket; relays to SW via chrome.runtime
const WS_URL = 'wss://your-server.example.com/ws';
const KEEPALIVE_MS = 20_000;
let ws = null, keepaliveTimer = null, reconnectAttempts = 0;

function connect() {
  ws = new WebSocket(WS_URL);

  ws.onopen = () => {
    reconnectAttempts = 0;
    chrome.runtime.sendMessage({ type: 'WS_STATE', state: 'connected' });
    keepaliveTimer = setInterval(() => {
      if (ws?.readyState === WebSocket.OPEN) ws.send(JSON.stringify({ type: 'ping' }));
    }, KEEPALIVE_MS);
  };

  ws.onmessage = (event) => {
    chrome.runtime.sendMessage({ type: 'WS_MESSAGE', data: event.data });
  };

  ws.onclose = () => {
    clearInterval(keepaliveTimer);
    reconnectAttempts++;
    chrome.runtime.sendMessage({ type: 'WS_STATE', state: 'disconnected' });
    // setTimeout is safe in offscreen doc (long-lived, unlike SW)
    const delay = Math.min(30_000, 1000 * Math.pow(2, reconnectAttempts - 1));
    setTimeout(connect, delay);
  };

  ws.onerror = () => ws?.close();
}

chrome.runtime.onMessage.addListener((msg) => {
  if (msg.type === 'WS_SEND' && ws?.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(msg.payload));
  }
});

connect();
```

### Service Worker (Offscreen Coordinator)

```javascript
// background.js
const OFFSCREEN_URL = 'offscreen.html';

async function ensureOffscreen() {
  const contexts = await chrome.runtime.getContexts({
    contextTypes: ['OFFSCREEN_DOCUMENT'],
    documentUrls: [chrome.runtime.getURL(OFFSCREEN_URL)],
  });
  if (contexts.length > 0) return;
  await chrome.offscreen.createDocument({
    url: OFFSCREEN_URL,
    reasons: ['WORKERS'],
    justification: 'Maintains persistent WebSocket connection to server',
  });
}

chrome.runtime.onMessage.addListener((msg) => {
  if (msg.type === 'WS_MESSAGE') handleServerMessage(msg.data);
  else if (msg.type === 'WS_STATE') chrome.storage.session.set({ wsState: msg.state });
});

function sendViaWebSocket(payload) {
  chrome.runtime.sendMessage({ type: 'WS_SEND', payload });
}

chrome.runtime.onStartup.addListener(ensureOffscreen);
chrome.runtime.onInstalled.addListener(ensureOffscreen);
```

---

## Socket.IO Configuration

Socket.IO requires DOM globals (`document`, `XMLHttpRequest`). It cannot run in a service worker. Host in an offscreen document.

### Socket.IO in an Offscreen Document

```javascript
import { io } from './vendor/socket.io-client.esm.js';

const socket = io('https://your-server.example.com', {
  transports: ['websocket'],   // skip HTTP polling — avoids CORS preflight issues
  reconnection: true,
  reconnectionAttempts: 15,
  reconnectionDelay: 1000,
  reconnectionDelayMax: 30000,
  randomizationFactor: 0.3,    // jitter to prevent thundering herd
  auth: (cb) => {
    chrome.runtime.sendMessage({ type: 'GET_AUTH_TOKEN' }, (response) => {
      cb({ token: response?.token || '' });
    });
  },
});

socket.on('connect', () => {
  chrome.runtime.sendMessage({ type: 'SOCKET_STATE', state: 'connected', id: socket.id });
});

socket.on('disconnect', (reason) => {
  chrome.runtime.sendMessage({ type: 'SOCKET_STATE', state: 'disconnected', reason });
  // 'io server disconnect' requires manual reconnect — all other reasons auto-reconnect
  if (reason === 'io server disconnect') socket.connect();
});

chrome.runtime.onMessage.addListener((msg) => {
  if (msg.type === 'SOCKET_EMIT') socket.emit(msg.event, msg.data);
});
```

### Key Socket.IO Settings for Extensions

| Option | Value | Why |
|---|---|---|
| `transports` | `['websocket']` | Skip polling; avoids CORS preflight in extension context |
| `reconnectionAttempts` | `10`–`20` | Enough retries for transient outages |
| `reconnectionDelayMax` | `30000` | 30s cap |
| `randomizationFactor` | `0.3` | Jitter to prevent thundering herd after server restart |
| `auth` | callback function | Fetch fresh token from SW on each connect attempt |

### Socket.IO Disconnect Reasons

| Reason | Auto-reconnects? | Action |
|---|---|---|
| `io server disconnect` | NO | Call `socket.connect()` manually |
| `io client disconnect` | NO | Intentional; no action |
| `ping timeout` | YES | Automatic |
| `transport close` | YES | Automatic |
| `transport error` | YES | Automatic |

---

## Connection State Machine

```
DISCONNECTED → CONNECTING → CONNECTED → DISCONNECTED
                               |                |
                               +→ RECONNECTING ←+
                                       |
                         (max retries) +→ FAILED
```

```javascript
// connection-state.js
export class ConnectionStateMachine {
  #state = 'disconnected';
  #listeners = new Set();

  static VALID_TRANSITIONS = {
    disconnected:  ['connecting'],
    connecting:    ['connected', 'reconnecting', 'disconnected'],
    connected:     ['disconnected', 'reconnecting'],
    reconnecting:  ['connecting', 'disconnected', 'failed'],
    failed:        ['connecting', 'disconnected'],
  };

  get state() { return this.#state; }
  get isConnected() { return this.#state === 'connected'; }

  transition(newState) {
    const allowed = ConnectionStateMachine.VALID_TRANSITIONS[this.#state];
    if (!allowed?.includes(newState)) {
      console.error(`[csm] Invalid: ${this.#state} → ${newState}`);
      return false;
    }
    const prev = this.#state;
    this.#state = newState;
    this.#listeners.forEach((fn) => fn(newState, prev));
    return true;
  }

  onChange(fn) { this.#listeners.add(fn); return () => this.#listeners.delete(fn); }
}
```

Persist state across SW restarts:

```javascript
const csm = new ConnectionStateMachine();
csm.onChange(async (newState) => {
  await chrome.storage.session.set({ connectionState: newState });
});

async function restoreState() {
  const { connectionState } = await chrome.storage.session.get('connectionState');
  if (connectionState && connectionState !== 'disconnected') await connect();
}
```

---

## Alarm-Based Polling (Alternative)

When true real-time is not required (updates every 30s–5min are acceptable):

```javascript
const POLL_ALARM = 'poll-updates';
chrome.alarms.create(POLL_ALARM, { delayInMinutes: 0.1, periodInMinutes: 1 });

chrome.alarms.onAlarm.addListener(async (alarm) => {
  if (alarm.name !== POLL_ALARM) return;
  try {
    const response = await fetch('https://api.example.com/updates?since=' + await getLastPollTs());
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const data = await response.json();
    if (data.updates?.length) { await processUpdates(data.updates); await setLastPollTs(Date.now()); }
  } catch (err) { console.error('[poll] Failed:', err); }
});
```

### Hybrid: Alarm Watchdog + WebSocket

```javascript
// Alarm acts as a watchdog — detects if the WebSocket silently died
chrome.alarms.create('ws-watchdog', { periodInMinutes: 1 });
chrome.alarms.onAlarm.addListener(async (alarm) => {
  if (alarm.name !== 'ws-watchdog') return;
  const { connectionState } = await chrome.storage.session.get('connectionState');
  if (connectionState !== 'connected') await connect();
});
```

---

## Durable Outbound Buffer

For messages that must survive SW termination:

```javascript
const BUFFER_KEY = 'ws_outbound_buffer';
const MAX_BUFFER_ITEMS = 500;

export async function enqueue(message) {
  const { [BUFFER_KEY]: buffer = [] } = await chrome.storage.session.get(BUFFER_KEY);
  if (buffer.length >= MAX_BUFFER_ITEMS) buffer.shift();
  buffer.push({ payload: message, ts: Date.now() });
  await chrome.storage.session.set({ [BUFFER_KEY]: buffer });
}

export async function flush(sendFn) {
  const { [BUFFER_KEY]: buffer = [] } = await chrome.storage.session.get(BUFFER_KEY);
  const failed = [];
  for (const item of buffer) {
    try { sendFn(item.payload); }
    catch { failed.push(item); }
  }
  await chrome.storage.session.set({ [BUFFER_KEY]: failed });
  return buffer.length - failed.length;
}
```

### Inbound Replay on Reconnect

```javascript
ws.onopen = async () => {
  const { lastEventId } = await chrome.storage.session.get('lastEventId');
  ws.send(JSON.stringify({ type: 'replay-request', since: lastEventId || null }));
};

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  if (msg.id) chrome.storage.session.set({ lastEventId: msg.id });
  handleServerMessage(msg);
};
```

---

## Anti-Patterns

| Anti-pattern | Fix |
|---|---|
| `setTimeout` for reconnection | Use `chrome.alarms` — `setTimeout` is lost when SW terminates |
| Connection state in global variables only | Persist to `chrome.storage.session` — globals reset on SW restart |
| No keepalive on quiet connections | Send a ping every 20s; server silence for 30s kills the SW |
| Socket.IO directly in the service worker | Socket.IO requires DOM globals; host in offscreen document |
| Forgetting the 5-minute port reset | Handle `port.onDisconnect` and re-establish the port |
| Creating multiple offscreen documents | Chrome enforces one per extension; check with `getContexts()` first |

---

## Troubleshooting

| Symptom | Cause | Fix |
|---|---|---|
| WebSocket connects then immediately closes | No keepalive (SW goes idle); CSP blocks WSS URL; server auth failure | Add 20s keepalive; add `connect-src wss://...` to CSP; check `event.code` + `reason` in `onclose` |
| Service worker keeps restarting | Normal idle termination (30s); unhandled exception crashes SW | Persist state to `chrome.storage.session`; add top-level error handler |
| Offscreen document silently disappears | Chrome reclaimed it; unhandled exception | Check via `chrome.runtime.getContexts()` before sending; re-create if missing |
| Socket.IO falls back to polling | `transports` not set | Set `transports: ['websocket']` |
| Messages lost during reconnection | No outbound buffer; server has no replay mechanism | Implement durable buffer; negotiate replay with event IDs |

---

## References

- [Chrome: Use WebSockets in Service Workers](https://developer.chrome.com/docs/extensions/how-to/web-platform/websockets)
- [Chrome: Extension Service Worker Lifecycle](https://developer.chrome.com/docs/extensions/develop/concepts/service-workers/lifecycle)
- [Chrome: Offscreen Documents in MV3](https://developer.chrome.com/blog/Offscreen-Documents-in-Manifest-v3)
- [Chrome: chrome.offscreen API Reference](https://developer.chrome.com/docs/extensions/reference/api/offscreen)
- [Chrome: Real-time Updates in Extensions](https://developer.chrome.com/docs/extensions/develop/concepts/real-time)
- [Chrome: chrome.alarms API Reference](https://developer.chrome.com/docs/extensions/reference/api/alarms)
- [Chrome 116 Extension Updates — WebSocket idle timer fix](https://developer.chrome.google.cn/blog/chrome-116-beta-whats-new-for-extensions)
