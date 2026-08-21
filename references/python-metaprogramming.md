# Python Metaprogramming & the Data Model

`lang-python` hub reference. The class-creation machinery and protocol model that
frameworks (Django ORM, SQLAlchemy, pydantic, dataclasses, attrs) are built on:
descriptors, the class-creation pipeline, `__init_subclass__`, metaclasses,
Protocols/ABCs, and runtime introspection.

**Scope:** decorators, closures, and everyday idioms are in `python-patterns`;
**static** typing of these constructs is in `python-static-type-checking`;
runtime *validation* models are `pydantic-v2` / `python-data-modeling`. This
spoke is the deep "code that customizes how classes are built and attributes are
accessed" layer.

---

## 1. The data model (dunder protocols)

Python behavior is defined by **special methods**; objects "are" what protocols
they implement (duck typing made concrete).

- **Attribute access:** `__getattr__` (only on miss), `__getattribute__` (every
  access; careful, recursion-prone), `__setattr__`, `__delattr__`, `__dir__`.
- **Callable / container / number protocols:** `__call__`, `__len__`,
  `__getitem__`/`__setitem__`, `__iter__`/`__next__`, `__contains__`, `__add__`,
  `__eq__`/`__hash__`, `__enter__`/`__exit__`, etc.
- **`__slots__`** — declare a fixed attribute set; drops the per-instance
  `__dict__` for big memory savings and faster attribute access. Trade-off: no
  dynamic attributes, multiple-inheritance constraints.

```python
class Proxy:
    def __init__(self, target): object.__setattr__(self, "_t", target)
    def __getattr__(self, name):           # only called when normal lookup fails
        return getattr(self._t, name)
```

---

## 2. Descriptors — the attribute-customization protocol

A **descriptor** is an object implementing `__get__`/`__set__`/`__delete__`. Put
one on a *class* and it intercepts attribute access on *instances*. This is the
mechanism behind `property`, methods (functions are descriptors → bound methods),
`classmethod`, `staticmethod`, and every ORM/typed-field system.

```python
class Field:                                # a typed, validating descriptor
    def __set_name__(self, owner, name):    # 3.6+: learns its own attribute name
        self.attr = f"_{name}"
    def __get__(self, obj, objtype=None):
        if obj is None: return self
        return getattr(obj, self.attr)
    def __set__(self, obj, value):
        if not isinstance(value, int): raise TypeError(value)
        setattr(obj, self.attr, value)

class Model:
    age = Field()                            # __set_name__(Model, "age") fires at class creation
```

- **Data descriptor** (defines `__set__`/`__delete__`) takes precedence over the
  instance `__dict__`; **non-data descriptor** (only `__get__`) is shadowed by an
  instance attribute. This precedence rule governs all attribute lookup.
- `__set_name__` removes the old "repeat the name" boilerplate
  (`age = Field("age")`).

---

## 3. The class-creation pipeline

Knowing the order is what lets you hook the right stage:

1. The class **body executes** in a namespace (a `metaclass.__prepare__` mapping,
   default `dict`).
2. `metaclass.__new__` builds the class object; `type.__new__` then **calls every
   descriptor's `__set_name__`**.
3. `__init_subclass__` is called on the **immediate parent** of the new class.
4. `metaclass.__init__` runs.
5. Class decorators (if any) wrap the finished class last.

Pick the lightest hook that does the job: class decorator < `__init_subclass__` <
metaclass.

---

## 4. `__init_subclass__` — the metaclass you usually want instead

A classmethod-like hook (3.6+) that runs whenever a class is subclassed. It
covers most "register/validate/inject on subclassing" needs **without** a
metaclass, and composes cleanly under inheritance.

```python
class Plugin:
    registry: dict[str, type] = {}
    def __init_subclass__(cls, /, key: str, **kw):   # subclass kwargs!
        super().__init_subclass__(**kw)
        Plugin.registry[key] = cls

class Csv(Plugin, key="csv"): ...          # auto-registered, no metaclass
```

Reach for it before metaclasses for: subclass registration, enforcing required
attributes/methods, default injection, validation at definition time.

---

## 5. Metaclasses — only when you must

A metaclass is the **class of a class** (default `type`). Use one only when you
need to control class creation *itself* in ways `__init_subclass__` and
descriptors cannot handle: e.g., rewriting the namespace via `__prepare__`, intercepting
*the class's own* instantiation via `__call__`, or affecting a class that has no
common base you control.

```python
class Singleton(type):
    _instances: dict = {}
    def __call__(cls, *a, **kw):            # intercepts MyClass(...)
        if cls not in cls._instances:
            cls._instances[cls] = super().__call__(*a, **kw)
        return cls._instances[cls]

class Config(metaclass=Singleton): ...
```

- `__new__(mcs, name, bases, namespace, **kw)` builds; `__init__` configures;
  `__call__` controls instance creation. `__prepare__` chooses the namespace type
  (e.g. ordered/enum behavior).
- Metaclass conflicts: a class's metaclass must be a (non-strict) subclass of all
  its bases' metaclasses; mixing two metaclassed hierarchies fails. This
  fragility is *why* `__init_subclass__`/`__set_name__` were added.

---

## 6. Protocols and ABCs — interfaces two ways

- **`typing.Protocol`** — **structural** typing ("static duck typing"): a class
  satisfies a Protocol by *shape*, no inheritance needed. Static checkers verify
  it; `@runtime_checkable` enables `isinstance` (shallow: method presence only).
  Prefer Protocols for interfaces in typed code.
  ```python
  from typing import Protocol, runtime_checkable
  @runtime_checkable
  class Closeable(Protocol):
      def close(self) -> None: ...
  ```
- **`abc.ABC` / `ABCMeta`** — **nominal** typing: subclasses must register or
  inherit; `@abstractmethod` blocks instantiation until overridden.
  `__subclasshook__` customizes `issubclass`. Use ABCs when you want enforced
  inheritance, abstract methods, or virtual-subclass registration.

Protocol = "anything shaped like this"; ABC = "anything declared to be this."

---

## 7. Runtime introspection

- **`inspect`** — `signature()`, `getmembers()`, `getsource()`, `get_annotations()`;
  the reflection toolkit.
- **`typing.get_type_hints(obj)`** — resolves annotations to real types
  (handles string/forward refs and PEP 563/649 deferred annotations), which is
  what pydantic/dataclasses read to build models.
- **`__annotations__`** — the raw annotation dict on classes/functions/modules.
  In 3.14, `annotationlib` + PEP 649/749 make annotations lazily evaluated; use
  `get_type_hints`/`annotationlib.get_annotations` rather than touching
  `__annotations__` directly.
- Dynamic creation: `type(name, bases, ns)` builds a class at runtime;
  `types.new_class` does it with a metaclass and kwargs.

---

## 8. Discipline — when *not* to metaprogram

Metaprogramming trades obviousness for power; overuse makes code unreadable and
breaks tooling/type-checkers. The hierarchy of restraint:

1. A plain function or class? Use it.
2. Need per-subclass behavior? `__init_subclass__` + descriptors.
3. Need typed attribute behavior? a descriptor (or `dataclass`/`pydantic`).
4. Genuinely need to control class creation/instantiation for a whole hierarchy,
   or rewrite the namespace? **then** a metaclass.

If a static type checker can't follow it, your teammates probably can't either.
Prefer the dunder/`__init_subclass__`/descriptor path that tools understand over
a clever metaclass.

---

## Sources

- [Python Data Model — Language Reference](https://docs.python.org/3/reference/datamodel.html) (`__set_name__`, `__init_subclass__`, descriptors, `__slots__`, metaclass `__prepare__`/`__new__`)
- [Descriptor HowTo Guide — Python docs](https://docs.python.org/3/howto/descriptor.html)
- [Python Metaclasses — Real Python](https://realpython.com/python-metaclasses/) and [Metaprogramming learning path](https://realpython.com/learning-paths/metaprogramming-in-python/)
- [typing.Protocol — Python docs](https://docs.python.org/3/library/typing.html#typing.Protocol) / [abc — Abstract Base Classes](https://docs.python.org/3/library/abc.html)
- [PEP 487 — Simpler customization of class creation](https://peps.python.org/pep-0487/) (`__init_subclass__`, `__set_name__`)
