<!-- hub-reference-banner -->
> **Reference file — part of the `da-applied-and-communication` hub.** Formerly the standalone `da-8-data-visualization` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-8-data-visualization
title: Data Visualization
version: "1.0.1"
updated: "2026-05-30"
category: data-analysis
origin: local
description: >-
  Data visualization — chart selection by analytical intent, perception (Cleveland & McGill,
  Gestalt, preattentive), Tufte data-ink/chartjunk, Knaflic storytelling, color theory
  (sequential/diverging/categorical, colorblind-safe), WCAG chart accessibility, dashboard design
  (Few, BAN tiles), grammar of graphics, libraries (D3, Plotly, Vega-Lite, matplotlib, ggplot2,
  seaborn), BI platforms (Tableau, Power BI, Looker, Atlas Charts), anti-patterns (3D,
  truncated/dual axis, pie overuse, rainbow). TRIGGER: which chart should I use; colorblind-safe
  palette; ggplot2 vs matplotlib vs plotly; Tableau vs Power BI; fix a misleading chart; BAN KPI
  tiles. SKIP: EDA → da-5-exploratory-data-analysis; hypothesis testing →
  da-1-4-statistical-inference-foundations; model fits → da-6-statistical-modeling; chart
  narrative → da-9-reporting-communication / writing-expert; Atlas Charts setup →
  mongodb-atlas-charts; report copy → technical-writing-craft.
triggers:
  - data visualization
  - chart selection
  - which chart should I use
  - bar chart vs pie chart
  - dashboard design
  - data-ink ratio
  - chartjunk
  - small multiples
  - sparkline
  - preattentive attributes
  - storytelling with data
  - declutter a chart
  - action title
  - BAN tile
  - color palette
  - sequential palette
  - diverging palette
  - categorical palette
  - Viridis
  - Cividis
  - ColorBrewer
  - colorblind-safe chart
  - WCAG chart
  - chart accessibility
  - alt text for chart
  - grammar of graphics
  - ggplot2 layered
  - matplotlib vs seaborn
  - D3 vs Plotly
  - Vega-Lite
  - Observable Plot
  - Tableau vs Power BI
  - Looker Studio
  - Metabase
  - Superset
  - Atlas Charts
  - truncated y-axis
  - dual y-axis
  - misleading chart
  - Gestalt principles
  - Cleveland and McGill
keywords:
  - data-visualization
  - chart-selection
  - Tufte
  - data-ink-ratio
  - chartjunk
  - small-multiples
  - sparkline
  - Cleveland-McGill
  - perceptual-ranking
  - Gestalt
  - preattentive
  - storytelling-with-data
  - Knaflic
  - dashboard-design
  - Few
  - BAN-tile
  - bullet-chart
  - color-palette
  - sequential
  - diverging
  - categorical
  - Viridis
  - Cividis
  - ColorBrewer
  - Okabe-Ito
  - colorblind-safe
  - WCAG
  - accessibility
  - alt-text
  - grammar-of-graphics
  - Wilkinson
  - Wickham
  - ggplot2
  - matplotlib
  - seaborn
  - altair
  - plotly
  - D3
  - Vega-Lite
  - Observable-Plot
  - Bokeh
  - Tableau
  - Power-BI
  - Looker-Studio
  - Metabase
  - Superset
  - Atlas-Charts
  - anti-patterns
  - 3D-chart
  - dual-axis
  - rainbow-colormap
  - jet-colormap
when_to_use:
  - Choosing a chart type for a data shape and analytical intent
  - Fixing a misleading, cluttered, or inaccessible chart
  - Picking a colorblind-safe palette or designing a dashboard
when_not_to_use:
  - You need to discover patterns in raw data — use da-5-exploratory-data-analysis
  - You need to run a hypothesis test — use da-1-4-statistical-inference-foundations
  - You need to fit a model — use da-6-statistical-modeling or da-7-machine-learning
  - You need to write the narrative around finished charts — use da-9-reporting-communication or writing-expert
  - You need product-specific Atlas Charts setup — use mongodb-atlas-charts
  - You need to author text-only report copy with no visual artifact — use writing-expert
related_skills:
  - writing-expert
  - da-applied-and-communication
  - da-analytical-methods
  - technical-writing-craft
---

# Data Visualization

**Taxonomy context:** Data Analysis > Data Visualization

Data visualization is the principled translation of data into visual form so that the human
visual system — which is fast, parallel, and pattern-seeking — can extract structure that the
verbal system alone cannot. Done well, a chart compresses a table of numbers into a single
glance and surfaces the comparison the audience needs to make a decision. Done poorly, it
distorts, distracts, or actively deceives. This skill is the principled middle: which
encodings the eye actually reads accurately, which design choices preserve truth, which
patterns survive accessibility and color-blindness, and which tool to reach for given the
audience, the medium, and the question.

Two modes of visualization show up repeatedly and are often conflated:

- **Exploratory visualization** — the analyst draws many quick charts to *find* something.
  Iteration speed matters more than polish. Tools: matplotlib, seaborn, ggplot2, altair,
  Observable Plot, BI tool exploration views.
- **Explanatory visualization** — the analyst has found the point and now needs the audience
  to see it. Polish, narrative, and accessibility matter more than iteration speed. Tools:
  the same plus Tableau / Power BI / Looker Studio dashboards, D3 / Plotly for the web,
  publication-grade matplotlib / ggplot2 figures.

The principles below apply to both, but the bar for explanatory work is higher: a reader
gets one pass, no follow-up question, and you do not get to clarify in voice.

---

## 1. Perceptual foundation — what the eye reads accurately

Cleveland and McGill (1984, *Journal of the American Statistical Association*) ran the
empirical work that anchors every modern chart-selection guide. They asked subjects to read
quantitative values off ten elementary perceptual tasks and ranked the encodings by accuracy
of magnitude estimation:

1. **Position along a common scale** — most accurate. Bar chart, dot plot, scatterplot on a
   shared x-axis.
2. **Position along non-aligned scales** — small-multiples bars on separate panels.
3. **Length, direction, angle** — slope of a line, length of a bar without a baseline,
   angle of a pie slice. Less accurate than position.
4. **Area** — bubble chart, treemap rectangle. Read with systematic underestimation.
5. **Volume, curvature** — 3D charts, density of curves. Worse than area.
6. **Shading, color saturation, hue** — least accurate for quantitative magnitude. Use for
   category, not for ordered comparison of unknown values.

**Practical implication:** when the audience must compare values, encode the comparison
along the *position* channel. Reach for area, angle, or color saturation only when position
is already taken by a more important dimension or when the chart's job is qualitative
pattern recognition rather than precise estimation.

This single ranking explains most of the chart-choice rules of thumb. Bar charts beat pie
charts because length on a common scale beats angle and area. Dot plots and slope graphs
beat connected polar charts. A heatmap is fine for spotting hotspots but bad if the audience
must read individual cell values.

---

## 2. Tufte's principles — data-ink, chartjunk, small multiples

Edward Tufte's *The Visual Display of Quantitative Information* (1983) is the canonical
treatment of the design economics of a chart. The book has roughly 29,000 citations as of
2025 and most of its vocabulary has entered the field.

**Data-ink ratio.** Tufte defines *data-ink* as the non-erasable core ink that represents
the data, and proposes maximizing the ratio of data-ink to total ink. Two erasing rules:
"erase non-data-ink, within reason" and "erase redundant data-ink, within reason." The
practical effect is removing chart borders, redundant tick marks, heavy gridlines, drop
shadows, gradient fills, and decorative backgrounds. The chart should *be* the data, not
contain it.

The data-ink ratio is a heuristic, not a fundamentalist commandment. A small amount of
non-data-ink — a clear axis label, a reference line, a brand mark — earns its space when
it reduces the cognitive cost of reading the chart. Frank Elavsky and others have argued
against fundamentalist minimalism in recent work, and accessibility requirements
(visible borders, contrast minimums) sometimes mandate "non-data" ink. Treat the ratio as
a question to ask, not a score to maximize.

**Chartjunk.** Tufte's term for unnecessary or distracting elements that do not contribute
to understanding — 3D bevels, animated gradients, decorative iconography overlaid on bars,
moiré patterns from heavy crosshatching. Chartjunk steals attention from the data and often
distorts it. The decorative pie chart embedded in a sales brochure with shadows and
exploding slices is the canonical example.

**Small multiples.** A grid of the same chart type repeated for each level of a
categorical variable, with shared axes and scales. Small multiples let the eye scan many
panels in parallel because the chart grammar is constant — only the data changes. They
replace stacked bars, multi-series line spaghetti, and large legends. Tufte may have been
the first to popularize the term, though Galileo's sunspot drawings are an early example.

**Sparklines.** Word-sized, axis-free, label-free trend lines designed to sit inline with
prose or in dashboard tables next to the most recent value. The point is not precision but
direction and shape recognition: "the last quarter looked like this." Sparklines have
become a standard control in BI tools.

**Graphical integrity.** Tufte's "Lie Factor" = (visual change shown in graphic) / (actual
change in data). A graphic with a Lie Factor materially different from 1 is misleading.
This formalizes the harm of truncated axes, area-encoded bar charts, and 3D perspective
charts that exaggerate the foreground.

---

## 3. Knaflic / Few — explanatory design and storytelling

Stephen Few (*Information Dashboard Design*, *Now You See It*) and Cole Nussbaumer Knaflic
(*Storytelling with Data*) translate Tufte's principles into operational rules for
business and dashboard audiences.

**Preattentive attributes.** Color, size, position, orientation, motion, and a handful of
other visual properties are processed by the visual system in roughly 200 ms — *before*
conscious attention. A single red bar in a field of grey bars is seen as red before the
viewer has decided to look at it. Knaflic's central technique: use preattentive attributes
to direct attention to the one thing the chart is about. Make the rest grey.

**Decluttering.** Knaflic demonstrates a six-step decluttering sequence applied to a
generic stock chart: remove the border, remove the gridlines, remove redundant tick marks,
soften the axis labels, push the legend or label directly onto the line, and finally apply
preattentive emphasis to the focal series. The result is dramatically more readable than
the default Excel output.

**Action titles.** Replace neutral descriptive titles ("Quarterly Revenue") with titles
that state the finding ("Revenue declined 12 percent in Q3, driven by APAC"). The title
carries the message. The chart is the evidence. If a reader sees only the title, they
should still know the point.

**Dashboard design — Few's rules.** *Information Dashboard Design* enumerates thirteen
common mistakes. The big ones: exceeding the boundaries of a single screen, requiring
scrolling, supplying inadequate context, displaying excessive detail or precision,
choosing inappropriate display media, encoding quantitatively when categorical encoding
would suffice, decorating gratuitously, misusing or overusing color, and arranging the
data poorly. A dashboard is not a report. It is a continuously updated, at-a-glance
monitoring surface, and every chart on it must earn its pixel budget.

**BAN tiles (Big-Ass Numbers).** Dashboard pattern: a single very large numeric value
(plus a small trend indicator or sparkline) representing the headline KPI. The number is
the primary visual element, sized so it can be read from across the room. BAN tiles
exploit position-zero attention and are the dashboard analog of an action title.

---

## 4. Gestalt principles — how grouping is perceived

The Gestalt school's perceptual grouping laws describe how the visual system parses a
scene into objects without conscious effort. Six laws matter for charts and dashboards:

- **Proximity** — items physically close are read as a group. Use whitespace between
  unrelated panels and tight spacing within a panel. Two charts that belong together
  should not be separated by more space than the elements inside each chart.
- **Similarity** — items with the same color, shape, or size are read as a group, even
  when separated. Use the same color for "revenue" across every chart in a dashboard so
  the reader can carry the meaning across panels.
- **Enclosure** — items inside a visible border are read as a group. Card layouts,
  panels, and quadrants exploit this. Borders are non-data-ink but earn their place when
  they replace ambiguity about grouping.
- **Closure** — open shapes are perceived as closed if a reasonable closure exists. Most
  charts can omit the top and right axis lines (the "L-shape" axis) and the eye still
  reads the chart area as bounded.
- **Continuity** — the eye follows lines and aligned edges. Charts arranged on a grid
  with consistent alignment scan more easily than a jumble.
- **Connection** — items joined by a line are read as more strongly grouped than items
  with the same color. The line in a line chart is therefore a strong assertion that the
  points are sequential and comparable.

Gestalt is the design vocabulary that turns Tufte's "data-ink" abstraction into concrete
layout decisions: *why* removing a border works (closure), *why* small multiples scan
quickly (similarity and proximity), *why* a stray red dot pops (similarity-violation).

---

## 5. Chart selection by intent

Andrew Abela's *Chart Chooser* organizes charts into four intents: **comparison**,
**composition**, **distribution**, **relationship**. Add **temporal** (change over time)
and **spatial** (geographic) and you have a working taxonomy. The matrix below pairs
intent with the standard go-to charts.

| Intent | Data shape | First-choice chart | Alternates | Avoid |
| --- | --- | --- | --- | --- |
| **Comparison among categories** | 1 categorical, 1 numeric | Horizontal bar chart, sorted | Dot plot, lollipop chart | Pie chart, 3D bar |
| **Comparison over time** | Time series, 1+ numeric | Line chart | Area chart (one series), slope graph, sparkline | Stacked area for >3 series |
| **Composition (parts of whole, static)** | 1 categorical, 1 numeric, summing to total | 100% stacked bar, treemap, waffle chart | Pie chart (only when ≤3 slices and gap is obvious) | Donut, exploded pie, 3D pie |
| **Composition over time** | Time + categorical + numeric | 100% stacked area, stream graph | Small-multiples lines | Stacked area with shifting baseline |
| **Distribution of one variable** | 1 numeric, many obs | Histogram, density plot, boxplot | Violin plot, strip plot, ECDF | Pie chart of binned ranges |
| **Distribution by group** | 1 numeric × 1 categorical | Small-multiples histograms, boxplots side-by-side, ridgeline | Violin plot | Overlapping density curves with many groups |
| **Relationship between two numerics** | 2 numeric | Scatterplot | Hexbin (large n), 2D density, contour | Connected scatter without clear order |
| **Relationship + third variable** | 3 numeric or 2 num + 1 cat | Scatter + color or size, small multiples scatter | Bubble chart (with care) | Tripled-up dual axis |
| **Spatial** | Has lat/lon or region | Choropleth (rates), proportional symbol (counts), dot density | Cartogram for population-weighted comparison | Choropleth of raw counts (small regions vanish) |
| **Hierarchy / part-of** | Tree structure with sizes | Treemap, sunburst, icicle | Indented bar | Pie of pies |
| **Flow** | Source→sink pairs with magnitude | Sankey, alluvial | Chord diagram | Spaghetti line chart pretending to be flow |
| **Single headline value** | 1 number with target / trend | BAN tile with sparkline + delta | Bullet chart (target + actual + range) | Gauge / speedometer |

Two heuristics that follow from Cleveland & McGill:

1. If the audience must read or compare *specific values*, encode along position (bar
   chart, dot plot). If the audience must spot *patterns* (clusters, outliers, trend
   shapes), encode along position and color (scatterplot, heatmap, line).
2. If the dataset has many series and stacking would lose individual trajectories, go to
   small multiples before reaching for a multi-color line spaghetti.

---

## 6. Color theory — palettes that survive colorblindness and grayscale

Color choice is where well-meaning charts die. The three palette types map to three data
types and must not be confused.

**Sequential palettes.** A single-hue ramp from light (low) to dark (high). Use for
*continuous, ordered* data with no meaningful midpoint — population density, request
latency, temperature in summer, cumulative count. Examples: Viridis, Cividis, ColorBrewer's
Blues / Greens / Oranges / Purples / Reds.

**Diverging palettes.** Two contrasting hues meeting at a neutral midpoint (often white or
pale yellow). Use for *continuous data with a meaningful zero or reference* — temperature
anomaly (above/below normal), profit/loss, percent change, vote margin, correlation
coefficient. Examples: ColorBrewer's RdBu, BrBG, PiYG, PuOr, Spectral. The midpoint must
correspond to the reference value, not the dataset median, or the chart will lie.

**Categorical (qualitative) palettes.** Distinct, unordered hues for nominal categories
where no ordering is implied. Examples: ColorBrewer's Set1 / Set2 / Set3 / Paired / Dark2;
Tableau 10; Okabe-Ito (the "Wong" palette designed for protanopia and deuteranopia).
Hard limit: about six to eight distinct hues. If the data has more categories than that,
group, facet, or use a different encoding (small multiples, position).

### Colorblind-safe defaults

Approximately 8% of men and 0.5% of women have some form of color vision deficiency
(predominantly red-green protanopia and deuteranopia). The defaults below survive both
forms and grayscale printing:

- **Sequential continuous:** Viridis is the standard recommendation. It is *perceptually
  uniform* (equal data steps map to equal perceived color steps) and was designed to be
  colorblind-safe. Magma, Inferno, and Plasma are perceptually uniform alternatives with
  different aesthetics.
- **Sequential, optimized for colorblindness:** Cividis is specifically tuned for
  protanopia and deuteranopia at the cost of slightly reduced range for normal vision.
- **Diverging:** ColorBrewer RdBu with a white midpoint is the safest default. PuOr also
  survives well.
- **Categorical:** Okabe-Ito (8 colors), ColorBrewer Set2 / Dark2 (8 colors each), or
  Tableau 10. Never put red and green side-by-side as the only distinguishing pair.

### Operational rules

- **Never use color as the only encoding.** Pair color with shape, pattern, direct label,
  or position. WCAG 1.4.1 (Use of Color) explicitly requires this.
- **Grayscale test.** Convert the chart to grayscale and verify the encoded distinctions
  still read. Viridis and Cividis pass. The "rainbow" / jet colormap does not.
- **Avoid the rainbow colormap.** Jet / rainbow is not perceptually uniform — yellow
  bands appear lighter than blue/red at the same data value, distorting magnitude reads.
- **Use color to highlight one thing.** Default to grey for non-focal series, color for
  the focal series (Knaflic). This applies equally to dashboards: shared color encodes
  shared meaning (revenue is always green; cost is always grey).

---

## 7. Accessibility — WCAG for data visualizations

Charts are non-text content and fall under WCAG 2.x success criteria. The minimum
compliance bar for any chart shipped to a customer or external audience:

| Criterion | Level | What it means for a chart |
| --- | --- | --- |
| 1.1.1 Non-text Content | A | Every chart needs a text alternative. Short alt text states the chart type, what it shows, and the headline finding. A complex chart additionally needs a longer description (`aria-describedby` or a `<figcaption>`) and/or a linked data table. |
| 1.4.1 Use of Color | A | Color may not be the only way information is conveyed. Pair color with shape, pattern, direct label, or text. |
| 1.4.3 Contrast (Minimum) | AA | Text in/around the chart needs contrast ratio ≥ 4.5:1 against its background (≥ 3:1 for large text ≥ 18 pt). |
| 1.4.11 Non-text Contrast | AA | Graphical objects required to understand the content — bars, lines, data points, axis lines, focus indicators — need contrast ratio ≥ 3:1 against adjacent colors. |
| 2.1.1 Keyboard | A | Interactive charts must be operable from keyboard alone (Tab into chart, arrow keys between data points, Enter to drill). |
| 2.4.6 Headings and Labels | AA | Chart titles, axis labels, and legend entries must describe topic or purpose. |

### Alt-text pattern

A working template for chart alt text:

```
[chart type] of [what is measured] by [grouping]. [Headline finding].
```

Example: "Line chart of monthly revenue by region, January 2024 to April 2026. APAC
overtook EMEA in Q4 2025 and now leads by 18 percent."

For complex charts, follow the alt text with a longer textual description and provide the
underlying data table as a fallback. The data table is mandatory for screen reader users
who cannot perceive the visual encoding at all.

### Practical checklist

- Alt text on every `<img>` chart. For SVG charts, include `<title>` and `<desc>` and an
  `aria-labelledby` reference.
- Provide a data table near the chart (visible, hidden behind a "view as table" toggle, or
  available via download).
- Check contrast with an automated tool (axe DevTools, WAVE, Stark, Chrome DevTools
  Lighthouse). Don't trust perceptual judgment — the AA thresholds are not intuitive.
- Test with at least one screen reader (VoiceOver, NVDA, JAWS).
- Test with at least one colorblind simulator (Colorblindly, Sim Daltonism, Chrome
  emulator).

---

## 8. Grammar of graphics — the layered model

Leland Wilkinson's *The Grammar of Graphics* (1999) and Hadley Wickham's *A Layered
Grammar of Graphics* (2010) describe a chart as a *composition* of orthogonal components
rather than a named template (bar chart, line chart, etc.). The model underlies ggplot2,
Vega-Lite, Altair, Plotly Express, and Observable Plot.

A layered grammar chart has, at minimum:

- **Data.** The table of observations.
- **Aesthetic mapping.** Which columns map to which visual channels (x, y, color, size,
  shape, fill, alpha).
- **Geometric object (geom).** What is drawn at each row — point, line, bar, area, ribbon,
  polygon, text.
- **Statistical transformation (stat).** Optional aggregation applied before drawing —
  identity (no transform), binning (histogram), counting (bar of counts), smoothing
  (loess), boxplot summary, density.
- **Position adjustment.** How overlapping marks are arranged — identity, dodge (side by
  side), stack, fill (100% stack), jitter.
- **Scales.** Map data values to visual values — linear/log axes, color ramps, size
  scales, breaks and limits.
- **Coordinate system.** Cartesian, polar, geographic projection, flipped axes.
- **Faceting.** Split the plot into small multiples by one or two categorical variables.

Why this matters: the named-chart taxonomy (bar / pie / line) is a *cache* of common
grammar configurations, not a primitive. A "stacked bar" is just `geom = bar, position =
stack`. A "scatterplot with regression line" is two layers (`geom = point` plus `geom =
line` with `stat = smooth`). Thinking in grammar lets you compose unfamiliar charts and
debug familiar ones (a stacked bar with the wrong total is usually a position-adjustment
error).

---

## 9. Tooling map — picking a library or platform

The visualization ecosystem splits along three orthogonal axes: *audience* (analyst
notebook vs production dashboard vs publication figure), *interactivity* (static export
vs interactive web), and *control* (high-level templates vs low-level primitives).

### Python notebook and static

- **matplotlib** — the imperative foundation. Maximum control, verbose API, the standard
  for publication-grade static figures. Backs seaborn, pandas `.plot()`, and many
  scientific packages. Reach for it when surgical pixel control matters.
- **seaborn** — statistical convenience layer on matplotlib. Boxplots, violin plots,
  pair plots, regression plots in one line. Sensible defaults for color and theme.
- **altair** — Vega-Lite bindings in Python. Declarative grammar of graphics, JSON output,
  excellent for reproducibility and embedded web rendering. Smaller dataset support out of
  the box (with workarounds for ≥ 5000 rows).
- **plotly express / plotly.py** — high-level grammar-of-graphics API producing
  interactive HTML. Hover, zoom, pan built in. Easy chart-to-dashboard via Dash.
- **plotnine** — ggplot2 grammar in Python for users coming from R.
- **bokeh** — interactive HTML output, server-aware (Bokeh server pushes streaming
  updates). Strong for in-browser interactive dashboards backed by Python.

### R notebook and static

- **ggplot2** — the reference implementation of the layered grammar. The de facto standard
  for publication-grade static graphics in R. Use with patchwork (composition), gganimate
  (animation), ggrepel (label collision), ggdist (distributions), and the broader
  tidyverse.
- **lattice** — predecessor with strong small-multiples support.

### JavaScript interactive

- **D3.js** — low-level DOM-driven primitives. Maximum control, steep learning curve.
  Reach for it only when you need a custom visualization that off-the-shelf libraries
  cannot express. Most production charts that *look* like D3 are actually higher-level
  libraries (Plot, Vega) built on top.
- **Observable Plot** — high-level grammar of graphics built by the D3 team. Concise API
  for the 90% case; drop down to D3 for the remaining 10%.
- **Vega-Lite** — declarative JSON grammar. The same spec runs in the browser, in
  notebooks (via Altair), in Tableau Public's underlying engine, and server-side. Best
  for reproducibility and config-driven chart generation.
- **Plotly.js** — high-level chart types, WebGL acceleration for large datasets, built-in
  zoom/pan/hover. Easier than D3, more interactive than static plotters.
- **Chart.js** — light, simple, opinionated. Good for "I just need a clean line chart on
  a marketing page." Bad for anything novel.
- **ECharts** — Apache-licensed, full-featured chart library with strong defaults for
  enterprise dashboards.

### Dashboard and BI platforms

- **Tableau** — the long-time category leader for self-service BI. Strong visual defaults,
  drag-and-drop chart builder, calculated fields, large connector library. Expensive at
  scale; licensing is per-user.
- **Microsoft Power BI** — default for Microsoft-shop organizations. Tight integration
  with Azure, SQL Server, Office 365, Fabric. DAX as the analytical layer. Competitive
  pricing.
- **Google Looker Studio** — free, Google-ecosystem-native (GA4, Ads, Search Console,
  BigQuery). The first BI tool many marketing teams adopt.
- **Looker (enterprise)** — semantic-modeling-first (LookML), governance, embedded
  analytics. Different product from Looker Studio despite the shared brand.
- **Metabase** — open-source, easy for non-technical users via its "Question" interface.
  Strong default for startups with a small data team.
- **Apache Superset** — open-source, enterprise-grade, ships with 40+ database
  connectors. Higher operational overhead than Metabase but more flexible.
- **MongoDB Atlas Charts** — native to Atlas, queries directly against MongoDB
  collections, document-aware. The right default when the data already lives in Atlas;
  use the `mongodb-atlas-charts` skill for product-specific setup.
- **Sigma, Mode, Hex, Hightouch / Census, Preset** — newer entrants worth comparing for
  modern-data-stack deployments.

### Selection heuristic

1. **Audience reads the chart in a notebook or paper** → matplotlib (control), seaborn
   (stats), or ggplot2 (R).
2. **Audience reads the chart on a public web page** → plotly, observable plot, or
   vega-lite for declarative; D3 for fully custom.
3. **Audience refreshes the chart on a dashboard repeatedly** → BI tool (Tableau /
   Power BI / Looker Studio / Metabase / Superset / Atlas Charts), not a hand-coded
   chart.
4. **Chart needs to be embedded in a product surface** → plotly, ECharts, observable plot,
   or a React-charting library (Recharts, Visx, Nivo) bound to the product's design
   system.

---

## 10. Dashboard design — Few's at-a-glance rules

A dashboard is "a visual display of the most important information needed to achieve one
or more objectives that has been consolidated on a single screen so the information can be
monitored at a glance" (Few). The rules below distill Few's thirteen common mistakes into
operational design choices.

**Single screen, no scroll.** If the most important information requires scrolling, it is
not a dashboard — it is a report. Either reduce the content or split into multiple
dashboards by audience or use case.

**Inverted pyramid of detail.** Top of the screen carries the headline KPIs (BAN tiles).
Middle shows trend and composition (line charts, small multiples, stacked bars). Bottom
holds the granular detail tables and outlier lists. The eye scans top-down, left-right.

**Encode by importance.** The most important number gets the largest tile and the loudest
preattentive attribute (highest contrast, brightest color). The least important number is
grey, small, and tucked into a footer. Equal-size tiles for unequal-importance numbers is
the most common dashboard sin.

**Consistent encoding across panels.** If revenue is teal in one chart, it is teal in
every chart. If the dashboard is left-to-right time-ordered, every chart on it is
left-to-right time-ordered. Inconsistency forces the reader to re-learn the chart at every
panel.

**Context, always.** A number alone is meaningless. Pair it with a target, a prior-period
comparison, a confidence interval, or a sparkline. Bullet charts (Few's invention) encode
actual + target + range in a single horizontal bar and are the canonical KPI tile.

**No gratuitous decoration.** Drop the gauges, the speedometers, the lipstick. A 2D bar
chart with a sorted y-axis communicates more accurately than a 3D pie with a glossy
finish.

**Color carries meaning, not personality.** Use color to encode (status: green = on track,
red = at risk; series: revenue is always green) and to direct attention. Do not use color
to "make the dashboard prettier." Maintain accessibility (3:1 graphical contrast, never
color-alone).

**Drill-down paths, not detail-everywhere.** The top-level dashboard answers "is anything
wrong." The drill-down answers "what specifically." Avoid putting the row-level detail on
the top-level screen.

**Refresh cadence on the screen.** A dashboard that doesn't tell you when it last
refreshed lies by omission. Show "Updated 12 minutes ago" or "Live."

---

## 11. Anti-patterns — common ways to mislead

Most misleading charts are not malicious; they are produced by default settings or by
trying too hard to fit a story. The anti-patterns below cover roughly 90% of what an
auditor will flag.

**Truncated y-axis on a bar chart.** Bar charts encode magnitude as bar *length* from a
zero baseline. Truncating the y-axis turns a 1% change into a tall bar and lies by
exaggeration. Rule: bar-chart y-axes start at zero. (Line charts are exempt — a line
encodes change rate, not magnitude, so a non-zero baseline can be appropriate when the
range is small relative to the absolute value. State the baseline explicitly.)

**3D charts.** 3D pie charts distort the foreground slice. 3D bar charts foreshorten the
back row. 3D scatterplots make depth ambiguous. The 3D effect adds chartjunk while
materially degrading readability. Use 2D unless the third dimension is the data.

**Pie charts with more than 3 slices.** Angle and area are low-accuracy encodings
(Cleveland & McGill). When slices are similar in size, ordinary readers cannot reliably
rank them. Default to a horizontal bar chart sorted by value. The one place a pie chart
is acceptable is a 2-or-3-slice composition where one slice is dominant and the message
is "X is the majority."

**Dual y-axis charts.** Two y-axes with independent scales lets the chart author choose
where the lines cross, fabricating an apparent relationship that does not exist. The
reader cannot tell whether the lines move together because the data correlates or because
the scales were tuned. Rule: avoid dual-axis. If you must show two series on the same
x-axis, use either small multiples, a connected scatter, or normalize both series to a
common index (= 100 at t=0) and plot on one axis.

**Rainbow / jet colormap for continuous data.** The rainbow colormap is not perceptually
uniform — the yellow band reads brighter than the data warrants and the green-blue
transition reads as a discontinuity. Use Viridis, Cividis, Magma, Inferno, or Plasma
instead. This is the single most common misencoding in scientific publishing.

**Color as the only encoding.** A red/green status indicator that fails for the 8% of
colorblind viewers. Fix: pair color with shape, pattern, text label, or position. WCAG
1.4.1 makes this mandatory for accessible products.

**Spaghetti line charts.** 12 lines on one axis, all the same line weight, with a tiny
legend in the corner. Nobody can read it. Fix: small multiples, or pick the focal line
and grey out the rest.

**Choropleth of raw counts.** A map shaded by absolute count of something (cases,
revenue, users) is dominated by population. California is "dark" because California is
populous, not because anything interesting is happening there. Fix: choropleth of *rates*
(cases per 100k, revenue per capita), or use a cartogram.

**Cumulative-only line charts hiding non-cumulative shape.** A cumulative count always
trends up. The interesting signal is the *rate*. Show the daily / weekly delta alongside
the cumulative line.

**Mixing absolute and relative scales without labeling.** "Revenue up 2x" can mean from
$1M to $2M or from $100 to $200. Tag every chart axis with units and base period.

**Stacked bar charts where the audience must read individual segment values.** Stacking
preserves the total but destroys the comparison of inner segments because each inner
segment starts at a different baseline. If the audience needs to compare segments, use a
grouped bar or small multiples.

**Pie or donut with percentage labels totaling more than 100% or with overlapping
categories.** A surprisingly common bug. The chart is technically misleading because the
visual encoding (parts of a whole) does not match the data (overlapping categories).

**Bubble charts with size encoding diameter instead of area.** Doubling diameter
quadruples area, exaggerating large values. Encode bubble *area* proportional to value.
Most modern libraries do this correctly; verify the library's default before shipping.

**"Lying" map projections.** Mercator inflates polar areas (Greenland looks larger than
Africa). For global statistical maps prefer equal-area projections (Robinson, Equal
Earth, Mollweide).

**Animated charts hiding the comparison.** An animated bar-chart race is fun to watch
once. It is terrible for actually comparing values at specific time points. Show the
small-multiples version next to it.

---

## 12. Workflow — from data to finished chart

A repeatable workflow for producing an explanatory chart:

1. **Define the question.** Write the headline finding as a one-sentence claim before
   touching the chart tool. "Q3 revenue is down 12% YoY, driven by APAC." If you can't
   write the claim, you can't pick the chart yet — go back to exploratory mode.
2. **Identify the audience and medium.** A board deck slide, a regulator filing, a
   self-service dashboard, and a tweet require different polish, density, and
   interactivity budgets.
3. **Pick the intent.** Comparison, composition, distribution, relationship, temporal,
   or spatial. The intent narrows the chart options to roughly three candidates.
4. **Pick the chart type.** Apply Cleveland & McGill: encode the focal comparison along
   position. Avoid pie / 3D / dual-axis unless the data really requires it.
5. **Draft and declutter.** Render the chart with library defaults. Apply Knaflic's
   decluttering: remove border, soften gridlines, soften axis labels, push the legend
   onto the line if possible.
6. **Apply preattentive emphasis.** Grey out non-focal series. Color the focal series.
   Add direct labels at the line ends, not in a corner legend.
7. **Write the action title.** Replace the descriptive title with the headline claim.
8. **Verify graphical integrity.** Y-axes start at zero on bar charts. Pie slices sum to
   100%. No 3D. No dual axis. No truncation.
9. **Accessibility pass.** Contrast ratios. Alt text. Color-alone check. Grayscale test.
   Screen-reader title and description.
10. **Show it to someone who has not seen the data.** If the headline doesn't land
    within 5 seconds, the chart is not done. Iterate.

For exploratory work, compress steps 4-7 into one pass, skip steps 7-9, and skip step 10
unless the finding is novel.

---

## 13. References

### Foundational books

- Tufte, E. (1983, 2001). *The Visual Display of Quantitative Information* (2nd ed.).
  Graphics Press.
- Tufte, E. (1990). *Envisioning Information*. Graphics Press.
- Tufte, E. (1997). *Visual Explanations: Images and Quantities, Evidence and Narrative*.
  Graphics Press.
- Few, S. (2013). *Information Dashboard Design: Displaying Data for At-a-Glance
  Monitoring* (2nd ed.). Analytics Press.
- Few, S. (2009). *Now You See It: Simple Visualization Techniques for Quantitative
  Analysis*. Analytics Press.
- Knaflic, C. N. (2015). *Storytelling with Data: A Data Visualization Guide for Business
  Professionals*. Wiley.
- Cairo, A. (2016). *The Truthful Art: Data, Charts, and Maps for Communication*. New
  Riders.
- Cairo, A. (2019). *How Charts Lie: Getting Smarter about Visual Information*. W. W.
  Norton.
- Wilkinson, L. (2005). *The Grammar of Graphics* (2nd ed.). Springer.
- Wickham, H. (2016). *ggplot2: Elegant Graphics for Data Analysis* (3rd ed., online at
  ggplot2-book.org).
- Munzner, T. (2014). *Visualization Analysis and Design*. CRC Press.

### Foundational papers

- Cleveland, W. S., & McGill, R. (1984). *Graphical Perception: Theory, Experimentation,
  and Application to the Development of Graphical Methods.* JASA, 79(387).
- Wickham, H. (2010). *A Layered Grammar of Graphics.* Journal of Computational and
  Graphical Statistics, 19(1).
- Heer, J., & Bostock, M. (2010). *Crowdsourcing Graphical Perception.* CHI 2010.

### Color and palettes

- Brewer, C. A. ColorBrewer 2.0 — https://colorbrewer2.org
- Van der Walt, S., & Smith, N. (2015). Viridis colormap, designed for perceptual
  uniformity and colorblind safety. matplotlib documentation.
- Nuñez, J. R., Anderton, C. R., & Renslow, R. S. (2018). *Cividis colormap.* PLoS One.
- Wong, B. (2011). *Color blindness.* Nature Methods 8, 441 (Okabe-Ito palette).

### Accessibility

- W3C *Web Content Accessibility Guidelines (WCAG) 2.2.* — wcag.com / w3.org/TR/WCAG22
- W3C *Making images accessible.* (Charts as complex images.)
- Lundgard, A., & Satyanarayan, A. (2022). *Accessible Visualization via Natural Language
  Descriptions.* IEEE VIS 2022.

### Tool documentation

- matplotlib — matplotlib.org
- seaborn — seaborn.pydata.org
- ggplot2 — ggplot2.tidyverse.org
- altair — altair-viz.github.io
- plotly express — plotly.com/python/plotly-express
- bokeh — bokeh.org
- vega-lite — vega.github.io/vega-lite
- observable plot — observablehq.com/plot
- d3.js — d3js.org
- tableau — tableau.com
- power bi — powerbi.microsoft.com
- looker studio — lookerstudio.google.com
- metabase — metabase.com
- apache superset — superset.apache.org
- mongodb atlas charts — mongodb.com/products/charts

### Web references consulted

- Cleveland & McGill perceptual ranking — https://homepage.divms.uiowa.edu/~luke/classes/STAT4580/percep.html
- Tufte principles overview — https://thedoublethink.com/tuftes-principles-for-visualizing-quantitative-information/
- Storytelling with Data summary — https://www.shortform.com/blog/cole-nussbaumer-knaflic-storytelling-with-data/
- Abela chart chooser — https://datavizblog.com/2013/04/29/andrew-abelas-chart-chooser/
- Colorblind-friendly data visualization — https://colorblind.io/guides/data-visualization
- WCAG chart accessibility — https://www.a11y-collective.com/blog/accessible-charts/
- D3 vs Plotly vs Vega-Lite vs Chart.js — https://npm-compare.com/chart.js,d3,plotly.js,vega-lite
- BI tools comparison 2026 — https://valiotti.com/blog/top-5-data-visualization-tools/
- Misleading graphs Wikipedia — https://en.wikipedia.org/wiki/Misleading_graph
- Layered grammar of graphics — https://vita.had.co.nz/papers/layered-grammar.html
- Gestalt principles for data viz — https://emeeks.github.io/gestaltdataviz/section1.html
