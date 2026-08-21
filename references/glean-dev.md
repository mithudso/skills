<!-- hub-reference-banner -->
> **Reference file — part of the `integration-clients` hub.** Formerly the standalone `glean-dev` skill.
> Sibling topics in this family are now reference files under the hubs (`integration-clients`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: glean-dev
version: 1.1.0
updated: 2026-05-29
description: >
  Glean enterprise developer reference — CLI, MCP, TypeScript/Python SDK clients,
  Client API, Indexing API, Agents API, Answers API, Webhooks, Admin APIs,
  permissions, and governance. Also covers Mozilla Glean telemetry SDK (a separate
  product) with clear separation.
  TRIGGER: user asks about Glean CLI, Glean API, Glean MCP, Glean SDK, indexing
  content into Glean, Glean search/chat integration, Glean agents, Glean answers,
  Glean permissions, or enterprise knowledge search. Also triggers for Mozilla
  Glean telemetry instrumentation.
  SKIP: questions about general enterprise search products (not Glean), Glean
  connector setup via Admin Console UI only (use glean-llm-client-patterns for
  code integration patterns), or Glean Analytics dashboards not involving the API.
category: developer
tags: [glean, enterprise-search, mcp, api-client, typescript, python, cli, indexing, agents]
related_skills: [glean-llm-client-patterns, mcp-servers, chrome-native-messaging]
whenToUse:
  - "How do I install the Glean CLI?"
  - "What Glean CLI commands are available?"
  - "How do I authenticate with the Glean API?"
  - "Glean TypeScript SDK setup"
  - "How do I index documents into Glean?"
  - "What is the difference between Glean Client API and Indexing API?"
  - "How do I run a Glean agent via the API?"
  - "Set up Glean MCP in Claude Code"
  - "Mozilla Glean telemetry SDK setup"
  - "Glean permissions model for custom datasources"
  - "Glean governance API"
  - "What are the Glean admin API scopes?"
whenNotToUse:
  - "Use glean-llm-client-patterns for Chrome extension Glean integration or auth fallback cascade patterns"
  - "Use mcp-servers for general MCP protocol details not specific to Glean"
origin: local
---

# Glean Developer Reference

Reference for building with Glean enterprise (search, chat, agents, MCP, indexing)
and Mozilla Glean (telemetry). These are **two different products** that share a
name — see Section 1 for the naming distinction.

## When to use this skill

- Working with the Glean enterprise CLI, SDK, or APIs
- Building Glean MCP integrations in IDE or agent hosts
- Indexing content into Glean via the Indexing API
- Administering Glean tokens, permissions, or governance
- Working with Mozilla Glean telemetry instrumentation

## When NOT to use this skill

- For Chrome extension Glean integration patterns (auth cascade, tab-based session,
  MCP JSON-RPC flow): use `glean-llm-client-patterns` instead.
- For general MCP protocol (not Glean-specific): use `mcp-servers`.

---

## 1. Critical naming distinction

Two different products share the "Glean" name:

| Product | Purpose | Main domains |
|---------|---------|-------------|
| **Glean enterprise** | Enterprise search, chat, agents, MCP, indexing, admin | `docs.glean.com`, `developers.glean.com`, `github.com/gleanwork/*` |
| **Mozilla Glean** | Cross-platform telemetry SDK and parser tooling | `mozilla.github.io/glean` |

Do **not** conflate them. Sections 2–12 cover Glean enterprise. Section 13 covers
Mozilla Glean.

---

## 2. Glean enterprise: capability map

| Surface | Use for | Main packages / tools |
|---------|--------|----------------------|
| **CLI** | Terminal search, chat, agents, knowledge management | `glean` CLI (`gleanwork/glean-cli`) |
| **MCP** | AI host integration (Claude, Cursor, VS Code, Copilot, ChatGPT) | Glean MCP server |
| **TypeScript client** | User-facing search/chat/agents in Node or web | `@gleanwork/api-client` v0.14.19+ |
| **Python client** | User-facing search/chat/agents in Python | `glean-api-client` v0.12.24+ |
| **Client API** | search, chat, agents, collections, user-facing workflows | `https://<server>/rest/api/v1/` |
| **Indexing API** | custom datasources, content ingestion, permissions | `https://<server>/api/index/v1/` |

---

## 3. Glean enterprise CLI

Repo: `gleanwork/glean-cli` — company knowledge, search, and AI from the terminal.

### 3.1 Installation

```bash
brew install gleanwork/tap/glean-cli
# or
curl -fsSL https://raw.githubusercontent.com/gleanwork/glean-cli/main/install.sh | sh
```

### 3.2 Quick start

```bash
glean auth login
glean search "vacation policy"
glean chat "Summarize our Q1 engineering goals"
```

### 3.3 Authentication model

```bash
glean auth login    # interactive login
glean auth status   # check current auth state
glean auth logout   # sign out
```

For CI/CD (credential resolution order: env vars → system keyring → `~/.glean/config.json`):

```bash
export GLEAN_API_TOKEN=your-token
export GLEAN_SERVER_URL=https://your-server-url
```

### 3.4 Agent-friendly design

- Structured JSON output on stdout; errors on stderr
- `--dry-run` to preview requests before executing write/delete operations
- `glean schema <command>` for machine-readable command schemas
- `--output ndjson` for streaming large result sets
- `--fields` to request only needed response fields

### 3.5 Core commands

Top-level: `search`, `chat`, `api`, `schema`, `auth`

Namespaces: `agents`, `answers`, `announcements`, `collections`, `documents`,
`entities`, `insights`, `messages`, `pins`, `shortcuts`, `tools`, `verification`,
`activity`

### 3.6 Best practices for LLM use of the CLI

1. Start with `glean schema` when unsure about a command.
2. Use `--dry-run` before write/delete operations.
3. Use `--fields` when you only need a subset of response fields.
4. Use `--output ndjson` for streaming or large-result workflows.
5. Keep secrets in env vars — never inline in committed scripts.

---

## 4. Glean enterprise MCP

### 4.1 Public developer-docs MCP

Endpoint for AI hosts to query Glean developer documentation:

**Claude Code:**
```bash
claude mcp add glean-developer-docs https://developers.glean.com/mcp --transport http --scope user
```

**Cursor:**
```json
{
  "mcpServers": {
    "glean-developer-docs": { "type": "http", "url": "https://developers.glean.com/mcp" }
  }
}
```

**VS Code:**
```bash
code --add-mcp '{"name":"glean-developer-docs","type":"http","url":"https://developers.glean.com/mcp"}'
```

### 4.2 Enterprise knowledge MCP

The Glean enterprise MCP bridges AI hosts to company knowledge with:
- Permission-aware access (queries return only documents the user can access)
- MCP host/server architecture with support for Claude, Cursor, VS Code, Copilot, ChatGPT, Windsurf, Goose
- Usage-based pricing; security inherited from Glean session and permission models

Effective prompt patterns for Glean MCP:
- State the data or action desired explicitly
- Provide document links or IDs when available
- Ask the assistant which tool it plans to use when debugging tool selection

### 4.3 Runtime Glean MCP tools (this environment)

This environment exposes three Glean knowledge tools:

| Need | Best tool |
|------|-----------|
| Find candidate internal docs fast | `glean_default-search` |
| Read a known internal URL exactly | `glean_default-read_document` |
| Synthesize across internal knowledge | `glean_default-chat` |

**`glean_default-search`** — targeted retrieval across company knowledge.
- Use **short, discriminative keywords** — avoid full sentences
- Avoid boolean logic in the main query
- Only use time filters if the user explicitly specifies a time range
- Use discovered dynamic filters from results instead of inventing them

**`glean_default-read_document`** — full content of known URLs.
- Accepts multiple URLs in one call
- Supports `raw_bytes` mode for original file bytes
- Use when you know the exact document URL

**`glean_default-chat`** — synthesis and analysis across enterprise context.
- Use for complex questions, unknown internal issues, contextual analysis
- Avoid for simple retrieval where `search` or `read_document` suffice

Recommended workflow: search → inspect results → read exact documents → synthesize only if needed.

---

## 5. Authentication essentials

### 5.1 Finding your server URL

1. Go to `https://app.glean.com/admin/about-glean`
2. Find **Server instance (QE)**
3. Strip the trailing slash

Example: displayed `https://acme-prod-be.glean.com/` → use `https://acme-prod-be.glean.com`

### 5.2 Token types

| Token type | Scope | Use `X-Glean-ActAs`? |
|-----------|-------|----------------------|
| User-scoped Client API | Single user's data and scopes | Must omit or leave empty |
| Global Client API | Any user (impersonation) | **Required** — include target user email |
| Indexing API | Full indexing access | N/A |

- User-scoped tokens: default safe choice
- Global tokens: Super Admin only; require `X-Glean-ActAs` for every request
- Scopes cannot be modified after token creation

See Section 11.1 for the full admin token type table (who can create each type and what capabilities each grants).

### 5.3 API base URLs

| API | Base URL |
|-----|---------|
| Client API | `https://<server>/rest/api/v1/` |
| Indexing API | `https://<server>/api/index/v1/` |

Auth choices: **OAuth** (recommended for Client API, user-facing flows) or **Glean-issued tokens**.

---

## 6. TypeScript client

```bash
npm install @gleanwork/api-client
```

**Bootstrap:**
```ts
import { Glean } from "@gleanwork/api-client";

const client = new Glean({
  apiToken: process.env.GLEAN_API_TOKEN,
  serverURL: process.env.GLEAN_SERVER_URL,
});
```

**Chat:**
```ts
const result = await client.client.chat.create({
  messages: [{ fragments: [{ text: "What are our company values?" }] }],
});
```

**Streaming chat:**
```ts
const stream = client.client.chat.stream({
  messages: [{ fragments: [{ text: "What are our priorities?" }] }],
});
for await (const chunk of stream) { console.log(chunk.text); }
```

**Search:**
```ts
const results = await client.client.search.search({
  query: "quarterly business review",
  pageSize: 10,
});
results.results?.forEach(r => console.log(`${r.title} — ${r.url}`));
```

**Global token + ActAs:**
```ts
const response = await client.client.chat.create(
  { messages: [{ fragments: [{ text: "Hello" }] }] },
  { headers: { "X-Glean-ActAs": "user@company.com" } }
);
```

---

## 7. Python client

```bash
pip install glean-api-client
```

**Bootstrap:**
```py
from glean.api_client import Glean
import os

client = Glean(
    api_token=os.getenv("GLEAN_API_TOKEN"),
    server_url=os.getenv("GLEAN_SERVER_URL"),
)
```

**Chat:**
```py
response = client.client.chat.create(
    messages=[{"fragments": [{"text": "What are our company values?"}]}],
    timeout_millis=30000,
)
```

**Search:**
```py
results = client.client.search.search(query="quarterly business review", page_size=10)
for result in results.results:
    print(f"{result.title} — {result.url}")
```

**Agents:**
```py
response = client.client.agents.create_and_wait_run(
    agent_id="your-agent-id",
    inputs={"query": "Analyze sales performance"},
)
```

---

## 8. Agent-building approaches

| Approach | Best for |
|---------|---------|
| Direct API | Maximum control, custom apps |
| MCP | IDE/host integration with low setup cost |
| LangChain | Python ecosystems |
| Glean Agent Toolkit | Cross-framework tooling |

Required scopes for agent execution: `agents` and `chat` (plus `search` depending on agent configuration).

### 8.1 Agents API (Beta — LangChain Agent Protocol subset)

**REST endpoints:**

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/rest/api/v1/agents/{agent_id}` | Retrieve agent details |
| GET | `/rest/api/v1/agents/{agent_id}/schemas` | Get input/output schemas |
| POST | `/rest/api/v1/agents/search` | Search agents by name |
| POST | `/rest/api/v1/agents/runs/wait` | Execute agent, return final response (blocking) |
| POST | `/rest/api/v1/agents/runs/stream` | Execute agent, return SSE stream |
| POST | `/rest/api/v1/agents/{agent_id}` | Create draft or publish agent |

Authentication: user-scoped Client API token with `AGENTS` scope.

**Run request schema:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `agent_id` | string | Yes | ID of the agent to run |
| `input` | Dict[str, Any] | Conditional | All form fields (required for form-triggered agents) |
| `messages` | List[Message] | No | Messages with `role` field (e.g., `"USER"`) |
| `metadata` | Dict[str, Any] | No | Metadata passed to the agent |

**Agent discovery:**
- Search by name: `POST /rest/api/v1/agents/search` with `{ "name": "HR Policy Agent" }`
- Agent IDs appear in Agent Builder UI URL: `/admin/agents/{agentId}`
- Get schemas: `GET /rest/api/v1/agents/{id}/schemas`

**TypeScript:**
```ts
const agent = await glean.client.agents.retrieve("<agent_id>");
const result = await glean.client.agents.run({
  agentId: "<id>",
  messages: [{ role: "USER" }],
});
```

**Python:**
```py
res = glean.client.agents.run(agent_id="<id>", messages=[{"role": "USER"}])
res = glean.client.agents.run_stream(agent_id="<id>", messages=[{"role": "USER"}])
```

**Best practices:**
- Retry with exponential backoff for HTTP 429
- Cache: MD5 hash of query+filters as key, 30-min TTL default
- Rate limiting: enforce requests-per-minute (default 60)
- Multi-turn: maintain last 10 messages for conversation history

---

## 9. Answers API

All endpoints use POST.

| Method | Path | SDK Method | Purpose |
|--------|------|------------|---------|
| POST | `/rest/api/v1/createanswer` | `answers.create()` | Create a Q&A entry |
| POST | `/rest/api/v1/editanswer` | `answers.update()` | Update an existing answer |
| POST | `/rest/api/v1/deleteanswer` | `answers.delete()` | Delete an answer |
| POST | `/rest/api/v1/getanswer` | `answers.retrieve()` | Get a specific answer by ID |
| POST | `/rest/api/v1/listanswers` | `answers.list()` | List answers (**deprecated**) |

Authentication: Client API token with `ANSWERS` scope.

**Answer object required fields:** `question` (string), `bodyText` (string).
**Optional:** `audienceFilters`, `addedRoles`, `roles`. Roles: `VIEWER`, `EDITOR`, `OWNER`, `ANSWER_MODERATOR`.
**Delete/retrieve identifier:** `id` (int) or `docId` (e.g., `"ANSWERS_answer_3"`).

**TypeScript example:**
```ts
const result = await glean.client.answers.create({
  data: {
    question: "Why is the sky blue?",
    bodyText: "Blue light is more strongly scattered...",
    addedRoles: [{
      person: { name: "Alice", obfuscatedId: "abc123" },
      role: "VIEWER",
    }],
  },
});
const answer = await glean.client.answers.retrieve({ id: 3 });
```

---

## 10. Webhooks / Events

**Glean does NOT expose outbound webhook subscriptions** (as of May 2026).

What exists:
- **Inbound webhook receivers:** Glean consumes webhooks FROM external sources (HubSpot, Greenhouse, SharePoint) to keep its index current.
- **Activity API (inbound):** `POST /rest/api/v1/activity` accepts events pushed TO Glean (requires ACTIVITY scope + global token).

Workaround patterns when you need Glean-originated event notifications:
1. **Polling:** use the Insights API periodically
2. **Agent-based:** run agents on a schedule and process outputs
3. **CLI scripting:** use `glean` CLI in cron jobs

---

## 11. Admin Tokens and Governance

### 11.1 Token types

| Token type | Created by | Capabilities |
|-----------|-----------|-------------|
| User-scoped (Client API) | Super Admin (any user); Admin/API Token Creator (self only) | Limited to specific user's data and scopes |
| Global (Client API) | Super Admin only | Impersonates any user via `X-Glean-ActAs` |
| Indexing API | Super Admin (any); API Token Creator (self only) | Full indexing access; optional datasource restriction |

### 11.2 Client API scopes

**Standard (17):** `ACTIVITY`, `AGENTS`, `ANNOUNCEMENTS`, `ANSWERS`, `CHAT`, `COLLECTIONS`, `DOCPERMISSIONS`, `DOCUMENTS`, `ENTITIES`, `FEEDBACK`, `INSIGHTS`, `PEOPLE`, `PINS`, `SEARCH`, `SHORTCUTS`, `SUMMARIZE`, `VERIFICATION`

**Admin-only (2):** `DATA_GOVERNANCE` (Super Admin or Sensitive Content Moderator), `CONTENT_HIDING` (Super Admin with visibility override)

Scopes cannot be modified after token creation.

### 11.3 Governance API

Requires `DATA_GOVERNANCE` or `CONTENT_HIDING` scope.

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/rest/api/v1/governance/data/policies` | List policies |
| POST | `/rest/api/v1/governance/data/policies` | Create policy |
| POST | `/rest/api/v1/governance/data/reports` | Create compliance report |
| GET | `/rest/api/v1/governance/data/reports/{id}/status` | Check report status |
| GET | `/rest/api/v1/governance/data/reports/{id}/download` | Download report CSV |
| POST | `/rest/api/v1/governance/documents/visibilityoverrides` | Set document hiding/unhiding |

### 11.4 Token rotation (Indexing API)

`POST /api/index/v1/rotatetoken` — minimum 1440-minute rotation period. Supports IP address restrictions via CIDR notation.

---

## 12. Permissions model

### 12.1 Permission modes per document

| Mode | Field | Behavior |
|------|-------|---------|
| Anonymous | `allowAnonymousAccess: true` | Any Glean user can find it |
| User-specific | `allowedUsers: [...]` | Only listed users |
| Datasource-wide | `allowAllDatasourceUsersAccess: true` | All indexed users in datasource |
| Group-based | `allowedGroups: [...]` | Members of listed groups |

### 12.2 Document permissions schema (Indexing API)

```json
{
  "permissions": {
    "allowAnonymousAccess": false,
    "allowedUsers": [{ "email": "user@company.com" }],
    "allowAllDatasourceUsersAccess": false,
    "allowedGroups": ["engineering-team"]
  }
}
```

### 12.3 Group management endpoints

| Path | Purpose |
|------|---------|
| `/api/index/v1/indexgroup` | Add/update a single group |
| `/api/index/v1/bulkindexgroups` | Replace all groups in a datasource |
| `/api/index/v1/deletegroup` | Delete a group |
| `/api/index/v1/indexmembership` | Add user or nested group to a group |
| `/api/index/v1/checkdocumentaccess` | Verify a user's access to a document |
| `/api/index/v1/updatedocumentpermissions` | Update permissions on existing documents |

Group naming: no whitespace, no "scio" prefix, non-empty.

### 12.4 ACL patterns for custom datasources

```json
{ "permissions": { "allowAnonymousAccess": true } }
{ "permissions": { "allowedGroups": ["engineering"] } }
{ "permissions": { "allowedUsers": [{"email": "vp@company.com"}], "allowedGroups": ["leadership-team"] } }
```

### 12.5 Permissions flow

1. **Indexing time:** documents pushed with permissions attached
2. **User/group resolution:** Glean resolves group memberships at query time
3. **Search time:** results filtered to documents the querying user can access
4. **Source ACL preservation:** Glean honors source system ACLs exactly
5. **IdP:** reads groups/members from IdPs but does NOT enforce IdP permission rules on document access

**Processing notes:**
- Permissions process asynchronously — visibility changes may be delayed
- Use `checkdocumentaccess` to verify access during development
- Users and groups must be indexed BEFORE being referenced in document permissions
- `isUserReferencedByEmail` on the datasource config controls the identity mechanism

---

## 13. Mozilla Glean telemetry SDK

Mozilla Glean is a **separate telemetry product** — use only for telemetry, metrics instrumentation, parser/metrics generation, and Mozilla-style analytics pipelines. Do **not** use it for enterprise search or company knowledge.

```bash
npm install @mozilla/glean
```

Entry points: `@mozilla/glean/web`, `@mozilla/glean/webext`, `@mozilla/glean/node`

**CLI (wraps `glean_parser` in a Python virtualenv):**
```bash
npx glean --help
VIRTUAL_ENV="my/venv/path" npx glean --help  # custom virtualenv
```

Version pinning:
```bash
pip install -U glean_parser==$(npx glean --glean-parser-version)
```

Current npm version: `@mozilla/glean` **5.0.8**

---

## 14. Decision guide

| Task | Recommended path |
|------|-----------------|
| Internal company knowledge retrieval | Glean MCP (`glean_default-search` / `glean_default-chat`) |
| IDE / MCP integration | Glean developer-docs MCP for docs; enterprise MCP for company knowledge |
| TypeScript app code integration | `@gleanwork/api-client` |
| Python app code integration | `glean-api-client` |
| Index custom content into Glean | Indexing API (`/api/index/v1/`) |
| Telemetry instrumentation | Mozilla Glean (`@mozilla/glean`) |
| "Glean CLI" disambiguation | Enterprise CLI = company knowledge/search/AI; Mozilla CLI = `glean_parser` wrapper |

---

## 15. Sources

1. Mozilla Glean JS SDK: https://mozilla.github.io/glean/book/language-bindings/javascript/index.html
2. Glean CLI GitHub: https://github.com/gleanwork/glean-cli
3. Glean docs: https://docs.glean.com/
4. Glean MCP platform docs: https://docs.glean.com/administration/platform/mcp/about
5. Glean developers: https://developers.glean.com/
6. Glean MCP guide: https://developers.glean.com/guides/mcp
7. Glean agent approaches: https://developers.glean.com/guides/agents/overview
8. Glean authentication: https://developers.glean.com/get-started/authentication
9. Glean API clients: https://developers.glean.com/libraries/api-clients
10. Glean TypeScript client: https://developers.glean.com/libraries/api-clients/typescript
11. Glean Python client: https://developers.glean.com/libraries/api-clients/python
12. Agents API: https://developers.glean.com/api/client-api/agents/overview
13. Answers API: https://developers.glean.com/api/client-api/answers/overview
14. Governance API: https://developers.glean.com/api/client-api/governance/overview
15. Permissions guide: https://developers.glean.com/api-info/indexing/documents/permissions
16. Authentication overview: https://developers.glean.com/get-started/authentication
17. API Tokens (Admin Console): https://docs.glean.com/administration/developer/api-tokens
