# Verification / Replay Mode (Mode C)

This document is the entry point for the **verification path** of the
`bf-triage` skill. It is NOT used in normal triage (Modes A and B). Read
this file only when the user is replaying / grading an old BF to test
the skill itself.

The main `SKILL.md` covers normal triage; everything in this file
applies on top of that. When the verification path delegates to the
normal triage rules (severity classification, log evidence, suspect
commits, etc.), it does so by reference — see `SKILL.md` Steps 4–9.

## When to invoke

Trigger phrases (case-insensitive substring match in the user prompt):

- "verify ... triage ..."
- "replay ... BF-NNNNN"
- "test the skill against BF-NNNNN"
- "grade ... triage report"
- "v2 verification"
- "held-in test"
- "score the skill on BF-NNNNN"

Anything else falls back to normal Mode A / Mode B. If the prompt is
ambiguous, ask the user once: *"Treat as verification/replay run, or
regular triage?"*

## Purpose

Verification mode runs the triage skill against a **closed historical
BF** while hiding the post-cutoff facts (the team's eventual
comments, resolution, the actual fix commit), so that we can measure
how well the skill would have done if it had been the triager. The
post-cutoff facts then become the grader's ground truth.

Two things make this hard to do in-process and force a multi-subagent
design:

1. The same context window cannot simultaneously hold the pre-cutoff
   slice **and** the held-out facts without leakage. We need process
   isolation between the triager and the grader.
2. The triager must not call `jira_*` against the BF under test — that
   would reveal the entire ticket including the held-out comments. We
   enforce this in the triager subagent's prompt.

## Actors and process isolation

There are **four distinct actors**:

```text
Coordinator (the parent agent the user invokes)
   │
   │ 1. Fetches raw Jira state via the gateway
   │ 2. Invokes run_held_in_test.sh for each BF
   │
   ├──► run_held_in_test.sh
   │       (slices raw JSON → sliced.md, heldout.md, cutoff.txt)
   │
   ├──► Triager subagent(s)   ← reads ONLY sliced.md
   │       (one per BF, run in parallel via the subagent-dispatch
   │        tool — Cursor: `Task`, Claude Code: `Agent` — see
   │        "C3 — Spawn triager subagents" for the runtime-
   │        specific parameters)
   │
   └──► Grader subagent(s)    ← reads triage report + heldout.md
           (one per BF, run after the matching triager finishes)
```

Strict isolation invariants:

| Actor | May read | Must NOT read |
| ----- | -------- | ------------- |
| Coordinator | All inputs and outputs | — |
| Slicer (script) | Pre-fetched Jira JSON | (it's just a script — no agent) |
| Triager subagent | `sliced.md` ONLY (for ticket data) | `heldout.md`, `cutoff.txt`, any sibling score file, the BF Jira page via `jira_*` or `web_*` |
| Grader subagent | Triage report + `heldout.md` + `cutoff.txt` + `sliced.md` for cross-reference | — (but never EDIT the triage report) |

The triager prompt template at
[`../templates/triager_prompt.md`](../templates/triager_prompt.md)
enforces the triager-side rules; the grader prompt at
[`../templates/grader_prompt.md`](../templates/grader_prompt.md)
enforces the grader-side rules.

## Coordinator workflow

```text
Verification Progress:
- [ ] C0. Confirm verification mode (trigger-phrase match or user-confirmed)
- [ ] C0.5. Pick run directory (NEVER reuse /tmp/bf-triage-test/<BF>/
       if it already exists — see below)
- [ ] C1. For each BF: pre-fetch raw Jira state via gateway     ── per BF
- [ ] C2. For each BF: invoke run_held_in_test.sh               ── per BF
- [ ] C3. Spawn N triager subagents in parallel (Cursor: `Task` tool;
       Claude Code: `Agent` tool)
- [ ] C4. After all triagers complete, spawn N grader subagents in parallel
- [ ] C5. Aggregate per-BF scorecards into a single verification report
- [ ] C6. Cleanup the THIS-RUN folder only (never the suffix-less folder
       belonging to another agent)
```

### C0.5 — Pick run directory (no-overwrite rule)

For each BF, the coordinator MUST pick a fresh work directory before
calling the gateway or the slicer. A pre-existing
`/tmp/bf-triage-test/<BF>/` means another agent (current or past) is
either still running or left artifacts behind for inspection — its
contents are **immutable** to this run.

Resolution:

```bash
BASE=/tmp/bf-triage-test
mkdir -p "$BASE"
if [[ ! -e "$BASE/<BF>" ]]; then
  RUN_DIR="$BASE/<BF>"            # first run for this BF
else
  RUN_DIR=$(mktemp -d "$BASE/<BF>-XXXXXX")   # race-safe new run folder
fi
```

`mktemp -d` is atomic and race-safe under concurrent agents (two
agents picking the same BF at the same time each get a distinct
random suffix). Use `$RUN_DIR` as the base path for every file the
workflow writes (`raw.json`, `comments.json`, `sliced.md`,
`heldout.md`, `cutoff.txt`, `triage.md`, `score.v2.md`). Pass it to
the slicer via `--output-dir $(dirname $RUN_DIR)` and
`--issue-json $RUN_DIR/raw.json` etc., OR by overriding the slicer's
`OUTPUT_DIR` env var. Substitute `$RUN_DIR` for `/tmp/bf-triage-test/<BF>/`
in every later step's path templates (C1, C2, C3 substitutions, C4
substitutions, C5 aggregation, C6 cleanup).

### C1 — Pre-fetch raw Jira state (gateway, two calls per BF)

The `devprod-mcp-gateway` does NOT expose a Jira changelog tool today
(verified by introspecting the gateway tool list — no `*changelog*`
tool is present).
We therefore make TWO calls per BF and let the harness stitch them:

```python
CallMcpTool devprod-mcp-gateway jira_get_issue {
  "issue_key": "<BF>"
}
# → save to $RUN_DIR/raw.json

CallMcpTool devprod-mcp-gateway jira_get_issue_comments {
  "issue_key": "<BF>",
  "limit": 200
}
# → save to $RUN_DIR/comments.json
```

`$RUN_DIR` comes from C0.5 — never hardcode `/tmp/bf-triage-test/<BF>/`.

Validate each file with `jq empty <path>` before proceeding.

If either gateway call fails, follow `SKILL.md` Hard rule 5 — do NOT
fall back to `user-atlassian`. The verification harness is read-only on
Jira so a single retry + UI-toggle is the entire recovery surface.

### C2 — Invoke the slicer

The slicer is at
[`../scripts/run_held_in_test.sh`](../scripts/run_held_in_test.sh).
For one BF:

```bash
bash "${BF_TRIAGE_SKILL_DIR}/scripts/run_held_in_test.sh" BF-NNNNN
```

It auto-detects the gateway shape (no top-level `changelogs`, has
`custom_fields`), stitches `raw.json` + `comments.json` into
`normalized.json` via
[`_stitch_gateway_jira.py`](../scripts/_stitch_gateway_jira.py),
strips nonce wrappers, and writes three artifacts:

- `sliced.md` — held-in (pre-cutoff) view for the triager.
- `heldout.md` — held-out (post-cutoff) view for the grader.
- `cutoff.txt` — audit log: cutoff timestamp, source, event.

For N BFs in parallel, use
[`../scripts/run_batch.sh`](../scripts/run_batch.sh).

**Cutoff resolution order**, baked into the script (no agent
judgment):

1. First changelog event setting `Assigned Teams = OWNING_TEAM`
   (legacy-shape JSON only — the gateway doesn't expose changelogs).
2. `custom_fields["Team Assigned (Effective Date)"]` value, gated on
   the current `custom_fields["Assigned Teams"]` containing
   OWNING_TEAM. This is the **default path** for gateway-fetched BFs.

If neither yields a timestamp, the slicer exits with code 1 and the BF
must be excluded from the corpus. `OWNING_TEAM` defaults to "Workload
Resilience"; override with `--owning-team` for BFs historically owned
by a former / sibling name.

**Approximation caveat** (gateway-shape inputs only): because we lack
a changelog, the held-in field snapshot is read from the BF's CURRENT
`custom_fields` and is therefore post-cutoff. For BFs whose fields
changed between routing and close (Severity Type reclassified
mid-investigation, etc.) the snapshot is inaccurate. The slicer
prints a prominent **APPROXIMATION WARNING** banner in `sliced.md`
when it operates in this mode; the cutoff timestamp itself remains
reliable.

### C3 — Spawn triager subagents

For each BF, spawn ONE triager subagent using the runtime's
subagent-dispatch tool. In both runtimes, emitting all N dispatch
calls in a single assistant message is what produces true
parallelism — the runtime executes them concurrently.

- **Cursor**: tool is named `Task`. Pass
  `subagent_type: "generalPurpose"` (camelCase) and
  `run_in_background: true` on each `Task` call.
- **Claude Code** (v2.1.63+): tool is named `Agent` — the legacy
  `Task` name no longer resolves to a distinct tool. Pass
  `subagent_type: "general-purpose"` (note the hyphen, NOT
  Cursor's `generalPurpose`) and `run_in_background: true` on
  each `Agent` call. Both parameters are supported.

Critical: the coordinator must pass **only the file paths**, not the
BF key, so the subagent has no straightforward way to fetch fresh Jira
state. The triager prompt template requires substitutions:

| Placeholder | Replace with |
| ----------- | ------------ |
| `{{SLICED_PATH}}` | `$RUN_DIR/sliced.md` |
| `{{OUTPUT_PATH}}` | `$RUN_DIR/triage.md` |

Read the prompt template from
[`../templates/triager_prompt.md`](../templates/triager_prompt.md) and
inline it into the subagent's prompt verbatim with the substitutions
applied.

Spawn ALL N triagers in a single tool-call batch (one parallel
dispatch call per BF — `Task` under Cursor, `Agent` under Claude
Code) so they run concurrently.

### C4 — Spawn grader subagents

After **all** triager subagents have completed, spawn ONE grader
subagent per BF. The grader template lives at
[`../templates/grader_prompt.md`](../templates/grader_prompt.md) and
requires substitutions:

| Placeholder | Replace with |
| ----------- | ------------ |
| `{{REPORT_PATH}}` | `$RUN_DIR/triage.md` |
| `{{HELDOUT_PATH}}` | `$RUN_DIR/heldout.md` |
| `{{SLICED_PATH}}` | `$RUN_DIR/sliced.md` |
| `{{CUTOFF_PATH}}` | `$RUN_DIR/cutoff.txt` |
| `{{SCORE_PATH}}` | `$RUN_DIR/score.v2.md` |
| `{{BF_KEY}}` | `BF-NNNNN` (header text only — grader is allowed to know it) |

Spawn all graders in a single tool-call batch using the same
runtime-specific dispatch tool as C3 (Cursor: `Task` with
`subagent_type: "generalPurpose"` + `run_in_background: true`;
Claude Code: `Agent` with `subagent_type: "general-purpose"` +
`run_in_background: true`). Multiple calls in one assistant
message gives concurrency in both runtimes.

### C5 — Aggregate

After all graders finish, the coordinator reads each `score.v2.md` and
writes a single aggregate report under the same directory as
Mode A / Mode B reports — i.e. `${BF_TRIAGE_OUTPUT_DIR:-./bf-reports}/`.
Default filename pattern:
`verification-<YYYYMMDDTHHMMSSZ>.md` for multi-BF runs, or
`verification-<BF-KEY>-<YYYYMMDDTHHMMSSZ>.md` for a single-BF replay.
Follow the Step 9 versioning rule (`-v2`, `-v3`, …) if the path
already exists. The aggregate should contain:

- Per-BF summary table: `| BF | Routing | Root-cause | Next steps | Overall | Leak audit |`.
- Pattern notes lifted from each grader's "Notes for the coordinator"
  section.
- Recommended skill refinements (skill rule / reference-file updates
  that would close the observed gaps).

### C6 — Cleanup

By default, leave `$RUN_DIR` in place so the user can inspect the
artifacts. Only delete `$RUN_DIR` itself when the user explicitly
asks. **NEVER touch any sibling `/tmp/bf-triage-test/<BF>*` folder
not owned by this run** — those belong to other concurrent or past
agents.

The triager follows `SKILL.md` Step 13 and cleans up its own
`/tmp/bf-triage-workdir-<BF>/` scratch dir for raw EVG logs / git
worktrees.

## Triager isolation rules (full list)

These live in [`../templates/triager_prompt.md`](../templates/triager_prompt.md)
and are repeated here as a quick reference. The triager **must not**:

1. Call any `jira_*` MCP tool against the BF under test (incl. BF key
   in JQL `text ~`).
2. Call `bb_get_bf` against the BF under test.
3. `web_search` / `web_fetch` of `https://jira.mongodb.org/browse/<BF>`
   or any URL that resolves to the BF Jira page.
4. `glean.search` queries containing the BF key.
5. Read any `heldout.md` / `score.v2.md` / any file path containing
   the substring `heldout` for any ticket.
6. Run any local `git` command against the user's working
   `10gen/mongo` / `10gen/dsi` checkout. When gateway `git_*` fails,
   the subagent MUST use a `setup_repos.sh` scratch clone — Step 1's
   normal "env vars → agent-provided workspaces → scratch clone"
   resolution order is overridden inside subagents to skip directly
   to the scratch clone. Rationale: subagents have no live human supervisor
   to interrupt a `git checkout` / `git reset` mistake; the user's
   working repo is unrecoverable damage. The triager prompt template
   carries the full forbidden-command list.

The triager **may**:

- Read `sliced.md` for ticket content.
- Use gateway `git_*` first; fall back to local `git` against a
  `setup_repos.sh` scratch clone (NOT the user's working repo — see
  forbidden item 6 above) on failure.
- Use `evergreen` CLI for log fetches (CLI-first per `SKILL.md` Step
  5.5).
- Call `bb_get_bfg_by_task` IFF the sliced file contains an Evergreen
  task ID.
- Call MCP-only EVG endpoints with task IDs that appear in the slice.
- Call `glean.search` for runbook / wiki search WITHOUT the BF key.

## Grader rubric (summary)

Full rubric and citation format in
[`../templates/grader_prompt.md`](../templates/grader_prompt.md). Each
report is scored on three dimensions, each PASS / PARTIAL / MISS /
INSUFFICIENT_SIGNAL:

1. **Routing accuracy** — does "Recommended next steps" point at the
   actual owning team or fix surface?
2. **Root-cause direction** — does "Root cause hypothesis" name the
   actual mechanism or a plausible parent of it?
3. **Actionable next steps** — are recommendations specific enough to
   act on, AND do they overlap with what actually happened?

Plus a **held-in leak audit**: any quote in the triage report that
could only have come from `heldout.md` is grounds to downgrade
Routing / Root-cause to MISS, regardless of accuracy.

## File map for verification

| File | Role |
| ---- | ---- |
| `scripts/run_held_in_test.sh` | Slice one BF's pre-fetched JSON into `sliced.md` / `heldout.md` / `cutoff.txt` |
| `scripts/run_batch.sh` | Parallel slicer for N BFs |
| `scripts/_slice_helper.py` | Internal markdown renderer |
| `scripts/_stitch_gateway_jira.py` | Normalizes gateway two-call shape into legacy single-blob shape (stop-gap until the gateway exposes `jira_get_issue_changelog`) |
| `scripts/README.md` | Coordinator pre-flight checklist + sequence diagram |
| `templates/triager_prompt.md` | Subagent prompt for held-in triager |
| `templates/grader_prompt.md` | Subagent prompt for grader |

## Promoting findings to skill content

When a verification run discovers a new pattern, classify it before
promoting to skill content. The classification rule is:

- **Shared file** (`SKILL.md`, `reference/workflow_overview.md`,
  `reference/log_patterns.md`, `reference/severity_types.md`,
  `reference/tools_inventory.md`) — patterns that hold true for
  **any** team using the skill. Concretely:
  - A new MongoDB log marker / fault string.
  - A new severity-type auto-detect identifier or manual identifier.
  - A new generic workflow rule (e.g. "always check field X before Y").
  - A new tool / CLI integration tip, gateway-failure recovery
    procedure, or descriptor-cache path.
  - A new generic re-routing matrix row (sub-system → owner team,
    where the owner is named by sub-system, not by the current
    skill user's team).
- **Team file** (`reference/team_knowledge.md`) — patterns that
  encode **who** owns a symptom **front-line**, **how the current
  team specifically routes it**, or **which downstream projects the
  current team links fix tickets in**. Concretely:
  - New front-line routing rules ("this team owns symptoms of shape
    X").
  - New error-class examples the current team specifically sees.
  - New variant-family names the current team handles.
  - New workload-path conventions or downstream project prefixes
    (`PERF-`, `TUNE-`, etc.).
  - Disposition defaults specific to this team's resolution policy.

### Heuristic checklist

Apply in order:

1. If the new rule names a team (any team, not just the current one)
   in its body or its conclusion → **team file** for the current team,
   or generic-matrix row in `workflow_overview.md` if the named team
   is a *sub-system* team (Replication, Sharding, Query, etc.) rather
   than the *current owning* team.
2. If the new rule names a non-`SERVER-` JIRA project prefix as a
   downstream → **team file** (the prefix is team-specific).
3. If the new rule names a specific workload, variant family, or DSI
   subdirectory → **team file** (the names are team-domain-specific).
4. If the new rule reads naturally without any team name, project
   prefix, or workload name in it → **shared file**.

### Workflow

After a verification run finds a new pattern:

1. Write a one-paragraph description of the new rule.
2. Apply the heuristic above to decide the target file.
3. **Ask the user before editing the skill files** — the agent must
   not silently mutate the skill in response to a single sample's
   evidence. Surface the proposed addition with target-file and the
   classification reason, and wait for approval.
4. On approval: edit the chosen file. If team-specific, also update
   any cross-references in shared files that should point at the new
   section.
5. Record the addition in the verification run's final report so the
   user can audit later what the skill learned.

This deliberately keeps the skill conservative: a single failed grade
on one BF should not auto-rewrite the skill's rules. Multiple
verification runs converging on the same new pattern is the bar for
promotion.

## Known limitations

1. **No changelog from the gateway.** Snapshots are best-effort from
   current `custom_fields`. The slicer's APPROXIMATION WARNING banner
   covers this; a structural fix requires a new gateway tool
   (`jira_get_issue_changelog`).
2. **Antithesis BFs lack an Evergreen task ID.** The triager will hit
   a dead end on Step 5 log fetch and degrade to the limited-evidence
   path. This is a real and rare case (≤10% of WR's queue, AI
   concluded from BF-37519); no special-casing in the skill.
3. **Verification quality depends on `heldout.md` density.** If the
   BF was closed with minimal follow-up (e.g. "Won't Fix" with no
   investigation notes), the grader may legitimately mark dimensions
   as `INSUFFICIENT_SIGNAL`. That is a data outcome, not a triager
   failure.
