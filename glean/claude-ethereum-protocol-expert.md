# ethereum-protocol-expert

**Category:** Databases, Data Engineering & Analytics
**Platform:** Claude
**Original Path:** claude/standalone/ethereum-protocol-expert

## Description
Ethereum protocol expert, how the chain works beneath the app layer (execution, consensus, scaling). TRIGGER: the account model (EOA vs contract accounts, nonce/balance/storageRoot/codeHash, EIP-7702), world state & the Merkle-Patricia trie (RLP, Verkle debate); the EVM (stack machine, opcodes, memory vs storage vs calldata, CALL/DELEGATECALL/STATICCALL/CREATE2, in-EVM gas, cold/warm access EIP-2929, precompiles); the EIP-1559 fee market (base fee + priority fee + burn, typed txs, mempool, MEV-Boost/PBS, tx & block lifecycle, receipts); Proof of Stake (Beacon Chain, the Merge, 32-ETH validators/staking/withdrawals, EIP-7251, slots/epochs/RANDAO, attestations, Gasper = Casper FFG + LMD-GHOST, slashing, finality/weak-subjectivity); data & scaling (EIP-4844 blobs/proto-danksharding, KZG ceremony, data-availability sampling, PeerDAS, danksharding, the rollup-centric roadmap); Layer-2 (optimistic vs ZK rollups, fraud vs validity proofs, sequencers, rollup vs validium, the 7-day challenge window, L2BEAT stages). SKIP: Solidity/Vyper/Foundry/ERCs & contract DEVELOPMENT → ethereum-smart-contract-development; contract auditing/exploits → contract-security; SNARK/STARK & circuit MATH → zero-knowledge-proofs (here: only validity proofs' ROLE in scaling); consensus theory (CAP/FLP/BFT proofs) → distributed-systems-consensus; crypto primitives (ECDSA, BLS, KZG math, Keccak) → blockchain-crypto-primitives; token-economics → blockchain-economics; other chains → per-chain skills; MongoDB → mongodb-* hubs.

---

# Ethereum Protocol Expert

Expert reference for how Ethereum works **beneath the application layer**, the execution model, the fee market, Proof-of-Stake consensus, the data-availability scaling roadmap, and Layer-2 rollups. This skill is protocol- and infrastructure-focused: it explains *how the machine runs*, not how to write the contracts that run on it.

> **Scope boundary (read first).** This skill covers the **protocol**. It deliberately defers adjacent domains to peer skills — see **SKIP** in the frontmatter and the **Peer skills & deferral** section below. In short: contract *development* (Solidity/Vyper/Foundry/ERCs) → `ethereum-smart-contract-development`; contract *security/auditing* → `contract-security`; ZK *proving-system math* → `zero-knowledge-proofs`; general *consensus theory* → `distributed-systems-consensus`; *cryptographic primitives* (ECDSA/BLS/KZG/Keccak math) → `blockchain-crypto-primitives`; *token economics* → `blockchain-economics`. **Note on peer status:** `distributed-systems-consensus` is installed and reciprocally defers per-chain Gasper back to this skill; the other named peers are planned siblings of this blockchain family that may not be installed yet. Where a peer is absent, that domain is simply out of scope here (the deferral is the honest boundary, not a broken link).

> **Volatility warning.** Ethereum ships hard forks roughly twice a year and several facts in this skill change at each fork. Every version-pinned claim is stamped `verified-as-of: 2026-06-16`. The fastest-moving facts — block gas limit, blob target/max counts, which fork is live, which L2 is at which L2BEAT stage, ePBS status — **must be re-checked against a primary source** (ethereum.org, the relevant EIP, L2BEAT's live dashboard, or a client release) before you rely on them. See each reference file's header for the specific volatile claims.

## How to use this skill

Five reference files, one per domain. Load the one matching the question:

| Question is about… | Read |
| --- | --- |
| Accounts, world state, the trie, **the EVM**, opcodes, in-EVM gas, precompiles, CALL/DELEGATECALL | `references/execution-state-and-evm.md` |
| Gas pricing, **EIP-1559** base/priority fee + burn, typed txs, the mempool, MEV/PBS, tx & block lifecycle, receipts | `references/gas-fee-market-and-transactions.md` |
| **Proof of Stake**, the Beacon Chain & Merge, validators/staking/withdrawals, attestations, **Gasper** (FFG + LMD-GHOST), slashing, finality | `references/proof-of-stake-consensus.md` |
| **EIP-4844 blobs**, the blob-gas market, KZG ceremony, **data-availability sampling**, **PeerDAS**, **danksharding**, the rollup-centric roadmap | `references/data-availability-and-danksharding.md` |
| **Layer-2 rollups**, optimistic vs ZK, fraud vs validity proofs, sequencers, rollup-vs-validium DA, challenge windows, **L2BEAT stages** | `references/layer-2-rollups.md` |

## The protocol in one picture (post-Merge)

Ethereum is a single, globally replicated state machine. Since **The Merge** (Sept 15, 2022) every node runs two coupled clients joined by the **Engine API**:

- **Execution layer (EL)**, holds the world state, runs the **EVM**, gossips transactions, prices gas (EIP-1559). Clients: Geth, Nethermind, Reth, Besu, Erigon.
- **Consensus layer (CL)**, runs **Proof of Stake** (the Beacon Chain): validators, slots/epochs, attestations, and **Gasper** finality. Clients: Prysm, Lighthouse, Teku, Nimbus, Lodestar.

A block proposer's CL asks its EL to build an **execution payload** (the transaction-execution half of the block), wraps it in a beacon block, and gossips it; attesters vote on it; Gasper finalizes it ~2 epochs later. Scaling happens **off** L1 in **rollups** that post their data to L1 **blobs** (EIP-4844) and inherit L1 security, the **rollup-centric roadmap**.

`verified-as-of: 2026-06-16`

## Core mental models (the load-bearing 12)

These are the ideas the references expand. Each is sourced in its reference file.

1. **Two account types, four fields.** Every account is `(nonce, balance, storageRoot, codeHash)`. EOAs are key-controlled and can *originate* transactions; contract accounts are code-controlled and only act when called. EIP-7702 (Pectra) lets an EOA *temporarily* point its code field at a contract — so `tx.origin == msg.sender` and "is this an EOA" checks are no longer reliable. → `execution-state-and-evm.md`
2. **The state root commits everything.** The world state isn't stored on-chain; only the `stateRoot` (root of the Merkle-Patricia trie) is in the block header. The MPT is a *hexary* trie with Merkle hashing — **not** a binary Merkle tree. Its replacement (Verkle vs binary tree) is an open roadmap question as of 2026. → `execution-state-and-evm.md`
3. **The EVM is a quasi-Turing-complete 256-bit stack machine.** "Quasi" because gas bounds every execution, solving the halting problem. Three data locations with wildly different cost/lifetime: **storage** (persistent, expensive), **memory** (per-call, cheap), **calldata** (read-only input, cheapest). → `execution-state-and-evm.md`
4. **Gas meters work; the fee market prices it.** *In-EVM* gas accounting (per-opcode cost, cold/warm access via EIP-2929, the 21,000 intrinsic floor) is separate from the *fee market* (EIP-1559) that decides the price per gas. → `execution-state-and-evm.md` (accounting) + `gas-fee-market-and-transactions.md` (pricing)
5. **EIP-1559 = a burned, algorithmic base fee + a tip.** The base fee is protocol-computed from the parent block, moves ≤±12.5%/block toward a 50%-full target, and is **burned**. The priority fee ("tip") goes to the proposer. EIP-1559 did **not** lower fees — L2s + blobs did. → `gas-fee-market-and-transactions.md`
6. **Transactions are typed (EIP-2718).** Five types are defined so far — 0x00 (legacy) through 0x04 (EIP-7702 set-code) — within EIP-2718's 0x00–0x7f envelope. Type 0x02 (dynamic-fee) is the default; Type 0x03 carries blobs. → `gas-fee-market-and-transactions.md`
7. **Most transactions no longer go through the public mempool.** A large and growing share route through private order flow / builders (MEV protection). Block building is dominated by **MEV-Boost** (out-of-protocol PBS); **enshrined PBS** is a roadmap direction, not yet live. → `gas-fee-market-and-transactions.md`
8. **PoS secures the chain by staked capital, not work.** 32 ETH activates a validator (EIP-7251 raised only the *ceiling* to 2048, not the floor). Validators **propose** and **attest**; correct behavior earns rewards, being offline causes small penalties (and a quadratic **inactivity leak** only during non-finality), provable equivocation causes **slashing**. → `proof-of-stake-consensus.md`
9. **Gasper = Casper FFG (finality) + LMD-GHOST (fork choice).** LMD-GHOST picks the head by accumulated attestation weight every slot; Casper FFG justifies/finalizes epoch-boundary checkpoints with a 2/3 stake supermajority. Finality is fast (~12.8 min) but **not instant and not objective** — hence weak-subjectivity checkpoints. → `proof-of-stake-consensus.md`
10. **Scaling is data-availability, not execution sharding.** Execution sharding was **abandoned**. The plan is to make L1 a cheap, guaranteed-available **data layer** for rollups. **Blobs are ephemeral** (~18 days), not permanent storage. → `data-availability-and-danksharding.md`
11. **Danksharding scales by sampling, not by downloading.** Proto-danksharding (EIP-4844) shipped the blob *format* and a separate blob-gas market; **data-availability sampling** (PeerDAS now, full danksharding later) lets nodes verify data is available by sampling small chunks of erasure-coded data instead of downloading all of it. → `data-availability-and-danksharding.md`
12. **A rollup inherits L1 security; a sidechain does not.** The test: you can exit your assets to L1 even if everyone else colludes. **Optimistic** rollups assume validity and allow fraud proofs (long withdrawal window, needs ≥1 honest watcher); **ZK/validity** rollups prove every batch (fast trustless withdrawal). The real trust boundary on most L2s today is the **centralized sequencer + the upgrade-key multisig**, captured by the **L2BEAT Stage 0/1/2** framework. → `layer-2-rollups.md`

## Cross-cutting "stale knowledge" traps (re-verify before relying)

The single most common way to be wrong about Ethereum is to quote a value that a fork has since changed. As of `2026-06-16` the live values are:

- **Block gas limit ≈ 60M** (EIP-7935, Fusaka, Dec 2025) — *not* the long-quoted 30M.
- **Blob target/max = 14/21** (BPO2, Jan 7, 2026) — *not* the Dencun-era 3/6 or Pectra-era 6/9.
- **EIP-7702 (Type-4 set-code EOAs) and EIP-7251 (MaxEB 2048) are live** (Pectra, May 7, 2025).
- **PeerDAS is live** (Fusaka, Dec 3, 2025); **full danksharding is NOT** shipped.
- **ePBS (enshrined PBS) is NOT live** — scheduled for a future fork (Glamsterdam direction).
- **No general-purpose L2 has reached L2BEAT Stage 2**; most secured value sits at Stage 1.
- **32 ETH is still the validator activation minimum** after Pectra.

Treat all of the above as volatile; the reference files carry the sourcing and the "what it was before" history.

## Peer skills & deferral

This skill is one spoke of a blockchain skill family. Route out when the question crosses a boundary:

| If the question is about… | Defer to | Why |
| --- | --- | --- |
| Writing/testing contracts — Solidity, Vyper, Foundry/Hardhat, ERC-20/721/1155, proxies, account-abstraction *implementation* | `ethereum-smart-contract-development` | This skill stops at the protocol; it explains DELEGATECALL and EIP-7702 mechanics, not how to write an upgradeable contract. |
| Contract vulnerabilities, audits, exploit classes (reentrancy, oracle manipulation), MEV *attacks* on apps | `contract-security` | Security review is a distinct discipline. |
| How a SNARK/STARK actually proves — circuits, arithmetization, proving systems, the math | `zero-knowledge-proofs` | This skill covers only the *role* validity proofs play in L2 scaling and trust, never the cryptography. |
| Consensus theory as a field — CAP/PACELC, FLP, BFT/3f+1, safety vs liveness proofs, Nakamoto/PoW-as-a-class, other chains' consensus | `distributed-systems-consensus` | That skill owns the *theory* (incl. the formal Gasper safety/liveness analysis); this skill covers Gasper at a *protocol-user* level. |
| Cryptographic primitives — ECDSA/secp256k1, BLS12-381, KZG polynomial commitments, Keccak, Merkle-proof math | `blockchain-crypto-primitives` | This skill *names and uses* these (e.g., the KZG point-evaluation precompile) but never derives them. |
| Token economics, monetary models, ETH valuation, "ultrasound money" as a thesis | `blockchain-economics` | This skill states the burn *mechanism* and notes the deflation claim is contested; the economic modeling lives there. |
| Bitcoin, Solana, or another whole chain | the per-chain skill (e.g., `bitcoin-protocol-expert`) | Per-chain protocol depth lives in per-chain skills. |
| MongoDB / Atlas | the `mongodb-*` hubs | Unrelated. |

When this skill needs a primitive (KZG, BLS, ECDSA) or a theory result (BFT bounds), it **references** the peer by name and keeps its own treatment at the protocol-user level.

## References

This SKILL.md is the index and the mental-model layer. The load-bearing detail, all inline-cited and stamped, lives in the five files under `references/`. Read the one that matches the question (table above); each opens with a Contents block and its own volatility header.