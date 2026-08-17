<!-- hub-reference-banner -->
> **Reference file — part of the `lang-python` hub.** Formerly the standalone `cpython-runtime-internals` skill.
> Sibling topics in this family are now reference files under the hubs (`lang-python`, `lang-js-ts`, `lang-go-and-mobile`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

name: cpython-runtime-internals
description: >
  CPython interpreter and runtime internals for the 3.12–3.14 concurrency/performance overhaul — the three PEPs that rewire how CPython executes: free-threaded no-GIL build (PEP 703), multiple interpreters in the stdlib / per-interpreter GIL (PEP 734 + PEP 684), and the experimental copy-and-patch JIT (PEP 744). Covers the shared runtime execution model (GIL, tier-1 specializing adaptive interpreter PEP 659, PyInterpreterState/thread state, immortal objects PEP 683), how to build/detect each mode, C-extension compatibility, and the isolation/parallelism model.
  TRIGGER: free-threaded Python, no-GIL, --disable-gil, python3.14t, PYTHON_GIL, sys._is_gil_enabled, Py_GIL_DISABLED; PEP 703 / PEP 779; subinterpreters, multiple interpreters, concurrent.interpreters, InterpreterPoolExecutor, per-interpreter GIL, PEP 734 / PEP 684 / PEP 554; CPython JIT, copy-and-patch, tier 2 micro-ops, --enable-experimental-jit, PYTHON_JIT, PEP 744; specializing adaptive interpreter (PEP 659); immortal objects (PEP 683); biased/deferred reference counting; making a C extension free-thread safe (Py_mod_gil, Py_MOD_GIL_NOT_USED, PyUnstable_Module_SetGIL); choosing free-threading vs subinterpreters vs multiprocessing for Python parallelism.
  SKIP: everyday modern-Python idioms, asyncio/TaskGroup, type hints, packaging, ruff/pytest → python-patterns; AI/ML Python frameworks → ai-agent-engineering; Node.js/libuv runtime concurrency → nodejs-concurrency-internals.
version: "1.0"
category: developer
updated: "2026-05-31"
tags:
  - python
  - cpython
  - free-threading
  - no-gil
  - subinterpreters
  - jit
  - concurrency
  - performance
  - pep-703
  - pep-734
  - pep-744
related_skills:
  - python-patterns
  - nodejs-concurrency-internals
  - performance-profiling-expert
---

# CPython Runtime Internals — Free-Threading, Subinterpreters, and the JIT

Deep reference for the 2023–2026 CPython runtime overhaul. Three Python Enhancement Proposals reshape how the interpreter executes, all touching the same machinery (the eval loop, reference counting, and per-interpreter state):

| PEP | Feature | First shipped | Status @ 3.14 |
|-----|---------|---------------|----------------|
| **703** | Free-threaded (no-GIL) build | 3.13 (experimental) | **Supported, opt-in** (per PEP 779) |
| **734** | `concurrent.interpreters` stdlib module (built on PEP 684 per-interpreter GIL) | 3.14 | Final |
| **744** | Copy-and-patch JIT (tier 2) | 3.13 (experimental) | Still experimental, off by default |

> These are three different answers to "how do I use more than one core / go faster in Python." Free-threading removes the lock; subinterpreters give each thread its own lock + isolated heap; the JIT speeds up single-threaded execution. They compose: a free-threaded build can also run subinterpreters and the JIT.

---

## 1. The shared runtime execution model (foundation — read first)

All three PEPs modify the same substrate. Understanding it once keeps the rest non-redundant.

**The eval loop has tiers (since 3.11–3.13):**
- **Tier 1 — specializing adaptive interpreter (PEP 659, "PEP 659" / 3.11+).** The bytecode interpreter rewrites hot bytecodes *in place* into type-specialized forms (e.g. `BINARY_OP` → `BINARY_OP_ADD_INT`) and collects profiling data (types, object layouts, hot paths). This is "free" speed with no compilation.
- **Tier 2 — micro-op (uop) IR (3.13+).** Hot regions are projected into a lower-level sequence of micro-ops. Enabled for experimentation with `-X uops` / `PYTHON_UOPS=1`. This is the IR the JIT consumes.
- **Tier 3 — JIT machine code (PEP 744, 3.13+).** Compiles tier-2 uops to native code (see §4).

**Per-interpreter and per-thread state.** Runtime state lives in three places: C globals, `_PyRuntimeState` (process-global), and `PyInterpreterState` (per interpreter). PEP 684's central engineering task was *moving* state out of globals/`_PyRuntimeState` into `PyInterpreterState` so interpreters can be isolated. Each thread has a `PyThreadState`.

**Immortal objects (PEP 683, 3.12).** Some objects (small ints, `None`, `True`/`False`, interned strings, code constants) are marked immortal: their refcount is never modified and they are never freed. This is the enabling primitive for *both* subinterpreters (immutable singletons can be shared across interpreters only if even the refcount never changes) *and* free-threading (avoids refcount cache-line contention on hot shared objects).

---

## 2. PEP 703 — Free-threaded (no-GIL) build

Adds a separate build of CPython where the GIL is disabled and the interpreter is made thread-safe, giving real multi-core parallelism for pure-Python threads.

### Status & timeline
- **3.13 (Oct 2024):** experimental free-threaded build.
- **3.14:** per **PEP 779**, free-threading is **officially supported but not the default** ("phase II"). You still opt in; the standard build keeps the GIL until the ecosystem catches up.

### Internals that change
- **Biased reference counting (BRC).** Fast path for objects owned by the current thread; slow (atomic) path for objects touched by other threads. Based on the observation that most objects are only ever touched by one thread.
- **Deferred reference counting.** Used for objects referenced constantly from many threads (module objects, functions, descriptors, `threading.local`): stack references aren't counted; cleanup is deferred to the GC. Reduces contention.
- **Per-thread reference counting** for some objects (heap types, code objects, module `__dict__`): each thread keeps a local count array, merged at safe points / `gc.collect()`.
- **Immortalization** (3.14): code constants and `sys.intern()`ed strings are immortal — never deallocated.
- **`mimalloc`** replaces `pymalloc` as the allocator; lock-free data structures use **QSBR** (quiescent-state-based reclamation), which delays freeing mimalloc "pages" until threads pass a quiescent state. `gc.collect()` forces deferred frees and per-thread count merges.
- **Per-object locks** and **`PyMutex`** (a 1-byte lock) protect mutable built-in containers (`list`, `dict`, `set`) so their operations stay atomic without a global lock.
- **Stop-the-world GC pauses** replace GIL-protected collection.
- **Object header grows**: non-GC objects like `None` go from 16→32 bytes on AMD64 (GC objects unchanged).

### Cost
Single-threaded overhead on the free-threaded build (pyperformance): ~1% on macOS aarch64, up to ~8% on x86-64 Linux. You trade single-thread speed for multi-core scaling.

### Building & detecting
```bash
./configure --disable-gil && make        # build from source
python -VV                                 # prints "free-threading build"
```
```python
import sys, sysconfig
sys._is_gil_enabled()                      # True/False at runtime
sysconfig.get_config_var("Py_GIL_DISABLED")  # 1 on a free-threaded build
```
Runtime GIL control (free-threaded build only):
```bash
python -X gil=0     # disable (default on FT build)   |  PYTHON_GIL=0
python -X gil=1     # force-enable                     |  PYTHON_GIL=1
```

### Behavioral changes on the FT build
- `sys.flags.thread_inherit_context` defaults **True** — new threads inherit a copy of the parent's `contextvars.Context`.
- `sys.flags.context_aware_warnings` defaults **True** — `warnings.catch_warnings` becomes thread-safe via context vars.
- `frame.f_locals` of a frame executing in another thread is **not** safe to read; sharing a single iterator across threads can drop/duplicate elements.

---

## 3. C-extension compatibility with free-threading

A C extension that doesn't declare free-thread support causes the interpreter to **auto-re-enable the GIL at import time with a warning** — silently negating no-GIL for the whole process. Declaring support is the single most important porting step.

**Multi-phase init (PEP 489) modules** — add a slot:
```c
static PyModuleDef_Slot slots[] = {
    {Py_mod_gil, Py_MOD_GIL_NOT_USED},   // "I am free-thread safe"
    {0, NULL}
};
```
**Single-phase init** (`PyModule_Create`) modules — call, guarded:
```c
#ifdef Py_GIL_DISABLED
    PyUnstable_Module_SetGIL(module, Py_MOD_GIL_NOT_USED);
#endif
```
- Wheels for the FT build use the **`cp314t`** ABI tag (the `t` = free-threaded). `abi3` / stable-ABI limited wheels need extra care.
- Ecosystem readiness: NumPy, Cython, pybind11, and PyO3 all ship FT-aware paths; track at `py-free-threading.github.io/tracking/` and `hugovk.github.io/free-threaded-wheels/`.
- Test under stress with a tiny `sys.setswitchinterval(...)` and validate native code with **ThreadSanitizer**.

---

## 4. PEP 734 — Multiple interpreters in the stdlib (subinterpreters)

CPython has had subinterpreters via the C API since 1.5 (1997); **PEP 684** (3.12) gave each interpreter its *own* GIL (per-interpreter GIL) by isolating runtime state into `PyInterpreterState`; **PEP 734** (3.14, final) exposes them to Python via the stdlib. This is "parallelism by isolation": each interpreter is a separate memory space with its own GIL, so N interpreters run on N cores — without the thread-safety hazards of the no-GIL build.

> PEP 734's predecessor draft was **PEP 554** (same idea, never landed). The module was renamed from `interpreters` to **`concurrent.interpreters`** before shipping.

### API (`concurrent.interpreters`)
```python
from concurrent import interpreters

interp = interpreters.create()          # -> Interpreter
interp.id                               # unique int
interp.is_running()
interp.prepare_main(x=10)               # bind objects into the target's globals
interp.exec("print(x)")                 # run source synchronously
interp.call(some_no_arg_callable)       # call a function (returns its result)
t = interp.call_in_thread(fn)           # run in a new threading.Thread
interp.close()

interpreters.get_current()              # the running interpreter
interpreters.list_all()                 # all live interpreters
```
**Cross-interpreter communication — a queue (not shared locks):**
```python
q = interpreters.create_queue(maxsize=0)
q.put(obj); q.get()                     # put/get/put_nowait/get_nowait/empty/full/qsize
```
**Sharing model:** nearly anything picklable can cross; most objects are *copied* (pickle), while `memoryview`/buffer-protocol objects share the underlying buffer directly. Interpreters are strictly isolated — synchronize by passing tokens through queues, not by sharing mutable objects.

**Exceptions:** `exec()` wraps uncaught errors in `ExecutionFailed` (with `.type`, `.msg`, `.snapshot`); `call()` propagates directly. Plus `InterpreterError`, `InterpreterNotFoundError`, `QueueEmpty`, `QueueFull`, etc.

**Pool:** `concurrent.futures.InterpreterPoolExecutor` — like `ThreadPoolExecutor` but each worker runs in its own subinterpreter+thread.

**Capability probe:** `sys.implementation.supports_isolated_interpreters` (bool).

### Caveats
Startup is heavier than a thread (each interpreter re-initializes); not every C extension is subinterpreter-safe (those relying on process-global C state). Best for CPU-bound, share-little workloads.

---

## 5. PEP 744 — The copy-and-patch JIT

A just-in-time compiler that turns hot **tier-2 micro-op** sequences into native machine code. New and experimental in 3.13; still off by default through 3.14.

**Copy-and-patch technique:** at *build time*, LLVM (Clang, needs `musttail`) compiles each micro-op into a machine-code **stencil**, dumped into a generated header. At *runtime*, the JIT compiles a hot uop sequence by copying each stencil's machine code almost verbatim and patching in the concrete operands — no per-run LLVM invocation, so JIT latency is tiny.

**Build / run:**
```bash
./configure --enable-experimental-jit          # yes | no | interpreter | yes-off
# yes-off  = build the JIT but run in interpreter mode unless enabled
```
```bash
PYTHON_JIT=1     # runtime toggle
```
- **Build-time** dependency: LLVM; adds ~3–60 s to the build. **Runtime** dependency: none, no API/ABI change.
- **Today's reality (3.13/3.14):** roughly *on par* with the specializing interpreter — "about as fast," with **10–20% memory overhead**. It's a foundation, not yet a free win.
- **Platforms:** x86-64 and aarch64 on Linux/macOS/Windows are tier-1.
- **Criteria to become non-experimental (PEP 744):** a meaningful (≈5%) speedup on at least one popular platform, deployable with minimal disruption, and Steering Council sign-off.

---

## 6. Choosing a parallelism / performance strategy

| Goal | Use | Why |
|------|-----|-----|
| CPU-bound, shared mutable state, threads | **Free-threading (PEP 703)** | Real shared-memory parallelism; but you own the locking, and FT-incompatible C ext re-enables the GIL |
| CPU-bound, little/no sharing, want isolation | **Subinterpreters (PEP 734)** | Per-interpreter GIL → multi-core, with isolation safety; pass data through queues |
| CPU-bound, hard isolation / crash containment | **`multiprocessing`** | Separate processes; highest overhead, strongest isolation |
| I/O-bound | **`asyncio` / threads** | GIL is released on I/O anyway → see `python-patterns` |
| Single-thread speed | **JIT (PEP 744)** + specializing interpreter | Transparent; modest gains today |

## Anti-patterns
- **Shipping a C extension without `Py_mod_gil`/`PyUnstable_Module_SetGIL`** → silently re-enables the GIL for the whole process on the FT build (with a warning many people miss).
- **Assuming no-GIL means thread-safe code.** Container *operations* are atomic (per-object locks), but multi-step invariants still need your own `threading.Lock`. Sharing one iterator across threads is unsafe.
- **Reading `frame.f_locals` of another thread's running frame** on the FT build — may crash.
- **Treating subinterpreters as cheap threads** — startup cost is real; don't spin one per tiny task (use `InterpreterPoolExecutor`).
- **Expecting the JIT to be a big speedup today** — it's ~parity with 10–20% memory overhead; benchmark before enabling.
- **Sharing mutable objects between subinterpreters** — only picklable copies and buffer-protocol shares are intended; reach for the queue.

## Troubleshooting
- *"GIL was re-enabled at runtime"* warning on the FT build → an imported C extension lacks the `Py_mod_gil` declaration; check `sys._is_gil_enabled()`, find the offender, set `PYTHON_GIL=0` only to confirm, then fix/replace the extension.
- *Objects not freed promptly* on the FT build → deferred/QSBR frees; call `gc.collect()`, or tune `MIMALLOC_PURGE_DELAY=0` (costs perf).
- *Subinterpreter import errors / crashes* → the extension keeps process-global C state and isn't multiple-interpreters-safe.
- *JIT shows no speedup* → expected at 3.13/3.14; verify it's actually built (`--enable-experimental-jit`) and on (`PYTHON_JIT=1`).

## References
- PEP 703 – Making the GIL Optional: https://peps.python.org/pep-0703/
- PEP 779 – Criteria for supported free-threaded Python: https://peps.python.org/pep-0779/
- Free-threading HOWTO (3.14): https://docs.python.org/3/howto/free-threading-python.html
- C-API extension support for free threading: https://docs.python.org/3/howto/free-threading-extensions.html
- Python Free-Threading Guide (porting, tracking): https://py-free-threading.github.io/
- PEP 734 – Multiple Interpreters in the Stdlib: https://peps.python.org/pep-0734/
- PEP 684 – A Per-Interpreter GIL: https://peps.python.org/pep-0684/
- PEP 554 – Multiple Interpreters (predecessor): https://peps.python.org/pep-0554/
- PEP 683 – Immortal Objects: https://peps.python.org/pep-0683/
- A per-interpreter GIL (LWN): https://lwn.net/Articles/941090/
- PEP 744 – JIT Compilation: https://peps.python.org/pep-0744/
- What is CPython's JIT Compiler (pydevtools): https://pydevtools.com/handbook/explanation/what-is-cpythons-jit-compiler/
- Following up on the Python JIT (LWN): https://lwn.net/Articles/1029307/
