<!-- hub-reference-banner -->
> **Reference file — part of the `lang-python` hub.** Formerly the standalone `python-testing` skill.
> Sibling topics in this family are now reference files under the hubs (`lang-python`, `lang-js-ts`, `lang-go-and-mobile`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

name: python-testing
description: >
  Python testing expert — pytest (fixtures, parametrization, markers, configuration),
  fixture lifecycle and scopes, mocking (monkeypatch, pytest-mock, unittest.mock),
  Hypothesis property-based and stateful testing, async testing (pytest-asyncio, AnyIO),
  coverage (pytest-cov), parallelism (pytest-xdist), and the plugin ecosystem.
  TRIGGER: writing or debugging pytest tests; designing fixtures or conftest.py; parametrizing
  tests; mocking/patching in Python; property-based testing with Hypothesis (@given, strategies,
  @composite, stateful RuleBasedStateMachine); testing async code; measuring coverage; speeding
  up a slow suite; choosing a pytest plugin.
  SKIP: JavaScript/TypeScript testing with Vitest/Jest — use testing-and-vitest-expert; browser
  E2E testing — use webapp-testing; general Python idioms/packaging — use python-patterns;
  language-agnostic debugging methodology — use software-engineering-patterns (debugging-strategies).
version: "1.0"
category: developer
updated: "2026-05-31"
tags:
  - python
  - testing
  - pytest
  - hypothesis
  - property-based-testing
---

# Python Testing — pytest, Fixtures, and Hypothesis

Practical reference for testing Python with **pytest** (the de-facto standard), its
**fixture** system, **parametrization**, **mocking**, **Hypothesis** property-based
testing, and the async/parallel/coverage plugin ecosystem. For general Python idioms,
type hints, and packaging, defer to `references/python-patterns.md` in this hub. For
exact API and version details, the official pytest and Hypothesis docs are the source
of truth.

## Overview

pytest replaces `unittest`'s class-based boilerplate with plain `assert` statements
(rewritten for rich introspection), function-based tests, and a powerful
dependency-injection model via **fixtures**. The modern Python testing stack is:
`pytest` + `pytest-cov` (coverage) + `pytest-mock` (mocking) + `hypothesis`
(property-based) + `pytest-asyncio`/`anyio` (async) + `pytest-xdist` (parallel).
Install: `pip install pytest pytest-cov pytest-mock hypothesis pytest-asyncio pytest-xdist`.

## Core Concepts

### 1. Fixtures — lifecycle, scope, and conftest.py
Fixtures are setup/teardown providers requested by name as test-function arguments
(dependency injection). The decorator is `@pytest.fixture`.

- **yield fixtures** are the preferred finalization pattern: code before `yield` is
  setup, code after `yield` is teardown, and **teardown runs even if the test fails**.
  Use `yield` instead of `return` whenever cleanup is needed.
- **`request.addfinalizer(fn)`** is the older, more verbose alternative — register a
  callable that runs at teardown. Note pytest will run a registered finalizer even if
  the fixture later raises.
- **Scope** controls how often a fixture runs: `function` (default), `class`, `module`,
  `package`, `session`. Wider scope = fewer setups = faster suites, but shared mutable
  state across tests is a flakiness risk. Choose the narrowest scope that is still cheap.
- **`conftest.py`** provides fixtures to every test in its directory (and subdirectories)
  **without imports** — pytest auto-discovers them. Nested `conftest.py` files compose:
  child directories inherit parent fixtures and can add their own.
- **Finalization order is first-in-last-out (LIFO):** the right-most/last-requested
  fixture tears down first. With multiple yield fixtures, teardown unwinds in reverse.
- **Autouse:** `@pytest.fixture(autouse=True)` applies a fixture to every test in scope
  without it being requested by name — use sparingly (hidden dependencies).
- **Built-in fixtures:** `tmp_path` (per-test temp dir, `pathlib.Path`), `tmp_path_factory`
  (session-scoped temp dirs), `capsys`/`capfd` (capture stdout/stderr), `caplog` (capture
  logging), `monkeypatch` (safe patching), `request` (test context/introspection).

### 2. Parametrization
`@pytest.mark.parametrize("arg", [v1, v2, ...])` runs a test once per value, each as a
separate reported test case.

- **Multiple params:** `@pytest.mark.parametrize("a,b,expected", [(1,2,3),(4,5,9)])`.
- **Stacking** two parametrize decorators produces the Cartesian product.
- **Custom IDs:** pass `ids=[...]` (or `pytest.param(value, id="name")`) so opaque values
  show readable names in output and are selectable with `-k`. By default pytest escapes
  non-ASCII in IDs.
- **Per-case marks:** wrap a value in `pytest.param(v, marks=pytest.mark.xfail)` to mark
  individual cases (e.g., `xfail`/`skip`).
- **Indirect parametrization** (`indirect=True`, or a list of param names): route values
  through a fixture *before* the test receives them. Use this to defer expensive setup
  (DB connections, subprocesses) to run time rather than collection time.
- **Dynamic parametrization:** implement the `pytest_generate_tests(metafunc)` hook and
  call `metafunc.parametrize(...)` to compute parameters at collection time.

### 3. Markers and configuration
- **Markers** tag tests: `@pytest.mark.slow`, `@pytest.mark.skip(reason=...)`,
  `@pytest.mark.skipif(cond, reason=...)`, `@pytest.mark.xfail(reason=..., strict=True)`.
  Select with `-m "slow and not network"`.
- **Register custom markers** in config to avoid warnings (and `--strict-markers` to make
  unregistered markers an error).
- **Configuration** lives in `pyproject.toml` (`[tool.pytest.ini_options]`), `pytest.ini`,
  `tox.ini`, or `setup.cfg`. Common keys: `testpaths`, `addopts`, `markers`, `filterwarnings`,
  `asyncio_mode`. Prefer `pyproject.toml` for modern projects.

### 4. Mocking and patching
- **`monkeypatch`** (built-in fixture): safely set/delete attributes, dict items, env vars
  (`monkeypatch.setenv`), and `sys.path` entries; all changes auto-revert at test end. Best
  for patching values/env, not for assertion-rich mocks.
- **`unittest.mock`** (stdlib): `Mock`/`MagicMock`, `patch()` (context manager / decorator),
  `return_value`, `side_effect`, `assert_called_with`. **Patch where the name is looked up**,
  not where it is defined (`patch("mymodule.dependency")`, not the origin module).
- **`pytest-mock`** provides the `mocker` fixture — a thin, auto-undoing wrapper over
  `unittest.mock`. Less boilerplate than `monkeypatch` (set `return_value` in one line) and
  adds `mocker.spy` (wrap and record calls on a real object) and `mocker.stub`. Teams commonly
  combine `monkeypatch` (env/values) + `mocker` (call assertions/spying).

### 5. Property-based testing with Hypothesis
Instead of hand-picked examples, you specify **properties** that should hold for *all*
inputs, and Hypothesis generates many inputs from **strategies** to try to falsify them.

- **`@given(...)`** decorates a test and injects generated arguments drawn from strategies
  (e.g., `@given(st.integers(), st.text())`). Integrates natively with pytest.
- **Strategies** (`hypothesis.strategies as st`): primitives (`integers(min_value=, max_value=)`,
  `floats`, `text`, `booleans`, `none`), collections (`lists`, `dictionaries`, `sets`, `tuples`),
  and composites. Transform/compose them with `.map(fn)`, `.filter(pred)`, and `.flatmap(fn)`
  (draw, turn the drawn value into a new strategy, draw again).
- **`@composite`** builds a custom strategy from many others: the function receives a `draw`
  callable; use `assume(cond)` inside to discard examples that don't meet preconditions
  (prefer narrowing the strategy over heavy filtering, which can be slow).
- **`@example(...)`** pins specific must-test inputs alongside generated ones (regressions,
  known edge cases).
- **Shrinking:** when a property fails, Hypothesis reduces the failing input to the
  **simplest reproducer**, dramatically easing debugging.
- **`settings(...)`**: tune `max_examples`, `deadline` (per-example time limit; raises on
  slow examples), `phases`, `suppress_health_check`. Hypothesis keeps a **database** of
  failing examples so a once-failing case is replayed first on the next run (great in CI;
  configure the DB directory for shared/CI persistence).
- **Stateful testing** (`RuleBasedStateMachine`): declare `@rule`-decorated operations and
  `@invariant` checks; Hypothesis generates random *sequences* of operations against your
  system (API, data structure, DB) and shrinks both the inputs and the operation order to a
  minimal failing trace. Ideal for systems with many interacting operations.
- **Advanced designs:** differential testing (compare against a reference implementation),
  metamorphic testing (relate outputs of related inputs), and targeted/exploration testing.

### 6. Async testing
- **`pytest-asyncio`:** mark coroutine tests with `@pytest.mark.asyncio`, or set
  `asyncio_mode = "auto"` in config to auto-collect `async def` tests. Supports async fixtures.
- **AnyIO:** mark with `@pytest.mark.anyio` and define an `anyio_backend` fixture returning
  `"asyncio"` (or `"trio"`) — runs the same tests across backends. **Conflict:** AnyIO auto
  mode clashes with pytest-asyncio auto mode; if using AnyIO, remove `asyncio_mode` (use
  pytest-asyncio strict mode) to avoid double-collection.
- Trio's structured concurrency (nurseries) prevents leaked tasks; AnyIO lets you test both
  backends from one suite.

### 7. Coverage and parallelism
- **`pytest-cov`** (wraps coverage.py): `pytest --cov=mypackage --cov-report=term-missing`
  shows uncovered lines; `--cov-report=html`/`xml` for reports/CI. Gate with
  `--cov-fail-under=N`. Branch coverage via `--cov-branch`.
- **`pytest-xdist`** runs tests across CPUs/hosts: `pytest -n auto` splits into worker
  processes. Use it for CPU-bound, **independent** tests; serialize tests that share mutable
  state or race. Note coverage needs combining across workers (pytest-cov handles this).
- **Concurrency rule of thumb:** sequential if tests race; `pytest-xdist` for CPU-bound
  non-async; `pytest-asyncio`/AnyIO (+ `pytest-subtests`) for IO-bound async.

## Practical Patterns

- **Arrange-Act-Assert** per test; one logical assertion target per test.
- Put shared fixtures in the nearest `conftest.py`; keep test-specific setup in the test file.
- Factory fixtures: return a function (a "fixture factory") when a test needs multiple
  configured instances.
- Use `tmp_path`/`tmp_path_factory` for filesystem tests — never write to the repo or CWD.
- Combine `parametrize` + Hypothesis: parametrize over modes/configs, let `@given` fuzz inputs.
- Pin Hypothesis regressions with `@example`; commit the Hypothesis DB or persist it in CI cache.
- Default CI invocation: `pytest -n auto --cov=pkg --cov-report=term-missing --strict-markers`.

## Anti-Patterns

- **Over-wide fixture scope** with mutable shared state → cross-test contamination and flaky,
  order-dependent failures.
- **Overusing `autouse`** → hidden dependencies that make tests hard to reason about.
- **Patching at the definition site** instead of the lookup site → the mock never takes effect.
- **Heavy `assume()`/`.filter()`** in Hypothesis → many discarded examples and slow tests;
  narrow the strategy instead.
- **Asserting on randomly generated values** in property tests → assert *properties/invariants*,
  not specific outputs.
- **Mixing AnyIO auto mode with pytest-asyncio auto mode** → double collection/conflicts.
- **`xfail` without `strict=True`** → silently passes when the bug is fixed, hiding regressions.
- **Catch-all `try/except` swallowing assertion errors** inside tests.

## Troubleshooting

- *Fixture "not found":* it must be in the test file, an imported plugin, or a `conftest.py`
  at/above the test's directory; check name spelling and scope visibility.
- *Teardown not running:* ensure you used `yield` (not `return`) or `addfinalizer`; a setup
  exception before `yield` skips teardown.
- *Async test "coroutine was never awaited" / skipped:* missing `@pytest.mark.asyncio` or
  `asyncio_mode`; for AnyIO, missing `anyio_backend` fixture.
- *Hypothesis `DeadlineExceeded`:* raise/disable `deadline` in `settings` for legitimately slow
  code, or speed up the example generation.
- *Flaky under `-n auto`:* tests share global/file/DB state — isolate with per-worker fixtures
  or serialize the offending tests.
- *Coverage lower than expected with xdist:* ensure pytest-cov is combining worker data
  (it does by default; verify no custom `.coveragerc` parallel misconfig).

## References

- pytest — How to use fixtures: https://docs.pytest.org/en/stable/how-to/fixtures.html
- pytest — Fixtures reference: https://docs.pytest.org/en/stable/reference/fixtures.html
- pytest — Parametrizing tests: https://docs.pytest.org/en/stable/how-to/parametrize.html
- pytest — monkeypatch/mock: https://docs.pytest.org/en/stable/how-to/monkeypatch.html
- pytest-mock (PyPI): https://pypi.org/project/pytest-mock/
- Hypothesis (GitHub): https://github.com/HypothesisWorks/hypothesis
- Hypothesis — Strategies reference: https://hypothesis.readthedocs.io/en/latest/reference/strategies.html
- Hypothesis — settings: https://hypothesis.readthedocs.io/en/latest/settings.html
- pytest-asyncio (PyPI): https://pypi.org/project/pytest-asyncio/
- AnyIO — testing: https://anyio.readthedocs.io/en/stable/testing.html
- Property-based testing guide (stateful/differential/metamorphic): https://www.marktechpost.com/2026/04/18/a-coding-guide-for-property-based-testing-using-hypothesis-with-stateful-differential-and-metamorphic-test-design/
- Hypothesis + pytest robust testing (Pytest with Eric): https://pytest-with-eric.com/pytest-advanced/hypothesis-testing-python/
- pytest plugins overview: https://pytest-with-eric.com/pytest-best-practices/pytest-plugins/
