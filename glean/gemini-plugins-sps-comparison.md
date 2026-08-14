# sps-comparison

**Category:** Frontend & Web Development
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/tooling/sps-comparison/skills/sps-comparison

## Description
Use when fetching SPS Performance Analyzer comparison results, checking multipatch comparison status, listing user comparisons, or creating new multipatches. Triggers on Performance Analyzer URLs or comparison IDs.

---

# SPS Performance Analyzer Comparison

Fetch, list, and create multipatch comparison results from the Signal Processing Service (SPS) API -- the backend powering the Performance Analyzer UI.

## Prerequisites

- Evergreen credentials in `~/.evergreen.yml` (`user` and `api_key` fields)
- Cloudflare WARP VPN active (required for `corp.mongodb.com` endpoints)

## Credential Setup

```bash
EVG_USER=$(grep '^user:' ~/.evergreen.yml | awk '{print $2}')
EVG_KEY=$(grep '^api_key:' ~/.evergreen.yml | awk '{print $2}')
```

Headers: `-H "Api-User: ${EVG_USER}" -H "Api-Key: ${EVG_KEY}"`

## Routing

| User intent | Action |
|---|---|
| Pastes a Performance Analyzer URL or comparison ID | Fetch results (Step 1) |
| "List my comparisons", "Show comparisons" | List comparisons (Step 2) |
| "Create a comparison for patch X", "Compare my patch" | Create comparison (Step 3) |
| "Create a multipatch", "Clone and compare patch X" | Create comparison (Step 3) |

Extract `comparison_id` from URLs like:
`https://performance-analyzer.server-tig.prod.corp.mongodb.com/perf-analyzer-viz/?comparison_id=<ID>`

## Step 1 -- Fetch Comparison Results

### 1a. Get metadata

```
GET https://performance-monitoring-api.corp.mongodb.com/comparisons/id/<comparison_id>
```

Extract: `comparison_type`, `comparison.status`, `comparison.title`, `comparison.created_by`, `comparison.compare_versions`, `comparison.base_versions`, `comparison.metrics`.

If `status` is not `success`, report current status and stop. Statuses: `submitted`, `waiting_on_evergreen`, `performing_comparison`, `failed`, `success`.

### 1b. Fetch results

Endpoint depends on `comparison_type`:

| Type | Endpoint |
|---|---|
| `version_comparison` | `GET /version_comparisons/<id>/results?skip=0&limit=-1` |
| `variant_comparison` | `GET /variant_comparisons/<id>/results?skip=0&limit=-1` |
| `perf_analyzer` | `GET /performance/analysis/<id>/details` |

Base URL: `https://performance-monitoring-api.corp.mongodb.com`

Use `limit=-1` to get all results.

### 1c. Display results

Report metadata table, then results sorted by `|z_score|` descending:

```markdown
## Comparison: <title>

| Field | Value |
|---|---|
| Comparison ID | `<id>` |
| Type | <comparison_type> |
| Status | <status> |
| Created by | <user> |
| Compare versions | <count> |
| Base version | <version_id> |
| [Performance Analyzer](<url>) | |

### Results

| Task | Test | Metric | Mean | CoV | % Change vs Stable | Z-Score | Significant? |
|---|---|---|---:|---:|---:|---:|---|
```

Flag rows where `abs(z_score) > 2` as significant. Respect `improvement_direction` -- for `down` metrics (latency), negative % change is an improvement.

See [references/sps-api.md](references/sps-api.md) for full result row field descriptions and additional endpoints.

## Step 2 -- List Comparisons

```
GET https://performance-monitoring-api.corp.mongodb.com/comparisons/user/<evg_user>
```

Display as table: comparison ID, type, title, status, created date.

## Step 3 -- Create Comparison from Evergreen Patch

When the user provides an Evergreen patch ID, Spruce URL, or asks to compare a patch.

### 3a. Resolve the version ID

For CLI patches, the patch ID is the version ID. Extract from URLs:
- `https://spruce.mongodb.com/version/<version_id>`
- `https://spruce.mongodb.com/patch/<patch_id>`
- `https://evergreen.mongodb.com/version/<version_id>`

Verify the patch exists:

```
GET https://evergreen.mongodb.com/rest/v2/versions/<version_id>
```

Use `Api-User` / `Api-Key` headers. Extract the `project` field for finding a base version.

### 3b. Find a base version (if user doesn't provide one)

Get a recent mainline version for the same project:

```
GET https://evergreen.mongodb.com/rest/v2/projects/<project>/versions?requester=mainline_commit&limit=5
```

Use the `id` field from the first result as `base_version_id`.

### 3c. Create the comparison

Ask the user which type they want:

**Simple version comparison** (fast, single run vs baseline):

```
POST https://performance-monitoring-api.corp.mongodb.com/version_comparisons/
{
  "compare_version_id": "<PATCH_VERSION_ID>",
  "base_version_id": "<MAINLINE_VERSION_ID>",
  "title": "<description>",
  "analyze_trend": false
}
```

**Managed multipatch** (clones the patch N times for statistical power, takes longer):

```
POST https://performance-monitoring-api.corp.mongodb.com/version_comparisons/managed_analysis
{
  "version": "<PATCH_VERSION_ID>",
  "clone_count": 3,
  "base_version": "<MAINLINE_VERSION_ID>",
  "metrics": ["ops_per_sec", "Latency95thPercentile"],
  "notification_preference": "none",
  "notification_target": ""
}
```

Note: the managed multipatch body uses `version`/`base_version` fields — not `compare_version_id`/`base_version_id` (which are for the simple version comparison endpoint above). For variant-level comparisons, use `/variant_comparisons/managed_analysis` with the same request body.

> **Safety limit:** The total number of tasks × `clone_count` must not exceed **600**. Before submitting, estimate the task count from the patch (via the Evergreen API) and validate `tasks × clone_count ≤ 600`. If the user requests a higher clone count, warn them and cap or reduce it. Requests exceeding this limit return **HTTP 422**.

> **Auth requirement:** Requests must be tied to a human user via `Api-User` / `Api-Key` headers from `~/.evergreen.yml`. Anonymous or service-identity-only requests (e.g., `spiffe://`) are blocked from creating comparisons.

> **Legacy endpoint `POST /performance/analysis/create` is disabled.** Do not use it.

### 3d. Handle errors

If the API returns **422**, the request exceeded the 600-task safety cap. Report this to the user with the calculated `tasks × clone_count` value and suggest reducing `clone_count` or narrowing the task scope.

### 3e. Report

Print the comparison ID and link:
`https://performance-analyzer.server-tig.prod.corp.mongodb.com/perf-analyzer-viz/?comparison_id=<ID>`

The comparison takes time to process. Suggest checking back with the comparison ID (Step 1).

## Quick One-Liner

Dump results to terminal (adjust endpoint per `comparison_type`):

```bash
COMPARISON_ID="<id>"
EVG_USER=$(grep '^user:' ~/.evergreen.yml | awk '{print $2}')
EVG_KEY=$(grep '^api_key:' ~/.evergreen.yml | awk '{print $2}')

curl -s "https://performance-monitoring-api.corp.mongodb.com/version_comparisons/${COMPARISON_ID}/results?skip=0&limit=-1" \
  -H "Api-User: ${EVG_USER}" -H "Api-Key: ${EVG_KEY}" \
  | python3 -c "
import json, sys, math
data = json.load(sys.stdin)
results = data.get('results', [])
print(f'Total rows: {len(results)}\n')
print(f\"{'Task':<55} {'Test':<40} {'Metric':<30} {'Mean':>12} {'CoV':>8} {'%ChgStable':>12} {'Z':>8}\")
print('-' * 167)
for r in sorted(results, key=lambda x: abs(float(x.get('z_score',0) or 0)), reverse=True):
    def f(v, d=2):
        try:
            v=float(v); return 'N/A' if math.isnan(v) else f'{v:.{d}f}'
        except: return 'N/A'
    cm=r.get('compare_mean',float('nan'))
    cc=r.get('compare_cov',float('nan'))
    pct=r.get('compare_mean_percent_change_against_stable_region',float('nan'))
    zs=r.get('z_score',float('nan'))
    sig=' *' if abs(float(zs or 0))>2 else ''
    print(f\"{r['task'][:55]:<55} {r['test'][:40]:<40} {r['measurement'][:30]:<30} {f(cm):>12} {f(cc,3):>8} {f(pct):>11}% {f(zs):>7}{sig}\")
"
```