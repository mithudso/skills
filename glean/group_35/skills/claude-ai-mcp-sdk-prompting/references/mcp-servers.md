<!-- hub-reference-banner -->
> **Reference file — part of the `ai-mcp-sdk-prompting` hub.** Formerly the standalone `mcp-servers` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: mcp-servers
description: >
  MCP server development expert — JSON-RPC 2.0 architecture, Tools/Resources/Prompts primitives,
  capability negotiation, session lifecycle, building servers in Python (FastMCP 3.x, decorators,
  context injection) and TypeScript (@modelcontextprotocol/sdk, Zod schemas), stdio vs Streamable HTTP
  transport, OAuth 2.1 + PKCE security, sampling/elicitation, deployment (Docker, Cloud Run, ECS,
  Vercel), and the 20K+ server ecosystem.
  TRIGGER: building or debugging an MCP server, FastMCP decorators or tool definitions, MCP
  architecture or session lifecycle, MCP transport guidance, MCP OAuth 2.1 or authorization,
  deploying an MCP server, MCP sampling or elicitation, discovering existing MCP servers.
  SKIP: designing a new MCP server from scratch (use mcp-builder for the full greenfield workflow);
  general REST/gRPC API design with no MCP involvement; MCP client-side integration only.
version: 1.1.0
last_updated: 2026-05-29
origin: local
category: developer
tags:
  - mcp
  - model-context-protocol
  - fastmcp
  - typescript-sdk
  - json-rpc
  - stdio
  - streamable-http
  - oauth2-pkce
  - mcp-tools
  - mcp-resources
  - mcp-prompts
related_skills:
  - mcp-builder
  - typescript-expert
  - python-patterns
  - docker-containers
  - aws-serverless
---

# MCP Server Development

Comprehensive reference for building, deploying, and integrating MCP servers. Assumes working knowledge of Python or TypeScript. Backed by `references/mcp-servers-context.md`.

*Last researched: May 2026. 34 cited sources in `references/mcp-servers-context.md`.*

## When to use this skill

- Building an MCP server (Python or TypeScript)
- FastMCP decorators, tool definitions, or context injection
- MCP architecture, primitives, or session lifecycle questions
- Transport guidance (stdio vs Streamable HTTP)
- MCP security (OAuth 2.1, PKCE, authorization flows)
- Deploying an MCP server (Docker, Cloud Run, Vercel, ECS)
- Sampling (server requesting LLM completions) or elicitation (user input)
- Discovering or evaluating existing MCP servers
- `.mcp.json` configuration or `claude_desktop_config.json`
- MCP Gateway patterns for enterprise multi-server orchestration
- Health checks, telemetry, or observability for an MCP server

## When NOT to use this skill

- Full greenfield MCP server design → use `mcp-builder` (workflow, security checklist, evaluation)
- General REST API or gRPC design with no MCP involvement
- MCP client-only questions → point to the client SDK docs

## Quick-start decision tree

| I want to... | Start here |
|---|---|
| Build a local tool for Claude Desktop | Python: FastMCP + stdio. See "Building servers in Python." |
| Build a remote/cloud MCP server | TypeScript or Python + Streamable HTTP. See "Transport protocols." |
| Add auth to an existing server | See "Security: OAuth 2.1 + PKCE." |
| Deploy to production | See "Deployment patterns." |
| Find an existing server | See "Ecosystem and discovery." |
| Debug a broken server | See "Debugging and testing." |

---

## Architecture

MCP uses JSON-RPC 2.0 with a three-layer model: **Host** (AI application) → **Client** (stateful session manager) → **Server** (exposes tools, resources, prompts).

### Three server primitives

| Primitive | Direction | Purpose |
|-----------|-----------|---------|
| **Tools** | Client → Server | Executable functions the LLM can invoke |
| **Resources** | Client → Server | Read-only contextual data (URI-based) |
| **Prompts** | Client → Server | Reusable interaction templates |

### Two client primitives

| Primitive | Direction | Purpose |
|-----------|-----------|---------|
| **Sampling** | Server → Client | Request LLM completions (deprecated in DRAFT-2026) |
| **Elicitation** | Server → Client | Request structured user input |

### Session lifecycle

1. **Initialize** — client sends `initialize` with capabilities; server responds with its own
2. **Operation** — bidirectional JSON-RPC; client calls `tools/call`, `resources/read`, `prompts/get`
3. **Shutdown** — graceful close via `shutdown` request or transport disconnect

A server that does not declare `tools` capability will not receive `tools/list` requests. This keeps the protocol extensible without breaking older implementations.

---

## Building servers in Python (FastMCP 3.x)

FastMCP is the high-level wrapper inside the official MCP Python SDK. 4M+ daily PyPI downloads as of 2026.

### Core decorators

```python
from fastmcp import FastMCP
mcp = FastMCP("my-server")

@mcp.tool
def search(query: str) -> str:
    """Search the database."""
    return do_search(query)

@mcp.resource("config://app")
def get_config() -> str:
    """Application configuration."""
    return json.dumps(config)

@mcp.prompt
def review_prompt(code: str) -> str:
    """Code review template."""
    return f"Review this code:\n{code}"
```

Type hints are automatically converted to JSON Schema. FastMCP supports all Pydantic-compatible types.

### Context injection

```python
from fastmcp import Context

@mcp.tool
async def analyze(data: str, ctx: Context) -> str:
    """Analyze data with progress reporting."""
    await ctx.info("Starting analysis...")
    await ctx.report_progress(0, 100)
    result = process(data)
    schema = await ctx.read_resource("schema://main")
    await ctx.report_progress(100, 100)
    return f"Analysis complete: {result}"
```

Context capabilities: `ctx.debug/info/warning/error()`, `ctx.report_progress(current, total)`, `ctx.read_resource(uri)`. Each request gets a new Context instance — state does not persist between requests.

### Lifespan and dependency injection

```python
from contextlib import asynccontextmanager

@asynccontextmanager
async def lifespan(server):
    db = await connect_database()
    try:
        yield {"db": db}
    finally:
        await db.close()

mcp = FastMCP("my-server", lifespan=lifespan)

@mcp.tool
async def query_db(sql: str, ctx: Context) -> str:
    db = ctx.request_context.lifespan_context["db"]
    return await db.execute(sql)
```

### Error handling and FastMCP 3.x changes

- Raise `ToolError` for user-visible errors the LLM should see
- Set `mask_error_details=True` to hide stack traces in production
- 3.x: decorators return the original function (tools directly callable in tests)
- 3.x: `fastmcp dev` became `fastmcp dev inspector`
- **Code mode (v3.1):** collapses large tool catalogs into 2 meta-tools, cutting token usage by up to 99%

---

## Building servers in TypeScript

### Tool registration with Zod schemas

```typescript
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { z } from "zod";

const server = new McpServer({ name: "my-server", version: "1.0.0" });

server.registerTool("search", {
  description: "Search the database",
  inputSchema: z.object({
    query: z.string().describe("Search query"),
    limit: z.number().optional().default(10),
  }),
  annotations: { readOnlyHint: true },
}, async ({ query, limit }) => {
  const results = await db.search(query, limit);
  return { content: [{ type: "text", text: JSON.stringify(results) }] };
});
```

Zod schemas are one of the most important security layers — they keep untrusted LLM-generated arguments from reaching your logic in unexpected shapes.

### Transport setup

```typescript
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { StreamableHTTPServerTransport } from "@modelcontextprotocol/sdk/server/streamableHttp.js";

// stdio — local development
const transport = new StdioServerTransport();
await server.connect(transport);

// Streamable HTTP — remote/production
import express from "express";
const app = express();
app.use("/mcp", createMcpExpressApp(server));
app.listen(3000);
```

**Never log to stdout** — MCP uses stdout exclusively for JSON-RPC messages. Use `console.error()` or a file logger.

---

## Transport protocols

### stdio

- Communication via stdin/stdout of a child process
- Zero network overhead, inherently single-client
- Default for local developer tools and desktop integrations
- No authentication needed — the process boundary is the security boundary
- Configuration in `claude_desktop_config.json` with `command`, `args`, `env` fields

### Streamable HTTP

Introduced in protocol version 2025-03-26, replacing the deprecated SSE transport.

| Verb | Direction | Purpose |
|------|-----------|---------|
| **POST** | Client → Server | Send JSON-RPC requests |
| **GET** | Server → Client | Open SSE stream for notifications |
| **DELETE** | Client → Server | Terminate a session |

**Stateful mode:** server issues `Mcp-Session-Id`; client includes it on subsequent requests.
**Stateless mode:** each request is independent — simpler horizontal scaling, no persistent subscriptions.

### Configuration files

- `.mcp.json` — project-level, committable to version control
- `~/.claude.json` — user-level for personal servers
- CLI: `claude mcp add --scope local|project|user`

---

## Security: OAuth 2.1 + PKCE

Mandatory for public remote MCP servers since November 2025.

### Authorization flow

1. Client hits MCP endpoint, receives `401` with `WWW-Authenticate` header
2. Client fetches `/.well-known/oauth-authorization-server` metadata
3. PKCE authorization with **S256 method** (the only safe code exchange)
4. Token exchange with `resource` parameter (RFC 8707) binding token to one MCP server
5. Client sends Bearer token on subsequent requests

### Key requirements

- **Resource indicators (RFC 8707):** bind tokens to one server — blocks cross-server replay
- **Refresh rotation:** for public clients, rotate refresh tokens on each use
- **Step-up authorization:** start with minimal scopes, request more on `403` responses
- For stdio/local servers, API keys or bearer tokens are acceptable

### Known vulnerabilities

- **CVE-2025-6514** — consent-bypass; mitigate with exact redirect URI matching
- **SSRF risk** — 36.7% of servers vulnerable; protect metadata endpoints
- **Credential reuse** — 53% of deployments use static credentials; replace with OAuth
- Never forward client tokens to backend services

---

## Sampling and elicitation

**Sampling (deprecated DRAFT-2026):** server requests LLM completions via `sampling/createMessage`. 1-year support window. Must be issued during an originating client request.

**Elicitation:** server requests structured user input via `elicitation/create`.
- **Form mode** — JSON Schema for input fields (flat objects, primitives only). Never collect passwords, API keys, or payment data via form mode.
- **URL mode** — redirect to external URL for sensitive interactions (OAuth consent, payment flows).

---

## Deployment patterns

### Docker

```dockerfile
FROM node:22-slim
WORKDIR /app
COPY package*.json ./
RUN npm ci --omit=dev
COPY dist/ ./dist/
USER node
EXPOSE 3000
HEALTHCHECK --interval=30s CMD curl -sf http://localhost:3000/mcp || exit 1
CMD ["node", "dist/index.js"]
```

Run as non-root user. Bind to `0.0.0.0` (not `127.0.0.1`) for container networking.

### Google Cloud Run

Built-in HTTP streaming support for Streamable HTTP. Auto-scales to zero for intermittent usage.

### AWS ECS on Fargate

Ideal for long-lived servers needing warm caches and persistent connections. Reported costs under $3/month for low-traffic servers.

### Vercel

```typescript
import { createMcpHandler } from "mcp-handler";
export const POST = createMcpHandler(server);
```

Hobby plan: 10s timeout. Use Fluid compute for longer operations. Pre-warm with Cron Job hitting `/api/mcp/sse` every 5 minutes (cuts p95 latency from 2.8s to 400ms).

### Production health reality

April 2026 audit of 2,181 remote endpoints: 52% completely dead, only 9% fully healthy. Main causes: upstream API changes, cold start timeouts, schema drift. Fix: health checks, minimum instances, pinned dependencies.

---

## Ecosystem and discovery

### Scale (May 2026)

- 97M monthly SDK downloads across Python and TypeScript
- 10,000–20,000+ public servers depending on registry
- Governed by the Agentic AI Foundation (AAIF) under Linux Foundation since December 2025

### Major registries

| Registry | Servers | Strengths |
|----------|---------|-----------|
| **Glama** (glama.ai) | 21,000-24,000 | Largest volume, daily updates, metaregistry |
| **mcp.so** | 19,700+ | Community coverage, unofficial tools |
| **Smithery** (smithery.ai) | 7,000+ | Hosted remote servers, CLI install |
| **awesome-mcp-servers** | Curated | GitHub-maintained quality list |

`npx mcpfinder` discovers and evaluates 25K+ servers across registries.

---

## Debugging and testing

### MCP Inspector

```bash
npx @modelcontextprotocol/inspector
```

Connects to your server, lists tools/resources/prompts, and lets you invoke them interactively.

### Testing FastMCP (Python)

```python
# Direct function call — no server needed (3.x behavior)
result = search("test query")
assert "expected" in result

# Full protocol test via Client
from fastmcp import Client
async with Client(mcp) as client:
    result = await client.call_tool("search", {"query": "test"})
    assert result[0].text == "expected"
```

### Testing TypeScript

```typescript
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { InMemoryTransport } from "@modelcontextprotocol/sdk/inMemory.js";

const [clientTransport, serverTransport] = InMemoryTransport.createLinkedPair();
await server.connect(serverTransport);
const client = new Client({ name: "test" });
await client.connect(clientTransport);
const result = await client.callTool({ name: "search", arguments: { query: "test" } });
```

---

## Best practices checklist

- [ ] Zod/Pydantic input validation — never trust raw LLM arguments
- [ ] Health checks (`/healthz` or `/mcp`) for production deployments
- [ ] `ToolError` (Python) or structured error responses (TypeScript) for user-visible errors
- [ ] Hide stack traces in production (`mask_error_details=True` in FastMCP)
- [ ] OAuth 2.1 + PKCE for any public-facing remote server
- [ ] Resource indicators (RFC 8707) to bind tokens to specific servers
- [ ] Progress reporting for tools that take more than a few seconds
- [ ] Lifespan/dependency injection for shared resources (DB connections, API clients)
- [ ] Pin dependencies and monitor for schema drift
- [ ] Test with MCP Inspector during development
- [ ] Tool annotations (`readOnlyHint`, `destructiveHint`) set correctly
- [ ] Streamable HTTP (not deprecated SSE) for all new remote servers
- [ ] Structured logging and telemetry for production observability

---

## Common mistakes

| Mistake | Why it fails | Fix |
|---------|-------------|-----|
| Deprecated SSE transport for new servers | Replaced by Streamable HTTP in March 2025 | Use Streamable HTTP |
| Skipping input validation | LLM arguments can be malformed or injected | Always validate with Zod/Pydantic |
| Running Docker containers as root | Container escape gives attacker root on host | `USER node` or `USER nobody` |
| Hardcoding secrets | Credentials leak into version control | Use environment variables or secret managers |
| Binding to `127.0.0.1` in containers | Other containers cannot reach the server | Bind to `0.0.0.0`; restrict via network policy |
| No health checks | Dead servers stay in rotation | Add `/mcp` or `/healthz` with liveness probes |
| Static API keys instead of OAuth | Tokens cannot be scoped, rotated, or audited | OAuth 2.1 + PKCE for production |
| Raw stack traces to LLM | Leaks paths, versions, DB schemas | `mask_error_details=True` or catch and wrap |
| Unpinned dependencies | Upstream changes break servers silently | Pin versions; monitor for schema drift |

---

## Sources

34 cited sources in `references/mcp-servers-context.md`. Key official references:
- [MCP Architecture](https://modelcontextprotocol.io/docs/learn/architecture) | [Authorization](https://modelcontextprotocol.io/docs/tutorials/security/authorization)
- [Python SDK (FastMCP)](https://github.com/modelcontextprotocol/python-sdk) | [TypeScript SDK](https://github.com/modelcontextprotocol/typescript-sdk)
- [FastMCP docs](https://gofastmcp.com/servers/context) | [TS SDK docs](https://ts.sdk.modelcontextprotocol.io/)
