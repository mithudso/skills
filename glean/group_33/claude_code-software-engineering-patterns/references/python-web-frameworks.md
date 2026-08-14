<!-- hub-reference-banner -->
> **Reference file — part of the `software-engineering-patterns` hub.** Authored 2026-05-31 from web research (`/dr`).
> Sibling topics in this family are reference files under the hubs (`programming-languages`, `software-engineering-patterns`) — **not** standalone
> skills. Ignore any "use the X skill" pointers that name a bare sibling; load that topic's `references/<name>.md` from the owning hub.
> For Python language idioms (typing, async syntax, packaging) see `programming-languages/references/python-patterns.md`; for Pydantic v2 model details see `programming-languages/references/pydantic-v2.md`; for the Node analog see `references/express-patterns.md`.

---
name: python-web-frameworks
version: "1.0.0"
updated: "2026-05-31"
category: developer
tags: [python, fastapi, django, flask, asgi, wsgi, uvicorn, gunicorn, async, web-framework, pydantic, drf, rest-api]
description: "Python web frameworks expert — FastAPI, Django, Flask, and the ASGI/WSGI server foundation they share. TRIGGER: choosing between FastAPI, Django, and Flask; building a FastAPI service (path operations, Depends DI, Pydantic v2 models, async vs sync routes, BackgroundTasks, lifespan, APIRouter); building a Django app (MVT, ORM, async views, DRF, ASGI deployment, Django 6.0 features); building a Flask app (application factory, blueprints, app/request context, extensions); deploying any of them (uvicorn, gunicorn, hypercorn, worker classes, lifespan protocol); diagnosing event-loop blocking, sync-in-async, or N+1 ORM problems. SKIP: Python language syntax/idioms → programming-languages/python-patterns; Pydantic v2 model deep-dive → programming-languages/pydantic-v2; Node.js Express → references/express-patterns; generic REST/GraphQL API surface design → references/api-design-patterns."
related_skills: [api-design-patterns, backend-patterns, web-auth-patterns, performance-profiling-expert]
---

# Python Web Frameworks: FastAPI, Django, Flask, and ASGI/WSGI

The three dominant Python web frameworks sit on two server-interface standards. **Pick the framework from the workload, not the hype.** This reference covers the shared WSGI/ASGI foundation first (because all three frameworks inherit from it), then each framework, then selection and production patterns.

> Cross-references: API surface design (REST/GraphQL/gRPC, versioning, resource modeling) → `references/api-design-patterns.md`. Auth flows (OAuth 2.1/PKCE, JWT, sessions, CSRF) → `references/web-auth-patterns.md`. Backend service architecture → `references/backend-patterns.md`. Pydantic v2 model/validator depth → `programming-languages/references/pydantic-v2.md`.

## 1. The shared foundation — WSGI vs ASGI (read this first)

Both standards are the contract between a **Python web application** and a **web/application server**. The framework you choose is mostly a choice of which standard it speaks.

| | WSGI (PEP 3333) | ASGI |
|---|---|---|
| Model | Synchronous, one request per worker thread/process | Asynchronous, `async`/`await`, event loop |
| Concurrency | Blocking; scale by adding workers/threads | Non-blocking I/O; one worker handles many in-flight requests |
| Protocols | HTTP/1.x request-response only | HTTP/1.x, HTTP/2, WebSockets, long-lived connections, server push |
| Lifespan | No startup/shutdown hooks in the protocol | **Lifespan protocol** — startup/shutdown events (init DB pools, caches) |
| Frameworks | Flask, Django (default), classic stacks | FastAPI/Starlette, Django (under ASGI), Quart |
| Servers | Gunicorn, uWSGI, mod_wsgi | Uvicorn, Hypercorn, Daphne, Granian, Gunicorn's ASGI worker |

Key facts:
- ASGI servers are roughly **2–4× faster than WSGI for async, high-concurrency workloads** (many simultaneous connections, WebSockets, long-polling). For purely CPU-bound or simple synchronous request/response work the gap is small or nonexistent — async wins on **concurrency**, not raw compute.
- An ASGI app can serve WSGI apps via an adapter (and Hypercorn can serve both), but you cannot run an async ASGI app on a pure WSGI server.

### Servers (the runtime you actually deploy)
- **Uvicorn** — the standard ASGI server (FastAPI/Starlette). In production run it **under Gunicorn** using the Uvicorn worker class for process management + graceful restarts: `gunicorn app:app -k uvicorn_worker.UvicornWorker -w 4`. (The worker class moved to the `uvicorn-worker` package; `uvicorn.workers.UvicornWorker` is deprecated.)
- **Gunicorn** — the reliable WSGI default for Django/Flask. Recent Gunicorn ships a **native ASGI worker**, so you can serve FastAPI/Starlette/Quart without a separate Uvicorn dependency. Worker rule of thumb: `(2 × CPU cores) + 1` for sync WSGI workers.
- **Hypercorn** — ASGI server with **HTTP/2 + WebSockets** and the ability to serve **both ASGI and WSGI** apps. The versatile choice when you need HTTP/2.
- **Daphne** — the original Django Channels ASGI server; still used for Channels/WebSocket-heavy Django.
- **Granian** — a newer Rust-based server (RSGI/ASGI/WSGI), gaining traction for performance.
- **2026 default recommendation:** unless you must squeeze raw req/sec, **Uvicorn under Gunicorn** is the lowest-operational-cost choice for async; **Gunicorn** for sync Django/Flask.

### The lifespan protocol
ASGI's lifespan lets the app run code at startup and shutdown — the right place to open/close DB connection pools, warm caches, and start background workers. FastAPI exposes it through the `lifespan` context manager (preferred over the deprecated `@app.on_event("startup"/"shutdown")`).

## 2. FastAPI

Modern, async-first, type-driven framework built on **Starlette** (ASGI toolkit) + **Pydantic v2** (validation). Auto-generates OpenAPI docs (`/docs`, `/redoc`).

**Core concepts:**
- **Path operations** — decorated functions (`@app.get`, `@app.post`, …). Use `async def` only when you `await` non-blocking I/O; use plain `def` for blocking work (FastAPI runs `def` routes in a threadpool automatically).
- **Pydantic models** for request/response. `response_model=` shapes and validates output and drives the OpenAPI schema. Input and output schemas can differ (e.g. a field with a default is optional on input but always present on output).
- **Dependency injection (`Depends`)** — the signature feature. Declare reusable dependencies (auth, DB sessions, pagination params, settings) and FastAPI resolves/injects them. Use `Annotated[T, Depends(fn)]` (modern style). Dependencies can be sync or async and can be stacked.
- **`BackgroundTasks`** for fire-and-forget work after the response is sent.
- **`lifespan`** context manager for startup/shutdown.
- **`APIRouter`** to split routes across modules (the basis of the standard project layout).

```python
from typing import Annotated
from fastapi import FastAPI, Depends

async def common_parameters(q: str | None = None, skip: int = 0, limit: int = 100):
    return {"q": q, "skip": skip, "limit": limit}

app = FastAPI()

@app.get("/items/")
async def read_items(commons: Annotated[dict, Depends(common_parameters)]):
    return commons
```

**Settings/config:** subclass `pydantic_settings.BaseSettings`, load env vars with type-checked validation, and inject via a `@lru_cache`-wrapped dependency so config is parsed once and overridable in tests.

**When FastAPI:** async-heavy APIs, microservices, ML model serving, high-concurrency I/O-bound endpoints, anywhere you want first-class typed contracts + auto OpenAPI.

## 3. Django

The batteries-included framework: ORM, admin, auth, migrations, forms, templating — **MVT** (Model-View-Template) architecture. Encourages rapid, secure development of large, feature-complete apps. **Django 6.0** is current (with 5.2 LTS as the conservative choice).

**Core concepts:**
- **ORM** — `models.Model` classes map to tables; rich query API. Watch for **N+1 queries** — use `select_related()` (joins) / `prefetch_related()` (separate queries) and `.only()`/`.defer()`.
- **MVT** — URLs → views → templates; the "controller" is the framework itself.
- **Migrations** — `makemigrations` / `migrate` for schema evolution.
- **Middleware** — request/response processing chain.
- **DRF (Django REST Framework)** — the standard add-on for building REST APIs (serializers, viewsets, routers, auth). Use this, not raw Django views, for JSON APIs.

**Async support:** Django supports **async views** (`async def`) and a fully async request stack **under ASGI**. Async views run under WSGI too but with a performance penalty and no efficient long-running requests. The ORM has an **async interface** (`a`-prefixed methods: `aget()`, `acreate()`, `afirst()`, `async for`), but the underlying DB operations are still synchronous — full native-async ORM is in progress. Use `sync_to_async()` to bridge to synchronous Django parts. Deploy async Django with `gunicorn myproject.asgi:application -k uvicorn_worker.UvicornWorker`.

```python
import asyncio
from django.http import HttpResponse

async def my_view(request):
    await asyncio.sleep(0.5)
    return HttpResponse("Hello, async world!")
```

**When Django:** content-heavy sites, admin-driven CRUD, apps needing auth/ORM/admin out of the box, large teams wanting convention over configuration, fast time-to-market with security defaults (CSRF, XSS, SQL-injection protections built in).

## 4. Flask

The micro-framework: a minimal **WSGI** core (Werkzeug + Jinja2) you extend yourself. Maximum flexibility, minimal opinions. Current line is **Flask 3.x**.

**Core concepts:**
- **Application factory** — define `create_app()` that builds the `Flask` app, loads config, initializes extensions, registers blueprints, and returns it. Enables multiple instances + clean testing.
- **Blueprints** — modular groups of routes/templates registered onto the app (`app.register_blueprint(...)`). The scaling unit for larger Flask apps.
- **App/request context** — at import time inside a blueprint the app object isn't bound, so use `current_app` and `g`/`request` proxies to reach the active instance/config.
- **Extensions** — design them unbound (`db = SQLAlchemy()`) and bind later with `db.init_app(app)` inside the factory so one extension object serves multiple app instances.

```python
def create_app(config_filename):
    app = Flask(__name__)
    app.config.from_pyfile(config_filename)
    from yourapp.model import db
    db.init_app(app)
    from yourapp.views.admin import admin
    app.register_blueprint(admin)
    return app
```

**Async:** Flask added `async def` view support (via an event loop per request on the WSGI worker), but it does **not** make Flask a true async framework — there's no shared event loop and no ASGI lifespan; for genuinely async workloads use **Quart** (the ASGI-native, Flask-API-compatible sibling) or FastAPI.

**When Flask:** small-to-mid apps, prototypes, services where you want to choose every component, learning/teaching, embedding a web layer in a larger app.

## 5. Choosing between them (selection rubric)

| Need | Pick |
|---|---|
| High-concurrency async I/O API, typed contracts, auto OpenAPI | **FastAPI** |
| ML/inference serving, microservices, WebSockets | FastAPI (or Quart) |
| Full-stack app with admin, ORM, auth, migrations out of the box | **Django** (+ DRF for APIs) |
| Content site, CRUD-heavy, large team, security defaults, fast TTM | **Django** |
| Minimal/flexible micro-service, prototype, pick-your-own-components | **Flask** |
| Flask ergonomics but truly async | **Quart** |

Rules of thumb: **async + API-first → FastAPI; batteries-included → Django; minimal + flexible → Flask.** All three are production-grade; the wrong choice is usually picking on raw benchmark numbers rather than ecosystem fit and team familiarity.

## 6. Practical production patterns
- **Project structure:** FastAPI — split by `APIRouter` modules + a `dependencies.py` + Pydantic `Settings`. Django — apps per bounded context, DRF serializers/viewsets. Flask — application factory + blueprints package.
- **Config:** typed settings (Pydantic `BaseSettings` for FastAPI; `settings.py`/env for Django; `app.config.from_*` for Flask). Catch misconfiguration at startup, override in tests.
- **DB sessions:** inject a per-request session via DI (FastAPI `Depends`) and reuse a **connection pool** — don't open a new connection per request.
- **Deployment:** Uvicorn-under-Gunicorn (FastAPI), Gunicorn (Django/Flask sync) or Gunicorn+Uvicorn-worker (async Django); set worker count from CPU; use the lifespan/startup hooks to open pools.
- **Testing:** FastAPI `TestClient`/`httpx.AsyncClient`; Django test runner + DRF `APIClient`; Flask `app.test_client()` with the factory.
- **Background/CPU work:** offload to **Celery** or a worker queue; never do CPU-bound work inline in an async route.

## 7. Anti-patterns and troubleshooting

**The #1 FastAPI pitfall — blocking the event loop:**
- `async def` does **not** parallelize blocking calls. Calling `requests.get()`, a sync DB driver, or heavy CPU work inside an `async def` route **freezes every other request** until it returns.
- Fixes: (a) put blocking work in a plain `def` route (FastAPI threadpools it); (b) use a real async client (`httpx.AsyncClient`, async DB driver); (c) offload with `asyncio.to_thread()` / `run_in_executor()`; (d) push CPU-bound work to Celery.
- Wrapping a sync function in `async def` and `await`-ing it does NOT make it async — it just blocks on the event loop.

**Other common problems:**
- **Django N+1 queries** — missing `select_related`/`prefetch_related`; profile with `django-debug-toolbar`.
- **Django async-under-WSGI** — async views work but pay a penalty and can't do efficient long-lived requests; deploy under ASGI for real async.
- **Flask "app object at import time"** — accessing `app`/config at module import inside a blueprint; use `current_app` instead.
- **Per-request connections** — exhausts the DB; reuse a pool.
- **Settings re-parsed every request** — wrap the settings provider in `@lru_cache`.
- **Deprecated FastAPI startup hooks** — replace `@app.on_event` with the `lifespan` context manager.
- **Deprecated Uvicorn worker import** — use `uvicorn_worker.UvicornWorker`, not `uvicorn.workers.UvicornWorker`.

## References
- FastAPI official docs (dependencies, async, Pydantic models, background tasks, lifespan) — fastapi.tiangolo.com / github.com/fastapi/fastapi
- Django 6.0 docs (async support, async ORM interface, ASGI deployment, ORM queries) — docs.djangoproject.com/en/6.0
- Flask 3.x docs (application factories, blueprints, app context, extensions) — flask.palletsprojects.com / github.com/pallets/flask
- Gunicorn ASGI worker + Uvicorn/Hypercorn/Daphne/Granian server comparisons — gunicorn.org/asgi, uvicorn.org, leapcell.io, deployhq.com (Python Application Servers 2026)
- FastAPI production best practices (async pitfalls, DI, Pydantic Settings) — github.com/zhanymkanov/fastapi-best-practices, orchestrator.dev FastAPI production patterns 2025, fastapi.tiangolo.com/async
- Framework comparison 2026 (architecture, performance, use cases) — zestminds.com, developersvoice.com, JetBrains PyCharm blog (Django/Flask/FastAPI)
