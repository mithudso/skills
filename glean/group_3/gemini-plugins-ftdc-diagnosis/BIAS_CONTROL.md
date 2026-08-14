# BIAS_CONTROL — editorial checklist for FTDC parser detectors

The FTDC parser is designed so that **detection algorithms (Go code) are
metric-agnostic** — they operate on numeric series without knowing what they
measure. **Metric classification (data files)** is MongoDB-aware policy.
This split is enforced by editorial discipline + CI lint, not by file
extension.

When adding or modifying ANY detector, the author AND the reviewer must
walk this checklist. CI lint catches some violations; the rest are caught
here.

## 1. Hardcoded MongoDB metric paths

### Forbidden
- `Kind*` const values containing MongoDB-specific nouns
  (e.g. `KindWTCachePressure`, `KindBufferReplBacklog`)
- `switch metricName { case "wiredTiger.cache.bytes": ... }`
- `if strings.HasPrefix(name, "wt.cache.") { ... }` *inside generic stages*
- `pattern_name` strings in Tier-R diagnosis catalog that name MongoDB
  metric paths (e.g. `"wt_cache_eviction_storm"` — NOT okay; structural
  shape name like `"primary_pali_log_pressure"` is okay because PALI is
  a Tier-5 grandfathered domain stage)

### Allowed (exceptions)
- `strings.HasPrefix(name, "disagg.")` and similar **structural prefix
  tests gating a Tier-5 domain stage**. These are intended architecture:
  the domain stage runs *only when* the prefix is present.
- MongoDB metric paths in **catalog data files** (`metric_docs.go`,
  `metricDocsCatalog`, `capacityRules`, `utilizationRules`). These are
  policy data, not branching logic.
- MongoDB metric paths in **Finding payloads** (the emitted `metric` field).
  Communicating which metric a Finding describes is structural fact,
  not pathology inference.

## 2. Bias vocabulary in Finding payloads

### Forbidden field names
- `class`, `subtype`, `priority`, `impact`, `severity`
- `verdict`, `diagnosis`, `recommendation`, `cause`, `role`
- `tip`, `hint`, `warning_level`, `urgency`

### Allowed alternatives
- Use structural facts: `metric_kind: "counter|gauge"`, `direction: "up|down"`
- Use numeric evidence: `welch_t`, `cohens_d`, `p_value`, `effect_size`,
  `confidence_rank` (NOT `confidence` — see B.13)
- For categorical structure: use `reason_code` with values that name the
  *measurement condition* (`capture_too_short`, `constant_zero`,
  `uint32_wraparound_suspected`), not the pathology

### Tier R exception
`kind:"diagnosis_candidate"` and `kind:"hypothesis_test"` deliberately
NAME pattern shapes (e.g. `pattern_name:"primary_pali_log_pressure"`).
These are constrained to:
1. Pattern names key off Finding *shapes* (which kinds co-occur), NOT
   metric paths.
2. Pattern names that reference Tier-5 domain stages (PALI, replication,
   sharded, host_os) are allowed because the stages are themselves
   structural exceptions.
3. Pattern names that reference MongoDB metric path roots
   (e.g. `wt_cache_*`, `wiredTiger_*`) are FORBIDDEN — same rule as
   generic-stage detectors.

## 3. The "show the choice" rule

Every algorithmic threshold that affects which Findings get emitted
**must** appear in the `run_metadata.thresholds` block. Hidden choices
are bias. The reviewer should grep new code for:
- Magic constants (literal floats / ints) inside emission gates
- New AutoConfig fields that lack a corresponding entry in run_metadata
- Threshold defaults defined as `const` without comments justifying the value

## 4. Symmetric alternatives

Detectors that take a direction or a comparison axis must emit BOTH sides:
- Movers emit `direction:"up"` AND `direction:"down"`
- Correlations run against a *vector* of followers, not a single chosen
  outcome metric
- Changepoint emit doesn't pre-classify the post-shift regime as
  "better" or "worse"

## 5. No pre-filtering by activity / shape

Every metric must enter at least one stage. Pre-filtering by CoV,
absolute magnitude, monotonicity, or any other "is this interesting"
heuristic excludes data from the analyzer's view based on parser-side
guesses. The analyzer can choose to ignore findings later; the parser
must not hide them.

## 6. CHANGELOG references — `detectIssues` and friends

The summary-mode `detectIssues()` function in `ftdc_parser.go` emits
labels like "Cache fill > 95%", "WT-13090", etc. These are explicit
verdicts and they are NOT echoed into the `--auto` output. The summary
mode is for human triage; `--auto` is the analyzer feed and must stay
clean.

If you find yourself wanting to echo a `detectIssues` label into a
Finding, stop. The structural evidence that produced the label belongs
in the Finding; the label itself does not.

## 7. Tier-7 narration scaffolds (`tunable_observation`)

Like Tier R (`diagnosis_candidate`, `hypothesis_test`), `tunable_observation` deliberately NAMES an editorial concept (a mongod tunable). It is allowed under strict guardrails:

1. The Finding cites `supporting_finding_ids` by ID (not metric paths).
2. It NEVER emits a recommended numeric value — only `direction` ∈ {`increase`, `decrease`, `disable`, `enable`, `investigate`, `refuse`}.
3. It is gated to opt-in via `--auto-recommend-tunables`; default `--auto` does not emit it.
4. It carries an explicit `disclaimer` field acknowledging it is editorial.
5. The catalog (`tunable_guardrails.go`) is pure data — no per-tunable Go logic. New entries follow the same review process as `diagnosisPatterns`.
6. Lint Level 1 ALLOWS `tunable`, `direction`, `disclaimer`, `guardrails`, `verification_recipe`, `verification_expected_if_helped`, `verification_expected_if_not_helped`, `anti_recommendation_codes`, `refusal_rationales` on `tunable_observation` Findings ONLY.

The associated anti-recommendation gates (`tunable_guardrails.go::antiRecommendationCatalog`) MAY refuse to recommend a tunable. The refusal is itself editorial but is bounded: each gate names a structural pattern that, when matched, causes the parser to emit `direction:"refuse"` + `anti_recommendation_codes` rather than committing to a direction. E.15 (PALI cluster A) is the prototype: composing on `seal_event + role_asymmetry` routes the user to a code-fix tracker rather than a tunable.

## CI lint

`tools/bias_check.sh` runs two passes:
- **Level 1**: forbidden vocabulary in payloads (existing pre-Tier-C ban).
- **Level 2 (Tier C)**: literal MongoDB metric path regexes in any
  generic-stage detector branch, `Kind*` const, or `pattern_name` catalog
  entry. Allowlisted: Tier-5 domain-stage files, catalog data files,
  Finding `metric` field.

Run via:

    ./tools/bias_check.sh
