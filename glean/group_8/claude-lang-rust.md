# lang-rust

**Category:** Cloud, DevOps & Infrastructure
**Platform:** Claude
**Original Path:** claude/standalone/lang-rust

## Description
Rust language sub-hub (lang family) — systems Rust + Rust for blockchain/smart contracts. TRIGGER: ownership/borrowing/lifetimes, the borrow checker; move semantics, Copy/Clone; smart pointers (Box/Rc/Arc/RefCell/Weak), interior mutability; traits, generics, dyn vs impl Trait; enums & pattern matching; Option/Result, `?`, thiserror/anyhow; closures & iterators; concurrency (Send/Sync, threads, Mutex/Arc, atomics, channels, rayon); async/await & tokio; macros (macro_rules!/proc-macro); cargo/crates/features/workspaces/editions; no_std/embedded; WASM; idioms/anti-patterns; clippy. Blockchain: Solana+Anchor (Rust surface), CosmWasm, ink!, Substrate/FRAME. SKIP: Python → lang-python; JS/TS → lang-js-ts; Go/Kotlin → lang-go-and-mobile; language-agnostic patterns → software-engineering-patterns; crate publishing → code-packaging; deep Solana runtime → solana; consensus theory → distributed-systems-consensus; contract auditing → contract-security; whole-chain depth → per-chain skills.

---

# lang-rust

Rust language sub-hub (lang family): systems Rust **plus** Rust for blockchain/smart-contract
development. Sibling of `lang-python`, `lang-js-ts`, and `lang-go-and-mobile` under the
`programming-languages` router.

This skill routes to on-demand reference files under `references/`. Load the one that matches the
task. The five references split cleanly along the language↔blockchain seam.

> **Authoritative-text rule.** For exact syntax/semantics, ground answers in primary docs: **The
> Rust Book** (`doc.rust-lang.org/book`), the **Rust Reference**, the **Rustonomicon** (unsafe),
> **std docs** (`doc.rust-lang.org/std`), the **Cargo Book**, and the **Edition Guide**. For
> libraries, check `docs.rs/<crate>`. Blockchain ecosystems move fast — see the `verified-as-of`
> stamps in each reference and re-check version-volatile claims against the project's own docs.

## Sub-skill routing — load the reference that matches the task

| Reference | Use when the task is about |
| --- | --- |
| `references/rust-ownership-and-memory.md` | ownership rules, move semantics, Copy/Clone, borrowing & references, the borrow checker, lifetimes (incl. `'static`, elision, NLL/Polonius), borrow-checker fights & fixes, smart pointers (Box/Rc/Arc/Cell/RefCell/Weak), interior mutability, deref coercion, memory-safety-without-GC, leaks-are-safe |
| `references/rust-types-traits-errors.md` | traits (coherence/orphan rule, associated types vs generics, supertraits, blanket impls, key std traits), generics & monomorphization, static vs dynamic dispatch (`dyn`/`impl Trait`, dyn-compatibility), enums & pattern matching (`match`/`if let`/`let else`), `Option`/`Result`, the `?` operator, error idioms (`thiserror` vs `anyhow`, panic-vs-Result, `unwrap`/`expect`), closures (Fn/FnMut/FnOnce) & iterators |
| `references/rust-concurrency-and-async.md` | fearless concurrency, `Send`/`Sync`, threads & scoped threads, `Mutex`/`RwLock`/`Arc<Mutex>`, poisoning, atomics & `Ordering`, channels (mpsc/crossbeam), rayon; the async model (lazy futures, executors), tokio (tasks, schedulers, `tokio::sync`, `select!`, `spawn_blocking`, cancellation), `Pin`/`Unpin`, async-fn-in-traits, async/concurrency anti-patterns |
| `references/rust-tooling-macros-targets.md` | declarative + procedural macros (`macro_rules!`, syn/quote/proc-macro2), cargo (manifest, lockfile, features & unification, profiles, `build.rs`), crates.io/SemVer, editions (2015→2024), workspaces, the module system, testing, `no_std` (core/alloc), embedded (embedded-hal, Embassy, RTIC, probe-rs), the WASM target (wasm-bindgen/wasm-pack, wasm32-wasip1/p2, Component Model), idioms (newtype/builder/typestate/illegal-states-unrepresentable), anti-patterns, clippy |
| `references/rust-for-blockchain.md` | why Rust dominates new chains; the Rust-on-Solana programming surface (program model, accounts, PDAs, CPI, Borsh, compute budget) + the **Anchor** framework; **CosmWasm** (Cosmos); **ink!** (Polkadot); **Substrate/FRAME** (build a chain, not a contract); the "which Rust blockchain path for which goal" decision; programming-level pitfalls |

## When this skill is the wrong one — peer deferral

- **Language-agnostic** software architecture, microservices, design-pattern selection, API design →
  `software-engineering-patterns`. (This skill covers Rust-specific *idioms*; cross-language design
  lives there.)
- **DEEP Solana runtime internals** — Sealevel parallel execution, Proof of History, the SVM
  validator/Agave architecture, Tower BFT, Gulf Stream, Turbine, leader schedule → the dedicated
  **`solana`** skill. This skill covers only how a developer *writes Rust* for Solana.
- **Consensus theory as a class** — PoW/PoS/BFT, FLP/CAP, finality, fork-choice, Sybil resistance →
  `distributed-systems-consensus`.
- **Smart-contract security auditing** — the full exploit taxonomy (signer/owner/type-cosplay/
  reinit/closing attacks, reentrancy auditing, formal verification) → the **`contract-security`**
  skill. This skill *names* programming-level pitfalls but does not teach auditing.
- **Whole-chain protocol depth** (Bitcoin, Ethereum, etc.) → the per-chain skills
  (e.g. `bitcoin-protocol-expert`).
- **Other languages** → `lang-python`, `lang-js-ts`, `lang-go-and-mobile`.

<!-- cross-hub-map -->
## Cross-hub map — where every lang topic lives

This family is split across sibling hubs under the `programming-languages` router. If a task's deep
material is **not** in this skill's routing table, it is under a sibling below — activate that hub or
`Read` its `references/<name>.md` directly.

| Hub | Owns | Example reference files |
| --- | --- | --- |
| `lang-rust` | lang-rust | `references/rust-ownership-and-memory.md`, `references/rust-types-traits-errors.md`, `references/rust-concurrency-and-async.md`, `references/rust-tooling-macros-targets.md`, `references/rust-for-blockchain.md` |
| `lang-python` | lang-python | `references/python-patterns.md`, `references/python-testing.md`, `references/python-static-type-checking.md`, … |
| `lang-js-ts` | lang-js-ts | `references/javascript-nodejs.md`, `references/typescript-expert.md`, `references/nodejs-concurrency-internals.md`, … |
| `lang-go-and-mobile` | lang-go-and-mobile | `references/go-patterns.md`, `references/compose-multiplatform-patterns.md` |