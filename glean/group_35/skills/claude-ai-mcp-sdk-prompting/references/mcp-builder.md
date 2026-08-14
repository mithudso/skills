<!-- hub-reference-banner -->
> **Reference file — part of the `ai-mcp-sdk-prompting` hub.** Formerly the standalone `mcp-builder` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: mcp-builder
description: >
  Guide for creating high-quality MCP (Model Context Protocol) servers that enable LLMs
  to interact with external services through well-designed tools. Covers server architecture,
  tool design, security, testing, deployment, and observability for Python (FastMCP) and
  TypeScript (MCP SDK).
  TRIGGER: building an MCP server, designing MCP tools, adding OAuth to an MCP server,
  testing with MCP Inspector, deploying an MCP server to Docker/Cloud Run/Vercel/Lambda,
  adding observability or health checks to an MCP server, gateway patterns for 3+ servers.
  SKIP: integrating an existing MCP server into a client application (use mcp-servers skill);
  general REST/GraphQL API design with no MCP involvement; extending an existing MCP codebase
  (read that codebase's conventions first, then use this as a reference).
version: 2.1.0
last_updated: 2026-05-29
category: developer
tags:
  - mcp
  - model-context-protocol
  - mcp-server
  - fastmcp
  - typescript-sdk
  - tool-design
  - oauth2
  - opentelemetry
  - mcp-inspector
  - mcp-gateway
triggers:
  - building an MCP server
  - MCP tool design
  - FastMCP Python server
  - TypeScript MCP SDK
  - MCP OAuth 2.1 PKCE
  - MCP Inspector testing
  - deploy MCP server Docker
  - MCP gateway proxy
  - MCP observability OpenTelemetry
  - MCP security checklist
  - MCP cold start optimization
  - MCP server discovery well-known
  - MCP sampling elicitation
  - MCP resources prompts primitives
related_skills:
  - mcp-servers
  - typescript-expert
  - python-patterns
  - docker-containers
  - aws-serverless
---

# MCP Server Development Guide

Create MCP (Model Context Protocol) servers that enable LLMs to interact with external services through well-designed tools. A high-quality MCP server is measured by how well it enables LLMs to accomplish real-world tasks.

## When NOT to use this guide

- **MCP client integration** — covers server development only. For client-side patterns, use the `mcp-servers` skill.
- **Extending an existing MCP server** — read that codebase's conventions first; this guide targets greenfield servers.
- **Non-MCP API development** — if the consumer is a REST/GraphQL client, not an LLM agent, standard API design guides apply.

---

## High-level workflow (4 phases)

### Phase 1 — Research and planning

**API coverage vs. workflow tools:** balance comprehensive endpoint coverage with specialized composite tools. When uncertain, prioritize comprehensive API coverage. Agents can compose basic tools; composites help with common workflows.

**Tool naming:** use consistent prefixes with action-oriented names — `github_create_issue`, `github_list_repos`. Clear names help agents find the right tool without guessing.

**Context management:** return focused, relevant data. Design tools that filter and paginate rather than dumping everything.

**Actionable errors:** every error message must guide agents toward a solution with specific next steps.

Load framework documentation during planning:
- **MCP Protocol:** start at `https://modelcontextprotocol.io/sitemap.xml`; fetch pages with `.md` suffix
- **TypeScript SDK:** `https://raw.githubusercontent.com/modelcontextprotocol/typescript-sdk/main/README.md`
- **Python SDK:** `https://raw.githubusercontent.com/modelcontextprotocol/python-sdk/main/README.md`
- See also: [MCP Best Practices](./reference/mcp_best_practices.md)

**Recommended stack:**
- **Language:** TypeScript — high-quality SDK support, static typing, broad LLM training data coverage, compatible with MCPB
- **Transport:** Streamable HTTP (stateless JSON) for remote servers; stdio for local servers

### Phase 2 — Implementation

Set up project structure per language guide, then implement:

1. Shared API client with authentication
2. Error handling helpers
3. Response formatting (JSON/Markdown)
4. Pagination support
5. Tools with Zod (TypeScript) or Pydantic (Python) input schemas

For each tool, set annotations: `readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`.

**Full implementation patterns:**
- [TypeScript guide](./reference/node_mcp_server.md) — project structure, Zod schemas, tool registration, working examples
- [Python/FastMCP guide](./reference/python_mcp_server.md) — decorators, context injection, lifespan, working examples

### Phase 3 — Review and test

Build and validate before shipping:

```bash
# TypeScript
npm run build
npx @modelcontextprotocol/inspector          # interactive UI

# Python
python -m py_compile your_server.py
npx @modelcontextprotocol/inspector          # works for both
```

**The five testing gates** (must pass before production):

| Gate | What it tests | Tool |
|------|--------------|------|
| 1. Schema validation | Tool schemas are valid JSON Schema | MCP Inspector `tools/list` |
| 2. Smoke tests | Each tool executes on happy-path inputs | Inspector CLI |
| 3. Error-path tests | Tools return structured errors for bad inputs | Custom harness |
| 4. Integration tests | Tools interact correctly with real/mock backends | Vitest/pytest |
| 5. Agent evaluation | LLM can solve realistic tasks using the tools | Evaluation XML |

**CI smoke test:**
```bash
npx @modelcontextprotocol/inspector --cli \
  --server "node dist/index.js" \
  tools/list | jq '.tools | length'
```

**Manifest snapshot testing** — capture the tool manifest and diff it in CI:
```bash
npx @modelcontextprotocol/inspector --cli \
  --server "node dist/index.js" \
  tools/list > tests/manifests/tools.json
git diff --exit-code tests/manifests/tools.json
```

See [Evaluation Guide](./reference/evaluation.md) for creating 10-question evaluation sets.

### Phase 4 — Evaluations

Create 10 complex, realistic questions to test whether an LLM can use your server effectively. Each question must be independent, read-only, verifiable, and stable over time. Full guidance: [Evaluation Guide](./reference/evaluation.md).

---

## Security checklist (non-negotiable)

A 2026 scan of 1,808 MCP servers found 66% had at least one security finding: 43% shell/command injection, 20% tooling infrastructure exploits, 13% auth bypasses, 10% path traversal.

| Gate | Check | Why |
|------|-------|-----|
| Input validation | Every parameter validated against schema constraints | Prevents injection and malformed requests |
| Enum over string | Known value sets use enums, not free strings | Cheapest attack-surface reduction |
| URL/path allow-listing | URLs to known hosts only | Prevents SSRF and open-redirect |
| Sanitize metadata | Validate every field from discovery responses | Blocks tool-poisoning via hidden instructions |
| No internal error leaks | Never expose stack traces to the model | Hides implementation details |
| Audit surface | Destructive tools must have an audit trail | Biggest single production incident risk |
| Trust boundaries | Document what the server can and cannot access | Prevents privilege escalation |

**Critical CVE:** CVE-2025-6514 (CVSS 9.6) — command injection in `mcp-remote` via unsanitized OAuth URLs. Always sanitize external URLs before passing to system calls.

### OAuth 2.1 + PKCE (mandatory for public remote servers)

- **PKCE is mandatory** (not optional) for all OAuth flows; use S256 code challenge method
- **Implicit grant is prohibited** — OAuth 2.1 removes it entirely
- For stdio/local servers, API keys or bearer tokens are acceptable
- **Resource indicators (RFC 8707):** bind each token to one MCP server to block cross-server replay

---

## Tool design checklist

- [ ] Name is action-oriented (`create_issue`, not `issue` or `issueManager`)
- [ ] Description is one sentence: when to use this tool
- [ ] Every parameter has a description with examples where helpful
- [ ] Enums used for known value sets — never free strings for finite options
- [ ] IDs are branded/patterned with regex constraints
- [ ] Pagination supported for any list operation (cursor-based preferred)
- [ ] Error responses are actionable — specific guidance, not generic messages
- [ ] Annotations set correctly (`readOnlyHint`, `destructiveHint`, `idempotentHint`)
- [ ] Destructive tools have confirmation or audit trail
- [ ] Response size is bounded — paginate or truncate large results
- [ ] Tool count is reasonable — fewer than 25 tools; prefer composite tools for common workflows
- [ ] No stdout logging (TypeScript) — use `console.error()` or file logger
- [ ] Schema snapshot exists in `tests/` for CI diffing

**Tool count sweet spot:** GitHub Copilot cut from 40 to 13 tools with measurable benchmark improvements. Block rebuilt its Linear MCP server from 30+ tools down to 2. Fewer tools, better descriptions, outcome-oriented design.

---

## The seven anti-patterns

| # | Anti-pattern | Problem | Fix |
|---|-------------|---------|-----|
| 1 | Under-specified schemas | Model cannot generalize from vague descriptions | Add constraints, examples, and enum values to every parameter |
| 2 | Retrofitted auth | Security gaps expensive to fix post-launch | Design auth from day one; OAuth 2.1 + PKCE for remote servers |
| 3 | Chatty request patterns | Accumulated latency from many small calls | Provide batch/composite tools for common multi-step workflows |
| 4 | God-tools | One tool that does everything fails in too many ways | Split into focused, single-purpose tools |
| 5 | Omnibus parameter blobs | Too many optional parameters confuse the model | Group related parameters; use separate tools for distinct operations |
| 6 | Indiscriminate error envelopes | Generic errors prevent agent self-correction | Return specific, actionable errors per failure mode |
| 7 | Missing audit gates | Destructive operations with no trail | Log every mutation; require confirmation for destructive actions |

**Three pre-launch blockers** (do not publish with any of these):
1. Destructive tools with no audit surface
2. No clear trust boundary documentation
3. A single god-tool that tries to do everything

---

## Error handling pattern (TypeScript)

```typescript
catch (error) {
  if (error instanceof AuthenticationError) {
    return {
      content: [{ type: "text", text: "Authentication failed. Verify credentials." }],
      isError: true,
    };
  }
  console.error("Unexpected error in tool_name:", error);
  return {
    content: [{ type: "text", text: "Internal error occurred. Retry or contact support." }],
    isError: true,
  };
}
```

**Error handling rules:**
1. Validate early — return structured errors before any processing
2. Catch specific exceptions first — targeted messages help the agent self-correct
3. Always return a valid `CallToolResult` — never break the JSON-RPC connection
4. Set `isError: true` so the agent knows the call failed
5. Include actionable guidance — "Check that MCP_API_KEY is set" beats "Auth error"
6. Log internally, sanitize externally — full details to server logs only
7. Include cleanup routines — clean up temporary state before returning

---

## Deployment

### Transport selection

| Transport | When to use | Scaling model |
|-----------|-------------|---------------|
| **stdio** | Local servers, CLI tools, single-user | Process-per-client |
| **Streamable HTTP (stateless JSON)** | Remote servers, multi-user, cloud | Standard HTTP scaling |
| **Streamable HTTP (SSE streaming)** | Long-running operations, progress | Requires sticky sessions |

**Default:** Streamable HTTP with stateless JSON for remote servers.

### Platform selection

| Platform | Best for | Cold start | Python |
|----------|----------|------------|--------|
| Docker/K8s | Full control, any runtime | None | Full |
| Cloudflare Workers | Edge-first, low latency | Near-zero | No (V8 only) |
| Vercel | Next.js ecosystem | 1-3s | Limited |
| AWS Lambda | Event-driven, auto-scale | 1-5s | Full |
| Google Cloud Run | Container serverless | 1-3s | Full |
| Fly.io | Global containers | Near-zero | Full |

### Production health reality

An April 2026 audit of 2,181 remote MCP endpoints: 52% completely dead, only 9% fully healthy. Main causes: expired credentials, upstream API changes breaking responses silently, cold start timeouts. Fix: health checks, pinned dependencies, schema-drift monitoring.

**Health endpoint:**
```typescript
app.get("/healthz", async (req, res) => {
  const checks = {
    server: "ok",
    upstreamApi: await checkUpstreamHealth(),
    timestamp: new Date().toISOString(),
  };
  const healthy = Object.values(checks).every(v => v === "ok" || typeof v === "string");
  res.status(healthy ? 200 : 503).json(checks);
});
```

---

## Observability

Instrument every tool invocation with OpenTelemetry traces:

```typescript
import { trace, SpanStatusCode } from "@opentelemetry/api";
const tracer = trace.getTracer("mcp-server");

async function instrumentedToolCall(toolName, args, handler) {
  return tracer.startActiveSpan(`tool.${toolName}`, async (span) => {
    span.setAttribute("tool.name", toolName);
    try {
      const result = await handler();
      span.setAttribute("tool.is_error", result.isError ?? false);
      span.setStatus({ code: SpanStatusCode.OK });
      return result;
    } catch (error) {
      span.setStatus({ code: SpanStatusCode.ERROR, message: String(error) });
      throw error;
    } finally { span.end(); }
  });
}
```

**Never include raw tool arguments in trace attributes** — they may contain secrets or PII. Redact before setting span attributes.

---

## Gateway pattern (3+ servers)

Once you operate 3+ MCP servers, centralized auth, routing, and observability become necessary.

| Role | What it does | Example |
|------|-------------|---------|
| **Proxy** | Forwards MCP traffic 1:1, bridges transports | mcp-proxy, Supergateway |
| **Aggregator** | Combines multiple servers behind one endpoint | MCPJungle, MetaMCP |
| **Gateway** | Aggregator + auth + rate limiting + audit | agentgateway (Linux Foundation), Cloudflare, Kong |

Benefits: single auth endpoint, tool-level access policies, lazy-loading tool descriptions (reduces context usage to 5-15% of unmanaged case), centralized audit trail.

---

## Server discovery (.well-known/mcp.json)

Expose a JSON document at `/.well-known/mcp/server-card.json` for automated discovery:

```json
{
  "name": "my-mcp-server",
  "transport": { "type": "streamable-http", "url": "/mcp" },
  "capabilities": { "tools": true, "resources": true },
  "tools": [
    { "name": "create_issue", "description": "Create a new issue" }
  ],
  "auth": { "type": "oauth2", "authorization_url": "..." }
}
```

---

## Debugging

| Layer | Tool | What it shows |
|-------|------|---------------|
| Protocol | MCP Inspector `--verbose` | Raw JSON-RPC messages |
| Transport | Network tab / curl | HTTP status, headers, timing |
| Application | Structured logs with trace IDs | Business logic errors |
| Observability | OTel traces + dashboards | End-to-end latency, error hotspots |

**"Tool not showing in client":** run `npx @modelcontextprotocol/inspector --cli --server "node dist/index.js" tools/list` and verify it appears.

**"stdout corrupted / JSON-RPC parse errors":** something is logging to stdout. MCP uses stdout exclusively for JSON-RPC. Switch all logging to `console.error()`.
