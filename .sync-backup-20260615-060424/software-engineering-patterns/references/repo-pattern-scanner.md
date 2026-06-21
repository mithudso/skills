<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `repo-pattern-scanner` skill.
> Sibling topics in this family are now reference files under the hubs (`software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: repo-pattern-scanner
description: "Scan a repository for reusable coding patterns and upload each to the mdb-context-hub shared library via tam_save_shared_library. TRIGGER: user says 'scan for patterns', 'find reusable code', 'sync patterns to hub', or triggers after cloning/onboarding a new repo, after a major refactor, or on a commit-triggered schedule. SKIP: single-file code review (use code-reviewer); security audit (use security-reviewer); full repo bootstrapping with meta-doc work (use repo-bootstrapper)."
version: "1.1"
updated: "2026-05-29"
category: developer
tags: [coding-patterns, repo-analysis, shared-library, mcp, automation]
whenToUse:
  - "scan this repo for patterns"
  - "find reusable code in this codebase"
  - "sync patterns to the hub"
  - "onboarding a new repo — extract coding patterns"
  - "after a major refactor, update the pattern catalog"
  - "run the pattern scanner on a schedule or commit trigger"
whenNotToUse:
  - "single-file code review — use code-reviewer"
  - "security audit — use security-reviewer"
  - "full repo bootstrapping with meta-doc work — use repo-bootstrapper"
related_skills: [repo-bootstrapper, repo-file-analyzer, software-engineering-patterns, security-reviewer]
---

# Repo Pattern Scanner

Scan a repository for reusable coding patterns and sync them to the mdb-context-hub shared library.

## When to use

- After cloning or onboarding a new repo
- After a major refactor that changes architectural patterns
- On a scheduled or commit-triggered basis to keep the hub current
- When the user says "scan for patterns", "find reusable code", or "sync patterns to hub"

## Inputs

The skill expects these arguments (passed via the agent prompt or directly):

- **repoPath** (required): Absolute path to the repo root
- **sourceRepo** (required): Human-readable repo name (e.g., `mdb-tam`, `mdb-case-assistant`)
- **runtime** (optional, default `node`): Runtime environment (`node`, `browser`, `python`)
- **dryRun** (optional, default `false`): When true, report findings without uploading

## Pattern catalog

Scan for these pattern families. Each match must include the file path, a short description of how the pattern is used, and the exported symbols (if any).

### Creational
- **Factory** — functions returning objects without `new`, typically named `create*` or `build*`
- **Builder** — method-chaining constructors (`.from().where().build()`)
- **Module singleton** — module-scoped instance exported directly

### Structural
- **Facade** — simplified wrapper over a complex subsystem (e.g., storage helpers over chrome.storage)
- **Adapter** — translates one interface to another (callback-to-promise, API normalization)
- **Proxy** — intercepts access (validation proxies, caching proxies, `new Proxy`)
- **Decorator** — function wrappers adding cross-cutting concerns (`withLogging`, `withRetry`)

### Behavioral
- **Observer / EventBus** — `.on()` / `.emit()` / `.subscribe()` patterns
- **Strategy** — interchangeable algorithms in a map or dispatch object
- **State machine** — explicit state + transitions objects
- **Command** — encapsulated operations with execute/undo
- **Chain of responsibility** — ordered handler pipelines
- **Dispatch table** — `const handlers = { TYPE: fn }` message routing

### Async / Resilience
- **Retry with backoff** — loops with exponential delay and jitter
- **Circuit breaker** — failure counting with open/closed/half-open states
- **Semaphore / concurrency limiter** — bounded parallel execution
- **Async queue** — sequential promise chaining
- **Single-flight / dedup** — inflight promise reuse to prevent duplicate calls

### Functional
- **Pipe / flow** — left-to-right function composition
- **Memoization** — cached pure-function results
- **Reducer** — `(state, action) => newState` immutable update functions

## Execution sequence

Follow these steps in order. Do not skip steps or reorder them.

### Step 1 — Discover source files

```bash
find {repoPath} -type f \( -name '*.js' -o -name '*.ts' -o -name '*.mjs' \) \
  ! -path '*/node_modules/*' ! -path '*/.git/*' ! -path '*/dist/*' \
  ! -path '*/build/*' ! -path '*/vendor/*' ! -path '*/coverage/*'
```

Apply quality gates immediately: skip files > 2000 lines, skip minified (avg line length > 200). Separate test files (`*test*`, `*spec*`, `__tests__/`) into a test-only set — patterns found exclusively there are excluded from upload.

### Step 2 — Keyword pre-filter

Grep each non-test file for pattern indicators. Only proceed to Step 3 for files matching at least one:

```
create[A-Z]|build[A-Z]
\.on\(|\.emit\(|\.subscribe\(|EventEmitter|EventBus|PubSub
new Proxy|handler.*set.*target
withRetry|retryWithBackoff|exponentialBackoff|backoff
circuitBreaker|CircuitBreaker|halfOpen|HALF_OPEN
stateMachine|StateMachine|transitions.*=
createSemaphore|Semaphore|concurrencyLimit
handlers\[msg|handlers\[type|handlers\[action
pipe\(|flow\(|compose\(
memoize|memo\(
reducer.*state.*action
createChain|chainOfResponsibility
\.build\(\)|\.where\(.*\.from\(
```

### Step 3 — Classify and extract

For each matching file, read the relevant section (first 300 lines or the matching region ± 50 lines). Identify which pattern(s) from the catalog it implements. For each pattern, extract:

| Field | Source |
|-------|--------|
| Pattern name | From the catalog above |
| Family | `creational`, `structural`, `behavioral`, `async-resilience`, or `functional` |
| File path | Relative to repo root |
| Exported symbols | `export function`, `export const`, `export class`, `module.exports` names |
| Description | One sentence: what the pattern does in this repo's context |
| Keywords | 3-8 search terms specific to this pattern instance |

**Gate**: Require at least one exported symbol. If the pattern has no exports (internal-only helper), skip it.

### Step 4 — Deduplicate

Group entries by pattern name. If the same pattern appears in multiple files, merge into one entry:
- `originalPath`: the primary (most complete) implementation
- `notes`: list all other file paths

### Step 5 — Upload to hub

For each unique pattern, call `tam_save_shared_library` with:

```
title:       "{Pattern Name} — {sourceRepo}"
sourceRepo:  the repo name
kind:        the family (e.g., "async-resilience")
runtime:     from input (default "node")
description: from Step 3
exports:     from Step 3
keywords:    from Step 3
tags:        ["coding-pattern", "{pattern-kebab}", "{family}"]
originalPath: primary file
useCases:    1-3 sentences on when to reuse
notes:       additional file paths if merged
overwrite:   true
```

If `dryRun` is true, print the payload instead of calling the MCP tool.

### Step 6 — Report

Print a markdown summary table as the final output:

```
## Pattern Scan Results — {sourceRepo}

Scanned {N} files, {M} matched pre-filter, {P} patterns uploaded.

| # | Pattern | Family | File(s) | Exports | Hub Status |
|---|---------|--------|---------|---------|------------|
| 1 | Retry with Backoff | async-resilience | src/lib/retry.js | retryWithBackoff | uploaded |
```

If `dryRun`, replace "uploaded" with "dry-run".

## Quality gates (summary)

- Skip files > 2000 lines (generated/bundled)
- Skip minified files (avg line length > 200)
- Skip patterns found only in test files
- Require at least one exported symbol per entry
- Deduplicate: one hub entry per pattern name per repo
