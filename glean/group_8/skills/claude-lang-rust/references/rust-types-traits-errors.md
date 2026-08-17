<!-- Provenance: reference under the `lang-rust` skill (lang family). Created 2026-06-16 via /dr deep-research. Synthesized from primary Rust docs; web sources treated as data. -->

# Rust Types, Traits, Enums & Error Handling

`verified-as-of: 2026-06-16` for version-volatile claims (RPIT-2024 capture, AFIT/RPITIT, let-else,
`try_trait_v2`). Editions referenced: 2021 (default through 2024) and 2024 (current).

The reference for Rust's type-system surface: traits, generics, trait objects, enums & pattern
matching, `Option`/`Result`, the `?` operator, error idioms, closures, and iterators.

## Contents
- Traits (coherence/orphan rule, associated types, supertraits, blanket impls)
- Key standard traits
- Generics & monomorphization
- Static vs dynamic dispatch; trait objects; `impl Trait`; dyn-compatibility
- Enums & pattern matching
- `Option` / `Result`
- The `?` operator & error idioms (`thiserror` vs `anyhow`, panic-vs-Result)
- Closures & iterators
- Anti-patterns

---

## 1. Traits

A trait defines shared behavior (method signatures) — Rust's interface mechanism — and may provide
**default method bodies** implementors override or inherit.[^book-traits] Implemented via
`impl Trait for Type`.

**Coherence & the orphan rule.** *Coherence* guarantees at most one implementation of a trait for a
given type. The **orphan rule** enforces it: an `impl Trait for Type` is allowed only if **either** the
trait **or** the type is local to your crate.[^orphan][^chalk] This is why you cannot
`impl Display for Vec<T>` (both foreign) — the standard workaround is the **newtype pattern** (wrap the
foreign type in a local tuple struct).

**Associated types vs generic type parameters.** A trait with a generic param (`trait Foo<T>`) can be
implemented *multiple times* for one type with different `T`; an **associated type** (`type Item;`)
permits only **one** impl per type, so the type is uniquely determined and needs no call-site
annotation.[^rfc-assoc][^book-traits] RFC framing: generic params are "input" types (they select which
impl applies); associated types are "output" types (determined by the impl). `Iterator` uses an
associated `Item` (a type iterates over one element type); `From<T>` uses a generic param (a type is
convertible *from* many sources). Associated types also evolve more gracefully.[^rfc-assoc]

**Supertraits.** `trait A: B` requires any implementor of `A` to also implement `B`, letting `A`'s
methods call `B`'s.[^book-traits]

**Marker traits.** `Send`, `Sync`, `Sized`, `Copy` are marker/auto traits with no methods — they
assert a property the compiler reasons about.[^pretzel]

**Blanket impls.** `impl<T: Bound> Trait for T` applies to *all* types meeting the bound; being the
impl, it isn't overridable and interacts with coherence (no overlapping blanket impls).[^pretzel]
Canonical std example: `impl<T, U: From<T>> Into<U> for T` — **implementing `From` automatically gives
you `Into`**.[^pretzel] Likewise `TryInto` is blanket-provided from `TryFrom`.

## 2. Key standard traits

The derivable set:[^derive][^effective-std] `Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Hash,
Default`. Notable semantics:
- **`Debug`** (`{:?}`) — programmer-facing; derive on virtually every type (free at runtime).
  **`Display`** (`{}`) is **not** derivable — hand-written for user-facing output.[^derive]
- **`PartialEq` vs `Eq`** — `PartialEq` allows `==` but permits `x == x` to be false (e.g. `f64::NAN`);
  `Eq` is a marker asserting full reflexivity. **`PartialOrd` vs `Ord`** — partial vs total order
  (floats are `PartialOrd` only).[^derive]
- **`Hash`** — with `Eq`, required for `HashMap`/`HashSet` keys; equal values must hash equal.
- **`From`/`Into`** — fallible-free conversion; implement `From`, get `Into` free; `?` hooks into this.
  **`TryFrom`/`TryInto`** — fallible, returning `Result`.[^pretzel]
- **`Drop`** — custom destructor; its pointer also lives in a trait object's vtable.
- **`Iterator`** — one required method `next() -> Option<Self::Item>`; everything else is default.

Guidance: derive the *minimum* needed, not everything.[^derive]

## 3. Generics & monomorphization

Generics are **monomorphized**: the compiler generates a specialized copy per concrete type, so
generic code has **zero runtime overhead** — the basis of Rust's "zero-cost abstraction"
claim.[^monomorph] **Tradeoff (negation):** monomorphization causes **code bloat** — each instantiation
emits distinct machine code, enlarging binaries and lengthening compile times.[^monomorph] LTO
deduplicates identical instantiations; a common manual technique keeps a thin generic outer fn calling
a non-generic inner fn. The dispatch axis: `dyn Trait` (runtime cost) vs `impl Trait`/generics (binary
bloat).[^monomorph] Syntax: `where` clauses move bounds out of the angle brackets; the **turbofish**
`::<>` disambiguates (`"42".parse::<i32>()`, `collect::<Vec<_>>()`).

## 4. Static vs dynamic dispatch; trait objects; `impl Trait`; dyn-compatibility

**Static dispatch** is the default (concrete type known at compile time, via monomorphization).
**Dynamic dispatch** uses **trait objects** when the type is unknown until runtime.[^dispatch]

**Trait objects (`dyn Trait`)** are **fat pointers**: pointer to data + pointer to a **vtable**. The
vtable holds the method function pointers *plus* the destructor pointer and the concrete type's
size/alignment.[^dispatch] `dyn Trait` is unsized, so it lives behind a pointer (`&dyn Trait`,
`Box<dyn Trait>`).

**`impl Trait`:**[^impl-trait] argument position (`fn f(x: impl Trait)`) is sugar for an anonymous
generic param (static dispatch); return position (`-> impl Trait`, "RPIT") returns an opaque concrete
type chosen by the callee (lets you return un-nameable types like closures/iterators).
- **RPIT lifetime capture — 2024 edition change (volatile):** in Rust 2021 and earlier, RPIT captured
  a lifetime only if it appeared in the bound; **starting Rust 2024, RPIT auto-captures *all* in-scope
  type and lifetime parameters**.[^rfc3498] **Rust 1.82** added explicit `use<'x, T>` bounds to control
  capture (usable in all editions).[^impl-capture] (`verified-as-of: 2026-06-16`.)

**dyn-compatibility (formerly "object safety"):** a trait is usable as `dyn Trait` only if it is
dyn-compatible — roughly no generic methods, no methods returning `Self` by value, no `Sized` self.
Notably, **traits with `async fn` or `-> impl Trait` methods are NOT dyn-compatible** and lack dynamic
dispatch.[^afit] (The "object safety" → "dyn compatibility" rename is recent; both name the same
property.)

**Performance cost of dynamic dispatch (negation):** two pointer dereferences (vtable → method),
frequent heap allocation (`dyn` is unsized → often `Box`ed), cache inefficiency, and it **inhibits
inlining**.[^dispatch] Use static dispatch when types are known and perf-critical; reach for `dyn` for
heterogeneous collections (`Vec<Box<dyn Trait>>`).

## 5. Enums & pattern matching

Rust enums are **algebraic data types (sum types)**: variants may carry differently-typed payloads,
and the compiler prevents using a variant's payload as the wrong type.[^enums] `match` is the primary
construct and must be **exhaustive** — every variant handled (or a `_`/binding catch-all) — catching
unhandled cases at compile time.[^enums] Forms:[^enums][^book-patterns]
- **`if let`** / **`if let ... else`** — match one pattern, ignore the rest, with optional fallback.
- **`let ... else`** — binds in the surrounding scope with a *diverging* `else` (must
  `return`/`break`/`panic`); the early-exit counterpart to `if let`. **Stabilized Rust 1.65 (2022),
  RFC 3137.**[^let-else]
- **`while let`** — loop while a pattern matches.
- **Match guards** (`if` on an arm); **binding `@`** (`n @ 1..=5`); **destructuring** of structs/
  tuples/nested enums.

## 6. `Option` / `Result`

These replace null and exceptions. `Option<T>` (`Some`/`None`) encodes presence/absence — there is no
null. `Result<T, E>` (`Ok`/`Err`) encodes success/failure as a value the caller must handle.[^result]
Both are ordinary enums consumed via `match`, `if let`, combinators (`map`, `unwrap_or`,
`unwrap_or_else`, `ok_or`), or `?`.

## 7. The `?` operator & error idioms

**`?`** propagates early-exit: on `Ok(v)`/`Some(v)` it unwraps to `v`; on `Err(e)`/`None` it returns
from the enclosing function with that error/none.[^try-trait][^ref-operator] **For `Result`, `?`
applies a `From` conversion on the error**: `expr?` on an error desugars to `return Err(From::from(e))`,
so a function returning `MyError` can use `?` on any error type that implements `From<ThatError> for
MyError`.[^try-trait] **Subtlety (volatile):** under RFC 3058 (`try_trait_v2`) the `From` conversion is
*Result-specific* — `?` on `Option` does **not** do a `From` conversion (it just returns `None`); the
user-implementable `Try` trait for custom carrier types remains **unstable as of
2026-06**.[^try-v2]

**`std::error::Error` trait.** The standard error interface; `source()` exposes the underlying error,
forming an **error chain**.[^std-error] **`Box<dyn Error>`** is the idiomatic type-erased error for
quick/glue code; `Box<dyn Error + Send + Sync + 'static>` is the thread-safe form.

**`thiserror` vs `anyhow`** — the crisp distinction is **"anyhow erases types, thiserror defines
types"** and **"anyhow for applications, thiserror for libraries":**[^thiserror-anyhow]
- **`thiserror`** — a `derive` macro generating the `Error` impl (incl. `Display`, `From`/`source` via
  `#[from]`/`#[source]`) for *your* structured error enums. **For libraries**, where callers match on
  specific variants. It *defines* types.[^thiserror-anyhow]
- **`anyhow`** — an opaque `anyhow::Error` (boxed `dyn Error`) with ergonomic `.context(...)` and easy
  `?` across heterogeneous errors. **For applications/binaries** where you report/log rather than
  branch. It *erases* types.[^thiserror-anyhow]
- **Negation:** **do not use `anyhow` in a library's public API** (callers can't distinguish errors);
  `thiserror` in app code where nothing matches the variants is needless boilerplate. They're commonly
  used *together* — `thiserror` in libraries, `anyhow` at the app boundary.[^anyhow-wrong]

**panic vs Result; unwrap/expect (negation).** Rust splits **recoverable** errors (`Result`) from
**unrecoverable** ones (`panic!`).[^panic] `unwrap()` panics with a generic message; `expect("msg")`
panics with *your* message — **prefer `expect`** (documents the invariant, eases debugging).[^panic]
Acceptable in examples/prototypes/tests and where a failure signals a genuinely broken internal
invariant; **anti-pattern** in library/recoverable paths — use `?`, `match`, or combinators. Default:
return `Result` from any function that can fail.[^panic]

## 8. Closures & iterators

**Closures** are anonymous functions that capture their environment; the compiler lowers each into an
anonymous struct holding the captures, implementing the appropriate call trait — so closure calls are
**static, non-virtual, zero-overhead**.[^closures] The three traits form a supertrait hierarchy:
**`FnOnce`** (callable once; may move captures out) ⊃ **`FnMut`** (callable repeatedly; may mutate) ⊃
**`Fn`** (repeatable; only reads). A bound of `FnOnce` therefore also accepts `FnMut`/`Fn`. **`move`**
forces capture *by value* — but `move` affects *how* values are captured, **not which `Fn` trait is
implemented** (that's decided by what the body *does*; a `move` closure can still be `Fn`).[^closures]

**Iterators.** The `Iterator` trait has one required method `next(&mut self) -> Option<Self::Item>`; the
adapters are default methods.[^iter-chains] Iterators are **lazy** — adapters (`map`, `filter`, `take`,
`zip`) do nothing until a consuming adapter (`collect`, `sum`, `fold`, `for_each`) runs them — and a
**zero-cost abstraction**: chains compile to the *same machine code* as a hand-written loop, with
bounds-check elision; iterators and `for` loops are **equivalent in
performance**.[^iter-chains][^book-iter-perf] **`for x in collection`** desugars to
`IntoIterator::into_iter`; the three idiomatic forms map to ownership: **`into_iter()`** (owned `T`),
**`iter()`** (`&T`), **`iter_mut()`** (`&mut T`), so `for x in v` / `&v` / `&mut v` pick
owning/borrow behavior.[^into-iter] **Negation (readability):** long/nested iterator chains can hurt
readability — choose between loop and iterator on **clarity, not performance** (they're equal in
speed).[^iter-chains]

## 9. Anti-patterns

- **`unwrap()`/`expect()` everywhere** — treat every `unwrap()` in non-test code as a code smell;
  prefer `?`/combinators/`match`, and `expect("why")` over `unwrap()`.[^panic]
- **Overusing generics / `dyn`** mis-chosen for the axis — generics bloat binaries; `dyn` adds runtime
  indirection. Pick per §3–4.
- **Stringly-typed code** — strings where an enum/newtype belongs (primitive obsession); use enums (to
  make invalid states unrepresentable) and newtypes. See `references/rust-tooling-macros-targets.md`
  for idioms.

## References
[^book-traits]: https://doc.rust-lang.org/book/ch10-02-traits.html + ch20-02-advanced-traits.html — defining/implementing traits, default methods, associated types, supertraits — book.
[^orphan]: https://td-bn.github.io/rust/coherance-and-orphan-rules.html — orphan rule / coherence; one-local-type rule — secondary.
[^chalk]: https://rust-lang.github.io/chalk/book/clauses/coherence.html — coherence formalization — rust-lang.
[^rfc-assoc]: https://rust-lang.github.io/rfcs/0195-associated-items.html — associated types "output" vs generic params "input"; evolution tradeoff — RFC.
[^pretzel]: https://github.com/pretzelhammer/rust-blog/blob/master/posts/tour-of-rusts-standard-library-traits.md — blanket impls, marker/auto traits, From/Into & TryFrom/TryInto blanket impls — secondary (widely cited).
[^derive]: https://doc.rust-lang.org/rust-by-example/trait/derive.html — derivable traits; Debug/Clone/Copy/PartialEq/Eq/Ord/Hash/Default semantics & constraints — Rust By Example.
[^effective-std]: https://effective-rust.com/std-traits.html — std-trait semantics, PartialEq vs Eq, derive guidance — Effective Rust.
[^monomorph]: https://www.rustfaq.org/en/what-is-monomorphization-and-how-does-it-affect-performance/ — monomorphization = zero-cost but code bloat / compile-time; LTO dedup; dyn-vs-impl axis — secondary.
[^dispatch]: https://www.eventhelix.com/rust/rust-to-assembly-static-vs-dynamic-dispatch/ — static vs dynamic dispatch, vtable contents (methods+drop+size+align), two-deref cost, inlining inhibition — secondary (assembly-level).
[^impl-trait]: https://doc.rust-lang.org/reference/types/impl-trait.html — impl Trait in argument vs return position — Reference.
[^rfc3498]: https://doc.rust-lang.org/edition-guide/rust-2024/rpit-lifetime-capture.html — RPIT auto-captures all in-scope params in Rust 2024; use<..> bounds — edition guide (+ RFC 3498).
[^impl-capture]: https://blog.rust-lang.org/2024/09/05/impl-trait-capture-rules/ — impl Trait capture changes; use<> in 1.82 — Rust Blog.
[^afit]: https://blog.rust-lang.org/2023/12/21/async-fn-rpit-in-traits/ — AFIT+RPITIT stabilized 1.75; not dyn-compatible; trait_variant/async-trait — Rust Blog.
[^enums]: https://microsoft.github.io/RustTraining/csharp-book/ch06-enums-and-pattern-matching.html — enums as ADTs, exhaustiveness, destructuring — secondary (Microsoft training).
[^book-patterns]: https://doc.rust-lang.org/book/ch06-02-match.html — match, if let, exhaustiveness — book.
[^let-else]: https://rust-lang.github.io/rfcs/3137-let-else.html — RFC 3137; let-else stabilized Rust 1.65 (2022) — RFC.
[^result]: https://doc.rust-lang.org/book/ch09-02-recoverable-errors-with-result.html — Result, recoverable vs unrecoverable — book.
[^try-trait]: https://rust-lang.github.io/rfcs/1859-try-trait.html — original `?` desugaring, From conversion on error — RFC.
[^ref-operator]: https://doc.rust-lang.org/reference/expressions/operator-expr.html — `?` operator desugaring semantics — Reference.
[^try-v2]: https://rust-lang.github.io/rfcs/3058-try-trait-v2.html — redesigned `?` desugaring; From is Result-specific (not Option); Try user-trait unstable as of 2026-06 — RFC (+ tracking issue #84277).
[^std-error]: https://doc.rust-lang.org/std/error/trait.Error.html — Error trait, source(), error chain, Box<dyn Error + Send + Sync> — std.
[^thiserror-anyhow]: https://www.rustfaq.org/en/what-is-the-difference-between-thiserror-and-anyhow/ — "erases vs defines"; app-vs-library; context; #[from]/#[source] — secondary (+ docs.rs/thiserror, docs.rs/anyhow).
[^anyhow-wrong]: https://oneuptime.com/blog/post/2026-01-25-error-types-thiserror-anyhow-rust/view — don't use anyhow in public lib API; thiserror boilerplate when misused; used together — secondary.
[^panic]: https://blog.logrocket.com/panic-vs-error-rust/ — panic vs Result, expect-over-unwrap, when unwrap/expect acceptable vs anti-pattern — secondary (+ Book ch9.3).
[^closures]: https://doc.rust-lang.org/book/ch13-01-closures.html — Fn/FnMut/FnOnce hierarchy, closure→struct lowering, move-affects-capture-not-trait — book (+ Reference types/closure.html).
[^iter-chains]: https://blog.jetbrains.com/rust/2024/03/12/rust-iterators-beyond-the-basics-part-ii-key-aspects/ — iterators lazy, zero-cost, == loop performance, readability negation — secondary (JetBrains).
[^book-iter-perf]: https://doc.rust-lang.org/book/ch13-04-performance.html — iterators compile to same code as loops; bounds-check elision — book.
[^into-iter]: https://doc.rust-lang.org/std/iter/trait.IntoIterator.html — for-loop desugaring; into_iter/iter/iter_mut ownership; impls for Vec/&Vec/&mut Vec — std (+ Book ch13.2).
