<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly the standalone `chrome-badge-metrics` skill.
> Sibling topics in this family are now reference files under the hubs (`chrome-extension-expert`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: chrome-badge-metrics
version: 1.1.0
category: developer
tags: [chrome-extension, badge, action-icon, mv3, service-worker, setBadgeText, setBadgeBackgroundColor, setBadgeTextColor, setIcon, OffscreenCanvas, state-machine, accessibility, metrics]
description: >-
  Chrome extension badge and action icon state management — chrome.action API (setBadgeText,
  setBadgeBackgroundColor, setBadgeTextColor, setIcon, setTitle, setPopup, enable/disable),
  per-tab vs global badge state, dynamic icon generation with OffscreenCanvas, badge state machines,
  metric-driven badge text, animation patterns, icon accessibility, and cross-context badge refresh.
  TRIGGER: implementing or reviewing extension toolbar badge/icon behavior, building metric dashboards
  on the badge, debugging badge not updating after SW restart, designing icon state machines, or
  per-tab vs global badge scoping questions.
  SKIP: extension popup UI (use chrome-dev); notification toasts (use chrome-notifications-patterns);
  side panel UI (use chrome-dev).
triggers:
  - "Show a count badge on my extension toolbar icon"
  - "Badge not restoring after service worker restarts"
  - "Per-tab vs global badge state in Chrome extension"
  - "Generate a dynamic icon with OffscreenCanvas in service worker"
  - "Badge state machine for multiple extension states"
  - "setBadgeText with a countdown timer"
  - "Extension badge accessibility — screen reader support"
  - "Animate the extension badge icon"
  - "chrome.action.setBadgeTextColor not working"
  - "Clear per-tab badge override on navigation"
related_skills:
  - mv3-service-worker-expert
  - chrome-dev
  - alarm-scheduler-patterns
---

# Chrome Badge Metrics

The `chrome.action` API controls the extension toolbar button: icon, badge text, badge colors, tooltip title, and popup. In MV3 it replaces `chrome.browserAction` and `chrome.pageAction`. Every setter is available from the service worker, popups, options pages, and extension-owned tabs.

The badge — up to 4 visible characters overlaid on the toolbar icon — is a persistent, always-visible status channel that requires no user interaction.

**When not to use:** popup UI design, side panel UI, or desktop notification toasts.

## Manifest prerequisite

The `chrome.action` API requires an `"action"` key in `manifest.json`. No separate permission is needed.

```json
{
  "manifest_version": 3,
  "action": {
    "default_icon": { "16": "icons/icon-16.png", "32": "icons/icon-32.png" },
    "default_title": "My Extension",
    "default_popup": "popup.html"
  }
}
```

To control everything via the API only: `"action": {}`.

## Core API

### setBadgeText

```javascript
await chrome.action.setBadgeText({ text: '3' });          // global
await chrome.action.setBadgeText({ text: '!', tabId });   // per-tab
await chrome.action.setBadgeText({ text: '' });           // clear global
await chrome.action.setBadgeText({ text: null, tabId }); // clear per-tab override → falls back to global
```

- Empty string `''` hides the badge.
- `null` + `tabId` clears the per-tab override (falls back to global).
- Text must be a **string** — passing a number throws or silently fails.
- Maximum 4 visible characters; longer strings are truncated.

### setBadgeBackgroundColor / setBadgeTextColor

```javascript
await chrome.action.setBadgeBackgroundColor({ color: '#0f62fe' });
await chrome.action.setBadgeBackgroundColor({ color: [0, 100, 200, 255], tabId });
await chrome.action.setBadgeTextColor({ color: '#ffffff' });  // Chrome 110+
```

- If `setBadgeTextColor` is not called, Chrome auto-selects white or black based on contrast — usually the right choice.
- Alpha value of 0 is rejected.

### setTitle

```javascript
await chrome.action.setTitle({ title: 'My Extension — 3 tracked cases' });
await chrome.action.setTitle({ title: null, tabId });  // clear per-tab → falls back to global
```

- `null` + `tabId` clears the per-tab override.
- `''` explicitly sets an empty title (removes tooltip — usually a mistake). Use `null` to clear.
- The title is the **accessible name** of the toolbar button. Screen readers announce it on focus. Always include meaningful state, not just the extension name.

### setIcon

```javascript
await chrome.action.setIcon({ path: { 16: 'icons/icon-16.png', 32: 'icons/icon-32.png' } });
await chrome.action.setIcon({ imageData: canvasImageData });  // from OffscreenCanvas
await chrome.action.setIcon({ path: 'icons/alert.png', tabId });
```

- SVG is **not** supported.
- Provide 16px and 32px for 1x/2x displays. Add 24px and 48px for unusual scale factors.
- Unpacked extensions must use PNG.

### setPopup / onClicked

```javascript
await chrome.action.setPopup({ popup: '' });         // disables popup → enables onClicked
await chrome.action.setPopup({ popup: 'popup.html' });

chrome.action.onClicked.addListener((tab) => {
  // Only fires when popup is '' (disabled)
  toggleSidePanel(tab.id);
});
```

### enable / disable

```javascript
await chrome.action.disable(tabId);   // grey out for specific tab
await chrome.action.enable(tabId);
const on = await chrome.action.isEnabled({ tabId });  // Chrome 110+
```

## Per-tab vs global state

Every setter accepts an optional `tabId`. Priority: tab-specific value → global value → manifest default.

- Set `tabId` when the badge reflects page-specific data (e.g., case count on a support page).
- Omit `tabId` for extension-wide state (e.g., total tracked cases).
- Tab-specific values are automatically cleared when the tab closes.

```javascript
// Clear per-tab badge on URL change
chrome.tabs.onUpdated.addListener((tabId, changeInfo) => {
  if (changeInfo.url) {
    chrome.action.setBadgeText({ text: null, tabId });
    chrome.action.setTitle({ title: null, tabId });
  }
});
```

## Metric-driven badge pattern

Separate pure computation from chrome.action calls to keep the computation testable:

```javascript
// Pure function — no chrome.* calls
export function buildBadgeDisplay({ metric, trackedCases = [], now = Date.now() }) {
  const cases = Array.isArray(trackedCases) ? trackedCases : [];
  switch (metric) {
    case 'tracked_cases':
      return {
        text: cases.length > 0 ? String(Math.min(999, cases.length)) : '',
        title: `${cases.length} tracked case${cases.length === 1 ? '' : 's'}.`,
      };
    case 'open_not_pending': {
      const count = cases.filter(c => !isResolved(c) && !isPending(c)).length;
      return { text: count > 0 ? String(count) : '', title: `${count} open cases not in customer pending.` };
    }
    default:
      return { text: '', title: 'No metric selected.' };
  }
}

// Renderer — calls chrome.action
async function refreshActionBadge({ options = null, trackingState = null } = {}) {
  if (!chrome.action?.setBadgeText) return null;
  const [opts, state] = await Promise.all([options ?? getOptions(), trackingState ?? getTrackingState()]);
  const display = buildBadgeDisplay({
    metric: opts.badgeMetric,
    trackedCases: Object.values(state.trackedCases || {}),
    now: Date.now(),
  });
  await Promise.all([
    chrome.action.setBadgeText({ text: display.text || '' }),
    chrome.action.setTitle({ title: `My Extension — ${display.title}` }),
    chrome.action.setBadgeBackgroundColor({ color: display.text ? '#0f62fe' : '#6b7280' }),
  ]).catch(() => {});
  return display;
}
```

### Badge text format helpers

```javascript
function formatCount(n)         { return Number.isFinite(n) && n > 0 ? String(Math.min(999, Math.round(n))) : ''; }
function formatSignedInt(n)     { if (!Number.isFinite(n)) return ''; const v = Math.max(-99, Math.min(99, Math.round(n))); return v > 0 ? `+${v}` : String(v); }
function formatDuration(ms) {
  if (!Number.isFinite(ms) || ms < 0) return '';
  const m = Math.round(ms / 60000);
  if (m < 100) return `${Math.max(0, m)}m`;
  const h = Math.round(ms / 3600000);
  if (h < 24) return `${h}h`;
  return `${Math.min(Math.round(ms / 86400000), 999)}d`;
}
```

## Dynamic icons with OffscreenCanvas

MV3 service workers have no DOM. Use `OffscreenCanvas` to generate icons:

```javascript
function createIconAtSize(size, count, color = '#0f62fe') {
  const canvas = new OffscreenCanvas(size, size);
  const ctx = canvas.getContext('2d');
  ctx.beginPath();
  ctx.arc(size / 2, size / 2, size / 2, 0, Math.PI * 2);
  ctx.fillStyle = color;
  ctx.fill();
  const fontSize = Math.round(size * 0.55);
  ctx.fillStyle = '#ffffff';
  ctx.font = `bold ${fontSize}px sans-serif`;
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';
  ctx.fillText(String(count), size / 2, size / 2 + 1);
  return ctx.getImageData(0, 0, size, size);
}

async function setDynamicIcon(count, color) {
  await chrome.action.setIcon({
    imageData: { 16: createIconAtSize(16, count, color), 32: createIconAtSize(32, count, color) },
  });
}
```

### Composite icon (base + status dot)

```javascript
async function createOverlayIcon(baseIconPath, dotColor) {
  const blob = await (await fetch(chrome.runtime.getURL(baseIconPath))).blob();
  const bitmap = await createImageBitmap(blob);
  const size = 32;
  const canvas = new OffscreenCanvas(size, size);
  const ctx = canvas.getContext('2d');
  ctx.drawImage(bitmap, 0, 0, size, size);
  const r = 5, x = size - r - 1, y = size - r - 1;
  ctx.beginPath(); ctx.arc(x, y, r, 0, Math.PI * 2);
  ctx.fillStyle = dotColor; ctx.fill();
  ctx.strokeStyle = '#ffffff'; ctx.lineWidth = 1.5; ctx.stroke();
  return ctx.getImageData(0, 0, size, size);
}
```

## Badge state machine

```javascript
const BADGE_STATES = Object.freeze({
  idle:        { text: '',     bg: '#6b7280', title: 'Idle' },
  loading:     { text: '...',  bg: '#f59e0b', title: 'Loading data...' },
  active:      { text: null,   bg: '#0f62fe', title: null },   // dynamic text
  error:       { text: 'ERR',  bg: '#e53e3e', title: 'Error — click for details' },
  authExpired: { text: 'AUTH', bg: '#dc2626', title: 'Authentication expired' },
  disabled:    { text: 'OFF',  bg: '#9ca3af', title: 'Extension disabled' },
});

let currentState = 'idle';

async function transitionBadge(newState, dynamicText = null) {
  const config = BADGE_STATES[newState];
  if (!config) return;
  currentState = newState;
  const text = dynamicText ?? config.text ?? '';
  await Promise.all([
    chrome.action.setBadgeText({ text }),
    chrome.action.setBadgeBackgroundColor({ color: config.bg }),
    chrome.action.setTitle({ title: `My Extension — ${config.title || text || 'Ready'}` }),
  ]).catch(() => {});
}
```

## Refresh triggers

```javascript
// 1. Data change
async function onTrackingStateChanged(newState) { await refreshActionBadge({ trackingState: newState }); }

// 2. Options change
chrome.storage.onChanged.addListener((changes, area) => {
  if (area === 'local' && changes.options) refreshActionBadge({ options: changes.options.newValue });
});

// 3. Periodic countdown refresh
chrome.alarms.create('badge-refresh', { periodInMinutes: 1 });
chrome.alarms.onAlarm.addListener(alarm => { if (alarm.name === 'badge-refresh') refreshActionBadge(); });

// 4. SW restart — CRITICAL: call at top level of service worker module
refreshActionBadge();
chrome.runtime.onInstalled.addListener(() => refreshActionBadge());
chrome.runtime.onStartup.addListener(() => refreshActionBadge());
```

## Animation

Animate briefly (2–5 seconds max). Design animations so the final state is safe if the SW suspends mid-animation — set the desired end state *before* starting, then animate, then confirm end state again.

```javascript
async function pulseBadge(normalColor, pulseColor, cycles = 3) {
  for (let i = 0; i < cycles; i++) {
    await chrome.action.setBadgeBackgroundColor({ color: pulseColor });
    await new Promise(r => setTimeout(r, 300));
    await chrome.action.setBadgeBackgroundColor({ color: normalColor });
    await new Promise(r => setTimeout(r, 300));
  }
}
```

## Accessibility

- **Title = accessible name.** Always include extension name + current state: `"My Extension — 3 tracked cases"`.
- **Update title with every badge change.** Badge text is visual-only; screen readers do not announce it.
- **Color-blind safety:** pair every color change with a text change — never use color as the sole state indicator.
- `setBadgeTextColor` auto-contrasts by default. Override only when the design requires a specific foreground color; if you do, ensure ≥ 4.5:1 contrast ratio (WCAG AA).

## Anti-patterns

| Anti-pattern | Fix |
|-------------|-----|
| No `refreshActionBadge()` at SW top level | Badge goes blank after SW restart — call at module top level |
| `setBadgeText({ text: 42 })` | Always pass a string: `String(42)` |
| Badge text > 4 chars | Use `formatCount()` — cap at `'999+'` |
| Continuous `setIcon` animation | Animate briefly (≤ 5s); restore static icon |
| Badge color change without text change | Add text so color-blind users get the signal |
| Per-tab state without cleanup | Add `tabs.onUpdated` navigation cleanup |

## Testing

`buildBadgeDisplay` is a pure function — test it without Chrome API mocks:

```javascript
import { describe, it, expect } from 'vitest';
import { buildBadgeDisplay } from '../../src/background/badge-metrics.js';

describe('buildBadgeDisplay', () => {
  it('counts tracked cases', () => {
    const r = buildBadgeDisplay({ metric: 'tracked_cases', trackedCases: [{}, {}, {}] });
    expect(r.text).toBe('3');
  });
  it('returns empty string for zero', () => {
    const r = buildBadgeDisplay({ metric: 'tracked_cases', trackedCases: [] });
    expect(r.text).toBe('');
  });
});
```

For integration tests, mock the API:

```javascript
function mockChromeAction() {
  const calls = { setBadgeText: [], setBadgeBackgroundColor: [], setTitle: [] };
  globalThis.chrome = {
    action: {
      setBadgeText: async (o) => calls.setBadgeText.push(o),
      setBadgeBackgroundColor: async (o) => calls.setBadgeBackgroundColor.push(o),
      setTitle: async (o) => calls.setTitle.push(o),
    },
  };
  return calls;
}
```
