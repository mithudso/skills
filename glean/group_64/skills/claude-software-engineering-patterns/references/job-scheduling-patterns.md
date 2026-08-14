<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `job-scheduling-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: job-scheduling-patterns
description: >
  Node.js job scheduling expert: cron expression syntax, in-process scheduling (node-cron,
  node-schedule), durable job queues (BullMQ, Agenda, SQS), distributed locks (MongoDB
  findAndModify, Redis Redlock), retry strategies (fixed/linear/exponential/jitter), dead
  letter queues, idempotent job design (natural keys, deduplication windows), graceful
  shutdown with in-flight drain, job lifecycle management, and monitoring/alerting.
  TRIGGER: user asks about cron jobs, scheduled tasks, or background job processing in
  Node.js; implementing or debugging node-cron, Agenda, or BullMQ; preventing overlapping
  job runs; distributed locks for jobs; retry strategies or DLQs; idempotent job design;
  choosing between BullMQ/Agenda/SQS/RabbitMQ; graceful shutdown with workers.
  SKIP: Chrome extension alarm scheduling (use alarm-scheduler-patterns); Python/Go/Ruby
  schedulers with no Node.js component; pure infrastructure scheduling (AWS EventBridge
  Scheduler, Kubernetes CronJob) with no application-level job logic.
version: 1.2.0
updated: "2026-05-29"
category: developer
tags:
  - node.js
  - cron
  - job-queue
  - bullmq
  - agenda
  - node-cron
  - redis
  - mongodb
  - distributed-systems
  - idempotency
  - retry
  - worker
  - graceful-shutdown
  - monitoring
  - scheduling
keywords:
  - node-cron
  - BullMQ
  - Agenda
  - node-schedule
  - cron expression
  - distributed lock
  - Redlock
  - dead letter queue
  - DLQ
  - idempotency key
  - exponential backoff
  - jitter
  - graceful shutdown
  - in-flight jobs
  - job lifecycle
  - stalled jobs
  - Bull Board
  - job deduplication
  - MongoDB job lock
  - Redis job queue
whenToUse:
  - User asks about cron jobs, scheduled tasks, or background job processing in Node.js
  - Implementing or debugging node-cron, node-schedule, Agenda, or BullMQ
  - Preventing overlapping job runs or implementing distributed locks
  - Designing retry strategies, dead letter queues, or failure recovery
  - Making jobs idempotent or handling deduplication
  - Graceful shutdown handling with in-flight workers
  - Choosing between job queue libraries (BullMQ vs Agenda vs SQS vs RabbitMQ)
  - Monitoring job queues or alerting on stuck/stalled jobs
whenNotToUse:
  - Chrome extension alarm scheduling — use alarm-scheduler-patterns
  - Python, Go, or Ruby schedulers with no Node.js component
  - Pure AWS infrastructure scheduling (EventBridge Scheduler, Step Functions) — use aws-serverless
related_skills:
  - backend-patterns
  - alarm-scheduler-patterns
  - mongodb-expert
  - ops-registry-patterns
  - javascript-nodejs
globs:
  - "**/jobs/**/*.js"
  - "**/cron/**/*.js"
  - "**/workers/**/*.js"
  - "**/queue/**/*.js"
  - "**/scheduler*.js"
  - "**/agenda.js"
  - "**/bull*.js"
---

# Job Scheduling Patterns

Node.js job scheduling spans two layers:

1. **In-process cron** (node-cron, node-schedule) — callbacks at wall-clock times within a single process. Zero persistence; dies with the process.
2. **Durable job queues** (BullMQ, Agenda, SQS) — jobs survive restarts, support retries, distributed workers, and observability dashboards.

The `server/src/jobs/` pattern in this codebase uses **node-cron** with MongoDB-persisted output. See the "This Repo's Pattern" section for details and graduation criteria.

## Library Selection

| Library       | Backend    | Persistence | Distributed | Best For |
|---------------|------------|-------------|-------------|----------|
| **node-cron** | none       | none        | no          | Simple in-process tasks, dev tooling |
| **node-schedule** | none   | none        | no          | Time/date-based one-off scheduling |
| **Agenda**    | MongoDB    | yes         | yes (lock)  | Apps on MongoDB needing durable jobs |
| **BullMQ**    | Redis      | yes         | yes         | High-throughput, priority queues, fan-out |
| **pg-boss**   | PostgreSQL | yes         | yes         | Postgres-native, one fewer dependency |
| **AWS SQS**   | AWS        | yes         | yes         | Cloud-native, managed, infinite scale |
| **RabbitMQ**  | AMQP       | yes         | yes         | Complex routing, polyglot, enterprise |

**Decision shortcuts:**
- MongoDB already in use + durable scheduled jobs → **Agenda**
- MongoDB already in use + simple fire-and-report (this repo's approach) → **node-cron**
- Redis already in use → **BullMQ**
- AWS-managed infrastructure → **SQS**
- Complex routing / polyglot → **RabbitMQ**
- Single-process, no restart survival needed → **node-cron**

## Cron Expression Quick Reference

```
┌──────── minute (0–59)
│ ┌────── hour (0–23)
│ │ ┌──── day of month (1–31)
│ │ │ ┌── month (1–12)
│ │ │ │ ┌ day of week (0–7, 0=Sunday)
* * * * *
```

| Expression    | Meaning               |
|---------------|-----------------------|
| `*/5 * * * *` | Every 5 minutes       |
| `15 5 * * *`  | 5:15 AM daily         |
| `0 9 * * 1-5` | 9 AM Mon–Fri          |
| `0 */6 * * *` | Every 6 hours         |
| `0 0 1 * *`   | Midnight on 1st/month |

**Timezone:** `cron.schedule('15 5 * * *', handler, { timezone: 'America/New_York' })` — test explicitly around DST transitions; use UTC for unambiguous schedules.

## Retry / Backoff Strategy Selection

| Strategy | When to Use | Formula |
|---|---|---|
| Fixed | Rate-limited APIs with known reset | `delay` |
| Exponential | Network errors, DB connection failures | `2^(attempt-1) * delay` |
| Exponential + full jitter | Thundering herd prevention (preferred) | `random(0, 2^(attempt-1) * delay)` |
| Custom (retry-after header) | Domain-specific | `function(attempts, err) → ms` |

**BullMQ exponential:** `backoff: { type: 'exponential', delay: 2000 }` — attempts 1→2s, 2→4s, 3→8s, etc.

**Note:** Custom backoff strategies are registered on the **Worker**, not the Queue.

## Idempotency Design

A job is idempotent if running it N times produces the same final state as running it once. Required for any at-least-once delivery system.

| Type | Example | Implementation |
|---|---|---|
| Natural | Setting a status field to a fixed value | `updateOne({ $set: { status: 'done' } })` — safe to repeat |
| Synthetic | Sending an email, charging a payment | Idempotency key + deduplication store |

**Key principle:** use a unique index on the idempotency key collection + catch `err.code === 11000` to detect duplicates.

**BullMQ deduplication:** pass a deterministic `jobId` — BullMQ skips enqueue if a job with that ID already exists in active/waiting state.

## Distributed Delivery Semantics

| Guarantee | What It Means | How to Achieve |
|---|---|---|
| At-most-once | May be dropped; never duplicated | Fire-and-forget, no retry |
| At-least-once | Will be delivered, may be duplicated | Retry + idempotent handlers |
| Exactly-once | Delivered exactly once | Not achievable without idempotency |

**Practical rule:** design all handlers as idempotent (at-least-once) and use deduplication keys to short-circuit duplicates.

## BullMQ Job Lifecycle

```
waiting → active → completed
        ↘ failed → (retry) → waiting
                 ↘ (exhausted) → failed (permanent)
        → delayed → waiting
        → stalled → waiting  (lock expired, worker crashed)
```

**Auto-cleanup:** `removeOnComplete: { count: 200, age: 7 * 24 * 3600 }` and `removeOnFail: { count: 500 }`.

## UI Monitoring Dashboards

| Tool | Libraries | Notes |
|---|---|---|
| **Bull Board** | BullMQ | Express middleware; inspect/retry/remove jobs |
| **Arena** | BullMQ, Bull | Separate server; interactive UI |
| **Taskforce.sh** | BullMQ | Hosted; real-time; built by BullMQ authors |

## This Repo's Pattern (node-cron + MongoDB runner)

```
registry.js            → maps name → report module; effectiveSchedule() applies env overrides
cron.js                → schedules via node-cron; REPORT_SCHEDULE_<NAME> env override
runner.js              → wraps run(), persists output to MongoDB, creates run_id
reports/daily-digest.js → run({ snapshot, logger, anthropic }) → { output, metadata }
```

**Key design choices:** `DISABLE_CRON=1` disables all scheduling; `REPORT_TIMEZONE` sets system TZ; `cron.validate()` guards bad expressions at startup; `stopCron()` supports `--watch` hot-reload.

**Graduate from node-cron to BullMQ/Agenda when:**
- Jobs regularly take > 30 seconds (need timeout + stall detection)
- Need job completion tracking across restarts
- Running multiple server instances (need distributed lock or queue)
- Need UI visibility into job history
- Need to enqueue jobs from external triggers (API, webhooks)

## Quick Reference Snippets

```js
// node-cron: validate + schedule
cron.validate('15 5 * * *')
cron.schedule('15 5 * * *', handler, { timezone: 'America/New_York' })

// BullMQ: add with retry + timeout + cleanup
queue.add('name', data, {
  attempts: 3, timeout: 5 * 60 * 1000,
  backoff: { type: 'exponential', delay: 1000 },
  removeOnComplete: { count: 100 }, removeOnFail: { count: 500 },
})

// BullMQ: graceful shutdown
await worker.close()  // stops new pickup, drains in-flight

// MongoDB distributed lock (requires unique index on jobName)
db.collection('job_locks').findOneAndUpdate(
  { jobName, $or: [{ lockedAt: null }, { lockedAt: { $lt: staleThreshold } }] },
  { $set: { lockedAt: new Date(), lockedBy: myId } },
  { upsert: true, returnDocument: 'after' }
)
// Then: if (!result || result.lockedBy !== myId) return;

// Idempotent natural-key upsert
db.collection('reports').updateOne({ reportId }, { $set: data }, { upsert: true })

// Synthetic idempotency key
try { await db.collection('idempotency_keys').insertOne({ key, ... }); }
catch (e) { if (e.code === 11000) return { skipped: true }; throw e; }
```

## Full Reference

For complete code examples — BullMQ queue+worker setup, Agenda configuration, Redis Redlock, DLQ pattern, in-flight drain shutdown, Prometheus alerting, Bull Board setup — read `references/job-scheduling-context.md`.
