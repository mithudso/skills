# WiredTiger API contract review

WiredTiger's public/internal API boundary: test code reaching into internal `_IMPL` structs, signature changes on externally-visible functions, wrong error codes for contract violations, and new methods added where a flag would preserve the existing API surface.

## Scope

In-scope: public header (`src/include/wiredtiger.h.in`), internal headers defining `*_IMPL` structs, test files under `test/suite/`, `test/csuite/`, `test/catch2/`, and implementation files. Skip anything outside those.

## Reference

Error-code conventions:

- `ENOTSUP` — the operation is not supported for this object type, cursor type, or configuration (e.g. `reserve` on a layered cursor).
- `EINVAL` — the caller passed a bad argument: wrong type, out-of-range value, incompatible flags.
- `WT_ERROR` / `WT_PANIC` — internal errors, not API contract violations.
- `EBUSY` — the resource is temporarily unavailable (lock held, transaction active).

## What to flag

### 1. Test code bypassing the public API

Test code in `test/suite/`, `test/csuite/`, or `test/catch2/` that directly accesses a `_IMPL` struct (`WT_SESSION_IMPL`, `WT_CURSOR_BTREE`, etc.) instead of the public API.

- **Signal:** includes `wt_internal.h`, casts `WT_SESSION *` to `WT_SESSION_IMPL *`, or assigns `session->txn.isolation = ...` instead of calling `session->reconfigure(...)`.
- **Fix:** use `session->reconfigure(session, ...)` or add a test-facing config string to `reconfigure`.

### 2. Public-API signature / ABI change

Changes to parameter lists, return types, or names in `wiredtiger.h.in` or the `WT_SESSION` / `WT_CURSOR` / `WT_CONNECTION` / `WT_STORAGE_SOURCE` vtable structs. Also `version` / generation fields added to in-memory (non-on-disk) structs — versioning is for persistence; in-memory should use flags, config strings, or NULL-checked function pointers.

- **Signal:** added param to `(*checkpoint)(...)`; or `uint32_t version;` in `WT_CURSOR_LAYERED` with no on-disk usage.
- **Fix:** add an `_ext` variant for signature changes; use append-only fields + NULL-check dispatch (`pl_complete_checkpoint` vs `pl_complete_checkpoint_ext`). Never shoehorn a value into an old signature via a function-pointer cast — it silently corrupts the call stack; dispatch via NULL-check on the function pointer instead, and guard every call site so alternative implementations can omit the method.
- Adding error codes to `wiredtiger.in` is itself a public-API change (no placeholders / test-only codes / `XXX` even in intermediate commits), and touches `api_err.py`, `api_strerror.c`, `wiredtiger.in`, and `error-handling.dox` together — run `dist/s_all` regeneration and grep all four before review.
- Public flag-bit ranges crossing boundaries must be value-stable: reserve ranges with header comments; avoid `static_assert(FLAG == 0x4, ...)` — it breaks every time the table grows.

### 3. Wrong error code for an API contract violation

Apply the Reference conventions: unsupported → `ENOTSUP`; bad argument → `EINVAL`; internal failure → `WT_ERROR`/`WT_PANIC`; retry-able unavailability → `EBUSY`.

- **Signal:** `WT_ERR_MSG(session, WT_ERROR, "... not supported ...")` → `ENOTSUP`; `WT_ERR_MSG(session, ENOTSUP, "invalid argument")` → `EINVAL`.
- **Fix:** match the error code to the caller's recourse. Enforce strict timing constraints (`early_load`, disagg-only) at the API boundary with `WT_ERR_MSG(... EINVAL ...)`, not a silent no-op when the connection isn't ready. Use `void` return for always-succeed-or-throw functions — `int` is misleading and callers waste cycles checking it.

### 4. New method where a flag would preserve the API surface

New public methods on `WT_CURSOR` / `WT_SESSION` / `WT_CONNECTION` whose sole purpose is an existing method with one behavioral toggle.

- **Signal:** `cursor->reserve_overwrite(cursor)` as a distinct vtable entry; or `open_cursor_v2` when `open_cursor` already accepts a config string.
- **Fix:** add a cursor flag (`WT_CURSTD_OVERWRITE`) or config string (`"new_feature=true"`).

### 5. Headers, extensions, naming, field shapes

- Internal struct definitions for new modules belong in `*_private.h`, never in a header public callers include. Application-visible numeric constants exposed via stats live in `wiredtiger.h.in` as `WT_*` macros (precedent: `WT_BACKUP_FILE` / `WT_BACKUP_RANGE`), not in internal headers.
- Heuristic config knobs must allow disabling (`min=0` or equivalent) so the feature can be field-reverted without a binary roll-back. For pluggable values (extension names like `page_log`), use `type='string'` in `api_data.py`, not `choices=[...]` — `choices` locks out runtime-loaded extensions.
- Loadable extensions (`/ext/`) depend only on `wiredtiger.h` — using `__wt_calloc_one`, `WT_ERR`, `WT_RET`, or `WT_ASSERT` ties them to internal headers; define local `<EXT>_KV_ERR` / `<EXT>_KV_RET` macros and use `calloc`/`free`/`assert`. Extension flag bits must not collide with public PALI flag bits — start internal flag domains past the public range and document the window.
- Public extension APIs use intent-named functions (`load_key`, `get_key`), never implementation-leaking names (`load_key_blob`). Document ownership in public header comments — who allocates, who frees, when memory becomes the other side's responsibility. Two-call APIs (size probe → data fill) must document that the implementer is forbidden from mutating the object between calls.
- Use `union { uint64_t lsn; int error; }` for either-or fields; overloading one `uint64_t` with documentation is rejected. New fields in public PUT/GET arg structs use the proper enum type (`WT_BTREE_STORAGE_TIER`), not raw `uint32_t`; prefer explicit `WT_BTREE_STORAGE_TIER_NONE` over empty-init. An out-parameter `const char *err_msg` is declared `const char *err_msg` at the caller and passed `&err_msg` — casting `char *` via `(const char **)&p` hides a real const-violation. Filesystem utilities must not know WiredTiger naming conventions (`"WiredTiger"` prefix, `.wt` suffix) — that layering belongs in connection/layered code.
