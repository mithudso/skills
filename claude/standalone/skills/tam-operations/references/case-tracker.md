<!-- hub-reference-banner -->
> **Reference file — part of the `tam-operations` hub.** Formerly the standalone `case-tracker` skill.
> Sibling topics in this family are now reference files under the hubs (`tam-operations`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: case-tracker
version: 1.1.2
updated: "2026-06-01"
description: >
  Authoritative reference for building, extending, and debugging the Customer Dashboard active
  case tracker. Covers the TS Tools Support API, case schema, severity model, customer
  subscription flows, LLM summarization patterns, diagnostic tool catalog, and cross-references
  to all MongoDB support skills.

  TRIGGER: user needs to build new case tracker features, extend the diagnostic tools registry,
  write API calls against the TS Tools Support API, triage a MongoDB support case using the
  diagnostic tool catalog, write or debug LLM prompts for case summarization, understand severity
  model or polling cadence, connect the tracker to new data sources, or extend the service-worker
  message API.

  SKIP: Atlas operational triage workflows without extension context (use atlas-diagnostics-expert),
  general Chrome extension architecture decisions (use chrome-extension-expert), or account-level
  artifact collection (use tam-operations (references/account-artifacts-collector.md)).
origin: local
category: custom
tags:
  - case-tracker
  - chrome-extension
  - support-api
  - llm-analysis
  - severity-model
  - service-worker
  - diagnostic-tools
triggers:
  - case tracker
  - case tracking
  - support API
  - ts-tools API
  - case analyzer
  - tracker polling
  - case severity
  - LLM case analysis
  - case summarization
  - diagnostic tools registry
related_skills:
  - atlas-diagnostics-expert
  - mongodb-kb
  - mongodb-expert
  - chrome-extension-expert
  - tam-operations (references/tstools-reference.md)
whenToUse:
  - "build new case tracker features or extend the diagnostic tools registry"
  - "write API calls against the TS Tools Support API"
  - "triage a MongoDB support case using the diagnostic tool catalog"
  - "write or debug LLM prompts for case summarization and analysis"
  - "understand the severity model, polling cadence, or customer subscription flows"
  - "connect the tracker to new data sources or extend the service-worker message API"
whenNotToUse:
  - "Atlas operational triage without extension context — use atlas-diagnostics-expert"
  - "general Chrome extension architecture — use chrome-extension-expert"
  - "account-level artifact collection — use tam-operations (references/account-artifacts-collector.md)"
---

# Customer Dashboard — Active Case Tracker Reference

Source files: `src/background/case-tracker.js` (3424 lines), `src/background/ts-tools-api.js` (389 lines), `src/background/service-worker.js`, `src/background/llm.js`

---

## 1. Case Tracker Lifecycle

### Overview

A *case tracker* is a persisted record in `chrome.storage.local` under `case_trackers` (object keyed by tracker ID) that monitors one external source — a CUB case, a Jira HELP ticket, or a generic URL — and periodically fetches its state to produce an LLM-generated analysis.

### Tracker ID format

`{scope}::{kind}::{ref_or_url}` where:
- `scope` = `account_id` or normalized `customer_name` (lowercased, non-alphanumeric → spaces) or `_global`
- `kind` = `cub` | `jira` | `generic`
- `ref` = uppercased case/ticket reference or lowercased URL

Example: `acct_abc123::cub::01234567`, `acct_abc123::jira::HELP-9999`

### Creation

```
parseTrackerSourceInput(sourceInput)
  → saveCaseTracker(record)        // writes to chrome.storage.local
  → ensureCaseTrackerAlarm()       // creates chrome.alarms entry
  → refreshCaseTracker(id, { force: true, reason: 'initial' })
```

`parseTrackerSourceInput` accepts: 8-digit case number, `HELP-NNNNN`, full Hub URL, full Jira URL. Throws `Error` on generic non-whitelisted URLs.

### Poll / Alarm Cycle

- Alarm name: `CASE_TRACKER_ALARM_NAME = 'case-tracker-poll'`
- Period: `CASE_TRACKER_DEFAULT_INTERVAL_MINUTES = 1` minute
- On each alarm: `pollDueCaseTrackers(Date.now())` → `syncAutomaticCaseTrackers()` → load enabled trackers → compute effective interval → skip if `(now - last_checked_at) < interval` → `refreshCaseTracker`
- Global `caseTrackerPollLock` prevents concurrent poll runs

### Effective interval per tracker

`max(1, floor(max(derivedInterval, storedInterval) || 1))` where `derivedInterval = getAutoTrackerPollIntervalMinutes(tracker.severity)`

### Refresh Flow

1. Load tracker. If not found, return early.
2. If archived (and not forced), skip.
3. If disabled (and not forced), skip.
4. **List-driven gate** (non-S1 CUB, non-force, non-initial): check if list-row signal (title, status, severity, last_updated) changed since last check. If unchanged, update `last_list_checked_at` and skip full refresh.
5. **CUB trackers**: try `extractHubApiTrackerSnapshotFromTab` on an existing open Hub tab.
6. **Fallback / all kinds**: `ensureBackgroundTrackerTab` → navigate to source URL → `extractTrackerSnapshotFromTab`. Timeout: `CASE_TRACKER_PAGE_LOAD_TIMEOUT_MS = 45000 ms`.
7. If snapshot is shell-only (< `TRACKER_MIN_CONTENT_CHARS = 300` chars), retry with `preferOriginTab: true`.
8. Pass to `applyTrackerSnapshotResult`.

### `applyTrackerSnapshotResult`

1. Check content hash — if identical to `last_content_hash`, update `last_checked_at`, return `{ changed: false }`.
2. If `snapshot.apiCaseState === 'closed'`, short-circuit: mark closed, skip LLM.
3. Call `analyzeTrackerSnapshot` → LLM analysis.
4. Merge analysis fields into tracker and save.
5. Return `{ changed, initial, tracker, result }`.

### S1 Handling

On first ingest (`item.initial === true`) of an S1 tracker:
- Broadcasts `CASE_TRACKER_S1_ALERT` with `{ trackerId, caseNumber, title, accountName }`
- Auto-triggers `orchestrateCaseDeepDive(caseNumber)` unless `deep_dive_pending` or `deep_dive_markdown` already set

---

## 2. TS Tools Support API

**Base URL:** `https://support-api.ts-tools.prod.corp.mongodb.com`

**Auth:** All calls use `chrome.scripting.executeScript` to inject fetch into an authenticated Hub tab at `https://hub.corp.mongodb.com`. Auth is carried by Hub session cookies via `credentials: 'include'`. On 401/403: calls `openHubForAuth()` then retries once with `retryOnAuth: false`.

### Endpoints (18 total)

| Function | Method | Path | Key parameters |
|----------|--------|------|---------------|
| `fetchAccountCasesSummary` | GET | `/account/{id}/cases/summary` | `accountId` |
| `fetchAccountTopics` | GET | `/account/{id}/topics` | `accountId` |
| `fetchAccountTeam` | GET | `/account/{id}/team` | `accountId` |
| `searchCasesByAccount` | POST | `/search/byAccount` | body: `{query, id}` |
| `fetchCaseComments` | POST | `/case/comments` | body: `{case_number, saved_paginate: {page, limit: 50}}`; paginates up to 20 pages |
| `fetchCaseAttachments` | POST | `/files/filter` | body: `{cases: [caseNumber], page: 1, limit: 25}` |
| `fetchAccountSummary` | GET | `/account/{id}/summary` | `accountId` |
| `fetchAccountContacts` | GET | `/account/{id}/contacts` | `accountId` |
| `fetchAccountProducts` | GET | `/account/{id}/products` | `accountId` |
| `fetchAccountRecommendations` | GET | `/account/{id}/recommendations` | `accountId` |
| `fetchAccountActivities` | GET | `/account/{id}/activities` | `accountId` |
| `fetchCaseDetail` | GET | `/case/{number}` | `caseNumber` |
| `fetchAwsSignedLink` | GET | `/aws/signedLink/{key}` | `key` — time-limited S3 pre-signed URL; **never persist** |
| `fetchAtlasClusters` | GET | `/atlas/clusters?projectId={projectId}` | `projectId` |
| `fetchPopularArticles` | GET | `/articles/popular/{count}` | `count` (max 50) |
| `fetchAccountCloudOrgs` | GET | `/account/{id}/cloud_orgs` | `accountId` |
| `escalateCase` | POST | `/case/escalate` | body: `{case_number, markdown_text}` |
| `reopenCase` | POST | `/case/reopen` | body: `{case_number, markdown_text}` |

`fetchCaseEnrichment(caseNumber, accountId)` fires all account + case calls via `Promise.allSettled` and returns all fields above as named keys.

---

## 3. Tracker Record Schema

Persisted in `chrome.storage.local` under `case_trackers[id]`:

```js
{
  // Identity
  id, tracker_id, source_kind: 'cub' | 'jira' | 'generic',
  source_ref, source_url, account_id, customer_name, label,

  // Configuration
  enabled, archived, auto_managed, provider, poll_interval_minutes,
  sound_enabled, needs_refresh,

  // State
  severity, severity_rank, current_title, current_status, change_summary,
  executive_summary, case_state: 'open' | 'closed' | 'unknown', acknowledged,

  // Refresh tracking
  last_checked_at, last_successful_refresh_at, last_content_hash,
  last_list_checked_at, last_list_ingest_at, list_change_key, last_list_change_at,

  // Full analysis output
  timeline: Array<{ date, event }>,
  people: Array<{ name, role, involvement }>,
  next_steps: string[], blockers: string[],
  proposed_solutions: string[], additional_solutions: string[],
  diagnostic_tools: string[],    // keys from DIAGNOSTIC_TOOLS_REGISTRY

  // Timestamps
  created_at, updated_at,

  // Deep dive
  deep_dive_pending: boolean, deep_dive_markdown: string
}
```

---

## 4. Severity System

### `getSeverityRank(severity)` → integer

| Rank | Recognized values (case-insensitive, trimmed) | Poll interval |
|------|----------------------------------------------|--------------|
| 1 | `sev1`, `sev 1`, `severity 1`, `critical` | 1 min |
| 2 | `sev2`, `sev 2`, `severity 2`, `high` | 1 min |
| 3 | `sev3`, `sev 3`, `severity 3`, `medium` | 60 min |
| 4 | `sev4`, `sev 4`, `severity 4`, `low` | 60 min |
| 5 | Any non-empty unrecognized string (including bare `S1`, `S2`) | 60 min |
| 999 | Empty string | 60 min |

**Important:** Bare `S1`, `S2` are intentionally NOT recognized as critical. Only `sev1`/`severity 1` forms are canonical. Adding bare S1/S2 would promote many trackers to 1-minute polling and saturate the service worker alarm budget.

`isCriticalSeverity(severity)` returns `true` if `getSeverityRank(severity) <= 2`.

### List-Driven Refresh Gate

CUB trackers with urgency rank > `CASE_TRACKER_FULL_REFRESH_SEVERITY_RANK = 1` skip the full Hub API refresh when the list-row signal is unchanged. Gate bypassed for: S1 cases, `force: true`, `reason === 'initial'`.

---

## 5. LLM Analysis Pipeline

### Provider selection

```
DEFAULT_PROVIDER_ORDER = ['copilot_cli', 'copilot', 'claude', 'gemini', 'gemini_cli', 'glean']
VALID_LLM_PROVIDERS    = ['copilot_cli', 'claude', 'gemini', 'gemini_cli', 'glean', 'copilot']
```

### Prompt construction

- **`buildCaseTrackerPrompt`** — stripped JSON-only base prompt
- **`buildCaseTrackerAnalysisPrompt`** — enriched: base + skill bundle + customer context + `buildPromptEnvelope`

### Retry logic

1. First attempt: enriched prompt.
2. On JSON parse failure: retry with bare `buildCaseTrackerPrompt`.
3. On double failure: `normalizeAnalysis({ executive_summary: errorText })` — all arrays empty, `case_state: 'unknown'`.

### `normalizeAnalysis` output schema

```json
{
  "title": "string",
  "case_state": "open | closed | unknown",
  "current_status": "string",
  "change_summary": "string",
  "executive_summary": "string",
  "timeline": [{ "date": "string", "event": "string" }],
  "people": [{ "name": "string", "role": "string", "involvement": "string" }],
  "next_steps": ["string"],
  "blockers": ["string"],
  "proposed_solutions": ["string"],
  "additional_solutions": ["string"],
  "diagnostic_tools": [{ "tool_id": "string", "variables": {}, "why": "string" }]
}
```

`normalizeCaseState` maps: `closed/resolved/won't fix` → `'closed'`; `open/new/in progress/working/pending/waiting` → `'open'`; anything else → `'unknown'`.

---

## 6. Diagnostic Tools Registry

15 tools. The `diagnostic_tools` array in analysis output contains 3–7 keys. `buildDiagnosticToolsCatalogForPrompt()` serializes the registry for LLM prompts.

| tool_id | Name | Category | When to use |
|---------|------|----------|------------|
| `atlas_realtime_performance` | Atlas Real-Time Performance Panel | atlas | slow queries, high CPU, latency spike, timeout, connection saturation |
| `atlas_performance_advisor` | Atlas Performance Advisor | atlas | slow queries, missing index, COLLSCAN, performance tuning |
| `atlas_query_profiler` | Atlas Query Profiler | atlas | slow queries, full collection scans, high docsExamined |
| `atlas_cluster_metrics` | Atlas Cluster Metrics | atlas | memory pressure, disk IOPS, storage full, replication lag |
| `atlas_logs` | Atlas Cluster Logs | atlas | error in logs, OOM, crash, election, network partition |
| `atlas_network_access` | Atlas Network Access | atlas | connection refused, IP allowlist, VPC peering, PrivateLink |
| `atlas_database_access` | Atlas Database Access | atlas | auth failed, permissions, LDAP, X.509 |
| `atlas_backup_restore` | Atlas Backup & Restore | atlas | backup failure, restore, snapshot, data loss |
| `mongosh_diagnostics` | MongoDB Shell Diagnostics | mongodb | live diagnostics, currentOp, serverStatus, explain plan |
| `mongodb_compass` | MongoDB Compass | mongodb | explain plan, schema analysis, index management |
| `mongodb_docs` | MongoDB Documentation | mongodb | documentation lookup, API reference |
| `mongodb_jira` | MongoDB Jira | mongodb | known bug, SERVER ticket, workaround |
| `hub_account_page` | Customer Hub Account Page | hub | account overview, open cases, contract, ARR |
| `ts_tools_case` | TS Tools Case View | ts-tools | case history, attachments, TS Tools, escalation |
| `atlas_search_explorer` | Atlas Search Explorer | search | Atlas Search, $search, Lucene, search index |

**Template variable inference:** `project_id` from Atlas project URL segment; `cluster_name` from page text or case subject; `case_number` from `tracker.source_ref` (8-digit `0XXXXXXX`); `account_id` from `tracker.account_id`; `ticket_id` from `tracker.source_ref` when `source_kind === 'jira'`. Variables that cannot be inferred must be left as `""`.

---

## 7. Automatic Sync and Passive Capture

### `syncAutomaticCaseTrackers`

Called at the start of every poll cycle and every batch refresh:
1. Load all open cases from corpus.
2. For each: create/update `cub` tracker; for every `HELP-NNNNN` reference in `case.full_text`, create/update `jira` tracker.
3. Archive auto-managed trackers whose case is no longer open.

### `ingestPassiveCaseTrackerSnapshot`

Origin policy: `hub-or-dashboard`. Process:
1. Extract case number from `snapshot.url` or `snapshot.caseNumber`.
2. Find all matching trackers; call `applyTrackerSnapshotResult` for each.
3. Returns `{ matched, processedCount, updatedCount, results }`.

### `markCaseTrackerStale`

Sets `needs_refresh: true` on all enabled, non-archived trackers matching `caseNumber`. Triggered by `TSTOOLS_CASE_STALE` event.

---

## 8. Skill Injection and LLM Context Assembly

### Pinned skills (always included)

```
PINNED_CASE_SKILL_IDS = ['mongodb-kb', 'mongodb-developer']
```

### Scored skill selection

Up to 4 additional skills on top of the pinned entries. `selectRelevantBundledSkillsForCaseTracker()` counts keyword hits across base prompt + page snapshot + customer context file. `mongodb-expert` gets a +2 bonus on `cub` and `jira` trackers.

### Cache key for prompt envelope

```
case_tracker:{account_id}:{comma-separated selected skill IDs}
```

### Case-relevant skills from `BUNDLED_SKILL_CATALOG`

| Skill ID | Representative keywords |
|----------|------------------------|
| `mongodb-kb` | knowledge base, KB, known issue, bug, workaround |
| `mongodb-developer` | driver, connection, SDK, CRUD, connection string |
| `mongodb-expert` | aggregation, schema, index, replica set, sharding |
| `mongodb-atlas-expert` | Atlas, cluster, project, M10, serverless, flex |
| `mongodb-performance-troubleshooting` | slow query, performance, explain, index scan, COLLSCAN |
| `tstools-reference` | ts-tools, support API, case, escalation, socket |

---

## 9. Service-Worker Message API

All messages require origin `extension-only` except `PASSIVE_CASE_TRACKER_SNAPSHOT` (requires `hub-or-dashboard`).

| Message type | Payload | Response |
|-------------|---------|---------|
| `CREATE_CASE_TRACKER` | `{ accountId, sourceInput, label?, provider?, pollIntervalMinutes?, soundEnabled? }` | `{ success, tracker, result }` |
| `UPDATE_CASE_TRACKER` | `{ trackerId, changes: Partial<TrackerRecord> }` | `{ success, tracker }` |
| `ACKNOWLEDGE_CASE_TRACKER` | `{ trackerId }` | `{ success, tracker }` |
| `DELETE_CASE_TRACKER` | `{ trackerId }` | `{ success }` |
| `REFRESH_CASE_TRACKER` | `{ trackerId, force?, provider? }` | `{ success, result, tracker }` |
| `REFRESH_CASE_TRACKERS` | `{ accountId?, force?, provider? }` | `{ success, results }` |
| `CANCEL_REFRESH_CASE_TRACKERS` | `{}` | `{ success }` |
| `OPEN_CASE_TRACKER_SOURCE` | `{ trackerId, active?, reload? }` | `{ success, tracker, tab }` |
| `OPEN_ALL_CASE_TRACKER_SOURCES` | `{ accountId?, active?, reload?, sourceKind? }` | `{ success, sourceKind, openedCount, skippedCount, results }` |
| `REMOVE_DUPLICATE_CASE_TRACKERS` | `{ accountId? }` | `{ success, groupsChecked, removedCount, results }` |
| `OPEN_CASE_TRACKER_OVERLAY` | `{ accountId? }` | `{ success, reused, windowId, tabId, url }` |
| `PASSIVE_CASE_TRACKER_SNAPSHOT` | `{ snapshot, provider? }` | `{ success, matched, processedCount, updatedCount, results }` |

**Broadcast-only messages** (no request/response):

| Type | Fields | Sent by |
|------|--------|---------|
| `CASE_TRACKER_PROGRESS` | `{ stage, trackerId?, accountId?, currentIndex?, total?, changedCount?, errorCount? }` | service-worker during refresh |
| `CASE_TRACKER_S1_ALERT` | `{ trackerId, caseNumber, title, accountName }` | service-worker on first S1 ingest |

---

## 10. Storage and Mutation Queuing

- All tracker records in `chrome.storage.local['case_trackers']` as `{ [trackerId]: trackerRecord }`.
- `queueCaseTrackerMutation(key, work)` serializes writes per tracker ID using a `Map<trackerId, Promise>` to prevent race conditions.
- **Limitation:** queue is in-memory only. Service worker suspension or Chrome restart drops all queued mutations.

---

## 11. Known Limitations

| Issue | Consequence | Mitigation |
|-------|------------|-----------|
| Auth dependency on readable tabs | Login page captured as snapshot text; produces misleading analysis | `shouldProcessPassiveTrackerSnapshot()` rejects snapshots with `missing-source-updated`; no automatic recovery — TAM must re-authenticate and force-refresh |
| Hub DOM extraction brittleness | Hub UI redesign breaks selector matching; `source_updated_at` and `source_text_full` may be empty | `source_text_hash` (FNV-1a) detects text changes independently of parsed timestamp |
| List-driven refresh gate delay | Up to corpus ingest interval delay between case update and tracker detecting it (non-S1 cases) | Key fields: `list_change_key`, `last_list_checked_at`, `last_list_ingest_at`, `last_list_change_at` |
| Duplicate tracker deduplication | Lower-quality duplicate persists until user removes it or next sync | `findDuplicateCaseTrackerGroups()` + `chooseCanonicalCaseTracker()` via quality score |
| URL normalization edge cases | Two Atlas URLs differing only in fragment treated as same tracker; `inferSourceKind` rejects non-Hub, non-Jira URLs | `normalizeTrackerUrl()` strips `#…` and trailing slashes |
| `HELP_TICKET_REGEX` extraction scope | `/\bHELP-\d+\b/gi` won't match SERVER tickets; TAM must create Jira tracker manually for SERVER cross-references | By design — avoid false positives from other ticket systems |
