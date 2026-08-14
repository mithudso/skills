<!-- hub-reference-banner -->
> **Reference file — part of the `blockchain` hub.** Formerly the standalone `bitcoin-protocol-expert` skill.
> Sibling topics in this family are now reference files under the hubs (`blockchain`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: bitcoin-protocol-expert
description: "Bitcoin protocol & ecosystem expert — base-layer mechanics, upgrades, and the L2/metaprotocol landscape. TRIGGER: the UTXO model & coin selection; tx structure (txid/wtxid, weight/vbytes, implicit fees); Bitcoin Script & addresses (P2PKH/P2SH/P2WPKH/P2WSH/P2TR, multisig, CLTV/CSV, bech32/bech32m); block/header structure & Merkle root; PoW mining, difficulty retargeting, the halving & 21M cap; mempool & fees (sat/vByte, RBF, CPFP, TRUC); SegWit, Taproot/Schnorr/MAST, soft-fork activation; Lightning (channels, HTLCs, routing, BOLTs); ordinals, BRC-20, Runes; BitVM & the Bitcoin-L2/restaking frontier (sidechains, drivechains, covenants, Babylon); full nodes, wallets, UTXO-vs-account design. SKIP: consensus theory (PoW-as-a-class, BFT, PoS, finality) → distributed-systems-consensus; crypto-primitive math (ECDSA/Schnorr/SHA-256) → blockchain-crypto-primitives; token-economics/valuation → blockchain-economics; other whole chains (Ethereum/Solana) → per-chain skills; MongoDB → mongodb-* hubs."
version: 1.0.1
updated: 2026-06-16
category: blockchain
whenToUse:
  - "Explaining or building with the UTXO model, transaction structure, or coin selection"
  - "Writing or debugging Bitcoin Script and addresses (P2PKH/P2SH/P2WPKH/P2WSH/P2TR, multisig, CLTV/CSV timelocks)"
  - "Reasoning about block/header structure, PoW mining, difficulty retargeting, the halving, or the 21M cap"
  - "Estimating fees or bumping a stuck transaction (RBF/full-RBF, CPFP, TRUC, sat/vByte)"
  - "Working with the SegWit or Taproot upgrades, or soft-fork activation (BIP9/BIP8/UASF/Speedy Trial)"
  - "Designing on the Lightning Network (payment channels, HTLCs, routing, BOLTs)"
  - "Explaining ordinals/inscriptions, BRC-20, or Runes, or evaluating a BitVM / Bitcoin-L2 / Babylon claim"
  - "Running a full node or choosing a wallet custody model (hardware/multisig/PSBT)"
  - "Comparing Bitcoin's UTXO design against account-model chains"
keywords:
  - bitcoin
  - UTXO
  - bitcoin script
  - P2WPKH P2TR taproot
  - proof of work mining difficulty halving
  - mempool fee RBF CPFP sat/vByte
  - segwit taproot schnorr MAST soft fork
  - lightning network HTLC payment channel BOLT
  - ordinals inscriptions runes BRC-20
  - BitVM bitcoin L2 drivechain covenant babylon staking
  - full node wallet custody multisig PSBT
  - account model vs UTXO
tags:
  - blockchain
  - bitcoin
  - cryptocurrency
  - lightning-network
  - layer-2
  - protocol
metadata:
  changelog:
    - "2026-06-16 sko v1.0.0->v1.0.1 — Pass H 10/10 pos, 0/10 neg (predicted); fixed 4 Medium (2 dangling citation anchors [^35]/[^15]; whenToUse prose->9-item list; description 1493->1000 chars under Glean cap); 0 banned terms; blind re-audit clean"
---

# Bitcoin Protocol & Ecosystem Expert

Expert reference for the **Bitcoin protocol and its ecosystem**: the base-layer data model and mechanics, the consensus-rule upgrades that shaped it, the Lightning Network, the ordinals/runes metaprotocols, the emerging L2/restaking landscape, and the practical operation (nodes, wallets) plus the comparative design against account-model chains.

> **Scope & peer deferral.** This is the *Bitcoin-specific* spoke of a blockchain skill family.
> - **General consensus theory** — Proof-of-Work *as a mechanism class*, BFT, Nakamoto consensus, Proof-of-Stake, longest-chain/fork-choice, finality, the scalability trilemma, Sybil resistance — belongs to **`distributed-systems-consensus`**. This skill covers only Bitcoin's *concrete* mining/validation mechanics. (`distributed-systems-consensus` explicitly defers whole-Bitcoin to here.)
> - **Cryptographic primitive internals** — ECDSA/Schnorr signature math, SHA-256/RIPEMD-160 internals, Merkle-proof construction math, BIP32 HD-derivation math — belong to **`blockchain-crypto-primitives`**. This skill references them by name and explains *what they do for Bitcoin*, never re-deriving the math.
> - **Token-economics / valuation theory** — stock-to-flow, monetary-premium models, fee-market valuation — belong to **`blockchain-economics`**. This skill covers only Bitcoin's *concrete issuance mechanics* (subsidy schedule, 21M cap).
> - Other whole chains (Ethereum, Solana, Cardano) → their per-chain skills. The account-vs-UTXO *comparison* lives here as design context.

## How to use this skill

The base-layer model is the **shared foundation** — read it first; everything else builds on it. Load the reference file matching the task:

| Topic | Reference |
|---|---|
| UTXO model, transaction structure, Script, addresses, timelocks | `references/utxo-transactions-and-script.md` |
| Block/header structure, Merkle root, PoW mining, difficulty, halving, 21M supply | `references/blocks-mining-and-monetary-policy.md` |
| Mempool, fee estimation, RBF/full-RBF, CPFP, TRUC/package relay; SegWit, Taproot, soft-fork activation | `references/mempool-fees-and-upgrades.md` |
| Lightning Network (channels, HTLCs/PTLCs, routing, BOLTs, implementations, limits) | `references/lightning-network.md` |
| Ordinals/inscriptions, BRC-20, Runes, BitVM, sidechains/drivechains, covenants debate, Babylon | `references/ordinals-runes-and-l2-landscape.md` |
| Running a full node, wallet custody models, UTXO-vs-account comparison | `references/nodes-wallets-and-account-model.md` |

## Core mental model (the load-bearing facts)

**Bitcoin is a UTXO ledger, not an account ledger.** There is no stored "balance" anywhere on-chain. Value exists as discrete **Unspent Transaction Outputs (UTXOs)**; a wallet's balance is the sum of every UTXO its keys can unlock.[^1] To spend, you consume whole UTXOs as **inputs** and create new **outputs**, returning the remainder to yourself as a **change output**. The leftover that isn't assigned to any output is the **fee** — fees are *implicit* (`sum(inputs) − sum(outputs)`), with no explicit fee field.[^2]

**Spending conditions are expressed in Bitcoin Script** — a stack-based, deliberately **non-Turing-complete** (no loops) language. Each output carries a locking script (`scriptPubKey`); spending it requires an unlocking script (`scriptSig` legacy, or the **witness** post-SegWit) that satisfies it.[^3] The standard output types and their addresses: **P2PKH** (`1…`, base58), **P2SH** (`3…`, base58), **P2WPKH/P2WSH** (`bc1q…`, bech32, SegWit v0), **P2TR** (`bc1p…`, bech32m, Taproot/SegWit v1).[^2][^3]

**Blocks chain via Proof-of-Work over an 80-byte header.** Miners compute double-SHA256 of the header (version, prev-block-hash, **Merkle root**, time, nBits/target, nonce) seeking a hash below the target.[^4] Difficulty **retargets every 2016 blocks** (~2 weeks) to hold a **10-minute** average block time, clamped to a 4× swing per period.[^5] The first transaction is the **coinbase**, paying the miner **subsidy + fees**; coinbase outputs are unspendable for **100 blocks**.[^4]

**Issuance is capped and disinflationary.** The block subsidy **halves every 210,000 blocks** (~4 years): 50 → 25 → 12.5 → 6.25 → **3.125 BTC after the April 2024 halving**.[^6] The supply asymptotically approaches but stays **just under 21 million** (≈20,999,999.9769 BTC — integer-satoshi truncation of the subsidy), with issuance effectively ending ~2140, after which miners earn fees only.[^6] *(verified-as-of 2026-06-16)*

**Two consensus upgrades dominate the modern protocol.** **SegWit** (2017) moved the witness out of the txid preimage — fixing third-party transaction malleability (a prerequisite for Lightning) and redefining the block limit as **4,000,000 weight units** with a witness discount.[^7] **Taproot** (Nov 2021) added **Schnorr signatures**, **Pay-to-Taproot** with key-path vs script-path spends, and **MAST**, improving privacy and multisig efficiency.[^7] Both are **soft forks** (backward-compatible rule *tightening*); activation methods are Bitcoin's governance flashpoint.[^7]

**Layers and metaprotocols sit on top.** The **Lightning Network** is the primary L2 for fast/cheap payments via bilateral payment channels and HTLC-routed multi-hop payments.[^8] **Ordinals/inscriptions, BRC-20, and Runes** are metaprotocols using the witness or OP_RETURN to embed data/tokens.[^9] **BitVM, sidechains/drivechains, and Babylon** are an *early, experimental, and heavily contested* L2/restaking frontier — most "Bitcoin L2" branding is marketing, and bridges are the recurring trust weak point.[^9] *(verified-as-of 2026-06-16)*

## Cross-domain routing notes

- A question about **whether PoW is a good consensus mechanism**, BFT bounds, or PoS comparisons → `distributed-systems-consensus`. A question about **how Bitcoin mining concretely works** → here (`references/blocks-mining-and-monetary-policy.md`).
- A question about **how Schnorr/ECDSA signing works mathematically** → `blockchain-crypto-primitives`. A question about **what Schnorr enabled in Taproot** (key aggregation, P2TR, PTLCs) → here.
- A question about **Bitcoin's price/monetary value** → `blockchain-economics`. A question about **the exact issuance schedule** → here.

## References (top-level sources)

The deep per-topic citations live in each reference file's own `## References`. Anchors used above:

[^1]: https://developer.bitcoin.org/devguide/transactions.html — Bitcoin Developer Guide: transaction anatomy, UTXO consumption, change, implicit fees.
[^2]: https://learnmeabitcoin.com/technical/transaction/ — transaction serialization, value-in-satoshis, fee-as-leftover, vbytes.
[^3]: https://en.bitcoin.it/wiki/Script — Bitcoin Script: stack-based, non-Turing-complete, locking vs unlocking; standard output templates.
[^4]: https://developer.bitcoin.org/reference/block_chain.html — block header fields, coinbase, double-SHA256, Merkle root.
[^5]: https://en.bitcoin.it/wiki/Difficulty — 2016-block retarget, 10-minute target, 4× clamp.
[^6]: https://en.bitcoin.it/wiki/Controlled_supply — halving schedule, ~20,999,999.9769 BTC asymptote, subsidy→0.
[^7]: https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki — SegWit (BIP141); Taproot per BIPs 340/341/342; soft-fork activation per BIP9/BIP8.
[^8]: https://lightning.network/lightning-network-paper.pdf — Lightning Network whitepaper (payment channels, penalty mechanism).
[^9]: https://docs.ordinals.com/ — Ordinal Theory; plus Runes (rodarmor.com/blog/runes/), BitVM (bitvm.org), Babylon (docs.babylonlabs.io) — see the L2 reference for the full, skeptically-sourced survey.
