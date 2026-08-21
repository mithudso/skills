<!-- hub-reference-banner -->
> **Reference file — part of the `atlas-diagnostics-expert` hub.** Formerly standalone `atlas-diagnostics-package` skill.
> Sibling MongoDB sub-topics now reference files under four hubs (`mongodb-expert`,
> `mongodb-atlas-expert`, `atlas-diagnostics-expert`, `mongodb-operations-expert`) — **not**
> standalone skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below naming
> bare `mongodb-*`/`atlas-*` skill; instead load that topic's `references/<name>.md`
> from owning hub (see hub's "Cross-hub map").

---

---
name: atlas-diagnostics-package
version: 1.1.0
updated: "2026-05-29"
description: >
  Expert reference for the @mdb-tam/atlas-diagnostics package and the diagnostic-recommendation
  layer built on top of it. Covers the tool registry data model, keyword-match highlighting,
  URL template interpolation, category filtering, LLM prompt catalog generation, ts-diag
  deep-link construction, scenario derivation, and the full case-diagnostic-recommendations
  pipeline used in the MDB Case Assistant Chrome extension.

  TRIGGER: user is adding/editing/removing entries from the diagnostic tool registry,
  building or modifying keyword-based tool highlighting, resolving Atlas/internal tool URLs
  from case context, generating the flat catalog string for LLM prompts, constructing ts-diag
  deep-link URLs, normalizing LLM-suggested tool arrays, extending the scenario-derivation or
  recommendation-ranking pipeline, or writing unit tests for any of the above.

  SKIP: general Atlas diagnostics workflows and operational triage — use atlas-diagnostics-expert.
  General Chrome extension architecture — use chrome-mv3-advanced. Case tracker lifecycle —
  use case-tracker.
origin: local
category: custom
tags:
  - atlas-diagnostics
  - chrome-extension
  - diagnostic-registry
  - tool-highlighting
  - case-recommendations
  - url-templates
  - ts-diag
triggers:
  - atlas diagnostics package
  - diagnostic tool registry
  - TOOL_REGISTRY
  - resolveToolUrl
  - getHighlightedDiagnosticTools
  - normalizeSuggestedDiagnosticTools
  - buildCaseDiagnosticRecommendations
  - getDiagnosticToolsForCategory
  - buildDiagnosticToolsCatalogForPrompt
  - diagnostic-registry
  - case-diagnostic-recommendations
  - ts-diag URL
  - tool highlighting
related_skills:
  - atlas-diagnostics-expert
  - case-tracker
  - mongodb-kb
  - 10gen
whenToUse:
  - "add, edit, or remove a tool from the diagnostic registry"
  - "build keyword-based tool highlighting against case text"
  - "resolve an Atlas tool URL from case context variables"
  - "generate the flat catalog string for LLM analysis prompts"
  - "construct a ts-diag snapshot deep-link URL"
  - "normalize an LLM-suggested diagnostic tool array"
  - "extend or debug the scenario-derivation pipeline"
  - "write unit tests for the atlas-diagnostics package"
whenNotToUse:
  - "operational Atlas triage workflows — use atlas-diagnostics-expert"
  - "general Chrome extension architecture decisions — use chrome-mv3-advanced"
  - "case tracker lifecycle, polling, or LLM analysis — use case-tracker"
---

# Atlas Diagnostics Package

Expert reference for `@mdb-tam/atlas-diagnostics` package and diagnostic-recommendation layer.

## Package location and structure

```
packages/atlas-diagnostics/
  package.json          # @mdb-tam/atlas-diagnostics, ESM, node:test runner
  src/index.js          # All exports — registry, helpers, public API
  tests/index.test.js   # Deterministic unit tests with node:test + assert
```

Re-exported for Chrome extension use:

```
src/background/diagnostic-registry.js   # re-exports everything from the package
```

Consumed by:

```
src/background/service-worker.js        # MCA_* message routing
src/background/case-tracker-analysis.js # tracked-case LLM analysis pipeline
src/panel/panel.js                      # overlay panel diagnostic tab
src/dashboard/dashboard.js              # case dashboard deep-dive
src/dashboard/case-diagnostic-recommendations.js  # full recommendation engine
src/shared/template-config.js           # LLM prompt template injection
```

---

## Tool registry data model

Every tool in `TOOL_REGISTRY` follows this shape:

```js
const TOOL_REGISTRY = {
  tool_id_snake_case: {
    name: 'Human-readable tool name',
    sourceCategory: 'Atlas Metrics',           // primary origin category
    displayCategories: ['Atlas Metrics', 'Profiler'],  // all UI filter categories
    guideUrl: 'https://docs.example.com/...',  // documentation link
    urlTemplate: 'https://cloud.mongodb.com/v2/{project_id}#/clusters/{cluster_name}/metrics',
    urlFallback: 'https://cloud.mongodb.com/', // used when variables are missing
    description: 'One-line description of what the tool surface shows.',
    whenToUse: 'comma, separated, keyword, phrases, for, matching',
    accessGuidance: [
      'Step to confirm access or install requirements.',
      'RBAC or auth requirements.'
    ],
    steps: [
      'Open the tool or surface',
      'Inspect the relevant data',
      'Correlate findings with the case'
    ]
  }
};
```

### Key design rules

1. `toolId` = object key (snake_case), stable identifier everywhere.
2. `sourceCategory` = single origin label; `displayCategories` = array for UI tab filtering.
3. `urlTemplate` uses `{variable_name}` placeholders. Interpolation function URI-encodes each substituted value.
4. `urlFallback` returned when any placeholder variable missing from case context.
5. `whenToUse` = comma-separated keyword phrases. Each matched independently against normalized case text during highlighting.
6. `steps` = short operator instructions shown in recommendation card UI.

### Current category order

```js
const CATEGORY_ORDER = [
  'Atlas Metrics', 'Profiler', 'Logs', 'Connection',
  'Replication', 'Storage', 'Security'
];
```

### Registered tools (16 total)

| Group | Tools |
|-------|-------|
| Atlas Metrics | `atlas_realtime_performance`, `atlas_performance_advisor`, `atlas_query_profiler`, `atlas_cluster_metrics` |
| Atlas Admin | `atlas_cluster_overview`, `atlas_logs`, `atlas_network_access`, `atlas_database_access`, `atlas_backup_restore`, `atlas_search_explorer` |
| Diagnostics | `mongosh_diagnostics`, `mongodb_compass` |
| Reference | `mongodb_docs`, `mongodb_jira` |
| Internal Tools | `hub_account_page`, `ts_tools_case` |

---

## Public API reference

### getDisplayCategories()

Returns defensive copy of `CATEGORY_ORDER`. Use for building category tabs.

```js
import { getDisplayCategories } from '../background/diagnostic-registry.js';
const tabs = getDisplayCategories();
// ['Atlas Metrics', 'Profiler', 'Logs', 'Connection', 'Replication', 'Storage', 'Security']
```

### getDiagnosticToolsRegistry()

Returns `structuredClone` of full `TOOL_REGISTRY`. Safe to mutate without affecting canonical data.

```js
const registry = getDiagnosticToolsRegistry();
const tool = registry.atlas_performance_advisor;
// tool.name === 'Atlas Performance Advisor'
```

### buildDiagnosticToolsCatalogForPrompt()

Returns single multi-line string listing every tool in `toolId | name | Use when: ...` format. Injected into LLM analysis prompts so model can suggest diagnostic tools by ID.

```js
const catalog = buildDiagnosticToolsCatalogForPrompt();
// atlas_realtime_performance | Atlas Real-Time Performance Panel [atlas metrics] | Use when: slow queries, high CPU, ...
// atlas_performance_advisor | Atlas Performance Advisor [atlas metrics] | Use when: slow queries, missing index, ...
```

### resolveToolUrl(toolId, variables)

Interpolates `urlTemplate` placeholders. Returns `{ url, warning, isDirect }`.

```js
// Full resolution
resolveToolUrl('atlas_cluster_metrics', { project_id: 'proj1', cluster_name: 'C0' });
// { url: 'https://cloud.mongodb.com/v2/proj1#/clusters/C0/metrics', warning: '', isDirect: true }

// Missing variables — falls back
resolveToolUrl('atlas_cluster_metrics', {});
// { url: 'https://cloud.mongodb.com/', warning: '...did not provide project_id, cluster_name...', isDirect: false }

// Unknown tool
resolveToolUrl('nonexistent_tool', {});
// { url: '', warning: 'Unknown diagnostic tool ID.', isDirect: false }
```

Always check `isDirect` before treating resolved URL as valid deep link. When `isDirect` is false, URL is fallback.

### getHighlightedDiagnosticTools(caseContext)

Scores every registered tool against case text, marks each `highlighted` (true/false) with `whyRelevant` explanation string.

```js
const caseContext = {
  subject: 'Slow query on Atlas Search index',
  description: 'Customer reports high latency on $search aggregation.',
  errorStrings: ['query targeting ratio exceeded threshold'],
  status: 'In Progress'
};

const tools = getHighlightedDiagnosticTools(caseContext);
const relevant = tools.filter(t => t.highlighted);
// Each element: all TOOL_REGISTRY fields + { toolId, highlighted: boolean, whyRelevant: string }
```

**Keyword matching algorithm:**
1. Combine subject + description + errorStrings + status into one haystack.
2. Normalize to lowercase, strip non-alphanumeric except spaces.
3. For each tool, split `whenToUse` on commas into keyword phrases.
4. Test each normalized keyword against haystack with word-boundary regex, falling back to substring inclusion.
5. Tool highlighted if at least one keyword matches.

### getDiagnosticToolsForCategory(category, caseContext)

Filters `getHighlightedDiagnosticTools` results by `displayCategories`. Pass `'All'` or empty string for unfiltered results.

### normalizeSuggestedDiagnosticTools(items)

Takes raw array from LLM analysis response (may contain `toolId` or `tool_id`, bad shapes, unknown IDs, or nulls), returns clean validated enriched array. Maximum 10 results.

```js
normalizeSuggestedDiagnosticTools([
  { toolId: 'atlas_cluster_metrics', variables: { project_id: 'p1' }, why: 'Check latency' },
  { tool_id: 'atlas_logs', why: 'Review errors' },   // tool_id also accepted
  { toolId: 'unknown_tool' },                         // filtered out
  null                                                  // filtered out
]);
// Returns enriched array: [{ toolId, name, url, isDirect, why, ...registryFields }, ...]
```

### buildTsDiagUrl(caseContext)

Constructs ts-diag snapshot deep-link URL. Missing fields become empty query parameters rather than omitted.

```js
buildTsDiagUrl({ caseNumber: '01581027', projectId: 'p1', orgId: 'o1', clusterId: 'c1', atlasHost: 'h1' });
// https://ts-diag.cloud-ops.prod.corp.mongodb.com/snapshot/project?project=p1&org=o1&case=01581027&cluster=c1&host=h1

// Also reads from nested snapshot object: caseContext.snapshot.projectId, snapshot.org_id, etc.
```

---

## Case diagnostic recommendations engine

Higher-level recommendation pipeline in `src/dashboard/case-diagnostic-recommendations.js`.

### buildCaseDiagnosticRecommendations(caseRecord, analysis)

Returns `{ scenarios, recommendations }`.

```js
const { scenarios, recommendations } = buildCaseDiagnosticRecommendations(caseRecord, analysis);

// scenarios: [{ id, score, label, reasons: string[] }]
//   id values: 'performance' | 'replication' | 'connection' | 'security' | 'search' | 'storage' | 'logs' | 'general'

// recommendations (max 7): [{ id, type, title, sourceLabel, summary, why, badges,
//   actionUrl, actionLabel, actionWarning, usage, install, command, steps }]
//   id format: 'tool:{toolId}' | 'repo:{repoId}' | 'command:{type}'
```

### Scenario derivation

Scenarios scored by:
1. **Keyword matches** against case evidence (subject 4×, analysis 3×, description 2×, comments 1×)
2. **Tool category mapping** — highlighted tools add 2 points per matching category; LLM-suggested tools add 3 points
3. Special boost for `atlas_search_explorer` mapping to `search` scenario

No scenario scores above zero → `general` fallback returned. UI must handle gracefully.

### Recommendation types

| Type | Source | Example |
|------|--------|---------|
| `tool` | Registry tools matched by highlighting or LLM suggestion | Atlas Performance Advisor |
| `repo` | 10gen verified related tools ranked by scenario tags | ts-diag, alexandria |
| `command` | Generated mongosh commands for top scenario | `rs.status()`, `serverStatus` |

### Recommendation ranking order

1. LLM-suggested tools (from `analysis.diagnosticTools`)
2. Internal tools (Hub account page, TS Tools case page) when IDs available
3. Repo recommendations ranked by scenario overlap (max 2)
4. Keyword-highlighted tools not already included (max 3)
5. Scenario-specific mongosh command (max 1)
6. Dedup and cap at 7 total

### URL resolution with case context

Variable resolution priority: `tool.variables` fields take precedence over `caseRecord` fields (`accountId`, `snapshot.clusterName`, `caseNumber`). Action labels inferred from resolved URL domain:

| Domain | Label |
|--------|-------|
| `cloud.mongodb.com` | "Open in Atlas" |
| `hub.corp.mongodb.com` | "Open in Hub" |
| `ts-tools.corp.mongodb.com` | "Open in TS Tools" |
| Other | "Open link" |
| Fallback URL | "Open Atlas home" or "Open tool home" |

---

## Adding a new diagnostic tool

1. Add entry to `TOOL_REGISTRY` in `packages/atlas-diagnostics/src/index.js` (all fields required).
2. New category needed → add to `CATEGORY_ORDER`.
3. New URL template variable needed → update `getToolVariables` in `case-diagnostic-recommendations.js`.
4. Run tests: `cd packages/atlas-diagnostics && node --test`.
5. Verify highlighting works for case text matching `whenToUse` keyword.

---

## Testing

```bash
cd packages/atlas-diagnostics && node --test
```

Test categories for any new tool or API change:
- URL resolution (direct + fallback + unknown tool ID)
- Keyword highlighting (positive match + negative non-match)
- Defensive copy (mutate clone, verify original unchanged)
- Normalization (valid + invalid + mixed input arrays)

---

## Anti-patterns

- **Do not** import `TOOL_REGISTRY` or `CATEGORY_ORDER` directly and mutate. Always use exported accessor functions returning clones/copies.
- **Do not** build diagnostic URLs by hand-concatenating strings. Always use `resolveToolUrl` so fallback handling, URI encoding, and warning generation stay consistent.
- **Do not** duplicate keyword-matching logic outside `getHighlightedDiagnosticTools`. Normalization and regex-fallback behavior must stay in one place.
- **Do not** add scenario-scoring logic in UI code. Scenario derivation belongs in `case-diagnostic-recommendations.js`.
- **Do not** skip `isDirect` check on resolved URLs. Displaying fallback URL without warning misleads operator.
- **Do not** hardcode tool IDs in UI rendering code. Use registry, let highlighting/recommendation pipeline select tools dynamically.
- **Do not** add new recommendation types without updating dedup `id` prefix scheme (`tool:`, `repo:`, `command:`).

## Common pitfalls

1. **Mutating registry.** Always use `getDiagnosticToolsRegistry()` (clones) or `getDisplayCategories()` (slices). Never import and modify `TOOL_REGISTRY` or `CATEGORY_ORDER` directly.
2. **Assuming URL templates complete.** Always check `isDirect` before treating resolved URL as valid deep link.
3. **Forgetting URI-encode.** `interpolateTemplate` helper already calls `encodeURIComponent`. Do not double-encode.
4. **Keyword ordering in `whenToUse`.** Put most specific phrases first. `whyRelevant` string lists matches in order — specificity first produces better explanations.
5. **Exceeding 10-tool cap.** `normalizeSuggestedDiagnosticTools` silently caps at 10. Extras dropped.
6. **Missing recommendation dedup.** Engine deduplicates by `id` (format `tool:{toolId}`, `repo:{repoId}`, or `command:{type}`). New recommendation sources must avoid collisions.
7. **Scenario fallback.** No scenario scores above zero → engine returns single `general` scenario. UI must handle gracefully.

---

## Atlas API patterns for Chrome extension diagnostics

Full Atlas API endpoint docs in `atlas-diagnostics-expert` skill. Key patterns:

| Endpoint pattern | Use |
|-----------------|-----|
| `GET /api/atlas/v2/groups/{groupId}/processes/{processId}/measurements?m=CONNECTIONS&...` | Metrics |
| `GET .../performanceAdvisor/suggestedIndexes` | Performance Advisor index suggestions |
| `GET .../performanceAdvisor/slowQueryLogs` | Slow query logs |
| `GET .../clusters/{clusterName}/logs/{logName}.gz` | Logs (30-day retention) |
| `GET .../alertConfigs` | Alert configs |
| `GET .../alerts?status=OPEN` | Open alerts |

**Auth in extensions:** Bearer tokens via `POST https://cloud.mongodb.com/api/oauth/token` (3600s TTL). MDB Case Assistant uses three-tier fallback: same-tab cookie fetch → extension fetch with cookies → bearer-token fallback.