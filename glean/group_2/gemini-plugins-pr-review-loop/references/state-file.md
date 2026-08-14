# State File Reference

Path: `tmp/pr-review-loop-state.json` (ensure `tmp/` is gitignored in your repo)

Stores state for the current branch only. On termination, the entire file is deleted.

## Schema

```json
{
  "branch": "<TICKET-ID>",
  "running_since": "2026-03-20T10:00:00Z",
  "last_push": {
    "sha": "abc123",
    "files_changed": ["Foo.java", "Bar.java"],
    "failures_fixed": ["task_id_1", "task_id_2"],
    "reviews_addressed": [3184327788, 3184327797],
    "flagged_items": [
      3184327799,
      { "comment_id": 3184327800, "decision": "go ahead with the refactor" }
    ],
    "timestamp": "2026-03-20T10:00:00Z"
  }
}
```

When `last_push` is null (no push has happened yet), `reviews_addressed`, `flagged_items`, and `timestamp` are persisted at the top level of the state file instead — same field names, same shapes:

```json
{
  "branch": "<TICKET-ID>",
  "running_since": "2026-03-20T10:00:00Z",
  "last_push": null,
  "reviews_addressed": [3184327788],
  "flagged_items": [3184327799],
  "timestamp": "2026-03-20T10:00:00Z"
}
```

Once a push occurs, these fields move into `last_push` and the top-level copies should be removed.

`flagged_items` — comments that required a human decision and were surfaced in the terse summary. Each entry is either a bare comment ID (no decision yet) or an object `{ "comment_id": N, "decision": "..." }` (decision recorded). Two ways to provide a decision:
- **In conversation**: tell the parent agent after seeing the FLAGGED summary — it writes the decision to this field immediately.
- **On the PR**: reply to the flagged thread on GitHub — the next Worker reads the reply as the decision.

## Lifecycle

1. **On start**: Read the state file. If `running_since` is set and less than 30 minutes old, **STOP** — another invocation is still running. If older than 30 minutes, treat as stale (previous run crashed) and proceed.
2. **On start**: Write `running_since` with the current timestamp.
3. **During Phase 3/4**: If `last_push` is not null, check `last_push.failures_fixed` — skip failures already handled. Check `last_push.files_changed` — if a *new* failure's stack trace points to a file this skill recently modified, flag it as a possible regression and do NOT auto-fix. Report it for human review instead. If `last_push` is null, treat all lists as empty (first invocation — nothing to de-duplicate).
4. **During Phase 5**: Read `reviews_addressed` and `flagged_items` from `last_push` if it is non-null, otherwise from the top level of the state file. Skip comments already in `reviews_addressed`. For each flagged item, check for a recorded `decision` or a new PR author reply; act on it and remove from `flagged_items`. Also check the GitHub API for existing replies (belt and suspenders).
5. **On push**: Update `last_push` with the new SHA, files, task IDs, comment IDs, and any remaining `flagged_items`.
6. **On exit (only if this invocation acquired the lock)**: Clear `running_since`. Do NOT clear if this invocation STOPped in step 1 because another invocation owns the lock.
7. **On termination (green build + all reviews addressed)**: Delete the state file entirely.
