// analyze.go - --auto mode orchestration and Finding schema.
//
// Emits structured Findings NDJSON describing the shape of an FTDC capture:
// movers, change-points, correlations, histogram percentiles, sample gaps,
// schema fingerprint, run health. No semantic tagging (no "severity" / "cause"
// / "role" / named-bug citations) — the analyzer skill's LLM + markdown
// catalog does interpretation.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	// SchemaVersion bumps only on breaking change to Finding shape.
	// Additive field changes do not bump.
	SchemaVersion = "1"

	// Default knobs. Each is overridable by a --auto-* flag and recorded in
	// the run_metadata Finding so the LLM downstream sees the choice.
	defaultAutoTopK       = 50
	defaultAutoDriftRatio = 5.0
	// Baseline / stress are full HALVES of the capture so every sample is in
	// exactly one window. The earlier 25%/25% split was blind to events
	// in the middle 50% of long captures — a mid-capture latency creep
	// would be excluded from both windows and miss Welch-score-driven
	// detection entirely. Full halves still catch end-of-capture shifts
	// (a shift in the last quarter affects the second-half mean) AND
	// mid-capture shifts (a shift partway through pulls the half means
	// apart).
	defaultBaselineFraction    = 0.50
	defaultStressFraction      = 0.50
	defaultGapThresholdSeconds = 5.0
	defaultCorrelationMinR     = 0.7
	defaultPELTMinSegment      = 60
	defaultMannKendallTauMin   = 0.30
	defaultMannKendallPMax     = 0.05
	defaultSpikeZ              = 5.0
)

// ParserVersion identifies which build of the parser produced an output
// stream. Bumped manually when emit semantics change (Finding kind set
// changes, schema-incompatible field additions, statistical-algorithm
// swaps). Surfaced in run_metadata.library_versions so the analyzer
// downstream can detect across-version comparisons. Also stamped onto
// rows in the Tier-D baseline file so a version bump invalidates stale
// entries.
const ParserVersion = "post-tier7+A"

// welfordWorkers is the pinned worker count for parallel-Welford
// reductions across all stats-bearing stages. Constant (not derived
// from NumCPU) so the reduction tree shape — and therefore the
// ULP-precise m2 sums — match across machines. See parallelWelfordWindow
// for the rationale.
const welfordWorkers = 8

// coreMetricsAlwaysAnalyzed is the list of metrics that ALWAYS get the
// full deep-analysis treatment (Mann-Kendall trend, ED-PELT change-point,
// variance-shift CUSUM, spike detection) regardless of their Stage-1 Welch
// shift score. A baseline-vs-stress mean comparison can miss mid-capture
// events on a metric whose endpoints look similar (e.g. a brief latency
// excursion that resolves before the stress window starts). Putting the
// canonical "did something break?" metrics here guarantees deep analysis
// runs on them unconditionally.
//
// Conservative list — names stable across recent MongoDB versions.
// Version-renamed metrics silently fall off the list; the analyzer skill
// catches drift via schema_fingerprint.
var coreMetricsAlwaysAnalyzed = []string{
	"serverStatus.opLatencies.reads.latency",
	"serverStatus.opLatencies.reads.ops",
	"serverStatus.opLatencies.writes.latency",
	"serverStatus.opLatencies.writes.ops",
	"serverStatus.opLatencies.commands.latency",
	"serverStatus.opLatencies.commands.ops",
	"serverStatus.opLatencies.transactions.latency",
	"serverStatus.opLatencies.transactions.ops",
	"serverStatus.queues.execution.write.out",
	"serverStatus.queues.execution.write.available",
	"serverStatus.queues.execution.read.out",
	"serverStatus.queues.execution.read.available",
	"serverStatus.connections.current",
	"serverStatus.tcmalloc.generic.current_allocated_bytes",
	"serverStatus.tcmalloc.generic.heap_size",
	"serverStatus.mem.resident",
	"serverStatus.mem.virtual",
	"serverStatus.wiredTiger.cache.bytes currently in the cache",
	"serverStatus.wiredTiger.cache.tracked dirty bytes in the cache",
	"serverStatus.wiredTiger.cache.pages evicted by application threads",
	"serverStatus.metrics.disagg.getPageRequest.numQueued",
	"serverStatus.metrics.disagg.paliHandleGet.paliHandleGetTotalLatencyUs",
	// waiters.replication grows +1/s when PALI is frozen because operations
	// queue up waiting for replication to complete. Welch scoring
	// under-rates this metric because it peaks mid-capture and returns to
	// zero (the metric's endpoints look similar), so it must be forced into
	// deep analysis unconditionally.
	"serverStatus.metrics.disagg.waiters.replication",
}

// tsFormat is the canonical 3-digit-ms UTC timestamp used throughout the parser.
const tsFormat = "2006-01-02T15:04:05.000Z"

// Finding kinds (stable strings; never renamed without bumping SchemaVersion).
const (
	KindRunMetadata              = "run_metadata"
	KindSchemaFingerprint        = "schema_fingerprint"
	KindMetricScoreTable         = "metric_score_table"
	KindMover                    = "mover"
	KindCorrelation              = "correlation"
	KindMonotonicTrend           = "monotonic_trend"
	KindTrendInconclusive        = "trend_inconclusive" // Tier B.6: short-capture refusal
	KindChangepoint              = "changepoint"
	KindChangepointOversegmented = "changepoint_oversegmented"
	KindHistogramPercentile      = "histogram_percentile"
	KindHistogramInvalid         = "histogram_invalid"
	KindGap                      = "gap"
	KindVarianceShift            = "variance_shift"
	KindSkipped                  = "skipped"
	KindHostOnlyMetrics          = "host_only_metrics"
	KindCrossHostDrift           = "cross_host_drift"
	KindCrossHostRoleInferred    = "cross_host_role_inferred" // Tier A.6: replaced by KindHostMetricSetAsymmetry; const kept to avoid breaking transitive references during the final-update migration window.
	KindHostMetricSetAsymmetry   = "host_metric_set_asymmetry"
	KindRunHealth                = "run_health"
	KindSpike                    = "spike"
	KindSuspiciousSmoothness     = "suspicious_smoothness"
	KindMissingSignal            = "missing_signal"
	KindChunkDecodeError         = "chunk_decode_error"
	KindCaptureQualityScore      = "capture_quality_score"
	KindApplySaturation          = "apply_saturation"
	KindElectionStorm            = "election_storm"
	KindCheckpointStall          = "checkpoint_stall"
	KindCounterStall             = "counter_stall"
	KindEvictionDivergence       = "eviction_divergence"
	KindTermChange               = "term_change"
	KindWorkloadSpec             = "workload_spec"
	KindSealEvent                = "seal_event"
	KindSealEventCandidate       = "seal_event_candidate"
)

// Reason codes for kind:"suspicious_smoothness" Findings (emitted by
// analyze_signal.go and consumed by analyze_capabilities.go).
const (
	ReasonCodeConstantZero       = "constant_zero"
	ReasonCodeVarianceTooLow     = "variance_too_low"
	ReasonCodePerfectlySteadyRate = "perfectly_steady_rate"
	ReasonCodeRoundNumberDensity  = "round_number_density"
)

// Reason codes for Tier-G detector Findings.
const (
	ReasonCodeMultipleTermChanges          = "multiple_term_changes_in_5min"
	ReasonCodeCheckpointExceededSyncdelay  = "checkpoint_exceeded_syncdelay_80pct"
	ReasonCodeReadIntoExceedsEvicted       = "read_into_exceeds_evicted"
)

// Finding is the structured fact unit. Heterogeneous fields by kind, so it's a
// map rather than a typed struct — but every Finding carries schema_version,
// kind, and id. host is included whenever single-host or multi-host context
// exists.
type Finding map[string]interface{}

// newFinding builds a Finding, fills mandatory fields, and computes a
// deterministic id. fields is merged on top of the mandatory ones; mandatory
// keys cannot be overridden.
func newFinding(kind, host string, idKey string, fields map[string]interface{}) Finding {
	f := Finding{
		"schema_version": SchemaVersion,
		"kind":           kind,
	}
	for k, v := range fields {
		if k == "schema_version" || k == "kind" || k == "id" {
			continue
		}
		f[k] = v
	}
	if host != "" {
		f["host"] = host
	}
	f["id"] = computeFindingID(kind, host, idKey)
	return f
}

// computeFindingID produces a short deterministic id like f-1a2b3c4d.
// idKey identifies the Finding within its kind (metric name + window for
// mover/changepoint/correlation; "global" for run_metadata; etc.).
func computeFindingID(kind, host, idKey string) string {
	h := sha256.New()
	fmt.Fprintf(h, "%s|%s|%s", kind, host, idKey)
	sum := h.Sum(nil)
	return "f-" + hex.EncodeToString(sum[:4])
}

// findingCollector gathers Findings produced during the pipeline, holds them
// in memory, then emits NDJSON sorted by (kind, metric or other stable key) at
// the end. Sorting is mandatory — the bias-control invariants forbid
// "interestingness" ordering by Go.
type findingCollector struct {
	items []Finding
}

func (c *findingCollector) add(f Finding) {
	c.items = append(c.items, f)
}

// emitAll writes Findings as NDJSON to w in deterministic order. The order key
// is (kind, secondary sort field per kind). Each Finding is one line of JSON.
// Map serialization in Go is non-deterministic in field order, so we marshal
// via a sorted-keys path.
func (c *findingCollector) emitAll(w io.Writer) error {
	sort.SliceStable(c.items, func(i, j int) bool {
		ki, _ := c.items[i]["kind"].(string)
		kj, _ := c.items[j]["kind"].(string)
		if ki != kj {
			return ki < kj
		}
		// Secondary: prefer a metric/host/follower field if present.
		si := stableSortKey(c.items[i])
		sj := stableSortKey(c.items[j])
		return si < sj
	})
	bw := w
	for _, f := range c.items {
		line, err := marshalFindingDeterministic(f)
		if err != nil {
			return err
		}
		if _, err := bw.Write(line); err != nil {
			return err
		}
		if _, err := bw.Write([]byte{'\n'}); err != nil {
			return err
		}
	}
	return nil
}

// stableSortKey picks a per-kind sort key. Falls back to id for kinds with no
// natural field.
func stableSortKey(f Finding) string {
	// Composite key: host + metric + leader + follower + at + window + segment_at.
	// All optional; missing fields become empty strings. The resulting key
	// totally orders Findings within a kind.
	var b strings.Builder
	for _, k := range []string{"host", "metric", "leader", "follower", "at", "segment_at"} {
		if v, ok := f[k]; ok {
			fmt.Fprintf(&b, "%v|", v)
		} else {
			b.WriteByte('|')
		}
	}
	// id is the tiebreaker.
	if v, ok := f["id"]; ok {
		fmt.Fprintf(&b, "%v", v)
	}
	return b.String()
}

// marshalFindingDeterministic serializes a Finding with sorted keys. json
// package's default map serialization sorts keys alphabetically, so we can
// rely on json.Marshal directly — but Go does NOT guarantee this for
// map[string]interface{} in all versions. As of go1.12+ it's documented to
// sort. We rely on that contract; if a future Go release breaks it, switch to
// a custom encoder.
func marshalFindingDeterministic(f Finding) ([]byte, error) {
	// Tier A.2: sanitize NaN / +Inf / -Inf before json.Marshal. The Go
	// json package refuses to encode these and returns an "unsupported
	// value" error, which would abort the whole emit. Replace each with
	// null and add a quality:"nan" annotation alongside so the analyzer
	// downstream can spot the substitution.
	return json.Marshal(map[string]interface{}(sanitizeFinding(f)))
}

// sanitizeFinding returns a copy of f with NaN/Inf float fields replaced by
// nil, recursing into nested map[string]interface{} payloads. Adds a top-
// level "quality": "nan" key if any substitution was performed.
func sanitizeFinding(f Finding) Finding {
	var dirty bool
	out := make(Finding, len(f))
	for k, v := range f {
		out[k] = sanitizeValue(v, &dirty)
	}
	if dirty {
		out["quality"] = "nan"
	}
	return out
}

func sanitizeValue(v interface{}, dirty *bool) interface{} {
	switch x := v.(type) {
	case float64:
		if math.IsNaN(x) || math.IsInf(x, 0) {
			*dirty = true
			return nil
		}
		return x
	case float32:
		f := float64(x)
		if math.IsNaN(f) || math.IsInf(f, 0) {
			*dirty = true
			return nil
		}
		return x
	case map[string]interface{}:
		m := make(map[string]interface{}, len(x))
		for k, val := range x {
			m[k] = sanitizeValue(val, dirty)
		}
		return m
	case []interface{}:
		a := make([]interface{}, len(x))
		for i, val := range x {
			a[i] = sanitizeValue(val, dirty)
		}
		return a
	case []string:
		// Deep-copy []string so downstream mutation in the caller can't
		// retroactively change the emitted Finding. Critical for fields
		// like supporting_finding_ids that Tier R emits.
		s := make([]string, len(x))
		copy(s, x)
		return s
	default:
		return v
	}
}

// AutoConfig carries all --auto-* flag values. Populated in main() before
// dispatch. Echoed verbatim in the run_metadata Finding so the LLM can audit
// every algorithmic choice.
type AutoConfig struct {
	// Mode and input
	Mode string // "single" | "multi" | "auto"

	// Window selection
	BaselineFrom time.Time
	BaselineTo   time.Time
	StressFrom   time.Time
	StressTo     time.Time
	HasBaseline  bool
	HasStress    bool

	// Output format toggles
	EmitOTel     bool
	EmitGolden   bool
	PrintSchema  bool
	BaselineFile string

	// Output budget + opt-in emission gates.
	BudgetTokens   int  // soft cap; 0 means unlimited
	EmitScoreTable bool // false = omit the ~190K-token metric_score_table

	// Tunables
	TopK            int
	DriftRatio      float64
	GapThresholdSec float64
	CorrelationMinR float64
	PELTMinSegment  int
	MKTauMin        float64
	MKPMax          float64
	SpikeZ          float64  // robust-z threshold for kind:"spike" emission
	FollowerMetrics []string // user-extended follower vector

	// ReAct CLI primitives
	AroundEvent       time.Time
	AroundWindow      time.Duration
	HasAroundEvent    bool
	SegmentByMetric   string
	CorrelateLeader   string
	CorrelateFollower string
	LagResolutionMS   int

	// Test harness
	CheckDir string // path passed to --auto-check

	// P2.3: --auto-compare path. When non-empty, the parser builds a
	// second seriesStore from this path and emits capture_diff +
	// capture_diff_summary Findings against the primary store.
	ComparePath string

	// P1.NEW-A: --auto-recommend-tunables enables the tunable_observation
	// stage. Default off — bare --auto stays clean. Tier-7 narration
	// scaffold (see BIAS_CONTROL.md).
	RecommendTunables bool

	// P1.6: --multi-host-parallel N. Number of per-host runs to execute
	// concurrently in multi-host mode. Default 1 (sequential) per the
	// solo-user memory constraint (decoded FTDC samples are GB-scale; N>1
	// requires N hosts' columns alive simultaneously and OOMs the host).
	// Documented in flag help with the warning.
	MultiHostParallel int

	// P0.3: input provenance for run_metadata. Populated in main() (or
	// programmatic callers) before runAuto. When InputFingerprint is nil,
	// emitRunMetadata omits the input_sha256 / input_files fields rather
	// than emitting empty strings — keeps the schema honest.
	InputPath        string            // displayed; not necessarily local-absolute
	InputFingerprint *InputFingerprint // streaming SHA256 over input bytes
	CmdLine          []string          // os.Args with paths basenamed

	// --auto-emit-t2: self-contained bundle writer. When AutoEmitT2Path is
	// non-empty, a findingTee wraps the main bufio.Writer so every Finding
	// line is captured and, after the normal output flushes, written into a
	// bundle directory + tarball that t2 can open by double-clicking
	// context.t2.
	AutoEmitT2Path      string // bundle directory path; tarball at <path>.tgz
	AutoEmitT2NoTarball bool   // skip the .tgz (dir only; for local iteration)
	AutoEmitT2TopK      int    // marker/descriptor budget (default 50)
	AutoEmitT2LinkFtdc  bool   // symlink FTDC files instead of copying
	AutoEmitT2Zoom      bool   // emit additional zoom contexts (off by default)
	AutoLogFile         string // optional mongod.log path; trimmed copy + excerpt written to bundle
}

// defaultAutoConfig returns a config with documented defaults. main() applies
// flag overrides on top.
func defaultAutoConfig() AutoConfig {
	return AutoConfig{
		Mode:            "auto",
		TopK:            defaultAutoTopK,
		DriftRatio:      defaultAutoDriftRatio,
		GapThresholdSec: defaultGapThresholdSeconds,
		CorrelationMinR: defaultCorrelationMinR,
		PELTMinSegment:  defaultPELTMinSegment,
		MKTauMin:        defaultMannKendallTauMin,
		MKPMax:          defaultMannKendallPMax,
		SpikeZ:          defaultSpikeZ,
		BudgetTokens:    defaultBudgetTokens,
		EmitScoreTable:  false, // opt-in via --auto-emit-score-table
		AutoEmitT2TopK:  50,    // matches A-Z auto-letter wrap in markers.ts:171-175
	}
}

// defaultBudgetTokens is the soft cap on `--auto` output size, expressed
// as approximate LLM tokens (1 token ≈ 4 output bytes). The full
// metric_score_table alone is ~190K tokens; with that opt-in by default
// the rest of the Findings comfortably fit a typical 100K context
// window. Override via --auto-budget-tokens.
const defaultBudgetTokens = 80000

// runAuto is the --auto entry point. Dispatched from main() after sample
// collection + time-range filter. Always processes a single host's
// pre-collected samples; multi-host dispatch happens in main() before this
// runs.
func runAuto(w io.Writer, path string, samples []Sample, metadata []MetadataDoc, cfg AutoConfig) error {
	return runAutoWithHost(w, hostIDFromPath(path), samples, metadata, cfg)
}

// runAutoWithHost is the legacy entry that takes a []Sample slice. It
// builds a columnar seriesStore from those samples (via the streaming
// builder) and dispatches to runAutoWithHostStore. Kept for backward
// compatibility with the multi-host dispatcher and react modes that
// already hold a filtered samples slice.
func runAutoWithHost(w io.Writer, host string, samples []Sample, metadata []MetadataDoc, cfg AutoConfig) error {
	store := buildSeriesStoreFromSamples(samples)
	return runAutoWithHostStore(w, host, store, metadata, cfg)
}

// runAutoWithHostStore is the column-store-native variant. Used directly
// by the streaming --auto path in ftdc_parser.go's main dispatch so the
// row-oriented []Sample slice never materialises globally.
func runAutoWithHostStore(w io.Writer, host string, store *seriesStore, metadata []MetadataDoc, cfg AutoConfig) error {
	collector := &findingCollector{}
	stats := &AutoStats{}

	// run_metadata is always the first Finding.
	emitRunMetadataStore(collector, host, store, metadata, cfg)

	// Stage 5 (gap + fingerprint) is cheap; emit unconditionally even when
	// later stages skip due to insufficient data.
	stats.timeStage(collector, "gap", 0, func() { emitGapFindingsStore(collector, host, store, cfg) })
	stats.timeStage(collector, "schema_fingerprint", 0, func() { emitSchemaFingerprintStore(collector, host, store, metadata) })

	N := store.nSamples()
	if N < 4 {
		// Too few samples to score; emit a single skipped Finding and bail.
		collector.add(newFinding(KindSkipped, host, "too_few_samples", map[string]interface{}{
			"reason": "fewer than 4 samples; statistical stages skipped",
			"n":      N,
		}))
		return collector.emitAll(w)
	}

	bStart, bEnd, sStart, sEnd, windowFinding := selectWindowsStore(store, cfg)
	if windowFinding != nil {
		collector.add(windowFinding)
	}

	// Stage 1: Welch shift score for every metric, using the store.
	stage1Start := time.Now()
	stage1Before := len(collector.items)
	scoreTable, movers, scores := scoreAndRankMetricsWithStoreCol(host, store, bStart, bEnd, sStart, sEnd, cfg)
	collector.add(scoreTable)
	for _, m := range movers {
		collector.add(m)
	}
	stats.recordStage("stage1_welch", stage1Start, len(collector.items)-stage1Before, 0)

	// Determine top-200 metric IDs (used by Stage 3 + variance_shift).
	// Materialize their flat float64 series once so Stage 3 doesn't have
	// to scan samples per metric.
	allTopIDs := sortedIDsByAbsScoreWithNames(scores, store.names)
	trendIDs := allTopIDs
	if len(trendIDs) > 200 {
		trendIDs = trendIDs[:200]
	}
	changepointIDs := allTopIDs
	if len(changepointIDs) > 30 {
		changepointIDs = changepointIDs[:30]
	}
	varianceIDs := changepointIDs

	// Force core metrics into ALL three pools regardless of Welch score.
	coreIDs := resolveCoreMetricIDs(store)
	trendIDs = mergeUnique(trendIDs, coreIDs)
	changepointIDs = mergeUnique(changepointIDs, coreIDs)
	varianceIDs = mergeUnique(varianceIDs, coreIDs)

	// Materialize series for the UNION of trendIDs + changepointIDs from
	// the column store. trendIDs is a superset since |trendIDs| >=
	// |changepointIDs|.
	store.materializeFromColumns(trendIDs)

	// Stage 2: lagged Pearson correlation against follower vector.
	stats.timeStage(collector, "stage2_correlation", 0, func() {
		for _, f := range correlateAgainstFollowersWithStoreCol(host, store, cfg) {
			collector.add(f)
		}
	})

	// Stage 3: Mann-Kendall on top-200 (cheap) and ED-PELT on top-30.
	stats.timeStage(collector, "stage3_trend", 0, func() {
		for _, f := range trendOnly(host, store, trendIDs, cfg) {
			collector.add(f)
		}
	})
	stats.timeStage(collector, "stage3_changepoint", 0, func() {
		for _, f := range changepointOnly(host, store, changepointIDs, cfg) {
			collector.add(f)
		}
	})

	// Stage 4: histogram percentile inference. Parallel across families.
	stats.timeStage(collector, "stage4_histogram", 0, func() {
		for _, f := range inferHistogramPercentilesStore(host, store, bStart, bEnd, sStart, sEnd, cfg) {
			collector.add(f)
		}
	})

	// Self-observability: variance-shift CUSUM on materialized top movers.
	stats.timeStage(collector, "variance_shift", 0, func() {
		for _, f := range detectVarianceShiftsWithStore(host, store, varianceIDs, cfg) {
			collector.add(f)
		}
	})

	// Spike detector.
	stats.timeStage(collector, "spike", 0, func() {
		for _, f := range detectSpikes(host, store, trendIDs, cfg) {
			collector.add(f)
		}
	})

	// Tier B.8: Hampel local-spike detector. Adds value on autocorrelated
	// series where the global-MAD detector loses sensitivity.
	stats.timeStage(collector, "local_spike", 0, func() {
		emitHampelLocalSpikes(collector, host, store, trendIDs, cfg)
	})

	// Capacity / near-ceiling Findings (structural; gated by
	// ceilings.go catalog presence in the seriesStore).
	stats.timeStage(collector, "capacity", 0, func() {
		emitCapacityPressure(collector, host, store, cfg)
		emitNearCeiling(collector, host, store, cfg)
	})

	// Tier B.7: counter wrap vs reset diagnostic stage.
	stats.timeStage(collector, "counter_events", 0, func() {
		emitCounterEvents(collector, host, store)
	})

	// Tier G.1: capability-extension detectors that read store columns
	// directly (apply_saturation, checkpoint_stall, eviction_divergence).
	// election_storm is moved later — it reads `term_change` Findings,
	// which only emit in the `replication` stage.
	stats.timeStage(collector, "apply_saturation", 0, func() {
		emitApplySaturation(collector, host, store)
	})
	stats.timeStage(collector, "checkpoint_stall", 0, func() {
		emitCheckpointStall(collector, host, store)
	})
	stats.timeStage(collector, "counter_stall", 0, func() {
		emitHeartbeatCounterStall(collector, host, store)
	})
	stats.timeStage(collector, "eviction_divergence", 0, func() {
		emitEvictionDivergence(collector, host, store)
	})

	// Cross-finding consistency check. Runs AFTER all
	// per-metric stages so it can compare mover and monotonic_trend on
	// the same metric.
	stats.timeStage(collector, "consistency", 0, func() {
		emitConsistencyWarning(collector, host)
	})

	// synthetic_event overlay clustering. Children remain in
	// place; this adds an overlay Finding per cluster.
	stats.timeStage(collector, "synthetic_event", 0, func() {
		emitSyntheticEvent(collector, host, store)
	})

	// Periodicity detection. Returns (emitted, skipped); we
	// surface BOTH so the analyzer can distinguish "stage ran, found
	// nothing" from "stage was disabled / metric absent".
	stats.timeStage(collector, "periodic", 0, func() {
		emitted, skipped := emitPeriodic(collector, host, store, trendIDs, cfg)
		_ = emitted
		_ = skipped
		// Skipped count is already accounted for in recordStage when we
		// thread it explicitly. Use the per-emitter return:
		stats.adjustLastSkipped(skipped)
	})

	// suspicious_smoothness.
	stats.timeStage(collector, "smoothness", 0, func() {
		emitted, skipped := emitSuspiciousSmoothness(collector, host, store, trendIDs)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// evidence_convergence.
	stats.timeStage(collector, "convergence", 0, func() {
		emitEvidenceConvergence(collector, host)
	})

	// Utilization (bounded ratios — cache, tickets, connections).
	stats.timeStage(collector, "utilization", 0, func() {
		emitted, skipped := emitUtilization(collector, host, store, bStart, bEnd, sStart, sEnd)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// alert_threshold (election / restart point events).
	stats.timeStage(collector, "alert_threshold", 0, func() {
		emitted, _ := emitAlertThreshold(collector, host, store)
		_ = emitted
	})

	// clock_skew (host-level data-quality guard).
	stats.timeStage(collector, "clock_skew", 0, func() {
		emitted, _ := emitClockSkew(collector, host, store)
		_ = emitted
	})

	// missing_signal (per-version expected-metric manifest).
	stats.timeStage(collector, "missing_signal", 0, func() {
		emitted, _ := emitMissingSignal(collector, host, store)
		_ = emitted
	})

	// Replication domain (gated on replSetGetStatus presence).
	stats.timeStage(collector, "replication", 0, func() {
		emitted, skipped := emitReplicationFindings(collector, host, store)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// Tier G.1 election_storm: must run AFTER replication so the
	// `term_change` Findings it overlays already exist in the collector.
	// Pre-fix this was grouped with the other Tier G.1 detectors above,
	// which silently caused it to never fire (no term_change Findings
	// existed yet at that point).
	stats.timeStage(collector, "election_storm", 0, func() {
		emitElectionStorm(collector, host)
	})

	// PALI / disagg-storage (gated on metrics.disagg.* presence).
	stats.timeStage(collector, "pali", 0, func() {
		emitted, skipped := emitPALIFindings(collector, host, store)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// P4.6 coverage_map: which code areas showed activity vs which were
	// silent. Always-on, no flag. Emits coverage_summary + per-gap
	// coverage_gap Findings.
	stats.timeStage(collector, "coverage_map", 0, func() {
		emitted, skipped := emitCoverageMap(collector, host, store)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// P3.3 mover_stratified: dimension-stratified Welch scoring across
	// families with a varying middle token (e.g. queues.execution.<pool>.*).
	// Recomputes per-window Welford internally — bounded cost. Runs after
	// the Stage 3 windows are known (bStart..sEnd from selectWindowsStore).
	stats.timeStage(collector, "mover_stratified", 0, func() {
		emitted, skipped := emitMoverStratified(collector, host, store, bStart, bEnd, sStart, sEnd)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// P3.1 quantile_trend: Theil-Sen slope on tail-latency-shaped
	// metrics. Independent of changepoint output (it's a robust slope
	// estimator, not a breakpoint detector) so order with the other
	// detector stages doesn't matter — runs here to be near the trend
	// signals analytically related to it.
	stats.timeStage(collector, "quantile_trend", 0, func() {
		emitted, skipped := emitQuantileTrend(collector, host, store)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// P3.2 joint_changepoint: cluster changepoint Findings by sample
	// index; emit when ≥5 align. Anchored on Stage-3 output so MUST run
	// AFTER changepoint stage. Cheap; pure composer.
	stats.timeStage(collector, "joint_changepoint", 0, func() {
		emitted, skipped := emitJointChangepoint(collector, host)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// P1.2 --auto-suggest-windows: cluster changepoint Findings into
	// (baseline, stress) window suggestions. MUST run AFTER changepoint
	// stage (whose Findings it reads). Pure post-stage composer, always-on
	// under --auto; the suggestions are advisory and never gate other
	// Findings.
	stats.timeStage(collector, "window_suggestion", 0, func() {
		emitted, skipped := emitWindowSuggestions(collector, host, store)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// P0.NEW-PALI-1: seal_event composer. Reads term_change, local_spike,
	// gap, counter_inversion Findings already emitted by upstream stages
	// and composes seal_event[/seal_event_candidate] per the D1 spec.
	// MUST run after replication (term_change), local_spike, gap, counter
	// stages, and PALI (which surfaces checkpoint_install_anomaly). Runs
	// before diagnosis_candidate so PALI-specific patterns can compose on
	// seal_event as a building block.
	stats.timeStage(collector, "seal_event", 0, func() {
		emitted, skipped := emitSealEvents(collector, host, store)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// Sharded-cluster (gated on shardingStatistics.* presence).
	stats.timeStage(collector, "sharded", 0, func() {
		emitted, skipped := emitShardedFindings(collector, host, store)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// host_os (gated on systemMetrics.* presence).
	stats.timeStage(collector, "host_os", 0, func() {
		emitted, skipped := emitHostOSFindings(collector, host, store)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// workload_spec (numeric workload profile from opcounters).
	stats.timeStage(collector, "workload_spec", 0, func() {
		emitted, skipped := emitWorkloadSpec(collector, host, store)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// precursor_candidates (cross-stage precedence reasoning).
	// Runs AFTER all detector stages so it sees the full Finding set.
	stats.timeStage(collector, "precursor_candidates", 0, func() {
		emitted, skipped := emitPrecursorCandidates(collector, host)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// memory_triplet (joint shape of tcmalloc + mem.resident).
	stats.timeStage(collector, "memory_triplet", 0, func() {
		emitted, skipped := emitMemoryTriplet(collector, host, store)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// recurrence_fingerprint (compact stable hash of the run
	// shape; persistable across captures for cross-run comparison).
	// Runs LAST among the detector tiers so it sees the full Finding
	// set including the overlay kinds.
	stats.timeStage(collector, "recurrence_fingerprint", 0, func() {
		emitRecurrenceFingerprint(collector, host)
	})

	// Tier D: cross-run regression + recurrence. No-op when
	// cfg.BaselineFile is unset. MUST run AFTER recurrence_fingerprint
	// (which it reads to update hit_count) AND AFTER all detector stages
	// (so the EWMA mean reflects everything emitted this run). Pre-fix
	// this stage ran BEFORE recurrence_fingerprint and silently never
	// incremented hit_count.
	stats.timeStage(collector, "baseline", 0, func() {
		_, skipped := runBaselineStage(collector, host, store, cfg)
		stats.adjustLastSkipped(skipped)
	})

	// Tier G.3: composite capture-quality score. MUST run last among
	// the count-reading stages so all detector emits + missing_signal
	// + chunk_decode_error + gap counts are available. Pre-fix this
	// ran BEFORE missing_signal/clock_skew/gap; counts were undercounted.
	stats.timeStage(collector, "capture_quality_score", 0, func() {
		emitCaptureQualityScore(collector, host)
	})

	// Tier R.1: diagnosis_candidate emission. Reads emitted Findings;
	// emits Finding-shape pattern matches. MUST run AFTER all detector
	// stages AND missing_signal AND capture_quality_score (all referenced
	// by pattern requiredKinds/supportingKinds). Pre-fix this ran BEFORE
	// missing_signal and capture_quality_score; patterns referencing
	// those kinds never matched.
	stats.timeStage(collector, "diagnosis_candidate", 0, func() {
		emitDiagnosisCandidates(collector, host, store)
	})

	// Tier R.2: hypothesis_test emission. Reads diagnosis_candidate
	// Findings emitted by R.1; emits copy-pasteable recipes for the
	// discriminating kinds each candidate still needs.
	stats.timeStage(collector, "hypothesis_test", 0, func() {
		emitHypothesisTests(collector, host)
	})

	// P2.2 diagnosis_uncertainty: per-candidate action list surfacing
	// what observation would push confidence higher. Runs AFTER
	// hypothesis_test so the user can pair the two; same-pass composer
	// (reads diagnosis_candidate, no new computation).
	stats.timeStage(collector, "diagnosis_uncertainty", 0, func() {
		emitted, skipped := emitDiagnosisUncertainty(collector, host)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// P2.4 followups_backlog_snapshot: append mid-confidence candidates
	// to a solo-user backlog at ~/claude-dump/ftdc-followups.ndjson and
	// emit a snapshot Finding so the user sees them in-stream.
	stats.timeStage(collector, "followups_backlog", 0, func() {
		emitted, skipped := emitFollowupsBacklog(collector, host)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// P2.1 --auto-recap: cross-session continuity keyed by input_sha256.
	// Reads ~/claude-dump/ftdc-recaps/<sha>.ndjson; emits kind:recap for
	// still-valid prior entries and kind:stored_inference_review for
	// quarantined ones. Appends the current run's compact summary at
	// the end. Runs late so it has access to diagnosis_candidate +
	// capture_quality_score + metric_score_table summaries.
	stats.timeStage(collector, "recap", 0, func() {
		emitted, skipped := emitRecap(collector, host)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// P1.NEW-A tunable_observation: opt-in (cfg.RecommendTunables) Tier-7
	// narration scaffold. Reads supporting Findings already in the
	// collector + the static catalog in tunable_guardrails.go. No-op
	// when the flag is unset.
	stats.timeStage(collector, "tunable_observation", 0, func() {
		emitted, skipped := emitTunableObservations(collector, host, store, cfg)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// P4.5 tunable_counterfactual: model-based prediction for each
	// emitted tunable_observation. Auto-triggered when cfg.RecommendTunables
	// is set; refuses when the tunable lacks an applicable model.
	stats.timeStage(collector, "tunable_counterfactual", 0, func() {
		emitted, skipped := emitTunableCounterfactuals(collector, host, store, cfg)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// P4.3 chaos-sweep auto-labeled comparison. Reads the fingerprint
	// library for sibling captures (no second seriesStore needed),
	// auto-labels HEALTHY/DEGRADED via 1-D k-means on per-capture
	// confidence × shape-set-size composite, and emits
	// sweep_discriminator + sweep_label_audit Findings.
	stats.timeStage(collector, "chaos_sweep", 0, func() {
		emitted, skipped := emitChaosSweep(collector, host)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// P4.1 autonomous diagnosis loop (planning only in v1): walks
	// sibling capture directories, extracts recipe + expected outcomes
	// from hypothesis_test, emits a structured converge_plan Finding.
	// Memory constraint blocks in-process execution; a future driver
	// can re-invoke ftdc_parser as a subprocess.
	stats.timeStage(collector, "converge_plan", 0, func() {
		emitted, skipped := emitConvergePlan(collector, host, store, cfg)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// P4.4 fingerprint library: persist a canonical shape-set fingerprint
	// for the current run + scan prior captures for high-Jaccard matches.
	// Solo + FTDC-only — corpus is the user's own past captures. All
	// stored entries carry provenance; stale + version-skewed entries
	// are quarantined and surfaced via stored_inference_review Findings.
	stats.timeStage(collector, "fingerprint_library", 0, func() {
		emitted, skipped := emitFingerprintMatches(collector, host)
		_ = emitted
		stats.adjustLastSkipped(skipped)
	})

	// Per-capture context Finding (per-capture pre-context the analyzer skill
	// cannot derive from FTDC values alone).
	stats.timeStage(collector, "metric_context", 0, func() {
		emitMetricContext(collector, host, store, metadata)
	})

	// Cited-metric glossary Finding (one-shot, keyed only on metrics cited
	// elsewhere in the output).
	stats.timeStage(collector, "metric_glossary", 0, func() {
		emitMetricGlossary(collector, host)
	})

	// Output gates: score-table opt-in + token-budget truncation.
	// Applied BEFORE run_health so the run_health Finding's
	// `findings_emitted_by_kind` reports POST-filter counts (matching
	// what's actually emitted). R5 final-acceptance review issue.
	applyOutputGates(collector, cfg)

	// End-of-run health Finding (must be emitted last; aggregates per-Finding
	// cross-validation outcomes over the post-filter Finding set). Final-acceptance
	// extension: attach per-stage observability stats so the analyzer can
	// distinguish "stage didn't run" from "stage ran and found nothing".
	rh := buildRunHealth(host, collector.items, N, len(metadata))
	if stages := stats.asJSON(); len(stages) > 0 {
		rh["stages"] = stages
	}
	collector.add(rh)

	if cfg.EmitOTel {
		return collector.otelEmitAll(w)
	}
	return collector.emitAll(w)
}

// resolveCoreMetricIDs translates the coreMetricsAlwaysAnalyzed name list
// into integer IDs against the current capture's seriesStore. Names not
// present in the capture (e.g. renamed in a future MongoDB version) are
// silently dropped — the analyzer skill catches the drift via the
// schema_fingerprint Finding.
func resolveCoreMetricIDs(store *seriesStore) []int {
	out := make([]int, 0, len(coreMetricsAlwaysAnalyzed))
	for _, name := range coreMetricsAlwaysAnalyzed {
		if id, ok := store.nameToID[name]; ok {
			out = append(out, id)
		}
	}
	return out
}

// mergeUnique returns the deduplicated union of two int slices preserving
// the relative order of the first slice followed by any new entries from
// the second. Used to add core-metric IDs to Stage-3 pools without
// duplicating IDs that already made the Welch top-K cut.
func mergeUnique(a, b []int) []int {
	seen := make(map[int]bool, len(a)+len(b))
	out := make([]int, 0, len(a)+len(b))
	for _, id := range a {
		if !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	for _, id := range b {
		if !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	return out
}

// parseRange parses a "FROM..TO" RFC3339 range argument used by --baseline
// and --stress. Returns ok=false on malformed input.
func parseRange(s string) (time.Time, time.Time, bool) {
	idx := strings.Index(s, "..")
	if idx < 0 {
		return time.Time{}, time.Time{}, false
	}
	from, err := time.Parse(time.RFC3339, s[:idx])
	if err != nil {
		return time.Time{}, time.Time{}, false
	}
	to, err := time.Parse(time.RFC3339, s[idx+2:])
	if err != nil {
		return time.Time{}, time.Time{}, false
	}
	return from, to, true
}

// autoJSONSchema returns the JSON Schema (Draft 2020-12) describing the
// Findings emitted by --auto mode. Implementation is in analyze_schema.go.
// Wired via stub here so the build succeeds before the schema doc is
// authored.
func autoJSONSchema() string {
	return jsonSchemaText()
}

// runAutoCheck is the test harness. Implementation is in analyze_check.go.
// Stub here for build wiring.
func runAutoCheck(w io.Writer, dir string, jobs int) error {
	return autoCheckRun(w, dir, jobs)
}

// hostIDFromPath derives a short host identifier from a path. For a directory,
// it's the basename. Used as the host field on Findings. Falls back to
// "default" for unusual inputs.
func hostIDFromPath(path string) string {
	// Strip trailing slashes.
	p := strings.TrimRight(path, "/")
	if p == "" {
		return "default"
	}
	idx := strings.LastIndex(p, "/")
	if idx < 0 {
		return p
	}
	return p[idx+1:]
}

// emitRunMetadata writes the run_metadata Finding. Echoes every algorithmic
// choice (window bounds, thresholds, library versions) so the LLM can audit.
// This is the bias-control "show the choice" invariant.
func emitRunMetadata(c *findingCollector, host string, samples []Sample, metadata []MetadataDoc, cfg AutoConfig) {
	var startTS, endTS string
	if len(samples) > 0 {
		startTS = samples[0].Timestamp.UTC().Format(time.RFC3339Nano)
		endTS = samples[len(samples)-1].Timestamp.UTC().Format(time.RFC3339Nano)
	}
	fields := map[string]interface{}{
		"sample_count":      len(samples),
		"metadata_count":    len(metadata),
		"window_start":      startTS,
		"window_end":        endTS,
		"baseline_fraction": defaultBaselineFraction,
		"stress_fraction":   defaultStressFraction,
		"lag_grid_seconds":  []int{0, 1, 5, 10, 30, 60, 120, 300},
		"library_versions": map[string]string{
			"changepoint": "pgregory.net/changepoint@v1.0.0",
			// Tier A.4 fix: runtime.Version() returns the build's Go version
			// (e.g. "go1.17.2") which differs across developer machines,
			// breaking byte-equality. Replace with a build-baked parser
			// version constant.
			"parser_version": ParserVersion,
		},
		"histogram_recognizers": histogramRecognizerNames(),
		"mode":                  cfg.Mode,
		// Bias-control invariant ("show the choice") — every algorithmic
		// threshold that affects which Findings get emitted is recorded
		// here so the LLM downstream can audit. Hidden choices = bias.
		"thresholds": map[string]interface{}{
			"top_k":                           cfg.TopK,
			"drift_ratio":                     cfg.DriftRatio,
			"gap_threshold_sec":               cfg.GapThresholdSec,
			"correlation_min_r":               cfg.CorrelationMinR,
			"correlation_per_follower_cap":    30,
			"pelt_min_segment":                cfg.PELTMinSegment,
			"pelt_downsample_cap":             4000,
			"mk_tau_min":                      cfg.MKTauMin,
			"mk_p_max":                        cfg.MKPMax,
			"slope_t_threshold":               2.0,
			"spike_z_threshold":               cfg.SpikeZ,
			"variance_shift_sigma_multiplier": 3.0,
			"saturation_weight":               4.0,
			// Tier-P0.5 audit-pass: previously-hidden choices lifted here for
			// bias-control "show the choice" invariant. Values mirror call-site
			// literals; a follow-up should unify via package-level consts to
			// avoid drift. Source files cited in the key suffix.
			"capacity_pressure_default_run_seconds":   10,     // analyze_capacity.go:187 fallback
			"capacity_pressure_default_pct_at":        0.05,   // analyze_capacity.go:191 fallback
			"capacity_pressure_event_radius_min_sec":  5,      // analyze_capacity.go:331 floor
			"capacity_pressure_max_synthetic_events":  50,     // analyze_capacity.go:388 cap
			"changepoint_max_breaks_per_metric":       25,     // analyze_changepoint.go:190 cap
			"changepoint_refinement_min_tolerance":    5,      // analyze_changepoint.go:610 floor
			"apply_saturation_latency_threshold_ms":   100.0,  // analyze_capabilities.go:56 const
			"apply_saturation_window_samples":         60,     // analyze_capabilities.go:57 const
			"apply_saturation_majority_fraction":      0.5,    // analyze_capabilities.go:58 const
			"eviction_divergence_ratio_threshold":     2.0,    // analyze_capabilities.go:239 const
			"baseline_ttl_days":                       30,     // analyze_baseline.go:50 const
			"baseline_regression_sigma":               3.0,    // analyze_baseline.go:53 const
			"periodicity_alpha":                       0.001,  // analyze_signal.go:97 Bonferroni-corrected
			"glossary_top_per_kind":                   3,      // analyze_glossary.go:349 const
			"selfcheck_min_window_size":               30,     // analyze_selfcheck.go:37 floor
			"fdr_alpha":                               0.05,   // P0.NEW-7 Benjamini-Hochberg target FDR
		},
		"topk_tiers": map[string]int{
			"score_mover_each_direction": cfg.TopK,
			"trend":                      200,
			"changepoint":                30,
			"variance_shift":             30,
			"spike":                      200,
		},
		// Core metrics always get deep analysis (Mann-Kendall, ED-PELT,
		// variance-shift CUSUM, spike) regardless of Stage-1 Welch shift
		// score. Echoed so the LLM downstream knows this layer exists.
		"core_metrics_always_analyzed": coreMetricsAlwaysAnalyzed,
	}
	// Echo the computed baseline/stress sample-index ranges (when not
	// explicitly user-set) so the LLM doesn't need to derive them from
	// fractions + total window.
	if n := len(samples); n >= 4 && !cfg.HasBaseline && !cfg.HasStress {
		bCount := int(float64(n) * defaultBaselineFraction)
		if bCount < 1 {
			bCount = 1
		}
		sStart := n - int(float64(n)*defaultStressFraction)
		if sStart >= n {
			sStart = n - 1
		}
		if sStart < bCount {
			sStart = bCount
		}
		fields["effective_baseline"] = map[string]interface{}{
			"sample_idx_lo": 0,
			"sample_idx_hi": bCount,
			"from":          samples[0].Timestamp.UTC().Format(time.RFC3339Nano),
			"to":            samples[bCount-1].Timestamp.UTC().Format(time.RFC3339Nano),
		}
		fields["effective_stress"] = map[string]interface{}{
			"sample_idx_lo": sStart,
			"sample_idx_hi": n,
			"from":          samples[sStart].Timestamp.UTC().Format(time.RFC3339Nano),
			"to":            samples[n-1].Timestamp.UTC().Format(time.RFC3339Nano),
		}
	}
	if cfg.HasBaseline {
		fields["explicit_baseline"] = map[string]string{
			"from": cfg.BaselineFrom.UTC().Format(time.RFC3339Nano),
			"to":   cfg.BaselineTo.UTC().Format(time.RFC3339Nano),
		}
	}
	if cfg.HasStress {
		fields["explicit_stress"] = map[string]string{
			"from": cfg.StressFrom.UTC().Format(time.RFC3339Nano),
			"to":   cfg.StressTo.UTC().Format(time.RFC3339Nano),
		}
	}
	if len(cfg.FollowerMetrics) > 0 {
		fields["user_follower_metrics"] = cfg.FollowerMetrics
	}
	// P0.3 input provenance. Both fields are advisory — a downstream consumer
	// may skip them — but they let a reviewer verify "same bytes, same
	// invocation" without re-hashing the inputs themselves.
	if cfg.InputFingerprint != nil {
		fields["input_sha256"] = cfg.InputFingerprint.TopSHA256
		if cfg.InputFingerprint.IsDirectory && len(cfg.InputFingerprint.PerFileSHA256) > 0 {
			fields["input_files"] = cfg.InputFingerprint.PerFileSHA256
		}
	}
	if cfg.InputPath != "" {
		fields["input_path_basename"] = filepathBase(cfg.InputPath)
	}
	if len(cfg.CmdLine) > 0 {
		fields["cmd_line"] = cfg.CmdLine
	}
	c.add(newFinding(KindRunMetadata, host, "global", fields))
}

// filepathBase wraps filepath.Base but tolerates empty strings cleanly so
// the emitter can call it unconditionally. Avoids importing filepath at
// the top of analyze.go for a single use.
func filepathBase(p string) string {
	if p == "" {
		return ""
	}
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' || p[i] == '\\' {
			return p[i+1:]
		}
	}
	return p
}

// ─────────────────────────────────────────────────────────────────────────
// Stubs for Stage 1 (Welch scoring), Stage 5 (gap + fingerprint), and
// window selection. Each is replaced by full implementation in its own
// task; the stubs let the build succeed while feature work proceeds
// incrementally.
// ─────────────────────────────────────────────────────────────────────────

// selectWindows picks baseline and stress sample slices. Stub returns
// first 25% / last 25% by sample count. Implementation lands in task #12.
func selectWindows(samples []Sample, cfg AutoConfig) ([]Sample, []Sample, Finding) {
	n := len(samples)
	if n < 4 {
		return nil, nil, nil
	}
	if cfg.HasBaseline || cfg.HasStress {
		// Honor explicit windows when provided. Defer full
		// implementation to task #12.
		_ = cfg
	}
	bCount := int(float64(n) * defaultBaselineFraction)
	if bCount < 1 {
		bCount = 1
	}
	sStart := n - int(float64(n)*defaultStressFraction)
	if sStart >= n {
		sStart = n - 1
	}
	if sStart < bCount {
		sStart = bCount
	}
	return samples[:bCount], samples[sStart:], nil
}

// metricAcc holds Welford streaming-variance accumulators for a single metric
// across baseline and stress windows, plus the auxiliary state needed to
// detect counter vs gauge and compute saturation deviation for gauges.
//
// Welford is numerically stable for streaming variance:
//
//	count += 1
//	delta  = x - mean
//	mean  += delta / count
//	delta2 = x - mean
//	m2    += delta * delta2
//	variance = m2 / (count - 1)
type metricAcc struct {
	// Kind detection (computed from full sample series, not per-window).
	isCounter    bool
	isDegenerate bool // min == max across the full series

	// Full-series bookkeeping for kind / saturation.
	observedMax int64
	observedMin int64
	sawDecrease bool
	seenAny     bool

	// Baseline window Welford (gauge values, or counter first-differences).
	baseN    int
	baseMean float64
	baseM2   float64

	// Stress window Welford.
	stressN    int
	stressMean float64
	stressM2   float64

	// Saturation component (gauges only): mean of value/observedMax per
	// window. Token-bucket gauges hover near max; a drop signals back-pressure
	// that pure CoV misses.
	baseSatSum   float64
	stressSatSum float64

	// Boundary state for counter rate within each window: the value of the
	// sample that immediately precedes the window (for the first window
	// sample's rate). nil-equivalent means "no prior available; skip i=0".
	basePrev      int64
	hasBasePrev   bool
	stressPrev    int64
	hasStressPrev bool
}

// scoreAndRankMetricsWithStore is the optimized Stage 1 that operates on
// a pre-built seriesStore. Avoids redundant classification passes and uses
// array indexing instead of map lookups for Welford accumulation. Target:
// 3-5s on a 5.5h × 6500-metric fixture vs the original 40+s.
func scoreAndRankMetricsWithStore(host string, store *seriesStore, baselineSamples, stressSamples []Sample, cfg AutoConfig) (Finding, []Finding, []float64) {
	M := len(store.names)
	if M == 0 || len(baselineSamples) == 0 || len(stressSamples) == 0 {
		return newFinding(KindMetricScoreTable, host, "global", map[string]interface{}{
			"entries":      []interface{}{},
			"empty_reason": "no metrics or empty window",
		}), nil, nil
	}

	// Per-metric per-window Welford. Both passes (baseline, stress) run
	// concurrently — they don't share state.
	baseN, baseMean, baseM2, baseSatSum := parallelWelfordWindow(store, baselineSamples, M)
	stressN, stressMean, stressM2, stressSatSum := parallelWelfordWindow(store, stressSamples, M)

	// Compose scores. For each metric, Welch t-statistic + saturation
	// boost (gauges only). The combined score ranks movers.
	scores := make([]float64, M)
	for id := 0; id < M; id++ {
		if store.isDegenerate[id] {
			continue
		}
		bSD := welfordSD(baseN[id], baseM2[id])
		sSD := welfordSD(stressN[id], stressM2[id])
		t := welchT(baseMean[id], bSD, baseN[id], stressMean[id], sSD, stressN[id])
		var satDelta float64
		if !store.isCounter[id] && baseN[id] > 0 && stressN[id] > 0 {
			satBase := baseSatSum[id] / float64(baseN[id])
			satStress := stressSatSum[id] / float64(stressN[id])
			satDelta = satBase - satStress
		}
		combined := t
		if !store.isCounter[id] {
			combined = t + 4.0*satDelta
		}
		scores[id] = combined
	}

	// Score table: deterministic alphabetical-by-name order.
	// Bias-control invariant ("no filtering"): include EVERY metric. Degenerate
	// constants (min == max across the capture) appear with kind:"constant"
	// and score 0 so the LLM downstream sees the FULL metric inventory and can
	// decide which "always-zero" or "always-fixed" metrics matter to its
	// diagnosis. Filtering would hide >5000 metrics on a typical fixture.
	type scored struct {
		id    int
		name  string
		kind  string
		score float64
	}
	scoredAll := make([]scored, 0, M)
	nDegenerate := 0
	for id := 0; id < M; id++ {
		kind := "gauge"
		if store.isDegenerate[id] {
			kind = "constant"
			nDegenerate++
		} else if store.isCounter[id] {
			kind = "counter"
		}
		scoredAll = append(scoredAll, scored{id, store.names[id], kind, scores[id]})
	}
	sort.Slice(scoredAll, func(i, j int) bool { return scoredAll[i].name < scoredAll[j].name })
	entries := make([]interface{}, 0, len(scoredAll))
	for _, s := range scoredAll {
		entries = append(entries, map[string]interface{}{
			"metric": s.name,
			"kind":   s.kind,
			"score":  roundFloat(s.score, 4),
		})
	}
	scoreTable := newFinding(KindMetricScoreTable, host, "global", map[string]interface{}{
		"entries":      entries,
		"n_metrics":    len(scoredAll),
		"n_degenerate": nDegenerate,
	})

	// Top-K each direction. Skip constants (kind:"constant", score=0) —
	// they're in the score table but a "mover" Finding for a metric that
	// didn't move would be misleading.
	movableScored := make([]scored, 0, len(scoredAll))
	for _, s := range scoredAll {
		if s.kind == "constant" {
			continue
		}
		movableScored = append(movableScored, s)
	}
	scoredAll = movableScored
	// Tier A.4 fix: stable sort with explicit (score desc, name asc) tiebreaker
	// so --jobs=1 vs --jobs=8 produce byte-identical output on ties.
	sort.SliceStable(scoredAll, func(i, j int) bool {
		if scoredAll[i].score != scoredAll[j].score {
			return scoredAll[i].score > scoredAll[j].score
		}
		return scoredAll[i].name < scoredAll[j].name
	})
	movers := make([]Finding, 0, 2*cfg.TopK)
	emitMover := func(s scored, direction string) {
		id := s.id
		bSD := welfordSD(baseN[id], baseM2[id])
		sSD := welfordSD(stressN[id], stressM2[id])
		t := welchT(baseMean[id], bSD, baseN[id], stressMean[id], sSD, stressN[id])
		// Tier B.9: emit Welch-Satterthwaite df so the analyzer can compute
		// a calibrated p-value rather than relying on |t|>2 ≈ p<0.05.
		// Tier B.10: emit p_value uniformly so downstream can compose findings.
		df := welchDF(bSD, baseN[id], sSD, stressN[id])
		var pValue float64
		if df > 0 {
			pValue = clampFloat01(2 * studentTSurvivalApprox(t, df)) // two-sided
		}
		// Tier B.5: emit Cohen's d on movers (already on changepoints).
		// Uses pooled SD via pooledStandardDeviation (Cohen's classical
		// convention). For unequal-variance regimes, Welch t + df gives
		// the better significance signal; d is for effect-size context.
		var cohenD float64
		if baseN[id] >= 2 && stressN[id] >= 2 {
			pooledSD := pooledStandardDeviation(bSD, baseN[id], sSD, stressN[id])
			if pooledSD > 0 {
				cohenD = (stressMean[id] - baseMean[id]) / pooledSD
			}
		}
		var satDelta float64
		if !store.isCounter[id] && baseN[id] > 0 && stressN[id] > 0 {
			satDelta = (baseSatSum[id] / float64(baseN[id])) - (stressSatSum[id] / float64(stressN[id]))
		}
		fields := map[string]interface{}{
			"metric":      s.name,
			"metric_kind": s.kind,
			"direction":   direction,
			"score":       roundFloat(s.score, 4),
			"welch_t":     roundFloat(t, 4),
			"welch_df":    roundFloat(df, 2),
			"p_value":     roundFloat(pValue, 6),
			"cohens_d":    roundFloat(cohenD, 4),
			"baseline": map[string]interface{}{
				"mean":  roundFloat(baseMean[id], 6),
				"stdev": roundFloat(bSD, 6),
				"n":     baseN[id],
			},
			"stress": map[string]interface{}{
				"mean":  roundFloat(stressMean[id], 6),
				"stdev": roundFloat(sSD, 6),
				"n":     stressN[id],
			},
		}
		if s.kind == "gauge" {
			fields["saturation_delta"] = roundFloat(satDelta, 6)
		}
		movers = append(movers, newFinding(KindMover, host, s.name+"|"+direction, fields))
	}
	k := cfg.TopK
	if k <= 0 {
		k = defaultAutoTopK
	}
	for i := 0; i < k && i < len(scoredAll); i++ {
		emitMover(scoredAll[i], "up")
	}
	for i := 0; i < k && i < len(scoredAll); i++ {
		emitMover(scoredAll[len(scoredAll)-1-i], "down")
	}
	return scoreTable, movers, scores
}

// welfordUpdate is the streaming-variance update step (Welford 1962).
// Stable for very long series and works incrementally.
func welfordUpdate(n *int, mean, m2 *float64, x float64) {
	*n++
	delta := x - *mean
	*mean += delta / float64(*n)
	delta2 := x - *mean
	*m2 += delta * delta2
}

// welfordSD returns the sample standard deviation given the streaming
// accumulators. Returns 0 when n < 2 (variance is undefined for one or zero
// observations); callers must guard against degenerate inputs.
func welfordSD(n int, m2 float64) float64 {
	if n < 2 {
		return 0
	}
	v := m2 / float64(n-1)
	if v < 0 {
		// Numerical underflow on near-constant input. Treat as zero.
		return 0
	}
	return sqrtFloat(v)
}

// parallelWelfordWindow runs Welford streaming variance per metric over a
// window of samples. Counters: rates (first differences). Gauges: raw + a
// saturation accumulator. Sample-chunk parallel; per-worker accumulators
// merged after via Chan's parallel-Welford combine formula:
//
//	delta = mean_b - mean_a
//	combined_n = n_a + n_b
//	combined_mean = mean_a + delta * n_b / combined_n
//	combined_m2 = m2_a + m2_b + delta² * n_a * n_b / combined_n
//
// For counters' first-difference series, each worker computes its own
// per-metric rate based on the prev value at the worker's start; we
// snapshot the value at chunk-boundary samples in a serial pre-pass so
// cross-chunk rate continuity holds (no dropped or duplicated rate
// observations at chunk seams).
func parallelWelfordWindow(store *seriesStore, samples []Sample, M int) (
	[]int, []float64, []float64, []float64,
) {
	if len(samples) == 0 {
		return make([]int, M), make([]float64, M), make([]float64, M), make([]float64, M)
	}
	// Tier A.4 fix: pin worker count to a deterministic constant
	// regardless of NumCPU. Pre-fix, defaultJobs() = min(NumCPU,8), so a
	// 4-core machine produced 4 workers and a 16-core machine 8, giving
	// different parallel-Welford reduction trees → ULP-divergent m2 →
	// divergent Welch t / Cohen's d / changepoint confidence across
	// machines. Chan's combine is exact regardless of nA/nB asymmetry,
	// so pinning the tree shape is sufficient for byte-equality.
	workers := welfordWorkers
	if len(samples) < workers*100 {
		workers = 1
	}
	chunkSize := (len(samples) + workers - 1) / workers

	// Boundary-snapshot pre-pass: serial walk to record, for each
	// (chunk_w, metric_id), the value of that metric at the LAST sample
	// before chunk_w begins. Workers seed their `prev[]` from this so
	// cross-chunk counter rates connect.
	boundaryPrev := make([][]int64, workers)
	boundaryHas := make([][]bool, workers)
	for w := 0; w < workers; w++ {
		boundaryPrev[w] = make([]int64, M)
		boundaryHas[w] = make([]bool, M)
	}
	runningPrev := make([]int64, M)
	runningHas := make([]bool, M)
	seenThisSample := make([]bool, M)
	for w := 0; w < workers; w++ {
		// Snapshot current running state into boundary for chunk w.
		// (For chunk 0 this snapshot is all-zeros/all-false, meaning the
		// first sample of the window correctly has no prior.)
		copy(boundaryPrev[w], runningPrev)
		copy(boundaryHas[w], runningHas)
		lo := w * chunkSize
		hi := lo + chunkSize
		if hi > len(samples) {
			hi = len(samples)
		}
		for i := lo; i < hi; i++ {
			for j := range seenThisSample {
				seenThisSample[j] = false
			}
			for name, v := range samples[i].Metrics {
				id, ok := store.nameToID[name]
				if !ok {
					continue
				}
				seenThisSample[id] = true
				runningPrev[id] = v
				runningHas[id] = true
			}
			// Reset runningHas for metrics absent from this sample so a
			// reappearing-after-gap counter starts a fresh rate sequence
			// rather than computing a single cross-gap rate observation.
			for id := 0; id < M; id++ {
				if !seenThisSample[id] && runningHas[id] {
					runningHas[id] = false
				}
			}
		}
	}

	type part struct {
		n      []int
		mean   []float64
		m2     []float64
		satSum []float64
	}
	parts := make([]part, workers)
	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		lo := w * chunkSize
		hi := lo + chunkSize
		if hi > len(samples) {
			hi = len(samples)
		}
		if lo >= len(samples) {
			break
		}
		parts[w] = part{
			n:      make([]int, M),
			mean:   make([]float64, M),
			m2:     make([]float64, M),
			satSum: make([]float64, M),
		}
		wg.Add(1)
		go func(w, lo, hi int) {
			defer wg.Done()
			pN := parts[w].n
			pMean := parts[w].mean
			pM2 := parts[w].m2
			pSat := parts[w].satSum
			// Seed prev[] from the chunk-boundary snapshot so the first
			// counter rate observation in this chunk connects to the
			// previous chunk's last value. Without this, (workers-1) rate
			// observations per counter metric are silently dropped.
			prev := make([]int64, M)
			hasPrev := make([]bool, M)
			copy(prev, boundaryPrev[w])
			copy(hasPrev, boundaryHas[w])
			seen := make([]bool, M)
			for i := lo; i < hi; i++ {
				for j := range seen {
					seen[j] = false
				}
				for name, v := range samples[i].Metrics {
					id, ok := store.nameToID[name]
					if !ok || store.isDegenerate[id] {
						continue
					}
					seen[id] = true
					if store.isCounter[id] {
						if hasPrev[id] {
							welfordUpdate(&pN[id], &pMean[id], &pM2[id], float64(v-prev[id]))
						}
						prev[id] = v
						hasPrev[id] = true
					} else {
						welfordUpdate(&pN[id], &pMean[id], &pM2[id], float64(v))
						if store.observedMax[id] > 0 {
							pSat[id] += float64(v) / float64(store.observedMax[id])
						}
					}
				}
				// Reset hasPrev for counters absent from this sample so
				// the next observation starts a fresh rate sequence
				// rather than producing a cross-gap rate. Mirrors the
				// fix in buildSeriesStore Phase 2 and materialize.
				for id := 0; id < M; id++ {
					if !seen[id] && hasPrev[id] && store.isCounter[id] {
						hasPrev[id] = false
					}
				}
			}
		}(w, lo, hi)
	}
	wg.Wait()

	// Combine partial Welford accumulators.
	n := make([]int, M)
	mean := make([]float64, M)
	m2 := make([]float64, M)
	satSum := make([]float64, M)
	for _, p := range parts {
		if p.n == nil {
			continue
		}
		for id := 0; id < M; id++ {
			if p.n[id] == 0 {
				continue
			}
			if n[id] == 0 {
				n[id] = p.n[id]
				mean[id] = p.mean[id]
				m2[id] = p.m2[id]
			} else {
				nA, meanA, m2A := n[id], mean[id], m2[id]
				nB, meanB, m2B := p.n[id], p.mean[id], p.m2[id]
				combinedN := nA + nB
				delta := meanB - meanA
				combinedMean := meanA + delta*float64(nB)/float64(combinedN)
				combinedM2 := m2A + m2B + delta*delta*float64(nA)*float64(nB)/float64(combinedN)
				n[id] = combinedN
				mean[id] = combinedMean
				m2[id] = combinedM2
			}
			satSum[id] += p.satSum[id]
		}
	}
	return n, mean, m2, satSum
}

// welchT computes Welch's t-statistic for two-sample comparison with
// potentially unequal variances. Returns 0 when either sample is empty or
// both standard errors are zero (no signal).
func welchT(mb, sb float64, nb int, ms, ss float64, ns int) float64 {
	if nb < 2 || ns < 2 {
		return 0
	}
	se := (sb*sb)/float64(nb) + (ss*ss)/float64(ns)
	if se <= 0 {
		return 0
	}
	return (ms - mb) / sqrtFloat(se)
}

// welchDF returns the Welch-Satterthwaite approximation to the degrees of
// freedom for the two-sample Welch t-test. Without this, |t| > 2 is only a
// rough z-approximation; with it, downstream can compute a calibrated
// p-value. Tier B.9 fix.
//
//   df = (s²₁/n₁ + s²₂/n₂)² / ((s²₁/n₁)²/(n₁−1) + (s²₂/n₂)²/(n₂−1))
//
// Returns 0 when either sample is too small or both variances are zero.
func welchDF(sb float64, nb int, ss float64, ns int) float64 {
	if nb < 2 || ns < 2 {
		return 0
	}
	vb := (sb * sb) / float64(nb)
	vs := (ss * ss) / float64(ns)
	num := vb + vs
	if num <= 0 {
		return 0
	}
	denom := (vb*vb)/float64(nb-1) + (vs*vs)/float64(ns-1)
	if denom <= 0 {
		return 0
	}
	return (num * num) / denom
}

// studentTSurvivalApprox approximates the one-sided p-value `Pr(T_df > |t|)`
// for Student-t with `df` degrees of freedom. Uses an Abramowitz-Stegun
// 26.7.4-derived rational approximation (NOT Hill 1970, despite earlier
// commentary). Accurate to ~3 significant figures for df ≥ 10. Returns 0.5
// when df ≤ 0. Used by Tier B.10 (uniform p_value emission).
func studentTSurvivalApprox(tStat, df float64) float64 {
	if df <= 0 || math.IsNaN(df) {
		return 0.5
	}
	at := tStat
	if at < 0 {
		at = -at
	}
	// Asymptotic normal: P(|Z| > t) for df → ∞
	// For finite df, use the well-known Abramowitz-Stegun 26.7.4
	// approximation with correction.
	// Z ≈ t / sqrt(1 + t²/(2·df)) is rough; refine with Hill term.
	w := at / sqrtFloat(df)
	z := at * (1.0 - 1.0/(4.0*df)) / sqrtFloat(1.0+w*w/2.0)
	return normalSurvival(z)
}

// normalSurvival returns 1 - Φ(z) using a Zelen-Severo rational approximation
// (Abramowitz-Stegun 26.2.17). Accurate to ~7 decimal places.
func normalSurvival(z float64) float64 {
	if z < 0 {
		return 1.0 - normalSurvival(-z)
	}
	const (
		p  = 0.2316419
		b1 = 0.319381530
		b2 = -0.356563782
		b3 = 1.781477937
		b4 = -1.821255978
		b5 = 1.330274429
	)
	t := 1.0 / (1.0 + p*z)
	density := math.Exp(-0.5*z*z) / sqrtFloat(2*math.Pi)
	cum := density * (b1*t + b2*t*t + b3*t*t*t + b4*t*t*t*t + b5*t*t*t*t*t)
	return cum
}

// sortedKeys returns map keys in deterministic alphabetical order.
func sortedKeys(m map[string]*metricAcc) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// sqrtFloat wraps math.Sqrt. Guards against negative input from numerical
// underflow on near-constant series.
func sqrtFloat(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return math.Sqrt(x)
}

// fullNOrDefault returns store.fullN[id] when available and > 0, otherwise 2.
// Used by effect-size calculations that require at least one degree of freedom
// for pooled SD; 2 is the minimum valid sample count for pooledStandardDeviation.
func fullNOrDefault(store *seriesStore, id int) int {
	if id < len(store.fullN) && store.fullN[id] > 0 {
		return store.fullN[id]
	}
	return 2
}

// clampFloat01 clamps x to [0, 1]. Used for p-values and similar probabilities
// where floating-point arithmetic can produce values slightly outside [0,1].
func clampFloat01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

// roundFloat rounds to n decimal places. Used to compact the score table for
// 6500 metrics (~200KB output cap).
func roundFloat(x float64, decimals int) float64 {
	if decimals < 0 {
		return x
	}
	mul := math.Pow(10, float64(decimals))
	return math.Round(x*mul) / mul
}

// emitGapFindings runs Stage 5 gap detection. Per forensics, the modal FTDC
// sample interval is 1.000s ±0.8% in real captures; a fixed 5s threshold
// avoids false positives on natural jitter while still surfacing real
// process stalls (which are usually multi-second).
func emitGapFindings(c *findingCollector, host string, samples []Sample, cfg AutoConfig) {
	if len(samples) < 2 {
		return
	}
	threshold := time.Duration(cfg.GapThresholdSec * float64(time.Second))
	if threshold <= 0 {
		threshold = 5 * time.Second
	}
	for i := 1; i < len(samples); i++ {
		gap := samples[i].Timestamp.Sub(samples[i-1].Timestamp)
		if gap > threshold {
			at := samples[i-1].Timestamp.UTC().Format(tsFormat)
			c.add(newFinding(KindGap, host, at, map[string]interface{}{
				"at":               at,
				"end_at":           samples[i].Timestamp.UTC().Format(tsFormat),
				"duration_seconds": roundFloat(gap.Seconds(), 3),
				"prior_sample_idx": i - 1,
				"next_sample_idx":  i,
			}))
		}
	}
}

// emitSchemaFingerprint produces a single kind:"schema_fingerprint" Finding
// per host. The fingerprint hashes the sorted (metric_name, inferred_kind)
// tuples — used by the analyzer skill to detect metric-set drift or
// gauge↔counter reclassification across runs and host pairs. Kind inference
// here is a quick monotonicity scan; correctness matters less than
// determinism since the fingerprint's value is detecting CHANGE across runs.
func emitSchemaFingerprint(c *findingCollector, host string, samples []Sample, metadata []MetadataDoc) {
	if len(samples) == 0 {
		return
	}
	keys := metricKeysFromSamples(samples)
	prev := map[string]int64{}
	hasPrev := map[string]bool{}
	sawDecrease := map[string]bool{}
	minV := map[string]int64{}
	maxV := map[string]int64{}
	for i := range samples {
		for name, v := range samples[i].Metrics {
			if hasPrev[name] && v < prev[name] {
				sawDecrease[name] = true
			}
			prev[name] = v
			hasPrev[name] = true
			if mv, ok := minV[name]; !ok || v < mv {
				minV[name] = v
			}
			if mv, ok := maxV[name]; !ok || v > mv {
				maxV[name] = v
			}
		}
	}
	sortedNames := make([]string, 0, len(keys))
	for k := range keys {
		sortedNames = append(sortedNames, k)
	}
	sort.Strings(sortedNames)

	h := sha256.New()
	for _, name := range sortedNames {
		kind := "gauge"
		if minV[name] == maxV[name] {
			kind = "constant"
		} else if !sawDecrease[name] {
			kind = "counter"
		}
		fmt.Fprintf(h, "%s|%s\n", name, kind)
	}
	sum := h.Sum(nil)
	fp := hex.EncodeToString(sum[:8])

	// Echo server version + buildInfo from metadata when present so the
	// analyzer skill's version-scoping notes can be applied. We don't parse
	// the bson.D here; just record the count and let the skill consult
	// --metadata for full details if needed.
	c.add(newFinding(KindSchemaFingerprint, host, "global", map[string]interface{}{
		"fingerprint":        fp,
		"metric_count":       len(sortedNames),
		"metadata_doc_count": len(metadata),
	}))
}
