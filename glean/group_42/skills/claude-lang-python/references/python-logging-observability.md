# Python Logging & Application Observability

`lang-python` hub reference. Emitting good logs from a Python app: the stdlib
`logging` model, structured logging (structlog / loguru), and wiring logs into
OpenTelemetry.

**Scope:** the Python-side logging libraries and patterns. The *platform* side
(running an OTel Collector, choosing a backend (Grafana/Datadog/etc.), traces and
metrics pillars) is `devops-observability`. Profiling a slow app is
`cpython-performance-profiling`. "What library matters less than the practices
around it" — get structure, context, and levels right and any library serves.

---

## 1. The stdlib `logging` model

Four object types, arranged as a pipeline:

- **Logger** — what you call (`logging.getLogger(__name__)`). Named, hierarchical
  by dotted name (`a.b` is a child of `a`).
- **Handler** — where records go (`StreamHandler`, `FileHandler`,
  `RotatingFileHandler`, `QueueHandler`).
- **Formatter** — how a record is rendered (text template or JSON).
- **Filter** — fine-grained include/exclude and record enrichment.

```python
import logging
log = logging.getLogger(__name__)        # module-level, named per module
log.info("user logged in", extra={"user_id": 42})
```

Key behaviors:
- **Levels:** DEBUG < INFO < WARNING < ERROR < CRITICAL; set per-logger and
  per-handler. WARNING is the default root level.
- **Propagation:** records bubble up the logger hierarchy to ancestors' handlers.
  Configure handlers at the root (or a top package) and let children propagate;
  don't add handlers everywhere.
- **`logging.getLogger(__name__)`** in every module, giving you the dotted
  hierarchy for free and per-subsystem verbosity tuning.

---

## 2. Configure once, and not in libraries

- **Applications** configure logging **once at the entry point**, ideally via
  `logging.config.dictConfig({...})` (declarative: formatters, handlers, levels,
  loggers). Avoid scattered `basicConfig` calls.
- **Libraries must not configure logging.** Add a `logging.NullHandler()` to your
  library's top logger and emit records only; let the *application* decide
  handlers/levels. Configuring in a library hijacks the host app's logging.

```python
# library top-level __init__.py
logging.getLogger("mylib").addHandler(logging.NullHandler())
```

---

## 3. Structured logging — why

Plain-text logs are hard for machines. **Structured logging** emits events as
key-value pairs (usually JSON) so they're queryable, aggregatable, and
correlatable in an observability backend. Prefer it for any service: consistent
fields + trace IDs are what make logs useful in production.

```
# instead of: "user 42 purchased item 99 for $12.50"
{"event": "purchase", "user_id": 42, "item_id": 99, "amount": 12.50, "level":"info"}
```

---

## 4. structlog — composable structured logging

Turns logging into a stream of structured events via a **processor pipeline**.
Preferred in large codebases for its composability.

```python
import structlog
structlog.configure(
    processors=[
        structlog.contextvars.merge_contextvars,     # pull in bound context
        structlog.processors.add_log_level,
        structlog.processors.TimeStamper(fmt="iso"),
        structlog.processors.JSONRenderer(),          # machine-readable out
    ],
)
log = structlog.get_logger()
log = log.bind(request_id="abc")                       # bound logger carries context
log.info("handled", path="/x", ms=12)
```

- **Bound loggers** carry context (`.bind()`); **`contextvars`** integration
  propagates request/trace context across `await` boundaries without threading it
  manually.
- Can sit **on top of stdlib** (`structlog.stdlib.*`) so third-party libs logging
  via stdlib still flow through your pipeline.
- **Keep processors short and pure**; defer expensive work (stack info, exception
  formatting) so the hot path stays fast.

---

## 5. loguru — zero-config structured logging

A single pre-configured `logger`; "it just works" with one import. Best when you
want structured logs quickly without wiring.

```python
from loguru import logger
logger.add("app.json", serialize=True, rotation="100 MB")  # JSON sink + rotation
logger.bind(user_id=42).info("logged in")
try: ...
except Exception:
    logger.exception("failed")                             # rich traceback w/ diagnose
```

- One global `logger`; `add()` registers **sinks** (files, stderr, callables) with
  per-sink level/format/rotation/retention/serialization.
- To reach stdlib-based libraries or export to OTel, route stdlib through an
  **`InterceptHandler`** into loguru, an extra indirection that stdlib+structlog
  don't need. That coupling is the main trade-off vs structlog.

---

## 6. OpenTelemetry logs — correlating logs with traces

The point of OTel logging is **log↔trace correlation**: stamp each log with the
active `trace_id`/`span_id` so a log line links to its distributed trace.

- `opentelemetry-sdk` provides a `LoggingHandler` that bridges stdlib records to
  the OTel log pipeline and exports via **OTLP**.
- Common production pattern: keep structlog/loguru writing **JSON to stdout** and
  let the **OpenTelemetry Collector** convert it to the OTLP log schema and
  forward it, decoupling the app from backend specifics.
- Correlation uses `contextvars` + the active span; structlog's
  `merge_contextvars` (or an OTel processor) injects the IDs.
- Collector setup, exporters, sampling, and the traces/metrics pillars →
  `devops-observability`.

---

## 7. Best practices & pitfalls

- **Log structured, log context, log levels**: these matter more than the
  library. Use consistent field names (`user_id`, not `uid`/`user`).
- **Propagate context with `contextvars`**, not function args (survives async).
- **Never log secrets/PII**: add a redaction processor/filter; this overlaps
  `python-supply-chain-security` hygiene.
- **Exceptions:** `log.exception(...)` / `exc_info=True` to capture tracebacks;
  don't `str(e)` and lose the stack.
- **Performance:** guard expensive debug payloads (`if log.isEnabledFor(logging.DEBUG)`),
  push heavy formatting off the hot path, consider `QueueHandler`/`QueueListener`
  to move I/O off request threads.
- **Don't use `print` for app logging**, and don't reconfigure the root logger
  from inside a library.
- **Pitfall — double logging:** adding handlers on a child logger *and* leaving
  propagation on duplicate records; configure at the root or set
  `propagate=False` deliberately.

### Choosing

| Want | Pick |
| --- | --- |
| Composable pipeline, big codebase, OTel-native context | **stdlib + structlog** |
| Fastest setup, small/standalone app, batteries included | **loguru** |
| Export logs as OTel signals correlated with traces | structlog/loguru → **OTel Collector** (OTLP) |

---

## Sources

- [Logging — Python stdlib docs](https://docs.python.org/3/library/logging.html) and [Logging HOWTO](https://docs.python.org/3/howto/logging.html) (hierarchy, dictConfig, NullHandler)
- [structlog documentation](https://www.structlog.org/) ; [Leveling up Python logs with structlog — Dash0](https://www.dash0.com/guides/python-logging-with-structlog)
- [Choosing a Python Logging Library in 2026 — Dash0](https://www.dash0.com/guides/python-logging-libraries)
- [Loguru documentation](https://loguru.readthedocs.io/) ; [Structured logging with Loguru — soumendrak](https://www.soumendrak.com/series/practical-observability-with-python/structured-logging/)
- [How to Structure Logs with OpenTelemetry (Python) — OneUptime](https://oneuptime.com/blog/post/2025-01-06-python-structured-logging-opentelemetry/view)
- [Python Logging Best Practices — SigNoz](https://signoz.io/guides/python-logging-best-practices/)
