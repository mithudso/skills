<!-- hub-reference-banner -->
> **Reference file — part of the `blockchain` hub.** Formerly content of the standalone `solana` skill. Sibling topics in this family are reference files under the `blockchain` hub — not standalone skills. -->

# Solana Account Model & On-Chain State

The account model is the **shared foundation** of everything else on Solana — parallel execution, fees, programs, and tokens all build on it. Read this first.

> Crypto-primitive math (Ed25519, SHA-256) → `blockchain-crypto-primitives`. The Rust *language* → `lang-rust`. This file covers Solana's account/state/rent/PDA model only.

## Contents

- Everything is an account (the five fields)
- Data accounts vs program accounts (code/state separation)
- Ownership rules
- Account creation (System Program)
- Rent & rent-exemption (and the deprecation nuance)
- Program Derived Addresses (PDAs) & canonical bump
- Anti-patterns & pitfalls

## Everything is an account

Solana stores all on-chain state in one giant key→value store: each key is a **32-byte address** (an Ed25519 public key, or an off-curve PDA), each value is an **account**. Wallets, programs, token balances, and program state are *all* accounts.[^1][^2]

Every account has exactly **five fields**:[^1][^2]

| Field | Type | Meaning |
|---|---|---|
| `lamports` | `u64` | Balance. **1 SOL = 1,000,000,000 lamports** (10⁹). |
| `data` | `Vec<u8>` | Raw byte array — holds program state, or (for programs) bytecode. |
| `owner` | `Pubkey` | The **program** with exclusive write access to this account. |
| `executable` | `bool` | `true` = this account is a program (loadable code). |
| `rent_epoch` | `Epoch` | Legacy rent field; now deprecated (set to `u64::MAX`). |

An account's `data` is capped at **10 MiB** (`MAX_PERMITTED_DATA_LENGTH` = 10,485,760 bytes); a single instruction can grow an account by at most 10 KiB.[^2][^7] *(verified-as-of 2026-06-16)*

## Data accounts vs program accounts (the code/state split)

The `executable` flag splits accounts into two kinds:[^2][^3]

- **Program accounts** (`executable = true`) — hold sBPF bytecode (or, under the upgradeable loader, a pointer to a separate program-data account).
- **Data accounts** (`executable = false`) — hold program-defined state.

**This is the headline contrast with Ethereum.** An Ethereum contract bundles its code *and* its storage in one account. Solana **separates code from state**: programs are **stateless** and operate on data accounts passed into each instruction by reference.[^3][^10][^12] One program deployment can serve unlimited data accounts, and — because each transaction names exactly the accounts it touches — this separation is what makes parallel execution possible (see `sealevel-and-svm.md`).

> Loader nuance (version-specific): under the **Upgradeable BPF Loader (loader-v3)** a program account's `data` does **not** hold the bytecode directly — it stores a pointer to a separate **program-data account** that holds the bytecode + upgrade metadata. loader-v4 differs. Treat "bytecode lives directly in the program account's data" as loader-specific, not universal.[^3] *(verified-as-of 2026-06-16)*

## Ownership rules

The runtime enforces a small set of invariants that are the bedrock of Solana security:[^6][^7][^11]

- **Only the owner program may modify an account's `data` or *debit* (subtract) its `lamports`.**
- **Any program may *credit* (add) lamports to any writable account.** (Credit is permissionless; debit is owner-only.)
- The `owner` field can be reassigned **only when the data is zeroed**.
- Accounts default to being owned by the **System Program** (`11111111111111111111111111111111`).

These invariants are *unanimous* across sources. A native program that forgets to check `account.owner` is the single most common Solana vulnerability class — see `testing-deploy-and-security.md`.

## Account creation (the System Program)

The **System Program** is the only program that can create accounts. New accounts start System-owned, with ownership typically reassigned during creation.[^6] Creation is conceptually three steps, usually combined by the `CreateAccount` instruction:[^6][^14]

1. **Allocate** — reserve `data` space.
2. **Transfer** — fund the account with enough lamports for rent-exemption.
3. **Assign** — set the `owner` to the controlling program.

`Assign` requires the account being assigned to **sign** the transaction.[^6][^14]

## Rent & rent-exemption

Every account must hold a minimum lamport balance proportional to its `data` size — the **rent-exempt minimum**, roughly **two years' worth of rent**. This functions as a **refundable deposit**: it is fully recovered when the account is closed and its lamports are withdrawn.[^1][^11][^12]

**Current mechanics (the deprecation nuance — get this right):**[^8][^9][^11]

- Periodic rent *collection* (the per-epoch deduction that used to shrink under-funded accounts) is **disabled/deprecated** — **SIMD-0084 "Disable rent fees collection"** (Implemented).
- It is now **impossible to create a rent-paying (non-exempt) account**; creation fails unless the account is rent-exempt from the start. The last rent-paying accounts were garbage-collected.
- `rent_epoch` is deprecated and set to `u64::MAX` for rent-exempt accounts (SIMD-0267).

> **Common misconception (disconfirming):** "rent is gone on Solana" is **wrong**. The correct statement is: *periodic rent deduction is gone; the mandatory upfront rent-exemption **deposit** remains and is collected at account creation.* Solana's own docs lagged this reality for a while.[^8][^11]

> **Volatile — formula representation is mid-migration.** The rent-exempt *minimum* is unchanged in absolute lamports, but its formula is being re-expressed. The legacy docs show `(account_data_len + 128) × 3,480 lamports/byte-year × 2.0 years`. **SIMD-0194** renames `lamports_per_byte_year`→`lamports_per_byte`, doubles `3,480 → 6,960`, and drops the `exemption_threshold` `2.0 → 1.0` (to remove `f64` math) — yielding the **same absolute minimum**. Flag any hard-coded `3480`-vs-`6960` claim as version-specific and verify against the live SDK constant / `getMinimumBalanceForRentExemption` RPC. *(verified-as-of 2026-06-16)*

## Program Derived Addresses (PDAs)

A **PDA** is a 32-byte address deterministically derived from **seeds + the owning program's ID + a one-byte bump**, hashed with SHA-256 and the constant string `"ProgramDerivedAddress"`, repeated until the result falls **off** the Ed25519 curve.[^4][^5]

Key properties (all FACT — unanimous):[^4][^5][^13]

- **Off-curve ⇒ no private key exists**, so no external keypair can sign for a PDA.
- A PDA is "signed for" by its **owning program** via the runtime's **`invoke_signed`**: the runtime re-derives the address from the supplied seeds + the *calling program's* ID and adds it to the valid-signer set before privilege checks. (See `program-model-rust.md` for CPI signing.)
- **Canonical bump** = the **first** valid bump found by `find_program_address`, searching **255 → 1** (it tries the highest byte first and decrements until the result is off-curve).
- Limits: **max 16 seeds**, **32 bytes per seed**, **max 16 PDA signers per CPI**.

**Use cases:** program-owned/autonomous-authority state, deterministic user-scoped addresses (e.g. a per-user account at a known address), and avoiding keypair management entirely.

### PDA security: always use the canonical bump

**(FACT, load-bearing security rule.)** Accepting a *non-canonical* bump means **multiple valid PDAs** exist for the same logical seeds — enabling account-substitution attacks: an attacker can fake the existence or contents of a "derived" account, or replay a supposedly one-time action up to 256 times (once per bump).[^4][^13] Mitigations:

- Use `find_program_address` (canonical bump) for derivation.
- **Store the bump in the account's data** and reuse it (re-derive with `create_program_address` against the stored bump) so you never accept an attacker-supplied bump.
- Anchor enforces the canonical bump on `init` and exposes `ctx.bumps` (see `anchor-framework.md`).

## Anti-patterns & pitfalls

- **Forgetting `owner` checks** — the #1 native-program bug (any program can pass you any account). See `testing-deploy-and-security.md`.
- **Accepting non-canonical PDA bumps** — account-substitution attacks (above).
- **Assuming "no rent"** — you still must fund the rent-exempt deposit at creation, or creation fails.
- **Hard-coding the rent constant** — `3480` vs `6960` is mid-migration; query `getMinimumBalanceForRentExemption` instead.
- **Treating a program account's `data` as its bytecode** — under loader-v3 it points to a separate program-data account.

## References

[^1]: https://solana.com/docs/core/accounts/account-structure — Solana Docs (Tier 1, docs): the five account fields, types/sizes, rent-exempt balance as a refundable deposit, `rent_epoch` deprecated.
[^2]: https://solana.com/docs/core/accounts — Solana Docs (Tier 1): accounts overview — structure, address, ownership, rent; limits table incl. 10 MiB and rent-exempt formula.
[^3]: https://solana.com/docs/core/accounts/account-types — Solana Docs (Tier 1): program vs data vs system accounts; loader-v3 program-data account; two-step state-account creation.
[^4]: https://solana.com/docs/core/pda/pda-derivation — Solana Docs (Tier 1): SHA-256 off-curve derivation, canonical bump (255→1), seed limits, non-canonical attack warning.
[^5]: https://solana.com/docs/core/pda — Solana Docs (Tier 1): PDA overview — deterministic/off-curve/no-private-key, `invoke_signed` signing flow, use cases.
[^6]: https://solana.com/docs/core/programs/builtin-programs — Solana Docs (Tier 1): System Program is the sole account creator; CreateAccount/Assign/Allocate/Transfer signer rules.
[^7]: https://solana.com/docs/core/accounts/account-runtime — Solana Docs (Tier 1): runtime ownership/loader validation; non-existent accounts default System-owned w/ `rent_epoch=u64::MAX`; `MAX_PERMITTED_DATA_LENGTH` and growth limits.
[^8]: https://github.com/solana-foundation/solana-improvement-documents/blob/main/proposals/0084-disable-rent-fees-collection.md — SIMD-0084 (Tier 1, protocol proposal, Implemented): rent collection disabled; cannot create rent-paying accounts.
[^9]: https://github.com/solana-foundation/solana-improvement-documents/blob/main/proposals/0194-deprecate-rent-exemption-threshold.md — SIMD-0194 (Tier 1, paper): formula re-expression (`lamports_per_byte` 6960, threshold→1.0, same absolute minimum); related SIMD-0267 sets `rent_epoch=u64::MAX`.
[^10]: https://www.helius.dev/blog/the-solana-programming-model-an-introduction-to-developing-on-solana — Helius (Tier 2, blog): "everything is an account," stateless programs, code/data separation, Ethereum contrast, ~2-year rent-exemption.
[^11]: https://www.quicknode.com/guides/solana-development/getting-started/an-introduction-to-the-solana-account-model — QuickNode (Tier 2, blog): owner-only writes/debits, any program credits, System Program role, rent refundable, "all new accounts must be rent-exempt."
[^12]: https://www.metaplex.com/docs/solana/understanding-solana-accounts — Metaplex (Tier 2, blog): Ethereum-vs-Solana account table, programs own data accounts, explicit creation, purge-below-minimum.
[^13]: https://secure-contracts.com/not-so-smart-contracts/solana/improper_pda_validation/ — Trail of Bits "Not So Smart Contracts" (Tier 2, security): non-canonical bump substitution attack mechanics (existence/content spoofing, 256× replay).
[^14]: https://github.com/solana-labs/solana/blob/master/sdk/program/src/system_instruction.rs — solana-labs source (Tier 1): three-step creation (allocate/transfer/assign), `MAX_PERMITTED_DATA_LENGTH=10,485,760`, assign-requires-signer.
