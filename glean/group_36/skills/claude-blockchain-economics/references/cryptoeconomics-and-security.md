<!-- Provenance: reference under the `blockchain-economics` standalone skill. Created 2026-06-16 via /dr deep-research. Synthesized from web sources; all load-bearing/volatile claims carry [^n] citations to the ## References section. -->

# Cryptoeconomics, Mechanism Design & Security Economics

verified-as-of: 2026-06-16 (volatile: staking yields, Lido share, restaking TVL, attack-cost dashboards, EIP-7716 status — re-verify before relying)

How incentives are designed so honest behavior is rational, how a chain pays for the cost of attacking it, and the economics of staking. Confidence tags: **[FACT]** = primary/peer-reviewed + multiply attested; **[QUALIFIED]** = 2 sources or model-dependent; **[TENTATIVE]** = single-source/draft. The *economic* security argument only — consensus safety/liveness proofs defer to `distributed-systems-consensus`.

## Contents

- Cryptoeconomics as a discipline
- The security budget & cost-of-attack (51% economics)
- Staking economics (yield, slashing, liquid staking, restaking)
- PoS vs PoW economic security
- Schelling points (SchellingCoin, oracles, TCRs)
- The verifier's dilemma
- References

## Cryptoeconomics as a discipline

**Origin & definition.** The term originated in the Ethereum community ~2014–2015; the earliest documented public use is **Vlad Zamfir's 2015 "What is Cryptoeconomics?"**[^1][^2] Operational definition: *"the application of incentive [mechanism] design to problems that are traditionally cryptography and distributed-systems problems"* — adding incentives to the cryptographer's toolkit.[^3][^4] **[FACT]**

**Placement.** Academic treatments place cryptoeconomics as **a branch of mechanism design, itself a branch of microeconomics** (Davidson et al. 2016; the 2024 arXiv survey "Cryptoeconomics and Tokenomics as Economics").[^1] Mechanism design is "reverse game theory": fix the desired outcome (honest consensus), then design payoffs that make it a Nash equilibrium.[^1][^2] **[FACT]**

**Incentive-compatibility** is the load-bearing property: following the protocol should be a (ideally dominant) strategy, robust against *coalitional* deviations, while *maximizing the cost of attack* — Zamfir's three Casper design targets.[^3] **[FACT for the framing; QUALIFIED that proving a unique live-protocol Nash equilibrium is aspirational.]**

**Vitalik's "design philosophy"** (2016): the best protocols are robust across *many* adversary models simultaneously (economic-rational individual & coordinated, fault-tolerant, BFT, even behavioral "we all cheat a little"), layering *economic incentives to deter cartels from acting anti-socially* with *anti-centralization incentives to deter cartels from forming at all.*[^5] A definitional contrast worth preserving: for Buterin cryptoeconomics is primarily a methodology for information-security guarantees; Zamfir's scope is the whole digital economy.[^2] **[FACT / QUALIFIED.]**

## The security budget & cost-of-attack (51% economics)

**The budget identity (PoW).** Miner revenue = block subsidy + transaction fees (coin × price). Revenue **caps** mining expenditure (miners won't spend more than they earn), and that expenditure is what makes the chain expensive to attack — so **honest revenue is the upper bound on what an attacker must outspend.**[^6][^7] At current Bitcoin hashrate the daily budget is roughly **$6–7M/day** (ignoring capex).[^7] **[FACT, volatile.]**

**Budish's "Economic Limits" (NBER 2018; QJE 2024)** — the canonical academic model. A three-equation argument (miner zero-profit + attack incentive-compatibility) implies the recurring **flow** payment to miners must be large relative to the one-off **stock** value of attacking: `p_block > V_attack / (A*·t(A*))`. Verdict: *"a very expensive form of trust"* — in his numerical example the **annual cost of trust must be ~400,000% of the value protected**, and *"in some scenarios the cost exceeds global GDP."* Implication: a pure-PoW chain *"would be majority attacked if it became sufficiently economically important."*[^8][^9] Crucially, the constraint *"would become much more attractive if an attacker were to lose the **stock** value of their capital in addition to the flow cost"* — which is exactly what **PoS slashing supplies** (see PoS-vs-PoW).[^9] **[FACT for model/quotes; QUALIFIED for the "would be attacked" conclusion, contested below.]**

**Buy vs rent.**
- **Buy hardware:** a physical Bitcoin attack ≈ **$5.5B+** in ASICs (conservative 2021 estimate) plus a coordination problem.[^10] **[QUALIFIED.]**
- **Rent hashpower:** for smaller PoW coins, renting (NiceHash) drives attacker *fixed* costs toward zero — "renters only need hashrate for the duration of the attack." A coin is high-risk when its **NiceHash-able % exceeds ~100%**; MIT DCI documents rented attacks as "break-even or profitable unless miners have large non-recoupable fixed costs." Quoted attack costs typically *exclude* block rewards earned during the attack, which can cut net cost ~80%.[^11][^12][^13][^14] **[FACT.]**

**The long-run "security-budget cliff."** Bitcoin's subsidy → 0 by ~2140; security must then be fee-funded.
- **Instability camp:** Carlsten et al. (CCS 2016) show fee-only mining introduces **undercutting** and sharpens selfish-mining incentives via fee *variance*; a 2026 model finds a fee-only deviation threshold "just 0.17%," with honesty today resting on subsidy dominance. Budish frames it as a *ratio* problem (budget-to-attack-value deteriorates as subsidy declines and secured value grows).[^6][^9][^15][^16][^17] **[QUALIFIED.]**
- **Sufficiency camp:** scarcity raises fees; layered settlement concentrates fee value; price appreciation buys time (subsidy halves in BTC, but if price doubles, USD security is flat). Post-2024-halving fees have run ~5–15% of reward — "no obvious crisis."[^7][^15] **[QUALIFIED — genuine open debate.]**

## Staking economics (yield, slashing, liquid staking, restaking)

**Nominal vs real yield.** *Nominal* = gross reward rate; *real* adjusts for **dilution** from new issuance. Dilution affects staked and unstaked holders equally, so it doesn't change the staking *incentive gap*, but it does reduce the *real* return — as staking ratio rises, much of the net incentive becomes "dilution protection," and in the limit (ratio→1) "real total staking yield only consists of MEV yield."[^19] Concrete (volatile): ETH nominal ≈3.08% / real ≈2.73% (late 2024); SOL nominal ≈11.5%; by 2026 ~33.4M ETH staked (~28%) at ~3.1% base APR.[^20][^21] **[FACT, volatile.]**

**Staking ratio & the yield curve.** Ethereum issuance follows an **inverse-square-root** curve (yield falls as more ETH stakes), which Wharton's Jermann shows is "qualitatively consistent with an optimal Ramsey policy" targeting a stable ratio (security vs dilution). Real-yield models derive negative-real-yield break-points (solo ~70M, institutional ~90M, LST ~100M ETH staked), beyond which "staking is a net cost for everyone" — motivating the issuance-curve-reduction debate (a lower nominal curve can yield *higher* real equilibrium yield).[^22][^23][^24] **[QUALIFIED — model-dependent, multiple analyses converge.]**

**Slashing & correlation penalties (Ethereum).** Three components for provably malicious acts (double-proposal, surround/double attestation):[^25][^26][^27]
1. **Initial slash** ("up to 1 ETH"), then a ~36-day forced exit during which the stake bleeds.
2. **Correlation penalty (~day 18)** scaling with total stake slashed in the surrounding window: `correlation_penalty = min(B, 3·S·B/T)` (B = validator effective balance, S = co-slashed sum, T = total active balance; **PROPORTIONAL_SLASHING_MULTIPLIER = 3** since Bellatrix). Intent: a lone fault is lightly punished; a mass event is punished up to the entire stake.
3. **Inactivity leak** when the chain fails to finalize.
Anti-correlation as a *decentralization tool*: **EIP-7716** (draft) extends correlation penalties to *mundane* missed attestations, because a single large operator on shared infra suffers correlated faults — penalizing correlation "reduces the economies of scale large operators enjoy."[^25][^28] **[FACT for current slashing; TENTATIVE for EIP-7716 (draft).]**

**Liquid staking (LSTs).** Issues a transferable receipt (Lido's **stETH**, rebasing) so capital earns yield and stays usable in DeFi — the largest DeFi category (~$58B TVL at peak).[^29][^30] Central risk = **centralization toward the ⅓ BFT threshold**: Lido peaked ~31.8% of staked ETH (May 2023), ~28.6% by 2026 — near the **33.3% finality-disruption** threshold. In a **June 2023 self-limit vote, ~99.8% of LDO voted *against* capping** Lido's share — read as evidence that "economic self-interest overrides decentralization concerns in a token-governed protocol." stETH is also the most-used DeFi collateral, so a depeg is *systemic* (it traded to **$0.935 on Curve in June 2022**, triggering ~$1.18B in liquidations).[^29][^30][^31][^32][^33] **[FACT for numbers; QUALIFIED that "winner-take-most is a fundamental law of PoS" is contested.]**

**Restaking (EigenLayer) — economic level.** Stakers **opt in to additional slashing conditions** to secure other "Actively Validated Services" (AVSs), extending Ethereum's security via **pooled security** rather than each service bootstrapping its own token; aggregating fragmented stake raises the **cost-of-corruption** and stakers earn extra fees.[^35][^36][^37] The **EIGEN** token handles **intersubjective faults** (not on-chain-objectively provable, but socially agreeable) via a **fork-and-slash** dispute mechanism (a generalization of SchellingCoin), with ETH (objective safety) + EIGEN (liveness) dual-staking.[^36] **[FACT for design.]** **Negation:** restaking introduces **systemic/contagion risk** — leveraged-restaking depeg cascades, correlated/cascading slashing, and critiques that "EIGEN yield is a mirage / intersubjective forking is untested"; the research frontier applies Eisenberg–Noe / Diamond–Dybvig network-stability models to LRTs.[^41] **[QUALIFIED — contested.]**

## PoS vs PoW economic security

**Core asymmetry — "security comes from value-at-loss."** Buterin: security comes not from burning energy but from putting up economic value-at-loss; a block has $X security if reverting it requires malicious validators to forfeit $X in in-protocol penalties (penalties "hundreds or thousands of times larger than rewards"). PoW relies on rewards, not penalties, so its defense is symmetric (attacker and defender both just spend on hashpower).[^5] **[FACT.]**

**Cost-to-attack, quantified ("Why Proof of Stake," 2020), per $1/day of rewards:**[^38]
- GPU PoW ≈ **$0.26** (→ near zero, attacker earns rewards while attacking).
- ASIC PoW ≈ **$486** (≈⅔ capital / ⅓ ongoing).
- PoS ≈ **$2,189** (→ possibly ~$10,000).
Bottom line: PoS gives a **~5–20× security-per-dollar** edge, because PoS is "almost entirely capital costs" (deposits don't depreciate and are returned), whereas PoW capital (ASICs) depreciates and must be re-funded by energy.[^38] **[FACT for the quoted figures; QUALIFIED — illustrative model numbers, not market measurements.]**

**"Stake can be slashed but hardware can't."** In PoS a large portion of the *attacker's* stake (and no one else's) is automatically destroyed; a PoW community can only brick ASICs by hard-forking the PoW, which bricks *honest* miners equally (no defender advantage). Recovery is asymmetric — the community is "back on its feet within days," and a repeat attack requires re-acquiring coins whose price the first attack depressed.[^5][^38] This is exactly the external-punishment channel Budish's PoW model says is missing.[^9] **[FACT — Buterin/Budish converge.]**

**Weak-subjectivity economics.** PoS can't be *fully* objective: once deposits are withdrawn, there's no incentive not to vote on a long-range fork (the long-range-attack problem). The remedy is **weak subjectivity** — a new/long-offline node obtains a recent *trusted checkpoint* and syncs forward from a state < N blocks old; within that window, locked deposits make rewriting history economically destructive. Buterin: weak subjectivity is "both sufficient and necessary" to sidestep nothing-at-stake, at a moderate one-time-trust cost.[^39][^40] **[FACT.]**

**Contradiction preserved — is PoS actually weaker?** A PoW-maximalist / DBA "Token Value Capture" counter: "sort-of-productive money" may be a *weakness* — if ETH/SOL grew into the trillions, "real yield would squish to near 0," whereas BTC's non-productive design sidesteps the yield-competition treadmill. And the Lido case argues PoS-in-practice may concentrate *more* than PoW at the liquid-staking layer.[^32][^41] **[QUALIFIED — genuine disagreement, no consensus.]**

## Schelling points (SchellingCoin, oracles, TCRs)

**The mechanism.** A **Schelling (focal) point** is the answer parties converge on without communication. **SchellingCoin** (Buterin, 2014) weaponizes it for decentralized oracles: everyone independently submits a value; submitters near the **25th–75th percentile (the median)** are rewarded, via **commit-reveal**, because "the truth is the most powerful Schelling point" (hard to coordinate on a specific lie). This is the lineage behind Augur-style oracles and EigenLayer's intersubjective forking.[^42][^43][^36] **[FACT.]**

**Token-curated registries (TCRs)** (Goldin, 2017): candidates post a deposit to be listed; "bad" listings can be **challenged**, resolved by token-weighted commit-reveal voting, with the loser's deposit redistributed. Incentive thesis: token holders have a *tactical* incentive to grab deposits but a *strategic* incentive to curate well (a quality list raises token value).[^45][^46] **[FACT.]**

**Negation / fragility:**
- **The P + ε attack** (Miller; Buterin 2015): an attacker credibly commits to pay `P + ε` to everyone who votes their way *only if the attack fails* — making defection a dominant strategy *regardless of beliefs*, so the majority defects and **the attacker pays nothing**. No epistemic takeover required.[^47] **[FACT — the canonical disconfirming result.]**
- **Bribery bidding wars:** costless bribes invite an auction "largely independent of the value at stake"; oracle security paradoxically *increases* when bidders have minimal information about each other.[^48] **[QUALIFIED.]**
- **Mitigations:** Sztorc/Truthcoin makes minority voters lose more as a vote gets contentious; *subjectivity* (nodes can ignore an attacker fork) reduces exploitability.[^44][^47] A conceptual critique (Gerbrandy): Schelling points strictly apply only where communication is *impossible*, so TCR security really rests on the *external* token-value incentive, not the focal-point game.[^49] **[QUALIFIED.]**

## The verifier's dilemma

**The dilemma** (Luu, Teutsch, Kulkarni, Saxena — "Demystifying Incentives in the Consensus Computer," CCS 2015). When validating a block is expensive, a rational miner faces an ill-fated choice: verify (and fall behind the race for the next block) or **skip verification** (gain a head start, risk building on an invalid block). The paper proves "rational miners are well-incentivized to accept unvalidated blockchains" — a transaction can be crafted to consume, e.g., 30% of the block interval to verify, letting skippers win the chain. Their resolution incentivizes correct execution only for scripts that are **cheap to verify**.[^50][^51] **[FACT.]**

**Strengthened / empirical:**
- **"Sluggish Mining" (FC 2019):** purpose-built slow contracts give an edge; the gas limit mitigates but doesn't eliminate it.[^52]
- **DSN 2020 (>300k contracts):** "it is often economically rational not to verify"; mitigations include parallelization and **active invalid-block insertion** (fraud-proof-like deterrent).[^53]
- **Real-world:** PoW pools mine **empty blocks** on a new block's hash before validating (SPV-style lazy validation); Ethereum's state-root commitment forces prev-block validation, making naive skipping harder than on Bitcoin.[^53][^54] **[FACT.]**

**Relation to fraud proofs / re-execution incentives.** The verifier's dilemma is *why optimistic systems can't assume everyone re-executes* — it motivates **fraud proofs** (re-execute only when challenged), **bonded challengers paid from the slashed bond** of a faulty proposer, and **validity (succinct) proofs** that move cost from every verifier to a single prover (echoing Luu et al.'s "make verification cheap").[^50][^53] **[QUALIFIED — the cited sources establish the incentive premise; the fraud-proof connection is well-established in the literature.]**

## Child concepts (future research)

1. MEV economics & PBS — see `references/mev.md` (the "real yield" sources keep deferring here).
2. Restaking systemic-risk & financial-contagion modeling (Eisenberg–Noe, Diamond–Dybvig applied to LRTs).
3. Bribery/collusion & credible-commitment economics (P+ε, dark-DAO vote-buying, time-bandit reorgs).
4. Token value-accrual & "REV" frameworks (the "internal tax" view of PoS issuance; how value capture interacts with the security budget).
5. Cost-of-corruption vs cost-of-attack formalism (EigenLayer's profit-from-corruption ≤ cost-of-corruption inequality).

## References

[^1]: Cryptoeconomics and Tokenomics as Economics: A Survey — https://arxiv.org/html/2407.15715 — academic survey — Zamfir 2015 origin; cryptoeconomics as a branch of mechanism design / microeconomics.
[^2]: Brekke & Alsindi, "Cryptoeconomics," Internet Policy Review — https://policyreview.info/pdf/policyreview-2021-2-1553.pdf — peer-reviewed — Buterin (info-security) vs Zamfir (digital-economy) definitions.
[^3]: Vlad Zamfir, "Cryptoeconomics in Casper" (CESC 2017) — https://www.youtube.com/watch?v=5ScY7ruD_eg — primary talk — incentive mechanism design applied to crypto/distributed-systems; 3 Casper targets.
[^4]: Tim Swanson on Zamfir's "What is Cryptoeconomics?" — https://www.ofnumbers.com/2015/01/30/cryptoeconomics-for-beginners-and-experts-alike/ — secondary — Nash-equilibrium-has-cryptoeconomic-security framing.
[^5]: Vitalik Buterin, "A Proof of Stake Design Philosophy" (2016) — https://medium.com/@VitalikButerin/a-proof-of-stake-design-philosophy-506585978d51 — essay (primary) — "security from value-at-loss"; multi-model robustness; two-layer defense.
[^6]: Geo Nicolaidis, "The Security Budget Cliff" (2026) — https://geonicolaidis.substack.com/p/issue-9-the-security-budget-cliff — newsletter — budget identity; Budish ratio framing.
[^7]: SatoshiBench, "Security Budget — Subsidy + Fees Math" (2026) — https://satoshibench.com/learn/security-budget/ — explainer — ~$6–7M/day; fees ~5–15% of reward; sufficiency vs instability.
[^8]: Eric Budish, "The Economic Limits of Bitcoin and the Blockchain" (NBER 24717, 2018) — https://www.nber.org/system/files/working_papers/w24717/w24717.pdf — working paper — flow>stock argument; "store of value akin to gold" attack conclusion.
[^9]: Eric Budish, QJE 2024 version — https://socialsciences.uchicago.edu/sites/default/files/2024-09/Economic%20Limits%20Crypto%20Blockchains%20-%20QJE%20Sept%202024.pdf — peer-reviewed — "400,000%"; "exceeds global GDP"; `p_block > V_attack/(A*·t(A*))`; external-punishment softens it.
[^10]: Braiins, "Cost to 51% Attack Bitcoin?" (2021) — https://braiins.com/blog/how-much-would-it-cost-to-51-attack-bitcoin — industry blog — ~$5.5B hardware estimate; NiceHash-able metric.
[^11]: Crypto51, "Cost of a 51% Attack" — https://www.crypto51.app/ — dashboard — per-coin attack costs; block rewards cut net cost ~80%.
[^12]: KuCoin, "PoW 51% Attack Cost in 2026" — https://www.kucoin.com/blog/hk-understanding-the-moat-what-is-pow-51-attack-cost-in-2026 — explainer — >100% NiceHash-able = rentable-attack risk.
[^13]: MIT Digital Currency Initiative, "51% Attacks" — https://www.dci.mit.edu/projects/51-percent-attacks — research — rental drives attacker fixed costs to zero; "break-even or profitable."
[^14]: CoinDesk, "51% Attacks for Rent" (2019) — https://www.coindesk.com/markets/2019/02/23/51-attacks-for-rent-the-trouble-with-a-liquid-mining-market — journalism — rent-a-miner "1000x cheaper" than buying.
[^15]: Bitcoin Institute, fee-only future / security-budget debate (2026) — https://bitcoin-institute.pages.dev/entries/analysis/2026-05-18-mining-reward-exhaustion-fee-only-future/ — explainer — Carlsten et al. CCS 2016 undercutting; sufficiency vs instability.
[^16]: "Bitcoin After Block Rewards" (2026) — https://arxiv.org/html/2606.05503 — preprint — fee-only deviation threshold "just 0.17%."
[^17]: Lopp et al., "A Model for Bitcoin's Security and the Declining Block Subsidy" — https://www.lopp.net/pdf/A-model-for-Bitcoins-security-and-the-declining-block-subsidy-v1.05.pdf — research — miners pre-commit ~2yr of coins; declining subsidy "biggest threat."
[^19]: Ethereum Research, "Endgame Staking Economics" (2024) — https://ethresear.ch/t/endgame-staking-economics-a-case-for-targeting/18751 — forum (core) — nominal vs real yield; dilution protection; real yield → MEV-only as ratio→1.
[^20]: Coin Metrics, "Understanding Staking Yields" (2024) — https://coinmetrics.substack.com/p/state-of-the-network-issue-288 — data research — ETH nominal 3.08%/real 2.73%; SOL ~11.5%.
[^21]: iBuidl, "Ethereum Validator Economics 2026" — https://ibuidl.org/blog/ethereum-validator-economics-2026-20260310 — analyst — 33.4M ETH staked (~28%); Lido 28.6% vs 33.3%; stETH most-used DeFi collateral.
[^22]: Urban Jermann, "Ethereum Issuance" (Wharton) — https://finance.wharton.upenn.edu/~jermann/ETHIssuanceSSRNv2.pdf — academic — inverse-√ curve "consistent with optimal Ramsey policy"; security-vs-dilution.
[^23]: pa7x1, ethereum-issuance model + "Shape of Issuance Curves" — https://ethresear.ch/t/the-shape-of-issuance-curves-to-come/20405 — forum/model — negative-real-yield break-points (solo ~70M, inst. ~90M, LST ~100M ETH).
[^24]: Ethereum Magicians, "Electra: Issuance Curve Adjustment" — https://ethereum-magicians.org/t/electra-issuance-curve-adjustment-proposal/18825 — governance forum — lower nominal curve → higher real equilibrium yield.
[^25]: ethereum.org, "PoS rewards and penalties" — https://ethereum.org/developers/docs/consensus-mechanisms/pos/rewards-and-penalties/ — official docs — slashable acts; "up to 1 ETH"; day-18 correlation penalty.
[^26]: Ben Edgington, eth2book §2.8.7 Slashing — https://eth2book.info/latest/part2/incentives/slashing/ — authoritative — `min(B, 3SB/T)`; PROPORTIONAL_SLASHING_MULTIPLIER=3 since Bellatrix.
[^27]: EIP-7716 + dapplion anti-correlation FAQ — https://eips.ethereum.org/EIPS/eip-7716 — EIP (draft) — extends anti-correlation to missed attestations; "reduce economies of scale."
[^28]: Vitalik Buterin / Ethereum Research, "anti-correlation incentives" (2024) — https://ethresear.ch/t/supporting-decentralized-staking-through-more-anti-correlation-incentives/19116/1 — forum (primary) — correlated faults reveal single large actors.
[^29]: Lido Governance, "Should Lido be limited to a fixed % of stake?" (2022) — https://research.lido.fi/t/should-lido-on-ethereum-be-limited-to-some-fixed-of-stake/2225/1 — primary governance — stETH mechanics; ⅓-threshold coercion risk; "winner-take-most."
[^30]: "SoK: Liquid Staking Tokens and Restaking" (2024) — https://arxiv.org/html/2404.00644v3 — academic SoK — ~$58B LST TVL; Lido ~31% via ~29 operators; shared-software risk.
[^31]: Merkle Labs, "Beyond 33%: Lido's Growing Stake" (2023) — https://merklelabs.substack.com/p/beyond-33-implications-of-lidos-growing — analyst — 33/50/66% threshold effects.
[^32]: DataField.Dev, "Lido's 33% Problem" (2026) — https://datafield.dev/blockchain-cryptocurrency/part-03/chapter-16/case-study-01.html — explainer — June 2023 vote ~99.8% against self-limiting.
[^33]: Coin Metrics, "Risks of Liquid Staking Derivatives" (2024) — https://coinmetrics.substack.com/p/state-of-the-network-issue-250 — data research — DVT mitigation; June-2022 stETH depeg to $0.935, ~$1.18B liquidations.
[^35]: EigenLayer Whitepaper — https://docs.eigenlayer.xyz/assets/files/EigenLayer_WhitePaper-88c47923ca0319870c611decd6e562ad.pdf — primary — restaking via extra slashing; pooled security; cost-of-corruption.
[^36]: EigenLayer, "EIGEN: Universal Intersubjective Work Token" (2024) — https://blog.eigencloud.xyz/eigen/ — primary — intersubjective fork-and-slash; ETH+EIGEN dual-staking.
[^37]: Consensys, "EigenLayer: A Restaking Primitive" — https://consensys.io/blog/eigenlayer-a-restaking-primitive — explainer — rehypothecation; cost-of-corruption rises as pools unify.
[^38]: Vitalik Buterin, "Why Proof of Stake (Nov 2020)" — https://vitalik.eth.limo/general/2020/11/06/pos2020.html — essay (primary) — per-$1/day attack costs GPU ~$0.26 / ASIC ~$486 / PoS ~$2,189; 5–20×; slashing vs brick-the-ASICs.
[^39]: Vitalik Buterin, "How I Learned to Love Weak Subjectivity" (2014) — https://blog.ethereum.org/2014/11/25/proof-stake-learned-love-weak-subjectivity — essay (primary) — weak-subjectivity definition; trusted recent checkpoint.
[^40]: Vitalik Buterin, "Proof of Stake FAQ" (2017) — https://vitalik.eth.limo/general/2017/12/31/pos_faq.html — FAQ (primary) — selfish-mining 25/33% bounds; economic finality; reduced alignment if non-base tokens stake.
[^41]: Jon Charbonneau (DBA), "L1 & L2 Token Value Capture" (2024) — https://dba.xyz/l1-l2-token-value-capture/ — analyst — PoS issuance = "internal tax"; trillion-$ PoS "squishes real yield to ~0"; BTC counterpoint; restaking-risk framing.
[^42]: Vitalik Buterin, "SchellingCoin" (2014) — https://blog.ethereum.org/2014/03/28/schellingcoin-a-minimal-trust-universal-data-feed — essay (primary) — "truth is the most powerful Schelling point"; reward near-median.
[^43]: Ethereum Blog, "Advanced Contract Programming: SchellingCoin" (2014) — https://blog.ethereum.org/2014/06/30/advanced-contract-programming-example-schellingcoin — primary — commit-reveal; 25th–75th-percentile reward; PoW isomorphism.
[^44]: Ethereum Blog, "The Subjectivity/Exploitability Tradeoff" (2015) — https://blog.ethereum.org/2015/02/14/subjectivity-exploitability-tradeoff — primary — Truthcoin contentious-vote penalty; subjectivity reduces exploitability.
[^45]: Mike Goldin, "Token-Curated Registries 1.0" (2017) — https://medium.com/@ilovebagels/token-curated-registries-1-0-61a232f8dac7 — primary spec — deposit/challenge/commit-reveal; tactical-vs-strategic incentive.
[^46]: Mike Goldin, "TCR 1.1, 2.0" (2017) — https://medium.com/@ilovebagels/token-curated-registries-1-1-2-0-tcrs-new-theory-and-dev-updates-34c9f079f33d — primary — variable "trust pools"; Schelling points around pool sizes.
[^47]: Vitalik Buterin, "The P + ε Attack" (2015) — https://blog.ethereum.org/2015/01/28/p-epsilon-attack — essay (primary) — conditional bribe → defection dominant → zero-cost takeover.
[^48]: Ethereum Research, "Schelling Coin Bribery Bidding Wars" (2020) — https://ethresear.ch/t/schelling-coin-bribery-bidding-wars/7484 — forum — costless bribes → bidding war independent of stake; security rises with bidder uncertainty.
[^49]: Jelle Gerbrandy, "Incentive alignment in TCRs" (2018) — https://medium.com/paratii/incentive-alignment-in-token-curated-registries-4d6e41652a9b — critique — Schelling points apply only when communication is impossible; security rests on external token value.
[^50]: Luu et al., "Demystifying Incentives in the Consensus Computer" (CCS 2015) — https://eprint.iacr.org/2015/702 — peer-reviewed — verifier's dilemma; "rational miners accept unvalidated blockchains."
[^51]: (same paper, PDF) — https://www.comp.nus.edu.sg/~prateeks/papers/VeriEther.pdf — peer-reviewed — 30%-of-interval verify cost lets skippers win; incentive incompatibility.
[^52]: Ferreira Torres et al., "Sluggish Mining" (FC 2019) — https://link.springer.com/chapter/10.1007/978-3-030-43725-1_6 — peer-reviewed — slow contracts give a mining edge; gas limit mitigates not eliminates.
[^53]: "Data-Driven Analysis of the Ethereum Verifier's Dilemma" (DSN 2020) — https://ar5iv.labs.arxiv.org/html/2004.12768 — peer-reviewed — "often rational not to verify"; parallelization + active invalid-block insertion.
[^54]: Ethereum Research, "alleviate the empty block problem" (2019) — https://ethresear.ch/t/a-proposal-to-alleviate-the-empty-block-problem/6191 — forum — SPV/empty-block lazy validation; Ethereum state-root forces prev-block validation.
