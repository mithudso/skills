<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Authored 2026-05-31 from web research (`/dr`).
> Sibling topics in this family are reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" pointers that name a bare sibling; load that topic's `references/<name>.md` from the owning hub.
> For the Express analog (middleware, routing, error handling) see `references/express-patterns.md` — that reference owns Express specifics and is the shared HTTP-framework foundation this file builds on. For Node language/runtime idioms see `programming-languages/references/javascript-nodejs.md`, `programming-languages/references/nodejs-concurrency-internals.md`, and `programming-languages/references/javascript-runtimes-deno-bun-edge.md` (Hono's edge/Bun/Deno targets). For TS schema validation depth see `zod-schema-validation`.

---
name: nodejs-backend-frameworks
version: "1.0.0"
updated: "2026-05-31"
category: developer
tags: [fastify, nestjs, hono, nodejs, typescript, backend, middleware, dependency-injection, edge, json-schema, zod, typebox, rpc, framework-selection]
description: "Node.js/TypeScript backend framework patterns beyond Express — Fastify, NestJS, and Hono. TRIGGER: building or reviewing a Fastify service (plugins, encapsulation context, fastify-plugin, lifecycle hooks onRequest→preHandler→onSend, JSON Schema validation + serialization, TypeBox/type-provider, decorators); a NestJS app (modules, hierarchical DI, providers/scopes, the Middleware→Guards→Interceptors→Pipes→Handler→Filters request pipeline, dynamic modules, forwardRef circular deps, Fastify adapter); a Hono app (Context object c, web-standard Request/Response, RegExpRouter, typed middleware/Variables/Bindings, @hono/zod-validator, RPC mode + hc client, edge runtimes — Cloudflare Workers/Deno/Bun/Lambda); or choosing among Fastify vs NestJS vs Hono (and vs Express). SKIP: Express.js specifics → references/express-patterns.md; Python FastAPI/Django/Flask → references/python-web-frameworks.md; generic REST/GraphQL/gRPC surface design → references/api-design-patterns.md; Node language/event-loop/runtime idioms → programming-languages; standalone Zod schema depth → zod-schema-validation."
related_skills: [express-patterns, api-design-patterns, backend-patterns, web-auth-patterns, zod-schema-validation, javascript-nodejs]
---

# Node.js Backend Frameworks: Fastify, NestJS, and Hono

Three post-Express Node.js/TypeScript backend frameworks, each occupying a distinct point on the same spectrum. They share the HTTP-handler foundation already captured in `references/express-patterns.md` (middleware chains, routing, error handling, graceful shutdown, security hardening) — **read that first; this file does not repeat it.** What follows is what each framework does *differently* from Express and from each other, then how to choose.

> Cross-references: Express baseline → `references/express-patterns.md`. API surface design (REST/GraphQL/gRPC, versioning) → `references/api-design-patterns.md`. Backend service architecture (caching, queues, DB) → `references/backend-patterns.md`. Auth flows (OAuth 2.1/PKCE, JWT, CSRF, cookies) → `references/web-auth-patterns.md`. Zod/TypeBox/Valibot schema depth → `zod-schema-validation`. Edge/Bun/Deno runtime constraints → `programming-languages/references/javascript-runtimes-deno-bun-edge.md`.

## 0. The spectrum — pick from the workload, not the hype

| Framework | Weight | Core identity | Runtimes |
| --- | --- | --- | --- |
| **NestJS** | Heavyweight | Opinionated architecture, DI container, enterprise modules (decorator/Angular-style) | Node only; uses Express **or** Fastify as the HTTP adapter; **no** edge runtimes |
| **Fastify** | Mid-weight | Express-like ergonomics + raw speed + built-in JSON-Schema validation; minimal opinions | Node only (Node-specific APIs) |
| **Hono** | Lightweight | ~14KB, Web-Standards (`Request`/`Response`) based, edge-native, tight type inference | Cloudflare Workers, Deno, Bun, AWS Lambda, Fastly, Node — multi-runtime |

Rough JSON-serialization throughput (synthetic, not real-world): Hono ~78k req/s, Fastify ~62k req/s. **Raw framework speed rarely decides real apps — DB queries and business logic dominate.** Choose on architecture, runtime target, and team conventions, not the benchmark.

## 1. Fastify

### 1.1 Plugins + encapsulation (the central concept)
Everything in Fastify is a **plugin**, and the defining feature is the **encapsulation context**: it governs which decorators, hooks, and child plugins a route can see. A plugin registered with `fastify.register(...)` creates a *child* context — its decorators/hooks/schemas are scoped to that subtree and its siblings/parent cannot see them. This replaces Express's "every request flows through one global middleware chain": hooks run only where you register them.

- **`fastify-plugin` wrapper** breaks encapsulation *on purpose* — wrap a plugin with `fp(...)` when its decorators/hooks must be visible to the parent and siblings (e.g. a shared DB connection, an auth decorator). The classic bug: "my decorator isn't defined on the parent instance" → you forgot `fastify-plugin`. `fp` also lets you declare a supported Fastify version range.
- Use plain (encapsulated) plugins for feature slices (a router subtree); use `fp`-wrapped plugins for cross-cutting capabilities.

### 1.2 Lifecycle hooks
Request hooks fire in a fixed order, and (except `onClose`) **all hooks are encapsulated** — register them where they should apply:

```
onRequest → preParsing → preValidation → preHandler → [route handler] → preSerialization → onSend → onResponse
```
Plus application hooks (`onReady`, `onClose`, `onRoute`, `onRegister`) and error hooks. `onError` is **read-only for logging** — to change an error response use `setErrorHandler()`, not `onError`.

### 1.3 JSON Schema validation + serialization (Fastify's superpower)
Attach a `schema` (`body`/`querystring`/`params`/`headers`/`response`) to a route. Fastify compiles it once: input is validated, and **response serialization is compiled from the response schema** (via `fast-json-stringify`) — this is a large part of why Fastify is fast and is something Express has no equivalent for. `fastify.addSchema()` registers reusable shared schemas (`$ref`-able); this API is **encapsulated** too.

- **TypeBox + `@fastify/type-provider-typebox`** (or the JSON-Schema-to-TS provider) gives end-to-end type safety: one schema is *both* the runtime validator and the TS type — request/response types are inferred from the schema, no duplicate interface.
- Async custom validators in a `preValidation` hook must **return `{ error }`** objects, not `throw` — a thrown error becomes an unhandled promise rejection that crashes the process.

### 1.4 Decorators
`decorate` / `decorateRequest` / `decorateReply` extend the instance/request/reply. Gotcha: `decorateRequest('user', {})` shares the **same object reference** across every request. For per-request state, decorate with `null` and assign in an `onRequest` hook.

### 1.5 Fastify anti-patterns / gotchas
- **Returning `undefined`** from an async handler → Fastify thinks "no response yet". Return the payload or call `reply.send()`.
- Using both `return value` **and** `reply.send(value)` — first wins, second is discarded with a warn log. Pick one.
- After `reply.send()` inside an **async** hook/handler, `return reply` (or `await reply`) to avoid "Reply already sent" / race conditions.
- **Arrow-function handlers** don't bind `this` to the Fastify instance — use `function` declarations when you need `this.<decorator>`.

## 2. NestJS

### 2.1 Modules + hierarchical DI
A Nest app is a tree of **modules**, each a `@Module({ imports, controllers, providers, exports })`. Modules **encapsulate their providers by default** — the `exports` array is the module's *public API*; another module sees a provider only by `imports`-ing the exporting module. `@Global()` makes a module's exports app-wide (register once, from the root/core module) — but **avoid making everything global**; explicit `imports` is the maintainable default.

### 2.2 Providers + scopes
Providers are `@Injectable()` classes (services, repositories, factories) resolved by the DI container. **Default scope is singleton** (one instance for the whole app). Other scopes:
- **`Scope.REQUEST`** — a fresh instance per incoming request, discarded after (carries a performance cost; it bubbles up to anything that injects it).
- **`Scope.TRANSIENT`** — a fresh instance per injection site.
Custom providers (`useClass` / `useValue` / `useFactory` / `useExisting`) encapsulate non-trivial initialization (config-driven instances, async factories).

### 2.3 The request pipeline (memorize the order)
```
Middleware → Guards → Interceptors(pre) → Pipes → [Route Handler] → Interceptors(post) → Exception Filters
```
- **Middleware** — earliest; raw req/res (Express/Fastify-style), before DI-aware components.
- **Guards** — authZ decision (roles/permissions/ACL): return `true`/`false` to admit the request. Run after middleware, before interceptors/pipes.
- **Interceptors** — wrap the whole handler (RxJS): transform request *and* response, add timing/caching/logging; code before `handle()` runs pre-handler, after runs post-handler.
- **Pipes** — validate/transform the inbound payload (e.g. `ValidationPipe` + class-validator/Zod) just before the handler.
- **Exception Filters** — catch thrown exceptions and shape the HTTP error response (built-in layer catches anything unhandled).

### 2.4 Dynamic modules + circular deps
- **Dynamic modules** (`forRoot()`/`forFeature()` returning `DynamicModule`) configure a module at import time — the pattern behind feature loaders and configurable infra modules.
- **Circular dependency** → `forwardRef()` on **both** sides (module-level in `imports`, service-level via `@Inject(forwardRef(() => X))`). Treat it as a smell: prefer refactoring the dependency direction (extract a shared module / mediator) over leaning on `forwardRef`.

### 2.5 HTTP adapter (Fastify under Nest)
Nest abstracts the HTTP layer behind an adapter. Swap `@nestjs/platform-express` for `@nestjs/platform-fastify` to get Fastify's throughput while keeping Nest's architecture (`NestFactory.create<NestFastifyApplication>(AppModule, new FastifyAdapter())`). Caveat: some Express-specific middleware/recipes need Fastify equivalents.

## 3. Hono

### 3.1 Web-Standards core + the Context object
Hono is built on the **WHATWG `Request`/`Response`** web standard, which is why it runs unchanged across Cloudflare Workers, Deno, Bun, Lambda, Fastly, and Node. The single thing to learn is the **Context `c`**: `c.req`, `c.json()/c.text()/c.html()`, `c.get()/c.set()`, `c.env`. Middleware is just `async (c, next) => { /* before */ await next(); /* after */ }` — the Koa-style onion model, with no separate req/res objects to juggle.

### 3.2 Routing + typed Variables/Bindings
The default **`RegExpRouter`** is the fastest router on Workers. Type the context generics so values are inferred end-to-end:
```ts
type Env = { Variables: { user: User }, Bindings: { DB: D1Database } }
const app = new Hono<Env>()
// c.get('user') is typed User; c.env.DB is typed D1Database in every middleware/handler
```

### 3.3 Validation + RPC (type-safe client)
First-party validators (`@hono/zod-validator`, plus Valibot/Typia/ArkType) run schema validation with full TS inference (`c.req.valid('json')` is typed). **RPC mode**: write a validator + chain routes, export the app *type*, then `hc<typeof app>(url)` on the client infers every path, argument, and return type — a type-safe client with no codegen and no shared OpenAPI step.

### 3.4 Batteries + edge deployment
Core package ships JWT, basic/bearer auth, CORS, CSRF, secure-headers, ETag, cache, compression, body-limit, IP restriction, timing, timeout, SSE, WebSockets, and JSX SSR. Real-world edge win comes mostly from **geography** (Workers run at the nearest PoP) plus Hono's low overhead — e.g. an Express-on-VM → Hono-on-Workers move reported p50 12ms→4ms, p99 45ms→9ms. Also the recommended path for a Node team moving to **Bun** that wants forward-compatibility.

## 4. Choosing among them (and vs Express)

| Choose… | When |
| --- | --- |
| **NestJS** | Enterprise architecture, team of 3+, long-lived backend, conventions matter; want first-party modules (auth, GraphQL, microservices, queues); coming from Angular/.NET; edge deployment **not** required. |
| **Fastify** | Standalone Node API; want Express-shaped ergonomics with better throughput **and** built-in schema validation/serialization; happy to choose your own ORM/auth/structure. |
| **Hono** | Deploying to Cloudflare Workers / Vercel Edge / Deno / Bun / Lambda; want the smallest, fastest, most type-inferred option; serverless/edge functions; Node→Bun migration. |
| **Express** (see `express-patterns.md`) | Maximal middleware ecosystem, maximum familiarity, no strong perf/validation/edge requirement. |

Combination worth knowing: **NestJS + Fastify adapter** = Nest's architecture with Fastify's speed. **NestJS does not run on edge runtimes** — if edge is a hard requirement, it's Hono (or Express/Fastify on a Node-compatible edge only where supported).

## 5. Troubleshooting quick table

| Symptom | Framework | Likely cause / fix |
| --- | --- | --- |
| Decorator/hook "not defined" on parent instance | Fastify | Plugin not wrapped in `fastify-plugin` — encapsulation kept it child-scoped. |
| Same object/state leaking across requests | Fastify | `decorateRequest('x', {})` shares one reference — init `null`, assign in `onRequest`. |
| Handler "hangs" / no response | Fastify | Async handler returned `undefined`, or `reply.send()` without `return reply`. |
| "Reply already sent" | Fastify | Mixed `return value` + `reply.send()`, or didn't `return reply` after async `send`. |
| `Nest can't resolve dependencies` | NestJS | Provider not in module `providers`/`exports`, or a circular dep needing `forwardRef()` on both sides. |
| Guard/validation runs in wrong order | NestJS | Pipeline order is fixed: Middleware → Guards → Interceptors → Pipes → Handler → Filters. |
| Request-scoped provider slow / unexpected new instances | NestJS | `Scope.REQUEST` bubbles up the injection chain — keep it narrow; default to singleton. |
| `c.env.X` / `c.get('y')` typed as `unknown` | Hono | Didn't parameterize `new Hono<{ Bindings, Variables }>()`. |
| RPC client types not inferred | Hono | Didn't `export type AppType = typeof app` or route lacks a validator; chain routes so types flow into `hc<AppType>()`. |

## 6. Sources
- Fastify docs — Encapsulation, Plugins, Hooks, Validation-and-Serialization, Decorators, Errors, Routes, TypeScript: https://fastify.dev/docs/latest/Reference/ ; Plugins Guide: https://fastify.dev/docs/latest/Guides/Plugins-Guide/
- Nearform — The complete guide to the Fastify plugin system: https://nearform.com/digital-community/the-complete-guide-to-fastify-plugin-system/
- Strapi — Build Production-Ready APIs with Fastify: https://strapi.io/blog/build-production-ready-apis-with-fastify
- NestJS docs — Modules, Circular dependency, Common errors, Performance (Fastify): https://docs.nestjs.com/modules , https://docs.nestjs.com/fundamentals/circular-dependency , https://docs.nestjs.com/techniques/performance
- DeepWiki — NestJS Request Processing Pipeline: https://deepwiki.com/nestjs/docs.nestjs.com/2.3-middleware-guards-pipes-and-interceptors
- Medium (Juel / Patel) — Middleware vs Guards vs Interceptors vs Pipes vs Filters in NestJS
- LogRocket — How to avoid circular dependencies in NestJS: https://blog.logrocket.com/avoid-circular-dependencies-nestjs/
- Hono docs — Concepts, RPC, Validation, Benchmarks, Stacks, Middleware: https://hono.dev/docs
- Cloudflare blog — The story of web framework Hono: https://blog.cloudflare.com/the-story-of-web-framework-hono-from-the-creator-of-hono/
- freeCodeCamp — Build Production-Ready Web Apps with Hono: https://www.freecodecamp.org/news/build-production-ready-web-apps-with-hono/
- Encore — NestJS vs Fastify vs Hono (2026): https://encore.dev/articles/nestjs-vs-fastify-vs-hono
- Better Stack — Hono vs Fastify / Fastify vs Express vs Hono: https://betterstack.com/community/guides/scaling-nodejs/hono-vs-fastify/
- HireNodeJS — Best Node.js Frameworks 2026: https://www.hirenodejs.com/blog/nodejs-frameworks-compared-2026
