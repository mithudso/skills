<!-- hub-reference-banner -->
> **Reference file — part of the `blockchain` hub.** Formerly the standalone `blockchain-economics` skill.
> Sibling topics in this family are now reference files under the hubs (`blockchain`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: blockchain-economics
description: "Blockchain economics & cryptoeconomics — the incentive/mechanism layer above protocol mechanics. TRIGGER: tokenomics (issuance/emission schedules, inflation/disinflation, burns/sinks, ultrasound-money debate, vesting & unlock cliffs, fair-launch vs VC, airdrops, utility/governance/security-token taxonomy + Howey, the token-velocity / MV=PQ problem); cryptoeconomics & mechanism design (incentive-compatibility, the security budget & 51%/cost-of-attack economics, staking economics — real yield/slashing/liquid staking/restaking, PoS-vs-PoW economic security, Schelling points, the verifier's dilemma); fee-market design (EIP-1559 base fee & why it's not a pure auction, priority fees, EIP-4844 blob fees, Solana local fee markets, first-vs-second-price TFM); MEV (arbitrage/sandwich/liquidation/JIT, the searcher→builder→relay→proposer supply chain, PBS & enshrined-PBS, Flashbots/MEV-Boost, order-flow auctions, intents/solvers, encrypted mempools); governance economics (token-weighted/quadratic/conviction voting, vote-buying & bribe markets, DAO treasuries). SKIP: DeFi protocol mechanics/AMM math → defi; consensus theory/safety proofs → distributed-systems-consensus; per-chain protocol detail → bitcoin-protocol-expert / ethereum-protocol-expert; crypto-primitive math → blockchain-crypto-primitives; crypto investing/trading → investing-and-retirement / trading-and-investing (mechanism economics, NOT trading advice)."
version: 1.0.0
updated: 2026-06-16
category: blockchain
whenToUse: Use when reasoning about the economics and incentive design of a blockchain or token system — token supply/issuance and value capture, the economic security budget and cost-of-attack, staking and restaking yield, transaction-fee mechanism design (EIP-1559, blobs, Solana local markets), MEV and its supply chain / mitigation, or governance and DAO-treasury economics. This is the mechanism/incentive layer; defer protocol mechanics, consensus theory, and trading advice to the peer skills named in SKIP.
keywords:
  - tokenomics
  - cryptoeconomics
  - MEV
  - maximal extractable value
  - EIP-1559
  - fee market
  - proposer-builder separation
  - staking economics
  - security budget
  - token velocity
  - vesting unlock
  - governance token
  - DAO treasury
  - quadratic voting
  - liquid staking
  - restaking
  - slashing
  - EIP-4844 blob fee
  - Flashbots MEV-Boost
  - bribe market
tags:
  - blockchain
  - blockchain-economics
  - cryptoeconomics
  - tokenomics
  - mev
  - defi-adjacent
  - mechanism-design
---

# Blockchain Economics & Cryptoeconomics

The **incentive and mechanism-design layer** of blockchains: how tokens are issued and capture value, how a chain pays for its own security, how transaction fees are priced, how value leaks to block producers (MEV), and how token-holders govern. This skill covers the *economics*; it defers protocol mechanics, consensus safety proofs, cryptographic primitives, and trading advice to peer skills.

> **Scope boundary (read first).** This skill answers "what are the incentives / what is the economic argument," not "how does the protocol work mechanically." See **SKIP** in the frontmatter and **Peer skills & deferral** below.

## When this skill applies

Reach for `blockchain-economics` when the question is about *incentives, value, supply, yield, fees-as-a-market, extractable value, or governance power* — e.g. "is this token's emission schedule inflationary," "how much does it cost to 51%-attack this chain," "why isn't EIP-1559 a first-price auction," "what is a sandwich attack and who captures the value," "why do DAO treasuries blow up," "what is the Curve-wars bribe market." If the question is "how does the EVM/UTXO/consensus actually work," route to the per-chain or consensus peer.

## How to use the references

This skill is split into five focused references. Load the one(s) matching the question; each is self-contained with its own citations and confidence tags.

| If the question is about… | Read |
|---|---|
| Token supply, issuance/emission, inflation/disinflation, burns/sinks, ultrasound money, vesting & unlocks, fair-launch vs VC, airdrops, utility/governance/security taxonomy + Howey, the token-velocity problem (MV=PQ) | `references/tokenomics.md` |
| Cryptoeconomics as a discipline, incentive-compatibility, the security budget & 51%/cost-of-attack economics, staking economics (real yield, slashing, liquid staking, restaking), PoS-vs-PoW economic security, Schelling points, the verifier's dilemma | `references/cryptoeconomics-and-security.md` |
| Transaction-fee mechanism design: EIP-1559 base fee & burn, why it's not a pure auction, priority fees, the EIP-4844 blob fee market, Solana localized/per-account fee markets, first-vs-second-price/VCG and the IC tradeoffs | `references/fee-markets.md` |
| MEV: definition & Flash Boys 2.0, the taxonomy (arbitrage/sandwich/liquidation/JIT/long-tail), the searcher→builder→relay→proposer supply chain, PBS & enshrined-PBS (EIP-7732), Flashbots/MEV-Boost & relay censorship, order-flow auctions, intents & solvers (CoW/UniswapX), mitigation (encrypted mempools, FCFS, MEV-Share/burn) | `references/mev.md` |
| Governance economics: token-weighted vs quadratic vs conviction voting, vote-buying & bribe markets (Curve wars / Votium / Hidden Hand), DAO treasuries & the native-token concentration doom-loop, governance attacks (flash-loan / quorum-apathy) | `references/governance-economics.md` |

## The five pillars in one screen

1. **Tokenomics** — supply schedule (fixed vs uncapped; the *flow*, not the cap, is what matters long-run), sinks (EIP-1559 burn, buyback-and-burn) and net issuance, vesting/unlock sell-pressure (≈90% of unlocks are net-negative; cliffs concentrate it), distribution (fair-launch vs VC; airdrop recipients overwhelmingly dump), the utility/governance/security taxonomy, and the **token-velocity problem** (doubling velocity halves price; "usage without value capture").

2. **Cryptoeconomics & security** — mechanism design as "reverse game theory" targeting incentive-compatibility; the **security budget** (miner/validator revenue caps the attack cost; Budish's flow-vs-stock bound); **staking economics** (nominal vs real/dilution-adjusted yield, slashing & correlation penalties, liquid staking centralization, restaking pooled security); **PoS vs PoW** ("security from value-at-loss," slashable stake vs un-slashable hashpower); Schelling points and the verifier's dilemma.

3. **Fee markets** — EIP-1559 is a **posted-price** mechanism (history-dependent burned base fee) **+ a first-price tip auction at the margin**, *not* a pure auction; the **burn** is load-bearing (it blocks the off-chain refund that would collapse it back to a first-price auction); blobs (EIP-4844) get an **independent** fee dimension; Solana prices contention **per-account**; and a true second-price/VCG auction is **not** miner-incentive-compatible (shill-bidding).

4. **MEV** — value from controlling transaction *ordering/inclusion*; a taxonomy (toxic sandwich/front-run vs benign arbitrage/liquidation, plus JIT and long-tail); a **supply chain** (searcher → builder → relay → proposer) that PBS exists to keep validator-decentralized while builder-specialized; off-chain PBS (MEV-Boost + trusted relays) → **enshrined PBS** (EIP-7732); and a mitigation menu (encrypted mempools, FCFS, intents/OFAs, MEV-Share/burn) — **MEV is intrinsic and can only be redistributed/minimized, not eliminated.**

5. **Governance economics** — coin (token-weighted) voting is plutocratic and only "feels neutral"; quadratic and conviction voting try to escape plutocracy but break on **sybil/identity** and **collusion/bribery**; **bribe markets** (Curve wars, Votium, Hidden Hand) make vote-buying an open, "legal," capital-efficient market; and **DAO treasuries** carry a native-token-concentration **doom-loop** (treasury collapses exactly when the protocol most needs it).

## Cross-cutting through-lines (cite these patterns)

- **Burning is a mechanism-design tool, not just monetary policy.** The same logic recurs: burn the base fee (EIP-1559) so the proposer can't manipulate it and users can't refund off-chain; burn the per-account write-lock fee (Solana SIMD-0110) so validators can't game it; burn MEV via a base-fee floor (MEV-burn) so value goes to all holders, not the lucky proposer. (`fee-markets`, `mev`)
- **"Security/value comes from value-at-loss."** Slashable stake (PoS), capital-at-risk attack cost (Budish), and vote-escrow lockups (veTokenomics) are the same idea: make the honest action cheap and the dishonest action forfeit real, non-recoverable value. (`cryptoeconomics-and-security`, `governance-economics`)
- **Centralization keeps reappearing one layer up.** PBS decentralizes proposing but concentrates *building*; intents/solvers internalize MEV for users but concentrate the *solver* market; liquid staking democratizes staking but concentrates *stake* (Lido near the ⅓ threshold). The honest read is "where did the concentration move," not "is it solved." (`mev`, `cryptoeconomics-and-security`)
- **Identity & collusion are the unsolved substrate.** Quadratic voting, SchellingCoin oracles, and airdrop fairness all assume one-person-one-identity and no side-payments; sybil attacks and the P+ε bribery attack break that assumption, which is why most anti-plutocratic mechanisms ship only in non-binding/funding contexts. (`governance-economics`, `cryptoeconomics-and-security`)

## Peer skills & deferral

This skill is a sibling in the **Blockchain protocols** family. It is the *economics* node; it defers:

| Topic | Defer to | Why |
|---|---|---|
| AMM/DEX math, lending-pool mechanics, DeFi protocol internals | `defi` (planned sibling) | This skill covers MEV's/economics' *interaction* with DeFi, not the protocol math. |
| Consensus theory & safety/liveness proofs (CAP/FLP/BFT, Gasper internals, fork-choice) | `distributed-systems-consensus` (installed; reciprocally defers token-economics here) | We use only the *economic* security argument (cost-of-attack, slashing-as-value-at-loss). |
| Per-chain protocol detail (EVM/gas accounting, UTXO, Bitcoin issuance code, the Merge/blobs as protocol) | `bitcoin-protocol-expert`, `ethereum-protocol-expert` (installed) | They state the *mechanism*; the *economic* abstraction (burn-as-sink, fee-as-market) lives here. |
| Cryptographic primitives (signatures, VRFs, KZG/hashing math) | `blockchain-crypto-primitives` (planned sibling) | Out of scope; only the role of primitives in mechanisms is noted. |
| Consumer crypto investing/trading, valuation as a buy/sell call | `investing-and-retirement`, `trading-and-investing` | **This skill is protocol/mechanism economics, not trading or investment advice.** |

**Peer status note.** `distributed-systems-consensus`, `bitcoin-protocol-expert`, and `ethereum-protocol-expert` are installed and the Ethereum/Bitcoin skills already point token-economics here. `defi` and `blockchain-crypto-primitives` are planned siblings; where absent, that domain is simply out of scope here (the deferral is the honest boundary, not a broken link).

## Anti-patterns this skill guards against

- **Treating tokenomics as price prediction.** The velocity/EoE frame is a *directional heuristic for value capture*, empirically shaky as a pricing engine; never present MV=PQ as a valuation that justifies a trade.
- **"Ultrasound money" as a settled fact.** ETH is net-deflationary only when burn > issuance; post-Dencun it has been largely net-inflationary. State the *mechanism*; flag the *direction* as volatile and contested.
- **"EIP-1559 lowered fees."** It changed *how* blockspace is priced (easier estimation, lower intra-block variance), not *how much* it costs. False as commonly stated.
- **"PoS is strictly more secure than PoW" (or vice-versa).** Preserve the genuine disagreement: capital-at-risk/slashing vs energy-cost; both camps agree on structure, disagree on the conclusion.
- **"MEV can be eliminated."** It is intrinsic to permissionless ordering; mitigation redistributes/minimizes/relocates it. Encrypted mempools are not a silver bullet (metadata leakage, liveness risk).
- **Quoting a single volatile number as fact.** Staking yields, relay/builder/solver concentration, bribe-market volumes, unlock calendars, and ePBS scheduling all move; carry the `verified-as-of` stamp and qualify.

## References (overview-level)

Each reference file carries its own full, numbered `## References`. Foundational anchors span the field:

- Vitalik Buterin, "A Proof of Stake Design Philosophy" (2016) — https://medium.com/@VitalikButerin/a-proof-of-stake-design-philosophy-506585978d51
- Daian et al., "Flash Boys 2.0" (arXiv:1904.05234; IEEE S&P 2020) — https://arxiv.org/abs/1904.05234
- Tim Roughgarden, "Transaction Fee Mechanism Design for the Ethereum Blockchain" (arXiv:2012.00854) — https://timroughgarden.org/papers/eip1559.pdf
- Eric Budish, "The Economic Limits of Bitcoin and the Blockchain" (NBER 24717 / QJE 2024) — https://www.nber.org/system/files/working_papers/w24717/w24717.pdf
- Kyle Samani (Multicoin), "Understanding Token Velocity" (2017) — https://multicoin.capital/2017/12/08/understanding-token-velocity/
