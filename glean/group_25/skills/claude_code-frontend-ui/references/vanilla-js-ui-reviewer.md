<!-- hub-reference-banner -->
> **Reference file — part of the `frontend-ui` hub.** Formerly the standalone `vanilla-js-ui-reviewer` skill.
> Sibling topics in this family are now reference files under the hubs (`frontend-ui`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: vanilla-js-ui-reviewer
title: Vanilla JS UI Reviewer
description: |
  Practical review reference for large plain-JavaScript UIs that manage DOM, events, focus, state, and accessibility without a framework.
  TRIGGER: reviewing vanilla JS UI code; event delegation bugs; innerHTML XSS audit; focus/keyboard accessibility in plain JS; listener leak or teardown review; ARIA state synchronization; DOM mutation performance; sidepanel or popup JS review.
  SKIP: React/Vue/Angular/Svelte component review (no framework abstractions here); backend JS with no DOM; automated accessibility scanning (use accessibility-ux-reviewer for that layer).
category: developer
version: "1.1.0"
updated: "2026-05-29"
keywords:
  - vanilla JavaScript
  - DOM
  - event delegation
  - accessibility
  - ARIA
  - focus management
  - innerHTML
  - listener teardown
  - AbortController
  - MutationObserver
  - sidepanel
  - popup
  - layout thrashing
  - keyboard navigation
when_to_use:
  - "review this vanilla JS UI"
  - "event delegation bug"
  - "innerHTML security review"
  - "focus trap or focus loss in plain JS"
  - "listener leak"
  - "ARIA state out of sync"
  - "keyboard navigation in custom widget"
  - "DOM performance in sidepanel"
  - "AbortController teardown pattern"
  - "modal focus management without a framework"
related_skills:
  - accessibility-ux-reviewer
  - chrome-extension-security-reviewer
  - platform-adapter-reviewer
origin: local
---

# Vanilla JS UI Reviewer

Practical review reference for large plain-JavaScript UIs that directly manage DOM, events, focus, state, and accessibility without a framework.

## How to use this skill

Start from the bundled context below. Defer to the cited official documentation for exact APIs and edge-case behavior. If the request falls outside vanilla-JS UI review, choose a more appropriate skill.

**Sources of truth:**
- **MDN** — DOM/event/HTML/CSS behavior
- **WAI-ARIA APG + WCAG understanding docs** — keyboard/focus/widget behavior
- **web.dev** — browser-performance guidance around DOM size and layout work

**Version note:** based on official pages accessed 2026-05-10, framed for this repo's large vanilla-JS sidepanel and popup UIs.

---

## Source scope

- **Event model, delegation, teardown, abortable listeners:** MDN `addEventListener`, `removeEventListener`, `AbortController`, `Event.target`, `Event.currentTarget`, `Element.closest`, `Event.composedPath()` ([MDN addEventListener](https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener), [MDN AbortController](https://developer.mozilla.org/en-US/docs/Web/API/AbortController))
- **Safe DOM updates and structure:** MDN `textContent`, `innerHTML`, `replaceChildren`, `appendChild`, `DocumentFragment`, and semantic elements `button`, `main`, `nav`, `section`, `ul`, `li` ([MDN textContent](https://developer.mozilla.org/en-US/docs/Web/API/Node/textContent), [MDN innerHTML](https://developer.mozilla.org/en-US/docs/Web/API/Element/innerHTML))
- **Focus, keyboard behavior, modal/widget semantics:** MDN `focus()`, `activeElement`, `tabIndex`, `:focus-visible`, `dialog`, `inert`; APG dialog/tabs/accordion/menu/treeview patterns; WCAG focus order and name/role/value ([APG keyboard interface](https://www.w3.org/WAI/ARIA/apg/practices/keyboard-interface/), [APG dialog modal](https://www.w3.org/WAI/ARIA/apg/patterns/dialog-modal/))
- **Performance:** web.dev DOM size/interactivity and layout thrashing ([web.dev DOM size](https://web.dev/articles/dom-size-and-interactivity), [web.dev layout thrashing](https://web.dev/avoid-large-complex-layouts-and-layout-thrashing/))
- **Repo-specific framing:** `src/dashboard/`, sidepanel, and popup files in this repo

## Quick review rules

1. **Prefer native semantic elements before custom widget code.** Native buttons, lists, landmarks, and dialogs beat hand-rolled `div` patterns.
2. **Use event delegation intentionally.** Review `target`, `currentTarget`, `closest()`, and composed-path assumptions together — most plain-JS UI bugs are mis-targeting bugs.
3. **Use `textContent` for text; treat `innerHTML` as a security and performance review hotspot.** HTML-string rendering requires explicit trust/sanitization scrutiny.
4. **Review focus and keyboard behavior together.** A UI is not accessible if it looks right but traps, loses, or hides focus.
5. **Prefer one-off or abortable listeners/observers/timers over immortal globals.** Large vanilla UIs often leak work through forgotten listeners or polling.
6. **Batch DOM reads before writes; be suspicious of large always-on sweeps.** DOM size and layout churn affect responsiveness directly.

## Review workflow

1. **Start with structure and semantics.** Identify landmarks, major sections, lists, headings, and interactive controls.
2. **Map event ownership.** Review direct vs delegated listeners, target resolution, and whether teardown exists.
3. **Audit DOM mutation surfaces.** Check for `innerHTML`, repeated full rerenders, and excessive polling or mutation-driven churn.
4. **Audit focus and keyboard flows.** Check initial focus, focus restoration, visible focus, tab order, and key bindings for custom widgets.
5. **Audit modal and composite-widget behavior.** Review dialog semantics, inert/background blocking, roving tabindex, `aria-selected`, `aria-expanded`, and state synchronization.
6. **Audit runtime work.** Review DOM size, repeated selectors, polling intervals, and layout-sensitive code for INP/responsiveness regressions.

## Review surfaces and checks

| Surface | Purpose | Review focus | Caveats |
|---|---|---|---|
| Event delegation | Reduce listener count and support dynamic DOM | `target`/`closest()` correctness, shadow DOM assumptions, cleanup | Delegation bugs often look like random UI bugs |
| `innerHTML` / HTML string rendering | Fast coarse-grained rendering | Trust boundary, escaping, rebind cost, focus loss | Easy XSS and state-loss hotspot |
| `textContent` / DOM node assembly | Safe text and fine-grained DOM updates | Simpler trust model, less parsing overhead | More verbose, but safer |
| Modal/dialog behavior | Block background work and focus | Initial focus, return focus, `aria-modal`, inert background | Native `<dialog>` still needs thoughtful usage |
| Composite widgets | Tabs/menus/accordions/tree controls | Tab vs arrow keys, roving tabindex, synced ARIA state | Prefer native controls when possible |
| Timers / sweeps / observers | Keep UI fresh | Visibility gating, teardown, selector cost, debounce | "Cheap" intervals add up in always-open panels |

## Standards and best practices

- Prefer **native controls and landmarks** over ARIA-heavy generic containers where HTML already models the concept.
- Keep **state, labels, and ARIA attributes synchronized** — a widget with stale `aria-expanded` or `aria-selected` is a real bug, not polish.
- Use **abortable listeners and explicit teardown** to keep long-lived panels from accumulating dead work.
- Review **periodic DOM sweeps and mutation observers** for visibility gating, debounce, and selector scope to avoid unnecessary CPU in an always-open sidepanel.
- For this repo, favor patterns that fit the existing large vanilla-JS sidepanel: readability, teardown, and semantic correctness matter more than framework-style abstractions.

## Known ambiguities

- APG patterns apply when custom widgets are truly necessary — do not force a complex ARIA pattern where a native element would be simpler.
- Performance guidance here is directional; confirm suspected hot paths with actual profiling.
- Some plain-JS UIs intentionally mix direct listeners, delegated listeners, and timers. The review question is whether that mix is controlled and comprehensible, not whether it is framework-like.
