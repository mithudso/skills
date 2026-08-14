<!-- hub-reference-banner -->
> **Reference file — part of the `misc-catch-all` hub.** Formerly the standalone `programmatic-charting-libraries` skill.
> Sibling topics in this family are now reference files under the hubs (`misc-catch-all`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: programmatic-charting-libraries
version: "1.0.0"
updated: "2026-06-19"
description: >-
  TRIGGER: implementing or rendering a chart in code — D3.js (selections, data join, scales, axes, line/area/arc/stack generators, transitions); declarative grammars (Vega, Vega-Lite JSON specs, marks/encodings/transforms/params, Observable Plot, ggplot2-mapping context); high-level JS libs (Plotly.js/Plotly Express, Apache ECharts, Chart.js) and when each fits; SVG-vs-Canvas-vs-WebGL substrate choice by data size; responsive/interactive charts; or SERVER-SIDE/headless/static rendering in a pipeline (vega-cli vg2svg/vl2png, node-canvas + Chart.js / chartjs-node-canvas, Plotly kaleido write_image, QuickChart). SKIP: which chart to use, perceptual/color theory, grammar-of-graphics THEORY, dashboard design → da-8-data-visualization. SKIP: graph/network layout algorithms → graph-network-3d-visualization.
---

# Programmatic Charting Libraries — Hands-On Implementation & Rendering

This skill is the **"how do I implement/render it in code"** layer. It assumes the *which chart and why* question is already answered — for chart selection by intent, perceptual/color theory, and grammar-of-graphics theory, **cite da-8-data-visualization**. For graph/network layout (force-directed, hierarchical, community layout), **cite graph-network-3d-visualization** — do not reproduce that here.

The charting stack spans three abstraction levels, plus an orthogonal rendering-substrate axis and a server-side/headless axis:

- **Low-level toolkit** — D3.js: maximum control, you build everything (scales, axes, shapes, DOM).
- **Declarative grammars** — Vega / Vega-Lite (JSON), Observable Plot (concise JS): describe *what*, the library figures out *how*. (ggplot2 is the R reference grammar.)
- **High-level chart libraries** — Plotly.js, Apache ECharts, Chart.js: pass config, get a finished interactive chart.

## Core concepts

### Library decision table (by use case)

| Situation | Pick | Why |
|---|---|---|
| Standard line/bar/pie/scatter, < ~10k points, fast ship | **Chart.js** | Tiny (~70 kB gz), simplest API, Canvas, mature docs, official React/Vue/Angular wrappers |
| Exploratory / declarative / multi-view, reproducible specs | **Vega-Lite** | JSON grammar; auto axes/legends/scales; layer/facet/concat/repeat; selections for interaction |
| Concise exploratory JS charts in notebooks/web | **Observable Plot** | D3-team grammar-of-marks; one-liners; sane defaults; returns an SVG element |
| Scientific/statistical/3D, interaction out-of-the-box | **Plotly.js / Plotly Express** | 3D, box/violin/contour, `scattergl` WebGL; zoom/pan/hover built in; Python/R/Julia bindings share the same JS core |
| Big data (10k–1M+ pts), rich chart catalog, dashboards | **Apache ECharts** | Canvas (or SVG), `series.sampling:'lttb'` downsampling, dataZoom/brush/toolbox, Sankey/treemap/heatmap/radar; Apache-2.0; node-canvas SSR |
| Fully bespoke, novel, non-standard visualization | **D3.js** | Total control; the substrate other libs are built on; highest effort |
| Real-time monitoring, 100k+ streaming points | **ECharts** (LTTB) or **Plotly `scattergl`** or raw **regl/deck.gl** (WebGL) | GPU rendering; element count stops mattering |

Sizing rule of thumb (size for your **95th-percentile real tenant, not demo data**): **< ~5k pts/chart → SVG libs are fine; 5k–100k → Canvas (Chart.js, ECharts); 100k+ → WebGL** (ECharts WebGL, Plotly `scattergl`, deck.gl).

**Why D3 is low-level:** D3 is not a charting library — it is a toolkit for binding data to the DOM and for building the *pieces* of a chart (scales, axes, shape generators). You write the rendering loop yourself. That is the source of both its unlimited flexibility and its cost; reach for it only when a higher-level library cannot express the visualization.

### D3 essentials — the data join, scales, generators

**The data join (enter / update / exit).** `selection.data(data, key?)` binds an array to selected elements and partitions them into three sub-selections: **enter** (data with no element yet — placeholders), **update** (elements that matched data), **exit** (elements whose data is gone). `selection.join()` is the modern convenience that replaces the manual enter/append/merge/exit/remove pattern. Always call `selectAll` *before* the join so entering elements have a parent. Provide a **key function** to keep object constancy across updates (and to minimize DOM churn / animate correctly).

```js
const x = d3.scaleBand().domain(data.map(d => d.name))
  .range([0, width]).padding(0.1);
const y = d3.scaleLinear().domain([0, d3.max(data, d => d.value)]).nice()
  .range([height, 0]);                       // note: range is flipped (SVG y grows down)

svg.append("g").attr("transform", `translate(0,${height})`)
  .call(d3.axisBottom(x));                    // axis generator renders ticks+labels
svg.append("g").call(d3.axisLeft(y));

svg.selectAll("rect")
  .data(data, d => d.name)                    // key function = object constancy
  .join(
    enter => enter.append("rect")
      .attr("x", d => x(d.name)).attr("width", x.bandwidth())
      .attr("y", height).attr("height", 0)    // start collapsed, then grow
      .call(e => e.transition().duration(750)
        .attr("y", d => y(d.value))
        .attr("height", d => height - y(d.value))),
    update => update.call(u => u.transition().duration(750)
        .attr("y", d => y(d.value))
        .attr("height", d => height - y(d.value))),
    exit => exit.call(x => x.transition().duration(750)
        .attr("height", 0).attr("y", height).remove())
  );
```

**Scales** map *data space → visual space*: `scaleLinear`, `scaleBand` (categorical bars, has `.bandwidth()`), `scaleOrdinal` (categories → colors), `scaleTime` (dates), `scaleLog`, `scaleSqrt`/`scalePow` (areas/radii), `scaleQuantize`/`scaleThreshold` (binning). `.nice()` rounds the domain; `.domain()`/`.range()` set in/out.

**Shape generators** turn data into SVG path strings: `d3.line()`, `d3.area()`, `d3.arc()` (+ `d3.pie()` for slices), `d3.stack()` (stacked bars/areas), `d3.curveX` interpolators. They return functions you assign to a path's `d` attribute. **Transitions** (`selection.transition().duration()`) interpolate attributes (default easing cubic-in-out); create them *inside* the enter/update/exit join functions to animate each phase.

### Vega-Lite — declarative grammar-of-graphics in JSON

A Vega-Lite spec is a JSON object mapping **data → marks → encoding channels**, with optional **transforms** and **params** (variables/selections). The compiler infers axes, legends, and scales automatically (overridable). Vega-Lite compiles down to **Vega** (a lower-level grammar with explicit signals/scales/marks); Vega in turn is rendered by its runtime. Marks: `bar, line, area, point, circle, square, tick, rule, text, geoshape, arc`. Channels: `x, y, color, size, shape, opacity, theta, tooltip, …` — each field def gives a `field` + `type` (`quantitative, nominal, ordinal, temporal`).

```json
{
  "$schema": "https://vega.github.io/schema/vega-lite/v5.json",
  "data": {"url": "data/seattle-weather.csv"},
  "transform": [{"filter": "datum.temp_max > 0"}],
  "params": [{"name": "brush", "select": {"type": "interval", "encodings": ["x"]}}],
  "mark": "point",
  "encoding": {
    "x": {"field": "date", "type": "temporal", "timeUnit": "yearmonth"},
    "y": {"field": "temp_max", "type": "quantitative", "aggregate": "mean"},
    "color": {"field": "weather", "type": "nominal"},
    "opacity": {"condition": {"param": "brush", "value": 1}, "value": 0.2}
  }
}
```

**Params** are the interaction grammar: a plain variable, or a **selection** (`point`/`interval`) bound to mouse input or to an input widget (slider/dropdown). Reference a selection in `condition`, `filter`, or a scale `domain` (binding a brush to a scale domain — `"scale": {"domain": {"param": "brush"}}` — gives superior interactive performance vs. filtering). Multi-view composition: `layer`, `facet`, `concat`, `repeat`. **ggplot2 context (R):** `ggplot(data) + aes(x, y, color) + geom_*()` is the same grammar — `aes()` ≈ encoding channels, `geom_*` ≈ marks, `stat_*`/`scale_*`/`facet_*` ≈ transforms/scales/facet; map this mental model onto Vega-Lite when porting.

**Observable Plot** is the concise JS cousin (same D3 team): you compose **marks** (`Plot.lineY`, `Plot.dot`, `Plot.barY`, `Plot.areaY`, `Plot.ruleY`) instead of picking a chart "type", and call `.plot()` (returns an SVG/HTML element). `Plot.lineY(aapl, {x: "Date", y: "Close", stroke: "Symbol"})` auto-derives scales/axes/legend. Shorthand: `Plot.lineY(numbers).plot()`.

### High-level JS libraries — when each fits

- **Chart.js** — config-object API (`{type, data, datasets, options}`), Canvas-rendered, smallest bundle, plugin ecosystem, great defaults. Best for standard dashboards and quick wins; weaker for exotic chart types and very large data without plugins.
- **Apache ECharts** — broadest catalog (Sankey, treemap, graph, heatmap, radar, candlestick, 3D), Canvas **and** SVG renderers, `dataZoom`/`brush`/`toolbox`, LTTB downsampling, node-canvas SSR. The Canvas-rendered "escape hatch" for big data and dense dashboards. Note v5 gotcha: `title:{text:'…'}` (object), not a string.
- **Plotly.js** — declarative `data` traces + `layout`; rich scientific/statistical/3D charts; zoom/pan/hover/select **built in**; `scattergl`/`scatter3d` WebGL traces for scale. Larger bundle. Plotly **Express** (`px.line`, `px.bar`, `px.scatter`) is the one-line Python API that emits the same figure JSON the JS core renders — so a Python analyst and a JS app share one rendering engine.

### Rendering substrates — SVG vs Canvas vs WebGL

| | SVG | Canvas | WebGL |
|---|---|---|---|
| Model | DOM nodes (retained) | single bitmap (immediate) | GPU bitmap (immediate) |
| Hit-testing/interaction | native, free | manual (`isPointInPath`, dirty-rects) | manual + GPU picking |
| Accessibility / CSS / styling | excellent | none built-in | none built-in |
| Server-side render | yes (just text) | via node-canvas | hard |
| Sweet spot | **< ~1–5k elements** | **~5k–100k+** | **100k–millions** |

Indicative draw times (one source's bench): at 10k points SVG ~40 ms / Canvas ~6 ms / WebGL ~0.8 ms; at 1M points SVG unusable / Canvas ~200 ms+ / WebGL ~8 ms. Caveats: SVG's cost is DOM node count and layout/paint, so it shines for **text-heavy, few-element, interactive, accessible** charts; Canvas wins on element count but you own hit-testing and state; WebGL wins on raw throughput but **text is the pain point** (render as textures or rebuild from primitives) and you must build the pipeline (buffers, shaders, draw-call batching, `OffscreenCanvas`). Many production systems **hybridize**: WebGL for the dense scene + Canvas/SVG/DOM for labels and UI. **Responsive/interactive:** use a `ResizeObserver` (or the library's `resize()`/`autosize`) to re-layout; debounce redraws; for SVG use `viewBox` + `preserveAspectRatio`; downsample (LTTB) before you switch substrate.

### Server-side / headless / static rendering

For pipelines, emails, PDFs, SMS, chatbots, CI artifacts, or reports — render charts to **PNG/SVG/PDF without a browser tab**.

**Vega / Vega-Lite CLI** (`npm i -g vega-cli vega-lite`, or run via `npx`). Vega specs: `vg2svg`, `vg2png`, `vg2pdf`. Vega-Lite specs: `vl2vg` (compile to Vega), `vl2svg`, `vl2png`, `vl2pdf`.
```bash
npx -p vega -p vega-lite vl2png chart.vl.json chart.png      # Vega-Lite → PNG
npx -p vega -p vega-lite vl2svg chart.vl.json > chart.svg
vg2svg -h spec.vg.json out.svg                                # -h adds XML header
cat spec.vl.json | vl2svg > out.svg                           # stdin → stdout, pipeable
```
PNG output and accurate font metrics require the **node-canvas** package (`vega-cli` bundles it; the bare `vega` package does not — add `canvas` yourself or you get `CanvasRenderer is missing a valid canvas`).

**node-canvas + Chart.js** via `chartjs-node-canvas` — render a Chart.js config to a Buffer/stream/data-URL with no DOM:
```js
const { ChartJSNodeCanvas } = require('chartjs-node-canvas');
const canvas = new ChartJSNodeCanvas({ width: 800, height: 400, backgroundColour: 'white' });
const buffer = await canvas.renderToBuffer({ type: 'bar', data: {...}, options: {...} }); // PNG by default
// SVG/PDF require the *Sync* API: renderToBufferSync(cfg, 'image/svg+xml' | 'application/pdf')
```

**Plotly static export (Kaleido)** — `pip install --upgrade "plotly[kaleido]"`:
```python
import plotly.express as px
fig = px.bar(px.data.gapminder().query("country=='Canada'"), x="year", y="pop")
fig.write_image("fig.png")            # format inferred from extension; png/jpg/webp/svg/pdf
fig.write_image("fig.svg")            # vector
img = fig.to_image(format="png")      # bytes, e.g. for an HTTP response
```
Key facts: Kaleido **v1 uses your installed Chrome/Chromium** (it no longer bundles Chrome); requires plotly.py ≥ 6.1.1. The legacy `engine=` arg and Orca are **deprecated/removed after Sept 2025** — drop `engine="orca"`/`"kaleido"`. For many figures use `plotly.io.write_images([...])` (faster than a loop).

**QuickChart** — hosted (or self-hostable) web service that renders a **Chart.js config** to an image via URL or POST; ideal for email/SMS/chatbots that can't run JS:
```
https://quickchart.io/chart?c={type:'bar',data:{labels:[2012,2013],datasets:[{data:[120,60]}]}}&width=500&height=300&format=png&version=4
```
Params: `chart`/`c`, `width`, `height`, `devicePixelRatio` (2=retina), `backgroundColor`, `version` (set `4` for Chart.js v4), `format` (png/webp/svg/pdf), `encoding` (url/base64). URL-encode the config for GET; for untrusted clients sign with HMAC (`sig`+`accountId`) instead of exposing the API key. **Security:** Chart.js configs can carry arbitrary JS — sandbox any self-hosted QuickChart and never trust client-supplied configs.

**Pipeline pattern:** generate a declarative spec (Vega-Lite JSON or Chart.js config) in your data job → shell out to `vl2png`/`vl2svg`, `chartjs-node-canvas`, or `kaleido` → attach the PNG/SVG to a report, email, or PDF. Prefer SVG for crisp/scalable/diff-able output and PNG when the consumer can't render SVG.

## Sources

- Joining data — D3 selection (data/join/enter/exit): https://d3js.org/d3-selection/joining
- How Selections Work — Mike Bostock: https://bost.ocks.org/mike/selection/
- Vega-Lite View Specification: https://vega.github.io/vega-lite/docs/spec.html
- Encoding — Vega-Lite (channels, scales): https://vega.github.io/vega-lite/docs/encoding.html
- Dynamic Behaviors with Parameters — Vega-Lite: https://vega.github.io/vega-lite/docs/parameter.html
- Marks / Plots — Observable Plot: https://observablehq.com/plot/features/marks
- Usage / Command Line Utilities — Vega (vg2svg/png, node-canvas): https://vega.github.io/vega/usage/
- Compiling Vega-Lite to Vega (vl2vg/png/svg/pdf): https://vega.github.io/vega-lite/usage/compile.html
- Static image export in Python — Plotly: https://plotly.com/python/static-image-export/
- plotly/Kaleido (headless static export, Chrome-based v1): https://github.com/plotly/kaleido
- QuickChart API parameters: https://quickchart.io/documentation/usage/parameters/
- Canvas vs SVG vs WebGL benchmark & thresholds: https://chartts.com/blog/canvas-vs-svg-vs-webgl
- The Best JavaScript Chart Libraries (decision table): https://www.usedatabrain.com/blog/javascript-chart-libraries
- SeanSobey/ChartjsNodeCanvas (server-side Chart.js render): https://github.com/SeanSobey/ChartjsNodeCanvas
