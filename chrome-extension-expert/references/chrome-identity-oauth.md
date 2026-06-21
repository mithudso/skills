<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly the standalone `chrome-identity-oauth` skill.
> Sibling topics in this family are now reference files under the hubs (`chrome-extension-expert`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: chrome-identity-oauth
description: >
  Chrome extension OAuth patterns — chrome.identity API (getAuthToken,
  launchWebAuthFlow), PKCE authorization code flow, token refresh and caching,
  Google and non-Google providers (GitHub, Okta, Azure AD, Auth0), MV3 service
  worker token management, and security best practices.
  TRIGGER: implementing or reviewing authentication or OAuth in a Chrome extension,
  token refresh logic, chrome.identity calls, extension redirect URI setup, PKCE
  flow for extensions, multi-provider OAuth in MV3.
  SKIP: server-side OAuth with no extension context, non-browser OAuth (mobile or
  desktop native), general web auth without a Chrome extension component.
  See also: chrome-mv3-advanced (offscreen docs, native messaging depth),
  web-auth-patterns (general OAuth 2.1/PKCE/JWT outside extensions),
  chrome-extension-security-reviewer (permissions audit, CSP, storage security).
version: 1.1.0
category: developer
tags: [chrome-extension, oauth, pkce, identity, authentication, google, mv3, github, okta, azure]
related_skills: [chrome-mv3-advanced, chrome-storage-patterns, web-auth-patterns, chrome-extension-security-reviewer]
updated: 2026-05-29
---

# Chrome Identity & OAuth for Extensions

## API Surface Overview

| Method | Best for | Manifest requirement | Provider support |
|---|---|---|---|
| `getAuthToken` | Google accounts only | `oauth2` block in manifest.json | Google only |
| `launchWebAuthFlow` | Any OAuth/OIDC provider | `identity` permission | Google, GitHub, Okta, Azure AD, Auth0, any OIDC |
| `getProfileUserInfo` | Signed-in Chrome user info | `identity.email` permission | Google profile only |
| `removeCachedAuthToken` | Invalidate cached Google token | `identity` permission | Google only |

**Manifest permissions:**

```json
{
  "permissions": ["identity"],
  "oauth2": {
    "client_id": "YOUR_CLIENT_ID.apps.googleusercontent.com",
    "scopes": ["openid", "email", "profile"]
  }
}
```

The `oauth2` block is only used by `getAuthToken`. For `launchWebAuthFlow`, omit it and manage client IDs in code.

## Decision Guide: getAuthToken vs launchWebAuthFlow

```
Need to auth with Google only?
  YES → Is cross-browser support needed (Brave, Edge, etc.)?
          YES → launchWebAuthFlow + PKCE
          NO  → getAuthToken (simplest path, zero token management code)
  NO  → launchWebAuthFlow + PKCE (only option for non-Google)
```

**Use `getAuthToken` when:** Google-only, Chrome-only, zero token management code is acceptable, and the `oauth2` manifest section is acceptable.

**Use `launchWebAuthFlow` when:** non-Google provider, cross-browser compatibility, need refresh tokens (getAuthToken manages them internally and does not expose them), or need full control over the PKCE flow.

## getAuthToken — Google-Only Fast Path

`getAuthToken` handles the full OAuth flow internally for Google accounts. It manages token caching, refresh, and consent prompts automatically.

```js
// --- service-worker.js (MV3) ---

async function getGoogleToken(interactive = true) {
  // Returns { token, grantedScopes } (Chrome 111+)
  const result = await chrome.identity.getAuthToken({ interactive });
  return result.token;
}

async function fetchGoogleProfile() {
  const token = await getGoogleToken(true);
  const resp = await fetch('https://www.googleapis.com/oauth2/v3/userinfo', {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (resp.status === 401) {
    // Token expired or revoked — clear cache and retry once
    await chrome.identity.removeCachedAuthToken({ token });
    const freshToken = await getGoogleToken(true);
    return fetch('https://www.googleapis.com/oauth2/v3/userinfo', {
      headers: { Authorization: `Bearer ${freshToken}` },
    }).then((r) => r.json());
  }
  return resp.json();
}
```

**Key behaviors:**
- `getAuthToken({ interactive: false })` is cheap and returns the cached token when valid
- Token refresh is automatic — when the cached token expires, the next call silently fetches a new one
- `getAuthToken` does **not** work in Brave or all Chromium forks — Chrome proper only
- The `oauth2.client_id` must be a Chrome App type credential in Google Cloud Console

**Incremental scopes:**
```js
chrome.identity.getAuthToken({
  interactive: true,
  scopes: ['https://www.googleapis.com/auth/calendar.readonly'],
}, (token) => { /* now has calendar scope */ });
```

## launchWebAuthFlow — Universal OAuth

`launchWebAuthFlow` opens a browser popup to any OAuth provider. When the provider redirects back to the extension's redirect URI, the popup closes and the final URL is returned.

### Redirect URI

```js
const redirectUri = chrome.identity.getRedirectURL();
// => "https://<extension-id>.chromiumapp.org/"
const redirectUriWithPath = chrome.identity.getRedirectURL('oauth2');
// => "https://<extension-id>.chromiumapp.org/oauth2"
```

Register this URI with your OAuth provider. For Google, use a **Web Application** credential (not Chrome App).

**Dev vs prod extension IDs:** An unpacked extension gets a different ID each load unless pinned with the `"key"` field in manifest.json. Register both redirect URIs during development, or pin the key:
```json
{ "key": "MIIBIjANBgkqh...<your-public-key>..." }
```

### Interactive vs silent mode

```js
// Interactive: shows login/consent UI — use for initial sign-in
chrome.identity.launchWebAuthFlow({ url, interactive: true });

// Silent: fails immediately if user interaction would be needed
// Use for background token refresh when cookies/session exist
chrome.identity.launchWebAuthFlow({ url, interactive: false });
```

When `interactive: false` and the provider uses JS redirects, add timeout options:
```js
chrome.identity.launchWebAuthFlow({
  url: authUrl.toString(),
  interactive: false,
  abortOnLoadForNonInteractive: false,
  timeoutMsForNonInteractive: 5000,
});
```

## PKCE Authorization Code Flow (Recommended)

PKCE (Proof Key for Code Exchange, RFC 7636) is the recommended flow for Chrome extensions — extensions cannot securely store a client secret.

### PKCE Helpers

```js
// --- pkce.js ---

function generateCodeVerifier(length = 64) {
  const array = new Uint8Array(length);
  crypto.getRandomValues(array);
  return base64UrlEncode(array);
}

async function generateCodeChallenge(verifier) {
  const data = new TextEncoder().encode(verifier);
  const digest = await crypto.subtle.digest('SHA-256', data);
  return base64UrlEncode(new Uint8Array(digest));
}

function base64UrlEncode(buffer) {
  // Use Array.from() to avoid call-stack limits on large typed arrays
  return btoa(Array.from(buffer, b => String.fromCharCode(b)).join(''))
    .replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

function generateState(length = 32) {
  const array = new Uint8Array(length);
  crypto.getRandomValues(array);
  return base64UrlEncode(array);
}
```

### Full PKCE Flow

```js
// --- auth.js (runs in service worker) ---

async function loginWithPKCE(config) {
  // config: { clientId, authEndpoint, tokenEndpoint, scopes }
  const codeVerifier = generateCodeVerifier();
  const codeChallenge = await generateCodeChallenge(codeVerifier);
  const state = generateState();
  const redirectUri = chrome.identity.getRedirectURL('oauth2');

  // Step 1: Build authorization URL
  const authUrl = new URL(config.authEndpoint);
  authUrl.searchParams.set('client_id', config.clientId);
  authUrl.searchParams.set('redirect_uri', redirectUri);
  authUrl.searchParams.set('response_type', 'code');
  authUrl.searchParams.set('scope', config.scopes.join(' '));
  authUrl.searchParams.set('code_challenge', codeChallenge);
  authUrl.searchParams.set('code_challenge_method', 'S256');
  authUrl.searchParams.set('state', state);

  // Step 2: Launch the auth popup
  const resultUrl = await chrome.identity.launchWebAuthFlow({
    url: authUrl.toString(),
    interactive: true,
  });

  // Step 3: Parse and verify the authorization code
  const responseUrl = new URL(resultUrl);
  if (responseUrl.searchParams.get('state') !== state) {
    throw new Error('OAuth state mismatch — possible CSRF attack');
  }
  const code = responseUrl.searchParams.get('code');
  if (!code) {
    throw new Error(`Authorization failed: ${responseUrl.searchParams.get('error') || 'no code returned'}`);
  }

  // Step 4: Exchange the code for tokens
  const tokenResponse = await fetch(config.tokenEndpoint, {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: new URLSearchParams({
      client_id: config.clientId,
      code,
      code_verifier: codeVerifier,
      grant_type: 'authorization_code',
      redirect_uri: redirectUri,
    }),
  });
  if (!tokenResponse.ok) {
    const err = await tokenResponse.json();
    throw new Error(`Token exchange failed: ${err.error_description || err.error}`);
  }

  const tokens = await tokenResponse.json();
  await saveTokens(tokens);
  return tokens;
}
```

## Token Storage and Refresh

### Storage Strategy

| Data | Store | Why |
|---|---|---|
| `refresh_token` | `chrome.storage.local` | Survives browser restart and SW termination |
| `access_token` + `expires_at` | `chrome.storage.session` | Cleared on browser close; fast access from SW |
| PKCE `code_verifier` | Local variable only | Never persisted; single-use within one auth flow |

**Content script access to session storage:** By default, `chrome.storage.session` is only accessible from the service worker. To allow content scripts to read session data:
```js
// In service-worker.js, at top level or in the 'install' handler:
chrome.storage.session.setAccessLevel({
  accessLevel: 'TRUSTED_AND_UNTRUSTED_CONTEXTS',
});
```

### Token Store Implementation

```js
// --- token-store.js ---

async function saveTokens(tokens) {
  const expiresAt = Date.now() + tokens.expires_in * 1000;
  if (tokens.refresh_token) {
    await chrome.storage.local.set({ oauth_refresh_token: tokens.refresh_token });
  }
  await chrome.storage.session.set({
    oauth_access_token: tokens.access_token,
    oauth_expires_at: expiresAt,
  });
}

async function getAccessToken() {
  const session = await chrome.storage.session.get([
    'oauth_access_token', 'oauth_expires_at',
  ]);
  // Return cached token if still valid (with 60s buffer)
  if (session.oauth_access_token && session.oauth_expires_at > Date.now() + 60_000) {
    return session.oauth_access_token;
  }
  return refreshAccessToken();
}

async function refreshAccessToken(tokenEndpoint, clientId) {
  const { oauth_refresh_token } = await chrome.storage.local.get('oauth_refresh_token');
  if (!oauth_refresh_token) {
    throw new Error('No refresh token — user must re-authenticate');
  }
  const resp = await fetch(tokenEndpoint, {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: new URLSearchParams({
      client_id: clientId,
      refresh_token: oauth_refresh_token,
      grant_type: 'refresh_token',
    }),
  });
  if (!resp.ok) {
    const err = await resp.json();
    if (err.error === 'invalid_grant') {
      await chrome.storage.local.remove('oauth_refresh_token');
      await chrome.storage.session.remove(['oauth_access_token', 'oauth_expires_at']);
      throw new Error('Refresh token revoked — user must re-authenticate');
    }
    throw new Error(`Token refresh failed: ${err.error_description || err.error}`);
  }
  const tokens = await resp.json();
  await saveTokens(tokens);
  return tokens.access_token;
}

async function logout() {
  await chrome.storage.local.remove(['oauth_refresh_token']);
  await chrome.storage.session.remove(['oauth_access_token', 'oauth_expires_at']);
  // If using getAuthToken, also clear Chrome's internal token cache.
  // MV3 promise form returns { token, grantedScopes } — extract the string.
  try {
    const result = await chrome.identity.getAuthToken({ interactive: false });
    if (result?.token) {
      await chrome.identity.removeCachedAuthToken({ token: result.token });
    }
  } catch {
    // Not signed in or token already cleared — nothing to do
  }
}
```

### Proactive Refresh with chrome.alarms

Service workers terminate after ~30 seconds of inactivity. Use `chrome.alarms` for proactive renewal:

```js
const REFRESH_ALARM = 'oauth-token-refresh';

async function scheduleRefreshAlarm() {
  const { oauth_expires_at } = await chrome.storage.session.get('oauth_expires_at');
  if (!oauth_expires_at) return;
  const refreshAt = oauth_expires_at - 5 * 60 * 1000; // 5 min before expiry
  const delayMinutes = Math.max(1, (refreshAt - Date.now()) / 60_000);
  chrome.alarms.create(REFRESH_ALARM, { delayInMinutes: delayMinutes });
}

chrome.alarms.onAlarm.addListener(async (alarm) => {
  if (alarm.name === REFRESH_ALARM) {
    try {
      await refreshAccessToken(TOKEN_ENDPOINT, CLIENT_ID);
      await scheduleRefreshAlarm();
    } catch (err) {
      console.error('[oauth] Proactive refresh failed:', err.message);
    }
  }
});
```

### Concurrency-Safe Refresh

When multiple contexts request a token simultaneously, coalesce refresh calls:

```js
let refreshPromise = null;

async function getAccessTokenSafe() {
  const session = await chrome.storage.session.get([
    'oauth_access_token', 'oauth_expires_at',
  ]);
  if (session.oauth_access_token && session.oauth_expires_at > Date.now() + 60_000) {
    return session.oauth_access_token;
  }
  if (!refreshPromise) {
    refreshPromise = refreshAccessToken(TOKEN_ENDPOINT, CLIENT_ID).finally(() => {
      refreshPromise = null;
    });
  }
  return refreshPromise;
}
```

## Non-Google Providers

### GitHub

GitHub does not support PKCE natively. Token exchange requires a client secret server-side.

```js
async function loginWithGitHub() {
  const state = generateState();
  const redirectUri = chrome.identity.getRedirectURL('github');
  const authUrl = new URL('https://github.com/login/oauth/authorize');
  authUrl.searchParams.set('client_id', GITHUB_CLIENT_ID);
  authUrl.searchParams.set('redirect_uri', redirectUri);
  authUrl.searchParams.set('scope', 'read:user user:email');
  authUrl.searchParams.set('state', state);

  const resultUrl = await chrome.identity.launchWebAuthFlow({
    url: authUrl.toString(), interactive: true,
  });

  const url = new URL(resultUrl);
  if (url.searchParams.get('state') !== state) throw new Error('State mismatch');
  const code = url.searchParams.get('code');

  // Exchange via your backend — GitHub requires client_secret server-side
  const tokens = await fetch('https://your-backend.example.com/api/github/token', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ code }),
  }).then((r) => r.json());

  await saveTokens(tokens);
  return tokens;
}
```

### Azure AD / Entra ID

Azure supports public clients — exchange directly without a client secret.

```js
async function loginWithAzure() {
  const codeVerifier = generateCodeVerifier();
  const codeChallenge = await generateCodeChallenge(codeVerifier);
  const state = generateState();
  const redirectUri = chrome.identity.getRedirectURL('azure');

  const authUrl = new URL(`https://login.microsoftonline.com/${AZURE_TENANT_ID}/oauth2/v2.0/authorize`);
  authUrl.searchParams.set('client_id', AZURE_CLIENT_ID);
  authUrl.searchParams.set('redirect_uri', redirectUri);
  authUrl.searchParams.set('response_type', 'code');
  authUrl.searchParams.set('scope', 'openid profile email offline_access');
  authUrl.searchParams.set('code_challenge', codeChallenge);
  authUrl.searchParams.set('code_challenge_method', 'S256');
  authUrl.searchParams.set('state', state);
  authUrl.searchParams.set('nonce', generateState()); // OIDC replay protection

  const resultUrl = await chrome.identity.launchWebAuthFlow({
    url: authUrl.toString(), interactive: true,
  });

  const url = new URL(resultUrl);
  if (url.searchParams.get('state') !== state) throw new Error('State mismatch');

  const resp = await fetch(
    `https://login.microsoftonline.com/${AZURE_TENANT_ID}/oauth2/v2.0/token`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: new URLSearchParams({
        client_id: AZURE_CLIENT_ID,
        code: url.searchParams.get('code'),
        code_verifier: codeVerifier,
        grant_type: 'authorization_code',
        redirect_uri: redirectUri,
      }),
    }
  );
  const tokens = await resp.json();
  await saveTokens(tokens);
  return tokens;
}
```

### Okta

Okta supports PKCE natively — exchange directly without a client secret.

```js
async function loginWithOkta() {
  const codeVerifier = generateCodeVerifier();
  const codeChallenge = await generateCodeChallenge(codeVerifier);
  const state = generateState();
  const redirectUri = chrome.identity.getRedirectURL('okta');

  const authUrl = new URL(`${OKTA_DOMAIN}/oauth2/default/v1/authorize`);
  authUrl.searchParams.set('client_id', OKTA_CLIENT_ID);
  authUrl.searchParams.set('redirect_uri', redirectUri);
  authUrl.searchParams.set('response_type', 'code');
  authUrl.searchParams.set('scope', 'openid profile email offline_access');
  authUrl.searchParams.set('code_challenge', codeChallenge);
  authUrl.searchParams.set('code_challenge_method', 'S256');
  authUrl.searchParams.set('state', state);
  authUrl.searchParams.set('nonce', generateState()); // OIDC replay protection

  const resultUrl = await chrome.identity.launchWebAuthFlow({
    url: authUrl.toString(), interactive: true,
  });

  const url = new URL(resultUrl);
  if (url.searchParams.get('state') !== state) throw new Error('State mismatch');

  const resp = await fetch(`${OKTA_DOMAIN}/oauth2/default/v1/token`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: new URLSearchParams({
      client_id: OKTA_CLIENT_ID,
      code: url.searchParams.get('code'),
      code_verifier: codeVerifier,
      grant_type: 'authorization_code',
      redirect_uri: redirectUri,
    }),
  });
  const tokens = await resp.json();
  await saveTokens(tokens);
  return tokens;
}
```

## Security Best Practices

| Practice | Rule |
|---|---|
| **PKCE over implicit** | Always use authorization code + S256. Implicit flow exposes tokens in URL fragments and is deprecated by OAuth 2.1. |
| **No client_secret in extension** | Extensions are public clients — never embed a secret. |
| **State parameter** | Always generate a random `state`, include it in the auth URL, and verify it in the response. |
| **Store `expires_at` not `expires_in`** | Absolute timestamps survive SW restarts; relative offsets do not. |
| **60-300s expiry buffer** | Check `expires_at > Date.now() + buffer` to avoid races. |
| **Redirect URI from API** | Always use `chrome.identity.getRedirectURL()` — never hardcode the URI. |
| **Scope minimization** | Request minimum scopes for current operation; add more via incremental authorization. |
| **Encrypt sensitive tokens** | For high-security needs, wrap the refresh token with AES-GCM before storing in `chrome.storage.local`. |

## Troubleshooting

| Symptom | Cause | Fix |
|---|---|---|
| `redirect_uri_mismatch` | URI not registered with OAuth provider | Add `chrome.identity.getRedirectURL()` output to your OAuth client config |
| `getAuthToken` returns `OAuth2 not granted or revoked` | User revoked consent | Call `removeCachedAuthToken`, then retry with `interactive: true` |
| `launchWebAuthFlow` hangs on non-Google provider | Provider uses JS redirect after page load | Set `abortOnLoadForNonInteractive: false` and `timeoutMsForNonInteractive: 5000` |
| Token gone after browser restart | Stored in `chrome.storage.session` | Move refresh_token to `chrome.storage.local` |
| SW dies mid-refresh | fetch() exceeds 30s or no activity | Keep token exchanges under 30s; use offscreen doc for complex flows |
| `getAuthToken` fails in Brave/Edge | Chrome-only API | Switch to `launchWebAuthFlow` for cross-browser compatibility |
| Duplicate refresh calls | Multiple contexts trigger refresh simultaneously | Use the coalescing lock pattern |
| `invalid_grant` on refresh | Refresh token revoked or expired | Clear stored tokens and force re-authentication |
| Redirect URI works in dev but not prod | Extension ID differs between unpacked and CWS | Pin `"key"` in manifest.json or register both redirect URIs |

## Manifest Quick Reference

### Google-only (getAuthToken)
```json
{
  "manifest_version": 3,
  "permissions": ["identity"],
  "oauth2": {
    "client_id": "123456789.apps.googleusercontent.com",
    "scopes": ["openid", "email", "profile"]
  }
}
```

### Any provider (launchWebAuthFlow)
```json
{
  "manifest_version": 3,
  "permissions": ["identity"],
  "host_permissions": [
    "https://accounts.google.com/*",
    "https://oauth2.googleapis.com/*"
  ]
}
```

`host_permissions` are needed for direct `fetch()` calls to token endpoints. If proxying through your own backend, list your backend's origin instead.
