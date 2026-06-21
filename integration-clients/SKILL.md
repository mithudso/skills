---
description: >-
  Third-party SaaS API & integration-client hub — Jira, Monday.com, Slack, Salesforce, Glean, Aha!. REST/GraphQL surfaces, auth (OAuth 2.0, cookie/session, API tokens), webhooks/events, SDKs/CLIs, app frameworks, calling from Chrome MV3. TRIGGER: client against a third-party SaaS API; Jira dev (REST v3, issues/search/workflows, agile) incl. from a Chrome extension; Monday.com GraphQL v2 + Monday apps (Apps Framework, mapps CLI); Slack apps (Web API, Events API, Socket Mode, Block Kit); Salesforce (REST/SOAP/Bulk 2.0/Composite/Tooling/Pub-Sub/GraphQL); Glean (CLI, MCP, SDK clients); Aha! API/webhooks. SKIP: vendor-neutral API DESIGN patterns → software-engineering-patterns; AI/LLM tool-calling → ai-mcp-sdk-prompting; generic Chrome OAuth plumbing → chrome-extension-expert; TAM ops audits (Monday board / Slack subscription) → tam-operations.
name: integration-clients
category: developer
tags:
  - developer
  - integration
  - api-client
  - jira
  - monday
  - slack
  - salesforce
  - glean
whenToUse:
  - "building a REST or GraphQL client against Jira, Monday, Slack, Salesforce, Glean, or Aha!"
  - "wiring OAuth 2.0, cookie/session, or API-token auth for a third-party SaaS platform"
  - "calling a vendor API from a Chrome MV3 extension (background/content-script context)"
  - "handling vendor webhooks, Events API, or Socket Mode event streams"
  - "developing a platform app: Monday Apps Framework, Slack app, Salesforce package"
  - "formatting vendor-specific payloads (Monday column-value JSON, Slack Block Kit)"
  - "using a vendor SDK or CLI (Glean SDK/MCP, Monday mapps, Salesforce tooling)"
  - "debugging auth, rate-limit, or pagination issues against an external SaaS API"
related_skills:
  - software-engineering-patterns
  - ai-agent-engineering
  - chrome-extension-expert
version: "1.1.1"
updated: "2026-05-31"
---

# Integration & API Clients

The hub for building and maintaining **integration clients against third-party
SaaS product platforms** — Jira, Monday.com, Slack, Salesforce, Glean, and Aha!.
It covers the concrete, vendor-specific work of talking to each platform: its
REST or GraphQL API surface, its auth model (OAuth 2.0, cookie/session, API
tokens), its webhook/event delivery, its SDKs and CLIs, and — where relevant —
the quirks of calling that API from a Chrome MV3 extension context rather than a
trusted server.

This is the vendor-detail layer. It assumes you have already decided *that* you
need to integrate; it answers *how* to do it correctly for a specific platform.

## How to use this skill

This skill consolidates **8 integration sub-skills** as on-demand reference
files under `references/`. Match the task to the routing table below and **Read
the listed `references/…md` file before answering deep questions** — the table
alone is not enough for endpoint, auth-flow, payload-shape, or CLI detail. For
exact field names, scopes, and version specifics, treat each vendor's official
developer docs as the source of truth.

## Sub-skill routing table

This hub absorbs 8 former standalone skills as on-demand reference files. When a task matches a row, **Read the listed `references/` file** before answering — do not rely on this table alone for depth.

| Sub-topic | When to load | Reference file |
| --- | --- | --- |
| `aha-api` | Use when integrating with Aha! product management platform via REST API, webhooks, or Chrome | `references/aha-api.md` |
| `glean-dev` | Glean enterprise developer reference — CLI, MCP, TypeScript/Python SDK clients, | `references/glean-dev.md` |
| `jira-developer-expert` | Jira Cloud developer expert: REST API v3 (issues, search, workflows), Jira Software agile API | `references/jira-developer-expert.md` |
| `jira-extension-client` | Jira REST API integration from Chrome MV3 extensions — cookie auth, OAuth 2.0 | `references/jira-extension-client.md` |
| `monday-api` | Monday.com GraphQL API v2 — queries, mutations, column value JSON formats, | `references/monday-api.md` |
| `monday-dev` | Monday.com platform architecture, Apps Framework, CLI (mapps), Vibe Design | `references/monday-dev.md` |
| `salesforce-developer-expert` | Salesforce platform development expert — APIs (REST, SOAP, Bulk 2.0, Composite, Tooling, Pub/Sub, GraphQL), | `references/salesforce-developer-expert.md` |
| `slack-dev` | Slack Platform developer reference — Web API, Events API, Socket Mode, Block | `references/slack-dev.md` |

## Cross-hub boundaries

This hub owns vendor-specific integration-client development. Hand off when the
task is not about a specific vendor's API:

- **General, vendor-neutral API design** (pagination strategy, idempotency keys,
  error-envelope conventions, REST-vs-RPC style, retry/backoff patterns not tied
  to one platform) → `software-engineering-patterns`.
- **AI / LLM / agent integration** (tool-calling, agent orchestration, RAG over
  these data sources) → `ai-agent-engineering`.
- **Generic Chrome extension OAuth/identity plumbing** (`chrome.identity`,
  `launchWebAuthFlow`, token caching) not tied to a specific vendor API →
  `chrome-extension-expert` (references/chrome-identity-oauth.md). This hub owns the vendor-API call; the extension
  identity flow that fronts it lives there.
- **TAM operational workflows that *use* these platforms** — Monday board
  audits, Slack subscription audits, account/report generation — stay in their
  standalone / tam-ops skills (`monday-board-audit`,
  `slack-subscription-auditor`, `tam-account-reports`). This hub builds the
  client; those skills run the recurring operational task on top of it.

Some topics straddle a boundary (e.g. a Monday GraphQL query embedded inside a
board-audit workflow). Lead with the hub that matches intent: if you are
building or fixing the *client/integration*, stay here; if you are running the
*operational task*, route to tam-ops.

<!-- cross-hub-map -->
## Cross-hub map — where every integration topic lives

This family is split across these hubs. If a task's deep material is **not** in this hub's Sub-skill
routing table, it is a reference file under a sibling hub below — **activate that hub or `Read` its
`references/<name>.md` directly**. Every former standalone skill in this family is now a reference under one
of these hubs (nothing was deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `integration-clients` | Integration & API Clients (Jira, Monday, Slack, Salesforce, Glean, Aha) | `references/aha-api.md`, `references/glean-dev.md`, `references/jira-developer-expert.md`, `references/jira-extension-client.md`, … |
