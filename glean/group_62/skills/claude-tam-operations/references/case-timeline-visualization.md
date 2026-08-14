<!-- hub-reference-banner -->
> **Reference file — part of the `tam-operations` hub.** Formerly the standalone `case-timeline-visualization` skill.
> Sibling topics in this family are now reference files under the hubs (`tam-operations`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: case-timeline-visualization
version: 1.1.1
updated: "2026-06-01"
description: >
  Rendering temporal event sequences in vanilla JS for dashboards — data modeling,
  vertical/horizontal layouts, CSS timeline patterns, SVG vs DOM rendering, zoom/scroll
  interactivity, accessibility, and performance.

  TRIGGER: user is building or reviewing a timeline or event-sequence visualization in vanilla
  JavaScript, choosing between DOM/SVG/Canvas rendering for temporal data, adding zoom/pan/
  scroll/filtering to a timeline, modeling event data structures for chronological display, or
  embedding a timeline inside a Chrome extension dashboard, iframe, or shadow DOM.

  SKIP: React/Vue/Svelte timeline components (use the relevant framework skill), general charting
  (bar/line/pie — use da-applied-and-communication), CSS-only decorative timelines with no data
  binding (use frontend-ui), or Gantt chart/project scheduling.
origin: local
category: developer
tags:
  - timeline
  - visualization
  - vanilla-js
  - dom
  - svg
  - canvas
  - accessibility
  - chrome-extension
triggers:
  - timeline visualization
  - event sequence
  - temporal data
  - timeline component
  - horizontal timeline
  - vertical timeline
  - case timeline
  - zoom pan timeline
related_skills:
  - case-tracker
  - frontend-ui
whenToUse:
  - "build a case event history timeline in vanilla JS"
  - "choose between DOM, SVG, and Canvas for a timeline component"
  - "add zoom, pan, or filtering to an event timeline"
  - "model event data for chronological display"
  - "embed a timeline in a Chrome extension dashboard or shadow DOM"
whenNotToUse:
  - "React/Vue/Svelte timeline components — use the relevant framework skill"
  - "general bar/line/pie charting — use da-applied-and-communication"
  - "CSS-only decorative timelines with no data binding — use frontend-ui"
  - "Gantt chart or project scheduling — use a dedicated scheduling skill"
---

# Case Timeline Visualization

Expert reference for building temporal event-sequence timelines in vanilla JavaScript and CSS. Covers data modeling, rendering strategies, layout orientation, interactivity (zoom, scroll, filtering), accessibility, and performance. Designed for support-case event histories, incident timelines, and operational audit trails.

---

## 1. Event data model

### Base event schema

```js
/**
 * @typedef {Object} TimelineEvent
 * @property {string}  id             - Unique identifier
 * @property {string}  timestamp      - ISO 8601 string ("2026-05-20T14:32:00Z")
 * @property {string}  label          - Short display text
 * @property {string}  [detail]       - Optional longer description
 * @property {string}  [type]         - Category key for color/icon mapping
 * @property {string}  [severity]     - "info" | "warning" | "error" | "critical"
 * @property {string}  [actor]        - Who or what caused this event
 * @property {Object}  [meta]         - Arbitrary payload for tooltips/detail panes
 * @property {string}  [endTimestamp] - For duration events (spans)
 */
```

### Data preparation helpers

```js
function sortEvents(events) {
  return [...events].sort((a, b) => new Date(a.timestamp) - new Date(b.timestamp));
}

function groupByDate(events) {
  const groups = new Map();
  for (const evt of events) {
    const key = evt.timestamp.slice(0, 10);
    if (!groups.has(key)) groups.set(key, []);
    groups.get(key).push(evt);
  }
  return groups;
}

function timeToPosition(ts, startMs, endMs, axisLengthPx) {
  const ratio = (new Date(ts).getTime() - startMs) / (endMs - startMs);
  return Math.round(ratio * axisLengthPx);
}

function isDurationEvent(evt) {
  return Boolean(evt.endTimestamp);
}
```

**Rule:** Store timestamps as ISO 8601 strings in the model. Create `Date` objects only at render time — `Date` objects are not serializable.

---

## 2. Rendering strategy: DOM vs SVG vs Canvas

| Criterion | DOM (HTML/CSS) | SVG | Canvas |
|-----------|---------------|-----|--------|
| Element count limit | ~5,000 nodes OK | ~3,000–5,000 nodes | Unlimited (bitmap) |
| Styling flexibility | Full CSS cascade | CSS + SVG attributes | Programmatic only |
| Hit detection | Native events | Native events | Manual (coordinates) |
| Accessibility | Native semantics | `<title>`/`<desc>` | Requires ARIA overlay |
| Text rendering | Excellent | Excellent | Blurry at low DPI |
| Animation | CSS transitions | SMIL / CSS | requestAnimationFrame |
| Debuggability | DevTools DOM panel | DevTools DOM panel | Opaque bitmap |

**Decision rule:**
- Start with DOM/CSS for typical case timelines (tens to low hundreds of events).
- Switch to SVG when you need precise geometric drawing (connector lines, arcs, span bars).
- Switch to Canvas only when DOM/SVG performance degrades (thousands of visible nodes).

---

## 3. Vertical timeline (DOM/CSS)

The most common layout for case event histories. Events stack top-to-bottom chronologically with a left-aligned spine.

### HTML structure

```html
<div class="tl" role="list" aria-label="Case timeline">
  <div class="tl-item" role="listitem" data-type="comment">
    <div class="tl-marker" aria-hidden="true"></div>
    <div class="tl-content">
      <time class="tl-time" datetime="2026-05-20T14:32:00Z">May 20, 2:32 PM</time>
      <p class="tl-label">Customer replied with diagnostic logs</p>
    </div>
  </div>
</div>
```

### Core CSS

```css
.tl {
  --tl-line-color: #d1d5db;
  --tl-marker-size: 12px;
  --tl-marker-color: #3b82f6;
  --tl-gap: 24px;
  --tl-line-width: 2px;
  position: relative;
  padding-left: calc(var(--tl-marker-size) / 2 + 16px);
  display: flex;
  flex-direction: column;
  gap: var(--tl-gap);
}

.tl::before {
  content: "";
  position: absolute;
  left: calc(var(--tl-marker-size) / 2);
  top: 0; bottom: 0;
  width: var(--tl-line-width);
  background: var(--tl-line-color);
}

.tl-item { position: relative; padding-left: 20px; }

.tl-marker {
  position: absolute;
  left: calc(-16px - var(--tl-marker-size) / 2);
  top: 4px;
  width: var(--tl-marker-size); height: var(--tl-marker-size);
  border-radius: 50%;
  background: var(--tl-marker-color);
  border: 2px solid #fff;
  box-shadow: 0 0 0 2px var(--tl-line-color);
  z-index: 1;
}

.tl-time { display: block; font-size: 0.75rem; color: #6b7280; margin-bottom: 2px; }
.tl-label { margin: 0; font-size: 0.875rem; line-height: 1.4; }

.tl-item[data-type="error"] .tl-marker   { background: #ef4444; }
.tl-item[data-type="warning"] .tl-marker  { background: #f59e0b; }
.tl-item[data-type="success"] .tl-marker  { background: #10b981; }
.tl-item[data-type="comment"] .tl-marker  { background: #3b82f6; }
.tl-item[data-type="system"] .tl-marker   { background: #6b7280; }

@media (prefers-reduced-motion: reduce) {
  .tl-item, .htl-event, .tl-filter-btn { transition: none; animation: none; }
}
```

### Rendering function

```js
function renderVerticalTimeline(container, events) {
  const sorted = sortEvents(events);
  if (!sorted.length) {
    const empty = document.createElement("p");
    empty.className = "tl-empty";
    empty.textContent = "No events to display.";
    container.replaceChildren(empty);
    return;
  }

  const frag = document.createDocumentFragment();
  const list = document.createElement("div");
  list.className = "tl";
  list.setAttribute("role", "list");
  list.setAttribute("aria-label", "Case timeline");

  for (const evt of sorted) {
    const item = document.createElement("div");
    item.className = "tl-item";
    item.setAttribute("role", "listitem");
    item.dataset.type = evt.type || "info";

    const marker = document.createElement("div");
    marker.className = "tl-marker";
    marker.setAttribute("aria-hidden", "true");

    const content = document.createElement("div");
    content.className = "tl-content";

    const time = document.createElement("time");
    time.className = "tl-time";
    time.setAttribute("datetime", evt.timestamp);
    time.textContent = formatTimestamp(evt.timestamp);

    const label = document.createElement("p");
    label.className = "tl-label";
    label.textContent = evt.label;  // textContent, not innerHTML — XSS prevention

    content.appendChild(time);
    content.appendChild(label);

    if (evt.detail) {
      const detail = document.createElement("p");
      detail.className = "tl-detail";
      detail.textContent = evt.detail;
      content.appendChild(detail);
    }

    item.appendChild(marker);
    item.appendChild(content);
    list.appendChild(item);
  }

  frag.appendChild(list);
  container.replaceChildren(frag);  // single DOM write
}

function formatTimestamp(iso) {
  const d = new Date(iso);
  return d.toLocaleDateString(undefined, { month: "short", day: "numeric" }) + ", " +
         d.toLocaleTimeString(undefined, { hour: "numeric", minute: "2-digit" });
}
```

---

## 4. Horizontal timeline (DOM/CSS)

Best for compact overviews showing event distribution across a date range.

### Core CSS

```css
.htl-wrapper {
  overflow-x: auto; overflow-y: hidden;
  position: relative;
  border: 1px solid #e5e7eb; border-radius: 6px;
  padding: 32px 16px 48px;
}
.htl-axis { position: relative; height: 4px; background: #d1d5db; border-radius: 2px; }
.htl-event {
  position: absolute; top: -8px;
  width: 16px; height: 16px; border-radius: 50%;
  background: #3b82f6; border: 2px solid #fff;
  box-shadow: 0 1px 3px rgba(0,0,0,0.2);
  transform: translateX(-50%);
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}
.htl-event:hover, .htl-event:focus-visible {
  transform: translateX(-50%) scale(1.3);
  box-shadow: 0 2px 8px rgba(0,0,0,0.25); z-index: 2;
}
.htl-span {
  position: absolute; top: -4px; height: 12px; border-radius: 6px;
  background: rgba(59, 130, 246, 0.3); border: 1px solid rgba(59, 130, 246, 0.6);
}
```

---

## 5. Zoom and scroll interactivity

### Wheel-zoom

```js
function attachZoom(wrapper, state, rerender) {
  wrapper.addEventListener("wheel", (e) => {
    if (!e.ctrlKey && !e.metaKey) return;
    e.preventDefault();
    const zoomFactor = e.deltaY > 0 ? 0.8 : 1.25;
    const rect = wrapper.getBoundingClientRect();
    const cursorRatio = (e.clientX - rect.left) / rect.width;
    const rangeMs = state.endMs - state.startMs;
    const newRange = rangeMs / zoomFactor;
    const cursorMs = state.startMs + rangeMs * cursorRatio;
    state.startMs = cursorMs - newRange * cursorRatio;
    state.endMs = cursorMs + newRange * (1 - cursorRatio);
    rerender();
  }, { passive: false });
}
```

### Drag-to-pan

```js
function attachPan(wrapper, state, rerender) {
  let dragging = false, lastX = 0;
  wrapper.addEventListener("pointerdown", (e) => {
    if (e.button !== 0) return;
    dragging = true; lastX = e.clientX;
    wrapper.setPointerCapture(e.pointerId);
    wrapper.style.cursor = "grabbing";
  });
  wrapper.addEventListener("pointermove", (e) => {
    if (!dragging) return;
    const dx = e.clientX - lastX; lastX = e.clientX;
    const pxPerMs = wrapper.offsetWidth / (state.endMs - state.startMs);
    const shiftMs = dx / pxPerMs;
    state.startMs -= shiftMs; state.endMs -= shiftMs;
    rerender();
  });
  wrapper.addEventListener("pointerup", () => {
    dragging = false; wrapper.style.cursor = "";
  });
}
```

### Debounced re-render

```js
function debounceRAF(fn) {
  let frameId = null;
  return (...args) => {
    if (frameId) cancelAnimationFrame(frameId);
    frameId = requestAnimationFrame(() => { fn(...args); frameId = null; });
  };
}
// const rerender = debounceRAF(() => renderHorizontalTimeline(...));
```

---

## 6. Event filtering

```js
function renderFilterBar(container, types, activeTypes, onChange) {
  const bar = document.createElement("div");
  bar.className = "tl-filters";
  bar.setAttribute("role", "group");
  bar.setAttribute("aria-label", "Filter events by type");

  for (const type of types) {
    const btn = document.createElement("button");
    btn.className = "tl-filter-btn";
    btn.textContent = type;
    btn.dataset.type = type;
    btn.setAttribute("aria-pressed", activeTypes.has(type) ? "true" : "false");
    btn.addEventListener("click", () => {
      if (activeTypes.has(type)) activeTypes.delete(type); else activeTypes.add(type);
      btn.setAttribute("aria-pressed", activeTypes.has(type) ? "true" : "false");
      onChange(activeTypes);
    });
    bar.appendChild(btn);
  }
  container.prepend(bar);
}
```

---

## 7. Accessibility checklist

1. Timeline container: `role="list"`. Each event: `role="listitem"`. Decorative elements: `aria-hidden="true"`.
2. Interactive elements are `<button>` or have `tabindex="0"` and keyboard handlers.
3. Each event node needs `aria-label` combining label and formatted timestamp.
4. Color is not the sole type indicator — add icon, text badge, or pattern fill (WCAG 1.4.1).
5. Wrap transition/animation CSS in `@media (prefers-reduced-motion: no-preference)`.
6. Filter toggle buttons use `aria-pressed`.
7. Tooltips use `role="tooltip"` associated via `aria-describedby` on the trigger.
8. After zoom changes, update an `aria-live="polite"` region with the new visible date range.

### Arrow-key navigation

```js
function attachArrowNav(axisEl) {
  axisEl.addEventListener("keydown", (e) => {
    if (e.key !== "ArrowLeft" && e.key !== "ArrowRight") return;
    const dots = [...axisEl.querySelectorAll(".htl-event")];
    const idx = dots.indexOf(document.activeElement);
    if (idx === -1) return;
    e.preventDefault();
    const next = e.key === "ArrowRight" ? Math.min(idx + 1, dots.length - 1) : Math.max(idx - 1, 0);
    dots[next].focus();
  });
}
```

---

## 8. Performance guidelines

- **Build off-DOM.** Construct the entire timeline in a `DocumentFragment`, then call `container.replaceChildren(frag)` once. Avoids incremental reflow.
- **Limit visible nodes.** For 500+ events, virtualize: render only events within the visible scroll window.
- **Batch reads/writes.** Never interleave `getBoundingClientRect()` calls with style mutations.
- **Use CSS containment.** Apply `contain: content` to timeline items to isolate layout recalculation.
- **SVG limit:** ~3,000–5,000 elements. Beyond that, switch to Canvas or cluster events.

### Type-to-color mapping

```js
const TYPE_COLORS = {
  error: "#ef4444", warning: "#f59e0b", success: "#10b981",
  comment: "#3b82f6", system: "#6b7280", info: "#8b5cf6",
};
function getTypeColor(type) { return TYPE_COLORS[type] || TYPE_COLORS.info; }
```

---

## 9. Dark mode and theming

```css
.tl, .htl-wrapper {
  --tl-bg: #ffffff; --tl-text: #1f2937;
  --tl-text-muted: #6b7280; --tl-line-color: #d1d5db;
}
@media (prefers-color-scheme: dark) {
  .tl, .htl-wrapper {
    --tl-bg: #111827; --tl-text: #f9fafb;
    --tl-text-muted: #9ca3af; --tl-line-color: #374151;
  }
}
```

For class-based theme toggle: replace `@media` with `.dark .tl, .dark .htl-wrapper { ... }`.

---

## 10. Shadow DOM embedding

When injecting into a host page (Chrome extension content script, third-party embed):

```js
function mountTimeline(hostEl, events) {
  const shadow = hostEl.attachShadow({ mode: "open" });
  const style = document.createElement("style");
  style.textContent = `/* all .tl-* CSS rules */`;
  shadow.appendChild(style);
  const root = document.createElement("div");
  shadow.appendChild(root);
  renderVerticalTimeline(root, events);
}
```

---

## 11. Responsive design

```css
@media (max-width: 768px) { .htl-wrapper { display: none; } .tl-vertical-fallback { display: flex; } }
@media (min-width: 769px) { .tl-vertical-fallback { display: none; } }

/* Container queries (modern browsers) */
@container timeline-host (max-width: 400px) {
  .tl-label { font-size: 0.75rem; }
  .tl-detail { display: none; }
}
```

Apply `container-type: inline-size; container-name: timeline-host;` on the timeline host element.

---

## 12. Anti-patterns

1. **String-based HTML insertion for event rendering.** XSS risk if event labels contain user input. Always use `textContent` or explicit DOM construction.
2. **Storing Date objects in the model.** Dates are not serializable. Use ISO 8601 strings; create Date objects only at render time.
3. **Re-rendering the entire DOM on every zoom tick.** Causes layout thrashing. Use `requestAnimationFrame` debouncing.
4. **Color-only event types.** Fails WCAG 1.4.1. Supplement color with icons, text, or patterns.
5. **Omitting keyboard navigation.** Dot markers on a horizontal axis without `tabindex` or button semantics are invisible to keyboard users.
6. **Pixel-absolute positioning without a time-to-position function.** Leads to events overlapping or drifting on container resize.
7. **Heavyweight tooltip libraries.** A single positioned `<div>` with show/hide toggling is sufficient.
8. **Forgetting cleanup.** In extension contexts, timelines may be created and destroyed multiple times. Use `AbortController` for event listeners; call `container.replaceChildren()` on teardown.

---

## 13. Review checklist

- [ ] Event model is flat, serializable, and uses ISO 8601 timestamps
- [ ] DOM is built in a fragment and inserted in a single write
- [ ] No string-based HTML insertion with untrusted content
- [ ] Timeline container has `role="list"`, items have `role="listitem"`
- [ ] Decorative elements have `aria-hidden="true"`
- [ ] Interactive elements are `<button>` or have `tabindex="0"` and keyboard handlers
- [ ] Color is not the sole type indicator
- [ ] Transitions/animations respect `prefers-reduced-motion`
- [ ] Zoom/pan is debounced via `requestAnimationFrame`
- [ ] Position calculations use `timeToPosition`, not hardcoded pixels
- [ ] Responsive behavior handles narrow containers (vertical fallback or scrollable)
- [ ] Event listeners are cleaned up on teardown
- [ ] CSS is scoped (shadow DOM, BEM, or `data-*` selectors) to avoid host page conflicts
