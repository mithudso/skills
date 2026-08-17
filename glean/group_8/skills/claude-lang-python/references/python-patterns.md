<!-- hub-reference-banner -->
> **Reference file — part of the `lang-python` hub.** Formerly the standalone `python-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`lang-python`, `lang-js-ts`, `lang-go-and-mobile`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: python-patterns
description: >
  Python patterns and best practices: modern idioms (3.12-3.14), async/await (asyncio, TaskGroup, Semaphore), PEP 695 type hints and generics, packaging (uv, hatch, pyproject.toml), ruff linter/formatter, pytest and Hypothesis property-based testing, dataclasses vs Pydantic v2, structural pattern matching, error handling (ExceptionGroup, except*), performance (free-threading, JIT), and anti-patterns.
  TRIGGER: user writes, reviews, or debugs Python code; asks about modern Python features (3.12, 3.13, 3.14, walrus operator, pattern matching); needs async/await guidance; asks about type hints, generics, PEP 695, pyright, mypy; needs packaging help (uv, hatch, pyproject.toml); configures ruff; writes pytest/Hypothesis tests; asks about dataclasses vs Pydantic; needs error handling or exception hierarchy design; asks about free-threading, the GIL, or profiling.
  SKIP: AI/ML frameworks (LangChain, PyTorch, TensorFlow) — use ai-languages; Django/Flask/FastAPI web framework specifics — use software-engineering-patterns (references/python-web-frameworks.md); MongoDB PyMongo patterns — use mongodb-developer; MCP server development with FastMCP — use mcp-servers; CI/CD pipeline configuration — use cicd-pipelines; Docker containerization — use docker-containers.
version: "1.1"
category: developer
updated: "2026-05-29"
tags:
  - python
  - async
  - type-hints
  - packaging
  - testing
  - pydantic
  - ruff
  - performance
  - anti-patterns
related_skills:
  - software-engineering-patterns
  - ai-agent-engineering
  - mongodb-expert
---

# Python Patterns and Best Practices

Comprehensive expert reference for modern Python development (3.12-3.14+). Covers language features, async patterns, type system, packaging, testing, data modeling, error handling, performance, and anti-patterns with actionable code examples and decision trees.

## When to use this skill

Activate when the user:

- writes, reviews, or debugs Python code and needs idiomatic patterns
- asks about modern Python features (3.12, 3.13, 3.14, walrus operator, pattern matching)
- needs async/await guidance, asyncio patterns, or concurrency architecture
- asks about Python type hints, generics, PEP 695, or static analysis (pyright, mypy)
- needs help with Python packaging (uv, hatch, pyproject.toml, publishing)
- configures ruff, asks about linting, or wants code quality tooling
- writes pytest tests, fixtures, parametrize, or Hypothesis property-based tests
- asks whether to use dataclasses, attrs, Pydantic, or NamedTuple
- needs error handling patterns, ExceptionGroup, or exception hierarchy design
- asks about Python performance, free-threading, the GIL, or profiling
- asks about structural pattern matching (match/case)
- wants to avoid Python anti-patterns or refactor legacy code

## When NOT to use this skill

- AI/ML frameworks (LangChain, PyTorch, TensorFlow) -- use `ai-languages` skill
- Django, Flask, FastAPI web framework specifics (ASGI/WSGI, DI, routing, ORM, deployment) -- see `software-engineering-patterns` → `references/python-web-frameworks.md`
- MongoDB Python driver (PyMongo) specifics -- use `mongodb-developer` skill
- MCP server development in Python (FastMCP) -- use `mcp-servers` or `mcp-builder` skill
- General coding conventions (naming, readability) without Python specifics -- use `coding-standards` skill
- CI/CD pipeline configuration -- use `cicd-pipelines` skill
- Docker containerization -- use `docker-containers` skill

## Cross-skill references

| Related skill | When to defer |
|---------------|---------------|
| `ai-languages` | Python AI/ML framework selection (LangChain, DSPy, etc.) |
| `mongodb-developer` | PyMongo driver patterns, connection pooling |
| `mcp-servers` | FastMCP server development in Python |
| `coding-standards` | Language-agnostic naming and readability conventions |
| `testing-and-vitest-expert` | JavaScript/TypeScript testing (not Python) |
| `backend-patterns` | Architecture patterns (CQRS, event sourcing) beyond Python |

---

## 1. Modern Python Idioms (3.12-3.14)

### 1.1 Walrus Operator (`:=`, PEP 572)

Use assignment expressions to reduce repeated computation:

```python
# GOOD: compute once, use the result in the condition
if (n := len(data)) > 10:
    print(f"Processing {n} items")

# GOOD: in while loops for sentinel patterns
while chunk := f.read(8192):
    process(chunk)

# GOOD: in list comprehensions with filtering
results = [clean for x in raw if (clean := normalize(x)) is not None]
```

**Anti-pattern**: Overusing `:=` in complex expressions -- readability always wins.

### 1.2 F-String Enhancements

Python 3.12 removed f-string limitations. Nested quotes and expressions are now unrestricted:

```python
# Python 3.12+: nested quotes, multiline expressions, backslashes all work
msg = f"Result: {"\n".join(items)}"
debug = f"Value: {obj["key"]!r}"
matrix = f"{'yes' if condition else 'no'}"

# Self-documenting expressions with =
print(f"{user.name=}, {user.age=}")  # user.name='Alice', user.age=30
```

### 1.3 Structural Pattern Matching (PEP 634-636, Python 3.10+)

```python
# Match on type + destructure
match command:
    case {"action": "move", "direction": str(d), "steps": int(n)}:
        move(d, n)
    case {"action": "quit"}:
        sys.exit(0)
    case _:
        print("Unknown command")

# Match with guards
match point:
    case Point(x, y) if x == y:
        print("On the diagonal")
    case Point(x, y) if x > 0 and y > 0:
        print("In quadrant I")

# OR patterns and capture
match status:
    case 200 | 201:
        handle_success()
    case 404:
        handle_not_found()
    case int(code) if 500 <= code < 600:
        handle_server_error(code)

# Matching dataclasses -- positional via __match_args__
@dataclass
class Point:
    x: float
    y: float

match origin:
    case Point(0, 0):
        print("Origin")
    case Point(x, 0):
        print(f"On x-axis at {x}")

# Matching Pydantic models -- named attributes
match pet:
    case Pet(species="dog", name=dog_name):
        print(f"Good dog, {dog_name}!")
```

**When to use**: Complex conditional dispatch, command parsing, protocol message handling, state machines. **Avoid**: Simple if/elif chains with 2-3 branches (use if/elif instead).

### 1.4 Exception Groups and `except*` (PEP 654, Python 3.11+)

```python
# Raise multiple exceptions from concurrent operations
async def fetch_all(urls):
    async with asyncio.TaskGroup() as tg:
        tasks = [tg.create_task(fetch(url)) for url in urls]
    # TaskGroup raises ExceptionGroup if any tasks fail

# Handle specific exception types from a group
try:
    results = await fetch_all(urls)
except* ConnectionError as eg:
    for exc in eg.exceptions:
        log.warning(f"Connection failed: {exc}")
except* TimeoutError as eg:
    for exc in eg.exceptions:
        retry_urls = [e.url for e in eg.exceptions]
except* ValueError as eg:
    raise  # re-raise unhandled types

# Python 3.14: PEP 765 -- disallow return/break/continue in finally
# that would silently suppress exceptions (enforced by linters)
```

### 1.5 Python 3.14 Highlights

- **Template strings (PEP 750)**: `t"Hello {name}"` -- creates `Template` objects for safe interpolation, preventing injection attacks in SQL/HTML
- **Deferred evaluation of annotations (PEP 649)**: Annotations stored as lazily evaluated functions, eliminating `from __future__ import annotations` need
- **Free-threading officially supported (PEP 779)**: No-GIL builds move from experimental to supported status
- **JIT compiler (Tier 2)**: Copy-and-patch JIT reduces overhead, especially for free-threaded builds
- **`except*` improvements**: Better integration with ExceptionGroup handling
- **Enhanced `finally` safety (PEP 765)**: Linters flag `return`/`break`/`continue` inside `finally` blocks
- **Incremental GC**: Thread-safe garbage collector that scans in small bursts, preventing latency spikes

---

## 2. Type Hints and Static Analysis

### 2.1 PEP 695 Generic Syntax (Python 3.12+)

The new syntax replaces verbose `TypeVar` declarations:

```python
# OLD (pre-3.12)
from typing import TypeVar, Generic
T = TypeVar("T")
class Box(Generic[T]):
    def __init__(self, value: T) -> None:
        self.value = value

def first(items: list[T]) -> T:
    return items[0]

# NEW (3.12+ PEP 695)
class Box[T]:
    def __init__(self, value: T) -> None:
        self.value = value

def first[T](items: list[T]) -> T:
    return items[0]

# Bounded type parameters
def serialize[T: (str, bytes)](data: T) -> T: ...

# Constrained generics with upper bounds
def process[T: Comparable](items: list[T]) -> T: ...
```

### 2.2 Type Aliases (PEP 695)

```python
# NEW: the `type` soft keyword
type Vector = list[float]
type Result[T] = tuple[T, Exception | None]
type JSON = dict[str, "JSON"] | list["JSON"] | str | int | float | bool | None
type Handler[**P, R] = Callable[P, Awaitable[R]]
```

### 2.3 Union Types with `|` (PEP 604, Python 3.10+)

```python
# Prefer | over Union
def process(data: str | bytes | None) -> str: ...

# isinstance checks
if isinstance(value, str | int): ...
```

### 2.4 Protocols and Structural Subtyping (PEP 544)

```python
from typing import Protocol, runtime_checkable

@runtime_checkable
class Renderable(Protocol):
    def render(self) -> str: ...

class HTMLWidget:
    def render(self) -> str:
        return "<div>Widget</div>"

# HTMLWidget satisfies Renderable without inheriting from it
def display(item: Renderable) -> None:
    print(item.render())

display(HTMLWidget())  # Works -- structural match
```

### 2.5 TypedDict, Required, NotRequired (PEP 655)

```python
from typing import TypedDict, Required, NotRequired

class UserConfig(TypedDict, total=False):
    name: Required[str]        # must be present
    email: Required[str]
    theme: NotRequired[str]    # optional
    locale: NotRequired[str]
```

### 2.6 Self Type (PEP 673, Python 3.11+)

```python
from typing import Self

class Builder:
    def set_name(self, name: str) -> Self:
        self.name = name
        return self

    @classmethod
    def create(cls) -> Self:
        return cls()
```

### 2.7 ParamSpec and Concatenate (PEP 612)

```python
from typing import ParamSpec, Concatenate, Callable

P = ParamSpec("P")

def with_logging[**P, R](func: Callable[P, R]) -> Callable[P, R]:
    def wrapper(*args: P.args, **kwargs: P.kwargs) -> R:
        print(f"Calling {func.__name__}")
        return func(*args, **kwargs)
    return wrapper
```

### 2.8 Static Analysis Tool Selection

| Tool | Speed | Inference | Use case |
|------|-------|-----------|----------|
| **pyright** | Fast (TS-based) | Best inference engine | New projects, VS Code |
| **mypy** | Moderate | Good, plugin ecosystem | Mature projects, Django plugin |
| **ruff** | Fastest | Lint rules only (no type checking) | Linting + formatting |

**Best practice**: Annotate function signatures, not local variables (pyright/mypy infer locals). Use `reveal_type()` for debugging inferred types.

---

## 3. Async/Await and Concurrency

### 3.1 When to Use What

| Pattern | Use case | GIL impact |
|---------|----------|------------|
| `asyncio` | I/O-bound (network, files, DB) | Single thread, no GIL contention |
| `threading` | Blocking libraries, legacy I/O | GIL limits CPU parallelism |
| `multiprocessing` | CPU-bound computation | Separate processes, no GIL |
| Free-threading (3.13t/3.14) | CPU-bound threads | No GIL, true parallelism |

### 3.2 Core Async Patterns

```python
import asyncio

# Basic coroutine
async def fetch_data(url: str) -> dict:
    async with aiohttp.ClientSession() as session:
        async with session.get(url) as response:
            return await response.json()

# TaskGroup (Python 3.11+) -- structured concurrency
async def fetch_all(urls: list[str]) -> list[dict]:
    results = []
    async with asyncio.TaskGroup() as tg:
        tasks = [tg.create_task(fetch_data(url)) for url in urls]
    return [t.result() for t in tasks]

# Always use asyncio.run() as entry point
if __name__ == "__main__":
    asyncio.run(main())
```

### 3.3 Rate Limiting with Semaphore

```python
async def fetch_with_limit(urls: list[str], max_concurrent: int = 10):
    semaphore = asyncio.Semaphore(max_concurrent)

    async def limited_fetch(url: str) -> dict:
        async with semaphore:
            return await fetch_data(url)

    async with asyncio.TaskGroup() as tg:
        tasks = [tg.create_task(limited_fetch(url)) for url in urls]
    return [t.result() for t in tasks]
```

### 3.4 Timeout Handling (Python 3.11+)

```python
# Modern timeout context manager
async def fetch_with_timeout(url: str) -> dict:
    async with asyncio.timeout(30.0):
        return await fetch_data(url)

# Reschedule deadline
async def adaptive_timeout():
    async with asyncio.timeout(10.0) as cm:
        data = await fast_operation()
        cm.reschedule(cm.when() + 20.0)  # extend if needed
        return await slow_operation(data)
```

### 3.5 Async Generators and Context Managers

```python
# Async generator for streaming
async def stream_events(channel: str):
    async with connect(channel) as conn:
        async for event in conn:
            yield event

# Async context manager
from contextlib import asynccontextmanager

@asynccontextmanager
async def managed_connection(url: str):
    conn = await create_connection(url)
    try:
        yield conn
    finally:
        await conn.close()
```

### 3.6 Async Iteration Patterns

```python
# Async comprehensions
results = [item async for item in aiter if item.valid]
total = sum([item.value async for item in aiter])

# async for with sentinel
async for message in websocket:
    if message.type == "close":
        break
    await handle(message)
```

### 3.7 Event Loop Best Practices

- **Always** use `asyncio.run()` -- never manage the event loop manually
- **Never** call `loop.run_until_complete()` from within a running loop
- **Store task references** -- untracked tasks can be garbage-collected mid-execution
- **Handle CancelledError** -- it is now a subclass of `BaseException` (3.9+)
- **Use `asyncio.to_thread()`** to call blocking functions from async code
- **Test with** `pytest-asyncio` and the `@pytest.mark.asyncio` decorator

---

## 4. Packaging and Tooling

### 4.1 uv -- The Modern Package Manager

uv (by Astral, written in Rust) replaces pip, pip-tools, virtualenv, and pyenv. It is 10-100x faster than pip.

```bash
# Install uv
curl -LsSf https://astral.sh/uv/install.sh | sh

# Project management
uv init myproject                  # Create new project
uv add requests httpx              # Add dependencies
uv add --dev pytest ruff mypy      # Add dev dependencies
uv remove requests                 # Remove dependency
uv lock                            # Generate/update lockfile
uv sync                            # Install from lockfile
uv run pytest                      # Run command in venv
uv run python script.py            # Run script

# Python version management
uv python install 3.14             # Install Python version
uv python pin 3.14                 # Pin for project

# Building and publishing
uv build                           # Build sdist + wheel
uv publish                         # Publish to PyPI
```

### 4.2 pyproject.toml -- Single Source of Truth

```toml
[project]
name = "mypackage"
version = "1.0.0"
description = "A modern Python package"
readme = "README.md"
license = "MIT"
requires-python = ">=3.12"
authors = [{ name = "Dev", email = "dev@example.com" }]
dependencies = [
    "httpx>=0.27",
    "pydantic>=2.0",
]

[project.optional-dependencies]
dev = ["pytest>=8.0", "ruff>=0.9", "pyright"]

[project.scripts]
myapp = "mypackage.cli:main"

[build-system]
requires = ["hatchling"]
build-backend = "hatchling.build"

[tool.hatch.build.targets.wheel]
packages = ["src/mypackage"]
```

### 4.3 uv.lock -- Deterministic Lockfile

- Cross-platform lockfile with exact resolved versions
- Committed to version control
- Replaces `requirements.txt` + `pip freeze` workflows
- Supports `--resolution lowest` for CI compatibility testing

### 4.4 Project Layout

```
myproject/
  pyproject.toml
  uv.lock                # committed to git
  src/
    mypackage/
      __init__.py
      core.py
      py.typed            # marker for type-checking consumers
  tests/
    conftest.py
    test_core.py
  .python-version         # pinned by uv python pin
```

### 4.5 Hatch vs Flit vs Setuptools vs PDM

| Backend | Strengths | Use case |
|---------|-----------|----------|
| **hatchling** | uv default, PEP-compliant, fast | New projects with uv |
| **flit-core** | Minimal, pure-Python only | Simple library packages |
| **setuptools** | Legacy support, C extensions | Projects needing setup.py compat |
| **pdm-backend** | PEP 582, lockfile | Alternative to uv |

---

## 5. Ruff -- Linter and Formatter

### 5.1 Configuration in pyproject.toml

```toml
[tool.ruff]
target-version = "py312"
line-length = 88

[tool.ruff.lint]
select = [
    "E",     # pycodestyle errors
    "W",     # pycodestyle warnings
    "F",     # pyflakes
    "I",     # isort
    "N",     # pep8-naming
    "UP",    # pyupgrade
    "B",     # flake8-bugbear
    "A",     # flake8-builtins
    "S",     # flake8-bandit (security)
    "T20",   # flake8-print
    "SIM",   # flake8-simplify
    "RUF",   # ruff-specific rules
]
ignore = [
    "E501",  # let formatter handle line length
]

[tool.ruff.lint.per-file-ignores]
"tests/**/*.py" = ["S101"]  # allow assert in tests

[tool.ruff.lint.isort]
known-first-party = ["mypackage"]

[tool.ruff.format]
quote-style = "double"
indent-style = "space"
docstring-code-format = true
```

### 5.2 Key Commands

```bash
ruff check .              # Lint
ruff check --fix .        # Lint + auto-fix
ruff format .             # Format (Black-compatible)
ruff check --select ALL . # Enable all rules (for auditing)
```

### 5.3 Important Rules

- **`select`** replaces defaults; **`extend-select`** adds to them
- **ISC001** (implicit string concat) can conflict with the formatter -- disable it if you see fighting
- Ruff replaces: flake8, Black, isort, pyupgrade, pydocstyle, autoflake, and 30+ plugins
- Performance: lints the entire pandas codebase in 0.5s vs flake8's 30s

---

## 6. Testing with pytest and Hypothesis

### 6.1 pytest Patterns

```python
# Fixtures with scope and cleanup
import pytest

@pytest.fixture(scope="module")
def db_connection():
    conn = create_connection()
    yield conn
    conn.close()

# Parametrize for table-driven tests
@pytest.mark.parametrize("input_val, expected", [
    ("hello", "HELLO"),
    ("", ""),
    ("  spaces  ", "  SPACES  "),
])
def test_uppercase(input_val, expected):
    assert input_val.upper() == expected

# Markers for test organization
@pytest.mark.slow
@pytest.mark.integration
def test_full_pipeline():
    ...

# Exception testing
def test_invalid_input():
    with pytest.raises(ValueError, match="must be positive"):
        process(-1)

# Approximate floating-point comparison
def test_calculation():
    assert compute() == pytest.approx(3.14159, rel=1e-5)
```

### 6.2 Fixture Composition

```python
@pytest.fixture
def user(db_connection):
    user = db_connection.create_user(name="test")
    yield user
    db_connection.delete_user(user.id)

@pytest.fixture
def authenticated_client(user):
    client = TestClient()
    client.login(user)
    yield client
    client.logout()
```

### 6.3 Async Test Support (pytest-asyncio)

```python
import pytest

@pytest.mark.asyncio
async def test_async_fetch():
    result = await fetch_data("https://api.example.com/data")
    assert result["status"] == "ok"

# Async fixtures
@pytest.fixture
async def async_db():
    db = await connect_db()
    yield db
    await db.close()
```

### 6.4 Hypothesis -- Property-Based Testing

```python
from hypothesis import given, strategies as st, assume, settings

# Basic property test
@given(st.lists(st.integers()))
def test_sort_is_idempotent(xs):
    assert sorted(sorted(xs)) == sorted(xs)

@given(st.lists(st.integers(), min_size=1))
def test_sort_preserves_length(xs):
    assert len(sorted(xs)) == len(xs)

# Composite strategies for complex data
@st.composite
def user_strategy(draw):
    name = draw(st.text(min_size=1, max_size=50, alphabet=st.characters(
        whitelist_categories=("L", "N", "Zs")
    )))
    age = draw(st.integers(min_value=0, max_value=150))
    email = draw(st.emails())
    return User(name=name, age=age, email=email)

@given(user_strategy())
def test_user_serialization_roundtrip(user):
    serialized = user.to_dict()
    restored = User.from_dict(serialized)
    assert restored == user

# Preconditions
@given(st.integers(), st.integers())
def test_division(a, b):
    assume(b != 0)
    assert (a * b) / b == pytest.approx(a)

# Settings for performance
@settings(max_examples=500, deadline=None)
@given(st.text())
def test_encode_decode_roundtrip(s):
    assert s.encode("utf-8").decode("utf-8") == s
```

### 6.5 Useful Hypothesis Strategies

```python
st.integers()                          # any int
st.floats(allow_nan=False)             # finite floats
st.text(min_size=1)                    # non-empty strings
st.emails()                            # valid email addresses
st.datetimes()                         # datetime objects
st.lists(st.integers(), max_size=100)  # bounded lists
st.dictionaries(st.text(), st.integers())  # dicts
st.builds(MyClass, name=st.text())     # build from class
st.from_type(MyDataclass)              # from type annotations
st.one_of(st.integers(), st.text())    # union types
st.sampled_from(["a", "b", "c"])       # enum-like
st.recursive(st.integers(), lambda children:  # recursive structures
    st.lists(children) | st.dictionaries(st.text(), children))
```

### 6.6 Coverage Configuration

```toml
[tool.coverage.run]
source = ["src"]
branch = true

[tool.coverage.report]
fail_under = 90
show_missing = true
exclude_lines = [
    "pragma: no cover",
    "if TYPE_CHECKING:",
    "if __name__ == .__main__.",
]
```

---

## 7. Data Modeling -- Dataclasses vs Pydantic vs attrs

### 7.1 Decision Tree

```
Does data cross a trust boundary (API, config, user input)?
  YES -> Pydantic v2
  NO  -> Is performance critical with minimal overhead?
    YES -> dataclasses (stdlib) or attrs (if you need validators)
    NO  -> Do you need JSON Schema / OpenAPI generation?
      YES -> Pydantic v2
      NO  -> dataclasses
```

### 7.2 Dataclasses (stdlib)

```python
from dataclasses import dataclass, field, asdict
from typing import ClassVar

@dataclass(frozen=True, slots=True)
class Point:
    x: float
    y: float

    def distance(self, other: "Point") -> float:
        return ((self.x - other.x)**2 + (self.y - other.y)**2) ** 0.5

@dataclass
class Config:
    host: str = "localhost"
    port: int = 8080
    tags: list[str] = field(default_factory=list)
    _internal: ClassVar[str] = "hidden"

    def __post_init__(self):
        if self.port < 0 or self.port > 65535:
            raise ValueError(f"Invalid port: {self.port}")
```

**Key flags**: `frozen=True` (immutable), `slots=True` (memory efficient, Python 3.10+), `kw_only=True` (keyword-only args, Python 3.10+), `match_args=True` (pattern matching support).

### 7.3 Pydantic v2

```python
from pydantic import BaseModel, Field, field_validator, model_validator
from pydantic import ConfigDict

class User(BaseModel):
    model_config = ConfigDict(
        strict=True,          # no coercion
        frozen=True,          # immutable
        str_strip_whitespace=True,
    )

    name: str = Field(min_length=1, max_length=100)
    email: str
    age: int = Field(ge=0, le=150)
    tags: list[str] = Field(default_factory=list)

    @field_validator("email")
    @classmethod
    def validate_email(cls, v: str) -> str:
        if "@" not in v:
            raise ValueError("Invalid email")
        return v.lower()

    @model_validator(mode="after")
    def check_consistency(self) -> "User":
        if self.age < 13 and "admin" in self.tags:
            raise ValueError("Underage admin not allowed")
        return self

# Serialization
user = User(name="Alice", email="ALICE@EXAMPLE.COM", age=30)
user.model_dump()          # -> dict
user.model_dump_json()     # -> JSON string
User.model_json_schema()   # -> JSON Schema

# Nested models
class Team(BaseModel):
    name: str
    members: list[User]
```

**Pydantic v2 performance**: Core rewritten in Rust (pydantic-core). Model creation is 5-50x faster than v1.

### 7.4 attrs (Lightweight Alternative)

```python
import attrs

@attrs.define(frozen=True, slots=True)
class Coordinate:
    x: float = attrs.field(validator=attrs.validators.instance_of(float))
    y: float = attrs.field(validator=attrs.validators.instance_of(float))
```

### 7.5 Comparison Table

| Feature | dataclasses | Pydantic v2 | attrs |
|---------|-------------|-------------|-------|
| Runtime validation | Manual (`__post_init__`) | Automatic | Via validators |
| JSON serialization | `asdict()` + json | Built-in | cattrs |
| JSON Schema | No | Yes | No |
| Performance | Fastest init | Fast (Rust core) | Very fast |
| Immutability | `frozen=True` | `frozen=True` | `frozen=True` |
| Slots | 3.10+ `slots=True` | Always slots | Default |
| Pattern matching | Automatic `__match_args__` | Named attributes | `__match_args__` |
| Stdlib | Yes | No (install) | No (install) |

---

## 8. Error Handling Patterns

### 8.1 Exception Hierarchy Design

```python
# Define a project exception hierarchy
class AppError(Exception):
    """Base exception for the application."""
    def __init__(self, message: str, code: str | None = None):
        self.code = code
        super().__init__(message)

class ConfigError(AppError):
    """Raised for configuration problems."""

class ValidationError(AppError):
    """Raised when input validation fails."""

class NotFoundError(AppError):
    """Raised when a resource is not found."""

class ExternalServiceError(AppError):
    """Raised when an external service call fails."""
    def __init__(self, message: str, service: str, status_code: int | None = None):
        self.service = service
        self.status_code = status_code
        super().__init__(message, code="EXTERNAL_ERROR")
```

### 8.2 Context Managers for Resource Safety

```python
from contextlib import contextmanager, suppress

# Custom context manager
@contextmanager
def managed_transaction(db):
    tx = db.begin()
    try:
        yield tx
        tx.commit()
    except Exception:
        tx.rollback()
        raise

# suppress specific exceptions
with suppress(FileNotFoundError):
    os.remove(tmp_path)
```

### 8.3 Exception Chaining

```python
# Preserve the original cause
try:
    data = parse_config(raw)
except json.JSONDecodeError as e:
    raise ConfigError(f"Invalid config format: {e}") from e

# Explicitly suppress the chain
try:
    result = legacy_function()
except OldError:
    raise NewError("Migrated") from None
```

### 8.4 Exception Notes (PEP 678, Python 3.11+)

```python
try:
    process(item)
except ValueError as e:
    e.add_note(f"Processing item: {item.id}")
    e.add_note(f"Batch: {batch_id}")
    raise
```

### 8.5 Anti-Patterns to Avoid

```python
# BAD: Bare except
try:
    do_work()
except:  # catches SystemExit, KeyboardInterrupt too!
    pass

# BAD: Catching too broadly
try:
    result = data["key"]
except Exception:
    result = default  # hides TypeError, AttributeError, etc.

# BAD: Silencing without logging
try:
    send_notification()
except Exception:
    pass  # failures are invisible

# BAD: Using exceptions for control flow
try:
    value = my_dict[key]
except KeyError:
    value = compute_default()
# GOOD:
value = my_dict.get(key) or compute_default()

# BAD: return/break/continue in finally (PEP 765)
try:
    result = risky_operation()
except SomeError:
    raise
finally:
    return None  # SILENTLY SUPPRESSES the exception!
```

### 8.6 Result Pattern (Functional Alternative)

```python
from dataclasses import dataclass
from typing import Generic

@dataclass(frozen=True)
class Ok[T]:
    value: T

@dataclass(frozen=True)
class Err[E]:
    error: E

type Result[T, E] = Ok[T] | Err[E]

def parse_int(s: str) -> Result[int, str]:
    try:
        return Ok(int(s))
    except ValueError:
        return Err(f"Cannot parse '{s}' as int")

# Usage
match parse_int(user_input):
    case Ok(value):
        process(value)
    case Err(msg):
        print(f"Error: {msg}")
```

---

## 9. Performance Optimization

### 9.1 Free-Threading (Python 3.13t / 3.14)

> **Deep dive:** for the full runtime-internals picture — biased/deferred reference counting, immortal objects (PEP 683), mimalloc/QSBR, `PyMutex`, subinterpreters (PEP 734 `concurrent.interpreters`, per-interpreter GIL PEP 684), the copy-and-patch JIT (PEP 744), C-extension free-thread compatibility (`Py_mod_gil`, `cp314t` wheels), and choosing free-threading vs subinterpreters vs multiprocessing — load `references/cpython-runtime-internals.md` in this hub.

```python
# Check if GIL is disabled
import sys
print(sys._is_gil_enabled())  # False in free-threaded builds

# True parallel threads for CPU-bound work
import threading

def cpu_work(data):
    return [x ** 2 for x in data]

threads = []
for chunk in split_data(data, n_threads=4):
    t = threading.Thread(target=cpu_work, args=(chunk,))
    threads.append(t)
    t.start()

for t in threads:
    t.join()
```

**Caveats**: Single-threaded perf penalty is 5-10% in 3.14. Not all C extensions are thread-safe yet. Check library compatibility.

### 9.2 Generators and Iterators for Memory

```python
# BAD: loads entire file into memory
lines = open("big.csv").readlines()

# GOOD: lazy iteration
def read_records(path):
    with open(path) as f:
        for line in f:
            yield parse_record(line)

# GOOD: itertools for pipeline processing
import itertools

first_100 = itertools.islice(read_records("big.csv"), 100)
batches = itertools.batched(read_records("big.csv"), 1000)  # Python 3.12+
```

### 9.3 `__slots__` for Memory Efficiency

```python
class Point:
    __slots__ = ("x", "y")

    def __init__(self, x: float, y: float):
        self.x = x
        self.y = y

# Or with dataclasses
@dataclass(slots=True)
class Point:
    x: float
    y: float
```

### 9.4 Caching and Memoization

```python
from functools import lru_cache, cache

@cache  # unbounded cache (Python 3.9+)
def fibonacci(n: int) -> int:
    if n < 2:
        return n
    return fibonacci(n - 1) + fibonacci(n - 2)

@lru_cache(maxsize=256)
def expensive_lookup(key: str) -> dict:
    return db.query(key)
```

### 9.5 Profiling Tools

```bash
# Built-in profiler
python -m cProfile -s cumtime my_script.py

# Line profiler (install: pip install line-profiler)
kernprof -l -v my_script.py

# Memory profiler
python -m memory_profiler my_script.py

# py-spy for production profiling (sampling, no overhead)
py-spy top -- python my_script.py
py-spy record -o profile.svg -- python my_script.py
```

### 9.6 String and Collection Performance

```python
# BAD: string concatenation in loop
result = ""
for item in items:
    result += str(item)

# GOOD: join
result = "".join(str(item) for item in items)

# Use collections for specialized needs
from collections import defaultdict, Counter, deque

word_counts = Counter(words)            # O(n) frequency counting
graph = defaultdict(list)               # auto-initializing dict
queue = deque(maxlen=1000)              # O(1) append/pop both ends

# Set operations for membership testing
valid_ids = frozenset(load_valid_ids()) # O(1) lookup, immutable
if user_id in valid_ids: ...
```

---

## 10. Itertools, Functools, and Standard Library Gems

### 10.1 itertools Patterns

```python
import itertools

# Group consecutive items
for key, group in itertools.groupby(sorted(items, key=keyfunc), keyfunc):
    process_group(key, list(group))

# Cartesian product
for combo in itertools.product(colors, sizes):
    create_variant(*combo)

# Pairwise (Python 3.10+)
for a, b in itertools.pairwise(sequence):
    check_transition(a, b)

# Batched (Python 3.12+)
for batch in itertools.batched(large_list, 100):
    process_batch(batch)

# Chain multiple iterables
all_items = itertools.chain(list_a, list_b, generator_c)

# Accumulate running totals
running_sums = list(itertools.accumulate(numbers))
```

### 10.2 functools Patterns

```python
from functools import partial, reduce, singledispatch, total_ordering

# Partial application
int_from_hex = partial(int, base=16)

# Single dispatch (function overloading by first arg type)
@singledispatch
def serialize(obj):
    raise TypeError(f"Cannot serialize {type(obj)}")

@serialize.register
def _(obj: str) -> bytes:
    return obj.encode("utf-8")

@serialize.register
def _(obj: dict) -> bytes:
    return json.dumps(obj).encode("utf-8")

# Total ordering -- define __eq__ and one comparison, get all six
@total_ordering
class Version:
    def __init__(self, major, minor, patch):
        self.major, self.minor, self.patch = major, minor, patch

    def __eq__(self, other):
        return (self.major, self.minor, self.patch) == (other.major, other.minor, other.patch)

    def __lt__(self, other):
        return (self.major, self.minor, self.patch) < (other.major, other.minor, other.patch)
```

---

## 11. Decorators and Metaprogramming

### 11.1 Decorator Patterns

```python
import functools
import time
from typing import Callable, ParamSpec, TypeVar

P = ParamSpec("P")
R = TypeVar("R")

# Decorator with arguments
def retry(max_attempts: int = 3, delay: float = 1.0):
    def decorator(func: Callable[P, R]) -> Callable[P, R]:
        @functools.wraps(func)
        def wrapper(*args: P.args, **kwargs: P.kwargs) -> R:
            for attempt in range(max_attempts):
                try:
                    return func(*args, **kwargs)
                except Exception:
                    if attempt == max_attempts - 1:
                        raise
                    time.sleep(delay * (2 ** attempt))
            raise RuntimeError("Unreachable")
        return wrapper
    return decorator

@retry(max_attempts=5, delay=0.5)
def fetch_api(url: str) -> dict: ...

# Class decorator
def singleton(cls):
    instances = {}
    @functools.wraps(cls)
    def get_instance(*args, **kwargs):
        if cls not in instances:
            instances[cls] = cls(*args, **kwargs)
        return instances[cls]
    return get_instance

# Decorator that works with and without arguments
def log(func=None, *, level="INFO"):
    def decorator(fn):
        @functools.wraps(fn)
        def wrapper(*args, **kwargs):
            print(f"[{level}] Calling {fn.__name__}")
            return fn(*args, **kwargs)
        return wrapper
    if func is not None:
        return decorator(func)
    return decorator

@log             # without arguments
def foo(): ...

@log(level="DEBUG")  # with arguments
def bar(): ...
```

### 11.2 Descriptor Protocol

```python
class Validated:
    def __init__(self, validator):
        self.validator = validator

    def __set_name__(self, owner, name):
        self.name = name
        self.storage_name = f"_{name}"

    def __get__(self, obj, objtype=None):
        if obj is None:
            return self
        return getattr(obj, self.storage_name, None)

    def __set__(self, obj, value):
        self.validator(value)
        setattr(obj, self.storage_name, value)

def positive(value):
    if value <= 0:
        raise ValueError("Must be positive")

class Product:
    price = Validated(positive)
    quantity = Validated(positive)
```

---

## 12. Common Anti-Patterns

### 12.1 Mutable Default Arguments

```python
# BAD
def append_to(item, target=[]):
    target.append(item)
    return target  # shared across calls!

# GOOD
def append_to(item, target=None):
    if target is None:
        target = []
    target.append(item)
    return target
```

### 12.2 Late Binding Closures

```python
# BAD: all lambdas capture the final value of i
funcs = [lambda: i for i in range(5)]
[f() for f in funcs]  # [4, 4, 4, 4, 4]

# GOOD: bind via default argument
funcs = [lambda i=i: i for i in range(5)]
[f() for f in funcs]  # [0, 1, 2, 3, 4]
```

### 12.3 Wildcard Imports

```python
# BAD
from module import *  # pollutes namespace, hides origin

# GOOD
from module import specific_function, SpecificClass
```

### 12.4 Ignoring Context Managers

```python
# BAD: resource leak risk
f = open("data.txt")
data = f.read()
# forgot f.close()

# GOOD
with open("data.txt") as f:
    data = f.read()
```

### 12.5 Type Checking at Runtime When Unnecessary

```python
# BAD: manual type checks in typed code
def process(data):
    if not isinstance(data, dict):
        raise TypeError("Expected dict")  # let the type checker handle this

# GOOD: annotate and trust the type system
def process(data: dict[str, Any]) -> None: ...
```

### 12.6 God Objects and Long Functions

```python
# BAD: 500-line function that does everything

# GOOD: decompose into focused functions
def process_order(order: Order) -> Receipt:
    validated = validate_order(order)
    priced = calculate_pricing(validated)
    receipt = finalize_order(priced)
    return receipt
```

### 12.7 LBYL vs EAFP

```python
# LBYL (Look Before You Leap) -- sometimes appropriate
if key in dictionary:
    value = dictionary[key]

# EAFP (Easier to Ask Forgiveness than Permission) -- Pythonic for hot paths
try:
    value = dictionary[key]
except KeyError:
    value = default

# But prefer .get() for simple dict lookups
value = dictionary.get(key, default)
```

---

## 13. Dependency Injection and Architecture

### 13.1 Protocol-Based DI

```python
from typing import Protocol

class EmailSender(Protocol):
    async def send(self, to: str, subject: str, body: str) -> None: ...

class SMTPSender:
    async def send(self, to: str, subject: str, body: str) -> None:
        # Real SMTP implementation
        ...

class MockSender:
    def __init__(self):
        self.sent: list[tuple] = []

    async def send(self, to: str, subject: str, body: str) -> None:
        self.sent.append((to, subject, body))

class NotificationService:
    def __init__(self, sender: EmailSender):
        self.sender = sender

    async def notify_user(self, user: User, message: str):
        await self.sender.send(user.email, "Notification", message)

# Production
service = NotificationService(SMTPSender())

# Testing
mock = MockSender()
service = NotificationService(mock)
await service.notify_user(user, "Hello")
assert len(mock.sent) == 1
```

### 13.2 Configuration with Environment Variables

```python
from pydantic_settings import BaseSettings
from pydantic import Field

class AppSettings(BaseSettings):
    model_config = ConfigDict(env_prefix="APP_")

    database_url: str
    redis_url: str = "redis://localhost:6379"
    debug: bool = False
    max_workers: int = Field(default=4, ge=1, le=32)
    allowed_origins: list[str] = ["http://localhost:3000"]

settings = AppSettings()  # reads from APP_DATABASE_URL, APP_REDIS_URL, etc.
```

---

## 14. Logging Best Practices

```python
import logging
import structlog

# stdlib logging -- structured with formatters
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(name)s %(message)s",
)
logger = logging.getLogger(__name__)
logger.info("Processing started", extra={"batch_id": batch_id})

# structlog -- structured logging (preferred for new projects)
structlog.configure(
    processors=[
        structlog.contextvars.merge_contextvars,
        structlog.processors.add_log_level,
        structlog.processors.TimeStamper(fmt="iso"),
        structlog.dev.ConsoleRenderer(),  # dev
        # structlog.processors.JSONRenderer(),  # prod
    ],
)
log = structlog.get_logger()
log.info("processing_started", batch_id=batch_id, item_count=len(items))
```

---

## 15. Quick Reference Checklists

### New Python Project Checklist

1. `uv init` with `pyproject.toml`
2. Set `requires-python = ">=3.12"`
3. Use `hatchling` build backend
4. Configure ruff (lint + format) in `pyproject.toml`
5. Set up pyright/mypy for type checking
6. Use `src/` layout with `py.typed` marker
7. Add pytest + hypothesis to dev dependencies
8. Pin Python version with `uv python pin`
9. Commit `uv.lock` to version control
10. Set up pre-commit hooks: `ruff check`, `ruff format --check`, `pyright`

### Code Review Checklist

- [ ] Type annotations on all public APIs
- [ ] No bare `except:` clauses
- [ ] No mutable default arguments
- [ ] Context managers for all resources (files, connections, locks)
- [ ] Generators for large data processing
- [ ] `__slots__` or `slots=True` on high-volume classes
- [ ] Specific exception types in except clauses
- [ ] `@functools.wraps` on all decorators
- [ ] No wildcard imports
- [ ] f-strings preferred over `.format()` or `%`

---

## 16. References

### Official Documentation

- [Python 3.14 What's New](https://docs.python.org/3/whatsnew/3.14.html)
- [Python 3.13 What's New](https://docs.python.org/3/whatsnew/3.13.html)
- [typing module reference](https://docs.python.org/3/library/typing.html)
- [PEP 695 -- Type Parameter Syntax](https://peps.python.org/pep-0695/)
- [PEP 636 -- Structural Pattern Matching Tutorial](https://peps.python.org/pep-0636/)
- [PEP 654 -- Exception Groups and except*](https://peps.python.org/pep-0654/)
- [PEP 8 -- Style Guide](https://peps.python.org/pep-0008/)
- [Free-Threading HOWTO](https://docs.python.org/3/howto/free-threading-python.html)
- [asyncio documentation](https://docs.python.org/3/library/asyncio.html)
- [Python Errors and Exceptions Tutorial](https://docs.python.org/3/tutorial/errors.html)

### Tooling

- [uv documentation](https://docs.astral.sh/uv/)
- [uv project guide](https://docs.astral.sh/uv/guides/projects/)
- [Ruff documentation](https://docs.astral.sh/ruff/)
- [Ruff configuration](https://docs.astral.sh/ruff/configuration/)
- [Ruff formatter](https://docs.astral.sh/ruff/formatter/)
- [pyright](https://github.com/microsoft/pyright)

### Testing

- [pytest documentation](https://docs.pytest.org/)
- [Hypothesis documentation](https://hypothesis.readthedocs.io/)
- [Hypothesis with pytest](https://pytest-with-eric.com/pytest-advanced/hypothesis-testing-python/)
- [Property-based testing guide](https://semaphore.io/blog/property-based-testing-python-hypothesis-pytest)

### Data Modeling

- [Pydantic v2 documentation](https://docs.pydantic.dev/latest/)
- [Pydantic v2 guide](https://devtoolbox.dedyn.io/blog/pydantic-complete-guide)
- [Python dataclasses guide](https://devtoolbox.dedyn.io/blog/python-dataclasses-guide)
- [attrs documentation](https://www.attrs.org/)

### Patterns and Practices

- [Real Python best practices](https://realpython.com/tutorials/best-practices/)
- [Modern Python 3.12+ features](https://dasroot.net/posts/2026/01/modern-python-312-features-type-hints-generics-performance/)
- [Python asyncio complete guide](https://dev.to/shehzan/mastering-python-async-patterns-a-complete-guide-to-asyncio-in-2026-10o6)
- [Modern Python best practices 2026](https://onehorizon.ai/blog/modern-python-best-practices-the-2026-definitive-guide)
- [Python packaging 2025: uv, Hatch](https://medium.com/@Modexa/python-packaging-2025-uv-hatch-and-the-end-of-it-works-locally-906264fc2aa5)
- [Ruff complete guide](https://pydevtools.com/handbook/explanation/ruff-complete-guide/)
- [Python error handling anti-patterns](https://github.com/charlax/antipatterns/blob/master/error-handling-antipatterns.md)
- [Enterprise Python error handling](https://www.augmentcode.com/guides/python-error-handling-10-enterprise-grade-tactics)
- [Async/await concurrent patterns](https://dasroot.net/posts/2026/04/mastering-python-async-await-patterns-concurrent-applications/)
- [Python 3.14 free-threading and JIT](https://blog.imseankim.com/python-3-14-free-threading-jit-compiler-gil-removal-2026/)
- [Astral blog: Python 3.14](https://astral.sh/blog/python-3.14)
