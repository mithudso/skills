<!-- hub-reference-banner -->
> **Reference file — part of the `integration-clients` hub.** Formerly the standalone `salesforce-developer-expert` skill.
> Sibling topics in this family are now reference files under the hubs (`integration-clients`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: salesforce-developer-expert
version: 1.1.0
category: developer
tags: [salesforce, apex, soql, lwc, api, flow-builder, bulk-api, platform-events, cdc, named-credentials, agentforce, ci-cd, governor-limits]
description: >-
  Salesforce platform development expert — APIs (REST, SOAP, Bulk 2.0, Composite, Tooling, Pub/Sub, GraphQL),
  Apex best practices (governor limits, bulkification, security, testing/mocking), LWC component architecture,
  SOQL optimization, Flow Builder automation, Platform Events and CDC, Salesforce DX CI/CD,
  Named Credentials, Agentforce/Einstein AI, Data Cloud, and deployment strategies.
  TRIGGER: implementing, reviewing, or troubleshooting Salesforce integrations, Apex code, platform automation,
  Service Cloud case CRUD, SOQL queries, governor limits, LWC components, or Salesforce DX pipelines.
  SKIP: Chrome extension scraping of Salesforce UI, shadow DOM traversal on Lightning pages, or content-script
  injection into Salesforce — use salesforce-scraping-patterns for those.
triggers:
  - "Write an Apex trigger for bulk record updates"
  - "How do I call the Salesforce REST API?"
  - "Bulk load records into Salesforce using Bulk API 2.0"
  - "What are the SOQL governor limits?"
  - "Build an LWC component to display account data"
  - "Set up Platform Events or CDC for real-time updates"
  - "Migrate from sfdx CLI to sf CLI"
  - "Secure a callout using Named Credentials"
  - "Flow builder before-save vs after-save trigger"
  - "SOQL query optimization for large data volumes"
  - "Deploy Salesforce metadata with scratch orgs"
  - "Apex security: stripInaccessible vs USER_MODE"
  - "Set up Agentforce or Einstein AI in Salesforce"
  - "Salesforce Composite API multi-step operation"
related_skills:
  - salesforce-scraping-patterns
origin:
  - docs/salesforce-developer-context.md
  - https://developer.salesforce.com/docs/apis
sources:
  - https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/resources_composite_batch.htm
  - https://developer.salesforce.com/docs/atlas.en-us.api_asynch.meta/api_asynch/bulk_api_2_0.htm
  - https://developer.salesforce.com/docs/atlas.en-us.apexcode.meta/apexcode/apex_gov_limits.htm
  - https://developer.salesforce.com/docs/platform/lwc/guide/get-started-best-practices.html
  - https://developer.salesforce.com/docs/atlas.en-us.change_data_capture.meta/change_data_capture/cdc_subscribe_pubsub_api.htm
  - https://architect.salesforce.com/decision-guides/event-driven
  - https://sfdcdevelopers.com/2026/05/06/salesforce-apis-complete-guide/
  - https://salesforcestack.com/optimizing-salesforce-soql-queries-best-practices-for-performance-and-limits/
  - https://salesforcemonday.com/2026/02/19/salesforce-flow-best-practices-2026/
  - https://www.salesforceben.com/12-salesforce-apex-best-practices/
  - https://blog.beyondthecloud.dev/blog/apex-security-and-sharing
  - https://www.salesforceben.com/salesforce-scratch-orgs/
  - https://www.echots.com/post/modern-salesforce-integrations-why-external-credentials-named-credentials-are-now-essential
  - https://marketingagent.blog/2026/03/25/salesforce-agentforce-complete-guide-deploy-ai-agents-in-2026/
  - https://xillentech.com/zero-copy-data-cloud-guide/
  - https://trailhead.salesforce.com/content/learn/modules/apex_testing/apex_testing_intro
  - https://help.salesforce.com/s/articleView?id=000387172&language=en_US&type=1
  - https://medium.com/@saurabh.samirs/lwc-best-practices-in-2025-performance-optimization-anti-patterns-to-avoid-25c315a38202
  - https://apex.kingsleymgbams.hashnode.dev/salesforce-summer-26-api-67-apex-user-mode-database-operations
---

# Salesforce Developer Expert

## When to use this skill

Use this skill when the task involves Salesforce development, including:

- choosing between Salesforce APIs
- OAuth and authenticated integration design
- Platform Events, Change Data Capture, Pub/Sub API, Streaming API, or Outbound Messaging
- webhook-equivalent integration patterns
- Service Cloud case retrieval, creation, update, upsert, delete, or search
- reports, dashboards, Analytics API, and case reporting
- Apex, triggers, sharing, CRUD/FLS, `stripInaccessible`, `USER_MODE`, testing, or governor limits
- Lightning Web Components design, performance, and data access
- SOQL/SOSL query optimization and large data volume handling
- Flow Builder automation (record-triggered, scheduled, screen, autolaunched)
- Salesforce DX, scratch orgs, CI/CD pipelines, and deployment strategies
- Named Credentials, External Credentials, and secure callout patterns
- Agentforce, Einstein AI, Data Cloud, and zero-copy integrations
- Bulk data processing, migration, and ETL patterns

## Skill guidance

1. Read `docs/salesforce-developer-context.md` (repo-relative) first and treat it as the primary retrieval context.
2. Prefer official Salesforce platform behavior over habit or analogy from other SaaS tools.
3. Choose the narrowest correct Salesforce surface:
   - REST for ordinary CRUD/query/search
   - Bulk API 2.0 for high-volume jobs
   - SOAP API when WSDL-oriented enterprise integrations matter
   - Analytics API for reports/dashboards
   - Platform Events / CDC / Pub/Sub API for realtime subscriptions
4. Be explicit about the webhook limitation:
   - Salesforce has event-driven outbound patterns, but not one universal native outbound HTTP webhook feature.
5. For case questions, distinguish:
   - raw Case object CRUD
   - case search/query
   - case reporting
   - case automation side effects such as assignment or escalation behavior
6. For Apex guidance, always preserve:
   - bulk-safe design
   - sharing model awareness
   - CRUD/FLS enforcement
   - `stripInaccessible`
   - `USER_MODE`
   - testing and governor-limit awareness
7. If exact payload shapes or headers matter and the context file does not include them, point back to the linked official guide rather than fabricating details.

## Practical defaults

- If the user says **webhook**, explain the tradeoffs among Outbound Messaging, Platform Events, CDC, Pub/Sub API, and Apex/Flow callouts.
- If the user says **case report** or **dashboard**, start with Analytics API instead of SOQL.
- If the user says **realtime subscription**, start with CDC or Platform Events and explain when Pub/Sub API is preferable.
- If the user says **update Salesforce data from code**, start with REST API unless scale or enterprise SOAP constraints change the answer.

## API selection guide

Salesforce exposes nine primary API surfaces. Select based on data volume, integration style, and latency needs.

| API | Protocol | Best for | Daily limit pool |
|-----|----------|----------|------------------|
| REST API | JSON/HTTPS | Standard CRUD, queries, search (<10k records) | Shared daily quota |
| SOAP API | XML/WSDL | Enterprise Java/.NET clients needing strict typing | Shared daily quota |
| Bulk API 2.0 | REST + CSV | Async ingest/query of 10k-150M records per job | Separate bulk pool |
| Composite API | REST | Multi-step operations in a single round trip | 1 API call per composite request |
| Tooling API | REST/SOAP | IDE integration, single-record metadata edits | Shared daily quota |
| Metadata API | SOAP/XML | Full package deployments between orgs | Separate metadata pool |
| Streaming / Pub/Sub API | CometD / gRPC | Real-time push notifications and event subscriptions | Event delivery limits |
| Connect REST API | JSON | Chatter, Experience Cloud, file management | Shared daily quota |
| Apex REST | JSON | Custom endpoints built in Apex | Shared daily quota |

### Daily API call quotas by edition

- **Developer Edition**: 15,000 calls/day
- **Enterprise Edition**: ~1,000,000 baseline + per-user scaling
- **Unlimited Edition**: ~10x Enterprise capacity

### API version retirement

API versions 21.0-30.0 were retired in June 2025. Integrations pinned below v31.0 will fail. Current version as of Summer 2026 is v67.0.

## Composite API patterns

Salesforce offers four composite resource types, each counted as a single API call:

### Composite (subrequest chaining)
- Up to 25 subrequests with `referenceId`-based chaining between them
- Maximum 5 queries or sObject collections per call
- Use `allOrNone: true` for transactional behavior
- Governor limits accumulate across ALL subrequests in the composite call

### Composite Graph
- Up to 500 subrequests organized into graphs
- Each graph is its own transaction
- Best for complex dependency trees that exceed the 25-subrequest limit

### Composite Batch
- Up to 25 independent subrequests executed sequentially
- No cross-reference between subrequests
- Simpler than Composite when operations are independent

### sObject Collections
- Up to 200 records per DML operation (insert, update, upsert, delete)
- Single-object batch operations
- Prefer over individual record calls for any multi-record DML

## Bulk API 2.0

### Ingest jobs
- Upload CSV data for insert, update, upsert, or delete
- Job lifecycle: `Open` -> upload data -> `UploadComplete` -> processing -> `JobComplete` or `Failed`
- Handles up to 150 million records per job
- Automatic batching and parallel processing by the platform
- Poll job status via `GET /jobs/ingest/{jobId}`

### Query jobs
- Asynchronous SOQL execution for large result sets (2,000+ records)
- Results returned as CSV via locator-based pagination
- Use for any query expected to return more than 2,000 rows
- Separate rate limit pool from REST API

### Best practices
- Use `columnDelimiter` and `lineEnding` parameters for non-standard CSV formats
- Monitor `numberRecordsProcessed` and `numberRecordsFailed` in job info
- Retrieve failed and unprocessed records for retry logic
- Prefer Bulk API 2.0 over legacy Bulk API (v1) for all new integrations

## Apex best practices

### Bulkification patterns

Triggers can process up to 200 records per invocation. All Apex code must be bulk-safe:

1. **Never SOQL/DML inside loops** -- collect IDs first, query once, process in collections
2. **Use Maps for lookups** -- `Map<Id, SObject>` from query results for O(1) access
3. **Aggregate DML** -- build lists of records, perform single insert/update/delete outside loops
4. **One trigger per object** -- delegate to handler classes using a trigger framework
5. **Recursion control** -- use a static Boolean flag in the handler to prevent re-entry

### Trigger handler pattern

```apex
// AccountTrigger.trigger
trigger AccountTrigger on Account (before insert, before update, after insert, after update) {
    AccountTriggerHandler.handle(Trigger.operationType, Trigger.new, Trigger.oldMap);
}

// AccountTriggerHandler.cls
public class AccountTriggerHandler {
    private static Boolean hasRun = false;
    public static void handle(TriggerOperation op, List<Account> newList, Map<Id, Account> oldMap) {
        if (hasRun) return;
        hasRun = true;
        switch on op {
            when BEFORE_INSERT { /* ... */ }
            when AFTER_UPDATE  { /* ... */ }
        }
    }
}
```

### Asynchronous Apex

| Mechanism | Max records | Use case | Governor reset |
|-----------|-------------|----------|----------------|
| `@future` | N/A (method args only) | Simple async callouts, non-urgent DML | Fresh limits per invocation |
| `Queueable` | Chainable | Complex async with state, job chaining | Fresh limits per execute() |
| `Batch Apex` | Up to 50M records | Large-volume processing with chunking | Fresh limits per execute() batch |
| `Scheduled Apex` | N/A | Cron-based periodic jobs | Fresh limits per execute() |

## Governor limits reference

Per-transaction limits enforce multi-tenant safety. Exceeding any limit throws an uncatchable `System.LimitException`.

| Resource | Synchronous limit | Asynchronous limit |
|----------|-------------------|-------------------|
| SOQL queries | 100 | 200 |
| Records retrieved by SOQL | 50,000 | 50,000 |
| DML statements | 150 | 150 |
| DML rows | 10,000 | 10,000 |
| CPU time | 10,000 ms | 60,000 ms |
| Heap size | 6 MB | 12 MB |
| Callouts (HTTP/Web service) | 100 | 100 |
| Callout timeout (single) | 120 seconds | 120 seconds |
| Total callout timeout | 120 seconds | 120 seconds |
| Future calls | 50 | 0 (not from async) |
| Queueable jobs added | 50 | 1 |
| Email invocations | 10 | 10 |
| SOSL queries | 20 | 20 |
| Describe calls | 100 | 100 |

### Runtime monitoring

Use the `Limits` class to check consumption mid-transaction:
```apex
System.debug('SOQL queries used: ' + Limits.getQueries() + ' of ' + Limits.getLimitQueries());
System.debug('DML statements used: ' + Limits.getDmlStatements() + ' of ' + Limits.getLimitDmlStatements());
System.debug('CPU time used: ' + Limits.getCpuTime() + ' of ' + Limits.getLimitCpuTime());
```

## Apex security model

### Sharing keywords

| Keyword | Record-level security | Use when |
|---------|----------------------|----------|
| `with sharing` | Enforced -- user sees only their records | Default for most classes; respects org-wide defaults, roles, sharing rules |
| `without sharing` | Not enforced -- all records accessible | System-level operations, background jobs, integrations |
| `inherited sharing` | Inherits from caller | Utility classes that should respect the caller's context |

### Field-level security enforcement

#### stripInaccessible (API v49.0+)
Graceful degradation -- strips fields the user cannot access and lets execution continue:
```apex
List<Account> accounts = [SELECT Id, Name, AnnualRevenue FROM Account];
SObjectAccessDecision decision = Security.stripInaccessible(AccessType.READABLE, accounts);
List<Account> sanitized = decision.getRecords();
// AnnualRevenue is removed if user lacks FLS
```

#### USER_MODE (API v49.0+, default in v67.0)
Enforces object-level, field-level, AND sharing rules on SOQL/DML:
```apex
List<Account> accounts = [SELECT Id, Name FROM Account WITH USER_MODE];
Database.insert(newAccounts, AccessLevel.USER_MODE);
```

#### API v67.0 change (Summer 2026)
User mode becomes the default for plain SOQL, SOSL, and DML operations. Existing code using `WITH SECURITY_ENFORCED` or `stripInaccessible` continues to work but `USER_MODE` is the recommended path forward.

### Security checklist for Apex code reviews
1. Every class declares an explicit sharing keyword
2. All SOQL/DML uses `USER_MODE` or `stripInaccessible` for FLS
3. No hard-coded IDs (RecordType, Profile, etc.) -- query dynamically
4. SOQL injection prevented via bind variables, never string concatenation
5. Callout credentials stored in Named Credentials, never in code
6. Sensitive data not written to debug logs in production

## SOQL optimization

### Selectivity and indexing

A query is **selective** when its filters target indexed fields, allowing the optimizer to use index-based access instead of full table scans.

**Standard indexed fields**: Id, Name, OwnerId, CreatedDate, SystemModstamp, RecordTypeId, lookup/master-detail fields, external ID fields.

**Custom indexes**: request via Salesforce Support for high-volume custom fields.

### Optimization rules

1. **Filter on indexed fields** in WHERE clauses
2. **Avoid functions on indexed fields** (e.g., `CALENDAR_YEAR(CreatedDate)` prevents index use)
3. **Avoid negative operators** on large tables (`!=`, `NOT LIKE`, `EXCLUDES` cause full scans)
4. **No leading wildcards** in LIKE clauses (`LIKE '%test'` is non-selective)
5. **Use LIMIT** to cap result size when full result set is not needed
6. **Use relationship queries** to fetch parent/child records in one query instead of multiple
7. **Use ORDER BY on indexed fields** to avoid sort overhead
8. **Use Query Plan Tool** in Developer Console to analyze cost of different query plans
9. **Selective field retrieval** -- list only needed fields, never use dynamic `SELECT *` equivalents

### Large data volume patterns

- For datasets exceeding 50,000 rows, use date-range filters or Batch Apex to chunk processing
- Use `Database.QueryLocator` in Batch Apex for queries up to 50 million rows
- Consider skinny tables (request via Support) for frequently queried wide objects
- Archive old records to reduce working set size

## Lightning Web Components (LWC)

### Architecture principles

1. **Small, focused components** -- each component serves a single purpose for reusability
2. **Unidirectional data flow** -- parent passes data down via `@api`, child communicates up via custom events
3. **Reactive data binding** -- `@wire` automatically retrieves and refreshes data when reactive parameters change
4. **Lightning Data Service (LDS)** first -- use for standard CRUD with automatic security enforcement, caching, and offline support; fall back to Apex for complex logic

### Data access patterns

| Pattern | When to use | Security |
|---------|------------|----------|
| Lightning Data Service (LDS) | Standard CRUD, single record | Automatic FLS, sharing, caching |
| `@wire` with Apex | Complex queries, aggregations | Manual FLS enforcement in Apex |
| Imperative Apex calls | User-initiated actions, conditional logic | Manual FLS enforcement in Apex |
| `lightning/uiGraphQLApi` (deprecated) | Complex related data | Automatic |
| `lightning/graphql` (current) | Complex related data, GA in Winter '26 | Automatic |

### Performance optimization

1. **Use `lwc:if` over CSS `display: none`** -- completely removes elements from DOM, preventing wire calls and lifecycle hooks
2. **Lazy load data** -- fetch only when needed, not on component initialization
3. **Infinite scrolling** -- load subsets with progressive fetching for large lists
4. **Debounce user input** -- avoid firing wire/apex calls on every keystroke
5. **Minimize reactive property changes** -- each change triggers re-render; batch updates when possible

### LWC anti-patterns
- Monolithic components with hundreds of lines
- Direct DOM manipulation instead of reactive data binding
- Imperative Apex calls where `@wire` would suffice
- Inline styles instead of SLDS utility classes or custom CSS
- Tight coupling between parent and child via shared state objects

### Recent LWC updates (2025-2026)
- **SLDS 2.0 GA** -- supports dark mode and custom themes
- **`lightning/graphql`** -- replaces `lightning/uiGraphQLApi` for complex data queries
- **Lightning Web Security (LWS) Trusted Mode** -- allows third-party scripts to navigate security restrictions
- **Aura phaseout** -- Salesforce actively migrating to LWC-only; avoid new Aura components

## Event-driven architecture

### Platform Events
- Custom event definitions you create (`MyEvent__e`)
- Publish from Apex, Flow, or external systems via REST
- Subscribe in Apex triggers, Flow, or external clients via Pub/Sub API
- High-volume events: up to 25M events/day (varies by edition)
- Replay: up to 72 hours of event replay via ReplayId
- Delivery: at-least-once; consumers must be idempotent

### Change Data Capture (CDC)
- Automatic change tracking on enabled standard/custom objects
- Channel format: `/data/{ObjectName}ChangeEvent`
- Emits CREATE, UPDATE, DELETE, UNDELETE, GAP events
- **Gap events**: header-only messages when payload exceeds 1 MB or internal cleanup occurs; best practice is to call back with record IDs to retrieve current state
- Use for data synchronization, audit logging, cache invalidation, and data lake streaming

### Pub/Sub API (recommended for new builds)
- gRPC-based bidirectional streaming
- Single API for publish, subscribe, and schema retrieval
- Higher throughput and lower latency than CometD-based Streaming API
- Supports both Platform Events and CDC channels
- Managed subscriptions with flow control (requested event count)
- Event replay from stored ReplayId

### When to use which

| Need | Recommended surface |
|------|-------------------|
| Custom business event signaling | Platform Events |
| Record change notifications | CDC |
| External system subscription | Pub/Sub API |
| Simple push notifications (legacy) | Streaming API (PushTopic/Generic) |
| Outbound HTTP to a single endpoint | Outbound Messaging or Apex/Flow callout |

### Streaming API deprecation note
CometD-based Streaming API (PushTopics, Generic Events) is in maintenance mode. Use Pub/Sub API for all new event-driven integrations.

## Flow Builder automation

### Flow types

| Type | Trigger | Use case |
|------|---------|----------|
| Record-triggered (before save) | Record create/update | Fast field updates on the same record (no extra DML) |
| Record-triggered (after save) | Record create/update/delete | Related record updates, callouts, event publishing |
| Scheduled-triggered | Time-based (cron) | Periodic batch processing |
| Screen flow | User interaction | Guided wizards, multi-step forms |
| Autolaunched | Invoked from Apex, other flows, or REST | Reusable sub-processes |
| Platform Event-triggered | Event received | Event-driven automation |

### Flow best practices (2026)

1. **Before-save for same-record updates** -- never use after-save to update the triggering record (wastes a DML operation)
2. **Avoid DML inside loops** -- Salesforce auto-bulkifies some elements, but DML in loops breaks this optimization
3. **Naming convention**: `ObjectName_TriggerContext_Purpose` (e.g., `Account_AfterUpdate_SyncToERP`)
4. **Centralized error handling** -- create one error-handling subflow that all flows call; log errors and send admin notifications
5. **Use sub-flows** for reusable logic blocks -- build once, invoke from multiple parent flows
6. **Flow Tests** -- native Flow testing (GA) provides persistent test scenarios for record-triggered flows; runs alongside Apex tests in CI/CD
7. **Entry conditions** -- always set precise entry conditions to avoid unnecessary executions
8. **Process Builder and Workflow Rules** -- reached end of support December 31, 2025; migrate all automation to Flow Builder

### Spring 2026 Flow features
- AI-powered Flow building via Agentforce (natural language to Flow)
- Enhanced debugging with detailed flow interview logs
- Improved sub-flow parameter passing

## Salesforce DX and CI/CD

### CLI migration
The `sfdx` CLI is deprecated. Use `sf` commands:
- `sfdx force:source:deploy` -> `sf project deploy start`
- `sfdx force:apex:test:run` -> `sf apex test run`
- `sfdx force:org:create` -> `sf org create scratch`
- `sfdx force:source:deploy --checkonly` -> `sf project deploy validate`

### Scratch orgs
- Ephemeral orgs for isolated development and testing (max 30-day lifespan)
- Defined by `project-scratch-def.json` with edition, features, and settings
- CI-friendly: create via CLI, run tests, tear down automatically
- Use for feature branches; sandboxes for UAT and staging

### Deployment pipeline pattern

```
Feature branch -> Scratch org (dev + unit tests)
         |
    Pull request -> CI scratch org (validation deploy + all tests)
         |
    Merge to main -> Deploy to Integration sandbox
         |
    Release branch -> Deploy to UAT sandbox
         |
    Release tag -> Deploy to Production (with RunLocalTests)
```

### CI/CD tools
- **Salesforce CLI** (`sf`) -- core deployment and testing commands
- **GitHub Actions / GitLab CI** -- pipeline orchestration
- **PMD for Apex** -- static analysis for code quality
- **SFDX Scanner** -- Salesforce-specific static analysis
- **Gearset / Copado / AutoRABIT** -- managed DevOps platforms with metadata comparison and rollback

### Deployment best practices
1. Always run `sf project deploy validate` before production deploys
2. Use `RunLocalTests` test level for production deployments (skips managed package tests)
3. Track metadata in version control (source format, not MDAPI)
4. Use `.forceignore` to exclude unnecessary metadata from deployments
5. Destructive changes require a separate `destructiveChanges.xml` manifest

## Named Credentials and External Services

### Modern architecture (External Credentials + Named Credentials)

The current model separates concerns:
- **External Credential**: encapsulates authentication config (OAuth flow, API keys, JWT, AWS SigV4, mTLS)
- **Principal**: defines which Permission Sets can use the credential
- **Named Credential**: binds an endpoint URL to an External Credential

### Supported auth protocols
- OAuth 2.0 (Authorization Code, Client Credentials, JWT Bearer)
- AWS Signature Version 4
- Basic Authentication
- API Key (header or query parameter)
- Mutual TLS (client certificate)
- Custom authentication headers

### Benefits
- Credentials never appear in Apex code, metadata, or version control
- Salesforce handles token refresh automatically
- One External Credential can serve multiple Named Credentials (e.g., Google Drive + Google Calendar share one OAuth config)
- Reusable across Apex, Flows, External Services, and HTTP callouts

### External Services
- Import OpenAPI/Swagger specs to auto-generate invocable actions
- Actions are usable in Flow Builder without writing Apex
- Automatically respects Named Credential authentication
- Use for straightforward REST API integrations where custom Apex logic is unnecessary

## Apex testing patterns

### Test class structure
```apex
@IsTest
private class AccountServiceTest {
    @TestSetup
    static void setupData() {
        // Create shared test data once; rolled back after all test methods
        insert new Account(Name = 'Test Corp');
    }

    @IsTest
    static void shouldCalculateRevenue() {
        // Arrange
        Account acc = [SELECT Id FROM Account LIMIT 1];
        // Act
        Test.startTest();
        Decimal result = AccountService.calculateRevenue(acc.Id);
        Test.stopTest();
        // Assert
        Assert.areEqual(100000, result, 'Revenue should match expected value');
    }
}
```

### Key testing practices

1. **AAA pattern** (Arrange-Act-Assert) for every test method
2. **`@TestSetup`** to create shared data once, reducing test execution time
3. **`Test.startTest()` / `Test.stopTest()`** to reset governor limits and force async execution
4. **Never use `@SeeAllData(true)`** -- creates dependency on org data, causes inconsistent failures
5. **No hard-coded IDs** -- query RecordTypes, Profiles, and other metadata dynamically
6. **Positive, negative, and bulk tests** -- test happy path, error conditions, and 200+ record batches
7. **Test naming**: `ClassName` + `Test` suffix (e.g., `AccountServiceTest`)
8. **75% minimum coverage** required for production deployment; aim for 90%+

### Mocking patterns

#### HttpCalloutMock
```apex
@IsTest
private class ExternalServiceTest {
    private class MockHttp implements HttpCalloutMock {
        public HttpResponse respond(HttpRequest req) {
            HttpResponse res = new HttpResponse();
            res.setStatusCode(200);
            res.setBody('{"status":"ok"}');
            return res;
        }
    }

    @IsTest
    static void shouldHandleCallout() {
        Test.setMock(HttpCalloutMock.class, new MockHttp());
        Test.startTest();
        String result = ExternalService.callApi();
        Test.stopTest();
        Assert.areEqual('ok', result);
    }
}
```

#### Stub API for dependency injection
```apex
@IsTest
private class OrderProcessorTest {
    private class MockInventory implements StubProvider {
        public Object handleMethodCall(Object stubbedObject, String methodName,
            Type returnType, List<Type> paramTypes, List<String> paramNames, List<Object> args) {
            if (methodName == 'checkStock') return true;
            return null;
        }
    }

    @IsTest
    static void shouldProcessOrder() {
        InventoryService mockSvc = (InventoryService) Test.createStub(InventoryService.class, new MockInventory());
        OrderProcessor processor = new OrderProcessor(mockSvc);
        Test.startTest();
        Boolean result = processor.process(new Order());
        Test.stopTest();
        Assert.isTrue(result);
    }
}
```

### What to mock
- HTTP callouts (always mock -- callouts fail in test context without mock)
- External service dependencies via Stub API
- Complex data scenarios where DML setup would be expensive
- System.now() via injectable clock pattern for time-dependent logic

## Agentforce and Einstein AI

### Agentforce (GA October 2025)
- Autonomous AI agents that execute multi-step workflows
- Built on Data Cloud + Einstein Trust Layer
- **Agent Builder**: low-code tool for creating custom agents with topics, instructions, and actions
- Agents can use Apex, Flow, Prompt Templates, and API actions
- 12,000+ customers as of 2026

### Einstein AI capabilities
- **Predictive**: lead scoring, opportunity insights, forecasting
- **Generative**: email generation, case summarization, knowledge article drafting
- **Einstein Copilot**: conversational assistant embedded in Lightning Experience

### Data Cloud
- Unified data platform ingesting data from any source
- 50+ trillion records managed (as of 2026)
- **Zero-Copy Partner Network**: query data in Snowflake, Databricks, BigQuery, Azure without replication
- Data harmonization via Data 360 Data Model Objects
- Foundation for Agentforce agent grounding and retrieval

### Einstein Trust Layer
- Prompt defense and toxicity filtering
- Data masking for PII before sending to LLMs
- Audit trail of all AI interactions
- Zero data retention policy with third-party LLM providers

## Anti-patterns to flag

| Anti-pattern | Why it is harmful | Fix |
|-------------|-------------------|-----|
| SOQL/DML inside loops | Hits governor limits at scale | Bulkify: collect, query once, process in collections |
| Hard-coded IDs | Breaks across orgs/sandboxes | Query metadata dynamically |
| No sharing keyword on class | Ambiguous security posture | Always declare `with sharing`, `without sharing`, or `inherited sharing` |
| `@SeeAllData(true)` in tests | Flaky, org-dependent tests | Create test data programmatically |
| String concatenation in SOQL | SOQL injection vulnerability | Use bind variables |
| After-save flow updating same record | Wastes DML and can cause recursion | Use before-save flow for same-record updates |
| Hardcoded credentials in Apex | Security risk, fails security review | Use Named Credentials |
| SELECT * equivalent queries | Retrieves unnecessary data, wastes resources | List only needed fields |
| Monolithic LWC components | Unmaintainable, poor reusability | Break into focused, single-purpose components |
| Using Aura for new components | Deprecated path, missing modern features | Use LWC exclusively for new development |
| CometD Streaming API for new builds | Maintenance mode, lower throughput | Use Pub/Sub API (gRPC) |
| Process Builder / Workflow Rules | End of support Dec 2025 | Migrate to Flow Builder |

## Quick-reference decision trees

### "How should I integrate with Salesforce?"

```
Need to read/write < 10k records? --> REST API
Need to process > 10k records?   --> Bulk API 2.0
Need real-time change events?    --> Pub/Sub API + CDC
Need to chain multiple operations? --> Composite API
Need strongly-typed WSDL?        --> SOAP API
Need to deploy metadata?         --> Metadata API / sf CLI
Need custom endpoint in SF?      --> Apex REST
```

### "How should I automate this?"

```
Same-record field update?         --> Before-save record-triggered Flow
Related record update?            --> After-save record-triggered Flow
Needs complex business logic?     --> Apex trigger + handler class
Scheduled periodic job?           --> Scheduled Flow or Scheduled Apex
User-facing wizard?               --> Screen Flow
Event-driven reaction?            --> Platform Event-triggered Flow
Large data batch processing?      --> Batch Apex
```

### "How should I enforce security?"

```
Need FLS on SOQL?                 --> WITH USER_MODE (preferred) or stripInaccessible
Need FLS on DML?                  --> Database.insert(records, AccessLevel.USER_MODE)
Need record-level security?       --> with sharing keyword on class
Need to skip sharing for system job? --> without sharing (document justification)
Utility/helper class?             --> inherited sharing
Storing external credentials?     --> Named Credentials (never hardcode)
```
