// analyze_crosshost.go - Cross-host multi-node analysis.
//
// PALI deployments are inherently multi-node (primary + N standbys) so
// single-host analysis is half-blind. When --auto is invoked with a path
// containing per-host FTDC subdirectories, each host's pipeline runs
// independently, then three cross-host Finding kinds are emitted:
//
//   host_only_metrics       — metrics unique to this host vs cluster
//                              intersection (e.g. standby.* lag metrics)
//   cross_host_drift        — metric in intersection where the ratio
//                              between any two hosts' window-aggregates
//                              exceeds the drift threshold (default 5×)
//   cross_host_role_inferred — heuristic role from host-only metrics
//                              (presence of standby.* → standby; election
//                              counters → primary)
//
// Window-aggregate only — point-wise alignment with clock-skew correction
// is explicitly out of scope (V2). 5-minute aggregation absorbs sub-second
// skew; larger skew is rare in healthy clusters.
package main

import (
	"sync"
)

import (
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"time"
)

// hostEntry pairs a stable host name (the original child directory under
// `path`) with the resolved FTDC root inside it. Multi-host dispatch must
// use the entry's name as the Findings `host` field — NOT
// filepath.Base(ftdcDir), since that would yield "diagnostic.data" for
// the common <host>/diagnostic.data/ layout and collide every host onto
// one key.
type hostEntry struct {
	name    string // original child dir name (e.g. "primary", "00-01")
	ftdcDir string // resolved FTDC source path (e.g. "<input>/primary/diagnostic.data")
}

// detectMultiHostInputDetailed returns per-host entries with the stable
// host name preserved across the FTDC-root resolution. Use this instead of
// detectMultiHostInput from main() so multi-host Findings get the correct
// host field.
func detectMultiHostInputDetailed(path string) []hostEntry {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}
	// Single-host bail-out.
	for _, e := range entries {
		if !e.IsDir() && isFTDCShaped(e.Name()) {
			return nil
		}
	}
	out := []hostEntry{}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		childDir := filepath.Join(path, name)
		resolved := resolveFTDCRoot(childDir)
		if resolved == "" {
			continue
		}
		out = append(out, hostEntry{name: name, ftdcDir: resolved})
	}
	if len(out) < 2 {
		return nil
	}
	return out
}

// resolveFTDCRoot looks for FTDC files inside `dir` itself, then inside
// `dir/diagnostic.data/`, then one more level deep (uncommon). Returns the
// path containing the FTDC files, or "" if none found.
func resolveFTDCRoot(dir string) string {
	if hasAnyFTDCFile(dir) {
		return dir
	}
	if hasAnyFTDCFile(filepath.Join(dir, "diagnostic.data")) {
		return filepath.Join(dir, "diagnostic.data")
	}
	// One more level: <host>/<something>/diagnostic.data/
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		nested := filepath.Join(dir, e.Name())
		if hasAnyFTDCFile(filepath.Join(nested, "diagnostic.data")) {
			return filepath.Join(nested, "diagnostic.data")
		}
		if hasAnyFTDCFile(nested) {
			return nested
		}
	}
	return ""
}

// hasAnyFTDCFile checks whether dir contains at least one regular file with
// an FTDC-shaped name (raw timestamp or "metrics." prefix).
func hasAnyFTDCFile(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if isFTDCShaped(e.Name()) {
			return true
		}
	}
	return false
}

// isFTDCShaped returns true for file names that look like FTDC samples:
// either raw Unix-timestamp digits ("1779068871") or the standard
// "metrics.YYYY-MM-DDT…" prefix.
func isFTDCShaped(name string) bool {
	if len(name) == 0 {
		return false
	}
	if name[0] >= '0' && name[0] <= '9' {
		// Plausible timestamp; require all-digit name to avoid false
		// positives on log files like "1foo.log".
		for i := 0; i < len(name); i++ {
			c := name[i]
			if c < '0' || c > '9' {
				return false
			}
		}
		return true
	}
	if len(name) > 8 && name[:8] == "metrics." {
		return true
	}
	return false
}

// runAutoMultiNamed is the host-name-aware multi-host orchestrator. Use this
// from main() so each host's Findings carry the correct host field even when
// the FTDC root is nested (e.g. <host>/diagnostic.data/). hosts[i].name is
// the stable host identifier; samplesByHost/metadataByHost are keyed by
// that name. Kept for backward compatibility with callers that already
// have the per-host samples slices materialised; new callers should prefer
// runAutoMultiStreaming so each host's data never accumulates globally.
func runAutoMultiNamed(w io.Writer, hosts []hostEntry, samplesByHost map[string][]Sample, metadataByHost map[string][]MetadataDoc, cfg AutoConfig) error {
	if len(hosts) == 0 {
		return nil
	}
	hostStores := map[string]*seriesStore{}
	for _, h := range hosts {
		samples := samplesByHost[h.name]
		metadata := metadataByHost[h.name]
		if len(samples) == 0 {
			continue
		}
		if err := runAutoWithHost(w, h.name, samples, metadata, cfg); err != nil {
			return err
		}
		hostStores[h.name] = buildSeriesStoreFromSamples(samples)
	}
	collector := &findingCollector{}
	emitCrossHostFindings(collector, hostStores, samplesByHost, cfg)
	return collector.emitAll(w)
}

// runAutoMultiStreaming orchestrates multi-host --auto by streaming each
// host's FTDC files into a per-host columnar seriesStore one at a time.
// After per-host analysis runs and emits Findings, the host's bulky
// columns + timestamps are released; only the lightweight aggregates
// (nameToID, fullMean, fullSD, isCounter, isDegenerate, observedMin/Max)
// are retained for the cross-host comparison pass.
//
// This eliminates the ~150 GB peak the previous all-hosts-at-once
// accumulation produced on 10 × 48h captures.
func runAutoMultiStreaming(w io.Writer, hosts []hostEntry, matcher *metricMatcher, jobs int, hasFrom, hasTo bool, fromTime, toTime time.Time, cfg AutoConfig) error {
	if len(hosts) == 0 {
		return nil
	}
	hostStores := map[string]*seriesStore{}
	hostRoles := map[string]PALIRole{}

	// P1.6 multi-host-parallel: when cfg.MultiHostParallel > 1, run host
	// pipelines through a bounded worker pool. DEFAULT is 1 (sequential)
	// because decoded FTDC samples are GB-scale per host and N parallel
	// stores trivially OOMs a normal box. We document this in
	// AutoConfig.MultiHostParallel and trust the user who opts into N>1
	// to know their memory budget.
	parallel := cfg.MultiHostParallel
	if parallel <= 0 {
		parallel = 1
	}
	if parallel == 1 {
		// Sequential path — identical to the pre-P1.6 implementation.
		for _, h := range hosts {
			if err := runHostPipeline(w, h, matcher, jobs, hasFrom, hasTo, fromTime, toTime, cfg, hostStores, hostRoles); err != nil {
				return err
			}
		}
	} else {
		// Bounded-pool variant. Serializes write access to w via a mutex
		// because runAutoWithHostStore writes the per-host Finding stream
		// directly to the shared bufio.Writer.
		type result struct {
			err error
		}
		sem := make(chan struct{}, parallel)
		var wgWait sync.WaitGroup
		var writeLock sync.Mutex
		var firstErr error
		var firstErrMu sync.Mutex
		for _, h := range hosts {
			h := h
			sem <- struct{}{}
			wgWait.Add(1)
			go func() {
				defer wgWait.Done()
				defer func() { <-sem }()
				// Each goroutine wraps w in a thin write lock; emit
				// stays one-host-at-a-time for output coherence.
				lockedW := &lockedWriter{mu: &writeLock, w: w}
				if err := runHostPipeline(lockedW, h, matcher, jobs, hasFrom, hasTo, fromTime, toTime, cfg, hostStores, hostRoles); err != nil {
					firstErrMu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					firstErrMu.Unlock()
				}
			}()
		}
		wgWait.Wait()
		if firstErr != nil {
			return firstErr
		}
	}

	collector := &findingCollector{}
	emitCrossHostFindings(collector, hostStores, nil, cfg)
	emitRoleAsymmetry(collector, hostStores, hostRoles, cfg)
	return collector.emitAll(w)
}

// runHostPipeline is the per-host body extracted from runAutoMultiStreaming.
// It builds the host's seriesStore, runs --auto, captures the inferred PALI
// role, drops the column bulk, and stashes the per-host aggregates into
// hostStores/hostRoles for the cross-host pass.
//
// hostStores + hostRoles are concurrently mutated under cfg.MultiHostParallel
// > 1. Callers serialize access via a mutex when parallel > 1.
var multiHostMu sync.Mutex

func runHostPipeline(w io.Writer, h hostEntry, matcher *metricMatcher, jobs int, hasFrom, hasTo bool, fromTime, toTime time.Time, cfg AutoConfig, hostStores map[string]*seriesStore, hostRoles map[string]PALIRole) error {
	store, meta, err := buildSeriesStoreStreaming(h.ftdcDir, matcher, jobs, hasFrom, hasTo, fromTime, toTime)
	if err != nil {
		return err
	}
	if store.nSamples() == 0 {
		return nil
	}
	if err := runAutoWithHostStore(w, h.name, store, meta, cfg); err != nil {
		return err
	}
	role, _ := paliRole(store)
	store.columns = nil
	store.timestamps = nil
	store.materializedSeries = nil

	// Acquire once to write both maps atomically — two separate locks leave
	// an inconsistency window where hostRoles has the entry but hostStores
	// does not, which would confuse emitRoleAsymmetry if it ran concurrently.
	multiHostMu.Lock()
	hostRoles[h.name] = role
	hostStores[h.name] = store
	multiHostMu.Unlock()
	debug.FreeOSMemory()
	return nil
}

// lockedWriter serializes Write calls on an io.Writer behind a mutex.
// Used by the multi-host parallel path so per-host Finding NDJSON
// doesn't interleave on the shared output stream.
type lockedWriter struct {
	mu *sync.Mutex
	w  io.Writer
}

func (lw *lockedWriter) Write(p []byte) (int, error) {
	lw.mu.Lock()
	defer lw.mu.Unlock()
	return lw.w.Write(p)
}

// applyAutoFilters mirrors the single-host time-range + nonzero filtering
// applied in main()'s default path. Extracted so the multi-host dispatch
// can use the same filter semantics without duplicating code.
func applyAutoFilters(samples []Sample, hasFrom, hasTo bool, fromTime, toTime time.Time, nonzero bool) []Sample {
	if hasFrom || hasTo {
		filtered := samples[:0]
		for _, s := range samples {
			if hasFrom && s.Timestamp.Before(fromTime) {
				continue
			}
			if hasTo && s.Timestamp.After(toTime) {
				continue
			}
			filtered = append(filtered, s)
		}
		samples = filtered
	}
	if nonzero {
		samples = stripZeroMetrics(samples)
	}
	return samples
}

// emitCrossHostFindings produces the 3 cross-host Finding kinds. The
// per-host pipelines have already emitted their own Findings to the
// caller's writer; we only emit the comparative ones here.
func emitCrossHostFindings(c *findingCollector, hostStores map[string]*seriesStore, samplesByHost map[string][]Sample, cfg AutoConfig) {
	if len(hostStores) < 2 {
		return
	}
	// Sorted host names for determinism.
	hosts := make([]string, 0, len(hostStores))
	for h := range hostStores {
		hosts = append(hosts, h)
	}
	sort.Strings(hosts)

	// host_only_metrics: per host, metrics not present in EVERY other host.
	// Compute intersection of metric names across all hosts, then per-host
	// difference.
	intersection := map[string]bool{}
	for n := range hostStores[hosts[0]].nameToID {
		intersection[n] = true
	}
	for _, h := range hosts[1:] {
		s := hostStores[h]
		for n := range intersection {
			if _, ok := s.nameToID[n]; !ok {
				delete(intersection, n)
			}
		}
	}
	for _, h := range hosts {
		s := hostStores[h]
		hostOnly := []string{}
		for n := range s.nameToID {
			if !intersection[n] {
				hostOnly = append(hostOnly, n)
			}
		}
		sort.Strings(hostOnly)
		if len(hostOnly) > 0 {
			c.add(newFinding(KindHostOnlyMetrics, h, h, map[string]interface{}{
				"host":           h,
				"metrics":        hostOnly,
				"intersection_n": len(intersection),
				"total_metrics":  len(s.nameToID),
			}))
		}
	}

	// host_metric_set_asymmetry: structural evidence of which metric paths
	// appear on this host but not in the cluster-wide intersection. The
	// analyzer SKILL.md can derive role inferences (primary/standby/etc.)
	// from the evidence; the parser does not name pathologies.
	//
	// Tier A.6 fix: pre-A.6 this emitted kind:"cross_host_role_inferred"
	// with a `role:"primary|standby|unknown"` field — a verdict-style
	// inference that violated the bias-control invariant (never carry
	// `role`/`severity`/`cause`). Replaced by neutral structural facts.
	// See analyzer SKILL.md addendum for the migration.
	for _, h := range hosts {
		s := hostStores[h]
		evidence := []string{}
		for n := range s.nameToID {
			if intersection[n] {
				continue
			}
			// Bucket the indicators we used to use for role inference.
			// Now we surface them as raw evidence rather than naming a role.
			if startsWith(n, "serverStatus.standby.") ||
				containsSubstring(n, "election") ||
				containsSubstring(n, "oplogTruncation") {
				evidence = append(evidence, n)
			}
		}
		sort.Strings(evidence)
		if len(evidence) > 10 {
			evidence = evidence[:10]
		}
		if len(evidence) == 0 {
			continue
		}
		c.add(newFinding(KindHostMetricSetAsymmetry, h, h, map[string]interface{}{
			"host":     h,
			"evidence": evidence,
		}))
	}

	// cross_host_drift: for each intersection metric, compare per-host
	// window-aggregate. Use fullMean from each host's store as the
	// aggregate (already computed during buildSeriesStore). Emit when
	// ratio between any pair of hosts exceeds the drift threshold.
	thresh := cfg.DriftRatio
	if thresh <= 0 {
		thresh = defaultAutoDriftRatio
	}
	intersectionList := make([]string, 0, len(intersection))
	for n := range intersection {
		intersectionList = append(intersectionList, n)
	}
	sort.Strings(intersectionList)
	for _, name := range intersectionList {
		// Drift is a magnitude question: which host's |value| is smallest
		// vs largest. Using signed min/max bugs out on metrics that
		// straddle zero (e.g. one host -1000, another +10 should report
		// 100× ratio not 0.01×). Track by absolute value, but record
		// the actual signed values for the LLM.
		var (
			minAbsVal, maxAbsVal float64
			minHost, maxHost     string
			minSigned, maxSigned float64
			first                = true
		)
		for _, h := range hosts {
			s := hostStores[h]
			id := s.nameToID[name]
			v := s.fullMean[id]
			a := absFloat(v)
			if first {
				minAbsVal, maxAbsVal = a, a
				minSigned, maxSigned = v, v
				minHost, maxHost = h, h
				first = false
				continue
			}
			if a < minAbsVal {
				minAbsVal = a
				minSigned = v
				minHost = h
			}
			if a > maxAbsVal {
				maxAbsVal = a
				maxSigned = v
				maxHost = h
			}
		}
		// Both magnitudes near zero ⇒ no meaningful drift.
		if maxAbsVal < 1e-9 {
			continue
		}
		// Pre-fix: applied a 1e-9 floor to minAbsVal when computing ratio.
		// That floor manufactured astronomical ratios on near-zero baselines
		// (e.g., minAbsVal=1e-11, maxAbsVal=1e-7 → ratio=1e2 even though
		// both values were noise floor). Skip the metric entirely when the
		// minimum side is too small to be a meaningful baseline.
		if minAbsVal < 1e-9 {
			continue
		}
		ratio := maxAbsVal / minAbsVal
		if ratio < thresh {
			continue
		}
		c.add(newFinding(KindCrossHostDrift, "", name, map[string]interface{}{
			"metric":    name,
			"hosts":     hosts,
			"min_host":  minHost,
			"min_value": roundFloat(minSigned, 6),
			"max_host":  maxHost,
			"max_value": roundFloat(maxSigned, 6),
			"ratio":     roundFloat(ratio, 4),
		}))
	}
}

func startsWith(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}

func containsSubstring(s, sub string) bool {
	if len(sub) == 0 {
		return true
	}
	if len(s) < len(sub) {
		return false
	}
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
