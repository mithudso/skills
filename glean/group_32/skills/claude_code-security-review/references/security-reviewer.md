<!-- hub-reference-banner -->
> **Reference file — part of the `security-review` hub.** Formerly the standalone `security-reviewer` skill.
> Sibling topics in this family are now reference files under the hubs (`security-review`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: security-reviewer
description: "Security review reference for web apps, Chrome extensions, and backend services. TRIGGER: auditing code for OWASP Top 10 risks; reviewing manifest permissions or host access in a Chrome extension; evaluating CSP headers, CORS config, or same-origin policy enforcement; checking secrets handling, data minimization, or user-data privacy; performing ASVS-style control verification; reviewing authentication, session handling, or supply-chain/release-path security; asking 'is this safe?', 'what permissions does this need?', or 'how do I harden this?'. SKIP: implementing new security features from scratch (use security-architecture); pure performance or functional review with no security angle (use code-reviewer)."
origin: local
version: "1.1"
updated: "2026-05-29"
category: developer
tags: [security, owasp, csp, cors, chrome-extension, privacy, permissions, asvs, xss, least-privilege]
whenToUse:
  - "audit this code for security issues"
  - "review manifest permissions for this Chrome extension"
  - "check the CSP headers on this response"
  - "is this CORS config safe?"
  - "review secrets handling and data minimization"
  - "ASVS-style control verification"
  - "check authentication and session security"
  - "supply chain and release-path security review"
  - "what permissions does this extension actually need?"
  - "harden this against XSS or prompt injection"
whenNotToUse:
  - "implementing new security features from scratch — use security-architecture"
  - "pure performance or functional review — use code-reviewer"
related_skills: [chrome-extension-expert, http-security-headers, software-engineering-patterns, mongodb-operations-expert]
---

# Security Reviewer

## How to use this skill

Treat OWASP as the primary framework for risk framing and verification, MDN as the source for
browser/web-platform security behavior, and the Chrome extension docs as the source for
extension-specific permission and privacy guidance.

Sources used: [OWASP Top 10 2025](https://owasp.org/www-project-top-ten/),
[OWASP Cheat Sheet Series](https://cheatsheetseries.owasp.org/),
[OWASP ASVS](https://owasp.org/www-project-application-security-verification-standard/),
[MDN Web Security](https://developer.mozilla.org/en-US/docs/Web/Security),
[MDN CSP](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP),
[Chrome stay secure](https://developer.chrome.com/docs/extensions/develop/security-privacy/stay-secure),
[Chrome user privacy](https://developer.chrome.com/docs/extensions/develop/security-privacy/user-privacy).

## Review Workflow

1. **Identify the trust boundary** — does the code cross origin, privilege, extension, or user-data boundaries?
2. **Classify the risk** — use OWASP Top 10 to frame the likely issue class.
3. **Verify the control** — use ASVS: what technical control should exist, and can it be tested?
4. **Check platform behavior** — confirm how the browser constrains or exposes the pattern (SOP, CORS, CSP).
5. **Check least privilege and privacy** — review requested permissions, collected data, and whether the design asks for more access than needed.
6. **Check secure defaults** — is the safer path the default, or does correct behavior require a fragile opt-in?

## Quick Rules

1. **Least privilege first** — only grant permissions, access, and data collection needed right now, not future capabilities.
2. **Security and privacy are distinct but linked** — good privacy requires good security; review both angles.
3. **HTTPS only** — flag any HTTP transport dependency; treat HTTP as interceptable.
4. **Check browser boundaries first** — SOP, CORS, and CSP often determine whether a risky pattern is actually exposed.
5. **Use ASVS for control verification + OWASP Top 10 for risk prioritization** — one structures the test, the other prioritizes common classes.
6. **CSP is a real control** — review whether directives are meaningfully restrictive, not just present.
7. **Prefer `activeTab` over `<all_urls>`** — use the narrower permission when it fits.
8. **Minimize post-compromise blast radius** — reduce permissions, roles, and exposed capabilities.

## Controls Inventory

| Control | Purpose | Key config | Mitigates | Common failure |
|---|---|---|---|---|
| OWASP Top 10 | Risk-awareness framework | Current 2025 release | Classify high-risk issue areas | Awareness only — not a complete checklist |
| OWASP ASVS | Verification basis for controls | Rigor level + control requirements | Gaps in implemented controls | Needs project-specific mapping |
| OWASP Cheat Sheet Series | Topic-specific implementation guidance | Pick relevant topic | Common secure-coding gaps | Too broad — pick relevant topics |
| Same-origin policy | Restrict cross-origin interaction by default | Origin boundaries | Unauthorized cross-origin access | Misunderstanding origin = false confidence |
| CORS | Controlled cross-origin access | Origin allowances | Cross-origin misuse | Over-opening for convenience |
| `Content-Security-Policy` header | Enforce resource-loading restrictions | `default-src`, `img-src`, etc. | XSS, risky resource loading | Meta-delivered CSP lacks full feature support |
| HTTPS-only transport | Protect data in transit | Use HTTPS instead of HTTP | Man-in-the-middle | HTTP assumed interceptable |
| Manifest `permissions` minimization | Limit extension privilege | Only required APIs/sites | Excess privilege | Future-proofing requests |
| `activeTab` | Temporary tab-scoped access | User-invoked; ends on tab exit | Broad persistent host access | Not a universal host-permission replacement |
| Data minimization | Limit collected/transmitted data | What is collected, shared, stored, deleted | Privacy abuse | Privacy claims meaningless without security |
| Developer-account protection | Protect publishing account | 2FA, correct member roles | Malicious extension takeover | Secure code insufficient if account is weak |

## Review Standards

### Authentication and session handling
Use ASVS as the basis for verifying auth/session controls. Review any path that crosses
identity, privilege, or user-data boundaries as a security-control path, not business logic.

### Validation and sanitization
Review user-controlled input paths together with rendering/consumption context and applicable
CSP. Use OWASP Cheat Sheet guidance for implementation detail.

### Extension permissions
- Prefer `activeTab` over broad host permissions when it fits the use case.
- Restrict cross-origin fetching to domains explicitly in the manifest.
- Remember: extensions are attractive targets because they have browser privileges.

### Supply-chain and release-path security
- Who can ship or update code matters as much as what the code does.
- Developer-account compromise can push malicious code to all users.

## Practical Defaults

- Start with: **what is the privilege boundary, what data is at risk, what control should exist?**
- For web code: check **SOP / CORS / CSP** before assuming a path is safe or unsafe.
- For extension code: check **manifest permissions, host access, data handling** before implementation details.

## Known Scope Notes

- OWASP Top 10 is an awareness document, not a complete test procedure.
- ASVS needs project-specific interpretation and control mapping.
- Meta-based CSP does not support all features — prefer header-based delivery.
- For implementation depth on any topic, follow the relevant OWASP Cheat Sheet.
