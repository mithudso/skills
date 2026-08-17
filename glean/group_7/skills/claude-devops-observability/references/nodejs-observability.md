<!-- hub-reference-banner -->
> **Reference file — part of the `devops-observability` hub.** Formerly the standalone `nodejs-observability` skill.
> Sibling topics in this family are now reference files under the hubs (`devops-linux-internals`, `devops-linux-admin`, `devops-containers-cicd`, `devops-observability`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: nodejs-observability
description: >
  Node.js observability expert: OpenTelemetry setup (zero-code auto-instrumentation, programmatic
  SDK, ESM support), distributed tracing (W3C Trace Context, manual spans, sampling, baggage),
  Prometheus metrics (prom-client, OTel Metrics API), health checks (/healthz /readyz /startupz),
  graceful shutdown, error tracking (Sentry), APM tools (Datadog/New Relic/Elastic/Grafana),
  Node.js runtime metrics (event loop lag/utilization, GC, memory leak detection), log-trace
  correlation, AsyncLocalStorage request context, and SLI/SLO alerting (RED/USE methods,
  error-budget burn-rate rules).
  TRIGGER: user is implementing or reviewing observability in a Node.js application; setting
  up OpenTelemetry, distributed tracing, Prometheus metrics, health check endpoints, Sentry,
  APM tools, event loop monitoring, graceful shutdown, log-trace correlation, SLO alerting,
  or error budget burn rate in Node.js.
  SKIP: observability for non-Node.js runtimes (Python, Go, Java); pure infrastructure
  monitoring (Kubernetes cluster metrics, cloud provider health dashboards) with no
  application-level instrumentation; Pino logger configuration details (use pino-structured-logging).
version: 1.2.0
updated: "2026-05-29"
category: developer
tags:
  - nodejs
  - observability
  - opentelemetry
  - tracing
  - metrics
  - apm
  - monitoring
  - prometheus
  - sentry
  - health-checks
  - sli
  - slo
keywords:
  - opentelemetry
  - otel
  - distributed tracing
  - spans
  - W3C trace context
  - prometheus
  - prom-client
  - health check
  - liveness
  - readiness
  - sentry
  - datadog
  - new relic
  - elastic apm
  - grafana tempo
  - loki
  - event loop lag
  - event loop utilization
  - GC metrics
  - memory leak
  - structured logging
  - trace correlation
  - SLI
  - SLO
  - RED method
  - USE method
  - error budget
  - burn rate
  - auto-instrumentation
  - graceful shutdown
  - AsyncLocalStorage
whenToUse:
  - Setting up OpenTelemetry tracing or metrics in a Node.js application
  - Adding /healthz, /readyz, or /startupz health check endpoints
  - Integrating Sentry, Datadog, New Relic, Elastic APM, or the Grafana OSS stack
  - Monitoring event loop lag, GC pauses, or memory growth in production
  - Implementing graceful shutdown with connection draining
  - Correlating log lines with distributed trace IDs
  - Setting up SLO-based alerting (error budget burn rate, RED method, USE method)
  - Choosing an APM vendor for a Node.js service
  - Debugging missing traces, missing spans, or metrics not appearing
whenNotToUse:
  - Observability for Python, Go, Java, or other non-Node.js runtimes
  - Pure infrastructure monitoring (Kubernetes cluster metrics, cloud dashboards) with no app instrumentation
  - Pino logger setup or configuration details — use pino-structured-logging
  - Generic Express middleware questions with no observability angle
related_skills:
  - pino-structured-logging
  - microservices-patterns
  - backend-patterns
  - sentry-monitoring
  - job-scheduling-patterns
triggers:
  - opentelemetry node
  - otel node setup
  - distributed tracing node
  - prometheus node metrics
  - event loop lag
  - health check readiness liveness
  - sentry node error tracking
  - apm node comparison
  - structured logging trace correlation
  - SLI SLO RED USE method
  - node observability
  - grafana tempo loki node
globs:
  - "**/instrumentation.{js,ts,mjs}"
  - "**/tracing.{js,ts,mjs}"
  - "**/metrics.{js,ts,mjs}"
  - "**/health*.{js,ts,mjs}"
  - "**/observability/**"
  - "**/telemetry/**"
  - "**/otel-config*.{yaml,yml}"
---

# Node.js Observability

Observability in Node.js spans three pillars — traces, metrics, and logs — plus health checks
and alerting. OpenTelemetry (OTel) is the 2025–2026 baseline standard. Instrument once with
OTel; export to any backend.

## Navigation

| I want to... | Start here |
|---|---|
| Add tracing + metrics with minimal code | [OpenTelemetry Setup](#quick-start) — zero-code path |
| Create custom spans for business logic | `references/nodejs-observability-context.md` — Manual Spans |
| Expose /metrics for Prometheus | `references/nodejs-observability-context.md` — Prometheus Client |
| Add /healthz and /readyz endpoints | `references/nodejs-observability-context.md` — Health Checks |
| Track errors in production | `references/nodejs-observability-context.md` — Sentry Integration |
| Choose an APM vendor | [APM Comparison](#apm-tools-comparison) below |
| Monitor event loop lag, GC, memory | `references/nodejs-observability-context.md` — Node.js Runtime Metrics |
| Link log lines to traces | `references/nodejs-observability-context.md` — Log-Trace Correlation |
| Set up SLOs and alerts | [SLI/SLO section](#slislo-alerting) below |
| Debug missing traces or metrics | [Troubleshooting](#troubleshooting) below |

## Quick Start

```bash
# Fastest path: zero-code OTel with OTLP export
npm install @opentelemetry/auto-instrumentations-node
OTEL_SERVICE_NAME=my-svc OTEL_EXPORTER_OTLP_ENDPOINT=http://collector:4318 \
  node --require @opentelemetry/auto-instrumentations-node/register app.js

# Prometheus metrics endpoint
npm install prom-client
# then: promClient.collectDefaultMetrics(); app.get('/metrics', ...)

# Health probes + graceful shutdown
npm install @godaddy/terminus
```

**ESM note:** add `--experimental-loader=@opentelemetry/instrumentation/hook.mjs` before `--require` for ESM applications.

## APM Tools Comparison

| Feature | Datadog | New Relic | Elastic APM | Grafana Stack |
|---|---|---|---|---|
| OTel native | Yes | Yes | Yes | Yes (Tempo) |
| Node.js agent | dd-trace | newrelic | elastic-apm-node | N/A (use OTel) |
| Auto-instrumentation | Excellent | Excellent | Good | Via OTel |
| Pricing | Per host + ingestion | Per-user + ingestion (100GB free) | Self-host free | Self-host free |
| Best for | Large orgs, multi-cloud | App-centric teams | Existing ELK users | Cost-sensitive, OSS-first |

**Quick decision:**
- Budget-constrained/OSS-first → Grafana stack (Tempo + Prometheus + Loki)
- Need turnkey deep Node.js APM → New Relic (best onboarding, 100GB free tier)
- Heavy infra + microservices → Datadog (best service discovery, 850+ integrations)
- Already running ELK → Elastic APM (native log correlation, self-hostable)

## Default Metrics from prom-client

| Metric | Type | Description |
|---|---|---|
| `nodejs_eventloop_lag_seconds` | Gauge | Event loop lag |
| `nodejs_eventloop_lag_p99_seconds` | Gauge | p99 event loop lag |
| `nodejs_gc_duration_seconds` | Histogram | GC pause duration by type |
| `nodejs_heap_size_used_bytes` | Gauge | V8 used heap size |
| `nodejs_version_info` | Info | Node.js version |

## SLI/SLO Alerting

### RED Method (request-oriented services)

| Signal | Metric | PromQL example |
|---|---|---|
| Rate | Requests/sec | `rate(http_requests_total[5m])` |
| Errors | Failed requests/sec | `rate(http_requests_total{status=~"5.."}[5m])` |
| Duration | Latency p99 | `histogram_quantile(0.99, sum(rate(http_request_duration_ms_bucket[5m])) by (le))` |

### USE Method (Node.js runtime resources)

| Signal | Node.js metric |
|---|---|
| Utilization | Event loop utilization (ELU) |
| Saturation | Event loop lag p99, active handles |
| Errors | GC failures, EMFILE errors |

### Recommended SLIs

| SLI | Target | Measurement |
|---|---|---|
| Availability | 99.9% | Successful responses / total |
| Latency p99 | < 500ms | `histogram_quantile(0.99, ...)` |
| Error rate | < 0.1% | 5xx / total |
| Event loop utilization | < 70% sustained | `perf_hooks.eventLoopUtilization()` |

## Anti-Patterns

1. **Fat liveness probes** — never query databases in `/healthz`; Kubernetes will restart pods on transient DB issues.
2. **Unbounded cardinality** — user IDs or UUIDs as metric labels OOM the collector; use bounded values.
3. **Sampling nothing** — 100% trace sampling in production is costly; use `ParentBasedSampler` at 1–10%.
4. **Sampling everything away** — 0.1% misses rare critical errors; use tail-based sampling to always capture error traces.
5. **Logging without trace context** — logs without `traceId` cannot be correlated; inject OTel context into your logger.
6. **Ignoring event loop** — CPU % misses the signal; monitor ELU before throughput drops.
7. **Alerting on raw error counts** — alert on error-budget burn rate over sliding windows instead.
8. **Missing graceful shutdown** — without SIGTERM handling, in-flight requests drop during deploys.
9. **Heap snapshots in production traffic** — `v8.writeHeapSnapshot()` blocks the event loop; only use on debug instances.
10. **No resource attributes** — traces without `service.name`, `service.version`, `deployment.environment` are unidentifiable.

## Troubleshooting

| Problem | Likely cause | Fix |
|---|---|---|
| No traces appearing | OTel SDK not loaded before app code | Move `--require instrumentation.js` before app entry |
| Missing library spans | Library not auto-instrumented | Add manual spans or find a community instrumentation package |
| Trace context lost across async | CLS context not propagating | Verify `@opentelemetry/context-async-hooks` registered (default in NodeSDK) |
| Metrics endpoint returns empty | `collectDefaultMetrics()` not called | Call at startup; verify `/metrics` route registered |
| High memory in collector | Unbounded metric cardinality | Audit label values; drop high-cardinality labels in collector config |
| No log-trace correlation | Logger not patched | Use OTel instrumentation for your logger or inject manually |
| OTLP export failures | Collector unreachable | Check `OTEL_EXPORTER_OTLP_ENDPOINT`, verify collector running, check grpc vs http |
| ESM auto-instrumentation fails | Missing loader hook | Add `--experimental-loader=@opentelemetry/instrumentation/hook.mjs` |

## Full Reference

For all code implementations — programmatic OTel SDK setup, manual spans, sampling strategies, baggage, Prometheus client, health check endpoints, Sentry integration, graceful shutdown, event loop/GC/memory metrics, log-trace correlation, AsyncLocalStorage, SLO alerting rules, Grafana/Tempo/Loki stack, OTel Collector config — read `references/nodejs-observability-context.md`.
