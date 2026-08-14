<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `ops-registry-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: ops-registry-patterns
description: Operations registry patterns — dispatch tables, retry with backoff/jitter, circuit breakers, auto-remediation, idempotency, saga compensation, and audit logging. Use when designing operation dispatch systems with reliability guarantees.
version: 1.1.0
category: developer
tags: [operations, registry, retry, remediation, circuit-breaker, saga, idempotency, concurrency, audit-logging]
aliases: [ops-reg, retry-patterns, operation-dispatch]
triggers:
  - User is designing a named operation dispatch system or command registry
  - User needs retry logic with backoff, jitter, or circuit breakers
  - User asks about idempotency keys or safe retry handling
  - User is implementing saga/compensation for multi-step workflows
  - User needs auto-remediation (detect-classify-fix-verify) patterns
  - User asks about concurrency control, semaphores, or rate limiting for async operations
  - User needs structured audit logging for operation state transitions
  - User mentions "operation queue", "dispatch table", or "handler registry"
---

# Operations Registry Patterns

## Instructions

When this skill is activated:
1. Identify which specific pattern(s) the user needs from the sections below.
2. Provide the relevant code example(s) adapted to their use case.
3. Highlight the Anti-Patterns section if the user's approach risks a known pitfall.
4. Recommend the Checklist section for greenfield implementations.
5. Combine patterns when appropriate (e.g., registry + retry + circuit breaker + audit logging is the standard composition).

## Overview

An operations registry is a centralized dispatch table that maps named operations to handler functions, wrapping each invocation with cross-cutting concerns: retry logic, circuit breaking, idempotency enforcement, concurrency control, and audit logging. Instead of scattering reliability logic across callers, the registry owns the full operation lifecycle (pending, running, success, failed, retry, remediated) and enforces consistent policies for every registered operation.

This skill covers the complete pattern set for building reliable operation dispatch systems in JavaScript/TypeScript.

## Core Concepts

### Operation Shape

Every operation flowing through the registry carries a standard envelope:

```js
const operation = {
  id: crypto.randomUUID(),          // unique execution ID
  idempotencyKey: null,             // optional — deduplicate retries
  name: 'deployService',           // registry key
  params: { serviceId: 'svc-42' }, // handler-specific payload
  status: 'pending',               // pending | running | success | failed | retry | remediated
  attempt: 0,                      // current attempt number
  maxRetries: 3,                   // per-operation override or registry default
  createdAt: Date.now(),
  updatedAt: Date.now(),
  result: null,                    // success payload
  error: null,                     // last error
  history: [],                     // audit trail of state transitions
};
```

### Operation Lifecycle

```
pending ──> running ──> success
                │
                ▼
             failed ──> retry ──> running ...
                │
                ▼
           remediated (auto-fix applied, re-queued or closed)
                │
                ▼
            terminal  (permanent failure, no further retries)
```

Each transition emits an audit event. The registry enforces that only valid transitions occur (e.g., you cannot move from `success` back to `running`).

## Registry & Dispatch

The registry is a dispatch table: a plain object or `Map` keyed by operation name, with each value being a handler descriptor.

```js
class OperationRegistry {
  #handlers = new Map();
  #defaults = { maxRetries: 3, backoff: 'exponential', timeoutMs: 30_000 };

  /**
   * Register a named operation handler.
   * @param {string} name        - unique operation name
   * @param {Function} handler   - async (params, context) => result
   * @param {object} [opts]      - per-operation overrides
   */
  register(name, handler, opts = {}) {
    if (this.#handlers.has(name)) {
      throw new Error(`Operation "${name}" is already registered`);
    }
    this.#handlers.set(name, {
      handler,
      opts: { ...this.#defaults, ...opts },
    });
  }

  /**
   * Dispatch an operation by name.
   * Returns the full operation envelope after execution.
   */
  async dispatch(name, params, { idempotencyKey } = {}) {
    const entry = this.#handlers.get(name);
    if (!entry) throw new Error(`Unknown operation: "${name}"`);

    const op = createOperation(name, params, entry.opts, idempotencyKey);
    return this.#execute(op, entry);
  }

  async #execute(op, entry) {
    const { handler, opts } = entry;

    while (op.attempt <= opts.maxRetries) {
      op.attempt++;
      transition(op, 'running');

      try {
        op.result = await withTimeout(handler(op.params, op), opts.timeoutMs);
        transition(op, 'success');
        return op;
      } catch (err) {
        op.error = serializeError(err);

        if (isTerminal(err) || op.attempt > opts.maxRetries) {
          transition(op, 'failed');
          return op;
        }

        transition(op, 'retry');
        const delay = computeDelay(op.attempt, opts);
        await sleep(delay);
      }
    }

    transition(op, 'failed');
    return op;
  }
}
```

### Handler Registration Patterns

Register handlers at startup, not lazily. This makes the operation surface area discoverable and prevents runtime registration races.

```js
// Feature modules export their handlers
// file: ops/deploy.js
export const deployHandler = async (params, op) => {
  const result = await deployService(params.serviceId);
  return { deploymentId: result.id, timestamp: Date.now() };
};

export const deployOpts = {
  maxRetries: 5,
  backoff: 'exponential',
  timeoutMs: 120_000,
  retriableErrors: ['ETIMEDOUT', 'ECONNRESET', 'SERVICE_UNAVAILABLE'],
};

// file: registry-init.js
import { registry } from './registry.js';
import { deployHandler, deployOpts } from './ops/deploy.js';
import { rollbackHandler, rollbackOpts } from './ops/rollback.js';

registry.register('deployService', deployHandler, deployOpts);
registry.register('rollbackService', rollbackHandler, rollbackOpts);
```

## Retry Strategies

### Exponential Backoff with Jitter

Pure exponential backoff across many clients synchronizes into a thundering herd. Always add jitter.

```js
/**
 * Compute retry delay with exponential backoff + full jitter.
 * @param {number} attempt   - current attempt (1-based)
 * @param {object} opts      - { baseDelayMs, maxDelayMs }
 * @returns {number}         - delay in milliseconds
 */
function computeDelay(attempt, opts = {}) {
  const { baseDelayMs = 200, maxDelayMs = 30_000 } = opts;
  const exponential = baseDelayMs * Math.pow(2, attempt - 1);
  const capped = Math.min(exponential, maxDelayMs);
  // Full jitter: uniform random in [0, capped]
  return Math.floor(Math.random() * capped);
}
```

**Jitter strategies compared:**

| Strategy | Formula | When to use |
|---|---|---|
| Full jitter | `random(0, cap)` | Default choice; best spread |
| Equal jitter | `cap/2 + random(0, cap/2)` | When you need a guaranteed minimum wait |
| Decorrelated | `random(baseDelay, prevDelay * 3)` | High-contention scenarios |

### Retry Budget

A retry budget caps the total retry traffic as a percentage of successful traffic, preventing retry storms system-wide.

```js
class RetryBudget {
  #window = [];
  #windowMs;
  #maxRetryRatio;

  constructor({ windowMs = 60_000, maxRetryRatio = 0.1 } = {}) {
    this.#windowMs = windowMs;
    this.#maxRetryRatio = maxRetryRatio;
  }

  record(type) {
    this.#window.push({ type, ts: Date.now() });
    this.#prune();
  }

  canRetry() {
    this.#prune();
    const successes = this.#window.filter(e => e.type === 'success').length;
    const retries = this.#window.filter(e => e.type === 'retry').length;
    if (successes === 0) return retries < 3; // bootstrap: allow a few retries
    return retries / successes < this.#maxRetryRatio;
  }

  #prune() {
    const cutoff = Date.now() - this.#windowMs;
    this.#window = this.#window.filter(e => e.ts > cutoff);
  }
}
```

## Circuit Breaker

The circuit breaker prevents repeated calls to a failing downstream by "opening" after a failure threshold. It operates in three states: **closed** (normal), **open** (all calls fail-fast), and **half-open** (one probe call allowed to test recovery).

```js
class CircuitBreaker {
  #state = 'closed';       // closed | open | half-open
  #failureCount = 0;
  #successCount = 0;
  #lastFailureTime = 0;
  #opts;

  constructor(opts = {}) {
    this.#opts = {
      failureThreshold: 5,
      successThreshold: 2,     // half-open successes needed to close
      resetTimeoutMs: 30_000,  // time before open -> half-open
      ...opts,
    };
  }

  get state() { return this.#state; }

  /**
   * Wrap an async function with circuit breaker protection.
   */
  async call(fn) {
    if (this.#state === 'open') {
      if (Date.now() - this.#lastFailureTime > this.#opts.resetTimeoutMs) {
        this.#state = 'half-open';
        this.#successCount = 0;
      } else {
        throw new CircuitOpenError(
          `Circuit open — retry after ${this.#opts.resetTimeoutMs}ms`
        );
      }
    }

    try {
      const result = await fn();
      this.#onSuccess();
      return result;
    } catch (err) {
      this.#onFailure();
      throw err;
    }
  }

  #onSuccess() {
    if (this.#state === 'half-open') {
      this.#successCount++;
      if (this.#successCount >= this.#opts.successThreshold) {
        this.#state = 'closed';
        this.#failureCount = 0;
      }
    } else {
      this.#failureCount = 0;
    }
  }

  #onFailure() {
    this.#failureCount++;
    this.#lastFailureTime = Date.now();
    if (this.#failureCount >= this.#opts.failureThreshold) {
      this.#state = 'open';
    }
  }
}

class CircuitOpenError extends Error {
  constructor(msg) {
    super(msg);
    this.name = 'CircuitOpenError';
    this.retriable = false;   // do not retry — wait for reset
  }
}
```

### Integrating Circuit Breaker with the Registry

The circuit breaker wraps the handler at the per-operation or per-downstream level:

```js
registry.register('fetchInventory', async (params, op) => {
  return inventoryBreaker.call(() => inventoryApi.get(params.itemId));
}, { maxRetries: 3 });
```

## Error Classification

Classifying errors as transient or terminal is the single most important decision in the retry path. Getting it wrong means either wasting retries on permanent failures or dropping recoverable operations.

```js
/**
 * Error classifier.
 * Returns true if the error is terminal (should NOT be retried).
 */
function isTerminal(err) {
  // Explicit opt-out from retry
  if (err.retriable === false) return true;

  // HTTP status codes
  const status = err.status || err.statusCode;
  if (status) {
    // 4xx client errors are usually terminal (except 408, 429)
    if (status >= 400 && status < 500 && status !== 408 && status !== 429) {
      return true;
    }
  }

  // Named error codes that are never transient
  const terminalCodes = new Set([
    'ERR_INVALID_ARG',
    'AUTH_FAILED',
    'FORBIDDEN',
    'NOT_FOUND',
    'VALIDATION_ERROR',
    'SCHEMA_MISMATCH',
  ]);
  if (terminalCodes.has(err.code)) return true;

  return false;
}

/**
 * Classify an error into a structured category.
 */
function classifyError(err) {
  if (isTerminal(err)) {
    return { category: 'terminal', retriable: false, action: 'abort' };
  }

  const status = err.status || err.statusCode;

  if (status === 429 || err.code === 'RATE_LIMITED') {
    return { category: 'rate_limit', retriable: true, action: 'backoff_extended' };
  }

  if (err.code === 'ETIMEDOUT' || err.code === 'ECONNRESET' || status === 503) {
    return { category: 'transient', retriable: true, action: 'retry_standard' };
  }

  if (err.code === 'CIRCUIT_OPEN') {
    return { category: 'circuit_open', retriable: false, action: 'wait_for_reset' };
  }

  // Default: treat unknown errors as transient (with limited retries)
  return { category: 'unknown', retriable: true, action: 'retry_cautious' };
}
```

## Auto-Remediation

Auto-remediation goes beyond retry: it detects a failure, classifies the root cause, applies a targeted fix, and verifies the fix before resuming the operation.

### Remediation Registry

```js
class RemediationRegistry {
  #remediators = new Map();

  /**
   * Register a remediator for a specific error pattern.
   * @param {string} errorPattern - error code or category to match
   * @param {object} remediator   - { detect, fix, verify, riskLevel }
   */
  register(errorPattern, remediator) {
    this.#remediators.set(errorPattern, {
      detect: remediator.detect,       // (err, op) => boolean
      fix: remediator.fix,             // (err, op) => Promise<void>
      verify: remediator.verify,       // (op) => Promise<boolean>
      riskLevel: remediator.riskLevel || 'low',  // low | medium | high
      ...remediator,
    });
  }

  /**
   * Attempt auto-remediation for a failed operation.
   * Returns true if remediation succeeded and the operation can be re-queued.
   */
  async attempt(err, op) {
    for (const [pattern, rem] of this.#remediators) {
      if (!rem.detect(err, op)) continue;

      // High-risk remediations require explicit approval
      if (rem.riskLevel === 'high') {
        auditLog(op, 'remediation_blocked', { pattern, reason: 'high_risk' });
        return false;
      }

      try {
        auditLog(op, 'remediation_started', { pattern });
        await rem.fix(err, op);

        const verified = await rem.verify(op);
        if (verified) {
          auditLog(op, 'remediation_verified', { pattern });
          transition(op, 'remediated');
          return true;
        } else {
          auditLog(op, 'remediation_verification_failed', { pattern });
          return false;
        }
      } catch (remErr) {
        auditLog(op, 'remediation_error', { pattern, error: remErr.message });
        return false;
      }
    }

    return false; // no matching remediator found
  }
}
```

### Example: Token Refresh Remediator

```js
remediationRegistry.register('AUTH_TOKEN_EXPIRED', {
  detect: (err) => err.code === 'AUTH_TOKEN_EXPIRED' || err.status === 401,
  fix: async (err, op) => {
    const newToken = await authClient.refreshToken();
    op.params._authToken = newToken;
  },
  verify: async (op) => {
    const token = op.params._authToken;
    return token && !isTokenExpired(token);
  },
  riskLevel: 'low',
});
```

### Integrating Remediation into the Execution Loop

Add a remediation attempt between failure classification and the retry/terminal decision:

```js
async #execute(op, entry) {
  const { handler, opts } = entry;

  while (op.attempt <= opts.maxRetries) {
    op.attempt++;
    transition(op, 'running');

    try {
      op.result = await withTimeout(handler(op.params, op), opts.timeoutMs);
      transition(op, 'success');
      return op;
    } catch (err) {
      op.error = serializeError(err);

      // Attempt auto-remediation before deciding retry/terminal
      const remediated = await this.#remediation.attempt(err, op);
      if (remediated) {
        op.attempt--; // do not count remediated attempt against retry budget
        continue;     // re-run with the fix applied
      }

      if (isTerminal(err) || op.attempt > opts.maxRetries) {
        transition(op, 'failed');
        return op;
      }

      transition(op, 'retry');
      await sleep(computeDelay(op.attempt, opts));
    }
  }

  transition(op, 'failed');
  return op;
}
```

## Idempotency

An operation is idempotent if executing it once or N times produces the same observable effect. Idempotency is the foundation of safe retries.

### Idempotency Key Store

```js
class IdempotencyStore {
  #store = new Map();  // In production: use Redis or a database
  #ttlMs;

  constructor({ ttlMs = 24 * 60 * 60 * 1000 } = {}) {
    this.#ttlMs = ttlMs;
  }

  /**
   * Check if this key was already processed.
   * Returns the cached result or null.
   */
  get(key) {
    const entry = this.#store.get(key);
    if (!entry) return null;
    if (Date.now() - entry.createdAt > this.#ttlMs) {
      this.#store.delete(key);
      return null;
    }
    return entry;
  }

  /**
   * Store the result for an idempotency key.
   */
  set(key, result, status) {
    this.#store.set(key, { result, status, createdAt: Date.now() });
  }
}
```

### Wrapping the Registry Dispatch with Idempotency

```js
async dispatch(name, params, { idempotencyKey } = {}) {
  const entry = this.#handlers.get(name);
  if (!entry) throw new Error(`Unknown operation: "${name}"`);

  // Idempotency check: return cached result if the key was already processed
  if (idempotencyKey) {
    const cached = this.#idempotencyStore.get(idempotencyKey);
    if (cached) {
      if (cached.status === 'running') {
        throw new ConflictError('Operation is already in progress');
      }
      return cached.result;
    }
  }

  const op = createOperation(name, params, entry.opts, idempotencyKey);

  if (idempotencyKey) {
    this.#idempotencyStore.set(idempotencyKey, null, 'running');
  }

  const result = await this.#execute(op, entry);

  if (idempotencyKey) {
    this.#idempotencyStore.set(idempotencyKey, result, result.status);
  }

  return result;
}
```

**Key rules for idempotency:**
1. The **client** generates the idempotency key (e.g., `payment_${orderId}_${amount}`).
2. The store check happens **before** any side effects.
3. In-progress operations return a conflict, not a stale result.
4. TTL must outlive the longest possible retry window.
5. Compensation (saga rollback) must also be idempotent.

## Saga / Compensation

The saga pattern breaks a multi-step operation into a sequence of local transactions, each with a compensating action that undoes its effect on failure.

### Saga Orchestrator

```js
class SagaOrchestrator {
  #steps = [];

  /**
   * Add a step with its forward action and compensating action.
   * Both must be idempotent.
   */
  step(name, execute, compensate) {
    this.#steps.push({ name, execute, compensate });
    return this; // chainable
  }

  /**
   * Run the saga. On failure, execute compensating actions
   * for all completed steps in reverse order.
   */
  async run(context = {}) {
    const completed = [];
    const audit = [];

    for (const step of this.#steps) {
      try {
        audit.push({ step: step.name, action: 'execute', ts: Date.now() });
        const result = await step.execute(context);
        context[step.name] = result;
        completed.push(step);
      } catch (err) {
        audit.push({ step: step.name, action: 'failed', error: err.message, ts: Date.now() });

        // Compensate in reverse order
        for (const done of completed.reverse()) {
          try {
            audit.push({ step: done.name, action: 'compensate', ts: Date.now() });
            await done.compensate(context);
          } catch (compErr) {
            // Compensation failure is critical — log and continue
            audit.push({
              step: done.name,
              action: 'compensate_failed',
              error: compErr.message,
              ts: Date.now(),
            });
          }
        }

        return { status: 'rolled_back', failedStep: step.name, error: err, audit };
      }
    }

    return { status: 'completed', context, audit };
  }
}
```

### Example: Order Fulfillment Saga

```js
const orderSaga = new SagaOrchestrator()
  .step(
    'reserveInventory',
    async (ctx) => inventoryService.reserve(ctx.orderId, ctx.items),
    async (ctx) => inventoryService.release(ctx.orderId, ctx.items),
  )
  .step(
    'chargePayment',
    async (ctx) => paymentService.charge(ctx.orderId, ctx.amount),
    async (ctx) => paymentService.refund(ctx.orderId, ctx.amount),
  )
  .step(
    'shipOrder',
    async (ctx) => shippingService.create(ctx.orderId, ctx.address),
    async (ctx) => shippingService.cancel(ctx.orderId),
  );

const result = await orderSaga.run({
  orderId: 'order-123',
  items: [{ sku: 'WIDGET-A', qty: 2 }],
  amount: 49_99,
  address: '123 Main St',
});
```

**Choreography vs. orchestration:** Use orchestration (a central coordinator) when you need clear visibility and control over the saga flow. Use choreography (event-driven, no coordinator) when services are independently owned and you want loose coupling. Orchestration is easier to debug; choreography scales better across team boundaries.

## Concurrency Control

### Semaphore for Concurrent Operation Limiting

```js
class Semaphore {
  #permits;
  #queue = [];

  constructor(permits) {
    this.#permits = permits;
  }

  async acquire() {
    if (this.#permits > 0) {
      this.#permits--;
      return;
    }
    // Wait for a permit to become available
    return new Promise(resolve => this.#queue.push(resolve));
  }

  release() {
    if (this.#queue.length > 0) {
      const next = this.#queue.shift();
      next();
    } else {
      this.#permits++;
    }
  }

  /**
   * Wrap an async function with semaphore-bounded execution.
   */
  async run(fn) {
    await this.acquire();
    try {
      return await fn();
    } finally {
      this.release();
    }
  }
}
```

### Rate-Limited Queue

```js
class RateLimitedQueue {
  #semaphore;
  #minIntervalMs;
  #lastRunAt = 0;

  constructor({ concurrency = 5, ratePerSecond = 10 } = {}) {
    this.#semaphore = new Semaphore(concurrency);
    this.#minIntervalMs = 1000 / ratePerSecond;
  }

  async enqueue(fn) {
    return this.#semaphore.run(async () => {
      const now = Date.now();
      const elapsed = now - this.#lastRunAt;
      if (elapsed < this.#minIntervalMs) {
        await sleep(this.#minIntervalMs - elapsed);
      }
      this.#lastRunAt = Date.now();
      return fn();
    });
  }
}
```

### Integrating Concurrency Control into the Registry

```js
// Per-operation concurrency limits
registry.register('fetchUserProfile', handler, {
  maxRetries: 2,
  concurrency: 10,   // max 10 concurrent executions of this operation
});

// The registry creates a semaphore per handler at registration time
register(name, handler, opts = {}) {
  const semaphore = opts.concurrency
    ? new Semaphore(opts.concurrency)
    : null;

  this.#handlers.set(name, { handler, opts, semaphore });
}
```

## Audit Logging

Every operation state transition should produce a structured audit event. This is non-negotiable for debugging distributed retry flows.

```js
/**
 * Record a state transition in the operation's audit history.
 */
function transition(op, newStatus) {
  const event = {
    from: op.status,
    to: newStatus,
    attempt: op.attempt,
    ts: Date.now(),
    error: op.error || null,
  };

  op.history.push(event);
  op.status = newStatus;
  op.updatedAt = Date.now();

  // Emit for external consumers (logging, metrics, alerting)
  emitAuditEvent({
    opId: op.id,
    opName: op.name,
    idempotencyKey: op.idempotencyKey,
    ...event,
  });
}

/**
 * Structured audit event emitter.
 * In production, pipe to your observability stack.
 */
function emitAuditEvent(event) {
  const entry = {
    level: event.to === 'failed' ? 'error' : 'info',
    msg: `op:${event.opName} ${event.from}->${event.to}`,
    opId: event.opId,
    opName: event.opName,
    idempotencyKey: event.idempotencyKey,
    fromStatus: event.from,
    toStatus: event.to,
    attempt: event.attempt,
    error: event.error,
    ts: new Date(event.ts).toISOString(),
  };

  // Structured JSON log line — parseable by any log aggregator
  console.log(JSON.stringify(entry));
}
```

**What to log per operation:**
- Operation ID, name, and idempotency key
- Every state transition with timestamps
- Error details (code, message, classification) on failure
- Retry delay and attempt number
- Remediation attempts and outcomes
- Saga step completions and compensations
- Circuit breaker state changes
- Concurrency slot acquisition/release times (for latency debugging)

## Helpers

Utility functions referenced throughout the examples:

```js
function createOperation(name, params, opts, idempotencyKey) {
  return {
    id: crypto.randomUUID(),
    idempotencyKey: idempotencyKey || null,
    name,
    params,
    status: 'pending',
    attempt: 0,
    maxRetries: opts.maxRetries,
    createdAt: Date.now(),
    updatedAt: Date.now(),
    result: null,
    error: null,
    history: [],
  };
}

function withTimeout(promise, ms) {
  return Promise.race([
    promise,
    new Promise((_, reject) =>
      setTimeout(() => reject(new TimeoutError(`Timed out after ${ms}ms`)), ms)
    ),
  ]);
}

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

function serializeError(err) {
  return {
    name: err.name,
    message: err.message,
    code: err.code,
    status: err.status || err.statusCode,
    stack: err.stack,
  };
}

class TimeoutError extends Error {
  constructor(msg) { super(msg); this.name = 'TimeoutError'; this.retriable = true; }
}
class ConflictError extends Error {
  constructor(msg) { super(msg); this.name = 'ConflictError'; this.retriable = false; }
}
```

## Decision Matrix: Which Pattern Do I Need?

| Symptom / Goal | Primary Pattern | Also Consider |
|---|---|---|
| Need to map named operations to handlers | Registry & Dispatch | Audit Logging |
| Operations fail intermittently | Retry with Backoff + Jitter | Error Classification |
| Downstream service goes down for minutes | Circuit Breaker | Retry Budget |
| Retried operations cause duplicates | Idempotency Keys | Error Classification |
| Auth tokens expire mid-operation | Auto-Remediation | Retry (after fix) |
| Multi-step workflow needs rollback on failure | Saga / Compensation | Idempotency |
| Too many concurrent calls overwhelm a resource | Concurrency Control (Semaphore) | Rate-Limited Queue |
| Cannot debug failed operations in production | Audit Logging | Structured Error Classification |
| Retry storms under load | Retry Budget | Circuit Breaker |
| Need all of the above composed together | Full Registry (see Checklist) | -- |

## Anti-Patterns

1. **Retrying terminal errors.** A 403 Forbidden will never succeed on retry. Classify errors before retrying.
2. **No jitter on backoff.** Pure exponential backoff across N clients creates a thundering herd at each retry interval.
3. **Unbounded retries.** Always set `maxRetries`. Infinite retry loops cause resource exhaustion and mask permanent failures.
4. **Retry without idempotency.** Retrying a non-idempotent operation (e.g., `POST /charge`) can double-charge. Always pair retries with idempotency keys.
5. **Ignoring circuit breaker in retry loops.** If the downstream is down, retrying 10 times per request multiplied by 1000 concurrent requests creates 10,000 wasted calls. Use a circuit breaker to fail fast.
6. **Compensation that is not idempotent.** If a saga compensation step runs twice (because the compensator itself failed mid-way), it must produce the same effect. Always check before acting (e.g., check if refund already issued before issuing another).
7. **Logging unstructured error strings.** `"Something went wrong"` is useless in production. Log structured fields: operation ID, error code, attempt number, classification.
8. **Shared mutable state across retries.** Each retry attempt should receive a fresh context or a deliberately carried-forward one. Do not accidentally accumulate side effects across attempts.
9. **Coupling remediation to the handler.** Remediation logic belongs in the remediation registry, not inside individual handlers. Handlers should throw; the registry decides what to do.
10. **No timeout on operations.** A handler that hangs forever holds a concurrency slot forever. Always wrap with `withTimeout`.

## Checklist: Building an Operations Registry

- [ ] Define the operation envelope shape with ID, idempotency key, status, and history
- [ ] Register handlers at startup with per-operation retry/timeout overrides
- [ ] Implement exponential backoff with full jitter (not pure exponential)
- [ ] Classify errors as transient vs terminal before deciding to retry
- [ ] Add a circuit breaker per downstream dependency
- [ ] Enforce idempotency keys for any operation with side effects
- [ ] Implement a retry budget to cap total retry traffic
- [ ] Add concurrency control (semaphore) per operation or globally
- [ ] Wire auto-remediation for known recoverable failure patterns
- [ ] Use saga orchestration for multi-step operations with compensating actions
- [ ] Emit structured audit events on every state transition
- [ ] Set timeouts on every handler invocation
- [ ] Test: verify retries stop on terminal errors
- [ ] Test: verify idempotent replay returns cached result
- [ ] Test: verify circuit breaker opens after threshold and recovers after reset

## References

**Registry & Dispatch**
- [Function Registry Pattern — iO TechHub](https://techhub.iodigital.com/articles/function-registry-pattern-react)
- [Registry Pattern — GeeksforGeeks](https://www.geeksforgeeks.org/system-design/registry-pattern/)
- [Registry Pattern in Programming — Mark Torres](https://markptorres.com/ai_workflows/2025-12-02-registry-pattern)

**Retry, Backoff & Circuit Breaker**
- [Node.js Advanced Retry Logic — V. Checha](https://v-checha.medium.com/advanced-node-js-patterns-implementing-robust-retry-logic-656cf70f8ee9)
- [Circuit Breaker with Exponential Backoff — Usama Uzair](https://medium.com/@usama19026/building-resilient-applications-circuit-breaker-pattern-with-exponential-backoff-fc14ba0a0beb)
- [Exponential Backoff with Jitter — Tito Adeoye](https://medium.com/@titoadeoye/requests-at-scale-exponential-backoff-with-jitter-with-examples-4d0521891923)
- [Circuit Breaker & Retry Patterns in Node.js 2026 — 1xAPI](https://1xapi.com/blog/resilient-api-circuit-breaker-bulkhead-retry-nodejs-2026)
- [Retry with Exponential Backoff in Node.js — OneUptime](https://oneuptime.com/blog/post/2026-01-06-nodejs-retry-exponential-backoff/view)
- [Error Classification for Retry — Alex Bogomol](https://medium.com/@bogomolalexander/designing-a-retry-strategy-for-transient-errors-a5cd8b4d0602)
- [Mastering Retry Logic Agents 2025 — SparkCo](https://sparkco.ai/blog/mastering-retry-logic-agents-a-deep-dive-into-2025-best-practices)

**Idempotency**
- [Idempotency Patterns in Distributed Systems — BackendBytes](https://backendbytes.com/articles/idempotency-patterns-distributed-systems/)
- [Idempotency Design Patterns Beyond Retry Safely — Dev.to](https://dev.to/aloknecessary/idempotency-in-distributed-systems-design-patterns-beyond-retry-safely-k66)

**Saga & Compensation**
- [Saga Pattern — Microservices.io](https://microservices.io/patterns/data/saga.html)
- [Saga Pattern Guide — OneUptime](https://oneuptime.com/blog/post/2026-01-24-saga-pattern-transactions/view)
- [Saga Design Pattern — Microsoft Azure](https://learn.microsoft.com/en-us/azure/architecture/patterns/saga)
- [Saga Patterns — AWS Prescriptive Guidance](https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/saga.html)

**Concurrency Control**
- [Semaphore Pattern in JS — Rohit Nandi](https://www.rohitnandi.com/blog/semaphore-pattern-js)
- [Concurrency Patterns in JS — Artem Khrienov](https://medium.com/@artemkhrenov/advanced-concurrency-patterns-in-javascript-semaphore-mutex-read-write-lock-deadlock-prevention-79e8bffb5b81)

**Auto-Remediation & Audit Logging**
- [Auto-Remediation in CI/CD — Cy5](https://www.cy5.io/blog/from-alerts-to-action-cspm-automated-remediation/)
- [AI Observability 2026: Auto-Remediation — Rootly](https://rootly.com/sre/ai-observability-2026-predictive-alerts-auto-remediation-577de)
- [Structured Logging in OpenTelemetry — OneUptime](https://oneuptime.com/blog/post/2025-08-28-how-to-structure-logs-properly-in-opentelemetry/view)
- [Audit Logs Guide — Middleware.io](https://middleware.io/blog/audit-logs/)
