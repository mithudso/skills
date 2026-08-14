# WiredTiger verbose logging and statistics naming review

Logging mistakes reviewers commonly flag: wrong category, wrong level, misleading stat names, missing outcome context, multi-line macro misuse. See `src/include/verbose.h` for categories and `src/support/stat.c` / `src/include/stat.h` for stat names. Walk each added or modified `__wt_verbose*` call and each new stat name against the checks below.

## Scope

Collect every added or modified `__wt_verbose*` call and every new stat name (in `src/support/stat.c` or `src/include/stat.h`). Apply the checks below to each.

## Reference

Verbose levels (lowest to highest severity):

- `WT_VERBOSE_ERROR` — actual errors that should always appear in the log.
- `WT_VERBOSE_WARNING` — something unexpected but recoverable.
- `WT_VERBOSE_NOTICE` — rare, operationally important events (e.g. recovery complete).
- `WT_VERBOSE_INFO` — routine informational events.
- `WT_VERBOSE_DEBUG_1` — default debug level; high-frequency paths should avoid this.
- `WT_VERBOSE_DEBUG_2` … `DEBUG_5` — increasingly detailed, for infrequent events.

`grep -nE 'WT_VERB_' src/include/verbose.h` for all categories.

## What to flag

### 1. Wrong `WT_VERB_*` category

Does the category match the subsystem? Wrong category adds noise and hides the message from where operators look.

- Prefer the operation's category (read during verify → `WT_VERB_READ`) unless the message is specifically about the operation.
- Common mismatches: read/`WT_VERB_VERIFY`, eviction/`WT_VERB_CHECKPOINT`, new subsystem/`WT_VERB_DEFAULT`, disagg/generic.
- **Category staleness:** a category named for one subsystem (`WT_VERB_CHECKPOINT_CLEANUP`) is wrong once the feature is shared (with compact/eviction) — pick a category that spans every caller, or split it.

### 2. Wrong verbosity level

Is severity proportional to event frequency and the operator attention needed?

- `ERROR`/`WARNING` only for actual errors — never retries, misses, or contention, and never at `DEBUG` (a failure logged at debug is invisible to customers).
- `NOTICE` for infrequent, significant events. `DEBUG_1` is the default; use `DEBUG_2+` for frequent paths, `DEBUG_3+` for per-page. Routine file enumeration is `DEBUG_1`, not `WARNING`.

### 3. Misleading or unjustified stat names

- A stat name containing `error` / `fail` / `err` for a normal event (retry, CAS race, contention, expected miss) — use neutral names (`*_race`, `*_retry`, `*_contention`, `*_miss`). Check new `WT_STAT_*` in `stat.c` / `stat.h`.
- Precise full words and consistent abbreviations (`URI` not `uri`); names persist forever.
- Justify new counters — add only with a distinguishable code path and a use case today, not "for future use".
- Fire the stat at the actual work site, not in a decision helper — otherwise it doesn't match observable behaviour.

### 4. Missing outcome context / format-string hygiene

- A log message that describes an outcome without distinguishing success from failure: `"checkpoint cleanup completed..."` — pass outcome as `%s`: `"completed (%s)...", ret == 0 ? "success" : "error"`.
- Never embed content in the format string: `"%s", str`, not `"content%s", "!"`; no space-before-colon (`"failed: X"`); never `"%" PRIu64 ""`.
- Include file/object identifiers (`fh->name`) in per-file messages — generic strings are useless in multi-file dumps.
- Error messages name the state, not the inability to check it: `"Ingest contains dirty content"`, not `"cannot be verified"`. Recovery-transcript logs say what was missing (`"colgroup missing"`, `"file missing"`).
- No multi-line string concatenation across `__wt_verbose_*` args (`"phase " NAME " completed" ...`) — use format args: `"phase %s completed in %"...`.

### Progress-logging mechanics

Progress logs fire on three OR'd triggers: `closing`, `__wt_counter_backoff()`, and `time_diff / PERIOD > msg_count`. Backoff alone starves slow ops; time alone floods fast ones. `msg_count` increments with every log — it both tracks logs and rate-limits the time trigger, so a missing increment breaks the time fallback.
