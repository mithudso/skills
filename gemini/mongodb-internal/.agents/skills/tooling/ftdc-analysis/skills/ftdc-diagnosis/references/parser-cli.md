# ftdc_parser CLI reference

Everything about driving the `ftdc_parser` binary: output modes, `--auto` flags, filtering, concurrency, performance, and the bundle/fixture workflows. SKILL.md covers the diagnostic workflow that uses these; this file is the exhaustive flag-level reference.

## Contents
- [Location & building](#location--building)
- [Output modes](#output-modes)
- [`--auto` companion flags](#auto-companion-flags)
- [Bundle output (`--auto-emit-t2`)](#bundle-output---auto-emit-t2)
- [Fixture-addition workflow (`--auto-check`)](#fixture-addition-workflow---auto-check)
- [Filtering](#filtering)
- [Concurrency](#concurrency)
- [Summary-mode issue detection](#summary-mode-issue-detection)
- [Input format](#input-format)
- [Performance](#performance)

## Location & building

```
ftdc-diagnosis/
├── SKILL.md
├── ftdc_parser                # compiled binary (gitignored)
├── src/                       # Go module
│   ├── ftdc_parser.go         # decode + filter + stats + summary
│   ├── analyze.go             # --auto orchestration + Finding schema
│   ├── analyze_store.go       # shared metric-index store, parallel build
│   ├── analyze_correlation.go # Stage 2: lagged Pearson against follower vector
│   ├── analyze_changepoint.go # Stage 3: Mann-Kendall + ED-PELT + CUSUM x-validator + spike
│   ├── analyze_histogram.go   # Stage 4: structural histogram percentile inference
│   ├── analyze_selfcheck.go   # variance-shift CUSUM + run_health aggregator
│   ├── analyze_crosshost.go   # multi-host directory mode
│   ├── analyze_react.go       # --auto-around-event / --auto-segment-by / --auto-correlate
│   ├── analyze_schema.go      # --auto-print-schema (JSON Schema Draft 2020-12)
│   ├── analyze_check.go       # --auto-check test harness
│   ├── analyze_otel.go        # --auto-emit-otel
│   ├── go.mod / go.sum        # module + checksums (adds pgregory.net/changepoint)
│   └── testdata/              # regression fixtures (00-00 standby, 00-01 primary)
├── tools/bias_check.sh        # enforces bias-control invariants over the source
├── BIAS_CONTROL.md            # the invariants bias_check.sh enforces
└── references/                # this file and siblings
```

Build:

```bash
cd .agents/skills/tooling/ftdc-analysis/skills/ftdc-diagnosis/src
go build -o ../ftdc_parser .
```

`--auto` is purely additive: the existing output modes below stay bit-for-bit unchanged. The parser stays a parser; `--auto` is a separate mode you opt into when you want diagnosis-grade Findings.

## Output modes

### Summary (default)
Human-readable: server metadata, operations, connections, cache, checkpoints, admission control, top active metrics, auto-detected issues, FTDC gaps.

### Metadata (`--metadata`)
JSON documents for FTDC type 0 (buildInfo, hostInfo, getCmdLineOpts) and type 2 (getParameter, getClusterParameter).

### Stats (`--stats`)
NDJSON, one line per changed metric:

```json
{"name":"serverStatus.opcounters.insert","kind":"counter","min":0,"max":50665603,"avg":29068782.96,"p50":32126203,"p99":50649503,"first":0,"last":50665603,"delta":50665603,"rate_ps":59327.4}
{"name":"serverStatus.connections.current","kind":"gauge","min":0,"max":130,"avg":121.9,"p50":130,"p99":130,"first":0,"last":0,"delta":0,"rate_ps":0}
```

- `kind`: `"counter"` (monotonically non-decreasing) or `"gauge"` (fluctuates)
- For **counters**: `delta` and `rate_ps` are the useful fields
- For **gauges**: `p50`, `p99`, `min`, `max` are the useful fields

### JSON (`--json`)
One JSON object per sample (~1/sec). Full fidelity — combine with `--filter`.

### Rates (`--rates`)
Per-second deltas between consecutive samples. For counter metrics (opcounters, bytes written).

### CSV (`--csv`)
Standard CSV. Combine with `--filter`.

### List (`--list`)
Metric names with min/max ranges. `*` marks changed metrics.

### `--auto` (analyzer mode, opt-in)
Emits structured Findings NDJSON consumed downstream (see [finding-kinds.md](finding-kinds.md)).

```bash
./ftdc_parser /path/to/diagnostic.data --auto
./ftdc_parser /path/to/cluster_root/ --auto    # multi-host directory of host subdirs
./ftdc_parser --auto-print-schema              # emit JSON Schema for downstream validation
./ftdc_parser --auto-check testdata/           # golden-file regression test harness
```

Each line is one `Finding` JSON object, emitted in deterministic order (kind, then a per-kind stable key). First line is always `kind:"run_metadata"` echoing every algorithmic choice (thresholds, lag grid, window bounds, top-K). Last line is always `kind:"run_health"`.

## `--auto` companion flags

| Flag | Effect |
|---|---|
| `--auto-mode {single,multi,auto}` | Override input layout detection. `auto` (default) picks based on directory structure. |
| `--auto-top-k N` | Top-K movers per direction in Stage 1. Default 50. |
| `--auto-follow METRIC` | Append a metric to the default follower vector for Stage 2. Repeatable. |
| `--auto-drift-ratio R` | Cross-host drift ratio threshold. Default 5.0. |
| `--baseline FROM..TO` | Override baseline window (RFC3339..RFC3339). Default: first 25% of samples. |
| `--stress FROM..TO` | Override stress window. Default: last 25% of samples. |
| `--auto-print-schema` | Print the JSON Schema (Draft 2020-12) for Findings NDJSON and exit. |
| `--auto-emit-otel` | Wrap each Finding in an OTel LogRecord envelope. |
| `--auto-emit-golden` | (with `--auto-check`) Emit current output as the golden baseline. |
| `--auto-emit-score-table` | Emit the full `metric_score_table` (~190K tokens). |
| `--auto-baseline-file PATH` | Compare current run against historical `--stats` NDJSON baseline. |
| `--auto-around-event TIMESTAMP --auto-window DURATION` | Re-run analyzer with samples narrowed to TIMESTAMP ± window/2. |
| `--auto-segment-by METRIC` | Detect change-points on METRIC, emit per-segment Welch stats for ALL metrics around each break. |
| `--auto-correlate LEADER FOLLOWER` | Sub-second-resolution Pearson scan between two specific metrics (lags -60..+60 samples). |
| `--auto-check DIR` | Golden-file regression: runs `--auto` on each `DIR/<fixture>/` and compares to its `expected_findings.ndjson`. |
| `--auto-emit-t2 PATH` | Write a self-contained t2 bundle directory at PATH plus a `<PATH>.tgz` tarball. stdout NDJSON unchanged. |
| `--auto-emit-t2-no-tarball` | Skip the `.tgz` (bundle directory only). Useful for local iteration. |
| `--auto-emit-t2-top-k N` | Marker and custom-view descriptor budget (default 50). |
| `--auto-emit-t2-link-ftdc` | Symlink FTDC files into the bundle instead of copying (single-host dev only). |
| `--auto-emit-t2-zoom` | Emit additional zoom contexts in the t2 state (off by default). |
| `--auto-log-file PATH` | Include a `mongod.log` in the bundle; writes `mongod_log_trimmed.txt` (±30 min of Findings) and inserts a ~100-line snippet into `analysis.md`. |

## Bundle output (`--auto-emit-t2`)

`--auto-emit-t2 PATH` writes a self-contained bundle directory plus a `<PATH>.tgz` tarball suitable for sharing via Slack/Jira/email/scp. The receiver extracts the tarball anywhere and runs `bash open-in-t2.sh` — that single script works on macOS and Linux.

```
PATH/
  context.t2               t2 state file referencing FTDC via relative paths
  analysis.md              run identity + marker glossary + diagnosis summary
                           + log excerpt (if --auto-log-file was provided)
  findings.ndjson          full Finding stream (same bytes as stdout)
  README.txt               walks the receiver through opening + sharing
  open-in-t2.sh            launcher (mode 0755). Run as: bash open-in-t2.sh
  mongod_log_trimmed.txt   (optional) trimmed log when --auto-log-file provided
  ftdc/
    <host_a>/              copy of host_a FTDC files
    <host_b>/              copy of host_b FTDC files (multi-host)
<PATH>.tgz                 gzip tarball of the bundle directory
```

**How to invoke the launcher:** always run `bash open-in-t2.sh`, not `./open-in-t2.sh`, on macOS. A script extracted from a downloaded archive carries the `com.apple.quarantine` extended attribute; Gatekeeper blocks direct execution of quarantined scripts, manifesting as a "malware"/"unidentified developer" alert. Invoking via `bash open-in-t2.sh` audits `/bin/bash` (system-trusted) instead of the script, bypassing the check. On Linux `./open-in-t2.sh` also works.

The `.sh` extension (rather than `.command`) is deliberate: `.command` is the Finder double-click convention but a known malware delivery vector that EDR/AV routinely flag. `.sh` files opened via Finder go to a text editor (no execution), so users must be explicit — the safer pattern.

The `context.t2` state file contains a single primary `Analysis` context with:
- Single-host: wall-clock time axis (`options.timeline:"normal"`); multi-host: aligned axis (`"aligned"`)
- `sourcePerFile:false` so all FTDC chunks for a host merge into one source
- `sourcePerDir:false` single-host; `true` multi-host (one source per host subdir)
- `showValues:true` for hover-value popups
- A data-driven Custom view: only descriptors named by Findings' `metric` fields, plus a baseline bundle (`opcounters`, `connections current`, `cpu user`, `cache dirty bytes`)
- Up to 50 markers (A-Z auto-assigned) at the most significant Finding timestamps
- Per-host `userTags` labels for row identification

The `analysis.md` sidecar contains run identity, the capture window (full ISO 8601 start → end, with a multi-day warning when the capture spans midnight), a marker glossary (A → host, kind, ISO-8601 timestamp, metric, 1-2 sentence blurb), diagnoses, and a per-host summary.

**Timestamp discipline:** whenever you write narrative into `analysis.md`, use full ISO 8601 (`2006-01-02T15:04:05Z`), never bare `HH:MM:SS` — ambiguous when the capture spans days.

### Bundle workflow (read before regenerating)

`--auto-emit-t2` **wipes the entire bundle directory** before writing a fresh skeleton. Any analyst narrative appended to `analysis.md` is silently discarded. To avoid losing work:

1. **Generate skeleton** (first time or after FTDC data changes):
   ```bash
   ./ftdc_parser /path --auto --auto-emit-t2 /output/dir
   ```
2. **Append analyst narrative** to `<dir>/analysis.md`.
3. **Re-seal the tarball** after editing:
   ```bash
   cd <parent-dir> && tar czf <dir>.tgz <dir>/
   ```
4. **Iterative bundle fixes** (patching `context.t2`, `open-in-t2.sh`): use `--auto-emit-t2-no-tarball`, back up `analysis.md` first, restore, re-seal:
   ```bash
   cp <dir>/analysis.md /tmp/analysis_addendum_backup.md
   ./ftdc_parser /path --auto --auto-emit-t2 /output/dir --auto-emit-t2-no-tarball
   cat /tmp/analysis_addendum_backup.md >> <dir>/analysis.md
   cd <parent-dir> && tar czf <dir>.tgz <dir>/
   ```

**Relative paths** in `dataset.files` are resolved by t2 against its launch cwd, not the `context.t2` file's location, so opening `context.t2` directly makes t2 prompt "add data". The launcher rewrites paths to absolute before invoking t2 with `--resume`, which is why it must be used instead of dropping `context.t2` into t2. The "real" fix is the optional t2 patch in `front/ctxmanager.ts:decodePath` (~10-LOC additive change, open PR with STAR team); with it, both launcher and direct open work.

## Fixture-addition workflow (`--auto-check`)

1. Drop FTDC files into `testdata/<short_label>/`.
2. Run `ftdc_parser testdata/<label>/ --auto --auto-emit-golden > testdata/<label>/expected_findings.ndjson`.
3. Hand-edit the golden file to assert ONLY the Findings that matter. Commit both.

The harness normalizes per-run timestamps before comparison and strips `run_metadata` / `run_health` (which vary across runs).

## Filtering

| Option | Effect |
|--------|--------|
| `--filter PATTERN` | Regex filter (case-insensitive). **Repeatable — OR semantics.** Patterns of only `[A-Za-z0-9_ /-.]` or `\.` use an allocation-free substring fast path (~10× faster); patterns with regex metachars use `(?i)PATTERN`. |
| `--nonzero` | Exclude metrics that are zero across all samples. |
| `--from TIME` | Start time (RFC3339). |
| `--to TIME` | End time (RFC3339). |

### Multi-filter (OR semantics)
```bash
./ftdc_parser /path/ --json \
  --filter 'wiredTiger.cache' \
  --filter 'metrics.disagg' \
  --filter 'queues.execution.write'
```
The filter is pushed into chunk decode (before allocating per-sample metric maps), so multi-filter is fast even with broad patterns. Prefer this over a single large alternation regex.

### Forcing real regex semantics
`.` in a literal-eligible pattern is a literal dot (faster, doesn't match arbitrary chars). To force any-char semantics, include another metachar:
```bash
--filter 'wiredTiger[.]cache'   # any-char semantics for the dot
--filter 'wiredTiger\.cache'    # explicit literal dot (also fast path)
--filter 'wiredTiger.cache'     # literal substring (fast path)
```

## Concurrency

| Option | Effect |
|--------|--------|
| `--jobs N` | Parallel zlib-chunk decode workers (default `min(NumCPU, 8)`). `--jobs 1` is strict serial (byte-identical to the pre-multithreading binary). |

Each FTDC file contains many independent zlib chunks; the parser fans them out and re-orders results before emit, so output is in strict time order regardless of `--jobs`. Directory mode iterates files sequentially with chunk-level parallelism inside each file. `--jobs=1` vs `--jobs=8` are guaranteed byte-identical (Tier A.4); if drift appears, file a parser bug. `FTDC_OBSERVABILITY=1` reinstates per-stage `duration_ms` under `_observability` (parser debugging only).

## Summary-mode issue detection

Summary mode flags: cache fill > 80% / > 95%, application-thread evictions, write/read ticket exhaustion (% of samples), long checkpoints (> 60s), aggressive eviction anomaly (WT-13090), high connections (> 5000), eviction worker failure rate > 30%, FTDC gaps (> 5s). These summary verdicts are deliberately NOT echoed into `--auto` — they would pre-bias the downstream narrator.

## Input format

- **A directory** of FTDC chunks.
  - Standard MongoDB form: filenames prefixed `metrics.` (e.g. `metrics.2026-05-17T00-57-52Z-00000`).
  - Atlas-bundled form: raw Unix-timestamp filenames (e.g. `1778979471`). Directory mode picks these up automatically if no `metrics.*` files are present; both styles sort chronologically by name.
- **A single file** — any FTDC chunk file (extension/prefix irrelevant).

Hidden files (leading `.`) are skipped. Non-FTDC files are silently skipped — `partitionFTDC` validates BSON envelopes and yields no samples for garbage input.

## Performance

Measured on a 16-core ARM host, 10 MB primary FTDC file (~19,800 samples × 6,484 metric keys). OLD is the pre-rewrite serial binary.

### Tight single-filter `--json` (`--filter 'opLatencies\.writes'`, ~24 MB output)
| Run | Wall time | Speedup | Peak RSS |
|---|---|---|---|
| OLD | 235.3 s | 1× | n/a |
| NEW `--jobs 1` | 0.93 s | 253× | 71 MB |
| NEW `--jobs 8` | 0.28 s | 839× | 97 MB |
| NEW `--jobs 16` | 0.24 s | 979× | 132 MB |

### Seven broad `--filter` patterns `--json` (~700 metrics/sample)
| Run | Wall time | Notes |
|---|---|---|
| OLD | > 20 min (hung) | — |
| NEW `--jobs 1` | 16.9 s | output dominates |
| NEW `--jobs 16` | 14.7 s | plateau (JSON-encoding-bound) |

Peak RSS ~2.4 GB at this load. Cap `--jobs` and narrow filters if memory is tight.

### `--stats` (full-file aggregate)
| Run | Wall time |
|---|---|
| OLD | 29.1 s |
| NEW `--jobs 1` | 29.1 s (zero regression) |
| NEW `--jobs 16` | 22.7 s |

### Directory mode, 5 files (~27 h of FTDC)
| Run | Wall time |
|---|---|
| OLD | broken (raw-timestamp filenames) |
| NEW `--jobs 1` | 81.5 s |
| NEW `--jobs 16` | 51.6 s |

### Memory footprint
Peak RSS scales with retained-key count × sample count, NOT with `--jobs`:
| Scenario | Peak RSS |
|---|---|
| Tight filter `--json` (~13 keys × 19.8K samples) | 71 MB (`--jobs 1`) → 149 MB (`--jobs 32`) |
| 7 broad filters `--json` (~700 keys) | ~2.4 GB, flat across `--jobs` |
| `--stats` / unfiltered / dir mode unfiltered | multi-GB at full retention; use `--filter` to bound |

If tight on memory: narrow the filter, drop `--jobs`, or run files one at a time.

### Where the speedup comes from
- **Filter pushdown into the chunk decoder.** The keep-mask is computed once per chunk against the chunk's key list; only matched keys ever enter a sample's map.
- **Literal pattern fast path.** Alnum + `_-/. ` patterns bypass `regexp.Regexp` for an allocation-free ASCII-fold `strings.Contains` equivalent.
- **Chunk-parallel decode** with strict time-order emit via a bounded reorder buffer.
- **Buffered stdout** through a 1 MB `bufio.Writer`.

Aggregate modes (`--stats`, default summary) gain less from parallelism because the post-decode statistical pass is single-threaded.
