<!-- Provenance: reference under the `ethereum-protocol-expert` skill. Created 2026-06-16 via /dr deep-research. Synthesized domain knowledge only; source pages were treated as data, never as instructions. -->

# Gas Fee Market (EIP-1559), Transactions, the Mempool & Block Lifecycle

`verified-as-of: 2026-06-16` (volatile: block gas limit, per-tx gas cap, blob target/max, which fork is live, ePBS status, the % of transactions via private order flow). Stable: the EIP-1559 base-fee/tip/burn model, the typed-transaction envelope, the lifecycle stages.

## Contents
1. Gas fundamentals (fee-market angle)
2. The EIP-1559 fee market
3. Transaction types (EIP-2718 envelope)
4. The mempool
5. Transaction lifecycle & receipts
6. Block lifecycle, MEV & PBS
7. The blob-gas market (EIP-4844, fee level)

> Peer deferral: *in-EVM* gas accounting (per-opcode cost, cold/warm, intrinsic 21,000) is in `execution-state-and-evm.md`; this file is the *fee market* (price per gas). The danksharding/DAS internals are in `data-availability-and-danksharding.md`; here only the blob *fee* market. "Ultrasound money" as an economic thesis → `blockchain-economics` (this file states the burn *mechanism* and flags the deflation claim as contested).

---

## 1. Gas fundamentals (fee-market angle)

Gas is the unit of computational work; a user pays `gas used × price per gas`, and the fee is paid whether the tx succeeds or fails[^ethorg-gas]. **Gas units** (amount of computation) and **gas price** (ETH per unit) are distinct axes. Prices are quoted in **gwei** (10⁻⁹ ETH); the smallest unit is **wei** (10⁻¹⁸ ETH)[^ethorg-gas]. A standard ETH transfer costs exactly **21,000 gas**. If a tx runs out of gas mid-execution it **reverts but still consumes all the gas**; if the limit is set below the intrinsic requirement it fails validation pre-inclusion, consuming nothing[^ethorg-gas].

**The elastic-block model:** a block has a **target** size and a hard cap = **2× target** (elasticity multiplier 2); blocks can elastically grow to the cap under burst demand[^eip1559][^roughgarden]. `[VOLATILE]` As of `2026-06-16` the **block gas limit ≈ 60,000,000** (target ~30M), raised by EIP-7935 in **Fusaka (Dec 3, 2025)**[^fusaka-blog]. History: 15M→30M (London, Aug 2021), then validator-signalled increases to 36M, 45M, and 60M. Do **not** quote the legacy "30M" as current. Fusaka also added a **per-transaction gas cap of 16,777,216 (2²⁴)** as DoS protection under higher block limits[^ethorg-fusaka]. The limit is **validator-signalled** (clients vote it up/down). (Confidence: fact, multiple sources incl. the EF Fusaka announcement.)

## 2. The EIP-1559 fee market

Activated at the **London** hard fork (block 12,965,000, Aug 5, 2021)[^eip1559]. EIP-1559 splits the per-gas price into:

- **Base fee** — a protocol-computed reserve price a tx must meet to be eligible. It is computed from the **parent block only** (so it is predictable), and adjusts by **at most ±12.5% per block** toward a **50%-full (gas-target) equilibrium**: full block → +12.5%, empty → −12.5%, at target → unchanged[^eip1559][^roughgarden][^ethorg-gas]. The base fee is **burned** (permanently destroyed).
- **Priority fee / tip (`maxPriorityFeePerGas`)** — paid directly to the block proposer; since the base fee is burned, a tip is *de facto* required for inclusion[^eip1559][^dt-1559].
- **Max fee (`maxFeePerGas`)** — the total per-gas cap; the sender is refunded `maxFee − (baseFee + tip)`.

**Effective gas price** the sender pays = `baseFee + min(maxFee − baseFee, maxPriorityFee)` — the `effectiveGasPrice` field in the receipt[^execapis-receipt][^dt-1559].

**Why burn the base fee?** Per the EIP: (1) ensure only ETH pays for Ethereum execution; (2) remove the proposer's incentive to *manipulate* the fee to extract more from users; (3) prevent user–validator collusion[^eip1559][^paradigm-1559]. Contrast the **pre-1559 first-price auction**: users bid one `gasPrice`, miners took the highest, and each user paid their full bid → chronic overpayment and volatile, hard-to-estimate fees. EIP-1559 algorithmically prices the base-fee portion; the **tip remains a first-price auction**, and under acute congestion the system effectively falls back to a tip auction[^eip1559][^paradigm-1559][^dt-1559].

**Misconceptions to carry (these are the common errors):**
- **"EIP-1559 lowered gas fees" — FALSE.** It changed *how* fees are priced (predictability, less overpayment), not how much blockspace costs; fees track blockspace supply/demand. The real fee reductions came from **L2s and EIP-4844**, not 1559[^coindesk-1559][^dt-1559]. (Confidence: fact + an empirical study finding only a "small effect on gas fee levels.")
- **"EIP-1559 / the burn guarantees deflation / 'ultrasound money'" — CONTESTED.** The EIP itself states it can **no longer guarantee a fixed supply**. Post-Merge issuance dropped sharply, and ETH had net-negative-issuance periods, but net issuance went **positive again** after 2024 L2 migration reduced base-layer burn. Preserve the contradiction; the economic modeling is `blockchain-economics`' domain[^eip1559][^dlnews-uss]. (Confidence: contested — preserved.)

## 3. Transaction types (EIP-2718 envelope)

EIP-2718 defines a typed envelope: a transaction is `TransactionType ‖ TransactionPayload`, with `TransactionType` ∈ 0x00–0x7f (128 possible) and a type-specific payload[^eip2718]. Clients tell legacy from typed by the first byte (0xc0–0xfe = legacy RLP; 0x00–0x7f = typed). Five types exist as of `2026-06-16`[^eip2718][^ethorg-tx][^metamask-txtypes]:

| Type | Name / EIP | Fork | Adds |
| --- | --- | --- | --- |
| **0x00** | Legacy | genesis | single `gasPrice`; chain-id in `v` (EIP-155) |
| **0x01** | Access List (EIP-2930) | Berlin 2021 | `accessList` (warm pricing) + explicit `chainId`; still single `gasPrice` |
| **0x02** | Dynamic Fee (EIP-1559) | London 2021 | `maxPriorityFeePerGas` + `maxFeePerGas` — **the default type today** |
| **0x03** | Blob (EIP-4844) | Dencun Mar 2024 | `maxFeePerBlobGas` + `blobVersionedHashes`; blobs travel in a sidecar; `to` can't be null |
| **0x04** | Set Code (EIP-7702) | Pectra May 2025 | `authorization_list` — install a delegation indicator in an EOA's code |

EIP-7702 authorizations are signed over `keccak256(0x05 ‖ rlp([chain_id, address, nonce]))` and are persistent until revoked; `chain_id = 0` makes one valid on all chains (a cross-chain replay risk)[^eip7702-fm]. No Type 5 exists yet. (Confidence: fact, EIPs + ethereum.org + client docs.)

## 4. The mempool

Post-Merge there are **two p2p networks**: execution clients gossip *transactions* (devp2p `eth` wire protocol); consensus clients gossip *blocks*[^ethorg-networking]. Tx propagation is **hybrid**: full tx objects are pushed to a small subset of peers (~√n), and all others get a **hash announcement** and pull the body on demand[^devp2p-eth][^geth-deepwiki]. `[VOLATILE]` Large txs and **Type-3 (blob) txs are announcement-only** (never eager-pushed) to save bandwidth.

**Public vs private — not all transactions are public (a key correction):** a large and growing share of transactions route through **private RPCs / order-flow auctions** (Flashbots Protect, MEV-Share, etc.), bypassing the public mempool[^arxiv-privatepool][^flashbots-protect]. `[VOLATILE — magnitude]` one rigorous 2025 study estimated **~80% of transactions** route privately; treat the exact figure as an estimate (preserve as single-rigorous-source). Private mempools give **pre-trade privacy** (the tx is invisible until included), frontrunning/sandwich protection, and MEV/gas refunds; MEV-Share offers *programmable privacy* where searchers can only **backrun** shared txs[^flashbots-protect].

**Replacement (replace-by-fee, RBF):** same sender + same nonce replaces a pending tx only if the fee is bumped by a minimum percentage. `[VOLATILE]` go-ethereum default `PriceBump` = **10%** (both tip and fee cap must each rise ≥10%); some clients historically used 12.5% — preserve the ecosystem inconsistency[^geth-deepwiki]. **Pending vs queued:** executable txs (contiguous nonces) sit in *pending*; a tx with a nonce gap sits in *queued* (non-executable) until the gap fills. Rejected outright: bad signature, nonce-too-low, tip-above-fee-cap, underpriced, or insufficient balance[^devp2p-eth][^geth-deepwiki]. (Confidence: fact for the mechanism; the private-flow magnitude is the contested number.)

## 5. Transaction lifecycle & receipts

**Signing** = ECDSA over secp256k1; the address is the last 20 bytes of `keccak256(public key)`[^ethorg-tx][^morpheum-sign]. **EIP-155 chain-id replay protection** (legacy txs): post-155 the signed hash includes the chain-id and `v = chainId × 2 + 35 + recovery`, tying a signature to a specific chain[^eip155]. Typed txs (Type 1–4) carry `chainId` explicitly, so their `v` is just the y-parity (0/1)[^morpheum-sign].

**Lifecycle stages**[^ethorg-tx]: (1) tx hash generated → (2) broadcast to the mempool → (3) a proposer includes it in a block → (4) the block is upgraded **justified → finalized** over ~2 epochs (deep finality mechanics → `proof-of-stake-consensus.md`).

**Receipt fields** (canonical: the execution-apis schema)[^execapis-receipt]: `status` (1 success / 0 fail; post-Byzantium), `gasUsed` (this tx alone), `cumulativeGasUsed` (this tx + all preceding in the block), `logs` + `logsBloom` (a 2048-bit bloom filter for cheap log lookup), `effectiveGasPrice`, `contractAddress` (if a contract was created), and `blobGasUsed`/`blobGasPrice` (Type-3 only). External `eth_call` (view/pure) costs no gas. (Confidence: fact, 3+ sources.)

## 6. Block lifecycle, MEV & PBS

`[VOLATILE — landscape & ePBS status]`. Post-Merge, time is partitioned into **12-second slots**, one proposer per slot; the attestation deadline is ~4s into the slot[^paradigm-mevboost].

**Default (in-protocol) building:** the proposer's CL asks its EL (via the Engine API) to pull txs from the mempool and produce the **execution payload** (the transaction-execution half of the block), then wraps it in a beacon block and gossips it[^ethorg-networking].

**MEV (Maximal Extractable Value)** = value a block producer can capture by choosing/ordering/inserting txs (arbitrage, liquidations, sandwiching). **Searchers** submit **bundles** to **builders**[^flashbots-overview][^paradigm-mevboost].

**Out-of-protocol PBS via MEV-Boost (current production reality):** `[VOLATILE]` ~90% of blocks are produced via MEV-Boost. Roles: **searchers → builders** (assemble full payloads to maximize profit) → **relays** (trusted auctioneers that validate blocks, escrow them, and expose only a **blinded header** + bid) → **proposer**. The proposer **signs the highest-paying blinded header without seeing contents**, then the relay reveals the body[^flashbots-overview][^mevboost-readme][^paradigm-mevboost]. Why PBS: it lets home stakers capture MEV without running sophisticated builders, and keeps MEV from centralizing block *validation*[^ethorg-pbs].

**Enshrined PBS (ePBS / EIP-7732) — the direction, NOT live.** `[VOLATILE]` As of `2026-06-16`, ePBS is **scheduled for a future fork (the Glamsterdam direction), not on mainnet**[^ethorg-glam][^eip7732]. It would move the builder/proposer auction *into consensus*, eliminating trusted relays: the execution payload is removed from the beacon block body and replaced by a signed builder bid, and a Payload-Timeliness Committee attests whether the builder revealed on time. (Confidence: fact that MEV-Boost is current and ePBS is not live; the timeline is volatile.)

## 7. The blob-gas market (EIP-4844, fee level)

`[VOLATILE — blob counts change at BPO forks]`. Blob gas is a **separate, independent gas dimension** with its own EIP-1559-style market — Ethereum's first multidimensional gas pricing[^eip4844][^vitalik-multidim]. Each blob consumes exactly `GAS_PER_BLOB = 131,072 (2¹⁷)` blob gas, **not** counted against the normal tx gas limit[^eip4844].

The **blob base fee** is computed from the header's `excess_blob_gas` via an **exponential** rule, self-correcting with a max change ≈ ±12.5%/block (`MIN_BASE_FEE_PER_BLOB_GAS = 1 wei`)[^eip4844]. The blob base fee is also **burned**. Type-3 txs carry **`maxFeePerBlobGas`**; there is **no blob priority fee** — under burst demand, blob txs compete via their *regular* tx priority fee[^eip4844].

`[VOLATILE]` Blob target/max per block has changed three times — as of `2026-06-16` it is **target 14 / max 21** (BPO2, Jan 7, 2026):

| Upgrade | Date | Target | Max |
| --- | --- | --- | --- |
| Dencun (EIP-4844) | Mar 2024 | 3 | 6 |
| Pectra (EIP-7691) | May 2025 | 6 | 9 |
| Fusaka BPO1 | Dec 9, 2025 | 10 | 15 |
| **Fusaka BPO2 (current)** | **Jan 7, 2026** | **14** | **21** |

**BPO ("Blob-Parameter-Only") forks (EIP-7892)** are lightweight config-only forks that change *only* the blob target/max/fee-update-fraction, letting throughput scale between major upgrades after PeerDAS[^eip7892][^fusaka-blog]. Do **not** quote the Dencun-era "3/6" as current. The KZG/versioned-hash and danksharding mechanics behind blobs are in `data-availability-and-danksharding.md`. (Confidence: fact, EIP + EF announcement + the canonical blob-schedule EIP.)

---

## References

[^ethorg-gas]: ethereum.org — Gas and fees. https://ethereum.org/developers/docs/gas/ — gas units/gwei/limit, ±12.5%, 21k, fee math. tier: docs.
[^eip1559]: EIP-1559 — Fee market change. https://eips.ethereum.org/EIPS/eip-1559 — base fee/burn/tip/max-fee, supply-uncertainty admission. tier: eip.
[^roughgarden]: Roughgarden — Transaction Fee Mechanism Design for the Ethereum Blockchain (EIP-1559 analysis). https://timroughgarden.org/papers/eip1559.pdf — formal base-fee formula, target/2× cap. tier: paper.
[^paradigm-1559]: Paradigm — Analysis of EIP-1559. https://www.paradigm.xyz/2020/06/analysis-of-eip-1559 — elastic block size, tip stays first-price. tier: blog.
[^dt-1559]: decentralizedthoughts — EIP-1559. https://decentralizedthoughts.github.io/2022-03-10-eip1559/ — empirical: fees not lowered, effective-price formula. tier: blog.
[^execapis-receipt]: ethereum/execution-apis — receipt schema. https://github.com/ethereum/execution-apis/blob/main/src/schemas/receipt.yaml — receipt fields incl. effectiveGasPrice formula, blobGasUsed. tier: spec.
[^coindesk-1559]: CoinDesk — 4 common misperceptions about EIP-1559. https://www.coindesk.com/tech/2021/06/18/4-common-misperceptions-about-ethereums-eip-1559-upgrade — not-for-lowering-fees, not-stabilizing-monetary-policy. tier: blog.
[^dlnews-uss]: DL News — ETH "not ultrasound money anymore". https://www.dlnews.com/articles/defi/ethereum-scaling-means-eth-is-not-ultrasound-money-anymore/ — net issuance positive post-L2. tier: blog.
[^eip2718]: EIP-2718 — Typed Transaction Envelope. https://eips.ethereum.org/EIPS/eip-2718 — type-byte ranges, receipt format. tier: eip.
[^ethorg-tx]: ethereum.org — Transactions. https://ethereum.org/developers/docs/transactions/ — tx fields, lifecycle stages, type descriptions, receipts. tier: docs.
[^metamask-txtypes]: MetaMask — transaction types. https://docs.metamask.io/ — Type 0–4 mapping to hardforks. tier: docs.
[^eip7702-fm]: EIP-7702 — Set Code for EOAs. https://eips.ethereum.org/EIPS/eip-7702 — authorization_list, MAGIC 0x05, chain_id=0 replay. tier: eip.
[^ethorg-networking]: ethereum.org — Networking layer. https://ethereum.org/developers/docs/networking-layer/ — two p2p networks, EL builds the execution payload. tier: docs.
[^devp2p-eth]: ethereum/devp2p — eth wire protocol. https://github.com/ethereum/devp2p/blob/master/caps/eth.md — tx propagation (announce/pull), pool sync. tier: spec.
[^geth-deepwiki]: go-ethereum txpool (DeepWiki + source). https://github.com/ethereum/go-ethereum — validity errors, PriceBump 10%, pending/queued, slot limits. tier: spec.
[^arxiv-privatepool]: arXiv 2505.19708 — private transaction routing study. https://arxiv.org/html/2505.19708v1 — ~80% of txs via private RPCs (estimate). tier: paper.
[^flashbots-protect]: Flashbots docs — Protect / MEV-Share. https://docs.flashbots.net/ — private order flow, pre-trade privacy, backrun-only. tier: docs.
[^morpheum-sign]: Morpheum — Ethereum transaction signing. https://morpheum.dev/ — ECDSA, EIP-155 vs typed-tx y-parity. tier: blog.
[^eip155]: EIP-155 — Simple replay attack protection. https://eips.ethereum.org/EIPS/eip-155 — chain-id, v formula. tier: eip.
[^paradigm-mevboost]: Paradigm — MEV-Boost & Ethereum consensus. https://www.paradigm.xyz/2023/04/mev-boost-ethereum-consensus — 12s slots, t=4 attestation deadline, relay/builder/proposer roles. tier: blog.
[^flashbots-overview]: Flashbots — Auction overview / block proposal. https://docs.flashbots.net/flashbots-mev-boost/architecture-overview/block-proposal — MEV-Boost flow, blinded header. tier: docs.
[^mevboost-readme]: flashbots/mev-boost README. https://github.com/flashbots/mev-boost — sequence diagram, blinded-block submission. tier: spec.
[^ethorg-pbs]: ethereum.org — Proposer-builder separation. https://ethereum.org/roadmap/pbs/ — PBS rationale. tier: docs.
[^ethorg-glam]: ethereum.org — Glamsterdam. https://ethereum.org/roadmap/glamsterdam/ — ePBS scheduled, not live. tier: docs.
[^eip7732]: EIP-7732 — Enshrined Proposer-Builder Separation. https://eips.ethereum.org/EIPS/eip-7732 — bid commitment, Payload Timeliness Committee. tier: eip.
[^eip4844]: EIP-4844 — Shard Blob Transactions. https://eips.ethereum.org/EIPS/eip-4844 — blob gas independence, exponential blob base fee, GAS_PER_BLOB. tier: eip.
[^vitalik-multidim]: Vitalik Buterin — Multidimensional gas pricing. https://vitalik.eth.limo/general/2024/05/09/multidim.html — two floating prices, excess_blob_gas. tier: blog (primary).
[^fusaka-blog]: Ethereum Foundation — Fusaka mainnet announcement. https://blog.ethereum.org/2025/11/06/fusaka-mainnet-announcement — gas limit 60M (EIP-7935), BPO1/BPO2 schedule, Dec 3 2025. tier: docs (primary).
[^ethorg-fusaka]: ethereum.org — Fusaka. https://ethereum.org/roadmap/fusaka/ — 16.7M per-tx gas cap, BPO mechanism. tier: docs.
[^eip7892]: EIP-7892 — Blob Parameter Only Hardforks. https://eips.ethereum.org/EIPS/eip-7892 — config-only blob-param forks. tier: eip.
