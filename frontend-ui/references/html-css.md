<!-- hub-reference-banner -->
> **Reference file — part of the `frontend-ui` hub.** Formerly the standalone `html-css` skill.
> Sibling topics in this family are now reference files under the hubs (`frontend-ui`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: html-css
version: "1.1.0"
updated: "2026-05-29"
description: >
  HTML (Living Standard) and CSS expert reference for markup structure, semantics, layout,
  specificity, and modern CSS features. TRIGGER: user asks about HTML elements or attributes,
  CSS selectors or properties, flexbox/grid layout, specificity, cascade, responsive design
  with media/container queries, CSS nesting, or document structure. SKIP: React/component
  architecture (use react-nextjs or frontend-design), accessibility audits (use
  accessibility-ux-reviewer), CSS-in-JS or Tailwind strategy (use frontend-design).
origin: local
whenToUse:
  - "How do I use flexbox/grid for layout?"
  - "CSS specificity question"
  - "Semantic HTML element choice"
  - "Responsive design with media queries or container queries"
  - "CSS cascade and inheritance"
  - "Form element markup and labeling"
  - "HTML document structure"
  - "CSS nesting, @layer, @scope"
  - "Box model and positioning"
  - "CSS animations and transitions"
tags:
  - html
  - css
  - html-living-standard
  - flexbox
  - grid
  - responsive-design
  - css-cascade
  - semantic-html
  - web-standards
related_skills:
  - frontend-design
  - accessibility-ux-reviewer
  - shadow-dom-component-authoring
---

# HTML and CSS Expert

Practical reference for HTML markup and CSS styling. Primary source: MDN. Canonical specs:
WHATWG HTML Living Standard and CSSWG drafts.

## When NOT to use this skill

- **Component architecture or CSS strategy (Tailwind, CSS Modules, design systems):** use `frontend-design`.
- **Accessibility audits:** use `accessibility-ux-reviewer`.
- **Shadow DOM and web components:** use `shadow-dom-component-authoring`.
- **React/Vue/Svelte component code:** use `react-nextjs` or `frontend-design`.

## Authoritative sources

| Source | Role |
|--------|------|
| [MDN HTML](https://developer.mozilla.org/en-US/docs/Web/HTML) | Primary authoring reference |
| [MDN CSS](https://developer.mozilla.org/en-US/docs/Web/CSS) | Primary authoring reference |
| [WHATWG HTML Living Standard](https://html.spec.whatwg.org/multipage/) | Canonical HTML spec (last updated May 2026) |
| [CSSWG drafts](https://drafts.csswg.org/) | Canonical CSS module specs |
| [MDN Specificity](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_cascade/Specificity) | Specificity algorithm reference |
| [MDN CSS layout](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/CSS_layout) | Layout system guide |

MDN is the faster author reference. Consult WHATWG/CSSWG when behavior is new, evolving,
or version-sensitive.

## Quick rules

1. Prefer **semantic HTML elements** over generic containers when the content has structure or meaning — `header`, `nav`, `main`, `article`, `section`, `aside`, `footer`.
2. Use **HTML for structure/meaning** and **CSS for presentation**. Never encode layout in markup.
3. Default to **lowercase tags** — MDN recommends this as the convention.
4. Keep **CSS selector weight low and intentional** — specificity is an algorithm, not intuition.
5. Choose **layout tools by job**: normal flow first → flexbox (one axis) → grid (two axes) → float (text wrap only) → positioning (overlay/offset).
6. Build **responsive layouts** with container queries and fluid techniques rather than fixed viewport assumptions.
7. Treat CSSWG/WHATWG docs as canonical for behavior; use MDN first for authoring decisions.

## HTML reference

### Document structure

Every HTML document needs this skeleton:

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Page Title</title>
    <link rel="stylesheet" href="styles.css">
  </head>
  <body>
    <header>...</header>
    <main>...</main>
    <footer>...</footer>
  </body>
</html>
```

- `html` — root; all elements descend from it.
- `head` — machine-readable metadata; not displayed as body content.
- `body` — visible document content; only one per document.

### Semantic sectioning elements

| Element | Landmark role | Use for |
|---------|--------------|---------|
| `<header>` | `banner` (top-level) | Site or section header, logo, nav |
| `<nav>` | `navigation` | Navigation links |
| `<main>` | `main` | Primary content; one per page |
| `<article>` | — | Self-contained, independently distributable content |
| `<section>` | `region` (if named) | Thematic grouping within a document |
| `<aside>` | `complementary` | Tangentially related content |
| `<footer>` | `contentinfo` (top-level) | Authorship, copyright, related links |

### Metadata elements

| Element | Purpose | Key attributes |
|---------|---------|---------------|
| `<meta>` | Encoding, viewport, descriptive metadata | `charset`, `name`/`content`, `http-equiv` |
| `<link>` | External resource relationship | `rel`, `href` |
| `<title>` | Tab/title bar; used by AT and search | Text only — HTML tags rendered as plain text |
| `<base>` | Base URL for relative links | `href`; only one allowed per document |
| `<style>` | Inline CSS | Prefer linked stylesheets for document-wide styles |

### Key interactive and media elements

| Element | Use for | Caveats |
|---------|---------|---------|
| `<button>` | Actions and commands | Built-in keyboard behavior; never use `<div>` as a button |
| `<a href>` | Navigation to a URL | `href` required for keyboard accessibility |
| `<form>`, `<input>`, `<select>`, `<textarea>` | Data entry | Always associate `<label for="id">` |
| `<dialog>` | Modal and non-modal dialogs | Native focus trapping; use `showModal()` for modals |
| `<details>`/`<summary>` | Disclosure widget | Native toggle; no JS needed |
| `<img>` | Images | Always provide `alt`; set `width`/`height` to prevent CLS |
| `<video>`, `<audio>` | Media | Provide captions/transcripts for accessibility |

### Modern HTML features (2025–2026)

- **`popover` attribute + `popovertarget`** — native popover behavior with light-dismiss, no JS required.
- **`<dialog>`** — native modal with built-in focus trap and `::backdrop` pseudo-element.
- **Customizable `<select>`** — `appearance: base-select` in CSS unlocks full custom styling.
- **`inert` attribute** — makes an element and all descendants non-interactive and invisible to AT.

## CSS reference

### Cascade, inheritance, and specificity

Specificity is calculated as a three-column value: **ID — CLASS — TYPE**.

| Selector type | Column affected | Example |
|--------------|----------------|---------|
| ID selector | ID (1-0-0) | `#nav` |
| Class, attribute, pseudo-class | CLASS (0-1-0) | `.active`, `[type="text"]`, `:hover` |
| Type selector, pseudo-element | TYPE (0-0-1) | `p`, `::before` |
| Universal `*`, `:where()`, combinators | None (0-0-0) | — |

Rules: higher specificity wins. On tie, later declaration wins. `!important` overrides all
specificity — use sparingly. `:is()` takes the specificity of its most specific argument;
`:where()` always contributes zero.

### Box model

```css
/* border-box is the recommended default */
*, *::before, *::after { box-sizing: border-box; }
```

- `content-box` (default): `width`/`height` excludes padding and border.
- `border-box`: `width`/`height` includes padding and border.

### Layout systems

**Normal flow** — default. Block elements stack vertically; inline elements flow horizontally. Understand this before reaching for other systems.

**Flexbox** — one-dimensional (row or column).

```css
.container {
  display: flex;
  gap: 16px;
  align-items: center;     /* cross axis */
  justify-content: space-between; /* main axis */
  flex-wrap: wrap;
}
.item { flex: 1 1 200px; } /* grow shrink basis */
```

**Grid** — two-dimensional (rows and columns).

```css
.container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 24px;
}
/* Named areas */
.layout {
  display: grid;
  grid-template-areas:
    "header header"
    "sidebar main"
    "footer footer";
}
```

**Float** — use only for wrapping text around content (original purpose). Do not use for primary page layout.

**Positioning** — for overlays and offsets from normal flow.

```css
.tooltip { position: absolute; top: 100%; left: 0; }
.sticky-nav { position: sticky; top: 0; }
.overlay { position: fixed; inset: 0; }
```

### Selectors

```css
/* Descendant */      .parent .child { }
/* Direct child */    .parent > .child { }
/* Adjacent */        h2 + p { }
/* General sibling */ h2 ~ p { }
/* Attribute */       input[type="email"] { }
/* :is() — any match */ :is(h1, h2, h3) { }
/* :has() — parent */ .card:has(img) { }   /* parent selector */
/* :not() */          li:not(.active) { }
```

### Responsive design

**Media queries** — respond to viewport or device characteristics.

```css
@media (min-width: 768px) { .sidebar { display: block; } }
@media (prefers-color-scheme: dark) { :root { --bg: #0f172a; } }
@media (prefers-reduced-motion: reduce) { * { animation: none !important; } }
```

**Container queries** — respond to the parent container's size. 93%+ browser support.

```css
.card-wrapper { container-type: inline-size; }

@container (min-width: 400px) {
  .card { flex-direction: row; }
}
```

**Fluid typography** — smooth scaling between min/max without breakpoint jumps.

```css
h1 { font-size: clamp(1.75rem, 1rem + 2.5vw, 3.5rem); }
```

### Modern CSS features (2025–2026)

| Feature | What it does | Support |
|---------|-------------|---------|
| `@layer` | Explicit cascade ordering; later layers win regardless of specificity | 96%+ |
| CSS nesting | Native `&` nesting without Sass/PostCSS | 93%+ |
| `@scope` | Limit selectors to a DOM subtree without Shadow DOM | 87%+ |
| `color-mix()` | Blend two colors in a given colorspace | 92%+ |
| `@property` | Register custom properties with type, initial value, inheritance | 90%+ |
| `:has()` | Parent/conditional selector — style based on children | 95%+ |
| Logical properties | `margin-inline`, `padding-block` — writing-mode aware | 95%+ |

### CSS cascade layers

Use `@layer` to define explicit priority instead of fighting specificity.

```css
@layer reset, base, tokens, components, utilities;

@layer components {
  .card { padding: var(--space-4); }
}
@layer utilities {
  .p-0 { padding: 0; } /* wins over components regardless of specificity */
}
/* Styles outside any @layer always win */
```

## Coding standards

### Semantic HTML

- Use elements for their **meaning**, not their default appearance.
- Never use `<div>` or `<span>` when a semantic element exists for the job.
- One `<h1>` per page; heading levels must be sequential (no skipping h2→h4).

### Forms

- Every `<input>` needs an associated `<label for="id">` — not just a placeholder.
- Use built-in HTML form elements (`<datalist>`, `<output>`, `<progress>`) when they match the interaction.
- Do not style away the semantic structure of form controls.

### CSS organization

- Keep specificity low; prefer class selectors over ID selectors for styling.
- Use `@layer` to manage cascade priority in large projects.
- Reach for flexbox first for component-level alignment; grid for page/section layout.

### Accessibility-minded markup

- Semantic HTML directly benefits screen readers, keyboard users, and search engines.
- Structure content with real headings, lists, and landmark elements — not just visually.

## Layout decision guide

| Problem | Tool |
|---------|------|
| Stack elements vertically | Normal flow (block) |
| Center an item in a container | Flexbox `align-items: center; justify-content: center` |
| Equal-width columns | Grid `repeat(N, 1fr)` |
| Item overlaps another | `position: absolute` on child, `position: relative` on parent |
| Element sticks while scrolling | `position: sticky` |
| Text wraps around an image | `float: left/right` |
| Component adapts to container width | Container queries |

## Known spec notes

- **HTML source of truth:** WHATWG HTML Living Standard (versioned HTML5 is a W3C artifact; WHATWG does not version HTML).
- **CSS source of truth:** CSSWG module drafts — CSS has no single spec; each feature lives in its own module.
- MDN is the faster author reference; WHATWG/CSSWG for behavior edge cases.
