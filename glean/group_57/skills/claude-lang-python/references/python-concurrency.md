# Python Concurrency — asyncio, Structured Concurrency, anyio/trio

`lang-python` hub reference. The **application-level** concurrency toolkit:
writing correct async code, scoping task lifetimes, and handling cancellation.

**Scope boundary:** the *runtime mechanism* of concurrency — the GIL, the
free-threaded (no-GIL) build, subinterpreters, the JIT, and the
free-threading-vs-subinterpreters-vs-multiprocessing-vs-asyncio **decision
matrix** — lives in `cpython-runtime-internals`. This spoke is the layer above:
how you actually *write* concurrent Python once you've chosen async. For raw
profiling of an async app, see `cpython-performance-profiling`.

---

## 1. The async execution model

- **Coroutine** — an `async def` function. Calling it returns a coroutine
  object; it does nothing until awaited or scheduled on a loop.
- **Awaitable** — anything usable with `await`: coroutines, `Task`s, `Future`s.
- **Task** — a coroutine wrapped and scheduled on the event loop
  (`asyncio.create_task`). Runs concurrently with other tasks; cooperative —
  control only yields at `await` points. CPU-bound code between awaits blocks the
  whole loop.
- **Event loop** — single-threaded scheduler. One loop per thread. `asyncio.run()`
  creates a loop, runs the coroutine, and closes the loop (don't nest it).

```python
import asyncio

async def fetch(n: int) -> int:
    await asyncio.sleep(n)        # yields control to the loop
    return n * 2

async def main() -> None:
    result = await fetch(1)      # sequential await
asyncio.run(main())              # the one entry point
```

Mental model: `await` = "I might suspend here; run something else." If a path has
no `await`, nothing else runs.

---

## 2. Structured concurrency — `asyncio.TaskGroup` (3.11+)

The modern default for running tasks concurrently. A `TaskGroup` is an async
context manager: every task created in it is awaited at block exit, and **if any
task raises, the rest are cancelled** and the errors surface together as an
`ExceptionGroup` (PEP 654). This is the structured-concurrency guarantee — child
task lifetimes are bound to the parent block; no orphans.

```python
async def main() -> None:
    async with asyncio.TaskGroup() as tg:
        t1 = tg.create_task(fetch(1))
        t2 = tg.create_task(fetch(2))
    # both guaranteed done here; results via t1.result(), t2.result()
    print(t1.result(), t2.result())
```

### TaskGroup vs `gather`

| | `TaskGroup` (3.11+) | `asyncio.gather` |
| --- | --- | --- |
| One task fails | cancels siblings, raises `ExceptionGroup` | by default others keep running; first exc propagates (or collected with `return_exceptions=True`) |
| Lifetime | bound to `async with` block | unbound — you hold the futures |
| Errors | aggregated `ExceptionGroup` | single exception or list |
| Use when | you want all-or-nothing, leak-free fan-out | you need per-result error handling / partial success |

Prefer `TaskGroup` for new code. Keep `gather(..., return_exceptions=True)` when
you genuinely want partial success and to inspect each outcome.

Handle aggregated errors with `except*`:

```python
try:
    async with asyncio.TaskGroup() as tg:
        tg.create_task(may_fail())
        tg.create_task(may_also_fail())
except* ValueError as eg:
    for exc in eg.exceptions: ...
except* ConnectionError as eg:
    ...
```

3.13 hardened TaskGroup: it now closes a coroutine if the group is no longer
active, and fixed simultaneous internal/external cancellation and cancel-count
preservation. 3.14: `create_task()` forwards all kwargs to `loop.create_task()`.

---

## 3. Cancellation and timeouts

Cancellation is **cooperative**: `task.cancel()` schedules a `CancelledError` to
be raised at the task's next `await`. Never swallow it blindly; catch only to
clean up, then **re-raise**:

```python
try:
    await long_op()
except asyncio.CancelledError:
    await cleanup()
    raise                         # propagate — do not absorb cancellation
```

### `asyncio.timeout()` (3.11+) — the modern timeout

A context manager that cancels the block on expiry; far cleaner than
`wait_for`. The cancellation is **isolated to the block** via the uncancel
mechanism, so outer code keeps running:

```python
async with asyncio.timeout(5.0):        # or timeout_at(asyncio.get_running_loop().time() + 5)
    await slow_operation()
# raises TimeoutError if it didn't finish in 5s
```

- `Task.uncancel()` (3.11) decrements the cancel count; `TaskGroup` and
  `timeout()` use it to keep a single timeout from leaking cancellation into the
  surrounding scope.
- `asyncio.shield()` protects an awaitable from cancellation — use sparingly;
  it breaks the structured model and is a common source of leaks.
- **Pitfall:** `asyncio.timeout(0)` can swallow a prior pending cancellation
  (CPython #134471). On 3.12, cancellation could leak out of a `TaskGroup` when
  using eager tasks (#128588). Stay current on patch releases.

---

## 4. The fire-and-forget garbage-collection trap

`asyncio.create_task()` returns a task the **loop holds only weakly**. If you
don't keep a reference, the task can be garbage-collected mid-flight and silently
vanish, a notorious footgun the docs now warn about explicitly.

```python
# WRONG — task may be GC'd before it finishes
asyncio.create_task(background_work())

# RIGHT — keep a strong reference; self-clean on completion
_background: set[asyncio.Task] = set()

def spawn(coro) -> None:
    t = asyncio.create_task(coro)
    _background.add(t)
    t.add_done_callback(_background.discard)
```

Better still: don't fire-and-forget at all. Put the task in a `TaskGroup` (or an
anyio task group) whose block outlives the work, so the lifetime is explicit.

**Eager task factory (3.12+):** `loop.set_task_factory(asyncio.eager_task_factory)`
runs the coroutine immediately up to its first real suspension; if it completes
without blocking, it finishes eagerly (skipping a loop iteration). A throughput win
for tasks that often return without awaiting, but it changes scheduling order, so
test cancellation paths.

---

## 5. `anyio` and `trio` — structured concurrency, done thoroughly

`asyncio.TaskGroup` is a narrow API: it can't list or cancel individual contained
tasks. **trio** pioneered the richer model; **anyio** brings it to asyncio.

- **trio** — a batteries-included async runtime that *enforces* structured
  concurrency. Tasks must live in a **nursery** (`async with trio.open_nursery()`);
  there is no fire-and-forget. Cancellation flows through a task **tree** via
  **cancel scopes** (`trio.move_on_after`, `trio.fail_after`), which are
  nestable, inspectable, and first-class.
- **anyio** — a portability layer exposing trio's model (`anyio.create_task_group`,
  `anyio.move_on_after`, `anyio.fail_after`, `CancelScope`) that runs **on either
  asyncio or trio**. Libraries should target anyio APIs to stay backend-agnostic
  (this is how FastAPI/Starlette, HTTPX, and others do it).

```python
import anyio

async def main() -> None:
    async with anyio.create_task_group() as tg:
        tg.start_soon(worker, 1)
        tg.start_soon(worker, 2)
        with anyio.move_on_after(5):      # cancel scope: soft timeout
            await slow()
```

Cancel scopes (anyio/trio) are strictly **level-scoped**: cancelling a scope
cancels everything inside it, and the cancellation stops at the scope boundary
(the same property `asyncio.timeout` reimplements via uncancel). Use anyio when
you want portable, composable cancellation and richer task-group control than
stdlib asyncio gives.

---

## 6. Bridging to threads and processes

Async is for **I/O-bound** concurrency. For blocking calls and CPU-bound work,
hand off so you don't stall the loop:

- `await asyncio.to_thread(func, *args)` (3.9+) — run a blocking function in the
  default thread pool. The simplest correct way to call sync I/O from async.
- `loop.run_in_executor(pool, func, ...)` — explicit `ThreadPoolExecutor` (I/O)
  or `ProcessPoolExecutor` (CPU-bound, sidesteps the GIL via processes).
- anyio: `anyio.to_thread.run_sync` / `anyio.to_process.run_sync`.

Choosing async vs threads vs processes vs free-threading is the **decision
matrix in `cpython-runtime-internals`** — this spoke assumes you've chosen async.

---

## 7. The async toolbox

- **Sync primitives** (`asyncio.Lock`, `Semaphore`, `Event`, `Condition`,
  `Queue`) — async-aware; never use `threading.Lock` in a coroutine. A
  `Semaphore` is the standard concurrency limiter for fan-out:
  ```python
  sem = asyncio.Semaphore(10)
  async def bounded(x):
      async with sem:
          return await fetch(x)
  ```
- **Async iterators / `async for`** — `__aiter__` / `__anext__`; stream results.
- **Async generators** (`async def` + `yield`) — but see PEP 789 below.
- **Async context managers** (`__aenter__` / `__aexit__`,
  `contextlib.asynccontextmanager`) — for async setup/teardown.
- **`asyncio.as_completed` / `asyncio.wait`** — lower-level fan-out when you need
  results as they land rather than the all-or-nothing TaskGroup contract.

---

## 8. PEP 789 and common pitfalls

- **PEP 789 — don't `yield` across a `TaskGroup`/cancel scope.** Suspending a
  frame (via `yield` in an async generator, or `contextlib.asynccontextmanager`)
  while a TaskGroup or `timeout()` is open lets the wrong task be cancelled,
  timeouts be ignored, and exceptions be mishandled. Keep task groups and cancel
  scopes inside a single coroutine frame; don't wrap them in an async generator.
- **Blocking the loop** — any synchronous CPU work or sync I/O (`requests`,
  `time.sleep`, blocking DB driver) between awaits freezes *all* tasks. Use async
  libraries or `to_thread`. `time.sleep` → `asyncio.sleep`.
- **Swallowing `CancelledError`** — catching it without re-raising breaks
  cancellation and structured concurrency. Always re-raise after cleanup.
- **Forgetting `asyncio.run` closes the loop** — don't call it twice or nest it;
  one entry point per program. Use `asyncio.Runner` (3.11+) for multiple
  top-level runs sharing setup.
- **Mixing event loops across threads** — a loop belongs to its thread; use
  `loop.call_soon_threadsafe` / `asyncio.run_coroutine_threadsafe` to cross.
- **Unbounded fan-out** — `TaskGroup` with thousands of tasks can exhaust
  sockets/memory; gate with a `Semaphore`.

---

## Sources

- [Coroutines and Tasks — Python 3.14 docs](https://docs.python.org/3/library/asyncio-task.html) — TaskGroup, create_task GC warning, eager_start, timeout, uncancel
- [PEP 789 – Preventing task-cancellation bugs by limiting yield in async generators](https://peps.python.org/pep-0789/)
- [Cancellation and timeouts — AnyIO docs](https://anyio.readthedocs.io/en/stable/cancellation.html) and [Why AnyIO over asyncio](https://anyio.readthedocs.io/en/stable/why.html)
- [AnyIO Task Groups & Structured Concurrency — DeepWiki](https://deepwiki.com/agronholm/anyio/2.2-cancellation-and-timeouts)
- [Fire and forget with Python's asyncio — Michael Kennedy](https://mkennedy.codes/posts/fire-and-forget-or-never-with-python-s-asyncio/) and [Garbage-collected asyncio.Task objects — discuss.python.org](https://discuss.python.org/t/whats-up-with-garbage-collected-asyncio-task-objects/29686)
- CPython issues [#128588 (eager-task cancellation leak)](https://github.com/python/cpython/issues/128588), [#134471 (timeout(0) swallows cancellation)](https://github.com/python/cpython/issues/134471)
