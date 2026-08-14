# Repo Bootstrapper — Coding Standards for External Calls

Every external call — HTTP fetch, MCP tool invocation, child_process spawn, MongoDB read/write,
message-bus publish, queue enqueue, third-party SDK call — must satisfy all five standards.
Audits fail when any one is missing.

## Standard 1 — CLI or MCP trigger

Every call must be invokable as a single named operation from a CLI or MCP tool. Bare in-process
invocations (route handler / cron / event listener only) are not sufficient.

- Each call registers in a single operations-registry module. Reference implementation:
  `mcp-server/src/operations-registry.ts` in `10gen/mdb-tam`. Other repos publish an equivalent
  file and reference its path in `docs/cli-and-operations.md`.
- Registry entry carries: `id` (kebab-case), `kind`, `transport`, target, default args.
- CLI exposes at least: `list`, `run <id>`, `status`, `history <id>`, `reset <id>`.
- Operations with side effects on customer data: mark `readOnly: false` with a warning string.

## Standard 2 — Centralized error logging

Every failed attempt must append one structured JSONL line to the central error log.
Convention: `~/.<repo-slug>/errors.jsonl`. Path recorded in `docs/error-monitoring-guide.md`.

Schema per line: `ts, opId, kind, transport, attempt, attemptsTotal, errorCode, errorMessage,
remediation, stack, repo, processPid, args` — secret-looking arg keys redacted.

- Error log is append-only; reset via explicit `errors-clear` command only.
- All error sources (dashboard backend, extension SW, MCP server, CLI) share one sink or a
  documented bridge.

## Standard 3 — Auto-remediation per error code

Every registry entry declares a `remediation` map keyed by every error code it can produce.

- **safe** = remediator agent may auto-apply (transient transport failure, rate limit, service degraded)
- **escalate** = must surface to human (auth missing, validation failed, config missing, customer-data side effect)

Audit cross-checks: every error code referenced in the call's tests, registry, or source must
appear as a key in the remediation map.

## Standard 4 — Dashboard card per call

**Minimum surface (required for audit pass):**
- Status: last-attempt outcome (ok / failed / unknown) + timestamp + duration
- Logs: most recent attempt's structured log lines visible in-card
- Paths: file:line, registry id, relay command or HTTP path, documentation link
- Resync: manually trigger the call (or "no manual trigger" notice)
- Clear history: flush this card's in-session attempt log

**Aspirational surface (when underlying call supports it):**
- Test: one-click run against mock or read-only target
- Stop: abort via AbortSignal (required for calls that can take >5s)
- Restart: re-issue stopped/failed invocation with same args

Cards must be auto-generated from the operations registry — never hand-rolled per-call markup.
Reference: `src/dashboard/external-calls-tab.js` in `10gen/mdb-case-assistant`.

## Standard 5 — Datastore-landed verification

Every call expected to produce persisted data must verify that data actually landed.
Verification is structural (row exists with expected id?), not semantic.

- After a feed run: read back at least one expected record from the typed collection.
- After an MCP tool call that writes state: read back and confirm the write is visible.
- Failed verification logs code `DATASTORE_VERIFY_FAILED` — document in `docs/error-monitoring-guide.md`.
- Calls that intentionally do not persist data: declare `verifyDatastore: false`.
- Calls omitting the field default to verification ON.

## Per-call audit table

| Standard | Pass criteria |
|---|---|
| 1. CLI/MCP trigger | Registered in operations registry; runnable via `<cli> run <id>` |
| 2. Centralized error log | Failure path writes one JSONL line matching the documented schema |
| 3. Auto-remediation | `remediation` map keyed by every error code; each entry marked safe/escalate |
| 4. Dashboard card (minimum) | Card visible with status / log / paths / docs / resync / clear-history |
| 5. Datastore verify | `verifyDatastore` field present; default-on calls produce `DATASTORE_VERIFY_FAILED` on read-back failure |

A repo passes the audit only when every call satisfies all five standards.
