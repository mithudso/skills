---
name: lang-python
description: >-
  Python language sub-hub (lang family). TRIGGER: Python idioms (3.12-3.14, pattern matching, exception groups, t-strings); async & structured concurrency (asyncio TaskGroup, anyio/trio, cancellation); CLI & TUI apps (argparse/Click/Typer, Rich, Textual); metaprogramming & data model (descriptors, metaclasses, Protocols/ABCs); AST/codegen (ast, LibCST, import hooks); data modeling (dataclasses/attrs/msgspec) & pydantic v2; pytest/Hypothesis testing; static type checking (mypy, pyright, ty, pyrefly); uv toolchain & packaging; app packaging/freezing (PyInstaller, Nuitka, shiv); logging & observability (structlog, loguru, OTel); CPython internals & profiling; supply-chain security; Python-in-browser/WASM. SKIP: JS/TS/Node/Deno/Bun → lang-js-ts; Go/Kotlin → lang-go-and-mobile; AI/ML frameworks → ai-* hubs; web frameworks (FastAPI/Django/Flask) → software-engineering-patterns; pandas/Polars/DuckDB → da-data-engineering-platform.
origin: local
---

> **Output rules:** No explanations — code only. Skip preamble. Don't recap, just proceed.

# lang-python

Python language sub-hub (lang family).

This hub routes to on-demand reference files under `references/`. See each spoke for depth.

## Sub-skill routing — load the reference that matches the task

| Reference | Use when the task is about |
| --- | --- |
| `references/python-patterns.md` | modern idioms (walrus, f-strings, pattern matching), exception groups/`except*`, t-strings, the decision-tree intros |
| `references/python-concurrency.md` | deep asyncio, structured concurrency (TaskGroup), cancellation/timeouts, anyio/trio, fire-and-forget pitfalls |
| `references/python-metaprogramming.md` | descriptors, metaclasses, `__init_subclass__`, the data model/dunders, Protocols vs ABCs, runtime introspection |
| `references/python-ast-codegen.md` | the `ast` module, LibCST codemods, import hooks, source transformation & code generation |
| `references/python-data-modeling.md` | choosing dataclasses vs attrs vs msgspec vs pydantic; NamedTuple/TypedDict/enum; serialization |
| `references/pydantic-v2.md` | pydantic v2 validation of untrusted/external data (BaseModel, validators, settings) |
| `references/python-cli-tui.md` | building CLIs (argparse/Click/Typer) and terminal UIs (Rich, Textual) |
| `references/python-app-packaging.md` | shipping standalone apps (PyInstaller, Nuitka, shiv/pex, zipapp, Briefcase) |
| `references/python-logging-observability.md` | the stdlib `logging` model, structlog, loguru, OpenTelemetry log correlation |
| `references/python-testing.md` | pytest (fixtures, parametrize, markers), Hypothesis property-based testing, coverage |
| `references/python-static-type-checking.md` | mypy / Pyright / ty / Pyrefly, the gradual-typing model, stubs, migrating a codebase |
| `references/uv-python-toolchain.md` | the uv toolchain — projects, the universal lockfile, Python install/pin, build/publish |
| `references/python-supply-chain-security.md` | pip-audit, SBOMs, bandit, sigstore/PEP 740 provenance, hash-pinned installs |
| `references/cpython-runtime-internals.md` | no-GIL free-threading, subinterpreters, the JIT, and the concurrency-model decision matrix |
| `references/cpython-performance-profiling.md` | profiling (cProfile/py-spy/Scalene/memray) and the native-acceleration ladder (Cython/PyO3/Numba) |
| `references/python-in-browser-wasm.md` | running CPython in the browser/WASM (Pyodide, PyScript, WASI) |

<!-- cross-hub-map -->
## Cross-hub map — where every lang topic lives

This family is split across these hubs. If a task's deep material is **not** in this hub's Sub-skill
routing table, it is a reference file under a sibling hub below — **activate that hub or `Read` its
`references/<name>.md` directly**. Every former standalone skill in this family is now a reference under one
of these hubs (nothing was deleted).

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `lang-python` | lang-python | `references/python-patterns.md`, `references/python-testing.md`, `references/python-static-type-checking.md`, `references/python-supply-chain-security.md`, … |
| `lang-js-ts` | lang-js-ts | `references/javascript-nodejs.md`, `references/javascript-runtimes-deno-bun-edge.md`, `references/javascript-node-html-css-debugging-expert.md`, `references/nodejs-concurrency-internals.md`, … |
| `lang-go-and-mobile` | lang-go-and-mobile | `references/go-patterns.md`, `references/compose-multiplatform-patterns.md` |
