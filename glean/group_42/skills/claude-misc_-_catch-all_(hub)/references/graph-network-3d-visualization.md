<!-- hub-reference-banner -->
> **Reference file — part of the `misc-catch-all` hub.** Formerly the standalone `graph-network-3d-visualization` skill.
> Sibling topics in this family are now reference files under the hubs (`misc-catch-all`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: graph-network-3d-visualization
version: "1.0.0"
updated: "2026-06-19"
description: >-
  Drawing and laying out graphs/networks and 3D node-link scenes in the browser — the rendering/layout layer, not analytics. TRIGGER: choosing a graph LAYOUT (force-directed/ForceAtlas2/d3-force, hierarchical/Sugiyama via dagre/ELK, circular, radial/tree, stress majorization) or seeding it for determinism; picking a 2D library (D3-force, Cytoscape.js, sigma.js+graphology, vis-network, React Flow) and SVG-vs-Canvas-vs-WebGL substrate by node count; building a 3D/WebGL/VR graph (three.js, 3d-force-graph, ngraph, deck.gl GraphLayer, reagraph, WebXR); scaling to 10k+ nodes (GPU instancing, LOD, semantic zoom, edge bundling, clustering). SKIP: graph ANALYTICS (centrality, community detection, PageRank) → da-27-network-graph-analytics; knowledge-graph/RDF semantics → da-41-knowledge-graphs-and-semantic-analytics; generic chart libraries (bar/line/pie) → programmatic-charting-libraries.
---

# Graph, Network & 3D Visualization Rendering

This skill covers **drawing** graphs/networks — computing node positions (layout) and painting nodes/edges to a screen (rendering) in 2D, 3D, and VR. It is the rendering counterpart to graph analytics: it answers "how do I lay this out and make it interactive at this scale," not "what does the structure mean." For centrality/community detection/PageRank see **da-27-network-graph-analytics**; for RDF/ontology semantics see **da-41-knowledge-graphs-and-semantic-analytics**; for bar/line/pie charts see **programmatic-charting-libraries**.

The two decisions that dominate every graph-viz project: **(1) which layout algorithm** computes positions, and **(2) which rendering substrate** (SVG / Canvas / WebGL) paints them — the latter is the single biggest determinant of how many nodes you can show.

## Core concepts

### 1. Layout-algorithm decision table

A *layout* maps graph topology to coordinates. Force-directed methods ("spring embedders") use only the graph's own structure — repulsion between all nodes, attraction along edges — minimizing a layout "energy"; they are intuitive and reveal symmetry/clusters but are **non-deterministic** (depend on initial state, can settle in local minima) and node positions have no Cartesian meaning (read proximity, not coordinates).

| Layout family | Algorithms / impls | Best for | Notes / determinism |
|---|---|---|---|
| Force-directed (energy) | Eades spring; **Fruchterman–Reingold** (FR); Kamada–Kawai (graph-distance springs); **ForceAtlas2** (Gephi); Yifan Hu; OpenOrd | General undirected graphs; revealing clusters/communities; social networks | FR: best quality but **O(n²)**, too slow on large graphs without Barnes–Hut. FA2 ≈ Yifan Hu in quality/speed and is *continuous* (tune live); LinLog mode = more cluster separation, slower. All non-deterministic — **seed** for reproducibility. |
| Velocity-Verlet force sim | **d3-force** (forceManyBody / forceLink / forceCenter / forceCollide / forceX,Y / forceRadial) | Interactive web graphs ≤ a few thousand nodes; custom force composition | Modular: a "force" is any fn mutating node velocity/position. Deterministic *given* fixed initial positions (d3 seeds NaN positions via phyllotaxis). For large graphs run the sim to stop in a **Web Worker** to avoid freezing the UI. |
| Hierarchical / layered (Sugiyama) | **dagre** (simple, fast); **ELK / elkjs** (`elk.layered`); d3-dag; MSAGL | DAGs, flowcharts, org charts, pipelines, dependency trees — anything with edge *direction* | Sugiyama = 5 phases: cycle-breaking → layer assignment → crossing minimization → node placement → edge routing. dagre is a near drop-in; ELK is far more configurable (ports, orthogonal/spline routing, **nested/compound** subflows) but async (Java→JS port, bigger bundle). |
| Tree / radial | d3-hierarchy (tidy tree, cluster, treemap, partition, pack); breadthfirst (Cytoscape) | Single-rooted trees; radial dendrograms | Deterministic. d3-hierarchy needs a true tree/forest. |
| Circular / concentric | Cytoscape `circle`, `concentric` | Showing per-node edge density (circle); relative importance by a metric (concentric rings) | Deterministic; cheap. |
| Stress majorization / MDS | OGDF Stress Minimization, Pivot MDS; graphology-layout-* | Faithful global distances; quality static layout; 3D | Optimizes graph-theoretic distance ↔ Euclidean distance; better global shape than basic spring models, more expensive. Common in VR/3D pipelines. |

Rule of thumb: **directed/hierarchical → dagre (simple) or ELK (complex); everything else → a force layout (try fCoSE/ForceAtlas2 first); trees → d3-hierarchy.** Always **pin a seed / fixed initial positions** when you need the same picture twice (tests, screenshots, user expectations).

### 2. 2D library landscape + substrate trade-offs

The rendering substrate sets the node-count ceiling more than anything else:

| Substrate | Practical ceiling (approx.) | Strengths | Weakness |
|---|---|---|---|
| **SVG** | ~1–2k nodes | Crisp text, per-element CSS/DOM events, easy custom shapes, accessibility | DOM node per element → slow past a couple thousand |
| **Canvas** (2D) | ~5k nodes/edges | Faster bulk draw, still CPU; immediate-mode | No DOM picking (hit-test manually); redraws whole frame |
| **WebGL** | ~10k usable; 50k–500k with tuning | GPU-batched draw calls, smooth pan/zoom on huge graphs | Custom shapes need GLSL; labels limited by texture-atlas size; harder to build |

(yWorks benchmark: SVG ≈ 2k, Canvas ≈ 5k, WebGL ≈ 10k+ usable on mid hardware; WebGL "shines from ~1,000 nodes up," and modern integrated GPUs handle ~50k.)

| Library | Substrate | Sweet spot |
|---|---|---|
| **D3-force / D3** | SVG or Canvas | Few-hundred to low-thousands; maximum custom rendering control; you wire your own render loop |
| **Cytoscape.js** | Canvas | Richest layout/algorithm set + styling; bio/network analysis; `hideEdgesOnViewport` to speed large graphs. Layouts: `grid`, `circle`, `concentric`, `breadthfirst`, `cose` → `cose-bilkent` → **`fcose`** (try fcose first for force layouts; supports compound graphs + constraints) |
| **sigma.js + graphology** | **WebGL** | Thousands→100k+ nodes. **Separation of concerns**: graphology owns the data model + algorithms (incl. ForceAtlas2, metrics); sigma owns WebGL rendering/interaction. Native rendering is circles+edges; custom node shapes need shaders/overlays |
| **vis-network** | Canvas | Out-of-box interactivity (physics drag/drop, clustering, zoom); good for ≤ few-thousand interactive diagrams; hard to deeply extend |
| **React Flow (xyflow)** | SVG/DOM | Node-based editors, flow builders, whiteboards; bring-your-own layout (recommends dagre for trees, elkjs for complex); not for huge graphs |
| **reagraph / reaviz** | WebGL (React/three.js) | React WebGL graphs with built-in 2D/3D force, radial, tree, concentric, clustering, edge bundling |

### 3. The 3D / WebGL stack (with a concrete snippet)

The 3D ecosystem centers on **three.js** for rendering and **d3-force-3d** or **ngraph** for physics. Vasco Asturiano's `3d-force-graph` is the standard high-level component (a web component wrapping `three-forcegraph` + a force engine); for VR there is `3d-force-graph-vr` / `aframe-forcegraph-component` (A-Frame + WebXR), and an AR variant. **deck.gl** brings GPU-accelerated layers (its community `GraphLayer`, plus `ScatterplotLayer`/`PathLayer`/`ArcLayer`/aggregation layers) when you want layout/aggregation pushed onto the GPU and map-grade scale.

Concrete `3d-force-graph` setup with a custom node object, DAG mode, and tuned forces:

```js
import ForceGraph3D from '3d-force-graph';
import * as THREE from 'three';

const Graph = new ForceGraph3D(document.getElementById('graph'))
  .graphData({ nodes, links })          // { nodes:[{id,...}], links:[{source,target}] }
  .nodeId('id')
  .forceEngine('d3')                      // 'd3' (d3-force-3d) or 'ngraph'
  .numDimensions(3)                       // 1 | 2 | 3
  .dagMode('radialout')                   // hierarchy for DAGs: td|bu|lr|rl|zout|zin|radialout|radialin
  .nodeThreeObject(node =>                // custom 3D node instead of default sphere
     new THREE.Mesh(
       new THREE.SphereGeometry(node.val ?? 4),
       new THREE.MeshLambertMaterial({ color: node.color ?? '#1f77b4' })
     ))
  .nodeThreeObjectExtend(false)           // true = decorate default sphere instead of replacing
  .linkDirectionalParticles(2);           // animated flow along links

// Tune the underlying d3-force-3d simulation:
Graph.d3Force('charge').strength(-120);   // many-body repulsion (more negative = more spread)
Graph.d3Force('link').distance(40);
```

Use `nodeThreeObject` for "3D mindmap"/knowledge-graph-in-3D patterns (sprites with text, icons, billboards). `dagMode` turns a force layout into a layered/radial hierarchy for DAGs. For VR, swap to `3d-force-graph-vr` (A-Frame) or feed precomputed x/y/z from a stress/organic layout into a WebXR (A-Frame + three.js) scene — VR research (Joos et al.; yFiles AR/VR) commonly precomputes positions with Stress Minimization / Pivot MDS / multi-level embedding rather than live force sim.

### 4. Large-graph performance (10k+ nodes)

- **Substrate first**: above ~5k go WebGL (sigma.js, reagraph, deck.gl, raw three.js). Canvas/SVG won't keep 60fps.
- **GPU instancing**: render N identical primitives (spheres/circles/arcs) in one instanced draw call — the deck.gl/three.js pattern; cost is ~proportional to vertex-shader invocations (≈ data length) and **fragment invocations (pixels drawn)**, so keep node radii small — overdraw on millions of points can hit billions of fragments/frame.
- **GPU compute for layout/aggregation**: deck.gl GPU aggregation crosses over ~100k objects (e.g. ~2.7× faster at 100k, ~11× at 1M vs CPU); GPGPU force layout (ParaGraphL-style) moves the spring simulation into shaders. Note browser limits: Chrome caps a single allocation at ~1GB, so layers crash somewhere between 10M–100M items unless you chunk.
- **Level of Detail (LOD)**: drop labels/detail when zoomed out; **switch substrates by zoom** — WebGL for the whole-graph overview, high-fidelity SVG when zoomed in past a threshold (yWorks "large graphs" pattern). Mind WebGL texture-atlas size limits for many labels.
- **Semantic zoom**: change *what* is shown by zoom level — bundled aggregate flows at macro scale, density hexbins at meso, detail at micro (deck.gl supply-chain pattern, 254× data reduction).
- **Edge bundling**: route edges along shared paths to cut clutter in dense graphs (skeleton-based SBEB, spanner-based S-EPB, force-directed bundling); run the iterative bundling in a **Web Worker** or on the GPU (WebGPU). reagraph/reaviz expose bundling out-of-box.
- **Clustering / aggregation**: collapse communities into meta-nodes (vis-network clustering; expand/collapse) to keep the on-screen element count bounded. Tile-pyramid / sleeve-routing approaches (MSAGL.js) treat the graph like an online map for client-side browsing of 30k+ nodes.

### 5. Interactivity

Core interactions: **zoom/pan** (constant in WebGL via GPU transform), node drag (re-heat the force sim or pin via `fx/fy`), hover/select highlighting, expand/collapse, lasso/path-finding. Hit-testing differs by substrate: SVG/DOM gives free event targets; Canvas/WebGL need an offscreen picking pass or spatial index. For drag in a live force sim, fix the dragged node's position during drag and release it (set `fx=null`) after. Keep interactions on the main thread responsive by offloading layout/bundling to Web Workers and only streaming back positions.

**Quick chooser**: directed/flow editor → React Flow + dagre/ELK · rich-styled analysis graph ≤ few-thousand → Cytoscape.js (fcose) · 10k–500k nodes 2D → sigma.js + graphology · 3D / knowledge-graph-in-3D / VR → 3d-force-graph (+ `-vr`) · GPU/map-scale or aggregation → deck.gl · full custom 2D render → D3-force.

## Sources

- ForceAtlas2 (Jacomy et al., PLOS ONE): https://journals.plos.org/plosone/article?id=10.1371/journal.pone.0098679
- Fruchterman & Reingold — force-directed placement: https://reingold.co/force-directed.pdf
- d3-force README (velocity Verlet, forces, worker): https://github.com/d3/d3-force/blob/v2.1.1/README.md
- Sigma.js official site (WebGL graph rendering, graphology): https://www.sigmajs.org/
- Cytoscape.js — built-in layouts & core API: https://js.cytoscape.org/
- Using layouts (Cytoscape.js blog): https://blog.js.cytoscape.org/2020/05/11/layouts/
- React Flow — Layouting overview (dagre/d3-hierarchy/d3-force/ELK): https://reactflow.dev/learn/layouting/layouting
- ELK Layered algorithm reference (Sugiyama 5-phase): https://eclipse.dev/elk/reference/algorithms/org-eclipse-elk-layered.html
- vasturiano/3d-force-graph (ThreeJS/WebGL): https://github.com/vasturiano/3d-force-graph
- deck.gl — Layer Rendering Performance: https://github.com/visgl/deck.gl/blob/master/docs/developer-guide/performance.md
- deck.gl-community — GraphLayer API: https://visgl.github.io/deck.gl-community/docs/modules/graph-layers/api-reference/layers/graph-layer
- deck.gl — Aggregation Layers (CPU vs GPU): https://deck.gl/docs/api-reference/aggregation-layers/overview
- yFiles WebGL2 rendering & Large Graphs demo: https://docs.yworks.com/yfiles-html/dguide/advanced/webgl2.html
- A Comparison of JavaScript Graph/Network Visualisation Libraries (Cylynx): https://www.cylynx.io/blog/a-comparison-of-javascript-graph-network-visualisation-libraries/
- reaviz/reagraph (WebGL React graph): https://github.com/reaviz/reagraph/
- 3D Force-Directed Graph in VR (A-Frame/WebXR): https://vasturiano.github.io/3d-force-graph-vr/
- Aesthetic-Driven Navigation for Node-Link Diagrams in VR (Joos et al.): https://dl.acm.org/doi/fullHtml/10.1145/3607822.3614537
