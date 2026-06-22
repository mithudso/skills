<!-- hub-reference-banner -->
> **Reference file — part of the `chrome-extension-expert` hub.** Formerly standalone `shadow-dom-component-authoring` skill.
> Sibling topics now reference files under hubs (`chrome-extension-expert`) — **not** standalone skills. Ignore "use the X skill" / `related_skills` / SKIP pointers naming bare sibling skills; load topic's `references/<name>.md` from owning hub (see hub's "Cross-hub map").

---

---
name: shadow-dom-component-authoring
version: "1.1.0"
updated: "2026-05-29"
description: >
  Shadow DOM component authoring for Chrome extension overlays and web components.
  TRIGGER: user building shadow DOM components, attaching shadow roots (open/closed),
  styling with adoptedStyleSheets or :host/:host-context/::slotted, composing with slots,
  debugging event retargeting or composedPath, theming across shadow boundaries with CSS
  custom properties, implementing custom element lifecycle callbacks, using declarative
  shadow DOM (shadowrootmode), or isolating Chrome extension content script overlays from
  host page CSS/JS. SKIP: React/Vue component styling (use frontend-design), general CSS
  reference (use html-css), Chrome extension messaging (use extension-message-bridge).
origin: local
whenToUse:
  - "attachShadow open vs closed mode"
  - "adoptedStyleSheets for shadow DOM styling"
  - ":host, :host-context, ::slotted selectors"
  - "Slot composition and slotchange events"
  - "Event retargeting and composedPath debugging"
  - "CSS custom properties crossing shadow boundaries"
  - "Custom element connectedCallback / disconnectedCallback"
  - "Declarative shadow DOM with shadowrootmode"
  - "Chrome extension content script overlay isolation"
  - "all: initial for host page CSS reset"
  - "::part() for external styling of shadow internals"
  - "iframe + shadow DOM hybrid for maximum isolation"
tags:
  - shadow-dom
  - web-components
  - custom-elements
  - chrome-extension
  - content-script
  - css-isolation
  - adoptedStyleSheets
  - slots
related_skills:
  - chrome-dev
  - chrome-extension-security-reviewer
  - html-css
  - vanilla-js-ui-reviewer
  - extension-message-bridge
  - dom-scraping-resilience
---

# Shadow DOM Component Authoring for Chrome Extension Overlays

Expert reference for encapsulated UI components using shadow DOM. Focus: Chrome extension content script overlays.

## When NOT to use this skill

- **React/Vue component styling:** use `frontend-design`.
- **CSS/HTML reference:** use `html-css`.
- **Chrome extension messaging (postMessage, runtime.sendMessage):** use `extension-message-bridge`.
- **Content script DOM scraping/extraction:** use `dom-scraping-resilience`.

## Skill guidance

- Extension overlays: always use **closed mode + `all: initial`** as baseline.
- Reusable web components (non-extension): **open mode** is standard default.
- Prefer `adoptedStyleSheets` over `<style>` when sharing styles across multiple shadow roots.
- Never use deprecated `/deep/` or `::shadow` selectors.
- Event handling across shadow boundaries: always check `composed` property.

---

## 1. attachShadow: open vs closed mode

### Open mode

Exposes shadow root via `element.shadowRoot`. External JS has full read/write access.

```js
class InfoCard extends HTMLElement {
  constructor() {
    super();
    const shadow = this.attachShadow({ mode: 'open' });
    shadow.innerHTML = `
      <style>:host { display: block; padding: 16px; } .title { font-weight: 700; }</style>
      <div class="title"><slot name="title">Default Title</slot></div>
      <div class="body"><slot></slot></div>
    `;
  }
}
customElements.define('info-card', InfoCard);

// External code CAN access internals:
document.querySelector('info-card').shadowRoot.querySelector('.title');
```

**Use open mode for:** reusable web components, design systems, test frameworks, accessibility auditors.

### Closed mode

Returns `null` from `element.shadowRoot`. Only code holding reference from `attachShadow()` can access tree.

```js
class SecureOverlay extends HTMLElement {
  #shadow; // private — only this class can access

  constructor() {
    super();
    this.#shadow = this.attachShadow({ mode: 'closed' });
    const styles = new CSSStyleSheet();
    styles.replaceSync(':host { all: initial; display: block; }');
    this.#shadow.adoptedStyleSheets = [styles];

    const root = document.createElement('div');
    root.className = 'overlay-root';
    this.#shadow.appendChild(root);
  }

  updateContent(text) {
    this.#shadow.querySelector('.overlay-root').textContent = text;
  }
}
customElements.define('secure-overlay', SecureOverlay);

document.querySelector('secure-overlay').shadowRoot; // null
```

**Use closed mode for:** Chrome extension content script overlays, security-sensitive UI.

**Closed mode limitations:**
- `element.shadowRoot` returns `null` — no external `querySelector`
- DevTools Elements panel still shows tree (closed blocks scripts, not user)
- `event.composedPath()` still reveals full path to receiving listeners

### Mode comparison

| Criterion | Open | Closed |
|-----------|------|--------|
| External JS via `.shadowRoot` | Yes | No (null) |
| DevTools visibility | Full | Full |
| Test framework access | Direct | Requires exposed API |
| Host page CSS interference | Blocked | Blocked |
| Host page JS interference | Possible | Blocked |
| Chrome extension overlays | Discouraged | Recommended |
| Reusable web components | Recommended | Unusual |

---

## 2. adoptedStyleSheets and CSSStyleSheet

### Constructable stylesheets

Create once, share across many shadow roots. Highest-perf approach for shared styles.

```js
const sharedStyles = new CSSStyleSheet();
sharedStyles.replaceSync(`
  :host { all: initial; display: block; font-family: system-ui, sans-serif; }
  .panel { background: #1e293b; color: #e2e8f0; border-radius: 12px; padding: 16px; }
  .btn  { background: #2563eb; color: #fff; border: none; padding: 8px 16px;
          border-radius: 6px; cursor: pointer; font-size: 14px; }
  .btn:hover { background: #1d4ed8; }
`);

class OverlayPanel extends HTMLElement {
  constructor() {
    super();
    const shadow = this.attachShadow({ mode: 'closed' });
    shadow.adoptedStyleSheets = [sharedStyles]; // shared reference — zero extra parse cost
    const panel = document.createElement('div');
    panel.className = 'panel';
    shadow.appendChild(panel);
  }
}
```

### Updating adopted stylesheets

```js
await sharedStyles.replace(':host { color: blue; }');     // async — non-blocking
sharedStyles.replaceSync(':host { color: red; }');         // sync — blocks briefly
sharedStyles.insertRule('.new { display: flex; }', sharedStyles.cssRules.length);
sharedStyles.deleteRule(0);
```

### adoptedStyleSheets vs `<style>` element

| Feature | adoptedStyleSheets | `<style>` element |
|---------|-------------------|-------------------|
| Shared across shadow roots | Yes (same object ref) | No (duplicated) |
| Parse cost | Once per sheet | Once per shadow root |
| Dynamic updates | `replace` / `replaceSync` | Modify `textContent` |
| Declarative shadow DOM | Not yet | Yes |
| Browser support | Chrome 73+, FF 101+, Safari 16.4+ | Universal |

**Rule:** use `adoptedStyleSheets` for extension overlays and component libraries; `<style>` for one-off components or declarative shadow DOM.

---

## 3. CSS selectors inside shadow DOM

### :host — style the host element

```css
:host { all: initial; display: block; position: fixed; bottom: 16px; right: 16px;
        z-index: 2147483646; font-family: Inter, system-ui, sans-serif; }

:host(.expanded)           { width: 400px; height: 600px; }
:host([data-collapsed])    { width: 48px; height: 48px; border-radius: 50%; }
```

### :host-context — style based on ancestor

Adapts component to external context without breaking encapsulation.

```css
:host-context(.dark-theme)  { --bg: #0f172a; --fg: #e2e8f0; }
:host-context(.light-theme) { --bg: #ffffff; --fg: #1e293b; }
```

**Browser support:** Chromium (always), Firefox 2024+, Safari 2025+.

### ::slotted — style projected light DOM content

Targets only top-level slotted elements. Cannot descend into children.

```css
::slotted(*)   { margin: 0; padding: 8px; }
::slotted(h2)  { font-size: 18px; font-weight: 700; color: var(--heading-color, #f8fafc); }

/* INVALID — cannot descend into slotted children: */
/* ::slotted(div) span { } */
```

### ::part — expose internals for external styling

```html
<my-button>
  <template shadowrootmode="open">
    <button part="base" class="btn"><slot></slot></button>
  </template>
</my-button>
```

```css
/* In host page or parent component */
my-button::part(base) { background: hotpink; border-radius: 999px; }
```

Use `::part()` for public styling API without exposing full shadow tree access.

---

## 4. Slot composition

### Named and default slots

```js
class TabPanel extends HTMLElement {
  constructor() {
    super();
    const shadow = this.attachShadow({ mode: 'open' });
    const styles = new CSSStyleSheet();
    styles.replaceSync(`
      .tabs { display: flex; gap: 4px; border-bottom: 1px solid #334155; }
      ::slotted([slot="tab"])        { padding: 8px 16px; cursor: pointer; background: transparent; }
      ::slotted([slot="tab"].active) { color: #fff; border-bottom: 2px solid #2563eb; }
    `);
    shadow.adoptedStyleSheets = [styles];

    shadow.innerHTML = `
      <div class="tabs"><slot name="tab"></slot></div>
      <div class="content"><slot></slot></div>
    `;
  }
}
customElements.define('tab-panel', TabPanel);
```

```html
<tab-panel>
  <button slot="tab" class="active">Info</button>
  <button slot="tab">History</button>
  <div>Default slot content.</div>
</tab-panel>
```

### Slot change detection

```js
const slot = shadow.querySelector('slot[name="tab"]');
slot.addEventListener('slotchange', () => {
  const assigned = slot.assignedElements();
  console.log('Tab slot now has', assigned.length, 'elements');
});
```

```js
slot.assignedElements();               // elements in this slot
slot.assignedNodes({ flatten: true }); // all nodes including text, follows fallback
```

---

## 5. Event retargeting and composed events

### How retargeting works

Events inside shadow tree retarget `event.target` to host element when bubbling past shadow root. External listeners never see internal details.

```js
hostElement.addEventListener('click', (e) => {
  console.log(e.target);         // host element (retargeted)
  console.log(e.composedPath()); // full path including shadow internals
});
```

### The `composed` property

Events need `composed: true` to cross shadow boundaries. Native UI events composed by default. **Custom events NOT composed by default.**

```js
// Stays inside shadow root
shadow.querySelector('.btn').dispatchEvent(
  new CustomEvent('internal-action', { bubbles: true, composed: false })
);

// Crosses shadow boundary
shadow.querySelector('.btn').dispatchEvent(
  new CustomEvent('overlay-action', {
    bubbles: true, composed: true,
    detail: { action: 'close', caseNumber: '12345' }
  })
);
```

### Composed events reference

| Event | composed | Notes |
|-------|----------|-------|
| click, mousedown/up, dblclick | true | All mouse events cross |
| pointerdown/up/move | true | All pointer events cross |
| keydown, keyup | true | All keyboard events cross |
| input, change | true | Form events cross |
| focus, blur, focusin, focusout | true | All focus events cross |
| scroll | false | Stays inside shadow root |
| Custom events | false (default) | Must set `composed: true` explicitly |
| slotchange | false | Stays inside shadow root |

**Slotted element note:** events on slotted elements (light DOM) NOT retargeted — elements physically live in light DOM.

---

## 6. CSS custom properties crossing shadow boundaries

CSS custom properties inherited, cross shadow DOM boundaries naturally — primary mechanism for external theming.

```js
const styles = new CSSStyleSheet();
styles.replaceSync(`
  :host {
    /* Defaults — overridable from outside */
    --overlay-bg: #0f172a;
    --overlay-fg: #e2e8f0;
    --overlay-accent: #2563eb;
    --overlay-radius: 14px;
  }
  .shell {
    background: var(--overlay-bg);
    color: var(--overlay-fg);
    border-radius: var(--overlay-radius);
  }
`);
```

```css
/* Host page overrides */
secure-overlay { --overlay-bg: #ffffff; --overlay-accent: #16a34a; }
.dark-mode secure-overlay { --overlay-bg: #020617; }
```

### Inheritance boundary

| Property type | Crosses shadow boundary? |
|---------------|------------------------|
| CSS custom properties (`--*`) | Yes (inherited) |
| `color`, `font-family`, `font-size` | Yes — blocked by `all: initial` |
| `background`, `border`, `margin`, `padding` | No (not inherited) |
| `display`, `position`, `width`, `height` | No (not inherited) |

**Rule:** use `all: initial` on `:host` for extension overlays, then expose `--*` custom properties as intentional theming contract.

---

## 7. Custom elements lifecycle

```js
class CaseOverlay extends HTMLElement {
  static observedAttributes = ['case-id', 'collapsed'];
  #shadow;

  constructor() {
    super();
    this.#shadow = this.attachShadow({ mode: 'closed' });
    this.#shadow.adoptedStyleSheets = [sharedStyles];
    this.#shadow.appendChild(document.createElement('div')).className = 'root';
  }

  connectedCallback() {
    // Added to DOM — safe to read attributes and set up listeners
    this.#render();
    this.#shadow.querySelector('.close-btn')
      ?.addEventListener('click', () => this.remove());
  }

  disconnectedCallback() {
    // Removed from DOM — clean up timers, observers, listeners
    clearInterval(this._refreshTimer);
  }

  attributeChangedCallback(name, oldVal, newVal) {
    if (oldVal === newVal) return;
    if (name === 'collapsed') this.#toggleCollapsed(newVal === 'true');
    if (name === 'case-id') this.#render();
  }

  adoptedCallback() {
    // Element moved to a new document (rare — iframe adoption)
  }

  #render() { /* build DOM */ }
  #toggleCollapsed(collapsed) {
    this.#shadow.querySelector('.root')?.classList.toggle('collapsed', collapsed);
  }
}
customElements.define('case-overlay', CaseOverlay);
```

### Lifecycle ordering

1. `constructor()` — attach shadow, set up initial DOM structure
2. `connectedCallback()` — added to document; safe to measure/render
3. `attributeChangedCallback()` — fires for each observed attribute
4. `disconnectedCallback()` — removed; clean up
5. `adoptedCallback()` — moved between documents (rare)

**Rule:** never read layout-dependent values (`offsetWidth`, `getBoundingClientRect`) in constructor. Wait for `connectedCallback`.

---

## 8. Declarative shadow DOM

```html
<info-card>
  <template shadowrootmode="open">
    <style>:host { display: block; } .title { font-weight: 700; }</style>
    <div class="title"><slot name="title">Untitled</slot></div>
    <div class="body"><slot></slot></div>
  </template>
  <span slot="title">Server Status</span>
  <p>All systems operational.</p>
</info-card>
```

Browser parses template, creates shadow root, removes template element — no JS needed for initial render.

**Hydration with custom elements:**
```js
class InfoCard extends HTMLElement {
  constructor() {
    super();
    if (!this.shadowRoot) {
      // No declarative root — create imperatively (fallback)
      this.attachShadow({ mode: 'open' });
    }
    // shadowRoot is ready either way
  }
}
```

**When declarative shadow DOM matters:**
- Server-side rendering — eliminates FOUC
- Static site generators — shadow DOM without JS bundle
- Progressive enhancement — structure first, interactivity on hydration

**Not applicable for Chrome extension overlays** — content scripts inject imperatively.

---

## 9. Chrome extension overlay isolation

### Standard isolation pattern

```js
(function injectOverlay() {
  const HOST_ID = 'my-extension-overlay-host';
  if (document.getElementById(HOST_ID)) return; // prevent double-injection

  const host = document.createElement('div');
  host.id = HOST_ID;
  const shadow = host.attachShadow({ mode: 'closed' }); // CLOSED mode

  const styles = new CSSStyleSheet();
  styles.replaceSync(`
    :host {
      all: initial !important;
      position: fixed !important;
      bottom: 16px !important;
      right: 16px !important;
      z-index: 2147483646 !important;
      font-family: system-ui, -apple-system, sans-serif !important;
      font-size: 14px !important;
      line-height: 1.5 !important;
      direction: ltr !important;
    }
    *, *::before, *::after { box-sizing: border-box; }
    .shell {
      width: 360px; height: 480px;
      background: #0f172a; border-radius: 14px;
      border: 1px solid rgba(148,163,184,0.2);
      box-shadow: 0 20px 50px rgba(0,0,0,0.35);
      display: flex; flex-direction: column; overflow: hidden;
    }
  `);
  shadow.adoptedStyleSheets = [styles];

  const shell = document.createElement('div');
  shell.className = 'shell';
  shadow.appendChild(shell);
  document.documentElement.appendChild(host);
})();
```

### Why `all: initial` is mandatory

Host pages routinely set global resets that break injected UI:

```css
html { font-size: 62.5%; }              /* Shrinks rem-based sizing */
* { box-sizing: content-box; }          /* Breaks flex layouts */
div { display: inline; }               /* Breaks block layout */
button { all: unset; }                  /* Strips button styling */
```

`all: initial` resets every CSS property. Use `!important` on `:host` in extension overlays to guard against host page `!important` rules on inherited properties.

### iframe + shadow DOM hybrid (maximum isolation)

```js
function injectOverlayWithIframe(panelUrl) {
  const host = document.createElement('div');
  const shadow = host.attachShadow({ mode: 'closed' });

  const styles = new CSSStyleSheet();
  styles.replaceSync(`
    :host { all: initial !important; position: fixed !important;
            bottom: 16px !important; right: 16px !important;
            z-index: 2147483646 !important; }
    .shell { width: 360px; height: 480px; border-radius: 14px; overflow: hidden;
             display: flex; flex-direction: column; background: #0f172a; }
    .header { height: 40px; background: #1e293b; display: flex;
              align-items: center; padding: 0 12px; color: #f8fafc; font: 700 13px system-ui; }
    iframe { flex: 1; border: none; width: 100%; }
  `);
  shadow.adoptedStyleSheets = [styles];

  const shell = document.createElement('div');
  shell.className = 'shell';
  const header = document.createElement('div');
  header.className = 'header';
  header.textContent = 'My Extension';
  const iframe = document.createElement('iframe');
  iframe.src = panelUrl;
  iframe.setAttribute('sandbox', 'allow-scripts allow-same-origin');
  shell.append(header, iframe);
  shadow.appendChild(shell);

  window.addEventListener('message', (e) => {
    if (e.source === iframe.contentWindow) handlePanelMessage(e.data);
  });

  document.documentElement.appendChild(host);
}
```

### Re-injection guard

```js
const observer = new MutationObserver(() => {
  if (!document.getElementById(HOST_ID)) document.documentElement.appendChild(host);
});
observer.observe(document.documentElement, { childList: true });
```

### Communication patterns

| Channel | Direction | When to use |
|---------|-----------|-------------|
| `chrome.runtime.sendMessage` | content/popup → service worker | Message routing to background |
| `postMessage` | parent ↔ iframe panel | Cross-origin iframe communication |
| `CustomEvent` on host element | shadow DOM → light DOM | Notify parent context of shadow events |
| DOM attributes (`data-*`) | light DOM → shadow DOM | Pass serialized context into component |
| `chrome.storage.session` | Any extension context | Shared volatile state |

### Common pitfalls and fixes

| Problem | Fix |
|---------|-----|
| Host page `font-size: 62.5%` shrinks overlay text | `all: initial` + explicit `font-size: 14px !important` |
| Host page `* { box-sizing: content-box }` breaks layout | `*, *::before, *::after { box-sizing: border-box; }` inside shadow |
| Host page stacking context hides overlay | `z-index: 2147483646` + `position: fixed` on `:host` |
| Click events caught by host page listener | `e.stopPropagation()` on shadow root for internal-only events |
| Host page `MutationObserver` removes overlay | Re-injection guard (pattern above) |
| Extension styles leak into host page | All styles must live inside shadow root — never `<style>` to host document |

---

## 10. Anti-patterns

1. **Open mode for extension overlays.** Host scripts can read/mutate via `element.shadowRoot`.
2. **Omitting `all: initial` on `:host`.** Host page inherited properties leak in.
3. **Deprecated `::shadow` or `/deep/` combinators.** Removed from all browsers.
4. **Custom events with `composed: true` by default.** Only set when event genuinely needs to cross boundary.
5. **Layout reads in `constructor`.** `offsetWidth`, `getBoundingClientRect` return zero before connection. Use `connectedCallback`.
6. **Missing `disconnectedCallback` cleanup.** Timers, observers, listeners must be cleaned up.
7. **`::slotted(div) span`.** `::slotted()` targets only top-level slotted elements.
8. **`<style>` added to host page document.** Extension styles must live exclusively inside shadow root.
9. **Closed mode as security sandbox.** DevTools still shows tree; `composedPath()` still exposes internal nodes.
10. **`innerHTML` with user-controlled strings.** Build DOM with `createElement`/`textContent`. Raw string injection is XSS vector even inside shadow DOM.

---

## 11. Review checklist

- [ ] `attachShadow({ mode })` — closed for extension overlays, open for reusable components
- [ ] `:host { all: initial; }` present for content script overlays
- [ ] `!important` on `:host` position/z-index/font-size in extension contexts
- [ ] `adoptedStyleSheets` used when sharing styles across multiple shadow roots
- [ ] No deprecated `::shadow` or `/deep/` selectors
- [ ] Custom events use `composed: true` only when cross-boundary propagation is intentional
- [ ] CSS custom properties documented as public theming API
- [ ] `connectedCallback` used for DOM-dependent setup (not `constructor`)
- [ ] `disconnectedCallback` cleans up timers, observers, and listeners
- [ ] `*, *::before, *::after { box-sizing: border-box }` set inside shadow root
- [ ] iframe `sandbox` attribute set appropriately for hybrid patterns
- [ ] Re-injection guard present if host page may remove overlay host element
- [ ] DOM built with safe methods (`createElement`/`textContent`), not raw `innerHTML`
- [ ] `::part()` used to expose intended styling hooks instead of full open access