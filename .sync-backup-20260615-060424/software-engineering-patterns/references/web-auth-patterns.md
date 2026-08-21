<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `web-auth-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: web-auth-patterns
version: 2.1.0
last_updated: 2026-05-29
description: >-
  Browser authentication patterns — cookie attributes, OAuth 2.1 + PKCE,
  JWT issuance and refresh rotation, session management, CSRF/XSS token
  protection, Chrome MV3 extension auth, passkeys/WebAuthn/FIDO2, OpenID
  Connect, SAML 2.0, MFA methods, zero trust architecture, token storage
  strategy, Content Security Policy for auth pages, CORS with credentials,
  Device Authorization Grant, and identity provider integration.
  TRIGGER: user sets cookie flags on a session or auth cookie; implements
  OAuth 2.0/2.1 or OIDC; designs JWT issuance/rotation/revocation; chooses
  server-side sessions vs stateless JWT; builds auth inside a Chrome MV3
  extension; protects against CSRF or XSS token theft; handles 401 errors or
  proactive token refresh; needs DPoP sender-constrained tokens; implements
  passkeys/WebAuthn/FIDO2; integrates OIDC for SSO; compares SAML 2.0 vs OIDC;
  implements TOTP or MFA; designs zero trust auth; chooses token storage
  (memory vs HttpOnly cookie vs localStorage); configures CSP for auth pages;
  configures CORS with credentials; implements Device Authorization Grant
  (RFC 8628); hardens JWT against algorithm confusion; integrates Auth0,
  Keycloak, Okta, or Entra ID.
  SKIP: backend IdP-only SAML with no browser component → okta-expert; general
  Chrome extension security audit → chrome-extension-security-reviewer; broad
  codebase security review → security-reviewer; API key management with no
  browser auth → api-design-patterns; HTTP security headers with no auth flow
  → http-security-headers.
category: developer
tags: [auth, cookies, oauth, jwt, sessions, csrf, chrome-extensions, security, passkeys, webauthn, fido2, oidc, saml, mfa, zero-trust, csp, cors]
whenToUse:
  - setting cookie attributes (Secure, HttpOnly, SameSite, __Host- prefix, CHIPS/Partitioned)
  - implementing OAuth 2.0/2.1 authorization code flow with PKCE
  - designing JWT issuance, access/refresh token rotation, or token revocation
  - choosing between server-side sessions and stateless JWT sessions
  - implementing sliding expiry, absolute timeout, or session fixation prevention
  - building authentication inside a Chrome MV3 extension
  - protecting against CSRF or XSS token theft
  - handling 401 errors, token expiry, or proactive refresh loops
  - implementing DPoP (RFC 9449) sender-constrained tokens
  - detecting or reacting to HttpOnly cookie expiry from client-side JavaScript
  - implementing passkeys, WebAuthn, or FIDO2 passwordless authentication
  - integrating OpenID Connect (OIDC) for SSO or identity federation
  - choosing between SAML 2.0 and OIDC for enterprise SSO
  - implementing multi-factor authentication (TOTP, WebAuthn, push)
  - designing zero trust authentication architecture
  - deciding token storage strategy (memory vs HttpOnly cookie vs localStorage)
  - configuring Content Security Policy for auth-related XSS prevention
  - configuring CORS with credentials for cross-origin auth
  - implementing Device Authorization Grant (RFC 8628) for CLI or IoT
  - hardening JWT against algorithm confusion or the 'none' algorithm
  - integrating Auth0, Keycloak, Okta, or Microsoft Entra ID
whenNotToUse:
  - purely backend IdP or SAML with no browser auth component → okta-expert
  - general Chrome extension security audit → chrome-extension-security-reviewer
  - broad codebase security review with no auth focus → security-reviewer
  - API key management with no browser auth component → api-design-patterns
  - HTTP security headers with no auth flow involved → http-security-headers
related_skills:
  - chrome-storage-patterns
  - chrome-extension-security-reviewer
  - okta-expert
  - http-security-headers
  - express-patterns
  - api-design-patterns
---

# Web Authentication Patterns

## Context file

Load and apply patterns from:

- `references/web-auth-patterns-context.md` — cookie mechanics and attributes, PKCE flow walkthrough, JWT structure and lifecycle, session design patterns, Chrome MV3 auth patterns, CSRF/XSS defenses, token recovery flows, 8 decision tables, 25 cited sources (2025).

Check the decision tables in that file before recommending an approach.

---

## Quick Decision Guide

| Question | Answer |
|---|---|
| Cookie attributes for an auth session? | `Secure; HttpOnly; SameSite=Lax; Path=/` + `__Host-` prefix |
| Third-party embedded widget needs cookies? | CHIPS: `SameSite=None; Secure; Partitioned` |
| Which OAuth flow in 2026? | Authorization Code + PKCE (S256) for all clients; implicit flow removed in OAuth 2.1 |
| Token storage in a browser SPA? | Access token in JS memory; refresh token in HttpOnly cookie |
| Chrome extension token storage? | Access tokens in `chrome.storage.session`; refresh tokens (encrypted) in `chrome.storage.local`; `localStorage` unavailable in service workers |
| CSRF defence strategy? | Synchronizer token (stateful app) or signed HMAC double-submit cookie (stateless); SameSite is defence-in-depth only, not a primary control |
| 401 recovery pattern? | Single-flight refresh; on reuse detection expire the token family and force re-authentication |
| Stolen-token protection? | DPoP (RFC 9449) — sender-constrains tokens to a client key pair; supported in Keycloak 26.4, Auth0, Okta |
| Passwordless authentication in 2026? | Passkeys (FIDO2/WebAuthn) — origin-bound credentials, phishing resistant, 1-tap biometric login |
| Enterprise SSO protocol? | OIDC for new integrations; SAML 2.0 when mandated by enterprise customer or compliance framework |
| CLI / IoT / smart TV auth? | Device Authorization Grant (RFC 8628) |
| MFA method ranking? | FIDO2/passkeys (phishing-resistant) > TOTP apps > push notifications > SMS OTP (weakest) |
| JWT algorithm choice? | RS256 or ES256 (asymmetric); never allow `none`; pin algorithm server-side |
| CORS with credentials? | `Access-Control-Allow-Credentials: true` with exact origin (never `*`); cookies need `SameSite=None; Secure` |
| CSP for auth pages? | `script-src 'nonce-{RANDOM}' 'strict-dynamic'; object-src 'none'; base-uri 'none'` |
| Session hijacking prevention? | Regenerate session ID after auth, bind session to IP/UA fingerprint, enforce HSTS, HttpOnly + Secure cookies |
| Zero trust authentication? | Continuous verification, least privilege, assume breach; start with strong IdP + MFA + device posture checks |

---

## Response Format

- **"Which approach should I use?"** — lead with the Quick Decision Guide row, then expand with rationale and trade-offs from the context file. Keep to 1–3 paragraphs unless depth is requested.
- **"How do I implement X?"** — provide a working code snippet with inline comments explaining the security-relevant choices, then list remaining security considerations.
- **Cross-domain questions** (e.g., auth + database storage) — answer the auth portion using this skill's patterns, then name the other skill that covers the remainder.
- Always name the specific RFC, OWASP section, or MDN page that backs each recommendation.

---

## Passkeys, WebAuthn, and FIDO2

### Registration ceremony

```javascript
const credential = await navigator.credentials.create({
  publicKey: {
    challenge: crypto.getRandomValues(new Uint8Array(32)),
    rp: { name: 'My App', id: 'example.com' },
    user: {
      id: Uint8Array.from(userId, c => c.charCodeAt(0)),
      name: 'user@example.com',
      displayName: 'Jane Doe'
    },
    pubKeyCredParams: [
      { alg: -7,   type: 'public-key' },   // ES256
      { alg: -257, type: 'public-key' }    // RS256
    ],
    authenticatorSelection: {
      authenticatorAttachment: 'platform',  // 'cross-platform' for security keys
      residentKey: 'required',              // discoverable credential
      userVerification: 'preferred'
    },
    timeout: 60000,
    attestation: 'none'
  }
});
// Server stores: credentialId, publicKey, signCount, transports
```

### Authentication ceremony

```javascript
const assertion = await navigator.credentials.get({
  publicKey: {
    challenge: serverChallenge,
    rpId: 'example.com',
    allowCredentials: [{ id: credentialIdFromServer, type: 'public-key', transports: ['internal', 'hybrid'] }],
    userVerification: 'preferred',
    timeout: 60000
  }
});
// Server verifies: signature over clientDataHash + authenticatorData
// Checks signCount > stored signCount (cloned authenticator detection)
```

### Conditional UI (passkey autofill)

```html
<input type="text" name="username" autocomplete="username webauthn">
```

```javascript
// mediation: 'conditional' shows passkeys in autofill dropdown without a modal
const credential = await navigator.credentials.get({
  publicKey: { challenge, rpId: 'example.com', userVerification: 'preferred' },
  mediation: 'conditional'
});
```

### PRF Extension

The PRF extension generates a 32-byte deterministic secret per credential, useful for deriving encryption keys. Browser support (2026): Chrome 132+, Safari 18+, Firefox 139+, Windows Hello 25H2+.

### Account recovery

- Register multiple authenticators per user (platform + roaming)
- Allow recovery via verified email with a time-limited link
- Offer a secondary TOTP as a fallback authenticator
- Never allow password-only recovery to bypass passkey enrollment

Source: [FIDO Alliance — Passkeys](https://fidoalliance.org/passkeys/), [W3C Web Authentication Level 3](https://www.w3.org/TR/webauthn-3/)

---

## OpenID Connect (OIDC)

### Key components

| Component | Purpose |
|---|---|
| ID Token | JWT with user identity claims (`sub`, `email`, `name`, `iss`, `aud`, `exp`, `nonce`) |
| UserInfo Endpoint | Optional endpoint returning additional user claims |
| Discovery Document | `/.well-known/openid-configuration` — auto-discovery of endpoints |
| JWKS Endpoint | Serves provider's public signing keys; clients fetch and cache |

### ID token validation checklist

1. Verify JWT signature using provider's JWKS keys
2. Confirm `iss` matches expected issuer URL exactly
3. Confirm `aud` contains your `client_id`
4. Check `exp` has not passed (allow 30–60s clock skew)
5. Validate `nonce` matches value sent in authorization request (prevents replay)
6. Check `iat` is not unreasonably old
7. For multi-tenant apps, verify `azp` (authorized party) if present

### OIDC logout (three layers)

1. **Application session**: Clear local session cookie/token
2. **IdP session**: Redirect to `end_session_endpoint`
3. **Refresh tokens**: Revoke via `/revoke` endpoint (RFC 7009)

### Discovery-based integration

```javascript
const config = await fetch('https://auth.example.com/.well-known/openid-configuration').then(r => r.json());
// config.authorization_endpoint, token_endpoint, jwks_uri, etc.
```

Source: [OpenID Connect Core 1.0](https://openid.net/specs/openid-connect-core-1_0.html)

---

## SAML 2.0 vs OIDC

| Aspect | SAML 2.0 | OIDC |
|---|---|---|
| Data format | XML assertions | JSON / JWT |
| Transport | Browser redirects, POST bindings | REST + JSON + OAuth 2.0 |
| Discovery | Manual metadata XML exchange | `/.well-known/openid-configuration` |
| Key rotation | Manual cert exchange with overlap window | Automatic via JWKS (`kid` rotation) |
| Mobile/SPA support | Poor (XML parsing, no native mobile flow) | Native (JSON, PKCE, device flow) |
| Library quality | Prone to XML Signature Wrapping CVEs | First-class, actively maintained |

**When to use each:**
- **OIDC**: Default for all new integrations
- **SAML 2.0**: When an enterprise customer provides a SAML metadata XML, or compliance mandates it
- **Both**: B2B SaaS selling up-market should ship both

**SAML security pitfalls:** XML Signature Wrapping attacks, XSLT injection, XML entity expansion (billion laughs) — disable external entity resolution. Manage certificate expiry with overlap windows.

---

## Multi-Factor Authentication

### MFA method security ranking

| Method | Phishing Resistant | Replay Resistant |
|---|---|---|
| FIDO2 / Passkeys / WebAuthn | Yes (origin-bound) | Yes |
| Hardware security keys (YubiKey) | Yes | Yes |
| TOTP authenticator apps | No (real-time phishing viable) | Partially (30s window) |
| Push notifications | No (MFA bombing) | Yes |
| SMS OTP | No (SIM swap, SS7) | Partially |

### TOTP implementation (RFC 6238)

```javascript
function verifyTOTP(secret, userCode, window = 1) {
  const time = Math.floor(Date.now() / 30000);
  for (let i = -window; i <= window; i++) {
    const expected = generateHOTP(secret, time + i);
    if (crypto.timingSafeEqual(Buffer.from(expected), Buffer.from(userCode))) return true;
  }
  return false;
}
```

**MFA blocks 99.9% of automated account attacks** (Microsoft 2025 security research). Even SMS OTP dramatically reduces account takeover risk.

---

## Zero Trust Authentication

### Core principles

1. **Verify explicitly**: Authenticate and authorize based on all available data points — user identity, device health, location, behavior anomalies
2. **Least privilege access**: Just-in-time and just-enough access; risk-based adaptive policies
3. **Assume breach**: End-to-end encryption, threat analytics, minimize blast radius

### Implementation phases

1. **Identity**: Deploy modern IdP, enforce MFA on all applications, implement SSO
2. **Network**: Micro-segmentation to isolate workloads, limit lateral movement
3. **Device**: Device posture checks, certificate-based device identity
4. **Continuous**: Risk-based adaptive auth, behavioral analytics, session binding

Source: [NIST SP 800-207 — Zero Trust Architecture](https://csrc.nist.gov/pubs/sp/800-207/final)

---

## Token Storage Security

| Storage | XSS Vulnerable | CSRF Vulnerable | Survives Reload | Best For |
|---|---|---|---|---|
| JavaScript memory | No | No | No | Short-lived access tokens |
| HttpOnly cookie | No | Yes (mitigate with SameSite + CSRF token) | Yes | Refresh tokens, session IDs |
| localStorage | Yes | No | Yes | Low-value preferences only |
| sessionStorage | Yes | No | No | Temporary non-sensitive data |
| `chrome.storage.session` | No | No | No | Extension access tokens |
| `chrome.storage.local` (encrypted) | No | No | Yes | Extension refresh tokens |

### Recommended architecture

```
Access Token:  JS memory (closure or React state)
               Short-lived (5–15 min). Lost on page refresh — acceptable because
               silent refresh via refresh token handles recovery.

Refresh Token: HttpOnly Secure SameSite=Strict cookie
               Long-lived (7–30 days). Server-side rotation with reuse detection.

Session ID:    HttpOnly Secure SameSite=Lax __Host- prefixed cookie
               Regenerated after authentication. Server-side session store.
```

**Never store tokens in localStorage for high-value applications** — any XSS vulnerability gives the attacker direct token access.

---

## Device Authorization Grant (RFC 8628)

For input-constrained devices without a browser: smart TVs, CLI tools, IoT, game consoles.

```
1. POST /device/code  →  { device_code, user_code: "WDJB-MJHT", verification_uri, interval: 5 }
2. Device shows user_code + verification_uri to user
3. User visits verification_uri on a browser-capable device, enters user_code, authenticates
4. Device polls at specified interval:
   POST /token  grant_type=urn:ietf:params:oauth:grant-type:device_code&device_code=...
   Responses: "authorization_pending" | "slow_down" | access_token
```

**Security requirements:** User codes must be short (≤8 chars, e.g., `WDJB-MJHT`). Device codes must be cryptographically random. Enforce the polling interval — respond with `slow_down` if client polls too fast. Set reasonable expiry (15–30 minutes).

Source: [RFC 8628](https://datatracker.ietf.org/doc/html/rfc8628)

---

## JWT Security Hardening

### Algorithm confusion attack

An attacker changes the JWT `alg` header from `RS256` to `HS256` and signs with the server's public key as the HMAC secret. If the server trusts the `alg` header, it verifies the forged token.

### Prevention checklist

1. Pin the algorithm server-side — never trust the JWT `alg` field
2. Reject `none` unconditionally — no production use case for unsigned JWTs
3. Use an allowlist: explicitly permit `RS256` or `ES256` only
4. Separate key objects for HMAC vs RSA — never pass an RSA public key to an HMAC function
5. Validate `iss` and `aud` — algorithm pinning alone doesn't prevent cross-realm reuse
6. Set short expiry — limits damage from a valid forged token
7. Use `jti` for revocation — maintain a blocklist of revoked token IDs

```javascript
// Node.js — jsonwebtoken library
const token = jwt.sign(payload, privateKey, {
  algorithm: 'RS256',   // pin algorithm
  expiresIn: '15m',
  issuer: 'https://auth.example.com',
  audience: 'https://api.example.com',
  jwtid: crypto.randomUUID()
});

const decoded = jwt.verify(token, publicKey, {
  algorithms: ['RS256'],  // ALLOWLIST — rejects HS256, none, etc.
  issuer: 'https://auth.example.com',
  audience: 'https://api.example.com',
  clockTolerance: 30
});
```

Source: [PortSwigger — Algorithm Confusion Attacks](https://portswigger.net/web-security/jwt/algorithm-confusion)

---

## CORS and Cross-Origin Authentication

### Server requirements for credentialed requests

```
Access-Control-Allow-Origin: https://app.example.com   (exact origin, never *)
Access-Control-Allow-Credentials: true
Access-Control-Allow-Methods: GET, POST, PUT, DELETE
Access-Control-Allow-Headers: Content-Type, Authorization
Access-Control-Max-Age: 86400
```

**The wildcard restriction:** `Access-Control-Allow-Origin: *` with `Access-Control-Allow-Credentials: true` is rejected by browsers. Use an exact origin or a validated dynamic origin.

### Cookie requirements for cross-origin

Cookies sent cross-origin must have `SameSite=None; Secure`.

### Common CORS auth pitfalls

| Pitfall | Fix |
|---|---|
| `origin: '*'` with `credentials: true` | Use explicit origin array or validator |
| Missing preflight handler | Add `app.options('*', cors(corsOptions))` |
| Cookie missing `SameSite=None` | Add `SameSite=None; Secure` to cookie |
| Reflecting `Origin` without validation | Validate against an allowlist |

Source: [MDN — CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS)

---

## CSP for Auth Pages

```http
Content-Security-Policy:
  script-src 'nonce-{RANDOM}' 'strict-dynamic';
  object-src 'none';
  base-uri 'none';
  frame-ancestors 'none';
  form-action 'self';
```

- Generate a new cryptographically random nonce per HTTP response (128+ bits)
- Inject nonce into every `<script>` tag: `<script nonce="{RANDOM}">`
- `strict-dynamic` propagates trust to scripts loaded by nonced scripts
- Never use `'unsafe-inline'` on auth pages

Source: [OWASP — CSP Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Content_Security_Policy_Cheat_Sheet.html)

---

## Session Hijacking Prevention

### Attack vectors and primary defenses

| Vector | Primary Defense |
|---|---|
| XSS token theft | HttpOnly cookies, CSP, output encoding |
| Session fixation | Regenerate session ID after authentication |
| Network sniffing | HTTPS everywhere, HSTS |
| Session sidejacking | HSTS preload, Secure cookie flag |

### Defense-in-depth checklist

1. Regenerate session ID after every authentication event or privilege escalation
2. Set all security cookie attributes: `Secure; HttpOnly; SameSite=Lax; __Host-`
3. Enforce HSTS: `Strict-Transport-Security: max-age=63072000; includeSubDomains; preload`
4. Bind session to client fingerprint: check IP range + User-Agent consistency
5. Implement absolute session timeout (banking: 30 min, general: 8 hours)
6. Monitor for anomalies: impossible travel, concurrent sessions from different geolocations
7. Provide session management UI: let users view and revoke active sessions

**Session hijacking bypasses MFA and passwords entirely** — growing 127% year-over-year in 2026.

---

## Identity Provider Integration

| IdP | Best For | Protocols |
|---|---|---|
| Auth0 (Okta) | B2C + B2B SaaS | OIDC, SAML, Social |
| Okta Workforce | Enterprise workforce | OIDC, SAML, SCIM |
| Microsoft Entra ID | Microsoft ecosystem | OIDC, SAML |
| Keycloak | Self-hosted, open source | OIDC, SAML, LDAP, DPoP (26.4+) |
| Google Identity | Google Workspace | OIDC |
| AWS Cognito | AWS-native apps | OIDC, SAML |

**Integration patterns:**
- **Direct**: App connects directly to one IdP
- **Identity brokering**: Keycloak or Auth0 acts as a federation hub
- **Multi-tenant federation**: Each customer tenant configures their own IdP

---

## OAuth 2.1 Status (2026)

OAuth 2.1 is at draft-ietf-oauth-v2-1-15 (March 2026). Not yet final RFC, but requirements are stable and widely adopted.

| Feature | OAuth 2.0 | OAuth 2.1 |
|---|---|---|
| PKCE | Optional | Mandatory for all Authorization Code clients |
| Implicit flow | Available | Removed |
| Resource Owner Password flow | Available | Removed |
| Redirect URI matching | Flexible | Exact string match required |
| Refresh token rotation | Recommended | Required for public clients |
| Bearer token in query strings | Allowed | Prohibited |

**RFC 9700** (OAuth 2.0 Security BCP, January 2025) mandates PKCE even for confidential server-side clients.

Source: [OAuth 2.1 Draft](https://oauth.net/2.1/), [RFC 9700](https://datatracker.ietf.org/doc/rfc9700/)

---

## Security Headers for Auth

```http
Strict-Transport-Security: max-age=63072000; includeSubDomains; preload
Content-Security-Policy: script-src 'nonce-{RANDOM}' 'strict-dynamic'; object-src 'none'; base-uri 'none'
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Referrer-Policy: strict-origin-when-cross-origin
Cache-Control: no-store
Pragma: no-cache
```

- `Cache-Control: no-store` on all auth endpoints — prevents caching of tokens or credentials
- `X-Frame-Options: DENY` on login pages — prevents clickjacking
- `Referrer-Policy: strict-origin-when-cross-origin` — prevents token leakage in referrer headers

---

## Cited Sources

1. [MDN — Using HTTP Cookies](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Cookies)
2. [MDN — CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS)
3. [MDN — CSP](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CSP)
4. [OWASP — CSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)
5. [OWASP — CSP Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Content_Security_Policy_Cheat_Sheet.html)
6. [RFC 9700 — OAuth 2.0 Security BCP](https://datatracker.ietf.org/doc/rfc9700/)
7. [RFC 9449 — OAuth 2.0 DPoP](https://datatracker.ietf.org/doc/html/rfc9449)
8. [RFC 8628 — OAuth 2.0 Device Authorization Grant](https://datatracker.ietf.org/doc/html/rfc8628)
9. [OAuth 2.1 Draft](https://oauth.net/2.1/)
10. [OpenID Connect Core 1.0](https://openid.net/specs/openid-connect-core-1_0.html)
11. [W3C — Web Authentication Level 3](https://www.w3.org/TR/webauthn-3/)
12. [FIDO Alliance — Passkeys](https://fidoalliance.org/passkeys/)
13. [NIST SP 800-207 — Zero Trust Architecture](https://csrc.nist.gov/pubs/sp/800-207/final)
14. [web.dev — Strict CSP](https://web.dev/articles/strict-csp)
15. [PortSwigger — Algorithm Confusion Attacks](https://portswigger.net/web-security/jwt/algorithm-confusion)
16. [PortSwigger — CORS and Access-Control-Allow-Origin](https://portswigger.net/web-security/cors/access-control-allow-origin)
