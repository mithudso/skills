<!-- Provenance: reference under the standalone `bitcoin-protocol-expert` skill. Created 2026-06-16 via /dr deep-research. Overview depth. Volatile claims stamped verified-as-of 2026-06-16. Peer deferral: Schnorr/adaptor-sig math → blockchain-crypto-primitives. -->

# The Lightning Network

Bitcoin's primary Layer-2 payment-channel network — **overview depth** (a working mental model, not an exhaustive spec).

**verified-as-of: 2026-06-16** — volatile claims (BOLT12 status, implementation landscape, capacity figures) stamped inline.

## Contents
- [The problem it solves](#the-problem-it-solves)
- [Payment channels](#payment-channels)
- [State enforcement & revocation](#state-enforcement--revocation)
- [HTLCs & multi-hop routing](#htlcs--multi-hop-routing)
- [Routing](#routing)
- [BOLT specs & invoices](#bolt-specs--invoices)
- [Implementations, watchtowers, LSPs, custody](#implementations-watchtowers-lsps-custody)
- [Honest limitations](#honest-limitations)
- [References](#references)

## The problem it solves

Bitcoin's base layer is deliberately throughput-constrained — ~1 MB/10-min limits it to an estimated **3.3–7 tps** (qualify the exact range) with multi-minute latency.[^1][^2][^15] The **Lightning Network (LN)**, from the 2016 Poon–Dryja whitepaper, moves individual payments **off-chain** into bilateral payment channels, using the base chain only as settlement/dispute court. A channel needs just **two on-chain transactions** (open, close) yet supports effectively unlimited instant low-fee payments between; LN is "block-size neutral" (no consensus change) (fact).[^1][^2][^15]

## Payment channels

A channel opens with a **2-of-2 multisig funding transaction** locking the **channel capacity** (fact).[^4][^5] Before broadcasting the funding tx, parties exchange signed-but-un-broadcast **commitment transactions** encoding the current balance split. Each payment is both sides signing a **new** commitment reflecting updated balances and **invalidating** the prior one — off-chain, instant, trustless.[^4][^5] Each party holds an **asymmetric** version (their own balance delayed/encumbered; the counterparty's immediately spendable).[^7][^9] Closing:[^5][^6][^7]
- **Cooperative (mutual) close** — both co-sign one clean closing tx; outputs immediately spendable, no timelocks.
- **Unilateral / force close** — one party broadcasts their latest commitment alone; their own output is timelock-encumbered, the counterparty's is immediately spendable. Force closes are expensive and lock funds for days — a frequent source of operator losses.

## State enforcement & revocation

Any old commitment a party signed is still technically valid; Bitcoin can't "un-sign" it, so LN makes broadcasting an old state **self-punishing** (fact).[^4][^8] On each advance, each party hands the other the **per-commitment secret** (`revoke_and_ack`) to reconstruct a **revocation key** for the now-old commitment. The encumbered output (BOLT #3) has two paths: a **revocation** branch (counterparty sweeps immediately with the revocation key) and a **`to_self_delay`** branch enforced by **CSV** (a relative timelock, commonly ~144 blocks ≈ 1 day) before the broadcaster can claim.[^5][^7][^9] If a cheater broadcasts a revoked commitment, the honest party has the CSV window to broadcast a **justice transaction (breach-remedy / penalty)** sweeping **all** channel funds — the cheater loses everything. This **LN-Penalty** model requires the honest party (or a watchtower) **online within the timelock**.[^4][^7][^8][^9]

**eltoo / LN-Symmetry** (name level) replaces penalties with "later state supersedes earlier" — publishing an old state costs only fees, simplifying backups and enabling multi-party channels. It needs a soft fork adding **SIGHASH_ANYPREVOUT (APO, BIP-118)**. **verified-as-of 2026-06-16: APO/BIP-118 is proposed-but-not-activated** (exercised on the Bitcoin Inquisition signet, not mainnet); eltoo is **not deployed on mainnet**.[^10][^11][^17]

## HTLCs & multi-hop routing

A single channel only lets two directly-connected peers pay each other. **Hashed Timelock Contracts (HTLCs)** extend this to trustless multi-hop payments (fact).[^12][^13] An HTLC is a conditional output with two paths: a **hashlock** (claimable by revealing the preimage `R` of payment hash `H = hash(R)`) and a **timelock** (refundable to the sender after a CLTV expiry if `R` is never revealed). The receiver generates `R`, shares only `H`; every hop sets up an HTLC locked to the **same `H`**. When the receiver claims by revealing `R`, each upstream hop learns `R` and claims its incoming HTLC — funds move **atomically**, and no intermediary can steal (they only pull funds in after paying forward). HTLC timelocks **decrease toward the destination** (`cltv_expiry_delta`), guaranteeing each hop time to settle on-chain if forced. On-chain HTLC scripts (HTLC-Success/Timeout, BOLT #3) are the backstop.[^5][^12][^13]

**PTLCs (Point Timelock Contracts)** (name level) — enabled by **Schnorr/Taproot** (live since Nov 2021), they replace the shared hash with a distinct EC **point** per hop via adaptor signatures, eliminating the payment-correlation leak (and wormhole vector) of HTLCs. Well-understood but **pending Lightning-layer spec/implementation work as of 2026-06-16**; no soft fork needed (Taproot suffices). *(Adaptor-signature math → `blockchain-crypto-primitives`.)*[^18][^19]

## Routing

LN uses **source-based routing**: the **sender** computes the entire path (intermediate nodes don't pick the next hop). Payments are wrapped in **onion routing** (the **Sphinx** mix-format, BOLT #4) — a fixed-size packet where each hop decrypts only its own layer, learning the previous and next hop but not the full path/origin/destination.[^14][^20] Nodes build the channel graph from the **gossip protocol** (BOLT #7): `channel_announcement` (proves an on-chain channel), `channel_update` (fees, `cltv_expiry_delta`, HTLC min/max), `node_announcement`. Only nodes with a public channel may gossip (anti-spam).[^14][^21][^22]

**Disconfirming — LN does NOT guarantee a route exists or succeeds** (fact). Gossip broadcasts channel **capacity but not the live balance split** (withheld for privacy/scalability), so the sender must **guess** which channels have liquidity on the correct side ("liquidity uncertainty"). A wrong guess fails the payment, forcing a retry; failure probability compounds with path length and payment size. This is a fundamental, much-discussed limitation, not an edge case.[^14][^20][^23][^25]

## BOLT specs & invoices

The protocol is the **BOLT** series — **Basis of Lightning Technology** (github.com/lightning/bolts) (fact).[^27][^28] Per BOLT #0's index: **#1** Base Protocol; **#2** Peer Protocol (Channel Management); **#3** Bitcoin Transaction & Script Formats; **#4** Onion Routing; **#5** On-chain Transaction Handling; **#7** Node & Channel Discovery (gossip); **#8** Encrypted Transport (Noise); **#9** Feature Flags; **#10** DNS Bootstrap; **#11** Invoice Protocol; **#12** Offers. *(No BOLT #6 — merged into #2.)*

**Invoices — BOLT11 vs BOLT12** (verified-as-of 2026-06-16):[^29][^30][^31]
- **BOLT11** — universal, bech32-encoded, QR-ready, **single-use** invoice (payee pubkey + amount + payment hash); the baseline every wallet supports.
- **BOLT12 "offers"** (Rusty Russell, **officially merged Sept 2024** — the first new BOLT since 2017) — a **static, reusable** payment code; the payer fetches a fresh invoice over the network via **onion messages**, with **blinded paths** for receiver privacy, plus refunds, recurring/subscription payments, and non-BTC denomination. **Additive, not a replacement** — BOLT11 remains the fallback and most users still transact with BOLT11.

## Implementations, watchtowers, LSPs, custody

Four major implementations follow the same BOLTs (verified-as-of 2026-06-16):[^32][^33][^41]
- **LND** (Lightning Labs, Go) — largest node share (commonly cited ~70% — qualify), best tooling; native BOLT12 via the LNDK sidecar.
- **Core Lightning / CLN** (Blockstream, C) — spec-leading, plugin architecture; **first to ship BOLT12**; most mature splicing.
- **Eclair** (ACINQ, Scala) — built for high-volume routing nodes; full BOLT12.
- **LDK** (Spiral/Block, Rust) — a **library** (not a node) for embedding LN into apps; full BOLT12.

**Watchtowers** watch the chain for revoked-commitment broadcasts and submit the justice tx while a client is offline (LND ships a private "altruist" tower; clients send encrypted blobs). Incentivized watchtowers remain largely unshipped.[^33][^34] **Lightning Service Providers (LSPs)** and the **LSPS0/1/2** spec abstract channel/liquidity management for end users.[^35] These create a **custodial ↔ non-custodial spectrum**: fully **custodial** wallets (e.g. Wallet of Satoshi) hold keys and see all metadata (convenient, reintroduces trust/censorship/hack risk); **LSP-assisted non-custodial** (e.g. Phoenix) — LSP sees payments not keys; **self-hosted nodes** — trustless in normal operation but operationally demanding.[^35][^36][^37]

## Honest limitations

- **Liquidity / inbound-capacity** — a freshly opened channel has all funds local; you can send but **cannot receive** until someone opens toward you or routes through. Fragmentation and **channel depletion** (funds pile one-sided) are persistent and partly-unsolved; mitigations (Loop, splicing, liquidity ads, LSP marketplaces) add cost.[^23][^25][^26][^38]
- **Routing reliability** — no guaranteed route (above). One industry source (Amboss, Feb 2025) claimed a ~$100 payment fails ~23% of the time; a 2025 preprint argued real topology can't reliably move >~100k sat between arbitrary nodes — **single-source figures, treat as indicative, not settled**; the qualitative problem is high-confidence.[^36]
- **On-chain fee dependence** — open/close/force-close cost on-chain fees and can lock funds for the CSV window.[^16]
- **Watchtower / always-online requirement** — the penalty model needs you or a tower online within the timelock; towers are themselves a liveness dependency.[^16][^40]
- **Channel-jamming / griefing** — a cheap DoS routing payments through a target and stalling, tying up HTLC **slots** (default max 483/channel) and liquidity; "essentially free" (failed payments pay no fees), **long-standing and unfixed in deployment**; proposed mitigations (unconditional fees + reputation; HTLC-GP) are studied, not standardized. Related: **replacement-cycling attacks** (Riard, CVE-2023-40231..40234).[^12][^42][^43][^44]
- **Privacy caveats** — improved over base-layer but **not anonymous**: the funding UTXO is public forever, BOLT11 exposes the receiver pubkey, balance probing can trace payments; BOLT12/PTLCs help but wallet support is uneven; custodial wallets see everything.[^36]
- **Centralization pressure** — open/close fees favor well-connected hubs; node/channel counts and capacity declined from their 2022/2025 peaks (specific percentages are single-source — qualify).[^36][^37][^39]

**Adoption context** (verified-as-of 2026-06-16, single-author current-state source — figures volatile): real usage exists (~$1.17B monthly volume Nov 2025; Coinbase/Binance/Kraken integrations; El Salvador Chivo; Nostr zaps; Taproot Assets bringing stablecoins to LN rails), but node decline, UX friction, and channelless competitors (Spark, Ark) are headwinds.[^36]

## Child sub-concepts (future research)
Splicing & channel factories; channelless L2 alternatives (Spark, Ark, statechains/VTXOs); Taproot Assets / multi-asset Lightning; channel-jamming mitigations & replacement-cycling security; the LSP spec & JIT channels; trampoline & Pickhardt-Richter pathfinding.

## References

[^1]: https://lightning.network/lightning-network-paper.pdf — paper (LN whitepaper): channels, commitment txs, breach-remedy penalty.
[^2]: https://en.wikipedia.org/wiki/Bitcoin_scalability_problem — encyclopedia: base-layer 3.3–7 tps, LN as L2.
[^4]: https://www.spark.money/research/payment-channels-concept-to-implementation — blog: channel lifecycle, justice tx, watchtower rationale.
[^5]: https://github.com/lightning/bolts/blob/master/05-onchain.md — spec (BOLT #5): cross-signed commitments, revocation seizes funds, on-chain resolution.
[^6]: https://www.spark.money/research/payment-channels-concept-to-implementation — blog: cooperative vs force close.
[^7]: https://bitclawd.com/learn/lightning/specs/bolt-03/ — blog (BOLT #3 walkthrough): revoke_and_ack, asymmetric commitments, OP_IF revocation + to_self_delay/CSV, HTLC-Success/Timeout.
[^8]: https://github.com/lnbook/lnbook/blob/develop/09_channel_operation.asciidoc — book: penalty/revocation construction, revoke_and_ack.
[^9]: https://developers.bitgo.com/docs/lightning-overview — docs: funding tx, commitments, revocation-key invalidation, closure modes.
[^10]: https://covenants.info/proposals/apo/ — wiki: eltoo/LN-Symmetry, no-penalty model, requires SIGHASH_ANYPREVOUT.
[^11]: https://github.com/bitcoin/bips/blob/master/bip-0118.mediawiki — spec: SIGHASH_NOINPUT→ANYPREVOUT (BIP-118).
[^12]: https://docs.lightning.engineering/the-lightning-network/multihop-payments/hash-time-lock-contract-htlc — docs (LND): HTLC anatomy, multi-hop atomicity, no-steal, 483-slot limit, cltv_expiry_delta.
[^13]: https://docs.lightning.engineering/the-lightning-network/multihop-payments — docs: atomic conditional payments, off-chain settlement.
[^14]: https://github.com/lightning/bolts/blob/master/04-onion-routing.md — spec (BOLT #4): Sphinx, fixed-size hop payloads, per-hop ECDH.
[^15]: https://www.clevelandfed.org/-/media/project/clevelandfedtenant/clevelandfedsite/publications/working-papers/2022/wp-2219-the-lightning-network.pdf — paper (Cleveland Fed): ~7 tps, two-tx open/close netting.
[^16]: https://www.spark.money/research/lightning-network-2026-state — blog (current state): force-close losses, LSP layer, watchtower context.
[^17]: https://petertodd.org/2024/covenant-dependent-layer-2-review — blog: APO exercised on Inquisition signet only; LN-Symmetry needs the (unactivated) soft fork.
[^18]: https://www.spark.money/research/ptlcs-next-evolution-lightning — blog: PTLCs via adaptor sigs, decorrelation, Taproot live since Nov 2021, APO not yet activated.
[^19]: https://github.com/BlockstreamResearch/scriptless-scripts/blob/master/md/multi-hop-locks.md — spec-adjacent: PTLC multi-hop locks, payment decorrelation.
[^20]: https://github.com/lnbook/lnbook/blob/develop/10_onion_routing.asciidoc — book: Sphinx (Danezis & Goldberg 2009), source-based onion routing.
[^21]: https://docs.lightning.engineering/the-lightning-network/the-gossip-network — docs: gossip announcements, anti-spam public-channel requirement.
[^22]: https://github.com/lnbook/lnbook/blob/develop/11_gossip_channel_graph.asciidoc — book: channel_announcement/update, fees & cltv_expiry_delta.
[^23]: https://www.spark.money/research/lightning-network-routing-deep-dive — blog: sender constructs full route, balance invisibility, hub concern.
[^25]: https://github.com/lnbook/lnbook/blob/develop/12_path_finding.asciidoc — book: liquidity uncertainty, balances unannounced, exponential failure with path length.
[^26]: https://protos.com/why-two-party-bitcoin-lightning-channels-keep-failing/ — blog: channel depletion (Pickhardt), capacity/channel decline.
[^27]: https://github.com/lightning/bolts/blob/master/00-introduction.md — spec (BOLT #0 index): full BOLT series + "Basis of Lightning Technology."
[^28]: https://github.com/lightning/bolts — spec (repo): canonical BOLT home.
[^29]: https://www.spark.money/research/lightning-invoices-bolt11-bolt12 — blog: BOLT12 merged Sept 2024, offer model, per-impl status, coexist.
[^30]: https://github.com/lightning/bolts/blob/master/12-offer-encoding.md — spec (BOLT #12): offers as invoice precursor, onion-message delivery, non-LN currency.
[^31]: https://bitcoinblog.de/2025/04/23/how-bolt12-can-enhance-the-lightning-user-experience/ — blog: offers vs single-use invoices, early wallet support.
[^32]: https://www.spark.money/tools/lightning-node-comparison — blog: LND/CLN/LDK/Eclair language/architecture/target/splicing.
[^33]: https://github.com/lightningnetwork/lnd/blob/master/docs/watchtower.md — docs: altruist vs reward watchtower, encrypted justice-tx blobs.
[^34]: https://docs.corelightning.org/docs/watchtowers — docs: watchtower role, always-online assumption.
[^35]: https://github.com/BitcoinAndLightningLayerSpecs/lsp — spec: the LSP specification (LSPS0/1/2 — transport, channel ordering, just-in-time channels).
[^36]: https://gajumaru.io/2026/04/02/why-the-bitcoin-lightning-network-failed/ — blog (skeptical, single-author): privacy caveats, custodial table, 23%/100k-sat failure claims, competitors. (Used for the limitation framing; figures flagged.)
[^37]: https://lightning.news/the-custodial-lightning-model/ — blog (advocacy): custodial risks, "10% of nodes / 80% liquidity" (cross-checked, flagged).
[^38]: https://blog.muun.com/the-inbound-capacity-problem-in-the-lightning-network/ — blog: inbound-capacity bootstrapping problem.
[^39]: https://thecurrencyanalytics.com/ — blog: capacity/channel decline figures (volatile, single-source).
[^40]: https://singulargrit.substack.com/p/lightnings-velvet-manacles-watchtowers — blog (opinion): watchtower liveness-dependency critique (framing only).
[^41]: https://github.com/FulgurVentures/Lightning-Implementation-Features — feature matrix: per-implementation feature support.
[^42]: https://eprint.iacr.org/2022/1454 — paper (Chaincode): jamming = cheap DoS, "longstanding and yet unfixed."
[^43]: https://research.chaincode.com/2022/11/15/unjamming-lightning/ — research blog: slot vs liquidity jamming, mitigation proposals.
[^44]: https://lists.linuxfoundation.org/pipermail/lightning-dev/2023-October/021999.html — mailing list: replacement-cycling CVE-2023-40231..40234.
