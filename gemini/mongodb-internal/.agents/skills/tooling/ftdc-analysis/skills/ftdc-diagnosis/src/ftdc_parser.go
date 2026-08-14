// ftdc_parser.go - MongoDB FTDC parser with filtering, rates, and smart summaries
// Build: cd .agents/skills/ftdc-parser && go build -o ftdc_parser .
// Usage: ./ftdc_parser <path> [options]
package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// signalSetup installs the process-wide signal handlers. Called from main()
// EARLY (before bw is created) for SIGPIPE; the SIGINT/SIGTERM hook is
// registered after main() owns bw, via installInterruptFlusher.
//
//   - SIGPIPE: ignored. Writes to a closed stdout return EPIPE; downstream
//     emit sites must check for it and exit cleanly (see usage in
//     analyze.go and other emit paths).
func signalSetup() {
	signal.Ignore(syscall.SIGPIPE)
}

// installInterruptFlusher wires SIGINT/SIGTERM to a goroutine that flushes
// `bw` BEFORE exiting, so partial output is preserved. Pre-Tier-A.2 this
// hook called os.Exit(130) directly, bypassing all deferred flushes.
//
// Idempotent: only the FIRST call registers a signal handler. main() has
// several early-exit branches that each create their own `bw`; without
// this guard, multiple signal.Notify registrations would stack, and the
// late-registered handlers would never observe the original `bw` (each
// goroutine captures the bw passed at its registration). Storing the
// installed bw in a sync.Once-guarded global keeps the flush target
// up-to-date as later main() phases replace bw.
var (
	interruptOnce sync.Once
	interruptBw   *bufio.Writer
	interruptMu   sync.Mutex
)

func installInterruptFlusher(bw *bufio.Writer) {
	interruptMu.Lock()
	interruptBw = bw
	interruptMu.Unlock()
	interruptOnce.Do(func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			sig := <-ch
			fmt.Fprintf(os.Stderr, "ftdc_parser: received signal %v, flushing and exiting\n", sig)
			interruptMu.Lock()
			cur := interruptBw
			interruptMu.Unlock()
			if cur != nil {
				_ = cur.Flush()
			}
			os.Exit(130)
		}()
	})
}

// isEPIPE returns true if `err` is or wraps syscall.EPIPE — signals a closed
// downstream pipe (e.g. `| head` exited). Emit sites use this to exit 0.
func isEPIPE(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, syscall.EPIPE)
}

type Sample struct {
	Timestamp time.Time        `json:"ts"`
	Metrics   map[string]int64 `json:"metrics"`
}

type MetadataDoc struct {
	DocType int    `json:"type"`
	Doc     bson.D `json:"doc"`
}

type RateSample struct {
	Timestamp time.Time          `json:"ts"`
	Rates     map[string]float64 `json:"rates"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Metric filter (literal fast-path + regex fallback)
// ──────────────────────────────────────────────────────────────────────────────

// metricMatcher decides whether a metric key passes a user --filter pattern.
// nil matcher means "accept everything" (no --filter supplied).
type metricMatcher struct {
	literals []string         // pre-lowered, fast ASCII-fold substring
	regexps  []*regexp.Regexp // patterns with regex metachars (case-insensitive)
}

func (m *metricMatcher) matches(name string) bool {
	if m == nil {
		return true
	}
	for _, l := range m.literals {
		if containsFold(name, l) {
			return true
		}
	}
	for _, r := range m.regexps {
		if r.MatchString(name) {
			return true
		}
	}
	return false
}

func compileMatcher(filters []string) (*metricMatcher, error) {
	if len(filters) == 0 {
		return nil, nil
	}
	m := &metricMatcher{}
	for _, f := range filters {
		if isLiteralPattern(f) {
			m.literals = append(m.literals, strings.ToLower(literalText(f)))
		} else {
			p, err := regexp.Compile("(?i)" + f)
			if err != nil {
				return nil, fmt.Errorf("invalid filter regex %q: %v", f, err)
			}
			m.regexps = append(m.regexps, p)
		}
	}
	return m, nil
}

// isLiteralPattern: pattern consists only of alnum, _, -, space, /, ., or \.
// Such patterns are equivalent to ASCII-fold substring search against
// FTDC metric paths (which only use `.` as a path separator). Anything
// containing other regex metachars (* + ? ( ) [ ] { } ^ $ |, or non-`.`
// backslash escapes) routes to the regexp engine for exact compatibility.
//
// Semantic note: in a real regex, `.` matches any character, so
// `wiredTiger.cache` would match a hypothetical metric `wiredTigerXcache`.
// The fast path treats `.` as a literal dot, which is what users almost
// always intend on FTDC paths (segments are separated by `.` only). Empirically
// no real FTDC key contains characters in positions that would cause the two
// to diverge; a `--filter` that requires regex `.` semantics can still force
// the regex engine by using a metachar like `wiredTiger[.]cache`.
func isLiteralPattern(p string) bool {
	for i := 0; i < len(p); i++ {
		c := p[i]
		if c == '\\' {
			if i+1 < len(p) && p[i+1] == '.' {
				i++
				continue
			}
			return false
		}
		switch {
		case c >= 'a' && c <= 'z':
		case c >= 'A' && c <= 'Z':
		case c >= '0' && c <= '9':
		case c == '_' || c == '-' || c == ' ' || c == '/' || c == '.':
		default:
			return false
		}
	}
	return true
}

// literalText extracts the literal text from a literal pattern: \. → .
func literalText(p string) string {
	var b strings.Builder
	b.Grow(len(p))
	for i := 0; i < len(p); i++ {
		c := p[i]
		if c == '\\' && i+1 < len(p) && p[i+1] == '.' {
			b.WriteByte('.')
			i++
			continue
		}
		b.WriteByte(c)
	}
	return b.String()
}

// containsFold reports whether the pre-lowered ASCII substring `sub` appears
// in `s`, comparing case-insensitively. Allocation-free.
func containsFold(s, sub string) bool {
	n := len(sub)
	if n == 0 {
		return true
	}
	if n > len(s) {
		return false
	}
	last := len(s) - n
	for i := 0; i <= last; i++ {
		ok := true
		for j := 0; j < n; j++ {
			c := s[i+j]
			if c >= 'A' && c <= 'Z' {
				c += 32
			}
			if c != sub[j] {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}

// ──────────────────────────────────────────────────────────────────────────────
// Varint + BSON flattening (unchanged semantics)
// ──────────────────────────────────────────────────────────────────────────────

func readVarint(r *bytes.Reader) (int64, error) {
	var result uint64
	var shift uint
	// Protobuf varints are at most 10 bytes for a 64-bit value; the 10th
	// byte contributes only 1 bit. Bounding the iteration count prevents
	// hostile inputs (e.g. an unterminated 0x80 stream) from looping
	// forever shifting past 64.
	for i := 0; i < 10; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		result |= uint64(b&0x7F) << shift
		if b&0x80 == 0 {
			return int64(result), nil
		}
		shift += 7
	}
	return 0, fmt.Errorf("varint too long")
}

// flattenBSONOrdered extracts keys and values from BSON doc preserving insertion order.
// This is CRITICAL for FTDC since delta decoding must match the exact order in the reference doc.
func flattenBSONOrdered(doc bson.D, prefix string, keys *[]string, vals *[]int64) {
	flattenBSONOrderedDepth(doc, prefix, keys, vals, 0)
}

// maxBSONDepth bounds recursion for hostile inputs. BSON spec allows arbitrary
// nesting but real serverStatus docs top out around 8 levels.
const maxBSONDepth = 100

func flattenBSONOrderedDepth(doc bson.D, prefix string, keys *[]string, vals *[]int64, depth int) {
	if depth >= maxBSONDepth {
		return
	}
	for _, elem := range doc {
		path := elem.Key
		if prefix != "" {
			path = prefix + "." + elem.Key
		}
		switch v := elem.Value.(type) {
		case int32:
			*keys = append(*keys, path)
			*vals = append(*vals, int64(v))
		case int64:
			*keys = append(*keys, path)
			*vals = append(*vals, v)
		case float64:
			*keys = append(*keys, path)
			// Match MongoDB C++ behavior: clamp NaN/Inf to safe int64 values
			var iv int64
			if math.IsNaN(v) {
				iv = 0
			} else if v >= float64(math.MaxInt64) {
				iv = math.MaxInt64
			} else if v <= float64(math.MinInt64) {
				iv = math.MinInt64
			} else {
				iv = int64(v)
			}
			*vals = append(*vals, iv)
		case bool:
			*keys = append(*keys, path)
			if v {
				*vals = append(*vals, 1)
			} else {
				*vals = append(*vals, 0)
			}
		case primitive.DateTime:
			*keys = append(*keys, path)
			*vals = append(*vals, int64(v))
		case primitive.Timestamp:
			*keys = append(*keys, path+".t")
			*vals = append(*vals, int64(v.T))
			*keys = append(*keys, path+".i")
			*vals = append(*vals, int64(v.I))
		case primitive.Decimal128:
			*keys = append(*keys, path)
			*vals = append(*vals, decimal128ToInt64(v))
		case bson.D:
			flattenBSONOrderedDepth(v, path, keys, vals, depth+1)
		case bson.A:
			for i, av := range v {
				if subdoc, ok := av.(bson.D); ok {
					flattenBSONOrderedDepth(subdoc, fmt.Sprintf("%s.%d", path, i), keys, vals, depth+1)
				}
			}
		}
	}
}

// decimal128ToInt64 converts a BSON Decimal128 to int64 using bounded
// scaling. Decimal128 has max exponent ±6144, but applying any exponent
// beyond ~18 overflows int64. We cap iterations at 34 (a safe ceiling that
// will never overflow a value with magnitude < 1) and clamp to ±MaxInt64
// on overflow. Negative exponents (fractional decimals) truncate toward
// zero by dividing — same convention as MongoDB's C++ FTDC sampler.
func decimal128ToInt64(v primitive.Decimal128) int64 {
	bi, exp, err := v.BigInt()
	if err != nil {
		return 0 // Inf, NaN, or unrepresentable
	}
	if exp == 0 {
		if !bi.IsInt64() {
			if bi.Sign() < 0 {
				return math.MinInt64
			}
			return math.MaxInt64
		}
		return bi.Int64()
	}
	if exp > 0 {
		const maxScale = 34
		if exp > maxScale {
			exp = maxScale
		}
		if !bi.IsInt64() {
			if bi.Sign() < 0 {
				return math.MinInt64
			}
			return math.MaxInt64
		}
		iv := bi.Int64()
		for e := 0; e < exp; e++ {
			next := iv * 10
			if iv != 0 && next/10 != iv {
				// overflow
				if iv < 0 {
					return math.MinInt64
				}
				return math.MaxInt64
			}
			iv = next
		}
		return iv
	}
	// exp < 0: divide by 10^(-exp), truncating toward zero.
	if !bi.IsInt64() {
		// Magnitude too large; fractional portion irrelevant.
		if bi.Sign() < 0 {
			return math.MinInt64
		}
		return math.MaxInt64
	}
	iv := bi.Int64()
	for e := 0; e > exp && iv != 0; e-- {
		iv /= 10
	}
	return iv
}

// ──────────────────────────────────────────────────────────────────────────────
// Chunk decoding with filter pushdown
// ──────────────────────────────────────────────────────────────────────────────

// chunkJob is a single zlib-compressed FTDC chunk awaiting decode.
type chunkJob struct {
	idx      int       // sequential position within the source file (defines time order)
	data     []byte    // compressed chunk bytes (sub-slice of file mmap/read)
	baseTime time.Time // BSON _id of the top-level type-1 doc
}

// decodeChunk decompresses one FTDC chunk and produces its samples.
// If matcher is non-nil, only matching metric keys are populated in each
// sample's Metrics map. The "start" key is always materialized in the
// internal matrix (timestamps depend on it) but is only added to
// Sample.Metrics when the matcher would have included it (preserving the
// pre-change contract that filtered-out "start" does not appear in output).
func decodeChunk(data []byte, baseTime time.Time, matcher *metricMatcher) []Sample {
	if len(data) < 4 {
		return nil
	}

	r, err := zlib.NewReader(bytes.NewReader(data[4:]))
	if err != nil {
		return nil
	}
	defer r.Close()

	// Bound decompressed size. A malicious or corrupted chunk could otherwise
	// inflate to gigabytes and OOM the process. 256 MB is ~10× the largest
	// natural FTDC chunk observed (24h × 300-sample × 10K-key captures).
	const maxChunkUncompressed = 256 << 20
	limited := &io.LimitedReader{R: r, N: maxChunkUncompressed + 1}
	uncompressed, _ := io.ReadAll(limited)
	if len(uncompressed) < 8 || int64(len(uncompressed)) > maxChunkUncompressed {
		return nil
	}

	refDocLen := int(binary.LittleEndian.Uint32(uncompressed))
	if refDocLen < 5 || refDocLen > len(uncompressed) {
		return nil
	}

	var refDoc bson.D
	if err := bson.Unmarshal(uncompressed[:refDocLen], &refDoc); err != nil {
		return nil
	}

	afterDoc := uncompressed[refDocLen:]
	if len(afterDoc) < 8 {
		return nil
	}

	nKeys := int(binary.LittleEndian.Uint32(afterDoc[0:4]))
	nDeltas := int(binary.LittleEndian.Uint32(afterDoc[4:8]))
	// FTDC's natural ceiling on samples-per-chunk is 300 (one 5-minute window
	// at 1Hz). Reject chunks that claim more than 10K — clearly hostile or
	// corrupted. This avoids huge matrix allocations downstream.
	const maxDeltasPerChunk = 10000
	if nDeltas > maxDeltasPerChunk {
		return nil
	}
	deltas := afterDoc[8:]

	var keys []string
	var baseVals []int64
	flattenBSONOrdered(refDoc, "", &keys, &baseVals)

	if nKeys > len(keys) {
		nKeys = len(keys)
	}

	// Compute keep-masks once for this chunk's key list.
	//   keepInMap   — whether key appears in output Sample.Metrics
	//   keepInMatrix — whether we accumulate delta values for this key
	// "start" is always kept in matrix (timestamps); only in map if matcher allows.
	startIdx := -1
	var keepInMap []bool
	var keepInMatrix []bool
	keepCount := 0
	if matcher != nil {
		keepInMap = make([]bool, nKeys)
		keepInMatrix = make([]bool, nKeys)
		for i := 0; i < nKeys; i++ {
			if keys[i] == "start" && startIdx < 0 {
				startIdx = i
			}
			if matcher.matches(keys[i]) {
				keepInMap[i] = true
				keepInMatrix[i] = true
				keepCount++
			}
		}
		if startIdx >= 0 {
			keepInMatrix[startIdx] = true
		}
	} else {
		// No filter: everything kept.
		for i := 0; i < nKeys; i++ {
			if keys[i] == "start" && startIdx < 0 {
				startIdx = i
				break
			}
		}
		keepCount = nKeys
	}

	// Allocate matrix rows only for kept-in-matrix keys.
	matrix := make([][]int64, nKeys)
	for i := 0; i < nKeys; i++ {
		if keepInMatrix == nil || keepInMatrix[i] {
			matrix[i] = make([]int64, nDeltas+1)
			matrix[i][0] = baseVals[i]
		}
	}

	// Walk the delta stream. We MUST consume every varint to maintain the
	// RLE zero-run state correctly, even for keys we don't keep.
	dr := bytes.NewReader(deltas)
	var zeroesCount uint64
	for i := 0; i < nKeys; i++ {
		row := matrix[i] // nil if this key isn't kept
		var prev int64
		if row != nil {
			prev = baseVals[i]
		}
		for j := 1; j <= nDeltas; j++ {
			var delta int64
			if zeroesCount > 0 {
				delta = 0
				zeroesCount--
			} else {
				d, err := readVarint(dr)
				if err != nil {
					break
				}
				delta = d
				if delta == 0 {
					z, err := readVarint(dr)
					if err == nil {
						zeroesCount = uint64(z)
					}
				}
			}
			if row != nil {
				prev += delta
				row[j] = prev
			}
		}
	}

	// Build the sample list. Cap each Metrics map at keepCount (with +1 slack
	// to allow the rare case where map auto-grows during the loop).
	mapCap := keepCount
	if mapCap < 1 {
		mapCap = 1
	}
	samples := make([]Sample, nDeltas+1)
	for j := 0; j <= nDeltas; j++ {
		samples[j].Metrics = make(map[string]int64, mapCap)
		samples[j].Timestamp = baseTime.Add(time.Duration(j) * time.Second)
		if startIdx >= 0 && matrix[startIdx] != nil {
			samples[j].Timestamp = time.UnixMilli(matrix[startIdx][j])
		}
		for i := 0; i < nKeys; i++ {
			if matrix[i] == nil {
				continue
			}
			if keepInMap != nil && !keepInMap[i] {
				continue // matrix kept only for timestamps; user filtered out
			}
			samples[j].Metrics[keys[i]] = matrix[i][j]
		}
	}
	return samples
}

// ──────────────────────────────────────────────────────────────────────────────
// File / directory ingestion with chunk-level parallelism
// ──────────────────────────────────────────────────────────────────────────────

// partitionFTDC walks the top-level BSON envelope and splits the file into
// chunk jobs (type-1 docs) + metadata docs (type-0, type-2). This is serial
// and cheap; the heavy work (zlib + varint + matrix) happens in decodeChunk.
func partitionFTDC(data []byte) ([]chunkJob, []MetadataDoc) {
	var jobs []chunkJob
	var metadata []MetadataDoc
	offset := 0
	for offset < len(data) {
		if len(data[offset:]) < 5 {
			break
		}
		docLen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
		if docLen < 5 || offset+docLen > len(data) {
			break
		}

		var doc bson.D
		if err := bson.Unmarshal(data[offset:offset+docLen], &doc); err != nil {
			offset += docLen
			continue
		}
		end := offset + docLen
		offset = end

		var docType int32
		var binData primitive.Binary
		var baseTime time.Time
		var innerDoc bson.D

		for _, elem := range doc {
			switch elem.Key {
			case "type":
				if v, ok := elem.Value.(int32); ok {
					docType = v
				}
			case "data":
				if v, ok := elem.Value.(primitive.Binary); ok {
					binData = v
				}
			case "doc":
				if v, ok := elem.Value.(bson.D); ok {
					innerDoc = v
				}
			case "_id":
				if dt, ok := elem.Value.(primitive.DateTime); ok {
					baseTime = dt.Time()
				}
			}
		}

		if docType == 1 && len(binData.Data) > 0 {
			jobs = append(jobs, chunkJob{
				idx:      len(jobs),
				data:     binData.Data,
				baseTime: baseTime,
			})
		} else if (docType == 0 || docType == 2) && innerDoc != nil {
			metadata = append(metadata, MetadataDoc{DocType: int(docType), Doc: innerDoc})
		}
	}
	return jobs, metadata
}

// chunkResult carries decoded samples back from a worker, tagged with the
// chunk index so the reorder buffer can emit in strict time order.
type chunkResult struct {
	idx     int
	samples []Sample
}

// processFTDCFile streams the chunks of one file in chunk-index (time) order.
// onBatch is called once per chunk's sample slice. Workers > 1 enables
// parallel chunk decode behind a bounded reorder buffer.
//
// onBatch contract: invoked exclusively from the calling goroutine and MUST
// NOT block indefinitely. A blocking onBatch would stall worker progress via
// resCh back-pressure (workers would fill resCh and block, the closer
// goroutine would wait on wg.Wait() forever). Today's callers just append
// to a slice; if you ever wire onBatch to streaming I/O, ensure it cannot
// deadlock on resource the workers also depend on.
func processFTDCFile(path string, matcher *metricMatcher, workers int, onBatch func([]Sample)) []MetadataDoc {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	jobs, metadata := partitionFTDC(raw)
	if len(jobs) == 0 {
		return metadata
	}

	if workers < 1 {
		workers = 1
	}
	if workers > len(jobs) {
		workers = len(jobs)
	}

	// Serial fast path: identical output, zero goroutine overhead.
	if workers == 1 {
		for _, j := range jobs {
			samples := decodeChunk(j.data, j.baseTime, matcher)
			if len(samples) > 0 {
				onBatch(samples)
			}
		}
		return metadata
	}

	// Parallel: fan-out workers, fan-in via reorder buffer.
	jobCh := make(chan int, workers)
	resCh := make(chan chunkResult, workers)

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			// Tier A.2: recover from panics in decode so one bad chunk
			// doesn't crash the whole process. Pre-A.2, a panic here
			// would unwind through the goroutine and terminate main.
			for idx := range jobCh {
				func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Fprintf(os.Stderr,
								"ftdc_parser: worker %d recovered from chunk %d panic: %v\n",
								workerID, idx, r)
							resCh <- chunkResult{idx: idx, samples: nil}
						}
					}()
					samples := decodeChunk(jobs[idx].data, jobs[idx].baseTime, matcher)
					resCh <- chunkResult{idx: idx, samples: samples}
				}()
			}
		}(w)
	}

	// Feeder: pushes job indices in order. Bounded buffer means at most
	// `workers` jobs are queued ahead of an idle worker.
	go func() {
		for i := range jobs {
			jobCh <- i
		}
		close(jobCh)
	}()

	// Closer: closes resCh once all workers exit.
	go func() {
		wg.Wait()
		close(resCh)
	}()

	// Reorder + emit. At most `workers` chunks can be pending behind the
	// next-to-emit cursor, so peak memory is O(workers × per-chunk samples).
	pending := make(map[int][]Sample, workers)
	nextEmit := 0
	for r := range resCh {
		pending[r.idx] = r.samples
		for {
			s, ok := pending[nextEmit]
			if !ok {
				break
			}
			if len(s) > 0 {
				onBatch(s)
			}
			delete(pending, nextEmit)
			nextEmit++
		}
	}
	return metadata
}

// listFTDCDir returns FTDC files in a directory in time order.
// Two-tier selection:
//  1. Prefer files with the standard "metrics." prefix (mongod $dbpath form).
//  2. If none are present, fall back to any non-hidden file (Atlas bundles
//     drop the prefix and use raw Unix-timestamp filenames).
// Both naming conventions sort chronologically by name.
func listFTDCDir(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var preferred, fallback []string
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || strings.HasPrefix(name, ".") {
			continue
		}
		full := filepath.Join(dir, name)
		if strings.HasPrefix(name, "metrics.") {
			preferred = append(preferred, full)
		} else {
			fallback = append(fallback, full)
		}
	}
	files := preferred
	if len(files) == 0 {
		files = fallback
	}
	sort.Strings(files)
	return files
}

// processFTDC dispatches a single file or a whole directory through
// processFTDCFile, preserving cross-file time order by processing files in
// name order. Chunk-level parallelism happens within each file.
func processFTDC(path string, matcher *metricMatcher, workers int, onBatch func([]Sample)) []MetadataDoc {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	var files []string
	if info.IsDir() {
		files = listFTDCDir(path)
	} else {
		files = []string{path}
	}
	var allMeta []MetadataDoc
	for _, f := range files {
		meta := processFTDCFile(f, matcher, workers, onBatch)
		allMeta = append(allMeta, meta...)
	}
	return allMeta
}

// ──────────────────────────────────────────────────────────────────────────────
// Post-decode filters (time range, all-zero)
// ──────────────────────────────────────────────────────────────────────────────

func stripZeroMetrics(samples []Sample) []Sample {
	if len(samples) == 0 {
		return samples
	}
	nonZero := make(map[string]bool)
	for i := range samples {
		for k, v := range samples[i].Metrics {
			if v != 0 {
				nonZero[k] = true
			}
		}
	}
	for i := range samples {
		filtered := make(map[string]int64, len(nonZero))
		for k, v := range samples[i].Metrics {
			if nonZero[k] {
				filtered[k] = v
			}
		}
		samples[i].Metrics = filtered
	}
	return samples
}

func computeRates(samples []Sample) []RateSample {
	if len(samples) < 2 {
		return nil
	}
	rates := make([]RateSample, len(samples)-1)
	for i := 1; i < len(samples); i++ {
		dtSec := samples[i].Timestamp.Sub(samples[i-1].Timestamp).Seconds()
		if dtSec <= 0 {
			dtSec = 1
		}
		r := make(map[string]float64)
		for k, v := range samples[i].Metrics {
			if prev, ok := samples[i-1].Metrics[k]; ok {
				// Negative deltas are preserved: gauges legitimately decrease, and counter
				// wraps/resets are shown as large negative values for downstream detection.
				delta := v - prev
				r[k] = float64(delta) / dtSec
			}
		}
		rates[i-1] = RateSample{
			Timestamp: samples[i].Timestamp,
			Rates:     r,
		}
	}
	return rates
}

type metricStats struct {
	name        string
	min, max    int64
	first, last int64
	totalDelta  int64
	sum         float64
	count       int
	changed     bool
}

func computeMetricStats(samples []Sample) []metricStats {
	if len(samples) == 0 {
		return nil
	}
	statsMap := make(map[string]*metricStats)

	for k, v := range samples[0].Metrics {
		statsMap[k] = &metricStats{name: k, min: v, max: v, first: v, last: v, sum: float64(v), count: 1}
	}

	for i := 1; i < len(samples); i++ {
		for k, v := range samples[i].Metrics {
			s, ok := statsMap[k]
			if !ok {
				statsMap[k] = &metricStats{name: k, min: v, max: v, first: v, last: v, sum: float64(v), count: 1}
				continue
			}
			if v < s.min {
				s.min = v
			}
			if v > s.max {
				s.max = v
			}
			s.last = v
			s.sum += float64(v)
			s.count++
		}
	}

	for _, s := range statsMap {
		s.totalDelta = s.last - s.first
		s.changed = s.min != s.max
	}

	result := make([]metricStats, 0, len(statsMap))
	for _, s := range statsMap {
		result = append(result, *s)
	}
	return result
}

// ──────────────────────────────────────────────────────────────────────────────
// Pretty printing (now all output goes through io.Writer for buffered stdout)
// ──────────────────────────────────────────────────────────────────────────────

func formatInt(v int64) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	s := fmt.Sprintf("%d", v)
	n := len(s)
	if n <= 3 {
		if neg {
			return "-" + s
		}
		return s
	}
	var buf strings.Builder
	rem := n % 3
	if rem > 0 {
		buf.WriteString(s[:rem])
	}
	for i := rem; i < n; i += 3 {
		if buf.Len() > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(s[i : i+3])
	}
	if neg {
		return "-" + buf.String()
	}
	return buf.String()
}

func formatBytes(v int64) string {
	if v < 1024 {
		return fmt.Sprintf("%d B", v)
	} else if v < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(v)/1024)
	} else if v < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(v)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", float64(v)/(1024*1024*1024))
	}
}

func formatRate(v float64) string {
	abs := math.Abs(v)
	if abs >= 1e6 {
		return fmt.Sprintf("%.1fM/s", v/1e6)
	} else if abs >= 1e3 {
		return fmt.Sprintf("%.1fK/s", v/1e3)
	} else if abs >= 1 {
		return fmt.Sprintf("%.1f/s", v)
	} else if abs >= 0.01 {
		return fmt.Sprintf("%.2f/s", v)
	}
	return fmt.Sprintf("%.4f/s", v)
}

func printSummary(w io.Writer, samples []Sample, metadata []MetadataDoc) {
	if len(samples) == 0 {
		fmt.Fprintln(w, "No samples found")
		return
	}

	first, last := samples[0], samples[len(samples)-1]
	durSec := last.Timestamp.Sub(first.Timestamp).Seconds()

	fmt.Fprintf(w, "FTDC Summary\n")
	fmt.Fprintf(w, "============\n")

	// Extract key metadata
	for _, md := range metadata {
		if md.DocType != 0 {
			continue
		}
		m := bsonToMap(md.Doc)
		if bi, ok := m["buildInfo"].(map[string]interface{}); ok {
			if v, ok := bi["version"]; ok {
				fmt.Fprintf(w, "Server:   %v\n", v)
			}
			if v, ok := bi["modules"]; ok {
				fmt.Fprintf(w, "Modules:  %v\n", v)
			}
		}
		if hi, ok := m["hostInfo"].(map[string]interface{}); ok {
			if sys, ok := hi["system"].(map[string]interface{}); ok {
				fmt.Fprintf(w, "Host:     %v cores, %v MB RAM, %v\n",
					sys["numCores"], sys["memSizeMB"], sys["cpuArch"])
			}
		}
		if opts, ok := m["getCmdLineOpts"].(map[string]interface{}); ok {
			if parsed, ok := opts["parsed"].(map[string]interface{}); ok {
				if sp, ok := parsed["setParameter"].(map[string]interface{}); ok {
					if v, ok := sp["disaggregatedStorageEnabled"]; ok {
						fmt.Fprintf(w, "Disagg:   %v\n", v)
					}
				}
				if st, ok := parsed["storage"].(map[string]interface{}); ok {
					if wt, ok := st["wiredTiger"].(map[string]interface{}); ok {
						if ec, ok := wt["engineConfig"].(map[string]interface{}); ok {
							if cs, ok := ec["configString"]; ok {
								fmt.Fprintf(w, "WT cfg:   %v\n", cs)
							}
						}
					}
				}
				if repl, ok := parsed["replication"].(map[string]interface{}); ok {
					if rs, ok := repl["replSetName"]; ok {
						fmt.Fprintf(w, "ReplSet:  %v\n", rs)
					}
				}
			}
		}
	}
	fmt.Fprintf(w, "Samples:  %s\n", formatInt(int64(len(samples))))
	fmt.Fprintf(w, "Duration: %.0fs (%.1f min)\n", durSec, durSec/60)
	fmt.Fprintf(w, "Start:    %s\n", first.Timestamp.UTC().Format(time.RFC3339))
	fmt.Fprintf(w, "End:      %s\n", last.Timestamp.UTC().Format(time.RFC3339))
	fmt.Fprintf(w, "Metrics:  %s total\n", formatInt(int64(len(first.Metrics))))

	stats := computeMetricStats(samples)

	changed := 0
	for _, s := range stats {
		if s.changed {
			changed++
		}
	}
	fmt.Fprintf(w, "          %s changed, %s static\n", formatInt(int64(changed)), formatInt(int64(len(stats)-changed)))

	// Helper closures
	get := func(s Sample, name string) (int64, bool) {
		v, ok := s.Metrics[name]
		return v, ok
	}
	delta := func(name string) int64 {
		l, lok := last.Metrics[name]
		f, fok := first.Metrics[name]
		if lok && fok {
			return l - f
		}
		return 0
	}
	rate := func(name string) float64 {
		if durSec <= 0 {
			return 0
		}
		return float64(delta(name)) / durSec
	}

	// Operations
	fmt.Fprintf(w, "\n--- Operations (total / avg rate) ---\n")
	ops := []struct{ name, label string }{
		{"serverStatus.opcounters.insert", "insert"},
		{"serverStatus.opcounters.query", "query"},
		{"serverStatus.opcounters.update", "update"},
		{"serverStatus.opcounters.delete", "delete"},
		{"serverStatus.opcounters.getmore", "getmore"},
		{"serverStatus.opcounters.command", "command"},
	}
	for _, op := range ops {
		d := delta(op.name)
		if d > 0 {
			fmt.Fprintf(w, "  %-10s %s  (%s)\n", op.label+":", formatInt(d), formatRate(rate(op.name)))
		}
	}

	// Connections
	fmt.Fprintf(w, "\n--- Connections ---\n")
	if v, ok := get(last, "serverStatus.connections.current"); ok {
		fmt.Fprintf(w, "  current: %s\n", formatInt(v))
	}
	if v, ok := get(last, "serverStatus.connections.totalCreated"); ok {
		fmt.Fprintf(w, "  totalCreated: %s\n", formatInt(v))
	}
	// Connection range
	minConn, maxConn := int64(math.MaxInt64), int64(0)
	for _, s := range samples {
		if v, ok := s.Metrics["serverStatus.connections.current"]; ok {
			if v < minConn {
				minConn = v
			}
			if v > maxConn {
				maxConn = v
			}
		}
	}
	if maxConn > 0 {
		fmt.Fprintf(w, "  range:   [%s, %s]\n", formatInt(minConn), formatInt(maxConn))
	}

	// Cache
	fmt.Fprintf(w, "\n--- WiredTiger Cache ---\n")
	if v, ok := get(last, "serverStatus.wiredTiger.cache.bytes currently in the cache"); ok {
		fmt.Fprintf(w, "  in cache:   %s\n", formatBytes(v))
	}
	if maxBytes, ok := get(last, "serverStatus.wiredTiger.cache.maximum bytes configured"); ok && maxBytes > 0 {
		fmt.Fprintf(w, "  configured: %s\n", formatBytes(maxBytes))
		// Show cache fill range
		minFill, maxFill := 100.0, 0.0
		for _, s := range samples {
			if cur, ok := s.Metrics["serverStatus.wiredTiger.cache.bytes currently in the cache"]; ok {
				pct := float64(cur) / float64(maxBytes) * 100
				if pct < minFill {
					minFill = pct
				}
				if pct > maxFill {
					maxFill = pct
				}
			}
		}
		fmt.Fprintf(w, "  fill %%:     %.1f%% - %.1f%%\n", minFill, maxFill)
	}
	if v, ok := get(last, "serverStatus.wiredTiger.cache.tracked dirty bytes in the cache"); ok {
		fmt.Fprintf(w, "  dirty:      %s\n", formatBytes(v))
	}
	appEvict := delta("serverStatus.wiredTiger.cache.pages evicted by application threads")
	if appEvict > 0 {
		fmt.Fprintf(w, "  app thread evictions: %s  ⚠️\n", formatInt(appEvict))
	}
	appEvictTime := delta("serverStatus.wiredTiger.thread-yield.application thread time evicting (usecs)")
	if appEvictTime > 0 {
		fmt.Fprintf(w, "  app thread eviction time: %.1fs  ⚠️\n", float64(appEvictTime)/1e6)
	}

	// Checkpoint
	fmt.Fprintf(w, "\n--- Checkpoints ---\n")
	if v, ok := get(last, "serverStatus.wiredTiger.checkpoint.generation"); ok {
		fmt.Fprintf(w, "  count:       %s\n", formatInt(v))
	}
	if v, ok := get(last, "serverStatus.wiredTiger.checkpoint.most recent time (msecs)"); ok {
		fmt.Fprintf(w, "  last:        %dms\n", v)
	}
	if v, ok := get(last, "serverStatus.wiredTiger.checkpoint.max time (msecs)"); ok {
		fmt.Fprintf(w, "  max:         %dms\n", v)
	}
	if v, ok := get(last, "serverStatus.wiredTiger.checkpoint.total time (msecs)"); ok {
		fmt.Fprintf(w, "  total:       %.1fs\n", float64(v)/1000)
	}

	// Tickets
	fmt.Fprintf(w, "\n--- Admission Control (last sample) ---\n")
	for _, kind := range []string{"read", "write"} {
		prefix := "serverStatus.queues.execution." + kind
		avail, aok := get(last, prefix+".available")
		out, ook := get(last, prefix+".out")
		total, tok := get(last, prefix+".totalTickets")
		if aok || ook || tok {
			fmt.Fprintf(w, "  %s: available=%s out=%s total=%s\n", kind, formatInt(avail), formatInt(out), formatInt(total))
		}
	}
	// Ticket exhaustion history
	for _, kind := range []string{"read", "write"} {
		prefix := "serverStatus.queues.execution." + kind + ".available"
		exhausted := 0
		for _, s := range samples {
			if v, ok := s.Metrics[prefix]; ok && v == 0 {
				exhausted++
			}
		}
		if exhausted > 0 {
			pct := float64(exhausted) / float64(len(samples)) * 100
			fmt.Fprintf(w, "  %s exhausted: %d/%d samples (%.1f%%)  ⚠️\n", kind, exhausted, len(samples), pct)
		}
	}

	// Replication
	if v, ok := get(last, "replSetGetStatus.myState"); ok {
		stateNames := map[int64]string{1: "PRIMARY", 2: "SECONDARY", 7: "ARBITER"}
		state := fmt.Sprintf("%d", v)
		if name, ok := stateNames[v]; ok {
			state = name
		}
		fmt.Fprintf(w, "\n--- Replication ---\n")
		fmt.Fprintf(w, "  state: %s\n", state)
	}

	// Document metrics
	fmt.Fprintf(w, "\n--- Document Operations (total) ---\n")
	docOps := []struct{ name, label string }{
		{"serverStatus.metrics.document.inserted", "inserted"},
		{"serverStatus.metrics.document.returned", "returned"},
		{"serverStatus.metrics.document.updated", "updated"},
		{"serverStatus.metrics.document.deleted", "deleted"},
	}
	for _, op := range docOps {
		d := delta(op.name)
		if d > 0 {
			fmt.Fprintf(w, "  %-10s %s\n", op.label+":", formatInt(d))
		}
	}

	// System
	fmt.Fprintf(w, "\n--- System ---\n")
	if v, ok := get(last, "systemMetrics.cpu.num_cpus"); ok {
		fmt.Fprintf(w, "  CPUs: %d\n", v)
	}
	if v, ok := get(last, "serverStatus.tcmalloc.generic.current_allocated_bytes"); ok {
		fmt.Fprintf(w, "  tcmalloc allocated: %s\n", formatBytes(v))
	}
	if v, ok := get(last, "serverStatus.tcmalloc.generic.heap_size"); ok {
		fmt.Fprintf(w, "  tcmalloc heap:      %s\n", formatBytes(v))
	}

	// Top 15 most-changed metrics by absolute delta
	isTimestampMetric := func(name string) bool {
		lower := strings.ToLower(name)
		for _, pat := range []string{
			"timestamp", "walltime", "lsn", "expiration time",
			"btree clean tree checkpoint expiration",
			"date", "millis",
		} {
			if strings.Contains(lower, pat) {
				return true
			}
		}
		return false
	}

	fmt.Fprintf(w, "\n--- Top 15 Most Active Metrics (by absolute delta) ---\n")
	// Tier A.4: stable sort with name tiebreaker for deterministic top-15.
	sort.SliceStable(stats, func(i, j int) bool {
		ai := stats[i].totalDelta
		if ai < 0 {
			ai = -ai
		}
		aj := stats[j].totalDelta
		if aj < 0 {
			aj = -aj
		}
		if ai != aj {
			return ai > aj
		}
		return stats[i].name < stats[j].name
	})
	shown := 0
	for _, s := range stats {
		if shown >= 15 {
			break
		}
		if !s.changed {
			break
		}
		if isTimestampMetric(s.name) {
			continue
		}
		shown++
		safeDurSec := durSec
		if safeDurSec <= 0 {
			safeDurSec = 1
		}
		avgRate := float64(s.totalDelta) / safeDurSec
		fmt.Fprintf(w, "  %s\n    delta=%s  rate=%s  range=[%s, %s]\n",
			s.name, formatInt(s.totalDelta), formatRate(avgRate),
			formatInt(s.min), formatInt(s.max))
	}

	// Detect potential issues
	issues := detectIssues(samples, stats)
	if len(issues) > 0 {
		fmt.Fprintf(w, "\n--- Potential Issues Detected ---\n")
		for _, issue := range issues {
			fmt.Fprintf(w, "  ⚠️  %s\n", issue)
		}
	}

	// Detect FTDC gaps
	gaps := detectGaps(samples)
	if len(gaps) > 0 {
		fmt.Fprintf(w, "\n--- FTDC Gaps (missing samples) ---\n")
		for _, g := range gaps {
			fmt.Fprintf(w, "  %s → %s  (%.0fs gap)\n",
				g.start.UTC().Format(time.RFC3339),
				g.end.UTC().Format(time.RFC3339),
				g.duration.Seconds())
		}
	}

	// Detect schema changes (metric count changes between samples)
	schemaChanges := 0
	prevCount := len(samples[0].Metrics)
	for i := 1; i < len(samples); i++ {
		curCount := len(samples[i].Metrics)
		if curCount != prevCount {
			if schemaChanges < 5 {
				fmt.Fprintf(w, "\n--- Schema Change ---\n")
				fmt.Fprintf(w, "  %s: metrics %d → %d\n",
					samples[i].Timestamp.UTC().Format(time.RFC3339),
					prevCount, curCount)
			}
			schemaChanges++
			prevCount = curCount
		}
	}
	if schemaChanges > 5 {
		fmt.Fprintf(w, "  ... and %d more schema changes\n", schemaChanges-5)
	}
}

type gap struct {
	start    time.Time
	end      time.Time
	duration time.Duration
}

func detectGaps(samples []Sample) []gap {
	var gaps []gap
	for i := 1; i < len(samples); i++ {
		d := samples[i].Timestamp.Sub(samples[i-1].Timestamp)
		if d > 5*time.Second {
			gaps = append(gaps, gap{
				start:    samples[i-1].Timestamp,
				end:      samples[i].Timestamp,
				duration: d,
			})
		}
	}
	return gaps
}

func detectIssues(samples []Sample, stats []metricStats) []string {
	if len(samples) < 2 {
		return nil
	}
	var issues []string

	last := samples[len(samples)-1]
	first := samples[0]
	durSec := last.Timestamp.Sub(first.Timestamp).Seconds()
	if durSec <= 0 {
		return nil
	}

	// Cache fill ratio
	if cur, ok := last.Metrics["serverStatus.wiredTiger.cache.bytes currently in the cache"]; ok {
		if max, ok := last.Metrics["serverStatus.wiredTiger.cache.maximum bytes configured"]; ok && max > 0 {
			pct := float64(cur) / float64(max) * 100
			if pct > 95 {
				issues = append(issues, fmt.Sprintf("Cache fill at %.1f%% (eviction trigger threshold)", pct))
			} else if pct > 80 {
				issues = append(issues, fmt.Sprintf("Cache fill at %.1f%% (above eviction target)", pct))
			}
		}
	}

	// App thread evictions
	for _, s := range stats {
		if s.name == "serverStatus.wiredTiger.cache.pages evicted by application threads" && s.totalDelta > 0 {
			issues = append(issues, fmt.Sprintf("Application thread evictions: %s — latency impact likely", formatInt(s.totalDelta)))
		}
	}

	// Write ticket exhaustion
	exhaustedCount := 0
	for _, s := range samples {
		if v, ok := s.Metrics["serverStatus.queues.execution.write.available"]; ok && v == 0 {
			exhaustedCount++
		}
	}
	if exhaustedCount > 0 {
		pct := float64(exhaustedCount) / float64(len(samples)) * 100
		issues = append(issues, fmt.Sprintf("Write tickets exhausted in %d/%d samples (%.1f%%)", exhaustedCount, len(samples), pct))
	}

	// Read ticket exhaustion
	exhaustedCount = 0
	for _, s := range samples {
		if v, ok := s.Metrics["serverStatus.queues.execution.read.available"]; ok && v == 0 {
			exhaustedCount++
		}
	}
	if exhaustedCount > 0 {
		pct := float64(exhaustedCount) / float64(len(samples)) * 100
		issues = append(issues, fmt.Sprintf("Read tickets exhausted in %d/%d samples (%.1f%%)", exhaustedCount, len(samples), pct))
	}

	// Long checkpoints
	for _, s := range stats {
		if s.name == "serverStatus.wiredTiger.checkpoint.max time (msecs)" && s.max > 60000 {
			issues = append(issues, fmt.Sprintf("Long checkpoint detected: max %dms (%.1fs)", s.max, float64(s.max)/1000))
		}
	}

	// Aggressive eviction mode (WT-13090). The metric is a boolean 0/1 in
	// modern mongod (1 when WT eviction is in "aggressive" mode). Pre-fix
	// checked `s.max > 100` which never fires for a 0/1 gauge — silent
	// dead detector. Trigger on `s.max == 1` (the mode was active at
	// least once during capture) for the boolean variant; keep the >100
	// check as fallback for older numeric variants.
	for _, s := range stats {
		if s.name == "serverStatus.wiredTiger.cache.eviction currently operating in aggressive mode" && (s.max == 1 || s.max > 100) {
			issues = append(issues, fmt.Sprintf("Aggressive eviction mode anomaly (max=%s) — possible WT-13090", formatInt(s.max)))
		}
	}

	// High connections
	maxConn := int64(0)
	for _, s := range samples {
		if v, ok := s.Metrics["serverStatus.connections.current"]; ok && v > maxConn {
			maxConn = v
		}
	}
	if maxConn > 5000 {
		issues = append(issues, fmt.Sprintf("High connection count: peak %s", formatInt(maxConn)))
	}

	// Eviction worker failure rate
	var evictAttempts, evictFailures int64
	for _, s := range stats {
		if s.name == "serverStatus.wiredTiger.cache.eviction worker thread evicting pages" {
			evictAttempts = s.totalDelta
		}
		if s.name == "serverStatus.wiredTiger.thread-yield.page acquire eviction blocked" {
			evictFailures = s.totalDelta
		}
	}
	if evictAttempts > 0 && evictFailures > 0 {
		failRate := float64(evictFailures) / float64(evictAttempts) * 100
		if failRate > 30 {
			issues = append(issues, fmt.Sprintf("High eviction failure rate: %.1f%% (%s/%s)", failRate, formatInt(evictFailures), formatInt(evictAttempts)))
		}
	}

	// FTDC gaps
	gapCount := 0
	for i := 1; i < len(samples); i++ {
		d := samples[i].Timestamp.Sub(samples[i-1].Timestamp)
		if d > 5*time.Second {
			gapCount++
		}
	}
	if gapCount > 0 {
		issues = append(issues, fmt.Sprintf("%d FTDC gaps detected (>5s between samples) — possible stalls", gapCount))
	}

	return issues
}

func printCSV(w io.Writer, samples []Sample) {
	if len(samples) == 0 {
		return
	}

	keySet := make(map[string]bool)
	for _, s := range samples {
		for k := range s.Metrics {
			keySet[k] = true
		}
	}
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Use encoding/csv to handle metric names with commas, quotes, or
	// newlines correctly. WT metric paths often contain spaces, parens,
	// and slashes (e.g., "wiredTiger.cache.bytes (currently in the cache)").
	// Pre-fix, raw fmt.Fprintf produced malformed CSV that some downstream
	// consumers (Python pandas, Excel) silently mis-parsed.
	cw := csv.NewWriter(w)
	defer cw.Flush()
	header := make([]string, 0, len(keys)+1)
	header = append(header, "timestamp")
	header = append(header, keys...)
	_ = cw.Write(header)
	row := make([]string, len(keys)+1)
	for _, s := range samples {
		row[0] = s.Timestamp.UTC().Format(time.RFC3339)
		for i, k := range keys {
			row[i+1] = strconv.FormatInt(s.Metrics[k], 10)
		}
		_ = cw.Write(row)
	}
}

func printMetricList(w io.Writer, samples []Sample) {
	if len(samples) == 0 {
		return
	}
	stats := computeMetricStats(samples)
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].name < stats[j].name
	})
	for _, s := range stats {
		marker := " "
		if s.changed {
			marker = "*"
		}
		fmt.Fprintf(w, "%s %s  [%s .. %s]\n", marker, s.name, formatInt(s.min), formatInt(s.max))
	}
}

func bsonToMap(doc bson.D) map[string]interface{} {
	m := make(map[string]interface{})
	for _, elem := range doc {
		switch v := elem.Value.(type) {
		case bson.D:
			m[elem.Key] = bsonToMap(v)
		case bson.A:
			arr := make([]interface{}, len(v))
			for i, av := range v {
				if subdoc, ok := av.(bson.D); ok {
					arr[i] = bsonToMap(subdoc)
				} else {
					arr[i] = av
				}
			}
			m[elem.Key] = arr
		case primitive.DateTime:
			m[elem.Key] = time.Time(v.Time()).UTC().Format(time.RFC3339)
		default:
			m[elem.Key] = v
		}
	}
	return m
}

func printMetadata(w io.Writer, metadata []MetadataDoc) {
	if len(metadata) == 0 {
		fmt.Fprintln(w, "No metadata found")
		return
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	for _, md := range metadata {
		out := map[string]interface{}{
			"_ftdc_type": md.DocType,
		}
		m := bsonToMap(md.Doc)
		for k, v := range m {
			out[k] = v
		}
		if err := enc.Encode(out); err != nil {
			if isEPIPE(err) {
				return
			}
			fmt.Fprintf(os.Stderr, "ftdc_parser: encode error: %v\n", err)
			return
		}
	}
}

func printStats(w io.Writer, samples []Sample) {
	if len(samples) < 2 {
		fmt.Fprintln(w, "Insufficient samples")
		return
	}

	first, last := samples[0], samples[len(samples)-1]
	durSec := last.Timestamp.Sub(first.Timestamp).Seconds()
	if durSec <= 0 {
		durSec = 1
	}

	stats := computeMetricStats(samples)
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].name < stats[j].name
	})

	type statEntry struct {
		Name  string  `json:"name"`
		Kind  string  `json:"kind"`
		Min   int64   `json:"min"`
		Max   int64   `json:"max"`
		Avg   float64 `json:"avg"`
		P50   int64   `json:"p50"`
		P99   int64   `json:"p99"`
		First int64   `json:"first"`
		Last  int64   `json:"last"`
		Delta int64   `json:"delta"`
		Rate  float64 `json:"rate_ps"`
	}

	enc := json.NewEncoder(w)

	// Print header with overview
	overview := map[string]interface{}{
		"_type":    "overview",
		"samples":  len(samples),
		"start":    first.Timestamp.UTC().Format(time.RFC3339),
		"end":      last.Timestamp.UTC().Format(time.RFC3339),
		"duration": int(durSec),
		"metrics":  len(stats),
	}
	if err := enc.Encode(overview); err != nil {
		if isEPIPE(err) {
			return
		}
		fmt.Fprintf(os.Stderr, "ftdc_parser: encode error: %v\n", err)
		return
	}

	// Collect all values per metric for percentile computation
	metricValues := make(map[string][]int64)
	for _, s := range stats {
		if !s.changed {
			continue
		}
		metricValues[s.name] = make([]int64, 0, len(samples))
	}
	for _, sample := range samples {
		for k, v := range sample.Metrics {
			if vals, ok := metricValues[k]; ok {
				metricValues[k] = append(vals, v)
			}
		}
	}

	// Classify and print each metric
	for _, s := range stats {
		if !s.changed {
			continue
		}
		vals := metricValues[s.name]
		sort.Slice(vals, func(a, b int) bool { return vals[a] < vals[b] })

		// Percentiles
		p50 := vals[(len(vals)-1)*50/100]
		p99idx := int(float64(len(vals)-1) * 0.99)
		if p99idx >= len(vals) {
			p99idx = len(vals) - 1
		}
		p99 := vals[p99idx]

		// Classify: "counter" if monotonically non-decreasing, "gauge" otherwise
		kind := "counter"
		for i := 1; i < len(samples); i++ {
			prev, pok := samples[i-1].Metrics[s.name]
			cur, cok := samples[i].Metrics[s.name]
			if pok && cok && cur < prev {
				kind = "gauge"
				break
			}
		}

		avg := s.sum / float64(s.count)
		entry := statEntry{
			Name:  s.name,
			Kind:  kind,
			Min:   s.min,
			Max:   s.max,
			Avg:   math.Round(avg*100) / 100,
			P50:   p50,
			P99:   p99,
			First: s.first,
			Last:  s.last,
			Delta: s.totalDelta,
			Rate:  math.Round(float64(s.totalDelta)/durSec*100) / 100,
		}
		if err := enc.Encode(entry); err != nil {
			if isEPIPE(err) {
				return
			}
			fmt.Fprintf(os.Stderr, "ftdc_parser: encode error: %v\n", err)
			return
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: ftdc_parser <path> [options]

Options:
  --json             Output one JSON object per sample (NDJSON)
  --rates            Output per-second rates instead of raw counters (implies --json)
  --csv              Output CSV format
  --metadata         Output FTDC metadata (server version, config, hostInfo) as JSON
  --stats            Output per-metric statistics as NDJSON (min/max/avg/first/last/delta/rate)
  --list             List all metric names with min/max ranges (* = changed)
  --filter PATTERN   Filter metrics by regex (can be repeated). Patterns with only
                     alnum/_/-/space/// or \. characters use an optimized substring
                     fast path; anything else uses (?i)PATTERN.
  --nonzero          Exclude metrics that are zero across all samples
  --from TIME        Start time filter (RFC3339, e.g. 2026-03-10T05:50:00Z)
  --to TIME          End time filter (RFC3339)
  --jobs N           Number of parallel chunk-decode workers (default: min(NumCPU, 8))

Auto-mode flags (see SKILL.md for full list):
  --auto                       Enable analyzer Findings stream
  --auto-mode {single,multi,auto}
  --auto-top-k N               Top-K movers per direction (default 50)
  --auto-spike-z FLOAT         Robust-z threshold for kind:"spike" (default 5.0)
  --auto-drift-ratio R         Cross-host drift ratio threshold (default 5.0)
  --auto-follow METRIC         Append to follower vector (repeatable)
  --auto-baseline-file PATH    Compare against historical baseline (Tier D)
  --reset-baseline             Wipe the baseline file before this run starts (Tier D)
  --auto-print-schema          Print JSON Schema for Findings and exit
  --auto-check DIR             Golden-file regression harness
  --auto-emit-otel             Wrap each Finding in OTel LogRecord
  --auto-log-file PATH         Include mongod.log in bundle: trimmed copy +
                               snippet around main event in analysis.md
  --baseline FROM..TO          Override baseline window (RFC3339)
  --stress   FROM..TO          Override stress window (RFC3339)

Env vars:
  FTDC_OBSERVABILITY=1         Include wall-clock duration_ms per stage in
                               run_health (default off for determinism)

Examples:
  ftdc_parser /path/to/mongod.0/                              # Smart summary
  ftdc_parser /path/to/mongod.0/ --json --nonzero             # JSON, skip zero-only metrics
  ftdc_parser /path/to/mongod.0/ --rates --filter opcounters  # Per-second op rates
  ftdc_parser /path/to/mongod.0/ --list --filter cache        # List cache metrics
  ftdc_parser /path/to/mongod.0/ --csv --filter "opcounters|connections"
  ftdc_parser /path/to/file       --json --jobs 8             # Force 8 workers
`)
	// Tier A.5 follow-up: when invoked via --help, exit 0 (success).
	// When invoked due to missing args, the caller passes wantHelp=false
	// (via usageErr).
	os.Exit(usageExitCode)
}

// usageExitCode is set to 0 when usage() is called from a --help branch,
// 1 when called from a missing-args branch. Set by the caller before calling
// usage().
var usageExitCode = 1

// resetBaselineOnlyMode returns true if the args contain --reset-baseline
// AND only flag-style args (no FTDC path). Used to dispatch the early-exit
// reset path regardless of `--reset-baseline` flag position.
func resetBaselineOnlyMode(args []string) bool {
	hasReset := false
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--reset-baseline":
			hasReset = true
		case "--auto-baseline-file":
			i++ // skip the value
		default:
			// Any other token (positional path, unrecognized flag) means
			// this is NOT a reset-only invocation; let main() handle it.
			return false
		}
	}
	return hasReset
}

// extractBaselineFileFlag scans args for `--auto-baseline-file VALUE` and
// returns VALUE, or "" if absent.
func extractBaselineFileFlag(args []string) string {
	for i := 0; i+1 < len(args); i++ {
		if args[i] == "--auto-baseline-file" {
			return args[i+1]
		}
	}
	return ""
}

func defaultJobs() int {
	n := runtime.NumCPU()
	if n > 8 {
		n = 8
	}
	if n < 1 {
		n = 1
	}
	return n
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	// Tier A.2: signal handling. SIGPIPE on broken downstream pipe (e.g.
	// `ftdc_parser ... | head -n 50` closes stdout once head exits) becomes
	// a no-op so we exit cleanly. SIGINT/SIGTERM cancel the root context;
	// workers cooperatively shut down.
	signalSetup()

	// P4.2: load user-authored diagnosis patterns from
	// ~/.config/ftdc-parser/patterns.d/ — autonomous (no flag). Soft-fails
	// on per-file errors; logs to stderr. Loaded patterns appended to
	// diagnosisPatterns and tracked in userLoadedPatternNames for the
	// confidence cap in emitDiagnosisCandidates.
	loadUserPatterns()

	// Allow path-less invocations for early-exit flags like
	// --auto-print-schema and --auto-check that don't read FTDC samples.
	path := os.Args[1]
	if path == "--help" || path == "-h" {
		usageExitCode = 0
		usage()
	}
	// Tier D follow-up: when --reset-baseline appears as the SOLE
	// path-position flag (i.e. the user wants to wipe and exit without
	// processing FTDC), wipe the baseline file and exit. Supports both
	// orderings: `--reset-baseline --auto-baseline-file PATH` and
	// `--auto-baseline-file PATH --reset-baseline`. Pre-fix, the second
	// ordering treated the path as an FTDC file and confused the user.
	if resetBaselineOnlyMode(os.Args[1:]) {
		bp := extractBaselineFileFlag(os.Args[1:])
		if bp == "" {
			bp = defaultBaselinePath()
		}
		if err := os.Remove(bp); err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "ftdc_parser: failed to reset baseline at %s: %v\n", bp, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "ftdc_parser: baseline reset at %s\n", bp)
		return
	}
	if path == "--auto-print-schema" {
		bw := bufio.NewWriterSize(os.Stdout, 1<<20)
		defer bw.Flush()
		installInterruptFlusher(bw)
		fmt.Fprintln(bw, autoJSONSchema())
		return
	}
	if path == "--auto-check" {
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: --auto-check requires a directory path")
			os.Exit(1)
		}
		bw := bufio.NewWriterSize(os.Stdout, 1<<20)
		defer bw.Flush()
		installInterruptFlusher(bw)
		if err := runAutoCheck(bw, os.Args[2], defaultJobs()); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	var (
		jsonOutput     bool
		ratesOutput    bool
		csvOutput      bool
		listOutput     bool
		statsOutput    bool
		metadataOutput bool
		nonzero        bool
		filters        []string
		fromTime       time.Time
		toTime         time.Time
		hasFrom        bool
		hasTo          bool
		jobs           int = -1 // -1 = unset → default

		// --auto and companions. All optional; --auto alone enables the
		// analyzer pipeline with sensible defaults echoed in run_metadata.
		autoOutput      bool
		autoCfg         = defaultAutoConfig()
		autoPrintSchema bool
		autoCheckDir    string
		autoEmitOTel    bool
		autoEmitGolden  bool
	)

	for i := 2; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--json":
			jsonOutput = true
		case "--rates":
			ratesOutput = true
		case "--csv":
			csvOutput = true
		case "--list":
			listOutput = true
		case "--stats":
			statsOutput = true
		case "--metadata":
			metadataOutput = true
		case "--nonzero":
			nonzero = true
		case "--filter":
			if i+1 < len(os.Args) {
				i++
				filters = append(filters, os.Args[i])
			}
		case "--from":
			if i+1 < len(os.Args) {
				i++
				t, err := time.Parse(time.RFC3339, os.Args[i])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing --from: %v\n", err)
					os.Exit(1)
				}
				fromTime = t
				hasFrom = true
			}
		case "--jobs":
			if i+1 < len(os.Args) {
				i++
				var n int
				_, err := fmt.Sscanf(os.Args[i], "%d", &n)
				if err != nil || n < 1 {
					fmt.Fprintf(os.Stderr, "Error: --jobs requires a positive integer (got %q)\n", os.Args[i])
					os.Exit(1)
				}
				jobs = n
			}
		case "--to":
			if i+1 < len(os.Args) {
				i++
				t, err := time.Parse(time.RFC3339, os.Args[i])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing --to: %v\n", err)
					os.Exit(1)
				}
				toTime = t
				hasTo = true
			}

		// --auto and companion flags. --auto enables the analyzer
		// pipeline; the others tune it. All flags are recorded in the
		// run_metadata Finding so the LLM downstream sees every choice.
		case "--auto":
			autoOutput = true
		case "--auto-mode":
			if i+1 < len(os.Args) {
				i++
				autoCfg.Mode = os.Args[i]
			}
		case "--auto-top-k":
			if i+1 < len(os.Args) {
				i++
				var n int
				if _, err := fmt.Sscanf(os.Args[i], "%d", &n); err == nil && n > 0 {
					autoCfg.TopK = n
				}
			}
		case "--auto-follow":
			if i+1 < len(os.Args) {
				i++
				autoCfg.FollowerMetrics = append(autoCfg.FollowerMetrics, os.Args[i])
			}
		case "--auto-drift-ratio":
			if i+1 < len(os.Args) {
				i++
				var v float64
				if _, err := fmt.Sscanf(os.Args[i], "%f", &v); err == nil && v > 0 {
					autoCfg.DriftRatio = v
				}
			}
		case "--auto-spike-z":
			if i+1 < len(os.Args) {
				i++
				var v float64
				if _, err := fmt.Sscanf(os.Args[i], "%f", &v); err == nil && v > 0 {
					autoCfg.SpikeZ = v
				}
			}
		case "--baseline":
			if i+1 < len(os.Args) {
				i++
				from, to, ok := parseRange(os.Args[i])
				if !ok {
					fmt.Fprintf(os.Stderr, "Error parsing --baseline %q (want FROM..TO RFC3339)\n", os.Args[i])
					os.Exit(1)
				}
				autoCfg.BaselineFrom = from
				autoCfg.BaselineTo = to
				autoCfg.HasBaseline = true
			}
		case "--stress":
			if i+1 < len(os.Args) {
				i++
				from, to, ok := parseRange(os.Args[i])
				if !ok {
					fmt.Fprintf(os.Stderr, "Error parsing --stress %q (want FROM..TO RFC3339)\n", os.Args[i])
					os.Exit(1)
				}
				autoCfg.StressFrom = from
				autoCfg.StressTo = to
				autoCfg.HasStress = true
			}
		case "--auto-check":
			if i+1 < len(os.Args) {
				i++
				autoCheckDir = os.Args[i]
			}
		case "--auto-emit-golden":
			autoEmitGolden = true
			autoCfg.EmitGolden = true
		case "--auto-print-schema":
			autoPrintSchema = true
		case "--auto-emit-otel":
			autoEmitOTel = true
			autoCfg.EmitOTel = true
		case "--auto-baseline-file":
			if i+1 < len(os.Args) {
				i++
				autoCfg.BaselineFile = os.Args[i]
			}
		case "--auto-compare":
			// P2.3: explicit two-capture comparison. The path is the
			// "secondary" capture; the positional path remains the
			// primary. capture_diff Findings carry direction language
			// (primary_higher / secondary_higher) that maps mechanically
			// to which side is which.
			if i+1 < len(os.Args) {
				i++
				autoCfg.ComparePath = os.Args[i]
			}
		case "--auto-recommend-tunables":
			// P1.NEW-A: opt-in tunable_observation stage. Tier-7
			// narration scaffold; bias-allowlisted in BIAS_CONTROL.md.
			autoCfg.RecommendTunables = true
		case "--multi-host-parallel":
			// P1.6: bounded per-host parallelism in multi-host mode.
			// Default 1 (sequential) — see memory-constraint note in
			// analyze.go's AutoConfig comment. Users who explicitly opt
			// into N>1 are acknowledging the OOM risk.
			if i+1 < len(os.Args) {
				i++
				if n, err := strconv.Atoi(os.Args[i]); err == nil && n > 0 {
					autoCfg.MultiHostParallel = n
				}
			}
		case "--reset-baseline":
			// Tier D: wipe the baseline file before this run starts.
			// Useful after a noisy first capture poisons the EWMA seed.
			if autoCfg.BaselineFile == "" {
				autoCfg.BaselineFile = defaultBaselinePath()
			}
			if err := os.Remove(autoCfg.BaselineFile); err != nil && !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "ftdc_parser: failed to reset baseline at %s: %v\n", autoCfg.BaselineFile, err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "ftdc_parser: baseline reset at %s\n", autoCfg.BaselineFile)
		case "--auto-around-event":
			if i+1 < len(os.Args) {
				i++
				t, err := time.Parse(time.RFC3339, os.Args[i])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing --auto-around-event: %v\n", err)
					os.Exit(1)
				}
				autoCfg.AroundEvent = t
				autoCfg.HasAroundEvent = true
			}
		case "--auto-window":
			if i+1 < len(os.Args) {
				i++
				d, err := time.ParseDuration(os.Args[i])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing --auto-window: %v\n", err)
					os.Exit(1)
				}
				autoCfg.AroundWindow = d
			}
		case "--auto-segment-by":
			if i+1 < len(os.Args) {
				i++
				autoCfg.SegmentByMetric = os.Args[i]
			}
		case "--auto-correlate":
			if i+2 < len(os.Args) {
				autoCfg.CorrelateLeader = os.Args[i+1]
				autoCfg.CorrelateFollower = os.Args[i+2]
				i += 2
			}
		case "--auto-lag-resolution":
			if i+1 < len(os.Args) {
				i++
				var v int
				if _, err := fmt.Sscanf(os.Args[i], "%d", &v); err == nil && v > 0 {
					autoCfg.LagResolutionMS = v
				}
			}
		case "--auto-budget-tokens":
			if i+1 < len(os.Args) {
				i++
				var v int
				if _, err := fmt.Sscanf(os.Args[i], "%d", &v); err == nil && v >= 0 {
					autoCfg.BudgetTokens = v
				}
			}
		case "--auto-emit-score-table":
			autoCfg.EmitScoreTable = true
		case "--auto-emit-t2":
			if i+1 < len(os.Args) {
				i++
				autoCfg.AutoEmitT2Path = os.Args[i]
			}
		case "--auto-emit-t2-no-tarball":
			autoCfg.AutoEmitT2NoTarball = true
		case "--auto-emit-t2-top-k":
			if i+1 < len(os.Args) {
				i++
				if n, err := strconv.Atoi(os.Args[i]); err == nil && n > 0 {
					autoCfg.AutoEmitT2TopK = n
				}
			}
		case "--auto-emit-t2-link-ftdc":
			autoCfg.AutoEmitT2LinkFtdc = true
		case "--auto-emit-t2-zoom":
			autoCfg.AutoEmitT2Zoom = true
		case "--auto-log-file":
			if i+1 < len(os.Args) {
				i++
				autoCfg.AutoLogFile = os.Args[i]
			}
		default:
			fmt.Fprintf(os.Stderr, "Unknown option: %s\n", os.Args[i])
			usage()
		}
	}

	matcher, err := compileMatcher(filters)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if jobs <= 0 {
		jobs = defaultJobs()
	}

	// Buffered stdout: cuts write syscalls and accelerates --json/--csv/--rates.
	bw := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer bw.Flush()
	installInterruptFlusher(bw)

	// Early-exit modes (no sample collection required).
	if autoPrintSchema {
		fmt.Fprintln(bw, autoJSONSchema())
		return
	}
	if autoCheckDir != "" {
		if err := runAutoCheck(bw, autoCheckDir, jobs); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	// --auto-emit-t2: wrap the single output writer in a findingTee so every
	// NDJSON line is captured for post-run bundle assembly. The tee is a
	// transparent pass-through; stdout content is bit-identical with or
	// without the flag. Registered as a defer here (after early exits) so it
	// runs only on auto-analysis paths. LIFO order: emitT2Bundle fires before
	// bw.Flush() — that's correct, since the tee holds its own copy of all
	// captured lines and doesn't require stdout to be flushed first.
	var autoWriter io.Writer = bw
	if autoCfg.AutoEmitT2Path != "" {
		tee := newFindingTee(bw)
		autoWriter = tee
		// Close over autoCfg by reference so the deferred call sees the
		// final cfg state (including InputPath set below in the P0.3 block).
		defer func() { emitT2Bundle(tee, autoCfg) }()
	}

	// P0.3: capture input provenance into autoCfg BEFORE any dispatch
	// branch — the streaming single-host path, the multi-host streaming
	// path, and the legacy aggregate-then-analyze path (used by ReAct
	// modes) all share this cfg. Hash failures are non-fatal — the
	// emitter omits the fields when InputFingerprint is nil.
	if autoOutput {
		autoCfg.InputPath = path
		if fp, err := computeInputFingerprint(path); err == nil {
			autoCfg.InputFingerprint = fp
		}
		autoCfg.CmdLine = normalizedCmdLine(os.Args)
	}

	// Multi-host detection MUST run BEFORE the single-host parse — when
	// the input is a directory-of-host-directories, the top-level parse
	// produces zero samples (listFTDCDir skips subdirectories) and the
	// "No samples found" branch below would early-exit before the
	// multi-host dispatch could even consider the input.
	//
	// Memory note (fix for SERVER-XXXXX): the previous loop accumulated
	// every host's full []Sample slice into a samplesByHost map BEFORE
	// invoking runAutoMultiNamed. On 10 hosts × 48h that's ~150 GB
	// resident. The new loop streams each host's data into a per-host
	// columnar seriesStore, runs the per-host pipeline, captures only the
	// per-host aggregates (nameToID, fullMean, fullSD, observedMin/Max,
	// isCounter, isDegenerate) needed for cross-host comparison, and
	// releases the full store before the next host starts.
	if autoOutput && autoCfg.Mode != "single" {
		hosts := detectMultiHostInputDetailed(path)
		if len(hosts) >= 2 {
			if err := runAutoMultiStreaming(autoWriter, hosts, matcher, jobs, hasFrom, hasTo, fromTime, toTime, autoCfg); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
	}

	// Single-host --auto: stream chunks into a columnar seriesStore so the
	// row-oriented []Sample never accumulates globally. Bypasses the
	// 16-25 GB peak that the original aggregate-then-analyze pipeline
	// produced on 48h captures. React modes (--auto-around-event,
	// --auto-segment-by, --auto-correlate) fall through to the legacy
	// aggregate-then-filter path below — their working windows are small.
	if autoOutput && !autoCfg.HasAroundEvent && autoCfg.SegmentByMetric == "" &&
		!(autoCfg.CorrelateLeader != "" && autoCfg.CorrelateFollower != "") {
		store, meta, serr := buildSeriesStoreStreaming(path, matcher, jobs, hasFrom, hasTo, fromTime, toTime)
		if serr != nil {
			fmt.Fprintln(os.Stderr, serr)
			os.Exit(1)
		}
		if store.nSamples() == 0 {
			if hasFrom || hasTo {
				fmt.Fprintln(bw, "No samples in specified time range")
			} else {
				fmt.Fprintln(bw, "No samples found")
			}
			return
		}
		host := hostIDFromPath(path)
		if err := runAutoWithHostStore(autoWriter, host, store, meta, autoCfg); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// P2.3 --auto-compare: after the primary --auto stream, run a
		// second pass on autoCfg.ComparePath and emit comparative Findings.
		// We reuse the same writer (bw) and emit through a fresh
		// findingCollector so capture_diff outputs slot AFTER the primary
		// run_health line. Memory peak is bounded — primary store's
		// columns are still alive when secondary builds, but per-host
		// store deallocation is reasonable on a typical host.
		if autoCfg.ComparePath != "" {
			cc := &findingCollector{}
			if err := runCaptureCompare(cc, host, store, autoCfg.ComparePath, jobs, hasFrom, hasTo, fromTime, toTime); err != nil {
				fmt.Fprintf(os.Stderr, "ftdc_parser: --auto-compare failed: %v\n", err)
			} else if err := cc.emitAll(autoWriter); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
		return
	}

	// Aggregate all samples in memory; this preserves the existing output
	// contract for every mode and lets --from/--to and --nonzero post-process.
	var samples []Sample
	metadata := processFTDC(path, matcher, jobs, func(batch []Sample) {
		samples = append(samples, batch...)
	})

	// Match the pre-rewrite binary's check ordering exactly:
	//   1. "No samples found" when the file had zero chunks (before time filter).
	//   2. Time-range filter is applied second.
	//   3. "No samples in specified time range" only when the range emptied a
	//      previously-non-empty sample set.
	//   4. Output-mode dispatch (including --metadata) only runs when samples > 0.
	if len(samples) == 0 {
		fmt.Fprintln(bw, "No samples found")
		return
	}

	// Time range filter
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
		if len(samples) == 0 {
			fmt.Fprintln(bw, "No samples in specified time range")
			return
		}
	}

	if nonzero {
		samples = stripZeroMetrics(samples)
	}

	// --auto takes precedence over the existing output modes. It is
	// intentionally additive: existing modes (--stats, --json, --csv,
	// --list, --metadata, --rates, default summary) are bit-for-bit
	// identical to the pre-auto binary. Multi-host dispatch is handled
	// upstream of the single-host parse (see above) to avoid the
	// "No samples found" early-exit on directory-of-hosts inputs.
	if autoOutput {

		var rerr error
		switch {
		case autoCfg.HasAroundEvent:
			rerr = runAutoAroundEvent(autoWriter, path, samples, metadata, autoCfg)
		case autoCfg.SegmentByMetric != "":
			rerr = runAutoSegmentBy(autoWriter, path, samples, metadata, autoCfg)
		case autoCfg.CorrelateLeader != "" && autoCfg.CorrelateFollower != "":
			rerr = runAutoCorrelate(autoWriter, path, samples, metadata, autoCfg)
		default:
			rerr = runAuto(autoWriter, path, samples, metadata, autoCfg)
		}
		if rerr != nil {
			fmt.Fprintln(os.Stderr, rerr)
			os.Exit(1)
		}
		_ = autoEmitOTel
		_ = autoEmitGolden
		return
	}

	switch {
	case metadataOutput:
		printMetadata(bw, metadata)
		return
	case statsOutput:
		printStats(bw, samples)
	case listOutput:
		printMetricList(bw, samples)
	case csvOutput:
		printCSV(bw, samples)
	case ratesOutput:
		rates := computeRates(samples)
		enc := json.NewEncoder(bw)
		for _, r := range rates {
			if err := enc.Encode(r); err != nil {
				if isEPIPE(err) {
					return
				}
				fmt.Fprintf(os.Stderr, "ftdc_parser: encode error: %v\n", err)
				return
			}
		}
	case jsonOutput:
		enc := json.NewEncoder(bw)
		for _, s := range samples {
			if err := enc.Encode(s); err != nil {
				if isEPIPE(err) {
					return
				}
				fmt.Fprintf(os.Stderr, "ftdc_parser: encode error: %v\n", err)
				return
			}
		}
	default:
		printSummary(bw, samples, metadata)
	}
}
