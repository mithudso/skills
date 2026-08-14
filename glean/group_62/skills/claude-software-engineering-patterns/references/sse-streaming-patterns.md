<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `sse-streaming-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: sse-streaming-patterns
version: 1.1.0
last_updated: 2026-05-29
description: >-
  Server-Sent Events (SSE) patterns — EventSource API, fetch+ReadableStream
  SSE client, Node.js server implementation, reconnection with Last-Event-ID
  replay, backpressure and flow control, browser connection limits and HTTP/2
  multiplexing, Chrome MV3 extension SSE patterns (offscreen document), and
  SSE vs WebSocket vs polling decision guidance.
  TRIGGER: user implements or reviews an SSE endpoint or client; uses
  EventSource; builds a text/event-stream server; streams LLM/AI token output;
  handles SSE reconnection or Last-Event-ID; hits browser connection limits;
  needs to choose between SSE, WebSocket, and polling; routes SSE through a
  Chrome extension offscreen document.
  SKIP: WebSocket or Socket.IO in a Chrome extension → websocket-extension-patterns;
  general real-time bidirectional communication → websocket-extension-patterns.
category: developer
tags: [sse, server-sent-events, eventsource, streaming, backpressure, real-time, http2, readablestream, fetch-streaming, chrome-extension]
whenToUse:
  - implementing an SSE server endpoint (Node.js/Express)
  - building an EventSource or fetch+ReadableStream SSE client
  - streaming LLM or AI token output to a browser
  - handling SSE reconnection, Last-Event-ID replay, or retry intervals
  - hitting the 6-connection browser limit and needing HTTP/2 or multiplexing solutions
  - choosing between SSE, WebSocket, and polling for a real-time feature
  - routing SSE through a Chrome MV3 extension offscreen document
  - debugging SSE that arrives in bursts (nginx buffering), drops every 60s, or loses events on reconnect
whenNotToUse:
  - WebSocket or Socket.IO in a Chrome MV3 extension → websocket-extension-patterns
  - bidirectional real-time communication (chat, gaming, collaborative editing) → websocket-extension-patterns
  - binary data streaming → WebSocket directly
related_skills:
  - websocket-extension-patterns
  - mv3-service-worker-expert
  - chrome-offscreen-documents
---

# SSE Streaming Patterns

## Overview

Server-Sent Events (SSE) provide a standardized HTTP-based mechanism for unidirectional server-to-client streaming. The server holds an HTTP connection open and writes plain-text frames (`text/event-stream`) as new data becomes available. The browser's `EventSource` API handles connection management, automatic reconnection, and event dispatch.

SSE is the standard pattern for streaming AI/LLM token responses, live dashboards, notification feeds, and any scenario where the server pushes updates and the client renders them. With HTTP/2 multiplexing, the historical 6-connection-per-domain limit is eliminated.

**When to use SSE vs alternatives:**

| Use SSE | Use WebSocket | Use Polling |
|---|---|---|
| Server pushes to client (one-way) | Client sends frequent messages during stream | Low frequency, simple state check |
| Streaming LLM/AI tokens | Chat, gaming, collaborative editing | < 1 update per minute |
| Auto-reconnect + event replay needed | Binary data streaming required | Infrastructure can't handle persistent connections |
| Must stay on standard HTTP | Sub-millisecond latency required | |

**Gray area:** If the client only sends occasional signals (cancel, pause), SSE + a regular HTTP POST for the reverse channel avoids WebSocket complexity.

---

## SSE Protocol

### Wire format

```
field: value\n
field: value\n
\n
```

### Fields

| Field | Purpose | Default |
|---|---|---|
| `data` | Event payload (multiple `data:` lines joined with `\n`) | (required) |
| `event` | Named event type (dispatched to `addEventListener(type, ...)`) | `"message"` |
| `id` | Sets `lastEventId`; sent as `Last-Event-ID` on reconnect | (none) |
| `retry` | Reconnection interval in milliseconds | ~3000 ms |

**Comment lines** (`:`) act as keepalive heartbeats:
```
: heartbeat\n
\n
```

### Example stream

```http
HTTP/1.1 200 OK
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
X-Accel-Buffering: no

retry: 5000

id: 1
event: update
data: {"temperature": 72.4, "unit": "F"}

id: 2
event: alert
data: {"message": "Threshold exceeded"}

: keepalive
```

---

## EventSource API

```js
const es = new EventSource('/api/stream');                   // GET only
const es = new EventSource('/api/stream', { withCredentials: true }); // with cookies

// Default event (event field absent or "message")
es.onmessage = (e) => { console.log('data:', e.data); };

// Named events — must use addEventListener, not onmessage
es.addEventListener('update', (e) => {
  const payload = JSON.parse(e.data);
});

es.onerror = (e) => {
  if (es.readyState === EventSource.CLOSED) {
    console.error('Permanently closed');
  }
  // else: reconnecting automatically
};

es.close(); // Sets readyState to CLOSED, stops reconnection
```

### EventSource limitations

| Limitation | Workaround |
|---|---|
| GET only, no custom headers | Use fetch + ReadableStream (see below) |
| No request body | Use fetch + ReadableStream |
| No error detail (status code) | Use fetch + ReadableStream |
| Text only, no binary | Use WebSocket |

---

## Reconnection Patterns

### Built-in auto-reconnect + Last-Event-ID

When the connection drops, `EventSource` auto-reconnects and sends `Last-Event-ID` with the most recent `id:` field. The server uses this to replay missed events.

**Critical rule:** Always assign an `id` to every event. Without `id`, replay silently breaks.

```js
// Server-side: replay missed events from Last-Event-ID
app.get('/api/stream', (req, res) => {
  const lastId = req.headers['last-event-id'];
  const startFrom = lastId ? parseInt(lastId, 10) + 1 : 0;

  res.writeHead(200, {
    'Content-Type': 'text/event-stream',
    'Cache-Control': 'no-cache',
    'Connection': 'keep-alive',
    'X-Accel-Buffering': 'no',
  });

  for (const evt of eventStore.getEventsFrom(startFrom)) {
    res.write(`id: ${evt.id}\ndata: ${JSON.stringify(evt.data)}\n\n`);
  }
  // continue streaming new events...
});
```

### Exponential backoff for persistent failures

EventSource auto-reconnect uses a fixed interval. For extended server downtime, implement manual exponential backoff:

```js
function createResilientStream(url, options = {}) {
  const { maxRetries = 10, baseDelay = 1000, maxDelay = 30000, jitter = true } = options;
  let retries = 0;
  let es = null;

  function connect() {
    es = new EventSource(url);
    es.onopen = () => { retries = 0; };
    es.onerror = () => {
      if (es.readyState === EventSource.CLOSED) return;
      es.close();
      if (retries >= maxRetries) { console.error('Max retries reached'); return; }
      let delay = Math.min(baseDelay * 2 ** retries, maxDelay);
      if (jitter) delay += Math.random() * delay * 0.1;
      retries++;
      setTimeout(connect, delay);
    };
    es.onmessage = (e) => { /* handle events */ };
    return es;
  }
  return { connect, close: () => es?.close() };
}
```

### Heartbeat timeout detection

```js
function monitoredStream(url, heartbeatTimeoutMs = 45000) {
  let heartbeatTimer = null;
  const es = new EventSource(url);

  function resetHeartbeat() {
    clearTimeout(heartbeatTimer);
    heartbeatTimer = setTimeout(() => {
      console.warn('No heartbeat received, reconnecting...');
      es.close();
      // reconnect logic here
    }, heartbeatTimeoutMs);
  }

  es.onopen = resetHeartbeat;
  es.onmessage = (e) => { resetHeartbeat(); /* process event */ };
  es.addEventListener('heartbeat', () => resetHeartbeat());
  return es;
}
```

---

## Fetch + ReadableStream (Modern Alternative)

Use when you need custom headers (`Authorization`), POST requests, or fine-grained stream control.

### Fetch-based SSE client

```js
async function fetchSSE(url, options = {}) {
  const { headers = {}, body, signal, onEvent } = options;

  const response = await fetch(url, {
    method: body ? 'POST' : 'GET',
    headers: { Accept: 'text/event-stream', ...(body ? { 'Content-Type': 'application/json' } : {}), ...headers },
    body: body ? JSON.stringify(body) : undefined,
    signal,
  });

  if (!response.ok) throw new Error(`SSE request failed: ${response.status}`);

  const reader = response.body.pipeThrough(new TextDecoderStream()).getReader();
  let buffer = '';

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;
    buffer += value;
    const events = buffer.split('\n\n');
    buffer = events.pop();
    for (const eventStr of events) {
      if (!eventStr.trim()) continue;
      const event = parseSSEEvent(eventStr);
      if (event) onEvent(event);
    }
  }
}

function parseSSEEvent(raw) {
  const event = { data: '', event: 'message', id: null, retry: null };
  for (const line of raw.split('\n')) {
    if (line.startsWith(':')) continue;
    const colonIdx = line.indexOf(':');
    if (colonIdx === -1) continue;
    const field = line.slice(0, colonIdx);
    const value = line.slice(colonIdx + 1).trimStart();
    switch (field) {
      case 'data': event.data += (event.data ? '\n' : '') + value; break;
      case 'event': event.event = value; break;
      case 'id': event.id = value; break;
      case 'retry': event.retry = parseInt(value, 10); break;
    }
  }
  return event.data ? event : null;
}
```

### Usage with authentication

```js
const controller = new AbortController();

fetchSSE('/api/ai/stream', {
  headers: { Authorization: `Bearer ${token}` },
  body: { prompt: 'Explain SSE' },
  signal: controller.signal,
  onEvent: (event) => {
    if (event.event === 'token') process.stdout.write(event.data);
    else if (event.event === 'done') console.log('\n[Stream complete]');
  },
}).catch((err) => { if (err.name !== 'AbortError') console.error(err); });

// Cancel: controller.abort();
```

### Async generator pattern

```js
async function* sseStream(url, fetchOptions = {}) {
  const response = await fetch(url, { ...fetchOptions, headers: { Accept: 'text/event-stream', ...fetchOptions.headers } });
  const reader = response.body.pipeThrough(new TextDecoderStream()).getReader();
  let buffer = '';
  try {
    while (true) {
      const { done, value } = await reader.read();
      if (done) return;
      buffer += value;
      const parts = buffer.split('\n\n');
      buffer = parts.pop();
      for (const part of parts) { const event = parseSSEEvent(part); if (event) yield event; }
    }
  } finally { reader.releaseLock(); }
}

for await (const event of sseStream('/api/stream')) {
  console.log(event.event, event.data);
}
```

---

## Server Implementation (Node.js/Express)

### Minimal Express SSE endpoint

```js
import express from 'express';
import crypto from 'crypto';

const app = express();
const clients = new Map();
let eventId = 0;

app.get('/api/events', (req, res) => {
  res.writeHead(200, {
    'Content-Type': 'text/event-stream',
    'Cache-Control': 'no-cache',
    'Connection': 'keep-alive',
    'X-Accel-Buffering': 'no',     // disable nginx buffering
  });

  // Replay missed events
  const lastId = req.headers['last-event-id'];
  if (lastId) {
    for (const evt of getEventsSince(parseInt(lastId, 10))) {
      res.write(`id: ${evt.id}\nevent: ${evt.type}\ndata: ${evt.payload}\n\n`);
    }
  }

  const clientId = crypto.randomUUID();
  clients.set(clientId, res);

  // Heartbeat every 25 seconds (proxy idle timeouts are typically 60s)
  const heartbeat = setInterval(() => res.write(': heartbeat\n\n'), 25000);

  req.on('close', () => { clearInterval(heartbeat); clients.delete(clientId); });
});

function broadcast(eventType, data) {
  eventId++;
  const frame = `id: ${eventId}\nevent: ${eventType}\ndata: ${JSON.stringify(data)}\n\n`;
  for (const [clientId, res] of clients) {
    const ok = res.write(frame);
    if (!ok) console.warn(`Backpressure on client ${clientId}`);
  }
}
```

### Required headers

| Header | Value | Purpose |
|---|---|---|
| `Content-Type` | `text/event-stream` | Triggers SSE parsing in browser |
| `Cache-Control` | `no-cache` | Prevents caching of the stream |
| `Connection` | `keep-alive` | Keeps TCP connection open (HTTP/1.1) |
| `X-Accel-Buffering` | `no` | Disables nginx proxy buffering |

### Disable compression on SSE routes

Gzip buffers output, delaying events. Disable compression specifically on SSE routes:

```js
import compression from 'compression';
app.use(compression({
  filter: (req, res) => {
    if (req.path === '/api/events') return false;
    return compression.filter(req, res);
  },
}));
```

---

## Backpressure and Flow Control

SSE has no built-in client-to-server backpressure. The server must handle it by checking `res.write()` return values.

### Server-side flow control (Node.js)

```js
async function sendEvent(res, id, data) {
  const chunk = `id: ${id}\ndata: ${JSON.stringify(data)}\n\n`;
  const canContinue = res.write(chunk);
  if (!canContinue) {
    await new Promise((resolve) => res.once('drain', resolve));
  }
}
```

### Buffering strategy comparison

| Strategy | When to use |
|---|---|
| Bounded ring buffer (drop oldest) | Dashboard/metrics where latest value wins |
| Drop on backpressure | High-frequency telemetry |
| Rate limiting per client | Multi-tenant systems |
| Aggregation (batch small updates) | Frequent small updates |
| Unbounded buffer | Never in production |

---

## Browser Limitations and Workarounds

### The 6-connection limit (HTTP/1.1)

Browsers limit concurrent HTTP/1.1 connections to ~6 per domain. Each `EventSource` consumes one. The fix is HTTP/2 multiplexing — it supports ~100+ concurrent streams over a single TCP connection, eliminating the bottleneck entirely.

**Action:** Serve your SSE endpoint over HTTP/2.

### Channel multiplexing (fallback for HTTP/1.1)

Use a single SSE connection that carries multiple logical channels via named events:

```js
// Server
function sendToChannel(res, channel, data, id) {
  res.write(`id: ${id}\nevent: ${channel}\ndata: ${JSON.stringify(data)}\n\n`);
}

// Client
const es = new EventSource('/api/stream');
es.addEventListener('notifications', (e) => handleNotification(JSON.parse(e.data)));
es.addEventListener('metrics',       (e) => handleMetrics(JSON.parse(e.data)));
es.addEventListener('alerts',        (e) => handleAlert(JSON.parse(e.data)));
```

---

## Chrome Extension Context

### The service worker problem

MV3 extension service workers terminate after ~30 seconds of inactivity. An SSE connection opened directly in the service worker will be killed. **Do NOT open SSE connections from the service worker.**

### Pattern 1: Offscreen document (recommended)

```js
// service-worker.js
async function ensureOffscreen() {
  const existing = await chrome.offscreen.hasDocument();
  if (!existing) {
    await chrome.offscreen.createDocument({
      url: 'offscreen/offscreen.html',
      reasons: ['WORKERS'],
      justification: 'Maintains persistent SSE connection for real-time updates',
    });
  }
}

chrome.runtime.onMessage.addListener((msg) => {
  if (msg.type === 'SSE_EVENT') handleSSEEvent(msg.payload);
});
```

```js
// offscreen/offscreen.js
const es = new EventSource('https://api.example.com/events');
es.onmessage = (e) => {
  chrome.runtime.sendMessage({ type: 'SSE_EVENT', payload: JSON.parse(e.data) });
};
```

### Pattern 2: Content script (same-origin pages only)

```js
// content-script.js
const es = new EventSource('/api/live-updates');
es.onmessage = (e) => {
  chrome.runtime.sendMessage({ type: 'LIVE_UPDATE', data: JSON.parse(e.data) });
};
```

### Pattern choice

| Pattern | Persistent | Custom headers | When to use |
|---|---|---|---|
| Offscreen + EventSource | Yes | No | Standard SSE with cookie auth |
| Offscreen + fetch | Yes | Yes | Needs Authorization header |
| Content script | Page-bound | No | Same-origin SSE, page-specific |
| Service worker | No (30s) | — | Never for long-lived SSE |

---

## Security

### Authentication options

1. **Cookie auth** — `withCredentials: true` for cross-origin cookies
2. **Token in URL** — `new EventSource('/api/stream?token=xxx')` — avoid for sensitive tokens (appears in server logs)
3. **Fetch-based SSE** — `Authorization: Bearer` header (recommended)

### CORS for SSE

```
Access-Control-Allow-Origin: https://your-client.com
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Last-Event-ID
```

Avoid `Access-Control-Allow-Origin: *` when using credentials — browsers reject this combination.

### DoS prevention

- Cap concurrent SSE connections per IP/user
- Enforce authentication before accepting SSE connections
- Set maximum connection duration (e.g., 24 hours) and force reconnect

---

## Anti-Patterns

| Anti-pattern | Fix |
|---|---|
| No `id:` field on events | Add monotonic `id:` to every event — required for replay |
| No heartbeat | Send comment heartbeat every 15–30 seconds; proxies close idle connections after 60–120s |
| Compression on SSE routes | Gzip buffers output; disable compression for SSE endpoints |
| Ignoring `res.write()` return `false` | Check return value and await `drain` event to prevent memory leaks |
| Opening EventSource in MV3 service worker | Connection dies in ~30s; use offscreen document |
| Multiple EventSource connections where one stream suffices | Each consumes a browser slot; multiplex via named events |
| Using EventSource when status code awareness is needed | `onerror` hides HTTP status; use fetch-based SSE |

---

## Troubleshooting

| Symptom | Likely Cause | Fix |
|---|---|---|
| Events arrive in bursts, not real-time | Proxy buffering (nginx, Cloudflare) | Set `X-Accel-Buffering: no`, disable compression |
| Connection drops every 60s | Proxy idle timeout | Send heartbeat comments every 15–30s |
| No reconnection after disconnect | Server returned 4xx/5xx | EventSource doesn't reconnect on HTTP errors; use fetch-based SSE |
| Events missing after reconnect | No `id:` field on events | Add monotonic `id:` to every event |
| 6th tab blocks all requests | HTTP/1.1 connection limit | Migrate to HTTP/2 or multiplex channels |
| Service worker SSE dies in 30s | MV3 lifecycle termination | Move SSE to offscreen document |
| Memory leak on server | Client map not cleaned up on disconnect | Listen for `req.on('close')` and remove client |
| Duplicate events after reconnect | Server replays from wrong offset | Validate sequence IDs with `Last-Event-ID` header |

---

## References

- [MDN: Using Server-Sent Events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events)
- [MDN: EventSource API](https://developer.mozilla.org/en-US/docs/Web/API/EventSource)
- [HTML Living Standard: Server-sent events](https://html.spec.whatwg.org/multipage/server-sent-events.html)
- [MDN: ReadableStream](https://developer.mozilla.org/en-US/docs/Web/API/ReadableStream)
- [Chrome: Offscreen Documents](https://developer.chrome.com/docs/extensions/reference/api/offscreen)
- [Chrome: Service Worker Lifecycle](https://developer.chrome.com/docs/extensions/develop/concepts/service-workers/lifecycle)
