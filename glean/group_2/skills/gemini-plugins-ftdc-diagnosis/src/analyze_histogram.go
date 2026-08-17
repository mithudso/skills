// analyze_histogram.go - Structural histogram detection + per-window
// percentile inference.
//
// Histograms are discovered by SHAPE, not by metric name. Two encoding
// patterns are recognized:
//
//	indexed-buckets pattern:
//	    <prefix>.histogram.<N>.count
//	    <prefix>.histogram.<N>.micros        (optional, used to derive
//	                                         per-bucket representative
//	                                         latency = micros/count)
//	  Found in opLatencies and various other places. Version-stable: the
//	  bucket bounds don't need hardcoding if .micros siblings are present.
//	  When .micros is absent we fall back to a unitless cumulative-count
//	  reading (counts per bucket only, no percentile latency estimate).
//
//	bracket-bound pattern:
//	    <prefix>.[lo, hi).count
//	    <prefix>.totalCount
//	  Found across disagg histogram families. Bounds are parsed from the
//	  metric names directly, so any future disagg histogram is detected
//	  automatically.
//
// Neither pattern hardcodes a metric prefix. New histogram families added
// to MongoDB are detected automatically as long as they follow one of the
// two shapes above. See run_metadata Finding for the structural recognizer
// names — what we look for, not which prefixes match.
package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

// histogramFamily is a single discovered histogram, ready for percentile
// inference. Same struct serves both indexed-buckets and bracket-bound.
type histogramFamily struct {
	prefix string // e.g. "serverStatus.opLatencies.writes"

	// pattern indicates how this family is encoded.
	pattern string // "indexed" or "bracket"

	// indexed-bucket fields:
	//   countKeys[i]   = "<prefix>.histogram.<i>.count"
	//   microsKeys[i]  = "<prefix>.histogram.<i>.micros" (empty if absent)
	countKeys  []string
	microsKeys []string

	// bracket-bound fields:
	//   buckets[i].lo / hi from the parsed key
	//   buckets[i].key is "<prefix>.[lo, hi).count"
	//   totalKey is "<prefix>.totalCount" if present
	buckets  []disaggBucketSchema
	totalKey string
}

// disaggBucketSchema is one parsed bracket-bound bucket schema entry
// (bounds + the metric name to read the cumulative count from).
type disaggBucketSchema struct {
	lo, hi float64
	key    string
}

// inferHistogramPercentiles runs Stage 4. Discovers histograms structurally
// from the metric set, then emits one histogram_percentile (or
// histogram_invalid) Finding per (family, window). Parallel across families
// — each family's window inference is independent.
func inferHistogramPercentiles(host string, baselineSamples, stressSamples, allSamples []Sample, cfg AutoConfig) []Finding {
	_ = cfg
	if len(baselineSamples) < 2 || len(stressSamples) < 2 {
		return nil
	}
	keys := metricKeysFromSamples(allSamples)
	families := discoverHistogramFamilies(keys)
	if len(families) == 0 {
		return nil
	}
	type famResult struct {
		baseline Finding
		stress   Finding
	}
	results := make([]famResult, len(families))
	var wg sync.WaitGroup
	for idx, fam := range families {
		idx, fam := idx, fam
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[idx].baseline = inferFamilyWindow(host, fam, "baseline", baselineSamples)
			results[idx].stress = inferFamilyWindow(host, fam, "stress", stressSamples)
		}()
	}
	wg.Wait()
	out := make([]Finding, 0, 2*len(families))
	for _, r := range results {
		if r.baseline != nil {
			out = append(out, r.baseline)
		}
		if r.stress != nil {
			out = append(out, r.stress)
		}
	}
	return out
}

// discoverHistogramFamilies scans metric keys for both encoding patterns
// and returns one histogramFamily per detected histogram. The function is
// pure / deterministic — same keys produce same families in same order.
func discoverHistogramFamilies(keys map[string]bool) []histogramFamily {
	// Pass 1: indexed-buckets. Group keys matching "<prefix>.histogram.<N>.count"
	// by <prefix>, recording bucket indices and whether a .micros sibling
	// exists.
	type indexedAcc struct {
		counts map[int]string
		micros map[int]string
	}
	indexed := map[string]*indexedAcc{}
	for k := range keys {
		prefix, idx, suffix, ok := parseIndexedHistogramKey(k)
		if !ok {
			continue
		}
		a, exists := indexed[prefix]
		if !exists {
			a = &indexedAcc{counts: map[int]string{}, micros: map[int]string{}}
			indexed[prefix] = a
		}
		switch suffix {
		case "count":
			a.counts[idx] = k
		case "micros":
			a.micros[idx] = k
		}
	}

	// Pass 2: bracket-bound. Group keys matching "<prefix>.[lo, hi).count"
	// by <prefix>. Also detect the optional <prefix>.totalCount sibling.
	type bracketAcc struct {
		buckets  []disaggBucketSchema
		totalKey string
	}
	bracket := map[string]*bracketAcc{}
	for k := range keys {
		prefix, lo, hi, ok := parseBracketHistogramKey(k)
		if !ok {
			continue
		}
		a, exists := bracket[prefix]
		if !exists {
			a = &bracketAcc{}
			bracket[prefix] = a
		}
		a.buckets = append(a.buckets, disaggBucketSchema{lo: lo, hi: hi, key: k})
	}
	for prefix, a := range bracket {
		totalKey := prefix + ".totalCount"
		if keys[totalKey] {
			a.totalKey = totalKey
		}
	}

	out := []histogramFamily{}
	// Emit indexed families. Require at least 2 indexed buckets to count as
	// a histogram (single .histogram.0 is just a counter); require .count
	// for each registered bucket index.
	prefixes := make([]string, 0, len(indexed))
	for p := range indexed {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)
	for _, p := range prefixes {
		a := indexed[p]
		if len(a.counts) < 2 {
			continue
		}
		idxs := make([]int, 0, len(a.counts))
		for i := range a.counts {
			idxs = append(idxs, i)
		}
		sort.Ints(idxs)
		countKeys := make([]string, len(idxs))
		microsKeys := make([]string, len(idxs))
		for j, i := range idxs {
			countKeys[j] = a.counts[i]
			microsKeys[j] = a.micros[i] // empty if absent
		}
		out = append(out, histogramFamily{
			prefix: p, pattern: "indexed",
			countKeys: countKeys, microsKeys: microsKeys,
		})
	}
	// Emit bracket families. Require at least 2 buckets.
	prefixes = prefixes[:0]
	for p := range bracket {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)
	for _, p := range prefixes {
		a := bracket[p]
		if len(a.buckets) < 2 {
			continue
		}
		// Sort buckets by lo bound for cumulative interpolation.
		sort.Slice(a.buckets, func(i, j int) bool { return a.buckets[i].lo < a.buckets[j].lo })
		out = append(out, histogramFamily{
			prefix: p, pattern: "bracket",
			buckets: a.buckets, totalKey: a.totalKey,
		})
	}
	return out
}

// parseIndexedHistogramKey matches "<prefix>.histogram.<N>.<suffix>" where
// <N> is a non-negative integer and <suffix> ∈ {"count", "micros"}. Returns
// the prefix (without ".histogram"), the bucket index, the suffix, and ok=true.
func parseIndexedHistogramKey(k string) (string, int, string, bool) {
	// Find the last ".count" or ".micros" suffix.
	var suffix string
	var rest string
	if strings.HasSuffix(k, ".count") {
		suffix = "count"
		rest = k[:len(k)-len(".count")]
	} else if strings.HasSuffix(k, ".micros") {
		suffix = "micros"
		rest = k[:len(k)-len(".micros")]
	} else {
		return "", 0, "", false
	}
	// rest should now end with ".histogram.<N>".
	dot := strings.LastIndex(rest, ".")
	if dot < 0 {
		return "", 0, "", false
	}
	idxStr := rest[dot+1:]
	idx, err := strconv.Atoi(idxStr)
	if err != nil || idx < 0 || idx > 1024 {
		return "", 0, "", false
	}
	before := rest[:dot]
	if !strings.HasSuffix(before, ".histogram") {
		return "", 0, "", false
	}
	prefix := before[:len(before)-len(".histogram")]
	return prefix, idx, suffix, true
}

// parseBracketHistogramKey matches "<prefix>.[lo, hi).count" where lo,hi
// are numeric (allowing decimals). Returns the prefix and parsed bounds.
func parseBracketHistogramKey(k string) (string, float64, float64, bool) {
	if !strings.HasSuffix(k, ").count") {
		return "", 0, 0, false
	}
	rest := k[:len(k)-len(").count")]
	// rest should now end with ".[lo, hi"
	openIdx := strings.LastIndex(rest, ".[")
	if openIdx < 0 {
		return "", 0, 0, false
	}
	body := rest[openIdx+2:]
	parts := strings.Split(body, ",")
	if len(parts) != 2 {
		return "", 0, 0, false
	}
	lo, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return "", 0, 0, false
	}
	hi, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return "", 0, 0, false
	}
	prefix := rest[:openIdx]
	return prefix, lo, hi, true
}

// inferFamilyWindow computes per-window cumulative-delta counts for one
// family and either emits a histogram_percentile Finding or, on counter
// non-monotonicity, a histogram_invalid Finding.
func inferFamilyWindow(host string, fam histogramFamily, windowName string, samples []Sample) Finding {
	first := samples[0].Metrics
	last := samples[len(samples)-1].Metrics

	switch fam.pattern {
	case "indexed":
		// MongoDB opLatencies histogram structure:
		//   serverStatus.opLatencies.<op>.histogram is a BSON array of
		//   {micros: <upper_bound_us>, count: <cumulative_count>}. After
		//   FTDC array-flattening this becomes:
		//     <prefix>.histogram.<i>.micros = bucket i's UPPER BOUND (µs)
		//     <prefix>.histogram.<i>.count  = cumulative count through bucket i
		//
		// The `micros` value is therefore STATIC per bucket (the bound),
		// not a cumulative time counter. Prior code treated it as a
		// cumulative microsecond accumulator and divided by count — that
		// produced 0 or near-0 for almost every bucket. The corrected
		// path uses micros[i] as the bucket upper bound and linearly
		// interpolates within [bounds[i-1], bounds[i]).
		counts := make([]int64, len(fam.countKeys))
		bucketBounds := make([]float64, len(fam.countKeys)) // upper bound µs
		hasBound := make([]bool, len(fam.countKeys))
		invalid := false
		for i, ck := range fam.countKeys {
			f, fok := first[ck]
			l, lok := last[ck]
			if !fok || !lok {
				continue
			}
			if l < f {
				invalid = true
				break
			}
			counts[i] = l - f
			if fam.microsKeys[i] != "" {
				if mv, ok := last[fam.microsKeys[i]]; ok {
					bucketBounds[i] = float64(mv)
					hasBound[i] = true
				}
			}
		}
		if invalid {
			return newFinding(KindHistogramInvalid, host, fam.prefix+"|"+windowName, map[string]interface{}{
				"metric":  fam.prefix,
				"window":  windowName,
				"pattern": "indexed",
				"reason":  "bucket counter decreased (counter reset or shard movement)",
			})
		}
		var total int64
		for _, c := range counts {
			total += c
		}
		if total < 4 {
			return nil
		}
		fields := map[string]interface{}{
			"metric":       fam.prefix,
			"window":       windowName,
			"pattern":      "indexed",
			"ops":          total,
			"buckets":      len(counts),
			"window_start": samples[0].Timestamp.UTC().Format(tsFormat),
			"window_end":   samples[len(samples)-1].Timestamp.UTC().Format(tsFormat),
		}
		// Emit percentile latency estimates if .micros siblings provided
		// bucket upper bounds; otherwise emit only the raw bucket-count
		// distribution.
		anyBound := false
		for _, h := range hasBound {
			if h {
				anyBound = true
				break
			}
		}
		if anyBound {
			// Build bounds array for interpolation. bounds[i] = upper
			// edge of bucket i. Bucket 0 spans [0, bounds[0]). Bucket i
			// (for i>0) spans [bounds[i-1], bounds[i]).
			bounds := make([]float64, 0, len(counts)+1)
			bounds = append(bounds, 0) // implicit lower edge for bucket 0
			for i := range counts {
				if hasBound[i] {
					bounds = append(bounds, bucketBounds[i])
				} else if len(bounds) > 0 {
					// Missing bound — extrapolate by doubling the prior
					// edge (preserves monotonicity for interpolation).
					// Tier A.5 fix: if the prior edge is 0 (e.g. bucket 0
					// has no .micros sibling), doubling yields a cascade
					// of zeros and interpolatePercentile collapses to 0.
					// Fall back to a geometric default of 1<<i µs so the
					// bucket has a non-degenerate width.
					prev := bounds[len(bounds)-1]
					if prev == 0 {
						bounds = append(bounds, float64(int64(1)<<uint(i+1)))
					} else {
						bounds = append(bounds, prev*2)
					}
				}
			}
			fields["p50_us"] = roundFloat(interpolatePercentile(counts, bounds, total, 0.50), 1)
			fields["p95_us"] = roundFloat(interpolatePercentile(counts, bounds, total, 0.95), 1)
			fields["p99_us"] = roundFloat(interpolatePercentile(counts, bounds, total, 0.99), 1)
			fields["bucket_upper_bounds"] = bucketBounds
		} else {
			// No bounds available; emit distribution shape only.
			dist := make([]int64, len(counts))
			copy(dist, counts)
			fields["bucket_counts"] = dist
		}
		return newFinding(KindHistogramPercentile, host, fam.prefix+"|"+windowName, fields)

	case "bracket":
		counts := make([]int64, len(fam.buckets))
		invalid := false
		var total int64
		for i, b := range fam.buckets {
			f, fok := first[b.key]
			l, lok := last[b.key]
			if !fok || !lok {
				continue
			}
			if l < f {
				invalid = true
				break
			}
			counts[i] = l - f
			total += counts[i]
		}
		if invalid {
			return newFinding(KindHistogramInvalid, host, fam.prefix+"|"+windowName, map[string]interface{}{
				"metric":  fam.prefix,
				"window":  windowName,
				"pattern": "bracket",
				"reason":  "bucket counter decreased (counter reset or shard movement)",
			})
		}
		if total < 4 {
			return nil
		}
		bounds := make([]float64, 0, len(fam.buckets)+1)
		bounds = append(bounds, fam.buckets[0].lo)
		for _, b := range fam.buckets {
			bounds = append(bounds, b.hi)
		}
		p50 := interpolatePercentile(counts, bounds, total, 0.50)
		p95 := interpolatePercentile(counts, bounds, total, 0.95)
		p99 := interpolatePercentile(counts, bounds, total, 0.99)
		// Bracket-pattern bounds carry the metric's own units (µs for
		// latency families, dimensionless for batch-size etc.). Use
		// the same _us key as opLatencies for consistency; consumers
		// can identify by metric name.
		return newFinding(KindHistogramPercentile, host, fam.prefix+"|"+windowName, map[string]interface{}{
			"metric":       fam.prefix,
			"window":       windowName,
			"pattern":      "bracket",
			"ops":          total,
			"buckets":      len(fam.buckets),
			"p50_us":       roundFloat(p50, 4),
			"p95_us":       roundFloat(p95, 4),
			"p99_us":       roundFloat(p99, 4),
			"window_start": samples[0].Timestamp.UTC().Format(tsFormat),
			"window_end":   samples[len(samples)-1].Timestamp.UTC().Format(tsFormat),
		})
	}
	return nil
}

// interpolatePercentile finds the bucket containing the target cumulative
// count and linearly interpolates within it. Open-ended last buckets are
// estimated as 2× their lower bound (a conservative heuristic).
func interpolatePercentile(counts []int64, bounds []float64, total int64, p float64) float64 {
	if total <= 0 || len(counts) == 0 {
		return 0
	}
	target := p * float64(total)
	cum := int64(0)
	for i, c := range counts {
		newCum := cum + c
		if float64(newCum) >= target && c > 0 {
			lo := bounds[i]
			var hi float64
			if i+1 < len(bounds) {
				hi = bounds[i+1]
			} else {
				hi = lo * 2
			}
			// Zero-width-bucket guard: if hi <= lo (corrupted bounds or
			// duplicated edges), skip the interpolation and return lo.
			// Pre-fix produced silent zero answers when bounds collided.
			if hi <= lo {
				return lo
			}
			within := target - float64(cum)
			frac := within / float64(c)
			return lo + frac*(hi-lo)
		}
		cum = newCum
	}
	if len(bounds) > 0 {
		return bounds[len(bounds)-1]
	}
	return 0
}

// metricKeysFromSamples returns the union of metric names appearing in any
// sample.
func metricKeysFromSamples(samples []Sample) map[string]bool {
	out := map[string]bool{}
	for i := range samples {
		for k := range samples[i].Metrics {
			out[k] = true
		}
	}
	return out
}

// histogramRecognizerNames returns the names of the structural recognizers
// for inclusion in run_metadata. These describe WHAT we look for, not WHICH
// prefixes match — so the LLM can audit our coverage independent of which
// histograms happen to be in the current capture.
func histogramRecognizerNames() []string {
	return []string{
		"indexed-buckets: <prefix>.histogram.<N>.count[+.micros]",
		"bracket-bound: <prefix>.[lo, hi).count [+ <prefix>.totalCount]",
	}
}
