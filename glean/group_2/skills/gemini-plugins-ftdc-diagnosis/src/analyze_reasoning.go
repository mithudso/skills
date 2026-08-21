// analyze_reasoning.go — cross-stage reasoning over already-emitted
// Findings. These stages run AFTER all detector stages so they can
// compose evidence across the pipeline.
//
// Deliverables:
//   - `kind:"workload_spec"`            — numeric workload profile derived
//                                         from opcounters + connection
//                                         churn + lock metrics. Replaces
//                                         the speculative "workload
//                                         classification" item with a
//                                         pure-numeric structural
//                                         summary.
//   - `kind:"precursor_candidates"`     — for each changepoint Finding,
//                                         ranks the metrics whose
//                                         changepoint/spike preceded it
//                                         within a configurable window.
//                                         Pure precedence + correlation,
//                                         no causality claim.
//
// Bias-control invariants preserved: workload_spec emits ONLY numeric
// rates and ratios (no read_pct/write_pct categorical label);
// precursor_candidates emits ordered precedence pairs with correlation
// magnitudes — no `cause`, no severity, no role.

package main

import (
	"sort"
	"time"
)

// emitWorkloadSpec builds a numeric workload profile per host. Fields
// echo the raw rates the analyzer would otherwise have to derive
// itself; this stage just packages them. No "category" label is
// emitted — the analyzer skill interprets the numbers.
func emitWorkloadSpec(c *findingCollector, host string, store *seriesStore) (emitted, skipped int) {
	if store.nSamples() < 2 {
		return 0, 1
	}
	// Cumulative-counter "rate over capture" helper.
	rate := func(metric string) (float64, bool) {
		id, ok := store.nameToID[metric]
		if !ok {
			return 0, false
		}
		col := store.columns[id]
		if col == nil || col.Len() < 2 {
			return 0, false
		}
		vals := col.ReadAllInt64(nil)
		first := vals[0]
		last := vals[len(vals)-1]
		// Duration in seconds (rough).
		dur := store.timestamps[len(store.timestamps)-1].Sub(store.timestamps[0]).Seconds()
		if dur <= 0 {
			return 0, false
		}
		if last < first {
			// Counter reset (process restart during capture); unavailable.
			return 0, false
		}
		return float64(last-first) / dur, true
	}
	// Gauge mean helper.
	gaugeMean := func(metric string) (float64, bool) {
		id, ok := store.nameToID[metric]
		if !ok {
			return 0, false
		}
		col := store.columns[id]
		if col == nil || col.Len() < 2 {
			return 0, false
		}
		vals := col.ReadAllInt64(nil)
		var sum float64
		for _, v := range vals {
			sum += float64(v)
		}
		return sum / float64(len(vals)), true
	}

	fields := map[string]interface{}{
		"host":         host,
		"sample_count": store.nSamples(),
	}

	// Op rates (per second). All entries are emitted as raw rates;
	// the analyzer derives ratios.
	for _, m := range []struct {
		key    string
		metric string
	}{
		{"insert_rate_ps", "serverStatus.opcounters.insert"},
		{"query_rate_ps", "serverStatus.opcounters.query"},
		{"update_rate_ps", "serverStatus.opcounters.update"},
		{"delete_rate_ps", "serverStatus.opcounters.delete"},
		{"command_rate_ps", "serverStatus.opcounters.command"},
		{"getmore_rate_ps", "serverStatus.opcounters.getmore"},
	} {
		if r, ok := rate(m.metric); ok {
			fields[m.key] = roundFloat(r, 2)
		}
	}

	// Connection churn rate.
	if r, ok := rate("serverStatus.connections.totalCreated"); ok {
		fields["conn_churn_rate_ps"] = roundFloat(r, 2)
	}
	// Concurrent client snapshot.
	if m, ok := gaugeMean("serverStatus.connections.current"); ok {
		fields["conn_current_mean"] = roundFloat(m, 1)
	}

	// Replication op rates (secondary writes mirror primary user
	// writes — surface separately so the analyzer can tell primary
	// from secondary by comparing).
	if r, ok := rate("serverStatus.opcountersRepl.insert"); ok {
		fields["repl_insert_rate_ps"] = roundFloat(r, 2)
	}
	if r, ok := rate("serverStatus.opcountersRepl.update"); ok {
		fields["repl_update_rate_ps"] = roundFloat(r, 2)
	}

	// Don't emit a Finding if every rate is missing.
	hasAny := false
	for k := range fields {
		if k != "host" && k != "sample_count" {
			hasAny = true
			break
		}
	}
	if !hasAny {
		return 0, 1
	}
	c.add(newFinding(KindWorkloadSpec, host, "global", fields))
	return 1, 0
}

// emitPrecursorCandidates walks the changepoint Findings already in
// the collector and, for each, finds the top-K Findings on OTHER
// metrics that fired EARLIER (precedence) within a configurable look-
// back window. Emits one `precursor_candidates` Finding per target
// changepoint with the ranked precursor list.
//
// Pure precedence + magnitude. NO causality claim — the analyzer can
// build a hypothesis from the ranking.
func emitPrecursorCandidates(c *findingCollector, host string) (emitted, skipped int) {
	type findingTime struct {
		f      Finding
		metric string
		t      time.Time
		mag    float64
	}
	var withTimes []findingTime
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		switch k {
		case KindChangepoint, KindSpike, KindVarianceShift:
			// eligible
		default:
			continue
		}
		// Multi-host guard: only include findings from the same host to
		// prevent cross-host temporal coincidences from producing false
		// precursor relationships.
		if fHost, _ := f["host"].(string); fHost != host {
			continue
		}
		atStr, _ := f["at"].(string)
		metric, _ := f["metric"].(string)
		if atStr == "" || metric == "" {
			continue
		}
		t, err := time.Parse(tsFormat, atStr)
		if err != nil {
			continue
		}
		withTimes = append(withTimes, findingTime{
			f: f, metric: metric, t: t, mag: findingMagnitude(f),
		})
	}
	if len(withTimes) == 0 {
		return 0, 1
	}
	// Sort by timestamp ascending so we can binary-search the
	// look-back window.
	sort.Slice(withTimes, func(i, j int) bool {
		return withTimes[i].t.Before(withTimes[j].t)
	})
	const lookbackSec = 300
	const topK = 5
	lookback := time.Duration(lookbackSec) * time.Second

	for _, target := range withTimes {
		// Only emit precursors for `changepoint` targets — the others
		// (spike, variance_shift) are too noisy as anchors.
		if k, _ := target.f["kind"].(string); k != KindChangepoint {
			continue
		}
		// Find candidates in (target.t - lookback, target.t).
		var candidates []findingTime
		for _, cand := range withTimes {
			if cand.metric == target.metric {
				continue // skip same metric
			}
			if !cand.t.Before(target.t) {
				continue
			}
			if target.t.Sub(cand.t) > lookback {
				continue
			}
			candidates = append(candidates, cand)
		}
		if len(candidates) == 0 {
			continue
		}
		// Rank by magnitude desc; ties broken by closest-in-time, then by
		// metric name (deterministic). Tier B fix: when two precursors
		// share both magnitude and timestamp (e.g. paired memory metrics
		// like nr_anon_transparent_hugepages vs AnonHugePages_kb), the
		// metric-name tiebreaker preserves byte-equality across --jobs.
		sort.SliceStable(candidates, func(i, j int) bool {
			if candidates[i].mag != candidates[j].mag {
				return candidates[i].mag > candidates[j].mag
			}
			if !candidates[i].t.Equal(candidates[j].t) {
				return candidates[i].t.After(candidates[j].t)
			}
			return candidates[i].metric < candidates[j].metric
		})
		if len(candidates) > topK {
			candidates = candidates[:topK]
		}
		precursors := make([]interface{}, 0, len(candidates))
		for _, p := range candidates {
			lagSec := target.t.Sub(p.t).Seconds()
			fid, _ := p.f["id"].(string)
			pkind, _ := p.f["kind"].(string)
			precursors = append(precursors, map[string]interface{}{
				"finding_id":   fid,
				"kind":         pkind,
				"metric":       p.metric,
				"magnitude":    roundFloat(p.mag, 4),
				"lead_seconds": roundFloat(lagSec, 1),
			})
		}
		targetID, _ := target.f["id"].(string)
		c.add(newFinding("precursor_candidates", host, targetID, map[string]interface{}{
			"target_finding_id": targetID,
			"target_metric":     target.metric,
			"target_at":         target.t.UTC().Format(tsFormat),
			"lookback_seconds":  lookbackSec,
			"precursors":        precursors,
		}))
		emitted++
	}
	if emitted == 0 {
		skipped++
	}
	return emitted, skipped
}
