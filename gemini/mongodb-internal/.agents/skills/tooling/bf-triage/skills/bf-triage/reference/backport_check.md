# Backport-presence check

Run this after Step 6 (Bisect) on **every** BF whose
`Evergreen Project` matches `mongodb-mongo-v8.*`,
`mongodb-mongo-v7.*`, or any other non-master release branch
(`sys-perf-8.*` included). Skip for `mongodb-mongo-master` /
`sys-perf` BFs — those have no backport to look for by definition.

## Why

Hot release-branch BFs are often the visible tail of a fix that
already landed on master and never got cherry-picked. Without the
check, the report recommends "Reproduce" or "Re-route" when the
correct disposition is "Backport `<sha>` from master".

## Procedure

1. Identify the failing test path / failing assertion source line
   from Step 6's `git_blame`.
2. Run a `git_log` (gateway or local fallback) on `10gen/mongo`
   master scoped to that path, **bounded by the failing branch's
   fork-point**:

   ```bash
   FORK_POINT=$(git -C "$MONGO_REPO" merge-base origin/master origin/v8.0)
   git -C "$MONGO_REPO" log "$FORK_POINT..origin/master" \
     --oneline -- <failing-source-path>
   ```

   Gateway equivalent: `git_log repo=10gen/mongo ref=master path=<failing-source-path> since=<fork-point-date>`.

3. For each commit in the result, read its message via `git_show`
   (or gateway `git_show`). A commit qualifies as a candidate
   backport if **any** of:
   - The message references the same `SERVER-XXXXX` ticket that the
     BF's standup comment or sibling-BF analysis named.
   - The message contains the same assertion text the BF cites
     (e.g. "fix invariant in admission ticket release").
   - The message touches the exact source line that `git_blame`
     attributed the failing assertion to.

4. If a candidate is found AND it is NOT present on the BF's
   release branch (`git -C "$MONGO_REPO" log origin/v8.0
   --grep='<commit-message-or-SERVER-key>'` returns nothing), the
   report's `Recommended next steps` includes a top-priority
   **Backport** entry naming the candidate SHA and the
   `SERVER-XXXXX` ticket. Set
   `disposition = keep_open_pending_fix` and
   `fix_location = mongo_server` (or `mongo_test` per the patch).

5. If no candidate is found, add a one-line note to the Appendix:
   `Backport-presence check: ran 'git log <fork-point>..master --
   <path>'; no commit references the failing path. Backport is
   not the explanation.`

## Tagging in the report

- The Recommended-next-steps Backport entry MUST cite the candidate
  SHA, the master-side commit date, the touched path, and the BF's
  branch (e.g. `v8.0`).
- Do NOT recommend the backport PR itself — the skill is read-only.
- The Appendix MUST list either the executed
  `git log --all --grep='<SERVER-key>'` + `git show <release-base>:<path>`
  calls (with their findings), or a single line stating why the check
  was skipped (e.g. `Skipped — no prior closed BF with a Fix Revision`).
  A report that omits both is non-conformant.
