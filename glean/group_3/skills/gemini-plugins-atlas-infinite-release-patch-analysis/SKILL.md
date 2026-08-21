---
name: atlas-infinite-release-patch-analysis
description: >-
  Use when running the Evergreen patch analysis for an upcoming Atlas Infinite
  (DSC/Disaggregated Storage Cluster) release. Triggers 3 Evergreen patch builds
  of a mongod master HEAD commit, classifies flaky vs. consistent failures across
  the three runs, traces each consistent failure to a suspect commit and BF
  ticket, and produces a patch-analysis release-confidence report. Run after the
  static-analysis skill reports clear.
argument-hint: <release-number>
arguments: release_version
allowed-tools: mcp__devprod-mcp-gateway__evg_get_patch_failed_jobs mcp__devprod-mcp-gateway__evg_get_task_log_summary mcp__devprod-mcp-gateway__evg_get_test_results_summary mcp__devprod-mcp-gateway__evg_get_test_results_detailed mcp__devprod-mcp-gateway__evg_list_user_recent_patches mcp__devprod-mcp-gateway__jira_search_issues Bash(evergreen *) Bash(git *) Bash(gh *)
source: 10gen/agent-skills
license: Internal
mongodb:
  team: RSSD
  owner: jaroslaw.kawa@mongodb.com
  internal: true
---

# Atlas Infinite Release — Patch Analysis

Triggers and analyzes the Evergreen patch builds for an upcoming Atlas Infinite (DSC) release.
Three parallel patch builds of master HEAD catch flaky vs. consistent failures; each consistent
failure is traced to a suspect commit and an existing BF ticket. Produces a patch-analysis
release-confidence score with a scoped recommendation.

This is **step 2 of the two-part release qualification**. Per the release runbook, run this
only after `atlas-infinite-release-static-analysis` has reported clear and the engineer has
reviewed its findings.

**Trigger examples:**
- `/atlas-infinite-release-patch-analysis 5` — analyze patches for release 5
- "Run the DSC release patch analysis for release 5"
- "Trigger the 3 qualification patches for master HEAD" — skill auto-detects the latest release branch

## Prerequisites

- **Evergreen auth** — run `evergreen login` and confirm it completes before invoking this
  skill. The CLI serializes OAuth token access via a lock file; an unauthenticated or expired
  session causes all patch submissions to hang or fail with "timed out waiting for OAuth token
  lock." Re-run `evergreen login` if symptoms appear (CLI prints "navigate to verification URI"
  and stalls). Only one `evergreen` process should run at a time while the token is being
  acquired — the skill submits patches sequentially for this reason.
- **Patch priority** — the Evergreen CLI does not expose a `--priority` flag. Without elevated
  priority, patches run at default priority and can take ~5 hours to complete. After the three
  patch URLs are reported in Phase 2, open each URL in Spruce and set priority to **100** via
  the UI (three-dot menu → Set Priority) before detaching. This keeps total wall-clock time
  under 2 hours.
- **Repo and branch** — must be run from the `mongodb/mongo` repo root on the `master` branch.

Create a todo list before starting with one item per phase. Mark each complete as it finishes.

## Phase 1 — Resolve release context and verify Evergreen auth

Establish the anchors before any other work:

0. **Verify Evergreen auth** — confirm the CLI has a valid session before doing anything else:
   ```bash
   evergreen list-patches -n 1
   ```
   - If the command succeeds and returns patch data, proceed.
   - If it hangs, prints a device-auth URL, or errors with "timed out waiting for OAuth token
     lock", **stop immediately** and return this message to the user:

   > **Evergreen authentication required.**
   > Your Evergreen CLI session has expired or was never set up.
   > Run the following command in your terminal and complete the browser login before re-running this skill:
   > ```
   > evergreen login
   > ```
   > Once the login completes successfully, restart the patch analysis.

1. **Verify branch and pull latest** — confirm the current branch is `master`, then pull the
   latest changes:
   ```bash
   git branch --show-current
   git pull origin master
   ```
   If the branch is not `master`, warn the user and stop — only master HEAD is valid to
   qualify. If the pull fails, surface the error and stop.

2. **Master HEAD SHA** — capture after the pull so the SHA reflects the latest commit:
   ```bash
   git rev-parse HEAD
   ```

3. **Prior release branch** — release branches follow the form `dsc-release-N` where N is an
   incrementing integer (e.g., `dsc-release-1`, `dsc-release-2`):
   - If `$release_version` was provided (e.g., `3`): the prior branch is
     `origin/dsc-release-<N-1>` (e.g., `origin/dsc-release-2`). If not found, list available
     branches and ask the user to confirm:
     ```bash
     git branch -r | grep dsc-release | sort -V
     ```
   - If no argument was provided: detect the most recently cut release branch:
     ```bash
     git branch -r | grep dsc-release | sort -V | tail -1
     ```
     Confirm the detected branch with the user before proceeding.

4. **Prior release cut commit** — tip of the prior release branch (used in Phase 4 to trace
   suspect commits):
   ```bash
   git rev-parse origin/dsc-release-<N-1>
   ```

Print a confirmation block and wait for user approval before triggering any patch:

```
Analyzing commit : <master HEAD SHA> (master)
Prior release    : origin/dsc-release-<N-1> @ <prior SHA>
Upcoming release : dsc-release-<N> (or "TBD — user to confirm")
```

## Phase 2 — Trigger 3 Evergreen patches (sequential, not parallel)

Run the following command **three times sequentially** — wait for each `Bash` tool call to
complete before issuing the next. Do **not** run them in parallel: the Evergreen CLI serializes
OAuth token access via a lock file, and concurrent invocations will race on that lock, causing
all but the first to time out. Running them one after another avoids the conflict while still
producing three independent patch IDs. Replace the placeholder with the actual command when
available:

```bash
evergreen patch -f -y -p mongodb-mongo-master \
  -v al2023-arm64-benchmarks \
  -v amazon-linux2023-arm64-static-compile \
  -v amazon-linux2023-arm64-crypt-compile \
  -v linux-arm64-debug-compile-required \
  -v linux-debug-aubsan-lite-all-feature-flags-required \
  -v enterprise-amazon-linux2023-arm64-all-feature-flags \
  -v enterprise-amazon-linux2023-arm64-all-feature-flags-extra-system-deps \
  -v enterprise-amazon-linux2023-arm64-all-feature-flags-extra-system-deps-sharded-clusters \
  -v generate-tasks-for-version \
  -t all \
  -d "DSC release patch analysis run <N> — master @ $(git rev-parse --short HEAD)"
```

Run N=1, N=2, N=3 sequentially. Each invocation produces output containing:
```
Created patch: <patch-id>
https://spruce.mongodb.com/version/<patch-id>
```

Capture the three patch IDs from the command output. **Do not write them to `/tmp` or any
local file.** `/tmp` is frequently a tmpfs that is wiped on reboot (and is subject to tmpfiles
cleanup), so over the 2–5h patch window a state file there can vanish before the resume fires.
Instead, carry all resume state inline in the `ScheduleWakeup` prompt below — that makes it
exactly as durable as the scheduled wake-up itself, with no extra failure point. The patch URLs
reported to the user below are the human-recoverable copy of the IDs.

Report the patch URLs and priority reminder to the user:
```
Patches submitted — Evergreen builds are running.
Set priority to 100 in Spruce for each patch (three-dot menu → Set Priority) to keep
wall-clock time under 2 hours. The patch-analysis report will be delivered once all three
patches complete.

Patch 1: https://spruce.mongodb.com/version/<id1>
Patch 2: https://spruce.mongodb.com/version/<id2>
Patch 3: https://spruce.mongodb.com/version/<id3>
```

Then use `ScheduleWakeup` to schedule Phase 3 in 5 hours (or 2 hours if priority was elevated).
Inline the captured patch IDs and release anchors into the prompt so the resume needs no local
file, and instruct it to query patches **by ID** rather than by recency:

```
ScheduleWakeup(
  delaySeconds: 18000,
  reason: "DSC patch analysis running — checking results in 5h",
  prompt: "Resume DSC release patch analysis. Patch IDs: <id1>, <id2>, <id3>. Anchors: HEAD_SHA=<head-sha>, PRIOR_SHA=<prior-sha>, RELEASE=<release_version>. Check each patch by its own ID with `evergreen list-patches -i <id> -j` (one call per ID — do NOT use `-n` recency, which other manual or CI patches can clobber). If all three are terminal (finish_time set, status succeeded or failed), run Phase 4 failure analysis and Phase 5 confidence report. If any is still running, reschedule with ScheduleWakeup(delaySeconds: 3600), re-passing this same prompt verbatim."
)
```

**Do not wait or poll manually.** The session ends here until the scheduled wake-up fires.

## Phase 3 — Resume: check patch completion (runs on wake-up)

The wake-up prompt carries the three patch IDs and the release anchors (HEAD_SHA, PRIOR_SHA,
RELEASE) inline — recover them from the prompt. No state file is written, so do not read `/tmp`.

Check each patch **by its specific ID**, never by recency. `evergreen list-patches -n <k>`
returns the k most recent patches for the user, which any other manual or CI-triggered patch
submitted in the interim can clobber; querying by ID is immune to that. Run one call per
captured ID (in parallel). The CLI is more reliable than MCP for user-owned patches on
`mongodb-mongo-master`:

```bash
evergreen list-patches -i <patch_id> -j
```

A patch is still running if `"status": "started"` or `"finish_time": null`.
A patch is terminal if `"status"` is `"succeeded"` or `"failed"` and `"finish_time"` is set.

If the MCP `evg_get_patch_failed_jobs` is available and has permissions, use it (keyed on the
same patch ID) as a supplement for per-task failure detail — but do not depend on it exclusively.

**If any patch is still running:** report a brief status update to the user, then reschedule
using `ScheduleWakeup` with `delaySeconds: 3600`, re-passing the same prompt verbatim so the IDs
and anchors carry forward. Do not proceed to Phase 4. The status update should follow this format:

```
⏳ DSC patch analysis still running — checking back in ~1 hour.

Patch 1: <status> — https://spruce.mongodb.com/version/<id1>
Patch 2: <status> — https://spruce.mongodb.com/version/<id2>
Patch 3: <status> — https://spruce.mongodb.com/version/<id3>

Time elapsed: ~<N>h since submission. The report will be delivered once all three complete.
```

**If all patches are terminal**, collect per patch via CLI and MCP:
- Total task count
- Failed task names and build variants
- Pass/fail summary

Then proceed immediately to Phase 4.

## Phase 4 — Failure analysis

Build a failure matrix across the 3 patches. For each unique `task + build_variant` combination:

| Task | Build Variant | Patch 1 | Patch 2 | Patch 3 | Classification |
|---|---|---|---|---|---|
| `<suite>` | `<variant>` | ❌ | ✅ | ❌ | Flaky |
| `<suite>` | `<variant>` | ❌ | ❌ | ❌ | Consistent |

**Classification rules:**
- Failed in **1/3** patches → **Flaky** — note but do not block release on this alone
- Failed in **2/3** patches → **Likely consistent** — investigate; treat as a soft blocker
- Failed in **3/3** patches → **Consistent failure** — high signal; treat as a hard blocker

For every **Consistent** or **Likely consistent** failure:

1. Fetch test-level detail:
   `mcp__devprod-mcp-gateway__evg_get_test_results_summary` → identify specific failing tests
2. Fetch task log:
   `mcp__devprod-mcp-gateway__evg_get_task_log_summary` → extract the error message/stack
3. Find likely responsible commit:
   ```bash
   git log <prior-SHA>..HEAD --oneline -- <affected-file-or-suite-path>
   ```
4. Search for an existing BF ticket:
   `mcp__devprod-mcp-gateway__jira_search_issues` with `project = BF AND summary ~ "<suite name>"`

Report each consistent/likely-consistent failure in this block format:

```
TASK:    <task name> / <build variant>
TESTS:   <failing test names>
ERROR:   <error snippet from logs>
SUSPECT: <SHA> — <commit message>
BF:      <existing BF key> or "none found — recommend filing"
```

## Phase 5 — Patch-analysis confidence report

Compute a **patch-analysis confidence score** from the failure matrix, starting from 100:

| Factor | Max deduction | Trigger |
|---|---|---|
| Consistent failures (3/3) | −60 | −20 per consistent failure |
| Likely-consistent failures (2/3) | −30 | −12 per likely-consistent failure |
| Flaky failures (1/3) | −10 | −3 per flaky task (capped at −10) |

Floor the score at 0. Round to the nearest integer.

**Output the report in this format:**

```
## Atlas Infinite Release — Patch Analysis Report
**Analyzing commit :** <SHA> (master)
**Prior release    :** dsc-release-<N-1> @ <SHA>
**Report date      :** <today>

### Evergreen Results (3 patches)
Consistent failures : N  [list task names]
Likely-consistent   : N  [list task names]
Flaky failures      : N  [list task names]
Total tasks run     : N

### Consistent / Likely-Consistent Failures
<the per-failure blocks from Phase 4, or "None">

### Patch Analysis Confidence: XX%

<2–3 sentence narrative explaining the main factors driving the score>

### Recommendation
✅ CLEAR / ⚠️ CONCERNS / 🔴 BLOCKERS

<Rationale: what failed consistently, suspected causes, and follow-up actions (file BFs,
re-run after a fix). Combine this with the static-analysis report for the final ship decision.>
```

Thresholds:
- **85–100%** → ✅ CLEAR — patches clean
- **65–84%** → ⚠️ CONCERNS — list the specific caveats the team should accept
- **< 65%** → 🔴 BLOCKERS — list the consistent failures that must be resolved before re-running
