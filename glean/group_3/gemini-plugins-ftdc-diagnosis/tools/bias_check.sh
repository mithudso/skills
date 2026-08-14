#!/usr/bin/env bash
# Tier C bias-control CI lint. Two passes.
#
# Level 1: forbidden vocabulary in Finding payloads (pre-existing
# CHANGELOG-ban). Greps for class/subtype/priority/impact/verdict/
# diagnosis/recommendation/tip/hint/severity as JSON keys in code that
# constructs Finding payloads. False positives are possible (a comment
# mentioning "verdict" is fine); the goal is fail-loud so a reviewer
# eyeballs the hit.
#
# Level 2: literal MongoDB metric paths in detector branching logic.
# Scans Kind* const values, switch-case-on-metric-name, and
# strings.HasPrefix tests inside generic stages. Allowlisted: Tier-5
# domain-stage files (analyze_domain.go, analyze_crosshost.go) where
# prefix-gating is the intended architecture; catalog data files
# (metric_docs.go); Finding `metric` field assignments.
#
# Usage:
#   ./tools/bias_check.sh
# Exit code: 0 = clean; 1 = violations found.

set -uo pipefail

cd "$(dirname "$0")/../src" || exit 1

VIOLATIONS=0

echo "=== Level 1: forbidden bias vocabulary in payload assignments ==="
# Grep for forbidden words appearing as JSON keys (quoted, followed by colon)
# inside Go map literals.
PATTERNS='"class"\s*:|"subtype"\s*:|"priority"\s*:|"impact"\s*:|"verdict"\s*:|"diagnosis"\s*:|"recommendation"\s*:|"tip"\s*:|"hint"\s*:|"severity"\s*:|"cause"\s*:|"role"\s*:|"warning_level"\s*:|"urgency"\s*:'
HITS=$(grep -nE "$PATTERNS" analyze*.go 2>/dev/null | grep -v "_test.go" | grep -v "BIAS_CONTROL")
if [ -n "$HITS" ]; then
    echo "VIOLATION: forbidden vocabulary in Finding payload(s):"
    echo "$HITS"
    VIOLATIONS=$((VIOLATIONS + 1))
fi

echo ""
echo "=== Level 2: literal MongoDB metric path regexes in detector logic ==="
# Files allowlisted as intentional Tier-5 domain stages where prefix-gating
# is the architecture.
ALLOWLIST="analyze_domain.go|analyze_crosshost.go|metric_docs.go|analyze_capacity.go|analyze_glossary.go|ftdc_parser.go|analyze_diagnosis.go|analyze_capabilities.go|analyze_baseline.go"

# MongoDB metric path prefixes that should NOT appear in generic-stage logic.
# These would be hardcoded pathology inferences.
MONGO_PREFIXES='"wt\.|"repl\.|"metrics\.|"asserts\.|"opcounters\.|"catalogCache\.|"sharding\.|"connections\.|"extra_info\.|"serverStatus\.metrics\.|"systemMetrics\.'

# Pattern A: switch-case-on-metric-name
HITS_A=$(grep -nE "case[[:space:]]+($MONGO_PREFIXES)" analyze*.go 2>/dev/null | \
    grep -vE "($ALLOWLIST)" | grep -v "_test.go")
if [ -n "$HITS_A" ]; then
    echo "VIOLATION: literal MongoDB path in switch-case:"
    echo "$HITS_A"
    VIOLATIONS=$((VIOLATIONS + 1))
fi

# Pattern B: strings.HasPrefix on hardcoded MongoDB path
HITS_B=$(grep -nE "strings\.HasPrefix\(.+,[[:space:]]*($MONGO_PREFIXES)" analyze*.go 2>/dev/null | \
    grep -vE "($ALLOWLIST)" | grep -v "_test.go")
if [ -n "$HITS_B" ]; then
    echo "VIOLATION: strings.HasPrefix on MongoDB path in non-domain stage:"
    echo "$HITS_B"
    VIOLATIONS=$((VIOLATIONS + 1))
fi

# Pattern C: Kind* const values containing MongoDB-specific nouns
HITS_C=$(grep -nE "^[[:space:]]*Kind[A-Z][a-zA-Z]*[[:space:]]*=[[:space:]]*\"[a-z_]*(wt|repl|opcounter|wiredtiger|cache|oplog|catalog|sharding|disagg)" analyze.go 2>/dev/null | \
    grep -v "_test.go" | grep -v "// Tier")
if [ -n "$HITS_C" ]; then
    echo "WARNING: Kind* const may contain MongoDB-specific noun (review):"
    echo "$HITS_C"
    # WARNING, not VIOLATION — Tier-5 grandfathered kinds intentionally
    # name MongoDB shapes (term_change, oplog_window_low, materialization_lag,
    # etc.). The grep is here to flag NEW kinds for human review.
fi

echo ""
if [ "$VIOLATIONS" -eq 0 ]; then
    echo "[PASS] bias_check.sh: 0 violations"
    exit 0
fi
echo "[FAIL] bias_check.sh: $VIOLATIONS violation(s)"
exit 1
