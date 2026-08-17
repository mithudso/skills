# WiredTiger style and idiom review

Style, idiom, and comment-quality violations — not bugs, refactors, or design. Source of truth is `CONTRIBUTING.rst` at the repo root: read it once at the start of every review; the rules evolve. Trust `dist/s_clang_format` for operator spacing and brace placement within a line — do not flag what it handles. If a finding can't be tied to `CONTRIBUTING.rst` or a `dist/s_*` script, don't raise it.

## Scope

In-scope: changed lines in `*.c` and `*.h` under `src/` and `test/csuite/`. Skip generated files — anything between `AUTOMATIC ... GENERATION START/STOP` markers, files under `dist/` produced by the Python scripts, and `src/include/*_auto.h`.

## What to flag

### 1. Naming
- `__wt_` (extern) / `__wti_` (internal-linkage-but-shared) / `__` (static) function prefixes.
- Output parameter `p` suffix (`resultp`, not `result`).
- Canonical short names (`session`, `conn`, `btree`, `ref`, `page`, `cbt`).
- `WT_*` macros, `_IMPL` / `WTI_` typedefs.
- Enum-as-state-machine: use enums (not `#define`) for 3+ states so `-Wswitch-enum` catches missing cases; two-valued `typedef enum` should be `bool`. Never add a `default:` to an exhaustive switch.
- Avoid `<` / `>` on enum values — extract a named predicate macro instead.
- No stdlib name shadowing (`time`, `order` → `checkpoint_time`, `checkpoint_order`).
- Positive-semantics boolean parameters: rename `ok_not_exist`/`no_error` to `must_exist` or split the function.
- Named sentinel constants over literal `0` / `0xff...` (`WT_BLOCK_INVALID_PAGE_ID`, `WT_TS_NONE`, `WT_TS_MAX`, `WT_PREPARED_ID_NONE`) — sentinel values can change and named forms are grep-able.

### 2. Declarations
- Block ordering: `struct` → `WT_*` alpha → declaration macros → built-ins by type, alpha within.
- `WT_DECL_RET`, scope, init-at-declaration, blank line after the block.
- `#ifdef HAVE_DIAGNOSTIC`-only variables must be declared inside the `#ifdef` block, not at function top.
- Compile-time constants: `constexpr static`, not `const`; `constexpr const` is redundant.

### 3. Comments

Walk every added or modified comment — lines starting with `/*`, ` *`, `//`, or inside `/*…*/` blocks. Quote the verbatim offending text and suggest a concrete replacement.

- Banner-form function headers (`/* __wt_foo -- ... */`) on every new exported function.
- No `//` comments in production code; use block style with opening/continuation/closing on their own lines.
- Full sentences ending in periods.
- **WHY content:** explain *why*, not *what*. If a reader can reconstruct what the code does just by reading it, the comment should explain why instead.
- **Staleness:** when a comment names a specific identifier, grep for it; if absent, flag as likely stale and describe the concept instead.
- **FIXME format:** `FIXME-WT-XXXXX:` (with colon) inside the function body on its own line, never in the doxygen header above it (that trips `s_all`). `TODO` → `/* FIXME-WT-NNNNN: <description>. */`. No closed Jira/PR references.
- **`\c` Doxygen markup** only in `dist/api_data.py` / Doxygen headers, never in `src/*.c` comments.
- **Cross-module docs** describe the contract, not internal function names on the other side ("must match what the disagg read path expects", not "what `__block_disagg_read_multiple` expects").

### 4. Error handling
- Use `WT_RET`/`WT_ERR`/`WT_TRET`/`WT_RET_NOTFOUND_OK` family — not raw `if (ret != 0) return ret;`.
- `ret` variable name; `if (0) { err: }` pattern (choice between macros is `wt-error-cleanup-reviewer`).
- Assertion choice: `WT_ASSERT` vs `_ALWAYS` vs `_OPTIONAL`. No raw `assert()`.

### 5. Memory
- `__wt_calloc`/`__wt_malloc`/`__wt_free` only (no raw `malloc`/`free`).
- `__wt_calloc_def`/`_one` shortcuts.
- Scratch-buffer pattern. No manual `= NULL` after `__wt_free`.

### 6. Flags
- `F_SET`/`F_CLR`/`F_ISSET` (and `FLD_*`, `LF_*`, `_ATOMIC_*` variants). No direct `flags |= mask`.
- No manual hex values inside auto-generated flag blocks.
- Reserve flag-bit ranges with a header comment; new flags at end of range with next free bit documented.

### 7. Logging & stats
- `__wt_verbose` with correct `WT_VERB_*` and level.
- `WT_STAT_CONN_INCR` vs `WT_STAT_DSRC_INCR`.
- No raw `printf` in library code.

### 8. Formatting
- 100-char line limit, 2-space indent, no tabs, no trailing whitespace.
- Single-statement bodies without braces.
- `return (value);` with parens.
- `/* FALLTHROUGH */` for intentional switch fallthrough.
- `s_all` requires the comma before a `#ifdef`-guarded parameter on its own line.

### 9. Forbidden / idiom
- `strcpy`, `strcat`, `sprintf`, raw `assert`, `//` comments, direct `__atomic_*` builtins, closed-ticket references.
- `WT_PREFIX_MATCH(uri, prefix)` instead of `strncmp(uri, prefix, strlen(prefix)) == 0` — correct on truncated strings and consistent.
- `WT_STRING_LIT_MATCH("LITERAL", buf, len)` for fixed-string compares; replace `strcpy`/`strcmp` with `__wt_strndup`/`snprintf` or a const-pointer-returning helper.

### Illustrative examples

Function header — wrong / right:

```c
// Free the page.
void
__wt_page_out(...)
```

```c
/*
 * __wt_page_out --
 *     Discard an in-memory page, freeing all its associated memory.
 */
void
__wt_page_out(...)
```

Declaration ordering:

```c
WT_BTREE *btree;
WT_CONNECTION_IMPL *conn;
WT_CURSOR *cursor;
WT_DECL_RET;
uint64_t now, start, timeout;
const char *key, *value;
bool release;
```

If a likely bug shows up here, mention it once and move on — it belongs to `wt-correctness-reviewer.md`. Don't review generated files.
