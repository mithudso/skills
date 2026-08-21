<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `performance-profiling-expert` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: performance-profiling-expert
version: 1.1.0
category: developer
tags: [performance, profiling, chrome-devtools, lighthouse, web-vitals, node-js, perf-hooks, runtime, browser]
description: >-
  Performance and profiling expert — Chrome DevTools Performance panel, Lighthouse audits,
  Core Web Vitals, Node.js perf_hooks (marks, measures, PerformanceObserver, eventLoopUtilization),
  RAIL model, CPU throttling, heap snapshots, and regression monitoring.
  TRIGGER: slow page loads, runtime jank, high INP/LCP/CLS scores, Node.js event-loop saturation,
  profiling a Chrome extension service worker, Lighthouse audit interpretation, or setting up
  performance measurement in a Node app.
  SKIP: MongoDB query performance (use mongodb-query-performance); general debugging not involving
  timing or profiling; infrastructure scaling.
triggers:
  - "Page feels janky — how do I profile it?"
  - "Lighthouse score is low — where do I start?"
  - "How do I measure LCP and INP?"
  - "Profile a Node.js service for slow response times"
  - "Use perf_hooks marks and measures in Node"
  - "Event loop utilization is high in my Node app"
  - "How do I use CPU throttling in Chrome DevTools?"
  - "Interpret a Chrome performance trace"
  - "Profile a Chrome extension service worker"
  - "Set up PerformanceObserver in Node"
related_skills:
  - chrome-dev
  - nodejs-observability
  - mongodb-query-performance
sources:
  - https://developer.chrome.com/docs/devtools/performance
  - https://developer.chrome.com/docs/lighthouse
  - https://web.dev/performance
  - https://developer.mozilla.org/en-US/docs/Web/Performance
  - https://nodejs.org/api/perf_hooks.html
---

# Performance and Profiling Expert

Performance work has two distinct problem classes — **load performance** (how fast a page becomes usable) and **runtime performance** (how smoothly it runs afterward). Use the right tool for each; they have different bottlenecks and different fixes.

**When not to use this skill:** MongoDB query optimization (use `mongodb-query-performance`), infrastructure capacity planning, or general debugging unrelated to timing and profiling.

## Quick triage rules

1. **Measure before changing** — you cannot improve performance without first measuring it.
2. Distinguish **load** from **runtime**: Chrome DevTools defines runtime performance as how the page performs while running, not while loading.
3. Use **CPU throttling** when profiling browser runtime so desktop measurements reflect mobile constraints.
4. Treat **Lighthouse as a prioritization aid**, not a final verdict — failing audits indicate where to investigate, not the root cause.
5. For Node.js, use `performance.mark()` / `measure()` and `PerformanceObserver` from `node:perf_hooks` rather than `Date.now()` timing.
6. Clear marks, measures, and resource timings between profiling sessions to avoid polluted timelines.
7. Treat performance as **ongoing monitoring**, not a one-time sprint.

## Diagnostic workflow

| Step | Action | Tool |
|------|--------|------|
| 1 | Identify symptom: loading, jank, input latency, backend | web.dev, MDN |
| 2 | Measure with the right tool (see table below) | Lighthouse / DevTools / perf_hooks |
| 3 | Profile in realistic conditions (CPU throttling, real network) | DevTools |
| 4 | Interpret trace or metrics to isolate the bottleneck | DevTools flame chart |
| 5 | Fix one bottleneck class at a time; re-measure to confirm | — |

## Tool selection

| Symptom | Start with | Then use |
|---------|-----------|----------|
| Slow initial page load / poor Core Web Vitals | Lighthouse / PageSpeed Insights | DevTools Performance (load trace) |
| Jank, dropped frames, slow interactions | DevTools Performance (runtime trace) | Flame chart, long-task analysis |
| Node.js slow response or high CPU | `perf_hooks` marks + measures | `PerformanceObserver`, event loop utilization |
| Memory growth / leak | DevTools Memory → Heap Snapshot | Comparison view, detached DOM nodes |
| Slow Chrome extension SW | DevTools → "Inspect views: service worker" | Flame chart for cold-start and event handlers |

## Browser profiling

### Runtime profiling workflow

1. Open DevTools in **Incognito** (reduces extension noise).
2. Performance panel → gear icon → set **CPU throttling** (4x for desktop-as-mobile approximation).
3. Click **Record**, reproduce the problematic interaction, click **Stop**.
4. Flame chart: look for long tasks (red corners) in the **Main** thread.
5. Focus on tasks > 50ms — the RAIL Response threshold for user-perceived latency.

### Timing thresholds (MDN)

| Metric | Threshold | Meaning |
|--------|-----------|---------|
| Input response | < 50ms | User perceives as instant |
| Animation frame | ≤ 16.7ms | 60 fps |
| Content load signal | ~1 second | Page feels loaded |

### Core Web Vitals (web.dev)

| Vital | Measures | Good threshold |
|-------|---------|----------------|
| LCP (Largest Contentful Paint) | Load performance | ≤ 2.5s |
| INP (Interaction to Next Paint) | Responsiveness | ≤ 200ms |
| CLS (Cumulative Layout Shift) | Visual stability | ≤ 0.1 |

### Lighthouse

```bash
# CLI
npx lighthouse https://example.com --output html --output-path report.html

# In DevTools: Lighthouse tab → choose categories → Analyze page load
```

Treat failing audits as indicators. Each audit links to a reference document explaining the fix. PageSpeed Insights runs Lighthouse against real-user data (field data) in addition to lab data.

## Node.js perf_hooks

### Basic marks and measures

```javascript
import { performance, PerformanceObserver } from 'node:perf_hooks';

// Mark start
performance.mark('db-query-start');

// ... async work ...

// Mark end and measure
performance.mark('db-query-end');
performance.measure('db-query', 'db-query-start', 'db-query-end');

// Read entries
const [entry] = performance.getEntriesByName('db-query');
console.log(`db-query: ${entry.duration.toFixed(2)}ms`);

// Clean up between sessions to avoid polluted timelines
performance.clearMarks();
performance.clearMeasures();
```

### PerformanceObserver (streaming)

```javascript
const obs = new PerformanceObserver((list) => {
  for (const entry of list.getEntries()) {
    console.log(`${entry.name}: ${entry.duration.toFixed(2)}ms`);
  }
});
obs.observe({ entryTypes: ['measure'] });

// ... run measured code ...

obs.disconnect();  // stop observing when done
```

### Event loop utilization

```javascript
// Snapshot 1 — take at process start or after warm-up
const elu1 = performance.eventLoopUtilization();

// ... after a load period ...

// Snapshot 2
const elu2 = performance.eventLoopUtilization(elu1);
console.log(`ELU: ${(elu2.utilization * 100).toFixed(1)}%`);
// > 80% sustained = event loop saturation risk
```

`eventLoopUtilization()` is Node-specific (not available in browsers). High utilization (> 80% sustained) indicates CPU-bound work blocking I/O responses.

## Memory profiling (browser)

### Heap snapshot workflow

1. DevTools → **Memory** panel → **Heap Snapshot** → **Take snapshot** (baseline).
2. Perform the suspected-leak operation N times.
3. Take a second snapshot.
4. Switch to **Comparison** view → sort by **# Delta** → positive delta = potential leak.
5. Filter by constructor (e.g., `Detached HTMLElement`) to find retained DOM nodes.

### Common extension memory leak patterns

| Pattern | Fix |
|---------|-----|
| `chrome.storage.onChanged` listeners added in popup, never removed | Remove in `window.unload` |
| Content scripts holding DOM references after navigation | Nullify in `beforeunload` |
| Offscreen document kept open indefinitely | Call `chrome.offscreen.closeDocument()` when idle |
| SW module globals holding stale tab objects | Store only `tab.id`; re-query on use |

## Standards and pitfalls

**Avoid misleading measurements:**
- Profile in Incognito to eliminate extension noise.
- Never generalize from a high-end desktop — use CPU throttling.
- Clear Node.js performance timeline state between profiling runs in the same process.

**Regression prevention:**
- Add `performance.mark()` / `measure()` instrumentation around critical paths in production code.
- Monitor Core Web Vitals via RUM (Real User Monitoring), not just lab Lighthouse scores.
- Treat regressions as a measurement problem first — re-instrument before changing code.
