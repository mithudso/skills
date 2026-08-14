---
name: check-epic-fixversions
description: Audit all Jira tickets in an Epic to verify fixVersion fields are correctly set. Checks resolution, merged GitHub PRs, and branch names. Strictly checks WT and SERVER tickets; skips other project types. Use when asked to audit or review Epic tickets for fixVersion compliance.
argument-hint: <epic-ticket-id>
arguments: epic
allowed-tools: mcp__jira__jira_search mcp__jira__jira_get_issue mcp__jira__jira_get_project_versions mcp__jira__jira_update_issue mcp__jira__jira_transition_issue mcp__jira__jira_get_transitions Bash(gh *)
source: 10gen/agent-skills
license: Internal
mongodb:
  team: storage-engines
  owner: luke.chen@mongodb.com
  internal: true
---

# Check Epic fixVersions: $epic

Audit every ticket in Epic **$epic** to verify `fixVersion` is correctly set.

## Step 1 — Fetch all issues in the Epic

Use `mcp__jira__jira_search` with JQL:
```
"Epic Link" = $epic ORDER BY key ASC
```
Request fields: `key,summary,status,resolution,fixVersions,issuetype,project`.

If that returns 0 results, retry with:
```
parent = $epic ORDER BY key ASC
```

Collect all issues. Note the total count.

## Step 2 — Categorize each ticket

For each issue, determine which check applies:

**Skip (no fixVersion required)** if any of:
- Resolution is: Won't Fix, Duplicate, Cannot Reproduce, Incomplete, Not a Bug, Works as Designed, Invalid, Obsolete
- Status is open (resolution not set)
- Project key is not `WT` or `SERVER` — other projects (HELP, BUILD, DEVPROD, etc.) may not have fixVersion; do not flag them

**Requires fixVersion** if:
- Resolution = "Fixed" AND project key is `WT` or `SERVER`

**Needs GitHub investigation** if:
- Resolution = "Done" AND project key is `WT` or `SERVER`

## Step 3 — GitHub check for "Done" WT/SERVER tickets

For each "Done" WT or SERVER ticket, search for merged PRs that reference the ticket key:

**WT tickets** — search `wiredtiger/wiredtiger`:
```bash
gh pr list --repo wiredtiger/wiredtiger --search "WT-XXXX in:title,body" --state merged --json number,title,baseRefName,url,mergedAt --limit 50
```

**SERVER tickets** — search `10gen/mongo`:
```bash
gh pr list --repo 10gen/mongo --search "SERVER-XXXX in:title,body" --state merged --json number,title,baseRefName,url,mergedAt --limit 50
```

Replace `WT-XXXX` / `SERVER-XXXX` with the actual ticket key.

**Interpret results:**

| PRs found? | Any PR merged into `develop` / `master` / `main`? | Verdict |
|------------|---------------------------------------------------|---------|
| No | — | ✅ OK — no code change, no fixVersion required |
| Yes | No (feature branches only) | ✅ OK (feature branch only) — no `WT<x>.<y>.<z>` fixVersion needed |
| Yes | Yes | ❌ Resolution should be "Fixed" (not "Done") AND fixVersion must be set |

For WT: `develop` is the release integration branch. A `WT<x>.<y>.<z>` fixVersion is only appropriate when a PR was merged to `develop`. Merging to a feature branch alone is not a release.

## Step 4 — Produce the findings report

Print a markdown table with one row per ticket:

| Ticket | Summary | Resolution | fixVersion | Status |
|--------|---------|------------|------------|--------|

Status values:
- ✅ **OK** — passes all checks
- ⚠️ **Wrong resolution** — "Done" but has merged code in develop; should be "Fixed"
- ❌ **Missing fixVersion** — "Fixed" but no fixVersion set
- ❌ **Missing fixVersion + wrong resolution** — "Done" with code in develop, fixVersion absent
- ℹ️ **Skipped** — non-WT/SERVER project, or resolution doesn't require fixVersion

Follow the table with summary counts:
```
Total tickets in epic:             N
✅ OK:                             N
❌ Missing fixVersion:             N
⚠️  Wrong resolution (Done→Fixed): N
ℹ️  Skipped:                       N
```

## Step 5 — Offer to fix (explicit confirmation required)

If any issues were found, present a numbered list of proposed changes, for example:

```
Proposed fixes:
1. WT-XXXX — change resolution "Done" → "Fixed"
2. WT-YYYY — set fixVersion (please specify the version to use)
3. WT-ZZZZ — change resolution "Done" → "Fixed" AND set fixVersion

Apply all? Or list the numbers you want applied. (Type "none" to skip.)
```

**Do NOT make any Jira changes until the user explicitly confirms.**

When fixing:
- **Change resolution**: call `mcp__jira__jira_get_transitions` to find the "Fixed" transition ID, then `mcp__jira__jira_transition_issue`.
- **Set fixVersion**: call `mcp__jira__jira_update_issue` with `{"fixVersions": [{"name": "<version>"}]}`. If the user has not specified a version, ask: "What fixVersion should I use? (e.g., WT7.0.0)"

After applying all changes, print a confirmation list of every update made.
