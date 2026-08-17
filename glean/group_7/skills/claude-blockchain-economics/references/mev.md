<!-- Provenance: reference under the `blockchain-economics` standalone skill. Created 2026-06-16 via /dr deep-research. Synthesized from web sources; all load-bearing/volatile claims carry [^n] citations to the ## References section. -->

# MEV — Maximal Extractable Value

verified-as-of: 2026-06-16 (volatile: ePBS/EIP-7732 scheduling, relay/builder/solver concentration, Shutter+Primev deployment, FOCIL timing — re-verify before relying)

Value from controlling transaction *ordering/inclusion*: its taxonomy, the supply chain that captures it, the proposer-builder-separation response, and the mitigation menu. Confidence tags: **(fact)** = multi-source primary; **(qualified)** = solid but contested/design-stage; **(tentative)** = single-source/fast-moving. AMM/DEX internals defer to the planned `defi` skill; the base-fee mechanism is `references/fee-markets.md`. **MEV is intrinsic to permissionless ordering and can only be redistributed/minimized/relocated, not eliminated.**

## Contents

- What MEV is, why it exists, Flash Boys 2.0
- MEV taxonomy
- The MEV supply chain
- Proposer-Builder Separation (PBS) & enshrined PBS
- Flashbots, MEV-Boost, relay censorship, SUAVE
- Order-flow auctions (OFAs)
- Intents & solver networks (CoW, UniswapX)
- MEV mitigation
- References

## What MEV is, why it exists, Flash Boys 2.0

**Definition.** "The maximum value that can be extracted from block production in excess of the standard block reward and gas fees by including, excluding, and changing the order of transactions in a block."[^1] Flashbots: "the total value extractable permissionlessly from the re-ordering, inclusion or censoring of transactions within a block."[^2] **(fact)**

**Maximal vs Miner.** Originally *Miner* Extractable Value (Daian et al. 2019); generalized to *Maximal* once Ethereum moved to PoS (proposers, not miners, order blocks) and once searchers/builders/wallets — not just the block producer — captured most of it.[^3][^1] **(fact)**

**Why it exists.** The block producer has unilateral control over inclusion and ordering within their slot; smart-contract state (esp. DeFi prices) is order-dependent, so ordering power is monetizable. A key insight from Daian et al.: this is a *concrete difference* between a payment system (Bitcoin — ordering generally unprofitable to manipulate) and a smart-contract system (Ethereum), so "analyses of Bitcoin miner economics fail to extend to Ethereum."[^1][^3] **(fact)**

**Flash Boys 2.0** (Daian et al., arXiv:1904.05234; IEEE S&P 2020) named MEV; documented arbitrage bots behaving like HFTs; formalized **Priority Gas Auctions (PGAs)** (bots bidding up gas in a continuous-time partial-information game); and identified the consensus threat — when MEV exceeds block rewards it enables **time-bandit attacks** (reorging past blocks to steal historical MEV).[^3] **(fact)**

**Scale (volatile, lower bounds).** Flashbots MEV-Explore ≥$314M from Jan 2020 (early 2021); >$720M in 2021; cumulative >$1.3B (Shutter/academic); ~99% of extraction is arbitrage by tx count.[^2][^11][^16] **(qualified — methodology-dependent lower bounds.)**

## MEV taxonomy

Two atomic primitives underlie everything: **front-running** (insert *before* a target) and **back-running** (insert *after*).[^6][^7]

| Strategy | Primitive | Mechanism | Atomic? | Toxicity |
|---|---|---|---|---|
| **DEX arbitrage** | Back-run | Exploit cross-pool price differences | Atomic | Benign (price discovery) |
| **CEX-DEX arbitrage** | Back-run | CEX↔DEX price gap | Non-atomic (inventory/execution risk) | Benign-ish |
| **Cross-chain arbitrage** | Back-run | Cross-chain price gaps | Non-atomic | Benign-ish |
| **Liquidations** | Back-run | Seize collateral / claim bonus on unhealthy positions | Atomic | Benign (protocol needs them) but gas wars |
| **Sandwich** | Front + Back | Buy before victim's swap (push price up), let victim fill worse, sell after | Atomic | **Toxic** |
| **JIT liquidity** | Front + Back | Add concentrated LP right before a large swap, remove right after | Atomic | Harmful to passive LPs; *helps* the swapper |
| **Long-tail MEV** | Various | Rare/event-based (NFT mints, fraud-proof frontrunning, bridge exploits) | Mixed | Mixed |
| **Time-bandit** | Reorg | Re-mine past blocks for historical MEV | Non-atomic | Toxic / consensus-destabilizing |

[^5][^7][^1][^16] **(fact)**

- **Sandwich** (canonical toxic) — the searcher profits from the victim's slippage; Torres et al. found $18.41M in front-running profits across 11M+ blocks.[^5][^7]
- **Liquidations** — broadly beneficial (protocols require timely liquidations) but the competition produces gas wars that bid away most of the bonus.[^4][^5]
- **JIT liquidity** — captures a large swap's fees while avoiding impermanent loss; gives the *swapper* better execution but steals yield from passive LPs.[^4][^6]
- **Toxic vs benign** — working convention: toxic = front-run/sandwich (value *diverted from* users); benign = back-run/arbitrage/liquidation (price discovery, cross-venue sync). Contested at the margins: a "benign" backrun becomes a frontrun if reordered; one view is toxicity ultimately stems from *validator reordering* being possible at all.[^8][^6] **(qualified — definitional.)**

## The MEV supply chain

Post-Merge pipeline, five roles, value flowing left→right via bids: **users → searchers → builders → relays → proposers (validators).**[^9][^10] **(fact)**

- **Searchers** run bots that detect opportunities and assemble **bundles** (ordered tx sets capturing profit atomically), submitting to builders with a bid — often a direct ETH transfer to `block.coinbase` conditioned on success (don't pay for failed bids), bypassing the public mempool for pre-trade privacy.[^10][^9]
- **Builders** ingest order flow (public txs + searcher bundles + private order flow), run continuous simulations to construct the most profitable block, and bid to buy the proposer's blockspace via relays. Builder edge comes from (1) bundle flow, (2) **private order flow** access, (3) vertical integration of their own searching.[^10][^14]
- **Relays** are "doubly-trusted" intermediaries — trusted by builders for fair routing and by proposers for validity/DA; they aggregate builder bids, validate them, and forward only the highest *valid* bid.[^17][^18]
- **Proposers (validators)** run **MEV-Boost** as a sidecar; it aggregates bids across relays and serves the single highest bid; the validator signs a *blinded* header (can't see contents).[^9][^18]

Payment standard: the builder sets itself as `feeRecipient`, then includes an end-of-block tx paying the proposer; block value = the fee-recipient balance delta.[^14] **(fact)**

**Concentration [mid-2026, volatile]:** top-2 builders (Beaverbuild ~46–49%, Titan ~41–46%) ≈87–95% of blocks; ~7–10 active relays; MEV-Boost mediates ~88–93% of all Ethereum blocks.[^19][^24] **(fact, volatile)**

## Proposer-Builder Separation (PBS) & enshrined PBS

**Rationale.** Profitable block-building needs sophisticated MEV infra most validators can't run; without PBS, MEV would centralize *staking*. PBS decouples block *building* (a few specialized builders) from *proposal* (many validators on consumer hardware) — the proposer just picks the highest bid, distributing the *reward* across all validators even if *extraction* is centralized. ethereum.org: "decentralization is a means, not an end — what we want are honest blocks" (prover-verifier asymmetry — centralized *production* is acceptable if a decentralized set can *verify*).[^20][^4] **(qualified — design philosophy; critics dispute that centralized building is acceptable.)**

**Off-chain PBS via MEV-Boost.** Monnot's "Why enshrine PBS?" is the canonical case: MEV-Boost provides permissionless builder-market access but "relies heavily on a small set of centralized relays acting as mutually-trusted auctioneers."[^21] **(fact)**

**Builder/relay centralization concern (the central critique):**
- **[mid-2026]** ~7–10 relays carry nearly all volume (the 2023 baseline was 6 relays / 5 entities = 99%).[^19][^24][^21] Relays are "antithetical to Ethereum's core values," unprofitable, and at shutdown risk; Flashbots' Dec 2024 migration to **BuilderNet** is "an admission that the relay architecture creates concentration risk."[^21][^25] **(fact / qualified)**
- The primary driver of builder concentration is **private/exclusive order flow** — a self-reinforcing moat.[^9][^14] **(qualified)**

**Enshrined PBS (ePBS / EIP-7732).** Brings the proposer-builder hand-off *into* consensus, eliminating trusted relays (authored by Potuz et al.).[^22] Mechanics: splits each block into a **consensus/beacon half** (broadcast at slot start) and an **execution-payload half** (revealed separately by the builder); builders become **staked on-chain entities** enabling *in-protocol, trustless* builder→proposer payment ("trust-free fair exchange"); adds a **Payload Timeliness Committee (PTC)** attesting on-time delivery, plus dual-deadline fork-choice (also widening the data-propagation window ~2s→~9s, a prerequisite for higher gas/blobs).[^22][^23][^26] **(fact)**

**ePBS status [mid-2026, volatile]:** EIP-7732 is the **consensus-layer headliner for Glamsterdam** (the hard fork after Fusaka), targeted **Q3/H2 2026** (Q4 a documented downside risk after an April 2026 slip), paired with EIP-7928 Block-Level Access Lists. **FOCIL** (Fork-Choice Enforced Inclusion Lists — the censorship-resistance mechanism) was **deferred out of Glamsterdam to the Hegota fork** (late-2026/early-2027) to avoid untested interactions with ePBS at mainnet scale.[^26][^27][^28] **(fact, volatile — verify on Forkcast/ACD notes.)**

## Flashbots, MEV-Boost, relay censorship, SUAVE

- **MEV-Geth (Jan 2021)** — Flashbots' first product: a forked Geth + relay forming a private tx pool + sealed-bid blockspace auction; searchers sent bundles to the relay, which forwarded valid bundles to *miners*, who included the most profitable bundle top-of-block (paid via `block.coinbase`).[^14][^29] **(fact)**
- **MEV-Boost (PoS)** — open-source middleware between the validator's consensus and execution clients; queries relays, profit-switches, falls back to local building if relays disconnect. The blinded header eliminates validator frontrunning; a commit-reveal scheme puts solo stakers and large pools on equal footing — but the validator must trust the relay to filter invalid payloads, report value accurately, and reveal contents after signing.[^30][^9] **(fact)**
- **Relay trust failure (concrete):** the Nov 10 2022 Flashbots relay vulnerability — a malicious builder sent bids with incorrect `timestamp`/`prev_randao`, causing ~350 blocks to fall back to local building (no slots missed); fixed by adding validation. Part of what pushed the community toward ePBS.[^31][^21] **(fact)**
- **Censorship / OFAC-compliant relays:** after Tornado Cash sanctions (Aug 2022), some relays filtered sanctioned txs; OFAC-compliant blocks spiked to ~51% (Oct 2022), with zero Tornado txs in 19,436 Flashbots-relay blocks. Censorship later declined (~66% by Dec 2022) as non-censoring relays gained share; Flashbots' "Cost of Resilience" proposed a `min-bid` parameter (~0.1% APR cost). FOCIL/inclusion lists are the structural fix (deferred to Hegota).[^32][^33][^34][^27] **(fact, historically volatile)**
- **SUAVE** (Single Unifying Auction for Value Expression) — Flashbots' decentralization endgame: an independent chain acting as a plug-and-play decentralized mempool + builder for any chain (MEVM + Kettles + Confidential Compute Requests). **(qualified — alpha-stage; status uncertain post-BuilderNet pivot.)**[^29][^35][^36]

## Order-flow auctions (OFAs)

OFAs auction the *right to execute (or backrun) an order* and return the winning bid to the user — explicitly analogized to **Payment For Order Flow (PFOF)** in TradFi.[^12][^15] Design space (Frontier Research): orders as **intents** vs **transactions**; settlement rights public/permissioned/private; mechanism families include exclusive batch auctions, transaction-execution services (share order info, compete on user cashback — "internalizes" MEV but carries non-execution risk), and blockspace aggregators.[^12] **MEV Blocker** (CoW): a builder-level OFA RPC providing privacy + frontrun protection, rebating users **up to ~90%** of backrun value.[^13] **MEV-Share** (Flashbots) — the redistribution OFA (see Mitigation). **(fact / qualified)**

## Intents & solver networks (CoW, UniswapX)

**Intent model.** Users sign a *desired outcome* ("≥X of B for Y of A") rather than an executable tx; competing solvers/fillers find optimal execution, **internalizing MEV** as price improvement.[^13][^38] **(fact)**

**CoW Protocol (CoW Swap)** — a dApp-level OFA using **batch auctions**: signed intents → off-chain orderbook → a **Fair Combinatorial Batch Auction** where **solvers** propose settlements in a fixed window (unfair batched bids filtered), winner settles on-chain (second-price rewards, CIP-67). MEV defenses are structural: **Uniform Directed Clearing Prices (UDCP)** make tx order irrelevant within the batch (undermining reordering MEV), and **Coincidence of Wants (CoW)** matches opposite-side users peer-to-peer without AMM liquidity. Solvers enforce EBBO ("Ethereum Best Bid and Offer").[^13][^39][^40] **(fact)**

**UniswapX** (whitepaper; launched Jul 2023) — a non-custodial **Dutch-auction** intent protocol: swappers sign off-chain orders (via **Permit2**) with a max price; broadcast to permissionless **Fillers** (searchers/MMs); a Dutch auction decays from high until a filler executes on-chain. Gas-free swaps (no gas on failed swaps), MEV internalized as price improvement, orders backstopped by Uniswap's Smart Order Router; `Reactor` contracts implement order types (incl. `ExclusiveDutchOrder`); extensible cross-chain; later versions use a two-stage RFQ.[^38][^41][^48] **(fact)**

## MEV mitigation

**Encrypted mempools / threshold encryption (Shutter).** Hide tx contents until the order is fixed → kills *content-dependent* MEV (frontrun/sandwich). **Shutter** uses **threshold encryption** (a decryption key split among *n* "Keypers"; ≥*t* must cooperate; live since Oct 2022 on Gnosis Chain); other techniques are delay encryption (time-lock) and TEEs (SGX). **[late-2025/2026]** Shutter + Primev announced the first threshold-encrypted mempool for Ethereum's out-of-protocol PBS: builders commit *blind* to encrypted txs backed by *slashable stake*; Keypers release keys only after ≥3 commitments (conditional decryption defeats decrypt-then-censor).[^11][^45][^46] **(fact, fast-moving)**

> **Negation — encrypted mempools are not a silver bullet:**
> - **Metadata leakage** (a16z): the encrypted payload is surrounded by *unencrypted* metadata (sender, nonce, gas, size, timing, IP) enabling **probabilistic/speculative MEV** — searchers don't decrypt, they guess; hiding metadata costs bandwidth/latency.[^42] **(fact)**
> - **Doesn't stop** censorship (proposer can exclude encrypted txs), latency games, or post-decryption visibility; goal is narrow ("keep intent hidden until order is fixed").[^44] **(qualified)**
> - **Reveal optionality** (LUCID critique): the sender can choose *not* to reveal — fundamental "as long as we lack sufficiently good cryptographic tools"; LUCID "does not prevent frontrunning — it only makes it harder."[^43] **(qualified)**
> - **Threshold-trust + liveness risk** (FRP-31; Charbonneau): *inconspicuous collusion* among decryptors is hard to detect; a high t/n threshold means an outage of half the committee can **halt** the chain (dangerous for liveness-favoring Gasper); encryption can preclude efficient orderings (timely liquidations) and incentivize spam.[^47][^49] **(fact)**
> - **Economic-vs-cryptographic enforcement** (Shutter+Primev): with economic enforcement, frontrunning is still possible *if MEV exceeds the builder's stake*; only on-chain cryptographic placement-validation prevents it unconditionally, but that needs consensus changes + higher gas.[^46] **(fact)**

**FCFS (first-come-first-served).** Several L2s order by arrival time, not priority fee, eliminating gas-auction frontrunning at the sequencer (IEEE survey arXiv:2407.19572 lists Arbitrum, Optimism, Starknet, zkSync, Scroll). Caveat: FCFS converts fee competition into a **latency race** (co-location advantage) and is spam-vulnerable; a centralized sequencer is itself a trust point.[^45a] **(fact / qualified)**

**Private order flow / private RPC.** Sending txs to a private RPC (Flashbots Protect, MEV Blocker) prevents frontrunning by hiding the tx until inclusion — but this is the **same mechanism that drives builder centralization** (private flow is the structural moat).[^13b][^9] **(qualified — mitigation for the user, centralization risk for the network.)**

**MEV redistribution — MEV-Share (Flashbots).** An open-source protocol: users send orderflow with **programmable privacy**; a **matchmaker** inserts users' private txs into searchers' partial bundles, forwarding matched bundles with a **validity condition** requiring an MEV payment back to the user — credibly neutral, "does not enshrine any individual builder."[^50][^51] **(fact)**

**MEV smoothing & MEV burn (consensus-layer redistribution).** Both target the *variance/spike* problem (MEV lotteries destabilize consensus and incentivize DoS'ing the lucky proposer):
- **MEV smoothing** (Monnot/Hitzig 2021) — reduce per-validator MEV variance toward a stake-proportional distribution; hard because MEV detection/attribution is difficult (a builder could bribe a pool to accept sub-optimal payments).[^52a] **(qualified)**
- **MEV burn** (Drake et al., an ePBS add-on) — attesters impose the median observed base-fee as a *floor* and **burn** it (EIP-1559-style), redistributing value to *all* holders rather than the lucky proposer.[^52][^53] **(qualified)**

> **Negation — limits & "MEV cannot be eliminated":**
> - **Builder collusion** breaks MEV-burn — colluding builders don't bid in the tipping window and split MEV among themselves (identical to individually-optimal private deals), suppressing the burn floor; a dishonest majority sets it to zero.[^52] Lotteries shrink but persist (simulations: 19 blocks still paid >0 ETH burn-tips).[^54] **(qualified)**
> - **Core principle:** "As long as validators choose which transactions to include and in what order, MEV exists." REV (Realized Extractable Value) ≤ MEV; MEV is an asymptotic bound. Mitigation **redistributes/minimizes/relocates** MEV (and minimizing finality time shrinks reorg/time-bandit MEV), never eliminates it; a16z: encrypted mempools "are unlikely to be the silver bullet."[^6][^21a][^42] **(fact)**
> - **Intent/solver centralization:** internalizing MEV *relocates* centralization "from the blockchain layer to the solver layer." Chitra et al. (arXiv:2403.02525): UniswapX top-3 solvers ≈79%; fillers collapsed 2,000+ → 12 (Jan 2024); role is *permissioned*; entry+effort costs → oligopoly (and a limited-solver market can be *welfare-superior* under congestion). Other measurements: on UniswapX two MMs >90%; 1inch Fusion top-3 ≈80%; CoW less concentrated.[^48a][^48] **(qualified — measurements vary by date/protocol; all show high concentration.)**

## Child concepts (future research)

1. Cross-domain / cross-chain & L2 MEV (shared sequencers, based rollups, cross-chain arbitrage atomicity).
2. Preconfirmations & the proposer-commitment economy (mev-commit/Primev, based-preconfs, execution tickets, APS).
3. MEV detection & attribution tooling (mev-inspect, EigenPhi, Zeromev, REV-vs-MEV methodology, bundle-toxicity classifiers).
4. Decentralized block-building (BuilderNet, TEE-based building SGX/TDX, searcher-builder vertical-integration economics).
5. Inclusion lists & censorship resistance (FOCIL/EIP-7805, multiplicity, the OFAC/neutrality policy layer — deferred to Hegota).

## References

[^1]: Maximal Extractable Value (MEV) — https://ethereum.org/developers/docs/mev/ — docs (canonical) — definition, maximal-vs-miner, examples, long-tail.
[^2]: Quantifying MEV / MEV-Explore v0 — Flashbots — https://writings.flashbots.net/quantifying-mev — primary — definition + ≥$314M cumulative lower bound.
[^3]: Flash Boys 2.0 (Daian et al.) — arXiv:1904.05234 / IEEE S&P 2020 — https://arxiv.org/abs/1904.05234 — foundational peer-reviewed — origin of "MEV," PGAs, time-bandit attacks, consensus-instability thesis.
[^4]: Maximal Extractable Value — DF3NDR — https://docs.df3ndr.com/Book/3/11/3-mev.html — security reference — arbitrage/liquidation/sandwich/backrun/JIT definitions + benign/toxic.
[^5]: MEV in DeFi: Taxonomy, Detection, Mitigation (SoK) — arXiv:2411.03327 — https://arxiv.org/html/2411.03327 — academic SoK — sandwich/liquidation/JIT taxonomy; value-diverting MEV.
[^6]: SoK: MEV Countermeasures — arXiv:2212.05111 — https://ar5iv.labs.arxiv.org/html/2212.05111 — academic SoK — front/back-run primitives; toxic-vs-benign; "as long as validators order, MEV exists."
[^7]: SoK: The Evolution of MEV, From Miners to Cross-Chain — arXiv:2603.07716 — https://arxiv.org/pdf/2603.07716 — academic SoK — strategy×primitive table (CEX-DEX/cross-chain non-atomic); empirical profits.
[^8]: On the toxicity classification of MEV transactions — Flashbots Collective — https://collective.flashbots.net/t/on-the-toxicity-classification-of-mev-transactions/521 — forum (primary) — toxic (front/sandwich) vs non-toxic (back/liquidation) + reordering caveat.
[^9]: MEV Supply Chain / PBS guide — https://cryptothreads.io/learn/mev-supply-chain-explained/ — explainer — five-role pipeline; private-orderflow-as-moat; concentration figures.
[^10]: Flashbots Auction Overview — https://docs.flashbots.net/flashbots-auction/overview — docs (primary) — searcher/builder/relay/validator roles; bundles; coinbase-transfer bidding.
[^11]: Shutter Network: Private Transactions from Threshold Cryptography — eprint 2024/1981 — https://eprint.iacr.org/2024/1981 — peer-reviewed — threshold-encryption protocol; >$1.3B cumulative MEV; production since Oct 2022.
[^12]: The Orderflow Auction Design Space — Frontier Research — https://frontier.tech/the-orderflow-auction-design-space — research firm — OFA = PFOF analogy; exclusive-batch / execution-service / blockspace-aggregator.
[^13]: Understanding Order Flow Auctions — CoW Protocol — https://blog.cow.fi/understanding-order-flow-auctions-5532d4ef2f40 — primary — OFA levels; MEV Blocker ~90% rebate; CoW intent flow.
[^13b]: Private Order Flow and OFA Mechanisms — Stengarl — https://stengarl.eth.limo/posts/mev-private-order-flow-and-ofa-mechanisms/ — independent — private-RPC taxonomy; front-run-protection-vs-concentration table.
[^14]: Block Builders (Flashbots Docs) + Unbundling the MEV Supply Chain Pt 1 (hack.vc) — https://docs.flashbots.net/flashbots-mev-boost/block-builders · https://blog.hack.vc/unbundling-mev-supplychain-part-1-history-and-evolution/ — docs/research — builder payment standard; MEV-Geth history; 3 builder-edge factors.
[^15]: Order Flow Auctions Explained — CoW DAO — https://cow.fi/learn/order-flow-auctions-explained — primary — OFA conceptual overview.
[^16]: An Introduction to MEV on Ethereum — EY — https://www.ey.com/content/dam/ey-unified-site/ey-com/en-us/insights/blockchain/documents/ey-an-introduction-to-maximal-extractable-value-on-ethereum.pdf — industry whitepaper — four MEV types incl. long-tail; ~99% arbitrage; >$720M in 2021.
[^17]: Relay Fundamentals — Flashbots Docs — https://docs.flashbots.net/flashbots-mev-boost/relay — docs — relays as "doubly-trusted" DA layer.
[^18]: flashbots/mev-boost (GitHub) — https://github.com/flashbots/mev-boost — spec/code — MEV-Boost = PBS middleware; relay-of-relays; profit-switching; local fallback.
[^19]: Rated Network PBS Landscape + mevboost.pics + relayscan.io — https://explorer.rated.network/builders · https://mevboost.pics/ · https://www.relayscan.io/ — live dashboards — [mid-2026] Beaverbuild ~46–49% / Titan ~41–46%; top-3 92–98%; ~7–10 relays; ~92% MEV-Boost share.
[^20]: Proposer-builder separation — ethereum.org — https://ethereum.org/roadmap/pbs/ — canonical — PBS rationale; prover-verifier asymmetry; "decentralization is a means."
[^21]: Why enshrine PBS? — Monnot, ethresear.ch — https://ethresear.ch/t/why-enshrine-proposer-builder-separation-a-viable-path-to-epbs/15710 — forum (primary) — MEV-Boost ≈90% of blocks; relays centralized (6/5 = 99%); OFAC censorship; case for ePBS.
[^21a]: Quantifying Realized Extractable Value — Flashbots — https://writings.flashbots.net/quantifying-rev — primary — REV ≤ MEV; MEV as asymptotic bound.
[^22]: EIP-7732: Enshrined Proposer-Builder Separation — https://eips.ethereum.org/EIPS/eip-7732 — docs (official EIP) — staked builders; trust-free fair exchange; CL-vs-EL payment enforcement.
[^23]: ethereum/epbs-security-analysis — https://github.com/ethereum/epbs-security-analysis — repo — two-phase block; PTC; removal of trusted relays.
[^24]: Rated Network Relay Landscape (30d) — https://explorer.rated.network/relays — live — [mid-2026] relay market shares.
[^25]: Glamsterdam Slips: ePBS Runs Late — BlockEden — https://blockeden.xyz/blog/2026/04/18/glamsterdam-epbs-delay-ethereum-mev-reform-slip/ — analysis — [2026] BuilderNet migration = concentration-risk admission; FOCIL → Hegota.
[^26]: Glamsterdam Hardfork Tracker — ETH Daily + ethereum.org — https://ethdaily.io/glamsterdam · https://ethereum.org/roadmap/glamsterdam/ — tracker — [mid-2026] ePBS + BALs headliners; PTC; H2 2026; FOCIL deferred.
[^27]: Ethereum Sets Q3 2026 Glamsterdam Timeline — Crypto Times — https://www.cryptotimes.io/2026/06/05/ethereum-glamsterdam-q3-2026/ — news — [Jun 2026] Q3 2026 target; >88% blocks via MEV-Boost; FOCIL → Hegota.
[^28]: EIP-7732 (ePBS) Selected as Glamsterdam Headliner — EtherWorld — https://etherworld.co/eip-7732-epbs-selected-as-glamsterdam-headliner/ — news — [Aug 2025] ACDC #162 SFI status; FOCIL = CFI.
[^29]: The Future of MEV is SUAVE — Flashbots — https://writings.flashbots.net/the-future-of-mev-is-suave — primary — MEV-Geth/relay history; commit-reveal trust relaxation; SUAVE vision.
[^30]: MEV-Boost: Merge-ready Flashbots Architecture — ethresear.ch — https://ethresear.ch/t/mev-boost-merge-ready-flashbots-architecture/11177 — design (primary) — middleware; blinded-header; relay-trust requirements.
[^31]: Post-mortem: relay vulnerability (Nov 10 2022) — Flashbots Collective — https://collective.flashbots.net/t/post-mortem-for-a-relay-vulnerability-affecting-flashbots-and-other-relays-impacting-3-of-mev-boost-relayed-validators/727 — primary — ~350 local-fallback blocks; timestamp/prev_randao/gas_limit fixes.
[^32]: MEV Watch + The Cost of Resilience — Flashbots — https://www.mevwatch.info/ · https://writings.flashbots.net/the-cost-of-resilience — primary — OFAC-compliant-block tracking; min-bid resilience at ~0.1% APR.
[^33]: At least 23% of blocks comply with US sanctions — The Block — https://www.theblock.co/post/173417/almost-half-of-all-ethereum-blocks-comply-with-us-sanctions — press — 0 Tornado txs in 19,436 Flashbots-relay blocks; ~51% peak.
[^34]: Is Ethereum's Censorship Problem Taking a Turn? — CoinDesk — https://www.coindesk.com/tech/2022/12/21/is-ethereums-censorship-problem-taking-a-turn — press — censorship decline to ~66% via relay diversity.
[^35]: What is SUAVE? — suave-alpha.flashbots.net — https://suave-alpha.flashbots.net/what-is-suave — primary — decentralized off-chain compute; MEVM, Kettles, Confidential Compute Requests.
[^36]: flashbots/suave-specs — https://github.com/flashbots/suave-specs — spec — SUAVE = platform for OFAs/builders; Centauri/Andromeda tracks.
[^38]: UniswapX Whitepaper — https://app.uniswap.org/whitepaper-uniswapx.pdf — protocol (primary) — Dutch-auction intents; Permit2; MEV-internalization; gas-free; cross-chain.
[^39]: CoW Fair Combinatorial Batch Auction + Coincidence of Wants — https://docs.cow.fi/cow-protocol/concepts/introduction/fair-combinatorial-auction — docs — UDCP makes tx-order irrelevant; CoW peer-to-peer matching.
[^40]: CoW Auction Mechanism + Solver Competition Rules — https://docs.cow.fi/cow-protocol/reference/core/auctions — docs — solver bids; unfair-bid filtering; second-price rewards (CIP-67); EBBO.
[^41]: Introducing UniswapX + Uniswap/UniswapX (GitHub) — https://blog.uniswap.org/uniswapx-protocol · https://github.com/Uniswap/UniswapX — primary — filler network; Smart-Order-Router backstop; Reactor/DutchOrder/ExclusiveDutchOrder.
[^42]: On the limits of encrypted mempools — a16z crypto — https://a16zcrypto.com/posts/article/limits-encrypted-mempools/ — research — NEGATION: metadata-driven probabilistic MEV; "not the silver bullet."
[^43]: A Criticism of LUCID and Encryption-Scheme-Agnostic Encrypted Mempools — ethresear.ch — https://ethresear.ch/t/a-criticism-of-lucid-and-encryption-scheme-agnostic-encrypted-mempools/25210 — forum — NEGATION: reveal-optionality fundamental; "does not prevent frontrunning, only makes it harder."
[^44]: Encrypted Mempools: Security Beyond Encryption — ZK/SEC — https://blog.zksecurity.xyz/posts/encrypted-mempool-security/ — analysis — NEGATION: "reduce MEV, do not eliminate it"; no censorship/latency/metadata protection.
[^45]: The Road Towards an Encrypted Mempool on Ethereum — Shutter docs — https://docs.shutter.network/docs/shutter/research/the_road_towards_an_encrypted_mempool_on_ethereum — docs — threshold-encryption mechanics; PBS integration.
[^45a]: MEV Mitigation Approaches in Ethereum and L2 (survey) — arXiv:2407.19572 — https://arxiv.org/html/2407.19572v1 — academic survey — FCFS table (Arbitrum/Optimism/Starknet/zkSync/Scroll); threshold/delay/TEE.
[^46]: Threshold Encrypted Mempools with mev-commit + Shutter blog — https://ethresear.ch/t/threshold-encrypted-mempools-with-mev-commit-preconfirmations/23588 · https://blog.shutter.network/the-first-encrypted-mempool-is-coming-to-pbs-on-ethereum/ — primary — [late-2025] Shutter+Primev; conditional decryption; economic-vs-cryptographic enforcement (NEGATION).
[^47]: FRP-31: Threshold Encrypted Mempools — Limitations — Flashbots Collective — https://collective.flashbots.net/t/frp-31-threshold-encrypted-mempools-limitations-and-considerations/2036 — primary — NEGATION: inconspicuous decryptor collusion; residual censorship; spam incentives.
[^48]: Execution Welfare Across Solver-based DEXes — arXiv:2503.00738 — https://arxiv.org/html/2503.00738v2 — academic — UniswapX RFQ; NEGATION: UniswapX 2 MMs >90%; 1inch top-3 ~80%; CoW less concentrated.
[^48a]: An Analysis of Intent-Based Markets (Chitra/Kulkarni et al.) — arXiv:2403.02525 — https://arxiv.org/pdf/2403.02525 — academic — NEGATION: UniswapX top-3 ≈79%; 2,000→12 fillers; permissioned; entry/effort → oligopoly.
[^49]: Encrypted Mempools — Jon Charbonneau — https://joncharbonneau.substack.com/p/encrypted-mempools — analysis — NEGATION: t/n liveness-halt risk (Gasper); metadata leakage; UX latency/cost.
[^50]: MEV-Share: programmably private orderflow — Flashbots Collective — https://collective.flashbots.net/t/mev-share-programmably-private-orderflow-to-share-mev-with-users/1264 — primary — matchmaker + validity conditions; MEV payment back to user; no enshrined builder.
[^51]: flashbots/mev-share (GitHub) — https://github.com/flashbots/mev-share — spec — open protocol for OFAs/redistribution.
[^52]: MEV burn — a simple design (Drake) — ethresear.ch — https://ethresear.ch/t/mev-burn-a-simple-design/15590 — forum (primary) — ePBS add-on; base-fee floor + burn; NEGATION: colluding-builders failure mode.
[^52a]: Committee-driven MEV smoothing (Monnot/Hitzig) — ethresear.ch — https://ethresear.ch/t/committee-driven-mev-smoothing/10408 — forum — variance reduction; detection/attribution difficulty; bribe-to-accept-suboptimal attack.
[^53]: The price is right: predictive MEV-burn — ethresear.ch — https://ethresear.ch/t/the-price-is-right-predictive-mev-burn/18656 — forum — MEV-burn borrows EIP-1559; pMEV-burn variant.
[^54]: In a post-MEV-Burn world — simulations (Wahrstätter) — ethresear.ch — https://ethresear.ch/t/in-a-post-mev-burn-world-some-simulations-and-stats/17092 — forum — NEGATION/qualifier: reduces inequality + DoS incentive, but lotteries persist (19 blocks >0 ETH burn-tip).
