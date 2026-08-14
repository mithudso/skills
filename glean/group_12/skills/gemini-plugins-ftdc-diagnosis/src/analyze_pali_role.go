// analyze_pali_role.go — PALI primary/standby role inference.
//
// D6 from the P0 plan: many downstream Findings (seal_event, role_asymmetry,
// tunable applicability filters, sweep label inference) need to know which
// PALI role a host is playing. The existing isPALICapture() gate tells us
// "this is a PALI capture" but not which side.
//
// Inference is FTDC-only and deterministic. The function is a pure
// function of the seriesStore — no environment lookups, no flag reads.
//
// Decision rule:
//
//   primary  ← `disagg.logServer.appendLogSuccess` is increasing (host
//              writes the log) AND the host appears to be the elected
//              primary (replSetUpdatePosition counter rising, indicating
//              secondaries push update positions to this node).
//   standby  ← `disagg.phylog.committedToMaterializedLagMillis` is
//              actively measured (host applies the log) AND the host
//              does NOT show primary-write activity AND is not the
//              elected primary.
//   unknown  ← otherwise (insufficient evidence; e.g. a non-PALI capture
//              or a brief window where the host wasn't doing PALI work).
//
// Bias-control: the role is a STRUCTURAL inference; Finding payloads that
// use it must STILL avoid `role` as a value field (per BIAS_CONTROL.md
// Tier A.6). Findings instead include the structural EVIDENCE that drove
// the inference and let downstream consumers re-derive the role.
package main

// PALIRole enumerates the inferable roles. Values are stable strings safe
// to embed in Finding payloads (when bias-control allows).
type PALIRole string

const (
	PALIRoleUnknown PALIRole = "unknown"
	PALIRolePrimary PALIRole = "primary"
	PALIRoleStandby PALIRole = "standby"
)

// PALIRoleEvidence is the bag of facts the inference reads. Returned
// alongside the role so callers can include the evidence verbatim in
// Findings, preserving the bias-control posture (we show the choice).
type PALIRoleEvidence struct {
	HasAppendLogSuccess     bool    // disagg.logServer.appendLogSuccess present
	AppendLogSuccessIncrPct float64 // fraction of consecutive samples where it increased
	AppendLogSuccessDelta   float64 // last-minus-first value (sanity check)
	HasPhylogMaterLag       bool    // disagg.phylog.committedToMaterializedLagMillis present
	PhylogMaterLagNonZeroN  int     // number of samples with non-zero materialization lag
	HasElectionID           bool    // serverStatus.repl.electionId present
	HasUpdatePositionCnt    bool    // metrics.repl.network.replSetUpdatePosition present
	UpdatePositionIncrPct   float64 // fraction of consecutive samples where it increased
}

// paliRole inspects the store and returns (role, evidence). It returns
// PALIRoleUnknown + a populated evidence struct if the metrics needed for
// inference are absent — never fabricates a role.
func paliRole(store *seriesStore) (PALIRole, PALIRoleEvidence) {
	ev := PALIRoleEvidence{}

	const (
		appendLogSuccessName  = "serverStatus.metrics.disagg.logServer.appendLogSuccess"
		phylogMaterLagName    = "serverStatus.metrics.disagg.phylog.committedToMaterializedLagMillis"
		electionIDName        = "serverStatus.repl.electionId"
		updatePositionCntName = "serverStatus.metrics.repl.network.replSetUpdatePosition"
	)

	if !isPALICapture(store) {
		return PALIRoleUnknown, ev
	}

	// Materialize the four candidate series. We only need lightweight
	// stats (last-minus-first delta, fraction-of-samples-increased) so
	// we can hand-roll them off the store columns without paying for a
	// full materializeFromColumns call.
	if id, ok := store.nameToID[appendLogSuccessName]; ok {
		ev.HasAppendLogSuccess = true
		ev.AppendLogSuccessIncrPct, ev.AppendLogSuccessDelta = colIncreaseStats(store, id)
	}
	if id, ok := store.nameToID[phylogMaterLagName]; ok {
		ev.HasPhylogMaterLag = true
		// PhylogMaterLagNonZeroN is informational evidence (returned to callers)
		// and is NOT a gate in standbyLikely — a fully caught-up standby has
		// lag == 0 throughout but is still a standby.
		ev.PhylogMaterLagNonZeroN = colNonZeroCount(store, id)
	}
	if _, ok := store.nameToID[electionIDName]; ok {
		ev.HasElectionID = true
	}
	if id, ok := store.nameToID[updatePositionCntName]; ok {
		ev.HasUpdatePositionCnt = true
		ev.UpdatePositionIncrPct, _ = colIncreaseStats(store, id)
	}

	// Primary criteria. We require BOTH that the host wrote to the
	// PALI log AND that secondaries are talking back (replSetUpdatePosition
	// increasing). Either alone could mis-fire on log-shipping standbys
	// or quiescent primaries.
	primaryLikely := ev.HasAppendLogSuccess &&
		ev.AppendLogSuccessIncrPct > 0.10 && // counter rose in ≥10% of consecutive samples
		ev.AppendLogSuccessDelta > 0 &&
		ev.HasUpdatePositionCnt &&
		ev.UpdatePositionIncrPct > 0.10

	// Standby criteria. The standby reads + materializes; lag metric is
	// non-zero for at least some samples. Crucially the standby does NOT
	// show appendLogSuccess writes (rules out primary).
	standbyLikely := ev.HasPhylogMaterLag &&
		!primaryLikely

	switch {
	case primaryLikely && !standbyLikely:
		return PALIRolePrimary, ev
	case standbyLikely && !primaryLikely:
		return PALIRoleStandby, ev
	default:
		return PALIRoleUnknown, ev
	}
}

// colIncreaseStats walks the column for metric id and returns:
//
//	incrPct: fraction of consecutive samples where the value strictly
//	         increased (0 if N<2). Counter-friendly metric.
//	delta:   last - first value across the column.
//
// The column is decoded via seriesStore.readColumnInt64 (cheap — values
// are already int64 in the column store) and freed at function exit.
func colIncreaseStats(store *seriesStore, id int) (incrPct, delta float64) {
	N := store.nSamples()
	if N < 2 {
		return 0, 0
	}
	col := store.readColumnInt64(id, nil)
	if len(col) < 2 {
		return 0, 0
	}
	first := col[0]
	last := col[len(col)-1]
	delta = float64(last - first)
	incr := 0
	for i := 1; i < len(col); i++ {
		if col[i] > col[i-1] {
			incr++
		}
	}
	incrPct = float64(incr) / float64(len(col)-1)
	return incrPct, delta
}

// colNonZeroCount returns the number of samples where the metric id's
// value was strictly greater than zero. Used to detect "the metric is
// being actively measured" vs "always zero" (which would imply the
// codepath isn't active).
func colNonZeroCount(store *seriesStore, id int) int {
	col := store.readColumnInt64(id, nil)
	n := 0
	for _, v := range col {
		if v > 0 {
			n++
		}
	}
	return n
}
