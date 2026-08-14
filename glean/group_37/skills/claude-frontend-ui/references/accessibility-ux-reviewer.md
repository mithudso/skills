<!-- hub-reference-banner -->
> **Reference file — part of the `frontend-ui` hub.** Formerly the standalone `accessibility-ux-reviewer` skill.
> Sibling topics in this family are now reference files under the hubs (`frontend-ui`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: accessibility-ux-reviewer
version: "1.1.0"
updated: "2026-05-29"
description: >
  Accessibility and UX review against WCAG 2.2, WAI-ARIA Authoring Practices Guide, and MDN
  guidelines. TRIGGER: user asks to audit, review, or check UI for accessibility; asks about
  WCAG compliance, keyboard navigation, focus management, ARIA usage, screen reader behavior,
  form labeling, landmark regions, or color contrast. SKIP: user wants to build a new component
  (use frontend-design), run automated Lighthouse scans (use chrome-devtools-mcp), or review
  pure visual design without accessibility concern.
origin: local
whenToUse:
  - "Audit this component for accessibility issues"
  - "Check keyboard navigation on this UI"
  - "Is this ARIA usage correct?"
  - "Review focus management in this modal"
  - "Does this meet WCAG 2.2?"
  - "Check color contrast"
  - "Review form labeling"
  - "Screen reader audit"
  - "Check landmark regions"
  - "Accessibility review before launch"
related_skills:
  - frontend-ui
  - chrome-devtools-mcp
---

# Accessibility and UX Reviewer

Expert reference for auditing UI against WCAG 2.2, WAI-ARIA APG, and MDN accessibility guidelines.
Produces structured findings with severity, location, and remediation.

## When NOT to use this skill

- **Building new components:** use `frontend-design` — it includes a pre-ship accessibility checklist.
- **Automated scans only:** use `chrome-devtools-mcp` for axe-core and Lighthouse runs.
- **Pure visual/aesthetic review without accessibility concern:** use `ui-ux-pro-max`.
- **Shadow DOM overlay isolation:** use `shadow-dom-component-authoring`.

## Review output format

For each finding, report:

```
[SEVERITY] <Location> — <Finding>
  Rule: <WCAG criterion or APG pattern>
  Fix: <specific remediation>
```

Severity levels: **Critical** (blocks AT users), **High** (impairs AT users significantly), **Medium** (reduces usability), **Low** (polish).

Example:
```
[Critical] <button class="icon-btn"> — Missing accessible name
  Rule: WCAG 2.4.6 / APG names-and-descriptions
  Fix: Add aria-label="Close dialog" or visually-hidden text
```

## Authoritative sources

| Source | Use for |
|--------|---------|
| [WCAG 2.2](https://www.w3.org/TR/WCAG22/) | Conformance criteria (A/AA/AAA) |
| [WAI APG](https://www.w3.org/WAI/ARIA/apg/) | Widget patterns, keyboard conventions |
| [APG keyboard interface](https://www.w3.org/WAI/ARIA/apg/practices/keyboard-interface/) | Tab order, focus management, internal widget navigation |
| [APG names and descriptions](https://www.w3.org/WAI/ARIA/apg/practices/names-and-descriptions/) | Accessible name techniques |
| [APG landmark regions](https://www.w3.org/WAI/ARIA/apg/practices/landmark-regions/) | Page structure for AT navigation |
| [MDN HTML accessibility](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Accessibility/HTML) | Semantic HTML baseline |
| [MDN CSS/JS accessibility](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Accessibility/CSS_and_JavaScript) | Styling and scripting pitfalls |
| [web.dev accessibility](https://web.dev/learn/accessibility) | Practical implementation reference |

**Version baseline:** WCAG 2.2 (current W3C recommendation). WCAG 3 drafts do not replace WCAG 2.2 for present-day audits.

## Review workflow

Run checks in this order — each step can surface Critical/High findings early:

1. **Semantic structure** — headings, landmarks, lists, buttons, links, form controls before looking at ARIA or scripts.
2. **Accessible names and labels** — all focusable interactive elements, dialogs, regions, and relevant containers.
3. **Keyboard behavior and focus** — tab order, visible focus, internal widget navigation, focus movement between components.
4. **Form behavior and status communication** — labels, instructions, error messages, dynamic UI changes.
5. **Responsive readability and zoom behavior** — text legibility, spacing, 400% zoom reflow.
6. **ARIA-heavy components** — only after steps 1–5 pass; evaluate against APG widget patterns.

## Quick review rules

| # | Rule | Authority |
|---|------|-----------|
| 1 | Prefer semantic HTML; only reach for ARIA when native elements can't express the relationship | MDN HTML accessibility |
| 2 | Every interactive element must be keyboard-operable | APG keyboard interface |
| 3 | Focus must be visible, predictable, and meaningful | APG keyboard interface |
| 4 | Every focusable interactive element requires an accessible name | APG names and descriptions |
| 5 | Use landmarks so AT users can navigate page structure | APG landmark regions |
| 6 | Use ARIA to add semantics for complex controls — not to paper over fixable HTML problems | MDN Accessibility |
| 7 | Preserve expected element appearance/behavior when restyling; never style a `div` to act as a button | MDN CSS/JS accessibility |
| 8 | WCAG compliance improves usability for all users, not just those with disabilities | WCAG 2.2 |

## Accessibility audit inventory

### Semantic structure

| Item | Key requirements | Failure mode |
|------|-----------------|--------------|
| Semantic HTML element choice | Use elements for their intended purpose | Generic containers where native elements exist |
| `<button>` | Use real button markup; activatable via keyboard | Non-button styled as button loses built-in behavior |
| Landmark roles | `main`, `nav`, `aside`, named `section`, `header`/`footer` at top level | Landmarks without meaningful structure add noise |
| Heading hierarchy | Sequential h1→h6, one `<h1>` per page, no level skip | Headings used for visual size rather than document structure |

### Names and labels

| Item | Techniques | Failure mode |
|------|------------|--------------|
| Accessible name | `aria-label`, `aria-labelledby`, visible child content | Missing name blocks AT users entirely |
| `aria-label` | Explicit string; use when no visible label text exists | Drifts out of sync with visible UI |
| `aria-labelledby` | References existing visible content by ID | Broken ID or weak source text |
| `aria-describedby` | Supplemental description; supports, not replaces, the name | Replacing a name with only a description |
| `<label for="id">` | Native form control naming; preferred over ARIA for inputs | Placeholder-only label loses on focus |
| `<legend>` for `<fieldset>` | Groups related controls; communicates group purpose | Missing group name harms form comprehension |

### Keyboard and focus

| Item | Requirements | Failure mode |
|------|-------------|--------------|
| Visible focus | Styling + logical movement; 2px+ indicator, 3:1 contrast (WCAG 2.4.13 AA) | `outline: none` with no replacement |
| Tab order | Matches visual/logical order; no positive `tabindex` | Unpredictable focus sequence |
| Composite widgets | Arrow-key internal navigation; Tab exits the widget | Custom widget with no keyboard support at all |
| Focus trapping | Modals trap Tab cycle; focus returns to trigger on close | Focus escapes to background or is lost |
| Custom ARIA widgets | Authors must implement keyboard behavior — browsers don't | ARIA widget with mouse-only interaction |

### Forms

| Item | Requirements | Failure mode |
|------|-------------|--------------|
| Labels | `<label>` explicitly associated to every input | Placeholder-only label |
| Error messages | Appear near the field; use `role="alert"` or `aria-live` for dynamic errors | Errors only at page top with no field connection |
| Required fields | Marked with `required` attribute and visual indicator | Asterisk with no explanation |
| Instructions | Present before the form, not only in placeholder | Context lost after user starts typing |

### Contrast and readability

| Target | Minimum | Enhanced |
|--------|---------|---------|
| Normal text | 4.5:1 (WCAG AA 1.4.3) | 7:1 (AAA) |
| Large text (18pt+ or 14pt bold+) | 3:1 (AA) | 4.5:1 (AAA) |
| UI components and focus indicators | 3:1 (WCAG 2.4.11/2.4.13) | — |

## Common failure modes and fixes

| Failure | Fix |
|---------|-----|
| `<div onClick>` instead of `<button>` | Replace with `<button>` |
| Icon-only button with no label | Add `aria-label` or visually-hidden text |
| Animated element with no `prefers-reduced-motion` | Add `@media (prefers-reduced-motion: reduce)` override |
| Custom widget with no keyboard support | Implement APG widget pattern keyboard conventions |
| `role="button"` on a `<div>` with no `tabindex="0"` | Use real `<button>` or add `tabindex="0"` + keyboard handler |
| Color as sole error indicator | Add icon or text label alongside color |
| Modal that doesn't trap focus | Implement focus trap on open; return focus to trigger on close |
| Dynamic content change not announced | Use `aria-live="polite"` for non-urgent or `aria-live="assertive"` for critical |

## Practical decision prompts

When starting a review, ask in sequence:

1. **Can this be a native HTML element instead of a custom widget?** If yes, use it.
2. **Can a keyboard-only user complete every flow with visible, predictable focus?** Test with Tab, Shift+Tab, Arrow keys, Enter, Escape.
3. **Does every meaningful interactive or structural region have the right name, landmark, or description?**
4. **Are dynamic updates (errors, loading states, notifications) announced to AT users?**

## Notes on WCAG levels

- **Level A:** Minimum compliance — critical barriers removed.
- **Level AA:** Legal standard in most jurisdictions (EAA, ADA, EN 301 549). Target for all production web content.
- **Level AAA:** Enhanced — not required for full conformance; apply where feasible.

APG is a pattern and authoring-practice guide, not a replacement for semantic HTML or WCAG conformance requirements.
