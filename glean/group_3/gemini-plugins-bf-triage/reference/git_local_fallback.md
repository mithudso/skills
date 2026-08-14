# Reference: Step 6 fallback — local history when gateway `git_*` is unreachable

> Read this file only when entering the Step 6 fallback branch (any
> gateway `git_*` call has failed after the Hard rule 5 retry). The main
> `SKILL.md` Step 6 section points here.

Enter this branch automatically whenever a gateway `git_*` call fails
after the Hard rule 5 retry. **No user permission is required** — git
data is always recoverable from local clones, unlike `bb_*` / `evg_*`.
The skill should silently switch to local for the remainder of Step 6
and continue.

1. Resolve a local clone. Use the values resolved in Step 1
   (`$MONGO_REPO_PATH` / `$DSI_REPO_PATH`). If Step 1 was not run yet
   (e.g. the failure was the very first gateway `git_*` call), run the
   Step-1 resolution algorithm now. The resolution order is:
   - Already-set `$MONGO_REPO_PATH` / `$DSI_REPO_PATH` env vars.
   - The user's working repos (origin URL contains `10gen/mongo` /
     `10gen/dsi`) — these are the canonical local sources when present.
     Use **read-only** commands only; never touch HEAD or the index
     (see step 4 below).
   - A `setup_repos.sh` scratch clone (typically under
     `/tmp/bf-triage-workdir-<random>/`) when the user's working
     copies are absent. **Pass a selection arg** based on the BF's
     `Evergreen Project` (from Step 2, which has already run by the
     time Step 6 fires): if the project is `sys-perf*` /
     `*-dsi*` / DSI-flavoured, run `setup_repos.sh both`; otherwise
     run `setup_repos.sh mongo` to skip the DSI clone (~tens of MB +
     a few seconds of network) since Step 6 will only ever consult
     `10gen/mongo` for that BF. The default `setup_repos.sh` with no
     arg still clones both — preserved for backward compatibility
     and for callers that haven't read the BF's project field yet.
2. Ensure the SHA is present in the local object DB:
   `git -C "$REPO" cat-file -e <sha>^{commit}`. If it returns non-zero,
   fetch it: `git -C "$REPO" fetch --depth=1 origin <sha>`. Use
   `--unshallow` only with explicit user approval (it can be slow).
3. Inspect the commit read-only — pick by need:
   - Single-file or full-diff inspection (preferred):
     - `git -C "$REPO" show <sha>:<path>` — content of one file at that SHA.
     - `git -C "$REPO" log -p -1 <sha>` — full diff + commit message.
     - `git -C "$REPO" blame <sha> -- <path>` — blame as of that SHA.
   - Multi-file `Grep` / `Read` over the tree at that revision:
     `git -C "$REPO" worktree add /tmp/bf-scratch-<sha> <sha>`,
     do the reads under that path, then clean up:
     `git -C "$REPO" worktree remove --force /tmp/bf-scratch-<sha>`.
     Always remove the worktree before finishing the report — leftover
     worktrees keep refs alive and confuse the user's `git status`.
4. Forbidden when `$REPO` is the user's working repo (i.e. anything
   resolved from steps 1.1 / 1.2 / 1.3 above — NOT a scratch clone):
   `git checkout`, `git checkout -b`, `git reset`, `git switch`,
   `git stash`, any index or HEAD-mutating command. Those would clobber
   the user's WIP. `git checkout <sha>` is only permitted inside a
   `setup_repos.sh` scratch clone (a `/tmp/bf-triage-workdir-*`
   directory), where you must `git switch -` and `git branch -D` (if
   any was created) before returning control. To tell the two apart,
   check whether the resolved `$REPO` path starts with the scratch-clone
   prefix `/tmp/bf-triage-workdir-`; if not, treat it as the user's
   working repo.
5. In the report, tag any commit citation produced via this fallback as
   `(local fallback)` next to the SHA so the engineer can see which
   evidence did not come from the gateway. The full "Limited evidence"
   banner is NOT required just because git fell back — only `bb_*` /
   `evg_*` outages trigger that banner (the banner says "no fresh `bb_*`,
   `evg_*`, or `git_*` data was fetched", which is no longer true when
   git work succeeded locally). Use a one-line note instead: "Note:
   gateway `git_*` was unreachable; commit evidence below is from local
   clones (tagged `(local fallback)`)."
6. Cleanup. Before writing the report, ensure no `/tmp/bf-scratch-*`
   worktrees remain (`git worktree list` + `git worktree remove --force`
   any stragglers). If a `setup_repos.sh` scratch clone was created
   solely for this fallback, run the script's `cleanup` subcommand
   unless `--keep-clones` was passed.
