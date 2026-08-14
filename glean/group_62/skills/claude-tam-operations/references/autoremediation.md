<!-- hub-reference-banner -->
> **Reference file — part of the `tam-operations` hub.** Formerly the standalone `autoremediation` skill.
> Sibling topics in this family are now reference files under the hubs (`tam-operations`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: autoremediation
description: "Design and implement self-healing systems: retry/circuit-breaker patterns, automated error recovery, graceful degradation, and AI-assisted code repair agents. TRIGGER: retry strategy design, circuit breaker state machine, fallback chain, canary rollback, LLM-driven repair loop, SLO burn-rate trigger, chaos engineering blast-radius, error classification tier, idempotency key, retry budget, runbook automation. SKIP: general exception handling with no recovery strategy, pure logging or telemetry setup, manual incident response without automation, load testing without failure injection."
version: "1.1.2"
updated: "2026-06-01"
category: developer
tags: [autoremediation, circuit-breaker, retry, self-healing, chaos-engineering, fallback, resilience, observability]
related_skills: [software-engineering-patterns, technical-writing-craft]
---

# Autoremediation and Self-Healing Systems

## When to Use This Skill

| Task | Examples |
|---|---|
| Retry logic design | Backoff strategy, jitter type, retry budget, idempotency guarantees |
| Circuit breaker | State machine, threshold tuning, half-open probe, fallback wiring |
| Graceful degradation | Fallback chain, feature-flag-gated degradation, SLA tiering |
| Automated rollback | Deployment health gates, canary abort triggers, migration reversals |
| Error classification | Auto-fixable vs escalate vs hard stop |
| Agent-based remediation | LLM-driven repair loops, plan-execute-verify cycles, safe code patching |
| Chaos engineering | Steady-state hypothesis, blast-radius control, observability coupling |
| Observability-driven response | Alert-to-action pipelines, SLO burn-rate triggers, runbook automation |

## When NOT to Use This Skill

- General exception handling or try/catch structuring (no recovery strategy needed)
- Logging or telemetry setup not tied to a remediation action
- Manual incident response (no automation component)
- Load testing or performance benchmarking without failure injection

---

## Execution Sequence

When asked to design or implement any autoremediation pattern, follow this sequence:

1. **Classify the failure** — use the Error Classification table below to determine the tier (auto-fix / escalate / hard stop)
2. **Select the pattern** — use the Decision Table to pick the right recovery strategy
3. **Configure thresholds** — apply the key config values from the Decision Table; for any value marked `[ASSUMED]`, state what workload-specific data (request rate, SLO target, dependency SLA) would change it
4. **Wire the fallback chain** — ensure fallbacks are independent, tested under load, and capped at ≤3 hops; verify each hop is idempotent before enabling async retry
5. **Define the exit condition** — specify when normal operation resumes (reset triggers, health check criteria)
6. **Validate** — confirm the pattern handles the failure class, does not introduce cascading risk, and has observable state

---

## Error Classification Tiers

| Error Class | Examples | Action |
|---|---|---|
| Transient / retriable | Network timeout, 429, 503, connection reset | Auto-retry with backoff + jitter |
| Degraded dependency | Slow DB, flaky upstream, partial outage | Circuit breaker + fallback |
| Bad input / client error | 400, 422, validation failure | No retry — return error immediately |
| Deployment regression | p99 spike, error rate jump post-deploy | Automated canary rollback |
| Known bug class | SAST-detected CVE, recurring null deref | AI-assisted patch with human gate |
| Unknown / novel | Unexpected panic, data corruption | Page on-call / open incident — do not auto-retry or auto-patch |

---

## Quick Reference: Decision Table

| Scenario | Recommended Pattern | Key Config |
|---|---|---|
| Transient network error | Exponential backoff + full jitter | max_attempts=5, base=100ms, cap=30s |
| Downstream service slow | Circuit breaker (count-based) | threshold=5 failures / 10s window |
| Downstream service flaky | Circuit breaker (rate-based) | >50% failure rate over 20-call window |
| Read from degraded DB | Fallback to read replica, then cache | Fallback chain length ≤ 3 |
| Write to unavailable service | Queue + async retry with DLQ | TTL-capped; all writes must carry an idempotency key |
| Bad deployment | Automated canary rollback | p99 latency +20% or error rate >1% |
| Known bug class (SAST) | AI-assisted code repair with human gate | Patch + test + human approval before merge |
| Cascading failure risk | Bulkhead + shed load | Separate thread pools per dependency |
| Recurring alert | Runbook automation / auto-remediation pipeline | Validate runbook via chaos test before enabling automation |
| High retry volume | Retry budget (token bucket) | Budget = 10% of total RPS; refill rate = 1 token/s per caller |

---

## Idempotency Requirement

Any operation placed on an async retry queue **must** be idempotent. Before enabling queue-based retry:

- Writes: include a client-generated `idempotency_key` (UUID or content hash) that the server deduplicates
- Reads: safe to retry without a key, but verify the read does not have side effects (e.g., audit logging that must not double-write)
- Mutations that cannot be made idempotent: do not queue — fail fast and escalate

---

## Required Output Format

When this skill produces a design or implementation recommendation, include:

- **Pattern chosen** and the rationale (one sentence)
- **Configuration values** with justification for threshold choices; mark workload-dependent values with `[ASSUMED]`
- **Fallback chain** described step-by-step, with exit/reset condition
- **Observable state** — what metric or signal confirms recovery
- **Idempotency guarantee** — state how each retried operation is made safe to replay
- Code or pseudo-code for the core retry/circuit-breaker/rollback logic when the request implies implementation

---

## Complementary Skills

- `software-engineering-patterns` — broader distributed-systems context (service mesh, saga, CQRS, queue/DLQ, idempotency, async job patterns)
- `ai-agent-engineering` — when the remediation loop involves LLM agents with tool use
- `incident-response` → `tam-operations` (references/incident-response/SKILL.md) — manual incident handling and on-call escalation once an unknown/novel failure is paged
