# Transactions, Locking & Security

> Reference for **postgresql-expert**. Isolation levels, the locking model, deadlock diagnosis, and the role/privilege/RLS access model. Verified-as-of 2026-06-30 (PG18 line).

## Isolation levels

PostgreSQL implements three of the four SQL isolation levels on top of MVCC (it has **no dirty reads** at any level):

| Level | Prevents | Notes |
|---|---|---|
| **Read Committed** (default) | Dirty reads | Each *statement* sees a fresh snapshot — two queries in one txn can see different data; non-repeatable & phantom reads possible |
| **Repeatable Read** | + non-repeatable & phantom reads | Snapshot taken at first statement, held for the txn. Stronger than the SQL-standard RR (PostgreSQL's RR is true snapshot isolation). Writers may hit `could not serialize access` (40001) |
| **Serializable** (SSI) | + serialization anomalies | **Serializable Snapshot Isolation** — true serializability via predicate-style read/write tracking; aborts one txn of a dangerous cycle with **40001** |

Key practical points:

- **Read Committed is per-statement**, which surprises people: a multi-statement transaction is *not* a consistent point-in-time view. Use Repeatable Read for read-consistent reports.
- **Serializable (SSI)** gives you correctness without explicit locking, but your application **must retry transactions on serialization failure (SQLSTATE 40001)**. Don't choose Serializable unless you implement retry. It also requires *all* concurrent transactions to be Serializable to get the guarantee.
- `default_transaction_isolation` sets the cluster/session default; `SET TRANSACTION ISOLATION LEVEL …` per txn.

## The locking model

Two broad lock families:

### Table-level locks (relation locks)
Eight modes from `ACCESS SHARE` (taken by `SELECT`) up to `ACCESS EXCLUSIVE` (taken by `DROP`, `TRUNCATE`, most `ALTER TABLE`, `VACUUM FULL`). The conflict matrix matters for **DDL during traffic**: an `ALTER TABLE` needing `ACCESS EXCLUSIVE` queues behind running queries *and blocks new ones behind it* — a classic production stall. Mitigations: `SET lock_timeout` before DDL, run migrations in low-traffic windows, prefer non-blocking variants (`ADD COLUMN` with default = metadata-only PG11+; `CREATE INDEX CONCURRENTLY`; `ADD CONSTRAINT … NOT VALID` then `VALIDATE`).

### Row-level locks
- `FOR UPDATE` / `FOR NO KEY UPDATE` / `FOR SHARE` / `FOR KEY SHARE` on `SELECT` lock specific rows. Row locks **don't block readers** (MVCC) — only conflicting writers.
- **`SELECT … FOR UPDATE SKIP LOCKED`** is the idiomatic, scalable **work-queue** pattern — concurrent workers each grab unlocked rows without blocking each other. `FOR UPDATE NOWAIT` fails fast instead of waiting.
- **Advisory locks** (`pg_advisory_lock`, `pg_try_advisory_lock`, xact-scoped variants) are application-defined locks not tied to any row — use for singleton jobs, leader election, or cross-session mutexes.

## Deadlocks & lock-wait diagnosis

A **deadlock** is a cycle of lock waits; PostgreSQL detects it after `deadlock_timeout` (default 1s) and aborts one transaction with SQLSTATE **40P01**. They are almost always caused by transactions acquiring locks on the same rows **in different orders**.

- **Prevent:** acquire locks in a consistent order across code paths; keep transactions short; touch rows in a deterministic sequence (e.g. `ORDER BY id` before update); avoid user think-time inside a transaction.
- **Diagnose live waits:** `pg_locks` joined to `pg_stat_activity` shows blocker/blocked pairs; `pg_blocking_pids(pid)` returns who's blocking a given backend. Set `log_lock_waits = on` to capture waits exceeding `deadlock_timeout` in the log.
- **`idle in transaction`** sessions holding locks (and xmin) are a frequent culprit — cap with `idle_in_transaction_session_timeout`; also `statement_timeout` and `lock_timeout` as guardrails.

## Roles, privileges & RLS

PostgreSQL has a unified **role** concept (a role can be a "user" with `LOGIN` and/or a "group" granted to other roles):

- **Privileges** are granted on objects (`GRANT SELECT ON … `, `GRANT USAGE ON SCHEMA`, `GRANT EXECUTE ON FUNCTION`). **`DEFAULT PRIVILEGES`** (`ALTER DEFAULT PRIVILEGES`) set what future objects inherit — essential to get right, or new tables won't be readable by your app role.
- **Schema `public`**: since PG15 the `public` schema no longer grants `CREATE` to all roles by default — a security hardening that surprises upgraders.
- **`SET ROLE` / `SECURITY DEFINER` functions** — run with the function owner's privileges; always `SET search_path` in `SECURITY DEFINER` functions to prevent search-path hijacking.
- **Predefined roles** (`pg_read_all_data`, `pg_monitor`, `pg_read_all_stats`, …) grant common bundles without superuser.

### Row-Level Security (RLS)
RLS filters which *rows* a role can see/modify — the foundation for **multi-tenant isolation** in a shared table:

```sql
ALTER TABLE docs ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON docs
  USING (tenant_id = current_setting('app.tenant_id')::int);
```

- `USING` filters reads (and the visible rows for update/delete); `WITH CHECK` constrains what rows may be written. Separate policies per command (`SELECT`/`INSERT`/`UPDATE`/`DELETE`) and per role.
- **Caveats:** table owners and `BYPASSRLS` roles skip RLS (use `FORCE ROW LEVEL SECURITY` to apply it to the owner too); policies are predicates the planner must include, so index them (e.g. index `tenant_id`); and set the tenant key from a trusted source (a `SET` from the pooler/app session, validated server-side) — never trust a client-supplied value without verification.
- RLS is enforcement-in-depth, not a substitute for getting the application's tenant scoping right; combine with connection-level role separation.

> **Mental model:** MVCC means reads and writes rarely block each other — contention is writer-vs-writer on the same rows, and the cure is short, consistently-ordered transactions. Choose Read Committed unless you specifically need snapshot consistency (RR) or true serializability (SSI, with retry). For tenancy, lean on RLS *plus* disciplined application scoping, not either alone.
