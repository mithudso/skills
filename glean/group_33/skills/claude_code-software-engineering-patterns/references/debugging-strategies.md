<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `debugging-strategies` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: debugging-strategies
version: 2.0.1
updated: "2026-05-29"
description: >
  Master-level debugging reference covering systematic methodology, root cause analysis (5 Whys, Fishbone, Fault Tree), profiling tools by language/runtime, AI-assisted debugging workflows, distributed system observability with OpenTelemetry, and advanced techniques (binary search, differential debugging, memory leak detection, flamegraphs). Use when investigating bugs, performance issues, memory leaks, production incidents, distributed system failures, or running postmortem root cause analysis.
  TRIGGER: "how do I debug", "debugging strategy", "root cause analysis", "postmortem", "flamegraph", "memory leak", "performance regression", "distributed tracing", "OpenTelemetry", "git bisect", "profiling", "production bug", "intermittent bug", "5 whys", "fishbone diagram", "AI debugging workflow".
  SKIP: specific single-bug investigation with known reproduction steps (use debugging); performance optimization after root cause is known (use performance-profiling-expert); production incident coordination requiring runbooks (use incident-response).
category: developer
tags: [debugging, root-cause-analysis, profiling, observability, distributed-systems, ai-debugging, flamegraph, opentelemetry]
related_skills:
  - debugging
  - performance-profiling-expert
  - nodejs-observability
  - testing-and-vitest-expert
  - incident-response
whenToUse:
  - choosing the right debugging technique for a symptom
  - running a structured postmortem root cause analysis
  - profiling performance issues (flamegraphs, CPU, memory)
  - debugging distributed systems with traces and logs
  - using AI coding agents for autonomous debugging
  - investigating memory leaks with heap snapshots
  - finding regressions with git bisect
whenNotToUse:
  - specific bug with known repro steps (use debugging for the 7-phase workflow)
  - post-profiling optimization work (use performance-profiling-expert)
  - production incident coordination requiring cross-team runbooks (use incident-response)
---

# Debugging Strategies

Transform debugging from guesswork into systematic problem-solving with proven strategies, tools, and methodologies.

## Decision Tree: Pick the Right Technique First

```
Symptom observed
  |
  +-- Reproducible? --------> YES --> Minimal repro exists?
  |                                    |
  |                                    +-- YES --> Debugger / breakpoints
  |                                    +-- NO  --> Binary search / comment-out
  |
  +-- Intermittent? --------> Structured logging + stress test
  |                           Check for race conditions, timing deps
  |
  +-- Performance regression? --> Profile first, never guess:
  |                               Frontend: Chrome DevTools Performance
  |                               Node.js: --inspect + flamegraph (0x, Clinic.js)
  |                               Python: cProfile / py-spy
  |                               Go: pprof
  |
  +-- Memory growth? -------> Heap snapshots over time, compare deltas
  |
  +-- Production-only? -----> Gather observability evidence first
  |                           (logs, traces, metrics) before local repro
  |
  +-- Distributed / cross-service? --> Trace-based debugging (OpenTelemetry)
  |
  +-- Known regression? ----> git bisect
  |
  +-- Unfamiliar codebase? -> AI-assisted exploration
```

---

## Core Principles

### The Scientific Method

1. **Observe** — what is the actual behavior?
2. **Hypothesize** — what could be causing it?
3. **Experiment** — test the hypothesis
4. **Analyze** — did it prove or disprove the theory?
5. **Repeat** — until you reach root cause

This maps to the four-phase debugging methodology: **See it** (reproduce reliably) → **Shrink it** (isolate) → **Understand it** (root-cause) → **Kill it forever** (fix + regression test).

### Mindset

- "It can't be X" — yes it can. Check it.
- "I didn't change Y" — check anyway.
- "It works on my machine" — find out why.
- Reproduce consistently before forming hypotheses.
- Take breaks when stuck — fresh perspective finds what exhaustion misses.

---

## Root Cause Analysis Frameworks

### 5 Whys

Start with a problem statement; ask "Why?" five times in sequence. Reaches root cause without special tooling.

```
Problem: Deploy failed at 2:14 AM
  Why? --> Health check timed out
  Why? --> App took 45s to start instead of 10s
  Why? --> A new migration ran at startup
  Why? --> Migration was added to boot sequence, not a pre-deploy hook
  Why? --> No team convention for migration placement
Root cause: Missing convention → Fix: Add migration hook + runbook entry
```

**Limitation:** misses parallel root causes. Pair with Fishbone when a problem has multiple simultaneous causes.

### Fishbone Diagram (Ishikawa)

Brainstorm causes across categories — People, Process, Technology, Environment, Data, Dependencies — then converge on likely contributors. The standard format for incident postmortems in software.

### Fault Tree Analysis (FTA)

Model how specific failures combine via Boolean logic (AND/OR gates) to produce a top-level failure event. Use for safety-critical systems or when you need to calculate recurrence probability.

**When to use which:**
- Most postmortems: Fishbone + 5 Whys (fast, sufficient)
- Safety-critical / regulatory: Fishbone to brainstorm → FTA to model critical paths formally

---

## Systematic Debugging Process

### Phase 1: Reproduce

- Can you reproduce it? Always? Sometimes? Under specific conditions?
- Create a minimal reproduction — simplest example, unrelated code removed
- Document exact steps, environment details, and error messages

### Phase 2: Gather Information

- Full stack trace and error codes
- OS, language/runtime version, dependency versions, environment variables
- Git history, deployment timeline, config changes
- Scope: all users or specific? All browsers or specific? Production only?

### Phase 3: Form Hypothesis

- What changed recently (code, dependencies, infrastructure)?
- What's different between working and broken environments?
- Where in the stack could this fail (input validation, business logic, data layer, external services)?

### Phase 4: Test and Verify

- **Binary search** — comment out half the code, narrow to the failing section
- **Logging** — trace variable values and execution flow at key decision points
- **Isolate components** — test each piece separately with mocked dependencies
- **Differential comparison** — diff configurations, environments, and data between working and broken cases

---

## Debugging Tools by Language

### JavaScript / TypeScript

```typescript
// Pause execution
debugger;

// Timing
console.time("operation"); /* code */ console.timeEnd("operation");

// Stack trace at any point
console.trace();

// Assertion (throws if false)
console.assert(value > 0, "Value must be positive");

// Performance marks
performance.mark("start"); /* code */ performance.mark("end");
performance.measure("label", "start", "end");
```

VS Code `launch.json` for Node.js:
```json
{
  "type": "node", "request": "launch", "name": "Debug",
  "program": "${workspaceFolder}/src/index.ts",
  "preLaunchTask": "tsc: build - tsconfig.json",
  "skipFiles": ["<node_internals>/**"]
}
```

### Python

```python
breakpoint()          # Python 3.7+ (replaces pdb.set_trace())

import pdb
pdb.post_mortem()     # Debug at the exception point after a failure

# Logging for debugging
import logging
logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)
```

### Go

```go
import "runtime/debug"
debug.PrintStack()    // Print current stack trace

// Memory / CPU profiling via net/http/pprof
// Visit http://localhost:6060/debug/pprof/
```

---

## Advanced Techniques

### Technique 1: Binary Search (git bisect)

```bash
git bisect start
git bisect bad                    # current commit is broken
git bisect good v1.0.0            # last known good version
# Git checks out the midpoint — test it, then:
git bisect good   # if it works
git bisect bad    # if it doesn't
git bisect reset  # when done
```

### Technique 2: Differential Debugging

Build a comparison table between working and broken:

| Aspect | Working | Broken |
|--------|---------|--------|
| Environment | Dev | Production |
| Node version | 18.16.0 | 18.15.0 |
| Data | Empty DB | 1M records |
| User | Admin | Regular |
| Time | Daytime | After midnight |

Spot the difference — form a hypothesis from it.

### Technique 3: Memory Leak Detection

The key technique: compare heap snapshots over time. Objects that grow between snapshots but are never collected point to the leak.

```typescript
// Node.js: connect via --inspect flag, open chrome://inspect
// Programmatic monitoring:
if (process.memoryUsage().heapUsed > 500 * 1024 * 1024) {
  console.warn("High memory:", process.memoryUsage());
  require("v8").writeHeapSnapshot();
}

// In tests — track per-test memory growth:
beforeEach(() => { global.gc?.(); beforeMem = process.memoryUsage().heapUsed; });
afterEach(() => {
  global.gc?.();
  const diff = process.memoryUsage().heapUsed - beforeMem;
  if (diff > 10 * 1024 * 1024) console.warn(`Possible leak: ${diff / 1024 / 1024}MB`);
});
```

### Technique 4: Flamegraphs

Each box is a function. Width = total CPU time. Height = call stack depth. Find the widest boxes — those are the hotspots.

```bash
# Node.js
npx 0x -- node src/index.js          # opens interactive flamegraph
npx clinic flame -- node src/index.js

# Python
py-spy record -o profile.svg -- python app.py

# Go
go tool pprof -http=:8080 cpu.prof
```

---

## Debugging Patterns by Issue Type

### Intermittent Bugs

1. Add structured logging — timing, state transitions, external interactions
2. Check for race conditions — concurrent access to shared state, async ops completing out of order
3. Check timing dependencies — setTimeout/setInterval, Promise resolution order
4. Stress test — run many times, vary timing, simulate load

### Performance Issues

1. **Profile first** — never optimize blindly; measure before and after
2. Common culprits: N+1 queries, unnecessary re-renders, large computation on main thread, synchronous I/O blocking event loop
3. Tool selection: Browser → Chrome DevTools Performance; Node.js → --inspect + 0x or Clinic.js; Python → cProfile + py-spy; Go → pprof

### Production Bugs

1. Gather evidence first — error tracking (Sentry), application logs, metrics, user reports
2. Reproduce locally using production data (anonymized) with matched environment
3. Safe investigation — don't change production; use feature flags; test fixes in staging

---

## AI-Assisted Debugging

### When to use AI agents

Claude Code's agentic loop: **Observe** (read file tree, git history, test results) → **Reason** (form a plan) → **Act** (write files, run commands) → **Verify** (check for errors, iterate if tests fail).

**Use AI agents for:**
- Explaining unfamiliar error messages and stack traces
- Searching large codebases for where a value is set or mutated
- Generating hypotheses and test cases
- Automating repetitive investigation (log parsing, bisecting)

**Fall back to manual debugging for:**
- Timing-sensitive race conditions requiring precise breakpoint placement
- Hardware or OS-level issues outside the agent's tool reach
- Security-sensitive investigations where you control exactly what runs

### Effective AI debugging prompts

Think in criteria, not descriptions:

- **Weak:** "Fix the login bug"
- **Strong:** "The `/api/auth/login` endpoint returns 500 when email contains a `+` character. The fix is correct when: (a) login succeeds for `user+tag@example.com`, (b) existing tests still pass, (c) a new test covers the `+` character case."

Let the agent iterate — intervene only if it goes in circles (3+ failed attempts on the same approach). Always verify the diff yourself before accepting.

---

## Distributed System Debugging

Traditional log-based debugging is insufficient for distributed systems. Trace-based observability reduces mean time to resolution from hours to minutes.

### OpenTelemetry (OTel)

The industry standard, supported by every major observability vendor (Datadog, New Relic, Grafana, Honeycomb, Dynatrace, Splunk).

| Pillar | Shows | Use for |
|--------|-------|---------|
| **Traces** | Request journey through services (spans) | Find which service is slow or failing |
| **Metrics** | Aggregated measurements (latency, error rate) | Detect anomalies, set alerts |
| **Logs** | Discrete events with context | Understand what happened at a specific moment |

Auto-instrumentation (captures HTTP, DB, messaging without code changes):

```bash
# Node.js
node --require @opentelemetry/auto-instrumentations-node/register app.js

# Python
opentelemetry-instrument python app.py

# Java
java -javaagent:opentelemetry-javaagent.jar -jar app.jar
```

### Distributed debugging workflow

1. **Start with the trace** — find the failing span: which service? what latency breakdown?
2. **Correlate with logs** — use the trace ID to pull logs from that service (OTel propagates trace IDs across boundaries)
3. **Check metrics** — error rate spike (deployment issue?) vs latency spike (resource contention?) vs normal metrics (data-dependent bug?)
4. **Reproduce locally** — once you know which service, extract the input and reproduce with unit/integration tests

---

## Best Practices

1. Reproduce first — you cannot fix what you cannot reproduce
2. Isolate the problem — minimal reproduction before diagnosing
3. Read error messages fully — stack traces are more helpful than they look
4. Check recent changes first — most bugs are recent; start with `git log`
5. Use version control — `git bisect`, `blame`, and history are debugging tools
6. Take breaks — fresh eyes find what exhausted eyes miss
7. Document findings — help future teammates and future you
8. Fix root cause — not symptoms; use 5 Whys or Fishbone to get there
9. Write a regression test — every fix should include a test that would have caught it
10. Think in criteria — define what "fixed" looks like in observable, testable terms before starting

## Common Mistakes

| Mistake | Impact | Fix |
|---------|--------|-----|
| Making multiple changes at once | Can't tell which change fixed it | Change one thing at a time |
| Not reading the full error message | Misses the actual cause | Read the full stack trace first |
| Assuming it's complex | Wastes time | It's often a typo, off-by-one, or missing import |
| Debug logging left in production | Noise, potential data leakage | Use structured logging with levels |
| Not using a debugger | Slower than necessary | Learn your debugger; `console.log` isn't always best |
| Not testing the fix | Fix might not work | Verify it actually resolves the reported symptom |
| Optimizing without profiling | May not fix the real bottleneck | Profile first; then optimize the measured hotspot |

---

## References

1. [Superpowers: systematic-debugging](https://www.ququ123.top/en/2026/02/superpowers-systematic-debugging-skill/) — Four-phase methodology
2. [Fishbone vs Fault Tree Analysis](https://5xwhys.com/articles/rca/fishbone-vs-fault-tree/) — RCA method selection guide
3. [Root Cause Analysis Methods Compared](https://fivewhys.ai/blog/root-cause-analysis-methods-compared) — 5 Whys vs Fishbone vs FTA
4. [Node.js Profiling with V8 Inspector](https://oneuptime.com/blog/post/2026-01-06-nodejs-profiling-v8-inspector-chrome-devtools/view) — CPU profiling walkthrough
5. [Node.js Memory Management](https://dev.to/_d7eb1c1703182e3ce1782/nodejs-memory-management-and-profiling-find-and-fix-memory-leaks-in-2026-od4) — Memory leak detection
6. [OpenTelemetry for Microservices](https://totalshiftleft.ai/blog/opentelemetry-microservices-observability) — OTel auto-instrumentation guide
7. [Claude Code Autonomous Debugging](https://ralphable.com/blog/claude-code-autonomous-debugging-atomic-tasks) — Structuring tasks for AI debugging agents
