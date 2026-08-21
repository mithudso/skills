<!-- Provenance: reference under the standalone `bitcoin-protocol-expert` skill. Created 2026-06-16 via /dr deep-research. Peer deferral: PoW-as-consensus-class → distributed-systems-consensus; SHA-256/Merkle-proof math → blockchain-crypto-primitives; valuation theory → blockchain-economics. -->

# Block Structure, Mining, Difficulty & Monetary Policy

How blocks are built and chained, how Proof-of-Work secures them, and how bitcoin is issued.

## Contents
- [Block & 80-byte header](#block--80-byte-header)
- [Coinbase & maturity](#coinbase--maturity)
- [Block size / weight limit](#block-size--weight-limit)
- [Merkle tree & root](#merkle-tree--root)
- [Proof of Work (Bitcoin-specific)](#proof-of-work-bitcoin-specific)
- [Difficulty adjustment](#difficulty-adjustment)
- [The halving](#the-halving)
- [Monetary policy & the 21M cap](#monetary-policy--the-21m-cap)
- [References](#references)

> **Scope note.** PoW *as a consensus-mechanism class* (Nakamoto consensus, Sybil resistance, 51% theory, fork-choice) → `distributed-systems-consensus`. SHA-256 and Merkle-proof math → `blockchain-crypto-primitives`. Valuation theory (stock-to-flow) → `blockchain-economics`. This file covers Bitcoin's concrete mechanics only.

## Block & 80-byte header

Every block is an **80-byte header** + the transaction list. The header is what gets hashed in mining, so its format is consensus-critical (fact — 4+ sources incl. primary docs):[^1][^2]

| Bytes | Field | Contents |
|---|---|---|
| 4 | **version** | block version / soft-fork signaling bits |
| 32 | **previous block header hash** | dSHA256 of the prior header (the chain link) |
| 32 | **merkle root** | commitment to all txs in the block |
| 4 | **time** | Unix timestamp when the miner started hashing |
| 4 | **nBits** | the target threshold in "compact" encoding |
| 4 | **nonce** | the value miners increment to search for a valid hash |

`4+32+32+4+4+4 = 80`. Hashes are internal byte order; integers little-endian.[^2]
- **nBits / target** — a 4-byte "compact" (base-256 scientific-notation) form of the full 256-bit target a valid block hash must fall below.[^1]
- **nonce** — only 32 bits (~4.3 B values); modern ASICs exhaust it in a fraction of a second, hence the **extranonce** (below).[^4]

## Coinbase & maturity

The first transaction in every block is the **coinbase**, which the miner constructs to pay themselves the **block reward = subsidy + sum of all transaction fees**; total coinbase output must not exceed subsidy + fees (fact).[^6][^7] A coinbase output is **unspendable until 100 blocks** are built on top of it (**coinbase maturity**) — protecting against reorgs spending a reward that later gets orphaned (fact).[^6][^8]

## Block size / weight limit

Pre-SegWit (2017): **1 MB** serialized size. SegWit replaced this with a **4,000,000 weight-unit (WU)** cap; non-witness bytes count 4 WU, witness bytes 1 WU. An all-legacy block still maxes near 1 MB; with witness data the serialized payload can reach ~**4 MB** (practically ~2–2.3 MB given real tx mixes). `vByte = weight ÷ 4`, so 4M WU = 1M vbytes. *(fact for the 4M-WU cap and weighting; qualified for the practical-size figure. SegWit internals → `references/mempool-fees-and-upgrades.md`.)*[^9][^10][^11]

## Merkle tree & root

Transactions are summarized into one 32-byte **merkle root** in the header. Construction (fact):[^12][^1]
1. List TXIDs, coinbase first.
2. Pair them; concatenate each pair (64 raw bytes) and **double-SHA256** to make the next row.
3. **Odd row → duplicate the last hash and pair it with itself.**
4. Repeat to one hash — the merkle root.

Changing or reordering any tx changes the root and invalidates the header's PoW.[^12]

**Role in SPV / light clients** (name level). The root enables **Simplified Payment Verification**: a light client downloads only the **80-byte headers** and verifies inclusion via a **merkle branch/proof** (the sibling hashes to the root) — e.g. ~384 bytes instead of ~75,000 for a full block.[^14][^12] *(Merkle-proof math → `blockchain-crypto-primitives`.)*

## Proof of Work (Bitcoin-specific)

A block is mined by computing **double-SHA256 of the 80-byte header** and checking whether the 256-bit result is **below the current target** (nBits). dSHA256 is uniform over `0…2²⁵⁶−1`, so if the target is fraction *p* of that range, each header hash succeeds with probability *p* (fact).[^4][^2]

**Search loop:** hardware iterates the 32-bit **nonce**; once exhausted, miners change an **extranonce** in the coinbase, which changes the coinbase TXID → the merkle root → a fresh header to search (the outer loop).[^4][^6] **Hashrate** is the aggregate hashes/sec the honest network computes — the economic wall an attacker must out-compute. **Mining pools** combine hashrate and pay participants by **shares** (near-miss hashes below an easier pool target); a real block's reward is split by shares.[^5]

## Difficulty adjustment

Difficulty retargets **every 2016 blocks** (~2 weeks at the 10-minute target: 2016 × 10 min ≈ 14 days). The protocol compares the actual span of the last 2016-block window to two weeks and rescales the target; faster → harder, slower → easier (fact).[^15][^16] A single retarget can change difficulty by **at most 4× (and at least ÷4)** — the 4× clamp.[^15][^16]

**Time-warp note** (fact; exact arithmetic qualified). The retarget uses the span between the **first and last** block of a window, but 2016 blocks have only **2015 intervals** — the code targets `600 × 2016 = 1,209,600 s` instead of `600 × 2015 = 1,209,000 s`, and adjustment periods **don't overlap**.[^16][^17] A majority-hashrate miner can exploit this **time warp attack** by setting most blocks to the minimum allowed timestamp and only the final block to true current time, making periods *appear* long so difficulty repeatedly halves; iterated, block production approaches the Median-Time-Past lower bound (~6 blocks/sec). A proposed **consensus-cleanup** fix requires the first block of a new period to be timestamped no earlier than 10 minutes before the previous period's last block.[^17]

## The halving

The subsidy **halves every 210,000 blocks** (~4 years) (fact):[^18][^19]

| Halving | Height | ~Date | Subsidy after |
|---|---|---|---|
| genesis era | 0 | Jan 2009 | 50 BTC |
| 1st | 210,000 | Nov 28, 2012 | 25 BTC |
| 2nd | 420,000 | Jul 9, 2016 | 12.5 BTC |
| 3rd | 630,000 | May 11, 2020 | 6.25 BTC |
| **4th** | **840,000** | **Apr 19–20, 2024** | **3.125 BTC** |
| 5th (proj.) | 1,050,000 | ~Mar–Apr 2028 | 1.5625 BTC |

Issuance continues until the subsidy rounds to zero around year **~2140**, after which miners are paid **only from transaction fees**.[^18][^20]

## Monetary policy & the 21M cap

Supply is governed by the **geometric subsidy schedule**, which asymptotically approaches ~21M (fact).[^21][^22]

**Disconfirming nuance — true asymptotic supply is *under* 21M** (fact). The subsidy is stored in **satoshis** and halved by **integer right-shift**, truncating fractional satoshis each step. Summing the truncated schedule yields **≈20,999,999.9769 BTC** (`2,099,999,997,690,000 sat`) — slightly less than 21M.[^21][^22] (Bitcoin Core's `MAX_MONEY` is `21,000,000 × COIN`, but *emitted* total falls short.) After ~33 halvings the right-shifted value rounds to 0, ending issuance.[^21][^18]

**Lost/unspendable coins** further reduce circulating supply (e.g., the genesis-block 50 BTC is unspendable by a code quirk; lost keys). Magnitudes are uncertain (qualify).

**Security-budget question (subsidy → 0)** (qualify — actively contested). Honest hashrate is funded largely by the subsidy today; whether **fees alone** will fund enough security after the subsidy vanishes is an open, genuinely contested debate — skeptics warn fees may be insufficient (a March-2025 data point put fees at ~1.25% of the block reward); proponents argue fee markets and price appreciation compensate.[^20][^23] *(Valuation theory → `blockchain-economics`.)*

## Child sub-concepts (future research)
Median-Time-Past (BIP113); the "Great Consensus Cleanup" soft fork; Stratum / Stratum V2 pool protocols; the witness-discount weight-accounting mechanism; nBits↔difficulty conversion (`difficulty_1`); fee-market / block-space-auction dynamics.

## References

[^1]: https://developer.bitcoin.org/reference/block_chain.html — docs: header fields, 80-byte serialization, merkle root, nBits compact format.
[^2]: https://learnmeabitcoin.com/technical/block/ — reference: six header fields, byte sizes/order, header-as-PoW-input.
[^4]: https://developer.bitcoin.org/devguide/mining.html — docs: dSHA256-of-header below target, nonce iteration, uniform distribution.
[^5]: https://bitcoinops.org/en/topics/pooled-mining/ — Optech: pooled mining, shares, extranonce.
[^6]: https://learnmeabitcoin.com/technical/mining/coinbase-transaction/ — reference: coinbase = subsidy + fees; 100-block maturity; extranonce.
[^7]: https://learnmeabitcoin.com/technical/transaction/fee/ — reference: fees as leftover claimed in the block reward.
[^8]: https://developer.bitcoin.org/devguide/block_chain.html — docs: 100-block coinbase maturity rationale (reorg protection).
[^9]: https://en.bitcoin.it/wiki/Block_weight — wiki: 4,000,000 WU cap, 4/1 weighting, ~2.3 MB all-segwit, discount rationale.
[^10]: https://jimmysong.medium.com/understanding-segwit-block-size-fd901b87c9d4 — blog: 1 MB → up-to-4 MB with witness data.
[^11]: https://bitcoinmagazine.com/guides/what-is-the-bitcoin-block-size-limit — blog: block-size-limit history, vbytes.
[^12]: https://learnmeabitcoin.com/technical/block/merkle-root/ — reference: merkle-root construction, odd-node duplication, proof size, light clients.
[^14]: https://github.com/bitcoin-sv/docs/blob/master/important-concepts/high-level/simplified-payment-verification-spv.md — docs: SPV, merkle branch to root, header-only verification.
[^15]: https://en.bitcoin.it/wiki/Difficulty — wiki: 2016-block retarget, 10-min target = 2 weeks, 4× clamp.
[^16]: https://21ideas.org/en/difficulty/ — blog: 2016-vs-2015-intervals off-by-one.
[^17]: https://bitcoinops.org/en/topics/time-warp/ — Optech: time-warp mechanics, non-overlapping periods, MTP, consensus-cleanup fix.
[^18]: https://www.theblock.co/post/289875/ — news: halving heights/dates, 50→25→12.5→6.25→3.125, subsidy→0 ~2140.
[^19]: https://www.coingecko.com/en/coins/bitcoin/bitcoin-halving — blog: halving schedule, 210,000-block interval, next-halving block.
[^20]: https://www.coingecko.com/learn/what-happens-last-bitcoin-mined — blog: last sat ~2140; fees-only thereafter; security-budget framing.
[^21]: https://en.bitcoin.it/wiki/Controlled_supply — wiki: 20,999,999.9769 BTC asymptote, satoshi/right-shift truncation, subsidy→0.
[^22]: https://stacker.news/items/574451 — forum: ~21M, MAX_MONEY math, 2,099,999,997,690,000 sat.
[^23]: https://danhedl.medium.com/bitcoins-security-is-fine-93391d9b61a8 — blog (opinion): counter-argument that fee markets sustain security.
