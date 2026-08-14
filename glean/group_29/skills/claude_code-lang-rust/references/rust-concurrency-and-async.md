<!-- Provenance: reference under the `lang-rust` skill (lang family). Created 2026-06-16 via /dr deep-research. Synthesized from primary Rust/tokio docs; web sources treated as data. -->

# Rust Concurrency & Async/Await

`verified-as-of: 2026-06-16` for version-volatile claims (scoped threads, AFIT, async-std status,
tokio specifics). Core `Send`/`Sync` and threading semantics are long-stable.

The reference for Rust's concurrency model (`Send`/`Sync`, threads, shared state, channels, rayon) and
its async/await model (futures, executors, tokio, `Pin`/`Unpin`).

## Contents
- Fearless concurrency; `Send`/`Sync`
- Threads & scoped threads
- Shared-state concurrency (Mutex/RwLock/Arc, poisoning, atomics, parking_lot)
- Message passing (channels)
- Data parallelism (rayon)
- The async/await model (lazy futures, executors, function coloring)
- tokio (tasks, schedulers, `tokio::sync`, `select!`, cancellation, `spawn_blocking`)
- async-fn-in-traits; `Pin`/`Unpin`
- Anti-patterns

---

## 1. Fearless concurrency; `Send`/`Sync`

Rust's central claim: **data races are caught at compile time**, by reusing the same ownership/borrow
machinery that governs memory safety.[^book-conc][^book-shared] A data race needs three ingredients
(two+ threads access the same location, ≥1 writes, no synchronization); the borrow checker forbids the
combination (many `&T` readers *or* one `&mut T` writer). **Caveat:** "fearless concurrency" is a
slogan — Rust prevents **data races** but does **not** prevent **deadlocks** or general logic
races.[^book-conc][^tokio-shared]

Most concurrency *primitives* live in the library (`std::sync`, `std::thread`); the language
contributes the borrow checker plus two marker traits — `Send` and `Sync`:[^book-send][^nomicon-send]
- **`Send`**: ownership of values can be safely **transferred to another thread**.
- **`Sync`**: it is safe to **share `&T` across threads**. The precise relationship: **`T: Sync` iff
  `&T: Send`**.[^book-send][^nomicon-send]
- Both are **marker/auto traits** (no methods), auto-implemented for any type composed entirely of
  `Send`/`Sync` members.[^nomicon-send]
- **`!Send`/`!Sync`:** `Rc<T>` is neither (unsynchronized refcount); `Cell`/`RefCell` are not `Sync`
  (runtime borrow tracking isn't thread-safe); raw pointers are neither; `MutexGuard` is
  `!Send`.[^book-send][^nomicon-send][^tokio-shared]
- **Manual implementation is `unsafe`** (the compiler can't verify it) — the mechanism behind
  hand-built primitives like `Arc` and `Mutex`.[^nomicon-send]

## 2. Threads & scoped threads

`std::thread::spawn` takes a (usually `move`) closure and returns a `JoinHandle<T>`; `.join()` blocks
until the thread finishes.[^std-scope] The closure is bounded by **`'static` + `Send`**, which is why
borrowed stack data historically couldn't be passed in. **Scoped threads** (`std::thread::scope`) were
**stabilized in Rust 1.63 (2022-08-11)** and let spawned threads **borrow from the local stack frame**
because the API auto-joins all un-joined threads at the end of the scope (bound `'scope`, not
`'static`).[^std-scope][^rust163] (`verified-as-of: 2026-06-16`.)

## 3. Shared-state concurrency

**`Mutex<T>` / `RwLock<T>` / `Arc`:**[^std-mutex][^std-rwlock][^arc-mutex]
- `Mutex<T>` gives exclusive access via a guard that unlocks on `Drop` (RAII). `RwLock<T>` allows many
  readers *or* one writer (use when reads ≫ writes).
- `Arc<T>` (Atomic Reference Counting) is the thread-safe analogue of `Rc`; the canonical
  shared-mutable-state pattern is **`Arc<Mutex<T>>`** (`Arc<RwLock<T>>` for read-heavy).
- **Poisoning:** if a thread **panics while holding** a `std` `Mutex`/`RwLock`, the lock becomes
  **poisoned** and later `lock()` calls return `Err(PoisonError)` (recover via `into_inner()`). An
  `RwLock` is poisoned only by a panic in a **writer**, not a reader.[^std-mutex][^std-rwlock]

**Atomics** (`std::sync::atomic`: `AtomicUsize`, `AtomicBool`, …) — lock-free primitives parameterized
by a memory **`Ordering`**: **`Relaxed`** (atomicity only), **`Acquire`** (loads), **`Release`**
(stores), **`AcqRel`**, **`SeqCst`** (a single total order).[^atomics] `Acquire`/`Release` are paired
(a `Release` store synchronizes-with a later `Acquire` load) — the lock acquire/release pattern.
(Rust's atomic memory model follows C++20's — widely stated; *qualify*.)[^atomics-book]

**`parking_lot`** — popular drop-in `Mutex`/`RwLock` alternatives: **no poisoning**, **1 byte** of
space, an inline uncontended fast path with adaptive spinning, and eventual fairness.[^parking-lot]
(*Qualify*: the throughput-vs-fairness micro-benchmark details rest largely on one strong source.)

## 4. Message passing (channels)

- **`std::sync::mpsc`** — **multi-producer, single-consumer**. `channel()` is unbounded;
  **`sync_channel(n)`** is bounded and `send` **blocks** when full (this is how std channels provide
  backpressure). Senders are `Clone`; there is one receiver.[^std-mpsc][^tokio-channels]
- **`crossbeam-channel`** — faster, **MPMC**, bounded/unbounded, with a powerful `select!`.[^crossbeam]
- **Production guidance** (repeated across sources): **prefer bounded channels** so a slow consumer
  exerts backpressure instead of letting memory grow without bound.[^crossbeam][^tokio-channels]

## 5. Data parallelism (rayon)

**Rayon** is the standard **data-parallelism** library: convert `.iter()` to **`.par_iter()`** to
parallelize an iterator chain.[^rayon] Built on a `join` primitive and a **work-stealing** scheduler;
it **statically guarantees data-race freedom** via the same `Send`/`Sync` bounds (misuse is a compile
error). Rayon targets **CPU-bound** parallelism — distinct from async (I/O concurrency).[^rayon]

## 6. The async/await model

- **Futures are lazy, poll-based state machines.** `async fn`/`async {}` makes the compiler generate a
  type implementing the **`Future`** trait; each `.await` is a state-machine transition. A Rust
  `Future` **does nothing until polled** — unlike a JS `Promise` (a common newcomer
  footgun).[^tokio-async][^async-book]
- **The `Future` trait** centers on `fn poll(self: Pin<&mut Self>, cx: &mut Context) -> Poll<T>`
  returning `Ready(T)` or `Pending`; when `Pending`, it arranges to be woken via the `Waker` in the
  `Context`.[^tokio-async]
- **Rust ships no runtime/executor** — it provides `Future` + `.await` but you pull an executor from
  the ecosystem (tokio, smol).[^tokio-async][^async-book] Rationale (Aaron Turon's "zero-cost
  futures"): the poll/readiness design lets composed futures compile to a single state-machine `enum`
  with **no per-combinator heap allocation and no forced dynamic dispatch** — a zero-cost
  abstraction.[^turon]
- **Cooperative scheduling:** tasks run until they voluntarily yield (return `Pending` at an `.await`);
  there is **no preemption**, so a future that never yields (CPU loop, blocking call) **starves the
  executor**.[^when-async]
- **Function coloring (criticism; contested — preserve both views):** per Bob Nystrom's "What Color Is
  Your Function?", `async` "infects" callers. The community is **split**: one camp says Rust async is
  "not really colored" (you can confine it to a subsystem and bridge with `block_on`); the other says
  it **is** colored (you can't pass an `async fn` where a plain `fn -> bool` is expected, e.g.
  `Iterator::filter`). Balanced reading: "Rust async **is** colored, and that's mostly **not a big
  deal** in practice."[^coloring-balanced][^coloring]

## 7. tokio — the dominant async runtime

- **Dominance is a fact:** tokio is the de-facto async runtime; the network/web ecosystem (hyper,
  tonic, axum) pulls it in transitively. It's a full ecosystem (`tokio::fs`, `net`, `io`, `time`,
  `sync`), not just a scheduler.[^corrode][^tokio-docs]
- **`#[tokio::main]`** sets up a runtime and runs `async fn main` (defaults to the multi-threaded
  scheduler).[^tokio-spawn]
- **Tasks (`tokio::spawn`)** are **green threads** — lightweight (~one allocation, on the order of tens
  of bytes), scheduled by tokio (no OS context switch), enabling thousands–millions of concurrent
  tasks. Spawned futures must be **`Send + 'static`** (the runtime may move a suspended task across
  threads at any `.await`) — the single most-cited source of async ergonomic pain and why `Arc`/`Mutex`
  proliferate.[^tokio-spawn][^tokio-tasks]
- **Multi-threaded vs current-thread scheduler:** the **multi-threaded, work-stealing** scheduler
  (`rt-multi-thread`) runs tasks across a worker pool (idle workers steal from busy ones); the
  **current-thread** scheduler (`rt`) runs all tasks on one thread (no `Send` requirement for
  `spawn_local`/`LocalSet`).[^tokio-spawn][^tokio-tasks]
- **`tokio::sync`:** async-aware `Mutex`/`RwLock`; channels — `mpsc` (bounded/unbounded), `oneshot`
  (single value), `broadcast` (fan-out), `watch` (latest-value); plus `Semaphore` and
  `Barrier`.[^tokio-spawn][^tokio-shared]
- **std `Mutex` vs tokio `Mutex` across `.await` (a top-tier expert point):** a `std`
  `MutexGuard` is **`!Send`**, so holding it across an `.await` is a **compile error** on a
  multi-threaded runtime — and even where it compiles, holding any lock across `.await` is a
  **deadlock hazard**. Guidance: **prefer `std::sync::Mutex` for short, non-`.await`-spanning critical
  sections** (cheaper); reach for **`tokio::sync::Mutex` only when you must hold the guard across
  `.await`**. Idiom: wrap the mutex and lock only inside **non-async** methods so the guard drops before
  any `.await`.[^tokio-shared]
- **`select!`:** `tokio::select!` waits on multiple branches concurrently, runs the handler of the
  **first** to complete, and **drops the rest**; it **randomly picks** among ready branches to avoid
  starvation.[^tokio-select]
- **`.await` cancellation & cancel-safety:** in async Rust, **cancellation = dropping the future**;
  when `select!` resolves one branch it drops the others, discarding in-flight state. A future is
  **cancel-safe** if "dropping and recreating it is a no-op" (no lost progress) — critical in a
  `select!` loop. `recv()` on a tokio channel is cancel-safe; an operation holding a partial read in a
  future-owned buffer is not.[^tokio-select]
- **Blocking work (`spawn_blocking`):** never run blocking I/O or heavy CPU on a runtime worker;
  `tokio::task::spawn_blocking` offloads to a dedicated blocking pool. For pure CPU parallelism, prefer
  **rayon** or plain threads over async.[^tokio-spawn][^when-async]
- **async-std / smol (volatile):** **`async-std` is discontinued** (end-of-maintenance announced
  ~2025-03-01; recorded as discontinued in RUSTSEC-2025-0052); the recommended migration target is
  **`smol`** (a small, modular executor). As of 2026-06 the ecosystem has consolidated on **tokio** for
  production, with **smol** for lightweight cases.[^corrode][^async-std-end] (`verified-as-of:
  2026-06-16`; verify the exact RUSTSEC date against the advisory DB if exactness matters.)

## 8. async-fn-in-traits; `Pin`/`Unpin`

- **AFIT/RPITIT** (async fn / return-position `impl Trait` in traits) were **stabilized in Rust 1.75
  (2023-12-21)**. **Limitations as of 2026-06:** such traits are **not dyn-compatible** (no `dyn Trait`
  dispatch), and there's no native way to add a `Send` bound to the returned future in a public trait —
  use `#[trait_variant::make]` to generate `Send`/non-`Send` variants, and the **`async-trait`** crate
  for `dyn`-dispatch needs.[^afit] (`verified-as-of: 2026-06-16`.)
- **`Pin`/`Unpin` (high level):** the state machine generated for an `async fn` is often
  **self-referential** (holds pointers into its own locals across `.await`); moving such a future would
  dangle them. Hence `poll` takes **`self: Pin<&mut Self>`** — `Pin<P>` prevents obtaining `&mut T` (and
  thus prevents moving) for types not safe to move. **`Unpin`** is an auto trait meaning "safe to move
  even when pinned"; **most types are `Unpin`**, so only self-referential generated futures (and a few
  manual types) are `!Unpin` and constrained. Application authors rarely touch `Pin` directly — they
  use `Box::pin`/`tokio::pin!`.[^pin]

## 9. Anti-patterns

- **Blocking the executor** — sync/blocking I/O, `std::thread::sleep`, or a CPU loop inside an async
  task stalls the whole worker (cooperative scheduler, no preemption). Use `spawn_blocking`, async I/O,
  `tokio::time::sleep`, or offload CPU work to rayon/threads.[^when-async][^tokio-spawn]
- **Holding a `std::sync::Mutex` (or any long-held lock) across `.await`** — compile error for std
  guards, deadlock hazard otherwise; drop the guard before awaiting or use `tokio::sync::Mutex` only
  when you must hold across await.[^tokio-shared]
- **Spawning unbounded tasks / unbounded channels under load** — no backpressure → unbounded memory
  growth / DoS. Bound with a `Semaphore`, bounded channels, or a bounded `JoinSet`.[^tokio-shared]
- **Overusing `Arc<Mutex<T>>`** (the most-cited concurrency design smell; *qualify* — sources largely
  practitioner) — one coarse mutex over a whole structure serializes all threads → contention. Prefer
  message passing, `RwLock` for read-heavy data, fine-grained/sharded locks, atomics for counters, or
  actor-style ownership.[^arc-mutex]
- **Reaching for async/tokio when threads are simpler** (*fact* across sources) — async shines for
  high-concurrency, I/O-bound, lots-of-waiting workloads; for CPU-bound work or a small number of
  long-running workers, OS threads (and `std::thread::scope`, rayon) are often faster and simpler, and
  avoid async's `'static`/`Send` constraints. Common advice: "learn synchronous Rust well; use async
  only when you actually need it."[^when-async][^corrode]

## References
[^book-conc]: https://doc.rust-lang.org/book/ch16-00-concurrency.html — fearless concurrency; ownership+types catch concurrency bugs at compile time; deadlocks not prevented — book.
[^book-shared]: https://doc.rust-lang.org/book/ch16-03-shared-state.html — shared-state, Mutex/Arc, data-race prevention — book.
[^book-send]: https://doc.rust-lang.org/book/ch16-04-extensible-concurrency-sync-and-send.html — Send/Sync, T:Sync⟺&T:Send, Rc/RefCell !Send/!Sync, auto-derivation, unsafe manual impl — book.
[^nomicon-send]: https://doc.rust-lang.org/nomicon/send-and-sync.html — authoritative Send/Sync; auto/unsafe traits; raw-ptr/UnsafeCell/Rc exceptions — Nomicon.
[^tokio-shared]: https://tokio.rs/tokio/tutorial/shared-state — std vs tokio Mutex, MutexGuard !Send, deadlock-across-await, lock-in-non-async-method idiom, bound concurrency — tokio docs.
[^std-scope]: https://doc.rust-lang.org/std/thread/fn.scope.html — std::thread::scope, ScopedJoinHandle, 'scope bound, auto-join — std.
[^rust163]: https://blog.rust-lang.org/2022/08/11/Rust-1.63.0/ — scoped threads stabilized 1.63.0 (2022-08-11) — Rust Blog.
[^std-mutex]: https://doc.rust-lang.org/std/sync/struct.Mutex.html — Mutex API, poisoning on panic, PoisonError recovery — std.
[^std-rwlock]: https://doc.rust-lang.org/std/sync/struct.RwLock.html — RwLock semantics, poisoned only on writer panic — std.
[^arc-mutex]: https://www.techbuddies.io/2026/03/20/how-to-choose-between-arc-mutex-and-rwlock-in-rust-concurrency/ — Arc<Mutex> overuse/contention antipattern & alternatives — secondary (practitioner).
[^atomics]: https://doc.rust-lang.org/std/sync/atomic/enum.Ordering.html — five orderings; acquire/release pairing — std.
[^atomics-book]: https://marabos.nl/atomics/memory-ordering.html — Mara Bos, Rust Atomics and Locks, memory ordering — book (std-libs lead).
[^parking-lot]: https://docs.rs/parking_lot — parking_lot no-poisoning, 1-byte, adaptive spin, fairness — crate docs.
[^std-mpsc]: https://doc.rust-lang.org/std/sync/mpsc/index.html — mpsc channel vs sync_channel (bounded, blocking send) — std.
[^crossbeam]: https://docs.rs/crossbeam-channel — crossbeam MPMC + select!; bounded-channel backpressure — crate docs.
[^tokio-channels]: https://tokio.rs/tokio/tutorial/channels — tokio channels; bounded backpressure guidance — tokio docs.
[^rayon]: https://github.com/rayon-rs/rayon — par_iter, join, work-stealing, statically guaranteed data-race freedom — official repo.
[^tokio-async]: https://tokio.rs/tokio/tutorial/async — Future trait, poll, lazy, state machine, executor drives poll, no language-level executor — tokio docs.
[^async-book]: https://rust-lang.github.io/async-book/02_execution/04_executor.html — executor/poll/wake mechanics — Async Book.
[^turon]: http://aturon.github.io/blog/2016/08/11/futures/ — zero-cost futures, demand-driven/readiness (poll) design, no allocation, no built-in runtime — foundational design author.
[^when-async]: https://www.rustfaq.org/en/performance-when-to-use-async-vs-threads-in-rust/ — async vs threads decision; CPU-bound needs threads/rayon; cooperative-scheduler starvation — secondary.
[^coloring]: https://journal.stuffwithstuff.com/2015/02/01/what-color-is-your-function/ — original coloring critique — essay.
[^coloring-balanced]: https://morestina.net/1686/rust-async-is-colored — "Rust async is colored, and that's not a big deal" (both sides) — blog.
[^corrode]: https://corrode.dev/blog/async/ — tokio dominance + Send+'static accidental complexity; async-std discontinued; smol; when threads win — secondary (well-regarded consultancy).
[^tokio-docs]: https://docs.rs/tokio — tokio runtime/ecosystem surface — crate docs.
[^tokio-spawn]: https://tokio.rs/tokio/tutorial/spawning — rt vs rt-multi-thread, work-stealing, tokio::sync inventory, spawn_blocking, Send+'static on spawn — tokio docs.
[^tokio-tasks]: https://docs.rs/tokio/latest/tokio/task/ — tasks as green threads, cooperative scheduling, Send/'static bounds — tokio docs.
[^tokio-select]: https://tokio.rs/tokio/tutorial/select — select! semantics, cancellation = drop, cancel-safety, recv() cancel-safe — tokio docs.
[^afit]: https://blog.rust-lang.org/2023/12/21/async-fn-rpit-in-traits/ — AFIT+RPITIT stabilized 1.75; no dyn; Send-bound problem; trait_variant/async-trait — Rust Blog.
[^pin]: https://doc.rust-lang.org/std/pin/index.html — Pin/Unpin rationale, self-referential futures, poll's Pin<&mut Self>, most types Unpin — std (+ blog.adamchalmers.com/pin-unpin).
[^async-std-end]: https://rustsec.org/advisories/RUSTSEC-2025-0052.html — async-std officially discontinued; migrate to smol — advisory (RustSec).
