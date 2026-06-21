<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `api-design-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: api-design-patterns
title: "API Design Patterns"
description: >-
  API design expert for REST, GraphQL, and gRPC. Covers resource modeling, Richardson
  Maturity Model, HATEOAS, GraphQL schema design, federation, N+1 prevention, Protocol
  Buffers, streaming, API versioning (URL/header/content negotiation), cursor/keyset
  pagination, RFC 9457 error responses, token-bucket/sliding-window rate limiting,
  OAuth 2.1, JWT, API keys, mTLS, idempotency keys, webhook design (HMAC, retry backoff,
  DLQ), OpenAPI 3.2, and BFF/gateway patterns.
  TRIGGER: designing or reviewing a REST/GraphQL/gRPC API; choosing between API paradigms;
  implementing pagination, error responses, versioning, rate limiting, or authentication;
  designing webhooks; writing OpenAPI specs; building an API gateway or BFF.
  SKIP: browser OAuth flows or OIDC login (use web-auth-patterns); Kafka/event streaming
  (use terraform-kafka-infra); database query design without an API surface.
version: "1.2.0"
updated: "2026-05-29"
category: developer
tags:
  - api
  - rest
  - graphql
  - grpc
  - protobuf
  - openapi
  - webhooks
  - authentication
  - pagination
  - error-handling
  - rate-limiting
  - versioning
  - hateoas
  - idempotency
keywords:
  - REST API
  - GraphQL schema
  - gRPC service
  - Protocol Buffers
  - API versioning
  - cursor pagination
  - keyset pagination
  - RFC 9457
  - problem details
  - HATEOAS
  - Richardson Maturity Model
  - OpenAPI 3.2
  - OAuth 2.1
  - JWT
  - API keys
  - rate limiting
  - token bucket
  - sliding window
  - webhook
  - HMAC signature
  - idempotency key
  - API gateway
  - BFF pattern
  - content negotiation
  - hypermedia
whenToUse:
  - "Designing a new REST, GraphQL, or gRPC API"
  - "Choosing between REST vs GraphQL vs gRPC for a project"
  - "Which HTTP status code should I return for this error?"
  - "How do I paginate this API response without offset drift?"
  - "Setting up API authentication with OAuth 2.1, JWT, or API keys"
  - "Designing webhook delivery with HMAC signatures and retry logic"
  - "Writing or reviewing an OpenAPI 3.x specification"
  - "Implementing rate limiting — token bucket vs sliding window"
  - "Building an API gateway or BFF layer for multiple clients"
  - "Reviewing an API design for anti-patterns like chatty calls or missing idempotency"
whenNotToUse:
  - "Browser-side OAuth / OIDC login flows — use web-auth-patterns"
  - "Kafka or event streaming architecture — use terraform-kafka-infra"
  - "Database schema or query design without an API layer"
related_skills:
  - web-auth-patterns
  - backend-patterns
  - microservices-patterns
  - http-security-headers
---

# API Design Patterns

Expert reference for designing, building, and reviewing APIs across REST, GraphQL, and gRPC, covering cross-cutting concerns: versioning, pagination, error handling, authentication, rate limiting, idempotency, and webhook delivery.

## 1. API Paradigm Selection

### 1.1 REST (Representational State Transfer)

REST remains the dominant API paradigm for public-facing APIs (80% of business APIs as of 2026). REST APIs expose resources via URLs and use standard HTTP methods and status codes.

**Richardson Maturity Model** (Leonard Richardson, via Martin Fowler):

| Level | Name | Description | Example |
|-------|------|-------------|---------|
| 0 | The Swamp of POX | Single URI, single verb (POST). RPC-over-HTTP | `POST /api` with action in body |
| 1 | Resources | Multiple URIs, one verb. Resource-oriented but no HTTP semantics | `POST /users`, `POST /orders` |
| 2 | HTTP Verbs | Proper use of GET/POST/PUT/DELETE + status codes. Where most "REST" APIs live | `GET /users/123`, `DELETE /orders/456` |
| 3 | Hypermedia (HATEOAS) | Responses include links to related actions. True REST | Response contains `_links` with `self`, `next`, `cancel` |

**REST Design Principles:**
- Use nouns for resources (`/users`, `/orders`), not verbs (`/getUser`)
- Use plural nouns consistently (`/users` not `/user`)
- Nest sub-resources max 2 levels deep (`/users/123/orders`, not `/users/123/orders/456/items/789`)
- Use HTTP methods semantically: GET (read), POST (create), PUT (full replace), PATCH (partial update), DELETE (remove)
- Return appropriate status codes: 200 (OK), 201 (Created), 204 (No Content), 400 (Bad Request), 401 (Unauthorized), 403 (Forbidden), 404 (Not Found), 409 (Conflict), 422 (Unprocessable Entity), 429 (Too Many Requests), 500 (Internal Server Error)
- Support filtering via query parameters (`?status=active&sort=-created_at`)
- Use `ETag` and `Last-Modified` headers for conditional requests and caching

> Sources: [Martin Fowler - Richardson Maturity Model](https://martinfowler.com/articles/richardsonMaturityModel.html), [REST API Tutorial](https://restfulapi.net/richardson-maturity-model/)

### 1.2 GraphQL

GraphQL lets clients request exactly the data they need in a single query, eliminating over-fetching and under-fetching. Adopted by 50%+ of enterprises in production as of 2026.

**Schema Design Principles:**
- Design schemas for client use cases, not database tables
- Use nullable fields by default; mark required fields with `!` intentionally
- Prefer object types over primitive fields for extensibility
- Name mutations as verb-noun (`createUser`, `cancelOrder`) and return the mutated object
- Use input types for mutation arguments
- Always paginate list fields (connection pattern or simple `[Type!]!` with `first`/`after`)
- Use DataLoader pattern to prevent N+1 queries

**Federation (Apollo Federation / GraphQL Composite Schema Spec):**
- Split schema by domain boundaries (DDD bounded contexts)
- Each subgraph owns its entity types and extends others via `@key` directives
- Entity resolution: define `@key(fields: "id")` on shared types
- Use `@external`, `@requires`, `@provides` for cross-subgraph field dependencies
- Gateway/router composes subgraphs into a supergraph and handles query planning

**Security:**
- Implement query depth limiting (max depth 7-10)
- Set query complexity analysis with cost-per-field weights
- Always paginate with limits to prevent DDoS via unbounded list queries
- Disable introspection in production
- Use persisted queries (query allowlisting) for public APIs

> Sources: [GraphQL Hive - Schema Design Best Practices Part 1](https://the-guild.dev/graphql/hive/blog/schema-design-best-practices-part-1), [Apollo - Schema Design Best Practices for Federation](https://www.apollographql.com/events/schema-design-best-practices-for-federated-graphql-apis)

### 1.3 gRPC and Protocol Buffers

gRPC uses HTTP/2 and Protocol Buffers for high-performance, type-safe inter-service communication. 5-10x faster than REST+JSON for service-to-service calls; 60% smaller payloads.

**Proto3 Design Best Practices:**
- Schema-first: write `.proto` files before implementation code
- Use wrapper types (`google.protobuf.StringValue`) when you need to distinguish "not set" from "empty"
- Reserve field numbers when removing fields (`reserved 3, 6;`) to prevent reuse
- Use `oneof` for mutually exclusive fields
- Prefer `repeated` fields over separate count + items patterns
- Name services as `<Domain>Service`, methods as `<Verb><Noun>` (`GetUser`, `ListOrders`)
- Use `google.protobuf.FieldMask` for partial updates (like PATCH in REST)

**Four RPC Patterns:**

| Pattern | Client | Server | Use Case |
|---------|--------|--------|----------|
| Unary | 1 request | 1 response | Standard request-response |
| Server streaming | 1 request | N responses | Real-time feeds, log tailing |
| Client streaming | N requests | 1 response | File uploads, bulk ingestion |
| Bidirectional streaming | N requests | N responses | Chat, collaborative editing |

**Error Handling:**
- Use standard gRPC status codes (`OK`, `INVALID_ARGUMENT`, `NOT_FOUND`, `PERMISSION_DENIED`, `INTERNAL`, etc.)
- Attach `google.rpc.Status` with detail messages for rich error context
- Never use `UNKNOWN` when a more specific code fits

**Deployment Strategy:**
- Internal: gRPC between microservices for performance
- External: REST or GraphQL with gRPC-gateway or Envoy transcoding

> Sources: [gRPC Official - Core Concepts](https://grpc.io/docs/what-is-grpc/core-concepts/)

### 1.4 Paradigm Decision Matrix

| Factor | REST | GraphQL | gRPC |
|--------|------|---------|------|
| **Best for** | Public APIs, CRUD, caching | Complex client queries, mobile BFF | Internal microservices, streaming |
| **Data format** | JSON (text) | JSON (text) | Protobuf (binary) |
| **Transport** | HTTP/1.1 or HTTP/2 | HTTP/1.1 or HTTP/2 | HTTP/2 required |
| **Schema** | Optional (OpenAPI) | Required (SDL) | Required (.proto) |
| **Caching** | Native HTTP caching | Requires client-side (Apollo, urql) | No built-in |
| **Browser support** | Native | Native | Requires gRPC-Web proxy |
| **Code generation** | Optional | Optional (codegen) | Built-in (protoc) |
| **Learning curve** | Low | Medium | Medium-High |
| **Real-time** | SSE, WebSockets | Subscriptions | Bidirectional streaming |

**Hybrid Architecture Pattern (2026 best practice):**
```
External clients --> REST (public API) or GraphQL (BFF)
                          |
                     API Gateway (Kong / Envoy / Spring Cloud Gateway)
                          |
              Internal services <--> gRPC (service mesh)
```

## 2. API Versioning

### 2.1 Versioning Strategies

| Strategy | Mechanism | Pros | Cons | Used By |
|----------|-----------|------|------|---------|
| **URL path** | `/v1/users` | Visible, cacheable, easy to test | Duplicates routes, not "pure REST" | AWS, Google, Twitter |
| **Header** | `X-API-Version: 2` | Clean URLs | Hidden, harder to test | Stripe |
| **Content negotiation** | `Accept: application/vnd.api.v2+json` | Most RESTful, supports format + version | Complex, poor tooling support | GitHub |
| **Query parameter** | `?version=2` | Easy to add | Pollutes query string, caching issues | Rarely recommended |

**Best Practices:**
- URL path versioning (`/v1/`) is the most practical default for most teams
- Only bump major versions for breaking changes; use additive changes otherwise
- Support maximum 2 concurrent versions (current + previous)
- Give consumers 6-12 months migration window before deprecating old versions
- Use `Sunset` header (RFC 8594) to signal deprecation timeline
- Centralize version routing in an API gateway
- For gRPC: use package versioning (`package myservice.v2;`)
- For GraphQL: prefer schema evolution (deprecate fields with `@deprecated`) over versioning

## 3. Pagination

### 3.1 Pagination Strategies Compared

| Strategy | Mechanism | Performance | Consistency | Page Jumping |
|----------|-----------|-------------|-------------|-------------|
| **Offset/Limit** | `?offset=20&limit=10` | Degrades with depth (DB scans skipped rows) | Inconsistent (inserts shift pages) | Yes |
| **Page-based** | `?page=3&per_page=10` | Same as offset | Same as offset | Yes |
| **Cursor-based** | `?after=eyJpZCI6MTIzfQ&limit=10` | Constant time (index seek) | Consistent (pointer-based) | No |
| **Keyset** | `?created_after=2025-01-01&id_after=123&limit=10` | Constant time (index seek) | Consistent | No |
| **Time-based** | `?since=2025-01-01T00:00:00Z` | Good with time index | Good for append-only | No |

### 3.2 Implementation Best Practices

**Cursor Pagination (recommended for most APIs):**
```json
{
  "data": [...],
  "pagination": {
    "has_more": true,
    "next_cursor": "eyJjcmVhdGVkX2F0IjoiMjAyNS0wMS0xNSIsImlkIjoiYWJjMTIzIn0=",
    "previous_cursor": "eyJjcmVhdGVkX2F0IjoiMjAyNS0wMS0xMCIsImlkIjoieHl6Nzg5In0="
  }
}
```

**Key rules:**
- Set sensible defaults: `limit` default 25, max 100 (or 1000 for bulk endpoints)
- Encode cursors as opaque Base64 strings (typically JSON of `{created_at, id}`)
- Never expose internal IDs or database column names in cursor values
- Use composite sort keys `(created_at, id)` for deterministic ordering when timestamps collide
- Always return `has_more` boolean or empty `next_cursor` to signal the final page
- Include total count only if cheap to compute (avoid `COUNT(*)` on large tables)
- For GraphQL, use the Relay Connection specification (`edges`, `node`, `pageInfo`, `cursor`)
- Index the sort columns together: `CREATE INDEX idx_pagination ON orders(created_at, id)`

**When to use offset:** Small, static datasets where users need page-jump navigation (admin tables, search results under 10K rows).

## 4. Error Handling

### 4.1 RFC 9457 Problem Details (successor to RFC 7807)

The standard error response format for HTTP APIs. Content-Type: `application/problem+json`.

```json
{
  "type": "https://api.example.com/errors/insufficient-funds",
  "title": "Insufficient Funds",
  "status": 422,
  "detail": "Account 12345 has a balance of $10.00 but the transfer requires $50.00.",
  "instance": "/transfers/abc-123",
  "balance": 1000,
  "currency": "USD"
}
```

**Standard Fields:**
| Field | Type | Required | Purpose |
|-------|------|----------|---------|
| `type` | URI | Yes (defaults to `about:blank`) | Machine-readable error category. Ideally a resolvable URL pointing to documentation |
| `title` | string | Yes | Short human-readable summary (should not change between occurrences) |
| `status` | integer | No | HTTP status code (mirrors response status for convenience) |
| `detail` | string | No | Human-readable explanation specific to this occurrence |
| `instance` | URI | No | URI reference identifying the specific occurrence |

**Best Practices:**
- Always use `application/problem+json` content type for error responses
- Make `type` URIs resolvable documentation links when possible
- Never leak stack traces, internal paths, or sensitive data in `detail`
- Keep `title` static per error type; put dynamic info in `detail`
- Use extension members for machine-actionable context (retry-after, field validation errors)
- For validation errors, include a `violations` array with `field`, `message`, `code`

**Validation Error Extension Pattern:**
```json
{
  "type": "https://api.example.com/errors/validation",
  "title": "Validation Error",
  "status": 422,
  "detail": "2 fields failed validation.",
  "violations": [
    { "field": "email", "message": "Must be a valid email address", "code": "INVALID_FORMAT" },
    { "field": "age", "message": "Must be at least 18", "code": "MIN_VALUE" }
  ]
}
```

### 4.2 HTTP Status Code Reference

| Range | Meaning | Common Codes |
|-------|---------|--------------|
| 2xx | Success | 200 OK, 201 Created, 202 Accepted, 204 No Content |
| 3xx | Redirection | 301 Moved Permanently, 304 Not Modified |
| 4xx | Client error | 400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found, 405 Method Not Allowed, 409 Conflict, 422 Unprocessable Entity, 429 Too Many Requests |
| 5xx | Server error | 500 Internal Server Error, 502 Bad Gateway, 503 Service Unavailable, 504 Gateway Timeout |

**Guidelines:**
- Use 201 for resource creation (include `Location` header)
- Use 202 for async operations (include status polling URL)
- Use 204 for successful DELETE (no body)
- Use 409 for state conflicts (optimistic concurrency violations)
- Use 422 for valid syntax but semantic errors (validation failures)
- Use 429 with `Retry-After` header for rate limiting

## 5. Rate Limiting

### 5.1 Algorithms

| Algorithm | How It Works | Burst Handling | Memory | Best For |
|-----------|-------------|----------------|--------|----------|
| **Token Bucket** | Tokens added at fixed rate; each request consumes one | Allows controlled bursts (up to bucket capacity) | Low | Public APIs (default choice) |
| **Leaky Bucket** | Requests queued and processed at fixed rate | Smooths all bursts | Low | Strict throughput control |
| **Fixed Window** | Counter per time window (e.g., per minute) | Boundary burst (2x limit possible at window edges) | Low | Simple implementations |
| **Sliding Window Log** | Stores timestamp of each request; counts in rolling window | No boundary burst | High | Security-sensitive endpoints |
| **Sliding Window Counter** | Weighted average of current + previous window | Minimal boundary burst | Low | High-scale distributed systems |

### 5.2 Implementation Best Practices

**Response Headers (standard):**
```
X-RateLimit-Limit: 1000          # Max requests per window
X-RateLimit-Remaining: 742       # Requests remaining
X-RateLimit-Reset: 1672531200    # Unix timestamp when window resets
Retry-After: 30                  # Seconds to wait (on 429 responses)
```

**Tiered Limits:**

| Tier | Rate Limit | Use Case |
|------|-----------|----------|
| Anonymous | 10 req/min | Public browsing |
| Free | 100 req/min | Registered developers |
| Pro | 1,000 req/min | Paid plans |
| Enterprise | 10,000 req/min | Custom SLA |

**Key practices:**
- Use Redis or a distributed store for rate limit counters in multi-instance deployments
- Token bucket is the strongest default for public APIs (predictable rate + controlled bursts)
- Implement dynamic rate limiting that adjusts based on server load
- Apply per-endpoint rate limits for expensive operations (search, reports)
- Always return `429 Too Many Requests` with `Retry-After` header
- Document rate limits in API documentation and OpenAPI spec
- Consider cost-based rate limiting: assign different costs to different endpoints

## 6. HATEOAS (Hypermedia as the Engine of Application State)

HATEOAS is Level 3 of the Richardson Maturity Model. API responses include hypermedia links that tell clients what actions they can take next, enabling self-documenting, evolvable APIs without hardcoded URL structures.

**HAL+JSON Example (`application/hal+json`):**
```json
{
  "id": "order-123",
  "status": "pending",
  "total": 59.99,
  "_links": {
    "self": { "href": "/orders/order-123" },
    "cancel": { "href": "/orders/order-123/cancel", "method": "POST" },
    "payment": { "href": "/orders/order-123/payment", "method": "PUT" },
    "items": { "href": "/orders/order-123/items" }
  }
}
```

**Key benefits:**
- **Discoverability:** Clients discover available actions from responses, not documentation
- **Decoupling:** Server can change URLs without breaking clients (clients follow links)
- **Evolvability:** Add new actions by adding links; remove actions by omitting links
- **State-driven UI:** Links reflect business state (e.g., `cancel` link only appears when order is cancellable)

**Media Types:**
- HAL (`application/hal+json`): most widely adopted
- JSON:API (`application/vnd.api+json`): includes pagination, sparse fieldsets
- Siren (`application/vnd.siren+json`): includes actions with field definitions

**AI agents (2025-2026):** AI agents can navigate HATEOAS APIs autonomously by following hypermedia links, making HATEOAS newly relevant for agent-driven API consumption.

## 7. Authentication and Authorization

> **Cross-reference:** For in-depth browser auth patterns (passkeys/WebAuthn, OIDC flows, SAML vs OIDC, MFA, token storage, CORS credentials, CSP, session hijacking prevention, JWT algorithm confusion attacks), see the `web-auth-patterns` skill.

### 7.1 Authentication Methods Compared

| Method | Best For | Security Level | Complexity |
|--------|----------|---------------|------------|
| **API Keys** | Server-to-server, internal tools | Low-Medium | Low |
| **OAuth 2.1 + PKCE** | End-user login, delegated access, third-party apps | High | High |
| **JWT (Bearer)** | Stateless auth, microservices | Medium-High | Medium |
| **mTLS** | Service mesh, zero-trust microservices | Very High | High |
| **Basic Auth** | Development/testing only | Low | Low |

### 7.2 OAuth 2.1 Key Changes

OAuth 2.1 (current as of 2026) consolidates OAuth 2.0 security best practices:
- **PKCE required** for all clients (not just public clients)
- **Implicit grant deprecated** — use Authorization Code + PKCE instead
- **Password grant deprecated** — use device authorization flow
- Exact string matching for redirect URIs (no wildcards)
- Refresh token rotation mandatory (one-time use)
- Access tokens: short-lived (5-15 minutes), JWT format recommended
- Scope to minimum permissions (principle of least privilege)

### 7.3 JWT Best Practices

- Use asymmetric signing (RS256 or ES256) for distributed verification
- Keep payloads minimal (user ID, roles, scopes — not full user profile)
- Short expiration (5-15 min access, 7-30 day refresh)
- Implement refresh token rotation with reuse detection
- Never store JWTs in localStorage (use httpOnly cookies or in-memory)
- Validate `iss`, `aud`, `exp`, `nbf` claims on every request
- Use `jti` claim for token revocation lists

### 7.4 API Key Best Practices

- Rotate keys quarterly at minimum
- Use key versioning so old keys remain valid during rotation window
- Hash keys at rest (treat like passwords)
- Scope keys to specific endpoints or operations
- Send via `Authorization: Bearer <key>` header (not query parameter)
- Implement per-key rate limits and usage tracking

## 8. Idempotency

### 8.1 HTTP Method Idempotency

| Method | Idempotent | Safe | Notes |
|--------|-----------|------|-------|
| GET | Yes | Yes | Always idempotent by spec |
| HEAD | Yes | Yes | Same as GET without body |
| PUT | Yes | No | Full replace is naturally idempotent |
| DELETE | Yes | No | Deleting twice yields same state |
| POST | **No** | No | Requires explicit idempotency design |
| PATCH | **No** | No | Depends on operation semantics |

### 8.2 Idempotency Key Pattern (for POST/PATCH)

**How it works (Stripe model):**
1. Client generates a UUID v4 idempotency key before sending the request
2. Client sends `Idempotency-Key: <uuid>` header with the request
3. Server processes on first call, stores result keyed by idempotency key
4. Subsequent requests with same key return the cached result (including errors)
5. Server expires idempotency keys after 24-48 hours

**Implementation rules:**
- Store the full response (status code + body) for the idempotency key
- Return cached result even if the original request returned an error (5xx)
- Use database-level locking or compare-and-swap to prevent race conditions
- Set idempotency key TTL to exceed the client's retry window
- Return `409 Conflict` if same key is used with different request body
- Apply to all state-changing operations, not just payment endpoints

## 9. Webhook Design

### 9.1 Webhook Architecture

Webhooks are server-to-server HTTP callbacks that push event notifications in near real-time, inverting the polling pattern.

**Delivery Flow:**
```
Event occurs --> Serialize payload --> Sign with HMAC --> HTTP POST to subscriber URL
    |                                                           |
    +--- If 2xx: mark delivered                                 |
    +--- If 5xx/timeout: enqueue for retry (exponential backoff)
    +--- After max retries: move to dead-letter queue (DLQ)
```

### 9.2 Security: HMAC Signature Verification

```javascript
// Provider (sending side)
const crypto = require('crypto');
const signature = crypto
  .createHmac('sha256', webhookSecret)
  .update(rawBody)
  .digest('hex');
// Send as header: X-Signature: sha256=<signature>

// Consumer (receiving side)
const expected = crypto
  .createHmac('sha256', webhookSecret)
  .update(rawBody)
  .digest('hex');
const isValid = crypto.timingSafeEqual(
  Buffer.from(signature),
  Buffer.from(expected)
);
```

**Key security rules:**
- Always use HTTPS (never HTTP)
- Sign payloads with HMAC-SHA256 using a shared secret
- Use timing-safe comparison to prevent timing attacks
- Include a timestamp in the signature to prevent replay attacks
- Rotate signing keys periodically (expose via JWKS endpoints)

### 9.3 Retry Strategy

**Standard exponential backoff schedule:**

| Attempt | Delay | Cumulative |
|---------|-------|------------|
| 1 | Immediate | 0s |
| 2 | 1 minute | 1m |
| 3 | 5 minutes | 6m |
| 4 | 15 minutes | 21m |
| 5 | 1 hour | 1h 21m |
| 6 | 2 hours | 3h 21m |

**Rules:**
- Only retry 5xx errors and network timeouts; discard 4xx (client error) immediately
- Add random jitter (0-1 second) to prevent thundering herd
- Cap at 5-6 retry attempts
- After max retries, move to a dead-letter queue for manual inspection
- Set a 5-second timeout per delivery attempt
- Verify signature, enqueue to background job, return 200 immediately

### 9.4 Consumer Best Practices

- Return 200 immediately after signature verification; process asynchronously
- Implement idempotency: deduplicate by event ID with a TTL exceeding the retry window
- Store raw webhook payloads for debugging and replay
- Monitor webhook processing lag and failure rates

## 10. OpenAPI Specification

### 10.1 Current State (2026)

- **OpenAPI 3.2.0** (September 2025) is the latest feature release
- New in 3.2: structured tag nesting, streaming media types (SSE, JSON Lines), native QUERY method, OAuth 2.0 Device Authorization Flow
- 82% of companies self-identify as API-first (Postman 2025 State of the API)

### 10.2 Design-First Workflow

```
1. Define OpenAPI spec (YAML/JSON)
       |
2. Generate mock server (Prism, Stoplight)
       |
3. Frontend + backend develop in parallel against the contract
       |
4. Generate server stubs + client SDKs (openapi-generator, Speakeasy)
       |
5. Validate implementation against spec in CI (Spectral, openapi-diff)
       |
6. Publish documentation (Redoc, Swagger UI, Scalar)
```

**Best Practices:**
- Write OpenAPI spec before implementation code (contract-first)
- Use `$ref` for shared schemas to avoid duplication
- Define reusable error responses as shared components
- Include examples for every endpoint (request + response)
- Add `description` to every field (not just complex ones)
- Use Spectral linting rules in CI to enforce API design standards
- Version your OpenAPI spec alongside your code

## 11. API Gateway and BFF Patterns

### 11.1 API Gateway

An API gateway sits between clients and backend services, handling cross-cutting concerns:
- **Routing:** Version-based, path-based, header-based
- **Authentication:** Validate tokens before requests reach services
- **Rate limiting:** Centralized enforcement
- **Protocol translation:** REST-to-gRPC, HTTP-to-WebSocket
- **Request aggregation:** Combine multiple service calls into one response
- **Caching:** Response caching with configurable TTLs
- **Observability:** Logging, metrics, distributed tracing

**Popular gateways:** Kong, Envoy, AWS API Gateway, Spring Cloud Gateway, Traefik, APISIX

### 11.2 Backend for Frontend (BFF) Pattern

Create a dedicated backend service for each frontend type (web, mobile, TV) that:
- Aggregates data from multiple microservices
- Transforms responses to match frontend needs
- Handles authentication flows specific to the platform
- Prevents over-fetching by returning exactly what the UI needs
- Eliminates chatty API anti-patterns

**When to use BFF:**
- Multiple clients with different data needs
- Mobile apps requiring optimized payloads
- Frontend teams needing autonomy over their API layer
- Complex aggregation that would be expensive on the client

## 12. Anti-Patterns and Troubleshooting

### 12.1 Common API Anti-Patterns

| Anti-Pattern | Problem | Solution |
|-------------|---------|----------|
| **Chatty API** | Client needs 10+ calls for one screen | Aggregate endpoints or BFF pattern |
| **God endpoint** | Single endpoint handles all operations | Split into proper REST resources |
| **Wrong HTTP verbs** | Using POST for reads, GET for deletes | Follow HTTP method semantics |
| **Missing pagination** | Returning unbounded lists | Always paginate collections |
| **No versioning** | Breaking changes break all clients | URL or header versioning |
| **Leaking internals** | DB column names, stack traces in responses | Map to public DTOs, sanitize errors |
| **No idempotency** | Duplicate charges on retry | Idempotency keys for state changes |
| **Over-fetching** | Returning 50 fields when client needs 3 | Sparse fieldsets or GraphQL |
| **Under-fetching** | Returning too little, forcing N+1 calls | Include related data or `?expand=` |
| **Ignoring caching** | No ETags, no Cache-Control | Add conditional request support |
| **String-typed everything** | Dates, enums, IDs all as untyped strings | Use proper types, ISO 8601 dates |
| **Inconsistent naming** | Mix of camelCase, snake_case, PascalCase | Pick one convention, enforce via linting |

### 12.2 Troubleshooting Checklist

- **Slow responses:** Check N+1 queries, missing indexes, unnecessary data fetching
- **Intermittent 5xx:** Check connection pool exhaustion, circuit breaker state, upstream timeouts
- **Rate limit exceeded:** Verify tier, check for retry loops, implement exponential backoff
- **Auth failures:** Validate token expiry, clock skew, audience/issuer claims
- **Pagination inconsistency:** Switch from offset to cursor-based pagination
- **Webhook delivery failures:** Check HTTPS, signature validation, response timeout, DLQ

## 13. References

**Standards:** [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457.html) | [RFC 8594 Sunset Header](https://www.rfc-editor.org/rfc/rfc8594) | [OpenAPI 3.2.0](https://spec.openapis.org/oas/v3.2.0.html) | [gRPC Docs](https://grpc.io/docs/) | [GraphQL Spec](https://spec.graphql.org/)

**REST:** [Richardson Maturity Model (Fowler)](https://martinfowler.com/articles/richardsonMaturityModel.html) | [REST API Tutorial](https://restfulapi.net/) | [Microsoft Web API Best Practices](https://learn.microsoft.com/en-us/azure/architecture/best-practices/api-design)

**GraphQL:** [Schema Design Pt1](https://the-guild.dev/graphql/hive/blog/schema-design-best-practices-part-1) | [Schema Design Pt2](https://the-guild.dev/graphql/hive/blog/schema-design-best-practices-part-2) | [Federation Design (Apollo)](https://www.apollographql.com/blog/backend/federation/federated-schema-design/)

**gRPC:** [Core Concepts](https://grpc.io/docs/what-is-grpc/core-concepts/)

**Auth & Security:** [API Security OAuth 2.1/JWT (daily.dev)](https://daily.dev/blog/dev-guide-api-security-oauth-2-1-jwt-vulnerabilities)

**Pagination:** [Developer's Guide (Gusto)](https://embedded.gusto.com/blog/api-pagination/) | [Cursor vs Offset vs Keyset (Design Gurus)](https://designgurus.substack.com/p/api-pagination-guide-cursor-vs-offset)

**Rate Limiting:** [Algorithms Guide (API7)](https://api7.ai/blog/rate-limiting-guide-algorithms-best-practices) | [Token Bucket vs Sliding Window (Arcjet)](https://blog.arcjet.com/rate-limiting-algorithms-token-bucket-vs-sliding-window-vs-fixed-window/)

**Webhooks:** [Idempotency (Hookdeck)](https://hookdeck.com/webhooks/guides/implement-webhook-idempotency) | [Idempotent Requests (Stripe)](https://docs.stripe.com/api/idempotent_requests)
