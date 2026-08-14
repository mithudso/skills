<!-- hub-reference-banner -->
> **Reference file — part of the `tam-operations` hub.** Formerly the standalone `account-health-scorer` skill.
> Sibling topics in this family are now reference files under the hubs (`tam-operations`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: account-health-scorer
version: 1.2.3
updated: "2026-06-01"
description: >
  Account health scoring algorithms — weighted composite scores, signal categories, grading
  heuristics, time-series trending, anomaly detection, and multi-dimensional scoring.

  TRIGGER: user is designing a new customer health scoring system, reviewing or refactoring an
  existing health score implementation, choosing signal weights or normalization strategies,
  adding time-series trending or anomaly detection to account dashboards, or normalizing scores
  across accounts of different sizes or tiers.

  SKIP: ML churn prediction model training from scratch (use a data science framework),
  financial credit scoring or fraud detection (different domain), real-time application
  performance monitoring (use APM tools); designing the intervention to move a low score
  by changing customer behavior — adoption, enablement, motivation (use applied-psychology).
category: developer
tags: [health-scoring, customer-success, algorithms, grading, analytics, churn-prediction]
triggers:
  - health score
  - account health
  - churn prediction
  - customer health
  - scoring algorithm
  - health grading
  - account scoring
  - trajectory analysis
  - anomaly detection health
related_skills:
  - account-artifacts-collector
  - tam-account-reports
  - tam-expertise
whenToUse:
  - "design a new customer health scoring system"
  - "review or refactor an existing health score implementation"
  - "choose signal weights, normalization strategies, or grading thresholds"
  - "add time-series trending or anomaly detection to account dashboards"
  - "normalize health scores across accounts of different sizes or tiers"
whenNotToUse:
  - "building ML churn prediction models from scratch — use a data science framework"
  - "financial credit scoring or fraud detection — different domain and regulations"
  - "real-time application performance monitoring — use APM tools"
---

# Account Health Scorer

## Overview

An account health score is a composite metric that blends behavioral, engagement, operational, and sentiment signals into a single predictive indicator of customer stability.

Key design goals:
- Predict churn 60–90 days before cancellation (industry benchmark: 63 days average lead time)
- Surface expansion opportunities alongside risk
- Normalize fairly across account sizes and tiers
- Provide actionable grades, not just numbers
- Support time-series trending for trajectory analysis

---

## Core Concepts

### Utility Functions

```javascript
function avg(arr) {
  return arr.length > 0 ? arr.reduce((sum, v) => sum + v, 0) / arr.length : 0;
}

function stdDev(arr) {
  if (arr.length < 2) return 0;
  const mean = avg(arr);
  const variance = arr.reduce((sum, v) => sum + (v - mean) ** 2, 0) / (arr.length - 1);
  return Math.sqrt(variance);
}

function daysBetween(dateA, dateB) {
  return Math.abs(new Date(dateB) - new Date(dateA)) / 86400000;
}

const GRADE_ORDER = ['A', 'B', 'C', 'D', 'F'];

function gradeToNumeric(grade) {
  return GRADE_ORDER.indexOf(grade);  // A=0 (best), F=4 (worst)
}

function isWorseGrade(gradeA, gradeB) {
  return gradeToNumeric(gradeA) > gradeToNumeric(gradeB);
}
```

### Composite Score Formula

```
HealthScore = sum(w_i * S_i)  for i in signals
where sum(w_i) = 1.0 and 0 <= S_i <= 100
```

Each signal S_i is independently normalized to 0–100 before weighting so raw value differences (ticket count vs. NPS rating) do not distort the composite.

### Signal Taxonomy

| Category | Type | Examples |
|----------|------|---------|
| Product usage | Leading | Login frequency, feature breadth, session depth |
| Engagement | Leading | Email opens, webinar attendance, community posts |
| Support behavior | Leading | Ticket velocity, escalation rate, CSAT trend |
| Sentiment | Leading | NPS delta, call sentiment score, survey responses |
| Relationship | Leading | Stakeholder breadth, exec sponsor activity |
| Financial | Lagging | Payment timeliness, contraction signals |
| Outcome | Lagging | Renewal event, cancellation request |

Leading indicators predict; lagging indicators confirm. A score built only on lagging signals detects churn after the decision is made.

---

## Signal Categories & Weights

### Recommended Starting Weights (SaaS B2B)

```javascript
const DEFAULT_WEIGHTS = {
  productUsage:    0.30,  // login freq, feature adoption, depth
  supportHealth:   0.25,  // ticket velocity, severity trend, resolution time
  engagement:      0.20,  // comms interaction, event attendance, community
  sentiment:       0.15,  // NPS, CSAT, call sentiment
  relationship:    0.10,  // stakeholder count, exec sponsor, breadth
};
```

These are starting points. Production weights should be derived from logistic regression or random forest feature-importance against historical churn labels.

### Signal Normalization

```javascript
// 1. Min-Max normalization (simple, sensitive to outliers)
function minMaxNorm(value, min, max) {
  return Math.max(0, Math.min(100, ((value - min) / (max - min)) * 100));
}

// 2. Percentile ranking (robust to outliers, relative)
function percentileNorm(value, sortedPopulation) {
  const rank = sortedPopulation.filter(v => v <= value).length;
  return (rank / sortedPopulation.length) * 100;
}

// 3. Z-score with sigmoid squash (handles Gaussian-distributed signals)
function zScoreNorm(value, mean, stdDev) {
  const z = (value - mean) / stdDev;
  return 100 / (1 + Math.exp(-z));  // sigmoid maps to 0-100
}
```

### Tier-Aware Normalization

A 10-person startup with 5 tickets/month is different from a 500-seat enterprise with 5 tickets/month. Normalize within cohorts:

```javascript
function tierAwareScore(account, signal, cohortStats) {
  const cohort = getCohort(account);
  const { mean, stdDev } = cohortStats[cohort][signal];
  return zScoreNorm(account[signal], mean, stdDev);
}

function getCohort(account) {
  if (account.seats <= 20) return 'SMB';
  if (account.seats <= 200) return 'MidMarket';
  return 'Enterprise';
}
```

---

## Scoring Algorithms

### 1. Linear Weighted Sum (Baseline)

```javascript
function computeHealthScore(account, weights, normalizers) {
  let score = 0;
  for (const [signal, weight] of Object.entries(weights)) {
    const rawValue = account.signals[signal];
    const normalized = normalizers[signal](rawValue);
    score += weight * normalized;
  }
  return Math.round(score);
}
```

- **Pros:** Interpretable, auditable, easy to explain to stakeholders.
- **Cons:** Assumes linear relationship between signals and health.

### 2. Exponential Decay Recency Weighting

Recent signals matter more than stale ones. Half-life guidelines:
- Support tickets: 14–21 days (fast decay)
- Product usage: 30–45 days (moderate)
- NPS/sentiment: 60–90 days (slow decay)
- Relationship signals: 90–120 days (slowest)

```javascript
function recencyWeight(daysSinceEvent, halfLifeDays) {
  const lambda = Math.LN2 / halfLifeDays;
  return Math.exp(-lambda * daysSinceEvent);
}

function timeWeightedSignal(events, halfLifeDays = 30) {
  let weightedSum = 0;
  let totalWeight = 0;
  for (const event of events) {
    const age = daysBetween(event.timestamp, Date.now());
    const w = recencyWeight(age, halfLifeDays);
    weightedSum += w * event.value;
    totalWeight += w;
  }
  return totalWeight > 0 ? weightedSum / totalWeight : 0;
}
```

### 3. Percentile Ranking (Relative Scoring)

```javascript
function percentileRank(accountScore, allAccountScores) {
  const sorted = [...allAccountScores].sort((a, b) => a - b);
  const position = sorted.findIndex(s => s >= accountScore);
  return Math.round((position / sorted.length) * 100);
}
```

Use when you need relative prioritization ("this account is in the bottom 10%") rather than absolute health assessment.

### 4. Multi-Model Ensemble (Production Grade)

```javascript
function ensembleScore(account) {
  const ruleBasedScore = linearWeightedScore(account);
  const mlScore = mlModel.predict(account.features);
  const trendScore = trendAnalysis(account.history);

  const blend = { ruleBased: 0.3, ml: 0.5, trend: 0.2 };

  return Math.round(
    blend.ruleBased * ruleBasedScore +
    blend.ml * mlScore +
    blend.trend * trendScore
  );
}
```

---

## Grading Heuristics

### Letter Grade System (A–F)

```javascript
const GRADE_THRESHOLDS = {
  A: { min: 85, max: 100, label: 'Healthy',   action: 'expand' },
  B: { min: 70, max: 84,  label: 'Stable',    action: 'nurture' },
  C: { min: 55, max: 69,  label: 'At Risk',   action: 'intervene' },
  D: { min: 40, max: 54,  label: 'Declining', action: 'escalate' },
  F: { min: 0,  max: 39,  label: 'Critical',  action: 'rescue' },
};

function assignGrade(score) {
  for (const [grade, { min, max }] of Object.entries(GRADE_THRESHOLDS)) {
    if (score >= min && score <= max) return grade;
  }
  return 'F';
}
```

### Threshold Calibration

Thresholds should not be arbitrary. Calibrate against historical outcomes:

1. Pull all accounts that churned in the last 12 months.
2. Record their scores 90 days before churn.
3. Set the C/D boundary at the 75th percentile of pre-churn scores.
4. Set the D/F boundary at the 50th percentile of pre-churn scores.
5. Validate with precision/recall analysis.

### Color-Code Mapping

```javascript
const GRADE_COLORS = {
  A: { bg: '#dcfce7', text: '#166534', badge: 'green' },
  B: { bg: '#dbeafe', text: '#1e40af', badge: 'blue' },
  C: { bg: '#fef9c3', text: '#854d0e', badge: 'yellow' },
  D: { bg: '#fed7aa', text: '#9a3412', badge: 'orange' },
  F: { bg: '#fecaca', text: '#991b1b', badge: 'red' },
};
```

---

## Time-Series Trending

### Trajectory Detection

A score of 65 that was 80 last month is worse than a score of 55 that was 45 last month.

```javascript
function computeTrajectory(scoreHistory, windowDays = 30) {
  if (scoreHistory.length < 2) return 'insufficient_data';

  const recent = scoreHistory.filter(s => daysBetween(s.date, Date.now()) <= windowDays);
  const prior = scoreHistory.filter(
    s => daysBetween(s.date, Date.now()) > windowDays &&
         daysBetween(s.date, Date.now()) <= windowDays * 2
  );

  if (recent.length === 0 || prior.length === 0) return 'insufficient_data';

  const delta = avg(recent.map(s => s.score)) - avg(prior.map(s => s.score));
  if (delta > 5) return 'improving';
  if (delta < -5) return 'declining';
  return 'stable';
}
```

### Rate of Change (Velocity)

```javascript
function scoreVelocity(scoreHistory) {
  if (scoreHistory.length < 3) return 0;
  const n = scoreHistory.length;
  const points = scoreHistory.map((s, i) => ({ x: i, y: s.score }));
  const xMean = (n - 1) / 2;
  const yMean = avg(points.map(p => p.y));

  let num = 0, den = 0;
  for (const p of points) {
    num += (p.x - xMean) * (p.y - yMean);
    den += (p.x - xMean) ** 2;
  }
  return den !== 0 ? num / den : 0;  // points per period
}
```

### Anomaly Detection (Z-Score)

```javascript
function detectAnomaly(currentScore, scoreHistory, threshold = 2.0) {
  const scores = scoreHistory.map(s => s.score);
  const mean = avg(scores);
  const std = stdDev(scores);
  if (std === 0) return { isAnomaly: false };

  const zScore = (currentScore - mean) / std;
  return {
    isAnomaly: Math.abs(zScore) > threshold,
    zScore,
    direction: zScore < -threshold ? 'sudden_drop' : zScore > threshold ? 'sudden_spike' : 'normal',
    severity: Math.abs(zScore) > 3 ? 'critical' : 'warning',
  };
}
```

### CUSUM (Cumulative Sum) for Shift Detection

```javascript
function cusumDetect(scoreHistory, targetMean, allowedDrift = 5, threshold = 15) {
  let cumSum = 0;
  const alerts = [];

  for (const entry of scoreHistory) {
    const deviation = targetMean - entry.score - allowedDrift;
    cumSum = Math.max(0, cumSum + deviation);

    if (cumSum > threshold) {
      alerts.push({
        date: entry.date,
        message: `Sustained downward shift detected (CUSUM=${cumSum.toFixed(1)})`,
      });
      cumSum = 0;
    }
  }
  return alerts;
}
```

---

## Multi-Dimensional Scoring

Single composite scores obscure root causes. Multi-dimensional scoring separates health into independently actionable dimensions.

### Three-Pillar Model

```javascript
const DIMENSIONS = {
  operational: {
    description: 'Product usage depth, feature adoption, performance',
    signals: ['loginFrequency', 'featureAdoption', 'apiCallVolume', 'errorRate'],
    weight: 0.40,
  },
  relationship: {
    description: 'Human engagement, sentiment, stakeholder breadth',
    signals: ['npsScore', 'meetingFrequency', 'stakeholderCount', 'execSponsor'],
    weight: 0.35,
  },
  adoption: {
    description: 'Expansion signals, seat utilization, workflow coverage',
    signals: ['seatUtilization', 'workflowCoverage', 'newFeatureAdoption', 'trainingCompletion'],
    weight: 0.25,
  },
};

function multiDimensionalScore(account) {
  const dimensionScores = {};
  let compositeScore = 0;

  for (const [dim, config] of Object.entries(DIMENSIONS)) {
    const signals = config.signals.map(s => account.normalizedSignals[s]);
    const dimScore = avg(signals);
    dimensionScores[dim] = Math.round(dimScore);
    compositeScore += config.weight * dimScore;
  }

  return {
    composite: Math.round(compositeScore),
    dimensions: dimensionScores,
    weakestDimension: Object.entries(dimensionScores).sort(([, a], [, b]) => a - b)[0][0],
    grade: assignGrade(Math.round(compositeScore)),
  };
}
```

---

## Anti-Patterns

1. **Scoring without segmentation.** Comparing a 5-seat startup to a 5000-seat enterprise on raw metrics produces meaningless grades. Always normalize within cohorts.
2. **Over-indexing on usage.** High usage with declining sentiment is a false positive. Include at least one sentiment signal.
3. **Static weights.** Weights derived once and never updated drift as the product evolves. Re-calibrate quarterly against actual churn data.
4. **Single score without dimensions.** A composite of 72 tells you nothing about where to act. Always expose sub-scores alongside the composite.
5. **Lagging-only signals.** Scoring on renewal outcomes and payment history detects churn after the decision is made. Prioritize leading indicators.
6. **Ignoring trajectory.** A score of 60 improving at +5/month is healthier than a score of 75 declining at -8/month. Always compute and display velocity.
7. **Binary thresholds without hysteresis.** Toggling between grades on every minor fluctuation creates alert fatigue. Use dead-band thresholds:
   ```javascript
   function hysteresisGrade(currentScore, previousGrade) {
     const upgrade = GRADE_THRESHOLDS[previousGrade].max + 3;
     const downgrade = GRADE_THRESHOLDS[previousGrade].min - 3;
     if (currentScore > upgrade) return gradeAbove(previousGrade);
     if (currentScore < downgrade) return gradeBelow(previousGrade);
     return previousGrade;
   }
   ```
8. **No implementation window.** Predicting churn 7 days out gives no time to act. Design for 60–90 day lead time by using features from an observation window separated from the prediction target.
9. **Treating all accounts equally.** Enterprise accounts need relationship signals weighted higher; SMB accounts need usage signals weighted higher. Use tier-specific weight profiles.
10. **Invisible scoring logic.** If CSMs cannot explain why an account scored 52, they will not trust or act on the score. Provide signal-level breakdowns with every grade.

---

## Full Health Score Engine

```javascript
class HealthScoreEngine {
  constructor(config) {
    this.weights = config.weights;
    this.thresholds = config.thresholds;
    this.decayRates = config.decayRates;
    this.cohortStats = null;
  }

  async initialize(accountPopulation) {
    this.cohortStats = this.computeCohortStats(accountPopulation);
  }

  computeCohortStats(accounts) {
    const cohorts = { SMB: {}, MidMarket: {}, Enterprise: {} };
    for (const cohort of Object.keys(cohorts)) {
      const members = accounts.filter(a => getCohort(a) === cohort);
      for (const signal of Object.keys(this.weights)) {
        const values = members.map(a => a.signals[signal]).filter(v => v != null);
        cohorts[cohort][signal] = {
          mean: avg(values),
          stdDev: stdDev(values),
          sorted: [...values].sort((a, b) => a - b),
        };
      }
    }
    return cohorts;
  }

  scoreAccount(account) {
    const cohort = getCohort(account);
    const normalized = {};
    let composite = 0;

    for (const [signal, weight] of Object.entries(this.weights)) {
      const raw = account.signals[signal];
      const stats = this.cohortStats[cohort][signal];
      normalized[signal] = zScoreNorm(raw, stats.mean, stats.stdDev);
      composite += weight * normalized[signal];
    }

    const score = Math.round(Math.max(0, Math.min(100, composite)));
    return {
      accountId: account.id,
      score,
      grade: assignGrade(score),
      trajectory: computeTrajectory(account.scoreHistory),
      anomaly: detectAnomaly(score, account.scoreHistory),
      signals: normalized,
      weakestSignal: Object.entries(normalized).sort(([, a], [, b]) => a - b)[0][0],
      computedAt: new Date().toISOString(),
    };
  }

  batchScore(accounts) {
    return accounts.map(a => this.scoreAccount(a));
  }

  detectAtRisk(accounts, options = {}) {
    const { maxGrade = 'C' } = options;
    const scored = this.batchScore(accounts);
    return scored.filter(s =>
      isWorseGrade(s.grade, maxGrade) || s.grade === maxGrade ||
      s.trajectory === 'declining' ||
      (s.anomaly?.isAnomaly && s.anomaly.direction === 'sudden_drop')
    );
  }
}
```

## Score History Storage Schema

```javascript
// MongoDB/document store schema
const scoreHistorySchema = {
  accountId: String,
  score: Number,         // 0-100
  grade: String,         // A-F
  dimensions: { operational: Number, relationship: Number, adoption: Number },
  signals: Object,       // { signalName: normalizedValue }
  trajectory: String,    // improving | stable | declining
  velocity: Number,
  anomaly: { isAnomaly: Boolean, zScore: Number, direction: String },
  metadata: { cohort: String, tier: String, computedAt: Date, engineVersion: String },
};
// Index: db.healthScores.createIndex({ accountId: 1, "metadata.computedAt": -1 })
```

## REST API Surface

```
GET  /api/accounts/:id/health              — current score with breakdown
GET  /api/accounts/:id/health/history?days=90 — time-series score history
GET  /api/accounts/:id/health/anomalies?days=30 — detected anomalies
GET  /api/health/at-risk?grade=D&trajectory=declining — filtered at-risk list
POST /api/health/recalculate               — batch recalculation for all accounts
```
