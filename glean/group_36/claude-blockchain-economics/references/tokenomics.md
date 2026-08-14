<!-- Provenance: reference under the `blockchain-economics` standalone skill. Created 2026-06-16 via /dr deep-research. Synthesized from web sources; all load-bearing/volatile claims carry [^n] citations to the ## References section. -->

# Tokenomics

verified-as-of: 2026-06-16 (volatile: net-issuance direction, ultrasound-money status, unlock calendars, SEC token-taxonomy status — re-verify before relying)

The economic design of a token: how it is **issued**, how supply is **removed**, how it is **distributed** and **vested**, how it is **classified**, and why "usage" so often fails to become **value**. Confidence tags: **[FACT]** = 3+ independent sources agree; **[QUALIFIED]** = 2 sources; **[TENTATIVE]** = single-source/contested. Mechanism-level only — protocol issuance code defers to the per-chain skills; price prediction is out of scope.

## Contents

- Supply & issuance/emission schedules
- Burns/sinks, net issuance & the ultrasound-money debate
- Vesting & unlock cliffs
- Token distribution (fair-launch vs VC, airdrops, sybil)
- Token taxonomy (utility/governance/security) + Howey
- The token-velocity problem (MV=PQ)
- References

## Supply & issuance/emission schedules

**Two orthogonal axes.** *Supply cap* (hard-capped like Bitcoin's 21M, or uncapped like Ethereum/Solana) and *issuance schedule* (fixed-rate, disinflationary, or discretionary) are independent. A coin can be uncapped yet disinflationary (Solana), or capped via a decaying schedule (Bitcoin).[^1][^5] **[FACT]**

- **Inflationary** — circulating supply grows materially; typically uncapped and/or perpetual issuance to fund mining/staking/liquidity incentives.[^5] **[FACT]**
- **Disinflationary** — supply still grows but the *rate* declines toward (never reaching) zero; the technical condition is that the per-period reward never *increases*.[^3] **[FACT]**
- **Deflationary** — net supply shrinks (a sink must exceed issuance).[^5] **[FACT]**

**Bitcoin halving as a disinflation curve (economic abstraction).** The block subsidy starts at 50 BTC and halves every 210,000 blocks (~4 yrs), the cumulative supply approaching but never exceeding ~21M (finite geometric series), truncating to zero after 33 halvings (~2140). Era schedule: ~50% of all BTC by 2012, 75% by 2016, 87.5% by 2020, 93.75% by 2024 (3.125 BTC subsidy).[^1][^2][^6] The economically salient point is a **pre-committed, mechanical, monotonically-decreasing issuance rate** — defer the satoshi-arithmetic/subsidy code to `bitcoin-protocol-expert`.[^1] **[FACT]**

**Tail emission ≠ inflationary (contested).** Monero/Grin use a perpetual fixed per-block reward. Peter Todd argues a constant nominal emission is *not* long-run inflationary: because coins are continuously lost, a fixed flow converges to a *stable* supply (∝ emission rate × loss probability) with the true inflation rate trending toward zero (Monero entered tail emission at ~0.9%, declining).[^4a] The relevant variable is the long-run *flow*, not the presence of a numeric cap.[^3][^4a] **[QUALIFIED]** The hard-money counter-position holds that *discretionary* issuance tends to inflate under political incentives, so a credible hard cap is "the only durable defense against debasement" — a value-laden position, not a fact.[^1] **[QUALIFIED — treat as a position.]**

## Burns/sinks, net issuance & the ultrasound-money debate

A **sink** permanently removes tokens (or locks demand) faster than they re-enter float. Two archetypes:

**(a) EIP-1559 base-fee burn (economic view).** The protocol-set base fee is *burned*, not paid to validators; only the tip goes to the proposer.[^7][^8] Economically this makes the fee a demand-proportional, usage-driven *sink* and ensures only the native asset can pay for blockspace ("cementing the economic value of ETH").[^7] Paradigm: burning the base fee removes the asset "no matter who pays," which makes a perpetual subsidy "much more palatable" because the burn counterbalances issuance.[^8] Defer the fee-market mechanics to `references/fee-markets.md`. **[FACT]**

**(b) Buyback-and-burn.** A protocol uses revenue to repurchase its token on the market and send it to a burn address (analogous to a stock buyback).[^11]
- **BNB Auto-Burn** — quarterly, formula-driven, targeting a 100M total-supply floor; Delphi found **no discernible burn↔price correlation**, possibly because treasury-sourced (vs open-market) buybacks neutralize buy pressure.[^11][^11a] **[QUALIFIED]**
- **MakerDAO** — MKR is *burned* via surplus auctions above a buffer and *minted* (diluting holders) via debt auctions to recapitalize after liquidations (e.g. Black Thursday 2020) — a **bidirectional** supply mechanism; a peer-reviewed study found MKR appreciates on surplus-burn and devalues on liquidation-mint.[^11b][^11c] **[FACT]**
- **Sink taxonomy** — the post-purchase fate can be (1) burned to null (irreversible), (2) accumulated on a fund balance (off-float but recoverable), or (3) recycled into incentives.[^11d] **[TENTATIVE — single source.]**

**Net issuance = issuance − burn; the "ultrasound money" thesis.** Post-Merge, ETH issuance fell from ~13,000 ETH/day (PoW) to ~1,700 ETH/day (PoS); if the EIP-1559 burn exceeds issuance, supply *shrinks*, making ETH "harder than gold" ("ultra sound money"). The design intent: pay for security with *predictable* issuance, let *volatile* fee burn act as "un-issuance."[^7][^7c][^7d] **[FACT]**

**Why it's contested (the thesis has materially broken — strong negation):**
- **Conditional, not guaranteed.** ETH is deflationary only when burn > issuance; quiet periods always net-inflate. As of 2026 ETH has been slightly net-*inflationary*.[^7b][^9] **[FACT]**
- **The Dencun/EIP-4844 break (Mar 2024).** Cheap L2 "blobs" cut rollup costs 10–100× and shifted activity off L1, **collapsing the base-layer burn** (daily burn fell to ~50–70 ETH; ~95% drop on some measures); CryptoQuant called the narrative "dead," and a core dev said "the meme of 'ultra sound money' is broken." Circulating supply has grown ~950k ETH since the Merge.[^9][^10][^12][^13] **[FACT]**
- **Structural critique.** The rollup-centric roadmap *structurally* weakens ETH value capture because L2s have their own tokens and push fees low; on **2026-02-03 Vitalik said rollups as "branded shards" of Ethereum "no longer makes sense."** EF response includes EIP-7918 (a blob-fee *floor* so the burn doesn't go dormant).[^7b][^10][^13] **[QUALIFIED — fast-moving.]**
- **Steelman.** Binance Research: the inflationary reversal is *expected* during scaling (blockspace outrunning demand); ETH issuance remains <1% (far below alt-L1s); cyclical demand should restore the burn.[^7a] **[QUALIFIED]**

## Vesting & unlock cliffs

**Definitions.** *Vesting* governs when allocated-but-locked tokens become transferable. *Cliff* = a lock after which a tranche unlocks at once (often then linear); *linear/graded* = gradual release with no single large event.[^15][^17][^19] **[FACT]**

**Sell-pressure findings (large quant studies converge):**
- **≈90% of token unlocks create negative price pressure**, largely regardless of size or type (Keyrock, n=16,000+; KuCoin corroboration).[^18][^16] **[FACT]**
- **Cliffs concentrate sell pressure; linear distributes it.** Echo Zero: cliff unlocks averaged ~18.3% drawdown vs ~7.2% for linear in the 7-day window around major unlocks; linear schedules show ~40–50% less unlock-period volatility.[^19] **[QUALIFIED — single-aggregator figures, directionally consistent.]**
- **The "1% rule" (6th Man Ventures, n=5,000).** Unlocks raising circulating supply by **0–1% had no meaningful price relationship**; the advice is "rethink the cliff" and prefer daily/weekly small unlocks. At >10% of supply, negative-impact probability ≈78%; at >20%, ≈91%.[^16][^19] **[FACT — multiple quant studies converge.]**
- **Impact starts *before* the unlock** (~30 days prior, as holders pre-empt) and stabilizes ~14 days after.[^18] **[QUALIFIED]**
- **Recipient type — literature *disagrees*.** Investor/VC unlocks correlate with large drops (seeking liquidity); some studies find team unlocks muted (lockup extensions/alignment), but Keyrock found team unlocks triggered the *worst* crashes (~−25%). **[CONTRADICTION — preserve both.]**[^18][^19]
- **Regime dominates.** Identical-size unlocks produced ~2–3× larger drawdowns in bear vs bull conditions.[^19] **[TENTATIVE — single source.]**

**Design tension.** Founder *retention* signals commitment (Fuchs & Momtaz 2024: every 10% reduction in tokens sold to the primary market ≈ 3% more capital raised), but that commitment's *cost is paid by the market at the unlock* — a 12-month cliff's sell-pressure-concentration index ≈3× continuous vesting.[^15] There is **no empirical basis** for the ubiquitous 12-month-cliff convention over 6/18-month alternatives — it's a norm partly driven by US compliance practice.[^15][^16] **[TENTATIVE.]** 2026 convention (volatile): team 12-mo cliff + 3–4yr linear; investors 6–12mo + 2–3yr; TGE unlock ~5–15% (conservative) vs 30–50% (aggressive); zero-cliff team/investor allocations are red flags.[^17] **[TENTATIVE — vendor data.]**

## Token distribution (fair-launch vs VC, airdrops, sybil)

**Four canonical buckets:** team/founders, investors (VC/private/public), community (airdrops, liquidity mining, ecosystem incentives), and treasury/foundation.[^20][^22][^25] **[FACT]**

- **Fair launch** — no pre-sale, no insider pre-allocation; all tokens earned via open participation. **Bitcoin is canonical** (all BTC mined, no premine); memecoins are the modern archetype. Trade-off: maximum credible-neutrality but **no project funding** and sybil-vulnerable.[^20][^23] **[FACT]**
- **VC/insider allocation** — pre-launch sale to investors; capital + expertise but concentrated power and unlock overhang. Typical: treasury 30–50%, team 15–25%, investors 15–25%, ecosystem 10–20%, public 5–15%.[^20][^24] **[FACT]**
- **Premine** = tokens allocated to insiders before public availability (the defining non-fair-launch feature).[^23] **[FACT]**
- **Canonical template — Uniswap UNI:** 60% community (incl. 15% retroactive airdrop), 21.51% team (4-yr vest), 17.80% investors (4-yr vest), 0.069% advisors.[^20] **[FACT]**

**Airdrops & the dump/sybil problem (strong data):**
- **Mechanism evolution:** flat retroactive (UNI 400/swapper; ENS) → criteria-based multi-round (Optimism) → **points-farming pre-TGE** (now dominant; Hyperliquid's ~$2.6B the largest).[^21][^24][^25] **[FACT]**
- **Recipients overwhelmingly dump.** UNI: >75% of wallets sold all within 7 days, ~93% within ~2 years, only ~1% increased their position, ~98% never voted. Generalized: >80% of airdropped tokens sold within 30 days (Arbitrum/Optimism); even criteria-based OP saw ~30% immediate sell-off — the **"mercenary capital"** problem.[^26][^27] **[FACT]**
- **Sybil attacks** — single entities farm criteria across hundreds of wallets, diluting genuine users; Optimism flagged/recovered tokens from ~17,000 sybil wallets; LayerZero's anti-sybil campaign was "prohibitively expensive."[^21][^25][^27] **[FACT]**
- **The "cobra effect"** — airdrops meant to bootstrap can backfire (most of 14 studied post-airdrop projects underperformed); mitigations like Jupiter's Active Staking Rewards reward stakers to *reduce* unlock sell pressure.[^27a] **[QUALIFIED]**

## Token taxonomy (utility/governance/security) + Howey

- **Utility token** — access to a product/service or participation right; generally no ownership/dividend.[^28][^28a] **[FACT]**
- **Governance token** — voting/oversight rights; academically a **subset of utility tokens** (the "utility" is governance participation); faces voter-participation challenges.[^28a] **[QUALIFIED]**
- **Security token** — blockchain-native representation of a traditional security/RWA; ownership/profit/dividend/asset claims; regulated as a security.[^28][^29] **[FACT]**
- **Caveat:** these are economic/functional categories, not legally dispositive — "marketing labels are secondary to economic reality."[^28] **[FACT]**

**Howey (high-level).** From *SEC v. W.J. Howey Co.*, 328 U.S. 293 (1946): an "investment contract" exists where there is (1) an investment of money, (2) in a common enterprise, (3) a reasonable expectation of profits, (4) derived from the efforts of others. Howey is flexible, substance-over-form, case-by-case, and remains binding on crypto.[^29][^30][^31] **[FACT]**

**Volatile — the 2025–2026 SEC token-taxonomy development.** Through 2025 the SEC largely "regulated by enforcement." In 2026 it issued an Interpretive Release (with CFTC complementarity) classifying crypto assets into ~five categories — **Digital Commodities** (not securities; "names names" incl. BTC/ETH/SOL/ADA), **Digital Collectibles**, **Digital Tools**, **Stablecoins** (GENIUS-Act payment stablecoins), and **Digital Securities** (securities). Key shifts: a non-security asset can still be *sold as part of* an investment contract; **investment-contract status is not perpetual** (rejecting "security forever"); protocol mining/staking/wrapping and certain airdrops may *not* be securities offerings; a safe-harbor follow-up is anticipated.[^29][^30][^31][^32] **[QUALIFIED — fast-moving, US-specific; verify status; defer detailed securities-law analysis to a securities-law resource.]**

## The token-velocity problem (MV=PQ)

**The frame.** Adapted from Fisher's equation of exchange: **M** = monetary base/token supply, **V** = velocity (turns per period), **P** = price level of the provisioned resource, **Q** = quantity provisioned; the supportable network value is **M = PQ / V**, and token value ≈ M / circulating supply.[^34][^35][^36] **[FACT]**

**Velocity as a headwind.** Multicoin's Kyle Samani ("Understanding Token Velocity," 2017) is canonical: holding M, P, Q fixed, **price and velocity are inversely proportional — doubling velocity halves price.** A proprietary payment token nobody wants to *hold* has velocity that grows with transaction volume, so "volume could grow a million-fold and network value remain constant" — the "usage without value capture" pattern.[^34][^36] (An asset with velocity 0 trades at a *premium* — the liquidity premium — so some velocity is healthy; the problem is *excess* velocity.)[^36] **[FACT]**

**Value-capture / velocity-reduction mechanisms (Multicoin's five):**[^36] (1) **profit-share / buy-and-burn** (falling price → rising yield → buy-and-hold); (2) **staking that locks the asset** (incl. PoS; e.g. Augur REP, Livepeer LPT must be staked to do paid work)[^34b]; (3) **balances / working capital** (require holding to use); (4) **gamification to encourage holding**; (5) **become a store of value** (the hardest; "defines Bitcoin today"). A sixth from the broader literature is **burn-and-mint equilibrium (BME)** — burn to consume, mint to reward producers — decoupling required-holding from velocity.[^34] **[FACT for the five; QUALIFIED for BME attribution.]**

**Burniske's valuation use.** "Cryptoasset Valuations" (2017): compute each year's **Current Utility Value** via MV=PQ, discount future utility value (often 30–50%) to a present value; market price = CUV + **Discounted Expected Utility Value** (speculative) — the "crypto J-curve." He expects crypto velocity to run much higher than traditional assets, depressing CUV.[^35][^35a] **[FACT]**

**Why the frame is contested (negation):**
- **Buterin calls the model brittle.** "On Medium-of-Exchange Token Valuations" (2017): the pure MoE valuation is "ultimately quite brittle" with "an unavoidable risk of collapsing at any time," sustained only by temporary zero-holding-cost equilibria.[^37] **[FACT]**
- **Holding-time math wrong.** A 2024 paper (arXiv:2403.04914) shows Buterin's H = 1/V "has been proven incorrect," and the EoE assumes a stable V that is in fact highly volatile.[^38a] **[QUALIFIED]**
- **"Not really a problem" for payment networks** (Zochowski/Logos): the problem rests on unrealistically high V; PQ reflects *total economic activity*, not just validation cost.[^38b] **[TENTATIVE — single contrarian source.]**
- **Austrian rejection** (Mises): the quantity theory is misapplied — "the value of money determines the flow of spending, not the other way around"; *saleableness* matters more than velocity.[^38c] **[TENTATIVE — heterodox.]**

> **Net read:** the velocity *framework* is the field's dominant qualitative lens for why pure utility/payment tokens struggle to accrue value, and the sink/value-capture mechanisms it motivates are real and widely deployed. The *quantitative* MV=PQ model is empirically shaky — a directional heuristic, never a pricing engine, and **never trading advice**.

## Child concepts (future research)

1. Burn-and-mint equilibrium & dual-token/work-token models (Factom/Helium; Multicoin "New Models for Utility Tokens").
2. veTokenomics & vote-escrow lockups (Curve veCRV, "Curve Wars," bribe markets) — see `references/governance-economics.md`.
3. FDV vs circulating-supply dynamics & the low-float/high-FDV critique.
4. Empirical token-valuation models beyond MV=PQ (Metcalfe/NVT, stock-to-flow and its debunking).
5. Staking economics & real vs inflationary yield — see `references/cryptoeconomics-and-security.md`.

## References

[^1]: Fixed supply vs adjustable money — Bitcoin Institute — https://bitcoin-institute.pages.dev/entries/analysis/2008-10-31-fixed-supply-vs-adjustable-money/ — blog — cap and schedule are orthogonal axes.
[^2]: Bitcoin monetary design — Bitcoin Institute — https://bitcoin-institute.pages.dev/entries/design/2009-01-03-bitcoin-monetary-design/ — blog — halving era subsidy/cumulative-supply table; 33 halvings to zero.
[^3]: A case for soft total supply — John Tromp — https://tromp.github.io/blog/2020/12/20/soft-supply — blog — disinflation definition; inflation→0 if flow never increases.
[^4a]: Surprisingly, Tail Emission Is Not Inflationary — Peter Todd — https://petertodd.org/2022/surprisingly-tail-emission-is-not-inflationary — essay — fixed perpetual reward → stable supply given loss; Monero ~0.9%.
[^5]: Inflationary vs Deflationary Crypto — VALR — https://blog.valr.com/blog/inflationary-vs-deflationary-crypto-an-explainer — blog — inflationary=uncapped/perpetual; deflationary=cap+halvings+burns.
[^6]: Bitcoin's Supply Schedule — Lightspark — https://lightspark.com/glossary/bitcoin-supply-schedule — glossary — 21M cap, 210k-block halving, 3.125 BTC post-2024, ~2140.
[^7]: EIP-1559 — Ethereum Improvement Proposals — https://eips.ethereum.org/EIPS/eip-1559 — docs (primary) — base fee burned, tip to proposer; burn counterbalances issuance; cements asset value.
[^7a]: The ETH Value Debate — Binance Research — https://public.bnbstatic.com/static/files/research/the-eth-value-debate.pdf — report — Dencun reshaped L1 fees; supply inflationary; issuance <1% defense.
[^7b]: What Is Ultrasound Money? — Cryptothreads — https://cryptothreads.io/learn/what-is-ultrasound-money/ — blog — ~1,700 ETH/day issuance threshold; post-Dencun burn 50–70 ETH; EIP-7918 floor; slightly inflationary 2026.
[^7c]: ETH is Ultra Sound Money — Bankless — https://www.bankless.com/eth-is-ultra-sound-money-market-monday — essay — origin; burn = "un-issuance."
[^7d]: Ultra Sound Money — Devcon wiki — https://www.devcon.wiki/Cryptoeconomics/Ultra%20Sound%20Money — wiki — net = burn − issuance.
[^8]: Analysis of EIP-1559 — Paradigm — https://www.paradigm.xyz/2020/06/analysis-of-eip-1559 — essay — burning base fee removes asset "no matter who pays"; makes perpetual subsidy palatable.
[^9]: Ethereum Scaling Has Damaged Its Tokenomics — Unchained — https://unchainedcrypto.com/ethereum-scaling-with-l2s-has-damaged-its-tokenomics-is-it-possible-to-fix-it/ — blog — CryptoQuant "narrative dead"; burn tied to revenue balance.
[^10]: The End of 'Ultrasound Money' — ForkLog — https://forklog.com/en/the-end-of-ultrasound-money-why-ethereum-is-losing-developers-and-whale-support/ — blog — Dencun turning point; Vitalik 2026-02-03 "branded shards no longer makes sense."
[^11]: Value Accrual Mechanisms — Delphi Digital — https://members.delphidigital.io/reports/more-than-just-governance-unpacking-value-accrual-mechanisms — report — buyback-and-burn; BNB Auto-Burn; MakerDAO SBE; no BNB burn↔price correlation.
[^11a]: What Is BNB Auto-Burn? — Binance Academy — https://www.binance.com/en/academy/articles/what-is-bnb-auto-burn — docs — 100M-supply target (primary issuer doc).
[^11b]: MKR Module — MakerDAO docs — https://docs.makerdao.com/smart-contract-modules/mkr-module — docs (primary) — MKR mint/burn; surplus/debt auctions.
[^11c]: Fundamentals of the MakerDAO Governance Token — Dagstuhl Tokenomics 2021 — https://drops.dagstuhl.de/storage/01oasics/oasics-vol097-tokenomics2021/OASIcs.Tokenomics.2021.11/OASIcs.Tokenomics.2021.11.pdf — peer-reviewed — surplus→burn→price↑; liquidations→mint→price↓.
[^11d]: Buyback Engineering — Giants/Tokenomics Lab — https://giantslabs.pro/models/buyback-engineering/ — blog — sink taxonomy (burn/accumulate/recycle); Hyperliquid burns MEV/priority fees.
[^12]: ETH is not ultrasound money anymore — DL News — https://www.dlnews.com/articles/defi/ethereum-scaling-means-eth-is-not-ultrasound-money-anymore/ — blog — Van Loon "meme broken"; Dencun cause.
[^13]: ETH Has a Narrative Problem — Blockhead — https://www.blockhead.co/2026/06/05/eth-has-a-narrative-problem/ — blog — structural critique; ~950k ETH supply growth since Merge.
[^15]: Cliff Periods Concentrate Sell Pressure — Token Strategy — https://token-strategy.com/blog/vesting-cliff-concentration-sell-pressure-research — blog — cites Fuchs & Momtaz 2024; cliff concentration ~3×; no basis for 12-mo norm.
[^16]: We Analyzed 5,000 Token Unlocks — 6th Man Ventures — https://6thman.ventures/writing/token-unlocks — quant study — the 1% rule; rethink the cliff; prefer daily/weekly.
[^17]: Token Vesting Benchmarks — Streamflow — https://streamflow.finance/blog/token-vesting-benchmark-report — vendor — 2026 convention table; cliff concentrates vs linear distributes; volatile unlock-$ figures.
[^18]: 16,000+ Token Unlocks — Keyrock — https://keyrock.com/from-locked-to-liquidity-what-16000-token-unlocks-teach-us/ — quant study — ~90% negative; impact starts ~30d before; team unlocks worst (−25%).
[^19]: Token Vesting Schedule Analysis — Echo Zero — https://blog.echozero.app/article/token-vesting-schedule-analysis-impact-on-price-action — blog — cliff 18.3% vs linear 7.2% drawdown; >10% 78% / >20% 91% impact probability.
[^20]: Fair Launch vs Investor Allocation — Tokenization Governance — https://tokenizationgovernance.com/token-design/governance-token-distribution/ — blog — fair-launch (Bitcoin canonical); typical allocation %; UNI breakdown.
[^21]: Do Airdrops Hurt More Than They Help? — Tiger Research — https://reports.tiger-research.com/p/do-airdrops-hurt-more-than-they-help-eng — report — UNI origin; sybil farming; LayerZero anti-sybil cost.
[^22]: Token Allocation in Tokenomics — Tokenomics.com — https://tokenomics.com/articles/token-allocation — blog — insider-pool vs 1:1 equity conversion; "fails from misaligned incentives."
[^23]: Token Distribution Guide 2026 — TokenMinds — https://tokenminds.co/blog/token-distribution — blog — fair launch = no insider pre-allocation; premine contrast.
[^24]: Token Allocation as a Supply Model — Giants/Tokenomics Lab — https://giantslabs.pro/models/allocation/ — blog — fair-launch vs sale table; Hyperliquid $2.6B; MC/FDV<10% trust loss.
[^25]: Tokenomics Deep Dive — Binance Research — https://research.binance.com/static/pdf/Tokenomics_Deep_Dive_Stefan_Piech_Shivam_Sharma.pdf — report — recipient-party table; Optimism ~17k sybils recovered.
[^26]: The Uniswap Airdrop — Lessons — Dune — https://dune.com/blog/uni-airdrop-analysis — on-chain data — 75% sold in 7d, ~93% sold all, ~1% increased, ~98% no governance.
[^27]: The Hidden Cost of Airdrops — ChainScore — https://chainscorelabs.com/blog/the-hidden-cost-of-airdrops-creating-mercenary-token-holders — blog — mercenary capital; >80% sold in 30d.
[^27a]: Airdrops and the Cobra Effect — Xangle — https://xangle.io/en/research/detail/2083 — report — 14 airdrop projects underperformed; Jupiter ASR to cut unlock sell pressure.
[^28]: Utility vs Governance vs Security Tokens — Cryip — https://cryip.co/utility-tokens-vs-governance-tokens-vs-security-tokens/ — blog — functional definitions; Howey applied; labels secondary.
[^28a]: Token classification — Cardozo Law Review — https://larc.cardozo.yu.edu/cgi/viewcontent.cgi?article=2412&context=clr — law review — governance tokens a subset of utility; Howey elements.
[^29]: SEC Fact Sheet — Application of Federal Securities Laws to Crypto Assets — https://www.sec.gov/files/33-11412-fact-sheet.pdf — docs (regulator) — 2026 token taxonomy; mining/staking/wrapping/airdrops; ends "regulating by enforcement."
[^30]: SEC Issues Interpretive Guidance on Crypto Asset Classification — Orrick — https://www.orrick.com/en/Insights/2026/04/SEC-Issues-Interpretive-Guidance-on-Crypto-Asset-Classification — law firm — Howey elements; status not perpetual; rejects "security forever."
[^31]: Howey's Cryptonite — K&L Gates — https://www.klgates.com/Howeys-Cryptonite-A-Deep-Dive-on-Digital-Asset-Classification-4-20-2026 — law firm — five categories "names names" (BTC/ETH/SOL/ADA digital commodities).
[^32]: SEC's Crypto Framework: The New Token Taxonomy — Dechert — https://www.dechert.com/knowledge/onpoint/2026/3/sec-s-crypto-framework--the-new-token-taxonomy.html — law firm — five categories; CFTC joins; supersedes 2019 Framework; safe-harbor anticipated.
[^34]: Token Velocity: Why High Velocity Destroys Value — Tokenomics.net — https://www.tokenomics.net/blog/token-velocity/ — blog — MV=PQ; doubling V halves P; "usage without value capture."
[^34b]: On Value Capture at Layers 1 and 2 — Multicoin (Samani) — https://multicoin.capital/2019/03/14/on-value-capture-at-layers-1-and-2.txt — essay (primary) — stateful vs stateless; Livepeer LPT work-requirement value.
[^35]: Cryptoasset Valuations — Chris Burniske — https://medium.com/@cburniske/cryptoasset-valuations-ac83479ffca7 — essay (primary) — M=PQ/V; CUV; discount 30–50%.
[^35a]: The Crypto J-Curve — Chris Burniske — https://medium.com/@cburniske/the-crypto-j-curve-be5fdddafa26 — essay — CUV + DEUV decomposition.
[^36]: Understanding Token Velocity — Multicoin (Samani) — https://multicoin.capital/2017/12/08/understanding-token-velocity/ — essay (canonical) — velocity, liquidity premium, the 5 value-capture mechanisms.
[^37]: On Medium-of-Exchange Token Valuations — Vitalik Buterin — https://vitalik.eth.limo/general/2017/10/17/moe.html — essay (primary) — MoE valuation "brittle," "risk of collapsing at any time."
[^38a]: Refining the Equation of Exchange for token valuation — arXiv:2403.04914 — https://arxiv.org/pdf/2403.04914 — paper — Buterin's H=1/V "proven incorrect"; V is volatile.
[^38b]: Debunking the Velocity "Problem" — Zochowski/Logos — https://medium.com/logos-network/cryptoasset-valuation-2-the-velocity-problem-8bbb4111c9c7 — essay — PQ = total economic activity, not validation cost.
[^38c]: Crypto Velocity Theorists & Austrian Economics — Mises Institute — https://mises.org/power-market/what-crypto-token-velocity-theorists-can-learn-austrian-economics — essay — Austrian rejection of MV=PQ; saleableness > velocity.
