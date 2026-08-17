// tunable_guardrails.go — P1.NEW-A static catalog of tunable observations.
//
// Each entry declares a mongod tunable the parser may emit a
// `tunable_observation` Finding for. The Finding NEVER carries a
// recommended numeric value — only direction + observed evidence +
// verification recipe. The bias-control contract treats this as a
// Tier-7 "narration scaffold" (analogous to Tier-R diagnosis_candidate).
//
// To add a tunable: append a tunableRule. To add an anti-recommendation:
// append to antiRecommendationCatalog. Both data shapes are pure structs;
// no per-tunable Go logic lives here.
package main

// tunableRule declares one tunable + its FTDC-derived evidence + the
// guardrails the analyzer must surface verbatim.
type tunableRule struct {
	tunable                 string   // setParameter / config name
	tunableKind             string   // "setParameter" | "config" | "startup_only"
	topologyApplicability   []string // subset of {standalone, replica_set_primary, replica_set_secondary, mongos, config_server, pali_primary, pali_standby}
	directionWhenObserved   string   // "increase" | "decrease" | "investigate"
	requiredFindingKinds    []string // ALL must be present in the run
	requiredMetricNames     []string // raw metric paths that must exist for the verification recipe to resolve
	verificationRecipe      string   // copy-pasteable; references metric families
	verificationHelped      string   // analyzer narrates verbatim
	verificationNotHelped   string   // analyzer narrates verbatim
	guardrails              []string // analyzer narrates verbatim
}

// antiRecommendationGate declares a structural pattern that BLOCKS
// emission of one or more tunables (or downgrades them to "refuse"
// direction). The Tunables field lists which tunables the gate applies
// to; an empty Tunables means "applies to all tunables this run."
type antiRecommendationGate struct {
	code                     string   // stable identifier (e.g. "max_connections_low_utilization")
	rationale                string   // narrated verbatim
	requiredKindsPresent     []string // all must fire
	requiredKindsAbsent      []string // none of these may fire
	tunablesBlocked          []string // empty = all
}

// tunableCatalog: the static set of tunables the parser knows. Order
// matters for emission determinism (alphabetical by tunable name).
var tunableCatalog = []tunableRule{
	{
		tunable:               "wiredTigerConcurrentWriteTransactions",
		tunableKind:           "setParameter",
		topologyApplicability: []string{"standalone", "replica_set_primary", "replica_set_secondary", "config_server", "pali_primary", "pali_standby"},
		directionWhenObserved: "increase",
		requiredFindingKinds:  []string{"capacity_pressure", "near_ceiling"},
		requiredMetricNames:   []string{"serverStatus.wiredTiger.concurrentTransactions.write.out"},
		verificationRecipe:    "ftdc_parser <next-capture> --auto --filter 'concurrentTransactions|opLatencies'",
		verificationHelped:    "capacity_pressure no longer fires on write tickets; opLatencies.writes p99 drops by ≥10%.",
		verificationNotHelped: "cache_at_ceiling shifts (capacity_pressure on wiredTiger.cache.bytes fires instead).",
		guardrails: []string{
			"Effect on cache pressure: more concurrent writers amplify cache fill rate.",
			"Effect on lock contention: no direct FTDC signal; consult lock latency histograms.",
			"On 7.0+ topologies prefer storageEngineConcurrencyAdjustmentAlgorithm=throughputProbing.",
		},
	},
	{
		tunable:               "wiredTigerConcurrentReadTransactions",
		tunableKind:           "setParameter",
		topologyApplicability: []string{"standalone", "replica_set_primary", "replica_set_secondary", "config_server", "pali_primary", "pali_standby"},
		directionWhenObserved: "increase",
		requiredFindingKinds:  []string{"capacity_pressure", "near_ceiling"},
		requiredMetricNames:   []string{"serverStatus.wiredTiger.concurrentTransactions.read.out"},
		verificationRecipe:    "ftdc_parser <next-capture> --auto --filter 'concurrentTransactions|opLatencies'",
		verificationHelped:    "capacity_pressure clears on read tickets; opLatencies.reads p99 improves.",
		verificationNotHelped: "CPU saturates (host_pathology subkind=cpu_steal fires).",
		guardrails: []string{
			"On 7.0+ prefer storageEngineConcurrencyAdjustmentAlgorithm=throughputProbing.",
		},
	},
	{
		tunable:               "net.maxIncomingConnections",
		tunableKind:           "config",
		topologyApplicability: []string{"standalone", "replica_set_primary", "replica_set_secondary", "mongos", "config_server", "pali_primary", "pali_standby"},
		directionWhenObserved: "increase",
		requiredFindingKinds:  []string{"near_ceiling"},
		requiredMetricNames:   []string{"serverStatus.connections.current", "serverStatus.connections.available"},
		verificationRecipe:    "ftdc_parser <next-capture> --auto --filter 'connections'",
		verificationHelped:    "connections.rejected delta is 0 and connect-time p99 stable.",
		verificationNotHelped: "ulimit/file-descriptor ceiling hit elsewhere; check kernel logs.",
		guardrails: []string{
			"Confirm OS ulimits permit the new value; mongod refuses values above the FD limit.",
			"Each connection holds a TCP socket + thread-stack memory (~1 MB).",
		},
	},
	{
		tunable:               "storage.wiredTiger.engineConfig.cacheSizeGB",
		tunableKind:           "config",
		topologyApplicability: []string{"standalone", "replica_set_primary", "replica_set_secondary", "config_server", "pali_primary"},
		directionWhenObserved: "increase",
		requiredFindingKinds:  []string{"capacity_pressure", "eviction_divergence"},
		requiredMetricNames:   []string{"serverStatus.wiredTiger.cache.bytes currently in the cache", "serverStatus.wiredTiger.cache.maximum bytes configured"},
		verificationRecipe:    "ftdc_parser <next-capture> --auto --filter 'wiredTiger\\.cache'",
		verificationHelped:    "capacity_pressure on cache.bytes clears; pages evicted by application threads rate drops ≥50%.",
		verificationNotHelped: "Cache pressure shifts to dirty bytes; eviction_divergence persists.",
		guardrails: []string{
			"Confirm host RAM headroom: cache size + OS page cache + non-WT mongod memory must fit.",
			"Containerized mongod: respect cgroup memory limits; OOM-killer can be sudden.",
		},
	},
	{
		tunable:               "storageEngineConcurrencyAdjustmentAlgorithm",
		tunableKind:           "setParameter",
		topologyApplicability: []string{"standalone", "replica_set_primary", "replica_set_secondary", "config_server", "pali_primary", "pali_standby"},
		directionWhenObserved: "investigate",
		// Variance shift on the ticket-out gauge is the structural cue; we
		// rely on metric_context to tell the analyzer which side (read/write).
		requiredFindingKinds:  []string{"variance_shift"},
		requiredMetricNames:   []string{"serverStatus.wiredTiger.concurrentTransactions.write.out"},
		verificationRecipe:    "ftdc_parser <next-capture> --auto --filter 'concurrentTransactions'",
		verificationHelped:    "variance_shift no longer fires on the ticket gauge; tickets-out variance stabilizes.",
		verificationNotHelped: "Oscillation continues — switching algorithm did not help; consider manual tickets.",
		guardrails: []string{
			"Requires mongod 7.0+.",
			"Switching algorithm at runtime is reversible.",
		},
	},
	{
		tunable:               "paliPageLogMaterializationThreadCount",
		tunableKind:           "setParameter",
		topologyApplicability: []string{"pali_standby"},
		directionWhenObserved: "increase",
		requiredFindingKinds:  []string{"materialization_lag"},
		requiredMetricNames:   []string{"serverStatus.metrics.disagg.phylog.committedToMaterializedLagMillis"},
		verificationRecipe:    "ftdc_parser <next-capture> --auto --filter 'disagg.phylog'",
		verificationHelped:    "materialization_lag p95 drops below 1000ms.",
		verificationNotHelped: "Primary upstream is the cause — check primary disagg.logServer rate.",
		guardrails: []string{
			"More applier threads compete with user workload on CPU-bound hosts.",
			"Confirm CPU has headroom (host_pathology subkind=cpu_steal NOT firing).",
		},
	},
	{
		tunable:               "paliPageServerLocalCacheSizeBytes",
		tunableKind:           "setParameter",
		topologyApplicability: []string{"pali_primary", "pali_standby"},
		directionWhenObserved: "increase",
		requiredFindingKinds:  []string{"io_routing"},
		requiredMetricNames:   []string{"serverStatus.metrics.disagg.pageServer.local.success", "serverStatus.metrics.disagg.pageServer.nonLocal.success"},
		verificationRecipe:    "ftdc_parser <next-capture> --auto --filter 'disagg.pageServer'",
		verificationHelped:    "io_routing non-local fraction drops below 0.05.",
		verificationNotHelped: "Hot working set genuinely larger than RAM; or workload is range scans across all collections.",
		guardrails: []string{
			"Local cache size is bounded by host RAM. Coordinate with cacheSizeGB.",
		},
	},
	{
		tunable:               "tcmallocReleaseRate",
		tunableKind:           "setParameter",
		topologyApplicability: []string{"standalone", "replica_set_primary", "replica_set_secondary", "mongos", "config_server", "pali_primary", "pali_standby"},
		directionWhenObserved: "increase",
		requiredFindingKinds:  []string{"memory_triplet"},
		requiredMetricNames:   []string{"serverStatus.tcmalloc.tcmalloc.pageheap_free_bytes", "serverStatus.mem.resident"},
		verificationRecipe:    "ftdc_parser <next-capture> --auto --filter 'tcmalloc|mem'",
		verificationHelped:    "pageheap_free_bytes p95 drops; mem.resident stable.",
		verificationNotHelped: "CPU cost outweighs memory benefit; tcmalloc.* CPU samples rise.",
		guardrails: []string{
			"Higher release rate = more decommit syscall traffic. Watch for sys-CPU rise.",
		},
	},
}

// antiRecommendationCatalog: structural gates that block tunable
// emission. Each gate names the codes it triggers and the tunables it
// blocks (empty = block all tunables this run).
var antiRecommendationCatalog = []antiRecommendationGate{
	{
		code:      "cache_at_ceiling_with_ticket_pressure",
		rationale: "Cache pressure AND ticket pressure both firing — increasing tickets amplifies cache fill. Fix the cache first.",
		// Hard to express as a single 'present' set without false positives;
		// the gate fires only when BOTH capacity_pressure (cache) AND
		// near_ceiling (tickets) are present.
		requiredKindsPresent: []string{"capacity_pressure", "near_ceiling"},
		requiredKindsAbsent:  []string{},
		tunablesBlocked:      []string{"wiredTigerConcurrentWriteTransactions", "wiredTigerConcurrentReadTransactions"},
	},
	{
		code:      "signal_below_noise_floor",
		rationale: "All supporting evidence has |cohens_d| < 0.2 — effect is in the noise.",
		requiredKindsPresent: []string{},
		requiredKindsAbsent:  []string{},
		tunablesBlocked:      []string{},
	},
	{
		code:      "capture_too_short",
		rationale: "Samples evaluated < 600 (~10 minutes of FTDC). Base rate of the signal is unstable.",
		requiredKindsPresent: []string{},
		requiredKindsAbsent:  []string{},
		tunablesBlocked:      []string{},
	},
	{
		code:      "pali_cluster_a_known_bug",
		rationale: "Co-occurrence of seal_event + role_asymmetry matches the PALI Cluster A bug pattern (see project_pali_cluster_a.md). Refuse to recommend tunables; this is a code fix, not a config change.",
		requiredKindsPresent: []string{KindSealEvent, "role_asymmetry"},
		requiredKindsAbsent:  []string{},
		tunablesBlocked:      []string{},
	},
}
