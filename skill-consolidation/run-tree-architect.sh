#!/usr/bin/env bash
# run-tree-architect.sh — durable DAILY entrypoint for the skill-tree-architect skill.
#
# Wired into the user's crontab (see HUB-STRATEGY.md scheduler + run-autohub.sh):
#   30 8 * * *  $HOME/.claude/skill-consolidation/run-tree-architect.sh \
#                 >> $HOME/.claude/skill-consolidation/tree-architect.log 2>&1
#
# Read-only first (cheap, every day): run the placement/balance audit + the new-family detector.
# Two tiers, mirroring run-autohub.sh's "only escalate on a real need" stance:
#   ESCALATE (hand off to the skill-tree-architect agent): an over-cap hub (lost trigger keywords =
#            real breakage) OR a >=8 unhubbed family (needs a new hub). These are structural needs.
#   WATCH (log only, never spins an agent): a hub >=95% of the 1536-char hard cap, a hub
#            description >1000 chars (Medium — Glean export cap; single definition: sko Pass M),
#            or a heuristic misplaced-spoke candidate. Surfaced in the daily log for review.
# If nothing escalates, log one line and exit — nothing is mutated, no agent is launched.
#
# Like run-autohub.sh, this does NOT bypass permission gates and does NOT auto-apply structural
# moves: skill-tree-architect surfaces folding / splits / registry-sync for review
# (~/.claude/skills is not git-backed). Zero-risk idempotent repairs only, on review.
set -euo pipefail

# cron's PATH (/usr/bin:/bin) has neither homebrew node nor ~/.local/bin claude — fix lookup for
# ALL invocations below (four bare `node` calls + the `claude` escalation handoff) in one line.
export PATH="/opt/homebrew/bin:$HOME/.local/bin:$PATH"
command -v node >/dev/null || { echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] FATAL: node not found on PATH"; exit 1; }

CONS_DIR="$HOME/.claude/skill-consolidation"
cd "$CONS_DIR"
TS="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

# Read-only signals, parsed without jq (each script prints machine JSON with --json).
# Two-tier description cap (single definition: sko Pass M): >1536 = High (harness truncation,
# ESCALATE via counts.overCapHubs at the default cap); >1000 = Medium (Glean export cap, WATCH).
SIGS="$(node audit-placement.mjs --json | node -e 'let d="";process.stdin.on("data",c=>d+=c).on("end",()=>{const j=JSON.parse(d);const c=j.counts||{};const near=(j.hubBalance||[]).filter(h=>h.capPct>=95).length;const over1000=(j.hubBalance||[]).filter(h=>(h.descLen||0)>1000).length;process.stdout.write(`${c.overCapHubs||0} ${near} ${c.misplacedSpokes||0} ${over1000}`)})')"
READY="$(node detect-candidates.mjs --json | node -e 'let d="";process.stdin.on("data",c=>d+=c).on("end",()=>{const j=JSON.parse(d);process.stdout.write(String((j.readyCandidates||[]).length))})')"
read -r OVERCAP NEARCAP MISPLACED OVER1000 <<< "$SIGS"
ESCALATE=$(( OVERCAP + READY ))

# WATCH-tier staleness signal (gap-staleness-consumer): LOCAL-FALLBACK estimate of /dr-researched
# concepts >90 days old, counting both the .concepts map and stray top-level concept nodes
# (documented dual-write drift). Approximate by design — the authoritative queue is
# tam_concept_tree_list(staleOnly:true), consumed by `/dr --refresh`; a count/queue mismatch means
# MCP-vs-local drift, not a consumer bug. Failure-isolated: "?" surfaces an unreadable local
# fallback as a visible watch item and can never abort the ESCALATE path. Never added to ESCALATE.
STALE_N=$(node -e 'const j=require(process.env.HOME+"/.claude/concept-tree.json");const vals=[...Object.values(j.concepts||{}),...Object.entries(j).filter(([k,v])=>k!=="concepts"&&v&&typeof v==="object"&&v.researchedAt).map(([,v])=>v)];const cut=Date.now()-90*864e5;process.stdout.write(String(vals.filter(n=>new Date(n.researchedAt)<cut).length))' 2>/dev/null || echo "?")
SUMMARY="overCap=$OVERCAP newFamilies>=8=$READY | watch: nearCap>=95%=$NEARCAP descGt1000=$OVER1000 misplaced=$MISPLACED stale=$STALE_N"
case "$STALE_N" in
  ''|*[!0-9]*|0) : ;;
  *) echo "[$TS] watch: ready to run: /dr --refresh --budget-minutes=30 ($STALE_N stale concepts)" ;;
esac

if [ "$ESCALATE" -eq 0 ]; then
  echo "[$TS] skill-tree-architect daily: no structural action needed ($SUMMARY)."
  exit 0
fi

echo "[$TS] skill-tree-architect daily: structural finding ($SUMMARY) — handing off to skill-tree-architect."
if command -v claude >/dev/null 2>&1; then
  claude -p "Invoke the skill-tree-architect skill. Run Phase 1 (analyze, read-only) over ~/.claude/skills, then emit the Phase 2 ranked rebalance plan. Apply ONLY zero-risk idempotent repairs; surface folding / splits / registry sync for review (the tree is not git-backed). Daily detector summary: $SUMMARY." \
    || echo "[$TS] claude headless returned non-zero (likely awaiting permission). Run the skill-tree-architect skill manually."
else
  echo "[$TS] 'claude' CLI not on PATH — run the skill-tree-architect skill manually."
fi
