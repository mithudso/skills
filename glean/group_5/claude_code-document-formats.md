# document-formats

**Category:** Databases, Data Engineering & Analytics
**Platform:** Claude Code
**Original Path:** claude-code/document-formats

## Description
Generate/parse/edit/convert document & data-file formats in Python & Node — PDF, Word(.docx), Excel(.xlsx), PowerPoint(.pptx), CSV/TSV, advanced JSON, draw.io(.drawio), Markdown. TRIGGER: PDF (pdf-lib, pypdf, ReportLab, HTML-to-PDF, extract text/tables); Word .docx (docx-js, TOC, tracked changes); Excel .xlsx (openpyxl/pandas, formulas, charts); PowerPoint .pptx (PptxGenJS/python-pptx, pptx→PDF); CSV/TSV (encoding/BOM, formula-injection CWE-1236, csvkit/qsv/DuckDB); advanced JSON (streaming, JSON Schema/Ajv/Zod, JSON Patch, JSONPath, MessagePack); draw.io (mxGraph XML, drawpyo); Markdown/CommonMark/GFM authoring + processing (remark/unified, markdown-it); reports/invoices/memos/decks/diagrams as files. SKIP: analytical/ETL tabular processing → da-analytical-methods / da-data-engineering-platform; extracting FROM aging/live-DOM sources → content-ingestion-extraction; in-browser markdown render+sanitize → markdown-rendering-browser.

---

# Document & File Formats

Hub for **programmatic document and data-file work** — creating, parsing, editing, and converting
the common office and data formats in Python and Node.js. Each former standalone format skill is now
an on-demand reference file: when a task matches a row in the routing table below, **`Read` that
`references/<name>.md` file before answering** — the routing descriptions are a dispatch index, not the
depth itself.

The boundary that defines this hub: it owns the **file format** — bytes in, bytes out, and the
libraries that manipulate them. When the real question is the *analysis* of the data, the *extraction*
of content from messy sources, or the *prose quality* of a written document, defer to the sibling hubs
named in the cross-hub map.

<!-- ROUTING TABLE: document-formats — auto-generated, edit descriptions as needed -->
## Sub-skill routing table

This hub provides 15 on-demand reference files (7 absorbed format skills + 8 Markdown/markup references). When a task matches a row, **Read the listed `references/` file** before answering — do not rely on this table alone for depth.

| Sub-topic | When to load | Reference file |
| --- | --- | --- |
| `pdf` | Create, parse, merge, split, watermark, encrypt, sign, validate, convert PDF (Python + Node); HTML→PDF; extract text/tables; AcroForms; PDF/A·PDF/UA | `references/pdf.md` |
| `docx` | Create, read, edit, manipulate Word .docx (docx-js + XML); TOC/headings/letterheads; tracked changes; comments; find-and-replace; images | `references/docx.md` |
| `xlsx` | Create, read, edit, fix Excel .xlsx (openpyxl + pandas); formulas, formatting, charts; clean/restructure messy tabular data | `references/xlsx.md` |
| `pptx` | Create, read, edit PowerPoint .pptx (PptxGenJS + python-pptx); slide masters/layouts/templates; charts/tables; pptx→PDF | `references/pptx.md` |
| `csv` | Parse, generate, validate, convert, stream CSV/TSV; encoding (BOM/UTF-8/1252); formula-injection (CWE-1236); csvkit/qsv/miller/DuckDB | `references/csv.md` |
| `json-advanced` | Streaming parsers; JSON Schema (Ajv/Zod/TypeBox); JSON Patch (RFC 6902); JSONPath; MessagePack/CBOR/BSON; NDJSON; JSON5/JSONC | `references/json-advanced.md` |
| `drawio-diagrams` | Programmatic .drawio creation/parse/transform; mxGraphModel/mxCell XML; export SVG/PNG/PDF; drawpyo/maxGraph; CI/CD diagram gen | `references/drawio-diagrams.md` |
| `markdown-authoring` | Write correct/portable Markdown; CommonMark 0.31.2 vs GFM vs Pandoc/Obsidian/MDX flavors; core syntax + GFM extensions (tables, task lists, footnotes, strikethrough), frontmatter (YAML/TOML), GitHub alerts `> [!NOTE]`, math `$…$`; portability cheat-sheet | `references/markdown-authoring.md` |
| `markdown-processing` | Parse/transform/analyze Markdown in code; choose marked vs markdown-it vs micromark vs unified/remark/rehype; the unist/mdast↔hast model, unist-util-visit, write remark/rehype plugins; md↔HTML, sanitize, recipes | `references/markdown-processing.md` |
| `mdx` | Markdown + JSX components; MDX 3 (ES2024, top-level await, block expressions); compile/evaluate, components prop / MDXProvider; Docusaurus/Astro/Next/Storybook; untrusted-input danger | `references/mdx.md` |
| `llms-txt` | `llms.txt` / `llms-full.txt` curated markdown index for LLM consumers; AGENTS.md/AGENTS.md agent context; clean-markdown endpoints; honest adoption caveats | `references/llms-txt.md` |
| `markdown-pandoc` | Pandoc universal conversion (md↔docx/PDF/HTML/epub/pptx/reST/…); reader→AST→writer; Lua/JSON filters; templates; citations (--citeproc); md→PDF engines | `references/markdown-pandoc.md` |
| `markdown-docs-as-code` | Docs-as-code workflow; static-site-generator selection (Starlight/Astro, Docusaurus, MkDocs Material, Hugo, VitePress, Eleventy, Jekyll, Sphinx); CI gates, preview deploys | `references/markdown-docs-as-code.md` |
| `lightweight-markup` | Markdown siblings & when to pick them: reStructuredText/Sphinx, AsciiDoc, Org-mode, MyST, Typst, Textile/wiki; selection heuristics; convert via Pandoc | `references/lightweight-markup-languages.md` |
| `markdown-linting` | Markdown quality gates: markdownlint(-cli2) rules/config, remark-lint presets, Vale prose lint, lychee link-check; pre-commit/CI integration | `references/markdown-linting.md` |

<!-- cross-hub-map -->
## Cross-hub map — where every document-formats topic lives

This family is split across these hubs. If a task's deep material is **not** in this hub's Sub-skill
routing table, it is a reference file under a sibling hub below — **activate that hub or `Read` its
`references/<name>.md` directly**. Every former standalone skill in this family is now a reference under one
of these hubs (nothing was deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `document-formats` | Document & File Formats (PDF, Word, Excel, PowerPoint, CSV, JSON, draw.io, Markdown) | `references/pdf.md`, `references/docx.md`, `references/xlsx.md`, `references/pptx.md`, `references/markdown-authoring.md`, `references/markdown-processing.md`, … |