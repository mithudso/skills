---
name: lsp-symbol-navigation
description: >-
  Use before any symbol navigation task — finding callers, references, usages,
  or blast radius of any method, function, field, or class. Required before
  renames, security audits, dead-code checks, and refactor planning. Grep tools
  (e.g. rg, grep) silently miss callers through interfaces and supertypes and
  must not be used alone. Also governs hybrid LSP + grep cases such as
  reflection and runtime annotation processing.
source: 10gen/mms
license: Internal
mongodb:
  team: devprod-bv
  owner: shannon.monasco@mongodb.com
  internal: true
---

# LSP Symbol Navigation

LSP-based tools resolve by receiver type — only callers of *this specific symbol on this specific type* are returned, unlike grep which matches text and hits unrelated types with the same name.

## Quick decision

- **Common name** (`getName`, `getId`, `toString`, `getStatus`, `isEnabled`)? → LSP only, never grep — even in an unindexed workspace (fix indexing first).
- **External library (declaration file lives in a third-party JAR, vendored dep, or SDK outside your source tree)?** → grep with warning; LSP cannot resolve cross-root callers regardless of indexing state. If uncertain whether this applies, run the canary check — a passing canary confirms indexing is fine but does not change the cross-root limitation.
- **Fresh workspace (LSP returns 0)?** → Run a canary check (see below); if the canary also returns 0, fix indexing before falling back to grep.
- **Reflection or framework string-binding suspected?** → LSP first, targeted string-literal pre-check, then grep; report both sets separately.
- **Security audit or dead-code check?** → LSP first, then grep regardless — missing a caller has higher cost; caveat findings as potentially incomplete (see security audits and dead-code check sections for details).
- **Everything else** → LSP only.

## Tool priority

**1. MCP `find_references` tool** — check your available tools for anything named `find_references`. Use it if present; it works in parent sessions *and* subagents. No line/column required — takes `file_path` and `symbol_name`.

⚠️ **Silent-0 caveat:** like native LSP, the MCP tool returns 0 silently when the language server hasn't finished indexing or when the declaration is outside the LSP root. Run a canary check before treating 0 as "no callers" (see *Canary check* under "When grep is the right choice").

⚠️ **Ambiguous-name caveat:** if multiple symbols in `file_path` share `symbol_name` (overloads, inner classes, nested functions), the tool resolves to one without warning, producing low-but-nonzero results that look authoritative. Before calling, scan `file_path` for multiple declarations matching `symbol_name`. If more than one exists and they differ by kind (method vs. field), pass `symbol_kind` to disambiguate. If they share kind (e.g. two overloaded methods — the common Java/Kotlin case), `symbol_kind` is insufficient — use `LSP(findReferences)` with explicit `line`/`character` on the intended declaration.

```
# Pseudocode — parameter names and value syntax vary by tool
# (Python uses True/False; JSON/YAML use true/false; parameter names differ per tool)
find_references(
  file_path: "path/to/File.ext",
  symbol_name: "symbolName",
  symbol_kind: "method",   # method | field | class | function | etc.
  include_declaration: true
)

# MMS/Java example using mcp__cclsp__find_references:
mcp__cclsp__find_references(
  file_path="server/src/main/com/xgen/cloud/user/_public/model/AppUser.java",
  symbol_name="getActorId",
  symbol_kind="method",
  include_declaration=True   # Python True; note includeDeclaration=True for native LSP()
)
```

**2. `LSP(operation="findReferences")`** — parent session only; not available in subagents ([#61210](https://github.com/anthropics/claude-code/issues/61210)). Parameters are 1-indexed. `character` must land on the symbol name token — landing on an adjacent token (e.g. the return type in Java, the receiver in Go, the `func` keyword in other languages) silently resolves that token's usages instead (no error, just wrong results). If the result count is implausibly large, you've landed on the wrong token — reposition `character` to the first character of the symbol name and rerun.

```
LSP(
  operation="findReferences",
  filePath="path/to/File.ext",
  line=<1-indexed>,
  character=<1-indexed, on the first character of the symbol name>,
  includeDeclaration=True
)
```

**3. Grep tool (e.g. `rg`, `grep`)** — fallback only. ⚠️ Always surface a warning before reporting results (see below). Recall is noticeably lower than LSP in any statically-typed codebase; misses calls through interface/supertype variables.

```
# Match direct calls and method references (adjust extension for your language):
rg -n 'symbolName\(\)|::symbolName' --include="*.<ext>"
grep -rn 'symbolName()' --include="*.<ext>" path/
```

## When grep is the right choice

⚠️ **Grep is never a silent operation.** Whenever a grep tool is used, surface a warning to the user that results may include false positives and will miss callers through interfaces or supertypes. Treat grep output as supplementary to LSP, never as a replacement.

- Declaration file not reachable by the language server (external library) — LSP returns 0 silently; confirm the file exists first, then fall back to grep.
- Fresh workspace where the language server hasn't finished indexing — LSP also returns 0 silently here; run the canary check before falling back to grep.
- Quick approximate count only → grep `--count-only` (or `-c`). Warning still required.

**Canary check** — before treating any 0-reference result as authoritative, query a symbol you can visually confirm has callers in the same file. If that returns 0 too, the language server is not indexed — do not trust any results until indexing completes. To fix indexing: trigger a build, or restart the language server if your tooling supports it (e.g. `mcp__cclsp__restart_server` in MMS).

If no symbol in the target file has visually confirmable callers, pick a symbol elsewhere in the workspace with a narrow enough name that you can confirm its callers exist via a quick grep or source scan — a domain-specific class name or narrowly-named utility, not a generic name like `get` or `run`. The goal is to confirm the language server is alive and indexed, not to validate the specific file. Note: a passing canary does not validate LSP coverage for declarations *outside the LSP root* — those require grep regardless (see Quick decision).

## When both LSP and grep are required

LSP resolves statically typed call sites only. Before reaching for grep as a complement (outside the security-audit case — see bullets below), verify that reflective or dynamic references actually exist — do not assume. Run a targeted check first:

```
rg -n '"symbolName"' path/
```

If nothing relevant appears near reflection APIs or framework binding code, LSP alone is sufficient. This pre-check only catches *literal* string references — dynamic name construction (e.g. `prefix + "SymbolName"`) is invisible to grep as well; if dynamic construction is plausible, document this limitation explicitly in your findings. Only proceed with grep if evidence of reflective usage turns up or the context strongly implies it. When grep is warranted, run LSP first, then supplement. Report both result sets separately and warn the user.

- **Reflection and introspection** — symbol name appears as a string literal in a runtime lookup (e.g. Java's `getDeclaredMethod("name")`, Python's `getattr(obj, "name")`), invisible to LSP.

- **Runtime annotation processing and framework string-binding** — frameworks that resolve method or class names at runtime produce no typed call site visible to LSP. JVM examples: Spring AOP pointcuts, `@EventListener`, Guice bindings; Spring DI by string bean name (`@Qualifier("foo")`, `@Bean(name = "foo")`, `getBean("foo")`, SpEL `#{foo.bar()}`, XML `<bean ref="foo">`).

- **Code generators and annotation processors** — generated sources are only indexed by the language server when the project has been built *and* the language server is configured to include generated source paths (e.g. `target/generated-sources`, `bazel-bin`). If both are true, LSP may cover generated callers — verify before skipping grep. If build state is uncertain or generated paths are excluded, warn the user and fall back to grep for generated source paths (e.g. `*_generated.go`, `*Factory.java`, `*MembersInjector.java`).

- **Test infrastructure** — some test runners and mocking frameworks invoke symbols reflectively.

- **Security audits** — supplement LSP with grep regardless of whether reflection is confirmed. The cost of a missed call site is higher; explicitly caveat all findings as potentially incomplete.

## Never use grep for common names

`getName`, `getId`, `toString`, `isEnabled`, `getStatus`, and similar ubiquitous names — grep produces 40×–4000× more hits than real callers. LSP only. If the workspace is unindexed, run the canary check (see above) and fix indexing before retrying — never fall back to grep for common names.

Exception: security audits and dead-code checks on common names *do* warrant grep supplementation (per the Quick decision flowchart), but scope it tightly — single package or file pattern — and explicitly caveat the high false-positive rate in findings. The "never grep" rule applies to ordinary navigation and refactor planning, not to audits where missing a caller is the more expensive error. The indexing precondition still applies: if LSP is unindexed (canary fails), fix indexing first — even scoped grep on common names is unusable as primary signal when LSP is dark.

## Dead-code checks

⚠️ A 0-reference result from LSP is not sufficient to declare dead code. The symbol may be unindexed (run canary check first — see above), called reflectively, or referenced outside the LSP root (cross-repo, external SDK). Confirm with the canary check, supplement with targeted grep, and explicitly caveat any dead-code finding as potentially incomplete.

## Interface/impl split

Many language servers unify interface and implementation callers automatically — querying either side returns the same complete set. If the count is lower than expected (e.g. 0 or near-0 when callers are known to exist), verify for your specific language server:

| Language server | Unifies? | Notes |
|---|---|---|
| jdtls / cclsp (Java) | ✅ Yes | Query either side; same complete set |
| gopls (Go) | ✅ Yes | Structural typing — interface callers included |
| tsserver (TypeScript) | ✅ Yes | Generally unified; same behavior both sides |
| rust-analyzer (Rust) | ✅ Yes | Trait method callers unified |
| Pyright/Pylance (Python) | ⚠️ Partial | Protocol-typed variables may not unify — query both sides |
| clangd (C/C++) | ⚠️ Partial | Same-TU virtual calls often unified; cross-TU may miss |
