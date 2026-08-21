<!-- hub-reference-banner -->
> **Reference file — part of the `writing-expert` hub.** Formerly the standalone `accessibility-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: accessibility-writing
description: "Writing craft for accessibility: heading structure (H1 → H2 → H3, no skipped levels), descriptive link text (no 'click here'), alt-text writing for decorative/informative/functional/complex images, form-label association, ARIA labels for icon-only controls, button labels by function (not visual), the 'skip to main content' link convention, color-independent information signaling, table captions and headers, language tags, captions and transcripts for video and audio, expandable/collapsible content patterns, and reading-level targeting per WCAG. Grounded in WCAG 2.2 (Level A/AA/AAA), Section 508, ADA Title III/II precedent, A11y Project, Deque, and WebAIM guidance. This is the WRITING side of accessibility — it covers what to put on the page in prose, alt attributes, labels, and link text. TRIGGER: user asks to write alt text, write accessible headings, rewrite a link/button label, write form labels, write a transcript or caption, write content that meets WCAG 2.2 / Section 508 / ADA, write descriptive button text, write accessible link text, replace 'click here' or 'read more', describe an image for a screen reader, structure headings on a page, write a 'skip to content' link. SKIP: code review of UI components (use accessibility-ux-reviewer); plain-language rewriting for general readability without an accessibility driver (use plain-language); WCAG compliance audit of an existing live site (use accessibility-ux-reviewer); designing color palettes for WCAG contrast (use frontend-design or ui-ux-pro-max); React/Vue ARIA component implementation (use frontend-design)."
origin: local
version: "1.0.0"
updated: "2026-05-29"
keywords:
  - accessibility writing
  - alt text
  - heading structure
  - descriptive link text
  - form labels
  - ARIA labels
  - WCAG 2.2
  - Section 508
  - WebAIM
  - A11y Project
  - Deque
  - skip to main content
  - color independence
  - captions
  - transcripts
  - reading level
  - lang attribute
tags:
  - writing
  - accessibility-writing
  - a11y
  - WCAG
  - inclusive
  - compliance
whenToUse:
  - User asks to write alt text for an image (decorative, informative, functional, or complex)
  - User asks to write or fix a heading hierarchy (no skipped levels, descriptive H2s)
  - User asks to rewrite "click here" / "read more" / "learn more" / URL-as-link-text into descriptive link text
  - User asks to write form-field labels, button labels for icon-only controls, or ARIA labels
  - User asks to write content that meets WCAG 2.2 Level A/AA, Section 508, or ADA requirements
  - User asks to write captions, transcripts, or audio descriptions for video/audio content
  - User asks to describe an image, chart, or diagram for a screen reader
  - User asks about language tags (`lang="en"`), reading-level targeting, or expandable-content patterns
  - User asks for a "skip to main content" or skip-navigation link
  - User asks how to convey information without relying on color alone
whenNotToUse:
  - User wants a code-level audit of a UI component or live site (use accessibility-ux-reviewer)
  - User wants plain-language simplification with no accessibility driver (use plain-language)
  - User wants color-palette contrast design (use frontend-design or ui-ux-pro-max)
  - User wants React/Vue/HTML ARIA component implementation patterns (use frontend-design)
  - User wants a general WCAG audit of an existing site (use accessibility-ux-reviewer)
related_skills:
  - accessibility-ux-reviewer
  - plain-language
  - inclusive-language
  - writing-expert
  - frontend-design
  - microcopy-and-ui-writing
---

# Accessibility Writing

Reference for writing content that is usable by screen-reader users, low-vision users, deaf and hard-of-hearing users, users with cognitive disabilities, and users on assistive tech. This is the **writing** side of accessibility: alt text, headings, link text, labels, captions. It is the partner to `accessibility-ux-reviewer` (which audits live UI) and to `plain-language` (which targets reading-grade level for general accessibility).

## When to use this skill

- User asks to write or fix alt text on an image
- User asks to fix a heading hierarchy (skipped levels, vague H2s)
- User asks to rewrite a link or button label
- User asks to write a form label, ARIA label, or icon-button label
- User asks to write a transcript or captions
- User asks how to describe a chart, diagram, or complex image for a screen reader
- User asks to meet WCAG 2.2, Section 508, or ADA Title II/III requirements
- User asks about the "skip to main content" link, color-independence, language tags, or reading-level targeting

## When NOT to use this skill

- Code-level component audit of a live site → `accessibility-ux-reviewer`
- General readability without an accessibility driver → `plain-language`
- Color contrast palette design → `frontend-design` / `ui-ux-pro-max`
- HTML/React/Vue ARIA component pattern implementation → `frontend-design`

## The 8-point accessibility-writing test

Every accessible draft should pass these eight tests:

1. **Headings form a proper outline.** Exactly one `<h1>` per page. `<h2>` for top-level sections. `<h3>` for subsections under each `<h2>`. No skipped levels (no `<h2>` followed by `<h4>`). Each heading describes what follows.
2. **Every image has an alt attribute.** Decorative images use `alt=""` (empty, not missing). Informative images get a 1–2 sentence description (≤140 characters typical). Functional images describe the action. Complex images get a short alt plus a long description.
3. **Every link reads as a self-contained label.** "Click here," "read more," "learn more," and bare URLs fail. A screen-reader user pulling up the page's link list must understand each link without surrounding context.
4. **Every form control has a programmatically associated label.** `<label for>` matched to input `id`, or wrapped `<label><input></label>`. ARIA labels only when a visible label is impossible.
5. **Every icon-only button has an accessible name.** `aria-label` describing the action ("Close dialog," not "X icon"). Decorative SVGs inside get `aria-hidden="true"`.
6. **Information is never conveyed by color alone.** Status indicated by "Error: …" prefix or icon plus color, not red text alone. Chart series differ by shape or pattern, not only by hue.
7. **The page declares its language.** `<html lang="en">` (or appropriate tag). Inline language changes wrapped with `lang` attribute (`<span lang="fr">`).
8. **Time-based media has alternatives.** Video has synchronized captions and a transcript. Audio has a transcript. Auto-playing audio has a pause control or does not auto-play.

## Core concepts

### 1. Alt text by image purpose

WCAG 1.1.1 (Non-text Content, Level A) requires a text alternative for every non-text content item. The rule branches by purpose:

**Decorative images** add no information. The image is purely visual flourish. The alt attribute is required but must be empty: `alt=""`. Missing alt attributes (no attribute at all) cause screen readers to read the file name or skip in inconsistent ways. The empty attribute tells the screen reader to ignore.

```html
<!-- Decorative divider -->
<img src="flourish.svg" alt="">
```

**Informative images** convey information that the surrounding text does not. Alt text should describe the information the image carries, not the image itself.

```html
<!-- Informative -->
<img src="cpu-chart.png" alt="CPU usage spikes to 95% at 14:30 UTC then returns to baseline by 14:35.">
```

**Functional images** are links or buttons. Alt text describes the action or destination, not the picture.

```html
<!-- Functional -->
<a href="/cart"><img src="cart-icon.svg" alt="View cart"></a>
<!-- NOT alt="shopping cart icon" -->
```

**Complex images** (graphs, charts, diagrams, infographics) cannot be summarized in 140 characters. Use a short alt for orientation plus a long description in adjacent text or via `aria-describedby`.

```html
<figure>
  <img src="latency-percentiles.png"
       alt="Latency percentile chart for production over the past 24 hours."
       aria-describedby="latency-desc">
  <figcaption id="latency-desc">
    The chart shows P50, P95, and P99 latency across 24 hours. P50 stays
    between 40 and 60 milliseconds. P95 trends around 180 milliseconds with
    two spikes at 02:00 and 14:00 UTC reaching 320 milliseconds. P99 follows
    the same pattern, peaking at 580 milliseconds.
  </figcaption>
</figure>
```

**Text-in-image** (a screenshot of a CLI command, a quote rendered as an image) must include the text in the alt or surrounding text. Better: render the text in actual text and style it, so screen readers, search engines, and copy/paste all work.

### 2. Heading structure

Headings are how screen-reader users navigate. They press `H` to jump to the next heading and `1` through `6` to jump to a specific level. A page with no headings or with arbitrary levels is unnavigable.

The rules:

- **One `<h1>` per page.** The page title.
- **No skipped levels going down.** After `<h2>` you may use another `<h2>` or an `<h3>`. You may not jump to `<h5>`.
- **You may skip levels going up.** After `<h4>` you may return to `<h2>` to start a new top-level section.
- **Headings are not for styling.** Don't use `<h3>` because it looks right; use the level that reflects the structural depth.
- **Descriptive headings.** "Introduction" and "More information" fail WCAG 2.4.6 (Headings and Labels, Level AA), which requires headings to describe the topic or purpose. "How to reset your password" passes.

A heading outline for a documentation page might be:

```
<h1>Setting up MongoDB Atlas
  <h2>Prerequisites
  <h2>Creating an Atlas cluster
    <h3>Choosing a cloud provider
    <h3>Sizing the cluster
  <h2>Connecting from your application
    <h3>Connection string format
    <h3>Authentication options
      <h4>SCRAM-SHA-256
      <h4>X.509 certificates
```

### 3. Link text

WCAG 2.4.4 (Link Purpose in Context, Level A) and 2.4.9 (Link Purpose Link Only, Level AAA) require that the purpose of a link be understandable from the link text. Screen-reader users frequently pull up a list of all links on a page using rotor navigation. In that list, only the link text appears.

Anti-patterns to rewrite:

| Anti-pattern | Why it fails | Rewrite |
|---|---|---|
| "Click here" | No context out of place | "View the full pricing table" |
| "Read more" | All "Read more" links sound identical | "Read more about MongoDB Atlas Search" |
| "Learn more" | Same problem | "Learn more about VPC peering setup" |
| "https://example.com/docs/atlas-cli" | URLs are read aloud character by character | "Atlas CLI documentation" |
| "this article" | Context-dependent | "the 2026 Atlas pricing changes article" |

Additional rules:

- Don't include the word "link" in link text. Screen readers already announce that an element is a link.
- For downloadable files, include the file format and (optionally) size: `<a href="/report.pdf">Q1 incident report (PDF, 1.2 MB)</a>`.
- For external links, indicate it ("opens in new tab") in the link text or via an aria-label, especially if behavior is different from same-tab links.
- Two links with the same text on the same page must point to the same destination. If two say "Read more" but go to different places, the screen-reader user cannot distinguish them.

### 4. Form labels

WCAG 1.3.1 (Info and Relationships, Level A) and 3.3.2 (Labels or Instructions, Level A) require every form control to have a visible, programmatically associated label.

The label-association options, in order of preference:

1. **`<label for>` matched to input `id`.** The most robust.
   ```html
   <label for="email">Email address</label>
   <input type="email" id="email" name="email">
   ```
2. **Wrapped `<label>`.** Works without explicit `for`/`id`.
   ```html
   <label>Email address <input type="email" name="email"></label>
   ```
3. **`aria-labelledby` pointing at a visible element.** Use when the label exists elsewhere on the page.
   ```html
   <span id="email-label">Email address</span>
   <input type="email" aria-labelledby="email-label">
   ```
4. **`aria-label`.** Last resort, only when no visible label is possible. The `aria-label` text becomes the accessible name but is invisible to sighted users.

For error messages, link the message to the field with `aria-describedby` and mark the field as invalid:

```html
<label for="email">Email address</label>
<input type="email" id="email" aria-describedby="email-err" aria-invalid="true">
<div id="email-err">Enter an email address in the format name@example.com.</div>
```

For required fields, use both the visible `*` (with a legend explaining what `*` means) and the `required` attribute. Don't rely on color alone for required status.

### 5. Icon-only buttons and accessible names

Icon-only buttons (a close X, a search magnifier, a hamburger menu) have no visible text. They need an accessible name that describes the **action**, not the **icon**.

```html
<!-- Wrong -->
<button><svg><!-- X icon --></svg></button>

<!-- Right -->
<button aria-label="Close dialog">
  <svg aria-hidden="true"><!-- X icon --></svg>
</button>
```

Two rules:

1. The `aria-label` describes what the button does ("Close dialog," "Search," "Open menu"). Not what it looks like ("X," "Magnifying glass," "Three lines").
2. The decorative SVG inside is hidden from screen readers with `aria-hidden="true"`. Otherwise the screen reader may read the SVG content in addition to the label.

For a button with both text and icon, the text usually carries the accessible name and the icon should be hidden:

```html
<button>
  <svg aria-hidden="true"><!-- Save icon --></svg>
  Save
</button>
```

### 6. Color independence (WCAG 1.4.1)

WCAG 1.4.1 (Use of Color, Level A) requires that color is never the only means of conveying information. Patterns that pass:

- **Status with a prefix word.** "Error: Email is invalid." "Warning: Unsaved changes." The color reinforces; the word carries meaning.
- **Status with an icon plus color.** A red X icon plus the word "Failed." A green check plus "Passed."
- **Links underlined as well as colored.** Removing the underline and signaling links by color alone fails for color-blind users.
- **Chart series differentiated by shape or pattern.** Line charts use solid vs dashed vs dotted lines, or markers of different shapes, in addition to color.
- **Form validation with text plus border color.** "Invalid email" message plus a red border. Not red border alone.

### 7. Captions, transcripts, audio descriptions

WCAG 1.2 requires alternatives for time-based media. The thresholds:

- **1.2.1 Audio-only and video-only (prerecorded), Level A.** Text alternative for the audio. For video-only, an audio or text description.
- **1.2.2 Captions (prerecorded), Level A.** Synchronized captions for prerecorded video with audio.
- **1.2.3 Audio description or media alternative (prerecorded), Level A.** Audio description or text transcript for prerecorded video.
- **1.2.4 Captions (live), Level AA.** Live captions for live audio in synchronized media.
- **1.2.5 Audio description (prerecorded), Level AA.** Audio description of prerecorded video.

**Captions** are synchronized text of dialogue plus non-speech information (music cues, laughter, sound effects). Captions go beyond subtitles — subtitles assume the viewer can hear the audio cues; captions don't.

**Transcripts** are the full text of dialogue and important sound, not synchronized. A transcript should include speaker labels, paragraph breaks, and timestamps every 30–60 seconds for navigation.

**Audio descriptions** narrate visual information during pauses in dialogue (a separate audio track, or an alternative version).

### 8. Language tags

WCAG 3.1.1 (Language of Page, Level A) and 3.1.2 (Language of Parts, Level AA) require that the language of the page and of any inline content in a different language is identified programmatically. Screen readers use this to switch pronunciation rules.

```html
<html lang="en">
  ...
  <p>The French phrase <span lang="fr">tour de force</span> means "feat of strength."</p>
```

Use IETF BCP 47 language tags: `en`, `en-US`, `en-GB`, `fr`, `es-MX`, `zh-Hans`, `zh-Hant`.

### 9. The "skip to main content" link

Keyboard and screen-reader users navigate the page sequentially. Without a skip link, every page requires tabbing through the entire navigation before reaching the main content. WCAG 2.4.1 (Bypass Blocks, Level A) requires a mechanism to skip repeated blocks.

The standard pattern:

```html
<body>
  <a href="#main" class="skip-link">Skip to main content</a>
  <header>...</header>
  <nav>...</nav>
  <main id="main">
    ...
  </main>
</body>
```

Important conventions:

- The skip link is the **first focusable element** on the page.
- It is **visually hidden until focused.** A common CSS pattern positions it off-screen and brings it on-screen when `:focus` applies. Hiding it permanently with `display: none` removes it from the focus order and breaks the feature.
- Its target is the `<main>` element or the first heading of main content.
- The text is clear: "Skip to main content" or "Skip navigation."

### 10. Reading-level targeting

WCAG 3.1.5 (Reading Level, Level AAA) suggests that when text requires a reading ability more advanced than lower secondary education level (around grade 9 in US terms), a supplemental or alternative version should be available. This is AAA, not AA, so it is aspirational for most sites — but legally adjacent for government and regulated content.

In practice:

- Aim for Flesch-Kincaid grade 8 for general audiences.
- Aim for grade 6 for health, financial, and legal information aimed at the general public.
- Use the `plain-language` skill for grade-level targeting craft.

### 11. Tables

Data tables need structure to be readable by screen readers:

- **`<caption>`** as the first child of `<table>`, describing what the table contains.
- **`<th scope="col">`** for column headers, **`<th scope="row">`** for row headers.
- For complex tables with multiple header levels, use `id` on the headers and `headers` on the cells.
- Avoid using tables for layout. If you do, add `role="presentation"`.

```html
<table>
  <caption>Quarterly revenue by region, 2026 Q1</caption>
  <thead>
    <tr>
      <th scope="col">Region</th>
      <th scope="col">Q1 2026</th>
      <th scope="col">YoY change</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <th scope="row">North America</th>
      <td>$4.2M</td>
      <td>+12%</td>
    </tr>
    ...
  </tbody>
</table>
```

### 12. Expandable/collapsible content (disclosure widgets)

For accordions, expand/collapse panels, and "Show more" patterns:

- The trigger is a `<button>` (not a `<div>`).
- The trigger has `aria-expanded="true"` or `"false"` reflecting current state, updated on toggle.
- The trigger has `aria-controls="<id of panel>"`.
- The panel is shown/hidden with `hidden` attribute or `display: none`. (CSS-only `visibility: hidden` may still be in the accessibility tree.)
- The label of the trigger remains the same when expanded or collapsed ("FAQ: Reset password," not "Show" vs "Hide").

## Templates

### Template 1: Alt-text rewrites

Take an image and write three alt-text candidates depending on purpose.

| Image | Purpose | Alt text |
|---|---|---|
| Photo of a smiling support agent next to a quote | Decorative (quote already in text) | `alt=""` |
| Photo of a smiling support agent illustrating "Our team is here for you" tagline | Informative (carries the message) | `alt="A support agent smiling and ready to help."` |
| Photo of a smiling support agent inside a link to /contact | Functional (the link is to contact support) | `alt="Contact support"` |
| Bar chart showing 30-day error rate trending down from 2.4% to 0.7% | Complex | Short: `alt="Error rate trend, last 30 days."` Long: include numbers in caption. |
| Screenshot of a CLI command: `mongosh --quiet --eval "db.serverStatus()"` | Text-in-image | `alt="Terminal command: mongosh --quiet --eval \"db.serverStatus()\""` (better: render as text) |
| Company logo in the header (linked to home) | Functional | `alt="Company Name — Home"` |
| Decorative gradient line dividing sections | Decorative | `alt=""` |

### Template 2: "Click here" rewrites

Each rewrite makes the link self-contained.

| Before | After |
|---|---|
| To see pricing, [click here](/pricing). | See [Atlas pricing](/pricing). |
| [Read more](/posts/123) | [Read about VPC peering setup](/posts/123) |
| Download the report [here](/q1.pdf). | Download the [Q1 2026 incident report (PDF, 1.2 MB)](/q1.pdf). |
| Documentation: [https://example.com/docs/atlas](https://example.com/docs/atlas) | See the [Atlas documentation](https://example.com/docs/atlas). |
| For details, see [this article](/help/x). | For details, see [How to reset your password](/help/x). |

### Template 3: Heading outline (rewrite)

Before (broken):

```
<h1>Documentation</h1>
<h3>Getting started</h3>
<h3>Advanced topics</h3>
<h5>Sharding</h5>
<h2>FAQ</h2>
```

Problems: skipped from `<h1>` to `<h3>`; skipped from `<h3>` to `<h5>`; vague heading "Documentation."

After:

```
<h1>MongoDB Atlas user guide</h1>
<h2>Getting started</h2>
<h3>Creating your first cluster</h3>
<h3>Connecting an application</h3>
<h2>Advanced topics</h2>
<h3>Sharding</h3>
<h3>Multi-region clusters</h3>
<h2>Frequently asked questions</h2>
```

### Template 4: Accessible form

```html
<form>
  <fieldset>
    <legend>Sign up for the newsletter</legend>

    <div>
      <label for="full-name">Full name</label>
      <input type="text" id="full-name" name="full-name" required
             aria-describedby="full-name-hint">
      <div id="full-name-hint">As you would like it to appear on emails.</div>
    </div>

    <div>
      <label for="email-addr">Email address</label>
      <input type="email" id="email-addr" name="email" required
             aria-describedby="email-hint email-err"
             aria-invalid="false">
      <div id="email-hint">We will not share your email.</div>
      <div id="email-err" role="alert"></div>
    </div>

    <div>
      <input type="checkbox" id="agree" name="agree" required>
      <label for="agree">
        I agree to the
        <a href="/terms">terms of service</a>.
      </label>
    </div>

    <button type="submit">Subscribe to newsletter</button>
  </fieldset>
</form>
```

### Template 5: Icon button library

```html
<!-- Close button -->
<button aria-label="Close dialog">
  <svg aria-hidden="true">...</svg>
</button>

<!-- Search button -->
<button aria-label="Search">
  <svg aria-hidden="true">...</svg>
</button>

<!-- Menu toggle -->
<button aria-label="Open navigation menu"
        aria-expanded="false"
        aria-controls="primary-nav">
  <svg aria-hidden="true">...</svg>
</button>

<!-- Favorite (toggle state) -->
<button aria-label="Add to favorites" aria-pressed="false">
  <svg aria-hidden="true">...</svg>
</button>

<!-- Delete -->
<button aria-label="Delete item">
  <svg aria-hidden="true">...</svg>
</button>
```

### Template 6: Status messages without relying on color

```html
<!-- Error -->
<div role="alert" class="error">
  <span aria-hidden="true">✗</span>
  <strong>Error:</strong> Email address is not valid. Use the format name@example.com.
</div>

<!-- Success -->
<div role="status" class="success">
  <span aria-hidden="true">✓</span>
  <strong>Success:</strong> Your changes have been saved.
</div>

<!-- Warning -->
<div role="alert" class="warning">
  <span aria-hidden="true">!</span>
  <strong>Warning:</strong> You have unsaved changes. Save before leaving.
</div>
```

### Template 7: Accessible video block

```html
<video controls preload="metadata">
  <source src="onboarding.mp4" type="video/mp4">
  <track kind="captions" src="onboarding.en.vtt" srclang="en" label="English" default>
  <track kind="descriptions" src="onboarding.descriptions.vtt" srclang="en" label="English audio descriptions">
</video>

<details>
  <summary>Read the transcript</summary>
  <div lang="en">
    <p><strong>[00:00] Narrator:</strong> Welcome to Atlas. In this video,
    we will walk through creating your first cluster.</p>
    <p><strong>[00:15] Narrator:</strong> First, sign in at cloud.mongodb.com…</p>
    ...
  </div>
</details>
```

### Template 8: Complex-image long description (chart)

```html
<figure>
  <img src="latency.png"
       alt="API latency over 24 hours. Detailed description below."
       aria-describedby="latency-long">
  <figcaption id="latency-long">
    The chart plots three latency lines (P50, P95, P99) over a 24-hour
    window from 2026-05-28 00:00 to 23:59 UTC.

    P50 stays between 38 and 62 ms throughout the window.

    P95 averages 180 ms with two notable spikes: a 320 ms spike at 02:00
    coinciding with the nightly backup job, and a 280 ms spike at 14:00
    coinciding with a traffic surge.

    P99 mirrors the P95 shape but peaks higher: 580 ms at 02:00 and 510 ms
    at 14:00.

    A horizontal SLO line at 250 ms is breached only by the two P95 spikes
    and by P99 between 01:30 and 03:00 and between 13:45 and 14:30.
  </figcaption>
</figure>
```

## Anti-patterns

1. **`alt="image of …"` or `alt="picture of …"`.** Screen readers already announce "image." The redundancy wastes time.
2. **`alt=""` for an informative image.** Empty alt is for decoration only. If the image carries any information, describe it.
3. **`alt` carrying the file name** (`alt="IMG_3421.png"`). Always meaningless. Always wrong.
4. **Heading levels chosen for visual size.** Use the correct semantic level and style with CSS.
5. **"Click here" / "read more" / bare URLs as link text.** Each fails the out-of-context test.
6. **Icon-only button with no `aria-label`.** Screen reader announces "button" with no name.
7. **Form field with no associated `<label>`.** Placeholder text disappears on input and is not a label.
8. **Color-only status indicators.** Red text without an icon or prefix word fails for color-blind users.
9. **`role="alert"` on every status message, including non-urgent ones.** `role="alert"` interrupts the screen reader. Reserve for genuinely urgent messages. Use `role="status"` for less urgent updates.
10. **A skip link that says "Skip" without telling you what you're skipping to.** "Skip to main content" is the convention.
11. **Captions vs subtitles confusion.** Captions include non-speech audio cues; subtitles assume hearing audio. Use captions for accessibility.
12. **Auto-playing video with audio.** Fails WCAG 1.4.2 (Audio Control, Level A) unless there is a mechanism to pause or stop within 3 seconds.
13. **Tables for layout without `role="presentation"`.** Screen readers announce "table with N rows and M columns" which is meaningless for layout.
14. **Long-form prose at college reading level for general-audience content.** Use the `plain-language` skill alongside this one.
15. **Trapping focus inside a modal without a close mechanism.** Modal must be dismissible via Escape key and a visible close button.

## Decision heuristics

| Situation | Choice |
|---|---|
| Image is purely decorative | `alt=""` |
| Image carries information not in surrounding text | Describe the information in 140 chars or fewer |
| Image is a link or button | Describe the action/destination, not the image |
| Image is a chart or complex diagram | Short alt + long description via `aria-describedby` or `<figcaption>` |
| Should I use `aria-label` on a form field? | No, if a visible `<label>` is possible. Yes only when impossible. |
| Should I use `role="alert"` on a status message? | Only if the user needs to be interrupted immediately (errors, warnings). Use `role="status"` otherwise. |
| Should I include "link" in link text? | No. Screen readers announce role separately. |
| Should the skip link be visually hidden? | Hidden until focused, then visible. Never `display: none`. |
| Should I disable autoplay? | Yes for video with audio. If autoplay is required, mute by default and provide a clear unmute control. |
| Should I support keyboard navigation? | Always. Every interactive element must be focusable and operable from keyboard. |
| Should I include audio descriptions? | Yes for prerecorded video where visual content is not conveyed in the dialogue (Level AA). |
| Should I aim for grade 8 reading level? | Yes for general audiences. Grade 6 for health, financial, legal. |

## Cross-skill notes

- **Use `accessibility-ux-reviewer` for code-level audits.** This skill writes content; that skill audits live UI for accessibility violations across HTML, CSS, and component behavior. They are complementary.
- **Use `plain-language` for grade-level targeting.** WCAG 3.1.5 reading-level targeting is a reading-grade problem; `plain-language` covers Flesch-Kincaid, common-word substitution, and sentence-length discipline. This skill handles structural accessibility (headings, alt, labels). Use both together for accessible long-form prose.
- **Use `inclusive-language` for terminology.** "Master/slave" → "primary/replica," "blacklist/whitelist" → "blocklist/allowlist," ableist metaphors, and gendered language fall under inclusive-language craft.
- **Use `microcopy-and-ui-writing` for button text and form copy** when accessibility is one consideration among many (brand voice, length, conversion). This skill is for the accessibility-specific subset of those decisions.
- **Use `frontend-design` for HTML/ARIA implementation patterns.** This skill tells you *what* to write; `frontend-design` covers *how* to wire it in React/Vue/HTML/CSS.

## References

1. W3C WCAG 2.2: https://www.w3.org/WAI/WCAG22/quickref/
2. W3C, *Understanding Success Criterion 2.4.4: Link Purpose (In Context)*: https://www.w3.org/WAI/WCAG21/Understanding/link-purpose-in-context.html
3. W3C, *Understanding Success Criterion 2.4.6: Headings and Labels*: https://www.w3.org/WAI/WCAG22/Understanding/headings-and-labels.html
4. W3C, *Alternative Text Tutorial*: https://www.w3.org/WAI/tutorials/images/
5. W3C, *Labeling Controls Tutorial*: https://www.w3.org/WAI/tutorials/forms/labels/
6. WebAIM, *Creating Accessible Forms — Advanced Form Labeling*: https://webaim.org/techniques/forms/advanced
7. The A11y Project, *Patterns*: https://www.a11yproject.com/patterns/
8. Deque University, *Accessibility Knowledge Base*: https://dequeuniversity.com/
9. Section 508, *Guide to Accessible Web Design & Development*: https://www.section508.gov/develop/guide-accessible-web-design-development/
10. MDN, *ARIA: aria-label attribute*: https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Reference/Attributes/aria-label

---

## Reading-grade formulas as accessibility instruments — Flesch-Kincaid, FOG, SMOG

**Rule.** Reading-grade formulas are an accessibility instrument. They estimate how educated a reader must be — measured in US school grades — to understand a passage on first reading. Lower grade = more accessible. Three formulas dominate, and they answer slightly different questions; an accessibility writer should know all three.

### The three formulas

| Formula | Year | Inputs | Predicts | Best for |
|---|---|---|---|---|
| Flesch-Kincaid Grade Level | 1975 (Kincaid et al., US Navy) | sentence length + syllables/word | US grade for *partial* comprehension | General-audience web content, business writing, default in MS Word and Hemingway |
| Gunning FOG Index | 1952 (Robert Gunning) | sentence length + % "hard words" (3+ syllables, common suffixes excluded) | US grade for *partial* comprehension | Business writing; older studies and many enterprise plain-language audits |
| SMOG (Simple Measure of Gobbledygook) | 1969 (Harry McLaughlin) | polysyllabic word density across 30 sentences | US grade for *complete* comprehension | Healthcare patient education (CDC, NIH, NCI), legal-adjacent and insurance disclosures, K-12 educational materials |

**Why SMOG specifically matters for accessibility.** SMOG is calibrated against *full* comprehension, not partial. If you ship content that a low-literacy reader must understand completely (medication instructions, evacuation guidance, eligibility for assistance, terms of service that gate a critical service), SMOG is the conservative choice. It typically scores 1–2 grades *higher* than Flesch-Kincaid on the same passage — not because the content is harder, but because the comprehension bar is stricter.

### Target grades by audience

| Audience / surface | Target US grade |
|---|---|
| General public (US adult literacy median is around grade 8) | 8 or below |
| Healthcare patient-facing materials (CDC, AMA) | 6–8 |
| K-12-aligned consumer materials | grade level of the target K-12 audience |
| Low-literacy populations (some adult literacy programs, ESL, disability-focused) | 5 or below |
| Government plain-language (PlainLanguage.gov; Plain Writing Act of 2010) | "as plain as the subject allows"; agencies often target 8 |
| Legal-adjacent (insurance, financial disclosures, ToS) | varies by jurisdiction; SEC plain-English rule uses Flesch Reading Ease ≥60, no formal grade target |
| Technical-writing for end users (Microsoft, Apple, Google docs guidance) | 8–10 typically |

### WCAG and reading-level requirements

WCAG 2.1/2.2 includes Success Criterion 3.1.5 (Reading Level) at Level AAA: *content should not require reading ability above the lower secondary education level (US grade 9 / around age 14–15) after removing proper names and titles*. This is a Level AAA criterion — not part of the standard AA conformance most jurisdictions require — but a useful target for any general-audience web content.

The current W3C **WCAG 3 working draft** (as of 2026) proposes more explicit reading-level thresholds as part of its Bronze/Silver/Gold conformance scheme. The draft has shifted multiple times; consult https://www.w3.org/TR/wcag-3.0/ for the current language. The high-level intent — *lower-secondary-education default, primary-education for safety-critical* — aligns with the CDC plain-language guidance and the AMA's grade-6 target for patient materials, so writing to those targets prepares you for WCAG 3 conformance regardless of where the draft lands.

### Operational workflow

1. **Score early, score often.** Run a draft through Flesch-Kincaid in MS Word or Hemingway during drafting, not after publication.
2. **Pair Flesch-Kincaid with SMOG for safety-critical text.** If FKGL says grade 7 and SMOG says grade 9, prefer SMOG — it is asking the stricter comprehension question. Aim for the *higher* of the two for high-stakes accessibility content.
3. **Score *after* removing proper names and titles** (WCAG 3.1.5 explicitly excludes these from the count). Most tools do not auto-exclude; do it manually or use a tool that supports the exclusion.
4. **Use the score as a flag, not a target.** A grade-8 passage with abstract content can be harder to understand than a grade-10 passage with concrete content. The score is necessary, not sufficient.

### Tools

- **Microsoft Word**: Review → Spelling & Grammar → enable readability statistics → report includes Flesch Reading Ease and Flesch-Kincaid Grade Level.
- **Hemingway Editor**: live grade-level meter, sentence-complexity highlighting. https://hemingwayapp.com/
- **Readable.com**: scores Flesch, FKGL, FOG, SMOG, ARI, Coleman-Liau in one pass. https://readable.com/
- **Python `textstat`**: programmatic access. `textstat.flesch_kincaid_grade(text)`, `textstat.gunning_fog(text)`, `textstat.smog_index(text)`. Useful for CI/CD readability gates on documentation pipelines.
- **CDC Clear Communication Index** (12-item rubric beyond grade-level): https://www.cdc.gov/ccindex/

### Worked example

Source passage (insurance benefits page):

> *Subject to the policy's eligibility provisions and exclusions, beneficiaries may submit claims for reimbursement of qualifying medical expenses incurred during the coverage period upon furnishing satisfactory documentation.*

Scores: FKGL ≈ 18, SMOG ≈ 16 — postgraduate. Inaccessible to >90% of US adults.

Rewrite to a grade-7 target:

> *If you have a covered medical expense during your policy period, you can ask us to pay you back. To do that, send us a copy of the bill and a short claim form. We will review it and send the payment within 14 days.*

Scores: FKGL ≈ 6, SMOG ≈ 7. Now within reach of the general adult population and broadly aligned with WCAG-AAA reading-level targets.

### When to break it

- *Technical reference documentation for trained professionals* (API references, developer guides, compiler manuals) does not have to hit consumer grade targets; the audience is trained.
- *Legal instruments where exact wording is mandated* (contracts, statutory text) cannot be simplified without changing meaning. Pair the original with a plain-language summary instead.
- *Proper names, technical terms, and titles* are excluded from the WCAG calculation but still count in raw FKGL/FOG/SMOG numbers — interpret accordingly.

### Composition with sibling skills

- Use `plain-language` for the deep rewrite techniques (common-word substitution, sentence-length discipline, jargon translation). That skill covers all three formulas and the rewrite craft in detail.
- This skill (accessibility-writing) handles the *structural* accessibility — heading hierarchy, alt text, form labels, link text — that grade-level scoring does not see. The two compose: a grade-6 passage with broken heading hierarchy is still inaccessible.
- For exact tools and library APIs, see this skill's `## Tools` section above.

### References

- Flesch, R. "A new readability yardstick." *Journal of Applied Psychology* 32(3), 1948, pp. 221–233.
- Kincaid, J. P., Fishburne, R. P., Rogers, R. L., & Chissom, B. S. *Derivation of New Readability Formulas for Navy Enlisted Personnel*. Naval Technical Training Command, Research Branch Report 8-75, 1975.
- Gunning, R. *The Technique of Clear Writing*. McGraw-Hill, 1952.
- McLaughlin, G. H. "SMOG Grading: A New Readability Formula." *Journal of Reading* 12(8), 1969, pp. 639–646.
- W3C, *Understanding Success Criterion 3.1.5: Reading Level (Level AAA)*. https://www.w3.org/WAI/WCAG21/Understanding/reading-level.html
- W3C, *WCAG 3 Working Draft*. https://www.w3.org/TR/wcag-3.0/
- CDC, *Simply Put — A Guide for Creating Easy-to-Understand Materials*. https://www.cdc.gov/healthliteracy/pdf/Simply_Put.pdf
- US Plain Writing Act of 2010 (Pub. L. 111-274) and PlainLanguage.gov guidelines.
