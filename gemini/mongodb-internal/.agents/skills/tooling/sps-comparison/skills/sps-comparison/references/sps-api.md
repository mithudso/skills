# SPS (Signal Processing Service) API Reference

## Base URLs and Auth

| Service | Base URL | Auth |
|---|---|---|
| Comparisons | `https://performance-monitoring-api.corp.mongodb.com/comparisons` | `Api-User` / `Api-Key` |
| Version comparisons | `https://performance-monitoring-api.corp.mongodb.com/version_comparisons` | `Api-User` / `Api-Key` |
| Variant comparisons | `https://performance-monitoring-api.corp.mongodb.com/variant_comparisons` | `Api-User` / `Api-Key` |
| Managed multipatch (version) | `https://performance-monitoring-api.corp.mongodb.com/version_comparisons/managed_analysis` | `Api-User` / `Api-Key` |
| Managed multipatch (variant) | `https://performance-monitoring-api.corp.mongodb.com/variant_comparisons/managed_analysis` | `Api-User` / `Api-Key` |
| Baseline/discovery | `https://performance-monitoring-api.corp.mongodb.com/time_series` | `Api-User` / `Api-Key` |
| Raw perf results | `https://performance-monitoring-api.corp.mongodb.com/raw_perf_results` | No auth |

Credentials come from `~/.evergreen.yml` (`user` and `api_key` fields).

Requires Cloudflare WARP VPN for `corp.mongodb.com` access.

## Comparison Endpoints

### Get metadata (all types)

```
GET /comparisons/id/<comparison_id>
```

Response fields: `comparison_type` (`version_comparison`, `variant_comparison`, `perf_analyzer`), `comparison.status`, `comparison.compare_versions`, `comparison.base_versions`, `comparison.metrics`, `comparison.title`, `comparison.created_by`.

### Get status

```
GET /version_comparisons/<comparison_id>/status
```

### Get available filter fields

```
GET /comparisons/<comparison_id>/fields
```

Returns variants, tasks, tests, measurements available in the comparison.

### List user comparisons

```
GET /comparisons/user/<evg_user>
```

### Fetch results by type

| Type | Endpoint |
|---|---|
| `version_comparison` | `GET /version_comparisons/<id>/results?skip=0&limit=-1` |
| `variant_comparison` | `GET /variant_comparisons/<id>/results?skip=0&limit=-1` |
| `perf_analyzer` | `GET /performance/analysis/<id>/details` |

## Result Row Fields

| Field | Description |
|---|---|
| `task` | Evergreen task name |
| `test` | Test/phase name (e.g., `Find-tuned_run`) |
| `measurement` | Metric name (e.g., `OperationThroughput`, `AverageLatency`) |
| `improvement_direction` | `up` (higher=better) or `down` (lower=better) |
| `compare_value_map` | Per-version raw values from patch runs |
| `compare_mean` | Mean across all patch runs |
| `compare_cov` | Coefficient of variation across patch runs |
| `compare_stable_mean` | Mean of the stable region (recent mainline) |
| `compare_mean_percent_change_against_stable_region` | % change vs stable baseline |
| `z_score` | Statistical significance (abs > 2 = significant) |
| `base_mean` | Mean of the base version values (if any) |

## Creating Comparisons

### Version comparison

```
POST /version_comparisons/
{
  "compare_version_id": "<PATCH_VERSION_ID>",
  "base_version_id": "<MAINLINE_VERSION_ID>",
  "title": "<title>",
  "analyze_trend": false
}
```

### Managed multipatch (clone + compare)

Use `/version_comparisons/managed_analysis` for version-level or `/variant_comparisons/managed_analysis` for variant-level:

```
POST /version_comparisons/managed_analysis
{
  "version": "<PATCH_VERSION_ID>",
  "clone_count": 3,
  "base_version": "<MAINLINE_VERSION_ID>",
  "metrics": ["ops_per_sec", "Latency95thPercentile"],
  "notification_preference": "none",
  "notification_target": ""
}
```

**Safety limit:** `tasks × clone_count` must not exceed **600**. Exceeding this returns HTTP 422.

**Auth:** Requests must carry `Api-User` / `Api-Key` tied to a human user. Anonymous and `spiffe://` identities are blocked.

**Legacy endpoint `POST /performance/analysis/create` is disabled — do not use.**

## Baseline Endpoints

### Discover test names

```
GET /time_series/?project_regex=sys-perf&variant_regex=<variant>&task_regex=<task>&measurement_regex=<measurement>&limit=50
```

Extract distinct test names from `data[].test`.

### Fetch stable region baseline

```
POST /time_series/stable_region
{
  "project": "sys-perf",
  "variant": "<variant>",
  "task": "<task>",
  "test": "<test>",
  "measurement": "<measurement>",
  "args": {},
  "last": true
}
```

`args` must use numeric values where applicable (e.g., `{"thread_level": 128}` not `{"thread_level": "128"}`).

## UI Links

- **Performance Analyzer**: `https://performance-analyzer.server-tig.prod.corp.mongodb.com/perf-analyzer-viz/?comparison_id=<ID>`
- **Comparison Creator**: `https://performance-monitoring-and-analysis.server-tig.prod.corp.mongodb.com/perf-comparison-creator`
- **Swagger Docs**: `https://performance-monitoring-service-rest.server-tig.prod.corp.mongodb.com/docs`
