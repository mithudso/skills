<!-- hub-reference-banner -->
> **Reference file — part of the `tam-operations` hub.** Formerly the standalone `customer-file-consolidator` skill.
> Sibling topics in this family are now reference files under the hubs (`tam-operations`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: customer-file-consolidator
version: "1.2.1"
updated: "2026-06-01"
description: >
  Collect, deduplicate, version-consolidate, and semantically analyze all local files for a named
  customer account, then produce a unified TAM briefing document and ingest it into the TAM corpus.

  TRIGGER: user says "consolidate files for <customer>", "gather all <customer> docs", "build a
  context file for <customer>", "run customer-file-consolidator", "merge my customer files",
  "create a briefing doc for <customer> from my downloads", "pull together everything on <account>",
  "find all files for <customer> and combine them", or invokes /customer-file-consolidator.

  SKIP: single-file analysis, general Atlas corpus queries, case analysis, TAM report generation,
  account health scoring, or any task that does not involve discovering and merging local files.
origin: local
category: custom
tags:
  - customer-files
  - consolidation
  - tam
  - briefing
  - corpus-ingestion
  - deduplication
related_skills:
  - tam-operations (references/tam-account-reports.md)
  - tam-operations (references/tam-reference.md)
  - tam-operations (references/account-artifacts-collector.md)
triggers:
  - consolidate files for
  - gather all customer docs
  - build a context file for
  - customer-file-consolidator
  - merge my customer files
  - create a briefing doc from downloads
  - pull together everything on
whenToUse:
  - "consolidate all files for a customer account into one briefing"
  - "gather and deduplicate customer documents from Downloads and Documents"
  - "build or refresh a customer context file from local files"
  - "ingest a consolidated customer briefing into the TAM corpus"
  - "merge customer files and archive older versions"
whenNotToUse:
  - "The customer name is ambiguous and matches more than three distinct companies — ask for clarification first"
  - "The target is a prospect, not an active customer account"
  - "The target is an internal MongoDB team or org — not a customer account"
  - "No local files exist for the account and the user is asking for account context from the corpus — use tam-operations (references/tam-account-reports.md) instead"
---

# Customer File Consolidator

Invoke with: `/customer-file-consolidator <customer_name>`

Example: `/customer-file-consolidator Goldman Sachs`

## Customer Resolution

On invocation:

1. **Parse the customer name** from the invocation arguments.
2. **Derive the customer key** by converting to snake_case (e.g., "Goldman Sachs" → "goldman_sachs").
3. **Look up aliases** in `~/.claude/monday-board-registry.json` if it exists — use the `display_name` and `account_id` fields. Common aliases can also be inferred (e.g., Goldman Sachs → GS, Disney → WDC).
4. **If ambiguous**, ask the user to confirm the customer key and any short aliases to search for.

## Output Discipline

Do not narrate actions, emit preamble, or add postscript. The only permitted outputs are: phase progress reports (one line per phase), the Phase 6 summary table, and the customer context file written to disk. No first-person commentary before, between, or after deliverables.

## Idempotency

The destination directory is `~/Downloads/<customer_key>/` (snake_case, e.g., `goldman_sachs`).

If it already exists with files from a prior run:
- Skip Phase 1 (file discovery across ~/Downloads and ~/Documents)
- Still run Phase 2 dedup on the destination contents (catches manually added files)
- Re-run Phase 3 version consolidation
- Re-run Phase 4 semantic analysis and synthesis (always regenerate the customer context file)
- Re-run Phase 5 corpus ingestion (upsert, not duplicate)

To force a full re-run from scratch, pass `--fresh` or delete the destination directory first.

---

## Phase 1 — Recursive File Discovery

### Alias Search Rules

| Alias length | Match scope |
|-------------|------------|
| 4+ chars (long) | All file types in filenames |
| 2–3 chars (short, e.g., "GS") | Document types only: `.pdf`, `.docx`, `.xlsx`, `.pptx`, `.html`, `.json`, `.md`, `.txt`, `.csv` |
| Directory match | Files inside directories named after the customer or aliases (e.g., `/GS/`, `/Goldman_Sachs/`) |

Short aliases use the document-type filter to avoid false positives from source code files containing the abbreviation.

### Excluded Directories

Always exclude:
- `~/Documents/dashboard/mdb-tam`
- `~/Documents/GitHub/mdb-case-assistant`
- `~/Documents/GitHub/mdb-context-hub`
- `node_modules/`, `.git/`

### Search Steps

1. Search `~/Downloads` and `~/Documents` recursively for files matching the customer name (case-insensitive) and each alias (with the short-alias document-type filter).
2. Report: total files found, breakdown by file type, breakdown by source directory (top 20).
3. **File-count guard:** if more than 100 files are found, pause and show the breakdown before proceeding. Ask: "Found N files matching `<customer>`. Proceed with all, or filter further?" Do not copy files until confirmed. If a single source directory accounts for >80% of matches (e.g., a `~/Documents/GitHub` repo), call it out as a likely false-positive source.

---

## Phase 2 — Consolidation & Deduplication

1. Create directory: `~/Downloads/<customer_key>/`
2. Copy all matched files. On name collision:
   - If MD5 hash matches existing file: skip (exact duplicate)
   - If hash differs: add numeric suffix (`_2`, `_3`, etc.)
3. **Confirmation required before deletion:** "About to remove N duplicate files. Show list and confirm?" Do not delete without explicit user confirmation. List the files to be removed (filename + reason: exact hash match).
4. Deduplicate by MD5 content hash: for each group of identical-hash files, keep the shortest/cleanest filename, remove the rest (after confirmation).
5. Report: files before dedup, duplicates removed, files after dedup.

---

## Phase 3 — Version Consolidation

Archive older versions into `~/Downloads/<customer_key>/archived/`:

- Files with `_2`, `_3`, `_4` numeric suffixes where the base file exists
- Files with `(1)`, `(2)`, `(3)` parenthetical copies where the base file exists
- Versioned documents (V2, V3, V3.1 when V4 exists) — keep the latest
- Dated variants of the same document — keep the most recent date, archive older
- Context module exports — keep the latest date set, archive older date sets

Report: files archived, files remaining, total in archive.

---

## Phase 4 — Semantic Analysis & Rollup

### File Reading

For all remaining files (excluding `archived/`, `.zip`, `.png`, `.jpg`, `.bak`, `.diff`, `.yaml`, `.DS_Store`, `.pptx`):

- `.yaml` excluded: typically contains IaC/config scaffolding, not narrative customer content.
- `.pptx` excluded: reliable text extraction requires LibreOffice or python-pptx. If extraction is needed, ask the user to export to PDF first.
- Password-protected files that cannot be opened: report as "excluded — password protected" in Sources and skip.

| Format | Extraction method |
|--------|------------------|
| `.md` | Read directly |
| `.docx` | `python3 -c "import mammoth, sys; r=mammoth.convert_to_markdown(open(sys.argv[1],'rb')); print(r.value)" <path>`; fallback: `unzip -p <path> word/document.xml \| python3 -c "import sys,re; print(re.sub('<[^>]+>','',sys.stdin.read()))"` |
| `.xlsx` | `python3 -c "import openpyxl,sys; wb=openpyxl.load_workbook(sys.argv[1]); [print(row) for ws in wb for row in ws.iter_rows(values_only=True)]" <path>`; if `openpyxl` missing: `python3 -m pip install openpyxl -q && <retry>` |
| `.pdf` | `pdftotext <path> -`; fallback: `python3 -m pdfminer.high_level <path>`; if both fail: report as "excluded — extraction failed" |
| `.html` | Strip tags, extract text body |
| `.json` | Parse and extract key metrics/summaries |
| `.csv` | Read headers and first 20 rows as sample |

### Signal Quality Rating

| Rating | Criteria | Examples |
|--------|----------|---------|
| **HIGH** | Primary authoritative source covering ≥3 account sections, dated within 12 months (or undated but clearly a canonical account doc) | Account reviews, support plans, playbooks, engagement overviews, weekly updates, customer context exports |
| **MEDIUM** | Supporting detail or partial coverage of 1–2 sections, or any age | Case analyses, meeting transcripts, contacts, structured JSON data, antipattern reports |
| **LOW** | Organizational/metadata only, navigation stubs, or confirmed duplicate content | Processed sidecars, duplicated transcripts, format variants, Google Drive shortcut metadata |

When a file has no date in its name, use the file system `mtime` as the proxy date. If `mtime` is unreliable (e.g., freshly synced from cloud), treat the file as undated.

### Customer Context File

The output artifact is the **customer context file** (the "rollup" produced by the synthesis passes below). Write to `~/Documents/dashboard/mdb-tam/customer-files/<slug>/<account>_customer_context_<YYYY-MM-DD>.md`
where `<slug>` is the kebab-case account slug (e.g., `goldman-sachs`) and `<account>` is the snake_case key (e.g., `goldman_sachs`).

This file is the primary handoff artifact consumed by tam-operations (references/tam-account-reports.md). That reference enforces a staleness check against `<YYYY-MM-DD>` and refuses files older than 72 hours, so the dated filename is required. Keep the path and structure stable.

**Quality bar:** a correct customer context file (a) cites at least one HIGH or MEDIUM signal source per populated section, (b) marks every unpopulated section with `[No data found]` rather than omitting it, and (c) contains no information not traceable to a source file — do not invent or infer facts beyond what the source files state.

```markdown
# <Customer Name> — Customer Context File

**AS-OF:** <today's date in YYYY-MM-DD> | **Author:** TAM AI Assistant (consolidated from <N> source files) | **Classification:** MongoDB Internal

## Account Overview
- Company profile (legal name, industry, employees, revenue, SFDC ID, segmentation)
- Commercial summary (ARR, subscription period, renewal window, blockers)
- MongoDB account team (table: name, role, responsibilities)
- Customer key contacts (table: name, function, notes)
- TAM engagement model (cadences, forums, constraints)

## Technical Landscape
- Atlas infrastructure summary (clusters, orgs, projects, versions, topology, backup)
- Atlas organizations (table: org, cluster count, workloads)
- Priority production clusters (table: project, cluster, tier, topology, notes)
- Product/component usage (table: product, status, notes)
- Competitive landscape

## Support Posture
- Open cases (table: case, severity, summary, status, owner, opened)
- 12-month case volume (table: severity, closed, open)
- HELP ticket chains (if applicable)
- Past incident history (key RCAs)

## Active Initiatives
- P1/Active current quarter (table: initiative, status, target, key details)
- Backlog/Watch (table: initiative, status, window)
- Recently completed (table: initiative, period, outcome)

## Meeting & Engagement History
- Key meeting summaries (most recent first, date + bullet format)
- Engagement patterns

## Risk Factors & Opportunities
- Active risks (priority-ordered, severity + mitigation)
- Expansion opportunities (numbered)
- Renewal indicators (positive, caution, action items)

## Customer Plans
- Short-term (0–6 months)
- Mid-term (6–18 months)
- Long-term (18+ months)

## Key Documents
- Table: document, location, purpose

## Sources
- HIGH Signal (numbered, filename + description)
- MEDIUM Signal (numbered, filename + description)
- LOW Signal (range-numbered, grouped by type)
- Excluded (with reason)
- Deduplication notes
```

### Convergence

Run iterative passes on the rollup document:
1. **Compile** all extracted data into the template
2. **Deduplicate** information across sources (same fact from multiple files = cite all sources but state once)
3. **Resolve contradictions** by favoring the most recent source
4. **Remove filler** (generic descriptions, boilerplate, AI-generated scaffolding)
5. **Verify** all sections are populated; mark empty sections with `[No data found]`

Stop when pass N produces a rollup byte-for-byte identical to pass N-1, or after 5 passes, whichever comes first.

**Post-synthesis self-check (run once after the final pass):**
1. Confirm all 8 sections (`Account Overview` through `Sources`) are present.
2. For every populated cell, verify it traces to a named source file — if not, remove it and mark the section `[No data found]`.
3. If any section is missing, add it with `[No data found]`.

---

## Phase 5 — Corpus Ingestion

Use the TAM MCP server (`mcp__mdb_tam_account_context__*`) — **NOT** the MongoDB MCP (`mcp__plugin_mongodb__*`), which targets the developer's local database, not the TAM corpus.

1. **Resolve the canonical `account_id`:**
   - Use the `account_id` (numeric/SFDC ID) from `~/.claude/monday-board-registry.json` if available.
   - If not available, use the snake_case `customer_key` as a fallback and note the limitation in the summary report.

2. **Upsert into `ingested_sources`** using the TAM MCP corpus tools. Document shape:
   ```json
   {
     "id": "customer_context_file_<customer_key>_<date>",
     "account_id": "<canonical_account_id>",
     "account_name": "<customer_name>",
     "source_type": "customer_context_file",
     "title": "<customer_name> — Customer Context File (Consolidated)",
     "ingest_ts": "<iso_timestamp>",
     "source_files_analyzed": <count>,
     "total_files_discovered": <count>,
     "unique_sources_cited": <count>,
     "lines": <count>
   }
   ```

3. **Before ingestion, redact:**
   - Connection strings containing passwords (e.g., `mongodb+srv://user:PASSWORD@...` → `mongodb+srv://[redacted]@...`)
   - API keys, tokens, or secrets (any value matching `[A-Za-z0-9_-]{32,}` adjacent to a key/token label)
   - Personal contact details beyond name + job title (personal phone numbers, personal email addresses)
   - HR or compensation data
   Then include the redacted customer context file markdown as the corpus document body.

4. **Failure handling:** If corpus ingestion fails (server unavailable, auth expired, network error):
   - The customer context file is still the primary deliverable — do not abort.
   - In the Phase 6 summary, mark corpus ingestion as FAILED with the error message.
   - Recovery: re-invoke the skill — idempotency rules will skip Phase 1 and re-run Phases 2–5, retrying corpus ingestion.

5. **Verify:** use `mcp__mdb_tam_account_context__mdb_tam_corpus_query` with `collection: "ingested_sources"` and `account_name: "<customer_name>"` to confirm the document is retrievable.

---

## Phase 6 — Summary Report

Output a summary table:

| Phase | Input | Output |
|-------|-------|--------|
| Discovery | Directories scanned | Files found (by type) |
| Dedup | Files before | Files after (duplicates removed) |
| Version Consolidation | Files before | Files remaining + archived |
| Semantic Analysis | Files analyzed | Sources by signal quality |
| Rollup | Sources cited | Lines written, convergence passes |
| Corpus Ingestion | Document ID | Verification status |
