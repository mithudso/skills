---
name: task-failure-triage
description: >-
  Use when an Evergreen task fails and the cause isn't obvious. Determines
  whether the failure is from your project's code/config or a genuine
  Evergreen infrastructure issue, before filing a DevProd ticket.
source: 10gen/evergreen
license: Internal
mongodb:
  team: devprod-bv
  owner: chaya.malik@mongodb.com
  internal: true
  allow_large_context: true
---

# Evergreen Task Failure Triage

## Purpose

Determine whether an Evergreen task failure is caused by **your project's code or configuration**, or by a **genuine Evergreen infrastructure issue** that requires attention from the DevProd team. Run this before filing a DevProd ticket.

Most task failures are caused by the project's own scripts, test code, or task configuration. This skill works through the available evidence — the Evergreen REST API, your project YAML, and the scripts your task runs — to reach a verdict.

## When to use

- Your Evergreen task failed and the cause isn't obvious
- You see missing logs, incomplete results, or an unclear error message
- You want to know if the issue is on your end before asking for help

## When NOT to use

- You already know the failure is in your test code
- The task is still running (wait for it to finish)
- The issue is with Evergreen's UI, permissions, or project settings (file a ticket directly)

## Hard rule: no writes, no escalation

This skill diagnoses only. It does not file tickets, restart tasks, or modify anything. All it produces is a verdict and a recommendation.

---

## Phases

### 1. Extract the task ID

Regex-scan the user's input for:
- Spruce URLs: `spruce.corp.mongodb.com/task/<task_id>` → capture `task_id`
- Parsley URLs: `parsley.corp.mongodb.com/evergreen/<task_id>` → capture `task_id`
- Bare task IDs: look like `<project>_<variant>_<task>_<sha>_<YY_MM_DD_HH_MM_SS>`

If no task ID is present, ask the user to provide a Spruce or Parsley URL.

### 2. Obtain an OAuth token

Static API keys (`~/.evergreen.yml`) are deprecated for human users. Use OAuth instead.

Get a token with the Evergreen CLI:

```bash
evergreen client get-oauth-token
```

If the command fails or returns an empty string, the session is expired. Refresh it first:

```bash
evergreen login
```

Then retry `evergreen client get-oauth-token`.

**Security rule — strictly enforced:** Never print, echo, or display the raw token value in any command output, response text, or tool call. Capture it in a shell variable and pass it inline to subsequent commands only.

All REST API calls must use `https://evergreen.corp.mongodb.com` (not `evergreen.mongodb.com`) and the `Authorization: Bearer` header:

```bash
python3 - <<'EOF'
import subprocess, json

token = subprocess.run(
    ["evergreen", "client", "get-oauth-token"],
    capture_output=True, text=True
).stdout.strip()

task_id = "<TASK_ID>"  # substitute the actual task ID
r = subprocess.run(
    ["curl", "-s",
     "-H", f"Authorization: Bearer {token}",
     f"https://evergreen.corp.mongodb.com/rest/v2/tasks/{task_id}"],
    capture_output=True, text=True
)
print(r.stdout)
EOF
```

### 3. Fetch task metadata from the Evergreen REST API

Always use the Python wrapper to keep the token out of shell history and output:

```bash
python3 - <<'EOF'
import subprocess, json

token = subprocess.run(
    ["evergreen", "client", "get-oauth-token"],
    capture_output=True, text=True
).stdout.strip()
headers = ["-H", f"Authorization: Bearer {token}"]
task_id = "<TASK_ID>"           # substitute actual task ID
execution = ""                  # set to "?execution=N" if the URL had that param

r = subprocess.run(["curl", "-s"] + headers +
    [f"https://evergreen.corp.mongodb.com/rest/v2/tasks/{task_id}{execution}"],
    capture_output=True, text=True)
task = json.loads(r.stdout)

# NOTE: the REST API uses "status_details" not "details"
fields = ["status", "display_status", "status_details", "start_time", "finish_time",
          "host_id", "distro_id", "build_variant", "display_name", "revision",
          "project_identifier", "execution", "logs", "parsley_logs"]
print(json.dumps({k: task.get(k) for k in fields}, indent=2))
EOF
```

**Critical fields to read first:**

The REST API returns failure info under `status_details` (not `details` — that's the internal DB field name).

| Field | What it means |
|---|---|
| `status_details.type` | `"test"` = agent alive, command exited non-zero. `"system"` = agent died or host lost. `"setup"` = pre-task failed. |
| `status_details.desc` | Agent's one-line description of why it failed. |
| `status_details.failing_command` | Which Evergreen function/step failed (format: `'<command>' in function '<name>' (step N of M)`). |
| `status_details.other_failing_commands` | Additional failures — often post-task commands. |
| `status_details.post_errored` | `true` means the post-task block also had failures. |
| `status_details.oom_tracker_info.detected` | `true` if the kernel OOM-killed a process during the task. |
| `status_details.resource_constraints` | `cpu_constrained` / `memory_constrained` flags from the agent's resource monitor. |
| `status_details.failure_metadata_tags` | Free-form string tags attached to the failure by the agent command (e.g. team tags like `assigned_to:my-team`). Empty list means no tags were set. |
| `logs` / `parsley_logs` | URLs to fetch task log content. |
| `start_time` / `finish_time` | Task duration — useful for spotting log gaps. |
| `revision` | The git commit SHA this task ran against. |
| `project_identifier` | The Evergreen project ID (needed for version history queries). |

### 3.5. Fetch task owner from foliage

Query the Evergreen GraphQL endpoint to get the foliage-assigned team owner:

```bash
python3 - <<'EOF'
import subprocess, json

token = subprocess.run(["evergreen","client","get-oauth-token"],capture_output=True,text=True).stdout.strip()
task_id = "<TASK_ID>"
execution = 1  # execution number from the URL

query = f"""
query {{
  task(taskId: "{task_id}", execution: {execution}) {{
    taskOwnerTeam {{ teamName jiraProject assignmentType messages }}
  }}
}}
"""
r = subprocess.run(["curl","-s","-X","POST",
    "-H", f"Authorization: Bearer {token}",
    "-H","Content-Type: application/json",
    "-d", json.dumps({"query": query}),
    "https://evergreen.corp.mongodb.com/graphql/query"],
    capture_output=True, text=True)
print(json.dumps(json.loads(r.stdout).get("data",{}).get("task",{}).get("taskOwnerTeam") or {}, indent=2))
EOF
```

If `assignmentType` is `none` or `teamName` is empty, foliage has no mapping for this task — fall back to `failure_metadata_tags` for team signal. Always include both in the verdict output.

### 3.6. Fetch recent distro AMI changes

Using the `distro_id` from Phase 3 and the task's `start_time`, check whether the distro's AMI was recently updated. A recent AMI change is one of the most common causes of sudden regressions that have no corresponding change in the project's own code.

This uses the Evergreen GraphQL API — accessible to any user with `DistroSettingsView` permission on the distro. No admin or MongoDB access required.

```bash
python3 - <<'EOF'
import subprocess, json
from datetime import datetime, timezone

token = subprocess.run(
    ["evergreen", "client", "get-oauth-token"],
    capture_output=True, text=True
).stdout.strip()
distro_id = "<DISTRO_ID>"         # from Phase 3: task["distro_id"]
task_start  = "<START_TIME>"      # from Phase 3: task["start_time"] (ISO string)

query = """
query DistroEvents($opts: DistroEventsInput!) {
  distroEvents(opts: $opts) {
    count
    eventLogEntries {
      before
      after
      timestamp
      user
    }
  }
}
"""

payload = json.dumps({"query": query, "variables": {"opts": {"distroId": distro_id, "limit": 20}}})
r = subprocess.run([
    "curl", "-s", "-X", "POST",
    "-H", f"Authorization: Bearer {token}",
    "-H", "Content-Type: application/json",
    "-d", payload,
    "https://evergreen.corp.mongodb.com/graphql/query"
], capture_output=True, text=True)

resp = json.loads(r.stdout)

# Surface permission errors explicitly — the directive returns a GraphQL error
# (not just an empty result) when the user lacks DistroSettingsView.
gql_errors = resp.get("errors") or []
if gql_errors:
    msgs = [e.get("message", "") for e in gql_errors]
    if any("does not have permission" in m or "Forbidden" in m for m in msgs):
        print(f"SKIP: insufficient permissions to read distro events for '{distro_id}'.")
        print("      (requires DistroSettingsView on this distro — ask a DevProd team member to run this check)")
    else:
        print(f"GraphQL error fetching distro events: {msgs}")
    raise SystemExit(0)

entries = (resp.get("data") or {}).get("distroEvents", {}).get("eventLogEntries") or []

def get_ami(distro_doc):
    psl = (distro_doc or {}).get("provider_settings_list") or (distro_doc or {}).get("provider_settings") or []
    return psl[0].get("ami") if psl and isinstance(psl[0], dict) else None

task_dt = datetime.fromisoformat(task_start.rstrip("Z")).replace(tzinfo=timezone.utc)
ami_changes = []
for e in entries:
    before_ami = get_ami(e.get("before"))
    after_ami  = get_ami(e.get("after"))
    if before_ami and after_ami and before_ami != after_ami:
        ami_changes.append({"timestamp": e["timestamp"], "user": e.get("user"), "before": before_ami, "after": after_ami})

if ami_changes:
    print(f"AMI changes for distro '{distro_id}' (newest first):")
    for ch in ami_changes:
        change_dt = datetime.fromisoformat(ch["timestamp"].rstrip("Z")).replace(tzinfo=timezone.utc)
        days_before = (task_dt - change_dt).days
        marker = f"  ⚠️  {days_before}d before task" if 0 <= days_before <= 14 else (
                  "  (after task)" if days_before < 0 else f"  ({days_before}d before task)")
        print(f"  {ch['timestamp']}  {ch['before']} → {ch['after']}  by {ch['user']}{marker}")
else:
    print(f"No AMI changes found in the last 20 events for distro '{distro_id}'.")
EOF
```

**How to interpret the output:**

| Scenario | Implication |
|---|---|
| AMI changed 0–14 days before `start_time` | Strong signal — carry this into Phase 8 as Pattern E |
| AMI changed >14 days before `start_time` | Weak signal — note it but don't lead with it |
| AMI changed after `start_time` | Not a factor for this failure |
| No AMI changes in last 20 events | Rule out this hypothesis |
| `SKIP: insufficient permissions` | Note it in the verdict as "AMI check skipped (no DistroSettingsView)" — do not treat absence of data as evidence against Pattern E |

Carry the most recent AMI change timestamp and IDs into the verdict regardless of whether it's within the 14-day window.

### 4. Classify by `status_details.type`

**`"test"` — agent was alive, command exited non-zero**

The Evergreen agent ran your task to completion. A command you configured exited with a non-zero code and the agent reported it. The Evergreen infrastructure functioned correctly. **This is almost always your code, your script, or your test suite.** Proceed to Phase 5 (trace the failing command).

**`"system"` — agent died or host was lost**

The Evergreen app server marked the task failed because it stopped receiving communication from the agent. This *can* be infrastructure, but is frequently caused by:
- Your task applying iptables/network rules that also block the agent's heartbeat
- Your task consuming all available memory (OOM kill — check `status_details.oom_tracker_info.detected`)
- Your task consuming all network bandwidth (preventing agent communication)

Check `status_details.resource_constraints` and look for network isolation in your scripts before assuming infrastructure failure. Proceed to Phase 8 (check known patterns) before rendering a verdict.

**`"setup"` — pre-task command failed**

The task failed before your main commands ran. Common causes: git clone failure, dependency installation failure, bad project YAML. Check your `pre` block and setup functions.

**empty / `""`  — task is still running or was aborted**

Wait for the task to finish, or check if it was manually aborted.

### 5. Trace the failing command in the project YAML

From `details.failing_command`, extract the function name (e.g. `'subprocess.exec' in function 'run my test suite'` → function name is `run my test suite`).

**Find the project YAML.** Check common locations:
- `.evergreen.yml` in the repo root
- `evergreen.yml`
- Files listed under `include:` in the main config

```bash
# Find the function definition
grep -A 20 '"<function name>"' .evergreen.yml
# or
grep -rA 20 '"<function name>"' evergreen/
```

**Read what the function does.** If it runs a shell script:
```bash
cat <path/to/script.sh>
```

Read the script. Look for: `exit 1` conditions, error traps, commands that could fail under the circumstances described in `details.desc`.

### 6. Fetch and read the task logs

Use the same safe Python wrapper for all log fetching to keep the token out of shell history and output:

```bash
python3 - <<'EOF'
import subprocess, json

token = subprocess.run(
    ["evergreen", "client", "get-oauth-token"],
    capture_output=True, text=True
).stdout.strip()
headers = ["-H", f"Authorization: Bearer {token}"]
task_id = "<TASK_ID>"

# Get log URL from task metadata
r = subprocess.run(["curl", "-s"] + headers +
    [f"https://evergreen.corp.mongodb.com/rest/v2/tasks/{task_id}"],
    capture_output=True, text=True)
task = json.loads(r.stdout)
log_url = task["logs"]["task_log"] + "&text=true"

# Fetch log content
r = subprocess.run(["curl", "-sL"] + headers + [log_url],
    capture_output=True, text=True)
lines = r.stdout.splitlines()
print("=== FIRST 20 LINES ===")
print("\n".join(lines[:20]))
print("\n=== LAST 30 LINES ===")
print("\n".join(lines[-30:]))
EOF
```

Note the **first** and **last** timestamps. If the last log timestamp is significantly earlier than `finish_time`, the logs were cut off mid-task (see known patterns below).

### 7. Find the last passing run and diff commits

This phase answers: **did something in your project change that could explain the regression?**

#### 7a. Get the failing task's revision and metadata

From the task metadata fetched in Phase 3, extract:
- `revision` — the git commit SHA this task ran against
- `project` — the Evergreen project identifier
- `build_variant` — the build variant name
- `display_name` — the task display name

#### 7b. Find recent runs of this task via the Evergreen REST API

Use the safe Python wrapper — never pass the token as a literal string in a shell command:

```bash
python3 - <<'EOF'
import subprocess, json

token = subprocess.run(
    ["evergreen", "client", "get-oauth-token"],
    capture_output=True, text=True
).stdout.strip()
headers = ["-H", f"Authorization: Bearer {token}"]
project = "<PROJECT_ID>"

r = subprocess.run(["curl", "-s"] + headers +
    [f"https://evergreen.corp.mongodb.com/rest/v2/projects/{project}/versions"
     "?requester=gitter_request&limit=30"],
    capture_output=True, text=True)
versions = json.loads(r.stdout)
for v in versions:
    print(v.get("revision","")[:12], v.get("create_time","")[:19], v.get("message","")[:60])
EOF
```

This gives you the 30 most recent mainline versions with their git SHAs and commit messages. Find the window around when the task started failing.

#### 7c. Check task status for candidate versions

```bash
python3 - <<'EOF'
import subprocess, json

token = subprocess.run(
    ["evergreen", "client", "get-oauth-token"],
    capture_output=True, text=True
).stdout.strip()
headers = ["-H", f"Authorization: Bearer {token}"]
version_id = "<VERSION_ID>"
build_variant = "<BUILD_VARIANT>"
display_name_substr = "<TASK_DISPLAY_NAME_SUBSTRING>"

r = subprocess.run(["curl", "-s"] + headers +
    [f"https://evergreen.corp.mongodb.com/rest/v2/versions/{version_id}/tasks"
     f"?variant={build_variant}"],
    capture_output=True, text=True)
tasks = json.loads(r.stdout)
for t in tasks:
    if display_name_substr in t.get("display_name", ""):
        print(t.get("status"), t.get("display_name"), t.get("task_id","")[:40])
EOF
```

Walk backwards through versions until you find one where the task had `status: success`. That version's `revision` is the last known good commit.

**Shortcut:** if the failing task ID contains a git SHA (e.g. `myproject_variant_taskname_<SHA>_YY_MM_DD`), that SHA is the failing revision. The task just before it in version history is a candidate for the last passing revision.

#### 7d. Diff commits between last passing and first failing

Once you have:
- `LAST_GOOD_SHA` — revision of the last passing run
- `FIRST_BAD_SHA` — revision of the failing task

Run in the user's project repo:

```bash
# All commits in the regression window
git log --oneline <LAST_GOOD_SHA>..<FIRST_BAD_SHA>

# Commits that touched files related to the failing task
# (substitute the paths your task's scripts live in)
git log --oneline <LAST_GOOD_SHA>..<FIRST_BAD_SHA> -- \
  .evergreen.yml \
  evergreen/ \
  scripts/evergreen/ \
  <any other relevant directories>

# Show what actually changed in those commits
git log -p <LAST_GOOD_SHA>..<FIRST_BAD_SHA> -- \
  .evergreen.yml evergreen/ scripts/evergreen/
```

**What to look for:**
- Changes to the failing function in the project YAML
- Changes to scripts that function calls
- Removed or renamed test fixtures, binaries, or configuration files
- New dependencies or version pins that the test environment doesn't satisfy
- Changes to Docker configuration, network setup, or binary download steps

#### 7e. Correlate commits with the failure

For each commit in the window, ask:
1. Does it touch the failing function or its scripts? → High suspicion.
2. Does it remove something the test previously depended on? → High suspicion.
3. Does the commit message reference removing, renaming, or upgrading something? → Read the diff carefully.
4. Is it a mechanical change (formatting, comments, unrelated feature)? → Lower suspicion.

If a specific commit stands out as the likely cause, include it in the verdict with the commit SHA and a one-line explanation of why.

If the repo isn't available locally, use `gh`:
```bash
gh api repos/<org>/<repo>/commits?sha=<FIRST_BAD_SHA>&per_page=20 \
  --jq '.[] | .sha[:12] + " " + .commit.message[:70]'
```

---

### 8. Check for known misleading patterns

Before rendering a verdict, check whether the symptoms match a known pattern. These patterns look like infrastructure failures but are caused by the project's own task setup.

---

#### Pattern A: Log cutoff + post-task network failures

**Symptoms:**
- Logs stop partway through the task (last log timestamp much earlier than `finish_time`)
- `post_errored: true`
- `other_failing_commands` includes `s3.put`, `attach.xunit_results`, or similar network-dependent post-task commands
- The main task still ran to completion (it has a `details.failing_command` or shows `status: failed` from the test itself)

**Cause:** Your task applies network isolation rules (iptables, nftables, Docker network policies, or similar) that block outbound connections. These rules also block the Evergreen agent's log streaming to S3 and the post-task commands that need network access.

**How to verify:** Search your task's scripts for:
```bash
grep -r "iptables\|nftables\|restrict.*internet\|block.*network\|docker network" .
```

**Verdict implication:** The missing logs and post-task failures are a side effect of your own network isolation setup. They are not an Evergreen infrastructure problem. The actual task failure (if any) is separate.

---

#### Pattern B: `"system"` type but task was resource-intensive

**Symptoms:**
- `details.type: "system"` with `details.desc` containing "heartbeat" or "system unresponsive"
- Your task runs CPU/memory/network-intensive workloads (large Docker workloads, compilation, data processing)

**Cause:** Your task consumed enough resources that the host became unresponsive, and the Evergreen app server marked it failed after the heartbeat timed out. This looks like infrastructure failure but is workload-induced.

**How to verify:** Look at the task's workload characteristics. Does it run Docker containers, compile large codebases, or saturate network I/O? If yes, this is likely resource exhaustion, not infrastructure.

**Verdict implication:** You may need a larger distro, resource limits in your task configuration, or to reduce your task's resource consumption.

---

#### Pattern C: "Bazel test(s) failed" / "process encountered problem: exit code N"

**Symptoms:**
- `details.desc` says "Bazel test(s) failed" or "process encountered problem: exit code N"
- `details.type: "test"`

**Cause:** These are agent-generated messages meaning your subprocess exited non-zero. The agent is correctly reporting what your command returned.

**Verdict implication:** Your test suite or script exited non-zero. Look at your actual test output in the logs for the root cause. This is never an Evergreen infrastructure issue.

---

#### Pattern D: Post-task failures only, main task succeeded

**Symptoms:**
- `status: success` but `post_errored: true`
- Failures only in `other_failing_commands` (all in the `post` block)

**Cause:** Your main task commands passed, but post-task cleanup commands failed. This is almost always either (a) network isolation still active during post-task, or (b) an artifact that was supposed to be generated wasn't found.

**Verdict implication:** Post-task failures alone rarely indicate an Evergreen infrastructure problem.

---

#### Pattern E: Recent distro AMI change coincides with regression onset

**Symptoms:**
- Task was passing consistently, then started failing without any project code changes in the regression window (Phase 7 found nothing relevant)
- Phase 3.6 found an AMI change on the task's distro within ~14 days before the task's `start_time`
- Failure is `"system"` type, OR test failures that look environmental (missing binary, changed path, unexpected OS behavior, different library versions)

**Cause:** The distro's base AMI was rotated. The new AMI may have a different OS version, different pre-installed packages, different kernel, or different default configuration compared to the previous AMI. Tasks that depended implicitly on a specific environment can break when the AMI changes under them.

**How to verify:** Compare the before and after AMI IDs from Phase 3.6 against the regression window from Phase 7. If the AMI changed during the window where the task went from passing to failing, and no project commit explains the failure, this pattern applies.

**Verdict implication:** This is a genuine environmental change, but it originates from Evergreen's distro management (AMI rotation is expected and automated). The fix is typically in the project — either pinning a dependency, using a more portable setup, or updating the task to work with the new AMI. File a DEVPROD ticket only if the new AMI broke something that was explicitly supported (e.g., a required binary is missing, or the AMI change was undocumented).

---

### 9. Render a verdict

Produce one of the following:

---

**VERDICT: YOUR CODE OR CONFIGURATION**

Owner: [teamName from foliage, or "unassigned" if none] | Jira: [jiraProject] | Assignment: [assignmentType]
Failure metadata tags: [list from status_details.failure_metadata_tags, or "none"]
Recent AMI change: [AMI before → AMI after, date — or "none in the last 20 events" — or "check skipped: no DistroSettingsView permission"]

Evidence: [cite details.type, details.failing_command, what the script does, log content]

This failure is in your project's code or task configuration. Evergreen infrastructure ran correctly and faithfully reported what your command returned.

Recommendation: [specific file or function to investigate, with reason]

---

**VERDICT: YOUR TASK SETUP (network isolation or resource exhaustion)**

Owner: [teamName from foliage, or "unassigned" if none] | Jira: [jiraProject] | Assignment: [assignmentType]
Failure metadata tags: [list from status_details.failure_metadata_tags, or "none"]
Recent AMI change: [AMI before → AMI after, date — or "none in the last 20 events" — or "check skipped: no DistroSettingsView permission"]

Evidence: [cite the pattern — log cutoff + post-task failures, or system type + intensive workload]

The symptoms look like an infrastructure problem but are caused by your task's own behavior. [Explain the specific mechanism — iptables rules, memory consumption, etc.]

Recommendation: [specific script or configuration to check; e.g. "search your scripts for iptables rules that may also be blocking the agent's log streaming"]

---

**VERDICT: POSSIBLE EVERGREEN INFRASTRUCTURE ISSUE**

Owner: [teamName from foliage, or "unassigned" if none] | Jira: [jiraProject] | Assignment: [assignmentType]
Failure metadata tags: [list from status_details.failure_metadata_tags, or "none"]
Recent AMI change: [AMI before → AMI after, date — or "none in the last 20 events" — or "check skipped: no DistroSettingsView permission"]

Evidence: [details.type: "system" with no obvious resource or network cause; unexpected host termination; no correlation with your task's behavior]

This looks like a genuine infrastructure issue. Before filing a ticket, collect:
- Task ID and execution number
- `start_time` and `finish_time` from the task metadata
- `details.desc` verbatim
- Whether the failure is reproducible (does it happen on every run, or intermittently?)

File a ticket in the DEVPROD Jira project with this information.

---

**VERDICT: DISTRO AMI CHANGE**

Owner: [teamName from foliage, or "unassigned" if none] | Jira: [jiraProject] | Assignment: [assignmentType]
Failure metadata tags: [list from status_details.failure_metadata_tags, or "none"]
AMI change: [before AMI] → [after AMI] on [date], [N] day(s) before this task ran

Evidence: [no project code changes in the regression window; AMI changed N days before the task started failing; describe how the failure symptoms look environmental]

The distro's base AMI was rotated before this task ran. The new AMI may differ in OS packages, kernel version, or default configuration. The task likely depended implicitly on something that changed in the new image.

Recommendation: Compare what changed between the two AMIs. Check whether required binaries, libraries, or paths are present on the new AMI.

---

**VERDICT: AMBIGUOUS — NEED MORE INFORMATION**

Owner: [teamName from foliage, or "unassigned" if none] | Jira: [jiraProject] | Assignment: [assignmentType]
Failure metadata tags: [list from status_details.failure_metadata_tags, or "none"]
Recent AMI change: [AMI before → AMI after, date — or "none in the last 20 events" — or "check skipped: no DistroSettingsView permission"]

Evidence: [what is inconclusive and why]

Before this can be classified, collect: [specific data — e.g. "check whether your task applies any iptables rules", "confirm the task's memory usage"]

---

## Quick reference: what the agent failure types mean

The REST API field is `status_details.type`. The `details` key (without `status_`) is the internal MongoDB field name and will be empty in API responses.

| `status_details.type` | Who decided the task failed | Most likely cause |
|---|---|---|
| `test` | The agent (your command exited non-zero) | Your code or scripts |
| `system` | The Evergreen app server (agent stopped responding) | Resource exhaustion in your task, or network isolation your task applied, or (less commonly) infrastructure |
| `setup` | The agent (pre-task command failed) | Bad project YAML, network issue fetching dependencies |
| empty | Not yet determined | Task still running or was aborted |

`status_details.type` is the single most important piece of information. A `"test"` failure is never an Evergreen infrastructure issue — it is your code.
