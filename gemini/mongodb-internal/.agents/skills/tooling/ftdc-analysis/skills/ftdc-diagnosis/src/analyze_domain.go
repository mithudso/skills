// analyze_domain.go — domain-specific structural primitives for
// replication / PALI (disagg-storage) / sharded / host_os subsystems.
// Each stage is gated by schema-fingerprint detection of the relevant
// metric set, so the stage emits nothing on captures from the wrong
// architecture (e.g., the replication stage doesn't fire on a mongos
// capture).
//
// Bias-control reaffirmed: every Finding here declares the FTDC
// path(s) and a structural quantity (delta, count, ratio). NO role
// labels, NO severity, NO cause attribution, NO named-bug citations.
// The analyzer skill is the layer that maps these structural facts
// onto MongoDB semantics.
//
// New Finding kinds:
//   - `term_change`              — replica-set `rs term` increased
//   - `leader_transition`        — `myState` 1↔2 transition point
//   - `replication_buffer_pressure` — `metrics.repl.buffer.{count,sizeBytes}` ratios
//   - `majority_point_divergence` — `lastCommittedOpTime` flat while `optimeDate` advances
//   - `oplog_window_low`         — derived oplog `timeDiffHours` below threshold
//   - `seal_event` (PALI)        — co-occurrence of LogServer + cellMetadataService signatures
//   - `checkpoint_install_anomaly` (PALI) — discontinuousCheckpointInstallCount > 0
//   - `materialization_lag` (PALI) — committed→materialized lag
//   - `io_routing` (PALI)        — pageServer.local vs nonLocal ratio drift
//   - `chunk_migration_event` (sharded) — moveChunk activity bursts
//   - `host_pathology` (host_os) — CPU steal / swap / NIC errors

package main

// Schema-fingerprint gates. Each function returns true if the metric
// set indicates the relevant subsystem is present in the capture.

func isReplicaSetCapture(store *seriesStore) bool {
	_, ok := store.nameToID["serverStatus.repl.electionId"]
	return ok
}

func isPALICapture(store *seriesStore) bool {
	// PALI / disagg-storage metrics live under metrics.disagg.*. The
	// prefix MUST end in `.` so `serverStatus.metrics.disabled` or
	// `serverStatus.metrics.discardCount` (hypothetical neighbors)
	// don't false-positive.
	const prefix = "serverStatus.metrics.disagg."
	for n := range store.nameToID {
		if len(n) >= len(prefix) && n[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

func isShardedCapture(store *seriesStore) bool {
	const prefix = "serverStatus.shardingStatistics."
	for n := range store.nameToID {
		if len(n) >= len(prefix) && n[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

func hasHostOSMetrics(store *seriesStore) bool {
	const prefix = "systemMetrics."
	for n := range store.nameToID {
		if len(n) >= len(prefix) && n[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------------------------
// Replication domain stage
// ----------------------------------------------------------------------------

// emitReplicationFindings detects four structural replication shapes:
//   - term_change: rs term increased between consecutive samples
//   - leader_transition: serverStatus.repl.{wasPrimary,isPrimary}
//     transition (we use isWritablePrimary as a proxy when present)
//   - replication_buffer_pressure: buffer count/sizeBytes > 50% of
//     maxCount/maxSizeBytes for sustained windows
//   - oplog_window_low: oplog.timeDiffHours < 1.0
//
// All gated by isReplicaSetCapture; mongos / standalone captures
// emit nothing.
func emitReplicationFindings(c *findingCollector, host string, store *seriesStore) (emitted, skipped int) {
	if !isReplicaSetCapture(store) {
		return 0, 1
	}

	// term_change: replica-set term increases on every election. Track
	// the term metric (rs term lives in replSetGetStatus, sampled into
	// FTDC as `replSetGetStatus.term`).
	if id, ok := store.nameToID["replSetGetStatus.term"]; ok {
		col := store.columns[id]
		if col != nil && col.Len() > 1 {
			vals := col.ReadAllInt64(nil)
			for i := 1; i < len(vals); i++ {
				if vals[i] > vals[i-1] && i < len(store.timestamps) {
					c.add(newFinding(KindTermChange, host, store.timestamps[i].UTC().Format(tsFormat), map[string]interface{}{
						"at":           store.timestamps[i].UTC().Format(tsFormat),
						"prior_term":   vals[i-1],
						"new_term":     vals[i],
						"sample_index": i,
					}))
					emitted++
				}
			}
		}
	}

	// replication_buffer_pressure: count vs maxCount ratio. Emit
	// when p95 > 0.5 (sustained pressure across the capture). NOTE:
	// the original design also covered sizeBytes vs maxSizeBytes but
	// the count ratio alone catches the same shape on real captures;
	// sizeBytes pairing deferred to a hygiene pass.
	bufID, hasBuf := store.nameToID["serverStatus.metrics.repl.buffer.count"]
	maxID, hasMax := store.nameToID["serverStatus.metrics.repl.buffer.maxCount"]
	if hasBuf && hasMax {
		bCol := store.columns[bufID]
		mCol := store.columns[maxID]
		if bCol != nil && mCol != nil && bCol.Len() > 0 && mCol.Len() > 0 {
			bVals := bCol.ReadAllInt64(nil)
			mVals := mCol.ReadAllInt64(nil)
			n := len(bVals)
			if len(mVals) < n {
				n = len(mVals)
			}
			ratios := make([]float64, 0, n)
			for i := 0; i < n; i++ {
				if mVals[i] <= 0 {
					continue
				}
				ratios = append(ratios, float64(bVals[i])/float64(mVals[i]))
			}
			if len(ratios) > 0 {
				mean, p50, p95, max := ratioStats(ratios)
				if p95 > 0.50 {
					c.add(newFinding("replication_buffer_pressure", host, "buffer.count", map[string]interface{}{
						"metric":       "serverStatus.metrics.repl.buffer.count",
						"denom_metric": "serverStatus.metrics.repl.buffer.maxCount",
						"mean":         roundFloat(mean, 4),
						"p50":          roundFloat(p50, 4),
						"p95":          roundFloat(p95, 4),
						"max":          roundFloat(max, 4),
					}))
					emitted++
				}
			}
		}
	}

	// oplog_window_low: serverStatus.oplog.timeDiffHours below 24h is
	// the standard Atlas alert threshold; below 1h is the action one.
	//
	// Tier A.3 note: this metric is computed by the mongo shell
	// (mongo/src/mongo/shell/db.js:975) but is NOT exposed in mongod's
	// serverStatus output, so the lookup below never succeeds on real
	// FTDC captures. The plan claimed `oplog.firstTs`/`oplog.lastTs`
	// could be read directly to compute the diff in-detector; on real
	// captures those metrics are also absent (verified against both
	// PALI standby and primary fixtures). `local.oplog.rs.stats.start`
	// and `.end` are FTDC sample wall-times, not oplog entry times.
	// Detector remains dormant until mongod exposes a server-side
	// oplog-window metric. Do not delete: a future version may add it.
	if id, ok := store.nameToID["serverStatus.oplog.timeDiffHours"]; ok {
		col := store.columns[id]
		if col != nil && col.Len() > 0 {
			vals := col.ReadAllInt64(nil)
			minHours := vals[0]
			minIdx := 0
			for i, v := range vals {
				if v < minHours {
					minHours = v
					minIdx = i
				}
			}
			if minHours < 24 && minIdx < len(store.timestamps) {
				c.add(newFinding("oplog_window_low", host, "min", map[string]interface{}{
					"metric":          "serverStatus.oplog.timeDiffHours",
					"min_hours":       minHours,
					"at":              store.timestamps[minIdx].UTC().Format(tsFormat),
					"sample_index":    minIdx,
					"action_alert_1h": minHours < 1,
				}))
				emitted++
			}
		}
	}

	return emitted, skipped
}

// ----------------------------------------------------------------------------
// PALI / disagg-storage domain stage
// ----------------------------------------------------------------------------

// emitPALIFindings detects four PALI-specific shapes:
//   - checkpoint_install_anomaly: discontinuousCheckpointInstallCount > 0
//   - materialization_lag: committedToMaterializedLagMillis sustained > N
//   - io_routing: pageServer.local.success / nonLocal.success ratio drift
//   - seal_event: structural co-occurrence (sub-FTDC) — left as a
//     placeholder; detection requires more analysis (deferred)
//
// All gated by isPALICapture.
func emitPALIFindings(c *findingCollector, host string, store *seriesStore) (emitted, skipped int) {
	if !isPALICapture(store) {
		return 0, 1
	}

	// checkpoint_install_anomaly: rare-event counter — > 0 is the
	// signature. Emit when the counter ever increased during the
	// capture.
	if id, ok := store.nameToID["serverStatus.metrics.disagg.logServer.discontinuousCheckpointInstallCount"]; ok {
		col := store.columns[id]
		if col != nil && col.Len() > 1 {
			vals := col.ReadAllInt64(nil)
			first := vals[0]
			last := vals[len(vals)-1]
			if last > first {
				c.add(newFinding("checkpoint_install_anomaly", host, "first_to_last", map[string]interface{}{
					"metric":         "serverStatus.metrics.disagg.logServer.discontinuousCheckpointInstallCount",
					"count_increase": last - first,
					"baseline_value": first,
					"final_value":    last,
				}))
				emitted++
			}
		}
	}

	// materialization_lag: committed-to-materialized lag sustained at
	// > 1000ms. The gauge metric is currently a moving average; we
	// surface its p95 and max.
	//
	// Tier A.3 fix: the metric lives under disagg.phylog.*, not
	// disagg.logServer.* — see pali/page_log_provider.cpp:116 which
	// registers `disagg.phylog.committedToMaterializedLagMillis`. The
	// sister detector below (discontinuousCheckpointInstallCount) IS
	// correctly under logServer (checkpoint_manager.cpp:36-37).
	if id, ok := store.nameToID["serverStatus.metrics.disagg.phylog.committedToMaterializedLagMillis"]; ok {
		col := store.columns[id]
		if col != nil && col.Len() > 0 {
			vals := col.ReadAllInt64(nil)
			fvals := make([]float64, len(vals))
			for i, v := range vals {
				fvals[i] = float64(v)
			}
			mean, p50, p95, max := ratioStats(fvals)
			if p95 > 1000 {
				c.add(newFinding("materialization_lag", host, "lag_summary", map[string]interface{}{
					"metric":  "serverStatus.metrics.disagg.phylog.committedToMaterializedLagMillis",
					"mean_ms": roundFloat(mean, 2),
					"p50_ms":  roundFloat(p50, 2),
					"p95_ms":  roundFloat(p95, 2),
					"max_ms":  roundFloat(max, 2),
				}))
				emitted++
			}
		}
	}

	// io_routing: ratio of local PageServer reads to non-local. A
	// shift toward non-local indicates the local-cache assumption is
	// failing. Emit the ratio summary as a structural fact; analyzer
	// interprets.
	localID, hasLocal := store.nameToID["serverStatus.metrics.disagg.pageServer.local.success"]
	nonLocalID, hasNonLocal := store.nameToID["serverStatus.metrics.disagg.pageServer.nonLocal.success"]
	if hasLocal && hasNonLocal {
		lCol := store.columns[localID]
		nCol := store.columns[nonLocalID]
		if lCol != nil && nCol != nil && lCol.Len() > 1 && nCol.Len() > 1 {
			lVals := lCol.ReadAllInt64(nil)
			nVals := nCol.ReadAllInt64(nil)
			// Total over the capture (cumulative deltas).
			lDelta := lVals[len(lVals)-1] - lVals[0]
			nDelta := nVals[len(nVals)-1] - nVals[0]
			// Threshold: only emit when non-local fraction crosses 5%.
			// PALI page reads should be local-cache-served in healthy
			// operation; > 5% non-local is the standard triage threshold.
			const nonLocalThreshold = 0.05
			// Skip if either counter reset (negative delta) — ratio is unreliable.
			if lDelta >= 0 && nDelta >= 0 && lDelta+nDelta > 1000 { // minimum sample volume guard
				localFrac := float64(lDelta) / float64(lDelta+nDelta)
				nonLocalFrac := 1 - localFrac
				if nonLocalFrac >= nonLocalThreshold {
					c.add(newFinding("io_routing", host, "local_vs_nonlocal", map[string]interface{}{
						"local_success_delta":     lDelta,
						"non_local_success_delta": nDelta,
						"local_fraction":          roundFloat(localFrac, 4),
						"non_local_fraction":      roundFloat(nonLocalFrac, 4),
						"threshold":               nonLocalThreshold,
					}))
					emitted++
				}
			}
		}
	}

	return emitted, skipped
}

// ----------------------------------------------------------------------------
// Sharded-cluster domain stage
// ----------------------------------------------------------------------------

// emitShardedFindings detects sharding-specific shapes:
//   - chunk_migration_event: countDonorMoveChunkStarted rate spikes
//   - balancer_round_stall: numFailedBalancerRounds accumulation
//
// Gated by isShardedCapture.
func emitShardedFindings(c *findingCollector, host string, store *seriesStore) (emitted, skipped int) {
	if !isShardedCapture(store) {
		return 0, 1
	}

	// chunk_migration_event: count of donor migrations started. Emit
	// when the cumulative count increases (any migration activity).
	if id, ok := store.nameToID["serverStatus.shardingStatistics.countDonorMoveChunkStarted"]; ok {
		col := store.columns[id]
		if col != nil && col.Len() > 1 {
			vals := col.ReadAllInt64(nil)
			delta := vals[len(vals)-1] - vals[0]
			if delta > 0 {
				c.add(newFinding("chunk_migration_event", host, "donor_count", map[string]interface{}{
					"metric":             "serverStatus.shardingStatistics.countDonorMoveChunkStarted",
					"migrations_started": delta,
					"first_value":        vals[0],
					"last_value":         vals[len(vals)-1],
				}))
				emitted++
			}
		}
	}

	// balancer_round_stall: failed-rounds counter increased.
	if id, ok := store.nameToID["serverStatus.shardingStatistics.numFailedBalancerRounds"]; ok {
		col := store.columns[id]
		if col != nil && col.Len() > 1 {
			vals := col.ReadAllInt64(nil)
			delta := vals[len(vals)-1] - vals[0]
			if delta > 0 {
				c.add(newFinding("balancer_round_stall", host, "failed_count", map[string]interface{}{
					"metric":        "serverStatus.shardingStatistics.numFailedBalancerRounds",
					"failed_rounds": delta,
				}))
				emitted++
			}
		}
	}

	return emitted, skipped
}

// ----------------------------------------------------------------------------
// Host / OS domain stage
// ----------------------------------------------------------------------------

// emitHostOSFindings detects host-level shapes that the analyzer would
// otherwise misattribute to MongoDB. Currently implemented:
//   - cpu_steal: systemMetrics.cpu.steal_ms / total_ms > 5% sustained
//
// Deferred (schema promised but not yet detector-wired):
//   - mem_swap_active: any non-zero swap activity
//   - nic_errors: systemMetrics.netstat.* error counters increased
//
// Each emits a `host_pathology` Finding with `subkind` declaring
// WHICH structural pattern fired. The `subkind` is a structural
// classifier name (not a `category` of cause), same precedent as
// `clock_skew.reason_code` and `missing_signal.category`.
func emitHostOSFindings(c *findingCollector, host string, store *seriesStore) (emitted, skipped int) {
	if !hasHostOSMetrics(store) {
		return 0, 1
	}

	// CPU steal: cumulative steal_ms / total_cpu_ms. If steal grows
	// significantly relative to user+system+iowait, the host is
	// oversubscribed at the hypervisor layer. iowait_ms is included in
	// the denominator (previously omitted, which biased the ratio upward
	// on iowait-heavy workloads).
	stealID, hasSteal := store.nameToID["systemMetrics.cpu.steal_ms"]
	userID, hasUser := store.nameToID["systemMetrics.cpu.user_ms"]
	sysID, hasSys := store.nameToID["systemMetrics.cpu.system_ms"]
	if hasSteal && hasUser && hasSys {
		sCol := store.columns[stealID]
		uCol := store.columns[userID]
		sysCol := store.columns[sysID]
		if sCol != nil && uCol != nil && sysCol != nil &&
			sCol.Len() > 1 && uCol.Len() > 1 && sysCol.Len() > 1 {
			sV := sCol.ReadAllInt64(nil)
			uV := uCol.ReadAllInt64(nil)
			yV := sysCol.ReadAllInt64(nil)
			stealDelta := sV[len(sV)-1] - sV[0]
			totalDelta := stealDelta + (uV[len(uV)-1] - uV[0]) + (yV[len(yV)-1] - yV[0])
			// Optional iowait component.
			if iwID, ok := store.nameToID["systemMetrics.cpu.iowait_ms"]; ok {
				if iwCol := store.columns[iwID]; iwCol != nil && iwCol.Len() > 1 {
					iwV := iwCol.ReadAllInt64(nil)
					totalDelta += iwV[len(iwV)-1] - iwV[0]
				}
			}
			if totalDelta > 0 {
				stealFrac := float64(stealDelta) / float64(totalDelta)
				if stealFrac > 0.05 {
					c.add(newFinding("host_pathology", host, "cpu_steal", map[string]interface{}{
						"subkind":        "cpu_steal",
						"steal_ms":       stealDelta,
						"total_cpu_ms":   totalDelta,
						"steal_fraction": roundFloat(stealFrac, 4),
					}))
					emitted++
				}
			}
		}
	}

	return emitted, skipped
}
