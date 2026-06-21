<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly the standalone `mv3-service-worker-expert` skill.
> Sibling topics in this family are now reference files under the hubs (`chrome-extension-expert`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: mv3-service-worker-expert
description: >
  Chrome MV3 extension service worker lifecycle — idle termination (30s), operation
  timeout (5 min), fetch timeout (30s), alarm-based wakeup, event-driven activation,
  state persistence patterns, WebSocket keepalive, and migration from MV2 background pages.
  TRIGGER: debugging SW eviction or unexpected termination, designing persistent scheduling
  in MV3 extensions, handling SW restart gracefully, reviewing MV3 extension architecture,
  replacing MV2 background page patterns with service worker equivalents.
  SKIP: offscreen documents (use chrome-offscreen-documents), native messaging lifecycle
  (use chrome-native-messaging), OAuth token management (use chrome-identity-oauth),
  general Web API service workers outside Chrome extensions (entirely different lifecycle).
version: 1.1.0
category: developer
tags: [chrome-extension, mv3, service-worker, lifecycle, alarms, persistence, keepalive, websocket, migration]
related_skills: [chrome-offscreen-documents, chrome-storage-patterns, chrome-mv3-advanced]
updated: 2026-05-29
---

# MV3 Service Worker Lifecycle Expert

Expert reference for Chrome Manifest V3 extension service worker lifecycle management, termination handling, and persistent scheduling patterns.

## Core Lifecycle Rules

### Termination Conditions (Chrome 120+)

| Condition | Timeout | What happens |
|---|---|---|
| **Idle** | 30 seconds of no events or API calls | SW terminates; next event re-initializes |
| **Operation** | 5 minutes per single request | SW killed mid-operation; work lost |
| **Fetch** | 30 seconds for fetch() response | SW killed; fetch aborted |

**Timer resets:** Receiving any event or calling any chrome.* API resets the 30-second idle timer. WebSocket messages also reset it (Chrome 116+).

**Historical note:** Chrome 109 and earlier terminated SWs 5 minutes after startup regardless of activity. Removed in Chrome 110.

### What Wakes a Service Worker

A dormant SW is revived by any registered event:
- `chrome.runtime.onMessage` / `chrome.runtime.onConnect`
- `chrome.alarms.onAlarm`
- `chrome.tabs.onUpdated` / `chrome.tabs.onCreated`
- `chrome.webNavigation.*` events
- `chrome.action.onClicked`
- `chrome.runtime.onInstalled` / `chrome.runtime.onStartup`

The SW is re-initialized from scratch — all module-level code runs again, all global variables reset to initial values.

## State Persistence

Global variables are lost on every SW restart (every 30 seconds of idle).

| Storage | Survives | Quota | Use for |
|---|---|---|---|
| `chrome.storage.session` | SW restart, NOT browser restart | 10 MB | Auth tokens, disposable caches |
| `chrome.storage.local` | Everything except uninstall | 10 MB (unlimited with permission) | Durable user data, tracking state |
| `chrome.storage.sync` | Everything + sync across devices | 100 KB | Small user preferences |
| IndexedDB | Everything | Device storage | Large structured data |

**Pattern:** Read state from storage on every SW wakeup; never rely on global variables for anything that must persist.

```js
// WRONG — state lost on restart
let cachedData = null;
chrome.runtime.onMessage.addListener((msg) => {
  if (msg.type === 'get') return cachedData; // null after restart!
});

// RIGHT — read from storage on every access
chrome.runtime.onMessage.addListener(async (msg, sender, sendResponse) => {
  if (msg.type === 'get') {
    const stored = await chrome.storage.session.get('cachedData');
    sendResponse(stored.cachedData || null);
    return true; // keep channel open for async response
  }
});
```

## Alarm-Based Scheduling

### Minimum Intervals (Chrome 120+)

`chrome.alarms.create(name, { periodInMinutes: 0.5 })` — minimum 30 seconds (was 1 minute before Chrome 120).

### Alarm vs setTimeout/setInterval

| API | Works in MV3? | Why |
|---|---|---|
| `setTimeout` | Unreliable | Canceled when SW terminates |
| `setInterval` | Unreliable | Canceled when SW terminates |
| `chrome.alarms` | Yes | Chrome manages timer externally |

Alarms fire even when the SW is dormant — Chrome wakes the SW to deliver the alarm event.

```js
// Register alarm on install
chrome.runtime.onInstalled.addListener(() => {
  chrome.alarms.create('refresh', { periodInMinutes: 1 });
});

// Also register on startup (SW re-initialization)
chrome.runtime.onStartup.addListener(() => {
  chrome.alarms.create('refresh', { periodInMinutes: 1 });
});

// Handle alarm — runs even if SW was dormant
chrome.alarms.onAlarm.addListener((alarm) => {
  if (alarm.name === 'refresh') {
    doWork();
  }
});
```

## Event Registration

Event listeners **must** be registered at the top level of the SW script, synchronously, on every initialization. Listeners registered inside an async function or after an `await` may not be found when Chrome delivers events.

```js
// WRONG — listener may be missed after SW restart
async function setup() {
  const config = await loadConfig();
  chrome.runtime.onMessage.addListener(handler); // registered too late
}
setup();

// RIGHT — register synchronously at top level
chrome.runtime.onMessage.addListener(handler);
async function handler(msg) {
  const config = await loadConfig(); // config loaded inside the handler
}
```

**The "first 5 seconds" rule:** Chrome gives the SW about 5 seconds after wakeup to register all event listeners.

## WebSocket Keepalive (Chrome 116+)

Active WebSocket connections extend SW lifetime. Both sending and receiving messages reset the 30-second idle timer (Chrome 116+).

```js
const ws = new WebSocket('wss://example.com');

ws.onopen = () => {
  // Send a ping every 25s to reset the idle timer before the 30s cutoff.
  // This works because the open WebSocket port keeps the SW event loop alive;
  // the setInterval itself would be cleared on SW termination, but the WS
  // connection prevents termination as long as pings flow.
  setInterval(() => ws.send(JSON.stringify({ type: 'ping' })), 25000);
};

ws.onmessage = (event) => {
  processMessage(event.data); // each inbound message also resets the 30s idle timer
};
```

**Caveat:** WebSocket keepalive resets the idle timer but does not prevent the 5-minute operation timeout. Keep individual message handlers short; offload heavy computation to an offscreen document.

## Long-Running Operations

Any single event handler running longer than 5 minutes causes termination. Mitigation strategies:

**Chunk work:** Break large operations into smaller pieces, each triggered by a separate alarm.

```js
// Process N items per alarm tick, checkpoint progress to storage
chrome.alarms.onAlarm.addListener(async (alarm) => {
  if (alarm.name !== 'batch-process') return;
  const state = await chrome.storage.session.get('batchState');
  const { items, cursor } = state.batchState || { items: [], cursor: 0 };
  const BATCH_SIZE = 10;
  const batch = items.slice(cursor, cursor + BATCH_SIZE);
  for (const item of batch) await processItem(item);
  if (cursor + BATCH_SIZE >= items.length) {
    chrome.alarms.clear('batch-process');
  } else {
    await chrome.storage.session.set({
      batchState: { items, cursor: cursor + BATCH_SIZE }
    });
  }
});
```

**Offscreen documents:** For CPU-intensive work, delegate to an offscreen document (Chrome 109+) — see `chrome-offscreen-documents` skill.

## Migration from MV2

| MV2 Pattern | MV3 Replacement |
|---|---|
| `"persistent": true` background page | Service worker (always non-persistent) |
| Global variables for state | `chrome.storage.session` / `.local` |
| `setTimeout` / `setInterval` | `chrome.alarms` |
| `XMLHttpRequest` | `fetch()` (with 30s timeout awareness) |
| Background page DOM access | Offscreen documents (Chrome 109+) |
| `chrome.browserAction` | `chrome.action` |
| `chrome.extension.getBackgroundPage()` | `chrome.runtime.sendMessage()` |

## Chrome Version Timeline

| Version | Change |
|---|---|
| Chrome 105 | Native messaging keeps SW alive during connection |
| Chrome 109 | Offscreen documents API added |
| Chrome 110 | Removed 5-min-after-start termination |
| Chrome 114 | Message-based keepalive for long-lived ports |
| Chrome 116 | WebSocket activity resets idle timer |
| Chrome 118 | Debugger sessions keep SW alive |
| Chrome 120 | Alarm minimum reduced to 30 seconds |

## Anti-Patterns

1. **Relying on global state** — variables reset on every restart (every 30s of idle)
2. **Registering listeners in async init** — Chrome may miss them on wakeup
3. **Using setTimeout for scheduling** — timers die with the SW
4. **Fetching without timeout handling** — 30s fetch timeout kills the SW
5. **Running operations > 5 minutes** — hard kill with no warning
6. **Assuming SW is always active** — it is event-driven, not persistent
7. **Not re-creating alarms on startup** — alarms persist but listener registration does not

## Troubleshooting

| Symptom | Likely cause | Fix |
|---|---|---|
| Event listener not firing after extension reload | Listener registered inside an async function | Move registration to top level, synchronously |
| State lost randomly | Global variables used instead of chrome.storage | Persist to session or local storage |
| Long operation gets killed | Operation exceeds 5 minutes | Use chunked processing with alarms |
| Extension slow to respond | SW waking from dormancy | Expected; ~100ms cold start on first event after idle |

## References

- [Extension Service Worker Lifecycle](https://developer.chrome.com/docs/extensions/develop/concepts/service-workers/lifecycle)
- [Migrate to a Service Worker](https://developer.chrome.com/docs/extensions/develop/migrate/to-service-workers)
- [What is Manifest V3?](https://developer.chrome.com/docs/extensions/develop/migrate/what-is-mv3)
