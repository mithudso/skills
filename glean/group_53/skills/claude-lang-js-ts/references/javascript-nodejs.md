<!-- hub-reference-banner -->
> **Reference file — part of the `lang-js-ts` hub.** Formerly the standalone `javascript-nodejs` skill.
> Sibling topics in this family are now reference files under the hubs (`lang-python`, `lang-js-ts`, `lang-go-and-mobile`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: javascript-nodejs
title: JavaScript and Node.js Expert
description: >
  JavaScript and Node.js expert reference for language semantics, runtime APIs,
  async patterns, modules, streams, error handling, and performance.
  TRIGGER: writing or reviewing JS/Node.js code; choosing between ESM and CommonJS;
  async/await or promise composition; Node.js streams, EventEmitter, fs APIs,
  process lifecycle, or perf_hooks; branching on error.code vs error.message;
  Map/Set vs plain objects. SKIP: TypeScript type system (typescript-expert);
  debugging tools or DevTools (javascript-node-html-css-debugging-expert);
  Vitest or test framework setup (testing-and-vitest-expert); shell scripting
  (shell-scripting); Python (python-patterns); MongoDB Node.js driver
  (mongodb-developer).
version: "1.1"
category: developer
tags:
  - javascript
  - nodejs
  - node
  - async
  - esm
  - commonjs
  - streams
  - modules
  - error-handling
  - performance
related_skills:
  - typescript-expert
  - javascript-node-html-css-debugging-expert
  - testing-and-vitest-expert
  - shell-scripting
---

# JavaScript and Node.js Expert

Behavior and API reference for generating or reviewing JavaScript and Node.js code. Treat **MDN JavaScript** and the **Node.js API docs** as the source of truth for semantics, API contracts, and version-sensitive behavior. For exhaustive member lists, use the linked reference indexes:
[MDN JavaScript Reference](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference) and [Node.js API index v26.1.0](https://nodejs.org/api/index.html).

> **Deep concurrency internals live elsewhere.** For the libuv event-loop phase model and
> microtask ordering, `UV_THREADPOOL_SIZE` tuning, "don't block the event loop", stream
> backpressure mechanics (`highWaterMark` / `write()` / `'drain'` / `pipeline`), or choosing
> between `worker_threads` / `cluster` / `child_process`, load
> `references/nodejs-concurrency-internals.md` (same hub). This file covers the broad language
> and runtime-API surface; that one covers ordering, blocking, throughput, and parallelism.

## When to use this skill

- Writing or reviewing JavaScript or Node.js code
- Choosing between ESM and CommonJS, or setting `package.json#"type"`
- Async/await, promise composition, or Node.js callback patterns
- Node.js streams, backpressure, or `stream.pipeline()`
- EventEmitter semantics and listener registration order
- Node.js `fs` API selection (callback vs promise vs sync)
- Branching on `error.code` vs `error.message`
- Map/Set vs plain-object tradeoffs
- `process` lifecycle, exit hooks, or the permission model
- Performance measurement with `node:perf_hooks`
- Node.js module resolution, `"exports"` maps, or `"imports"`

## When NOT to use this skill

- TypeScript type system, generics, or tsconfig → `typescript-expert`
- Debugging tools, Chrome DevTools, heap analysis → `javascript-node-html-css-debugging-expert`
- Vitest, test frameworks, or coverage configuration → `testing-and-vitest-expert`
- Shell scripting or CLI toolchain → `shell-scripting`
- MongoDB Node.js driver patterns → `mongodb-developer`
- Python code → `python-patterns`

---

## Quick rules

1. Prefer **explicit module markers**: `.mjs` / `.cjs` or `package.json#"type"`; Node explicitly recommends being explicit, especially for packages ([Node.js ESM](https://nodejs.org/api/esm.html), [Node.js Packages](https://nodejs.org/api/packages.html)).
2. Use **`async` / `await` with `try` / `catch`** for promise-based APIs; async functions always return promises and uncaught exceptions reject them ([MDN async function](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/async_function)).
3. Choose the `fs` API form intentionally: **callback APIs for maximal performance**, promise APIs for ergonomics, sync APIs only at startup or in CLI tools where blocking the event loop is not a problem ([Node.js fs](https://nodejs.org/api/fs.html)).
4. Use **`stream.pipeline()`** (not manual `.pipe()` chains) for multi-stream flows; `pipeline()` handles cleanup and error forwarding where `.pipe()` does not ([Node.js Stream](https://nodejs.org/api/stream.html)).
5. `EventEmitter` listeners run **synchronously in registration order**; move work async with `setImmediate()` or `process.nextTick()` when needed ([Node.js Events](https://nodejs.org/api/events.html)).
6. Match Node errors by **`error.code`**, not `error.message` — message text can change across releases. Example: `if (err.code === 'ENOENT')` not `if (err.message.includes('no such file'))` ([Node.js Errors](https://nodejs.org/api/errors.html)).
7. Prefer **static `Object.*` utilities** over instance calls to non-polymorphic `Object.prototype` methods ([MDN Object](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object)).
8. Use **`Map`** and **`Set`** when you need keyed collections, uniqueness, insertion order, or faster membership checks than array scans ([MDN Map](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Map), [MDN Set](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Set)).
9. Avoid `eval()` unless you truly need dynamic script evaluation ([MDN eval](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/eval)).
10. Use **`node:perf_hooks`** for timing and event-loop measurements instead of ad hoc timing ([Node.js perf_hooks](https://nodejs.org/api/perf_hooks.html)).

---

## Language and runtime reference

### JavaScript language model

- Dynamic, prototype-based, garbage-collected language; standards basis is **ECMAScript (ECMA-262)** plus **ECMA-402** ([MDN JavaScript](https://developer.mozilla.org/en-US/docs/Web/JavaScript)).
- MDN guide/reference covers grammar and types, control flow, functions, expressions, collections, iterators/generators, promises, and modules ([MDN Guide](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide), [MDN Reference](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference)).

### Core built-ins

| Built-in | Key use | Caveats |
|----------|---------|---------|
| `Object` | Static utilities (`keys`, `values`, `entries`, `defineProperty`) | Prefer static helpers over overridden instance methods |
| `Array` | Ordered, resizable indexed collections | Not associative; copy operations are shallow |
| `Promise` | Async results and composition | Pending → fulfilled or rejected; settled = either |
| `Map` | Keyed collections with insertion-order iteration | Key equality is SameValueZero; objects compared by identity |
| `Set` | Unique-value collections, average sublinear membership checks | Prefer over array scans when uniqueness and membership dominate |

### Async model

- `async function` returns a **new Promise** every call; `await` suspends until fulfillment or rejection, making `try` / `catch` viable around async code ([MDN async function](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/async_function)).
- In Node, sync APIs generally throw; async APIs may reject promises, pass error-first callbacks, or emit `'error'` events depending on the API style ([Node.js Errors](https://nodejs.org/api/errors.html)).

### Modules and packages

- **CommonJS**: `require()` / `module.exports`; each file wrapped as a private module ([Node.js Modules](https://nodejs.org/api/modules.html)).
- **ESM**: `import` / `export`; explicit markers are `.mjs`, `.cjs`, or `package.json#"type"` ([Node.js ESM](https://nodejs.org/api/esm.html)).
- Relative ESM import specifiers **require file extensions**.
- Use the `"exports"` field to control public package API surface ([Node.js Packages](https://nodejs.org/api/packages.html)).

### Node runtime areas to know cold

| Module | Purpose | Key note |
|--------|---------|----------|
| `node:process` | Process info, env, exit lifecycle, permissions | `'exit'` listeners must do only synchronous work |
| `node:events` | `EventEmitter`, synchronous listener dispatch | Listeners run sync in registration order |
| `node:fs` / `node:fs/promises` | Callback, promise, and sync filesystem APIs | Callback form is fastest; sync blocks event loop |
| `node:stream` / `node:stream/promises` | Readable/writable/transform, `pipeline()` | Prefer `pipeline()` over manual pipe chains |
| `node:console` | Logging to stdout/stderr | Sync/async behavior depends on backing stream/platform |
| `node:test` + `node:assert/strict` | Built-in test runner and strict assertions | Don't mix promise-based and `done`-callback in same test |
| `node:perf_hooks` | Performance marks, measures, event-loop utilization | Use `eventLoopUtilization()` for real pressure metrics |

---

## API inventory (condensed)

For exhaustive member lists, use [MDN JavaScript Reference](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference) and [Node.js API index v26.1.0](https://nodejs.org/api/index.html).

### JavaScript core

| API | Method / form | Purpose | Caveats |
|-----|--------------|---------|---------|
| `Promise` | `.then(onFulfilled, onRejected)` | Chain handlers | Returns new promise; handlers run even after settlement |
| `Promise` | `.catch(onRejected)` | Handle rejection | Unhandled rejections throw in Node 15+ |
| `Promise` | `.finally(onFinally)` | Cleanup regardless of outcome | Do not use to transform resolved values |
| `Promise` | `.all(iterable)` | Fail-fast parallel wait — rejects on first rejection | Use `.allSettled()` when you need all results regardless |
| `Promise` | `.allSettled(iterable)` | Wait for all, collect outcomes | Returns `{status, value\|reason}` per promise |
| `Promise` | `.race(iterable)` | Settle with first resolved/rejected | Useful for timeouts paired with a sentinel promise |
| `Promise` | `.any(iterable)` | Resolve with first fulfillment | Throws `AggregateError` only if all reject |
| `async function` | `async function name(...) {}` | Promise-returning functions with `await` | Always returns a promise |
| `Object` | `.keys()` / `.values()` / `.entries()` | Enumerate own enumerable props | Use static helpers, not instance calls |
| `Object` | `.defineProperty(obj, key, descriptor)` | Explicit property definition | Prefer over `__defineGetter__`/`__defineSetter__` |
| `Map` | `.set()` / `.get()` / `.delete()` | Key-value store | SameValueZero equality; object keys by identity |
| `Set` | `.add()` / `.has()` / `.delete()` | Unique-value collection | Prefer over array scans for membership |
| `eval` | `eval(script)` | Evaluate string as script | Avoid unless dynamic evaluation is truly required |
| `AbortController` | `new AbortController()` → `.signal` / `.abort()` | Cancellation token for fetch, streams, and async ops | Pass `.signal` to `fetch`, `pipeline`, or `setTimeout`; check `signal.aborted` in loops |

### Node.js core

| API | Method / form | Purpose | Caveats |
|-----|--------------|---------|---------|
| `node:fs/promises` | `readFile()` / `writeFile()` | Promise-based file I/O | Concurrent mutations of same file are not synchronized |
| `node:fs` | callback APIs (`unlink(path, cb)` etc.) | Error-first callback I/O | Preferred for maximal performance |
| `node:fs` | sync APIs (`readFileSync()` etc.) | Blocking I/O | Blocks the event loop |
| `FileHandle` | `.close()` / `.appendFile()` | Operate on open files | Always close explicitly |
| `EventEmitter` | `.on(event, listener)` | Register listener | Runs sync in registration order |
| `EventEmitter` | `.once(event, listener)` | One-time listener | Good for lifecycle/readiness events |
| `EventEmitter` | `.emit(event, ...args)` | Trigger listeners | Return values from listeners are ignored |
| `node:stream` | `pipeline(source, ...transforms, dest)` | Compose streams with teardown/error forwarding | Prefer over manual pipe chains |
| `node:stream` | `finished(stream, ...)` | Detect stream end/error | Use instead of ad hoc listeners |
| `Readable` | `Readable.from(iterable)` | Readable from iterables | Bridge between iterator and stream code |
| `process` | `.on('exit' \| 'beforeExit' \| ...)` | Observe lifecycle events | `'exit'` handlers must be synchronous |
| `process` | `.exit(code)` / `.exitCode = n` | Control termination | Explicit exit skips async continuation |
| `process.permission` | `.has(scope[, reference])` | Runtime permission check | Only exists when `--permission` flag is set |
| `require()` / `module.exports` | CommonJS load/export | CommonJS modules | Always uses CommonJS loader |
| `import` / `export` | ESM load/export | ES modules | Relative specifiers require file extensions |
| `node:test` | `test(name, fn)` / `t.test(...)` | Define tests and subtests | Don't mix promise-based and `done`-callback |
| `node:assert/strict` | `assert.equal` / `.deepStrictEqual` | Strict assertions | Prefer strict mode; message can be string, Error, or function |
| `node:perf_hooks` | `performance.mark()` / `.measure()` | Performance instrumentation | Clear marks when done to keep timelines manageable |
| `node:perf_hooks` | `performance.eventLoopUtilization()` | Event-loop pressure metric | Use snapshots for interval measurements |

---

## Coding standards

### Error handling

- `try` / `catch` for sync exceptions and `await`-based promise flows.
- Understand the API's error channel before using: sync throw, promise rejection, callback-first error, or emitted `'error'` event ([Node.js Errors](https://nodejs.org/api/errors.html)).
- Match Node errors by **`error.code`**, not message text.

### Async / await

- Prefer `async` / `await` when consuming promise-based APIs — simpler than chaining while preserving `try` / `catch` ergonomics.
- Every async function call returns a new promise, even when returning a non-promise value.
- In tests, use either promise-based async **or** callback-style — not both in the same test.

### Validation and sanitization

- Validate inputs at boundaries; avoid dynamic code evaluation (`eval()`).
- Prefer explicit APIs over implicit coercion-heavy paths.
- `eval()` accepts strings or `TrustedScript` and can be restricted by Trusted Types/CSP.

### Performance

- Callback `fs` APIs for **maximal performance**; promise APIs cost more in allocations.
- Be explicit about module type to avoid Node's syntax-detection fallback.
- Use `Map`/`Set` for workloads with frequent membership checks or uniqueness needs.
- Use `node:perf_hooks` for timing; do not guess.

### Security

- Node's permission model reduces accidental access by trusted code; it **does not protect against malicious code** ([Node.js Permissions](https://nodejs.org/api/permissions.html)).
- Avoid `eval()`; be careful with prototype mutation (prototype pollution risk).

---

## Practical defaults

| Decision | Default |
|----------|---------|
| New Node projects | ES modules (`"type": "module"` or `.mjs`) |
| Async I/O | Promise APIs + `async` / `await` |
| Hot-path file I/O | Callback fs APIs |
| Keyed collections | `Map` over plain-object dictionaries |
| Deduplication / membership | `Set` over arrays |
| Multi-stream flows | `stream.pipeline()` |
| Testing | `node:test` + `node:assert/strict` |
| Performance measurement | `node:perf_hooks` |

---

## Version notes

- **Node version in scope:** v26.1.0 ([Node.js API index v26.1.0](https://nodejs.org/api/index.html)).
- **MDN is a living reference** — confirm runtime support for the target environment before using recently-added features.
- **Permission model** requires `--permission` flag; available permissions are Node-version-sensitive.
- **Module resolution** is version- and package-sensitive: `"type"`, `"exports"`, `"imports"`, file extensions, and syntax detection all affect behavior. Be explicit rather than relying on fallback heuristics.
