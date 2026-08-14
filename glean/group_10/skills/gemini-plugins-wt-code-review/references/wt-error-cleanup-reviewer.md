# WiredTiger error-handling and resource-cleanup review

Error-handling macro correctness and resource cleanup in WiredTiger C: `WT_RET` leaking allocations, `WT_TRET` swallowing context, `WT_NOTFOUND` propagated as an error, errors silently dropped in cleanup helpers, unused scratch buffers, redundant NULL guards, missing `WT_CLEAR`, and functions tangling cleanup with real work. All relevant macros are in `src/include/error.h`.

## Scope

In-scope: `+` lines in `*.c` and `*.h`. For each modified function, read enough to see the whole function — what allocations exist before each macro call, what `WT_DECL_ITEM` / `WT_DECL_RET` declarations exist, whether an `err:` label is present, what `ret` holds at the point of use, and what the function returns at the end.

## Reference

Source: `src/include/error.h`. Key macros:

- **`WT_RET(a)`** — returns immediately if `a` non-zero. Does NOT touch `ret` or reach `err:`. Safe only before any resources are acquired.
- **`WT_ERR(a)`** — sets `ret`, jumps to `err:` if non-zero. Use whenever the function owns resources that must be freed on failure.
- **`WT_TRET(a)`** — updates `ret` only if the new error "wins" (`WT_PANIC`, or current `ret == 0 / WT_DUPLICATE_KEY / WT_NOTFOUND / WT_RESTART`). Use in cleanup paths to accumulate errors. Does NOT log.
- **`WT_TRET_MSG(session, a, ...)`** — same winning logic, also calls `__wt_err` when the error wins.
- **`WT_RET_NOTFOUND_OK(a)`** — returns only on non-zero errors other than `WT_NOTFOUND`.
- **`WT_ERR_NOTFOUND_OK(a, keep)`** — jumps to `err:` only on non-`WT_NOTFOUND` errors. `keep=true` leaves `ret == WT_NOTFOUND`; `keep=false` clears `ret` to 0.
- **`WT_RET_PANIC` / `WT_ERR_PANIC`** — call `__wt_panic` then return or jump to `err:`.

`dist/s_void` is the canonical check for a missing `WT_RET` wrapper — fix the call site, don't add suppressions.

## What to flag

### 1. `WT_RET` after an allocation — resource leak

Any `WT_RET(...)` after a heap pointer or `WT_ITEM`/scratch buffer is assigned bypasses `err:` cleanup. `WT_RET` is correct only before any allocation or in a function with no cleanup needs.

- Flag: `WT_RET(__some_helper(...))` after `WT_ERR(__wt_calloc_one(session, &buf))`. Fix: change to `WT_ERR(...)`.

### 2. `WT_TRET` where error context is lost

In cleanup paths needing an always-on diagnostic use `WT_TRET_MSG`; flag extra local error variables used instead of the standard pattern. `WT_TRET_MSG` logs at error level via `__wt_err` — a path that deliberately wants a verbose-level message (`__wt_verbose_warning` for an expected failure) legitimately uses `WT_TRET` plus its own log call, so don't reflexively demand `WT_TRET_MSG`.

- Flag: `if ((ret2 = __wt_some_func(...)) != 0) { ... ret = ret2; }` — fix to `WT_TRET(...)` / `WT_ERR` / `WT_TRET_MSG`.
- Use distinct error variables (`stable_ret`, `ingest_ret`) only when combining genuinely parallel sub-op errors — reusing one `ret` there silently loses the first error.

### 3. `WT_RET` / `WT_ERR` where `WT_NOTFOUND` is expected

- Flag: `WT_RET(__wt_*_lookup(...))` / `WT_ERR(...)` where a miss (lookup, end-of-scan, metadata-scan exit `next()`) is normal — use `WT_RET_NOTFOUND_OK(...)` to ignore silently, or `WT_ERR_NOTFOUND_OK(..., true)` then test `ret == WT_NOTFOUND`.
- Flag: `WT_ERR_NOTFOUND_OK(expr, false)` followed by `if (ret != 0)` — `keep=false` clears `ret`, so the check never fires.
- After `WT_ERR_NOTFOUND_OK(call, true)` that intentionally keeps `ret == WT_NOTFOUND`, explicitly clear `ret = 0` before falling into `err:` — otherwise subsequent ops livelock on the stale NOTFOUND.
- A bool-returning helper that internally calls `WT_ERR_NOTFOUND_OK(..., true)` silently swallows non-`WT_NOTFOUND` errors — convert it to `int` with a `bool *out` parameter so real failures propagate. Same for `WT_IGNORE_RET` on a function that walks open btrees or applies a callback (`__wt_conn_btree_apply`) — make the wrapper `int` and propagate with `WT_RET` even when current callees always return 0.

### 4. Inline `__wt_panic` instead of `WT_RET_PANIC` / `WT_ERR_PANIC`

- Flag: `return (__wt_panic(session, WT_ERROR, "..."));` → `WT_RET_PANIC(session, WT_ERROR, "...");` or `WT_ERR_PANIC(...)` inside a function with `err:`.
- `WT_RET_ONLY(call, WT_PANIC)` is the pattern for "swallow most errors but never `WT_PANIC`" (e.g. eviction-side restore helpers).

### 5. Cleanup helper that always returns 0 — silently drops errors

- Flag: a function ending `return (0);` / `return ret = 0;` that contains `WT_TRET(...)` / `WT_TRET_MSG(...)` earlier or in a `CLEANUP:` / `err:` block. Fix: return `ret`; update the caller to check or propagate it.

### 6. Mixed cleanup and real work / control-flow

- Flag: an `err:` body containing both `__wt_free` calls AND substantive work like `WT_ERR(__wt_log_write(...))` — extract a `__do_the_work(...)` helper.
- Flag: `ret = 0; goto err;` on a success path — `err:` means failure in WiredTiger convention. Use a `done:` success label, or inline/call `__cleanup_resources(...)` directly.
- An `err:` label without an explicit `return (ret)` falls through to whatever follows — adding the label is not enough.
- Cleanup must run in reverse-allocation (LIFO) order when an extracted `__init` / `__deinit` pair exists — forward-order error paths are a latent leak.
- Swapping a foundational macro from `WT_RET` to `WT_ERR` (e.g. `WT_SPIN_INIT_TRACKED`) silently changes every caller's control flow; callers without an `err:` label fall through — include a per-caller audit in the same PR.
- A macro with an output `ret` parameter must initialise `(ret) = 0` at macro entry before any conditional branch — otherwise a path that doesn't assign leaks the caller's prior `ret`.

### 7. Unused scratch buffer declared but never allocated

`WT_DECL_ITEM(name)` with no subsequent `__wt_scr_alloc` / `__wt_buf_*` call, but `__wt_scr_free` at `err:` — both the declaration and the free are dead code. Fix: remove both.

### 8. Redundant NULL guard / missing `WT_CLEAR` / allocation hygiene

- `__wt_scr_free` and `__wt_buf_free` already check NULL — flag `if (tmp != NULL) __wt_scr_free(session, &tmp);`.
- Use `WT_CLEAR(x)` for zeroing stack structs, not `memset(&x, 0, sizeof(x))` / `= {0}` in non-declaration contexts.
- Initialise every local pointer at declaration (`const char *x = NULL`) so the `err:` chain of `__wt_free` calls is safe even when the first allocation fails. After `__wt_buf_grow()` any saved pointer into the buffer's `mem` is invalid — recompute write pointers from the buffer base.
- No assignments inside compound conditions — leaves variables uninitialised on the short-circuit branch.
- Boolean extraction from `WT_CONFIG_ITEM` is `cval.val != 0`, not `cval.len != 0` (`.len` is the value's string length, not a presence indicator). Gate diagnostic-only parameters with `#ifdef HAVE_DIAGNOSTIC` in the signature and delete redundant `if (param != NULL)` guards inside diagnostic blocks.
