<!-- hub-reference-banner -->
> **Reference file — part of the `da-data-engineering-platform` hub.** Formerly the standalone `da-14-streaming-analytics` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-14-streaming-analytics
title: Streaming Analytics
version: "1.0.0"
updated: "2026-05-30"
category: data-analysis
origin: local
description: >
  Streaming analytics — continuous, low-latency processing of unbounded event
  streams. Covers the Kafka ecosystem (brokers, partitions, consumer groups,
  Kafka Streams, ksqlDB), Apache Flink (stateful stream processing,
  checkpoints, savepoints, DataStream vs Table API), Spark Structured
  Streaming (micro-batch vs continuous mode, Delta as streaming sink), Apache
  Beam (unified batch/stream model, runner abstraction), exactly-once
  semantics (idempotent producers, transactions, two-phase commit),
  windowing (tumbling, sliding, session, global), watermarks, allowed
  lateness, and real-time decisioning (feature stores, online inference,
  real-time A/B tests).
  TRIGGER: user asks about Kafka topology, partition strategy, consumer
  groups, Kafka Streams vs ksqlDB, Flink jobs / checkpoints / savepoints,
  Spark Structured Streaming, Apache Beam, exactly-once semantics, stream
  windowing, watermarks, late data handling, feature stores for online
  inference, real-time A/B tests, or any "process this event stream
  continuously" problem.
  SKIP: building MongoDB change-stream consumers (use mongodb-change-streams);
  MongoDB Kafka Connector source/sink configuration (use mongodb-kafka-connector);
  Atlas Stream Processing operator authoring (use mongodb-atlas-stream-processing);
  pure batch ETL (use da-3-data-acquisition-sampling); generic ML training
  (use da-7-machine-learning).
triggers:
  - streaming analytics
  - Kafka topology
  - Kafka partition
  - consumer group
  - Kafka Streams
  - ksqlDB
  - Apache Flink
  - Flink checkpoint
  - Flink savepoint
  - DataStream API
  - Table API
  - Spark Structured Streaming
  - micro-batch
  - continuous mode
  - Apache Beam
  - Beam runner
  - exactly-once
  - idempotent producer
  - tumbling window
  - sliding window
  - session window
  - watermark
  - allowed lateness
  - late data
  - feature store
  - online inference
  - real-time A/B test
tags:
  - data-analysis
  - streaming
  - kafka
  - flink
  - spark
  - beam
  - real-time
when_to_use:
  - Designing a Kafka topology (topics, partitions, replication, consumer groups)
  - Choosing Kafka Streams vs ksqlDB vs Flink vs Spark Structured Streaming vs Beam
  - Tuning Flink checkpoint interval, state backend, or savepoint strategy
  - Switching Spark Streaming between micro-batch and continuous mode
  - Picking the right Beam runner (Dataflow, Flink, Spark) for portability
  - Engineering exactly-once with idempotent producers and Kafka transactions
  - Selecting a windowing strategy (tumbling, sliding, session, global)
  - Setting watermarks and allowed lateness for late-arriving events
  - Productionizing a feature store for online inference (Feast, Tecton)
  - Running real-time A/B tests or contextual bandits on a live stream
when_not_to_use:
  - MongoDB change-stream consumer or resume-token mechanics — use mongodb-change-streams
  - MongoDB Kafka Connector source/sink config — use mongodb-kafka-connector
  - Atlas Stream Processing operator authoring — use mongodb-atlas-stream-processing
  - Pure batch ETL or one-time extraction — use da-3-data-acquisition-sampling
  - General ML training pipelines without streaming feature serving — use da-7-machine-learning
  - A/B test design that is not real-time / continuous — use da-12-ab-testing-causal-inference
related_skills:
  - da-3-data-acquisition-sampling
  - da-13-modern-data-stack
  - da-7-machine-learning
  - da-12-ab-testing-causal-inference
  - mongodb-change-streams
  - mongodb-kafka-connector
  - mongodb-atlas-stream-processing
---

# Streaming Analytics

Streaming analytics is the discipline of processing unbounded sequences of events
in (near-)real time, producing continuous answers instead of one-shot batch
results. Unlike batch analytics — where the input is a finite, frozen dataset —
streaming systems must reason about time, ordering, lateness, state, and
failure on an input that never ends. The shift from "compute once over yesterday's
data" to "compute continuously over the live wire" introduces a different set of
primitives: partitions, watermarks, windows, checkpoints, exactly-once.

This skill covers the modern streaming stack: the Kafka ecosystem as the
durable backbone; Flink, Spark Structured Streaming, and Apache Beam as the
processing engines; the semantics that distinguish them (delivery guarantees,
windowing, late-data handling); and the real-time decisioning layer that
turns continuous computation into product behavior (feature stores, online
inference, real-time experiments).

## Core mental model

A streaming system has four layers, and confusion almost always comes from
collapsing them:

1. **Transport** — durable, replayable log of events. Kafka is the de-facto
   standard; Pulsar, Kinesis, and Pub/Sub are the alternatives. The log is the
   source of truth.
2. **Processing** — stateful computation over the log: filtering, joining,
   aggregating, windowing. Kafka Streams, ksqlDB, Flink, Spark Structured
   Streaming, Beam.
3. **Serving / sink** — where results land for downstream use: a database, a
   feature store, a Delta/Iceberg lake, another Kafka topic, an inference
   service.
4. **Decisioning** — code that turns a continuous answer into action: a
   real-time model, a rule engine, an experiment assignment, a feature flag.

Every streaming question maps to one of these layers. "Why am I losing events?"
is almost always a transport-or-processing exactly-once question. "Why are my
windows wrong?" is a watermark-or-allowed-lateness question. "Why is my model
stale?" is a feature-store-freshness question.

## 1. The Kafka ecosystem

Apache Kafka is a distributed, replicated, append-only commit log. It is the
single most important piece of infrastructure in modern streaming because it
decouples producers from consumers in time (replay) and in space (independent
scaling) while providing strong durability and ordering guarantees per
partition.

### Topology

- **Broker** — a Kafka server process. A cluster is N brokers (typically
  3 / 5 / 7) running in a quorum. Since KRaft (KIP-500) replaced Zookeeper,
  brokers run their own Raft-based metadata quorum; this is now the default
  in all 3.x and 4.x deployments. Zookeeper mode was removed in Kafka 4.0.
- **Topic** — a named, partitioned, replicated log. Topics are append-only;
  retention is by time, size, or compaction.
- **Partition** — the unit of parallelism. Each partition is an ordered log
  on disk, replicated to N brokers (typically RF=3). Order is guaranteed
  *within* a partition; across partitions, order is undefined.
- **Producer** — writes records (key, value, headers, timestamp). The key
  determines the partition (`hash(key) mod numPartitions` by default), which
  is how you achieve per-entity ordering: same key → same partition →
  ordered.
- **Consumer group** — set of consumers that cooperatively read a topic.
  Each partition is assigned to exactly one consumer in the group, so adding
  consumers up to `numPartitions` scales out; beyond that, additional
  consumers idle.
- **Offset** — monotonic position of a consumer in a partition; committed
  back to Kafka (`__consumer_offsets`) so a restarting consumer resumes
  where it left off.

### Partition strategy is the design decision

Pick the partition count and key carefully — it's the lever you cannot easily
change later without resharding. Heuristics:

- Partitions = max parallelism for a single consumer group. Plan for 2-4×
  your current peak so you can scale without repartitioning.
- Too few → consumer-side bottleneck. Too many → controller overhead, longer
  rebalances, more open file handles per broker. Modern Kafka handles
  hundreds of thousands per cluster, but per-topic stay under a few thousand
  unless you have a specific reason.
- Key choice = ordering boundary. If you need per-user ordering, key by
  user_id. If you key by something low-cardinality (e.g., region), you'll
  get hot partitions.

### Kafka Streams vs ksqlDB

Both are stream-processing libraries that run *inside* the JVM and read /
write Kafka topics — no external cluster.

- **Kafka Streams** — Java/Scala library. You write a DSL or Processor API
  topology, package it as a normal JVM app, and run it anywhere. State is
  held in embedded RocksDB and changelog-replicated back to Kafka. Best when
  you want full control, custom logic, and tight integration with an
  existing JVM service.
- **ksqlDB** — SQL surface over Kafka Streams. You write `CREATE STREAM` /
  `CREATE TABLE` / `SELECT` statements and ksqlDB compiles them into Kafka
  Streams topologies that run on a ksqlDB server cluster. Best when the
  transformations are SQL-shaped (filter, project, join, window aggregate)
  and the team prefers SQL over Java.

Both are bounded to a single Kafka cluster as both source and sink — you
can't easily process a non-Kafka source. If you need cross-system processing
or higher-level state primitives, reach for Flink.

## 2. Apache Flink — stateful stream processing

Apache Flink is a distributed stream processor that treats batch as a special
case of streaming. It is the reference implementation of the Dataflow model
and is favored when low latency, sophisticated state, and exactly-once
semantics across heterogeneous sources/sinks matter more than ecosystem
breadth.

### Architecture

- **JobManager** — coordinates: schedules tasks, triggers checkpoints,
  recovers on failure. Runs in HA mode with a standby JobManager and a
  leader election service (ZK or K8s).
- **TaskManager** — worker process. Each TM has N task slots; each slot
  runs one parallel sub-task of an operator. Parallelism is configured per
  operator or per job.
- **State backend** — where operator state lives. Options:
  - `hashmap` (heap-based) — fast, bounded by TM heap, used for small state.
  - `rocksdb` — embedded LSM tree on local disk, used for large state
    (TB-scale). Slightly slower per access, but unbounded by heap.
- **Checkpoint storage** — durable storage for checkpoint snapshots
  (S3, GCS, HDFS, ABFS).

### Checkpoints and savepoints — the difference matters

Both are consistent snapshots of operator state plus source offsets, but they
serve different roles:

- **Checkpoint** — automatic, periodic, owned by the running job. Used for
  failure recovery. Retention is short (last N), and on job cancellation
  they're deleted by default. Algorithm: Chandy-Lamport-style asynchronous
  barriers flow through the dataflow; each operator snapshots its state when
  it receives the barrier on all input channels.
- **Savepoint** — manual, durable, owned by you. Used for planned operations:
  job upgrades, parallelism changes, A/B comparisons. A savepoint is a
  self-contained image you can restore a new (potentially different) job
  from. Always take a savepoint before upgrading Flink or your job DAG.

Tune `execution.checkpointing.interval` based on tolerable recovery cost,
not too aggressively — every checkpoint costs CPU and I/O. 10-60 s is
typical for production jobs; sub-second is exotic. Use
`execution.checkpointing.mode = EXACTLY_ONCE` (the default) unless you've
proved you don't need it.

### DataStream API vs Table API

- **DataStream API** — the imperative core. You write Java/Scala/Python that
  manipulates `DataStream<T>` objects with operators (`map`, `keyBy`,
  `window`, `process`, `connect`). Most expressive; gives you full access
  to state, timers, side outputs, and `ProcessFunction` for fine-grained
  control. Use when the logic is not easily expressed in SQL.
- **Table API / Flink SQL** — declarative. Define tables over streams
  (`CREATE TABLE ... WITH ('connector' = 'kafka', ...)`), write `SELECT
  ... GROUP BY TUMBLE(...)`. The planner compiles to DataStream. Use when
  the work is shaped like SQL — joins, aggregations, windowed analytics —
  and you want optimizer-level rewrites and operator fusion.

The two APIs are interoperable: convert a `DataStream` to a `Table` and back
inside a single job. Flink CDC connectors (the upstream of Debezium-style
change data capture in many architectures) are first-class Table API
sources.

## 3. Spark Structured Streaming

Spark Structured Streaming is the streaming layer of Apache Spark. It reuses
the DataFrame / Dataset / SQL API, so the same Spark code runs on bounded
batch and unbounded streams with minimal change. This is its primary
appeal: one engine, one mental model, one team skill set for both regimes.

### Micro-batch vs continuous processing

Two execution modes coexist:

- **Micro-batch (default)** — the engine repeatedly slices the input into
  small batches and runs them as ordinary Spark jobs. Latency is bounded by
  trigger interval, typically 100 ms - seconds. End-to-end latency is
  usually 1-5 s. This is what virtually every production Structured
  Streaming job runs.
- **Continuous processing (experimental since Spark 2.3)** — a record-at-a-
  time engine that drops latency to ~1 ms but supports a much smaller
  surface of operations (essentially map-like, no stateful aggregations).
  In practice, teams rarely choose continuous mode; if you need
  sub-100-ms latency, reach for Flink.

Triggers:

- `Trigger.ProcessingTime("30 seconds")` — fixed cadence.
- `Trigger.AvailableNow()` — process all currently-available data, then
  stop. Useful for streaming-style ingest scheduled by Airflow/cron — gives
  you idempotent restartable ingest with exactly-once semantics inherited
  from the engine.
- `Trigger.Continuous("1 second")` — continuous mode checkpoint interval.

### State, watermarks, output modes

Stateful operations (`groupBy(...).count()`, joins between streams, drop-
duplicates) require a watermark to bound state. Set it with
`.withWatermark("event_time", "10 minutes")` — state for windows older than
`max(event_time) - 10 min` is evicted.

Output modes:
- `append` — only new rows are written; required for sinks that don't
  support updates, supported only when state can be finalized (e.g., after
  watermark closes a window).
- `update` — only changed aggregation rows are written; supported by
  Kafka, Delta, Iceberg.
- `complete` — full result table is rewritten every trigger; only viable
  for small aggregations.

### Delta Lake as the canonical streaming sink

Delta Lake (and to a lesser extent Iceberg / Hudi) is the de-facto streaming
sink for Structured Streaming. Delta provides:

- ACID transactions on object storage, so concurrent batch + streaming
  writers don't corrupt the table.
- Streaming both *from* and *to* a Delta table — a Delta table can be the
  source of another stream, enabling the "medallion" (bronze / silver / gold)
  pattern of progressively-refined streaming pipelines.
- `MERGE INTO` for streaming upserts (CDC-style change data application).
- Change Data Feed (CDF), so downstream streams can subscribe to row-level
  changes on a Delta table.

Pair Structured Streaming + Delta when your organization already runs Spark
and wants storage, batch, and streaming on one substrate. Choose Flink + a
plain Kafka topic when latency and per-record state matter more than the
unified analytics surface.

## 4. Apache Beam — the unified model

Apache Beam is not a runtime; it's a portable programming model with
pluggable runners. You write your pipeline once against the Beam SDK
(Java, Python, Go) and run it on Google Cloud Dataflow, Apache Flink,
Apache Spark, or the direct (local) runner.

### Why Beam exists

Beam was originally Google's "FlumeJava + MillWheel + Cloud Dataflow"
contribution. The thesis: batch and streaming are the same problem viewed
through different time-completeness lenses, and a single SDK should target
many engines.

The model is the Dataflow model (Akidau et al., 2015): every record has an
event time, every operation is conceptually applied within a window, and
correctness is parameterized by *what* is being computed, *where* in event
time, *when* in processing time, and *how* refinements relate.

### Core primitives

- **PCollection** — a (potentially-unbounded) collection of elements.
- **PTransform** — a transformation: `ParDo` (per-element), `GroupByKey`,
  `Combine`, `Window`, `Flatten`, side inputs.
- **Pipeline** — a DAG of PTransforms.
- **Runner** — the execution engine. The same pipeline runs on Dataflow
  (managed GCP), Flink, Spark, Samza, or DirectRunner (local testing).

### When to pick Beam

- You want portability across engines or clouds.
- You're on GCP and Dataflow is the path of least resistance for managed
  streaming (autoscaling, Streaming Engine, Dataflow Prime).
- You need a single SDK across batch + streaming, but don't want to commit
  to Spark or Flink's surface specifically.

When to skip:
- Engine-specific features (Flink savepoint introspection, Spark's
  full DataFrame API, ksqlDB SQL) are unavailable through Beam's lowest-
  common-denominator surface.
- For Kafka-only, Kafka-in-Kafka-out work, Kafka Streams or ksqlDB is
  simpler.

## 5. Exactly-once semantics

Three delivery guarantees, in increasing strength:

1. **At-most-once** — messages may be lost, never duplicated. Cheap; rarely
   what you want.
2. **At-least-once** — messages are never lost, may be duplicated. The
   default for naive retry loops. Tolerable if downstream is idempotent.
3. **Exactly-once** — each message has exactly one effect, end-to-end.
   Requires coordination between producer, broker, processor, and sink.

### Exactly-once in Kafka

Two building blocks, introduced in Kafka 0.11:

- **Idempotent producer** (`enable.idempotence=true`, the default since
  Kafka 3.0). The producer assigns a producer ID and sequence number to
  every record. The broker deduplicates retries within a producer session
  on a single partition.
- **Transactions** (`transactional.id`). A producer can batch writes to
  multiple partitions inside a transaction; the broker commits or aborts
  atomically. Consumers with `isolation.level=read_committed` see only
  committed records. The producer ID survives restarts via the
  `transactional.id`, enabling fencing — the broker rejects writes from a
  stale producer instance.

### Exactly-once across processor + Kafka

Kafka Streams sets `processing.guarantee=exactly_once_v2` and the runtime
handles transaction boundaries: it reads a batch of records, updates state
in changelog topics, writes outputs, and commits offsets — all in one
Kafka transaction. Failure means abort + retry from the last committed
position.

Flink uses **two-phase commit (2PC)** sinks for end-to-end exactly-once.
Phase 1 (pre-commit) happens at checkpoint barrier; the sink writes data
in a pending state and persists transaction state in the checkpoint.
Phase 2 (commit) happens after the checkpoint is acknowledged globally.
On failure, the job restores the checkpoint and re-runs the pending
commits. The KafkaSink and FileSink ship with built-in 2PC support;
custom sinks implement `TwoPhaseCommitSinkFunction`.

Spark Structured Streaming achieves exactly-once via idempotent sinks
(Delta, Iceberg, idempotent Kafka) plus checkpointed source offsets — the
engine guarantees each input offset is processed exactly once *into the
sink*, with the sink responsible for de-duping on retry. This is weaker
than Flink's 2PC for arbitrary sinks: if your sink can't be idempotent
(e.g., a REST POST that triggers a side effect), exactly-once is not
achievable in Structured Streaming without external dedup.

### Real-world caveats

- "Exactly-once" means *within the system*. A non-transactional external
  side effect (an HTTP POST, an email send) is always at-least-once
  unless the receiver is idempotent.
- Producer transactions have throughput cost: ~5-15% in practice. Measure
  before assuming you need them — many pipelines are fine with
  at-least-once + idempotent downstream writes.

