# Pandoc — Universal Document Conversion (Markdown ↔ everything)

Pandoc is the swiss-army converter for markdown: **read one format, write another** through a single
internal document AST. This reference owns Pandoc and format conversion. For plain-markdown
syntax/flavors load `references/markdown-authoring.md`; for the JS remark/unified AST ecosystem load
`references/markdown-processing.md` (Pandoc is the Haskell-world counterpart, with its own AST).

---

## 1. The model: reader → AST → writer

Pandoc parses any input format into the **Pandoc AST** (a typed document tree), then renders it to any
output format. Conversion = pick a reader + a writer.

```bash
pandoc input.md -o output.docx              # format inferred from extensions
pandoc -f gfm -t html5 input.md -o out.html # explicit -f (from) / -t (to)
pandoc README.md -o readme.pdf              # md → PDF (needs a LaTeX engine, see §4)
```

- **Readers** include: `markdown` (pandoc's extended flavor), `commonmark`, `commonmark_x`, `gfm`,
  `html`, `latex`, `docx`, `odt`, `epub`, `rst`, `org`, `mediawiki`, `textile`, `ipynb`, `csv`, `bibtex`.
- **Writers** include all of the above plus `pdf`, `pptx`, `beamer` (slides), `revealjs`/`s5` (HTML slides),
  `man`, `asciidoc`, `docbook`, `jats`, `icml`, `typst`.

## 2. Pandoc's Markdown flavor

The richest markdown flavor — a superset of CommonMark with toggleable extensions via `+ext`/`-ext`:

```bash
pandoc -f markdown+hard_line_breaks-yaml_metadata_block input.md -t html
pandoc -t commonmark_x input.md             # CommonMark + pandoc's extension set
```

Notable extensions: `pipe_tables`, `grid_tables`, `multiline_tables`, `footnotes`, `definition_lists`,
`fenced_divs` (`::: warning … :::`), `bracketed_spans` (`[text]{.class}`), `attributes` (`# H {#id .cls}`),
`tex_math_dollars` (`$…$`), `citations` (`[@key]`), `yaml_metadata_block`. Use `pandoc --list-extensions`
to see all and their on/off defaults per format.

## 3. Filters, templates, metadata

- **Lua filters** (`--lua-filter=f.lua`) — transform the AST in-process (no extra runtime); the modern,
  fast default for custom transforms. **JSON filters** (`--filter=prog`) — any language that reads/writes
  the Pandoc JSON AST on stdin/stdout (e.g. `pandoc-crossref`).
- **Templates** (`--template=tpl`, `-V key=val`) — control the output skeleton; `pandoc -D html5` dumps
  the default template to customize. Partials + variables drive layout.
- **Metadata** — YAML frontmatter or `--metadata-file`; drives title/author/date and template variables.
- **Citations** — `--citeproc` + `--bibliography=refs.bib` + a CSL style (`--csl=style.csl`) renders
  `[@key]` citations and a reference list (replaces the old `pandoc-citeproc`).

## 4. Markdown → PDF (the common ask)

PDF output goes through an intermediate engine — pick with `--pdf-engine`:

| Engine | Path | Use when |
| --- | --- | --- |
| `pdflatex`/`xelatex`/`lualatex` | LaTeX (needs a TeX install, e.g. TinyTeX) | best typography; xelatex/lualatex for Unicode/fonts |
| `weasyprint` | HTML/CSS → PDF (Python) | style with CSS, no LaTeX |
| `wkhtmltopdf` | HTML/CSS → PDF (WebKit) | legacy HTML→PDF |
| `typst` | Typst (fast, modern) | newer, fast, no TeX |
| `context` | ConTeXt | advanced layout |

```bash
pandoc doc.md -o doc.pdf --pdf-engine=xelatex -V mainfont="Helvetica"
pandoc doc.md -o doc.pdf --pdf-engine=weasyprint --css=style.css
```

## 5. Handy recipes

```bash
pandoc -s in.md -o out.html                       # -s = standalone (full doc, not a fragment)
pandoc -s --toc --toc-depth=2 in.md -o out.html    # table of contents
pandoc in.md -o slides.pptx                         # markdown → PowerPoint (see references/pptx.md)
pandoc in.md -t revealjs -s -o slides.html          # markdown → HTML slide deck
pandoc in.docx -t gfm -o out.md                     # Word → GitHub markdown (great for migrations)
pandoc in.md --embed-resources --standalone -o self-contained.html   # inline CSS/images/JS
pandoc in.ipynb -o notebook.md                      # Jupyter → markdown
```

## 6. When Pandoc vs the JS ecosystem

- **Pandoc** — cross-format conversion (docx/odt/epub/latex/pptx), academic/publishing pipelines,
  citations, PDF/print, CLI/CI batch jobs. Heavier dependency (Haskell binary + maybe TeX).
- **remark/unified** (`references/markdown-processing.md`) — in-JS transforms, web pipelines, custom
  plugins, when you're already in a Node toolchain and outputting HTML.
- Both have a real AST; choose by ecosystem and target formats.

## Sources
- [pandoc.org](https://pandoc.org/) · [MANUAL](https://pandoc.org/MANUAL.html) · [Markdown flavor + extensions](https://pandoc.org/MANUAL.html#pandocs-markdown) · [Lua filters](https://pandoc.org/lua-filters.html) · [Citations / --citeproc](https://pandoc.org/MANUAL.html#citations) · [Creating a PDF](https://pandoc.org/MANUAL.html#creating-a-pdf)
