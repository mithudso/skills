# Python AST, Codegen & Source Transformation

`lang-python` hub reference. **Compile/source-level** metaprogramming: parsing,
analyzing, transforming, and generating Python *source* via the `ast` module,
LibCST, and the import system.

**Scope:** the runtime side of metaprogramming (descriptors, metaclasses,
`__init_subclass__`) is `python-metaprogramming`. This spoke is the **source**
side — code that reads or rewrites code as a tree. It powers linters, formatters,
codemods, instrumentation, and DSLs.

---

## 1. The `ast` module (stdlib)

Parse source to an abstract syntax tree, walk or rewrite it, compile it back.

```python
import ast

tree = ast.parse(source)                      # str -> AST (mode="exec"/"eval")

class StripAsserts(ast.NodeTransformer):      # rewrite nodes
    def visit_Assert(self, node): return None # drop all asserts

new = ast.fix_missing_locations(StripAsserts().visit(tree))
code = compile(new, "<gen>", "exec")          # AST -> code object
exec(code, namespace)                          # or ast.unparse(new) -> source (3.9+)
```

- **`ast.NodeVisitor`** — read-only walk (`visit_<NodeType>` methods, `generic_visit`).
- **`ast.NodeTransformer`** — mutating walk; return a new node, a list, or `None`
  to replace/remove.
- **`ast.fix_missing_locations`** — copy line/col info onto synthesized nodes
  (required before `compile`, or you get errors).
- **`ast.unparse(tree)`** (3.9+) — AST back to source (lossy: comments/formatting
  gone). **`ast.literal_eval`** — safe eval of literals only (never `eval` on
  untrusted input).
- **`ast.dump(tree, indent=2)`** — inspect node structure while developing.

---

## 2. AST is lossy — when you need a CST

The stdlib AST **discards comments, whitespace, and exact formatting** — like a
JPEG, it's lossy. If you parse → modify → reprint, you can't recover the original
layout. For tools that must edit code and keep it human-readable (codemods,
linters that autofix), use a **Concrete Syntax Tree**.

---

## 3. LibCST — lossless, codemod-friendly

LibCST (Instagram) is a CST that "looks and feels like an AST" but **preserves
every formatting detail** (comments, whitespace, parens). Parses Python 3.0 →
3.14. The tool of choice for automated refactoring and lint-autofix.

```python
import libcst as cst

class RenameFoo(cst.CSTTransformer):
    def leave_Name(self, orig, updated):
        return updated.with_changes(value="bar") if updated.value == "foo" else updated

module = cst.parse_module(source)
new_source = module.visit(RenameFoo()).code     # round-trips with formatting intact
```

- **Matchers** (`libcst.matchers`) — declarative node matching instead of big
  `isinstance` trees.
- **Metadata providers** — scope/position/type info attached to nodes.
- **Codemod framework** — batch-apply transforms across a repo with a CLI
  (`python -m libcst.tool`). Use LibCST for refactors you ship back to source;
  use stdlib `ast` for analysis or generate-and-exec where formatting doesn't
  matter.

---

## 4. Import hooks — transform at import time

The import system lets you intercept and rewrite modules as they load — how
coverage tools, `pytest` assertion rewriting, and "durable execution" frameworks
work.

- Prepend a **finder** to `sys.meta_path` (`importlib.abc.MetaPathFinder`) and a
  **loader** that returns transformed code; PEP 302 / PEP 451 (`ModuleSpec`).
- Typical loader: read source → `ast.parse` → transform → `compile` → `exec` into
  the module. `importlib.util` has helpers; subclass `importlib.machinery.SourceFileLoader`
  and override `source_to_code` for the common case.
- Powerful but invasive — debuggability and startup cost suffer; prefer a build
  step or explicit decorator when you can.

---

## 5. Code generation

- **From AST:** build nodes programmatically and `ast.unparse`, or `compile` +
  `exec` for runtime code. `ast.fix_missing_locations` first. (Note: the stdlib
  `dataclasses` and `namedtuple` generate `__init__` via `exec` on *string*
  templates, not AST node construction — string templating is simpler when you
  control the output entirely.)
- **Legacy/3rd-party:** `astor`/`astunparse` (pre-3.9 unparsing), `libcst` for
  format-preserving generation.
- **String templating** (Jinja, f-strings) is simpler when you're emitting code
  from scratch and don't need to *analyze* existing code — reserve AST/CST for
  when you must understand the input.

---

## 6. Use cases

- **Linters / formatters** (flake8/ruff-style analysis, Black-style formatting).
- **Codemods / large-scale refactoring** (rename APIs across a monorepo) — LibCST.
- **Instrumentation** (inject tracing/coverage/asserts) — import hooks + ast.
- **DSLs / compilers** (rewrite a subset of Python; e.g. numba/triton-style, or
  the t-string-based DSLs enabled by PEP 750).
- **Static analysis / metrics** (count complexity, find patterns) — `ast.walk`.

---

## 7. Pitfalls

- **Version drift:** the AST grammar changes across Python versions (new node
  types, field changes). Code that walks the AST must be tested per target
  version; LibCST tracks the grammar for you.
- **Locations:** forgetting `fix_missing_locations` → `compile` errors on
  synthesized nodes.
- **Security:** never `exec`/`eval`/`pickle.load` untrusted source; AST analysis
  of untrusted code is fine, executing it is not. `ast.literal_eval` for literals.
- **Round-trip loss:** don't use stdlib `ast` to edit-and-reprint code you intend
  humans to keep reading — use LibCST.
- **Over-engineering:** import-time AST rewriting is rarely worth its
  debuggability cost — prefer decorators, a build step, or runtime
  metaprogramming (`python-metaprogramming`) unless you genuinely must transform
  arbitrary source.

---

## Sources

- [ast — Abstract Syntax Trees, Python docs](https://docs.python.org/3/library/ast.html) (NodeVisitor/NodeTransformer, unparse, literal_eval, fix_missing_locations)
- [LibCST — GitHub](https://github.com/Instagram/LibCST) and [Why LibCST](https://libcst.readthedocs.io/en/latest/why_libcst.html)
- [importlib — import system / finders & loaders](https://docs.python.org/3/library/importlib.html) (PEP 302/451)
- [Hacking the import system & rewriting the AST for durable execution — AutoKitteh](https://autokitteh.com/technical-blog/hacking-the-import-system-and-rewriting-the-ast-for-durable-execution/)
- [Parse, modify, and write back Python source using AST — sqlpey](https://sqlpey.com/python/solved-how-to-parse-modify-and-write-back-python-source-code-using-ast/)
