<!-- hub-reference-banner -->
> **Reference file — part of the `tam-operations` hub.** Formerly the standalone `account-artifacts-collector` skill.
> Sibling topics in this family are now reference files under the hubs (`tam-operations`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: account-artifacts-collector
version: 1.1.2
updated: "2026-06-01"
description: >
  Orchestrates parallel data collection across all MCP and local sources for a customer account,
  persisting typed JSON artifacts to disk for downstream report generation, account reviews, and
  agent consumption.

  TRIGGER: any account review, QBR/EBR prep, report generation (case_analysis_report,
  daily_catchup_report, tam_action_report), account health assessment, shift handoff, meeting
  prep, or explicit /collect-account-artifacts invocation. Also triggered when the
  account-data-collector agent dispatches collection.

  SKIP: single-case lookups that do not need full account context, or when artifacts already
  exist and are fresh (check manifest.json collected_at timestamps against staleness thresholds
  below before deciding to re-collect).
origin: local
category: custom
tags:
  - account-data
  - artifact-collection
  - parallel-agents
  - report-prep
  - QBR
  - EBR
  - account-review
  - data-pipeline
triggers:
  - collect account artifacts
  - account artifacts
  - account review
  - QBR prep
  - EBR prep
  - report generation
  - account data collection
  - /collect-account-artifacts
related_skills:
  - case-tracker
  - case-mcp-server-guide
  - tam-reference
  - tam-expertise
  - operator-report-generator
  - account-health-scorer
  - ts-tools-support-api
whenToUse:
  - "collect all account data before generating a report or account review"
  - "QBR or EBR preparation requiring full account context"
  - "shift handoff requiring a fresh snapshot of cases, meetings, and initiatives"
  - "meeting prep requiring cases, Slack, stakeholders, and commercial data"
  - "running /collect-account-artifacts"
whenNotToUse:
  - "single-case lookup — use mdb_case_get_case directly"
  - "checking server health or MCP status — use mdb_case_get_server_status"
  - "artifacts already exist and are within staleness thresholds (check manifest.json first)"
  - "running a diagnostic on a single Atlas cluster — use atlas-diagnostics-expert"
---

# Account Artifacts Collector

Orchestrates parallel data collection from every available MCP and local source for a single
customer account. Persists structured JSON artifacts to disk so downstream agents and report
templates can consume them without re-fetching.

## Invocation

```
/collect-account-artifacts <account_name>
```

Or triggered automatically by the `account-data-collector` agent, any report template requiring full account context, or account review / QBR/EBR / meeting prep workflows.

---

## Artifact Storage Layout

```
~/.customer-dashboard/artifacts/<account_slug>/<YYYY-MM-DD>/
```

`<account_slug>`: lowercase, whitespace runs → single hyphen, strip `[^a-z0-9-]`, collapse consecutive hyphens, trim leading/trailing hyphens. Example: `"Goldman Sachs & Co."` → `"goldman-sachs-co"`.

| File | Contents |
|------|----------|
| `cases.json` | Open/in-progress/recently-closed cases; severity, status, dates, next actions, comments |
| `sfdc.json` | SFDC account record, contacts, opportunities, renewal data, competitors |
| `atlas-infrastructure.json` | Cluster configs, backup jobs, monitoring metrics, DB/collection counts, limits |
| `initiatives.json` | Active/backlog/completed initiatives from Monday boards and corpus |
| `stakeholders.json` | Customer + MongoDB team contacts and roles from SFDC + corpus |
| `documents.json` | TAM Engagement Overview, Support Plan, JIMP summaries, weekly summaries |
| `slack-activity.json` | Recent Slack messages with timestamps and channel context |
| `meetings.json` | Meeting notes, action items, attendees, dates |
| `help-tickets.json` | JIRA HELP ticket chain: status, assignee, resolution, last comment |
| `commercial.json` | ARR, renewal dates, competitors, spend tiers, opportunity pipeline |
| `manifest.json` | Metadata: account, date, source versions, per-source collection timestamps, errors |

## Artifact Envelope Format

Every `.json` file:

```json
{
  "collected_at": "2026-05-26T14:30:00.000Z",
  "source": "mdb_case_assistant",
  "account": "Account Name",
  "account_id": "acct_abc123",
  "account_slug": "account-name",
  "data": { }
}
```

All five envelope fields (`collected_at`, `source`, `account`, `account_id`, `account_slug`) are required in every artifact file.

## Staleness Thresholds

| Source | Stale after |
|--------|------------|
| Cases, HELP tickets | 15 minutes |
| Slack activity | 30 minutes |
| Atlas infrastructure | 1 hour |
| SFDC, commercial, stakeholders | 4 hours |
| Meetings, documents, initiatives | 6 hours |

## Incremental Update Strategy

Before collecting, check for existing artifacts at the same `<account_slug>/<date>/` path:

1. Read `manifest.json` if it exists.
2. For each source, compare `collected_at` timestamp against the staleness threshold above.
3. Re-collect only sources that are stale or missing.
4. Merge new data: replace the entire `data` field and update `collected_at`.
5. Update `manifest.json` with new per-source timestamps.

If `--force` flag is passed (or `force: true` from an agent), skip staleness checks and re-collect everything.

---

## Phase 0 — Preparation

1. **Determine today's date** in `YYYY-MM-DD` format.
2. **Resolve the account identity:**
   - Call `mdb_case_get_account` with the account name.
   - If that fails, search with `mdb_case_search` using the account name as query.
   - If both fail, abort: `"Account not found: <account_name>. Verify the account name matches SFDC or the Support portal."` Do not proceed to Phase 1.
   - Store `account_id`, `account_name`, and `account_slug` for all subsequent calls.
3. **Create the output directory:**
   ```bash
   mkdir -p ~/.customer-dashboard/artifacts/<account_slug>/<YYYY-MM-DD>/
   ```
4. **Load existing manifest** (if present) to determine which sources need refresh.

---

## Phase 1 — Parallel Data Collection

Dispatch all collectors **in parallel**. Each collector is independent and writes its own artifact file. If a source errors, log the error in the manifest and continue — never let one failure block the others. Individual collector timeout: 60 seconds. Overall orchestration timeout: 120 seconds.

### Collector 1: Cases (Case Assistant MCP)

**Source:** `mdb_case_assistant` MCP | **Output:** `cases.json`

1. Call `mdb_case_list_account_cases` with the resolved account ID.
2. Partition: open (status = Open, In Progress, Waiting on MongoDB) vs recently closed (Closed AND last update within 7 days).
3. For each open case (up to 30), call in parallel: `mdb_case_get_case`, `mdb_case_get_case_next_action`, `mdb_case_get_case_stage`, `mdb_case_get_case_comments` (limit 5).
4. For each recently closed case, call `mdb_case_get_case` only.
5. Compute derived fields: `days_open` (today − creation date), `days_since_update` (today − last comment or status change), `severity_weight` (S1=4, S2=3, S3=2, S4=1).

### Collector 2: SFDC Account Data (Glean MCP)

**Source:** `glean_default` MCP | **Output:** `sfdc.json`

Search Glean for: SFDC account record, contacts, opportunities/renewals, TAM engagement overview, support plan. For promising URLs, call `mcp__glean_default__read_document` for full content.

### Collector 3: Atlas Infrastructure (Server Feeds)

**Source:** `127.0.0.1:8787` | **Output:** `atlas-infrastructure.json`

```
GET /api/feeds/atlas-cluster-configs?account=<account_slug>
GET /api/feeds/atlas-backup-jobs?account=<account_slug>
GET /api/feeds/atlas-monitoring?account=<account_slug>
GET /api/feeds/atlas-db-collection-counts?account=<account_slug>
GET /api/feeds/atlas-limits?account=<account_slug>
```
Auth: `Authorization: Bearer ${DASHBOARD_API_TOKEN}`. If server unreachable, log error and skip.

### Collector 4: Corpus Context (TAM MCP)

**Source:** `mdb_tam` MCP | **Output:** `documents.json` and `meetings.json`

Search `mcp__tam_mcp__tam_search_all` for: `<account_name>`, `<account_name> meeting`, `<account_name> weekly summary`, `<account_name> engagement overview`, `<account_name> support plan`, `<account_name> JIMP`. For high-relevance results, call `mcp__tam_mcp__tam_get_file_analysis` for full content. Split results: meeting notes → `meetings.json`; TAM docs → `documents.json`.

### Collector 5: Stakeholders (Case Assistant MCP)

**Source:** `mdb_case_assistant` MCP | **Output:** `stakeholders.json`

Call `mdb_case_get_account` to get team information and customer contacts. SFDC contact enrichment (from Collector 2) is deferred to Phase 2 assembly to avoid a parallel write race.

### Collector 6: HELP Tickets (Server Feed + Case Assistant)

**Source:** `127.0.0.1:8787` + `mdb_case_assistant` | **Output:** `help-tickets.json`

```
GET /api/feeds/jira-help-tickets?account=<account_slug>
```
Also search via `mdb_case_search` with `"HELP <account_name>"`. For each HELP ticket key found, call `mdb_case_get_help_ticket`.

### Collector 7: SFDC Commercial Data (Server Feed + Glean)

**Source:** `127.0.0.1:8787` + `glean_default` | **Output:** `commercial.json`

```
GET /api/feeds/sfdc-team-and-opps?account=<account_slug>
```
Also search Glean for: `<account_name> ARR renewal contract`, `<account_name> competitor spend tier`.

### Collector 8: Slack Activity (Local MongoDB)

**Source:** Local mongod via `mcp__plugin_mongodb_mongodb__find` | **Output:** `slack-activity.json`

```
database: "dashboard_corpus", collection: "slack_messages"
filter: { "account": "<account_name>" }, sort: { "timestamp": -1 }, limit: 50
```
Fallback if local mongod unavailable: `GET http://127.0.0.1:8787/api/corpus?collection=slack_messages&account=<account_slug>&limit=50&sort=-timestamp`

### Collector 9: Initiatives (Local MongoDB + Monday)

**Source:** Local mongod `initiatives` collection + Monday boards | **Output:** `initiatives.json`

Query local mongod for initiatives filtered by account name. If Monday MCP is available, also query Monday boards for non-Done items. Merge and deduplicate from both sources.

### Collector 10: Case History Metrics (Server Feeds)

**Source:** `127.0.0.1:8787` | **Note:** Enriches `cases.json` but stored in memory until Phase 2

```
GET /api/feeds/case-history-mttr?account=<account_slug>
GET /api/feeds/gs-open-case-recheck?account=<account_slug>
GET /api/feeds/lapsed-action-checker?account=<account_slug>
```
Store raw results in memory. Do **not** write to `cases.json` yet — merge happens in Phase 2.

---

## Phase 2 — Assembly and Manifest

After all collectors complete (or after 120-second timeout):

1. **Post-collection merges:**
   - **Stakeholders + SFDC cross-reference:** If both Collectors 5 and 2 succeeded, merge SFDC contact data (title, email, role) into `stakeholders.json`. Deduplicate by name.
   - **Cases + metrics enrichment:** If both Collectors 1 and 10 succeeded, merge MTTR stats, lapsed-action flags, and recheck results into `cases.json`.

2. **Write each artifact file** with `JSON.stringify(data, null, 2)`.

3. **Build `manifest.json`:**
   ```json
   {
     "collected_at": "<now ISO>",
     "source": "account-artifacts-collector/1.1.2",
     "account": "<account_name>",
     "account_id": "<account_id>",
     "account_slug": "<account_slug>",
     "collection_date": "<YYYY-MM-DD>",
     "sources": {
       "cases": { "status": "ok|error|partial|skipped", "collected_at": "...", "record_count": 12, "error": null }
     },
     "report_templates_compatible": ["case_analysis_report", "daily_catchup_report", "tam_action_report"],
     "incremental": false,
     "duration_ms": 8423
   }
   ```

4. **Print a summary** to the user showing each artifact file, record count, and status.

---

## Phase 3 — Report Template Compatibility

| Template | Primary input | Enrichment |
|----------|--------------|-----------|
| `case_analysis_report.md` | `cases.json` | `meetings.json`, `help-tickets.json` |
| `daily_catchup_report.md` | `cases.json`, `meetings.json` | `documents.json`, `initiatives.json` |
| `tam_action_report.md` | `cases.json`, `meetings.json`, `initiatives.json` | `documents.json`, `commercial.json` |

All three templates use `{{TODAY}}` = `collection_date`, `{{ACCOUNT_NAME}}` = account. Report templates live at: `~/Documents/GitHub/mdb-context-hub/reports-library/`.

---

## Error Handling

| Error type | Behavior |
|-----------|---------|
| MCP server unavailable | Log `"status": "error"` + message in manifest; continue with remaining collectors |
| Partial data | Write what was collected; set `"status": "partial"` + `"warning"` in manifest |
| Auth failure | Log distinctly from connectivity issues; check `DASHBOARD_API_TOKEN` for server feeds; check `mdb_case_get_support_auth_status` for case assistant |
| Timeout | Individual collectors time out after 60s; overall orchestration at 120s — write whatever has been collected |
| Empty results | Write artifact with empty `data` array/object and `"status": "ok"` with `"record_count": 0` — not an error |

---

## Cross-References

- **Report templates:** `~/Documents/GitHub/mdb-context-hub/reports-library/`
- **Server feeds:** `server/src/jobs/feeds/` (atlas-cluster-configs, atlas-backup-jobs, atlas-monitoring, atlas-db-collection-counts, atlas-limits, case-history-mttr, gs-open-case-recheck, jira-help-tickets, lapsed-action-checker, sfdc-team-and-opps)
- **Corpus store:** `src/background/corpus-store/dual-write-corpus-store.js` (IndexedDB primary, local mongod mirror via server `/api/corpus`)
