<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `playbook-matcher-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: playbook-matcher-patterns
version: 2.1.0
last_updated: 2026-05-29
description: >
  Rule-based matching of support cases to diagnostic playbooks and KB articles — keyword scoring,
  confidence tiers, haystack construction, synonym expansion, TF-IDF-lite ranking, decision table
  design, multi-signal fusion, fuzzy matching, external rule engine integration, and KB article
  recommendation. Grounded in the production playbook-matcher.js and kb-index.js from
  mdb-case-assistant.
  TRIGGER: building or extending rule-based case-to-playbook matching, adding new playbook patterns
  or KB articles to static indexes, implementing keyword scoring or synonym expansion, designing
  diagnostic decision trees, reviewing why a case matched or failed to match a pattern, adding
  fuzzy matching for typos or abbreviations, evaluating rule engine library vs hand-rolled matcher,
  building decision tables or routing rules for support case classification.
  SKIP: LLM prompt construction or Glean integration (use glean-llm-client-patterns);
  TS Tools API data fetching (use ts-tools-support-api); DOM scraping or case field extraction
  (use dom-scraping-resilience); overlay/panel UI rendering (use vanilla-js-ui-reviewer);
  general MongoDB troubleshooting knowledge (use mongodb-kb or atlas-diagnostics-expert);
  full-text search with inverted indexes at scale (use a dedicated search engine).
category: developer
tags:
  - playbook-matching
  - keyword-scoring
  - case-triage
  - decision-table
  - tfidf
  - fuzzy-matching
  - rule-engine
  - kb-articles
  - diagnostic-playbook
  - confidence-scoring
globs:
  - src/background/playbook-matcher.js
  - src/background/kb-index.js
  - src/background/case-enricher.js
triggers:
  - playbook matching
  - case-to-playbook mapping
  - KB article recommendation
  - keyword scoring
  - diagnostic playbook
  - symptom-based triage
  - case routing rules
  - matchPlaybookPatterns
  - searchKbArticles
  - rule engine pattern
  - decision table matching
  - support case classification
  - weighted keyword match
  - confidence scoring
  - case triage automation
  - PLAYBOOK_PATTERNS
  - KB_ARTICLES
  - buildDiagnosticSuggestion
  - getTopMatch
  - TOKEN_SYNONYMS
  - deriveKeywords
  - TF-IDF scoring
  - json-rules-engine
  - fuzzy matching
related_skills:
  - diagnostic-registry-patterns
  - atlas-diagnostics-expert
  - atlas-diagnostics-package
  - case-tracker
  - mongodb-kb
  - ops-registry-patterns
  - glean-llm-client-patterns
---

# Playbook Matcher Patterns

Rule-based matching of support cases to diagnostic playbooks and KB articles. All core patterns are grounded in the production `src/background/playbook-matcher.js` and `src/background/kb-index.js` from the mdb-case-assistant Chrome extension.

## When to use this skill

- Building or extending rule-based case-to-playbook matching
- Adding new playbook patterns or KB articles to the static indexes
- Implementing keyword scoring, synonym expansion, or weighted search
- Designing diagnostic decision trees for guided operator triage
- Reviewing or debugging why a case matched (or failed to match) a pattern
- Adding fuzzy matching to improve recall on typos or abbreviations
- Evaluating rule engine library vs. hand-rolled matcher
- Building decision tables or routing rules for support case classification

## When NOT to use this skill

- LLM prompt construction or Glean integration → use `glean-llm-client-patterns`
- TS Tools API data fetching → use `ts-tools-support-api` or `tstools-reference`
- DOM scraping or case field extraction → use `dom-scraping-resilience`
- Overlay/panel UI rendering of match results → use `vanilla-js-ui-reviewer`
- General MongoDB troubleshooting knowledge → use `mongodb-kb` or `atlas-diagnostics-expert`
- Full-text search at scale (500+ articles) → use MeiliSearch, Typesense, or Elasticsearch

---

## 1. Core architecture: extract, score, rank

```
Case record (subject, description, errorStrings, status)
  |
  +-- Stage 1: Build haystack — flatten text fields into a single searchable string
  +-- Stage 2: Score — evaluate every pattern/article against the haystack
  +-- Stage 3: Rank — sort by match count or weighted score, apply confidence tiers
```

### Production matcher

```js
export function matchPlaybookPatterns(caseRecord = {}) {
  const snapshot = caseRecord.snapshot || caseRecord;
  const haystack = [
    snapshot.subject, snapshot.description,
    ...(snapshot.errorStrings || []), snapshot.status
  ].filter(Boolean).join(' ').toLowerCase();

  if (!haystack.trim()) return [];

  return PLAYBOOK_PATTERNS
    .map(pattern => {
      const matched = pattern.keywords.filter(kw => haystack.includes(kw.toLowerCase()));
      return {
        ...pattern, matchedKeywords: matched, matchCount: matched.length,
        confidence: matched.length >= 3 ? 'high' : matched.length >= 1 ? 'medium' : 'none'
      };
    })
    .filter(p => p.matchCount > 0)
    .sort((a, b) => b.matchCount - a.matchCount);
}
```

**Key design decisions:**
- Haystack normalization: lowercase only — no stemming or lemmatization. Keeps matching deterministic and fast.
- Keyword matching is substring-based (`String.includes`), not token-boundary-aware. `"auth"` matches `"authentication"` and `"unauthorized"`.
- Fixed confidence thresholds: 3+ keywords = high, 1-2 = medium.
- No external dependencies — entire engine runs synchronously in a service worker.

**Haystack construction rules:**
- Always include `subject` and `description` — highest signal
- Include extracted error strings separately — often contain driver codes not in the prose
- Include `status` or `category` when the source system provides them
- Do NOT include case ID, timestamps, or agent names — noise only
- Call `.filter(Boolean)` before `.join()` to strip null/undefined

---

## 2. Pattern registry design

```js
{
  id: 'connection-timeout',
  label: 'Connection Timeout / Network Issue',
  keywords: [
    'connection timeout', 'timed out', 'ECONNREFUSED',
    'connection refused', 'network error', 'socket hang up',
    'ETIMEDOUT', 'connection reset'
  ],
  runbook: 'Check firewall rules, IP allowlisting, DNS resolution...',
  playbookSection: 'Connection Troubleshooting',
  severity: 'varies',
  diagnosticToolIds: ['ts-diag']
}
```

### Guidelines for new patterns

1. **Keep keywords specific.** Prefer `"ECONNREFUSED"` over `"error"`. Generic terms cause false positives.
2. **Mix phrase keywords with error codes.** Multi-word phrases give precision; driver error codes catch structured error output.
3. **Include MongoDB-specific identifiers.** `"COLLSCAN"`, `"StaleConfig"`, `"TransactionTooLargeForCache"` are high-signal.
4. **Link diagnostic tools.** `diagnosticToolIds` connects patterns to runnable tools for one-click triage.
5. **Use severity ranges.** `"S1"`, `"S2-S4"` sets operator expectations.
6. **Aim for 5-10 keywords per pattern.** Fewer than 3 makes `high` confidence unreachable; more than 15 dilutes signal.

### Confidence tier calibration

| Tier | Threshold | Meaning |
|------|-----------|---------|
| `high` | 3+ keyword hits | Strong signal — safe for auto-suggestion |
| `medium` | 1-2 keyword hits | Plausible — show as secondary recommendation |
| `none` | 0 hits | No match — filter out |

For patterns with only 3-4 total keywords, lower `high` to 2. For patterns with 15+ keywords, use percentage thresholds (e.g., 10% = medium, 30% = high).

---

## 3. KB article recommendation engine

### Keyword derivation pipeline

```js
function deriveKeywords(entry) {
  const baseTokens = new Set([
    ...tokenize(entry.title), ...tokenize(entry.summary),
    ...tokenize(entry.category),
    ...entry.products.flatMap(p => tokenize(p)),
    entry.articleId,
  ]);

  // Category-level keyword injection
  for (const kw of CATEGORY_KEYWORDS[entry.category] || []) {
    baseTokens.add(normalizeText(kw));
  }

  // Synonym expansion at index time
  for (const token of Array.from(baseTokens)) {
    for (const syn of TOKEN_SYNONYMS[token] || []) {
      baseTokens.add(normalizeText(syn));
    }
  }

  return Array.from(baseTokens).filter(Boolean).sort();
}
```

### Category keyword injection

```js
const CATEGORY_KEYWORDS = {
  'Atlas':                    ['atlas', 'cluster', 'ops', 'deployment'],
  'Connectivity & Timeouts':  ['connectivity', 'timeouts', 'driver', 'dns', 'srv', 'tls'],
  'Replication':              ['replication', 'oplog', 'secondary', 'failover', 'election'],
  'Performance & Indexes':    ['performance', 'indexes', 'slow query', 'memory', 'storage'],
  'Security & Auth':          ['security', 'authentication', 'authorization', 'tls', 'ldap', 'oidc'],
  'Networking':               ['networking', 'dns', 'peering', 'privatelink', 'ip access list']
};
```

When adding a new KB category, add a corresponding entry here or articles in that category will miss category-level search hits.

### Synonym expansion

```js
const TOKEN_SYNONYMS = {
  auth:    ['authentication', 'authorization', 'login', 'access'],
  tls:     ['ssl', 'certificate', 'x509'],
  timeout: ['timeouts', 'sockettimeout', 'brokenpipe'],
  network: ['networking', 'latency', 'privatelink', 'peering'],
  performance: ['slow', 'latency', 'throughput'],
  dns:     ['srv', 'txt', 'seedlist']
};
```

Apply expansion at **index time** (when deriving article keywords), not at query time. Index-time expansion means the search text is wider, so simple `includes()` checks on the query side still match. For bidirectional synonyms (`tls` ↔ `ssl`), ensure both directions are covered.

### Weighted search scoring

```js
const score = queryTokens.reduce((acc, token) =>
  acc
  + (entry.title.toLowerCase().includes(token)   ? 3 : 0)
  + (entry.keywords.includes(token)              ? 2 : 0)
  + (entry.summary.toLowerCase().includes(token) ? 1 : 0),
  0
);
```

| Field | Weight | Rationale |
|-------|--------|-----------|
| Title | 3x | Highest editorial signal; concise and intentional |
| Keywords (derived) | 2x | Includes synonyms and category terms |
| Summary | 1x | Broader context but noisier |

---

## 4. Multi-signal fusion

```js
export function buildDiagnosticSuggestion(caseRecord = {}) {
  const matches = matchPlaybookPatterns(caseRecord);
  if (matches.length === 0) return null;

  const top = matches[0];
  return {
    pattern: top.label, confidence: top.confidence,
    matchedKeywords: top.matchedKeywords, runbook: top.runbook,
    playbookSection: top.playbookSection,
    additionalMatches: matches.slice(1, 3).map(m => ({ pattern: m.label, confidence: m.confidence }))
  };
}
```

### Full triage composition with KB search

```js
function triageCase(caseRecord) {
  const playbookMatches = matchPlaybookPatterns(caseRecord);
  const topPlaybook = playbookMatches[0] || null;
  const snapshot = caseRecord.snapshot || caseRecord;

  // Feed matched keywords into KB search for contextual articles
  const query = topPlaybook
    ? topPlaybook.matchedKeywords.join(' ')
    : [snapshot.subject, ...(snapshot.errorStrings || [])].filter(Boolean).join(' ');

  const kbResults = searchKbArticles(query, { products: snapshot.products || [], limit: 5 });

  return {
    playbook: topPlaybook ? {
      pattern: topPlaybook.label, confidence: topPlaybook.confidence,
      runbook: topPlaybook.runbook, section: topPlaybook.playbookSection,
      diagnosticTools: topPlaybook.diagnosticToolIds
    } : null,
    kbArticles: kbResults.map(a => ({ articleId: a.articleId, title: a.title, url: a.url, visibility: a.visibility })),
    additionalPatterns: playbookMatches.slice(1, 4).map(m => ({ pattern: m.label, confidence: m.confidence }))
  };
}
```

**Fusion rules:**
- Lead with the top playbook — it provides runbook/action guidance
- Feed matched keywords from the top playbook into KB search to avoid re-parsing case text
- Cap secondary playbooks at 2-3 to avoid overwhelming the operator
- When no playbook matches, fall back to KB-only search using case subject + errors

---

## 5. Diagnostic decision trees

Use when a single keyword set cannot distinguish sub-problems.

```js
const DECISION_TREE = {
  id: 'connection-issue',
  children: [
    {
      id: 'dns-failure',
      matchKeywords: ['NXDOMAIN', 'dns', 'srv', 'getaddrinfo', 'ENOTFOUND'],
      children: [
        { id: 'srv-misconfigured', matchKeywords: ['srv', 'txt record'], leaf: true,
          recommendation: 'Validate SRV/TXT DNS records with dig/nslookup.' },
        { id: 'standard-dns', matchKeywords: ['direct', 'standard'], leaf: true,
          recommendation: 'Check /etc/resolv.conf and VPC DNS settings.' }
      ]
    },
    {
      id: 'tcp-timeout',
      matchKeywords: ['ETIMEDOUT', 'connection timeout', 'ECONNREFUSED'],
      children: [
        { id: 'private-endpoint', matchKeywords: ['privatelink', 'vpc endpoint'], leaf: true,
          recommendation: 'Verify security groups, endpoint status, and DNS.' },
        { id: 'public-internet', matchKeywords: ['public', 'ip access'], leaf: true,
          recommendation: 'Check Atlas IP access list and firewall rules.' }
      ]
    }
  ]
};

const MAX_TREE_DEPTH = 10;

function walkDecisionTree(node, haystack, path = [], depth = 0) {
  path.push({ id: node.id, question: node.question || null });
  if (node.leaf) return { path, recommendation: node.recommendation };
  if (!node.children || depth >= MAX_TREE_DEPTH) {
    return { path, recommendation: null, ...(depth >= MAX_TREE_DEPTH && { truncated: true }) };
  }
  const scored = node.children
    .map(child => ({ child, score: (child.matchKeywords || []).filter(kw => haystack.includes(kw.toLowerCase())).length }))
    .sort((a, b) => b.score - a.score);
  if (scored[0].score === 0) return { path, recommendation: null, ambiguous: true };
  return walkDecisionTree(scored[0].child, haystack, path, depth + 1);
}
```

**Always set `MAX_TREE_DEPTH`.** Malformed or circular definitions cause unbounded recursion without it.

### When to use trees vs flat patterns

| Criterion | Flat pattern match | Decision tree |
|---|---|---|
| Problem categories are distinct | Preferred | Overkill |
| Sub-categories share keywords | Produces ambiguous results | Preferred |
| Operator needs guided workflow | Returns a single answer | Walks through questions |

---

## 6. Decision table pattern

```js
const ROUTING_TABLE = [
  { severity: 'S1', hasErrorCode: true,  category: 'outage', action: 'route-to-critical' },
  { severity: 'S2', hasErrorCode: true,  category: 'perf',   action: 'route-to-perf-team' },
  { severity: '*',  hasErrorCode: true,  category: '*',      action: 'attach-error-playbook' },
  { severity: '*',  hasErrorCode: '*',   category: '*',      action: 'route-to-general' }
];

function evaluateDecisionTable(caseContext, table) {
  for (const rule of table) {
    const match = Object.entries(rule)
      .filter(([key]) => key !== 'action')
      .every(([key, expected]) => expected === '*' || caseContext[key] === expected);
    if (match) return rule.action;
  }
  return null;
}
```

**Decision table design rules:**
- Order rows from most specific to least specific — first match wins
- Use `'*'` as a wildcard for "don't care" conditions
- Keep the table declarative and separate from evaluation logic
- Decision tables are easy to serialize to JSON and store externally — editable without code changes
- Ideal for 2-5 conditions. Beyond 5, use a decision tree instead.

---

## 7. Advanced scoring (upgrade paths)

When simple keyword counting produces too many ties or you need fuzzy matching for typos/abbreviations, see the full implementations in [references/advanced-scoring.md](./references/advanced-scoring.md):

- **TF-IDF-lite:** `computeIdf()` + `tfidfScore()` — upgrade when KB index exceeds ~200 articles or common terms dominate scores
- **Levenshtein fuzzy matching:** `fuzzyMatch()` with 0.5x weight for approximate hits — always guard with minimum keyword length check (skip keywords ≤ 3 chars)
- **Bigram similarity:** `bigramSimilarity()` for multi-word phrase matching
- **Text normalization / stop words:** `normalizeText()` + `tokenize()` — keep stop word lists minimal; domain terms like `"can"` may matter in error messages

---

## 8. External rule engine (json-rules-engine)

For projects that outgrow hand-rolled matchers — rules editable by non-developers, stored in a database, or needing AND/OR/NOT combinators. Full integration example in [references/advanced-scoring.md](./references/advanced-scoring.md).

**Adopt when:** rules change frequently without code deployments; need declarative AND/OR/NOT logic.
**Stay hand-rolled when:** < 20 patterns maintained by developers; zero runtime dependencies needed (Chrome extension, edge function); synchronous evaluation required (json-rules-engine is async).

---

## 9. Scaling guidance

| Corpus size | Recommended approach |
|---|---|
| < 50 playbooks, < 100 KB articles | Static arrays, `includes()` matching, field-weighted scoring |
| 50-500 articles | Add TF-IDF-lite scoring and precomputed search text |
| 500-5000 articles | Build an inverted index; use `Map<token, Set<articleId>>` |
| 5000+ articles | Move to a dedicated search engine (MeiliSearch, Typesense, Elasticsearch) |

The current mdb-case-assistant sits in the first tier (~8 playbooks, ~70 KB articles). The static-array approach is correct at this scale.

---

## 10. Testing patterns

```js
// Unit-test a pattern match
test('connection-timeout matches ECONNREFUSED', () => {
  const result = matchPlaybookPatterns({
    subject: 'Cannot connect to Atlas cluster',
    description: 'Getting ECONNREFUSED when connecting from Lambda',
    errorStrings: ['ECONNREFUSED', 'connection refused']
  });
  expect(result[0].id).toBe('connection-timeout');
  expect(result[0].confidence).toBe('high');
  expect(result[0].matchedKeywords).toContain('connection refused');
});

// Regression suite with anonymized case snapshots
const REGRESSION_CASES = [
  {
    name: 'S1 cluster outage',
    input: { subject: 'Production cluster down',
             description: 'Primary unavailable, election failing, heartbeat timeouts' },
    expectedTopId: 'cluster-outage', expectedConfidence: 'high'
  },
  {
    name: 'auth failure with LDAP',
    input: { subject: 'Users cannot log in',
             description: 'LDAP auth failed, access denied for all users' },
    expectedTopId: 'auth-failure', expectedConfidence: 'high'
  }
];

describe('regression suite', () => {
  for (const tc of REGRESSION_CASES) {
    it(tc.name, () => {
      const top = getTopMatch(tc.input);
      expect(top.id).toBe(tc.expectedTopId);
      expect(top.confidence).toBe(tc.expectedConfidence);
    });
  }
});
```

**Edge cases to cover:**
- Empty case record returns `[]`
- Case with no matching keywords returns `[]`
- Multiple patterns can match the same case (multi-match)
- Confidence tiers transition at exactly 1 and 3 keyword matches
- KB search with stop-word-only query returns broad results, not empty

---

## 11. Anti-patterns

1. **Over-broad keywords.** Adding `"error"` or `"issue"` matches nearly every case. Use specific error codes or MongoDB terminology.
2. **Hardcoded thresholds without context.** The 1/3 thresholds work for 5-10 keyword patterns. For 20+ keywords, recalibrate to percentages.
3. **Mutating the frozen index at runtime.** Both `PLAYBOOK_PATTERNS` and `KB_ARTICLES` are exported as frozen arrays. Copy before modifying.
4. **Skipping synonym expansion.** Without synonyms, `"ssl"` and `"tls"` are separate terms. The `TOKEN_SYNONYMS` map bridges these.
5. **Ignoring visibility filtering.** KB articles have `Public` and `Internal` visibility. Always pass the `visibility` option when the consumer context restricts access.
6. **Fuzzy matching very short keywords.** Running Levenshtein on 1-2 character keywords with `maxDistance = 2` matches almost any word. Guard with a minimum keyword length check.
7. **Stale keyword lists.** Lists not updated when new products, error codes, or driver versions ship cause recall to degrade silently. Schedule quarterly reviews.
8. **Missing haystack fields.** If the haystack omits a field carrying matching signal (e.g., `errorStrings`), the matcher will never fire for cases where the keyword appears only in that field.
9. **No fallback for zero matches.** Degrade gracefully — suggest a general triage playbook, offer a KB search, or prompt the operator for more context.

---

## 12. Checklist: adding a new playbook pattern

1. Choose a unique kebab-case `id` that does not collide with existing patterns
2. Write 5-10 keywords mixing natural-language phrases and error codes
3. Verify keywords are lowercase and do not overlap heavily with other patterns
4. Add a `runbook` string with actionable first steps — not just "investigate"
5. Set `severity` to the expected range for this pattern
6. Link relevant `diagnosticToolIds` from the diagnostic registry
7. Write a unit test with a case record that should produce a `high` confidence match
8. Write a negative test with a case record that should NOT match this pattern
9. Run the full regression suite to ensure no existing mappings broke
10. Verify that relevant KB articles appear in `searchKbArticles` results when queried with the pattern's keywords

---

## 13. Source files

- `src/background/playbook-matcher.js` — production pattern registry and matcher
- `src/background/kb-index.js` — static KB article index with weighted search
- `src/background/case-enricher.js` — normalizes DOM/API data before matching
- `src/background/service-worker.js` — orchestrates triage via `MCA_*` messages
- `src/background/llm-client.js` — consumes matcher output to build LLM/Glean prompts
