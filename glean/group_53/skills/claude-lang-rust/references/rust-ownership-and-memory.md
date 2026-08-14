<!-- Provenance: reference under the `lang-rust` skill (lang family). Created 2026-06-16 via /dr deep-research. Synthesized from primary Rust docs; web sources treated as data. -->

# Rust Ownership & Memory Model

`verified-as-of: 2026-06-16` for version-volatile claims (NLL/Polonius status, edition behavior). Core
ownership semantics are stable since Rust 1.0.

The reference for ownership, borrowing, lifetimes, the borrow checker, smart pointers, and interior
mutability — the part of Rust that is genuinely new to most programmers and the source of most
early friction.

## Contents
- Ownership rules, scope-based drop, RAII
- Move semantics; Copy vs Clone
- Borrowing & references; aliasing-XOR-mutability
- The borrow checker & lifetimes (elision, `'static`, NLL, Polonius)
- Common borrow-checker fights and idiomatic fixes
- Smart pointers (Box/Rc/Arc/Cell/RefCell/Weak) & interior mutability
- Deref coercion
- Memory safety without a GC; leaks are safe
- Anti-patterns

---

## 1. Ownership rules, scope-based drop, RAII

Rust manages memory through a **compiler-checked ownership** discipline — a third path distinct from
both garbage collection and manual `malloc`/`free`. Rule violations are compile errors, and the
checks impose **no runtime cost**.[^book-own] The three rules: each value has an **owner**; there is
exactly **one owner at a time**; when the owner **goes out of scope, the value is dropped**.[^book-own]

"Dropped" means Rust calls the value's destructor (`Drop`) at the end of the owner's scope — the same
pattern C++ calls **RAII (Resource Acquisition Is Initialization)**.[^book-own] The Reference
formalizes it: each variable/temporary has a **drop scope**; values drop in **reverse order of
declaration** (but struct/tuple/enum *fields* drop in declaration order), and a partially-initialized
value drops only its initialized fields.[^ref-destructors] This gives **deterministic destruction** at
compile-time-known points — unlike GC finalization.

## 2. Move semantics; Copy vs Clone

**Move is the default** for non-`Copy` types: assigning or passing a value transfers ownership and
the source binding becomes invalid (using it after is a compile error).[^book-own] For heap-managing
types (`String`, `Vec`, `Box`) only the stack-resident pointer/len/cap moves — the heap buffer is not
duplicated — which is exactly what prevents C/C++ double-free: after `let s2 = s1;`, `s1` is invalid
so only `s2`'s drop frees the buffer.[^book-own] The Rustonomicon frames it sharply: Rust has no
move/copy/assignment constructors because **"move semantics are the only semantics"** — at most
`x = y` memcpys the bits.[^nomicon-ctor]

**`Copy` vs `Clone`:**[^std-copy][^std-clone][^comprehensive]
- `Copy` = implicit, always a trivial **bitwise** copy; not overloadable; the original stays valid
  after assignment. `Copy` can only be implemented for types that do **not** implement `Drop` and
  whose fields are all `Copy` (e.g. integers, floats, `bool`, `char`, and tuples/arrays of `Copy`).
- `Clone` is **explicit** (`x.clone()`), may run arbitrary (possibly expensive) code. `String::clone`
  deep-copies the heap buffer — which is why `String` is `Clone` but **not** `Copy` (a bitwise copy
  would alias the buffer → double free).
- `Clone` is a **supertrait of `Copy`**: every `Copy` type must also implement `Clone`.
- Shared references `&T` are `Copy`; **mutable references `&mut T` are not** (copying one would
  violate exclusivity). For `Rc`/`Arc`, `.clone()` increments the refcount rather than deep-copying.

**Partial moves:** you can move individual fields out of a composite, leaving the rest usable; the
drop machinery accounts for partially-initialized values.[^ref-destructors] (Confidence: the drop
mechanics are *fact*; the named "partial move" pattern has one direct authoritative anchor.)

## 3. Borrowing & references; aliasing-XOR-mutability

A **reference** (`&T`) accesses a value without owning it; the compiler guarantees a reference is
**always valid** (never dangling) by ensuring the referent outlives the reference. Creating one is
**borrowing**.[^book-refs] The core rule: **at any time you may have either one mutable reference
(`&mut T`) or any number of immutable references (`&T`), and references must always be
valid**.[^book-refs] A `&mut` cannot coexist with **any** other reference to the same value. The
Rust Project calls this **"mutation xor sharing"** (aliasing-XOR-mutability): shared-immutable access
*or* exclusive-mutable access, never both.[^roadmap-bck] This forbids **data races at compile
time**[^book-refs] and is also an **optimization enabler** — because `&mut` can't be aliased, the
compiler may assume two `&mut` don't overlap.[^effective-borrows]

## 4. The borrow checker & lifetimes

A **lifetime** is a *region of code (a region of validity) for which a reference must be valid* —
**not** a duration and **not** how long the data lives.[^nomicon-lifetimes][^rbe-lifetime] Rust By
Example draws the distinction: a borrow's **scope** is where the reference is *used*, and it must end
before the referent is destroyed — "lifetimes and scopes … are not the same."[^rbe-lifetime]

**Lifetime elision** — three compiler rules so you don't annotate common cases (these are *compiler*
cases, not rules for programmers):[^nomicon-elision][^book-lifetimes]
1. Each elided **input** lifetime becomes a **distinct** parameter.
2. If there is **exactly one** input lifetime, it is assigned to **all** elided **output** lifetimes.
3. If there are multiple inputs but one is `&self`/`&mut self`, `self`'s lifetime is assigned to all
   elided outputs. Otherwise eliding an output lifetime is an **error** (`E0106`), fixed with explicit
   `'a` (e.g. `fn longest<'a>(x: &'a str, y: &'a str) -> &'a str`).

**`'static` has two distinct meanings** (a precision point): `&'static T` is a reference valid for the
**entire program** (string literals are `'static`); the **bound** `T: 'static` means the type contains
**no short-lived borrows** (it is self-contained / owns its data) — it does **not** mean "lives
forever."[^book-lifetimes] A compiler suggestion to add `'static` usually signals an attempted
dangling reference; fix the lifetime mismatch, don't slap on `'static`.[^book-lifetimes]

**The borrow checker is a compile-time static analysis** that runs on the compiler's **MIR** over the
control-flow graph; it is **sound but conservative** — it never accepts unsafe code but does reject
some safe code.[^roadmap-bck][^towardsdev]

**NLL (Non-Lexical Lifetimes)** — a borrow is alive from creation to its **last use**, not to the end
of its lexical block (CFG/MIR-based, RFC 2094).[^nomicon-lifetimes][^nll-blog] NLL became enabled
**by default for all editions in Rust 1.63 (2022-08-05)**; the old AST borrow checker was removed in
late 2019.[^nll-blog] (`verified-as-of: 2026-06-16`.)

**Polonius (volatile):** the next-gen borrow checker. Its motivating case is "Problem Case #3":
returning a reference from a conditional branch.[^nll-blog] As of 2026-06 a **"Polonius alpha"**
(location-sensitive, accepts conditional borrows and lending iterators) is **pursuing stabilization in
2026**, available on nightly via `-Zpolonius` but "not ready for widespread use"; full flow-sensitivity
remains future work.[^polonius] (TENTATIVE, dated.)

## 5. Common borrow-checker fights and idiomatic fixes

Recurring fights and the expert resolutions (the consistent message: a rejection often signals
*restructure your data flow*, not reach for `.clone()`/`Rc<RefCell>`/`unsafe`):[^nomicon-split][^fighting]
- **Disjoint struct fields.** Borrowck *understands* that distinct fields can be borrowed at once —
  borrow individual fields (`&self.players` + `&mut self.scores`), not the whole `&self` +
  `&mut self`.[^nomicon-split][^fighting]
- **Disjoint slice/array elements.** Borrowck does **not** understand element disjointness —
  `&cache[i]` + `&mut cache[j]` is rejected even with `assert!(i != j)`. Idiomatic fix:
  **`slice::split_at_mut`** (returns two non-overlapping `&mut` halves; uses `unsafe` internally but
  the API is safe).[^nomicon-split][^fighting]
- **Separate read and write phases.** Collect reads first, then mutate; or scope the read in a `{ }`
  block so the reference drops early.[^fighting]
- **Use indices instead of references.** Store `usize` offsets instead of references for graph-shaped
  / self-referential data. **Tradeoff (contested):** Effective Rust warns indices become "pseudo-
  pointers" the borrow checker no longer protects, and offers `Rc`+`RefCell` as the alternative — the
  two camps genuinely disagree.[^fighting][^roadmap-bck]
- **Purpose-built APIs.** `HashMap::entry`, `Vec::retain`, `drain`, `split_at_mut` exist precisely to
  express these patterns safely.[^fighting]

## 6. Smart pointers & interior mutability

| Type | Owners | Thread-safe | Mutability | Use when |
|---|---|---|---|---|
| `Box<T>` | 1 | if `T: Send` | via `&mut` | heap alloc, trait objects, recursive types |
| `Rc<T>` | N | **No** | none (wrap in `Cell`/`RefCell`) | shared ownership, single-thread |
| `Arc<T>` | N | **Yes** | none (wrap in `Mutex`/`RwLock`) | shared ownership across threads |
| `Cell<T>` | — | No | `.get()`/`.set()` (never panics) | interior mutability for `Copy` types |
| `RefCell<T>` | 1 | No | `.borrow()`/`.borrow_mut()` (**panics** on violation) | interior mutability, runtime borrow check |

- **`Box<T>`** — the simplest heap allocation; single ownership; drops its contents on scope exit; for
  `T: Sized` it's a single pointer (ABI-compatible with C `T*`). Used for recursive types, trait
  objects, moving large data.[^std-box]
- **`Rc<T>`** — single-threaded reference counting; multiple owners; **immutable** access only; **not
  `Send`/`Sync`**.[^std-rc] **`Arc<T>`** — **atomic** refcount, the thread-safe analogue; the atomics
  make it slightly costlier than `Rc`, which is why both exist.[^smartptr-table]
- **Interior mutability (`Cell`/`RefCell`).** The std `cell` module lets you mutate **through a shared
  reference `&T`** ("interior mutability"), versus the normal **"inherited mutability"** (mutate only
  via a unique `&mut T`).[^std-cell] `RefCell` does **dynamic borrowing**: `borrow()`/`borrow_mut()`
  enforce the same "many-shared-XOR-one-exclusive" rule **at runtime**, **panicking** on
  violation.[^book-refcell][^std-cell] `Cell` never panics (swaps whole values) but works only with
  `Copy` types or via `replace`/`swap`/`take`. Neither is `Sync` — the multithreaded equivalents are
  `Mutex`/`RwLock`.
- **The `Rc<RefCell<T>>` pattern** — the canonical way to get **shared ownership + mutability** in
  single-threaded code (`Rc` gives many owners but immutable access; nesting `RefCell` restores
  mutability). `Arc<Mutex<T>>` is the multithreaded analogue.[^book-refcell]
- **`Weak<T>` & reference cycles** — refcounting **cannot free cycles** (A→B→A leaks). `Weak` is a
  **non-owning** handle that does not increment the **strong** count; `upgrade()` returns
  `Option<Rc<T>>` (`None` if dropped). Rule of thumb: **strong (`Rc`/`Arc`) for ownership edges, `Weak`
  for back-references/parent-pointers/caches**.[^smartptr-table][^std-weak]

## 7. Deref coercion

`Deref`/`DerefMut` make smart pointers behave like references. **Deref coercion** auto-converts
`&T`→`&U` (and mutable variants) when `T: Deref<Target=U>`, inserting as many `deref` calls as needed
to match a parameter type (e.g. `&String`→`&str`).[^book-deref][^std-deref] It is resolved **entirely
at compile time → zero runtime cost** and also drives method resolution. The std docs **warn**:
implement `Deref` only when pointer-like behavior is genuinely intended (`Box` deliberately exposes
almost no inherent methods to avoid collisions).[^std-deref] Abusing `Deref` to fake OO inheritance is
a named anti-pattern (see `references/rust-tooling-macros-targets.md`).

## 8. Memory safety without a GC; "fearless"

Rust encodes memory-safety rules **into the type system**, verified at compile time with **no runtime
overhead** — distinct from both manual management and GC.[^book-own] The model draws on **linear/affine
type theory** (resources used at most once), eliminating use-after-free, double-free, and dangling
pointers in safe code. The *same* ownership+type machinery prevents **data races**, making many
concurrency errors compile-time errors ("fearless concurrency" — see
`references/rust-concurrency-and-async.md`).[^book-refs]

**Leaks are SAFE (the "Leakpocalypse").** Memory leaks do **not** violate Rust's memory safety, and
Rust does not try to prevent them. Safe Rust can leak via `Rc`/`Arc` reference cycles or unbounded
collection growth.[^nomicon-unsafe] Consequently `std::mem::forget` (drops ownership *without* running
the destructor) was **made safe** (RFC 1066), with the language stating **destructors are not
guaranteed to run**.[^rfc-forget][^std-forget] Crucial expert distinction: **`unsafe` guards
undefined behavior (use-after-free, data races, dangling/aliasing violations); it does *not* guard
leaks, deadlocks, or logic races** — those are all "safe."[^nomicon-unsafe]

## 9. Anti-patterns

- **Premature `Rc<RefCell<T>>`** to dodge the borrow checker — converts compile-time errors into
  runtime `BorrowMutError` panics. It is the *correct* tool for genuine shared-mutable graphs/trees;
  the smell is *premature/spread* use, not all use. Precise consensus: it is **not** a smell to *have*
  shared mutable ownership; it *is* a smell to spread `Rc<RefCell<_>>` across many APIs — **isolate**
  it to a small part and keep most code on plain `&T`/`&mut T`.[^smell][^refcell-smell] (Contradiction
  preserved: "overused/last-resort" vs "legitimate, just isolate it.")
- **`.clone()` spam** existing only to silence a borrow error — forces fresh allocations and masks
  unclear ownership. Note `Rc::clone`/`Arc::clone` is *intentional* (bumps a refcount).
- **Severity of borrow-checker over-restriction is contested:** that it over-rejects sound programs is
  *fact*; how much it matters ranges from "productivity crater" to "goes away with experience — experts
  no longer fight it."[^challenges][^towardsdev]

## References
[^book-own]: https://doc.rust-lang.org/book/ch04-01-what-is-ownership.html — ownership rules, drop, RAII, move invalidation, no runtime cost — book.
[^ref-destructors]: https://doc.rust-lang.org/reference/destructors.html — drop scopes, reverse-order drop, field order, partial-init drop — Reference.
[^nomicon-ctor]: https://doc.rust-lang.org/nomicon/constructors.html — "move semantics are the only semantics" — Nomicon.
[^std-copy]: https://doc.rust-lang.org/std/marker/trait.Copy.html — Copy = implicit bitwise copy; excludes Drop; Clone supertrait — std.
[^std-clone]: https://doc.rust-lang.org/std/clone/trait.Clone.html — Clone explicit/arbitrary; String not Copy; Rc/Arc clone = refcount — std.
[^comprehensive]: https://google.github.io/comprehensive-rust/memory-management/copy-types.html — Copy vs Clone; &mut not Copy — book (Google).
[^book-refs]: https://doc.rust-lang.org/book/ch04-02-references-and-borrowing.html — borrowing, &/&mut rules, no dangling, data-race prevention, NLL last-use — book.
[^roadmap-bck]: https://rust-lang.github.io/rust-project-goals/2026/roadmap-borrow-checker-within.html — "mutation xor sharing," over-rejection middle ground — project goals.
[^effective-borrows]: https://effective-rust.com/borrows.html — reference rules; aliasing→optimization; Rc/RefCell/Arc/Mutex roles — Effective Rust.
[^nomicon-lifetimes]: https://doc.rust-lang.org/nomicon/lifetimes.html — lifetimes as named regions of code; last-use — Nomicon.
[^rbe-lifetime]: https://doc.rust-lang.org/rust-by-example/scope/lifetime.html — lifetime vs scope distinction — Rust By Example.
[^nomicon-elision]: https://doc.rust-lang.org/nomicon/lifetime-elision.html — the three elision rules verbatim — Nomicon.
[^book-lifetimes]: https://doc.rust-lang.org/book/ch10-03-lifetime-syntax.html — elision are compiler cases; E0106; 'static two meanings — book.
[^towardsdev]: https://towardsdev.com/rust-borrowing-the-borrow-checker-part-3-why-rust-rejects-safe-code-bfa464d2528e — sound-but-conservative; undecidability — blog.
[^nll-blog]: https://blog.rust-lang.org/2022/08/05/nll-by-default/ — NLL default all editions in 1.63 (2022-08-05); AST borrowck removed 2019 — Rust Blog.
[^polonius]: https://rust-lang.github.io/rust-project-goals/2026/polonius.html — Polonius alpha pursuing stabilization 2026; conditional borrows + lending iterators; full flow-sensitivity future — project goals.
[^nomicon-split]: https://doc.rust-lang.org/nomicon/borrow-splitting.html — borrowck understands disjoint struct fields, not slices; split_at_mut — Nomicon.
[^fighting]: https://www.atharvapandey.com/post/rust/rust-own-fighting-borrow-checker/ — split borrows, read/write phases, entry/split_at_mut/retain, "restructure don't hack" — blog.
[^std-box]: https://doc.rust-lang.org/std/boxed/index.html — Box = simplest heap alloc, single ownership, drops on scope exit — std.
[^std-rc]: https://doc.rust-lang.org/std/rc/struct.Rc.html — Rc single-threaded refcount; immutable shared access; not Send/Sync — std.
[^smartptr-table]: https://microsoft.github.io/RustTraining/rust-patterns-book/ch09-smart-pointers-and-interior-mutability.html — smart-pointer comparison; Arc atomic; Weak for cycles; Cell vs RefCell — book (Microsoft).
[^std-cell]: https://doc.rust-lang.org/std/cell/ — interior vs inherited mutability; RefCell dynamic borrowing at runtime; Cell vs RefCell; single-thread only — std.
[^book-refcell]: https://doc.rust-lang.org/book/ch15-05-interior-mutability.html — RefCell runtime borrow check, panic on violation; Rc<RefCell> pattern; compile-vs-runtime tradeoff — book.
[^std-weak]: https://doc.rust-lang.org/std/rc/struct.Weak.html — Weak non-owning, doesn't increment strong count, upgrade→Option, breaks cycles — std.
[^book-deref]: https://doc.rust-lang.org/book/ch15-02-deref.html — deref coercion &String→&str, recursive, zero runtime cost — book.
[^std-deref]: https://doc.rust-lang.org/std/ops/trait.Deref.html — smart-pointer definition; coercion warning; Box exposes no methods — std.
[^nomicon-unsafe]: https://doc.rust-lang.org/nomicon/what-unsafe-does.html — exhaustive UB causes; safe = leak/deadlock/race-condition/overflow/abort — Nomicon.
[^rfc-forget]: https://rust-lang.github.io/rfcs/1066-safe-mem-forget.html — mem::forget made safe; destructors may not run — RFC.
[^std-forget]: https://doc.rust-lang.org/std/mem/fn.forget.html — forget safe because no dtor-run guarantee; ManuallyDrop preferred — std.
[^smell]: https://users.rust-lang.org/t/is-rc-refcell-a-code-smell/27366 — NOT a smell to have shared mutable ownership; IS a smell to spread it; isolate; Servo example — forum.
[^refcell-smell]: https://users.rust-lang.org/t/why-do-all-docs-say-refcell-is-bad/37086 — RefCell as last-resort tool; runtime panic risk — forum.
[^challenges]: https://blog.rust-lang.org/2026/03/20/rust-challenges/ — experts no longer fight borrow checker; beginners do — Rust Blog.
