---
name: atlas-infinite-release-static-analysis
description: >-
  Use when running the static (pre-patch) analysis for an upcoming Atlas Infinite
  (DSC/Disaggregated Storage Cluster) release. Audits the SDAP (Security,
  Durability, Availability, Performance) git diff of a mongod master HEAD commit
  against the prior release branch, reviews the diff for introduced security
  vulnerabilities and likely bugs, and scans Jira and the #buildbaron Slack
  channel for open DSC issues. Produces a static-analysis release-confidence
  report. Does not trigger Evergreen patches.
argument-hint: <release-number>
arguments: release_version
allowed-tools: mcp__devprod-mcp-gateway__jira_search_issues mcp__devprod-mcp-gateway__jira_get_issue mcp__claude_ai_Glean_via_MCP__search Bash(git *) Bash(gh *)
source: 10gen/agent-skills
license: Internal
mongodb:
  team: RSSD
  owner: jaroslaw.kawa@mongodb.com
  internal: true
---

# Atlas Infinite Release — Static Analysis

Runs the static, pre-patch analysis for an upcoming Atlas Infinite (DSC) release. Audits the
SDAP git diff against the prior release branch, reviews the diff for introduced security
vulnerabilities and likely bugs, and scans Jira and #buildbaron for open DSC issues. Produces
a static-analysis release-confidence score with a scoped recommendation on whether to advance
to patch analysis.

This is **step 1 of the two-part release qualification**. Per the release runbook: run this
skill, review the confidence report, and only if it reports clear (or the concerns are
accepted) run `atlas-infinite-release-patch-analysis` to trigger and analyze the Evergreen
patch builds.

**Trigger examples:**
- `/atlas-infinite-release-static-analysis 5` — analyze for release 5
- "Run the static analysis for DSC release 5"
- "Analyze the DSC release diff and open issues" — skill auto-detects the latest release branch

## Prerequisites

- **Repo and branch** — must be run from the `mongodb/mongo` repo root on the `master` branch.
- **Jira and Glean access** — via the devprod MCP gateway and the Glean MCP. If a scan returns
  a permissions error, note it in the report rather than failing the whole run.

Create a todo list before starting with one item per phase. Mark each complete as it finishes.

## Phase 1 — Resolve release context

Establish the diff anchors before any other work:

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

4. **Prior release cut commit** — tip of the prior release branch:
   ```bash
   git rev-parse origin/dsc-release-<N-1>
   ```

Print a confirmation block and wait for user approval before proceeding:

```
Analyzing commit : <master HEAD SHA> (master)
Prior release    : origin/dsc-release-<N-1> @ <prior SHA>
Upcoming release : dsc-release-<N> (or "TBD — user to confirm")
```

## Phase 2 — SDAP git diff analysis

Run this concurrently with Phase 4 (signal gathering) — do not block one on the other.

**Collect the commit log and diff:**
```bash
git log <prior-SHA>..HEAD --oneline --no-merges > /tmp/dsc-static-commits.txt
git diff <prior-SHA>..HEAD -- src/ > /tmp/dsc-static-diff.patch
```

Scan `/tmp/dsc-static-commits.txt` and `/tmp/dsc-static-diff.patch` for changes in the
following SDAP areas. For each area, report the relevant commit SHAs, file paths touched, and
a risk assessment:

| Area | File/symbol signals to look for |
|---|---|
| **Security** | `src/mongo/db/auth/`, `src/mongo/util/net/ssl*`, `src/mongo/db/audit*`, `encryption`, `TLS`, `x509`, `LDAP` |
| **Durability** | `src/mongo/db/repl/oplog*`, `src/mongo/db/storage/`, `journal`, `checkpoint`, `crash_recovery`, `WT_TXN` |
| **Availability** | `src/mongo/db/modules/atlas/` — DSC replication, leader/follower state transitions, membership changes, failover, availability-related coordinator logic |
| **Performance** | `src/mongo/db/exec/`, `src/mongo/db/query/`, `src/mongo/db/index/`, `eviction`, `cache`, `bulk_write` |

Also flag any commits whose messages contain: `revert`, `CRIT`, `data loss`, `corruption`,
`regression`, `security fix`, or `CVE`.

Output a risk table:

| SDAP Area | # Commits | Key Files Changed | Risk |
|---|---|---|---|
| Security | N | ... | 🟢 Low / 🟡 Medium / 🔴 High |
| Durability | N | ... | ... |
| Availability | N | ... | ... |
| Performance | N | ... | ... |

Risk calibration:
- 🔴 **High** — direct changes to critical-path code with no test coverage visible in the diff, or a revert/regression commit
- 🟡 **Medium** — changes that touch the area but appear guarded, tested, or behind a feature flag
- 🟢 **Low** — no changes, or only comments/cosmetic changes

## Phase 3 — Security & bug-introduction review

Beyond the SDAP area mapping, review the diff content in `/tmp/dsc-static-diff.patch` for
vulnerabilities and bugs *introduced* by these changes — independent of which area they touch.

Look for:
- **Security** — injection (command / JS / format-string), authentication or authorization
  bypass, missing validation on external input, unsafe deserialization, secrets or credentials
  committed in code, weakened TLS/crypto, integer overflow feeding an allocation or index, and
  use-after-free / double-free / buffer overrun in C++.
- **Correctness** — unchecked error returns, inverted or off-by-one conditions, missing
  null/bounds checks, a lock acquired but not released on an early return, data races on shared
  state, resource leaks, and changed invariants whose assertions were not updated.

For each finding capture file:line, the issue, why it is exploitable or wrong, and severity:

| File:line | Finding | Severity |
|---|---|---|
| ... | ... | 🔴 High / 🟡 Suspected / 🟢 Note |

Severity:
- 🔴 **High** — likely exploitable, or can cause data loss / corruption / crash on a reachable path
- 🟡 **Suspected** — a plausible bug or weakness that needs author confirmation
- 🟢 **Note** — minor or defensive-only; informational

If the diff is large, prioritize files in the SDAP-flagged areas from Phase 2 and any C++
touching memory, concurrency, or parsing. **Note explicitly if any portion was not fully
reviewed** so the score is not read as full coverage.

## Phase 4 — Signal gathering (run concurrently with Phase 2)

Run all four queries below in parallel — three Jira searches and one Glean search. None depend
on each other or on Phase 2/3 completing first.

### Jira scan

Use `mcp__devprod-mcp-gateway__jira_search_issues` for each query:

**Build Failures:**
```
project = BF AND statusCategory != Done AND (summary ~ "disagg" OR summary ~ "DSC" OR text ~ "disaggregated")
```

**Automation Failures / Server bugs:**
```
project = SERVER AND statusCategory != Done AND (summary ~ "disagg" OR summary ~ "DSC" OR labels in ("DSC", "disagg", "disaggregated-storage"))
```

**HELP tickets:**
```
project = HELP AND statusCategory != Done AND (summary ~ "disagg" OR summary ~ "DSC")
```

For each ticket returned, fetch full details with `mcp__devprod-mcp-gateway__jira_get_issue`
and capture: key, summary, status, priority, assignee, and affected versions.

Build a table per category and assign a release risk label:

| Ticket | Summary | Status | Priority | Assignee | Release Risk |
|---|---|---|---|---|---|

Release risk labels:
- 🔴 **Blocking** — high/critical priority, actively impacting DSC stability, no workaround
- 🟡 **Watch** — medium priority, or has a known workaround, or already being fixed
- 🟢 **Informational** — low priority, cosmetic, or not affecting the release scope

### Build Baron Slack scan

Use `mcp__claude_ai_Glean_via_MCP__search` to scan #buildbaron for signals that could affect
the release.

**Determine the time window** before running the searches:
- If the prior release cut commit is known from Phase 1, resolve its date:
  ```bash
  git log -1 --format=%cd --date=format:'%Y-%m-%d' <prior-SHA>
  ```
  Use that date as `<scan-start-date>`.
- If the date cannot be resolved, default `<scan-start-date>` to one month before today.

Run the following searches in parallel:

**Reverts flagged in buildbaron since scan start:**
```
channel:buildbaron revert after:<scan-start-date>
```

**DSC / disagg flags raised in buildbaron since scan start:**
```
channel:buildbaron DSC disagg atlas infinite after:<scan-start-date>
```

For each result, capture: message date, author, a brief summary of what was flagged, and
whether it appears to be resolved or still active.

Produce a table of notable findings:

| Date | Author | Signal | Status |
|---|---|---|---|

Signal risk labels (same scale as Jira):
- 🔴 **Blocking** — active revert or flag directly targeting DSC code, unresolved
- 🟡 **Watch** — flag raised but appears resolved, or adjacent to DSC rather than directly affecting it
- 🟢 **Informational** — noise, unrelated, or already reflected in a Jira ticket found above

If Glean returns no results, note that explicitly — it may mean the channel wasn't indexed or
that there were genuinely no flags in the window.

## Phase 5 — Static-analysis confidence report

Combine Phases 2–4 into a **static-analysis confidence score** starting from 100:

| Factor | Max deduction | Trigger |
|---|---|---|
| SDAP diff risk | −40 | −20 per 🔴 High area, −10 per 🟡 Medium area |
| Security / bug-introduction findings | −25 | −15 per 🔴 High finding, −7 per 🟡 Suspected finding |
| Blocking Jira issues | −25 | −12 per 🔴 Blocking ticket, −5 per 🟡 Watch ticket |
| Build Baron flags | −10 | −7 per 🔴 Blocking signal, −3 per 🟡 Watch signal |

Floor the score at 0. Round to the nearest integer.

**Output the report in this format:**

```
## Atlas Infinite Release — Static Analysis Report
**Analyzing commit :** <SHA> (master)
**Prior release    :** dsc-release-<N-1> @ <SHA>
**Report date      :** <today>

### 🔴 Blocking findings
<one line per 🔴 item across SDAP / security / Jira / Build Baron; or "None">

### SDAP Diff Risk
Security     : 🟢/🟡/🔴  (<N commits>)
Durability   : 🟢/🟡/🔴  (<N commits>)
Availability : 🟢/🟡/🔴  (<N commits>)
Performance  : 🟢/🟡/🔴  (<N commits>)

### Security & Bug-Introduction Findings
High      : N  [list file:line — finding]
Suspected : N  [list]
Notes     : N
Coverage  : <full / partial — note any area not fully reviewed>

### Open Jira Issues
Blocking     : N
Watch        : N
Informational: N

### Build Baron (#buildbaron)
Blocking     : N  [list signals]
Watch        : N  [list signals]
Informational: N

### Static Analysis Confidence: XX%

<2–3 sentence narrative explaining the main factors driving the score>

### Recommendation
✅ CLEAR / ⚠️ CONCERNS / 🔴 BLOCKERS

<Rationale: what's clean, what's at risk, follow-up actions. State the next step: if CLEAR
or the CONCERNS are accepted, run `atlas-infinite-release-patch-analysis` to trigger the
Evergreen patch builds.>
```

Thresholds:
- **85–100%** → ✅ CLEAR — safe to advance to patch analysis
- **65–84%** → ⚠️ CONCERNS — list the specific caveats to weigh before advancing
- **< 65%** → 🔴 BLOCKERS — list what must be resolved before triggering patches
