<!-- hub-reference-banner -->
> **Reference file — part of the `da-3-data-acquisition-sampling` hub.** Formerly the standalone `da-3-2-6-apis-data-feeds` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-3-2-6-apis-data-feeds
description: |
  Domain knowledge for collecting data via APIs and structured data feeds for analysis purposes.
  Covers REST, GraphQL, streaming feeds, authentication, pagination, rate limiting, push vs pull
  patterns, incremental extraction, and data quality considerations — all from the analyst's
  perspective of acquiring reliable data for downstream analysis.

  TRIGGER: when the user is collecting data from an external API or data feed for analysis or
  reporting purposes; when designing or debugging an API-based data ingestion pipeline; when
  choosing between polling and webhooks for data freshness; when handling pagination, rate limits,
  or incremental sync from an API; when evaluating REST vs GraphQL for data acquisition; when
  asking about public/open data APIs (government, financial, social) as data sources; when
  troubleshooting data quality issues originating from upstream API schema changes; when asking
  which HTTP error codes to retry and how to handle auth failures in a data pipeline.

  SKIP: defer to `api-design-patterns` for questions about designing or building API servers and
  contracts; defer to `da-14-streaming-analytics` or `mongodb-change-streams` for CDC / streaming
  ingestion architectures and event-processing pipelines; defer to `sse-streaming-patterns` or
  `websocket-extension-patterns` for implementing streaming transports in application code; defer
  to `backend-patterns` or `express-patterns` for server-side API implementation engineering;
  defer to `web-auth-patterns` for implementing OAuth or auth flows inside your own application
  (not for calling a third-party API from a data pipeline); defer to
  `da-13-data-engineering-and-pipelines` for pipeline orchestration, scheduling architecture, and
  multi-source data platform design.
---

# APIs & Data Feeds as Data Collection Methods

## 1. What APIs and Data Feeds Are (in the Analyst Context)

An **API (Application Programming Interface)** is a programmatic interface that allows a data consumer to request structured data from an external system on demand. In the data collection context, the analyst or pipeline acts as the *client*: it sends requests, the API server returns records, and those records become raw material for analysis.
([Integrate.io](https://www.integrate.io/blog/top-rest-api-best-practices-for-data-integration/))

A **data feed** is a regularly updated stream or delivery of records from a source system, typically following a fixed format and schedule. Feeds may be push-based (the source delivers records when they are ready) or pull-based (the consumer fetches them on a schedule). Common feed formats include JSON, XML, CSV, Atom, and RSS. While APIs let you query and filter dynamically, feeds tend to deliver complete or pre-scoped snapshots on a schedule.
([PredictLeads](https://blog.predictleads.com/2026/05/26/company-data-api-flat-files-webhooks-mcp))

The analyst's concern is not how to build an API server but how to *reliably acquire data* from an API or feed with acceptable completeness, freshness, and cost.

---

## 2. API Styles Relevant to Data Collection

### 2.1 REST APIs

REST is the dominant style, powering ~83% of web services as of 2024.
([jsonconsole.com](https://jsonconsole.com/blog/rest-api-vs-graphql-statistics-trends-performance-comparison-2025))
Resources are addressed by URL, verbs are HTTP methods (GET for reads), and responses are typically JSON. REST is well-suited for:
- Simple or well-understood data shapes where one endpoint per resource type is enough.
- Scenarios where HTTP caching matters (GET responses are cacheable).
- Broad third-party provider support (virtually every external data source exposes REST).

REST's weakness for analysts: **over-fetching** (endpoints return more fields than needed) and **under-fetching** (multiple round-trips needed to assemble a record). Both inflate bandwidth and slow pipelines.
([API7.ai](https://api7.ai/blog/graphql-vs-rest-api-comparison-2025))

### 2.2 GraphQL APIs

GraphQL lets the client specify exactly which fields to return in a single request. This eliminates over/under-fetching and typically reduces bandwidth 30–50% for complex, multi-relationship queries.
([jsonconsole.com](https://jsonconsole.com/blog/rest-api-vs-graphql-statistics-trends-performance-comparison-2025))

Tradeoffs to weigh for data collection:
- GraphQL parsing and resolution costs 20–35% more CPU on the server than REST for simple lookups, so simple full-table extractions may run slower. ([Postman Blog](https://blog.postman.com/graphql-vs-rest/))
- HTTP-level caching is harder; GET-based query caching requires persisted queries.
- Pagination and rate limiting are less standardized — each provider implements them differently (cursor-based, `after`/`first`, etc.).

GraphQL is the right choice when the provider exposes it and you need to pull denormalized data spanning multiple entity types in one query (e.g., GitHub's GraphQL API for issue + comment + label in a single request).

### 2.3 Streaming Feeds (SSE / WebSocket)

For data that changes continuously — financial ticks, social media events, sensor readings — polling REST endpoints wastes resources and introduces freshness lag. Two options:
- **Server-Sent Events (SSE):** one-way server-to-client stream over HTTP. Simple, reconnects automatically, works through most proxies. About 60% of developers chose SSE for lightweight streaming in 2025, citing ~25% lower resource use than WebSocket for read-only streams. ([Medium/serifcolakel](https://medium.com/@serifcolakel/real-time-data-streaming-with-server-sent-events-sse-9424c933e094))
- **WebSocket:** bidirectional. Use when the data source requires the client to send messages (e.g., subscription filters). Higher overhead than SSE for pure read cases. ([VideoSDK](https://www.videosdk.live/developer-hub/websocket/websocket-streaming))

Streaming feeds are a *collection-side* concern here (receiving a feed), not building one.

---

## 3. Authentication Patterns

Every external API requires proof of identity before delivering data. Three patterns dominate data-collection contexts:

| Pattern | When to use | Key analyst concern |
|---|---|---|
| **API key** (header or query param) | Service-to-service; simple scripts; public-data APIs | Never embed in source code or URLs; store in env vars or a secrets manager; rotate every 90 days. ([Google Cloud](https://docs.cloud.google.com/docs/authentication/api-keys-best-practices)) |
| **OAuth 2.0** | User-data APIs (Google Analytics, Salesforce, social) where delegated access is required | Access tokens expire (often 1 h); use a token refresh flow; credentials never reach the analyst's code. ([Axway](https://blog.axway.com/learning-center/digital-security/keys-oauth/api-keys-oauth)) |
| **JWT / Bearer token** | Microservice-to-microservice; internal data platforms | Stateless; verify signature and expiry before caching; use short TTLs. |

Practical rules for pipeline authentication:
- Always use HTTPS; every auth method is broken without encrypted transport.
- Issue per-pipeline credentials so compromising one key limits blast radius.
- Avoid long-lived keys; shorter validity windows reduce the exposure window.
- Follow OAuth 2.1 guidance for new integrations: require PKCE for authorization code flows, deprecate Implicit Grant.

([APIsec](https://www.apisec.ai/blog/master-api-authentication-and-authorization-best-practices-for-security), [NCSC](https://www.ncsc.gov.uk/collection/securing-http-based-apis/2-api-authentication-and-authorisation))

---

## 4. Pagination

APIs rarely return all records in one response. Pagination splits results into pages; analysts must traverse all pages to collect complete datasets.

### Pagination strategies

**Offset / page-number pagination**
`GET /records?page=3&page_size=100`
- Simple to implement; lets you jump to any page.
- Breaks under concurrent writes: if rows are inserted between page 2 and page 3 requests, records shift and you may skip or double-count rows.
- Acceptable for small, static datasets.

**Cursor-based pagination**
`GET /records?after=eyJpZCI6MTAwfQ&limit=100`
- The server returns an opaque cursor token pointing to the next page. The client echoes it back.
- Stable under concurrent writes: uses a database index position, not a row count.
- Preferred for large, frequently-changing datasets.
([merge.dev](https://www.merge.dev/blog/api-pagination-best-practices))

**Keyset pagination**
`GET /records?since_id=12345&limit=100`
- Uses a sortable column value (ID, timestamp) as the boundary.
- Predictable and index-friendly; requires the sort key to be stable and indexed.

### Pagination best practices for analysts

- Parse the `next` / `has_more` / `Link` header in each response; do not assume a fixed number of pages.
- Cache cursor tokens with their expiry metadata — expired cursors restart the full extraction.
- Set `limit` to the provider's allowed maximum to minimize round-trips, but respect the server's default page size if not stated.
- Add a 100–500 ms delay between page requests to stay within rate limits without triggering backoff (tune to the provider's documented requests-per-second ceiling).

([merge.dev](https://www.merge.dev/blog/api-pagination-best-practices), [cybergarden.au](https://cybergarden.au/blog/rest-api-pagination-best-practices))

---

## 5. Rate Limiting and Throttling

Providers impose rate limits to protect their infrastructure. Violating them yields HTTP 429 (`Too Many Requests`).

**Rate limiting** is a hard cap: once the window's quota is consumed, further requests are rejected until the window resets.
**Throttling** is a soft approach: the server queues or slows requests rather than rejecting them outright.
([digitalapi.ai](https://www.digitalapi.ai/blogs/api-throttling), [getknit.dev](https://www.getknit.dev/blog/10-best-practices-for-api-rate-limiting-and-throttling))

### Reading rate-limit headers proactively

Most APIs expose remaining quota before you hit the ceiling. Read these on every response:

| Header | Meaning |
|---|---|
| `X-RateLimit-Limit` | Total requests allowed in the current window |
| `X-RateLimit-Remaining` | Requests left before hitting the limit |
| `X-RateLimit-Reset` | Unix timestamp when the window resets |
| `Retry-After` | Seconds to wait (present on 429 responses) |

When `X-RateLimit-Remaining` drops below 10% of `X-RateLimit-Limit`, slow the request rate proactively rather than waiting for a 429.

### Handling 429s: exponential backoff with jitter

When a 429 is received:
1. Read the `Retry-After` header for the exact wait time the server specifies.
2. If absent, apply exponential backoff: wait `min(base * 2^attempt, max_wait)` seconds.
3. Add random jitter (± 10–30% of the computed wait) to prevent a *thundering herd* — multiple pipeline workers retrying simultaneously.
4. Only retry on 429, 500, and 503; other status codes indicate a client-side problem that retrying will not fix.

([OpenAI Cookbook](https://developers.openai.com/cookbook/examples/how_to_handle_rate_limits), [Redis Blog](https://redis.io/blog/api-throttling-algorithms-patterns/))

### Analyst-side rate management
- Track requests-per-window against known limits; shed load proactively before hitting the ceiling.
- Use concurrency controls (semaphores, queues) when parallelizing page requests.
- For bulk historical pulls, schedule them outside business hours (e.g., 2–6 AM local time for the provider's region) to avoid competing with production traffic against shared rate-limit quotas.

---

## 6. HTTP Error Handling for Data Pipelines

Not all HTTP errors warrant the same response. This table covers the codes analysts encounter most when collecting data:

| Status | Meaning | Pipeline action |
|---|---|---|
| `200 OK` | Success | Ingest and advance watermark |
| `206 Partial Content` | Partial result (range request) | Accumulate chunks; advance only after final chunk |
| `400 Bad Request` | Malformed request (bad filter, wrong param type) | Log and fail the run; fix query before retrying |
| `401 Unauthorized` | Auth token missing or expired | Refresh token / re-authenticate; then retry once |
| `403 Forbidden` | Token valid but scope insufficient | Alert operator; do not retry automatically |
| `404 Not Found` | Resource does not exist | Log as missing; skip record; do not fail the run |
| `429 Too Many Requests` | Rate limit hit | Exponential backoff with jitter (see §5) |
| `500 Internal Server Error` | Transient server fault | Retry with backoff (up to 3 attempts) |
| `503 Service Unavailable` | Server overloaded / maintenance | Retry with backoff; respect `Retry-After` if present |

Key rule: only retry on 429, 500, and 503. A 400 or 403 will not self-resolve — retrying wastes quota and masks the underlying misconfiguration.

---

## 7. Pull (Polling) vs. Push (Webhooks)

### Polling
The pipeline calls the API on a schedule (every N minutes) and fetches whatever is new.

- Full control over timing and volume.
- Simple to implement and reason about.
- Wastes resources when data is sparse: most poll cycles may return zero new records.
- Freshness is bounded by the poll interval.

### Webhooks (push)
The source system sends an HTTP POST to a listener URL the analyst registers whenever an event occurs.

- Near-real-time delivery; no wasted requests.
- Requires the analyst's system to expose a publicly reachable endpoint.
- Needs careful error handling: verify signatures (e.g., HMAC), acknowledge immediately (return 200), process asynchronously, handle retries idempotently.

Practical rule: if you find yourself polling every few seconds to detect changes, that is a strong signal to switch to webhooks.
([authgear.com](https://www.authgear.com/post/webhooks-vs-apis-difference/), [PredictLeads](https://blog.predictleads.com/2026/05/26/company-data-api-flat-files-webhooks-mcp))

---

## 8. Incremental Extraction

Full extraction (re-downloading all records every run) is expensive and unnecessary for large, append-heavy datasets. Incremental extraction fetches only records created or modified since the last successful run.

### High-watermark pattern
1. At the start of each run, read the stored watermark (e.g., `last_updated_at = 2025-03-01T12:00:00Z`).
2. Query the API with a filter: `updated_since=<watermark>`.
3. Process and write the new records.
4. On successful completion, advance the watermark to the maximum `updated_at` seen in this batch.
5. Use an overlap window (e.g., subtract 5 minutes from the watermark) to catch late-arriving records and clock skew.

([Microsoft Tech Community](https://techcommunity.microsoft.com/blog/fasttrackforazureblog/robust-data-ingestion-with-high-watermarking/3707480), [dlt Docs](https://dlthub.com/docs/general-usage/incremental/cursor))

The target API must expose a reliable sort/filter field — `updated_at`, `created_at`, or a monotonic integer ID. Not all APIs do; check documentation before assuming incremental sync is possible.

### Idempotency
Network failures mid-run can cause partial writes. Make writes idempotent: use upserts on a stable primary key so that reprocessing the same records produces the same final state rather than duplicates.
([Medium/towards-data-engineering](https://medium.com/towards-data-engineering/building-idempotent-data-pipelines-a-practical-guide-to-reliability-at-scale-2afc1dcb7251))

---

## 9. Data Feed Types and Public Sources

### Feed format taxonomy

| Format | Common use | Notes |
|---|---|---|
| **REST JSON** | Most modern APIs | Human-readable; widely supported |
| **GraphQL** | Complex graph data; selective fields | Provider-dependent pagination |
| **RSS / Atom** | News, blog, content syndication | XML; well-standardized; date-sortable |
| **CSV flat file** | Bulk historical data, government datasets | Simple; no pagination; schema-fragile |
| **Streaming (SSE/WS)** | Financial ticks, IoT, live events | Low latency; requires persistent connection |

### Notable public API data sources for analysts

- **US Census Bureau API** — American Community Survey, decennial census, economic indicators; free; well-documented; requires registration for an API key. ([census.gov](https://www.census.gov/data/developers/data-sets.html))
- **Data.gov** — ~526,000 federal datasets as of early 2026; heterogeneous formats; reliability varies by agency. ([data.gov](https://data.gov/open-gov/))
- **NOAA / NWS** — weather and climate data; REST and bulk file downloads.
- **Financial market APIs** (Alpha Vantage, Financial Modeling Prep, EODHD) — equities, forex, fundamentals; free tiers rate-limited; paid tiers for intraday and bulk history. ([apilayer.com](https://blog.apilayer.com/12-best-financial-market-apis-for-real-time-data-in-2026/))
- **Social / news APIs** — increasingly restricted; Twitter/X deprecated free academic access; Reddit API repriced in 2023; factor licensing costs and access tier into the data collection plan.

---

## 10. Data Quality Pitfalls When Collecting via APIs

### Schema drift
Upstream teams modify their APIs (rename fields, drop attributes, change types) without always issuing a formal breaking-change notice. A `SELECT *` query against a REST response that gains a new column can silently corrupt downstream transformations.
([Monte Carlo](https://montecarlo.ai/blog-5-ways-to-stop-software-engineers-from-causing-data-quality-challenges/), [DataExpert.io](https://www.dataexpert.io/blog/backward-compatibility-schema-evolution-guide))

Mitigations:
- Define explicit field lists in extraction code; never expand `*`.
- Validate response schema at ingestion time against a stored contract; alert on deviation.
- Pin to a versioned API endpoint (`/v2/`) and monitor deprecation announcements.

### Partial reads and data completeness
Pagination bugs, auth expiry mid-run, or quota exhaustion can produce partial datasets that appear complete. Common symptoms: record counts that vary between runs, gaps in timestamps, missing foreign-key joins.

Mitigations:
- Record the total count returned by the API's metadata field (`total_count`, `x-total-count` header) and compare to rows actually stored.
- Implement a post-extraction assertion: query a known-stable count (e.g., records from 30 days ago) and alert if it changes.

### Deduplication
Retries, overlapping incremental windows, and webhook re-delivery all introduce duplicates. Guarantee uniqueness at write time with upserts on a stable primary key, or deduplicate during transformation using `ROW_NUMBER() OVER (PARTITION BY id ORDER BY updated_at DESC)`.
([OpenMeter](https://openmeter.io/blog/usage-deduplication), [python-bloggers.com](https://python-bloggers.com/2024/12/mastering-idempotency-in-data-analytics-ensuring-reliable-pipelines/))

### Timezone and clock skew
Incremental watermarks based on timestamps assume all parties use the same timezone and synchronized clocks. APIs may return UTC, local time, or ambiguous epoch seconds. Normalize to UTC at ingestion; add a buffer window to the watermark.

---

## 11. Choosing the Right Collection Method

| Situation | Recommended approach |
|---|---|
| Periodic reporting from a third-party SaaS | REST polling with incremental watermark; cursor pagination if available |
| Real-time event alerts (payments, sign-ups) | Webhook listener; fallback poll for missed events |
| Complex multi-entity joins in one call | GraphQL (if available from provider) |
| Large historical backfill | Bulk flat-file export + API for incremental updates going forward |
| Financial tick data | Streaming feed (WebSocket or SSE from exchange/vendor) |
| Government / census data | REST or CSV download; check update frequency before building a real-time pipeline |
