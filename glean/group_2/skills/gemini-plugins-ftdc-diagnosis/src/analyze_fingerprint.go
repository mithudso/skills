// analyze_fingerprint.go — P4.4 local fingerprint library + self-doubt.
//
// Persists a per-capture "shape fingerprint" — a set of canonicalized
// Finding tuples (kind, direction or pattern, magnitude bucket) — to
//   ~/claude-dump/ftdc-fingerprints/<input_sha256>.json
// On every new --auto run, scans the library and emits
//   kind:"similar_prior_capture"
// for past captures with high Jaccard similarity on the shape set.
//
// This is solo + FTDC-only: the corpus is Gregory's own past captures.
// No JIRA, no embeddings, no hosted service.
//
// Stored-inference self-doubt:
//   - Files older than 90 days are quarantined unless parser_version +
//     pattern_catalog_hash match exactly.
//   - Confidence-rank ≥ 0.5 required at write time for diagnosis-shape
//     entries to be considered "authoritative" matches; entries from
//     low-confidence runs are surfaced as informational only.
//   - Every read emits kind:"stored_inference_review" describing the
//     verification performed.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// fingerprintStoreDir returns where per-capture shape fingerprints live.
func fingerprintStoreDir() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		home = "/tmp"
	}
	return filepath.Join(home, "claude-dump", "ftdc-fingerprints")
}

// fingerprintStaleAge mirrors the stored-inference self-doubt window.
const fingerprintStaleAge = 90 * 24 * time.Hour

// fingerprintTopK is the maximum number of nearest-neighbor matches
// emitted per run. Three is plenty for a solo user; surfacing more
// just adds noise.
const fingerprintTopK = 3

// fingerprintMinJaccard suppresses very weak matches. Below this
// similarity score, the alleged "match" is mostly random Finding-set
// overlap rather than real shape resemblance.
const fingerprintMinJaccard = 0.30

// fingerprintEntry is the stored file's top-level shape.
type fingerprintEntry struct {
	Provenance   fingerprintProvenance `json:"provenance"`
	ShapeSet     []string              `json:"shape_set"`
	ShapeDigest  string                `json:"shape_digest,omitempty"`
	NotesPath    string                `json:"notes_path,omitempty"`
}

type fingerprintProvenance struct {
	WrittenAt           time.Time `json:"written_at"`
	ParserVersion       string    `json:"parser_version"`
	PatternCatalogHash  string    `json:"pattern_catalog_hash"`
	InputSHA256         string    `json:"input_sha256"`
	InputPathBasename   string    `json:"input_path_basename,omitempty"`
	MaxConfidenceAtWrite float64  `json:"max_confidence_at_write"`
}

// emitFingerprintMatches reads run_metadata to get input_sha256, builds
// the current run's shape set, searches the library for neighbors, and
// emits similar_prior_capture Findings for matches. At end, writes the
// current fingerprint to the library file.
func emitFingerprintMatches(c *findingCollector, host string) (emitted int, skipped int) {
	var inputSHA256, inputBasename, parserVer string
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		if k != "run_metadata" {
			continue
		}
		if s, ok := f["input_sha256"].(string); ok {
			inputSHA256 = s
		}
		if s, ok := f["input_path_basename"].(string); ok {
			inputBasename = s
		}
		if lv, ok := f["library_versions"].(map[string]interface{}); ok {
			if pv, ok := lv["parser_version"].(string); ok {
				parserVer = pv
			}
		}
		break
	}
	if inputSHA256 == "" {
		return 0, 1
	}
	catalogHash := patternCatalogHash()
	current := buildCurrentFingerprint(c, host, parserVer, catalogHash, inputSHA256, inputBasename)

	// Search the library.
	storeDir := fingerprintStoreDir()
	matches, quarantined := searchFingerprintLibrary(storeDir, current, inputSHA256, parserVer, catalogHash)
	for _, m := range matches {
		c.add(newFinding("similar_prior_capture", host, m.fileBasename, map[string]interface{}{
			"prior_input_sha256":     m.entry.Provenance.InputSHA256,
			"prior_input_basename":   m.entry.Provenance.InputPathBasename,
			"prior_written_at":       m.entry.Provenance.WrittenAt.UTC().Format(time.RFC3339),
			"prior_parser_version":   m.entry.Provenance.ParserVersion,
			"prior_max_confidence":   m.entry.Provenance.MaxConfidenceAtWrite,
			"jaccard_similarity":     roundFloat(m.jaccard, 4),
			"matching_shape_tuples":  m.matchingShapes,
			"prior_fingerprint_path": m.filePath,
			"evidence_basis":         "prior_captures",
			"interpretation": "This past capture shares a high-similarity Finding-shape set with the current capture. Advisory — review the prior capture's notes/recap before assuming the same root cause.",
		}))
		emitted++
	}
	// Self-doubt audit Findings for any quarantined library entries.
	for _, q := range quarantined {
		c.add(newFinding("stored_inference_review", host, q.fileBasename+"|"+q.reason, map[string]interface{}{
			"action":              "neighbor_quarantined",
			"quarantined_reason":  q.reason,
			"file":                q.filePath,
			"entry_written_at":    q.writtenAt,
			"entry_parser_version": q.parserVer,
		}))
		emitted++
	}

	// Write/update the current fingerprint.
	_ = writeFingerprint(storeDir, current)
	return emitted, 0
}

// buildCurrentFingerprint canonicalizes the current run's Findings into
// a sorted shape-tuple set + provenance.
func buildCurrentFingerprint(c *findingCollector, host, parserVer, catalogHash, inputSHA256, inputBasename string) fingerprintEntry {
	tuples := map[string]bool{}
	maxConf := 0.0
	for _, f := range c.items {
		k, _ := f["kind"].(string)
		hf, _ := f["host"].(string)
		if hf != "" && hf != host {
			continue
		}
		switch k {
		case "mover":
			dir, _ := f["direction"].(string)
			bucket := magnitudeBucketFor(absFloatField(f["cohens_d"]))
			// Don't include metric path — that would prevent matching
			// across different captures with similar shapes but slightly
			// different metric names.
			tuples[k+":"+dir+":"+bucket] = true
		case "diagnosis_candidate", "diagnosis_uncertainty":
			pattern, _ := f["pattern_name"].(string)
			cr, _ := f["confidence_rank"].(float64)
			bucket := confidenceBucketFor(cr)
			tuples[k+"::"+pattern+":"+bucket] = true
			if cr > maxConf {
				maxConf = cr
			}
		case KindSealEvent, KindSealEventCandidate:
			cr, _ := f["confidence_rank"].(float64)
			tuples[k+":"+confidenceBucketFor(cr)] = true
		case "role_asymmetry", "node_pair_asymmetry":
			dir, _ := f["direction"].(string)
			bucket := magnitudeBucketFor(absFloatField(f["cohens_d"]))
			tuples[k+":"+dir+":"+bucket] = true
		case "joint_changepoint":
			size, _ := toIntField(f["cluster_size"])
			tuples[k+":size_bucket:"+sizeBucketFor(size)] = true
		case "checkpoint_install_anomaly", "materialization_lag", "io_routing", "election_storm",
			"checkpoint_stall", "eviction_divergence", "apply_saturation", "capacity_pressure",
			"near_ceiling", "memory_triplet", "counter_inversion", "counter_wrap",
			KindTermChange, "coverage_gap", "host_pathology":
			// Kinds whose mere presence is the signal.
			tuples[k] = true
		}
	}
	if len(tuples) == 0 {
		return fingerprintEntry{Provenance: fingerprintProvenance{
			WrittenAt:          time.Now().UTC(),
			ParserVersion:      parserVer,
			PatternCatalogHash: catalogHash,
			InputSHA256:        inputSHA256,
			InputPathBasename:  inputBasename,
		}}
	}
	shapes := make([]string, 0, len(tuples))
	for t := range tuples {
		shapes = append(shapes, t)
	}
	sort.Strings(shapes)
	return fingerprintEntry{
		Provenance: fingerprintProvenance{
			WrittenAt:            time.Now().UTC(),
			ParserVersion:        parserVer,
			PatternCatalogHash:   catalogHash,
			InputSHA256:          inputSHA256,
			InputPathBasename:    inputBasename,
			MaxConfidenceAtWrite: maxConf,
		},
		ShapeSet: shapes,
	}
}

// magnitudeBucketFor coarsens a numeric magnitude so two captures with
// "similar but not identical" effect sizes match. Buckets: low (<0.5),
// medium ([0.5, 1.0)), high ([1.0, 2.0)), very_high (>=2.0).
func magnitudeBucketFor(x float64) string {
	x = math.Abs(x)
	switch {
	case x < 0.5:
		return "low"
	case x < 1.0:
		return "medium"
	case x < 2.0:
		return "high"
	default:
		return "very_high"
	}
}

// confidenceBucketFor coarsens a confidence_rank.
func confidenceBucketFor(c float64) string {
	switch {
	case c >= 0.9:
		return "strong"
	case c >= 0.5:
		return "plausible"
	default:
		return "weak"
	}
}

// sizeBucketFor coarsens a cluster size.
func sizeBucketFor(n int) string {
	switch {
	case n < 5:
		return "small"
	case n < 15:
		return "medium"
	default:
		return "large"
	}
}

// absFloatField returns |x| from a Finding-map field value (interface{}),
// treating an absent/non-numeric field as zero. Companion to absFloat in
// analyze_correlation.go which takes a typed float64 — different call
// sites, different signatures.
func absFloatField(v interface{}) float64 {
	x, _ := v.(float64)
	return math.Abs(x)
}

// fingerprintMatch represents one nearest neighbor.
type fingerprintMatch struct {
	jaccard        float64
	matchingShapes []string
	entry          fingerprintEntry
	fileBasename   string
	filePath       string
}

type fingerprintQuarantine struct {
	reason       string
	fileBasename string
	filePath     string
	writtenAt    string
	parserVer    string
}

// searchFingerprintLibrary scans the library dir, computes Jaccard
// similarity against the current shape set, applies self-doubt
// filters, and returns top-K matches + quarantine summary.
func searchFingerprintLibrary(dir string, current fingerprintEntry, currentSHA, currentParserVer, currentCatalogHash string) (matches []fingerprintMatch, quarantined []fingerprintQuarantine) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil
	}
	currentSet := map[string]bool{}
	for _, s := range current.ShapeSet {
		currentSet[s] = true
	}
	now := time.Now()
	var pool []fingerprintMatch
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		// Skip the current capture's own file — comparing against self
		// is uninteresting.
		if strings.HasPrefix(e.Name(), currentSHA) {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var entry fingerprintEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}
		// Reject entries with zero written_at: now.Sub(zero) exceeds
		// fingerprintStaleAge and could pass the version check silently.
		if entry.Provenance.WrittenAt.IsZero() {
			quarantined = append(quarantined, fingerprintQuarantine{
				reason:       "missing_written_at",
				fileBasename: e.Name(),
				filePath:     path,
				parserVer:    entry.Provenance.ParserVersion,
			})
			continue
		}
		writtenAtStr := entry.Provenance.WrittenAt.UTC().Format(time.RFC3339)
		if now.Sub(entry.Provenance.WrittenAt) > fingerprintStaleAge {
			if entry.Provenance.ParserVersion != currentParserVer || entry.Provenance.PatternCatalogHash != currentCatalogHash {
				quarantined = append(quarantined, fingerprintQuarantine{
					reason:       "stale_with_version_mismatch",
					fileBasename: e.Name(),
					filePath:     path,
					writtenAt:    writtenAtStr,
					parserVer:    entry.Provenance.ParserVersion,
				})
				continue
			}
		}
		// Jaccard.
		match := 0
		var matchedShapes []string
		for _, s := range entry.ShapeSet {
			if currentSet[s] {
				match++
				if len(matchedShapes) < 12 {
					matchedShapes = append(matchedShapes, s)
				}
			}
		}
		union := len(currentSet) + len(entry.ShapeSet) - match
		if union == 0 {
			continue
		}
		j := float64(match) / float64(union)
		if j < fingerprintMinJaccard {
			continue
		}
		pool = append(pool, fingerprintMatch{
			jaccard:        j,
			matchingShapes: matchedShapes,
			entry:          entry,
			fileBasename:   e.Name(),
			filePath:       path,
		})
	}
	sort.Slice(pool, func(i, j int) bool {
		if pool[i].jaccard != pool[j].jaccard {
			return pool[i].jaccard > pool[j].jaccard
		}
		return pool[i].entry.Provenance.WrittenAt.After(pool[j].entry.Provenance.WrittenAt)
	})
	if len(pool) > fingerprintTopK {
		pool = pool[:fingerprintTopK]
	}
	return pool, quarantined
}

// writeFingerprint serializes the current entry into the library.
func writeFingerprint(dir string, entry fingerprintEntry) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if entry.Provenance.InputSHA256 == "" {
		return nil
	}
	path := filepath.Join(dir, entry.Provenance.InputSHA256+".json")
	// Compute a deterministic tail digest of the shape set so we can
	// surface it later as a recurrence-fingerprint analog.
	if len(entry.ShapeSet) > 0 {
		h := sha256.Sum256([]byte(strings.Join(entry.ShapeSet, "|")))
		entry.ShapeDigest = hex.EncodeToString(h[:8])
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
