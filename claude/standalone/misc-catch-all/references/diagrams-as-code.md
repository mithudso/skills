<!-- hub-reference-banner -->
> **Reference file — part of the `misc-catch-all` hub.** Formerly the standalone `diagrams-as-code` skill.
> Sibling topics in this family are now reference files under the hubs (`misc-catch-all`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: diagrams-as-code
version: "1.0.0"
updated: "2026-06-19"
description: >-
  TRIGGER: authoring or rendering text-to-diagram DSLs — Mermaid (flowchart/sequence/class/state/ER/gantt/gitgraph/mindmap/timeline/v11 architecture), Graphviz/DOT, PlantUML, D2; init directives, themes, securityLevel, ELK-vs-dagre layout, maxTextSize/maxEdges caps; headless rendering with @mermaid-js/mermaid-cli (mmdc), mermaid.render/parse JS API, puppeteerConfigFile/executablePath pitfalls, Kroki, mermaid-ink; SVG→PNG rasterization (sharp/resvg); rendering in CI/CD; large-graph "Too many edges"/"Maximum text size" failures; embedding limits in GitHub/markdown; prompting an LLM to emit valid Mermaid/DOT with a validate-and-repair loop. SKIP: .drawio/diagrams.net file format → drawio-diagrams (this skill only references it for conversion); chart-type selection and perception theory → da-8-data-visualization; mindmap STRUCTURE, interchange formats (OPML/XMind/.mm) and generation → mindmaps-and-concept-maps (this skill owns only Mermaid mindmap rendering).
---

Diagrams-as-code means committing a plain-text source (a DSL) and rendering it to SVG/PNG/PDF deterministically — diffable, reviewable, CI-renderable. This skill covers the four major DSLs, the unified Kroki server, the headless-browser rendering toolchain, the node/text caps that break large graphs, and LLM-driven diagram generation.

## Core concepts

### DSL selection

| DSL | Engine / runtime | Strengths | Layout engines | Output via |
|---|---|---|---|---|
| **Mermaid** | JS (browser/Node + headless Chrome) | Markdown-native (GitHub/GitLab render it inline), huge diagram-type catalog, easiest LLM target | dagre (default), **ELK** (`layout: elk`), tidy-tree, cose-bilkent | mmdc, mermaid.render, Kroki, mermaid-ink |
| **Graphviz/DOT** | Native C binary | Best automatic graph layout, mature, fast, scriptable | `dot` (layered/directed, default), `neato`/`fdp`/`sfdp` (force-directed), `circo` (circular), `twopi` (radial), `osage`/`patchwork` (clusters) | `dot -Tsvg`, Kroki |
| **PlantUML** | Java JAR (needs Graphviz, or bundled) | Rich UML (sequence/class/component/state), C4 via stdlib | Graphviz (default), **Smetana** (pure-Java port, no Graphviz dep), VizJs, ELK (`!pragma layout elk`) | `java -jar plantuml.jar`, PlantUML server, Kroki |
| **D2** (Terrastruct) | Go binary | Modern syntax, best-looking software-architecture diagrams, themes | **dagre** (default, bundled), **ELK** (bundled), **TALA** (paid, software-arch-tuned; `D2_LAYOUT=tala`) | `d2 in.d2 out.svg`, Kroki |

Rule of thumb: Mermaid for docs/READMEs and LLM output; Graphviz for large auto-laid-out graphs; PlantUML for formal UML; D2 for polished architecture diagrams.

### Mermaid essentials
- Diagram types: `flowchart`/`graph`, `sequenceDiagram`, `classDiagram`, `stateDiagram-v2`, `erDiagram`, `gantt`, `gitGraph`, `mindmap`, `timeline`, `journey`, `quadrantChart`, and **v11 `architecture-beta`** (cloud/service architecture with groups, services, junctions, icon packs).
- **Configuration precedence** (lowest→highest): default config → site config (`mermaid.initialize()` / `setSiteConfig()`) → **frontmatter** (`--- config: ... ---`, preferred, per-diagram) → **init directive** (`%%{init: {...}}%%`, JSON-only, **deprecated v10.5+** in favor of frontmatter). Frontmatter has highest priority and only affects its diagram.
- Key top-level keys: `theme` (`default|dark|forest|neutral|base`), `themeVariables`, `look` (`classic|handDrawn`), `layout` (`dagre|elk`), `fontFamily`, `logLevel`, `htmlLabels` (set globally now, not under `flowchart`), `startOnLoad`.
- **`securityLevel`**: `strict` (default — encodes tags, blocks click/JS), `loose` (allows interaction/HTML — XSS risk with untrusted input), `antiscript`, `sandbox` (renders in a sandboxed iframe, safest for third-party content). The `secure` array lists keys that directives/frontmatter cannot override (e.g. `securityLevel`, `maxTextSize`).
- **Node caps** (see large-graph section): `maxTextSize` (default 50000 chars) and `maxEdges` (default ~500; hard parser guard historically threw at 280/"Too many edges").
- ELK config example:
```yaml
---
config:
  layout: elk
  elk:
    mergeEdges: false
    nodePlacementStrategy: BRANDES_KOEPF   # SIMPLE | NETWORK_SIMPLEX | LINEAR_SEGMENTS | BRANDES_KOEPF
---
flowchart LR
  A --> B
```

### Graphviz/DOT essentials
- Pick the engine by structure: `dot` for hierarchy/DAGs, `neato`/`fdp`/`sfdp` for undirected/force-directed (sfdp scales to large graphs), `circo` circular, `twopi` radial. Override via `dot -Kfdp` / `-Tsvg`, or `layout=fdp` graph attribute.
- Attribute scopes: graph (`rankdir=LR`, `ranksep`, `overlap`, `splines`, `nodesep`, `mclimit` for dot), node (`shape`, `style`, `fillcolor`), edge (`label`, `samehead`/`sametail`). Many attributes are engine-specific (e.g. `overlap`/`sep`/`start` only apply to neato/fdp/sfdp/circo/twopi; `ranksep`/`rankdir` to dot/twopi).
- `dot -Tsvg in.dot -o out.svg` (also `-Tpng -Tpdf`). `overlap=prism`/`scale` removes node overlap in non-dot engines.

### PlantUML essentials
- Source between `@startuml`/`@enduml` (or `@startmindmap`, `@startgantt`, etc.). CLI: `java -jar plantuml.jar -tsvg file.puml` (`-tpng` default, `-tpdf`, `-ttxt` ASCII for sequence only).
- Default layout needs Graphviz installed; use **Smetana** (`!pragma layout smetana` or `-Playout=smetana`) for a zero-dependency pure-Java layout, or ELK (`!pragma layout elk`). The **LGPL** jar ships without embedded Graphviz.
- Validate before render: `plantuml -checkonly` (or `-checkmetadata`).
- Server URL scheme: `/plantuml/svg/<deflate+base64>` and `/plantuml/png/<encoded>`.

### Headless rendering toolchain
**mermaid-cli (mmdc)** drives headless Chrome via Puppeteer:
```sh
npm i -g @mermaid-js/mermaid-cli
mmdc -i input.mmd -o output.svg            # also .png, .pdf, .md (renders embedded fences)
mmdc -i in.mmd -o out.png -t dark -b transparent -c config.json
echo "graph TD; A-->B" | mmdc -i - -o out.svg   # stdin via heredoc/pipe
```
- **CI/Docker pitfall — "Failed to launch chrome … Running as root without --no-sandbox":** create `puppeteer-config.json` and pass it:
```json
{ "args": ["--no-sandbox", "--disable-setuid-sandbox", "--disable-dev-shm-usage"] }
```
```sh
mmdc -p puppeteer-config.json -i in.mmd -o out.svg
```
- **Use an already-installed Chromium** (avoid the bundled download): `--puppeteerConfigFile` accepts any `puppeteer.launch` option — `executablePath`, `product` (`chrome`|`firefox`), `timeout`, `headless`, `dumpio`:
```json
{ "executablePath": "/usr/bin/chromium", "args": ["--no-sandbox"] }
```
  Set `PUPPETEER_SKIP_DOWNLOAD=1` (older: `PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=1`) at install to skip the download; or `PUPPETEER_EXECUTABLE_PATH=...` env var. Note: Puppeteer is only guaranteed against its bundled Chromium version — mismatched browsers may fail. Official Docker images (`minlag/mermaid-cli`, `ghcr.io/mermaid-js/mermaid-cli`) already wire this up, looking for files in `/data`:
```sh
docker run --rm -u $(id -u):$(id -g) -v "$PWD":/data minlag/mermaid-cli -i diagram.mmd -o diagram.svg
```

**mermaid JS API** (in-browser/jsdom):
```js
import mermaid from 'mermaid';
mermaid.initialize({ startOnLoad: false, theme: 'dark', securityLevel: 'strict' });
const { svg, bindFunctions, diagramType } = await mermaid.render('id', 'graph TB\nA-->B');
element.innerHTML = svg; bindFunctions?.(element);
```
Prefer the high-level `mermaid` export over the deprecated `mermaidAPI`; `render()`/`parse()` are auto-queued (serial) to avoid race conditions. `mermaid.render` requires a DOM, so headless Node needs jsdom or a real browser.

**Kroki** — one HTTP API for ~30 DSLs (GraphViz, Mermaid, PlantUML/C4, D2, DBML, BlockDiag, Excalidraw, Vega, Structurizr, and **diagrams.net** experimentally — the conversion bridge to `.drawio`, owned by `drawio-diagrams`). POST `{diagram_source, diagram_type, output_format, diagram_options}` to `/`, or GET `/<type>/<format>/<deflate+base64>`; or POST plain text to `/<type>` with `Accept: image/svg+xml`. Self-host with `docker run yuzutech/kroki` (gateway bundles PlantUML/GraphViz/D2/etc.); **Mermaid, BPMN, Excalidraw, diagrams.net need companion containers** (`yuzutech/kroki-mermaid`, etc.) wired via `KROKI_MERMAID_HOST`. The public kroki.io makes no remote calls. **mermaid-ink** is a lighter hosted Mermaid-only PNG/SVG service (`https://mermaid.ink/img/<base64>`).

**SVG→PNG rasterization** (DSLs emit SVG; rasterize for embeds/email):
- **sharp** (libvips/librsvg): `sharp(svgBuffer, { density: 300 }).png().toFile(...)`. Increase `density` for crisp output (default 96 dpi); set `unlimited:true`/`limitInputPixels:false` for big SVGs. Known librsvg failures: `NoMemory` and dimension-dependent "SVG rendering failed" on large outputs.
- **resvg-js** (Rust, zero-dep, faster than sharp in benchmarks, custom-font support, WASM build): `new Resvg(svg, { fitTo: { mode:'width', value:1200 } }).render().asPng()`.
- node-canvas/svg2img also work but are slower. For LaTeX/print, prefer PDF output directly (Graphviz `-Tpdf`, PlantUML `-tpdf`).

### Node-cap & large-graph troubleshooting
- **"Maximum text size in diagram exceeded"** → source exceeds `maxTextSize` (50000 chars). Raise via `mermaid.initialize({ maxTextSize: 1000000 })` or config file passed to `mmdc -c config.json`. Keep it as low as feasible; large values let untrusted input hang the browser tab (the reason for the cap).
- **"Too many edges"** → exceeds `maxEdges` (default ~500; a hard `gk.length < 280` guard existed pre-fix). Set `{ "maxEdges": 1000 }` in the config and pass with `-c`. On mermaid.live, paste `{"maxTextSize":1000000,"maxEdges":1000}` into the Config tab and reload.
- For genuinely large graphs, switch Mermaid to `layout: elk`, or move to Graphviz `sfdp` (scalable force-directed) / `dot` with `overlap=prism`.
- **GitHub/markdown embedding limits:** GitHub renders Mermaid in fenced ```` ```mermaid ```` blocks but on its own (older) Mermaid version with its own caps; very large diagrams silently fail to render. GitLab defers/limits by source length and total rendered elements per page. Gitea caps via `MERMAID_MAX_SOURCE_CHARACTERS` (default raised to 50000). When a diagram won't render in a host, render it offline to SVG/PNG and embed the image instead.

### LLM-to-diagram generation
LLMs are strong at emitting Mermaid/DOT but routinely produce **syntax-invalid** output (bracket/shape mismatches like `A["Check Status?"]}`, where a `?` lures a lightweight model into a `}` decision-node close while the line opened with `[`; invalid `classDef`/annotation syntax; bad arrow/relationship tokens). Build a **generate→validate→repair loop**:
1. Prompt the model with a tight grammar reminder and few-shot valid examples (a "comment-first" protocol — declaring node type+symbol in a `%%` comment before each node — measurably cuts error rates).
2. **Validate** without rendering: `await mermaid.parse(src, { suppressErrors: true })` returns `false` (or `{ diagramType }` if valid) instead of throwing — cheaper than `render()`. For other DSLs use `plantuml -checkonly`, `dot` exit code, or `d2` parse. mmdc compilation is the heaviest validator.
3. **Repair:** feed the exact parser error message (`Parse error on line N: … Expecting 'NEWLINE','SQS',…`) back to the model and retry, capped at ~3 attempts to avoid infinite loops (PocketFlow self-healing-mermaid, GenAIScript `system.diagrams` repairer pattern).
4. Cheap client-side regex pre-checks (balanced brackets, known diagram-type keyword) catch most errors before invoking an expensive agentic repair pass.

## Sources

- Mermaid docs — Configuration & Directives: https://mermaid.js.org/config/directives.html
- Mermaid docs — Layout engines (ELK vs dagre): https://mermaid-js-mermaid.mintlify.app/configuration/layouts
- Mermaid docs — rendering options / mermaid.render API: https://www.mintlify.com/mermaid-js/mermaid/advanced/rendering-options
- Mermaid docs — mermaid.parse() & error handling: https://www.mintlify.com/mermaid-js/mermaid/advanced/error-handling
- mermaid-cli — Linux sandbox issue (puppeteer-config.json --no-sandbox): https://github.com/mermaid-js/mermaid-cli/blob/HEAD/docs/linux-sandbox-issue.md
- mermaid-cli — use already-installed Chromium (executablePath): https://github.com/mermaid-js/mermaid-cli/blob/HEAD/docs/already-installed-chromium.md
- Graphviz Layout Engines: https://graphviz.org/docs/layouts/
- Graphviz Attributes reference: https://graphviz.org/doc/info/attrs.html
- PlantUML command line: https://plantuml.com/command-line
- D2 layout engines: https://d2lang.com/tour/layouts/
- Kroki usage & supported diagram types: https://docs.kroki.io/kroki/setup/usage/
- Kroki install (self-host, companion containers): https://docs.kroki.io/kroki/setup/install/
- mermaid #5042: Too many edges / flowchart.maxEdges: https://github.com/mermaid-js/mermaid/issues/5042
- resvg-js — SVG→PNG (Rust, zero-dep): https://github.com/thx/resvg-js
- GenAIScript: Mermaids Unbroken (LLM diagram auto-repair): https://microsoft.github.io/genaiscript/blog/mermaids/
- PocketFlow self-healing-mermaid cookbook: https://github.com/The-Pocket/PocketFlow/tree/main/cookbook/pocketflow-self-healing-mermaid
