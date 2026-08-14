// metric_docs.go — Tier C.3: extracted from analyze_glossary.go.
//
// metricDocsCatalog is MongoDB-aware policy data. It maps fully-qualified
// FTDC metric paths to short, neutral, structural docs (units, accumulation
// semantics, ratio pairings). It is NOT detector logic — no branching
// reads from it during a detection pass. Stage emit checks the catalog
// only AFTER a Finding has been formed, to attach a `metric_doc` field
// for the analyzer's downstream context.
//
// Editorial rules for entries:
//   - Counter / gauge / histogram classification
//   - Units (bytes, microseconds, milliseconds, count, fraction, %)
//   - Ratio pairings ("pair with X for Y")
//   - NO verdict vocabulary (no "high is bad", "approaching ceiling means…")
//   - NO MongoDB version-specific historical context
//
// Adding entries: just append to the map literal. CI lint (Tier C.2) does
// NOT scan this file's entries — paths in DATA files are policy, not
// branching logic. Reviewer checks against the editorial rules above.
//
// Per Tier C: this is the ONLY catalog-data Go file. capacityRules and
// utilizationRules stay as typed structs in analyze_capacity.go (a YAML
// migration would lose go-to-definition for negligible governance gain
// on a single-engineer project).

package main

var metricDocsCatalog = map[string]string{
	// --- Top-level connections / network ---
	"serverStatus.connections.current":      "Gauge. Currently-open client connections.",
	"serverStatus.connections.available":    "Gauge. Connections still acceptable before the configured cap. (current + available) = max.",
	"serverStatus.connections.totalCreated": "Counter. Cumulative connections ever created.",
	"serverStatus.connections.rejected":     "Counter. Cumulative connections refused.",
	"serverStatus.network.bytesIn":          "Counter. Cumulative bytes received over all client sockets.",
	"serverStatus.network.bytesOut":         "Counter. Cumulative bytes sent over all client sockets.",
	"serverStatus.network.numRequests":      "Counter. Cumulative protocol-level requests received.",
	"serverStatus.network.tcp.retransmits":  "Counter. Cumulative TCP retransmits.",

	// --- opcounters / opLatencies ---
	"serverStatus.opcounters.insert":                "Counter. User-initiated insert ops (excludes replication-applied writes).",
	"serverStatus.opcounters.query":                 "Counter. User-initiated find() ops. Does NOT include aggregate(), which is counted in opcounters.command.",
	"serverStatus.opcounters.update":                "Counter. User-initiated update ops.",
	"serverStatus.opcounters.delete":                "Counter. User-initiated delete ops.",
	"serverStatus.opcounters.command":               "Counter. Commands including aggregate, find-via-command, isMaster, etc.",
	"serverStatus.opcountersRepl.insert":            "Counter. Replicated insert ops applied on secondaries.",
	"serverStatus.opcountersRepl.update":            "Counter. Replicated update ops on secondaries.",
	"serverStatus.opcountersRepl.delete":            "Counter. Replicated delete ops on secondaries.",
	"serverStatus.opLatencies.reads.ops":            "Counter. Read-op count paired with .latency to derive mean latency.",
	"serverStatus.opLatencies.reads.latency":        "Counter. Cumulative read-op latency in microseconds. latency/ops = mean μs.",
	"serverStatus.opLatencies.writes.ops":           "Counter. Write-op count paired with .latency.",
	"serverStatus.opLatencies.writes.latency":       "Counter. Cumulative write-op latency in microseconds.",
	"serverStatus.opLatencies.commands.ops":         "Counter. Command-op count.",
	"serverStatus.opLatencies.commands.latency":     "Counter. Cumulative command-op latency in microseconds.",
	"serverStatus.opLatencies.transactions.ops":     "Counter. Transaction-commit count.",
	"serverStatus.opLatencies.transactions.latency": "Counter. Cumulative transaction-commit latency in microseconds.",

	// --- Ticket queues / concurrency ---
	"serverStatus.queues.execution.write.out":                        "Gauge. Write tickets currently held. (out + available) = total.",
	"serverStatus.queues.execution.write.available":                  "Gauge. Write tickets available.",
	"serverStatus.queues.execution.read.out":                         "Gauge. Read tickets currently held.",
	"serverStatus.queues.execution.read.available":                   "Gauge. Read tickets available.",
	"serverStatus.wiredTiger.concurrentTransactions.read.out":        "Gauge. WT read tickets currently held (legacy ticket model).",
	"serverStatus.wiredTiger.concurrentTransactions.read.available":  "Gauge. WT read tickets available.",
	"serverStatus.wiredTiger.concurrentTransactions.write.out":       "Gauge. WT write tickets currently held.",
	"serverStatus.wiredTiger.concurrentTransactions.write.available": "Gauge. WT write tickets available.",

	// --- WiredTiger cache ---
	"serverStatus.wiredTiger.cache.bytes currently in the cache":                     "Gauge. WT cache occupancy in bytes. Pair with `maximum bytes configured` for ratio.",
	"serverStatus.wiredTiger.cache.maximum bytes configured":                         "Gauge. Configured WT cache ceiling in bytes (set via wiredTigerCacheSizeGB).",
	"serverStatus.wiredTiger.cache.tracked dirty bytes in the cache":                 "Gauge. Dirty (modified-not-flushed) bytes in cache. WT default eviction_dirty_trigger = 20% of cache.",
	"serverStatus.wiredTiger.cache.pages evicted by application threads":             "Counter. Cumulative pages evicted by application threads (vs background workers).",
	"serverStatus.wiredTiger.cache.eviction worker thread evicting pages":            "Counter. Cumulative pages evicted by background eviction workers.",
	"serverStatus.wiredTiger.cache.eviction server unable to reach eviction goal":    "Counter. Cumulative cycles where the eviction server did not meet its target.",
	"serverStatus.wiredTiger.cache.pages selected for eviction unable to be evicted": "Counter. Pages the eviction server picked but could not free (locked, dirty, or refcount > 0).",

	// --- WiredTiger checkpoints / reconciliation ---
	"serverStatus.wiredTiger.transaction.transaction checkpoint currently running":        "Gauge. 0/1 — 1 while a checkpoint is in flight.",
	"serverStatus.wiredTiger.transaction.transaction checkpoint most recent time (msecs)": "Gauge. Wall time the most recent checkpoint took, in ms.",
	"serverStatus.wiredTiger.transaction.transaction checkpoint max time (msecs)":         "Gauge. Maximum checkpoint duration since process start, in ms (monotonic).",

	// --- Memory / tcmalloc ---
	"serverStatus.mem.resident":                              "Gauge. Process resident set size in MB.",
	"serverStatus.mem.virtual":                               "Gauge. Process virtual-address-space size in MB.",
	"serverStatus.tcmalloc.generic.current_allocated_bytes":  "Gauge. Bytes currently allocated to MongoDB per tcmalloc accounting.",
	"serverStatus.tcmalloc.generic.heap_size":                "Gauge. Total tcmalloc heap including unused pages.",
	"serverStatus.tcmalloc.tcmalloc.pageheap_free_bytes":     "Gauge. Free pages in tcmalloc pageheap (returnable to OS).",
	"serverStatus.tcmalloc.tcmalloc.pageheap_unmapped_bytes": "Gauge. Pages returned to OS by tcmalloc.",

	// --- Replication ---
	"serverStatus.repl.electionId":                        "Gauge. Replica-set election identifier; changes indicate a new election occurred.",
	"serverStatus.metrics.repl.buffer.count":              "Gauge. Oplog buffer entry count on secondary. Pair with maxCount.",
	"serverStatus.metrics.repl.buffer.sizeBytes":          "Gauge. Oplog buffer bytes on secondary. Pair with maxSizeBytes.",
	"serverStatus.metrics.repl.buffer.maxCount":           "Gauge. Configured oplog-buffer entry-count cap.",
	"serverStatus.metrics.repl.buffer.maxSizeBytes":       "Gauge. Configured oplog-buffer byte cap.",
	"serverStatus.metrics.repl.apply.batches.num":         "Counter. Cumulative apply batches on secondary.",
	"serverStatus.metrics.repl.apply.batches.totalMillis": "Counter. Cumulative time secondary spent applying batches, in ms.",

	// --- Sharding ---
	"serverStatus.shardingStatistics.countDonorMoveChunkStarted":          "Counter. Migrations started where this shard was the donor.",
	"serverStatus.shardingStatistics.numBalancerRounds":                   "Counter. Balancer rounds executed (emits only on config server primary).",
	"serverStatus.shardingStatistics.numFailedBalancerRounds":             "Counter. Balancer rounds that returned an error.",
	"serverStatus.shardingStatistics.catalogCache.countStaleConfigErrors": "Counter. StaleConfig errors hitting the routing cache.",

	// --- PALI / disagg storage ---
	"serverStatus.metrics.disagg.getPageRequest.numQueued":                      "Gauge. PageServer get-page requests currently queued.",
	"serverStatus.metrics.disagg.paliHandleGet.paliHandleGetTotalLatencyUs":     "Gauge. Running-mean of paliHandleGet latency in μs (MovingAverage, smoothed).",
	"serverStatus.metrics.disagg.logServer.lastWrittenLSN":                      "Gauge. Most-recently-written LSN observed by the LogServer.",
	"serverStatus.metrics.disagg.logServer.lastReadOplogLSN":                    "Gauge. Most-recently-read oplog LSN observed by the standby reader.",
	"serverStatus.metrics.disagg.phylog.committedToMaterializedLagMillis":       "Gauge. Time delta in ms between committed and materialized LSN.",
	"serverStatus.metrics.disagg.logServer.discontinuousCheckpointInstallCount": "Counter. Rare-event counter for checkpoint-install gaps.",

	// --- Host / OS ---
	"systemMetrics.cpu.user_ms":            "Counter. Cumulative user-mode CPU milliseconds summed across cores.",
	"systemMetrics.cpu.system_ms":          "Counter. Cumulative system-mode CPU milliseconds.",
	"systemMetrics.cpu.iowait_ms":          "Counter. Cumulative iowait CPU milliseconds.",
	"systemMetrics.cpu.steal_ms":           "Counter. Cumulative steal CPU milliseconds (hypervisor scheduling).",
	"systemMetrics.memory.MemAvailable_kb": "Gauge. Linux MemAvailable estimate in KB.",
}
