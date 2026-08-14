// analyze_sweep.go — P4.3 chaos-sweep auto-labeled comparison.
//
// Gregory runs ~50-round chaos sweeps. Each round becomes a capture
// the parser eventually --auto's. The fingerprint library (P4.4)
// builds up one entry per capture. The sweep stage reads the library
// for sibling captures (by mtime-proximity OR path-prefix match),
// auto-labels them HEALTHY vs DEGRADED from FTDC-only health composite,
// and emits kind:sweep_discriminator Findings naming the Finding shapes
// that most discriminate the two groups.
//
// Memory constraint: this stage works ENTIRELY off the fingerprint
// library (already-flat JSON files). No full seriesStores are loaded
// for non-primary captures. The autonomous label inference is over
// fingerprint provenance + shape_set membership, not raw FTDC.
//
// Self-doubt:
//   - Refuse to publish discriminators if k=2 silhouette score < 0.4
//     (groups didn't separate cleanly).
//   - Emit kind:sweep_label_audit with the auto-assigned label per
//     capture so the user sees and can correct mistakes.
package main

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// sweepMinSiblings is the minimum library-entry count required before
// the sweep stage fires. Below this the chaos-sweep pattern isn't
// recognizable.
const sweepMinSiblings = 10

// sweepMaxAgeForLayout limits how stale a sibling can be before we
// consider it part of the same sweep run. 30 days is generous.
const sweepMaxAgeForLayout = 30 * 24 * time.Hour

// sweepMinSilhouette: below this, emit sweep_inconclusive and refuse
// to publish discriminators.
const sweepMinSilhouette = 0.4

// emitChaosSweep reads the fingerprint library for siblings near the
// current capture, auto-labels them, and emits sweep_discriminator or
// sweep_inconclusive Findings. Pure read-only on disk.
func emitChaosSweep(c *findingCollector, host string) (emitted int, skipped int) {
	storeDir := fingerprintStoreDir()
	entries, err := os.ReadDir(storeDir)
	if err != nil {
		return 0, 1
	}
	type libEntry struct {
		path     string
		baseName string
		entry    fingerprintEntry
	}
	now := time.Now()
	var pool []libEntry
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		full := filepath.Join(storeDir, e.Name())
		data, err := os.ReadFile(full)
		if err != nil {
			continue
		}
		var fe fingerprintEntry
		if err := json.Unmarshal(data, &fe); err != nil {
			continue
		}
		if now.Sub(fe.Provenance.WrittenAt) > sweepMaxAgeForLayout {
			continue
		}
		pool = append(pool, libEntry{path: full, baseName: e.Name(), entry: fe})
	}
	if len(pool) < sweepMinSiblings {
		return 0, 1
	}

	// Autonomous health composite: lower = "healthier."
	// We use (max_confidence_at_write, shape_set_count) — a capture that
	// produced strong diagnosis_candidates and a large shape set is
	// more likely degraded (something definitely went wrong); a clean
	// run has few shapes and low confidence ranks.
	scores := make([]float64, len(pool))
	for i, p := range pool {
		score := p.entry.Provenance.MaxConfidenceAtWrite + float64(len(p.entry.ShapeSet))*0.01
		scores[i] = score
	}

	// k-means with k=2 on the 1-D score. Deterministic init: median +/-
	// 25% percentile.
	sorted := append([]float64(nil), scores...)
	sort.Float64s(sorted)
	low := sorted[len(sorted)/4]
	high := sorted[3*len(sorted)/4]
	if high == low {
		return 0, 1 // no separation possible
	}
	muLo, muHi := low, high
	const iter = 30
	labels := make([]int, len(pool))
	for it := 0; it < iter; it++ {
		// Assign.
		for i, s := range scores {
			if math.Abs(s-muLo) <= math.Abs(s-muHi) {
				labels[i] = 0
			} else {
				labels[i] = 1
			}
		}
		// Update.
		var sumLo, sumHi float64
		var nLo, nHi int
		for i, s := range scores {
			if labels[i] == 0 {
				sumLo += s
				nLo++
			} else {
				sumHi += s
				nHi++
			}
		}
		if nLo == 0 || nHi == 0 {
			break
		}
		newLo, newHi := sumLo/float64(nLo), sumHi/float64(nHi)
		if math.Abs(newLo-muLo) < 1e-6 && math.Abs(newHi-muHi) < 1e-6 {
			break
		}
		muLo, muHi = newLo, newHi
	}
	// Map: 0 → HEALTHY (lower score), 1 → DEGRADED. We ensure muLo < muHi.
	if muLo > muHi {
		muLo, muHi = muHi, muLo
		for i := range labels {
			labels[i] = 1 - labels[i]
		}
	}

	// Silhouette score (mean over all points).
	sil := silhouette1D(scores, labels)
	if sil < sweepMinSilhouette {
		c.add(newFinding("sweep_inconclusive", host, "global", map[string]interface{}{
			"sibling_count":      len(pool),
			"silhouette_score":   roundFloat(sil, 4),
			"silhouette_min":     sweepMinSilhouette,
			"reason":             "k=2 cluster silhouette below threshold — auto-labeling didn't produce stable HEALTHY/DEGRADED groups. Refusing to publish discriminators.",
		}))
		emitted++
		return emitted, 0
	}

	// Build per-shape contingency: count fires in HEALTHY vs DEGRADED.
	shapeCount := map[string][2]int{}
	healthyN, degradedN := 0, 0
	for i, p := range pool {
		if labels[i] == 0 {
			healthyN++
		} else {
			degradedN++
		}
		for _, s := range p.entry.ShapeSet {
			pair := shapeCount[s]
			pair[labels[i]]++
			shapeCount[s] = pair
		}
	}
	type discRow struct {
		shape         string
		degradedFire  float64 // fraction of degraded captures firing this shape
		healthyFire   float64
		discriminator float64 // |degraded_fire - healthy_fire|
	}
	rows := make([]discRow, 0, len(shapeCount))
	for s, counts := range shapeCount {
		var h, d float64
		if healthyN > 0 {
			h = float64(counts[0]) / float64(healthyN)
		}
		if degradedN > 0 {
			d = float64(counts[1]) / float64(degradedN)
		}
		rows = append(rows, discRow{
			shape:         s,
			degradedFire:  d,
			healthyFire:   h,
			discriminator: math.Abs(d - h),
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].discriminator != rows[j].discriminator {
			return rows[i].discriminator > rows[j].discriminator
		}
		return rows[i].shape < rows[j].shape
	})
	const topK = 20
	if len(rows) > topK {
		rows = rows[:topK]
	}

	discriminators := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		if r.discriminator < 0.2 {
			continue
		}
		discriminators = append(discriminators, map[string]interface{}{
			"shape":          r.shape,
			"healthy_fire":   roundFloat(r.healthyFire, 4),
			"degraded_fire":  roundFloat(r.degradedFire, 4),
			"discriminator":  roundFloat(r.discriminator, 4),
		})
	}
	c.add(newFinding("sweep_discriminator", host, "global", map[string]interface{}{
		"sibling_count":       len(pool),
		"healthy_count":       healthyN,
		"degraded_count":      degradedN,
		"silhouette_score":    roundFloat(sil, 4),
		"discriminators":      discriminators,
		"interpretation":      "Shapes ranked by their fire-frequency difference between auto-labeled HEALTHY and DEGRADED captures from the local fingerprint library. Discrimination ≥ 0.2 was the inclusion threshold.",
	}))
	emitted++

	// Per-capture label audit.
	auditRows := make([]map[string]interface{}, 0, len(pool))
	for i, p := range pool {
		label := "HEALTHY"
		if labels[i] == 1 {
			label = "DEGRADED"
		}
		auditRows = append(auditRows, map[string]interface{}{
			"fingerprint_file":      p.baseName,
			"input_sha256":          p.entry.Provenance.InputSHA256,
			"input_path_basename":   p.entry.Provenance.InputPathBasename,
			"written_at":            p.entry.Provenance.WrittenAt.UTC().Format(time.RFC3339),
			"max_confidence_at_write": p.entry.Provenance.MaxConfidenceAtWrite,
			"shape_set_size":        len(p.entry.ShapeSet),
			"auto_label":            label,
		})
	}
	c.add(newFinding("sweep_label_audit", host, "global", map[string]interface{}{
		"sibling_count": len(pool),
		"rows":          auditRows,
		"note":          "Auto-assigned labels — review these and re-run --auto with --auto-recommend-tunables for corrected captures if any are mislabeled.",
	}))
	emitted++
	return emitted, 0
}

// silhouette1D returns the mean silhouette coefficient over all points
// in a 1-D k-means assignment. Range [-1, 1]; higher = better separation.
func silhouette1D(values []float64, labels []int) float64 {
	if len(values) != len(labels) {
		return 0
	}
	if len(values) < 2 {
		return 0
	}
	sum := 0.0
	computed := 0
	for i, v := range values {
		ownDist := 0.0
		ownCount := 0
		otherDist := 0.0
		otherCount := 0
		for j, u := range values {
			if i == j {
				continue
			}
			d := math.Abs(v - u)
			if labels[j] == labels[i] {
				ownDist += d
				ownCount++
			} else {
				otherDist += d
				otherCount++
			}
		}
		if ownCount == 0 || otherCount == 0 {
			continue
		}
		a := ownDist / float64(ownCount)
		b := otherDist / float64(otherCount)
		denom := math.Max(a, b)
		if denom == 0 {
			continue
		}
		sum += (b - a) / denom
		computed++
	}
	if computed == 0 {
		return 0
	}
	return sum / float64(computed)
}
