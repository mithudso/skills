<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `spec-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: spec-writing
description: Engineering implementation spec craft — the contract-shaped genre that says WHAT a system must do, distinct from PRDs (which describe what to build at the product level) and RFCs (which propose architecture). Covers API specs (OpenAPI 3.1 for HTTP, AsyncAPI 3.0 for events), behavior specs (Gherkin / Given-When-Then), data specs (schema contracts), Joel Spolsky's "Painless Functional Specifications" tradition, the spec-as-contract mindset, versioning conventions (SemVer for APIs, deprecation windows), and the spec-file-vs-code-doc decision. TRIGGER when user says "write a spec", "API spec", "OpenAPI spec", "AsyncAPI", "Gherkin scenarios", "BDD spec", "functional spec", "behavior spec", "spec this endpoint", "contract for X", "what should this API return", "spec versioning", or asks to define the contract/behavior of a component independent of implementation. SKIP for product-level requirements (use prd-writing), architecture proposals and tradeoff analysis (use rfc-and-design-docs), code-task plans (use code-plan-writing), software architecture decisions (use software-architect), REST/GraphQL design choices (use api-design-patterns — design vs documentation), reference docs published to end-users (use api-docs-craft), and user-facing acceptance criteria standalone (use user-story-and-acceptance-criteria).
category: custom
tags: [writing, engineering-writing, spec, api, contract]
---

# Spec Writing

## Overview

An engineering spec is a contract. It says **WHAT** a system, component, endpoint, or message must do — independent of **HOW** it is implemented. Specs are read by the engineers who build the thing, the engineers who consume the thing, the QA engineers who validate the thing, and the future engineers who modify the thing six months from now.

Three forms dominate modern practice:

- **API specs** — machine-readable contracts for HTTP (OpenAPI 3.1) and event-driven (AsyncAPI 3.0) interfaces
- **Behavior specs** — executable narratives (Gherkin / Given-When-Then) describing how the system behaves under specific conditions
- **Data specs** — schema contracts (JSON Schema, Avro, Protobuf) defining the shape and constraints of data at rest or in transit

Joel Spolsky's 2000 series "Painless Functional Specifications" introduced the discipline to mass programming culture. His core distinction still holds: a **functional spec** describes how a product works from the user's perspective; a **technical spec** describes internal implementation. This skill is about the functional/contract layer — what other code can assume about your code.

The modern shift: specs are increasingly **machine-readable** (OpenAPI, AsyncAPI, JSON Schema) and **executable** (Gherkin scenarios drive Cucumber tests). Prose specs still exist for non-API behavior, but the bias is toward formats that tooling can validate, mock, and test against.

## Core Concepts

### 1. Spec says WHAT, not HOW

Joel Spolsky: "A functional specification describes how a product will work entirely from the user's perspective. It doesn't care how the thing is implemented." This applies one layer down too — the consumer of your API is a "user" of it.

A spec that says "internally we'll use a Redis cache with 5-minute TTL" has leaked implementation. The consumer cannot rely on that detail; it may change. A spec that says "responses are guaranteed fresh within 5 minutes" is a contract — the consumer can plan around it.

Litmus test: if you could swap out the implementation entirely (language, database, infrastructure) and the spec still holds, the spec is at the right level.

### 2. Contract-first design

Contract-first development: agree on the contract before writing code. The OpenAPI / AsyncAPI document is the source of truth. Server and client implementations are derived from (or validated against) it.

Benefits:
- Parallel work — frontend mocks against the spec while backend implements
- Pre-merge contract testing — implementation diffs against the spec in CI
- Documentation comes free — the spec IS the docs
- Cross-team coordination — one canonical artifact

When NOT to do contract-first: pure internal services with one team owning both sides, exploratory prototypes still finding their shape, or specs being reverse-engineered from existing systems (write code-first then capture the spec).

### 3. OpenAPI 3.1 for HTTP APIs

OpenAPI 3.1 is the dominant standard for REST/HTTP APIs. Key top-level fields:

- `openapi` — version (3.1.0)
- `info` — name, version, description, contact, license
- `servers` — base URLs per environment
- `paths` — endpoints, verbs, parameters, request/response shapes
- `components` — reusable schemas, parameters, responses, security schemes
- `security` — auth applied globally
- `tags` — grouping for docs

OpenAPI 3.1 is fully aligned with JSON Schema 2020-12 (3.0.x was not), so schemas can be shared with validation libraries directly.

YAML is preferred over JSON for human authoring — more compact, supports comments. JSON for tooling output.

### 4. AsyncAPI 3.0 for events

AsyncAPI is the event-driven equivalent of OpenAPI. v3.0 introduced a major restructuring:

- **Channels** are addressable destinations (Kafka topic, AMQP queue, WebSocket path) and define the messages they carry — **decoupled** from operations
- **Operations** describe what an application does on a channel using `action: send` or `action: receive` (replacing v2's confusing `publish`/`subscribe`)
- **Messages** are defined once and referenced from channels and operations

The decoupling means one channel can be reused across operations (one app sends, another receives the same channel) — much cleaner for Kafka/topic-style architectures.

### 5. Gherkin / Given-When-Then for behavior specs

Gherkin is the structured plain-text format used by Cucumber and BDD frameworks. The Given-When-Then triplet, developed by Daniel Terhorst-North and Chris Matts:

- **Given** — preconditions; state of the world before the behavior
- **When** — the triggering action or event
- **Then** — the expected observable outcome

Why it works as a spec format:

- Readable by non-engineers (PMs, QA, support, legal)
- Executable — Cucumber, Behave, SpecFlow run the steps as automated tests
- Composes via `Background:` (shared preconditions) and `Scenario Outline:` (parameterized cases)
- Captures edge cases as separate scenarios rather than hand-waving them in prose

Use Gherkin for: business-logic behavior (pricing rules, eligibility checks, state transitions), workflows that span multiple systems, anything where the spec doubles as the test plan.

Avoid Gherkin for: API shape (use OpenAPI), data format (use JSON Schema), simple CRUD without conditional logic.

### 6. Spec-as-contract mindset

A spec is a promise. Once published and consumed by another team or external user, breaking it is a contract violation. This forces three disciplines:

1. **Versioning** — non-breaking changes bump minor (1.1 → 1.2); breaking changes bump major (1.x → 2.0) and run in parallel for a deprecation window
2. **Backward compatibility** — additive changes (new optional field, new endpoint) are safe; removals, renames, and type changes are breaking
3. **Deprecation policy** — published timeline for retiring old versions (e.g., 12 months notice for major version sunset)

Stripe is the canonical example: every API change is dated, old versions remain accessible by clients pinning to that date, and clients see a defined deprecation window.

### 7. The "implementation discussion" anti-pattern

Specs frequently rot into design discussions. A correct OpenAPI spec says `responses: 200: schema: User`. An incorrect one has comments like `// TODO: should this be cached at edge?` or `// we may switch to gRPC later`.

When you find yourself wanting to discuss HOW in the spec, that text belongs in:

- An **RFC / design doc** if it's an architecture decision being proposed
- A **code comment** if it's an implementation note
- A **changelog** if it's a historical decision worth preserving
- A **README** if it's onboarding context

The spec stays clean and contractual.

### 8. Versioning conventions

For HTTP APIs:

- **URL versioning** (`/v1/users`) — explicit, cacheable, but proliferates URLs
- **Header versioning** (`Accept: application/vnd.acme.v1+json`) — clean URLs, harder to debug
- **Date versioning** (Stripe: `Stripe-Version: 2024-04-10`) — fine-grained, pins client behavior
- **Query parameter** (`?version=2`) — discouraged; pollutes the URL

For events (AsyncAPI):

- Version in the topic name (`orders.v1`, `orders.v2`) — common in Kafka
- Version in the message envelope (`{ "version": 2, "data": {...} }`) — common in CloudEvents-style schemas
- Schema registry with explicit compatibility modes (BACKWARD, FORWARD, FULL) — Confluent Schema Registry pattern

For data schemas:

- SemVer on schema file (`user.schema.v1.2.0.json`)
- Schema registry with evolution rules
- Field-level deprecation (`"deprecated": true` flags)

### 9. Spec file vs code doc decision

When does behavior live in a spec file vs an inline code comment vs a README?

**Spec file (separate artifact):**
- Contract crosses a team boundary
- More than one implementation must conform
- The artifact is published to consumers (SDK generation, public docs)
- Tooling validates against it (CI checks, contract tests)

**Inline code doc / JSDoc / docstring:**
- Behavior is purely internal
- Single implementation, single team
- Changes track 1:1 with code changes

**README / runbook:**
- Operational behavior (how to deploy, how to debug)
- Setup/usage rather than contract

A common pattern: the spec file is the source of truth for the contract; code comments link to the spec section; the README links to the spec for consumers.

### 10. Examples are part of the spec

OpenAPI's `examples` field, AsyncAPI's `payload.examples`, and Gherkin's `Examples:` table are not decorative — they are the spec. Reviewers and consumers learn faster from a concrete example than from a JSON Schema. CI tools also validate that examples conform to the declared schema, catching drift.

Rule: every endpoint, every message, every Gherkin scenario outline should have at least one realistic example. Use realistic IDs, plausible values, and at least one error case.

## Templates and Examples

### OpenAPI 3.1 minimal endpoint (HTTP)

```yaml
openapi: 3.1.0
info:
  title: Orders API
  version: 1.4.0
  description: Order lifecycle management for the storefront.
servers:
  - url: https://api.example.com/v1
    description: Production
  - url: https://api-staging.example.com/v1
    description: Staging
paths:
  /orders/{orderId}:
    get:
      operationId: getOrder
      summary: Retrieve a single order by ID.
      parameters:
        - name: orderId
          in: path
          required: true
          schema:
            type: string
            format: uuid
          examples:
            standard:
              value: "8a2e1b7c-9f4d-4e3a-b1c2-3d4e5f6a7b8c"
      responses:
        "200":
          description: Order found.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Order"
              examples:
                paid:
                  $ref: "#/components/examples/PaidOrder"
        "404":
          description: Order not found.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Error"
        "401":
          description: Unauthenticated.
      security:
        - bearerAuth: []
components:
  schemas:
    Order:
      type: object
      required: [id, status, total, currency, createdAt]
      properties:
        id:
          type: string
          format: uuid
        status:
          type: string
          enum: [pending, paid, shipped, refunded, cancelled]
        total:
          type: integer
          description: Amount in smallest currency unit (cents).
          minimum: 0
        currency:
          type: string
          pattern: "^[A-Z]{3}$"
        createdAt:
          type: string
          format: date-time
    Error:
      type: object
      required: [code, message]
      properties:
        code: { type: string }
        message: { type: string }
  examples:
    PaidOrder:
      value:
        id: "8a2e1b7c-9f4d-4e3a-b1c2-3d4e5f6a7b8c"
        status: "paid"
        total: 2999
        currency: "USD"
        createdAt: "2026-05-29T14:23:00Z"
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
```

### AsyncAPI 3.0 event spec (event-driven)

```yaml
asyncapi: 3.0.0
info:
  title: Order Events
  version: 2.1.0
servers:
  production:
    host: kafka.example.com:9092
    protocol: kafka
channels:
  orderPlaced:
    address: orders.placed.v2
    messages:
      OrderPlaced:
        $ref: "#/components/messages/OrderPlaced"
operations:
  publishOrderPlaced:
    action: send
    channel:
      $ref: "#/channels/orderPlaced"
    summary: Emit an event when a new order is placed.
  consumeOrderPlaced:
    action: receive
    channel:
      $ref: "#/channels/orderPlaced"
    summary: Consume order-placed events for fulfilment.
components:
  messages:
    OrderPlaced:
      name: OrderPlaced
      title: Order Placed
      contentType: application/json
      payload:
        type: object
        required: [orderId, customerId, total, currency, placedAt]
        properties:
          orderId: { type: string, format: uuid }
          customerId: { type: string, format: uuid }
          total: { type: integer, minimum: 0 }
          currency: { type: string, pattern: "^[A-Z]{3}$" }
          placedAt: { type: string, format: date-time }
        examples:
          - orderId: "8a2e1b7c-9f4d-4e3a-b1c2-3d4e5f6a7b8c"
            customerId: "1b2c3d4e-5f6a-7b8c-9d0e-1f2a3b4c5d6e"
            total: 2999
            currency: USD
            placedAt: "2026-05-29T14:23:00Z"
```

### Gherkin behavior spec (BDD)

```gherkin
Feature: Apply discount codes at checkout
  As a shopper
  I want to apply a discount code
  So that I receive the advertised price reduction

  Background:
    Given a cart with one item of $50.00

  Scenario: Valid percentage discount applies
    Given a discount code "SUMMER10" worth 10%
    When I apply the code at checkout
    Then the cart total is "$45.00"
    And the discount line shows "SUMMER10 (-$5.00)"

  Scenario: Expired code is rejected
    Given a discount code "EXPIRED" that expired yesterday
    When I apply the code at checkout
    Then the cart total remains "$50.00"
    And an error "This code has expired" is shown

  Scenario Outline: Stacking rules
    Given a discount code "<code>" worth <percent>%
    And the cart already has code "FIRST5" applied
    When I apply the code at checkout
    Then the second code is <outcome>

    Examples:
      | code      | percent | outcome   |
      | STACK10   | 10      | rejected  |
      | STACKOK   | 5       | accepted  |
```

### Functional spec skeleton (Spolsky-style, prose)

For non-API behavior that doesn't fit OpenAPI or Gherkin (e.g., a CLI tool, a background job, a sync algorithm):

```markdown
# [Component] Functional Spec

**Status:** Draft / Approved · **Version:** 1.0 · **Owner:** [name]
**Last updated:** YYYY-MM-DD

## Overview
[1–2 paragraphs: what this thing does, who calls it, why it exists.]

## Inputs
| Name | Type | Required | Description | Example |
|---|---|---|---|---|
| ... | ... | ... | ... | ... |

## Outputs
| Name | Type | Description | Example |
|---|---|---|---|
| ... | ... | ... | ... |

## Behavior
- **Happy path:** [step-by-step from input to output]
- **Edge cases:**
  - When [condition], then [behavior]
  - When [condition], then [behavior]
- **Error cases:**
  - On [failure mode], emit [error code], do [action]

## Invariants
- [Thing that is always true after a successful call]
- [State guarantee — e.g., idempotency]

## Performance / SLOs
- p50 latency: [...]
- p99 latency: [...]
- Throughput: [...]

## Non-goals
- [What this spec does NOT cover]

## Versioning
- This spec: v1.0
- Compatibility: backward-compatible additions only; removals require v2.0
- Deprecation window: 12 months for major versions

## Open issues
- [Unresolved decision] — decide by: [date]
```

## Anti-patterns

1. **Implementation leakage** — spec mentions Redis, Postgres, language choice, or framework. Consumers may build on assumptions that break when implementation changes. *Fix:* describe only observable behavior.

2. **Architecture argument in the spec** — spec text reads like an RFC debating approaches. *Fix:* move tradeoff analysis to an RFC; spec captures the decision only.

3. **No examples** — schemas without concrete payloads. Consumers misinterpret types, formats, and edge cases. *Fix:* every endpoint, message, and scenario gets at least one realistic example.

4. **Hand-waved errors** — "returns an error on failure" without listing codes, shapes, or conditions. Clients can't handle them. *Fix:* explicit error response schemas with code list and trigger conditions.

5. **Missing versioning policy** — no statement of what counts as breaking, no deprecation window, no version field. *Fix:* publish a versioning policy in `info.description` or a top-level README.

6. **Gherkin abuse** — Gherkin scenarios that describe UI clicks (`When I click the blue button`) instead of behavior (`When I apply the discount code`). Couples spec to UI implementation. *Fix:* describe the user's intent and the system's response, not the mechanics.

7. **Spec-test drift** — spec says one thing, implementation does another, no CI check. *Fix:* contract tests that validate implementation against spec on every PR.

8. **Optional everything** — every field is optional, every status code is `default`. The spec is useless as a contract. *Fix:* be strict; mark required fields required; enumerate status codes explicitly.

9. **Living-spec-only-in-code** — spec exists only as Swagger annotations buried in handlers, and the YAML is generated. Spec is not reviewed independent of code; design-first impossible. *Fix:* check in the spec file; treat it as source of truth; generate stubs from it.

10. **One-spec-per-codebase** — single 5000-line OpenAPI file for 80 endpoints. Becomes unmaintainable. *Fix:* split with `$ref` into per-resource files; use a build step to bundle if tools demand a single file.

## Decision Heuristics

| Situation | Use |
|---|---|
| HTTP / REST API, public or internal | OpenAPI 3.1 |
| GraphQL API | GraphQL SDL (the schema IS the spec) |
| gRPC service | Protobuf `.proto` file |
| Kafka / AMQP / WebSocket events | AsyncAPI 3.0 |
| Webhook receivers | OpenAPI 3.1 (model the receiver as a callback) |
| Business logic with branching rules | Gherkin / Given-When-Then |
| Internal library / CLI / job | Functional spec (Spolsky-style prose) |
| Pure data shape | JSON Schema / Avro / Protobuf |

**Spec-file vs code-doc test:** does the contract cross a team or service boundary, AND will another team consume it independently? Yes → separate spec file. No → inline code doc.

**Contract-first vs code-first test:** is more than one team or developer building against this AND can they start in parallel? Yes → contract-first. No → code-first then capture spec.

**Prose vs Gherkin vs OpenAPI test:**
- Question is "what does the endpoint accept and return?" → OpenAPI/AsyncAPI
- Question is "what happens when condition X meets condition Y?" → Gherkin
- Question is "how does this background sync algorithm behave?" → prose functional spec

**Versioning scheme by audience:**
- External public API → URL or Date versioning (Stripe-style)
- Internal microservices → URL versioning + Schema Registry
- Event topics → topic-name versioning (`orders.v2`)

## Cross-references

- **Adjacent skills:** `api-design-patterns` (the design phase that produces the spec inputs), `api-docs-craft` (reference docs derived from specs), `rfc-and-design-docs` (architecture decisions that precede specs), `prd-writing` (product-level requirements that motivate specs), `coding-standards` (style for the implementation that conforms).
- **Pairs with:** `user-story-and-acceptance-criteria` — user stories often translate directly into Gherkin scenarios.
- **Upstream:** PRD → RFC → spec. Downstream: spec → implementation → contract tests.

## References

1. Joel Spolsky, "Painless Functional Specifications, Part 2: What's a Spec?" — the seminal what-not-how distinction. https://www.joelonsoftware.com/2000/10/03/painless-functional-specifications-part-2-whats-a-spec/
2. OpenAPI Initiative, "Best Practices" — official guidance for OpenAPI documents. https://learn.openapis.org/best-practices.html
3. OpenAPI Specification v3.1 — formal spec. https://spec.openapis.org/oas/v3.1.0
4. AsyncAPI Initiative, "AsyncAPI 3.0.0 Specification" — formal spec. https://www.asyncapi.com/docs/reference/specification/v3.0.0
5. AsyncAPI Initiative, "AsyncAPI 3.0.0 Release Notes" — changes from v2 (channel/operation decoupling, send/receive). https://www.asyncapi.com/blog/release-notes-3.0.0
6. Martin Fowler, "GivenWhenThen" — bliki entry by the framer of the BDD style. https://martinfowler.com/bliki/GivenWhenThen.html
7. Cucumber Project, "Gherkin Reference" — canonical syntax docs. https://cucumber.io/docs/gherkin/reference/

## When to Use This Skill

Use when the user asks to: write or critique an API spec (OpenAPI, AsyncAPI), draft Gherkin/BDD scenarios, write a functional spec for a CLI/library/job, define a data contract, decide between contract-first and code-first, set a versioning policy, or clarify the boundary between spec and RFC.

Skip when the user asks for: product-level requirements (use prd-writing), architecture tradeoff analysis (use rfc-and-design-docs), an engineering task plan (use code-plan-writing), software architecture decisions (use software-architect), API design decisions (REST vs GraphQL, resource modeling — use api-design-patterns), public-facing reference docs (use api-docs-craft), or standalone user stories (use user-story-and-acceptance-criteria).
