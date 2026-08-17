<!-- Provenance: reference under the standalone `bitcoin-protocol-expert` skill. Created 2026-06-16 via /dr deep-research. Volatile claims stamped verified-as-of 2026-06-16. Peer deferral: ECDSA/HD-derivation math → blockchain-crypto-primitives; consensus theory → distributed-systems-consensus; the Ethereum side of the comparison → ethereum-protocol-expert. -->

# Running a Node, Wallets & UTXO-vs-Account Design

Practical operation of Bitcoin, plus the comparative-design payoff against account-model chains.

**verified-as-of: 2026-06-16** — volatile claims (chain size, Bitcoin Core version) stamped inline.

## Contents
- [Full nodes](#full-nodes)
- [Bitcoin Core & sync](#bitcoin-core--sync)
- [Storage: archival vs pruned](#storage-archival-vs-pruned)
- [Node vs miner; node vs SPV](#node-vs-miner-node-vs-spv)
- [Do you need a node?](#do-you-need-a-node)
- [Wallets manage keys](#wallets-manage-keys)
- [Wallet types](#wallet-types)
- [Seeds, HD, PSBT (name level)](#seeds-hd-psbt-name-level)
- [UTXO vs account model](#utxo-vs-account-model)
- [References](#references)

## Full nodes

A full node **independently and fully validates every transaction and block** against Bitcoin's consensus rules — trusting no third party (fact).[^1][^2][^6] Beyond validation and relay, most full nodes **store and serve the blockchain** and **serve lightweight clients**; if too few do so, light clients are pushed onto centralized services.[^1] Validation checks: every input spends a valid unspent UTXO; signatures satisfy locking scripts; outputs don't exceed inputs (no inflation); the subsidy doesn't exceed the schedule; all soft-fork rules (SegWit, Taproot, CLTV, CSV) hold.[^3][^6]

**"Don't trust, verify" / economic-node sovereignty** (fact). A full node **rejects an invalid block absolutely — "even if every other node thinks it is valid."** No 51% attack can force an invalid block onto a full validator; hashrate can only reorder or censor, not change the rules.[^3][^6] The Block Size War (2015–17) proved it: user-run economic nodes rejected the SegWit2x hard fork and miners backed down. To actually contribute economic weight you must **use the node as your wallet on hardware you control** — a node on a cloud VM you don't transact through adds little.[^2][^7][^8]

## Bitcoin Core & sync

Bitcoin Core is the reference implementation; running `bitcoind` (daemon) or `bitcoin-qt` (GUI) makes you a full node.[^6][^7] **verified-as-of 2026-06-16:** the latest major release is **Bitcoin Core 31.0 (April 19, 2026)** (GUI on Qt 6.8), with 30.x/29.x/28.x maintenance lines current; v29.0 (Apr 2025) migrated the build system to **CMake**.[^9][^10][^11]

- **Initial Block Download (IBD)** — a new node downloads and validates every block from genesis (Jan 3, 2009) to the tip; the slow part is signature validation across 900M+ inputs (hours to >1 day).[^12][^13]
- **assumevalid** (fact; exact % qualified) — Bitcoin Core ships a hash of a recent known-valid block; for its ancestors the node **skips script/signature verification** while still checking structure, PoW, inflation, and UTXO-set consistency. It is an optimization, **not a trust shortcut** (identical resulting UTXO set), cutting IBD roughly in half (sources range "half" to "~80%").[^13][^14]
- **AssumeUTXO** (fact) — loads a pre-built UTXO-set snapshot at a height (hash hardcoded), becoming usable in minutes while **validating full history from genesis in the background** (`loadtxoutset`/`dumptxoutset`/`getchainstates`); the redundant chainstate is removed once background validation reaches the snapshot.[^15][^16][^17]

## Storage: archival vs pruned

**Chain size (verified-as-of 2026-06-16):** real-time trackers converge on **hundreds of GB, roughly 685–745 GB** full chain (YCharts 745.11 GB on June 4, 2026; Blockchair 685.59 GB; bitinfocharts 491.50 GB is a different/older methodology). Growth ~50–80 GB/year. **Order of magnitude: ~600–750 GB archival.** (fact for order-of-magnitude; exact figure qualified — methodologies differ by tens of GB.)[^18][^19][^20]

**Pruning (`-prune`)** (since v0.11, 2015) deletes old raw block/undo data after validation; minimum **550 MiB** (keeps ≥288 blocks ≈ 2 days), real footprints ~**10–20 GB** vs ~600–750 GB archival.[^21][^22][^23] **Key point: a pruned node is still fully-validating** — it applies every consensus rule and retains the complete UTXO set, identical to archival.[^4][^22][^24] **Pruned limitations** (disconfirming): cannot serve historical blocks to IBD peers, run `-txindex`, fully rescan old wallets, or back a block explorer; enabling/disabling pruning or any reindex requires re-downloading the whole chain.[^21][^22][^25]

## Node vs miner; node vs SPV

**Validation vs block production are separate roles** (fact). A node *verifies* (and rejects non-compliant blocks); a miner *proposes* (assembles txs, competes via PoW). Miners can only reorder/remove txs and **cannot decide the rules or force invalid blocks** — a miner's block joins the chain only if nodes accept it ("nodes are referees, miners are participants"). Many miners also run nodes, but the roles are distinct.[^6][^26][^27] *(Mining mechanics → `references/blocks-mining-and-monetary-policy.md`; mining-as-consensus → `distributed-systems-consensus`.)*

**Node vs SPV/light client (the trust trade-off)** (fact). SPV wallets download only block headers (~60 MB) and verify a tx's *inclusion* via Merkle proofs + PoW — but **cannot verify whether the tx's contents are actually valid** (only full validators can), so they effectively trust the most-PoW chain and can be fed fabricated txs by dishonest full nodes. If most of the economy ran light clients, evil miners could change rules and light clients would follow — hence the argument that the economy should be backed by full nodes.[^2][^3][^6]

## Do you need a node?

**Disconfirming — no, most people use Bitcoin via exchanges/wallets without running a node** (fact). It is "optional but valuable"; the payoff is **privacy, self-verification (sovereignty), and supporting decentralization — not money**.[^28][^29][^30] Caveat: the value only materializes if you **use the node as your wallet on hardware you control**; merely leaving a node on a cloud VM is a common misunderstanding ("bandwidth isn't the problem; trust, security and privacy are").[^31][^8]

## Wallets manage keys

A wallet **does not store bitcoin; it stores the private keys** that prove ownership of coins recorded on-chain (the coins exist as UTXOs) (fact).[^32][^33][^34] **Custodial vs non-custodial (self-custody)** (fact): custodial (an exchange) holds your keys — you hold an IOU/claim, like a bank balance; self-custody means you hold the keys directly — no one can freeze/censor your funds, but lose the keys/seed and the coins are gone, with no recovery. The mantra: **"Not your keys, not your coins."**[^33][^34][^35]

**Disconfirming — is self-custody worth it / the risks** (fact). Self-custody removes counterparty risk but concentrates operational risk on the user; **key loss is the dominant cause of permanently lost Bitcoin**. Failure modes: lost/damaged seed or passphrase (irreversible), seed compromise/theft, physical coercion ("$5-wrench attack"), single-point-of-failure backups (paper burns, metal plates lost in rubble — CNBC's 2025 California-wildfire reporting), silent failures found only at recovery. There is **no universally best method**; common advice: start custodial, graduate to self-custody as holdings grow, and a documented recovery plan is not optional.[^37][^38][^39][^40]

## Wallet types

(fact)[^34][^36][^37][^41]
- **Software (desktop/mobile)** — convenient, online ("hot"), more exposed.
- **Hardware / cold (air-gapped)** — devices (Coldcard, Trezor, NGRAVE) that store keys offline and sign without exposing them; highest security; data via MicroSD/QR.
- **Watch-only (xpub)** — loaded with only an extended public key / output descriptor: generates addresses, shows balances, builds (but can't sign) transactions; the coordinator in PSBT flows.
- **Multisig (m-of-n, e.g. 2-of-3)** — needs m of n keys, eliminating any single device/location/person as a point of failure; Coldcard supports up to 15 co-signers; commonly coordinated by watch-only software like Sparrow.
- **Hot vs cold** — internet-connected (convenient) vs offline (secure).

## Seeds, HD, PSBT (name level)

Wallets are **HD (BIP-32)** — one seed derives a tree of keypairs, so backing up the seed backs up every address.[^45] **BIP-39** is the **mnemonic seed phrase** (commonly 12 or 24 words from a 2048-word list; entropy → words with a checksum → binary seed); **BIP-44** adds a standard derivation path.[^46][^47] *(The derivation math → `blockchain-crypto-primitives`.)* **PSBT (BIP-174)** is the interchange format for not-yet-fully-signed transactions plus signing metadata — designed for hardware wallets, multisig, and CoinJoin, letting air-gapped signers sign without the UTXO set; it defines roles (Creator/Updater/Signer/Combiner/Finalizer/Extractor). In a 2-of-3 hardware-multisig flow, a watch-only coordinator creates the PSBT, each device adds a partial signature, and a combiner merges them before broadcast.[^42][^50][^41] *(See `references/utxo-transactions-and-script.md` for the BIP family.)*

## UTXO vs account model

The comparative-design payoff (fact — 3+ sources incl. Stanford EE374, ethereum.org, Cardano docs).[^51][^52][^53][^55] *(The Ethereum side in depth → `ethereum-protocol-expert`.)*

| Property | **UTXO (Bitcoin)** | **Account model (Ethereum)** |
|---|---|---|
| State | a set of discrete unspent outputs; state only *extended* | a global mutable state tree: address → {balance, nonce, code, storage} |
| Tx references | specific prior outputs (inputs) it consumes | sender + nonce; reads/writes shared global state |
| Balance | no stored balance — summed client-side | a single mutable value |
| Replay defense | each UTXO spent once (implicit via outpoints); **no nonce** | **per-account nonce**, strict ordering |
| Parallel validation | **natural** — independent UTXOs | **limited** — same-account txs nonce-ordered; shared contract state serializes |
| Change | leftover returns as change outputs | balance decremented; no change |
| Privacy | better with a new address per tx | reused static addresses link easily |
| Programmability | deliberately limited scripting ("verification model") | Turing-complete contracts ("computational model") |
| Determinism | stateless/deterministic; predictable fees; can't fail mid-execution | stateful; tx can fail mid-execution and still cost gas; MEV/reordering exposure |

**Why Bitcoin chose UTXO** — auditability/simple verification (check each referenced UTXO is unspent → solves double-spend), parallelism, predictability/determinism (fully-determined effect and predictable fee; no mid-execution failures), minimized attack surface + address-rotation privacy. **What it gives up:** rich on-chain programmability — UTXO is poorly suited to stateful applications.[^57][^60][^61]

**Why account-model chains chose otherwise** — Turing-complete smart contracts at static, callable addresses are far simpler for stateful dApps; more intuitive. Trade-offs accepted: nonces for replay defense, sequential per-account processing limiting parallelism, contract-security/computational-error risk, and indeterministic execution (unpredictable fees, MEV).[^51][^57][^61]

**Disconfirming — UTXO is NOT "primitive"/strictly worse** (fact). A genuine engineering trade-off, not a hierarchy. The critique (Buterin/ConsenSys) calls UTXO "unnecessarily complicated," stateless, and notes real wallets get ~null parallelism gains by keeping a single change output.[^60a] The rebuttal: **Cardano's Extended UTXO (eUTXO)** (Alonzo, Sept 2021) keeps UTXO's advantages while adding expressiveness — validators see the *specific* transaction's inputs/outputs + context (not arbitrary global state), giving **"equivalent computational power similar to the account model while maintaining stronger security guarantees"**: determinism/no partial failures, precisely predictable fees, **no reentrancy** (UTXOs consumed atomically), and parallelism from locality.[^53][^62][^63][^64] Academic framing (HackMD/arXiv): UTXO and accounts are equivalent under a unified model — UTXO's strict access lists enable greater parallelism, and Ethereum-style contracts *can* be built on UTXO via covenants — a **trade-off, not a hierarchy**.[^56][^60d][^67] Calling UTXO "primitive" misreads an intentional design choice.

## Child sub-concepts (future research)
Bitcoin Script & the non-Turing-completeness/covenant debate; output descriptors & the modern descriptor-wallet architecture; eUTXO smart-contract design patterns & the Cardano concurrency problem; node-in-a-box stacks (Umbrel, Start9, RaspiBlitz; electrs/Fulcrum); stateless clients & state-bloat solutions (Utreexo accumulators vs Verkle trees); multisig & collaborative-custody operational design (Shamir vs multisig, inheritance planning).

## References

[^1]: https://bitcoin.org/en/full-node — docs: full-node definition, relay/serve, IBD behavior, use-as-wallet.
[^2]: https://wiki.21ideas.org/en/practice/running-a-node — wiki: sovereignty, Block Size War / economic nodes, SPV trust.
[^3]: https://bitcoin.org/en/bitcoin-core/features/validation — docs: validates without trusting miners, 21M limit, SPV fabricated-tx risk.
[^4]: https://bitcoin-institute.pages.dev/entries/design/2009-01-03-bitcoin-storage-design/ — blog: storage layout, ~650 GB archival, 50–80 GB/yr, archival-vs-pruned.
[^6]: https://en.bitcoin.it/wiki/Full_node — wiki: trustless validation, absolute rejection of invalid blocks, miners' limited power, light-node risk.
[^7]: https://en.bitcoin.it/wiki/Full_node — wiki: use a full node as wallet to add economic weight.
[^8]: https://en.bitcoin.it/wiki/Clearing_Up_Misconceptions_About_Full_Nodes — wiki: must use node as wallet on own hardware; bandwidth misconception.
[^9]: https://github.com/bitcoin/bitcoin/releases — docs: Bitcoin Core 31.0 latest; 30.2/29.3/28.4 lines (verified-as-of 2026-06-16).
[^10]: https://bitcoincore.org/en/2026/04/19/release-31.0/ — docs: Bitcoin Core 31.0 (April 19, 2026), Qt 6.8.
[^11]: https://bitcoincore.org/en/releases/29.0/ — docs: 29.0 (Apr 2025), Autotools→CMake.
[^12]: https://elementarybitcoin.org/chapters/21-node-optimizations.html — book: IBD steps, 900M+ inputs, assumevalid, AssumeUTXO two-phase.
[^13]: https://www.spark.money/glossary/assume-valid — blog: assumevalid mechanism, ~half IBD time, distinct from AssumeUTXO.
[^14]: https://www.spark.money/glossary/assume-valid — blog: assumevalid vs AssumeUTXO distinction.
[^15]: https://github.com/bitcoin/bitcoin/blob/master/doc/assumeutxo.md — docs (primary): assumeutxo fast bootstrap, loadtxoutset/dumptxoutset/getchainstates.
[^16]: https://bitcoin-bitcoin.mintlify.app/advanced/assumeutxo — docs (mirror): immediate usability + background genesis validation.
[^17]: https://gist.github.com/jdlcdl/7a853ab9ff0867b060b073f22432df3c — blog: watch-only via xpub/descriptor, PSBT create/sign/combine/finalize, air-gapped QR signers.
[^18]: https://ycharts.com/indicators/bitcoin_blockchain_size — data: 745.11 GB on June 4, 2026, +12.27% YoY.
[^19]: https://blockchair.com/bitcoin — explorer: blockchain size 685.59 GB.
[^20]: https://bitinfocharts.com/bitcoin/ — stats: v31.0 release, chain "database size" 491.50 GB (outlier methodology).
[^21]: https://bitcoincore.org/en/releases/0.11.0/ — release notes (primary): pruning since v0.11, `-prune=N`, 550 MB minimum, 288-block retention.
[^22]: https://www.spark.money/glossary/pruned-node — blog: pruned fully validates, 740+ GB vs 12–20 GB, limitations.
[^23]: https://bitcoin-bitcoin.mintlify.app/operations/pruning — docs (mirror): pruning disk table, 550 MiB minimum.
[^24]: https://docs.start9.com/bitcoin-guides/1.0.x/archival-vs-pruned.html — docs: both modes fully validate; pruned can't run Mempool explorer.
[^25]: https://blog.blockstream.com/education/nodes/should-i-run-a-bitcoin-full-node-or-a-pruned-node/ — blog: recommends archival unless storage-constrained.
[^26]: https://bitcoin.stackexchange.com/questions/59220/ — Q&A: miner creates blocks vs node validates/stores.
[^27]: https://crypto.com/en/crypto/learn/whats-bitcoin-node-how-to-run-one — blog: node verifies / miner proposes; miners don't decide rules.
[^28]: https://www.bitcoin.diy/learn/run-a-node — blog: "Do I need a node? No"; optional; privacy/security/decentralization.
[^29]: https://coinbureau.com/guides/how-to-run-a-bitcoin-node — blog: optional; payoff is sovereignty + reduced data leakage.
[^30]: https://fiatisfake.org/run-a-node/ — blog: don't need a node to hold/use BTC; needed for max privacy/sovereignty.
[^31]: https://en.bitcoin.it/wiki/Clearing_Up_Misconceptions_About_Full_Nodes — wiki: must use node as wallet; "trust, security, privacy" are the issues.
[^32]: https://developer.bitcoin.org/devguide/wallets.html — docs: wallets create keys to receive/spend; "collection of private keys"; watch-only.
[^33]: https://getuntaught.com/learn/what-is-a-bitcoin-wallet — blog: keys not coins; custodial vs self-custody; "not your keys, not your coins."
[^34]: https://www.bitcoin.diy/learn/bitcoin-wallets-explained — blog: keys not coins; hot vs cold; custodial vs non-custodial; counterparty risk.
[^35]: https://www.ledger.com/academy/topics/security/custodial-vs-non-custodial-wallets — vendor blog: custodial = claim against provider; non-custodial = sole signing authority.
[^36]: https://loveisbitcoin.com/not-your-keys-not-your-coins-why-do-we-say-that/ — blog: self-custody = key control; wallet types.
[^37]: https://www.spark.money/research/self-custodial-vs-custodial-wallets — blog: HD/BIP-32, BIP-39 12/24 words; key loss = dominant cause of lost BTC; recovery planning.
[^38]: https://www.aureobitcoin.com/en/blog/custody-risks-comparison — blog: third-party vs self-management risk; no universally best method; collaborative custody.
[^39]: https://www.cnbc.com/2025/04/06/bitcoin-self-custody-crypto-risks.html — news (CNBC): wildfire seed-phrase losses; single-key single-point-of-failure; multisig.
[^40]: https://custodystress.com/memos/self-custody-bitcoin-risks — blog: single-operator fragility, silent failures.
[^41]: https://coldcard.com/guides/advanced/coldcard-multisig-setup — vendor docs: 2-of-3 multisig, xpub export, Sparrow watch-only coordinator, PSBT signing.
[^42]: https://github.com/bitcoin/bitcoin/blob/master/doc/psbt.md — docs (primary): PSBT (BIP-174) since 0.17; hardware/multisig/CoinJoin; roles.
[^45]: https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki — BIP-32 (primary): HD wallets, tree of keypairs from one seed, change keychains.
[^46]: https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki — BIP-39 (primary): mnemonic generation, 128/256 bits → 12/24 words, checksum, PBKDF2 seed.
[^47]: https://arcsign.io/en/blog/key-derivation-bip39-44 — blog: BIP-32/39/44 relationship at name level.
[^50]: https://github.com/bitcoin/bips/blob/master/bip-0174.mediawiki — BIP-174 (primary): PSBT purpose, air-gapped signing without UTXO-set access.
[^51]: https://www.alchemy.com/docs/utxo-vs-account-models — docs: UTXO non-fungible/spent-once/change; accounts = balance; trade-off table.
[^52]: https://web.stanford.edu/class/ee374/lec_notes/lec10.pdf — lecture notes: formal UTXO vs accounts validation; nonce solves replay.
[^53]: https://developers.cardano.org/docs/learn/core-concepts/eutxo/ — docs (Cardano): UTXO stateless/deterministic; "dumb contracts"; eUTXO ≈ account-model power with stronger security.
[^55]: https://ethereum.org/en/developers/docs/accounts/ — docs (ethereum.org): EOA vs contract accounts, account fields, nonce prevents replay.
[^56]: https://hackmd.io/@adlerjohn/Hy9ydLNSL — technical blog: UTXO parallelism via strict access lists; UTXO≡accounts; contracts on UTXO via covenants.
[^57]: https://blockonomi.com/utxo-vs-account-based-transaction-models/ — blog: UTXO parallelism + pseudonymity; Ethereum chose accounts for Turing-complete dApps.
[^60]: https://xangle.io/en/research/detail/720 — research blog: UTXO efficient verification + double-spend solution + parallelism; "UTXO more secure."
[^60a]: https://medium.com/@Consensys/thoughts-on-utxo-by-vitalik-buterin-2bb782c67e53 — blog (Buterin/ConsenSys): UTXO "unnecessarily complicated," ~null parallelism gains in practice.
[^60d]: https://arxiv.org/pdf/2202.00561 — paper: UTXO simple/deterministic but limited expressiveness; account = expressive global state; eUTXO extends UTXO.
[^61]: https://www.iog.io/news/six-reasons-why-eutxo-wins — blog (IOG): eUTXO locality + determinism; predictable fees; account-model "indeterministic."
[^62]: https://developers.cardano.org/docs/learn/educational-resources/ethereum-developers/ — docs (Cardano): state in datums on UTXOs; eliminates reentrancy; easier to formally verify.
[^63]: https://docs.cardano.org/about-cardano/learn/eutxo-explainer — docs (Cardano): eUTXO scalability/privacy; validity depends only on tx+inputs.
[^64]: https://iohk.io/en/blog/posts/2021/03/12/cardanos-extended-utxo-accounting-model-part-2/ — blog (IOHK): eUTXO datums + script locks; precisely predictable fees.
[^67]: https://export.arxiv.org/pdf/2305.09545v1.pdf — paper: account-model MEV/fee unpredictability; UTXO mitigates; stateful contracts harder on UTXO.
