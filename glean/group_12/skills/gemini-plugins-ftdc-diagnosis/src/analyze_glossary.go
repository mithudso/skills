// analyze_glossary.go - Foundation: metric_glossary + metric_context
// Findings, plus the token-budget gate and metric_score_table opt-in.
//
// Goals (from /home/ubuntu/.claude/plans/how-to-make-ftdc-parser-graceful-sky.md):
//
//   * `kind:"metric_glossary"`: single Finding at end of run, keyed only on
//     metrics actually cited by other Findings. Cuts cost from O(findings)
//     to O(distinct_metrics_cited) — typically 30-80 entries vs. duplicating
//     per-finding docs across ~2000 Findings.
//
//   * `kind:"metric_context"`: per-capture context the analyzer skill cannot
//     know — MongoDB version, storage engine, tcmalloc variant, host count,
//     architectural shape. Extracted from FTDC metadata documents.
//
//   * `--auto-budget-tokens N` gate: after all stages emit, if cumulative
//     output exceeds the budget, switch to summary mode (per-kind counts +
//     top-3 per kind + a `kind:"truncated"` Finding listing elided kinds).
//
//   * `--auto-emit-score-table` opt-in: by default, the ~190K-token
//     metric_score_table Finding is filtered out at emit time.
//
// All four mechanisms run AFTER stages emit Findings; they're post-pass
// transforms on the collector's Findings list.

package main

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"sort"
	"strings"
)

// metricDocs returns the one-line description (and inferred unit category)
// for a given canonical FTDC metric path. Returns (doc, ok=false) when the
// metric is not in the catalog. Catalog is intentionally narrow: ~80 of the
// most-cited metrics across replication, storage, query, network, host, and
// PALI subsystems. Maintained inline rather than YAML so the build pipeline
// has no external file dependency; future expansion can move it to YAML
// when the catalog crosses ~300 entries.
//
// Naming convention: keys are exact FTDC paths (post-flattening). The doc
// text is one short sentence describing what the metric measures and, where
// useful, the consumer-facing implication.
func metricDocs(name string) (string, bool) {
	d, ok := metricDocsCatalog[name]
	return d, ok
}

// metricDocsCatalog: structural descriptions of metric semantics ONLY.
//
// Bias-control rule (per ftdc-parser/SKILL.md:161-165): no causes, no
// severity, no "bottleneck"/"stall"/"leak suspect"/"pressure"/etc.
// language. Each entry states what the metric measures (gauge/counter,
// units, accumulation semantics) and the structural pairing for ratios
// where applicable. Interpretive language lives in the analyzer skill,
// not in parser-emitted Findings.
//
// Tier C.3: the catalog map literal moved to metric_docs.go so editing
// documentation never touches detector logic. The variable is declared
// there; consumers in this file just reference `metricDocsCatalog`.

// The map literal previously here is now in metric_docs.go.

// emitMetricContext writes ONE Finding per capture summarizing pre-capture
// context the analyzer skill cannot derive from FTDC values alone:
// MongoDB version, storage engine, tcmalloc variant, host count, replica-set
// vs sharded vs standalone, etc. Extracted from FTDC metadata documents
// (type-0 buildInfo, hostInfo, getCmdLineOpts).
func emitMetricContext(c *findingCollector, host string, store *seriesStore, metadata []MetadataDoc) {
	fields := map[string]interface{}{
		"host":         host,
		"metric_count": len(store.names),
	}
	// MongoDB version + storage engine from buildInfo / getCmdLineOpts.
	for _, m := range metadata {
		if m.DocType != 0 {
			continue
		}
		// Best-effort: scan known buildInfo / hostInfo key paths.
		walkMetadataMap(m.Doc, func(key string, val interface{}) {
			switch key {
			case "version":
				if s, ok := val.(string); ok {
					fields["mongo_version"] = s
				}
			case "storageEngine.name":
				if s, ok := val.(string); ok {
					fields["storage_engine"] = s
				}
			case "system.numCores":
				if v, ok := toInt(val); ok {
					fields["host_num_cores"] = v
				}
			case "system.memSizeMB":
				if v, ok := toInt(val); ok {
					fields["host_memory_mb"] = v
				}
			case "system.cpuArch":
				if s, ok := val.(string); ok {
					fields["host_cpu_arch"] = s
				}
			case "os.name":
				if s, ok := val.(string); ok {
					fields["host_os"] = s
				}
			case "os.version":
				if s, ok := val.(string); ok {
					fields["host_os_version"] = s
				}
			}
		})
	}
	// Topology inference from metric names (already in the store).
	topology := inferTopology(store.names)
	if topology != "" {
		fields["topology"] = topology
	}
	// tcmalloc variant: presence of `tcmalloc.tcmalloc.metadata_bytes`
	// indicates google-tcmalloc (gperftools lacks this path). Used by
	// analyzer to handle the 7.0+ semantic shift.
	if _, ok := store.nameToID["serverStatus.tcmalloc.tcmalloc.metadata_bytes"]; ok {
		fields["tcmalloc_variant"] = "google"
	} else if _, ok := store.nameToID["serverStatus.tcmalloc.generic.heap_size"]; ok {
		fields["tcmalloc_variant"] = "gperftools_or_compat"
	}
	c.add(newFinding("metric_context", host, "global", fields))
}

// inferTopology returns "standalone" / "replica_set" / "sharded_router" /
// "sharded_member" based on metric name presence. Pure structural test,
// no name-based recognizers beyond what's already encoded in the FTDC
// schema fingerprint.
func inferTopology(names []string) string {
	if len(names) == 0 {
		return ""
	}
	hasSharding := false
	hasRepl := false
	hasWT := false
	hasMongos := true // mongos lacks WT and repl; assume true until contradicted
	for _, n := range names {
		switch {
		case strings.HasPrefix(n, "serverStatus.shardingStatistics."):
			hasSharding = true
		case strings.HasPrefix(n, "serverStatus.repl."):
			hasRepl = true
			hasMongos = false
		case strings.HasPrefix(n, "serverStatus.wiredTiger."):
			hasWT = true
			hasMongos = false
		}
	}
	switch {
	case hasMongos && !hasWT && !hasRepl:
		return "sharded_router"
	case hasSharding && hasWT:
		return "sharded_member"
	case hasRepl && hasWT:
		return "replica_set"
	case hasWT:
		return "standalone"
	}
	return ""
}

// emitMetricGlossary scans all already-emitted Findings for metric names
// and emits ONE `kind:"metric_glossary"` Finding with entries for the
// distinct metrics cited that have docs in the catalog. Token cost: O(K)
// where K = distinct cited metrics (typically 30-80), instead of
// O(findings × per-finding-docs).
func emitMetricGlossary(c *findingCollector, host string) {
	cited := map[string]bool{}
	for _, f := range c.items {
		for _, k := range []string{"metric", "leader", "follower"} {
			if v, ok := f[k]; ok {
				if s, ok := v.(string); ok && s != "" {
					cited[s] = true
				}
			}
		}
	}
	if len(cited) == 0 {
		return
	}
	names := make([]string, 0, len(cited))
	for n := range cited {
		names = append(names, n)
	}
	sort.Strings(names)
	entries := make([]interface{}, 0, len(names))
	// Tier C.4: track uncatalogued cited metrics so we can emit a
	// glossary_gap Finding listing them. Pre-C.4 the silent-continue
	// here hid the gap; growing the catalog became guesswork. Now the
	// gap surfaces in every run's output, capped at top-50 by citation
	// frequency.
	uncatalogued := make([]string, 0)
	for _, n := range names {
		doc, ok := metricDocs(n)
		if !ok {
			uncatalogued = append(uncatalogued, n)
			continue
		}
		entries = append(entries, map[string]interface{}{
			"metric": n,
			"docs":   doc,
		})
	}
	if len(entries) > 0 {
		c.add(newFinding("metric_glossary", host, "global", map[string]interface{}{
			"entries":     entries,
			"cited_count": len(names),
			"covered":     len(entries),
			"note":        "single-Finding glossary covering only metrics cited in this run; full catalog ~80 metrics",
		}))
	}
	if len(uncatalogued) > 0 {
		// Cap at 50 so a capture with thousands of cited-but-uncatalogued
		// metrics doesn't drown the output. Names are already sorted (parent
		// sort.Strings call above).
		const glossaryGapTopK = 50
		shown := uncatalogued
		if len(shown) > glossaryGapTopK {
			shown = shown[:glossaryGapTopK]
		}
		c.add(newFinding("glossary_gap", host, "global", map[string]interface{}{
			"uncatalogued_metrics": shown,
			"uncatalogued_count":   len(uncatalogued),
			"truncated":            len(uncatalogued) > glossaryGapTopK,
			"note":                 "metrics cited by other Findings that are not in metricDocsCatalog; growing the catalog from real captures is the intended use",
		}))
	}
}

// applyOutputGates removes Findings the user opted out of (e.g. score
// table) and applies the token-budget truncation. Called immediately
// before emitAll.
//
// Score-table gate: by default the ~190K-token metric_score_table is
// dropped from output. Opt-in via --auto-emit-score-table.
//
// Budget gate: when the cumulative byte count × 0.25 (≈ tokens) exceeds
// cfg.BudgetTokens, switch to summary mode — keep run_metadata,
// run_health, schema_fingerprint, metric_context, metric_glossary,
// changepoint_oversegmented, and the top-3-by-magnitude per remaining
// kind. Append a single `kind:"truncated"` Finding listing the elided
// kinds and counts.
func applyOutputGates(c *findingCollector, cfg AutoConfig) {
	// 1. Score-table gate.
	if !cfg.EmitScoreTable {
		filtered := c.items[:0]
		for _, f := range c.items {
			if k, _ := f["kind"].(string); k == KindMetricScoreTable {
				continue
			}
			filtered = append(filtered, f)
		}
		c.items = filtered
	}

	// 2. Budget gate.
	if cfg.BudgetTokens <= 0 {
		return
	}
	budgetBytes := cfg.BudgetTokens * 4
	totalBytes := 0
	for _, f := range c.items {
		if b, err := json.Marshal(map[string]interface{}(f)); err == nil {
			totalBytes += len(b) + 1
		}
	}
	if totalBytes <= budgetBytes {
		return
	}

	// Over budget: build summary-mode output. The keepAll set covers
	// every Finding kind that the analyzer skill cannot reconstruct
	// from the cumulative top-3 view of the truncated kinds:
	// - run_metadata, run_health, schema_fingerprint, metric_context,
	//   metric_glossary: per-capture context the analyzer reads first.
	// - changepoint_oversegmented: data-quality warning per metric.
	// - synthetic_event: overlay clustering that ties other Findings
	//   together. Losing it would mean N child Findings with no
	//   cross-reference signal.
	// - consistency_warning: data-quality flag on specific metrics;
	//   analyzer needs every one to demote dependent Findings.
	keepAll := map[string]bool{
		KindRunMetadata:              true,
		KindRunHealth:                true,
		KindSchemaFingerprint:        true,
		"metric_context":             true,
		"metric_glossary":            true,
		KindChangepointOversegmented: true,
		"synthetic_event":            true,
		"consistency_warning":        true,
		// Signal stage:
		"suspicious_smoothness": true, // data-quality flag; analyzer needs every instance
		"evidence_convergence":  true, // overlay-on-overlay; lose this, lose multi-detector signal
		// Health stage:
		"clock_skew":      true, // data-quality flag; affects every cross-host claim
		"missing_signal":  true, // data-quality flag; affects every absence-based claim
		"alert_threshold": true, // point events (election/restart) are rare and high-value
		// Domain stage — emits at most a handful per capture:
		KindTermChange:              true, // election event marker
		"oplog_window_low":           true, // rare per-capture
		"checkpoint_install_anomaly": true, // rare per-capture
		"host_pathology":             true, // host-level pathology overrides MongoDB-layer claims
		// Reasoning stage — analyzer reads these BEFORE forming hypotheses:
		"workload_spec":        true, // pre-context numeric profile
		"precursor_candidates": true, // cross-stage reasoning overlay
		// Workflow stage — features used across captures:
		"memory_triplet":         true, // joint memory shape — analyzer maps to fragmentation/leak
		"recurrence_fingerprint": true, // single-Finding hash; analyzer compares across runs
		// Cross-host pre-existing kinds (Final-acceptance review): one
		// per host, low cardinality; capping at top-3 silently drops
		// data on multi-host runs with ≥4 hosts.
		KindHostOnlyMetrics:        true,
		KindCrossHostRoleInferred:  true,
		KindHostMetricSetAsymmetry: true,
		// Tier B new kinds — all advisory; keepAll so the budget gate
		// doesn't silently truncate diagnostically useful Findings.
		KindTrendInconclusive:   true,
		"counter_wrap":          true,
		"counter_inversion":     true,
		"local_spike":           true,
		"glossary_gap":          true, // Tier C.4
		"regression":            true, // Tier D
		"baseline_diagnostic":   true, // Tier D
		"apply_saturation":      true, // Tier G.1
		"election_storm":        true, // Tier G.1
		"checkpoint_stall":      true, // Tier G.1
		"eviction_divergence":   true, // Tier G.1
		"capture_quality_score": true, // Tier G.3
		"chunk_decode_error":    true, // Tier A.2
		"diagnosis_candidate":   true, // Tier R.1
		"hypothesis_test":       true, // Tier R.2
		// Round 7 fix — previously missing from keepAll; under budget
		// pressure their top-3 cap silently truncated diagnostically
		// important kinds (e.g., a balancer thrash with 30+ migrations).
		"balancer_round_stall":        true,
		"chunk_migration_event":       true,
		"io_routing":                  true,
		"near_ceiling":                true,
		"periodic":                    true,
		"replication_buffer_pressure": true,
		"utilization":                 true,
		"capacity_pressure":           true,
		"materialization_lag":         true,
		"truncated":                   true, // self-evident: don't drop the truncation notice
		// P0.NEW-PALI-1: seal event Findings are structural anchors; always keep.
		KindSealEvent:          true,
		KindSealEventCandidate: true,
	}
	const topPerKind = 3
	byKind := map[string][]Finding{}
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		byKind[k] = append(byKind[k], f)
	}
	kept := make([]Finding, 0, len(c.items))
	elided := map[string]int{}
	for k, lst := range byKind {
		if keepAll[k] {
			kept = append(kept, lst...)
			continue
		}
		// Keep top-3 by absolute score / effect / r when available.
		sort.SliceStable(lst, func(i, j int) bool {
			return findingMagnitude(lst[i]) > findingMagnitude(lst[j])
		})
		if len(lst) > topPerKind {
			kept = append(kept, lst[:topPerKind]...)
			elided[k] = len(lst) - topPerKind
		} else {
			kept = append(kept, lst...)
		}
	}
	c.items = kept
	if len(elided) > 0 {
		summary := make(map[string]interface{}, len(elided)+3)
		summary["kinds_elided"] = elided
		summary["budget_tokens"] = cfg.BudgetTokens
		summary["pre_budget_bytes"] = totalBytes
		summary["note"] = "output truncated to --auto-budget-tokens; per-kind top-3 retained; full set available via --auto-budget-tokens 0"
		var hostStr string
		if len(kept) > 0 {
			if hostAny, ok := kept[0]["host"]; ok {
				hostStr, _ = hostAny.(string)
			}
		}
		c.items = append(c.items, newFinding("truncated", hostStr, "global", summary))
	}
}

// findingMagnitude returns a per-Finding score used to rank elidability
// under the budget gate. Higher = more important. The priority list
// covers magnitude-like fields across all Finding kinds — adding here
// when a new kind ships keeps the budget gate's top-K sort meaningful.
// Field order: most-magnitude-discriminating first.
func findingMagnitude(f Finding) float64 {
	for _, k := range []string{
		"effect_size",         // changepoint
		"score",               // mover
		"r",                   // correlation
		"welch_t",             // mover
		"slope_t",             // monotonic_trend
		"z_score",             // spike
		"cusum_peak",          // variance_shift
		"duration_seconds",    // gap
		"longest_run_seconds", // capacity_pressure
		"last_ratio",          // near_ceiling
		"pct_at_ceiling",      // capacity_pressure
		"snr",                 // periodic
	} {
		if v, ok := f[k]; ok {
			if x, ok := v.(float64); ok {
				return absFloat(x)
			}
		}
	}
	// synthetic_event: rank by member count (cluster size).
	// consistency_warning: no natural magnitude — keep static
	// non-zero rank so the sort is stable.
	if k, _ := f["kind"].(string); k == "synthetic_event" {
		if mc, ok := f["member_count"].(int); ok {
			return float64(mc)
		}
	}
	return 0
}

// walkMetadataMap traverses a BSON document (in bson.D form post-Unmarshal)
// emitting (dotted-path, value) pairs to fn. Used by emitMetricContext to
// extract structured fields from buildInfo/hostInfo without hand-coding
// every path. Non-leaf containers (nested bson.D / map / array) are
// recursed; scalars are emitted.
func walkMetadataMap(doc interface{}, fn func(string, interface{})) {
	walkMetadataMapRec(doc, "", fn)
}
func walkMetadataMapRec(node interface{}, prefix string, fn func(string, interface{})) {
	switch v := node.(type) {
	case bson.D:
		for _, e := range v {
			key := e.Key
			if prefix != "" {
				key = prefix + "." + key
			}
			walkMetadataMapRec(e.Value, key, fn)
		}
	case bson.M:
		for k, val := range v {
			key := k
			if prefix != "" {
				key = prefix + "." + k
			}
			walkMetadataMapRec(val, key, fn)
		}
	case map[string]interface{}:
		for k, val := range v {
			key := k
			if prefix != "" {
				key = prefix + "." + k
			}
			walkMetadataMapRec(val, key, fn)
		}
	case bson.A:
		// Arrays recurse with index-suffixed keys (versionArray.0, .1).
		// Without this case, arrays would fall through to default and
		// emit as a single leaf, masking nested D/M elements inside
		// (e.g. config-array buildInfo fields).
		for i, val := range v {
			walkMetadataMapRec(val, prefix+"."+intToString(i), fn)
		}
	default:
		if prefix != "" {
			fn(prefix, node)
		}
	}
}

// toInt extracts an int from common BSON numeric types.
func toInt(v interface{}) (int, bool) {
	switch x := v.(type) {
	case int32:
		return int(x), true
	case int64:
		return int(x), true
	case int:
		return x, true
	case float64:
		return int(x), true
	}
	return 0, false
}

// intToString avoids importing strconv into this file.
func intToString(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = '0' + byte(i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
