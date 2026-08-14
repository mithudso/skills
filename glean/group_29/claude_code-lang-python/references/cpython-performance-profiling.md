<!-- hub-reference-banner -->
> **Reference file — part of the `lang-python` hub.** Formerly the standalone `cpython-performance-profiling` skill.
> Sibling topics in this family are now reference files under the hubs (`lang-python`, `lang-js-ts`, `lang-go-and-mobile`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

name: cpython-performance-profiling
description: >
  CPython performance profiling and acceleration — the measure-then-optimize toolchain for Python
  hotspots. Covers deterministic profiling (cProfile/profile + pstats, SortKey, ncalls/tottime/cumtime,
  calibration, snakeviz/gprof2dot), statistical/sampling profilers that attach to a live process with
  no code changes (py-spy record/top/dump + --native/--gil/--subprocesses, Austin), Scalene (line-level
  CPU+GPU+memory profiling that separates Python vs native vs system time, copy-volume, AI suggestions),
  memray (Bloomberg allocation-level memory profiler with native tracking, flame/table/tree reporters,
  live mode, leaks/temporal, pytest-memray), line_profiler/kernprof line-level CPU, rigorous benchmarking
  (timeit, pyperf, pytest-benchmark), flame-graph interpretation (self vs cumulative), and the native
  acceleration ladder (Cython cdef/typed-memoryviews/nogil/prange, Numba @njit, mypyc, PyO3/Rust, ctypes/cffi).
  TRIGGER: a Python program is slow or memory-heavy and you need to find/fix the bottleneck; choosing a
  Python profiler (cProfile vs py-spy vs Scalene vs memray); profiling a running/production Python process;
  reading a flame graph; benchmarking Python code reliably; deciding whether/how to drop to a native
  extension (Cython, Numba, mypyc, PyO3) to speed up a hotspot.
  SKIP: language-agnostic debugging methodology / root-cause analysis → software-engineering-patterns
  (references/debugging-strategies.md); JS/web/Chrome-DevTools performance → software-engineering-patterns
  (references/performance-profiling-expert.md); the no-GIL/JIT/subinterpreter runtime mechanics themselves
  → references/cpython-runtime-internals.md; general Python idioms → references/python-patterns.md.
category: developer
tags:
  - developer
  - python
  - performance
  - profiling
related_skills:
  - programming-languages
  - software-engineering-patterns
version: "1.0.0"
updated: "2026-06-01"
---

# CPython Performance Profiling and Acceleration

Performance work in Python is two phases, in this order: **measure** (profile and benchmark to find
the real bottleneck) then **accelerate** (fix it, native last). The single most-violated rule is to
optimize before profiling — and the second is to trust measurements that the profiler's own overhead
has distorted. Pick the tool by the *question* you are asking; do not reach for a native rewrite
before the profile proves where the time goes.

## Pick the profiler by the question

| Question | Tool | Why |
| --- | --- | --- |
| "Which functions/call-paths cost the most?" | **cProfile + pstats** | Built-in, deterministic call graph, exact `ncalls`/`tottime`/`cumtime`. |
| "What is a *running / production* process doing right now?" | **py-spy** (or Austin) | Samples another process's memory; **no code changes, negligible overhead**. |
| "Is this hotspot Python, native, or I/O — and which *line*?" | **Scalene** | Line-level; separates Python vs native vs system vs GPU time. |
| "What is allocating memory (incl. C extensions)?" | **memray** | Tracks every allocation, Python + native + interpreter. |
| "Which *line* of this hot function is slow?" | **line_profiler / kernprof** | Per-line CPU timing inside `@profile` functions. |
| "Is implementation A faster than B, reliably?" | **pyperf / timeit / pytest-benchmark** | Statistically rigorous micro/macro benchmarks. |

## 1. Deterministic profiling — cProfile / profile / pstats

Deterministic profiling monitors **every** call/return/exception with precise interval timing. Use
the C-extension **`cProfile`** (low overhead); `profile` is the pure-Python, hookable, much slower
twin used mainly for extension/calibration.

```bash
python -m cProfile -o out.prof -s cumtime script.py     # -s sorts stdout; -o writes a file
python -m cProfile -m yourpackage.module                # profile a module
```

```python
import cProfile, pstats
from pstats import SortKey

with cProfile.Profile() as pr:        # context manager, Python 3.8+
    run_workload()
# also: pr.enable()/pr.disable(), pr.runcall(func, *args), cProfile.run("expr", "out.prof")

stats = pstats.Stats(pr).strip_dirs()
stats.sort_stats(SortKey.CUMULATIVE).print_stats(15)   # expensive call chains
stats.sort_stats(SortKey.TIME).print_stats(15)         # hot leaf functions
stats.print_callers(.5, "parse")                        # who calls the hot fn
stats.dump_stats("out.prof")
```

**Reading the columns** (the crux of cProfile):
- `ncalls` — call count (shows `total/primitive` when recursion is present).
- `tottime` — time **in the function itself, excluding** subcalls → **sort by this to find hot loops / leaf work**.
- `cumtime` — **cumulative**, function + everything it calls → sort by this to find expensive
  high-level operations.
- two `percall` columns — `tottime/ncalls` and `cumtime/primitive_calls`.

**SortKey** enum (3.7+; prefer over strings for error-checking): `CALLS`, `CUMULATIVE`, `FILENAME`,
`LINE`, `NAME`, `NFL`, `PCALLS`, `STDNAME`, `TIME`. `Stats` also supports `add()` (merge multiple
runs), `print_callees()`, `reverse_order()`, and `get_stats_profile()` (3.9+, structured dataclass).

**Calibration / bias** — the pure-Python `profile` (not cProfile) can subtract its own overhead:
`bias = profile.Profile().calibrate(10000)`, then `profile.Profile(bias=bias)` or set the class
`bias`. cProfile's overhead is small enough that calibration is rarely needed.

**Visualize** the `.prof`: **snakeviz** (interactive sunburst/icicle in the browser),
**gprof2dot** (`gprof2dot -f pstats out.prof | dot -Tsvg > out.svg` call-graph), or `tuna`.

**Tradeoff:** deterministic profiling perturbs timings by adding per-call overhead, which especially
distorts workloads dominated by many tiny calls. For long-running services and production, prefer a
sampling profiler.

## 2. Statistical / sampling profilers — py-spy, Austin

**py-spy** (Rust, by Ben Frederickson; rbspy lineage) is the go-to for **profiling a process you
cannot or do not want to instrument**, including production. It runs as a *separate process* and
**reads the target's memory** (`process_vm_readv` on Linux, `vm_read` on macOS, `ReadProcessMemory`
on Windows), extracting Python call stacks from interpreter structures — **zero code changes, very
low overhead**.

```bash
py-spy record -o profile.svg --pid 12345               # attach to a running PID → flamegraph SVG
py-spy record -o profile.json --format speedscope -- python prog.py   # spawn + speedscope JSON
py-spy top --pid 12345                                  # live top-style view of hot functions
py-spy dump --pid 12345                                 # print every thread's stack (find a hang)
```

Key flags: `--rate`, `--duration`, **`--native`** (interleave C/C++/Cython frames — Linux
i686/x86-64/ARM/Aarch64, Windows x86-64, macOS, FreeBSD), **`--gil`** (only threads holding the GIL —
true on-CPU Python time), `--subprocesses`, `--idle` (include sleeping threads),
`--nonblocking` (don't pause the target), `--locals` (locals in `dump`).
Output: flamegraph SVG (default), `--format speedscope` (time-ordered + left-heavy views in
speedscope.app), or `raw`.

**Permissions:** spawning a process is unprivileged; **attaching** to an existing one needs `sudo` /
ptrace on Linux (relax via `kernel.yama.ptrace_scope`), root on macOS, and the `SYS_PTRACE` capability
in Docker/Kubernetes.

**Austin** is a sibling frame-stack sampler (C, minimal dependencies) that emits collapsed stacks for
the same flamegraph/speedscope pipeline; it has VS Code and web-viewer integrations. Use it where you
want a tiny dependency-free sampler.

## 3. Scalene — line-level CPU + GPU + memory, with native separation

Scalene (plasma-umass) is a high-precision, low-overhead profiler whose distinguishing feature is
**separating time spent in Python vs. native (C/C++ libraries) vs. system (I/O/blocking)** — plus
GPU (NVIDIA) and memory — at **per-line and per-function** granularity. That separation immediately
answers "is this even optimizable in Python?": if a line is mostly *system* time it is I/O-bound and
no Python rewrite will help; if it is mostly *native* time the cost is inside a C library.

```bash
scalene run your_prog.py                 # profile → scalene-profile.json
scalene view --html                      # or --cli, --standalone, --json
scalene run your_prog.py --- --arg1      # pass args to the program after ---
```

- **Modes/flags:** `--cpu-only` (fastest), `--gpu`, `--memory`; `--reduced-profile` (only lines >1%
  CPU or >100 allocations); `--profile-only PATH` / `--profile-exclude PATH`; `--profile-interval N`;
  `--cpu-percent-threshold`, `--cpu-sampling-rate`, `--malloc-threshold`.
- **Targeting:** decorate functions with `@profile` (only those are reported), or programmatically
  `from scalene import scalene_profiler; scalene_profiler.start()/stop()`, or
  `with scalene.scalene_profiler.enable_profiling():`.
- **Copy volume (MB/s)** — a Scalene-original metric that flags costly *silent* copies when Python
  converts between C and Python representations or between CPU and GPU.
- **Low overhead via sampling + signal handlers + native stack stitching** (not per-line
  instrumentation): typically 10–20%, ~35% on heavy benchmarks — among the lowest for a profiler this
  detailed.
- **AI suggestions:** ⚡ (line) / 💥 (region) icons generate optimization proposals via
  Bedrock/Azure/OpenAI/Ollama. Experimental `--memory-leak-detector`.

## 4. memray — allocation-level memory profiling (Bloomberg)

memray tracks memory allocations in **Python code, native extension modules, and the interpreter
itself** by intercepting the allocators — the right tool for "what is eating memory" and "where do
numpy/pandas/C-extension allocations come from." **Linux and macOS only (no Windows).**

```bash
memray run [--native] [--follow-fork] [--trace-python-allocators] script.py   # → capture.bin
memray run --live script.py            # interactive terminal UI while it runs
memray flamegraph capture.bin          # interactive HTML flame graph (default reporter)
memray table  capture.bin              # HTML table of peak allocations
memray tree   capture.bin              # terminal tree
memray summary / stats capture.bin     # terminal high-level report
memray flamegraph --leaks capture.bin  # show allocations never freed (leak hunting)
memray flamegraph --temporal capture.bin   # time-resolved allocation view
```

- `--native` adds C/C++ frames (auto-consumed by every reporter — essential for numpy/pandas/etc.).
- `--trace-python-allocators` resolves allocations through Python's own allocator (pymalloc) for
  finer attribution.
- **pytest-memray:** add `--memray` to the pytest invocation for per-test memory reports; enforce
  caps with `@pytest.mark.limit_memory("100 MB")` (fails the test if exceeded).
- Default reporting shows the **high-watermark** (peak) allocation snapshot; `--leaks` switches to
  not-freed allocations, `--temporal` to over-time.

## 5. line_profiler / kernprof — per-line CPU

When a profiler points at a hot function and you need to know **which line** inside it costs the time
(typical for loops and numeric inner loops), use **line_profiler**:

```python
@profile                      # injected by kernprof; no import needed
def hot(): ...
```
```bash
kernprof -l -v script.py      # -l line-level, -v print results; writes script.py.lprof
python -m line_profiler script.py.lprof
```

It reports `Hits`, `Time`, `Per Hit`, `% Time`, and the source line. **py-heat** renders the same
line timings as a heatmap. line_profiler adds real overhead, so scope it to the one function you are
investigating.

## 6. Benchmarking — measure the fix, not the noise

- **timeit** — micro-snippets: `python -m timeit "your_code()"` or `timeit.timeit(stmt, number=10000)`.
  Good for one-liners; weak isolation.
- **pyperf** (PSF) — rigorous: runs the benchmark in **multiple processes**, warms up (skips the first
  value per process), reports **mean ± stdev**, does **not** disable the GC, and can **tune the
  system** (`pyperf system tune`) to suppress outliers from noisy neighbors. Use this for any "A vs B"
  claim that matters.
- **pytest-benchmark** — wires benchmarks into the test suite, tracks regressions, and compares runs;
  pair with CI **perf budgets** so a slowdown fails the build.

## 7. Flame-graph interpretation

- **Width = cost** (time for CPU flame graphs, bytes for memray). The x-axis is **not** time order in
  a classic flame graph — it is grouped/sorted stacks.
- **Self time** (the frame's own bar minus its children) vs **cumulative/total** (the whole stack
  width). A wide frame with narrow children = the work is *here*; a wide frame with wide children =
  the cost is *below* it.
- **speedscope** offers three views: *Time Order* (chronological), *Left Heavy* (merged, sorted —
  best for finding the biggest contributors), and *Sandwich* (callers/callees of one frame).
- py-spy `--gil` flame graphs show *real on-CPU Python*; without it, sleeping/IO-waiting threads
  inflate the picture.

## 8. The native-acceleration ladder

Climb in order; each rung is cheaper and safer than the next. **Native is the last resort**, only
after the profile proves a CPU-bound Python hotspot that algorithm/data-structure/vectorization work
cannot fix.

1. **Algorithm / data structure** — the biggest wins are almost always here (O(n²)→O(n log n),
   `set`/`dict` membership, generators to avoid materializing).
2. **Builtins / vectorization** — push loops into C: comprehensions, `str.join`, `itertools`,
   `functools`, and especially **NumPy** vectorized ops instead of Python-level loops.
3. **Concurrency** — `asyncio`/threads for I/O-bound; processes (or free-threaded 3.13t+, see
   `references/cpython-runtime-internals.md`) for CPU-bound parallelism.
4. **Native compilation** of the proven hotspot:

| Tool | Model | Best for | Cost / caveat |
| --- | --- | --- | --- |
| **Numba** `@njit`/`@jit` | LLVM **JIT** at runtime; infers types | Numerical / NumPy inner loops, `prange` parallel, `@vectorize` | First call pays compile cost; not a packaging-free win for arbitrary Python. |
| **Cython** | Python superset → C; `cdef`/`cpdef` types, **typed memoryviews**, `nogil`, **`prange()`** (OpenMP) | CPU-bound numeric/array code you want as a **shippable compiled wheel** (no end-user toolchain) | Build step; needs type annotations to actually be fast; parallel blocks must be `nogil`. |
| **mypyc** | Compiles **type-annotated** Python → C extension | Codebases already using type hints (e.g. mypy/black) — minimal source change | Gains scale with annotation coverage; subset of Python features. |
| **PyO3 / Rust** (+ **maturin**) | Rust compiled to a native module | New high-performance code, memory safety, releasing the GIL | Rust toolchain + bindings learning curve. |
| **ctypes / cffi** | Call an **existing** C library | Wrapping a pre-built native lib | No speedup unless most time is spent inside the C; FFI overhead per call. |

**Cython practical notes:** use **typed memoryviews** (`double[:, ::1] arr`) for fast contiguous
array access — they also unlock `nogil`. Wrap hot loops in `with nogil:` and use `prange(...,
nogil=True)` for OpenMP parallelism (the loop body must touch no Python objects). Run **`cython -a
file.pyx`** and open the HTML: **yellow lines mean residual Python-object interaction** — drive them
white to get C speed. **Pure-Python mode** (`import cython` + `@cython.cclass`/`cython.int`
annotations in a `.py`) keeps the source runnable as plain Python while still compiling.

## Anti-patterns and gotchas

- **Optimizing before profiling.** Intuition about Python hotspots is usually wrong; profile first.
- **Trusting overhead-distorted numbers.** cProfile inflates many-small-call workloads; line_profiler
  inflates the line under test. Cross-check hotspots with a sampling profiler (py-spy/Scalene) before
  committing to a rewrite.
- **Optimizing the wrong layer.** A line that is mostly *system* time (I/O) or *native* time (inside a
  C library) will not get faster from Python changes — Scalene's Python/native/system split is how you
  catch this.
- **Micro-benchmarking without warmup or isolation.** Single timeit runs on a noisy laptop mislead;
  use pyperf's multi-process + warmup + system tuning for decisions.
- **Reaching for native too early.** A native rewrite is expensive to build, ship, and maintain
  (compilers, wheels, ABI). Exhaust algorithm/vectorization/concurrency first.
- **Confusing wall-clock vs CPU time.** A "slow" function dominated by `sleep`/network is wall-clock,
  not CPU — don't try to JIT it.
- **Forgetting `--native`.** Without it, py-spy/Scalene/memray hide the C-extension frames where the
  real cost often lives (numpy/pandas/torch).
- **memray on Windows.** Not supported — use py-spy or tracemalloc there for memory questions.

## References (researched 2026-06-01)

- Python docs — The Python Profilers (cProfile/profile, columns, calibration): https://docs.python.org/3/library/profile.html
- Python docs — pstats / SortKey: https://docs.python.org/3/library/pstats.html
- py-spy (GitHub, benfred): https://github.com/benfred/py-spy
- Scalene (GitHub, plasma-umass): https://github.com/plasma-umass/scalene
- Berger et al., "Triangulating Python Performance Issues with Scalene" (arXiv): https://arxiv.org/pdf/2212.07597
- memray (GitHub / docs, Bloomberg): https://github.com/bloomberg/memray • https://bloomberg.github.io/memray/
- pytest-memray (GitHub): https://github.com/bloomberg/pytest-memray
- Cython — Parallelism (prange/OpenMP/nogil): https://cython.readthedocs.io/en/latest/src/userguide/parallelism.html
- Cython — Typed Memoryviews: https://docs.cython.org/en/latest/src/userguide/memoryviews.html
- scikit-learn — Cython Best Practices (annotation HTML): https://scikit-learn.org/stable/developers/cython.html
- pyperf (PSF docs / GitHub): https://pyperf.readthedocs.io/ • https://github.com/psf/pyperf
- Boost Python Performance with Cython, Numba, and PyO3 (Witt): https://wittgeo.medium.com/boost-python-performance-with-cython-numba-and-pyo3-486d59d8c2c6
