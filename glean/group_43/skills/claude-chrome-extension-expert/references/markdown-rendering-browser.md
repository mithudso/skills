<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly the standalone `markdown-rendering-browser` skill.
> Sibling topics in this family are now reference files under the hubs (`chrome-extension-expert`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: markdown-rendering-browser
title: "Markdown Rendering in the Browser"
description: >
  Browser-side markdown rendering with marked.js, DOMPurify XSS sanitization,
  syntax highlighting (highlight.js + marked-highlight), CSP compliance for Chrome MV3
  extensions, and shadow DOM isolation patterns.
  TRIGGER: user needs to render markdown to HTML in a browser page or Chrome extension;
  sanitize HTML output from a markdown parser to prevent XSS; add syntax highlighting to
  fenced code blocks; comply with Chrome extension MV3 CSP (script-src self) when loading
  marked.js or DOMPurify; render markdown inside a shadow DOM to prevent CSS bleed;
  configure DOMPurify hooks for safe links or HTTPS-only images.
  SKIP: server-side markdown rendering (use javascript-nodejs or python-patterns);
  Markdown parsing in Node.js without browser/DOM context; MDX or React-specific
  markdown rendering (use react-nextjs).
category: developer
version: "1.1.0"
updated: "2026-05-29"
tags: [markdown, browser, chrome-extension, security, xss, csp, shadow-dom]
keywords:
  - marked.js browser rendering
  - DOMPurify XSS sanitization
  - highlight.js syntax highlighting
  - marked-highlight extension
  - Chrome MV3 CSP script-src self
  - shadow DOM markdown isolation
  - ALLOWED_TAGS ALLOWED_ATTR DOMPurify
  - Trusted Types RETURN_TRUSTED_TYPE
  - marked.parse marked v18
  - bundle marked DOMPurify locally
whenToUse:
  - Rendering markdown to safe HTML in a browser page or Chrome extension panel
  - Sanitizing parsed HTML with DOMPurify before inserting into the DOM
  - Adding syntax highlighting to fenced code blocks with highlight.js
  - Bundling marked.js and DOMPurify locally to comply with Chrome MV3 CSP
  - Rendering markdown inside a shadow DOM to prevent CSS bleed with the host page
  - Configuring DOMPurify hooks to enforce HTTPS links or images
  - Handling Trusted Types enforcement in strict CSP environments
whenNotToUse:
  - Server-side or Node.js markdown rendering without a DOM (use javascript-nodejs)
  - React/MDX component-based markdown (use react-nextjs)
  - Static site generation with markdown (use python-patterns or javascript-nodejs)
related_skills:
  - chrome-dev
  - chrome-mv3-advanced
  - mv3-service-worker-expert
  - dom-scraping-resilience
  - javascript-nodejs
triggers:
  - markdown
  - marked.js
  - DOMPurify
  - sanitize HTML
  - XSS prevention
  - markdown shadow DOM
  - syntax highlighting markdown
  - CSP markdown
  - render markdown browser
---

# Markdown Rendering in the Browser

Expert reference for rendering markdown to safe, styled HTML in browser contexts — vanilla pages, Chrome extensions, shadow DOM overlays, and CSP-restricted environments. A response from this skill is correct when it always sanitizes parsed HTML with DOMPurify before DOM insertion, bundles libraries locally for Chrome extensions, and injects highlight.js styles into the shadow root rather than the document head.

> **Staleness note:** Library versions (marked 18.x, DOMPurify 3.4.x, highlight.js 11.11.x) were current as of May 2026. Verify current versions from npm or jsDelivr before use.

**Navigation by task:**
- Basic marked.js setup and options → §1 Core Libraries → marked.js
- DOMPurify setup and ALLOWED_TAGS/ALLOWED_ATTR configuration → §1 Core Libraries → DOMPurify
- Safe rendering pipeline (marked → DOMPurify → DOM) → §2 The Safe Rendering Pipeline
- XSS prevention checklist, link enforcement, image HTTPS hooks → §3 XSS Prevention Checklist
- Chrome MV3 CSP compliance, local bundling, content script caveats → §4 CSP Compliance for Chrome Extensions
- Syntax highlighting with highlight.js and marked-highlight → §5 Syntax Highlighting
- Rendering markdown inside a shadow root, highlight.js in shadow DOM → §6 Shadow DOM Rendering
- Full Chrome extension panel integration example → §7 Full Integration Example
- Performance (lazy loading, debouncing, pre-instantiation) → §8 Performance Considerations
- Common pitfalls table → §9 Common Pitfalls
- Version reference table → §10 Version Reference

## 1. Core Libraries

### marked.js (v18+)

Lightweight markdown parser and compiler. Zero dependencies, works in
browsers and Node.js.

**CDN (UMD -- works everywhere, assigns `window.marked`):**

```html
<script src="https://cdn.jsdelivr.net/npm/marked/lib/marked.umd.js"></script>
```

**CDN (ESM -- for `<script type="module">` or bundlers):**

```html
<script type="module">
  import { marked } from 'https://cdn.jsdelivr.net/npm/marked/lib/marked.esm.js';
</script>
```

**Basic usage:**

```js
const html = marked.parse('# Hello\n\nParagraph with **bold**.');
```

**Key options:**

```js
marked.setOptions({
  breaks: true,       // GFM line breaks
  gfm: true,          // GitHub Flavored Markdown (default true)
  headerIds: false,    // disable auto-generated IDs (reduces DOM surface)
  mangle: false,       // do not mangle email autolinks
  pedantic: false,     // do not conform to original markdown.pl
});
```

### DOMPurify (v3.4+)

DOM-only XSS sanitizer. Parses HTML through the browser's own DOM parser,
then walks the tree and strips anything dangerous. Orders of magnitude more
reliable than regex-based sanitizers.

**CDN:**

```html
<script src="https://cdn.jsdelivr.net/npm/dompurify/dist/purify.min.js"></script>
```

**ESM:**

```js
import DOMPurify from 'https://cdn.jsdelivr.net/npm/dompurify/dist/purify.es.mjs';
```

## 2. The Safe Rendering Pipeline

**Rule: never inject parsed markdown into the DOM without sanitization.**

```js
// UNSAFE -- XSS if markdown contains attacker-controlled content
container.textContent = ''; // clear first
const parsed = marked.parse(untrustedMarkdown);
// Directly inserting `parsed` as HTML is dangerous without sanitization

// SAFE -- DOMPurify strips dangerous nodes/attributes before DOM insertion
const safeHtml = DOMPurify.sanitize(marked.parse(untrustedMarkdown));
// Now safe to inject `safeHtml` into the DOM
```

### Complete safe render function

```js
/**
 * Render markdown to sanitized HTML and inject into a container.
 * Uses DOMPurify.sanitize() to prevent XSS before any DOM insertion.
 * @param {string} md        - raw markdown string
 * @param {HTMLElement} el   - target container
 * @param {Object} [opts]    - optional DOMPurify config overrides
 */
function renderMarkdown(md, el, opts = {}) {
  const dirty = marked.parse(md);
  const clean = DOMPurify.sanitize(dirty, {
    ALLOWED_TAGS: [
      'h1','h2','h3','h4','h5','h6','p','br','hr',
      'ul','ol','li','blockquote','pre','code',
      'a','strong','em','del','img','table','thead',
      'tbody','tr','th','td','span','div','details','summary',
    ],
    ALLOWED_ATTR: ['href','src','alt','title','class','id','open'],
    ALLOW_DATA_ATTR: false,
    ...opts,
  });
  // Safe: content has been sanitized by DOMPurify above
  el.replaceChildren();
  el.insertAdjacentHTML('afterbegin', clean);
}
```

### What DOMPurify strips by default

- `<script>`, `<iframe>`, `<object>`, `<embed>`, `<form>` tags
- `on*` event handler attributes (`onerror`, `onclick`, etc.)
- `javascript:` and `data:` URI schemes in `href`/`src`
- SVG/MathML-based mXSS vectors
- DOM clobbering payloads (`name="location"`, etc.)

## 3. XSS Prevention Checklist

1. **Always sanitize after parsing.** `marked.parse()` faithfully converts
   markdown (including embedded HTML) to HTML. DOMPurify removes the
   dangerous parts.
2. **Restrict `ALLOWED_TAGS` and `ALLOWED_ATTR`.** A tight allowlist is
   stronger than relying on the default blocklist alone.
3. **Disable `ALLOW_DATA_ATTR`.** Data attributes can be used for DOM
   clobbering and exfiltration.
4. **Use `RETURN_TRUSTED_TYPE: true` when Trusted Types are enforced.**

   ```js
   const clean = DOMPurify.sanitize(dirty, { RETURN_TRUSTED_TYPE: true });
   ```

5. **Never pass sanitized HTML back through dynamic code evaluation,
   string-based timers, or template literal injection.** These re-introduce
   code execution paths that sanitization was designed to prevent.
6. **Validate links.** Use a DOMPurify hook to enforce `https:` only:

   ```js
   DOMPurify.addHook('afterSanitizeAttributes', (node) => {
     if (node.tagName === 'A') {
       const href = node.getAttribute('href') || '';
       if (!/^https?:\/\//i.test(href) && !href.startsWith('#')) {
         node.removeAttribute('href');
       }
       node.setAttribute('target', '_blank');
       node.setAttribute('rel', 'noopener noreferrer');
     }
   });
   ```

7. **Lock images to HTTPS sources:**

   ```js
   DOMPurify.addHook('afterSanitizeAttributes', (node) => {
     if (node.tagName === 'IMG') {
       const src = node.getAttribute('src') || '';
       if (!src.startsWith('https://')) {
         node.removeAttribute('src');
       }
     }
   });
   ```

## 4. CSP Compliance for Chrome Extensions (MV3)

Chrome MV3 extensions enforce a strict default CSP that bans dynamic code
evaluation, inline scripts, and remote script sources on extension-owned
pages (popup, options, dashboard).

### Extension page CSP (manifest.json)

```jsonc
{
  "content_security_policy": {
    "extension_pages": "script-src 'self'; object-src 'none';"
  }
}
```

**What this means for markdown rendering:**

| Approach | Works on extension pages? | Notes |
|---|---|---|
| `<script src="cdn.jsdelivr.net/...">` | No | Remote scripts blocked by `script-src 'self'` |
| Bundle marked.js + DOMPurify locally | Yes | Copy the UMD files into the extension and reference them with relative paths |
| `insertAdjacentHTML` with sanitized content | Yes | Injecting DOMPurify-sanitized content is permitted |
| Dynamic code evaluation (`unsafe-eval`) | No | Blocked by CSP; marked.js does not require it |
| Trusted Types enforcement | Yes | DOMPurify supports `RETURN_TRUSTED_TYPE` natively |

### How to bundle for a Chrome extension

```
extension-root/
  lib/
    marked.umd.js          # copy from node_modules/marked/lib/
    purify.min.js           # copy from node_modules/dompurify/dist/
    highlight.min.js        # copy from node_modules/highlight.js/
    highlight-github.css    # copy a theme CSS
  src/
    panel/
      panel.html
      panel.js
```

```html
<!-- panel.html -->
<script src="../../lib/marked.umd.js"></script>
<script src="../../lib/purify.min.js"></script>
<script src="../../lib/highlight.min.js"></script>
<link rel="stylesheet" href="../../lib/highlight-github.css">
<script src="panel.js"></script>
```

### Content scripts and host page CSP

Content scripts run in an isolated world and are NOT subject to the
extension's CSP -- they inherit the host page's CSP. However, injecting
`<script>` tags into the host page IS subject to the host CSP. Best
practice: do all markdown rendering inside the extension context (e.g.,
an iframe loaded from an extension URL) rather than in the host page DOM.

## 5. Syntax Highlighting with highlight.js

### marked-highlight extension

The official `marked-highlight` package bridges marked.js and highlight.js.

**CDN (UMD):**

```html
<script src="https://cdn.jsdelivr.net/npm/marked/lib/marked.umd.js"></script>
<script src="https://cdn.jsdelivr.net/npm/marked-highlight/lib/index.umd.js"></script>
<script src="https://cdn.jsdelivr.net/npm/highlight.js/lib/common.js"></script>
<link rel="stylesheet"
      href="https://cdn.jsdelivr.net/npm/highlight.js/styles/github.min.css">
```

**Integration:**

```js
const { Marked } = globalThis.marked || marked;
const { markedHighlight } = globalThis.markedHighlight;

const md = new Marked(
  markedHighlight({
    emptyLangClass: 'hljs',
    langPrefix: 'hljs language-',
    highlight(code, lang) {
      const language = hljs.getLanguage(lang) ? lang : 'plaintext';
      return hljs.highlight(code, { language }).value;
    },
  })
);

const html = md.parse('```js\nconsole.log("hello");\n```');
```

### Minimal manual approach (no marked-highlight)

If you do not want the extra dependency, override the `renderer.code`
method directly:

```js
const renderer = new marked.Renderer();
renderer.code = function (code, lang) {
  const language = hljs.getLanguage(lang) ? lang : 'plaintext';
  const highlighted = hljs.highlight(code, { language }).value;
  return `<pre><code class="hljs language-${language}">${highlighted}</code></pre>`;
};

marked.setOptions({ renderer });
```

### Theme selection

highlight.js ships 40+ themes. Pick one and load its CSS:

```html
<!-- Light -->
<link rel="stylesheet" href="highlight.js/styles/github.min.css">

<!-- Dark -->
<link rel="stylesheet" href="highlight.js/styles/github-dark.min.css">

<!-- Adaptive (prefers-color-scheme) -->
<link rel="stylesheet" href="highlight.js/styles/github.min.css"
      media="(prefers-color-scheme: light)">
<link rel="stylesheet" href="highlight.js/styles/github-dark.min.css"
      media="(prefers-color-scheme: dark)">
```

## 6. Shadow DOM Rendering

Shadow DOM isolates rendered markdown from the host page's styles and
prevents CSS bleed in both directions. This is critical for Chrome
extension overlays and web components.

### Rendering markdown inside a shadow root

```js
function createMarkdownShadow(hostEl, markdownText) {
  const shadow = hostEl.attachShadow({ mode: 'open' });

  // Inject scoped styles
  const style = document.createElement('style');
  style.textContent = `
    :host { all: initial; display: block; font-family: system-ui, sans-serif; }
    h1, h2, h3 { margin: 0.5em 0 0.25em; }
    p { margin: 0.4em 0; line-height: 1.5; }
    pre { background: #f6f8fa; padding: 0.75em; border-radius: 6px; overflow-x: auto; }
    code { font-family: 'SFMono-Regular', Consolas, monospace; font-size: 0.9em; }
    a { color: #0969da; text-decoration: none; }
    a:hover { text-decoration: underline; }
    blockquote {
      border-left: 3px solid #d0d7de; margin: 0.5em 0;
      padding-left: 1em; color: #57606a;
    }
    table { border-collapse: collapse; width: 100%; }
    th, td { border: 1px solid #d0d7de; padding: 6px 12px; text-align: left; }
    img { max-width: 100%; height: auto; }
  `;
  shadow.appendChild(style);

  // Render sanitized markdown into the shadow root
  const wrapper = document.createElement('div');
  const safeHtml = DOMPurify.sanitize(marked.parse(markdownText));
  wrapper.insertAdjacentHTML('afterbegin', safeHtml);
  shadow.appendChild(wrapper);

  return shadow;
}
```

### DOMPurify inside shadow DOM

DOMPurify v3+ automatically handles shadow DOM contexts. When you call
`DOMPurify.sanitize()` it uses the correct document context. No special
configuration is required -- just call it normally before inserting content
inside the shadow root.

### highlight.js inside shadow DOM

highlight.js styles must be injected INTO the shadow root because external
stylesheets do not cross the shadow boundary:

```js
function injectHighlightStyles(shadowRoot) {
  const link = document.createElement('link');
  link.rel = 'stylesheet';
  link.href = chrome.runtime.getURL('lib/highlight-github.css');
  shadowRoot.appendChild(link);
}
```

For Chrome extensions, use `chrome.runtime.getURL()` to resolve the path
to the bundled CSS file from inside the shadow DOM.

## 7. Full Integration Example (Chrome Extension Panel)

This combines all patterns: local bundling, safe rendering, syntax
highlighting, and shadow DOM isolation.

```js
// panel.js -- loaded inside an extension iframe (panel.html)

/** One-time setup */
const md = new marked.Marked(
  markedHighlight({
    emptyLangClass: 'hljs',
    langPrefix: 'hljs language-',
    highlight(code, lang) {
      const language = hljs.getLanguage(lang) ? lang : 'plaintext';
      return hljs.highlight(code, { language }).value;
    },
  })
);

// Enforce safe links
DOMPurify.addHook('afterSanitizeAttributes', (node) => {
  if (node.tagName === 'A') {
    node.setAttribute('target', '_blank');
    node.setAttribute('rel', 'noopener noreferrer');
  }
});

/**
 * Render a markdown string into the target element.
 * All content is sanitized by DOMPurify before DOM insertion.
 * @param {string} raw       - markdown text (may be untrusted)
 * @param {HTMLElement} el   - container element
 */
function renderPanel(raw, el) {
  const dirty = md.parse(raw);
  const clean = DOMPurify.sanitize(dirty, {
    ALLOWED_TAGS: [
      'h1','h2','h3','h4','h5','h6','p','br','hr',
      'ul','ol','li','blockquote','pre','code',
      'a','strong','em','del','img','table','thead',
      'tbody','tr','th','td','span','div','details','summary',
    ],
    ALLOWED_ATTR: ['href','src','alt','title','class','id','open'],
    ALLOW_DATA_ATTR: false,
  });
  el.replaceChildren();
  el.insertAdjacentHTML('afterbegin', clean);
}

// Listen for markdown payloads from the service worker
window.addEventListener('message', (e) => {
  if (e.data?.type === 'RENDER_MARKDOWN') {
    renderPanel(e.data.markdown, document.getElementById('content'));
  }
});
```

## 8. Performance Considerations

- **Lazy load highlight.js.** Only import it when the markdown actually
  contains fenced code blocks. Check for triple-backtick before loading.
- **Throttle re-renders.** If markdown updates frequently (e.g., streaming
  LLM output), debounce `renderPanel` calls (150-300ms).
- **Use `marked.parse()` not `marked()`.** The function-call form is
  deprecated in v18+.
- **Pre-instantiate `new Marked()` with extensions.** Creating the
  instance once avoids re-initializing the highlight extension on every
  parse call.
- **DOMPurify is fast.** Typical sanitize calls take under 1ms for
  reasonable document sizes. Do not skip it for "performance."

## 9. Common Pitfalls

| Pitfall | Fix |
|---|---|
| Injecting unsanitized marked output into the DOM | Always wrap with `DOMPurify.sanitize()` before DOM insertion |
| Loading marked/DOMPurify from CDN on extension pages | Bundle locally; CDN scripts violate `script-src 'self'` |
| highlight.js styles not applying inside shadow DOM | Inject stylesheet into the shadow root, not the document head |
| `marked()` deprecation warning in v18+ | Use `marked.parse()` instead |
| Links open inside the extension iframe | Add `target="_blank"` and `rel="noopener noreferrer"` via DOMPurify hook |
| Trusted Types violation in strict CSP environments | Pass `{ RETURN_TRUSTED_TYPE: true }` to DOMPurify |
| Markdown images load HTTP resources | Use a DOMPurify hook to strip non-HTTPS `src` attributes |
| Host page CSS bleeds into rendered markdown | Render inside a shadow DOM with `all: initial` on `:host` |

## 10. Version Reference (as of May 2026)

| Library | Latest | CDN (jsDelivr) |
|---|---|---|
| marked | 18.0.x | `cdn.jsdelivr.net/npm/marked/lib/marked.umd.js` |
| DOMPurify | 3.4.x | `cdn.jsdelivr.net/npm/dompurify/dist/purify.min.js` |
| highlight.js | 11.11.x | `cdn.jsdelivr.net/npm/highlight.js/lib/common.js` |
| marked-highlight | 2.x | `cdn.jsdelivr.net/npm/marked-highlight/lib/index.umd.js` |
