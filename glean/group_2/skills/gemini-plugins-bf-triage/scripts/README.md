# bf-triage skill — scripts/

Three groups of scripts:

- **Runtime helpers** — used by the published skill at user-invocation time
  (`setup_repos.sh` for auto-clones; `list_active_bfs.sh` for Mode B
  team-query prompt building; `_make_parsley_url.py` for Log-evidence
  deep-links; `md_to_pdf.py` for Step 9b PDF rendering;
  `attach_pdf_to_jira.py` for the Step 11b opt-in PDF attachment).
- **`automation/`** — optional, macOS-oriented unattended-scheduling
  wrappers (NOT part of the published skill payload): drive a headless
  `claude -p` Mode B run on a launchd schedule, with proxy-token
  preflight and notifications. See each file's header for setup.
- **v2 verification harness** — used by a coordinator agent when
  testing or verifying the skill (`run_held_in_test.sh`,
  `run_batch.sh`, `_slice_helper.py`). Implements the "Verification
  Methodology v2 — Process Isolation" pattern: the coordinator
  pre-fetches each BF ticket, the slicer splits it into held-in /
  held-out files at a cutoff timestamp, then the coordinator spawns
  isolated triager and grader subagents in parallel. The scripts are
  self-contained — see the leading comments in each file for usage.
  Verification reports (if you generate any) are intended to live
  alongside the BF reports under `${BF_TRIAGE_OUTPUT_DIR:-./bf-reports}/`
  in whichever workspace you run them from. (Mode C dev artifacts —
  sliced / heldout / score files — go to `/tmp/bf-triage-test/` and
  are not controlled by `BF_TRIAGE_OUTPUT_DIR`; see
  `reference/verification_mode.md`.)

## File map

| File | Role | Run by |
|---|---|---|
| [`setup_repos.sh`](setup_repos.sh) | Clones `10gen/mongo` + `10gen/dsi` (history-preserving, blob-less) for the published skill's bisect step. Accepts `both` / `mongo` / `dsi` (default `both`). | The published skill at runtime, **only** when local clones are missing |
| [`list_active_bfs.sh`](list_active_bfs.sh) | Mode B helper: resolves the team-name / limit / statuses / extra-jql config and prints the JQL + a ready-to-paste agent prompt. **Does NOT call MCP** — shells can't reach the gateway directly. | User, before invoking the skill in Mode B |
| [`md_to_pdf.py`](md_to_pdf.py) | Converts a bf-triage Markdown report to a styled PDF (letter, 0.75" margins, blue table headers, amber blockquotes, h1 page-breaks). Adapted from the `ftdc-analysis` skill. Requires `pip3 install markdown weasyprint`. Exits 2 if the deps are missing so the caller can degrade gracefully. | The published skill at runtime, Step 9b, **only** when `BF_TRIAGE_GENERATE_PDF=1` or the user asked for PDF |
| [`run_held_in_test.sh`](run_held_in_test.sh) | Slices ONE pre-fetched BF JSON into `sliced.md` / `heldout.md` / `cutoff.txt`. Accepts either the gateway two-call shape (`--issue-json` + `--comments-json`) or a legacy single-blob shape (`--input-json`). | Coordinator (or `run_batch.sh`) |
| [`run_batch.sh`](run_batch.sh) | Fans out `run_held_in_test.sh` for N BFs in parallel (background + `wait`) | Coordinator |
| [`_stitch_gateway_jira.py`](_stitch_gateway_jira.py) | Normalizes the gateway two-call shape (`jira_get_issue` + `jira_get_issue_comments`) into the legacy single-blob shape the slicer understands. Strips nonce wrappers. Stop-gap until the gateway exposes a `jira_get_issue_changelog` tool. | Not directly — invoked by the slicer |
| [`_slice_helper.py`](_slice_helper.py) | Internal Python rendering helper called by `run_held_in_test.sh` | Not directly — invoked by the slicer |
| [`_make_parsley_url.py`](_make_parsley_url.py) | Composes a filtered + bookmarked Parsley deep-link URL for the report's Log-evidence section (offline; takes task id, execution, line numbers, filter patterns). | The published skill at runtime, Step 5 |
| [`attach_pdf_to_jira.py`](attach_pdf_to_jira.py) | Uploads a rendered PDF to a BF via the Jira REST API (stdlib only; resolves a `JIRA_PERSONAL_TOKEN` from env or the Cursor/Claude MCP config; never echoes it). Exits 2 if no token so the caller degrades gracefully. | The published skill at runtime, Step 11b, **only** on explicit PDF-attach opt-in |
| [`automation/run-bf-triage.sh`](automation/run-bf-triage.sh) | macOS unattended wrapper: proxy-token preflight, then a headless `claude -p` Mode B batch-triage under a hard timeout, with success/failure notifications. NOT part of the published skill. | launchd / cron / manual (macOS) |
| [`automation/refresh-mcp-auth.sh`](automation/refresh-mcp-auth.sh) | Re-establishes the `devprod-mcp-gateway` (kanopy-oidc) session via device-flow (no local browser) and seeds the proxy token file. Run when the session lapses (e.g. each morning / after a weekend). | User, before a scheduled run |
| [`automation/install-launchd.sh`](automation/install-launchd.sh) | Generates a machine-specific LaunchAgent from the plist template and loads it into the per-user GUI domain. Idempotent; derives all paths from the clone + `$HOME`. | User, once, to schedule the job (macOS) |
| [`automation/com.example.bf-triage.plist.template`](automation/com.example.bf-triage.plist.template) | LaunchAgent template consumed by `install-launchd.sh` (placeholders filled at install time). Not run directly. | Not directly — consumed by `install-launchd.sh` |

## Mode B helper — quick usage

```bash
# Pick the install path matching your agent runtime, OR export
# BF_TRIAGE_SKILL_DIR to override:
SKILL_DIR=~/.cursor/skills/bf-triage   # Cursor
# SKILL_DIR=~/.claude/skills/bf-triage  # Claude Code

# Default (limit 5, statuses "Needs Triage,Open"):
"$SKILL_DIR/scripts/list_active_bfs.sh" --team "Workload Resilience"

# Override limit + add JQL clause:
"$SKILL_DIR/scripts/list_active_bfs.sh" --team "Query Execution" \
    --limit 3 --extra-jql 'Temperature ~ "hot"'

# Override via env var:
BF_TRIAGE_TEAM_LIMIT=10 \
"$SKILL_DIR/scripts/list_active_bfs.sh" --team "DevProd Performance Infrastructure"
```

The script prints the resolved JQL plus a one-line invocation prompt
to paste into your agent chat (Cursor / Claude Code). The skill itself
does the `jira_search_issues` call
(via the gateway) and then fans out per-BF triagers in parallel — see
`SKILL.md` Step 0.5 for the dispatch logic.

## Division of labour for v2 verification

The verification flow has **four distinct actors** that must not share
context. They are split as follows:

```text
+-----------+   +---------------+   +--------------+   +-----------+
| Coord-    |   | run_batch.sh  |   | Triager      |   | Grader    |
| inator    |   | (this dir)    |   | subagent(s)  |   | sub-      |
| (you)     |   |               |   |              |   | agent(s)  |
+-----------+   +---------------+   +--------------+   +-----------+
      |                |                  |                  |
      | 1. fetch raw   |                  |                  |
      |   JSON via MCP |                  |                  |
      | (gateway       |                  |                  |
      |  jira_get_issue|                  |                  |
      |  expand=       |                  |                  |
      |  changelog)    |                  |                  |
      |                |                  |                  |
      | 2. save to     |                  |                  |
      |   /tmp/bf-     |                  |                  |
      |   triage-test/ |                  |                  |
      |   <BF>/raw.json|                  |                  |
      |                |                  |                  |
      | 3. invoke ---->| 4. fan out       |                  |
      |   for ALL BFs  |    in parallel,  |                  |
      |   in one shell |    one slicer    |                  |
      |   call         |    per BF        |                  |
      |                |                  |                  |
      |                | 5. write         |                  |
      |                |    sliced.md +   |                  |
      |                |    heldout.md +  |                  |
      |                |    cutoff.txt    |                  |
      |                |                  |                  |
      | 6. spawn N     |                  | 7. read sliced.md|
      |   triager sub- |                  |    only; write   |
      |   agents in    |                  |    triage.v2.md  |
      |   parallel via |                  |                  |
      |   Task tool    |                  |                  |
      |   (run_in_back-|                  |                  |
      |    ground:true)|                  |                  |
      |                |                  |                  |
      | 8. spawn N     |                  |                  | 9. read
      |   grader sub-  |                  |                  |    triage.v2.md
      |   agents in    |                  |                  |    + heldout.md
      |   parallel     |                  |                  |    + cutoff.txt;
      |                |                  |                  |    write
      |                |                  |                  |    score.v2.md
      |                |                  |                  |
      | 10. aggregate  |                  |                  |
      |   scores into  |                  |                  |
      |   verifica-    |                  |                  |
      |   tion.md      |                  |                  |
```

## Key invariants enforced by the harness

1. **Slicing happens outside the triager.** The triager subagent never
   sees `raw.json` or any post-cutoff facts — only `sliced.md`. The
   slicer drops every comment / changelog entry whose timestamp is
   `>= cutoff` (including those at the exact cutoff timestamp).
2. **Cutoff resolution order**, computed by `jq` in `run_held_in_test.sh`
   (no agent judgment):
   1. First changelog event setting `Assigned Teams = OWNING_TEAM`
      (legacy-shape input only).
   2. `custom_fields["Team Assigned (Effective Date)"]` value
      (gateway-shape input, used because the gateway exposes no changelog
      tool today). Requires the current `custom_fields["Assigned Teams"]`
      to include OWNING_TEAM.

   If neither yields a timestamp, the slicer exits with code 1 and the BF
   must be excluded from the corpus. Gateway-shape inputs also trigger a
   prominent **APPROXIMATION WARNING** banner in `sliced.md` because
   field snapshots fall back to current `custom_fields` (post-cutoff)
   instead of value-at-cutoff.
3. **The triager is parameterised by file path, NOT by BF key.** The
   coordinator must pass `{{SLICED_PATH}}` and `{{OUTPUT_PATH}}` to the
   subagent prompt template (see
   [`../templates/triager_prompt.md`](../templates/triager_prompt.md)) and
   must NOT mention the BF key. The triager prompt itself forbids fetching
   the BF directly via `jira_*` / `bb_*` / web tools.
4. **Coordinator never reads `heldout.md`.** Only the grader subagent does.

## Coordinator pre-flight checklist (P10)

Before invoking `run_batch.sh`:

- [ ] For each BF in the corpus, make **two** gateway calls:
      - `devprod-mcp-gateway.jira_get_issue { "issue_key": "<BF>" }`
        → write to `${OUTPUT_DIR}/<BF>/raw.json`
      - `devprod-mcp-gateway.jira_get_issue_comments
         { "issue_key": "<BF>", "limit": 200 }`
        → write to `${OUTPUT_DIR}/<BF>/comments.json`
- [ ] Do NOT fall back to `user-atlassian` — if the gateway fails, retry
      once and then stop per SKILL.md Hard rule 5.
- [ ] Confirm both files are valid JSON: `jq empty <path>`.
- [ ] `run_held_in_test.sh` auto-detects the gateway shape (no top-level
      `changelogs`, has `custom_fields`) and stitches via
      `_stitch_gateway_jira.py` into `${OUTPUT_DIR}/<BF>/normalized.json`,
      using `Team Assigned (Effective Date)` as the cutoff fallback. No
      coordinator-side stitching needed.

After `run_batch.sh` returns:

- [ ] Confirm cutoff for each BF matches the value recorded in any
      prior verification aggregate under
      `${BF_TRIAGE_OUTPUT_DIR:-./bf-reports}/verification-*.md`
      (round-1 / round-2 cutoffs are listed there); investigate any drift.
- [ ] For each BF, spawn a triager subagent in the SAME tool-call batch
      using the `Task` tool (Cursor: with `subagent_type: "generalPurpose"`
      and `run_in_background: true`; Claude Code: no extra params, just
      multiple `Task` calls in one assistant message). Pass
      `{{SLICED_PATH}}=${OUTPUT_DIR}/<BF>/sliced.md` and
      `{{OUTPUT_PATH}}=${OUTPUT_DIR}/<BF>/triage.v2.md`. Do NOT include
      the BF key in the prompt.
- [ ] After all triagers complete, spawn N grader subagents in the same
      tool-call batch (same Cursor/Claude Code parameter difference),
      passing `{{REPORT_PATH}}`, `{{HELDOUT_PATH}}`, `{{SLICED_PATH}}`,
      `{{CUTOFF_PATH}}`, `{{SCORE_PATH}}`, and `{{BF_KEY}}` for the
      report header.

## What `run_batch.sh` does NOT do

- It does not fetch the BF from Jira (no MCP from shell).
- It does not launch agents (only the coordinator can do that, via the
  `Task` tool).
- It does not aggregate scorecards. The coordinator (or a single dedicated
  synthesis subagent) writes the final v2 section in `verification.md`.
