<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly the standalone `platform-adapter-reviewer` skill.
> Sibling topics in this family are now reference files under the hubs (`chrome-extension-expert`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: platform-adapter-reviewer
title: Platform Adapter Reviewer
description: |
  Practical review reference for browser-page adapters that intercept third-party-site traffic and bridge normalized data into a Chrome extension.
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

Practical review reference for browser-page adapters that intercept third-party-site traffic and bridge normalized data into a Chrome extension.

## How to use this skill

Start from the bundled context below. Defer to the cited official documentation for exact APIs and edge-case behavior. If the request falls outside platform adapter review, choose a more appropriate skill.

**Sources of truth:**
- **Chrome extension docs** — content-script contexts, script injection, permissions, messaging, WAR
- **MDN** — `CustomEvent`, `postMessage`, Fetch/XHR/WebSocket semantics, `MutationObserver`, SPA navigation
- **Repo architecture docs/code** — how these patterns are applied in this repo

**Version note:** based on official pages accessed 2026-05-10, framed for this repo's MAIN-world adapter + ISOLATED-world bridge pattern across Sniffies, Grindr, Adam4Adam, Gmail, and Yahoo.

---

## Source scope

- **Execution contexts, injection, messaging, permissions:** Chrome docs for content scripts, scripting, messaging, permissions, and WAR ([Chrome content scripts](https://developer.chrome.com/docs/extensions/develop/concepts/content-scripts), [Chrome messaging](https://developer.chrome.com/docs/extensions/develop/concepts/messaging), [Chrome WAR](https://developer.chrome.com/docs/extensions/reference/manifest/web-accessible-resources))
- **Cross-world and web-platform API semantics:** MDN docs for `CustomEvent`, `postMessage`, Fetch API, `Response.clone()`, `XMLHttpRequest`, `WebSocket`, `MutationObserver`, `pushState()`, `popstate` ([MDN CustomEvent](https://developer.mozilla.org/en-US/docs/Web/API/CustomEvent), [MDN Fetch API](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API), [MDN MutationObserver](https://developer.mozilla.org/en-US/docs/Web/API/MutationObserver))
- **Repo-specific patterns:** `docs/ARCHITECTURE.md`, manifest entries, content scripts, bridges, and service-worker handlers

## Quick review rules

1. **Patch page APIs only in the MAIN world.** Fetch/XHR/WebSocket interception in the isolated world will not see the page's own network stack.
2. **Keep page response semantics intact.** Clone/read safely; never consume or mutate the page's original response objects.
3. **Use the narrowest bridge that fits the boundary.** `CustomEvent` for same-window handoff; `postMessage` only when crossing world boundaries; `chrome.runtime.sendMessage` only from extension-owned contexts.
4. **Never expose raw credentials or reusable privileged helpers on `window.*`.** MAIN-world convenience globals are page-readable and become security bugs.
5. **Treat SPA navigation and DOM observation as correctness code, not glue.** `pushState`, `popstate`, mutation observers, and polling loops are all bug-prone on third-party apps.
6. **Treat selector and payload drift as normal.** Resilient adapters need conservative parsers, fallbacks, and failure reporting because the host app can change without notice.
7. **Keep host/resource exposure minimal.** Page-injected files require WAR exposure and should be reviewed as part of the security surface.

## Review workflow

1. **Map the context split.** Identify what runs in MAIN, what runs in ISOLATED, and what crosses through the service worker.
2. **Review interception points.** Check fetch/XHR/WebSocket patching for ordering, cloning, parse guards, and preservation of host-page behavior.
3. **Review bridge messaging.** Check serialization, message-type validation, context invalidation handling, and whether the bridge leaks too much power across the boundary.
4. **Review SPA/DOM survival logic.** Check navigation hooks, observers, polling, and selector fallback logic against likely host-app churn.
5. **Review privacy/security exposure.** Audit `window.*`, localStorage reads, WAR entries, and host permissions together.
6. **Review platform-specific resilience.** Check whether the adapter handles partial payloads, repeated events, reconnection, throttling, and context reloads gracefully.

## Review surfaces and checks

| Surface | Purpose | Review focus | Caveats |
|---|---|---|---|
| MAIN-world adapter | Intercept page-owned APIs and parse data | Patch timing, response cloning, no page breakage, no secret leakage | Page can inspect globals and events |
| ISOLATED bridge | Relay to/from extension APIs | Message narrowing, context invalidation, minimal DOM assumptions | Bridge code often becomes accidental logic hub |
| `CustomEvent` relay | Same-window message handoff | Structured payload shape, no implicit trust | Does not cross every boundary developers assume |
| `postMessage` relay | Cross-world/page message handoff | Target/origin assumptions, payload filtering | Safe only when scoped to the actual boundary in use |
| Fetch/XHR interception | Capture API payloads | Clone-before-read, content-type checks, original response untouched | Easy to break page behavior by consuming bodies |
| WebSocket interception | Capture live events | Framing/parsing resilience, reconnection assumptions | Host apps can change opaque frame formats |
| MutationObserver / polling | Survive SPA DOM churn | Scope, debounce, teardown, selector drift | Overbroad observers/polls become CPU leaks |
| SPA navigation hooks | Sync route state | `pushState` vs `popstate`, fallback behavior | `pushState()` does not itself fire `popstate` |

## Standards and best practices

- Prefer **explicit MAIN vs ISOLATED separation** — keep each side responsible only for what that context can safely do.
- Preserve the host page's behavior: interception code should be **transparent first, observant second**.
- Keep **bridge protocols typed and narrow**. Message names should reflect intent, not expose a generic "run arbitrary adapter command" backchannel.
- Make **context invalidation a first-class case** in long-lived bridges and panels — extension reloads happen during development and updates.
- For this repo: MAIN-world scripts keep page-only auth/state inside closure scope; bridges request narrow operations rather than reading raw auth out of the page (`docs/ARCHITECTURE.md`).

## Known ambiguities

- Third-party sites do not publish stable contracts for most private APIs. Assume host drift is inevitable.
- `postMessage('*')` can be acceptable for same-page world bridging in a tightly-scoped extension context, but should not be generalized into cross-origin messaging guidance.
- Adapter review is both a correctness and security discipline: a parser fallback that "mostly works" can still be a privacy or boundary bug if it broadens what crosses contexts.
