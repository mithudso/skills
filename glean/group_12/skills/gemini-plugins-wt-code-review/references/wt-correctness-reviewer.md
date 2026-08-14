# WiredTiger correctness review

Logic bugs and undefined behaviour that escape pre-PR review — visibility/timestamp semantics, lifecycle order, and general C UB.

## Scope

Classify the change by which subsystem(s) it touches: `src/cursor/`, `src/txn/`, `src/checkpoint/`, `src/btree/`, `src/conn/conn_layered.c`, `src/rollback_to_stable/`, `src/session/`, etc. Each subsystem has its own dominant failure modes.

Apply the diff-triage order below to prioritise reading (earlier categories have higher bug density):

1. Visibility / timestamp semantics — any new use of `__wt_txn_visible*`, read of `txn_global->oldest_timestamp`, new RTS comparison, new disagg / layered code path.
2. Lifecycle order — session close, cursor reset-before-reuse, snapshot release before iterating `txn->mod[]`, page-modify init/mutate/set order.
3. General C / UB — integer underflow on unsigned, format-string mismatches, `sizeof` of pointer, switch fallthrough, uninitialized locals before `goto err`.

## Reference

Verify against current source before reporting (names durable; line numbers drift). Useful greps:

- Visibility helpers: `grep -nE '__wt_txn_visible' src/include/txn_inline.h`
- Timestamp constants: `grep -nE '^#define WT_TS_' src/include/txn.h`
- Page-modify discipline: `grep -nE '__wt_page_modify_(init|set)' src/include/btree_inline.h`
- RTS cutoff: `grep -rn 'upd_durable_ts' src/rollback_to_stable/`

## What to flag

### 1. Visibility, timestamps, storage-engine semantics

- **Visibility helper.** `__wt_txn_visible` answers per-txn "can this reader see this update?"; `__wt_txn_visible_all` answers "is this obsolete to every possible reader?" Cleanup / eviction / RTS / HS-pruning need `_all`; read paths need the per-txn variant.
- **RTS comparison.** RTS aborts when `rollback_timestamp < upd->upd_durable_ts`. Any comparison against `upd_start_ts` or `upd_commit_ts` is wrong.
- **Timestamp scope.** Inside checkpoint or on a disagg follower, reads of `txn_global->{oldest,stable,pinned}_timestamp` are suspect — use `checkpoint_*_timestamp` or `disaggregated_storage.last_checkpoint_*_timestamp`. Prefer `__wt_txn_pinned_stable_timestamp()` over a hand-rolled "checkpoint running ? checkpoint_ts : stable_ts".
- **Stable-timestamp fallback.** `__wt_get_stable_timestamp()` falls back to `recovery_timestamp` — gate on `has_stable_timestamp`, not `stable_timestamp == WT_TS_NONE`.
- **After restart, txn ids are unreliable** for obsolete-eligibility — use `write_gen < S2C(session)->base_write_gen`.
- **Keep visibility-helper inputs transparent.** Conditional `prepare_ts` vs `start_ts` accessors must not flow into `__wt_txn_visible` — branch at the caller, not inside the field access.
- **Disagg leader guard.** New mutation paths on disagg-eligible btrees need `!F_ISSET(btree, WT_BTREE_DISAGGREGATED) || layered_table_manager.leader`. Compare sibling writes in `src/cursor/cur_layered.c`.
- **Layered tombstone.** Layered ingest tables use `__clayered_deleted_encode`; `cursor->remove()` via a plain file cursor corrupts the format.
- **`WT_CONFIG_ITEM`.** `.str` is not null-terminated — use `WT_CONFIG_MATCH` / `WT_STRING_MATCH`.
- **Timestamp constants.** Prefer `WT_TS_NONE` over `0`, `WT_TS_MAX` over `UINT64_MAX`. Guard timestamp sources against `WT_TS_NONE` before min/max aggregation — a `0` treated as a valid min produces false positives in parent/child validation.

Time windows and prepared cells:

- **Adding a field to `WT_TIME_WINDOW`** requires touching `INIT`, `IS_EMPTY`, `WINDOWS_EQUAL`, and a `HAS_*` macro if sentinelled (equality is the most-missed). Use `HAS_START_PREPARE` / `HAS_STOP_PREPARE` at decision points — "both ends prepared" uses `HAS_START` because a prepared `stop_ts == WT_TXN_MAX`.
- **`*_prepare_ts` and `*_prepared_id` are paired** — `WT_PREPARED_ID_NONE` is distinct from `WT_TXN_NONE`; don't init prepared-id fields with `WT_TS_NONE`. Never write zero prepare fields, and set `WT_CELL_TS_DURABLE_*` when packing the field — losing the flag corrupts unpacking.
- **Feature-gate prepared_id on `WT_CONN_PRESERVE_PREPARED`**, never `WT_BTREE_DISAGGREGATED` — independent flags.
- **Temp `stop_ts` must initialise to `WT_TS_MAX`**, not `WT_TS_NONE` — defer interpretation until all fields are read.
- **Non-timestamped HS updates** (`start_ts == stop_ts == WT_TS_NONE`) must skip delete-and-reinsert — otherwise they spuriously return EBUSY.

Cells, pages, reconciliation:

- **`WT_CELL_FOREACH_ADDR` / `WT_CELL_FOREACH_DELTA_*` iterate cells, not key/value pairs** — divide by 2 for row-leaf entry counts.
- **Latest-wins k-way merge of delta streams must iterate backwards** (highest index = newest) with a *signed* index type so the `-1` sentinel works — `size_t` / `uint32_t` silently breaks termination.
- **Write-conflict vs cursor-positioning over a truncate list use inverted visibility predicates** — write-conflict tests entries *not visible*, positioning tests *visible*. Easy to flip during refactors.
- **Page-state TOCTOU:** push read-and-decide into `__wt_evict_copy_page_state` — don't read, branch in the caller, then write back.

Metadata, recovery, verify:

- **For known exact metadata URIs** (e.g. `file:foo.wt`), use `__wt_metadata_search` directly — faster than a cursor scan and avoids the cursor-loop error-handling failure mode.
- **Checkpoint identity is `(order, time)`, not name** — name-based discard decisions corrupt data on followers.

### 2. Lifecycle order

- **Session close** (`__wt_session_close_internal`): clear `WT_SESSION_CACHE_CURSORS` → close cursors → rollback txn → release snapshot → close metadata → `__wt_hazard_close` → teardown. A `goto err` before `__wt_hazard_close` leaks hazard pointers.
- **Txn rollback:** snapshot released *before* iterating `txn->mod[]`. New code walking `txn->mod[]` while holding the snapshot pins data unnecessarily.
- **Page modify** (`__wt_page_modify_init` / `__wt_page_modify_set`): init → mutate → set. Calling `_modify_set` early races the eviction thread.
- **Cursor reuse.** Needs `reset` first — position state (`WT_CBT_ACTIVE`, key/value SET flags) leaks across uses.

### 3. General C / undefined behaviour

- **Unsigned underflow.** `v - 1` where `v` is unsigned and can be `0` wraps to `UINT*_MAX`.
- **64-bit truncation.** Timestamps, txn ids, LSNs must stay `uint64_t` — flag storage in `uint32_t` / `int` / `unsigned`. Use `header->header_size` for payload offsets and checksum ranges, not `sizeof(struct)` — downgrades from newer formats read garbage.
- **Off-by-one.** `<=` vs `<` against array bounds; byte length vs element count; inclusive vs exclusive ranges.
- **Format specifiers.** `%d` for `uint64_t`, `%lu` for `size_t` — use `PRIu64`, `WT_SIZET_FMT`.
- **`sizeof(p)` of a pointer** where struct size was intended (e.g. `__wt_calloc(1, sizeof(p))`).
- **Macro double evaluation.** `WT_MIN(x++, y)` — confirm each argument is evaluated once.
- **Side effects in `WT_ASSERT`.** Compiled out in non-diagnostic builds.
- **`switch` fallthrough** without `/* FALLTHROUGH */` or `__attribute__((fallthrough))`.
- **`WT_DECL_ITEM` before alloc.** NULL until `__wt_scr_alloc`; field access before that segfaults.
- **Strict-aliasing cast.** `(uint64_t *)char_ptr` is UB — use `memcpy` or WT pack/unpack helpers.
- **64-bit parses.** Use `strtoll()` / `strtoull()` for offset/size parses — `strtol()` truncates on 32-bit-`long` platforms.
