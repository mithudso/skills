<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Formerly the standalone `a2a-interop` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: a2a-interop
version: 1.1.0
updated: 2026-05-29
description: >
  A2A protocol and agent interoperability expert — A2A architecture (Agent Cards,
  8-state task lifecycle, JSON-RPC 2.0), building A2A agents, A2A vs MCP
  (horizontal vs vertical), other standards (ACP merged into A2A, ANP, AAIF
  governance), and cross-framework bridging (LangGraph/CrewAI/ADK/Claude SDK
  via A2A+MCP). TRIGGER: user asks about A2A protocol, Agent Cards, agent-to-agent
  communication, multi-vendor agent interop, or combining A2A with MCP.
  SKIP: pure MCP tooling questions with no cross-agent delegation; general HTTP
  API design unrelated to agent protocols.
origin: local
tags: [a2a, mcp, agent-interop, protocols, multi-agent, aaif]
related_skills: [agent-ecosystem, mcp-servers, mcp-builder]
---

# A2A Protocol & Agent Interoperability

Reference for the Agent-to-Agent protocol and cross-framework agent communication.
Deep spec details are in `references/a2a-interop-context.md`.

## When to use this skill

Activate when the user:

- asks about the A2A protocol or agent-to-agent communication
- wants to publish an Agent Card at `.well-known/agent.json`
- needs to implement A2A task lifecycle handling
- asks how A2A and MCP work together
- wants cross-framework agent interop (LangGraph <-> CrewAI <-> ADK <-> Claude SDK)
- asks about AAIF, ACP, ANP, or agent protocol governance
- needs to design a multi-vendor agent system

## When NOT to use this skill

- Pure MCP server/tool questions with no cross-agent delegation — use `mcp-servers` or `mcp-builder`
- General agent framework selection without A2A requirements — use `agent-ecosystem`
- Generic HTTP API or microservices design — use `api-design-patterns`

## Quick reference

### A2A vs MCP

| Aspect | MCP | A2A |
| --- | --- | --- |
| Direction | Vertical (agent -> tools) | Horizontal (agent -> agent) |
| Abstraction | Stateless function calls | Stateful multi-turn tasks |
| Discovery | Server capabilities | Agent Cards at `.well-known/agent.json` |
| Transport | stdio / Streamable HTTP | JSON-RPC 2.0 over HTTPS + SSE |
| State management | None | 8-state task lifecycle |
| Governance | AAIF (Linux Foundation) | AAIF (Linux Foundation) |
| Use together | Each agent uses MCP for its tools | A2A for inter-agent coordination |

### Task lifecycle (8 states)

```
SUBMITTED -> WORKING ----+----> COMPLETED (terminal)
               |         |----> FAILED (terminal)
               |         +----> CANCELED (terminal)
               |
               +----> INPUT_REQUIRED (interrupted) --+
               |                                      |
               +----> AUTH_REQUIRED (interrupted) ----+-> client sends new message
               |
               +----> REJECTED (terminal)
```

Terminal states accept no further messages. Interrupted states allow the client to continue by sending additional messages.

### Protocol stack (2026)

| Layer | Protocol | Purpose |
| --- | --- | --- |
| Tool access | MCP | Agent connects to tools/data |
| Agent collaboration | A2A | Agents discover and delegate to each other |
| Open discovery | ANP | Decentralized peer-to-peer (emerging) |

### How MCP + A2A work together

```
User Request
    |
[Orchestrator Agent]
    +--- A2A ---> [Research Agent]
    |                 +--- MCP ---> Web Search Tool
    |                 +--- MCP ---> Document Store
    +--- A2A ---> [Analysis Agent]
    |                 +--- MCP ---> Database
    +--- A2A ---> [Report Agent]
                      +--- MCP ---> Email Service
```

Each agent uses MCP for its own tools. Agents coordinate with each other via A2A.

## Agent Card

An Agent Card is a JSON metadata document published at `/.well-known/agent.json`.

**Minimal required fields:** `name`, `description`, `version`, `default_input_modes`, `default_output_modes`, `capabilities`, `supported_interfaces`, `skills`.

**Key capabilities flags:** `streaming`, `push_notifications`, `extended_agent_card`.

**Security:** Declare schemes (OAuth2, mTLS, API key, OIDC) in `security_schemes`; reference them in `security`. Sign the card with JWS (`signature` field) in zero-trust environments.

Extended card (authenticated endpoint) can expose additional skills or capabilities not in the public card.

See `references/a2a-interop-context.md` for the full Agent Card JSON schema.

## Core JSON-RPC operations

| Operation | Method | Purpose |
| --- | --- | --- |
| SendMessage | `tasks/send` | Initiate or continue a task |
| SendStreamingMessage | `tasks/sendSubscribe` | Real-time SSE updates |
| GetTask | `tasks/{id}` | Retrieve current task state |
| ListTasks | `tasks` | Query tasks with filters + pagination |
| CancelTask | `tasks/{id}/cancel` | Request cancellation (idempotent) |
| SubscribeToTask | `tasks/{id}/subscribe` | Stream updates for existing task |

**Required header:** `A2A-Version: Major.Minor` in all requests. Missing version defaults to 0.3; mismatch returns `VersionNotSupportedError`.

## Error codes

| Error | When |
| --- | --- |
| `TaskNotFoundError` | Task ID invalid or inaccessible |
| `TaskNotCancelableError` | Task in a non-cancellable state |
| `PushNotificationNotSupportedError` | Agent didn't declare push capability |
| `ContentTypeNotSupportedError` | Media type rejected |
| `VersionNotSupportedError` | Protocol version incompatible |

## Enterprise orchestration models

| Model | Best for |
| --- | --- |
| Hub-and-spoke | Small-to-medium agent networks |
| Mesh | Decentralized, fault-tolerant systems |
| Hierarchical | Large enterprise with domain boundaries |
| Hybrid | Most real-world deployments |

**N-squared problem:** Direct peer-to-peer HTTP/gRPC connections grow O(N²). Mitigation: central registry, API gateway, message broker, or hierarchical orchestration.

## Cross-framework support

| Framework | A2A | MCP | Bridge pattern |
| --- | --- | --- |  --- |
| Google ADK | Native | Via MCP toolkit | ADK agents expose A2A natively |
| LangGraph | Via adapter | Native | Wrap graphs as A2A servers |
| CrewAI | Via adapter | Via tool wrapping | Expose crews as A2A endpoints |
| Claude SDK | Via A2A client | Native MCP | Call A2A servers as remote agents |
| OpenAI Agents SDK | Via A2A client | Native MCP | Same as Claude SDK |

**General bridge pattern:** (1) Create A2A server with Agent Card, (2) route `tasks/send` to the framework's execution engine, (3) map internal state to A2A task states, (4) return results as artifacts. A2A treats agents as opaque — internal implementation doesn't matter.

## Governance (2026)

- **AAIF** (Agentic AI Foundation, Linux Foundation): neutral governance body for A2A and MCP. Co-founded by OpenAI, Anthropic, Google, Microsoft, AWS, and Block.
- **ACP**: merged into A2A under AAIF.
- **ANP**: emerging decentralized peer-to-peer discovery layer above A2A.
- **A2A v1.0**: gRPC support, signed Agent Cards, multi-tenancy.

## Common pitfalls

1. **7 vs 8 states**: A2A has 8 task states (adds AUTH_REQUIRED to the original 7).
2. **Terminal states**: Completed, Failed, Canceled, Rejected accept no further messages.
3. **Missing Agent Card**: Without one, other agents cannot discover yours.
4. **Stateless thinking**: A2A tasks are stateful; design for multi-turn.
5. **No INPUT_REQUIRED handler**: Client agents must respond when a remote agent requests more info.
6. **Missing version header**: Always send `A2A-Version`; omitting it defaults to 0.3.
7. **N-squared scaling**: Direct peer-to-peer doesn't scale; use registries or gateways.

## Implementation checklists

### Publishing an Agent Card
- [ ] All skills have descriptions and example inputs
- [ ] Accurate `input_modes` / `output_modes` MIME types
- [ ] Security schemes declared and match your auth infrastructure
- [ ] Served at `/.well-known/agent.json` over HTTPS
- [ ] JWS signature if zero-trust environment
- [ ] Validated against A2A Agent Card schema

### Implementing an A2A server
- [ ] `tasks/send` as minimum viable operation
- [ ] `tasks/sendSubscribe` if streaming declared in capabilities
- [ ] All 8 task states handled correctly
- [ ] Proper A2A error codes (not generic HTTP errors)
- [ ] `A2A-Version` header in all responses
- [ ] `tasks/cancel` for long-running operations
- [ ] TLS + authentication per Agent Card declaration

## References

- [A2A Protocol Specification](https://a2a-protocol.org/latest/specification/)
- [A2A GitHub Repository](https://github.com/a2aproject/A2A)
- [Google ADK A2A Quickstart](https://google.github.io/adk-docs/a2a/quickstart-exposing/)
- [Awesome A2A](https://github.com/ai-boost/awesome-a2a)
- Full schema examples and Python/TypeScript code samples: `references/a2a-interop-context.md`
