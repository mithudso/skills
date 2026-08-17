# WiredTiger hot-path performance review

Performance anti-patterns reviewers flag in PRs: allocations smuggled into allocation-free hot paths, per-page atomic RMWs turning a cache line into a bottleneck, cheap predicates buried after expensive ones, and magic thresholds that hide tuning reasoning.

## Scope

In-scope: `+` lines in `*.c` and `*.h`. For each modified function, decide whether it lives on a hot path — eviction, reconciliation, metadata parsing, checkpointing, or inside a per-page / per-record loop. Only apply the checks below to those.

## What to flag

### 1. Heap allocation in an allocation-free hot path

Flag any `__wt_malloc` / `__wt_calloc` / `__wt_realloc` / `__wt_scr_alloc` / `WT_DECL_ITEM` inside a per-page loop, metadata parsing helper (`src/meta/`), or eviction fast path (`src/evict/`). Allocations add latency jitter and stall under memory pressure. Fix: hoist before the loop, or restructure to check the common/cheap format first and allocate only on the slow path.

### 2. Per-page or per-record contended atomic RMW

Flag `WT_STAT_*_INCR` or `__wt_atomic_add` inside `for`/`while` loops walking pages/records with concurrent workers (sweep, eviction, verify, RTS). Hot atomic counters become serialization points as all workers fight one cache line. Fix: accumulate in a per-thread local (`uint64_t pages_walked = 0; ++pages_walked;`) and flush with `WT_STAT_*_INCRV` at scan end or checkpoint boundary. For aggregate cross-btree counters, reset at checkpoint start and recompute by walking handles — scattered per-call-site CAS updates miss untouched handles.

### 3. Expensive predicate before cheap predicate

Flag `strcmp` / `strncmp` / `__wt_config_getones` / hash lookup / non-trivial call placed before a cheaper integer comparison, flag test, or NULL check that could short-circuit it. Fix: cheapest test that eliminates the most cases runs first — e.g. compare cheap `uint64_t` fields (`order`, `time`) before `strcmp(checkpoint_name, ...)`.

Related config cost: reparsing config in extracted hot-loop helpers is measurable on schema-heavy workloads. Extract boolean values once via `cval.val` (O(1)); pass either the raw `cfg[]` or a pre-parsed struct, never both (duplicate parsing hides as convenience). Don't extract a single-use atomic load into a named local purely for naming — the snapshot diverges from the next read; extract only when reusing the same loaded value. Use `__bit_ffs` / `__bit_ffc` for bitmap searches instead of bit-by-bit iteration.

### 4. Magic threshold without a named constant

Flag raw integer literals used as thresholds (eviction aggressiveness, cache fill %, retry counts, step sizes) without `#define` / `static const` / explanatory comment — units and reasoning are invisible. Fix: `#define WT_CACHE_EVICT_TRIGGER_PCT 90` and use the constant; keep tuning rationale in a comment near the definition. Also flag threshold *changes* with no stated why.

### Benchmarking note

For state-mutating ops (insert/update/modify), benchmark with instruction counts (`perf_event_open` / `rdpmc`), not wall-clock averages — update chains lengthen and trigger eviction, so wall-clock is noisy.
