<!-- hub-reference-banner -->
> **Reference file — part of the `tam-operations` hub.** Formerly the standalone `tstools-reference` skill.
> Sibling topics in this family are now reference files under the hubs (`tam-operations`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: tstools-reference
description: >
  TS Tools Support API Integration Reference — all 92 endpoints grouped by domain, full
  request/response shapes for case lifecycle, search, attachments, account, project, and page
  bootstrap routes, authentication (cookie vs Bearer), Socket.IO realtime event model, and
  integration patterns for the Customer Dashboard extension.
  TRIGGER: questions about TS Tools API endpoint shapes, case lifecycle REST routes, attachment
  upload flow, account endpoint behavior, socket event model, TS Tools auth schemes (cookie vs
  Bearer), page initial-load bootstrap payload, or any Support Portal / Customer Hub backend
  integration constraint.
  SKIP: Chrome extension architecture (use chrome-extension-expert, references/chrome-dev.md), Jira
  API (use integration-clients, references/jira-extension-client.md), MongoDB query/performance (use
  mongodb-expert), Glean/LLM pipelines (use integration-clients, references/glean-dev.md),
  implementation patterns and auth cascade code (use tam-operations, references/ts-tools-support-api.md
  — it covers the same domain with executable code patterns and a checklist).
version: 1.1.1
updated: "2026-06-01"
origin: local
category: developer
tags:
  - ts-tools
  - support-api
  - case-management
  - customer-hub
  - rest-api
  - socket-io
  - authentication
triggers:
  - TS Tools API
  - support-api endpoint
  - case lifecycle endpoint
  - GET /case
  - POST /case
  - socket.io case update
  - hub.corp.mongodb.com
  - page initial load
  - account endpoint
  - attachment upload TS Tools
related_skills:
  - ts-tools-support-api
  - case-tracker
  - case-mcp-server-guide
  - chrome-extension-expert
---

# TS Tools Support API Integration Reference

Generated from `docs/tstools-context.md` in `10gen/mdb-tam`. Start from the bundled context below; defer to that source file for exact runtime assumptions, auth behavior, and integration constraints.

## When to use this skill

Use when the request needs Support Portal / Customer Hub / TS Tools backend integration work — case lifecycle, attachment, search, summarization, socket-driven case freshness, or endpoint discovery.

**Use `tam-operations` (references/ts-tools-support-api.md) instead** when you need implementation patterns (auth cascade code, response normalizers, pagination helpers, error handling, checklists). This file is the endpoint reference; that one is the implementation guide.

**Do not use** for:
- Chrome extension architecture → `chrome-extension-expert` (references/chrome-dev.md)
- Jira API integration → `integration-clients` (references/jira-extension-client.md)
- MongoDB query/performance → `mongodb-expert`
- Glean/LLM analysis pipelines → `integration-clients` (references/glean-dev.md)

---

## Bundled Reference

Source: `docs/tstools-context.md`

---
title: TS Tools Support API Integration Reference
source_documents:
  - tstoolsapi.txt
  - support-api-swagger.json
synthesis_version: 2
synthesis_date: 2026-05-12
dedupe_policy: prefer_concrete_runtime_details
audience: LLM retrieval and case-management backend integration guidance
---

## 0. Scope

This document is the authoritative integration reference for the TS Tools Support API backend. It covers all 92 endpoints discovered from the live OpenAPI spec at `/public/v1.0/api-docs` (staging), grouped by domain, with full request/response shapes for the case-management routes most relevant to the Customer Dashboard extension.

## 1. Service Identity

| Property | Value |
|---|---|
| Internal production host | `https://support-api.ts-tools.prod.corp.mongodb.com` |
| Staging host | `https://support-api.ts-tools.staging.corp.mongodb.com` |
| Public REST base | `//api.support.mongodb.com/public/v1.0` |
| OpenAPI spec | `GET /public/v1.0/api-docs` (requires auth; `/docs` renders ReDoc UI) |
| Realtime socket | `wss://support-api.ts-tools.prod.corp.mongodb.com/socket.io/?EIO=4&transport=websocket` |
| Contact | tstools@mongodb.com, Slack `#ts-supporttools`, Jira `TSTOOLS` |

## 2. Authentication

Two auth schemes are accepted on nearly every route:

| Scheme | Header / Cookie |
|---|---|
| Bearer JWT | `Authorization: Bearer <token>` |
| Cookie | `support`, `stg_support`, `dev_support`, `HUB`, `stg_HUB`, `dev_HUB`, `internal_user` |

### `POST /auth`
Authenticate with an API username + key to obtain a Bearer token.
```json
Request:  { "username": "test@email.com", "key": "EREPIJ..." }
Response: { "token": "..." }
```

### `GET /status/auth`
Check if the current session/token is valid. Returns `200` when authenticated.

**Integration note:** For first-party integrations running inside an authenticated Support Portal / Customer Hub session, prefer **cookie auth** — the session cookies are already present and do not require a separate token exchange.

## 3. Domain Overview (all 92 paths)

| Domain | Count | Key paths |
|---|---|---|
| Case lifecycle | 10 | `/case`, `/case/{number}`, `/case/resolve`, `/case/reopen`, `/case/escalate` |
| Case search/filter | 3 | `/cases`, `/cases/search`, `/case/comments` |
| Attachments/files | 7 | `/case/attachments`, `/case/add/attachments`, `/case/delete/attachment`, `/case/{number}/attachments`, `/files/bycases`, `/files/filter`, `/aws/*` |
| Account | 11 | `/account/{id}/summary`, `/account/{id}/cases/summary`, `/account/{id}/activities`, `/account/{id}/team`, `/account/{id}/contacts`, `/account/{id}/products`, `/account/{id}/resources`, `/account/{id}/recommendations`, `/account/{id}/topics`, `/account/{id}/request-services`, `/accounts`, `/accounts/search` |
| Project | 13 | `/project/{id}`, `/project/{id}/casecreation`, `/project/{id}/productsandservices`, `/project/{id}/users`, `/project/{id}/mailinglist*`, `/projects`, `/projects/search`, `/activeProjects` |
| Search | 5 | `/search`, `/search/byAccount`, `/cases/search`, `/projects/search`, `/accounts/search` |
| Summarization | 2 | `/account/{id}/summary`, `/internal/account/{id}/summary` |
| User/profile | 10 | `/user/profile`, `/user/create`, `/user/{id}`, `/user/{id}/projects`, `/user/{id}/addprojects`, `/user/{id}/removeprojects`, `/users`, etc. |
| Atlas | 1 | `/atlas/clusters` |
| Articles/KB | 4 | `/article/{number}`, `/articles/recent`, `/articles/popular`, `/articles/recent/{count}` |
| AI/chat | 2 | `/ai/chat/{id}/sendchatfeedback`, `/ai/chat/{id}/sendsearchfeedback` |
| Alerts | 1 | `/alerts` |
| AWS/S3 | 7 | `/aws/signedUpload`, `/aws/signedRequest`, `/aws/signedLinks`, `/aws/signedLink/{key}`, `/aws/signedProcessedLink/{key}`, `/aws/filenameversion/{key}`, `/aws/pixel`, `/aws/cli` |
| Page | 1 | `/page/initial` |
| Misc | 5 | `/coveo/token/generate`, `/survey/*`, `/redirect/case/{id}`, `/metadata/phones` |

## 4. Case Lifecycle — Full Endpoint Reference

### `POST /case` — Create Case
**Auth:** Bearer or Cookie
```json
{
  "project_id": "IdPattern",
  "groupname": "string",
  "product": "string",
  "severity": "S1|S2|S3|S4",
  "subject": "string",
  "description": "string",
  "triagecomponents": "string",
  "storage": "string",
  "mdbDrivers": "string",
  "version": "string"
}
```

### `GET /case/{number}` — View Case
**Auth:** Bearer or Cookie. Path param: `number` (CaseNumber pattern, e.g. `00000000`). Returns full case object.

### `POST /case/{number}` — Update Case
**Auth:** Bearer or Cookie. All body fields optional:
```json
{
  "casenumber": "CaseNumber",
  "project": "string",
  "severity": "S1|S2|S3|S4",
  "subject": "string",
  "description": "string",
  "options": {}
}
```

### `POST /case/resolve` — Resolve a Case
**Auth:** Bearer or Cookie
```json
{
  "case_number": "CaseNumber",
  "reason": "CaseClosedReason",
  "markdown_text": "string",
  "userLocale": "string"
}
```

### `POST /case/reopen` — Reopen a Case
**Auth:** Bearer or Cookie
```json
{
  "case_number": "CaseNumber",
  "markdown_text": "string",
  "userLocale": "string"
}
```

### `POST /case/escalate` — Escalate a Case
**Auth:** Bearer or Cookie
```json
{
  "case_number": "CaseNumber",
  "markdown_text": "string",
  "userLocale": "string",
  "phone": "string"
}
```

### `POST /case/{number}/comment` — Create Case Comment
**Auth:** Bearer or Cookie
```json
{
  "case_number": "CaseNumber",
  "markdown_text": "string",
  "partner_visibility": "partner|public",
  "paging": {},
  "userLocale": "string"
}
```

### `POST /case/comments` — Paged Case Comments (read)
**Auth:** Bearer or Cookie
```json
{
  "case_number": "CaseNumber",
  "saved_paginate": "PaginateObject"
}
```

## 5. Case Search and Filter

### `POST /cases` — Filter Cases (paged list)
```json
{
  "paging": {
    "rows": 25,
    "page": 1,
    "sortBy": "string",
    "status": "OPEN",
    "type": "SUPPORT",
    "action": "next"
  },
  "filter": { "text": "string" }
}
```

### `POST /cases/search` — Full-text Case Search
**Request body:** `{ "query": "search string" }` — Returns array.

### `POST /search` — Global Search (Accounts + Projects + Cases)
**Request body:** `{ "query": "search string" }` — Returns array.

### `POST /search/byAccount` — Search Within Account
**Request body:** `{ "query": "string", "id": "IdPattern" }` — Returns array.

## 6. Attachment and File Endpoints

### `POST /case/attachments` — List Attachments
```json
{ "case_number": "CaseNumber", "paging": "PagingObject", "filter": "FilterObject" }
```

### `POST /case/add/attachments` — Add Attachments
```json
{ "case_number": "CaseNumber", "attached": [] }
```

### `POST /case/delete/attachment` — Delete Attachment
```json
{ "case_number": "CaseNumber", "attachmentId": "IdPattern", "fileName": "string" }
```

### AWS S3 Upload Flow
1. `POST /aws/signedUpload` — get pre-signed S3 upload fields
2. PUT directly to S3 with the signed fields
3. `POST /aws/signedRequest` — confirm/process the upload
4. `GET /aws/signedLink/{key}` — time-limited download URL
5. `GET /aws/signedProcessedLink/{key}` — download link for a processed/extracted file

## 7. Account Endpoints

| Path | Purpose |
|---|---|
| `GET /account/{id}/cases/summary` | Case count/status breakdown for dashboard widgets |
| `GET /account/{id}/summary` | Structured account health data (public) |
| `GET /internal/account/{id}/summary` | Richer data, internal auth required |
| `GET /account/{id}/team` | MongoDB internal team (TAM, CSM) assigned to account |
| `GET /account/{id}/contacts` | Support contacts |
| `GET /account/{id}/products` | Products and services |
| `GET /account/{id}/resources` | Account resources |
| `GET /account/{id}/recommendations` | Account recommendations (optional `?limit=N`) |
| `GET /account/{id}/topics` | Common themes across support cases |
| `GET /accounts` | List accessible accounts |
| `GET /accounts/search` | Search accounts by query |

## 8. Project Endpoints

| Path | Purpose |
|---|---|
| `GET /project/{id}` | View project details |
| `GET /project/{id}/casecreation` | Form options for case creation (products, severities) |
| `GET /project/{id}/productsandservices` | Product/services options |
| `POST /project/{id}/users` | Get project users with paging/filter |
| `GET /activeProjects` | List active projects |
| `POST /projects` | Get projects (paginated) |
| `POST /projects/search` | Search projects |

## 9. Page Initial Load

### `POST /page/initial`
Get all data needed to bootstrap a support portal page:
```json
{
  "case": "00000000",
  "file": "myfilename.zip",
  "project": "5003000000D8cuI",
  "account": "5003000000D8cuI"
}
```
Returns a consolidated payload (case, project, account, user, options).

## 10. Data Schemas

### Core Identifiers
- **IdPattern** — Salesforce 15/18-char ID
- **MongoIdPattern** — MongoDB ObjectId string
- **CaseNumber** — numeric string, 8 digits with leading zero (e.g. `"00000000"`)

### Case (`FullCase`)

| Field | Type | Required |
|---|---|---|
| case_number | CaseNumber | ✓ |
| description | string | ✓ |
| subject | string | ✓ |
| severity | Severity | ✓ |
| project | IdPattern | ✓ |
| project_name | string | ✓ |
| status | CaseStatus | ✓ |
| product_id | string | |
| group_name | string | |
| created_date | ISO8601 string | |
| last_modified_date | ISO8601 string | |
| updated_date | ISO8601 string | |

### Enums
- **Severity:** `S1`, `S2`, `S3`, `S4`
- **CaseStatus:** `New`, `On Hold`, `In Progress`, `Waiting for Customer`, `Waiting for Development`, `Waiting for Feedback`, `Closed`, `Resolved`
- **CaseClosedReason:** `This issue is resolved`, `This issue has gone away`, `Other`
- **Role:** `Admin`, `User`, `Read-only User`, `Restricted User`, `Owner`

### PagingObject (list endpoints)
```json
{
  "rows": 25, "page": 1, "sortBy": "string",
  "status": "OPEN", "type": "SUPPORT", "action": "next",
  "hasNextPage": true, "hasPrevPage": false,
  "totalDocs": 142, "totalPages": 6
}
```

### PaginateObject (comment endpoints)
```json
{
  "hasNextPage": true, "hasPrevPage": false, "limit": 50,
  "nextPage": 2, "prevPage": null, "page": 1,
  "pagingCounter": 1, "returnedDocs": 50,
  "totalDocs": 127, "totalPages": 3
}
```

### AttachmentObject
```json
{
  "attachmentId": "IdPattern",
  "name": "string",
  "contentType": "string",
  "fileSize": 0,
  "createdDate": "ISO8601",
  "lastModifiedDate": "ISO8601",
  "createdBy": "string",
  "isPrivate": false,
  "isIncluded": false,
  "s3file": {}
}
```

## 11. Realtime / Event Model

- **Transport:** Socket.IO / Engine.IO at `wss://support-api.ts-tools.prod.corp.mongodb.com/socket.io/?EIO=4&transport=websocket`
- **Upstream channel:** Salesforce platform event `/event/Case_Updated__e`
- **Events:** `Case_Updated` (a case field changed), `notifyCaseViewers` (presence), `ALL_Updated` (broad signal)
- **Cookie names accepted:** `dev_support`, `stg_support`, `support`, `internal_user`, `dev_HUB`, `stg_HUB`, `HUB`

**Flow:**
1. Connect with authenticated support session cookie
2. Emit `joinCase` (single) or `joinCases` (array of case numbers)
3. On `Case_Updated`: use as stale signal, then REST-fetch authoritative state via `GET /case/{number}`

**Critical:** Socket events are stale signals only. Never update UI state directly from socket events.

## 12. Integration Patterns for the Customer Dashboard

### 12.1 Case Refresh Flow
1. Connect to socket, emit `joinCases` for tracked case numbers
2. On `Case_Updated`: mark tracker stale, call `GET /case/{number}` to refresh
3. Use `POST /cases` with `status: "OPEN"` + paging for bulk list sync

### 12.2 Account Context Loading
1. `POST /page/initial` with `account` ID — gets consolidated bootstrap payload
2. `GET /account/{id}/cases/summary` — dashboard case-count widget
3. `GET /account/{id}/team` — show assigned MongoDB team members (TAM, CSM)
4. `GET /account/{id}/topics` — surface common case themes

### 12.3 Case Creation Flow
1. `GET /project/{id}/casecreation` — fetch valid product/severity/group options
2. `POST /case` — create with required fields
3. `POST /aws/signedUpload` + PUT to S3 + `POST /case/add/attachments` — attach diagnostics

### 12.4 Attachment / Diagnostic Workflow
1. `POST /files/filter` — find existing diagnostic files for the case
2. `POST /aws/signedUpload` — get S3 pre-signed fields
3. PUT file to S3
4. `POST /case/add/attachments` — attach the uploaded file to the case
5. `GET /aws/signedLink/{key}` — share a time-limited download URL

## 13. Operational Constraints

- `/docs` (ReDoc UI) and `/public/v1.0/api-docs` require auth — no anonymous discovery in production
- Production docs page may be disabled; keep route references explicit in code
- The websocket endpoint is a transport endpoint, not a REST base URL
- `POST /auth` with `username`+`key` is for API key auth; first-party integrations should use existing portal cookies
- Case number format: `0xxxxxxx` (8 digits, leading zero for CUB)
- All comment text (`markdown_text`) accepts Markdown formatting
- `partner_visibility` on comments: `"partner"` = MongoDB-internal only; `"public"` = visible to customers

## 14. Pre-Build Checklist

Before building any TS Tools integration:

1. Identify whether the workflow needs REST, socket, or both
2. Verify the exact auth mode (cookie vs Bearer) for the target route
3. Confirm a valid internal/portal session is present in the browser context
4. Validate with `GET /status/auth` before making case-mutating calls
5. Use `GET /project/{id}/casecreation` to get valid enum values before `POST /case`
6. Wrap all diagnostic uploads in the full AWS S3 flow, not a direct attachment POST
