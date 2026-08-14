<!-- hub-reference-banner -->
> **Reference file — part of the `tam-operations` hub.** Formerly the standalone `operator-report-generator` skill.
> Sibling topics in this family are now reference files under the hubs (`tam-operations`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: operator-report-generator
description: >
  Operator report generation expert — shift handoff (SBAR), meeting prep documents, quality
  validation frameworks, data freshness scoring, BLUF writing, and template-driven report engines.
  TRIGGER: building or reviewing automated report generation, shift handoff systems, meeting prep
  automation, executive summary pipelines, data freshness scoring, BLUF writing, template-driven
  report engines, quality validation for operational reports.
  SKIP: general dashboard or metrics visualization (use frontend-ui for chart/UI work, devops-infra
  for monitoring/observability); ad-hoc data analysis without a repeating report structure; raw data
  pipeline design without a document-rendering layer.
version: 1.2.2
updated: 2026-06-01
category: developer
tags:
  - reports
  - handoff
  - meeting-prep
  - quality-validation
  - nlg
  - templates
  - SBAR
  - BLUF
  - freshness
  - confidence-scoring
triggers:
  - shift handoff
  - meeting prep
  - report generation
  - report template
  - SBAR
  - BLUF
  - data freshness
  - quality validation
  - executive summary
  - operator report
  - handoff report
  - pre-read document
related_skills:
  - devops-infra
  - tam-operations (references/incident-response)
---

# Operator Report Generator

Expert reference for automated operator report generation: shift handoff reports (SBAR), meeting prep documents, quality validation, data freshness scoring, BLUF executive writing, and template-driven report engines.

## When to use this skill

- Building or reviewing a shift handoff report generator
- Designing meeting prep document automation
- Adding quality validation or data freshness scoring to any report pipeline
- Writing executive summaries that follow BLUF conventions
- Creating template engines for recurring operational reports
- Reviewing report output for completeness, staleness, or coverage gaps

## When NOT to use this skill

- General dashboard or metrics visualization → use `frontend-ui` (chart/UI work) or `devops-infra` (monitoring/observability)
- Ad-hoc data analysis without a repeating report structure
- Raw data pipeline design without a document-rendering layer

---

## Report pipeline

```
Data sources → Collection → Enrichment → Validation → Rendering → Delivery
    |              |             |             |            |           |
  cases         merge &      add derived   freshness &  template   channel
  metrics        dedupe       fields        completeness  engine    routing
  alerts                      confidence    checks
  comments                    scores
```

**Key principles:**
- **Deterministic core, flexible shell.** Validation and scoring logic must be deterministic and testable. Template rendering can vary by audience.
- **Fail loud on data gaps.** A report that hides missing data is worse than no report. Surface gaps explicitly.
- **Freshness over completeness.** A timely report with flagged gaps beats a delayed report waiting for stale data.
- **Audience-aware formatting.** Same data renders differently for shift handoff (detailed, action-oriented) vs. executive pre-read (BLUF, high-level).

---

## Canonical report data schema

All templates and validators operate on this shared shape:

```javascript
const ReportData = {
  // Metadata (required)
  reportType:   'shift-handoff' | 'meeting-prep' | 'executive-summary',
  generatedAt:  '2026-05-26T14:30:00Z',
  generatedBy:  'operator-report-generator/1.2.1',

  // Cases
  cases: [{
    caseId: '12847', title: 'Atlas cluster unreachable',
    severity: 'S1', status: 'In Progress', owner: 'alice',
    age: 172800000,  // ms since creation
    lastUpdate: '2026-05-26T13:45:00Z',
    escalated: false,
    slaDeadline: '2026-05-26T16:00:00Z',
  }],
  openCaseCount: 5, critUrgentCount: 2,

  // Derived metrics
  newEscalations: 1,
  slaDeadlines: [{ caseId: '12847', deadline: '...', hoursRemaining: 2 }],
  severityTrend: { s1: { start: 1, end: 2, delta: +1 } },

  // Assessment
  riskLevel: 'elevated',  // low | moderate | elevated | critical
  riskRationale: 'S1 count doubled during shift',
  bluf: 'string, 1-3 sentences',
  topAction: 'Prioritize case #12847 — SLA breach in 2h',
  actionItems: [{ priority: 'P1', description: '...', owner: 'bob' }],

  // Source freshness
  sources: {
    cases:   { fetchedAt: 1748268600000, fromCache: false },
    alerts:  { fetchedAt: 1748268590000, fromCache: false },
    metrics: { fetchedAt: 1748268000000, fromCache: false },
  },

  // Quality (computed after validation)
  qualityScore: 85, freshnessScore: 92, overallFreshnessLabel: 'Fresh',
};
```

---

## Shift handoff reports (SBAR)

SBAR originated in the US Navy, adopted by healthcare as the gold standard for handoff communication. Studies show up to 70% of sentinel events trace to ineffective handoffs.

| Section | Purpose | Operations translation |
|---------|---------|----------------------|
| **S** — Situation | What is happening right now? | Active incidents, open cases, current alert state |
| **B** — Background | What context does the receiver need? | Recent changes, escalation history, SLA deadlines |
| **A** — Assessment | What do you think is going on? | Root cause hypotheses, severity trends, risk assessment |
| **R** — Recommendation | What should happen next? | Prioritized action items, watch items, follow-ups |

### SBAR template (abbreviated)

```markdown
# Shift Handoff Report
**Outgoing:** {{outgoing_operator}} | **Incoming:** {{incoming_operator}}
**Shift:** {{shift_start}} - {{shift_end}}
**Data freshness:** {{overall_freshness_label}} ({{freshness_score}}/100)

## BLUF
{{bluf_statement}}

## Situation
{{#active_incidents}}
- **[{{severity}}] {{title}}** ({{case_id}}) — {{status}} | {{owner}} | {{age_display}}
{{/active_incidents}}
**Open:** {{open_case_count}} | **Critical/Urgent:** {{crit_urgent_count}}

## Background
{{#recent_changes}}- {{timestamp_relative}}: {{description}}{{/recent_changes}}
**SLA deadlines in next 4 hours:** {{#upcoming_sla_deadlines}}...{{/upcoming_sla_deadlines}}

## Assessment
{{assessment_narrative}}
**Risk:** {{risk_level}} — {{risk_rationale}}

## Recommendation
{{#action_items}}- [ ] **{{priority}}** {{description}} (assign: {{owner}}){{/action_items}}

---
**Quality:** {{quality_score}}/100 | **Gaps:** {{coverage_gap_summary}}
```

---

## Meeting prep structure

Meeting prep documents optimize for pre-read consumption — stakeholders with 5-10 minutes before a live discussion.

**Agenda item automation** — auto-generate items from:
1. Unresolved action items from the previous meeting
2. Cases with severity changes since the last meeting
3. SLA violations or near-misses
4. Metric anomalies that exceed threshold
5. New high-severity cases
6. Pending decisions by case status or workflow state

**Meeting prep data sources:**

| Source | What it provides | Freshness SLA |
|--------|-----------------|---------------|
| Case management | Open cases, status, severity, owner | < 5 min |
| Metrics/monitoring | KPIs, trends, SLA compliance | < 15 min |
| Previous meeting notes | Prior action items, decisions | Static |
| Alert system | Active/recent alerts | < 5 min |
| Escalation log | Recent escalations, customer impact | < 10 min |

---

## Quality validation framework

Every generated report must pass validation across five dimensions before delivery:

| Dimension | What it checks | Fail condition |
|-----------|---------------|----------------|
| **Completeness** | All required sections populated | Any required section empty or contains template markers |
| **Freshness** | Source data within acceptable age | Any data source exceeds its freshness SLA |
| **Accuracy** | Cross-referenced counts and totals | Case count in summary differs from detailed list |
| **Coverage** | Report spans expected scope | Known cases missing; entire severity tier absent |
| **Consistency** | No contradictions between sections | BLUF says "no critical issues" but Situation lists S1 cases |

```javascript
function validateReport(report, sources, config) {
  const findings = [];

  // 1. Completeness
  for (const section of config.requiredSections) {
    if (!report[section]?.trim()) {
      findings.push({ dimension: 'completeness', severity: 'error',
        message: `Required section "${section}" is empty` });
    }
  }

  // 2. Freshness
  for (const [sourceName, sourceData] of Object.entries(sources)) {
    const ageMs = Date.now() - sourceData.fetchedAt;
    const slaMs = config.freshnessSLA[sourceName] ?? config.defaultFreshnessSLA;
    if (ageMs > slaMs) {
      findings.push({ dimension: 'freshness',
        severity: ageMs > slaMs * 2 ? 'error' : 'warning',
        message: `Source "${sourceName}" is ${formatAge(ageMs)} old (SLA: ${formatAge(slaMs)})` });
    }
  }

  // 3. Accuracy — cross-check derived totals
  const listedCount = report.situations?.length ?? 0;
  if (report.openCaseCount !== listedCount) {
    findings.push({ dimension: 'accuracy', severity: 'error',
      message: `Summary says ${report.openCaseCount} open cases but ${listedCount} listed` });
  }

  // 4. Coverage
  const coveredSevs = new Set(report.situations?.map(s => s.severity));
  for (const sev of config.expectedSeverities) {
    if (!coveredSevs.has(sev) && sources.cases.bySeverity[sev]?.length > 0) {
      findings.push({ dimension: 'coverage', severity: 'warning',
        message: `Severity "${sev}" has cases but none appear in report` });
    }
  }

  // 5. Consistency — BLUF vs body contradiction
  if (report.bluf && report.critUrgentCount > 0) {
    const blufLower = report.bluf.toLowerCase();
    if (blufLower.includes('no critical') || blufLower.includes('all clear')) {
      findings.push({ dimension: 'consistency', severity: 'error',
        message: `BLUF claims no critical issues but ${report.critUrgentCount} critical/urgent cases exist` });
    }
  }

  const errorCount   = findings.filter(f => f.severity === 'error').length;
  const warningCount = findings.filter(f => f.severity === 'warning').length;
  return { score: Math.max(0, 100 - (errorCount * 20) - (warningCount * 5)),
           passed: errorCount === 0, findings };
}
```

### Quality score interpretation

| Score | Label | Action |
|-------|-------|--------|
| 90-100 | High confidence | Deliver automatically |
| 70-89 | Moderate confidence | Deliver with quality warnings visible |
| 50-69 | Low confidence | Flag for human review before delivery |
| 0-49 | Undeliverable | Block delivery, alert the operator |

---

## Data freshness model

Freshness measures the gap between when an event happens and when the report reflects it. A dashboard can render in 50ms while showing two-hour-old data. Freshness and latency are distinct concerns.

### Freshness SLA configuration

```javascript
const FRESHNESS_SLA = {
  'cases':          5 * 60 * 1000,   // 5 minutes
  'alerts':         60 * 1000,        // 1 minute
  'metrics':        15 * 60 * 1000,   // 15 minutes
  'escalation_log': 10 * 60 * 1000,   // 10 minutes
  'comments':       5 * 60 * 1000,    // 5 minutes
  'prior_actions':  Infinity,         // static reference
};
```

### Confidence scoring per section

```javascript
function computeSectionConfidence(section, sources) {
  let score = 100;

  // Freshness penalty: linear decay from SLA boundary to 2x SLA
  for (const src of section.dependsOn) {
    const ageMs = Date.now() - sources[src].fetchedAt;
    const slaMs = FRESHNESS_SLA[src] ?? 300000;
    if (ageMs > slaMs) {
      const overageRatio = Math.min((ageMs - slaMs) / slaMs, 1.0);
      score -= overageRatio * 30;  // max 30-point penalty per source
    }
  }

  // Completeness penalty
  const presentFields = section.fields.filter(f => f.value != null).length;
  score -= (1 - presentFields / section.fields.length) * 40;

  // Cache penalty
  for (const src of section.dependsOn) {
    if (sources[src].fromCache) score -= 15;
  }

  return Math.max(0, Math.round(score));
}
```

**Freshness labels in rendered reports:**
- **Fresh** (green): all sources within SLA
- **Aging** (yellow): one or more sources between 1x and 2x SLA
- **Stale** (red): one or more sources beyond 2x SLA

---

## BLUF writing patterns

BLUF (Bottom Line Up Front) is a US military communication standard placing the conclusion at the very beginning. A BLUF is 1-3 sentences maximum.

**Formula:** `[Current state] + [Key change/risk] + [Recommended action or decision needed]`

**Examples:**

*Shift handoff:* "Shift ends with 3 open S1 cases (up from 1 at shift start). Cluster migration for Acme Corp stalled at 60% and needs manual intervention. Incoming operator should prioritize case #12847 — customer SLA breach in 2 hours."

*Executive summary:* "All critical cases resolved. Support response time improved 12% week-over-week. One emerging risk: 3 accounts on legacy driver versions will lose compatibility after the October release."

**BLUF anti-patterns:**

| Anti-pattern | Example | Fix |
|-------------|---------|-----|
| Burying the lead | "During the past week, the team worked on..." | Lead with the outcome, not the timeline |
| Vague hedging | "Things are mostly okay with some concerns" | Use counts, severities, deadlines |
| Missing the action | "We have 5 open cases" | Add what should happen: "...need to escalate 2 by EOD" |
| Too long | Full paragraph of context | 1-3 sentences; context goes in Background |

### BLUF generation algorithm

```javascript
function generateBLUF(reportData) {
  const parts = [];

  // State: current position
  parts.push(reportData.critUrgentCount > 0
    ? `${reportData.critUrgentCount} critical/urgent case${reportData.critUrgentCount > 1 ? 's' : ''} active`
    : 'No critical or urgent cases active');

  // Change: what shifted since last report
  if (reportData.newEscalations > 0) {
    parts.push(`${reportData.newEscalations} new escalation${reportData.newEscalations > 1 ? 's' : ''} since last report`);
  }

  // Risk: upcoming SLA
  const imminentSLA = reportData.slaDeadlines
    .filter(d => d.hoursRemaining < 4)
    .sort((a, b) => a.hoursRemaining - b.hoursRemaining);
  if (imminentSLA.length > 0) {
    parts.push(`nearest SLA breach in ${imminentSLA[0].hoursRemaining}h (${imminentSLA[0].caseId})`);
  }

  // Action: primary recommendation
  if (reportData.topAction) parts.push(reportData.topAction);

  return parts.join('. ') + '.';
}
```

---

## Template engine design

Separate data from presentation — same underlying data renders to multiple formats (markdown, HTML, Slack blocks, PDF) for different audiences.

```
Template Registry
    +-- shift-handoff.md.hbs
    +-- meeting-prep.md.hbs
    +-- executive-summary.md.hbs
    +-- slack-digest.json.hbs

renderReport(templateId, data, format)
    → resolve template → inject helpers → render → validate → return
```

**Template customization points:**

| Customization | Mechanism | Example |
|--------------|-----------|---------|
| Section visibility | Config flags | `showAppendix: false` |
| Section ordering | Ordered section list | `['bluf', 'situation', 'recommendation']` |
| Field inclusion | Per-section allowlist | Show `severity` but not `internalNotes` |
| Audience variant | Template variant | "Executive" hides technical detail |
| Conditional sections | Data-driven predicates | Only show "Escalations" if escalations exist |

---

## Data collection with error handling

```javascript
async function collectSources(sourceConfigs, { timeoutMs = 10000 } = {}) {
  const results = {};
  await Promise.allSettled(
    Object.entries(sourceConfigs).map(async ([name, config]) => {
      const controller = new AbortController();
      const timer = setTimeout(() => controller.abort(), timeoutMs);
      try {
        const data = await config.fetcher({ signal: controller.signal });
        results[name] = { data, fetchedAt: Date.now(), fromCache: false, error: null };
      } catch (err) {
        const cached = await config.getCache?.();
        results[name] = {
          data: cached?.data ?? null, fetchedAt: cached?.fetchedAt ?? 0,
          fromCache: true, error: err.name === 'AbortError' ? 'timeout' : err.message,
        };
      } finally { clearTimeout(timer); }
    })
  );
  return results;
}
```

---

## Anti-patterns

| Anti-pattern | Risk | Remedy |
|-------------|------|--------|
| Silent data gaps | Report looks complete but source failed silently | Always surface source status; never swallow fetch errors |
| Stale without indicator | Two-hour-old data with no freshness label | Every section must display its data age |
| One-size-fits-all | Same verbose report sent to operators AND executives | Audience-specific templates |
| Hardcoded thresholds | "Critical" means > 5 cases, baked into code | Make thresholds configurable |
| Unbounded growth | Report length grows with case count | Cap section length; overflow to appendix |
| LLM hallucination in summaries | AI-generated narrative invents case details | Ground generated text against structured data |
| No validation before send | Broken template sends `{{undefined}}` | Run `validateReport()` before every delivery |
| Missing operator acknowledgement | Handoff sent but incoming operator never confirms | Require explicit acknowledgement; escalate if unacknowledged |

---

## Implementation checklist

- [ ] All sources identified, fetch errors handled, results timestamped
- [ ] Heterogeneous source data mapped to canonical report schema
- [ ] Freshness SLAs defined per source, enforced at validation, displayed in output
- [ ] BLUF auto-generated from structured data, validated against body content
- [ ] SBAR sections each have a clear data source mapping, never silently empty
- [ ] Quality validation runs before every delivery across all five dimensions
- [ ] Confidence scores computed per section and overall, displayed to consumer
- [ ] Templates versioned and tested against fixtures; multiple formats supported
- [ ] Audience variants use separate templates (shift handoff, meeting prep, executive)
- [ ] Delivery confirmation tracked; failures alerted; handoffs acknowledged
- [ ] Edge cases tested: zero items, maximum items, all-stale sources, partial data
- [ ] No AI-isms in generated text ("delve", "robust", "seamless", "paradigm", "leverage")
