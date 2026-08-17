<!-- hub-reference-banner -->
> **Reference file — part of the `ai-agents-orchestration` hub.** Formerly the standalone `agent-state-and-durable-execution` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

<!-- FOLDED REFERENCE under the ai-agent-engineering hub — agent state management & durable execution (the orchestration runtime). Sibling refs: autonomous-loops (loop patterns), agent-harness-construction (action/observation contracts), multi-agent-orchestration (topologies), agent-planning-patterns (control flow). Any "use the X" pointer refers to a sibling reference under this hub. -->

# Agent State & Durable Execution

An agent runtime needs two separable things: a way to **manage state** (the data it reasons over) and a way to **survive failure** over long horizons. Conflating them is the common mistake.

## 1. State management (the in-process layer)
**LangGraph** is the dominant answer: an explicit `StateGraph` with typed **channels** and **reducers** (how each update merges into state). The graph runs in **super-steps**; a pluggable **checkpointer** snapshots state after each super-step, keyed by **`thread_id`**. That one mechanism powers conversation memory, resume, human-in-the-loop, and time-travel.
- State vs scratchpad vs long-term memory (memory architecture is the `agent-memory-architecture` reference; this is the *execution* state).

## 2. Human-in-the-loop & time-travel
From the checkpointer you get: **`interrupt()` / resume** (approval gates, edit-state-and-continue), and **time-travel** — replay from or **fork** at any past checkpoint (Git-like history of the run). Essential for approvals, debugging, and "what-if" re-runs.

## 3. Durable execution (the crash-survival layer)
**Checkpointing ≠ durable execution.** A checkpointer saves *between* super-steps only — no in-node recovery, no cross-process exactly-once coordination. For long-running, side-effecting agents you need a durable-execution engine that **journals every step and replays after a crash**, so completed steps (and their paid, non-idempotent side effects) are **not re-run**:
- **Temporal** (activities + workflow replay), **Restate** (journaled durable handlers), **DBOS** (Postgres-backed, in-process library), **Inngest** (`step.run` memoization, serverless).
- **The determinism constraint:** replay requires the workflow body be deterministic, so **all nondeterminism (LLM calls, tool I/O) must be quarantined into recorded steps/activities**. Exactly-once holds only for registered resources — add **idempotency keys** for arbitrary external APIs.

## 4. Architecture spectrum & the "use both" pattern
In-process library (LangGraph checkpointer; DBOS — Postgres only) → serverless step engine (Inngest, Restate) → external orchestration server (Temporal).
- **Minutes-long, low-failure, single-process →** a checkpointer suffices.
- **Multi-hour, side-effecting, distributed →** a durable backend, often the **"use both"** pattern: **Temporal for the macro lifecycle, LangGraph for the micro reasoning loop**.
- **LangGraph `durability` modes** trade performance for recovery granularity: `"exit"` (least) → `"async"` → `"sync"` (most durable). Stream intermediate state/tokens with the event/message stream; support `cancel`/interrupt for long runs.

## Decision guide
- Need conversation memory / resume / approvals / debugging replay → **LangGraph state + checkpointer** (pick `async`/`sync` durability per cost tolerance).
- Agent runs for hours/days, calls paid or non-idempotent tools, must survive process crashes → add a **durable-execution engine**; push every LLM/tool call into a recorded step and add idempotency keys.
- Both (production long-running agent that also needs a rich reasoning loop) → **durable engine wrapping LangGraph**.

## Anti-patterns
Treating a checkpointer as crash-safety for side-effecting agents; nondeterministic code inside a durable workflow body (breaks replay); assuming exactly-once for un-registered external APIs (needs idempotency keys); a full durable backend for a 30-second single-process agent (over-engineering).

## References
LangGraph persistence, durable-execution (`exit`/`async`/`sync`), and time-travel docs (docs.langchain.com) · Temporal "durable execution meets AI" · Restate durable agents · DBOS durable execution for AI agents · Inngest durable workflows · Diagrid, "Checkpoints are not durable execution" · vendor-neutral determinism explainers (Vanlightly). Verify fast-moving API signatures against live official docs.
