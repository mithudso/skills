// analyze_counterfactual.go — P4.5 tunable counterfactual simulation.
//
// Auto-triggered after tunable_observation emission. For each
// observation whose targeted tunable has a closed-form queuing model
// in the static registry, the parser simulates the predicted effect
// of the suggested direction. Outputs:
//
//   kind:"tunable_counterfactual"          — prediction with 90% CI
//   kind:"tunable_counterfactual_refused"  — no model applies; reason cited
//
// Self-doubt: every prediction discloses `assumptions_made` (Poisson
// arrivals, FCFS, stationary, etc). KS test failure on the arrival
// series downgrades confidence_rank by 0.3. Bootstrap CI > ±25% of
// point estimate caps confidence_rank at 0.5.
//
// Refuses on tunables outside the queuing-theory scope (cache size,
// compressors, featureFlags, anything requiring a process restart).
package main

import (
	"math"
	"math/rand"
	"sort"
)

// modelKind identifies which counterfactual model applies to a tunable.
type modelKind int

const (
	modelNone modelKind = iota
	modelMMc            // M/M/c queuing (tickets, connections)
	modelLinearExtrap   // linear extrapolation under saturation (cache hit rate)
)

// counterfactualModelRegistry maps tunable name → model kind. Anything
// not present here refuses with `reason_code:"no_applicable_model"`.
var counterfactualModelRegistry = map[string]modelKind{
	"wiredTigerConcurrentWriteTransactions": modelMMc,
	"wiredTigerConcurrentReadTransactions":  modelMMc,
	"net.maxIncomingConnections":            modelMMc,
}

// emitTunableCounterfactuals reads tunable_observation Findings emitted
// earlier in the run + the store, and emits one counterfactual or
// counterfactual_refused Finding per observation.
func emitTunableCounterfactuals(c *findingCollector, host string, store *seriesStore, cfg AutoConfig) (emitted int, skipped int) {
	if !cfg.RecommendTunables {
		return 0, 0
	}
	// Walk existing tunable_observation Findings.
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		if k != "tunable_observation" {
			continue
		}
		hf, _ := f["host"].(string)
		if hf != "" && hf != host {
			continue
		}
		tunable, _ := f["tunable"].(string)
		direction, _ := f["direction"].(string)
		// Refusal direction already declines — no point modeling.
		if direction == "refuse" {
			continue
		}
		parentID, _ := f["id"].(string)

		mk := counterfactualModelRegistry[tunable]
		if mk == modelNone {
			c.add(newFinding("tunable_counterfactual_refused", host, tunable, map[string]interface{}{
				"tunable":              tunable,
				"parent_finding_id":    parentID,
				"reason_code":          "no_applicable_model",
				"reason":               "No closed-form queuing/extrapolation model is registered for this tunable. The parser refuses to predict.",
				"models_supported":     []string{"M/M/c (tickets, connections)", "linear extrapolation (cache, with prior captures)"},
			}))
			emitted++
			continue
		}

		switch mk {
		case modelMMc:
			emitted += runMMcCounterfactual(c, host, store, tunable, parentID)
		case modelLinearExtrap:
			// Linear extrapolation requires the fingerprint library to
			// supply a prior at-different-cache-size data point. Refuse
			// in v1 (the library doesn't carry that yet); future P4.4
			// extension can revisit.
			c.add(newFinding("tunable_counterfactual_refused", host, tunable, map[string]interface{}{
				"tunable":           tunable,
				"parent_finding_id": parentID,
				"reason_code":       "model_assumptions_unmet",
				"reason":            "Linear extrapolation requires a prior capture at a different cache size; the fingerprint library does not yet record cache_size in provenance.",
			}))
			emitted++
		}
	}
	return emitted, 0
}

// runMMcCounterfactual fits an M/M/c-shaped model to the ticket queue
// metrics, predicts the effect of doubling capacity (or halving on
// "decrease" direction), bootstraps the arrival series for a 90% CI,
// and emits one tunable_counterfactual Finding.
func runMMcCounterfactual(c *findingCollector, host string, store *seriesStore, tunable, parentID string) int {
	// Pick the right arrival + service metrics per tunable.
	var (
		arrivalMetric  string
		serviceMetric  string
		serviceOpsMetric string // companion ops counter for latency first-difference
		capacityField  string
		capacityOutField string // "out" gauge to add to "available" for true capacity
	)
	switch tunable {
	case "wiredTigerConcurrentWriteTransactions":
		arrivalMetric    = "serverStatus.opcounters.update"
		serviceMetric    = "serverStatus.opLatencies.writes.latency"
		serviceOpsMetric = "serverStatus.opLatencies.writes.ops"
		capacityField    = "serverStatus.wiredTiger.concurrentTransactions.write.available"
		capacityOutField = "serverStatus.wiredTiger.concurrentTransactions.write.out"
	case "wiredTigerConcurrentReadTransactions":
		arrivalMetric    = "serverStatus.opcounters.query"
		serviceMetric    = "serverStatus.opLatencies.reads.latency"
		serviceOpsMetric = "serverStatus.opLatencies.reads.ops"
		capacityField    = "serverStatus.wiredTiger.concurrentTransactions.read.available"
		capacityOutField = "serverStatus.wiredTiger.concurrentTransactions.read.out"
	case "net.maxIncomingConnections":
		arrivalMetric    = "serverStatus.connections.totalCreated"
		serviceMetric    = "" // no clean "service time" for connection lifetime
		capacityField    = "serverStatus.connections.available"
		capacityOutField = "serverStatus.connections.current"
	}

	arrivalID, hasArrival := store.nameToID[arrivalMetric]
	if !hasArrival {
		emitCounterfactualRefused(c, host, tunable, parentID, "missing_arrival_metric",
			"Arrival-rate metric "+arrivalMetric+" not present in the capture.")
		return 1
	}
	arrCol := store.readColumnInt64(arrivalID, nil)
	if len(arrCol) < 60 {
		emitCounterfactualRefused(c, host, tunable, parentID, "series_too_short",
			"Need ≥60 samples to fit a queuing model.")
		return 1
	}
	// First-difference the counter to get per-sample arrivals.
	deltas := make([]float64, 0, len(arrCol)-1)
	for i := 1; i < len(arrCol); i++ {
		d := arrCol[i] - arrCol[i-1]
		if d < 0 {
			continue // counter reset; skip
		}
		deltas = append(deltas, float64(d))
	}
	if len(deltas) < 60 {
		emitCounterfactualRefused(c, host, tunable, parentID, "series_too_short",
			"After dropping counter resets, fewer than 60 arrival deltas.")
		return 1
	}
	meanArrivals := sliceMean(deltas)
	if meanArrivals == 0 {
		emitCounterfactualRefused(c, host, tunable, parentID, "no_arrivals",
			"Arrival rate is zero across the capture; nothing to simulate.")
		return 1
	}

	// Capacity estimate: median of (available[i] + out[i]) per sample.
	// Using available alone understates capacity when tickets are in use;
	// true capacity = available + in-use (out).
	capacity := 128.0
	if capID, ok := store.nameToID[capacityField]; ok {
		capCol := store.readColumnInt64(capID, nil)
		if len(capCol) > 0 {
			totals := make([]int64, len(capCol))
			copy(totals, capCol)
			// Add "out" (in-use) if available for this tunable.
			if capacityOutField != "" {
				if outID, ok2 := store.nameToID[capacityOutField]; ok2 {
					outCol := store.readColumnInt64(outID, nil)
					if len(outCol) == len(totals) {
						for i := range totals {
							totals[i] += outCol[i]
						}
					}
				}
			}
			sort.Slice(totals, func(i, j int) bool { return totals[i] < totals[j] })
			capacity = float64(totals[len(totals)/2])
			if capacity <= 0 {
				capacity = 128
			}
		}
	}

	// Service rate proxy: first-difference the cumulative latency and ops
	// counters, then compute mean service latency as sum(Δlatency)/sum(Δops).
	// The old code averaged the raw cumulative latency counter values, which
	// grows monotonically and does not represent per-op latency.
	if serviceMetric == "" {
		emitCounterfactualRefused(c, host, tunable, parentID, "no_service_metric",
			"M/M/c model needs a per-op latency metric; connections does not expose one.")
		return 1
	}
	srvID, hasSrv := store.nameToID[serviceMetric]
	if !hasSrv {
		emitCounterfactualRefused(c, host, tunable, parentID, "missing_service_metric",
			"Service-time metric "+serviceMetric+" not present.")
		return 1
	}
	srvOpsID, hasSrvOps := store.nameToID[serviceOpsMetric]
	if !hasSrvOps {
		emitCounterfactualRefused(c, host, tunable, parentID, "missing_service_ops_metric",
			"Service-ops counter "+serviceOpsMetric+" not present; cannot compute mean latency.")
		return 1
	}
	latCol := store.readColumnInt64(srvID, nil)
	opsColSrv := store.readColumnInt64(srvOpsID, nil)
	if len(latCol) < 60 || len(opsColSrv) < 60 {
		emitCounterfactualRefused(c, host, tunable, parentID, "series_too_short",
			"Service-time series too short.")
		return 1
	}
	if len(latCol) != len(opsColSrv) {
		emitCounterfactualRefused(c, host, tunable, parentID, "series_length_mismatch",
			"Latency and ops series have different lengths.")
		return 1
	}
	// First-difference both counters and compute mean service latency as
	// sum(Δlatency) / sum(Δops) (ratio-of-means, robust to variable op-rate).
	var sumLat, sumOps float64
	for i := 1; i < len(latCol); i++ {
		dOps := float64(opsColSrv[i] - opsColSrv[i-1])
		if dOps > 0 {
			sumLat += float64(latCol[i] - latCol[i-1])
			sumOps += dOps
		}
	}
	if sumOps == 0 {
		emitCounterfactualRefused(c, host, tunable, parentID, "no_ops_in_capture",
			"All per-interval op-count deltas are zero; cannot compute service latency.")
		return 1
	}
	srvMeanMicros := sumLat / sumOps
	if srvMeanMicros <= 0 {
		emitCounterfactualRefused(c, host, tunable, parentID, "service_time_zero",
			"Service time metric averages to zero across the capture.")
		return 1
	}
	// service rate (ops/sec per server) = 1 / (mean_latency_seconds)
	muPerServer := 1.0 / (srvMeanMicros / 1e6)
	// Current utilization ρ = λ / (c * μ).
	rho := meanArrivals / (capacity * muPerServer)
	if rho > 0.99 {
		rho = 0.99
	}

	// Predicted utilization at doubled capacity.
	newCapacity := capacity * 2
	rhoNew := meanArrivals / (newCapacity * muPerServer)
	predictedDropPct := 0.0
	if rho > 0 {
		predictedDropPct = (rho - rhoNew) / rho * 100
	}

	// Bootstrap 90% CI on rho via resampling the deltas.
	rng := rand.New(rand.NewSource(42)) // deterministic for byte-equality
	const iter = 500
	rhos := make([]float64, 0, iter)
	for i := 0; i < iter; i++ {
		s := 0.0
		for j := 0; j < len(deltas); j++ {
			s += deltas[rng.Intn(len(deltas))]
		}
		bootMean := s / float64(len(deltas))
		bootRho := bootMean / (capacity * muPerServer)
		if bootRho > 0.99 {
			bootRho = 0.99
		}
		rhos = append(rhos, bootRho)
	}
	sort.Float64s(rhos)
	ciLo := rhos[int(0.05*iter)]
	ciHi := rhos[int(0.95*iter)]

	// KS test against exponential (simplified). For deltas that look
	// exponential, max |empirical CDF - exponential CDF| should be small.
	ks := ksExponentialDistance(deltas, meanArrivals)
	assumptions := []string{"M/M/c", "Poisson arrivals", "FCFS service", "stationary"}
	modelViolated := []string{}
	if ks > 0.2 {
		modelViolated = append(modelViolated, "non_poisson_arrivals")
	}

	// Confidence rank starts at 0.7; reductions per the plan.
	confidence := 0.7
	ciSpan := ciHi - ciLo
	if rho > 0 && ciSpan/(rho+1e-9) > 0.5 {
		if confidence > 0.5 {
			confidence = 0.5
		}
	}
	if len(modelViolated) > 0 {
		confidence -= 0.3
	}
	if confidence < 0.0 {
		confidence = 0.0
	}

	c.add(newFinding("tunable_counterfactual", host, tunable, map[string]interface{}{
		"tunable":                       tunable,
		"parent_finding_id":             parentID,
		"model":                         "M/M/c",
		"current_capacity_observed":     capacity,
		"proposed_capacity":             newCapacity,
		"current_utilization":           roundFloat(rho, 4),
		"predicted_utilization":         roundFloat(rhoNew, 4),
		"predicted_utilization_drop_pct": roundFloat(predictedDropPct, 2),
		"current_utilization_ci90":      []float64{roundFloat(ciLo, 4), roundFloat(ciHi, 4)},
		"ks_distance_vs_exponential":    roundFloat(ks, 4),
		"assumptions_made":              assumptions,
		"model_assumption_violated":     modelViolated,
		"confidence_rank":               roundFloat(confidence, 4),
		"disclaimer":                    "Model-based prediction. Assumes Poisson arrivals (KS-tested) and FCFS service. Bootstrap 90% CI shown. Confidence is downgraded when assumptions are violated or CI is wide.",
	}))
	return 1
}

// emitCounterfactualRefused is a small helper for the multiple bail-out
// paths in runMMcCounterfactual.
func emitCounterfactualRefused(c *findingCollector, host, tunable, parentID, code, reason string) {
	c.add(newFinding("tunable_counterfactual_refused", host, tunable, map[string]interface{}{
		"tunable":           tunable,
		"parent_finding_id": parentID,
		"reason_code":       code,
		"reason":            reason,
	}))
}

// ksExponentialDistance returns max |F_emp(x) - F_exp(x)| where F_exp is
// the CDF of an exponential with rate 1/mean(xs). A coarse Kolmogorov-
// Smirnov stat: ≤0.05 is "looks exponential"; >0.2 suggests violation.
func ksExponentialDistance(xs []float64, m float64) float64 {
	if m <= 0 || len(xs) < 2 {
		return 1.0
	}
	sorted := make([]float64, len(xs))
	copy(sorted, xs)
	sort.Float64s(sorted)
	n := float64(len(sorted))
	maxD := 0.0
	for i, x := range sorted {
		empCDF := float64(i+1) / n
		expCDF := 1.0 - math.Exp(-x/m)
		d := empCDF - expCDF
		if d < 0 {
			d = -d
		}
		if d > maxD {
			maxD = d
		}
	}
	return maxD
}
