<!-- hub-reference-banner -->
> **Reference file — part of the `frontend-ui` hub.** The OBJECTIVE / quantitative-scoring layer of the
> visual-design-critique family: computational aesthetics + automated UI-quality & accessibility metrics.
> It complements the *subjective* critique references (`web-design`, `ui-ux-pro-max`,
> `accessibility-ux-reviewer`), it does not replace them. Sibling topics in this family are reference files
> under the hubs (`frontend-ui`), **not** standalone skills. Ignore any "use the X skill" /
> `related_skills` / SKIP pointers below that name a bare sibling skill; load that topic's
> `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: computational-aesthetics-ui-metrics
title: Computational Aesthetics & Automated UI-Quality / Accessibility Metrics
description: >-
  The OBJECTIVE / quantitative scoring layer for visual-design & UI critique, machine-computable metrics
  that feed, never replace, human/LLM judgment (they explain only ~48-49% web / ~32% mobile aesthetic
  variance; flag outliers and quantify deltas, never rank or auto-score). TRIGGER: score/measure a design
  objectively; colourfulness / visual-complexity / clutter / symmetry / white-space / grid metric; Aalto
  Interface Metrics (AIM); saliency or attention heatmap for a UI; compute a WCAG contrast ratio, run
  axe-core, or read a Lighthouse a11y score as a numeric signal; target-size 2.5.8 check; before/after
  visual-quality delta. SKIP: subjective visual critique / design principles → web-design, ui-ux-pro-max;
  WCAG/ARIA compliance review (not metric computation) → accessibility-ux-reviewer; chart/data-viz
  perception or colour-blind / palette choice → da-8-data-visualization; prompting a VLM to critique a
  screenshot → ai-mcp-sdk-prompting.
category: developer
tags:
  - developer
  - frontend
  - ui
  - ux
  - design
  - accessibility
  - metrics
  - computational-aesthetics
keywords:
  - computational aesthetics
  - visual complexity
  - colorfulness
  - Hasler-Susstrunk
  - color harmony
  - visual balance symmetry
  - feature congestion clutter
  - edge density
  - white space ratio
  - grid quality
  - Aalto Interface Metrics AIM
  - visual saliency attention prediction
  - saliency metrics AUC NSS CC SIM
  - UMSI UEyes
  - axe-core
  - Lighthouse accessibility score
  - WCAG contrast ratio
  - target size 2.5.8
  - APCA
  - aesthetic metric validity limits
whenToUse:
  - "you want an objective, reproducible number to complement a subjective design critique"
  - "computing or interpreting colourfulness, visual complexity, clutter, symmetry, white-space, or grid metrics on a screenshot"
  - "choosing or reading the Aalto Interface Metrics (AIM) toolkit output for a UI"
  - "generating or interpreting a visual-saliency / attention-prediction heatmap for a UI or design"
  - "computing a WCAG contrast ratio, running axe-core, or reading a Lighthouse accessibility score as a signal"
  - "measuring a before/after or A/B visual-quality delta for the same design lineage"
  - "deciding how much weight to put on an automated aesthetic or accessibility score"
whenNotToUse:
  - "subjective visual critique, design principles, style/palette/typography guidance, use web-design / ui-ux-pro-max"
  - "WCAG/WAI-ARIA compliance review and remediation (not metric computation), use accessibility-ux-reviewer"
  - "chart / data-visualization perception and chart-type choice, use da-8-data-visualization"
  - "prompting a vision model to critique a screenshot in prose, use ai-mcp-sdk-prompting (vision-model UI critique)"
related_skills:
  - web-design
  - ui-ux-pro-max
  - accessibility-ux-reviewer
  - da-8-data-visualization
version: "1.0.0"
updated: "2026-06-16"
---

# Computational Aesthetics & Automated UI-Quality / Accessibility Metrics

> **verified-as-of: 2026-06-16** for all volatile claims (AIM metric count/version & live-service status;
> axe-core coverage %; WCAG 2.2 status; APCA draft status; commercial attention-tool accuracy figures).
> Re-verify these against the cited sources on refresh; do not date-bump without a re-fetch.

## Overview

This reference is the **objective / quantitative scoring layer** of the visual-design-critique family. It
catalogs metrics that a machine can compute from a screenshot (or a rendered DOM) and that quantify
properties of a UI or graphic design: how colourful, how complex/cluttered, how symmetric/balanced, how
grid-aligned, how much white space, where attention is predicted to go, and whether it clears objective
accessibility floors.

**The single most important framing (read before using any metric below).** These metrics correlate only
**weakly-to-moderately** with human aesthetic judgment. The strongest published UI-aesthetics models
explain roughly **48–49% of the variance** in web-page aesthetic ratings and only **~32%** for mobile
apps, leaving half to two-thirds of the judgment unexplained, and the web figure leans partly on
*demographic* predictors, not pixels alone.[^reinecke13][^miniukovich15] Automated accessibility tools
catch only a limited fraction of real issues (Deque's own measure: ~57% by issue volume; the older
cross-tool consensus: ~30–40% by criteria coverage).[^deque57][^b13] Therefore: **use these as
complementary objective signals that feed a human or LLM critique. Never treat them as ground truth, never
rank designs as "objectively better," never auto-score quality.** The validity section (§6) is not a
footnote; it is the operating contract for everything above it.

What these metrics are *good* for: flagging candidates, surfacing things humans under-weight (contrast
failures, extreme clutter, gross imbalance), and quantifying **deltas** within the *same* design lineage
(before/after, A/B) where confounds are held constant. What they are *not* for: declaring a winner across
unrelated designs, or substituting for the semantic, brand, content, task, and emotional judgment that
lives in the other critique references.

## Core Concepts

### 1. Image-statistic aesthetic metrics (pixel-computable)

These need only a rendered image, no DOM.

**Colourfulness — Hasler & Süsstrunk (2003).** The canonical metric, from "Measuring Colourfulness in
Natural Images" (IS&T/SPIE Electronic Imaging vol. 5007).[^hasler] It works in a simple **opponent color
space** derived directly from RGB: `rg = R − G` and `yb = ½(R + G) − B`. The headline metric is:

```
M = sqrt(σ_rg² + σ_yb²) + 0.3 · sqrt(μ_rg² + μ_yb²)
```

where σ and μ are the standard deviation and mean of each opponent channel; the **0.3** weight on the mean
term is the paper's empirically fitted coefficient.[^hasler][^pyimage] The metric achieved >90% correlation
with subjective colourfulness ratings in the original study.[^hasler] Hasler & Süsstrunk also define a
"simplified" form pooling the combined opponent magnitudes (`C ≈ σ_RGYB + 0.3·μ_RGYB`), *tentative on the
exact pooled constants; the structure is well-attested but byte-verification of the simplified equation
failed.*[^hasler] Reinecke et al. (CHI 2013) reused a colourfulness model on website screenshots, and AIM
implements colourfulness as metric `cp6` citing Hasler & Süsstrunk.[^reinecke13][^aim]

**Colour count vs. colourfulness are distinct.** Colourfulness is a perceptual *spread* statistic;
**colour count** (unique RGB values) and **dominant/quantized colour count** are cardinality measures —
not interchangeable.[^hasler][^bvdart] Dominant-colour extraction reduces millions of pixels to ~8–16
representative colours via clustering (K-means, median-cut).[^bvdart] AIM exposes a family of colour
metrics: PNG file size, unique RGB count, HSV stats, unique HSV, LAB colours, colourfulness (cp6), static
and dynamic colour clustering, luminance SD, **WAVE**, and colour harmony.[^aim]

**Colour harmony (Itten / Matsuda / Cohen-Or).** Cohen-Or et al., "Color Harmonization" (ACM TOG /
SIGGRAPH 2006), formalizes harmony as **template-fitting on the hue wheel**, harmonic colours are defined
by *relationships* (positions in hue space), not specific colours; rotatable templates (types i, V, L, I,
T, Y, X) derive from Matsuda/Tokumaru harmonic schemes, with Itten's schemes as the art-theory
antecedent.[^cohenor] **WAVE (Weighted Affective Valence Estimates)**, Palmer & Schloss (PNAS 2010), is a
colour-*preference* predictor built from human object-association ratings; it explained ~80% of colour-
preference variance, beating cone-contrast/appearance/emotion models. Note WAVE is human-rating-derived,
**not** a pure pixel statistic.[^palmer]

**Visual complexity.** Computable proxies: **compression/file size** (JPEG/PNG size — more complex images
compress less → larger files), **quadtree decomposition leaf count** (recursively split until blocks are
homogeneous; more leaves = more complexity), **edge density**, fractal dimension, and information-theoretic
**entropy**.[^saraee][^karjus][^quadtree] A feature model combining entropy + edge density + JPEG ratio
reached r ≈ 0.70 against human complexity ratings.[^saraee] Miniukovich & De Angeli decompose interface
complexity into computable factors: visual clutter, contour/edge density, figure-ground contrast, colour
variability, symmetry.[^miniukovich14][^miniukovich15]

**Edge density / feature congestion / clutter — Rosenholtz, Li & Nakano (Journal of Vision 2007,
"Measuring visual clutter").**[^rosenholtz] Three measures: **Feature Congestion** (clutter from local
variability in colour, contrast, and orientation across scale), **Subband Entropy** (clutter via entropy
of wavelet subband coefficients), and **Edge Density** (apply **Canny** edge detection; take the
*proportion of edge pixels*, implementable in ~two lines of MATLAB).[^rosenholtz][^clutterpy] Edge density
generalizes as (edge pixels)/(total pixels) and is widely reused for GUI contour density.[^miniukovich15]

**Pixel symmetry & figure-ground contrast.** Pixel **mirror-symmetry** compares an image against its flip
across vertical, horizontal, and (for square images) diagonal axes, reported as a percentage.[^toolbox]
Beware polarity conventions: some toolboxes report "mean symmetry near 0% = balanced" (Wilson &
Chatterjee) and others "higher % = more symmetric", flag when comparing across tools.[^toolbox]
**Figure-ground contrast** measures how distinctly foreground separates from background; it is consistently
named as a Miniukovich complexity factor, though the exact pixel formula is *tentative* here.[^miniukovich15]

### 2. The Aalto Interface Metrics (AIM) toolkit & the computational-aesthetics-of-UI research line

**AIM** (Aalto Interface Metrics) is an open-source web service and codebase from the User Interfaces group
at Aalto University (Oulasvirta lab) that pools empirically validated models of human perception and
attention into one tool. You submit a GUI as a **URL or uploaded screenshot**, pick metrics, and get
numeric results plus visualizations and comparison against a reference website set.[^aim][^aimportal] The
canonical paper is Oulasvirta et al., "Aalto Interface Metrics (AIM): A Service and Codebase for
Computational GUI Evaluation," **UIST 2018 Adjunct** (DOI 10.1145/3266037.3266087).[^aim] Service:
`interfacemetrics.aalto.fi` (a JS SPA; live as of 2026-06-16, no deprecation notice found). Code:
`github.com/aalto-ui/aim`, **MIT** license (Python/Tornado + MongoDB + Vue; Selenium/Chrome capture).[^aimrepo]

AIM organizes metrics into **four categories**: **Colour Perception, Perceptual Fluency, Visual Guidance,
Accessibility**.[^aimrepo][^aimportal] The 2018 paper cites **"17 metrics"**, treat 17 as the *launch*
figure; the codebase has since grown (≈21 user-facing metrics in v1 plus an `aim2`/v2.3 line adding
segmentation metrics), so the exact live count is **version-dependent, use "17+".**[^aim][^aimrepo]
Verified AIM metric IDs worth knowing (category → ID → metric → basis):[^aimrepo]

| Category | ID | Metric | Basis |
| --- | --- | --- | --- |
| Colour Perception | `cp6` | Colourfulness | Hasler-Süsstrunk |
| Colour Perception | `cp10` | WAVE | Palmer & Schloss |
| Perceptual Fluency | `pf1` | Edge density | Rosenholtz-family |
| Perceptual Fluency | `pf2` | Edge congestion | Rosenholtz-family |
| Perceptual Fluency | `pf3` | JPEG file size | complexity proxy |
| Perceptual Fluency | `pf4` | Figure-ground contrast | Miniukovich & De Angeli |
| Perceptual Fluency | `pf5` | Pixel symmetry | (no external cite) |
| Perceptual Fluency | `pf6` | Quadtree decomposition (balance/symmetry/equilibrium/leaves) | Ngo / Reinecke |
| Perceptual Fluency | `pf7` | White space | Miniukovich & De Angeli |
| Perceptual Fluency | `pf8` | Grid quality | Miniukovich & De Angeli |
| Visual Guidance | `vg1` | Saliency | Itti-Koch |
| Visual Guidance | `vg2` | Visual search (experimental) | Jokinen et al. |
| Accessibility | `ac1` | Colour-blindness simulation | Machado 2009 |

AIM is essentially an **implementation of the research line below**: its metric files cite Miniukovich &
De Angeli, Reinecke, Ngo, and Hasler & Süsstrunk directly.[^aimrepo]

**The research line (2013–2015) — what each established:**

- **Reinecke et al., CHI 2013** — "Predicting Users' First Impressions of Website Aesthetics with a
  Quantification of Perceived Visual Complexity and Colorfulness." 450 websites, 548 participants; models
  of **complexity + colourfulness** plus demographics explained **~48% (≈ half)** of aesthetic appeal
  judged after **500 ms**.[^reinecke13]
- **Miniukovich & De Angeli, AVI 2014** ("Quantification of Interface Visual Complexity") + **NordiCHI
  2014**, structural complexity determinants across information *amount* (clutter, colour variability),
  *organization* (symmetry, grid, grouping), and *discriminability* (contour density, figure-ground
  contrast); ~40% of subjective complexity variance.[^miniukovich14][^nordichi]
- **Miniukovich & De Angeli, CHI 2015** ("Computation of Interface Aesthetics," Honorable Mention) —
  **eight** automatic GUI-aesthetics metrics (incl. symmetry, colour range, visual clutter, white space,
  grid quality); best models explained **up to 49% (web)** and **~32% (iPhone apps)** of aesthetics
  variance.[^miniukovich15] These eight are the direct lineage of AIM's perceptual-fluency set.
- **Zen & Vanderdonckt, RCIS 2014** — region-based GUI aesthetic model with seven geometric metrics
  (density, symmetry, balance, proportionality, uniformity, simplicity, sequence), implemented in a web
  tool; descends from Ngo's classic measures.[^zen]

### 3. Layout / grid / white-space metrics (need element boxes, DOM or segmentation)

**Grid quality = count of distinct alignment lines.** AIM's grid-quality metric collects every element's
x-start/x-end into a set and counts distinct values (vertical alignment lines), does the same for
y-start/y-end (horizontal), and sums them; **fewer distinct alignment lines = tighter gridness** (cites
Miniukovich & De Angeli 2015).[^aimrepo][^aimportal] "Number of distinct element sizes" is a related
*heuristic*, not a validated standalone metric.[^aimrepo]

**White-space ratio** = `(image_area − Σ element_box_area) / image_area` (AIM `pf7`, cites Miniukovich &
De Angeli 2015).[^aimrepo] Harrington et al.'s document-aesthetics lineage refines this with **white-space
fraction** *and* **white-space free-flow** (how connected the empty area is, scattered gaps read worse
than one coherent margin).[^toolbox]

**Ngo et al. (2003), "Modelling interface aesthetics"** — **fourteen** measures (balance, equilibrium,
symmetry, sequence, cohesion, unity, proportion, simplicity, density, regularity, economy, homogeneity,
rhythm, order), each normalized ~[0,1] and **computed from element bounding boxes** (position, size, area,
center).[^ngo][^toolbox] Key layout-level formulas (verified):

- **Balance (BM):** `BM_vertical = (w_L − w_R)/max(|w_L|,|w_R|)`, `BM_horizontal = (w_T − w_B)/max(|w_T|,|w_B|)`,
  where each side-weight `w_j = Σ_i (area_i · distance_i)` of element centers from the layout center line.
  AIM returns `BM = 1 − (|IB_top_bottom| + |IB_left_right|)/2`.[^toolbox][^aimrepo]
- **Equilibrium (EM):** difference between the element center-of-mass and the physical center, normalized
  (1 = centered).[^ngo][^aimrepo]
- **Symmetry (SM):** axial duplication about vertical/horizontal/diagonal axes via per-quadrant sums.[^ngo][^aimrepo]
- **Density (DM):** screen coverage by objects, with an *optimal mid-range* (not monotonic).[^ngo]
- **Regularity (RM):** consistency of spacing + number of distinct alignment points.[^ngo][^toolbox]

**Element density / leaf count.** Reinecke et al. segment a screenshot via **quadtree** (split on minimum
colour/intensity entropy) and **space-based** decomposition, returning **number of leaves** plus text-group
/ image-area counts as density proxies.[^reinecke13] AIM's quadtree returns "Number of Leaves" alongside
balance/symmetry/equilibrium.[^aimrepo]

**Modern/ML layout scoring (flag maturity).** **Rico** (Deka et al., UIST 2017) — ~66k–72k Android UI
screens with view hierarchies + screenshots; the authors train an **autoencoder to a 64-d layout
embedding** for similarity search/generation.[^rico] Rico scores layout *similarity*, not aesthetic
quality; learned *aesthetic* scoring is **research-grade, not turnkey**, treat with caution.[^rico]

**DOM-vs-pixel distinction.** Box-based metrics (Ngo, white space, grid quality) want a *list of element
rectangles*, free from a real DOM (`getBoundingClientRect`, Rico view hierarchies); without one you must
**segment the screenshot into boxes** first. AIM uses a *quadtree* as a stand-in for elements precisely
because better segmentation was unavailable, so AIM's balance/symmetry/equilibrium are **pixel-derived
approximations** of the box-based Ngo formulas.[^aimrepo][^reinecke13] Cleanest pipeline: render tree for
box-based metrics; quadtree/leaf metrics for image-only cases.

### 4. Visual saliency / attention prediction

The conceptual move: from **free-viewing natural-image saliency** ("where the eye is *pulled* by low-level
conspicuity") to **task/communication-driven visual *importance*** ("what the designer made matter"). These
are *not the same signal*, and conflating them is the core error in applying saliency to UIs.[^bylinskii17][^umsi][^agd]

**Classic bottom-up.** **Itti, Koch & Niebur (1998)** — the foundational model: multiscale colour/
intensity/orientation features via Gaussian pyramids → **center-surround** differences → a single
topographical **saliency map** → winner-take-all selection. Models rapid *task-free* viewing.[^itti] **GBVS
(Harel et al., NIPS 2006)** improved on it via Markov chains over feature maps (the chain's equilibrium
distribution = saliency).[^gbvs]

**The metric zoo (and why it matters).** Saliency evaluation splits into **location-based** (AUC-Judd,
AUC-Borji, **shuffled AUC / sAUC** which penalizes center bias, **NSS**) and **distribution-based** (**CC**
linear correlation, **SIM** histogram intersection, **KL**-divergence, **EMD**).[^transalnet][^salmetrics]
Bylinskii et al. (IEEE PAMI 2019, "What do different evaluation metrics tell us about saliency models?")
showed **metrics rank models differently**, they recommend **NSS and CC** as general-purpose, **sAUC**
when center-bias invariance matters. **Practical rule: a saliency model's "accuracy" is meaningless without
naming the metric.**[^bylinskii19]

**Deep saliency.** The **SALICON** project (Jiang et al., CVPR 2015) crowdsourced "free-viewing" attention
via a mouse-contingent click method on MS-COCO and seeded an early deep-saliency line.[^salicon]
**DeepGaze II/IIE**, **ML-Net**, and
**SAM** fine-tune ImageNet backbones; no single model wins every metric.[^borji]

**For UIs/designs specifically (the key part).**

- **Bylinskii et al., UIST 2017** ("Learning Visual Importance for Graphic Designs and Data
  Visualizations"), nets predict **relative element importance** on graphic designs (GDI annotations) and
  visualizations (BubbleView mouse clicks); drives retargeting/thumbnailing.[^bylinskii17]
- **BubbleView** (Kim et al., TOCHI 2017) — a mouse-contingent moving-window proxy where clicks reveal
  sharp "bubbles"; clicks accounted for **>75% of eye fixations** with 10–15 participants. A scalable
  *proxy*, not real gaze.[^bubbleview]
- **UMSI** (Fosco et al., UIST 2020, "Predicting Visual Importance Across Graphic Design Types") — one deep
  model (Xception + ASPP) jointly trained on posters/infographics/mobile-UIs **and** natural-image
  saliency, with auto input-type classification; introduces the **Imp1k** dataset (`predimportance.mit.edu`).
  Reported comparable to the best natural-image saliency models (SOTA), *self-reported, tentative*.[^umsi]
- **UEyes** (Jiang et al., CHI 2023) — the key UI **ground-truth** resource: a real eye tracker recorded
  **62 participants viewing 1,980 UIs** (495 each: webpage/mobile/desktop/poster), yielding ~20k eye-
  movement sequences + multi-duration saliency maps. Prior UI work leaned on mouse/manual proxies; UEyes
  gives true fixations and shows colour/location/gaze-direction differ across UI types.[^ueyes]

**How saliency feeds critique.** Overlay an attention/importance heatmap on a mockup to **validate visual
hierarchy**, does the eye land where intended; is the CTA/hero in a high-saliency region; what likely
draws **first fixation**.[^bylinskii17][^umsi] **Commercial predictive-attention tools** (Attention
Insight, Neurons/Predict) ship instant CNN heatmaps without live participants and market 90%+ accuracy vs.
eye-tracking, **vendor self-reported, predictions not measurements, lacking standardized third-party
validation** (*verified-as-of 2026-06-16; treat accuracy claims as tentative*).[^attninsight][^neurons][^heatmapcrit]

### 5. Automated accessibility / UI-quality checkers as objective signals

These are **deterministic, reproducible, machine-computable**, the same input yields the same number every
time, which makes them ideal *objective inputs* to critique. But every authoritative source frames them as
**necessary but not sufficient**.

**axe-core (Deque)** — the open-source rules engine under most automated a11y testing (and the engine
Lighthouse uses).[^axeapi][^axerepo][^lhscoring] Powers axe DevTools, `@axe-core/playwright`,
`@axe-core/cli`, `jest-axe`. **Rule tags** select standards: `wcag2a`, `wcag2aa`, `wcag21aa`, `wcag22aa`,
`best-practice`, plus regime tags (`section508`, `EN-301-549`, `ACT`), tags do **not** roll up, so cover A
*and* AA by listing both.[^axeapi][^mabl] **Impact levels** are Deque-assigned (not WCAG-defined): `minor`,
`moderate`, `serious`, `critical`.[^axeapi][^mabl] Results come in **four arrays**, `violations`, `passes`,
`incomplete` ("needs review", where axe parks anything it can't decide deterministically, e.g. *whether*
alt text is meaningful), `inapplicable`, which is what makes axe a *structured signal*.[^axeapi]
**Design principle: zero false positives** — axe pushes uncertainty into `incomplete` rather than guessing;
the deliberate cost is narrow coverage.[^axeapi][^accessproof]

**The limited-fraction claim (load-bearing — cite both figures, don't average).** Deque's *own* measured
figure: **~57%** of accessibility issues covered by automated testing (2,000+ audits, 13k+ pages),
explicitly reframing the older "20–30% of coverage" belief, Deque measures by *issue volume*.[^deque57]
The broader cross-tool research consensus: **~30–40%** of WCAG issues by *criteria coverage*, the rest
needing human judgment + assistive-tech testing.[^b13][^crosscheck] These are different measurement bases,
not a contradiction. *verified-as-of 2026-06-16.*

**Lighthouse accessibility audit (Google)** — runs a **subset of axe-core rules** and produces a **0–100
score** that is a **weighted average** of audits, weighted by **axe user-impact**; scoring is **binary per
audit** (no partial credit).[^lhscoring] **Lighthouse's own guidance is explicit:** "Only a subset of
accessibility issues can be automatically detected so manual testing is also encouraged"; un-automatable
checks appear under "Additional Items To Manually Check" and **do not affect the score**. **A score of 100
does NOT mean the page is accessible**, it means the automated checks passed.[^lhscoring][^afixt][^a11ytest]

**Programmatic WCAG contrast ratio (fully deterministic).** The formula:

```
contrast = (L1 + 0.05) / (L2 + 0.05)        # L1 = lighter, L2 = darker; range 1:1 .. 21:1
```

**Relative luminance** per WCAG: for each sRGB channel normalized to [0,1],
`c_lin = c/12.92 if c ≤ 0.03928 else ((c + 0.055)/1.055)^2.4`, then
`L = 0.2126·R_lin + 0.7152·G_lin + 0.0722·B_lin`.[^g17][^bickford] **Thresholds** (look up by success
criterion):[^contrastmin][^nontext][^targetmin]

| SC | Level | Threshold | Applies to |
| --- | --- | --- | --- |
| 1.4.3 Contrast (Minimum) | AA | 4.5:1 | normal text |
| 1.4.3 Contrast (Minimum) | AA | 3:1 | large text (≥18pt, or ≥14pt bold ≈ 24px / 18.5px) |
| 1.4.11 Non-Text Contrast | AA | 3:1 | UI components & graphical objects |
| 1.4.6 Contrast (Enhanced) | AAA | 7:1 / 4.5:1 | normal / large text |
| 2.5.8 Target Size (Minimum) | AA | 24×24 CSS px | pointer targets (spacing/inline/essential exceptions) |
| 2.5.5 Target Size (Enhanced) | AAA | 44×44 CSS px | pointer targets |

**APCA** (Accessible Perceptual Contrast Algorithm) is the **WCAG 3 candidate** that replaces the fixed
ratio with a lightness-contrast model outputting **Lc** (≈0 to ±106); visual-contrast was moved *out* of
the WCAG 3 Working Draft in **July 2023** and the WCAG 3 algorithm is **still TBD**, APCA is a
**candidate, non-normative, draft** as of 2026-06-16 (*tentative, single source-cluster*).[^apca][^apca2]

**Focus/tab order and the programmatic boundary.** Beyond contrast and target size (table above), other
machine-checkable signals include focus/tab order vs. tabindex (2.4.3), focus-visible (2.4.7/2.4.11),
reflow at 320px / 400% zoom (1.4.10), plus name/role/value, ARIA-validity, label-presence,
heading-structure, and alt-*presence*.[^axeapi] **What CANNOT be checked programmatically (needs a
human):** whether alt text is *meaningful*, whether reading/focus order is *logical*, whether captions are
*accurate*, whether link/heading text is *descriptive*. That boundary *is* the
57%-vs-rest / 30–40% split.[^accessproof][^lhscoring]

### 6. Validity, limits & the correct operating stance (the contract)

This section governs §1–§5. It is the most-cited and least-optional part of this reference.

**The variance ceiling.** Best published UI-aesthetics models cap at **R² ≈ 0.48–0.49 (web)** and **~0.32
(mobile)**, half to two-thirds of aesthetic judgment is **unexplained**, and the web figure leans partly
on demographics.[^reinecke13][^miniukovich15] No automated UI metric approaches a human rating panel's
reliability. Variance-explained also depends entirely on *which construct* you name (Miniukovich's
*complexity* sub-construct: ~40%).[^miniukovich15]

**Individual differences & context.** Reinecke & Gajos (CHI 2014, "Quantifying Visual Preferences Around
the World," **2.4M ratings**, ~40k participants) showed the complexity/colourfulness level at which appeal
*peaks* varies systematically by **country, gender, and education**, preference is **population-relative,
not universal**.[^reinecke14] Complexity preference is **non-monotonic** (Berlyne's inverted-U: moderate
preferred over very-low or very-high), so "minimize complexity" is wrong as a monotonic rule, and the
inverted-U is itself **boundary-dependent/contested** across domains.[^berlyne][^althuizen] **Prototypicality**
is a co-equal first-impression driver (Tuch et al., 17–50 ms) that surface metrics miss.[^tuch] The
**aesthetics↔usability halo** (Tractinsky et al. 2000, "what is beautiful is usable") means aesthetic
ratings are contaminated by, and contaminate, usability perception, not a clean isolated
signal.[^tractinsky][^halo]

**The classic literature — the cautionary origin story.** **Birkhoff (1933)**, `M = O/C` (order over
complexity), was the first formula for beauty and **largely failed empirically**, psychological studies
"largely failed to corroborate" it; Eysenck even argued for `M = O × C`.[^birkhoff][^birkhoff2] **Lavie &
Tractinsky (2004)** showed perceived aesthetics is **two factors**, *classical* (orderly, clean,
symmetrical) and *expressive* (creative, original, convention-breaking); a single scalar score conflates
them, and expressive aesthetics is nearly invisible to order-tuned metrics.[^lavie]

**The construct-validity gap (strongest disconfirming evidence).** Metrics often **do not measure the
construct they're named after**: studies found **no correlation between objective "balance" and subjective
balance**, and Ngo-style symmetry/balance showed **null** correlation with perceived aesthetics, overall
*shape* mattered more than computed center-of-mass.[^twosides][^equilibrium] AIM's own source notes its
symmetry metric "has not proven significant" in two studies and that quantifying symmetry "might be
problematic in HCI."[^aimrepo] Reinecke et al. *dropped* balance/symmetry/equilibrium as non-significant;
the surviving predictors were density/leaf counts.[^reinecke13]

**Why complementary, not ground truth.** Metrics measure the **visual surface only** — they miss
semantics, brand fit, content/copy quality, task/usability context, emotional resonance, trends, and
accessibility-of-meaning, which is where most of the unexplained variance lives.[^tuch][^lavie] **Goodhart's
law:** because metrics are proxies missing ~half the variance, **optimizing the metric can degrade the true
goal** (*inference grounded in the construct-validity gap + non-monotonicity + population-relativity;
tentative*).[^reinecke14][^berlyne][^twosides]

**Maturity & reproducibility honesty.** There is **no consensus** on which measure/formula/interpretation
to adopt.[^zen] Tools are **research artifacts**: AIM pools validated models but its README markets metrics
as "verified" *without* ground-truth or predictions-vs-measurements disclaimers, its framing outruns its
evidence base; treat outputs as **research-grade predictions, not measurements**.[^aim][^aimrepo]
Correlations are **dataset-dependent** and may not generalize to a new product/audience/year without
re-validation.[^miniukovich15][^reinecke14]

## Methodology, a metric-assisted critique pipeline

1. **Decide the input regime.** Have a live DOM? → box-based metrics (Ngo, white space, grid quality) are
   cheap and exact. Only a screenshot? → pixel metrics (colourfulness, clutter, edge density) + quadtree
   approximations; accept they are coarser.
2. **Run the deterministic floor first (accessibility).** Compute WCAG contrast ratios, run axe-core
   (`@axe-core/playwright`), read the Lighthouse a11y score. Treat results as a **falsifiable floor**, not
   a conformance claim. Triage axe `violations` by impact; route `incomplete` to human review.
3. **Compute aesthetic signals.** Colourfulness, visual complexity/clutter (Feature Congestion + Subband
   Entropy + edge density), white-space ratio, grid-quality (alignment-line count), Ngo balance/symmetry
   (if boxes available). Run AIM if a hosted toolkit is acceptable.
4. **Compute attention prediction.** Generate an importance/saliency heatmap (UMSI-class model for designs;
   note free-viewing-saliency-≠-UI-importance). Check whether high-importance regions match design intent.
5. **Attach correlation strength to every number, in a fixed output shape.** Never report a metric bare.
   Emit each signal as `<metric>: <value> [<unit/range>], <correlation strength vs. its construct> —
   <caveat>`. Worked: `colourfulness (cp6): 38, part of the complexity+colourfulness model that explains
   ~48% of 500 ms web-appeal variance, weak signal`; `Ngo balance: 0.72, ~null correlation with perceived
   balance, informational only, do not action`; `contrast: 3.9:1, deterministic, hard AA floor, FAILS
   1.4.3 (needs 4.5:1)`. **Decision rule for items flagged *tentative* in §1–§5** (the simplified
   colourfulness form, figure-ground contrast, UMSI/commercial accuracy, APCA status): report them as
   *direction-only* ("higher/lower than baseline"), never as a precise value cited to a formula or a
   vendor accuracy figure.
6. **Feed signals into the human/LLM critique**, alongside the subjective references (`web-design`,
   `ui-ux-pro-max`) and semantic/brand/content/task judgment. Report deltas, not verdicts.

## Practical Patterns

- **Quantify deltas, not absolutes.** Use metrics for **before/after** and **A/B** on the *same* design
  lineage where confounds are held constant, this is where they are most trustworthy.[^reinecke13]
- **Surface the high-confidence tails.** Let metrics flag extreme clutter, contrast failures, gross
  imbalance, colour overload, the things humans under-weight or miss.[^rosenholtz][^g17]
- **Make the a11y floor non-negotiable and falsifiable.** A 3.9:1 label is *measurably* below the 4.5:1 AA
  bar; a 20×20px tap target is *measurably* below 24px. Spend human judgment on the semantic layer axe
  leaves as `incomplete`.[^g17][^accessproof]
- **Above-the-fold vs full-page matters.** Complexity predicts aesthetics better for above-the-fold than
  full-page screenshots, be explicit about the capture region.[^reinecke13]
- **Open-source building blocks:** `visual-clutter` (Python, Feature Congestion + Subband Entropy);
  PyImageSearch colourfulness (OpenCV); AIM codebase (`aalto-ui/aim`); `@axe-core/*`; Lighthouse CLI;
  Canny via OpenCV for edge density.[^clutterpy][^pyimage][^aimrepo][^axerepo]

## Anti-Patterns

- **Ranking unrelated designs "objectively."** No UI metric explains more than ~half of aesthetic variance,
  and key drivers (prototypicality, expressive aesthetics, semantics, brand) are out of scope. Never crown
  a winner from a score.[^reinecke13][^lavie][^twosides]
- **Trusting a metric's name as its construct.** "Balance"/"symmetry" may not predict *perceived* balance —
  for Ngo-style measures the correlation was effectively null.[^twosides][^aimrepo]
- **Optimizing a single metric (Goodhart).** A higher colourfulness/balance score is not a more-appealing
  design and may be worse for the actual audience.[^reinecke14][^berlyne]
- **Treating tool output as measurement.** AIM and predictive-attention tools emit dataset-bound
  *predictions*; a "100" Lighthouse score is "passed automated checks," not "accessible."[^aim][^lhscoring][^afixt]
- **"Minimize complexity" / "maximize" any single dimension.** Preference is non-monotonic and
  population-relative; recommend *moderate*, audience-aware targets.[^berlyne][^reinecke14]
- **Comparing metric values across tools without checking convention.** Symmetry/balance polarity and
  normalization differ (e.g. "0% = balanced" vs "higher = more symmetric").[^toolbox]
- **Accessibility overlays.** They can block axe/Lighthouse from running and return zero errors —
  manufacturing a *worse* false sense of security than no testing.[^accessproof]

## Operating stance (adopt near-verbatim)

**DO:** use objective metrics to flag candidates and quantify within-lineage deltas; surface high-confidence
tails (clutter, contrast, gross imbalance); pair every metric with its known correlation strength + an
explicit caveat; feed metrics as *one input* to a human/LLM critique; state which population a threshold
reflects; respect non-monotonicity (recommend moderate complexity).

**DON'T:** rank designs as "objectively better" or auto-score quality from metrics alone; trust a metric's
name as its construct; optimize a single metric; treat AIM/saliency/Lighthouse output as measurement or
ground truth; generalize a published R² to a new product/audience/era without re-validation.

**One-line contract:** *Every computational aesthetic metric is a weak-to-moderate proxy (best-case
R² ≈ 0.48–0.49 web, ~0.32 mobile, with half-plus the variance unexplained), use it to flag outliers and
quantify deltas, always paired with its correlation strength and human/LLM judgment; never to rank designs
or auto-score quality.*

## References

[^hasler]: Hasler, D. & Süsstrunk, S., "Measuring Colourfulness in Natural Images," Proc. IS&T/SPIE Electronic Imaging 2003 (HVEI VIII), vol. 5007, pp. 87–95. DOI 10.1117/12.477378. https://infoscience.epfl.ch/record/33994 (paper)
[^pyimage]: Rosebrock, A. (PyImageSearch), "Computing image 'colorfulness' with OpenCV and Python" (implements rg=R−G, yb=½(R+G)−B, M=√(σ_rg²+σ_yb²)+0.3·√(μ_rg²+μ_yb²)). https://pyimagesearch.com/2017/06/05/computing-image-colorfulness-with-opencv-and-python/ (blog)
[^bvdart]: BVDART, "Dominant Color Extraction: A Practical Guide" (resize → quantize 8–16 colours → sort). https://bvdart.nl/en/articles/dominant-color-extraction-in-practice (blog)
[^cohenor]: Cohen-Or, D., Sorkine, O., Gal, R., Leyvand, T. & Xu, Y.-Q., "Color Harmonization," ACM TOG (SIGGRAPH) 25(3):624–630, 2006. https://igl.ethz.ch/projects/color-harmonization/harmonization.pdf (paper)
[^palmer]: Palmer, S. E. & Schloss, K. B., "An ecological valence theory of human color preference," PNAS 107:8877–8882, 2010 (WAVE). https://www.pnas.org/doi/10.1073/pnas.0906172107 (paper)
[^saraee]: Saraee, Jalal & Betke / Miniukovich-adjacent, "Data Compression Algorithms in Analysis of UI Layouts Visual Complexity" (entropy + edge density + JPEG ratio, r≈0.70). https://link.springer.com/chapter/10.1007/978-3-030-37487-7_14 (paper)
[^karjus]: Karjus, A. et al., "Compression ensembles quantify aesthetic complexity and the evolution of visual art," EPJ Data Science, 2023. https://epjdatascience.springeropen.com/articles/10.1140/epjds/s13688-023-00397-3 (paper)
[^quadtree]: York, T., "Quadtrees for Image Compression" (leaf-block decomposition explainer). https://medium.com/@tannerwyork/quadtrees-for-image-processing-302536c95c00 (blog)
[^rosenholtz]: Rosenholtz, R., Li, Y. & Nakano, L., "Measuring visual clutter," Journal of Vision 7(2):17, 2007. DOI 10.1167/7.2.17. https://pubmed.ncbi.nlm.nih.gov/18217832/ ; https://persci.mit.edu/research/clutter/ (paper)
[^clutterpy]: kargaranamir, `visual-clutter`: Python implementation of Feature Congestion & Subband Entropy (JoV 2007). https://github.com/kargaranamir/visual-clutter (docs)
[^toolbox]: "A toolbox for calculating objective image properties in aesthetics research," arXiv:2408.10616 (mirror symmetry %, Wilson & Chatterjee balance, edge/contour density, self-similarity; Ngo bounding-box distinction; Harrington white-space fraction & free-flow). https://arxiv.org/pdf/2408.10616 (paper)
[^miniukovich14]: Miniukovich, A. & De Angeli, A., "Quantification of Interface Visual Complexity," AVI 2014. https://dl.acm.org/doi/10.1145/2598153.2598173 (paper)
[^nordichi]: Miniukovich, A. & De Angeli, A., "Visual Impressions of Mobile App Interfaces," NordiCHI 2014. https://dl.acm.org/doi/pdf/10.1145/2639189.2641219 (paper)
[^miniukovich15]: Miniukovich, A. & De Angeli, A., "Computation of Interface Aesthetics," CHI 2015 (Honorable Mention; 8 metrics; up to 49% web / ~32% iPhone-app variance; N=62 web, 53 apps; 150 ms & 4 s). https://dl.acm.org/doi/10.1145/2702123.2702575 (paper)
[^reinecke13]: Reinecke, K., Yeh, T., Miratrix, L., Mardiko, R., Zhao, Y., Liu, J. & Gajos, K., "Predicting Users' First Impressions of Website Aesthetics with a Quantification of Perceived Visual Complexity and Colorfulness," CHI 2013 (450 sites, 548 participants; ~48%/500 ms; balance/symmetry dropped as n.s.). https://iis.seas.harvard.edu/papers/2013/reinecke13aesthetics.pdf (paper)
[^aim]: Oulasvirta, A., De Pascale, S., Koch, J., Langerak, T., Jokinen, J., Todi, K., Laine, M., et al., "Aalto Interface Metrics (AIM): A Service and Codebase for Computational GUI Evaluation," UIST 2018 Adjunct. DOI 10.1145/3266037.3266087. https://dl.acm.org/doi/abs/10.1145/3266037.3266087 (paper)
[^aimrepo]: aalto-ui/aim, GitHub repository (MIT; `aim_metrics/{colour_perception,perceptual_fluency,visual_guidance,accessibility}`; `aim_frontend/src/config/metrics.js` metric IDs + per-metric citations; `aim2_metrics`; v2.3 release; in-code symmetry/equilibrium caveats). https://github.com/aalto-ui/aim (docs/code)
[^aimportal]: AIM service + Aalto research-portal entry (four categories; "17 metrics"; Grid Quality = number of alignment lines). https://interfacemetrics.aalto.fi/ ; https://research.aalto.fi/en/publications/aalto-interface-metrics-aim-a-service-and-codebase-for-computatio/ (docs)
[^zen]: Zen, M. & Vanderdonckt, J., "Towards an Evaluation of Graphical User Interfaces Aesthetics based on Metrics," IEEE RCIS 2014 (7 geometric metrics; "no consensus on which measure/formula/interpretation"). https://dial.uclouvain.be/pr/boreal/object/boreal:152890 (paper)
[^ngo]: Ngo, D.C.L., Teo, L.S. & Byrne, J.G., "Modelling interface aesthetics," Information Sciences 152:25–46, 2003 (14 measures from element bounding boxes). https://www.sciencedirect.com/science/article/abs/pii/S0020025502004048 (paper)
[^rico]: Deka, B., Huang, Z., Franzen, C., Hibschman, J., Afergan, D., Li, Y., Nichols, J. & Kumar, R., "Rico: A Mobile App Dataset for Building Data-Driven Design Applications," UIST 2017 (~66k screens + view hierarchies; 64-d layout autoencoder). https://dl.acm.org/doi/pdf/10.1145/3126594.3126651 ; http://www.interactionmining.org/rico (paper/docs)
[^itti]: Itti, L., Koch, C. & Niebur, E., "A Model of Saliency-Based Visual Attention for Rapid Scene Analysis," IEEE TPAMI 20(11):1254–1259, 1998. http://ilab.usc.edu/publications/doc/Itti_etal98pami.pdf (paper)
[^gbvs]: Harel, J., Koch, C. & Perona, P., "Graph-Based Visual Saliency," NIPS 2006. http://papers.neurips.cc/paper/3095-graph-based-visual-saliency.pdf (paper)
[^bylinskii19]: Bylinskii, Z., Judd, T., Oliva, A., Torralba, A. & Durand, F., "What Do Different Evaluation Metrics Tell Us About Saliency Models?," IEEE TPAMI 41(3):740–757, 2019 (NSS/CC recommended; sAUC for center-bias). https://arxiv.org/abs/1604.03605 (paper)
[^transalnet]: Lou, J. et al., "TranSalNet," Neurocomputing 2022 / arXiv:2110.03593 (metric taxonomy context). https://arxiv.org/pdf/2110.03593 (paper)
[^salmetrics]: Saliency-metrics literature (location- vs distribution-based; sAUC center-bias; CC symmetric FP/FN; SIM sensitive to false negatives), arXiv:2502.05554. https://arxiv.org/pdf/2502.05554 (paper)
[^salicon]: Jiang, M., Huang, S., Duan, J. & Zhao, Q., "SALICON: Saliency in Context," CVPR 2015 (mouse-contingent crowdsourcing; early deep model). https://openaccess.thecvf.com/content_cvpr_2015/papers/Jiang_SALICON_Saliency_in_2015_CVPR_paper.pdf (paper)
[^borji]: Borji, A., "Saliency Prediction in the Deep Learning Era: Successes and Limitations," IEEE TPAMI / arXiv:1810.03716 (SALICON/DeepGaze II/ML-Net/SAM; no model wins every metric). https://arxiv.org/pdf/1810.03716 (paper)
[^bylinskii17]: Bylinskii, Z., Kim, N.W., O'Donovan, P., Alsheikh, S., Madan, S., Pfister, H., Durand, F., Russell, B. & Hertzmann, A., "Learning Visual Importance for Graphic Designs and Data Visualizations," UIST 2017. https://arxiv.org/abs/1708.02660 ; https://github.com/cvzoya/visimportance (paper)
[^bubbleview]: Kim, N.W., Bylinskii, Z., Borkin, M.A., Gajos, K.Z., Oliva, A., Durand, F. & Pfister, H., "BubbleView," ACM TOCHI 24(5), 2017 (clicks ≈ >75% of fixations w/ 10–15 participants; a proxy). https://arxiv.org/abs/1702.05150 (paper)
[^umsi]: Fosco, C., Casser, V., Bedi, A.K., O'Donovan, P., Hertzmann, A. & Bylinskii, Z., "Predicting Visual Importance Across Graphic Design Types," UIST 2020 (UMSI; Imp1k; Xception+ASPP). https://arxiv.org/abs/2008.02912 ; https://predimportance.mit.edu/ (paper/docs)
[^ueyes]: Jiang, Y., Leiva, L.A., Rezazadegan Tavakoli, H., R.G. Houssel, P., Kylmälä, J. & Oulasvirta, A., "UEyes: Understanding Visual Saliency across User Interface Types," CHI 2023 (62 participants, 1,980 UIs, real eye tracker). https://dl.acm.org/doi/10.1145/3544548.3581096 ; https://github.com/YueJiang-nj/UEyes-CHI2023 (paper)
[^agd]: "Predicting Visual Attention in Graphic Design Documents," arXiv:2407.02439, 2024 (natural-image attention models "do not generalize well to graphic design documents"; design-specific AGD model beats SAM-ResNet). https://arxiv.org/html/2407.02439v1 (paper)
[^attninsight]: Attention Insight, predictive eye-tracking technology & accuracy claims (vendor self-reported ~92.5–96%). https://attentioninsight.com/technology/ (blog/vendor)
[^neurons]: Neurons, "How Neurons Attention Prediction works" (vendor self-reported ~95–98%). https://knowledge.neuronsinc.com/how-neurons-attention-prediction (docs/vendor)
[^heatmapcrit]: Attention/heatmap-validity criticism (no standardized eval; task-blind), arXiv:2407.02484. https://arxiv.org/html/2407.02484v1 (paper)
[^axeapi]: Deque, Axe API Documentation (rule tags; impact levels; violations/passes/incomplete/inapplicable; rendered-content / zero-false-positive principle). https://www.deque.com/axe/core-documentation/api-documentation/ (docs)
[^axerepo]: dequelabs/axe-core, GitHub (engine overview; integration surfaces). https://github.com/dequelabs/axe-core (docs)
[^mabl]: mabl, "Accessibility rules and tags" (tag semantics; tags don't roll up; impact mapping). https://help.mabl.com/hc/en-us/articles/25101592214804 (docs/vendor)
[^deque57]: Deque, "Automated Testing Study Identifies 57 Percent of Digital Accessibility Issues" (57% by issue volume; 2,000+ audits / 13k+ pages; reframes the 20–30% belief). https://www.deque.com/blog/automated-testing-study-identifies-57-percent-of-digital-accessibility-issues/ (blog/vendor)
[^b13]: b13, "Why Automated Accessibility Testing Isn't Enough" (30–40% figure; rest needs human judgment). https://b13.com/blog/why-automated-accessibility-testing-isnt-enough (blog)
[^crosscheck]: Crosscheck, "axe vs WAVE vs Pa11y (2026)" (30–40% catch rate; context/semantics not evaluable). https://crosscheck.cloud/blogs/axe-vs-wave-vs-pa11y-accessibility-testing/ (blog)
[^lhscoring]: Google / Chrome for Developers, "Lighthouse accessibility scoring" (weighted average of axe audits; weighting by axe user impact; binary per audit; manual-testing encouraged; subset only). https://developer.chrome.com/docs/lighthouse/accessibility/scoring (docs)
[^accessproof]: AccessProof, "What is axe-core and Why It Powers Real Accessibility Audits" (zero-false-positive; can't judge alt-text meaning/reading order; overlay risk; WCAG-2.2-target-size only). https://access-proof.com/blog/what-is-axe-core-evidence-based-audits (blog)
[^afixt]: AFixt, "Why Your Lighthouse Score of 100 Means Almost Nothing" (100 = automated checks passed, not accessible). https://afixt.com/why-your-lighthouse-score-of-100-means-almost-nothing/ (blog)
[^a11ytest]: accessibility-test.org, "Lighthouse Accessibility Score, Insights and Limitations" ("No automated tool can guarantee conformance with WCAG"). https://accessibility-test.org/blog/testing-tools/lighthouse-accessibility-score-insights-and-limitations/ (blog)
[^g17]: W3C, "G17 / Contrast technique" (ratio (L1+0.05)/(L2+0.05); luminance linearization 0.03928/12.92/1.055/2.4; coefficients 0.2126/0.7152/0.0722; 7:1 AAA). https://www.w3.org/TR/WCAG20-TECHS/G17.html (docs)
[^bickford]: Bickford, N., "Computing WCAG Contrast Ratios" (independent confirmation of formula & constants). https://www.neilbickford.com/blog/2020/10/18/computing-wcag-contrast-ratios/ (blog)
[^contrastmin]: W3C/WAI, "Understanding SC 1.4.3 Contrast (Minimum)" (4.5:1 normal; 3:1 large; large = 18pt / 14pt bold ≈ 24px / 18.5px). https://www.w3.org/WAI/WCAG21/Understanding/contrast-minimum.html (docs)
[^nontext]: Deque University, "1.4.11 Non-Text Contrast (AA)" (3:1 for UI components & graphical objects). https://dequeuniversity.com/resources/wcag2.1/1.4.11-non-text-contrast (docs)
[^targetmin]: W3C / community, WCAG 2.2 SC 2.5.8 Target Size (Minimum) AA = 24×24 CSS px + exceptions; SC 2.5.5 Enhanced AAA = 44×44. https://www.w3.org/WAI/WCAG22/Understanding/target-size-minimum.html ; https://tetralogical.com/blog/2022/12/20/foundations-target-size/ (docs)
[^apca]: Myndex, SAPC-APCA (APCA = WCAG 3 candidate; Lc model; derived from SAPC). https://github.com/Myndex/SAPC-APCA (docs)
[^apca2]: APCA documentation + Dan Hollick, "WCAG 3 and APCA" (Lc ±106; contrast moved out of WCAG 3 WD July 2023, "TBD"). https://git.apcacontrast.com/documentation/README (docs/blog)
[^reinecke14]: Reinecke, K. & Gajos, K.Z., "Quantifying Visual Preferences Around the World," CHI 2014 (2.4M ratings, ~40k participants; preference peaks vary by country/gender/education). https://www.eecs.harvard.edu/~kgajos/papers/2014/reinecke14visual.pdf (paper)
[^berlyne]: Marin, M.M. et al., "Berlyne Revisited: Evidence for the Multifaceted Nature of Hedonic Tone," Frontiers in Human Neuroscience, 2016 (inverted-U; arousal potential). https://www.frontiersin.org/articles/10.3389/fnhum.2016.00536/full (paper)
[^althuizen]: Althuizen, N., "Revisiting Berlyne's inverted U-shape relationship between complexity and liking," Psychology & Marketing, 2021 (N>1,800; scant inverted-U evidence). https://onlinelibrary.wiley.com/doi/abs/10.1002/mar.21449 (paper)
[^tuch]: Tuch, A.N., Presslaber, E.E., Stöcklin, M., Opwis, K. & Bargas-Avila, J.A., "The role of visual complexity and prototypicality regarding first impression of websites," IJHCS 70(11), 2012 (17–50 ms; prototypicality co-driver). https://www.sciencedirect.com/science/article/abs/pii/S1071581912001127 (paper)
[^tractinsky]: Tractinsky, N., Katz, A.S. & Ikar, D., "What is beautiful is usable," Interacting with Computers 13:127–145, 2000. https://www.ise.bgu.ac.il/faculty/noam/papers/00_nt_ask_di_iwc.pdf (paper)
[^halo]: Sutcliffe-tradition review, "Is beautiful really usable?," Computers in Human Behavior, 2012 (hedonic halo). https://www.sciencedirect.com/science/article/abs/pii/S0747563212000908 (paper)
[^birkhoff]: Douchová, V., "Birkhoff's aesthetic measure" (empirical failure; Eysenck M=O×C). https://www.researchgate.net/publication/323296865 (paper)
[^birkhoff2]: Secondary summaries of Birkhoff's Aesthetic Measure (1933), M=O/C and limited empirical support. https://www.historymath.com/birkhoff/ (blog/secondary)
[^lavie]: Lavie, T. & Tractinsky, N., "Assessing dimensions of perceived visual aesthetics of web sites," IJHCS 60, 2004 (classical vs expressive factors). https://www.ise.bgu.ac.il/faculty/noam/papers/04_tl_nt_ijhcs.pdf (paper)
[^twosides]: "Objective and Subjective Measures of Visual Aesthetics of Website Interface Design" (no correlation between objective and subjective balance; symmetry/balance null vs perceived aesthetics). https://www.academia.edu/105248847/ (paper)
[^equilibrium]: "Judging whether it is aesthetic: Does equilibrium compensate for the lack of symmetry?," PMC. https://pmc.ncbi.nlm.nih.gov/articles/PMC3690416/ (paper)
