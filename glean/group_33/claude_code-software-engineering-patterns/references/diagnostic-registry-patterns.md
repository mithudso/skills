<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Formerly the standalone `diagnostic-registry-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: diagnostic-registry-patterns
version: 1.0.0
description: >
  Diagnostic tool registry design -- static registries with keyword-based symptom matching,
  URL template resolution, category filtering, recommendation scoring, conditional execution,
  LLM prompt catalog generation, and case-driven suggestion normalization. Grounded in the
  atlas-diagnostics package from mdb-case-assistant.
triggers:
  - diagnostic tool registry
  - symptom mapping
  - tool recommendation
  - whenToUse keyword matching
  - diagnostic category filter
  - URL template resolution
  - case-driven tool suggestion
  - recommendation scoring
  - tool catalog for LLM
  - diagnostic deep-link
  - getHighlightedDiagnosticTools
  - normalizeSuggestedDiagnosticTools
  - buildDiagnosticToolsCatalogForPrompt
related-skills:
  - atlas-diagnostics-expert
  - atlas-diagnostics-package
  - case-tracker
  - ops-registry-patterns
---

# Diagnostic Registry Patterns

Expert reference for building diagnostic tool registries that map symptoms to tools,
resolve context-aware URLs, filter by category, score recommendations, and generate
LLM-ready catalogs. All patterns are grounded in the real `packages/atlas-diagnostics/src/index.js`
implementation from the mdb-case-assistant Chrome extension.

**Source file:** `packages/atlas-diagnostics/src/index.js` (re-exported via `src/background/diagnostic-registry.js`)
**Test file:** `packages/atlas-diagnostics/tests/index.test.js`

## When to Use This Skill

**Use when:**
- Building or extending a diagnostic tool registry that maps case symptoms to recommended tools
- Adding keyword-based matching, category filtering, or URL template resolution to a registry
- Generating LLM prompt catalogs from a tool registry
- Normalizing tool suggestions returned by an LLM or external API
- Reviewing code that touches `TOOL_REGISTRY`, `getHighlightedDiagnosticTools`, or `normalizeSuggestedDiagnosticTools`

**Do NOT use when:**
- The task is about Atlas diagnostic workflows or FTDC/log analysis (use `atlas-diagnostics-expert` instead)
- The task is about the Chrome extension overlay or panel UI that consumes registry data (use `atlas-diagnostics-package` instead)
- The task is general registry/service-locator design without a diagnostic/symptom-mapping component (use `ops-registry-patterns` instead)

## Public API Summary

| Export | Purpose |
|--------|---------|
| `getDiagnosticToolsRegistry()` | Returns a deep clone of the full `TOOL_REGISTRY` object |
| `getDisplayCategories()` | Returns the ordered category list (defensive copy) |
| `getHighlightedDiagnosticTools(caseContext)` | Keyword-matches case text against every tool's `whenToUse` |
| `getDiagnosticToolsForCategory(category, caseContext)` | Filters highlighted tools to a single display category |
| `resolveToolUrl(toolId, variables)` | Interpolates a tool's URL template or returns fallback |
| `normalizeSuggestedDiagnosticTools(items)` | Validates, enriches, and caps LLM/external tool suggestions |
| `buildDiagnosticToolsCatalogForPrompt()` | Serializes the registry into a compact LLM prompt catalog |
| `buildTsDiagUrl(caseContext)` | Constructs a ts-diag deep-link from case identifiers |

## 1. Static Tool Registry Design

A diagnostic registry is a flat object keyed by stable tool IDs. Each entry carries
enough metadata for display, matching, linking, and operator guidance.

### Registry Entry Schema

Every tool entry MUST include these fields:

| Field              | Type       | Purpose                                              |
|--------------------|------------|------------------------------------------------------|
| `name`             | string     | Human-readable display name                          |
| `sourceCategory`   | string     | Canonical owner category (single value)              |
| `displayCategories`| string[]   | All categories where the tool appears in filtered UI |
| `guideUrl`         | string     | Documentation or how-to link                         |
| `urlTemplate`      | string     | Parameterized deep-link with `{placeholder}` tokens  |
| `urlFallback`      | string     | Safe landing page when variables are missing         |
| `description`      | string     | One-sentence purpose statement                       |
| `whenToUse`        | string     | Comma-separated symptom keywords for matching        |
| `accessGuidance`   | string[]   | Ordered setup/auth prerequisites                     |
| `steps`            | string[]   | Ordered operator workflow steps                      |

### Real Example -- Registry Entry

```js
const TOOL_REGISTRY = {
  atlas_realtime_performance: {
    name: 'Atlas Real-Time Performance Panel',
    sourceCategory: 'Atlas Metrics',
    displayCategories: ['Atlas Metrics'],
    guideUrl: 'https://www.mongodb.com/docs/atlas/monitor-cluster-metrics/',
    urlTemplate: 'https://cloud.mongodb.com/v2/{project_id}#/clusters/{cluster_name}/metrics',
    urlFallback: 'https://cloud.mongodb.com/',
    description: 'Live view of active operations, connections, opcounters, CPU, memory, and network.',
    whenToUse: 'slow queries, high CPU, latency spike, timeout, connection saturation, performance degradation, ops/sec',
    accessGuidance: [
      'No local install is required; open the Atlas web UI in a browser.',
      'Use an Atlas role that can view the target project and cluster metrics.'
    ],
    steps: [
      'Log in to Atlas at https://cloud.mongodb.com',
      'Select the organization and project that owns the affected cluster',
      'Open the cluster Metrics page and switch to Real-Time',
      'Compare active connections, opcounters, CPU, and network spikes to the incident window'
    ]
  }
};
```

### Design Rules

1. **Stable keys.** Tool IDs are snake_case and never change once shipped. UI and LLM integrations reference them by ID.
2. **Flat registry, no nesting.** Every tool is a top-level entry. Categories are metadata on the entry, not a grouping hierarchy.
3. **Immutable at runtime.** Expose the registry through `structuredClone()` or `.slice()` to prevent mutation by callers.
4. **Single source of truth.** The registry object is the only place tool metadata lives. UI, matching, URL resolution, and prompt generation all read from it.

```js
// Defensive copy -- callers cannot mutate the canonical registry
export function getDiagnosticToolsRegistry() {
  return structuredClone(TOOL_REGISTRY);
}
```

## 2. Symptom-to-Tool Keyword Matching

The core matching engine maps free-text case context (subject, description, error strings)
to relevant tools using the `whenToUse` keyword phrases on each registry entry.

### Matching Algorithm

1. **Build haystack.** Concatenate case subject, description, error strings, and status into a single normalized string.
2. **Tokenize keywords.** Split each tool's `whenToUse` on commas into individual keyword phrases.
3. **Match each phrase.** Test each normalized keyword against the haystack using word-boundary regex with an `includes()` fallback for multi-word phrases.
4. **Tag results.** Tools with one or more matches are marked `highlighted: true` with a `whyRelevant` explanation string.

```js
function normalizeText(value = '') {
  return String(value || '')
    .toLowerCase()
    .replace(/[^a-z0-9$]+/g, ' ')
    .trim();
}

function escapeRegExp(value = '') {
  return String(value || '').replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

function tokenizeKeywords(whenToUse = '') {
  return String(whenToUse || '')
    .split(',')
    .map(item => item.trim())
    .filter(Boolean);
}

export function getHighlightedDiagnosticTools(caseContext = {}) {
  const ctx = caseContext || {};
  const haystack = normalizeText([
    ctx.subject,
    ctx.description,
    ...(ctx.errorStrings || []),
    ctx.status
  ].join(' '));

  return Object.entries(TOOL_REGISTRY).map(([toolId, tool]) => {
    const matches = tokenizeKeywords(tool.whenToUse).filter(keyword => {
      const normalizedKeyword = normalizeText(keyword);
      if (!normalizedKeyword) return false;
      return new RegExp(`\\b${escapeRegExp(normalizedKeyword)}\\b`, 'i').test(haystack)
        || haystack.includes(normalizedKeyword);
    });
    return {
      toolId,
      ...tool,
      highlighted: matches.length > 0,
      whyRelevant: matches.length
        ? `Matches: ${matches.map(m => `'${m}'`).join(', ')}`
        : ''
    };
  });
}
```

### Matching Best Practices

- **Normalize aggressively.** Strip all non-alphanumeric characters except `$` (for MongoDB operators like `$search`). This avoids false negatives from punctuation or casing.
- **Escape user-facing regex.** Always run `escapeRegExp()` on keyword phrases before building a `RegExp` to prevent injection from characters like `(`, `)`, `*`.
- **Dual-strategy matching.** Word-boundary regex catches exact tokens; `includes()` catches compound phrases that span word boundaries.
- **Return all tools, not just matches.** UI consumers need the full list to render "highlighted vs dimmed" views. Filtering to matches-only is a caller decision.

## 3. Recommendation Scoring (Extension Pattern)

> **Note:** This section describes an extension pattern not present in the current source.
> The shipped implementation uses binary highlighting via `getHighlightedDiagnosticTools`.
> Use this pattern when the UI needs ranked/sorted recommendations rather than highlight-only.

The base matching engine produces a binary highlighted/not-highlighted signal. For
ranked recommendations, layer a scoring function over the match results.

### Score Calculation Pattern

```js
function scoreTool(tool, caseContext) {
  const haystack = normalizeText([
    caseContext.subject,
    caseContext.description,
    ...(caseContext.errorStrings || [])
  ].join(' '));

  const keywords = tokenizeKeywords(tool.whenToUse);
  let score = 0;

  for (const keyword of keywords) {
    const normalized = normalizeText(keyword);
    if (!normalized) continue;

    // Exact word-boundary match scores higher
    if (new RegExp(`\\b${escapeRegExp(normalized)}\\b`, 'i').test(haystack)) {
      score += 2;
    } else if (haystack.includes(normalized)) {
      score += 1;
    }
  }

  // Boost tools whose sourceCategory matches case metadata
  if (caseContext.category && tool.sourceCategory === caseContext.category) {
    score += 3;
  }

  return score;
}

function getRankedTools(caseContext) {
  return Object.entries(TOOL_REGISTRY)
    .map(([toolId, tool]) => ({ toolId, ...tool, score: scoreTool(tool, caseContext) }))
    .filter(t => t.score > 0)
    .sort((a, b) => b.score - a.score);
}
```

### Scoring Guidelines

- **Weight exact boundary matches higher than substring matches.** A boundary match on "slow queries" is a stronger signal than a substring hit on "queries" inside "query-optimizer".
- **Boost on category alignment.** If the case already has a known category, tools in that category get a static bonus.
- **Cap returned results.** Use `.slice(0, N)` to prevent overwhelming the operator. The real implementation caps at 10 via `normalizeSuggestedDiagnosticTools`.
- **Never score on absence.** Only positive signals count. Negative scoring (penalizing tools that do not match) leads to brittle ranking.

## 4. Category Filtering

The registry supports a fixed category order for consistent UI rendering and a filter
function that intersects category membership with symptom matching.

### Category Order

```js
const CATEGORY_ORDER = [
  'Atlas Metrics', 'Profiler', 'Logs', 'Connection',
  'Replication', 'Storage', 'Security'
];

export function getDisplayCategories() {
  return CATEGORY_ORDER.slice();  // defensive copy
}
```

### Filtered Query

```js
export function getDiagnosticToolsForCategory(category = '', caseContext = {}) {
  const normalizedCategory = String(category || '').trim();
  const tools = getHighlightedDiagnosticTools(caseContext || {});
  if (!normalizedCategory || normalizedCategory === 'All') return tools;
  return tools.filter(tool => tool.displayCategories.includes(normalizedCategory));
}
```

### Category Design Rules

- **Tools appear in multiple categories.** A single tool can belong to several `displayCategories` (e.g., `atlas_cluster_metrics` appears under Atlas Metrics, Replication, and Storage).
- **sourceCategory is singular.** It represents canonical ownership. `displayCategories` represents visibility.
- **"All" is the default.** When no filter is active, return every tool with highlight status intact.
- **Preserve category order.** The UI renders tabs/filters in `CATEGORY_ORDER` sequence. Never sort categories alphabetically.

## 5. URL Template Resolution and Conditional Execution

Tools carry parameterized URL templates with `{placeholder}` tokens. Resolution
interpolates case variables and falls back gracefully when data is incomplete.

### Template Interpolation

```js
function interpolateTemplate(template = '', variables = {}) {
  return String(template || '').replace(/\{([^}]+)\}/g, (_, key) => {
    const value = String(variables?.[key] || '').trim();
    return value ? encodeURIComponent(value) : '';
  });
}
```

### Conditional URL Resolution

```js
export function resolveToolUrl(toolId, variables = {}) {
  const tool = TOOL_REGISTRY[toolId];
  if (!tool) return { url: '', warning: 'Unknown diagnostic tool ID.', isDirect: false };

  const resolved = interpolateTemplate(tool.urlTemplate, variables);
  const missingPlaceholders = Array.from(tool.urlTemplate.matchAll(/\{([^}]+)\}/g))
    .map(match => match[1])
    .filter(key => !String(variables?.[key] || '').trim());

  if (missingPlaceholders.length) {
    return {
      url: tool.urlFallback,
      warning: `Using fallback URL because the case did not provide ${missingPlaceholders.join(', ')}. Verify before use.`,
      isDirect: false
    };
  }

  return { url: resolved, warning: '', isDirect: true };
}
```

### Resolution Rules

1. **Always return an object with `url`, `warning`, and `isDirect`.** Callers branch on `isDirect` to decide whether to show a warning badge.
2. **Detect missing variables explicitly.** Do not silently produce a URL with empty path segments. Fall back to `urlFallback` and list what was missing.
3. **Encode values.** Use `encodeURIComponent()` on every interpolated value to prevent URL corruption from special characters.
4. **Unknown tool IDs return an empty result.** Never throw. The caller gets a warning string to surface.

## 6. LLM Prompt Catalog Generation

The registry can be serialized into a compact text catalog for injection into LLM system
prompts, giving the model awareness of available diagnostic tools.

```js
export function buildDiagnosticToolsCatalogForPrompt() {
  return Object.entries(TOOL_REGISTRY)
    .map(([toolId, tool]) =>
      `${toolId} | ${tool.name} [${tool.sourceCategory.toLowerCase()}] | Use when: ${tool.whenToUse}`
    )
    .join('\n');
}
```

### Catalog Format

Each line follows: `tool_id | Display Name [category] | Use when: keyword, keyword, ...`

This format is optimized for LLM consumption:
- **Pipe-delimited** for clear field separation without JSON overhead.
- **Lowercase category** to reduce token noise.
- **whenToUse inline** so the model can match case symptoms to tools without a second lookup.

## 7. Suggestion Normalization from External Sources

When tools are suggested by an LLM or external system, the raw suggestions must be
validated, enriched, and capped before display.

```js
export function normalizeSuggestedDiagnosticTools(items = []) {
  return (Array.isArray(items) ? items : [])
    .map(item => ({
      toolId: String(item?.toolId || item?.tool_id || '').trim(),
      variables: item?.variables && typeof item.variables === 'object' ? item.variables : {},
      why: String(item?.why || '').trim()
    }))
    .filter(item => item.toolId && TOOL_REGISTRY[item.toolId])
    .slice(0, 10)
    .map(item => ({
      ...item,
      ...TOOL_REGISTRY[item.toolId],
      ...resolveToolUrl(item.toolId, item.variables)
    }));
}
```

### Normalization Rules

- **Accept both `toolId` and `tool_id`.** LLMs produce inconsistent casing. Normalize on input.
- **Reject unknown tool IDs.** If the ID is not in `TOOL_REGISTRY`, silently drop it. Never fabricate metadata.
- **Cap at 10 results.** Operator attention is finite. More than 10 diagnostic links causes decision paralysis.
- **Enrich from registry.** Merge the full tool metadata (name, description, steps, accessGuidance) onto each suggestion so the UI has everything it needs in one pass.
- **Resolve URLs per suggestion.** Each suggestion may carry its own `variables` (e.g., different project IDs), so URL resolution runs per item.

## 8. Deep-Link Construction (ts-diag)

For internal diagnostic platforms, build parameterized deep-links from case context,
tolerating missing fields gracefully.

```js
const TS_DIAG_BASE_URL = 'https://ts-diag.cloud-ops.prod.corp.mongodb.com/snapshot/project';

export function buildTsDiagUrl(caseContext = {}) {
  const snapshot = caseContext.snapshot || {};
  const params = new URLSearchParams({
    project: String(caseContext.projectId || snapshot.projectId || snapshot.project_id || '').trim(),
    org: String(caseContext.orgId || snapshot.orgId || snapshot.org_id || '').trim(),
    case: String(caseContext.caseNumber || '').trim(),
    cluster: String(caseContext.clusterId || snapshot.clusterId || snapshot.cluster_id || '').trim(),
    host: String(caseContext.atlasHost || snapshot.atlasHost || snapshot.atlas_host || '').trim()
  });
  return `${TS_DIAG_BASE_URL}?${params.toString()}`;
}
```

### Deep-Link Rules

- **Cascade through field aliases.** Case context may arrive from DOM scraping (`project_id`) or API normalization (`projectId`). Try both.
- **Leave missing params blank.** The receiving tool handles empty params more gracefully than omitted params or placeholder text.
- **Use `URLSearchParams`.** It handles encoding automatically and produces a canonical query string.

## 9. Testing Checklist

When building or extending a diagnostic registry, cover these scenarios:

- [ ] `getDisplayCategories()` returns the correct order and a defensive copy
- [ ] `getDiagnosticToolsRegistry()` returns a deep clone (mutating the clone does not affect the source)
- [ ] `resolveToolUrl()` returns `isDirect: true` with complete variables, `isDirect: false` with fallback URL and warning when variables are missing
- [ ] `getHighlightedDiagnosticTools()` marks matching tools as `highlighted: true` and leaves unrelated tools as `false`
- [ ] `getDiagnosticToolsForCategory('All')` returns all tools; specific categories filter correctly
- [ ] `normalizeSuggestedDiagnosticTools()` drops unknown IDs, accepts both `toolId` and `tool_id`, and caps at 10
- [ ] `buildDiagnosticToolsCatalogForPrompt()` includes every registry entry with `Use when:` keywords
- [ ] `buildTsDiagUrl()` preserves available identifiers and leaves missing ones blank

## 10. Anti-Patterns

| Anti-Pattern | Why It Fails | Correct Approach |
|---|---|---|
| Nested category hierarchies | Coupling display to data shape; hard to multi-categorize | Flat registry with `displayCategories` array |
| Throwing on unknown tool IDs | Crashes callers when LLM hallucinates an ID | Return empty result with warning string |
| Mutating the registry at runtime | Race conditions in concurrent reads | `structuredClone()` / `.slice()` on every read |
| Scoring on absence | "Tool X does not mention replication" penalizes unfairly | Score only on positive keyword matches |
| Alphabetical category sort | Destroys deliberate priority ordering | Use explicit `CATEGORY_ORDER` array |
| Unescaped regex from user text | ReDoS or false matches from special characters | `escapeRegExp()` before `new RegExp()` |
| Hardcoded URL construction | Breaks when context fields are missing | Template interpolation with fallback |
