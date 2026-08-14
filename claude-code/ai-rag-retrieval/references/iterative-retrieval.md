<!-- hub-reference-banner -->
> **Reference file — part of the `ai-rag-retrieval` hub.** Formerly the standalone `iterative-retrieval` skill.
> Sibling topics in this family are now reference files under the hubs (`ai-agents-orchestration`, `ai-rag-retrieval`, `ai-llm-model-layer`, `ai-mcp-sdk-prompting`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: iterative-retrieval
version: "1.1.1"
updated: "2026-05-31"
description: >
  Iterative context refinement pattern for multi-agent workflows. Solves the cold-start
  context problem where subagents don't know what files or patterns they need until they
  start working. Runs a 4-phase DISPATCH→EVALUATE→REFINE→LOOP cycle (max 3 iterations)
  to progressively build a high-relevance context set, stopping when ≥3 files score
  ≥0.7 relevance and no critical gaps remain.
  TRIGGER: user asks about subagent context, cold-start retrieval, progressive file
  discovery, relevance scoring for agents, or "how should my agent find the right files."
  SKIP: simple single-file lookups; questions about RAG pipeline architecture (→ ai-agent-engineering);
  deep research workflows (→ deep-research).
category: developer
tags:
  - agents
  - retrieval
  - context-engineering
  - subagents
  - relevance-scoring
keywords:
  - iterative retrieval
  - cold-start context
  - subagent context
  - progressive file discovery
  - relevance scoring
  - context refinement
whenToUse:
  - "How should my subagent find the right files for a task?"
  - "Subagent cold-start context problem"
  - "Progressive file discovery with relevance scoring"
  - "Build a focused context set without exceeding token limits"
origin: local
related_skills:
  - deep-research-methods
  - rag-architecture
---

# Iterative Retrieval Pattern

Solves the cold-start context problem in multi-agent workflows: a subagent cannot know which files or patterns it needs until it begins working. Rather than guessing or flooding with everything, use a 4-phase cycle to progressively refine context toward high-relevance results.

## When to Use

- A subagent needs to discover which codebase files are relevant to its task
- Initial keyword searches are missing because the codebase uses different terminology
- You need to build a focused context set (3–5 high-relevance files) without exceeding token limits
- Multi-agent workflows where later queries should build on earlier findings

## When NOT to Use

- **Single known file** — read it directly; no iteration needed
- **RAG pipeline design** — use `ai-agent-engineering`
- **Web research workflows** — use `deep-research`

## The Problem

Subagents launch with limited context. They don't know:
- Which files contain relevant code
- What patterns exist in the codebase
- What terminology the project uses

Standard approaches fail:
- **Send everything**: exceeds context limits
- **Send nothing**: agent lacks critical information
- **Guess what's needed**: frequently wrong

## Solution: 4-Phase Iterative Cycle

```
DISPATCH → EVALUATE → REFINE → LOOP (max 3 cycles)
     ↑                              |
     └──────────────────────────────┘
```

### Phase 1: DISPATCH

Start broad. Collect candidate files using high-level intent keywords:

```javascript
const initialQuery = {
  patterns: ['src/**/*.ts', 'lib/**/*.ts'],
  keywords: ['authentication', 'user', 'session'],
  excludes: ['*.test.ts', '*.spec.ts']
};
const candidates = await retrieveFiles(initialQuery);
```

### Phase 2: EVALUATE

Score each candidate for relevance to the task (0–1 scale):

```javascript
function evaluateRelevance(files, task) {
  return files.map(file => ({
    path: file.path,
    relevance: scoreRelevance(file.content, task),   // 0.0–1.0
    reason: explainRelevance(file.content, task),
    missingContext: identifyGaps(file.content, task)
  }));
}
```

Relevance thresholds:
- **0.8–1.0** — directly implements the target feature
- **0.5–0.7** — contains related patterns or types
- **0.2–0.4** — indirectly related
- **0.0–0.2** — not relevant; exclude on the next cycle

### Phase 3: REFINE

Update the query based on what the evaluation revealed:

```javascript
function refineQuery(evaluation, previousQuery) {
  return {
    // Add patterns found in high-relevance files
    patterns: [...previousQuery.patterns, ...extractPatterns(evaluation)],
    // Add terminology discovered in the codebase
    keywords: [...previousQuery.keywords, ...extractKeywords(evaluation)],
    // Exclude confirmed irrelevant paths
    excludes: [
      ...previousQuery.excludes,
      ...evaluation.filter(e => e.relevance < 0.2).map(e => e.path)
    ],
    // Target specific gaps identified in evaluation
    focusAreas: evaluation.flatMap(e => e.missingContext).filter(unique)
  };
}
```

### Phase 4: LOOP

Repeat with the refined query. Stop early when context is sufficient:

```javascript
async function iterativeRetrieve(task, maxCycles = 3) {
  let query = createInitialQuery(task);
  let bestContext = [];

  for (let cycle = 0; cycle < maxCycles; cycle++) {
    const candidates = await retrieveFiles(query);
    const evaluation = evaluateRelevance(candidates, task);
    const highRelevance = evaluation.filter(e => e.relevance >= 0.7);

    // Stop when ≥3 high-relevance files found with no critical gaps
    if (highRelevance.length >= 3 && !hasCriticalGaps(evaluation)) {
      return highRelevance;
    }

    query = refineQuery(evaluation, query);
    bestContext = mergeContext(bestContext, highRelevance);
  }

  return bestContext;
}
```

## Worked Examples

### Example 1: Bug Fix Context

```
Task: "Fix authentication token expiry bug"

Cycle 1:
  DISPATCH: search src/** for "token", "auth", "expiry"
  EVALUATE: auth.ts (0.9), tokens.ts (0.8), user.ts (0.3)
  REFINE: add "refresh", "jwt" keywords; exclude user.ts

Cycle 2:
  DISPATCH: search with refined terms
  EVALUATE: session-manager.ts (0.95), jwt-utils.ts (0.85) found
  STOP: ≥3 high-relevance files, no critical gaps

Result: auth.ts, tokens.ts, session-manager.ts, jwt-utils.ts
```

### Example 2: Terminology Mismatch

```
Task: "Add rate limiting to API endpoints"

Cycle 1:
  DISPATCH: search routes/** for "rate", "limit", "api"
  EVALUATE: zero matches — codebase uses "throttle" not "rate limit"
  REFINE: add "throttle", "middleware" keywords

Cycle 2:
  DISPATCH: search with "throttle", "middleware"
  EVALUATE: throttle.ts (0.9), middleware/index.ts (0.7)
  REFINE: router patterns still missing

Cycle 3:
  DISPATCH: search for "router", "express" patterns
  EVALUATE: router-setup.ts (0.8) found
  STOP: sufficient context

Result: throttle.ts, middleware/index.ts, router-setup.ts
```

## Agent Prompt Integration

Add this instruction block when dispatching a subagent:

```markdown
When gathering context for this task:
1. Start with broad keyword searches on the high-level intent
2. Score each file's relevance (0–1) and explain why
3. Identify what context is still missing
4. Refine search criteria based on terminology and patterns found
5. Stop when ≥3 files score ≥0.7 and no critical gaps remain (max 3 cycles)
6. Return only files with relevance ≥0.7
```

## Best Practices

1. **Start broad, narrow progressively** — don't over-specify the initial query
2. **Learn the codebase's vocabulary** — first cycle often surfaces naming conventions
3. **Track gaps explicitly** — naming what's missing drives focused refinement
4. **Stop at "good enough"** — 3 high-relevance files beat 10 marginal ones
5. **Exclude with confidence** — files scoring < 0.2 won't become relevant on re-query

## Stopping Criteria Summary

| Condition | Action |
|-----------|--------|
| ≥3 files score ≥0.7 AND no critical gaps | Stop immediately |
| Max 3 cycles reached | Return best context so far |
| New cycle adds no files ≥0.7 (stagnation) | Stop; return current best |
