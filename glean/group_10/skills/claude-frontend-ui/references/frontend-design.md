<!-- hub-reference-banner -->
> **Reference file — part of the `frontend-ui` hub.** Formerly the standalone `frontend-design` skill.
> Sibling topics in this family are now reference files under the hubs (`frontend-ui`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: frontend-design
version: "1.1.0"
updated: "2026-05-29"
description: >
  Production-grade frontend engineering: design systems and tokens, component architecture
  (React compound components, custom hooks, composition), CSS architecture (Tailwind, CSS
  Modules, @layer, modern CSS), responsive design (container queries, fluid typography),
  accessibility (WCAG 2.2, ARIA, focus management), animation (scroll-driven, View
  Transitions), state management (Zustand, Jotai, TanStack Query), Core Web Vitals
  optimization, Storybook testing, and micro-frontends. TRIGGER: user asks to build a web
  component, page, or app; asks for frontend architecture, CSS strategy, or component design
  review. SKIP: visual design ideation only (use web-design), UI/UX palette and style catalog
  (use ui-ux-pro-max), plain HTML/CSS reference (use html-css), accessibility audit checklist
  (use accessibility-ux-reviewer), React 19/Next.js 15 App Router specifics (use react-nextjs).
license: Complete terms in LICENSE.txt
whenToUse:
  - "Build a new page or component"
  - "Choose a CSS architecture strategy"
  - "Make a component responsive with container queries"
  - "Add scroll-driven animations or View Transitions"
  - "Pass an accessibility audit"
  - "Choose a state management library"
  - "Fix slow Core Web Vitals (LCP, INP, CLS)"
  - "Set up Storybook component testing"
  - "Architect a micro-frontend with Module Federation"
  - "Design system tokens and Atomic Design structure"
related_skills:
  - react-nextjs
  - html-css
  - accessibility-ux-reviewer
  - ui-ux-pro-max
  - shadow-dom-component-authoring
---

## Overview

This skill guides creation of distinctive, production-grade frontend interfaces. It covers the
full frontend engineering surface: design systems, component architecture (React-focused,
framework-agnostic principles), CSS architecture (Tailwind, CSS Modules, Cascade Layers, modern
CSS), responsive design (container queries, fluid typography), accessibility (WCAG 2.2),
animation (scroll-driven, View Transitions), state management (Zustand, Jotai, TanStack Query),
performance (Core Web Vitals), testing (Storybook, visual regression), and micro-frontends.

## When NOT to use this skill

| Goal | Use instead |
|------|-------------|
| Visual/graphic design ideation without code | `web-design` |
| UI/UX style catalog, palette generation, design-system search | `ui-ux-pro-max` |
| HTML/CSS element and property reference | `html-css` |
| Accessibility audit with WCAG criterion citations | `accessibility-ux-reviewer` |
| React 19 / Next.js 15 App Router, Server Components, server actions | `react-nextjs` |
| Shadow DOM web components and Chrome extension overlays | `shadow-dom-component-authoring` |

## Quick decision guide

| Scenario | Start here |
|----------|------------|
| Build a new page or component | Design Thinking → Component Architecture |
| Choose a CSS strategy | CSS Architecture (Tailwind vs. Modules vs. @layer) |
| Make a component responsive | Responsive Design (container queries, clamp) |
| Add animations | Animation & Motion (scroll-driven, View Transitions) |
| Pass an accessibility audit | Accessibility (WCAG 2.2 checklist) |
| Choose a state management library | State Management (category table) |
| Fix slow Core Web Vitals | Performance (LCP/INP/CLS optimization) |
| Set up component testing | Testing UI (Storybook, visual regression) |
| Multi-team frontend at scale | Micro-Frontends (Module Federation) |

## Design thinking

Before writing code, commit to a direction:

- **Purpose**: What problem does this interface solve? Who uses it?
- **Tone**: Choose a clear aesthetic — brutally minimal, maximalist, retro-futuristic, organic, luxury, playful, editorial, brutalist, art deco, soft/pastel, industrial. Execute it with precision.
- **Constraints**: Framework, performance budget, accessibility requirements.
- **Differentiation**: What makes this memorable? One strong design choice executed well beats ten timid ones.

**CRITICAL**: Bold maximalism and refined minimalism both work. Intentionality matters more than intensity.

Deliver working code (HTML/CSS/JS, React, Vue, etc.) that is production-grade, visually striking, and cohesive. Meticulously refined in every detail.

## Frontend aesthetics guidelines

- **Typography**: Choose distinctive, characterful fonts. Avoid generic families (Arial, Inter, Roboto, system fonts). Pair a display font with a refined body font.
- **Color**: Commit to a cohesive system via CSS variables. Dominant colors with sharp accents outperform evenly-distributed timid palettes.
- **Motion**: Use CSS animations for micro-interactions. For React, use Motion library. One well-orchestrated page-load stagger creates more impact than scattered micro-interactions.
- **Spatial composition**: Unexpected layouts, asymmetry, overlap, generous negative space or controlled density.
- **Backgrounds**: Gradient meshes, noise textures, geometric patterns, layered transparencies, dramatic shadows — create atmosphere, not just solid fills.

NEVER use generic AI-generated aesthetics: overused font families (Inter, Roboto, Space Grotesk), clichéd purple gradients, predictable layouts, cookie-cutter patterns.

**Match complexity to vision.** Maximalist designs need extensive animations and effects. Minimalist designs need restraint, precision, and careful spacing. Elegance comes from executing the vision well.

## Design systems

### Design tokens

Organize in three tiers:
- **Primitive tokens**: raw values — `blue-500: #3B82F6`, `space-4: 16px`
- **Semantic tokens**: contextual meaning — `color-primary: {blue-500}`, `spacing-element: {space-4}`
- **Component tokens**: scoped — `button-bg: {color-primary}`, `card-padding: {spacing-element}`

W3C DTCG Spec 2025.10 provides a vendor-neutral JSON format for sharing tokens across tools. Style Dictionary v4+ has first-class DTCG support. CSS custom properties (`--color-blue-500`) enable runtime theming without rebuilds.

### Atomic Design

Atoms → Molecules → Organisms → Templates → Pages. Document every level in Storybook.

## Component architecture

### React patterns

**Compound Components** — two or more components sharing implicit state via React Context. Best for related UI groups like `<Select>`, `<Tabs>`, `<Accordion>`.

```tsx
const TabsContext = createContext<TabsState | null>(null);

export function Tabs({ children, defaultTab }: TabsProps) {
  const [active, setActive] = useState(defaultTab);
  return (
    <TabsContext value={{ active, setActive }}>
      <div>{children}</div>
    </TabsContext>
  );
}

export function Tab({ value, children }: TabProps) {
  const { active, setActive } = use(TabsContext)!;
  return (
    <button
      role="tab"
      aria-selected={active === value}
      onClick={() => setActive(value)}
    >
      {children}
    </button>
  );
}
```

**Custom Hooks** — extract stateful logic into reusable functions (`useDebounce`, `useLocalStorage`, `useMediaQuery`). Replaced HOCs and render props for 90% of use cases.

**Composition over configuration** — break components into smaller pieces consumers assemble. Prefer `children` and slot patterns over deeply nested prop APIs.

**Presentational vs. Container** — in 2026 this maps to Server Components (data) and Client Components (interactivity).

### Framework-agnostic principles

- Single Responsibility: each component does one thing well
- Open/Closed: extend via composition, not modification
- Dependency Inversion: depend on abstractions (interfaces, context)

## CSS architecture

### Tailwind CSS (dominant in 2026)

Zero-runtime, React Server Components compatible, AI tools produce consistent output. Tailwind v4 has native CSS-first config with no PostCSS dependency.

**Pragmatic stack**: Tailwind for 90% of styles, CSS Modules for complex selectors or keyframe-heavy animation. They coexist without conflict.

### CSS Modules

Scope class names to the component file. Zero-runtime like Tailwind. Use for component-specific styles unwieldy as utilities.

### CSS Cascade Layers (`@layer`)

Explicit priority ordering independent of selector complexity. Styles in later layers always beat earlier layers. 96%+ browser support.

```css
@layer reset, base, tokens, components, utilities;

@layer components { .card { padding: var(--space-4); } }
@layer utilities  { .p-0  { padding: 0; } /* beats components */ }
/* Unlayered CSS always wins */
```

### Modern CSS (2026)

| Feature | What it does |
|---------|-------------|
| `@scope` | Style boundaries without Shadow DOM |
| CSS Nesting | Native `&` nesting, no preprocessor |
| `color-mix()` | Blend colors in CSS |
| `@property` | Register custom properties with type + animation support |
| Logical properties | `margin-inline`, `padding-block` for writing-mode-aware layouts |
| `:has()` | Parent/conditional selector |

## Responsive design

### Container queries (93%+ support)

Style based on container size rather than viewport. Enables truly reusable components.

```css
.card-container { container-type: inline-size; }

@container (min-width: 400px) { .card { flex-direction: row; } }
@container (max-width: 399px) { .card { flex-direction: column; } }
```

### Fluid typography with `clamp()`

```css
h1   { font-size: clamp(1.75rem, 1rem + 2.5vw, 3.5rem); }
body { font-size: clamp(1rem, 0.875rem + 0.5vw, 1.25rem); }
```

Use `rem` for boundaries, `vw` for the scaling factor.

### Responsive images

```html
<picture>
  <source type="image/avif" srcset="hero.avif">
  <source type="image/webp" srcset="hero.webp">
  <img src="hero.jpg" alt="Hero" width="1200" height="600"
       loading="eager" fetchpriority="high">
</picture>
```

## Accessibility (a11y)

Accessibility is a legal requirement. The European Accessibility Act is enforced; ADA deadline was April 2026. 96.3% of the top 1M websites still fail basic accessibility tests.

### WCAG 2.2 POUR principles

1. **Perceivable** — alt text, captions, 4.5:1 contrast for normal text, 3:1 for large text
2. **Operable** — keyboard nav, 24×24px minimum touch targets (WCAG 2.5.8)
3. **Understandable** — predictable behavior, clear labels
4. **Robust** — interpretable by assistive technologies

### Semantic HTML first

```html
<!-- Use <button> not <div onClick> -->
<button type="button">Save</button>

<!-- Use <a href> not <span onClick> -->
<a href="/dashboard">Dashboard</a>

<!-- Heading hierarchy: one h1, sequential levels -->
<h1>Page Title</h1>
<h2>Section</h2>
<h3>Subsection</h3>
```

### ARIA — use only when native HTML can't express the relationship

```tsx
// Dialog with proper ARIA
<div role="dialog" aria-modal="true" aria-labelledby="dialog-title">
  <h2 id="dialog-title">Confirm Delete</h2>
  ...
</div>

// Live region for dynamic updates
<div aria-live="polite" aria-atomic="true">
  {statusMessage}
</div>
```

### Focus management

```css
/* Never outline: none without replacement */
:focus-visible {
  outline: 2px solid var(--color-focus);
  outline-offset: 2px;
}
```

Focus trapping in modals: Tab cycle stays within dialog until dismissed. Skip navigation: first focusable element, visible on focus, jumps to `<main>`.

### Pre-ship accessibility checklist

- [ ] All images have descriptive `alt` text (or `alt=""` for decorative)
- [ ] Color contrast: 4.5:1 normal text, 3:1 large text (WCAG AA)
- [ ] All interactive elements keyboard-reachable and operable
- [ ] Focus indicators visible with 3:1 contrast ratio
- [ ] Touch targets ≥ 24×24px
- [ ] One `<h1>` per page, headings sequential
- [ ] Forms have associated `<label>` elements
- [ ] Dynamic updates use `aria-live` regions
- [ ] Modals trap focus and return focus on close
- [ ] `prefers-reduced-motion` respected for all animations
- [ ] axe-core returns zero critical/serious violations
- [ ] Keyboard-only navigation completes all primary flows

## Animation & motion

### Scroll-driven animations (CSS native, universal support in 2026)

Runs on compositor thread — 60fps guaranteed even under main-thread load.

```css
.progress-bar {
  animation: grow-width linear;
  animation-timeline: scroll();
}

.reveal-card {
  animation: fade-slide-up linear both;
  animation-timeline: view();
  animation-range: entry 0% entry 100%;
}
```

### View Transitions API

```js
document.startViewTransition(() => {
  updateDOM(); // Your DOM mutation
});
```

### Animation best practices

| Rule | Detail |
|------|--------|
| Animate only `transform` and `opacity` | Avoids layout triggers |
| Respect `prefers-reduced-motion` | Always provide instant-state alternative |
| Duration: micro-interactions 150–300ms | Page transitions 300–500ms |
| Use `cubic-bezier()` for natural easing | Avoid linear for UI transitions |
| Motion should communicate meaning | State change, hierarchy, spatial relationship — not decoration |

## State management (2026)

| Category | Solution | Size |
|----------|----------|------|
| Server state | TanStack Query or Server Component fetch | — |
| Client state (stores) | Zustand | ~3KB |
| Client state (atoms) | Jotai | ~4KB |
| Form state | React Hook Form + Zod | — |
| URL state | Router params + `useSearchParams` | — |
| Derived state | Selectors, `useMemo`, computed atoms | — |

**Decision rule:**
1. Is it data from an API/DB? → TanStack Query or Server Component
2. Is it form state? → React Hook Form + `useActionState`
3. Is it local to one component? → `useState`/`useReducer`
4. Is it many independent atoms? → Jotai
5. Is it one coherent store? → Zustand
6. Is it enterprise-scale with strict patterns? → Redux Toolkit

Never store server-fetched data in Zustand/Jotai — use a dedicated server state solution.

## Performance & Core Web Vitals

### Targets (2026)

| Metric | Good | Poor |
|--------|------|------|
| LCP | ≤ 2.5s | > 4.0s |
| INP | ≤ 200ms | > 500ms |
| CLS | ≤ 0.1 | > 0.25 |

43% of sites fail the 200ms INP threshold — the most commonly failed metric.

### LCP optimization

```html
<!-- Preload LCP image -->
<link rel="preload" as="image" href="hero.webp">
```

- Inline critical CSS in `<head>` (first 14KB)
- `fetchpriority="high"` on LCP element
- Server-render above-the-fold content

### INP optimization

- Break long tasks (>50ms) with `scheduler.yield()` or `requestIdleCallback`
- `content-visibility: auto` for off-screen sections
- Virtualize long lists (TanStack Virtual)
- Debounce rapid input handlers

### CLS optimization

- Set explicit `width` and `height` on all images/videos/iframes
- `aspect-ratio` to reserve space for dynamic content
- `font-display: optional` or `font-display: swap` with `size-adjust`

## Testing UI

### Component testing stack

| Layer | Tool | What it catches |
|-------|------|-----------------|
| Unit | Vitest + Testing Library | Logic, state, rendering |
| Interaction | Storybook play functions | User flows within components |
| Visual | Chromatic / Playwright screenshots | Pixel-level regressions |
| Accessibility | axe-core + Storybook a11y addon | WCAG violations |
| E2E | Playwright / Cypress | Cross-page user journeys |

### Testing best practices

- Test behavior, not implementation. Query by role, label, or text — not class names.
- Use `userEvent` over `fireEvent` for realistic interaction simulation.
- Mock network requests at the boundary (MSW) rather than mocking individual functions.

## Micro-frontends

Use when multiple teams own features with independent release cadences. Do not use for single-team products.

**Module Federation 3.0 (2026):** runtime code sharing, Server-Side Module Federation for stitching Remote Server Components, TypeScript type safety across federated boundaries.

**Security**: Use `Cross-Origin-Opener-Policy` (COOP) and `Cross-Origin-Embedder-Policy` (COEP) to sandbox micro-frontends from each other.

## Anti-patterns

### Component

| Anti-pattern | Fix |
|---|---|
| Prop drilling through 5+ levels | Context, state management, or composition |
| God components doing everything | Focused sub-components via composition |
| Boolean prop explosion `<Button primary small loading>` | Variant enums + compound components |
| State hoisted too high | Colocate state with the owning component |

### CSS

| Anti-pattern | Fix |
|---|---|
| Specificity wars and `!important` escalation | `@layer` cascade architecture |
| No naming convention or scoping strategy | CSS Modules or Tailwind |
| Unused CSS accumulation | PurgeCSS or Tailwind tree-shaking |

### Performance

| Anti-pattern | Fix |
|---|---|
| Sequential data fetching waterfalls | `Promise.all()` or framework loaders |
| Unoptimized images | `srcset`, AVIF/WebP, responsive sizing |
| Synchronous third-party scripts | `async`/`defer` or post-interaction load |

## References

- [W3C Design Tokens Community Group Spec 2025.10](https://www.designtokens.org/)
- [CSS Cascade Layers Guide — CSS-Tricks](https://css-tricks.com/css-cascade-layers/)
- [Container Queries Guide 2026](https://csscodelab.com/the-ultimate-guide-to-css-container-queries-in-2026/)
- [Scroll-Driven Animations — Josh W. Comeau](https://www.joshwcomeau.com/animation/scroll-driven-animations/)
- [View Transition API — MDN](https://developer.mozilla.org/en-US/docs/Web/API/View_Transition_API)
- [Core Web Vitals 2026 — Digital Applied](https://www.digitalapplied.com/blog/core-web-vitals-2026-inp-lcp-cls-optimization-guide)
- [State Management 2026 — Dev.to](https://dev.to/jsgurujobs/state-management-in-2026-zustand-vs-jotai-vs-redux-toolkit-vs-signals-2gge)
- [React Design Patterns 2026](https://www.turbodocx.com/blog/react-design-patterns)
- [Micro-Frontends & Module Federation 2026](https://kawaldeepsingh.medium.com/microfrontends-module-federation-in-2026-a-practical-playbook-for-frontend-teams-4445c93fe61f)
- [WCAG 2 Overview — W3C WAI](https://www.w3.org/WAI/standards-guidelines/wcag/)
- [Storybook Visual Testing](https://storybook.js.org/docs/writing-tests/visual-testing)
- [Figma Design System Guide 2026](https://muz.li/blog/how-to-build-a-design-system-in-figma-a-practical-guide-2026/)
