<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly the standalone `chrome-notifications-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`chrome-extension-expert`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: chrome-notifications-patterns
description: >
  Chrome extension notification and alert patterns — chrome.notifications API
  (create/update/clear), notification templates (basic/image/list/progress),
  click and button handlers, notification permissions, in-page toast/snackbar
  via shadow DOM, alert popup windows (chrome.windows.create), badge and
  notification coordination, notification sound via offscreen documents,
  alarm-driven recurring notifications, notification ID management and
  deduplication, priority and requireInteraction options.
  TRIGGER: implementing notifications in Chrome extensions, coordinating badge
  and notification state, building in-page toasts from content scripts,
  debugging notification delivery issues, scheduling recurring notifications
  with alarms, deduplicating notifications by ID.
  SKIP: web push notifications outside extensions (no chrome.notifications API),
  browser native notification permission (Notification API) without chrome.notifications,
  server-side push or FCM without an extension context.
version: 1.1.0
category: developer
tags: [chrome-extension, notifications, mv3, service-worker, badge, toast, snackbar, offscreen, alarms, shadow-dom, popup-window, requireInteraction, priority, deduplication]
related_skills: [mv3-service-worker-expert, chrome-offscreen-documents, chrome-storage-patterns, chrome-tabs-management]
updated: 2026-05-29
---

# Chrome Notifications Patterns

## Overview

Chrome MV3 extensions have four notification surfaces. Choosing the wrong surface causes notification fatigue or missed alerts.

| Surface | Urgency | Persists? | User input? | Works when tab closed? |
|---------|---------|-----------|-------------|----------------------|
| System notification (`chrome.notifications`) | Medium–High | Until dismissed or timeout | Click / 2 buttons max | Yes |
| Badge (`chrome.action`) | Low (passive) | Until cleared | None | Yes |
| In-page toast (content script) | Low–Medium | Auto-dismiss 3–8s | Optional dismiss | No |
| Popup window (`chrome.windows.create`) | High (critical) | Until window closed | Full HTML form | Yes |

## Manifest Requirements

```json
{
  "manifest_version": 3,
  "permissions": ["notifications", "alarms", "offscreen"],
  "background": { "service_worker": "service-worker.js" },
  "action": { "default_icon": "icons/icon-16.png" }
}
```

- `"notifications"` — required for `chrome.notifications` API and all notification events
- `"alarms"` — required only for recurring/scheduled notifications
- `"offscreen"` — required only if playing notification sounds via offscreen audio
- `"action"` key — required for badge text/color via `chrome.action.setBadgeText`

## chrome.notifications API

### Template Types

All templates require `type`, `iconUrl`, `title`, and `message`.

```javascript
// Basic — simple informational notification
await chrome.notifications.create('case-update-123', {
  type: 'basic',
  iconUrl: chrome.runtime.getURL('icons/icon-128.png'),
  title: 'Case Updated',
  message: 'Case 01234567 moved to Waiting for Customer.',
  priority: 1,
  requireInteraction: false
});

// List — multi-line summary
await chrome.notifications.create('batch-update', {
  type: 'list',
  iconUrl: chrome.runtime.getURL('icons/icon-128.png'),
  title: '3 Case Updates',
  message: '3 items',
  items: [
    { title: '01234567', message: 'Waiting for Customer' },
    { title: '01234568', message: 'Escalated' },
    { title: '01234569', message: 'Resolved' }
  ]
});

// Progress — shows a 0–100 progress bar
await chrome.notifications.create('sync-progress', {
  type: 'progress',
  iconUrl: chrome.runtime.getURL('icons/icon-128.png'),
  title: 'Syncing Cases',
  message: 'Fetching case data...',
  progress: 42
});
// Update progress: chrome.notifications.update('sync-progress', { progress: 75 });
```

### Notification Options Reference

| Option | Type | Default | Notes |
|--------|------|---------|-------|
| `type` | string | required | `'basic'`, `'image'`, `'list'`, `'progress'` |
| `iconUrl` | string | required | Use `chrome.runtime.getURL()` for bundled assets |
| `title` | string | required | Primary heading |
| `message` | string | required | Body text |
| `priority` | number | 0 | Range −2 to 2; values below 0 may not display on some platforms |
| `requireInteraction` | boolean | false | Stays until user dismisses (use for critical alerts) |
| `silent` | boolean | false | Suppresses OS notification sound |
| `buttons` | array | [] | Up to 2 buttons: `[{ title: 'View' }, { title: 'Dismiss' }]` |
| `contextMessage` | string | — | Secondary text below the message |
| `eventTime` | number | — | Timestamp in ms; displayed in the notification |
| `items` | array | — | Required for `list` type |
| `progress` | number | — | Required for `progress` type (0–100) |
| `imageUrl` | string | — | Required for `image` type |

### Notification ID Management and Deduplication

Calling `create()` with an ID that already exists **replaces** the existing notification — use this as the primary deduplication strategy.

```javascript
// Deterministic ID prevents duplicate notifications for the same event
function buildNotificationId(caseNumber, eventType) {
  return `mca-${caseNumber}-${eventType}`;
}

// Clear all notifications matching a prefix
async function clearNotificationsByPrefix(prefix) {
  const all = await chrome.notifications.getAll();
  await Promise.all(
    Object.keys(all)
      .filter(id => id.startsWith(prefix))
      .map(id => chrome.notifications.clear(id))
  );
}
```

## Event Handling

All notification events fire in the service worker. Register listeners at the top level so they survive SW restarts.

```javascript
// service-worker.js

// Notification body clicked (not a button)
chrome.notifications.onClicked.addListener((notificationId) => {
  chrome.tabs.create({ url: buildCaseUrl(notificationId) });
  chrome.notifications.clear(notificationId);
});

// Button clicked — route by notification ID prefix + button index
chrome.notifications.onButtonClicked.addListener(async (notificationId, buttonIndex) => {
  if (notificationId.startsWith('escalation-')) {
    const caseNumber = notificationId.replace('escalation-', '');
    if (buttonIndex === 0) {
      await chrome.tabs.create({ url: `https://support.mongodb.com/case/${caseNumber}` });
    } else {
      await markEscalationAcknowledged(caseNumber);
    }
  }
  await chrome.notifications.clear(notificationId);
});

// Notification dismissed by user or system
chrome.notifications.onClosed.addListener((notificationId, byUser) => {
  if (byUser) recordDismissal(notificationId);
});
```

**Always call `chrome.notifications.clear(id)` after handling a click or button event** — the notification does not auto-clear on interaction.

## Badge + Notification Coordination

Badge text provides a persistent count indicator that complements transient system notifications.

```javascript
// service-worker.js
const BADGE_COLORS = {
  critical: '#DC2626',
  warning:  '#F59E0B',
  info:     '#3B82F6',
  clear:    '#22C55E'
};

async function updateBadge(count, severity = 'info') {
  await Promise.all([
    chrome.action.setBadgeText({ text: count > 0 ? String(count) : '' }),
    chrome.action.setBadgeBackgroundColor({ color: BADGE_COLORS[severity] || BADGE_COLORS.info })
  ]);
}

// Decrement badge when notification is dismissed
chrome.notifications.onClosed.addListener(async (notificationId, byUser) => {
  if (byUser) {
    const { pendingCount = 0 } = await chrome.storage.session.get('pendingCount');
    const newCount = Math.max(0, pendingCount - 1);
    await chrome.storage.session.set({ pendingCount: newCount });
    await updateBadge(newCount);
  }
});

// Per-tab badge
async function updateTabBadge(tabId, count, severity = 'info') {
  await Promise.all([
    chrome.action.setBadgeText({ text: count > 0 ? String(count) : '', tabId }),
    chrome.action.setBadgeBackgroundColor({ color: BADGE_COLORS[severity], tabId })
  ]);
}
```

## Alarm-Driven Recurring Notifications

Service workers are ephemeral. Use `chrome.alarms` to schedule notification checks that survive SW restarts.

```javascript
// service-worker.js

// Register on install; also re-register on startup in case alarms were lost
chrome.runtime.onInstalled.addListener(() => {
  chrome.alarms.create('check-case-updates', {
    delayInMinutes: 1,
    periodInMinutes: 5
  });
});

chrome.runtime.onStartup.addListener(async () => {
  const existing = await chrome.alarms.get('check-case-updates');
  if (!existing) {
    chrome.alarms.create('check-case-updates', { delayInMinutes: 1, periodInMinutes: 5 });
  }
});

chrome.alarms.onAlarm.addListener(async (alarm) => {
  if (alarm.name !== 'check-case-updates') return;
  try {
    const updates = await fetchCaseUpdates();
    for (const update of updates) {
      await chrome.notifications.create(
        buildNotificationId(update.caseNumber, update.eventType),
        {
          type: 'basic',
          iconUrl: chrome.runtime.getURL('icons/icon-128.png'),
          title: `Case ${update.caseNumber}`,
          message: `Event: ${update.eventType}`,
          priority: 1
        }
      );
    }
  } catch (err) {
    console.error('[alarm:check-case-updates] failed:', err.message);
  }
});
```

**Minimum interval:** `chrome.alarms` enforces a minimum of 1 minute in production and 30 seconds for unpacked extensions.

## In-Page Toast / Snackbar (Content Script)

For notifications that belong inside the page (not system tray), inject a toast via content script using shadow DOM isolation.

```javascript
// content-script.js
let _toastShadow = null;

function createToastHost() {
  if (_toastShadow) return _toastShadow;
  const host = document.createElement('div');
  host.id = 'mca-toast-host';
  document.body.appendChild(host);
  const shadow = host.attachShadow({ mode: 'closed' });
  _toastShadow = shadow;

  const style = document.createElement('style');
  style.textContent = `
    .mca-toast-container {
      position: fixed; bottom: 20px; right: 20px; z-index: 2147483647;
      display: flex; flex-direction: column-reverse; gap: 8px;
      pointer-events: none;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    }
    .mca-toast {
      pointer-events: auto; max-width: 360px; padding: 12px 16px;
      border-radius: 8px; color: #fff; font-size: 13px; line-height: 1.4;
      box-shadow: 0 4px 12px rgba(0,0,0,0.25);
      opacity: 0; transform: translateY(12px);
      animation: mca-slide-in 0.25s ease forwards;
    }
    .mca-toast.info    { background: #1e40af; }
    .mca-toast.success { background: #15803d; }
    .mca-toast.warning { background: #b45309; }
    .mca-toast.error   { background: #b91c1c; }
    .mca-toast.exiting { animation: mca-slide-out 0.2s ease forwards; }
    .mca-toast-title   { font-weight: 600; margin-bottom: 2px; }
    .mca-toast-dismiss { float: right; cursor: pointer; opacity: 0.7; margin-left: 12px; }
    .mca-toast-dismiss:hover { opacity: 1; }
    @keyframes mca-slide-in  { to { opacity: 1; transform: translateY(0); } }
    @keyframes mca-slide-out { to { opacity: 0; transform: translateY(12px); } }
  `;
  const container = document.createElement('div');
  container.className = 'mca-toast-container';
  shadow.appendChild(style);
  shadow.appendChild(container);
  return shadow;
}

function showToast({ title, message, level = 'info', duration = 5000 }) {
  const shadow = createToastHost();
  const container = shadow.querySelector('.mca-toast-container');
  const toast = document.createElement('div');
  toast.className = `mca-toast ${level}`;

  // Build DOM via textContent — no innerHTML to avoid XSS
  const dismissBtn = document.createElement('span');
  dismissBtn.className = 'mca-toast-dismiss';
  dismissBtn.textContent = '×';
  dismissBtn.addEventListener('click', () => dismissToast(toast));
  toast.appendChild(dismissBtn);

  if (title) {
    const titleEl = document.createElement('div');
    titleEl.className = 'mca-toast-title';
    titleEl.textContent = title;
    toast.appendChild(titleEl);
  }
  const messageEl = document.createElement('div');
  messageEl.textContent = message;
  toast.appendChild(messageEl);
  container.appendChild(toast);

  if (duration > 0) setTimeout(() => dismissToast(toast), duration);
  return toast;
}

function dismissToast(toast) {
  if (!toast || toast.classList.contains('exiting')) return;
  toast.classList.add('exiting');
  toast.addEventListener('animationend', () => toast.remove(), { once: true });
}

// Listen for messages from service worker
chrome.runtime.onMessage.addListener((msg) => {
  if (msg.type === 'MCA_SHOW_TOAST') showToast(msg.payload);
});
```

**Triggering from service worker:**
```javascript
async function sendToastToActiveTab(payload) {
  const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
  if (tab?.id) {
    chrome.tabs.sendMessage(tab.id, { type: 'MCA_SHOW_TOAST', payload })
      .catch(() => {
        // Content script not injected — fall back to system notification
        showBasicNotification(`toast-fallback-${Date.now()}`, payload.title, payload.message);
      });
  }
}
```

## Notification Sound via Offscreen Document

MV3 service workers cannot play audio directly. Use an offscreen document:

```javascript
// service-worker.js
async function playNotificationSound(soundFile = 'sounds/alert.mp3', volume = 0.5) {
  const existing = await chrome.runtime.getContexts({
    contextTypes: ['OFFSCREEN_DOCUMENT']
  });
  if (existing.length === 0) {
    await chrome.offscreen.createDocument({
      url: chrome.runtime.getURL('offscreen/audio.html'),
      reasons: ['AUDIO_PLAYBACK'],
      justification: 'Play notification sound effects'
    });
  }
  await chrome.runtime.sendMessage({
    type: 'MCA_PLAY_SOUND',
    target: 'offscreen',
    soundUrl: chrome.runtime.getURL(soundFile),
    volume
  });
}

// offscreen/audio-player.js
let audio = null;
chrome.runtime.onMessage.addListener((msg) => {
  if (msg.type !== 'MCA_PLAY_SOUND' || msg.target !== 'offscreen') return;
  if (audio) { audio.pause(); audio.currentTime = 0; }
  audio = new Audio(msg.soundUrl);
  audio.volume = msg.volume ?? 0.5;
  audio.play().catch(err => console.warn('[offscreen] audio play failed:', err.message));
});
```

**Note:** Offscreen documents created with `AUDIO_PLAYBACK` auto-close after 30 seconds of silence. Always check for an existing document before sending a play message.

## Alert Popup Windows

For critical alerts requiring dedicated UI (acknowledgement forms, multi-field input):

```javascript
// service-worker.js
async function openAlertPopup(alertData) {
  await chrome.storage.session.set({ pendingAlert: alertData });
  const width = 420, height = 340;
  await chrome.windows.create({
    url: chrome.runtime.getURL('alerts/alert.html'),
    type: 'popup', width, height, focused: true
  });
}

// alerts/alert.js
document.addEventListener('DOMContentLoaded', async () => {
  const { pendingAlert } = await chrome.storage.session.get('pendingAlert');
  if (!pendingAlert) { window.close(); return; }

  document.getElementById('alert-title').textContent = pendingAlert.title;
  document.getElementById('alert-message').textContent = pendingAlert.message;

  document.getElementById('btn-acknowledge').addEventListener('click', async () => {
    await chrome.runtime.sendMessage({ type: 'MCA_ALERT_ACKNOWLEDGED', alertId: pendingAlert.id });
    await chrome.storage.session.remove('pendingAlert');
    window.close();
  });
  document.getElementById('btn-dismiss').addEventListener('click', async () => {
    await chrome.storage.session.remove('pendingAlert');
    window.close();
  });
});
```

**Limitations:** Not modal, cannot stay on top of other windows, slower to open than system notifications. Use only when the alert requires user input beyond a simple click/dismiss.

## Anti-Patterns

### Random notification IDs (causes spam)
```javascript
// BAD: new UUID every time = no deduplication
chrome.notifications.create(crypto.randomUUID(), options);

// GOOD: deterministic ID from event source
chrome.notifications.create(`case-${caseNumber}-${eventType}`, options);
```

### Not clearing after button clicks
```javascript
// BAD: notification lingers after user acted
chrome.notifications.onButtonClicked.addListener((id, idx) => { doSomething(id); });

// GOOD: always clear after handling
chrome.notifications.onButtonClicked.addListener((id, idx) => {
  doSomething(id);
  chrome.notifications.clear(id);
});
```

### Forgetting to re-register alarms on startup
```javascript
// BAD: only register on install — alarm lost after browser restart
chrome.runtime.onInstalled.addListener(() => {
  chrome.alarms.create('check', { periodInMinutes: 5 });
});

// GOOD: also register on startup
chrome.runtime.onStartup.addListener(async () => {
  if (!await chrome.alarms.get('check')) {
    chrome.alarms.create('check', { periodInMinutes: 5 });
  }
});
```

### Rate-limit state in JS variables (lost on SW restart)
```javascript
// BAD: in-memory Map lost on SW termination
const lastNotified = new Map();

// GOOD: use chrome.storage.session for rate-limit state
async function shouldNotify(caseNumber) {
  const key = `notif-ts-${caseNumber}`;
  const stored = await chrome.storage.session.get(key);
  const lastTime = stored[key] || 0;
  if (Date.now() - lastTime < 5 * 60 * 1000) return false;
  await chrome.storage.session.set({ [key]: Date.now() });
  return true;
}
```

## Platform Differences

| Behavior | Windows | macOS | Linux | ChromeOS |
|----------|---------|-------|-------|----------|
| Priority −2 to −1 | Not shown | Not shown | Not shown | Shown in notification center |
| `requireInteraction` | Stays on screen | Stays in notification center | Varies by DE | Stays on screen |
| Buttons | Up to 2 inline | Up to 2 inline | May not render | Up to 2 inline |
| Image template | Full image | Thumbnail only | Varies | Full image |

## Permission Check

The `"notifications"` permission in `manifest.json` is declared permission — it does not guarantee the user has granted it (users can revoke it in OS settings). Check before creating:

```javascript
async function safeNotify(id, options) {
  const granted = await chrome.permissions.contains({ permissions: ['notifications'] });
  if (!granted) {
    // Fall back to badge-only or in-page toast
    await updateBadge(1, 'warning');
    return false;
  }
  await chrome.notifications.create(id, options);
  return true;
}
```

## Quick Recipes

### Notify only if user opted in
```javascript
async function notifyIfEnabled(id, options) {
  const { notificationsEnabled = true } = await chrome.storage.local.get('notificationsEnabled');
  if (!notificationsEnabled) return;
  await chrome.notifications.create(id, options);
}
```

### Batch multiple updates into one notification
```javascript
let pendingUpdates = [];
let batchTimer = null;

function queueUpdate(update) {
  pendingUpdates.push(update);
  if (!batchTimer) batchTimer = setTimeout(flushUpdates, 3000);
}

async function flushUpdates() {
  batchTimer = null;
  const batch = pendingUpdates.splice(0);
  if (batch.length === 0) return;
  if (batch.length === 1) {
    await showBasicNotification(buildNotificationId(batch[0].caseNumber, 'update'),
      `Case ${batch[0].caseNumber}`, batch[0].summary);
  } else {
    await chrome.notifications.create('batch-update', {
      type: 'list',
      iconUrl: chrome.runtime.getURL('icons/icon-128.png'),
      title: `${batch.length} Case Updates`,
      message: `${batch.length} items`,
      items: batch.map(u => ({ title: u.caseNumber, message: u.summary }))
    });
  }
}
```
