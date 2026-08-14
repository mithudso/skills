# wt-code-review-quick

**Category:** Science, Biology & Medicine
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/storage-engines/wt-code-review-quick

## Description
Use when the user wants a fast review of a WiredTiger change without the full multi-specialist pipeline. Triggers on the same phrasing as wt-code-review ("review this PR", "check this diff", "audit my changes") combined with a speed hint such as "quick", "fast", "just the basics", or "don't need the full review". The top-level agent walks the diff directly against a distilled checklist that spans all major review surfaces — no subagents, no triage, no verifier. Choose this over wt-code-review when turnaround time matters more than deep specialist coverage; choose wt-code-review when the diff is large, novel, or touches risky surfaces such as a new state machine or new disagg invariant.

---

# Quick WiredTiger code review

Three phases. No subagents — the top-level agent does the whole review. If the diff is large, novel, or touches a risky surface (new state machine, new visibility helper, new disagg invariant), redirect to `wt-code-review` for the full pipeline.

If the diff is not WiredTiger (no `src/` WT paths), stop — these checks assume the WiredTiger codebase.

## Phase 1 — Acquire the diff

Resolve the user's request to a concrete diff and write it to `/tmp/wt-quick-<short-id>.diff`:

- PR number → `gh pr diff <num>`
- Branch → `git diff develop...HEAD` (or `main...HEAD`)
- Working tree → `git diff` (staged + unstaged)
- Pasted diff → write to disk verbatim

If the diff is empty, stop and ask what the user wants reviewed.

## Phase 2 — Walk the checklist

Read the diff. For each modified function, read enough surrounding source to judge candidate findings (use `Read`/`Grep` as needed). Apply the checks below, scoped to whichever ones the diff actually touches. **Investigate, don't speculate** — every finding must point to a line in the diff. Skip categories the diff doesn't exercise.

### Correctness — visibility, lifecycle, C/UB

- `__wt_txn_visible` (per-txn snapshot) vs `__wt_txn_visible_all` (obsolete to every reader, uses `durable_ts`). Cleanup / eviction / RTS / HS-pruning need `_all`; read paths need the per-txn variant.
- RTS compares `rollback_ts < upd->upd_durable_ts`. Any new comparison against `upd_start_ts` or `upd_commit_ts` is wrong.
- Inside checkpoint (`F_ISSET(session->txn, WT_TXN_IS_CHECKPOINT)`) or on a disagg follower, reads of `txn_global->{oldest,stable,pinned}_timestamp` are suspect — use the session-scoped `checkpoint_*_timestamp` or `disaggregated_storage.last_checkpoint_*_timestamp`.
- `WT_CONFIG_ITEM.str` is NOT null-terminated — no `strcmp` / `strlen` / `%s`; use `WT_CONFIG_MATCH` / `WT_STRING_MATCH`.
- Prefer `WT_TS_NONE` / `WT_TS_MAX` over `0` / `UINT64_MAX`.
- Session close order: clear `WT_SESSION_CACHE_CURSORS` → close cursors → rollback txn → release snapshot → close metadata → `__wt_hazard_close`. A new `goto err` before `__wt_hazard_close` leaks hazard pointers.
- Txn rollback releases snapshot *before* iterating `txn->mod[]`.
- Page modify order: `__wt_page_modify_init` → mutate → `__wt_page_modify_set` (the release-store publishes to eviction).
- Cursor reuse needs `reset` first.
- Unsigned underflow: `v - 1` where `v` can be `0` wraps to `UINT*_MAX`.
- Timestamps / txn ids / LSNs stay `uint64_t` end-to-end; flag any storage in `uint32_t` / `int`.
- Format: `PRIu64` for `uint64_t`, `WT_SIZET_FMT` for `size_t`. No `%d` / `%lu`.
- `sizeof(p)` of a pointer where struct size was intended (especially `__wt_calloc(1, sizeof(p))`).
- No side effects in `WT_ASSERT` conditions (compiled out in release).
- `WT_DECL_ITEM(buf)` is NULL until `__wt_scr_alloc` — accessing fields before alloc segfaults.
- `switch` fallthrough without `/* FALLTHROUGH */`.

### Concurrency and memory ordering

- New shared field accessed without a lock needs the `wt_shared` annotation.
- Plain `=` or `++` on a `wt_shared` field is a bug — use `__wt_atomic_store_*` / `__wt_atomic_add_*`.
- Publish/consume pairs: store with `_release`, load with `_acquire`. Mismatched pairs (release store + relaxed load, or relaxed store + acquire load) are bugs on ARM/POWER.
- Pointer publish: `__wt_atomic_store_ptr_release` after writing pointed-to data; `__wt_atomic_load_ptr_acquire` before dereferencing.
- A `wt_shared` flag that gates access to other data must be loaded with at least `_acquire`, not `_relaxed`.
- Spin loops need `WT_PAUSE()` in the body; re-reads use `WT_READ_ONCE` or a relaxed atomic load, not a plain local read.
- Operations on `flags_atomic` fields use `F_SET_ATOMIC_*` / `F_CLR_ATOMIC_*` / `F_ISSET_ATOMIC_*`, not the non-atomic variants.
- Deprecated: `WT_ACQUIRE_READ_WITH_BARRIER` → `__wt_atomic_load_*_acquire`; `WT_RELEASE_WRITE_WITH_BARRIER` → `__wt_atomic_store_*_release`.
- Lock ordering: schema first; no read→write upgrade on handle-list or table; no mixing table-read with handle-list.
- Use `WT_WITH_*_LOCK` macros, not raw `__wt_spin_lock` + `WT_RET` (raw lock + early return leaks the lock). No I/O, malloc, or callouts inside a spinlock.
- Condvar wait must re-check predicate after wakeup; signaler must publish state before signaling.
- `__wt_session_gen_enter` paired with `_leave` on every return path; nothing blocking (cond_wait, I/O) inside.
- `__wt_hazard_set` paired with `_clear`; a missed clear leaves the page un-evictable.
- New background thread: `WT_THREAD_GROUP` lifecycle, private session, `WT_THREAD_RUN` cleared on shutdown, join after unlock.
- CAS loops: bounded or use `__wt_spin_backoff`. Inner retries can use `_relaxed`; final CAS provides ordering.
- x86 ≠ ARM. A pattern that "works" on x86 (TSO) may be broken on weak-memory targets.

### Error handling and resource cleanup

- `WT_RET` after an allocation (or any acquired resource without `err:` label) leaks. Change to `WT_ERR`.
- `WT_TRET` does not log — use `WT_TRET_MSG` in cleanup paths where a diagnostic is needed.
- No `ret2` / `tmp_ret` / `truncate_ret` locals — reuse `ret` via `WT_TRET`.
- `WT_RET(__wt_*_lookup(...))` where a miss is normal → `WT_RET_NOTFOUND_OK` (silent) or `WT_ERR_NOTFOUND_OK(..., true)` (handle).
- `WT_ERR_NOTFOUND_OK(expr, false)` followed by `if (ret != 0)` — `keep=false` clears `ret` on not-found, so the subsequent check never fires.
- `return (__wt_panic(...))` or `ret = __wt_panic(...); goto err;` → use `WT_RET_PANIC` / `WT_ERR_PANIC`.
- Cleanup function that accumulates errors via `WT_TRET`/`WT_TRET_MSG` then `return (0)` silently discards them. Return `ret`.
- Function whose `err:` label mixes `__wt_free`/`__wt_scr_free` with substantive work like `WT_ERR(__wt_log_write(...))` violates single-responsibility — extract a helper.
- `goto err` on a success path is wrong; `err:` is for failure only.
- `WT_DECL_ITEM(tmp)` declared but never passed to `__wt_scr_alloc`, only `__wt_scr_free`-d at `err:` — delete the declaration entirely.
- Redundant `if (tmp != NULL) __wt_scr_free(session, &tmp)` — the macros handle NULL.
- Missing `WT_CLEAR(x)` — prefer it over `memset(&x, 0, sizeof(x))` for stack zeroing in non-declaration contexts.

### Assertions and defensive checks

- Default to `WT_ASSERT`; escalate to `_ALWAYS` only when a silent violation in production would corrupt data and the check is cheap; `_OPTIONAL` for expensive opt-in checks.
- Raw `assert()` from `<assert.h>` is forbidden.
- No side effects in `WT_ASSERT` condition — it can compile to nothing.
- Assert the *positive* invariant. Polarity bugs (`cacheable || !bulk` when the rule is "bulk cursors must not be cacheable" → `!cacheable || !bulk`).
- Split asserts on unrelated arms, or branch on path first — a single `A || B` assert that fires doesn't tell you which arm failed.
- Always-on macros (`_ALWAYS` / `_OPTIONAL` / `_ERR_ASSERT` / `_RET_ASSERT` / `_PANIC*`) take a message — include the offending value and state the invariant in words, not "Assertion failed".
- Don't assert non-NULL right before dereferencing — the deref crashes anyway. Don't NULL-check before `__wt_free` / `__wt_scr_free`. Don't hand-roll `if (ret != 0) { __wt_err(...); return (ret); }` — use `WT_RET_MSG`.
- Reachability classification: internal contract bug → assert. API misuse → `EINVAL`. OS/external → matching errno (`ENOMEM`, `EBUSY`, `WT_NOTFOUND`). Unrecoverable → `WT_RET_PANIC` / `WT_ERR_PANIC`.

### API contract

- Test code in `test/suite/` / `test/csuite/` / `test/catch2/` directly accessing `_IMPL` structs (`WT_SESSION_IMPL`, `WT_CONNECTION_IMPL`, etc.) — extend the public API instead.
- Changing parameter list, return type, or name of a function in `src/include/wiredtiger.h.in` or in `WT_SESSION` / `WT_CURSOR` / `WT_CONNECTION` / `WT_STORAGE_SOURCE` vtables is ABI-breaking — add an `_ext` variant.
- Adding a `version` / `generation` field to a non-on-disk struct — versioning is for persistence; in-memory should use flags or config strings.
- Wrong error code: `ENOTSUP` for unsupported, `EINVAL` for bad argument, `EBUSY` for temporarily unavailable, `WT_ERROR`/`WT_PANIC` for internal errors.
- New method that's `existing_method + one boolean` should be a flag or config string, not a new vtable entry.

### Disagg storage

Only apply if the diff touches `src/block_disagg/`, `src/conn/conn_layered*`, `__wt_disagg_*`, `WT_BTREE_DISAGGREGATED`, `block_manager=disagg`, `type=layered`, `role=leader|follower`, or related symbols.

- Mutation, reconciliation, or `plh_put` on a disagg-eligible btree needs a leader guard like `!F_ISSET(btree, WT_BTREE_DISAGGREGATED) || layered_table_manager.leader`. Compare to siblings in `src/cursor/cur_layered.c`.
- Refactor that deletes, weakens (`&&` → `||`), or moves a leader guard past a side-effectful call is almost always wrong.
- Layered ingest tables use an encoded tombstone via `__clayered_deleted_encode`; `cursor->remove()` through a plain file cursor corrupts the format.
- L1 is in-memory only — never written to shared storage. On step-up, L1 must drain into L0 *before* the leader exposes itself as writer.
- Address cookie invariants: `page_id` non-zero, `lsn >= base_lsn`, `size > 0`. Checkpoint cookie = root address cookie.
- `plh_put` must be durable before `plh_discard` of the same version; async paths that let discard race ahead are corruption bugs.
- New discard path on a disagg btree must consult `__wt_btree_can_discard` / `__wt_materialization_check`. Don't evict dirty internal pages.
- Precise checkpoint reconciliation must not write unstable updates. No RTS on disagg restart.
- Delta chain bound: `block_meta->delta_count < page_delta.max_consecutive_delta`; next write at the limit must be a full page.
- No local turtle, no local free-block lists on disagg — metadata flows through `__wt_disagg_enqueue_metadata_operation` under the schema lock.

### Performance — hot paths only

Apply only if the modified function lives on a hot path: eviction, reconciliation, metadata parsing, checkpointing, or inside a per-page / per-record loop.

- New `__wt_malloc` / `__wt_calloc` / `__wt_scr_alloc` / `WT_DECL_ITEM` inside a per-page or per-record loop — hoist before the loop or restructure to allocate only on the slow path.
- `WT_STAT_*_INCR` or `__wt_atomic_add` inside a per-page loop with concurrent workers becomes a cache-line bottleneck — accumulate locally and flush with `WT_STAT_*_INCRV` at scan end.
- Expensive predicate before cheap predicate: `strcmp` / `__wt_config_gets` before an integer or flag comparison that could short-circuit it.
- Magic threshold without a named constant: `if (cache_used > 90)` — what unit, why 90? Use `#define WT_CACHE_EVICT_TRIGGER_PCT 90`.

### Verbose logging and statistics

- Wrong `WT_VERB_*` category: a read-path message under `WT_VERB_VERIFY`, an eviction message under `WT_VERB_CHECKPOINT`, a disagg message under a generic category instead of `WT_VERB_DISAGGREGATED_STORAGE` / `WT_VERB_LAYERED`. Prefer the operation's category over the caller's.
- Wrong verbosity level: `WT_VERBOSE_ERROR` only for actual errors (not contention, retries, expected misses); `WT_VERBOSE_DEBUG_1` is the default — per-page or per-record events go to `DEBUG_2` or higher.
- Stat name with `error` / `fail` / `err` for non-error events (retries, CAS races, contention, expected misses). Use `_race`, `_retry`, `_contention`, `_miss`.
- Log message on a path that's reached on both success and failure must indicate outcome (e.g. `"completed (%s)", ret == 0 ? "success" : "error"`).
- Open string concatenation across `__wt_verbose_*` arguments — use `%s` and pass the name.

### Style and idiom

`CONTRIBUTING.rst` is the source of truth. Trust `dist/s_clang_format` for operator spacing and brace placement — don't flag what it handles. Skip generated files (between `AUTOMATIC ... GENERATION START/STOP` markers, `dist/`, `src/include/*_auto.h`).

- `//` comments in production C code — use `/* ... */` block style with each continuation line starting `" * "`.
- New exported function without a banner header (`/* __wt_foo -- ... */`).
- Comments explaining *what* instead of *why* on non-obvious logic.
- Comment naming a specific identifier (function, macro, struct field) that no longer exists — grep for it; if absent, flag as stale.
- `TODO` / bare `FIXME` / closed-ticket reference — use `FIXME-WT-NNNNN:` with an open Jira ticket.
- Raw `malloc` / `free` / `strcpy` / `strcat` / `sprintf` / `printf` / `assert` — use the WT wrappers.
- Manual `flags |= mask` — use `F_SET` / `F_CLR` / `F_ISSET` (or `F_*_ATOMIC_*` on `flags_atomic`).
- Manual `= NULL` after `__wt_free` — the macro handles it.
- `return value;` without parens; `dist/s_clang_format` enforces `return (value);`.
- Naming: `__wt_` (extern), `__wti_` (internal-shared), `__` (static); output-param `p` suffix (`resultp`); canonical short names (`session`, `conn`, `btree`, `ref`, `page`, `cbt`).
- Declaration ordering inside a block: `struct` → `WT_*` alpha → declaration macros (`WT_DECL_RET`, `WT_DECL_ITEM`) → built-ins by type, alpha within. Blank line after.

### Tests

Apply only if the diff touches `test/suite/*.py`, `test/catch2/`, `test/csuite/`, `test/**/*.cpp`.

- `time.sleep()` / `sleep()` / wall-clock polling as a synchronization mechanism. Assert on a WiredTiger stat instead.
- Boundary test uses a value safely beyond the boundary (e.g. `ts=25` when feature triggers at `ts=20`) — add a test at the exact boundary and one just before.
- New regression test without a `# Fails without the fix` comment — reviewers will ask.
- New test duplicating coverage already in an existing file — extend the existing test instead. A Python and a C++ test covering exactly the same ground is also flagged.
- Forward-only traversal (`cursor.next()` loop) without a `cursor.prev()` companion on range scans, `search_near`, and truncation boundary checks.

## Phase 3 — Summarize

Produce a single response with two sections:

1. **Findings**, grouped by severity (`BUG` → `WARNING` → `NIT`) then by file. Cite `file:line` on every entry. For each finding give one sentence on what's wrong and one on the fix. If nothing was found, say so in one sentence.
2. **Coverage**, one line per checklist section that you actually exercised on this diff (skip ones the diff doesn't touch). This lets the user see what was *not* checked.

End with a one-line **Suggested next steps** — typically `cd dist && ./s_fast` if style flagged anything, plus "For deeper coverage run `wt-code-review`" if the diff touches disagg, concurrency, memory-ordering, or anything you flagged as risky.

Keep the summary scannable in under 30 seconds.

## Guardrails

- **Single agent.** Do not spawn subagents — this skill is intentionally the lightweight path. If the diff is large or risky enough to need parallel specialists, redirect to `wt-code-review`.
- **Investigate, don't speculate.** Every finding must point to a line in the diff. If a check requires reading more source than you have time for, note the area as "needs deeper review" rather than guessing.
- **Severity calibration.** `BUG` = would fail tests, corrupt data, or break ABI. `WARNING` = likely wrong, needs author's eyes. `NIT` = minor. Over-flagging trains the reader to ignore findings.
- **Stay in scope.** Skip checklist sections the diff doesn't touch. Skip generated files. Don't critique architecture beyond the diff level.
- **Stop on empty diff.** Ask before doing anything.
- **Respect user scope hints.** "Just style" → run only the style section. "Look at the disagg parts" → emphasise that section.