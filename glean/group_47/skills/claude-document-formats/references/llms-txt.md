# llms.txt & Markdown-for-LLMs

`llms.txt` — **curated, markdown-formatted index site publishes for LLM consumers**, plus broader practice of authoring markdown *for* language models (CLAUDE.md / AGENTS.md, `.md` API endpoints). Markdown **file-format convention** → lives in `document-formats`; deeper *LLM context-engineering* (system prompt content, caching, retrieval) → `ai-mcp-sdk-prompting` hub.

> **Honesty note (carry into any recommendation):** `llms.txt` *proposed* convention, not ratified standard. **No major AI platform committed to reading it as first-class input**. Nov-2025 SERanking study across ~300k domains: does **not measurably improve AI citations**. Recommend as low-cost, well-formed metadata — not proven ranking/visibility lever.

---

## 1. `llms.txt` — what it is

Proposed by **Jeremy Howard (Answer.AI), Sept 2024**. Markdown file at site root (`/llms.txt`) giving LLM **hand-curated, priority-ordered map** of high-value content. To LLM retrieval what `robots.txt` is to crawlers — but *additive guidance*, not access control.

**Spec shape** (from llmstxt.org) — markdown with required H1 + optional blockquote summary, then H2-grouped link lists:

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

- `## Optional` section special: links may be **skipped** when consumer needs shorter context.
- Links ideally point to **clean markdown** versions (LLM-friendly, no nav chrome).

## 2. `llms.txt` vs `llms-full.txt` vs neighbors

| File | Purpose |
| --- | --- |
| `/llms.txt` | Curated **index** — links + one-line descriptions; small, entry point. |
| `/llms-full.txt` | **Entire docs inlined** into one big markdown file — paste-the-whole-thing context. Can be huge; watch token cost. |
| `robots.txt` | *Restricts* crawler access (access control). Orthogonal to llms.txt. |
| `sitemap.xml` | Enumerates *every* URL for search (exhaustive, machine-only). llms.txt is curated + human-readable. |

Common pattern: serve markdown twin of each HTML page (content negotiation, or `.md` suffix — "append `.md` to any docs URL" pattern several dev-tool sites ship), list important ones in `llms.txt`. Adopters: Anthropic, Stripe, Zapier, Cloudflare, many dev-tool docs.

## 3. Authoring markdown *for* LLMs (broader practice)

- **`CLAUDE.md` / `AGENTS.md`** — repo-root markdown instructing coding agents (project conventions, build/test commands, do/don't). Keep terse, imperative, high-signal; agents read every session.
- **Clean-markdown endpoints** — strip nav/ads/chrome; stable headings (become anchors + retrieval boundaries); short, front-loaded sections (mitigates "lost in the middle").
- **Determinism & token economy** — prefer plain CommonMark (tables/code fences render predictably); avoid renderer-specific syntax LLM consumer won't interpret; budget `llms-full.txt` size.
- *How LLM uses context* (caching, retrieval, injection hardening, context-window budgeting) → **`ai-mcp-sdk-prompting`**, not this format reference.

## 4. Generating & validating

- Generators: docs frameworks increasingly emit `llms.txt`/`llms-full.txt` (Mintlify, Starlight plugins, Docusaurus plugins, `llmstxt`-style site plugins); or generate from sitemap + markdown source.
- Validate: just markdown — lint with `markdownlint`/`remark-lint` (see `references/markdown-linting.md`), confirm H1 + link-list structure; check linked `.md` URLs resolve.

## Sources
- [llmstxt.org (the proposal/spec)](https://llmstxt.org/) · [Answer.AI announcement](https://www.answer.ai/posts/2024-09-03-llmstxt.html)
- [Is llms.txt Dead? Adoption in 2025 (llms-txt.io)](https://llms-txt.io/blog/is-llms-txt-dead) · [SERanking 300k-domain study coverage](https://codersera.com/blog/llms-txt-complete-guide-2026/)
- Related practice: `CLAUDE.md`/`AGENTS.md` agent-context conventions (see `ai-mcp-sdk-prompting` hub for context engineering).