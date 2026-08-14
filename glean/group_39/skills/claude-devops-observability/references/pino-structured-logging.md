<!-- hub-reference-banner -->
> **Reference file — part of the `devops-observability` hub.** Formerly the standalone `pino-structured-logging` skill.
> Sibling topics in this family are now reference files under the hubs (`devops-linux-internals`, `devops-linux-admin`, `devops-containers-cicd`, `devops-observability`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: pino-structured-logging
description: >
  Pino v9+ structured logging for Node.js: core API, child loggers, serializers, redaction,
  transports (worker-thread architecture, pino-pretty, pino/file, multiple targets, pipeline,
  custom transport), performance optimization, Express/Fastify/NestJS integration with pino-http,
  AsyncLocalStorage for automatic request context, OpenTelemetry trace correlation (instrumentation
  approach and transport approach), TypeScript types, browser logging, and production configuration.
  TRIGGER: user is implementing or reviewing structured logging in a Node.js application using
  Pino; configuring pino-pretty, pino-http, child loggers, redaction, serializers, or OTel
  integration with Pino; reviewing logging code in a Node.js project that imports pino.
  SKIP: Winston, Bunyan, or other non-Pino loggers; generic log aggregation infrastructure
  (Elasticsearch, Loki ingestion config) with no Pino-specific question; Node.js observability
  beyond logging (use nodejs-observability for tracing/metrics/health checks).
version: 1.2.0
updated: "2026-05-29"
category: developer
tags:
  - pino
  - logging
  - structured-logging
  - nodejs
  - observability
  - json
  - typescript
  - opentelemetry
keywords:
  - pino
  - pino v9
  - structured logging
  - child logger
  - serializer
  - redaction
  - pino-pretty
  - pino-http
  - pino transport
  - worker thread transport
  - pino/file
  - pino-roll
  - pino-elasticsearch
  - pino-loki
  - pino-opentelemetry-transport
  - AsyncLocalStorage
  - correlation ID
  - request context
  - trace context
  - OTel instrumentation-pino
  - NestJS pino
  - Fastify logger
  - TypeScript pino
  - browser pino
  - log level
  - NDJSON
whenToUse:
  - Setting up or configuring Pino logging in a Node.js application
  - Implementing child loggers or correlation ID patterns with pino-http
  - Configuring pino-pretty for development or pino/file for production
  - Adding log redaction or custom serializers to sanitize PII
  - Integrating Pino with OpenTelemetry for trace context injection
  - Reviewing logging code in a Node.js project that uses Pino
  - Setting up multi-target transports (stdout + file + remote)
  - Configuring Pino in Express, Fastify, or NestJS
whenNotToUse:
  - Winston, Bunyan, Morgan, or other non-Pino loggers — advise migration or handle inline
  - Generic log aggregation infrastructure (Elasticsearch index config, Loki pipeline setup) with no Pino question
  - Node.js observability beyond logging (tracing, metrics, health checks) — use nodejs-observability
related_skills:
  - nodejs-observability
  - backend-patterns
  - javascript-nodejs
  - express-patterns
globs:
  - "**/logger.{js,ts,mjs}"
  - "**/logging.{js,ts,mjs}"
  - "**/pino*.{js,ts,mjs}"
  - "**/log-config*.{js,ts,mjs}"
alwaysApply: false
---

# Pino Structured Logging

Pino is a high-performance, structured JSON logger for Node.js. It outputs NDJSON to stdout and
offloads all formatting, transport, and I/O to worker threads, keeping the main event loop free.
**5–8x faster than Winston and Bunyan** in throughput benchmarks.

**Version:** Pino v9.x (current stable as of 2026). Transport API stable since v7.

```bash
npm install pino                    # v9.x
npm install -D pino-pretty          # dev-only pretty printer
npm install pino-http               # Express/Koa HTTP logging middleware
```

## Key Design Principles

- **Minimal main-thread work** — serialize to JSON string, write to stdout, return immediately.
- **Worker-thread transports** — pretty-printing, file I/O, and remote shipping run off the main thread via `pino.transport()`.
- **Child logger context** — bind request IDs once; every subsequent log includes them automatically. Child creation uses prototype inheritance — no data copying.
- **Level guard before expensive serialization** — use `logger.isLevelEnabled('debug')` before computing expensive objects.

## Log Levels

| Level  | Number | When to use |
|--------|--------|-------------|
| trace  | 10     | Fine-grained debug during development |
| debug  | 20     | Operational debug, variable state |
| info   | 30     | Normal operational events |
| warn   | 40     | Recoverable issues, degraded state |
| error  | 50     | Errors requiring attention |
| fatal  | 60     | Process-killing failures |

Set `level` to the lowest target needed. If a transport target uses `debug` but the root logger is `info`, debug messages are filtered before reaching the transport.

## Decision Guide

| Need | Solution |
|---|---|
| Human-readable output in development | `pino-pretty` transport (development only) |
| Per-request context without parameter drilling | `logger.child()` + AsyncLocalStorage |
| Sanitize PII before logging | `redact` paths or custom serializer |
| Write to multiple destinations | Multi-target transport with per-target level filters |
| Inject OTel trace IDs into logs | `@opentelemetry/instrumentation-pino` (recommended) or `pino-opentelemetry-transport` |
| Log rotation | `pino-roll` transport |
| Ship to Elasticsearch/Loki/Datadog | Community transport packages |
| Express HTTP request logging | `pino-http` middleware |

## Essential Patterns

### Basic Usage

```js
const pino = require('pino');
const logger = pino({ level: 'info' });

logger.info({ port: 3000 }, 'server listening');
// => {"level":30,"time":...,"port":3000,"msg":"server listening"}

try { dangerousOp(); } catch (err) {
  logger.error({ err }, 'operation failed');
}

// Guard expensive serialization
if (logger.isLevelEnabled('debug')) {
  logger.debug({ stats: computeExpensiveStats() }, 'tick');
}
```

### Child Logger Pattern

```js
// Per-request context via child (prototype-based, zero copy)
const reqLogger = logger.child({ requestId: 'abc-123', userId: 456 });
reqLogger.info('processing request');
// requestId and userId appear on every subsequent log line
```

### pino-http for Express (recommended)

```js
app.use(pinoHttp({
  logger,
  genReqId: (req, res) => {
    const id = req.headers['x-request-id'] || randomUUID();
    res.setHeader('x-request-id', id);
    return id;
  },
  customLogLevel: (req, res, err) => {
    if (res.statusCode >= 500 || err) return 'error';
    if (res.statusCode >= 400) return 'warn';
    return 'info';
  },
  autoLogging: { ignore: (req) => req.url === '/health' }
}));

app.get('/orders/:id', (req, res) => {
  req.log.info({ orderId: req.params.id }, 'fetching order'); // requestId auto-attached
});
```

### Production Configuration (minimal copy-paste)

```js
const logger = pino({
  level: process.env.LOG_LEVEL || 'info',
  base: { service: process.env.SERVICE_NAME || 'unknown', env: process.env.NODE_ENV },
  timestamp: pino.stdTimeFunctions.isoTime,
  formatters: { level(label) { return { level: label }; } },
  redact: {
    paths: ['password', 'token', 'accessToken', 'headers.authorization', 'headers.cookie', '*.apiKey'],
    censor: '[REDACTED]'
  },
  serializers: { err: pino.stdSerializers.err, req: pino.stdSerializers.req, res: pino.stdSerializers.res },
  transport: process.env.NODE_ENV === 'development'
    ? { target: 'pino-pretty', options: { colorize: true, translateTime: 'SYS:standard' } }
    : undefined
});
```

### OTel Trace Context Injection (recommended approach)

```js
const { NodeSDK } = require('@opentelemetry/sdk-node');
const { PinoInstrumentation } = require('@opentelemetry/instrumentation-pino');

// Must start SDK BEFORE requiring pino
const sdk = new NodeSDK({
  instrumentations: [new PinoInstrumentation()]
});
sdk.start();

// Every log line now automatically includes trace_id + span_id
const logger = require('pino')();
logger.info('this line carries trace context');
```

## Anti-Patterns

1. **`pino-pretty` in production** — 10x slower; produces non-parseable output. Development only.
2. **Logging inside serializers** — calling `logger.info()` from within a serializer causes infinite recursion.
3. **New logger per request** — use `logger.child()` instead; it's prototype-based and negligible cost.
4. **Skipping `flushSync()` on shutdown** — async destinations buffer writes; last log lines (including fatal errors) are lost without an explicit flush before exit.
5. **Root level higher than target level** — if root is `info` but a file transport target is `debug`, debug messages are filtered before reaching the transport.
6. **Logging raw `req`/`res` without serializers** — dumps circular references, auth headers, and body payloads into logs.
7. **String concatenation in messages** — `logger.info('User ' + id + ' logged in')` loses structured fields; use `logger.info({ userId: id }, 'user logged in')`.
8. **Leaking secrets via child bindings** — `logger.child({ token: req.headers.authorization })` permanently attaches the token to every log line.
9. **`PinoInstrumentation` registered after `require('pino')`** — OTel must patch pino before it is imported; register SDK first.

## Troubleshooting

| Symptom | Cause | Fix |
|---|---|---|
| Logs not appearing | Level set too high | Check `logger.level`; ensure root level ≤ target level |
| `[Redacted]` in unexpected fields | Wildcard redact too broad | Narrow redact paths; test with `logger.info({...})` |
| Transport module not found | Relative path without `./` | Use absolute path or `require.resolve()` |
| Logs lost on crash | Async destination not flushed | Add `handler.flushSync()` in `uncaughtException` + `beforeExit` handlers |
| Circular reference error | Raw req/res without serializer | Add `pino.stdSerializers.req` / `.res` |
| Child logger missing parent fields | Bindings mutated after child creation | Treat bindings as immutable after `.child()` call |
| OTel `trace_id` missing from logs | Instrumentation registered after pino import | Register `PinoInstrumentation` SDK before `require('pino')` |

## Full Reference

For complete code examples — all transport configurations, multiple targets, pipeline, custom transport authoring, async destination, log rotation, AsyncLocalStorage context propagation, OTel transport approach, TypeScript types, browser logging, environment-aware factory, NestJS/Fastify integration — read `references/pino-context.md`.
