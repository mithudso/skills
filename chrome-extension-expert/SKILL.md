---
description: >-
  Chrome extension / MV3 development hub — build, package, secure, test browser extensions. 17 references. TRIGGER: MV3 service-worker lifecycle/idle termination; manifest.json fields, permissions, CRX packaging, Web Store publishing; chrome.runtime/chrome.tabs cross-context messaging; native messaging hosts (Python/Node, 4-byte framing); chrome.identity getAuthToken/launchWebAuthFlow OAuth; chrome.storage local/session/sync; chrome.tabs query/update/create; chrome.notifications; offscreen documents & keepalive; chrome.action badge/icon; content scripts & shadow-DOM overlays; WebSocket/Socket.IO in a service worker; Playwright E2E for unpacked extensions; marked.js + DOMPurify markdown rendering; extension security/permission/CSP review. SKIP: general frontend/UI/CSS → frontend-ui; JS/TS language → lang-js-ts; language-agnostic patterns → software-engineering-patterns.
name: chrome-extension-expert
category: developer
tags:
  - developer
  - chrome
  - extension
  - mv3
whenToUse:
  - "building or debugging a Chrome/Chromium MV3 extension end to end"
  - "MV3 service-worker lifecycle, idle termination, and keepalive questions"
  - "manifest.json fields, permissions, CRX packaging, or Web Store publishing"
  - "cross-context messaging (chrome.runtime / chrome.tabs) or native messaging hosts"
  - "chrome.identity OAuth, chrome.storage surfaces, tabs, notifications, or badge state"
  - "offscreen documents, content scripts, shadow-DOM overlays, or websockets in an extension"
  - "Playwright E2E testing of an unpacked extension or in-browser markdown rendering"
  - "security review of extension code for context-boundary, permission, or CSP issues"
related_skills:
  - frontend-ui
  - programming-languages
  - software-engineering-patterns
version: "1.0.0"
updated: "2026-05-31"
---
# Chrome Extension Expert

Dev hub for Chrome/Chromium MV3 browser extensions. First-choice reference for extension-specific work: MV3 model, service-worker lifecycle, messaging boundaries (worker, content scripts, popup, options, offscreen docs), `chrome.*` API surface (identity, storage, tabs, notifications, action), native-messaging bridges, packaging, Web Store distribution, in-page UI (content scripts, shadow DOM, markdown rendering), websockets in worker, E2E testing unpacked extension, extension security review.

Use when task: how extension built, packaged, wired, or secured — not general web UI, language syntax, or generic architecture (see sibling hubs below).

## How to use this skill

17 Chrome-extension sub-skills as on-demand reference files. Match task to routing table below and **Read listed `references/<name>.md` before answering deep questions** — table alone not enough for depth. For exact API signatures, manifest keys, Chrome version behavior: defer to official Chrome for Developers extension docs.

<!-- ROUTING TABLE: chrome-extension-expert — auto-generated, edit descriptions as needed -->
## Sub-skill routing table

Hub absorbs 17 former standalone skills as on-demand reference files. Task matches row → **Read listed `references/` file** before answering — do not rely on table alone for depth.

| Sub-topic | When to load | Reference file |
| --- | --- | --- |
| `chrome-badge-metrics` | Chrome extension badge and action icon state management — chrome.action API (setBadgeText, | `references/chrome-badge-metrics.md` |
| `chrome-dev` | Chrome MV3 extension development reference — architecture fundamentals, full Chrome API index, | `references/chrome-dev.md` |
| `chrome-extension-packaging` | Chrome extension packaging and distribution -- build pipelines (webpack/vite/CRXJS), manifest validation, Chrome Web Store publishing, version management… | `references/chrome-extension-packaging.md` |
| `chrome-extension-security-reviewer` | Focused MV3 extension security review reference for auditing Chrome extension code: context-boundary mistakes, permission/privacy overreach, content-script… | `references/chrome-extension-security-reviewer.md` |
| `chrome-identity-oauth` | Chrome extension OAuth patterns — chrome.identity API (getAuthToken, | `references/chrome-identity-oauth.md` |
| `chrome-mv3-advanced` | Advanced Chrome MV3 extension APIs — deep reference for offscreen documents, | `references/chrome-mv3-advanced.md` |
| `chrome-native-messaging` | Chrome extension native messaging — writing a native host in Python or Node.js, | `references/chrome-native-messaging.md` |
| `chrome-notifications-patterns` | Chrome extension notification and alert patterns — chrome.notifications API | `references/chrome-notifications-patterns.md` |
| `chrome-offscreen-documents` | Chrome MV3 offscreen documents — creating, closing, silent audio keepalive, | `references/chrome-offscreen-documents.md` |
| `chrome-storage-patterns` | Chrome extension storage patterns — chrome.storage.local vs session vs sync, | `references/chrome-storage-patterns.md` |
| `chrome-tabs-management` | Chrome extension tabs API patterns — chrome.tabs query/get/update/create/remove, | `references/chrome-tabs-management.md` |
| `extension-e2e-testing` | Chrome extension E2E testing with Playwright — launchPersistentContext, unpacked extension loading, | `references/extension-e2e-testing.md` |
| `extension-message-bridge` | Chrome MV3 extension cross-context messaging — chrome.runtime.sendMessage, chrome.tabs.sendMessage, | `references/extension-message-bridge.md` |
| `markdown-rendering-browser` | Browser-side markdown rendering with marked.js, DOMPurify XSS sanitization, | `references/markdown-rendering-browser.md` |
| `mv3-service-worker-expert` | Chrome MV3 extension service worker lifecycle — idle termination (30s), operation | `references/mv3-service-worker-expert.md` |
| `shadow-dom-component-authoring` | Shadow DOM component authoring for Chrome extension overlays and web components. | `references/shadow-dom-component-authoring.md` |
| `websocket-extension-patterns` | WebSocket and Socket.IO patterns in Chrome MV3 extensions — service worker | `references/websocket-extension-patterns.md` |
| `dexie-indexeddb-local-first-reviewer` | Dexie.js and IndexedDB local-first architecture reviewer. Audits schema/index design, | `references/dexie-indexeddb-local-first-reviewer.md` |
| `sticky-notes-local-data` | Chrome extension sticky notes with local data persistence — chrome.storage.local note CRUD | `references/sticky-notes-local-data.md` |
| `live-hub-toolkit` | Real-time hub data discovery patterns for Chrome extensions monitoring support portals — | `references/live-hub-toolkit/SKILL.md` |
| `platform-adapter-reviewer` | Practical review reference for browser-page adapters that intercept third-party-site traffic and bridge normalized data into a Chrome extension. | `references/platform-adapter-reviewer.md` |

## Cross-hub boundaries

Hub owns Chrome-extension-specific development. Hand off when task falls into sibling hub:

- **General frontend / UI / CSS** not extension-specific → `frontend-ui`.
- **General JavaScript / TypeScript language** → `programming-languages`.
- **General software patterns / architecture** → `software-engineering-patterns`.
- **This project's own extension code and runtime** (mdb-tam dashboard extension) → treat as user's repo, not this hub. Use repo's own docs (`docs/ARCHITECTURE.md`, `src/background/README.md`) and source.

Some topics touch two hubs (e.g., shadow-DOM overlay = extension content-script concern here + general UI in `frontend-ui`). Lead with hub matching user's intent — extension-specific stays here.

<!-- cross-hub-map -->
## Cross-hub map — where every chrome-extension topic lives

Family split across hubs. Task's deep material **not** in this hub's sub-skill routing table → it's a reference file under sibling hub below — **activate that hub or `Read` its `references/<name>.md` directly**. Every former standalone skill in family now a reference under one hub (nothing deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `chrome-extension-expert` | Chrome Extension Development (MV3, APIs, packaging, security) | `references/chrome-badge-metrics.md`, `references/chrome-dev.md`, `references/chrome-extension-packaging.md`, `references/chrome-extension-security-reviewer.md`, … |