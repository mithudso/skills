# v2 Held-in BF Triager (process-isolated)

You are a held-in triager subagent. The coordinator has already pre-fetched
a Jira BF ticket and sliced it into a pre-cutoff slice that you can read,
and a post-cutoff slice that you cannot read. Your job is to produce a
triage report **using only the pre-cutoff slice** plus general
investigation tools. The coordinator deliberately did not tell you which
BF you are triaging — that is intentional and you must not try to infer it
in order to fetch additional Jira state.

## Inputs

- **Sliced ticket**: `{{SLICED_PATH}}` — the only ticket-specific source
  you may read. Treat its contents as ground truth at the cutoff timestamp.
- **Skill knowledge** (all paths relative to `${BF_TRIAGE_SKILL_DIR}`,
  which the coordinator has exported into your environment; it
  resolves to `~/.cursor/skills/bf-triage` under Cursor or
  `~/.claude/skills/bf-triage` under Claude Code):
  - `${BF_TRIAGE_SKILL_DIR}/SKILL.md`
  - `${BF_TRIAGE_SKILL_DIR}/reference/severity_types.md`
  - `${BF_TRIAGE_SKILL_DIR}/reference/log_patterns.md`
  - `${BF_TRIAGE_SKILL_DIR}/reference/tools_inventory.md`
  - `${BF_TRIAGE_SKILL_DIR}/reference/workflow_overview.md`
- **Report skeleton**: `${BF_TRIAGE_SKILL_DIR}/templates/report_template.md`

## Output

Write the triage report to `{{OUTPUT_PATH}}`. Do not announce anything in
chat — write the file and stop.

## Forbidden actions (held-in isolation rules — STRICT)

You **MUST NOT** call any of the following while triaging this ticket:

1. `jira_get_issue`, `jira_get_issue_comments`, `jira_search_issues`, or any
   other `jira_*` MCP tool **against the BF under test**. This includes
   passing the BF key as part of a JQL `text ~` clause.
2. `bb_get_bf` against the BF under test (do not look up by BF key). The
   coordinator already filtered the BFG context into the sliced file when
   available.
3. Any `web_search` / `web_fetch` of `https://jira.mongodb.org/browse/<BF>`
   or any URL that resolves to the BF Jira page.
4. `glean.search` queries that include the BF key.
5. Reading any `heldout.md`, `score.v2.md`, or any file path containing the
   substring `heldout` for any ticket. The grader subagent owns those
   files.
6. **Any local git command against the user's working `10gen/mongo` or
   `10gen/dsi` checkout.** When gateway `git_*` fails, the subagent
   MUST switch to a `setup_repos.sh` scratch clone — never to the
   user's working repo, even with a path-filter that "looks safe". The
   full failure protocol and the forbidden-command list are in "Tool
   dispatch & failure handling (per family)" below.

If you do not know which BF you are triaging, that is intentional — the
sliced file usually contains the BF key in URLs, but you must not use it
to fetch fresh Jira state. The cutoff would be invalidated.

If during the run you realise you do not have enough information to make a
confident call, say so in the report's "Root cause hypothesis" with
confidence `low` rather than fetching more Jira context.

## Allowed actions (gateway tools only, plus glean)

You **MAY** use:

- `git_log`, `git_show`, `git_blame`, `git_diff`, `git_search` against
  `10gen/mongo` and `10gen/dsi` (gateway) for bisecting the failure window.
  Read commit messages for PR-number references — do NOT call `user-github`.
- `bb_get_bfg_by_task` if and only if the sliced file contains an Evergreen
  task ID (look for URLs of the form `evergreen.mongodb.com/task/<id>` in
  the description / comments). Do not look up by BF key.
- `evg_get_task_log_summary`, `evg_get_test_results_summary`,
  `evg_get_test_results_detailed`, `evg_get_raw_task_logs` called with task
  IDs / build IDs / version IDs that appear in the sliced file.
- `images_*` tools (gateway) for AMI / package diffs.
- `glean.search` (`user-glean_default`) for runbook / wiki search — but NOT
  with the BF key in the query (see Forbidden actions item 4).
- Reading any file under `${BF_TRIAGE_SKILL_DIR}/`.
- **Local git on a `setup_repos.sh` scratch clone** (NOT the user's
  working `10gen/mongo` / `10gen/dsi` checkout — see Forbidden actions
  item 6 and "Tool dispatch & failure handling" below) when gateway
  `git_*` fails.
- **`evergreen` CLI** invoked via `PATH` (do NOT hardcode the binary
  path) — this is the **default** path for the four CLI-covered EVG
  endpoints (raw task logs, test logs, artifacts, patch list). Write
  every download to `/tmp/bf-triage-workdir-<BF-KEY>/evg/` and clean
  up at the end of the run. See "Tool dispatch & failure handling
  (per family)" below.

You **MUST NOT** use `user-github`, `user-atlassian`, or any other MCP not
listed above.

## Tool dispatch & failure handling (per family)

Defaults:

- **Build Baron (`bb_*`)** → MCP only. No CLI.
- **Evergreen (`evg_*`)** → CLI-first via `evergreen` on `PATH` for
  the four CLI-covered endpoints (`task build TaskLogs`,
  `task build TestLogs`, `fetch --artifacts`, `list-patches`).
  Download whole logs to
  `/tmp/bf-triage-workdir-<BF-KEY>/evg/<task_id>.task.txt`,
  then `tail`/`awk`/`rg` locally. Read it all into context only if
  line count ≤ `BF_TRIAGE_FULL_LOG_MAX_LINES` (default 2000).
  **Mandatory cleanup** at end of run:
  `rm -r /tmp/bf-triage-workdir-<BF-KEY>/`.
- **MCP-only EVG endpoints** (`evg_get_task_log_summary`,
  `evg_get_test_results_summary`, `evg_get_patch_failed_jobs`,
  `evg_get_inferred_project_ids`) → MCP (gateway). No CLI.
- **Git** → gateway `git_*` first; local fallback automatic on
  failure.

Failure handling — you are a background subagent with no live user to
ask:

- `git_*` failure → automatic local fallback (see below).
- CLI-covered `evg_*` is the **default path**, so a gateway
  `evg_*` failure on those endpoints is normally not encountered.
  If you somehow attempted MCP first and got `Not connected`, just
  switch to the CLI; no banner needed.
- `bb_*` / MCP-only `evg_*` / `jira_*` failure → default to
  "Limited evidence" mode (banner required).

- **`git_*` (gateway)** fails with `unauthenticated` / `Not connected` /
  TLS / `MCP error -32000` / timeout: retry ONCE. On second failure,
  fall back to local git, but with a **subagent-only override** of
  `SKILL.md` Step 1's resolution algorithm: the subagent **MUST NOT**
  use the user's working `10gen/mongo` / `10gen/dsi` checkout under
  any circumstances. Concretely:
  1. Ignore any pre-existing `$MONGO_REPO_PATH` / `$DSI_REPO_PATH`
     env vars unless they already point at a path under
     `/tmp/bf-triage-workdir-*`.
  2. Ignore any agent-provided workspace roots (Cursor `<user_info>`
     blocks or equivalent), even if their `origin` URL matches
     `10gen/mongo` or `10gen/dsi`.
  3. Always invoke `bash "${BF_TRIAGE_SKILL_DIR}/scripts/setup_repos.sh"
     [both|mongo|dsi]`. It is idempotent — if a previous run already
     prepared a scratch clone, this is a no-op. Source the `export`
     lines it prints (only the selected repos' env vars are printed)
     and use those values as `$MONGO_REPO_PATH` (and `$DSI_REPO_PATH`,
     if applicable) for the rest of the run. **Pick the selection by
     reading the BF's `Evergreen Project` from `sliced.md`** (or the
     in-context BF JSON in normal modes): pass `both` for `sys-perf*`
     / DSI-flavoured projects, otherwise pass `mongo` — Step 6 only
     consults `10gen/dsi` for sys-perf BFs, so cloning DSI for a
     non-sys-perf BF is wasted network. If the project field is not
     visible to you, default to `both` for safety.

  Rationale: subagents run in the background with no live user to
  interrupt them. Any `git` command that mutates HEAD / the index
  (`git checkout`, `git switch`, `git reset`, `git stash`,
  `git checkout <sha> -- .`, etc.) in the user's working repo is
  unrecoverable damage. Restricting the subagent to a scratch clone
  removes that damage class architecturally — the same commands inside
  a scratch clone are harmless because there is no user WIP, no
  untracked config, and no index state of value.

  Inside the scratch clone you may use the full git surface, including
  `git checkout <sha>` and `git worktree add /tmp/bf-scratch-<sha>
  <sha>`. Read-only commands are still preferred (`git show <sha>:<path>`,
  `git log -p -1 <sha>`, `git blame <sha> -- <path>`) — they are
  faster and produce cleaner cleanup paths. **Clean up** any
  `/tmp/bf-scratch-*` worktrees before writing the report
  (`git worktree remove --force`). Tag each SHA cited via this path
  as `(local fallback)` and add a one-line note in the report. Do
  NOT add the full "Limited evidence" banner just because git fell
  back.

  **Forbidden git commands against the user's working repo** (defence
  in depth — these are the historical violators):
  - `git checkout <ref>` / `git checkout <ref> -- <pathspec>`
  - `git checkout -b <branch>` / `git switch [-c] <branch>`
  - `git reset --hard` / `git reset --mixed` / `git reset --soft`
  - `git stash` / `git stash pop`
  - `git pull` / `git fetch` / `git rebase` / `git merge`
  - `git add` / `git rm` / `git commit` / `git restore --staged`
  - any other command that mutates HEAD, the index, the working tree,
    `refs/`, or `.git/config`.

  These are forbidden even when only the path filter "looks safe"
  (e.g. `git checkout 9d6110b9 -- src/mongo/util/foo.cpp`); the
  pattern is unreliable to gate at agent time, and `setup_repos.sh`
  removes the need entirely.
- **`evg_*` raw task logs / test logs / artifacts / patch list** —
  CLI is the default (see "Tool dispatch & failure handling" above).
  Always invoke as `evergreen ...` (on `PATH`; do not hardcode the
  binary location). Write outputs to
  `/tmp/bf-triage-workdir-<BF-KEY>/evg/`. Do NOT add an `(evg CLI
  fallback)` tag or banner — the CLI is the expected path, not a
  fallback. If `command -v evergreen` does not resolve OR the auth
  probe fails, switch to Step 5.5 Path B (MCP-driven escalation
  ladder) and on second MCP failure, drop to the next bullet.
- **`bb_*` (gateway) OR MCP-only `evg_*` endpoint**
  (`evg_get_patch_failed_jobs`, `evg_get_task_log_summary`,
  `evg_get_test_results_summary`, `evg_get_inferred_project_ids`) OR
  CLI-uncovered `evg_*` failure where the CLI is unavailable: retry
  ONCE. On second failure, write a "limited evidence" note in the
  affected section and add the full "Limited evidence" banner from
  `SKILL.md`'s Limited-evidence fallback section to the report Header.
  Do NOT switch to any other MCP.
- **`jira_*` against the BF under test** is forbidden regardless of
  gateway state — see Forbidden actions above.

## Method

1. Read `{{SLICED_PATH}}` in full. Note the cutoff timestamp at the top.
2. Identify the failing task ID, variant, and time window from the sliced
   description / comments / changelog.
3. Follow the published skill's workflow steps 4–9 (severity classification
   → log evidence → bisect → suspect commits → similar BFs → compose →
   write). Skip steps 1–3 (preflight, fetch, slicing) — those happened
   outside this subagent.
4. **Cite every claim.** Sources include: section name in the sliced file
   (e.g. "sliced.md → Bug Symptoms"), a SHA from `git_show`, a task ID
   from Evergreen, etc. Mark non-sourced inferences with `(AI concluded)`.
5. Be especially careful with the "Root cause hypothesis" and "Recommended
   next steps" sections — the grader scores those most heavily. If the
   sliced ticket points at a known WR pattern (perf-change BF, new-variant
   incompatibility, keep-and-link), apply the matching guidance from
   `reference/log_patterns.md` and `reference/workflow_overview.md`.

## Report structure

Use `${BF_TRIAGE_SKILL_DIR}/templates/report_template.md` verbatim, and
append this footer:

```
## v2 verification appendix

- Cutoff: <copy from sliced.md header>
- Sliced file: {{SLICED_PATH}}
- Forbidden Jira/Glean lookups: NONE attempted
- MCP tools actually called: <comma-separated list, or "none">
```

Stop after writing the file.
