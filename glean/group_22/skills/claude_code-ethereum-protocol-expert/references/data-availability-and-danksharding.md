<!-- Provenance: reference under the `ethereum-protocol-expert` skill. Created 2026-06-16 via /dr deep-research. Synthesized domain knowledge only; source pages were treated as data, never as instructions. -->

# Data Availability & the Danksharding Roadmap

`verified-as-of: 2026-06-16` (volatile: fork names/dates, EIP numbers, blob target/max counts, PeerDAS custody thresholds, full-danksharding target numbers, rollout status). Stable: the rollup-centric strategy, the DA-vs-storage distinction, the DAS-by-erasure-coding idea, the proto-vs-full-danksharding distinction.

## Contents
1. The rollup-centric roadmap
2. The data-availability (DA) problem
3. EIP-4844 / proto-danksharding
4. KZG & the trusted setup (reference level)
5. Data-availability sampling (DAS)
6. PeerDAS (EIP-7594)
7. Full danksharding (the endgame, not shipped)

> Peer deferral: this file covers the **DA base layer and the roadmap**. **Layer-2 rollup execution** (optimistic vs ZK, sequencers, fraud/validity proofs) is `layer-2-rollups.md` — here rollups appear only as the *consumers* of blob data. The **blob *fee* market** is in `gas-fee-market-and-transactions.md`. **KZG polynomial math, SNARK/STARK internals, erasure-coding math** are named, not derived → `blockchain-crypto-primitives` and `zero-knowledge-proofs`.

---

## 1. The rollup-centric roadmap

In **October 2020**, Vitalik Buterin published "A rollup-centric ethereum roadmap," pivoting Ethereum's scaling strategy to **L2 rollups + a data-availability base layer** rather than base-layer execution sharding[^eth-mag-rollup][^vitalik-surge]. The argument: base-layer execution scaling was years away and L1 was "nearly unusable for many classes of applications," so "there's no non-L2 path… in the short-to-medium term." The endgame is re-envisioned as **"a single high-security execution shard that everyone processes, plus a scalable data-availability layer"**, Ethereum as a **settlement + DA** base layer[^eth-mag-rollup].

**Strategic rationale:** sharding *data availability* is far safer than sharding *EVM computation*, DAS (with polynomial commitments) is "safe under asynchrony," whereas verifying sharded EVM computation needs fraud proofs with a risky synchrony assumption[^eth-mag-rollup][^paradigm-das]. **"The Surge"** is the data-scaling chapter of the roadmap; Vitalik's stated goals (Oct 2024): 100,000+ TPS across L1+L2, preserve L1 decentralization, and maximize L2 interoperability ("Ethereum should feel like one ecosystem")[^vitalik-surge]. (Confidence: fact, the primary roadmap post + Vitalik + multiple.)

## 2. The data-availability (DA) problem

**DA = the guarantee that a block proposer actually *published* all the block's data so anyone can download and verify it**, it is about *publication at block time*, not long-term storage[^ethorg-da][^chainlink-da][^status-da]. **The key conceptual split, often confused:** DA ≠ data retrievability/storage. ethereum.org states it explicitly, DA is the assurance that nodes *could access and verify* the data; it does **not** follow that the data is "accessible forever." Long-term retrievability (serving history) can be handled by a small set of archive nodes or decentralized storage (e.g., the Portal Network); "the core Ethereum protocol is primarily concerned with data availability, not data retrievability"[^ethorg-da].

**Why rollups need it:** rollups execute off-chain and post compressed batches to L1. If a sequencer *withholds* the underlying data, no one can generate a fraud proof (optimistic) or reconstruct state to withdraw (any rollup), "a malicious sequencer could finalize invalid state changes." DA is the link that lets an L2 inherit L1 security[^chainlink-da][^status-da]. **A rollup needs DA during its challenge/proof window, not forever.** The original problem: rollups posted data as L1 **calldata**, which all nodes store *permanently* and which was expensive, historically **>90% of a rollup's cost was L1 data posting**[^ethorg-scaling][^vitalik-multidim]. (Confidence: fact, 3+ converging sources.)

## 3. EIP-4844 / proto-danksharding

`[VOLATILE — blob counts]`. **EIP-4844 ("Shard Blob Transactions," a.k.a. proto-danksharding)** went live in the **Cancun-Deneb ("Dencun") upgrade on March 13, 2024**[^ethorg-danksharding][^eip4844-da][^consensys-4844]. Key facts:

- **Blob-carrying transactions are Type 3** and carry `blob_versioned_hashes` + a `max_fee_per_blob_gas` bid. **The tx does NOT contain the blob data**, only the versioned-hash reference; blobs propagate separately as **sidecars** on the consensus layer (not the EVM-accessible execution layer)[^eip4844-da][^consensys-4844].
- **Blob size:** each blob = 4096 field elements × 32 bytes = **128 KiB** (≈125 KB usable, since field elements must be below the BLS modulus)[^eip4844-da][^vitalik-multidim].
- **A separate blob-gas fee market** (the blob base fee), fully independent of the execution gas market, see `gas-fee-market-and-transactions.md` for the fee mechanics[^eip4844-da].
- **KZG → versioned hashes → point-evaluation precompile:** each blob is committed with a KZG commitment; the EVM sees only a **versioned hash** (version byte `0x01` + a hash of the commitment), exposed via the **`BLOBHASH` opcode** (0x49). A **point-evaluation precompile at 0x0a** verifies a KZG proof that a blob evaluates to a claimed value at a point, this is how rollups prove blob contents inside the EVM[^eip4844-da][^consensys-4844].
- `[VOLATILE]` **Blob pruning / retention ≈ 18 days** (`MIN_EPOCHS_FOR_BLOB_SIDECARS_REQUESTS = 4096` epochs)[^ethorg-danksharding][^chainlink-blob]. This comfortably covers optimistic rollups' ~7-day challenge window, so their security model is unchanged. (Drafts once said "30 days"; the shipped value is ~18 days, use the shipped value.)
- **Why blobs are cheaper than calldata:** (1) blob data is **not accessible to the EVM** (no execution cost) and (2) it is **not stored permanently** (pruned after ~18 days). Realized effect: rollups became **~10–100× cheaper** (the 100× tail is mostly ZK rollups with better compression)[^chainlink-blob][^vitalik-multidim].
- `[VOLATILE]` Counts have been raised several times since Dencun (via EIP-7691 in Pectra, then **BPO forks** after Fusaka). The canonical blob target/max history table lives in **one place**: see `gas-fee-market-and-transactions.md` §7. The conceptual point here is that PeerDAS (§6) is what lets the count grow safely without raising per-node load.

(Confidence: fact, EIP + ethereum.org + Consensys + Chainlink.)

## 4. KZG & the trusted setup (reference level)

**KZG** (Kate-Zaverucha-Goldberg) is a polynomial/vector commitment scheme giving **constant-size commitments and constant-size proofs**, fast to verify (including inside SNARKs) and forward-compatible with DAS[^kzg-faq][^ethorg-danksharding]. A blob is committed with a single 48-byte commitment; a point-evaluation proof lets anyone cheaply verify the polynomial's value at a point without the full data. **Polynomial-commitment internals are deferred to `blockchain-crypto-primitives`.**

**The KZG Ceremony** (a "Powers of Tau" trusted setup) produced the Structured Reference String the commitments need. It ran ~8 months, concluded Jan 2024, and was the **largest MPC trusted setup of its kind: 141,416 contributions**[^kzg-ceremony][^kzg-faq]. **Trust model = 1-of-N:** the setup is secure as long as *at least one* participant discarded their secret; to break it, *every* contributor would have to collude[^kzg-faq][^vitalik-trustedsetup]. (Confidence: fact, EF + ceremony site + ethereum.org.)

## 5. Data-availability sampling (DAS)

**Core idea:** instead of downloading an entire block to confirm availability, each node **randomly samples a few small chunks**; if all samples return, the node is *probabilistically* confident the whole dataset is available, making DA verification scale sublinearly[^paradigm-das][^vitalik-das][^ethorg-danksharding].

**Erasure coding (Reed-Solomon) is the enabling trick**, it turns "≥50% available" into "100% recoverable": data is encoded so that **any 50% of the encoded chunks reconstruct all of it**[^vitalik-das][^paradigm-das][^ethorg-peerdas]. Checking full availability would require downloading everything; checking ≥50% needs only a few random samples. An adversary must withhold **>50%** to make data unrecoverable, and random sampling catches that with overwhelming probability (each sample roughly halves the attacker's escape odds). The **2D scheme** (the danksharding target) Reed-Solomon-extends in both dimensions, enabling cell-level sampling and stronger reconstruction. Ethereum uses **KZG commitments** (not plain Merkle roots) as the binding commitment, because they let you verify a coded chunk directly against the commitment to the *uncoded* data, "no room for invalid encoding"[^paradigm-das][^vitalik-das]. (The coding/commitment math is `blockchain-crypto-primitives`' domain.) (Confidence: fact, 3+ sources.)

## 6. PeerDAS (EIP-7594)

`[VOLATILE — live but fast-moving]`. **PeerDAS (Peer Data-Availability Sampling, EIP-7594) shipped in the Fusaka upgrade, live on mainnet Dec 3, 2025**[^fusaka-blog-da][^ethorg-peerdas][^eip7594]. It is the **first step toward full DAS** and uses **1D** erasure coding (not the full 2D scheme):

- Each blob is independently Reed-Solomon-extended (size doubles), and the matrix of extended blobs is split into **128 columns**; a column is the unit of custody/sampling, and **any 64 of the 128 columns** reconstruct the full data[^eip4844-peerdas][^teku-peerdas].
- **Column-based custody:** a node **deterministically derives its custody column set from its node ID** (and validator balance), subscribes only to those gossip subnets, custodies its columns, and **samples** others over the p2p network[^eip7594][^ethorg-peerdas].
- `[VOLATILE]` **Custody tiers:** a default full node custodies a small subset (ethereum.org: "at least 8 of 128 subnets"); **supernodes** (validators totaling ≥4096 ETH) custody all 128 columns and "heal" gaps. Exact thresholds shifted across devnets — treat as volatile[^ethorg-peerdas][^teku-peerdas].
- **Why this lets blob count grow safely:** with sampling, bandwidth/storage no longer scale linearly with blob count, so throughput can rise via **BPO forks** (EIP-7892) without raising per-node requirements[^optimism-peerdas][^ethorg-peerdas]. `[VOLATILE]` Post-Fusaka projections target blob max well past 9 (e.g., toward ~48–72) through 2026 — present as a *direction*, not a fixed number.

`[VOLATILE] ROLLOUT NOTE:` As of `2026-06-16` PeerDAS is **live (1D sampling)**; the **full 2D / FullDAS** work (validator custody on all clients, distributed blob publishing, higher blob counts) was still in testing on devnets[^optimism-peerdas]. (Confidence: fact that PeerDAS is live; specific thresholds/projections volatile.)

## 7. Full danksharding (the endgame, not shipped)

`[VOLATILE]`. **Full danksharding completes the data-scaling vision and is NOT yet shipped** — proto-danksharding (4844) and PeerDAS are the stepping stones, and it is "several years away"[^ethorg-danksharding][^vitalik-surge].

- **Target capacity — express in MB/slot, since sources differ on blob count.** Vitalik's authoritative endgame figure is **~16 MB per slot** (~1.33 MB/s), which with rollup compression yields **~58,000 TPS**[^vitalik-surge]. Blob-count framing varies and is volatile: ethereum.org says full danksharding expands blobs "to **64**," while Vitalik/a16z/ethresear.ch frame the 1D-sampling ceiling as **~128 target / 256 max ≈ 16–32 MB**[^ethorg-danksharding][^a16z-danksharding]. Resolution: 64 is a near-term/illustrative milestone; the deeper ceiling is ~256 blobs — anchor on the 16 MB/slot figure and preserve both blob counts.
- **Mechanism:** a **single proposer** selects all blobs in a block; blobs are extended with a **2D** Reed-Solomon + KZG scheme; **every node performs DAS** (cell-level sampling in the 2D scheme), so no node downloads the full block[^ethorg-danksharding][^a16z-danksharding]. **PBS is a prerequisite** — specialized builders do the expensive commitment/proof computation while the proposer just picks the most profitable valid header (preserving home-staker decentralization)[^ethorg-danksharding].
- **PeerDAS vs full danksharding:** PeerDAS = **1D** erasure coding, column custody, "the practical bridge"; full danksharding = **2D** coding across the whole matrix, cell-level sampling, stronger redundancy. Going from one to the other is largely "an efficient implementation of the 2D cryptography" plus a parameter bump[^ethorg-peerdas][^ethresear-fulldas].

**Misconceptions to carry (the corrections this skill must make)** — all fact, 3+ sources:
1. **"Danksharding is traditional sharding" — FALSE.** Shard chains were **removed** from the roadmap; danksharding scales *data* via sampling, keeping a **single proposer view** of one block — it does not split base-layer execution into many user-facing sub-chains[^ethorg-danksharding][^cube-danksharding].
2. **"Blobs are permanent storage" — FALSE.** Blobs are ephemeral (~18 days); the EVM keeps only the versioned hash. Long-term storage is for rollup operators/indexers/Portal Network, not the L1 protocol[^ethorg-da][^chainlink-blob].
3. **"Data availability = data storage" — FALSE** (see §2)[^ethorg-da][^status-da].
4. **"Proto-danksharding gave us DAS" — FALSE.** Under EIP-4844 every consensus node still downloads *all* blob data; 4844 only built the scaffolding (KZG commitments, sidecars, separate blob space) so DAS could be added later[^vitalik-blobs][^eip4844-da].
5. **"PeerDAS = full danksharding" — FALSE** (1D vs 2D; PeerDAS is only a step)[^ethorg-peerdas][^optimism-peerdas].
6. **Blobs are *native* L1 DA on the consensus layer** — not "calldata 2.0," and distinct from external DA layers (Celestia/EigenDA/Avail) a rollup may opt into instead (which is a weaker-trust choice — see `layer-2-rollups.md`)[^chainlink-da].

(Confidence: fact for the conceptual structure; the endgame numbers and rollout status are volatile.)

---

## References

[^eth-mag-rollup]: Vitalik Buterin — A rollup-centric ethereum roadmap. https://ethereum-magicians.org/t/a-rollup-centric-ethereum-roadmap/4698 — the 2020 strategic pivot. tier: forum (primary).
[^vitalik-surge]: Vitalik Buterin — Possible futures part 2: The Surge. https://vitalik.eth.limo/general/2024/10/17/futures2.html — 100k TPS goals, 16 MB/slot endgame, ~58k TPS, convergence story. tier: blog (primary).
[^paradigm-das]: Paradigm — Data Availability Sampling: From Basics to Open Problems. https://www.paradigm.xyz/2022/08/das — DAS, Reed-Solomon, KZG-vs-Merkle. tier: research blog.
[^ethorg-da]: ethereum.org — Data availability. https://ethereum.org/developers/docs/data-availability/ — DA vs retrievability, ~18-day window. tier: docs.
[^chainlink-da]: Chainlink — Data availability in smart contracts / DA layers. https://chain.link/article/data-availability-layers — DA definition, onchain vs external DA, "secured by Ethereum." tier: blog.
[^status-da]: Status — What is data availability. https://status.network/blog/what-is-data-availability-blockchain-layer-2 — DA vs storage, DAS, committee = weakest. tier: blog.
[^ethorg-scaling]: ethereum.org — Scaling. https://ethereum.org/roadmap/scaling/ — >90% of rollup cost is data, blobs temporary, two-stage plan. tier: docs.
[^vitalik-multidim]: Vitalik Buterin — Multidimensional gas pricing. https://vitalik.eth.limo/general/2024/05/09/multidim.html — why a separate blob lane, ~100× cheaper, ~125 KB usable. tier: blog (primary).
[^ethorg-danksharding]: ethereum.org — Danksharding. https://ethereum.org/roadmap/danksharding/ — not-traditional-sharding, 6→64 blobs, PBS prerequisite, KZG explainer. tier: docs.
[^eip4844-da]: EIP-4844 — Shard Blob Transactions. https://eips.ethereum.org/EIPS/eip-4844 — all constants (Type 3, 4096, 2¹⁷, precompile 0x0a, 4096-epoch retention), "downloaded by all consensus nodes." tier: eip.
[^consensys-4844]: Consensys — Dencun part 5: EIP-4844. https://consensys.io/blog/ethereum-evolved-dencun-upgrade-part-5-eip-4844 — Type-3 anatomy, two precompiles, 128 KB. tier: blog.
[^chainlink-blob]: Chainlink — EIP-4844 blob storage. https://chain.link/article/eip-4844-blob-storage — temporary blobs ~18 days, EVM-inaccessible, separate fee market. tier: blog.
[^fusaka-blog-da]: Ethereum Foundation — Fusaka mainnet announcement. https://blog.ethereum.org/2025/11/06/fusaka-mainnet-announcement — Fusaka live Dec 3 2025, BPO1 10/15, BPO2 14/21. tier: docs (primary).
[^kzg-faq]: ethereum/kzg-ceremony — FAQ. https://github.com/ethereum/kzg-ceremony/blob/main/FAQ.md — KZG = Kate-Zaverucha-Goldberg, constant-size, Powers of Tau, 1-of-N. tier: spec/docs.
[^kzg-ceremony]: Ethereum Foundation — Wrapping up the KZG Ceremony. https://blog.ethereum.org/2024/01/23/kzg-wrap — 141,416 contributions, largest of its kind. tier: docs (primary).
[^vitalik-trustedsetup]: Vitalik Buterin — How do trusted setups work? https://vitalik.eth.limo/general/2022/03/14/trustedsetup.html — 1-of-N, Zcash-6 contrast. tier: blog (primary).
[^vitalik-das]: Vitalik Buterin — DAS notes. https://hackmd.io/@vbuterin/das — 50%→100% trick, subnets, KZG-over-Merkle rationale. tier: forum (primary).
[^ethorg-peerdas]: ethereum.org — PeerDAS. https://ethereum.org/roadmap/fusaka/peerdas/ — 128 columns, 8/128 subnets, supernode ≥4096 ETH, 1D-vs-2D. tier: docs.
[^eip7594]: EIP-7594 — PeerDAS. https://eips.ethereum.org/EIPS/eip-7594 — 1D extension, columns, custody-by-node-ID. tier: eip.
[^eip4844-peerdas]: ethpandaops — Live custody monitoring (PeerDAS). https://ethpandaops.io/posts/live-custody-monitoring-peerdas/ — 128 columns @2KiB, 64-column reconstruction threshold. tier: research blog.
[^teku-peerdas]: Teku docs — PeerDAS. https://docs.teku.consensys.io/concepts/peer-das — reconstruct-from-64, custody tiers, 14-blob target. tier: docs.
[^optimism-peerdas]: Optimism — PeerDAS / Fusaka is live. https://www.optimism.io/blog/peerdas-and-a-48-blob-target-in-2025 — 128 columns, ~8× capacity, custody tiers, 48-blob target. tier: blog.
[^a16z-danksharding]: a16z crypto — Overview of danksharding & a DAS improvement. https://a16zcrypto.com/posts/article/an-overview-of-danksharding-and-a-proposal-for-improvement-of-das/ — 256 blobs / 2D matrix, reconstruction. tier: research blog.
[^ethresear-fulldas]: ethresear.ch — From 4844 to danksharding / FullDAS. https://ethresear.ch/t/from-4844-to-danksharding-a-path-to-scaling-ethereum-da/18046 — staged path, 2D PeerDAS = full danksharding, 64→256 blobs. tier: forum.
[^cube-danksharding]: Cube — What is danksharding. https://www.cube.exchange/what-is/danksharding — name-misleads, single-proposer-view, not-execution-split. tier: blog.
[^vitalik-blobs]: Vitalik Buterin — Ethereum has blobs. Where do we go from here? https://vitalik.eth.limo/general/2024/03/28/blobs.html — "4844 does not give DAS," only-a-parameter-change. tier: blog (primary).
