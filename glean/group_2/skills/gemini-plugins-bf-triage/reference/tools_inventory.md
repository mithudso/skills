# MCP Tools Inventory for `bf-triage`

Condensed for the skill. The skill is restricted to **two** MCP servers:

- **`devprod-mcp-gateway`** — Jira, Build Baron, and four structured Evergreen endpoints, plus Git (internal repos). **Required.** Per `SKILL.md` Hard rule 3, the skill does NOT fall back to `user-atlassian`, `user-github`, or any other MCP for these capabilities.
- **`evergreen` CLI on `PATH`** — **default** path for the four CLI-covered Evergreen endpoints (raw task logs, test logs, artifacts, patch list). See "Evergreen CLI — default path for log fetch" below. Uses the user's own `~/.evergreen.yml` credentials.
- Two tool families have **automatic local handling** with no user permission required: `git_*` (local git clones) and CLI-covered `evg_*` endpoints. For `bb_*`, the four MCP-only EVG structured endpoints, and `jira_*` the skill stops and asks the user to fix access (Hard rule 5).
- **`user-glean_default`** — internal wiki / runbook search. Used in Step 7 only.

No other MCPs are permitted. `user-github` is intentionally excluded — PR-level metadata (title, body, reviews, check runs) is outside the gateway's surface today, and the skill uses `git_show` commit messages to identify PRs instead.

Always familiarise yourself with a gateway tool's parameter shape
**before** the first call to any tool you have not used in the current
run. The canonical way is to call the tool once with a minimal payload
and read the error / response — both Cursor and Claude Code surface
MCP tool schemas through their standard tool-call channels.

Cursor additionally exposes a filesystem cache of tool descriptors at
`~/.cursor/projects/<project-id>/mcps/<server>/tools/<tool>.json` (handy
for offline reading, but it lags behind the live gateway). Claude Code
does not expose an analogous cache; use the call-and-read pattern, or
fall back to this file's hand-maintained catalogue below.

## Jira tools (`devprod-mcp-gateway`)

| Tool | Use | Key params |
| ---- | --- | ---------- |
| `jira_get_issue` | Full BF details, all custom fields | `issue_key` |
| `jira_get_issue_comments` | Comment timeline (body, author, timestamp) | `issue_key`, `max_results` |
| `jira_search_issues` | JQL search, max 50 per call. Used in Mode B (team query) and Step 7 (similar-BF search). | `jql`, `max_results` |
| `jira_get_issue_transitions` | Discover valid status transitions | `issue_key` |
| `jira_add_comment` | **Write — never call without explicit user approval** | `issue_key`, `body` |

### Mode B team-query JQL pattern

For batch-triaging a team's active BFs (`SKILL.md` Step 0.5), use:

```text
project = BF
  AND "Assigned Teams" = "<team>"
  AND status in ("Needs Triage", "Open")
  [AND <extra-jql>]
ORDER BY created DESC
```

Defaults (overridable per `SKILL.md` "Configuration defaults (Mode B)"):

- `<team>` — required, taken from the user prompt or wrapper.
- statuses — `"Needs Triage", "Open"` (env `BF_TRIAGE_TEAM_STATUSES`,
  comma-separated).
- `<extra-jql>` — empty by default (env `BF_TRIAGE_TEAM_JQL_EXTRA`).
- `max_results` — `5` (env `BF_TRIAGE_TEAM_LIMIT`). Tool hard cap is
  50.

Examples (full JQL strings the skill might send to `jira_search_issues`):

```text
project = BF AND "Assigned Teams" = "Workload Resilience" AND status in ("Needs Triage", "Open") ORDER BY created DESC

project = BF AND "Assigned Teams" = "DevProd Performance Infrastructure" AND status in ("Needs Triage", "Open") AND Temperature ~ "hot" ORDER BY created DESC

project = BF AND "Assigned Teams" = "Query Execution" AND status in ("Needs Triage", "Open") AND created <= -48h ORDER BY created DESC
```

The wrapper helper `${BF_TRIAGE_SKILL_DIR}/scripts/list_active_bfs.sh`
(under `~/.cursor/skills/bf-triage/` on Cursor or `~/.claude/skills/bf-triage/`
on Claude Code) prints the resolved JQL + a ready-to-paste skill
invocation prompt.
The shell cannot itself call MCP tools — it only builds the prompt.

### Custom fields the skill cares about

- `Assigned Teams` — current team(s) responsible
- `Evergreen Project` — project identifier
- `Temperature` — hot / warm / cold
- `Performance Change Type` — Regression / Improvement / etc.
- `Bug Symptoms` — short failure signature(s)
- `Severity Type` — one of the entries in `severity_types.md`

## Build Baron tools (`devprod-mcp-gateway`)

No CLI exists; MCP is the only interface.

| Tool | Use | Key params |
| ---- | --- | ---------- |
| `bb_get_bf` | BF: variants, tasks, tests, BFGs, severity, time range | `bf_key` (e.g. `BF-43272`) |
| `bb_get_bfg` | Single BFG: extracted faults, log snippets, severity, BF suggestions, hosts, event timeline | `bfg_key` |
| `bb_get_bfg_by_task` | BFG for a specific Evergreen task | `task_id`, `execution` (default 0) |
| `bb_search_bfgs` | BFG search with rich filters | `bf_key`, `start_date` (MM-DD-YYYY), `end_date`, `projects`, `tasks`, `tests`, `variants`, `terms`, `status`, `max_results` |

## Evergreen tools (`devprod-mcp-gateway`)

| Tool | Use | Key params |
| ---- | --- | ---------- |
| `evg_get_patch_failed_jobs` | Up to 20 failed tasks for a patch | `patch_id` |
| `evg_get_raw_task_logs` | Full raw task log (default tail 500 lines, 0 = full) | `task_id`, `tail_lines`, `execution` |
| `evg_get_task_log_summary` | Structured log entries via GraphQL | `task_id`, `execution` |
| `evg_get_test_results_summary` | Up to 30 test results (status, duration, log URLs) | `task_id`, `test_name`, `execution` |
| `evg_get_test_results_detailed` | Raw log for a specific test (tail 500 default) | `task_id`, `test_name` (Job0/Job1 etc.), `tail_limit`, `execution` |
| `evg_download_task_artifacts` | List artifacts (names, URLs, content types) | `task_id`, `execution` |

`evg_get_test_results_detailed` `test_name` is `Job0` / `Job1` not the test
path — pick the runner-output job number, not the `.js` filename.

### Evergreen CLI — default path for log fetch

For the four CLI-covered Evergreen endpoints, the `evergreen` CLI is
the **default**, not a fallback. The CLI:

- Has no token / context cost (writes output to disk).
- Uses the user's own `~/.evergreen.yml` credentials, which usually
  has **better** project-scoped access than the gateway SA (sidesteps
  the DEVPROD-30722 authorization gap that blocks the gateway on
  `sys-perf` and similar projects).
- Empirically downloads a full sys-perf task log (~3000 lines, ~1 MB)
  in ~2 seconds — small enough that downloading whole and grepping
  locally is faster and cheaper than tail-then-escalate.

| Gateway MCP tool | CLI equivalent | Coverage | Default path |
| ---------------- | -------------- | -------- | ------------ |
| `evg_get_raw_task_logs` | `evergreen task build TaskLogs --task_id <id> --tail_limit 0 --out <file>` | Full | CLI |
| `evg_get_test_results_detailed` | `evergreen task build TestLogs --task_id <id> --log_path <test> --out <file>` | Full | CLI |
| `evg_download_task_artifacts` | `evergreen fetch --artifacts --task <id> --dir <dir>` | Full | CLI |
| `evg_list_user_recent_patches` | `evergreen list-patches -n <limit> --json` | Full | CLI |
| `evg_get_task_log_summary` | — | None (GraphQL summary is gateway-only) | MCP |
| `evg_get_test_results_summary` | — | None | MCP |
| `evg_get_patch_failed_jobs` | — | None | MCP |
| `evg_get_inferred_project_ids` | — | None | MCP |

Operationally (see `SKILL.md` Step 5.5 Path A):

1. Download the full log once to
   `/tmp/bf-triage-workdir-<BF-KEY>/evg/<task_id>.task.txt`.
2. `wc -l` it. If ≤ `BF_TRIAGE_FULL_LOG_MAX_LINES` (default 2000),
   read it all; otherwise keep on disk and use local `tail`, `awk`
   time-window slicing, and `rg` keyword search.
3. **Clean up** `/tmp/bf-triage-workdir-<BF-KEY>/` at the end of
   triage (workflow Step 13).

Binary location: assume the `evergreen` binary is on `PATH`; do NOT
hardcode an install location — it may live at
`/usr/local/bin/evergreen`, `~/cli_bin/evergreen`, a `cipd` path, etc.

Tagging: under Path A, do NOT tag values with `(evg CLI fallback)` —
the CLI is the expected path, not a fallback. Reserve the
`(evg CLI fallback)` tag and the one-line report note for the
specific Path-B-rescued-by-CLI case where MCP was attempted first
(e.g. some other agent's compatibility flow) and then the CLI
rescued the call.

If `command -v evergreen` does not resolve OR the auth probe fails,
the skill switches to Step 5.5 Path B (MCP-driven 500-tail
escalation ladder). MCP-only structured endpoints fall to "Limited
evidence" only when both their MCP path AND any alternative are
exhausted.

## Gateway Git tools (`devprod-mcp-gateway`)

Operate on pre-registered internal repos with no local clone needed:

| Tool | Use |
| ---- | --- |
| `git_log` | Commit history (`since`/`until`/`author`/`grep`/`path`) |
| `git_diff` | Diff between two refs |
| `git_blame` | Line-level blame |
| `git_show` | One commit (message + diff) |
| `git_search` | `git grep` regex content search at a ref |

Pre-registered repos include `10gen/mongo`, `10gen/dsi`, `10gen/mms`,
`10gen/mms-automation`, `10gen/buildhost-post-config`,
`10gen/buildhost-configuration`, `10gen/toolchain-builder`,
`10gen/distro-settings`, `10gen/signal-processing-service`,
`10gen/performance-tooling-docs`, `evergreen-ci/evergreen`.

### Local fallback for old commits (automatic when gateway `git_*` fails)

Documented in `SKILL.md` "Step 6 fallback — local history". Triggered
**automatically** after the Hard rule 5 retry has failed for any
`git_*` call. **No user permission required** — git data is always
recoverable locally, unlike `bb_*` / `evg_*`. Pick the read-only
command for the inspection style you need:

| Goal | Command (run in user repo first, scratch clone as fallback) |
| ---- | --- |
| Verify SHA is local | `git -C "$REPO" cat-file -e <sha>^{commit}` |
| Fetch a missing SHA | `git -C "$REPO" fetch --depth=1 origin <sha>` |
| Show one file at SHA | `git -C "$REPO" show <sha>:<path>` |
| Full diff + message | `git -C "$REPO" log -p -1 <sha>` |
| Blame at SHA | `git -C "$REPO" blame <sha> -- <path>` |
| Multi-file `Grep`/`Read` at SHA | `git -C "$REPO" worktree add /tmp/bf-scratch-<sha> <sha>` → reads → `git -C "$REPO" worktree remove --force /tmp/bf-scratch-<sha>` |

Prefer the user's working repos resolved in `SKILL.md` Step 1 as
`$MONGO_REPO_PATH` / `$DSI_REPO_PATH` (they already exist and are
usually current). Fall back to a `setup_repos.sh` scratch clone only
if Step-1 resolution returns nothing.

Hard constraint: when `$REPO` is the user's working repo (anything not
under the `/tmp/bf-triage-workdir-*` scratch-clone prefix), never run
`git checkout`, `git switch`, `git checkout -b`, `git reset`,
`git stash`, or any other command that mutates the index or HEAD — they
would clobber the user's WIP. `git checkout` is only allowed inside a
`setup_repos.sh` scratch clone, where the branch must be deleted
before returning control.

Cleanup invariant: before the report is written, run `git -C "$REPO"
worktree list` and remove any `/tmp/bf-scratch-*` leftovers with
`git worktree remove --force`. Leftover worktrees keep refs alive and
confuse the user's next `git status`.

Tag any commit citation that came from this branch as `(local fallback)`
in the report. Do NOT add the full "Limited evidence" banner just for
a git fallback — that banner is reserved for `bb_*` / `evg_*` outages.
Use a one-line note instead.

## GitHub access

The skill does NOT use `user-github`. PR-level data is obtained from
gateway `git_show` commit messages (which include the PR number / link
trailer) and from `git_log` / `git_blame` / `git_diff` on `10gen/mongo` and
`10gen/dsi`.

**Known gap vs. `user-github`** (intentional, no fallback):

- PR title / body / reviewers / status checks
- `search_pull_requests` by commit SHA or label
- `get_review_comments` / `get_check_runs`

If PR-level metadata becomes essential for a future triage class, the right
fix is to extend the gateway rather than to re-enable `user-github` here.

## Glean MCP (`user-glean_default`)

| Tool | Use |
| ---- | --- |
| `search` | Internal wiki / runbook / design-doc search |
| `read_document` | Full content of a doc found via search |
| `chat` | Synthesis question across many sources |

Useful for finding past triage decisions, team runbooks, known-failure pages.

## Tool-sequence template (gateway only, except step 5 glean)

```text
0. PRECHECK          (gateway — fail-fast, parallel)
   ├─ jira_list_projects                       (cheapest JIRA ping)
   ├─ bb_get_bf <BF-KEY>                       (cached for step 1)
   ├─ git_log repo=10gen/mongo limit=1         (cheapest Git ping)
   └─ evg_list_user_recent_patches limit=1     (user-scoped, no project perm)
   ↳ on any failure → Hard rule 5(b): stop + ask user to toggle MCP
     in Cursor Settings UI → Tools & MCP → devprod-mcp-gateway

1. UNDERSTAND        (gateway)
   ├─ jira_get_issue           (BF custom fields, current Assigned Teams)
   ├─ jira_get_issue_comments
   └─ bb_get_bf                 (variants, tasks, tests, time range, BFGs)
                                — skip if already cached from step 0

2. CLASSIFY          (in-skill reference + gateway)
   ├─ severity_types.md          → match faults
   └─ bb_search_bfgs             → Hot/Cold counts in last 30 days

3. ANALYSE FAILURE INSTANCE  (gateway)
   ├─ bb_get_bfg or bb_get_bfg_by_task   (extracted faults)
   ├─ evg_get_task_log_summary            (structured log entries)
   ├─ evg_get_raw_task_logs               (deep dive when needed)
   └─ evg_get_test_results_summary        (per-test status / duration)

4. BISECT            (gateway)
   ├─ git_log     (commits in failure window on 10gen/mongo or 10gen/dsi)
   ├─ git_show    (inspect suspicious commits — commit message has PR link)
   ├─ git_blame   (who last touched the failing code)
   └─ git_search  (regex search at a ref for related symbols)

5. SCOPE             (gateway + glean)
   ├─ bb_search_bfgs            (other instances, branch/variant scope)
   ├─ jira_search_issues        (similar BFs by symptom text)
   └─ glean.search              (past triage notes / runbook)

6. (only on explicit user approval) ACT  (gateway)
   └─ jira_add_comment with the report's executive summary
```

## Gateway failure handling

Per `SKILL.md` Hard rule 5. Failures are classified **per tool family**:

- The skill must run the Step 0 precheck first; the precheck calls one
  representative tool from each family in parallel.
- Retry once on `unauthenticated` / `forbidden` / `Not connected` /
  `MCP error -32000` / TLS verification error / timeout / connection
  refused.
- On second failure, branch by tool family:
  - **`git_*`** → silently fall back to local git (see "Local fallback
    for old commits" above). No user permission required. Tag SHAs as
    `(local fallback)`.
  - **`evg_*` CLI-covered endpoints** → usually a non-event because
    Step 5.5 Path A is CLI-first. Path B (the MCP escalation ladder)
    is only entered when the CLI itself is unavailable, in which case
    a gateway `evg_*` failure means **stop** and ask the user to
    reconnect (with optional "Limited evidence" mode on explicit
    approval).
  - **`evg_*` MCP-only endpoints** (`evg_get_patch_failed_jobs`,
    `evg_get_task_log_summary`, `evg_get_test_results_summary`,
    `evg_get_inferred_project_ids`) → **stop** and ask the user to
    reconnect. No CLI equivalent. Only with explicit user approval,
    switch to "Limited evidence" mode and add the corresponding
    banner.
  - **`jira_*`** → **stop** and ask the user to reconnect the MCP from
    **Cursor Settings UI → Tools & MCP → `devprod-mcp-gateway`** (toggle
    off, wait ~2 s, toggle on). Never fall back to `user-atlassian`.
  - **`bb_*`** → **stop** and ask the user to reconnect (same UI hint).
    No CLI exists for Build Baron. Only with explicit user approval,
    switch to "Limited evidence" mode and add the corresponding
    banner.
- The `devprod-mcp-proxy` stdio process can crash (observed: Go
  `singleflight` panic → transport_closed → FSM stuck in `closed`; also
  observed: TLS verify error `x509: certificate signed by unknown
  authority` on the proxy's HTTPS client). In both cases the UI toggle
  is the documented recovery path. There is no `cursor-agent` / `cursor`
  CLI subcommand for MCP restart on this remote host, no auto-retry in
  Cursor's V2 FSM, and no MCP-disconnect hook.
- `evg_*` project-scoped tools (`evg_get_task_log_summary`,
  `evg_get_raw_task_logs`, `evg_get_test_results_summary`,
  `evg_get_test_results_detailed`, `evg_get_patch_failed_jobs`,
  `evg_download_task_artifacts`) returning *"does not have permission to
  'view tasks' for the project '<X>'"* is a per-project authorization
  gap on the gateway's SPIFFE service account (DEVPROD-30722), NOT a
  connectivity failure. Confirmed blocked: `sys-perf`. The `evergreen`
  CLI uses the user's own `~/.evergreen.yml` and usually has access —
  try the CLI for the affected endpoint first (CLI-covered endpoints
  only). If the CLI also lacks access or there is no CLI equivalent,
  mark EVG-derived evidence on that project as "Limited evidence" in
  the report.
- `evg_get_raw_task_logs` returning empty / 404 is a data outcome (raw
  logs expired), not a connectivity failure — proceed and note in the
  Log evidence section.
- `bb_search_bfgs` returning 0 in the 30-day window: widen to 90 days
  and note the window used in the report.
