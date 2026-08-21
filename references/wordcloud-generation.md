<!-- hub-reference-banner -->
> **Reference file — part of the `misc-catch-all` hub.** Formerly the standalone `wordcloud-generation` skill.
> Sibling topics in this family are now reference files under the hubs (`misc-catch-all`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: wordcloud-generation
version: "1.0.0"
updated: "2026-06-19"
description: >-
  TRIGGER: building or rendering a wordcloud/tag cloud; turning text into weighted terms for a cloud (tokenize, stopwords, frequency vs TF-IDF, bigram clouds, max-words/min-count); the Wordle/d3-cloud spiral-placement + sprite-collision layout; font-size scaling (linear/sqrt/log); generation toolchains (Python wordcloud, d3-cloud, wordcloud2.js, R ggwordcloud, HTML/CSS clouds); shape/image-masked or semantic clouds; or deciding whether a wordcloud is the right chart vs a frequency bar chart. SKIP: TF-IDF math/theory and deep NLP/topic-modeling → da-36-text-analytics-nlp (summarize the weighting choice and cross-ref); general chart-selection theory → da-8-data-visualization.
---

# Wordclouds & Tag Clouds: Generation and Rendering

A wordcloud (a.k.a. tag cloud) maps a set of weighted terms to text glyphs whose **font size** encodes weight, packed into a 2D region. It is a *decorative summary* of which terms are prominent — not a precision chart. This skill covers the full pipeline (text → weighted terms → layout → render), the major toolchains with correct code, advanced masking/semantic variants, and — critically — the well-documented reasons wordclouds mislead and when to reach for a bar chart instead.

The lineage: Jonathan Feinberg's **Wordle** (2008, IBM; Viégas-Wattenberg-Feinberg, *Participatory Visualization with Wordle*, IEEE TVCG 2009) defined the aesthetic; Jason Davies' **d3-cloud** is the canonical open-source reimplementation of that layout. (TF-IDF math/theory and deep NLP belong to **da-36-text-analytics-nlp**; chart-selection theory to **da-8-data-visualization**.)

## Core concepts

### 1. From text to weighted terms (the input pipeline)

A cloud is only as good as the `{term: weight}` map you feed it. Steps, in order:

1. **Tokenize** — split text into tokens. Most libraries use a regex; Python `wordcloud` defaults to `regexp=r"\w[\w']+"` (drops single chars). For meaningful clouds, tokenize on word boundaries and decide whether to keep numbers (`include_numbers`) and a `min_word_length`.
2. **Casefold** — lowercase so `The`/`the` merge. Optionally **lemmatize/stem** (`running`→`run`) to collapse inflections; `wordcloud` does a lightweight `normalize_plurals=True` (merges trailing-`s`), not true lemmatization. Real lemmatization needs spaCy/NLTK upstream.
3. **Remove stopwords** — strip high-frequency, low-signal words (`the, and, of`). `wordcloud` ships a built-in `STOPWORDS` set; you almost always extend it with domain terms. **Caution:** stopword removal is also a *source of lies* — Jacob Harris' classic example showed a cloud where the largest words were `like`/`policy` only because `don't` was a stopword, inverting the actual sentiment.
4. **Weight the terms.** Two main choices:
   - **Raw / relative frequency** — simplest; what Wordle and most tools do by default. Biases toward common-but-uninformative words.
   - **TF-IDF** — down-weights terms common across a corpus, surfacing distinctive terms. This is usually the *better* analytical weighting for comparing documents. (TF-IDF theory/computation is owned by **da-36-text-analytics-nlp** — cross-ref it; here the practical rule is: pass TF-IDF scores as the `frequencies` dict to `generate_from_frequencies`.)
5. **N-grams / collocations** — single tokens lose phrases. `wordcloud` sets `collocations=True` by default, detecting **bigrams** above a `collocation_threshold` (default 30, a likelihood score) so `"machine learning"` can appear as one entry. This directly answers Harris' "yes we can"/"car bomb" critique that unigram clouds destroy phrase meaning. Note: collocations are **ignored** when you call `generate_from_frequencies`.
6. **Threshold** — cap with `max_words` (default 200) and a `min_count`/min-frequency floor so the long tail of one-off terms doesn't clutter the canvas.

### 2. The layout algorithm (Wordle / d3-cloud)

Feinberg's own description (verbatim shape): count words, drop boring ones, sort by count descending, keep top N, assign **font size proportional to count**, then place greedily largest-first:

```
for each word in decreasing frequency order:
    place the word where it "wants" to be (near center / a random start)
    while it intersects any previously-placed word:
        move it one step along an ever-increasing spiral
```

Key mechanics:

- **Archimedean spiral.** Each candidate position steps outward along a spiral from the start point; d3-cloud supports `"archimedean"` (default) and `"rectangular"` spirals, or a custom polar generator `function(size){ return function(t){ return [x,y]; }; }`.
- **Greedy & order-dependent.** Big words get prime central real estate; small words fill gaps. Different random seeds → different layouts (set a seed for reproducibility).
- **Collision detection must be fast.** Naive O(n²) bounding-box tests are too slow. The tricks (per Feinberg & d3-cloud):
  - **Sprite/bitmap masks** — d3-cloud renders each word once to an off-screen HTML5 canvas, reads back its pixels, and packs them into a 32-bit-per-row **sprite bitmask**. Collision becomes a cheap bitwise-AND of pixel rows, far cheaper than pixel-by-pixel.
  - **Bounding-box prefilter** — only do the expensive sprite test when axis-aligned bounding boxes overlap.
  - Feinberg additionally used **hierarchical bounding boxes, a quadtree spatial index, and last-hit caching**.
- **Dropped words.** If a word can't be placed anywhere along the spiral it is **silently omitted** from the layout (d3-cloud documents this explicitly) — a subtle data-integrity hazard.
- **Rotation.** `prefer_horizontal` (wordcloud, default 0.9 = 90% horizontal) or an explicit angle set / `rotate` callback (d3-cloud) gives the playful mixed-orientation look. Rotation hurts readability, so most "serious" clouds keep words horizontal.
- **Font-size scaling function** — maps weight → px. Options:
  - **Linear** — Wordle's choice; Feinberg/Viégas noted linear "seemed more dynamic." Exaggerates the largest term.
  - **Square-root** — traditional tag clouds (and `wordcloud`'s `relative_scaling`) use sqrt so that *area* tracks frequency better, since a doubled font roughly quadruples ink.
  - **Log** — compresses heavy-tailed Zipfian frequency so mid-frequency words stay legible.
  - `wordcloud`'s `relative_scaling` (0–1) blends rank-only vs frequency-proportional sizing; `0` = size by rank only, `1` = strictly proportional, `"auto"` picks based on `repeat`.

### 3. Generation toolchains

**Python `wordcloud` (amueller/word_cloud)** — the de-facto Python tool. PIL-based, supports masks, colormaps, image-coloring.

```python
from wordcloud import WordCloud, STOPWORDS
import matplotlib.pyplot as plt

stop = set(STOPWORDS) | {"said", "company"}      # extend built-ins

# A) straight from raw text (does its own tokenize/stopword/collocation work)
wc = WordCloud(width=800, height=400, background_color="white",
               stopwords=stop, max_words=200, min_word_length=3,
               collocations=True, collocation_threshold=30,
               colormap="viridis", prefer_horizontal=0.9,
               relative_scaling=0.5, random_state=42).generate(text)

# B) from precomputed weights (e.g. TF-IDF) — bypasses tokenizer & stopwords
weights = {"mongodb": 0.42, "atlas": 0.31, "sharding": 0.18}   # term -> score
wc = WordCloud(background_color="white").generate_from_frequencies(weights)

plt.imshow(wc, interpolation="bilinear"); plt.axis("off"); plt.show()
wc.to_file("cloud.png")
```

**d3-cloud (JS, Jason Davies)** — computes a layout (positions/sizes/rotations); *you* render with D3/SVG. Asynchronous; you draw in the `end` callback.

```js
import cloud from "d3-cloud";

const layout = cloud()
  .size([800, 400])
  .words(data.map(d => ({ text: d.term, size: 10 + d.freq * 90 })))  // your scale
  .padding(2)
  .rotate(() => (~~(Math.random() * 2)) * 90)   // 0 or 90 deg
  .spiral("archimedean")
  .font("Impact")
  .fontSize(d => d.size)
  .on("end", words => {
    d3.select("svg").append("g")
      .attr("transform", "translate(400,200)")
      .selectAll("text").data(words).join("text")
        .style("font-size", d => d.size + "px")
        .attr("text-anchor", "middle")
        .attr("transform", d => `translate(${d.x},${d.y}) rotate(${d.rotate})`)
        .text(d => d.text);
  });
layout.start();
```
Gotcha: web fonts must be **loaded before** `start()` or bounding boxes are computed for the wrong glyphs.

**wordcloud2.js (timdream)** — renders directly to a canvas/DOM; takes `[word, size]` pairs and supports polar **shapes** and image masks. It draws each word to a scratch canvas, reads pixels to find free space.

```js
WordCloud(document.getElementById('canvas'), {
  list: [['foo', 12], ['bar', 6]],
  weightFactor: size => Math.pow(size, 1.2) * 4,   // your scaling fn
  gridSize: 8,            // larger grid = bigger gaps between words
  shape: 'cardioid',      // circle|cardioid|diamond|square|triangle|pentagon|star, or a polar fn
  rotateRatio: 0.5,
  color: 'random-dark',
  drawOutOfBound: false
});
```

**R — `wordcloud` and `ggwordcloud`.** `wordcloud::wordcloud(words, freq, ...)` is the classic. `ggwordcloud` (Le Pennec) is a ggplot2 geom whose C++ placement is a hybrid of `wordcloud` + wordcloud2.js, spiral-around-spawn with `geom_text_repel`-style non-overlap:

```r
library(ggwordcloud)
ggplot(df, aes(label = word, size = freq)) +
  geom_text_wordcloud_area(shape = "circle", rm_outside = TRUE) +  # _area: area proportional to size
  scale_size_area(max_size = 24) + theme_minimal()
```
`geom_text_wordcloud` makes font size proportional to `size`; `geom_text_wordcloud_area` makes the printed **area** proportional to `size` (the perceptually fairer default, since it corrects for long-words-look-bigger). `eccentricity` (default 0.65) stretches the spiral horizontally.

**HTML/CSS tag clouds** — for a small, accessible, link-bearing cloud (navigation tags), skip canvas entirely: emit `<a>` elements with `font-size` (and optional `color`) set from weight buckets. This is the most accessible form (real text, real links, screen-reader friendly) but offers no packing — words flow as inline text.

### 4. Advanced: masks, image-coloring, semantic positioning, color

- **Shape/image masks.** Constrain the cloud to a silhouette. In `wordcloud`, pass a `mask` ndarray; **white (#FFFFFF) pixels = "do not draw"**, everything else is paintable; `contour_width`/`contour_color` outline the shape:
  ```python
  import numpy as np; from PIL import Image
  mask = np.array(Image.open("shape.png"))
  wc = WordCloud(mask=mask, background_color="white",
                 contour_width=3, contour_color="steelblue").generate(text)
  ```
  wordcloud2.js and `ggwordcloud` take a mask image similarly.
- **Image-based coloring.** `ImageColorGenerator` colors each word by the **mean color of the source image under that word's bounding box**, so the cloud "paints" a picture:
  ```python
  from wordcloud import ImageColorGenerator
  img = np.array(Image.open("color_source.png"))
  wc = WordCloud(mask=img, background_color="white").generate(text)
  wc.recolor(color_func=ImageColorGenerator(img))
  ```
- **Custom color via `color_func`** — a callable `(word, font_size, position, orientation, font_path, random_state) -> color`; overrides `colormap`. Single color: `color_func=lambda *a, **k: "white"`. Encoding a *second* variable (e.g. sentiment) in color is one of the few ways a cloud carries more than one dimension — but color is itself a weak perceptual channel for magnitude.
- **Semantic / embedding-positioned clouds.** Standard layouts place words by *packing convenience*, so **position is meaningless**. Semantic clouds instead position related words near each other — e.g., embed terms (word2vec/transformer), reduce to 2D (t-SNE/UMAP), then place. Grouped/"organized" layouts with white-space zones of meaning have been shown more effective for analytic/topic-summarization tasks (NSF/IBM study) — at the cost of the playful Wordle look.

### 5. The critique — when NOT to use a wordcloud

Wordclouds are among the most-criticized visualizations. The canonical case is Jacob Harris, *"Word Clouds Considered Harmful"* (2011). Known perceptual failures:

1. **Size encoding is ambiguous.** Is weight encoded by font height, glyph width, or bounding-box **area**? Different fonts and different word *lengths* produce different visual weight for the *same* frequency. Readers can't reliably decode it (SCITEPRESS 2022 study: even area-corrected clouds gave wildly inaccurate, high-variance frequency estimates).
2. **Longer words look bigger.** A long word at the same font size occupies more area/ink than a short one, inflating its apparent importance — a pure length artifact.
3. **Frequency ≠ importance.** The most frequent word is rarely the most meaningful (Obama-speech cloud: "America" largest; the actual message "yes we can" — a bigram — invisible to a unigram cloud).
4. **No comparison, no quantities.** You can't tell if the top term is 2× or 10× the next, and comparing two clouds is near-impossible. Controlled studies (IBM/NSF) found Wordles **worse than column/row layouts** for word lookup and frequency comparison.
5. **Stopword/preprocessing sensitivity & no context.** Removing a word (or a different tokenizer) can flip the whole picture; isolated words lose the sentence context that gives them meaning.
6. **Confirmation-bias magnet.** Ask three people for the "top themes" of a cloud and you get three answers.

**Decision rule — prefer a frequency bar chart (or treemap/dot plot) when:** the goal is to *rank, compare, or quantify* term frequencies, or support a decision. **A wordcloud is acceptable when:** the purpose is a low-stakes, eye-catching *gist/decoration* (a banner, an icebreaker, an at-a-glance "what's this corpus roughly about"), the audience won't make decisions from it, and you've at least used sqrt/area scaling, a horizontal layout, and honest stopwords. Even then, pair it with the underlying counts. Better text alternatives: ranked **bar chart of frequencies** (the standard fix), dot/lollipop plot, treemap, or grouped/semantic layouts for topic summarization.

## Quick reference

| Need | Reach for |
|---|---|
| Python, masks, image-color | `wordcloud` (amueller) |
| Interactive web SVG, custom spiral | d3-cloud |
| Web canvas with polar shapes/masks | wordcloud2.js |
| ggplot2 / reproducible R, area-fair | `ggwordcloud` (`_area`) |
| Accessible tag/nav cloud with links | HTML/CSS `<a>` + font-size |
| Rank/compare frequencies | **bar chart, not a cloud** |

## Sources

- jasondavies/d3-cloud (Wordle-inspired layout, sprite masks, spiral): https://github.com/jasondavies/d3-cloud
- Feinberg: how Wordle actually works (SO answer): https://stackoverflow.com/questions/342687/algorithm-to-implement-a-word-cloud-like-wordle
- Viégas, Wattenberg & Feinberg, 'Participatory Visualization with Wordle' (IEEE TVCG 2009): http://hint.fm/papers/wordle_final2.pdf
- wordcloud (amueller) WordCloud parameters & generate_from_frequencies: https://amueller.github.io/word_cloud/generated/wordcloud.WordCloud.html
- wordcloud masked & image-colored examples: https://amueller.github.io/word_cloud/auto_examples/colored.html
- wordcloud2.js API (timdream): https://github.com/timdream/wordcloud2.js/blob/HEAD/API.md
- ggwordcloud (Le Pennec): https://lepennec.github.io/ggwordcloud/articles/ggwordcloud.html
- Jacob Harris, 'Word Clouds Considered Harmful' (2011): https://jacobharr.is/published/word-clouds
- SCITEPRESS 2022 — empirical accuracy study of word clouds: http://www.scitepress.org/publishedPapers/2022/107787/pdf/index.html
- NSF/IBM study — Wordles vs column/row layouts: https://par.nsf.gov/servlets/purl/10196691
- Displayr — 7 alternatives to word clouds: https://www.displayr.com/alternatives-word-cloud/
- Thematic — the 5 major faults of word clouds: https://getthematic.com/insights/word-clouds-and-their-5-major-faults
