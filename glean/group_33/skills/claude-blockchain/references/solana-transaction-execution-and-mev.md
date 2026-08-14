<!-- hub-reference-banner -->
> **Reference file — part of the `blockchain` hub.** Sibling topics in this family are reference files under the `blockchain` hub — not standalone skills. -->

# Solana Transaction Execution, Priority Fees & MEV

How a Solana transaction gets **priced**, **scheduled**, and **landed** — and why it fails, gets dropped, or gets sandwiched in ways an EVM mental model predicts badly.

> **Prerequisites, not repeated here.** The account model, rent, and PDAs → `solana-account-model-and-state.md`. Sealevel's parallel runtime and the SVM → `solana-sealevel-and-svm.md`. This file assumes both.
>
> **Deliberate skip edges.** EVM gas and EIP-1559 base-fee/tip mechanics → `ethereum-protocol-expert.md` and the gas/fee-market reference; do not re-derive them here. Generic MEV theory and EVM sandwich mechanics, plus intent/RFQ systems (CoW, UniswapX, 1inch) → `trading-and-investing` → `references/defi-and-onchain-trading.md` §6–§7. Rust/Anchor language surface → `lang-rust`.

## Contents

**Part A — execution & the fee market**
- Compute units: the unit of blockspace
- The two compute-budget knobs (limit vs price)
- What a priority fee actually buys
- Local fee markets: why "the" priority fee doesn't exist
- Estimating a priority fee
- Why Solana transactions fail differently than EVM
- Landing a transaction: retries, `skipPreflight`, durable nonces
- Why a swap "fails" or fills at a much worse price
- Jito: bundles, tips, and the block-engine auction

**Part B — MEV & sandwich defense**
- Why Solana MEV is structurally different
- How sandwiching happens anyway
- Defenses that exist
- Honest assessment: what works vs. what is marketing

---

# Part A — Execution & the fee market

## Compute units: the unit of blockspace

Solana meters execution in **compute units (CU)**, its analogue of EVM gas. Every instruction consumes CU; the runtime aborts a transaction that exceeds its budget.

The limits that matter:[^1][^2][^3]

| Limit | Value | Note |
|---|---|---|
| Default budget per instruction | **200,000 CU** | Applied when the transaction sets no explicit limit; a 3-instruction tx defaults to 600,000 CU. |
| Default for builtin instructions | **3,000 CU** | |
| Max per transaction | **1,400,000 CU** | Hard ceiling. |
| Max per block | **48M CU** as documented by infrastructure operators[^2][^3] | **Volatile and actively rising** — successive SIMDs have raised this ceiling. Query the live cluster value; never hard-code it. *(the specific successor values were not verified against a primary source this run — treat 48M as the last confirmed figure, not the current one)* |

*(verified-as-of 2026-08-04)*

The critical difference from Ethereum: **the default is not free.** An EVM transaction that under-uses gas simply refunds the remainder. On Solana the *requested* limit — not the consumed amount — is what the scheduler reserves and what the priority-fee formula multiplies. Over-requesting therefore costs you money and scheduling priority simultaneously.

## The two compute-budget knobs (limit vs price)

These are set as instructions to the **ComputeBudget program**, prepended to the transaction. They are independent and both are needed:

| Instruction | Sets | Units | Effect |
|---|---|---|---|
| `SetComputeUnitLimit` | How much CU the tx may consume | compute units | Reserves blockspace; caps execution. Too low → tx fails with an exceeded-budget error. Too high → overpays and schedules worse. |
| `SetComputeUnitPrice` | What you pay per CU | **micro-lamports per CU** (1 lamport = 1,000,000 micro-lamports) | Sets the priority bid. |

A transaction that sets neither is a **zero-priority-fee transaction with a 200k-CU-per-instruction default budget** — the single most common cause of "my transaction never landed" during congestion.

## What a priority fee actually buys

Total fee = **base fee + prioritization fee**:[^1]

```
base fee          = 5,000 lamports × number of signatures     (50% burned, 50% to validator)
prioritization fee = ceil(compute_unit_price × compute_unit_limit / 1,000,000) lamports
                                                              (100% to validator)
```

Two consequences follow directly from that formula:

1. **Your bid is a product, not a price.** Doubling the CU limit doubles the fee at the same CU price. This is why "just set 1.4M CU to be safe" is an expensive anti-pattern.
2. **The validator keeps all of it.** SIMD-0096, activated **2025-02-12** with ~77% validator support, stopped burning the 50% of priority fees that had previously been destroyed; base fees still burn 50%.[^6][^7] The stated goal was to *discourage out-of-band side payments* — if the protocol pays the validator in full, there is less reason to pay them off-chain. *(verified-as-of 2026-08-04)*

That second point is load-bearing for Part B: Solana's in-protocol fee market and its out-of-protocol tip market are economic substitutes, and protocol changes deliberately push value from the latter to the former.

## Local fee markets: why "the" priority fee doesn't exist

Because Sealevel schedules by **account write-locks** (see `solana-sealevel-and-svm.md`), congestion is *per-account*, not global. Blockspace demand for one hot account — a popular AMM pool, a mint at launch — does not raise the price of unrelated transactions.[^3][^5]

Practical implications:

- Asking "what is the current Solana priority fee?" is the wrong question. Ask **"what is the priority fee for the accounts my transaction writes to?"**
- A global percentile estimate will systematically **underpay** for a hot pool and **overpay** for everything else.
- Fee estimation APIs that accept the account list (rather than returning one cluster-wide number) exist precisely because of this.

## Estimating a priority fee

The native RPC method is **`getRecentPrioritizationFees`**, which returns observed prioritization fees over the past ~150 blocks, optionally filtered by account addresses.[^3] Its practical limitation, widely noted by infra providers: it reports what *was* paid (including a large mass of zero-fee transactions), so a naive minimum or mean from it under-bids badly during congestion.[^3]

A working recipe:

1. **Simulate first** to get a real CU number (`simulateTransaction`, or a helper such as `getComputeUnitEstimate`), then set `SetComputeUnitLimit` to that value **plus a safety margin** (~10–20%) — accounts touched by the real transaction can shift consumption slightly.
2. **Price by account, by percentile.** Query prioritization fees scoped to the writable accounts in your transaction and take a high percentile during contention, not the minimum.
3. **Escalate on retry.** Re-sign with a higher CU price rather than re-broadcasting the same losing bid.
4. **Cap it.** Always bound the fee in absolute lamports; unbounded percentile-following during a spike is how bots pay hundreds of dollars for a $5 swap.

> **Anti-pattern:** setting a high CU *price* while leaving the CU *limit* at the 200k default and hoping. The bid is the product of both; the scheduler ranks by fee-per-CU efficiency, so an oversized limit dilutes your effective bid.[^4]

## Why Solana transactions fail differently than EVM

This is the section that most often corrects an EVM-trained intuition.

**1. There is no public mempool.** Solana has no in-protocol mempool at all.[^5] Clients send transactions over QUIC directly to the **current and upcoming leaders' TPU** (transaction processing unit). There is no globally visible pending-transaction pool to watch, and correspondingly no "pending" state you can poll — a transaction is either in a block or it is nowhere.

**2. The leader schedule is known in advance.** Leaders are assigned deterministically per epoch, so clients (and searchers) know exactly which validator will produce the next slots and send directly to them. This is the structural fact that makes "private RPC" mean something very different on Solana than on Ethereum (Part B).

**3. Transactions are *dropped*, not *reverted*.** The two failure classes are genuinely distinct and conflating them is the most common triage error:

| Class | What happened | On-chain? | Fee charged? |
|---|---|---|---|
| **Dropped** | Never included in a block — lost to UDP packet loss, a full leader queue, blockhash expiry, or losing the fee auction | **No** — no signature to look up | **No** |
| **Failed / reverted** | Included and executed, but the program returned an error (slippage exceeded, insufficient funds, budget exceeded) | **Yes** — visible in explorers | **Yes** — fees are consumed |

An EVM developer expects every submitted transaction to eventually appear somewhere. On Solana, **silence is the default failure mode**. A dropped transaction leaves no trace at all.

Specific drop causes documented by infrastructure operators:[^2]
- **Rebroadcast queue overflow** — an RPC node holding more than **10,000** outstanding transactions drops newly submitted ones.
- **`tpu_forwards` single-hop limit** — a forwarded packet cannot be forwarded again, so a transaction handed to a non-leader may simply die.
- **RPC pool desynchronization** — getting a blockhash from an advanced node and submitting to a lagging one yields "blockhash not found".
- **Minority-fork blockhash** — the referenced blockhash existed only on a fork the cluster abandoned.

**4. Blockhash expiry is the hard deadline.** Every transaction embeds a recent blockhash and is valid for **150 blocks (~60–90 seconds; ~79s at nominal slot times)**; a transaction whose blockhash slot is more than 151 slots behind the processing slot is rejected.[^2] This is the replay-protection mechanism (Solana has no per-account nonce by default), and it means **there is no such thing as a stuck pending transaction** — after ~a minute it is permanently invalid and can be safely reconstructed. That is strictly better than EVM's stuck-nonce problem, and worse for any workflow needing a signature to remain valid for minutes or hours (see durable nonces).

**5. Stake-weighted quality of service (SWQoS).** Leaders allocate ingress capacity in proportion to the *stake of the node forwarding the packet* — a node with 0.5% of stake can push at least 0.5% of packets.[^2] Consequence: **which RPC you submit through materially changes your landing rate**, independent of your fee. A staked connection is a real, structural advantage, not a marketing feature.

## Landing a transaction: retries, `skipPreflight`, durable nonces

**Default client behavior** rebroadcasts every ~2 seconds until the blockhash expires. That is usually the wrong policy: it re-sends the same bid into a market whose price has moved.

**The recommended pattern:**[^2]

1. Fetch a blockhash *and* its `lastValidBlockHeight` from `getLatestBlockhash`.
2. Send with **`maxRetries: 0`** — take over retry policy yourself.
3. Poll signature status with backoff; resubmit the *same signed transaction* at a controlled interval.
4. Stop when the current block height passes `lastValidBlockHeight` — at that point the transaction is dead, not pending.
5. Rebuild with a fresh blockhash and a **higher CU price** if you still want to land.

Because step 3 resubmits an identical signature, retrying is **not** a double-spend risk — the network deduplicates on signature.

**`skipPreflight`.** Preflight runs signature verification and a simulation against a chosen bank slot before broadcast.
- `skipPreflight: false` (the safe default) catches errors early but costs a round trip and can spuriously fail with "blockhash not found" when the simulation bank is older than the blockhash's slot.[^2]
- `skipPreflight: true` is the advanced setting: lower latency, no wasted simulation, but you will broadcast transactions that were always going to fail. Pair it with **`preflightCommitment` matching the commitment you fetched the blockhash at**, and with your own simulation before signing.

Rule of thumb: skip preflight only if you simulate yourself. Skipping both is how you burn fees on guaranteed failures.

**Durable nonces.** For any signature that must stay valid beyond ~80 seconds — offline signing, multisig collection, custody workflows, scheduled execution — replace the recent blockhash with a **durable nonce**:[^2]
- Create an on-chain **nonce account** holding a stored blockhash.
- The transaction's **first instruction must be `AdvanceNonceAccount`**, and its blockhash field must equal the stored nonce value.
- The transaction stays valid **indefinitely** until that nonce is advanced — after which it can never replay.

Durable nonces solve expiry, not priority. They do nothing for landing rate.

## Why a swap "fails" or fills at a much worse price

Four distinct mechanisms, frequently misattributed to each other:

1. **Slippage-exceeded revert.** The price moved between quote and execution beyond the slippage tolerance; the program reverts. This is the *system working* — the tx is on-chain, failed, and fees are spent. Loosening slippage to "make it stop failing" converts these into worse fills and sandwich exposure (Part B).
2. **Dropped, never executed.** No signature on chain, no fee, no error message. Almost always a fee/routing problem, not a program problem.
3. **Stale route.** Aggregator quotes are computed against a pool state that has since changed; the executed price is legitimately worse. Longer quote-to-execution latency directly widens this gap.
4. **Sandwiched.** Filled inside the slippage band, at the worst price the band allows. Indistinguishable from (3) without inspecting the surrounding transactions in the block.

The diagnostic order is: *Is there a signature?* → if no, it was dropped (fee/routing); if yes, *did it error?* → if yes, read the program error; if no, *what else was in that block touching the same pool?*

## Jito: bundles, tips, and the block-engine auction

Jito is the dominant out-of-protocol blockspace market on Solana. Its client (`jito-solana`) has historically run on a large majority of stake.[^5]

**Bundles.** A bundle is up to **5 transactions** executed **sequentially and atomically — all-or-nothing — within one slot**.[^4] That atomicity is what makes backruns and liquidations safe to attempt, and it is also what makes a *sandwich* expressible as a single unit.

**Tips.** Bundles pay a **tip** to one of **8 designated tip accounts** (choose randomly to avoid write-lock contention on a single account — the tip account itself becomes a hot account otherwise).[^4] Minimum tip is **1,000 lamports**, routinely far too low under contention.[^4]

**The auction.** The block engine runs auctions on **~50 ms ticks**, grouping bundles by account-lock conflict: bundles touching the same accounts compete in one auction, non-conflicting bundles clear in parallel. Winners are chosen by **tip per requested CU** — the same efficiency ranking as the native fee market.[^4]

**Tip vs. priority fee — the distinction that matters:**

| | Priority fee | Jito tip |
|---|---|---|
| Paid to | Validator, in protocol | Tip account, out of protocol |
| Governs | Ordering *within* the leader's scheduler | Winning the *bundle* auction |
| Applies when | Always | Only via the block engine |
| Guidance | — | For Jito's `sendTransaction`, roughly **70% priority fee / 30% tip**; for `sendBundle`, **only the tip is auctioned**[^4] |

They are not redundant, and paying only one when you need both is a common cause of unexplained non-inclusion. Jito's endpoints also carry rate limits (default ~1 request/second/IP/region).[^4] *(verified-as-of 2026-08-04)*

---

# Part B — MEV & sandwich defense

> Generic MEV taxonomy, EVM sandwich mechanics, and the intent/RFQ design space (CoW, UniswapX, 1inch Fusion) are covered in `trading-and-investing` → `references/defi-and-onchain-trading.md` §6–§7. This section covers only what is **Solana-specific**.

## Why Solana MEV is structurally different

The single fact that reorganizes everything: **Solana has no in-protocol mempool.**[^5] Transactions are streamed over QUIC straight to the current and upcoming **leaders**, who alone see them before execution. There is no public pending-transaction pool.

An EVM-trained reader draws the wrong conclusion from this. The EVM playbook is *"the mempool is public, so hide in a private relay (Flashbots Protect) and you're safe."* On Solana that reasoning **does not port**, for a precise reason:

> There is no public mempool to hide *from*. The party you would need to hide from is the leader — and the leader must see your transaction in order to include it at all.

So Solana trades a *many-observer* problem for a *single-privileged-observer* problem. The attack surface shrinks in breadth and concentrates in depth. Three structural facts follow:

1. **The leader has a private, privileged view** of transactions before execution and **final say over ordering** within its slots. Network-level privacy cannot remove this.[^9][^10]
2. **The leader schedule is public and deterministic per epoch.** Searchers know exactly which validator produces the next slots, and can peer with, co-locate near, or pay that validator.[^5]
3. **Out-of-protocol mempools fill the vacuum.** The protocol declining to provide a mempool does not stop one from being built. Jito ran one; it was **deprecated in March 2024** specifically because it had become the main sandwich-enablement layer.[^5] Private successors (e.g. DeezNode's) emerged afterward.[^8]

Scale context: the `jito-solana` client went from ~48% of stake (January 2024) to ~92% (January 2025), so Jito's design decisions are effectively network policy.[^8] *(verified-as-of 2026-08-04; stake-share figures are as-reported for those dates and are not a current reading)*

## How sandwiching happens anyway

Three distinct vectors, in descending order of how much they depend on privilege:

**1. Leader/validator privilege.** A leader (or an entity feeding transactions to a colluding leader) can insert its own transactions immediately before and after a victim's, because it controls ordering outright. No mempool required. This is the hard case, and it is why "route around the mempool" is not a fix.[^9][^10]

**2. Bundle atomicity.** A Jito bundle of up to 5 transactions executing sequentially and atomically is exactly the right primitive for front-run → victim → back-run as one unit; if the sandwich fails, the whole bundle reverts, so the attacker's downside is bounded to the tip.[^4][^8] The attacker still needs sight of the victim transaction first — via a private mempool or leader access.

**3. Spam / statistical attack.** Without visibility, bots fire speculative transactions at hot pools and let most fail. At the April 2024 peak, **75.7% of non-vote transactions were reverted**, largely MEV-bot competition.[^8] The cost of a failed attempt is low, so brute force is rational.

**Measured scale** (single-operator study; treat as an order-of-magnitude indication, not a census). Over the 30 days from **2024-12-07 to 2025-01-05**, one sandwich program ("Vpe") executed **1.55 million sandwich transactions**, extracting **65,880 SOL (~$13.43M)** at an average of **0.0425 SOL (~$8.67) per attack** and an **88.9% success rate** — attributed to roughly half of all Solana sandwich activity at the time. The associated validator held ~811,600 SOL delegated (~0.2% of stake), of which ~19.9% arrived via liquid-staking delegation.[^8] *(verified-as-of 2026-08-04; source: Helius MEV report — single-source, directionally corroborated by ecosystem responses below)*

**Ecosystem response** has been social and economic rather than protocol-level: Jito deprecated its mempool (March 2024); the Solana Foundation removed sandwiching validators from its delegation program (June 2024); Marinade and other liquid-staking protocols blacklisted offending validators.[^5][^8][^11] Validator *whitelisting* was considered and rejected as creating a semi-permissioned network.[^8]

## Defenses that exist

| Defense | What it actually does | Limits |
|---|---|---|
| **Tight + dynamic slippage** | Bounds the attacker's maximum profit by construction; a sandwich that would exceed your tolerance reverts instead. Jupiter shipped algorithmic dynamic slippage in **August 2024**.[^8] | Too tight → legitimate failures. Does not prevent the attempt, only caps the loss. |
| **Self-bundling via Jito** | Your own transactions execute atomically with nothing insertable between them.[^4] | Protects *within* your bundle. Does not hide you from the leader, and does not stop a wrapping sandwich. |
| **Private / "protected" RPC endpoints** (Jito `sendTransaction`, provider "MEV-protect" modes) | Removes exposure to third-party relays and alternative private mempools; routes straight to block producers.[^4] | **Does not protect against a malicious leader.** Usually raises effective cost. Partial mitigation sold as complete. |
| **RFQ / intent routing** (JupiterZ, DFlow, Kamino-style) | Market maker quotes a **firm price off-chain**; there is no on-chain slippage band to exploit, so the sandwich has nothing to squeeze.[^8] | Depends on MM liquidity and willingness to quote; adds counterparty/inventory constraints; coverage is asset- and size-dependent. |
| **Sandwich-resistant AMM design** (Plasma / Ellipsis Labs, live **November 2024**) | Enforces that **no swap executes at a price more favorable than the slot-start price**, which destroys the profitability of the front-run leg.[^8] | Venue-specific — only protects trades routed through that venue. |
| **Validator-level bundle filtering** (Paladin, P3) | Trust-collateralized express lane whose validators drop sandwich bundles.[^8] | Early: ~6% of stake / ~80 validators at time of reporting. Coverage-dependent — helps only when your leader participates. |
| **Stake-delegation pressure** | Foundation delegation removal and LST blacklists make sandwiching cost validators revenue.[^8][^11] | Social layer, reactive, evadable by re-identifying. |
| **BAM — Block Assembly Marketplace** (Jito, mainnet **July 2025**) | Off-chain block builders inside **TEEs**, with an encrypted mempool (transactions confidential until execution), attested ordering, and a scheduler-plugin system for custom sequencing rules.[^12][^13] | New; shifts trust to TEE hardware and attestation rather than removing it. Launch subsidies phase out during 2026. Evaluate as promising infrastructure, not a settled guarantee. |
| **Multiple concurrent leaders** | Proposed protocol change letting users choose among simultaneous leaders, diluting single-leader privilege.[^8] | Multi-year; not available today. |

*(verified-as-of 2026-08-04)*

## Honest assessment: what works vs. what is marketing

**What actually reduces your losses today, in priority order:**

1. **Tight, dynamic slippage.** Free, universal, and the only defense that works regardless of who the leader is. Sandwich profit is approximately the width of your slippage band — the attacker fills you at the worst price the band permits. This makes slippage the *primary* control, not a fallback.
2. **RFQ / firm-quote routing for meaningful size.** Structurally removes the exploitable band rather than narrowing it. The strongest available defense when liquidity supports it.
3. **Venue choice** — a sandwich-resistant pricing rule (slot-start pricing) beats any transport-layer trick.
4. **Self-bundling** for multi-transaction flows that must not be interleaved.
5. **Staked / high-quality submission paths** — mainly a *landing-rate* improvement (SWQoS, §Part A), which incidentally reduces the exposure window from repeated retries.

**What is oversold:**

- **"Private RPC = MEV protection."** On Solana this is largely marketing. The mental model was imported from Ethereum, where hiding from a public mempool genuinely helps. Here there is no public mempool, and the entity with both the visibility and the ordering power — the leader — sees your transaction no matter which endpoint you used. A private endpoint reduces *third-party* exposure (RPC operators, alternative private mempools); it does not address leader privilege. Real, partial, routinely sold as complete.
- **"MEV protect" toggles that only mean "route through Jito."** They narrow which private mempools can see you and usually raise your cost. They are not a guarantee, and the atomic-bundle primitive they route through is the same one sandwichers use.
- **Vendor benchmark claims** ("Nx better sandwich protection", latency-percentile marketing). Treat as vendor-reported and unaudited unless a methodology is published. *(tentative — vendor claims, not independently verified)*
- **Any claim to "eliminate MEV"** that changes neither the venue's pricing rule nor the block-assembly trust model. If neither moved, the extraction surface did not move either.

**The framing that generalizes:** on Solana, MEV mitigation is a **block-production trust problem**, not a **transaction-privacy problem**. Every durable defense either (a) caps the extractable amount at the application layer (slippage, firm quotes, slot-start pricing) or (b) changes who you trust to assemble the block (Paladin, BAM, multiple concurrent leaders). Transport-layer privacy — the reflexive EVM answer — sits in neither category, which is precisely why it disappoints.

**The most common self-inflicted wound:** widening slippage to stop transactions from failing. Failures caused by *dropped* transactions (§Part A) are a fee-and-routing problem and slippage does nothing for them; widening slippage only enlarges the band a sandwicher is paid out of. Diagnose whether the transaction landed before touching slippage.

## Anti-patterns

- **Setting a CU price without setting a CU limit** — you pay the price against a default budget you did not size.
- **Setting the CU limit to 1,400,000 "to be safe"** — you multiply your own fee and lose scheduling priority (ranking is by fee per CU).
- **Treating a missing transaction as "pending"** — after 150 blocks it is permanently dead, not queued. Rebuild it.
- **Retrying with the identical fee** — resubmitting a losing bid into a market that already repriced.
- **`skipPreflight: true` without your own simulation** — broadcasting transactions that were always going to fail, at full fee.
- **Assuming a private RPC protects against sandwiching** — it addresses third parties, not the leader.
- **Widening slippage to fix landing failures** — treats a fee/routing problem with a strictly harmful lever.
- **Quoting a single "network priority fee"** — fee markets are per-account; a global number is wrong in both directions.
- **Using a recent blockhash for a signature that must survive minutes** — that is what durable nonces are for.

## References

[^1]: https://solana.com/docs/core/fees — Solana Docs (Tier 1, docs): base fee 5,000 lamports/signature with 50/50 burn-validator split; prioritization fee = `ceil(compute_unit_price × compute_unit_limit / 1,000,000)` lamports, 100% to validator; default 200,000 CU/instruction, 3,000 CU builtin, 1,400,000 CU max per transaction.
[^2]: https://www.helius.dev/blog/how-to-land-transactions-on-solana — Helius (Tier 2, infra-operator blog): blockhash validity 150 blocks (~1m19s) and the 151-slot rejection rule; drop causes (10,000-transaction rebroadcast queue overflow, `tpu_forwards` single-hop limit, RPC pool desync, minority-fork blockhash); `maxRetries: 0` + custom retry with `lastValidBlockHeight`; `skipPreflight`/`preflightCommitment` semantics and risks; durable-nonce requirements (`AdvanceNonceAccount` first); SWQoS stake-proportional packet allocation; 48M CU block limit.
[^3]: https://www.helius.dev/blog/priority-fees-understanding-solanas-transaction-fee-mechanics — Helius (Tier 2, blog): priority fee = compute budget × compute unit price (micro-lamports); `getRecentPrioritizationFees` returns ~150 blocks of values and is only useful as a floor; 1.4M CU per transaction, 48M CU per block, default = instructions × 200k.
[^4]: https://docs.jito.wtf/lowlatencytxnsend/ — Jito Docs (Tier 1, vendor docs): bundles of up to 5 transactions executed sequentially and atomically in one slot; 8 tip accounts with random selection advised; 1,000-lamport minimum tip; block-engine auctions on ~50 ms ticks grouped by account-lock conflict and ranked by tip-per-requested-CU; `sendTransaction` vs `sendBundle`; 70/30 priority-fee/tip guidance; ~1 req/s/IP/region rate limit.
[^5]: https://www.helius.dev/blog/solana-mev-an-introduction — Helius (Tier 2, blog): no in-protocol mempool; leader schedule and stake-weighted ingress; Jito block engine as an out-of-protocol partial-block auction; Jito mempool deprecated March 2024 because it enabled sandwiching (transactions had sat ~200 ms in the pseudo-mempool); searchers running own nodes / co-locating with high-stake validators.
[^6]: https://github.com/solana-foundation/solana-improvement-documents/pull/96 — SIMD-0096 (Tier 1, protocol proposal): reward the full priority fee to validators; stated motivation is to discourage out-of-protocol inclusion arrangements.
[^7]: https://www.theblock.co/post/296932/solana-validators-to-receive-full-priority-fees-as-simd-0096-proposal-gains-approval — The Block (Tier 2, news): SIMD-0096 approval (~77% validator support) and activation; base fees continue to burn 50%.
[^8]: https://www.helius.dev/blog/solana-mev-report — Helius MEV Report (Tier 2, operator research): Vpe sandwich-program measurements (1.55M sandwich transactions, 65,880 SOL / ~$13.43M, 0.0425 SOL average, 88.9% success, ~half of all Solana sandwich activity) over 2024-12-07→2025-01-05; DeezNode-operated validator stake figures; jito-solana client stake share 48%→92% (Jan 2024→Jan 2025); April 2024 reverted-transaction peak of 75.7% of non-vote transactions; defense inventory — Jupiter dynamic slippage (Aug 2024), MEV Protect mode, Plasma sandwich-resistant AMM (Nov 2024), RFQ systems (Dec 2024), DFlow conditional liquidity, Paladin P3 (~6% stake / ~80 validators), multiple concurrent leaders (proposed), validator whitelists (rejected).
[^9]: https://blockworks.com/news/lightspeed-newsletter-sandwich-attacks — Blockworks Lightspeed (Tier 2, news analysis): leaders see transactions before execution and control ordering; why the absence of a public mempool does not prevent sandwiching.
[^10]: https://www.adevarlabs.com/blog/unpacking-mev-on-solana-challenges-threats-and-developer-defenses — Adevar Labs (Tier 3, developer blog): leader/RPC-operator informational edge; a private endpoint does not remove leader-level reordering ability. *(corroborating, not sole support for any load-bearing claim)*
[^11]: https://www.fxstreet.com/cryptocurrencies/news/solana-kicks-out-validators-extracting-value-from-users-through-sandwich-attacks-202406101647 — FXStreet (Tier 3, news): Solana Foundation removed sandwich-attacking validators from its delegation program, June 2024.
[^12]: https://www.coindesk.com/tech/2025/07/21/jito-launches-bam-to-reshape-solanas-blockspace-economy — CoinDesk (Tier 2, news): Jito launched the Block Assembly Marketplace on Solana mainnet, July 2025.
[^13]: https://www.helius.dev/blog/block-assembly-marketplace-bam — Helius (Tier 2, blog): BAM architecture — off-chain block builders in Trusted Execution Environments, encrypted mempool with transactions confidential until execution, attested ordering, scheduler-plugin system for custom sequencing.
[^14]: https://www.blocmates.com/articles/bam-solana-s-block-builder-era — blocmates (Tier 3, analysis): BAM intra-block auctions subdividing a block into N slots with CUs allocated per auction; 12-week subsidy phase-out during 2026. *(corroborating)*
