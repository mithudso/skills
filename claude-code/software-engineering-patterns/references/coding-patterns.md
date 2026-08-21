<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `coding-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: coding-patterns
description: >
  Software design patterns for JavaScript/TypeScript: creational (factory, builder, singleton, prototype), structural (facade, adapter, proxy, decorator, composite), behavioral (observer, strategy, state machine, command, chain of responsibility, mediator), event-driven (pub/sub, dispatch tables, event sourcing, postMessage bridge), async (retry/backoff, circuit breaker, semaphore, async queue), functional (pipe/flow, memoization, currying, reducers), and module patterns. Use when designing system architecture, choosing between structural approaches, or refactoring code into maintainable patterns.
  TRIGGER: "design pattern", "factory pattern", "observer pattern", "state machine", "circuit breaker", "retry", "dispatch table", "pub/sub", "memoize", "dependency injection", "facade", "adapter", "how should I structure", "which pattern", "refactor this module".
  SKIP: naming conventions or readability rules (use coding-standards); React component patterns or hooks (use frontend-design); REST/API endpoint design (use backend-patterns or api-design-patterns); performance profiling or memory debugging (use performance-profiling-expert).
version: 1.0.1
category: developer
tags: [design-patterns, javascript, typescript, architecture, event-driven, async, functional]
updated: "2026-05-29"
related_skills:
  - coding-standards
  - code-reviewer
  - frontend-design
  - backend-patterns
  - mv3-service-worker-expert
whenToUse:
  - designing a new module, service, or subsystem and choosing its internal structure
  - reviewing whether a codebase applies patterns correctly or uses the wrong pattern
  - refactoring tangled code into a recognized, maintainable shape
  - adding resilience (retry, circuit breaker, backoff) to async operations
  - building message-routing or event-driven dispatch layers
  - introducing functional composition, memoization, or pipeline transforms
whenNotToUse:
  - naming, formatting, or linting rules (use coding-standards)
  - React component patterns, hooks, or rendering (use frontend-design)
  - REST/API endpoint design (use backend-patterns or api-design-patterns)
  - performance profiling or memory debugging (use performance-profiling-expert)
---

# Coding Patterns

Software design patterns for JavaScript and TypeScript, with emphasis on Chrome MV3 extensions and event-driven architectures.

**Scope:** structural design decisions — how objects are created, composed, and communicate. For naming, readability, and code-quality review, use `coding-standards`.

---

## 1. Creational Patterns

### 1a. Factory

Returns new objects without exposing instantiation logic. Callers depend on an interface, not a concrete class.

```js
function createPaymentProvider(type) {
  const providers = {
    stripe: () => ({ charge(amount) { return stripeSDK.charges.create({ amount }); } }),
    paypal: () => ({ charge(amount) { return paypalSDK.payment.create({ amount }); } }),
  };
  const build = providers[type];
  if (!build) throw new Error(`Unknown provider: ${type}`);
  return build();
}
```

**Pick when:** multiple concrete implementations behind a shared interface; the caller should not know which concrete type it holds.

### 1b. Builder

Constructs complex objects step-by-step, separating construction from representation.

```js
class QueryBuilder {
  #parts = { select: '*', from: '', where: [], order: '' };

  from(table)       { this.#parts.from = table; return this; }
  select(...cols)   { this.#parts.select = cols.join(', '); return this; }
  where(condition)  { this.#parts.where.push(condition); return this; }
  orderBy(col, dir) { this.#parts.order = `${col} ${dir}`; return this; }

  build() {
    const { select, from, where, order } = this.#parts;
    let sql = `SELECT ${select} FROM ${from}`;
    if (where.length) sql += ` WHERE ${where.join(' AND ')}`;
    if (order) sql += ` ORDER BY ${order}`;
    return sql;
  }
}
```

**Pick when:** object requires many optional configuration steps.

### 1c. Singleton (Module-Scoped)

In JS/TS a module is evaluated once — export an instance from a module for a singleton without GoF ceremony. A shared `logger` or `config` object exported from `logger.js` is the canonical example.

**Pick when:** exactly one shared instance (loggers, config holders, connection pools). Avoid when you need per-context isolation.

### 1d. Prototype / Object.create

Clone an existing object as a template via `Object.create(baseHandler)`, then override specific methods. **Pick when:** well-defined base template with lightweight variants; class hierarchies would be overkill.

---

## 2. Structural Patterns

### 2a. Facade

Provides a simplified interface to a complex subsystem.

```js
export const storage = {
  async get(key) {
    const result = await chrome.storage.local.get(key);
    return result[key] ?? null;
  },
  async set(key, value) { await chrome.storage.local.set({ [key]: value }); },
  async remove(key)     { await chrome.storage.local.remove(key); },
  onChange(key, callback) {
    chrome.storage.onChanged.addListener((changes, area) => {
      if (area === 'local' && changes[key]) callback(changes[key].newValue, changes[key].oldValue);
    });
  },
};
```

**Pick when:** a subsystem has a broad API but callers need only a narrow slice.

### 2b. Adapter

Translates one interface into another that the caller expects.

```js
function adaptCallbackToPromise(legacyFn) {
  return (...args) =>
    new Promise((resolve, reject) => {
      legacyFn(...args, (err, result) => { if (err) reject(err); else resolve(result); });
    });
}
const readFile = adaptCallbackToPromise(fs.readFile);
```

**Pick when:** integrating third-party or legacy code whose interface doesn't match your expected contract.

### 2c. Proxy

Wraps an object to intercept and control access.

```js
function createValidatedConfig(defaults) {
  const validators = {
    port:    (v) => Number.isInteger(v) && v > 0 && v < 65536,
    retries: (v) => Number.isInteger(v) && v >= 0,
    baseUrl: (v) => typeof v === 'string' && v.startsWith('http'),
  };
  return new Proxy({ ...defaults }, {
    set(target, prop, value) {
      const validate = validators[prop];
      if (validate && !validate(value)) throw new TypeError(`Invalid value for ${prop}: ${value}`);
      target[prop] = value;
      return true;
    },
  });
}
```

**Pick when:** lazy loading, access control, validation, logging, or caching around an existing object.

### 2d. Decorator (Function Wrapping)

Adds behavior to a function without changing its signature.

```js
function withLogging(fn, label) {
  return async function (...args) {
    const start = performance.now();
    try {
      const result = await fn.apply(this, args);
      console.log(`[${label}] returned in ${(performance.now() - start).toFixed(1)}ms`);
      return result;
    } catch (err) {
      console.error(`[${label}] threw`, err);
      throw err;
    }
  };
}
```

**Pick when:** cross-cutting concerns (logging, timing, auth checks) that apply to many functions.

### 2e. Composite

Treats individual objects and groups uniformly through a shared interface. Both `MenuItem` (leaf) and `MenuGroup` (branch) expose the same `execute()` method; `MenuGroup.execute()` delegates to each child. **Pick when:** tree structures (menus, file systems, UI component trees).

---

## 3. Behavioral Patterns

### 3a. Observer / EventEmitter

One-to-many dependency: subject changes notify all subscribers.

```js
class EventBus {
  #listeners = new Map();
  on(event, fn) {
    if (!this.#listeners.has(event)) this.#listeners.set(event, new Set());
    this.#listeners.get(event).add(fn);
    return () => this.#listeners.get(event)?.delete(fn); // returns unsubscribe handle
  }
  emit(event, data) {
    for (const fn of this.#listeners.get(event) ?? []) {
      try { fn(data); } catch (err) { console.error(`[EventBus] ${event}:`, err); }
    }
  }
}
```

**Pick when:** multiple parts of the system need to react to changes in one place without tight coupling.

### 3b. Strategy

Encapsulates interchangeable algorithms behind a common interface.

```js
const sortStrategies = {
  byDate:     (a, b) => new Date(b.created) - new Date(a.created),
  bySeverity: (a, b) => a.severity - b.severity,
  bySubject:  (a, b) => a.subject.localeCompare(b.subject),
};

function sortCases(cases, strategyName) {
  const compare = sortStrategies[strategyName];
  if (!compare) throw new Error(`Unknown sort: ${strategyName}`);
  return [...cases].sort(compare);
}
```

**Pick when:** multiple algorithms for the same task, selected at runtime.

### 3c. Command

Encapsulates a request as an object, enabling undo, queuing, and logging.

```js
class CommandHistory {
  #stack = [];
  execute(cmd) { cmd.execute(); this.#stack.push(cmd); }
  undo()       { const cmd = this.#stack.pop(); cmd?.undo(); }
}
```

**Pick when:** operations that need undo/redo, transaction logging, or deferred execution.

### 3d. State Machine

Governs a finite set of states and legal transitions.

```js
function createCaseStateMachine(initialState = 'open') {
  const transitions = {
    open:        { assign: 'assigned', close: 'closed' },
    assigned:    { start: 'in_progress', unassign: 'open', close: 'closed' },
    in_progress: { resolve: 'resolved', escalate: 'escalated', close: 'closed' },
    escalated:   { resolve: 'resolved', close: 'closed' },
    resolved:    { reopen: 'open', close: 'closed' },
    closed:      { reopen: 'open' },
  };
  let state = initialState;
  return {
    get state() { return state; },
    can(action) { return !!transitions[state]?.[action]; },
    send(action) {
      const next = transitions[state]?.[action];
      if (!next) throw new Error(`Cannot ${action} from ${state}`);
      const prev = state; state = next;
      return { prev, next: state, action };
    },
  };
}
```

**Pick when:** workflows with strict lifecycle rules (case tracking, connection states, auth flows, UI wizards).

### 3e. Chain of Responsibility

`createChain(...handlers)` iterates handlers in order; returns the first non-`undefined` result. Throws if nothing matches. **Pick when:** ordered fallback logic (auth: session → bearer → reject), middleware pipelines, or plugin hooks.

### 3f. Mediator

A single coordinator routes messages between N parties so no party holds a reference to any other. In Chrome MV3 the service worker is the natural mediator: `CONTENT_SCRAPED` → enrich → `PANEL_UPDATE`; `PANEL_REQUEST_REFRESH` → `RESCRAPE` tab message. **Pick when:** many-to-many context communication (content scripts ↔ popup ↔ panel ↔ offscreen) where direct coupling would create a dependency web.

---

## 4. Event-Driven Patterns

### 4a. Pub/Sub with Topics

```js
class PubSub {
  #subs = new Map();
  subscribe(topic, fn) {
    if (!this.#subs.has(topic)) this.#subs.set(topic, new Set());
    this.#subs.get(topic).add(fn);
    return () => this.#subs.get(topic).delete(fn);
  }
  publish(topic, data) {
    for (const fn of this.#subs.get(topic) ?? []) queueMicrotask(() => fn(data));
  }
}
```

### 4b. Type-Based Dispatch Table

The dominant pattern in Chrome MV3 extensions: a flat map from message type to handler.

```js
const handlers = {
  MCA_GET_CASE:       handleGetCase,
  MCA_UPSERT_CONTEXT: handleUpsertContext,
  MCA_TRACK_CASE:     handleTrackCase,
  MCA_SUMMARIZE:      handleSummarize,
};

chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
  const handler = handlers[msg.type];
  if (!handler) {
    logger.warn('Unhandled message type', { type: msg.type });
    return false;
  }
  handler(msg, sender)
    .then(result => sendResponse({ ok: true, data: result }))
    .catch(err  => sendResponse({ ok: false, error: err.message }));
  return true; // keep channel open for async
});
```

**Why dispatch tables beat switch/case:** adding a handler is one line; handlers are independently testable; you can iterate `Object.keys(handlers)` for logging or validation.

### 4c. Event Sourcing Lite

Store events (`{ type, payload, ts }`) in an append-only array; derive current state by reducing them through a `caseReducer`. Enables full replay, undo, and audit trails. Use when you need reproducible state from an event stream.

### 4d. postMessage Bridge (Extension-Specific)

Listen on `window.addEventListener('message', ...)`, guard on `event.source === window && event.data?.source === 'MCA_HOST'`, then forward to `iframe.contentWindow.postMessage(...)`. Bridges host-page events into the extension's shadow-DOM iframe without exposing internal message types.

---

## 5. Async Patterns

Full implementations in [`references/async-patterns.md`](references/async-patterns.md).

| Pattern | What it does | Pick when |
|---------|-------------|-----------|
| **5a. Retry + Backoff** | Retries a failing fn with exponential delay × jitter (0.75–1.25×) | Transient network failures; never use tight-loop retry |
| **5b. Circuit Breaker** | `closed → open → half-open` state machine; fails fast when a service is down | Prevent cascading failures; resume probing after cooldown |
| **5c. Semaphore** | Limits concurrent async invocations to N; queues excess | API rate limits; bounded `Promise.all` over large arrays |
| **5d. Async Queue** | Chains tasks onto a single promise so they run sequentially | Ordered writes to storage; serial mutation of shared state |

---

## 6. Functional Patterns

### 6a. Pipe / Flow

```js
const pipe = (...fns) => (x) => fns.reduce((v, fn) => fn(v), x);

const pipeAsync = (...fns) => (x) => fns.reduce((chain, fn) => chain.then(fn), Promise.resolve(x));
```

### 6b. Memoization

```js
function memoize(fn, { maxSize = 100, keyFn = JSON.stringify } = {}) {
  const cache = new Map();
  function memoized(...args) {
    const key = keyFn(args);
    if (cache.has(key)) return cache.get(key);
    const result = fn.apply(this, args);
    cache.set(key, result);
    if (cache.size > maxSize) cache.delete(cache.keys().next().value);
    return result;
  }
  memoized.clear = () => cache.clear();
  return memoized;
}
```

Only memoize pure functions with small, hashable argument signatures. Avoid `JSON.stringify` on large payloads.

### 6c. Currying

Transforms a multi-argument function into a sequence of single-argument functions: `const add = curry((a, b) => a + b); const increment = add(1); increment(5); // 6`. Enables partial application and reusable utility factories.

### 6d. Immutable State Updates (Reducer)

Pure `(state, action) => newState` function using spread to return new objects. The `UNTRACK` case uses destructuring to omit a key: `const { [id]: _, ...rest } = state.tracked`. Use with `useReducer` or a standalone store.

---

## 7. Module Patterns

### 7a. Revealing Module

IIFE that closes over private state (`authToken`, `BASE_URL`) and returns only the public API (`setToken`, `getCase`, `search`). Zero dependencies — works in any JS environment without a class.
```

### 7b. Barrel Exports

```js
// src/utils/index.js
export { formatDate, parseISO }    from './date-utils.js';
export { slugify, truncate }       from './string-utils.js';
export { deepMerge, pick }         from './object-utils.js';
```

### 7c. Dependency Injection Lite

Pass dependencies as arguments instead of importing directly; aids testing.

```js
function createCaseEnricher({ storage, apiClient, logger }) {
  return {
    async enrich(caseId) {
      const raw    = await apiClient.getCase(caseId);
      const cached = await storage.get(`case:${caseId}`);
      const merged = { ...cached, ...raw, lastEnriched: Date.now() };
      await storage.set(`case:${caseId}`, merged);
      return merged;
    },
  };
}
```

---

## 8. Pattern Selection Guide

| Problem / Scenario | Recommended Pattern | Avoid |
|---|---|---|
| Multiple implementations behind one interface | Factory | Giant if/else chains |
| Object needs many optional config steps | Builder | Constructors with 10+ params |
| Exactly one shared instance (logger, config) | Module-scoped Singleton | Class-based Singleton with `getInstance()` |
| Simplify a complex subsystem for callers | Facade | Exposing all internal APIs directly |
| Integrate code with an incompatible interface | Adapter | Modifying third-party source |
| Intercept property access (validate, cache) | Proxy | Monkey-patching the target object |
| Add cross-cutting behavior to functions | Decorator (fn wrapping) | Copy-pasting logic into every function |
| Multiple subscribers react to state changes | Observer / EventBus | Direct function calls from the subject |
| Select algorithm at runtime | Strategy | switch/case sprawl per algorithm |
| Encapsulate undoable operations | Command | Inline mutations without history |
| Strict lifecycle with legal transitions | State Machine | Boolean flags + nested if/else |
| Pass request through ordered processors | Chain of Responsibility | One giant handler with all logic |
| Central coordinator, multi-party comms | Mediator | Every component talking to every other |
| Chrome MV3 message routing | Dispatch Table | switch/case in onMessage |
| Host page / shadow DOM / iframe messaging | postMessage Bridge | Direct DOM manipulation across boundaries |
| Transient network failure recovery | Retry with Backoff + Jitter | Immediate retry in a tight loop |
| Prevent cascading failure to a bad service | Circuit Breaker | Retrying forever against a down service |
| Limit concurrent requests | Semaphore | Unbounded `Promise.all` on large arrays |
| Guarantee sequential async processing | Async Queue | Hoping concurrent writes don't conflict |
| Transform data through multi-step pipeline | Pipe / Flow | Deeply nested function calls |
| Cache expensive pure-function results | Memoization | Recomputing on every call |
| Partial application for reusable utilities | Currying | Duplicating near-identical functions |
| Predictable state transitions | Reducer (immutable updates) | Direct state mutation |
| Testable code with swappable dependencies | Dependency Injection Lite | Hard-coded imports in every function |

---

## 9. Anti-Patterns

| Anti-Pattern | Symptom | Fix |
|---|---|---|
| **God Object / God Module** | One file handles routing, storage, auth, enrichment, and rendering | Split into focused modules with single responsibilities |
| **Premature Abstraction** | Factory + Strategy + Observer for a feature with one implementation | Start concrete; extract when the second or third variant appears |
| **Pattern Soup** | Every GoF pattern in a 500-line module | Apply the simplest pattern that solves the problem |
| **Callback Pyramid** | Deeply nested callbacks | Refactor to async/await or pipe/flow |
| **Mutable Singleton State** | Module-scoped state mutated by multiple callers without synchronization | Use a reducer or async queue to serialize mutations |
| **Dispatch Table Without Validation** | Unknown message types silently dropped | Log or throw on unrecognized types |
| **Retry Without Backoff** | Tight retry loop amplifies load on a struggling service | Always add exponential delay with jitter |
| **Over-Memoizing** | Memoizing large-payload or side-effectful functions | Only memoize pure functions with small, hashable args |

---

## 10. References

- Refactoring Guru — Design Patterns in TypeScript: https://refactoring.guru/design-patterns/typescript
- Chrome Developers — Extension Service Workers: https://developer.chrome.com/docs/extensions/mv3/service_workers/events/
- Chrome Developers — Message Passing: https://developer.chrome.com/docs/extensions/mv3/messaging/
- XState — JavaScript State Machines: https://xstate.js.org/
- Opossum — Circuit Breaker for Node.js: https://github.com/nodeshift/opossum
