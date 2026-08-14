<!-- hub-reference-banner -->
> **Reference file — part of the `frontend-ui` hub.** Formerly the standalone `ui-ux-pro-max` skill.
> Sibling topics in this family are now reference files under the hubs (`frontend-ui`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: ui-ux-pro-max
version: "1.1.1"
updated: "2026-06-16"
description: >
  UI/UX design intelligence for web and mobile: 50+ styles, 161 color palettes, 57 font
  pairings, 161 product types, 99 UX guidelines, 25 chart types across React, Next.js, Vue,
  Svelte, SwiftUI, React Native, Flutter, Tailwind, shadcn/ui, and HTML/CSS stacks. TRIGGER:
  user asks for design system recommendations, color palette selection, font pairings, UI style
  choices, UX pattern guidance, chart type selection, or pre-launch UI quality review for any
  web or mobile product. SKIP: building a specific React/Next.js feature (use react-nextjs),
  accessibility audits with WCAG citations (use accessibility-ux-reviewer), iOS SwiftUI
  implementation (use mobile-ios-design), CSS/HTML reference (use html-css).
whenToUse:
  - "Designing a new page — landing page, dashboard, admin panel, SaaS, mobile app"
  - "Creating or refactoring UI components — buttons, modals, forms, tables, charts"
  - "Choosing color schemes, typography systems, or layout systems"
  - "Reviewing UI code for UX quality or visual consistency"
  - "Implementing navigation structures, animations, or responsive behavior"
  - "Making product-level design decisions — style, information hierarchy, brand"
  - "Improving perceived quality, clarity, or usability of an interface"
  - "Aligning cross-platform design — web, iOS, Android"
  - "Building design systems or reusable component libraries"
  - "Recommend color palette for a fintech / healthcare / SaaS product"
  - "What font pairing fits a luxury brand?"
  - "Which chart type should I use for time-series data?"
related_skills:
  - frontend-ui
  - accessibility-ux-reviewer
---

# UI/UX Pro Max — Design Intelligence

Comprehensive design guide for web and mobile applications. Contains 50+ styles, 161 color
palettes, 57 font pairings, 161 product types with reasoning rules, 99 UX guidelines, and 25
chart types across 10 technology stacks.

## When NOT to use this skill

- **Building a specific React/Next.js feature:** use `react-nextjs`.
- **Accessibility audits with WCAG criterion citations:** use `accessibility-ux-reviewer`.
- **iOS SwiftUI implementation patterns:** use `mobile-ios-design`.
- **CSS/HTML element reference:** use `html-css`.
- **Pure frontend component architecture:** use `frontend-design`.
- **Pure backend, API, database, infrastructure work:** not applicable.

**Decision rule:** if the task changes how a feature looks, feels, moves, or is interacted with, use this skill.

## Rule categories by priority

| Priority | Category | Impact | Key checks | Anti-patterns |
|----------|----------|--------|------------|---------------|
| 1 | Accessibility | CRITICAL | Contrast 4.5:1, alt text, keyboard nav, aria-labels | Removing focus rings, icon-only buttons without labels |
| 2 | Touch & Interaction | CRITICAL | Min 44×44px targets, 8px+ spacing, loading feedback | Hover-only interactions, instant state changes (0ms) |
| 3 | Performance | HIGH | WebP/AVIF, lazy loading, reserve space (CLS < 0.1) | Layout thrashing, cumulative layout shift |
| 4 | Style Selection | HIGH | Match product type, consistency, SVG icons (no emoji) | Mixing flat and skeuomorphic, emoji as icons |
| 5 | Layout & Responsive | HIGH | Mobile-first, no horizontal scroll, viewport meta | Fixed px containers, disable zoom |
| 6 | Typography & Color | MEDIUM | Base 16px, line-height 1.5, semantic color tokens | Body text < 12px, gray-on-gray, raw hex in components |
| 7 | Animation | MEDIUM | 150–300ms, motion conveys meaning, spatial continuity | Decorative-only animation, animating width/height |
| 8 | Forms & Feedback | MEDIUM | Visible labels, error near field, progressive disclosure | Placeholder-only label, errors only at top |
| 9 | Navigation Patterns | HIGH | Predictable back, bottom nav ≤ 5, deep linking | Overloaded nav, broken back behavior |
| 10 | Charts & Data | LOW | Legends, tooltips, accessible color palettes | Color as the only data differentiator |

---

## Quick reference

### 1. Accessibility (CRITICAL)

- `color-contrast` — 4.5:1 for normal text, 3:1 for large text
- `focus-states` — visible focus rings, 2–4px, on all interactive elements
- `alt-text` — descriptive alt text for meaningful images
- `aria-labels` — `aria-label` for icon-only buttons
- `keyboard-nav` — tab order matches visual order; full keyboard support
- `form-labels` — `<label for="id">` on every input
- `skip-links` — skip to main content for keyboard users
- `heading-hierarchy` — sequential h1→h6, no level skip
- `color-not-only` — add icon or text alongside color to convey meaning
- `dynamic-type` — support system text scaling; avoid truncation as text grows
- `reduced-motion` — respect `prefers-reduced-motion`
- `escape-routes` — cancel/back in modals and multi-step flows

### 2. Touch & Interaction (CRITICAL)

- `touch-target-size` — min 44×44pt (iOS) / 48×48dp (Android); extend hit area beyond visual
- `touch-spacing` — min 8px/8dp gap between touch targets
- `hover-vs-tap` — primary interactions on tap/click, not hover
- `loading-buttons` — disable during async; show spinner or progress
- `tap-delay` — `touch-action: manipulation` to remove 300ms delay
- `gesture-conflicts` — avoid horizontal swipe on main content; prefer vertical scroll
- `press-feedback` — visual feedback within 80–150ms
- `haptic-feedback` — confirmations and important actions; avoid overuse
- `safe-area-awareness` — keep targets away from notch, Dynamic Island, gesture bar

### 3. Performance (HIGH)

- `image-optimization` — WebP/AVIF, `srcset`/`sizes`, lazy load non-critical assets
- `image-dimension` — explicit `width`/`height` or `aspect-ratio` to prevent CLS
- `font-loading` — `font-display: swap`/`optional`; preload only critical fonts
- `lazy-loading` — lazy load non-hero components with dynamic import / route splitting
- `content-jumping` — reserve space for async content (CLS < 0.1)
- `virtualize-lists` — virtualize lists of 50+ items
- `debounce-throttle` — debounce scroll, resize, and rapid input handlers

### 4. Style Selection (HIGH)

- `style-match` — match style to product type and industry
- `consistency` — same style across all pages; no style mixing
- `no-emoji-icons` — use SVG icons (Heroicons, Lucide); never emoji as UI icons
- `effects-match-style` — shadows, blur, radius aligned with chosen style
- `dark-mode-pairing` — design light/dark variants together
- `icon-style-consistent` — one icon set, consistent stroke weight and corner radius
- `primary-action` — each screen has one primary CTA; secondary actions subordinate

### 5. Layout & Responsive (HIGH)

- `viewport-meta` — `width=device-width, initial-scale=1`; never disable zoom
- `mobile-first` — design mobile-first, then scale up
- `readable-font-size` — min 16px body on mobile (avoids iOS auto-zoom)
- `line-length-control` — 35–60 chars/line mobile; 60–75 desktop
- `spacing-scale` — 4pt/8dp incremental system
- `container-width` — consistent max-width on desktop (max-w-6xl / 7xl)
- `viewport-units` — prefer `min-h-dvh` over `100vh` on mobile

### 6. Typography & Color (MEDIUM)

- `line-height` — 1.5–1.75 for body text
- `font-scale` — consistent type scale (e.g. 12 14 16 18 24 32)
- `weight-hierarchy` — bold headings (600–700), regular body (400), medium labels (500)
- `color-semantic` — semantic color tokens (`primary`, `error`, `surface`); no raw hex in components
- `color-dark-mode` — desaturated/lighter tonal variants in dark mode; test contrast separately
- `truncation-strategy` — prefer wrapping; when truncating use ellipsis + tooltip/expand
- `number-tabular` — tabular figures for data columns, prices, timers

### 7. Animation (MEDIUM)

- `duration-timing` — 150–300ms micro-interactions; complex transitions ≤ 400ms
- `transform-performance` — animate `transform`/`opacity` only; never `width`/`height`/`top`/`left`
- `easing` — ease-out for entering, ease-in for exiting
- `motion-meaning` — every animation expresses cause-effect; no decoration-only motion
- `exit-faster-than-enter` — exit ~60–70% of enter duration
- `interruptible` — animations must be interruptible; user tap cancels in-progress animation
- `no-blocking-animation` — never block user input during animation

### 8. Forms & Feedback (MEDIUM)

- `input-labels` — visible label per input (never placeholder-only)
- `error-placement` — errors below the related field, not only at page top
- `inline-validation` — validate on blur, not keystroke
- `progressive-disclosure` — reveal complex options progressively
- `input-type-keyboard` — semantic input types (`email`, `tel`, `number`) for correct mobile keyboard
- `undo-support` — "Undo delete" toast for destructive actions
- `error-clarity` — error messages state cause and recovery path, not just "Invalid input"
- `multi-step-progress` — show step indicator; allow back navigation

### 9. Navigation Patterns (HIGH)

- `bottom-nav-limit` — max 5 items; icon + text label
- `back-behavior` — predictable, consistent; preserves scroll and state
- `deep-linking` — all key screens reachable via URL/deep link
- `modal-escape` — clear close affordance; swipe-down to dismiss on mobile
- `state-preservation` — navigating back restores scroll, filters, and input
- `adaptive-navigation` — sidebar on ≥ 1024px; bottom/top nav on smaller screens
- `focus-on-route-change` — move focus to main content region after page transition

### 10. Charts & Data (LOW)

- `chart-type` — trend → line, comparison → bar, proportion → pie/donut (≤ 5 categories)
- `color-guidance` — accessible palettes; supplement color with patterns/shapes for colorblind users
- `legend-visible` — always show legend; position near chart
- `tooltip-on-interact` — hover (web) or tap (mobile) shows exact values
- `responsive-chart` — reflow or simplify on small screens
- `empty-data-state` — meaningful empty state with guidance, not a blank chart
- `screen-reader-summary` — `aria-label` describing key insight for screen readers

---

## How to use the search scripts

This skill includes Python search scripts for generating design systems from the bundled
databases. The scripts are optional — use the quick reference above when scripts are
unavailable (no Python, CI environment, etc.).

**Prerequisites:** Python 3.8+ installed (`python3 --version`).

**Step 1 — Generate a design system (recommended starting point):**

```bash
python3 skills/ui-ux-pro-max/scripts/search.py "<product_type> <industry> <keywords>" \
  --design-system [-p "Project Name"]
```

Returns: style, color palette, typography, effects, and anti-patterns for the product context.

**Step 2 — Deep-dive a specific domain:**

```bash
python3 skills/ui-ux-pro-max/scripts/search.py "<keyword>" --domain <domain> [-n <max_results>]
```

**Step 3 — Persist for cross-session retrieval:**

```bash
python3 skills/ui-ux-pro-max/scripts/search.py "<query>" --design-system --persist -p "Project Name"
```

Creates:
- `design-system/MASTER.md` — global source of truth
- `design-system/pages/<page>.md` — page-specific overrides (with `--page "<page>"`)

When building a specific page, check `design-system/pages/<page>.md` first; fall back to `MASTER.md`.

**When scripts are unavailable:** use the Quick Reference sections above directly. They cover all CRITICAL and HIGH priority rules without requiring script execution.

### Available domains

| Domain | Use for | Example keywords |
|--------|---------|-----------------|
| `product` | Product type patterns | SaaS, e-commerce, healthcare, beauty |
| `style` | UI styles and effects | glassmorphism, minimalism, dark mode, brutalism |
| `typography` | Font pairings | elegant, playful, professional, modern |
| `color` | Palettes by product/industry | fintech, healthcare, beauty, entertainment |
| `landing` | Page structure and CTA strategies | hero, testimonial, pricing, social-proof |
| `chart` | Chart types and library recommendations | trend, comparison, timeline, funnel |
| `ux` | Best practices and anti-patterns | animation, accessibility, loading, z-index |
| `google-fonts` | Individual Google Fonts lookup | sans-serif, variable, popular |
| `react` | React/Next.js performance patterns | waterfall, suspense, memo, bundle |
| `web` | App interface guidelines (iOS/Android/RN) | touch targets, safe areas, Dynamic Type |

---

## Common professional UI rules

### Icons and visual elements

| Rule | Do | Avoid |
|------|----|-------|
| No emoji as structural icons | Vector icons (Lucide, Heroicons, SF Symbols) | 🎨 🚀 ⚙️ as nav or system controls |
| Vector-only assets | SVG / platform vector icons | Raster PNG icons that blur |
| Consistent icon sizing | Design tokens: `icon-sm`, `icon-md` (24pt), `icon-lg` | Arbitrary sizes (20pt / 24pt / 28pt) mixed |
| Stroke consistency | Same stroke width within a visual layer (1.5px or 2px) | Mixed thick and thin strokes |
| Filled vs outline discipline | One icon style per hierarchy level | Mixed filled and outline at same level |
| Touch target minimum | 44×44pt interactive area; use `hitSlop` if icon is smaller | Small icons without expanded tap area |

### Light/dark mode contrast

| Rule | Do | Avoid |
|------|----|-------|
| Text contrast (light) | Body text ≥ 4.5:1 on light surfaces | Low-contrast gray text |
| Text contrast (dark) | Primary ≥ 4.5:1; secondary ≥ 3:1 on dark surfaces | Text that blends into background |
| Token-driven theming | Semantic color tokens mapped per theme | Hardcoded hex values per screen |
| State contrast parity | Hover/focus/disabled states equally clear in both themes | States defined for one theme only |
| Scrim legibility | Modal scrim 40–60% opacity | Weak scrim leaving background competing |

### Layout and spacing

| Rule | Do | Avoid |
|------|----|-------|
| Safe-area compliance | Respect top/bottom safe areas for fixed UI | Content under notch or gesture bar |
| 8dp spacing rhythm | 4/8dp spacing system for padding, gaps, sections | Random increments |
| Adaptive gutters | Increase horizontal insets at larger widths / landscape | Same narrow gutter on all devices |
| Scroll and fixed coexistence | Add content insets so lists are not hidden behind fixed bars | Scroll content obscured by sticky headers |

---

## Pre-delivery checklist

### Visual quality
- [ ] No emoji used as icons — SVG only
- [ ] All icons from a consistent family and style
- [ ] Pressed-state visuals do not shift layout or cause jitter
- [ ] Semantic theme tokens used consistently — no ad-hoc hex colors

### Interaction
- [ ] All tappable elements have clear pressed feedback
- [ ] Touch targets ≥ 44×44pt iOS / 48×48dp Android
- [ ] Micro-interaction timing 150–300ms with native-feeling easing
- [ ] Disabled states visually clear and non-interactive
- [ ] Screen reader focus order matches visual order; labels are descriptive

### Light/dark mode
- [ ] Primary text contrast ≥ 4.5:1 in both modes
- [ ] Secondary text contrast ≥ 3:1 in both modes
- [ ] Dividers/borders distinguishable in both modes
- [ ] Modal scrim 40–60% opacity in both modes
- [ ] Both themes tested — not inferred from one

### Layout
- [ ] Safe areas respected for headers, tab bars, bottom CTAs
- [ ] Scroll content not hidden behind fixed/sticky bars
- [ ] Verified on small phone, large phone, tablet (portrait + landscape)
- [ ] 4/8dp spacing rhythm maintained throughout
- [ ] Long-form text readable on large devices (no full-width paragraphs on tablets)

### Accessibility
- [ ] All meaningful images/icons have accessibility labels
- [ ] Form fields have labels, hints, and clear error messages
- [ ] Color is not the sole indicator for any information
- [ ] Reduced motion and dynamic text size supported without layout breakage
- [ ] Accessibility roles, states (selected, disabled, expanded) announced correctly

## Related references (frontend-ui hub)

When the task is not "produce/style a UI" but "**evaluate** one against a named
instrument", hand off:

- **Named usability heuristic evaluation + Laws of UX** (Nielsen's 10 heuristics
  with detection→severity→fix, the 0–4 severity scale, and the Laws of UX —
  Hick/Fitts/Jakob/Miller/Gestalt/etc.) → `usability-heuristics-laws-of-ux`.
- **Visual-design critique rubric & methodology** (hierarchy, gestalt, C.R.A.P.,
  scanning patterns; I-like/I-wish/what-if) → `visual-design-principles-and-critique`.
- **WCAG 2.2 / WAI-ARIA accessibility review** → `accessibility-ux-reviewer`.
