<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `express-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: express-patterns
version: "1.1.0"
updated: "2026-05-29"
category: developer
tags: [express, nodejs, middleware, routing, error-handling, rate-limiting, cors, security, rest-api, supertest, graceful-shutdown]
description: "Express.js production patterns expert — middleware, routing, error handling, rate limiting, CORS, security hardening, and testing for Express 4.x and 5.x. TRIGGER: user is writing or reviewing Express.js middleware (custom, async, composed, execution order); designing REST API routes (versioning, modular routers, nested resources, param validation); implementing centralized error handling; adding rate limiting (express-rate-limit, Redis store, sliding window); configuring CORS (credentials mode, dynamic allowlists, preflight caching); hardening an Express app (Helmet, trust proxy, input sanitization, request size limits); setting up production patterns (graceful shutdown, health checks, pino-http, request IDs, compression, ETags); writing tests with supertest; migrating from Express 4 to Express 5. SKIP: user is working with a different Node.js framework (Fastify, Koa, Hapi, NestJS) or asking a general JavaScript/Node.js question not specific to Express routing or middleware."
related_skills: [web-auth-patterns, pino-structured-logging, javascript-nodejs, testing-and-vitest-expert]
---

# Express.js Production Patterns

> **Sibling frameworks:** For Fastify, NestJS, and Hono (and Fastify-vs-NestJS-vs-Hono selection) — the post-Express Node.js/TypeScript backend frameworks — load `references/nodejs-backend-frameworks.md` in this same hub.

> **Cross-reference:** For auth-specific patterns used with Express (OAuth 2.1 PKCE flows, JWT refresh rotation, session fixation prevention, CSRF token middleware, cookie security attributes, CORS credentials configuration, CSP headers, passkeys/WebAuthn integration), see the `web-auth-patterns` skill.

## Prerequisites to check first

Check Express version before recommending async error patterns:

- Check `package.json` for the `"express"` version field.
- **Express 5.x (≥5.0.0)**: native async — no wrapper needed. Requires Node 18+. Released Oct 2024.
- **Express 4.x (<5.0.0)**: wrap every async route with `asyncHandler` or install `express-async-errors` at the top of the entry file.

Edge case: if the project uses `express-async-errors` AND is being upgraded to Express 5, remove the package — it double-catches errors on Express 5.

## Context reference

Load this file for concrete patterns, decision tables, and cited sources:

```
@~/.claude/skills/express-patterns/references/express-patterns-context.md
```

## Key decisions

### 1. Async error handling (Express 4 vs 5)

- **Express 4**: without `asyncHandler` or `express-async-errors`, a rejected promise silently hangs the request.
- **Express 5**: native promise rejection forwarding built in. Do NOT add `express-async-errors`.

### 2. Centralized error handler placement

The most common mistake is wrong registration order. The required order is:

```
routes → 404 notFoundHandler → errorHandler (LAST, exactly 4 params)
```

- Must have exactly 4 params: `(err, req, res, next)` — Express detects this by arity.
- Missing the 4th param (`next`) causes Express to treat it as normal middleware and skip it.
- Register AFTER all routes, AFTER notFoundHandler.
- Distinguish operational errors (`isOperational: true`) from programmer errors.
- Log programmer errors at `error`/`fatal`; operational at `warn`/`info`.
- Never leak stack traces or internal details in production responses.

### 3. Rate limiting

| Deployment | Store | Notes |
|-----------|-------|-------|
| Single process / dev | `MemoryStore` (default) | Fine |
| Multi-process / Kubernetes / PM2 cluster | `RedisStore` | Required — in-memory store gives each pod its own counter, making limits ineffective |

- Auth endpoints (login, password reset): separate stricter limiter, `skipSuccessfulRequests: true`.
- Authenticated routes: key by `req.user.id` not IP (shared office IPs, NAT).
- Use `express-slow-down` alongside hard limits for graduated throttling on public APIs.

### 4. CORS credentials trap

The single most common CORS misconfiguration:

- **Never** combine `credentials: true` with `origin: '*'` — browsers reject it silently.
- When `credentials: true`, enumerate exact origins in an array or validator function.
- Add `app.options('*', cors(corsOptions))` to handle preflight for all routes.
- Use `maxAge: 86400` to cache preflight for 24 hours.

### 5. Trust proxy

Required when Express sits behind Nginx, AWS ALB, Heroku, or any reverse proxy. Set **before** rate limiter middleware:

```js
app.set('trust proxy', 1); // trust one proxy hop
```

Without this, `req.ip` returns the proxy IP — rate limiting, geo-blocking, and audit logs all break.

### 6. Health check split

Kubernetes uses two separate probes. Register these **before** auth middleware:

| Endpoint | Purpose | Checks |
|----------|---------|--------|
| `/healthz` | Liveness: is the process alive? | Returns 200 with no dependency checks |
| `/readyz` | Readiness: can the process serve traffic? | DB ping + Redis ping; returns 503 if any dependency is down |

## Quick code references

### asyncHandler wrapper (Express 4)

```js
const asyncHandler = (fn) => (req, res, next) =>
  Promise.resolve(fn(req, res, next)).catch(next);

router.get('/users/:id', asyncHandler(async (req, res) => {
  const user = await db.findById(req.params.id);
  if (!user) throw new NotFoundError('User');
  res.json(user);
}));
```

### Custom error classes

```js
class AppError extends Error {
  constructor(message, statusCode, code) {
    super(message);
    this.statusCode = statusCode;
    this.code = code;              // machine-readable, e.g. 'USER_NOT_FOUND'
    this.isOperational = true;
    Error.captureStackTrace(this, this.constructor);
  }
}
class NotFoundError extends AppError {
  constructor(resource) { super(`${resource} not found`, 404, 'NOT_FOUND'); }
}
class ValidationError extends AppError {
  constructor(errors) {
    super('Validation failed', 422, 'VALIDATION_ERROR');
    this.errors = errors;
  }
}
```

### Centralized error handler + 404 handler

```js
// 404 — after all routes, before error handler
app.use((req, res, next) => {
  next(new NotFoundError(`${req.method} ${req.path}`));
});

// Error handler — LAST middleware, exactly 4 params
function errorHandler(err, req, res, next) {
  const statusCode = err.isOperational ? (err.statusCode || 500) : 500;
  res.status(statusCode).json({
    error: {
      code: err.code || 'INTERNAL_ERROR',
      message: err.isOperational ? err.message : 'An unexpected error occurred',
      ...(err.errors && { details: err.errors }),
      ...(process.env.NODE_ENV !== 'production' && { stack: err.stack }),
    }
  });
}
app.use(errorHandler); // LAST
```

### Rate limiter with Redis

```js
const { rateLimit } = require('express-rate-limit');
const { RedisStore } = require('rate-limit-redis');

const limiter = rateLimit({
  windowMs: 15 * 60 * 1000,
  limit: 100,
  standardHeaders: 'draft-7',
  legacyHeaders: false,
  store: new RedisStore({
    sendCommand: (...args) => redisClient.sendCommand(args),
    prefix: 'rl:',
  }),
  keyGenerator: (req) => req.user?.id ?? req.ip,
});
```

### pino-http with request ID

```js
// Request ID middleware MUST come before pinoHttp
app.use((req, res, next) => {
  req.id = req.headers['x-request-id'] || uuidv4();
  res.setHeader('X-Request-ID', req.id);
  next();
});
app.use(pinoHttp({
  logger,
  genReqId: (req) => req.id,
  autoLogging: { ignore: (req) => req.url === '/healthz' },
}));
// In route handlers: req.log.info('message') — child logger with reqId bound
```

### Graceful shutdown

```js
const server = app.listen(PORT);

process.on('SIGTERM', () => {
  server.close((err) => {
    if (err) { logger.error(err); process.exit(1); }
    db.close().then(() => process.exit(0));
  });
  setTimeout(() => process.exit(1), 30_000).unref(); // force-kill after 30s
});
```

### Supertest app isolation pattern

```js
// app.js — export app WITHOUT calling listen
const app = express();
// ... middleware and routes ...
module.exports = app;

// server.js — entry point only
const app = require('./app');
app.listen(process.env.PORT || 3000);

// users.test.js
const request = require('supertest');
const app = require('../app');

describe('GET /api/v1/users/:id', () => {
  it('returns 200 with user data', async () => {
    const res = await request(app)
      .get('/api/v1/users/abc-123')
      .set('Authorization', `Bearer ${testToken}`)
      .expect(200)
      .expect('Content-Type', /json/);
    expect(res.body.id).toBe('abc-123');
  });
  it('returns 401 without token', async () => {
    await request(app).get('/api/v1/users/abc-123').expect(401);
  });
  it('returns 404 for unknown user', async () => {
    const res = await request(app)
      .get('/api/v1/users/does-not-exist')
      .set('Authorization', `Bearer ${testToken}`)
      .expect(404);
    expect(res.body.error.code).toBe('NOT_FOUND');
  });
});
```

## Common edge cases and failure modes

| Symptom | Cause | Fix |
|---------|-------|-----|
| All errors silently ignored or cause 500 with no body | 4-param error handler placed before routes | Move `app.use(errorHandler)` to the very last `app.use()` call |
| Double error forwarding on Express 5 | `express-async-errors` still installed | Remove the package when on Express 5 |
| All clients share the same rate limit counter | Rate limiter without trust proxy set | `app.set('trust proxy', 1)` before limiter |
| Browser silently drops credentialed requests | `credentials: true` + `origin: '*'` in CORS | Use explicit origin array or validator function |
| `req.params.parentId` is undefined in child router | Nested router missing `mergeParams` | `express.Router({ mergeParams: true })` |
| Error handler ignored despite correct position | 3 params instead of 4 | Add `next` as 4th param even if unused |
| Kubernetes liveness probe fails, pod restart loop | Health check registered behind auth middleware | Register `/healthz` and `/readyz` before auth middleware |
| `req.body` undefined (Express 5) | Body-parser not registered | Add `app.use(express.json())` before routes; Express 5 no longer defaults to `{}` |
