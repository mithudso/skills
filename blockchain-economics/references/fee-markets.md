<!-- Provenance: reference under the `blockchain-economics` standalone skill. Created 2026-06-16 via /dr deep-research. Synthesized from web sources; all load-bearing/volatile claims carry [^n] citations to the ## References section. -->

# Fee Markets — Transaction-Fee Mechanism Design

verified-as-of: 2026-06-16 (volatile: blob target/max parameters across forks, Solana SIMD-0110 status — re-verify before relying)

The **economic mechanism design** of pricing blockspace. Per-opcode gas accounting defers to `ethereum-protocol-expert`; the MEV/PBS supply chain is the sibling `references/mev.md`. Confidence tags: **[FACT]** = primary spec + Roughgarden + independent corroboration; **[QUALIFIED]** = 2 sources / active research; **[TENTATIVE]** = single-source/in-flux.

## Contents

- EIP-1559: base-fee dynamics, the burn, why it's not a first-price auction
- Priority fees (the auction at the margin)
- The blob fee market (EIP-4844)
- Solana localized (per-account) fee markets
- First-price vs second-price/VCG and the IC tradeoffs
- Disconfirming notes
- References

## EIP-1559: base-fee dynamics, the burn, why it's not a first-price auction

**The mechanism (live London hard fork, Aug 5 2021).** EIP-1559 replaced the legacy first-price (pay-as-bid) gas auction with a **posted-price / dynamic-pricing** mechanism: each block carries a protocol-computed `baseFeePerGas` — a **reserve price** every included tx must pay per gas unit, **burned** (not paid to the proposer) — plus a **priority fee (tip)** to the proposer. Typed (EIP-2718 type-2) txs set `maxFeePerGas` (total ceiling) and `maxPriorityFeePerGas` (tip ceiling); the sender pays base fee + `min(maxPriorityFeePerGas, maxFee − baseFee)`, the rest refunded.[^1] **[FACT]**

**Base-fee adjustment — exact parameters (from spec):**[^2]
- `ELASTICITY_MULTIPLIER = 2` → `gas_limit = gas_target × 2`, so the chain targets **50% block fullness**.
- `BASE_FEE_MAX_CHANGE_DENOMINATOR = 8` → base fee changes ≤ **1/8 = 12.5% per block** (100%-full → +12.5%; at target → no change; empty → −12.5%); **linear in deviation** from target between those bounds.
- Update: `base_fee_delta = parent_base_fee × (gas_used − gas_target)/gas_target / 8`.
- `INITIAL_BASE_FEE = 1 gwei`; integer math floors at 7 wei.
- Depends only on the **parent block** — this is what makes it a *posted price* and **decouples a user's optimal bid from other users' bids**. Convergence is tâtonnement/AIMD-style.[^1] (A max-size block then an empty block leaves the base fee at 9/8 × 7/8 ≈ 98.4% — slightly contractive.)[^4] **[FACT]**

**Why the base fee is BURNED (load-bearing, not an afterthought):**[^1][^4][^11][^12]
1. *Removes proposer manipulation* — if the proposer received the base fee, they could mine emptier/fuller blocks to push it around. Burning severs that link.
2. *Prevents the off-chain collusion that would collapse it back to a first-price auction* — Roughgarden's "Takeaway 5": if base-fee revenue went to the miner, a user wanting to pay only 60 (base fee 90) could pay 90 on-chain and be **refunded 30 off-chain**, costlessly simulating a first-price auction. Burning makes the base fee non-refundable, so "only ETH can pay for transactions." This is why "easy fee estimation and fee burning are inextricably linked through the threat of off-chain agreements."
3. *Monetary side effects (secondary, contested)* — the burn counterbalances issuance; the spec concedes it can't assert whether ETH ends up inflationary or deflationary. (See `references/tokenomics.md`.) **[FACT]**

**Why it is NOT a pure first-price auction.** Most of the fee (the base fee) is a **history-dependent posted price**, not what competing bidders bid in the same block. Roughgarden proves that *outside demand spikes*, EIP-1559 behaves as a posted-price mechanism and is therefore **DSIC** — users set `maxFee` = true value and a minimal tip, with no need to guess others' bids (unlike a first-price auction). The qualification: during a **sudden demand spike** (more willing payers than fit, base fee not yet caught up), it **reverts to a first-price auction on the tip** until the base fee climbs (each full block +12.5%, so a few minutes).[^5][^9][^10][^12] **[FACT]**

## Priority fees (the auction at the margin)

The tip **restores an auction for ordering and marginal inclusion.** Roughgarden's takeaways:[^9][^1][^12]
- **Disincentivize empty blocks** (Takeaway 7) — base fee is burned, so without a tip the proposer is indifferent to including txs; the tip pays them to fill blocks.
- **Compensate orphan/uncle and marginal inclusion cost** — canonical default tip ≈1 gwei; ~2 gwei common in practice due to MEV (MEV's deeper effect on tips/ordering → `references/mev.md`).
- **The first-price auction lives here** (Takeaway 8) — tips are user-specified so high-value txs can outbid for priority during congestion.
- **Tips go entirely to the proposer, never burned** (Takeaway 9) — burning any part would push the tip market off-chain.

So the architecture is a **hybrid**: posted-price base fee (most of the fee, burned, DSIC) + first-price tip auction at the margin (paid to proposer). Roughgarden formalizes first-price auctions, the 1559 mechanism, and the tipless mechanism as "first-price auctions with different *restricted bid spaces*."[^9] **[FACT]**

## The blob fee market (EIP-4844)

EIP-4844 ("proto-danksharding," live Dencun, Mar 13 2024) introduced **blobs** (~128 KB data packets, pruned ~18 days, used by L2 rollups), creating Ethereum's **first multidimensional fee market** — blob gas priced **independently** of execution gas.[^6][^7][^8] **[FACT]**

**Launch parameters (from spec):**[^6]
- `GAS_PER_BLOB = 2¹⁷ = 131,072`.
- `TARGET_BLOB_GAS_PER_BLOCK = 393,216` = **3 blobs** target; `MAX_BLOB_GAS_PER_BLOCK = 786,432` = **6 blobs** max (target = ½ max, mirroring 1559 elasticity-2).
- `MIN_BASE_FEE_PER_BLOB_GAS = 1` wei; `BLOB_BASE_FEE_UPDATE_FRACTION = 3,338,477`.

**Adjustment rule — exponential (EIP-1559-analogous but distinct):**[^6][^7]
- `base_fee_per_blob_gas = MIN_BASE_FEE_PER_BLOB_GAS × e^(excess_blob_gas / UPDATE_FRACTION)` (integer `fake_exponential`).
- `excess_blob_gas` accumulates `(blob_gas_used − TARGET)` each block, floored at 0.
- Update fraction chosen so max per-block change ≈ e^(target/fraction) ≈ **±12.5%** (same headline rate as 1559). **Functional form differs**: blob fee is *exponential in cumulative excess*; the EVM base fee is *linear-in-deviation* AIMD. **[FACT for params; QUALIFIED that both are "1559-like controllers," not identical.]**

**Why blobs get their own fee market.** Blob data (DA bandwidth) and EVM execution (computation) stress **different non-fungible resources**; bundling them hardcodes a relative price and blocks granular price discovery. A separate dimension lets a DA-spike on rollups *not* price out ordinary EVM txs.[^8][^13] **Structural difference: blobs have NO priority fee** — blob txs set only `max_fee_per_blob_gas`; there is only an auto-adjusted blob base fee (the tx's *execution* tip still orders it).[^8] **[FACT / QUALIFIED.]**

**Volatile — superseded by EIP-7691 (Pectra, 2025): target 6 / max 9 blobs**, `UPDATE_FRACTION_PRAGUE = 5,007,716`; this **broke the target=½max symmetry** (2:3), making the fee more responsive to under-target ("two full-blob blocks to make up for one empty one"). EIP-7840 added a per-fork blob schedule; EIP-7892 ("BPO" hardforks) and a generalized update-fraction formula are the forward path toward PeerDAS/danksharding.[^14] **[FACT for 4844+7691; TENTATIVE for post-Pectra proposals in flux.]** Documented failure mode: the **"target demand paradox"** — when target sits outside the range of effective demand, the market can't reach equilibrium (post-launch the blob base fee sat pinned at ~1 wei while rollup uptake was below target).[^8] **[QUALIFIED.]**

## Solana localized (per-account) fee markets

**The contrast.** Ethereum's base fee is a **single global** price for one fungible resource — a hotspot on one contract raises the price for *everyone*. Solana's goal is **localized (per-writable-account) fee markets** so contention on a "hot" account raises *that account's* cost without pricing out unrelated txs, enabled by **Sealevel** requiring every tx to declare the accounts it reads/writes (which also permits parallel execution).[^15][^16] **[FACT]**

**Solana fee structure today:**[^15][^18][^19]
- **Base fee:** fixed **5,000 lamports per signature**.
- **Prioritization fee:** optional, via two Compute Budget instructions — `SetComputeUnitLimit` (max CUs; default 200k/instruction, **max 1,400,000 CUs/tx**) and `SetComputeUnitPrice` (micro-lamports/CU); `prioritization_fee = ceil(cu_price × cu_limit / 1,000,000)` lamports, charged on the **requested** (not actual) CU limit, **100% to the validator** (since SIMD-0096, May 2024, which moved priority fees from 50/50-burn to 100%-validator).

**De facto vs deterministic local markets (often gotten wrong):**[^16][^17]
- Today's locality is an **emergent, non-deterministic scheduler byproduct** — a higher priority fee makes inclusion *more likely* and contested state *tends* to need higher fees, but ordering is non-deterministic; **there is no consensus-enforced per-account base fee in production.**
- **SIMD-0110 ("Exponential fee for write-lock accounts,"** Zhu/Yakovenko, Jan 2023) would make local markets *deterministic*: a bounded LRU cache (≤2048 accounts) tracks each contentious account's **EMA of CU utilization**, raising a **write-lock cost rate** above a target (initial: target ~25% of max CU, 1,000 µlamports/CU, ±1%/slot); the write-lock fee is **100% burned** (EIP-1559 burn logic, per-account). The authors note it is "like EIP-1559... but far more granular, setting per-account local base fee markets," and the bounded cache answers the resource-exhaustion objection. **Status (volatile): SIMD-0110 is inactive/closed** — Solana has *not* shipped consensus-enforced local fee markets.[^15][^17] **[FACT it's not live; TENTATIVE on what ships next.]**

**Theory bridge.** Both 4844 (resource level) and Solana (per-account state level) are **multidimensional fee markets**; the foundational treatment is Diamandis–Evans–Chitra–Angeris (arXiv:2208.07919), deriving an EIP-1559-like controller from a network-designer loss function and extending it to per-account pricing.[^8][^13] **[QUALIFIED — active research.]**

## First-price vs second-price/VCG and the IC tradeoffs

The theoretical spine is **Roughgarden's "Transaction Fee Mechanism Design for the Ethereum Blockchain"** (arXiv:2012.00854; EC'21 summary 2106.01340).[^5][^4]

**Three incentive-compatibility notions:**[^5]
1. **DSIC (user)** — truthful bidding is dominant; no need to reason about others' bids.
2. **MMIC (myopic miner)** — a profit-maximizing proposer maximizes single-block revenue by following the rule (notably: **no profit from inserting fake transactions**).
3. **OCA-proofness** — no proposer+users coalition can strike an off-chain deal that strictly increases joint utility.

**Why naive first-price auctions are inefficient.** Winners *pay their bid*, so the optimal bid depends on everyone else's bids — users must **estimate** the clearing price, producing chronic overpayment in calm periods and stuck txs/volatility in spikes. First-price auctions are **not DSIC** — but *are* MMIC (the intended allocation already maximizes miner revenue).[^5][^4][^10] **[FACT]**

**Why a true second-price/VCG auction is NOT miner-incentive-compatible (the central impossibility intuition):**[^5][^4]
- A VCG/Vickrey (second-price) auction *is* DSIC for users.
- **But it is not MMIC** — because each winner's payment depends on *other* bids, the proposer can **insert a fake (shill) transaction** to inflate the price real winners pay. Roughgarden's Example 3.5: bids 10, 8, 3 for 3 slots → honest second-price → each pays 3 → revenue 9; insert a fake bid of 8, include the two real top txs + fake → winners pay 8 → revenue 16. The proposer profits from shilling.
- Deeper reason: MMIC essentially **requires a "separable" payment rule** — a tx's payment must not depend on the other included txs' bids (Thm 5.2 / Cor 5.3). First-price, 1559, and tipless all satisfy this; second-price/VCG inherently violates it. This is why on-chain TFMs **cannot** simply adopt VCG: the trustless, miner-controlled setting breaks VCG's auctioneer-honesty assumption.

**EIP-1559's "report card":**[^5][^4] **MMIC ✓** (separable: base fee independent of other txs), **OCA-proof ✓** (because the base fee is **burned**), **DSIC ✓ except during demand spikes** (reverts to a first-price auction on the tip). Roughgarden's **"tipless mechanism"** is an *incomparable* alternative: MMIC + DSIC always, OCA-proof except in a spike. **[FACT]**

**The impossibility frontier.** Roughgarden left open whether *any* TFM satisfies all three for finite blocks; later work answers largely **no** under the standard model — Chung & Shi propose relaxed notions + a randomized "burning second-price auction" and prove **randomness is necessary** for useful collusion-resistant mechanisms; a "dream TFM" is generally impossible without relaxations (~2.3% of Ethereum blocks are congested, so the first-price-fallback regime is not merely theoretical).[^5] **[QUALIFIED.]**

## Disconfirming notes

- **"EIP-1559 lowered gas fees" — FALSE.** Empirical study (Liu/Zhao et al., arXiv:2201.05574, CCS'22): 1559 had "only a small effect on gas fee levels"; it **did not lower the fee level**, but **made fee estimation easier** and **reduced intra-block gas-price variance** and waiting times. "High gas fee is a scalability issue, not a mechanism-design issue." The post-4844 fee decline came from L2 migration, not 1559.[^10] **[FACT]**
- **Base-fee manipulation by non-myopic miners.** MMIC is for *myopic* miners; the DISC'23 paper (Azouvi et al.) shows **non-myopic** proposers can profitably manipulate (mine emptier blocks to drive the base fee toward 0, then harvest tips), weakening 1559's guarantees against forward-looking miners.[^10] **[QUALIFIED.]**
- **Estimation-harder concern resolved.** An earlier worry that base-fee volatility makes estimation harder is contradicted by the intra-block-variance (IQR) data favoring easier estimation.[^10] **[QUALIFIED — resolved contradiction.]**
- **The 1559 burn/deflation claim is conditional** — deflationary only when burn > issuance; quiet/post-L2-migration periods net-inflate. (See `references/tokenomics.md`.)[^10][^11] **[FACT]**

## Child concepts (future research)

1. TFM impossibility & collusion-resilience frontier (Chung–Shi randomized burning second-price auction; necessity of randomness).
2. Multidimensional fee markets (Diamandis–Evans–Chitra–Angeris; per-resource vs per-account; when multidimensionality is worth the complexity).
3. Solana congestion-control / scheduler economics (de-facto-vs-deterministic local markets; Agave scheduler; stake-weighted QoS; SIMD-0096/0110 successors).
4. EIP-1559 base-fee controller dynamics (AIMD/tâtonnement stability, variable learning rates, the 50%-target equilibrium as a control object).
5. Blob fee market evolution & the rollup-data economy (EIP-7691→7840→7892 BPO schedule; the target-demand paradox; PeerDAS/danksharding endgame).

## References

[^1]: EIP-1559: Fee market change — https://eips.ethereum.org/EIPS/eip-1559 — docs (primary spec) — base fee = burned reserve price; max/priority fee semantics; ~12.5% max change; parent-only dependence.
[^2]: EIP-1559 spec constants — https://eips.ethereum.org/EIPS/eip-1559 — docs — ELASTICITY_MULTIPLIER=2, BASE_FEE_MAX_CHANGE_DENOMINATOR=8, INITIAL_BASE_FEE=1 gwei, update formula, 7-wei floor.
[^4]: Roughgarden, "Transaction Fee Mechanism Design…EIP-1559" (arXiv:2012.00854) — https://timroughgarden.org/papers/eip1559.pdf — paper (canonical) — base fee = burnt reserve price; 9/8·7/8 example; burn-prevents-off-chain-first-price-auction.
[^5]: Roughgarden, EC'21/SIGecom summary (arXiv:2106.01340) — https://arxiv.org/html/2106.01340v3 — paper — DSIC/MMIC/OCA-proof; Example 3.5 (second-price not MMIC, 9 vs 16 shill); 1559 report card; tipless; restricted-bid-space framing.
[^6]: EIP-4844: Shard Blob Transactions — https://eips.ethereum.org/EIPS/eip-4844 — docs (primary spec) — GAS_PER_BLOB=2¹⁷, TARGET=393216 (3 blobs), MAX=786432 (6 blobs), MIN_BASE_FEE=1, UPDATE_FRACTION=3338477, fake_exponential; no blob priority fee.
[^7]: "On Blob Markets, Base Fee Adjustments" — Ethereum Research — https://ethresear.ch/t/on-blob-markets-base-fee-adjustments-and-optimizations/21024 — forum (core) — exponential self-adjusting blob fee; ±12.5%; fraction = target/ln(1.125).
[^8]: "Impact of EIP-4844 on…Blob Gas Fee Markets" (arXiv:2405.03183) — https://arxiv.org/pdf/2405.03183 — study — first multidimensional fee market; blob market base-fee-only; blob fee pinned ~1 wei early (target-demand paradox).
[^9]: Analysis of EIP-1559 — Paradigm — https://www.paradigm.xyz/2020/06/analysis-of-eip-1559 — research — tip = first-price auction at the margin; burn prevents off-chain refund; FPA fallback in congestion.
[^10]: Liu/Zhao et al., "Empirical EIP-1559" (arXiv:2201.05574, CCS'22) + "EIP-1559 In Retrospect" — https://export.arxiv.org/pdf/2201.05574v4.pdf · https://decentralizedthoughts.github.io/2022-03-10-eip1559/ — peer-reviewed + summary — small effect on fee levels; lower intra-block IQR; easier estimation; "did not lower fees"; non-myopic manipulation (DISC'23).
[^11]: EIP-1559 burn rationale + EF blog — https://blog.ethereum.org/2020/06/16/eth1x-1559 — docs — burn removes miner manipulation incentive; "only ETH pays for tx"; concedes loss of fixed supply.
[^12]: Off-chain-agreement / burn linkage — Paradigm (as [^9]) + Roughgarden — https://www.paradigm.xyz/2020/06/analysis-of-eip-1559 — research — pay-90-refund-30 off-chain attack absent the burn; burn makes the base fee enforceable.
[^13]: Diamandis, Evans, Chitra, Angeris, "Dynamic Pricing for Non-fungible Resources" (arXiv:2208.07919; AFT'23) — https://ar5iv.labs.arxiv.org/html/2208.07919 — paper — one-dim gas hardcodes relative prices; derives 1559-like controller from a loss function; extends to per-account; Sealevel enables parallel exec + local markets.
[^14]: EIP-7691: Blob throughput increase (Pectra) — https://eips.ethereum.org/EIPS/eip-7691 — docs (spec) — target 6 / max 9 blobs; UPDATE_FRACTION_PRAGUE=5,007,716; 2:3 target:max breaks symmetry. (Volatile.)
[^15]: "The Truth about Solana Local Fee Markets" — Helius — https://www.helius.dev/blog/solana-local-fee-markets — specialist infra — base fee 5,000 lamports/sig; SIMD-0096→100% validator; SIMD-0110 EMA write-lock fee (100% burnt); SIMD-0110 inactive/closed.
[^16]: "Solana Fees, Part 1" — Umbra Research — https://umbraresearch.xyz/writings/solana-fees-part-1 — research — priority fee on *requested* CUs; state declaration → parallel exec; today's locality is non-deterministic, not consensus-enforced.
[^17]: SIMD-0110 PR + Eclipse analysis — https://github.com/solana-foundation/solana-improvement-documents/pull/110 · https://www.eclipse.xyz/articles/solana-ethereum-transaction-fee-mechanisms-proposals-to-improve-solanas-tfm — primary proposal + analysis — exponential write-lock fee; EMA; bounded cache (≤2048); "like EIP-1559 but per-account"; 100% burnt.
[^18]: Solana Docs — Compute Budget — https://solana.com/docs/core/fees/compute-budget — docs (primary) — SetComputeUnitLimit/Price; default 200k CU/instruction, max 1,400,000 CU/tx.
[^19]: Solana Docs — Fee Structure — https://solana.com/docs/core/fees/fee-structure — docs (primary) — base 5,000 lamports/sig; prioritization_fee = ceil(cu_price × cu_limit / 1e6); priority 100% to validator; on requested CUs.
