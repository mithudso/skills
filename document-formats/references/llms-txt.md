# llms.txt & Markdown-for-LLMs

The `llms.txt` convention — a **curated, markdown-formatted index a site publishes for LLM
consumers**, plus the broader practice of authoring markdown *for* language models (CLAUDE.md /
AGENTS.md, `.md` API endpoints). This is a markdown **file-format convention**, so it lives in
`document-formats`; the deeper *LLM context-engineering* discipline (what to put in a system prompt,
caching, retrieval) belongs to the `ai-mcp-sdk-prompting` hub — defer there for prompt/context design.

> **Honesty note (carry this into any recommendation):** `llms.txt` is a *proposed* convention, not a
> ratified standard, and **no major AI platform has committed to reading it as a first-class input**. A
> Nov-2025 SERanking study across ~300k domains found it does **not measurably improve AI citations**.
> Recommend it as low-cost, well-formed metadata — not as a proven ranking/visibility lever.

---

## 1. `llms.txt` — what it is

Proposed by **Jeremy Howard (Answer.AI), Sept 2024**. A markdown file at the site root (`/llms.txt`)
that gives an LLM a **hand-curated, priority-ordered map** of the site's high-value content. It is to
LLM retrieval what `robots.txt` is to crawlers — but *additive guidance*, not access control.

**Spec shape** (from llmstxt.org) — markdown with a required H1 + optional blockquote summary, then
H2-grouped link lists:

```markdown
# Project Name

> One-paragraph summary of what this project/site is.

Optional context paragraphs (no headings).

## Docs
- [Quickstart](https://example.com/quickstart.md): Get running in 5 minutes
- [API reference](https://example.com/api.md): Full endpoint list

## Optional
- [Changelog](https://example.com/changelog.md)
```

- The `## Optional` section is special: its links may be **skipped** when the consumer needs a shorter context.
- Links ideally point to **clean markdown** versions of pages (LLM-friendly, no nav chrome).

## 2. `llms.txt` vs `llms-full.txt` vs neighbors

| File | Purpose |
| --- | --- |
| `/llms.txt` | Curated **index** — links + one-line descriptions; small, an entry point. |
| `/llms-full.txt` | The **entire documentation inlined** into one big markdown file — paste-the-whole-thing context. Can be huge; watch token cost. |
| `robots.txt` | *Restricts* what crawlers may fetch (access control). Orthogonal to llms.txt. |
| `sitemap.xml` | Enumerates *every* URL for search indexing (exhaustive, machine-only). llms.txt is curated + human-readable. |

Common pattern: serve a markdown twin of each HTML page (e.g. content negotiation, or a `.md` suffix —
the "append `.md` to any docs URL" pattern several dev-tool sites now ship), and list the important
ones in `llms.txt`. Adopters include Anthropic, Stripe, Zapier, Cloudflare, and many dev-tool docs.

## 3. Authoring markdown *for* LLMs (the broader practice)

- **`CLAUDE.md` / `AGENTS.md`** — repo-root markdown that instructs coding agents (project conventions,
  build/test commands, do/don't). Keep terse, imperative, high-signal; agents read it as context every session.
- **Clean-markdown endpoints** — strip nav/ads/chrome; stable headings (they become anchors and
  retrieval boundaries); short, front-loaded sections (mitigates "lost in the middle").
- **Determinism & token economy** — prefer plain CommonMark (tables/code fences render predictably in
  LLM output); avoid renderer-specific syntax an LLM consumer won't interpret; budget `llms-full.txt` size.
- For *how an LLM uses* this context (caching, retrieval, injection hardening, context-window budgeting),
  that's context-engineering → **`ai-mcp-sdk-prompting`**, not this format reference.

## 4. Generating & validating

- Generators: docs frameworks increasingly emit `llms.txt`/`llms-full.txt` (Mintlify, Starlight plugins,
  Docusaurus plugins, `llmstxt`-style site plugins); or generate from your sitemap + markdown source.
- Validate: it's just markdown — lint with `markdownlint`/`remark-lint` (see `references/markdown-linting.md`)
  and confirm the H1 + link-list structure; check linked `.md` URLs resolve.

## Sources
- [llmstxt.org (the proposal/spec)](https://llmstxt.org/) · [Answer.AI announcement](https://www.answer.ai/posts/2024-09-03-llmstxt.html)
- [Is llms.txt Dead? Adoption in 2025 (llms-txt.io)](https://llms-txt.io/blog/is-llms-txt-dead) · [SERanking 300k-domain study coverage](https://codersera.com/blog/llms-txt-complete-guide-2026/)
- Related practice: `CLAUDE.md`/`AGENTS.md` agent-context conventions (see the `ai-mcp-sdk-prompting` hub for context engineering).
