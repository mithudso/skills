---
name: firecrawl-deep-research
description: Run multi-source deep research with Firecrawl using the Deep Research / Research API. Use when the user asks to research a topic, compare perspectives, produce a sourced briefing, investigate a technical or market question, or synthesize web evidence across many sources (or invokes /dr).
license: ISC
metadata:
  author: firecrawl
  version: "0.2.0"
  homepage: https://www.firecrawl.dev
  source: https://github.com/firecrawl/firecrawl-workflows
inputs:
  - name: FIRECRAWL_API_KEY
    description: Firecrawl API key for hosted Firecrawl requests.
    required: true
---

# Firecrawl Deep Research (/dr)

Use this when the user wants a sourced research report, not raw search results, or invokes the `/dr` command.

## Onboarding Interview

Infer the topic, depth, and output format from context. If the topic is clear, proceed immediately.

Ask at most 1-3 concise questions only if blocked, such as the research topic, required depth, or a critical angle/source constraint.

## Firecrawl Research API Execution Plan

Primary method: Use the **Firecrawl Research API / Deep Research API** (`firecrawl_deep_research` or `firecrawl deep-research` / `POST /v1/deep-research`) to conduct autonomous multi-step web research.

### API Invocation

- **MCP Tool**: `firecrawl_deep_research(query: "<topic>", max_depth: <depth>, max_urls: <count>)`
- **CLI Command**: `firecrawl deep-research "<topic>" [options]`
- **REST Endpoint**: `POST /v1/deep-research` with `{ "query": "<topic>", "depth": "<quick|thorough|exhaustive>" }`

### Depth Settings

- **Quick**: `max_depth: 1-2`, `max_urls: 5-10` — fast overview of key facts and sources.
- **Thorough**: `max_depth: 3-4`, `max_urls: 15-25` — multi-angle investigation across technical, market, and primary sources.
- **Exhaustive**: `max_depth: 5+`, `max_urls: 25+` — deep iteration including technical papers, contrarian views, and official documentation.

### Academic & Technical Research Index API

For paper-focused or academic research topics, supplement or use the **Firecrawl Research Index API**:
- MCP: `firecrawl_research_search_papers(query)`, `firecrawl_research_related_papers(seed_ids, intent)`, `firecrawl_research_read_paper(id, question)`
- CLI: `firecrawl research search-papers "<query>"`, `firecrawl research related-papers <seedIds...>`

### Fallback / Manual Execution

If the managed Research API endpoint is unauthenticated or unavailable, fall back to manual iterative collection using `firecrawl search` and `firecrawl scrape`:
- Quick: search 3-5 queries and scrape 5-10 high-quality sources.
- Thorough: search 5-10 queries from different angles and scrape 15-25 sources.
- Exhaustive: search 10+ queries and scrape 25+ sources, including primary sources, research papers, expert views, and contrarian sources.

Avoid re-scraping URLs already returned with full content from a search-with-scrape result.

## Parallel Work

If appropriate, use sub-agents or equivalent parallel task runners by research angle when using manual fallback or multi-query synthesis:

- overview and definitions
- technical or implementation details
- market and industry context
- contrarian views, risks, and limitations
- primary sources and official docs

Each researcher should return claims, source URLs, source quality notes, and uncertainty.

## Final Deliverable

Default structure:

```markdown
# Deep Research: [Topic]

## Executive Summary
[2-3 paragraphs]

## Key Findings
[Numbered findings with source links]

## Detailed Analysis
[Themes, evidence, and synthesis]

## Contrarian Views And Risks
[Counterarguments, limitations, failure modes]

## Open Questions
[What remains uncertain]

## Sources
[Every URL used with a one-line note]

## Rerun Inputs
workflow: firecrawl-deep-research
topic: [topic]
depth: [quick/thorough/exhaustive]
output: [markdown/json/brief]
```

## Quality Bar

- Cite sources for factual claims.
- Prefer primary sources when available.
- Flag uncertainty and conflicting evidence.
- Synthesize instead of listing scrape summaries.

