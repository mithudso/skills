# BF Severity Types & Hot/Cold Criteria

Condensed for the `bf-triage` skill. Use this when classifying a BF's severity
type and when filling the **Severity classification** and **Frequency & scope**
sections of the report template.

## Severity Types — log identifiers

The skill should grep extracted faults / raw task logs for these strings to
classify the BF. Auto-detected identifiers come from the Build Baron
`task_log_parser_config.yml`.

| Severity Type | Auto-detected log identifiers | Manual identifiers |
| ------------- | ----------------------------- | ------------------ |
| **Data Inconsistency** | `assert failed : dbhash mismatch`, `"codeName" : "DataCorruptionDetected"`, `collection validation failed` | `dbhash mismatch`, `Failed to apply operation due to missing collection`, `detected one or more invalid documents`, `Unexpected number of documents for collection`, `Dbhash check failed`, `Collection counts did not match` |
| **Memory Pressure** | `Out of memory: Kill process (\d+) \((\w+)\) (.*)` | `OSError - [Errno 12] Cannot allocate memory` |
| **Sanitizer Failure** | `WARNING: ThreadSanitizer`, `: runtime error`, `AddressSanitizer`, `LeakSanitizer` | — |
| **Server Crash** | `Segmentation fault`, `Fatal assertion`, `Invariant failure`, `BACKTRACE:`, `----- BEGIN BACKTRACE -----`, `Unhandled exception`, `^mongod\.exe` | `tassert` does NOT warrant Server Crash severity |
| **Server Hang** | `Cycle detected` | `assert.soon failed`, similar assertions from background hooks |
| **Shell Crash** | (never auto-applied) | `Segmentation fault`, `Got signal: 22 (SIGABRT)`, `DBClientCursor::kill`, `Assertion: BSONObjectTooLarge`, `Got signal: 6`, `Invariant failure`, `unhandled exception` |
| **Query Correctness Failure** | (sometimes auto) | `data consistency checks (point-in-time) failed`, `Failed to find all results`, `The above documents are different`, `Version 2 returned a result set`, `Version 1 returned a result set` |
| **Tripwire Assertion** | `tassert`, `Tripwire assertion` | — |
| **Blocking Commit Queue** | Tasks in the `commit-queue` variant | — |
| **Performance Regression** | (never auto, applied by perf Barons only) | — |
| **Release Infrastructure** | (never auto) | failures mentioning `MacOS notary` |
| **Security Vulnerability** | (never auto, Security Team only) | CVSS >= 4.0, libfuzzer findings |

The skill MUST cite the exact log line that triggered the severity-type
classification in the report's "Severity classification" section.

## Severity Level (customer impact)

| Severity | Definition | Example |
| -------- | ---------- | ------- |
| Low | Uncommon operation + stepdown/planned maintenance → [Severity Type] | tassert during getDatabaseVersion while sharding metadata op holds critical section |
| Medium | Common operation + stepdown/planned maintenance → [Severity Type] | server crash during replication rollback; won't recur on restart |
| High | Common/uncommon operation without stepdown → [Severity Type] OR data inconsistency (uncommon op) | range deleter + elections → orphan docs revived; `$setWindowFields` crash from uninitialized memory |
| Critical | Common/uncommon operation → repeated server crash requiring manual intervention OR data inconsistency (common operation) | server crash in FTDC threads from WiredTiger eviction |

## Hot vs Cold criteria

The Build Baron service auto-scores BFs as Hot (200) or Cold (0). The skill
should evaluate each criterion as a checklist using `bb_search_bfgs` counts and
linked-BFG metadata, and surface the results in the **Frequency & scope**
section.

| Criterion | Threshold |
| --------- | --------- |
| Releasability | At least 1 serious-severity label → release blocker |
| Mainline waterfall | 10+ BFGs linked in past 30 days |
| Frequency | 3+ BFGs AND same task/variant failed 10%+ of runs in 30 days |
| Patch builds | 10+ patch failures from 5+ authors in 30 days |
| Released version | 3+ BFGs from nightly/released AND 10%+ failure rate in 30 days |

## Customer-impact cheat sheet

When triaging without a reproduction:

| Indicator | Severity weight |
| --------- | --------------- |
| Failed in `:ValidateCollections` or `:CheckReplDBHash` suffix hooks | Potential dataloss — never ignore |
| Different results with vs without index | Far-reaching client impact |
| Server crash | Potential availability bug |
| Other server bug | Severity depends on context |
| Test issue or system issue | Not as critical, but still blocks CI |
