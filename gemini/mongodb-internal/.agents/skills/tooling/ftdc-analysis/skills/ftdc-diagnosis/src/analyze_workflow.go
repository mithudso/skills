// analyze_workflow.go — workflow primitives used across captures
// (not just within a single capture):
//
//   - kind:"memory_triplet" — joint shape of the three canonical
//     memory metrics (tcmalloc.current_allocated_bytes, mem.resident,
//     tcmalloc.heap_size) with their slopes and pairwise ratios. NO
//     "fragmentation vs leak" classification (that's analyzer-side);
//     this stage emits the three structural slopes and the analyzer
//     maps them onto interpretation.
//
//   - kind:"recurrence_fingerprint" — compact deterministic hash of
//     the run's event shape. Computed from the sorted list of
//     (synthetic_event.metric_family_count, top member kinds,
//     workload_spec dominant op family). Identical captures produce
//     identical fingerprints; near-identical captures produce
//     fingerprints with measurable hamming distance (analyzer-side).
//
// Bias-control: every field is a raw numeric / structural quantity.
// No verdict, no severity, no role.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
)

// emitMemoryTriplet emits a single Finding per host with the three
// memory metrics' slopes (per second) and pairwise ratios. Gated on
// the presence of at least 2 of the 3 metrics — single-metric
// captures don't have enough info for the joint shape.
func emitMemoryTriplet(c *findingCollector, host string, store *seriesStore) (emitted, skipped int) {
	if store.nSamples() < 2 {
		return 0, 1
	}
	dur := store.timestamps[len(store.timestamps)-1].Sub(store.timestamps[0]).Seconds()
	if dur <= 0 {
		return 0, 1
	}

	type triplet struct {
		name string
		key  string
	}
	metrics := []triplet{
		{"current_allocated_bytes", "serverStatus.tcmalloc.generic.current_allocated_bytes"},
		{"mem_resident_mb", "serverStatus.mem.resident"},
		{"heap_size_bytes", "serverStatus.tcmalloc.generic.heap_size"},
	}

	fields := map[string]interface{}{"host": host}
	presentCount := 0
	values := map[string][2]int64{} // name -> [first, last]
	for _, m := range metrics {
		id, ok := store.nameToID[m.key]
		if !ok {
			continue
		}
		col := store.columns[id]
		if col == nil || col.Len() < 2 {
			continue
		}
		vals := col.ReadAllInt64(nil)
		first := vals[0]
		last := vals[len(vals)-1]
		values[m.name] = [2]int64{first, last}
		slopePerSec := float64(last-first) / dur
		fields[m.name+"_first"] = first
		fields[m.name+"_last"] = last
		fields[m.name+"_slope_per_sec"] = roundFloat(slopePerSec, 4)
		presentCount++
	}
	if presentCount < 2 {
		return 0, 1
	}

	// Pairwise ratios for the three canonical pairings the analyzer
	// uses to distinguish fragmentation / leak / cache growth.
	if cur, hasCur := values["current_allocated_bytes"]; hasCur {
		if res, hasRes := values["mem_resident_mb"]; hasRes {
			// mem.resident is in MB; convert to bytes for ratio.
			lastResBytes := res[1] * 1024 * 1024
			if lastResBytes > 0 {
				fields["allocated_to_resident_ratio_last"] = roundFloat(float64(cur[1])/float64(lastResBytes), 4)
			}
		}
		if heap, hasHeap := values["heap_size_bytes"]; hasHeap && heap[1] > 0 {
			fields["allocated_to_heap_ratio_last"] = roundFloat(float64(cur[1])/float64(heap[1]), 4)
		}
	}

	c.add(newFinding("memory_triplet", host, "global", fields))
	return 1, 0
}

// emitRecurrenceFingerprint computes a compact stable hash of the
// run's event shape, drawn from the already-emitted Findings. Two
// runs of the same workload pattern should produce the same
// fingerprint; subtle shifts produce a different fingerprint.
//
// Input shape (sorted before hashing for determinism):
//   - synthetic_event metric_family_counts (one per event)
//   - top member kinds per synthetic_event (sorted, comma-joined)
//   - workload_spec dominant op family (the highest *_rate_ps field)
//   - missing_signal classes present
//
// Output: 8-hex-char fingerprint + the input shape (for human audit).
func emitRecurrenceFingerprint(c *findingCollector, host string) (emitted int) {
	type shapeBit struct {
		key string
		val string
	}
	var bits []shapeBit

	for _, f := range c.items {
		k, _ := f["kind"].(string)
		switch k {
		case "synthetic_event":
			if mfc, ok := f["metric_family_count"].(int); ok {
				bits = append(bits, shapeBit{key: "synthetic_event.mfc", val: fmt.Sprintf("%d", mfc)})
			}
			// top member kinds: pull from the members array.
			if members, ok := f["members"].([]interface{}); ok {
				kinds := map[string]bool{}
				for _, mi := range members {
					m, _ := mi.(map[string]interface{})
					if mk, ok := m["kind"].(string); ok {
						kinds[mk] = true
					}
				}
				var ksorted []string
				for kk := range kinds {
					ksorted = append(ksorted, kk)
				}
				sort.Strings(ksorted)
				bits = append(bits, shapeBit{key: "synthetic_event.kinds", val: joinStrings(ksorted, ",")})
			}
		case "workload_spec":
			// Pick the highest *_rate_ps field. Suffix "_rate_ps" is
			// 8 characters; len(k2) must be strictly > 8 so the
			// substring slice doesn't include the whole key.
			//
			// Init bestVal to -Inf so a workload_spec with all-negative
			// rates (rare; signed counter deltas can be negative) still
			// picks a dominant_op. Pre-fix, var bestVal float64 = 0.0
			// suppressed dominant_op selection on negative-rate inputs.
			const suffix = "_rate_ps"
			var bestKey string
			bestVal := math.Inf(-1)
			for k2, v := range f {
				if len(k2) > len(suffix) && k2[len(k2)-len(suffix):] == suffix {
					if vf, ok := v.(float64); ok && (vf > bestVal || (vf == bestVal && k2 < bestKey)) {
						bestVal = vf
						bestKey = k2
					}
				}
			}
			if bestKey != "" {
				bits = append(bits, shapeBit{key: "workload.dominant_op", val: bestKey})
			}
		case "missing_signal":
			// Tier C renamed the field from "class" to "category"; read the
			// new name. Fall back to the old name on legacy Findings so a
			// mid-flight schema change doesn't break recurrence hashing.
			if cat, ok := f["category"].(string); ok {
				bits = append(bits, shapeBit{key: "missing_signal.category", val: cat})
			} else if class, ok := f["class"].(string); ok {
				bits = append(bits, shapeBit{key: "missing_signal.category", val: class})
			}
		}
	}

	if len(bits) == 0 {
		return 0
	}

	// Stable sort by key:val.
	sort.Slice(bits, func(i, j int) bool {
		if bits[i].key != bits[j].key {
			return bits[i].key < bits[j].key
		}
		return bits[i].val < bits[j].val
	})

	h := sha256.New()
	shape := make([]interface{}, 0, len(bits))
	for _, b := range bits {
		fmt.Fprintf(h, "%s=%s\n", b.key, b.val)
		shape = append(shape, map[string]interface{}{
			"key":   b.key,
			"value": b.val,
		})
	}
	sum := h.Sum(nil)
	fp := hex.EncodeToString(sum[:4])

	c.add(newFinding("recurrence_fingerprint", host, "global", map[string]interface{}{
		"host":        host,
		"fingerprint": fp,
		"shape_bits":  len(bits),
		"shape":       shape,
	}))
	return 1
}

// joinStrings concatenates with sep; tiny helper to avoid pulling in
// strings package for a one-liner.
func joinStrings(s []string, sep string) string {
	out := ""
	for i, v := range s {
		if i > 0 {
			out += sep
		}
		out += v
	}
	return out
}
