<!-- hub-reference-banner -->
> **Reference file — part of the `lang-python` hub.** Formerly the standalone `python-static-type-checking` skill.
> Sibling topics in this family are now reference files under the hubs (`lang-python`, `lang-js-ts`, `lang-go-and-mobile`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

name: python-static-type-checking
description: >
  Python static type checkers and the gradual-typing model they implement — the four-checker landscape (mypy, Pyright/Pylance/BasedPyright, Astral ty, Meta Pyrefly) plus the shared PEP-defined foundation (PEP 483/484 gradual typing, PEP 561 stub distribution, the "gradual guarantee", the typing-spec conformance suite). Covers how each checker handles unannotated code, inference depth (declared-vs-inferred, the `list[int]` vs `list[Unknown]` divergence), strictness modes, config-file format ([tool.mypy]/[tool.pyright]/[tool.ty]/[tool.pyrefly]), plugins vs stubs, language-server/IDE story, performance (Rust checkers 10-80x faster), conformance scores, and tool-selection guidance.
  TRIGGER: choosing or configuring a Python type checker; "mypy vs pyright", "ty", "pyrefly", basedpyright, pylance, pyre; mypy.ini / setup.cfg / [tool.mypy] / [tool.pyright] / [tool.ty] / [tool.pyrefly]; --strict, disallow_untyped_defs, reportXxx rules, strictness levels; type stubs (.pyi, PEP 561, py.typed, types-* packages); mypy plugins (pydantic.mypy, sqlalchemy); "type checker is slow"/CI typing gate; reveal_type / # type: ignore / reportXxx suppression; gradual typing, the gradual guarantee, typing-spec conformance; migrating an untyped codebase to typed.
  SKIP: type-hint *syntax* itself (PEP 695 generics, Annotated, ParamSpec, TypeVar) → python-patterns; runtime validation/coercion of data → pydantic-v2 (BaseModel) or zod-schema-validation (TS); TypeScript's compiler/type system → typescript-expert / typescript-advanced-types; Node native TS type-stripping (which does NO type-checking) → nodejs-typescript-and-runtime-features.
version: "1.0"
category: developer
updated: "2026-05-31"
tags:
  - python
  - type-checking
  - mypy
  - pyright
  - ty
  - pyrefly
  - gradual-typing
  - static-analysis
---

# Python Static Type Checking — mypy, Pyright, ty, Pyrefly

Python type checkers are **external static-analysis tools**, not part of the interpreter. The CPython runtime ignores annotations (beyond storing them in `__annotations__`); a separate tool reads the same annotations a human reads and proves — or disproves — type consistency *before* the code runs. This file covers the shared model all four checkers implement, then each checker, then config/migration/anti-patterns.

As of mid-2026 the landscape is a **two-generation split**: the established Python-implemented checkers (**mypy**, **Pyright**) and the new **Rust-implemented** checkers (**Astral ty**, **Meta Pyrefly**) that are 10-80x faster. Pyrefly reached stable **1.0.0 (May 2026)**; ty is still **beta (0.x)**.

---

## 1. The shared foundation: Python's gradual-typing model

Every checker is an implementation of one specification, so learn the model once.

- **PEP 483 (Theory of Type Hints)** + **PEP 484 (Type Hints, 2015)** define *gradual typing*: type hints are **optional** and coexist with dynamic typing, so a codebase can be annotated incrementally without breaking. PEP 484 explicitly says the goal is **not** runtime enforcement — annotations exist for external tools (mypy was the reference implementation; Pyright et al. followed).
- **`Any` is the escape hatch.** `Any` is *consistent with* every type (assignable to and from anything). It is the formal seam between the static and dynamic worlds. Unannotated code is effectively `Any`-typed.
- **The "gradual guarantee".** A well-formed program should not gain *new* type errors merely because you **removed** an annotation (or, conversely, adding annotations should only ever *narrow* the set of errors, never invent new ones). This principle is the deepest design fault line between the checkers (see §3 inference divergence). ty and Pyright honor it strictly; mypy and Pyrefly infer more aggressively and can violate it.
- **PEP 561 (Distributing Type Information)** defines *how a package ships its types*: an inline-typed package drops a marker file **`py.typed`** in its root; third-party stub-only packages are named **`types-<pkg>`** (e.g. `types-requests`, `types-PyYAML`) or **`<pkg>-stubs`**; and it fixes the **resolution order** a checker uses to find types (stubs → inline → typeshed). Stub files are **`.pyi`**.
- **typeshed** is the community stdlib + third-party stub repository every checker consumes.
- **The typing-spec conformance suite** (python/typing repo) is the shared correctness benchmark. Pyrefly 1.0.0 reports **>90%** conformance, scoring above ty and mypy on the suite per its release notes; Pyright is historically the conformance leader since it tracks the spec closely.

**Two kinds of types a checker reasons about:** *declared* types (you wrote the annotation) vs *inferred* types (the checker deduced it from a literal or from later usage). The checkers agree on declared types and **disagree on inference** — that disagreement is the whole story in §3.

---

## 2. The four checkers at a glance

| | **mypy** | **Pyright** | **ty** | **Pyrefly** |
|---|---|---|---|---|
| Author | Python core / community | Microsoft | Astral (uv, ruff) | Meta |
| Language | Python (mypyc-compiled) | TypeScript/Node | **Rust** | **Rust** |
| Maturity (2026) | mature, reference impl | mature | **beta, 0.x** (no stable API) | **stable 1.0.0** (May 2026) |
| Default on unannotated funcs | **skips** them | **checks** them (inferred) | checks everything | checks aggressively |
| Inference philosophy | aggressive (`list[int]`) | conservative (`list[Unknown]`) | conservative, **gradual guarantee** | **aggressive** (`list[int]`) |
| Config table | `[tool.mypy]` / `mypy.ini` | `[tool.pyright]` | `[tool.ty]` | `[tool.pyrefly]` |
| Extensibility | **plugins** + stubs | **stubs only** (no plugins) | stubs (ruff/uv integration planned) | stubs |
| Language server | none built-in (3rd-party) | **best-in-class** (Pylance) | full LSP | full LSP (+ infer-to-source) |
| Relative speed | baseline (slow) | 2-3x slower than mypy on some, faster on big | **fastest** (~80x pyright) | ~2-3x slower than ty, still crushes mypy/pyright |
| Novel type features | — | early PEP support | **intersection & negation types** | strong generic inference |

CLI is uniform: `mypy <path>`, `pyright <path>`, `ty check <path>`, `pyrefly check <path>`.

---

## 3. The inference divergence (the single most important behavioral difference)

The canonical example. Given:

```python
my_list = []
my_list.append(1)
```

- **mypy, Pyrefly** analyze later usage and infer `list[int]`. A subsequent `my_list.append("foo")` is then an **error** (assumes homogeneous list).
- **Pyright, ty** infer `list[Unknown]` (ty) / partial-unknown, stay permissive, and **allow** `append("foo")` — honoring the gradual guarantee that already-working dynamic code shouldn't break under a checker.

Generic resolution shows the same split — `c: C[int] = C()` reveals `C[int]` under Pyrefly but `C[Unknown]` under ty. **Practical consequence:** mypy/Pyrefly catch more real bugs in *loosely-typed* code but produce more false positives and can violate the gradual guarantee; Pyright/ty produce fewer surprises during incremental adoption. Use `reveal_type(x)` (a checker-only pseudo-builtin) anywhere to print a variable's inferred type during a check.

---

## 4. mypy — the reference implementation

- **Default behavior:** *skips the bodies of unannotated functions.* This is the #1 gotcha — large swaths of an untyped codebase are silently unchecked until you turn on flags. `--check-untyped-defs` (or `--strict`) makes mypy check inside unannotated functions too.
- **Strictness is flag-by-flag.** `--strict` is a bundle that turns on all optional checks at once (the single most powerful config line). Key individual flags: `disallow_untyped_defs`, `disallow_incomplete_defs`, `disallow_any_generics`, `warn_return_any`, `warn_unused_ignores`, `no_implicit_optional`, `warn_redundant_casts`.
- **Config:** `mypy.ini`, `setup.cfg`, or `[tool.mypy]` in `pyproject.toml`. Per-module overrides via `[[tool.mypy.overrides]]` / `[mypy-<module>.*]` sections — this is how you do **strict-by-default, loose-for-legacy** gradual migration:

```toml
[tool.mypy]
strict = true
plugins = ["pydantic.mypy"]

[[tool.mypy.overrides]]
module = ["legacy.*", "vendor.*"]
disallow_untyped_defs = false
ignore_missing_imports = true
```

- **Plugins (unique to mypy).** A Python plugin API lets mypy *understand dynamic patterns* checkers can't see from types alone — ORM query builders, metaclasses, dataclass-like decorators. Workhorses: `pydantic.mypy`, `sqlalchemy.ext.mypy.plugin`. No other checker has this.
- **Missing third-party types:** `ignore_missing_imports = true` (or install the `types-*` stub package). mypy 2.0 added experimental parallel checking (`--num-workers`, ~1.3x).

---

## 5. Pyright — conservative, IDE-first

- **Checks all code by default**, including unannotated functions, using inferred types — so it surfaces obvious bugs with zero annotations (opposite of mypy's default).
- **Five strictness levels** set globally or per-file: `off`, `basic`, `standard` (default), `strict`, `all`. `strict` enables ~30 extra rules (annotations required on every param/return) and commonly produces a **~10x** jump in reported errors vs `basic`.
- **Per-rule control** via `reportXxx` keys, each settable to `"none" | "warning" | "error"` — e.g. `reportUnknownParameterType`, `reportMissingTypeStubs`, `reportUnusedImport`. Config in `[tool.pyright]` or `pyrightconfig.json`.

```toml
[tool.pyright]
typeCheckingMode = "strict"
reportUnknownMemberType = "warning"
```

- **No plugin system** — relies entirely on type stubs (a deliberate design choice for predictability).
- **Language server is the differentiator:** Pyright powers **Pylance** (Microsoft's closed-source VS Code extension, best-in-class hover/completion). **BasedPyright** is an open-source fork that adds Pylance-style features (and stricter defaults) without the license restriction. Inline suppression: `# pyright: ignore[reportXxx]`.

---

## 6. ty — Astral's Rust checker (beta)

- From the **uv + ruff** team; written in Rust; uses the **Salsa** fine-grained incremental framework (same engine family as rust-analyzer), so it recomputes only what changed — re-diagnosing a PyTorch file in **~4.7 ms (~80x faster than Pyright's ~386 ms)**.
- **Conservative inference + strict gradual guarantee:** uses an explicit **`Unknown`** type for inference uncertainty (so `[]` is `list[Unknown]`), and guarantees removing annotations never adds errors.
- **Novel type features:** the only checker (currently) with **intersection types and negation types** — after `if not isinstance(obj, MySubclass):`, `reveal_type(obj)` yields `MyClass & ~MySubclass`. Errors are deliberately **concise/structured** vs the verbose mypy/Pyright style.
- **Status:** **beta, `0.0.x` versioning, no stable API** — breaking changes (including to diagnostics) can land between any two releases. Astral dogfoods it internally and recommends it for *motivated* users willing to file issues. Config table `[tool.ty]`; full LSP with VS Code / Neovim / Zed / PyCharm extensions. Roadmap: integration with ruff/uv for type-aware linting and dead-code elimination.
- CLI: `ty check [path]`.

---

## 7. Pyrefly — Meta's Rust checker (stable 1.0.0)

- The Rust **successor to Pyre** (Meta's OCaml checker that ran Instagram). Open-sourced under **MIT, May 2025**; reached stable **1.0.0 in May 2026**. Custom incremental engine with **module-level** incrementalization + multi-threaded parallel checking. "Hard problems first" design — generics, overloads, wildcard imports built in from day one.
- **Aggressive inference** (same camp as mypy): infers return/local/container types, strong generic inference (`C()` → `C[int]`). This can violate the gradual guarantee, the trade-off for catching more in untyped code.
- **`pyrefly infer`** writes inferred annotations *directly into source files* (params, returns, container element types, plus the needed imports) — a migration accelerator no other checker offers.
- **Performance:** typechecks the ~20M-LOC Instagram codebase in **13.4 s** (vs 100+ s for Pyre); PyTorch in ~2.4 s vs Pyright ~35 s, mypy ~48 s. **Conformance >90%**, above ty and mypy. Deployed at Instagram; adoption at PyTorch and JAX. Full LSP with hyperlinked hover showing the *source* of an inference. Config table `[tool.pyrefly]`; CLI `pyrefly check [path]`.

---

## 8. Choosing a checker

- **Existing project with an established config / CI gate:** stay on **mypy** — mature, plugin ecosystem (pydantic, SQLAlchemy), most tutorials. Accept the slower runs.
- **Best IDE experience / VS Code:** **Pyright** (via Pylance) or **BasedPyright** for an OSS, stricter-by-default variant. Strong for early PEP support.
- **New project wanting maximum catch + speed and you're OK with aggressive inference:** **Pyrefly** (now stable 1.0.0, highest conformance, `pyrefly infer` for migration).
- **Want predictable incremental adoption (gradual guarantee) + the fastest checks + you tolerate beta churn:** **ty**.
- **CI strategy that scales:** run a Rust checker (ty/Pyrefly) as the fast pre-commit/local gate and optionally keep mypy in CI for plugin-dependent code, or commit fully to one. Don't run two strict checkers as *blocking* gates — their inference disagreements (§3) will fight each other.

---

## 9. Migrating an untyped codebase

1. **Start permissive, ratchet strictness.** Begin with non-strict, get a clean baseline, then enable strict flags incrementally.
2. **Strict-by-default, loose-for-legacy** via per-module overrides (mypy `[[tool.mypy.overrides]]`; Pyright per-file `# pyright: strict`). New modules pay the full tax; legacy modules are exempted explicitly.
3. **Fix missing third-party types:** install `types-<pkg>` stub packages or set `ignore_missing_imports` for the offending modules — don't blanket-ignore.
4. **Auto-annotate** with `pyrefly infer` (writes annotations to source) to bootstrap a large codebase, then review.
5. **Mark your own package typed:** ship a `py.typed` marker (PEP 561) so downstream checkers trust your inline annotations.
6. **Gate in CI** once a baseline is clean; track error count downward rather than requiring zero on day one.

---

## 10. Anti-patterns

- **Trusting mypy's default coverage.** mypy *skips unannotated function bodies* by default — a green run can mean "almost nothing was checked." Turn on `--check-untyped-defs` / `--strict` and verify with `reveal_type`.
- **Blanket `# type: ignore` / bare `reportXxx` suppression.** Always scope it: `# type: ignore[arg-type]` (mypy) / `# pyright: ignore[reportArgumentType]`. Set `warn_unused_ignores = true` so stale ignores get flagged when the underlying issue is fixed.
- **Running two strict checkers as blocking gates.** Their inference models disagree (§3); pick one source of truth.
- **Expecting runtime enforcement.** Type checkers never run your code or validate data at runtime. For *runtime* data validation use **Pydantic v2** (`references/pydantic-v2.md`) — annotations alone do nothing at runtime.
- **`Any` creep.** `Any` is consistent with everything and silently disables checking through it; prefer `object`, a `Protocol`, or a precise union. Use `disallow_any_generics` / `warn_return_any`.
- **Assuming ty's diagnostics are stable.** ty is `0.x` with no stable API — pin the version in CI and expect diagnostic changes between releases.

---

## 11. Troubleshooting

- **"Skipping analyzing '<pkg>': module is installed, but missing library stubs"** → install `types-<pkg>` or set `ignore_missing_imports` for that module; check the package ships `py.typed`.
- **"X works in mypy but errors in Pyright" (or vice-versa)** → almost always the inference divergence (§3) or mypy plugin behavior Pyright can't replicate (no plugins). Don't expect parity.
- **Slow mypy runs** → enable the mypy cache (incremental, default on), try `--num-workers` (2.0), or move the fast local gate to ty/Pyrefly.
- **Pydantic/SQLAlchemy models "untyped" under Pyright/ty/Pyrefly** → those rely on mypy *plugins*; on plugin-less checkers use the libraries' native typing support (Pydantic v2 is natively typed) rather than expecting plugin magic.
- **Floods of errors after enabling strict** → expected (Pyright strict ≈ 10x basic). Ratchet per-module, don't enable globally on day one.

---

## References

- PEP 483 — The Theory of Type Hints — https://peps.python.org/pep-0483/
- PEP 484 — Type Hints — https://peps.python.org/pep-0484/
- PEP 561 — Distributing and Packaging Type Information — https://peps.python.org/pep-0561/
- mypy — The mypy configuration file — https://mypy.readthedocs.io/en/stable/config_file.html
- mypy — Getting started — https://mypy.readthedocs.io/en/stable/getting_started.html
- Pyright — mypy comparison (official) — https://github.com/microsoft/pyright/blob/main/docs/mypy-comparison.md
- Astral ty (repo) — https://github.com/astral-sh/ty
- pydevtools — ty beta announcement — https://pydevtools.com/blog/ty-beta/
- pydevtools — How do mypy, pyright, and ty compare? — https://pydevtools.com/handbook/explanation/how-do-mypy-pyright-and-ty-compare/
- InfoQ — Meta Open-Sources Pyrefly — https://www.infoq.com/news/2025/05/meta-pyrefly-python-typechecker/
- pydevtools — Pyrefly reference — https://pydevtools.com/handbook/reference/pyrefly/
- Edward Li — Pyrefly vs. ty (deep comparison) — https://blog.edward-li.com/tech/comparing-pyrefly-vs-ty/
- InfoWorld — Pyrefly and Ty compared — https://www.infoworld.com/article/4005961/pyrefly-and-ty-two-new-rust-powered-python-type-checking-tools-compared.html
- Rob's Blog — How well do new Python type checkers conform? (ty, Pyrefly, Zuban) — https://sinon.github.io/future-python-type-checkers/
- Pyrefly — Why typed Python — https://pyrefly.org/blog/why-typed-python/
