// analyze_coverage.go — P4.6 FTDC coverage map.
//
// FTDC contains structural signal about which mongod code paths were
// actually exercised during a capture. "I ran my chaos test but the
// PALI seal-event count metric stayed at zero" means the test didn't
// actually trigger seals — a coverage gap, not a passing test.
//
// This stage cross-references the metric families that showed activity
// during the capture against a static manifest of "metric family → code
// area" mappings, and against a per-topology expected-active set. Emits:
//
//   kind:"coverage_summary"  — always; lists active + silent areas
//   kind:"coverage_gap"      — when a topology-applicable area was silent
//
// The "expected active" baseline in v1 is a small hardcoded map per
// topology (RS, PALI, sharded, standalone). P4.4 fingerprint-library
// support can replace it with a data-driven baseline derived from
// prior captures of the same topology.
package main

import (
	"sort"
	"strings"
)

// codeAreaManifest maps a metric-family prefix to a human-friendly code
// area name. Order matters: longer (more-specific) prefixes come FIRST
// so the first-match scan picks the most specific area for a metric.
var codeAreaManifest = []struct {
	prefix string
	area   string
}{
	// PALI subsystems (most specific first)
	{"serverStatus.metrics.disagg.logServer.", "pali_log_server"},
	{"serverStatus.metrics.disagg.phylog.", "pali_phylog"},
	{"serverStatus.metrics.disagg.pageServer.", "pali_page_server"},
	{"serverStatus.metrics.disagg.getPageRequest.", "pali_get_page_request"},
	{"serverStatus.metrics.disagg.", "pali_other"},
	// Sharded
	{"serverStatus.shardingStatistics.balancer.", "sharded_balancer"},
	{"serverStatus.shardingStatistics.countDonorMoveChunkStarted", "sharded_balancer"},
	{"serverStatus.shardingStatistics.", "sharded_other"},
	// Replication
	{"serverStatus.metrics.repl.network.", "repl_network"},
	{"serverStatus.metrics.repl.buffer.", "repl_buffer"},
	{"serverStatus.metrics.repl.apply.", "repl_apply"},
	{"serverStatus.metrics.repl.", "repl_other"},
	{"serverStatus.repl.", "repl_state"},
	// WiredTiger
	{"serverStatus.wiredTiger.cache.", "wt_cache"},
	{"serverStatus.wiredTiger.transaction.transaction checkpoint", "wt_checkpoint"},
	{"serverStatus.wiredTiger.log.", "wt_log"},
	{"serverStatus.wiredTiger.concurrentTransactions.", "wt_tickets"},
	{"serverStatus.wiredTiger.", "wt_other"},
	// Core server
	{"serverStatus.opcounters.", "core_opcounters"},
	{"serverStatus.opLatencies.", "core_oplatencies"},
	{"serverStatus.connections.", "core_connections"},
	{"serverStatus.queues.execution.", "core_queues"},
	{"serverStatus.mem.", "core_mem"},
	{"serverStatus.tcmalloc.", "core_tcmalloc"},
	{"serverStatus.metrics.ttl.", "core_ttl"},
	{"serverStatus.metrics.cursor.", "core_cursors"},
	{"serverStatus.metrics.timeseries.", "core_timeseries"},
	{"serverStatus.metrics.changeStreams.", "core_changestreams"},
	// Host OS
	{"systemMetrics.cpu.", "host_cpu"},
	{"systemMetrics.memory.", "host_memory"},
	{"systemMetrics.disks.", "host_disks"},
	{"systemMetrics.vmstat.", "host_vmstat"},
	{"systemMetrics.", "host_other"},
}

// expectedActiveByTopology lists code areas that SHOULD show activity
// when the named topology is detected. Areas in this list that ended
// up silent will emit a coverage_gap.
var expectedActiveByTopology = map[string][]string{
	// PALI captures should exercise all four PALI subsystems.
	"pali": {"pali_log_server", "pali_phylog", "pali_page_server", "core_opcounters", "core_oplatencies"},
	// Sharded.
	"sharded": {"sharded_balancer", "core_opcounters"},
	// Generic replica set or standalone.
	"replica_set": {"core_opcounters", "core_oplatencies", "core_queues", "repl_state"},
	"standalone":  {"core_opcounters", "core_oplatencies", "core_queues"},
}

// emitCoverageMap emits coverage_summary + coverage_gap Findings. Pure
// function of the store; no globals.
func emitCoverageMap(c *findingCollector, host string, store *seriesStore) (emitted int, skipped int) {
	if len(store.names) == 0 {
		return 0, 1
	}

	// Activity per area: count metrics whose observed range is nonzero
	// (max > min) — that's our cheap "active" proxy.
	activity := map[string]int{}
	totalByArea := map[string]int{}
	for id, name := range store.names {
		area := codeAreaFor(name)
		if area == "" {
			continue
		}
		totalByArea[area]++
		if store.observedMax[id] != store.observedMin[id] {
			activity[area]++
		}
	}

	// Build the sorted active/silent split.
	var activeAreas, silentAreas []string
	for area, total := range totalByArea {
		if activity[area] > 0 {
			activeAreas = append(activeAreas, area)
		} else if total > 0 {
			silentAreas = append(silentAreas, area)
		}
	}
	sort.Strings(activeAreas)
	sort.Strings(silentAreas)

	// Determine topology heuristic.
	var topo string
	switch {
	case isPALICapture(store):
		topo = "pali"
	case isShardedCapture(store):
		topo = "sharded"
	case isReplicaSetCapture(store):
		topo = "replica_set"
	default:
		topo = "standalone"
	}

	// Always emit coverage_summary so the analyzer has a stable place to
	// look. idKey is the topology so multi-host captures emit one per
	// host, deterministically sortable.
	c.add(newFinding("coverage_summary", host, topo, map[string]interface{}{
		"topology_inferred": topo,
		"active_areas":      activeAreas,
		"silent_areas":      silentAreas,
		"area_totals":       totalByArea,
		"area_active":       activity,
	}))
	emitted++

	// Emit per-area coverage_gap for any expected-active area that ended
	// up in the silent set. Only fire when the expected list exists for
	// this topology — otherwise we have no opinion.
	expected := expectedActiveByTopology[topo]
	silentSet := map[string]bool{}
	for _, a := range silentAreas {
		silentSet[a] = true
	}
	for _, a := range expected {
		inSilent := silentSet[a]
		_, hasMetrics := totalByArea[a]
		if !inSilent && hasMetrics {
			continue
		}
		reason := "all_static"
		if !hasMetrics {
			reason = "entirely_absent"
		}
		c.add(newFinding("coverage_gap", host, topo+"|"+a, map[string]interface{}{
			"area":              a,
			"topology_inferred": topo,
			"reason":            reason,
			"rationale": "Code area expected to show activity on this topology, but no metric within the area moved during the capture window. The test/workload may not have exercised this path.",
		}))
		emitted++
	}
	return emitted, 0
}

// codeAreaFor returns the area name for a metric, or "" if no manifest
// entry matches. First-match wins; codeAreaManifest is order-sensitive.
func codeAreaFor(name string) string {
	for _, e := range codeAreaManifest {
		if strings.HasPrefix(name, e.prefix) {
			return e.area
		}
	}
	return ""
}

