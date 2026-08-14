// analyze_wave9_fixes_test.go — regression tests for Wave 9 confirmed bugs.
//
// W9-1: analyze_recap.go bh_significant_count type mismatch (int vs float64).
// W9-2: analyze_role_asymmetry.go direction label based on alphabetical order
//       instead of actual primary/standby role.
package main

import (
	"testing"
)

// makeRoleStore builds a minimal seriesStore for a single PALI metric,
// used by the role_asymmetry direction tests.
func makeRoleStore(mean float64) *seriesStore {
	return &seriesStore{
		names:        []string{"serverStatus.metrics.disagg.logServer.appendLogTotal"},
		nameToID:     map[string]int{"serverStatus.metrics.disagg.logServer.appendLogTotal": 0},
		isCounter:    []bool{true},
		isDegenerate: []bool{false},
		fullMean:     []float64{mean},
		fullSD:       []float64{50.0},
		fullN:        []int{100},
	}
}

// Fix W9-1: bh_significant_count is stored as int in the metric_score_table
// Finding but was being asserted as float64, which always fails. The recap
// stage's SignificantCount was always 0.
func TestRecapBHSignificantCountType(t *testing.T) {
	c := &findingCollector{}
	c.add(newFinding("metric_score_table", "h1", "at", map[string]interface{}{
		"bh_significant_count": 42, // int, as returned by countSignificantBH
		"fdr_alpha":            0.05,
	}))
	c.add(newFinding(KindMover, "h1", "at", map[string]interface{}{
		"metric":    "serverStatus.opcounters.insert",
		"direction": "up",
	}))

	// Call the production function — not a reimplementation.
	entry := buildCurrentRecap(c, "h1", "test-v1", "abc12345", "sha256test")
	if entry.SignificantCount != 42 {
		t.Errorf("SignificantCount = %d, want 42 (pre-fix float64 assertion yielded 0)", entry.SignificantCount)
	}
}

// Fix W9-2: role_asymmetry direction must reflect actual primary/standby role,
// not alphabetical host order. When the standby host sorts before the primary
// alphabetically, "primary_higher" was being emitted when the standby was the
// elevated host.
func TestRoleAsymmetryDirectionByRole(t *testing.T) {
	// Scenario: standby="aardvark-host", primary="zebra-host"
	// Alphabetically: aardvark < zebra, so hA="aardvark" (standby), hB="zebra" (primary).
	// Metric has meanA=200, meanB=100 → d > 0 → standby is higher.
	// Pre-fix: direction = "primary_higher" (wrong — labeled by alphabetical position).
	// Post-fix: direction = "standby_higher" (correct).

	standbyHost := "aardvark-host"
	primaryHost := "zebra-host"

	hostStores := map[string]*seriesStore{
		standbyHost: makeRoleStore(200.0), // standby has higher mean
		primaryHost: makeRoleStore(100.0), // primary has lower mean
	}
	hostRoles := map[string]PALIRole{
		standbyHost: PALIRoleStandby,
		primaryHost: PALIRolePrimary,
	}

	c := &findingCollector{}
	emitted, _ := emitRoleAsymmetry(c, hostStores, hostRoles, AutoConfig{})

	if emitted == 0 {
		t.Fatal("emitRoleAsymmetry emitted 0 findings; expected at least 1 role_asymmetry")
	}

	for _, f := range c.items {
		k, _ := f["kind"].(string)
		if k != "role_asymmetry" {
			continue
		}
		dir, _ := f["direction"].(string)
		if dir != "standby_higher" {
			t.Errorf("direction = %q, want %q (standby='aardvark' sorts as hA; pre-fix code labeled hA-wins as primary_higher)", dir, "standby_higher")
		}
		// Sanity: roles_by_host must show standby as the alphabetically-first host.
		roles, _ := f["roles_by_host"].(map[string]string)
		if roles[standbyHost] != string(PALIRoleStandby) {
			t.Errorf("roles_by_host[%s] = %q, want %q", standbyHost, roles[standbyHost], PALIRoleStandby)
		}
	}
}

// Fix W9-2 (part 2): when the primary is alphabetically before the standby
// and the primary has the higher mean, direction must be "primary_higher".
// This case was coincidentally correct before the fix; it verifies the fix
// didn't break the already-correct ordering.
func TestRoleAsymmetryDirectionPrimaryHigher(t *testing.T) {
	// primary="alpha-primary", standby="zeta-standby"
	// Alphabetically: alpha < zeta → hA="alpha" (primary), hB="zeta" (standby).
	// meanA=300, meanB=100 → d > 0 → primary is higher.
	primaryHost := "alpha-primary"
	standbyHost := "zeta-standby"

	hostStores := map[string]*seriesStore{
		primaryHost: makeRoleStore(300.0),
		standbyHost: makeRoleStore(100.0),
	}
	hostRoles := map[string]PALIRole{
		primaryHost: PALIRolePrimary,
		standbyHost: PALIRoleStandby,
	}

	c := &findingCollector{}
	emitted, _ := emitRoleAsymmetry(c, hostStores, hostRoles, AutoConfig{})

	if emitted == 0 {
		t.Fatal("emitRoleAsymmetry emitted 0 findings; expected at least 1 role_asymmetry")
	}

	for _, f := range c.items {
		k, _ := f["kind"].(string)
		if k != "role_asymmetry" {
			continue
		}
		dir, _ := f["direction"].(string)
		if dir != "primary_higher" {
			t.Errorf("direction = %q, want %q", dir, "primary_higher")
		}
	}
}

// Fix W9-2 (part 3): when roles are unknown, node_pair_asymmetry must use
// neutral host-name-based direction labels, not primary/standby vocabulary.
func TestRoleAsymmetryNodePairDirection(t *testing.T) {
	// Both hosts have PALIRoleUnknown — emitRoleAsymmetry should emit
	// node_pair_asymmetry with direction "<hostA>_higher" or "<hostB>_higher".
	hostA := "aardvark-host" // alphabetically first → hA
	hostB := "zebra-host"

	hostStores := map[string]*seriesStore{
		hostA: makeRoleStore(200.0), // hA has higher mean → d > 0 → hA_higher
		hostB: makeRoleStore(100.0),
	}
	hostRoles := map[string]PALIRole{
		hostA: PALIRoleUnknown,
		hostB: PALIRoleUnknown,
	}

	c := &findingCollector{}
	emitted, _ := emitRoleAsymmetry(c, hostStores, hostRoles, AutoConfig{})

	if emitted == 0 {
		t.Fatal("emitRoleAsymmetry emitted 0 findings; expected at least 1 node_pair_asymmetry")
	}

	for _, f := range c.items {
		k, _ := f["kind"].(string)
		if k != "node_pair_asymmetry" {
			continue
		}
		dir, _ := f["direction"].(string)
		if dir == "primary_higher" || dir == "standby_higher" {
			t.Errorf("direction = %q; role vocabulary must not appear in node_pair_asymmetry", dir)
		}
		if dir != hostA+"_higher" {
			t.Errorf("direction = %q, want %q (hA has higher mean, d>0)", dir, hostA+"_higher")
		}
	}
}
