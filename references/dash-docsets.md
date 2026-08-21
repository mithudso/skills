<!-- hub-reference-banner -->
> **Reference file — part of the `misc-catch-all` hub.** Formerly the standalone `dash-docsets` skill.
> Sibling topics in this family are now reference files under the hubs (`misc-catch-all`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: dash-docsets
description: >-
  Dash (Kapeli) docset ecosystem — the offline API-documentation format and its tooling.
  Covers the .docset bundle anatomy (Info.plist, docSet.dsidx SQLite searchIndex, ~75 entry
  types, dashAnchor TOC), generating docsets (doc2dash for Sphinx/MkDocs, dashing CSS-selector
  mapping, scrape-to-docset pipelines), publishing (Dash-User-Contributions PR flow, docset.json,
  feed XML, dash-feed://), integration surfaces (dash:// scheme, editor plugins, dasht CLI,
  Zeal/Velocity), and the Dash 8 AI surfaces (native MCP support, local HTTP API,
  Kapeli/dash-mcp-server install for Claude). TRIGGER: "create a docset", "doc2dash or dashing",
  "docset format / dsidx schema", "publish to Dash User Contributions", "hook Dash up to
  Claude/MCP", "offline docs for AI agents", "Zeal docsets". SKIP: MongoDB Manual lookups via the
  local docset → mongodb-docset-lookup; Plotly Dash (Python dashboard framework) →
  da-8-data-visualization / frontend-ui; general web scraping → content-ingestion-extraction.
origin: local
version: "1.0.0"
category: developer
updated: 2026-06-10
whenToUse:
  - "build or generate a Dash docset from HTML or Sphinx/MkDocs docs"
  - "docSet.dsidx searchIndex schema and entry types"
  - "publish a docset to Dash User Contributions"
  - "connect Dash docsets to Claude via MCP"
  - "choose between Dash, Zeal, DevDocs for offline docs"
related_skills:
  - mongodb-expert
  - content-ingestion-extraction
tags:
  - developer
  - documentation
  - tooling
---

# Dash Docsets

The Dash (Kapeli) docset ecosystem: format internals, generation tooling, publishing, integrations, and the 2025–2026 AI/MCP second act. Confidence tags: [HIGH] = 3+ independent sources or official primary; [MEDIUM] = thinner sourcing. Research date 2026-06-10. **Disambiguation:** this is Kapeli's macOS documentation browser and its portable docset format — NOT Plotly Dash and NOT the dash shell.

## Overview

Dash is an API Documentation Browser and Code Snippet Manager for macOS (solo developer Bogdan Popescu): instant offline fuzzy search over 200+ first-party docsets, 100+ cheat sheets, user-contributed docsets, and on-demand registry-generated docsets (PyPI/Sphinx, docs.rs, pkg.go.dev, hexdocs, Maven, RubyGems, Swift Package Index, etc.). The format is deliberately portable — **Zeal** (Windows/Linux, alive: v0.8.x releases Feb–Apr 2026) and **Velocity** (Windows, dormant) consume the same docsets; DevDocs.io is the browser-native alternative (39k stars, "searching for maintainers"). **Dash 8 (2025-09-13) added native MCP support and a local JSON HTTP API** — the ecosystem's positioning inverted from "AI makes offline docs obsolete" to "agents need grounded, version-pinned, local docs" (<10ms local lookups vs 100–500ms cloud doc APIs, no rate limits, privacy). Licensing moved to subscription with Dash 7 (2023); current exact pricing unverified [MEDIUM]. Net assessment [MEDIUM]: a stable niche with an AI-driven second act — flat as a human browsing habit, growing as machine-readable retrieval infrastructure.

## Core Concepts

### Docset bundle anatomy [HIGH]

```
<Name>.docset/
  icon.png / icon@2x.png        (16/32px, optional)
  Contents/
    Info.plist
    Resources/
      docSet.dsidx              (SQLite 3 search index)
      Documents/                (the HTML)
```

Search index — one table:

```sql
CREATE TABLE searchIndex(id INTEGER PRIMARY KEY, name TEXT, type TEXT, path TEXT);
CREATE UNIQUE INDEX anchor ON searchIndex (name, type, path);
INSERT OR IGNORE INTO searchIndex(name, type, path) VALUES ('printf', 'Function', 'stdio.html#printf');
```

`type` must be one of ~75 fixed entry types (Class, Function, Method, Guide, Command, Option, Operator, Section, … — new types by request to Kapeli). `path` is relative into Documents/, may carry `#anchor` or even be an http URL. The format descends from Apple's Xcode docsets — hence the TOC anchor syntax:

```html
<a name="//apple_ref/cpp/<EntryType>/<percent-encoded name>" class="dashAnchor"></a>
```
(requires `DashDocSetFamily = dashtoc` in Info.plist).

Key Info.plist keys: `dashIndexFilePath` (landing page), `DashDocSetFallbackURL` (online-page redirection — keep `searchIndex.path` congruent with the online URL structure; extensionless docs sites need extensionless files), `isJavaScriptEnabled` (external JS off by default), `DashDocSetDefaultFTSEnabled`, `DashDocSetPlayURL`.

### Generation tooling [HIGH]

| Source | Tool | Notes |
|---|---|---|
| Sphinx / MkDocs+mkdocstrings / pydoctor | **doc2dash** (Python, Hynek Schlawack) | Canonical for intersphinx; extensible parsers; active (deps/CI current 2026) |
| Arbitrary clean HTML | **dashing** (Go) | `dashing.json` maps CSS selectors → entry types; maintenance mode since 2019, still works |
| Arbitrary websites / local HTML | Dash's built-in Docset Generator | Settings > Downloads |
| Javadoc / Doxygen / Rust / Scala / Go | javadocset, doxygen2docset, cargo-docset, sbt-dash, godocdash | per-toolchain |

Scrape-to-docset pipelines are normal practice for sites without clean exports: discover pages (e.g. from a `_toctree.yml`), rewrite internal links relative, inject dashAnchors from stable element ids, strip site chrome, build the SQLite index, write Info.plist, tar. 2026 user-contribution PRs are increasingly AI-co-authored (Claude session links in commit messages).

### Publishing [HIGH]

- **Dash User Contributions** (github.com/Kapeli/Dash-User-Contributions — 2.1k stars, 502 contributors, active May 2026): fork → copy `Sample_Docset` → add `<name>.tgz` + icons + README + `docset.json` (`name`, `version`, `archive`, `author`, `aliases`, optional `specific_versions`) → PR. Archives migrate to Kapeli's CDN. Custom-generator authors can request a free Dash license.
- **Self-hosted feed XML:** `<entry><version>1.2.0</version><url>https://example.com/MyDocs.tgz</url></entry>`; archive with `tar --exclude='.DS_Store' -cvzf`; share via `dash-feed://<percent-encoded feed URL>`.

### Integration surfaces [HIGH]

- URL schemes: `dash://?query=php:printf` (keyword-scoped search); `dash-plugin://` for multi-docset plugin calls [MEDIUM on exact grammar].
- Editor/launcher plugins: VS Code (`deerawan.vscode-dash`), JetBrains family, Alfred (reworked in Dash 8), Raycast, Vim/Neovim, Emacs (`helm-dash` works without Dash installed).
- CLI: **dasht** queries docSet.dsidx directly from the terminal — docsets without Dash.
- **AI/MCP (the 2025–2026 development) [HIGH]:** Dash 8 Settings > Integration enables MCP; official server `Kapeli/dash-mcp-server` (Python, fronts Dash's HTTP API): `claude mcp add dash-api -- uvx --from "git+https://github.com/Kapeli/dash-mcp-server.git" "dash-mcp-server"`. Tools: `list_installed_docsets`, `search_documentation` → take `load_url` → `load_documentation_page`, `enable_docset_fts`. Third-party: enhanced-dash-mcp (scans the full `~/Library/Application Support/Dash/` tree, project-context aware), DocsetMCP. Local example in this skill tree: `mongodb-docset-lookup` queries a version-pinned MongoDB Manual docset.

## Practical Patterns

Minimal hand-rolled pipeline:

```bash
mkdir -p MyDocs.docset/Contents/Resources/Documents/
cp -R built-html/* MyDocs.docset/Contents/Resources/Documents/
# Info.plist → MyDocs.docset/Contents/ (sample: kapeli.com/resources/Info.plist)
sqlite3 MyDocs.docset/Contents/Resources/docSet.dsidx \
  "CREATE TABLE searchIndex(id INTEGER PRIMARY KEY, name TEXT, type TEXT, path TEXT); \
   CREATE UNIQUE INDEX anchor ON searchIndex (name, type, path);"
# walk the HTML, INSERT OR IGNORE one row per symbol
tar --exclude='.DS_Store' -cvzf MyDocs.tgz MyDocs.docset
```

Tool selection: Sphinx-family → doc2dash; clean HTML → dashing or built-in generator; per-toolchain generators otherwise. Prefer User Contributions over self-hosted feeds for anything publicly useful.

## Anti-Patterns

- Re-adding a docset without removing it first after Info.plist/icon changes (changes don't take effect).
- `searchIndex.path` diverging from the online URL structure — breaks `DashDocSetFallbackURL` "Open Online Page".
- Assuming external JS runs (off unless `isJavaScriptEnabled=true`).
- Quoting docset counts interchangeably (Dash "200+" first-party ≠ Zeal "981" incl. user-contributed ≠ a user's installed tree).
- Treating MCP-directory blurbs as independent sources — they derive from the GitHub READMEs.

## Troubleshooting

| Symptom | Fix |
|---|---|
| Entries don't appear in search | Check entry `type` is in the supported list; `INSERT OR IGNORE` + unique index dedupe |
| TOC not showing | dashAnchors present? `DashDocSetFamily=dashtoc` set? Names percent-encoded? |
| "Open Online Page" 404s | Align paths with `DashDocSetFallbackURL` (extensionless sites → extensionless files) |
| MCP server can't see docsets | Requires Dash 8 + Settings > Integration enabled; official server fronts the local HTTP API |
| Full-text search missing | User-opt-in per docset unless `DashDocSetDefaultFTSEnabled=true` (or call `enable_docset_fts` via MCP) |

## References

Access 2026-06-10. [HIGH] unless noted.

1. kapeli.com/docsets — official format spec + generation guide (bundle anatomy, dsidx schema, entry types, dashtoc, feed XML). [docs]
2. kapeli.com/dash; blog.kapeli.com/dash-8 (2025-09-13, MCP + HTTP API); blog.kapeli.com/dash-7 (2023-08-17, subscription). [official]
3. github.com/Kapeli/Dash-User-Contributions — contribution flow, docset.json schema, activity (2026-05 commits); PR #5578 (Transformers scrape pipeline, AI-co-authored). [official repo]
4. github.com/Kapeli/dash-mcp-server — official MCP server (tools, install, Dash 8 requirement). [official repo]
5. github.com/hynek/doc2dash + doc2dash.hynek.me; github.com/technosophos/dashing (maintenance badge). [tool primaries]
6. zealdocs.org + github.com/zealdocs/zeal (v0.8.0 2026-02-28, v0.8.1 2026-04-04); issue #1594 (Kapeli's docset deprecations — decline signal); issue #247 (dash:// vs dash-plugin://). [primary]
7. freeCodeCamp/devdocs (maintainer search); velocity.silverlakesoftware.com (dormant) [MEDIUM].
8. github.com/sunaku/dasht; marketplace.visualstudio.com deerawan.vscode-dash; github.com/joshuadanpeterson/enhanced-dash-mcp (2025). [tool primaries]
9. dev.to "Local-First Documentation: What It Is and Why Your AI Agent Needs It" (2026-02-19) — local-docs-for-agents argument. [practitioner]
10. mjtsai.com/blog/2023/08/18/dash-7 — independent pricing analysis; HN threads 36137521 / 25396809 / 15508087 (sentiment, aggregate-only) [MEDIUM].
