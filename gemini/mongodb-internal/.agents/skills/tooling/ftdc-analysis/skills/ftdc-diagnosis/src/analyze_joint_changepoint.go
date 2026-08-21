// analyze_joint_changepoint.go — P3.2 joint-changepoint clustering.
//
// Stage 3 ED-PELT emits a changepoint Finding per metric, independently.
// Five metrics breaking at the same sample index is much stronger
// evidence of a true regime shift than five independent univariate
// breakpoints scattered across the capture. This stage post-processes
// the changepoint Findings already in the collector and emits a
// `joint_changepoint` Finding wherever ≥5 breakpoints align within a
// ±3-sample window.
//
// The synthetic_event Finding clusters co-occurring detector firings
// across kinds (gap + variance_shift + spike + …). joint_changepoint
// is narrower and stronger: it clusters within a SINGLE kind, anchored
// on what's structurally the same event (a multi-metric regime shift).
package main

import (
	"math"
	"sort"
	"strconv"
)

// jointCPClusterSpan: same ±3-sample width as window_suggestion. Two
// changepoints whose sample indices differ by more than this are NOT
// considered part of the same regime shift.
const jointCPClusterSpan = 3

// jointCPMinCluster: minimum number of independent metric changepoints
// in a cluster before we emit. Stage 3's univariate detector already
// catches single-metric breakpoints — we only emit when there's joint
// evidence.
const jointCPMinCluster = 5

// emitJointChangepoint reads changepoint Findings from the collector,
// clusters them by sample_index, and emits joint_changepoint Findings
// for clusters meeting the size threshold. Pure post-stage composer.
func emitJointChangepoint(c *findingCollector, host string) (emitted int, skipped int) {
	type cpEntry struct {
		fid       string
		idx       int
		metric    string
		direction string
		cohensD   float64
	}
	var entries []cpEntry
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		if k != KindChangepoint {
			continue
		}
		hf, _ := f["host"].(string)
		if hf != "" && hf != host {
			continue
		}
		idx, ok := toIntField(f["sample_index"])
		if !ok || idx < 0 {
			continue
		}
		fid, _ := f["id"].(string)
		metric, _ := f["metric"].(string)
		cd, _ := f["cohens_d"].(float64)
		// cohens_d is stored as absFloat(effect) — always ≥ 0. Delegate
		// direction to inferFindingDirection, which uses baseline vs stressed
		// means and returns "" when both are equal or missing.
		dir := inferFindingDirection(f)
		entries = append(entries, cpEntry{
			fid:       fid,
			idx:       idx,
			metric:    metric,
			direction: dir,
			cohensD:   cd,
		})
	}
	if len(entries) < jointCPMinCluster {
		skipped++
		return 0, skipped
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].idx != entries[j].idx {
			return entries[i].idx < entries[j].idx
		}
		return entries[i].metric < entries[j].metric
	})

	// Greedy single-pass cluster sweep, same algorithm as
	// window_suggestion. We avoid sharing code so each composer stays
	// self-contained — different thresholds, different outputs.
	type cluster struct {
		minIdx, maxIdx int
		members        []cpEntry
	}
	var clusters []cluster
	var cur cluster
	flush := func() {
		if len(cur.members) >= jointCPMinCluster {
			clusters = append(clusters, cur)
		}
		cur = cluster{}
	}
	for _, e := range entries {
		if len(cur.members) == 0 {
			cur.minIdx, cur.maxIdx = e.idx, e.idx
			cur.members = []cpEntry{e}
			continue
		}
		if e.idx-cur.maxIdx <= jointCPClusterSpan {
			cur.maxIdx = e.idx
			cur.members = append(cur.members, e)
		} else {
			flush()
			cur.minIdx, cur.maxIdx = e.idx, e.idx
			cur.members = []cpEntry{e}
		}
	}
	flush()

	if len(clusters) == 0 {
		skipped++
		return 0, skipped
	}

	for _, cl := range clusters {
		// Build per-metric direction vector + supporting Finding IDs.
		// Order by metric name for determinism.
		sort.Slice(cl.members, func(i, j int) bool {
			return cl.members[i].metric < cl.members[j].metric
		})
		dirVec := make([]map[string]string, 0, len(cl.members))
		supportingIDs := make([]string, 0, len(cl.members))
		nUp, nDown := 0, 0
		var sumAbsD float64
		for _, m := range cl.members {
			dirVec = append(dirVec, map[string]string{
				"metric":    m.metric,
				"direction": m.direction,
			})
			if m.fid != "" {
				supportingIDs = append(supportingIDs, m.fid)
			}
			switch m.direction {
			case "up":
				nUp++
			case "down":
				nDown++
			}
			sumAbsD += math.Abs(m.cohensD)
		}
		idKey := "joint_cp|" + strconv.Itoa(cl.minIdx)
		c.add(newFinding("joint_changepoint", host, idKey, map[string]interface{}{
			"min_sample_index":       cl.minIdx,
			"max_sample_index":       cl.maxIdx,
			"cluster_size":           len(cl.members),
			"n_up":                   nUp,
			"n_down":                 nDown,
			"sum_abs_cohens_d":       roundFloat(sumAbsD, 4),
			"direction_vector":       dirVec,
			"supporting_finding_ids": supportingIDs,
		}))
		emitted++
	}
	return emitted, skipped
}

