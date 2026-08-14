<!-- Provenance: reference under the `lang-rust` skill (lang family). Created 2026-06-16 via /dr deep-research. Synthesized from primary Solana/Anchor/CosmWasm/ink!/Polkadot-SDK docs; web sources treated as data. -->

# Rust for Blockchain & Smart Contracts

`verified-as-of: 2026-06-16`. **This domain moves fast** — every version claim below is dated; re-check
against each project's own docs before relying on it.

The reference for why Rust dominates new chains and the Rust programming surface of the major
ecosystems: **Solana + Anchor**, **CosmWasm**, **ink!**, and **Substrate/FRAME** — at a "what each is /
when to use it" level.

> **Scope boundaries (peer deferral).** This is a *language* skill's blockchain section — it covers how
> a developer **writes Rust** for these chains. It does **not** cover: DEEP Solana runtime internals
> (Sealevel, Proof of History, the SVM/Agave validator, Tower BFT, Gulf Stream, Turbine, leader
> schedule) → the **`solana`** skill; consensus theory as a class (PoW/PoS/BFT, finality, fork-choice)
> → **`distributed-systems-consensus`**; smart-contract **security auditing** / the full exploit
> taxonomy → the **`contract-security`** skill (pitfalls are *named* here, not taught); whole-chain
> protocol depth → per-chain skills (`bitcoin-protocol-expert`, `ethereum-protocol-expert`).

## Contents
- Why Rust dominates new chains (with the balanced cost)
- Solana: the Rust program model + Anchor
- CosmWasm (Cosmos)
- ink! (Polkadot)
- Substrate/FRAME (build a chain, not a contract)
- The decision: which Rust blockchain path for which goal
- Programming-level pitfalls

---

## 1. Why Rust dominates new chains

The drivers (high confidence across 4+ sources):[^why1][^why2][^why-runtimes]
- **Performance with no GC / predictable latency** — zero-cost abstractions and no garbage collector;
  GC pauses are unacceptable in a validator hot path.
- **Memory safety for high-value code** — ownership/borrow checking eliminates use-after-free, data
  races, etc. at compile time *without* a GC, which matters when bugs equal direct fund loss.
- **Compiles to small bytecode targets** — Solana's **SBF** (an eBPF derivative) and the **WASM** used
  by Polkadot/CosmWasm/NEAR/Stellar Soroban; on-chain space is expensive, so tight bytecode matters.
- **Strong type system + `no_std`** — catches bugs at compile time and fits constrained on-chain
  runtimes.
- **Ecosystem network effect** — Rust is the primary choice for non-EVM high-performance chains
  (Solana, Polkadot, NEAR, Aptos, Stellar Soroban); Cargo/crates + a shared talent pool compound it.

**Balanced / negation view** (multiple independent 2025 sources):[^neg-tooling][^neg-state][^why2]
- **Steep learning curve shrinks the dev pool** — Rust is stricter and more complex than Solidity
  (which is closer to JS); the talent pool is smaller and costlier; Solana Stack Exchange is "less
  populated than its Ethereum counterpart."
- **Long compile times + brittle toolchain versioning** — a 6-month-old repo can require specific
  toolchain versions due to breaking changes across Rust/Anchor/runtime.
- **Historically immature tooling vs EVM** (debugging, simulation, formal verification) — *qualify:*
  the gap is narrowing fast (LiteSVM, Surfpool, Codama, Pinocchio).

## 2. Solana: the Rust program model + Anchor

### 2a. The program model at the Rust API level

- **Programs are stateless Rust crates compiled to SBF.** A program holds no state; **all state lives
  in accounts** passed into each instruction, and the program *owns* the accounts holding its
  state.[^sol-structure] Each `AccountInfo` exposes `.owner`, `.is_signer`, `.is_writable`, `.key`,
  `.lamports`, `.data`.
- **The entrypoint.** The `entrypoint!(process_instruction)` macro emits the boilerplate the runtime
  calls and deserializes the raw input into the conventional signature:[^sol-structure]
  ```rust
  pub fn process_instruction(
      program_id: &Pubkey,
      accounts: &[AccountInfo],
      instruction_data: &[u8],
  ) -> ProgramResult
  ```
  Native code typically deserializes a Borsh-derived instruction enum, then `match`-routes to per-
  instruction functions, pulling accounts in order via `next_account_info`.
- **PDAs (Program Derived Addresses).** 32-byte addresses deterministically derived from **seeds +
  program ID + a bump**, SHA-256-hashed and **forced off the Ed25519 curve** (so no private key
  exists); **max 16 seeds, 32 bytes each**.[^sol-pda] `Pubkey::find_program_address(seeds, program_id)`
  returns `(pubkey, bump)` (the canonical bump is the first valid value counting down from 255). A
  program "signs" for its PDA via **`invoke_signed`**, passing `signer_seeds` (seeds + bump); the
  runtime re-derives and, on match, treats the account as signed — only the program whose ID was used
  in derivation can sign.[^sol-pda][^sol-cpi]
- **Cross-Program Invocation (CPI).** One program invokes another mid-instruction; the caller supplies
  the `Instruction` and the `AccountInfo`s for every required account. `invoke()` for normal signers,
  **`invoke_signed()`** when a PDA must sign. CPI depth is limited (commonly cited max 4) and a CPI can
  allocate at most 10,240 bytes of new account data.[^sol-cpi] (*Qualify* the exact limits — version-
  dependent.)
- **Rent / rent-exemption (dev-facing).** Accounts must hold a minimum lamport balance proportional to
  their data size to be **rent-exempt**; devs fund this at creation via
  `Rent::get()?.minimum_balance(space)`.[^sol-structure]
- **Compute budget / units (a coding constraint).** Every transaction has a compute-unit budget
  (default 200,000 CU per instruction, raisable to 1.4M via the Compute Budget program); each CPI and
  byte copied consumes CUs, so devs favor zero-copy reads and minimal CPIs.[^sol-cpi][^sol-optimize]
  (*Version-volatile* — verify defaults against the current Agave.)
- **Serialization with Borsh** (`#[derive(BorshSerialize, BorshDeserialize)]`); zero-copy
  (bytemuck/`AccountLoader`) is the CU-saving alternative for large fixed-layout
  accounts.[^sol-structure][^sol-cpi]

### 2b. The Anchor framework

Anchor is a Rust eDSL (macro framework) + toolchain + TypeScript client generator. It is *the* dominant
Solana framework and the recommended starting point for new developers as of 2026-06; it massively
reduces native boilerplate by making account validation **declarative**.[^anchor-helius][^anchor-book]
- **Core surface:**[^anchor-book] `declare_id!("…")` (program ID); `#[program]` (module of instruction
  handlers, each public fn an instruction); handlers take **`Context<T>`** first
  (`ctx.accounts`/`ctx.program_id`/`ctx.bumps`); `#[derive(Accounts)]` (generates account parsing +
  validation); `#[account]` (state struct — sets the program as owner, prepends the **8-byte
  discriminator**, auto-(de)serializes via Borsh); `#[error_code]` (custom errors).
- **Account constraints** inside `#[account(...)]`:[^anchor-constraints] `init`, `mut`, `payer = …` +
  `space = …`, `seeds = [...]` + `bump`, `has_one = field`, `constraint = <expr>`, `init_if_needed`,
  `close`, `owner`, `token::…`. The `seeds`/`bump` constraints double as PDA-correctness checks.
- **The 8-byte discriminator = the "free" security.** Anchor assigns each account type an 8-byte
  discriminator (first 8 bytes of `SHA256("account:<Name>")`) and checks it on deserialization,
  automatically **preventing "type cosplay"/account-confusion attacks**. Anchor also adds signer and
  owner checks that native code must write by hand — which is why the canonical advice is "use Anchor
  unless you have a good reason not to."[^anchor-book][^sol-security] (*Qualify*: as of Anchor 0.31 the
  discriminator is configurable; 8-byte SHA256 is the default.)
- **IDL + TS client.** Anchor emits an **IDL** (JSON: instructions, accounts, errors) used to
  auto-generate a typed TS client.[^anchor-helius] **CPI** via Anchor uses generated `cpi` helpers +
  `CpiContext` (type-checked). **Testing:** `anchor test` (a TS Mocha/Chai harness by default); **as of
  2026-06 `bankrun` is deprecated (~Mar 2025) — use LiteSVM** (in-process SVM; `anchor init` defaults
  to a LiteSVM template).[^anchor-litesvm] **`Anchor.toml`** is the manifest (program IDs per cluster,
  provider, `[toolchain]` version pins).
- **When to use Anchor vs native vs newer options (2026-06):**[^anchor-pinocchio][^sol-optimize]
  - **Anchor** — default; best DX, fastest to ship, free security checks, IDL/client generation. Cost:
    **CU + binary-size overhead** from eager `Account<T>` Borsh deserialization (`LazyAccount` mitigates).
  - **Pinocchio** (Anza) — a **zero-dependency, `no_std` replacement for the `solana-program` crate**
    (a *lower* layer, not an Anchor competitor); zero-copy; claims large CU savings vs Anchor
    (*qualify* the exact figure). **Maturity caveat:** as of mid-2025 unaudited, not at feature parity,
    not beginner-friendly.
  - **Steel** — a modular, less-opinionated wrapper over `solana-program`; native-like performance; **no
    IDL generation** (its main gap); assumes solid Rust.
  - **Rule of thumb:** ship on **Anchor**; once stable + high-volume, optimize hot programs with
    Pinocchio/Steel/native.

### 2c. The broader Solana-Rust toolchain (date-flagged 2026-06)

- **Build:** `cargo build-sbf` is current; **`cargo build-bpf` is removed** (Agave v2.0+ errors).
- **Validator client:** **Agave** (by **Anza**, formerly Solana Labs); **CLI:** the `solana` tool suite
  (recent line ~2.2.x); **Anchor:** 0.31.0 is the relevant modern line.[^anchor-031]
- **SBF/sBPF naming:** BPF → **SBF** (Solana Bytecode Format); the VM is **sBPF**. Use "SBF/sBPF," not
  "BPF," in current writing.[^sbf]
- **`solana-program` monolith → modular crates (important churn):** historically monolithic
  `solana-program`/`solana-sdk` split into focused crates (`solana-pubkey`, `solana-account-info`,
  `solana-cpi`, …) which `solana-program` re-exports. `solana-program` reached **3.0.0 (2025-08-11)**
  and **4.0.0 (2026-02-17)** — pin versions and check Anchor compatibility before upgrading.[^sol-crates]
  (*Tentative* on the exact 4.0.0 date.)
- **IDL/client gen for non-Anchor:** **Shank** (annotate → IDL) + **Codama** (IDL → typed
  clients).[^anchor-pinocchio]

## 3. CosmWasm (Cosmos)

**What it is.** A framework for writing smart contracts in Rust that compile to **WASM**, running on
Cosmos-SDK chains via the **`x/wasm`** module (the `wasmd` reference app) — pluggable into an existing
chain without changing its logic.[^cw-entry][^cw-wasmd] Because it sits on Cosmos-SDK + **IBC**, the
same contract logic can deploy across IBC-connected chains. Deployment is two-phase: a `store code` tx
uploads bytecode → a `code_id`, then one or more *instances* are instantiated (contrast Solana, which
has no multi-instance model).[^cw-vs-sol]

**The actor model (no reentrancy by design).** Each contract is an **actor**: it processes one message
at a time, can only modify its own private state, and affects other contracts only indirectly through
messages.[^cw-actor][^cw-cto] When A wants to call B, A finishes its state changes and *returns a
message* to invoke B rather than synchronously calling in — guaranteeing "a consistent view … in
storage" and **avoiding reentrancy attacks** by construction. This is the deliberate contrast with
EVM/Solidity synchronous calls (which enabled reentrancy, e.g. the DAO); CosmWasm trades synchronous
composability for safety.[^cw-cto]

**Entry points & types.**[^cw-entry][^cw-actor] Required: `instantiate`; common set:
`instantiate`/`execute`/`query` plus optional `migrate`/`sudo`/`reply`. `instantiate`/`execute` take
`(DepsMut, Env, MessageInfo, msg)`; `query` takes `(Deps, Env, msg)` and is read-only. APIs are Rust
enums — `InstantiateMsg`, `ExecuteMsg`, `QueryMsg` — serialized as **JSON**. `Response` (returned by
state-changing entry points) bundles events/attributes, (sub)messages to re-dispatch, and an optional
data blob. **State:** `cw-storage-plus` gives typed `Item<T>` and `Map<K,V>`; `cosmwasm-std` is the
stdlib; **`cw-multi-test`** is an in-Rust multi-contract simulator; `cosmwasm-schema` generates JSON
schemas.[^cw-storage][^cw-multitest] **CW standards:** `cw20` (fungible ≈ ERC-20), `cw721` (NFT ≈
ERC-721).[^cw-plus] **Versions (date-flag):** CosmWasm **2.0** (Mar 2024, gas cut ~1000×), **2.2** (Dec
2024), **3.0** (Jun 2025, IBCv2 entry points).[^cw-versions]

**When to use:** building dApps on the Cosmos ecosystem / app-chains where **IBC** interoperability
matters, or wanting Rust+WASM contracts pluggable into a Cosmos-SDK chain.

## 4. ink! (Polkadot)

**What it is.** An **eDSL for Rust** for writing smart contracts that run on Substrate-based chains
including the contracts pallet.[^ink-sdk] **Macros:**[^ink-macros] `#[ink::contract]`,
`#[ink(storage)]`, `#[ink(constructor)]`, `#[ink(message)]`, `#[ink(event)]`, `#[ink(payable)]`. It uses
the **SCALE codec** (`parity-scale-codec`) — a compact binary format — for serialization (the key
contrast vs CosmWasm's JSON).[^ink-scale]

**THE BIG DATE-FLAG (2026-06): WASM → RISC-V/PolkaVM.** ink! is migrating from `pallet-contracts` +
WebAssembly (run in the `wasmi` interpreter) to **`pallet-revive`** + **RISC-V** (run on
**PolkaVM**).[^ink-sdk][^ink-riscv] Parity forked `pallet-contracts` into `pallet-revive` in late 2024;
contracts uploaded to it are **RISC-V (PVM) bytecode**, not WASM. **ink! v6** is the version targeting
RISC-V/PVM; `pallet-revive` is **decoupled from the language** and supports **both ink! (Rust) and
Solidity** (via Parity's `revive`/`resolc` compiler), introducing **20-byte EVM-style accounts**.
`pallet-revive` merged with an enactment date of **2026-01-27**; contracts deploy on **Asset Hub**.
`pallet-contracts` is **legacy** (no formal removal date found — *qualify*).[^ink-riscv][^revive-status]

**⚠️ CONTRADICTION — preserve both halves (date-flag 2026-06).** Simultaneously, **Parity announced in
Jan 2026 it could no longer actively maintain or develop ink!** after funding lapsed; ink! became a
**community project** (the **ink! Alliance**) funded by **Polkadot Treasury grants** that delivered the
v6/PolkaVM port.[^ink-discontinue][^ink-alliance] One ecosystem source says the Alliance itself was
under funding strain in early 2026, conflicting with the active v6-delivery referenda — treat the
Alliance's current health as **tentative/contested** (the documented Treasury referenda are the
stronger signal). Accurate statement: *Parity ceased first-party ink! maintenance; the language is
alive as a Treasury-funded community project that shipped v6.*

**ink! contracts vs runtime pallets.** "Smart contracts sit on top of parachains" — sandboxed,
deployable by anyone, constrained by what the chain allows, and "can never be as fast as a native
pallet."[^ink-sdk] (See §5.) **When to use:** deploying contracts (deployed-by-anyone, sandboxed) on
Polkadot-ecosystem chains — historically Astar, Aleph Zero, now Polkadot **Asset Hub** via
`pallet-revive`; choose ink! over Solidity-on-revive for Rust's safety + the macro ergonomics.

## 5. Substrate/FRAME (build a chain, not a contract)

**What it is.** A Rust framework (from Parity) for building **application-specific blockchains**
(app-chains / Polkadot parachains) — you build the chain's **runtime** itself, not a contract on
someone else's chain.[^polkadot-sdk-runtime] The runtime is the **state-transition function**, compiled
to a **WASM runtime stored on-chain**, which enables **forkless runtime upgrades**.

**FRAME (the pallet system).** FRAME composes a runtime from **pallets** (units of encapsulated logic —
balances, staking, governance, your custom module), defined with `#[frame_support::pallet]`; the
runtime is the composition of chosen pallets.[^polkadot-sdk-runtime]

**The key mental model — contract vs runtime/pallet** (the most important distinction in this
cluster):[^polkadot-sdk-runtime]

| Dimension | Smart contract (ink!/CosmWasm/Solana program) | Runtime module / FRAME pallet |
|---|---|---|
| Who deploys | **Anyone** (permissionless) | Privileged; upgrades via **on-chain governance** |
| Isolation | Sandboxed | Integrated into the chain core |
| Resource cost | Metered dynamically (gas/CU) | **Weighed** — fixed cost known in advance (benchmarked) |
| Privilege | Cannot modify chain rules | **Direct access to chain state** |
| Performance | Slower (VM overhead) | Faster, optimized for the chain |
| Blast radius | That contract + its users | A bug "can compromise the entire blockchain" |

Authoritative TL;DR: **"If you need to create a Blockchain, then write a runtime. If you need to create
a DApp, then write a Smart Contract."**[^polkadot-sdk-runtime] **Polkadot SDK (date-flag):** Substrate is
now folded into the **Polkadot SDK** monorepo (`paritytech/polkadot-sdk`) alongside FRAME and
Cumulus.[^polkadot-sdk-runtime][^substrate-sdk] **When to use:** when you need to **build your own
L1/parachain** (full sovereignty, custom consensus, native-fast privileged logic) rather than deploy a
contract.

## 6. The decision — which Rust blockchain path for which goal

Two orthogonal axes organize the whole space:[^polkadot-sdk-runtime][^cw-vs-sol]

**Axis A — deploy code onto a chain (contract) vs build a whole chain (runtime):**
- **Contract** (permissionless, sandboxed, gas/CU-metered, slower): **Solana programs/Anchor**,
  **CosmWasm**, **ink!**.
- **Chain/runtime** (privileged, native-fast, governance-gated): **Substrate/FRAME** (the Cosmos peer
  for building a chain is the **Cosmos SDK in Go**, not a Rust contract framework).

**Axis B — contract execution target:**
- **WASM contracts:** CosmWasm; ink! through v5.
- **RISC-V/PolkaVM (new, 2026):** ink! v6 + Solidity on `pallet-revive`.
- **SBF/eBPF:** Solana programs.
- **Native WASM runtime (not a contract VM):** Substrate/FRAME pallets.

**Ecosystem alignment (the fast routing rule):**
- **Cosmos** → **CosmWasm** for contracts; **Cosmos SDK (Go)** to build the chain.
- **Polkadot** → **ink!** (or Solidity-on-revive) for contracts on Asset Hub/parachains;
  **Substrate/FRAME** to build the parachain.
- **Solana** → **Anchor / native Rust** programs.

**The common thread:** all chose Rust for compile-time memory safety with **no GC**, `no_std`/
small-runtime builds, a **WASM (now also RISC-V)** compile target at near-native speed, and determinism
(sandboxing, overflow checks) — exactly the constraints of on-chain code (see §1).[^why-runtimes]

## 7. Programming-level pitfalls (auditing deferred to `contract-security`)

Named here, not taught — the full exploit taxonomy and audit methodology belong to the
**`contract-security`** skill:[^sol-security]
1. **Missing signer check** — verifying a pubkey without checking `is_signer` lets anyone impersonate
   an authority.
2. **Missing owner/account-ownership check** — accepting an account owned by a malicious program;
   native code must assert `account.owner == program_id`.
3. **Account confusion / "type cosplay"** — passing a different account type than expected (Anchor's
   8-byte discriminator prevents this for free).
4. **Integer overflow in on-chain math** — Rust **does not check overflow in release mode, which Solana
   uses by default**; use `checked_add/sub/mul` (or `overflow-checks = true`).
5. **Reinventing checks a framework gives free** — hand-rolling validation Anchor's constraints/types
   already enforce, introducing gaps.
6. **Exceeding the compute budget** / **large stack usage** (SBF has small 4 KB frames — eager
   deserialization can overflow it; favor zero-copy/boxing).
7. **(CosmWasm/ink!)** submessage/`reply` error-handling bugs and storage-layout/upgrade hazards — also
   deferred to `contract-security`.

## References
[^why1]: https://www.developernation.net/blog/a-comprehensive-guide-to-rust-programming-language-for-smart-contracts-development-web3/ — why Rust for smart contracts (perf, no-GC, memory safety, adoption) — blog.
[^why2]: https://www.blockchain-council.org/smart-contracts/solidity-vs-rust-for-smart-contracts-ethereum-vs-solana/ — Rust vs Solidity; non-EVM adoption; learning-curve cost — blog.
[^why-runtimes]: https://www.quillaudits.com/blog/smart-contract/wasm-smart-contracts — Rust memory safety/no-GC/no_std/WASM near-native, determinism — blog.
[^neg-tooling]: https://shubhendukumar125.medium.com/deep-dive-of-the-state-of-developer-tooling-on-solana-july-2025-7e0e9d6e0e0e — toolchain version brittleness, opaque debugging, maturing 2025 — blog (NEGATION).
[^neg-state]: https://blog.superteam.fun/p/deep-dive-of-the-state-of-developer — onboarding debt, talent pool, tooling gaps vs EVM — blog (NEGATION).
[^sol-structure]: https://solana.com/docs/programs/rust/program-structure — Solana Rust program structure: entrypoint, process_instruction, accounts, AccountInfo, Borsh, rent — official-docs.
[^sol-pda]: https://solana.com/docs/core/pda — PDA derivation, seeds/bump, off-curve, 16-seed/32-byte limits, find_program_address — official-docs.
[^sol-cpi]: https://www.quicknode.com/guides/solana-development/anchor/what-are-cpis — CPI, invoke/invoke_signed, 10,240-byte limit, depth — blog (corroborated by docs).
[^sol-optimize]: https://www.helius.dev/blog/optimizing-solana-programs — Anchor CU/overhead, native vs Anchor, eager deserialization, compute budget — blog (authoritative; NEGATION).
[^anchor-helius]: https://www.helius.dev/blog/an-introduction-to-anchor-a-beginners-guide-to-building-solana-programs — Anchor as DSL/toolchain/TS-client; IDL — blog (authoritative).
[^anchor-book]: https://www.anchor-lang.com/docs/basics/program-structure — declare_id!, #[program], #[derive(Accounts)], Context, #[account], 8-byte discriminator (SHA256 "account:<Name>") — book (official).
[^anchor-constraints]: https://www.quicknode.com/guides/solana-development/anchor/how-to-use-constraints-in-anchor — account constraints (init/mut/has_one/seeds/bump/payer/space/constraint) — blog.
[^sol-security]: https://blog.neodyme.io/posts/solana_common_pitfalls — signer/owner/type-cosplay/overflow pitfalls; "use Anchor" advice — blog (+ solana.com program-security course).
[^anchor-litesvm]: https://www.quicknode.com/guides/solana-development/tooling/litesvm — bankrun deprecated ~Mar 2025, LiteSVM default — blog.
[^anchor-pinocchio]: https://www.helius.dev/blog/pinocchio — Pinocchio zero-dep/zero-copy, CU savings, unaudited/feature-gaps; Shank/Codama — blog (authoritative).
[^anchor-031]: https://www.anchor-lang.com/docs/updates/release-notes/0-31-0 — Anchor 0.31, Solana 2.1.0 rec, agave-install — book (official).
[^sbf]: https://github.com/anza-xyz/sbpf — sBPF/SBF naming; rBPF archived 2025-01-10 — official-repo.
[^sol-crates]: https://crates.io/crates/solana-program — modular crate split; 3.0.0 (2025-08-11)/4.0.0 (2026-02-17) — official-registry.
[^cw-entry]: https://book.cosmwasm.com/basics/entry-points.html — entry points (instantiate required; #[entry_point]; Deps/DepsMut/Env/MessageInfo/Response) — book (official).
[^cw-wasmd]: https://cosmwasm.cosmos.network/wasmd — x/wasm module wiring into a Cosmos-SDK chain — official-docs.
[^cw-vs-sol]: https://rustopian.dev/article/from-cosmwasm-to-solana-rust-blockchain-development — store-code→code_id→instances; WASM-sandbox determinism/gas vs Solana SBF/compute-units — blog (practitioner).
[^cw-actor]: https://book.cosmwasm.com/actor-model/contract-as-actor.html — actor model, one-message-at-a-time, message types as enums, Response — book (official).
[^cw-cto]: https://medium.com/cosmwasm/cosmwasm-for-ctos-i-the-architecture-59a3e52d9b9c — actor model avoids reentrancy; A-calls-B via returned message — blog (core team).
[^cw-storage]: https://github.com/CosmWasm/cw-storage-plus — Item & Map typed storage — official-repo.
[^cw-multitest]: https://github.com/CosmWasm/cw-multi-test — cw-multi-test in-Rust simulator — official-repo.
[^cw-plus]: https://github.com/CosmWasm/cw-plus — cw20 (fungible) / cw721 (NFT) standards — official-repo.
[^cw-versions]: https://github.com/CosmWasm/cosmwasm/releases — version timeline 2.0 (Mar 2024)/2.2 (Dec 2024)/3.0 (Jun 2025, IBCv2) — official-repo.
[^ink-sdk]: https://use.ink/docs/v6/background/polkadot-sdk/ — ink! is an eDSL on Rust; contracts-on-parachains; pallet-revive; contract vs pallet — official-docs.
[^ink-macros]: https://use.ink/docs/v5/macros-attributes/ — #[ink::contract]/#[ink(storage/constructor/message/event)]/payable/selector — official-docs.
[^ink-scale]: https://use.ink/faq/ — SCALE codec; #[ink::scale_derive]; Decode requirement — official-docs (+ parity-scale-codec).
[^ink-riscv]: https://use.ink/docs/v6/background/why-riscv-and-polkavm-for-smart-contracts/ — pallet-contracts (wasmi/WASM) → pallet-revive (RISC-V/PolkaVM), 20-byte accounts, dual ink!+Solidity, cargo-contract ≥v6 — official-docs.
[^revive-status]: https://forum.polkadot.network/t/revive-smart-contracts-status-update/16366 — pallet-revive merged, enactment 2026-01-27, dual-VM, pallet-contracts legacy — official forum.
[^ink-discontinue]: https://forum.polkadot.network/t/discontinuation-of-ink-rust-smart-contract-language/16849 — Parity ceasing active ink! maintenance, Jan 2026 — official forum.
[^ink-alliance]: https://use.ink/about/ — ink! Alliance community maintenance, Treasury-funded v6 delivery — official-docs.
[^polkadot-sdk-runtime]: https://paritytech.github.io/polkadot-sdk/master/polkadot_sdk_docs/reference_docs/runtime_vs_smart_contract/index.html — authoritative runtime-vs-contract table; FRAME pallets+runtime; Substrate in Polkadot SDK; "build a chain → runtime, build a DApp → contract" — official-docs (Polkadot SDK).
[^substrate-sdk]: https://openguild.wtf/blog/polkadot/polkadot-from-substrate-to-polkadot-sdk — Substrate→Polkadot SDK consolidation; residual complexity — secondary.
