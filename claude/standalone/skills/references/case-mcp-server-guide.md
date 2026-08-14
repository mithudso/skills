<!-- hub-reference-banner -->
> **Reference file — part of the `misc-catch-all` hub.** Formerly the standalone `case-mcp-server-guide` skill.
> Sibling topics in this family are now reference files under the hubs (`misc-catch-all`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: case-mcp-server-guide
title: Case MCP Server Guide
description: |
  Guidance for starting, using, and troubleshooting the MDB Case Assistant local case MCP server and its mdb_case_* tool surface.
  TRIGGER: deciding whether to use the local case MCP server; starting the case MCP server; choosing the right mdb_case_* tool for a case/account/HELP workflow; understanding the trust/auth model; troubleshooting helper offline, worker offline, Hub auth, Jira auth, or vault-locked tracking state; finding the right local commands for the helper relay or case MCP server.
  SKIP: popup/options/dashboard/overlay DOM inspection (use Chrome DevTools MCP or Playwright MCP); browser console or network debugging; code changes inside the repo (read the code directly); worker architecture, message routing, or implementation details that are better answered by reading the code.
category: custom
version: "1.1.0"
updated: "2026-05-29"
keywords:
  - case MCP server
  - mdb_case
  - MDB Case Assistant
  - case lookup
  - account lookup
  - HELP ticket
  - Hub auth
  - Jira auth
  - helper relay
  - mcp case server
  - case workflow
  - tracking state
  - vault locked
  - worker offline
when_to_use:
  - "start the case MCP server"
  - "which mdb_case tool should I use"
  - "case MCP server is offline"
  - "helper relay not responding"
  - "Hub auth not working in case MCP"
  - "Jira auth failure on mdb_case_get_help_ticket"
  - "vault locked in case tracking"
  - "worker not connected to case MCP"
  - "mdb_case_get_server_status"
  - "decide whether to use case MCP or another workflow"
related_skills:
  - chrome-devtools-mcp:chrome-devtools
  - tam-operations
origin: local
---

# Case MCP Server Guide

Guidance for starting, using, and troubleshooting the MDB Case Assistant local case MCP server and its `mdb_case_*` tool surface.

## When not to use this skill

- **DOM/UI inspection** → use Chrome DevTools MCP or Playwright MCP
- **Browser console or network debugging** → use Chrome DevTools MCP
- **Code changes inside the repo** → read the code directly
- **Worker architecture, message routing, or implementation details** → read the code directly

## Skill guidance

- Treat `docs/case-mcp-server-context.md` in the `mdb-case-assistant` repo as the primary reference.
- Keep the local case MCP server separate from the shipped extension runtime in explanations.
- State plainly that the case MCP server is **read-only today** — do not invent MCP mutation tools.
- Do not describe browser-debug MCP configs as runtime dependencies.

## Server dependencies

The case MCP server requires all of the following to be running/available:

1. `npm run dev:extension` — the extension helper relay
2. The unpacked extension loaded in Chrome
3. At least one extension surface opened so the service worker connects
4. Authenticated Hub / Support / Jira browser sessions as needed by the specific tool

## Preferred workflow rules

Use tools in this order when diagnosing or querying:

| Step | Tool | When |
|---|---|---|
| 1 | `mdb_case_get_server_status` | Always first when diagnosing availability |
| 2 | `mdb_case_get_support_auth_status` | Before assuming case lookup is broken |
| 3 | `mdb_case_get_case` | For best-available case context |
| 4 | `mdb_case_get_case_comments` | When comment history matters — respect page/limit |
| 5 | `mdb_case_get_case_stage` + `mdb_case_get_case_next_action` | For precise workflow state |
| 6 | `mdb_case_get_account` | First when you only have a query or need account resolution |
| 7 | `mdb_case_list_account_cases` | When you already have a confirmed account ID |
| 8 | `mdb_case_search` | For Support-backed discovery across accounts, projects, and cases |
| 9 | `mdb_case_get_help_ticket` | Only for explicit HELP keys |
| 10 | `mdb_case_build_evidence_snapshot` | When tracker-style evidence packaging is needed |
| 11 | `mdb_case_get_tracked_case_analysis` | Only for already-stored tracked analysis |

## Startup commands

```bash
npm install
npm run dev:extension
npm run mcp:case-server
npm run mcp:case-server:build
npm run dev:helper -- status
npm run dev:helper -- get-logs --level error --limit 50
```

## Troubleshooting

| Symptom | Diagnosis | Fix |
|---|---|---|
| Helper offline | `mdb_case_get_server_status` shows helper down | Start `npm run dev:extension` |
| Worker offline | Helper up but no worker connected | Open popup, options page, dashboard, or a supported case page |
| Missing Hub / Support auth | `mdb_case_get_support_auth_status` returns auth failure | Sign into Hub / Support in the same Chrome profile |
| Missing Jira auth | `mdb_case_get_help_ticket` returns auth-style failure | Sign into `jira.mongodb.org` in the same Chrome profile |
| Vault locked | `mdb_case_get_tracking_state` fails | Unlock the repo vault through the extension workflow |
| Not found vs auth vs connectivity | Ambiguous error | Classify separately — cite the exact `mdb_case_*` tool and error returned |

## Bundled context

Source: `docs/case-mcp-server-context.md` in the `mdb-case-assistant` repo
