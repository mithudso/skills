<!-- Provenance: reference under the `splunk-platform-spl` standalone skill. Created 2026-06-18 via /dr deep-research (multi-source web research, ≥3 independent sources/concept). Volatile claims stamped verified-as-of: 2026-06-18. -->

# Splunk Knowledge Objects, Dashboards & Alerts

`verified-as-of: 2026-06-18` for Dashboard Studio direction and the Simple-XML deprecation status (volatile — verify). Sourced mostly to Splunk Docs (10.4 line), Splunk Community, Splunk Lantern.

## Contents
- Knowledge objects: fields & extraction, tags/event types, macros, lookups, calculated fields/aliases, permissions & knowledge bundle
- Dashboards: Simple XML, Dashboard Studio, base/post-process searches
- Alerts / scheduled searches
- Disconfirming findings / gotchas

## Knowledge objects
"Knowledge" is the collective term for objects associated with event data — event types, transactions, tags, saved searches, lookups, fields — that let you "interpret, classify, enrich, and normalize" data.[^1] Almost all of it is applied at **search time**, in a fixed sequence: field filters → inline `EXTRACT` → `REPORT` transforms → automatic KV → **field aliasing → calculated fields → lookups** → event types → tags.[^2] So an alias can feed a calculated field, which can feed a lookup, but not the reverse.

### Fields and field extraction
Splunk distinguishes index-time fields (extracted as data is written) from search-time "extracted fields" (computed when a search runs). The repeated official guidance: **prefer search-time extraction** — index-time extraction enlarges the index and slows search, recommended "only … when it is absolutely necessary" (the inverse of many users' intuition).[^3][^4][^5]

**The three `props.conf` extraction types:**[^3]
- `TRANSFORMS-<class>` → **index-time** extractions (references a transform in `transforms.conf`).
- `REPORT-<class>` → **search-time**, references a field transform in `transforms.conf`.
- `EXTRACT-<class>` → **search-time**, "inline," regex directly in `props.conf` (no `transforms.conf` needed).

**EXTRACT vs REPORT:** most of the time inline `EXTRACT` suffices; use `REPORT` (a field transform) when you reuse one regex across multiple extractions, apply multiple regexes to one config, or need delimiter-based handling of structured data.[^3] `transforms.conf` requires `REGEX` for index-time transforms; `FORMAT` is required at index time, optional at search time.[^5]

**The Field Extractor (IFX)** is a Splunk Web wizard that builds search-time extractions (Select sample → Select method → Select/Rename → Validate → Save) via **two methods**: **regular expression** (for *unstructured* data — generates a regex from highlighted sample values, refinable, with "required text" to scope it) or **delimiter** (for *structured* CSV-like data with a consistent separator; fails when field *positions* vary event-to-event — use regex then).[^8][^9] The ad-hoc SPL equivalent is `rex` (inline) or the `extract` command (works on `_raw`, can call `transforms.conf`).[^11]

### Tags and event types
- **Event type** = a label based on a **search expression**; saving a search (no time range / transforming commands) as an event type means any matching event gets `eventtype=<name>` (a multivalue field). Defined in `eventtypes.conf`; supports wildcards.[^12][^14]
- **Tag** = a name applied to a specific **field=value pair** (not necessarily unique); apply many tags to one field/value, defined collaboratively across `tags.conf` files; you can tag any field/value including eventtype/host/source/sourcetype.[^13][^14]
- "An event is a single instance of data … an event type is a classification used to label and group events."[^14] At scale (tens of thousands of items), docs recommend **lookups** over tags for categorization performance.[^13]

### Macros
Reusable chunks of SPL inserted into other searches — "can be any part of a search … and do not need to be a complete command."[^15] Invoked with **backticks** (`` `mymacro` ``); defined in `macros.conf` (or Settings → Advanced search). A definition is a string with `$arg$` substitutions; arguments declared with `args =` and the stanza carries an arity (`[foobar(2)]`); named args supported (`` `foobar(foo=1,bar=2)` ``).[^15][^16] Support validation expressions, `iseval=true` (eval-computed definitions), and nesting with cycle detection.[^16] Gotcha: if a macro expands to SPL beginning with a generating command (`from`/`search`/`tstats`/…), put a leading pipe before the macro call. **Why they reduce duplication:** one central, parameterized definition replaces repeated SPL across many searches/dashboards — change once, propagate everywhere.[^15][^16]

### Lookups
Lookups add fields from an external source to events based on existing field values.[^17] Official taxonomy = **four types by data source**:[^18]
1. **CSV (file-based / "static")** — values from a CSV; best for small, relatively static tables; stanza in `transforms.conf` (`filename=`).[^19]
2. **KV Store** — matches events to an App Key Value Store collection; for **large or frequently-updated** tables. Docs note: "investigate whether a CSV lookup will do … CSV lookups are easier to implement and suffice for the majority of cases."[^20]
3. **External / scripted** — a Python/binary script populates fields (the bundled example does DNS); input/output is CSV-formatted; **`outputlookup` is not supported** for external lookups.[^21]
4. **Geospatial** — maps coordinates to regions via a KMZ/KML file.[^18]

**Automatic lookups** are a *behavior overlay* (not a 5th type): an extra `props.conf` config runs a lookup in the background at search time on all matching events, scoped to a host/source/sourcetype. Splunk does **not** support nested automatic lookups.[^18][^22]

**Lookups vs joins for enrichment:** the `join` command is performance-limited (default **50,000-row / 60-second** subsearch caps).[^24] Mental model from SQL: **lookup ≈ LEFT OUTER JOIN, join ≈ INNER JOIN.** Repeated practitioner guidance: for enrichment prefer **lookups** (or a `stats`-by-common-field pattern) — "Forget join — that is not a Splunk way of doing things."[^23][^25] Invocation commands: `lookup` (add fields), `inputlookup` (read a table as source), `outputlookup` (write results to CSV/KV Store).[^18]

### Calculated fields & field aliases (brief)
- **Calculated fields** — `EVAL-<field> = <eval expr>` in `props.conf` (or Settings → Fields); moves commonly-used eval out of the search string, computed at search time. Caveat: all `EVAL-` configs in one stanza process **in parallel, not sequentially** — you cannot chain one calculated field into another within the same stanza.[^26][^28]
- **Field aliases** — `FIELDALIAS-<class> = <orig> AS <new>`; an alternate name (original not removed); a field can have many aliases, but one alias maps to one field.[^29]

### Permissions & knowledge bundle (brief)
- **Sharing levels:** **Private** (owner), **App** (current app), or **Global** (all apps). By default only **power** and **admin** roles can share/promote.[^30]
- **Knowledge bundle replication:** in distributed search, the SH packages its knowledge objects into a **knowledge bundle** and replicates it to its search peers (indexers), which are "ignorant of any local knowledge objects." First push to a new peer is the **full** bundle; subsequent pushes are usually the **delta**. Three policies: **Classic** (SH→each peer directly; default; not optimal beyond ~15–20 peers), **Cascading** (peers relay to peers for parallelism), **Mounted** (peers mount a shared bundle directory, eliminating replication).[^32][^33][^34]
- **Gotcha:** **KV Store collections are NOT bundle-replicated to indexers by default** — KV Store lookups run on the SH. To use a KV Store collection in an *automatic* lookup running on indexers, you must explicitly enable replication for that collection.[^22]

## Dashboards
### Simple XML (Classic)
Source is **Simple XML**; structure is `<dashboard>`/`<form>` → `<row>` → `<panel>` → visualization (`<chart>`/`<table>`/`<single>`/`<event>`/`<map>`/`<viz>`/`<html>`); layout is **row-column.**[^35][^36] Each panel can have its own inline search; a dashboard becomes a **form** when you add `<input>` elements (dropdown, radio, checkbox, text, time, multiselect).[^37] **Tokens** are dashboard variables — adding an input auto-generates a token; tokens carry selections into searches, drilldowns, set/unset, and conditional show/hide; predefined tokens exist for click events (`$click.value2$`).[^38][^39] **Drilldown:** clicking a visualization can run a secondary search, set tokens, open another dashboard, or link to a URL.[^40]

### Dashboard Studio
Splunk's **newer** framework; source is **JSON** (not XML), with a visual editor that wires visualizations to "data sources."[^35][^41] **Layouts:** **absolute** (free-form, pixel-perfect, shapes/lines/icons/custom backgrounds) and **grid** (panels snap to rows, auto-scales to the browser, surrounding panels resize when one is hidden — closest to classic).[^42][^36] **Direction (VOLATILE):** Dashboard Studio ships as a default app in Enterprise and Cloud, and active development lands there, not in XML — recent additions (10.x/2026) include a custom-visualizations framework, conditional panel show/hide, moving XML properties to JSON, and directly adding reports/saved searches. Splunk provides an automatic **Simple-XML→Studio converter** (clones, doesn't destroy) and recommends re-converting each release.[^44][^45][^46]

### Simple XML vs Dashboard Studio — and the deprecation question (disconfirming)
A common claim is "Simple XML is deprecated." **The precise, sourced reality (as of 10.4):** Splunk's "Deprecated and removed" pages do **NOT** list Simple XML dashboards themselves as deprecated. What *is* deprecated/removed: **HTML dashboards** (deprecated 8.2), **PDF export/printing for Classic Simple XML** (deprecated 9.3/9.4), specific old Simple XML tags, and **jQuery 2 support removed in 10.4.2604** (Simple XML dashboards with `version="1.0"` no longer load — change to `version="1.1"`). So the accurate statement is "Splunk steers new development to Studio and chips away at Classic's edges (HTML, PDF, jQuery), but Simple XML is not formally deprecated."[^48][^49]

**Documented Dashboard Studio limitations vs Simple XML:** **no custom JavaScript and no custom CSS** in Studio (Splunk's own comparison); **radio, checkbox, and link-list inputs not supported**; tokens "limited" relative to Classic. The converter marks custom CSS/JS, third-party visualizations, and event handlers **"Not converted / Unsupported"** (HTML panels only semi-convert → Markdown). Studio *advantages*: better out-of-box viz, working drag-and-drop, a dashboard **defaults** section, conditional formatting.[^50][^51][^52][^45]

### Base searches / post-process searches (dashboard efficiency)
Every concurrently-running search consumes a CPU; a 25–30-panel dashboard with one search per panel is very resource-intensive under many viewers.[^53] **Base + post-process** solves this: one "heavy" **base search** (`<search id="...">`) hits the indexes once; lightweight **post-process searches** (`<search base="...">`) filter/aggregate the cached results per panel.[^54][^53] In **Dashboard Studio** the equivalent is **base + chain searches** (up to 10 chains off one base).[^55] **Caveats (disconfirming the "always faster" assumption):** post-process is limited to non-transforming-before-the-pipe results; passing too many results/fields causes timeouts; stopping a post-process stops its base search **everywhere** it's used; and generic base searches returning large raw-event sets "perform very poorly" and eat disk quota — best for **aggregated** summary panels, not raw-event tables.[^54][^55][^56]

## Alerts / scheduled searches
**The relationship:** a **saved search** is the underlying object (one stanza per saved search in `savedsearches.conf`); a **report** is a saved search/pivot saved for reuse (schedulable, can trigger an action); an **alert** is a saved search with **alerting** turned on. So reports, alerts, and scheduled searches are all the **same saved-search object** with different settings layered on.[^57][^58][^59]

**Scheduled searches & cron:** scheduling uses `cron_schedule` (e.g., `*/5 * * * *`); the pre-4.0 `schedule =` is deprecated. Best practice: **stagger** schedules (offsets like `03,23,43 * * * *`) to avoid load spikes, and align the search time range to the interval (run every 20 min → `-20m`) to avoid gaps/overlaps.[^57][^60]

**Alert types (by search timing):**[^61]
- **Scheduled alert** — runs on cron/interval; triggers when results meet conditions. Use when real-time isn't required.
- **Real-time alert** — searches **continuously**; for immediate monitoring.

**Trigger conditions:**[^63][^61]
- **Per-result** (real-time only by default): fires **once for every matching result.**
- **Number-of-results / field-count**: fires once when results match a count condition; available for **scheduled** alerts and **real-time with a rolling time window** (e.g., >10 results in 5 minutes).
- **Custom condition**: any condition as a secondary search on the result set. (Trigger condition ≠ throttling.)

**Throttling:** after a trigger, suppresses further triggering for a set period to cut noise; specify a suppression **time period**, and (for per-result/real-time) optional **field values** so suppression is per-entity (e.g., once/hour per user).[^61][^62]

**Alert actions:** email (`action.email`, tokens like `$name$`), **webhook** (HTTP POST JSON; in 9.0+ the URL must be on the **webhook allow list**), output to a CSV lookup, log events, Triggered Alerts list, Send to Splunk Mobile, and **custom alert actions** (the modern extension point). **The simple "Run a script" alert action is DEPRECATED** — use custom alert actions.[^66][^67]

## Disconfirming findings / gotchas
1. **Index-time extraction is discouraged**, not encouraged — it bloats the index and slows search; default to search-time.[^4][^5]
2. **Simple XML is not formally deprecated** (only HTML dashboards / its PDF export / jQuery 2 are); correct the common overstatement.[^48][^49]
3. **Dashboard Studio cannot do custom CSS/JS** and lacks radio/checkbox/link-list inputs — a real migration blocker for complex Classic dashboards; conversion is lossy.[^50][^51][^52]
4. **Base/post-process searches aren't universally faster** — generic base searches returning large raw-event sets perform poorly and eat disk quota; best for aggregated summary panels only.[^56]
5. **Real-time alerts are costly** — "Real-time alerts can be costly … consider using a scheduled alert when possible"; multiple real-time searches negatively impact indexing capacity (real-time is optimized for sparse/rare-term searches, sacrificing indexing capacity for latency). Mitigations: **indexed real-time search**, restricting via the `rtsearch` capability, or converting to 1–5-min scheduled alerts ("nearly the same responsiveness"). Per-result real-time on an HA deployment can silently return incomplete results if a peer is down.[^68][^69][^70][^71]
6. **The "Run a script" alert action is deprecated** — use custom alert actions.[^66]
7. **KV Store collections are NOT bundle-replicated to indexers by default** — KV Store lookups run on the SH unless you enable collection replication.[^22]
8. **Calculated fields in one stanza run in parallel** — you can't chain `EVAL-` fields within the same stanza.[^28]

## Adjacent / frontier concepts
KV Store as a subsystem (collections.conf, REST CRUD, replication, accelerated fields); data models + `pivot` + acceleration (the higher-order knowledge object); the CIM (normalization at scale — see `spl-language-and-commands.md`); dashboard token + eval/conditional drilldowns; Splunk Web Framework / SplunkJS (the JS layer Classic exposes and Studio removes); workflow actions (linking events to external URLs/searches); summary indexing & report acceleration as enrichment/reporting alternatives (see `search-performance-optimization.md`); `distsearch.conf` replication tuning (classic/cascading/mounted).

## References
[^1]: Splexicon: Knowledge — Splunk Docs — https://docs.splunk.com/Splexicon:Knowledge — defines knowledge objects.
[^2]: The sequence of search-time operations — Splunk Docs (Knowledge Mgmt Manual) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.2/about-fields/the-sequence-of-search-time-operations — the search-time operation order.
[^3]: props.conf reference — Splunk Docs (Admin Manual 10.4) — https://help.splunk.com/en/splunk-enterprise/administer/admin-manual/10.4/configuration-file-reference/props.conf — TRANSFORMS/REPORT/EXTRACT types, EXTRACT-vs-REPORT.
[^4]: Configure custom fields at search time — Splunk Docs (Knowledge Mgmt 10.4) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/configure-fields-at-search-time — search-time vs index-time, prefer search-time.
[^5]: transforms.conf reference — Splunk Docs (Admin Manual 10.4) — https://help.splunk.com/en/splunk-enterprise/administer/admin-manual/10.4/configuration-file-reference/transforms.conf — index-time perf warning, REGEX/FORMAT.
[^8]: Build field extractions with the field extractor — Splunk Docs (Knowledge Mgmt 9.4) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/9.4/build-field-extractions-with-the-field-extractor — IFX wizard steps, two methods.
[^9]: Field Extractor: Select Method step — Splunk Docs (Knowledge Mgmt 10.4) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/build-field-extractions-with-the-field-extractor — regex (unstructured) vs delimiter (structured), delimiter failure mode.
[^11]: extract (SPL command) — Splunk Docs (Search Reference 10.4) — https://help.splunk.com/en/splunk-enterprise/search/spl-search-reference/10.4/search-commands/extract — `extract` works on `_raw`, can call transforms.conf.
[^12]: Basic difference between tags and event types — Splunk Community — https://community.splunk.com/t5/Splunk-Search/Basic-difference-between-tags-and-eventtypes/m-p/35381 — event type = search-expression label; tag = single key=value.
[^13]: About tags and aliases — Splunk Docs (Knowledge Mgmt 10.4) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/about-tags-and-aliases — tags = names for field/value combos; use lookups at scale.
[^14]: Classify and group similar events — Splunk Docs (Search Manual 10.4) — https://help.splunk.com/en/splunk-enterprise/search/search-manual/10.4/classify-and-group-similar-events/about-event-types — event vs event type; eventtype multivalue field.
[^15]: Use search macros in searches — Splunk Docs (Knowledge Mgmt 10.4) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/use-search-macros-in-searches — macros = reusable SPL, backticks, args, leading-pipe gotcha.
[^16]: macros.conf reference — Splunk Docs (Admin Manual 10.4) — https://help.splunk.com/en/splunk-enterprise/administer/admin-manual/10.4/configuration-file-reference/macros.conf — args/definition/$sub$, named args, iseval, nesting+cycle detection.
[^17]: Introduction to lookup configuration — Splunk Docs (Knowledge Mgmt 10.4) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/introduction-to-lookup-configuration — lookups add fields; four types.
[^18]: About lookups — Splunk Docs (Knowledge Mgmt 10.4) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/about-lookups — four lookup types, automatic lookups, lookup/inputlookup/outputlookup.
[^19]: Configure CSV lookups — Splunk Docs (Knowledge Mgmt 10.4) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/configure-csv-lookups — CSV/static lookup, transforms.conf filename.
[^20]: Configure KV Store lookups — Splunk Docs (Knowledge Mgmt 10.4) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/configure-kv-store-lookups — KV Store for large/frequently-updated; CSV-first guidance.
[^21]: Define an external lookup in Splunk Web — Splunk Docs (Knowledge Mgmt 10.4) — https://help.splunk.com/en/splunk-cloud-platform/manage-knowledge-objects/knowledge-management-manual/10.4/define-an-external-lookup — scripted/external lookups, Python/binary, CSV in/out, no outputlookup.
[^22]: Define an automatic lookup — Splunk Docs (Knowledge Mgmt 10.4) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/define-an-automatic-lookup — automatic-lookup scoping, no nesting, KV Store not bundle-replicated by default.
[^23]: Differences between join and lookup — Splunk Community — https://community.splunk.com/t5/Splunk-Search/Difference-between-join-and-lookup/m-p/659206 — lookup≈LEFT OUTER, join≈INNER; prefer lookup.
[^24]: join — Splunk Docs (Search Reference 10.4) — https://help.splunk.com/en/splunk-enterprise/spl-search-reference/10.4/search-commands/join — 50k-row/60s subsearch limits.
[^25]: SPL: Using lookup field values in tstats — Splunk Community — https://community.splunk.com/t5/Splunk-Search/SPL-Using-lookup-field-values-in-tstats/m-p/690921 — "forget join"; lookup vs subsearch.
[^26]: eval (calculated fields) — Splunk Docs (Search Reference 10.4) — https://help.splunk.com/en/splunk-cloud-platform/search/search-reference/10.4.2604/search-commands/eval — calculated fields move eval out of search.
[^28]: About calculated fields — Splunk Docs (Knowledge Mgmt) — https://help.splunk.com/en/splunk-cloud-platform/manage-knowledge-objects/knowledge-management-manual/10.3.2512/about-calculated-fields — EVAL- configs parallel not sequential; ordering after aliasing, before lookups.
[^29]: Configure field aliases with props.conf — Splunk Docs (Knowledge Mgmt 10.4) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/configure-field-aliases — FIELDALIAS syntax; alias doesn't remove original.
[^30]: Manage knowledge object permissions — Splunk Docs (Knowledge Mgmt 10.4) — https://help.splunk.com/en/splunk-enterprise/manage-knowledge-objects/knowledge-management-manual/10.4/manage-knowledge-object-permissions — private/app/global, power/admin roles.
[^32]: What search heads send to search peers — Splunk Docs (Distributed Search 10.4) — https://help.splunk.com/en/splunk-enterprise/administer/distributed-search/10.4/what-search-heads-send-to-search-peers — knowledge bundle definition & purpose.
[^33]: Knowledge bundle replication overview — Splunk Docs (Distributed Search 10.4) — https://help.splunk.com/en/splunk-enterprise/administer/distributed-search/10.4/knowledge-bundle-replication-overview — full vs delta, classic/cascading/mounted policies.
[^34]: Cascading/Mounted knowledge bundle replication — Splunk Docs (Distributed Search 10.4) — https://help.splunk.com/en/splunk-enterprise/administer/distributed-search/10.4/cascading-knowledge-bundle-replication — policy mechanics, ~15-20 peer classic limit.
[^35]: Simple XML reference / Getting started — Splunk Docs (Simple XML 10.4) — https://help.splunk.com/en/splunk-enterprise/create-dashboards/simple-xml/10.4 — XML structure; Studio uses JSON.
[^36]: Compare absolute and grid layouts — Splunk Docs (Dashboard Studio 10.4) — https://help.splunk.com/en/splunk-enterprise/create-dashboards/dashboard-studio/10.4/build-and-edit-the-layout/compare-absolute-and-grid-layouts — absolute vs grid; classic = row-column.
[^37]: Create and edit forms — Splunk Docs (Simple XML) — https://help.splunk.com/en/splunk-enterprise/create-dashboards/simple-xml/10.0/create-and-edit-forms — inputs convert dashboard→form; token per input.
[^38]: Token usage in dashboards — Splunk Docs (Simple XML 10.4) — https://help.splunk.com/en/splunk-cloud-platform/create-dashboards/simple-xml/10.4.2604/token-usage-in-dashboards — tokens as variables; predefined click tokens.
[^39]: Token reference — Splunk Docs (Simple XML 10.4) — https://help.splunk.com/en/splunk-enterprise/create-dashboards/simple-xml/10.4/token-reference — token types: form input, drilldown, conditional, set-destination.
[^40]: Event Handler Reference — Splunk Docs (Simple XML 10.4) — https://help.splunk.com/en/splunk-enterprise/create-dashboards/simple-xml/10.4/event-handler-reference — drilldown, condition, set, link event handlers.
[^41]: What is Splunk Dashboard Studio? — Splunk Docs (Dashboard Studio 10.2) — https://help.splunk.com/en/splunk-enterprise/create-dashboards/dashboard-studio/10.2/introduction-to-dashboard-studio/what-is-splunk-dashboard-studio — JSON source, visual editor, default app.
[^42]: Compare absolute and grid layouts (panel show/hide) — Splunk Docs (Dashboard Studio 10.4) — https://help.splunk.com/en/splunk-enterprise/create-dashboards/dashboard-studio/10.4/build-and-edit-the-layout/compare-absolute-and-grid-layouts — absolute pixel-perfect; grid auto-scale, panel resize on hide.
[^44]: Splunk Dashboard Studio vs Simple XML — reportcraft.app (blog, 2026) — https://reportcraft.app/blog/splunk-dashboard-studio-vs-simple-xml — "being deprecated" practitioner sentiment, Studio JSON complexity (qualify).
[^45]: What's new in Dashboard Studio — Splunk Docs (Dashboard Studio 10.4 Cloud) — https://help.splunk.com/en/splunk-cloud-platform/create-dashboards/dashboard-studio/10.4.2604/introduction-to-dashboard-studio/whats-new — custom-viz framework, XML→JSON property moves, conditional show/hide, add reports/saved searches.
[^46]: About conversion from Simple XML to Dashboard Studio — Splunk Docs (Dashboard Studio 10.4) — https://help.splunk.com/en/splunk-enterprise/create-dashboards/dashboard-studio/10.4/convert-dashboards/about-conversion-from-simple-xml-to-dashboard-studio — fully/semi/not-converted; re-convert each version.
[^48]: Deprecated and removed in version 10.4 — Splunk Docs (Release Notes 10.4) — https://help.splunk.com/en/splunk-enterprise/administer/release-notes/10.4/known-issues/deprecated-and-removed-features — HTML dashboards, Simple XML PDF export, jQuery 2 removal, version 1.0 won't load.
[^49]: Deprecated and removed in Splunk Cloud Platform — Splunk Docs (Cloud Release Notes 10.4.2604) — https://help.splunk.com/en/splunk-cloud-platform/administer/release-notes/10.4.2604/deprecated-and-removed-features — corroborates 10.4 deprecations.
[^50]: About conversion (custom CSS/JS unsupported) — Splunk Docs (Dashboard Studio 10.4) — https://help.splunk.com/en/splunk-enterprise/create-dashboards/dashboard-studio/10.4/convert-dashboards/about-conversion-from-simple-xml-to-dashboard-studio — custom CSS/JS/3rd-party viz/event-handlers Not converted/Unsupported; HTML→Markdown.
[^51]: HTML/CSS/JS limitations in Studio — Splunk Community — https://community.splunk.com/t5/Splunk-Dev/Dashboard-Studio-custom-CSS-JS/m-p/650000 — custom CSS/JS not supported in Studio.
[^52]: Dashboard Studio versus SimpleXML options — Splunk Community — https://community.splunk.com/t5/Dashboards-Visualizations/Dashboard-Studio-versus-SimpleXML/m-p/640000 — no radio/checkbox/link-list inputs, no JS/CSS, tokens limited.
[^53]: Post Process Searching — How to Optimize Dashboards — sp6.io (blog) — https://www.sp6.io/post-process-searching-how-to-optimize-splunk-dashboards — 1 CPU/search, base+post-process, reduction.
[^54]: Searches power dashboards and forms — Splunk Docs (Simple XML 10.4) — https://help.splunk.com/en/splunk-enterprise/create-dashboards/simple-xml/10.4/searches-power-dashboards-and-forms — base/post-process, `base` attribute, timeout/stop caveats.
[^55]: Chain searches with a base search — Splunk Docs (Dashboard Studio 10.4) — https://help.splunk.com/en/splunk-enterprise/create-dashboards/dashboard-studio/10.4/add-and-format-data-sources/chain-searches — Studio base+chain, 10-chain limit, "like post-process".
[^56]: Does using base searches increase performance? — Splunk Community — https://community.splunk.com/t5/Dashboards-Visualizations/Does-using-base-searches-increase-performance/m-p/640123 — base searches help only if aggregated; disk/quota cost for raw events.
[^57]: savedsearches.conf reference — Splunk Docs (Admin Manual 10.4) — https://help.splunk.com/en/splunk-enterprise/administer/admin-manual/10.4/configuration-file-reference/savedsearches.conf — cron_schedule, deprecated schedule, stagger guidance.
[^58]: Configure alerts in savedsearches.conf — Splunk Docs (Alerting) — https://help.splunk.com/en/splunk-cloud-platform/alerts/10.1/configure-alerts-in-savedsearches-conf — alerts = saved search; action.email/cron/counttype/quantity/relation.
[^59]: About reports — Splunk Docs (Reporting Manual) — https://help.splunk.com/en/splunk-enterprise/create-dashboards-and-reports/reporting-manual/9.0/about-reports — reports = saved search/pivot; scheduled reports trigger actions.
[^60]: Create scheduled alerts — Splunk Docs (Alerting 10.2) — https://help.splunk.com/en/splunk-enterprise/alerts/10.2/create-scheduled-alerts — earliest/latest align to interval, cron option.
[^61]: Alert types / comparison — Splunk Docs (Alerting) — https://help.splunk.com/en/splunk-cloud-platform/alerts/10.1/alert-types — scheduled vs real-time; per-result vs rolling-window; throttling options.
[^62]: Alert type and triggering scenarios — Splunk Docs (Alerting 10.2) — https://help.splunk.com/en/splunk-enterprise/alerts/10.2/alert-type-and-triggering-scenarios — scenario examples; rolling window; per-entity throttle.
[^63]: Configure alert trigger conditions — Splunk Docs (Alerting 10.4) — https://help.splunk.com/en/splunk-enterprise/alerts/10.4/configure-alert-trigger-conditions — trigger-condition table, once vs per-result, secondary-search semantics.
[^66]: Set up alert actions — Splunk Docs (Alerting 10.4) — https://help.splunk.com/en/splunk-enterprise/alerts/10.4/set-up-alert-actions — email/webhook/CSV-lookup/log-events/mobile; script action deprecated, custom actions.
[^67]: Use a webhook alert action — Splunk Docs (Alerting 10.0) — https://help.splunk.com/en/splunk-enterprise/alerts/10.0/use-a-webhook-alert-action — HTTP POST JSON, allow list (9.0+).
[^68]: Expected performance and known limitations of real-time searches — Splunk Docs (Search Manual 10.4) — https://help.splunk.com/en/splunk-cloud-platform/search/search-manual/10.4.2604/real-time-searches-and-reports/expected-performance-and-known-limitations-of-real-time-searches-and-reports — indexing-capacity impact, sparse-search optimization, indexed real-time.
[^69]: Create real-time alerts — Splunk Docs (Alerting 10.4) — https://help.splunk.com/en/splunk-enterprise/alerts/10.4/create-real-time-alerts — "costly … consider scheduled"; rolling window most demanding; HA caveat.
[^70]: How to restrict usage of real-time search — Splunk Docs (Search Manual 10.4) — https://help.splunk.com/en/splunk-enterprise/search/search-manual/10.4/real-time-searches-and-reports/how-to-restrict-usage-of-real-time-search — rtsearch capability, limits.conf caps, enableRealtimeSearch=false.
[^71]: Real-Time Alerts Performance Impact — thinkcloudly.com (blog, 2026) — https://thinkcloudly.com/blog/real-time-alerts-performance — convert RT→scheduled, specific searches (corroborates Tier-1 cost claims).
