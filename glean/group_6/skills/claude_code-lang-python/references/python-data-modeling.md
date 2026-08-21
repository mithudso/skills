# Python Data Modeling — dataclasses, attrs, msgspec

`lang-python` hub reference. Choosing how to model structured data in Python and
the performance/feature trade-offs across the four main tools.

**Scope:** this is the **decision and the lightweight options**. Runtime
validation of **untrusted external data** (coercion, JSON Schema, settings) is
`pydantic-v2` — this spoke covers when to reach for it vs. the lighter tools, and
the stdlib/`attrs`/`msgspec` alternatives. Class-creation machinery underneath
all of these is `python-metaprogramming`.

---

## 1. The decision in one line

> **Internal, trusted data** → `dataclasses` (or `attrs`). **External, untrusted
> data** needing validation/coercion → `pydantic`. **Hot path / millions of
> objects / fast serde** → `msgspec`.

Performance tiers (object init, approx.): **msgspec ~95 ns** < attrs ~315 ns ≈
dataclasses ~370 ns ≪ **pydantic ~1.8 µs**. Pydantic's cost buys validation;
don't pay it for internal containers.

---

## 2. `dataclasses` (stdlib) — the default for plain records

```python
from dataclasses import dataclass, field

@dataclass(slots=True, frozen=True, kw_only=True)   # slots: 3.10+, kw_only: 3.10+
class Point:
    x: int
    y: int = 0
    tags: list[str] = field(default_factory=list)   # NEVER tags: list = []
    def __post_init__(self): ...                      # post-construction hook
```

- Generates `__init__`/`__repr__`/`__eq__`; `frozen=True` adds immutability +
  hashability; `order=True` adds comparisons; `slots=True` drops `__dict__`.
- **`field()`**: `default_factory` (mutable defaults), `compare=`, `repr=`,
  `kw_only=`, `metadata=`.
- `asdict()` / `astuple()` for shallow conversion (no JSON, no validation).
- No runtime type checking — annotations are hints only. Use when the data is
  already trustworthy.

---

## 3. `attrs` — the powerful predecessor

`dataclasses` was inspired by `attrs`; `attrs` still offers more: composable
**validators**, **converters**, slots by default, and `cattrs` for fast
structuring/unstructuring (serde).

```python
from attrs import define, field, validators as v

@define                                   # slots + init/repr/eq, by default
class Account:
    email: str = field(validator=v.matches_re(r".+@.+"))
    balance: int = field(converter=int, default=0)
```

- `@define` (modern API) = slotted + sensible defaults; `@frozen` for immutable.
- **Validators/converters** run at init — lightweight validation without
  pydantic's weight. No native JSON/coercion of untrusted input → add **`cattrs`**
  (`structure`/`unstructure`) for serde.
- Best for performance-sensitive internal models that still want a little
  validation.

---

## 4. `msgspec` — fastest validation + serialization

`msgspec.Struct` is the speed king: ~4× faster init than dataclasses/attrs, ~17×
faster than pydantic, and decoding 10–20× faster than pydantic v2. Every `Struct`
is automatically `__slots__`-ed.

```python
import msgspec

class User(msgspec.Struct):
    id: int
    name: str
    roles: list[str] = []

data = msgspec.json.decode(raw_bytes, type=User)     # validate + parse in one pass
out  = msgspec.json.encode(data)                      # ultra-fast encode
```

- Built-in, schema-driven **JSON / MessagePack** encode/decode with type
  validation — like a lean, fast pydantic for high-throughput services and
  message parsing. (YAML and TOML support exists in recent versions but is
  experimental/version-dependent — verify against your target msgspec release.)
- Less ecosystem/feature breadth than pydantic (no rich `Field` constraints DSL,
  smaller validator vocabulary) — choose it when throughput dominates.

---

## 5. pydantic — when validation is the point

Reach for **`pydantic-v2`** (own spoke) when data crosses a trust boundary and
you need coercion, rich constraints, JSON Schema, settings management, or FastAPI
integration. Don't use it for internal data structures in tight loops — that's
the "pydantic everywhere" performance anti-pattern.

---

## 6. stdlib `typing` data shapes

- **`typing.NamedTuple`** — immutable, tuple-compatible, lightweight typed record;
  good for fixed return values.
- **`TypedDict`** — type a **dict's shape** (keys → value types) without changing
  runtime behavior; for JSON-ish dicts you don't want to wrap in a class.
  `Required`/`NotRequired` (PEP 655), `total=False`.
- **`enum`** — `Enum`, `IntEnum`, `StrEnum` (3.11+), `Flag`/`IntFlag` for
  bitfields; the right way to model closed sets of constants instead of bare
  strings.

---

## 7. Serialization (brief)

- **`json`** (stdlib) — universal, slow-ish; `default=`/`object_hook` for custom
  types.
- **`orjson`** — fast Rust-backed JSON (bytes in/out, handles dataclasses/datetimes); drop-in
  for speed when you don't need schema validation.
- **`msgspec`** — fastest *typed* serde (see §4); also MessagePack for binary.
- **`pickle`** — Python-only, **insecure on untrusted input** (arbitrary code
  execution on load); never unpickle data you didn't produce.

---

## 8. Choosing & pitfalls

| Situation | Pick |
| --- | --- |
| Internal record, no deps | **dataclasses** (`slots=True`) |
| Internal model + validators/converters, perf | **attrs** (+ cattrs for serde) |
| Untrusted input, coercion, JSON Schema, FastAPI | **pydantic-v2** |
| Hot path, millions of objects, fast JSON/msgpack | **msgspec** |
| Fixed typed tuple / dict shape / constant set | NamedTuple / TypedDict / enum |

Pitfalls:
- **Mutable default** (`x: list = []`) shares state across instances — always
  `field(default_factory=list)`.
- **`slots=True` + inheritance / `cached_property`** can conflict; test before
  enabling on a base class.
- `frozen=True` dataclasses are hashable only if all fields are hashable.
- Don't reach for pydantic by reflex — match the tool to whether the data is
  trusted and how hot the path is.

---

## Sources

- [dataclasses — Python stdlib docs](https://docs.python.org/3/library/dataclasses.html)
- [attrs documentation](https://www.attrs.org/) / [cattrs](https://catt.rs/)
- [msgspec — benchmarks](https://jcristharif.com/msgspec/benchmarks.html) and [docs](https://jcristharif.com/msgspec/)
- [dataclasses vs Pydantic vs attrs guide — TildAlice](https://tildalice.io/python-dataclasses-pydantic-attrs/)
- [msgspec vs Pydantic advantages — Hrekov](https://hrekov.com/blog/msgspec-vs-pydantic-advantages)
- [typing — NamedTuple/TypedDict](https://docs.python.org/3/library/typing.html) / [enum — Python docs](https://docs.python.org/3/library/enum.html)
