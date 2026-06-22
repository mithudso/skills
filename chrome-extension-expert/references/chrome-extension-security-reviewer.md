The task asks to return the fixed compressed file content. Here it is:

> **Reference file — part of the `chrome-extension-expert` hub.** Formerly standalone `chrome-extension-security-reviewer` skill.
> Sibling topics now reference files under hubs (`chrome-extension-expert`) — **not** standalone skills. Ignore "use the X skill" / `related_skills` / SKIP pointers naming bare sibling skills; load from `references/<name>.md` in owning hub.

---

---
name: chrome-extension-security-reviewer
title: Chrome Extension Security Reviewer
description: |
  Focused MV3 extension security review reference for auditing Chrome extension code: context-boundary mistakes, permission/privacy overreach, content-script exposure, unsafe resource exposure, message validation, and secret/data handling.
  TRIGGER: security audit of a Chrome MV3 extension; reviewing content-script trust boundaries; manifest permission overreach; web_accessible_resources exposure; runtime.sendMessage validation; chrome.storage secret handling; MAIN-world script window.* exposure; externally_connectable misconfiguration; CSP/remote-code audit; reviewing service worker, bridge, popup, or sidepanel security posture.
  SKIP: platform adapter logic review without a security angle (use platform-adapter-reviewer); vanilla JS UI review with no security concern (use vanilla-js-ui-reviewer); general web application security not specific to Chrome extensions (use security-reviewer); LLM prompt injection in extension context (use llm-integration-reviewer).
category: developer
version: "1.1.0"
updated: "2026-05-29"
keywords:
  - Chrome extension
  - MV3
  - security
  - content script
  - MAIN world
  - ISOLATED world
  - trust boundary
  - permissions
  - host permissions
  - web_accessible_resources
  - runtime.sendMessage
  - chrome.storage
  - CSP
  - service worker
  - secret handling
  - externally_connectable
when_to_use:
  - "security audit of Chrome extension"
  - "content script trust boundary review"
  - "manifest permission overreach"
  - "web_accessible_resources too broad"
  - "runtime.sendMessage validation"
  - "chrome.storage secret handling"
  - "window.* exposure in MAIN world script"
  - "externally_connectable misconfiguration"
  - "CSP violation in MV3 extension"
  - "service worker privilege escalation"
related_skills:
  - platform-adapter-reviewer
  - vanilla-js-ui-reviewer
  - chrome-mv3-advanced
  - chrome-storage-patterns
origin: local
---

# Chrome Extension Security Reviewer

Focused MV3 extension security review: context-boundary mistakes, permission/privacy overreach, content-script exposure, unsafe resource exposure, message validation, secret/data handling.

## How to use this skill

Start from bundled context. Defer to cited official docs for exact APIs and edge-case behavior. If request falls outside Chrome extension security review, use different skill.

**Sources of truth:**
- **Chrome extension security/privacy docs** — MV3-specific controls and threat surfaces (primary)
- **MDN** — cross-browser web-platform semantics only where Chrome docs narrower

**Version note:** based on official pages accessed 2026-05-10, framed for this repo's MV3 architecture with MAIN-world adapters, ISOLATED bridges, service worker, sidepanel.

---

## Source scope

- **Content-script execution contexts, isolated worlds, API access:** [Chrome content scripts](https://developer.chrome.com/docs/extensions/develop/concepts/content-scripts), [MDN content scripts](https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/Content_scripts)
- **Messaging semantics and async-response rules:** [Chrome messaging](https://developer.chrome.com/docs/extensions/develop/concepts/messaging), [MDN runtime.onMessage](https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/API/runtime/onMessage)
- **Permissions, host permissions, activeTab:** [Chrome declare permissions](https://developer.chrome.com/docs/extensions/develop/concepts/declare-permissions), [Chrome activeTab](https://developer.chrome.com/docs/extensions/develop/concepts/activeTab), [MDN host_permissions](https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/manifest.json/host_permissions)
- **Storage, secret handling, CSP, remote-code rules, WAR, external connectivity:** [Chrome storage API](https://developer.chrome.com/docs/extensions/reference/api/storage), [Chrome storage and cookies](https://developer.chrome.com/docs/extensions/develop/concepts/storage-and-cookies), [Chrome CSP manifest](https://developer.chrome.com/docs/extensions/reference/manifest/content-security-policy), [Chrome improve security](https://developer.chrome.com/docs/extensions/develop/migrate/improve-security), [Chrome WAR](https://developer.chrome.com/docs/extensions/reference/manifest/web-accessible-resources), [Chrome externally connectable](https://developer.chrome.com/docs/extensions/reference/manifest/externally-connectable)
- **Repo-specific framing:** `docs/ARCHITECTURE.md`, `CLAUDE.md`, `manifest.json`, content-script/bridge/service-worker code

## Quick review rules

1. **Start from trust boundary.** Content scripts, page scripts, bridges, service worker: unequal trust and capability — data from less-privileged or page-facing contexts = attacker-controlled.
2. **Minimize permissions and host permissions.** Every persistent permission expands blast radius if extension or message path compromised.
3. **Never expose reusable privileged helpers or secrets on `window.*`.** Page can read them — prefer narrow request/response flows.
4. **Validate runtime messages by sender, shape, and intent.** Messaging not a trust boundary itself; Chrome warns hostile inputs flow through extension messaging.
5. **Treat `web_accessible_resources` as security exposure, not packaging.** Expose minimum files to minimum origins.
6. **Don't store secrets where page code or long-lived disk state can expose them.** `chrome.storage` not encrypted secret vault by default.
7. **Prefer HTTPS and least-privilege networking.** Broad host permissions + page-facing message paths = common extension security footgun.

## Review workflow

1. **Identify contexts and privilege levels.** Classify code as service worker, bridge, content script, page-injected MAIN script, popup, or sidepanel.
2. **Map data/control flow across boundaries.** Check what page-originated data reaches `chrome.runtime.sendMessage`, storage, network calls, or privileged APIs.
3. **Review manifest privilege surface.** Audit permissions, host permissions, WAR entries, `externally_connectable` vs actual need.
4. **Audit message validation and command shape.** Look for commands passing raw URLs, selectors, payloads, or action names from weakly-trusted contexts into privileged code.
5. **Audit secret and privacy boundaries.** Review storage usage, auth-token handling, whether sensitive values cross into page-readable contexts.
6. **Review CSP / remote code / resource exposure.** Confirm no MV2-era remote-code patterns or unnecessary page-exposed assets.

## Review surfaces and checks

| Surface | Purpose | Review focus | Caveats |
|---|---|---|---|
| Content script / isolated world | Page DOM access with limited extension APIs | Page-facing data trust, direct DOM sinks, bridge handoff | Isolated world ≠ safe data |
| MAIN-world injected script | Patch page JS APIs | `window.*` exposure, page-readable state, auth leakage | Page can read globals and reusable helpers |
| `runtime.sendMessage` / `onMessage` | Cross-context command path | Sender checks, schema validation, narrow command vocabulary | Async response handling has Chrome-version caveats |
| Manifest permissions / host permissions | Privilege declaration | Least privilege, optional vs persistent access, warning surface | Broad host permissions amplify any other bug |
| `activeTab` | Temporary user-invoked tab privilege | Whether it can replace persistent hosts for a feature | Only applies after user gesture |
| `web_accessible_resources` | Expose extension files to page origins | Minimal origin scoping, no unnecessary scripts/assets | Exposed files aid fingerprinting and abuse |
| `externally_connectable` | Allow external pages/extensions to message extension | Tight allowlists, necessity check | Over-broad allowlists = external attack surface |
| `chrome.storage.*` | Persist settings/data | Secret placement, page accessibility assumptions, retention | Not a secret manager by default |

## Standards and best practices

- Treat **less-privileged contexts as untrusted inputs** — keep privileged ops in service worker or other extension-owned code.
- Prefer **optional permissions or `activeTab`** when feature doesn't need permanent ambient access.
- Keep **manifest and resource exposure minimal** — unused keys and exposed assets are review findings.
- Keep **message protocols explicit and narrow** — enumerated actions with structured payload validation over open-ended "execute this" or "fetch that URL" shapes.
- Preserve repo invariant: **no new data or capabilities exposed on `window.*` from MAIN-world scripts** (`docs/ARCHITECTURE.md`, `CLAUDE.md`).

## Known ambiguities

- Chrome docs primary; MDN useful for cross-browser semantics but don't override Chrome MV3-specific behavior.
- Some extension security issues architectural not code-local — harmless-looking message handler still dangerous depending on who can invoke it.
- `storage.session` reduces persistence risk but doesn't remove boundary review — critical question still which context can read or forward the value.