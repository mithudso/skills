# ftdc-metric-query

**Category:** Databases, Data Engineering & Analytics
**Platform:** Gemini / Plugins
**Original Path:** gemini/mongodb-internal/.agents/skills/tooling/ftdc-analysis/skills/ftdc-metric-query

## Description
How to use the `alexandria` tool and PromQL to explore MongoDB FTDC (diagnostic.data) and reconstruct the state of mongod and WiredTiger at a point in time. Use this skill whenever the user has FTDC / diagnostic.data / `metrics.*` files and wants to query them, graph a metric, find a metric's exact name, translate a T2 / Atlas metric to its Prometheus name, run alexandria's rules to detect known problems, characterize cache / eviction / checkpoint / ticket / replication behavior over a window, or answer "what was this node doing at time T". Triggers on phrasing like "analyze this FTDC", "query the diagnostic data", "what was WiredTiger cache doing at 14:12", "graph the eviction rate", "run alexandria on this", "what's the prom name for this T2 metric", "summarize this metrics file", or any time `alexandria`, `FTDC`, `diagnostic.data`, or a `metrics.YYYY-...` file is mentioned. For a full guided incident root-cause workflow use `mongod-incident-investigation`; this skill is the underlying tool-usage layer.

---

# Analysing FTDC with alexandria and PromQL

`alexandria` is a Swiss-army knife for FTDC analysis. FTDC (Full-Time Diagnostic Data Capture) is the binary `metrics.*` stream MongoDB writes to `diagnostic.data/`; it is a per-second snapshot of `serverStatus`, WiredTiger statistics, replication state, system metrics, and more. `alexandria` parses FTDC and exposes it as a [Prometheus](https://prometheus.io/docs/prometheus/latest/querying/basics/) time series, so any metric can be queried with PromQL, windowed to a point in time, and summarized. It also ships a rules engine that flags known failure patterns.

This skill covers using the tool to **explore** FTDC and **reconstruct mongod / WiredTiger state at a point in time** — the mechanical layer. For a full guided incident root-cause workflow (notes files, hypotheses, source grounding, RCA report), use the `mongod-incident-investigation` skill, which builds on top of this one.

## Build / locate

Use the `alexandria` on PATH if present (`command -v alexandria`). Otherwise build it from `~/alexandria`, cloning the `10gen/alexandria` repo first if that checkout is missing. Go 1.26 panics on the Makefile's `GOEXPERIMENT`, so build plain:

```bash
[ -d ~/alexandria ] || git clone https://github.com/10gen/alexandria.git ~/alexandria
cd ~/alexandria && GOEXPERIMENT= go build -o alexandria cmd/alexandria/alexandria.go
```

Set `ALEX_DISABLE_AUTO_UPDATE=1` to skip the per-run update check (also required for clean JSON output).

`alexandria` accepts a single `metrics.*` file, a whole `diagnostic.data/` directory, or a `.tar.gz` / `.tar.bz2` archive (it decompresses in place). A sharded-node export has **multiple FTDC streams** (`shardServer/`, `router/`, `mongotune/`); the mongod (shard) metrics are in `shardServer/`.

## The four modes

```bash
alexandria <dir-or-file>                          # run all rules → markdown findings (known issues)
alexandria --dump <file>                          # list every metric name in the file
alexandria -rosetta "ss wt cache ..."             # translate T2 ↔ Prometheus ↔ FTDC path
alexandria --human --query "<promql>" <file>      # query one metric → ISO-timestamped CSV
```

1. **Rules run** (`alexandria <path>`) is the fastest way to triage: it slices FTDC into ~15-minute windows, runs the ruleset, ranks failures, and prints a markdown table of the single worst instance per condition (memory pressure, eviction stalls, write-progress stalls, I/O wait, replication, etc.). Start here to see if a known signature fires, then drill in with `--query`.

2. **`--dump`** lists the exact metric names available in *that* file. Names vary by MongoDB version, so never guess — confirm with `--dump <file> | rg -i <term>`.

3. **`-rosetta`** is the metric translator. The same series has three names: the T2 / Atlas display name ("ss wt cache bytes currently in the cache"), the Prometheus name (`mongodb_wiredTiger_cache_bytes_currently_in_the_cache`), and the FTDC path (`ftdc.serverStatus.wiredTiger.cache...`). Rosetta auto-detects which you typed (spaces → T2, underscores → Prom, dots → FTDC, otherwise substring search) and prints all three plus the derived **Prom Expr** (the rate/ratio/scaling T2 applies) and the units. Use it to go from "the metric T2 shows" to "the PromQL I can run".

4. **`--query`** runs PromQL against the file and prints CSV. `--human` gives ISO timestamps (otherwise epoch). Pass `--query` multiple times for multiple series (aligned on a shared time column).

## Exploration workflow

To answer a question about FTDC, work in this order:

1. **Triage with a rules run** (`alexandria <dir>`) to see which known conditions fire and roughly when.
2. **Name the metric.** Translate the user's concept to a Prometheus name with `-rosetta`, or grep the file's metrics with `--dump <file> | rg -i <term>`.
3. **Query it**, converting counters to rates where needed (see PromQL below).
4. **Window to the point in time** with `--from` / `--to`.
5. **Summarize or correlate** — totals, percentiles, or Pearson correlation between two series.
6. **Interpret** the trajectory against the known signatures in `references/signatures.md`.

### Windowing to a point in time

```bash
alexandria --human \
  --from 2026-06-23T22:38:00Z --to 2026-06-23T23:42:00Z \
  --query "<promql>" <file>
```

`--from` / `--to` accept ISO (`2006-01-02T15:04:05.000Z`), epoch seconds or millis, and the keywords `start` / `end` / `now` / `today` / `yesterday`. `--to` also accepts a Go duration (e.g. `13h67m92s`) as an offset. To read state **at** an instant, window a minute or two around it and read the last sample; to read a **trajectory**, window the whole episode and watch the series climb or pin.

### Reading totals and summary statistics

For a counter, the delta over the window (how many ops happened) and the rate distribution (how fast) are both useful:

```bash
# total delta of a counter across the file (first vs last value)
alexandria --query metric1 --query metric2 <file> | (head -3 | tail -2; tail -1) | \
  awk -F, 'NR==1{print} NR==2{for(i=1;i<=NF;i++) a[i]=$i} NR==3{for(i=1;i<=NF;i++) printf "%s%s", $i-a[i], (i==NF?ORS:FS)}'

# min / max / mean / p95 of a counter's rate
alexandria --dm --query "irate(<metric>[60s:1s])" <file> | sed 1d | datamash min 2 max 2 mean 2 perc:95 2
```

`--dm` emits a bare numeric CSV (no human formatting) suited to piping into `datamash` / `awk`.

## PromQL essentials

FTDC counters are monotonically increasing; convert to a per-second rate before reasoning about activity:

```promql
irate(mongodb_wiredTiger_cache_pages_evicted_by_application_threads[120s:1s])
```

Ratios and arithmetic work directly — cache fill % is a ratio of two gauges:

```promql
100 * mongodb_wiredTiger_cache_bytes_currently_in_the_cache / mongodb_wiredTiger_cache_maximum_bytes_configured
```

For grouping a family of metrics by regex, defaulting missing series to zero, boolean thresholds, and Pearson correlation between two series, see `references/promql-patterns.md`.

## Reading mongod / WiredTiger state at a point in time

The metrics that characterize a node's storage-engine state — cache fill / dirty %, eviction walks / workers / aggressive mode, checkpoint-prepare running, tickets and queues, active writers, cursor rates, dhandles — are best found per version with `-rosetta` and `--dump` (above) rather than memorized; names drift across releases. Query them across the window, then map the trajectory to a signature. `references/signatures.md` names the specific series that confirm each known state.

The single most important correctness rule: **an FTDC filename is the rotation/creation epoch, not the first-sample time**, and a frozen stream can keep emitting stale 1-per-second repeats. Before claiming "X stopped at T", verify true sample coverage with the freeze detector (`mongodb_uptimeMillis` — a monotonic field that stops advancing marks where real collection ended). The full set of data-completeness traps is in `references/ftdc-gotchas.md` — read it before drawing any conclusion from FTDC.

## Additional resources

### Reference files

- **`references/promql-patterns.md`** — PromQL recipes: counter→gauge, boolean thresholds, default-to-zero, regex grouping of metric families, and a generator for the Pearson-correlation query between two metrics.