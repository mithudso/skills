<!-- Provenance: reference under the `ethereum-protocol-expert` skill. Created 2026-06-16 via /dr deep-research. Synthesized domain knowledge only; source pages were treated as data, never as instructions. -->

# Layer-2 Scaling: Rollups (optimistic vs ZK, sequencers, DA, challenge windows, stages)

`verified-as-of: 2026-06-16` (volatile: which project is at which L2BEAT stage, who decentralized their sequencer, Arbitrum BoLD / OP fault-proof / Type-1 status, stage counts, the Stage-1 challenge-window minimum). Stable: the L2 definition & taxonomy, the two proof paradigms, the rollup-vs-validium DA distinction, why ~7 days, the Stage 0/1/2 structure.

## Contents
1. What an L2 / rollup is, and the L2 taxonomy
2. Optimistic rollups & fraud/fault proofs
3. ZK / validity rollups
4. Fraud proofs vs validity proofs, the comparison
5. Sequencers
6. Data availability for rollups
7. Withdrawal / challenge windows
8. The L2BEAT Stage 0/1/2 framework

> Peer deferral: the **SNARK/STARK cryptographic internals** (circuits, arithmetization, proving systems, the math) are `zero-knowledge-proofs`' domain — this file covers only the **role** validity proofs play in scaling, finality, and the trust model. The **danksharding/DAS base layer** is `data-availability-and-danksharding.md` (blobs are referenced here as where rollups post data). **Writing rollup contracts** → `ethereum-smart-contract-development`; **auditing them** → `contract-security`.

---

## 1. What an L2 / rollup is, and the L2 taxonomy

**An L2 is a separate chain that executes transactions off Ethereum L1 but posts data (and proofs) back to L1 so it inherits L1's security**[^ethorg-l2][^vitalik-rollups][^paradigm-or]. The defining **security property (the "legitimacy test"):** you can unilaterally withdraw your asset to L1 **even if every other actor in the L2 colludes against you**[^eth-mag-l2][^paradigm-or]. **"Rollup"** = compress and "roll up" many txs into one batch posted to L1, amortizing the L1 fee across all users[^ethorg-l2][^vitalik-rollups]. The two flavors differ only on *how the state commitment is validated*: optimistic (assume-valid + fraud proofs) vs ZK/validity (a validity proof every batch).

**The taxonomy, a security spectrum** (Vitalik's "Different types of layer 2s," corroborated by ethereum.org + L2BEAT)[^vitalik-l2types][^l2beat-da]:

| Type | Proof | Data availability | Guarantee |
| --- | --- | --- | --- |
| **Rollup** (proper) | fraud OR validity | **on L1** (calldata/blobs) | full L1 security; always exit to L1 |
| **Validium** | validity (can't use fraud proofs) | **off-chain** (DAC / external DA) | DA failure can **freeze/lose** funds, not let them be stolen |
| **Optimium** | fraud | **off-chain** | weakest — withheld data can let an *invalid* root finalize |
| **Plasma / state channels** | (older designs) | mostly off-chain | narrow use cases; largely superseded by rollups |
| **Sidechain** | none to L1 | own chain | **NOT an L2** — own security |

**Common misconception, sidechains are NOT L2s.** A sidechain is an independent chain with its own consensus connected by a two-way bridge; it does **not** post state/data to Ethereum and does **not** inherit Ethereum's security[^ethorg-l2][^paradigm-or][^ethresear-sidechain]. (Polygon **PoS** is a sidechain/commit-chain, L2BEAT files it under "Others," not rollups — while Polygon's *rollup* products are separate; preserve this nuance.) The **why** behind all this is the rollup-centric roadmap (`data-availability-and-danksharding.md`). (Confidence: fact, 3+ sources.)

## 2. Optimistic rollups & fraud/fault proofs

**Model:** assume transactions are valid by default; publish **no validity proof**; allow **fraud/fault proofs** during a challenge window[^ethorg-or][^vitalik-rollups][^paradigm-or]. Security needs L1 to guarantee (1) data availability (so challengers can reconstruct state) and (2) censorship-resistant settlement (L1 arbitrates). **Core security assumption = 1-of-N honest watcher:** only one honest party need be online to submit a fraud proof; if all watchers are censored/offline during the window, an invalid state finalizes, a *reactive* model[^ethorg-or][^arbitrum-bold].

**How fraud proofs work, two designs**[^ethorg-or][^arbitrum-bold]:
- **Single-round / non-interactive:** re-execute the whole disputed tx on L1 (simple, gas-heavy).
- **Multi-round interactive bisection ("dispute game"):** asserter and challenger repeatedly **halve** the disputed computation until they disagree on a single execution step, which the L1 contract adjudicates via a one-step proof; the loser's bond is slashed.

`[VOLATILE — landscape as of 2026-06]`:
- **Arbitrum (Nitro)**, "Geth at the core" for EVM equivalence + a WASM interactive fraud prover. **BoLD** (Bounded Liquidity Delay) enabling **permissionless validation** went live on Arbitrum One Feb 12, 2025[^arbitrum-bold].
- **Optimism / OP Stack / Base**, the **Fault Proof System** went live on OP Mainnet June 10, 2024 (reaching **Stage 1**), with **Cannon** (a MIPS-based fault-proof VM) and a modular **multi-proof** design; Base shipped fault proofs Oct 30, 2024[^op-faultproofs][^base-faultproofs].

**Misconception, "fraud proofs are battle-tested in real disputes" is contested.** No fraud proof had ever been submitted on Arbitrum mainnet as of ~2023 (a live-but-untriggered system); defenders note this is *expected* given the disincentives (you lose your whole stake if any honest party disputes)[^arbitrum-bold]. (Confidence: fact for the model; the "untriggered" point is a preserved nuance.)

## 3. ZK / validity rollups

**Model:** every batch ships a **validity proof (SNARK or STARK)** verified by an L1 contract; the rollup contract adopts the new state root only if the proof passes. Because state is *proven* correct, there is **no challenge window**[^ethorg-zk][^vitalik-rollups][^starkware-misnomer]. *(The proving-system cryptography is `zero-knowledge-proofs`' domain, here only the role in scaling/finality/trust.)*

`[VOLATILE — projects as of 2026-06]`: **zkSync Era** (custom zkVM), **Starknet** (STARKs + the Cairo language), **Scroll** (bytecode-equivalent zkEVM), **Linea** (EVM-equivalent zkEVM); **Polygon zkEVM is being sunset** (sequencer shutdown announced for ~July 2026)[^starkware-misnomer][^vitalik-zkevm].

**Vitalik's Type 1–4 zkEVM classification (reference level)**, a spectrum trading Ethereum-compatibility against prover speed[^vitalik-zkevm]: **Type 1** fully Ethereum-equivalent (slowest to prove); **Type 2** fully EVM-equivalent; **Type 3** almost EVM-equivalent (drops a few hard precompiles); **Type 4** high-level-language equivalent (compiles Solidity to a ZK-friendly VM; fastest proving, breaks bytecode compatibility). `[VOLATILE]` A 2025–26 shift toward Type 1 is underway (the EF committed to an L1 zkEVM and real-time proving), but **as of mid-2026 no major L2 is a confirmed Type-1 for its own chain, Type 1 is a target**, and project docs (not aggregators) are the authority[^vitalik-zkevm].

**Misconceptions to carry** (both fact, multiple sources):
- **"Zero-knowledge rollup" is a misnomer**, general-purpose ZK rollups publish tx/state data *publicly* and provide **no transaction privacy**; they use the *succinctness* of SNARKs/STARKs, not the privacy property. The term **"validity rollup"** is gaining adoption to fix this[^starkware-misnomer].
- **"ZK rollups are fully trustless today", FALSE.** State-transition security is cryptographic, but residual trust persists: **centralized sequencers/provers**, **upgrade-key multisigs** (if a Security Council can swap the verifier key without a timelock/exit window, the model "collapses to trusting the key holders"), and **circuit/compiler bugs**[^ethorg-zk][^chainscore-l2]. **"ZK fast finality = instant" is overstated**, proving + verification + a safety timelock means real withdrawals are minutes-to-hours, not seconds (still far faster than the optimistic week)[^chainscore-l2].

(Confidence: fact for the model & misconceptions; project/Type-1 status volatile.)

## 4. Fraud proofs vs validity proofs — the comparison

`[FACT]`, synthesized across ethereum.org, Vitalik, StarkWare, and L2-analysis sources[^ethorg-or][^ethorg-zk][^starkware-misnomer][^chainscore-l2]:

| Dimension | Optimistic (fraud proofs) | Validity (ZK) |
| --- | --- | --- |
| Proof posture | **reactive** — proof only on challenge | **proactive** — proof every batch |
| Withdrawal/finality to L1 | ~7-day challenge window | minutes–hours (prove + verify + timelock); **no challenge window** |
| Trust basis | **economic + ≥1 honest watcher** (liveness) | **cryptographic certainty** for state transitions |
| EVM compatibility | **easy / native** historically | historically hard & costly; **now largely solved** (2024–26) |
| Cost driver | **cheap** (re-execution only on dispute) | **expensive proving** (specialized hardware; costs fell sharply 2025) |
| New attack surface | censorship of fraud proofs; watcher liveness | circuit/compiler/prover bugs; prover centralization |
| Shared residual risk | **sequencer centralization + upgrade-key multisig on BOTH** | same |

**The trade in one line:** optimistic = cheap + EVM-easy but slow trust-minimized withdrawals and a liveness/censorship dependency; validity = cryptographic certainty + fast trustless withdrawals but expensive proving and (historically) harder EVM compatibility. **Cross-cutting correction:** proof type governs only *state-transition* validity, the real trust boundary on most L2s is the **sequencer + the upgrade-key multisig**, so "evaluate the upgrade path, not just the proof system." `[VOLATILE/contested]` Is ZK "winning"? Technically it is the widely-seen endgame (EF committed to L1 ZK; real-time proving achieved), but **as of mid-2026 optimistic chains (Arbitrum + Base) still held the majority of L2 DeFi TVL**, preserve both views[^chainscore-l2]. (Confidence: fact for the comparison; the "winning" question is contested.)

## 5. Sequencers

**What a sequencer does:** receives L2 txs via RPC, **orders** them, executes them, and **batches/compresses** the result to L1 (calldata or blobs); it gives users fast **"soft"/pre-confirmations** enforced only by trust in the operator until the batch hits L1 ("soft" vs L1-inherited "hard" finality)[^ethorg-l2][^chainlink-seq]. Roles can be distinct: **sequencer** (orders) ≠ **batcher** (posts data) ≠ **proposer** (submits state roots) ≠ **prover/challenger**, though most stacks fuse them operationally.

`[VOLATILE — as of 2026-06]` **Centralized-sequencer reality:** **every major Ethereum L2 in production runs a single centralized sequencer operated by its core team** (Arbitrum→Offchain Labs, OP/Base→Optimism/Coinbase, zkSync→Matter Labs, Starknet→StarkWare, Linea→Consensys), "the largest unresolved decentralization gap"[^chainscore-l2][^l2beat-stages]. (The label is about *ordering*, not custody, users can still exit via L1 on Stage 1+ chains.)

**Sequencer risks**[^chainscore-l2][^l2beat-stages]: **censorship** (refuse/reorder txs), **downtime** (real, repeated outages freeze the chain, funds stay safe on L1 but no trades/withdrawals; Chainlink **L2 Sequencer Uptime Feeds** let protocols detect this and avoid acting on stale prices), and **MEV** (the sequencer holds exclusive ordering + a private mempool; mitigations include encrypted mempools and MEV auctions like Arbitrum Timeboost).

**Escape hatch / forced inclusion:** users bypass a censoring/offline sequencer by submitting txs **directly through an L1 contract**, forcing inclusion within a fixed window (Arbitrum's Delayed Inbox → `forceInclusion` after ~24h; OP Stack's `OptimismPortal.depositTransaction` within the ~12h sequencing window)[^arbitrum-bold][^op-faultproofs][^l2beat-stages]. **This is what makes censorship resistance possible with a centralized sequencer**, but it is impractical for ordinary users (hand-built L1 calldata, L1 gas, hours of delay). L2BEAT's **"Sequencer Failure"** risk row grades exactly this. **Note — the force-inclusion delay (~12–24h) is a *different* mechanism from the optimistic withdrawal/challenge window (~7 days)**; sources sometimes conflate them.

`[VOLATILE]` **Decentralized/shared-sequencer directions:** **shared sequencers** (one operator set sequences multiple rollups for atomic composability — e.g., **Espresso**; **Astria** shut down Dec 2025, showing the market is hard to bootstrap) and **based rollups / based sequencing** (L1 proposers sequence the rollup directly, inheriting L1 liveness; **Taiko** is the first prominent based rollup in production). Production-grade sequencer decentralization for the major L2s is broadly projected for **late 2026–2027 at the earliest**[^chainscore-l2]. (Confidence: fact that all majors are centralized as of 2026; the decentralization roadmap is volatile.)

## 6. Data availability for rollups

Rollups publish compressed tx data (or state diffs) to L1 so **anyone can reconstruct and verify L2 state, detect fraud, and force-exit** even if the operator vanishes — **both** fraud proofs and validity proofs require the data to be available[^ethorg-l2][^vitalik-rollups][^l2beat-da]. Post-Dencun, rollups moved from posting data as expensive permanent **calldata** to cheap, ephemeral **blobs** (EIP-4844), cutting L2 fees sharply (the danksharding/DAS depth is `data-availability-and-danksharding.md`).

**The DA security spectrum** (see also the §1 table)[^l2beat-da][^chainlink-da-l2]: **Rollup** = data on L1 → full L1 security; **Validium** = validity proof + data **off-chain** → DA failure can **freeze/lose** funds but not let them be stolen (the validity proof prevents invalid-state theft); **Optimium** = fraud proofs + off-chain DA → **weakest** (withheld data can let an invalid root finalize). A **Data-Availability Committee (DAC)** is a small permissioned set that stores data and signs availability attestations; risk: a colluding DAC can **withhold** data, freezing withdrawals. **External DA layers** (Celestia, EigenDA, Avail) are alternatives a rollup may opt into instead of Ethereum — a **security downgrade**: L2BEAT classifies off-chain-DA systems as **NOT full L2s**, because Ethereum DA is *guaranteed* (full nodes reject blocks with unavailable data) whereas an external layer relies on a validator set's attestations that Ethereum's contracts cannot verify[^l2beat-da][^chainlink-da-l2]. Vitalik (Jan 2024): an Ethereum L2 using external DA is **a validium, not a genuine rollup**. (Confidence: fact, 3+ sources.)

## 7. Withdrawal / challenge windows

**Optimistic rollups impose a ~7-day challenge/withdrawal window** (Arbitrum's default is precisely ~6.4 days; OP Stack/Base ~7 days)[^ethorg-or][^kelvinfichter-7day][^l2beat-stages]. Process: burn on L2 → batch to L1 → user computes a Merkle inclusion proof → wait out the window → finalize on L1.

**Deposits (L1→L2) are fast and asymmetric to withdrawals.** A user deposits by calling a function on the rollup's L1 bridge/inbox contract (e.g., OP Stack's `OptimismPortal.depositTransaction`, Arbitrum's Delayed Inbox); the sequencer observes the L1 event and credits the user on L2 within minutes (no challenge window applies, because the deposit is an L1 transaction the rollup simply mirrors). Only the L2→L1 *withdrawal* of an optimistic rollup is slow, for the reason below. **Why 7 days specifically:** the window must exceed *the maximum time an attacker can censor all challengers* from getting a challenge included; the driving threat is **repeatable/forking censorship** (a 51%-style attack), so 7 days is a **conservative Schelling-point safety margin** leaving time for the social layer (detect within ~24h, coordinate a response within ~5d, play the challenge in <1d) — **not a tuned cryptographic constant**[^kelvinfichter-7day][^l2beat-stages].

**Fast-withdrawal / liquidity-provider bridges** front the user funds on L1 **immediately for a fee** and take over the pending canonical withdrawal, getting reimbursed ~7 days later (Hop, Across, etc.)[^ethorg-l2][^kelvinfichter-7day]. **They are NOT trustless** — they reintroduce LP solvency/counterparty risk or trust in an external messaging set (an LP *can* itself verify the withdrawal before fronting, bounding its own risk, but from the user's view it's a swap of rollup-trust for LP-trust). **Canonical/native bridges** are trustless by design (built by the rollup team, no new validator set) — only downside is the slow withdrawal; **third-party bridges** add a separate trust layer and are historically the most-exploited crypto infrastructure. `[VOLATILE]` **ZK contrast:** validity rollups need **no** challenge window — once the proof verifies on L1, withdrawals settle in ~30 min–few hours. L2BEAT historically required ≥7-day challenge periods for Stage 1 and in 2026 **lowered the Stage-1 minimum to ~5 days**[^l2beat-stages]. (Confidence: fact; the exact Stage-1 minimum is volatile.)

## 8. The L2BEAT Stage 0/1/2 framework

The framework — inspired by Vitalik's "proposed milestones for rollups taking off training wheels" — categorizes rollups by how much they rely on **"training wheels" (centralized overrides)**, keyed on *when a Security Council may override the trustless proof system*[^l2beat-stages][^vitalik-stages]:

- **Stage 0 — full training wheels.** Operators effectively run the rollup; a proof system may exist but the **Security Council can override it with a simple majority** (the proof system is "advisory"). Software must allow state reconstruction from L1 data.
- **Stage 1 — limited training wheels.** Governed by smart contracts with a **fully functional, permissionless-submission proof system**; the Security Council remains a bug backstop but can only override with a **high supermajority (≥75%)**, with a quorum-blocking subset outside the primary org.
- **Stage 2 — no training wheels.** Fully proof-governed; the proof system is **permissionless**, users get **≥30 days to exit unwanted upgrades**, and the Security Council can act **only on provable on-chain bugs** (e.g., two redundant proof systems disagreeing).

The **Security Council** is a multisig safeguard; the 2025 Stage-1 principle is that "the only way (other than bugs) to indefinitely block or push an invalid L2→L1 message is by compromising ≥75% of the Security Council"[^l2beat-stages]. Vitalik's "math of when stage 1/2 make sense" argues **Stage 0 is rarely justified — launch at least at Stage 1**[^vitalik-stages].

`[VOLATILE — as of 2026-06]`: **No general-purpose L2 has reached Stage 2** (the few Stage-2 entries are app-specific or shut down); the **bulk of secured value sits at Stage 1**, a minority at Stage 0[^l2beat-stages]. **Secondary trackers disagree on the exact Stage-1 roster** (one lists Arbitrum One / OP Mainnet / Base / Scroll / Starknet; another lists far fewer) — this reflects genuine volatility plus secondary-source error: **the authoritative source is L2BEAT's live dashboard; verify there before relying**[^l2beat-stages]. Stage 2 is hard mostly because teams retain Security-Council emergency powers for rapid bug remediation; broadly a multi-year (~2026–2027+) target. (Confidence: fact for the framework structure and "no Stage 2 / mostly Stage 1"; the exact roster is volatile.)

---

## References

[^ethorg-l2]: ethereum.org — Layer 2 / Scaling. https://ethereum.org/developers/docs/scaling/ — L2 definition, taxonomy, fast-withdrawal LP bridges, sidechains-not-L2. tier: docs.
[^vitalik-rollups]: Vitalik Buterin — An Incomplete Guide to Rollups. https://vitalik.eth.limo/general/2021/01/05/rollup.html — rollup definition, optimistic vs ZK, DA need. tier: blog (primary).
[^paradigm-or]: Paradigm (gakonst) — Almost everything about optimistic rollup / "Sidechains are not L2." https://www.paradigm.xyz/2021/01/almost-everything-you-need-to-know-about-optimistic-rollup — fraud-proof model, L2-security==L1-security, sidechain distinction. tier: blog.
[^eth-mag-l2]: Vitalik Buterin — rollup-centric roadmap (unilateral-exit framing). https://ethereum-magicians.org/t/a-rollup-centric-ethereum-roadmap/4698 — "withdraw even if everyone colludes." tier: forum (primary).
[^vitalik-l2types]: Vitalik Buterin — Different types of layer 2s. https://vitalik.eth.limo/general/2023/10/31/l2types.html — rollup/validium/optimium/disconnected taxonomy. tier: blog (primary).
[^ethresear-sidechain]: ethresear.ch — Understanding sidechains / DA policies. https://ethresear.ch/t/taking-a-closer-look-at-data-availability-policies-in-ethereum-rollups/17594 — sidechain-vs-rollup, "better-security sidechains reduce to optimistic rollup." tier: forum.
[^ethorg-or]: ethereum.org — Optimistic rollups. https://ethereum.org/developers/docs/scaling/optimistic-rollups/ — assume-valid, fraud proofs, 1-of-N watcher, 7-day window. tier: docs.
[^arbitrum-bold]: Arbitrum docs — Nitro / BoLD. https://docs.arbitrum.io/ — Geth core + WASM prover, interactive bisection, BoLD permissionless validation (live Feb 12 2025). tier: docs.
[^op-faultproofs]: Optimism docs/specs — Fault Proof System / Cannon. https://docs.optimism.io/ — fault proofs live June 10 2024, Cannon MIPS VM, OptimismPortal forced txs, multi-proof. tier: docs.
[^base-faultproofs]: Base — fault proofs / multiproofs. https://blog.base.org/ — Base fault proofs Oct 30 2024, Stage 1, TEE+ZK multiproof direction. tier: blog.
[^ethorg-zk]: ethereum.org — Zero-knowledge rollups. https://ethereum.org/developers/docs/scaling/zk-rollups/ — validity proof per batch, no challenge window, residual trust. tier: docs.
[^vitalik-zkevm]: Vitalik Buterin — The different types of ZK-EVMs. https://vitalik.eth.limo/general/2022/08/04/zkevm.html — Type 1–4 classification, "everything becomes Type 1 over time." tier: blog (primary).
[^starkware-misnomer]: StarkWare — "validity rollup" / ZK-is-a-misnomer. https://starkware.co/ — succinctness not privacy, validity-rollup terminology, projects. tier: blog.
[^chainscore-l2]: ChainScore — L2 trust, sequencers, ZK fast-finality reality, TVL split. https://chainscore.finance/ — residual trust, centralized-sequencer reality, withdrawal-time reality, optimistic TVL majority. tier: blog.
[^chainlink-seq]: Chainlink — L2 sequencers & uptime feeds. https://chain.link/education-hub/sequencer — sequencer role, soft confirmations, L2 Sequencer Uptime Feeds. tier: blog.
[^l2beat-da]: L2BEAT — Data availability / "validiums are not L2s". https://l2beat.com/scaling/data-availability — DA classification, off-chain-DA trust downgrade, Ethereum-DA guarantee. tier: tracker.
[^chainlink-da-l2]: Chainlink — Data availability layers. https://chain.link/article/data-availability-layers — DAC trust model, Celestia/EigenDA/Avail, "secured by Ethereum." tier: blog.
[^kelvinfichter-7day]: kelvinfichter — Why is the challenge period 7 days? https://kelvinfichter.com/ — 7-day = Schelling-point margin for the social layer, repeatable censorship. tier: blog.
[^l2beat-stages]: L2BEAT — Stages framework & risk model (live dashboard + forum). https://l2beat.com/scaling/stages — Stage 0/1/2 definitions, Security Council ≥75%, Sequencer-Failure row, stage distribution, 7d→5d Stage-1 update. tier: tracker.
[^vitalik-stages]: Vitalik Buterin — Proposed milestones for rollups / "the math of when stage 1 and stage 2 make sense." https://ethereum-magicians.org/t/proposed-milestones-for-rollups-taking-off-training-wheels/11571 — training-wheels framing, "launch at least at Stage 1." tier: forum (primary).
