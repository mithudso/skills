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

Hub for building/maintaining **integration clients against third-party SaaS platforms** — Jira, Monday.com, Slack, Salesforce, Glean, Aha!. Covers vendor-specific API surfaces: REST/GraphQL, auth (OAuth 2.0, cookie/session, API tokens), webhook/event delivery, SDKs/CLIs, and Chrome MV3 extension quirks.

Vendor-detail layer. Assumes you decided *that* to integrate; answers *how* for specific platform.

## How to use this skill

Consolidates **8 integration sub-skills** as on-demand reference files under `references/`. Match task to routing table below. **Read listed `references/…md` file before answering deep questions** — table alone not enough for endpoint, auth-flow, payload-shape, or CLI detail. For exact field names, scopes, version specifics: vendor's official developer docs = source of truth.

## Sub-skill routing table

Hub absorbs 8 former standalone skills as on-demand reference files. Task matches row → **Read listed `references/` file** before answering — table alone insufficient for depth.

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

Hub owns vendor-specific integration-client dev. Hand off when task not about specific vendor's API:

- **General, vendor-neutral API design** (pagination strategy, idempotency keys, error-envelope conventions, REST-vs-RPC style, retry/backoff not tied to one platform) → `software-engineering-patterns`.
- **AI / LLM / agent integration** (tool-calling, agent orchestration, RAG over these data sources) → `ai-agent-engineering`.
- **Generic Chrome extension OAuth/identity plumbing** (`chrome.identity`, `launchWebAuthFlow`, token caching) not tied to specific vendor API → `chrome-extension-expert` (references/chrome-identity-oauth.md). Hub owns vendor-API call; extension identity flow lives there.
- **TAM operational workflows that *use* these platforms** — Monday board audits, Slack subscription audits, account/report generation — stay in standalone / tam-ops skills (`monday-board-audit`, `slack-subscription-auditor`, `tam-account-reports`). Hub builds client; those skills run recurring operational task on top.

Some topics straddle boundary (e.g. Monday GraphQL query inside board-audit workflow). Lead with hub matching intent: building/fixing *client/integration* → stay here; running *operational task* → route to tam-ops.

<!-- cross-hub-map -->
## Cross-hub map — where every integration topic lives

Family split across these hubs. Task's deep material **not** in this hub's sub-skill routing table → reference file under sibling hub below — **activate that hub or `Read` its `references/<name>.md` directly**. Every former standalone skill now reference under one hub (nothing deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `integration-clients` | Integration & API Clients (Jira, Monday, Slack, Salesforce, Glean, Aha) | `references/aha-api.md`, `references/glean-dev.md`, `references/jira-developer-expert.md`, `references/jira-extension-client.md`, … |