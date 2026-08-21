# JIRA reference index for FTDC bug shapes

> **Bias-control note (do not delete):** this file is intentionally NOT loaded by the ftdc-diagnosis skill. It is a manual-lookup reference. Loading it into a narration context would prime the LLM with prior-bug verdicts and undo the bias-control posture of the analyzer. Read it yourself when you want to map a structural Finding to a likely existing JIRA; do not surface it pre-emptively in analyzer output.

## Selected mappings (non-exhaustive)

| Symptom | Likely JIRA | What to look for |
|---|---|---|
| Application threads evicting > 30% of WT pages | WT-13090 | `kind:"capacity_pressure"` on cache + `mover.welch_t` on `pages evicted by application threads` |
| RSTL acquisition timeout + ticket exhaustion | SERVER-75205 | `near_ceiling` on `queues.execution` + log evidence |
| Sharded primary memory leak across migrations | SERVER-73915 | `monotonic_trend` on `mem.resident` + `chunk_migration_event` |
| 30-second oplog stalls | SERVER-92554 | `gap` Findings at 30 s spacing |
| Slow TTL deletes from high index-key count | SERVER-XXXXX | `monotonic_trend` on `oplog.timeDiffHours` (when present) + slow TTL log lines |
| Excessive WT log generation | WT-12012 | `monotonic_trend` on `wiredTiger.log.bytes written` + disk-fill |
| API version mismatch during failover | SERVER-106075/SERVER-106138 | `term_change` + log evidence of version negotiation |

## Rule for using this file

1. Surface structural Findings first (movers, changepoints, correlations, diagnosis_candidate, hypothesis_test).
2. Only after the user explicitly asks "is this $TICKET?" or "does this look like a known issue?" consult this index.
3. Cite the JIRA as a *candidate* matched on structural shape, never as the verdict. The parser's bias-control rules forbid the analyzer from issuing verdicts.
