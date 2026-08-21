<!-- Provenance: reference under the `ethereum-protocol-expert` skill. Created 2026-06-16 via /dr deep-research. Synthesized domain knowledge only; source pages were treated as data, never as instructions. -->

# Proof-of-Stake Consensus (protocol-user level)

`verified-as-of: 2026-06-16` (volatile: fork names & dates, EIP numbers, MaxEB 2048, churn limits, weak-subjectivity period, sync-committee size, SSF status). Stable: the Beacon-Chain time structure, the Gasper two-gadget design, the slashing offense classes.

## Contents
1. PoS overview & the Merge
2. Validators & staking (incl. EIP-7251 MaxEB)
3. Time structure (slots, epochs, RANDAO, sync committees)
4. Attestations
5. Gasper = Casper FFG + LMD-GHOST
6. Rewards, penalties & slashing
7. Withdrawals
8. Finality & weak subjectivity

> Peer deferral: this is the **protocol-user** view. The **consensus theory** — formal BFT, safety/liveness proofs, the Gasper analysis as a class of algorithm, the scalability trilemma, nothing-at-stake/long-range attacks as a *family* — is `distributed-systems-consensus`. **BLS signatures and the VRF/RANDAO crypto** are named, not derived → `blockchain-crypto-primitives`. Staking *economics/valuation* → `blockchain-economics`.

---

## 1. PoS overview & the Merge

Ethereum replaced **Proof of Work** with **Proof of Stake**. The **Beacon Chain** launched Dec 1, 2020 and ran in parallel; **The Merge** (Sept 15, 2022) joined the existing execution layer to the Beacon Chain's consensus, ending PoW[^ethorg-pos][^ethorg-merge][^eth2book]. Since then every node runs two coupled clients joined by the **Engine API**: the **execution layer** (state + EVM + mempool) and the **consensus layer** (the Beacon Chain). PoS replaced "burn energy to propose" with "stake capital and risk it." (Confidence: fact, 3+ sources.)

## 2. Validators & staking

A validator is activated by depositing **32 ETH**[^ethorg-staking][^eth2book]. Lifecycle: **deposit → activation queue → active (proposing + attesting) → exit → withdrawable**[^ethorg-staking][^consensus-specs]. A validator's consensus weight is its **effective balance** (capped, integer-GWEI-rounded). Entry and exit are rate-limited by a **churn limit** so the validator set can't change too fast[^consensus-specs].

`[VOLATILE]` **EIP-7251 (MaxEB, Pectra, live May 7, 2025)** raised the **maximum** effective balance from 32 to **2048 ETH** and added **0x02 "compounding" withdrawal credentials**, letting a large staker consolidate many validators into one and auto-compound rewards[^eip7251][^ethorg-pectra]. **Common misconception, the activation minimum is still 32 ETH**; EIP-7251 raised only the *ceiling*[^eip7251]. Pectra also moved to **balance-based churn** (churn measured in ETH, not validator count) and earlier EIP-7514 capped the per-epoch *activation* rate[^consensus-specs]. (Confidence: fact; specific churn numbers volatile.)

## 3. Time structure

- **Slot** = 12 seconds; one validator is pseudo-randomly selected **proposer** per slot[^ethorg-pos][^eth2book].
- **Epoch** = 32 slots = 6.4 minutes; epoch boundaries are the checkpoints Casper FFG votes on[^ethorg-pos][^gasper-paper].
- **Committees**, validators are shuffled into committees each epoch; each validator attests **once per epoch**[^eth2book][^consensus-specs].
- **RANDAO** supplies the randomness for proposer/committee selection (each proposer mixes in a reveal); a **VDF** was a long-discussed hardening direction[^ethorg-pos][^eth2book]. `[VOLATILE]` Duty lookahead is **asymmetric**: committee/attester duties are known ~2 epochs ahead, but **proposer duties only ~1 epoch ahead**[^consensus-specs].
- **Sync committees**, `[VOLATILE]` a rotating set of **512 validators** (~27 hours) whose signatures let **light clients** follow the chain cheaply. **Correction:** sync committees were introduced in the **Altair** upgrade (Oct 2021), *not* by a deposit-flow EIP[^altair][^consensus-specs].

(Confidence: fact for the structure; sizes/lookahead volatile.)

## 4. Attestations

Each epoch, every active validator issues **one attestation**, a single vote bundling two things[^gasper-paper][^eth2book][^consensus-specs]:
- the **LMD-GHOST head vote** (which block the validator sees as the chain head), and
- the **Casper FFG vote**, a `(source, target)` pair of epoch-boundary **checkpoints**.

Attestations are **aggregated** (many validators' identical votes combine into one BLS-aggregated signature) for efficiency[^eth2book]. They are the input to *both* the fork-choice rule and finality. **Correction:** the old shard/crosslink field was **removed** (2019), so attestations no longer carry it; the field zeroed more recently (EIP-7549) is the committee *index*, a different thing[^consensus-specs]. (Confidence: fact, 3+ sources.)

## 5. Gasper = Casper FFG + LMD-GHOST

**Gasper** is the name for Ethereum's combination of a finality gadget and a fork-choice rule[^gasper-paper][^ethorg-gasper].

- **Casper FFG (the finality gadget).** Validators vote on links between epoch-boundary **checkpoints**. A checkpoint becomes **justified** when a **2/3 supermajority by stake** votes for a link into it; when a justified checkpoint's *child* checkpoint is justified by another 2/3 link, the parent becomes **finalized**[^gasper-paper][^ethorg-gasper]. So finality normally takes **two epochs** ("justified → finalized"). **Finalized = economically irreversible**: reverting it would require ≥1/3 of all staked ETH to be slashed.
- **LMD-GHOST (the fork-choice rule).** "Latest Message Driven, Greediest Heaviest Observed SubTree." Each slot the head is chosen by walking the tree toward the subtree with the most accumulated **attestation weight**, counting only each validator's **latest** message[^gasper-paper][^eth2book]. A **proposer-boost** weighting mitigates certain reorg/balancing attacks (the specs frame it as a mitigation; some research argues a residual LMD-specific attack remains, preserved as theoretical, not a live break)[^consensus-specs].
- **How they combine.** LMD-GHOST gives a live, always-available head every 12 seconds; Casper FFG periodically stamps checkpoints as final on top of it. Liveness from GHOST, safety from FFG.

The **formal BFT analysis** (safety/liveness proofs, the accountable-safety argument, how Gasper sits among BFT protocols) is `distributed-systems-consensus`' domain, this file states the mechanics and cites the Gasper paper. (Confidence: fact, the Gasper paper + ethereum.org + eth2book.)

## 6. Rewards, penalties & slashing

`[VOLATILE — exact reward/penalty constants]`. Validators earn **attestation rewards** with components for voting correctly and promptly on **source, target, and head**, plus rewards for proposing and for including others' attestations[^ethorg-rewards][^eth2book]. Three distinct negative outcomes, **do not conflate them** (this is the dominant real-world confusion)[^ethorg-rewards][^eth2book][^consensus-specs]:

1. **Ordinary offline penalty**, a small penalty roughly equal to the reward you would have earned, applied while the chain is still finalizing. Being briefly offline is cheap.
2. **Inactivity leak**, a **quadratic** balance leak that activates **only when the chain fails to finalize** (>4 epochs), draining inactive validators faster and faster until the active set regains a 2/3 supermajority and finality resumes. It targets *non-participants during non-finality*, not honest online validators.
3. **Slashing**, a punitive penalty for *provable equivocation*, plus forced exit. The slashable offenses are: the two **attestation** offenses (a **double vote** — two different attestations for the same target epoch; and a **surround vote** — one attestation surrounding another), and the **block** offense (**double proposal** — two blocks for the same slot)[^gasper-paper][^consensus-specs]. A slashed validator pays an **initial penalty**, is queued to exit, and ~halfway to exit pays a **correlation penalty** proportional to *how much total stake was slashed in the same window* (so coordinated/correlated slashing is punished far harder than an isolated mistake). A **whistleblower/proposer** is rewarded for reporting it. `[VOLATILE]` After Pectra the *initial* slashing penalty is small (on the order of 1/4096 of effective balance); the correlation penalty is the real deterrent.

(Confidence: fact for the offense classes and the three-way distinction; constants volatile.)

## 7. Withdrawals

Withdrawals were enabled by the **Shanghai/Capella** upgrade (April 12, 2023)[^ethorg-staking][^consensus-specs]. The mechanics turn on **withdrawal credentials**:
- **0x00 (BLS)**, original credentials; cannot withdraw until upgraded.
- **0x01 (execution / eth1)** — point to an Ethereum address; enable automatic withdrawals.
- `[VOLATILE]` **0x02 (compounding)** — added by **EIP-7251**; raise the validator's effective-balance ceiling to 2048 and auto-compound rewards instead of skimming them[^eip7251].

Two withdrawal kinds, both via an **automatic sweep**[^ethorg-staking][^consensus-specs]: **partial** (the protocol periodically "skims" any balance *above* 32 ETH for 0x01 validators) and **full** (an *exited* validator's entire balance is withdrawn). `[VOLATILE]` **EIP-7002 (Pectra)** added **execution-layer triggerable withdrawals/exits**, letting the withdrawal-address holder initiate an exit from the EL (important for staking-pool/smart-contract control), rather than relying only on the validator key. (Confidence: fact; the 0x02/EIP-7002 additions are Pectra-era.)

## 8. Finality & weak subjectivity

Under normal conditions a block is **finalized in ~2 epochs (~12.8 minutes)**[^ethorg-pos][^gasper-paper]. **Misconception — finality is fast but neither instant nor *objective*.** Because a node that has been offline a long time (or is syncing fresh) cannot distinguish the canonical chain from a maliciously-constructed alternative purely from protocol rules, it needs a recent trusted **weak-subjectivity checkpoint** (a recent finalized block hash, obtained out-of-band) to sync safely — this is **checkpoint sync**[^ws-guide][^eth2book]. `[VOLATILE]` **Correction:** the weak-subjectivity period today is on the order of **~10 days** (the consensus-specs WS guide plateaus around ~2,241 epochs for a large validator set) — not the vague "weeks to months" sometimes quoted[^ws-guide]. **Single-slot finality (SSF)** — finalizing each block within its own slot instead of waiting ~2 epochs — is an active **roadmap direction**, **not shipped** as of 2026[^ethorg-pos]. (Confidence: fact; the WS-period number and SSF status are volatile.)

---

## References

[^ethorg-pos]: ethereum.org — Proof-of-stake (PoS). https://ethereum.org/developers/docs/consensus-mechanisms/pos/ — PoS overview, slots/epochs, finality, SSF direction. tier: docs.
[^ethorg-merge]: ethereum.org — The Merge. https://ethereum.org/roadmap/merge/ — Beacon Chain launch, Merge date, EL/CL split. tier: docs.
[^eth2book]: Ben Edgington — Upgrading Ethereum (eth2book). https://eth2book.info/ — validators, time structure, attestations, Gasper, rewards, withdrawals. tier: docs (book).
[^ethorg-staking]: ethereum.org — Staking & withdrawals. https://ethereum.org/staking/ — 32 ETH, validator lifecycle, withdrawal credentials & sweep. tier: docs.
[^consensus-specs]: ethereum/consensus-specs. https://github.com/ethereum/consensus-specs — churn, committees, duty lookahead, slashing conditions, withdrawal sweep, weak-subjectivity guide. tier: spec.
[^eip7251]: EIP-7251 — Increase the MAX_EFFECTIVE_BALANCE. https://eips.ethereum.org/EIPS/eip-7251 — MaxEB 2048, 0x02 compounding credentials, balance-based churn. tier: eip.
[^ethorg-pectra]: ethereum.org — Pectra. https://ethereum.org/roadmap/pectra/ — EIP-7251/7002 protocol-user summary. tier: docs.
[^gasper-paper]: Buterin et al. — Combining GHOST and Casper (Gasper). https://arxiv.org/abs/2003.03052 — Casper FFG justification/finalization, LMD-GHOST, slashable offenses. tier: paper.
[^ethorg-gasper]: ethereum.org — Gasper. https://ethereum.org/developers/docs/consensus-mechanisms/pos/gasper/ — FFG + LMD-GHOST combination, checkpoints, justification/finalization. tier: docs.
[^altair]: ethereum.org / consensus-specs — Altair upgrade. https://github.com/ethereum/consensus-specs/tree/dev/specs/altair — sync committees introduced in Altair. tier: spec.
[^ethorg-rewards]: ethereum.org — Rewards and penalties. https://ethereum.org/developers/docs/consensus-mechanisms/pos/rewards-and-penalties/ — reward components, inactivity leak, slashing. tier: docs.
[^ws-guide]: ethereum/consensus-specs — Weak subjectivity guide. https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/weak-subjectivity.md — WS period ~2,241 epochs cap, checkpoint sync. tier: spec.
