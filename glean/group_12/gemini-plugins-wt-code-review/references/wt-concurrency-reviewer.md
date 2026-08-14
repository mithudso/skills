# WiredTiger concurrency and memory-ordering review

Everything that breaks when multiple threads touch the same state: locking, thread lifecycle, inter-thread coordination, atomic memory ordering. Many concurrency bugs only manifest on weak-memory architectures (ARM, POWER, RISC-V) even when the code appears to work on x86 — focus on the C11 memory model, not whether it happens to work on TSO.

## Scope

**In scope:** spinlocks, rwlocks, `WT_WITH_*_LOCK` macros, condition variables (`WT_CONDVAR`), futexes, thread groups, generations, hazard pointers, session lifecycle, CAS loops, state machines, shutdown/cleanup paths, lock ordering, atomic loads/stores (`__wt_atomic_*`), memory barriers (`WT_FULL_BARRIER`, `WT_ACQUIRE_BARRIER`, `WT_RELEASE_BARRIER`), `wt_shared` annotations, `WT_READ_ONCE` / `WT_WRITE_ONCE`, `F_*_ATOMIC_*` flag macros, publish/consume pairs.

## Reference

Source of truth: `src/include/gcc.h`, `src/include/hardware.h`, `src/include/mutex.h`, `src/support/generation.c`, `src/support/hazard.c`.

Key primitives:
- `__wt_atomic_load_<T>_acquire` / `__wt_atomic_store_<T>_release` — standard publish/consume pair.
- `__wt_atomic_cas_<T>` / `_relaxed` — CAS with/without ordering.
- `_relaxed` variants — stats, counters, fields inside exclusive locks.
- `wt_shared` — required annotation on every field accessed lock-free across threads.
- `WT_PAUSE()` — required inside every spin loop body.
- Lock macros: `WT_WITH_SCHEMA_LOCK`, `WT_WITH_METADATA_LOCK`, `WT_WITH_HANDLE_LIST_{READ,WRITE}_LOCK`, `WT_WITH_TABLE_{READ,WRITE}_LOCK`, `WT_WITH_CHECKPOINT_LOCK`, `WT_WITH_HOTBACKUP_READ_LOCK`. Use the macro and propagate the captured error (`WT_WITH_SCHEMA_LOCK(session, ret = __schema_op(...)); WT_ERR(ret);`) — never a raw lock + `WT_RET`.

Greps: `grep -nE 'WT_ATOMIC_FUNC' src/include/hardware.h`; `grep -nE 'WT_WITH_' src/include/mutex.h`; `grep -nE 'WT_GEN_' src/include/generation.h`.

Deprecated: `WT_ACQUIRE_READ_WITH_BARRIER` → `__wt_atomic_load_<T>_acquire`; `WT_RELEASE_WRITE_WITH_BARRIER` → `__wt_atomic_store_<T>_release`. Flag on sight — `WT_ACQUIRE_READ_WITH_BARRIER` is a TSAN suppressant, not a real ordering primitive; demand real `acquire` / `relaxed` atomics.

Common memory-ordering bug patterns:

| Pattern | Problem | Fix |
|---------|---------|-----|
| `if (flag) { ... use data ... }` with relaxed flag load | No ordering between flag and data reads | Load flag with `_acquire` |
| Release store then immediately acquire load in same thread | Wrong direction — does nothing | Review the synchronization intent |
| `__ATOMIC_SEQ_CST` everywhere | Unnecessarily expensive on ARM | Downgrade to acquire/release where sufficient |
| `__ATOMIC_RELAXED` on a flag that gates other data | Load must be at least acquire | Use `_acquire` load |
| `volatile` without atomics | Stops compiler reordering, not CPU reordering on ARM | Add an atomic primitive |
| CAS retry loop with seq_cst inner ops | Costly | Inner retries `_relaxed`; final CAS provides ordering |

## What to flag

### Locks, threads, generations, hazards

1. **Lock order.** What locks are acquired and in what order? Cross-check the asserted hierarchy — schema lock must be acquired before checkpoint lock (`WT_ASSERT` fires on `session->lock_flags`). Flag any added/removed lock not justified in the commit message; silently dropping a lock (e.g. checkpoint lock disappearing from a disagg pickup path) is a flag.
2. **Unlock on every exit.** Every locked region — does every exit path unlock? `WT_ERR` (not `WT_RET`) inside a raw lock/unlock pair; lock release belongs in `err:`. Inside `WT_WITH_*_LOCK` it's safe. Lock-macro ret variables must be `__`-prefixed (`__lock_ret`, `__table_lock_ret`); `WT_CONFLICT_*` sub-codes are valid only for `EBUSY` — `EAGAIN`/`EINVAL` from a locking primitive must report `WT_NONE`.
3. **Dhandle-list iteration.** Hold the read lock across the *entire* loop — releasing mid-iteration lets entries be added/removed.
4. **New background thread.** Verify `WT_THREAD_GROUP` lifecycle, private session, `WT_THREAD_RUN` cleared on shutdown, join after unlock.
5. **`__wt_session_gen_enter`.** Matching `_leave` on every return path; nothing blocking in between.
6. **`__wt_hazard_set`.** Matching `_clear`; scan paths hold `WT_GEN_HAZARD`.
7. **Condvar wait.** Predicate re-checked after wakeup? Signaler publishes state before signaling?
8. **CAS loop.** Bounded or `__wt_spin_backoff`. State-transition diagram covered? Inner retries can use `_relaxed`; the final CAS provides ordering.
9. **State machine field.** All transitions valid? No intermediate state read as if final? A `WT_*_PENDING` intermediate state plus release semantics is *not* a substitute for the right lock or atomic discipline — gating on `WT_TRUNCATE_COMMIT_INIT` rather than `WT_TRUNCATE_COMMIT_PUBLISHED` is reading a struct before it's published.
10. **Lock-choice and assertions.** Read-heavy-per-cursor-op with an infrequent writer → rwlock, not spin (switching late is cheap, leaving the spinlock is expensive). Assert lock-held invariants explicitly on functions touching shared queues rather than relying on caller context; use `__wt_rwlock_islocked()` / `WT_ASSERT_SPINLOCK_OWNED` for the assertion — recording "I hold the read lock" via a bit on the shared structure races. The `caller_locked` / `locked` bool parameter is the documented pattern for a function callable both inside and outside a lock (better than a `_no_lock` twin); for reentrancy use `__wt_spin_owned(session, lock)` plus a conditional acquire/release rather than threading a bool through every helper.
11. **Misc.** `__wt_failpoint` is probabilistic — when every-time injection is needed use `timing_stress_flags`. `PANIC` is sometimes load-bearing in shutdown/restart paths — don't downgrade it casually. Per-thread independent interval counters are wrong for shared progress reports — coordinate a single writer per period via CAS on a shared `last_report_clock`. Mock / demo implementations of a thread-safe contract still need locks/CAS — silently racing teaches the wrong contract. Avoid timeout-based concurrency tests; prefer deterministic flags.

### Atomic memory ordering

12. **Locate shared fields.** Flag any newly added field accessed lock-free across threads that's missing the `wt_shared` annotation. Single-writer fields don't need lock-free annotations *unless* there's an unsynchronized reader — justify based on the actual reader/writer set. Inside an exclusive lock (`txn_global->rwlock`, RTS exclusive access) still use `_relaxed` atomics on shared fields — the atomic call marks the field shared and prevents future-bug regressions.
13. **Trace each shared access.** For every read/write of a `wt_shared` field: wrapped in an appropriate `__wt_atomic_*` call? Ordering strength (relaxed / acquire / release / seq_cst) correct for the use? Plain `=` or `++` instead of an atomic op is a bug. **Default to relaxed** for counters, phase fields, progress tracking; escalate to acquire/release/seq_cst only with a one-line justification near the operation. Futex / wait-wake need `__ATOMIC_SEQ_CST` — WT's default wrappers are relaxed and insufficient; use `__atomic_load_n(..., __ATOMIC_SEQ_CST)` with a comment.
14. **Publish/consume pairs.** Data writes must complete before the flag store (`_release` store, or `WT_RELEASE_BARRIER` before a relaxed store); the consumer loads the flag `_acquire` before reading the data. Mismatched pairs (release store / relaxed load, or relaxed store / acquire load) are bugs on ARM/POWER. Flag-gated paired state: store the value relaxed, then the flag with `_release`; load the flag `_acquire`, then the value relaxed — relaxed on the flag breaks the pair. A single-writer / many-reader publish (`wt_shared bool committed` with `store_bool_release` on the writer, `load_bool_acquire` on the reader) makes subsequent reads of paired fields safe as relaxed loads. Update a timestamp and its `_is_pinned` companion under the same release-store epoch — separated updates leave readers seeing `pinned=stale, value=new`.
15. **Pointer publish.** Store the pointed-to data first (any ordering), then `__wt_atomic_store_ptr_release` the pointer; load with `__wt_atomic_load_ptr_acquire` before dereferencing.
16. **Spin loops.** Every spin loop includes `WT_PAUSE()`; loops that re-read a shared value use at least `WT_READ_ONCE` or a relaxed atomic load, not a plain local read.
17. **Flag bitmask ops.** `flags_atomic` fields use `F_SET_ATOMIC_*` / `F_CLR_ATOMIC_*` / `F_ISSET_ATOMIC_*`, not the non-atomic variants.
18. **Calibrated relaxation.** Eviction-only paths can use relaxed loads on `upd->txnid` / `upd->prepare_state` — narrower concurrency than general reader paths makes acquire overkill. But use `WT_ACQUIRE_READ` on `upd->txnid` before downstream loads of `upd->upd_saved_txnid` — prepared rollback writes the tombstone then marks the update aborted, and a non-acquire load races that transition. Don't extract atomic loads into named locals purely for naming — the snapshot diverges from the next read at the same call site; extract only when reusing the same loaded value.
19. **x86 ≠ ARM.** A pattern that "works" on x86 (TSO) may be broken on weak-memory targets. Reason about C11 guarantees, not observed x86 behaviour. Stable-timestamp lag windows can race commit-ts validation — guard the comparison.
