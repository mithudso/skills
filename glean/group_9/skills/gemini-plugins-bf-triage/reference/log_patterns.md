# Log Analysis Patterns & Diagnosis Flows

Condensed for the `bf-triage` skill. Use this when extracting evidence from
`evg_get_raw_task_logs`, `evg_get_test_results_detailed`, and `bb_get_bfg`
fault snippets to fill the **Log evidence** and **Root cause hypothesis**
sections of the report.

## Log fetch & search policy

The full policy lives in `SKILL.md` → **Step 5.5 — Raw-log fetch
policy (CLI-first, MCP backup)**. Quick mental model for this file:

- **Path A — CLI available (default)**:
  1. `evergreen task build TaskLogs --task_id <id> --tail_limit 0
     --print_time --out /tmp/bf-triage-workdir-<BF-KEY>/evg/<task>.txt` —
     one-shot full download, ~2 s, no token cost.
  2. `wc -l` it. If ≤ `BF_TRIAGE_FULL_LOG_MAX_LINES` (default 2000)
     read it all; otherwise keep on disk and use `tail`, `awk`
     time-window slice (timestamps comparable lexically when
     `--print_time` is set), `rg` keyword search.
  3. Per-test logs and artifact bundles also via the CLI into the
     same scratch dir.
  4. **Mandatory cleanup**:
     `rm -r /tmp/bf-triage-workdir-<BF-KEY>/` at end of run.
- **Path B — MCP-only (CLI unavailable / unauthenticated)**: 500-tail
  → full-log-if-small / widen-tail → pivot to per-test log / artifacts.
- **No tool on any path supports server-side keyword search** (verified
  from MCP descriptors + `evergreen ... --help`). Path A's local
  `rg` on the dumped file is the canonical "search the log" mechanism.
  The magic-string and decision-tree tables below are the recommended
  `rg` queries.

## Magic strings — find the error in logs

| Context | Search string |
| ------- | ------------- |
| Parallel test failure | `********** Parallel Test FAILED` |
| Normal test failure | `failed to load` |
| Fuzzer test failure | `expected to be running in teardown(), but wasn't.` |
| Task timeout (find dumps) | `dumping` |
| Test started (regex) | `0000 Running (.*?\.js)\.\.\.` |
| Test finished (regex) | `0000 (.*?\.js) ran in \d+\.\d+ seconds` |
| Stack trace | `BACKTRACE` |
| Assertion | `assert:` |
| Equality failure | `are not equal :` |
| Command failure | `Error: command failed:` |
| Unexpected success | `Error: command worked when it should have failed:` |

## Connection error decision tree

When a test fails with "network error while attempting to run command", the
cause is exactly one of:

1. **Server crash** — `mongod`/`mongos` crashed (invariant failure, fassert)
2. **OOM kill** — process killed by the OS OOM killer
3. **Unexpected stepdown** — `mongod` stepped down from primary

### Tracing in logs

Find the last activity of the server on the port mentioned in the error:

- Deployment via `resmoke.py` → search for `[MongoDFixture:job<jobno>]`,
  `[ReplicaSetFixture:job<jobno>:primary]`,
  `[ShardedClusterFixture:job<jobno>:mongos]`
- Deployment via `MongoRunner` → search for `<port>|`

Then in that process's logs, search for: `segmentation`, `fatal`, `violation`,
`terminate()`, `uncaught`, `invariant`, `No space left on device`.

### Confirming OOM kill

- `Out of memory: Kill process <pid> (mongod) score <N> or sacrifice child`
- Track the PID — if it becomes "defunct" with high recorded memory → OOM

### Confirming unexpected election

- Search replica-set member logs for `dry run`
- Or in the process logs: regex `state transition.*SECONDARY`
- Slow-op symptom: regex `durationMillis":\d{5,}`

## StopError diagnosis

`MongoRunner.StopError` means the process exited nonzero on stop:

1. Process crashed/was killed before `MongoRunner` tried to stop it
2. Process exited nonzero on termination (e.g. exit code 23 = ASan memory leak)

Search for `<port>|` to find that process's messages and look for the fatal
strings above. Return code -9 on POSIX is **not** an error (SIGKILL after
SIGTERM timeout is normal).

## JavaScript assertion failures

1. Open the failing test file
2. Use the assertion message + JS engine backtrace to locate the offending op
3. Identify the connection over which the command was issued
4. Trace backwards using `[conn<conn_num>]` to see prior operations on that
   thread

For concurrency-suite failures, also fetch the FSM workload list and the
debugging guide.

## Fixture teardown failures

When `resmoke.py` logs "Teardown of <fixture_type> (Job #<jobno>) was not
successful":

| Search | Meaning |
| ------ | ------- |
| `was expected to be running in teardown(), but wasn't` | Process already exited (crashed/killed before teardown) |
| `exited with code` (in task logs) | Process returned nonzero on termination |

After identifying the process, follow the connection-error flow above using the
last test's logs.

## Task timeouts

Two timeout types:

- **Command timeout** (30 min on required builders, 2 hours otherwise) — single
  test took too long
- **Task timeout** (6 hours) — cumulative test execution exceeded limit

### Finding the hanging test (single-test timeout)

1. Search task logs for `Running task-timeout commands` — anything after this
   timestamp may be the hang_analyzer attaching debuggers
2. Find unfinished tests: `Running <test_name>...` without matching
   `<test_name> ran in <N> seconds.`
3. Look for `jsTest.log()` / `print()` lines to locate where progress stopped
4. Search hang_analyzer output for `Running Print JavaScript Stack Supplement`
   (Linux/gdb) to get the JS stack
5. In hang_analyzer output, search for `handleRequest` (or `assembleResponse`
   on 3.4 or earlier) to find active client commands, skipping background jobs

### Concurrency-suite timeouts

- Search for the last `Workload(s) started` line to find which workloads were
  still running

### `basic.js` / `basicPlus.js` timeouts

- Find last `S0 Test : <test_name> ...` without matching
  `S0 Test : <test_name> <num_millis>ms`. Repeat for S1, S2, S3.

### `jstestfuzz` suite timeouts

- Last `Top-level statement <##> completed in <###> ms` — the next statement
  in the generated file was hanging
- Long-running statements (post hang_analyzer): regex `completed in [0-9]{5,} ms`
- Slow progress: regex `completed in [0-9]{4,} ms`
- Generated test files are downloadable from "Generated Tests - Execution 0" in
  the Evergreen Files section

### 6-hour task timeout

- Slow disk: search for `took ` in logs; if disk ~100% utilized and
  <60MB/s write throughput (per `mongo-diskstats.tgz`) → likely TIG-407
- Slow commands: regex matching `<millis>ms`

## Disappearing hosts

Spot instances can be decommissioned mid-task. Check host logs (sir-xx link)
for `HOST_TERMINATED_EXTERNALLY`. Should auto-retry as "system unresponsive"
but occasionally surfaces as a task failure.

## Cascading test failures

Tests within a task CAN be dependent. A failure in test A can leave the system
in a bad state, causing test B to fail. When this happens:

- File/dedupe a BF for the original failure in test A
- File a TIG ticket about preventing the cascade

## Sharding DDL-lock deadlock

Signature: a client/test fails with `Failed to acquire DDL lock ... after
300000 ms` / `LockBusy` (code 46), naming another op (often
`reshardCollection`) as the lock holder. This is usually the *symptom*, not
the cause. Diagnosis:

- The `LockBusy` line is the **blocked** op. Trace **what the lock holder is
  itself waiting on** — look in the fixture mongod log (`job0/global.log`,
  see SKILL.md Step 5.5) for the holder's timeline.
- Look for a back-edge forming a cycle, e.g. reshard ↔ moveChunk migration ↔
  range-deletion (orphan cleanup) on the same namespace: migration waits on
  cleanup, cleanup is paused while resharding is in progress, resharding
  waits on the migration to drain.
- Frame the fix as **lock-ordering / drain-before-acquire** (acquire the DDL
  lock only after in-flight migrations have drained), not as a timeout bump.

## Change-stream / CSM cursor-visibility hang

Signature: a `getMore` loop that keeps returning `docsExamined: 0` (cursor
never advances past an optime) together with a **parked-baton** backtrace
(`SharedStateBase::wait`, `AsioNetworkingBaton`, `ParkingLot::parkOne`) that
gets SIGABRT'd by the harness watchdog. Read it as a hang aborted by the
timer, not an active crash. Suspect order:

- **Cursor visibility first** — majority-commit point not advanced, or a lost
  `awaitData` wakeup, so the change-stream / CSM cursor never observes the
  next entry.
- **Completion-promise plumbing second** — a future/promise that never
  resolves.

## Performance-change BF interpretation

Perf-change BFs are auto-generated by the change-point detector
(`Mongo Perf User` reporter, `Performance Change Type =
Regression|Improvement` field, summary starts with
"Performance changes in sys-perf"). They look different from
Build-Baron-derived BFs. Read them as follows:

| Signal in summary / description | Triage interpretation |
| ------------------------------- | --------------------- |
| Summary contains "**fixed by `<sha>`, `<date>`**" before triage | The change-point detector already located a fix-by commit. The detection rule is universal; team-specific disposition (when to "Accept-as-known / Gone away" versus pursue further) lives in [`team_knowledge.md`](team_knowledge.md) § "Performance-change BF disposition rules". |
| Description has roughly equal counts of regressions and improvements on the same task | Strong indicator of a **behavioural / measurement-floor change** in workload setup or rate-limiter, not real performance loss. Lower the severity-level expectation to Low. |
| `Regression: inf%` (or `Improvement: inf%`) in `ErrorsTotal` / `ErrorRate` measurements | Indicates an error-counting metric flipping ON or OFF (test started reporting errors where it previously reported zero, or vice-versa). Almost always a workload / rate-limiter configuration change rather than a server bug. |
| Description: "Overall, no high value workloads were affected" | Suggests the BF is not release-blocking; consider de-escalation. |
| `First Failing Revision` set + the commit message contains a non-`SERVER-` JIRA prefix | The fix may live outside `10gen/mongo` (DSI workload code, infra-tuning repos, etc.). The exact projects to check first are team-dependent — see [`team_knowledge.md`](team_knowledge.md) § "Performance-change BF disposition rules" for the current team's defaults. |
| `First Failing Revision` carries a `SERVER-` prefix BUT the failing variant is sys-perf-only / Mongotune-only / disagg-only AND a known sibling variant on the same task is unaffected | The cause is likely **per-variant config divergence in DSI** (`evergreen/system_perf/master/variants.yml` differences, missing `locust_poetry_override`, missing setParameter, missing/extra task in the variant's `tasks:` list), not a server bug. The perf-change-detector bisects against the mongo waterfall only and so attributes by accident to the first dispatched mongo SHA in the gap. Audit the failing variant's `expansions:` block against the unaffected sibling's block. See [`SKILL.md`](../SKILL.md) § "Step 6 — DSI sub-checklist for sys-perf BFs" and § "Step 6 — Inert-mongo-diff redirect rule". |
| Candidate mongo SHA's diff is inert in this configuration (default-`false` server parameter not enabled by any setParameter on the failing variant; or in a code path the workload doesn't exercise) | Stop further mongo bisecting; redirect to DSI immediately. Do NOT wait for the mongo author's comment to "rule out" the SHA — the gating evidence is sufficient on its own. The next variable to check is the DSI module commit pulled into the same Evergreen run. |
| `ErrorRate` / `ErrorsTotal` regression of two-or-more orders of magnitude (e.g. 0.15 → 28.59 ≈ 190×) on a workload whose open-loop spike pattern is calibrated against a specific client driver | High prior on a **client driver / workload-config change**, not a server bug. Check `configurations/test_control/<failing-task>.yml` for `locust_poetry_override` add/remove and `evergreen/system_perf/master/variants.yml` for variant-level pymongo / driver pins. |

For controller-tuning-artefact cases (older controller hard-codes a
value that the new policy expects to be configurable) the BF is often
closed as `Type of fix = Non-test code change, Resolution = Gone away`
even though no obvious revert happened. Whether to recommend that
disposition by default is team-specific; see `team_knowledge.md`.
