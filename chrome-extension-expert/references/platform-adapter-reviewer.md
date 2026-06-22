<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly standalone `platform-adapter-reviewer` skill.
> Sibling topics now reference files under hubs (`chrome-extension-expert`) — **not** standalone skills. Ignore "use the X skill" / `related_skills` / SKIP pointers naming bare sibling skills; load `references/<name>.md` from owning hub (see hub's "Cross-hub map").

---

---
name: platform-adapter-reviewer
title: Platform Adapter Reviewer
description: |
  Review ref for browser-page adapters that intercept third-party traffic and bridge normalized data into Chrome extension.
  TRIGGER: reviewing content-script adapters; fetch/XHR/WebSocket interception in MAIN world; ISOLATED-to-MAIN bridge messaging; CustomEvent or postMessage relay patterns; SPA navigation hooks (pushState/popstate); MutationObserver polling for DOM churn; WAR (web_accessible_resources) exposure audit; adapter resilience to host-site payload drift; reviewing Sniffies/Grindr/Adam4Adam/Gmail/Yahoo adapters in this repo.
  SKIP: general Chrome extension security audit (use security-reviewer or chrome-extension-expert); vanilla JS UI review with no adapter/interception concern (use frontend-ui); backend API integration with no browser content-script layer.
category: developer
version: "1.1.1"
updated: "2026-05-31"
keywords:
  - platform adapter
  - content script
  - MAIN world
  - ISOLATED world
  - fetch interception
  - XHR interception
  - WebSocket interception
  - CustomEvent
  - postMessage
  - MutationObserver
  - SPA navigation
  - pushState
  - popstate
  - web accessible resources
  - bridge messaging
  - payload drift
when_to_use:
  - "review this content-script adapter"
  - "fetch interception in MAIN world"
  - "bridge between MAIN and ISOLATED world"
  - "CustomEvent relay pattern"
  - "postMessage across world boundary"
  - "MutationObserver polling teardown"
  - "SPA navigation hook"
  - "web_accessible_resources audit"
  - "adapter breaking after host site update"
  - "WebSocket frame parsing in extension"
related_skills:
  - chrome-extension-expert
  - frontend-ui
origin: local
---

# Platform Adapter Reviewer

Review ref for browser-page adapters that intercept third-party traffic and bridge normalized data into Chrome extension.

## How to use this skill

Start from bundled context below. Defer to cited official docs for exact APIs and edge-case behavior. If request falls outside platform adapter review, choose different skill.

**Sources of truth:**
- **Chrome extension docs** — content-script contexts, script injection, permissions, messaging, WAR
- **MDN** — `CustomEvent`, `postMessage`, Fetch/XHR/WebSocket semantics, `MutationObserver`, SPA navigation
- **Repo architecture docs/code** — patterns applied in repo

**Version note:** official pages accessed 2026-05-10; framed for repo's MAIN-world adapter + ISOLATED-world bridge pattern across Sniffies, Grindr, Adam4Adam, Gmail, Yahoo.

---

## Source scope

- **Execution contexts, injection, messaging, permissions:** Chrome docs — content scripts, scripting, messaging, permissions, WAR ([Chrome content scripts](https://developer.chrome.com/docs/extensions/develop/concepts/content-scripts), [Chrome messaging](https://developer.chrome.com/docs/extensions/develop/concepts/messaging), [Chrome WAR](https://developer.chrome.com/docs/extensions/reference/manifest/web-accessible-resources))
- **Cross-world and web-platform API semantics:** MDN — `CustomEvent`, `postMessage`, Fetch API, `Response.clone()`, `XMLHttpRequest`, `WebSocket`, `MutationObserver`, `pushState()`, `popstate` ([MDN CustomEvent](https://developer.mozilla.org/en-US/docs/Web/API/CustomEvent), [MDN Fetch API](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API), [MDN MutationObserver](https://developer.mozilla.org/en-US/docs/Web/API/MutationObserver))
- **Repo-specific patterns:** `docs/ARCHITECTURE.md`, manifest entries, content scripts, bridges, service-worker handlers

## Quick review rules

1. **Patch page APIs only in MAIN world.** Fetch/XHR/WebSocket interception in isolated world won't see page's network stack.
2. **Keep page response semantics intact.** Clone/read safely; never consume or mutate page's original response objects.
3. **Use narrowest bridge for boundary.** `CustomEvent` for same-window handoff; `postMessage` only crossing world boundaries; `chrome.runtime.sendMessage` only from extension-owned contexts.
4. **Never expose raw credentials or reusable privileged helpers on `window.*`.** MAIN-world convenience globals are page-readable → security bugs.
5. **Treat SPA navigation and DOM observation as correctness code, not glue.** `pushState`, `popstate`, mutation observers, polling loops — all bug-prone on third-party apps.
6. **Treat selector and payload drift as normal.** Resilient adapters need conservative parsers, fallbacks, failure reporting; host app can change without notice.
7. **Keep host/resource exposure minimal.** Page-injected files require WAR exposure; review as part of security surface.

## Review workflow

1. **Map context split.** Identify what runs in MAIN, ISOLATED, and what crosses service worker.
2. **Review interception points.** Check fetch/XHR/WebSocket patching: ordering, cloning, parse guards, host-page behavior preservation.
3. **Review bridge messaging.** Check serialization, message-type validation, context invalidation handling, bridge power leakage across boundary.
4. **Review SPA/DOM survival logic.** Check navigation hooks, observers, polling, selector fallback against likely host-app churn.
5. **Review privacy/security exposure.** Audit `window.*`, localStorage reads, WAR entries, host permissions together.
6. **Review platform-specific resilience.** Check adapter handles partial payloads, repeated events, reconnection, throttling, context reloads.

## Review surfaces and checks

| Surface | Purpose | Review focus | Caveats |
|---|---|---|---|
| MAIN-world adapter | Intercept page-owned APIs, parse data | Patch timing, response cloning, no page breakage, no secret leakage | Page can inspect globals and events |
| ISOLATED bridge | Relay to/from extension APIs | Message narrowing, context invalidation, minimal DOM assumptions | Bridge code often becomes accidental logic hub |
| `CustomEvent` relay | Same-window message handoff | Structured payload shape, no implicit trust | Doesn't cross every boundary devs assume |
| `postMessage` relay | Cross-world/page message handoff | Target/origin assumptions, payload filtering | Safe only scoped to actual boundary in use |
| Fetch/XHR interception | Capture API payloads | Clone-before-read, content-type checks, original response untouched | Easy to break page behavior by consuming bodies |
| WebSocket interception | Capture live events | Framing/parsing resilience, reconnection assumptions | Host apps can change opaque frame formats |
| MutationObserver / polling | Survive SPA DOM churn | Scope, debounce, teardown, selector drift | Overbroad observers/polls become CPU leaks |
| SPA navigation hooks | Sync route state | `pushState` vs `popstate`, fallback behavior | `pushState()` does not itself fire `popstate` |

## Standards and best practices

- Prefer **explicit MAIN vs ISOLATED separation** — each side responsible only for what context can safely do.
- Preserve host page behavior: interception code **transparent first, observant second**.
- Keep **bridge protocols typed and narrow**. Message names reflect intent; don't expose generic "run arbitrary adapter command" backchannel.
- Make **context invalidation first-class** in long-lived bridges and panels — extension reloads happen during dev and updates.
- Repo: MAIN-world scripts keep page-only auth/state in closure scope; bridges request narrow operations, not raw auth reads from page (`docs/ARCHITECTURE.md`).

## Known ambiguities

- Third-party sites don't publish stable contracts for private APIs. Assume host drift inevitable.
- `postMessage('*')` acceptable for same-page world bridging in tightly-scoped extension context, but don't generalize into cross-origin messaging guidance.
- Adapter review is both correctness and security discipline: parser fallback that "mostly works" can still be privacy/boundary bug if it broadens what crosses contexts.