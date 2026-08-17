// analyze_emit_t2.go - --auto-emit-t2 bundle writer.
//
// Emitted when --auto-emit-t2 PATH is set. Captures every Finding NDJSON line
// via findingTee (a transparent pass-through wrapper around the single
// bufio.Writer in main), then after the normal auto output has flushed, writes
// a self-contained bundle directory:
//
//	PATH/
//	  context.t2      t2 state file with relative FTDC paths
//	  analysis.md     marker glossary + diagnosis summary
//	  findings.ndjson full Finding stream (copy of stdout)
//	  README.txt
//	  open-in-t2.sh   launcher script. Run as: bash open-in-t2.sh
//	                  (.sh extension avoids macOS Gatekeeper quarantine-bit
//	                  blocking that affects .command files from downloaded
//	                  archives; explicit `bash` invocation is always safe.)
//	  ftdc/
//	    <host>/...     FTDC file copies (or symlinks)
//
// A .tgz of the bundle is produced alongside by default; skip with
// --auto-emit-t2-no-tarball.
//
// Zero-behavior-change constraint: when the flag is unset, findingTee is never
// constructed and this file contributes no allocations or code paths.
package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// ── findingTee ────────────────────────────────────────────────────────────────

// findingTee wraps an io.Writer and accumulates complete NDJSON lines.
// Write is a transparent pass-through: it always returns the inner writer's
// (n, err), so callers see no behavioral difference.
// The tee is goroutine-safe; lockedWriter in multi-host mode calls Write from
// multiple goroutines concurrently.
type findingTee struct {
	inner io.Writer
	mu    sync.Mutex
	lines []json.RawMessage
	buf   []byte // partial-line accumulation across Write calls
}

func newFindingTee(inner io.Writer) *findingTee {
	return &findingTee{inner: inner}
}

func (t *findingTee) Write(p []byte) (int, error) {
	n, err := t.inner.Write(p)
	if err != nil {
		return n, err
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.buf = append(t.buf, p...)
	for {
		idx := bytes.IndexByte(t.buf, '\n')
		if idx < 0 {
			break
		}
		line := bytes.TrimSpace(t.buf[:idx])
		if len(line) > 0 {
			cp := make([]byte, len(line))
			copy(cp, line)
			t.lines = append(t.lines, json.RawMessage(cp))
		}
		t.buf = t.buf[idx+1:]
	}
	return n, nil
}

// capturedLines returns a snapshot of the accumulated lines under mu.
func (t *findingTee) capturedLines() []json.RawMessage {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]json.RawMessage, len(t.lines))
	copy(out, t.lines)
	return out
}

// ── emitT2Bundle entry point ──────────────────────────────────────────────────

// emitT2Bundle is called (deferred) after the auto dispatch completes.
// It is best-effort: errors go to stderr; the parser's exit code is never changed.
func emitT2Bundle(tee *findingTee, cfg AutoConfig) {
	if err := doEmitT2Bundle(tee, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "ftdc_parser: --auto-emit-t2: %v\n", err)
	}
}

func doEmitT2Bundle(tee *findingTee, cfg AutoConfig) error {
	rawLines := tee.capturedLines()

	// Parse findings.
	findings := make([]Finding, 0, len(rawLines))
	for _, raw := range rawLines {
		var f Finding
		if err := json.Unmarshal(raw, &f); err == nil {
			findings = append(findings, f)
		}
	}

	hosts := resolveHostsForEmit(cfg.InputPath)
	bundleDir := cfg.AutoEmitT2Path

	// Write into a temp directory first so a failure leaves no partial bundle.
	parentDir := filepath.Dir(bundleDir)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("mkdir parent %s: %w", parentDir, err)
	}
	tmpDir, err := os.MkdirTemp(parentDir, ".t2bundle-")
	if err != nil {
		return fmt.Errorf("mkdirtemp: %w", err)
	}
	success := false
	defer func() {
		if !success {
			os.RemoveAll(tmpDir)
		}
	}()

	// Copy FTDC files into bundle.
	topK := cfg.AutoEmitT2TopK
	if topK <= 0 {
		topK = 50
	}
	relPaths, err := copyFTDCIntoBundle(tmpDir, hosts, cfg.AutoEmitT2LinkFtdc)
	if err != nil {
		return fmt.Errorf("copy ftdc: %w", err)
	}

	// Collect markers before building JSON (both need the marker list).
	markerEntries := collectMarkerEntries(findings, topK)
	sort.Slice(markerEntries, func(i, j int) bool {
		return markerEntries[i].epochMs < markerEntries[j].epochMs
	})

	// Build t2 state JSON.
	t2JSON, err := buildT2StateJSON(findings, relPaths, markerEntries, topK)
	if err != nil {
		return fmt.Errorf("build t2 json: %w", err)
	}

	// Process optional log file before writing analysis.md.
	var logRes *logResult
	if cfg.AutoLogFile != "" {
		mainEventTime := deriveMainEventTime(findings)
		if res, err := processLogFile(cfg.AutoLogFile, findings, tmpDir, mainEventTime); err == nil {
			logRes = res
		}
	}

	// Write bundle files.
	analysisMd := buildAnalysisMd(findings, markerEntries, hosts, cfg)
	if logRes != nil {
		analysisMd += buildLogSection(logRes, cfg.AutoLogFile)
	}
	for name, content := range map[string][]byte{
		"context.t2":      t2JSON,
		"analysis.md":     []byte(analysisMd),
		"findings.ndjson": buildFindingsNDJSON(rawLines),
		"README.txt":      []byte(bundleReadme()),
	} {
		if err := os.WriteFile(filepath.Join(tmpDir, name), content, 0644); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
	}
	// Single launcher. Written 0755 so `./open-in-t2.sh` works on Linux
	// (no Gatekeeper there). On macOS, instruct users to run `bash
	// open-in-t2.sh` instead of `./open-in-t2.sh` — explicit bash invocation
	// bypasses the quarantine-bit check that blocks downloaded executables;
	// the checked binary is /bin/bash (system-trusted), not the script itself.
	if err := os.WriteFile(filepath.Join(tmpDir, "open-in-t2.sh"), []byte(bundleOpenScript()), 0755); err != nil {
		return fmt.Errorf("write open-in-t2.sh: %w", err)
	}

	// Rename temp dir to final bundle dir (atomic on same filesystem).
	os.RemoveAll(bundleDir)
	if err := os.Rename(tmpDir, bundleDir); err != nil {
		// Cross-device: fall back to recursive copy.
		if err2 := copyDirRecursive(tmpDir, bundleDir); err2 != nil {
			return fmt.Errorf("rename bundle dir: %w (fallback copy: %v)", err, err2)
		}
		os.RemoveAll(tmpDir)
	}
	success = true

	// Tarball (default; disabled by --auto-emit-t2-no-tarball).
	if !cfg.AutoEmitT2NoTarball {
		tgzPath := bundleDir + ".tgz"
		if err := createTarball(bundleDir, tgzPath); err != nil {
			fmt.Fprintf(os.Stderr, "ftdc_parser: --auto-emit-t2: tarball: %v\n", err)
		}
	}
	return nil
}

// resolveHostsForEmit returns the host entries to bundle. For multi-host
// inputs (directory-of-directories), returns the full set. For single-host
// inputs, wraps the resolved FTDC root in a one-element slice.
func resolveHostsForEmit(inputPath string) []hostEntry {
	if hosts := detectMultiHostInputDetailed(inputPath); len(hosts) > 0 {
		return hosts
	}
	// Single host: resolve the FTDC root and use the input path's basename as
	// the host name.
	ftdcDir := resolveFTDCRoot(inputPath)
	if ftdcDir == "" {
		// Input is a single file or already at an FTDC root with no sub-dir layout.
		ftdcDir = inputPath
		if info, err := os.Stat(inputPath); err == nil && !info.IsDir() {
			ftdcDir = filepath.Dir(inputPath)
		}
	}
	return []hostEntry{{name: hostIDFromPath(inputPath), ftdcDir: ftdcDir}}
}

// ── FTDC file copy ────────────────────────────────────────────────────────────

// copyFTDCIntoBundle copies (or symlinks) each host's FTDC files into
// <bundleDir>/ftdc/<host>/ and returns a map[hostName][]relPath where
// relPath is of the form "ftdc/<host>/<filename>" — the relative paths
// written into dataset.files in the t2 state.
func copyFTDCIntoBundle(bundleDir string, hosts []hostEntry, link bool) (map[string][]string, error) {
	result := make(map[string][]string, len(hosts))
	for _, h := range hosts {
		destDir := filepath.Join(bundleDir, "ftdc", h.name)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return nil, fmt.Errorf("mkdir %s: %w", destDir, err)
		}
		entries, err := os.ReadDir(h.ftdcDir)
		if err != nil {
			// ftdcDir might itself be a single FTDC file.
			if info, statErr := os.Stat(h.ftdcDir); statErr == nil && !info.IsDir() {
				rel, copyErr := copyOrLinkFile(h.ftdcDir, destDir, link)
				if copyErr != nil {
					return nil, copyErr
				}
				result[h.name] = append(result[h.name], "ftdc/"+h.name+"/"+rel)
			}
			continue
		}
		for _, e := range entries {
			if e.IsDir() || strings.HasPrefix(e.Name(), ".") {
				continue
			}
			if !isFTDCShaped(e.Name()) {
				continue
			}
			src := filepath.Join(h.ftdcDir, e.Name())
			rel, err := copyOrLinkFile(src, destDir, link)
			if err != nil {
				return nil, err
			}
			result[h.name] = append(result[h.name], "ftdc/"+h.name+"/"+rel)
		}
	}
	return result, nil
}

func copyOrLinkFile(src, destDir string, link bool) (string, error) {
	name := filepath.Base(src)
	dest := filepath.Join(destDir, name)
	if link {
		if err := os.Symlink(src, dest); err != nil && !os.IsExist(err) {
			return "", fmt.Errorf("symlink %s: %w", src, err)
		}
		return name, nil
	}
	return name, copyFile(src, dest)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func copyDirRecursive(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return copyFile(path, target)
	})
}

// ── Descriptor name conversion ────────────────────────────────────────────────

// ftdcPathToDescriptorName converts a dot-separated FTDC metric path (as used
// in Finding.metric) to the t2 descriptor name used in "view: Custom" override
// maps. The naming rules are derived from descriptors.ts prefix helpers:
//
//	serverStatus.wiredTiger.X.Y.Z  →  "ss wt X Y Z"
//	serverStatus.X.Y.Z              →  "ss X Y Z"   (X ≠ "wiredTiger")
//	systemMetrics.X.Y.Z             →  "system X Y Z"
//	replSetGetStatus.X.Y            →  "rs X Y"
//
// Returns ("", false) for unknown root segments; callers drop those entries.
// Descriptors with hand-crafted names in descriptors.ts (e.g. custom `name:`
// fields) may not match — those produce no-op "include" keys that t2 silently
// ignores per views.ts:80-86 contains().
func ftdcPathToDescriptorName(metric string) (string, bool) {
	dot := strings.IndexByte(metric, '.')
	if dot < 0 {
		return "", false
	}
	root := metric[:dot]
	rest := metric[dot+1:]
	switch root {
	case "serverStatus":
		idx := strings.IndexByte(rest, '.')
		if idx < 0 {
			return "ss " + rest, true
		}
		first := rest[:idx]
		tail := rest[idx+1:]
		if first == "wiredTiger" {
			return "ss wt " + strings.ReplaceAll(tail, ".", " "), true
		}
		return "ss " + strings.ReplaceAll(rest, ".", " "), true
	case "systemMetrics":
		return "system " + strings.ReplaceAll(rest, ".", " "), true
	case "replSetGetStatus":
		return "rs " + strings.ReplaceAll(rest, ".", " "), true
	}
	return "", false
}

// ── Custom view picker ────────────────────────────────────────────────────────

// pickCustomViewDescriptors builds the "view: Custom" override map by reading
// metric fields directly from Findings (data-driven; no hand-curated table).
// Returns map[descriptorName]"include".
func pickCustomViewDescriptors(findings []Finding, topK int) map[string]string {
	scores := make(map[string]float64)

	// Index findings by ID for one-level anchor-ref resolution.
	byID := make(map[string]Finding, len(findings))
	for _, f := range findings {
		if id, ok := f["id"].(string); ok {
			byID[id] = f
		}
	}

	addMetric := func(metric string, score float64) {
		if name, ok := ftdcPathToDescriptorName(metric); ok {
			scores[name] += score
		}
	}

	for _, f := range findings {
		kind, _ := f["kind"].(string)
		score := findingViewScore(f)

		// Direct metric carrier.
		if m, ok := f["metric"].(string); ok && m != "" {
			addMetric(m, score)
		}
		// capacity_pressure carries both metric and ceiling_metric.
		if m, ok := f["ceiling_metric"].(string); ok && m != "" {
			addMetric(m, score)
		}
		// seal_event / seal_event_candidate carry local_spike_metrics[] in
		// evidence_summary — these are the seal-trigger descriptors, NOT
		// downstream checkpoint metrics.
		if es, ok := f["evidence_summary"].(map[string]interface{}); ok {
			if lsm, ok := es["local_spike_metrics"].([]interface{}); ok {
				for _, m := range lsm {
					if s, ok := m.(string); ok {
						addMetric(s, score)
					}
				}
			}
		}
		// joint_changepoint carries member_metrics[].metric.
		if mm, ok := f["member_metrics"].([]interface{}); ok {
			for _, entry := range mm {
				if em, ok := entry.(map[string]interface{}); ok {
					if m, ok := em["metric"].(string); ok {
						addMetric(m, score)
					}
				}
			}
		}

		// Anchor-ref carriers: walk supporting_finding_ids one level deep.
		switch kind {
		case "diagnosis_candidate", "hypothesis_test", "precursor_candidates",
			"evidence_convergence", "joint_changepoint":
			for _, field := range []string{
				"supporting_finding_ids", "required_finding_ids",
				"member_finding_ids", "synthetic_event_id",
			} {
				refs := t2BundleStringSlice(f, field)
				for _, refID := range refs {
					if refF, ok := byID[refID]; ok {
						if m, ok := refF["metric"].(string); ok {
							addMetric(m, score)
						}
					}
				}
			}
		case "host_pathology":
			if sub, ok := f["subkind"].(string); ok {
				for _, name := range hostPathologyViewDescriptors(sub) {
					scores[name] += score
				}
			}
		}
	}

	// Augmentation: when any WiredTiger metric is in scope, also include the
	// PALI disagg context metrics and the WT checkpoint health metrics. The
	// checkpoint descriptors were historically absent from Custom because the
	// checkpoint_stall detector emitted 0 Findings (wrong metric path bug);
	// adding them here ensures they always appear when WT is in scope.
	for name := range scores {
		if strings.Contains(name, "metrics disagg") || strings.Contains(name, "ss wt") {
			for _, extra := range disaggGeneralViewDescriptors() {
				if _, exists := scores[extra]; !exists {
					scores[extra] = 0.1
				}
			}
			for _, extra := range wtCheckpointViewDescriptors() {
				if _, exists := scores[extra]; !exists {
					scores[extra] = 0.1
				}
			}
			break
		}
	}

	// Always include baseline context descriptors.
	for _, b := range t2BaselineDescriptors() {
		if _, exists := scores[b]; !exists {
			scores[b] = 0
		}
	}

	// Rank and cap at budget (baseline always retained).
	type kv struct {
		name  string
		score float64
	}
	ranked := make([]kv, 0, len(scores))
	for n, s := range scores {
		ranked = append(ranked, kv{n, s})
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].score != ranked[j].score {
			return ranked[i].score > ranked[j].score
		}
		return ranked[i].name < ranked[j].name
	})

	baselineSet := make(map[string]bool)
	for _, b := range t2BaselineDescriptors() {
		baselineSet[b] = true
	}
	budget := topK + len(t2BaselineDescriptors())
	if budget > 80 {
		budget = 80
	}
	result := make(map[string]string, len(ranked))
	nonBaselineCount := 0
	for _, e := range ranked {
		if baselineSet[e.name] {
			result[e.name] = "include"
		} else if nonBaselineCount < budget {
			result[e.name] = "include"
			nonBaselineCount++
		}
	}
	return result
}

// findingViewScore returns a comparable priority for the Finding when ranking
// descriptor candidates.
func findingViewScore(f Finding) float64 {
	kind, _ := f["kind"].(string)
	switch kind {
	case "changepoint", "diagnosis_candidate", "seal_event", "seal_event_candidate":
		if cr, ok := f["confidence_rank"].(float64); ok {
			return cr + 1.0
		}
	case "spike", "local_spike":
		if z, ok := f["z_score"].(float64); ok {
			return z / 10.0
		}
	case "regression":
		if ds, ok := f["deviation_sigmas"].(float64); ok {
			return ds / 3.0
		}
	case "mover":
		if s, ok := f["score"].(float64); ok {
			return s
		}
	case "variance_shift":
		if cp, ok := f["cusum_peak"].(float64); ok {
			return cp / 5.0
		}
	case "checkpoint_stall", "counter_stall":
		// Stall Findings are high priority — score proportional to stalled duration,
		// capped so they don't crowd out all other markers.
		if s, ok := f["stalled_duration_s"].(float64); ok {
			score := s / 300.0 // 5 min stall = score 1.0; 1 h stall = score 12
			if score > 15 {
				score = 15
			}
			return score
		}
		return 5.0 // duration-threshold path (checkpoint_stall only)
	}
	return 1.0
}

func hostPathologyViewDescriptors(subkind string) []string {
	switch subkind {
	case "cpu_steal":
		return []string{"system cpu steal", "system cpu user", "system cpu system"}
	case "tcmalloc":
		return []string{
			"ss tcmalloc generic current_allocated_bytes",
			"ss tcmalloc generic heap_size",
			"ss tcmalloc pageheap_free_bytes",
		}
	case "network":
		return []string{"system network lo bytesIn", "system network lo bytesOut"}
	}
	return nil
}

func disaggGeneralViewDescriptors() []string {
	return []string{
		"ss metrics disagg logServer logSealedDetailsApplied",
		"ss metrics disagg logServer logSealedDetailsFallback",
		"ss metrics disagg materializedLsnAdvanceCount",
		"ss metrics disagg getPageRequest numQueued",
		"ss metrics disagg paliHandleGet paliHandleGetTotalLatencyUs",
	}
}

// wtCheckpointViewDescriptors are added to the Custom view whenever any WT
// metric is in scope. These are the key signals for checkpoint health and were
// historically absent from Custom because checkpoint_stall emitted 0 Findings
// (bug: wrong metric path). Descriptor names verified against live FTDC data
// (serverStatus.wiredTiger.checkpoint.* subsystem).
func wtCheckpointViewDescriptors() []string {
	return []string{
		"ss wt checkpoint generation",
		"ss wt checkpoint progress state",
		"ss wt checkpoint most recent time (msecs)",
		"ss wt checkpoint total succeed number of checkpoints",
	}
}

func t2BaselineDescriptors() []string {
	return []string{
		"ss opcounters insert",
		"ss opcounters update",
		"ss opcounters delete",
		"ss opcounters query",
		"ss connections current",
		"ss wt cache bytes currently in the cache",
		"ss wt cache tracked dirty bytes in the cache",
		"system cpu user",
	}
}

// t2BundleStringSlice extracts a []string from a Finding field that may be
// []interface{} (JSON-decoded) or a direct []string (already typed).
func t2BundleStringSlice(f Finding, key string) []string {
	v, ok := f[key]
	if !ok {
		return nil
	}
	switch tv := v.(type) {
	case []string:
		return tv
	case []interface{}:
		out := make([]string, 0, len(tv))
		for _, elem := range tv {
			if s, ok := elem.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

// ── Markers ───────────────────────────────────────────────────────────────────

type t2MarkerEntry struct {
	epochMs int64
	host    string
	kind    string
	at      string // ISO 8601 as emitted by the parser
	metric  string
	score   float64
	blurb   string // 1-2 sentence "why this marker matters"
}

// markerBlurbForFinding returns a 1-2 sentence explanation of why a Finding
// matters, suitable for the marker glossary in analysis.md. The blurb is
// kind-keyed and pulls live numeric fields where available.
func markerBlurbForFinding(f Finding) string {
	kind, _ := f["kind"].(string)
	metric, _ := f["metric"].(string)

	switch kind {
	case "changepoint":
		dir, _ := f["direction"].(string)
		base := "Abrupt, persistent level shift (ED-PELT + CUSUM cross-validated). The metric crossed a threshold and stayed there — a behavioral mode change."
		if dir == "up" {
			base += " Direction: up."
		} else if dir == "down" {
			base += " Direction: down."
		}
		return base
	case "variance_shift":
		base := "Sustained volatility increase (CUSUM on rolling-median deviation). Instability without a mean shift often signals resource contention or an oscillating feedback loop."
		if strings.Contains(metric, "dirty bytes") {
			base += " A simultaneous checkpoint-generation stall confirms dirty data is accumulating without being flushed."
		}
		if cp, ok := f["cusum_peak"].(float64); ok {
			base += fmt.Sprintf(" CUSUM peak: %.1f.", cp)
		}
		return base
	case "spike", "local_spike":
		if z, ok := f["z_score"].(float64); ok {
			return fmt.Sprintf("Momentary outlier burst (robust z-score %.1fσ). Transient — may indicate lock contention, GC pause, or a brief eviction storm.", z)
		}
		return "Momentary outlier burst. Transient — may indicate lock contention, GC pause, or a brief eviction storm."
	case "checkpoint_stall":
		rc, _ := f["reason_code"].(string)
		if rc == "generation_counter_stalled" {
			stalledS, _ := f["stalled_duration_s"].(float64)
			return fmt.Sprintf("WiredTiger checkpoint generation counter stopped advancing for %.0f seconds. No new dirty pages are being flushed; dirty bytes accumulate unboundedly until the stall resolves or the process restarts.", stalledS)
		}
		maxMs, _ := f["max_duration_ms"].(float64)
		return fmt.Sprintf("WiredTiger checkpoint took %.0f ms — exceeding 80%% of the configured syncdelay. Sustained long checkpoints block writes and may cascade into eviction pressure.", maxMs)
	case "counter_stall":
		subsystem, _ := f["subsystem"].(string)
		stalledS, _ := f["stalled_duration_s"].(float64)
		switch subsystem {
		case "pali_log_producer":
			blurb := fmt.Sprintf("PALI log producer stopped generating new entries for %.0f seconds while the local oplog kept advancing. This is the primary signal for a disaggregated-storage producer-thread freeze: the materializer stalls, WT checkpoints stop advancing, and dirty bytes accumulate unboundedly.", stalledS)
			if peak, ok := f["waiter_peak"].(float64); ok && peak > 0 {
				blurb += fmt.Sprintf(" Corroborated: waiters.replication peaked at %.0f (operations queued waiting for replication that PALI was not draining).", peak)
			}
			return blurb
		case "pali_rate_limiter":
			return fmt.Sprintf("PALI rate limiter stopped admitting new log entries for %.0f seconds. Advances in lockstep with the log producer — a simultaneous freeze with phylog.generatedMillis confirms a producer-thread stall rather than throttling.", stalledS)
		case "repl_apply":
			return fmt.Sprintf("Replication apply batch counter stopped advancing for %.0f seconds. The secondary is not processing new oplog entries — replication lag is accumulating.", stalledS)
		}
		return fmt.Sprintf("Subsystem heartbeat counter (%s) stopped advancing for %.0f seconds while the process remained alive. The subsystem has frozen.", subsystem, stalledS)
	case "seal_event", "seal_event_candidate":
		return "PALI log seal: the standby materializer was instructed to freeze application at a stable LSN before a topology change (step-up or reconfiguration)."
	case "gap":
		if as, ok := f["at_start"].(string); ok {
			if ae, ok := f["at_end"].(string); ok {
				_ = as
				_ = ae
			}
		}
		return "FTDC gap: consecutive sample timestamps differ by more than 5 seconds. The mongod process was stalled (kernel scheduling, IO, GC) or the FTDC writer was blocked during this interval."
	case "diagnosis_candidate":
		if r, ok := f["rationale"].(string); ok && r != "" {
			return r
		}
		if pn, ok := f["pattern_name"].(string); ok {
			return fmt.Sprintf("Parser-identified structural pattern: %s.", pn)
		}
		return "Parser-identified diagnosis candidate — see diagnoses section for detail."
	case "monotonic_trend":
		dir, _ := f["direction"].(string)
		if dir == "up" {
			return "Metric is trending steadily upward across the capture — may indicate a slow accumulation (memory leak, growing connection pool, index bloat)."
		}
		return "Metric is trending steadily downward across the capture — may indicate resource depletion or a subsystem slowing over time."
	case "evidence_convergence", "synthetic_event":
		return "Multiple independent signals converge at this timestamp, reinforcing that something structurally changed at this point."
	case "election_storm":
		return "Three or more term changes detected within a 5-minute window — indicates replica set instability or repeated failovers."
	case "host_pathology":
		sub, _ := f["subkind"].(string)
		switch sub {
		case "cpu_steal":
			return "CPU steal time elevated — the hypervisor is scheduling other VMs on this host's physical cores, causing unpredictable latency spikes."
		case "tcmalloc":
			return "TCMalloc heap shows fragmentation or leak pattern — physical memory consumed by the allocator exceeds what MongoDB explicitly requested."
		case "network":
			return "Network throughput anomaly detected on the loopback interface."
		}
		return "Host-level resource pathology detected."
	}
	return "Significant anomaly detected in this metric at this timestamp."
}

// collectMarkerEntries extracts timestamp+score entries from Findings, capped
// at topK after sorting by score descending (tiebreak: kind+at lexicographic).
func collectMarkerEntries(findings []Finding, topK int) []t2MarkerEntry {
	type candidate struct {
		me    t2MarkerEntry
		score float64
	}
	var cands []candidate

	tryAt := func(f Finding, atStr string, score float64) {
		if atStr == "" {
			return
		}
		t, err := time.Parse(tsFormat, atStr)
		if err != nil {
			t, err = time.Parse(time.RFC3339, atStr)
			if err != nil {
				t, err = time.Parse(time.RFC3339Nano, atStr)
				if err != nil {
					return
				}
			}
		}
		host, _ := f["host"].(string)
		kind, _ := f["kind"].(string)
		metric, _ := f["metric"].(string)
		cands = append(cands, candidate{t2MarkerEntry{
			epochMs: t.UnixMilli(),
			host:    host,
			kind:    kind,
			at:      atStr,
			metric:  metric,
			score:   score,
			blurb:   markerBlurbForFinding(f),
		}, score})
	}

	for _, f := range findings {
		kind, _ := f["kind"].(string)
		score := findingViewScore(f)
		switch kind {
		case "changepoint", "variance_shift", "spike", "local_spike",
			"counter_wrap", "counter_inversion", "term_change",
			"oplog_window_low", "alert_threshold", "seal_event",
			"seal_event_candidate", "checkpoint_stall", "counter_stall":
			if at, ok := f["at"].(string); ok {
				tryAt(f, at, score)
			}
		case "gap", "evidence_convergence", "synthetic_event":
			if at, ok := f["at_start"].(string); ok {
				tryAt(f, at, score)
			}
			if at, ok := f["at_end"].(string); ok {
				tryAt(f, at, score*0.5)
			}
		case "election_storm":
			if at, ok := f["window_start"].(string); ok {
				tryAt(f, at, score)
			}
		case "diagnosis_candidate":
			if at, ok := f["at"].(string); ok {
				tryAt(f, at, score)
			}
		}
	}

	sort.Slice(cands, func(i, j int) bool {
		if cands[i].score != cands[j].score {
			return cands[i].score > cands[j].score
		}
		ki := cands[i].me.kind + cands[i].me.at
		kj := cands[j].me.kind + cands[j].me.at
		return ki < kj
	})
	if len(cands) > topK {
		cands = cands[:topK]
	}
	result := make([]t2MarkerEntry, len(cands))
	for i, c := range cands {
		result[i] = c.me
	}
	return result
}

// ── t2 state JSON builder ─────────────────────────────────────────────────────

// allT2ViewNames is the complete list of views registered in views.ts. We emit
// an empty {} for each so t2's restoreState gets a full undo snapshot. Extra
// keys t2 doesn't recognise are silently ignored.
var allT2ViewNames = []string{
	"Custom",
	"All metrics",
	"Basic",
	"Memory",
	"CPU",
	"Bottlenecks",
	"Network",
	"Cache Overflow Usage",
	"Parallel Operations",
	"Compact",
	"WT Checkpoint",
	"WT Cursors",
	"WT Eviction Efficiency",
	"WT Handle Counts",
	"WT IO",
	"WT Hot Pages",
	"WT Workload, Cache Hit Rate, and IO",
	"Query",
	"DX: Colocated Data and Journal",
	"DX: Pinning Single Oplog Batch",
	"DX: tcmalloc Heap Allocation",
	"DX: Table Drop Performance",
	"Replication",
	"Timeseries",
	"Index Builds",
	"Mongot",
}

func buildT2StateJSON(findings []Finding, relPaths map[string][]string, markerEntries []t2MarkerEntry, topK int) ([]byte, error) {
	// Derive time range from run_metadata window_start / window_end fields.
	// Multi-host runs emit one run_metadata per host; take min start, max end.
	var tMin, tMax float64
	for _, f := range findings {
		if f["kind"] != "run_metadata" {
			continue
		}
		if ws, ok := f["window_start"].(string); ok {
			if t, err := time.Parse(time.RFC3339Nano, ws); err == nil {
				ms := float64(t.UnixMilli())
				if tMin == 0 || ms < tMin {
					tMin = ms
				}
			}
		}
		if we, ok := f["window_end"].(string); ok {
			if t, err := time.Parse(time.RFC3339Nano, we); err == nil {
				ms := float64(t.UnixMilli())
				if ms > tMax {
					tMax = ms
				}
			}
		}
	}
	if tMin == 0 && tMax == 0 {
		// Fallback: scan all marker entries.
		for _, me := range markerEntries {
			ms := float64(me.epochMs)
			if tMin == 0 || ms < tMin {
				tMin = ms
			}
			if ms > tMax {
				tMax = ms
			}
		}
	}
	// 5% padding on each side for visual breathing room.
	if tMax > tMin {
		pad := (tMax - tMin) * 0.05
		tMin -= pad
		tMax += pad
	} else if tMin > 0 {
		// Single-point fallback: ±5 minutes.
		tMin -= 5 * 60 * 1000
		tMax += 5 * 60 * 1000
	}

	// dataset.files and userTags, in stable host-name order.
	var files []string
	userTags := make(map[string]interface{})
	hostNames := make([]string, 0, len(relPaths))
	for h := range relPaths {
		hostNames = append(hostNames, h)
	}
	sort.Strings(hostNames)
	for _, h := range hostNames {
		for _, p := range relPaths[h] {
			files = append(files, p)
			userTags[p] = h
		}
	}
	if files == nil {
		files = []string{}
	}

	// Markers as float64 epoch-ms (JSON number).
	markerTS := make([]float64, len(markerEntries))
	for i, me := range markerEntries {
		markerTS[i] = float64(me.epochMs)
	}

	customView := pickCustomViewDescriptors(findings, topK)

	// Resolve t2 source/timeline options based on the input topology.
	//
	// sourcePerFile must be FALSE — when true, each file becomes its own
	// SourceName, defeating the per-host merge that t2's dataset.Load does
	// off EmbeddedSource (the hostInfo.system.hostname from FTDC metadata).
	// For an Atlas single-host capture (N file chunks of one continuous
	// host) sourcePerFile:true splits the chunks into N pseudo-sources and,
	// combined with timeline:"aligned", time-shifts them onto a common
	// origin so only the first file's data is visible at any real wall
	// timestamp. Setting sourcePerFile:false lets t2 group by hostname.
	//
	// sourcePerDir TRUE only for multi-host bundles, where each subdir
	// under ftdc/ is one host. For single-host bundles it must be FALSE.
	//
	// timeline: keep "normal" (wall-clock) unless we're overlaying multiple
	// hosts that may have differing capture starts; "aligned" only earns
	// its keep in multi-host mode.
	multiHost := len(relPaths) > 1
	timeline := "normal"
	if multiHost {
		timeline = "aligned"
	}

	// Build context as a plain map so dynamic "view: X" keys sort predictably.
	//
	// showValues:true matches t2's default and enables the hover-value
	// popups (per-timeseries readout at the cursor). Setting it false here
	// suppressed that for everyone who opened the bundle.
	ctx := map[string]interface{}{
		"ctxName": "Analysis",
		"options": map[string]interface{}{
			"merge":         false,
			"showValues":    true,
			"commonScale":   true,
			"autoScale":     true,
			"sourcePerFile": false,
			"sourcePerDir":  multiHost,
			"chartSize":     "medium",
			"timeline":      timeline,
		},
		"viewSelector": map[string]interface{}{"viewName": "Custom"},
		"ctx": map[string]interface{}{
			"range": map[string]interface{}{
				"$type": "TimeRange",
				"tMin":  tMin,
				"tMax":  tMax,
			},
		},
		"dataset": map[string]interface{}{
			"files":       files,
			"defaultTags": map[string]interface{}{},
			"userTags":    userTags,
			"hide":        map[string]interface{}{},
			"shifts":      map[string]interface{}{},
		},
		"markers":         markerTS,
		"sectionHeadings": map[string]interface{}{},
		"datasources":     []interface{}{},
		"view: Custom":    customView,
	}
	for _, v := range allT2ViewNames {
		if v == "Custom" {
			continue
		}
		ctx["view: "+v] = map[string]interface{}{}
	}

	state := map[string]interface{}{
		"version":       1,
		"stateFilePath": "",
		"contexts":      []interface{}{ctx},
	}
	return json.MarshalIndent(state, "", "  ")
}

// ── Sidecar markdown ──────────────────────────────────────────────────────────

func buildAnalysisMd(findings []Finding, markers []t2MarkerEntry, hosts []hostEntry, cfg AutoConfig) string {
	var sb strings.Builder

	// H1: run identity.
	sb.WriteString("# FTDC Analysis\n\n")
	if cfg.InputPath != "" {
		sb.WriteString(fmt.Sprintf("**Input**: `%s`\n\n", cfg.InputPath))
	}
	if cfg.InputFingerprint != nil {
		sb.WriteString(fmt.Sprintf("**Fingerprint**: `%s`\n\n", cfg.InputFingerprint.TopSHA256))
	}
	sb.WriteString(fmt.Sprintf("**Hosts**: %d  \n", len(hosts)))
	sb.WriteString(fmt.Sprintf("**Parser version**: `%s`\n\n", ParserVersion))

	// Emit capture date range from run_metadata so timestamps elsewhere in this
	// document (especially in diagnosis rationale strings) can be unambiguously
	// interpreted. This is especially important when the capture spans midnight.
	var wStart, wEnd string
	for _, f := range findings {
		if f["kind"] != "run_metadata" {
			continue
		}
		if ws, ok := f["window_start"].(string); ok && wStart == "" {
			wStart = ws
		}
		if we, ok := f["window_end"].(string); ok {
			wEnd = we
		}
	}
	if wStart != "" || wEnd != "" {
		sb.WriteString(fmt.Sprintf("**Capture window**: `%s` → `%s`\n\n", wStart, wEnd))
		// Warn if the capture crosses midnight so readers know bare HH:MM:SS timestamps
		// are ambiguous without their date prefix.
		if wStart != "" && wEnd != "" {
			ts, _ := time.Parse(time.RFC3339Nano, wStart)
			te, _ := time.Parse(time.RFC3339Nano, wEnd)
			if !ts.IsZero() && !te.IsZero() && ts.Day() != te.Day() {
				sb.WriteString("> **Note**: this capture spans multiple calendar days. " +
					"All timestamps in this document use full ISO 8601 format " +
					"(e.g. `2006-01-02T15:04:05Z`) — bare `HH:MM:SS` times " +
					"without a date are ambiguous and should not be used.\n\n")
			}
		}
	}

	// H2: Marker glossary — map letter → finding details.
	sb.WriteString("## Marker glossary\n\n")
	if len(markers) == 0 {
		sb.WriteString("_No significant markers._\n\n")
	} else {
		for i, me := range markers {
			letter := string(rune('A' + i%26))
			sb.WriteString(fmt.Sprintf("- **%s** host=`%s` kind=`%s` at=`%s`",
				letter, me.host, me.kind, me.at))
			if me.metric != "" {
				sb.WriteString(fmt.Sprintf(" metric=`%s`", me.metric))
			}
			sb.WriteString("\n")
			if me.blurb != "" {
				sb.WriteString(fmt.Sprintf("  > %s\n", me.blurb))
			}
		}
	}
	sb.WriteString("\n")

	// H2: Diagnoses (if any).
	var diags []Finding
	for _, f := range findings {
		if f["kind"] == "diagnosis_candidate" {
			diags = append(diags, f)
		}
	}
	if len(diags) > 0 {
		sb.WriteString("## Diagnoses\n\n")
		for _, d := range diags {
			name, _ := d["pattern_name"].(string)
			rationale, _ := d["rationale"].(string)
			cr, _ := d["confidence_rank"].(float64)
			sb.WriteString(fmt.Sprintf("### %s (confidence %.2f)\n\n%s\n\n", name, cr, rationale))
		}
	}

	// H2: Per-host summary.
	sb.WriteString("## Per-host summary\n\n")
	for _, h := range hosts {
		sb.WriteString(fmt.Sprintf("### %s\n\n", h.name))
		var letters []string
		for i, me := range markers {
			if me.host == h.name {
				letters = append(letters, string(rune('A'+i%26)))
			}
		}
		if len(letters) > 0 {
			sb.WriteString(fmt.Sprintf("Markers on this host: %s\n\n", strings.Join(letters, ", ")))
		} else {
			sb.WriteString("_No markers on this host._\n\n")
		}
	}

	return sb.String()
}

// ── NDJSON passthrough ────────────────────────────────────────────────────────

func buildFindingsNDJSON(lines []json.RawMessage) []byte {
	var buf bytes.Buffer
	buf.Grow(len(lines) * 256)
	for _, line := range lines {
		buf.Write(line)
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

// ── Bundle README / launcher scripts ──────────────────────────────────────────

// bundleReadme is the human-facing README placed at the top of the bundle.
// It walks the recipient through opening the bundle in t2 on either macOS or
// Linux, and points at the sidecar analysis.md / findings.ndjson for the
// analyst writeup and raw findings.
func bundleReadme() string {
	return `FTDC analysis bundle — self-contained package.

Contents:
  context.t2       t2 state file (loaded by the launcher; do not
                   open directly in t2 — relative paths won't resolve).
  analysis.md      analyst writeup: marker glossary, diagnoses,
                   detailed narrative. Shareable with other engineers.
  findings.ndjson  full Finding stream produced by ftdc_parser --auto.
  ftdc/<host>/     the FTDC files referenced by context.t2.
  open-in-t2.sh    launcher (macOS and Linux). Absolutizes the relative
                   FTDC paths in context.t2 to a temp state file, then
                   launches t2 with --resume so t2 treats the file as a
                   state, not a data file.

How to open in t2:

  macOS:  Open Terminal, cd into this directory, then run:
            bash open-in-t2.sh
          (Use 'bash open-in-t2.sh', not './open-in-t2.sh' — macOS
          quarantines downloaded executables from archives. Invoking
          via bash uses the system-trusted /bin/bash and bypasses the
          quarantine check entirely.)

  Linux:  cd into this directory and run:
            bash open-in-t2.sh
          or:
            ./open-in-t2.sh
          Both require t2 and python3 on PATH.

  Manual (advanced): rewrite the relative ftdc/* paths in context.t2 to
          absolute paths, save as e.g. context.absolute.t2, then run:
            t2 --resume context.absolute.t2          (Linux)
            open -a t2 --args --resume context.absolute.t2   (macOS)
          The launcher script automates exactly this. Do NOT pass
          context.t2 positionally to t2 — without --resume t2 sniffs
          it as a data file, fails the sniffer, and panics in the log
          reader.

How to share this bundle:
  Send the .tgz archive (created alongside this directory by
  ftdc_parser --auto-emit-t2). The recipient extracts it anywhere on
  disk and follows the macOS / Linux instructions above. The bundle is
  fully self-contained — no other files are needed.

Generated by: ftdc_parser --auto-emit-t2
`
}

// bundleOpenScript is the cross-platform launcher (open-in-t2.sh).
//
// Why this script exists: context.t2 references FTDC via relative paths.
// Stock t2 resolves those against its process cwd, which on macOS is "/"
// when t2.app is launched via `open -a` (launchd ignores the shell's cwd).
// So simply `cd`ing before `open` does NOT help on macOS. The launcher
// instead REWRITES the relative paths in context.t2 to absolute paths
// (based on the bundle's runtime location) into a temp file, then opens
// that — which makes the bundle portable: the recipient can extract the
// .tgz anywhere on disk and the launcher resolves paths correctly without
// needing the optional t2 decodePath patch.
//
// Extension choice (.sh, not .command):
//   .command is the macOS Finder double-click convention but is flagged as
//   potential malware by EDR tools when extracted from a downloaded archive
//   (quarantine xattr + executable shell script = suspicious pattern).
//   .sh is ignored by Finder (opens in a text editor on double-click),
//   so users must run it explicitly. Explicit `bash open-in-t2.sh` bypasses
//   the quarantine-bit check entirely: the OS audits /bin/bash (trusted),
//   not the script. Both macOS and Linux use the same invocation.
//
// Dependency: python3 (preinstalled on macOS with Xcode CLT and on every
// modern Linux distro). Used only to rewrite a JSON file portably; no
// other Python is invoked.
func bundleOpenScript() string {
	return `#!/bin/bash
# Resolve absolute paths in context.t2 to the bundle's current location,
# then launch t2 on the rewritten file. Works regardless of whether t2 is
# launched via macOS launchd ('open -a') or Linux exec — both see absolute
# paths.
set -e
BUNDLE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$BUNDLE_DIR"
if [ ! -f context.t2 ]; then
  echo "context.t2 not found in $BUNDLE_DIR" >&2
  exit 1
fi
if ! command -v python3 >/dev/null 2>&1; then
  echo "python3 not found in PATH; cannot rewrite paths in context.t2." >&2
  echo "Install python3 and retry, or open context.t2 manually after " >&2
  echo "launching t2 from $BUNDLE_DIR." >&2
  exit 1
fi

# Rewrite relative ftdc/* paths in context.t2 to absolute paths under
# $BUNDLE_DIR. The result is written to a temp file so the original is
# left untouched (the bundle stays movable / re-sharable).
TMP_CONTEXT="$(mktemp -t t2_context.XXXXXX)"
mv "$TMP_CONTEXT" "${TMP_CONTEXT}.t2"
TMP_CONTEXT="${TMP_CONTEXT}.t2"
export BUNDLE_DIR TMP_CONTEXT
python3 - <<'PY'
import json, os, sys
bundle = os.environ["BUNDLE_DIR"]
out = os.environ["TMP_CONTEXT"]
with open(os.path.join(bundle, "context.t2")) as f:
    d = json.load(f)
def absolutize(p):
    return p if os.path.isabs(p) else os.path.join(bundle, p)
for ctx in d.get("contexts", []):
    ds = ctx.get("dataset", {})
    if "files" in ds:
        ds["files"] = [absolutize(p) for p in ds["files"]]
    if "userTags" in ds:
        ds["userTags"] = {absolutize(k): v for k, v in ds["userTags"].items()}
with open(out, "w") as f:
    json.dump(d, f, indent=2)
PY

case "$(uname -s)" in
  Darwin)
    # macOS: 'open -a t2 --args --resume <file>' passes --resume to t2.
    # Without --args, t2 treats the path as a positional data file (.t2 is
    # not auto-recognized as a state file from argv), sniff-fails on csv/
    # ftdc/iostat, then panics in the log reader.
    exec /usr/bin/open -a t2 --args --resume "$TMP_CONTEXT"
    ;;
  *)
    if command -v t2 >/dev/null 2>&1; then
      exec t2 --resume "$TMP_CONTEXT"
    fi
    echo "t2 not found in PATH." >&2
    echo "Launch t2 manually with: t2 --resume $TMP_CONTEXT" >&2
    exit 1
    ;;
esac
`
}

// ── Tarball ───────────────────────────────────────────────────────────────────

func createTarball(bundleDir, tgzPath string) error {
	f, err := os.Create(tgzPath)
	if err != nil {
		return err
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	bundleName := filepath.Base(bundleDir)
	return filepath.WalkDir(bundleDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(bundleDir, path)
		tarPath := filepath.Join(bundleName, rel)

		// Resolve symlinks so the tarball is self-contained (important when
		// --auto-emit-t2-link-ftdc is set: the receiver won't have the originals).
		info, err := os.Stat(path) // follows symlinks
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = tarPath
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()
		_, err = io.Copy(tw, src)
		return err
	})
}
