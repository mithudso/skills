<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `microservices-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: microservices-patterns
description: >
  Microservices architecture expert: service decomposition (bounded contexts, DDD, strangler fig),
  inter-service communication (REST/gRPC/async messaging), message brokers (Kafka/RabbitMQ/SQS/NATS),
  service discovery, API gateway patterns (BFF, aggregation), event-driven architecture (event
  sourcing, CQRS), saga pattern (choreography vs orchestration), service mesh (Istio/Linkerd),
  distributed tracing (OpenTelemetry), database-per-service, eventual consistency, and deployment
  patterns (sidecar, ambassador).
  TRIGGER: user is designing, building, or migrating a microservices system; asking about
  service decomposition, inter-service communication, message broker selection, distributed
  transactions, event-driven design, breaking up a monolith, Kafka vs RabbitMQ, saga pattern,
  CQRS, service mesh, or how services should talk to each other — even without the word "microservices".
  SKIP: single-service internal implementation with no inter-service concerns — use backend-patterns;
  early-stage/small-team products where a monolith is correct; AWS-specific Lambda/EventBridge/Step
  Functions questions — use aws-serverless; purely infrastructure questions (Kubernetes cluster
  config, cloud networking) with no service-communication angle.
version: "1.1.0"
updated: "2026-05-29"
category: developer
tags:
  - microservices
  - event-driven
  - kafka
  - rabbitmq
  - saga
  - cqrs
  - event-sourcing
  - service-mesh
  - distributed-systems
  - opentelemetry
  - api-gateway
  - ddd
keywords:
  - microservices
  - bounded context
  - DDD
  - strangler fig
  - saga
  - CQRS
  - event sourcing
  - Kafka
  - RabbitMQ
  - SQS
  - NATS
  - gRPC
  - service mesh
  - Istio
  - Linkerd
  - BFF
  - distributed tracing
  - OpenTelemetry
  - database per service
  - eventual consistency
  - circuit breaker
  - sidecar
  - ambassador
  - monolith decomposition
whenToUse:
  - Designing or reviewing a microservices architecture from scratch
  - Choosing between REST, gRPC, and async messaging for inter-service communication
  - Selecting a message broker (Kafka vs RabbitMQ vs SQS vs NATS)
  - Implementing a saga pattern for distributed transactions
  - Applying CQRS or event sourcing to a service
  - Breaking up a monolith into services (strangler fig migration)
  - Configuring a service mesh (Istio, Linkerd) or API gateway
  - Planning observability for distributed services (distributed tracing, structured logs)
  - Asking how services should talk to each other, even without saying "microservices"
whenNotToUse:
  - Single-service internal implementation with no inter-service concerns — use backend-patterns
  - Early-stage or small-team (< 5 engineers) products where a monolith is correct
  - AWS-specific serverless event-driven questions (Lambda/EventBridge/Step Functions) — use aws-serverless
  - Purely infrastructure questions (Kubernetes cluster config, cloud networking) with no service-communication angle
related_skills:
  - backend-patterns
  - aws-serverless
  - nodejs-observability
  - job-scheduling-patterns
  - kubernetes-networking
---

# Microservices Architecture Patterns

Comprehensive reference for designing, building, and operating microservices architectures.
Covers decomposition strategies, communication patterns, data management, resilience,
observability, and deployment.

*Last researched: May 2026. 25+ cited sources in `references/microservices-patterns-context.md`.*

## When NOT to use this skill

Recommend simpler approaches when:

- Team is < 5 engineers or product is early-stage — microservices overhead overwhelms thin teams
- Domain is not yet understood; a monolith lets you refactor without cross-service contracts
- Question is purely about single-service internals — use `backend-patterns` instead
- No Kubernetes/container orchestration exists and the team has no ops capacity for it
- Question is about the *consensus algorithm* itself — Paxos/Raft, quorum intersection, linearizability vs eventual consistency, or the CAP/FLP theory beneath eventual-consistency and saga tradeoffs — use `distributed-systems-consensus`

## How to Respond

**Before recommending a pattern**, establish context if not already known:
- Team size and deployment environment (Kubernetes, managed cloud, bare metal)?
- Greenfield design or migration from an existing system?
- Primary constraints: latency, throughput, operational simplicity, cost?

Ask at most 2–3 targeted questions. If the user says "just answer", use available context and note assumptions.

**Output format by question type:**

| Question type | Output shape |
|---|---|
| "Which pattern should I use?" | Decision rationale + table row from Quick Decision Tables + 1–2 trade-off bullets |
| "How do I design X?" | Read context file first; numbered steps + key decision points |
| "Compare A vs B" | Side-by-side table (Criterion / A / B) + recommendation with stated assumptions |
| "Help me migrate / decompose" | Phased plan (identify boundaries → extract → route → cut over) + strangler-fig note |
| "What's wrong with my design?" | Flag matching anti-patterns; explain why each is a problem; suggest alternative |

For trade-off questions: state the key tension explicitly before recommending. Example: "you need replay, so Kafka; but your team is on AWS and wants zero ops, so consider MSK or pull toward SQS for simpler queues that don't need replay."

For depth questions ("how do I implement X?", "what are the trade-offs of Y?"): read `references/microservices-patterns-context.md` first and cite sources when making non-obvious recommendations.

## Quick Decision Tables

### Communication Pattern Selection

| Scenario | Recommended Pattern |
|---|---|
| Public external API | REST (OpenAPI) |
| High-performance internal RPC | gRPC |
| Fire-and-forget notifications | Async messaging (any broker) |
| Workflow / multi-step processes | Orchestrated saga (Temporal) |
| Domain events to many consumers | Event bus (Kafka or NATS pub/sub) |
| Streaming data pipeline | Kafka |
| Request requires immediate response | Synchronous (REST or gRPC) |

### Message Broker Selection

| Need | Choose |
|---|---|
| High-throughput event streaming + replay | Kafka |
| Flexible routing / task queues | RabbitMQ |
| AWS-native, zero-ops | SQS/SNS |
| Ultra-low latency service mesh messaging | NATS |
| Ordered log / audit trail | Kafka (log compaction) |

## Core Concepts

| Concept | Summary |
|---|---|
| Bounded Context | Each service owns one business domain; anti-corruption layers protect domain models at boundaries |
| Saga | Sequence of local transactions with compensating actions; choreography for simple flows, orchestration for complex ones |
| CQRS | Separate read (query) and write (command) models; enables independent scaling and specialized read stores |
| Event Sourcing | Persist all state changes as immutable events; replay to reconstruct state; pairs with CQRS |
| Service Mesh | Infrastructure layer handling mTLS, load balancing, retries, circuit breaking via sidecar proxies |
| BFF | Dedicated backend per client type (web, mobile, third-party); aggregates downstream services into client-optimized responses |

## Anti-Patterns to Flag

- Shared database across services (breaks service autonomy)
- Synchronous chains of 4+ service calls (latency amplification, cascading failures)
- Anemic event payloads (forces consumers to call back for state — defeats async decoupling)
- God gateway with business logic (should be routing and cross-cutting only)
- Sagas without compensating transactions defined upfront
- Missing distributed trace context propagation across async message boundaries
