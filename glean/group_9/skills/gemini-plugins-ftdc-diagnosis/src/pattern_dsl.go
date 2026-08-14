// pattern_dsl.go — P4.2 user-editable pattern DSL.
//
// Drop one or more JSON files into ~/.config/ftdc-parser/patterns.d/
// and the parser will merge them into the built-in `diagnosisPatterns`
// catalog at startup. Each file is a single pattern. Schema:
//
//   {
//     "name": "my_custom_shape",
//     "required_kinds":  ["seal_event", "role_asymmetry"],
//     "supporting_kinds": ["term_change"],
//     "excluded_kinds":   ["counter_inversion"],
//     "discriminating_kinds_needed": ["counter_wrap"],
//     "domain_stage_gate": "pali",
//     "rationale_template": "Custom shape — %d supporting; %d excluded."
//   }
//
// Validation rules (lint-style, applied at load):
//   - `name` must not collide with a built-in pattern. Collisions skip
//     the user pattern + log to stderr (mongod-style soft failure).
//   - All kinds referenced must appear in the parser's Finding-kind
//     enum (validated against the loaded analyze_schema.go literals).
//   - `domain_stage_gate`, when non-empty, must be one of the known
//     stages: "", "pali", "replication", "sharded", "host_os".
//
// Per the plan's autonomy contract: pattern files are picked up
// automatically — no flag, no restart beyond the next run. User-
// authored patterns are tagged with `pattern_source:"user_yaml"` (kept
// the field name even though we use JSON; it's the conceptual
// distinction) and the analyzer narration caps their confidence_rank
// at 0.85 until validated against fixtures.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// userPatternsDir returns where the user's pattern files live.
func userPatternsDir() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		home = "/tmp"
	}
	return filepath.Join(home, ".config", "ftdc-parser", "patterns.d")
}

// userPatternConfidenceCap is the maximum confidence_rank a user
// pattern can receive. The parser's diagnosis_candidate emitter
// should respect this when emitting candidates whose `name` matches a
// user-loaded pattern.
const userPatternConfidenceCap = 0.85

// userLoadedPatternNames is populated by loadUserPatterns at startup
// and read by emitDiagnosisCandidates to apply the confidence cap.
// Package-level mutable state is unusual in this codebase; we tolerate
// it because pattern loading is one-shot at process start.
var userLoadedPatternNames []string

// loadUserPatterns reads every *.json file in the patterns dir and
// appends valid entries to diagnosisPatterns. Returns the slice of
// user-loaded pattern names so emitDiagnosisCandidates can apply the
// confidence cap.
//
// Soft-fails on individual file errors (logs to stderr, continues).
func loadUserPatterns() []string {
	dir := userPatternsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	knownKinds := knownKindSet()
	knownStages := map[string]bool{"": true, "pali": true, "replication": true, "sharded": true, "host_os": true}
	builtinNames := map[string]bool{}
	for _, p := range diagnosisPatterns {
		builtinNames[p.name] = true
	}
	var loadedNames []string
	// Sort filenames so load order is deterministic.
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)
	for _, n := range names {
		full := filepath.Join(dir, n)
		data, err := os.ReadFile(full)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ftdc_parser: user pattern read error %s: %v\n", full, err)
			continue
		}
		var raw struct {
			Name                      string   `json:"name"`
			RequiredKinds             []string `json:"required_kinds"`
			SupportingKinds           []string `json:"supporting_kinds"`
			ExcludedKinds             []string `json:"excluded_kinds"`
			DiscriminatingKindsNeeded []string `json:"discriminating_kinds_needed"`
			DomainStageGate           string   `json:"domain_stage_gate"`
			RationaleTemplate         string   `json:"rationale_template"`
		}
		if err := json.Unmarshal(data, &raw); err != nil {
			fmt.Fprintf(os.Stderr, "ftdc_parser: user pattern parse error %s: %v\n", full, err)
			continue
		}
		if raw.Name == "" {
			fmt.Fprintf(os.Stderr, "ftdc_parser: user pattern %s missing `name`; skipped\n", full)
			continue
		}
		if builtinNames[raw.Name] {
			fmt.Fprintf(os.Stderr, "ftdc_parser: user pattern %s name %q collides with a built-in pattern; skipped\n", full, raw.Name)
			continue
		}
		if len(raw.RequiredKinds) == 0 {
			fmt.Fprintf(os.Stderr, "ftdc_parser: user pattern %s: required_kinds must not be empty; skipped\n", full)
			continue
		}
		if !knownStages[raw.DomainStageGate] {
			fmt.Fprintf(os.Stderr, "ftdc_parser: user pattern %s domain_stage_gate=%q not recognized; skipped\n", full, raw.DomainStageGate)
			continue
		}
		// Validate every referenced kind.
		validateErr := func(label string, kinds []string) error {
			for _, k := range kinds {
				if !knownKinds[k] {
					return fmt.Errorf("%s references unknown Finding kind %q", label, k)
				}
			}
			return nil
		}
		if err := validateErr("required_kinds", raw.RequiredKinds); err != nil {
			fmt.Fprintf(os.Stderr, "ftdc_parser: user pattern %s: %v; skipped\n", full, err)
			continue
		}
		if err := validateErr("supporting_kinds", raw.SupportingKinds); err != nil {
			fmt.Fprintf(os.Stderr, "ftdc_parser: user pattern %s: %v; skipped\n", full, err)
			continue
		}
		if err := validateErr("excluded_kinds", raw.ExcludedKinds); err != nil {
			fmt.Fprintf(os.Stderr, "ftdc_parser: user pattern %s: %v; skipped\n", full, err)
			continue
		}
		if err := validateErr("discriminating_kinds_needed", raw.DiscriminatingKindsNeeded); err != nil {
			fmt.Fprintf(os.Stderr, "ftdc_parser: user pattern %s: %v; skipped\n", full, err)
			continue
		}
		template := raw.RationaleTemplate
		if template == "" {
			template = raw.Name + " shape: %d supporting; %d excluded."
		}
		diagnosisPatterns = append(diagnosisPatterns, diagnosisPattern{
			name:                      raw.Name,
			requiredKinds:             raw.RequiredKinds,
			supportingKinds:           raw.SupportingKinds,
			excludedKinds:             raw.ExcludedKinds,
			discriminatingKindsNeeded: raw.DiscriminatingKindsNeeded,
			domainStageGate:           raw.DomainStageGate,
			rationaleTemplate:         template,
		})
		// Mark the name as seen so a second user-pattern file with the
		// same name is caught by the collision check above.
		builtinNames[raw.Name] = true
		loadedNames = append(loadedNames, raw.Name)
	}
	userLoadedPatternNames = loadedNames
	return loadedNames
}

// knownKindSet returns the set of Finding-kind strings the parser will
// accept in a user pattern. Pulled from the JSON Schema enum so the
// source of truth is one place.
func knownKindSet() map[string]bool {
	out := map[string]bool{}
	// Extract the enum from analyze_schema.go's literal. The schema is
	// a single Go string; we scan it for `"some_kind"` entries inside
	// the "enum" block. This avoids hardcoding the list in two places.
	schemaText := autoFindingsJSONSchema
	// Locate the enum array start.
	idx := strings.Index(schemaText, `"enum"`)
	if idx < 0 {
		return out
	}
	tail := schemaText[idx:]
	// Find the closing bracket of this enum.
	end := strings.Index(tail, "]")
	if end < 0 {
		return out
	}
	enum := tail[:end]
	// Pull quoted strings.
	for {
		q1 := strings.Index(enum, `"`)
		if q1 < 0 {
			break
		}
		enum = enum[q1+1:]
		q2 := strings.Index(enum, `"`)
		if q2 < 0 {
			break
		}
		v := enum[:q2]
		// Skip the "enum" string itself.
		if v != "enum" {
			out[v] = true
		}
		enum = enum[q2+1:]
	}
	return out
}

// userPatternConfidenceCapApplies returns true when the given pattern
// name was loaded from a user file (used by emitDiagnosisCandidates to
// cap confidence at userPatternConfidenceCap).
func userPatternConfidenceCapApplies(patternName string, userLoaded []string) bool {
	for _, n := range userLoaded {
		if n == patternName {
			return true
		}
	}
	return false
}
