# PromQL patterns for FTDC

`alexandria --query` takes arbitrary PromQL. FTDC metrics are mostly monotonically increasing counters sampled ~1/second; gauges (cache bytes, queue depths, active threads) are absolute. These recipes cover the transformations that come up repeatedly.

## Counter → gauge (per-second rate)

Counters must be rate-converted before they mean anything. `irate` uses the last two samples in the range; the `[120s:1s]` subquery resamples to 1s steps so the rate is stable:

```promql
irate(counter_metric[120s:1s])
```

A shorter window (`[60s:1s]`) reacts faster but is noisier; a longer one (`[300s:1s]`) smooths.

## Ratios and scaling

Divide two gauges for a percentage; multiply by 100 to match T2 display:

```promql
100 * mongodb_wiredTiger_cache_bytes_currently_in_the_cache / mongodb_wiredTiger_cache_maximum_bytes_configured
```

## Boolean threshold

`bool` turns a comparison into 1 (true) / 0 (false) — useful for marking when a condition held:

```promql
irate(counter_metric[60s:1s]) > bool 0
```

## Default to zero when a series is missing

A metric that does not exist in this version (or has no samples in the window) yields an empty result, which breaks `and`-ing conditions together. Fall back to a zero vector:

```promql
irate(counter_metric[60s:1s]) > 0 or on() vector(0)
```

## Group a family of metrics by regex

Sum every metric matching a name pattern — e.g. all lock-acquire times, all disk metrics, all CPU counters:

```promql
sum(irate(label_join({__name__=~"mongodb_locks_.*_timeAcquiringMicros.*"}, "metric_name", "", "__name__")[60s:1s]))
```

## Pearson correlation between two metrics

A value from 0 to 1 showing how correlated two series are (≥0.7 is high). Useful for testing whether, say, eviction rate tracks cache fill. Generate the query:

```bash
python3 - << 'EOF'
def generate_correlation_query(metric_a, metric_b, timespan="5m"):
    return f"""(
 avg_over_time(({metric_a} * {metric_b})[{timespan}]) -
 (avg_over_time({metric_a}[{timespan}]) * avg_over_time({metric_b}[{timespan}]))
) / (
 stddev_over_time({metric_a}[{timespan}]) * stddev_over_time({metric_b}[{timespan}])
)"""
print(generate_correlation_query("metric_a", "metric_b"))
EOF
```

The formula is `r = cov(X,Y) / (stddev(X) * stddev(Y))`, with `cov(X,Y) = E[XY] - E[X]E[Y]`. Feed the printed expression to `alexandria --query`.
