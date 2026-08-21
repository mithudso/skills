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

Expert reference for **Bitcoin protocol and ecosystem**: base-layer data model, consensus-rule upgrades, Lightning Network, ordinals/runes metaprotocols, emerging L2/restaking landscape, node/wallet operation, and design comparison vs account-model chains.

> **Scope & peer deferral.** Bitcoin-specific spoke of blockchain skill family.
> - **General consensus theory** — PoW as mechanism class, BFT, Nakamoto consensus, PoS, longest-chain/fork-choice, finality, scalability trilemma, Sybil resistance — belongs to **`distributed-systems-consensus`**. Skill covers only Bitcoin's *concrete* mining/validation mechanics. (`distributed-systems-consensus` defers whole-Bitcoin here.)
> - **Cryptographic primitive internals** — ECDSA/Schnorr math, SHA-256/RIPEMD-160, Merkle-proof math, BIP32 HD-derivation math — belong to **`blockchain-crypto-primitives`**. Skill references by name, explains *what they do for Bitcoin*, never re-derives math.
> - **Token-economics / valuation theory** — stock-to-flow, monetary-premium models, fee-market valuation — belong to **`blockchain-economics`**. Skill covers only Bitcoin's *concrete* issuance mechanics (subsidy schedule, 21M cap).
> - Other chains (Ethereum, Solana, Cardano) → per-chain skills. Account-vs-UTXO *comparison* lives here as design context.

## How to use this skill

Base-layer model = **shared foundation** — read first; everything builds on it. Load reference file matching task:

| Topic | Reference |
|---|---|
| UTXO model, transaction structure, Script, addresses, timelocks | `references/utxo-transactions-and-script.md` |
| Block/header structure, Merkle root, PoW mining, difficulty, halving, 21M supply | `references/blocks-mining-and-monetary-policy.md` |
| Mempool, fee estimation, RBF/full-RBF, CPFP, TRUC/package relay; SegWit, Taproot, soft-fork activation | `references/mempool-fees-and-upgrades.md` |
| Lightning Network (channels, HTLCs/PTLCs, routing, BOLTs, implementations, limits) | `references/lightning-network.md` |
| Ordinals/inscriptions, BRC-20, Runes, BitVM, sidechains/drivechains, covenants debate, Babylon | `references/ordinals-runes-and-l2-landscape.md` |
| Running a full node, wallet custody models, UTXO-vs-account comparison | `references/nodes-wallets-and-account-model.md` |

## Core mental model (the load-bearing facts)

**Bitcoin = UTXO ledger, not account ledger.** No stored "balance" on-chain. Value exists as discrete **Unspent Transaction Outputs (UTXOs)**; wallet balance = sum of every UTXO its keys can unlock.[^1] Spending: consume whole UTXOs as **inputs**, create new **outputs**, return remainder as **change output**. Leftover not assigned to output = **fee** — fees *implicit* (`sum(inputs) − sum(outputs)`), no explicit fee field.[^2]

**Spending conditions expressed in Bitcoin Script** — stack-based, deliberately **non-Turing-complete** (no loops). Each output carries locking script (`scriptPubKey`); spending requires unlocking script (`scriptSig` legacy, or **witness** post-SegWit) that satisfies it.[^3] Standard output types and addresses: **P2PKH** (`1…`, base58), **P2SH** (`3…`, base58), **P2WPKH/P2WSH** (`bc1q…`, bech32, SegWit v0), **P2TR** (`bc1p…`, bech32m, Taproot/SegWit v1).[^2][^3]

**Blocks chain via Proof-of-Work over 80-byte header.** Miners compute double-SHA256 of header (version, prev-block-hash, **Merkle root**, time, nBits/target, nonce) seeking hash below target.[^4] Difficulty **retargets every 2016 blocks** (~2 weeks) to hold **10-minute** average block time, clamped to 4× swing per period.[^5] First transaction = **coinbase**, paying miner **subsidy + fees**; coinbase outputs unspendable for **100 blocks**.[^4]

**Issuance capped and disinflationary.** Block subsidy **halves every 210,000 blocks** (~4 years): 50 → 25 → 12.5 → 6.25 → **3.125 BTC after April 2024 halving**.[^6] Supply asymptotically approaches but stays **just under 21 million** (≈20,999,999.9769 BTC — integer-satoshi truncation), issuance ends ~2140, miners earn fees only after.[^6] *(verified-as-of 2026-06-16)*

**Two consensus upgrades dominate modern protocol.** **SegWit** (2017) moved witness out of txid preimage — fixes third-party transaction malleability (Lightning prerequisite), redefines block limit as **4,000,000 weight units** with witness discount.[^7] **Taproot** (Nov 2021) added **Schnorr signatures**, **Pay-to-Taproot** with key-path vs script-path spends, and **MAST**, improving privacy and multisig efficiency.[^7] Both = **soft forks** (backward-compatible rule *tightening*); activation methods = Bitcoin's governance flashpoint.[^7]

**Layers and metaprotocols sit on top.** **Lightning Network** = primary L2 for fast/cheap payments via bilateral payment channels and HTLC-routed multi-hop payments.[^8] **Ordinals/inscriptions, BRC-20, Runes** = metaprotocols using witness or OP_RETURN to embed data/tokens.[^9] **BitVM, sidechains/drivechains, Babylon** = *early, experimental, contested* L2/restaking frontier — most "Bitcoin L2" branding is marketing, bridges = recurring trust weak point.[^9] *(verified-as-of 2026-06-16)*

## Cross-domain routing notes

- **Whether PoW is good consensus mechanism**, BFT bounds, PoS comparisons → `distributed-systems-consensus`. **How Bitcoin mining concretely works** → here (`references/blocks-mining-and-monetary-policy.md`).
- **How Schnorr/ECDSA signing works mathematically** → `blockchain-crypto-primitives`. **What Schnorr enabled in Taproot** (key aggregation, P2TR, PTLCs) → here.
- **Bitcoin price/monetary value** → `blockchain-economics`. **Exact issuance schedule** → here.

## References (top-level sources)

Deep per-topic citations in each reference file's own `## References`. Anchors used above:

[^1]: https://developer.bitcoin.org/devguide/transactions.html — Bitcoin Developer Guide: transaction anatomy, UTXO consumption, change, implicit fees.
[^2]: https://learnmeabitcoin.com/technical/transaction/ — transaction serialization, value-in-satoshis, fee-as-leftover, vbytes.
[^3]: https://en.bitcoin.it/wiki/Script — Bitcoin Script: stack-based, non-Turing-complete, locking vs unlocking; standard output templates.
[^4]: https://developer.bitcoin.org/reference/block_chain.html — block header fields, coinbase, double-SHA256, Merkle root.
[^5]: https://en.bitcoin.it/wiki/Difficulty — 2016-block retarget, 10-minute target, 4× clamp.
[^6]: https://en.bitcoin.it/wiki/Controlled_supply — halving schedule, ~20,999,999.9769 BTC asymptote, subsidy→0.
[^7]: https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki — SegWit (BIP141); Taproot per BIPs 340/341/342; soft-fork activation per BIP9/BIP8.
[^8]: https://lightning.network/lightning-network-paper.pdf — Lightning Network whitepaper (payment channels, penalty mechanism).
[^9]: https://docs.ordinals.com/ — Ordinal Theory; plus Runes (rodarmor.com/blog/runes/), BitVM (bitvm.org), Babylon (docs.babylonlabs.io) — see L2 reference for full, skeptically-sourced survey.