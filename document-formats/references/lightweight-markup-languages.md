# Lightweight Markup Languages (the markdown siblings: reST · AsciiDoc · Org · Textile · wiki)

Markdown's siblings — the other plain-text markup languages — and **when to choose one over markdown**.
This reference owns the comparison/selection of markup *languages*. For markdown's own syntax load
`references/markdown-authoring.md`; for converting between them load `references/markdown-pandoc.md`
(Pandoc reads/writes most of these).

---

## 1. Why markdown isn't always the answer

Markdown optimizes for *easy writing* and *web rendering*, but it is **deliberately minimal** — no
native cross-references, includes, attributes, or a rich semantic model. For large technical manuals,
single-source publishing, or structured authoring, a heavier lightweight-markup language is often better.

## 2. The siblings

| Language | Origin / ecosystem | Strengths | Weaknesses | Choose when |
| --- | --- | --- | --- | --- |
| **reStructuredText (reST)** | Python/Docutils, **Sphinx** | Rich directives & roles (`.. note::`, `:ref:`), cross-refs, includes, extensible; the Python-docs standard | Stricter/fiddlier syntax (indentation, underline-length for headings) | Python projects; API docs via Sphinx autodoc; large cross-referenced manuals |
| **AsciiDoc** | AsciiDoctor (Ruby/JVM/JS) | Powerful yet readable; includes, attributes, conditionals, admonitions, tables; great for books/manuals (DocBook/PDF) | Smaller ecosystem than markdown; another syntax to learn | Technical books, O'Reilly-style manuals, single-source to PDF/HTML/EPUB |
| **Org-mode** | Emacs | Outliner + docs + literate programming + agenda/tasks; babel code execution | Tied to Emacs for full power | Emacs users; literate programming; notes+tasks+docs in one |
| **Textile** | Older web CMSs | Compact inline syntax | Largely superseded by markdown | Maintaining legacy Textile content |
| **Wiki markups** (MediaWiki, Creole) | Wikis | Templates/transclusion; familiar to wiki editors | Per-engine dialects | Wiki platforms (MediaWiki etc.) |
| **MyST Markdown** | Sphinx/Jupyter | **Markdown with reST's power** — directives/roles in markdown syntax | Sphinx-centric | Want Sphinx/reST features but markdown ergonomics |
| **Typst** | Newer typesetting | Markup + programmable typesetting; fast; LaTeX alternative | Young ecosystem | Modern print/PDF typesetting without LaTeX |

## 3. Quick syntax contrasts

```rst
.. reStructuredText
Section Title
=============
.. note::
   An admonition with a directive.
See :ref:`other-section` for cross-reference.
```

```asciidoc
// AsciiDoc
= Document Title
== Section
[NOTE]
====
An admonition block.
====
<<other-section,See here>> for cross-reference.
:attribute: reusable value
```

```org
# Org-mode
* Top heading
** Sub heading
#+BEGIN_SRC python :results output
print("literate code that can execute")
#+END_SRC
```

Markdown equivalents (admonitions, cross-refs) are renderer-specific extensions, not core — which is
exactly the gap reST/AsciiDoc fill natively.

## 4. Selection heuristics

- **Default to markdown/GFM** for READMEs, web content, chat, LLM context, most project docs.
- **reST + Sphinx** for Python API docs and big cross-referenced manuals (or **MyST** to keep markdown ergonomics with Sphinx power).
- **AsciiDoc** for books/manuals and serious single-source → PDF/EPUB/HTML.
- **Org-mode** if you live in Emacs and want notes + tasks + literate docs unified.
- **Typst** for modern, programmable print/PDF without LaTeX.
- **Converting** between any of these → Pandoc (`references/markdown-pandoc.md`): e.g. `pandoc -f rst -t gfm`, `pandoc -f docbook -t asciidoc`.

## Sources
- [reStructuredText / Docutils](https://docutils.sourceforge.io/rst.html) · [Sphinx](https://www.sphinx-doc.org/) · [MyST Markdown](https://myst-parser.readthedocs.io/)
- [AsciiDoc / Asciidoctor](https://asciidoc.org/) · [Org mode](https://orgmode.org/) · [Typst](https://typst.app/docs/) · [Pandoc format support](https://pandoc.org/MANUAL.html#general-options)
