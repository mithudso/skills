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

**Incentive and mechanism-design layer** of blockchains: token issuance + value capture, chain security payment, fee pricing, MEV, governance. Covers *economics*; defers protocol mechanics, consensus proofs, crypto primitives, trading advice to peer skills.

> **Scope boundary (read first).** Skill answers "what are incentives / what is economic argument," not "how does protocol work mechanically." See **SKIP** in frontmatter and **Peer skills & deferral** below.

## When this skill applies

Use `blockchain-economics` when question is about *incentives, value, supply, yield, fees-as-market, extractable value, or governance power* — e.g. "is this token's emission schedule inflationary," "how much to 51%-attack this chain," "why isn't EIP-1559 first-price auction," "what is sandwich attack and who captures value," "why do DAO treasuries blow up," "what is Curve-wars bribe market." If question is "how does EVM/UTXO/consensus work," route to per-chain or consensus peer.

## How to use the references

Five focused references. Load matching one(s); each self-contained with own citations and confidence tags.

| If question is about… | Read |
|---|---|
| Token supply, issuance/emission, inflation/disinflation, burns/sinks, ultrasound money, vesting & unlocks, fair-launch vs VC, airdrops, utility/governance/security taxonomy + Howey, token-velocity problem (MV=PQ) | `references/tokenomics.md` |
| Cryptoeconomics as discipline, incentive-compatibility, security budget & 51%/cost-of-attack, staking economics (real yield, slashing, liquid staking, restaking), PoS-vs-PoW economic security, Schelling points, verifier's dilemma | `references/cryptoeconomics-and-security.md` |
| Fee mechanism design: EIP-1559 base fee & burn, why not pure auction, priority fees, EIP-4844 blob fee market, Solana localized/per-account fee markets, first-vs-second-price/VCG and IC tradeoffs | `references/fee-markets.md` |
| MEV: definition & Flash Boys 2.0, taxonomy (arbitrage/sandwich/liquidation/JIT/long-tail), searcher→builder→relay→proposer supply chain, PBS & enshrined-PBS (EIP-7732), Flashbots/MEV-Boost & relay censorship, order-flow auctions, intents & solvers (CoW/UniswapX), mitigation (encrypted mempools, FCFS, MEV-Share/burn) | `references/mev.md` |
| Governance economics: token-weighted vs quadratic vs conviction voting, vote-buying & bribe markets (Curve wars / Votium / Hidden Hand), DAO treasuries & native-token concentration doom-loop, governance attacks (flash-loan / quorum-apathy) | `references/governance-economics.md` |

## The five pillars in one screen

1. **Tokenomics** — supply schedule (fixed vs uncapped; *flow*, not cap, matters long-run), sinks (EIP-1559 burn, buyback-and-burn) and net issuance, vesting/unlock sell-pressure (≈90% of unlocks net-negative; cliffs concentrate it), distribution (fair-launch vs VC; airdrop recipients overwhelmingly dump), utility/governance/security taxonomy, **token-velocity problem** (doubling velocity halves price; "usage without value capture").

2. **Cryptoeconomics & security** — mechanism design as "reverse game theory" targeting incentive-compatibility; **security budget** (miner/validator revenue caps attack cost; Budish's flow-vs-stock bound); **staking economics** (nominal vs real/dilution-adjusted yield, slashing & correlation penalties, liquid staking centralization, restaking pooled security); **PoS vs PoW** ("security from value-at-loss," slashable stake vs un-slashable hashpower); Schelling points and verifier's dilemma.

3. **Fee markets** — EIP-1559 is **posted-price** mechanism (history-dependent burned base fee) **+ first-price tip auction at margin**, *not* pure auction; **burn** is load-bearing (blocks off-chain refund that would collapse it back to first-price auction); blobs (EIP-4844) get **independent** fee dimension; Solana prices contention **per-account**; true second-price/VCG auction is **not** miner-incentive-compatible (shill-bidding).

4. **MEV** — value from controlling tx *ordering/inclusion*; taxonomy (toxic sandwich/front-run vs benign arbitrage/liquidation, plus JIT and long-tail); **supply chain** (searcher → builder → relay → proposer) — PBS keeps validators decentralized while builders specialize; off-chain PBS (MEV-Boost + trusted relays) → **enshrined PBS** (EIP-7732); mitigation menu (encrypted mempools, FCFS, intents/OFAs, MEV-Share/burn) — **MEV intrinsic, only redistributable/minimizable, not eliminable.**

5. **Governance economics** — coin (token-weighted) voting plutocratic, only "feels neutral"; quadratic and conviction voting escape plutocracy but break on **sybil/identity** and **collusion/bribery**; **bribe markets** (Curve wars, Votium, Hidden Hand) make vote-buying open, "legal," capital-efficient; **DAO treasuries** carry native-token-concentration **doom-loop** (treasury collapses exactly when protocol most needs it).

## Cross-cutting through-lines (cite these patterns)

- **Burning is mechanism-design tool, not just monetary policy.** Same logic recurs: burn base fee (EIP-1559) so proposer can't manipulate + users can't refund off-chain; burn per-account write-lock fee (Solana SIMD-0110) so validators can't game it; burn MEV via base-fee floor (MEV-burn) so value goes to all holders, not lucky proposer. (`fee-markets`, `mev`)
- **"Security/value comes from value-at-loss."** Slashable stake (PoS), capital-at-risk attack cost (Budish), vote-escrow lockups (veTokenomics) — same idea: make honest action cheap, dishonest action forfeit real non-recoverable value. (`cryptoeconomics-and-security`, `governance-economics`)
- **Centralization reappears one layer up.** PBS decentralizes proposing but concentrates *building*; intents/solvers internalize MEV for users but concentrate *solver* market; liquid staking democratizes staking but concentrates *stake* (Lido near ⅓ threshold). Honest read: "where did concentration move," not "is it solved." (`mev`, `cryptoeconomics-and-security`)
- **Identity & collusion are unsolved substrate.** Quadratic voting, SchellingCoin oracles, airdrop fairness all assume one-person-one-identity and no side-payments; sybil attacks and P+ε bribery attack break that — why most anti-plutocratic mechanisms ship only in non-binding/funding contexts. (`governance-economics`, `cryptoeconomics-and-security`)

## Peer skills & deferral

Sibling in **Blockchain protocols** family. *Economics* node; defers:

| Topic | Defer to | Why |
|---|---|---|
| AMM/DEX math, lending-pool mechanics, DeFi protocol internals | `defi` (planned sibling) | Skill covers MEV's/economics' *interaction* with DeFi, not protocol math. |
| Consensus theory & safety/liveness proofs (CAP/FLP/BFT, Gasper internals, fork-choice) | `distributed-systems-consensus` (installed; reciprocally defers token-economics here) | Only *economic* security argument used (cost-of-attack, slashing-as-value-at-loss). |
| Per-chain protocol detail (EVM/gas accounting, UTXO, Bitcoin issuance code, the Merge/blobs as protocol) | `bitcoin-protocol-expert`, `ethereum-protocol-expert` (installed) | They state *mechanism*; *economic* abstraction (burn-as-sink, fee-as-market) lives here. |
| Cryptographic primitives (signatures, VRFs, KZG/hashing math) | `blockchain-crypto-primitives` (planned sibling) | Out of scope; only role of primitives in mechanisms noted. |
| Consumer crypto investing/trading, valuation as buy/sell call | `investing-and-retirement`, `trading-and-investing` | **Skill is protocol/mechanism economics, not trading or investment advice.** |

**Peer status note.** `distributed-systems-consensus`, `bitcoin-protocol-expert`, `ethereum-protocol-expert` installed; Ethereum/Bitcoin skills already point token-economics here. `defi` and `blockchain-crypto-primitives` planned siblings; where absent, domain is out of scope (deferral is honest boundary, not broken link).

## Anti-patterns this skill guards against

- **Treating tokenomics as price prediction.** Velocity/EoE frame is *directional heuristic for value capture*, empirically shaky as pricing engine; never present MV=PQ as valuation justifying a trade.
- **"Ultrasound money" as settled fact.** ETH net-deflationary only when burn > issuance; post-Dencun largely net-inflationary. State *mechanism*; flag *direction* as volatile and contested.
- **"EIP-1559 lowered fees."** Changed *how* blockspace priced (easier estimation, lower intra-block variance), not *how much* it costs. False as commonly stated.
- **"PoS strictly more secure than PoW" (or vice-versa).** Preserve genuine disagreement: capital-at-risk/slashing vs energy-cost; both camps agree on structure, disagree on conclusion.
- **"MEV can be eliminated."** Intrinsic to permissionless ordering; mitigation redistributes/minimizes/relocates it. Encrypted mempools not silver bullet (metadata leakage, liveness risk).
- **Quoting single volatile number as fact.** Staking yields, relay/builder/solver concentration, bribe-market volumes, unlock calendars, ePBS scheduling all move; carry `verified-as-of` stamp and qualify.

## References (overview-level)

Each reference file carries own full, numbered `## References`. Foundational anchors:

- Vitalik Buterin, "A Proof of Stake Design Philosophy" (2016) — https://medium.com/@VitalikButerin/a-proof-of-stake-design-philosophy-506585978d51
- Daian et al., "Flash Boys 2.0" (arXiv:1904.05234; IEEE S&P 2020) — https://arxiv.org/abs/1904.05234
- Tim Roughgarden, "Transaction Fee Mechanism Design for the Ethereum Blockchain" (arXiv:2012.00854) — https://timroughgarden.org/papers/eip1559.pdf
- Eric Budish, "The Economic Limits of Bitcoin and the Blockchain" (NBER 24717 / QJE 2024) — https://www.nber.org/system/files/working_papers/w24717/w24717.pdf
- Kyle Samani (Multicoin), "Understanding Token Velocity" (2017) — https://multicoin.capital/2017/12/08/understanding-token-velocity/