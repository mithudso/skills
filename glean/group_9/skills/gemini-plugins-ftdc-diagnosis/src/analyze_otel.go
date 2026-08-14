// analyze_otel.go - OpenTelemetry log-record envelope for --auto Findings.
//
// When --auto-emit-otel is set, each Finding is wrapped in an OTel
// LogRecord-shaped JSON envelope:
//
//	{
//	  "event": {
//	    "name": "ftdc.finding.<kind>",
//	    "observed_time_unix_nano": <int>,
//	    "severity_number": 9,                  // INFO
//	    "severity_text": "INFO",
//	    "body": <raw Finding JSON>,
//	    "attributes": {
//	      "ftdc.host": "<host>",
//	      "ftdc.kind": "<kind>",
//	      "ftdc.finding_id": "f-..."
//	    }
//	  }
//	}
//
// This is a pragmatic shape compatible with OTel log pipelines (Collector
// receivers, OTLP) without depending on the actual OTel Go libraries —
// we emit JSON directly. Downstream consumers wanting strict OTLP can pipe
// through their preferred adapter.
package main

import (
	"encoding/json"
	"sort"
	"time"
)

// otelWrap converts a Finding to an OTel LogRecord-shaped map. The
// resulting object marshals to NDJSON via the standard json package.
// observed_time_unix_nano comes from the Finding's `at` field when present,
// else the current wall clock.
func otelWrap(f Finding) map[string]interface{} {
	kind, _ := f["kind"].(string)
	host, _ := f["host"].(string)
	id, _ := f["id"].(string)
	var tsNanos int64
	if atStr, ok := f["at"].(string); ok && atStr != "" {
		if t, err := time.Parse(tsFormat, atStr); err == nil {
			tsNanos = t.UnixNano()
		}
	}
	if tsNanos == 0 {
		// Tier A.4 determinism: leaking wall clock breaks byte-equality
		// across runs. Use 0 by default; only fall back to time.Now()
		// when FTDC_OBSERVABILITY=1 explicitly requests it.
		// observabilityEnabled is cached at init from the env var.
		if observabilityEnabled {
			tsNanos = time.Now().UnixNano()
		}
	}
	return map[string]interface{}{
		"event": map[string]interface{}{
			"name":                    "ftdc.finding." + kind,
			"observed_time_unix_nano": tsNanos,
			"severity_number":         9, // OTLP SeverityNumber.INFO
			"severity_text":           "INFO",
			"body":                    map[string]interface{}(f),
			"attributes": map[string]interface{}{
				"ftdc.host":       host,
				"ftdc.kind":       kind,
				"ftdc.finding_id": id,
			},
		},
	}
}

// otelEmitAll writes a collector's Findings as OTel log records (NDJSON).
// Used when --auto-emit-otel is set. Bypasses the standard findingCollector
// emit path; the same deterministic ordering applies.
func (c *findingCollector) otelEmitAll(w jsonWriter) error {
	// Sort with the same key the standard path uses, to keep output
	// deterministic between -otel and not.
	c.sortDeterministic()
	for _, f := range c.items {
		line, err := json.Marshal(otelWrap(f))
		if err != nil {
			return err
		}
		if _, err := w.Write(line); err != nil {
			return err
		}
		if _, err := w.Write([]byte{'\n'}); err != nil {
			return err
		}
	}
	return nil
}

// sortDeterministic is a helper that performs the same sort emitAll does.
// Extracted so otelEmitAll can call it without duplicating the comparator.
func (c *findingCollector) sortDeterministic() {
	// Re-uses the same comparator as emitAll. The current emitAll sorts
	// in place too, so calling this multiple times is idempotent.
	items := c.items
	// Insertion sort here would be cheap for typical N (~hundreds of
	// Findings), but reuse the stdlib sort for safety.
	stableSortFindings(items)
}

// stableSortFindings sorts in-place by (kind, stableSortKey).
func stableSortFindings(items []Finding) {
	// Forward to the same comparator emitAll uses.
	stableSortFindingsImpl(items)
}

func stableSortFindingsImpl(items []Finding) {
	type pair struct {
		idx int
		key string
	}
	pairs := make([]pair, len(items))
	for i, f := range items {
		kind, _ := f["kind"].(string)
		pairs[i] = pair{i, kind + "\x00" + stableSortKey(f)}
	}
	// Tier A.5 fix: insertion sort was O(N²) on finding counts that can
	// reach thousands; the comment claiming "N is small" no longer holds
	// for full --auto runs. sort.SliceStable is the right tool.
	sort.SliceStable(pairs, func(i, j int) bool {
		return pairs[i].key < pairs[j].key
	})
	sorted := make([]Finding, len(items))
	for i, p := range pairs {
		sorted[i] = items[p.idx]
	}
	copy(items, sorted)
}

// jsonWriter is an io.Writer alias to avoid an import in this file.
type jsonWriter interface {
	Write([]byte) (int, error)
}
