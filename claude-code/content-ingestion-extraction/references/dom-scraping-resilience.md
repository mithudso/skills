<!-- hub-reference-banner -->
> **Reference file — part of the `content-ingestion-extraction` hub.** Formerly the standalone `dom-scraping-resilience` skill.
> Sibling topics in this family are now reference files under the hubs (`content-ingestion-extraction`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: dom-scraping-resilience
description: >
  DOM scraping resilience patterns for Chrome extension content scripts — selector fallback chains,
  MutationObserver, waitForElement, shadow DOM traversal (open and closed via chrome.dom API),
  SPA navigation detection, and graceful partial extraction.
  TRIGGER: user is building or reviewing a content script that extracts DOM data; debugging
  "element not found" or null-extraction errors in a content script; adding MutationObserver or
  waitForElement logic; handling SPA pushState navigation in an extension; traversing shadow DOM;
  implementing partial extraction with confidence scoring.
  SKIP: general web scraping without a Chrome extension context (use programming-languages);
  Chrome extension architecture questions not about DOM extraction (use chrome-extension-expert).
version: "1.2.1"
updated: "2026-05-31"
category: developer
tags: [chrome-extension, content-script, dom, scraping, mutation-observer, selectors, shadow-dom, spa]
keywords:
  - content script DOM extraction
  - selector fallback chain
  - MutationObserver waitForElement
  - shadow DOM traversal
  - SPA navigation detection
  - chrome.dom openOrClosedShadowRoot
  - partial extraction confidence
  - schema drift detection
whenToUse:
  - Building a content script that scrapes or extracts DOM data from a web page
  - Reviewing extractor code for selector brittleness or timing issues
  - Debugging "element not found" or null-extraction failures in a content script
  - Adding MutationObserver or waitForElement logic to handle async DOM rendering
  - Handling SPA pushState/replaceState navigation in a Chrome extension
  - Traversing open or closed shadow DOM roots from a content script
  - Adding schema versioning and drift detection to an extraction pipeline
whenNotToUse:
  - General web scraping outside a Chrome extension (use programming-languages)
  - Chrome extension architecture questions unrelated to DOM extraction (use chrome-extension-expert)
  - Service worker message handling or storage (use chrome-extension-expert)
related_skills:
  - chrome-extension-expert
  - programming-languages
triggers:
  - building a content script that scrapes or extracts DOM data
  - reviewing extractor code for resilience or brittleness
  - debugging "element not found" or null-extraction issues in a content script
  - adding MutationObserver or waitForElement logic
  - handling SPA navigation in an extension
  - traversing shadow DOM from a content script
globs:
  - "**/content/**"
  - "**/content-script*"
  - "**/extractor*"
  - "**/hub-extractor*"
  - "**/scraper*"
---

# DOM Scraping Resilience

Expert reference for resilient DOM extraction in Chrome MV3 content scripts. A response from this skill is correct when it applies the appropriate timing strategy, provides a selector fallback chain, handles SPA navigation, and returns partial results rather than crashing on missing elements.

> **Staleness note:** Chrome extension APIs (chrome.dom, MV3 `run_at` values, `world: "MAIN"` execution) were current as of May 2026. Verify `chrome.dom.openOrClosedShadowRoot` availability for your extension's minimum Chrome version.

**Navigation by task:**
- Content script injection timing (`run_at`, `document_idle` vs `document_start`) → Content Script Injection Timing
- Selector fallback chains and text-content fallbacks → Selector Fallback Chains
- MutationObserver (basic, debounced, attribute) → MutationObserver Patterns
- `waitForElement`, cancellation, retry, parallel waits → waitForElement Implementation
- Open and closed shadow DOM traversal → Shadow DOM Traversal
- SPA navigation detection (history override, webNavigation, polling) → SPA Navigation Detection
- Partial extraction with confidence scoring → Error Boundaries and Partial Extraction
- Schema versioning and drift detection → Schema Versioning and Drift Detection
- Anti-patterns table and quick-reference decision tree → Anti-Patterns / Quick Reference

## Overview

Content scripts that scrape page DOM are inherently fragile. Sites redesign, frameworks
randomize class names, SPAs swap entire subtrees without page reloads, and shadow DOM
hides elements behind closed roots. This skill codifies the patterns that keep extraction
working across all of those failure modes.

The core principle: **never trust a single selector, never assume timing, always degrade
gracefully.** Every extraction should produce partial results rather than crash, and every
selector should have a fallback chain rather than a single brittle path.

---

## Content Script Injection Timing

Chrome MV3 provides three `run_at` values that control when your content script executes.
Choosing the wrong one is the most common source of "element not found" bugs.

| Value            | When it fires                                       | Use case                            |
|------------------|-----------------------------------------------------|-------------------------------------|
| `document_start` | Before any DOM is constructed                       | Intercept network, patch globals    |
| `document_end`   | DOM parsed, sub-resources still loading             | Early DOM read, before images load  |
| `document_idle`  | After `document_end`, browser picks optimal moment  | **Default. Best for most scrapers** |

```jsonc
// manifest.json
"content_scripts": [{
  "matches": ["https://example.com/*"],
  "js": ["src/content/extractor.js"],
  "run_at": "document_idle"          // safe default
}]
```

**Rule of thumb:** use `document_idle` (the default) for extraction. Only drop to
`document_start` when you need to monkey-patch `history.pushState` or intercept
fetch/XHR before the page runs. `document_end` is rarely the right choice -- it fires
before deferred scripts execute, so framework-rendered content is usually missing.

For programmatic injection from the service worker, use `chrome.scripting.executeScript`
or `chrome.scripting.registerContentScripts` with the same `runAt` options:

```js
// service-worker.js -- inject on demand
chrome.scripting.executeScript({
  target: { tabId },
  files: ['src/content/extractor.js'],
  injectImmediately: false   // equivalent to document_idle
});
```

---

## Selector Fallback Chains

A single CSS selector is a single point of failure. Sites change class names on every
deploy (CSS-in-JS hashes), restructure DOM during A/B tests, and differ between
desktop and mobile layouts. A fallback chain tries selectors from most-specific to
most-generic, returning the first match.

### Selector priority hierarchy

1. **`data-testid` / `data-qa` attributes** -- stable across deploys, added for testing
2. **`aria-label` / `role` attributes** -- semantic, unlikely to change
3. **ID selectors** -- unique but sometimes generated
4. **Structural selectors** -- `main > section:first-child h1` -- survive class renames
5. **Class-name selectors** -- last resort, most brittle

### Implementation

```js
/**
 * Try selectors in order, return first match or null.
 * @param {Element|Document} root
 * @param {string[]} selectors
 * @returns {Element|null}
 */
function queryFallback(root, selectors) {
  for (const sel of selectors) {
    try {
      const el = root.querySelector(sel);
      if (el) return el;
    } catch {
      // invalid selector string -- skip silently
    }
  }
  return null;
}

/**
 * queryAll variant -- returns first non-empty NodeList.
 */
function queryAllFallback(root, selectors) {
  for (const sel of selectors) {
    try {
      const els = root.querySelectorAll(sel);
      if (els.length > 0) return Array.from(els);
    } catch { /* skip */ }
  }
  return [];
}
```

### Usage with a selector map

```js
const SELECTORS = {
  caseNumber: [
    '[data-testid="case-number"]',
    '[aria-label="Case Number"]',
    '#case-number',
    '.case-header .case-id',
    'h1.title > span:first-child',
  ],
  severity: [
    '[data-testid="severity-badge"]',
    '[role="status"][aria-label*="severity" i]',
    '.severity-indicator',
    '.badge.severity',
  ],
};

function extractCase(root) {
  return {
    caseNumber: queryFallback(root, SELECTORS.caseNumber)?.textContent?.trim() ?? null,
    severity:   queryFallback(root, SELECTORS.severity)?.textContent?.trim() ?? null,
  };
}
```

### Text-content fallback

When no structural selector works, search by visible text as a last-resort heuristic:

```js
function findByText(root, tagName, textPattern) {
  const els = root.querySelectorAll(tagName);
  for (const el of els) {
    if (textPattern.test(el.textContent)) return el;
  }
  return null;
}

// fallback: find the heading that contains "Case #"
const heading = findByText(document, 'h1,h2,h3', /Case\s*#?\s*\d+/i);
```

---

## MutationObserver Patterns

Modern SPAs render content asynchronously. A `MutationObserver` watches for DOM changes
and fires your callback when the target element appears.

### Basic: watch for a container to appear

```js
function observeForElement(selector, callback, opts = {}) {
  const {
    root = document.body,
    timeout = 15_000,
    once = true,
  } = opts;

  // Already present?
  const existing = root.querySelector(selector);
  if (existing) {
    callback(existing);
    if (once) return { disconnect() {} };
    // fall through to start observer for future matches
  }

  let timer;
  const observer = new MutationObserver(() => {
    const el = root.querySelector(selector);
    if (el) {
      callback(el);
      if (once) {
        observer.disconnect();
        clearTimeout(timer);
      }
    }
  });

  observer.observe(root, { childList: true, subtree: true });

  if (timeout > 0) {
    timer = setTimeout(() => {
      observer.disconnect();
      callback(null); // signal timeout
    }, timeout);
  }

  return observer;
}
```

### Advanced: debounced batch observer

SPAs like React can trigger dozens of mutations per render cycle. Debouncing prevents
your extractor from running on every intermediate state:

```js
function observeDebounced(selector, callback, {
  root = document.body,
  debounceMs = 300,
  timeout = 15_000,
} = {}) {
  let debounceTimer;
  let timeoutTimer;

  const observer = new MutationObserver(() => {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      const el = root.querySelector(selector);
      if (el) {
        observer.disconnect();
        clearTimeout(timeoutTimer);
        callback(el);
      }
    }, debounceMs);
  });

  observer.observe(root, { childList: true, subtree: true });

  timeoutTimer = setTimeout(() => {
    observer.disconnect();
    clearTimeout(debounceTimer);
    callback(null);
  }, timeout);

  return observer;
}
```

### Attribute change observer

Some frameworks toggle visibility via attributes rather than adding/removing nodes:

```js
function observeAttribute(element, attrName, callback) {
  const observer = new MutationObserver((mutations) => {
    for (const m of mutations) {
      if (m.type === 'attributes' && m.attributeName === attrName) {
        callback(element.getAttribute(attrName), m.oldValue);
      }
    }
  });
  observer.observe(element, {
    attributes: true,
    attributeFilter: [attrName],
    attributeOldValue: true,
  });
  return observer;
}
```

---

## waitForElement Implementation

A Promise-based utility that combines MutationObserver with timeout. This is the
single most reused primitive in resilient content scripts.

```js
/**
 * Wait for an element matching `selector` to appear in the DOM.
 * @param {string}  selector   CSS selector
 * @param {object}  opts
 * @param {Element} opts.root  Observation root (default: document.body)
 * @param {number}  opts.timeout  Max wait in ms (default: 10000)
 * @returns {Promise<Element>}  Resolves with the element, rejects on timeout
 */
function waitForElement(selector, { root = document.body, timeout = 10_000 } = {}) {
  return new Promise((resolve, reject) => {
    // Fast path: already in DOM
    const existing = root.querySelector(selector);
    if (existing) return resolve(existing);

    let timer;
    const observer = new MutationObserver(() => {
      const el = root.querySelector(selector);
      if (el) {
        observer.disconnect();
        clearTimeout(timer);
        resolve(el);
      }
    });

    observer.observe(root, { childList: true, subtree: true });

    timer = setTimeout(() => {
      observer.disconnect();
      reject(new Error(`waitForElement("${selector}") timed out after ${timeout}ms`));
    }, timeout);
  });
}
```

### Cancellation with AbortController

When SPA navigation occurs, you should cancel any pending `waitForElement` calls from
the previous page. Wrap the observer in an AbortController-compatible pattern:

```js
function waitForElementAbortable(selector, { root = document.body, timeout = 10_000, signal } = {}) {
  return new Promise((resolve, reject) => {
    if (signal?.aborted) return reject(new DOMException('Aborted', 'AbortError'));

    const existing = root.querySelector(selector);
    if (existing) return resolve(existing);

    let timer;
    const observer = new MutationObserver(() => {
      const el = root.querySelector(selector);
      if (el) { cleanup(); resolve(el); }
    });

    function cleanup() {
      observer.disconnect();
      clearTimeout(timer);
      signal?.removeEventListener('abort', onAbort);
    }

    function onAbort() {
      cleanup();
      reject(new DOMException('Aborted', 'AbortError'));
    }

    signal?.addEventListener('abort', onAbort);
    observer.observe(root, { childList: true, subtree: true });
    timer = setTimeout(() => {
      cleanup();
      reject(new Error(`waitForElement("${selector}") timed out after ${timeout}ms`));
    }, timeout);
  });
}

// Usage: cancel all pending waits on SPA navigation
let abortCtrl = new AbortController();
onNavigation(() => {
  abortCtrl.abort();                  // cancel stale waits
  abortCtrl = new AbortController();  // fresh controller for new page
  scheduleReExtraction(abortCtrl.signal);
});
```

### waitForElement with retry and exponential backoff

When the target element appears only after an async data fetch (not a simple DOM mount),
combine `waitForElement` with retries:

```js
async function waitForElementRetry(selector, {
  root = document.body,
  maxRetries = 4,
  baseDelay = 500,
  timeout = 5_000,
} = {}) {
  for (let attempt = 0; attempt <= maxRetries; attempt++) {
    try {
      return await waitForElement(selector, { root, timeout });
    } catch {
      if (attempt === maxRetries) {
        throw new Error(
          `waitForElementRetry("${selector}") failed after ${maxRetries + 1} attempts`
        );
      }
      const delay = baseDelay * Math.pow(2, attempt) + Math.random() * 100; // jitter
      await new Promise(r => setTimeout(r, delay));
    }
  }
}
```

### waitForElements: wait for multiple selectors in parallel

```js
async function waitForElements(selectorMap, opts = {}) {
  const entries = Object.entries(selectorMap);
  const results = await Promise.allSettled(
    entries.map(([key, sel]) =>
      waitForElement(sel, opts).then(el => [key, el])
    )
  );

  const found = {};
  const missing = [];
  for (let i = 0; i < results.length; i++) {
    const [key] = entries[i];
    if (results[i].status === 'fulfilled') {
      found[key] = results[i].value[1];
    } else {
      missing.push(key);
    }
  }
  return { found, missing };
}

// Usage:
const { found, missing } = await waitForElements({
  caseHeader:  '[data-testid="case-header"]',
  commentList: '.comment-thread',
  sidebar:     '#case-sidebar',
}, { timeout: 8_000 });

if (missing.length > 0) {
  console.warn('Partial extraction -- missing:', missing);
}
```

---

## Shadow DOM Traversal

Web components hide their internals behind shadow roots. Content scripts need special
handling to reach inside them.

### Open shadow roots

Open shadow roots are accessible via `element.shadowRoot`:

```js
function queryShadow(root, selector) {
  // Try the light DOM first
  const light = root.querySelector(selector);
  if (light) return light;

  // Walk into open shadow roots
  const walk = (node) => {
    if (node.shadowRoot) {
      const found = node.shadowRoot.querySelector(selector);
      if (found) return found;
      // Recurse into shadow root children
      for (const child of node.shadowRoot.querySelectorAll('*')) {
        const deep = walk(child);
        if (deep) return deep;
      }
    }
    for (const child of node.querySelectorAll('*')) {
      if (child.shadowRoot) {
        const deep = walk(child);
        if (deep) return deep;
      }
    }
    return null;
  };

  return walk(root);
}
```

### Closed shadow roots with chrome.dom

Chrome extensions have privileged access to closed shadow roots through the `chrome.dom`
API. This requires the `"dom"` permission in your manifest (available in MV3).

```jsonc
// manifest.json
"permissions": ["dom"]
```

```js
// Content script -- access closed shadow root
function getClosedShadowRoot(element) {
  if (typeof chrome !== 'undefined' && chrome.dom?.openOrClosedShadowRoot) {
    return chrome.dom.openOrClosedShadowRoot(element);
  }
  // Fallback to open shadow root
  return element.shadowRoot;
}

function queryDeepShadow(root, selector) {
  const light = root.querySelector(selector);
  if (light) return light;

  const walk = (node) => {
    const shadow = getClosedShadowRoot(node);
    if (shadow) {
      const found = shadow.querySelector(selector);
      if (found) return found;
      for (const child of shadow.querySelectorAll('*')) {
        const deep = walk(child);
        if (deep) return deep;
      }
    }
    for (const child of node.querySelectorAll('*')) {
      const deep = walk(child);
      if (deep) return deep;
    }
    return null;
  };

  return walk(root);
}
```

### Performance guard (recommended approach)

Deep shadow DOM traversal can be expensive. **Prefer this bounded version** over the
unbounded `queryShadow`/`queryDeepShadow` variants above. Limit the depth to avoid
runaway recursion:

```js
function queryShadowBounded(root, selector, maxDepth = 5) {
  if (maxDepth <= 0) return null;
  const light = root.querySelector(selector);
  if (light) return light;

  const shadows = root.querySelectorAll('*');
  for (const el of shadows) {
    const sr = el.shadowRoot || getClosedShadowRoot?.(el);
    if (sr) {
      const found = queryShadowBounded(sr, selector, maxDepth - 1);
      if (found) return found;
    }
  }
  return null;
}
```

---

## SPA Navigation Detection

SPAs update the URL via `history.pushState` / `replaceState` without triggering a page
load. Content scripts injected at `document_idle` will NOT re-run on SPA navigations.
You must detect them yourself.

### Method 1: Override history methods (injected into page context)

This must run in the **page world**, not the isolated content script world, because
`history.pushState` lives on the page's `History` prototype:

```js
// inject-navigation-hook.js -- runs in MAIN world
(function () {
  const originalPush = history.pushState.bind(history);
  const originalReplace = history.replaceState.bind(history);

  function dispatch(type) {
    window.dispatchEvent(new CustomEvent('mca:navigation', {
      detail: { type, url: location.href },
    }));
  }

  history.pushState = function (...args) {
    originalPush(...args);
    dispatch('pushState');
  };

  history.replaceState = function (...args) {
    originalReplace(...args);
    dispatch('replaceState');
  };

  window.addEventListener('popstate', () => dispatch('popstate'));
})();
```

```jsonc
// manifest.json
"content_scripts": [{
  "matches": ["https://example.com/*"],
  "js": ["src/content/inject-navigation-hook.js"],
  "run_at": "document_start",
  "world": "MAIN"
}]
```

Your isolated content script listens on the custom event:

```js
// content-script.js (ISOLATED world)
window.addEventListener('mca:navigation', (e) => {
  console.log('SPA navigated:', e.detail.type, e.detail.url);
  scheduleReExtraction();
});
```

### Method 2: chrome.webNavigation API (service worker)

The service worker can listen for SPA-style navigations and message the content script:

```js
// service-worker.js
chrome.webNavigation.onHistoryStateUpdated.addListener((details) => {
  chrome.tabs.sendMessage(details.tabId, {
    type: 'MCA_SPA_NAVIGATION',
    url: details.url,
    frameId: details.frameId,
  });
}, { url: [{ hostContains: 'example.com' }] });
```

### Method 3: URL polling fallback

For edge cases where neither history patching nor webNavigation fires (hash-based
routing, custom navigation libraries):

```js
function pollForUrlChange(callback, intervalMs = 500) {
  let lastUrl = location.href;
  const timer = setInterval(() => {
    if (location.href !== lastUrl) {
      const prev = lastUrl;
      lastUrl = location.href;
      callback({ prev, current: lastUrl });
    }
  }, intervalMs);
  return () => clearInterval(timer);
}
```

### Combining all three for maximum resilience

```js
function onNavigation(callback) {
  // Method 1: custom event from MAIN-world hook
  window.addEventListener('mca:navigation', (e) => callback(e.detail));

  // Method 2: message from service worker
  chrome.runtime.onMessage.addListener((msg) => {
    if (msg.type === 'MCA_SPA_NAVIGATION') callback(msg);
  });

  // Method 3: polling fallback
  pollForUrlChange(callback);
}

onNavigation(({ url }) => {
  console.log('Navigation detected, re-extracting...');
  scheduleReExtraction();
});
```

---

## Error Boundaries and Partial Extraction

The cardinal rule: **never let one failed selector crash the entire extraction.**
Return whatever you can, annotate what is missing, and let the consumer decide whether
partial data is acceptable.

### Extraction envelope pattern

```js
/**
 * @typedef {Object} ExtractionResult
 * @property {string}   schemaVersion  Extraction schema version
 * @property {number}   timestamp      When extraction ran
 * @property {Object}   data           Extracted fields (nulls where missing)
 * @property {string[]} missing        Field names that failed to extract
 * @property {Object[]} errors         Per-field error details
 * @property {number}   confidence     0-1 score based on fields found
 */
function extractWithEnvelope(root, fieldExtractors) {
  const data = {};
  const missing = [];
  const errors = [];
  let found = 0;
  const total = Object.keys(fieldExtractors).length;

  for (const [field, extractFn] of Object.entries(fieldExtractors)) {
    try {
      const value = extractFn(root);
      if (value != null && value !== '') {
        data[field] = value;
        found++;
      } else {
        data[field] = null;
        missing.push(field);
      }
    } catch (err) {
      data[field] = null;
      missing.push(field);
      errors.push({ field, message: err.message, stack: err.stack });
    }
  }

  return {
    schemaVersion: '1.0.0',
    timestamp: Date.now(),
    data,
    missing,
    errors,
    confidence: found / total,
  };
}
```

### Usage

```js
const result = extractWithEnvelope(document, {
  caseNumber: (root) => queryFallback(root, SELECTORS.caseNumber)?.textContent?.trim(),
  severity:   (root) => queryFallback(root, SELECTORS.severity)?.textContent?.trim(),
  assignee:   (root) => queryFallback(root, SELECTORS.assignee)?.textContent?.trim(),
  createdAt:  (root) => queryFallback(root, SELECTORS.createdAt)?.getAttribute('datetime'),
});

if (result.confidence < 0.5) {
  console.warn('Low-confidence extraction:', result.missing);
  // optionally skip sending to background, or flag for manual review
}

// Send partial result anyway -- consumer decides
chrome.runtime.sendMessage({ type: 'MCA_UPSERT_CASE_CONTEXT', payload: result });
```

### Field-level retry

For high-value fields, retry with a short delay before accepting null:

```js
async function extractWithRetry(root, fieldExtractors, {
  retryFields = [],
  retryDelay = 1_000,
  maxRetries = 2,
} = {}) {
  let result = extractWithEnvelope(root, fieldExtractors);

  for (let attempt = 0; attempt < maxRetries; attempt++) {
    const stillMissing = result.missing.filter(f => retryFields.includes(f));
    if (stillMissing.length === 0) break;

    await new Promise(r => setTimeout(r, retryDelay * (attempt + 1)));

    for (const field of stillMissing) {
      try {
        const value = fieldExtractors[field](root);
        if (value != null && value !== '') {
          result.data[field] = value;
          result.missing = result.missing.filter(f => f !== field);
          result.confidence = (Object.keys(fieldExtractors).length - result.missing.length)
            / Object.keys(fieldExtractors).length;
        }
      } catch { /* still missing */ }
    }
  }

  return result;
}
```

---

## Schema Versioning and Drift Detection

When the target site redesigns, selectors silently return null. Detect this early
rather than shipping garbage data downstream.

### Schema version tagging

Tag every extraction with the selector schema version. Bump the version when you
update any selector in the map:

```js
const SCHEMA_VERSION = '2.3.0';  // bump on any selector change
const SCHEMA_DATE = '2026-05-26';

function tagResult(result) {
  return {
    ...result,
    schemaVersion: SCHEMA_VERSION,
    schemaDate: SCHEMA_DATE,
  };
}
```

### Drift detection: field-count anomaly

Compare current extraction results against historical norms:

```js
function detectDrift(result, historicalAvg) {
  const {
    avgConfidence = 0.9,
    minFields = 4,
  } = historicalAvg;

  const foundCount = Object.keys(result.data).length - result.missing.length;
  const warnings = [];

  if (result.confidence < avgConfidence * 0.7) {
    warnings.push(
      `Confidence ${result.confidence.toFixed(2)} is well below historical avg ${avgConfidence}`
    );
  }
  if (foundCount < minFields) {
    warnings.push(`Only ${foundCount} fields extracted (expected >= ${minFields})`);
  }
  if (result.errors.length > 0) {
    warnings.push(`${result.errors.length} extraction error(s)`);
  }

  return {
    drifted: warnings.length > 0,
    warnings,
  };
}
```

### Selector health check

Periodically validate that your selectors still match something on the page,
independent of data extraction:

```js
function selectorHealthCheck(selectorMap) {
  const report = {};
  for (const [field, selectors] of Object.entries(selectorMap)) {
    const hits = selectors.map(sel => ({
      selector: sel,
      matches: document.querySelectorAll(sel).length,
    }));
    const anyHit = hits.some(h => h.matches > 0);
    report[field] = { hits, healthy: anyHit };
  }
  return report;
}

// Log health on each extraction cycle
const health = selectorHealthCheck(SELECTORS);
const broken = Object.entries(health)
  .filter(([, v]) => !v.healthy)
  .map(([k]) => k);

if (broken.length > 0) {
  console.error('Broken selector groups:', broken);
  // Optionally send telemetry to background for alerting
}
```

---

## Anti-Patterns

| Anti-pattern | Why it fails | What to do instead |
|---|---|---|
| Single hardcoded selector | One deploy breaks it | Fallback chain of 3-5 selectors |
| `setTimeout(fn, 2000)` for timing | Arbitrary; too slow or too fast | `waitForElement` with MutationObserver |
| Observing `document` instead of a subtree | Performance: fires on every DOM change | Scope observer to closest stable ancestor |
| No `observer.disconnect()` | Memory leak, runs forever | Always disconnect after match or timeout |
| `.innerHTML` parsing for data | XSS risk, fragile | `.textContent` + attribute reads |
| Ignoring closed shadow roots | Misses web-component content | `chrome.dom.openOrClosedShadowRoot()` |
| Re-running full extraction on every mutation | Perf disaster on large pages | Debounce mutations (200-500ms) |
| Throwing on missing optional fields | Kills entire extraction | Return null, log to `missing[]` |
| Class-name selectors on CSS-in-JS sites | Hashes change every build | Use `data-testid`, ARIA, structural selectors |
| No schema version on extracted data | Cannot detect silent degradation | Tag with version + date, alert on drift |

---

## Troubleshooting

### Element exists in DevTools but not in content script
- Check `run_at` timing. The element may not exist when your script runs.
- The element may be inside a shadow root (open or closed).
- The element may be in an iframe; content scripts run per-frame, not across frames.

### Extraction works on first load but not on SPA navigation
- Content scripts do not re-inject on `pushState`/`replaceState`.
- Add the SPA navigation detection pattern (history override + webNavigation listener).
- After detecting navigation, add a short delay or `waitForElement` before re-extracting.

### MutationObserver callback never fires
- Confirm the `observe()` target is already in the DOM when you call `observe()`.
- Check that `childList: true` and `subtree: true` are set.
- The change may be an attribute change (add `attributes: true`).

### Extraction returns null for everything after a site update
- Run `selectorHealthCheck()` to identify which selector groups are broken.
- Check `detectDrift()` warnings.
- Update the selector map and bump `SCHEMA_VERSION`.

### Performance degradation from observers
- Scope the observer to the smallest possible subtree.
- Debounce the callback (200-500ms).
- Disconnect the observer as soon as the target element is found.
- Avoid `querySelectorAll('*')` inside observer callbacks on large pages.

### Shadow DOM traversal returns null
- Verify `"permissions": ["dom"]` is in manifest.json for closed root access.
- Confirm `chrome.dom.openOrClosedShadowRoot` is available (MV3 only).
- Limit traversal depth to avoid exponential blowup on deeply nested components.

---

## Quick Reference: Choosing the Right Pattern

```
Need to find an element?
  +-- Already in DOM?  -->  queryFallback() with selector chain
  +-- Appears later?   -->  waitForElement() with timeout
  +-- Appears after SPA nav?  -->  onNavigation() + waitForElement()
  +-- Inside shadow DOM?  -->  queryShadow() / queryDeepShadow()

Need to watch for changes?
  +-- Element added/removed  -->  MutationObserver (childList + subtree)
  +-- Attribute changed      -->  MutationObserver (attributes + filter)
  +-- URL changed (SPA)      -->  history override + webNavigation + polling

How to handle failures?
  +-- Single field missing   -->  Return null, add to missing[]
  +-- Most fields missing    -->  Flag low confidence, warn, still return partial
  +-- All fields missing     -->  Schema drift alert, do NOT send garbage downstream
```

---

## References

- [Chrome Content Scripts docs](https://developer.chrome.com/docs/extensions/develop/concepts/content-scripts)
- [MutationObserver MDN](https://developer.mozilla.org/en-US/docs/Web/API/MutationObserver)
- [chrome.dom API](https://developer.chrome.com/docs/extensions/reference/api/dom)
- [Navigation API for SPAs](https://developer.chrome.com/docs/web-platform/navigation-api)
- [chrome.webNavigation API](https://developer.chrome.com/docs/extensions/reference/api/webNavigation)
- [Shadow DOM mode (open vs closed)](https://developer.mozilla.org/en-US/docs/Web/API/ShadowRoot/mode)
- [Manifest content_scripts reference](https://developer.chrome.com/docs/extensions/reference/manifest/content-scripts)
