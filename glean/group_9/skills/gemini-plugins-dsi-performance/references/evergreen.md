# Evergreen Integration for DSI Performance Tests

## Submitting a Performance Patch

From 10gen/dsi checkout:
```bash
evergreen patch -p dsi -d "perf test: <description>" -y
```

## Monitoring Results

- View results: https://evergreen.mongodb.com/waterfall/dsi
- Performance dashboard: https://performance.mongodb.com
- Perf regression alerts: check #performance-alerts Slack channel

## Common Evergreen Tasks

- `industry_benchmarks` — TPC-C, YCSB workloads
- `sys-perf` — System performance suite
- `linkbench` — Link benchmark suite

## Fetching Task Artifacts

```bash
evergreen fetch -t <task-id> --artifacts
```
