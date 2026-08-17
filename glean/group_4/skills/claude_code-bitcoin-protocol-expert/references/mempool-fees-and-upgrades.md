<!-- Provenance: reference under the standalone `bitcoin-protocol-expert` skill. Created 2026-06-16 via /dr deep-research. Volatile claims stamped verified-as-of 2026-06-16. Peer deferral: Schnorr/ECDSA math → blockchain-crypto-primitives; soft-fork-as-governance theory referenced not re-derived. -->

# Mempool, Fees & Protocol Upgrades

The fee market and fee-bumping mechanics, plus the SegWit and Taproot soft forks and how soft forks activate.

**verified-as-of: 2026-06-16** — all volatile claims (RBF default state, TRUC/package-relay status, min-relay-fee, adoption %) are stamped inline.

## Contents
- [The mempool](#the-mempool)
- [Minimum relay fee](#minimum-relay-fee)
- [Fee estimation, vbytes & the witness discount](#fee-estimation-vbytes--the-witness-discount)
- [Block-template selection](#block-template-selection)
- [Replace-By-Fee (RBF → full RBF)](#replace-by-fee-rbf--full-rbf)
- [TRUC / v3 & package relay](#truc--v3--package-relay)
- [Child-Pays-For-Parent (CPFP)](#child-pays-for-parent-cpfp)
- [SegWit](#segwit)
- [Taproot](#taproot)
- [Soft-fork mechanics & activation](#soft-fork-mechanics--activation)
- [References](#references)

## The mempool

The mempool is each node's local store of valid-but-unconfirmed transactions — **per-node, not global**; there is no canonical network-wide mempool, and operators run different policies, so mempools diverge (fact).[^1][^2][^3] Bitcoin Core's default size limit is **300 MB** (`-maxmempool`); on overflow the node **evicts the lowest-*descendant-feerate* package first** (descendant scoring avoids evicting a low-fee parent with a high-fee CPFP child) and raises its **dynamic `mempoolminfee`** to the last-evicted feerate (decaying back over time). Unconfirmed txs **expire after 336 h (14 days)** (`-mempoolexpiry`).[^1][^3][^4]

## Minimum relay fee

Nodes reject and don't forward txs below the **minimum relay feerate** (`-minrelaytxfee`) — a DoS floor, not a fee paid to the node.[^1] **verified-as-of 2026-06-16:** the historic default was **1 sat/vB (1,000 sat/kvB)** for over a decade; **Bitcoin Core v29.1 (Sept 4, 2025, PR #33106)** lowered the default `-minrelaytxfee` and `-incrementalrelayfee` to **100 sat/kvB (0.1 sat/vB)** (and `-blockmintxfee` to 1 sat/kvB), backported to 28.x/29.x. The release note warns lower-fee txs aren't guaranteed to propagate unless the lower defaults are widely adopted; wallet feerates were unchanged.[^5][^6]

## Fee estimation, vbytes & the witness discount

**Weight units & vbytes (BIP141).** Block limit = **4,000,000 WU**; each non-witness byte = 4 WU, each witness byte = 1 WU (witness discounted 4×). **vsize = weight ÷ 4** = `non-witness bytes + 0.25 × witness bytes`. *Worked example:* 226 raw bytes (116 non-witness + 110 witness) = 116×4 + 110×1 = **574 WU = 143.5 vbytes**.[^7][^8][^9] *Why discount?* Witness data needn't be retained by all nodes long-term and (arguably) costs less; the discount also incentivizes UTXO consolidation. BIP141 chose a single composite weight limit so miners face a linear knapsack.[^7]

**Fee rate is sat/vByte.** `estimatesmartfee` uses BIP141 vsize (witness discounted).[^10] Bitcoin Core's `CBlockPolicyEstimator` is **backward-looking**: it records feerate buckets and how many blocks later they confirmed, exponentially decaying old data, and returns the lowest bucket meeting a success threshold. Modes: **CONSERVATIVE** (95%, longer horizon, stickier) vs **ECONOMICAL** (85%, more responsive).[^10][^12] **verified-as-of 2026-06-16:** Bitcoin Core **v28.0 changed the default mode from `conservative` to `economical`**, reducing over-estimation when RBF is available. A newer mempool-based forecaster was in development (qualify).[^13][^12]

## Block-template selection

Miners maximize **fees per weight**. Since PR #7600 (2016), Bitcoin Core selects by **ancestor feerate** = (tx + all unconfirmed ancestors' fees) ÷ (combined vsize): greedily pick the highest-ancestor-feerate package, add it and its unconfirmed ancestors, recompute descendant feerates, repeat until the **4,000,000 WU** limit (fact).[^14][^15][^16] Calculation is cheap (inherited ancestor totals; chains bounded at depth 25).

## Replace-By-Fee (RBF → full RBF)

**BIP125 opt-in RBF** (Bitcoin Core 0.12.0, 2016): a tx signals replaceability if **any input has nSequence < 0xfffffffe**; a replacement must pay a higher absolute fee *and* higher feerate, and pay its own bandwidth.[^17][^18][^19]

**verified-as-of 2026-06-16 — full RBF is now standard and non-optional** (fact — primary release notes):[^17][^21][^22]
- **v24.0 (2022)** — `-mempoolfullrbf` option added, default `false`.
- **v28.0 (Oct 2, 2024)** — default flipped to `mempoolfullrbf=1`.
- **v29.0 (Apr 14, 2025)** — the option was **removed entirely** (PR #30592): full RBF is the standard behavior, BIP125 signaling no longer required.

**Controversy** (negation). A decade-plus debate: **0-conf merchants** rely on "first-seen" behavior that full RBF undermines; opponents accused Core of "forcing" RBF (Bitrefill cited "fewer than 1 in a million" 0-conf double-spends — single interested source, flag); proponents argued 0-conf was already insecure. Optech notes this is a **mempool-policy** change with **no chain-split risk** (unlike a consensus change). A separate open critique: feerate-only replacement enables **transaction pinning**, motivating TRUC.[^17][^23][^24][^25]

## TRUC / v3 & package relay

**TRUC (Topologically Restricted Until Confirmation) / "v3 policy" (BIP431):** opt in with **nVersion=3**. Rules: a TRUC tx **always signals replaceability**; its unconfirmed ancestry must all be TRUC; **1-parent-1-child topology** (cluster ≤ 2); size cap 10,000 vB (child ≤ 1,000 vB); may be **below min-relay feerate (even 0 fee)** if evaluated in a qualifying package.[^27][^28][^29] *Rationale:* cluster-size-2 makes ancestor score an accurate incentive measure, hardening RBF/CPFP and letting contracting protocols (Lightning) presign 0-fee txs bumped by one anchor — mitigating pinning.

**verified-as-of 2026-06-16 — status:** TRUC shipped in **Bitcoin Core v28.0 (Oct 2024)**; **limited package RBF** (clusters ≤ 2) also v28.0. **Opportunistic 1-parent-1-child (1p1c) package relay** shipped v28.0 (`submitpackage` RPC; opportunistically pairs a too-low-feerate parent with its child); release note warns it is **limited and not yet reliable under adversarial conditions**. Hardening followed (PRs #31397, #31829). **Full package relay (BIP331) was NOT merged into mainline as of 2026** (qualify).[^21][^30]

## Child-Pays-For-Parent (CPFP)

CPFP lets a **high-fee child pull a low-fee parent into a block**: miners evaluate **packages**, so a high-feerate child raises `(parent.fee + child.fee) ÷ (parent.vsize + child.vsize)`, lifting the bundle. The child spends the parent's change (sender) or a received output (recipient, e.g. an exchange bumping a withdrawal). The parent must confirm for the child to be valid — exactly the dependency that forces inclusion (fact).[^14][^31][^32] **Limits:** `-limitancestorcount/size` and `-limitdescendantcount/size` default to 25 / 101,000 vB; once a parent has 25 unconfirmed descendants, no further CPFP — and an attacker can attach a large low-feerate descendant to **pin** a victim (the classic Lightning force-close vector TRUC + package RBF address).[^32]

## SegWit

(BIP141/143/144, activated Aug 2017.) Segregated Witness moves the **witness** into a structure committed separately from the tx Merkle tree, so it is **not part of the txid preimage** (fact):[^7]
- **BIP141** — witness commitment; block limit redefined as **weight ≤ 4,000,000 WU**; the witness discount.
- **BIP143** — new sighash for v0 witness programs, fixing O(n²) signature hashing.
- **BIP144** — P2P relay/serialization for witnesses (marker byte; `MSG_WITNESS_*`; wtxid).

**Fixes third-party malleability** (txid no longer changes when signing changes), **enabling Lightning** (which signs txs spending not-yet-broadcast parents).[^7] Backward-compatible **soft fork**: old nodes enforce only the witness-stripped ~1 MB; effective capacity rose to ~1.6–2.3 MB.[^7][^8] *Activation (block-size-war context):* the BIP9 deployment stalled on withheld miner signaling; **BIP148 (UASF)** threatened a flag-day rejecting non-signaling blocks, and **BIP91 (MASF, 80%/336-block on bit 4)** locked in at **block 477,120 (July 22, 2017)**; SegWit activated at **block 481,824 (Aug 24, 2017)**.[^34][^36][^37]

## Taproot

(BIPs 340/341/342, activated **block 709,632, Nov 14, 2021**.) The largest upgrade since SegWit; defines **SegWit v1** outputs (fact):[^38][^39][^40]
- **BIP340 — Schnorr signatures** for secp256k1: **linearity** enables **key aggregation / MuSig** (N keys → one key, N sigs → one sig, indistinguishable from single-party); 32-byte x-only keys; key-prefixing. *(The EC math → `blockchain-crypto-primitives`.)*
- **BIP341 — Taproot:** the **P2TR** output (`bc1p…`) commits to a tweaked key encoding **two spend paths** — **key-path** (one Schnorr sig; cheapest/most private; never reveals whether a script path existed) and **script-path** (reveals one leaf from a **MAST** with a Merkle proof; other branches stay hidden).
- **BIP342 — Tapscript:** adds **OP_CHECKSIGADD**, Schnorr sig checks, **disables OP_CHECKMULTISIG(VERIFY)**, reserves **OP_SUCCESSx** for future soft forks.

**Privacy:** a multisig treasury, timelocked contract, or Lightning close using the key-path looks like an ordinary single-key payment; multisig is also cheaper (one aggregated sig).[^40][^42] *Activation: Speedy Trial* (BIP9 with a lowered 90% threshold + short window); threshold met **block 687,284 (June 12, 2021)**; activated block 709,632.[^43][^45]

**Adoption** (negation; verified-as-of 2026-06-16): slow and volatile — ~1–2% in year one, a peak above ~40% in early 2024 driven by Ordinals/BRC-20/Runes, then settling to ~13–20% (Glassnode reported ~13.1% adoption / 15.1% utilization on June 13, 2026 — single source, flag). Critics (e.g. Jimmy Song) note it is mostly inscription-driven, not organic, with poor multisig UX. By contrast **SegWit adoption is mature (~85–90% of txs spend a SegWit input as of early 2026** — single strong aggregator, qualify).[^46][^47][^48][^49]

## Soft-fork mechanics & activation

**Soft fork** *tightens* rules so new-rule blocks stay valid to old nodes (backward-compatible; SegWit and Taproot are soft forks). **Hard fork** *loosens/changes* rules so old nodes reject new blocks (must upgrade or split).[^7] Activation methods (fact):[^51][^52][^53][^54]
- **Flag day / IsSuperMajority** — early hard-coded height/version; one proposal at a time.
- **BIP9 version bits** — `nVersion` as a bit field (bits 0–28, ~29 concurrent forks); **95% over a 2016-block period** → LOCKED_IN → ACTIVE after one more period; failure frees the bit. *Weakness:* a miner minority can stall a broadly-supported upgrade.
- **BIP148 (UASF)** — economic nodes enforce a flag-day rejecting non-signaling blocks, removing the miner veto.
- **BIP8** — BIP9 by block height + a **LockInOnTimeout (LOT)** flag: **LOT=true** forces lock-in before timeout; **LOT=false** ≈ BIP9. The LOT debate is about whether miners get a veto.
- **Speedy Trial** (used for Taproot) — BIP9 with a 90% threshold + short window; signals → activate after delay; else expires harmlessly.

**Governance/contention.** Activation choice is Bitcoin's governance flashpoint: the 2015–17 block-size war and the SegWit UASF/BIP91 standoff demonstrated that **economic full nodes — not just miners — hold ultimate authority**. Taproot's Speedy Trial was engineered to avoid re-litigating it. The `mempoolfullrbf` fight is the *policy-layer* analog (contentious, but no chain-split risk).[^45][^36][^51][^23]

## Child sub-concepts (future research)
Transaction pinning & anti-pinning; cluster mempool; ephemeral dust & Pay-to-Anchor (P2A); MuSig2 / FROST threshold signatures; the full block-size war history (Segwit2x, BCH split); the mempool-based fee forecaster.

## References

[^1]: https://bitcoinops.org/en/blog/waiting-for-confirmation/ — Optech: per-node mempool, min-relay-fee DoS floor, descendant-score eviction, fee market.
[^2]: https://github.com/bitcoin/bitcoin/blob/master/src/txmempool.h — Core source: acceptance/eviction, rolling minimum fee.
[^3]: https://www.spark.money/research/bitcoin-mempool-congestion-economics — blog: 300 MB limit, descendant-score eviction, 14-day expiry, v29.1 minrelaytxfee 1→0.1 sat/vB.
[^4]: https://www.spark.money/tools/bitcoin-mempool-visualization-guide — blog: eviction + dynamic-minfee, descendant feerate/CPFP, mempoolexpiry.
[^5]: https://bitcoincore.org/en/releases/29.1/ — release notes: minrelaytxfee & incrementalrelayfee → 100 sat/kvB, PR #33106, propagation warning.
[^6]: https://github.com/bitcoin/bitcoin/commit/6da5de58cabc4133c379baa50845e30e5bc6b3e4 — source diff: DEFAULT_MIN_RELAY_TX_FEE/DEFAULT_INCREMENTAL_RELAY_FEE 1000→100; backported.
[^7]: https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki — BIP141: witness commitment, 4M-WU limit, witness discount, malleability fix; (BIP143/144 companions).
[^8]: https://en.bitcoin.it/wiki/Block_weight — wiki: weight units, 1/4 witness weighting, ~2.3 MB all-segwit, discount rationale.
[^9]: https://learnmeabitcoin.com/technical/transaction/size/ — reference: weight/vsize formulas, 574 WU / 143.5 vB example.
[^10]: https://developer.bitcoin.org/reference/rpc/estimatesmartfee.html — docs: BIP141 vsize, CONSERVATIVE vs ECONOMICAL.
[^12]: https://bitcoinindex.net/blog/how-bitcoin-core-s-fee-estimator-works/ — blog: CBlockPolicyEstimator buckets, exponential decay, MempoolForecaster.
[^13]: https://bitcoincore.org/en/releases/28.0/ — release notes: estimatesmartfee default conservative→economical; mempoolfullrbf default→1; TRUC; limited package RBF; 1p1c relay.
[^14]: https://bitcoin.stackexchange.com/questions/122216/ — expert answer: ancestor-set-feerate block-template selection.
[^15]: https://www.spark.money/glossary/block-template — blog: ancestor-feerate selection, PR #7600, 4M-WU limit.
[^16]: https://bitcoin.stackexchange.com/questions/60344/ — SE: CPFP ancestor-feerate selection, depth-25 cap.
[^17]: https://bitcoinops.org/en/topics/replace-by-fee/ — Optech: opt-in vs full RBF; mempoolfullrbf timeline (v24/v28/v29); pinning critique.
[^18]: https://github.com/bitcoin/bitcoin/blob/master/doc/policy/mempool-replacements.md — Core docs: RBF rules; nSequence<0xfffffffe and nVersion=3 signaling.
[^19]: https://github.com/bitcoin/bips/blob/master/bip-0125.mediawiki — BIP125: opt-in RBF signaling, 0.12.0 rules.
[^21]: https://bitcoincore.org/en/releases/28.0/ — release notes (full RBF default in v28.0; TRUC; 1p1c relay).
[^22]: https://bitcoincore.org/en/releases/29.0/ — release notes: removal of -mempoolfullrbf in v29.0; signaling no longer required.
[^23]: https://bitcoin.stackexchange.com/questions/116347/ — SE/Murch: mempoolfullrbf debate, 0-conf opposition, "no chain-split risk."
[^24]: https://protos.com/bitcoin-core-devs-accused-of-forcing-replace-by-fee-transactions/ — blog: "forcing RBF" accusation.
[^25]: https://github.com/glozow/bitcoin-notes/blob/full-rbf/full-rbf.md — Core dev notes: full-RBF arguments, Bitrefill "<1 in a million."
[^27]: https://github.com/bitcoin/bips/blob/master/bip-0431.mediawiki — BIP431: TRUC rules (nVersion=3, all-TRUC ancestry, 1p1c, ancestor-score, sub-min-fee-in-package).
[^28]: https://bips.dev/431/ — BIP431 mirror: rationale, cluster-size-2, package-relay dependency.
[^29]: https://github.com/bitcoin/bitcoin/blob/master/src/policy/truc_policy.h — Core source: TRUC ancestor/descendant limits = 2, version=3.
[^30]: https://bitcoinops.org/en/topics/package-relay/ — Optech: v28.0 1p1c, submitpackage, opportunistic/limited, hardening PRs, BIP331 not merged.
[^31]: https://learn.txid.uk/en/ideas/cpfp-child-pays-for-parent/ — blog: CPFP package mechanics, parent-must-confirm dependency.
[^32]: https://satoshibench.com/learn/cpfp/ — blog: ancestor-feerate "leaderboard," 25/101kvB limits, pinning/Lightning vector.
[^34]: https://gist.github.com/ajtowns/1c5e3b8bdead01124c04c45f01c817bc — blog: segwit activation history (lockin/active heights).
[^36]: https://www.coindesk.com/markets/2017/07/18/ — news: BIP91 mechanics (bit 4, 80%, 336-block window).
[^37]: https://news.bitcoin.com/bip91-activates-while-fork-still-looms-in-the-backdrop/ — news: BIP91 activated block 477,120, July 22 2017.
[^38]: https://en.bitcoin.it/wiki/Taproot — wiki: Taproot = BIPs 340/341/342; block 709,632 Nov 14 2021; slow-adoption caveats.
[^39]: https://github.com/bitcoin/bips/blob/master/bip-0341.mediawiki — BIP341: key-path hides script path; Schnorr key aggregation.
[^40]: https://www.spark.money/research/taproot-schnorr-signatures-explained — blog: BIP table; Tapscript OP_CHECKSIGADD, OP_CHECKMULTISIG disabled; bc1p; Speedy Trial.
[^42]: https://wiki.21ideas.org/en/concepts/taproot — blog: key-path vs script-path, privacy table, OP_SUCCESS.
[^43]: https://btc.network/learn/what-is-taproot — blog: P2TR key/script-path, MAST example, Speedy Trial block 709,632.
[^45]: https://www.spark.money/tools/bitcoin-soft-fork-history — blog: Speedy Trial = BIP9 + 90% + short window; threshold met block 687,284; LOT explanation.
[^46]: https://www.spark.money/tools/bitcoin-segwit-adoption-tracker — blog: SegWit ~85–90% (2026); Taproot peak >40% early-2024, settled 15–20%.
[^47]: https://phemex.com/news/article/bitcoin-taproot-usage-drops/45802 — blog: Taproot usage 42%→20% (Glassnode).
[^48]: https://studio.glassnode.com/charts/transactions.TaprootAdoption — Glassnode: Taproot 13.1% adoption / 15.1% utilization June 13 2026.
[^49]: https://cointelegraph.com/news/taproot-devs-overlooked-trolling-value-added-bitcoin — blog: Jimmy Song critique (multisig UX, social attack surface).
[^51]: https://bitcoinops.org/en/topics/soft-fork-activation/ — Optech: BIP9/BIP8/LOT/UASF/Speedy Trial; governance framing.
[^52]: https://bitcoin.stackexchange.com/questions/100490/ — SE/Murch: BIP8 vs BIP9, LOT=true forces lock-in.
[^53]: https://github.com/bitcoin/bips/blob/master/bip-0009.mediawiki — BIP9: version-bits semantics, 2016-block retarget.
[^54]: https://bitcoincore.org/en/2016/06/08/version-bits-miners-faq/ — bitcoincore.org: BIP9 95% threshold, LOCKED_IN→ACTIVE.
