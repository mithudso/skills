<!-- hub-reference-banner -->
> **Reference file — part of the `blockchain` hub.** Formerly content of the standalone `solana` skill. Sibling topics in this family are reference files under the `blockchain` hub — not standalone skills. -->

# Sealevel Parallel Execution & the SVM

Sealevel and the SVM are the two halves of how Solana executes programs: **Sealevel is the parallel *scheduling* layer; the SVM is the per-program *execution* environment.**

> Consensus theory → `distributed-systems-consensus`. The Rust *language* → `lang-rust`. Builds on `account-model-and-state.md` (read that first).

## Contents

- Sealevel: the parallel runtime
- Why up-front account declaration is the enabler
- The locking rule (read-many / write-exclusive)
- The SVM (narrow vs broad meaning); sBPF bytecode
- The Banking Stage & the scheduler (evolving)
- The SVM-beyond-L1 frontier
- Disconfirming nuances

## Sealevel: the parallel runtime

**Sealevel** is Solana's parallel smart-contract runtime, introduced by Anatoly Yakovenko in 2019.[^1] Its defining mechanism: **every Solana transaction declares, up front, all the accounts it will read and write**, each flagged read-only or writable. Because the runtime knows the full read/write set *before* execution, it can schedule **non-overlapping transactions concurrently** across as many CPU cores as the validator has.[^1][^7][^11]

A transaction carries an instruction vector; each instruction names a `program_id`, an `accounts` array, and opaque `data`.

## Why up-front account declaration is the enabler

The EVM (and similar single-threaded VMs) execute one contract against global state at a time; they can parallelize only signature verification.[^1][^11] Declaring the full dependency set up front lets the runtime ask a precise question — *does this transaction's write set overlap any running transaction's read/write set?* — and run everything that doesn't conflict in parallel.[^1] This is the direct payoff of the **code/state separation** in the account model.

Yakovenko's original post framed a GPU/SIMD vision (run the same instruction over many inputs on a GPU multiprocessor).[^1] **Disconfirming/historical:** there is no Tier-1 confirmation that GPU-SIMD execution is the live production path. Production **Agave parallelizes across CPU cores** and JIT-compiles to x86-64. Treat the GPU framing as historical/aspirational, not current architecture. *(verified-as-of 2026-06-16)*

## The locking rule (read-many / write-exclusive)

Account locking is database-style:[^12][^13]

- Multiple transactions may hold a **read lock** on the same account simultaneously.
- A **write lock is exclusive** — a writable account cannot be locked if it is already read- or write-locked.
- A lock conflict returns **`AccountInUse`**; transactions touching disjoint account sets run in parallel, conflicting writes serialize.

### Parallelism is not unlimited — hot accounts serialize (disconfirming, important)

Sealevel makes contention *visible*; it does not abolish it.[^12][^16][^17][^18][^19] Any design with a single shared **writable** account — an AMM pool, a global counter, a liquidation queue — **serializes at that hotspot**. There is a per-account write-CU ceiling per block (commonly cited ~12M CU; **volatile** — see `compute-units-and-fees.md`). The real lever is **state granularity**: split state across many accounts so transactions touch disjoint sets. (Concrete example: `spl-token` deliberately does *not* write-lock the **mint** during transfers, so transfers stay parallel.) The tradeoff is more accounts + lock-table overhead + deadlock-avoidance complexity.[^19]

## The SVM (Solana Virtual Machine)

**"SVM" is overloaded** — disambiguate it explicitly:[^4][^9]

- **Narrow:** the sandboxed eBPF-derived VM — register-based (11 general-purpose registers, five bounded memory regions), compute-unit metered, with syscalls for CPI/crypto/logging.
- **Broad:** the whole decoupled transaction-processing pipeline driven by the validator's `Bank`.

Programs are written in Rust/C/Zig and compiled by LLVM to an ELF binary containing **sBPF (Solana Bytecode Format)** — a custom, non-standard eBPF variant.[^2][^4] The **rBPF** VM (forked from Quentin Monnet's project, now maintained by Anza) runs the bytecode either by **interpreter** or, in production Agave, by **JIT compilation to x86-64**.[^4][^9] Static **verification** of the bytecode runs at load; the JIT also links syscall references.[^4]

> **Naming history (FACT):** BPF → SBF was a deliberate rename (target triple `sbf-solana-solana`, GitHub #19113, 2021). The "SBF" acronym was later regretted (#30566), and **"sBPF"** is now common usage (Anza, Firedancer).[^14][^15]

### "SVM execution is parallel" is imprecise (disconfirming)

Sealevel *schedules* non-conflicting transactions in parallel, but **within a single transaction the instructions execute *sequentially***, and the per-program rBPF VM is a sequential bytecode executor.[^3][^8] The parallelism lives in the **scheduler / Banking Stage**, not inside the VM. A common conflation worth flagging.

## The Banking Stage & the scheduler (evolving)

Account locking happens in the **Banking Stage** of the leader's TPU (see `transaction-pipeline.md`).[^7] Scheduler evolution:[^6][^7][^13]

- **Original scheduler:** 4 banking threads, each greedily looping the queue — prone to lock contention.
- **Central Scheduler** (Agave v1.18, opt-in, 2024): a single scheduling thread feeds worker threads, using a **prio-graph** (lazily-evaluated DAG) with `Priority = fees / (cost + 1)` ordering; conflicting transactions process in priority order on the same thread.

> **Volatile:** scheduler internals are actively evolving (v1.18, 2024); the default may have changed. Different scheduler variants also differ on intra-slot ordering guarantees (the unified scheduler describes per-address FIFO queues, while optimistic variants allow non-deterministic retry order). Verify the current default before relying on ordering semantics. *(verified-as-of 2026-06-16)*

## The SVM-beyond-L1 frontier

The **`solana-svm` crate / SVM API** (Anza, July 2024) decoupled the SVM from `Bank`, exposing a `TransactionBatchProcessor` usable outside a full validator.[^5][^8] This enables **SVM rollups** and "SVM-beyond-Solana" execution layers — Eclipse (Ethereum SVM L2), SOON (optimistic SVM rollup), and others — as the active frontier.[^20][^21][^22] The register-based, smaller ISA is also argued to be easier to prove in zero-knowledge than the EVM.[^20]

## Disconfirming nuances (summary)

1. **Hot accounts force serialization** — parallelism scales with state granularity, not for free.
2. **The VM is sequential** — parallelism is at the scheduler layer, not inside rBPF.
3. **Intra-slot ordering can be non-deterministic** and is scheduler-variant-dependent.
4. **The GPU/SIMD framing is historical**, not the live production path (CPU cores + x86-64 JIT).
5. **"SVM" is a contested/overloaded term** — narrow (the VM) vs broad (the pipeline).

## References

[^1]: https://medium.com/solana-labs/sealevel-parallel-processing-thousands-of-smart-contracts-d814b378192 — Yakovenko (Tier 1, founder, origin post): up-front account declaration, OS-device analogy, GPU/SIMD vision.
[^2]: https://solana.com/docs/core/programs/program-execution — Solana Docs (Tier 1): sBPF via LLVM, ELF binaries, ProgramCache, JIT, loaders.
[^3]: https://solana.com/docs/core/transactions/transaction-pipeline — Solana Docs (Tier 1): instructions execute *sequentially* within a transaction; sigverify parallel.
[^4]: https://www.anza.xyz/blog/the-solana-ebpf-virtual-machine — Anza (Tier 1): rBPF lineage, interpreter vs JIT, "SVM is a misnomer," SVM ISA, verification.
[^5]: https://www.anza.xyz/blog/anzas-new-svm-api — Anza (Tier 1): `solana-svm` crate, `TransactionBatchProcessor`, decoupling from `Bank`, external use cases.
[^6]: https://www.anza.xyz/blog/introducing-the-central-scheduler-an-optional-feature-of-agave-v1-18 — Anza (Tier 1): scheduling thread, old greedy 4-thread problem, priority formula.
[^7]: https://docs.anza.xyz/architecture/ — Agave docs (Tier 1): Sealevel as the parallel engine, GreedyScheduler vs PrioGraphScheduler, Bank/SVM integration, Banking Stage.
[^8]: https://github.com/anza-xyz/agave/blob/master/svm/doc/spec.md — Agave SVM spec (Tier 1, source): SVM = the transaction-execution components, stand-alone library, load+execute stages, `Bank` not part of SVM.
[^9]: https://www.helius.dev/blog/solana-virtual-machine — Helius (Tier 2): narrow vs broad SVM; sBPF VM internals (11 registers, 5 memory regions); BPF loaders; JIT.
[^10]: https://www.helius.dev/blog/solana-vs-sui-transaction-lifecycle — Helius (Tier 2): runtime shared by TPU/TVU; account-centric dynamic conflict detection in Banking Stage.
[^11]: https://ackee.xyz/solana/book/latest/chapter2/sealevel/ — Ackee Blockchain book (Tier 2, educational): sort+schedule overview; BPF/eBPF/rBPF lineage; other chains single-threaded.
[^12]: https://dev.cube.exchange/what-is/sealevel — Cube (Tier 2): read/write lock rules, `AccountInUse`, hot-account contention limits, lock-table overhead.
[^13]: https://docs.rs/solana-unified-scheduler-logic — docs.rs (Tier 1, source docs): `SchedulingStateMachine`, per-address FIFO queues, `UsageQueue`, parallel-safe task guarantees.
[^14]: https://github.com/solana-labs/solana/issues/19113 — solana-labs (Tier 1, source): SBF = Solana Bytecode Format, `sbf-solana-solana` triple, BPF→SBF rename.
[^15]: https://github.com/solana-labs/solana/issues/30566 — solana-labs (Tier 1, source): the SBF-acronym regret, sBPF adoption, Firedancer uses sBPF.
[^16]: https://rpcfast.com/blog/solana-same-block-execution-rpc-fast — RPC Fast (Tier 2, disconfirming): optimistic parallel execution, no per-account FIFO, non-deterministic intra-slot ordering, retry windows.
[^17]: https://arxiv.org/html/2505.05358v1 — arXiv (Tier 2, academic, disconfirming): empirical transaction-conflict analysis; longest-chain-of-conflicts bounds parallelism.
[^18]: https://solana.stackexchange.com/questions/7206 — Solana Stack Exchange (Tier 2, disconfirming): per-account write-CU cap; spl-token avoids write-locking the mint to keep transfers parallel.
[^19]: https://public.bnbstatic.com/static/files/research/technical-deep-dive-parallel-execution.pdf — Binance Research (Tier 2): memory-lock model, Mutex lock table, state-granularity/parallelism tradeoff, deadlock risk.
[^20]: https://www.eclipse.xyz/articles/introducing-eclipse-mainnet-the-ethereum-svm-l2 — Eclipse (Tier 2, frontier/vendor): SVM as a parallel execution layer beyond L1; smaller ISA easier to prove in ZK.
[^21]: https://medium.com/@soon_SVM/why-and-how-to-decouple-svm-execution-layer-for-an-optimistic-rollup-8609e0fd8e01 — SOON (Tier 2, frontier/vendor): decoupled TPU (SigVerify+Banking+SVM+Entry), Anza SVM API.
[^22]: https://www.hyperlane.xyz/post/svm-expansion-the-landscape-beyond-solana — Hyperlane (Tier 2): Eclipse/SOON landscape, SVM-beyond-Solana framing.
