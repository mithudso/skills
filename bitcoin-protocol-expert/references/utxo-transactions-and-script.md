<!-- Provenance: reference under the standalone `bitcoin-protocol-expert` skill. Created 2026-06-16 via /dr deep-research. Peer deferral: crypto math → blockchain-crypto-primitives; consensus theory → distributed-systems-consensus. -->

# UTXO Model, Transaction Structure & Bitcoin Script

The shared foundation of the Bitcoin protocol. Everything else (mining, fees, upgrades, Lightning) builds on this layer.

## Contents
- [The UTXO model](#the-utxo-model)
- [Coin selection, change & dust](#coin-selection-change--dust)
- [Transaction structure](#transaction-structure)
- [txid vs wtxid & malleability](#txid-vs-wtxid--malleability)
- [Weight, virtual bytes & implicit fees](#weight-virtual-bytes--implicit-fees)
- [Bitcoin Script](#bitcoin-script)
- [Standard output types & addresses](#standard-output-types--addresses)
- [Multisig](#multisig)
- [Timelocks](#timelocks)
- [Address encodings](#address-encodings)
- [HD wallets (name level)](#hd-wallets-name-level)
- [References](#references)

## The UTXO model

A **UTXO (Unspent Transaction Output)** is a discrete, indivisible chunk of bitcoin value created as an *output* of one transaction and consumed *in full* as an *input* of a later one. Each is spendable exactly once. Bitcoin stores **no balance** anywhere on-chain; a wallet's balance is the **sum of every UTXO its keys can unlock**.[^1][^2] The set of all currently-spendable outputs is the **UTXO set**, which full nodes track to validate spends.[^2]

**UTXO vs account/balance model** (fact — 3+ sources). The *account model* (Ethereum, Tezos) keeps a mutable per-address balance, debited/credited like a bank ledger. The *UTXO model* (Bitcoin, Litecoin, Cardano) treats each output like a distinct banknote: to pay, select specific "notes," consume them whole, receive change.[^1][^3] Consequences:[^3][^4]
- **Statelessness & parallelism** — independent UTXOs share no state, so transactions spending *different* UTXOs can be validated in parallel. The account model needs a per-address **nonce** to prevent replay, forcing same-sender transactions into strict order.
- **Privacy** — value fragments across many outputs/addresses rather than one running balance.
- **Trade-offs** (qualified — UTXO is *not* strictly superior) — statelessness is awkward for stateful smart contracts; Cardano's **Extended UTXO (eUTXO)** exists precisely to carry datums/state on outputs. Tiny-output proliferation ("dust") can bloat the set. See `references/nodes-wallets-and-account-model.md` for the full comparison.

## Coin selection, change & dust

**Coin selection** (fact). Since outputs are consumed whole, a wallet chooses *which* UTXOs to spend to cover amount + fee. Bitcoin Core uses **Branch and Bound (BnB)**, **Knapsack**, and **Single Random Draw**, picking among solutions via a "waste" metric.[^7] BnB searches UTXOs by *effective value* (value minus the fee to spend that input) for an exact-match set that **avoids a change output** — lowering current and future fees, shrinking the transaction graph, and reducing address-linkage privacy leaks.[^7][^8]

**Change outputs** (fact). Inputs rarely sum to exactly (amount + fee), so the remainder returns to the sender as a **change output** (usually a fresh internal address). Anything not assigned to an output is the fee.[^1]

**Dust** (qualify — exact constants are version-dependent policy). A "dust" output costs more in fees to spend than it is worth. Bitcoin Core defines dust relative to **`dustRelayFee`** (default 3000 sat/kvB): an output is dust if spending it would cost more than **1/3** of its value.[^9][^10] Concretely a legacy P2PKH output is dust below **546 satoshis**; a SegWit P2WPKH output below **294 satoshis**.[^9][^10] These are *relay/standardness policy*, not consensus rules.

## Transaction structure

A SegWit transaction serializes in this order (fact — 3+ sources incl. BIP141):[^11][^12][^14]
1. **version** (4 B, LE) — format/feature version; v2 enables BIP68 relative locktime.
2. **marker** (`0x00`) + **flag** (`0x01`) — *SegWit serialization only*; the `0x00` marker disambiguates from a legacy tx.
3. **input count** (compactSize).
4. **inputs**, each: **outpoint** = previous **txid** (32 B) + **vout** (4 B output index); **scriptSig** (unlocking script; empty for native SegWit); **nSequence** (4 B — relative locktime per BIP68, and RBF/locktime opt-in signaling).
5. **output count** (compactSize).
6. **outputs**, each: **value** (8 B, LE) in **satoshis** (1 BTC = 100,000,000 sat); **scriptPubKey** (locking script).
7. **witness** — *SegWit only*; one field per input (signatures, pubkeys, scripts).
8. **nLockTime** (4 B) — earliest height/time the tx may be mined; a height if < 500,000,000, else a Unix timestamp. A tx is "final" if locktime is past **or** every input's sequence is `0xffffffff`.

## txid vs wtxid & malleability

Both are **double-SHA256** hashes over different serializations (fact — BIP141):[^14][^16]
- **txid** = dSHA256 of the *legacy* serialization `[version][inputs][outputs][locktime]` — **excludes** marker, flag, witness.
- **wtxid** = dSHA256 of the *full* serialization including witness.
- For a tx with no witness data, **wtxid == txid**. The coinbase's witness commitment is the Merkle root over wtxids.[^14]

**Transaction malleability** (fact — confirms txid *was* malleable). Pre-SegWit the txid is malleable: a third party could alter `scriptSig` encoding (non-minimal pushes, ECDSA `S`→`-S`) to produce a different-but-valid tx with a different txid, because the signature is taken over a form that excludes the input scripts.[^17] Low-S enforcement (Bitcoin Core PR #6769, 2015) curbed one vector; **BIP62 was proposed then withdrawn**. **SegWit (BIP141)** structurally fixes it by moving the witness out of the txid preimage — a prerequisite for safe pre-signed transaction chains (e.g., Lightning). **BIP147 (NULLDUMMY)** separately fixed the `OP_CHECKMULTISIG` extra-element malleability.[^14][^17][^18]

## Weight, virtual bytes & implicit fees

SegWit replaced the 1 MB size cap with a **4,000,000 weight-unit (WU)** cap (fact):[^19][^20]
- **Non-witness (base) bytes count 4 WU each; witness bytes count 1 WU each** — i.e. witness data is discounted 4×. `weight = base_size × 4 + witness_size × 1` (legacy: `size × 4`).
- **Virtual size (vsize) = weight ÷ 4**; **1 vByte = 4 WU**; block limit = **1,000,000 vbytes**.
- **Fees are implicit:** `fee = sum(input values) − sum(output values)` — no explicit fee field; the leftover goes to the miner. **Fee rate is quoted in sat/vByte**; `total fee ≈ vsize × fee_rate`.

## Bitcoin Script

Bitcoin Script is a **stack-based, Forth-like** language (fact): operands and **opcodes** push onto a stack; opcodes pop arguments, compute, and push results.[^21][^22] It is **stateless** (no storage between scripts) and **intentionally not Turing-complete** — **no loops, no recursion**.

**Why non-Turing-complete is a design choice, not a flaw** (fact — disconfirms the "limitation/accident" framing). No loops guarantees bounded, predictable termination, preventing infinite-loop DoS against validators and making spending conditions statically analyzable.[^22][^23] The trade-off (limited contract expressiveness) is real and deliberate — which is why richer-but-riskier systems chose differently.

**Locking vs unlocking scripts** (fact):[^11][^21][^24]
- **scriptPubKey** = locking script on an output; sets spend conditions.
- **scriptSig** (legacy) / **witness** (SegWit) = unlocking script supplying satisfying data.
- **Execution:** the unlocking script runs first, leaving items on the stack; then the locking script runs against that stack. Valid iff execution ends with a single truthy value and no error. *(Historical note: very early Bitcoin concatenated the two; they were split into separate scripts to prevent the unlocking script injecting operators into the locking script.)*[^24]

## Standard output types & addresses

(fact — 3+ sources)[^11][^21][^25][^26]

| Type | scriptPubKey template | Unlock via | Prefix | Encoding |
|---|---|---|---|---|
| **P2PKH** (legacy) | `OP_DUP OP_HASH160 <20-B pkh> OP_EQUALVERIFY OP_CHECKSIG` | `<sig> <pubkey>` | `1…` | Base58Check |
| **P2SH** | `OP_HASH160 <20-B script hash> OP_EQUAL` | inputs + redeemScript | `3…` | Base58Check |
| **P2WPKH** (SegWit v0) | `OP_0 <20-B pkh>` | witness `<sig> <pubkey>` | `bc1q…` (42) | Bech32 |
| **P2WSH** (SegWit v0) | `OP_0 <32-B script hash>` | witness inputs + witnessScript | `bc1q…` (62) | Bech32 |
| **P2TR** (Taproot, v1) | `OP_1 <32-B tweaked x-only pubkey>` | key-path sig **or** script-path | `bc1p…` (62) | Bech32m |

A **P2WPKH** witness is verified as if by a P2PKH script (sig + pubkey).[^11][^24] **P2TR** uses **Schnorr** (x-only) keys and commits to an optional script tree — spendable by the tweaked key *or* by revealing a committed script (see Taproot in `references/mempool-fees-and-upgrades.md`).[^25] *(The signature math itself → `blockchain-crypto-primitives`.)*

## Multisig

Bare multisig is `OP_<m> <pubkey1> … <pubkeyN> OP_<n> OP_CHECKMULTISIG` (m-of-n), in practice wrapped in **P2SH** or **P2WSH** (fact, incl. the well-known bug):[^18][^27] `OP_CHECKMULTISIG` pops **one extra element** beyond what it should, so spenders prepend a **dummy element** (conventionally `OP_0`). The bug is consensus-frozen (fixing it needs a hard fork); **BIP147/NULLDUMMY** now requires that dummy be the empty byte array.

## Timelocks

(fact — BIP65/BIP112)[^15][^28]
- **Transaction-level fields:** **nLockTime** (absolute — earliest height/time the *whole tx* is valid) and **nSequence** (per-input; encodes a *relative* locktime since BIP68).
- **Script opcodes:**
  - **OP_CHECKLOCKTIMEVERIFY (CLTV, BIP65)** — *absolute* UTXO-level timelock; soft-forked Dec 2015. Fails unless the spending tx's **nLockTime ≥** the stack value.
  - **OP_CHECKSEQUENCEVERIFY (CSV, BIP112)** — *relative* timelock; soft-forked May 2016. Fails unless the input's **nSequence** indicates the required relative time/blocks have elapsed since the spent output was mined.

These opcodes underpin Lightning (HTLCs, `to_self_delay`) and vaults.

## Address encodings

(fact — BIPs/Optech)[^25][^26][^30]
- **Base58Check** — legacy `1…`/`3…`: version byte + payload + 4-byte dSHA256 checksum, Base58-encoded (no `0`,`O`,`I`,`l`; case-sensitive).
- **Bech32 (BIP173)** — SegWit **v0** (P2WPKH/P2WSH): lowercase, human-readable part (`bc`/`tb`), separator `1`, witness-version char, BCH checksum (polymod constant `1`).
- **Bech32m (BIP350)** — SegWit **v1+** (**P2TR**): after a flaw in bech32's error-detection for non-zero witness versions, bech32m swaps the checksum constant (reported as `0x2bc830a3` — single-source, treat as detail); otherwise identical. Wallets supporting only bech32 must upgrade to pay Taproot.

## HD wallets (name level)

Modern wallets are **hierarchical-deterministic (HD)**: one seed derives a tree of keys.[^31] **BIP32** (derivation), **BIP39** (mnemonic seed phrase), **BIP44** (5-level path `m/purpose'/coin_type'/account'/change/address_index`), **BIP84** (purpose `84'` → P2WPKH `bc1q…`), **BIP86** (purpose `86'` → P2TR `bc1p…`), BIP49 (`49'` → nested-SegWit). *(The derivation math → `blockchain-crypto-primitives`.)* See `references/nodes-wallets-and-account-model.md` for wallet operation.

## Child sub-concepts (future research)
PSBT (BIP174/370); SIGHASH flags & the signature-commitment model (incl. BIP143/BIP341 sighash); Miniscript & output descriptors; full Taproot internals; the Script opcode table & disabled-opcode/covenant re-enablement debates; the coin-selection waste-metric math.

## References

[^1]: https://river.com/learn/bitcoins-utxo-model/ — blog: UTXO definition, balance = sum of UTXOs, change, coin selection.
[^2]: https://en.wikipedia.org/wiki/Unspent_transaction_output — wiki: UTXO + UTXO-set concept.
[^3]: https://www.kaleido.io/blockchain-blog/utxo-vs-account-model — blog: UTXO vs account, privacy, parallelism, nonce, disadvantages.
[^4]: https://docs.cardano.org/about-cardano/learn/eutxo-explainer — docs (Cardano): eUTXO as response to UTXO's stateless smart-contract limits.
[^7]: https://bitcoinops.org/en/topics/coin-selection/ — Optech: BnB/Knapsack/SRD, waste metric, change avoidance.
[^8]: https://blog.summerofbitcoin.org/coin-selection-part-1/ — blog: BnB by effective value, exact-match/no-change, privacy.
[^9]: https://github.com/bitcoin/bitcoin/blob/master/src/policy/policy.cpp — Core source: dust via dustRelayFee; spend-cost formulas.
[^10]: https://www.lightspark.com/glossary/dust-limit — blog: 546 sat legacy / 294 sat SegWit thresholds, 1/3-of-value rule, 3000 sat/kvB default.
[^11]: https://developer.bitcoin.org/devguide/transactions.html — docs: tx anatomy, outpoint, scriptSig/scriptPubKey, implicit fees, P2PKH/P2SH/multisig.
[^12]: https://learnmeabitcoin.com/technical/transaction/ — reference: field-by-field serialization order and byte sizes.
[^14]: https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki — BIP141: SegWit serialization, marker/flag, txid vs wtxid, witness commitment, malleability fix, 4M-WU limit.
[^15]: https://github.com/bitcoin/bips/blob/master/bip-0112.mediawiki — BIP112: OP_CHECKSEQUENCEVERIFY (CSV) relative-timelock semantics; nSequence per BIP68.
[^16]: https://learnmeabitcoin.com/technical/transaction/wtxid/ — reference: wtxid = dSHA256 of witness serialization; equals txid when no witness.
[^17]: https://en.bitcoin.it/wiki/Transaction_malleability — wiki: malleability sources, low-S/PR#6769, BIP62 withdrawn, SegWit fix.
[^18]: https://en.bitcoin.it/wiki/OP_CHECKMULTISIG — wiki: off-by-one bug, dummy OP_0, hard-fork freeze, BIP147 NULLDUMMY.
[^19]: https://en.bitcoin.it/wiki/Weight_units — wiki: weight = base×4 + witness×1, 4,000,000 WU cap, vsize = weight/4.
[^20]: https://learnmeabitcoin.com/technical/transaction/size/ — reference: bytes vs weight vs vbytes; fee = vbytes × sat/vByte.
[^21]: https://learnmeabitcoin.com/technical/script/ — reference: stack-based, locking/unlocking, execution, opcodes.
[^22]: https://en.bitcoin.it/wiki/Script — wiki: Forth-like, stateless, not Turing-complete, no loops, DoS rationale.
[^23]: https://www.0xkishan.com/blogs/satoshis-ingenious-decision-behind-bitcoins-non-turing-completeness — blog: non-Turing-completeness as deliberate security choice; trade-offs.
[^24]: https://dev.to/tvpeter/bitcoin-transaction-validation-p2pkh-338 — blog: unlocking-then-locking execution, P2PKH validation, P2WPKH "as if P2PKH."
[^25]: https://learnmeabitcoin.com/technical/keys/address/ — reference: Base58Check vs Bech32/Bech32m, prefixes, P2TR Schnorr/tweaked-key.
[^26]: https://bitcoinops.org/en/preparing-for-taproot/ — Optech: bech32 (v0) vs bech32m (Taproot/v1), Schnorr, prefixes.
[^27]: https://github.com/jimmysong/programmingbitcoin/blob/master/ch08.asciidoc — book: bare multisig template, OP_CHECKMULTISIG off-by-one, P2SH wrapping.
[^28]: https://github.com/bitcoin/bips/blob/master/bip-0065.mediawiki — BIP65: OP_CHECKLOCKTIMEVERIFY semantics.
[^30]: https://bitcoinops.org/en/topics/bech32/ — Optech: BIP173 bech32 vs BIP350 bech32m, error-detection weakness, checksum constant difference.
[^31]: https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki — BIP32: HD key derivation (plus BIP39/44/84/86 family at name level).
