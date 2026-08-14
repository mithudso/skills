<!-- hub-reference-banner -->
> **Reference file — part of the `aws-cloud` hub.** Formerly the standalone `aws-serverless` skill.
> Sibling topics in this family are now reference files under the hubs (`aws-cloud`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: aws-serverless
description: >
  AWS serverless, event-driven, and container orchestration expert. Covers Lambda (cold starts,
  concurrency, SnapStart, Managed Instances), EventBridge (buses, pipes, scheduler, schema
  registry), Step Functions (standard vs express, Distributed Map), API Gateway (REST vs HTTP
  vs WebSocket), ECS Fargate, EKS (Auto Mode, Karpenter, managed node groups), CDK, SAM, and
  CloudFormation. Delivers decision tables, annotated code, or structured comparisons.
  TRIGGER: user asks about Lambda, EventBridge, Step Functions, ECS, EKS, Fargate, CDK, SAM,
  API Gateway, serverless patterns (saga, CQRS, fan-out, idempotency), cold starts, DLQs,
  container orchestration, or IaC tool selection for AWS.
  SKIP: pure EC2/VPC/S3/RDS/IAM questions with no serverless component — use aws-core; AWS
  ML/AI services (Bedrock, SageMaker) — use aws-ai-ml; non-AWS serverless (Cloudflare Workers,
  Vercel, GCP Cloud Run); Redshift/Athena/Glue analytics not paired with a Lambda/EventBridge trigger.
version: "1.2.0"
updated: "2026-05-29"
category: developer
tags:
  - aws
  - serverless
  - lambda
  - eventbridge
  - step-functions
  - ecs
  - eks
  - fargate
  - cdk
  - sam
  - cloudformation
  - api-gateway
  - event-driven
  - containers
keywords:
  - lambda
  - eventbridge
  - step functions
  - ecs fargate
  - eks
  - karpenter
  - cdk
  - sam
  - cloudformation
  - api gateway
  - cold start
  - provisioned concurrency
  - snapstart
  - dlq dead letter queue
  - saga pattern
  - fan-out
  - cqrs
  - idempotency
  - event sourcing
  - distributed map
  - lambda managed instances
  - eks auto mode
whenToUse:
  - Designing or reviewing Lambda functions, triggers, or concurrency settings
  - Choosing between ECS Fargate, EKS, and Lambda for a workload
  - Selecting between CloudFormation, CDK, and SAM for IaC
  - Implementing serverless patterns: fan-out, saga, CQRS, event sourcing, idempotency
  - Troubleshooting cold starts, concurrency limits, DLQs, or dead-letter failures
  - EKS cluster design, node group vs Fargate tradeoffs, Karpenter or Auto Mode configuration
  - API Gateway design: REST vs HTTP vs WebSocket, usage plans, JWT auth
  - EventBridge rule/pipe/scheduler design
  - CDK or CloudFormation template validation
whenNotToUse:
  - Pure EC2, VPC, S3, RDS, or IAM questions with no serverless or event-driven component — use aws-core
  - AWS AI/ML: Bedrock, SageMaker, foundation models — use aws-ai-ml
  - Non-AWS serverless (Cloudflare Workers, Vercel Functions, GCP Cloud Run)
  - Redshift, Athena, Glue, or analytics services not explicitly paired with a Lambda/EventBridge trigger
related_skills:
  - aws-core
  - aws-ai-ml
  - software-engineering-patterns
---

# AWS Serverless

Expert guidance on AWS serverless, event-driven, and container orchestration services.
Produces decision tables, annotated code snippets, or structured comparisons.

## How to Use This Skill

1. For comparison, design, or troubleshooting questions: read `references/aws-serverless-context.md` first for full decision tables and pattern detail. If the file is missing, use the inline shortcuts below and ask 1–2 targeted questions to narrow scope (e.g., "What runtime? What traffic pattern?").
2. For routing questions ("what AWS service handles X?"): the inline shortcuts below are sufficient.
3. Apply the relevant decision table for service selection. Present as a table when comparing 3+ options. If no table row matches, ask 2–3 targeted questions (team size, traffic pattern, existing infrastructure) before recommending.
4. For design questions ("how do I design X?"): produce numbered steps + key decision points. For IaC questions: use the CDK vs SAM vs CFN table, factoring in team profile (devs → CDK, ops → raw CFN, Lambda-only → SAM).
5. For Lambda cold start questions: check SnapStart runtime support (Java 11+, Python 3.12+, .NET 8+) before recommending it. If the runtime is unsupported, fall back to Provisioned Concurrency.
6. For EKS system add-ons (CoreDNS, kube-proxy, VPC CNI): recommend managed node groups over Fargate — Fargate pods don't exist until a workload is scheduled, which breaks add-on initialization.
7. For CDK/CloudFormation validation: use the `deploy-on-aws` MCP tools when available: `search_cdk_documentation`, `validate_cloudformation_template`, `check_cloudformation_template_compliance`.

**Output rule:** produce the most actionable artifact — a decision table, annotated snippet, or numbered design steps. For trade-off questions, state the key tension explicitly before recommending (example: "you need replay so EventBridge Archive fits, but if your team wants zero ops on AWS, SQS is simpler for queues that don't need replay").

> **Note:** CDKTF was deprecated by HashiCorp in 2025 and is not a valid escape hatch between CDK and Terraform.

## Quick Decision Tables

### Compute

| Workload profile | Recommended service |
|---|---|
| Event-driven, short bursts, variable traffic | Lambda |
| Long-running containers, predictable sustained load | ECS Fargate |
| Full K8s ecosystem, multi-team clusters, complex scheduling | EKS with Auto Mode |
| EKS with manual node autoscaling or custom NodePool logic | EKS + self-managed Karpenter |
| Cost-optimized batch with spot, GPU, or packed bin | ECS EC2 |
| Lambda with EC2 pricing model + multi-concurrency | Lambda Managed Instances (re:Invent 2025) |

### IaC Tool

| Situation | Tool |
|---|---|
| Multi-cloud required today or within 18 months | Terraform |
| Lambda-heavy serverless, small team, quick start | SAM |
| Complex AWS-native stack, developer team (TypeScript/Python) | CDK |
| Ops/SRE team, AWS-only, moderate complexity | Raw CloudFormation |

### API Gateway Type

| Need | Type |
|---|---|
| API keys, usage plans, per-client rate limiting, request validation | REST API |
| Low-latency, JWT/Cognito auth, Lambda proxy (up to ~70% cheaper than REST) | HTTP API |
| Bidirectional real-time, stateful connections | WebSocket API |

### Step Functions Workflow Type

| Need | Type |
|---|---|
| Long-running (up to 1 year), durable, exactly-once, auditable | Standard |
| High-volume, short-lived, at-least-once, cost-sensitive | Express |
| Large-scale parallel over S3 data (up to 10,000 concurrent child executions) | Distributed Map (Standard only) |

### EventBridge Component

| Need | Component |
|---|---|
| Route events from many producers to many consumers | Event Bus + Rules |
| Point-to-point source → target with filter/enrich/transform (no code) | EventBridge Pipes |
| Trigger actions on a cron or one-time schedule (millions of schedules) | EventBridge Scheduler |
| Discover and share event schemas across teams | Schema Registry |
| Replay or archive events for debugging | Archive + Replay |

## Core Reference

Read `references/aws-serverless-context.md` for:

- Cold start mitigation: SnapStart, Provisioned Concurrency, Lambda Managed Instances
- EventBridge component guide with configuration examples
- Step Functions patterns: saga, Distributed Map, error handling
- EKS compute options: Auto Mode, managed node groups, Fargate, Karpenter
- Serverless patterns: fan-out, CQRS, event sourcing, idempotency, DLQ design
- re:Invent 2025 announcements: Lambda Managed Instances, Lambda Durable Functions, EKS Auto Mode

If the file is not found, proceed with the inline shortcuts above and note the missing context file.
