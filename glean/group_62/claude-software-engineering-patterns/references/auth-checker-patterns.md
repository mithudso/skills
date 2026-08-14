<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `auth-checker-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: auth-checker-patterns
title: "Auth Checker Patterns"
description: >-
  Multi-service auth state monitoring for Chrome MV3 extensions. Covers endpoint
  registry design, probing strategies (cookie/token/tab-session), expiry countdown,
  floating notification widgets, alarm-based polling, storage-driven cross-context sync,
  state machines, remediation UIs, and bootstrap sequencing.
  TRIGGER: building an auth dashboard, token health check, or multi-provider auth
  monitor in a Chrome extension; proactive cookie or token expiry detection; service
  worker alarm-based auth polling; auth bootstrap sequences on install/startup.
  SKIP: single-provider OAuth (use chrome-identity-oauth or web-auth-patterns);
  server-side session management without a browser extension; token refresh logic
  entirely in a backend.
version: "1.1.0"
updated: "2026-05-29"
category: developer
tags:
  - chrome-extension
  - auth
  - mv3
  - service-worker
  - cookie
  - token
  - polling
  - storage
keywords:
  - auth checker
  - auth dashboard
  - cookie expiry
  - token expiry
  - chrome alarms
  - chrome storage
  - auth state machine
  - cross-context sync
  - auth bootstrap
  - floating widget
  - HMAC
  - bearer token
  - service worker polling
whenToUse:
  - "Multi-service auth status dashboard or floating indicator in a Chrome extension"
  - "Proactive token or cookie expiry detection with countdown warnings"
  - "Service worker alarm-based auth polling and cross-context sync"
  - "Auth bootstrap sequences on install or startup"
  - "Auth health checks across cookie, token, and tab-session providers"
whenNotToUse:
  - "Single-provider OAuth — use chrome-identity-oauth or web-auth-patterns"
  - "Server-side session management without a browser extension"
  - "Token refresh logic entirely in a backend"
related_skills:
  - chrome-storage-patterns
  - alarm-scheduler-patterns
  - mv3-service-worker-expert
  - chrome-identity-oauth
  - web-auth-patterns
---

# Auth Checker Patterns — Multi-Service Auth Monitoring for Chrome Extensions

## 1. Endpoint registry pattern

```js
// Replace message types and URLs with your extension's values.
const AUTH_ENDPOINTS = [
  {
    id: 'api',
    label: 'Primary API',
    checkType: 'APP_CHECK_AUTH',          // message to service worker
    extract: (r) => r?.result?.api || {},
    okField: 'ok',
    openUrl: 'https://app.example.com',   // remediation: open login page
    severity: 'critical'                  // critical = blocks core features
  },
  {
    id: 'llm',
    label: 'LLM Provider',
    checkType: 'APP_CHECK_AUTH',
    extract: (r) => r?.result?.llm || {},
    okField: 'ok',
    openType: 'APP_OPEN_LLM_AUTH',       // remediation: send message
    severity: 'warn'                     // warn = degrades but does not block
  },
  {
    id: 'worker',
    label: 'Service Worker',
    checkType: 'APP_GET_OPTIONS',         // lightweight ping
    extract: () => ({ ok: true }),
    severity: 'critical'
  }
];
```

### Design rules

- **One entry per logical service** even if they share a cookie domain.
- **Severity drives UI ordering.** Critical = red, first; warn = yellow, after.
- **Remediation fields are mutually exclusive:** `openUrl` OR `openType`, never both.
- **`extract` normalizes** diverse response shapes into `{ ok, message?, ... }`.

## 2. Shared helpers

```js
function clearChildren(el) {
  while (el.firstChild) el.removeChild(el.firstChild);
}

async function sendMessage(type, payload = {}) {
  return new Promise((resolve) => {
    chrome.runtime.sendMessage({ type, ...payload }, (response) => {
      if (chrome.runtime.lastError) {
        resolve({ ok: false, error: chrome.runtime.lastError.message });
        return;
      }
      resolve(response || { ok: false });
    });
  });
}
```

**Always check `chrome.runtime.lastError`** before reading the response. Failing to do so causes silent uncaught errors.

## 3. Auth probing strategies

### 3a. Cookie-based detection (fastest, no network)

```js
async function hasCookieSession(urls, cookieNames) {
  for (const name of cookieNames) {
    for (const url of urls) {
      try {
        const cookie = await chrome.cookies.get({ url, name });
        if (cookie) return true;
      } catch { /* cookie API unavailable or permission missing */ }
    }
  }
  return false;
}
```

**Requires:** `"permissions": ["cookies"]` and matching `"host_permissions"`.
**Limitation:** Cookie presence does not guarantee a valid server-side session. Pair with a lightweight API probe when accuracy matters.

### 3b. Cookie expiry inspection

```js
async function getEarliestCookieExpiry(urls, cookieNames) {
  let earliest = Infinity;
  let found = false;
  for (const name of cookieNames) {
    for (const url of urls) {
      try {
        const cookie = await chrome.cookies.get({ url, name });
        if (!cookie) continue;
        found = true;
        if (cookie.expirationDate) {
          const ms = cookie.expirationDate * 1000; // API uses seconds, not ms
          if (ms < earliest) earliest = ms;
        }
      } catch { /* skip */ }
    }
  }
  return found && earliest < Infinity ? earliest : null;
}
```

**Key detail:** `cookie.expirationDate` is in **seconds** since epoch. Always multiply by 1000 before comparing with `Date.now()`.

### 3c. Bearer token + cookie fallback

```js
async function probeAuthentication({ token = '', cookieUrls, cookieNames }) {
  const hasCookie = await hasCookieSession(cookieUrls, cookieNames).catch(() => false);
  const hasToken = Boolean(token);
  if (hasCookie || hasToken) {
    return { ok: true, code: 'ok', hasCookie, hasToken };
  }
  return { ok: false, code: 'auth_required', message: 'Sign in or provide a token.' };
}
```

**Auth order precedence:** Cookie-backed fetch first, then bearer-token fallback. This avoids unnecessary token exposure when a valid browser session exists.

### 3d. Tab-session detection

```js
async function findServiceTab(endpoint) {
  if (!endpoint) return null;
  const origin = new URL(endpoint).origin;
  const tabs = await chrome.tabs.query({ url: origin + '/*' });
  return tabs[0] || null;
}
```

### 3e. Indirect probing

When no dedicated auth check exists, send a lightweight message and interpret the result.

```js
async function probeIndirect(messageType, payload = {}) {
  try {
    const result = await sendMessage(messageType, payload);
    return { ok: !result?.error, detail: result?.error ? 'Sign in required' : 'Session active' };
  } catch {
    return { ok: false, detail: 'Could not verify' };
  }
}
```

## 4. Auth state machine

```
UNKNOWN -> CHECKING -> AUTHENTICATED -> EXPIRING -> EXPIRED
              |                            |           |
              v                            v           v
           FAILED <------------------------+-----------+
              |
              v
        REMEDIATING -> CHECKING (re-probe)
```

| State | UI | Trigger in |
|---|---|---|
| UNKNOWN | Gray | Initial load |
| CHECKING | Spinner | Startup, alarm, storage change |
| AUTHENTICATED | Green | Probe ok, no imminent expiry |
| EXPIRING | Yellow | ok but within warning window |
| EXPIRED | Red | Expiry timestamp in past |
| FAILED | Red | Probe not-ok or threw |
| REMEDIATING | Blue | User clicked fix; re-checks after 3-5s |

## 5. Proactive expiry tracking

```js
const AUTH_EXPIRY_KEY = 'app_auth_expiry_v1';
const AUTH_EXPIRY_WARNING_MS = 15 * 60 * 1000;

async function recordAuthExpiry(id, { expiresAt, label, refreshUrl }) {
  const stored = await chrome.storage.local.get(AUTH_EXPIRY_KEY);
  const map = stored?.[AUTH_EXPIRY_KEY] || {};
  map[id] = { expiresAt, label, refreshUrl, checkedAt: Date.now() };
  await chrome.storage.local.set({ [AUTH_EXPIRY_KEY]: map });
}

async function evaluateExpiry(id, expiresAt, { label, refreshUrl }) {
  await recordAuthExpiry(id, { expiresAt, label, refreshUrl });
  const remainingMs = expiresAt - Date.now();
  if (remainingMs <= 0) {
    await recordAuthFailure(id + '-expired', {
      label: label + ' expired', detail: 'Sign in again.',
      action: 'Sign in', url: refreshUrl, severity: 'critical'
    });
  } else if (remainingMs <= AUTH_EXPIRY_WARNING_MS) {
    const mins = Math.ceil(remainingMs / 60000);
    await recordAuthFailure(id + '-expiring', {
      label: `${label} expiring in ${mins} min`, detail: 'Refresh to extend session.',
      action: 'Refresh', url: refreshUrl, severity: 'warn'
    });
  } else {
    await clearAuthFailure(id + '-expiring');
    await clearAuthFailure(id + '-expired');
  }
}
```

**Rules:** 15-min default warning window. Clear stale warnings on refresh. Store in `chrome.storage.local` (survives worker restarts).

## 6. Auth failure recording

```js
const AUTH_FAILURE_KEY = 'app_auth_failures_v1';
const FAILURE_TTL_MS = 300_000;

async function recordAuthFailure(id, { label = '', detail = '', action = '',
  url = '', severity = 'warn', tabId = null } = {}) {
  const stored = await chrome.storage.local.get(AUTH_FAILURE_KEY);
  const map = stored?.[AUTH_FAILURE_KEY] || {};
  map[id] = { label, detail, action, url, severity, tabId, ts: Date.now() };
  await chrome.storage.local.set({ [AUTH_FAILURE_KEY]: map });
}

async function clearAuthFailure(id) {
  const stored = await chrome.storage.local.get(AUTH_FAILURE_KEY);
  const map = stored?.[AUTH_FAILURE_KEY] || {};
  delete map[id];
  await chrome.storage.local.set({ [AUTH_FAILURE_KEY]: map });
}
```

**Staleness rule:** Skip entries older than `FAILURE_TTL_MS` when reading, to prevent stale banners.

## 7. Service worker polling with alarms

Use `chrome.alarms` (survives worker suspension) instead of `setInterval` (killed on suspension).

```js
// Register on both lifecycle events to cover installs and restarts
chrome.runtime.onInstalled.addListener(async () => {
  await chrome.alarms.create('app-refresh', { periodInMinutes: 5 });
  await runStartupAuthChecks('onInstalled');
});
chrome.runtime.onStartup.addListener(async () => {
  await chrome.alarms.create('app-refresh', { periodInMinutes: 5 });
  await runStartupAuthChecks('onStartup');
});
chrome.alarms.onAlarm.addListener(async (alarm) => {
  if (alarm.name === 'app-refresh') void checkCookieExpiry();
});
```

**Rules:** Minimum period is 30s in MV3; use 1-5 min for production. Piggyback auth checks onto existing alarms.

## 8. Auth bootstrap sequence

Cache startup probe results in session storage to skip redundant checks.

```js
const AUTH_BOOTSTRAP_KEY = 'app_auth_bootstrap_v1';
async function runStartupAuthChecks(source) {
  const version = chrome.runtime.getManifest().version;
  const cache = (await chrome.storage.session.get(AUTH_BOOTSTRAP_KEY))?.[AUTH_BOOTSTRAP_KEY] || {};
  if (cache.version === version) return cache.results;
  // probeEndpoint() dispatches to per-service probes from Section 3
  const results = await Promise.all(
    AUTH_ENDPOINTS.filter(e => e.checkType).map(e => probeEndpoint(e))
  );
  await chrome.storage.session.set({
    [AUTH_BOOTSTRAP_KEY]: { version, checkedAt: new Date().toISOString(), results }
  });
  return results;
}
```

Use `chrome.storage.session` (cleared on restart) for bootstrap cache; `chrome.storage.local` for durable failure/expiry records.

## 9. Cross-context sync via storage.onChanged

```js
chrome.storage.onChanged.addListener((changes, area) => {
  if (area === 'local' && (changes[AUTH_FAILURE_KEY] || changes[AUTH_EXPIRY_KEY])) {
    renderAuthStatus(); // service worker writes, UI contexts re-render reactively
  }
});
```

## 10. Auth status UI rendering

One renderer powers both floating widgets and dashboard panels. The `variant` param controls CSS class prefix.

```js
function renderAuthStatus(container, issues, { variant = 'float' } = {}) {
  if (issues.length === 0) { container.hidden = true; return; }
  container.hidden = false;
  issues.sort((a, b) =>
    (a.severity === 'critical' ? 0 : 1) - (b.severity === 'critical' ? 0 : 1)
  );

  clearChildren(container);
  for (const issue of issues) {
    const row = document.createElement('div');
    row.className = `auth-${variant}-row is-${issue.severity}`;

    const dot = document.createElement('span');
    dot.className = 'auth-dot';
    dot.textContent = issue.severity === 'critical' ? '✗' : '⚠';
    dot.style.color = issue.severity === 'critical' ? '#ef4444' : '#eab308';

    const info = document.createElement('div');
    const labelEl = document.createElement('div');
    labelEl.className = 'auth-label';
    labelEl.textContent = issue.label;
    const detailEl = document.createElement('div');
    detailEl.className = 'auth-detail';
    detailEl.textContent = issue.detail;
    info.append(labelEl, detailEl);

    const btn = document.createElement('button');
    btn.className = 'auth-action';
    btn.textContent = issue.action;
    btn.addEventListener('click', () => handleRemediation(issue, btn));

    row.append(dot, info, btn);
    container.appendChild(row);
  }
}

async function handleRemediation(issue, btn) {
  if (issue.tabId) {
    try { await chrome.tabs.update(issue.tabId, { active: true }); return; }
    catch { /* tab may have closed */ }
  }
  if (issue.url) { window.open(issue.url, '_blank'); return; }
  if (issue.messageType) {
    await sendMessage(issue.messageType);
    btn.textContent = 'Opening...';
    setTimeout(() => { btn.textContent = issue.action; }, 3000);
  }
}
```

**Floating variant:** Fixed-position container with dismiss button; use `storage.onChanged` for re-checks.
**Dashboard variant:** Toggle on button click; add per-service "Recheck" buttons.

## 11. Anti-patterns

| Anti-pattern | Why it fails | Fix |
|---|---|---|
| `setInterval` for auth polling | Worker suspension kills timers | `chrome.alarms` |
| Auth state only in memory | Lost on worker restart | `chrome.storage.local` |
| Sequential service checks | One timeout blocks all | `Promise.all` |
| Stale failure banners | User fixed auth externally | TTL on failure records (5 min) |
| Auto-opening auth pages | Browsers block popups | Only open on user click |
| Hardcoded cookie names | Names change across environments | Configurable list |
| Ignoring `lastError` | Silent uncaught errors | Always check before reading response |

## 12. CSS patterns

```css
.auth-row, .auth-float-row { display: flex; align-items: center; gap: 10px; padding: 8px 0; }
.auth-float { position: fixed; bottom: 16px; right: 16px; z-index: 9999;
  background: #1a1a2e; border: 1px solid #333; border-radius: 10px;
  padding: 12px 16px; max-width: 380px; }
.is-critical { border-left: 3px solid #ef4444; padding-left: 8px; }
.is-warn { border-left: 3px solid #eab308; padding-left: 8px; }
.auth-action { background: #3b82f6; color: #fff; border: none;
  border-radius: 4px; padding: 4px 10px; cursor: pointer; }
```

## 13. Testing

```js
// Vitest — mock chrome APIs and test your probe dispatcher
import { vi, describe, it, expect } from 'vitest';
globalThis.chrome = {
  runtime: { sendMessage: vi.fn(), lastError: null },
  storage: { local: { get: vi.fn(), set: vi.fn() } }
};

describe('probeAuthentication', () => {
  it('returns ok when cookie session exists', async () => {
    // Mock hasCookieSession to return true
    const result = await probeAuthentication({
      token: '', cookieUrls: ['https://app.example.com'],
      cookieNames: ['session']
    });
    expect(result.ok).toBe(true);
  });
});
```

## 14. Checklist for adding a new monitored service

1. Add an entry to `AUTH_ENDPOINTS` with `id`, `label`, `checkType`, `extract`, `okField`, and remediation fields.
2. Implement `probeXxxAuth()` returning `{ ok, code, message }`.
3. Wire into `checkAllAuth()` and `runStartupAuthChecks()`.
4. Call `recordAuthFailure()` / `clearAuthFailure()` from the result handler.
5. If cookie-based, add expiry check and wire to the alarm handler.
6. Add remediation URL/message to the "Open all logins" button.
7. Add a unit test mocking the probe response.
