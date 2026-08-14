<!-- hub-reference-banner -->
> **Reference file — part of the `tam-operations` hub.** Formerly the standalone `slack-subscription-auditor` skill.
> Sibling topics in this family are now reference files under the hubs (`tam-operations`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: slack-subscription-auditor
version: "1.3.0"
updated: "2026-06-01"
description: >-
  Audit Slack workspaces and apps for health and compliance: slash-command
  inventory, dead-command detection, event-subscription lifecycle, orphaned-
  subscription detection, bot-token scope auditing (OAuth least privilege),
  rate-limit handling, Enterprise Grid Audit Logs API, webhook validation, and
  required-channel compliance checks.
  TRIGGER: "audit Slack slash commands", "detect orphaned Slack subscriptions",
  "check bot token scopes for least privilege", "build a Slack app health
  checker", "paginated sweep of Slack APIs", "Slack compliance channel audit",
  "scope drift on bot tokens", "Socket Mode app health check".
  SKIP: building a Slack bot, calling Web API, Block Kit, Events API, or Bolt SDK
  → integration-clients (references/slack-dev.md); general SSE/webhook
  infrastructure → software-engineering-patterns (references/sse-streaming-patterns.md);
  Chrome extension integrating with Slack
  → chrome-extension-expert (references/chrome-dev.md).
category: developer
tags: [slack, audit, subscriptions, commands, events-api, compliance, oauth, bot-token, webhook, enterprise-grid]
whenToUse:
  - "audit Slack slash command registrations across a workspace"
  - "detect stale or orphaned Slack event subscriptions"
  - "build an automated Slack app health checker or auditor bot"
  - "review bot token scopes for excessive permissions"
  - "check for scope drift on bot tokens week-over-week"
  - "design compliance checks for required channels and notification routing"
  - "implement paginated audit sweeps against Slack APIs"
  - "detect dead or abandoned slash commands"
  - "audit Socket Mode app health and WebSocket connection validity"
  - "sweep Enterprise Grid Audit Logs API for app install/scope-change events"
whenNotToUse:
  - "building a Slack bot, writing API calls, or implementing Block Kit → integration-clients (references/slack-dev.md)"
  - "general SSE/webhook infrastructure not Slack-specific → software-engineering-patterns (references/sse-streaming-patterns.md)"
  - "Chrome extension integrating with Slack → chrome-extension-expert (references/chrome-dev.md)"
related_skills:
  - integration-clients
  - software-engineering-patterns
  - chrome-extension-expert
---

# Slack Subscription Auditor

## Overview

Autonomous audit patterns for Slack workspaces and apps. The audit enumerates registered slash commands, validates event subscriptions, detects dead or orphaned integrations, checks bot-token scopes against least-privilege, handles rate limits during large sweeps, and verifies compliance with notification-routing rules.

## Prerequisites for the Auditor Bot

Request only the scopes the audit requires:

| Scope | Why the auditor needs it |
|---|---|
| `commands` | Validate slash command registrations |
| `channels:read` | Enumerate public channels for compliance |
| `groups:read` | Enumerate private channels (if auditing private) |
| `users:read` | Check app owner active/deactivated status |
| `team:read` | Get workspace metadata |
| `chat:write` | Post audit reports to an alert channel |

For Enterprise Grid audits, also request:
- `auditlogs:read` scope (org-level install by an Org Owner)
- `admin.apps:read` scope for listing approved/restricted apps

---

## Core Concepts

### Slack App Architecture (Audit-Relevant Surfaces)

```
Slack Workspace
  |
  +-- Installed Apps (each has a bot token, scopes, manifest)
  |     +-- Slash Commands (registered via manifest or API)
  |     +-- Event Subscriptions (HTTP or Socket Mode)
  |     +-- Interactive Components (shortcuts, modals, actions)
  |     +-- OAuth Scopes (bot token + user token)
  |
  +-- Audit Logs API (Enterprise Grid only)
  |     +-- Actor / Action / Entity / Context model
  |     +-- App install/uninstall/scope-change events
  |
  +-- Admin APIs
        +-- admin.apps.approved.list
        +-- admin.apps.restricted.list
        +-- admin.apps.requests.list
```

### Token Types and Audit Relevance

| Token prefix | Type | Audit concern |
|---|---|---|
| `xoxb-` | Bot token | Scope creep, unused permissions |
| `xoxp-` | User token | Over-privileged user delegation |
| `xoxe-` | Enterprise token | Org-wide admin access |
| `xapp-` | App-level token | Socket Mode connections |
| legacy workspace token | Deprecated | Must not exist in modern apps; flag any non-`xoxb-`/`xoxp-`/`xapp-` token |

### Key API Methods for Auditing

| Method | Tier | Purpose |
|---|---|---|
| `auth.test` | Tier 4 | Validate token, get bot/team identity |
| `team.info` | Tier 3 | Get workspace metadata |
| `users.list` | Tier 2 | Enumerate workspace members |
| `conversations.list` | Tier 2 | Enumerate channels for compliance |
| `apps.manifest.export` | Tier 3 | Export app manifest (commands, scopes) |
| `admin.apps.approved.list` | Tier 2 | List approved apps (admin) |
| `admin.apps.restricted.list` | Tier 2 | List restricted apps (admin) |

---

## Command Audit Workflow

### Step 1: Build a Command Inventory

```javascript
async function getRegisteredCommands(appToken) {
  const res = await fetch('https://slack.com/api/apps.manifest.export', {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${appToken}`, 'Content-Type': 'application/json' },
    body: JSON.stringify({ app_id: APP_ID }),
  });
  const data = await res.json();
  if (!data.ok) throw new Error(`manifest export failed: ${data.error}`);

  const commands = data.manifest?.features?.slash_commands || [];
  return commands.map(cmd => ({
    command: cmd.command,
    url: cmd.url,
    description: cmd.description,
    usage_hint: cmd.usage_hint || null,
  }));
}
```

### Step 2: Validate Each Command Endpoint

```javascript
async function validateCommandEndpoint(command) {
  const result = { command: command.command, url: command.url, reachable: false, sslValid: false, respondsIn3s: false, status: 'unknown' };
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 3000);

  try {
    const res = await fetch(command.url, {
      method: 'POST',
      signal: controller.signal,
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: new URLSearchParams({ _audit_probe: '1' }),
    });
    result.reachable = true;
    result.sslValid = true;  // fetch throws on bad SSL
    result.respondsIn3s = true;
    result.status = res.status < 500 ? 'healthy' : 'server_error';
  } catch (err) {
    if (err.name === 'AbortError') { result.reachable = true; result.status = 'timeout'; }
    else if (err.code === 'UNABLE_TO_VERIFY_LEAF_SIGNATURE') { result.status = 'ssl_invalid'; }
    else { result.status = 'unreachable'; }
  } finally { clearTimeout(timeout); }

  return result;
}
```

### Step 3: Cross-Reference with Usage Data (Enterprise Grid)

```javascript
async function findDeadCommands(registeredCommands, auditClient) {
  const thirtyDaysAgo = Math.floor(Date.now() / 1000) - (30 * 86400);
  const usedCommands = new Set();
  let cursor = '';

  do {
    const logs = await auditClient.logs({ action: 'user_channel_slash_command', oldest: thirtyDaysAgo, cursor, limit: 200 });
    for (const entry of logs.entries || []) {
      if (entry.entity?.command) usedCommands.add(entry.entity.command);
    }
    cursor = logs.response_metadata?.next_cursor || '';
  } while (cursor);

  return registeredCommands
    .filter(cmd => !usedCommands.has(cmd.command))
    .map(cmd => ({ ...cmd, status: 'dead_command', lastUsed: null }));
}
```

### Command Audit Checklist

- [ ] All registered commands have reachable HTTPS endpoints
- [ ] All endpoints respond within 3 seconds (Slack timeout)
- [ ] SSL certificates are valid and not expiring within 30 days
- [ ] No commands registered in manifest that lack a handler in code
- [ ] No commands in code that are missing from the manifest
- [ ] Usage data confirms each command has been invoked in last 90 days
- [ ] Commands owned by deactivated users are identified and reassigned

---

## Event Subscription Lifecycle

### Subscription States

```
CREATED → ACTIVE → STALE → ORPHANED
              |
              +→ HEALTHY (>95% ack rate)
              +→ DEGRADED (5–95% ack rate)
              +→ DISABLED (auto, <5% success rate)
```

### Challenge Verification Audit

```javascript
async function auditChallengeVerification(requestUrl) {
  const challenge = crypto.randomUUID();
  try {
    const res = await fetch(requestUrl, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ type: 'url_verification', challenge, token: 'audit_probe' }),
    });
    const body = await res.text();
    let parsed;
    try { parsed = JSON.parse(body); } catch { parsed = null; }
    const echoedChallenge = parsed?.challenge || body.trim();
    return {
      url: requestUrl,
      status: echoedChallenge === challenge ? 'pass' : 'fail',
      httpStatus: res.status,
      detail: echoedChallenge === challenge ? 'Challenge echoed correctly' : `Expected "${challenge}", got "${echoedChallenge}"`,
    };
  } catch (err) { return { url: requestUrl, status: 'fail', detail: err.message }; }
}
```

### Event Delivery Health Monitoring

Slack retries failed event deliveries 3 times with exponential backoff. If success rate drops below 5% over 60 minutes, Slack automatically disables the subscription.

```javascript
function computeDeliveryHealth(eventLog, windowMinutes = 60) {
  const cutoff = Date.now() - windowMinutes * 60 * 1000;
  const recent = eventLog.filter(e => e.timestamp > cutoff);
  const total = recent.length;
  const successful = recent.filter(e => e.ackStatus === 200).length;
  const rate = total > 0 ? (successful / total) * 100 : 0;

  return {
    totalEvents: total, successfulAcks: successful,
    successRate: rate.toFixed(1) + '%',
    status: rate >= 95 ? 'healthy' : rate >= 5 ? 'degraded' : 'critical_will_be_disabled',
    recommendation: rate < 5
      ? 'URGENT: Success rate below 5%. Slack will auto-disable delivery.'
      : rate < 95 ? 'WARNING: Investigate failed event acknowledgements.' : 'OK',
  };
}
```

### Detecting Orphaned Subscriptions

```javascript
async function detectOrphanedSubscriptions(apps, slackAdmin) {
  const findings = [];
  for (const app of apps) {
    const checks = { appId: app.id, appName: app.name, ownerActive: false, endpointReachable: false, hasActiveTokens: false };

    try { const owner = await slackAdmin.users.info({ user: app.owner_id }); checks.ownerActive = !owner.user.deleted; }
    catch { checks.ownerActive = false; }

    if (app.event_subscriptions?.request_url) {
      try { const res = await fetch(app.event_subscriptions.request_url, { method: 'HEAD', signal: AbortSignal.timeout(5000) }); checks.endpointReachable = res.status < 500; }
      catch { checks.endpointReachable = false; }
    }

    try { const auth = await slackAdmin.auth.test({ token: app.bot_token }); checks.hasActiveTokens = auth.ok; }
    catch { checks.hasActiveTokens = false; }

    if (!checks.ownerActive || !checks.endpointReachable || !checks.hasActiveTokens) {
      findings.push({
        ...checks, status: 'orphaned',
        reasons: [!checks.ownerActive && 'owner_deactivated', !checks.endpointReachable && 'endpoint_unreachable', !checks.hasActiveTokens && 'token_invalid'].filter(Boolean),
      });
    }
  }
  return findings;
}
```

### Signing Secret Validation Audit

Every event handler and slash command endpoint must validate `X-Slack-Signature`. Skipping this is a **critical** finding.

```javascript
const crypto = require('crypto');

function verifySlackSignature(signingSecret, req) {
  const timestamp = req.headers['x-slack-request-timestamp'];
  const signature = req.headers['x-slack-signature'];

  // Reject requests older than 5 minutes (replay attack prevention)
  if (parseInt(timestamp, 10) < Math.floor(Date.now() / 1000) - 300) {
    return { valid: false, reason: 'timestamp_expired' };
  }

  const sigBasestring = `v0:${timestamp}:${req.rawBody}`;
  const hmac = crypto.createHmac('sha256', signingSecret);
  hmac.update(sigBasestring);
  const expected = `v0=${hmac.digest('hex')}`;
  const valid = crypto.timingSafeEqual(Buffer.from(signature), Buffer.from(expected));
  return { valid, reason: valid ? 'ok' : 'signature_mismatch' };
}
```

### Socket Mode Audit

```javascript
async function auditSocketModeApp(appToken) {
  if (!appToken.startsWith('xapp-')) {
    return { status: 'not_socket_mode', detail: 'Token is not an app-level token' };
  }
  try {
    const res = await fetch('https://slack.com/api/apps.connections.open', {
      method: 'POST', headers: { 'Authorization': `Bearer ${appToken}` }
    });
    const data = await res.json();
    if (!data.ok) {
      return { status: 'connection_failed', error: data.error,
        detail: data.error === 'invalid_auth' ? 'App-level token is revoked or invalid' : `Connection failed: ${data.error}` };
    }
    return { status: 'healthy', wsUrl: data.url ? 'present' : 'missing' };
  } catch (err) { return { status: 'error', detail: err.message }; }
}
```

Socket Mode audit checklist:
- [ ] App-level token (`xapp-`) is valid and not revoked
- [ ] `apps.connections.open` returns a WebSocket URL
- [ ] No public request URL is configured (Socket Mode replaces it)
- [ ] Reconnection logic handles token refresh (URLs expire ~1 hour)

### Event Subscription Audit Checklist

- [ ] All event request URLs pass challenge verification
- [ ] Signing secret validation is implemented on every endpoint (critical)
- [ ] Event delivery success rate is above 95% per app
- [ ] No subscriptions are in auto-disabled state
- [ ] Socket Mode apps have active WebSocket connections
- [ ] Event types subscribed match actual handler implementations
- [ ] `X-Slack-No-Retry: 1` header set when intentionally dropping events
- [ ] App owners are active workspace members

---

## Dead Command Detection

### Manifest-to-Code Diff

```javascript
function diffManifestVsHandlers(manifestCommands, implementedRoutes) {
  const registered = new Set(manifestCommands.map(c => c.command));
  const implemented = new Set(implementedRoutes);
  return {
    registeredButNotImplemented: [...registered].filter(c => !implemented.has(c)),
    implementedButNotRegistered: [...implemented].filter(c => !registered.has(c)),
    fullyMatched: [...registered].filter(c => implemented.has(c)),
  };
}
// Example: diffManifestVsHandlers(['/deploy','/status','/legacy-report'], ['/deploy','/status','/new-feature'])
// registeredButNotImplemented: ['/legacy-report']
// implementedButNotRegistered: ['/new-feature']
```

### Endpoint Liveness with Slack Signature

```javascript
function generateSlackSignature(signingSecret, timestamp, body) {
  const hmac = crypto.createHmac('sha256', signingSecret);
  hmac.update(`v0:${timestamp}:${body}`);
  return `v0=${hmac.digest('hex')}`;
}

async function probeCommandHandler(url, signingSecret, command) {
  const timestamp = Math.floor(Date.now() / 1000).toString();
  const body = new URLSearchParams({ command, text: '__audit_probe__', user_id: 'UAUDITBOT', team_id: 'TAUDIT', response_url: 'https://hooks.slack.com/audit/noop', trigger_id: '0000.0000.noop' }).toString();
  const signature = generateSlackSignature(signingSecret, timestamp, body);
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded', 'X-Slack-Request-Timestamp': timestamp, 'X-Slack-Signature': signature },
    body,
  });
  return { command, httpStatus: res.status, alive: res.status === 200 };
}
```

### Usage Frequency Analysis

```javascript
function classifyCommandHealth(invocationCounts, windowDays = 90) {
  const classifications = {};
  for (const [command, dailyCounts] of Object.entries(invocationCounts)) {
    const total = dailyCounts.reduce((a, b) => a + b, 0);
    const recentThird = dailyCounts.slice(-Math.ceil(windowDays / 3));
    const recentTotal = recentThird.reduce((a, b) => a + b, 0);
    const trend = total > 0 ? recentTotal / (total / 3) : 0;
    classifications[command] = {
      totalInvocations: total, trend: trend.toFixed(2),
      status: total === 0 ? 'dead' : trend < 0.1 ? 'dying' : trend < 0.5 ? 'declining' : 'active',
    };
  }
  return classifications;
}
```

---

## Scope Auditing

### Principle of Least Privilege

Bot tokens (`xoxb-`) should carry only the scopes the app actively uses.

```javascript
const REQUIRED_SCOPES = new Set(['chat:write', 'commands', 'channels:read']);

const HIGH_RISK_SCOPES = new Set([
  'admin', 'admin.apps:write', 'admin.conversations:write', 'admin.users:write',
  'channels:history', 'groups:history', 'im:history', 'mpim:history',
  'users.profile:write', 'files:write',
]);

function auditScopes(grantedScopes) {
  const granted = new Set(grantedScopes);
  return {
    missing: [...REQUIRED_SCOPES].filter(s => !granted.has(s)),
    unnecessary: [...granted].filter(s => !REQUIRED_SCOPES.has(s)),
    highRisk: [...granted].filter(s => HIGH_RISK_SCOPES.has(s)),
    assessment: [...granted].filter(s => HIGH_RISK_SCOPES.has(s)).length > 0
      ? 'REVIEW_REQUIRED'
      : [...granted].filter(s => !REQUIRED_SCOPES.has(s)).length > 2
      ? 'OVER_PROVISIONED' : 'ACCEPTABLE',
  };
}
```

### Scope Drift Detection

```javascript
function detectScopeDrift(previousScopes, currentScopes) {
  const prev = new Set(previousScopes);
  const curr = new Set(currentScopes);
  const added = [...curr].filter(s => !prev.has(s));
  const removed = [...prev].filter(s => !curr.has(s));
  return {
    added, removed, driftDetected: added.length > 0,
    escalation: added.some(s => HIGH_RISK_SCOPES.has(s)),
    summary: added.length === 0 && removed.length === 0 ? 'No scope changes' : `+${added.length} added, -${removed.length} removed`,
  };
}
```

### Scope Audit Checklist

- [ ] Bot token carries only scopes the app actively uses
- [ ] No legacy umbrella `bot` scope (indicates unmigrated app)
- [ ] High-risk scopes (`admin.*`, `*.history`, `users.profile:write`) have documented justification
- [ ] User tokens (`xoxp-`) avoided unless user-context actions are required
- [ ] Scope set compared weekly; any additions trigger review
- [ ] No deprecated legacy workspace tokens exist anywhere in the workspace

---

## Rate Limiting and Pagination

### Rate Limit Tiers

| Tier | Requests/min | Typical methods |
|---|---|---|
| Tier 1 | 1 | `admin.apps.approve` |
| Tier 2 | 20 | `conversations.list`, `users.list` |
| Tier 3 | 50 | `chat.postMessage`, `auth.test` |
| Tier 4 | 100 | `reactions.add` |

### Cursor-Based Pagination

Cursors expire within minutes — do not persist them across sessions.

```javascript
async function paginateSlackAPI(client, method, params = {}, limit = 200) {
  const results = [];
  let cursor = undefined;
  do {
    const res = await client[method]({ ...params, limit, cursor });
    if (!res.ok) throw new Error(`${method} failed: ${res.error}`);
    const dataKey = Object.keys(res).find(k => Array.isArray(res[k]) && k !== 'response_metadata');
    if (dataKey) results.push(...res[dataKey]);
    cursor = res.response_metadata?.next_cursor || undefined;
  } while (cursor);
  return results;
}
```

### Rate Limit Handler with Exponential Backoff

```javascript
async function rateLimitedFetch(fn, maxRetries = 5) {
  for (let attempt = 0; attempt <= maxRetries; attempt++) {
    try {
      return await fn();
    } catch (err) {
      if (err.data?.error === 'ratelimited' || err.status === 429) {
        const retryAfter = parseInt(err.headers?.['retry-after'] || '1', 10);
        const jitter = Math.random() * 500;
        const backoff = retryAfter * 1000 + (attempt > 0 ? 1000 * Math.pow(2, attempt - 1) : 0) + jitter;
        if (attempt === maxRetries) throw err;
        await new Promise(r => setTimeout(r, backoff));
      } else { throw err; }
    }
  }
}
```

### Audit Logs API Pagination (Enterprise Grid)

Rate-limited at Tier 3 (50 req/min) on an org-wide basis.

```javascript
async function sweepAuditLogs(orgToken, action, oldestTimestamp) {
  const entries = [];
  let cursor = '';
  do {
    const res = await rateLimitedFetch(async () => {
      const params = new URLSearchParams({ action, oldest: String(oldestTimestamp), limit: '200' });
      if (cursor) params.set('cursor', cursor);
      const response = await fetch(`https://api.slack.com/audit/v1/logs?${params}`, {
        headers: { 'Authorization': `Bearer ${orgToken}` }
      });
      if (response.status === 429) {
        const err = new Error('ratelimited');
        err.status = 429; err.headers = Object.fromEntries(response.headers.entries());
        throw err;
      }
      return response.json();
    });
    entries.push(...(res.entries || []));
    cursor = res.response_metadata?.next_cursor || '';
  } while (cursor);
  return entries;
}
```

---

## Compliance Validation

### Required Channel Audit

```javascript
async function auditRequiredChannels(client, requiredChannels) {
  const allChannels = await paginateSlackAPI(client, 'conversations.list', { types: 'public_channel' });
  const channelMap = new Map(allChannels.map(c => [c.name, c]));
  const findings = [];

  for (const required of requiredChannels) {
    const channel = channelMap.get(required.name);
    if (!channel) {
      findings.push({ channel: required.name, status: 'missing', severity: required.critical ? 'critical' : 'warning' });
      continue;
    }
    try {
      const members = await paginateSlackAPI(client, 'conversations.members', { channel: channel.id });
      const botInfo = await client.auth.test();
      const botIsMember = members.includes(botInfo.user_id);
      findings.push({ channel: required.name, channelId: channel.id, status: botIsMember ? 'compliant' : 'bot_not_member', archived: channel.is_archived, severity: !botIsMember ? 'warning' : 'ok' });
    } catch (err) {
      findings.push({ channel: required.name, status: 'access_denied', severity: 'warning', detail: err.message });
    }
  }
  return findings;
}
```

### Notification Routing Validation

```javascript
const ROUTING_RULES = [
  { pattern: /deploy|release/i,  channel: '#deployments',    critical: true  },
  { pattern: /incident|outage/i, channel: '#incidents',      critical: true  },
  { pattern: /security|vuln/i,   channel: '#security-alerts', critical: true },
];

function validateNotificationRouting(recentMessages) {
  const misrouted = [];
  for (const msg of recentMessages) {
    for (const rule of ROUTING_RULES) {
      if (rule.pattern.test(msg.text) && msg.channel !== rule.channel) {
        misrouted.push({ messageTs: msg.ts, text: msg.text.substring(0, 80), actualChannel: msg.channel, expectedChannel: rule.channel, critical: rule.critical });
      }
    }
  }
  return misrouted;
}
```

---

## Full Audit Orchestrator

```javascript
async function runFullSlackAudit(config) {
  const report = {
    timestamp: new Date().toISOString(), workspace: null,
    commands: { total: 0, dead: 0, findings: [] },
    scopes: { assessment: null, findings: [] },
    compliance: { channels: [], routing: [] },
    summary: { status: 'pass', criticalFindings: 0, warnings: 0 },
  };

  const identity = await config.client.auth.test();
  const team = await config.client.team.info();
  report.workspace = { id: team.team.id, name: team.team.name };

  const commands = await getRegisteredCommands(config.appToken);
  report.commands.total = commands.length;
  for (const cmd of commands) {
    const validation = await rateLimitedFetch(() => validateCommandEndpoint(cmd));
    if (validation.status !== 'healthy') { report.commands.findings.push(validation); report.commands.dead++; }
  }

  report.scopes = auditScopes(identity.scopes || []);
  report.compliance.channels = await auditRequiredChannels(config.client, config.requiredChannels || []);

  const critical = [
    ...report.commands.findings.filter(f => f.status === 'unreachable'),
    ...report.compliance.channels.filter(f => f.severity === 'critical'),
    ...(report.scopes.assessment === 'REVIEW_REQUIRED' ? [report.scopes] : []),
  ];
  const warnings = [
    ...report.commands.findings.filter(f => f.status === 'timeout'),
    ...report.compliance.channels.filter(f => f.severity === 'warning'),
  ];
  report.summary.criticalFindings = critical.length;
  report.summary.warnings = warnings.length;
  report.summary.status = critical.length > 0 ? 'fail' : warnings.length > 0 ? 'warn' : 'pass';
  return report;
}
```

---

## Anti-Patterns

1. **Ignoring rate limits during audits** — a single audit sweep can hit hundreds of endpoints. Always use `rateLimitedFetch` with exponential backoff and respect `Retry-After` headers.

2. **Persisting cursors across sessions** — Slack cursors expire within minutes. If an audit is interrupted, restart pagination from the beginning.

3. **Using legacy tokens for new audits** — deprecated legacy workspace tokens bundle excessive permissions. Use granular `xoxb-` bot tokens with only the scopes needed.

4. **Hardcoding scope lists** — maintain `REQUIRED_SCOPES` and `HIGH_RISK_SCOPES` as configuration, not constants. Slack adds new scopes regularly.

5. **Skipping SSL verification on probes** — invalid SSL is a finding, not an obstacle. Never disable certificate validation.

6. **Auditing without owner context** — apps owned by deactivated users are a top source of orphaned subscriptions. Always check `app.owner_id` against `users.info`.

7. **Treating Socket Mode and HTTP identically** — Socket Mode apps have no public request URL. Verify `apps.connections.open` sessions instead of endpoint reachability.

8. **Missing compliance baseline** — running audits without defined required channels, scopes, and routing rules produces noise without actionable findings. Define the baseline first.

---

## References

- [Implementing Slash Commands](https://docs.slack.dev/interactivity/implementing-slash-commands/)
- [Events API](https://docs.slack.dev/apis/events-api/)
- [Scopes Reference](https://docs.slack.dev/reference/scopes/)
- [Tokens](https://docs.slack.dev/authentication/tokens/)
- [Rate Limits](https://docs.slack.dev/apis/web-api/rate-limits/)
- [Pagination](https://docs.slack.dev/apis/web-api/pagination/)
- [App Manifests](https://docs.slack.dev/app-manifests/)
- [Audit Logs API](https://docs.slack.dev/reference/audit-logs-api/methods-actions-reference/)
