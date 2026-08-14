# WiredTiger assertion and defensive-check review

Covers the invariant-checking surface: every new or modified `WT_ASSERT*`, `WT_*_PANIC*`, `WT_ERR_ASSERT`, `WT_RET_ASSERT`, `WT_VERIFY_*`, and every hand-rolled `if (cond) return err;` that exists to catch a "should-never-happen" condition. Walk each added check against the six principles below.

## Scope

Find the surface by grepping the diff for assertion / panic / defensive-return shapes:

```
WT_ASSERT|WT_(ERR|RET)_(ASSERT|PANIC)|WT_VERIFY|^\+.*if .*return (WT_ERROR|WT_PANIC|EINVAL)
```

Include hand-rolled defensive returns, not just macros. For every match, read the surrounding function so the check can be classified before applying the rules.

## Reference

**Source of truth:** `src/include/error.h` (macro definitions) and `CONTRIBUTING.rst` (assertion policy, forbidden raw `assert()`). Read both at the start of every review — definitions evolve. If a finding can't be tied to one of these or a concrete invariant, don't raise it.

Starter grep: `grep -nE '#define WT_(ASSERT|RET_ASSERT|ERR_ASSERT|.*PANIC)' src/include/error.h`

**Macro summary:**

- `WT_ASSERT(session, cond)` — DIAGNOSTIC builds only; zero cost in release. Default for internal invariants.
- `WT_ASSERT_ALWAYS(session, cond, msg, ...)` — always fires; abort acceptable; use for corruption / data-loss invariants.
- `WT_ASSERT_OPTIONAL(session, category, cond, msg, ...)` — fires when a runtime diagnostic category is enabled; for expensive invariants.
- `WT_ERR_ASSERT` / `WT_RET_ASSERT(session, category, cond, errcode, msg, ...)` — DIAG: abort; release: propagate. Use when a condition shouldn't fail but can legitimately, and you want abort in dev, propagation in prod.
- `WT_RET_PANIC_ASSERT(...)` — DIAG: abort; release: `WT_PANIC`. For unrecoverable invariants.
- `WT_RET_PANIC` / `WT_ERR_PANIC` — always; engine has detected its own corruption.
- Raw `assert()` from `<assert.h>` — **forbidden** (CONTRIBUTING.rst).

Frequency in `src/` (calibration): `WT_ASSERT` ~1,200, `WT_ASSERT_ALWAYS` ~240, `WT_ASSERT_OPTIONAL` ~10. Default to `WT_ASSERT`; escalations are deliberate.

**Assert vs defensive return:** can the condition occur in a correctly-coded system with valid inputs? No → assert. API misuse → `WT_RET_MSG(EINVAL, ...)`. OS/external → defensive return (`ENOMEM`, `EBUSY`, `WT_NOTFOUND`). Unrecoverable → `WT_RET_PANIC` / `WT_ERR_PANIC`. When the call site is internal-only, prefer an assert — `EINVAL` means the API is being misused, an assertion failure means a bug in our own code.

## What to flag

### 1. Effect-preservation — the check must not change program behaviour

`WT_ASSERT` and `WT_ASSERT_OPTIONAL` may evaluate to nothing under release / diagnostics-disabled. Anything the condition does — function call, increment, pointer mutation, lock acquisition, RNG advance — disappears with it. Read the condition as if it were `(void)0` and ask whether the surrounding code still works.

```c
/* WRONG: helper_that_mutates() is compiled out under !HAVE_DIAGNOSTIC */
WT_ASSERT(session, helper_that_mutates(session) == 0);

/* RIGHT */
WT_DECL_RET;
ret = helper_that_mutates(session);
WT_ASSERT(session, ret == 0);
```

### 2. Reachability classification — the macro must match how the condition can fail

- **Programmer error inside the engine** (caller-violated internal contract, dead branch, ordering invariant): assert. Default to `WT_ASSERT`; escalate to `WT_ASSERT_ALWAYS` only if a silent production violation would corrupt data or hand back wrong answers and the check is cheap.
- **API misuse from outside the engine**: runtime error (typically `EINVAL`).
- **OS / allocator / IO / other-process condition**: runtime error matching the failure mode (`ENOMEM`, `EBUSY`, `WT_NOTFOUND`).
- **Engine-detected corruption with no recovery**: `WT_RET_PANIC` / `WT_ERR_PANIC` / `WT_RET_PANIC_ASSERT`.
- **Can-legitimately-fail-but-shouldn't, want abort in dev + propagation in prod**: `WT_ERR_ASSERT` / `WT_RET_ASSERT`.

Two symmetric mistakes: a runtime error guarding an internal-only condition no caller can reach (should be an assert or deleted), and an assert on a condition reachable from outside the engine (should be a runtime error). Reference-count guards like `if (rc == 0 || rc + 1 == 0) return EINVAL` should be `WT_ASSERT(session, rc > 0 && rc < UINT32_MAX)` — overflow and zero-on-add are programmer bugs. Don't assert externalities you don't own — a key-provider mock should not assert LSN-increasing; LSN ordering is WT's contract, not the extension's. Asserts check invariants, not policy — `WT_ASSERT(session, strcmp(name, name_new) == 0)` under `HAVE_DIAGNOSTIC` is fine to cross-check a decision, but if it can fire on a legitimate difference it's the wrong assertion. Use `WT_ASSERT_ALWAYS` for corruption / data-loss invariants (queue-drain on shutdown, invalid live-restore bitmap, turtle-metadata parse failure); plain `WT_ASSERT` under `#ifdef HAVE_DIAGNOSTIC` for debug-only verification.

### 3. Condition shape — the predicate must say what it means

- Assert the *positive* invariant — the thing that must be true. Inverting it changes the failure semantics.

```c
/* WRONG: polarity inverted */          /* RIGHT: bulk cursors must not be cacheable */
WT_ASSERT(session, cacheable || !bulk);  WT_ASSERT(session, !cacheable || !bulk);
```

- One condition per assert when the arms protect unrelated paths; when `A || B` fires you can't tell which arm — split on mutually-exclusive paths: `if (!__wt_ref_is_root(ref)) WT_ASSERT(session, WT_REF_GET_STATE(ref) == WT_REF_LOCKED);`. Don't OR across mutually-exclusive states — `WT_ASSERT(session, WT_REC_EVICT | WT_REC_HS)` is bitwise nonsense; use separate asserts per branch.
- Each branch of a mutually-exclusive multi-way decode needs its own precondition assertion set (e.g. three prepared-info branches: `temp_stop_ts == WT_TS_NONE` for both-prepared, `temp_durable_stop_ts != WT_TS_NONE` for stop-only). Missing one branch is how misinterpretation slips through.
- Test the invariant directly, not an incidental consequence. Use sentinel macros in conditions — `ts == WT_TS_NONE`, not `ts == 0`; `!WT_TIME_WINDOW_HAS_PREPARE(&tw)`, not `... == 0` — the macro form survives sentinel-value changes.

### 4. Information content — a failed check must give the next reader something to work with

- Always-on macros (`_ALWAYS`, `_OPTIONAL`, `_ERR_ASSERT`, `_RET_ASSERT`, `_PANIC*`) take a message; use it. Include the offending value via `%` formatting and state the invariant in words, not "Assertion failed".
- When the intent of an assert is not obvious from local context, a one-line comment pays for itself. A comment near the check and the check itself must agree.

### 5. Non-redundancy — a check must add information

- Don't assert a condition the very next statement already checks by crashing (asserting a pointer non-NULL immediately before dereferencing it — the deref crashes anyway; this includes `strcmp`).
- Don't guard a call against an input the callee already handles (`NULL`-checking before `__wt_free` / `__wt_scr_free`).
- Don't write a check whose condition can't fire because of code immediately above it.
- Don't hand-roll a check an existing macro expresses (`WT_RET_MSG`, `WT_ERR_MSG_CHK`, `WT_RET_PANIC`, `WT_ERR_NOTFOUND_OK`).
- Feature-gate asserts (`WT_ASSERT(session, <feature_flag>)`) become noise once the flag is permanent — remove them.

### 6. Coverage — load-bearing invariants should be checked when cheap

- A function callable only while holding a specific lock should assert ownership at entry with `WT_ASSERT_SPINLOCK_OWNED` (or `__wt_spin_owned` in a `WT_ASSERT_ALWAYS`) — add it whenever you split out a helper that assumes locked entry.
- A function valid only in a specific role / mode / phase (leader vs follower, during checkpoint, on a particular dhandle) should assert that context at entry. Assert isolation at entry (`session->txn->isolation == WT_ISO_SNAPSHOT`) on snapshot-only paths rather than re-checking inside loops.
- Push feature-flag gating to the highest-level entry point and use entry asserts in deeper helpers — not duplicated `if (!flag) return 0` guards in every callee.
- Don't delete an invariant assertion to accommodate a new path — condition it: `WT_ASSERT(session, F_ISSET(btree, WT_BTREE_DISAGGREGATED) || <old invariant>)` keeps the non-disagg invariant enforced.
- For role-asymmetric invariants prefer the flat form `WT_ASSERT_ALWAYS(session, !leader || <leader-only-check>, ...)` over `if (leader) WT_ASSERT(...)` — keeps the condition in one line and avoids skipping `__wt_buf_fmt` / lookup setup on followers where it's still needed.
- A shared field an assert reads should be read with the appropriate atomic load helper (the ordering analysis itself belongs to `wt-concurrency-reviewer.md`).
