<!-- hub-reference-banner -->
> **Reference file — part of the `devops-observability` hub.** Formerly the standalone `sentry-monitoring` skill.
> Sibling topics in this family are now reference files under the hubs (`devops-linux-internals`, `devops-linux-admin`, `devops-containers-cicd`, `devops-observability`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: sentry-monitoring
version: 1.1.0
updated: 2026-05-29
description: >
  Sentry error monitoring and performance tracing for Node.js — @sentry/node v8+
  SDK setup (CJS and ESM), captureException, spans (startSpan/startInactiveSpan),
  breadcrumbs, source maps (Vite/Webpack/CLI), Express and Fastify middleware,
  release tracking, sampling strategies (tracesSampleRate/tracesSampler),
  PII scrubbing (beforeSend), and alert configuration patterns.
  TRIGGER: adding Sentry to a Node.js application; debugging Sentry SDK
  initialization order; configuring tracing sample rates; setting up source maps
  for Node.js; integrating Sentry with Express or Fastify; scrubbing PII from
  Sentry events; configuring Sentry alerts for error rate or latency.
  SKIP: Sentry for browser JavaScript (different SDK — @sentry/browser);
  Sentry for Python, Ruby, or other non-Node runtimes; Datadog, New Relic, or
  other APM tools (unrelated products).
category: developer
tags: [sentry, monitoring, node.js, tracing, observability, express, fastify, error-tracking]
related_skills: [pino-structured-logging, nodejs-observability, express-patterns, cicd-pipelines]
whenToUse:
  - "Add Sentry error monitoring to a Node.js app"
  - "Sentry SDK must be initialized before other imports"
  - "Configure Sentry tracing sample rate for production"
  - "Set up source maps with Sentry Vite plugin"
  - "Sentry Express error handler setup"
  - "Scrub PII from Sentry events with beforeSend"
  - "Sentry alert rules for error rate spike"
  - "Sentry v8 startSpan vs startInactiveSpan"
  - "Sentry profiling setup with @sentry/profiling-node"
whenNotToUse:
  - "Sentry for browser JS (use @sentry/browser docs)"
  - "Datadog or New Relic instrumentation (different products)"
  - "Python, Ruby, or other non-Node Sentry SDKs"
---

# Sentry Monitoring for Node.js

## Overview

Sentry is an application monitoring platform that captures errors, traces performance, and tracks releases in production. The `@sentry/node` SDK (v8+) uses OpenTelemetry under the hood; initialization **must happen before any other require/import** so auto-instrumentation can wrap modules at load time.

## When Not to Use

- `NODE_ENV=test` or CI pipelines — set `enabled: false` in init to suppress noise.
- Local dev hot reloads that spam Sentry quota — use `tracesSampleRate: 0` or `enabled: false`.
- Already using a full-featured APM (Datadog, New Relic) with SDK-level Node.js tracing — adding Sentry on top creates double instrumentation overhead.

## Quick Reference

| Goal | API / Config key |
|---|---|
| Capture handled error | `Sentry.captureException(err)` |
| Capture message | `Sentry.captureMessage("msg", "warning")` |
| Measure a block | `Sentry.startSpan({ name, op }, callback)` |
| Add breadcrumb | `Sentry.addBreadcrumb({ category, message, level })` |
| Set user context | `Sentry.setUser({ id, email })` |
| Set tag | `Sentry.setTag("key", "value")` |
| Request-scoped context | `Sentry.withScope(scope => { ... })` |
| Scrub before send | `beforeSend` / `beforeSendTransaction` callbacks |
| Control sampling | `tracesSampleRate` or `tracesSampler` function |
| Profile CPU | `profileSessionSampleRate` (requires tracing enabled) |

---

## 1. SDK Setup and Initialization

Install: `npm install @sentry/node`

**Critical:** create `instrument.js` (or `instrument.ts`) and load it **first**, before any other module. The file format (CJS vs ESM) must match your app's module system.

**CJS `instrument.js`:**
```js
// instrument.js (CommonJS) — require() this as the very first line of your entrypoint
const Sentry = require("@sentry/node");

Sentry.init({
  dsn: process.env.SENTRY_DSN,
  release: process.env.npm_package_version,     // e.g. "1.4.2"
  environment: process.env.NODE_ENV ?? "production",

  tracesSampleRate: 1.0,  // lower in production (see §8)

  // debug: true,          // uncomment to log SDK activity during setup
});
```

**CJS entrypoint:**
```js
// index.js
require("./instrument");      // MUST be the first require()
const express = require("express");
// ...
```

**ESM `instrument.mjs` (Node ≥ 18.19):**
```js
// instrument.mjs
import * as Sentry from "@sentry/node";

Sentry.init({
  dsn: process.env.SENTRY_DSN,
  release: process.env.npm_package_version,
  environment: process.env.NODE_ENV ?? "production",
  tracesSampleRate: 1.0,
});
```

**ESM entrypoint — use `--import` flag, not `import` inside code:**
```bash
node --import ./instrument.mjs src/index.js
```

**Adding profiling** (separate optional package — extend your existing init, do not call init twice):
```bash
npm install @sentry/profiling-node
```
```js
// instrument.js — merge these fields into your existing Sentry.init() call:
const { nodeProfilingIntegration } = require("@sentry/profiling-node");
Sentry.init({
  dsn: process.env.SENTRY_DSN,
  // ...your other config...
  integrations: [nodeProfilingIntegration()],   // add this
  profileSessionSampleRate: 0.1,                // add this
  profileLifecycle: "trace",                    // add this
});
```

Sources: [Sentry Node.js docs](https://docs.sentry.io/platforms/javascript/guides/node/), [Late Initialization (ESM/CJS)](https://docs.sentry.io/platforms/javascript/guides/node/install/late-initialization/)

---

## 2. Error Capture

```js
// Unhandled exceptions + promise rejections: captured automatically.

// Handled errors — in try/catch:
try {
  await riskyOperation();
} catch (err) {
  Sentry.captureException(err, {
    tags: { component: "payment-processor" },
    extra: { orderId },
  });
  throw err;   // re-throw if caller needs it
}

// Informational message with severity level:
// levels: "fatal" | "error" | "warning" | "log" | "info" | "debug"
Sentry.captureMessage("Quota approaching 80%", "warning");

// Attach user context (call once per request / session):
Sentry.setUser({ id: user.id, email: user.email });
Sentry.setTag("tenant", req.hostname);
```

Sources: [Capturing Errors](https://docs.sentry.io/platforms/javascript/guides/node/usage/), [Capturing Errors (backend)](https://docs.sentry.io/product/sentry-basics/integrate-backend/capturing-errors/)

---

## 3. Performance Tracing — Spans

In SDK v8, transactions are gone; everything is a **span**. Use `startSpan` for the active scope pattern.

```js
// Active span — ends automatically when callback resolves
const result = await Sentry.startSpan(
  { name: "process-order", op: "function" },
  async () => {
    // nested child span
    await Sentry.startSpan(
      { name: "db.query", op: "db" },
      () => db.query("SELECT …"),
    );
    return processOrder(data);
  },
);

// Manual span (for fire-and-forget or streaming):
const span = Sentry.startInactiveSpan({ name: "upload", op: "http.client" });
try {
  await uploadFile();
} finally {
  span.end();
}
```

**Auto-instrumented** (no code needed when `tracesSampleRate` > 0):
- HTTP (http/https), Express/Fastify routes, pg, mysql2, mongodb, redis, graphql, Prisma, AWS SDK

Sources: [Set Up Tracing](https://docs.sentry.io/platforms/javascript/guides/node/tracing/), [Auto-Instrumentation](https://docs.sentry.io/platforms/javascript/guides/node/tracing/instrumentation/automatic-instrumentation/)

---

## 4. Breadcrumbs

Breadcrumbs are a trail of events leading up to an error. Many are collected automatically (HTTP calls, DB queries, console output). Add custom ones for business logic events.

```js
// Custom breadcrumb
Sentry.addBreadcrumb({
  category: "auth",
  message: `User ${user.email} authenticated`,
  level: "info",           // fatal | error | warning | log | info | debug
  data: { method: "oauth" },
});

// Filter or mutate automatic breadcrumbs — must be inside Sentry.init():
Sentry.init({
  // ...other config...
  beforeBreadcrumb(breadcrumb, hint) {
    // Drop noisy health-check requests
    if (breadcrumb.category === "http" &&
        breadcrumb.data?.url?.includes("/healthz")) {
      return null;  // null = discard; returning breadcrumb = keep
    }
    return breadcrumb;
  },
});
```

Sources: [Breadcrumbs](https://docs.sentry.io/platforms/javascript/guides/node/enriching-events/breadcrumbs/), [Sentry Blog — Node Breadcrumbs](https://blog.sentry.io/node-breadcrumbs/)

---

## 5. Source Maps

Source maps let Sentry display original TypeScript/transpiled source in stack traces.

**Vite (recommended):**
```bash
npm install --save-dev @sentry/vite-plugin
```
```js
// vite.config.js
import { sentryVitePlugin } from "@sentry/vite-plugin";
export default {
  build: { sourcemap: true },
  plugins: [
    // Sentry plugin must come LAST in the plugins array
    sentryVitePlugin({
      authToken: process.env.SENTRY_AUTH_TOKEN,
      org: "my-org",
      project: "my-node-api",
    }),
  ],
};
```

**Webpack:**
```bash
npm install --save-dev @sentry/webpack-plugin
```
```js
// webpack.config.js
const { sentryWebpackPlugin } = require("@sentry/webpack-plugin");
module.exports = {
  devtool: "source-map",
  plugins: [
    sentryWebpackPlugin({
      authToken: process.env.SENTRY_AUTH_TOKEN,
      org: "my-org",
      project: "my-node-api",
    }),
  ],
};
```

**CLI fallback (esbuild, tsc, custom):**
```bash
npx @sentry/cli sourcemaps upload --auth-token $SENTRY_AUTH_TOKEN \
  --org my-org --project my-node-api ./dist
```

Sources: [Source Maps](https://docs.sentry.io/platforms/javascript/guides/node/sourcemaps/), [Vite plugin](https://docs.sentry.io/platforms/javascript/guides/node/sourcemaps/uploading/vite/)

---

## 6. Express and Fastify Integration

### Express

In SDK v8, `setupExpressErrorHandler(app)` replaces the old `Handlers.requestHandler()` / `Handlers.errorHandler()` middleware pair. Call it after all routes are registered.

```js
// instrument.js MUST be required first (see §1)
const Sentry = require("@sentry/node");
const express = require("express");

const app = express();

// --- define your routes here ---
app.get("/", (req, res) => res.send("ok"));

// Register Sentry error handler AFTER all routes, BEFORE your own error middleware.
// This automatically captures request data AND errors (replaces v7 requestHandler + errorHandler).
Sentry.setupExpressErrorHandler(app, {
  shouldHandleError(err) {
    return err.status >= 400;    // default is >= 500
  },
});

// Your own error middleware goes after Sentry's:
app.use((err, req, res, next) => {
  res.status(err.status ?? 500).json({ error: err.message });
});
```

### Fastify

```js
import Fastify from "fastify";
import * as Sentry from "@sentry/node";
import { setupFastifyErrorHandler } from "@sentry/node";

const app = Fastify();
setupFastifyErrorHandler(app);   // captures 5xx by default (also ≤2xx)

// routes as normal…
```

### Request-scoped context (async isolation)

In SDK v8, Sentry uses AsyncLocalStorage (via OpenTelemetry) to automatically propagate scope per request — you don't need `withScope` for basic use. Use `withScope` only when you need a *temporary, isolated* scope for a specific operation within a request (e.g., capturing an error with extra context without polluting the request scope):

```js
// ✅ Per-request user tagging — set directly on the current scope:
app.use((req, res, next) => {
  Sentry.setUser({ id: req.user?.id });
  Sentry.setTag("route", req.path);
  next();   // SDK's AsyncLocalStorage keeps this isolated per-request automatically
});

// ✅ Temporary isolated scope for a specific operation within a request:
async function processPayment(orderId) {
  await Sentry.withScope(async (scope) => {
    scope.setTag("orderId", orderId);
    scope.setLevel("warning");
    try {
      await charge();
    } catch (err) {
      Sentry.captureException(err);  // captured with orderId tag, won't affect other spans
    }
  });
}
```

> Note: Do not call `next()` inside a `withScope` callback in Express middleware — the scope ends when the synchronous callback returns, so async downstream handlers run outside it.

Sources: [Express](https://docs.sentry.io/platforms/javascript/guides/express/), [Fastify Error Handler](https://docs.sentry.io/platforms/javascript/guides/fastify/features/error-handler/)

---

## 7. Release Tracking and Deploy Notifications

Set `release` in `Sentry.init` (see §1). Then use sentry-cli in CI/CD:

```bash
# Install
npm install --save-dev @sentry/cli

# Full release workflow
export SENTRY_AUTH_TOKEN=...
export SENTRY_ORG=my-org
export SENTRY_PROJECT=my-node-api
VERSION=$(npm pkg get version --workspaces=false | tr -d '"')

sentry-cli releases new $VERSION
sentry-cli releases set-commits $VERSION --auto
sentry-cli releases finalize $VERSION
sentry-cli releases deploys $VERSION new --env production
```

> Tip: store `SENTRY_ORG`, `SENTRY_PROJECT`, and `SENTRY_URL` in `.sentryclirc` at the project root instead of exporting them in every script.

**GitHub Actions shortcut:**
```yaml
- uses: getsentry/action-release@v1
  env:
    SENTRY_AUTH_TOKEN: ${{ secrets.SENTRY_AUTH_TOKEN }}
    SENTRY_ORG: my-org
    SENTRY_PROJECT: my-node-api
  with:
    environment: production
```

Sources: [Releases & Health](https://docs.sentry.io/platforms/javascript/guides/node/configuration/releases/), [Releases Setup](https://docs.sentry.io/product/releases/setup/)

---

## 8. Sampling Strategies

```js
Sentry.init({
  // Simple uniform rate (0–1). Start at 1.0 for dev.
  tracesSampleRate: 0.1,   // sample 10% in production

  // OR — dynamic sampler (takes precedence over tracesSampleRate when set):
  tracesSampler(samplingContext) {
    const { name, parentSampled } = samplingContext;

    // Respect upstream sampling decision in distributed traces
    if (parentSampled !== undefined) return parentSampled;

    // Always trace critical paths
    if (name?.includes("/checkout")) return 1.0;

    // Drop health checks
    if (name?.includes("/healthz")) return 0;

    return 0.05;   // default 5%
  },

  // Profiling (must be ≤ tracesSampleRate)
  profileSessionSampleRate: 0.1,
  profileLifecycle: "trace",
});
```

**Rule of thumb:** production backend → `tracesSampleRate: 0.05–0.1`; use `tracesSampler` when you need per-route control.

Sources: [Sampling](https://docs.sentry.io/platforms/node/configuration/sampling/), [Sentry Blog — Sampling Strategy](https://blog.sentry.io/sampling-strategy-sentry/), [Profiling](https://docs.sentry.io/platforms/javascript/guides/node/profiling/)

---

## 9. Privacy and Data Scrubbing

**Goal:** ensure PII never leaves your infrastructure.

```js
Sentry.init({
  // 1. beforeSend — mutate or drop error events
  //    Return null (not undefined) to discard an event entirely.
  beforeSend(event, hint) {
    // Strip auth headers
    if (event.request?.headers) {
      delete event.request.headers["authorization"];
      delete event.request.headers["cookie"];
    }
    // Drop events from bots
    if (event.request?.headers?.["user-agent"]?.includes("bot")) return null;
    return event;
  },

  // 2. beforeSendTransaction — same for perf events
  beforeSendTransaction(event) {
    // Remove query strings that may carry tokens
    if (event.request?.url) {
      event.request.url = event.request.url.split("?")[0];
    }
    return event;
  },

  // 3. denyUrls — never send errors whose top stack frame matches these patterns
  denyUrls: [/node_modules/, /^chrome-extension:\/\//],

  // 4. allowUrls — only send errors from these (useful for library-heavy apps)
  // allowUrls: [/https:\/\/yourapp\.com/],

  // 5. sendDefaultPii: false (default) — disable IP + user-agent auto-capture
  sendDefaultPii: false,
});
```

**Server-side scrubbing** (no redeploy needed): Sentry UI → Project Settings → Security & Privacy → Data Scrubber. Use for last-mile cleanup of fields you can't easily catch client-side.

Sources: [Scrubbing Sensitive Data](https://docs.sentry.io/platforms/javascript/guides/node/data-management/sensitive-data/), [Options Reference](https://docs.sentry.io/platforms/javascript/configuration/options/)

---

## 10. Alert Configuration Patterns

Alerts live in the Sentry UI (Project → Alerts → Create Alert Rule) or via the Sentry API. Two types:

| Type | Trigger |
|---|---|
| **Issue alert** | New event matches conditions (first seen, regression, frequency) |
| **Metric alert** | Aggregate metric crosses threshold (error rate, p95 latency, Apdex) |

**Recommended alert rules for a Node.js service:**

1. **First-seen unhandled error** — Issue alert: "A new issue is created" → notify Slack `#alerts`.
2. **Error rate spike** — Metric alert: `error_count()` > 100/min for 5 min → PagerDuty.
3. **High p95 latency** — Metric alert: `p95(transaction.duration)` > 2000ms → Slack.
   *(Note: `transaction.duration` is the Sentry UI metric name — this is separate from the SDK v8 spans API.)*
4. **Crash rate** — Metric alert: `failure_rate()` > 0.05 over 10 min → page on-call.

**Mute noisy issues with fingerprinting:**
```js
Sentry.init({
  beforeSend(event) {
    // Group all timeout errors under one fingerprint
    if (event.exception?.values?.[0]?.type === "TimeoutError") {
      event.fingerprint = ["timeout-error", "{{ default }}"];
    }
    return event;
  },
});
```

Sources: [Alerts](https://docs.sentry.io/product/alerts/), [Metric Alert Configuration](https://docs.sentry.io/product/alerts/create-alerts/metric-alert-config/)

---

## Common Mistakes

| Mistake | Fix |
|---|---|
| Importing app code before `instrument.js` | `require("./instrument")` must be the first line of the entrypoint |
| Using ESM `import` in `instrument.js` then `require()`-ing it | ESM and CJS instrument files are incompatible — match the file to your app's module system |
| `tracesSampleRate: 1.0` left in production | Lower to 0.05–0.1 or use `tracesSampler` |
| Forgetting `release` in `Sentry.init` | Source maps and commit tracking won't resolve without it |
| `beforeSend` returns `undefined` instead of `null` | Return `null` (not `undefined`) to drop an event |
| Using `Sentry.expressErrorHandler()` (does not exist in v8) | Use `Sentry.setupExpressErrorHandler(app)` |
| Placing Express error handler before routes | `setupExpressErrorHandler(app)` must be called after all routes are defined |
| Using `startTransaction` (v7 API) in v8 SDK | Use `startSpan` / `startInactiveSpan` |
| Calling `next()` inside a `withScope` callback in Express middleware | Scope ends when the synchronous callback returns — set user/tags directly; SDK v8 AsyncLocalStorage isolates them per-request automatically |
| Source map plugin placed before bundler plugins | Sentry Vite/Webpack plugin must be last in plugins array |
| `sendDefaultPii: true` without scrubbing | IPs and user agents flow to Sentry — add `beforeSend` guards |
