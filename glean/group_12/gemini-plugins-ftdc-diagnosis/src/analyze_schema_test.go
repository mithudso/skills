package main

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestAutoJSONSchemaIsValidJSON verifies that autoJSONSchema() returns
// valid JSON that can be unmarshaled without error.
func TestAutoJSONSchemaIsValidJSON(t *testing.T) {
	schema := autoJSONSchema()
	if schema == "" {
		t.Fatal("autoJSONSchema() returned empty string")
	}

	var m map[string]interface{}
	err := json.Unmarshal([]byte(schema), &m)
	if err != nil {
		t.Fatalf("autoJSONSchema() is not valid JSON: %v", err)
	}

	// Verify the schema has the expected top-level structure.
	if _, ok := m["$schema"]; !ok {
		t.Fatal("schema missing $schema field")
	}
	if _, ok := m["properties"]; !ok {
		t.Fatal("schema missing properties field")
	}
}

// TestSchemaEnumOneOfBijection verifies that:
//  1. Every kind in the enum has a corresponding oneOf entry.
//  2. Every kind in oneOf is also in the enum (no ghost entries).
//
// This is a structural test on autoFindingsJSONSchema; it parses the schema
// JSON and checks the two arrays directly.
func TestSchemaEnumOneOfBijection(t *testing.T) {
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(autoFindingsJSONSchema), &schema); err != nil {
		t.Fatalf("schema is not valid JSON: %v", err)
	}

	// Extract enum values from properties.kind.enum.
	props, _ := schema["properties"].(map[string]interface{})
	kindProp, _ := props["kind"].(map[string]interface{})
	enumRaw, _ := kindProp["enum"].([]interface{})
	enumSet := make(map[string]bool, len(enumRaw))
	for _, v := range enumRaw {
		if s, ok := v.(string); ok {
			enumSet[s] = true
		}
	}

	// Extract kind values from oneOf entries.
	oneOfRaw, _ := schema["oneOf"].([]interface{})
	oneOfSet := make(map[string]bool, len(oneOfRaw))
	for _, entry := range oneOfRaw {
		m, _ := entry.(map[string]interface{})
		entryProps, _ := m["properties"].(map[string]interface{})
		kindConst, _ := entryProps["kind"].(map[string]interface{})
		if cv, ok := kindConst["const"].(string); ok {
			oneOfSet[cv] = true
		}
	}

	// Every enum kind must have a oneOf entry.
	for kind := range enumSet {
		if !oneOfSet[kind] {
			t.Errorf("kind %q is in enum but has no oneOf entry", kind)
		}
	}
	// Every oneOf kind must be in the enum.
	for kind := range oneOfSet {
		if !enumSet[kind] {
			t.Errorf("kind %q is in oneOf but missing from enum", kind)
		}
	}
}

// TestEveryKindConstHasSchemaEntry verifies that every Kind* constant
// declared in analyze.go appears in the schema's kind enum.
func TestEveryKindConstHasSchemaEntry(t *testing.T) {
	schema := autoJSONSchema()

	// List of all Kind* const strings from analyze.go.
	kindConsts := []string{
		KindRunMetadata,
		KindSchemaFingerprint,
		KindMetricScoreTable,
		KindMover,
		KindCorrelation,
		KindMonotonicTrend,
		KindTrendInconclusive,
		KindChangepoint,
		KindChangepointOversegmented,
		KindHistogramPercentile,
		KindHistogramInvalid,
		KindGap,
		KindVarianceShift,
		KindSkipped,
		KindHostOnlyMetrics,
		KindCrossHostDrift,
		KindCrossHostRoleInferred,
		KindHostMetricSetAsymmetry,
		KindRunHealth,
		KindSpike,
	}

	for _, kind := range kindConsts {
		// Check if the kind appears as a literal string in the schema.
		searchStr := `"` + kind + `"`
		if !strings.Contains(schema, searchStr) {
			t.Errorf("Kind constant %q not found in autoJSONSchema()", kind)
		}
	}
}
