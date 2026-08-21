// Tier R.2: emit kind:"hypothesis_test" Findings — copy-pasteable recipes
// that confirm or rule out diagnosis_candidate Findings.
//
// Runs AFTER analyze_diagnosis.go has emitted diagnosis_candidate Findings.
// For each candidate with confidence_rank < 0.9 AND discriminating_kinds_needed
// non-empty, emits one or more hypothesis_test Findings.
//
// Recipes come from a pre-written catalog; the parser NEVER invents
// recipes. Recipe strings are parametrized by the matched metric names
// from supporting_finding_ids (looked up via the collector by id).

package main

import (
	"fmt"
	"strings"
)

// hypothesisRecipe defines the test for a particular discriminating kind.
// The recipe is a copy-pasteable parser invocation; the expected_outcome
// strings describe what the user should look for in the next capture.
type hypothesisRecipe struct {
	discriminatingKind         string
	recipeTemplate             string // sprintf'd with %s metric-glob
	metricGlob                 string // glob pattern substituted into the recipe's --filter
	expectedOutcomeIfConfirmed string
	expectedOutcomeIfRuledOut  string
}

var hypothesisRecipes = map[string][]hypothesisRecipe{
	"primary_pali_log_pressure": {
		{
			discriminatingKind:         "counter_wrap",
			recipeTemplate:             "ftdc_parser /path/to/next/capture/ --auto --filter '%s'",
			metricGlob:                 "disagg.phylog.*",
			expectedOutcomeIfConfirmed: "counter_wrap Finding on phylog rate, OR materialization_lag with higher confidence_rank",
			expectedOutcomeIfRuledOut:  "no counter_wrap; phylog counters flat or slowly increasing",
		},
		{
			discriminatingKind:         "regression",
			recipeTemplate:             "ftdc_parser --auto --auto-baseline-file ~/.local/share/ftdc-parser/baseline.jsonl /path/to/next/capture/",
			metricGlob:                 "",
			expectedOutcomeIfConfirmed: "regression Finding on phylog or memory metrics; pattern is a drift, not a steady-state",
			expectedOutcomeIfRuledOut:  "no regression Finding; current capture matches historical baseline; pattern is steady-state",
		},
	},
	"primary_step_down": {
		{
			discriminatingKind:         "alert_threshold",
			recipeTemplate:             "ftdc_parser /path/to/next/capture/ --auto --filter '%s'",
			metricGlob:                 "repl.*",
			expectedOutcomeIfConfirmed: "alert_threshold:election Finding around the term_change timestamp",
			expectedOutcomeIfRuledOut:  "no election-class alert; step-down was administrative, not failover",
		},
	},
	"secondary_apply_lag": {
		{
			discriminatingKind:         "checkpoint_stall",
			recipeTemplate:             "ftdc_parser /path/to/next/capture/ --auto --filter '%s'",
			metricGlob:                 "wiredTiger.transaction.transaction checkpoint*",
			expectedOutcomeIfConfirmed: "checkpoint_stall Finding co-occurring; apply backlog is bound by storage write throughput",
			expectedOutcomeIfRuledOut:  "no checkpoint_stall; apply backlog is CPU- or oplog-fetch-bound, not write-IO-bound",
		},
	},
	"wt_eviction_pressure": {
		{
			discriminatingKind:         "counter_wrap",
			recipeTemplate:             "ftdc_parser /path/to/next/capture/ --auto --filter '%s'",
			metricGlob:                 "wiredTiger.cache.*",
			expectedOutcomeIfConfirmed: "counter_wrap on a uint32 cache counter (unlikely; modern WT is int64) — pattern is then artifactual",
			expectedOutcomeIfRuledOut:  "no counter_wrap; eviction-divergence ratio reflects real cache-fill faster than eviction can drain",
		},
		{
			discriminatingKind:         "synthetic_event",
			recipeTemplate:             "ftdc_parser --auto --auto-around-event TIMESTAMP --auto-window 30m /path/to/capture/",
			metricGlob:                 "",
			expectedOutcomeIfConfirmed: "synthetic_event clustering around the eviction-divergence onset; multi-metric pressure, not isolated",
			expectedOutcomeIfRuledOut:  "no co-occurring multi-metric event; eviction pressure is isolated to the cache layer",
		},
	},
	"host_resource_starvation": {
		{
			discriminatingKind:         "regression",
			recipeTemplate:             "ftdc_parser --auto --auto-baseline-file ~/.local/share/ftdc-parser/baseline.jsonl /path/to/next/capture/",
			metricGlob:                 "",
			expectedOutcomeIfConfirmed: "regression Findings on memory or CPU metrics confirm trend versus historical baseline",
			expectedOutcomeIfRuledOut:  "no regression; resource consumption is steady-state, not drifting",
		},
	},
	"sharded_balancer_thrash": {
		{
			discriminatingKind:         "regression",
			recipeTemplate:             "ftdc_parser --auto --auto-baseline-file ~/.local/share/ftdc-parser/baseline.jsonl /path/to/next/capture/",
			metricGlob:                 "",
			expectedOutcomeIfConfirmed: "regression on chunk_migration counters confirms accelerating balancer activity",
			expectedOutcomeIfRuledOut:  "no regression; balancer activity is at steady-state for this cluster",
		},
	},
}

func emitHypothesisTests(c *findingCollector, host string) (emitted, skipped int) {
	// Walk diagnosis_candidate Findings.
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		if k != "diagnosis_candidate" {
			continue
		}
		pname, _ := f["pattern_name"].(string)
		fid, _ := f["id"].(string)
		score, _ := f["confidence_rank"].(float64)
		// discriminating_kinds_needed is stored as []string by emitDiagnosisCandidates.
		// asserting to []interface{} always fails (Go type system is exact).
		needed, _ := f["discriminating_kinds_needed"].([]string)
		if pname == "" || score >= 0.9 || len(needed) == 0 {
			skipped++
			continue
		}
		recipes, ok := hypothesisRecipes[pname]
		if !ok {
			skipped++
			continue
		}
		neededSet := map[string]bool{}
		for _, s := range needed {
			if s != "" {
				neededSet[s] = true
			}
		}
		for _, r := range recipes {
			if !neededSet[r.discriminatingKind] {
				continue
			}
			recipe := r.recipeTemplate
			if strings.Contains(recipe, "%s") {
				recipe = fmt.Sprintf(recipe, r.metricGlob)
			}
			c.add(newFinding("hypothesis_test", host, pname+"|"+r.discriminatingKind, map[string]interface{}{
				"parent_finding_id":             fid,
				"parent_pattern_name":           pname,
				"discriminating_kind":           r.discriminatingKind,
				"recipe":                        recipe,
				"metric_glob":                   r.metricGlob,
				"expected_outcome_if_confirmed": r.expectedOutcomeIfConfirmed,
				"expected_outcome_if_ruled_out": r.expectedOutcomeIfRuledOut,
			}))
			emitted++
		}
	}
	return emitted, skipped
}
