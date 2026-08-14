<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `glean-llm-client-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: glean-llm-client-patterns
version: 1.1.0
updated: 2026-05-29
description: >
  Glean AI platform API integration patterns for JavaScript clients, Chrome MV3
  extensions, and agent workflows. Covers auth (bearer tokens, browser-session
  cookies, OAuth), chat/search REST endpoints, MCP JSON-RPC session flow, streaming
  SSE, rate limiting, error handling, manual prompt fallback, and Chrome extension
  tab-based auth.
  TRIGGER: building or reviewing code that calls Glean APIs; configuring Glean auth
  in a Chrome extension or Node app; debugging Glean 401/403/429 errors; implementing
  Glean MCP JSON-RPC flow; designing fallback when Glean is unavailable.
  SKIP: Glean platform administration or connector setup (use glean-dev); Mozilla
  Glean telemetry SDK (unrelated product); general MCP protocol not specific to
  Glean (use mcp-servers).
category: developer
tags: [glean, chrome-extension, authentication, bearer-token, mcp, streaming, rate-limiting, service-worker]
related_skills: [glean-dev, chrome-mv3-advanced, mcp-servers, chrome-native-messaging]
whenToUse:
  - "Building or debugging Glean API calls from a Chrome extension"
  - "Glean auth not working in my service worker"
  - "How do I implement Glean tab-based browser session auth?"
  - "Glean MCP JSON-RPC initialize flow"
  - "Implement streaming from Glean chat endpoint"
  - "Handle Glean 429 rate limit with exponential backoff"
  - "Fallback to manual prompt when Glean is unavailable"
  - "Glean bearer token vs session cookie auth"
  - "Chrome extension Glean integration patterns"
whenNotToUse:
  - "For Glean CLI, SDK setup, Indexing API, admin/governance: use glean-dev"
  - "For Mozilla Glean telemetry instrumentation: use glean-dev (Section 13)"
  - "For general MCP protocol (not Glean-specific): use mcp-servers"
---

# Glean LLM Client Patterns

Expert reference for integrating with the Glean enterprise AI platform from
JavaScript clients, Chrome extensions, and agent workflows.

## When to use this skill

- Building or debugging Glean API calls (search, chat, agents, MCP)
- Configuring Glean authentication (bearer tokens, browser-session cookies, OAuth)
- Implementing streaming responses from Glean chat/agent endpoints
- Handling Glean rate limits, error codes, and retry logic
- Designing manual prompt fallback when Glean is unavailable
- Integrating Glean into Chrome MV3 extensions via service workers and tab-based auth
- Connecting to the Glean MCP remote server from any MCP-compatible host

## When NOT to use this skill

- For Glean CLI, SDK setup, Indexing API, or admin/governance: use `glean-dev`.
- For Mozilla Glean telemetry SDK: use `glean-dev` (Section 13).
- For general MCP protocol details not specific to Glean: use `mcp-servers`.

---

## 1. Glean API landscape

Glean exposes three API surfaces:

| Surface | Base URL pattern | Auth | Use case |
|---------|-----------------|------|----------|
| Client REST API | `https://{instance}-be.glean.com/rest/api/v1/*` | Bearer token or OAuth | Search, chat, documents, agents, people |
| Indexing API | `https://{instance}-be.glean.com/api/index/v1/*` | Indexing token | Push content, permissions, datasource config |
| MCP Remote Server | `https://{instance}-be.glean.com/mcp/default` | Bearer or browser-session | JSON-RPC tool invocation (search, chat, read) |

The `-be` suffix stands for "backend" and is required for API calls. The instance name is your organization slug (e.g., `mongodb-be.glean.com`).

---

## 2. Authentication

### 2.1 Token types

Glean issues two scopes of Client API tokens from the Admin Console (Platform > API tokens):

**Global tokens** can act as any user. Every request MUST include the `X-Scio-ActAs` header:

```js
const headers = {
  Authorization: `Bearer ${GLEAN_GLOBAL_TOKEN}`,
  'X-Scio-ActAs': 'operator@company.com',
  'Content-Type': 'application/json'
};
```

**User-scoped tokens** are bound to a single user at creation time. Omit `X-Scio-ActAs`:

```js
const headers = {
  Authorization: `Bearer ${GLEAN_USER_TOKEN}`,
  'Content-Type': 'application/json'
};
```

Token secrets are shown once at creation. Set an expiry date and rotate proactively.

### 2.2 OAuth access tokens

For production integrations, Glean recommends OAuth access tokens issued by your SSO provider. These are scoped to the individual user and avoid storing long-lived API keys.

### 2.3 Browser-session authentication (Chrome extension pattern)

When no bearer token is available, a Chrome extension can piggyback on an authenticated Glean tab's session by executing fetch calls inside that tab via `chrome.scripting.executeScript`.

Resolution order:

1. Check for a configured tab ID matching the endpoint origin
2. Look up a previously persisted managed auth tab from `chrome.storage.session`
3. Find any open tab on the endpoint origin (`chrome.tabs.query`)
4. Find any open tab on `*.glean.com`
5. Open a background pinned tab to warm the session
6. If warm-up fails, open a visible tab and prompt the user to sign in

```js
// Simplified tab-based fetch (runs inside the Glean tab's context)
async function fetchViaTab(tabId, url, body) {
  const [result] = await chrome.scripting.executeScript({
    target: { tabId },
    args: [url, body],
    func: async (requestUrl, requestBody) => {
      try {
        const response = await fetch(requestUrl, {
          method: 'POST',
          credentials: 'include',  // sends session cookies
          headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(requestBody)
        });
        const text = await response.text();
        if (!response.ok) return { ok: false, status: response.status, text };
        return { ok: true, payload: JSON.parse(text) };
      } catch (err) {
        return { ok: false, status: 0, text: String(err.message) };
      }
    }
  });
  return result?.result;
}
```

`credentials: 'include'` is required — it carries the Glean browser-session cookies.

### 2.4 Auth fallback cascade

A robust client tries multiple auth paths before failing:

1. **Bearer token** — if configured, send `Authorization: Bearer {token}` directly from the service worker
2. **Tab-based session** — if bearer fails with 401/403, retry through an authenticated tab
3. **MCP flow** — if the endpoint responds with JSON-RPC, switch to the full MCP initialize/tools/call flow
4. **Interactive sign-in** — open a Glean tab and throw `auth_interaction_required` so the UI can prompt the user

```js
function buildAuthHeaders(token = '') {
  const normalized = token.trim();
  return normalized ? { Authorization: `Bearer ${normalized}` } : {};
}
```

---

## 3. Chat API

### 3.1 REST endpoint

`POST https://{instance}-be.glean.com/rest/api/v1/chat`

Requires a CHAT-scoped bearer token. Global tokens cover all scopes but require `X-Scio-ActAs`.

### 3.2 Request schema

```js
const chatRequest = {
  messages: [
    {
      author: 'USER',
      fragments: [{ text: 'Summarize our Q4 OKR progress' }]
    }
  ],
  chatId: 'existing-chat-id-here',  // optional: continue a multi-turn conversation
  stream: true                       // optional: enable streaming SSE response
};
```

### 3.3 Response schema (non-streaming)

```json
{
  "chatId": "abc123",
  "messages": [
    {
      "author": "GLEAN_AI",
      "fragments": [
        { "text": "Based on your company's documents..." },
        {
          "citation": {
            "sourceDocument": {
              "title": "Q4 OKR Tracker",
              "url": "https://docs.google.com/...",
              "snippet": "Engineering completed 85% of..."
            }
          }
        }
      ]
    }
  ]
}
```

### 3.4 Streaming responses (SSE)

When `stream: true`, the endpoint returns `text/event-stream`:

```js
async function streamGleanChat(endpoint, headers, messages) {
  const response = await fetch(endpoint, {
    method: 'POST',
    headers: { ...headers, Accept: 'text/event-stream' },
    body: JSON.stringify({ messages, stream: true })
  });

  if (!response.ok) {
    const text = await response.text();
    throw new Error(`Glean chat stream failed (${response.status}): ${text}`);
  }

  const reader = response.body.getReader();
  const decoder = new TextDecoder();
  let buffer = '';

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    buffer += decoder.decode(value, { stream: true });
    const lines = buffer.split('\n');
    buffer = lines.pop();  // keep incomplete line in buffer

    for (const line of lines) {
      if (!line.startsWith('data: ')) continue;
      const data = line.slice(6);
      if (data === '[DONE]') return;
      try {
        const event = JSON.parse(data);
        handleFragment(event);
      } catch { /* skip malformed SSE lines */ }
    }
  }
}
```

---

## 4. Search API

### 4.1 REST endpoint

`POST https://{instance}-be.glean.com/rest/api/v1/search`

### 4.2 Request schema

```js
const searchRequest = {
  query: 'MongoDB Atlas connection timeout',
  pageSize: 10,
  requestOptions: {
    datasourcesFilter: ['confluence', 'jira', 'slack'],
    facetFilters: [
      { fieldName: 'author', values: [{ value: 'jane@company.com' }] }
    ],
    facetBucketSize: 5
  }
};
```

### 4.3 Response shape

```json
{
  "results": [
    {
      "document": {
        "title": "Atlas Connectivity Troubleshooting",
        "url": "https://confluence.example.com/...",
        "datasource": "confluence"
      },
      "snippets": [{ "snippet": "...connection timeout usually indicates..." }]
    }
  ],
  "facetResults": [],
  "totalCount": 42
}
```

---

## 5. MCP Remote Server integration

Glean exposes a built-in MCP remote server — the preferred integration path for agent frameworks.

### 5.1 Endpoint

`https://{instance}-be.glean.com/mcp/default`

### 5.2 JSON-RPC session flow

The MCP server uses JSON-RPC 2.0 over HTTP with session tracking via `Mcp-Session-Id`:

```js
async function gleanMcpSession(endpoint, authHeaders, prompt) {
  let sessionId = '';
  let requestId = 1;

  async function transport(payload) {
    const headers = {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      ...(sessionId ? { 'Mcp-Session-Id': sessionId } : {}),
      ...authHeaders
    };
    const response = await fetch(endpoint, {
      method: 'POST',
      headers,
      body: JSON.stringify(payload)
    });
    const text = await response.text();
    if (!response.ok) throw new Error(`MCP error ${response.status}: ${text}`);
    const nextSession = response.headers.get('mcp-session-id');
    if (nextSession) sessionId = nextSession;
    return JSON.parse(text);
  }

  function request(method, params) {
    return transport({ jsonrpc: '2.0', id: requestId++, method,
      ...(params !== undefined ? { params } : {}) });
  }

  function notify(method, params) {
    return transport({ jsonrpc: '2.0', method,
      ...(params !== undefined ? { params } : {}) }).catch(() => {});
  }

  // Step 1: Initialize
  await request('initialize', {
    protocolVersion: '2025-06-18',
    capabilities: {},
    clientInfo: { name: 'my-app', version: '1.0.0' }
  });
  await notify('notifications/initialized', {});

  // Step 2: Discover tools
  const toolsResult = await request('tools/list', {});
  const tools = toolsResult?.tools || [];

  // Step 3: Select and invoke the chat tool
  const chatTool = tools.find(t => t.name.includes('chat'));
  if (!chatTool) throw new Error('No chat tool found');

  // Step 4: Call the tool
  return request('tools/call', {
    name: chatTool.name,
    arguments: buildToolArgs(chatTool, prompt)
  });
}
```

### 5.3 Tool selection heuristic

Priority order used when selecting from available tools:

1. Exact name `glean_default-chat`
2. Exact name `chat`
3. Name ending with `chat`
4. Name containing `chat` anywhere

### 5.4 Dynamic argument mapping

Different Glean MCP tool versions accept different parameter names — map dynamically:

```js
function buildToolArgs(tool, prompt) {
  const props = tool?.inputSchema?.properties || {};
  if (props.message) return { message: prompt };
  if (props.query)   return { query: prompt };
  if (props.prompt)  return { prompt };
  for (const [key, spec] of Object.entries(props)) {
    if (spec?.type === 'string') return { [key]: prompt };
  }
  return { message: prompt };  // last resort
}
```

### 5.5 JSON-RPC error handling

```js
function unwrapJsonRpcResult(payload) {
  if (payload?.jsonrpc && payload?.error) {
    const err = new Error(payload.error.message || 'MCP request failed');
    err.code = payload.error.code;
    err.data = payload.error.data;
    err.jsonRpc = true;
    throw err;
  }
  return payload?.result ?? payload;
}
```

---

## 6. Payload format auto-detection

When the endpoint type is unknown, try three formats in sequence:

```js
const attempts = [
  { query: prompt },                                    // Glean search-style
  { prompt },                                           // Glean direct prompt
  { model: 'gpt-4.1-mini',                             // OpenAI-compatible
    messages: [{ role: 'user', content: prompt }] }
];

for (const body of attempts) {
  try {
    const result = await postJson(endpoint, headers, body);
    if (result?.jsonrpc) {
      return await gleanMcpSession(endpoint, headers, prompt);
    }
    return result;
  } catch (err) {
    if (err.jsonRpc) break;
    lastError = err;
  }
}
```

This lets a single client work against Glean REST, Glean MCP, and OpenAI-compatible endpoints without explicit configuration.

---

## 7. Response text extraction

Glean responses vary by endpoint and tool version — use recursive extraction:

```js
function extractText(value, seen = new Set()) {
  if (!value || seen.has(value)) return [];
  if (typeof value === 'string') return [value.trim()].filter(Boolean);
  if (typeof value !== 'object') return [];
  seen.add(value);

  const parts = [];
  const push = (s) => { if (typeof s === 'string' && s.trim()) parts.push(s.trim()); };
  const visit = (v) => parts.push(...extractText(v, seen));

  push(value.answer);
  push(value.response);
  push(value.text);
  push(value.content);
  push(value.output_text);

  for (const key of ['content', 'parts', 'messages', 'choices', 'output', 'candidates']) {
    if (Array.isArray(value[key])) value[key].forEach(visit);
  }
  for (const key of ['message', 'delta', 'data', 'result']) {
    if (value[key]) visit(value[key]);
  }

  return parts;
}

function getResponseText(payload) {
  return extractText(payload).join('\n').trim();
}
```

This handles Glean chat fragments, OpenAI choices, Gemini candidates, and MCP tool call results uniformly.

---

## 8. Rate limiting

Glean uses token-bucket rate limiters at both user and endpoint levels.

| Signal | Action |
|--------|--------|
| HTTP 429 | Exponential backoff: 1s, 2s, 4s, 8s, cap at 30s |
| `Retry-After` header | Respect the value (seconds or HTTP-date) |
| Consistent 429s | Reduce concurrency; batch or queue requests |

```js
async function fetchWithRetry(url, options, maxRetries = 3) {
  for (let attempt = 0; attempt <= maxRetries; attempt++) {
    const response = await fetch(url, options);
    if (response.status !== 429) return response;

    const retryAfter = response.headers.get('Retry-After');
    const delayMs = retryAfter
      ? (Number(retryAfter) || 1) * 1000
      : Math.min(1000 * 2 ** attempt, 30000);
    await new Promise(r => setTimeout(r, delayMs));
  }
  throw new Error('Rate limited after max retries');
}
```

---

## 9. Error handling and classification

### 9.1 HTTP error mapping

```js
function classifyGleanError(error) {
  const status = Number(error?.status || 0);

  if (error?.authRequired) return {
    code: 'auth_interaction_required',
    message: 'Opened a Glean sign-in tab. Finish signing in, then retry.',
    recoverable: true
  };
  if (status === 0)   return { code: 'network_error', message: 'Cannot reach Glean. Check endpoint URL and network.', recoverable: true };
  if (status === 401) return { code: 'unauthorized', message: 'Token rejected. Sign in again or update the bearer token.', recoverable: true };
  if (status === 403) return { code: 'forbidden', message: 'Authenticated but access denied. Check account authorization.', recoverable: false };
  if (status === 404) return { code: 'not_found', message: 'Endpoint not found. Verify the Glean API path.', recoverable: false };
  if (status === 429) return { code: 'rate_limited', message: 'Too many requests. Back off and retry.', recoverable: true };
  if (status >= 500)  return { code: 'server_error', message: 'Glean server error. Retry shortly.', recoverable: true };
  return { code: 'endpoint_mismatch', message: 'Unexpected response. Verify the Glean API endpoint.', recoverable: false };
}
```

### 9.2 Error enrichment pattern

```js
function createRequestError(response, text = '') {
  const error = new Error(
    `Provider request failed (${response.status}): ${text || response.statusText}`
  );
  error.status = response.status;
  error.statusText = response.statusText;
  error.responseText = text;
  return error;
}
```

---

## 10. Manual prompt fallback

When Glean is unavailable, degrade gracefully rather than block the operator:

```js
async function analyzeCase({ provider, options, caseContext }) {
  const prompt = buildAnalysisPrompt(caseContext);

  if (provider === 'manual') {
    return { mode: 'manual', prompt, analysis: buildFallbackAnalysis(caseContext) };
  }

  try {
    const payload = await requestGlean(prompt, options);
    const text = getResponseText(payload);
    return { mode: 'glean', prompt, analysis: parseAnalysis(text) };
  } catch (error) {
    return {
      mode: provider,
      prompt,
      analysis: buildFallbackAnalysis(caseContext, { fallbackReason: error.message })
    };
  }
}
```

The fallback preserves the case context and explains why the automated path failed, so the operator can take the generated prompt to any alternative LLM.

---

## 11. Chrome extension integration patterns

### 11.1 Session storage for auth state

Store disposable auth data in `chrome.storage.session` (cleared on browser close):

```js
const GLEAN_AUTH_SESSION_KEY = 'mca_glean_auth_session_v1';

async function getGleanAuthSession() {
  const data = await chrome.storage.session.get(GLEAN_AUTH_SESSION_KEY);
  return data?.[GLEAN_AUTH_SESSION_KEY] || null;
}

async function setGleanAuthSession(value) {
  await chrome.storage.session.set({ [GLEAN_AUTH_SESSION_KEY]: value });
}

async function clearGleanAuthSession() {
  await chrome.storage.session.remove(GLEAN_AUTH_SESSION_KEY);
}
```

### 11.2 Managed auth tab lifecycle

```js
async function openManagedAuthTab(authUrl, { interactive = false } = {}) {
  const existing = await getStoredAuthTab();
  if (existing) {
    return interactive ? await focusTab(existing) : existing;
  }

  const tab = await chrome.tabs.create({
    url: authUrl,
    active: interactive,
    pinned: !interactive
  });

  await waitForTabLoad(tab.id, 15000);
  return tab;
}
```

### 11.3 Tab load wait pattern

```js
async function waitForTabLoad(tabId, timeoutMs = 15000) {
  const existing = await chrome.tabs.get(tabId);
  if (existing.status === 'complete') return;

  return new Promise((resolve, reject) => {
    let timer;
    const cleanup = (err) => {
      clearTimeout(timer);
      chrome.tabs.onUpdated.removeListener(onUpdated);
      chrome.tabs.onRemoved.removeListener(onRemoved);
      err ? reject(err) : resolve();
    };
    const onUpdated = (id, info) => {
      if (id === tabId && info.status === 'complete') cleanup();
    };
    const onRemoved = (id) => {
      if (id === tabId) cleanup(new Error('Auth tab closed'));
    };
    chrome.tabs.onUpdated.addListener(onUpdated);
    chrome.tabs.onRemoved.addListener(onRemoved);
    timer = setTimeout(() => cleanup(new Error('Auth tab load timeout')), timeoutMs);
  });
}
```

### 11.4 Connection test probe

```js
async function probeGleanConnection({ endpoint, token, gleanTab }) {
  const headers = token ? { Authorization: `Bearer ${token}` } : {};
  const payload = await requestGleanViaMcp('Reply with OK.', {
    gleanEndpoint: endpoint,
    gleanToken: token
  }, gleanTab, headers);
  const text = getResponseText(payload);
  if (!text) throw new Error('Probe returned empty response');
  return { ok: true, text };
}
```

### 11.5 Status reporting pattern

```js
function reportStatus(reporter, payload) {
  if (typeof reporter !== 'function') return;
  reporter({ timestamp: new Date().toISOString(), ...payload });
}

// Usage in async flow
reportStatus(statusReporter, {
  step: 'auth-resolution-start',
  message: 'Resolving Glean authentication...',
  level: 'info',
  details: { endpoint, authMode: token ? 'bearer' : 'tab' }
});
```

---

## 12. TypeScript SDK (npm)

For Node.js or non-extension web apps:

```bash
npm add @gleanwork/api-client
```

```js
import { Glean } from '@gleanwork/api-client';

const glean = new Glean({
  apiToken: process.env.GLEAN_API_TOKEN,
  instance: 'your-instance'
});

// Chat
const chatResponse = await glean.client.chat.create({
  messages: [{ fragments: [{ text: 'What is our PTO policy?' }] }]
});

// Search
const searchResponse = await glean.client.search.query({
  query: 'deployment runbook',
  pageSize: 5
});
```

---

## 13. Agents API (SSE streaming)

`POST https://{instance}-be.glean.com/rest/api/v1/agents/{agentId}/runs`

```js
async function runGleanAgent(endpoint, headers, agentId, input) {
  const response = await fetch(
    `${endpoint}/rest/api/v1/agents/${agentId}/runs`,
    {
      method: 'POST',
      headers: { ...headers, Accept: 'text/event-stream' },
      body: JSON.stringify({ input, stream: true })
    }
  );

  if (!response.ok) {
    const text = await response.text();
    throw new Error(`Agent run failed (${response.status}): ${text}`);
  }

  const reader = response.body.getReader();
  const decoder = new TextDecoder();
  let buffer = '';
  const events = [];

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;
    buffer += decoder.decode(value, { stream: true });

    const lines = buffer.split('\n');
    buffer = lines.pop();
    for (const line of lines) {
      if (!line.startsWith('data: ')) continue;
      const data = line.slice(6);
      if (data === '[DONE]') return events;
      try { events.push(JSON.parse(data)); } catch { /* skip malformed */ }
    }
  }

  return events;
}
```

---

## 14. Integration review checklist

- [ ] Auth: token stored in `chrome.storage.session` or env vars — never in source
- [ ] Auth: global tokens always include `X-Scio-ActAs` header
- [ ] Auth: user-scoped tokens never include `X-Scio-ActAs` header
- [ ] Auth: token rotation and expiry handling implemented
- [ ] Auth: browser-session fallback uses `credentials: 'include'`
- [ ] Endpoints: base URL uses `{instance}-be.glean.com` pattern
- [ ] MCP: session ID tracked via `Mcp-Session-Id` header across requests
- [ ] MCP: `initialize` + `notifications/initialized` sent before `tools/list`
- [ ] MCP: tool selection handles schema changes gracefully
- [ ] Errors: 429 handled with exponential backoff and `Retry-After` respect
- [ ] Errors: 401/403 trigger auth refresh or tab-based fallback, not hard failure
- [ ] Errors: structured error objects carry `status`, `code`, and `responseText`
- [ ] Fallback: manual prompt path available when all automated paths fail
- [ ] Streaming: SSE `[DONE]` sentinel handled to close the stream reader
- [ ] Extension: auth state in `chrome.storage.session`, not `chrome.storage.local`
- [ ] Extension: managed auth tabs cleaned up when session ends
- [ ] Logging: request/response lifecycle logged with structured fields, not raw payloads

---

## 15. Anti-patterns

| Anti-pattern | Fix |
|-------------|-----|
| Hardcoding the endpoint | Always read from config — instances vary by org |
| Ignoring JSON-RPC detection | An MCP endpoint returns `jsonrpc` payloads even for raw POST — detect and switch |
| Retrying 401/403 without re-auth | These mean the credential is wrong, not that the server is busy |
| Logging bearer tokens | Log `hasToken: Boolean(token)`, never the token value |
| Blocking on Glean failure | Always provide a manual fallback so the operator workflow continues |
| `chrome.storage.local` for session tokens | Session data must not survive browser restart |
| Skipping the `initialize` handshake | MCP servers reject `tools/call` without a prior session |
| Assuming a fixed tool schema | Glean MCP tool names and `inputSchema` properties can change between versions |
