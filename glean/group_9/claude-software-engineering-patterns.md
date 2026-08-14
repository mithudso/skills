# software-engineering-patterns

**Category:** Security, Auth & Diagnostics
**Platform:** Claude
**Original Path:** claude/standalone/software-engineering-patterns

## Description
Language-agnostic software design/architecture/practice hub. TRIGGER: API design (REST/GraphQL/gRPC, versioning); backend & microservices patterns (decomposition, bounded contexts, strangler fig); Express/Fastify/NestJS/Hono selection; Python web frameworks (FastAPI/Django/Flask, ASGI/WSGI, uvicorn/gunicorn); design patterns (creational/structural/behavioral); coding standards; antipatterns & code smells; architecture docs (ISO 42010, ADRs, C4); code review (OWASP checklist); debugging & root-cause (5 Whys); performance profiling; SSE/streaming; job scheduling / MV3 alarms; auth-checker/web-auth; diagnostic registry; Vitest testing (vi mocks, coverage, snapshots, Jest→Vitest); Playwright web/UI testing; repo reusable-pattern scanning; automated program repair & code auto-remediation (GenProg, patch generation, fault localization, RepairAgent, Agentless, patch overfitting). SKIP: language syntax/idioms → lang-python/lang-js-ts/lang-go-and-mobile; frontend/UI/CSS → frontend-ui; SaaS API clients (Jira/Slack/Salesforce/Monday) → integration-clients.

---

> **Output rules:** No explanations — code only. Skip preamble. Don't recap, just proceed.

# Software Engineering Patterns

The hub for language-agnostic software design, architecture, and engineering
practice — the part of building software that is about *how to structure,
review, debug, and operate* a system rather than the syntax of any one
language. It spans API and backend/microservices design, the classic design
patterns, software-architecture methodology and documentation, code review,
debugging methodology and performance profiling, and the family of operational
patterns this workspace leans on (SSE/streaming, job/alarm scheduling,
diagnostic/ops/playbook registries, auth checks, templating, IndexedDB, and
integration-client scaffolding).

Use it when the task is about *design decisions, architecture, review, or
debugging methodology*. When the task is really about language syntax, the
frontend, AI/agent design, infrastructure, or a specific vendor's API, hand off
to the sibling hub named in the cross-hub note below.

## How to use this skill

This skill consolidates 25 sub-skills as on-demand reference files under
`references/`. Match the task to the routing table below and **Read the listed
`references/…md` file before answering deep questions** — the routing table
alone is not enough for depth. For framework-, language-, or vendor-exact
details, defer to the relevant sibling hub and the official docs as the source
of truth.

## Sub-skill routing table

This hub absorbs 25 former standalone skills as on-demand reference files. When a
task matches a row, **Read the listed `references/` file** before answering — do
not rely on this table alone for depth.

| Sub-topic | When to load | Reference file |
| --- | --- | --- |
| `api-design-patterns` | API design expert for REST, GraphQL, and gRPC. Covers resource modeling, Richardson | `references/api-design-patterns.md` |
| `backend-patterns` | Backend architecture patterns for Node.js, Express, and Next.js API routes. Covers | `references/backend-patterns.md` |
| `microservices-patterns` | Microservices architecture expert: service decomposition (bounded contexts, DDD, strangler fig), | `references/microservices-patterns.md` |
| `express-patterns` | Express.js production patterns expert — middleware, routing, error handling, rate limiting, CORS, security hardening, and testing for Express 4.x and 5.x. | `references/express-patterns.md` |
| `nodejs-backend-frameworks` | Fastify / NestJS / Hono patterns beyond Express — Fastify plugins+encapsulation+lifecycle hooks+JSON-Schema/TypeBox; NestJS modules/DI/providers + Middleware→Guards→Interceptors→Pipes→Filters pipeline + Fastify adapter; Hono Context/RegExpRouter/typed Bindings/zod-validator/RPC/edge runtimes; and Fastify-vs-NestJS-vs-Hono selection. | `references/nodejs-backend-frameworks.md` |
| `python-web-frameworks` | Python web frameworks — FastAPI (Depends DI, Pydantic v2, async/sync routes, lifespan, APIRouter), Django (MVT, ORM, async views, DRF, 6.0), Flask (app factory, blueprints, context), and the shared ASGI/WSGI + uvicorn/gunicorn/hypercorn server foundation; framework selection and event-loop-blocking diagnosis. | `references/python-web-frameworks.md` |
| `coding-patterns` | Software design patterns for JavaScript/TypeScript: creational (factory, builder, singleton, prototype), structural (facade, adapter, proxy, decorator… | `references/coding-patterns.md` |
| `coding-standards` | Baseline cross-project coding conventions for naming, readability, immutability, error handling, and code-quality review. Covers JavaScript, TypeScript, and… | `references/coding-standards.md` |
| `development-antipatterns` | Antipattern catalog with detection/remediation — classics (God object, spaghetti, big ball of mud, premature optimization), architecture (distributed monolith, microservice envy, lock-in), process (cargo cult, death march, bikeshedding), AI-era (AI-amplified duplication, slopsquatting, prompt-and-pray, complacency with AI code), the empirical evidence for/against smell catalogs, and when an antipattern is the rational choice | `references/development-antipatterns.md` |
| `debugging` | Systematic 7-phase debugging workflow for diagnosing and fixing bugs: symptom collection, environment analysis, unit test execution, integration testing… | `references/debugging.md` |
| `debugging-strategies` | Master-level debugging reference covering systematic methodology, root cause analysis (5 Whys, Fishbone, Fault Tree), profiling tools by language/runtime… | `references/debugging-strategies.md` |
| `software-architect` | Software architecture methodology and documentation expert — ISO/IEC/IEEE 42010 concepts, | `references/software-architect.md` |
| `code-reviewer` | Practical code-review reference covering reviewer mindset, workflow, comment quality, GitHub PR mechanics, and OWASP security checklists. Use when reviewing… | `references/code-reviewer.md` |
| `performance-profiling-expert` | Performance and profiling expert — Chrome DevTools Performance panel, Lighthouse audits, | `references/performance-profiling-expert.md` |
| `sse-streaming-patterns` | Server-Sent Events (SSE) patterns — EventSource API, fetch+ReadableStream | `references/sse-streaming-patterns.md` |
| `job-scheduling-patterns` | Node.js job scheduling expert: cron expression syntax, in-process scheduling (node-cron, | `references/job-scheduling-patterns.md` |
| `alarm-scheduler-patterns` | Chrome MV3 extension alarm scheduling expert: tiered severity-based polling, cooldowns, | `references/alarm-scheduler-patterns.md` |
| `auth-checker-patterns` | Multi-service auth state monitoring for Chrome MV3 extensions. Covers endpoint | `references/auth-checker-patterns.md` |
| `diagnostic-registry-patterns` | Diagnostic tool registry design -- static registries with keyword-based symptom matching, | `references/diagnostic-registry-patterns.md` |
| `ops-registry-patterns` | Operations registry patterns — dispatch tables, retry with backoff/jitter, circuit breakers, auto-remediation, idempotency, saga compensation, and audit… | `references/ops-registry-patterns.md` |
| `playbook-matcher-patterns` | Rule-based matching of support cases to diagnostic playbooks and KB articles — keyword scoring, | `references/playbook-matcher-patterns.md` |
| `template-config-patterns` | Configurable report and prompt template patterns for JavaScript operator tools, Chrome | `references/template-config-patterns.md` |
| `web-auth-patterns` | Browser authentication patterns — cookie attributes, OAuth 2.1 + PKCE, | `references/web-auth-patterns.md` |
| `indexeddb-patterns` | IndexedDB usage patterns for Chrome MV3 extensions and web apps — object stores, indexes, | `references/indexeddb-patterns.md` |
| `glean-llm-client-patterns` | Glean AI platform API integration patterns for JavaScript clients, Chrome MV3 | `references/glean-llm-client-patterns.md` |
| `salesforce-scraping-patterns` | Salesforce data extraction patterns for Chrome extensions — Lightning DOM traversal (LWC synthetic vs | `references/salesforce-scraping-patterns.md` |
| `testing-and-vitest-expert` | Vitest testing — vi mocks/spies, coverage (V8/Istanbul), snapshots, flaky-test diagnosis, vitest.config, Jest→Vitest migration | `references/testing-and-vitest-expert.md` |
| `webapp-testing` | Playwright local web-app testing & UI automation — DOM inspection, screenshots, console-log capture, multi-server setups | `references/webapp-testing.md` |
| `repo-pattern-scanner` | Scan a repository for reusable coding patterns and sync each to the shared pattern library | `references/repo-pattern-scanner.md` |
| `automated-program-repair` | Automated program repair & code auto-remediation — GenProg, patch generation, fault localization, RepairAgent, Agentless, patch-overfitting | `references/automated-program-repair.md` |

## Cross-hub boundaries

This hub owns *design, architecture, review, and debugging methodology* — the
language-agnostic engineering layer. Hand off when the task is really about
something a sibling hub owns:

- **Language syntax and idioms** — Python, Go, TypeScript, JS/Node, Kotlin/Compose
  → `programming-languages`.
- **Frontend, UI, CSS, component design** → `frontend-ui`.
- **AI / agent / LLM system design and prompt engineering** → `ai-agent-engineering`.
- **Infrastructure, CI/CD, containers, observability** → `devops-infra`.
- **Vendor API client specifics** (endpoints, auth flows, SDK quirks) →
  `integration-clients`.
- **Deep multi-pass review that also *applies* Medium+ fixes and verifies via build/lint/tests, looping a whole file or repo to convergence** → `code-deep-optimizer`. This hub's `references/code-reviewer.md` is single-pass review *guidance*; the apply-and-verify convergence loop is cdo's job.

**The patterns ↔ languages boundary (the most common overlap):** this hub holds
the **DESIGN, architecture, review, and debugging-*methodology*** — *should this
be one service or three? what does a clean retry/circuit-breaker look like? how
do I review this PR? what is a systematic way to isolate this bug?* The
**language-specific SYNTAX and idioms** that implement those decisions —
`asyncio.TaskGroup` mechanics, Go channel patterns, a TypeScript conditional
type, the exact `try/except` form — live in `programming-languages`. Lead with
the hub that matches intent: *how to structure or reason about it* stays here;
*how to spell it in language X* goes to `programming-languages`. Many tasks
touch both — start here for the design, cross-load `programming-languages` for
the implementation detail.

<!-- cross-hub-map -->
## Cross-hub map — where every software topic lives

This family is split across these hubs. If a task's deep material is **not** in this hub's Sub-skill
routing table, it is a reference file under a sibling hub below — **activate that hub or `Read` its
`references/<name>.md` directly**. Every former standalone skill in this family is now a reference under one
of these hubs (nothing was deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `programming-languages` | Programming Languages (Python, Go, TypeScript, JS/Node, Kotlin/Compose) | `references/python-patterns.md`, `references/go-patterns.md`, `references/javascript-nodejs.md`, `references/typescript-expert.md`, … |
| `software-engineering-patterns` | Software Engineering Patterns & Practices (architecture, APIs, debugging, reviews) | `references/api-design-patterns.md`, `references/backend-patterns.md`, `references/microservices-patterns.md`, `references/express-patterns.md`, … |