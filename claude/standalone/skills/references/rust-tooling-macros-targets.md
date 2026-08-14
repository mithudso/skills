<!-- Provenance: reference under the `lang-rust` skill (lang family). Created 2026-06-16 via /dr deep-research. Synthesized from primary Rust docs; web sources treated as data. -->

# Rust Macros, Cargo Toolchain, Compilation Targets & Idioms

`verified-as-of: 2026-06-16` for version-volatile claims (2024 edition, MSRV-aware resolver, WASI
targets, Embassy/embedded-hal status, Component Model). Macro and cargo fundamentals are long-stable.

The reference for Rust metaprogramming (declarative + procedural macros), the cargo/crates/editions/
workspace toolchain, `no_std`/embedded, the WASM target, idioms, anti-patterns, and clippy.

## Contents
- Declarative macros (`macro_rules!`)
- Procedural macros (derive/attribute/function-like; syn/quote/proc-macro2)
- cargo (manifest, lockfile, features & unification, profiles, build.rs)
- crates.io & SemVer; editions; workspaces; modules; testing
- `no_std` + embedded
- The WASM target
- Idioms; anti-patterns; clippy

---

## 1. Declarative macros (`macro_rules!`)

Declarative macros ("macros by example") work like a `match` over *Rust source tokens*: patterns match
the token structure passed in, and each arm's body replaces the invocation; they expand at compile
time.[^ref-mbe][^book-macros] **Fragment specifiers** (`$name:spec`) bind a captured fragment — `expr`,
`ident`, `ty`, `tt` (one TokenTree), `pat`, `path`, `block`, `stmt`, `item`, `literal`, `meta`, `vis`,
`lifetime`.[^ref-mbe] **AST opacity (expert gotcha):** capturing as anything other than `tt`/`ident`/
`lifetime` makes the fragment **opaque** (a single pre-parsed AST node), so reliable recursive macros
capture as `tt` and re-parse internally.[^tlborm] **Repetition** is `$( … )<sep><op>` with `*` (0+),
`+` (1+), or `?` (0/1; **forbids a separator**).[^ref-mbe]

**Hygiene — partial.** `macro_rules!` is hygienic for **local variables, labels, and `$crate`** —
nothing else. Each expansion gets a fresh syntax context, so a macro's internal `let a` can't capture
the caller's `a`. Hygiene **leaks** for generic type parameters, lifetimes, and macro-defined items;
`$crate` expands to an absolute path to the *defining* crate.[^tlborm] **Scoping quirk:** unlike
functions, `macro_rules!` macros must be defined **textually before** the call site.[^book-macros]

**When to reach for them.** Consensus: **functions first, generics/traits next, macros last** — a
feature of last resort. Use only for what the type system can't express: variadic args (`println!`,
`vec!`), DSLs, syntactic patterns, keeping disparate code in sync; **avoid hidden nonlocal control
flow** (return/loop) inside macros.[^book-macros][^effective-macros] (Macros 2.0 / the `macro`-keyword
system remains **unstable** as of 2026-06; work lands incrementally on `macro_rules!` instead.)[^macro2]

## 2. Procedural macros

The three kinds — all run **at compile time**, transforming `TokenStream`→`TokenStream`:[^book-macros][^dtolnay]
1. **Derive** — `#[derive(Name)]`; **appends** new items (cannot modify the annotated item).
2. **Attribute** — `#[route(GET, "/")]`-style; **replaces** the item; takes two inputs (attr args +
   item).
3. **Function-like** — `sql!(…)`; input need not be valid Rust.

**Crate requirement (load-bearing):** proc macros **must live in their own crate** with
`proc-macro = true`; the `proc_macro` types are usable only while a proc macro is running.[^book-macros]
**The dtolnay ecosystem:**[^dtolnay] **`proc-macro2`** (wraps the compiler API; makes proc-macro logic
unit-testable in `build.rs`/`main.rs`), **`syn`** (parses a token stream into a Rust syntax tree;
`syn::DeriveInput`, `parse_macro_input!`), **`quote`** (the `quote!` quasi-quotation macro turns syntax
back into tokens, with `#var` interpolation and `#(...)*` repetition).

**Costs (negation):** compile-time has three parts — compiling the proc-macro crate's deps
(`syn`/`quote`), *running* the macros, and **compiling the generated code** (often underappreciated;
post-expansion code can balloon).[^nethercote] Proc macros are **essentially unhygienic** and don't
recognize `$crate` — emit fully-qualified paths (`::my_crate::Thing`) and manage `Span` manually
(`call_site`/`def_site`, `quote_spanned!`). Debug with `cargo expand` and `trybuild`.[^dtolnay]

## 3. cargo

**Dependency tables:**[^cargo-deps] `[dependencies]` (compiled in, propagate downstream);
`[dev-dependencies]` (tests/examples/benches only, not propagated); `[build-dependencies]` (for
`build.rs` only). **Cargo.lock** records exact resolved versions + checksums for deterministic builds;
**current guidance (changed Aug 2023) is to commit it for *any* project** — but a library's lockfile is
excluded from the published `.crate` and does **not** constrain consumers (only `Cargo.toml`
does).[^cargo-lock] (Contradiction preserved: the old "commit for binaries, not libraries" split is
superseded but still widely cited.)

**Features** — named flags (`[features]`) gating code (`#[cfg(feature="…")]`) and/or optional deps
(`dep:name`); `default` lists auto-on features (opt out with `default-features = false`).[^cargo-features]
**The ADDITIVE rule** is the central constraint: enabling a feature must only *add* functionality —
because of **feature unification**, when multiple packages depend on the same crate Cargo builds it
**once with the union** of all requested features. Pitfalls: mutually-exclusive features are
unsupported; feature-gating *public* API breaks unification; in a workspace, features unify across
**all members**.[^cargo-features][^feature-pitfalls]

**Profiles:**[^cargo-profiles] `dev` (opt-level 0, debug on) vs `release` (`--release`; opt-level 3).
**LTO** = `"fat"`/`true` (whole-program), `"thin"` (faster), `false`/`"off"`. `panic` = `'unwind'`
(default) vs `'abort'` (smaller). **Commands:**[^cargo-commands] `build`, `check` (type-check, **skips
codegen** — the fast loop), `test` (unit + doc + integration), `run`, `clippy`, `fmt`. **Build scripts
(`build.rs`)** run **before** the package compiles (codegen, native C/C++, locating libs); they
communicate via `cargo::` directives and write to `OUT_DIR`; narrow re-runs with
`cargo::rerun-if-changed`.[^cargo-build]

## 4. crates.io & SemVer; editions; workspaces; modules; testing

**crates.io & SemVer:**[^cargo-semver] the default registry; the default version requirement is a
**caret** (`"1.2.3"` ≡ `^1.2.3` ≡ `>=1.2.3, <2.0.0`); compatibility keys off the **left-most non-zero**
component (`^0.1.2` = `>=0.1.2, <0.2.0`). **MSRV** is declared via `rust-version`; the **MSRV-aware
resolver** (opt-in since Rust 1.84, Jan 2025) prefers dep versions whose `rust-version ≤` yours and
became the **default in edition 2024** (resolver 3).[^msrv] (`verified-as-of: 2026-06-16`.)

**Editions** — Rust's mechanism for **opt-in, per-crate, backwards-incompatible language changes
without fragmenting the ecosystem**:[^editions] crates of different editions **interoperate
without friction** (a crate's edition is private to it); the **standard library is shared** and **one
`rustc`** handles all editions; migration is automated via `cargo fix --edition`. Editions: **2015**
(Rust 1.0), **2018**, **2021**, **2024**. The **2024 edition stabilized in Rust 1.85 (2025-02-20)** —
the largest to date (RPIT lifetime capture change, more `unsafe` requirements, `Future`/`IntoFuture` in
the prelude, resolver 3 default).[^rust185][^editions] (`verified-as-of: 2026-06-16`.)

**Workspaces** — one or more "members" sharing **one `Cargo.lock`** and **one `target/`** at the
root.[^workspaces][^book-workspaces] `[workspace]` keys: `members`, `exclude`, `default-members`,
`package` (inheritable metadata), `dependencies` (inheritable — members opt in with `foo.workspace =
true`), `resolver`. `[patch]`/`[profile.*]` are honored only in the root. **Inter-member deps are NOT
implicit** — a member depends on a sibling via explicit `path = "../sibling"`. **Pitfall:** features
unify across members, so a crate can compile *inside* the workspace yet fail standalone; resolver
v2/v3 mitigates build/dev/target cases but does **not** disable unification between regular
members.[^feature-pitfalls]

**Module system (concise):**[^book-modules] compilation starts at the crate root (`lib.rs`/`main.rs`),
a module named `crate`. `mod foo;` declares a module (`src/foo.rs` or `src/foo/mod.rs`; declare once,
refer by path). Everything is private by default; `pub`, `pub(crate)`, `pub(super)`, `pub(in path)`
widen visibility. Absolute paths start `crate::`; relative with `self::`/`super::`; `use` brings paths
into scope.

**Testing (concise):**[^book-testing] `#[test]` + `cargo test`; `assert!`/`assert_eq!`/`assert_ne!`;
`#[should_panic(expected="…")]`; tests may return `Result` (enables `?`); `#[ignore]`. **Unit tests** in
`#[cfg(test)] mod tests` (same file; **can test private items**). **Integration tests** in top-level
`tests/` (each file its own crate; **only the public API**). **Doc-tests** — code blocks in `///`
comments run as tests. A binary-only crate can't be integration-tested — put logic in `lib.rs`.

## 5. `no_std` + embedded

**`#![no_std]`** links only `core` (+ optionally `alloc`), not `std` (which assumes a host OS).[^no-std]
Use cases: **embedded/MCUs** (dominant), OS kernels, bootloaders, **blockchain/smart contracts**
(on-chain size is cost-metered), and **WASM**. The three-layer split:[^alloc]
- **`core`** — always available; no heap, no OS. `Option`/`Result`/primitives/iterators/slices/`str`.
  Excludes `Vec`/`String`/`Box`.
- **`alloc`** — heap subset (`Vec`, `String`, `Box`, `Rc`, `Arc`, `BTreeMap`); requires a **global
  allocator** (`#[global_allocator]`) and an allocation-error handler.
- **`std`** = `core` + `alloc` + OS integration.

**What you LOSE in no_std:** no heap *by default*; **`HashMap`/`HashSet` unavailable even with
`alloc`** (need RNG seeding — use `BTreeMap`/`BTreeSet`); no filesystem/net/threads.[^no-std-pain]
**`#![no_main]`** — bare metal has no C-runtime to call `main`; embedded crates use **`#[entry]`** from
`cortex-m-rt`. **`#[panic_handler]`** is **required** (one per binary; pull in `panic-halt`
etc.).[^embedonomicon][^nomicon-panic]

**Embedded ecosystem (high level):**
- **`embedded-hal`** — *traits* defining the HAL↔driver↔app contract for platform-agnostic drivers.
  **v1.0 shipped 2024-01-09** (stable GPIO/SPI/I2C/UART traits).[^embedded-hal] Layering: **PAC**
  (auto-generated register access via `svd2rust`) → **HAL** (ergonomic, implements embedded-hal) →
  **BSP** (board wiring).[^embedded-book]
- **Embassy** — the leading **async/await** embedded framework; each `async fn` becomes a state machine
  on a single shared stack (no per-task stacks, no heap allocator). **Maturity (2026-06):** runs on
  **stable Rust** (MSRV 1.75) with HALs for STM32/nRF/RP2040/ESP — **but the Embassy framework crates
  have explicitly unstable (pre-1.0) APIs**: stable *compiler* + stable *embedded-hal traits
  underneath*, but Embassy's own API still moves.[^embassy] **RTIC** is the lower-level **preemptive**
  alternative (tasks as interrupt handlers).[^rtic] **`probe-rs`** flashes/debugs (`cargo run` over RTT,
  with `defmt` deferred formatting).[^probe-rs]
- **Negation:** a peer-reviewed report quantifies embedded friction — **`unsafe` far heavier (48.5% of
  embedded crates vs 23.6% non-embedded)**, thin hardware coverage (concentrated on Cortex-M), painful
  C interop, doc inadequacy.[^embedded-paper] **Contradiction preserved:** that paper (fragmented,
  unsafe-heavy base) vs 2025-26 "converging fast" retrospectives (embedded-hal 1.0, AFIT, Embassy on
  stable). Reconciliation: improving rapidly off a still-incomplete base; strongest on Cortex-M.

## 6. The WASM target

**Two target families:**[^wasm-targets]
- **`wasm32-unknown-unknown`** — "bare wasm," zero environment assumptions; the **browser/JS target**
  (paired with `wasm-bindgen`). Its **std is a stub** (`println!` does nothing, `std::fs` errors,
  `thread::spawn` panics).[^wasm-targets]
- **WASI family** — adds a system interface. **The WASI rename (Rust 1.78, 2024-05-02):**
  `wasm32-wasi` (Preview 1) → **`wasm32-wasip1`**; new **`wasm32-wasip2`** (Preview 2 / Component
  Model). `wasm32-wasi` was **removed from stable in 1.84** (~Jan 2025).[^wasi-rename] Tier status
  (2026-06): `wasm32-wasip1` = **Tier 2** (since 1.78); `wasm32-wasip2` = **Tier 2 since Rust 1.82
  (2024-10-17)**.[^wasip2-tier] (`verified-as-of: 2026-06-16`; wasip2 std coverage was still filling in
  through 1.82-1.83 — verify the current surface.)

**`wasm-bindgen` + `wasm-pack`:**[^wasm-bindgen][^wasm-pack] `wasm-bindgen` generates JS↔Rust bindings
(pass strings/objects/closures, not just integers); companions **`js-sys`** (ECMAScript globals) and
**`web-sys`** (raw bindings to all Web APIs incl. DOM, each behind a Cargo feature). `wasm-pack build`
packages WASM for JS → a `pkg/` with `.wasm` + JS wrapper + `package.json` (`--target` bundler/web/
nodejs/no-modules/deno).

**Why Rust is a top WASM language:** low-level control + **no GC pauses**; **small `.wasm`** (no
runtime/GC bloat — "pay only for what you use"); plays well with JS tooling.[^wasm-book] (vs Go, which
bundles its runtime + GC → multi-MB WASM; direction is *fact*, specific byte figures are *qualify*.)
2024 State-of-Rust: ~23% target browser WASM.[^wasm-survey] **Component Model + WASI 0.2** rebased WASI
on **WIT** (interface IDL) + a component type system (`wit-bindgen`, `cargo-component`); WASI 0.2 is
production-ready for server/edge, but the **Component Model is not yet 1.0 and has no native browser
implementation as of mid-2026** (browsers use `jco transpile`) — *tentative/volatile*.[^component]

**WASM criticism (negation):** naive builds are **"too big" by design** (the wasm-bindgen guide
concedes this; rely on the CLI + `wasm-opt -Oz` to strip dead code — optimized Rust WASM is genuinely
small);[^wasm-opt] `wasm-bindgen` crate and CLI versions **must match exactly** (cryptic panics
otherwise);[^wasm-bindgen] **DOM access is not free** — every op crosses the JS boundary, so for
DOM-heavy UI optimized JS often beats routing through WASM. WASM's sweet spot is **compute, not DOM
churn**.[^wasm-dom]

## 7. Idioms

- **Newtype** — a single-field tuple struct wrapping another type, for **type safety** (`Miles(f64)` ≠
  `Kilometres(f64)`), the **orphan-rule workaround** (impl a foreign trait on a foreign type by wrapping
  it), and **encapsulation**; **zero-cost**.[^newtype] Cost: boilerplate (re-forward inner methods);
  mitigate with `derive_more`.
- **Builder** — chainable setters + a terminal `build()`; idiomatic precisely because Rust has **no
  function overloading and no default/named args**.[^builder] Mitigate boilerplate with `typed_builder`/
  `bon`/`derive_builder`.
- **Typestate** — encode run-time state into the **compile-time type** so illegal transitions fail to
  compile (a struct per state, or generic `Foo<S>` with a `PhantomData` marker); transitions consume
  `self`; **zero runtime cost**. Canonical use: embedded HAL pins (`Pin<Input>` vs
  `Pin<Output>`).[^typestate]
- **"Make illegal states unrepresentable"** — design types so invalid states can't be constructed
  (Minsky → Alexis King "Parse, Don't Validate"): **enums instead of boolean flags**; **newtype with a
  private field + validating constructor** (`fn new(...) -> Result<Self,_>`). Practical middle path:
  validate **once at the boundary**, then trust the type.[^illegal-states]
- **`Default` + `From`/`Into`** (per the Rust API Guidelines) — implement `Default` where no-arg
  construction is sensible; implement `From` (get `Into` free); name conversions `as_`/`to_`/`into_` by
  cost/ownership; eagerly derive `Copy/Clone/Eq/Debug/Default` where sensible.[^api-guidelines]

## 8. Anti-patterns

- **`unwrap()`/`expect()` everywhere** — every `unwrap()` in non-test code is a smell; prefer
  `?`/combinators/`match` and `expect("why")` over `unwrap()`.[^api-guidelines]
- **Premature `Rc<RefCell<T>>`** to dodge the borrow checker — see
  `references/rust-ownership-and-memory.md` (the smell is *premature* use, not all use).
- **`.clone()` spam** to satisfy the borrow checker — a named anti-pattern; forces allocations and
  masks ownership. (`Rc::clone`/`Arc::clone` is intentional.)[^design-patterns]
- **Stringly-typed code** — strings where an enum/newtype belongs (primitive obsession).[^design-patterns]
- **Unnecessary `unsafe`** — minimize, localize, document, hide beneath a safe interface.
- **Deref polymorphism** — abusing `Deref` to fake OO inheritance; a named anti-pattern (traits on the
  target do **not** transfer to the wrapper).[^design-patterns]

## 9. clippy

The official Rust linter (`cargo clippy`, 800+ lints; `cargo clippy --fix` auto-applies many).[^clippy]
Lint groups + default levels: **correctness** (deny), **suspicious**/**style**/**complexity**/**perf**
(warn) — together `clippy::all`; **pedantic**/**restriction**/**nursery**/**cargo** (allow, opt-in).
Configure with `#![warn(clippy::pedantic)]`/`#![deny(clippy::all)]`/`#![allow(clippy::…)]` or CLI
`-W`/`-D`/`-A`; `-D clippy::all` (errors in CI). It directly catches the anti-patterns above
(`redundant_clone`, `unnecessary_unwrap`, opt-in `unwrap_used`/`expect_used`). Never blanket-enable the
whole `restriction` group — it's a per-lint house-style tool.[^clippy]

## References
[^ref-mbe]: https://doc.rust-lang.org/reference/macros-by-example.html — macros-by-example, fragment specifiers, repetition — Reference.
[^book-macros]: https://doc.rust-lang.org/book/ch20-05-macros.html — declarative vs procedural macros; three proc-macro kinds; proc-macro crate requirement; textual-scope quirk — book.
[^tlborm]: https://veykril.github.io/tlborm/ — The Little Book of Rust Macros: AST opacity, hygiene (partial), $crate — book.
[^effective-macros]: https://effective-rust.com/macros.html — macros as last resort; avoid hidden nonlocal control flow — Effective Rust.
[^macro2]: https://rust-lang.github.io/rust-project-goals/2025h1/macro-improvements.html — macro-keyword system unstable; incremental macro_rules! work — project goals.
[^dtolnay]: https://docs.rs/syn — syn/quote/proc-macro2 ecosystem; parse_macro_input!; quote! interpolation; unit-testability; unhygienic/Span management — crate docs (dtolnay).
[^nethercote]: https://nnethercote.github.io/2025/06/26/how-much-code-does-that-proc-macro-generate.html — proc-macro compile-time, generated-code blowup — secondary (perf team).
[^cargo-deps]: https://doc.rust-lang.org/cargo/reference/specifying-dependencies.html — dependency tables; dev/build-dependencies — Cargo Book.
[^cargo-lock]: https://doc.rust-lang.org/cargo/faq.html — Cargo.lock; commit-as-default guidance; library lockfile excluded from publish — Cargo Book.
[^cargo-features]: https://doc.rust-lang.org/cargo/reference/features.html — features, additive rule, feature unification, dep:/optional — Cargo Book.
[^feature-pitfalls]: https://nickb.dev/blog/cargo-workspace-and-the-feature-unification-pitfall/ — feature unification across workspace members; standalone-vs-workspace build divergence — secondary.
[^cargo-profiles]: https://doc.rust-lang.org/cargo/reference/profiles.html — dev/release profiles, opt-level, LTO, panic — Cargo Book.
[^cargo-commands]: https://doc.rust-lang.org/cargo/commands/ — build/check/test/run; check skips codegen — Cargo Book.
[^cargo-build]: https://doc.rust-lang.org/cargo/reference/build-scripts.html — build.rs runs first; cargo:: directives; OUT_DIR; rerun-if-changed — Cargo Book.
[^cargo-semver]: https://doc.rust-lang.org/cargo/reference/specifying-dependencies.html — caret default; left-most-non-zero compatibility — Cargo Book.
[^msrv]: https://doc.rust-lang.org/edition-guide/rust-2024/cargo-resolver.html — MSRV-aware resolver (opt-in 1.84, default edition 2024 / resolver 3) — edition guide.
[^editions]: https://doc.rust-lang.org/edition-guide/editions/index.html — editions opt-in/per-crate; clean cross-edition interop; shared std; one rustc; cargo fix --edition — edition guide.
[^rust185]: https://blog.rust-lang.org/2025/02/20/Rust-1.85.0/ — 2024 edition stabilized in 1.85 (2025-02-20); resolver 3 default; RPIT capture; prelude additions — Rust Blog.
[^workspaces]: https://doc.rust-lang.org/cargo/reference/workspaces.html — workspace keys; shared lock/target; root-only patch/profile; inheritance — Cargo Book.
[^book-workspaces]: https://doc.rust-lang.org/book/ch14-03-cargo-workspaces.html — explicit path deps between members; shared target — book.
[^book-modules]: https://doc.rust-lang.org/book/ch07-02-defining-modules-to-control-scope-and-privacy.html — crate root, mod, visibility, paths, use — book.
[^book-testing]: https://doc.rust-lang.org/book/ch11-01-writing-tests.html — #[test], assert macros, should_panic, unit vs integration vs doc tests — book.
[^no-std]: https://doc.rust-lang.org/reference/names/preludes.html#the-no_std-attribute — #![no_std] links core not std — Reference (+ rust-lang no_std docs).
[^alloc]: https://doc.rust-lang.org/alloc/ — alloc crate; heap subset; requires global allocator — std/alloc docs.
[^no-std-pain]: https://docs.rust-embedded.org/book/intro/no-std.html — what no_std loses; HashMap unavailable; heapless — Embedded Book.
[^embedonomicon]: https://docs.rust-embedded.org/embedonomicon/ — #![no_main], #[entry], cortex-m-rt — Embedonomicon.
[^nomicon-panic]: https://doc.rust-lang.org/nomicon/panic-handler.html — #[panic_handler] required, one per binary — Nomicon.
[^embedded-hal]: https://blog.rust-embedded.org/embedded-hal-v1/ — embedded-hal v1.0 (2024-01-09), stable traits — Rust Embedded blog.
[^embedded-book]: https://docs.rust-embedded.org/book/ — PAC→HAL→BSP layering; cortex-m — Embedded Book.
[^embassy]: https://embassy.dev/ — async embedded framework; single shared stack; runs on stable; framework API pre-1.0/unstable — official docs.
[^rtic]: https://rtic.rs/2/book/en/ — RTIC preemptive, tasks-as-interrupts — official docs.
[^probe-rs]: https://probe.rs/ — flash/debug over RTT; defmt deferred formatting — official docs.
[^embedded-paper]: https://arxiv.org/html/2311.05063v2 — peer-reviewed "Rust for Embedded": unsafe 48.5% vs 23.6%, thin coverage, C-interop pain — academic paper.
[^wasm-targets]: https://doc.rust-lang.org/rustc/platform-support.html — wasm32-unknown-unknown stub std; WASI targets — rustc docs.
[^wasi-rename]: https://blog.rust-lang.org/2024/04/09/updates-to-rusts-wasi-targets/ — wasm32-wasi → wasm32-wasip1; wasm32-wasip2 added; wasm32-wasi removed in 1.84 — Rust Blog.
[^wasip2-tier]: https://blog.rust-lang.org/2024/11/26/wasip2-tier-2/ — wasm32-wasip2 promoted to Tier 2 in Rust 1.82 — Rust Blog.
[^wasm-bindgen]: https://wasm-bindgen.github.io/wasm-bindgen/ — JS↔Rust bindings; js-sys/web-sys; version-match requirement — official book.
[^wasm-pack]: https://rustwasm.github.io/docs/wasm-pack/ — wasm-pack build; pkg/ output; --target modes — official docs.
[^wasm-book]: https://rustwasm.github.io/docs/book/ — why Rust for WASM; small binaries, no GC; optimization — official book.
[^wasm-survey]: https://blog.rust-lang.org/2025/02/13/2024-State-Of-Rust-Survey-Results/ — ~23% target browser WASM — Rust Blog.
[^component]: https://component-model.bytecodealliance.org/ — Component Model + WASI 0.2, WIT, wit-bindgen/cargo-component; not 1.0 / no native browser as of mid-2026 — Bytecode Alliance docs.
[^wasm-opt]: https://rustwasm.github.io/docs/book/reference/code-size.html — naive builds large by design; wasm-opt -Oz strips; optimized WASM small — official book.
[^wasm-dom]: https://rustwasm.github.io/docs/book/game-of-life/implementing.html — DOM access crosses JS boundary; WASM sweet spot is compute — official book.
[^newtype]: https://rust-unofficial.github.io/patterns/patterns/behavioural/newtype.html — newtype: type safety, orphan-rule workaround, encapsulation, zero-cost — Rust Design Patterns.
[^builder]: https://rust-unofficial.github.io/patterns/patterns/creational/builder.html — builder pattern; no overloading/default args — Rust Design Patterns.
[^typestate]: https://cliffle.com/blog/rust-typestate/ — typestate: compile-time state, PhantomData, consume-self transitions, zero cost — canonical blog.
[^illegal-states]: https://corrode.dev/blog/illegal-state/ — make illegal states unrepresentable; enums/validating constructors; validate-at-boundary middle path — primary-adjacent.
[^api-guidelines]: https://rust-lang.github.io/api-guidelines/ — Default/From/Into; conversion naming; eager derives; error-handling guidance — official.
[^design-patterns]: https://rust-unofficial.github.io/patterns/ — clone-to-satisfy-borrowck, deref polymorphism, primitive-obsession anti-patterns — Rust Design Patterns.
[^clippy]: https://doc.rust-lang.org/clippy/ — clippy lint groups + default levels; configuration; --fix — official.
