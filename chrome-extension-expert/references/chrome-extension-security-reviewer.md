<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly the standalone `chrome-extension-security-reviewer` skill.
> Sibling topics in this family are now reference files under the hubs (`chrome-extension-expert`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

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

Focused MV3 extension security review reference for context-boundary mistakes, permission/privacy overreach, content-script exposure, unsafe resource exposure, message validation, and secret/data handling.

## How to use this skill

Start from the bundled context below. Defer to cited official documentation for exact APIs and edge-case behavior. If the request falls outside Chrome extension security review, choose a more appropriate skill.

**Sources of truth:**
- **Chrome extension security/privacy docs** — MV3-specific controls and threat surfaces (primary)
- **MDN** — cross-browser web-platform semantics only where Chrome docs are narrower

**Version note:** based on official pages accessed 2026-05-10, framed for this repo's MV3 architecture with MAIN-world adapters, ISOLATED bridges, a service worker, and a sidepanel.

---

## Source scope

- **Content-script execution contexts, isolated worlds, API access:** [Chrome content scripts](https://developer.chrome.com/docs/extensions/develop/concepts/content-scripts), [MDN content scripts](https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/Content_scripts)
- **Messaging semantics and async-response rules:** [Chrome messaging](https://developer.chrome.com/docs/extensions/develop/concepts/messaging), [MDN runtime.onMessage](https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/API/runtime/onMessage)
- **Permissions, host permissions, activeTab:** [Chrome declare permissions](https://developer.chrome.com/docs/extensions/develop/concepts/declare-permissions), [Chrome activeTab](https://developer.chrome.com/docs/extensions/develop/concepts/activeTab), [MDN host_permissions](https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/manifest.json/host_permissions)
- **Storage, secret handling, CSP, remote-code rules, WAR, external connectivity:** [Chrome storage API](https://developer.chrome.com/docs/extensions/reference/api/storage), [Chrome storage and cookies](https://developer.chrome.com/docs/extensions/develop/concepts/storage-and-cookies), [Chrome CSP manifest](https://developer.chrome.com/docs/extensions/reference/manifest/content-security-policy), [Chrome improve security](https://developer.chrome.com/docs/extensions/develop/migrate/improve-security), [Chrome WAR](https://developer.chrome.com/docs/extensions/reference/manifest/web-accessible-resources), [Chrome externally connectable](https://developer.chrome.com/docs/extensions/reference/manifest/externally-connectable)
- **Repo-specific framing:** `docs/ARCHITECTURE.md`, `CLAUDE.md`, `manifest.json`, and content-script/bridge/service-worker code

## Quick review rules

1. **Start from the trust boundary.** Content scripts, page scripts, bridges, and the service worker do not have equal trust or capability — treat data from less-privileged or page-facing contexts as attacker-controlled.
2. **Minimize permissions and host permissions.** Every persistent permission expands blast radius if the extension or a message path is compromised.
3. **Never expose reusable privileged helpers or secrets on `window.*`.** The page can read them — prefer narrow request/response flows instead.
4. **Validate runtime messages by sender, shape, and intent.** Messaging is not a trust boundary by itself; Chrome explicitly warns that hostile inputs can flow through extension messaging.
5. **Treat `web_accessible_resources` as a security exposure, not just packaging.** Expose only the minimum files to the minimum origins.
6. **Do not store secrets in places page code or long-lived disk state can casually expose.** `chrome.storage` is not an encrypted secret vault by default.
7. **Prefer HTTPS and least-privilege networking.** Broad host permissions combined with page-facing message paths are a common extension security footgun.

## Review workflow

1. **Identify contexts and privilege levels.** Classify code as service worker, bridge, content script, page-injected MAIN script, popup, or sidepanel.
2. **Map data/control flow across boundaries.** Check what page-originated data can reach `chrome.runtime.sendMessage`, storage, network calls, or privileged APIs.
3. **Review manifest privilege surface.** Audit permissions, host permissions, WAR entries, and `externally_connectable` settings against actual need.
4. **Audit message validation and command shape.** Look for commands that pass raw URLs, selectors, payloads, or action names from weakly-trusted contexts into privileged code.
5. **Audit secret and privacy boundaries.** Review storage usage, auth-token handling, and whether sensitive values ever cross into page-readable contexts.
6. **Review CSP / remote code / resource exposure.** Confirm no MV2-era remote-code patterns or unnecessary page-exposed assets have crept in.

## Review surfaces and checks

| Surface | Purpose | Review focus | Caveats |
|---|---|---|---|
| Content script / isolated world | Page DOM access with limited extension APIs | Page-facing data trust, direct DOM sinks, bridge handoff | Isolated world is not the same as safe data |
| MAIN-world injected script | Patch page JS APIs | `window.*` exposure, page-readable state, auth leakage | Page can read globals and reusable helpers |
| `runtime.sendMessage` / `onMessage` | Cross-context command path | Sender checks, schema validation, narrow command vocabulary | Async response handling has Chrome-version caveats |
| Manifest permissions / host permissions | Privilege declaration | Least privilege, optional vs persistent access, warning surface | Broad host permissions amplify impact of any other bug |
| `activeTab` | Temporary user-invoked tab privilege | Whether it can replace persistent hosts for a feature | Only applies after user gesture |
| `web_accessible_resources` | Expose extension files to page origins | Minimal origin scoping, no unnecessary scripts/assets | Exposed files aid fingerprinting and abuse |
| `externally_connectable` | Allow external pages/extensions to message the extension | Tight allowlists, necessity check | Over-broad allowlists create external attack surface |
| `chrome.storage.*` | Persist settings/data | Secret placement, page accessibility assumptions, retention | Not a secret manager by default |

## Standards and best practices

- Treat **less-privileged contexts as untrusted inputs** — keep privileged operations in the service worker or other extension-owned code.
- Prefer **optional permissions or `activeTab`** when a feature does not need permanent ambient access.
- Keep **manifest and resource exposure minimal** — unused keys and exposed assets are themselves review findings.
- Keep **message protocols explicit and narrow** — prefer enumerated actions with structured payload validation over open-ended "execute this" or "fetch that URL" shapes.
- Preserve the repo invariant: **no new data or capabilities should be exposed on `window.*` from MAIN-world scripts** (`docs/ARCHITECTURE.md`, `CLAUDE.md`).

## Known ambiguities

- Chrome docs are primary here; MDN is useful for cross-browser semantics but should not override Chrome MV3-specific behavior.
- Some extension security issues are architectural rather than code-local — a harmless-looking message handler can still be dangerous depending on who can invoke it.
- `storage.session` reduces some persistence risk but does not remove the need for boundary review — the critical question is still which context can read or forward the value.
