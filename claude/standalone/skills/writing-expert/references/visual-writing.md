<!-- hub-reference-banner -->
> **Reference file — part of the `writing-expert` hub.** Formerly the standalone `visual-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: visual-writing
version: "1.0.0"
updated: "2026-05-29"
category: custom
tags:
  - writing
  - visual-writing
  - alt-text
  - captions
  - data-viz
  - accessibility
description: >
  Writing the words that travel with images, charts, infographics, and video.
  Covers image captions, photo captions, video captions, sub-captions, alt text
  (the "describe what's there, why it's there, what it shows" triple), the W3C
  alt-decision-tree categories (decorative / informational / functional / image-
  of-text / complex), infographic copy, chart titles, axis labels, data-viz
  annotations, Edward Tufte's principles (data-ink ratio, chartjunk, sparklines,
  small multiples), Cole Nussbaumer Knaflic's "action title" rule ("title tells
  the story, axis labels confirm"), the "no chart without a takeaway" rule, the
  "describe the trend, then the value" caption pattern for accessible data-viz,
  and WCAG 1.1.1 alt-text patterns. References W3C WAI, WebAIM, Tufte, Knaflic,
  and the UK GDS accessibility blog.
  TRIGGER: "alt text", "image alt", "caption this", "photo caption", "chart
  caption", "chart title", "axis label", "infographic copy", "data viz
  annotation", "describe this chart", "accessible chart", "WCAG image", "alt
  decision tree", "decorative image", "screen-reader description", "what should
  I say about this image", "sparkline label", "figure caption", "long
  description".
  SKIP: pure accessibility audits of full pages or apps (use
  accessibility-ux-reviewer); broad UX/UI critique of a layout (use
  ui-ux-pro-max); writing prose around the image rather than for the image
  (use writing-expert); ARIA semantics or focus order debugging (use
  accessibility-ux-reviewer).
triggers:
  - "alt text"
  - "image alt"
  - "caption this"
  - "photo caption"
  - "chart caption"
  - "chart title"
  - "axis label"
  - "infographic copy"
  - "data viz annotation"
  - "describe this chart"
  - "accessible chart"
  - "WCAG image"
  - "alt decision tree"
  - "decorative image"
  - "screen-reader description"
  - "sparkline label"
  - "figure caption"
  - "long description"
skip:
  - full-page accessibility audit → use accessibility-ux-reviewer
  - layout or UX critique → use ui-ux-pro-max
  - prose around the image, not for the image → use writing-expert
  - ARIA roles or focus order → use accessibility-ux-reviewer
related:
  - accessibility-ux-reviewer
  - writing-expert
  - plain-language
  - storytelling-and-narrative
  - mongodb-atlas-charts
---

# Visual Writing

Reference for the words that ship with images, charts, and infographics: alt
text, captions, chart titles, axis labels, and annotations. The image is half
the message. This skill is the other half.

Deliver all responses in a direct, plain register. Avoid hedging and meta-commentary.

---

## When to use this skill

Activate when the user:

- Asks for alt text, image alt, or a screen-reader description
- Writes a caption for a photo, figure, video, or sub-caption
- Builds a chart and needs a title, axis labels, or annotations
- Wants an infographic reviewed for copy clarity
- Asks "what should I say about this image" for any accessibility purpose
- Needs the "describe the trend, then the value" pattern for an accessible chart
- Has a sparkline, small multiple, or in-text micro-chart that needs labelling
- Is writing the long description (sometimes called the `aria-describedby`
  body) for a complex chart

Skip when:

- The work is a full-page or app accessibility audit (use
  `accessibility-ux-reviewer`)
- The work is layout or UX critique (use `ui-ux-pro-max`)
- The work is body prose around the image (use `writing-expert` or
  `technical-writing-craft`)
- The work is ARIA roles, focus order, or keyboard nav (use
  `accessibility-ux-reviewer`)

---

## The one rule: visual writing makes images legible without the image

Every other rule in this skill follows from one premise: a reader who cannot
see the image must still receive the load-bearing information. That reader may
be a blind user with a screen reader, a sighted user on a slow network with
images failing to load, a sighted user skimming for the chart's takeaway, or a
search engine indexing the page.

The alt text serves the screen reader. The caption serves the skimmer. The
chart title serves the takeaway. The annotation serves the trend. They are
four jobs, not one — and the same content repeated four times is wrong.

---

## Core concept 1 — The alt-text triple

Alt text answers three questions in this order:

1. **What's there?** — concrete, observable content
2. **Why is it there?** — the function the image serves in the page
3. **What does it show?** — the specific information a sighted reader gets

The triple is contextual. Two pages can use the same image and need different
alt text. The W3C WAI tutorial uses the canonical example: a parrot photo
needs short alt on a "city parks" page ("a green parrot on a branch") and
long, species-specific alt on an ornithology page ("juvenile male Eclectus
parrot, *Eclectus roratus*, with green plumage and orange beak, perched on a
branch in profile").

**Rewrite pairs:**

Bad: `alt="image"`
Worse: `alt="picture of bird"`
Better: `alt="green parrot on a branch"`
Best for the bird-ID page: `alt="juvenile male Eclectus parrot, *Eclectus
roratus*, perched in profile; green plumage, orange beak"`

---

## Core concept 2 — The W3C alt-decision-tree categories

The W3C WAI tutorial categorizes every image into one of these buckets. Apply
the right pattern for the bucket:

| Image type | Alt text strategy | Example |
|------------|-------------------|---------|
| **Decorative** | Empty alt (`alt=""`) | A flourish, a hero-section gradient, a divider line |
| **Informational** | Describe the information | "Bar chart: 2024 revenue rose 18% over 2023" |
| **Functional** | Describe the function, not the icon | `alt="Search"` for the magnifying-glass icon, not `alt="magnifying glass"` |
| **Image of text** | Repeat the text verbatim | A logotype reading "Acme Corp" → `alt="Acme Corp"` |
| **Complex** | Short alt + long description nearby | `alt="Cluster topology diagram (long description below)"` then a `<details>` or `aria-describedby` body |
| **Group / collage** | One alt for the meaningful whole, empty alts on parts | A logo strip = `alt="our customers"` once; not five alts |
| **CAPTCHA** | Identify the test, not the content | `alt="visual CAPTCHA test; audio alternative available"` |

---

## Core concept 3 — Caption vs alt text vs long description

Three distinct artifacts. Do not duplicate.

- **Alt text** is for users who cannot see the image. It lives in `alt=""` or
  the figure's accessible name. It is short — typically 5 to 15 words.
- **Caption** is visible to everyone. It lives in `<figcaption>` or directly
  beneath the figure. It adds context, attribution, or commentary that the
  image alone does not carry.
- **Long description** is a structured prose or tabular alternative for
  complex images. It lives near the image or behind `aria-describedby`. It is
  used for diagrams, multi-series charts, schematics, and infographics where
  alt text cannot reasonably carry the load.

**Anti-pattern: duplication.** If `alt="Sales rose 18% in Q4"` and the caption
reads "Sales rose 18% in Q4," the screen-reader user hears the same sentence
twice. Either remove the alt (use `alt=""`) and let the caption do the work
via `aria-labelledby`, or split: alt for the trend, caption for the source
attribution.

---

## Core concept 4 — Tufte's principles

Edward Tufte's *The Visual Display of Quantitative Information* still anchors
visual writing. Five rules with direct copywriting implications:

1. **Above all else, show the data.** Title and labels exist to disclose the
   data, not to advertise the designer.
2. **Maximize the data-ink ratio.** Erase ink that doesn't carry data. If a
   label can be removed without losing information, remove it.
3. **Erase redundant data-ink.** A bar chart with a legend, a title, axis
   labels, and a redundant data table embeds the same fact four times. Pick one.
4. **Reject chartjunk.** Heavy gridlines, moiré fills, 3-D effects, drop
   shadows, and skeuomorphic icons add ink without adding data.
5. **Use sparklines for in-line trend.** A sparkline is a word-sized graphic
   embedded in a sentence. Caption a sparkline at the same density as a
   word — usually with just the latest value, e.g.
   "*latency last 30 days* ▁▂▃▅▇ **142 ms**".

For the writer specifically: every label is data-ink. Test each label by
removing it. If the chart is still legible, the label was junk.

---

## Core concept 5 — The Knaflic "action title" rule

Cole Nussbaumer Knaflic's *Storytelling with Data*: the chart title states the
takeaway. The axis labels confirm it.

**Generic title (descriptive only):**

> Quarterly Revenue, 2023–2024

**Action title (takeaway-bearing):**

> Q4 2024 revenue exceeded plan by 18% — the strongest quarter on record

The action title makes one assertion. The chart proves it. The reader does
not have to decode the chart to extract the message; the message is the
headline.

**Rule of thumb: no chart without a takeaway.** If you cannot state in one
sentence what the chart is for, the chart is decoration. Either find the
takeaway or delete the chart.

---

## Core concept 6 — Axis labels

Axis labels confirm the takeaway and add precision. Three rules:

1. **Label the units.** "Revenue (USD millions)" not "Revenue."
2. **Date axes get human dates.** "Jan 2024", "Apr 2024" — not raw
   timestamps, not "Q1".
3. **Avoid axis label rotation.** If labels overflow horizontally, the chart
   probably has too many categories. Aggregate or filter; don't rotate to 45°.

---

## Core concept 7 — Annotations

An annotation is a written assertion attached to a specific data point. Use
annotations for:

- The single highest or lowest point ("**peak: 8,420 — Aug 14**")
- A regime change ("← deploy of v2.3")
- An anomaly worth calling out ("backfill artifact; ignore")
- A trend that needs a name ("retention plateau begins here")

Annotations are journalism inside a chart. They commit to a claim. Empty
annotations ("note this") are wasted ink.

---

## Core concept 8 — Accessible data-viz captioning

The "describe the trend, then the value" pattern, from the UK GDS
accessibility blog and W3C complex-image guidance:

**Pattern: trend → key value → outlier.**

> Sales increased steadily from January to July, then declined through
> December. The peak was 8,420 units in July. The lowest month was December
> at 3,100 units.

Three sentences. Trend, peak, trough. A screen-reader user gets the same
top-line read a sighted user gets from a glance.

**Anti-pattern: read every value.**

> January was 4,200, February was 4,800, March was 5,100, April was…

This is data-table territory, not chart-alt-text territory. If every value
matters, provide a `<table>` and skip the chart alt-text recital.

---

## Core concept 9 — Infographic copy

Infographics are the worst-case for visual writing because they bundle
several charts plus copy in a single image file, often with no underlying
HTML. Three rules:

1. **Never publish an infographic as an image alone.** Always provide an
   HTML or PDF equivalent with extractable text.
2. **Each sub-element gets its own takeaway.** An infographic is N charts; it
   needs N takeaway titles, not one umbrella headline.
3. **Limit the copy budget to roughly 70 words per panel.** Beyond that, the
   reader stops reading and starts skimming the pictures.

---

## Core concept 10 — Sub-captions and video captions

**Sub-captions** are the second-tier line under a primary caption. Use them
for source attribution, date, or method note:

> *Figure 3. Q4 2024 revenue exceeded plan by 18%.*
> Source: internal finance, Feb 5 2026. Excludes the EMEA reseller channel.

**Video captions** (closed captions / subtitles) are a different domain. Key
rules:

- Caption every spoken word and every non-trivial sound effect.
- Use brackets for non-speech audio: `[laughter]`, `[applause]`, `[music
  swells]`.
- Identify speakers when off-screen: `MARCUS: We pushed the migration at
  midnight.`
- Line length: ≤ 42 characters per line, ≤ 2 lines on screen, ≥ 1 second per
  caption.
- Match the cadence of speech, not the script. People pause; captions break
  at the pause.

---

## Templates

### Alt-text decision shortcut

```text
Step 1: Is the image purely decorative? → alt=""
Step 2: Is it a functional icon (button, link)? → describe the FUNCTION
Step 3: Is it an image of text? → repeat the TEXT
Step 4: Is it informational and simple? → describe what's shown in 5-15 words
Step 5: Is it complex (chart, diagram, infographic)?
        → short alt + long description nearby
Step 6: Is it part of a group (logo strip, gallery)?
        → one accessible name for the whole, empty alts on parts
```

### Alt-text formula for charts

```text
[Chart type]: [trend / takeaway]. [Key value]. [Anomaly or outlier].

Examples:
"Line chart: API latency rose steadily from May to August, peaking at 480 ms
on Aug 14, then dropped back to 140 ms after the Aug 22 deploy."

"Bar chart: Q4 2024 revenue exceeded plan by 18% — the strongest quarter on
record."

"Stacked area chart: storage cost is dominated by hot-tier through Q3, then
shifts to cold-tier after the Sep migration."
```

### Chart takeaway title — before / after

| Generic descriptive | Action title |
|---------------------|--------------|
| "Quarterly Revenue" | "Q4 2024 was the strongest quarter on record" |
| "Customer Churn by Month" | "Churn doubled in March after the price change" |
| "API Latency, P95" | "Latency stabilized after the Aug 22 deploy" |
| "Headcount by Department" | "Engineering hiring drove 60% of 2024 growth" |
| "Cost Breakdown" | "Compute is 72% of spend; everything else combined is 28%" |

### Long description for a complex chart (template)

```html
<figure>
  <img src="topology.png"
       alt="Cluster topology diagram (long description below)">
  <figcaption>Figure 4. Production cluster topology, Feb 2026.</figcaption>
</figure>

<details>
  <summary>Long description of Figure 4</summary>
  <p>The cluster has three regions: us-east-1, us-west-2, and eu-west-1.
  Each region runs a 3-node replica set. The us-east-1 primary holds the
  write lease; us-west-2 and eu-west-1 are read-only secondaries with
  asynchronous replication. A separate analytics node in us-east-1 receives
  change-stream events for downstream BI workloads. Backups run from the
  us-east-1 hidden member at 03:00 UTC daily.</p>
</details>
```

### Photo caption template

```text
[Subject + active verb + context]. [Optional second-level: attribution, date,
note].

Examples:
"Migrant workers in California's Central Valley, 1936. Photo: Dorothea
Lange, Library of Congress."

"Engineering team at the v2.0 launch, Feb 14 2026."

"The decommissioned data center in Phoenix, days before demolition. The
site had been in continuous service since 1998."
```

### Sparkline + value (inline)

```text
"Latency last 30 days ▁▂▃▅▇ 142 ms (was 88 ms)"
"P95 ▁▂▃▂▁ 240 ms — stable"
"DAU ▂▃▅▇▆▇▇ 12.3k (+8% MoM)"
```

### Video caption block (CC convention)

```text
00:00:04.500 --> 00:00:08.000
MARCUS: We pushed the migration at midnight.

00:00:08.500 --> 00:00:12.000
[phone rings off-screen]

00:00:12.500 --> 00:00:16.000
MARCUS: And — of course — that's the on-call rotation.
```

---

## Anti-patterns

| Anti-pattern | Why it fails | Fix |
|--------------|--------------|-----|
| `alt="image"` or `alt="picture of X"` | Adds no information; screen reader already announces "image" | Describe what's shown |
| Alt text identical to the visible caption | Screen reader hears it twice | `alt=""` and `aria-labelledby` the caption, or split jobs |
| Reading every chart value in alt text | Long, useless to most users | Trend + key value + outlier; provide a data table if needed |
| Descriptive chart titles only ("Quarterly Revenue") | Reader has to extract the meaning | Action title with takeaway |
| Chartjunk: 3-D bars, drop shadows, moiré fills | Adds ink, hides data | Strip to data-ink |
| Axis labels in raw timestamps | "1717248000" is illegible | Human dates: "Jun 1 2024" |
| Legends with 7+ items | Reader cannot match color to series | Direct label on each series |
| Infographic published as single image only | Inaccessible, unsearchable | Provide HTML or PDF equivalent |
| One umbrella headline for a multi-chart infographic | Each chart's takeaway is buried | One takeaway per sub-chart |
| Captioning a video with only the script, not pauses | Captions lag or lead the audio | Caption breaks follow speech pauses |
| Sub-caption that repeats the caption | Wasted line | Use sub-caption for source / date / note |
| Functional icon alt'd as the icon ("magnifying glass") | Misses the function | `alt="Search"` |

---

## Decision heuristics

**Decorative vs informational:**

- Does the image carry information not in the surrounding text? → Informational
- Would a sighted reader miss anything if the image disappeared? → No → Decorative

**Caption vs alt:**

- Does the same info already appear in body prose? → Alt only; no visible caption
- Does the figure stand alone as a referenced object ("see Fig. 3")? → Caption + alt
- Is it complex? → Short alt + long description, optional caption for source

**Chart title:**

- Can you state the takeaway in one sentence? → Use it as the title
- You can't? → The chart needs work, not a title

**Annotation:**

- Will the reader's eye go to that data point on its own? → Maybe no annotation needed
- Is there a non-obvious story (deploy, anomaly, regime change)? → Annotate it

**Sparkline labelling:**

- Is the chart embedded in prose at word-density? → Label with last value only
- Standalone? → Full chart treatment, not a sparkline

**Long description trigger:**

- More than 3 series, more than 12 data points, or a diagram → Long description
- Single-series trend → Alt text alone is enough

---

## References

- W3C WAI: [Images
  Tutorial](https://www.w3.org/WAI/tutorials/images/) and the [Alt
  Decision Tree](https://www.w3.org/WAI/tutorials/images/decision-tree/)
- WebAIM: Alternative Text article (`webaim.org/techniques/alttext/`)
- Edward Tufte, *The Visual Display of Quantitative Information* (2nd ed.,
  2001) and *Beautiful Evidence* (2006); essays at
  `edwardtufte.com/notebook/chartjunk/`
- Cole Nussbaumer Knaflic, *Storytelling with Data* (2015) and *Storytelling
  with Data: Before & After*
- UK GDS: [Text descriptions for data
  visualisations](https://accessibility.blog.gov.uk/2023/04/13/text-descriptions-for-data-visualisations/)
- WCAG 2.1, Success Criterion 1.1.1 (Non-text Content)
- TPGi / Vispero: ["Making data visualizations
  accessible"](https://www.tpgi.com/making-data-visualizations-accessible/)
- A11Y Collective: ["The Ultimate Checklist for Accessible Data
  Visualisations"](https://www.a11y-collective.com/blog/accessible-charts/)

---

## Related skills

- `accessibility-ux-reviewer` — for full-page or app-level accessibility
- `writing-expert` — for prose around the image
- `plain-language` — when captions must hit a reading-grade target
- `storytelling-and-narrative` — for the takeaway under an action title
- `mongodb-atlas-charts` — for Atlas-specific chart-building inside Charts
