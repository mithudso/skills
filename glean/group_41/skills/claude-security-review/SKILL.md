---
name: security-review
version: "1.3.1"
updated: "2026-07-20"
description: >-
  Application & infrastructure security-review hub. TRIGGER: web app / Chrome extension / backend security review — OWASP, ASVS, CSP/CORS, session/auth (security-reviewer); compliance audit — secrets, supply-chain, PII/logging (security-compliance-auditor); security headers — Helmet.js, CSP, HSTS (http-security-headers); Okta platform — Identity Engine, OAuth/OIDC, management APIs (okta-expert); Okta admin/tenant hardening — ThreatInsight, breakglass (okta-admin-hardening); phishing-resistant auth — FastPass, passkeys (okta-phishing-resistant-auth); Zero Trust & device posture — 800-207, device assurance (okta-zero-trust-device); Okta ITDR — MFA fatigue, stolen-cookie replay, Universal Logout (okta-itdr-session-security); OIG — certification campaigns, SoD, SOX evidence (okta-identity-governance). SKIP: Web Crypto/vault review → webcrypto-vault-reviewer; auth-flow design → software-engineering-patterns (references/web-auth-patterns.md); AI-agent identity → agent-identity-authz-payments; LLM/RAG prompt-injection defense → rag-prompt-injection-defense; embedding inversion / vector-store leakage → embedding-inversion-threat-model.
origin: local
metadata:
  changelog:
    - "2026-06-11 sko v1.2.0->v1.3.0 — Pass H 5/10->10/10 pos, 10/10->10/10 neg (predicted); 1 high + 5 medium fixed"
---

# security-review

Application & infrastructure security-review hub.

This hub routes to on-demand reference files under `references/`. Load the matching reference for depth on any spoke.

## Sub-skill routing table

| Topic | Load this reference |
| --- | --- |
| Security review of web apps / Chrome extensions / backend services (OWASP Top 10, ASVS, CSP/CORS, secrets, session/auth) | `references/security-reviewer.md` |
| Codebase security & compliance audit — secrets/credentials, supply-chain, PII/logging, AI-policy, repo governance/SSDLC gates | `references/security-compliance-auditor.md` |
| HTTP security headers — Helmet.js, Content-Security-Policy, HSTS, COOP/COEP, X-Frame-Options | `references/http-security-headers.md` |
| Okta identity platform — Identity Engine, OAuth/OIDC, auth servers, management APIs, platform-level security posture (deep Okta security management → the five okta-* rows below) | `references/okta-expert.md` |
| Okta tenant/admin-console hardening — admin MFA + session binding, HealthInsight, ThreatInsight/zones, custom admin roles, API-token governance, breakglass design, Oct-2023 breach lessons | `references/okta-admin-hardening.md` |
| Phishing-resistant authentication rollout on Okta — FastPass, passkeys (synced vs device-bound, AAGUID policy), authentication method chains, NIST AALs, rollout sequencing, recovery without phishable factors | `references/okta-phishing-resistant-auth.md` |
| Zero Trust with Okta as policy decision point — NIST 800-207 mapping, device assurance/EDR posture signals, Okta Device Access (Desktop MFA), network zones vs device context, ZTNA integration | `references/okta-zero-trust-device.md` |
| Okta ITDR & session security — System Log detections (MFA fatigue, stolen-cookie/AiTM replay, credential stuffing, impossible travel), response loop (clear sessions/revoke/suspend/Universal Logout), Identity Threat Protection, DPoP, HAR sanitization | `references/okta-itdr-session-security.md` |
| Okta Identity Governance — certification campaign design, Access Requests/approvals, entitlement management, SoD rules, SOX/SOC2/ISO audit evidence, OIG-vs-dedicated-IGA boundary, anti-patterns | `references/okta-identity-governance.md` |

<!-- cross-hub-map -->
## Cross-hub map — where every security-review topic lives

All security-review spokes are reference files under this hub. If a topic is not in the routing table above, it belongs to a peer hub — activate that hub or `Read` its `references/<name>.md` directly.

| Hub | Owns | Reference files |
| --- | --- | --- |
| `security-review` | Application & infrastructure security review + Okta security management | `references/security-reviewer.md`, `references/security-compliance-auditor.md`, `references/http-security-headers.md`, `references/okta-expert.md`, `references/okta-admin-hardening.md`, `references/okta-phishing-resistant-auth.md`, `references/okta-zero-trust-device.md`, `references/okta-itdr-session-security.md`, `references/okta-identity-governance.md` |

Fix-loop handoff: a security review that must also apply and verify code fixes across a file/repo to convergence (not review-only) → `code-deep-optimizer`.
