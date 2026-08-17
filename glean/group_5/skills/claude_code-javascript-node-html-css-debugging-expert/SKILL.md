---
name: javascript-node-html-css-debugging-expert
description: >-
  JavaScript, Node.js, HTML, and CSS debugging expert: breakpoints (all 8 Chrome DevTools types),
  memory leaks and heap snapshots, async/promise debugging, CPU flame charts, Long Animation
  Frames (LoAF), Core Web Vitals, CSS specificity and layout debugging, network waterfall and
  CORS, source maps, advanced console API, production error monitoring (Sentry, TrackJS),
  diagnostics_channel. TRIGGER: debug JS/Node behavior; browser runtime issues; HTML/CSS
  troubleshooting; memory leaks; CPU/rendering profiling; network/CORS/source-map debugging. SKIP:
  framework-specific debugging (React/Vue DevTools) → the framework skill; TypeScript type errors
  → typescript-expert; pure Vitest test failures → testing-and-vitest-expert; Node CPU/heap
  profiling, heap-snapshot leak hunting, flame graphs (clinic.js/0x), perf_hooks,
  diagnostics_channel/TracingChannel → nodejs-diagnostics-profiling (keeps breakpoints, DevTools,
  HTML/CSS, Web Vitals/LoAF, network).
version: "1.1"
category: developer
updated: "2026-05-29"
tags:
  - javascript
  - nodejs
  - debugging
  - chrome-devtools
  - performance
  - memory
  - css
  - html
  - network
  - source-maps
  - profiling
related_skills:
  - lang-js-ts
  - software-engineering-patterns
  - devops-linux-admin
---

# JavaScript / Node / HTML / CSS Debugging Expert

Practical debugging reference for JavaScript, Node.js, HTML, and CSS. Covers all Chrome DevTools breakpoint types, memory profiling, async diagnostics, CPU flame charts, CSS layout debugging, network/CORS, source maps, and production error monitoring.

## When to use this skill

- Debugging JavaScript or Node.js behavior (logic errors, race conditions, unhandled rejections)
- Investigating browser runtime problems (crashes, hangs, unexpected behavior)
- Troubleshooting HTML or CSS issues (layout, specificity, rendering)
- Choosing between breakpoints, probes, reports, audits, validators, or profilers
- Hunting memory leaks in browser or Node.js applications
- Diagnosing async/promise-related bugs
- Profiling CPU and rendering performance (flame charts, LoAF, Core Web Vitals)
- Debugging network requests and CORS issues
- Working with source maps in production
- Using advanced console API methods
- Setting up production error monitoring (Sentry, TrackJS)
- Using Node.js `diagnostics_channel` for structured observability

## When NOT to use this skill

- Framework-specific debugging (React DevTools, Vue DevTools) — use the framework skill
- TypeScript type errors — use `typescript-expert`
- Test failures in Vitest — use `testing-and-vitest-expert`
- Linux system-level process or memory issues — use `linux-sysadmin`

## Skill guidance

- Start from the decision tree at the end of this file to pick the right tool for each problem type.
- Prefer high-signal techniques: probe mode, diagnostic reports, heap snapshots, performance traces, advanced breakpoint types, and validator/audit stacks.
- For exact commands, flags, and APIs, follow the cited official documentation links throughout this file.

---

## Chrome DevTools Breakpoints -- Complete Reference

Chrome DevTools provides eight breakpoint types for pausing JavaScript execution at precise points.
[Source: developer.chrome.com/docs/devtools/javascript/breakpoints]

### Line-of-code breakpoints
Click a line number in the Sources panel. A blue marker appears. Execution pauses before that line runs. Use when you know the exact region of code to investigate.

### Conditional breakpoints
Right-click a line number > "Add conditional breakpoint" > enter a JS expression. Execution pauses only when the expression evaluates truthy. Eliminates most `if (x) console.log()` patterns, especially useful inside loops.

### Logpoints
Right-click a line number > "Add logpoint" > enter a message template with `{expression}` interpolation. Logs to the Console without modifying source code and without pausing. Use for printf-style debugging without touching the codebase.

### DOM change breakpoints
In the Elements panel, right-click a node > "Break on" > choose one of:
- **Subtree modifications** -- pauses when a child node is added, removed, or moved.
- **Attribute modifications** -- pauses when an attribute on the node changes.
- **Node removal** -- pauses when the node itself is removed.

### XHR/fetch breakpoints
In Sources > XHR/Fetch Breakpoints, click "+" and enter a URL substring. Pauses when any XHR or fetch request URL contains that string. Use to find the code path that triggers a specific API call.

### Event listener breakpoints
In Sources > Event Listener Breakpoints, expand a category (Mouse, Keyboard, Timer, Animation, etc.) and check specific events. Pauses when that event fires. Useful for tracking which handler responds to user input.

### Exception breakpoints
Click the "Pause on exceptions" icon in Sources (octagon with pause). Toggle "Pause on caught exceptions" to also break inside try/catch blocks. Essential for finding swallowed errors.

### Function breakpoints
In the Console, call `debug(functionName)` to set a breakpoint at the first line of that function. Call `undebug(functionName)` to remove it. Works with any function in scope.

### Stepping controls
- **Step over** (F10) -- execute the current line and move to the next line in the same function.
- **Step into** (F11) -- dive into the function call on the current line.
- **Step out** (Shift+F11) -- run the rest of the current function and pause at the caller.
- **Resume** (F8) -- continue execution until the next breakpoint.

### Blackboxing / Ignore list
Settings > Ignore List > add patterns (e.g., `/node_modules/`). After blackboxing, Step Into skips library internals and the call stack only shows your code. Right-click a script in Sources > "Add script to ignore list" for ad-hoc blackboxing.

---

## Node.js Debugging -- Inspector Protocol and Tools

[Sources: nodejs.org/learn/getting-started/debugging, nodejs.org/api/debugger.html, code.visualstudio.com/docs/nodejs/nodejs-debugging]

### The --inspect flag family

| Flag | Behavior |
|------|----------|
| `--inspect` | Activate the V8 inspector on default port 9229. Process runs immediately. |
| `--inspect-brk` | Activate inspector and break before the first line of user code. |
| `--inspect=0.0.0.0:9229` | Bind to all interfaces (use only in trusted networks). |
| `--inspect-brk=PORT` | Break on first line, custom port. |

After starting with `--inspect`, open `chrome://inspect` in Chrome, or attach VS Code with a launch.json configuration:

```json
{
  "name": "Attach to Node",
  "type": "node",
  "request": "attach",
  "port": 9229,
  "continueOnAttach": true
}
```

The `continueOnAttach` option resumes a process launched with `--inspect-brk` as soon as the debugger connects.

### VS Code Auto Attach
Enable Auto Attach via the Command Palette (`Toggle Auto Attach`). In `onlyWithFlag` mode, VS Code automatically attaches to any node process launched in the integrated terminal with `--inspect` or `--inspect-brk`. After enabling, restart the terminal (click the warning icon or create a new one).

### ndb -- Enhanced Node.js debugging
Google ChromeLabs ndb provides an improved debugging experience without needing `--inspect` flags. Install with `npm install -g ndb`, then run `ndb node index.js` or `ndb npm test`. Key advantages: automatic child process debugging, local `node_modules` editing with instant reload, and a standalone DevTools window.
[Source: github.com/GoogleChromeLabs/ndb]

### Node.js Diagnostic Reports
Generate a JSON report containing stack traces, heap statistics, platform info, resource usage, and loaded libraries:

```bash
# Generate report on fatal error (OOM, assertion)
node --report-on-fatalerror app.js

# Generate report on unhandled exception
node --report-uncaught-exception app.js

# Generate report on signal (default SIGUSR2)
node --report-on-signal app.js

# Generate report programmatically
process.report.writeReport('./report.json');
```

The report includes: JavaScript and native stack traces, heap statistics, system information, resource usage, and libuv handles. Use `--report-directory` to control output location.

### Increasing heap memory
When hitting "JavaScript heap out of memory" (FATAL ERROR: Reached heap limit):

```bash
# Increase to 4GB
node --max-old-space-size=4096 app.js

# Via environment variable
NODE_OPTIONS="--max-old-space-size=4096" node app.js
```

Default V8 heap: ~1.5 GB on 64-bit systems. For production, investigate the root cause rather than just increasing the limit.

---

## Node.js diagnostics_channel API

The `node:diagnostics_channel` module (Stable since Node.js 20+) provides a low-overhead publish/subscribe mechanism for structured diagnostic events without monkey-patching.
[Sources: nodejs.org/api/diagnostics_channel.html, sentry.engineering/blog/from-monkey-patching-to-tracing-channels]

### Channel API basics

```js
import { channel } from 'node:diagnostics_channel';

const ch = channel('my-app:request');

// Publisher
ch.publish({ url: '/api/data', method: 'GET' });

// Subscriber
ch.subscribe((message, name) => {
  console.log(`[${name}]`, message);
});
```

### TracingChannel -- lifecycle instrumentation
TracingChannel wraps five channels for a single traceable action: `start`, `end`, `asyncStart`, `asyncEnd`, `error`.

```js
import { tracingChannel } from 'node:diagnostics_channel';

const tracing = tracingChannel('my-app:db-query');

// Subscribe to lifecycle events
tracing.subscribe({
  start(message)      { /* query begins */ },
  end(message)        { /* sync portion completes */ },
  asyncStart(message) { /* async callback begins */ },
  asyncEnd(message)   { /* async callback completes */ },
  error(message)      { /* error thrown */ }
});
```

### Built-in HTTP channels
Since Node.js 20.12+, undici (the built-in HTTP client) emits diagnostics_channel events. Additional framework support:
- **undici**: `undici:request:create`, `undici:request:headers`, `undici:request:trailers`
- **fastify**: Native diagnostics_channel support
- **nitro/h3**: Native tracing channel support
- **mysql2**: TracingChannel support
- **ioredis / node-redis**: TracingChannel support

This replaces the old monkey-patching approach for observability libraries like Sentry, Datadog, and New Relic.

---

## Memory Leak Detection and Heap Analysis

[Sources: developer.chrome.com/docs/devtools/memory-problems, nodejs.org/learn/diagnostics/memory/using-heap-snapshot, oneuptime.com/blog/post/2026-01-26-nodejs-memory-leak-profiling, bun.com/blog/debugging-memory-leaks]

> **GC mechanics & tuning vs. leak hunting.** This section covers *finding* leaks (heap snapshots, retainer
> chains, DevTools workflow). For the underlying V8 garbage collector — Orinoco generational GC (Scavenger
> vs Mark-Compact), why memory grows or pauses spike, and the Node heap-sizing flags (`--max-old-space-size`,
> `--max-semi-space-size`, `--trace-gc`) — load the `references/v8-engine-internals.md` hub reference. Rule of
> thumb: if RSS climbs forever it's a *leak* (use this section); if GC is just expensive or pauses are long
> it's a *tuning* problem (use v8-engine-internals). The same file also covers hidden classes / inline caches
> and the JIT tiering pipeline (Ignition/Sparkplug/Maglev/TurboFan) behind `--trace-opt`/`--trace-deopt`.

### Symptom recognition
A memory leak shows continuously growing memory consumption that never plateaus even when workload stays constant. Signs include:
- Progressively slower page performance over time
- Chrome Task Manager showing growing "JavaScript Memory" column
- Node.js process RSS growing without bound in `process.memoryUsage()`

### Chrome DevTools Memory panel workflow

**1. Heap Snapshot comparison (primary technique)**
1. Open DevTools > Memory tab > select "Heap snapshot"
2. Take a snapshot (baseline)
3. Perform the suspected leaking operation several times
4. Take a second snapshot
5. Select the second snapshot, switch view from "Summary" to "Comparison"
6. Sort by "Delta" or "Alloc. Size" to find growing object types
7. Drill into retained paths to find what holds the reference

**2. Allocation Timeline**
1. Memory tab > select "Allocation instrumentation on timeline"
2. Click Start, interact with the app, click Stop
3. Blue bars = allocations; gray bars = freed memory
4. Persistent blue bars that never turn gray indicate retained objects
5. Click a blue bar to see the object and its retaining tree

**3. Allocation Sampling**
Lower overhead than timeline. Good for production-like profiling. Shows which functions allocate the most memory.

### Common JavaScript memory leak patterns

| Pattern | Description | Fix |
|---------|-------------|-----|
| Detached DOM trees | Removed from document but held by JS variable | Null the reference after removal |
| Forgotten event listeners | addEventListener without removeEventListener | Clean up in cleanup/unmount lifecycle |
| Closures over large scopes | Inner function captures entire outer scope | Extract only needed variables |
| Global variables | Accidental globals via missing `const`/`let` | Use strict mode, ESLint no-undef |
| Uncleared timers | setInterval never cleared | Store interval ID, clearInterval on cleanup |
| Growing Maps/Sets/Arrays | Collections that only grow, never prune | Implement size limits or use WeakMap/WeakSet |
| Console references | `console.log(bigObject)` retains reference in DevTools | Remove console.log in production |

### WeakRef and FinalizationRegistry (ES2021+)

```js
// WeakRef -- holds reference without preventing GC
const cache = new Map();
function getCached(key, compute) {
  const ref = cache.get(key);
  if (ref) {
    const val = ref.deref();  // returns undefined if GC'd
    if (val !== undefined) return val;
  }
  const value = compute();
  cache.set(key, new WeakRef(value));
  return value;
}

// FinalizationRegistry -- cleanup callback when object is GC'd
const registry = new FinalizationRegistry((heldValue) => {
  console.log(`Object for ${heldValue} was garbage collected`);
  // Clean up external resources (file handles, connections)
});
registry.register(someObject, 'resource-id');
```

**Caveats**: GC timing is non-deterministic. Never rely on WeakRef/FinalizationRegistry for critical logic. Use them for caches and resource cleanup hints only.
[Sources: developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/WeakRef, v8.dev/features/weak-references]

### Node.js heap snapshots in production

```js
// Using v8 module (built-in)
const v8 = require('v8');
const fs = require('fs');

const snapshotFile = v8.writeHeapSnapshot();
console.log(`Heap snapshot written to ${snapshotFile}`);

// Using heapdump module
const heapdump = require('heapdump');
heapdump.writeSnapshot('/tmp/heap-' + Date.now() + '.heapsnapshot');
```

Open `.heapsnapshot` files in Chrome DevTools Memory panel for analysis.

---

## Async Debugging

[Sources: draconianoverlord.com/2025/04/17/fixing-async-stack-traces, dev.to/alex_aslam/tackling-asynchronous-bugs-in-javascript]

### Async stack traces in Chrome DevTools
Chrome DevTools captures async stack traces by default. In the Sources panel, the Call Stack shows the full async chain (e.g., `setTimeout` caller, `Promise.then` originator, `fetch` initiator). Ensure "Async" checkbox is enabled in the Call Stack section.

### Common async bugs and fixes

**1. Unhandled promise rejections**
```js
// BAD -- rejection silently swallowed
fetchData();

// GOOD -- always handle rejections
fetchData().catch(handleError);

// Or use await with try/catch
try {
  await fetchData();
} catch (err) {
  handleError(err);
}
```

In Node.js, listen for unhandled rejections:
```js
process.on('unhandledRejection', (reason, promise) => {
  console.error('Unhandled Rejection at:', promise, 'reason:', reason);
  // In production: log to error tracker, exit gracefully
});
```

Since Node.js 15+, unhandled rejections throw by default (`--unhandled-rejections=throw`).

**2. Race conditions**
```js
// BAD -- last call wins, may show stale data
input.addEventListener('input', async (e) => {
  const results = await search(e.target.value);
  showResults(results);
});

// GOOD -- abort previous request
let controller;
input.addEventListener('input', async (e) => {
  controller?.abort();
  controller = new AbortController();
  try {
    const results = await search(e.target.value, {
      signal: controller.signal
    });
    showResults(results);
  } catch (e) {
    if (e.name !== 'AbortError') throw e;
  }
});
```

**3. Floating promises (ESLint detection)**
Use `@typescript-eslint/no-floating-promises` to catch promises that are neither awaited nor returned nor caught. This prevents silent failures where errors disappear.

**4. processTicksAndRejections in Node.js stack traces**
The presence of `processTicksAndRejections` in a stack trace indicates that the Promise rejection was not handled in the same microtask and Node deferred handling to the next event loop tick. This usually means a missing `await` or `.catch()`.

---

## Performance Profiling

[Sources: developer.chrome.com/docs/devtools/performance, developer.chrome.com/blog/profiling-cpu, nodejs.org/learn/diagnostics/flame-graphs, blog.platformatic.dev/introducing-next-gen-flamegraphs-for-nodejs]

### Chrome DevTools Performance panel
1. Open DevTools > Performance tab
2. Click Record, interact with the page, click Stop
3. Analyze the flame chart, summary, bottom-up, call tree, and event log tabs

Key sections of a Performance recording:
- **Frames** -- visual frame rate, red frames indicate jank
- **Main thread** -- flame chart showing JS execution, layout, paint, and composite
- **Network** -- request waterfall during the recording
- **Timings** -- User Timing API marks and measures, plus Web Vitals

### CPU flame charts
The flame chart displays function execution over time. The horizontal axis is time; each row in the stack is a function call. Wider bars = more time spent in that function. Look for:
- **Tall stacks** -- deeply nested calls
- **Wide bars** -- expensive individual functions
- **Repeated patterns** -- functions called too frequently

### Node.js flame graphs
```bash
# Generate a CPU profile
node --prof app.js
# Process the log
node --prof-process isolate-*.log > profile.txt

# OR use 0x for flame graphs
npx 0x app.js
# Opens interactive flame graph in browser

# OR use clinic.js
npx clinic flame -- node app.js
```

### Long Animation Frames API (LoAF)
Replaces the Long Tasks API with richer data. A long animation frame is any rendering update delayed beyond 50ms. Available since Chrome 123.
[Sources: developer.chrome.com/docs/web-platform/long-animation-frames, developer.mozilla.org/en-US/docs/Web/API/Performance_API/Long_animation_frame_timing, requestmetrics.com/web-performance/long-animation-frame-loaf]

```js
// Observe long animation frames
const observer = new PerformanceObserver((list) => {
  for (const entry of list.getEntries()) {
    console.log('LoAF duration:', entry.duration, 'ms');
    console.log('Block duration:', entry.blockingDuration, 'ms');
    console.log('Render start:', entry.renderStart);
    // Inspect scripts that contributed
    for (const script of entry.scripts) {
      console.log('  Script:', script.sourceURL);
      console.log('  Function:', script.sourceFunctionName);
      console.log('  Duration:', script.duration, 'ms');
      console.log('  Invoker:', script.invoker);
      console.log('  Type:', script.invokerType);
    }
  }
});
observer.observe({ type: 'long-animation-frame', buffered: true });
```

**PerformanceLongAnimationFrameTiming** properties: `duration`, `blockingDuration`, `renderStart`, `styleAndLayoutStart`, `firstUIEventTimestamp`, `scripts[]`.

**PerformanceScriptTiming** properties: `sourceURL`, `sourceFunctionName`, `sourceCharPosition`, `invoker`, `invokerType` (user-callback, event-listener, resolve-promise, script-block), `duration`, `executionStart`, `forcedStyleAndLayoutDuration`, `windowAttribution`.

LoAF helps diagnose poor INP (Interaction to Next Paint) scores by identifying exactly which scripts cause rendering delays.

### Core Web Vitals in DevTools
Since Chrome 132 (January 2025), the Performance panel defaults to showing live Core Web Vitals metrics:
[Source: debugbear.com/blog/2025-in-web-performance]

| Metric | Target | What it measures |
|--------|--------|-----------------|
| **LCP** (Largest Contentful Paint) | < 2.5s | How fast main content loads |
| **INP** (Interaction to Next Paint) | < 200ms | How quickly page responds to interactions |
| **CLS** (Cumulative Layout Shift) | < 0.1 | Visual stability during load |

DevTools surfaces Lighthouse-style performance insights inline, including LCP subpart analysis, request discovery issues, third-party script impact, and legacy/duplicated JavaScript detection.

Since December 2025, you can throttle individual network requests by right-clicking in the Network panel instead of throttling the entire page.

---

## CSS Debugging

[Sources: developer.chrome.com/docs/devtools/css/reference, developer.chrome.com/docs/devtools/css/issues, developer.mozilla.org/en-US/docs/Learn_web_development/Core/Styling_basics/Debugging_CSS, css-snacks.com/guide/advanced-css-specificity]

### Styles panel
Select an element in the Elements panel to see all CSS rules applied in the Styles tab. Rules are listed by specificity order. Crossed-out declarations are overridden by higher-specificity or later-declared rules.

**Filter box**: Type a property name (e.g., "margin") to filter all rules to those affecting that property.

**Specificity tooltip**: Hover over any selector to see its specificity weight displayed as a three-part tuple (e.g., `(0, 1, 2)`).

### Computed panel
Shows the final resolved CSS values after all inheritance, cascading, and specificity rules are applied. Key uses:
- See the actual pixel value for relative units (`30vw` resolved to `786px`)
- Click the arrow next to any computed value to jump to the declaration that set it
- Check the "Show all" checkbox to see inherited properties

### Box model diagram
Visual representation of margin, border, padding, and content dimensions. Click any number in the diagram to edit it live.

### CSS issue detection
DevTools marks CSS issues with icons:
- **Inactive declarations** (grayed, with info icon) -- property has no effect in current context (e.g., `width` on an inline element)
- **Overridden declarations** (crossed out) -- another rule with higher specificity or later in cascade wins
- **Invalid values** (yellow warning) -- the property value is not recognized
- **Non-inherited properties shown on child** -- the property won't inherit from the parent

### Flexbox and Grid debugging
- **Flexbox overlay**: Click the `flex` badge on a flex container in Elements to toggle the overlay showing flex lines, item boundaries, and available space.
- **Grid overlay**: Click the `grid` badge to see grid lines, track sizes, area names, and gap spacing.
- **Grid editor**: In the Styles panel, click the grid icon next to `grid-template-columns` or `grid-template-rows` to visually edit tracks.

### CSS specificity debugging workflow
1. Select the element showing unexpected styles
2. In the Styles panel, find the expected rule -- if it is crossed out, another rule overrides it
3. Hover over both selectors to compare specificity weights
4. Use the Computed tab to confirm which declaration wins
5. Fix by increasing specificity, using `!important` (last resort), or restructuring selectors

---

## Network Debugging

[Sources: developer.chrome.com/docs/devtools/network/reference, httptoolkit.com/blog/how-to-debug-cors-errors, headersnap.com/blog/debug-cors-errors-chrome, knowledgelib.io/software/debugging/browser-cors-errors/2026]

### Network panel waterfall
The waterfall shows request timing broken into phases:
- **Queueing** -- waiting for an available connection
- **Stalled** -- waiting after queueing
- **DNS Lookup** -- resolving the domain
- **Initial connection / SSL** -- TCP handshake, TLS negotiation
- **TTFB (Time to First Byte)** -- waiting for server response
- **Content Download** -- receiving the response body

Sort by total duration. The lighter portion represents waiting; the darker portion represents downloading.

### Throttling
- **Page-level**: Network tab > throttle dropdown (Fast 3G, Slow 3G, Offline)
- **Per-request** (Chrome 132+): Right-click a request > "Throttle request" to simulate slow responses for a single resource without affecting the rest

### CORS debugging workflow

**Step 1 -- Identify the error**: Console shows `Access-Control-Allow-Origin` errors. Note the origin, method, and whether it is a preflight failure.

**Step 2 -- Find the preflight**: Filter Network requests by "OPTIONS" method. In Chrome 79+, preflight requests may be hidden; enable "Show all network requests" in DevTools settings.

**Step 3 -- Inspect headers on the failed request**:
- **Request headers**: Verify `Origin` header is present and correct
- **Response headers**: Check for:
  - `Access-Control-Allow-Origin` -- must match the request origin or be `*`
  - `Access-Control-Allow-Methods` -- must include the HTTP method used
  - `Access-Control-Allow-Headers` -- must include any custom headers sent
  - `Access-Control-Allow-Credentials` -- must be `true` if cookies are sent

**Step 4 -- Common fixes**:
- Server not returning CORS headers: Configure the server middleware
- Wildcard `*` with credentials: Not allowed; must specify exact origin
- Missing preflight response: Server must respond to OPTIONS with 204 and the correct CORS headers
- Mixed content: HTTP API called from HTTPS page

### Request blocking
Network panel > right-click any request > "Block request URL" or "Block request domain". Useful for testing how the page behaves when a third-party script or API is unavailable.

---

## Source Maps

[Sources: developer.chrome.com/docs/devtools/javascript/source-maps, thisdot.co/blog/understanding-sourcemaps-from-development-to-production, polarsignals.com/blog/posts/2025/11/04/javascript-source-maps-internals, hardcoreprawn.github.io/tech-content-curator/posts/2025-11-04-mastering-source-maps]

### Enabling source maps in DevTools
Settings > Preferences > Sources > check "JavaScript source maps" and "CSS source maps". With source maps enabled, you see original TypeScript/JSX/SCSS files in the Sources panel and can set breakpoints in them directly.

### Source map types

| Type | Description | Use case |
|------|-------------|----------|
| **Inline** | Base64-encoded map embedded in the JS file via `//# sourceMappingURL=data:...` | Small projects, dev only |
| **External** | Separate `.map` file referenced by `//# sourceMappingURL=file.js.map` | Standard production + debugging |
| **Hidden** | External `.map` file with no reference in the JS bundle | Production with server-side error tracking (Sentry) |

### Build tool configuration

**Webpack**:
```js
// Development -- fast rebuilds, good quality
devtool: 'eval-source-map'

// Production -- full quality, separate file
devtool: 'source-map'

// Production (hidden) -- maps exist but not referenced
devtool: 'hidden-source-map'
```

**Vite** (esbuild/Rollup):
```js
build: {
  sourcemap: true,       // external .map files
  // sourcemap: 'hidden' // hidden source maps
}
```

### Security considerations
- Never serve source maps publicly if you want to protect proprietary code
- Use hidden source maps + error tracking services (Sentry, Datadog) that apply maps server-side
- Upload source maps to error tracking during CI/CD, then delete from the deployment artifact

---

## Advanced Console API

[Sources: developer.mozilla.org/en-US/docs/Web/API/console, developer.chrome.com/docs/devtools/console/api]

### Beyond console.log

| Method | Purpose | Example |
|--------|---------|---------|
| `console.table(data, columns?)` | Display array/object as sortable table | `console.table(users, ['name', 'email'])` |
| `console.dir(obj, {depth})` | Interactive object tree (especially DOM nodes) | `console.dir(document.body, {depth: 2})` |
| `console.trace(label?)` | Print stack trace from call site | `console.trace('called from')` |
| `console.group(label)` | Start collapsible group | Nest related logs |
| `console.groupCollapsed(label)` | Start collapsed group | Reduce log noise |
| `console.groupEnd()` | End current group | |
| `console.time(label)` | Start timer | `console.time('fetch')` |
| `console.timeEnd(label)` | Stop timer, print duration | `console.timeEnd('fetch')` |
| `console.timeLog(label)` | Print elapsed without stopping | Intermediate timing |
| `console.count(label)` | Count how many times called | `console.count('render')` |
| `console.countReset(label)` | Reset counter | |
| `console.assert(cond, msg)` | Log error only when condition is false | `console.assert(x > 0, 'x must be positive')` |
| `console.warn(msg)` | Yellow warning-level log | |
| `console.error(msg)` | Red error-level log with stack trace | |
| `console.clear()` | Clear the console | |

### Console utilities (DevTools only, not available in scripts)
- `$0` -- currently selected element in Elements panel
- `$_` -- result of the last evaluated expression
- `$('selector')` -- shortcut for `document.querySelector`
- `$$('selector')` -- shortcut for `document.querySelectorAll`
- `copy(value)` -- copy any value to clipboard
- `monitor(fn)` -- log every call to a function with arguments
- `monitorEvents(el, events?)` -- log DOM events on an element
- `getEventListeners(el)` -- list all event listeners on an element
- `queryObjects(Constructor)` -- find all instances of a constructor in the heap

### String substitution in console.log
```js
console.log('User %s has %d items (%f MB)', name, count, size);
console.log('%cStyled text', 'color: blue; font-weight: bold');
console.log('%o', domElement);  // interactive DOM element
console.log('%O', jsObject);    // interactive JS object
```

---

## HTML Validation and Accessibility Auditing

[Sources: rocketvalidator.com/axe-core, developer.chrome.com/docs/lighthouse/accessibility/scoring, screamingfrog.co.uk/seo-spider/tutorials/how-to-perform-a-web-accessibility-audit]

### HTML validation tools
- **W3C Nu Validator** (validator.w3.org) -- checks HTML conformance to the spec. Catches unclosed tags, invalid nesting, deprecated attributes, and missing required attributes.
- **DevTools Elements panel** -- highlights parser errors with red underline in the DOM tree.
- **HTMLHint** / **html-validate** -- CLI/CI linters for HTML best practices.

### Accessibility testing stack

| Tool | Scope | Coverage |
|------|-------|----------|
| **Lighthouse** | Browser/CI | ~56 axe rules, performance + SEO + PWA bundled |
| **axe-core** | Library/CI | 80+ WCAG 2.1 rules, ~57% of WCAG issues automatable |
| **WAVE** | Browser extension | Visual overlay of accessibility issues |
| **axe DevTools** | Browser extension | Full axe rule set + guided manual tests |
| **pa11y** | CLI/CI | Automated WCAG/Section 508 checks |

### Lighthouse accessibility audit workflow
1. DevTools > Lighthouse tab > check "Accessibility" > Analyze
2. Review the score (0-100) and individual audit results
3. Each finding links to the axe rule documentation with fix guidance
4. For issues that require human judgment (43% of WCAG criteria): test keyboard navigation, screen reader output, logical focus order, and meaningful alt text

### WCAG 2.1 Level AA quick checklist
- All images have descriptive `alt` text (or empty `alt=""` for decorative images)
- Color contrast ratio >= 4.5:1 for normal text, >= 3:1 for large text
- All interactive elements are keyboard accessible
- Focus indicators are visible
- Page has a logical heading hierarchy (h1 > h2 > h3)
- Form inputs have associated `<label>` elements
- ARIA roles and properties are valid and complete

---

## Production Error Monitoring

[Sources: sentry.io/for/javascript, trackjs.com, docs.sentry.io/platforms/javascript, inspectlet.com/guides/best-javascript-error-tracking-tools]

### Error capture listeners
```js
// Global error handler (synchronous errors)
window.addEventListener('error', (event) => {
  reportError({
    message: event.message,
    source: event.filename,
    line: event.lineno,
    column: event.colno,
    stack: event.error?.stack
  });
});

// Unhandled promise rejections
window.addEventListener('unhandledrejection', (event) => {
  reportError({
    message: event.reason?.message || String(event.reason),
    stack: event.reason?.stack
  });
});

// Resource load errors (images, scripts, stylesheets)
window.addEventListener('error', (event) => {
  if (event.target !== window) {
    reportResourceError(event.target.src || event.target.href);
  }
}, true);  // capture phase to catch resource errors
```

### Sentry integration
```js
import * as Sentry from '@sentry/browser';

Sentry.init({
  dsn: 'https://key@sentry.io/project',
  integrations: [
    Sentry.browserTracingIntegration(),
    Sentry.replayIntegration()
  ],
  tracesSampleRate: 0.1,        // 10% of transactions
  replaysSessionSampleRate: 0.01, // 1% of sessions
  replaysOnErrorSampleRate: 1.0,  // 100% of error sessions
});
```

Features: source map upload (via `@sentry/cli` or `@sentry/webpack-plugin`), session replay, breadcrumbs trail, release tracking, AI-powered error analysis ("Seer").

### TrackJS integration
Lightweight (<10KB gzipped) error monitoring focused on client-side JavaScript. Provides a timeline view of user actions leading up to each error, giving behavioral context beyond raw stack traces.

### Key production debugging patterns
1. **Source map upload**: Upload maps during CI/CD to error tracking service, remove from deployment
2. **Release tagging**: Tag each deployment with a version/commit hash for correlation
3. **Breadcrumbs**: Automatically capture console logs, network requests, and DOM interactions before an error
4. **Session replay**: Video-like reproduction of the user session for visual context
5. **Alert tuning**: Configure ignore rules to filter noise from third-party scripts and browser extensions

---

## Debugging Decision Tree

Use this to pick the right tool for the job:

```
What are you debugging?
|
+-- JavaScript logic error?
|   +-- Know the file/line? --> Line-of-code breakpoint
|   +-- Know the function? --> debug(fn) in Console
|   +-- Related to a DOM change? --> DOM breakpoint
|   +-- Related to a network request? --> XHR/fetch breakpoint
|   +-- Related to user input? --> Event listener breakpoint
|   +-- Intermittent / conditional? --> Conditional breakpoint + logpoint
|
+-- Memory issue?
|   +-- Browser? --> Heap snapshot comparison + allocation timeline
|   +-- Node.js? --> --inspect + heap snapshot via v8 module
|   +-- OOM crash? --> --report-on-fatalerror + --max-old-space-size
|
+-- Performance issue?
|   +-- Browser rendering? --> Performance panel flame chart
|   +-- Long tasks / janky UI? --> Long Animation Frames API (LoAF)
|   +-- Core Web Vitals? --> Performance panel live metrics + Lighthouse
|   +-- Node.js CPU? --> --prof + 0x / clinic flame
|
+-- CSS styling issue?
|   +-- Wrong value applied? --> Computed tab
|   +-- Overridden rule? --> Styles tab (look for crossed-out)
|   +-- Layout problem? --> Flexbox/Grid overlays + box model
|   +-- Specificity conflict? --> Hover selectors for specificity weights
|
+-- Network issue?
|   +-- Slow requests? --> Network waterfall timing breakdown
|   +-- CORS errors? --> Console error + OPTIONS request inspection
|   +-- Missing requests? --> Request blocking test
|
+-- Async / Promise issue?
|   +-- Swallowed error? --> Pause on exceptions (including caught)
|   +-- Race condition? --> AbortController pattern
|   +-- Unhandled rejection? --> process.on('unhandledRejection') or window event
|
+-- Production error?
    +-- Need stack traces? --> Source maps + Sentry/TrackJS
    +-- Need reproduction? --> Session replay
    +-- Need Node.js diagnostics? --> Diagnostic reports + diagnostics_channel
```

---

## Sources

1. [Chrome DevTools: Pause your code with breakpoints](https://developer.chrome.com/docs/devtools/javascript/breakpoints)
2. [Chrome DevTools: Debug JavaScript](https://developer.chrome.com/docs/devtools/javascript)
3. [Node.js: Debugging Guide](https://nodejs.org/learn/getting-started/debugging)
4. [Node.js: diagnostics_channel API](https://nodejs.org/api/diagnostics_channel.html)
5. [Node.js: Using Heap Snapshot](https://nodejs.org/learn/diagnostics/memory/using-heap-snapshot)
6. [Chrome DevTools: Fix memory problems](https://developer.chrome.com/docs/devtools/memory-problems)
7. [Chrome DevTools: CSS features reference](https://developer.chrome.com/docs/devtools/css/reference)
8. [Chrome DevTools: Find CSS issues](https://developer.chrome.com/docs/devtools/css/issues)
9. [Chrome DevTools: Long Animation Frames API](https://developer.chrome.com/docs/web-platform/long-animation-frames)
10. [MDN: Long animation frame timing](https://developer.mozilla.org/en-US/docs/Web/API/Performance_API/Long_animation_frame_timing)
11. [MDN: Console API](https://developer.mozilla.org/en-US/docs/Web/API/console)
12. [Chrome DevTools: Console API reference](https://developer.chrome.com/docs/devtools/console/api)
13. [Chrome DevTools: Source maps](https://developer.chrome.com/docs/devtools/javascript/source-maps)
14. [Sentry: From Monkey-Patching to Tracing Channels](https://sentry.engineering/blog/from-monkey-patching-to-tracing-channels)
15. [VS Code: Node.js Debugging](https://code.visualstudio.com/docs/nodejs/nodejs-debugging)
16. [MDN: WeakRef](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/WeakRef)
17. [V8: Weak references and finalizers](https://v8.dev/features/weak-references)
18. [DebugBear: 2025 in Web Performance](https://www.debugbear.com/blog/2025-in-web-performance)
19. [Chrome DevTools: Performance panel](https://developer.chrome.com/docs/devtools/performance)
20. [Node.js: Flame Graphs](https://nodejs.org/learn/diagnostics/flame-graphs)
21. [Platformatic: Next-Gen Flamegraphs for Node.js](https://blog.platformatic.dev/introducing-next-gen-flamegraphs-for-nodejs)
22. [Sentry: JavaScript Error Monitoring](https://sentry.io/for/javascript/)
23. [TrackJS: JavaScript Error Monitoring](https://trackjs.com/)
24. [HTTP Toolkit: How to Debug CORS Errors](https://httptoolkit.com/blog/how-to-debug-cors-errors/)
25. [MDN: Debugging CSS](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Styling_basics/Debugging_CSS)
26. [Fixing Async Stack Traces (2025)](https://www.draconianoverlord.com/2025/04/17/fixing-async-stack-traces.html/)
