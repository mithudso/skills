<!-- hub-reference-banner -->
> **Reference file — part of the `security-review` hub.** Formerly the standalone `http-security-headers` skill.
> Sibling topics in this family are now reference files under the hubs (`security-review`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: http-security-headers
title: "HTTP Security Headers"
description: >-
  HTTP security headers expert. Covers Helmet.js middleware configuration, Content-Security-Policy
  (CSP directives, nonce/hash strategies, report-to, report-only mode), HSTS (preload ramp-up,
  preload list requirements), X-Frame-Options, X-Content-Type-Options, Referrer-Policy,
  Permissions-Policy, COOP/COEP/CORP, Chrome extension CSP (MV3 extension_pages vs sandbox),
  and testing/validation tooling (securityheaders.com, Mozilla Observatory, CSP Evaluator).
  TRIGGER: configuring, reviewing, or debugging HTTP security headers on any web server or
  Express/Node.js app; setting up Helmet.js; writing or tightening a Content-Security-Policy;
  configuring HSTS or its preload list; implementing COOP/COEP for SharedArrayBuffer or
  cross-origin isolation; setting Chrome extension CSP in manifest.json.
  SKIP: application-level auth/authz and API access control (use security-compliance-auditor);
  network-layer firewall rules and cloud security groups (use aws-core);
  full security audit of a codebase (use security-compliance-auditor).
version: "1.1.1"
updated: "2026-05-31"
category: developer
tags:
  - security
  - csp
  - hsts
  - helmet
  - http-headers
  - x-frame-options
  - chrome-extension
  - cors
  - permissions-policy
keywords:
  - Content-Security-Policy
  - CSP nonce
  - CSP hash
  - strict-dynamic
  - Helmet.js
  - HSTS preload
  - Strict-Transport-Security
  - X-Frame-Options
  - X-Content-Type-Options
  - Referrer-Policy
  - Permissions-Policy
  - Cross-Origin-Opener-Policy
  - Cross-Origin-Embedder-Policy
  - Cross-Origin-Resource-Policy
  - COOP
  - COEP
  - CORP
  - report-to
  - report-uri
  - unsafe-inline
  - frame-ancestors
  - Chrome extension CSP
  - MV3 CSP
whenToUse:
  - "Configuring Helmet.js on an Express or Node.js server"
  - "Writing or tightening a Content-Security-Policy"
  - "Setting up HSTS and evaluating preload list eligibility"
  - "Implementing COOP/COEP for SharedArrayBuffer or cross-origin isolation"
  - "Setting Chrome extension CSP in manifest.json for MV3"
  - "Debugging a CSP violation in the browser console"
  - "Reviewing HTTP response headers for security grade A or higher"
  - "Migrating from X-Frame-Options to CSP frame-ancestors"
  - "Setting up CSP violation reporting with report-to"
whenNotToUse:
  - "Application-level auth/authz and API access control — use security-compliance-auditor"
  - "Network-layer firewall rules and cloud security groups — use aws-core"
  - "Full codebase security audit — use security-compliance-auditor"
related_skills:
  - software-engineering-patterns
  - security-compliance-auditor
---

# HTTP Security Headers

Comprehensive reference for HTTP security headers: Helmet.js middleware, CSP, HSTS, X-Frame-Options, Referrer-Policy, Permissions-Policy, COOP/COEP, Chrome extension CSP, and validation tooling.

## 1. Quick Reference Table

| Header | Purpose | Recommended Value | Notes |
|---|---|---|---|
| `Content-Security-Policy` | XSS / injection mitigation | Strict nonce- or hash-based policy | Most important header; see section 3 |
| `Strict-Transport-Security` | Force HTTPS | `max-age=63072000; includeSubDomains; preload` | HTTPS-only; preload is permanent |
| `X-Content-Type-Options` | Prevent MIME sniffing | `nosniff` | Always set; no configuration needed |
| `X-Frame-Options` | Clickjacking protection | `DENY` or `SAMEORIGIN` | Legacy; prefer CSP `frame-ancestors` |
| `Referrer-Policy` | Control referrer leakage | `strict-origin-when-cross-origin` | Prevents leaking paths/query params |
| `Permissions-Policy` | Disable unused browser APIs | `camera=(), microphone=(), geolocation=()` | Disable everything you do not use |
| `Cross-Origin-Opener-Policy` | Isolate browsing context | `same-origin` | Prevents XS-Leaks |
| `Cross-Origin-Embedder-Policy` | Require CORS for embeds | `require-corp` | Required for `SharedArrayBuffer` |
| `Cross-Origin-Resource-Policy` | Restrict who loads resource | `same-origin` or `same-site` | Protects images/scripts/fonts |
| `X-Permitted-Cross-Domain-Policies` | Block Flash/PDF cross-domain | `none` | Legacy but still recommended |

## 2. Helmet.js (Express Middleware)

Helmet sets sensible security headers with one line. Current versions (v7+/v8) set 13+ headers by default.

### 2.1 Basic Setup

```javascript
import express from "express";
import helmet from "helmet";

const app = express();

// Apply all defaults (recommended starting point)
app.use(helmet());
```

### 2.2 Default Headers Set by helmet()

| Middleware function | Header set |
|---|---|
| `contentSecurityPolicy` | `Content-Security-Policy` |
| `crossOriginEmbedderPolicy` | `Cross-Origin-Embedder-Policy` |
| `crossOriginOpenerPolicy` | `Cross-Origin-Opener-Policy` |
| `crossOriginResourcePolicy` | `Cross-Origin-Resource-Policy` |
| `dnsPrefetchControl` | `X-DNS-Prefetch-Control` |
| `frameguard` | `X-Frame-Options` |
| `hidePoweredBy` | Removes `X-Powered-By` |
| `hsts` | `Strict-Transport-Security` |
| `ieNoOpen` | `X-Download-Options` |
| `noSniff` | `X-Content-Type-Options` |
| `originAgentCluster` | `Origin-Agent-Cluster` |
| `permittedCrossDomainPolicies` | `X-Permitted-Cross-Domain-Policies` |
| `referrerPolicy` | `Referrer-Policy` |
| `xssFilter` | `X-XSS-Protection` (set to `0`) |

### 2.3 Default CSP from Helmet

```
default-src 'self';
base-uri 'self';
font-src 'self' https: data:;
form-action 'self';
frame-ancestors 'self';
img-src 'self' data:;
object-src 'none';
script-src 'self';
script-src-attr 'none';
style-src 'self' https: 'unsafe-inline';
upgrade-insecure-requests
```

### 2.4 Custom Configuration

```javascript
app.use(
  helmet({
    // Disable a specific header
    crossOriginEmbedderPolicy: false,

    // Configure CSP with custom directives
    contentSecurityPolicy: {
      directives: {
        defaultSrc: ["'self'"],
        scriptSrc: ["'self'", "'nonce-abc123'"],
        styleSrc: ["'self'", "https://fonts.googleapis.com"],
        fontSrc: ["'self'", "https://fonts.gstatic.com"],
        imgSrc: ["'self'", "data:", "https://cdn.example.com"],
        connectSrc: ["'self'", "https://api.example.com"],
        frameSrc: ["'none'"],
        objectSrc: ["'none'"],
        upgradeInsecureRequests: [],
      },
    },

    // Configure HSTS
    hsts: {
      maxAge: 63072000, // 2 years
      includeSubDomains: true,
      preload: true,
    },

    // Configure Referrer-Policy
    referrerPolicy: {
      policy: "strict-origin-when-cross-origin",
    },
  })
);
```

### 2.5 Nonce-Based CSP with Helmet

```javascript
import crypto from "node:crypto";

app.use((req, res, next) => {
  // Generate a unique nonce per request
  res.locals.cspNonce = crypto.randomBytes(32).toString("base64");
  next();
});

app.use(
  helmet({
    contentSecurityPolicy: {
      directives: {
        defaultSrc: ["'self'"],
        scriptSrc: [
          "'self'",
          (req, res) => `'nonce-${res.locals.cspNonce}'`,
        ],
        styleSrc: [
          "'self'",
          (req, res) => `'nonce-${res.locals.cspNonce}'`,
        ],
      },
    },
  })
);

// In your templates, use the nonce on inline scripts/styles:
// <script nonce="<%= cspNonce %>">...</script>
```

### 2.6 Report-Only Mode

```javascript
app.use(
  helmet({
    contentSecurityPolicy: {
      directives: {
        ...helmet.contentSecurityPolicy.getDefaultDirectives(),
        "report-to": "csp-endpoint",
      },
      reportOnly: true, // CSP violations logged, not blocked
    },
  })
);
```

### 2.7 Development vs Production

```javascript
const isDev = process.env.NODE_ENV !== "production";

app.use(
  helmet({
    contentSecurityPolicy: isDev
      ? false // Disable CSP in dev to avoid localhost HTTPS issues
      : {
          directives: {
            ...helmet.contentSecurityPolicy.getDefaultDirectives(),
          },
        },
    hsts: isDev ? false : { maxAge: 63072000 },
  })
);
```

**Warning:** Helmet's default CSP includes `upgrade-insecure-requests`, which causes browsers to upgrade `http://` to `https://`. Safari upgrades `http://localhost` to `https://localhost`, which breaks local dev. Disable CSP or remove that directive in development.

## 3. Content-Security-Policy (CSP)

### 3.1 Directive Reference

| Directive | Controls | Example |
|---|---|---|
| `default-src` | Fallback for all fetch directives | `'self'` |
| `script-src` | JavaScript sources | `'self' 'nonce-abc'` |
| `script-src-elem` | `<script>` elements only | `'self'` |
| `script-src-attr` | Inline event handlers (onclick) | `'none'` |
| `style-src` | CSS sources | `'self' 'nonce-abc'` |
| `style-src-elem` | `<style>` and `<link rel=stylesheet>` | `'self'` |
| `style-src-attr` | Inline style attributes | `'unsafe-inline'` |
| `img-src` | Image sources | `'self' data: https:` |
| `font-src` | Font sources | `'self' https://fonts.gstatic.com` |
| `connect-src` | XHR, fetch, WebSocket, EventSource | `'self' https://api.example.com` |
| `media-src` | Audio and video sources | `'self'` |
| `object-src` | object, embed, applet elements | `'none'` |
| `frame-src` | iframe sources | `'none'` |
| `worker-src` | Worker, SharedWorker, ServiceWorker | `'self'` |
| `frame-ancestors` | Who can embed this page | `'self'` |
| `base-uri` | Restrict base element URLs | `'self'` |
| `form-action` | Form submission targets | `'self'` |
| `report-uri` | Legacy violation report endpoint | `/csp-report` (deprecated) |
| `report-to` | Modern violation report endpoint | `csp-endpoint` |
| `upgrade-insecure-requests` | Auto-upgrade HTTP to HTTPS | (no value) |

### 3.2 Source Values

| Value | Meaning |
|---|---|
| `'self'` | Same origin (scheme + host + port) |
| `'none'` | Block everything |
| `'unsafe-inline'` | Allow inline scripts/styles (AVOID) |
| `'unsafe-eval'` | Allow string-to-code execution APIs (AVOID) |
| `'wasm-unsafe-eval'` | Allow WebAssembly compilation |
| `'strict-dynamic'` | Trust scripts loaded by already-trusted scripts |
| `'nonce-<base64>'` | Allow elements with matching nonce attribute |
| `'sha256-<hash>'` | Allow elements with matching content hash |
| `https:` | Any HTTPS URL |
| `data:` | data: URIs |
| `blob:` | blob: URIs |
| `*.example.com` | Wildcard subdomain |

### 3.3 Nonce-Based CSP (Recommended for Dynamic Apps)

```
Content-Security-Policy:
  default-src 'self';
  script-src 'nonce-4AEemGb0xJptoIGFP3Nd' 'strict-dynamic';
  style-src 'nonce-4AEemGb0xJptoIGFP3Nd';
  object-src 'none';
  base-uri 'self';
```

**Requirements:**
- Generate a cryptographically random nonce per response (minimum 128 bits / 16 bytes)
- Attach the nonce to every inline script and style element
- `strict-dynamic` propagates trust to scripts loaded by nonced scripts
- When nonce/hash is present, browsers ignore `unsafe-inline` (backward-compatible)

### 3.4 Hash-Based CSP (Recommended for Static Pages)

```
Content-Security-Policy:
  script-src 'sha256-RFWPLDbv2BY+rCkDzsE+0fr8ylGr2R2faWMhq4lfEQc=';
```

Generate the hash:
```bash
echo -n 'console.log("hello")' | openssl dgst -sha256 -binary | openssl base64
```

### 3.5 CSP Reporting

**Modern approach (report-to + Reporting-Endpoints):**

```
Reporting-Endpoints: csp-endpoint="https://reports.example.com/csp"
Content-Security-Policy: default-src 'self'; report-to csp-endpoint
```

**Dual compatibility (supports both old and new browsers):**

```
Reporting-Endpoints: csp-endpoint="https://reports.example.com/csp"
Content-Security-Policy:
  default-src 'self';
  report-uri /csp-report;
  report-to csp-endpoint
```

Reports are POST-ed as `application/reports+json` (report-to) or `application/csp-report` (report-uri). Endpoints must use HTTPS.

### 3.6 Report-Only Mode

```
Content-Security-Policy-Report-Only:
  default-src 'self';
  script-src 'self';
  report-to csp-endpoint
```

Use report-only first: log violations, fix them, then switch to enforcing.

## 4. Strict-Transport-Security (HSTS)

### 4.1 Recommended Header

```
Strict-Transport-Security: max-age=63072000; includeSubDomains; preload
```

### 4.2 Directives

| Directive | Purpose | Recommended |
|---|---|---|
| `max-age` | Duration (seconds) browser remembers HTTPS-only | `63072000` (2 years) |
| `includeSubDomains` | Apply to all subdomains | Always include |
| `preload` | Request browser preload list inclusion | Only when fully committed |

### 4.3 Deployment Ramp-Up Strategy

Roll out gradually to catch issues before committing to preload:

```
# Stage 1: 5 minutes (test)
Strict-Transport-Security: max-age=300

# Stage 2: 1 week
Strict-Transport-Security: max-age=604800; includeSubDomains

# Stage 3: 1 month
Strict-Transport-Security: max-age=2592000; includeSubDomains

# Stage 4: 2 years (preload-eligible)
Strict-Transport-Security: max-age=63072000; includeSubDomains; preload
```

### 4.4 Preload List Requirements

To submit to https://hstspreload.org:

1. Valid HTTPS certificate
2. Redirect HTTP to HTTPS on the same host
3. Serve HSTS header over HTTPS with `max-age >= 31536000`
4. Include `includeSubDomains`
5. Include `preload`
6. **All subdomains** must support HTTPS

**Warning:** Preload is effectively permanent. Removal from the preload list takes months and requires a browser release cycle. Only preload when certain all current and future subdomains support HTTPS.

## 5. X-Frame-Options

| Value | Effect |
|---|---|
| `DENY` | Page cannot be framed at all |
| `SAMEORIGIN` | Page can only be framed by same origin |

`ALLOW-FROM uri` was never widely supported and is deprecated.

### Migration to CSP frame-ancestors

X-Frame-Options is legacy. Prefer `frame-ancestors` in CSP:

```
# Equivalent to X-Frame-Options: DENY
Content-Security-Policy: frame-ancestors 'none'

# Equivalent to X-Frame-Options: SAMEORIGIN
Content-Security-Policy: frame-ancestors 'self'

# Allow specific origins (not possible with X-Frame-Options)
Content-Security-Policy: frame-ancestors 'self' https://trusted.example.com
```

Set both headers for backward compatibility with older browsers.

## 6. Other Security Headers

### 6.1 X-Content-Type-Options

```
X-Content-Type-Options: nosniff
```

Prevents browsers from MIME-sniffing the Content-Type. Always set this header.

### 6.2 Referrer-Policy

| Policy | Behavior |
|---|---|
| `no-referrer` | Never send referrer |
| `no-referrer-when-downgrade` | Send for same-protocol; drop on HTTPS to HTTP |
| `origin` | Send origin only (no path) |
| `origin-when-cross-origin` | Full URL same-origin; origin only cross-origin |
| `same-origin` | Send only for same-origin requests |
| `strict-origin` | Origin only; drop on HTTPS to HTTP |
| `strict-origin-when-cross-origin` | **Recommended.** Full URL same-origin; origin cross-origin; drop on downgrade |
| `unsafe-url` | Always send full URL (AVOID) |

### 6.3 Permissions-Policy

Disable browser features you do not use. Syntax: `feature=(allowlist)`.

```
Permissions-Policy: camera=(), microphone=(), geolocation=(), payment=(),
  usb=(), bluetooth=(), accelerometer=(), gyroscope=(),
  magnetometer=(), midi=(), screen-wake-lock=(),
  display-capture=(), document-domain=()
```

Allow specific origins:

```
Permissions-Policy: geolocation=(self "https://maps.example.com"), camera=()
```

### 6.4 Cross-Origin-Opener-Policy (COOP)

| Value | Effect |
|---|---|
| `unsafe-none` | Default; no isolation |
| `same-origin` | Isolate browsing context from cross-origin documents |
| `same-origin-allow-popups` | Isolate but allow popups to retain opener reference |

```
Cross-Origin-Opener-Policy: same-origin
```

Prevents cross-origin documents from accessing your window object. Mitigates XS-Leaks and Spectre-class attacks.

### 6.5 Cross-Origin-Embedder-Policy (COEP)

| Value | Effect |
|---|---|
| `unsafe-none` | Default; no restriction |
| `require-corp` | All cross-origin resources must opt in via CORS or CORP |
| `credentialless` | Cross-origin no-cors requests sent without credentials |

```
Cross-Origin-Embedder-Policy: require-corp
```

Combined with COOP `same-origin`, enables `crossOriginIsolated` context (required for SharedArrayBuffer, high-resolution timers).

### 6.6 Cross-Origin-Resource-Policy (CORP)

```
Cross-Origin-Resource-Policy: same-origin
```

Set on static assets (images, scripts, fonts) to control who can load them. Values: `same-origin`, `same-site`, `cross-origin`.

## 7. Chrome Extension CSP (Manifest V3)

### 7.1 Manifest Configuration

```json
{
  "manifest_version": 3,
  "content_security_policy": {
    "extension_pages": "script-src 'self' 'wasm-unsafe-eval'; object-src 'self';",
    "sandbox": "sandbox allow-scripts allow-forms allow-popups; script-src 'self' 'unsafe-inline'; child-src 'self';"
  }
}
```

### 7.2 Extension Pages vs Sandbox

| Scope | Applies to | Can relax? |
|---|---|---|
| `extension_pages` | Popup, background worker, extension tabs, iframes opened by extension | No — minimum enforced: `script-src 'self' 'wasm-unsafe-eval'; object-src 'self'` |
| `sandbox` | Pages listed in `manifest.sandbox.pages` | Yes — can use unsafe-inline and dynamic code execution |

### 7.3 MV3 Restrictions

- `extension_pages` **cannot** use unsafe-inline or dynamic code execution keywords for script-src
- `extension_pages` **cannot** allowlist remote script hosts
- `wasm-unsafe-eval` is explicitly allowed for WebAssembly
- Sandboxed pages have no access to extension APIs but can run inline scripts
- Default if unspecified: `script-src 'self'; object-src 'self';`

### 7.4 Common Extension CSP Patterns

**Minimal (default):**
```json
"content_security_policy": {
  "extension_pages": "script-src 'self'; object-src 'self';"
}
```

**With WebAssembly:**
```json
"content_security_policy": {
  "extension_pages": "script-src 'self' 'wasm-unsafe-eval'; object-src 'self';"
}
```

**With sandbox for template engines:**
```json
"content_security_policy": {
  "sandbox": "sandbox allow-scripts; script-src 'self' 'unsafe-inline'; object-src 'self';"
},
"sandbox": {
  "pages": ["sandbox.html"]
}
```

### 7.5 Content Scripts CSP

Content scripts injected into web pages run under the **host page's CSP**, not the extension's CSP. If the host page blocks inline scripts, your content script cannot inject inline script elements. Use `chrome.scripting.executeScript()` from the service worker instead.

## 8. Common Pitfalls

### 8.1 CSP Anti-Patterns

| Pitfall | Problem | Fix |
|---|---|---|
| `script-src 'unsafe-inline'` | Allows any inline script, including XSS payloads | Use nonces or hashes |
| `script-src 'unsafe-eval'` | Allows string-to-code execution APIs | Refactor to avoid dynamic code generation |
| `default-src *` or `script-src *` | Permits loading from any origin | Use explicit allowlists |
| `script-src https:` | Any HTTPS host can serve scripts | Allowlist specific domains |
| `script-src 'self' data:` | data: URIs can contain executable scripts | Remove data: from script-src |
| Missing `object-src 'none'` | Flash/plugin-based XSS bypasses | Always set object-src to none |
| Missing `base-uri` | Attacker can inject base element to redirect relative URLs | Set base-uri to self |

### 8.2 HSTS Pitfalls

- Setting HSTS over HTTP (header is ignored; must be sent over HTTPS)
- Enabling preload before all subdomains support HTTPS
- Short max-age values providing inadequate protection
- Forgetting that HSTS preload removal takes months

### 8.3 Helmet Pitfalls

- `upgrade-insecure-requests` breaks `http://localhost` in Safari during development
- Helmet's default CSP allows style-src unsafe-inline; tighten for production
- COEP require-corp can break third-party images/scripts that lack CORS headers
- Calling helmet() after route handlers (must be early middleware)

### 8.4 General Header Pitfalls

- Setting security headers only on HTML responses; APIs and static assets need them too
- Misconfiguring CORS alongside COOP/COEP causing broken cross-origin requests
- Relying solely on X-Frame-Options when CSP frame-ancestors is more flexible
- Not testing in report-only mode before enforcing CSP

## 9. Testing and Validation

### 9.1 Online Scanners

| Tool | URL | What it checks |
|---|---|---|
| Security Headers | https://securityheaders.com | HTTP response header grading (A+ to F) |
| Mozilla Observatory | https://developer.mozilla.org/en-US/observatory | Headers, cookies, TLS, CORS, redirects |
| CSP Evaluator | https://csp-evaluator.withgoogle.com | CSP policy strength; finds bypasses |
| HSTS Preload Check | https://hstspreload.org | Preload eligibility verification |

### 9.2 CLI Validation

```bash
# Check response headers
curl -sI https://example.com | grep -iE '(content-security|strict-transport|x-frame|x-content-type|referrer-policy|permissions-policy|cross-origin)'

# Check CSP only
curl -sI https://example.com | grep -i content-security-policy

# Validate with Mozilla Observatory CLI (httpobs)
pip install httpobs-cli
httpobs example.com
```

### 9.3 CI/CD Integration

```yaml
# GitHub Actions: scan security headers after deploy
- name: Check security headers
  run: |
    GRADE=$(curl -s "https://securityheaders.com/?q=https://example.com&followRedirects=on" \
      | grep -oP 'class="grade_[A-F]"' | head -1)
    echo "Security Headers Grade: $GRADE"
```

### 9.4 Automated Testing (Node.js)

```javascript
import { describe, it, expect } from "vitest";

describe("security headers", () => {
  it("should return required security headers", async () => {
    const res = await fetch("http://localhost:3000");

    expect(res.headers.get("content-security-policy")).toBeTruthy();
    expect(res.headers.get("strict-transport-security")).toMatch(/max-age=\d+/);
    expect(res.headers.get("x-content-type-options")).toBe("nosniff");
    expect(res.headers.get("x-frame-options")).toMatch(/DENY|SAMEORIGIN/);
    expect(res.headers.get("referrer-policy")).toBeTruthy();
  });

  it("should not contain unsafe CSP directives", async () => {
    const res = await fetch("http://localhost:3000");
    const csp = res.headers.get("content-security-policy");

    expect(csp).not.toContain("unsafe-inline");
    expect(csp).not.toContain("unsafe-eval");
    expect(csp).not.toContain("script-src *");
  });
});
```

## 10. Production-Ready CSP Templates

### 10.1 Strict Nonce-Based (SPA with API)

```
Content-Security-Policy:
  default-src 'none';
  script-src 'nonce-{RANDOM}' 'strict-dynamic';
  style-src 'nonce-{RANDOM}';
  img-src 'self' data: https://cdn.example.com;
  font-src 'self' https://fonts.gstatic.com;
  connect-src 'self' https://api.example.com wss://realtime.example.com;
  frame-ancestors 'none';
  base-uri 'self';
  form-action 'self';
  object-src 'none';
  upgrade-insecure-requests;
  report-to csp-endpoint
```

### 10.2 Static Site (Hash-Based)

```
Content-Security-Policy:
  default-src 'none';
  script-src 'sha256-{HASH_OF_INLINE_SCRIPT}';
  style-src 'self';
  img-src 'self';
  font-src 'self';
  frame-ancestors 'none';
  base-uri 'self';
  form-action 'none';
  object-src 'none';
  upgrade-insecure-requests
```

### 10.3 Full Helmet Production Config

```javascript
import express from "express";
import helmet from "helmet";
import crypto from "node:crypto";

const app = express();

app.use((req, res, next) => {
  res.locals.cspNonce = crypto.randomBytes(32).toString("base64");
  next();
});

app.use(
  helmet({
    contentSecurityPolicy: {
      directives: {
        defaultSrc: ["'none'"],
        scriptSrc: [
          "'self'",
          "'strict-dynamic'",
          (req, res) => `'nonce-${res.locals.cspNonce}'`,
        ],
        styleSrc: [
          "'self'",
          (req, res) => `'nonce-${res.locals.cspNonce}'`,
        ],
        imgSrc: ["'self'", "data:", "https://cdn.example.com"],
        fontSrc: ["'self'", "https://fonts.gstatic.com"],
        connectSrc: ["'self'", "https://api.example.com"],
        frameAncestors: ["'none'"],
        baseUri: ["'self'"],
        formAction: ["'self'"],
        objectSrc: ["'none'"],
        upgradeInsecureRequests: [],
        reportTo: "csp-endpoint",
      },
    },
    hsts: {
      maxAge: 63072000,
      includeSubDomains: true,
      preload: true,
    },
    referrerPolicy: {
      policy: "strict-origin-when-cross-origin",
    },
    crossOriginOpenerPolicy: { policy: "same-origin" },
    crossOriginEmbedderPolicy: { policy: "require-corp" },
    crossOriginResourcePolicy: { policy: "same-origin" },
  })
);

// Set Reporting-Endpoints header (Helmet does not set this)
app.use((req, res, next) => {
  res.setHeader(
    "Reporting-Endpoints",
    'csp-endpoint="https://reports.example.com/csp"'
  );
  next();
});

// Set Permissions-Policy (Helmet does not set this)
app.use((req, res, next) => {
  res.setHeader(
    "Permissions-Policy",
    "camera=(), microphone=(), geolocation=(), payment=(), usb=()"
  );
  next();
});
```

## 11. Decision Checklist

When reviewing or setting up security headers, verify:

- [ ] CSP is present and uses nonce or hash (not unsafe-inline / unsafe-eval)
- [ ] object-src is set to none and base-uri is set to self in CSP
- [ ] CSP is tested in report-only mode before enforcing
- [ ] HSTS is set with max-age >= 63072000 over HTTPS only
- [ ] HSTS preload is only added when all subdomains support HTTPS
- [ ] X-Content-Type-Options: nosniff is present
- [ ] X-Frame-Options or CSP frame-ancestors is set
- [ ] Referrer-Policy is set to strict-origin-when-cross-origin or stricter
- [ ] Permissions-Policy disables all unused browser features
- [ ] COOP/COEP are set if cross-origin isolation is needed
- [ ] X-Powered-By is removed (Helmet does this automatically)
- [ ] Headers are validated with securityheaders.com or Mozilla Observatory
- [ ] Chrome extension CSP does not attempt to relax extension_pages beyond MV3 minimum
- [ ] Content scripts handle host page CSP restrictions gracefully
