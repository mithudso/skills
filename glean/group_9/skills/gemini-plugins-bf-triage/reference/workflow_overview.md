# BF / BFG Workflow

Condensed for the `bf-triage` skill. Defines the terminology, the
lifecycle a BF moves through, and the generic team-level expectations
the skill encodes when it recommends next steps.

Team-specific routing rules, priority order, and pattern recognition
live in [`team_knowledge.md`](team_knowledge.md) and are read in by
the skill at decision points. The default content shipped is for the
**Workload Resilience** team; other teams should clear that file's
content and replace it with their own.

## Core terminology

- **BFG** (Build Failure Genesis) — a single failure instance on a specific
  commit / variant / task in Evergreen. BFGs are tracked in the Build Baron UI
  and live in the BF database, not as Jira tickets.
- **BF** (Build Failure) — a unique defect, tracked as a Jira ticket in the
  `BF` project. A BF can be linked to one or many BFGs (the "Count of Linked
  BFGs" field shows the count over the past 30 days).
- **AF** (Atlas Failure) — a failure observed in Atlas production clusters,
  created from log-ingestion rules (fatal asserts, tripwire asserts, invariant
  failures, crashes). Limited diagnostics; logs expire.
- **ARR** (Auto-Resolution Rule) — a Build Baron rule that auto-attributes new
  BFGs to an existing BF based on test name regex or fault search terms.
- **Block-on-Red** — master-branch lockdown: only BF-fix commits allowed when
  the team or global BF threshold is hit.
- **Rapid-response** — top-priority cross-team BFs, P0, daily comment update
  required.

## Lifecycle (mainline / stable branches)

```text
Evergreen task fails ─► BFG auto-created
        │
        ▼
ARR matches an existing BF? ── yes ──► BFG attached to that BF (closed)
        │ no
        ▼
Build Baron team triages the BFG manually:
  - Find or create BF
  - Set Bug Symptoms, Severity Type(s), Severity Level, Assigned Teams
  - Add an ARR for future deduplication
        │
        ▼
BF in "Needs Triage" or "Open" with Assigned Teams set
        │
        ▼
Owning team's BF Triager does a 15-minute sanity check:
  - Cold that should be Hot? Escalate
  - Wrong team? Re-assign
  - Related to a known BF? Comment to help the eventual assignee
        │
        ▼
Team assigns an engineer:
  - One Hot BF or AF max per engineer
  - Cold BFs handled in sprint planning unless threshold is breached
        │
        ▼
Engineer investigates, comments diagnosis, links a SERVER fix ticket,
moves BF to "Waiting for Bug Fix"
        │
        ▼
Fix lands (and any backports). When ALL deps are closed, BF can be closed.
```

## Team triage priority

The owning team's triager picks the next BF to work according to a
priority queue. The default queue (Workload Resilience) is documented
in [`team_knowledge.md`](team_knowledge.md) §
"Triage priority queue". Other teams should rewrite that section with
their own ordering — common variations are weight on `block-on-red`
thresholds, AF prioritisation, and per-component sub-queues.

## What a good triage report has

Per `bf_triage_debugging_guide.md`:

1. Failure summary — what failed, on which variant / task / test
2. Severity classification — which type and why (cite log evidence)
3. Frequency & scope — how often, which branches / variants, BFG count
4. Root-cause hypothesis — based on log analysis, git history, pattern matching
5. Relevant commits — git log around the failure window, suspicious changes
6. Similar / duplicate BFs — by symptom text or shared task/variant
7. Environmental factors — AMI changes, machine issues, infrastructure
   anomalies
8. Recommended next steps — reproduce, revert, re-assign team, accept-as-known

Severity-type classification matrix lives in
[`severity_types.md`](severity_types.md). Log magic strings and
per-symptom decision trees live in [`log_patterns.md`](log_patterns.md).

## When a revert is warranted

Per the Revert Policy, the BF Triager (or any engineer) is empowered to revert
any commit that caused a failure **beyond a reasonable doubt**. The skill
should recommend "revert" only when:

1. Past executions of the same task / variant show the failure starts
   precisely at one commit; AND
2. The commit's title / diff plausibly relates to the failing test.

Otherwise the recommendation should be "reproduce" or "re-assign team", not
"revert".

## Generic re-routing decision matrix

The matrix below is **owner-team-agnostic** — it routes by failing
subsystem alone. Team-specific extensions (front-line patterns the
current team owns rather than re-routes away) live in
[`team_knowledge.md`](team_knowledge.md) § "Team front-line routing
rules".

| Symptom in logs | Likely owner |
| --------------- | ------------ |
| Replication / rollback / oplog crash | Replication |
| Sharding metadata, range deleter, balancer | Sharding |
| Query plan correctness, `$setWindowFields`, plan cache | Query |
| WiredTiger eviction / cache, dbpath corruption | Storage Engines |
| TLS / network handshake, TCP backlog | Networking & Observability |
| EC2 host vanished, AMI mismatch, package upgrade | DevProd Build (Runtime Environments) |
| Evergreen scheduler / task runner | DevProd Services & Integrations |

For each row the skill's "Recommended next steps" should pick the
listed team explicitly when re-routing is recommended, and explain the
routing decision with cited log evidence.

When the symptom in the BF does **not** match any row above — for
example a workload-runner error class, an admission-control source-code
change, an infra-tuning regression, or a workload-config schema mismatch
— consult [`team_knowledge.md`](team_knowledge.md) for the current
team's front-line ownership rules and keep-and-link patterns before
defaulting to re-route.

## Where new findings should land (Mode C learning)

When a verification run (Mode C) discovers a new pattern, classify it
before promoting to skill content:

- **Shared file** (this file, `log_patterns.md`, `severity_types.md`,
  `SKILL.md`) — patterns that hold true for **any** owning team:
  log markers, severity-type identifiers, generic workflow rules,
  tool / CLI integration tips.
- **Team file** ([`team_knowledge.md`](team_knowledge.md)) — patterns
  that name a team, encode "who owns this symptom front-line", or
  cite a team-specific downstream project (e.g. `TUNE-`, `PERF-`).

See [`verification_mode.md`](verification_mode.md) §
"Promoting findings to skill content" for the Mode C workflow that
applies this classification.
