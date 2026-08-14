<!-- hub-reference-banner -->
> **Reference file — part of the `misc-catch-all` hub.** Formerly the standalone `mindmaps-and-concept-maps` skill.
> Sibling topics in this family are now reference files under the hubs (`misc-catch-all`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: mindmaps-and-concept-maps
version: "1.0.0"
updated: "2026-06-19"
description: >-
  Structural and programmatic reference for hierarchical knowledge diagrams: Buzan radial single-root mind maps vs. Novak labeled-edge concept maps (propositions, cross-links) vs. argument maps, plus generation and interchange. TRIGGER: choosing mind map vs concept map vs argument map; generating a mindmap from markdown/JSON/outline/AST (Markmap markmap-cli/markmap-lib/markmap-view, Mermaid mindmap, jsMind, GoJS); interchange formats (OPML 2.0, FreeMind/Freeplane .mm, XMind zip, CmapTools CXL) and converting between them; radial/tidy-tree (Reingold-Tilford/Buchheim) layout and left/right balancing; tool landscape (XMind, Freeplane, Obsidian/Logseq plugins). SKIP: 3D/WebGL/force-directed graph rendering -> graph-network-3d-visualization; concept-mapping pedagogy/learning-theory depth -> instructional-design-course-architecture; PKM/note-taking systems -> research-note-taking-pkm; headless rendering of Mermaid mindmaps to PNG/SVG -> diagrams-as-code.
---

# Mindmaps & Concept Maps: Programmatic Knowledge Visualization

A practical reference for *generating and exchanging* hierarchical knowledge diagrams in code — not for hand-drawing them in a GUI. The three families below look similar but have different formal structures (tree vs. DAG vs. argument graph), which dictates which file format, which generator, and which layout algorithm you can use. Pick the structure first; the toolchain follows.

Cross-references: 3D / WebGL / force-directed graph rendering is owned by **graph-network-3d-visualization**. Concept-mapping *pedagogy* / learning theory is owned by **instructional-design-course-architecture**. PKM / note-taking systems are owned by **research-note-taking-pkm**.

## Core concepts

### 1. Mind map vs. concept map vs. argument map (decision)

These differ in their **graph topology** and in **whether edges carry meaning**. Get this right before choosing a format.

**Mind map (Buzan, 1974).** A *rooted tree* with exactly **one central node** and a single parent per node; branches radiate outward, organic/colorful, often one keyword per branch. Edges are **unlabeled** — connection means "is associated with / belongs under," nothing more. There are no cross-links, so a mind map cannot express many-to-many relationships natively (NN/g: "no definition of relationships… all edges are unlabeled"). Best for: *expanding a single topic*, brainstorming, divergent capture. Substrate: any tree format (OPML, `.mm`, markmap markdown, Mermaid `mindmap`, jsMind JSON).

**Concept map (Novak & Cañas, IHMC, 1972).** A *directed graph* (often a DAG, not a pure tree): nodes are concepts, **edges are labeled with linking phrases** (verbs/prepositions: "leads to," "results from," "is part of"). Two concepts + a linking phrase form a **proposition** — a readable, verifiable sentence ("Topic A *causes* Topic B"). A node may have **multiple parents** and **cross-links** between branches, so it expresses many-to-many relationships. Usually organized top-down from a **focus question**. Best for: *explaining relationships among several concepts*, assessment of understanding, knowledge models. Substrate: CmapTools **CXL** (the only mainstream format that natively models labeled edges as first-class `connection`/`linking-phrase` elements).

**Argument map.** A directed graph specialized for *reasoning*: nodes are claims/premises/objections, edges are **support/rebut/undercut** relations. It is a concept map whose edge vocabulary is fixed to argumentation roles. Best for: debate structure, critical-thinking analysis, decision rationale. (Davies 2011 surveys all three.)

Quick rule: **one root + unlabeled edges → mind map; labeled edges + cross-links → concept map; support/attack edges → argument map.** If you only have a tree, a mind map format suffices and every generator below works. If edges must carry meaning, you need CXL (or a general graph format) and most mindmap generators will *lose* your edge labels on round-trip.

### 2. Programmatic generation toolchains

**Markmap (markdown → mindmap).** The dominant markdown-driven toolchain. Three packages:
- `markmap-lib` — `Transformer.transform(markdown)` parses markdown into a `{ root, features }` node tree plus an `assets` object (CSS/JS to inject). Headings and nested bullet lists become the hierarchy.
- `markmap-view` — renders the transformed `{ root }` into an interactive **SVG** via `Markmap.create('#svg', options, root)` (depends on D3 v7). Supports `deriveOptions(jsonOptions)` for a portable JSON subset of options.
- `markmap-cli` — terminal wrapper. Quick start:
  ```
  npx markmap-cli markmap.md            # generate + open HTML
  npx markmap-cli --offline -o out.html markmap.md   # inline all assets
  markmap --no-open --no-toolbar -w notes.md         # watch mode (dev)
  ```
  Also surfaces via VSCode extension, `coc-markmap` (Vim/Neovim), and an autoloader script. Pipeline: `markmap-lib` (markdown→data) → `markmap-render`/`markmap-view` (data→SVG HTML).

**Mermaid `mindmap`.** Indentation-driven, embeddable in any Mermaid host (GitHub, docs sites). Syntax (still flagged **experimental**; icon integration unstable):
```
mindmap
  Root
    A
      B
      C
```
Indentation depth (relative, not absolute) sets parent/child; on ambiguous indentation Mermaid walks back to "the first node with smaller indentation" and treats the node as that node's child. Node shapes: square `[]`, rounded `()`, circle `(())`, bang `))((`, cloud `)`, hexagon `{{}}`. Supports `::icon()` and `:::class`, markdown strings (`**bold**`, auto-wrap), and a registerable **Tidy Tree** layout. Practical limits: it is **single-root** (true mind map, no labeled edges, no cross-links), and very large maps degrade — for big trees prefer markmap or jsMind. (For headless rendering of Mermaid mindmaps to PNG/SVG see **diagrams-as-code**.)

**jsMind.** Pure-JS editor/renderer on HTML5 **canvas or SVG** (`view.engine: 'canvas' | 'svg'`; switch to SVG for large maps for big perf gains). Loads/exports **three data formats**: `node_tree` (nested), `node_array` (flat with `parentid`), and `freemind` (embedded `.mm` XML string) — making it a useful in-browser converter. Node carries `id`, `topic`, `direction (left|center|right)`, `expanded`, `children`. `jm.get_data(format)` exports in any of the three.

**GoJS / mxGraph (draw.io).** Heavier commercial/diagramming libraries with mindmap samples; use when you need a full editable diagramming surface with custom layouts, undo, and serialization rather than a read-mostly render. (See **drawio-diagrams** for draw.io mxGraph XML.)

**Generating from JSON / outline / AST.** The common pattern: produce a nested `{ id, topic|text, children:[] }` tree, then feed it to the renderer. From markdown, parse to an AST (remark/markdown-it) and map heading/list depth to tree depth (this is exactly what markmap-lib does). From an outline, parse indentation. From arbitrary data, emit the renderer's native JSON (jsMind `node_tree`, or markmap's root node).

### 3. Interchange formats (and a minimal OPML example)

**OPML 2.0 (Outline Processor Markup Language).** XML; an outline is "a tree where each node has named string attributes." Root `<opml version="2.0">` with one `<head>` (optional `title`, `dateCreated`, `ownerName`, `expansionState`, window geometry) and one `<body>` of nested `<outline>` elements. **Every `<outline>` requires a `text` attribute**; optional `type`, `_note`, `isComment`, `created`, `category`. `type="link"`/`type="include"` reference external URLs/OPML. Lingua franca for outliners and feed readers; most mindmap tools import/export it. Minimal example:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<opml version="2.0">
  <head><title>Demo</title></head>
  <body>
    <outline text="Root">
      <outline text="A">
        <outline text="B"/>
        <outline text="C"/>
      </outline>
    </outline>
  </body>
</opml>
```
Note: OPML is a pure tree with unlabeled parent/child links — it maps cleanly to *mind maps* but cannot represent concept-map labeled edges.

**FreeMind / Freeplane `.mm`.** XML; root `<map version="...">` containing one root `<node>`, recursively nested `<node>` children. Key `<node>` attributes: `TEXT`, `ID`, `LINK`, `FOLDED`, `POSITION` (`left|right` — only meaningful for direct children of root, this is how left/right balancing is *persisted*), `COLOR`, `CREATED`/`MODIFIED`. Child elements: `<edge>`, `<font>`, `<icon>`, `<cloud>`, `<richcontent>` (HTML notes), and `<arrowlink destination="...">` — the **arrowlink is the one place a `.mm` file can express a cross-link / non-tree edge** (graphical link to another node by ID). Freeplane extends the FreeMind schema (richer styles, `STYLE`, format patterns). FreeMind writes no XML declaration.

**XMind.** A **ZIP archive** with `.xmind` extension. Modern XMind (Zen/2020+) uses `content.json` (a `[{ sheet → rootTopic → children.attached[] }]` tree); legacy XMind 8 uses `content.xml` (`<xmap-content>` → `<sheet>` → `<topic>` → `<children><topics type="attached">`). Companions: `styles.xml`/styling, `comments.xml`, `META-INF/manifest.xml`, optional `Attachments/`, `Revisions/`, `Thumbnails/`. Topics carry `id`, `title`, `children` (`attached`/`detached`/`summary`), `markers`, `labels`, `notes` (dual plain + XHTML); `<relationship>`s are stored at *sheet* level (not inside topics) — XMind's cross-link mechanism. Python SDKs (`xmindltd/xmind-sdk-python`, `zhuifengshen/xmind`) and a Go `xmind-mcp` parse the archive into an object tree; preserve unknown keys via an `extra`/passthrough field for safe round-trips.

**CmapTools CXL.** IHMC's XML schema (`cmap.xsd`) for *concept maps* — the format that models **labeled edges as first-class data**: separate `<concept>`, `<linking-phrase>`, and `<connection from=… to=…>` elements, plus a focus question and appearance lists. Use CXL (not OPML/`.mm`/XMind) when propositions and many-to-many cross-links must survive.

**Converting between them.** OPML is the usual hub because it is the simplest common tree:
- `pandoc --from=opml --to=markdown_mmd file.opml` (OPML ↔ markdown; markdown then feeds markmap/Mermaid).
- FreeMind ships `mm2opml.xsl` (XSLT) to emit OPML; community scripts (`adxsoft/OPMLtoMM`) go OPML→`.mm`, mapping each `<outline>` to a `<node>` and `_note` to `<richcontent>`.
- XMind and Freeplane both import OPML and export `.mm`, so a common chain is **XMind → FreeMind `.mm` → OPML (XSLT) → markdown (pandoc)**.
- Lossy-conversion warning: tree formats (OPML, `.mm` without arrowlinks) **drop concept-map edge labels and cross-links**. Round-trip concept maps only through CXL or a labeled-graph format.

### 4. Layout algorithms

The visual heart of a mind/concept map renderer is a **tree (or graph) layout** that assigns coordinates.

**Reingold–Tilford "tidy tree."** The standard tidy node-link layout: lay out each subtree independently, then push sibling subtrees together to minimal non-overlapping separation via a postorder contour walk; a preorder pass converts relative offsets to absolute coordinates. Guarantees isomorphic subtrees draw identically and parents center over children. The naive form is O(n²) in the worst case; **Buchheim, Jünger & Leipert** give a **linear-time** variant, which is what D3's `d3.tree()` and Mermaid's Tidy Tree layout implement.

**Radial layout.** Reuse the tidy-tree by reinterpreting its coordinate system: treat node `x` as **angle** and `y` as **radius**. In D3: `d3.tree().size([2*Math.PI, radius])` then map with `radialPoint(x,y) = [y·cos(x−π/2), y·sin(x−π/2)]` and `d3.linkRadial()`. For radial maps use the radius-scaled separation `(a.parent==b.parent ? 1 : 2) / a.depth` so deep rings don't over-spread. Radial is the classic *Buzan look* (root centered, branches sweeping out).

**Left/right balancing.** A true mind map splits root's children into two halves so the diagram grows both directions and stays compact — markmap and most editors auto-balance; FreeMind/Freeplane persist the choice per-node via `POSITION="left|right"`. Auto-balancing typically alternates or bin-packs children by subtree weight to equalize the two sides' heights.

**Rendering substrate.** **SVG** (markmap, jsMind-svg, d3) gives crisp text, hit-testing, and CSS styling, and scales best for large interactive maps; **HTML/DOM** nodes (some editors) ease rich content but cost layout perf; **canvas** (jsMind default) is fast to draw but needs manual hit-testing and re-render. Rule of thumb: many nodes + interactivity → SVG; thousands of static nodes → canvas. (Substrate trade-offs for general graphs: **graph-network-3d-visualization**.)

### 5. Tool landscape (pointer level)

- **XMind** — polished commercial app; `.xmind` zip format, OPML/FreeMind import-export, AI features; good as a conversion hub and for structured maps.
- **Freeplane / FreeMind** — open-source Java; `.mm` XML, scripting, arrowlinks/cross-links, conditional styles; Freeplane is the actively maintained successor.
- **CmapTools (IHMC)** — the reference *concept-map* tool; CXL format, CmapServer + KEA web services, collaborative editing, resource linking into "Knowledge Models."
- **MindMeister / Mindomo / Coggle** — web SaaS mind mappers (collaboration, OPML/`.mm`/image export).
- **Obsidian** — Mind Map / Enhanced Mind Map community plugins render a note's markdown headings as a markmap-style mindmap; Canvas gives a free-form graph surface.
- **Logseq** — built-in document/mindmap view of the outliner's block tree (its outline *is* the map).
- **Markmap / Mermaid** — the code-first, text-as-source options for docs and pipelines (covered above).

When in doubt: text-driven and version-controllable → markmap or Mermaid; rich interactive editing → XMind/Freeplane; labeled-edge concept maps → CmapTools/CXL; in-app embedding → jsMind or GoJS.

## Sources

- markmap docs (lib/cli/view): https://markmap.js.org/docs/markmap
- Cognitive Maps, Mind Maps, and Concept Maps: Definitions (NN/g): https://www.nngroup.com/articles/cognitive-mind-concept/
- Davies — Concept mapping, mind mapping and argument mapping: https://philpapers.org/archive/DAVCMM.pdf
- Mermaid mindmap syntax: https://docs.mermaidchart.com/mermaid-oss/syntax/mindmap.html
- OPML 2.0 specification: http://2005.opml.org/spec2.html
- OPML 2.0 (Library of Congress digital formats): https://www.loc.gov/preservation/digital/formats/fdd/fdd000554.shtml
- FreeMind file format & Freemind.xsd: https://freemind.sourceforge.io/wiki/index.php/File_format
- Current Freeplane File Format (XSD): https://dpolivaev.github.io/freeplane-docs/attic/old-mediawiki-content/Current_Freeplane_File_Format.html
- XMind File Format (DeepWiki): https://deepwiki.com/zhuifengshen/xmind/1.2-xmind-file-format
- CXL: An XML-based language for Cmaps (IHMC): https://cmap.ihmc.us/xml/CXL.html
- Novak & Cañas — Theory Underlying Concept Maps (IHMC): https://cmap.ihmc.us/docs/theory-of-concept-maps
- Reingold & Tilford — Tidier Drawings of Trees: https://reingold.co/tidier-drawings.pdf
- D3 d3-hierarchy tree (Reingold-Tilford/Buchheim, radial): https://d3js.org/d3-hierarchy/tree
- jsMind data formats & usage docs: https://github.com/hizzgdev/jsmind/blob/master/docs/en/1.usage.md
