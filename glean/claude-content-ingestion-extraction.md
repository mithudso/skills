# content-ingestion-extraction

**Category:** Cloud, DevOps & Infrastructure
**Platform:** Claude
**Original Path:** claude/standalone/content-ingestion-extraction

## Description
Acquire, extract, audit & restructure content from heterogeneous sources — aging docs, doc folders/KBs, live web DOM, meeting-audio devices, one-off docs to templatize. TRIGGER: doc archaeology (staleness, dead links, process drift, deprecate/archive decisions); doc-store bootstrapping (numbered taxonomy, _meta indexes, archive policy); Chrome content-script DOM extractor (selector fallback, MutationObserver/waitForElement, shadow-DOM, SPA nav, partial extraction); Granola transcript/notes fetch (REST, native-host, ProseMirror, polling sync); Plaud recorder ecosystem (devices, Developer API/OAuth, MCP, AutoFlow); reverse-engineer a doc into a reconstruction prompt with {{placeholders}} + data-manifest. SKIP: producing/parsing a finished file FORMAT → document-formats; general Chrome-extension architecture → chrome-extension-expert; prose/voice → writing-expert.

---

# Content Ingestion & Extraction

Hub for **getting content in and pulling content out** — messy front of pipeline, before clean file or report exists. Spans auditing aging docs, standing up doc stores, resilient web-DOM extraction, meeting-audio transcription ingestion, reverse-engineering finished doc into reusable template. Each former standalone skill now on-demand reference: when task matches row below, **`Read` that `references/<name>.md` file before answering**.

Boundary: hub owns **acquisition and extraction** of content from heterogeneous or live sources. Once content in hand and question becomes *which file format to emit* or *is prose any good*, defer to sibling hubs in cross-hub map.

<!-- ROUTING TABLE: content-ingestion-extraction -->
## Sub-skill routing table

Hub absorbs 6 former standalone skills as on-demand reference files. When task matches row, **Read listed `references/` file** before answering — don't rely on table alone for depth.

| Sub-topic | When to load | Reference file |
| --- | --- | --- |
| `doc-archaeology` | Audit aging doc for staleness, dead links, process drift, phantom deps, obsolete examples; salvage decision (update/restructure/deprecate/archive/delete) | `references/doc-archaeology.md` |
| `doc-store-bootstrapper` | Organize/upgrade doc folder (Drive engagement, KB) to mdb-tam standard — numbered taxonomy (00–99), `_meta/` indexes, archive policy, resumable memory | `references/doc-store-bootstrapper.md` |
| `dom-scraping-resilience` | Resilient Chrome content-script DOM extraction — selector fallback chains, MutationObserver/waitForElement, shadow-DOM traversal, SPA nav, partial extraction | `references/dom-scraping-resilience.md` |
| `granola-transcription` | Granola meeting transcripts — REST API, ProseMirror notes, polling sync, native-host bridge, corpus wiring, 401/403/429 handling | `references/granola-transcription.md` |
| `plaud-integration` | Plaud recorder ecosystem — devices, Developer Platform API (OAuth2/keys), MCP servers, AutoFlow/Zapier, privacy/compliance, extension integration | `references/plaud-integration.md` |
| `document-deconstructor` | Reverse-engineer doc into `{{placeholder}}` reconstruction prompt + `data-manifest.json` mapping each placeholder to source/method/last-value (templatize) | `references/document-deconstructor.md` |
| `using-plaud-mcp` | Plaud MCP quick-start — recordings, transcripts, memory_search/ingest, processing status, presigned audio URLs | `references/using-plaud-mcp.md` |

<!-- cross-hub-map -->
## Cross-hub map — where every content-ingestion-extraction topic lives

Family split across hubs. If task's deep material **not** in this hub's Sub-skill routing table, it's reference file under sibling hub below — **activate that hub or `Read` its `references/<name>.md` directly**. Every former standalone skill in family now reference under one of these hubs (nothing deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `content-ingestion-extraction` | Content Ingestion & Extraction (doc audit, scraping, transcription, templatizing) | `references/doc-archaeology.md`, `references/doc-store-bootstrapper.md`, `references/dom-scraping-resilience.md`, `references/granola-transcription.md`, … |