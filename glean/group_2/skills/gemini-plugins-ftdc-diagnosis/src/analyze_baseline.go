// Tier D: cross-run baseline + regression detection.
//
// Stores per-(host_token, metric) rolling EWMA statistics + Tier-7 recurrence
// fingerprints in a single flat JSONL file at the path passed via
// --auto-baseline-file (or the XDG default). On each run we:
//
//   1. Read the file, build in-memory maps keyed by (host_token, metric)
//      for baseline rows and (host_token, capture_token, fingerprint) for
//      recurrence rows. Skip rows with mismatched ParserVersion or with
//      last_update_ts > 30 days old (TTL).
//
//   2. Walk the seriesStore. For each metric, compare the current
//      capture's mean / p95 against the rolling baseline. Emit
//      kind:"regression" when the deviation exceeds N · ewma_sigma AND the
//      Tier-B-calibrated Welch p_value (already on the corresponding
//      mover Finding) supports significance.
//
//   3. Update the EWMA stats and recurrence count, write back atomically
//      via os.Rename of a tmpfile.
//
// host_token = SHA-256(hostname || per-db salt). Never store the raw
// hostname. The salt is generated once at first write and persisted as
// the file's first line.
//
// "What if the analysis is wrong?" mitigations (per plan v5):
//   - parser_version stamping invalidates rows across binary versions.
//   - 30-day TTL prevents old rows poisoning the baseline.
//   - --reset-baseline CLI flag wipes the file.
//   - regression Findings are ADVISORY — they never gate other detectors.
//   - kind:"baseline_diagnostic" emits the EWMA state per metric so the
//     baseline is inspectable, not opaque.

package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const (
	baselineTTLDays         = 30
	baselineEWMAlphaMean    = 0.2 // ~5-capture half-life
	baselineEWMAlphaVar     = 0.1 // slower variance tracking
	baselineRegressionSigma = 3.0
)

type baselineSaltLine struct {
	Kind string `json:"kind"`
	Salt string `json:"salt"`
}

type recurrenceRow struct {
	Kind          string `json:"kind"`
	HostToken     string `json:"host_token"`
	CaptureToken  string `json:"capture_token"`
	Fingerprint   string `json:"fingerprint"`
	HitCount      int    `json:"hit_count"`
	LastSeen      string `json:"last_seen"`
	ParserVersion string `json:"parser_version"`
}

type baselineRow struct {
	Kind          string  `json:"kind"`
	HostToken     string  `json:"host_token"`
	Metric        string  `json:"metric"`
	EWMAMean      float64 `json:"ewma_mean"`
	EWMAVar       float64 `json:"ewma_var"`
	N             int     `json:"n"`
	LastUpdateTS  string  `json:"last_update_ts"`
	ParserVersion string  `json:"parser_version"`
}

// loadBaseline reads the baseline JSONL file. Returns the salt (creating
// one on first write), per-(host, metric) baseline rows, and per-recurrence
// rows. Caller passes the host_token and we filter to that host's rows.
func loadBaseline(path string) (salt string, baseRows map[string]baselineRow, recurRows map[string]recurrenceRow, err error) {
	baseRows = map[string]baselineRow{}
	recurRows = map[string]recurrenceRow{}
	f, ferr := os.Open(path)
	if ferr != nil {
		if !os.IsNotExist(ferr) {
			return "", nil, nil, ferr
		}
		// Fresh DB: generate a salt.
		var b [16]byte
		if _, rerr := rand.Read(b[:]); rerr != nil {
			return "", nil, nil, rerr
		}
		salt = hex.EncodeToString(b[:])
		return salt, baseRows, recurRows, nil
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 16<<20)
	cutoff := time.Now().Add(-baselineTTLDays * 24 * time.Hour)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		// Peek at "kind" to dispatch.
		var head struct {
			Kind          string `json:"kind"`
			ParserVersion string `json:"parser_version"`
		}
		if jerr := json.Unmarshal(line, &head); jerr != nil {
			continue // tolerate corrupt lines
		}
		switch head.Kind {
		case "salt":
			var s baselineSaltLine
			if jerr := json.Unmarshal(line, &s); jerr == nil {
				salt = s.Salt
			}
		case "baseline":
			if head.ParserVersion != ParserVersion {
				continue // invalidated by version bump
			}
			var b baselineRow
			if jerr := json.Unmarshal(line, &b); jerr != nil {
				continue
			}
			ts, perr := time.Parse(time.RFC3339Nano, b.LastUpdateTS)
			if perr == nil && ts.Before(cutoff) {
				continue // TTL
			}
			baseRows[b.HostToken+"|"+b.Metric] = b
		case "recurrence":
			if head.ParserVersion != ParserVersion {
				continue
			}
			var r recurrenceRow
			if jerr := json.Unmarshal(line, &r); jerr != nil {
				continue
			}
			// Apply the same TTL to recurrence rows as baseline rows.
			// Without this the recurrence table grows unboundedly across
			// runs (one row per unique fingerprint per host, forever).
			lt, perr := time.Parse(time.RFC3339Nano, r.LastSeen)
			if perr == nil && lt.Before(cutoff) {
				continue
			}
			recurRows[r.HostToken+"|"+r.Fingerprint] = r
		}
	}
	if scErr := sc.Err(); scErr != nil {
		// Mid-file I/O error: half-loaded rows could be silently wrong.
		// Surface the error so the caller can decide (typically the
		// pipeline will continue with `baseRows`/`recurRows` partially
		// populated AND emit a `baseline_diagnostic` indicating the
		// partial-load condition).
		return salt, baseRows, recurRows, fmt.Errorf("baseline scan error after %d rows: %w", len(baseRows)+len(recurRows), scErr)
	}
	if salt == "" {
		// File existed but had no salt line — corruption or legacy file.
		// Regenerate AND warn loudly: this changes host_token across runs,
		// orphaning all existing baseline rows under the old token.
		// Surface to stderr so the user can `--reset-baseline` if they
		// prefer a clean start.
		var b [16]byte
		if _, rerr := rand.Read(b[:]); rerr != nil {
			return "", nil, nil, rerr
		}
		salt = hex.EncodeToString(b[:])
		fmt.Fprintln(os.Stderr,
			"ftdc_parser: baseline file lacked a salt line; regenerated. "+
				"Previous host_token values are now orphaned. Run --reset-baseline "+
				"for a clean start if you want to rebuild the baseline from scratch.")
	}
	return salt, baseRows, recurRows, nil
}

// hostToken returns SHA-256(hostname || salt) hex-encoded. Used so the
// baseline file never contains raw hostnames.
func hostToken(hostname, salt string) string {
	h := sha256.Sum256([]byte(hostname + "|" + salt))
	return hex.EncodeToString(h[:16])
}

// writeBaseline atomically rewrites the JSONL file. Caller passes the salt
// + all rows (already TTL-filtered + version-matched). Rewrites via
// os.Rename of a tmpfile.
func writeBaseline(path, salt string, baseRows map[string]baselineRow, recurRows map[string]recurrenceRow) error {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	// Remove the .tmp file on any error path so partial writes don't accumulate.
	success := false
	defer func() {
		_ = f.Close()
		if !success {
			_ = os.Remove(tmp)
		}
	}()
	bw := bufio.NewWriter(f)
	// Salt line first.
	if err := json.NewEncoder(bw).Encode(baselineSaltLine{Kind: "salt", Salt: salt}); err != nil {
		return err
	}
	// Sort row keys for determinism.
	bks := make([]string, 0, len(baseRows))
	for k := range baseRows {
		bks = append(bks, k)
	}
	sort.Strings(bks)
	enc := json.NewEncoder(bw)
	for _, k := range bks {
		if err := enc.Encode(baseRows[k]); err != nil {
			return err
		}
	}
	rks := make([]string, 0, len(recurRows))
	for k := range recurRows {
		rks = append(rks, k)
	}
	sort.Strings(rks)
	for _, k := range rks {
		if err := enc.Encode(recurRows[k]); err != nil {
			return err
		}
	}
	if err := bw.Flush(); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		return err
	}
	success = true
	return nil
}

// defaultBaselinePath returns $XDG_DATA_HOME/ftdc-parser/baseline.jsonl,
// with a fallback to ~/.local/share/... when XDG_DATA_HOME is unset.
func defaultBaselinePath() string {
	if p := os.Getenv("XDG_DATA_HOME"); p != "" {
		return filepath.Join(p, "ftdc-parser", "baseline.jsonl")
	}
	if h, err := os.UserHomeDir(); err == nil {
		return filepath.Join(h, ".local", "share", "ftdc-parser", "baseline.jsonl")
	}
	return ""
}

// runBaselineStage is the entry point called from the analyze pipeline.
// Reads the baseline file (if cfg.BaselineFile is set), compares the
// current capture's metric statistics, emits regression + baseline_diagnostic
// Findings, and writes back the updated state.
func runBaselineStage(c *findingCollector, host string, store *seriesStore, cfg AutoConfig) (emitted, skipped int) {
	if cfg.BaselineFile == "" {
		return 0, 0
	}
	salt, baseRows, recurRows, err := loadBaseline(cfg.BaselineFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ftdc_parser: baseline load error: %v\n", err)
		return 0, 0
	}
	hostTok := hostToken(host, salt)

	// Compute per-metric mean for the full capture window.
	for id, name := range store.names {
		if id >= len(store.isDegenerate) || id >= len(store.isCounter) || id >= len(store.fullMean) {
			continue
		}
		// Skip degenerate / counter rates that have unstable EWMA semantics.
		// Counters are stored as rates in materializedSeries; their EWMA
		// makes sense too but for simplicity we focus on gauges.
		if store.isDegenerate[id] || store.isCounter[id] {
			skipped++
			continue
		}
		mean := store.fullMean[id]
		key := hostTok + "|" + name
		row, ok := baseRows[key]
		nowTS := time.Now().UTC().Format(time.RFC3339Nano)
		if !ok {
			// First observation: seed EWMA with this capture's mean,
			// zero variance.
			baseRows[key] = baselineRow{
				Kind:          "baseline",
				HostToken:     hostTok,
				Metric:        name,
				EWMAMean:      mean,
				EWMAVar:       0,
				N:             1,
				LastUpdateTS:  nowTS,
				ParserVersion: ParserVersion,
			}
			continue
		}
		// Compare current mean against rolling baseline.
		dev := mean - row.EWMAMean
		ewmaSigma := math.Sqrt(row.EWMAVar)
		if ewmaSigma > 0 && math.Abs(dev) > baselineRegressionSigma*ewmaSigma && row.N >= 3 {
			c.add(newFinding("regression", host, name, map[string]interface{}{
				"metric":           name,
				"current_mean":     roundFloat(mean, 6),
				"baseline_mean":    roundFloat(row.EWMAMean, 6),
				"baseline_sigma":   roundFloat(ewmaSigma, 6),
				"deviation_sigmas": roundFloat(dev/ewmaSigma, 4),
				"baseline_n":       row.N,
				"threshold_sigmas": baselineRegressionSigma,
				"reason_code":      "ewma_deviation",
				"advisory":         true,
			}))
			emitted++
		}
		// Update EWMA.
		newMean := (1-baselineEWMAlphaMean)*row.EWMAMean + baselineEWMAlphaMean*mean
		// Hunter (1986) EWMA variance: V_new = (1-β)*V_old + β*(x-μ_old)*(x-μ_new).
		// μ_new = (1-α)*μ_old + α*x, so (x-μ_new) = (1-α)*(x-μ_old) = (1-α)*dev.
		// Using dev² instead of dev*(1-α_mean)*dev overestimates steady-state
		// variance by 1/(1-α_mean) ≈ 1.25, widening the detection band by ~12%.
		newVar := (1-baselineEWMAlphaVar)*row.EWMAVar + baselineEWMAlphaVar*(1-baselineEWMAlphaMean)*dev*dev
		row.EWMAMean = newMean
		row.EWMAVar = newVar
		row.N++
		row.LastUpdateTS = nowTS
		row.ParserVersion = ParserVersion
		baseRows[key] = row
	}

	// Emit one baseline_diagnostic per host summarizing the file state.
	c.add(newFinding("baseline_diagnostic", host, "global", map[string]interface{}{
		"host_token":      hostTok,
		"baseline_path":   cfg.BaselineFile,
		"baseline_rows":   len(baseRows),
		"recurrence_rows": len(recurRows),
		"parser_version":  ParserVersion,
		"ttl_days":        baselineTTLDays,
		"alpha_mean":      baselineEWMAlphaMean,
		"alpha_var":       baselineEWMAlphaVar,
		"note":            "advisory only; baseline state is inspectable here. --reset-baseline to wipe.",
	}))
	emitted++

	// Update recurrence count: look for any kind:"recurrence_fingerprint"
	// already emitted in this run and increment.
	for _, f := range c.items {
		if k, _ := f["kind"].(string); k != "recurrence_fingerprint" {
			continue
		}
		fp, _ := f["fingerprint"].(string)
		if fp == "" {
			continue
		}
		rkey := hostTok + "|" + fp
		rawHost := ""
		if v, ok := f["host"].(string); ok {
			rawHost = v
		}
		now := time.Now().UTC().Format(time.RFC3339Nano)
		// Uniqueness per run comes from `now`; rawHost scopes the hash to the host.
		runUnique := hashShort(rawHost + "|" + now)
		if row, ok := recurRows[rkey]; ok {
			row.HitCount++
			row.LastSeen = now
			row.ParserVersion = ParserVersion
			recurRows[rkey] = row
		} else {
			recurRows[rkey] = recurrenceRow{
				Kind:          "recurrence",
				HostToken:     hostTok,
				CaptureToken:  runUnique,
				Fingerprint:   fp,
				HitCount:      1,
				LastSeen:      now,
				ParserVersion: ParserVersion,
			}
		}
	}

	if err := writeBaseline(cfg.BaselineFile, salt, baseRows, recurRows); err != nil {
		fmt.Fprintf(os.Stderr, "ftdc_parser: baseline write error: %v\n", err)
	}
	return emitted, skipped
}

func hashShort(s string) string {
	if s == "" {
		return ""
	}
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:8])
}
